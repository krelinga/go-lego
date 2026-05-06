package pod

import (
	"cmp"
	"iter"
	"slices"

	"github.com/krelinga/go-libs/zero"
)

// FixedVec is a read-only view of a vector. It provides methods to access the elements, but does not allow mutation.
type FixedVec[T any] interface {
	// Len returns the number of elements in the FixedVec.
	Len() int

	// Get returns the element at the specified index in the FixedVec. If the index is out of bounds, it will panic.
	Get(i int) T

	// Vals returns an in-order sequence of values in the FixedVec.
	Vals() iter.Seq[T]

	// IdxVals returns an in-order sequence of indexed values in the FixedVec. Each element is a pair of the form (index, value), where index is the position of the value in the FixedVec and value is the corresponding element from the FixedVec.
	IdxVals() iter.Seq2[int, T]

	// RevVals returns a reverse-order sequence of values in the FixedVec.
	RevVals() iter.Seq[T]

	// RevIdxVals returns a reverse-order sequence of indexed values in the FixedVec. Each element is a pair of the form (index, value), where index is the position of the value in the FixedVec and value is the corresponding element from the FixedVec.
	RevIdxVals() iter.Seq2[int, T]
}

// Vec is a mutable vector that implements the FixedVec interface. It allows setting, clearing, pushing, popping, and resizing elements.
type Vec[T any] interface {
	FixedVec[T]

	// Set sets the element at the specified index in the Vec to the given value. If the index is out of bounds, it will panic.
	Set(i int, value T)

	// Clear removes all elements from the Vec, leaving it empty.
	Clear()

	// Push adds an element to the end of the Vec.
	Push(value T)

	// Pop removes and returns the last element of the Vec. If the Vec is empty, it will panic.
	Pop() T

	// Resize changes the length of the Vec to the specified new length. If the new length is greater than the current length, the new elements will be initialized to the zero value of the element type. If the new length is less than the current length, the Vec will be truncated.
	Resize(newLen int)
}

// AsVec creates a FixedVec from a slice. It provides a read-only view of the slice, and any changes to the slice will be reflected in the FixedVec.
func AsVec[V ~[]T, T any](v V) FixedVec[T] {
	return asVec[V, T]{v: v}
}

type asVec[V ~[]T, T any] struct {
	v V
}

func (v asVec[V, T]) Len() int {
	return len(v.v)
}

func (v asVec[V, T]) Get(i int) T {
	return v.v[i]
}

func (v asVec[V, T]) Vals() iter.Seq[T] {
	return slices.Values(v.v)
}

func (v asVec[V, T]) IdxVals() iter.Seq2[int, T] {
	return slices.All(v.v)
}

func (v asVec[V, T]) RevVals() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(v.v) - 1; i >= 0; i-- {
			if !yield(v.v[i]) {
				return
			}
		}
	}
}

func (v asVec[V, T]) RevIdxVals() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := len(v.v) - 1; i >= 0; i-- {
			if !yield(i, v.v[i]) {
				return
			}
		}
	}
}

// Slice is a mutable vector that implements the Vec interface. It is a wrapper around a slice, so it is possible to create literals of Slice as follows:
//
//	s := &Slice[int]{1, 2, 3}
type Slice[T any] []T

// NewSlice creates a new Slice with the specified length. The elements of the Slice will be initialized to the zero value of the element type.
func NewSlice[T any](len int) *Slice[T] {
	slice := make([]T, len)
	return (*Slice[T])(&slice)
}

// NewSliceCap creates a new Slice with the specified length and capacity. The elements of the Slice will be initialized to the zero value of the element type.
func NewSliceCap[T any](len, cap int) *Slice[T] {
	slice := make([]T, len, cap)
	return (*Slice[T])(&slice)
}

// NewSliceOf creates a new Slice from a Bag. It collects the values from the Bag into a slice and returns a pointer to the Slice.
// The resulting Slice will be independent of the Bag, so changes to the Bag will not affect the Slice, and vice versa.
func NewSliceOf[T any](b Bag[T]) *Slice[T] {
	slice := make([]T, b.Len())
	i := 0
	for v := range b.Elems() {
		slice[i] = v
		i++
	}
	return (*Slice[T])(&slice)
}

// NewSliceOfCap creates a new Slice from a Bag with the specified capacity. It collects the values from the Bag into a slice and returns a pointer to the Slice. The resulting Slice will have the specified capacity, which can help improve performance if you know in advance how many elements will be added to the Slice after it is created. The resulting Slice will be independent of the Bag, so changes to the Bag will not affect the Slice, and vice versa.
func NewSliceOfCap[T any](b Bag[T], cap int) *Slice[T] {
	slice := make([]T, b.Len(), cap)
	i := 0
	for v := range b.Elems() {
		slice[i] = v
		i++
	}
	return (*Slice[T])(&slice)
}

// Len returns the number of elements in the Slice.
func (v *Slice[T]) Len() int {
	return len(*v)
}

// Get returns the element at the specified index in the Slice. If the index is out of bounds, it will panic.
func (v *Slice[T]) Get(i int) T {
	return (*v)[i]
}

// Vals returns an in-order sequence of values in the Slice.
func (v *Slice[T]) Vals() iter.Seq[T] {
	return slices.Values(*v)
}

// IdxVals returns an in-order sequence of indexed values in the Slice. Each element is a pair of the form (index, value), where index is the position of the value in the Slice and value is the corresponding element from the Slice.
func (v *Slice[T]) IdxVals() iter.Seq2[int, T] {
	return slices.All(*v)
}

// RevVals returns a reverse-order sequence of values in the Slice.
func (v *Slice[T]) RevVals() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(*v) - 1; i >= 0; i-- {
			if !yield((*v)[i]) {
				return
			}
		}
	}
}

