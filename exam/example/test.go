package example_test

import (
	"testing"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/match"
	"github.com/krelinga/go-lego/order"
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

func FooEqual(a, b FooStruct) bool {
	return a.Foo == b.Foo
}

func TestExample(t *testing.T) {
	e := exam.New(t)

	e.Run("simple int comparison", func(e exam.E) {
		exam.Equal(e, 1, 1).Must()
		exam.Match(e, 1, match.GreaterThan(0))
	})
	
	e.Run("struct comparison", func(e exam.E) {
		a := FooStruct{Foo: 1, Bar: "a"}
		b := FooStruct{Foo: 1, Bar: "b"}
		exam.Match(e, b, match.GreaterThan(a).Using(FooOrder))
		exam.Match(e, a, match.Equal(b).Using(FooEqual))
	})
}
