package v2

type Adder[V any] interface {
	Add(V)
}

type Reserver interface {
	Reserve(int)
}

func AddAll[V any](a Adder[V], r Range[V]) {
	if res, ok := a.(Reserver); ok {
		res.Reserve(r.Length())
	}
	for v := range r.All() {
		a.Add(v)
	}
}

func AddAllSlice[V any](a Adder[V], vs []V) {
	slice := Slice[V](vs)
	AddAll(a, ValuesFrom(&slice))
}
