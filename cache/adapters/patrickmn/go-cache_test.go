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
