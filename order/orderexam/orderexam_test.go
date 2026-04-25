package orderexam_test

import (
	"cmp"
	"testing"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/order/orderexam"
)

type Struct struct {
	Int int
}

func StructOrder(a, b Struct) int {
	return cmp.Compare(a.Int, b.Int)
}

func Test(t *testing.T) {
	cases := []struct {
		Name string
		Loc exam.Loc
		A, B Struct
		Func func(Struct, Struct, any) *exam.Failure2
		WantPass bool
	}{
		{
			Name: "greater pass",
			Loc: exam.Here(),
			A: Struct{Int: 2},
			B: Struct{Int: 1},
			Func: orderexam.GreaterFunc[Struct],
			WantPass: true,
		},
		{
			Name: "greater fail",
			Loc: exam.Here(),
			A: Struct{Int: 1},
			B: Struct{Int: 2},
			Func: orderexam.GreaterFunc[Struct],
			WantPass: false,
		},
		{
			Name: "equal pass",
			Loc: exam.Here(),
			A: Struct{Int: 1},
			B: Struct{Int: 1},
			Func: orderexam.EqualFunc[Struct],
			WantPass: true,
		},
		{
			Name: "equal fail",
			Loc: exam.Here(),
			A: Struct{Int: 1},
			B: Struct{Int: 2},
			Func: orderexam.EqualFunc[Struct],
			WantPass: false,
		},
		{
			Name: "less pass",
			Loc: exam.Here(),
			A: Struct{Int: 1},
			B: Struct{Int: 2},
			Func: orderexam.LessFunc[Struct],
			WantPass: true,
		},
		{
			Name: "less fail",
			Loc: exam.Here(),
			A: Struct{Int: 2},
			B: Struct{Int: 1},
			Func: orderexam.LessFunc[Struct],
			WantPass: false,
		},
		{
			Name: "less equal pass when less",
			Loc: exam.Here(),
			A: Struct{Int: 1},
			B: Struct{Int: 2},
			Func: orderexam.LessEqualFunc[Struct],
			WantPass: true,
		},
		{
			Name: "less equal pass when equal",
			Loc: exam.Here(),
			A: Struct{Int: 1},
			B: Struct{Int: 1},
			Func: orderexam.LessEqualFunc[Struct],
			WantPass: true,
		},
		{
			Name: "less equal fail",
			Loc: exam.Here(),
			A: Struct{Int: 2},
			B: Struct{Int: 1},
			Func: orderexam.LessEqualFunc[Struct],
			WantPass: false,
		},
		{
			Name: "greater equal pass when greater",
			Loc: exam.Here(),
			A: Struct{Int: 2},
			B: Struct{Int: 1},
			Func: orderexam.GreaterEqualFunc[Struct],
			WantPass: true,
		},
		{
			Name: "greater equal pass when equal",
			Loc: exam.Here(),
			A: Struct{Int: 1},
			B: Struct{Int: 1},
			Func: orderexam.GreaterEqualFunc[Struct],
			WantPass: true,
		},
		{
			Name: "greater equal fail",
			Loc: exam.Here(),
			A: Struct{Int: 1},
			B: Struct{Int: 2},
			Func: orderexam.GreaterEqualFunc[Struct],
			WantPass: false,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			result := c.Func(c.A, c.B, StructOrder)
			if c.WantPass {
				exam.Must(t, exam.Nil(result))
			} else {
				exam.Must(t, exam.NotNil(result))
			}
		})
	}
}