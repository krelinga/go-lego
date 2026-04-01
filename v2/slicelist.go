package v2

import (
	"iter"
	"slices"
)

type sliceList[V any] struct {
	slice []V
}

func (l *sliceList[V]) Length() int {
	return len(l.slice)
}

func (l *sliceList[V]) Get(k int) (V, bool) {
	if k < 0 || k >= len(l.slice) {
		var zero V
		return zero, false
	}
	return l.slice[k], true
}

func (l *sliceList[V]) All() iter.Seq2[int, V] {
	return slices.All(l.slice)
}

func (l *sliceList[V]) First() (int, bool) {
	return 0, len(l.slice) > 0
}

func (l *sliceList[V]) Last() (int, bool) {
	if len(l.slice) == 0 {
		return 0, false
	}
	return len(l.slice) - 1, true
}

func (l *sliceList[V]) InsertBefore(k int, v V) {
	if k < 0 || k > len(l.slice) {
		panic("index out of bounds")
	}
	l.slice = append(l.slice[:k], append([]V{v}, l.slice[k:]...)...)
}

func (l *sliceList[V]) InsertAfter(k int, v V) {
	if k < 0 || k >= len(l.slice) {
		panic("index out of bounds")
	}
	l.slice = append(l.slice[:k+1], append([]V{v}, l.slice[k+1:]...)...)
}

func (l *sliceList[V]) Remove(k int) {
	if k < 0 || k >= len(l.slice) {
		panic("index out of bounds")
	}
	l.slice = append(l.slice[:k], l.slice[k+1:]...)
}

func NewSliceList[V any](entries ...V) List[int, V] {
	return &sliceList[V]{
		slice: entries,
	}
}
