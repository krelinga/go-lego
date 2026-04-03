package v2

type Adder[V any] interface {
	Add(V)
}

type Reserver interface {
	Reserve(int)
}

func AddAll[V any](a Adder[V], vs ...V) {
	if r, ok := a.(Reserver); ok {
		r.Reserve(len(vs))
	}
	for _, v := range vs {
		a.Add(v)
	}
}
