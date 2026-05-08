package pod

import "iter"

// Bag is a reference type that provides a view of a collection of elements. It does not own the
// elements, but provides methods to access them. Any mutations to the underlying collection will be
// reflected in the Bag.
type Bag[T any] struct {
	lenFn   func() int
	elemsFn func() iter.Seq[T]
}

// Len returns the number of elements in the Bag. If the underlying collection is nil, it returns 0.
func (b Bag[T]) Len() int {
	if b.lenFn == nil {
		return 0
	}
	return b.lenFn()
}

// Elems returns a sequence of elements in the Bag. If the underlying collection is nil, it returns
// an empty sequence.
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

// Bag2 is a reference type that provides a view of a collection of pairs. It does not own the
// elements, but provides methods to access them. Any mutations to the underlying collection will be
// reflected in the Bag2.
type Bag2[T1, T2 any] struct {
	lenFn   func() int
	elemsFn func() iter.Seq2[T1, T2]
}

// Len returns the number of elements in the Bag2. If the underlying collection is nil, it returns
// 0.
func (b Bag2[T1, T2]) Len() int {
	if b.lenFn == nil {
		return 0
	}
	return b.lenFn()
}

// Elems returns a sequence of pairs in the Bag2. If the underlying collection is nil, it returns an
// empty sequence.
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

// Vals is an interface that provides a view of a collection of values.
type Vals[T any] interface {
	Len() int
	Vals() iter.Seq[T]
}

// RevVals is an interface that provides a view of a collection of values in reverse order.
type RevVals[T any] interface {
	Len() int
	RevVals() iter.Seq[T]
}

// Keys is an interface that provides a view of a collection of keys.
type Keys[K any] interface {
	Len() int
	Keys() iter.Seq[K]
}

// KeyVals is an interface that provides a view of a collection of key-value pairs.
type KeyVals[K, V any] interface {
	Len() int
	KeyVals() iter.Seq2[K, V]
}

// IdxVals is an interface that provides a view of a collection of indexed values.
type IdxVals[T any] interface {
	Len() int
	IdxVals() iter.Seq2[int, T]
}

// IdxValsOf creates a Bag2 from an IdxVals. The Bag2 will reflect any mutations to the underlying
// collection.
func IdxValsOf[T any](v IdxVals[T]) Bag2[int, T] {
	return newBag2(v.Len, v.IdxVals)
}

// IdxValsOfFunc creates a Bag2 from an IdxVals, applying a transformation function to each element.
// The Bag2 will reflect any mutations to the underlying collection.
func IdxValsOfFunc[T1, T2 any](v IdxVals[T1], f func(int, T1) (int, T2)) Bag2[int, T2] {
	return newBag2(
		v.Len,
		func() iter.Seq2[int, T2] {
			return func(yield func(int, T2) bool) {
				for i, t1 := range v.IdxVals() {
					if !yield(f(i, t1)) {
						return
					}
				}
			}
		},
	)
}

// RevIdxVals is an interface that provides a view of a collection of indexed values in reverse
// order.
type RevIdxVals[T any] interface {
	Len() int
	RevIdxVals() iter.Seq2[int, T]
}

// RevIdxValsOf creates a Bag2 from a RevIdxVals. The Bag2 will reflect any mutations to the
// underlying collection.
func RevIdxValsOf[T any](v RevIdxVals[T]) Bag2[int, T] {
	return newBag2(v.Len, v.RevIdxVals)
}

// RevIdxValsOfFunc creates a Bag2 from a RevIdxVals, applying a transformation function to each
// element. The Bag2 will reflect any mutations to the underlying collection.
func RevIdxValsOfFunc[T1, T2 any](v RevIdxVals[T1], f func(int, T1) (int, T2)) Bag2[int, T2] {
	return newBag2(
		v.Len,
		func() iter.Seq2[int, T2] {
			return func(yield func(int, T2) bool) {
				for i, t1 := range v.RevIdxVals() {
					if !yield(f(i, t1)) {
						return
					}
				}
			}
		},
	)
}

// RevValsOf creates a Bag from a RevVals. The Bag will reflect any mutations to the underlying
// collection.
func RevValsOf[T any](v RevVals[T]) Bag[T] {
	return newBag(v.Len, v.RevVals)
}

// RevValsOfFunc creates a Bag from a RevVals, applying a transformation function to each element.
// The Bag will reflect any mutations to the underlying collection.
func RevValsOfFunc[T1, T2 any](v RevVals[T1], f func(T1) T2) Bag[T2] {
	return newBag(
		v.Len,
		func() iter.Seq[T2] {
			return func(yield func(T2) bool) {
				for t1 := range v.RevVals() {
					if !yield(f(t1)) {
						return
					}
				}
			}
		},
	)
}

// ValsOf creates a Bag from a Vals. The Bag will reflect any mutations to the underlying
// collection.
func ValsOf[T any](v Vals[T]) Bag[T] {
	return newBag(v.Len, v.Vals)
}

// ValsOfFunc creates a Bag from a Vals, applying a transformation function to each element. The Bag
// will reflect any mutations to the underlying collection.
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

// KeysOf creates a Bag from a Keys. The Bag will reflect any mutations to the underlying
// collection.
func KeysOf[K any](k Keys[K]) Bag[K] {
	return newBag(k.Len, k.Keys)
}

// KeysOfFunc creates a Bag from a Keys, applying a transformation function to each element. The Bag
// will reflect any mutations to the underlying collection.
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

// KeyValsOf creates a Bag2 from a KeyVals. The Bag2 will reflect any mutations to the underlying
// collection.
func KeyValsOf[K, V any](kv KeyVals[K, V]) Bag2[K, V] {
	return newBag2(kv.Len, kv.KeyVals)
}

// KeyValsOfFunc creates a Bag2 from a KeyVals, applying a transformation function to each element.
// The Bag2 will reflect any mutations to the underlying collection.
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
