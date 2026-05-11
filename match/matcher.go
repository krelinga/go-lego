package match

import (
	"reflect"
)

type Matcher interface {
	Match(reflect.Value) (*Result, error)
}

type FuncMatcher func(val reflect.Value) (*Result, error)

func (f FuncMatcher) Match(val reflect.Value) (*Result, error) {
	return f(val)
}
