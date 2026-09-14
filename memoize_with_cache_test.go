package memoize

import (
	"sync"
	"testing"
	"time"

	"github.com/emad-elsaid/memoize/cache"
)

func TestMemoizerWithCache(t *testing.T) {
	counters := map[string]int{}
	var l sync.Mutex
	inc := func(k string) int {
		l.Lock()
		defer l.Unlock()
		counters[k]++

		return counters[k]
	}

	mem := NewWithCache(cache.New[string, int](), inc)

	concurrency := 100
	var wg sync.WaitGroup
	wg.Add(concurrency)

	routine := func() {
		r := mem("key1")
		assertEqual(t, 1, r)

		r = mem("key2")
		assertEqual(t, 1, r)

		r = mem("key3")
		assertEqual(t, 1, r)
		wg.Done()
	}

	for range concurrency {
		go routine()
	}

	wg.Wait()

	expected := map[string]int{
		"key1": 1,
		"key2": 1,
		"key3": 1,
	}

	assertEqualMaps(t, expected, counters)
}

func TestMemoizerWithCacheErr(t *testing.T) {
	counters := map[string]int{}
	var l sync.Mutex
	inc := func(k string) (int, error) {
		l.Lock()
		defer l.Unlock()

		counters[k]++

		return counters[k], nil
	}

	mem := NewWithCacheErr(cache.New[string, Pair[int]](), inc)

	concurrency := 100
	var wg sync.WaitGroup
	wg.Add(concurrency)

	routine := func() {
		r, _ := mem("key1")
		assertEqual(t, 1, r)

		r, _ = mem("key2")
		assertEqual(t, 1, r)

		r, _ = mem("key3")
		assertEqual(t, 1, r)
		wg.Done()
	}

	for range concurrency {
		go routine()
	}

	wg.Wait()

	expected := map[string]int{
		"key1": 1,
		"key2": 1,
		"key3": 1,
	}

	assertEqualMaps(t, expected, counters)
}

// A panic inside Fun must not leave the per-key in-flight mutex locked forever.
// Without the fix, every later call for the same key hangs on Lock().
func TestMemoizerWithCachePanicDoesNotPoisonKey(t *testing.T) {
	m := &MemoizerWithCache[string, int]{Cache: cache.New[string, int]()}

	m.Fun = func(k string) int { panic("boom") }
	func() {
		defer func() { _ = recover() }()
		m.Do("k")
	}()

	// Replace Fun with a working one and verify the key is usable again.
	m.Fun = func(k string) int { return 42 }
	done := make(chan int, 1)
	go func() { done <- m.Do("k") }()

	select {
	case v := <-done:
		assertEqual(t, 42, v)
	case <-time.After(1 * time.Second):
		t.Fatal("key poisoned: Do hangs forever after a panic in Fun")
	}
}

// Same guarantee for the error-returning variant.
func TestMemoizerWithCacheErrPanicDoesNotPoisonKey(t *testing.T) {
	m := &MemoizerWithCacheErr[string, int]{Cache: cache.New[string, Pair[int]]()}

	m.Fun = func(k string) (int, error) { panic("boom") }
	func() {
		defer func() { _ = recover() }()
		m.Do("k")
	}()

	m.Fun = func(k string) (int, error) { return 42, nil }
	done := make(chan int, 1)
	go func() {
		v, _ := m.Do("k")
		done <- v
	}()

	select {
	case v := <-done:
		assertEqual(t, 42, v)
	case <-time.After(1 * time.Second):
		t.Fatal("key poisoned: Do hangs forever after a panic in Fun")
	}
}

// Single-flight: Fun must never execute more than once concurrently for the
// same key, even when the underlying cache can miss right after a Store
// (Noop drops every write, forcing every critical section to recompute).
//
// Regression test for the delete-after-unlock race: with the old code, the
// in-flight mutex was removed from the map while callers were still queued on
// it, so callers arriving after the delete minted a fresh mutex and ran Fun
// concurrently with the still-draining queued callers (maxConcurrent = 2).
func TestMemoizerWithCacheSingleFlight(t *testing.T) {
	var mu sync.Mutex
	concurrent, maxConcurrent := 0, 0

	mem := NewWithCache(cache.Noop[string, int]{}, func(k string) int {
		mu.Lock()
		concurrent++
		if concurrent > maxConcurrent {
			maxConcurrent = concurrent
		}
		mu.Unlock()

		time.Sleep(15 * time.Millisecond)

		mu.Lock()
		concurrent--
		mu.Unlock()
		return 1
	})

	var wg sync.WaitGroup

	// Wave 1: callers that queue on the in-flight mutex while the first
	// computation (spawned a moment earlier) is still running.
	time.Sleep(2 * time.Millisecond)
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); mem("k") }()
	}

	// Wave 2: callers arriving after the first computation finished. With the
	// buggy delete-after-unlock these bypass the old mutex and execute Fun
	// concurrently with the still-draining wave 1 callers.
	time.Sleep(40 * time.Millisecond)
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); mem("k") }()
	}

	wg.Wait()

	if maxConcurrent != 1 {
		t.Fatalf("single-flight violated: Fun ran %d times concurrently for the same key (want 1)", maxConcurrent)
	}
}

// Same guarantee for the error-returning variant.
func TestMemoizerWithCacheErrSingleFlight(t *testing.T) {
	var mu sync.Mutex
	concurrent, maxConcurrent := 0, 0

	mem := NewWithCacheErr(cache.Noop[string, Pair[int]]{}, func(k string) (int, error) {
		mu.Lock()
		concurrent++
		if concurrent > maxConcurrent {
			maxConcurrent = concurrent
		}
		mu.Unlock()

		time.Sleep(15 * time.Millisecond)

		mu.Lock()
		concurrent--
		mu.Unlock()
		return 1, nil
	})

	var wg sync.WaitGroup

	time.Sleep(2 * time.Millisecond)
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); mem("k") }()
	}

	time.Sleep(40 * time.Millisecond)
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); mem("k") }()
	}

	wg.Wait()

	if maxConcurrent != 1 {
		t.Fatalf("single-flight violated: Fun ran %d times concurrently for the same key (want 1)", maxConcurrent)
	}
}
