package cache

// WriteOnlyCache wraps another cacher making it's Load method Noop. useful for
// systems that you would like to write to without reading from. For example you
// can have the last cache of your system write to Kafka. But doesn't attempt to
// read from it.
type WriteOnlyCache[K, V any] struct {
	Storer[K, V]
}

// Load always respond with false
func (w *WriteOnlyCache[K, V]) Load(key K) (value V, ok bool) { return value, false }
