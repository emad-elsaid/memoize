# Memoize

Golang Memoize implementation for functions of type `func(any) any` and
`func(any) (any, error)`. Safe for concurrent use

Simple usecase

```go
var NewPage = memoize.New(func(name string) (p Page) {
    // your normal logic here
})
```

in the previous snippet will make sure the function runs only once for each
value of `name` and will return the same value for the same input always.


Also function that returns 2 values can be memoized:

```go
var NewPage = memoize.NewWithErr(func(name string) (p Page, err error) {
    // your normal logic here
})
```

> [!IMPORTANT]
> This package is still being tested, use it with caution

# Brenchmarks results

```
➜ go test --bench .
goos: linux
goarch: amd64
pkg: github.com/emad-elsaid/memoize
cpu: AMD Ryzen 7 2700X Eight-Core Processor
BenchmarkMemoizer/Keys:10-16               	25258276	       43.08 ns/op
BenchmarkMemoizer/Keys:100-16              	25672358	       46.79 ns/op
BenchmarkMemoizer/Keys:1000-16             	27347005	       42.03 ns/op
BenchmarkMemoizer/Keys:10000-16            	21916824	       55.37 ns/op
BenchmarkMemoizer/Keys:100000-16           	14022218	       84.96 ns/op
BenchmarkMemoizerParallel/Keys:10-16       	45314600	       26.33 ns/op
BenchmarkMemoizerParallel/Keys:100-16      	45316003	       25.73 ns/op
BenchmarkMemoizerParallel/Keys:1000-16     	41009624	       26.64 ns/op
BenchmarkMemoizerParallel/Keys:10000-16    	37614056	       28.19 ns/op
BenchmarkMemoizerParallel/Keys:100000-16   	2024397	      537.5 ns/op
PASS
ok  	github.com/emad-elsaid/memoize	29.049s
```
