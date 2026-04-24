package mirror_test

import (
	"reflect"
	"testing"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/mirror"
)

func add(x, y int) int {
	return x + y
}

type fooInterface interface {
	Foo() int
}

func fooFunc(f fooInterface) int {
	if f == nil {
		return 0
	}
	return f.Foo()
}

type fooImpl int

func (f fooImpl) Foo() int {
	return int(f)
}

var invalidValue = reflect.Value{}
var nilFunc func(int, int) int = nil
var nilFooInterface fooInterface = nil

func TestWrap(t *testing.T) {
	cases := []struct {
		Name      string
		Loc       exam.Loc
		F         reflect.Value
		Args      []reflect.Value
		InTypes   []reflect.Type
		OutTypes  []reflect.Type
		WantError bool
		WantOut   []reflect.Value
	}{
		{
			Name:      "invalid func",
			Loc:       exam.Here(),
			F:         invalidValue,
			Args:      []reflect.Value{mirror.ValueFor(1), mirror.ValueFor(2)},
			InTypes:   []reflect.Type{reflect.TypeFor[int](), reflect.TypeFor[int]()},
			OutTypes:  []reflect.Type{reflect.TypeFor[int]()},
			WantError: true,
		},
		{
			Name:      "nil func",
			Loc:       exam.Here(),
			F:         mirror.ValueFor(nilFunc),
			Args:      []reflect.Value{mirror.ValueFor(1), mirror.ValueFor(2)},
			InTypes:   []reflect.Type{reflect.TypeFor[int](), reflect.TypeFor[int]()},
			OutTypes:  []reflect.Type{reflect.TypeFor[int]()},
			WantError: true,
		},
		{
			Name:      "not a function",
			Loc:       exam.Here(),
			F:         mirror.ValueFor(42),
			Args:      []reflect.Value{mirror.ValueFor(1), mirror.ValueFor(2)},
			InTypes:   []reflect.Type{reflect.TypeFor[int](), reflect.TypeFor[int]()},
			OutTypes:  []reflect.Type{reflect.TypeFor[int]()},
			WantError: true,
		},
		{
			Name:      "wrong number of args",
			Loc:       exam.Here(),
			F:         mirror.ValueFor(add),
			Args:      []reflect.Value{mirror.ValueFor(1)},
			InTypes:   []reflect.Type{reflect.TypeFor[int]()},
			OutTypes:  []reflect.Type{reflect.TypeFor[int]()},
			WantError: true,
		},
		{
			Name:      "wrong arg type",
			Loc:       exam.Here(),
			F:         mirror.ValueFor(add),
			Args:      []reflect.Value{mirror.ValueFor(1), mirror.ValueFor("2")},
			InTypes:   []reflect.Type{reflect.TypeFor[int](), reflect.TypeFor[string]()},
			OutTypes:  []reflect.Type{reflect.TypeFor[int]()},
			WantError: true,
		},
		{
			Name:      "wrong number of out types",
			Loc:       exam.Here(),
			F:         mirror.ValueFor(add),
			Args:      []reflect.Value{mirror.ValueFor(1), mirror.ValueFor(2)},
			InTypes:   []reflect.Type{reflect.TypeFor[int](), reflect.TypeFor[int]()},
			OutTypes:  []reflect.Type{},
			WantError: true,
		},
		{
			Name:      "wrong out type",
			Loc:       exam.Here(),
			F:         mirror.ValueFor(add),
			Args:      []reflect.Value{mirror.ValueFor(1), mirror.ValueFor(2)},
			InTypes:   []reflect.Type{reflect.TypeFor[int](), reflect.TypeFor[int]()},
			OutTypes:  []reflect.Type{reflect.TypeFor[string]()},
			WantError: true,
		},
		{
			Name:      "valid wrapping",
			Loc:       exam.Here(),
			F:         mirror.ValueFor(add),
			Args:      []reflect.Value{mirror.ValueFor(1), mirror.ValueFor(2)},
			InTypes:   []reflect.Type{reflect.TypeFor[int](), reflect.TypeFor[int]()},
			OutTypes:  []reflect.Type{reflect.TypeFor[int]()},
			WantError: false,
			WantOut:   []reflect.Value{mirror.ValueFor(3)},
		},
		{
			Name:      "interface wrap",
			Loc:       exam.Here(),
			F:         mirror.ValueFor(fooFunc),
			Args:      []reflect.Value{mirror.ValueFor(fooImpl(42))},
			InTypes:   []reflect.Type{reflect.TypeFor[fooInterface]()},
			OutTypes:  []reflect.Type{reflect.TypeFor[int]()},
			WantError: false,
			WantOut:   []reflect.Value{mirror.ValueFor(42)},
		},
		{
			Name:      "nil interface wrap",
			Loc:       exam.Here(),
			F:         mirror.ValueFor(fooFunc),
			Args:      []reflect.Value{mirror.ValueFor(nilFooInterface)},
			InTypes:   []reflect.Type{reflect.TypeFor[fooInterface]()},
			OutTypes:  []reflect.Type{reflect.TypeFor[int]()},
			WantError: false,
			WantOut:   []reflect.Value{mirror.ValueFor(0)},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			w, err := mirror.WrapFunc(c.F, c.InTypes, c.OutTypes)
			if c.WantError {
				exam.Try(t, exam.NotNil(err))
			} else {
				exam.Must(t, exam.Nil(err))
				out := w.Call(c.Args)
				exam.Must(t, exam.Equal(len(out), len(c.WantOut)))
				for i := range out {
					exam.Must(t, exam.True(out[i].CanInterface()))
					exam.Must(t, exam.Equal(out[i].Interface(), c.WantOut[i].Interface()))
				}
			}
		})
	}
}

func TestWrapFunc2In1Out(t *testing.T) {
	t.Run("valid call", func(t *testing.T) {
		f, err := mirror.WrapFunc2In1Out[int, int, int](add)
		exam.Must(t, exam.Nil(err))
		exam.Try(t, exam.Equal(f(1, 2), 3))
	})

	t.Run("invalid func", func(t *testing.T) {
		var notFunc any = 42
		_, err := mirror.WrapFunc2In1Out[int, int, int](notFunc)
		exam.Try(t, exam.NotNil(err))
	})
}
