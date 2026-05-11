package match

import (
	"fmt"
	"reflect"
)

func Not(m Matcher) Matcher {
	meta := MetaHere()
	return FuncMatcher(func(val reflect.Value) (*Result, error) {
		helper := &Helper{
			Meta: meta,
			Val: val,
		}
		childAccept, err := helper.Child("", val, m)
		if err != nil {
			return nil, err
		}
		if childAccept {
			return helper.Reject("child matcher accepted"), nil
		}
		return helper.Accept("child matcher rejected"), nil
	})
}

func AllOf(matchers ...Matcher) Matcher {
	meta := MetaHere()
	return FuncMatcher(func(val reflect.Value) (*Result, error) {
		helper := &Helper{
			Meta: meta,
			Val: val,
		}
		accept := true
		for i, m := range matchers {
			childAccept, err := helper.Child(fmt.Sprintf("matcher #%d", i), val, m)
			if err != nil {
				return nil, err
			}
			if !childAccept {
				accept = false
			}
		}
		if accept {
			return helper.Accept("all child matchers accepted"), nil
		}
		return helper.Reject("at least one child matcher rejected"), nil
	})
}