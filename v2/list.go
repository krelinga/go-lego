package v2

import "iter"

type ListSeq[P, V any] iter.Seq2[P, V]

func (s ListSeq[P, V]) Positions() iter.Seq[P] {
	return func(yield func(P) bool) {
		for p := range s {
			if !yield(p) {
				return
			}
		}
	}
}

func (s ListSeq[P, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

type FixedList[P, V any] interface {
	Length() int
	Get(P) (V, bool)
	List() ListSeq[P, V]
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
