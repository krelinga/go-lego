package exam_test

import (
	"strings"
	"testing"

	"github.com/krelinga/go-libs/exam"
)

type IFace interface {
	IFaceMethod()
}

type ImplementsIFace struct{}

func (i ImplementsIFace) IFaceMethod() {}

func TestImplements(t *testing.T) {
	cases := []struct {
		Name             string
		Loc              exam.Loc
		Failure          *exam.Failure
		WantNilFailure   bool
		WantWrappedError bool
	}{
		{
			Name:             "Implements",
			Loc:              exam.Here(),
			Failure:          exam.Implements[ImplementsIFace, IFace](),
			WantNilFailure:   true,
			WantWrappedError: false,
		},
		{
			Name:             "DoesNotImplement",
			Loc:              exam.Here(),
			Failure:          exam.Implements[int, IFace](),
			WantNilFailure:   false,
			WantWrappedError: false,
		},
		{
			Name:             "NonInterface",
			Loc:              exam.Here(),
			Failure:          exam.Implements[int, int](),
			WantNilFailure:   false,
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

func TestImplementsArgFormat(t *testing.T) {
	failure := exam.Implements[int, IFace]()
	exam.Must(t, exam.NotNil(failure))
	exam.Must(t, exam.Equal(len(failure.Args), 2))

	gotStr, err := failure.Args[0].ToString()
	exam.Must(t, exam.Nil(err))
	if strings.Contains(gotStr, "reflect.rtype") {
		t.Errorf("Args[0] should not contain reflect.rtype, got: %s", gotStr)
	}
	exam.Try(t, exam.Equal(gotStr, "Got: int"))

	ifaceStr, err := failure.Args[1].ToString()
	exam.Must(t, exam.Nil(err))
	if strings.Contains(ifaceStr, "reflect.rtype") {
		t.Errorf("Args[1] should not contain reflect.rtype, got: %s", ifaceStr)
	}
	exam.Try(t, exam.Equal(ifaceStr, "IFace: exam_test.IFace"))
}

func TestPanicsWith(t *testing.T) {
	cases := []struct {
		Name        string
		Loc         exam.Loc
		Run         func()
		Check       func(any) *exam.Failure
		WantFailure bool
	}{
		{
			Name: "Panics With Matching Failure",
			Loc:  exam.Here(),
			Run: func() {
				panic(int(10))
			},
			Check: func(p any) *exam.Failure {
				v, ok := p.(int)
				if !ok {
					return exam.NewFailure1("panic value type", p)
				}
				return exam.Equal(v, 10)
			},
			WantFailure: false,
		},
		{
			Name: "Panics With Non-Matching Failure",
			Loc:  exam.Here(),
			Run: func() {
				panic(int(20))
			},
			Check: func(p any) *exam.Failure {
				v, ok := p.(int)
				if !ok {
					return exam.NewFailure1("panic value type", p)
				}
				return exam.Equal(v, 10)
			},
			WantFailure: true,
		},
		{
			Name: "Does Not Panic",
			Loc:  exam.Here(),
			Run: func() {
				// do nothing
			},
			Check: func(p any) *exam.Failure {
				return nil
			},
			WantFailure: true,
		},
		{
			Name: "Panics With nil Check",
			Loc:  exam.Here(),
			Run: func() {
				panic("panic with nil check")
			},
			Check:       nil,
			WantFailure: false,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			f := exam.PanicsWith(c.Run, c.Check)
			if c.WantFailure {
				exam.Try(t, exam.NotNil(f))
			} else {
				exam.Try(t, exam.Nil(f))
			}
		})
	}
}
