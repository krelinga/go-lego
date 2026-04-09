package pod

type CanAdd[V any] interface {
	Add(V)
}

type CanReserve interface {
	Reserve(int)
}

func AddAll[V any](a CanAdd[V], r Range[V]) {
	if res, ok := a.(CanReserve); ok {
		res.Reserve(r.Length())
	}
	for v := range r.All() {
		a.Add(v)
	}
}
