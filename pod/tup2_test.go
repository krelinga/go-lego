package pod_test

import (
	"testing"

	"github.com/krelinga/go-libs/exam"
	"github.com/krelinga/go-libs/pod"
)

func TestTup2ImplementsFixedTup2(t *testing.T) {
	exam.Try(t, exam.Implements[pod.Tup2[int, string], pod.FixedTup2[int, string]]())
}

func TestTup2Fields(t *testing.T) {
	cases := []struct {
		Name  string
		Loc   exam.Loc
		Input pod.Tup2[int, string]
		WantA int
		WantB string
	}{
		{
			Name:  "zero values",
			Loc:   exam.Here(),
			Input: pod.Tup2[int, string]{},
			WantA: 0,
			WantB: "",
		},
		{
			Name:  "non-zero values",
			Loc:   exam.Here(),
			Input: pod.Tup2[int, string]{A: 42, B: "hello"},
			WantA: 42,
			WantB: "hello",
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			exam.Try(t, exam.Equal(c.Input.A, c.WantA))
			exam.Try(t, exam.Equal(c.Input.B, c.WantB))
		})
	}
}

func TestTup2GetA(t *testing.T) {
	cases := []struct {
		Name  string
		Loc   exam.Loc
		Input pod.Tup2[int, string]
		Want  int
	}{
		{
			Name:  "zero value",
			Loc:   exam.Here(),
			Input: pod.Tup2[int, string]{},
			Want:  0,
		},
		{
			Name:  "non-zero value",
			Loc:   exam.Here(),
			Input: pod.Tup2[int, string]{A: 7, B: "x"},
			Want:  7,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			exam.Try(t, exam.Equal(c.Input.GetA(), c.Want))
		})
	}
}

func TestTup2GetB(t *testing.T) {
	cases := []struct {
		Name  string
		Loc   exam.Loc
		Input pod.Tup2[int, string]
		Want  string
	}{
		{
			Name:  "zero value",
			Loc:   exam.Here(),
			Input: pod.Tup2[int, string]{},
			Want:  "",
		},
		{
			Name:  "non-zero value",
			Loc:   exam.Here(),
			Input: pod.Tup2[int, string]{A: 7, B: "world"},
			Want:  "world",
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			exam.Try(t, exam.Equal(c.Input.GetB(), c.Want))
		})
	}
}

func TestNewTup2(t *testing.T) {
	cases := []struct {
		Name  string
		Loc   exam.Loc
		A     int
		B     string
		WantA int
		WantB string
	}{
		{
			Name:  "zero values",
			Loc:   exam.Here(),
			A:     0,
			B:     "",
			WantA: 0,
			WantB: "",
		},
		{
			Name:  "non-zero values",
			Loc:   exam.Here(),
			A:     99,
			B:     "go-lego",
			WantA: 99,
			WantB: "go-lego",
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			got := pod.NewTup2(c.A, c.B)
			exam.Try(t, exam.Equal(got.A, c.WantA))
			exam.Try(t, exam.Equal(got.B, c.WantB))
			exam.Try(t, exam.Equal(got.GetA(), c.WantA))
			exam.Try(t, exam.Equal(got.GetB(), c.WantB))
		})
	}
}
