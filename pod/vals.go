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

type KeyVals[K, V any] interface {
	Vals[V]
	Keys() iter.Seq[K]
	KeyVals() iter.Seq2[K, V]
}

func ConcatKeyVals[K, V any](keyVals ...KeyVals[K, V]) KeyVals[K, V] {
	return &concatKeyVals[K, V]{keyVals: keyVals}
}

type concatKeyVals[K, V any] struct {
	keyVals []KeyVals[K, V]
}

func (c *concatKeyVals[K, V]) Len() int {
	total := 0
	for _, kv := range c.keyVals {
		total += kv.Len()
	}
	return total
}

func (c *concatKeyVals[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, kv := range c.keyVals {
			for k := range kv.Keys() {
				if !yield(k) {
					return
				}
			}
		}
	}
}

func (c *concatKeyVals[K, V]) KeyVals() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, kv := range c.keyVals {
			for k, v := range kv.KeyVals() {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

func (c *concatKeyVals[K, V]) Vals() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, kv := range c.keyVals {
			for v := range kv.Vals() {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func KeyAsVal[K, V any](keyVals KeyVals[K, V]) Vals[K] {
	return &keyAsVal[K, V]{keyVals: keyVals}
}

type keyAsVal[K, V any] struct {
	keyVals KeyVals[K, V]
}

func (k *keyAsVal[K, V]) Len() int {
	return k.keyVals.Len()
}

func (k *keyAsVal[K, V]) Vals() iter.Seq[K] {
	return k.keyVals.Keys()
}
