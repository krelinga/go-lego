package exam

import (
	"cmp"
	"testing"
)

type E interface {
	// methods proxied from *testing.T.
	Helper()

	// TODO: test failure reporting methods, etc.
	Run(string, func(E)) Result
}

// TODO: this should probably be a struct instead of an interface, but this is easier for prototyping.
type Result interface {
	OK() bool
	Must()
}

func New(t *testing.T) E {
	return nil  // TODO
}

func Equal(e E, got, want any) Result {
	return nil  // TODO
}

func GreaterThan[T cmp.Ordered](e E, x, t T) Result {
	return nil  // TODO
}

func Pred2[T any](f func(T, T) bool) func(E, T, T) Result {
	return func(e E, got T, want T) Result {
		return nil // TODO
	}
}