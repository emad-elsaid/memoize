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
