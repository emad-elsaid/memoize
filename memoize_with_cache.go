package memoize

import (
	"sync"

	"github.com/emad-elsaid/memoize/cache"
)

type MemoizerWithCache[In, Out any] struct {
	inFlight cache.Cache[In, *sync.Mutex]
	Cache    cache.Cacher[In, Out]
	Fun      func(In) Out
}

// Do calls the memoized function with input i and memoize the result and return it
func (m *MemoizerWithCache[In, Out]) Do(i In) Out {
	inFlight, _ := m.inFlight.LoadOrStore(i, new(sync.Mutex))
	inFlight.Lock()

	val, ok := m.Cache.Load(i)

	if !ok {
		val = m.Fun(i)
		m.Cache.Store(i, val)
	}

	inFlight.Unlock()
	m.inFlight.Delete(i)

	return val
}

// NewWithCache takes a cacher and a function to memoize and returns, creates a MemoizerWithCache for it that uses the cacher c and returns its Do method
func NewWithCache[In any, Out any](c cache.Cacher[In, Out], fun func(In) Out) func(In) Out {
	m := MemoizerWithCache[In, Out]{
		Cache: c,
		Fun:   fun,
	}

	return m.Do
}
