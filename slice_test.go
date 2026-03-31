package lego_test

import (
	"testing"

	"github.com/krelinga/go-lego"
)

func TestSliceImplements(t *testing.T) {
	implements[*lego.Slice[string], lego.FixedSlice[string]](t)
	implements[*lego.Slice[string], lego.LenLister[lego.Pair[int, string]]](t)
	implements[*lego.Slice[string], lego.Adder[string]](t)
}

func TestSlice(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var slice *lego.Slice[string]
		panics(t, func() { slice.Len() })
		panics(t, func() { slice.Get(0) })
		panics(t, func() {
			for range slice.List() {
			}
		})
		panics(t, func() { slice.Reserve(1) })
		panics(t, func() { slice.Add("a") })
	})

	t.Run("empty", func(t *testing.T) {
		slice := &lego.Slice[string]{}
		if slice.Len() != 0 {
			t.Errorf("Expected length 0, got %d", slice.Len())
		}
		panics(t, func() { slice.Get(0) })
		for range slice.List() {
			t.Errorf("Expected no elements, but got some")
		}
	})

	t.Run("add_to_empty", func(t *testing.T) {
		slice := &lego.Slice[string]{}
		slice.Add("a")
		if slice.Len() != 1 {
			t.Errorf("Expected length 1, got %d", slice.Len())
		}
		if val := slice.Get(0); val != "a" {
			t.Errorf("Expected Get(0) to return 'a', got '%s'", val)
		}
		panics(t, func() { slice.Get(-1) })
		panics(t, func() { slice.Get(1) })
		for p := range slice.List() {
			i, x := p.GetKey(), p.GetValue()
			if i != 0 {
				t.Errorf("Expected only one element at index 0, got index %d", i)
			}
			if x != "a" {
				t.Errorf("Expected element to be 'a', got '%s'", x)
			}
		}
	})

	t.Run("new", func(t *testing.T) {
		slice := lego.NewSlice[string](3)
		if slice.Len() != 3 {
			t.Errorf("Expected length 3, got %d", slice.Len())
		}
		for i := 0; i < 3; i++ {
			if val := slice.Get(i); val != "" {
				t.Errorf("Expected Get(%d) to return '', got '%s'", i, val)
			}
		}
		panics(t, func() { slice.Get(-1) })
		panics(t, func() { slice.Get(3) })
	})

	t.Run("add_to_non_empty", func(t *testing.T) {
		slice := &lego.Slice[string]{"a"}
		slice.Add("b")
		if slice.Len() != 2 {
			t.Errorf("Expected length 2, got %d", slice.Len())
		}
		if val := slice.Get(0); val != "a" {
			t.Errorf("Expected Get(0) to return 'a', got '%s'", val)
		}
		if val := slice.Get(1); val != "b" {
			t.Errorf("Expected Get(1) to return 'b', got '%s'", val)
		}
		panics(t, func() { slice.Get(-1) })
		panics(t, func() { slice.Get(2) })
		for p := range slice.List() {
			i, x := p.GetKey(), p.GetValue()
			switch i {
			case 0:
				if x != "a" {
					t.Errorf("Expected first element to be 'a', got '%s'", x)
				}
			case 1:
				if x != "b" {
					t.Errorf("Expected second element to be 'b', got '%s'", x)
				}
			default:
				t.Errorf("Expected only 2 elements, got more: %d", i+1)
			}
		}
	})
}

func TestViewSlice(t *testing.T) {
	slice := &ExampleSlice{
		Slice: lego.Slice[*Example]{
			&Example{
				String: "a",
				Int:    1,
			},
			&Example{
				String: "b",
				Int:    2,
			},
		},
	}
	slice.Add(&Example{
		String: "c",
		Int:    3,
	})
	var view lego.FixedSlice[ExampleView] = lego.ViewSlice(slice)
	if view.Len() != 3 {
		t.Errorf("Expected length 3, got %d", view.Len())
	}
	if val := view.Get(0); val.String() != "a" || val.Int() != 1 {
		t.Errorf("Expected first element to be ('a', 1), got ('%s', %d)", val.String(), val.Int())
	}
	if val := view.Get(1); val.String() != "b" || val.Int() != 2 {
		t.Errorf("Expected second element to be ('b', 2), got ('%s', %d)", val.String(), val.Int())
	}
	if val := view.Get(2); val.String() != "c" || val.Int() != 3 {
		t.Errorf("Expected third element to be ('c', 3), got ('%s', %d)", val.String(), val.Int())
	}
	panics(t, func() { view.Get(-1) })
	panics(t, func() { view.Get(3) })
	for p := range view.List() {
		i, x := p.GetKey(), p.GetValue()
		switch i {
		case 0:
			if x.String() != "a" || x.Int() != 1 {
				t.Errorf("Expected first element to be ('a', 1), got ('%s', %d)", x.String(), x.Int())
			}
		case 1:
			if x.String() != "b" || x.Int() != 2 {
				t.Errorf("Expected second element to be ('b', 2), got ('%s', %d)", x.String(), x.Int())
			}
		case 2:
			if x.String() != "c" || x.Int() != 3 {
				t.Errorf("Expected third element to be ('c', 3), got ('%s', %d)", x.String(), x.Int())
			}
		default:
			t.Errorf("Expected only 3 elements, got more: %d", i+1)
		}
	}
}

func TestEqualSlice(t *testing.T) {
	s1 := &ExampleSlice{lego.Slice[*Example]{
		&Example{String: "a", Int: 1},
		&Example{String: "b", Int: 2},
	}}
	s2 := &ExampleSlice{lego.Slice[*Example]{
		&Example{String: "a", Int: 1},
		&Example{String: "b", Int: 2},
	}}
	s3 := &ExampleSlice{lego.Slice[*Example]{
		&Example{String: "a", Int: 1},
		&Example{String: "b", Int: 3},
	}}

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

	if !vs1.Equal(s2.View()) {
		t.Errorf("Expected vs1 to equal s2.View(), but they are not equal")
	}
	if vs1.Equal(s3.View()) {
		t.Errorf("Expected vs1 to not equal s3.View(), but they are equal")
	}
}