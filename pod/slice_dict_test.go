package pod_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/krelinga/go-libs/exam"
	"github.com/krelinga/go-libs/pod"
	"github.com/krelinga/go-libs/tuple"
)

func TestSliceDictImplements(t *testing.T) {
	exam.Try(t, exam.Implements[*pod.SliceDict[string, int], pod.Dict[string, int]]())
}

func TestNewSliceDict(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Entries  []tuple.View2[string, int]
		WantLen  int
		WantData map[string]int
	}{
		{
			Name:     "no entries",
			Loc:      exam.Here(),
			Entries:  nil,
			WantLen:  0,
			WantData: map[string]int{},
		},
		{
			Name:     "distinct entries",
			Loc:      exam.Here(),
			Entries:  []tuple.View2[string, int]{tuple.New2("a", 1), tuple.New2("b", 2)},
			WantLen:  2,
			WantData: map[string]int{"a": 1, "b": 2},
		},
		{
			Name:     "duplicate key keeps last",
			Loc:      exam.Here(),
			Entries:  []tuple.View2[string, int]{tuple.New2("a", 1), tuple.New2("a", 99)},
			WantLen:  1,
			WantData: map[string]int{"a": 99},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			d := pod.NewSliceDict(c.Entries...)
			exam.Try(t, exam.Equal(d.Len(), c.WantLen))
			for k, wantV := range c.WantData {
				gotV, ok := d.Get(k)
				exam.Try(t, exam.True(ok))
				exam.Try(t, exam.Equal(gotV, wantV))
			}
		})
	}
}

func TestNewSliceDictFunc(t *testing.T) {
	caseInsensitive := func(a, b string) bool {
		return strings.EqualFold(a, b)
	}
	t.Run("distinct keys", func(t *testing.T) {
		d := pod.NewSliceDictFunc(caseInsensitive,
			tuple.New2("a", 1),
			tuple.New2("b", 2),
		)
		exam.Try(t, exam.Equal(d.Len(), 2))
		v, ok := d.Get("A")
		exam.Try(t, exam.True(ok))
		exam.Try(t, exam.Equal(v, 1))
	})
	t.Run("duplicate key by custom equality keeps last", func(t *testing.T) {
		d := pod.NewSliceDictFunc(caseInsensitive,
			tuple.New2("a", 1),
			tuple.New2("A", 99),
		)
		exam.Try(t, exam.Equal(d.Len(), 1))
		v, ok := d.Get("a")
		exam.Try(t, exam.True(ok))
		exam.Try(t, exam.Equal(v, 99))
	})
}

func TestNewSliceDictOf(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Source   *pod.Map[string, int]
		WantLen  int
		WantData map[string]int
	}{
		{
			Name:     "empty source",
			Loc:      exam.Here(),
			Source:   pod.NewMap[string, int](),
			WantLen:  0,
			WantData: map[string]int{},
		},
		{
			Name:     "non-empty source",
			Loc:      exam.Here(),
			Source:   &pod.Map[string, int]{"a": 1, "b": 2},
			WantLen:  2,
			WantData: map[string]int{"a": 1, "b": 2},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			d := pod.NewSliceDictOf(pod.KeyValsOf(c.Source))
			exam.Try(t, exam.Equal(d.Len(), c.WantLen))
			for k, wantV := range c.WantData {
				gotV, ok := d.Get(k)
				exam.Try(t, exam.True(ok))
				exam.Try(t, exam.Equal(gotV, wantV))
			}
		})
	}
}

func TestNewSliceDictOfFunc(t *testing.T) {
	caseInsensitive := func(a, b string) bool {
		return strings.EqualFold(a, b)
	}
	source := &pod.Map[string, int]{"x": 10, "y": 20}
	d := pod.NewSliceDictOfFunc(caseInsensitive, pod.KeyValsOf(source))
	exam.Try(t, exam.Equal(d.Len(), 2))
	v, ok := d.Get("X")
	exam.Try(t, exam.True(ok))
	exam.Try(t, exam.Equal(v, 10))
}

