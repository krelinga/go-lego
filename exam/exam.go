package exam

import (
	"cmp"
	"errors"
	"fmt"
	"reflect"
	"strings"
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

type ErrorParam struct {
	Name  string
	Value any
}

func Param(name string, value any) ErrorParam {
	return ErrorParam{Name: name, Value: value}
}

type Params = []ErrorParam

type Error struct {
	Name   string
	Params []ErrorParam
}

func (e Error) Error() string {
	paramsFmt := make([]string, len(e.Params))
	for i, p := range e.Params {
		paramsFmt[i] = fmt.Sprintf("  %s: %v", p.Name, p.Value)
	}
	return fmt.Sprintf("exam: %s failed:\n%s", e.Name, strings.Join(paramsFmt, "\n"))
}

func Equal[T comparable](e E, actual, expected T, opts ...Option) bool {
	return NewPred2(func(got, want T) error {
		if got == want {
			return nil
		}
		return Error{
			Name:   "Equal",
			Params: Params{
				Param("actual", got),
				Param("expected", want),
			},
		}
	})(e, actual, expected, opts...)
}

func GreaterThan[T cmp.Ordered](e E, value, threshold T, opts ...Option) bool {
	return NewPred2(func(got, want T) error {
		if got > want {
			return nil
		}
		return Error{
			Name:   "GreaterThan",
			Params: Params{
				Param("value", got),
				Param("threshold", want),
			},
		}
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
		return Error{
			Name:   "Nil",
			Params: Params{Param("value", got)},
		}
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
		return Error{
			Name:   "NotNil",
			Params: Params{Param("value", got)},
		}
	})(e, got, opts...)
}

func PodMapIsSubset[K, V comparable](e E, subset, superset pod.MapView[K, V], opts ...Option) bool {
	return NewPred2(func(sub, super pod.MapView[K, V]) error {

		for k, v := range sub.KeyVals() {
			superVal, ok := super.Get(k)
			if !ok || superVal != v {
				return Error{
					Name:   "PodMapIsSubset",
					Params: Params{
						Param("subset", sub),
						Param("superset", super),
					},
				}
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
				return Error{
					Name:   "PodMapIsSubsetFunc",
					Params: Params{
						Param("subset", sub),
						Param("superset", super),
					},
				}
			}
		}
		return nil
	})(e, subset, superset, opts...)
}
