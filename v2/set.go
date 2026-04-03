package v2

import "iter"

type FixedSet[V any] interface {
	Length() int
	Has(V) bool
	Values() iter.Seq[V]
	String() string
}

type Set[V any] interface {
	FixedSet[V]
	Add(V)
	Remove(V)
}
