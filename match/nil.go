package match

import "reflect"

func Nil() Matcher {
	meta := MetaHere()
	return FuncMatcher(func(val reflect.Value) (*Result, error) {
		helper := &Helper{
			Meta: meta,
			Val: val,
		}
		if err := helper.CheckValid(); err != nil {
			return nil, err
		}
		switch val.Kind() {
		case reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
			if val.IsNil() {
				return helper.Accept("value is nil"), nil
			}
			return helper.Reject("value is not nil"), nil
		default:
			return nil, helper.Fatalf("value of kind %s cannot be nil", val.Kind())
		}
	})
}