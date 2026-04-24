package example_test

import (
	"testing"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/order"
	"github.com/krelinga/go-lego/order/orderexam"
	"github.com/krelinga/go-lego/pod"
)

type FooStructView interface {
	GetFoo() int
	GetBar() string
}

type FooStruct struct {
	Foo int
	Bar string
}

func (f FooStruct) GetFoo() int {
	return f.Foo
}

func (f FooStruct) GetBar() string {
	return f.Bar
}

func FooOrder(a, b FooStructView) int {
	return order.Using(a, b,
		order.Get(FooStructView.GetBar),
		order.Get(FooStructView.GetFoo),
	)
}

func FooEqual(a, b FooStructView) bool {
	return a.GetFoo() == b.GetFoo()
}

func TestExample(t *testing.T) {
	t.Run("simple int comparison", func(t *testing.T) {
		one := 1
		zero := 0
		exam.Must(t, exam.Equal(one, 1))
		exam.Try(t, exam.Greater(one, zero))
	})

	t.Run("struct comparison", func(t *testing.T) {
		a := FooStruct{Foo: 1, Bar: "a"}
		b := FooStruct{Foo: 1, Bar: "b"}
		exam.Must(t, orderexam.LessFunc(a, b, FooOrder))
		exam.Must(t, exam.EqualFunc(a, b, FooEqual))
	})

	t.Run("Map of struct comparison", func(t *testing.T) {
		a := &pod.Map[string, FooStructView]{
			"a": FooStruct{Foo: 1, Bar: "a"},
			"b": FooStruct{Foo: 2, Bar: "b"},
		}
		b := &pod.Map[string, FooStructView]{
			"a": FooStruct{Foo: 1, Bar: "a"},
			"b": FooStruct{Foo: 2, Bar: "b"},
		}
		// TODO: add a pod helper for this pattern of comparing two maps with a custom equality function.
		exam.Must(t, exam.True(pod.MapEqualFunc(a, b, FooEqual)), a, b)
	})

	t.Run("nil comparison", func(t *testing.T) {
		var a *int
		var b []int
		var c map[string]int
		exam.Try(t, exam.Nil(a))
		exam.Try(t, exam.Nil(b))
		exam.Try(t, exam.Nil(c))
	})

	cases := []struct {
		Name    string
		Loc     exam.Loc
		A, B, C int
	}{
		{Name: "case 1", Loc: exam.Here(), A: 1, B: 2, C: 3},
		{
			Name: "case 2",
			Loc:  exam.Here(),
			A:    4,
			B:    5,
			C:    9,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			exam.Try(t, exam.Equal(c.A+c.B, c.C+1))
		})
	}
}
