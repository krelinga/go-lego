package exam

import (
	"cmp"
	"fmt"

	"github.com/krelinga/go-lego/mirror"
)

type Failure struct {
	Args []any
}

func (e Failure) Error() string {
	args := make([]any, len(e.Args)+1)
	args[0] = "assertion failed with args:"
	copy(args[1:], e.Args)
	return fmt.Sprint(args...)
}

func Equal[T comparable](x, y T) error {
	if x == y {
		return nil
	}
	return Failure{
		Args: []any{x, y},
	}
}

func EqualFunc[T any](x, y T, f any) error {
	eq, err := mirror.WrapFunc2In1Out[T, T, bool](f)
	if err != nil {
		return err
	} else if eq(x, y) {
		return nil
	}
	return Failure{
		Args: []any{x, y},
	}
}

func NotEqual[T comparable](x, y T) error {
	if x != y {
		return nil
	}
	return Failure{
		Args: []any{x, y},
	}
}

func NotEqualFunc[T any](x, y T, f any) error {
	eq, err := mirror.WrapFunc2In1Out[T, T, bool](f)
	if err != nil {
		return err
	} else if !eq(x, y) {
		return nil
	}
	return Failure{
		Args: []any{x, y},
	}
}

func Greater[T cmp.Ordered](x, y T) error {
	if x > y {
		return nil
	}
	return Failure{
		Args: []any{x, y},
	}
}

func Less[T cmp.Ordered](x, y T) error {
	if x < y {
		return nil
	}
	return Failure{
		Args: []any{x, y},
	}
}

func GreaterEqual[T cmp.Ordered](x, y T) error {
	if x >= y {
		return nil
	}
	return Failure{
		Args: []any{x, y},
	}
}

func LessEqual[T cmp.Ordered](x, y T) error {
	if x <= y {
		return nil
	}
	return Failure{
		Args: []any{x, y},
	}
}

func True(condition bool) error {
	if condition {
		return nil
	}
	return Failure{
		Args: []any{condition},
	}
}

func Nil[T any](x T) error {
	isNil, err := mirror.IsNil(x)
	if err != nil {
		return err
	} else if isNil {
		return nil
	}
	return Failure{
		Args: []any{x},
	}
}

func NotNil[T any](x T) error {
	isNil, err := mirror.IsNil(x)
	if err != nil {
		return err
	} else if !isNil {
		return nil
	}
	return Failure{
		Args: []any{x},
	}
}
