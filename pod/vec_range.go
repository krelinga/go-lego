package pod

import "iter"

func VecRange[T any](parent VecView[T], fromIdx, toIdx int) VecView[T] {
	if fromIdx < 0 || toIdx > parent.Len() || fromIdx > toIdx {
		panic("invalid range")
	}
	return &vecRange[T]{
		parent: parent,
		fromIdx:   fromIdx,
		toIdx:     toIdx,
	}
}

func VecRangeFrom[T any](parent VecView[T], fromIdx int) VecView[T] {
	return VecRange(parent, fromIdx, parent.Len())
}

func VecRangeTo[T any](parent VecView[T], toIdx int) VecView[T] {
	return VecRange(parent, 0, toIdx)
}

type vecRange[T any] struct {
	parent VecView[T]
	fromIdx   int
	toIdx     int
}

func (r vecRange[T]) checkParent() {
	if r.parent.Len() < r.toIdx {
		panic("parent vector is too short for range")
	}
}

func (r vecRange[T]) Len() int {
	r.checkParent()
	return r.toIdx - r.fromIdx
}

func (r vecRange[T]) Get(i int) T {
	r.checkParent()
	if i < 0 || i >= r.Len() {
		panic("index out of range")
	}
	return r.parent.Get(r.fromIdx + i)
}

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
