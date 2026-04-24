package order

// Desc wraps an ordering function and reverses its result, producing a descending order.
func Desc[T any](order Func[T]) Func[T] {
	return func(x, y T) int {
		return -order(x, y)
	}
}
