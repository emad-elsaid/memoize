package dgraphio

import (
	"testing"
)

type Mock map[string]int

func (m Mock) Get(k string) (int, bool) {
	v, ok := m[k]
	return v, ok
}

func (m Mock) Set(k string, v int, _ int64) bool {
	m[k] = v
	return true
}

func TestRistretto(t *testing.T) {
	m := Mock{}
	cacher := Ristretto(m)
	cacher.Store("k1", 1)

	v, ok := cacher.Load("k1")
	if v != 1 {
		t.Error("Expected", 1, "got", v)
	}

	if !ok {
		t.Error("Expected", true, "got", ok)
	}
}
