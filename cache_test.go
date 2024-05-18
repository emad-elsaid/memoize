package memoize

import "testing"

func TestCacheWithFallback(t *testing.T) {
	t.Run("Load", func(t *testing.T) {
		var c1, c2 Cache[string, int]
		c := CacheWithFallback[string, int]{
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
		c := CacheWithFallback[string, int]{
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

func TestReadOnlyCache(t *testing.T) {
	var c1 Cache[string, int]
	ro := ReadOnlyCache[string, int]{&c1}

	c1.Store("k1", 1)
	v, ok := ro.Load("k1")
	assertEqual(t, 1, v)
	assertEqual(t, true, ok)

	// Write new value and check again
	ro.Store("k1", 2)
	v, ok = ro.Load("k1")
	assertEqual(t, 1, v)
	assertEqual(t, true, ok)
}

func TestWriteOnlyCache(t *testing.T) {
	var c1 Cache[string, int]
	wo := WriteOnlyCache[string, int]{&c1}

	c1.Store("k1", 1)
	v, ok := wo.Load("k1")
	assertEqual(t, 0, v)
	assertEqual(t, false, ok)

	// Write new value and check again
	wo.Store("k1", 2)
	v, ok = wo.Load("k1")
	assertEqual(t, 0, v)
	assertEqual(t, false, ok)
}
