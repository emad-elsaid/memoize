// Package bradfitz includes cache adapter for https://github.com/bradfitz/gomemcache package
package bradfitz

import (
	"encoding/json"

	"github.com/emad-elsaid/memoize/cache"
)

// Memcacher cache interface (the part memoize needs)
type Memcacher interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

// Memcache creates a new cacher from the bradfitz/gomemcache
// As gomemcache supports only string keys and byte values, the returned cacher
// uses JSON encoding/decoding for values.
//
// For example:
//
//	mc := memcache.New("localhost:11211")
//	memoize.NewWithCache(
//	    bradfitz.Memcache[int](mc),
//	    func(s string) int { return len(s) },
//	)
func Memcache[V any](c Memcacher) cache.Cacher[string, V] {
	return &memcache[V]{c}
}

type memcache[V any] struct {
	Memcacher
}

func (m *memcache[V]) Load(key string) (value V, loaded bool) {
	data, err := m.Memcacher.Get(key)
	if err != nil {
		return value, false
	}

	var val V
	if err := json.Unmarshal(data, &val); err != nil {
		return value, false
	}

	return val, true
}

func (m *memcache[V]) Store(key string, value V) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	m.Memcacher.Set(key, data)
}
