package exam

import (
	"fmt"

	"github.com/krelinga/go-libs/mirror"
	"github.com/krelinga/go-libs/order"
)

func callOrderFunc[T any](x, y T, orderFn any, want func(int) bool) *Failure {
	failure := NewFailure2("x", x, "y", y)
	wrapped, err := mirror.WrapFunc2In1Out[T, T, int](orderFn)
	if err != nil {
		return failure.Wrap(fmt.Errorf("failed to wrap order function: %w", err))
	} else if want(wrapped(x, y)) {
		return nil
	}
	return failure
}

func GreaterFunc[T any](x, y T, f any) *Failure {
	return callOrderFunc(x, y, f, order.Greater)
}

func LessFunc[T any](x, y T, f any) *Failure {
	return callOrderFunc(x, y, f, order.Less)
}

func OrderEqualFunc[T any](x, y T, f any) *Failure {
	return callOrderFunc(x, y, f, order.Equal)
}

func LessEqualFunc[T any](x, y T, f any) *Failure {
	return callOrderFunc(x, y, f, order.LessEqual)
}

func GreaterEqualFunc[T any](x, y T, f any) *Failure {
	return callOrderFunc(x, y, f, order.GreaterEqual)
}
