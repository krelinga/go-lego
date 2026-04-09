package v2

import (
	"fmt"
	"iter"
	"maps"
	"strings"
)

type MapSetValue struct{}

type MapSet[V comparable] map[V]MapSetValue

func (s *MapSet[V]) Length() int {
	return len(*s)
}

func (s *MapSet[V]) Has(v V) bool {
	_, ok := (*s)[v]
	return ok
}

func (s *MapSet[T]) Values() iter.Seq[T] {
	return maps.Keys(*s)
}

func (s *MapSet[V]) String() string {
	sb := strings.Builder{}
	sb.WriteString("{")
	first := true
	for v := range *s {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		fmt.Fprintf(&sb, "%v", v)
	}
	sb.WriteString("}")
	return sb.String()
}

func (s *MapSet[V]) Add(v V) {
	(*s)[v] = MapSetValue{}
}

func (s *MapSet[V]) Remove(v V) {
	delete(*s, v)
}

func (s *MapSet[V]) Clear() {
	*s = nil
}
