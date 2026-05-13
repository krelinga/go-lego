package match

import (
	"fmt"
	"math"
	"reflect"
)

type EqualFunc struct {
	f any
}

func NewEqualFunc[T any](f func(T, T) bool) EqualFunc {
	return EqualFunc{f: f}
}

func (e *EqualFunc) checkInit() error {
	if e.f == nil {
		return fmt.Errorf("match.EqualFunc must be created with match.NewEqualFunc")
	}
	return nil
}

func (e *EqualFunc) checkType(t reflect.Type) error {
	if err := e.checkInit(); err != nil {
		return err
	}
	wantType := reflect.TypeOf(e.f).In(0)
	if !t.AssignableTo(wantType) {
		return fmt.Errorf("equal function expects a type assignable to %s but got type %s", wantType, t)
	}
	return nil
}

func (e *EqualFunc) equal(a, b any) (bool, error) {
	if err := e.checkType(reflect.TypeOf(a)); err != nil {
		return false, err
	}
	if err := e.checkType(reflect.TypeOf(b)); err != nil {
		return false, err
	}
	fVal := reflect.ValueOf(e.f)
	result := fVal.Call([]reflect.Value{reflect.ValueOf(a), reflect.ValueOf(b)})
	return result[0].Bool(), nil
}

type Approx struct {
	approx any
}

func NewApprox[T interface{ float32 | float64 }](approx T) Approx {
	return Approx{approx: approx}
}

func (e *Approx) checkInit() error {
	if e.approx == nil {
		return fmt.Errorf("match.Approx must be created with match.NewApprox")
	}
	return nil
}

func (e *Approx) checkType(t reflect.Type) error {
	if err := e.checkInit(); err != nil {
		return err
	}
	wantType := reflect.TypeOf(e.approx)
	if !t.ConvertibleTo(wantType) {
		return fmt.Errorf("approx value must be convertible to %s but got type %s", wantType, t)
	}
	return nil
}

func (e *Approx) equal(a, b any) (bool, error) {
	if err := e.checkType(reflect.TypeOf(a)); err != nil {
		return false, err
	}
	if err := e.checkType(reflect.TypeOf(b)); err != nil {
		return false, err
	}
	aVal := reflect.ValueOf(a).Convert(reflect.TypeOf(e.approx)).Float()
	bVal := reflect.ValueOf(b).Convert(reflect.TypeOf(e.approx)).Float()
	approxVal := reflect.ValueOf(e.approx).Float()
	diff := aVal - bVal
	return math.Abs(diff) <= approxVal, nil
}

type Equal struct {
	// The expected value to compare against.  Requirements:
	//   - Must be non-nil.
	//   - Must be a comparable type (unless Func is non-nil).
	Want any

	// Optional function to compare the expected and actual values.  If non-nil, this function will be used instead of the default equality check.  Requirements:
	//   - Must be a function type that accepts two parameters of the same type and returns a bool.
	//   - The type Want must be assignable to the parameter type of the function.
	Func *EqualFunc

	// Optional format function to use in the Result.  If non-nil, this function will be used to format the expected and actual values in the Result.  Requirements:
	//   - Must be a function type that accepts a single parameter and returns a string.
	//   - The parameter type of the function must be assignable from the type of Want and the type of the actual value.
	Format *Format

	// If non-nil, this value will be used as a tolerance when comparing floating point numbers.  Requirements:
	//   - Must be a floating point numeric type (e.g. float32, float64).
	//   - The type of the value must be assignable to the type of the actual value.
	//   - If Func is non-nil, this value must be nil.
	//   - If Want is not a floating point numeric type, this value must be nil.
	Approx *Approx
}

func (e Equal) Validate() error {
	if e.Want == nil {
		return fmt.Errorf("Equal matcher requires a non-nil Want value")
	}
	wantType := reflect.TypeOf(e.Want)
	if e.Approx != nil {
		if e.Func != nil {
			return fmt.Errorf("Equal matcher cannot have both Func and Approx set")
		}
		if err := e.Approx.checkType(wantType); err != nil {
			return err
		}
	}
	if e.Func != nil {
		if err := e.Func.checkType(wantType); err != nil {
			return err
		}
	} else if !wantType.Comparable() {
		return fmt.Errorf("Equal matcher requires a comparable type for Want when Func is not set, but got type %s", wantType)
	}
	return nil
}

func (e *Equal) Match(val any) (*Result, error) {
	h := &Helper{
		Meta: MetaHere(),
		Val:  val,
	}
	h.Context("expected", e.Want)
	if err := e.Validate(); err != nil {
		return nil, h.Fatal(err)
	}
	if e.Func != nil {
		return e.matchWithFunc(h, val)
	} else if e.Approx != nil {
		return e.matchApprox(h, val)
	} else {
		return e.matchWithEquality(h, val)
	}
}

func (e *Equal) matchWithFunc(h *Helper, val any) (*Result, error) {
	accept, err := e.Func.equal(e.Want, val)
	if err != nil {
		return nil, h.Fatal(err)
	}
	if accept {
		return h.Accept("values are equal according to custom function"), nil
	}
	return h.Reject("values are not equal according to custom function"), nil
}

func (e Equal) matchApprox(h *Helper, val any) (*Result, error) {
	accept, err := e.Approx.equal(e.Want, val)
	if err != nil {
		return nil, h.Fatal(err)
	}
	if accept {
		return h.Accept("values are approximately equal"), nil
	}
	return h.Reject("values are not approximately equal"), nil
}

func (e *Equal) matchWithEquality(h *Helper, val any) (*Result, error) {
	gotType := reflect.TypeOf(val)
	wantType := reflect.TypeOf(e.Want)

	if gotType != wantType {
		return nil, h.Fatalf("type mismatch: expected type %s but got type %s", wantType, gotType)
	}
	if val != e.Want {
		return h.Reject("values are not equal"), nil
	}
	return h.Accept("values are equal"), nil
}
