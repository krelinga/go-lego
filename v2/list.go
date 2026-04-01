package v2

import "iter"

type FixedList[P, V any] interface {
	Length() int
	Get(P) (V, bool)
	List() iter.Seq2[P, V]
	First() (P, V, bool)
	Last() (P, V, bool)
}

type List[P, V any] interface {
	FixedList[P, V]
	Set(P, V)
	InsertBefore(P, V) P
	InsertAfter(P, V) P
	Append(V) P
	Remove(P)
}
