# Memoize

[![Go Reference](https://pkg.go.dev/badge/github.com/emad-elsaid/memoize.svg)](https://pkg.go.dev/github.com/emad-elsaid/memoize)
[![Go Report Card](https://goreportcard.com/badge/github.com/emad-elsaid/memoize)](https://goreportcard.com/report/github.com/emad-elsaid/memoize)
[![codecov](https://codecov.io/gh/emad-elsaid/memoize/graph/badge.svg?token=QBXTR1XRD6)](https://codecov.io/gh/emad-elsaid/memoize)

Golang Memoize implementation for functions of type `func(any) any` and
`func(any) (any, error)`. Safe for concurrent use.

## Guarantees

Memoize premetives offers the following guarantees for `Memoizer` and `MemoizerWithErr`:
* Duplicate function call suppression
* 0 allocation
* Dependency only on stdlib
* Low overhead

And for `MemoizerWithCache*` premitives:
* Duplicate function call suppression
* Dependency only on stdlib


The reason `MemoizerWithCache` doesn't guarantee 0 allocation is that it needs extra synchronization mechanism to guarantee the duplicate function call suppression which adds around 5-8 allocations/op and adds an overhead around 400-800ns/op

## Usage

```go
var NewPage = memoize.New(func(name string) (p Page) {
    // your normal logic here
})
```

The previous snippet will make sure the function runs only once for each
value of `name` and will return the same value for the same input always.


Also function that returns 2 values can be memoized:

```go
var NewPage = memoize.NewWithErr(func(name string) (p Page, err error) {
    // your normal logic here
})
```

Also the package offers the same interface but with the option to provide your own cache

```go
var NewPage = memoize.NewWithCache(cache, func(name string) (p Page) {
    // your normal logic here
})
var NewPage = memoize.NewWithCacheErr(cache, func(name string) (p Page, err error) {
    // your normal logic here
})
```

where `cache` implements the interface:
```go
type Cacher[K any, V any] interface {
	Store(key K, value V)
	Load(key K) (value V, ok bool)
}
```

# Note on Cacher interface

* `memoize` require the cache interface to implement two simple `Load` and `Store` functions
* So you can adapt any other caching library to `memoize`
* This will give you a very powerful memoization patterns where you can store your cache in memory, file, remote system or have failovers with multiple layers of caching.
* This also means `memoize` package will not remove any items from the cache, it's the cache implementation responsibility to manage it's size, TTL, and communication with remote systems
* `memoize` comes with 1 concurrency safe implementation of the cache `Cache`, stored in memory and uses Go generics. packed by `sync.Map`

# Cache subpackage

* `memoize` include a subpackage `cache` which provides several implementations for the `Cacher` interface
* `Cache` is a simple in-memory forever cacher
* `WithFallback`, `WithReadOnly`, and `WithWriteOnly`...etc wraps a cacher or more to provide or supress functionality


> [!IMPORTANT]
> This package is still being tested, use it with caution

# Brenchmarks

Each struct is tested with two benchmarks:
* Sequencial executions
* Parallel executions

And each benchmark is repeated for different key space sizes to see how it reacts to smaller vs. large key spaces

```
pkg: github.com/emad-elsaid/memoize
cpu: 11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz
BenchmarkMemoizer/Keys:10-8               	52843860	       22.98 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizer/Keys:100-8              	48128289	       24.24 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizer/Keys:1000-8             	38378518	       30.21 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizer/Keys:10000-8            	30157987	       38.10 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizer/Keys:100000-8           	23937156	       49.28 ns/op	      1 B/op	      0 allocs/op
BenchmarkMemoizerParallel/Keys:10-8       	91730824	       12.99 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizerParallel/Keys:100-8      	88690077	       13.41 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizerParallel/Keys:1000-8     	70860706	       14.15 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizerParallel/Keys:10000-8    	79582891	       15.17 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizerParallel/Keys:100000-8   	48657999	       23.51 ns/op	      1 B/op	      0 allocs/op
BenchmarkMemoizerWithCache/Keys:10-8      	2408961	      511.0 ns/op	    336 B/op	      8 allocs/op
BenchmarkMemoizerWithCache/Keys:100-8     	2546571	      454.8 ns/op	    336 B/op	      8 allocs/op
BenchmarkMemoizerWithCache/Keys:1000-8    	1976749	      560.7 ns/op	    336 B/op	      8 allocs/op
BenchmarkMemoizerWithCache/Keys:10000-8   	2341108	      491.9 ns/op	    336 B/op	      8 allocs/op
BenchmarkMemoizerWithCache/Keys:100000-8  	2158876	      549.7 ns/op	    342 B/op	      8 allocs/op
BenchmarkMemoizerWithCacheParallel/Keys:10-8         	2487337	      479.9 ns/op	    313 B/op	      7 allocs/op
BenchmarkMemoizerWithCacheParallel/Keys:100-8        	2412850	      495.6 ns/op	    299 B/op	      7 allocs/op
BenchmarkMemoizerWithCacheParallel/Keys:1000-8       	2473747	      492.4 ns/op	    188 B/op	      5 allocs/op
BenchmarkMemoizerWithCacheParallel/Keys:10000-8      	2268688	      528.3 ns/op	    151 B/op	      5 allocs/op
BenchmarkMemoizerWithCacheParallel/Keys:100000-8     	1733913	      717.1 ns/op	    195 B/op	      5 allocs/op
BenchmarkMemoizerWithCacheErr/Keys:10-8              	2456088	      514.3 ns/op	    336 B/op	      8 allocs/op
BenchmarkMemoizerWithCacheErr/Keys:100-8             	2184484	      485.3 ns/op	    336 B/op	      8 allocs/op
BenchmarkMemoizerWithCacheErr/Keys:1000-8            	2106156	      529.6 ns/op	    336 B/op	      8 allocs/op
BenchmarkMemoizerWithCacheErr/Keys:10000-8           	2386858	      503.9 ns/op	    336 B/op	      8 allocs/op
BenchmarkMemoizerWithCacheErr/Keys:100000-8          	2040483	      648.2 ns/op	    343 B/op	      8 allocs/op
BenchmarkMemoizerWithCacheErrParallel/Keys:10-8      	1978743	      590.1 ns/op	    314 B/op	      7 allocs/op
BenchmarkMemoizerWithCacheErrParallel/Keys:100-8     	1960455	      602.6 ns/op	    297 B/op	      7 allocs/op
BenchmarkMemoizerWithCacheErrParallel/Keys:1000-8    	2046734	      581.2 ns/op	    182 B/op	      5 allocs/op
BenchmarkMemoizerWithCacheErrParallel/Keys:10000-8   	1821135	      644.7 ns/op	    149 B/op	      5 allocs/op
BenchmarkMemoizerWithCacheErrParallel/Keys:100000-8  	1314518	      931.3 ns/op	    203 B/op	      6 allocs/op
BenchmarkMemoizerWithErr/Keys:10-8                   	43579936	       25.06 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizerWithErr/Keys:100-8                  	46920016	       28.24 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizerWithErr/Keys:1000-8                 	36951790	       33.61 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizerWithErr/Keys:10000-8                	27627970	       41.76 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizerWithErr/Keys:100000-8               	20259408	       57.24 ns/op	      1 B/op	      0 allocs/op
BenchmarkMemoizerWithErrParallel/Keys:10-8           	90412166	       20.23 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizerWithErrParallel/Keys:100-8          	47339029	       24.99 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizerWithErrParallel/Keys:1000-8         	43452195	       25.43 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizerWithErrParallel/Keys:10000-8        	55278070	       22.31 ns/op	      0 B/op	      0 allocs/op
BenchmarkMemoizerWithErrParallel/Keys:100000-8       	29371357	       37.95 ns/op	      2 B/op	      0 allocs/op
PASS
ok  	github.com/emad-elsaid/memoize	67.529s
```
