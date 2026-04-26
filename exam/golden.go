package exam

import (
	"flag"
	"fmt"
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

var examGoldensDiffPath = flag.String("exam_goldens_diff_path", "", "Path to write golden diffs to.  If set, the golden data assertions will always pass.")
var examGoldensMu = sync.Mutex{}

func GoldenEqual(actual string, expected Golden) *Failure {
	actual = "\n" + actual
	failure := NewFailure2("actual", actual, "expected", expected.text)

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