func TestSliceDictLen(t *testing.T) {
	cases := []struct {
		Name    string
		Loc     exam.Loc
		Dict    *pod.SliceDict[string, int]
		WantLen int
	}{
		{
			Name:    "empty dict",
			Loc:     exam.Here(),
			Dict:    pod.NewSliceDict[string, int](),
			WantLen: 0,
		},
		{
			Name:    "one entry",
			Loc:     exam.Here(),
			Dict:    pod.NewSliceDict(tuple.New2("a", 1)),
			WantLen: 1,
		},
		{
			Name:    "three entries",
			Loc:     exam.Here(),
			Dict:    pod.NewSliceDict(tuple.New2("a", 1), tuple.New2("b", 2), tuple.New2("c", 3)),
			WantLen: 3,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			exam.Try(t, exam.Equal(c.Dict.Len(), c.WantLen))
		})
	}
}

func TestSliceDictGet(t *testing.T) {
	cases := []struct {
		Name      string
		Loc       exam.Loc
		Dict      *pod.SliceDict[string, int]
		Key       string
		WantVal   int
		WantFound bool
	}{
		{
			Name:      "existing key",
			Loc:       exam.Here(),
			Dict:      pod.NewSliceDict(tuple.New2("a", 1), tuple.New2("b", 2)),
			Key:       "a",
			WantVal:   1,
			WantFound: true,
		},
		{
			Name:      "missing key",
			Loc:       exam.Here(),
			Dict:      pod.NewSliceDict(tuple.New2("a", 1)),
			Key:       "z",
			WantVal:   0,
			WantFound: false,
		},
		{
			Name:      "empty dict",
			Loc:       exam.Here(),
			Dict:      pod.NewSliceDict[string, int](),
			Key:       "a",
			WantVal:   0,
			WantFound: false,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			gotVal, gotFound := c.Dict.Get(c.Key)
			exam.Try(t, exam.Equal(gotFound, c.WantFound))
			exam.Try(t, exam.Equal(gotVal, c.WantVal))
		})
	}
}

func TestSliceDictKeyVals(t *testing.T) {
	cases := []struct {
		Name      string
		Loc       exam.Loc
		Dict      *pod.SliceDict[string, int]
		WantPairs []tuple.T2[string, int]
	}{
		{
			Name:      "empty dict",
			Loc:       exam.Here(),
			Dict:      pod.NewSliceDict[string, int](),
			WantPairs: nil,
		},
		{
			Name: "preserves insertion order",
			Loc:  exam.Here(),
			Dict: pod.NewSliceDict(
				tuple.New2("a", 1),
				tuple.New2("b", 2),
				tuple.New2("c", 3),
			),
			WantPairs: []tuple.T2[string, int]{
				tuple.New2("a", 1),
				tuple.New2("b", 2),
				tuple.New2("c", 3),
			},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			var got []tuple.T2[string, int]
			for k, v := range c.Dict.KeyVals() {
				got = append(got, tuple.New2(k, v))
			}
			exam.Try(t, exam.SliceEqualFunc(got, c.WantPairs, func(a, b tuple.T2[string, int]) bool {
				return a.A == b.A && a.B == b.B
			}))
		})
	}
}

func TestSliceDictKeys(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Dict     *pod.SliceDict[string, int]
		WantKeys []string
	}{
		{
			Name:     "empty dict",
			Loc:      exam.Here(),
			Dict:     pod.NewSliceDict[string, int](),
			WantKeys: nil,
		},
		{
			Name: "preserves insertion order",
			Loc:  exam.Here(),
			Dict: pod.NewSliceDict(
				tuple.New2("a", 1),
				tuple.New2("b", 2),
				tuple.New2("c", 3),
			),
			WantKeys: []string{"a", "b", "c"},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			got := slices.Collect(c.Dict.Keys())
			exam.Try(t, exam.SliceEqual(got, c.WantKeys))
		})
	}
}

func TestSliceDictVals(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Dict     *pod.SliceDict[string, int]
		WantVals []int
	}{
		{
			Name:     "empty dict",
			Loc:      exam.Here(),
			Dict:     pod.NewSliceDict[string, int](),
			WantVals: nil,
		},
		{
			Name: "preserves insertion order",
			Loc:  exam.Here(),
			Dict: pod.NewSliceDict(
				tuple.New2("a", 1),
				tuple.New2("b", 2),
				tuple.New2("c", 3),
			),
			WantVals: []int{1, 2, 3},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			got := slices.Collect(c.Dict.Vals())
			exam.Try(t, exam.SliceEqual(got, c.WantVals))
		})
	}
}

