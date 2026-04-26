// exam_update_goldens runs the tests identified by a -run pattern and applies any
// golden-file differences that the tests emit back into the source files.
//
// Usage:
//
//	go tool exam_update_goldens ./...
package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/krelinga/go-lego/exam/internal"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <test-run-pattern>\n", os.Args[0])
		os.Exit(1)
	}
	pattern := os.Args[1]

	// Temporary file that test binaries write golden diffs into.
	tmpFile, err := os.CreateTemp("", "exam_goldens_diff_*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: creating temp file: %v\n", err)
		os.Exit(1)
	}
	tmpFile.Close()
	diffPath := tmpFile.Name()
	defer os.Remove(diffPath)

	// Run the tests.  -p=1 keeps them sequential so multiple test binaries
	// don't race when appending to the shared golden-diff file.
	cmd := exec.Command("go", "test",
		"-count=1",
		"-p=1",
		pattern,
		fmt.Sprintf("-exam_goldens_diff_path=%s", diffPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: go test failed: %v\n", err)
		os.Exit(1)
	}

	// Read golden entries written during the test run.
	entries, err := internal.ReadGoldenEntries(diffPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: reading golden entries: %v\n", err)
		os.Exit(1)
	}
	if len(entries) == 0 {
		fmt.Println("No golden diffs to apply.")
		return
	}

	// Group entries by source file and sort each group by line number so that
	// we can track how inserted/removed lines shift subsequent entries.
	byFile := make(map[string][]internal.GoldenEntry)
	for _, e := range entries {
		byFile[e.Path] = append(byFile[e.Path], e)
	}
	for path, fileEntries := range byFile {
		sort.Slice(fileEntries, func(i, j int) bool {
			return fileEntries[i].Line < fileEntries[j].Line
		})
		if err := applyDiffs(path, fileEntries); err != nil {
			fmt.Fprintf(os.Stderr, "error: applying diffs to %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("updated %s (%d diff(s))\n", path, len(fileEntries))
	}
}

// applyDiffs patches the raw-string literal arguments of GoldenHere() calls in
// path according to entries.  Entries must already be sorted by Line ascending.
func applyDiffs(path string, entries []internal.GoldenEntry) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	src := string(data)
	// lineOffset accumulates the net change in line count from all patches
	// applied so far, so that later entries (whose Line fields refer to the
	// original file) can be located correctly in the already-patched src.
	lineOffset := 0

	for _, entry := range entries {
		targetLine := entry.Line + lineOffset // still 1-indexed

		lineStart, err := findLineStart(src, targetLine)
		if err != nil {
			return fmt.Errorf("locating line %d: %w", targetLine, err)
		}

		// The opening backtick of the raw string literal that is the argument
		// to GoldenHere() must appear at or after the GoldenHere( call on this
		// line.  Find it.
		relOpen := strings.Index(src[lineStart:], "`")
		if relOpen < 0 {
			return fmt.Errorf("no opening backtick found on line %d", targetLine)
		}
		openContent := lineStart + relOpen + 1 // byte position just after the opening `

		// Find the closing backtick.
		relClose := strings.Index(src[openContent:], "`")
		if relClose < 0 {
			return fmt.Errorf("no closing backtick found for GoldenHere on line %d", targetLine)
		}
		closeBacktick := openContent + relClose // byte position of the closing `

		// Measure the line-count delta so we can adjust subsequent entries.
		oldContent := src[openContent:closeBacktick]
		lineOffset += strings.Count(entry.Text, "\n") - strings.Count(oldContent, "\n")

		// Splice in the new text (leave the two backticks in place).
		src = src[:openContent] + entry.Text + src[closeBacktick:]
	}

	return os.WriteFile(path, []byte(src), 0644)
}

// findLineStart returns the byte offset of the first byte on the given
// 1-indexed line within src.
func findLineStart(src string, lineNum int) (int, error) {
	pos := 0
	for line := 1; line < lineNum; line++ {
		idx := strings.Index(src[pos:], "\n")
		if idx < 0 {
			return 0, fmt.Errorf("file has fewer than %d lines", lineNum)
		}
		pos += idx + 1
	}
	return pos, nil
}
