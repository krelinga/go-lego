package lego

import "iter"

type SliceView[T any] interface {
	Len() int
	List() iter.Seq[Pair[int, T]]

	Get(int) T
}

type Slice[T any] []T

func (s Slice[T]) Len() int {
	return len(s)
}

func (s Slice[T]) List() iter.Seq[Pair[int, T]] {
	return func(yield func(Pair[int, T]) bool) {
		for i, v := range s {
			if !yield(NewPair(i, v)) {
				return
			}
		}
	}
}

func NewSlice[S ~[]T, T any](slice S) Slice[T] {
	return Slice[T](slice)
}

type SliceViewer[V1 Viewer[V2], V2 any] Slice[V1]

func (s SliceViewer[V1, V2]) Get(i int) V2 {
	return s[i].View()
}

func (s SliceViewer[V1, V2]) List() iter.Seq[Pair[int, V2]] {
	return func(yield func(Pair[int, V2]) bool) {
		for i, v := range s {
			if !yield(NewPair(i, v.View())) {
				return
			}
		}
	}
}

func NewSliceViewer[S ~[]V1, V1 Viewer[V2], V2 any](slice S) SliceViewer[V1, V2] {
	return SliceViewer[V1, V2](slice)
}