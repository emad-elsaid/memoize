package memoize

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoizer(t *testing.T) {
	counters := map[string]int{}
	inc := func(k string) (int, error) {
		counters[k]++

		return 0, nil
	}

	mem := MemoryMemoizerWithErr(inc)

	concurrency := 100
	var wg sync.WaitGroup
	wg.Add(concurrency)

	routine := func() {
		mem("key1")
		mem("key2")
		mem("key3")
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

	assert.Equal(t, expected, counters)
}
