package cache

// Loader is an interface that can be used to lookup key value
type Loader[K, V any] interface {
	Load(key K) (value V, ok bool)
}

// Storer is an interface that can be used to store key value
type Storer[K, V any] interface {
	Store(key K, value V)
}

// Cacher an interface needed to save and read data
type Cacher[K, V any] interface {
	Loader[K, V]
	Storer[K, V]
}
