package match

import "reflect"

type NilMatcher struct {
	meta   Meta
	gotFmt Fmt[any]
	valid  bool
}

func (m *NilMatcher) Match(val any) (*Result, error) {
	h := &Helper{
		Meta: m.meta,
		Val:  val,
	}
	if !m.valid {
		return nil, h.Fatalf("Nil matcher must be created with Nil()")
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
	return nil, h.Fatalf("XNil matcher only accepts nil-able types, but got type %s", rVal.Type())
}

func (m *NilMatcher) GotFmt(t Fmt[any]) *NilMatcher {
	m.gotFmt = t
	return m
}

func Nil() *NilMatcher {
	return &NilMatcher{
		meta:  MetaHere(),
		valid: true,
	}
}
