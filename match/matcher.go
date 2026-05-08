package match

import (
	"fmt"
	"reflect"
)

type Matcher interface {
	Match(reflect.Value) *Result
}

func NewFunc[T any](m Meta, f func(*Result, T)) Matcher {
	return funcMatcher[T]{
		Meta: m,
		F:    f,
	}
}

type funcMatcher[T any] struct {
	Meta Meta
	F    func(*Result, T)
}

func (m funcMatcher[T]) Match(val reflect.Value) (result *Result) {
	result = &Result{
		Meta: m.Meta,
		Val:  val,
	}

	if !val.IsValid() {
		result.Err = fmt.Errorf("invalid value")
		return
	}

	var typedVal T
	tType := reflect.TypeFor[T]()
	if !val.Type().AssignableTo(tType) {
		result.Err = fmt.Errorf("expected type %s but got %s", tType, val.Type())
		return
	}
	reflect.ValueOf(&typedVal).Elem().Set(val)

	m.F(result, typedVal)
	return
}

func NewValFunc(m Meta, f func(*Result, reflect.Value)) Matcher {
	return funcValMatcher{
		Meta: m,
		F:    f,
	}
}

type funcValMatcher struct {
	Meta Meta
	F    func(*Result, reflect.Value)
}

func (m funcValMatcher) Match(val reflect.Value) (result *Result) {
	result = &Result{
		Meta: m.Meta,
		Val:  val,
	}

	if !val.IsValid() {
		result.Err = fmt.Errorf("invalid value")
		return
	}

	m.F(result, val)
	return
}