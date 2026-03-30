package lego_test

import (
	"testing"

	"github.com/krelinga/go-lego"
)

func TestSetImplements(t *testing.T) {
	implements[*lego.Set[string], lego.FixedSet[string]](t)
	implements[*lego.Set[string], lego.LenLister[string]](t)
	implements[*lego.Set[string], lego.Adder[string]](t)

	implements[lego.GoSet[string], lego.FixedSet[string]](t)
	implements[lego.GoSet[string], lego.LenLister[string]](t)
}

func TestGoSet(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var set lego.GoSet[string]
		if set.Len() != 0 {
			t.Errorf("Expected length 0, got %d", set.Len())
		}
		if set.Has("a") {
			t.Errorf("Expected Has('a') to return false, got true")
		}
		for range set.List() {
			t.Errorf("Expected no elements in the set, got some")
		}
	})

	t.Run("non_nil", func(t *testing.T) {
		set := lego.GoSet[string]{"a": {}, "b": {}, "c": {}}
		if set.Len() != 3 {
			t.Errorf("Expected length 3, got %d", set.Len())
		}
		if !set.Has("a") {
			t.Errorf("Expected Has('a') to return true, got false")
		}
		if !set.Has("b") {
			t.Errorf("Expected Has('b') to return true, got false")
		}
		if !set.Has("c") {
			t.Errorf("Expected Has('c') to return true, got false")
		}
		if set.Has("d") {
			t.Errorf("Expected Has('d') to return false, got true")
		}
		for v := range set.List() {
			switch v {
			case "a", "b", "c":
				// expected
			default:
				t.Errorf("Expected only 'a', 'b', and 'c' in the set, got '%s'", v)
			}
		}
	})
}

func TestSet(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		set := &lego.Set[string]{}
		if set.Len() != 0 {
			t.Errorf("Expected length 0, got %d", set.Len())
		}
		if set.Has("a") {
			t.Errorf("Expected Has('a') to return false, got true")
		}
		for range set.List() {
			t.Errorf("Expected no elements in the set, got some")
		}
	})

	t.Run("add_to_empty", func(t *testing.T) {
		set := &lego.Set[string]{}
		set.Add("a")
		if set.Len() != 1 {
			t.Errorf("Expected length 1, got %d", set.Len())
		}
		if !set.Has("a") {
			t.Errorf("Expected Has('a') to return true, got false")
		}
		if set.Has("b") {
			t.Errorf("Expected Has('b') to return false, got true")
		}
		for v := range set.List() {
			switch v {
			case "a":
				// expected
			default:
				t.Errorf("Expected only 'a' in the set, got '%s'", v)
			}
		}
	})

	t.Run("add_to_non_empty", func(t *testing.T) {
		set := lego.NewSet("a")
		set.Add("b")
		if set.Len() != 2 {
			t.Errorf("Expected length 2, got %d", set.Len())
		}
		if !set.Has("a") {
			t.Errorf("Expected Has('a') to return true, got false")
		}
		if !set.Has("b") {
			t.Errorf("Expected Has('b') to return true, got false")
		}
		if set.Has("c") {
			t.Errorf("Expected Has('c') to return false, got true")
		}
		for v := range set.List() {
			switch v {
			case "a", "b":
				// expected
			default:
				t.Errorf("Expected only 'a' and 'b' in the set, got '%s'", v)
			}
		}
	})
}
