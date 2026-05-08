package exam

import (
	"strings"

	"github.com/krelinga/go-libs/valid"
)

// MustValidate takes a Validator and checks if it is valid. If the Validator is valid, it returns
// the original value. If the Validator is invalid, it reports a fatal error using the provided
// testing interface, including information about the validation error and the context of the
// assertion.
func MustValidate[V valid.Validator](t T, v V) V {
	if validErr := v.Validate(); validErr != nil {
		b := &strings.Builder{}
		context, err := loadAssertionContext(here(1))
		if err != nil {
			context = "unknown assertion context"
		}
		b.WriteString("FATAL: ")
		b.WriteString(context)
		b.WriteString("\n")
		b.WriteString("validation failed: ")
		b.WriteString(validErr.Error())
		t.Fatal(b.String())
	}
	return v
}
