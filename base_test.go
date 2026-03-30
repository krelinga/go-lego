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

type ExampleView interface {
	GetString() string
	GetInt() int
}

type Example struct {
	String string
	Int    int
}

func (e Example) GetString() string {
	return e.String
}

func (e Example) GetInt() int {
	return e.Int
}

func (e Example) View() ExampleView {
	return e
}

type ExampleSlice struct {
	*lego.Slice[Example]
}

func (s ExampleSlice) View() lego.FixedSlice[ExampleView] {
	return lego.ViewSlice(s)
}