package pod

import (
	"iter"
	"maps"

	"github.com/krelinga/go-libs/zero"
)

// DictView is a read-only view of a dictionary. It provides methods to access the keys and values, but does not allow mutation.
type DictView[K, V any] interface {
	Len() int
	Get(key K) (V, bool)
	KeyVals() iter.Seq2[K, V]
	Keys() iter.Seq[K]
	Vals() iter.Seq[V]
}

// Dict is a mutable dictionary that implements the DictView interface. It allows adding, removing, and clearing key-value pairs.
type Dict[K, V any] interface {
	DictView[K, V]
	Put(key K, value V)
	Clear()
	Del(key K)
}

// AsDict creates a DictView from a map. It provides a read-only view of the map, and any changes to the map will be reflected in the DictView.
func AsDict[M ~map[K]V, K comparable, V any](m M) DictView[K, V] {
	return dictView[M, K, V]{m: m}
}

type dictView[M ~map[K]V, K comparable, V any] struct {
	m M
}

func (m dictView[M, K, V]) Len() int {
	return len(m.m)
}

func (m dictView[M, K, V]) Get(key K) (V, bool) {
	value, ok := m.m[key]
	return value, ok
}

func (m dictView[M, K, V]) KeyVals() iter.Seq2[K, V] {
	return maps.All(m.m)
}

func (m dictView[M, K, V]) Keys() iter.Seq[K] {
	return maps.Keys(m.m)
}

func (m dictView[M, K, V]) Vals() iter.Seq[V] {
	return maps.Values(m.m)
}

// Map is a mutable map that implements the Dict interface. It provides methods to add, remove, and clear key-value pairs, as well as access the keys and values.
// It is a wrapper around the built-in map type, so it is possible to create literals of Map as follows:
//
//	m := &Map[string, int]{"a": 1, "b": 2}
type Map[K comparable, V any] map[K]V

// NewMap creates a new empty Map with the specified key and value types.
func NewMap[K comparable, V any]() *Map[K, V] {
	m := make(map[K]V)
	return (*Map[K, V])(&m)
}

// NewMapHint creates a new Map with the specified key and value types, and reserves space for the specified number of elements.
func NewMapHint[K comparable, V any](hint int) *Map[K, V] {
	m := make(map[K]V, hint)
	return (*Map[K, V])(&m)
}

// NewMapOf creates a new Map from a Bag2 of key-value pairs. It collects the pairs into a map and returns a pointer to the Map.
func NewMapOf[K comparable, V any](b Bag2[K, V]) *Map[K, V] {
	m := maps.Collect(b.Elems())
	return (*Map[K, V])(&m)
}

// Len returns the number of key-value pairs in the Map.
func (m *Map[K, V]) Len() int {
	return len(*m)
}

// Get retrieves the value associated with the given key. It returns the value and a boolean indicating whether the key was found in the Map.
func (m *Map[K, V]) Get(key K) (V, bool) {
	value, ok := (*m)[key]
	return value, ok
}

// KeyVals returns a sequence of key-value pairs in the Map.
func (m *Map[K, V]) KeyVals() iter.Seq2[K, V] {
	return maps.All(*m)
}

