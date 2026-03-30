package lego_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/krelinga/go-lego"
)

type FailureReporter interface {
	Error(...any)
	Errorf(string, ...any)
	Fatal(...any)
	Fatalf(string, ...any)
	Helper()
}

func implements[T any, I any](r FailureReporter) bool {
	r.Helper()
	iType := reflect.TypeFor[I]()
	if iType.Kind() != reflect.Interface {
		r.Fatalf("I must be an interface type, got %s", iType)
	}
	tType := reflect.TypeFor[T]()
	if !tType.Implements(iType) {
		r.Errorf("%s does not implement %s", tType, iType)
	}
	return true
}

func panics(r FailureReporter, f func()) (panicked bool) {
	r.Helper()
	defer func() {
		r.Helper()
		if recover() == nil {
			r.Errorf("Expected function to panic, but it did not")
		} else {
			panicked = true
		}
	}()
	f()
	return
}

type ExampleView struct {
	e *Example
}

func (v ExampleView) String() string {
	return v.e.String
}

func (v ExampleView) Int() int {
	return v.e.Int
}

func (v ExampleView) Equal(other ExampleView) bool {
	return lego.EqualComparer(v, other)
}

func (v ExampleView) Combine() string {
	return fmt.Sprintf("%s%d", v.String(), v.Int())
}

func (v ExampleView) Compare(other ExampleView) int {
	return lego.CompareUsing(v, other,
		lego.NewCmpFuncGo(ExampleView.String),
		lego.NewCmpFuncGo(ExampleView.Int),
	)
}

type Example struct {
	String string
	Int    int
}

func (e *Example) Equal(other *Example) bool {
	return lego.EqualViewer(e, other)
}

func (e *Example) View() ExampleView {
	return ExampleView{e}
}

func (e *Example) Combine() string {
	return e.View().Combine()
}

func (e *Example) Compare(other *Example) int {
	return lego.CompareViewer(e, other)
}

type ExampleSliceView struct {
	lego.FixedSlice[ExampleView]
}

func (v ExampleSliceView) Equal(other ExampleSliceView) bool {
	return lego.EqualSlice(v, other)
}

func (v ExampleSliceView) Sum() int {
	sum := 0
	for p := range v.List() {
		sum += p.GetValue().Int()
	}
	return sum
}

type ExampleSlice struct {
	*lego.Slice[*Example]
}

func (s *ExampleSlice) Sum() int {
	return s.View().Sum()
}

func (s *ExampleSlice) View() ExampleSliceView {
	return ExampleSliceView{lego.ViewSlice(s)}
}

func (s *ExampleSlice) Equal(other *ExampleSlice) bool {
	return lego.EqualViewer(s, other)
}

func NewExampleSlice(elements ...*Example) *ExampleSlice {
	return &ExampleSlice{lego.NewSlice(elements...)}
}

func TestExample(t *testing.T) {
	t.Run("implements", func(t *testing.T) {
		implements[*Example, lego.Equaler[*Example]](t)
		implements[*Example, lego.Viewer[ExampleView]](t)
	})

	t.Run("equal", func(t *testing.T) {
		e1 := &Example{String: "a", Int: 1}
		e2 := &Example{String: "a", Int: 1}
		e3 := &Example{String: "a", Int: 2}

		if !e1.Equal(e2) {
			t.Errorf("Expected e1 to equal e2, but they are not equal")
		}
		if e1.Equal(e3) {
			t.Errorf("Expected e1 to not equal e3, but they are equal")
		}

		e1v := e1.View()
		e2v := e2.View()
		e3v := e3.View()

		if !e1v.Equal(e2v) {
			t.Errorf("Expected e1v to equal e2v, but they are not equal")
		}
		if e1v.Equal(e3v) {
			t.Errorf("Expected e1v to not equal e3v, but they are equal")
		}
	})

	t.Run("view", func(t *testing.T) {
		e := &Example{String: "a", Int: 1}
		v := e.View()
		if v.String() != "a" {
			t.Errorf("Expected String() to return 'a', got '%s'", v.String())
		}
		if v.Int() != 1 {
			t.Errorf("Expected Int() to return 1, got %d", v.Int())
		}

		e.Int = 2
		if v.Int() != 2 {
			t.Errorf("Expected Int() to return 2 after modification, got %d", v.Int())
		}
	})

	t.Run("combine", func(t *testing.T) {
		e := &Example{String: "a", Int: 1}
		if e.Combine() != "a1" {
			t.Errorf("Expected Combine() to return 'a1', got '%s'", e.Combine())
		}
		v := e.View()
		if v.Combine() != "a1" {
			t.Errorf("Expected Combine() to return 'a1', got '%s'", v.Combine())
		}
	})
}

