package allegro

import (
	"errors"
	"testing"
)

type Mock map[string][]byte

func (m Mock) Get(key string) ([]byte, error) {
	v, ok := m[key]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (m Mock) Set(key string, entry []byte) error {
	m[key] = entry
	return nil
}

func TestBigCache_StoreAndLoad(t *testing.T) {
	m := Mock{}
	cacher := BigCache[int](m)

	cacher.Store("k1", 42)

	v, ok := cacher.Load("k1")
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	if v != 42 {
		t.Errorf("Expected 42, got %d", v)
	}
}

func TestBigCache_LoadNonExistent(t *testing.T) {
	m := Mock{}
	cacher := BigCache[int](m)

	v, ok := cacher.Load("nonexistent")
	if ok {
		t.Error("Expected ok to be false, got true")
	}
	if v != 0 {
		t.Errorf("Expected 0, got %d", v)
	}
}

func TestBigCache_LoadInvalidJSON(t *testing.T) {
	m := Mock{
		"invalid": []byte("not json"),
	}
	cacher := BigCache[int](m)

	v, ok := cacher.Load("invalid")
	if ok {
		t.Error("Expected ok to be false for invalid JSON, got true")
	}
	if v != 0 {
		t.Errorf("Expected 0, got %d", v)
	}
}

type FailMarshalType struct {
	Ch chan int
}

func TestBigCache_StoreMarshalError(t *testing.T) {
	m := Mock{}
	cacher := BigCache[FailMarshalType](m)

	// Channel types cannot be marshaled to JSON
	cacher.Store("k1", FailMarshalType{Ch: make(chan int)})

	// Should not store anything
	_, exists := m["k1"]
	if exists {
		t.Error("Expected key to not exist after marshal error")
	}
}

func TestBigCache_StoreSlice(t *testing.T) {
	m := Mock{}
	cacher := BigCache[[]int](m)

	expected := []int{1, 2, 3, 4, 5}
	cacher.Store("numbers", expected)

	v, ok := cacher.Load("numbers")
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	if len(v) != 5 {
		t.Errorf("Expected length 5, got %d", len(v))
	}
	for i, num := range expected {
		if v[i] != num {
			t.Errorf("At index %d: expected %d, got %d", i, num, v[i])
		}
	}
}