func TestSliceDictPut(t *testing.T) {
	t.Run("add to empty dict", func(t *testing.T) {
		d := pod.NewSliceDict[string, int]()
		d.Put("a", 1)
		exam.Try(t, exam.Equal(d.Len(), 1))
		v, ok := d.Get("a")
		exam.Try(t, exam.True(ok))
		exam.Try(t, exam.Equal(v, 1))
	})
	t.Run("add new key preserves order", func(t *testing.T) {
		d := pod.NewSliceDict(tuple.New2("a", 1))
		d.Put("b", 2)
		exam.Try(t, exam.Equal(d.Len(), 2))
		got := slices.Collect(d.Keys())
		exam.Try(t, exam.SliceEqual(got, []string{"a", "b"}))
	})
	t.Run("update existing key moves to end", func(t *testing.T) {
		d := pod.NewSliceDict(tuple.New2("a", 1), tuple.New2("b", 2), tuple.New2("c", 3))
		d.Put("b", 99)
		exam.Try(t, exam.Equal(d.Len(), 3))
		got := slices.Collect(d.Keys())
		exam.Try(t, exam.SliceEqual(got, []string{"a", "c", "b"}))
		v, ok := d.Get("b")
		exam.Try(t, exam.True(ok))
		exam.Try(t, exam.Equal(v, 99))
	})
}

func TestSliceDictClear(t *testing.T) {
	cases := []struct {
		Name string
		Loc  exam.Loc
		Dict *pod.SliceDict[string, int]
	}{
		{
			Name: "clear non-empty dict",
			Loc:  exam.Here(),
			Dict: pod.NewSliceDict(tuple.New2("a", 1), tuple.New2("b", 2)),
		},
		{
			Name: "clear already empty dict",
			Loc:  exam.Here(),
			Dict: pod.NewSliceDict[string, int](),
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			c.Dict.Clear()
			exam.Try(t, exam.Equal(c.Dict.Len(), 0))
			_, ok := c.Dict.Get("a")
			exam.Try(t, exam.Equal(ok, false))
		})
	}
}

func TestSliceDictDel(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Dict     *pod.SliceDict[string, int]
		DelKey   string
		WantLen  int
		WantKeys []string
		WantVals []int
	}{
		{
			Name:     "delete first entry",
			Loc:      exam.Here(),
			Dict:     pod.NewSliceDict(tuple.New2("a", 1), tuple.New2("b", 2), tuple.New2("c", 3)),
			DelKey:   "a",
			WantLen:  2,
			WantKeys: []string{"b", "c"},
			WantVals: []int{2, 3},
		},
		{
			Name:     "delete middle entry",
			Loc:      exam.Here(),
			Dict:     pod.NewSliceDict(tuple.New2("a", 1), tuple.New2("b", 2), tuple.New2("c", 3)),
			DelKey:   "b",
			WantLen:  2,
			WantKeys: []string{"a", "c"},
			WantVals: []int{1, 3},
		},
		{
			Name:     "delete last entry",
			Loc:      exam.Here(),
			Dict:     pod.NewSliceDict(tuple.New2("a", 1), tuple.New2("b", 2), tuple.New2("c", 3)),
			DelKey:   "c",
			WantLen:  2,
			WantKeys: []string{"a", "b"},
			WantVals: []int{1, 2},
		},
		{
			Name:     "delete missing key is a no-op",
			Loc:      exam.Here(),
			Dict:     pod.NewSliceDict(tuple.New2("a", 1)),
			DelKey:   "z",
			WantLen:  1,
			WantKeys: []string{"a"},
			WantVals: []int{1},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			c.Dict.Del(c.DelKey)
			exam.Try(t, exam.Equal(c.Dict.Len(), c.WantLen))
			exam.Try(t, exam.SliceEqual(slices.Collect(c.Dict.Keys()), c.WantKeys))
			exam.Try(t, exam.SliceEqual(slices.Collect(c.Dict.Vals()), c.WantVals))
		})
	}
}

func TestSliceDictReserve(t *testing.T) {
	t.Run("reserve then put", func(t *testing.T) {
		d := pod.NewSliceDict[string, int]()
		d.Reserve(10)
		d.Put("a", 1)
		exam.Try(t, exam.Equal(d.Len(), 1))
		v, ok := d.Get("a")
		exam.Try(t, exam.True(ok))
		exam.Try(t, exam.Equal(v, 1))
	})
}
