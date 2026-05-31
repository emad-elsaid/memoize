package rueidis

import (
	"context"
	"errors"
	"testing"
)

type Mock map[string]string

func (m Mock) Get(ctx context.Context, key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func (m Mock) Set(ctx context.Context, key string, value string) error {
	m[key] = value
	return nil
}

func TestRueidis_StoreAndLoad(t *testing.T) {
	m := Mock{}
	ctx := context.Background()
	cacher := Rueidis[int](ctx, m)

	cacher.Store("k1", 42)

	v, ok := cacher.Load("k1")
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	if v != 42 {
		t.Errorf("Expected 42, got %d", v)
	}
}

func TestRueidis_LoadNonExistent(t *testing.T) {
	m := Mock{}
	ctx := context.Background()
	cacher := Rueidis[int](ctx, m)

	v, ok := cacher.Load("nonexistent")
	if ok {
		t.Error("Expected ok to be false, got true")
	}
	if v != 0 {
		t.Errorf("Expected 0, got %d", v)
	}
}

func TestRueidis_LoadInvalidJSON(t *testing.T) {
	m := Mock{
		"invalid": "not json",
	}
	ctx := context.Background()
	cacher := Rueidis[int](ctx, m)

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

func TestRueidis_StoreMarshalError(t *testing.T) {
	m := Mock{}
	ctx := context.Background()
	cacher := Rueidis[FailMarshalType](ctx, m)

	// Channel types cannot be marshaled to JSON
	cacher.Store("k1", FailMarshalType{Ch: make(chan int)})

	// Should not store anything
	_, exists := m["k1"]
	if exists {
		t.Error("Expected key to not exist after marshal error")
	}
}

func TestRueidis_StoreComplexType(t *testing.T) {
	m := Mock{}
	ctx := context.Background()
	cacher := Rueidis[[]string](ctx, m)

	expected := []string{"a", "b", "c"}
	cacher.Store("list", expected)

	v, ok := cacher.Load("list")
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	if len(v) != 3 || v[0] != "a" || v[1] != "b" || v[2] != "c" {
		t.Errorf("Expected %v, got %v", expected, v)
	}
}
