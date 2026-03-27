package lego

import "iter"

type MapView[K comparable, V any] interface {
	Len() int
	List() iter.Seq[Pair[K, V]]

	Get(K) (V, bool)
	Has(K) bool
	GetZero(K) V
}

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

func (m Map[K, V]) Has(key K) bool {
	_, ok := m[key]
	return ok
}

func (m Map[K, V]) GetZero(key K) V {
	return m[key]
}

func NewMap[M ~map[K]V, K comparable, V any](m M) Map[K, V] {
	return Map[K, V](m)
}

type MapViewer[K comparable, V1 Viewer[V2], V2 any] Map[K, V1]

func (m MapViewer[K, V1, V2]) Get(key K) (V2, bool) {
	v1, ok := m[key]
	if !ok {
		var zero V2
		return zero, false
	}
	return v1.View(), true
}

func (m MapViewer[K, V1, V2]) Has(key K) bool {
	_, ok := m[key]
	return ok
}

func (m MapViewer[K, V1, V2]) GetZero(key K) V2 {
	v1 := m[key]
	return v1.View()
}

func NewMapViewer[M ~map[K]V1, K comparable, V1 Viewer[V2], V2 any](m M) MapViewer[K, V1, V2] {
	return MapViewer[K, V1, V2](m)
}