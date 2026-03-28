package lego

import (
	"iter"
	"reflect"
)

// A FixedMap is a map that does not allow adding or removing elements, but which may still allow modifying the elements in the map (for example, if the values are pointers).
type FixedMap[K comparable, V any] interface {
	Len() int
	List() iter.Seq[Pair[K, V]]

	Get(K) (V, bool)
}

// A Map is a mutable map that wraps Go's built-in map type.
// It implements the [FixedMap] interface.
type Map[K comparable, V any] map[K]V

func (m Map[K, V]) Len() int {
	return len(m)
}

func (m Map[K, V]) List() iter.Seq[Pair[K, V]] {
	return func(yield func(Pair[K, V]) bool) {
		for k, v := range m {
			if !yield(NewPair(k, v)) {
				return
			}
		}
	}
}

func (m Map[K, V]) Get(key K) (V, bool) {
	v, ok := m[key]
	return v, ok
}

func NewMap[M ~map[K]V, K comparable, V any](m M) Map[K, V] {
	return Map[K, V](m)
}

// A ViewerMap is a special case of [Map] that stores values that implement the [Viewer] interface, and provides a method to get a view of the map.
type ViewerMap[K comparable, V1 Viewer[V2], V2 any] Map[K, V1]

func (m ViewerMap[K, V1, V2]) Len() int {
	return Map[K, V1](m).Len()
}

func (m ViewerMap[K, V1, V2]) List() iter.Seq[Pair[K, V1]] {
	return Map[K, V1](m).List()
}

func (m ViewerMap[K, V1, V2]) Get(key K) (V1, bool) {
	return Map[K, V1](m).Get(key)
}

// View returns a view of the map. It panics if the map has pointer keys, since pointer keys are mutable and would violate the immutability guarantee of the view.
func (m ViewerMap[K, V1, V2]) View() FixedMap[K, V2] {
	t := reflect.TypeFor[K]()
	if t.Kind() == reflect.Pointer {
		panic("cannot create a view of a map with pointer keys")
	}
	return viewerMapView[K, V1, V2]{m: m}
}

func NewViewerMap[M ~map[K]V1, K comparable, V1 Viewer[V2], V2 any](m M) ViewerMap[K, V1, V2] {
	return ViewerMap[K, V1, V2](m)
}

type viewerMapView[K comparable, V1 Viewer[V2], V2 any] struct {
	m ViewerMap[K, V1, V2]
}

func (v viewerMapView[K, V1, V2]) Len() int {
	return v.m.Len()
}

func (v viewerMapView[K, V1, V2]) List() iter.Seq[Pair[K, V2]] {
	return func(yield func(Pair[K, V2]) bool) {
		for pair := range v.m.List() {
			if !yield(NewPair(pair.GetKey(), pair.GetValue().View())) {
				return
			}
		}
	}
}

func (v viewerMapView[K, V1, V2]) Get(key K) (V2, bool) {
	v1, ok := v.m.Get(key)
	if !ok {
		var zero V2
		return zero, false
	}
	return v1.View(), true
}

type Getter[K, V any] interface {
	Get(K) (V, bool)
}

func Get[G Getter[K, V], K, V any](g G, key K) (V, bool) {
	return g.Get(key)
}

func GetOr[G Getter[K, V], K, V any](g G, key K, or V) V {
	v, ok := g.Get(key)
	if !ok {
		return or
	}
	return v
}

func GetZero[G Getter[K, V], K, V any](g G, key K) V {
	v, _ := g.Get(key)
	return v
}

func Has[G Getter[K, V], K, V any](g G, key K) bool {
	_, ok := g.Get(key)
	return ok
}

func GetPanic[G Getter[K, V], K, V any](g G, key K) V {
	v, ok := g.Get(key)
	if !ok {
		panic("key not found")
	}
	return v
}
