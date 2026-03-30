package lego

type Equaler[T any] interface {
	Equal(T) bool
}
