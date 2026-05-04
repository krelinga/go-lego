package podexam

import (
	"fmt"

	"github.com/krelinga/go-libs/exam"
	"github.com/krelinga/go-libs/mirror"
	"github.com/krelinga/go-libs/pod"
)

func DictEqualFunc[K comparable, V any](x, y pod.DictView[K, V], f any) *exam.Failure {
	failure := exam.NewFailure2("x", x, "y", y)
	valueEqual, err := mirror.WrapFunc2In1Out[V, V, bool](f)
	if err != nil {
		return failure.Wrap(fmt.Errorf("failed to wrap value equality function: %w", err))
	} else if pod.DictEqualFunc(x, y, valueEqual) {
		return nil
	}
	return failure
}

func VecEqualFunc[T any](x, y pod.VecView[T], f any) *exam.Failure {
	failure := exam.NewFailure2("x", x, "y", y)
	valueEqual, err := mirror.WrapFunc2In1Out[T, T, bool](f)
	if err != nil {
		return failure.Wrap(fmt.Errorf("failed to wrap value equality function: %w", err))
	} else if pod.VecEqualFunc(x, y, valueEqual) {
		return nil
	}
	return failure
}
