package match_test

import (
	"testing"

	m "github.com/krelinga/go-libs/match"
)

func TestNil(t *testing.T) {
	tryMatch(t, nil, m.Nil())
	tryMatch(t, &struct{}{}, m.Not(m.Nil()))
}
