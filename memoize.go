package memoize

import "sync"

type Cacher[K any, V any] interface {
	LoadOrStore(key K, value V) (actual V, loaded bool)
	Store(key K, value V)
	Load(key K) (value V, ok bool)
	Delete(key K)
}

type Memoizer[In any, Out any, F func(In) Out] struct {
	onces MemoryCache[In, func() Out]
	Cache Cacher[In, Out]
	Fun   F
}

func (m *Memoizer[In, Out, F]) Do(i In) Out {
	if val, ok := m.Cache.Load(i); ok {
		return val
	}

	once, _ := m.onces.LoadOrStore(i,
		sync.OnceValue(
			func() Out {
				val := m.Fun(i)
				_, loaded := m.Cache.LoadOrStore(i, val)
				// Checking loaded to make sur ethe order of Store/Delete is correct
				// otherwise the compiler reorder them in reverse
				if !loaded {
					m.onces.Delete(i)
				}

				return val
			},
		),
	)

	return once()
}

func New[In any, Out any, F func(In) Out](fun F) F {
	m := Memoizer[In, Out, F]{
		Cache: &MemoryCache[In, Out]{},
		Fun:   fun,
	}

	return m.Do
}

type Pair[T any] struct {
	V   T
	err error
}
type MemoizerWithErr[In any, Out any, F func(In) (Out, error)] struct {
	onces MemoryCache[In, func() (Out, error)]
	Cache Cacher[In, Pair[Out]]
	Fun   F
}

func (m *MemoizerWithErr[In, Out, F]) Do(i In) (Out, error) {
	if val, ok := m.Cache.Load(i); ok {
		return val.V, val.err
	}

	once, _ := m.onces.LoadOrStore(i,
		sync.OnceValues(
			func() (Out, error) {
				v, err := m.Fun(i)
				pair := Pair[Out]{V: v, err: err}
				_, loaded := m.Cache.LoadOrStore(i, pair)
				if !loaded {
					m.onces.Delete(i)
				}

				return v, err
			},
		),
	)

	return once()
}

func NewWithErr[In any, Out any, F func(In) (Out, error)](fun F) F {
	m := MemoizerWithErr[In, Out, F]{
		Cache: &MemoryCache[In, Pair[Out]]{},
		Fun:   fun,
	}

	return m.Do
}
