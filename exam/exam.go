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
	fset  *token.FileSet
	file  *ast.File
	lines []string
}

var (
	parsedFileCache   = map[string]*parsedFile{}
	parsedFileCacheMu sync.Mutex
)

func getParsedFile(path string) (*parsedFile, error) {
	parsedFileCacheMu.Lock()
	defer parsedFileCacheMu.Unlock()

	if pf, ok := parsedFileCache[path]; ok {
		return pf, nil
	}

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	pf := &parsedFile{fset: fset, file: astFile, lines: lines}
	parsedFileCache[path] = pf
	return pf, nil
}

func errorPrefix(sb *strings.Builder, fatal bool, context string) {
	if fatal {
		sb.WriteString("FATAL: ")
	}
	sb.WriteString(context)
}

func makeError(fatal bool, context string, args ...any) string {
	sb := &strings.Builder{}
	errorPrefix(sb, fatal, context)
	if len(args) > 0 {
		sb.WriteString("\n")
		fmt.Fprint(sb, args...)
	}
	return sb.String()
}

func makeErrorf(fatal bool, context string, format string, args ...any) string {
	sb := &strings.Builder{}
	errorPrefix(sb, fatal, context)
	if len(args) > 0 {
		sb.WriteString("\n")
		fmt.Fprintf(sb, format, args...)
	}
	return sb.String()
}

func Try(t *testing.T, condition bool, args ...any) bool {
	t.Helper()
	if !condition {
		loc := here(1)
		context, err := loadAssertionContext(loc)
		if err != nil {
			context = "unknown assertion context"
		}
		t.Error(makeError(false, context, args...))
	}
	return condition
}

func Tryf(t *testing.T, condition bool, format string, args ...any) bool {
	t.Helper()
	if !condition {
		loc := here(1)
		context, err := loadAssertionContext(loc)
		if err != nil {
			context = "unknown assertion context"
		}
		t.Error(makeErrorf(false, context, format, args...))
	}
	return condition
}

func Must(t *testing.T, condition bool, args ...any) bool {
	t.Helper()
	if !condition {
		loc := here(1)
		context, err := loadAssertionContext(loc)
		if err != nil {
			context = "unknown assertion context"
		}
		t.Fatal(makeError(true, context, args...))
	}
	return condition
}

func Mustf(t *testing.T, condition bool, format string, args ...any) bool {
	t.Helper()
	if !condition {
		loc := here(1)
		context, err := loadAssertionContext(loc)
		if err != nil {
			context = "unknown assertion context"
		}
		t.Fatal(makeErrorf(true, context, format, args...))
	}
	return condition
}

type Loc struct {
	File string
	Line int
}

func (l Loc) String() string {
	basename := filepath.Base(l.File)
	return fmt.Sprintf("%s:%d", basename, l.Line)
}

func Here() Loc {
	return here(1)
}

func here(n int) Loc {
	_, file, line, _ := runtime.Caller(n + 1)
	return Loc{File: file, Line: line}
}

func loadAssertionContext(loc Loc) (string, error) {
	pf, err := getParsedFile(loc.File)
	if err != nil {
		return "", err
	}
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
	pf, err := getParsedFile(loc.File)
	if err != nil {
		return "", err
	}
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
