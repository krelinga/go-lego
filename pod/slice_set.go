package pod

import (
	"fmt"
	"iter"
)

type SliceSet[T any] struct {
	slice *Slice[T]
	equal func(a, b T) bool
}

func NewSliceSetFunc[T any](equal func(a, b T) bool, vals ...T) *SliceSet[T] {
	for i := range vals {
		for j := range i {
			if equal(vals[i], vals[j]) {
				panic(fmt.Sprintf("duplicate value at index %d and %d", i, j))
			}
		}
	}
	mySlice := &Slice[T]{}
	CloneValsIntoVec(AsVec(vals), mySlice)
	return &SliceSet[T]{
		slice: mySlice,
		equal: equal,
	}
}

func NewSliceSet[T comparable](vals ...T) *SliceSet[T] {
	return NewSliceSetFunc(func(a, b T) bool { return a == b }, vals...)
}

func (s *SliceSet[T]) Len() int {
	return s.slice.Len()
}

func (s *SliceSet[T]) Has(value T) bool {
	for v := range s.slice.Vals() {
		if s.equal(v, value) {
			return true
		}
	}
	return false
}

func (s *SliceSet[T]) Vals() iter.Seq[T] {
	return s.slice.Vals()
}

func (s *SliceSet[T]) Put(value T) {
	if s.Has(value) {
		return
	}
	s.slice.Push(value)
}

func (s *SliceSet[T]) Clear() {
	s.slice.Clear()
}

func (s *SliceSet[T]) Del(value T) {
	valueIdx := -1
	for i, v := range s.slice.IdxVals() {
		if s.equal(v, value) {
			valueIdx = i
			break
		}
	}
	if valueIdx == -1 {
		return
	}
	before := VecRangeTo(s.slice, valueIdx)
	after := VecRangeFrom(s.slice, valueIdx+1)
	vals := ConcatVals(before, after)
	CloneValsIntoVec(vals, s.slice)
}

func (s *SliceSet[T]) Reserve(n int) {
	s.slice.Reserve(n)
}
