package order

import "cmp"

func Using[T any](a, b T, orders ...Func[T]) int {
	for _, order := range orders {
		if cmp := order(a, b); cmp != 0 {
			return cmp
		}
	}
	return 0
}

func Get[T any, U cmp.Ordered](get func(T) U) Func[T] {
	return func(x, y T) int {
		return cmp.Compare(get(x), get(y))
	}
}

func GetFunc[T any, U any](get func(T) U, order func(U, U) int) Func[T] {
	return func(x, y T) int {
		return order(get(x), get(y))
	}
}