// Package patrickmn includes cache adapter for https://github.com/patrickmn/go-cache package
package patrickmn


import (
	"github.com/emad-elsaid/memoize/cache"
)

// GoCacher cache interface (the part memoize need)
type GoCacher interface {
	Get(k string) (any, bool)
	SetDefault(k string, x any)
}

// GoCache creates a new cacher from the patrickmn/go-cache
// For example:
//
//	  c := cache.New(5*time.Minute, 10*time.Minute)
//	  memoize.NewWithCache(
//		   patrickmn.GoCache[int](c),
//		   func(s string) int { return len(s) },
//	  )
func GoCache[V any](c GoCacher) cache.Cacher[string, V] {
	return &goCache[V]{c}
}

type goCache[V any] struct {
	GoCacher
}

func (h *goCache[V]) Load(key string) (value V, loaded bool) {
	out, loaded := h.GoCacher.Get(key)
	val, ok := out.(V)

	if !ok {
		return value, false
	}

	return val, true
}

func (h *goCache[V]) Store(key string, value V) {
	h.GoCacher.SetDefault(key, value)
}
