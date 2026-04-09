package v2

import (
	"fmt"
	"iter"
	"strings"
)

type FixedList[P, V any] interface {
	Length() int
	Get(P) (V, bool)
	All() iter.Seq2[P, V]
	Positions() iter.Seq[P]
	Values() iter.Seq[V]
	First() (P, V, bool)
	Last() (P, V, bool)
	String() string
}

type List[P, V any] interface {
	FixedList[P, V]
	Set(P, V)
	InsertBefore(P, V) P
	InsertAfter(P, V) P
	Add(V)
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

func listStringHelper[P, V any](l FixedList[P, V]) string {
	var sb strings.Builder
	sb.WriteString("[")
	first := true
	for v := range l.Values() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		fmt.Fprintf(&sb, "%v", v)
	}
	sb.WriteString("]")
	return sb.String()
}
