package cache

// WithReadOnly wraps another cacher making it's Store method Noop. useful with
// combination of CacheWithFallback. you can supress writes to some levels of
// the caches for example a list of [memory, redis, database] will store values
// in all of them on write, wrapping `database` in WithReadOnly will stop it
// from being modified. So the CacheWithFallback will bring up the values from
// the database to faster cachers without modifying it. Also useful if your
// database is readonly.
type WithReadOnly[K, V any] struct {
	Loader[K, V]
}

// Store method in this case is a Noop
func (r *WithReadOnly[K, V]) Store(key K, value V) {}
