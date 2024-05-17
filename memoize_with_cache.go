package memoize

import "sync"

type Cacher[K any, V any] interface {
	Store(key K, value V)
	Load(key K) (value V, ok bool)
}

type MemoizerWithCache[In, Out any] struct {
	inFlight Cache[In, *sync.Mutex]
	Cache    Cacher[In, Out]
	Fun      func(In) Out
}

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

func NewWithCache[In any, Out any](c Cacher[In, Out], fun func(In) Out) func(In) Out {
	m := MemoizerWithCache[In, Out]{
		Cache: c,
		Fun:   fun,
	}

	return m.Do
}
