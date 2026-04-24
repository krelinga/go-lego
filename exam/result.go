package exam

import (
	"cmp"
	"errors"
	"fmt"
	"reflect"

	"github.com/krelinga/go-lego/mirror"
)

var ErrFailed = errors.New("failed")

func AsError(ok bool) error {
	if ok {
		return nil
	}
	return ErrFailed
}

type Result struct {
	Error error
	Args []any
}

func (r Result) isFatal() bool {
	return r.Error != nil && r.Error != ErrFailed
}

func (r Result) passed() bool {
	return r.Error == nil
}

func Equal[T comparable](x, y T) Result {
	return Result{
		Error: AsError(x == y),
		Args: []any{x, y},
	}
}

func EqualFunc[T any](x, y T, f any) Result {
	eq, err := mirror.Call2In1Out[bool](f, x, y)
	if err == nil && !eq {
		err = ErrFailed
	}
	return Result{
		Error: err,
		Args: []any{x, y},
	}
}

func NotEqual[T comparable](x, y T) Result {
	return Result{
		Error: AsError(x != y),
		Args: []any{x, y},
	}
}

func NotEqualFunc[T any](x, y T, f any) Result {
	eq, err := mirror.Call2In1Out[bool](f, x, y)
	if err == nil && eq {
		err = ErrFailed
	}
	return Result{
		Error: err,
		Args: []any{x, y},
	}
}

func Greater[T cmp.Ordered](x, y T) Result {
	return Result{
		Error: AsError(x > y),
		Args: []any{x, y},
	}
}

func Less[T cmp.Ordered](x, y T) Result {
	return Result{
		Error: AsError(x < y),
		Args: []any{x, y},
	}
}
func GreaterEqual[T cmp.Ordered](x, y T) Result {
	return Result{
		Error: AsError(x >= y),
		Args: []any{x, y},
	}
}

func LessEqual[T cmp.Ordered](x, y T) Result {
	return Result{
		Error: AsError(x <= y),
		Args: []any{x, y},
	}
}

func True(condition bool) Result {
	return Result{
		Error: AsError(condition),
		Args: []any{condition},
	}
}

func isNil(x any) (bool, error) {
	if x == nil {
		return true, nil
	}
	val := reflect.ValueOf(x)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return val.IsNil(), nil
	default:
		return false, fmt.Errorf("value of type %T cannot be nil", x)
	}
}

func Nil(x any) Result {
	isNil, err := isNil(x)
	if err == nil && !isNil {
		err = ErrFailed
	}
	return Result{
		Error: err,
		Args:  []any{x},
	}
}

func NotNil(x any) Result {
	isNil, err := isNil(x)
	if err == nil && isNil {
		err = ErrFailed
	}
	return Result{
		Error: err,
		Args:  []any{x},
	}
}