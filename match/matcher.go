package match

import (
	"errors"
	"fmt"
	"reflect"
)

type Matcher interface {
	Match(reflect.Value) error
}

func New[T any](f func(T) error) Matcher {
	return typedFuncMatcher[T]{f: f}
}

type typedFuncMatcher[T any] struct {
	f func(T) error
}

func (m typedFuncMatcher[T]) Match(val reflect.Value) error {
	if !val.IsValid() {
		return errors.New("invalid value")
	}
	tType := reflect.TypeFor[T]()
	if !val.Type().AssignableTo(tType) {
		return fmt.Errorf("value is of type %s, expected a type assignable to %s", val.Type(), tType)
	}
	if !val.CanInterface() {
		return fmt.Errorf("value of type %s cannot be interfaced", val.Type())
	}
	var out T
	outPtrVal := reflect.ValueOf(&out)
	outPtrVal.Elem().Set(val)
	return m.f(out)
}
