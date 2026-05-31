package goburrow

import (
	"testing"
)

type Mock[K, V any] map[any]any

func (m Mock[K, V]) GetIfPresent(key K) (V, bool) {
	v, ok := m[key]
	if !ok {
		var zero V
		return zero, false
	}
	val, ok := v.(V)
	return val, ok
}

func (m Mock[K, V]) Put(key K, value V) {
	m[key] = value
}

func TestMango_StoreAndLoad(t *testing.T) {
	m := make(Mock[string, int])
	cacher := Mango[string, int, string, int](m)

	cacher.Store("k1", 42)

	v, ok := cacher.Load("k1")
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	if v != 42 {
		t.Errorf("Expected 42, got %d", v)
	}
}

func TestMango_LoadNonExistent(t *testing.T) {
	m := make(Mock[string, int])
	cacher := Mango[string, int, string, int](m)

	v, ok := cacher.Load("nonexistent")
	if ok {
		t.Error("Expected ok to be false, got true")
	}
	if v != 0 {
		t.Errorf("Expected 0, got %d", v)
	}
}

func TestMango_LoadKeyTypeMismatch(t *testing.T) {
	m := make(Mock[string, int])
	cacher := Mango[string, int, int, int](m) // Key type mismatch

	v, ok := cacher.Load(123)
	if ok {
		t.Error("Expected ok to be false for key type mismatch, got true")
	}
	if v != 0 {
		t.Errorf("Expected 0, got %d", v)
	}
}

func TestMango_StoreKeyTypeMismatch(t *testing.T) {
	m := make(Mock[string, int])
	cacher := Mango[string, int, int, int](m) // Key type mismatch

	cacher.Store(123, 42)

	// Should not store anything due to type mismatch
	if len(m) != 0 {
		t.Error("Expected cache to be empty after key type mismatch")
	}
}

func TestMango_StoreValueTypeMismatch(t *testing.T) {
	m := make(Mock[string, int])
	cacher := Mango[string, int, string, string](m) // Value type mismatch

	cacher.Store("k1", "value")

	// Should not store anything due to type mismatch
	if len(m) != 0 {
		t.Error("Expected cache to be empty after value type mismatch")
	}
}

func TestMango_LoadValueTypeMismatch(t *testing.T) {
	m := make(Mock[string, int])
	m["k1"] = "wrong type" // Store wrong type directly
	cacher := Mango[string, int, string, int](m)

	v, ok := cacher.Load("k1")
	if ok {
		t.Error("Expected ok to be false for value type mismatch, got true")
	}
	if v != 0 {
		t.Errorf("Expected 0, got %d", v)
	}
}
