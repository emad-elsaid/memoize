package memoize

import "sync"

type Memoizer[In any, Out any, F func(In) Out] struct {
	cache *Cache[In, func() Out]
	Fun   F
}

func (m *Memoizer[In, Out, F]) Do(i In) Out {
	once, ok := m.cache.Load(i)
	if !ok {
		once, _ = m.cache.LoadOrStore(i,
			sync.OnceValue(
				func() Out { return m.Fun(i) },
			),
		)
	}

	return once()
}

func New[In any, Out any, F func(In) Out](fun F) F {
	m := Memoizer[In, Out, F]{
		cache: &Cache[In, func() Out]{},
		Fun:   fun,
	}

	return m.Do
}
