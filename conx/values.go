package conx

import "iter"

type Values[T any] interface {
	Len() int
	Values() iter.Seq[T]
}

type OrderedValues[T any] interface {
	Values[T]
	ReverseValues() iter.Seq[T]
}
