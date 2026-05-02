package pod

import "iter"

type Bag[T any] struct {
	lenFn   func() int
	elemsFn func() iter.Seq[T]
}

func (b Bag[T]) Len() int {
	if b.lenFn == nil {
		return 0
	}
	return b.lenFn()
}

func (b Bag[T]) Elems() iter.Seq[T] {
	if b.elemsFn == nil {
		return func(yield func(T) bool) {}
	}
	return b.elemsFn()
}

func newBag[T any](lenFn func() int, elemsFn func() iter.Seq[T]) Bag[T] {
	return Bag[T]{
		lenFn:   lenFn,
		elemsFn: elemsFn,
	}
}

type Bag2[T1, T2 any] struct {
	lenFn   func() int
	elemsFn func() iter.Seq2[T1, T2]
}

func (b Bag2[T1, T2]) Len() int {
	if b.lenFn == nil {
		return 0
	}
	return b.lenFn()
}

func (b Bag2[T1, T2]) Elems() iter.Seq2[T1, T2] {
	if b.elemsFn == nil {
		return func(yield func(T1, T2) bool) {}
	}
	return b.elemsFn()
}

func newBag2[T1, T2 any](lenFn func() int, elemsFn func() iter.Seq2[T1, T2]) Bag2[T1, T2] {
	return Bag2[T1, T2]{
		lenFn:   lenFn,
		elemsFn: elemsFn,
	}
}

type Vals[T any] interface {
	Len() int
	Vals() iter.Seq[T]
}

type Keys[K any] interface {
	Len() int
	Keys() iter.Seq[K]
}

type KeyVals[K, V any] interface {
	Len() int
	KeyVals() iter.Seq2[K, V]
}

func ValsOf[T any](v Vals[T]) Bag[T] {
	return newBag(v.Len, v.Vals)
}

func ValsOfFunc[T1, T2 any](v Vals[T1], f func(T1) T2) Bag[T2] {
	return newBag(
		v.Len,
		func() iter.Seq[T2] {
			return func(yield func(T2) bool) {
				for t1 := range v.Vals() {
					if !yield(f(t1)) {
						return
					}
				}
			}
		},
	)
}

func KeysOf[K any](k Keys[K]) Bag[K] {
	return newBag(k.Len, k.Keys)
}

func KeysOfFunc[K1, K2 any](k Keys[K1], f func(K1) K2) Bag[K2] {
	return newBag(
		k.Len,
		func() iter.Seq[K2] {
			return func(yield func(K2) bool) {
				for k1 := range k.Keys() {
					if !yield(f(k1)) {
						return
					}
				}
			}
		},
	)
}

func KeyValsOf[K, V any](kv KeyVals[K, V]) Bag2[K, V] {
	return newBag2(kv.Len, kv.KeyVals)
}

func KeyValsOfFunc[K1, V1, K2, V2 any](kv KeyVals[K1, V1], f func(K1, V1) (K2, V2)) Bag2[K2, V2] {
	return newBag2(
		kv.Len,
		func() iter.Seq2[K2, V2] {
			return func(yield func(K2, V2) bool) {
				for k1, v1 := range kv.KeyVals() {
					if !yield(f(k1, v1)) {
						return
					}
				}
			}
		},
	)
}
