package match

import (
	"cmp"
	"reflect"

	"github.com/krelinga/go-lego/pod"
)

func LessThan[T cmp.Ordered](expected T) CanMatch {
	return lessThan[T]{expected: expected}
}

type lessThan[T cmp.Ordered] struct {
	expected T
}

func (m lessThan[T]) Match(actual any, r Reporter) bool {
	actualValue, ok := actual.(T)
	if !ok {
		if r != nil {
			r.TypeMismatch(reflect.TypeOf(m.expected), reflect.TypeOf(actual))
		}
		return false
	}
	if !(actualValue < m.expected) {
		if r != nil {
			r.ValueMismatch(m.expected, actualValue)
		}
		return false
	}
	return true
}

func LessThanPod[T pod.CanCompare[T]](expected T) CanMatch {
	return lessThanPod[T]{expected: expected}
}

type lessThanPod[T pod.CanCompare[T]] struct {
	expected T
}

func (m lessThanPod[T]) Match(actual any, r Reporter) bool {
	actualValue, ok := actual.(T)
	if !ok {
		if r != nil {
			r.TypeMismatch(reflect.TypeOf(m.expected), reflect.TypeOf(actual))
		}
		return false
	}
	if !(pod.LessThan(m.expected, actualValue)) {
		if r != nil {
			r.ValueMismatch(m.expected, actualValue)
		}
		return false
	}
	return true
}