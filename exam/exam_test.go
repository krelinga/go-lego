package exam_test

import (
	"fmt"
	"testing"

	"github.com/krelinga/go-libs/exam"
)

type FakeT struct {
	T                    *testing.T
	ErrorBody, FatalBody string
}

func (t *FakeT) Helper() {}

func (t *FakeT) Error(args ...any) {
	if len(args) != 1 {
		t.T.Fatalf("Error called with %d arguments, expected 1", len(args))
	} else if len(t.ErrorBody) != 0 {
		t.T.Fatalf("Error called multiple times: previous call had body %q, new call has body %q", t.ErrorBody, fmt.Sprint(args...))
	}
	t.ErrorBody = fmt.Sprint(args...)
}

func (t *FakeT) Fatal(args ...any) {
	if len(args) != 1 {
		t.T.Fatalf("Fatal called with %d arguments, expected 1", len(args))
	} else if len(t.FatalBody) != 0 {
		t.T.Fatalf("Fatal called multiple times: previous call had body %q, new call has body %q", t.FatalBody, fmt.Sprint(args...))
	}
	t.FatalBody = fmt.Sprint(args...)
}

func (t *FakeT) Run(name string, f func(t *testing.T)) bool {
	t.T.Fatal("Run should not be called on FakeT")
	return false
}

func TestTryMust(t *testing.T) {
	t.Run("NoFailureMust", func(t *testing.T) {
		fakeT := &FakeT{T: t}
		exam.Must(fakeT, nil)
		if fakeT.ErrorBody != "" {
			t.Errorf("Must unexpectedly called Error: %q", fakeT.ErrorBody)
		}
		if fakeT.FatalBody != "" {
			t.Errorf("Must unexpectedly called Fatal: %q", fakeT.FatalBody)
		}
	})
	t.Run("NoFailureTry", func(t *testing.T) {
		fakeT := &FakeT{T: t}
		exam.Try(fakeT, nil)
		if fakeT.ErrorBody != "" {
			t.Errorf("Try unexpectedly called Error: %q", fakeT.ErrorBody)
		}
		if fakeT.FatalBody != "" {
			t.Errorf("Try unexpectedly called Fatal: %q", fakeT.FatalBody)
		}
	})

	cases := []struct {
		Name       string
		MustGolden exam.Golden
		TryGolden  exam.Golden
		Failure    *exam.Failure
		Args       []any
		Fatal      bool
	}{
		{
			Name: "SingleArgFailure",
			MustGolden: exam.GoldenHere(`
FATAL: exam.Must(fakeT, c.Failure, c.Args...)
arg 0 Foo: "bar"`),
			TryGolden: exam.GoldenHere(`
exam.Try(fakeT, c.Failure, c.Args...)
arg 0 Foo: "bar"`),
			Failure: exam.NewFailure1("Foo", "bar"),
			Fatal:   false,
		},
		{
			Name: "MultiArgFailure",
			MustGolden: exam.GoldenHere(`
FATAL: exam.Must(fakeT, c.Failure, c.Args...)
arg 0 Foo: "bar"
arg 1 Baz: 42`),
			TryGolden: exam.GoldenHere(`
exam.Try(fakeT, c.Failure, c.Args...)
arg 0 Foo: "bar"
arg 1 Baz: 42`),
			Failure: exam.NewFailure2("Foo", "bar", "Baz", 42),
			Fatal:   false,
		},
		{
			Name: "WrappedFailure",
			MustGolden: exam.GoldenHere(`
FATAL: exam.Must(fakeT, c.Failure, c.Args...)
STRUCTURAL ERROR: wrapped error
arg 0 Foo: "bar"`),
			TryGolden: exam.GoldenHere(`
FATAL: exam.Try(fakeT, c.Failure, c.Args...)
STRUCTURAL ERROR: wrapped error
arg 0 Foo: "bar"`),
			Failure: exam.NewFailure1("Foo", "bar").Wrap(fmt.Errorf("wrapped error")),
			Fatal:   true,
		},
		{
			Name: "ArgWithStringIndentLinesFmt",
			MustGolden: exam.GoldenHere(`
FATAL: exam.Must(fakeT, c.Failure, c.Args...)
arg 0 Foo: line1
	line2
	line3`),
			TryGolden: exam.GoldenHere(`
exam.Try(fakeT, c.Failure, c.Args...)
arg 0 Foo: line1
	line2
	line3`),
			Failure: func() *exam.Failure {
				f := exam.NewFailure1("Foo", "line1\nline2\nline3")
				f.Args[0].Fmt = exam.StringIndentLines
				return f
			}(),
			Fatal: false,
		},
		{
			Name: "ExtrasPrintedOnePerLine",
			MustGolden: exam.GoldenHere(`
FATAL: exam.Must(fakeT, c.Failure, c.Args...)
arg 0 Foo: "bar"
baz
biff`),
			TryGolden: exam.GoldenHere(`
exam.Try(fakeT, c.Failure, c.Args...)
arg 0 Foo: "bar"
baz
biff`),
			Failure: exam.NewFailure1("Foo", "bar"),
			Args:    []any{"baz", "biff"},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.MustGolden.GetLoc(), func(t *testing.T) {
			func() {
				fakeT := &FakeT{T: t}
				exam.Must(fakeT, c.Failure, c.Args...)
				if fakeT.ErrorBody != "" {
					t.Fatalf("Must unexpectedly called Error: %q", fakeT.ErrorBody)
				}
				exam.Try(t, exam.GoldenEqual(fakeT.FatalBody, c.MustGolden))
			}()
			func() {
				fakeT := &FakeT{T: t}
				exam.Try(fakeT, c.Failure, c.Args...)
				var got string
				if c.Fatal {
					if fakeT.ErrorBody != "" {
						t.Fatalf("Try unexpectedly called Error: %q", fakeT.ErrorBody)
					}
					got = fakeT.FatalBody
				} else if !c.Fatal {
					if fakeT.FatalBody != "" {
						t.Fatalf("Try unexpectedly called Fatal: %q", fakeT.FatalBody)
					}
					got = fakeT.ErrorBody
				}
				exam.Try(t, exam.GoldenEqual(got, c.TryGolden))
			}()
		})
	}
}
