package match

import "reflect"

func Pointee(m Matcher) Matcher {
	meta := MetaHere()
	return FuncMatcher(func(val reflect.Value) (*Result, error) {
		helper := &Helper{
			Meta: meta,
			Val:  val,
		}
		if err := helper.CheckValid(); err != nil {
			return nil, err
		}
		if val.Kind() != reflect.Ptr {
			return nil, helper.Fatalf("value of type %s is not a pointer", val.Type())
		}
		if val.IsNil() {
			return helper.Reject("pointer is nil"), nil
		}
		if accepted, err := helper.Child("", val.Elem(), m); err != nil {
			return nil, err
		} else if accepted {
			return helper.Accept("child accepted pointee"), nil
		} else {
			return helper.Reject("child rejected pointee"), nil
		}
	})
}