package exam

import (
	"cmp"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/krelinga/go-lego/pod"
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

var ErrCritical = errors.New("exam: critical failure")

type errorBuilderParam struct {
	name  string
	value any
}

type ErrorBuilder struct {
	op     string
	params []errorBuilderParam
}

func (e *ErrorBuilder) AddParam(name string, value any) *ErrorBuilder {
	e.params = append(e.params, errorBuilderParam{name: name, value: value})
	return e
}

func (e *ErrorBuilder) Error() error {
	return fmt.Errorf("exam: %s failed: %v", e.op, e.params)
}

func NewErrorBuilder(op string) *ErrorBuilder {
	return &ErrorBuilder{op: op}
}

func Equal[T comparable](e E, actual, expected T, opts ...Option) bool {
	return NewPred2(func(got, want T) error {
		if got == want {
			return nil
		}
		return NewErrorBuilder("Equal").
			AddParam("actual", got).
			AddParam("expected", want).
			Error()
	})(e, actual, expected, opts...)
}

func GreaterThan[T cmp.Ordered](e E, value, threshold T, opts ...Option) bool {
	return NewPred2(func(got, want T) error {
		if got > want {
			return nil
		}
		return NewErrorBuilder("GreaterThan").
			AddParam("value", got).
			AddParam("threshold", want).
			Error()
	})(e, value, threshold, opts...)
}

type Pred2[T any] func(E, T, T, ...Option) bool

func NewPred2[T any](f func(T, T) error) Pred2[T] {
	return func(e E, got T, want T, opts ...Option) bool {
		return false // TODO
	}
}

type Pred[T any] func(E, T, ...Option) bool

func NewPred[T any](f func(T) error) Pred[T] {
	return func(e E, got T, opts ...Option) bool {
		return false // TODO
	}
}

type options struct{}

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
	return NewPred(func(got any) error {
		isNil, err := isNil(got)
		if err != nil {
			return err
		} else if isNil {
			return nil
		}
		return NewErrorBuilder("Nil").
			AddParam("value", got).
			Error()
	})(e, got, opts...)
}

func NotNil(e E, got any, opts ...Option) bool {
	return NewPred(func(got any) error {
		isNil, err := isNil(got)
		if err != nil {
			return err
		} else if !isNil {
			return nil
		}
		return NewErrorBuilder("NotNil").
			AddParam("value", got).
			Error()
	})(e, got, opts...)
}

func PodMapIsSubset[K, V comparable](e E, subset, superset pod.MapView[K, V], opts ...Option) bool {
	return NewPred2(func(sub, super pod.MapView[K, V]) error {

		for k, v := range sub.KeyVals() {
			superVal, ok := super.Get(k)
			if !ok || superVal != v {
				return NewErrorBuilder("PodMapIsSubset").
					AddParam("subset", sub).
					AddParam("superset", super).
					Error()
			}
		}
		return nil
	})(e, subset, superset, opts...)
}

func PodMapIsSubsetFunc[K, V any](e E, subset, superset pod.MapView[K, V], equal func(V, V) bool, opts ...Option) bool {
	return NewPred2(func(sub, super pod.MapView[K, V]) error {
		for k, v := range sub.KeyVals() {
			superVal, ok := super.Get(k)
			if !ok || !equal(superVal, v) {
				return NewErrorBuilder("PodMapIsSubsetFunc").
					AddParam("subset", sub).
					AddParam("superset", super).
					Error()
			}
		}
		return nil
	})(e, subset, superset, opts...)
}
