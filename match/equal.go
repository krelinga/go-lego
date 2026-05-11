package match

import (
	"fmt"
	"reflect"
)

func Equal[T comparable](expected T) Matcher {
	meta := MetaHere()
	return FuncMatcher(func(val reflect.Value) (*Result, error) {
		helper := &Helper{
			Meta: meta,
			Val: val,
		}
		tVal, err := As[T](helper, val)
		if err != nil {
			return nil, err
		}
		if tVal != expected {
			return helper.Reject(fmt.Sprintf("%v != %v", tVal, expected)), nil
		}
		return helper.Accept(fmt.Sprintf("%v == %v", tVal, expected)), nil
	})
}
