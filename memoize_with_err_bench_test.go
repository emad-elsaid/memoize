package memoize

import (
	"fmt"
	"math"
	"testing"
)

func BenchmarkMemoizerWithErr(b *testing.B) {
	for i := 1; i < maxKeySpace; i++ {
		keysN := math.Pow(10, float64(i))
		name := fmt.Sprintf("Keys:%.f", keysN)

		b.Run(name, func(b *testing.B) {

			mem := NewWithErr(func(k string) (int, error) { return len(k), nil })

			keys := []string{}
			for i := range int64(keysN) {
				keys = append(keys, string(fmt.Sprintf("Key%d", i)))
			}

			b.ResetTimer()
			for i := range b.N {
				k := keys[i%len(keys)]
				r, _ := mem(k)
				returns = r
			}

		})
	}
}

func BenchmarkMemoizerWithErrParallel(b *testing.B) {
	for i := 1; i < maxKeySpace; i++ {
		keysN := math.Pow(10, float64(i))
		name := fmt.Sprintf("Keys:%.f", keysN)

		b.Run(name, func(b *testing.B) {

			mem := NewWithErr(func(k string) (int, error) { return len(k), nil })

			keys := []string{}
			for i := range int64(keysN) {
				keys = append(keys, string(fmt.Sprintf("Key%d", i)))
			}

			b.ResetTimer()
			b.RunParallel(func(b *testing.PB) {
				i := 0
				for b.Next() {
					k := keys[i%len(keys)]
					r, _ := mem(k)
					returns = r
					i++
				}
			})

		})
	}
}
