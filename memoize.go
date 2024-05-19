package memoize

import (
	"sync"

	"github.com/emad-elsaid/memoize/cache"
)

// Memoizer memoizes func(In) Out function
type Memoizer[In any, Out any] struct {
	cache cache.Cache[In, func() Out]
	Fun   func(In) Out
}

// Do calls the memoized function with input i and memoize the result and return it
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

// New creates a new memoizer wrapping func and returning the Memoizer Do function directly
func New[In, Out any](fun func(In) Out) func(In) Out {
	m := Memoizer[In, Out]{
		Fun: fun,
	}

	return m.Do
}
