package v2

import "iter"

type FixedDict[K, V any] interface {
	Length() int
	Get(K) (V, bool)
	All() iter.Seq2[K, V]
	Keys() iter.Seq[K]
	Values() iter.Seq[V]
	KVs() iter.Seq[KV[K, V]]
}

type Dict[K, V any] interface {
	FixedDict[K, V]
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
