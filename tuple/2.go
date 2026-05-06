package tuple

// View2 is a view of a two-tuple (pair) of values of types A and B.
type View2[A, B any] interface {
	GetA() A
	GetB() B
}

// T2 is a two-tuple (pair) of values of types A and B.
type T2[A, B any] struct {
	A A
	B B
}

func (t T2[A, B]) GetA() A {
	return t.A
}

func (t T2[A, B]) GetB() B {
	return t.B
}

// New2 creates a new T2 with the given values.
func New2[A, B any](a A, b B) T2[A, B] {
	return T2[A, B]{A: a, B: b}
}
