package exam

import (
	"testing"

	"github.com/krelinga/go-lego/match"
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

func Match(e E, got any, matcher match.Matcher) Result {
	return nil  // TODO
}