package mirror

import (
	"reflect"

	"github.com/krelinga/go-libs/zero"
)

func As[T any](val reflect.Value) (T, bool) {
	if !val.IsValid() || !val.Type().AssignableTo(reflect.TypeFor[T]()) {
		return zero.For[T](), false
	}
	return val.Interface().(T), true
}
