package view

type Can[T any] interface {
	View() T
}

func Direct[T any](v T) Can[T] {
	return direct[T]{v: v}
}

type direct[T any] struct {
	v T
}

func (d direct[T]) View() T {
	return d.v
}