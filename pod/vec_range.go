package pod

import "iter"

// VecRange creates a new FixedVec that represents a contiguous range of elements from the parent FixedVec, starting from fromIdx (inclusive) and ending at toIdx (exclusive). The resulting FixedVec will reflect any changes made to the parent FixedVec within the specified range. If the provided indices are out of bounds or if fromIdx is greater than toIdx, it will panic.
func VecRange[T any](parent FixedVec[T], fromIdx, toIdx int) FixedVec[T] {
	if fromIdx < 0 || toIdx > parent.Len() || fromIdx > toIdx {
		panic("invalid range")
	}
	return &vecRange[T]{
		parent:  parent,
		fromIdx: fromIdx,
		toIdx:   toIdx,
	}
}

// VecRangeFrom creates a new FixedVec that represents a contiguous range of elements from the parent FixedVec, starting from fromIdx (inclusive) and ending at the end of the parent FixedVec. The resulting FixedVec will reflect any changes made to the parent FixedVec within the specified range. If fromIdx is out of bounds, it will panic.
func VecRangeFrom[T any](parent FixedVec[T], fromIdx int) FixedVec[T] {
	return VecRange(parent, fromIdx, parent.Len())
}

// VecRangeTo creates a new FixedVec that represents a contiguous range of elements from the parent FixedVec, starting from the beginning of the parent FixedVec and ending at toIdx (exclusive). The resulting FixedVec will reflect any changes made to the parent FixedVec within the specified range. If toIdx is out of bounds, it will panic.
func VecRangeTo[T any](parent FixedVec[T], toIdx int) FixedVec[T] {
	return VecRange(parent, 0, toIdx)
}

type vecRange[T any] struct {
	parent  FixedVec[T]
	fromIdx int
	toIdx   int
}

func (r vecRange[T]) checkParent() {
	if r.parent.Len() < r.toIdx {
		panic("parent vector is too short for range")
	}
}

// Len returns the number of elements in the VecRange, which is the difference between toIdx and fromIdx. If the parent FixedVec has been modified such that it is shorter than toIdx, it will panic.
func (r vecRange[T]) Len() int {
	r.checkParent()
	return r.toIdx - r.fromIdx
}

// Get returns the element at the specified index within the VecRange. The index is relative to the start of the range (i.e., index 0 corresponds to fromIdx in the parent FixedVec). If the index is out of bounds (less than 0 or greater than or equal to the length of the range), it will panic. If the parent FixedVec has been modified such that it is shorter than toIdx, it will panic.
func (r vecRange[T]) Get(i int) T {
	r.checkParent()
	if i < 0 || i >= r.Len() {
		panic("index out of range")
	}
	return r.parent.Get(r.fromIdx + i)
}

// Vals returns a sequence of values in the VecRange. The order of the values is the same as their order in the parent FixedVec. If the parent FixedVec has been modified such that it is shorter than toIdx, it will panic.
func (r vecRange[T]) Vals() iter.Seq[T] {
	r.checkParent()
	return func(yield func(T) bool) {
		for i := r.fromIdx; i < r.toIdx; i++ {
			if !yield(r.parent.Get(i)) {
				return
			}
		}
	}
}

// IdxVals returns a sequence of indexed values in the VecRange. Each element is a pair of the form (index, value), where index is the position of the value within the range (starting from 0) and value is the corresponding element from the parent FixedVec. The order of the indexed values is the same as their order in the parent FixedVec. If the parent FixedVec has been modified such that it is shorter than toIdx, it will panic.
func (r vecRange[T]) IdxVals() iter.Seq2[int, T] {
	r.checkParent()
	return func(yield func(int, T) bool) {
		for i := r.fromIdx; i < r.toIdx; i++ {
			if !yield(i-r.fromIdx, r.parent.Get(i)) {
				return
			}
		}
	}
}

// RevVals returns a sequence of values in the VecRange in reverse order. The order of the values is the reverse of their order in the parent FixedVec. If the parent FixedVec has been modified such that it is shorter than toIdx, it will panic.
func (r vecRange[T]) RevVals() iter.Seq[T] {
	r.checkParent()
	return func(yield func(T) bool) {
		for i := r.toIdx - 1; i >= r.fromIdx; i-- {
			if !yield(r.parent.Get(i)) {
				return
			}
		}
	}
}

// RevIdxVals returns a sequence of indexed values in the VecRange in reverse order. Each element is a pair of the form (index, value), where index is the position of the value within the range (starting from 0) and value is the corresponding element from the parent FixedVec. The order of the indexed values is the reverse of their order in the parent FixedVec. If the parent FixedVec has been modified such that it is shorter than toIdx, it will panic.
func (r vecRange[T]) RevIdxVals() iter.Seq2[int, T] {
	r.checkParent()
	return func(yield func(int, T) bool) {
		for i := r.toIdx - 1; i >= r.fromIdx; i-- {
			if !yield(i-r.fromIdx, r.parent.Get(i)) {
				return
			}
		}
	}
}
