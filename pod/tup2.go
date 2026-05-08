package pod

// FixedTup2 is a read-only interface for a two-tuple (pair) of values of types A and B. It provides
// methods to access the values but does not allow modification. If either A or B is a pointer then
// the values can be modified through the pointer, but the references to A and B themselves cannot
// be changed.
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

// WrapTup2A takes a FixedTup2 and a function that transforms the A value, and returns a new
// FixedTup2 with the transformed A value and the original B value. The resulting FixedTup2 will
// keep a reference to the original FixedTup2, so if the original value is modified the changes will
// be reflected in the wrapped FixedTup2.
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

// WrapTup2B takes a FixedTup2 and a function that transforms the B value, and returns a new
// FixedTup2 with the original A value and the transformed B value. The resulting FixedTup2 will
// keep a reference to the original FixedTup2, so if the original value is modified the changes will
// be reflected in the wrapped FixedTup2.
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

// ViewTup2A takes a FixedTup2 who's A value implements the Viewer interface and returns a new
// FixedTup2 with the A value transformed using the View method of the Viewer interface, and the
// original B value. The resulting FixedTup2 will keep a reference to the original FixedTup2, so if
// the original value is modified the changes will be reflected in the wrapped FixedTup2.
func ViewTup2A[A Viewer[AA], B, AA any](f FixedTup2[A, B]) FixedTup2[AA, B] {
	return WrapTup2A(f, A.View)
}

// ViewTup2B takes a FixedTup2 who's B value implements the Viewer interface and returns a new
// FixedTup2 with the original A value and the B value transformed using the View method of the
// Viewer interface. The resulting FixedTup2 will keep a reference to the original FixedTup2, so if
// the original value is modified the changes will be reflected in the wrapped FixedTup2.
func ViewTup2B[A, B Viewer[BB], BB any](f FixedTup2[A, B]) FixedTup2[A, BB] {
	return WrapTup2B(f, B.View)
}
