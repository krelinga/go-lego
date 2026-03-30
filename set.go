package lego

import (
	"iter"
	"reflect"
)

// A FixedSet is a set that does not allow adding or removing elements, but which may still allow modifying the elements in the set (for example, if the elements are pointers).
type FixedSet[V comparable] interface {
	Len() int
	List() iter.Seq[V]

	Has(V) bool
}

// A GoSetValue is a placeholder value used in the implementation of GoSet, since Go's built-in maps do not allow sets directly.
type GoSetValue struct{}

// A GoSet is a set that wraps Go's built-in map type with a placeholder value.
// It implements the [FixedSet] interface.
// It does not implement the [Adder] interface since Go's built-in maps may be nil and thus not safe to add to without reassignment, which this type does not support.
type GoSet[V comparable] map[V]GoSetValue

func (s GoSet[V]) Len() int {
	return len(s)
}

func (s GoSet[V]) List() iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

func (s GoSet[V]) Has(v V) bool {
	_, ok := s[v]
	return ok
}

// A Set is a mutable set.
// It implements the [FixedSet] interface and the [Adder] interface.
type Set[V comparable] struct {
	GoSet[V]
}

func (s *Set[V]) Add(v V) {
	if s.GoSet == nil {
		s.GoSet = GoSet[V]{}
	}
	s.GoSet[v] = GoSetValue{}
}

// Reserve reserves space for n elements in the set. This is a best-effort operation and will do nothing if the set already contains some values, since Go's built-in maps do not support reserving space after initialization.
func (s *Set[V]) Reserve(n int) {
	if s.GoSet == nil {
		s.GoSet = make(GoSet[V], n)
	}
}

func NewSet[V comparable](values ...V) *Set[V] {
	m := make(GoSet[V], len(values))
	for _, v := range values {
		m[v] = GoSetValue{}
	}
	return &Set[V]{GoSet: m}
}

func DeepCopySet[S FixedSet[V], V comparable](s S) *Set[V] {
	if t := reflect.TypeFor[V](); t.Kind() == reflect.Pointer || t.Kind() == reflect.Interface {
		panic("cannot deep copy a set with pointer or interface elements")
	}
	var out Set[V]
	out.Reserve(s.Len())
	for v := range s.List() {
		out.Add(v)
	}
	return &out
}
func ShallowCopySet[S FixedSet[V], V comparable](s S) *Set[V] {
	var out Set[V]
	out.Reserve(s.Len())
	for v := range s.List() {
		out.Add(v)
	}
	return &out
}