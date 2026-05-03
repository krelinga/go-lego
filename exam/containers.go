package exam

import "github.com/krelinga/go-lego/mirror"

func MapEqual[K, V comparable](a, b map[K]V) *Failure {
	return MapEqualFunc(a, b, func(x, y V) bool { return x == y })
}

func MapEqualFunc[K comparable, V any](a, b map[K]V, eq any) *Failure {
	failure := NewFailure2("a", a, "b", b)
	eqFunc, err := mirror.WrapFunc2In1Out[V, V, bool](eq)
	if err != nil {
		return failure.Wrap(err)
	}
	if len(a) != len(b) {
		return failure
	}
	for k, vA := range a {
		vB, ok := b[k]
		if !ok || !eqFunc(vA, vB) {
			return failure
		}
	}
	return nil
}

func SliceEqual[V comparable](a, b []V) *Failure {
	return SliceEqualFunc(a, b, func(x, y V) bool { return x == y })
}

func SliceEqualFunc[V any](a, b []V, eq any) *Failure {
	failure := NewFailure2("a", a, "b", b)
	eqFunc, err := mirror.WrapFunc2In1Out[V, V, bool](eq)
	if err != nil {
		return failure.Wrap(err)
	}
	if len(a) != len(b) {
		return failure
	}
	for i := range a {
		if !eqFunc(a[i], b[i]) {
			return failure
		}
	}
	return nil
}
