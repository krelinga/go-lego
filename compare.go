package lego

import (
	"cmp"
)

// A Comparer is a type that can be compared to other values of the same type using the Compare method, which returns an integer indicating the relative order of the values.
// The Compare method should return a negative integer if the value is less than the other value, a positive integer if the value is greater than the other value, and zero if the two values are equal.
type Comparer[T any] interface {
	Compare(T) int
}

func CompareViewer[V Viewer[T], T Comparer[T]](a, b V) int {
	return a.View().Compare(b.View())
}

// A CmpFunc is a function that compares two values of the same type and returns an integer indicating their relative order.
// The function should return a negative integer if the first value is less than the second value, a positive integer if the first value is greater than the second value, and zero if the two values are equal.
type CmpFunc[T any] func(T, T) int

// CompareUsing compares two values using the given comparison functions.
// The comparison functions are tried in order until one of them returns a non-zero value, which is returned by this function.
// If all comparison functions return zero, then this function returns zero.
//
// This function is useful for implementing the [Comparer] interface on a type that has multiple fields that need to be compared, since it allows you to write the comparison logic for each field separately and then combine them easily.
func CompareUsing[T any](a, b T, funcs ...CmpFunc[T]) int {
	for _, f := range funcs {
		if x := f(a, b); x != 0 {
			return x
		}
	}
	return 0
}

// NewCmpFunc creates a CmpFunc that compares two values by applying the given function to each value and then comparing the results using the Compare method of the results.
func NewCmpFunc[T any, F Comparer[F]](g func(T) F) CmpFunc[T] {
	return func(a, b T) int {
		return g(a).Compare(g(b))
	}
}

// NewCmpFuncGo creates a CmpFunc that compares two values by applying the given function to each value and then comparing the results using the cmp.Compare function from the cmp package.
func NewCmpFuncGo[T any, F cmp.Ordered](g func(T) F) CmpFunc[T] {
	return func(a, b T) int {
		return cmp.Compare(g(a), g(b))
	}
}

// NewCmpFuncUsing creates a CmpFunc that compares two values using the given comparison functions, which are tried in order until one of them returns a non-zero value, which is returned by the resulting CmpFunc. If all comparison functions return zero, then the resulting CmpFunc returns zero.
func NewCmpFuncUsing[T any](funcs ...CmpFunc[T]) CmpFunc[T] {
	return func(a, b T) int {
		return CompareUsing(a, b, funcs...)
	}
}

// NewCmpFuncReverse creates a CmpFunc that compares two values by applying the given CmpFunc and then negating the result, so that the order of the values is reversed.
func NewCmpFuncReverse[T any](f CmpFunc[T]) CmpFunc[T] {
	return func(a, b T) int {
		return -f(a, b)
	}
}

// Less returns true if a is less than b according to the Compare method of the [Comparer] interface, and false otherwise.
func Less[T Comparer[T]](a, b T) bool {
	return a.Compare(b) < 0
}

// LessEqual returns true if a is less than or equal to b according to the Compare method of the [Comparer] interface, and false otherwise.
func LessEqual[T Comparer[T]](a, b T) bool {
	return a.Compare(b) <= 0
}

// Greater returns true if a is greater than b according to the Compare method of the [Comparer] interface, and false otherwise.
func Greater[T Comparer[T]](a, b T) bool {
	return a.Compare(b) > 0
}

// GreaterEqual returns true if a is greater than or equal to b according to the Compare method of the [Comparer] interface, and false otherwise.
func GreaterEqual[T Comparer[T]](a, b T) bool {
	return a.Compare(b) >= 0
}
