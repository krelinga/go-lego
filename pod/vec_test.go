package pod_test

import (
	"slices"
	"testing"

	"github.com/krelinga/go-libs/exam"
	"github.com/krelinga/go-libs/pod"
)

func TestVecImplements(t *testing.T) {
	exam.Try(t, exam.Implements[*pod.Slice[int], pod.Vec[int]]())
}

func TestFixedVec(t *testing.T) {
	cases := []struct {
		Name     string
		Loc      exam.Loc
		FixedVec pod.FixedVec[int]
		WantVals []int
	}{
		{
			Name:     "empty Slice",
			Loc:      exam.Here(),
			FixedVec: &pod.Slice[int]{},
			WantVals: nil,
		},
		{
			Name:     "non-empty Slice",
			Loc:      exam.Here(),
			FixedVec: &pod.Slice[int]{1, 2, 3},
			WantVals: []int{1, 2, 3},
		},
		{
			Name:     "AsVec nil",
			Loc:      exam.Here(),
			FixedVec: pod.AsVec([]int(nil)),
			WantVals: nil,
		},
		{
			Name:     "AsVec empty",
			Loc:      exam.Here(),
			FixedVec: pod.AsVec([]int{}),
			WantVals: nil,
		},
		{
			Name:     "AsVec non-empty",
			Loc:      exam.Here(),
			FixedVec: pod.AsVec([]int{1, 2, 3}),
			WantVals: []int{1, 2, 3},
		},
		{
			Name: "Wrapped Empty",
			Loc:  exam.Here(),
			FixedVec: pod.WrapVecVals(
				pod.AsVec([]float64(nil)),
				func(x float64) int { return int(x) },
			),
			WantVals: nil,
		},
		{
			Name: "Wrapped non-empty",
			Loc:  exam.Here(),
			FixedVec: pod.WrapVecVals(
				pod.AsVec([]float64{1.0, 2.0, 3.0}),
				func(x float64) int { return int(x) },
			),
			WantVals: []int{1, 2, 3},
		},
		{
			Name:     "VecRange Full",
			Loc:      exam.Here(),
			FixedVec: pod.VecRange(pod.AsVec([]int{1, 2, 3}), 0, 3),
			WantVals: []int{1, 2, 3},
		},
		{
			Name:     "VecRange Partial",
			Loc:      exam.Here(),
			FixedVec: pod.VecRange(pod.AsVec([]int{1, 2, 3}), 1, 3),
			WantVals: []int{2, 3},
		},
		{
			Name:     "VecRange Empty",
			Loc:      exam.Here(),
			FixedVec: pod.VecRange(pod.AsVec([]int{1, 2, 3}), 1, 1),
			WantVals: nil,
		},
		{
			Name:     "VecRangeFrom to end",
			Loc:      exam.Here(),
			FixedVec: pod.VecRangeFrom(pod.AsVec([]int{1, 2, 3}), 1),
			WantVals: []int{2, 3},
		},
		{
			Name:     "VecRangeTo from start",
			Loc:      exam.Here(),
			FixedVec: pod.VecRangeTo(pod.AsVec([]int{1, 2, 3}), 2),
			WantVals: []int{1, 2},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			exam.Must(t, exam.Equal(c.FixedVec.Len(), len(c.WantVals)))

			for i := range c.WantVals {
				exam.Try(t, exam.Equal(c.FixedVec.Get(i), c.WantVals[i]))
			}

			gotVals := slices.Collect(c.FixedVec.Vals())
			exam.Try(t, exam.SliceEqual(gotVals, c.WantVals))
			wantIdx := 0
			for idx, val := range c.FixedVec.IdxVals() {
				exam.Try(t, exam.Equal(idx, wantIdx))
				exam.Try(t, exam.Equal(val, c.WantVals[idx]))
				wantIdx++
			}
			exam.Try(t, exam.Equal(wantIdx, len(c.WantVals)))

			gotRevVals := slices.Collect(c.FixedVec.RevVals())
			wantRevVals := slices.Clone(c.WantVals)
			slices.Reverse(wantRevVals)
			exam.Try(t, exam.SliceEqual(gotRevVals, wantRevVals))
			wantIdx = len(c.WantVals) - 1
			for idx, val := range c.FixedVec.RevIdxVals() {
				exam.Try(t, exam.Equal(idx, wantIdx))
				exam.Try(t, exam.Equal(val, c.WantVals[idx]))
				wantIdx--
			}
			exam.Try(t, exam.Equal(wantIdx, -1))
		})
	}
}
