package lego

import "iter"

// A FixedSet is a set that does not allow adding or removing elements, but which may still allow modifying the elements in the set (for example, if the elements are pointers).
type FixedSet[V comparable] interface {
	Len() int
	List() iter.Seq[V]

	Has(V) bool
}

// A Set is a mutable set that wraps Go's built-in map type.
// It implements the [FixedSet] interface.
type Set[V comparable] struct {
	m map[V]struct{}
}

func (s Set[V]) Len() int {
	return len(s.m)
}

func (s Set[V]) List() iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range s.m {
			if !yield(v) {
				return
			}
		}
	}
}

func (s Set[V]) Has(v V) bool {
	_, ok := s.m[v]
	return ok
}

func NewSet[V comparable](values ...V) Set[V] {
	m := make(map[V]struct{}, len(values))
	for _, v := range values {
		m[v] = struct{}{}
	}
	return Set[V]{m: m}
}
