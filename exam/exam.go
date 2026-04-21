package exam

import (
	"cmp"
	"fmt"
	"reflect"
	"testing"
)

type E interface {
	// methods proxied from *testing.T.
	Helper()

	// TODO: test failure reporting methods, etc.
	Run(string, func(E)) bool
}

func New(t *testing.T) E {
	return nil // TODO
}

func Equal[T comparable](e E, actual, expected T, opts ...Option) bool {
	return NewPred2("Equal", "actual", "expected", func(got, want T) bool {
		return got == want
	})(e, actual, expected, opts...)
}

func GreaterThan[T cmp.Ordered](e E, value, threshold T, opts ...Option) bool {
	return NewPred2("GreaterThan", "value", "threshold", func(got, want T) bool {
		return got > want
	})(e, value, threshold, opts...)
}

type Pred2[T any] func(E, T, T, ...Option) bool

func NewPred2[T any](op, p1, p2 string, f func(T, T) bool) Pred2[T] {
	return func(e E, got T, want T, opts ...Option) bool {
		return false // TODO
	}
}

func NewPred2Err[T any](op, p1, p2 string, f func(T, T) (bool, error)) Pred2[T] {
	return func(e E, got T, want T, opts ...Option) bool {
		return false // TODO
	}
}

type Pred[T any] func(E, T, ...Option) bool

func NewPred[T any](op, p string, f func(T) bool) Pred[T] {
	return func(e E, got T, opts ...Option) bool {
		return false // TODO
	}
}

func NewPredErr[T any](op, p string, f func(T) (bool, error)) Pred[T] {
	return func(e E, got T, opts ...Option) bool {
		return false // TODO
	}
}

type options struct {}

type Option func(*options)

func Must() Option {
	return nil // TODO
}

func Log(items ...any) Option {
	return nil // TODO
}

func Logf(format string, args ...any) Option {
	return nil // TODO
}

func Quiet() Option {
	return nil // TODO
}

func isNil(v any) (bool, error) {
	if v == nil {
		return true, nil
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return val.IsNil(), nil
	default:
		return false, fmt.Errorf("value of type %T cannot be nil", v)
	}
}

func Nil(e E, got any, opts ...Option) bool {
	return NewPredErr("Nil", "value", isNil)(e, got, opts...)
}

func NotNil(e E, got any, opts ...Option) bool {
	return NewPredErr("NotNil", "value", func(v any) (bool, error) {
		isNil, err := isNil(v)
		if err != nil {
			return false, err
		}
		return !isNil, nil
	})(e, got, opts...)
}