package v2

import "iter"

type CanView[V any] interface {
	View() V
}

type DictViewEmbed[K any, T CanView[V], V any] struct {
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

func (m DictViewEmbed[K, T, V]) String() string {
	return m.fm.String()
}

func NewDictViewEmbed[K any, T CanView[V], V any](m FixedDict[K, T]) DictViewEmbed[K, T, V] {
	return DictViewEmbed[K, T, V]{m}
}

type ListViewEmbed[P any, T CanView[V], V any] struct {
	fl FixedList[P, T]
}

func (l ListViewEmbed[P, T, V]) Length() int {
	if l.fl == nil {
		return 0
	}
	return l.fl.Length()
}

func (l ListViewEmbed[P, T, V]) Get(p P) (V, bool) {
	if l.fl == nil {
		var zero V
		return zero, false
	}
	viewer, ok := l.fl.Get(p)
	if !ok {
		var zero V
		return zero, false
	}
	return viewer.View(), true
}

func (l ListViewEmbed[P, T, V]) All() iter.Seq2[P, V] {
	return func(yield func(P, V) bool) {
		if l.fl == nil {
			return
		}
		for p, v := range l.fl.All() {
			if !yield(p, v.View()) {
				return
			}
		}
	}
}

func (l ListViewEmbed[P, T, V]) Positions() iter.Seq[P] {
	return func(yield func(P) bool) {
		if l.fl == nil {
			return
		}
		for p := range l.fl.Positions() {
			if !yield(p) {
				return
			}
		}
	}
}

func (l ListViewEmbed[P, T, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		if l.fl == nil {
			return
		}
		for v := range l.fl.Values() {
			if !yield(v.View()) {
				return
			}
		}
	}
}

func (l ListViewEmbed[P, T, V]) String() string {
	return l.fl.String()
}

func NewListViewEmbed[P any, T CanView[V], V any](l FixedList[P, T]) ListViewEmbed[P, T, V] {
	return ListViewEmbed[P, T, V]{l}
}

