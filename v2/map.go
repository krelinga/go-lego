package v2

import "iter"

type FixedMap[K, V any] interface {
	Length() int
	Get(K) (V, bool)
	All() iter.Seq2[K, V]
	Keys() iter.Seq[K]
	Values() iter.Seq[V]
	KVs() iter.Seq[KV[K, V]]
}

type Map[K, V any] interface {
	FixedMap[K, V]
	Set(K, V)
	Remove(K)
}

type KV[K, V any] struct {
	K K
	V V
}

func NewKV[K, V any](k K, v V) KV[K, V] {
	return KV[K, V]{k, v}
}
