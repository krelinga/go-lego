package order

type Op[T any] = func(a, b T) bool

func OpLess[T any](order Func[T]) Op[T] {
	return func(a, b T) bool {
		return order(a, b) < 0
	}
}

func OpGreater[T any](order Func[T]) Op[T] {
	return func(a, b T) bool {
		return order(a, b) > 0
	}
}

func OpEqual[T any](order Func[T]) Op[T] {
	return func(a, b T) bool {
		return order(a, b) == 0
	}
}

func OpLessEqual[T any](order Func[T]) Op[T] {
	return func(a, b T) bool {
		return order(a, b) <= 0
	}
}

func OpGreaterEqual[T any](order Func[T]) Op[T] {
	return func(a, b T) bool {
		return order(a, b) >= 0
	}
}