package exam_test

import "testing"
import "github.com/krelinga/go-lego/exam"

type IFace interface {
	IFaceMethod()
}

type ImplementsIFace struct{}

func (i ImplementsIFace) IFaceMethod() {}

func TestImplements(t *testing.T) {
	cases := []struct {
		Name string
		Loc exam.Loc
		Failure *exam.Failure
		WantNilFailure bool
		WantWrappedError bool
	}{
		{
			Name: "Implements",
			Loc: exam.Here(),
			Failure: exam.Implements[ImplementsIFace, IFace](),
			WantNilFailure: true,
			WantWrappedError: false,
		},
		{
			Name: "DoesNotImplement",
			Loc: exam.Here(),
			Failure: exam.Implements[int, IFace](),
			WantNilFailure: false,
			WantWrappedError: false,
		},
		{
			Name: "NonInterface",
			Loc: exam.Here(),
			Failure: exam.Implements[int, int](),
			WantNilFailure: false,
			WantWrappedError: true,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			if c.WantNilFailure {
				exam.Try(t, exam.Nil(c.Failure))
			} else {
				exam.Must(t, exam.NotNil(c.Failure))
				if c.WantWrappedError {
					exam.Try(t, exam.NotNil(c.Failure.Wrapped))
				}
			}
		})
	}
}