package match

import "reflect"

type Nil struct{}

func (n Nil) Match(val any) (*Result, error) {
	h := &Helper{
		Meta: MetaHere(),
		Val:  val,
	}
	const acceptStr = "value is nil"

	if val == nil {
		return h.Accept(acceptStr), nil
	}
	rVal := reflect.ValueOf(val)
	switch rVal.Kind() {
	case reflect.Interface, reflect.Pointer, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		if rVal.IsNil() {
			return h.Accept(acceptStr), nil
		} else {
			return h.Reject("value is not nil"), nil
		}
	}
	return nil, h.Fatalf("Nil matcher only accepts nil-able types, but got type %s", rVal.Type())
}
