package coocood

import (
	"errors"
	"testing"
)

type Mock struct {
	data          map[string][]byte
	expireSeconds int
}

func NewMock() *Mock {
	return &Mock{
		data: make(map[string][]byte),
	}
}

func (m *Mock) Get(key []byte) ([]byte, error) {
	v, ok := m.data[string(key)]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (m *Mock) Set(key, value []byte, expireSeconds int) error {
	m.data[string(key)] = value
	m.expireSeconds = expireSeconds
	return nil
}

func TestFreeCache_StoreAndLoad(t *testing.T) {
	m := NewMock()
	cacher := FreeCache[int](m, 3600)

	cacher.Store("k1", 42)

	v, ok := cacher.Load("k1")
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	if v != 42 {
		t.Errorf("Expected 42, got %d", v)
	}
	if m.expireSeconds != 3600 {
		t.Errorf("Expected expireSeconds to be 3600, got %d", m.expireSeconds)
	}
}

func TestFreeCache_LoadNonExistent(t *testing.T) {
	m := NewMock()
	cacher := FreeCache[int](m, 3600)

	v, ok := cacher.Load("nonexistent")
	if ok {
		t.Error("Expected ok to be false, got true")
	}
	if v != 0 {
		t.Errorf("Expected 0, got %d", v)
	}
}

func TestFreeCache_LoadInvalidJSON(t *testing.T) {
	m := NewMock()
	m.data["invalid"] = []byte("not json")
	cacher := FreeCache[int](m, 3600)

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

func TestFreeCache_StoreMarshalError(t *testing.T) {
	m := NewMock()
	cacher := FreeCache[FailMarshalType](m, 3600)

	// Channel types cannot be marshaled to JSON
	cacher.Store("k1", FailMarshalType{Ch: make(chan int)})

	// Should not store anything
	_, exists := m.data["k1"]
	if exists {
		t.Error("Expected key to not exist after marshal error")
	}
}

func TestFreeCache_StoreStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	m := NewMock()
	cacher := FreeCache[Person](m, 7200)

	expected := Person{Name: "Alice", Age: 30}
	cacher.Store("person", expected)

	v, ok := cacher.Load("person")
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	if v.Name != "Alice" || v.Age != 30 {
		t.Errorf("Expected %v, got %v", expected, v)
	}
	if m.expireSeconds != 7200 {
		t.Errorf("Expected expireSeconds to be 7200, got %d", m.expireSeconds)
	}
}
