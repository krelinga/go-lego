package pod_test

import (
	"maps"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/pod"
)

func TestDictImplements(t *testing.T) {
	exam.Try(t, exam.Implements[*pod.Map[int, string], pod.Dict[int, string]]())
}

func TestAsDict(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Input    map[string]int
		WantLen  int
		WantData map[string]int
	}{
		{
			Name:     "nil map",
			Loc:      exam.Here(),
			Input:    nil,
			WantLen:  0,
			WantData: map[string]int{},
		},
		{
			Name:     "empty map",
			Loc:      exam.Here(),
			Input:    map[string]int{},
			WantLen:  0,
			WantData: map[string]int{},
		},
		{
			Name:     "non-empty map",
			Loc:      exam.Here(),
			Input:    map[string]int{"a": 1, "b": 2},
			WantLen:  2,
			WantData: map[string]int{"a": 1, "b": 2},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			d := pod.AsDict(c.Input)
			exam.Try(t, exam.Equal(d.Len(), c.WantLen))
			for k, wantV := range c.WantData {
				gotV, ok := d.Get(k)
				exam.Try(t, exam.True(ok))
				exam.Try(t, exam.Equal(gotV, wantV))
			}
			_, ok := d.Get("missing")
			exam.Try(t, exam.Equal(ok, false))
		})
	}
	t.Run("reflects underlying map mutations", func(t *testing.T) {
		m := map[string]int{"a": 1}
		d := pod.AsDict(m)
		exam.Try(t, exam.Equal(d.Len(), 1))
		m["b"] = 2
		exam.Try(t, exam.Equal(d.Len(), 2))
		v, ok := d.Get("b")
		exam.Try(t, exam.True(ok))
		exam.Try(t, exam.Equal(v, 2))
	})
}

func TestNewMap(t *testing.T) {
	m := pod.NewMap[string, int]()
	exam.Try(t, exam.Equal(m.Len(), 0))
	m.Put("a", 1)
	v, ok := m.Get("a")
	exam.Try(t, exam.True(ok))
	exam.Try(t, exam.Equal(v, 1))
}

func TestNewMapHint(t *testing.T) {
	m := pod.NewMapHint[string, int](10)
	exam.Try(t, exam.Equal(m.Len(), 0))
	m.Put("a", 1)
	v, ok := m.Get("a")
	exam.Try(t, exam.True(ok))
	exam.Try(t, exam.Equal(v, 1))
}

func TestNewMapOf(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Source   *pod.Map[string, int]
		WantData map[string]int
	}{
		{
			Name:     "empty source",
			Loc:      exam.Here(),
			Source:   pod.NewMap[string, int](),
			WantData: map[string]int{},
		},
		{
			Name:     "non-empty source",
			Loc:      exam.Here(),
			Source:   &pod.Map[string, int]{"a": 1, "b": 2},
			WantData: map[string]int{"a": 1, "b": 2},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			m := pod.NewMapOf(pod.KeyValsOf(c.Source))
			exam.Try(t, exam.MapEqual(map[string]int(*m), c.WantData))
		})
	}
}

func TestMapLen(t *testing.T) {
	cases := []struct {
		Name    string
		Loc     exam.Loc
		Map     *pod.Map[string, int]
		WantLen int
	}{
		{
			Name:    "empty map",
			Loc:     exam.Here(),
			Map:     pod.NewMap[string, int](),
			WantLen: 0,
		},
		{
			Name:    "one element",
			Loc:     exam.Here(),
			Map:     &pod.Map[string, int]{"a": 1},
			WantLen: 1,
		},
		{
			Name:    "three elements",
			Loc:     exam.Here(),
			Map:     &pod.Map[string, int]{"a": 1, "b": 2, "c": 3},
			WantLen: 3,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			exam.Try(t, exam.Equal(c.Map.Len(), c.WantLen))
		})
	}
}

