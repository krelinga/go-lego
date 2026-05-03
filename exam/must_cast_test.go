package exam_test

import (
	"testing"

	"github.com/krelinga/go-lego/exam"
)

func TestMustCast(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		fakeT := &FakeT{T: t}
		value := 42
		castedValue := exam.MustCast[int](fakeT, value)
		exam.Try(t, exam.Equal(castedValue, value))
		exam.Try(t, exam.Equal(fakeT.FatalBody, ""))
		exam.Try(t, exam.Equal(fakeT.ErrorBody, ""))
	})

	t.Run("Failure", func(t *testing.T) {
		fakeT := &FakeT{T: t}
		value := "not an int"
		castedValue := exam.MustCast[int](fakeT, value)
		exam.Try(t, exam.Equal(castedValue, 0))
		golden := exam.GoldenHere(`
FATAL: castedValue := exam.MustCast[int](fakeT, value)
expected value that could be cast to type int, got string`)
		exam.Try(t, exam.GoldenEqual(fakeT.FatalBody, golden))
		exam.Try(t, exam.Equal(fakeT.ErrorBody, ""))
	})
}