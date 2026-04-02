package v2

type Viewer[V any] interface {
	View() V
}
