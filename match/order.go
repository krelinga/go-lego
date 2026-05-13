package match

import (
	"cmp"
	"fmt"
	"reflect"
)

type OrderOp struct {
	op     func(int) bool
	str    string
	invStr string
}

func (o *OrderOp) checkInit() error {
	if o.op == nil {
		return fmt.Errorf("OrderOp must be created with OrderOpLt, OrderOpGt, OrderOpLte, OrderOpGte, or OrderOpEq")
	}
	return nil
}

func (o *OrderOp) evaluate(i int) (bool, error) {
	if err := o.checkInit(); err != nil {
		return false, err
	}
	return o.op(i), nil
}

func OrderOpLt() *OrderOp {
	return &OrderOp{
		op: func(i int) bool {
			return i < 0
		},
		str:    "<",
		invStr: ">=",
	}
}

func OrderOpGt() *OrderOp {
	return &OrderOp{
		op: func(i int) bool {
			return i > 0
		},
		str:    ">",
		invStr: "<=",
	}
}

func OrderOpLte() *OrderOp {
	return &OrderOp{
		op: func(i int) bool {
			return i <= 0
		},
		str:    "<=",
		invStr: ">",
	}
}

func OrderOpGte() *OrderOp {
	return &OrderOp{
		op: func(i int) bool {
			return i >= 0
		},
		str:    ">=",
		invStr: "<=",
	}
}

func OrderOpEq() *OrderOp {
	return &OrderOp{
		op: func(i int) bool {
			return i == 0
		},
		str:    "==",
		invStr: "!=",
	}
}

type OrderFunc struct {
	f any
}

func NewOrderFunc[T any](f func(T, T) int) *OrderFunc {
	if f == nil {
		panic("match.NewOrderFunc must be called with a non-nil function")
	}
	return &OrderFunc{f: f}
}

func (o *OrderFunc) checkInit() error {
	if o.f == nil {
		return fmt.Errorf("match.OrderFunc must be created with match.NewOrderFunc")
	}
	return nil
}

func (o *OrderFunc) checkType(t reflect.Type) error {
	if err := o.checkInit(); err != nil {
		return err
	}
	wantType := reflect.TypeOf(o.f).In(0)
	if !t.AssignableTo(wantType) {
		return fmt.Errorf("order function expects a type assignable to %s but got type %s", wantType, t)
	}
	return nil
}

func (o *OrderFunc) order(a, b any) (int, error) {
	if err := o.checkType(reflect.TypeOf(a)); err != nil {
		return 0, err
	}
	if err := o.checkType(reflect.TypeOf(b)); err != nil {
		return 0, err
	}
	fVal := reflect.ValueOf(o.f)
	result := fVal.Call([]reflect.Value{reflect.ValueOf(a), reflect.ValueOf(b)})
	return int(result[0].Int()), nil
}

type Order struct {
	Limit  any
	Op     *OrderOp
	Func   *OrderFunc
	Format *Format
}

func (o *Order) Validate() error {
	if o.Limit == nil {
		return fmt.Errorf("order limit must be specified")
	}
	if o.Op == nil {
		return fmt.Errorf("order operator must be specified")
	}
	if o.Format != nil {
		if err := o.Format.checkType(reflect.TypeOf(o.Limit)); err != nil {
			return err
		}
	}
	if o.Func != nil {
		if err := o.Func.checkType(reflect.TypeOf(o.Limit)); err != nil {
			return err
		}
	} else {
		// If no custom order function is provided, we can only support types that are ordered by Go's built-in comparison operators.
		limitType := reflect.TypeOf(o.Limit)
		switch limitType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.String:
			// These types are supported by Go's built-in comparison operators.
		default:
			return fmt.Errorf("type %s does not support ordering and no custom order function was provided", limitType)
		}
	}
	return nil
}

func (o *Order) Match(val any) (*Result, error) {
	h := &Helper{
		Meta: MetaHere(),
		Val:  val,
	}
	h.Context("limit", o.Limit)
	if err := o.Validate(); err != nil {
		return nil, h.Fatal(err)
	}
	var orderResult int
	if o.Func != nil {
		var err error
		orderResult, err = o.Func.order(val, o.Limit)
		if err != nil {
			return nil, h.Fatal(err)
		}
	} else {
		switch v := val.(type) {
		case int:
			orderResult = cmp.Compare(v, o.Limit.(int))
		case int8:
			orderResult = cmp.Compare(v, o.Limit.(int8))
		case int16:
			orderResult = cmp.Compare(v, o.Limit.(int16))
		case int32:
			orderResult = cmp.Compare(v, o.Limit.(int32))
		case int64:
			orderResult = cmp.Compare(v, o.Limit.(int64))
		case uint:
			orderResult = cmp.Compare(v, o.Limit.(uint))
		case uint8:
			orderResult = cmp.Compare(v, o.Limit.(uint8))
		case uint16:
			orderResult = cmp.Compare(v, o.Limit.(uint16))
		case uint32:
			orderResult = cmp.Compare(v, o.Limit.(uint32))
		case uint64:
			orderResult = cmp.Compare(v, o.Limit.(uint64))
		case float32:
			orderResult = cmp.Compare(v, o.Limit.(float32))
		case float64:
			orderResult = cmp.Compare(v, o.Limit.(float64))
		case string:
			orderResult = cmp.Compare(v, o.Limit.(string))
		default:
			return nil, h.Fatalf("unsupported type %T", val)
		}
	}
	if ok, err := o.Op.evaluate(orderResult); err != nil {
		return nil, h.Fatal(err)
	} else if ok {
		return h.Accept(fmt.Sprintf("value %s limit", o.Op.str)), nil
	} else {
		return h.Reject(fmt.Sprintf("value %s limit", o.Op.invStr)), nil
	}
}

func OrderCmpLt[T comparable](limit T) *Order {
	return &Order{
		Limit: limit,
		Op:    OrderOpLt(),
	}
}

func OrderCmpGt[T comparable](limit T) *Order {
	return &Order{
		Limit: limit,
		Op:    OrderOpGt(),
	}
}

func OrderCmpLte[T comparable](limit T) *Order {
	return &Order{
		Limit: limit,
		Op:    OrderOpLte(),
	}
}

func OrderCmpGte[T comparable](limit T) *Order {
	return &Order{
		Limit: limit,
		Op:    OrderOpGte(),
	}
}

func OrderCmpEq[T comparable](limit T) *Order {
	return &Order{
		Limit: limit,
		Op:    OrderOpEq(),
	}
}
