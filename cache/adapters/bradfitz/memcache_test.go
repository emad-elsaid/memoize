package bradfitz

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

func (m Mock) Set(key string, value []byte) error {
	m[key] = value
	return nil
}

func TestMemcache_StoreAndLoad(t *testing.T) {
	m := Mock{}
	cacher := Memcache[int](m)

	cacher.Store("k1", 42)

	v, ok := cacher.Load("k1")
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	if v != 42 {
		t.Errorf("Expected 42, got %d", v)
	}
}

func TestMemcache_LoadNonExistent(t *testing.T) {
	m := Mock{}
	cacher := Memcache[int](m)

	v, ok := cacher.Load("nonexistent")
	if ok {
		t.Error("Expected ok to be false, got true")
	}
	if v != 0 {
		t.Errorf("Expected 0, got %d", v)
	}
}

func TestMemcache_LoadInvalidJSON(t *testing.T) {
	m := Mock{
		"invalid": []byte("not json"),
	}
	cacher := Memcache[int](m)

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

func TestMemcache_StoreMarshalError(t *testing.T) {
	m := Mock{}
	cacher := Memcache[FailMarshalType](m)

	// Channel types cannot be marshaled to JSON
	cacher.Store("k1", FailMarshalType{Ch: make(chan int)})

	// Should not store anything
	_, exists := m["k1"]
	if exists {
		t.Error("Expected key to not exist after marshal error")
	}
}