func TestExampleSlice(t *testing.T) {
	t.Run("implements", func(t *testing.T) {
		implements[*ExampleSlice, lego.Equaler[*ExampleSlice]](t)
		implements[*ExampleSlice, lego.Viewer[ExampleSliceView]](t)
		implements[*ExampleSlice, lego.LenLister[lego.Pair[int, *Example]]](t)
		implements[*ExampleSlice, lego.FixedSlice[*Example]](t)

		implements[ExampleSliceView, lego.Equaler[ExampleSliceView]](t)
		implements[ExampleSliceView, lego.FixedSlice[ExampleView]](t)
		implements[ExampleSliceView, lego.LenLister[lego.Pair[int, ExampleView]]](t)
	})

	t.Run("sum", func(t *testing.T) {
		slice := NewExampleSlice(
			&Example{String: "a", Int: 1},
			&Example{String: "b", Int: 2},
			&Example{String: "c", Int: 3},
		)
		if slice.Sum() != 6 {
			t.Errorf("Expected Sum() to return 6, got %d", slice.Sum())
		}
		if slice.View().Sum() != 6 {
			t.Errorf("Expected Sum() to return 6, got %d", slice.View().Sum())
		}
	})

	t.Run("sort", func(t *testing.T) {
		e1 := &Example{String: "a", Int: 2}
		e2 := &Example{String: "a", Int: 1}
		e3 := &Example{String: "c", Int: 3}
		slice := NewExampleSlice(e1, e2, e3)
		lego.Sort(slice.Slice)
		if !slice.Equal(NewExampleSlice(e2, e1, e3)) {
			t.Errorf("Expected slice to be sorted by String then Int, but it is not")
		}
	})
}

type ExampleSliceInt struct {
	*lego.Slice[int]
}

func (s *ExampleSliceInt) View() ExampleSliceIntView {
	return ExampleSliceIntView{s.Slice}
}

func (s *ExampleSliceInt) Equal(other *ExampleSliceInt) bool {
	return lego.EqualViewer(s, other)
}

type ExampleSliceIntView struct {
	lego.FixedSlice[int]
}

func (v ExampleSliceIntView) Equal(other ExampleSliceIntView) bool {
	return lego.EqualSliceGo(v, other)
}

func TestExampleSliceInt(t *testing.T) {
	s1 := &ExampleSliceInt{lego.NewSlice(1, 2, 3)}
	s2 := &ExampleSliceInt{lego.NewSlice(1, 2, 3)}
	s3 := &ExampleSliceInt{lego.NewSlice(1, 2, 4)}

	if !s1.Equal(s2) {
		t.Errorf("Expected s1 to equal s2, but they are not equal")
	}
	if s1.Equal(s3) {
		t.Errorf("Expected s1 to not equal s3, but they are equal")
	}

	vs1 := s1.View()
	vs2 := s2.View()
	vs3 := s3.View()

	if !vs1.Equal(vs2) {
		t.Errorf("Expected vs1 to equal vs2, but they are not equal")
	}
	if vs1.Equal(vs3) {
		t.Errorf("Expected vs1 to not equal vs3, but they are equal")
	}
}
