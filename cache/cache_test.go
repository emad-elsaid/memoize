package cache

import "testing"

func TestCache(t *testing.T) {
	t.Run("LoadOrStore", func(t *testing.T) {
		var c Cache[string, int]
		v, loaded := c.LoadOrStore("k1", 1)
		assertEqual(t, 1, v)
		assertEqual(t, false, loaded)

		v, loaded = c.LoadOrStore("k1", 2)
		assertEqual(t, 1, v)
		assertEqual(t, true, loaded)
	})

	t.Run("Load/Store", func(t *testing.T) {
		var c Cache[string, int]

		v, ok := c.Load("k1")
		assertEqual(t, 0, v)
		assertEqual(t, false, ok)

		c.Store("k1", 1)
		v, ok = c.Load("k1")
		assertEqual(t, 1, v)
		assertEqual(t, true, ok)
	})

	t.Run("Delete", func(t *testing.T) {
		var c Cache[string, int]

		c.Store("k1", 1)
		c.Delete("k1")

		v, ok := c.Load("k1")
		assertEqual(t, 0, v)
		assertEqual(t, false, ok)
	})
}
