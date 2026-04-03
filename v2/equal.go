package v2

type Equaler[V any] interface {
	Equal(V) bool
}

type DictEqualEmbed[O FixedDict[K, T], K comparable, T Equaler[T]] struct {
	fd FixedDict[K, T]
}

func (d DictEqualEmbed[O, K, T]) Equal(other O) bool {
	return dictEqualHelper(d.fd, other, func(a, b T) bool {
		return a.Equal(b)
	})
}

func NewDictEqualEmbed[O FixedDict[K, T], K comparable, T Equaler[T]](d FixedDict[K, T]) DictEqualEmbed[O, K, T] {
	return DictEqualEmbed[O, K, T]{d}
}

type DictEqualEmbedComparable[O FixedDict[K, T], K comparable, T comparable] struct {
	fd FixedDict[K, T]
}

func (m DictEqualEmbedComparable[O, K, T]) Equal(other O) bool {
	return dictEqualHelper(m.fd, other, func(a, b T) bool {
		return a == b
	})
}

func NewDictEqualEmbedComparable[O FixedDict[K, T], K comparable, T comparable](d FixedDict[K, T]) DictEqualEmbedComparable[O, K, T] {
	return DictEqualEmbedComparable[O, K, T]{d}
}

type DictEqualEmbedComparer[O FixedDict[K, T], K comparable, T Comparer[T]] struct {
	fd FixedDict[K, T]
}

func (m DictEqualEmbedComparer[O, K, T]) Equal(other O) bool {
	return dictEqualHelper(m.fd, other, func(a, b T) bool {
		return a.Compare(b) == 0
	})
}

func NewDictEqualEmbedComparer[O FixedDict[K, T], K comparable, T Comparer[T]](d FixedDict[K, T]) DictEqualEmbedComparer[O, K, T] {
	return DictEqualEmbedComparer[O, K, T]{d}
}

func dictEqualHelper[K comparable, V any](a, b FixedDict[K, V], equalFunc func(V, V) bool) bool {
	if a == nil && b == nil {
		return true
	}
	if (a == nil) != (b == nil) {
		return false
	}
	if a.Length() != b.Length() {
		return false
	}
	for k, v := range a.All() {
		otherV, ok := b.Get(k)
		if !ok || !equalFunc(v, otherV) {
			return false
		}
	}
	return true
}

type ListEqualEmbed[O FixedList[P, T], P any, T Equaler[T]] struct {
	fl FixedList[P, T]
}

func (l ListEqualEmbed[O, P, T]) Equal(other O) bool {
	return listEqualHelper(l.fl, other, func(a, b T) bool {
		return a.Equal(b)
	})
}

func NewListEqualEmbed[O FixedList[P, T], P any, T Equaler[T]](l FixedList[P, T]) ListEqualEmbed[O, P, T] {
	return ListEqualEmbed[O, P, T]{l}
}

type ListEqualEmbedComparable[O FixedList[P, T], P any, T comparable] struct {
	fl FixedList[P, T]
}

func (l ListEqualEmbedComparable[O, P, T]) Equal(other O) bool {
	return listEqualHelper(l.fl, other, func(a, b T) bool {
		return a == b
	})
}

func NewListEqualEmbedComparable[O FixedList[P, T], P any, T comparable](l FixedList[P, T]) ListEqualEmbedComparable[O, P, T] {
	return ListEqualEmbedComparable[O, P, T]{l}
}

type ListEqualEmbedComparer[O FixedList[P, T], P any, T Comparer[T]] struct {
	fl FixedList[P, T]
}

func (l ListEqualEmbedComparer[O, P, T]) Equal(other O) bool {
	return listEqualHelper(l.fl, other, func(a, b T) bool {
		return a.Compare(b) == 0
	})
}

func NewListEqualEmbedComparer[O FixedList[P, T], P any, T Comparer[T]](l FixedList[P, T]) ListEqualEmbedComparer[O, P, T] {
	return ListEqualEmbedComparer[O, P, T]{l}
}

func listEqualHelper[P any, V any](a, b FixedList[P, V], equalFunc func(V, V) bool) bool {
	if a == nil && b == nil {
		return true
	}
	if (a == nil) != (b == nil) {
		return false
	}
	if a.Length() != b.Length() {
		return false
	}
	for i, vA := range a.All() {
		vb, ok := b.Get(i)
		if !ok || !equalFunc(vA, vb) {
			return false
		}
	}
	return true
}
