package match

type Matcher interface {
	Match(got any) error
}

type GreaterThanMatcher struct {
	threshold any
	using     any
}

func (m GreaterThanMatcher) Match(got any) error {
	return nil // TODO
}

func (m GreaterThanMatcher) Using(using any) GreaterThanMatcher {
	m.using = using
	return m
}

func GreaterThan(threshold any) GreaterThanMatcher {
	return GreaterThanMatcher{threshold: threshold}
}

type EqualMatcher struct {
	expected any
	using    any
}

func (m EqualMatcher) Match(got any) error {
	return nil // TODO
}

func (m EqualMatcher) Using(using any) EqualMatcher {
	m.using = using
	return m
}

func Equal(expected any) EqualMatcher {
	return EqualMatcher{expected: expected}
}