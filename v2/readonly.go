package v2

type ReadonlyFunc[I, O any] func(I) O

func Identity[I any](i I) I {
	return i
}

type Readonlyer[O any] interface {
	Readonly() O
}

func Readonly[I Readonlyer[O], O any](i I) O {
	return i.Readonly()
}