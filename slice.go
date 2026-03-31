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

// A Slice is a mutable slice type.
// It implements the [FixedSlice] interface and the [Adder] interface.
type Slice[T any] []T

func (s *Slice[T]) Len() int {
	return len(*s)
}

func (s *Slice[T]) List() iter.Seq[Pair[int, T]] {
	return func(yield func(Pair[int, T]) bool) {
		for i, v := range *s {
			if !yield(NewPair(i, v)) {
				return
			}
		}
	}
}

func (s *Slice[T]) Get(i int) T {
	return (*s)[i]
}

func (s *Slice[T]) Add(value T) {
	*s = append(*s, value)
}

// Reserve reserves space for n elements in the slice. This is a best-effort operation and will do nothing if the slice already contains some values, since Go's built-in slices do not support reserving space after initialization.
func (s *Slice[T]) Reserve(n int) {
	if *s == nil {
		*s = make([]T, 0, n)
	}
}

func NewSlice[T any](l int) *Slice[T] {
	s := make([]T, l)
	return (*Slice[T])(&s)
}

func NewSliceCap[T any](l, c int) *Slice[T] {
	s := make([]T, l, c)
	return (*Slice[T])(&s)
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

func DeepCopySlice[S FixedSlice[T], T DeepCopier[T]](s S) *Slice[T] {
	var out Slice[T]
	out.Reserve(s.Len())
	for v := range Values(s).List() {
		out.Add(v.DeepCopy())
	}
	return &out
}

func ShallowCopySlice[S FixedSlice[T], T any](s S) *Slice[T] {
	var out Slice[T]
	out.Reserve(s.Len())
	for v := range Values(s).List() {
		out.Add(v)
	}
	return &out
}

// Sort sorts the elements of the given [Slice] in place using the Compare method of the [Comparer] interface to determine the order of the elements.
func Sort[T Comparer[T]](s *Slice[T]) {
	slices.SortFunc(*s, func(a, b T) int {
		return a.Compare(b)
	})
}

// SortFunc sorts the elements of the given [Slice] in place using the given CmpFunc to determine the order of the elements.
func SortFunc[T any](s *Slice[T], compare CmpFunc[T]) {
	slices.SortFunc(*s, compare)
}

// SortGo sorts the elements of the given [Slice] in place using the natural order of the elements.
func SortGo[T cmp.Ordered](s *Slice[T]) {
	slices.Sort(*s)
}

// Reverse reverses the elements of the given [Slice] in place.
func Reverse[T any](s *Slice[T]) {
	slices.Reverse(*s)
}

func EqualSlice[S FixedSlice[T], T Equaler[T]](a, b S) bool {
	return equalSliceImpl(a, b, func(x, y T) bool {
		return x.Equal(y)
	})
}

func EqualSliceGo[S FixedSlice[T], T comparable](a, b S) bool {
	return equalSliceImpl(a, b, func(x, y T) bool {
		return x == y
	})
}

func equalSliceImpl[S FixedSlice[T], T any](a, b S, equal func(T, T) bool) bool {
	if a.Len() != b.Len() {
		return false
	}
	for i := 0; i < a.Len(); i++ {
		if !equal(a.Get(i), b.Get(i)) {
			return false
		}
	}
	return true
}