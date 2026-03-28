package lego

// FixedPair is a pair that does not allow callers to reassign the key or value fields, but which may still allow modifying the key and value (for example, if the key or value is a pointer).
type FixedPair[K any, V any] interface {
	GetKey() K
	GetValue() V
}

// A Pair is a mutable pair that implements the [FixedPair] interface.
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

// A ViewerPair is a special case of [Pair] that stores a value that implements the [Viewer] interface, and provides methods to get the key and a view of the value.
type ViewerPair[K any, V1 Viewer[V2], V2 any] Pair[K, V1]

func (p ViewerPair[K, V1, V2]) GetKey() K {
	return Pair[K, V1](p).GetKey()
}

func (p ViewerPair[K, V1, V2]) GetValue() V1 {
	return Pair[K, V1](p).GetValue()
}

func (p ViewerPair[K, V1, V2]) View() FixedPair[K, V2] {
	return viewerPairView[K, V1, V2]{p: p}
}

func NewViewerPair[K any, V1 Viewer[V2], V2 any](key K, value V1) ViewerPair[K, V1, V2] {
	return ViewerPair[K, V1, V2]{Key: key, Value: value}
}

type viewerPairView[K any, V1 Viewer[V2], V2 any] struct {
	p ViewerPair[K, V1, V2]
}

func (v viewerPairView[K, V1, V2]) GetKey() K {
	return v.p.GetKey()
}

func (v viewerPairView[K, V1, V2]) GetValue() V2 {
	return v.p.GetValue().View()
}
