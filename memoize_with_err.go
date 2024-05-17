package memoize

import "sync"

// MemoizerWithErr wraps a function func(In) (Out, error)
type MemoizerWithErr[In any, Out any] struct {
	cache Cache[In, func() (Out, error)]
	Fun   func(In) (Out, error)
}

// Do calls the memoized function with input i and memoize the result and return it
func (m *MemoizerWithErr[In, Out]) Do(i In) (Out, error) {
	once, ok := m.cache.Load(i)
	if !ok {
		once, _ = m.cache.LoadOrStore(i,
			sync.OnceValues(
				func() (Out, error) { return m.Fun(i) },
			),
		)
	}

	return once()
}

// NewWithErr wraps a function fun that returns a value and error in a MemoizerWithErr and returns its Do method
func NewWithErr[In any, Out any](fun func(In) (Out, error)) func(In) (Out, error) {
	m := MemoizerWithErr[In, Out]{
		Fun: fun,
	}

	return m.Do
}
