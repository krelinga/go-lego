package lego

import "reflect"

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

func ViewPair[P FixedPair[K, V1], K any, V1 Viewer[V2], V2 any](p P) FixedPair[K, V2] {
	if t := reflect.TypeFor[K](); t.Kind() == reflect.Pointer {
		panic("cannot create a view of a pair with pointer keys")
	}
	return pairView[P, K, V1, V2]{p: p}
}

type pairView[P FixedPair[K, V1], K any, V1 Viewer[V2], V2 any] struct {
	p P
}

func (v pairView[P, K, V1, V2]) GetKey() K {
	return v.p.GetKey()
}

func (v pairView[P, K, V1, V2]) GetValue() V2 {
	return v.p.GetValue().View()
}
