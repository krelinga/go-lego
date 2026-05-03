package pod_test

import (
	"testing"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/pod"
)

func TestDictImplements(t *testing.T) {
	exam.Try(t, exam.Implements[*pod.Map[int, string], pod.Dict[int, string]]())
}