package v2

import "maps"

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

func (m *HashMap[K, V]) Range() MapSeq[K, V] {
	return MapSeq[K, V](maps.All(m.m))
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
