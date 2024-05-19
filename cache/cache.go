package cache

import "sync"

// Cache wraps sync.Map to use generic types of key and values instead of any type. It inherits the benefits of sync.Map. its methods can be used concurrently
type Cache[K any, V any] struct {
	sync.Map
	empty V
}

// Store See sync.Map#Store
func (m *Cache[K, V]) Store(key K, value V) {
	m.Map.Store(key, value)
}

// LoadOrStore See sync.Map#LoadOrStore
func (m *Cache[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, loaded := m.Map.LoadOrStore(key, value)
	return a.(V), loaded
}

// Load See sync.Map#Load
func (m *Cache[K, V]) Load(key K) (value V, ok bool) {
	a, ok := m.Map.Load(key)
	if a == nil {
		return m.empty, ok
	}
	return a.(V), ok
}

// Delete See sync.Map#Delete
func (m *Cache[K, V]) Delete(key K) {
	m.Map.Delete(key)
}
