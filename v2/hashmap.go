package v2

import (
	"iter"
	"maps"
)

type HashMap[K comparable, V any] map[K]V

func (m *HashMap[K, V]) Length() int {
	return len(*m)
}

func (m *HashMap[K, V]) Get(k K) (V, bool) {
	v, ok := (*m)[k]
	return v, ok
}

func (m *HashMap[K, V]) All() iter.Seq2[K, V] {
	return maps.All(*m)
}

func (m *HashMap[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(*m)
}

func (m *HashMap[K, V]) Values() iter.Seq[V] {
	return maps.Values(*m)
}

func (m *HashMap[K, V]) KVs() iter.Seq[KV[K, V]] {
	return func(yield func(KV[K, V]) bool) {
		for k, v := range *m {
			if !yield(NewKV(k, v)) {
				return
			}
		}
	}
}

func (m *HashMap[K, V]) Set(k K, v V) {
	if *m == nil {
		*m = make(map[K]V)
	}
	(*m)[k] = v
}

func (m *HashMap[K, V]) Remove(k K) {
	delete(*m, k)
}

func (m *HashMap[K, V]) Reserve(n int) {
	if *m == nil {
		*m = make(map[K]V, n)
	}
}

func (m *HashMap[K, V]) Add(kv KV[K, V]) {
	m.Set(kv.K, kv.V)
}
