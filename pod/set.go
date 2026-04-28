package pod

import (
	"iter"
	"maps"
)

type SetView[T any] interface {
	Len() int
	Has(value T) bool
	Vals() iter.Seq[T]
}

func AsSet[S ~map[T]V, T comparable, V any](s S) SetView[T] {
	return setView[S, T, V]{s: s}
}

type setView[S ~map[T]V, T comparable, V any] struct {
	s S
}

func (s setView[S, T, V]) Len() int {
	return len(s.s)
}

func (s setView[S, T, V]) Has(value T) bool {
	_, ok := s.s[value]
	return ok
}

func (s setView[S, T, V]) Vals() iter.Seq[T] {
	return maps.Keys(s.s)
}

func SetOfMapKeys[T comparable, V any](m DictView[T, V]) SetView[T] {
	return setOfMapKeys[T, V]{m: m}
}

type setOfMapKeys[T comparable, V any] struct {
	m DictView[T, V]
}

func (s setOfMapKeys[T, V]) Len() int {
	return s.m.Len()
}

func (s setOfMapKeys[T, V]) Has(value T) bool {
	_, ok := s.m.Get(value)
	return ok
}

func (s setOfMapKeys[T, V]) Vals() iter.Seq[T] {
	return s.m.Keys()
}

type Set[T comparable] map[T]struct{}

func CloneSet[T comparable](set SetView[T]) *Set[T] {
	return CloneSetFunc(set, func(x T) T { return x })
}

func CloneSetFunc[T any, U comparable](set SetView[T], valueFunc func(T) U) *Set[U] {
	s := &Set[U]{}
	s.Reserve(set.Len())
	for value := range set.Vals() {
		s.Add(valueFunc(value))
	}
	return s
}

func (s *Set[T]) Len() int {
	return len(*s)
}

func (s *Set[T]) Has(value T) bool {
	_, ok := (*s)[value]
	return ok
}

func (s *Set[T]) Vals() iter.Seq[T] {
	return maps.Keys(*s)
}

func (s *Set[T]) Add(value T) {
	if *s == nil {
		*s = make(map[T]struct{})
	}
	(*s)[value] = struct{}{}
}

func (s *Set[T]) Clear() {
	*s = nil
}

func (s *Set[T]) Reserve(n int) {
	if *s == nil {
		*s = make(map[T]struct{}, n)
	}
}

func (s *Set[T]) Delete(value T) {
	delete(*s, value)
}

func WrapSetVals[T, V any](set SetView[T], wrap func(T) V, unwrap func(V) T) SetView[V] {
	return wrappedSetVals[T, V]{
		set:    set,
		wrap:   wrap,
		unwrap: unwrap,
	}
}

type wrappedSetVals[T, V any] struct {
	set    SetView[T]
	wrap   func(T) V
	unwrap func(V) T
}

func (w wrappedSetVals[T, V]) Len() int {
	return w.set.Len()
}

func (w wrappedSetVals[T, V]) Has(value V) bool {
	return w.set.Has(w.unwrap(value))
}

func (w wrappedSetVals[T, V]) Vals() iter.Seq[V] {
	return func(yield func(V) bool) {
		for x := range w.set.Vals() {
			if !yield(w.wrap(x)) {
				return
			}
		}
	}
}

func SetEqualFunc[T any](a, b SetView[T], eq func(T, T) bool) bool {
	if a.Len() != b.Len() {
		return false
	}
	for x := range a.Vals() {
		if !b.Has(x) {
			return false
		}
	}
	return true
}

func SetEqual[T comparable](a, b SetView[T]) bool {
	return SetEqualFunc(a, b, func(x, y T) bool {
		return x == y
	})
}
