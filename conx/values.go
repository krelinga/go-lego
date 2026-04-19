package conx

import "iter"

type Vals[T any] interface {
	Len() int
	Vals() iter.Seq[T]
}

type OrdVals[T any] interface {
	Vals[T]
	RevVals() iter.Seq[T]
}
