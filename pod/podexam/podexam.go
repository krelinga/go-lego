package podexam

import (
	"fmt"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/mirror"
	"github.com/krelinga/go-lego/pod"
)

func MapEqualFunc[K comparable, V any](x, y pod.MapView[K, V], f any) error {
	valueEqual, err := mirror.WrapFunc2In1Out[V, V, bool](f)
	if err != nil {
		return fmt.Errorf("failed to wrap value equality function: %w", err)
	} else if pod.MapEqualFunc(x, y, valueEqual) {
		return nil
	}
	return exam.Failure{
		Args: []any{x, y},
	}
}

func VecEqualFunc[T any](x, y pod.VecView[T], f any) error {
	valueEqual, err := mirror.WrapFunc2In1Out[T, T, bool](f)
	if err != nil {
		return fmt.Errorf("failed to wrap value equality function: %w", err)
	} else if pod.VecEqualFunc(x, y, valueEqual) {
		return nil
	}
	return exam.Failure{
		Args: []any{x, y},
	}
}
