package pod

import (
	"iter"
	"maps"
)

type MapView[K, V any] interface {
	Len() int
	Get(key K) (V, bool)
	KeyVals() iter.Seq2[K, V]
	Keys() iter.Seq[K]
	Vals() iter.Seq[V]
}

func AsMap[M ~map[K]V, K comparable, V any](m M) MapView[K, V] {
	return mapView[M, K, V]{m: m}
}

type mapView[M ~map[K]V, K comparable, V any] struct {
	m M
}

func (m mapView[M, K, V]) Len() int {
	return len(m.m)
}

func (m mapView[M, K, V]) Get(key K) (V, bool) {
	value, ok := m.m[key]
	return value, ok
}

func (m mapView[M, K, V]) KeyVals() iter.Seq2[K, V] {
	return maps.All(m.m)
}

func (m mapView[M, K, V]) Keys() iter.Seq[K] {
	return maps.Keys(m.m)
}

func (m mapView[M, K, V]) Vals() iter.Seq[V] {
	return maps.Values(m.m)
}

type Map[K comparable, V any] map[K]V

func CloneMap[K comparable, V any](m MapView[K, V]) *Map[K, V] {
	return CloneMapFunc(m, func(k K) K { return k }, func(v V) V { return v })
}

func CloneMapFunc[K any, KK comparable, V, VV any](m MapView[K, V], keyFunc func(K) KK, valueFunc func(V) VV) *Map[KK, VV] {
	c := &Map[KK, VV]{}
	c.Reserve(m.Len())
	for k, v := range m.KeyVals() {
		c.Set(keyFunc(k), valueFunc(v))
	}
	return c
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

func (m *Map[K, V]) Set(key K, value V) {
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

func (m *Map[K, V]) Delete(key K) {
	delete(*m, key)
}

func WrapMapVals[K, V, W any](m MapView[K, V], wrap func(V) W) MapView[K, W] {
	return wrappedMapVals[K, V, W]{
		m: m,
		wrap: wrap,
	}
}

type wrappedMapVals[K, V, W any] struct {
	m MapView[K, V]
	wrap func(V) W
}

func (w wrappedMapVals[K, V, W]) Len() int {
	return w.m.Len()
}

func (w wrappedMapVals[K, V, W]) Get(key K) (W, bool) {
	value, ok := w.m.Get(key)
	if !ok {
		var zero W
		return zero, false
	}
	return w.wrap(value), true
}

func (w wrappedMapVals[K, V, W]) KeyVals() iter.Seq2[K, W] {
	return func(yield func(K, W) bool) {
		for k, v := range w.m.KeyVals() {
			if !yield(k, w.wrap(v)) {
				return
			}
		}
	}
}

func (w wrappedMapVals[K, V, W]) Keys() iter.Seq[K] {
	return w.m.Keys()
}

func (w wrappedMapVals[K, V, W]) Vals() iter.Seq[W] {
	return func(yield func(W) bool) {
		for v := range w.m.Vals() {
			if !yield(w.wrap(v)) {
				return
			}
		}
	}
}

func WrapMapKeys[K, V, W any](m MapView[K, V], wrap func(K) W, unwrap func(W) K) MapView[W, V] {
	return wrappedMapKeys[K, V, W]{
		m:   m,
		wrap:   wrap,
		unwrap: unwrap,
	}
}

type wrappedMapKeys[K, V, W any] struct {
	m   MapView[K, V]
	wrap   func(K) W
	unwrap func(W) K
}

func (w wrappedMapKeys[K, V, W]) Len() int {
	return w.m.Len()
}

func (w wrappedMapKeys[K, V, W]) Get(key W) (V, bool) {
	return w.m.Get(w.unwrap(key))
}

func (w wrappedMapKeys[K, V, W]) KeyVals() iter.Seq2[W, V] {
	return func(yield func(W, V) bool) {
		for k, v := range w.m.KeyVals() {
			if !yield(w.wrap(k), v) {
				return
			}
		}
	}
}

func (w wrappedMapKeys[K, V, W]) Keys() iter.Seq[W] {
	return func(yield func(W) bool) {
		for k := range w.m.Keys() {
			if !yield(w.wrap(k)) {
				return
			}
		}
	}
}

func (w wrappedMapKeys[K, V, W]) Vals() iter.Seq[V] {
	return w.m.Vals()
}

func MapEqual[K, V comparable](a, b MapView[K, V]) bool {
	return MapEqualFunc(a, b, func(vA, vB V) bool {
		return vA == vB
	})
}

func MapEqualFunc[K, V any](a, b MapView[K, V], eq func(V, V) bool) bool {
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
