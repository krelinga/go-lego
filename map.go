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

// A GoMap is map that wraps Go's built-in map type.
// It implements the [FixedMap] interface.
// It does not implement the [Adder] interface since Go's built-in maps may be nil and thus not safe to add to without reassignment, which this type does not support.
type GoMap[K comparable, V any] map[K]V

func (m GoMap[K, V]) Len() int {
	return len(m)
}

func (m GoMap[K, V]) List() iter.Seq[Pair[K, V]] {
	return func(yield func(Pair[K, V]) bool) {
		for k, v := range m {
			if !yield(NewPair(k, v)) {
				return
			}
		}
	}
}

func (m GoMap[K, V]) Get(key K) (V, bool) {
	v, ok := m[key]
	return v, ok
}

// A Map is a mutable map.
// It implements the [FixedMap] interface and the [Adder] interface.
type Map[K comparable, V any] struct {
	m GoMap[K, V]
}

func (m Map[K, V]) Len() int {
	return len(m.m)
}

func (m Map[K, V]) List() iter.Seq[Pair[K, V]] {
	return m.m.List()
}

func (m Map[K, V]) Get(key K) (V, bool) {
	return m.m.Get(key)
}

func (m Map[K, V]) Add(pair Pair[K, V]) {
	if m.m == nil {
		m.m = GoMap[K, V]{}
	}
	m.m[pair.GetKey()] = pair.GetValue()
}

// Reserve reserves space for n elements in the map. This is a best-effort operation and will do nothing if the map already contains some values, since Go's built-in maps do not support reserving space after initialization.
func (m Map[K, V]) Reserve(n int) {
	if m.m == nil {
		m.m = make(GoMap[K, V], n)
	}
}

func NewMap[M ~map[K]V, K comparable, V any](m M) Map[K, V] {
	return Map[K, V]{m: GoMap[K, V](m)}
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

// ViewMap creates a view of a map that allows viewing the values of the map as a different type, without modifying the original map. It panics if the map has pointer keys, since pointer keys are mutable and would violate the immutability guarantee of the view.
func ViewMap[M FixedMap[K, V1], K comparable, V1 Viewer[V2], V2 any](m M) FixedMap[K, V2] {
	t := reflect.TypeFor[K]()
	if t.Kind() == reflect.Pointer {
		panic("cannot create a view of a map with pointer keys")
	}
	return mapView[M, K, V1, V2]{m: m}
}

type mapView[M FixedMap[K, V1], K comparable, V1 Viewer[V2], V2 any] struct {
	m M
}

func (v mapView[M, K, V1, V2]) Len() int {
	return v.m.Len()
}

func (v mapView[M, K, V1, V2]) List() iter.Seq[Pair[K, V2]] {
	return func(yield func(Pair[K, V2]) bool) {
		for pair := range v.m.List() {
			if !yield(NewPair(pair.GetKey(), pair.GetValue().View())) {
				return
			}
		}
	}
}

func (v mapView[M, K, V1, V2]) Get(key K) (V2, bool) {
	v1, ok := v.m.Get(key)
	if !ok {
		var zero V2
		return zero, false
	}
	return v1.View(), true
}

func DeepCopyMap[M FixedMap[K, V], K comparable, V DeepCopier[V]](m M) Map[K, V] {
	if t := reflect.TypeFor[K](); t.Kind() == reflect.Pointer || t.Kind() == reflect.Interface {
		panic("cannot deep copy a map with pointer or interface keys")
	}
	var out Map[K, V]
	out.Reserve(m.Len())
	for pair := range m.List() {
		out.Add(NewPair(pair.GetKey(), pair.GetValue().DeepCopy()))
	}
	return out
}
