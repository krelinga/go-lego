package lego

import "iter"

type MapView[K comparable, V any] interface {
	Len() int
	List() iter.Seq[Pair[K, V]]

	Get(K) (V, bool)
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

func NewMap[M ~map[K]V, K comparable, V any](m M) Map[K, V] {
	return Map[K, V](m)
}

type MapViewer[K comparable, V1 Viewer[V2], V2 any] Map[K, V1]

func (m MapViewer[K, V1, V2]) Len() int {
	return Map[K, V1](m).Len()
}

func (m MapViewer[K, V1, V2]) List() iter.Seq[Pair[K, V2]] {
	return func(yield func(Pair[K, V2]) bool) {
		for k, v1 := range m {
			v2 := v1.View()
			if !yield(NewPair(k, v2)) {
				return
			}
		}
	}
}

func (m MapViewer[K, V1, V2]) Get(key K) (V2, bool) {
	v1, ok := m[key]
	if !ok {
		var zero V2
		return zero, false
	}
	return v1.View(), true
}

func NewMapViewer[M ~map[K]V1, K comparable, V1 Viewer[V2], V2 any](m M) MapViewer[K, V1, V2] {
	return MapViewer[K, V1, V2](m)
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