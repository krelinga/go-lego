// Package exam provides a lightweight test assertion library built on top of the standard [testing]
// package.
//
// Assertions are expressed as functions that return a [Failure] value: nil means the assertion
// passed, a non-nil value means it failed. Failures are reported by passing them to [Try] or
// [Must]:
//
//   - [Try] records a non-fatal error via t.Error and continues the test.
//   - [Must] records a fatal error via t.Fatal and stops the test immediately.
//
// Assertion functions ([Equal], [Greater], [Nil], …) are generic where possible. Custom equality
// or ordering logic can be supplied via the corresponding *Func variants ([EqualFunc],
// [NotEqualFunc], …).
//
// The package also supports table-driven tests through [Here] and [Run]. [Here] captures the source
// location of a test-case struct literal; [Run] logs that source text at the start of each subtest
// so that failures can be traced back to the originating data row.
//
// Golden-file testing is provided via [GoldenHere] and [GoldenEqual]. When the
// -exam_goldens_diff_path flag is set (typically by the update_goldens binary), mismatches are
// written to a diff file for bulk updating instead of failing the test.
package exam

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

// T is the subset of [testing.T] that exam assertion functions accept. *testing.T satisfies T
// directly. This interface allows for tests that use exam assertions to be unit-tested with a fake
// implementation of T.
type T interface {
	Helper()
	Error(args ...any)
	Fatal(args ...any)
	Run(name string, f func(t *testing.T)) bool
}

type parsedFile struct {
	// It is not safe to concurrently access the same *ast.File, even for seemingly read-only
	// operations, so we use a mutex to ensure that only one goroutine is accessing it at a time.
	mu    sync.Mutex
	fset  *token.FileSet
	file  *ast.File
	lines []string
}

var (
	parsedFileCache   = map[string]*parsedFile{}
	parsedFileCacheMu sync.Mutex
)

func getParsedFile(path string) (*parsedFile, func(), error) {
	parsedFileCacheMu.Lock()
	defer parsedFileCacheMu.Unlock()

	if pf, ok := parsedFileCache[path]; ok {
		pf.mu.Lock()
		return pf, pf.mu.Unlock, nil
	}

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	pf := &parsedFile{fset: fset, file: astFile, lines: lines}
	parsedFileCache[path] = pf
	pf.mu.Lock()
	return pf, pf.mu.Unlock, nil
}

func handleResult(t T, must bool, failure *Failure, extras ...any) bool {
	t.Helper()
	if failure == nil {
		return true
	}

	loc := here(2)
	context, err := loadAssertionContext(loc)
	if err != nil {
		context = "unknown assertion context"
	}
	fatal := must || failure.IsCritical()
	sb := &strings.Builder{}
	if fatal {
		sb.WriteString("FATAL: ")
	}
	sb.WriteString(context)
	if failure.Wrapped != nil {
		fmt.Fprintf(sb, "\nSTRUCTURAL ERROR: %v", failure.Wrapped)
	}
	for i, arg := range failure.Args {
		argStr, err := arg.ToString()
		if err != nil {
			fmt.Fprintf(sb, "\nERROR FORMATTING ARG %d (%s): %v", i, arg.Name, err)
		} else {
			fmt.Fprintf(sb, "\narg %d %s", i, argStr)
		}
	}
	for i := range extras {
		fmt.Fprintf(sb, "\n%v", extras[i])
	}
	if fatal {
		t.Fatal(sb.String())
	} else {
		t.Error(sb.String())
	}
	return false
}

// Try reports a non-fatal test error via t.Error when failure is non-nil. The failure message
// includes the source text of the assertion call site; any extras are appended to it. Returns true
// when failure is nil.
func Try(t T, failure *Failure, extras ...any) bool {
	t.Helper()
	return handleResult(t, false, failure, extras...)
}

// Must reports a fatal test error via t.Fatal when failure is non-nil, stopping the test
// immediately. Failures that wrap an underlying error (see Failure.IsCritical) are also fatal when
// passed to Try. The failure message includes the source text of the assertion call site; any
// extras are appended to it. Returns true when failure is nil.
func Must(t T, failure *Failure, extras ...any) bool {
	t.Helper()
	return handleResult(t, true, failure, extras...)
}

// Loc identifies a specific location in a source file. Obtain a Loc at a particular call site using
// Here().
type Loc struct {
	// File is the absolute path to the source file.
	File string
	// Line is the 1-based line number within File.
	Line int
}

