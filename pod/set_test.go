package pod_test

import (
	"testing"

	"github.com/krelinga/go-libs/exam"
	"github.com/krelinga/go-libs/pod"
)

func TestSetImplements(t *testing.T) {
	exam.Try(t, exam.Implements[*pod.MapSet[int], pod.Set[int]]())

	exam.Try(t, exam.Implements[*pod.SliceSet[int], pod.Set[int]]())
}
