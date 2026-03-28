package lego

import "iter"

// An Adder is any container that values can be added to.
// The semantics of adding values to an Adder are not defined by this interface, and may vary between different implementations.
// For example, a Map is an Adder that may replace the value associated with a key if a value with the same key is added, while a Set is an Adder that does not allow duplicate values and ignores values that are already in the set.
type Adder[T any] interface {
	Add(T)
}

// Add adds the values from the given sequence to the given Adder.
func Add[A Adder[T], T any](values iter.Seq[T], a A) {
	for value := range values {
		a.Add(value)
	}
}

// Collect collects the values from the given sequence into a new Adder of the specified type, and returns the Adder.
func Collect[A Adder[T], T any](values iter.Seq[T]) A {
	var a A
	Add(values, a)
	return a
}