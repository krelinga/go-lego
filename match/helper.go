package match

import (
	"fmt"
	"reflect"
)

type Helper struct {
	Meta Meta
	Val  reflect.Value
	children []Child
}

func (h *Helper) Child(name string, val reflect.Value, m Matcher) (accepted bool, err error) {
	result, err := m.Match(val)
	if err != nil {
		return false, &FatalError{
			Meta: h.Meta,
			Val: val,
			Err: &ChildError{
				Name: name,
				Err:  err,
			},
		}
	}
	h.children = append(h.children, Child{
		Name: name,
		Result: result,
	})
	return result.Accepted, nil
}

func (h *Helper) Fatalf(format string, args ...any) error {
	return &FatalError{
		Meta: h.Meta,
		Val: h.Val,
		Err: fmt.Errorf(format, args...),
	}
}

func (h *Helper) CheckValid() error {
	if !h.Val.IsValid() {
		return h.Fatalf("value is invalid")
	}
	return nil
}

func (h *Helper) Accept(why string) *Result {
	return &Result{
		Meta: h.Meta,
		Val:  h.Val,
		Accepted: true,
		Why: why,
		Children: h.children,
	}
}

func (h *Helper) Reject(why string) *Result {
	return &Result{
		Meta: h.Meta,
		Val:  h.Val,
		Accepted: false,
		Why: why,
		Children: h.children,
	}
}