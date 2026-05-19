package match

import (
	"cmp"
	"fmt"
)

type OrderMatcher[T, TT any] struct {
	meta Meta
	limit T
	limitFmt Fmt[T]
	gotFmt Fmt[any]
	order func(TT, TT) int
	check func(int) bool
	acceptStr string
	rejectStr string
	valid bool
}

func (o *OrderMatcher[T, TT]) Match(val any) (*Result, error) {
	h := &Helper{
		Meta: o.meta,
		Val:  val,
	}
	if !o.valid {
		return nil, h.Fatalf("OrderMatcher must be created with OrderLt OrderLtFunc, OrderGt, OrderGtFunc, OrderLte, OrderLteFunc, OrderGte, OrderGteFunc, OrderEq, or OrderEqFunc")
	}
	h.Context("limit", o.limit)
	limitTT, ok := any(o.limit).(TT)
	if !ok {
		return nil, h.Fatalf("limit of type %T is not assignable to order function type %T", o.limit, limitTT)
	}
	gotTT, err := As[TT](h, val)
	if err != nil {
		return nil, err
	}
	result := o.order(gotTT, limitTT)
	if o.check(result) {
		return h.Accept(fmt.Sprintf("got %s limit", o.acceptStr)), nil
	} else {
		return h.Reject(fmt.Sprintf("got %s limit", o.rejectStr)), nil
	}
}

func (o *OrderMatcher[T, TT]) LimitFmt(t Fmt[T]) *OrderMatcher[T, TT] {
	o.limitFmt = t
	return o
}

func (o *OrderMatcher[T, TT]) GotFmt(t Fmt[any]) *OrderMatcher[T, TT] {
	o.gotFmt = t
	return o
}

func OrderLt[T cmp.Ordered](limit T) *OrderMatcher[T, T] {
	return &OrderMatcher[T, T]{
		meta:      MetaHere(),
		limit:     limit,
		order:     cmp.Compare[T],
		check:     func(i int) bool { return i < 0 },
		acceptStr: "<",
		rejectStr: ">=",
		valid:     true,
	}
}

func OrderLtFunc[T, TT any](limit T, orderFunc func(TT, TT) int) *OrderMatcher[T, TT] {
	return &OrderMatcher[T, TT]{
		meta:      MetaHere(),
		limit:     limit,
		order:     orderFunc,
		check:     func(i int) bool { return i < 0 },
		acceptStr: "<",
		rejectStr: ">=",
		valid:     true,
	}
}

func OrderGt[T cmp.Ordered](limit T) *OrderMatcher[T, T] {
	return &OrderMatcher[T, T]{
		meta:      MetaHere(),
		limit:     limit,
		order:     cmp.Compare[T],
		check:     func(i int) bool { return i > 0 },
		acceptStr: ">",
		rejectStr: "<=",
		valid:     true,
	}
}

func OrderGtFunc[T, TT any](limit T, orderFunc func(TT, TT) int) *OrderMatcher[T, TT] {
	return &OrderMatcher[T, TT]{
		meta:      MetaHere(),
		limit:     limit,
		order:     orderFunc,
		check:     func(i int) bool { return i > 0 },
		acceptStr: ">",
		rejectStr: "<=",
		valid:     true,
	}
}

func OrderLte[T cmp.Ordered](limit T) *OrderMatcher[T, T] {
	return &OrderMatcher[T, T]{
		meta:      MetaHere(),
		limit:     limit,
		order:     cmp.Compare[T],
		check:     func(i int) bool { return i <= 0 },
		acceptStr: "<=",
		rejectStr: ">",
		valid:     true,
	}
}

func OrderLteFunc[T, TT any](limit T, orderFunc func(TT, TT) int) *OrderMatcher[T, TT] {
	return &OrderMatcher[T, TT]{
		meta:      MetaHere(),
		limit:     limit,
		order:     orderFunc,
		check:     func(i int) bool { return i <= 0 },
		acceptStr: "<=",
		rejectStr: ">",
		valid:     true,
	}
}

func OrderGte[T cmp.Ordered](limit T) *OrderMatcher[T, T] {
	return &OrderMatcher[T, T]{
		meta:      MetaHere(),
		limit:     limit,
		order:     cmp.Compare[T],
		check:     func(i int) bool { return i >= 0 },
		acceptStr: ">=",
		rejectStr: "<",
		valid:     true,
	}
}

func OrderGteFunc[T, TT any](limit T, orderFunc func(TT, TT) int) *OrderMatcher[T, TT] {
	return &OrderMatcher[T, TT]{
		meta:      MetaHere(),
		limit:     limit,
		order:     orderFunc,
		check:     func(i int) bool { return i >= 0 },
		acceptStr: ">=",
		rejectStr: "<",
		valid:     true,
	}
}

func OrderEq[T cmp.Ordered](limit T) *OrderMatcher[T, T] {
	return &OrderMatcher[T, T]{
		meta:      MetaHere(),
		limit:     limit,
		order:     cmp.Compare[T],
		check:     func(i int) bool { return i == 0 },
		acceptStr: "==",
		rejectStr: "!=",
		valid:     true,
	}
}

func OrderEqFunc[T, TT any](limit T, orderFunc func(TT, TT) int) *OrderMatcher[T, TT] {
	return &OrderMatcher[T, TT]{
		meta:      MetaHere(),
		limit:     limit,
		order:     orderFunc,
		check:     func(i int) bool { return i == 0 },
		acceptStr: "==",
		rejectStr: "!=",
		valid:     true,
	}
}
