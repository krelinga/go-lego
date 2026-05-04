package order_test

import (
	"cmp"
	"testing"

	"github.com/krelinga/go-libs/exam"
	"github.com/krelinga/go-libs/order"
)

func TestDesc(t *testing.T) {
	cases := []struct {
		Name string
		Loc  exam.Loc
		A, B int
		Want func(int) bool
	}{
		{
			Name: "a < b",
			Loc:  exam.Here(),
			A:    1,
			B:    2,
			Want: order.Greater,
		},
		{
			Name: "a > b",
			Loc:  exam.Here(),
			A:    2,
			B:    1,
			Want: order.Less,
		}, {
			Name: "a == b",
			Loc:  exam.Here(),
			A:    1,
			B:    1,
			Want: order.Equal,
		},
	}

	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			exam.Try(t, exam.True(c.Want(order.Desc(cmp.Compare[int])(c.A, c.B))))
		})
	}
}
