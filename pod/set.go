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

type Set[T any] interface {
	SetView[T]
	Add(value T)
	Clear()
	Delete(value T)
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

func SetOfDictKeys[T, V any](m DictView[T, V]) SetView[T] {
	return setOfDictKeys[T, V]{m: m}
}

type setOfDictKeys[T, V any] struct {
	m DictView[T, V]
}

func (s setOfDictKeys[T, V]) Len() int {
	return s.m.Len()
}

func (s setOfDictKeys[T, V]) Has(value T) bool {
	_, ok := s.m.Get(value)
	return ok
}

func (s setOfDictKeys[T, V]) Vals() iter.Seq[T] {
	return s.m.Keys()
}

type MapSet[T comparable] map[T]struct{}

func CloneSetInto[T any](in SetView[T], out Set[T]) {
	CloneSetIntoFunc(in, out, func(x T) T { return x })
}

func CloneSetIntoFunc[T, U any](in SetView[T], out Set[U], valueFunc func(T) U) {
	out.Clear()
	if canReserve, ok := out.(CanReserve); ok {
		canReserve.Reserve(in.Len())
	}
	for value := range in.Vals() {
		out.Add(valueFunc(value))
	}
}

func (s *MapSet[T]) Len() int {
	return len(*s)
}

func (s *MapSet[T]) Has(value T) bool {
	_, ok := (*s)[value]
	return ok
}

func (s *MapSet[T]) Vals() iter.Seq[T] {
	return maps.Keys(*s)
}

func (s *MapSet[T]) Add(value T) {
	if *s == nil {
		*s = make(map[T]struct{})
	}
	(*s)[value] = struct{}{}
}

func (s *MapSet[T]) Clear() {
	*s = nil
}

func (s *MapSet[T]) Reserve(n int) {
	if *s == nil {
		*s = make(map[T]struct{}, n)
	}
}

func (s *MapSet[T]) Delete(value T) {
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
