package exam_test

import (
	"cmp"
	"testing"

	"github.com/krelinga/go-libs/exam"
)

type orderStruct struct {
	Int int
}

func orderStructOrder(a, b orderStruct) int {
	return cmp.Compare(a.Int, b.Int)
}

func TestOrder(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		A, B     orderStruct
		Func     func(orderStruct, orderStruct, any) *exam.Failure
		WantPass bool
	}{
		{
			Name:     "greater pass",
			Loc:      exam.Here(),
			A:        orderStruct{Int: 2},
			B:        orderStruct{Int: 1},
			Func:     exam.GreaterFunc[orderStruct],
			WantPass: true,
		},
		{
			Name:     "greater fail",
			Loc:      exam.Here(),
			A:        orderStruct{Int: 1},
			B:        orderStruct{Int: 2},
			Func:     exam.GreaterFunc[orderStruct],
			WantPass: false,
		},
		{
			Name:     "equal pass",
			Loc:      exam.Here(),
			A:        orderStruct{Int: 1},
			B:        orderStruct{Int: 1},
			Func:     exam.OrderEqualFunc[orderStruct],
			WantPass: true,
		},
		{
			Name:     "equal fail",
			Loc:      exam.Here(),
			A:        orderStruct{Int: 1},
			B:        orderStruct{Int: 2},
			Func:     exam.OrderEqualFunc[orderStruct],
			WantPass: false,
		},
		{
			Name:     "less pass",
			Loc:      exam.Here(),
			A:        orderStruct{Int: 1},
			B:        orderStruct{Int: 2},
			Func:     exam.LessFunc[orderStruct],
			WantPass: true,
		},
		{
			Name:     "less fail",
			Loc:      exam.Here(),
			A:        orderStruct{Int: 2},
			B:        orderStruct{Int: 1},
			Func:     exam.LessFunc[orderStruct],
			WantPass: false,
		},
		{
			Name:     "less equal pass when less",
			Loc:      exam.Here(),
			A:        orderStruct{Int: 1},
			B:        orderStruct{Int: 2},
			Func:     exam.LessEqualFunc[orderStruct],
			WantPass: true,
		},
		{
			Name:     "less equal pass when equal",
			Loc:      exam.Here(),
			A:        orderStruct{Int: 1},
			B:        orderStruct{Int: 1},
			Func:     exam.LessEqualFunc[orderStruct],
			WantPass: true,
		},
		{
			Name:     "less equal fail",
			Loc:      exam.Here(),
			A:        orderStruct{Int: 2},
			B:        orderStruct{Int: 1},
			Func:     exam.LessEqualFunc[orderStruct],
			WantPass: false,
		},
		{
			Name:     "greater equal pass when greater",
			Loc:      exam.Here(),
			A:        orderStruct{Int: 2},
			B:        orderStruct{Int: 1},
			Func:     exam.GreaterEqualFunc[orderStruct],
			WantPass: true,
		},
		{
			Name:     "greater equal pass when equal",
			Loc:      exam.Here(),
			A:        orderStruct{Int: 1},
			B:        orderStruct{Int: 1},
			Func:     exam.GreaterEqualFunc[orderStruct],
			WantPass: true,
		},
		{
			Name:     "greater equal fail",
			Loc:      exam.Here(),
			A:        orderStruct{Int: 1},
			B:        orderStruct{Int: 2},
			Func:     exam.GreaterEqualFunc[orderStruct],
			WantPass: false,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			result := c.Func(c.A, c.B, orderStructOrder)
			if c.WantPass {
				exam.Must(t, exam.Nil(result))
			} else {
				exam.Must(t, exam.NotNil(result))
			}
		})
	}
}
