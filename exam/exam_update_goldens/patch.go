package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/krelinga/go-libs/exam/internal"
)

// validateEntryPath returns an error if path should not be edited.  Two
// checks are applied:
//  1. The file must have a ".go" extension.
//  2. The file must be located within workspaceRoot (i.e. the resolved
//     relative path must not start with "..").
func validateEntryPath(path, workspaceRoot string) error {
	if filepath.Ext(path) != ".go" {
		return fmt.Errorf("refusing to edit %s: not a .go file", path)
	}
	rel, err := filepath.Rel(workspaceRoot, path)
	if err != nil {
		return fmt.Errorf("refusing to edit %s: cannot determine path relative to workspace: %w", path, err)
	}
	if strings.HasPrefix(rel, "..") {
		return fmt.Errorf("refusing to edit %s: not within workspace root %s", path, workspaceRoot)
	}
	return nil
}

// applyDiffsToSrc applies a set of golden-entry patches to the Go source text
// src and returns the updated text.  entries must be sorted by Line ascending
// and must all refer to line numbers within src (1-indexed).
//
// Each entry locates the GoldenHere( call on the entry's source line, finds
// the extent of its single argument (which may be a raw-string literal, an
// interpreted string literal, the exam.TODO constant, or a concatenated
// expression), and replaces that argument with a new Go string expression
// produced by generateStringExpr.  A running lineOffset tracks how prior
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

		// Find "GoldenHere(" on this line.
		relGoldenHere := strings.Index(src[lineStart:], "GoldenHere(")
		if relGoldenHere < 0 {
			return "", fmt.Errorf("no GoldenHere( found on line %d", targetLine)
		}
		// argStart is the byte position of the first byte of the argument.
		argStart := lineStart + relGoldenHere + len("GoldenHere(")

		// Find the matching closing ')' of GoldenHere(...).
		argEnd, err := findArgEnd(src, argStart)
		if err != nil {
			return "", fmt.Errorf("finding argument end for GoldenHere on line %d: %w", targetLine, err)
		}

		oldArg := src[argStart:argEnd]
		newArg := generateStringExpr(entry.Text)

		// Update the running offset by the net line-count delta of this patch.
		lineOffset += strings.Count(newArg, "\n") - strings.Count(oldArg, "\n")

		// Splice in the new argument, leaving the surrounding parentheses in place.
		src = src[:argStart] + newArg + src[argEnd:]
	}

	return src, nil
}

// findArgEnd returns the byte position of the closing ')' that closes the
// GoldenHere( call, given that argStart is the first byte after the opening
// '('.  It properly handles nested parentheses, raw string literals (backtick-
// delimited), interpreted string literals (double-quote-delimited), and rune
// literals (single-quote-delimited).
func findArgEnd(src string, argStart int) (int, error) {
	depth := 0
	i := argStart
	for i < len(src) {
		switch src[i] {
		case '`':
			// Raw string literal: scan forward until the closing backtick.
			i++
			for i < len(src) && src[i] != '`' {
				i++
			}
			if i >= len(src) {
				return 0, fmt.Errorf("unterminated raw string literal")
			}
			i++ // skip closing backtick
		case '"':
			// Interpreted string literal: scan until an unescaped '"'.
			i++
			closed := false
			for i < len(src) {
				if src[i] == '\\' {
					if i+1 < len(src) {
						i += 2
					} else {
						i++
					}
					continue
				}
				if src[i] == '"' {
					i++
					closed = true
					break
				}
				i++
			}
			if !closed {
				return 0, fmt.Errorf("unterminated interpreted string literal")
			}
		case '\'':
			// Rune literal: scan until an unescaped '\''.
			i++
			closed := false
			for i < len(src) {
				if src[i] == '\\' {
					if i+1 < len(src) {
						i += 2
					} else {
						i++
					}
					continue
				}
				if src[i] == '\'' {
					i++
					closed = true
					break
				}
				i++
			}
			if !closed {
				return 0, fmt.Errorf("unterminated rune literal")
			}
		case '(':
			depth++
			i++
		case ')':
			if depth == 0 {
				return i, nil
			}
			depth--
			i++
		default:
			i++
		}
	}
	return 0, fmt.Errorf("no matching closing parenthesis found")
}

// generateStringExpr returns a Go expression that evaluates to text.
// When text contains no backtick characters a single raw string literal is
// returned.  When text does contain backticks the result is a concatenation
// of raw string literals (for the non-backtick segments) and double-quoted
// string literals (for the backtick characters themselves), so that the
// overall expression is valid Go regardless of the content.
func generateStringExpr(text string) string {
	if !strings.ContainsRune(text, '`') {
		return "`" + text + "`"
	}

	parts := strings.Split(text, "`")
	var exprs []string
	for i, part := range parts {
		if part != "" {
			exprs = append(exprs, "`"+part+"`")
		}
		if i < len(parts)-1 {
			// Represent the backtick character as a double-quoted string literal.
			exprs = append(exprs, "\"`\"")
		}
	}
	return strings.Join(exprs, " + ")
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
