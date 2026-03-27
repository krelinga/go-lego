package lego

import (
	"cmp"
	"slices"
)

type CmpBooklet[V any] interface {
	Cmp(V, V) int
}

type OrderedCmpBooklet[V cmp.Ordered] struct{}

func (_ OrderedCmpBooklet[V]) Cmp(a, b V) int {
	return cmp.Compare(a, b)
}

type WithLegoBooklet[V any] interface {
	LegoBooklet() V
}

func Sort[L Lister[V], V WithLegoBooklet[B], B CmpBooklet[V]](in L) []V {
	var out []V
	if l, ok := any(in).(Lener); ok && l.Len() > 0 {
		out = make([]V, l.Len())
	}
	for v := range in.List() {
		out = append(out, v)
	}
	var booklet B
	slices.SortFunc(out, booklet.Cmp)
	return out
}
