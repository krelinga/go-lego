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

func NewMap[O Map[K, V], K, V any](kvs ...KV[K, V]) *O {
	var m O
	for _, kv := range kvs {
		m.Set(kv.K, kv.V)
	}
	return &m
}
