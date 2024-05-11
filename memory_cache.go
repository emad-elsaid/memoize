package memoize

import "sync"

type MemoryCache[K any, V any] struct{ sync.Map }

func (m *MemoryCache[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, loaded := m.Map.LoadOrStore(key, value)
	return a.(V), loaded
}
