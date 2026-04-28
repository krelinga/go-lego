package pod

import (
	"fmt"
	"iter"
	"slices"
)

type SliceSet[T any] struct {
	slice []T
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
	return &SliceSet[T]{
		slice: vals,
		equal: equal,
	}
}

func NewSliceSet[T comparable](vals ...T) *SliceSet[T] {
	return &SliceSet[T]{
		slice: vals,
		equal: func(a, b T) bool { return a == b },
	}
}

func (s *SliceSet[T]) Len() int {
	return len(s.slice)
}

func (s *SliceSet[T]) Has(value T) bool {
	for _, v := range s.slice {
		if s.equal(v, value) {
			return true
		}
	}
	return false
}

func (s *SliceSet[T]) Vals() iter.Seq[T] {
	return slices.Values(s.slice)
}

func (s *SliceSet[T]) Add(value T) {
	if s.Has(value) {
		return
	}
	s.slice = append(s.slice, value)
}

func (s *SliceSet[T]) Clear() {
	s.slice = nil
}

func (s *SliceSet[T]) Delete(value T) {
	for i, v := range s.slice {
		if s.equal(v, value) {
			s.slice = append(s.slice[:i], s.slice[i+1:]...)
			return
		}
	}
}

func (s *SliceSet[T]) Reserve(n int) {
	if cap(s.slice) < n {
		newSlice := make([]T, len(s.slice), n)
		copy(newSlice, s.slice)
		s.slice = newSlice
	}
}
