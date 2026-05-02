package pod

import (
	"iter"
)

// InOrderSet is a Set implementation that preserves the order of insertion.
// Iteration over the set will yield elements in the order they were added.
// Set equality is determined by a user-provided equality function, allowing for custom comparison logic.
type InOrderSet[T any] struct {
	slice *Slice[T]
	equal func(a, b T) bool
}

// NewInOrderSetFunc creates a new InOrderSet with a custom equality function and initial values.
// If any of the initial values are considered equal according to the equality function, only the first occurrence will be added to the set.
func NewInOrderSetFunc[T any](equal func(a, b T) bool, vals ...T) *InOrderSet[T] {
	s := &InOrderSet[T]{
		slice: &Slice[T]{},
		equal: equal,
	}
	s.Reserve(len(vals))
	for _, v := range vals {
		s.Put(v)
	}
	return s
}

// NewInOrderSet creates a new InOrderSet with the default equality function (using Go's built-in equality) and initial values.
// If any of the initial values are considered equal, only the first occurrence will be added to the set.
func NewInOrderSet[T comparable](vals ...T) *InOrderSet[T] {
	return NewInOrderSetFunc(func(a, b T) bool { return a == b }, vals...)
}

// NewInOrderSetOfFunc creates a new InOrderSet from a Bag, using a custom equality function.
// The set will contain all unique elements from the Bag, determined by the equality function.
// The order of elements in the set will match the order they appear in the Bag.
// The new Set will be independent of the original Bag; changes to the Bag's underlying collection will not affect the Set, and vice versa.
func NewInOrderSetOfFunc[T any](equal func(a, b T) bool, b Bag[T]) *InOrderSet[T] {
	s := NewInOrderSetFunc(equal)
	s.Reserve(b.Len())
	for v := range b.Elems() {
		s.Put(v)
	}
	return s
}

// NewInOrderSetOf creates a new InOrderSet from a Bag, using the default equality function (using Go's built-in equality).
// The set will contain all unique elements from the Bag.
// The order of elements in the set will match the order they appear in the Bag.
// The new Set will be independent of the original Bag; changes to the Bag's underlying collection will not affect the Set, and vice versa.
func NewInOrderSetOf[T comparable](b Bag[T]) *InOrderSet[T] {
	return NewInOrderSetOfFunc(func(a, b T) bool { return a == b }, b)
}

// Len returns the number of unique elements in the InOrderSet.
// This is a constant-time operation.
func (s *InOrderSet[T]) Len() int {
	return s.slice.Len()
}

// Has checks if the given value is present in the InOrderSet.
// This operation has a time complexity of O(n) in the worst case, as it may need to iterate through all elements to find a match.
func (s *InOrderSet[T]) Has(value T) bool {
	for v := range s.slice.Vals() {
		if s.equal(v, value) {
			return true
		}
	}
	return false
}

// Vals returns a sequence of all unique elements in the InOrderSet, in the order they were added.
func (s *InOrderSet[T]) Vals() iter.Seq[T] {
	return s.slice.Vals()
}

// Put adds a value to the InOrderSet if it is not already present.
// If the value is already in the set (as determined by the equality function), it will not be added again, preserving the uniqueness of elements in the set.
// This operation has a time complexity of O(n) in the worst case, as it may need to iterate through all elements to check for duplicates before adding a new value.
func (s *InOrderSet[T]) Put(value T) {
	if s.Has(value) {
		return
	}
	s.slice.Push(value)
}

// Clear removes all elements from the InOrderSet, resulting in an empty set.
func (s *InOrderSet[T]) Clear() {
	s.slice.Clear()
}

// Del removes a value from the InOrderSet if it is present.
// If the value is not in the set, this operation has no effect.
// This operation has a time complexity of O(n) in the worst case, as it may need to iterate through all elements to find the value to remove.
func (s *InOrderSet[T]) Del(value T) {
	// TODO: implement.
}

// Reserve pre-allocates space for at least n elements in the InOrderSet.
// This can improve performance when adding a large number of elements, as it reduces the need for multiple allocations as the set grows.
func (s *InOrderSet[T]) Reserve(n int) {
	s.slice.Reserve(n)
}
