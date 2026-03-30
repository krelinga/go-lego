package lego

type Viewer[T any] interface {
	View() T
}

func EqualViewer[V Viewer[T], T Equaler[T]](a, b V) bool {
	return a.View().Equal(b.View())
}