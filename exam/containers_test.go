package exam_test

import (
	"strings"
	"testing"

	"github.com/krelinga/go-lego/exam"
)

func TestMapEqual(t *testing.T) {
	cases := []struct {
		Name       string
		Loc        exam.Loc
		Map1, Map2 map[string]int
		WantEqual  bool
	}{
		{
			Name:      "Equal",
			Loc:       exam.Here(),
			Map1:      map[string]int{"a": 1, "b": 2},
			Map2:      map[string]int{"a": 1, "b": 2},
			WantEqual: true,
		},
		{
			Name:      "DifferentKeys",
			Loc:       exam.Here(),
			Map1:      map[string]int{"a": 1, "b": 2},
			Map2:      map[string]int{"a": 1, "c": 2},
			WantEqual: false,
		},
		{
			Name:      "DifferentValues",
			Loc:       exam.Here(),
			Map1:      map[string]int{"a": 1, "b": 2},
			Map2:      map[string]int{"a": 1, "b": 3},
			WantEqual: false,
		},
		{
			Name:      "NilVsEmpty",
			Loc:       exam.Here(),
			Map1:      map[string]int(nil),
			Map2:      map[string]int{},
			WantEqual: true,
		},
		{
			Name:      "DifferentLengths",
			Loc:       exam.Here(),
			Map1:      map[string]int{"a": 1, "b": 2},
			Map2:      map[string]int{"a": 1},
			WantEqual: false,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			failure := exam.MapEqual(c.Map1, c.Map2)
			if c.WantEqual {
				exam.Try(t, exam.Nil(failure))
			} else {
				exam.Must(t, exam.NotNil(failure))
			}
		})
	}
}

func TestMapEqualFunc(t *testing.T) {
	cases := []struct {
		Name       string
		Loc        exam.Loc
		Map1, Map2 map[string]string
		EqualFunc  func(string, string) bool
		WantEqual  bool
	}{
		{
			Name:      "Equal",
			Loc:       exam.Here(),
			Map1:      map[string]string{"a": "1", "b": "2"},
			Map2:      map[string]string{"a": "1", "b": "2"},
			EqualFunc: func(s1, s2 string) bool { return s1 == s2 },
			WantEqual: true,
		},
		{
			Name:      "CaseInsensitiveEqual",
			Loc:       exam.Here(),
			Map1:      map[string]string{"a": "a", "b": "B"},
			Map2:      map[string]string{"a": "A", "b": "b"},
			EqualFunc: func(s1, s2 string) bool { return strings.EqualFold(s1, s2) },
			WantEqual: true,
		},
		{
			Name:      "DifferentKeys",
			Loc:       exam.Here(),
			Map1:      map[string]string{"a": "1", "b": "2"},
			Map2:      map[string]string{"a": "1", "c": "2"},
			EqualFunc: func(s1, s2 string) bool { return s1 == s2 },
			WantEqual: false,
		},
		{
			Name:      "DifferentValues",
			Loc:       exam.Here(),
			Map1:      map[string]string{"a": "1", "b": "2"},
			Map2:      map[string]string{"a": "1", "b": "3"},
			EqualFunc: func(s1, s2 string) bool { return s1 == s2 },
			WantEqual: false,
		},
		{
			Name:      "NilVsEmpty",
			Loc:       exam.Here(),
			Map1:      map[string]string(nil),
			Map2:      map[string]string{},
			EqualFunc: func(s1, s2 string) bool { return s1 == s2 },
			WantEqual: true,
		},
		{
			Name:      "DifferentLengths",
			Loc:       exam.Here(),
			Map1:      map[string]string{"a": "1", "b": "2"},
			Map2:      map[string]string{"a": "1"},
			EqualFunc: func(s1, s2 string) bool { return s1 == s2 },
			WantEqual: false,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			failure := exam.MapEqualFunc(c.Map1, c.Map2, c.EqualFunc)
			if c.WantEqual {
				exam.Try(t, exam.Nil(failure))
			} else {
				exam.Must(t, exam.NotNil(failure))
			}
		})
	}

	t.Run("EqualAsAny", func(t *testing.T) {
		Map1 := map[string]string{"a": "1", "b": "2"}
		Map2 := map[string]string{"a": "1", "b": "2"}
		failure := exam.MapEqualFunc(Map1, Map2, func(s1, s2 any) bool {
			v1 := exam.MustCast[string](t, s1)
			v2 := exam.MustCast[string](t, s2)
			return v1 == v2
		})
		exam.Try(t, exam.Nil(failure))
	})
}

