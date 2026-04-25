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

func handleResult(t *testing.T, must bool, failure *Failure, extra string) bool {
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
		fmt.Fprintf(sb, "\narg %d %s", i, arg)
	}
	if extra != "" {
		fmt.Fprintf(sb, "\n%s", extra)
	}
	if fatal {
		t.Fatal(sb.String())
	} else {
		t.Error(sb.String())
	}
	return false
}

// Try checks condition and, if false, records a non-fatal test failure using t.Error.
// The failure message includes the source text of the call site and any optional args.
// Returns true if the condition passed, false otherwise.
func Try(t *testing.T, failure *Failure, args ...any) bool {
	t.Helper()
	return handleResult(t, false, failure, fmt.Sprint(args...))
}

// Tryf checks condition and, if false, records a non-fatal test failure using t.Error.
// The failure message includes the source text of the call site and a formatted message
// constructed from format and args. Returns true if the condition passed, false otherwise.
func Tryf(t *testing.T, failure *Failure, format string, args ...any) bool {
	t.Helper()
	return handleResult(t, false, failure, fmt.Sprintf(format, args...))
}

// Must checks condition and, if false, records a fatal test failure using t.Fatal,
// stopping the test immediately. The failure message includes the source text of the
// call site and any optional args. Returns true if the condition passed, false otherwise.
func Must(t *testing.T, failure *Failure, args ...any) bool {
	t.Helper()
	return handleResult(t, true, failure, fmt.Sprint(args...))
}

// Mustf checks condition and, if false, records a fatal test failure using t.Fatal,
// stopping the test immediately. The failure message includes the source text of the
// call site and a formatted message constructed from format and args.
// Returns true if the condition passed, false otherwise.
func Mustf(t *testing.T, failure *Failure, format string, args ...any) bool {
	t.Helper()
	return handleResult(t, true, failure, fmt.Sprintf(format, args...))
}

// Loc identifies a specific location in a source file by file path and line number.
// It is typically obtained by calling Here() at the point of interest.
type Loc struct {
	File string
	Line int
}

// String returns a short human-readable representation of the location in the form
// "filename.go:line".
func (l Loc) String() string {
	basename := filepath.Base(l.File)
	return fmt.Sprintf("%s:%d", basename, l.Line)
}

// Here returns a Loc representing the source location of the caller.
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
			// Go's builtin testing framework will automatically indent subsequent lines in log messages,
			// so we trim one level of indentation if it's present to avoid double-indenting.
			line = line[1:]
		}
		sb.WriteString(strings.TrimRight(line, " \t"))
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String()), nil
}

// Run executes a named subtest using t.Run and logs the source text of the table-driven
// test case identified by loc. loc should be obtained by calling Here() on the same line
// as the struct literal that defines the test case data.
func Run(t *testing.T, name string, loc Loc, f func(*testing.T)) bool {
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
