package v2

import "iter"

type MapSeq[K, V any] iter.Seq2[K, V]

func (s MapSeq[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range s {
			if !yield(k) {
				return
			}
		}
	}
}

func (s MapSeq[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

type FixedMap[K, V any] interface {
	Length() int
	Get(K) (V, bool)
	Range() MapSeq[K, V]
}

type Map[K, V any] interface {
	FixedMap[K, V]
	Set(K, V)
	Remove(K)
}
