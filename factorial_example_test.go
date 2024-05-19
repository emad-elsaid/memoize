package memoize

import (
	"fmt"
)

func ExampleNew_factorial() {
	var factorial func(uint64) uint64
	factorial = New(func(i uint64) uint64 {
		if i == 0 {
			return 1
		}

		return i * factorial(i-1)
	})

	fmt.Println(factorial(10))
	// Output: 3628800
}
