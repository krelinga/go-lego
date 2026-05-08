package mirror

import "reflect"

// ValueFor returns a reflect.Value for the given value v.
//
// It is similar to reflect.ValueOf, with some subtle differences:
//   - Interface types are not stripped away, so if v is an interface value then ValueFor will
//     return a reflect.Value of the interface type, not the concrete type.
//   - nil values of interface types will return a valid reflect.Value, and reflect.Value.IsNil()
//     will return true. In contrast, reflect.ValueOf(nil) returns an invalid reflect.Value, and
//     calling IsNil() on it will panic.
//
// This function is useful in generic settings where we can capture information about the type of v.
func ValueFor[T any](v T) reflect.Value {
	typ := reflect.TypeFor[T]()
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Interface {
		if !val.IsValid() {
			return reflect.Zero(typ)
		} else {
			n := reflect.New(typ)
			n.Elem().Set(val)
			return n.Elem()
		}
	}
	return val
}
