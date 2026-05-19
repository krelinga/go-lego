package match

type NotMatcher struct {
	meta    Meta
	matcher Matcher
	fmtGot  Fmt[any]
	valid   bool
}

func (n *NotMatcher) Match(val any) (*Result, error) {
	h := &Helper{
		Meta: n.meta,
		Val:  val,
	}
	if !n.valid {
		return nil, h.Fatalf("Not matcher must be created with Not()")
	}
	childAccepted, err := h.Child("", val, n.matcher)
	if err != nil {
		return nil, err
	}
	if childAccepted {
		return h.Reject("child matcher accepted"), nil
	} else {
		return h.Accept("child matcher rejected"), nil
	}
}

func (n *NotMatcher) GotFmt(t Fmt[any]) *NotMatcher {
	n.fmtGot = t
	return n
}

func Not(m Matcher) Matcher {
	return &NotMatcher{
		meta:    MetaHere(),
		matcher: m,
		valid:   true,
	}
}
