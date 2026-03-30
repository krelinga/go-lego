package lego_test

import (
	"reflect"
	"testing"

	"github.com/krelinga/go-lego"
)

type FailureReporter interface {
	Error(...any)
	Errorf(string, ...any)
	Fatal(...any)
	Fatalf(string, ...any)
	Helper()
}

func implements[T any, I any](r FailureReporter) bool {
	r.Helper()
	iType := reflect.TypeFor[I]()
	if iType.Kind() != reflect.Interface {
		r.Fatalf("I must be an interface type, got %s", iType)
	}
	tType := reflect.TypeFor[T]()
	if !tType.Implements(iType) {
		r.Errorf("%s does not implement %s", tType, iType)
	}
	return true
}

func panics(r FailureReporter, f func()) (panicked bool) {
	r.Helper()
	defer func() {
		r.Helper()
		if recover() == nil {
			r.Errorf("Expected function to panic, but it did not")
		} else {
			panicked = true
		}
	}()
	f()
	return
}

type ExampleView struct {
	e *Example
}

func (v ExampleView) String() string {
	return v.e.String
}

func (v ExampleView) Int() int {
	return v.e.Int
}

func (v ExampleView) Equal(other ExampleView) bool {
	return v.String() == other.String() && v.Int() == other.Int()
}

type Example struct {
	String string
	Int    int
}

func (e *Example) Equal(other *Example) bool {
	return lego.EqualViewer(e, other)
}

func (e *Example) View() ExampleView {
	return ExampleView{e}
}

type ExampleSliceView struct {
	lego.FixedSlice[ExampleView]
}

func (v ExampleSliceView) Equal(other ExampleSliceView) bool {
	return lego.EqualSlice(v, other)
}

type ExampleSlice struct {
	*lego.Slice[*Example]
}

func (s *ExampleSlice) View() ExampleSliceView {
	return ExampleSliceView{lego.ViewSlice(s)}
}

func (s *ExampleSlice) Equal(other *ExampleSlice) bool {
	return lego.EqualViewer(s, other)
}

func NewExampleSlice(elements ...*Example) *ExampleSlice {
	return &ExampleSlice{lego.NewSlice(elements...)}
}

func TestExample(t *testing.T) {
	t.Run("implements", func(t *testing.T) {
		implements[*Example, lego.Equaler[*Example]](t)
		implements[*Example, lego.Viewer[ExampleView]](t)
	})

	t.Run("equal", func(t *testing.T) {
		e1 := &Example{String: "a", Int: 1}
		e2 := &Example{String: "a", Int: 1}
		e3 := &Example{String: "a", Int: 2}

		if !e1.Equal(e2) {
			t.Errorf("Expected e1 to equal e2, but they are not equal")
		}
		if e1.Equal(e3) {
			t.Errorf("Expected e1 to not equal e3, but they are equal")
		}

		e1v := e1.View()
		e2v := e2.View()
		e3v := e3.View()

		if !e1v.Equal(e2v) {
			t.Errorf("Expected e1v to equal e2v, but they are not equal")
		}
		if e1v.Equal(e3v) {
			t.Errorf("Expected e1v to not equal e3v, but they are equal")
		}
	})

	t.Run("view", func(t *testing.T) {
		e := &Example{String: "a", Int: 1}
		v := e.View()
		if v.String() != "a" {
			t.Errorf("Expected String() to return 'a', got '%s'", v.String())
		}
		if v.Int() != 1 {
			t.Errorf("Expected Int() to return 1, got %d", v.Int())
		}
	})
}

func TestExampleSlice(t *testing.T) {
	t.Run("implements", func(t *testing.T) {
		implements[*ExampleSlice, lego.Equaler[*ExampleSlice]](t)
		implements[*ExampleSlice, lego.Viewer[ExampleSliceView]](t)
		implements[*ExampleSlice, lego.LenLister[lego.Pair[int, *Example]]](t)
		implements[*ExampleSlice, lego.FixedSlice[*Example]](t)

		implements[ExampleSliceView, lego.Equaler[ExampleSliceView]](t)
		implements[ExampleSliceView, lego.FixedSlice[ExampleView]](t)
		implements[ExampleSliceView, lego.LenLister[lego.Pair[int, ExampleView]]](t)
	})
}
