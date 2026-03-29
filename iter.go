package lego

import "iter"

func Seq[K comparable, V any](s iter.Seq2[K, V]) iter.Seq[Pair[K, V]] {
	return func(yield func(Pair[K, V]) bool) {
		for k, v := range s {
			if !yield(NewPair(k, v)) {
				return
			}
		}
	}
}

func Seq2[P FixedPair[K, V], K comparable, V any](s iter.Seq[P]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for pair := range s {
			if !yield(pair.GetKey(), pair.GetValue()) {
				return
			}
		}
	}
}
