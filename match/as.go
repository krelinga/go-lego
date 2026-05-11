package match

import (
	"reflect"

	"github.com/krelinga/go-libs/mirror"
	"github.com/krelinga/go-libs/zero"
)

func As[T any](h *Helper, val reflect.Value) (T, error) {
	if err := h.CheckValid(); err != nil {
		return zero.For[T](), err
	}
	tVal, ok := mirror.As[T](val)
	if !ok {
		return tVal, h.Fatalf("value is not of expected type %T", tVal)
	}
	return tVal, nil
}