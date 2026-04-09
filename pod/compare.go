package v2

import "cmp"

type CanCompare[V any] interface {
	Compare(V) int
}

func LessThan[V CanCompare[V]](a, b V) bool {
	return a.Compare(b) < 0
}

func GreaterThan[V CanCompare[V]](a, b V) bool {
	return a.Compare(b) > 0
}

func LessThanEqual[V CanCompare[V]](a, b V) bool {
	return a.Compare(b) <= 0
}

func GreaterThanEqual[V CanCompare[V]](a, b V) bool {
	return a.Compare(b) >= 0
}

type Comparator[V any] func(V, V) int

func CompareUsing[V any](a, b V, funcs ...Comparator[V]) int {
	for _, f := range funcs {
		if cmp := f(a, b); cmp != 0 {
			return cmp
		}
	}
	return 0
}

func NewComparator[V any, F CanCompare[F]](g func(V) F) Comparator[V] {
	return func(a, b V) int {
		return g(a).Compare(g(b))
	}
}

func NewComparatorOrdered[V any, F cmp.Ordered](g func(V) F) Comparator[V] {
	return func(a, b V) int {
		return cmp.Compare(g(a), g(b))
	}
}

func NewComparatorReversed[V any](f Comparator[V]) Comparator[V] {
	return func(a, b V) int {
		return -f(a, b)
	}
}
