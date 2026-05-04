package mirror_test

import (
	"reflect"
	"testing"

	"github.com/krelinga/go-libs/exam"
	"github.com/krelinga/go-libs/mirror"
)

func TestValueFor(t *testing.T) {
	t.Run("nil interface", func(t *testing.T) {
		var i any
		v := mirror.ValueFor(i)
		exam.Must(t, exam.True(v.IsValid()))
		i2 := v.Interface()
		exam.Try(t, exam.Nil(i2))
		exam.Try(t, exam.Equal(v.Type(), reflect.TypeFor[any]()))
	})

	t.Run("non-nil interface", func(t *testing.T) {
		var i any = int(42)
		v := mirror.ValueFor(i)
		exam.Must(t, exam.True(v.IsValid()))
		i2 := v.Interface()
		exam.Try(t, exam.Equal(i2, any(int(42))))
		exam.Try(t, exam.Equal(v.Type(), reflect.TypeFor[any]()))
	})

	t.Run("non-interface", func(t *testing.T) {
		x := int(42)
		v := mirror.ValueFor(x)
		exam.Must(t, exam.True(v.IsValid()))
		exam.Try(t, exam.Equal(v.Interface(), any(int(42))))
		exam.Try(t, exam.Equal(v.Type(), reflect.TypeFor[int]()))
	})
}
