package v2

type Equaler[V any] interface {
	Equal(V) bool
}

func DictEqualValues[K comparable, V Equaler[V]](a, b FixedDict[K, V]) bool {
	return dictEqualHelper(a, b, func(vA, vB V) bool {
		return vA.Equal(vB)
	})
}

func DictEqualValuesComparable[K comparable, V comparable](a, b FixedDict[K, V]) bool {
	return dictEqualHelper(a, b, func(vA, vB V) bool {
		return vA == vB
	})
}

func DictEqualValuesComparer[K comparable, V Comparer[V]](a, b FixedDict[K, V]) bool {
	return dictEqualHelper(a, b, func(vA, vB V) bool {
		return vA.Compare(vB) == 0
	})
}

func ListEqualValues[P any, V Equaler[V]](a, b FixedList[P, V]) bool {
	return listEqualHelper(a, b, func(vA, vB V) bool {
		return vA.Equal(vB)
	})
}

func ListEqualValuesComparable[P any, V comparable](a, b FixedList[P, V]) bool {
	return listEqualHelper(a, b, func(vA, vB V) bool {
		return vA == vB
	})
}

func ListEqualValuesComparer[P any, V Comparer[V]](a, b FixedList[P, V]) bool {
	return listEqualHelper(a, b, func(vA, vB V) bool {
		return vA.Compare(vB) == 0
	})
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
