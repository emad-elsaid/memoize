// Package redis includes cache adapter for https://github.com/redis/go-redis package
package redis

import (
	"context"
	"encoding/json"

	"github.com/emad-elsaid/memoize/cache"
)

// GoRediser cache interface (the part memoize needs)
type GoRediser interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any) error
}

// GoRedis creates a new cacher from the redis/go-redis
// As go-redis supports only string keys, the returned cacher uses
// JSON encoding/decoding for values.
//
// For example:
//
//	rdb := redis.NewClient(&redis.Options{
//	    Addr: "localhost:6379",
//	})
//	memoize.NewWithCache(
//	    redis.GoRedis[int](context.Background(), rdb),
//	    func(s string) int { return len(s) },
//	)
func GoRedis[V any](ctx context.Context, c GoRediser) cache.Cacher[string, V] {
	return &goRedis[V]{c, ctx}
}

type goRedis[V any] struct {
	GoRediser
	ctx context.Context
}

func (r *goRedis[V]) Load(key string) (value V, loaded bool) {
	data, err := r.GoRediser.Get(r.ctx, key)
	if err != nil {
		return value, false
	}

	var val V
	if err := json.Unmarshal([]byte(data), &val); err != nil {
		return value, false
	}

	return val, true
}

func (r *goRedis[V]) Store(key string, value V) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	r.GoRediser.Set(r.ctx, key, string(data))
}