func TestMapGet(t *testing.T) {
	cases := []struct {
		Name      string
		Loc       exam.Loc
		Map       *pod.Map[string, int]
		Key       string
		WantVal   int
		WantFound bool
	}{
		{
			Name:      "existing key",
			Loc:       exam.Here(),
			Map:       &pod.Map[string, int]{"a": 1, "b": 2},
			Key:       "a",
			WantVal:   1,
			WantFound: true,
		},
		{
			Name:      "missing key",
			Loc:       exam.Here(),
			Map:       &pod.Map[string, int]{"a": 1},
			Key:       "z",
			WantVal:   0,
			WantFound: false,
		},
		{
			Name:      "empty map",
			Loc:       exam.Here(),
			Map:       pod.NewMap[string, int](),
			Key:       "a",
			WantVal:   0,
			WantFound: false,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			gotVal, gotFound := c.Map.Get(c.Key)
			exam.Try(t, exam.Equal(gotFound, c.WantFound))
			exam.Try(t, exam.Equal(gotVal, c.WantVal))
		})
	}
}

func TestMapKeyVals(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Map      *pod.Map[string, int]
		WantData map[string]int
	}{
		{
			Name:     "empty map",
			Loc:      exam.Here(),
			Map:      pod.NewMap[string, int](),
			WantData: map[string]int{},
		},
		{
			Name:     "non-empty map",
			Loc:      exam.Here(),
			Map:      &pod.Map[string, int]{"a": 1, "b": 2},
			WantData: map[string]int{"a": 1, "b": 2},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			got := maps.Collect(c.Map.KeyVals())
			exam.Try(t, exam.MapEqual(got, c.WantData))
		})
	}
}

func TestMapKeys(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Map      *pod.Map[string, int]
		WantKeys []string
	}{
		{
			Name:     "empty map",
			Loc:      exam.Here(),
			Map:      pod.NewMap[string, int](),
			WantKeys: nil,
		},
		{
			Name:     "non-empty map",
			Loc:      exam.Here(),
			Map:      &pod.Map[string, int]{"a": 1, "b": 2, "c": 3},
			WantKeys: []string{"a", "b", "c"},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			got := slices.Collect(c.Map.Keys())
			slices.Sort(got)
			exam.Try(t, exam.SliceEqual(got, c.WantKeys))
		})
	}
}

func TestMapVals(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Map      *pod.Map[string, int]
		WantVals []int
	}{
		{
			Name:     "empty map",
			Loc:      exam.Here(),
			Map:      pod.NewMap[string, int](),
			WantVals: nil,
		},
		{
			Name:     "non-empty map",
			Loc:      exam.Here(),
			Map:      &pod.Map[string, int]{"a": 1, "b": 2, "c": 3},
			WantVals: []int{1, 2, 3},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			got := slices.Collect(c.Map.Vals())
			slices.Sort(got)
			exam.Try(t, exam.SliceEqual(got, c.WantVals))
		})
	}
}

func TestMapPut(t *testing.T) {
	t.Run("add to empty map", func(t *testing.T) {
		m := pod.NewMap[string, int]()
		m.Put("a", 1)
		exam.Try(t, exam.Equal(m.Len(), 1))
		v, ok := m.Get("a")
		exam.Try(t, exam.True(ok))
		exam.Try(t, exam.Equal(v, 1))
	})
	t.Run("overwrite existing key", func(t *testing.T) {
		m := &pod.Map[string, int]{"a": 1}
		m.Put("a", 99)
		exam.Try(t, exam.Equal(m.Len(), 1))
		v, ok := m.Get("a")
		exam.Try(t, exam.True(ok))
		exam.Try(t, exam.Equal(v, 99))
	})
	t.Run("put into nil-value map", func(t *testing.T) {
		m := new(pod.Map[string, int]) // *m is nil (zero value of map type)
		m.Put("a", 1)
		exam.Try(t, exam.Equal(m.Len(), 1))
		v, ok := m.Get("a")
		exam.Try(t, exam.True(ok))
		exam.Try(t, exam.Equal(v, 1))
	})
}

func TestMapClear(t *testing.T) {
	cases := []struct {
		Name string
		Loc  exam.Loc
		Map  *pod.Map[string, int]
	}{
		{
			Name: "clear non-empty map",
			Loc:  exam.Here(),
			Map:  &pod.Map[string, int]{"a": 1, "b": 2},
		},
		{
			Name: "clear already empty map",
			Loc:  exam.Here(),
			Map:  pod.NewMap[string, int](),
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			c.Map.Clear()
			exam.Try(t, exam.Equal(c.Map.Len(), 0))
		})
	}
}

