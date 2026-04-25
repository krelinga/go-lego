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

type Failure struct {
	Args    Args
	Wrapped error
}

func (f *Failure) clone() *Failure {
	return &Failure{
		Args:    f.Args,
		Wrapped: f.Wrapped,
	}
}

func (f *Failure) Wrap(err error) *Failure {
	out := f.clone()
	out.Wrapped = err
	return out
}

func (f *Failure) IsCritical() bool {
	return f.Wrapped != nil
}

func NewFailure1[T1 any](n1 string, t1 T1) *Failure {
	return &Failure{
		Args: NewArgs1(n1, t1),
	}
}

func NewFailure2[T1, T2 any](n1 string, t1 T1, n2 string, t2 T2) *Failure {
	return &Failure{
		Args: NewArgs2(n1, t1, n2, t2),
	}
}

func Equal[T comparable](x, y T) *Failure {
	if x == y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

func EqualFunc[T any](x, y T, f any) *Failure {
	failure := NewFailure2("x", x, "y", y)
	eq, err := mirror.WrapFunc2In1Out[T, T, bool](f)
	if err != nil {
		return failure.Wrap(err)
	} else if eq(x, y) {
		return nil
	}
	return failure
}

func NotEqual[T comparable](x, y T) *Failure {
	if x != y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

func NotEqualFunc[T any](x, y T, f any) *Failure {
	failure := NewFailure2("x", x, "y", y)
	eq, err := mirror.WrapFunc2In1Out[T, T, bool](f)
	if err != nil {
		return failure.Wrap(err)
	} else if !eq(x, y) {
		return nil
	}
	return failure
}

func Greater[T cmp.Ordered](x, y T) *Failure {
	if x > y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

func Less[T cmp.Ordered](x, y T) *Failure {
	if x < y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

func GreaterEqual[T cmp.Ordered](x, y T) *Failure {
	if x >= y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

func LessEqual[T cmp.Ordered](x, y T) *Failure {
	if x <= y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

func True(condition bool) *Failure {
	if condition {
		return nil
	}
	return NewFailure1("condition", condition)
}

func Nil[T any](x T) *Failure {
	failure := NewFailure1("x", x)
	isNil, err := mirror.IsNil(x)
	if err != nil {
		return failure.Wrap(err)
	} else if isNil {
		return nil
	}
	return failure
}

func NotNil[T any](x T) *Failure {
	failure := NewFailure1("x", x)
	isNil, err := mirror.IsNil(x)
	if err != nil {
		return failure.Wrap(err)
	} else if !isNil {
		return nil
	}
	return failure
}
