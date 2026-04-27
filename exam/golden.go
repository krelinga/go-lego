package exam

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/krelinga/go-lego/exam/internal"
)

// Golden holds an expected string and the source location of the GoldenHere
// call that produced it. Use GoldenHere to create one at a call site.
type Golden struct {
	loc  Loc
	text string
}

// GetLoc returns the source location of the GoldenHere call that created g.
func (g Golden) GetLoc() Loc {
	return g.loc
}

// TODO is a placeholder value for golden data that has not yet been generated.
// Its value is "\n" — a single newline — which satisfies the requirement that
// golden text must begin with a newline while being short enough that
// GoldenEqual will always record a mismatch and prompt the tool to fill in the
// real content.  Pass it to GoldenHere as a placeholder until the golden data
// can be filled in, for example by running the exam_update_goldens tool.
const TODO = "\n"

// GoldenHere creates a Golden whose expected text is text and whose source
// location is the line where GoldenHere is called. text must begin with a
// newline; by convention the closing backtick of the raw string literal
// appears on its own line so that the content is unambiguous. As a
// convenience, exam.TODO may be passed as a placeholder until real golden
// data is available.
func GoldenHere(text string) Golden {
	return Golden{
		loc:  here(1),
		text: text,
	}
}

// StringIndentLines is a FmtFunc that inserts a tab after every newline in
// val's string representation, making multi-line values readable in failure
// output. It returns an error if val is not of kind string.
func StringIndentLines(val reflect.Value) (string, error) {
	if !val.IsValid() {
		return "", fmt.Errorf("invalid value")
	} else if val.Kind() != reflect.String {
		return "", fmt.Errorf("value is of kind %s, expected string", val.Kind())
	}
	valStr := fmt.Sprintf("%s", val.String())
	return strings.ReplaceAll(valStr, "\n", "\n\t"), nil
}

var examGoldensDiffPath = flag.String("exam_goldens_diff_path", "", "Path to write golden diffs to.  If set, the golden data assertions will always pass.")
var examGoldensMu = sync.Mutex{}

// GoldenEqual asserts that actual matches the expected text stored in expected.
// actual is prepended with a newline before comparison so that multi-line
// values align naturally with the raw string literal in the test source.
//
// When -exam_goldens_diff_path is set, any mismatch is recorded to that file
// for later bulk-updating and the assertion always passes. Without the flag,
// a mismatch returns a Failure.  Instead of using this flag directly, prefer
// running tests with the exam_update_goldens tool, which sets the flag and
// applies the recorded diffs after the test run completes.
func GoldenEqual(actual string, expected Golden) *Failure {
	actual = "\n" + actual
	failure := NewFailure2("actual", actual, "expected", expected.text)
	for i := range failure.Args {
		failure.Args[i].Fmt = StringIndentLines
	}

	if len(expected.text) == 0 || expected.text[0] != '\n' {
		return failure.Wrap(fmt.Errorf("expected text must start with a newline"))
	}

	if *examGoldensDiffPath != "" {
		if actual == expected.text {
			return nil
		}

		examGoldensMu.Lock()
		defer examGoldensMu.Unlock()
		if err := internal.WriteGoldenEntry(*examGoldensDiffPath, internal.GoldenEntry{
			Path: expected.loc.File,
			Line: expected.loc.Line,
			Text: actual,
		}); err != nil {
			return failure.Wrap(fmt.Errorf("writing golden entry: %w", err))
		}
		return nil
	}

	if actual == expected.text {
		return nil
	}
	return failure
}
