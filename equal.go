package lego

type Equaler[T any] interface {
	Equal(T) bool
}

func EqualComparer[T Comparer[T]](a, b T) bool {
	return a.Compare(b) == 0
}
