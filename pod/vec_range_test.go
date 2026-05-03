package pod_test

import (
	"slices"
	"testing"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/pod"
)

func TestVecRange(t *testing.T) {
	cases := []struct {
		Name      string
		Loc       exam.Loc
		Parent    pod.VecView[int]
		FromIdx   int
		ToIdx     int
		WantLen   int
		WantPanic bool
	}{
		{
			Name:    "full range",
			Loc:     exam.Here(),
			Parent:  pod.AsVec([]int{1, 2, 3}),
			FromIdx: 0,
			ToIdx:   3,
			WantLen: 3,
		},
		{
			Name:    "partial range",
			Loc:     exam.Here(),
			Parent:  pod.AsVec([]int{1, 2, 3}),
			FromIdx: 1,
			ToIdx:   3,
			WantLen: 2,
		},
		{
			Name:    "empty range",
			Loc:     exam.Here(),
			Parent:  pod.AsVec([]int{1, 2, 3}),
			FromIdx: 2,
			ToIdx:   2,
			WantLen: 0,
		},
		{
			Name:      "fromIdx negative",
			Loc:       exam.Here(),
			Parent:    pod.AsVec([]int{1, 2, 3}),
			FromIdx:   -1,
			ToIdx:     2,
			WantPanic: true,
		},
		{
			Name:      "toIdx beyond length",
			Loc:       exam.Here(),
			Parent:    pod.AsVec([]int{1, 2, 3}),
			FromIdx:   0,
			ToIdx:     4,
			WantPanic: true,
		},
		{
			Name:      "fromIdx greater than toIdx",
			Loc:       exam.Here(),
			Parent:    pod.AsVec([]int{1, 2, 3}),
			FromIdx:   2,
			ToIdx:     1,
			WantPanic: true,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			if c.WantPanic {
				exam.Try(t, exam.PanicsWith(func() {
					pod.VecRange(c.Parent, c.FromIdx, c.ToIdx)
				}, func(p any) *exam.Failure {
					return exam.Equal(exam.MustCast[string](t, p), "invalid range")
				}))
			} else {
				r := pod.VecRange(c.Parent, c.FromIdx, c.ToIdx)
				exam.Try(t, exam.Equal(r.Len(), c.WantLen))
			}
		})
	}
}

func TestVecRangeFrom(t *testing.T) {
	cases := []struct {
		Name      string
		Loc       exam.Loc
		Parent    pod.VecView[int]
		FromIdx   int
		WantLen   int
		WantPanic bool
	}{
		{
			Name:    "from start",
			Loc:     exam.Here(),
			Parent:  pod.AsVec([]int{1, 2, 3}),
			FromIdx: 0,
			WantLen: 3,
		},
		{
			Name:    "from middle",
			Loc:     exam.Here(),
			Parent:  pod.AsVec([]int{1, 2, 3}),
			FromIdx: 1,
			WantLen: 2,
		},
		{
			Name:    "from end produces empty range",
			Loc:     exam.Here(),
			Parent:  pod.AsVec([]int{1, 2, 3}),
			FromIdx: 3,
			WantLen: 0,
		},
		{
			Name:      "fromIdx negative",
			Loc:       exam.Here(),
			Parent:    pod.AsVec([]int{1, 2, 3}),
			FromIdx:   -1,
			WantPanic: true,
		},
		{
			Name:      "fromIdx beyond length",
			Loc:       exam.Here(),
			Parent:    pod.AsVec([]int{1, 2, 3}),
			FromIdx:   4,
			WantPanic: true,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			if c.WantPanic {
				exam.Try(t, exam.PanicsWith(func() {
					pod.VecRangeFrom(c.Parent, c.FromIdx)
				}, func(p any) *exam.Failure {
					return exam.Equal(exam.MustCast[string](t, p), "invalid range")
				}))
			} else {
				r := pod.VecRangeFrom(c.Parent, c.FromIdx)
				exam.Try(t, exam.Equal(r.Len(), c.WantLen))
			}
		})
	}
}

func TestVecRangeTo(t *testing.T) {
	cases := []struct {
		Name      string
		Loc       exam.Loc
		Parent    pod.VecView[int]
		ToIdx     int
		WantLen   int
		WantPanic bool
	}{
		{
			Name:    "to end",
			Loc:     exam.Here(),
			Parent:  pod.AsVec([]int{1, 2, 3}),
			ToIdx:   3,
			WantLen: 3,
		},
		{
			Name:    "to middle",
			Loc:     exam.Here(),
			Parent:  pod.AsVec([]int{1, 2, 3}),
			ToIdx:   2,
			WantLen: 2,
		},
		{
			Name:    "to start produces empty range",
			Loc:     exam.Here(),
			Parent:  pod.AsVec([]int{1, 2, 3}),
			ToIdx:   0,
			WantLen: 0,
		},
		{
			Name:      "toIdx negative",
			Loc:       exam.Here(),
			Parent:    pod.AsVec([]int{1, 2, 3}),
			ToIdx:     -1,
			WantPanic: true,
		},
		{
			Name:      "toIdx beyond length",
			Loc:       exam.Here(),
			Parent:    pod.AsVec([]int{1, 2, 3}),
			ToIdx:     4,
			WantPanic: true,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			if c.WantPanic {
				exam.Try(t, exam.PanicsWith(func() {
					pod.VecRangeTo(c.Parent, c.ToIdx)
				}, func(p any) *exam.Failure {
					return exam.Equal(exam.MustCast[string](t, p), "invalid range")
				}))
			} else {
				r := pod.VecRangeTo(c.Parent, c.ToIdx)
				exam.Try(t, exam.Equal(r.Len(), c.WantLen))
			}
		})
	}
}

