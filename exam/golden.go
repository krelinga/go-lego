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

var examUpdateGoldens = flag.Bool("exam_update_goldens", false, "update golden data instead of comparing it")

func GoldenEqual(actual string, expected Golden) *Failure {
	expectedText := expected.text
	if len(expectedText) > 0 && expectedText[0] == '\n' {
		expectedText = expectedText[1:]
	}

	if *examUpdateGoldens {
		if actual == expectedText {
			return nil
		}
		// TODO: implement updating the golden data in the source file at expected.loc with actual
		return nil
	}
	return nil // TODO: implement comparison and failure reporting
}