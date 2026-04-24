package podexam

import (
	"fmt"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/mirror"
	"github.com/krelinga/go-lego/pod"
)

func MapEqualFunc[K comparable, V any](x, y pod.MapView[K, V], f any) (result exam.Result) {
	result = exam.Result{
		Args: []any{x, y},
	}
	valueEqual, err := mirror.WrapFunc2In1Out[V, V, bool](f)
	if err != nil {
		result.Error = fmt.Errorf("failed to wrap value equality function: %w", err)
		return
	}
	result.Error = exam.AsError(pod.MapEqualFunc(x, y, valueEqual))
	return
}

func VecEqualFunc[T any](x, y pod.VecView[T], f any) (result exam.Result) {
	result = exam.Result{
		Args: []any{x, y},
	}
	valueEqual, err := mirror.WrapFunc2In1Out[T, T, bool](f)
	if err != nil {
		result.Error = fmt.Errorf("failed to wrap value equality function: %w", err)
		return
	}
	result.Error = exam.AsError(pod.VecEqualFunc(x, y, valueEqual))
	return
}
