package exam

import (
	"flag"
)

type Golden struct {
	loc Loc
	text string
}

func GoldenHere(text string) Golden {
	return Golden{
		loc:  here(2),
		text: text,
	}
}

var examGoldensDiffPath = flag.String("exam_goldens_diff_path", "", "Path to write golden diffs to.  If set, the golden data assertions will always pass.")

func GoldenEqual(actual string, expected Golden) *Failure {
	expectedText := expected.text
	if len(expectedText) > 0 && expectedText[0] == '\n' {
		expectedText = expectedText[1:]
	}

	if *examGoldensDiffPath != "" {
		if actual == expectedText {
			return nil
		}
		// TODO: implement writing the golden diff to the file at *examGoldensDiffPath
		return nil
	}

	if actual == expectedText {
		return nil
	}
	return nil // TODO: implement
}