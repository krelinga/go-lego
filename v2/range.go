package v2

import "iter"

type CanLength interface {
	Length() int
}

type Range[V any] interface {
	Length() int
	All() iter.Seq[V]
}

type rangeImpl[C CanLength, V any] struct {
	container C
	iterFunc  func(C) iter.Seq[V]
}

func (r rangeImpl[C, V]) Length() int {
	return r.container.Length()
}

func (r rangeImpl[C, V]) All() iter.Seq[V] {
	return r.iterFunc(r.container)
}

func newRangeImpl[C CanLength, V any](container C, iterFunc func(C) iter.Seq[V]) rangeImpl[C, V] {
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
