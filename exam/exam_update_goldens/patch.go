package main

import (
	"fmt"
	"strings"

	"github.com/krelinga/go-lego/exam/internal"
)

// applyDiffsToSrc applies a set of golden-entry patches to the Go source text
// src and returns the updated text.  entries must be sorted by Line ascending
// and must all refer to line numbers within src (1-indexed).
//
// Each entry locates the first backtick on the entry's source line and
// replaces the raw-string literal content between that backtick and the next
// one in the file with entry.Text.  A running lineOffset tracks how prior
// patches shift the line numbers of later entries.
func applyDiffsToSrc(src string, entries []internal.GoldenEntry) (string, error) {
	// lineOffset accumulates the net change in line count from patches applied
	// so far, so that later entries (whose Line fields refer to the original
	// source) can be located correctly in the already-patched text.
	lineOffset := 0

	for _, entry := range entries {
		targetLine := entry.Line + lineOffset // 1-indexed within current src

		lineStart, err := findLineStart(src, targetLine)
		if err != nil {
			return "", fmt.Errorf("locating line %d: %w", targetLine, err)
		}

		// The opening backtick of the GoldenHere() raw-string argument must
		// appear at or after the GoldenHere( token on this line.
		relOpen := strings.Index(src[lineStart:], "`")
		if relOpen < 0 {
			return "", fmt.Errorf("no opening backtick found on line %d", targetLine)
		}
		openContent := lineStart + relOpen + 1 // byte position just after the opening `

		// Find the matching closing backtick.
		relClose := strings.Index(src[openContent:], "`")
		if relClose < 0 {
			return "", fmt.Errorf("no closing backtick found for GoldenHere on line %d", targetLine)
		}
		closeBacktick := openContent + relClose // byte position of the closing `

		// Update the running offset by the net line-count delta of this patch.
		oldContent := src[openContent:closeBacktick]
		lineOffset += strings.Count(entry.Text, "\n") - strings.Count(oldContent, "\n")

		// Splice in the new content, leaving both backticks in place.
		src = src[:openContent] + entry.Text + src[closeBacktick:]
	}

	return src, nil
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
