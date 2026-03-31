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

type FluidSet[V comparable] interface {
	FixedSet[V]

	// Reserve reserves space for n elements in the set. This is a best-effort operation and will do nothing if the set already contains some values, since Go's built-in maps do not support reserving space after initialization.
	Reserve(int)

	// Add adds an element to the set. If the set already contains the given element, it will be replaced with the new value.
	Add(V)

	// LegoSet returns the underlying *Set[V] that implements the FluidSet interface. This is used internally by functions like Sort and SortFunc to access the underlying set for sorting.
	// It is helpful because several other types are expected to embed a *Set[V] to implement the FluidSet interface, and this method provides a consistent way to access the underlying set regardless of the embedding type.
	LegoSet() *Set[V]
}

// A SetValue is a placeholder value used in the implementation of Set, since Go's built-in maps do not allow sets directly.
type SetValue struct{}

// A Set is a mutable set.
// It implements the [FixedSet] interface and the [Adder] interface.
type Set[V comparable] map[V]SetValue

func (s *Set[V]) Len() int {
	return len(*s)
}

func (s *Set[V]) List() iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range *s {
			if !yield(v) {
				return
			}
		}
	}
}

func (s *Set[V]) Has(v V) bool {
	_, ok := (*s)[v]
	return ok
}

func (s *Set[V]) Add(v V) {
	if *s == nil {
		*s = make(Set[V])
	}
	(*s)[v] = SetValue{}
}

// Reserve reserves space for n elements in the set. This is a best-effort operation and will do nothing if the set already contains some values, since Go's built-in maps do not support reserving space after initialization.
func (s *Set[V]) Reserve(n int) {
	if *s == nil {
		*s = make(Set[V], n)
	}
}

func (s *Set[V]) LegoSet() *Set[V] {
	return s
}

func NewSet[V comparable]() *Set[V] {
	var s Set[V]
	return &s
}

func NewSetHint[V comparable](n int) *Set[V] {
	var s Set[V]
	s.Reserve(n)
	return &s
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
