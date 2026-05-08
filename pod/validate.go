package pod

import (
	"errors"

	"github.com/krelinga/go-libs/valid"
)

// ValidateAll takes a Bag of Validators and returns an error if any of the Validators are invalid.
// It collects all validation errors and returns them as a single error using errors.Join.
func ValidateAll[T valid.Validator](vals Bag[T]) error {
	var errs []error
	for v := range vals.Elems() {
		errs = append(errs, v.Validate())
	}
	return errors.Join(errs...)
}