// Keys returns a sequence of keys in the Map.
func (m *Map[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(*m)
}

// Vals returns a sequence of values in the Map.
func (m *Map[K, V]) Vals() iter.Seq[V] {
	return maps.Values(*m)
}

// Put adds a key-value pair to the Map. If the key already exists, it replaces the value.
func (m *Map[K, V]) Put(key K, value V) {
	if *m == nil {
		*m = make(map[K]V)
	}
	(*m)[key] = value
}

// Clear removes all key-value pairs from the Map, leaving it empty.
func (m *Map[K, V]) Clear() {
	*m = nil
}

// Reserve allows callers to pre-allocate space for a certain number of key-value pairs in an empty map. If the Map is already initialized, it does nothing.
func (m *Map[K, V]) Reserve(n int) {
	if *m == nil {
		*m = make(map[K]V, n)
	}
}

// Del removes the key-value pair associated with the given key from the Map. If the key does not exist, it does nothing.
func (m *Map[K, V]) Del(key K) {
	delete(*m, key)
}

// WrapDictVals creates a new DictView that wraps the values of the given DictView with the provided wrap function. The keys remain unchanged.
func WrapDictVals[K, V, W any](d DictView[K, V], wrap func(V) W) DictView[K, W] {
	return wrappedDictVals[K, V, W]{
		d:    d,
		wrap: wrap,
	}
}

type wrappedDictVals[K, V, W any] struct {
	d    DictView[K, V]
	wrap func(V) W
}

func (w wrappedDictVals[K, V, W]) Len() int {
	return w.d.Len()
}

func (w wrappedDictVals[K, V, W]) Get(key K) (W, bool) {
	value, ok := w.d.Get(key)
	if !ok {
		return zero.For[W](), false
	}
	return w.wrap(value), true
}

func (w wrappedDictVals[K, V, W]) KeyVals() iter.Seq2[K, W] {
	return func(yield func(K, W) bool) {
		for k, v := range w.d.KeyVals() {
			if !yield(k, w.wrap(v)) {
				return
			}
		}
	}
}

func (w wrappedDictVals[K, V, W]) Keys() iter.Seq[K] {
	return w.d.Keys()
}

func (w wrappedDictVals[K, V, W]) Vals() iter.Seq[W] {
	return func(yield func(W) bool) {
		for v := range w.d.Vals() {
			if !yield(w.wrap(v)) {
				return
			}
		}
	}
}

// WrapDictKeys creates a new DictView that wraps the keys of the given DictView with the provided wrap function.
// The values remain unchanged. The unwrap function is used to convert wrapped keys back to their original form
// when accessing the underlying DictView in cases like Get.
//
// For example, if you have a DictView with string keys and you want to wrap the keys as integers (e.g., by parsing the strings), you could use WrapDictKeys with a wrap function that converts strings to integers and an unwrap function that converts integers back to strings.
//
// The wrap function must be one-to-one (i.e., it should not map two different keys to the same wrapped key) to ensure that the resulting DictView behaves correctly. The unwrap function must be the inverse of the wrap function for the keys to ensure that Get and other methods work as expected.
func WrapDictKeys[K, V, W any](d DictView[K, V], wrap func(K) W, unwrap func(W) K) DictView[W, V] {
	return wrappedDictKeys[K, V, W]{
		d:      d,
		wrap:   wrap,
		unwrap: unwrap,
	}
}

type wrappedDictKeys[K, V, W any] struct {
	d      DictView[K, V]
	wrap   func(K) W
	unwrap func(W) K
}

func (w wrappedDictKeys[K, V, W]) Len() int {
	return w.d.Len()
}

func (w wrappedDictKeys[K, V, W]) Get(key W) (V, bool) {
	return w.d.Get(w.unwrap(key))
}

func (w wrappedDictKeys[K, V, W]) KeyVals() iter.Seq2[W, V] {
	return func(yield func(W, V) bool) {
		for k, v := range w.d.KeyVals() {
			if !yield(w.wrap(k), v) {
				return
			}
		}
	}
}

func (w wrappedDictKeys[K, V, W]) Keys() iter.Seq[W] {
	return func(yield func(W) bool) {
		for k := range w.d.Keys() {
			if !yield(w.wrap(k)) {
				return
			}
		}
	}
}

func (w wrappedDictKeys[K, V, W]) Vals() iter.Seq[V] {
	return w.d.Vals()
}

// DictEqual checks if two DictViews are equal by comparing their key-value pairs. It returns true if both DictViews have the same keys and corresponding values, and false otherwise. The values are compared using the equality operator (==), so the value type must be comparable.
func DictEqual[K, V comparable](a, b DictView[K, V]) bool {
	return DictEqualFunc(a, b, func(vA, vB V) bool {
		return vA == vB
	})
}

// DictEqualFunc checks if two DictViews are equal by comparing their key-value pairs using a custom equality function for the values. It returns true if both DictViews have the same keys and corresponding values that are considered equal by the provided function, and false otherwise.
func DictEqualFunc[K, V any](a, b DictView[K, V], eq func(V, V) bool) bool {
	if a.Len() != b.Len() {
		return false
	}
	for k, vA := range a.KeyVals() {
		vB, ok := b.Get(k)
		if !ok || !eq(vA, vB) {
			return false
		}
	}
	return true
}
