package lego

type PairView[K any, V any] interface {
	GetKey() K
	GetValue() V
}

type Pair[K any, V any] struct {
	Key   K
	Value V
}

func (p Pair[K, V]) GetKey() K {
	return p.Key
}

func (p Pair[K, V]) GetValue() V {
	return p.Value
}

func NewPair[K any, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{Key: key, Value: value}
}

type PairViewer[K any, V1 Viewer[V2], V2 any] Pair[K, V1]

func (p PairViewer[K, V1, V2]) GetKey() K {
	return Pair[K, V1](p).GetKey()
}

func (p PairViewer[K, V1, V2]) GetValue() V2 {
	return p.Value.View()
}

func NewPairViewer[K any, V1 Viewer[V2], V2 any](key K, value V1) PairViewer[K, V1, V2] {
	return PairViewer[K, V1, V2]{Key: key, Value: value}
}