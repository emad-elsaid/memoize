package redis

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

func (m Mock) Set(ctx context.Context, key string, value any) error {
	if str, ok := value.(string); ok {
		m[key] = str
		return nil
	}
	return errors.New("invalid type")
}

func TestGoRedis_StoreAndLoad(t *testing.T) {
	m := Mock{}
	ctx := context.Background()
	cacher := GoRedis[int](ctx, m)

	cacher.Store("k1", 42)

	v, ok := cacher.Load("k1")
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	if v != 42 {
		t.Errorf("Expected 42, got %d", v)
	}
}

func TestGoRedis_LoadNonExistent(t *testing.T) {
	m := Mock{}
	ctx := context.Background()
	cacher := GoRedis[int](ctx, m)

	v, ok := cacher.Load("nonexistent")
	if ok {
		t.Error("Expected ok to be false, got true")
	}
	if v != 0 {
		t.Errorf("Expected 0, got %d", v)
	}
}

func TestGoRedis_LoadInvalidJSON(t *testing.T) {
	m := Mock{
		"invalid": "not json",
	}
	ctx := context.Background()
	cacher := GoRedis[int](ctx, m)

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

func TestGoRedis_StoreMarshalError(t *testing.T) {
	m := Mock{}
	ctx := context.Background()
	cacher := GoRedis[FailMarshalType](ctx, m)

	// Channel types cannot be marshaled to JSON
	cacher.Store("k1", FailMarshalType{Ch: make(chan int)})

	// Should not store anything
	_, exists := m["k1"]
	if exists {
		t.Error("Expected key to not exist after marshal error")
	}
}

func TestGoRedis_StoreComplexType(t *testing.T) {
	m := Mock{}
	ctx := context.Background()
	cacher := GoRedis[map[string]int](ctx, m)

	expected := map[string]int{"a": 1, "b": 2}
	cacher.Store("complex", expected)

	v, ok := cacher.Load("complex")
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	if len(v) != 2 || v["a"] != 1 || v["b"] != 2 {
		t.Errorf("Expected %v, got %v", expected, v)
	}
}
