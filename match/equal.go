package match

import (
	"math"
)

type EqualMatcher[T comparable] struct {
	meta    Meta
	want    T
	wantFmt Fmt[T]
	gotFmt  Fmt[any]
	valid   bool
}

func (m *EqualMatcher[T]) Match(val any) (*Result, error) {
	h := &Helper{
		Meta: m.meta,
		Val:  val,
	}
	if !m.valid {
		return nil, h.Fatalf("EqualMatcher must be created with Equal")
	}
	h.Context("expected", m.want)
	asT, err := As[T](h, val)
	if err != nil {
		return nil, err
	}
	if asT != m.want {
		return h.Reject("got != expected"), nil
	}
	return h.Accept("got == expected"), nil
}

func (m *EqualMatcher[T]) WantFmt(t Fmt[T]) *EqualMatcher[T] {
	m.wantFmt = t
	return m
}

func (m *EqualMatcher[T]) GotFmt(t Fmt[any]) *EqualMatcher[T] {
	m.gotFmt = t
	return m
}

func Equal[T comparable](want T) *EqualMatcher[T] {
	return &EqualMatcher[T]{
		meta:  MetaHere(),
		want:  want,
		valid: true,
	}
}

type FloatingPoint interface {
	~float32 | ~float64
}

type EqualApproxMatcher[T FloatingPoint] struct {
	meta      Meta
	want      T
	approx    T
	wantFmt   Fmt[T]
	gotFmt    Fmt[any]
	approxFmt Fmt[T]
	valid     bool
}

func (m *EqualApproxMatcher[T]) Match(val any) (*Result, error) {
	h := &Helper{
		Meta: m.meta,
		Val:  val,
	}
	if !m.valid {
		return nil, h.Fatalf("EqualApproxMatcher must be created with EqualApprox")
	}
	h.Context("expected", m.want)
	h.Context("tolerance", m.approx)
	asT, err := As[T](h, val)
	if err != nil {
		return nil, err
	}
	diff := asT - m.want
	if math.Abs(float64(diff)) > math.Abs(float64(m.approx)) {
		return h.Reject("got is not approximately equal to expected"), nil
	}
	return h.Accept("got is approximately equal to expected"), nil
}

func (m *EqualApproxMatcher[T]) WantFmt(t Fmt[T]) *EqualApproxMatcher[T] {
	m.wantFmt = t
	return m
}

func (m *EqualApproxMatcher[T]) ApproxFmt(t Fmt[T]) *EqualApproxMatcher[T] {
	m.approxFmt = t
	return m
}

func (m *EqualApproxMatcher[T]) GotFmt(t Fmt[any]) *EqualApproxMatcher[T] {
	m.gotFmt = t
	return m
}

func EqualApprox[T FloatingPoint](want T, approx T) *EqualApproxMatcher[T] {
	return &EqualApproxMatcher[T]{
		meta:   MetaHere(),
		want:   want,
		approx: approx,
		valid:  true,
	}
}

type EqualFuncMatcher[T, TT any] struct {
	meta      Meta
	want      T
	equalFunc func(TT, TT) bool
	wantFmt   Fmt[T]
	gotFmt    Fmt[any]
	valid     bool
}

func (m *EqualFuncMatcher[T, TT]) Match(val any) (*Result, error) {
	h := &Helper{
		Meta: m.meta,
		Val:  val,
	}
	if !m.valid {
		return nil, h.Fatalf("EqualFuncMatcher must be created with EqualFunc")
	}
	h.Context("expected", m.want)
	wantTT, ok := any(m.want).(TT)
	if !ok {
		return nil, h.Fatalf("expected value of type %T cannot be used with equal function that expects type %T", m.want, wantTT)
	}
	gotTT, err := As[TT](h, val)
	if err != nil {
		return nil, err
	}
	if !m.equalFunc(wantTT, gotTT) {
		return h.Reject("got != expected according to equal function"), nil
	}
	return h.Accept("got == expected according to equal function"), nil
}

func (m *EqualFuncMatcher[T, TT]) WantFmt(t Fmt[T]) *EqualFuncMatcher[T, TT] {
	m.wantFmt = t
	return m
}

func (m *EqualFuncMatcher[T, TT]) GotFmt(t Fmt[any]) *EqualFuncMatcher[T, TT] {
	m.gotFmt = t
	return m
}

func EqualFunc[T, TT any](want T, equalFunc func(TT, TT) bool) *EqualFuncMatcher[T, TT] {
	return &EqualFuncMatcher[T, TT]{
		meta:      MetaHere(),
		want:      want,
		equalFunc: equalFunc,
		valid:     true,
	}
}
