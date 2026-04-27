package exam

import "testing"

// FakeLog holds the variadic arguments from a single Error or Fatal call on FakeT.
type FakeLog []any

// FakeT is a test double for T. It captures Error and Fatal calls so that
// exam-based assertions can be unit-tested without a real *testing.T.
type FakeT struct {
	// Errors accumulates the arguments from each Error call.
	Errors []FakeLog
	// Fatals accumulates the arguments from each Fatal call.
	Fatals []FakeLog
}

func (t *FakeT) Helper() {}

func (t *FakeT) Error(args ...any) {
	t.Errors = append(t.Errors, args)
}

func (t *FakeT) Fatal(args ...any) {
	t.Fatals = append(t.Fatals, args)
}

func (t *FakeT) Run(name string, f func(t *testing.T)) bool {
	panic("FakeT does not support subtests")
}
