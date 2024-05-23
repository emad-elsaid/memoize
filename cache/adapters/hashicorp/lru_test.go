package hashicorp

import (
	"testing"
)

type Mock map[string]int

func (m Mock) Get(k string) (int, bool) {
	v, ok := m[k]
	return v, ok
}

func (m Mock) Add(k string, v int) bool {
	m[k] = v
	return true
}

func TestLRU(t *testing.T) {
	m := Mock{}
	cacher := LRU(m)
	cacher.Store("k1", 1)

	v, ok := cacher.Load("k1")
	if v != 1 {
		t.Error("Expected", 1, "got", v)
	}

	if !ok {
		t.Error("Expected", true, "got", ok)
	}
}
