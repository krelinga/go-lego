package exam

import "testing"

type FakeLog []any

type FakeT struct {
	Errors []FakeLog
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
