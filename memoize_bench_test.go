package memoize

import (
	"fmt"
	"math"
	"testing"
	"time"
)

var returns int

func benchmarkMemoizer(keysN int, b *testing.B) {
	mem := MemoizerWithErr[string, int, func(string) (int, error)]{
		Cache: &MemoryCache[string, func() (int, error)]{},
		Fun: func(k string) (int, error) {
			return len(k), nil
		},
	}

	keys := []string{}
	for i := range keysN {
		keys = append(keys, string(fmt.Sprintf("Key%d", i)))
	}

	b.ResetTimer()
	for i := range b.N {
		k := keys[i%len(keys)]
		r, _ := mem.Fun(k)
		returns = r
	}
}

func BenchmarkMemoizer(b *testing.B) {
	for i := 1; i < 6; i++ {
		keyN := math.Pow(10, float64(i))
		b.Run(fmt.Sprintf("Keys:%.f", keyN), func(b *testing.B) {
			benchmarkMemoizer(int(keyN), b)
		})
	}
}
func benchmarkMemoizerPar(keysN int, wait time.Duration, b *testing.B) {
	mem := MemoizerWithErr[string, int, func(string) (int, error)]{
		Cache: &MemoryCache[string, func() (int, error)]{},
		Fun: func(k string) (int, error) {
			time.Sleep(wait)
			return len(k), nil
		},
	}

	keys := []string{}
	for i := range keysN {
		keys = append(keys, string(fmt.Sprintf("Key%d", i)))
	}

	b.ResetTimer()
	b.RunParallel(func(b *testing.PB) {
		i := 0
		for b.Next() {
			k := keys[i%len(keys)]
			r, _ := mem.Fun(k)
			returns = r
			i++
		}
	})
}

func BenchmarkMemoizerPar(b *testing.B) {
	durations := map[string]time.Duration{
		"Micro":    time.Microsecond,
		"10Micro":  time.Microsecond * 10,
		"100Micro": time.Microsecond * 100,
		"Milli":    time.Millisecond,
		"10Milli":  time.Millisecond * 10,
		"100Milli": time.Millisecond * 100,
		"Sec":      time.Second,
	}

	for i := 1; i < 6; i++ {
		keyN := math.Pow(10, float64(i))
		for s, d := range durations {
			d := d
			keyN := keyN
			b.Run(fmt.Sprintf("Keys:%.fWait:%s", keyN, s), func(b *testing.B) {
				benchmarkMemoizerPar(int(keyN), d, b)
			})
		}
	}
}