func TestVecRangeLen(t *testing.T) {
	cases := []struct {
		Name    string
		Loc     exam.Loc
		Range   pod.VecView[int]
		WantLen int
	}{
		{
			Name:    "full range",
			Loc:     exam.Here(),
			Range:   pod.VecRange(pod.AsVec([]int{1, 2, 3}), 0, 3),
			WantLen: 3,
		},
		{
			Name:    "partial range",
			Loc:     exam.Here(),
			Range:   pod.VecRange(pod.AsVec([]int{1, 2, 3}), 1, 2),
			WantLen: 1,
		},
		{
			Name:    "empty range",
			Loc:     exam.Here(),
			Range:   pod.VecRange(pod.AsVec([]int{1, 2, 3}), 1, 1),
			WantLen: 0,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			exam.Try(t, exam.Equal(c.Range.Len(), c.WantLen))
		})
	}
	t.Run("parent too short", func(t *testing.T) {
		parent := &pod.Slice[int]{1, 2, 3}
		r := pod.VecRange(parent, 0, 3)
		parent.Resize(1)
		exam.Try(t, exam.PanicsWith(func() {
			r.Len()
		}, func(p any) *exam.Failure {
			return exam.Equal(exam.MustCast[string](t, p), "parent vector is too short for range")
		}))
	})
}

func TestVecRangeGet(t *testing.T) {
	cases := []struct {
		Name      string
		Loc       exam.Loc
		Range     pod.VecView[int]
		Idx       int
		WantVal   int
		WantPanic string
	}{
		{
			Name:    "first element of sub-range",
			Loc:     exam.Here(),
			Range:   pod.VecRange(pod.AsVec([]int{10, 20, 30}), 1, 3),
			Idx:     0,
			WantVal: 20,
		},
		{
			Name:    "last element of sub-range",
			Loc:     exam.Here(),
			Range:   pod.VecRange(pod.AsVec([]int{10, 20, 30}), 1, 3),
			Idx:     1,
			WantVal: 30,
		},
		{
			Name:      "index negative",
			Loc:       exam.Here(),
			Range:     pod.VecRange(pod.AsVec([]int{10, 20, 30}), 0, 3),
			Idx:       -1,
			WantPanic: "index out of range",
		},
		{
			Name:      "index equal to range length",
			Loc:       exam.Here(),
			Range:     pod.VecRange(pod.AsVec([]int{10, 20, 30}), 0, 3),
			Idx:       3,
			WantPanic: "index out of range",
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			if c.WantPanic != "" {
				exam.Try(t, exam.PanicsWith(func() {
					c.Range.Get(c.Idx)
				}, func(p any) *exam.Failure {
					return exam.Equal(exam.MustCast[string](t, p), c.WantPanic)
				}))
			} else {
				exam.Try(t, exam.Equal(c.Range.Get(c.Idx), c.WantVal))
			}
		})
	}
	t.Run("parent too short", func(t *testing.T) {
		parent := &pod.Slice[int]{10, 20, 30}
		r := pod.VecRange(parent, 0, 3)
		parent.Resize(1)
		exam.Try(t, exam.PanicsWith(func() {
			r.Get(0)
		}, func(p any) *exam.Failure {
			return exam.Equal(exam.MustCast[string](t, p), "parent vector is too short for range")
		}))
	})
}

func TestVecRangeVals(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Range    pod.VecView[int]
		WantVals []int
	}{
		{
			Name:     "full range",
			Loc:      exam.Here(),
			Range:    pod.VecRange(pod.AsVec([]int{1, 2, 3}), 0, 3),
			WantVals: []int{1, 2, 3},
		},
		{
			Name:     "partial range",
			Loc:      exam.Here(),
			Range:    pod.VecRange(pod.AsVec([]int{1, 2, 3}), 1, 3),
			WantVals: []int{2, 3},
		},
		{
			Name:     "empty range",
			Loc:      exam.Here(),
			Range:    pod.VecRange(pod.AsVec([]int{1, 2, 3}), 1, 1),
			WantVals: nil,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			got := slices.Collect(c.Range.Vals())
			exam.Must(t, exam.Equal(len(got), len(c.WantVals)))
			for i := range c.WantVals {
				exam.Try(t, exam.Equal(got[i], c.WantVals[i]))
			}
		})
	}
	t.Run("parent too short", func(t *testing.T) {
		parent := &pod.Slice[int]{1, 2, 3}
		r := pod.VecRange(parent, 0, 3)
		parent.Resize(1)
		exam.Try(t, exam.PanicsWith(func() {
			r.Vals()
		}, func(p any) *exam.Failure {
			return exam.Equal(exam.MustCast[string](t, p), "parent vector is too short for range")
		}))
	})
}

