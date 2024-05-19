package cache

import "testing"

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
