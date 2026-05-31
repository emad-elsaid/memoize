// Package coocood includes cache adapter for https://github.com/coocood/freecache package
package coocood

import (
	"encoding/json"

	"github.com/emad-elsaid/memoize/cache"
)

// FreeCacher cache interface (the part memoize needs)
type FreeCacher interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte, expireSeconds int) error
}

// FreeCache creates a new cacher from the coocood/freecache
// As freecache supports only byte keys and values, the returned cacher uses
// string keys and JSON encoding/decoding for values.
// The expireSeconds parameter sets the TTL for all cached items.
//
// For example:
//
//	cacheSize := 100 * 1024 * 1024 // 100MB
//	fc := freecache.NewCache(cacheSize)
//	memoize.NewWithCache(
//	    coocood.FreeCache[int](fc, 3600), // 1 hour TTL
//	    func(s string) int { return len(s) },
//	)
func FreeCache[V any](c FreeCacher, expireSeconds int) cache.Cacher[string, V] {
	return &freeCache[V]{c, expireSeconds}
}

type freeCache[V any] struct {
	FreeCacher
	expireSeconds int
}

func (f *freeCache[V]) Load(key string) (value V, loaded bool) {
	data, err := f.FreeCacher.Get([]byte(key))
	if err != nil {
		return value, false
	}

	var val V
	if err := json.Unmarshal(data, &val); err != nil {
		return value, false
	}

	return val, true
}

func (f *freeCache[V]) Store(key string, value V) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	f.FreeCacher.Set([]byte(key), data, f.expireSeconds)
}
