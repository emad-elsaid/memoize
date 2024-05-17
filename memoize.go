package memoize

import "sync"

type Memoizer[In any, Out any] struct {
	cache *Cache[In, func() Out]
	Fun   func(In) Out
}

func (m *Memoizer[In, Out]) Do(i In) Out {
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

func New[In any, Out any](fun func(In) Out) func(In) Out {
	m := Memoizer[In, Out]{
		cache: &Cache[In, func() Out]{},
		Fun:   fun,
	}

	return m.Do
}
