package cache

import "testing"

func TestWithReadOnly(t *testing.T) {
	var c1 Cache[string, int]
	ro := WithReadOnly(&c1)

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
