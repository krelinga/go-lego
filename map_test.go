package lego_test

import (
	"testing"

	"github.com/krelinga/go-lego"
)

func TestMapImplements(t *testing.T) {
	implements[*lego.Map[string, int], lego.FixedMap[string, int]](t)
	implements[*lego.Map[string, int], lego.LenLister[lego.Pair[string, int]]](t)
	implements[*lego.Map[string, int], lego.Adder[lego.Pair[string, int]]](t)
	
	implements[lego.GoMap[string, int], lego.FixedMap[string, int]](t)
	implements[lego.GoMap[string, int], lego.LenLister[lego.Pair[string, int]]](t)
}

func TestGoMap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var m lego.GoMap[string, int]
		if m.Len() != 0 {
			t.Errorf("Expected length 0, got %d", m.Len())
		}
		if _, ok := m.Get("a"); ok {
			t.Errorf("Expected Get('a') to return false, got true")
		}
		for range m.List() {
			t.Errorf("Expected no elements in the map, got some")
		}
	})

	t.Run("non_nil", func(t *testing.T) {
		m := lego.GoMap[string, int]{"a": 1, "b": 2, "c": 3}
		if m.Len() != 3 {
			t.Errorf("Expected length 3, got %d", m.Len())
		}
		if v, ok := m.Get("a"); !ok || v != 1 {
			t.Errorf("Expected Get('a') to return (1, true), got (%d, %t)", v, ok)
		}
		if v, ok := m.Get("b"); !ok || v != 2 {
			t.Errorf("Expected Get('b') to return (2, true), got (%d, %t)", v, ok)
		}
		if v, ok := m.Get("c"); !ok || v != 3 {
			t.Errorf("Expected Get('c') to return (3, true), got (%d, %t)", v, ok)
		}
		if _, ok := m.Get("d"); ok {
			t.Errorf("Expected Get('d') to return false, got true")
		}
		for pair := range m.List() {
			k, v := pair.GetKey(), pair.GetValue()
			switch k {
			case "a":
				if v != 1 {
					t.Errorf("Expected value for 'a' to be 1, got %d", v)
				}
			case "b":
				if v != 2 {
					t.Errorf("Expected value for 'b' to be 2, got %d", v)
				}
			case "c":
				if v != 3 {
					t.Errorf("Expected value for 'c' to be 3, got %d", v)
				}
			default:
				t.Errorf("Expected only keys 'a', 'b', and 'c' in the map, got '%s'", k)
			}
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var m *lego.Map[string, int]
		panics(t, func() { m.Len() })
		panics(t, func() { m.Get("a") })
		panics(t, func() {
			for range m.List() {
			}
		})
		panics(t, func() { m.Reserve(1) })
		panics(t, func() { m.Add(lego.NewPair("a", 1)) })
		panics(t, func() { m.Insert("a", 1) })
	})

	t.Run("empty", func(t *testing.T) {
		m := &lego.Map[string, int]{}
		if m.Len() != 0 {
			t.Errorf("Expected length 0, got %d", m.Len())
		}
		if _, ok := m.Get("a"); ok {
			t.Errorf("Expected Get('a') to return false, got true")
		}
		for range m.List() {
			t.Errorf("Expected no elements in the map, got some")
		}
	})

	t.Run("add_to_empty", func(t *testing.T) {
		m := &lego.Map[string, int]{}
		m.Add(lego.NewPair("a", 1))
		if m.Len() != 1 {
			t.Errorf("Expected length 1, got %d", m.Len())
		}
		if v, ok := m.Get("a"); !ok || v != 1 {
			t.Errorf("Expected Get('a') to return (1, true), got (%d, %t)", v, ok)
		}
		if _, ok := m.Get("b"); ok {
			t.Errorf("Expected Get('b') to return false, got true")
		}
		for pair := range m.List() {
			k, v := pair.GetKey(), pair.GetValue()
			if k != "a" || v != 1 {
				t.Errorf("Expected only key 'a' with value 1 in the map, got key '%s' with value %d", k, v)
			}
		}
	})

	t.Run("insert_to_empty", func(t *testing.T) {
		m := &lego.Map[string, int]{}
		m.Insert("a", 1)
		if m.Len() != 1 {
			t.Errorf("Expected length 1, got %d", m.Len())
		}
		if v, ok := m.Get("a"); !ok || v != 1 {
			t.Errorf("Expected Get('a') to return (1, true), got (%d, %t)", v, ok)
		}
		if _, ok := m.Get("b"); ok {
			t.Errorf("Expected Get('b') to return false, got true")
		}
		for pair := range m.List() {
			k, v := pair.GetKey(), pair.GetValue()
			if k != "a" || v != 1 {
				t.Errorf("Expected only key 'a' with value 1 in the map, got key '%s' with value %d", k, v)
			}
		}
	})

	t.Run("add_to_non_empty", func(t *testing.T) {
		m := lego.NewMap(lego.NewPair("a", 1))
		m.Add(lego.NewPair("b", 2))
		if m.Len() != 2 {
			t.Errorf("Expected length 2, got %d", m.Len())
		}
		if v, ok := m.Get("a"); !ok || v != 1 {
			t.Errorf("Expected Get('a') to return (1, true), got (%d, %t)", v, ok)
		}
		if v, ok := m.Get("b"); !ok || v != 2 {
			t.Errorf("Expected Get('b') to return (2, true), got (%d, %t)", v, ok)
		}
		if _, ok := m.Get("c"); ok {
			t.Errorf("Expected Get('c') to return false, got true")
		}
		for pair := range m.List() {
			k, v := pair.GetKey(), pair.GetValue()
			switch k {
			case "a":
				if v != 1 {
					t.Errorf("Expected value for 'a' to be 1, got %d", v)
				}
			case "b":
				if v != 2 {
					t.Errorf("Expected value for 'b' to be 2, got %d", v)
				}
			default:
				t.Errorf("Expected only keys 'a' and 'b' in the map, got '%s'", k)
			}
		}
	})

	t.Run("insert_to_non_empty", func(t *testing.T) {
		m := lego.NewMap(lego.NewPair("a", 1))
		m.Insert("b", 2)
		if m.Len() != 2 {
			t.Errorf("Expected length 2, got %d", m.Len())
		}
		if v, ok := m.Get("a"); !ok || v != 1 {
			t.Errorf("Expected Get('a') to return (1, true), got (%d, %t)", v, ok)
		}
		if v, ok := m.Get("b"); !ok || v != 2 {
			t.Errorf("Expected Get('b') to return (2, true), got (%d, %t)", v, ok)
		}
		if _, ok := m.Get("c"); ok {
			t.Errorf("Expected Get('c') to return false, got true")
		}
		for pair := range m.List() {
			k, v := pair.GetKey(), pair.GetValue()
			switch k {
			case "a":
				if v != 1 {
					t.Errorf("Expected value for 'a' to be 1, got %d", v)
				}
			case "b":
				if v != 2 {
					t.Errorf("Expected value for 'b' to be 2, got %d", v)
				}
			default:
				t.Errorf("Expected only keys 'a' and 'b' in the map, got '%s'", k)
			}
		}
	})
}