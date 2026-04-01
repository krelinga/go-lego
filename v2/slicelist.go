package v2

import (
	"iter"
	"slices"
)

type SliceListView[V any] interface {
	ListView[int, V]
	ReverseList() iter.Seq2[int, V]
}

type SliceList[V any] struct {
	slice []V
}

func (l *SliceList[V]) Length() int {
	return len(l.slice)
}

func (l *SliceList[V]) Get(k int) (V, bool) {
	if k < 0 || k >= len(l.slice) {
		var zero V
		return zero, false
	}
	return l.slice[k], true
}

func (l *SliceList[V]) List() iter.Seq2[int, V] {
	return slices.All(l.slice)
}

func (l *SliceList[V]) ReverseList() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i := len(l.slice) - 1; i >= 0; i-- {
			if !yield(i, l.slice[i]) {
				return
			}
		}
	}
}

func (l *SliceList[V]) First() (int, bool) {
	return 0, len(l.slice) > 0
}

func (l *SliceList[V]) Last() (int, bool) {
	if len(l.slice) == 0 {
		return 0, false
	}
	return len(l.slice) - 1, true
}

func (l *SliceList[V]) Set(k int, v V) {
	if k < 0 || k >= len(l.slice) {
		panic("index out of bounds")
	}
	l.slice[k] = v
}

func (l *SliceList[V]) InsertBefore(k int, v V) {
	if k < 0 || k > len(l.slice) {
		panic("index out of bounds")
	}
	l.slice = append(l.slice[:k], append([]V{v}, l.slice[k:]...)...)
}

func (l *SliceList[V]) InsertAfter(k int, v V) {
	if k < 0 || k >= len(l.slice) {
		panic("index out of bounds")
	}
	l.slice = append(l.slice[:k+1], append([]V{v}, l.slice[k+1:]...)...)
}

func (l *SliceList[V]) Append(v V) {
	l.slice = append(l.slice, v)
}

func (l *SliceList[V]) Remove(k int) {
	if k < 0 || k >= len(l.slice) {
		panic("index out of bounds")
	}
	l.slice = append(l.slice[:k], l.slice[k+1:]...)
}

func NewSliceList[V any](entries ...V) *SliceList[V] {
	return &SliceList[V]{
		slice: entries,
	}
}
