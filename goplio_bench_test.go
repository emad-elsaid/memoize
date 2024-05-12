package memoize

import (
	"fmt"
	"math"
	"testing"
)

func BenchmarkGoplio(b *testing.B) {
	for i := 1; i < maxKeySpace; i++ {
		keysN := math.Pow(10, float64(i))
		name := fmt.Sprintf("Keys:%.f", keysN)

		b.Run(name, func(b *testing.B) {

			mem := New(func(k string) (any, error) { return len(k), nil })

			keys := []string{}
			for i := range int64(keysN) {
				keys = append(keys, string(fmt.Sprintf("Key%d", i)))
			}

			b.ResetTimer()
			for i := range b.N {
				k := keys[i%len(keys)]
				r, _ := mem.Get(k)
				returns = r.(int)
			}

		})
	}
}

func BenchmarkGoplioParallel(b *testing.B) {
	for i := 1; i < maxKeySpace; i++ {
		keysN := math.Pow(10, float64(i))
		name := fmt.Sprintf("Keys:%.f", keysN)

		b.Run(name, func(b *testing.B) {

			mem := New(func(k string) (any, error) { return len(k), nil })

			keys := []string{}
			for i := range int64(keysN) {
				keys = append(keys, string(fmt.Sprintf("Key%d", i)))
			}

			b.ResetTimer()
			b.RunParallel(func(b *testing.PB) {
				i := 0
				for b.Next() {
					k := keys[i%len(keys)]
					r, _ := mem.Get(k)
					returns = r.(int)
					i++
				}
			})

		})
	}
}
