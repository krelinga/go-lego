package pod

import (
	"cmp"
	"iter"
	"slices"
)

type VecView[T any] interface {
	Len() int
	Get(i int) T
	Vals() iter.Seq[T]
	IdxVals() iter.Seq2[int, T]
	RevVals() iter.Seq[T]
	RevIdxVals() iter.Seq2[int, T]
}

type Vec[T any] interface {
	VecView[T]
	Set(i int, value T)
	Clear()
	Push(value T)
	Pop() T
	PushVals(vals Vals[T])
}

func AsVec[V ~[]T, T any](v V) VecView[T] {
	return vecView[V, T]{v: v}
}

type vecView[V ~[]T, T any] struct {
	v V
}

func (v vecView[V, T]) Len() int {
	return len(v.v)
}

func (v vecView[V, T]) Get(i int) T {
	return v.v[i]
}

func (v vecView[V, T]) Vals() iter.Seq[T] {
	return slices.Values(v.v)
}

func (v vecView[V, T]) IdxVals() iter.Seq2[int, T] {
	return slices.All(v.v)
}

func (v vecView[V, T]) RevVals() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(v.v) - 1; i >= 0; i-- {
			if !yield(v.v[i]) {
				return
			}
		}
	}
}

func (v vecView[V, T]) RevIdxVals() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := len(v.v) - 1; i >= 0; i-- {
			if !yield(i, v.v[i]) {
				return
			}
		}
	}
}

type Slice[T any] []T

func CloneValsIntoVec[T any](vals Vals[T], out Vec[T]) {
	CloneValsIntoVecFunc(vals, out, func(x T) T { return x })
}

func CloneValsIntoVecFunc[T any, U any](vals Vals[T], out Vec[U], valueFunc func(T) U) {
	out.Clear()
	if canReserve, ok := out.(CanReserve); ok {
		canReserve.Reserve(vals.Len())
	}
	for v := range vals.Vals() {
		out.Push(valueFunc(v))
	}
}

func (v *Slice[T]) Len() int {
	return len(*v)
}

func (v *Slice[T]) Get(i int) T {
	return (*v)[i]
}

func (v *Slice[T]) Vals() iter.Seq[T] {
	return slices.Values(*v)
}

func (v *Slice[T]) IdxVals() iter.Seq2[int, T] {
	return slices.All(*v)
}

func (v *Slice[T]) RevVals() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(*v) - 1; i >= 0; i-- {
			if !yield((*v)[i]) {
				return
			}
		}
	}
}

func (v *Slice[T]) RevIdxVals() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := len(*v) - 1; i >= 0; i-- {
			if !yield(i, (*v)[i]) {
				return
			}
		}
	}
}

func (v *Slice[T]) Set(i int, value T) {
	(*v)[i] = value
}

func (v *Slice[T]) Clear() {
	*v = nil
}

func (v *Slice[T]) Reserve(n int) {
	if cap(*v) < n {
		newData := make([]T, len(*v), n)
		copy(newData, *v)
		*v = newData
	}
}

func (v *Slice[T]) Push(value T) {
	*v = append(*v, value)
}

func (v *Slice[T]) PushVals(vals Vals[T]) {
	v.Reserve(v.Len() + vals.Len())
	for value := range vals.Vals() {
		v.Push(value)
	}
}

func (v *Slice[T]) Pop() T {
	value := (*v)[len(*v)-1]
	*v = (*v)[:len(*v)-1]
	return value
}

func WrapVecVals[T, V any](vec VecView[T], wrap func(T) V) VecView[V] {
	return wrappedVecVals[T, V]{
		vec:  vec,
		wrap: wrap,
	}
}

type wrappedVecVals[T, V any] struct {
	vec  VecView[T]
	wrap func(T) V
}

func (w wrappedVecVals[T, V]) Len() int {
	return w.vec.Len()
}

func (w wrappedVecVals[T, V]) Get(i int) V {
	return w.wrap(w.vec.Get(i))
}

func (w wrappedVecVals[T, V]) Vals() iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range w.vec.Vals() {
			if !yield(w.wrap(v)) {
				return
			}
		}
	}
}

func (w wrappedVecVals[T, V]) IdxVals() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i, v := range w.vec.IdxVals() {
			if !yield(i, w.wrap(v)) {
				return
			}
		}
	}
}

func (w wrappedVecVals[T, V]) RevVals() iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range w.vec.RevVals() {
			if !yield(w.wrap(v)) {
				return
			}
		}
	}
}

func (w wrappedVecVals[T, V]) RevIdxVals() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i, v := range w.vec.RevIdxVals() {
			if !yield(i, w.wrap(v)) {
				return
			}
		}
	}
}

func VecEqualFunc[T any](a, b VecView[T], eq func(T, T) bool) bool {
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

func VecEqual[T comparable](a, b VecView[T]) bool {
	return VecEqualFunc(a, b, func(x, y T) bool {
		return x == y
	})
}

func VecSort[T cmp.Ordered](vec Vec[T]) {
	VecSortFunc(vec, cmp.Compare[T])
}

func VecSortFunc[T any](vec Vec[T], order func(T, T) int) {
	if slice, ok := vec.(*Slice[T]); ok {
		slices.SortFunc(*slice, order)
		return
	}

	// TODO: make this more-efficient if/when we have non-Slice Vec implementations.
	temp := make([]T, vec.Len())
	for i := 0; i < vec.Len(); i++ {
		temp[i] = vec.Get(i)
	}
	slices.SortFunc(temp, order)
	for i := 0; i < vec.Len(); i++ {
		vec.Set(i, temp[i])
	}
}
