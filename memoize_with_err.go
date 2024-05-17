package memoize

import "sync"

type MemoizerWithErr[In any, Out any, F func(In) (Out, error)] struct {
	Cache *Cache[In, func() (Out, error)]
	Fun   F
}

func (m *MemoizerWithErr[In, Out, F]) Do(i In) (Out, error) {
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

func NewWithErr[In any, Out any, F func(In) (Out, error)](fun F) F {
	m := MemoizerWithErr[In, Out, F]{
		Cache: &Cache[In, func() (Out, error)]{},
		Fun:   fun,
	}

	return m.Do
}
