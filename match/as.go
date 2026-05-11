package match

import (
	"reflect"

	"github.com/krelinga/go-libs/mirror"
)

func As[T any](h *Helper, val reflect.Value) (T, error) {
	tVal, err := mirror.As[T](val)
	if err != nil {
		return tVal, h.Fatal(err)
	}
	return tVal, nil
}
