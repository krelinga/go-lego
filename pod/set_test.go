package pod_test

import (
	"testing"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/pod"
)

func TestSetImplements(t *testing.T) {
	exam.Try(t, exam.Implements[*pod.MapSet[int], pod.Set[int]]())
}