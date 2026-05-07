package pod

// FixedTup2 is a view of a two-tuple (pair) of values of types A and B.
type FixedTup2[A, B any] interface {
	GetA() A
	GetB() B
}

// Tup2 is a two-tuple (pair) of values of types A and B.
type Tup2[A, B any] struct {
	A A
	B B
}

func (t Tup2[A, B]) GetA() A {
	return t.A
}

func (t Tup2[A, B]) GetB() B {
	return t.B
}

// NewTup2 creates a new Tup2 with the given values.
func NewTup2[A, B any](a A, b B) Tup2[A, B] {
	return Tup2[A, B]{A: a, B: b}
}

func WrapTup2A[A, B, AA any](f FixedTup2[A, B], w func(A) AA) FixedTup2[AA, B] {
	return wrapTup2A[A, B, AA]{f: f, w: w}
}

type wrapTup2A[A, B, AA any] struct {
	f FixedTup2[A, B]
	w func(A) AA
}

func (w wrapTup2A[A, B, AA]) GetA() AA {
	return w.w(w.f.GetA())
}

func (w wrapTup2A[A, B, AA]) GetB() B {
	return w.f.GetB()
}

func WrapTup2B[A, B, BB any](f FixedTup2[A, B], w func(B) BB) FixedTup2[A, BB] {
	return wrapTup2B[A, B, BB]{f: f, w: w}
}

type wrapTup2B[A, B, BB any] struct {
	f FixedTup2[A, B]
	w func(B) BB
}

func (w wrapTup2B[A, B, BB]) GetA() A {
	return w.f.GetA()
}

func (w wrapTup2B[A, B, BB]) GetB() BB {
	return w.w(w.f.GetB())
}

func ViewTup2A[A Viewer[AA], B, AA any](f FixedTup2[A, B]) FixedTup2[AA, B] {
	return WrapTup2A(f, A.View)
}

func ViewTup2B[A, B Viewer[BB], BB any](f FixedTup2[A, B]) FixedTup2[A, BB] {
	return WrapTup2B(f, B.View)
}
