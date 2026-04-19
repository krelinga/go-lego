package conx

import "iter"

type Vals[T any] interface {
	Len() int
	Vals() iter.Seq[T]
}
