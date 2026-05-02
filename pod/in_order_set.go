package pod

import (
	"iter"
)

type InOrderSet[T any] struct {
	slice *Slice[T]
	equal func(a, b T) bool
}

func NewInOrderSetFunc[T any](equal func(a, b T) bool, vals ...T) *InOrderSet[T] {
	s := &InOrderSet[T]{
		slice: &Slice[T]{},
		equal: equal,
	}
	for _, v := range vals {
		s.Put(v)
	}
	return s
}

func NewInOrderSet[T comparable](vals ...T) *InOrderSet[T] {
	return NewInOrderSetFunc(func(a, b T) bool { return a == b }, vals...)
}

func (s *InOrderSet[T]) Len() int {
	return s.slice.Len()
}

func (s *InOrderSet[T]) Has(value T) bool {
	for v := range s.slice.Vals() {
		if s.equal(v, value) {
			return true
		}
	}
	return false
}

func (s *InOrderSet[T]) Vals() iter.Seq[T] {
	return s.slice.Vals()
}

func (s *InOrderSet[T]) Put(value T) {
	if s.Has(value) {
		return
	}
	s.slice.Push(value)
}

func (s *InOrderSet[T]) PutVals(vals Vals[T]) {
	for value := range vals.Vals() {
		s.Put(value)
	}
}

func (s *InOrderSet[T]) Clear() {
	s.slice.Clear()
}

func (s *InOrderSet[T]) Del(value T) {
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

func (s *InOrderSet[T]) Reserve(n int) {
	s.slice.Reserve(n)
}
