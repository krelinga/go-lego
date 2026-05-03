package pod

import (
	"iter"
)

// SliceSet is a Set implementation that preserves the order of insertion.
// Iteration over the set will yield elements in the order they were added.
// Set equality is determined by a user-provided equality function, allowing
// for custom comparison logic.
//
// Many of the operations on SliceSet require iterating through all elements
// to check for equality, resulting in O(n) time complexity for operations like
// Has, Put, and Del. This makes SliceSet less efficient than a hash-based set
// for large collections, but it provides the benefit of maintaining insertion
// order and allowing for custom equality logic.
type SliceSet[T any] struct {
	slice *Slice[T]
	equal func(a, b T) bool
}

// NewSliceSetFunc creates a new SliceSet with a custom equality function and initial values.
// If any of the initial values are considered equal according to the equality function, only the first occurrence will be added to the set.
func NewSliceSetFunc[T any](equal func(a, b T) bool, vals ...T) *SliceSet[T] {
	s := &SliceSet[T]{
		slice: &Slice[T]{},
		equal: equal,
	}
	s.Reserve(len(vals))
	for _, v := range vals {
		s.Put(v)
	}
	return s
}

// NewSliceSet creates a new SliceSet with the default equality function (using Go's built-in equality) and initial values.
// If any of the initial values are considered equal, only the first occurrence will be added to the set.
func NewSliceSet[T comparable](vals ...T) *SliceSet[T] {
	return NewSliceSetFunc(func(a, b T) bool { return a == b }, vals...)
}

// NewSliceSetOfFunc creates a new SliceSet from a Bag, using a custom equality function.
// The set will contain all unique elements from the Bag, determined by the equality function.
// The order of elements in the set will match the order they appear in the Bag.
// The new Set will be independent of the original Bag; changes to the Bag's underlying collection will not affect the Set, and vice versa.
func NewSliceSetOfFunc[T any](equal func(a, b T) bool, b Bag[T]) *SliceSet[T] {
	s := NewSliceSetFunc(equal)
	s.Reserve(b.Len())
	for v := range b.Elems() {
		s.Put(v)
	}
	return s
}

// NewSliceSetOf creates a new SliceSet from a Bag, using the default equality function (using Go's built-in equality).
// The set will contain all unique elements from the Bag.
// The order of elements in the set will match the order they appear in the Bag.
// The new Set will be independent of the original Bag; changes to the Bag's underlying collection will not affect the Set, and vice versa.
func NewSliceSetOf[T comparable](b Bag[T]) *SliceSet[T] {
	return NewSliceSetOfFunc(func(a, b T) bool { return a == b }, b)
}

// Len returns the number of unique elements in the SliceSet.
// This is a constant-time operation.
func (s *SliceSet[T]) Len() int {
	return s.slice.Len()
}

// Has checks if the given value is present in the SliceSet.
// This operation has a time complexity of O(n) in the worst case, as it may need to iterate through all elements to find a match.
func (s *SliceSet[T]) Has(value T) bool {
	for v := range s.slice.Vals() {
		if s.equal(v, value) {
			return true
		}
	}
	return false
}

// Vals returns a sequence of all unique elements in the SliceSet, in the order they were added.
func (s *SliceSet[T]) Vals() iter.Seq[T] {
	return s.slice.Vals()
}

// Put adds a value to the SliceSet if it is not already present.
// If the value is already in the set (as determined by the equality function), it will not be added again, preserving the uniqueness of elements in the set.
// This operation has a time complexity of O(n) in the worst case, as it may need to iterate through all elements to check for duplicates before adding a new value.
func (s *SliceSet[T]) Put(value T) {
	if s.Has(value) {
		return
	}
	s.slice.Push(value)
}

// Clear removes all elements from the SliceSet, resulting in an empty set.
func (s *SliceSet[T]) Clear() {
	s.slice.Clear()
}

// Del removes a value from the SliceSet if it is present.
// If the value is not in the set, this operation has no effect.
// This operation has a time complexity of O(n) in all cases.
func (s *SliceSet[T]) Del(value T) {
	for i := range s.slice.Len() {
		if s.equal(s.slice.Get(i), value) {
			for j := i + 1; j < s.slice.Len(); j++ {
				s.slice.Set(j-1, s.slice.Get(j))
			}
			s.slice.Pop()
			return
		}
	}
}

// Reserve pre-allocates space for at least n elements in the SliceSet.
// This can improve performance when adding a large number of elements, as it reduces the need for multiple allocations as the set grows.
func (s *SliceSet[T]) Reserve(n int) {
	s.slice.Reserve(n)
}
