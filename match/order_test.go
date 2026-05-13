package match_test

import (
	"testing"

	m "github.com/krelinga/go-libs/match"
)

func TestOrder(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		tryMatch(t, 5, &m.Order{
			Op:    m.OrderOpGt(),
			Limit: 3,
		})
		tryMatch(t, 5, m.OrderGt(3))
	})
}
