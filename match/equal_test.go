package match_test

import (
	"fmt"
	"strings"
	"testing"

	m "github.com/krelinga/go-libs/match"
)

type StringStringer string

func (s StringStringer) String() string {
	return string(s)
}

func tryMatch(t *testing.T, val any, m m.Matcher) {
	t.Helper()
	r, err := m.Match(val)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Accepted {
		t.Fatalf("expected match to be accepted but it was rejected: %s", r.Why)
	}
}

func TestEqual(t *testing.T) {
	t.Run("basic equality", func(t *testing.T) {
		m := m.Equal(42)
		r, err := m.Match(42)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !r.Accepted {
			t.Fatalf("expected match to be accepted but it was rejected: %s", r.Why)
		}
	})

	t.Run("equal function", func(t *testing.T) {
		caselessEq := func(a, b fmt.Stringer) bool {
			aStr, bStr := a.String(), b.String()
			return strings.EqualFold(aStr, bStr)
		}
		m := m.EqualFunc(StringStringer("Hello, World!"), caselessEq)
		r, err := m.Match(StringStringer("hello, world!"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !r.Accepted {
			t.Fatalf("expected match to be accepted but it was rejected: %s", r.Why)
		}
	})

	t.Run("approx", func(t *testing.T) {
		m := m.EqualApprox(3.14, 0.01)
		r, err := m.Match(3.1415)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !r.Accepted {
			t.Fatalf("expected match to be accepted but it was rejected: %s", r.Why)
		}
	})

	t.Run("with tryMatch", func(t *testing.T) {
		tryMatch(t, 42, m.Equal(42))
		tryMatch(t, "hello", m.Equal("hello"))
		tryMatch(t, "Hello World", m.EqualFunc("hello world", strings.EqualFold))
		tryMatch(t, 3.14159, m.EqualApprox(3.14, 0.01))
	})
}
