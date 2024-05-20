package cache

// WithFallback On read it will return the first available value in the
// cachers then write to all the caches missing the value in the way.
//
// On store it will store to all caches.
func WithFallback[K, V any](c ...Cacher[K, V]) Cacher[K, V] {
	return &withFallback[K, V]{
		Caches: c,
	}
}

type withFallback[K, V any] struct {
	Caches []Cacher[K, V]
}

// Load will ask every cacher for the key in order. keep track of all caches missing the key.
// When the key is found it writes the value to all the caches that was asked for the key and didn't have it
// effectively bringing the key to cachers at the top of the list.
func (c *withFallback[K, V]) Load(key K) (value V, ok bool) {
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
func (c *withFallback[K, V]) Store(key K, value V) {
	for i := range c.Caches {
		c.Caches[i].Store(key, value)
	}
}