func TestMapReserve(t *testing.T) {
	t.Run("reserve on nil-value map allows subsequent Put", func(t *testing.T) {
		m := new(pod.Map[string, int])
		m.Reserve(5)
		m.Put("a", 1)
		exam.Try(t, exam.Equal(m.Len(), 1))
		v, ok := m.Get("a")
		exam.Try(t, exam.True(ok))
		exam.Try(t, exam.Equal(v, 1))
	})
	t.Run("reserve on initialized map is no-op", func(t *testing.T) {
		m := &pod.Map[string, int]{"a": 1}
		m.Reserve(100)
		exam.Try(t, exam.Equal(m.Len(), 1))
		v, ok := m.Get("a")
		exam.Try(t, exam.True(ok))
		exam.Try(t, exam.Equal(v, 1))
	})
}

func TestMapDel(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Map      *pod.Map[string, int]
		DelKey   string
		WantLen  int
		WantData map[string]int
	}{
		{
			Name:     "delete existing key",
			Loc:      exam.Here(),
			Map:      &pod.Map[string, int]{"a": 1, "b": 2},
			DelKey:   "a",
			WantLen:  1,
			WantData: map[string]int{"b": 2},
		},
		{
			Name:     "delete missing key is a no-op",
			Loc:      exam.Here(),
			Map:      &pod.Map[string, int]{"a": 1},
			DelKey:   "z",
			WantLen:  1,
			WantData: map[string]int{"a": 1},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			c.Map.Del(c.DelKey)
			exam.Try(t, exam.Equal(c.Map.Len(), c.WantLen))
			exam.Try(t, exam.MapEqual(map[string]int(*c.Map), c.WantData))
		})
	}
}

func TestWrapDictVals(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Input    *pod.Map[string, int]
		WantData map[string]string
	}{
		{
			Name:     "empty map",
			Loc:      exam.Here(),
			Input:    pod.NewMap[string, int](),
			WantData: map[string]string{},
		},
		{
			Name:     "non-empty map",
			Loc:      exam.Here(),
			Input:    &pod.Map[string, int]{"a": 1, "b": 2},
			WantData: map[string]string{"a": "1", "b": "2"},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			wrapped := pod.WrapDictVals(c.Input, strconv.Itoa)
			exam.Try(t, exam.Equal(wrapped.Len(), len(c.WantData)))

			for k, wantV := range c.WantData {
				gotV, ok := wrapped.Get(k)
				exam.Try(t, exam.True(ok))
				exam.Try(t, exam.Equal(gotV, wantV))
			}
			_, ok := wrapped.Get("missing")
			exam.Try(t, exam.Equal(ok, false))

			gotKV := maps.Collect(wrapped.KeyVals())
			exam.Try(t, exam.MapEqual(gotKV, c.WantData))

			gotKeys := slices.Collect(wrapped.Keys())
			slices.Sort(gotKeys)
			wantKeys := slices.Collect(maps.Keys(c.WantData))
			slices.Sort(wantKeys)
			exam.Try(t, exam.SliceEqual(gotKeys, wantKeys))

			gotVals := slices.Collect(wrapped.Vals())
			slices.Sort(gotVals)
			wantVals := slices.Collect(maps.Values(c.WantData))
			slices.Sort(wantVals)
			exam.Try(t, exam.SliceEqual(gotVals, wantVals))
		})
	}
}

