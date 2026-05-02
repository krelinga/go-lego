package pod

import (
	"iter"
	"maps"
)

type DictView[K, V any] interface {
	Len() int
	Get(key K) (V, bool)
	KeyVals() iter.Seq2[K, V]
	Keys() iter.Seq[K]
	Vals() iter.Seq[V]
}

type Dict[K, V any] interface {
	DictView[K, V]
	Put(key K, value V)
	Clear()
	Del(key K)
}

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

type Map[K comparable, V any] map[K]V

func NewMap[K comparable, V any]() *Map[K, V] {
	m := make(map[K]V)
	return (*Map[K, V])(&m)
}

func NewMapHint[K comparable, V any](hint int) *Map[K, V] {
	m := make(map[K]V, hint)
	return (*Map[K, V])(&m)
}

func (m *Map[K, V]) Len() int {
	return len(*m)
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	value, ok := (*m)[key]
	return value, ok
}

func (m *Map[K, V]) KeyVals() iter.Seq2[K, V] {
	return maps.All(*m)
}

func (m *Map[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(*m)
}

func (m *Map[K, V]) Vals() iter.Seq[V] {
	return maps.Values(*m)
}

func (m *Map[K, V]) Put(key K, value V) {
	if *m == nil {
		*m = make(map[K]V)
	}
	(*m)[key] = value
}

func (m *Map[K, V]) Clear() {
	*m = nil
}

func (m *Map[K, V]) Reserve(n int) {
	if *m == nil {
		*m = make(map[K]V, n)
	}
}

func (m *Map[K, V]) Del(key K) {
	delete(*m, key)
}

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
		var zero W
		return zero, false
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

func DictEqual[K, V comparable](a, b DictView[K, V]) bool {
	return DictEqualFunc(a, b, func(vA, vB V) bool {
		return vA == vB
	})
}

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
