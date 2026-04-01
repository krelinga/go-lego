package v2

import "iter"

type ListView[P, V any] interface {
	Length() int
	Get(P) (V, bool)
	List() iter.Seq2[P, V]
	First() (P, V, bool)
	Last() (P, V, bool)
}

type List[P, V any] interface {
	ListView[P, V]
	Set(P, V)
	InsertBefore(P, V)
	InsertAfter(P, V)
	Append(V)
	Remove(P)
}
