package order_test

import (
	"cmp"
	"testing"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/order"
	"github.com/krelinga/go-lego/pod"
)

type Struct struct {
	Foo int
	Bar string
}

func (s Struct) GetFoo() int {
	return s.Foo
}

func (s Struct) GetBar() string {
	return s.Bar
}

func StructOrder(a, b Struct) int {
	return order.Using(a, b,
		order.Get(Struct.GetFoo),
		order.Get(Struct.GetBar),
	)
}

func TestUsingGet(t *testing.T) {
	cases := []struct {
		Name string
		Loc exam.Loc
		A, B Struct
		Want func(int) bool
	}{
		{
			Name: "a < b",
			Loc: exam.Here(),
			A: Struct{Foo: 1, Bar: "a"},
			B: Struct{Foo: 2, Bar: "b"},
			Want: order.Less,
		},
		{
			Name: "a > b",
			Loc: exam.Here(),
			A: Struct{Foo: 2, Bar: "b"},
			B: Struct{Foo: 1, Bar: "a"},
			Want: order.Greater,
		},
		{
			Name: "a == b",
			Loc: exam.Here(),
			A: Struct{Foo: 1, Bar: "a"},
			B: Struct{Foo: 1, Bar: "a"},
			Want: order.Equal,
		},
		{
			Name: "a < b by Bar",
			Loc: exam.Here(),
			A: Struct{Foo: 1, Bar: "a"},
			B: Struct{Foo: 1, Bar: "b"},
			Want: order.Less,
		},
		{
			Name: "a > b by Bar",
			Loc: exam.Here(),
			A: Struct{Foo: 1, Bar: "b"},
			B: Struct{Foo: 1, Bar: "a"},
			Want: order.Greater,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			exam.Must(t, c.Want(StructOrder(c.A, c.B)), c.A, c.B)
		})
	}
}

type Struct2 struct {
	Foo int
	Bar *pod.Vec[string]
}

func (s Struct2) GetFoo() int {
	return s.Foo
}

func (s Struct2) GetBar() pod.VecView[string] {
	return s.Bar
}

func VecOrder(a, b pod.VecView[string]) int {
	for i := 0; i < a.Len() && i < b.Len(); i++ {
		if idxRes := cmp.Compare(a.At(i), b.At(i)); idxRes != 0 {
			return idxRes
		}
	}
	return cmp.Compare(a.Len(), b.Len())
}

func Struct2Order(a, b Struct2) int {
	return order.Using(a, b,
		order.Get(Struct2.GetFoo),
		order.GetFunc(Struct2.GetBar, VecOrder),
	)
}

func TestUsingGetFunc(t *testing.T) {
	cases := []struct {
		Name string
		Loc exam.Loc
		A, B Struct2
		Want func(int) bool
	}{
		{
			Name: "a < b",
			Loc: exam.Here(),
			A: Struct2{Foo: 1, Bar: &pod.Vec[string]{"a"}},
			B: Struct2{Foo: 2, Bar: &pod.Vec[string]{"b"}},
			Want: order.Less,
		},
		{
			Name: "a > b",
			Loc: exam.Here(),
			A: Struct2{Foo: 2, Bar: &pod.Vec[string]{"b"}},
			B: Struct2{Foo: 1, Bar: &pod.Vec[string]{"a"}},
			Want: order.Greater,
		},
		{
			Name: "a == b",
			Loc: exam.Here(),
			A: Struct2{Foo: 1, Bar: &pod.Vec[string]{"a"}},
			B: Struct2{Foo: 1, Bar: &pod.Vec[string]{"a"}},
			Want: order.Equal,
		},
		{
			Name: "a < b by Bar",
			Loc: exam.Here(),
			A: Struct2{Foo: 1, Bar: &pod.Vec[string]{"a"}},
			B: Struct2{Foo: 1, Bar: &pod.Vec[string]{"b"}},
			Want: order.Less,
		},
		{
			Name: "a > b by Bar",
			Loc: exam.Here(),
			A: Struct2{Foo: 1, Bar: &pod.Vec[string]{"b"}},
			B: Struct2{Foo: 1, Bar: &pod.Vec[string]{"a"}},
			Want: order.Greater,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			exam.Must(t, c.Want(Struct2Order(c.A, c.B)), c.A, c.B)
		})
	}
}