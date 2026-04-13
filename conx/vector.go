package conx

import "iter"

type VectorView[T any] interface {
	Len() int
	Get(i int) T
}

type Vector[T any] struct {
	data []T
}

func NewVector[T any](data ...T) *Vector[T] {
	return &Vector[T]{data: data}
}

func CloneVector[T any](vec VectorView[T]) *Vector[T] {
	v := &Vector[T]{}
	v.Reserve(vec.Len())
	for i := 0; i < vec.Len(); i++ {
		v.Push(vec.Get(i))
	}
	return v
}

func (v Vector[T]) Len() int {
	return len(v.data)
}

func (v Vector[T]) Get(i int) T {
	return v.data[i]
}

func (v *Vector[T]) Set(i int, value T) {
	v.data[i] = value
}

func (v *Vector[T]) Clear() {
	v.data = nil
}

func (v *Vector[T]) Reserve(n int) {
	if cap(v.data) < n {
		newData := make([]T, len(v.data), n)
		copy(newData, v.data)
		v.data = newData
	}
}

func (v *Vector[T]) Push(value T) {
	v.data = append(v.data, value)
}

func (v *Vector[T]) Pop() T {
	value := v.data[len(v.data)-1]
	v.data = v.data[:len(v.data)-1]
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

func VectorAll[T any](vec VectorView[T]) Range2[int, T] {
	return vectorAllRange[T]{vec: vec}
}

type vectorAllRange[T any] struct {
	vec VectorView[T]
}

func (r vectorAllRange[T]) Len() int {
	return r.vec.Len()
}

func (r vectorAllRange[T]) Range() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := 0; i < r.vec.Len(); i++ {
			if !yield(i, r.vec.Get(i)) {
				return
			}
		}
	}
}

func VectorValues[T any](vec VectorView[T]) Range[T] {
	return vectorValuesRange[T]{vec: vec}
}

type vectorValuesRange[T any] struct {
	vec VectorView[T]
}

func (r vectorValuesRange[T]) Len() int {
	return r.vec.Len()
}

func (r vectorValuesRange[T]) Range() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < r.vec.Len(); i++ {
			if !yield(r.vec.Get(i)) {
				return
			}
		}
	}
}

func VectorReverseAll[T any](vec VectorView[T]) Range2[int, T] {
	return vectorReverseAllRange[T]{vec: vec}
}

type vectorReverseAllRange[T any] struct {
	vec VectorView[T]
}

func (r vectorReverseAllRange[T]) Len() int {
	return r.vec.Len()
}

func (r vectorReverseAllRange[T]) Range() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := r.vec.Len() - 1; i >= 0; i-- {
			if !yield(i, r.vec.Get(i)) {
				return
			}
		}
	}
}

func VectorReverseValues[T any](vec VectorView[T]) Range[T] {
	return vectorReverseValuesRange[T]{vec: vec}
}

type vectorReverseValuesRange[T any] struct {
	vec VectorView[T]
}

func (r vectorReverseValuesRange[T]) Len() int {
	return r.vec.Len()
}

func (r vectorReverseValuesRange[T]) Range() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := r.vec.Len() - 1; i >= 0; i-- {
			if !yield(r.vec.Get(i)) {
				return
			}
		}
	}
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
