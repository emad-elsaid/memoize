package cache

import "testing"

func TestNoop(t *testing.T) {
	t.Run("Implements Cacher", func(t *testing.T) {
		var _ Cacher[any, any] = Noop[any, any]{}
	})

	t.Run("Discards writes", func(t *testing.T) {
		c := Noop[string, int]{}
		c.Store("k1", 1)
		v, loaded := c.Load("k1")
		assertEqual(t, 0, v)
		assertEqual(t, false, loaded)
	})
}
