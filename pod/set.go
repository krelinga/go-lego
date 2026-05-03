package pod

import (
	"iter"
	"maps"
)

// SetView is a read-only view of a set. It provides methods to access the values, but does not allow mutation.
type SetView[T any] interface {
	Len() int
	Has(value T) bool
	Vals() iter.Seq[T]
}

// Set is a mutable set that implements the SetView interface. It allows adding, removing, and clearing values.
type Set[T any] interface {
	SetView[T]
	Put(value T)
	Clear()
	Del(value T)
}

// AsSet creates a SetView from the keys of a map. It provides a read-only view of the keys of the map as a set, and any changes to the map will be reflected in the SetView.
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

// SetOfDictKeys creates a SetView of the keys of a DictView. It provides a read-only view of the keys of the DictView as a set, and any changes to the DictView will be reflected in the SetView.
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

// MapSet is a mutable set that uses a map for storage. It is a wrapper around a map with empty struct values, so it is possible to create literals of MapSet as follows:
//
//	s := &MapSet[string]{"a": {}, "b": {}}
type MapSet[T comparable] map[T]struct{}

// NewMapSet creates a new empty MapSet.
func NewMapSet[T comparable]() *MapSet[T] {
	s := make(map[T]struct{})
	return (*MapSet[T])(&s)
}

// NewMapSetHint creates a new empty MapSet with a hint for the initial capacity. This can help improve performance if you know in advance how many elements will be added to the set.
func NewMapSetHint[T comparable](hint int) *MapSet[T] {
	s := make(map[T]struct{}, hint)
	return (*MapSet[T])(&s)
}

// NewMapSetOf creates a new MapSet from a Bag of values. It collects the values into a set and returns a pointer to the MapSet.
// The returned MapSet will contain all the unique values from the Bag, and any duplicates in the Bag will be ignored.
// The resulting MapSet will be separate from the Bag, so changes to the Bag will not affect the MapSet, and vice versa.
func NewMapSetOf[T comparable](b Bag[T]) *MapSet[T] {
	s := NewMapSetHint[T](b.Len())
	for v := range b.Elems() {
		s.Put(v)
	}
	return s
}

// Len returns the number of elements in the MapSet.
func (s *MapSet[T]) Len() int {
	return len(*s)
}

// Has checks if the given value is present in the MapSet. It returns true if the value is found, and false otherwise.
func (s *MapSet[T]) Has(value T) bool {
	_, ok := (*s)[value]
	return ok
}

// Vals returns a sequence of values in the MapSet. The order of the values is not guaranteed, as it depends on the internal ordering of the map keys.
func (s *MapSet[T]) Vals() iter.Seq[T] {
	return maps.Keys(*s)
}

// Put adds a value to the MapSet. If the value is already present, it does nothing.
func (s *MapSet[T]) Put(value T) {
	if *s == nil {
		*s = make(map[T]struct{})
	}
	(*s)[value] = struct{}{}
}

// Clear removes all values from the MapSet, leaving it empty.
func (s *MapSet[T]) Clear() {
	*s = nil
}

// Reserve allows callers to pre-allocate space for a certain number of values in an empty MapSet. If the MapSet is already initialized, it does nothing.
func (s *MapSet[T]) Reserve(n int) {
	if *s == nil {
		*s = make(map[T]struct{}, n)
	}
}

// Del removes a value from the MapSet. If the value is not present, it does nothing.
func (s *MapSet[T]) Del(value T) {
	delete(*s, value)
}

// WrapSetVals creates a new SetView that wraps the values of the given SetView with the provided wrap and unwrap functions. The wrap function is used to convert values from the original type to the wrapped type when yielding values, and the unwrap function is used to convert values from the wrapped type back to the original type when checking for membership with Has.
//
// For example, if you have a SetView of integers and you want to wrap the values as strings (e.g., by converting the integers to their string representations), you could use WrapSetVals with a wrap function that converts integers to strings and an unwrap function that converts strings back to integers.
//
// The wrap function must be one-to-one (i.e., it should not map two different values to the same wrapped value) to ensure that the resulting SetView behaves correctly. The unwrap function must be the inverse of the wrap function for the values to ensure that Has and other methods work as expected.
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

// SetEqual checks if two SetViews are equal by comparing their elements. It returns true if both SetViews contain the same elements, and false otherwise. The elements are compared using the Set's Has method.
func SetEqual[T any](a, b SetView[T]) bool {
	if a.Len() != b.Len() {
		return false
	}
	for v := range a.Vals() {
		if !b.Has(v) {
			return false
		}
	}
	return true
}
