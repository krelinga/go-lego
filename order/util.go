package order

func Less[T any](a, b T, order Func[T]) bool {
	return order(a, b) < 0
}

func Greater[T any](a, b T, order Func[T]) bool {
	return order(a, b) > 0
}

func Equal[T any](a, b T, order Func[T]) bool {
	return order(a, b) == 0
}

func LessEqual[T any](a, b T, order Func[T]) bool {
	return order(a, b) <= 0
}

func GreaterEqual[T any](a, b T, order Func[T]) bool {
	return order(a, b) >= 0
}
