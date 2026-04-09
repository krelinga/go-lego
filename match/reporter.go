package match

import "reflect"

type Reporter interface {
	TypeMismatch(expected, actual reflect.Type)
	ValueMismatch(expected, actual any)

	Struct(reflect.Type) (StructReporter, func())
	Container(reflect.Type) (ContainerReporter, func())
	Pointer(reflect.Type) (PointerReporter, func())
}

type StructReporter interface {
	Embedded(reflect.StructField) (Reporter, func())
	Field(reflect.StructField) (Reporter, func())
	Method(reflect.Method) (Reporter, func())
}

type MethodReporter interface {
	NumOutsMismatch(expected, actual int)

	Out(int) (Reporter, func())
}

type ContainerReporter interface {
	LengthMismatch(expected, actual int)
	Extra(any)
	Missing(any)

	Element() (Reporter, func())
}

type PointerReporter interface {
	AddressMismatch(expected, actual uintptr)
	ExpectedNil()
	ActualNil()

	Dereference() (Reporter, func())
}