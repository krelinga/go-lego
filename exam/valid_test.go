package exam_test

import (
	"fmt"
	"testing"

	"github.com/krelinga/go-libs/exam"
)

type StringMustStartWithA string

func (s StringMustStartWithA) Validate() error {
	if len(s) == 0 || s[0] != 'A' {
		return fmt.Errorf("string must start with 'A', got %q", s)
	}
	return nil
}

func TestMustValidate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		fakeT := &FakeT{T: t}
		value := StringMustStartWithA("Apple")
		validatedValue := exam.MustValidate(fakeT, value)
		exam.Try(t, exam.Equal(validatedValue, value))
		exam.Try(t, exam.Equal(fakeT.FatalBody, ""))
		exam.Try(t, exam.Equal(fakeT.ErrorBody, ""))
	})

	t.Run("Failure", func(t *testing.T) {
		fakeT := &FakeT{T: t}
		value := StringMustStartWithA("Banana")
		validatedValue := exam.MustValidate(fakeT, value)
		exam.Try(t, exam.Equal(validatedValue, value))
		golden := exam.GoldenHere(`
FATAL: validatedValue := exam.MustValidate(fakeT, value)
validation failed: string must start with 'A', got "Banana"`)
		exam.Try(t, exam.GoldenEqual(fakeT.FatalBody, golden))
		exam.Try(t, exam.Equal(fakeT.ErrorBody, ""))
	})
}