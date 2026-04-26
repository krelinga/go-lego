package exam_test

import (
	"fmt"
	"testing"

	"github.com/krelinga/go-lego/exam"
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

func TestMust(t *testing.T) {
	t.Run("NoFailure", func(t *testing.T) {
		fakeT := &FakeT{T: t}
		exam.Must(fakeT, nil)
		if fakeT.ErrorBody != "" {
			t.Errorf("Must unexpectedly called Error: %q", fakeT.ErrorBody)
		}
		if fakeT.FatalBody != "" {
			t.Errorf("Must unexpectedly called Fatal: %q", fakeT.FatalBody)
		}
	})

	cases := []struct {
		Name        string
		FatalGolden exam.Golden
		Failure     *exam.Failure
	}{
		{
			Name: "SingleArgFailure",
			FatalGolden: exam.GoldenHere(`
`),
			Failure: exam.NewFailure1("Foo", "bar"),
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			fakeT := &FakeT{T: t}
			exam.Must(fakeT, c.Failure)
			if fakeT.ErrorBody != "" {
				t.Errorf("Must unexpectedly called Error: %q", fakeT.ErrorBody)
			}
			exam.Try(t, exam.GoldenEqual(fakeT.FatalBody, c.FatalGolden))
		})
	}
}
