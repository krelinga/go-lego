package lego

type Viewer[T any] interface {
	View() T
}