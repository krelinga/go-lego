package mirror

import (
	"fmt"
	"reflect"
)

func Call(f reflect.Value, args []reflect.Value, outTypes []reflect.Type) ([]reflect.Value, error) {
	if !f.IsValid() {
		return nil, fmt.Errorf("f is not a valid reflect.Value")
	}
	for i, arg := range args {
		if !arg.IsValid() {
			return nil, fmt.Errorf("arg %d is not a valid reflect.Value", i)
		}
	}
	for i, outType := range outTypes {
		if outType == nil {
			return nil, fmt.Errorf("outTypes[%d] is nil", i)
		}
	}

	if f.Kind() != reflect.Func {
		return nil, fmt.Errorf("f must be a function")
	}
	if f.IsNil() {
		return nil, fmt.Errorf("f must be a non-nil function")
	}
	if f.Type().NumIn() != len(args) {
		return nil, fmt.Errorf("f expects %d arguments but got %d", f.Type().NumIn(), len(args))
	}
	if f.Type().NumOut() != len(outTypes) {
		return nil, fmt.Errorf("f returns %d values but expected %d", f.Type().NumOut(), len(outTypes))
	}

	for i, arg := range args {
		if !arg.Type().AssignableTo(f.Type().In(i)) {
			return nil, fmt.Errorf("argument %d of type %s is not assignable to function input of type %s", i, arg.Type(), f.Type().In(i))
		}
	}
	for i, outType := range outTypes {
		if f.Type().Out(i) != outType {
			return nil, fmt.Errorf("function output %d is of type %s but expected %s", i, f.Type().Out(i), outType)
		}
	}
	return f.Call(args), nil
}

func Call2In1Out[O, I1, I2 any](f any, i1 I1, i2 I2) (O, error) {
	args := []reflect.Value{ValueFor(i1), ValueFor(i2)}
	outTypes := []reflect.Type{reflect.TypeFor[O]()}
	results, err := Call(reflect.ValueOf(f), args, outTypes)
	if err != nil {
		var zero O
		return zero, err
	}
	return results[0].Interface().(O), nil
}
