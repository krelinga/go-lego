package match

import (
	"reflect"

	"github.com/krelinga/go-libs/mirror"
)

func As[T any](h *Helper, val reflect.Value) (T, error) {
	tVal, ok := mirror.As[T](val)
	if !ok {
		return tVal, h.Fatalf("value is not of expected type %T", tVal)
	}
	return tVal, nil
}