package dgraphio

import "github.com/emad-elsaid/memoize/cache"

type RistrettoCache[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V, int64) bool
}

func Ristretto[K comparable, V any](c RistrettoCache[K, V]) cache.Cacher[K, V] {
	return &ristretto[K, V]{c}
}

type ristretto[K comparable, V any] struct {
	RistrettoCache[K, V]
}

func (h *ristretto[K, V]) Load(key K) (value V, loaded bool) {
	return h.RistrettoCache.Get(key)
}

func (h *ristretto[K, V]) Store(key K, value V) {
	h.RistrettoCache.Set(key, value, 0)
}
