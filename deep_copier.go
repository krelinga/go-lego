package lego

type DeepCopier[T any] interface {
	DeepCopy() T
}