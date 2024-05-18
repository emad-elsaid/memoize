package memoize

import "testing"

func assertEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if actual != expected {
		t.Errorf("Value not equal:\n\tExpected: %#v\n\tActual: %#v", expected, actual)
	}
}

func assertEqualMaps[K, V comparable](t *testing.T, expected, actual map[K]V) {
	t.Helper()
	for k, v := range expected {
		if v != actual[k] {
			t.Errorf("Key: %v value is not equal:\n\tExpected: %#v\n\tActual: %#v", k, v, actual[k])
		}
	}
	for k, v := range actual {
		if _, ok := expected[k]; !ok {
			t.Errorf("Actual map has extra key\n\tKey: %#v\n\tValue: %#v", k, v)
		}
	}
}
