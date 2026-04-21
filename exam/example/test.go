package example_test

import (
	"testing"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/order"
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
	e := exam.New(t)

	e.Run("simple int comparison", func(e exam.E) {
		exam.Equal(e, 1, 1).Must()
		exam.GreaterThan(e, 1, 0)
	})

	e.Run("struct comparison", func(e exam.E) {
		a := FooStruct{Foo: 1, Bar: "a"}
		b := FooStruct{Foo: 1, Bar: "b"}
		exam.NewPred2(func(a, b FooStructView) bool {
			return order.Greater(a, b, FooOrder)
		})(e, b, a).Must()
		exam.NewPred2(FooEqual)(e, a, b).Must()
	})

	e.Run("Map of struct comparison", func(e exam.E) {
		a := &pod.Map[string, FooStructView]{
			"a": FooStruct{Foo: 1, Bar: "a"},
			"b": FooStruct{Foo: 2, Bar: "b"},
		}
		b := &pod.Map[string, FooStructView]{
			"a": FooStruct{Foo: 1, Bar: "a"},
			"b": FooStruct{Foo: 2, Bar: "b"},
		}
		exam.NewPred2(func(a, b pod.MapView[string, FooStructView]) bool {
			return pod.MapEqualFunc(a, b, FooEqual)
		})(e, a, b).Must()
	})
}
