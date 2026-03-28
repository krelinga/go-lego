package lego

import "iter"

type Lister[V any] interface {
	List() iter.Seq[V]
}

func Keys[L Lister[P], P FixedPair[K, V], K any, V any](l L) iter.Seq[K] {
	return func(yield func(K) bool) {
		for pair := range l.List() {
			if !yield(pair.GetKey()) {
				return
			}
		}
	}
}

func Values[L Lister[P], P FixedPair[K, V], K any, V any](l L) iter.Seq[V] {
	return func(yield func(V) bool) {
		for pair := range l.List() {
			if !yield(pair.GetValue()) {
				return
			}
		}
	}
}

func All[L Lister[V], V any](l L) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range l.List() {
			if !yield(v) {
				return
			}
		}
	}
}

func All2[L Lister[P], P FixedPair[K, V], K any, V any](l L) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for pair := range l.List() {
			if !yield(pair.GetKey(), pair.GetValue()) {
				return
			}
		}
	}
}
