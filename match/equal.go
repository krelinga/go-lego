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
