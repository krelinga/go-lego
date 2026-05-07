package view

type Viewer[T any] interface {
	View() T
}
