package exam

import (
	"cmp"
	"fmt"
	"reflect"

	"github.com/krelinga/go-lego/mirror"
)

// FmtFunc is the signature of a custom formatting function used by Arg.
// It converts a reflect.Value to its string representation, or returns an error.
type FmtFunc = func(reflect.Value) (string, error)

// Arg is a named value captured as part of a Failure. It records the argument
// name, its reflected value, and an optional custom formatting function.
type Arg struct {
	// Name is the argument label shown in failure messages.
	Name string
	// Value is the reflected argument value.
	Value reflect.Value
	// Fmt, if non-nil, is called to format Value in failure messages.
	// When nil, fmt.Sprintf("%#v", ...) is used instead.
	Fmt FmtFunc
}

// ToString formats the argument as "name: value" using Fmt if set, or %#v otherwise.
func (a Arg) ToString() (string, error) {
	var valueStr string
	if a.Fmt != nil {
		var err error
		valueStr, err = a.Fmt(a.Value)
		if err != nil {
			return "", fmt.Errorf("formatting argument %#v: %w", a.Name, err)
		}
	} else {
		valueStr = fmt.Sprintf("%#v", a.Value)
	}
	return fmt.Sprintf("%s: %s", a.Name, valueStr), nil
}

// Args is a slice of Arg values attached to a Failure.
type Args []Arg

// NewArgs1 constructs a single-element Args from a named value.
func NewArgs1[T any](n1 string, v1 T) Args {
	return Args{
		{Name: n1, Value: mirror.ValueFor(v1)},
	}
}

// NewArgs2 constructs a two-element Args from two named values.
func NewArgs2[T1 any, T2 any](n1 string, v1 T1, n2 string, v2 T2) Args {
	return Args{
		{Name: n1, Value: mirror.ValueFor(v1)},
		{Name: n2, Value: mirror.ValueFor(v2)},
	}
}

// Failure represents a failed assertion. A nil *Failure indicates success.
type Failure struct {
	// Args holds the named values that were compared when the assertion failed.
	Args Args
	// Wrapped, if non-nil, holds an underlying error that caused a structural
	// problem (e.g. an incompatible function passed to EqualFunc). A Failure
	// with a non-nil Wrapped is treated as fatal by both Try and Must.
	Wrapped error
}

func (f *Failure) clone() *Failure {
	return &Failure{
		Args:    f.Args,
		Wrapped: f.Wrapped,
	}
}

// Wrap returns a copy of f with Wrapped set to err, making the failure critical.
func (f *Failure) Wrap(err error) *Failure {
	out := f.clone()
	out.Wrapped = err
	return out
}

// IsCritical reports whether the failure wraps an underlying error.
// Critical failures are always reported with t.Fatal regardless of whether
// Try or Must was used.
func (f *Failure) IsCritical() bool {
	return f.Wrapped != nil
}

// NewFailure1 returns a Failure carrying a single named argument.
func NewFailure1[T1 any](n1 string, t1 T1) *Failure {
	return &Failure{
		Args: NewArgs1(n1, t1),
	}
}

// NewFailure2 returns a Failure carrying two named arguments.
func NewFailure2[T1, T2 any](n1 string, t1 T1, n2 string, t2 T2) *Failure {
	return &Failure{
		Args: NewArgs2(n1, t1, n2, t2),
	}
}

// Equal returns nil when x == y, or a Failure containing both values otherwise.
func Equal[T comparable](x, y T) *Failure {
	if x == y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

// EqualFunc returns nil when f reports x and y as equal, or a Failure otherwise.
// f must be a func(U, U) bool (where T is assignable to U); a structural error is
// wrapped into the Failure if it is not.
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

// NotEqual returns nil when x != y, or a Failure containing both values otherwise.
func NotEqual[T comparable](x, y T) *Failure {
	if x != y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

// NotEqualFunc returns nil when f reports x and y as unequal, or a Failure otherwise.
// f must be a func(U, U) bool (where T is assignable to U); a structural error is
// wrapped into the Failure if it is not.
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

// Greater returns nil when x > y, or a Failure containing both values otherwise.
func Greater[T cmp.Ordered](x, y T) *Failure {
	if x > y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

// Less returns nil when x < y, or a Failure containing both values otherwise.
func Less[T cmp.Ordered](x, y T) *Failure {
	if x < y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

// GreaterEqual returns nil when x >= y, or a Failure containing both values otherwise.
func GreaterEqual[T cmp.Ordered](x, y T) *Failure {
	if x >= y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

// LessEqual returns nil when x <= y, or a Failure containing both values otherwise.
func LessEqual[T cmp.Ordered](x, y T) *Failure {
	if x <= y {
		return nil
	}
	return NewFailure2("x", x, "y", y)
}

// True returns nil when condition is true, or a Failure otherwise.
func True(condition bool) *Failure {
	if condition {
		return nil
	}
	return NewFailure1("condition", condition)
}

// Nil returns nil when x is nil, or a Failure otherwise.
// x must be any nilable type: pointer, interface, slice, map, channel, or function.
// A structural error is returned for non-nilable types.
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

// NotNil returns nil when x is non-nil, or a Failure otherwise.
// x must be any nilable type: pointer, interface, slice, map, channel, or function.
// A structural error is returned for non-nilable types.
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

// Implements returns nil when Got implements IFace, or a Failure otherwise.
func Implements[Got, IFace any]() *Failure {
	failure := NewFailure2("Got", reflect.TypeFor[Got](), "IFace", reflect.TypeFor[IFace]())
	if reflect.TypeFor[IFace]().Kind() != reflect.Interface {
		return failure.Wrap(fmt.Errorf("type parameter IFace must be an interface"))
	}
	var gotZero Got
	_, ok := any(gotZero).(IFace)
	if ok {
		return nil
	}
	return failure
}
