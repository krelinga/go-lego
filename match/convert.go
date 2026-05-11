package match

import (
	"reflect"
)

func Convert[T any](m Matcher) Matcher {
	meta := MetaHere()
	return FuncMatcher(func(val reflect.Value) (*Result, error) {
		helper := &Helper{
			Meta: meta,
			Val:  val,
		}
		if err := helper.CheckValid(); err != nil {
			return nil, err
		}
		tType := reflect.TypeFor[T]()
		if !val.CanConvert(tType) {
			return nil, helper.Fatalf("value of type %s cannot be converted to %s", val.Type(), tType)
		}
		accept, err := helper.Child("", val.Convert(tType), m)
		if err != nil {
			return nil, err
		}
		if accept {
			return helper.Accept("child accepted converted value"), nil
		}
		return helper.Reject("child rejected converted value"), nil
	})
}