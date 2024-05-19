package memoize

import (
	"fmt"
)

// a function that calculates factorial recursively and memoize each input/output pair
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

// Create a log function that prints the message only
// once. no matter how many times it's called for the same message
func ExampleNew_uniqueLogger() {
	log := New(func(msg string) bool {
		fmt.Println(msg)
		return true
	})

	log("Hello World!")
	log("Hello World!")
	log("We came in peace")
	log("We came in peace")

	// Output: Hello World!
	// We came in peace
}
