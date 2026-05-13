package match_test

import (
	"testing"

	m "github.com/krelinga/go-libs/match"
)

func TestOrder(t *testing.T) {
	tryMatch := func(t *testing.T, val any, m m.Matcher) {
		t.Helper()
		if err := m.Validate(); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
		r, err := m.Match(val)
		if err != nil {
			t.Fatalf("unexpected match error: %v", err)
		}
		if !r.Accepted {
			t.Error("expected match to succeed, but it failed")
		}
	}
	t.Run("basic", func(t *testing.T) {
		tryMatch(t, 5, &m.Order{
			Op:    m.OrderOpGt(),
			Limit: 3,
		})
		tryMatch(t, 5, m.OrderGt(3))
	})
}
