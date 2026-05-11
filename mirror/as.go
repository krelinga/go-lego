package mirror

import (
	"fmt"
	"reflect"

	"github.com/krelinga/go-libs/zero"
)

func As[T any](val reflect.Value) (T, error) {
	if !val.IsValid() {
		return zero.For[T](), fmt.Errorf("value is invalid")
	}
	tType := reflect.TypeFor[T]()
	if !val.Type().AssignableTo(tType) {
		return zero.For[T](), fmt.Errorf("value of type %s cannot be assigned to type %s", val.Type(), tType)
	}
	if !val.CanInterface() {
		return zero.For[T](), fmt.Errorf("value cannot be interfaced")
	}
	return val.Interface().(T), nil
}
