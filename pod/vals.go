package pod

import "iter"

type Vals[T any] interface {
	Len() int
	Vals() iter.Seq[T]
}

type CanReserve interface {
	Reserve(n int)
}

func ConcatVals[T any](vals ...Vals[T]) Vals[T] {
	return &concatVals[T]{vals: vals}
}

type concatVals[T any] struct {
	vals []Vals[T]
}

func (c *concatVals[T]) Len() int {
	total := 0
	for _, v := range c.vals {
		total += v.Len()
	}
	return total
}

func (c *concatVals[T]) Vals() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range c.vals {
			for curVal := range v.Vals() {
				if !yield(curVal) {
					return
				}
			}
		}
	}
}