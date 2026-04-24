package exam

import (
	"cmp"
	"errors"

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

func Nil[T any](x T) Result {
	isNil, err := mirror.IsNil(x)
	if err == nil && !isNil {
		err = ErrFailed
	}
	return Result{
		Error: err,
		Args:  []any{x},
	}
}

func NotNil[T any](x T) Result {
	isNil, err := mirror.IsNil(x)
	if err == nil && isNil {
		err = ErrFailed
	}
	return Result{
		Error: err,
		Args:  []any{x},
	}
}