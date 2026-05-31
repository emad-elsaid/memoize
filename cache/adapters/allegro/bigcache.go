// Package allegro includes cache adapter for https://github.com/allegro/bigcache package
package allegro

import (
	"encoding/json"

	"github.com/emad-elsaid/memoize/cache"
)

// BigCacher cache interface (the part memoize needs)
type BigCacher interface {
	Get(key string) ([]byte, error)
	Set(key string, entry []byte) error
}

// BigCache creates a new cacher from the allegro/bigcache
// As bigcache supports only string keys and byte values, the returned cacher uses
// JSON encoding/decoding for values.
//
// For example:
//
//	bc, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
//	memoize.NewWithCache(
//	    allegro.BigCache[int](bc),
//	    func(s string) int { return len(s) },
//	)
func BigCache[V any](c BigCacher) cache.Cacher[string, V] {
	return &bigCache[V]{c}
}

type bigCache[V any] struct {
	BigCacher
}

func (b *bigCache[V]) Load(key string) (value V, loaded bool) {
	data, err := b.BigCacher.Get(key)
	if err != nil {
		return value, false
	}

	var val V
	if err := json.Unmarshal(data, &val); err != nil {
		return value, false
	}

	return val, true
}

func (b *bigCache[V]) Store(key string, value V) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	b.BigCacher.Set(key, data)
}
