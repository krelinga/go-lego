package v2

import (
	"iter"
	"slices"
)

type Slice[V any] []V

func (l *Slice[V]) Length() int {
	return len(*l)
}

func (l *Slice[V]) Get(p int) (V, bool) {
	if p < 0 || p >= len(*l) {
		var zero V
		return zero, false
	}
	return (*l)[p], true
}

func (l *Slice[V]) All() iter.Seq2[int, V] {
	return slices.All(*l)
}

func (l *Slice[V]) Positions() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range *l {
			if !yield(i) {
				return
			}
		}
	}
}

func (l *Slice[V]) Values() iter.Seq[V] {
	return slices.Values(*l)
}

func (l *Slice[V]) ReverseAll() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i := len(*l) - 1; i >= 0; i-- {
			if !yield(i, (*l)[i]) {
				return
			}
		}
	}
}

func (l *Slice[V]) ReversePositions() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := len(*l) - 1; i >= 0; i-- {
			if !yield(i) {
				return
			}
		}
	}
}

func (l *Slice[V]) ReverseValues() iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := len(*l) - 1; i >= 0; i-- {
			if !yield((*l)[i]) {
				return
			}
		}
	}
}

func (l *Slice[V]) First() (int, V, bool) {
	if len(*l) == 0 {
		var zero V
		return 0, zero, false
	}
	return 0, (*l)[0], true
}

func (l *Slice[V]) Last() (int, V, bool) {
	if len(*l) == 0 {
		var zero V
		return 0, zero, false
	}
	lastIdx := len(*l) - 1
	return lastIdx, (*l)[lastIdx], true
}

func (l *Slice[V]) String() string {
	return listStringHelper(l)
}

func (l *Slice[V]) Set(p int, v V) {
	if p < 0 || p >= len(*l) {
		panic("index out of bounds")
	}
	(*l)[p] = v
}

func (l *Slice[V]) InsertBefore(p int, v V) int {
	if p < 0 || p > len(*l) {
		panic("index out of bounds")
	}
	*l = append((*l)[:p], append([]V{v}, (*l)[p:]...)...)
	return p
}

func (l *Slice[V]) InsertAfter(p int, v V) int {
	if p < 0 || p >= len(*l) {
		panic("index out of bounds")
	}
	*l = append((*l)[:p+1], append([]V{v}, (*l)[p+1:]...)...)
	return p + 1
}

func (l *Slice[V]) Add(v V) {
	*l = append(*l, v)
}

func (l *Slice[V]) Remove(p int) {
	if p < 0 || p >= len(*l) {
		panic("index out of bounds")
	}
	*l = append((*l)[:p], (*l)[p+1:]...)
}

func (l *Slice[V]) Reserve(n int) {
	if *l == nil {
		*l = make([]V, 0, n)
	}
}
