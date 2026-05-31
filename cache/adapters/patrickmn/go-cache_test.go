package patrickmn

import (
	"testing"
)

type Mock map[string]any


func (m Mock) Get(k string) (any, bool) {
	v, ok := m[k]
	return v, ok
}

func (m Mock) SetDefault(k string, x any) {
	m[k] = x
}

func TestGoCache(t *testing.T) {
	m := Mock{}
	cacher := GoCache[int](m)
	cacher.Store("k1", 1)

	v, ok := cacher.Load("k1")
	if v != 1 {
		t.Error("Expected", 1, "got", v)
	}

	if !ok {
		t.Error("Expected", true, "got", ok)
	}
}

func TestGoCache_LoadNonExistent(t *testing.T) {
	m := Mock{}
	cacher := GoCache[int](m)

	v, ok := cacher.Load("nonexistent")
	if ok {
		t.Error("Expected ok to be false, got true")
	}
	if v != 0 {
		t.Errorf("Expected 0, got %d", v)
	}
}

// TestGoCache_LoadTypeMismatch verifies defensive behavior when the underlying
// cache contains values of the wrong type (e.g., if stored through direct cache
// access rather than through the adapter's Store method).
func TestGoCache_LoadTypeMismatch(t *testing.T) {
	m := Mock{
		"wrongtype": "string value",
	}
	cacher := GoCache[int](m)

	v, ok := cacher.Load("wrongtype")
	if ok {
		t.Error("Expected ok to be false for type mismatch, got true")
	}
	if v != 0 {
		t.Errorf("Expected 0, got %d", v)
	}
}
