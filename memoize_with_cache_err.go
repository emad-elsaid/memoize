package memoize

import "sync"

type Pair[T any] struct {
	V   T
	err error
}
type MemoizerWithCacheErr[In any, Out any] struct {
	inFlight Cache[In, *sync.Mutex]

	Cache Cacher[In, Pair[Out]]
	Fun   func(In) (Out, error)
}

func (m *MemoizerWithCacheErr[In, Out]) Do(i In) (Out, error) {
	inFlight, _ := m.inFlight.LoadOrStore(i, new(sync.Mutex))
	inFlight.Lock()

	val, ok := m.Cache.Load(i)

	if !ok {
		v, err := m.Fun(i)
		val = Pair[Out]{V: v, err: err}
		m.Cache.Store(i, val)
	}

	inFlight.Unlock()
	m.inFlight.Delete(i)

	return val.V, val.err
}

func NewWithCacheErr[In any, Out any](c Cacher[In, Pair[Out]], fun func(In) (Out, error)) func(In) (Out, error) {
	m := MemoizerWithCacheErr[In, Out]{
		Cache: c,
		Fun:   fun,
	}

	return m.Do
}
