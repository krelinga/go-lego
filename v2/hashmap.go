package v2

import (
	"iter"
	"maps"
)

type HashMap[K comparable, V any] struct {
	m map[K]V
}

func NewHashMap[K comparable, V any](entries ...KV[K, V]) *HashMap[K, V] {
	m := &HashMap[K, V]{m: make(map[K]V)}
	for _, entry := range entries {
		m.m[entry.K] = entry.V
	}
	return m
}

func NewHashMapHint[K comparable, V any](n int, entries ...KV[K, V]) *HashMap[K, V] {
	m := &HashMap[K, V]{m: make(map[K]V, n)}
	for _, entry := range entries {
		m.m[entry.K] = entry.V
	}
	return m
}

func (m *HashMap[K, V]) Length() int {
	return len(m.m)
}

func (m *HashMap[K, V]) Get(k K) (V, bool) {
	v, ok := m.m[k]
	return v, ok
}

func (m *HashMap[K, V]) All() iter.Seq2[K, V] {
	return maps.All(m.m)
}

func (m *HashMap[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(m.m)
}

func (m *HashMap[K, V]) Values() iter.Seq[V] {
	return maps.Values(m.m)
}

func (m *HashMap[K, V]) KVs() iter.Seq[KV[K, V]] {
	return func(yield func(KV[K, V]) bool) {
		for k, v := range m.m {
			if !yield(NewKV(k, v)) {
				return
			}
		}
	}
}

func (m *HashMap[K, V]) Set(k K, v V) {
	if m.m == nil {
		m.m = make(map[K]V)
	}
	m.m[k] = v
}

func (m *HashMap[K, V]) Remove(k K) {
	delete(m.m, k)
}
