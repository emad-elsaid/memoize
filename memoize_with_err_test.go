package memoize

import (
	"sync"
	"testing"
)

func TestMemoizerWithErr(t *testing.T) {
	counters := map[string]int{}
	inc := func(k string) (int, error) {
		counters[k]++

		return counters[k], nil
	}

	mem := NewWithErr(inc)

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
