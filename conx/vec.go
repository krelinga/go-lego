package conx

import (
	"cmp"
	"slices"
)

type VecView[T any] interface {
	Len() int
	At(i int) T
}

type Vec[T any] []T

func CloneVec[T any](vec VecView[T]) *Vec[T] {
	return CloneVecFunc(vec, func(x T) T { return x })
}

func CloneVecFunc[T any, U any](vec VecView[T], valueFunc func(T) U) *Vec[U] {
	v := &Vec[U]{}
	v.Reserve(vec.Len())
	for i := 0; i < vec.Len(); i++ {
		v.Push(valueFunc(vec.At(i)))
	}
	return v
}

func (v *Vec[T]) Len() int {
	return len(*v)
}

func (v *Vec[T]) At(i int) T {
	return (*v)[i]
}

func (v *Vec[T]) Set(i int, value T) {
	(*v)[i] = value
}

func (v *Vec[T]) Clear() {
	*v = nil
}

func (v *Vec[T]) Reserve(n int) {
	if cap(*v) < n {
		newData := make([]T, len(*v), n)
		copy(newData, *v)
		*v = newData
	}
}

func (v *Vec[T]) Push(value T) {
	*v = append(*v, value)
}

func (v *Vec[T]) Pop() T {
	value := (*v)[len(*v)-1]
	*v = (*v)[:len(*v)-1]
	return value
}

func WrapVecValues[T, V any](vec VecView[T], wrap func(T) V) VecView[V] {
	return wrappedVecValues[T, V]{
		vec:  vec,
		wrap: wrap,
	}
}

type wrappedVecValues[T, V any] struct {
	vec  VecView[T]
	wrap func(T) V
}

func (w wrappedVecValues[T, V]) Len() int {
	return w.vec.Len()
}

func (w wrappedVecValues[T, V]) At(i int) V {
	return w.wrap(w.vec.At(i))
}

func VecEqualFunc[T any](a, b VecView[T], eq func(T, T) bool) bool {
	if a.Len() != b.Len() {
		return false
	}
	for i := 0; i < a.Len(); i++ {
		if !eq(a.At(i), b.At(i)) {
			return false
		}
	}
	return true
}

func VecEqual[T comparable](a, b VecView[T]) bool {
	return VecEqualFunc(a, b, func(x, y T) bool {
		return x == y
	})
}

func VecSort[T cmp.Ordered](vec *Vec[T]) {
	VecSortFunc(vec, cmp.Compare[T])
}

func VecSortFunc[T any](vec *Vec[T], order func(T, T) int) {
	slices.SortFunc(*vec, order)
}