func TestVecRangeIdxVals(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Range    pod.VecView[int]
		WantVals []int
	}{
		{
			Name:     "full range",
			Loc:      exam.Here(),
			Range:    pod.VecRange(pod.AsVec([]int{1, 2, 3}), 0, 3),
			WantVals: []int{1, 2, 3},
		},
		{
			Name:     "partial range",
			Loc:      exam.Here(),
			Range:    pod.VecRange(pod.AsVec([]int{1, 2, 3}), 1, 3),
			WantVals: []int{2, 3},
		},
		{
			Name:     "empty range",
			Loc:      exam.Here(),
			Range:    pod.VecRange(pod.AsVec([]int{1, 2, 3}), 1, 1),
			WantVals: nil,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			wantIdx := 0
			for idx, val := range c.Range.IdxVals() {
				exam.Try(t, exam.Equal(idx, wantIdx))
				exam.Try(t, exam.Equal(val, c.WantVals[wantIdx]))
				wantIdx++
			}
			exam.Try(t, exam.Equal(wantIdx, len(c.WantVals)))
		})
	}
	t.Run("parent too short", func(t *testing.T) {
		parent := &pod.Slice[int]{1, 2, 3}
		r := pod.VecRange(parent, 0, 3)
		parent.Resize(1)
		exam.Try(t, exam.PanicsWith(func() {
			r.IdxVals()
		}, func(p any) *exam.Failure {
			return exam.Equal(exam.MustCast[string](t, p), "parent vector is too short for range")
		}))
	})
}

func TestVecRangeRevVals(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Range    pod.VecView[int]
		WantVals []int
	}{
		{
			Name:     "full range",
			Loc:      exam.Here(),
			Range:    pod.VecRange(pod.AsVec([]int{1, 2, 3}), 0, 3),
			WantVals: []int{3, 2, 1},
		},
		{
			Name:     "partial range",
			Loc:      exam.Here(),
			Range:    pod.VecRange(pod.AsVec([]int{1, 2, 3}), 1, 3),
			WantVals: []int{3, 2},
		},
		{
			Name:     "empty range",
			Loc:      exam.Here(),
			Range:    pod.VecRange(pod.AsVec([]int{1, 2, 3}), 1, 1),
			WantVals: nil,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			got := slices.Collect(c.Range.RevVals())
			exam.Must(t, exam.Equal(len(got), len(c.WantVals)))
			for i := range c.WantVals {
				exam.Try(t, exam.Equal(got[i], c.WantVals[i]))
			}
		})
	}
	t.Run("parent too short", func(t *testing.T) {
		parent := &pod.Slice[int]{1, 2, 3}
		r := pod.VecRange(parent, 0, 3)
		parent.Resize(1)
		exam.Try(t, exam.PanicsWith(func() {
			r.RevVals()
		}, func(p any) *exam.Failure {
			return exam.Equal(exam.MustCast[string](t, p), "parent vector is too short for range")
		}))
	})
}

func TestVecRangeRevIdxVals(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		Range    pod.VecView[int]
		WantVals []int // forward-ordered; indices emitted by RevIdxVals correspond to positions in this slice
	}{
		{
			Name:     "full range",
			Loc:      exam.Here(),
			Range:    pod.VecRange(pod.AsVec([]int{1, 2, 3}), 0, 3),
			WantVals: []int{1, 2, 3},
		},
		{
			Name:     "partial range",
			Loc:      exam.Here(),
			Range:    pod.VecRange(pod.AsVec([]int{1, 2, 3}), 1, 3),
			WantVals: []int{2, 3},
		},
		{
			Name:     "empty range",
			Loc:      exam.Here(),
			Range:    pod.VecRange(pod.AsVec([]int{1, 2, 3}), 1, 1),
			WantVals: nil,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			wantIdx := len(c.WantVals) - 1
			for idx, val := range c.Range.RevIdxVals() {
				exam.Try(t, exam.Equal(idx, wantIdx))
				exam.Try(t, exam.Equal(val, c.WantVals[wantIdx]))
				wantIdx--
			}
			exam.Try(t, exam.Equal(wantIdx, -1))
		})
	}
	t.Run("parent too short", func(t *testing.T) {
		parent := &pod.Slice[int]{1, 2, 3}
		r := pod.VecRange(parent, 0, 3)
		parent.Resize(1)
		exam.Try(t, exam.PanicsWith(func() {
			r.RevIdxVals()
		}, func(p any) *exam.Failure {
			return exam.Equal(exam.MustCast[string](t, p), "parent vector is too short for range")
		}))
	})
}
