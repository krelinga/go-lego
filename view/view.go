package view

type Viewer[T any] interface {
	View() T
}

func Direct[T any](v T) Viewer[T] {
	return direct[T]{v: v}
}

type direct[T any] struct {
	v T
}

func (d direct[T]) View() T {
	return d.v
}
