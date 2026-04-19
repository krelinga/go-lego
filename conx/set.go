package conx

import (
	"iter"
	"maps"
)

type SetView[T any] interface {
	Len() int
	Has(value T) bool
	Values() iter.Seq[T]
}

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](values ...T) (s *Set[T]) {
	s = &Set[T]{}
	if len(values) == 0 {
		return
	}
	s.Reserve(len(values))
	for _, value := range values {
		s.Add(value)
	}
	return
}

func CloneSet[T comparable](set SetView[T]) *Set[T] {
	return CloneSetFunc(set, func(x T) T { return x })
}

func CloneSetFunc[T any, U comparable](set SetView[T], valueFunc func(T) U) *Set[U] {
	s := &Set[U]{}
	s.Reserve(set.Len())
	for value := range set.Values() {
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

func (s *Set[T]) Values() iter.Seq[T] {
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

func WrapSetValues[T, V any](set SetView[T], wrap func(T) V, unwrap func(V) T) SetView[V] {
	return wrappedSetValues[T, V]{
		set:    set,
		wrap:   wrap,
		unwrap: unwrap,
	}
}

type wrappedSetValues[T, V any] struct {
	set    SetView[T]
	wrap   func(T) V
	unwrap func(V) T
}

func (w wrappedSetValues[T, V]) Len() int {
	return w.set.Len()
}

func (w wrappedSetValues[T, V]) Has(value V) bool {
	return w.set.Has(w.unwrap(value))
}

func (w wrappedSetValues[T, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for x := range w.set.Values() {
			if !yield(w.wrap(x)) {
				return
			}
		}
	}
}

func NewSetViewFromMapKeys[M ~map[K]V, K comparable, V any](m M) SetView[K] {
	return setViewFromMapKeys[M, K, V]{m: m}
}

type setViewFromMapKeys[M ~map[K]V, K comparable, V any] struct {
	m M
}

func (w setViewFromMapKeys[M, K, V]) Len() int {
	return len(w.m)
}

func (w setViewFromMapKeys[M, K, V]) Has(value K) bool {
	_, ok := w.m[value]
	return ok
}

func (w setViewFromMapKeys[M, K, V]) Values() iter.Seq[K] {
	return maps.Keys(w.m)
}

func SetEqualFunc[T any](a, b SetView[T], eq func(T, T) bool) bool {
	if a.Len() != b.Len() {
		return false
	}
	for x := range a.Values() {
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