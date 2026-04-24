package orderexam

import (
	"fmt"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/mirror"
	"github.com/krelinga/go-lego/order"
)

func callOrderFunc[T any](x, y T, order any, want func(int) bool) (result exam.Result) {
	result = exam.Result{
		Args: []any{x, y},
	}
	wrapped, err := mirror.WrapFunc2In1Out[T, T, int](order)
	if err != nil {
		result.Error = fmt.Errorf("failed to call order function: %w", err)
		return
	}
	if !want(wrapped(x, y)) {
		result.Error = exam.ErrFailed
	}
	return
}

func GreaterFunc[T any](x, y T, f any) exam.Result {
	return callOrderFunc(x, y, f, order.Greater)
}

func LessFunc[T any](x, y T, f any) exam.Result {
	return callOrderFunc(x, y, f, order.Less)
}

func EqualFunc[T any](x, y T, f any) exam.Result {
	return callOrderFunc(x, y, f, order.Equal)
}

func LessEqualFunc[T any](x, y T, f any) exam.Result {
	return callOrderFunc(x, y, f, order.LessEqual)
}

func GreaterEqualFunc[T any](x, y T, f any) exam.Result {
	return callOrderFunc(x, y, f, order.GreaterEqual)
}
