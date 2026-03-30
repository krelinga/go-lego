package lego

import (
	"cmp"
	"iter"
	"slices"
)

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
	GoSlice[T]
}

func (s *Slice[T]) Add(value T) {
	s.GoSlice = append(s.GoSlice, value)
}

// Reserve reserves space for n elements in the slice. This is a best-effort operation and will do nothing if the slice already contains some values, since Go's built-in slices do not support reserving space after initialization.
func (s *Slice[T]) Reserve(n int) {
	if s.GoSlice == nil {
		s.GoSlice = make(GoSlice[T], 0, n)
	}
}

func NewSlice[T any](elements ...T) *Slice[T] {
	return &Slice[T]{GoSlice: GoSlice[T](elements)}
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
	out.Reserve(s.Len())
	for v := range Values(s).List() {
		out.Add(v.DeepCopy())
	}
	return out
}

func ShallowCopySlice[S FixedSlice[T], T any](s S) Slice[T] {
	var out Slice[T]
	out.Reserve(s.Len())
	for v := range Values(s).List() {
		out.Add(v)
	}
	return out
}

// Sort sorts the elements of the given [Slice] in place using the Compare method of the [Comparer] interface to determine the order of the elements.
func Sort[T Comparer[T]](s *Slice[T]) {
	slices.SortFunc(s.GoSlice, func(a, b T) int {
		return a.Compare(b)
	})
}

// SortFunc sorts the elements of the given [Slice] in place using the given CmpFunc to determine the order of the elements.
func SortFunc[T any](s *Slice[T], compare CmpFunc[T]) {
	slices.SortFunc(s.GoSlice, compare)
}

// SortGo sorts the elements of the given [Slice] in place using the natural order of the elements, which must implement the [cmp.Ordered] interface.
func SortGo[T cmp.Ordered](s *Slice[T]) {
	slices.Sort(s.GoSlice)
}

// Reverse reverses the elements of the given [Slice] in place.
func Reverse[T any](s *Slice[T]) {
	slices.Reverse(s.GoSlice)
}
