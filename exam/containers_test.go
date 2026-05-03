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
