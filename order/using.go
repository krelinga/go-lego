package order

import "cmp"

// Using compares x and y by applying each ordering function in sequence. It returns the result of
// the first function that produces a non-zero value, or zero if all functions consider x and y
// equal.
func Using[T any](x, y T, orders ...Func[T]) int {
	for _, order := range orders {
		if cmp := order(x, y); cmp != 0 {
			return cmp
		}
	}
	return 0
}

// Get returns an ordering function for type T that compares values by extracting a cmp.Ordered
// field or property using the provided get function.
func Get[T any, U cmp.Ordered](get func(T) U) Func[T] {
	return func(x, y T) int {
		return cmp.Compare(get(x), get(y))
	}
}

// GetFunc returns an ordering function for type T that extracts a value using get and then compares
// those values using the provided order function.
func GetFunc[T any, U any](get func(T) U, order Func[U]) Func[T] {
	return func(x, y T) int {
		return order(get(x), get(y))
	}
}
