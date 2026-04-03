package v2

import (
	"iter"
	"slices"
)

type SliceList[V any] struct {
	slice []V
}

func (l *SliceList[V]) Length() int {
	return len(l.slice)
}

func (l *SliceList[V]) Get(p int) (V, bool) {
	if p < 0 || p >= len(l.slice) {
		var zero V
		return zero, false
	}
	return l.slice[p], true
}

func (l *SliceList[V]) All() iter.Seq2[int, V] {
	return slices.All(l.slice)
}

func (l *SliceList[V]) Positions() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range l.slice {
			if !yield(i) {
				return
			}
		}
	}
}

func (l *SliceList[V]) Values() iter.Seq[V] {
	return slices.Values(l.slice)
}

func (l *SliceList[V]) ReverseAll() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i := len(l.slice) - 1; i >= 0; i-- {
			if !yield(i, l.slice[i]) {
				return
			}
		}
	}
}

func (l *SliceList[V]) ReversePositions() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := len(l.slice) - 1; i >= 0; i-- {
			if !yield(i) {
				return
			}
		}
	}
}

func (l *SliceList[V]) ReverseValues() iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := len(l.slice) - 1; i >= 0; i-- {
			if !yield(l.slice[i]) {
				return
			}
		}
	}
}

func (l *SliceList[V]) First() (int, V, bool) {
	if len(l.slice) == 0 {
		var zero V
		return 0, zero, false
	}
	return 0, l.slice[0], true
}

func (l *SliceList[V]) Last() (int, V, bool) {
	if len(l.slice) == 0 {
		var zero V
		return 0, zero, false
	}
	lastIdx := len(l.slice) - 1
	return lastIdx, l.slice[lastIdx], true
}

func (l *SliceList[V]) Set(p int, v V) {
	if p < 0 || p >= len(l.slice) {
		panic("index out of bounds")
	}
	l.slice[p] = v
}

func (l *SliceList[V]) InsertBefore(p int, v V) int {
	if p < 0 || p > len(l.slice) {
		panic("index out of bounds")
	}
	l.slice = append(l.slice[:p], append([]V{v}, l.slice[p:]...)...)
	return p
}

func (l *SliceList[V]) InsertAfter(p int, v V) int {
	if p < 0 || p >= len(l.slice) {
		panic("index out of bounds")
	}
	l.slice = append(l.slice[:p+1], append([]V{v}, l.slice[p+1:]...)...)
	return p + 1
}

func (l *SliceList[V]) Append(v V) int {
	l.slice = append(l.slice, v)
	return len(l.slice) - 1
}

func (l *SliceList[V]) Remove(p int) {
	if p < 0 || p >= len(l.slice) {
		panic("index out of bounds")
	}
	l.slice = append(l.slice[:p], l.slice[p+1:]...)
}

func (l *SliceList[V]) Reserve(n int) {
	if l.slice == nil {
		l.slice = make([]V, 0, n)
	}
}

func (l *SliceList[V]) Add(v V) {
	l.Append(v)
}

func NewSliceList[V any](entries ...V) *SliceList[V] {
	return &SliceList[V]{
		slice: entries,
	}
}
