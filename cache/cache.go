package cache

import "sync"

// Loader is an interface that can be used to lookup key value
type Loader[K, V any] interface {
	Load(key K) (value V, ok bool)
}

// Storer is an interfade that can be used to store key value
type Storer[K, V any] interface {
	Store(key K, value V)
}

// Cacher an interface needed to save and read data
type Cacher[K, V any] interface {
	Loader[K, V]
	Storer[K, V]
}

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

// CacheWithFallback On read it will return the first available value in the
// cachers then write to all the caches missing the value in the way. On write
// it will write to all caches.
type CacheWithFallback[K, V any] struct {
	Caches []Cacher[K, V]
}

// Load will ask every cacher for the key in order. keep track of all caches missing the key.
// When the key is found it writes the value to all the caches that was asked for the key and didn't have it
// effectively bringing the key to cachers at the top of the list.
func (c *CacheWithFallback[K, V]) Load(key K) (value V, ok bool) {
	for i := range c.Caches {
		if value, ok = c.Caches[i].Load(key); ok {
			// From the previous cache backward until we hit the head of the slice
			// Store the value of this key
			for i--; i >= 0; i-- {
				c.Caches[i].Store(key, value)
			}

			return
		}
	}

	return
}

// Store will store the key, value pair in all the cachers in order
func (c *CacheWithFallback[K, V]) Store(key K, value V) {
	for i := range c.Caches {
		c.Caches[i].Store(key, value)
	}
}

// ReadOnlyCache wraps another cacher making it's Store method Noop. useful with
// combination of CacheWithFallback. you can supress writes to some levels of
// the caches for example a list of [memory, redis, database] will store values
// in all of them on write, wrapping `database` in ReadOnlyCache will stop it
// from being modified. So the CacheWithFallback will bring up the values from
// the database to faster cachers without modifying it. Also useful if your
// database is readonly.
type ReadOnlyCache[K, V any] struct {
	Loader[K, V]
}

// Store method in this case is a Noop
func (r *ReadOnlyCache[K, V]) Store(key K, value V) {}

// WriteOnlyCache wraps another cacher making it's Load method Noop. useful for
// systems that you would like to write to without reading from. For example you
// can have the last cache of your system write to Kafka. But doesn't attempt to
// read from it.
type WriteOnlyCache[K, V any] struct {
	Storer[K, V]
}

// Load always respond with false
func (w *WriteOnlyCache[K, V]) Load(key K) (value V, ok bool) { return value, false }
