package pod

import "iter"

func VecRange[T any](parent VecView[T], from, to int) VecView[T] {
	if from < 0 || to > parent.Len() || from > to {
		panic("invalid range")
	}
	return &vecRange[T]{
		parent: parent,
		from:   from,
		to:     to,
	}
}

func VecRangeFrom[T any](parent VecView[T], from int) VecView[T] {
	return VecRange(parent, from, parent.Len())
}

func VecRangeTo[T any](parent VecView[T], to int) VecView[T] {
	return VecRange(parent, 0, to)
}

type vecRange[T any] struct {
	parent VecView[T]
	from   int
	to     int
}

func (r vecRange[T]) checkParent() {
	if r.parent.Len() < r.to {
		panic("parent vector is too short for range")
	}
}

func (r vecRange[T]) Len() int {
	r.checkParent()
	return r.to - r.from
}

func (r vecRange[T]) At(i int) T {
	r.checkParent()
	if i < 0 || i >= r.Len() {
		panic("index out of range")
	}
	return r.parent.At(r.from + i)
}

func (r vecRange[T]) Vals() iter.Seq[T] {
	r.checkParent()
	return func(yield func(T) bool) {
		for i := r.from; i < r.to; i++ {
			if !yield(r.parent.At(i)) {
				return
			}
		}
	}
}

func (r vecRange[T]) IdxVals() iter.Seq2[int, T] {
	r.checkParent()
	return func(yield func(int, T) bool) {
		for i := r.from; i < r.to; i++ {
			if !yield(i-r.from, r.parent.At(i)) {
				return
			}
		}
	}
}

func (r vecRange[T]) RevVals() iter.Seq[T] {
	r.checkParent()
	return func(yield func(T) bool) {
		for i := r.to - 1; i >= r.from; i-- {
			if !yield(r.parent.At(i)) {
				return
			}
		}
	}
}

func (r vecRange[T]) RevIdxVals() iter.Seq2[int, T] {
	r.checkParent()
	return func(yield func(int, T) bool) {
		for i := r.to - 1; i >= r.from; i-- {
			if !yield(i-r.from, r.parent.At(i)) {
				return
			}
		}
	}
}
