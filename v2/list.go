package v2

import "iter"

type ListView[K1, V1, K2, V2 any] interface {
	Length() int
	Get(K1) (V1, bool)
	All() iter.Seq2[K1, V1]
	First() (K1, bool)
	Last() (K1, bool)
	Readonly() ListView[K2, V2, K2, V2]
}

type List[K1, V1, K2, V2 any] interface {
	ListView[K1, V1, K2, V2]
	InsertBefore(K1, V1)
	InsertAfter(K1, V1)
	Remove(K1)
}
