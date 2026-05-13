package match

type not struct {
	Matcher Matcher
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
