package memoize

import "sync"

type Cacher[K any, V any] interface {
	// LoadOrStore returns the existing value for the key if present. Otherwise,
	// it stores and returns the given value. The loaded result is true if the
	// value was loaded, false if stored.
	LoadOrStore(key K, value V) (actual V, loaded bool)
	Load(key K) (value V, ok bool)
}

type Memoizer[In any, Out any, F func(In) Out] struct {
	Cache Cacher[In, func() Out]
	Fun   F
}

func (m *Memoizer[In, Out, F]) Do(i In) Out {
	once, _ := m.Cache.LoadOrStore(i,
		sync.OnceValue(
			func() Out { return m.Fun(i) },
		),
	)

	return once()
}

func MemoryMemoizer[In any, Out any, F func(In) Out](fun F) F {
	m := Memoizer[In, Out, F]{
		Cache: &MemoryCache[In, func() Out]{},
		Fun:   fun,
	}

	return m.Do
}

type MemoizerWithErr[In any, Out any, F func(In) (Out, error)] struct {
	Cache Cacher[In, func() (Out, error)]
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

func MemoryMemoizerWithErr[In any, Out any, F func(In) (Out, error)](fun F) F {
	m := MemoizerWithErr[In, Out, F]{
		Cache: &MemoryCache[In, func() (Out, error)]{},
		Fun:   fun,
	}

	return m.Do
}
