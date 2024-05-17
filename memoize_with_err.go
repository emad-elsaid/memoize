package memoize

import "sync"

type MemoizerWithErr[In any, Out any] struct {
	Cache *Cache[In, func() (Out, error)]
	Fun   func(In) (Out, error)
}

func (m *MemoizerWithErr[In, Out]) Do(i In) (Out, error) {
	once, ok := m.Cache.Load(i)
	if !ok {
		once, _ = m.Cache.LoadOrStore(i,
			sync.OnceValues(
				func() (Out, error) { return m.Fun(i) },
			),
		)
	}

	return once()
}

func NewWithErr[In any, Out any](fun func(In) (Out, error)) func(In) (Out, error) {
	m := MemoizerWithErr[In, Out]{
		Cache: &Cache[In, func() (Out, error)]{},
		Fun:   fun,
	}

	return m.Do
}
