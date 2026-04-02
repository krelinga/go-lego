package v2

import "iter"

type MapSeq[K comparable, V any] iter.Seq2[K, V]

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

type FixedMap[K comparable, V any] interface {
	Length() int
	Get(K) (V, bool)
	Range() MapSeq[K, V]
}

type Map[K comparable, V any] interface {
	FixedMap[K, V]
	Set(K, V)
	Remove(K)
}

type KV[K comparable, V any] struct {
	K K
	V V
}