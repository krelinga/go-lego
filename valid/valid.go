// Package valid provides a simple interface for validating types and a helper function to enforce
// validation.
package valid

// Validator is an interface that requires implementing types to provide a Validate method that
// checks their validity.
type Validator interface {
	// Validate checks the validity of the implementing type and returns an error if it is invalid.
	Validate() error
}

// Must takes a Validator and panics if the validation fails. It returns the original Validator if
// it is valid. This function is useful for enforcing validation in situations where you want to
// ensure that a value is valid before proceeding, such as during initialization or when
// constructing complex types.
func Must[T Validator](v T) T {
	if err := v.Validate(); err != nil {
		panic(err)
	}
	return v
}
