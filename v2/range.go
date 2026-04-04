package v2

import "iter"

type canLength interface {
	Length() int
}

type Range[V any] interface {
	Length() int
	All() iter.Seq[V]
}

type rangeImpl[C canLength, V any] struct {
	container C
	iterFunc  func(C) iter.Seq[V]
}

func (r rangeImpl[C, V]) Length() int {
	return r.container.Length()
}

func (r rangeImpl[C, V]) All() iter.Seq[V] {
	return r.iterFunc(r.container)
}

func newRangeImpl[C canLength, V any](container C, iterFunc func(C) iter.Seq[V]) rangeImpl[C, V] {
	return rangeImpl[C, V]{container, iterFunc}
}

type KeysContainer[K any] interface {
	Length() int
	Keys() iter.Seq[K]
}

func KeysFrom[K any](k KeysContainer[K]) Range[K] {
	return newRangeImpl(k, KeysContainer[K].Keys)
}

type ValuesContainer[V any] interface {
	Length() int
	Values() iter.Seq[V]
}

func ValuesFrom[V any](v ValuesContainer[V]) Range[V] {
	return newRangeImpl(v, ValuesContainer[V].Values)
}

type KVsContainer[K, V any] interface {
	Length() int
	KVs() iter.Seq[KV[K, V]]
}

func KVsFrom[K, V any](kv KVsContainer[K, V]) Range[KV[K, V]] {
	return newRangeImpl(kv, KVsContainer[K, V].KVs)
}

type PositionsContainer[P any] interface {
	Length() int
	Positions() iter.Seq[P]
}

func PositionsFrom[P any](p PositionsContainer[P]) Range[P] {
	return newRangeImpl(p, PositionsContainer[P].Positions)
}

type ReversePositionsContainer[P any] interface {
	Length() int
	ReversePositions() iter.Seq[P]
}

func ReversePositionsFrom[P any](p ReversePositionsContainer[P]) Range[P] {
	return newRangeImpl(p, ReversePositionsContainer[P].ReversePositions)
}

type ReverseValuesContainer[V any] interface {
	Length() int
	ReverseValues() iter.Seq[V]
}

func ReverseValuesFrom[V any](v ReverseValuesContainer[V]) Range[V] {
	return newRangeImpl(v, ReverseValuesContainer[V].ReverseValues)
}

func RangeFrom[V any](vs ...V) Range[V] {
	slice := Slice[V](vs)
	return ValuesFrom(&slice)
}