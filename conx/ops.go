package conx

type CanLen interface {
	Len() int
}

type CanGetIdx[T any] interface {
	GetIdx(index int) T
}

type CanSetIdx[T any] interface {
	SetIdx(index int, value T)
}

type CanClear interface {
	Clear()
}

type CanAppend[T any] interface {
	Append(value T)
}

type CanReserve interface {
	Reserve(capacity int)
}

type CanFind[K, V any] interface {
	Find(key K) (V, bool)
}

type CanSet[K, V any] interface {
	Set(key K, value V)
}