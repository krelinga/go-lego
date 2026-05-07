package tuple

import "github.com/krelinga/go-libs/view"

// Fixed2 is a view of a two-tuple (pair) of values of types A and B.
type Fixed2[A, B any] interface {
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

func Wrap2A[A, B, AA any](f Fixed2[A, B], w func(A) AA) Fixed2[AA, B] {
	return wrap2A[A, B, AA]{f: f, w: w}
}

type wrap2A[A, B, AA any] struct {
	f Fixed2[A, B]
	w func(A) AA
}

func (w wrap2A[A, B, AA]) GetA() AA {
	return w.w(w.f.GetA())
}

func (w wrap2A[A, B, AA]) GetB() B {
	return w.f.GetB()
}

func Wrap2B[A, B, BB any](f Fixed2[A, B], w func(B) BB) Fixed2[A, BB] {
	return wrap2B[A, B, BB]{f: f, w: w}
}

type wrap2B[A, B, BB any] struct {
	f Fixed2[A, B]
	w func(B) BB
}

func (w wrap2B[A, B, BB]) GetA() A {
	return w.f.GetA()
}

func (w wrap2B[A, B, BB]) GetB() BB {
	return w.w(w.f.GetB())
}

func View2A[A view.Viewer[AA], B, AA any](f Fixed2[A, B]) Fixed2[AA, B] {
	return Wrap2A(f, A.View)
}

func View2B[A, B view.Viewer[BB], BB any](f Fixed2[A, B]) Fixed2[A, BB] {
	return Wrap2B(f, B.View)
}
