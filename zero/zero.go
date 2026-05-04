package zero

// For returns the zero value of the specified type T.
func For[T any]() T {
	var zero T
	return zero
}