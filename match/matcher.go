package match

type Matcher interface {
	Match(got any) error
}

type Equal struct {
	X any
	Use any
}

func (e Equal) Match(got any) error {
	return nil  // TODO
}

type GreaterThan struct {
	X any
	Use any
}

func (g GreaterThan) Match(got any) error {
	return nil  // TODO
}