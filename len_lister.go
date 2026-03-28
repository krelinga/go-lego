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
func Keys[L LenLister[P], P FixedPair[K, V], K any, V any](l L) LenLister[K] {
	return keysLenLister[L, P, K, V]{l: l}
}

type keysLenLister[L LenLister[P], P FixedPair[K, V], K any, V any] struct {
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
func Values[L LenLister[P], P FixedPair[K, V], K any, V any](l L) LenLister[V] {
	return valuesLenLister[L, P, K, V]{l: l}
}

type valuesLenLister[L LenLister[P], P FixedPair[K, V], K any, V any] struct {
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

// ViewLenLister is a special case of [LenLister] that also implements the [Viewer] interface, so that it can provide a view of itself.
// V1 is the type of the elements in the LenLister, and V2 is the type of the elements in the view of the LenLister.
type ViewLenLister[V1, V2 any] interface {
	LenLister[V1]
	Viewer[LenLister[V2]]
}

// ViewerValues is a special case of [Values] that takes a [ViewLenLister] of pairs, and returns a [ViewLenLister] of the values in the pairs.
// The returned ViewLenLister will provide a view of the values in the pairs.
func ViewerValues[L ViewLenLister[P1, P2], P1 FixedPair[K, V1], P2 FixedPair[K, V2], K any, V1 Viewer[V2], V2 any](l L) ViewLenLister[V1, V2] {
	return viewerValuesLenLister[L, P1, P2, K, V1, V2]{l: l}
}

type viewerValuesLenLister[L ViewLenLister[P1, P2], P1 FixedPair[K, V1], P2 FixedPair[K, V2], K any, V1 Viewer[V2], V2 any] struct {
	l L
}

func (v viewerValuesLenLister[L, P1, P2, K, V1, V2]) Len() int {
	return v.l.Len()
}

func (v viewerValuesLenLister[L, P1, P2, K, V1, V2]) List() iter.Seq[V1] {
	return Values(v.l).List()
}

func (v viewerValuesLenLister[L, P1, P2, K, V1, V2]) View() LenLister[V2] {
	return Values(v.l.View())
}