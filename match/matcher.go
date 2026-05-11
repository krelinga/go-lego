package match

import (
	"fmt"
	"reflect"
)

type Matcher interface {
	Match(reflect.Value) (*Result, error)
}

type FuncMatcher func(val reflect.Value) (*Result, error)

func (f FuncMatcher) Match(val reflect.Value) (*Result, error) {
	return f(val)
}

type Format func(val reflect.Value) (string, error)

func WithFormat(m Matcher, format Format) Matcher {
	if fm, ok := m.(*formatMatcher); ok {
		fm.f = format
		return fm
	}
	return &formatMatcher{
		m: m,
		f: format,
	}
}

type formatMatcher struct {
	m Matcher
	f Format
}

func (fm *formatMatcher) format(val reflect.Value) string {
	if fm.f == nil {
		return defaultFormat(val)
	}
	if s, err := fm.f(val); err != nil {
		return fmt.Sprintf("%s (error formatting value: %v)", defaultFormat(val), err)
	} else {
		return s
	}
}

func (fm *formatMatcher) Match(val reflect.Value) (*Result, error) {
	result, err := fm.m.Match(val)
	if result != nil {
		result.format = fm.format
	}
	if fatalError, ok := err.(*FatalError); ok {
		fatalError.format = fm.format
	}
	return result, err
}

func defaultFormat(val reflect.Value) string {
	if !val.IsValid() {
		return "<invalid value>"
	}
	if !val.CanInterface() {
		return fmt.Sprintf("<uninterfaceable value of type %s>", val.Type())
	}
	if val.Kind() == reflect.String {
		return fmt.Sprintf("%q", val.String())
	}
	return fmt.Sprintf("%v", val.Interface())
}