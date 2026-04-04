package v2

import (
	"fmt"
	"iter"
	"strings"
)

type FixedDict[K, V any] interface {
	Length() int
	Get(K) (V, bool)
	All() iter.Seq2[K, V]
	Keys() iter.Seq[K]
	Values() iter.Seq[V]
	KVs() iter.Seq[KV[K, V]]
	String() string
}

type Dict[K, V any] interface {
	FixedDict[K, V]
	Set(K, V)
	Remove(K)
	Clear()
}

type KV[K, V any] struct {
	K K
	V V
}

func NewKV[K, V any](k K, v V) KV[K, V] {
	return KV[K, V]{k, v}
}

func dictStringHelper[K comparable, V any](d FixedDict[K, V]) string {
	var sb strings.Builder
	sb.WriteString("{")
	first := true
	for k, v := range d.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		fmt.Fprintf(&sb, "%v: %v", k, v)
	}
	sb.WriteString("}")
	return sb.String()
}