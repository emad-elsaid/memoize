// Package rueidis includes cache adapter for https://github.com/redis/rueidis package
package rueidis

import (
	"context"
	"encoding/json"

	"github.com/emad-elsaid/memoize/cache"
)

// Rueidiser cache interface (the part memoize needs)
type Rueidiser interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}

// Rueidis creates a new cacher from the redis/rueidis
// As rueidis supports only string keys, the returned cacher uses
// JSON encoding/decoding for values.
//
// For example:
//
//	client, err := rueidis.NewClient(rueidis.ClientOption{
//	    InitAddress: []string{"localhost:6379"},
//	})
//	memoize.NewWithCache(
//	    rueidis.Rueidis[int](context.Background(), client),
//	    func(s string) int { return len(s) },
//	)
func Rueidis[V any](ctx context.Context, c Rueidiser) cache.Cacher[string, V] {
	return &rueidis[V]{c, ctx}
}

type rueidis[V any] struct {
	Rueidiser
	ctx context.Context
}

func (r *rueidis[V]) Load(key string) (value V, loaded bool) {
	data, err := r.Rueidiser.Get(r.ctx, key)
	if err != nil {
		return value, false
	}

	var val V
	if err := json.Unmarshal([]byte(data), &val); err != nil {
		return value, false
	}

	return val, true
}

func (r *rueidis[V]) Store(key string, value V) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	r.Rueidiser.Set(r.ctx, key, string(data))
}
