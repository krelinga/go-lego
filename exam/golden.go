package exam

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/krelinga/go-lego/exam/internal"
)

type Golden struct {
	loc  Loc
	text string
}

func (g Golden) GetLoc() Loc {
	return g.loc
}

func GoldenHere(text string) Golden {
	return Golden{
		loc:  here(1),
		text: text,
	}
}

// StringIndentLines is a FmtFunc that formats a string value by indenting each line with a tab.  It returns an error if the value is not a string.
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
