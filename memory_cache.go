package memoize

import "sync"

type Cache[K any, V any] struct {
	sync.Map
	empty V
}

func (m *Cache[K, V]) Store(key K, value V) {
	m.Map.Store(key, value)
}

func (m *Cache[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, loaded := m.Map.LoadOrStore(key, value)
	return a.(V), loaded
}

func (m *Cache[K, V]) Load(key K) (value V, ok bool) {
	a, ok := m.Map.Load(key)
	if a == nil {
		return m.empty, ok
	}
	return a.(V), ok
}

func (m *Cache[K, V]) Delete(key K) {
	m.Map.Delete(key)
}