// String returns the location as "filename.go:line", using only the base name of File.
func (l Loc) String() string {
	basename := filepath.Base(l.File)
	return fmt.Sprintf("%s:%d", basename, l.Line)
}

// Here returns a Loc for the source line where Here is called.
func Here() Loc {
	return here(1)
}

func here(n int) Loc {
	_, file, line, _ := runtime.Caller(n + 1)
	return Loc{File: file, Line: line}
}

func loadAssertionContext(loc Loc) (string, error) {
	pf, unlock, err := getParsedFile(loc.File)
	if err != nil {
		return "", err
	}
	defer unlock()
	fset := pf.fset

	var found *ast.CallExpr
	ast.Inspect(pf.file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if fset.Position(call.Pos()).Line != loc.Line {
			return true
		}
		// prefer the outermost (largest span) call starting on this line
		if found == nil {
			found = call
		} else {
			foundEnd := fset.Position(found.End()).Line
			callEnd := fset.Position(call.End()).Line
			if callEnd > foundEnd {
				found = call
			}
		}
		return true
	})

	if found == nil {
		return "", fmt.Errorf("no call expression found at line %d in %s", loc.Line, loc.File)
	}

	startLine := fset.Position(found.Pos()).Line
	endLine := fset.Position(found.End()).Line
	lines := pf.lines[startLine-1 : endLine]

	// find minimum indentation among non-empty lines
	minIndent := -1
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if trimmed == "" {
			continue
		}
		indent := len(line) - len(trimmed)
		if minIndent < 0 || indent < minIndent {
			minIndent = indent
		}
	}
	if minIndent < 0 {
		minIndent = 0
	}

	var sb strings.Builder
	for i, line := range lines {
		if len(line) >= minIndent {
			line = line[minIndent:]
		}
		if i > 0 && strings.HasPrefix(line, "\t") {
			// Go's builtin testing framework will automatically indent subsequent lines in log messages, so
			// we trim one level of indentation if it's present to avoid double-indenting.
			line = line[1:]
		}
		sb.WriteString(strings.TrimRight(line, " \t"))
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String()), nil
}

// Run executes f as a named subtest via t.Run. loc must be obtained by calling Here() on the same
// source line as the table-driven test case struct literal; its source text is logged at the start
// of the subtest so that failures can be traced back to the originating data row. Returns the bool
// result of t.Run.
func Run(t T, name string, loc Loc, f func(*testing.T)) bool {
	t.Helper()
	return t.Run(name, func(t *testing.T) {
		t.Helper()
		context, err := loadTableContext(loc)
		if err != nil {
			context = "unknown table context"
		}

		t.Logf("using data from %s\n%s", loc, context)
		f(t)
	})
}

func loadTableContext(loc Loc) (string, error) {
	pf, unlock, err := getParsedFile(loc.File)
	if err != nil {
		return "", err
	}
	defer unlock()
	fset := pf.fset

	var found *ast.CompositeLit
	ast.Inspect(pf.file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		start := fset.Position(lit.Pos()).Line
		end := fset.Position(lit.End()).Line
		if loc.Line >= start && loc.Line <= end {
			// prefer the innermost (smallest span) literal containing the line
			if found == nil {
				found = lit
			} else {
				foundStart := fset.Position(found.Pos()).Line
				foundEnd := fset.Position(found.End()).Line
				if (end - start) < (foundEnd - foundStart) {
					found = lit
				}
			}
		}
		return true
	})

	if found == nil {
		return "", fmt.Errorf("no struct literal found containing line %d in %s", loc.Line, loc.File)
	}

	startLine := fset.Position(found.Pos()).Line
	endLine := fset.Position(found.End()).Line
	lines := pf.lines[startLine-1 : endLine]

	// find minimum indentation among non-empty lines
	minIndent := -1
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if trimmed == "" {
			continue
		}
		indent := len(line) - len(trimmed)
		if minIndent < 0 || indent < minIndent {
			minIndent = indent
		}
	}
	if minIndent < 0 {
		minIndent = 0
	}

	var sb strings.Builder
	for _, line := range lines {
		if len(line) >= minIndent {
			line = line[minIndent:]
		}
		sb.WriteString(strings.TrimRight(line, " \t"))
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String()), nil
}
