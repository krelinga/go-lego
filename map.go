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

type FluidMap[K comparable, V any] interface {
	FixedMap[K, V]

	// Reserve reserves space for n elements in the map. This is a best-effort operation and will do nothing if the map already contains some values, since Go's built-in maps do not support reserving space after initialization.
	Reserve(int)

	// Add adds a key-value pair to the map. If the map already contains a value for the given key, it will be replaced with the new value.
	Add(Pair[K, V])

	// Set stores the given value in the map under the given key, replacing any existing value for that key.
	Set(K, V)

	// LegoMap returns the underlying *Map[K, V] that implements the FluidMap interface. This is used internally by functions like Sort and SortFunc to access the underlying map for sorting.
	// It is helpful because several other types are expected to embed a *Map[K, V] to implement the FluidMap interface, and this method provides a consistent way to access the underlying map regardless of the embedding type.
	LegoMap() *Map[K, V]
}

// A Map is a mutable map.
// It implements the [FixedMap] interface and the [Adder] interface.
type Map[K comparable, V any] map[K]V

func (m *Map[K, V]) Len() int {
	return len(*m)
}

func (m *Map[K, V]) List() iter.Seq[Pair[K, V]] {
	return func(yield func(Pair[K, V]) bool) {
		for k, v := range *m {
			if !yield(NewPair(k, v)) {
				return
			}
		}
	}
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	v, ok := (*m)[key]
	return v, ok
}


func (m *Map[K, V]) Add(pair Pair[K, V]) {
	if *m == nil {
		*m = make(Map[K, V])
	}
	(*m)[pair.GetKey()] = pair.GetValue()
}

// Set stores the given value in the map under the given key, replacing any existing value for that key.
func (m *Map[K, V]) Set(key K, value V) {
	if *m == nil {
		*m = make(Map[K, V])
	}
	(*m)[key] = value
}

// Reserve reserves space for n elements in the map. This is a best-effort operation and will do nothing if the map already contains some values, since Go's built-in maps do not support reserving space after initialization.
func (m *Map[K, V]) Reserve(n int) {
	if *m == nil {
		*m = make(Map[K, V], n)
	}
}

func (m *Map[K, V]) LegoMap() *Map[K, V] {
	return m
}

// NewMap creates a new map with the given entries.
func NewMap[K comparable, V any]() *Map[K, V] {
	m := make(Map[K, V])
	return &m
}

// NewMapReserve creates a new map with the given entries, and reserves space for n elements in the map.
func NewMapHint[K comparable, V any](n int) *Map[K, V] {
	m := make(Map[K, V], n)
	return &m
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

func MustGet[G Getter[K, V], K, V any](g G, key K) V {
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

func DeepCopyMap[M FixedMap[K, V], K comparable, V DeepCopier[V]](m M) *Map[K, V] {
	if t := reflect.TypeFor[K](); t.Kind() == reflect.Pointer || t.Kind() == reflect.Interface {
		panic("cannot deep copy a map with pointer or interface keys")
	}
	var out Map[K, V]
	out.Reserve(m.Len())
	for pair := range m.List() {
		out.Add(NewPair(pair.GetKey(), pair.GetValue().DeepCopy()))
	}
	return &out
}

func ShallowCopyMap[M FixedMap[K, V], K comparable, V any](m M) *Map[K, V] {
	var out Map[K, V]
	out.Reserve(m.Len())
	for pair := range m.List() {
		out.Add(pair)
	}
	return &out
}
