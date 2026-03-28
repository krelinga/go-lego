package lego

import "iter"

// A FixedSlice is a slice that does not allow adding or removing elements, but which may still allow modifying the elements in the slice (for example, if the elements are pointers).
type FixedSlice[T any] interface {
	Len() int
	List() iter.Seq[Pair[int, T]]

	Get(int) T
}

// A Slice is a mutable slice that wraps Go's built-in slice type.
// It implements the [FixedSlice] interface.
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

func (s Slice[T]) Get(i int) T {
	return s[i]
}

func NewSlice[S ~[]T, T any](slice S) Slice[T] {
	return Slice[T](slice)
}

// A ViewerSlice is a special case of [Slice] that stores values that implement the [Viewer] interface, and provides a method to get a view of the slice.
type ViewerSlice[V1 Viewer[V2], V2 any] Slice[V1]

func (s ViewerSlice[V1, V2]) Get(i int) V1 {
	return Slice[V1](s).Get(i)
}

func (s ViewerSlice[V1, V2]) List() iter.Seq[Pair[int, V1]] {
	return Slice[V1](s).List()
}

func (s ViewerSlice[V1, V2]) Len() int {
	return Slice[V1](s).Len()
}

func (s ViewerSlice[V1, V2]) View() FixedSlice[V2] {
	return viewerSliceView[V1, V2]{s: s}
}

type viewerSliceView[V1 Viewer[V2], V2 any] struct {
	s ViewerSlice[V1, V2]
}

func (v viewerSliceView[V1, V2]) Len() int {
	return v.s.Len()
}

func (v viewerSliceView[V1, V2]) List() iter.Seq[Pair[int, V2]] {
	return func(yield func(Pair[int, V2]) bool) {
		for i, v1 := range v.s {
			if !yield(NewPair(i, v1.View())) {
				return
			}
		}
	}
}

func (v viewerSliceView[V1, V2]) Get(i int) V2 {
	return v.s.Get(i).View()
}

func NewViewerSlice[S ~[]V1, V1 Viewer[V2], V2 any](slice S) ViewerSlice[V1, V2] {
	return ViewerSlice[V1, V2](slice)
}