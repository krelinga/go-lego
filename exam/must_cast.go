package exam

import (
	"fmt"
	"reflect"
	"strings"
)

// MustCast attempts to cast inVal value to the Out type. If the cast is successful, it returns the casted value. If the cast fails, it reports a fatal error using the provided testing interface, including information about the expected and actual types, as well as the context of the assertion.
func MustCast[Out, In any](t T, inVal In) Out {
	t.Helper()
	value, ok := any(inVal).(Out)
	if !ok {
		inName := reflect.TypeFor[In]().String()
		outName := reflect.TypeFor[Out]().String()
		context, err := loadAssertionContext(here(1))
		if err != nil {
			context = "unknown assertion context"
		}
		sb := &strings.Builder{}
		fmt.Fprintf(sb, "FATAL: %s\n", context)
		fmt.Fprintf(sb, "expected value that could be cast to type %s, got %s", outName, inName)
		t.Fatal(sb.String())
	}
	return value
}
