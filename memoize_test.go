package memoize

import (
	"sync"
	"testing"
)

func TestMemoizer(t *testing.T) {
	counters := map[string]int{}
	inc := func(k string) int {
		counters[k]++

		return counters[k]
	}

	mem := New(inc)

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
