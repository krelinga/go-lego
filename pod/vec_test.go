package pod_test

import (
	"slices"
	"testing"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/pod"
)

func TestVec(t *testing.T) {
	t.Run("literals", func(t *testing.T) {
		cases := []struct {
			Name        string
			Loc         exam.Loc
			Vec         *pod.Vec[int]
			WantVals    []int
			WantRevVals []int
		}{
			{
				Name:        "empty",
				Loc:         exam.Here(),
				Vec:         &pod.Vec[int]{},
				WantVals:    nil,
				WantRevVals: nil,
			},
		}
		for _, c := range cases {
			exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
				exam.Must(t, len(c.WantVals) == len(c.WantRevVals), len(c.WantVals), len(c.WantRevVals))

				exam.Must(t, c.Vec.Len() == len(c.WantVals), c.Vec.Len(), len(c.WantVals))

				for i := range c.WantVals {
					exam.Try(t, c.Vec.At(i) == c.WantVals[i], c.Vec.At(i), c.WantVals[i])
				}

				gotVals := slices.Collect(c.Vec.Vals())
				for i := range gotVals {
					exam.Try(t, gotVals[i] == c.WantVals[i], gotVals[i], c.WantVals[i])
				}

				gotRevVals := slices.Collect(c.Vec.RevVals())
				for i := range gotRevVals {
					exam.Must(t, gotRevVals[i] == c.WantRevVals[i], gotRevVals[i], c.WantRevVals[i])
				}
			})
		}
	})
}
