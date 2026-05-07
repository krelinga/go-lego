package pod

import (
	"errors"

	"github.com/krelinga/go-libs/valid"
)

func ValidateAll[T valid.Validator](vals Bag[T]) error {
	var errs []error
	for v := range vals.Elems() {
		errs = append(errs, v.Validate())
	}
	return errors.Join(errs...)
}