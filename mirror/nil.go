package mirror

import (
	"fmt"
	"reflect"
)

func IsNil[T any](x T) (bool, error) {
	val := ValueFor(x)
	if !val.IsValid() {
		return false, fmt.Errorf("x is not a valid value")
	}
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return val.IsNil(), nil
	default:
		return false, fmt.Errorf("value of type %s cannot be nil", val.Type())
	}
}