package orderexam

import (
	"fmt"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/mirror"
	"github.com/krelinga/go-lego/order"
)

func callOrderFunc[T any](x, y T, order any, want func(int) bool) *exam.Failure {
	failure := exam.NewFailure2("x", x, "y", y)
	wrapped, err := mirror.WrapFunc2In1Out[T, T, int](order)
	if err != nil {
		return failure.Wrap(fmt.Errorf("failed to wrap order function: %w", err))
	} else if want(wrapped(x, y)) {
		return nil
	}
	return failure
}

func GreaterFunc[T any](x, y T, f any) *exam.Failure {
	return callOrderFunc(x, y, f, order.Greater)
}

func LessFunc[T any](x, y T, f any) *exam.Failure {
	return callOrderFunc(x, y, f, order.Less)
}

func EqualFunc[T any](x, y T, f any) *exam.Failure {
	return callOrderFunc(x, y, f, order.Equal)
}

func LessEqualFunc[T any](x, y T, f any) *exam.Failure {
	return callOrderFunc(x, y, f, order.LessEqual)
}

func GreaterEqualFunc[T any](x, y T, f any) *exam.Failure {
	return callOrderFunc(x, y, f, order.GreaterEqual)
}
