package v2

import (
	"iter"
	"slices"
)

type sliceList[V1, V2 any] struct {
	slice []V1
	valueReadonly ReadonlyFunc[V1, V2]
}

func (l *sliceList[V1, V2]) Length() int {
	return len(l.slice)
}

func (l *sliceList[V1, V2]) Get(k int) (V1, bool) {
	if k < 0 || k >= len(l.slice) {
		var zero V1
		return zero, false
	}
	return l.slice[k], true
}

func (l *sliceList[V1, V2]) All() iter.Seq2[int, V1] {
	return slices.All(l.slice)
}

func (l *sliceList[V1, V2]) First() (int, bool) {
	return 0, len(l.slice) > 0
}

func (l *sliceList[V1, V2]) Last() (int, bool) {
	if len(l.slice) == 0 {
		return 0, false
	}
	return len(l.slice) - 1, true
}

func (l *sliceList[V1, V2]) Readonly() ListView[int, V2, int, V2] {
	return sliceListReadonly[V1, V2]{
		sliceList: l,
	}
}

func (l *sliceList[V1, V2]) InsertBefore(k int, v V1) {
	if k < 0 || k > len(l.slice) {
		panic("index out of bounds")
	}
	l.slice = append(l.slice[:k], append([]V1{v}, l.slice[k:]...)...)
}

func (l *sliceList[V1, V2]) InsertAfter(k int, v V1) {
	if k < 0 || k >= len(l.slice) {
		panic("index out of bounds")
	}
	l.slice = append(l.slice[:k+1], append([]V1{v}, l.slice[k+1:]...)...)
}

func (l *sliceList[V1, V2]) Remove(k int) {
	if k < 0 || k >= len(l.slice) {
		panic("index out of bounds")
	}
	l.slice = append(l.slice[:k], l.slice[k+1:]...)
}

func NewSliceList[V1 any](entries ...V1) List[int, V1, int, V1] {
	return &sliceList[V1, V1]{
		slice: entries,
		valueReadonly: Identity[V1],
	}
}

func NewSliceListReadonly[V1, V2 any](valueReadonly ReadonlyFunc[V1, V2], entries ...V1) List[int, V1, int, V2] {
	return &sliceList[V1, V2]{
		slice: entries,
		valueReadonly: valueReadonly,
	}
}

type sliceListReadonly[V1, V2 any] struct {
	*sliceList[V1, V2]
}

func (l sliceListReadonly[V1, V2]) Get(k int) (V2, bool) {
	v1, ok := l.sliceList.Get(k)
	if !ok {
		var zero V2
		return zero, false
	}
	return l.valueReadonly(v1), true
}

func (l sliceListReadonly[V1, V2]) All() iter.Seq2[int, V2] {
	return func(yield func(int, V2) bool) {
		for i, v1 := range l.sliceList.slice {
			v2 := l.valueReadonly(v1)
			if !yield(i, v2) {
				return
			}
		}
	}
}

func (l sliceListReadonly[V1, V2]) Readonly() ListView[int, V2, int, V2] {
	return l
}