package valid

type Validator interface {
	// Validate checks the validity of the implementing type and returns an error if it is invalid.
	Validate() error
}

func Must[T Validator](v T) T {
	if err := v.Validate(); err != nil {
		panic(err)
	}
	return v
}
