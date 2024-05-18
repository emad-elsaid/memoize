package memoize

import (
	"sync"
	"testing"
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

	mem := NewWithCache(&Cache[string, int]{}, inc)

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

	mem := NewWithCacheErr(&Cache[string, Pair[int]]{}, inc)

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
