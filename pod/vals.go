package pod

import "iter"

type Vals[T any] interface {
	Len() int
	Vals() iter.Seq[T]
}

type CanReserve interface {
	Reserve(n int)
}