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
		m := m.Equal{
			Want: 42,
		}
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
		m := m.Equal{
			Func: m.NewEqualFunc(caselessEq),
			Want: StringStringer("Hello, World!"),
		}
		r, err := m.Match(StringStringer("hello, world!"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !r.Accepted {
			t.Fatalf("expected match to be accepted but it was rejected: %s", r.Why)
		}
	})

	t.Run("EqualCmp", func(t *testing.T) {
		r, err := m.EqualCmp(3.14).Match(3.14)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !r.Accepted {
			t.Fatalf("expected match to be accepted but it was rejected: %s", r.Why)
		}
	})

	t.Run("with tryMatch", func(t *testing.T) {
		tryMatch(t, 42, m.EqualCmp(42))
		tryMatch(t, "hello", &m.Equal{Want: "hello"})
		tryMatch(t, "Hello World", &m.Equal{
			Func: m.NewEqualFunc(strings.EqualFold),
			Want: "hello world",
		})
	})
}
