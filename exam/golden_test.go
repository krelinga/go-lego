package exam_test

import (
	"path/filepath"
	"testing"

	"github.com/krelinga/go-lego/exam"
)

func TestGoldenHere(t *testing.T) {
	const goldenLine = 12
	golden := exam.GoldenHere(`
test`)
	if filepath.Base(golden.GetLoc().File) != "golden_test.go" {
		t.Errorf("unexpected file: %s", golden.GetLoc().File)
	}
	if golden.GetLoc().Line != goldenLine {
		t.Errorf("unexpected line: %d", golden.GetLoc().Line)
	}
}