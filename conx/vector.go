package conx

type Vector[T any] struct {
	data []T
}

func NewVector[T any](data ...T) *Vector[T] {
	return &Vector[T]{data: data}
}

func (v *Vector[T]) Len() int {
	return len(v.data)
}

func (v *Vector[T]) GetIdx(index int) T {
	return v.data[index]
}

func (v *Vector[T]) SetIdx(index int, value T) {
	v.data[index] = value
}

func (v *Vector[T]) Clear() {
	v.data = nil
}

func (v *Vector[T]) Append(value T) {
	v.data = append(v.data, value)
}

func (v *Vector[T]) Reserve(capacity int) {
	if cap(v.data) < capacity {
		newData := make([]T, len(v.data), capacity)
		copy(newData, v.data)
		v.data = newData
	}
}