func TestWrapDictKeys(t *testing.T) {
	wrap := strconv.Itoa
	unwrap := func(s string) int {
		v, _ := strconv.Atoi(s)
		return v
	}

	cases := []struct {
		Name     string
		Loc      exam.Loc
		Input    *pod.Map[int, string]
		WantData map[string]string
	}{
		{
			Name:     "empty map",
			Loc:      exam.Here(),
			Input:    pod.NewMap[int, string](),
			WantData: map[string]string{},
		},
		{
			Name:     "non-empty map",
			Loc:      exam.Here(),
			Input:    &pod.Map[int, string]{1: "one", 2: "two"},
			WantData: map[string]string{"1": "one", "2": "two"},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			wrapped := pod.WrapDictKeys(c.Input, wrap, unwrap)
			exam.Try(t, exam.Equal(wrapped.Len(), len(c.WantData)))

			for k, wantV := range c.WantData {
				gotV, ok := wrapped.Get(k)
				exam.Try(t, exam.True(ok))
				exam.Try(t, exam.Equal(gotV, wantV))
			}
			_, ok := wrapped.Get("99")
			exam.Try(t, exam.Equal(ok, false))

			gotKV := maps.Collect(wrapped.KeyVals())
			exam.Try(t, exam.MapEqual(gotKV, c.WantData))

			gotKeys := slices.Collect(wrapped.Keys())
			slices.Sort(gotKeys)
			wantKeys := slices.Collect(maps.Keys(c.WantData))
			slices.Sort(wantKeys)
			exam.Try(t, exam.SliceEqual(gotKeys, wantKeys))

			gotVals := slices.Collect(wrapped.Vals())
			slices.Sort(gotVals)
			wantVals := slices.Collect(maps.Values(c.WantData))
			slices.Sort(wantVals)
			exam.Try(t, exam.SliceEqual(gotVals, wantVals))
		})
	}
}

func TestDictEqual(t *testing.T) {
	cases := []struct {
		Name      string
		Loc       exam.Loc
		A         pod.DictView[string, int]
		B         pod.DictView[string, int]
		WantEqual bool
	}{
		{
			Name:      "both empty",
			Loc:       exam.Here(),
			A:         pod.NewMap[string, int](),
			B:         pod.NewMap[string, int](),
			WantEqual: true,
		},
		{
			Name:      "equal non-empty maps",
			Loc:       exam.Here(),
			A:         &pod.Map[string, int]{"a": 1, "b": 2},
			B:         &pod.Map[string, int]{"a": 1, "b": 2},
			WantEqual: true,
		},
		{
			Name:      "different values same keys",
			Loc:       exam.Here(),
			A:         &pod.Map[string, int]{"a": 1},
			B:         &pod.Map[string, int]{"a": 99},
			WantEqual: false,
		},
		{
			Name:      "different keys same length",
			Loc:       exam.Here(),
			A:         &pod.Map[string, int]{"a": 1},
			B:         &pod.Map[string, int]{"b": 1},
			WantEqual: false,
		},
		{
			Name:      "different lengths",
			Loc:       exam.Here(),
			A:         &pod.Map[string, int]{"a": 1, "b": 2},
			B:         &pod.Map[string, int]{"a": 1},
			WantEqual: false,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			exam.Try(t, exam.Equal(pod.DictEqual(c.A, c.B), c.WantEqual))
		})
	}
}

func TestDictEqualFunc(t *testing.T) {
	caseInsensitiveEq := func(a, b string) bool {
		return strings.EqualFold(a, b)
	}

	cases := []struct {
		Name      string
		Loc       exam.Loc
		A         pod.DictView[int, string]
		B         pod.DictView[int, string]
		WantEqual bool
	}{
		{
			Name:      "equal with case-insensitive comparison",
			Loc:       exam.Here(),
			A:         &pod.Map[int, string]{1: "hello"},
			B:         &pod.Map[int, string]{1: "HELLO"},
			WantEqual: true,
		},
		{
			Name:      "not equal with case-insensitive comparison",
			Loc:       exam.Here(),
			A:         &pod.Map[int, string]{1: "hello"},
			B:         &pod.Map[int, string]{1: "world"},
			WantEqual: false,
		},
		{
			Name:      "different lengths",
			Loc:       exam.Here(),
			A:         &pod.Map[int, string]{1: "hello", 2: "world"},
			B:         &pod.Map[int, string]{1: "hello"},
			WantEqual: false,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			exam.Try(t, exam.Equal(pod.DictEqualFunc(c.A, c.B, caseInsensitiveEq), c.WantEqual))
		})
	}
}
