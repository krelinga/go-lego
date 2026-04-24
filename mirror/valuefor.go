package mirror

import "reflect"

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
	return reflect.ValueOf(v)
}
