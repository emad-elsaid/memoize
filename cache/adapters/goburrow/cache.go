package goburrow

import "github.com/emad-elsaid/memoize/cache"

type MangoCache[K, V any] interface {
	GetIfPresent(K) (V, bool)
	Put(K, V)
}

func Mango[K, V, OK, OV any](c MangoCache[K, V]) cache.Cacher[OK, OV] {
	return &mangoCache[K, V, OK, OV]{c}
}

type mangoCache[K, V, OK, OV any] struct {
	MangoCache[K, V]
}

func (h *mangoCache[K, V, OK, OV]) Load(key OK) (value OV, loaded bool) {
	k, ok := any(key).(K)
	if !ok {
		return value, false
	}

	v, loaded := h.MangoCache.GetIfPresent(k)
	if !loaded {
		return value, false
	}

	value, loaded = any(v).(OV)

	return value, loaded
}

func (h *mangoCache[K, V, OK, OV]) Store(key OK, value OV) {
	k, ok := any(key).(K)
	if !ok {
		return
	}

	v, ok := any(value).(V)
	if !ok {
		return
	}

	h.MangoCache.Put(k, v)
}
