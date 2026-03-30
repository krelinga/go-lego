package lego_test

import (
	"reflect"

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

func (e Example) Equal(other Example) bool {
	return lego.EqualViewer(e, other)
}

func (e Example) View() ExampleView {
	return ExampleView{&e}
}

type ExampleSliceView struct {
	lego.FixedSlice[ExampleView]
}

func (v ExampleSliceView) Equal(other ExampleSliceView) bool {
	return lego.EqualSlice(v, other)
}

type ExampleSlice struct {
	*lego.Slice[Example]
}

func (s ExampleSlice) View() ExampleSliceView {
	return ExampleSliceView{lego.ViewSlice(s)}
}

func (s ExampleSlice) Equal(other ExampleSlice) bool {
	return lego.EqualViewer(s, other)
}

func NewExampleSlice(elements ...Example) ExampleSlice {
	return ExampleSlice{lego.NewSlice(elements...)}
}
