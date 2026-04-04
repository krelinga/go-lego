package v2

import "iter"

type Lengther interface {
	Length() int
}

type Range[V any] interface {
	Length() int
	Iterate() iter.Seq[V]
}

type rangeImpl[C Lengther, V any] struct {
	container C
	iterFunc func(C) iter.Seq[V]
}

func (r rangeImpl[C, V]) Length() int {
	return r.container.Length()
}

func (r rangeImpl[C, V]) Iterate() iter.Seq[V] {
	return r.iterFunc(r.container)
}

func NewRange[C Lengther, V any](container C, iterFunc func(C) iter.Seq[V]) Range[V] {
	return rangeImpl[C, V]{container, iterFunc}
}
