package v2

import "iter"

type FixedList[P, V any] interface {
	Length() int
	Get(P) (V, bool)
	All() iter.Seq2[P, V]
	Positions() iter.Seq[P]
	Values() iter.Seq[V]
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

type FixedReversibleList[P, V any] interface {
	FixedList[P, V]
	ReverseAll() iter.Seq2[P, V]
	ReversePositions() iter.Seq[P]
	ReverseValues() iter.Seq[V]
}

type ReversibleList[P, V any] interface {
	List[P, V]
	FixedReversibleList[P, V]
}