package v2

import "maps"

type GoMap[K comparable, V any] struct {
	m map[K]V
}

func (m *GoMap[K, V]) Length() int {
	return len(m.m)
}

func (m *GoMap[K, V]) Get(k K) (V, bool) {
	v, ok := m.m[k]
	return v, ok
}

func (m *GoMap[K, V]) Range() MapSeq[K, V] {
	return MapSeq[K, V](maps.All(m.m))
}

func (m *GoMap[K, V]) Set(k K, v V) {
	if m.m == nil {
		m.m = make(map[K]V)
	}
	m.m[k] = v
}

func (m *GoMap[K, V]) Remove(k K) {
	delete(m.m, k)
}
