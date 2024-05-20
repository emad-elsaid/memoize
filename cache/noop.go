package cache

// Noop discards all writes and load always fails. useful for tests for example.
type Noop[K, V any] struct{}

func (_ Noop[K, V]) Load(_ K) (value V, loaded bool) { return }
func (_ Noop[K, V]) Store(_ K, _ V)                  {}
