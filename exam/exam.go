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

func Pred2[T1, T2 any](f func(T1, T2) bool) func(E, T1, T2) Result {
	return func(e E, got T1, want T2) Result {
		return nil // TODO
	}
}