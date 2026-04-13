package orders

func Desc[T any](order Func[T]) Func[T] {
	return func(x, y T) int {
		return -order(x, y)
	}
}