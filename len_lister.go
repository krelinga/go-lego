package lego

import "iter"

// A LenLister is any type of container that has both a Len() and List() method.
// The container must have a defined length, so that multiple calls to Len() will return the same value, and the length must be non-negative.
// It must be safe to call List() multiple times on the container, and each call to List() must return the same sequence of elements (although not necessarily in the same order).
type LenLister[V any] interface {
	// Returns the number of elements in the container.
	Len() int

	// Returns a sequence of the elements in the container.
	List() iter.Seq[V]
}

// Keys returns a LenLister of the keys in the given LenLister of pairs.
func Keys[L LenLister[P], P FixedPair[K, V], K comparable, V any](l L) LenLister[K] {
	return keysLenLister[L, P, K, V]{l: l}
}

type keysLenLister[L LenLister[P], P FixedPair[K, V], K comparable, V any] struct {
	l L
}

func (k keysLenLister[L, P, K, V]) Len() int {
	return k.l.Len()
}

func (k keysLenLister[L, P, K, V]) List() iter.Seq[K] {
	return func(yield func(K) bool) {
		for pair := range k.l.List() {
			if !yield(pair.GetKey()) {
				return
			}
		}
	}
}

// Values returns a LenLister of the values in the given LenLister of pairs.
func Values[L LenLister[P], P FixedPair[K, V], K comparable, V any](l L) LenLister[V] {
	return valuesLenLister[L, P, K, V]{l: l}
}

type valuesLenLister[L LenLister[P], P FixedPair[K, V], K comparable, V any] struct {
	l L
}

func (v valuesLenLister[L, P, K, V]) Len() int {
	return v.l.Len()
}

func (v valuesLenLister[L, P, K, V]) List() iter.Seq[V] {
	return func(yield func(V) bool) {
		for pair := range v.l.List() {
			if !yield(pair.GetValue()) {
				return
			}
		}
	}
}

func ViewLenLister[L LenLister[V1], V1 Viewer[V2], V2 any](l L) LenLister[V2] {
	return lenListerView[L, V1, V2]{l: l}
}

type lenListerView[L LenLister[V1], V1 Viewer[V2], V2 any] struct {
	l L
}

func (v lenListerView[L, V1, V2]) Len() int {
	return v.l.Len()
}

func (v lenListerView[L, V1, V2]) List() iter.Seq[V2] {
	return func(yield func(V2) bool) {
		for elem := range v.l.List() {
			if !yield(elem.View()) {
				return
			}
		}
	}
}