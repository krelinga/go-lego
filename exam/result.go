package exam

import (
	"cmp"
	"fmt"
	"reflect"

	"github.com/krelinga/go-lego/mirror"
)

type Arg struct {
	Name  string
	Value reflect.Value
}

func (a Arg) String() string {
	return fmt.Sprintf("%s: %#v", a.Name, a.Value)
}

type Args []Arg

func NewArgs1[T any](n1 string, v1 T) Args {
	return Args{
		{Name: n1, Value: mirror.ValueFor(v1)},
	}
}

func NewArgs2[T1 any, T2 any](n1 string, v1 T1, n2 string, v2 T2) Args {
	return Args{
		{Name: n1, Value: mirror.ValueFor(v1)},
		{Name: n2, Value: mirror.ValueFor(v2)},
	}
}

type Failure2 struct {
	Args Args
	Wrapped error
}

func (f *Failure2) clone() *Failure2 {
	return &Failure2{
		Args: f.Args,
		Wrapped: f.Wrapped,
	}
}

func (f *Failure2) Wrap(err error) *Failure2 {
	out := f.clone()
	out.Wrapped = err
	return out
}

func (f *Failure2) IsCritical() bool {
	return f.Wrapped != nil
}

func NewFailure1[T1 any](n1 string, t1 T1) *Failure2 {
	return &Failure2{
		Args: NewArgs1(n1, t1),
	}
}

func NewFailure2[T1, T2 any](n1 string, t1 T1, n2 string, t2 T2) *Failure2 {
	return &Failure2{
		Args: NewArgs2(n1, t1, n2, t2),
	}
}

type Failure struct {
	Args []any
}

func (e Failure) Error() string {
	args := make([]any, len(e.Args)+1)
	args[0] = "assertion failed with args:"
	copy(args[1:], e.Args)
	return fmt.Sprint(args...)
}

func Equal[T comparable](x, y T) *Failure2 {
	if x == y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

func EqualFunc[T any](x, y T, f any) *Failure2 {
	failure := NewFailure2("x", x, "y", y)
	eq, err := mirror.WrapFunc2In1Out[T, T, bool](f)
	if err != nil {
		return failure.Wrap(err)
	} else if eq(x, y) {
		return nil
	}
	return failure
}

func NotEqual[T comparable](x, y T) *Failure2 {
	if x != y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

func NotEqualFunc[T any](x, y T, f any) *Failure2 {
	failure := NewFailure2("x", x, "y", y)
	eq, err := mirror.WrapFunc2In1Out[T, T, bool](f)
	if err != nil {
		return failure.Wrap(err)
	} else if !eq(x, y) {
		return nil
	}
	return failure
}

func Greater[T cmp.Ordered](x, y T) *Failure2 {
	if x > y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

func Less[T cmp.Ordered](x, y T) *Failure2 {
	if x < y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

func GreaterEqual[T cmp.Ordered](x, y T) *Failure2 {
	if x >= y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

func LessEqual[T cmp.Ordered](x, y T) *Failure2 {
	if x <= y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

func True(condition bool) *Failure2 {
	if condition {
		return nil
	}
	return NewFailure1("condition", condition)
}

func Nil[T any](x T) *Failure2 {
	failure := NewFailure1("x", x)
	isNil, err := mirror.IsNil(x)
	if err != nil {
		return failure.Wrap(err)
	} else if isNil {
		return nil
	}
	return failure
}

func NotNil[T any](x T) *Failure2 {
	failure := NewFailure1("x", x)
	isNil, err := mirror.IsNil(x)
	if err != nil {
		return failure.Wrap(err)
	} else if !isNil {
		return nil
	}
	return failure
}
