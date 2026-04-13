package conx

import "iter"

type Range[T any] interface {
	Len() int
	Range() iter.Seq[T]
}

type Range2[T, U any] interface {
	Len() int
	Range() iter.Seq2[T, U]
}

func RangeTransform[T, U any](r Range[T], f func(T) U) Range[U] {
	return rangeTransform[T, U]{r: r, f: f}
}

type rangeTransform[T, U any] struct {
	r Range[T]
	f func(T) U
}

func (r rangeTransform[T, U]) Len() int {
	return r.r.Len()
}

func (r rangeTransform[T, U]) Range() iter.Seq[U] {
	return func(yield func(U) bool) {
		for x := range r.r.Range() {
			if !yield(r.f(x)) {
				return
			}
		}
	}
}

func Range2Transform[T, U, V, W any](r Range2[T, U], f func(T, U) (V, W)) Range2[V, W] {
	return range2Transform[T, U, V, W]{r: r, f: f}
}

type range2Transform[T, U, V, W any] struct {
	r Range2[T, U]
	f func(T, U) (V, W)
}

func (r range2Transform[T, U, V, W]) Len() int {
	return r.r.Len()
}

func (r range2Transform[T, U, V, W]) Range() iter.Seq2[V, W] {
	return func(yield func(V, W) bool) {
		for x, y := range r.r.Range() {
			v, w := r.f(x, y)
			if !yield(v, w) {
				return
			}
		}
	}
}

func ToKVs[K, V any](kvs Range2[K, V]) Range[KV[K, V]] {
	return toKVsRange[K, V]{kvs: kvs}
}

type toKVsRange[K, V any] struct {
	kvs Range2[K, V]
}

func (r toKVsRange[K, V]) Len() int {
	return r.kvs.Len()
}

func (r toKVsRange[K, V]) Range() iter.Seq[KV[K, V]] {
	return func(yield func(KV[K, V]) bool) {
		for k, v := range r.kvs.Range() {
			if !yield(NewKV(k, v)) {
				return
			}
		}
	}
}

func FromKVs[K, V any](kvs Range[KV[K, V]]) Range2[K, V] {
	return fromKVsRange[K, V]{kvs: kvs}
}

type fromKVsRange[K, V any] struct {
	kvs Range[KV[K, V]]
}

func (r fromKVsRange[K, V]) Len() int {
	return r.kvs.Len()
}

func (r fromKVsRange[K, V]) Range() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for kv := range r.kvs.Range() {
			if !yield(kv.K, kv.V) {
				return
			}
		}
	}
}
