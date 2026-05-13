package match

import (
	"fmt"
	"reflect"

	"github.com/krelinga/go-libs/valid"
)

type Matcher interface {
	valid.Validator

	Match(any) (*Result, error)
}

type Format struct {
	f any
}

func NewFormat[T any](f func(T) (string, error)) *Format {
	return &Format{f: f}
}

func (fa *Format) checkInit() error {
	if fa.f == nil {
		return fmt.Errorf("match.Format must be created with match.NewFormat")
	}
	return nil
}

func (fa *Format) checkType(t reflect.Type) error {
	if err := fa.checkInit(); err != nil {
		return err
	}
	wantType := reflect.TypeOf(fa.f).In(0)
	if !t.AssignableTo(wantType) {
		return fmt.Errorf("format function expects a type assignable to %s but got type %s", wantType, t)
	}
	return nil
}

func (fa *Format) format(a any) (string, error) {
	if a == nil {
		return "nil", nil
	}
	val := reflect.ValueOf(a)
	if err := fa.checkType(val.Type()); err != nil {
		return "", err
	}
	fVal := reflect.ValueOf(fa.f)
	result := fVal.Call([]reflect.Value{val})
	strVal := result[0].String()
	errVal := result[1].Interface().(error)
	return strVal, errVal
}
