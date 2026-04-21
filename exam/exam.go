package exam

import (
	"cmp"
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

func Equal[T comparable](e E, got, want T, opts ...Option) bool {
	return false // TODO
}

func GreaterThan[T cmp.Ordered](e E, x, t T, opts ...Option) bool {
	return false // TODO
}

type Pred2[T any] func(E, T, T, ...Option) bool

func NewPred2[T any](f func(T, T) bool) Pred2[T] {
	return func(e E, got T, want T, opts ...Option) bool {
		return false // TODO
	}
}

type Pred[T any] func(E, T, ...Option) bool

func NewPred[T any](f func(T) bool) Pred[T] {
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