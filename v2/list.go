package v2

import "iter"

type ListView[K, V any] interface {
	Length() int
	Get(K) (V, bool)
	List() iter.Seq2[K, V]
	First() (K, V, bool)
	Last() (K, V, bool)
}

type List[K, V any] interface {
	ListView[K, V]
	Set(K, V)
	InsertBefore(K, V)
	InsertAfter(K, V)
	Append(V)
	Remove(K)
}
