package mirror_test

import (
	"testing"

	"github.com/krelinga/go-lego/mirror"
)

func TestNil(t *testing.T) {
	t.Run("nil slice", func(t *testing.T) {
		var s []int
		isNil, err := mirror.IsNil(s)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !isNil {
			t.Errorf("expected nil slice to be nil")
		}
	})

	t.Run("non-nil slice", func(t *testing.T) {
		s := []int{}
		isNil, err := mirror.IsNil(s)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if isNil {
			t.Errorf("expected non-nil slice to not be nil")
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		var p *int
		isNil, err := mirror.IsNil(p)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !isNil {
			t.Errorf("expected nil pointer to be nil")
		}
	})

	t.Run("non-nil pointer", func(t *testing.T) {
		x := 42
		p := &x
		isNil, err := mirror.IsNil(p)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if isNil {
			t.Errorf("expected non-nil pointer to not be nil")
		}
	})

	t.Run("non-nil value", func(t *testing.T) {
		x := 42
		_, err := mirror.IsNil(x)
		if err == nil {
			t.Fatalf("expected error for non-nil value")
		}
	})

	t.Run("nil interface", func(t *testing.T) {
		var x any
		isNil, err := mirror.IsNil(x)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !isNil {
			t.Errorf("expected nil interface to be nil")
		}
	})

	t.Run("non-nil interface", func(t *testing.T) {
		var x any = 42
		isNil, err := mirror.IsNil(x)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if isNil {
			t.Errorf("expected non-nil interface to not be nil")
		}
	})

	t.Run("nil function", func(t *testing.T) {
		var f func()
		isNil, err := mirror.IsNil(f)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !isNil {
			t.Errorf("expected nil function to be nil")
		}
	})

	t.Run("non-nil function", func(t *testing.T) {
		f := func() {}
		isNil, err := mirror.IsNil(f)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if isNil {
			t.Errorf("expected non-nil function to not be nil")
		}
	})

	t.Run("nil map", func(t *testing.T) {
		var m map[string]int
		isNil, err := mirror.IsNil(m)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !isNil {
			t.Errorf("expected nil map to be nil")
		}
	})

	t.Run("non-nil map", func(t *testing.T) {
		m := map[string]int{}
		isNil, err := mirror.IsNil(m)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if isNil {
			t.Errorf("expected non-nil map to not be nil")
		}
	})

	t.Run("nil channel", func(t *testing.T) {
		var ch chan int
		isNil, err := mirror.IsNil(ch)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !isNil {
			t.Errorf("expected nil channel to be nil")
		}
	})

	t.Run("non-nil channel", func(t *testing.T) {
		ch := make(chan int)
		isNil, err := mirror.IsNil(ch)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if isNil {
			t.Errorf("expected non-nil channel to not be nil")
		}
	})
}