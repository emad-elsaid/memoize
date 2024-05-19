package cache

import "testing"

func TestWithFallback(t *testing.T) {
	t.Run("Load", func(t *testing.T) {
		var c1, c2 Cache[string, int]
		c := WithFallback[string, int]{
			Caches: []Cacher[string, int]{&c1, &c2},
		}

		// Looking up non-existent
		v, ok := c.Load("non-existent")
		assertEqual(t, 0, v)
		assertEqual(t, false, ok)

		// Loading key that exists in first level
		c1.Store("k1", 1)
		v, ok = c.Load("k1")
		assertEqual(t, 1, v)
		assertEqual(t, true, ok)

		// Loading key that exists in second level
		c2.Store("k2", 2)
		v, ok = c.Load("k2")
		assertEqual(t, 2, v)
		assertEqual(t, true, ok)
		// writes to the first level
		v, ok = c1.Load("k2")
		assertEqual(t, 2, v)
		assertEqual(t, true, ok)
	})

	t.Run("Store", func(t *testing.T) {
		var c1, c2 Cache[string, int]
		c := WithFallback[string, int]{
			Caches: []Cacher[string, int]{&c1, &c2},
		}

		c.Store("k1", 1)

		v, ok := c.Load("k1")
		assertEqual(t, 1, v)
		assertEqual(t, true, ok)

		v, ok = c1.Load("k1")
		assertEqual(t, 1, v)
		assertEqual(t, true, ok)

		v, ok = c2.Load("k1")
		assertEqual(t, 1, v)
		assertEqual(t, true, ok)
	})
}
