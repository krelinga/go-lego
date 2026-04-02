package v2

type Viewer[V any] interface {
	View() V
}

type ViewMixinTrivial[P, V any] struct {
	fixedList FixedList[P, V]
}

func (m ViewMixinTrivial[P, V]) View() FixedList[P, V] {
	return m.fixedList
}

type ViewMixinListOfDirectValues[P, V any, FL FixedList[P, V]] func() FL

func (m ViewMixinListOfDirectValues[P, V, FL]) View() FL {
	return m()
}
