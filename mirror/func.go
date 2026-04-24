package mirror

import (
	"fmt"
	"reflect"
)

func WrapFunc(f reflect.Value, inTypes []reflect.Type, outTypes []reflect.Type) (reflect.Value, error) {
	if !f.IsValid() {
		return reflect.Value{}, fmt.Errorf("f is not a valid reflect.Value")
	}
	for i, inType := range inTypes {
		if inType == nil {
			return reflect.Value{}, fmt.Errorf("inTypes[%d] is nil", i)
		}
	}
	for i, outType := range outTypes {
		if outType == nil {
			return reflect.Value{}, fmt.Errorf("outTypes[%d] is nil", i)
		}
	}

	if f.Kind() != reflect.Func {
		return reflect.Value{}, fmt.Errorf("f must be a function")
	}
	if f.IsNil() {
		return reflect.Value{}, fmt.Errorf("f must be a non-nil function")
	}
	if f.Type().NumIn() != len(inTypes) {
		return reflect.Value{}, fmt.Errorf("f expects %d arguments but got %d", f.Type().NumIn(), len(inTypes))
	}
	if f.Type().NumOut() != len(outTypes) {
		return reflect.Value{}, fmt.Errorf("f returns %d values but expected %d", f.Type().NumOut(), len(outTypes))
	}

	for i, inType := range inTypes {
		if !inType.AssignableTo(f.Type().In(i)) {
			return reflect.Value{}, fmt.Errorf("input type %d of type %s is not assignable to function input of type %s", i, inType, f.Type().In(i))
		}
	}
	for i, outType := range outTypes {
		if f.Type().Out(i) != outType {
			return reflect.Value{}, fmt.Errorf("function output %d is of type %s but expected %s", i, f.Type().Out(i), outType)
		}
	}

	wrapperType := reflect.FuncOf(inTypes, outTypes, false)
	wrapperFunc := reflect.MakeFunc(wrapperType, func(args []reflect.Value) []reflect.Value {
		return f.Call(args)
	})
	return wrapperFunc, nil
}

func WrapFunc2In1Out[O, I1, I2 any](f any) (func(I1, I2) O, error) {
	inTypes := []reflect.Type{reflect.TypeFor[I1](), reflect.TypeFor[I2]()}
	outTypes := []reflect.Type{reflect.TypeFor[O]()}
	wrapperFunc, err := WrapFunc(reflect.ValueOf(f), inTypes, outTypes)
	if err != nil {
		return nil, err
	}
	return wrapperFunc.Interface().(func(I1, I2) O), nil
}
