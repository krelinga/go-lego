package v2

import "iter"

type Viewer[V any] interface {
	View() V
}

type DictViewEmbed[K any, T Viewer[V], V any] struct {
	fm FixedDict[K, T]
}

func (m DictViewEmbed[K, T, V]) Length() int {
	if m.fm == nil {
		return 0
	}
	return m.fm.Length()
}

func (m DictViewEmbed[K, T, V]) Get(k K) (V, bool) {
	if m.fm == nil {
		var zero V
		return zero, false
	}
	viewer, ok := m.fm.Get(k)
	if !ok {
		var zero V
		return zero, false
	}
	return viewer.View(), true
}

func (m DictViewEmbed[K, T, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m.fm == nil {
			return
		}
		for k, v := range m.fm.All() {
			if !yield(k, v.View()) {
				return
			}
		}
	}
}

func (m DictViewEmbed[K, T, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		if m.fm == nil {
			return
		}
		for k := range m.fm.Keys() {
			if !yield(k) {
				return
			}
		}
	}
}

func (m DictViewEmbed[K, T, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		if m.fm == nil {
			return
		}
		for v := range m.fm.Values() {
			if !yield(v.View()) {
				return
			}
		}
	}
}

func (m DictViewEmbed[K, T, V]) KVs() iter.Seq[KV[K, V]] {
	return func(yield func(KV[K, V]) bool) {
		if m.fm == nil {
			return
		}
		for k, v := range m.fm.All() {
			if !yield(NewKV(k, v.View())) {
				return
			}
		}
	}
}

func NewMapVieEmbed[K any, T Viewer[V], V any](m FixedDict[K, T]) DictViewEmbed[K, T, V] {
	return DictViewEmbed[K, T, V]{m}
}
