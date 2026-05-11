package match

import (
	"fmt"
	"reflect"
)

type Helper struct {
	Meta     Meta
	Val      reflect.Value
	children []*Child
	context  []*Context
}

func (h *Helper) Context(name string, val reflect.Value) {
	h.context = append(h.context, &Context{
		Name: name,
		Val:  val,
	})
}

func (h *Helper) Child(name string, val reflect.Value, m Matcher) (accepted bool, err error) {
	result, err := m.Match(val)
	if err != nil {
		return false, &FatalError{
			Meta: h.Meta,
			Val:  val,
			Err: &ChildError{
				Name: name,
				Err:  err,
			},
		}
	}
	h.children = append(h.children, &Child{
		Name:   name,
		Result: result,
	})
	return result.Accepted, nil
}

func (h *Helper) Fatalf(format string, args ...any) error {
	return h.Fatal(fmt.Errorf(format, args...))
}

func (h *Helper) Fatal(err error) error {
	return &FatalError{
		Meta: h.Meta,
		Val:  h.Val,
		Err:  err,
	}
}

func (h *Helper) CheckValid() error {
	if !h.Val.IsValid() {
		return h.Fatalf("value is invalid")
	}
	return nil
}

func (h *Helper) result(accepted bool, why string) *Result {
	return &Result{
		Meta:     h.Meta,
		Val:      h.Val,
		Accepted: accepted,
		Why:      why,
		Children: h.children,
		Context:  h.context,
	}
}

func (h *Helper) Accept(why string) *Result {
	return h.result(true, why)
}

func (h *Helper) Reject(why string) *Result {
	return h.result(false, why)
}
