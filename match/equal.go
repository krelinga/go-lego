package match

import (
	"reflect"

	"github.com/krelinga/go-libs/mirror"
)

func Equal[T comparable](expected T) Matcher {
	meta := MetaHere()
	return FuncMatcher(func(val reflect.Value) (*Result, error) {
		helper := &Helper{
			Meta: meta,
			Val:  val,
		}
		helper.Context("expected", mirror.ValueFor(expected))
		tVal, err := As[T](helper, val)
		if err != nil {
			return nil, err
		}
		if tVal != expected {
			return helper.Reject("values are not equal"), nil
		}
		return helper.Accept("values are equal"), nil
	})
}

type EqualAny struct {
	// The expected value to compare against.  Requirements:
	//   - Must be non-nil.
	//   - Must be a comparable type (unless Func is non-nil).
	Want any

	// Optional function to compare the expected and actual values.  If non-nil, this function will be used instead of the default equality check.  Requirements:
	//   - Must be a function type that accepts two parameters of the same type and returns a bool.
	//   - The type Want must be assignable to the parameter type of the function.
	Func any

	// Optional format function to use in the Result.  If non-nil, this function will be used to format the expected and actual values in the Result.  Requirements:
	//   - Must be a function type that accepts a single parameter and returns a string.
	//   - The parameter type of the function must be assignable from the type of Want and the type of the actual value.
	Format any

	// If non-nil, this value will be used as a tolerance when comparing floating point numbers.  Requirements:
	//   - Must be a floating point numeric type (e.g. float32, float64).
	//   - The type of the value must be assignable to the type of the actual value.
	//   - If Func is non-nil, this value must be nil.
	//   - If Want is not a floating point numeric type, this value must be nil.
	Approx any
}

func (e *EqualAny) Match(val reflect.Value) (*Result, error) {
	h := &Helper{
		Meta: MetaHere(),
		Val:  val,
	}
	h.Context("expected", reflect.ValueOf(e.Want))
	if err := h.CheckValid(); err != nil {
		return nil, err
	}
	if e.Func != nil {
		return e.matchWithFunc(h, val)
	} else if e.Approx != nil {
		return e.matchApprox(h, val)
	} else if reflect.TypeOf(e.Want).Comparable() {
		return e.matchWithEquality(h, val)
	}
}

func (e *EqualAny) matchWithFunc(h *Helper, val reflect.Value) (*Result, error) {
	return nil, nil // TODO
}

func (e *EqualAny) matchWithEquality(h *Helper, val reflect.Value) (*Result, error) {
	return nil, nil // TODO
}

func (e EqualAny) matchApprox(h *Helper, val reflect.Value) (*Result, error) {
	return nil, nil // TODO
}