func TestSliceEqual(t *testing.T) {
	cases := []struct {
		Name           string
		Loc            exam.Loc
		Slice1, Slice2 []int
		WantEqual      bool
	}{
		{
			Name:      "Equal",
			Loc:       exam.Here(),
			Slice1:    []int{1, 2, 3},
			Slice2:    []int{1, 2, 3},
			WantEqual: true,
		},
		{
			Name:      "DifferentValues",
			Loc:       exam.Here(),
			Slice1:    []int{1, 2, 3},
			Slice2:    []int{1, 2, 99},
			WantEqual: false,
		},
		{
			Name:      "DifferentLengths",
			Loc:       exam.Here(),
			Slice1:    []int{1, 2, 3},
			Slice2:    []int{1, 2},
			WantEqual: false,
		},
		{
			Name:      "NilVsEmpty",
			Loc:       exam.Here(),
			Slice1:    []int(nil),
			Slice2:    []int{},
			WantEqual: true,
		},
		{
			Name:      "BothNil",
			Loc:       exam.Here(),
			Slice1:    nil,
			Slice2:    nil,
			WantEqual: true,
		},
		{
			Name:      "DifferentOrder",
			Loc:       exam.Here(),
			Slice1:    []int{1, 2, 3},
			Slice2:    []int{3, 2, 1},
			WantEqual: false,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			failure := exam.SliceEqual(c.Slice1, c.Slice2)
			if c.WantEqual {
				exam.Try(t, exam.Nil(failure))
			} else {
				exam.Must(t, exam.NotNil(failure))
			}
		})
	}
}

func TestSliceEqualFunc(t *testing.T) {
	cases := []struct {
		Name           string
		Loc            exam.Loc
		Slice1, Slice2 []string
		EqualFunc      func(string, string) bool
		WantEqual      bool
	}{
		{
			Name:      "Equal",
			Loc:       exam.Here(),
			Slice1:    []string{"a", "b"},
			Slice2:    []string{"a", "b"},
			EqualFunc: func(s1, s2 string) bool { return s1 == s2 },
			WantEqual: true,
		},
		{
			Name:      "CaseInsensitiveEqual",
			Loc:       exam.Here(),
			Slice1:    []string{"Hello", "World"},
			Slice2:    []string{"hello", "WORLD"},
			EqualFunc: func(s1, s2 string) bool { return strings.EqualFold(s1, s2) },
			WantEqual: true,
		},
		{
			Name:      "DifferentValues",
			Loc:       exam.Here(),
			Slice1:    []string{"a", "b"},
			Slice2:    []string{"a", "z"},
			EqualFunc: func(s1, s2 string) bool { return s1 == s2 },
			WantEqual: false,
		},
		{
			Name:      "DifferentLengths",
			Loc:       exam.Here(),
			Slice1:    []string{"a", "b"},
			Slice2:    []string{"a"},
			EqualFunc: func(s1, s2 string) bool { return s1 == s2 },
			WantEqual: false,
		},
		{
			Name:      "NilVsEmpty",
			Loc:       exam.Here(),
			Slice1:    []string(nil),
			Slice2:    []string{},
			EqualFunc: func(s1, s2 string) bool { return s1 == s2 },
			WantEqual: true,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			failure := exam.SliceEqualFunc(c.Slice1, c.Slice2, c.EqualFunc)
			if c.WantEqual {
				exam.Try(t, exam.Nil(failure))
			} else {
				exam.Must(t, exam.NotNil(failure))
			}
		})
	}

	t.Run("EqualAsAny", func(t *testing.T) {
		slice1 := []string{"a", "b"}
		slice2 := []string{"a", "b"}
		failure := exam.SliceEqualFunc(slice1, slice2, func(s1, s2 any) bool {
			v1 := exam.MustCast[string](t, s1)
			v2 := exam.MustCast[string](t, s2)
			return v1 == v2
		})
		exam.Try(t, exam.Nil(failure))
	})
}
