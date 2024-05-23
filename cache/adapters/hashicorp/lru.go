// Package hashicorp includes LRU adapter for https://github.com/hashicorp/golang-lru/v2 package
package hashicorp

import (
	"github.com/emad-elsaid/memoize/cache"
)

// Hashicorp LRU cache interface (the part we need)
type HashicorpLRU[K comparable, V any] interface {
	Get(K) (V, bool)
	Add(K, V) bool
}

// LRU creates a new cacher from the hashicorp cache
// For example:
//
//	  lru, _ := lru.New[string, int](1000)
//	  memoize.NewWithCache(
//		   hashicorp.LRU(lru),
//		   func(s string) int { return len(s) },
//	  )
func LRU[K comparable, V any](c HashicorpLRU[K, V]) cache.Cacher[K, V] {
	return &lru[K, V]{c}
}

type lru[K comparable, V any] struct {
	HashicorpLRU[K, V]
}

func (h *lru[K, V]) Load(key K) (value V, loaded bool) {
	return h.HashicorpLRU.Get(key)
}

func (h *lru[K, V]) Store(key K, value V) {
	h.HashicorpLRU.Add(key, value)
}
