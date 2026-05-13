package match

import "fmt"

type not struct {
	Matcher Matcher
}

func (n *not) Validate() error {
	if n.Matcher == nil {
		return fmt.Errorf("not matcher must have a non-nil Matcher field")
	}
	return n.Matcher.Validate()
}

func (n *not) Match(val any) (*Result, error) {
	h := &Helper{
		Meta: MetaHere(),
		Val:  val,
	}
	childAccepted, err := h.Child("", val, n.Matcher)
	if err != nil {
		return nil, err
	}
	if childAccepted {
		return h.Reject("child matcher accepted"), nil
	} else {
		return h.Accept("child matcher rejected"), nil
	}
}

func Not(m Matcher) Matcher {
	return &not{
		Matcher: m,
	}
}
