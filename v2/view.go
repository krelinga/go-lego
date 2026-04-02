package v2

type Viewer[V any] interface {
	View() V
}

func ViewHelperListOfViewers[P any, T Viewer[V], V any](in FixedList[P, T]) FixedList[P, V] {
	return viewHelperListOfViewers[P, T, V]{
		FixedList: in,
	}
}

type viewHelperListOfViewers[P any, T Viewer[V], V any] struct {
	FixedList[P, T]
}

func (h viewHelperListOfViewers[P, T, V]) Get(p P) (V, bool) {
	viewer, ok := h.FixedList.Get(p)
	if !ok {
		var zero V
		return zero, false
	}
	return viewer.View(), true
}

func (h viewHelperListOfViewers[P, T, V]) List() ListSeq[P, V] {
	return func(yield func(P, V) bool) {
		for p, viewer := range h.FixedList.List() {
			if !yield(p, viewer.View()) {
				return
			}
		}
	}
}

func (h viewHelperListOfViewers[P, T, V]) First() (P, V, bool) {
	p, viewer, ok := h.FixedList.First()
	if !ok {
		var zero V
		return p, zero, false
	}
	return p, viewer.View(), true
}

func (h viewHelperListOfViewers[P, T, V]) Last() (P, V, bool) {
	p, viewer, ok := h.FixedList.Last()
	if !ok {
		var zero V
		return p, zero, false
	}
	return p, viewer.View(), true
}
