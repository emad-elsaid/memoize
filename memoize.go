package memoize

import "sync"

type Cacher[K any, V any] interface {
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
	val, ok := m.Cache.Load(i)
	if !ok {
		once, _ := m.onces.LoadOrStore(i,
			sync.OnceValue(
				func() Out { return m.Fun(i) },
			),
		)
		val = once()
		m.Cache.Store(i, val)
	}

	return val
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
	val, ok := m.Cache.Load(i)
	if ok {
		return val.V, val.err
	}

	once, _ := m.onces.LoadOrStore(i,
		sync.OnceValues(
			func() (Out, error) {
				v, err := m.Fun(i)
				pair := Pair[Out]{V: v, err: err}
				m.Cache.Store(i, pair)
				m.onces.Delete(i)

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
