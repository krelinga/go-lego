package conx

import (
	"cmp"
	"slices"
)

type VectorView[T any] interface {
	Len() int
	Get(i int) T
}

type Vector[T any] []T

func CloneVector[T any](vec VectorView[T]) *Vector[T] {
	return CloneVectorFunc(vec, func(x T) T { return x })
}

func CloneVectorFunc[T any, U any](vec VectorView[T], valueFunc func(T) U) *Vector[U] {
	v := &Vector[U]{}
	v.Reserve(vec.Len())
	for i := 0; i < vec.Len(); i++ {
		v.Push(valueFunc(vec.Get(i)))
	}
	return v
}

func (v *Vector[T]) Len() int {
	return len(*v)
}

func (v *Vector[T]) Get(i int) T {
	return (*v)[i]
}

func (v *Vector[T]) Set(i int, value T) {
	(*v)[i] = value
}

func (v *Vector[T]) Clear() {
	*v = nil
}

func (v *Vector[T]) Reserve(n int) {
	if cap(*v) < n {
		newData := make([]T, len(*v), n)
		copy(newData, *v)
		*v = newData
	}
}

func (v *Vector[T]) Push(value T) {
	*v = append(*v, value)
}

func (v *Vector[T]) Pop() T {
	value := (*v)[len(*v)-1]
	*v = (*v)[:len(*v)-1]
	return value
}

func WrapVectorValues[T, V any](vec VectorView[T], wrap func(T) V) VectorView[V] {
	return wrappedVectorValues[T, V]{
		vec:  vec,
		wrap: wrap,
	}
}

type wrappedVectorValues[T, V any] struct {
	vec  VectorView[T]
	wrap func(T) V
}

func (w wrappedVectorValues[T, V]) Len() int {
	return w.vec.Len()
}

func (w wrappedVectorValues[T, V]) Get(i int) V {
	return w.wrap(w.vec.Get(i))
}

func NewVectorViewFromSlice[C ~[]T, T any](slice C) VectorView[T] {
	return sliceVectorView[C, T]{slice: slice}
}

type sliceVectorView[C ~[]T, T any] struct {
	slice C
}

func (s sliceVectorView[C, T]) Len() int {
	return len(s.slice)
}

func (s sliceVectorView[C, T]) Get(i int) T {
	return s.slice[i]
}

func VectorEqualFunc[T any](a, b VectorView[T], eq func(T, T) bool) bool {
	if a.Len() != b.Len() {
		return false
	}
	for i := 0; i < a.Len(); i++ {
		if !eq(a.Get(i), b.Get(i)) {
			return false
		}
	}
	return true
}

func VectorEqual[T comparable](a, b VectorView[T]) bool {
	return VectorEqualFunc(a, b, func(x, y T) bool {
		return x == y
	})
}

func VectorSort[T cmp.Ordered](vec *Vector[T]) {
	VectorSortFunc(vec, cmp.Compare[T])
}

func VectorSortFunc[T any](vec *Vector[T], order func(T, T) int) {
	slices.SortFunc(*vec, order)
}