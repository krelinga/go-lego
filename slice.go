package lego

import "iter"

// A FixedSlice is a slice that does not allow adding or removing elements, but which may still allow modifying the elements in the slice (for example, if the elements are pointers).
type FixedSlice[T any] interface {
	Len() int
	List() iter.Seq[Pair[int, T]]

	Get(int) T
}

// A GoSlice is a wrapper around Go's built-in slice type that implements the [FixedSlice] interface.
// It does not implement the [Adder] interface, since Go's built-in slices do not allow adding elements without creating a new slice.
type GoSlice[T any] []T

func (s GoSlice[T]) Len() int {
	return len(s)
}

func (s GoSlice[T]) List() iter.Seq[Pair[int, T]] {
	return func(yield func(Pair[int, T]) bool) {
		for i, v := range s {
			if !yield(NewPair(i, v)) {
				return
			}
		}
	}
}

func (s GoSlice[T]) Get(i int) T {
	return s[i]
}

// A Slice is a mutable slice type.
// It implements the [FixedSlice] interface and the [Adder] interface.
type Slice[T any] struct {
	s GoSlice[T]
}

func (s Slice[T]) Len() int {
	return len(s.s)
}

func (s Slice[T]) Add(value T) {
	s.s = append(s.s, value)
}

func (s Slice[T]) List() iter.Seq[Pair[int, T]] {
	return s.s.List()
}

func (s Slice[T]) Get(i int) T {
	return s.s.Get(i)
}

func NewSlice[S ~[]T, T any](slice S) Slice[T] {
	return Slice[T]{s: GoSlice[T](slice)}
}

// ViewSlice creates a view of a slice that allows viewing the elements of the slice as a different type, without modifying the original slice.
func ViewSlice[S FixedSlice[V1], V1 Viewer[V2], V2 any](s S) FixedSlice[V2] {
	return sliceView[S, V1, V2]{s: s}
}

type sliceView[S FixedSlice[V1], V1 Viewer[V2], V2 any] struct {
	s S
}

func (v sliceView[S, V1, V2]) Len() int {
	return v.s.Len()
}

func (v sliceView[S, V1, V2]) List() iter.Seq[Pair[int, V2]] {
	return func(yield func(Pair[int, V2]) bool) {
		for p := range v.s.List() {
			if !yield(NewPair(p.GetKey(), p.GetValue().View())) {
				return
			}
		}
	}
}

func (v sliceView[S, V1, V2]) Get(i int) V2 {
	return v.s.Get(i).View()
}

func DeepCopySlice[S FixedSlice[T], T DeepCopier[T]](s S) Slice[T] {
	var out Slice[T]
	for p := range s.List() {
		out.Add(p.GetValue().DeepCopy())
	}
	return out
}
