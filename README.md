# Memoize

Golang Memoize implementation for functions of type `func(any) any` and
`func(any) (any, error)`. Safe for concurrent use

Simple usecase

```go
var NewPage = MemoryMemoizer(func(name string) (p Page) {
    // your normal logic here
})
```

in the previous snippet will make sure the function runs only once for each
value of `name` and will return the same value for the same input always.


Also function that returns 2 values can be memoized:

```go
var NewPage = MemoryMemoizerWithErr(func(name string) (p Page, err error) {
    // your normal logic here
})
```

> [!IMPORTANT]
> This package is still being tested, use it with caution

# Brenchmarks results against gopl.io

```
go test --bench .
goos: linux
goarch: amd64
pkg: github.com/emad-elsaid/memoize
cpu: AMD Ryzen 7 2700X Eight-Core Processor
BenchmarkGoplio/Keys:10-16      	1000000	     1364 ns/op
BenchmarkGoplio/Keys:100-16     	 787778	     1459 ns/op
BenchmarkGoplio/Keys:1000-16    	 870450	     1449 ns/op
BenchmarkGoplio/Keys:10000-16   	 859154	     1396 ns/op
BenchmarkGoplio/Keys:100000-16  	 620088	     1733 ns/op
BenchmarkGoplioParallel/Keys:10-16         	1000000	     1090 ns/op
BenchmarkGoplioParallel/Keys:100-16        	1000000	     1165 ns/op
BenchmarkGoplioParallel/Keys:1000-16       	 999837	     1176 ns/op
BenchmarkGoplioParallel/Keys:10000-16      	 985710	     1200 ns/op
BenchmarkGoplioParallel/Keys:100000-16     	 879453	     1379 ns/op
BenchmarkMemoizer/Keys:10-16               	2594858	      439.7 ns/op
BenchmarkMemoizer/Keys:100-16              	2828440	      459.8 ns/op
BenchmarkMemoizer/Keys:1000-16             	2630336	      455.1 ns/op
BenchmarkMemoizer/Keys:10000-16            	2582623	      456.1 ns/op
BenchmarkMemoizer/Keys:100000-16           	2357390	      464.5 ns/op
BenchmarkMemoizerParallel/Keys:10-16       	14447314	       84.29 ns/op
BenchmarkMemoizerParallel/Keys:100-16      	14607079	       83.06 ns/op
BenchmarkMemoizerParallel/Keys:1000-16     	14245002	       84.10 ns/op
BenchmarkMemoizerParallel/Keys:10000-16    	12791991	       91.26 ns/op
BenchmarkMemoizerParallel/Keys:100000-16   	2490600	      450.4 ns/op
PASS
ok  	github.com/emad-elsaid/memoize	33.902s
```
