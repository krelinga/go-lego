package match

import (
	"reflect"

	"github.com/krelinga/go-lego/pod"
)

func Equal[T comparable](expected T) CanMatch {
	return equal[T]{expected: expected}
}

type equal[T comparable] struct {
	expected T
}

func (m equal[T]) Match(actual any, r Reporter) bool {
	actualValue, ok := actual.(T)
	if !ok {
		if r != nil {
			r.TypeMismatch(reflect.TypeOf(m.expected), reflect.TypeOf(actual))
		}
		return false
	}
	if m.expected != actualValue {
		if r != nil {
			r.ValueMismatch(m.expected, actualValue)
		}
		return false
	}
	return true
}

func EqualPod[T pod.CanEqual[T]](expected T) CanMatch {
	return equalPod[T]{expected: expected}
}

type equalPod[T pod.CanEqual[T]] struct {
	expected T
}

func (m equalPod[T]) Match(actual any, r Reporter) bool {
	actualValue, ok := actual.(T)
	if !ok {
		if r != nil {
			r.TypeMismatch(reflect.TypeOf(m.expected), reflect.TypeOf(actual))
		}
		return false
	}
	if !m.expected.Equal(actualValue) {
		if r != nil {
			r.ValueMismatch(m.expected, actualValue)
		}
		return false
	}
	return true
}