// RevIdxVals returns a reverse-order sequence of indexed values in the Slice. Each element is a pair of the form (index, value), where index is the position of the value in the Slice and value is the corresponding element from the Slice.
func (v *Slice[T]) RevIdxVals() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := len(*v) - 1; i >= 0; i-- {
			if !yield(i, (*v)[i]) {
				return
			}
		}
	}
}

// Set sets the element at the specified index in the Slice to the given value. If the index is out of bounds, it will panic.
func (v *Slice[T]) Set(i int, value T) {
	(*v)[i] = value
}

// Clear removes all elements from the Slice, leaving it empty.
func (v *Slice[T]) Clear() {
	*v = nil
}

// Reserve allows callers to pre-allocate space for a certain number of elements in an empty Slice.
// If the slice is already initialized and has a capacity of less than n, it will create a new underlying array with the specified capacity and copy the existing elements to it. If the slice is already initialized and has a capacity of n or more, it does nothing.
func (v *Slice[T]) Reserve(n int) {
	if cap(*v) < n {
		newData := make([]T, len(*v), n)
		copy(newData, *v)
		*v = newData
	}
}

// Push adds a value to the end of the Slice, increasing its length by 1. If the underlying array does not have enough capacity to accommodate the new element, it will create a new underlying array with increased capacity and copy the existing elements to it before adding the new element.
func (v *Slice[T]) Push(value T) {
	*v = append(*v, value)
}

// Resize changes the length of the Slice to newLen. If newLen is less than the current length, it truncates the Slice. If newLen is greater than the current length, it extends the Slice and initializes the new elements to the zero value of the element type. If newLen is equal to the current length, it does nothing.
func (v *Slice[T]) Resize(newLen int) {
	if newLen < len(*v) {
		*v = (*v)[:newLen]
	} else if newLen > len(*v) {
		v.Reserve(newLen)
		for i := len(*v); i < newLen; i++ {
			*v = append(*v, zero.For[T]())
		}
	}
}

// Pop removes and returns the last element of the Slice, decreasing its length by 1. If the Slice is empty, it will panic.
func (v *Slice[T]) Pop() T {
	value := (*v)[len(*v)-1]
	*v = (*v)[:len(*v)-1]
	return value
}

// WrapVecVals creates a new FixedVec that wraps the values of the given FixedVec with the provided wrap function. The resulting FixedVec will reflect any changes to the underlying FixedVec, and the wrap function will be applied to each value when accessed through the new FixedVec. The indices of the elements remain unchanged.
func WrapVecVals[T, V any](vec FixedVec[T], wrap func(T) V) FixedVec[V] {
	return wrappedVecVals[T, V]{
		vec:  vec,
		wrap: wrap,
	}
}

type wrappedVecVals[T, V any] struct {
	vec  FixedVec[T]
	wrap func(T) V
}

func (w wrappedVecVals[T, V]) Len() int {
	return w.vec.Len()
}

func (w wrappedVecVals[T, V]) Get(i int) V {
	return w.wrap(w.vec.Get(i))
}

func (w wrappedVecVals[T, V]) Vals() iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range w.vec.Vals() {
			if !yield(w.wrap(v)) {
				return
			}
		}
	}
}

func (w wrappedVecVals[T, V]) IdxVals() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i, v := range w.vec.IdxVals() {
			if !yield(i, w.wrap(v)) {
				return
			}
		}
	}
}

func (w wrappedVecVals[T, V]) RevVals() iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range w.vec.RevVals() {
			if !yield(w.wrap(v)) {
				return
			}
		}
	}
}

func (w wrappedVecVals[T, V]) RevIdxVals() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i, v := range w.vec.RevIdxVals() {
			if !yield(i, w.wrap(v)) {
				return
			}
		}
	}
}

// VecEqualFunc checks if two FixedVecs are equal by comparing their lengths and corresponding elements using the provided equality function. It returns true if the FixedVecs are of the same length and all corresponding elements are equal according to the equality function, and false otherwise.
func VecEqualFunc[T any](a, b FixedVec[T], eq func(T, T) bool) bool {
	if a.Len() != b.Len() {
		return false
	}
	for i := 0; i < a.Len(); i++ {
		if !eq(a.Get(i), b.Get(i)) {
			return false
		}
	}
	return true
}

// VecEqual checks if two FixedVecs are equal by comparing their lengths and corresponding elements using the equality operator (==). It returns true if the FixedVecs are of the same length and all corresponding elements are equal, and false otherwise. The element type must be comparable for this function to work.
func VecEqual[T comparable](a, b FixedVec[T]) bool {
	return VecEqualFunc(a, b, func(x, y T) bool {
		return x == y
	})
}

// VecSort sorts the elements of the given Vec in-place using the natural ordering of the element type. The element type must be ordered (i.e., it must satisfy the cmp.Ordered type set) for this function to work.
func VecSort[T cmp.Ordered](vec Vec[T]) {
	VecSortFunc(vec, cmp.Compare[T])
}

// VecSortFunc sorts the elements of the given Vec in-place using the provided comparison function. The order of the elements will be determined by the comparison function, which should return a negative value if the first argument is less than the second, zero if they are equal, and a positive value if the first argument is greater than the second.
func VecSortFunc[T any](vec Vec[T], order func(T, T) int) {
	if slice, ok := vec.(*Slice[T]); ok {
		slices.SortFunc(*slice, order)
		return
	}

	// TODO: make this more-efficient if/when we have non-Slice Vec implementations.
	temp := make([]T, vec.Len())
	for i := 0; i < vec.Len(); i++ {
		temp[i] = vec.Get(i)
	}
	slices.SortFunc(temp, order)
	for i := 0; i < vec.Len(); i++ {
		vec.Set(i, temp[i])
	}
}
