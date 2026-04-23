package exam

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func errorPrefix(sb *strings.Builder, fatal bool, context string) {
	if fatal {
		sb.WriteString("FATAL: ")
	}
	sb.WriteString(context)
}

func makeError(fatal bool, context string, args ... any) string {
	sb := &strings.Builder{}
	errorPrefix(sb, fatal, context)
	if len(args) > 0 {
		sb.WriteString("\n")
		fmt.Fprint(sb, args...)
	}
	return sb.String()
}

func makeErrorf(fatal bool, context string, format string, args ... any) string {
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
	f, err := os.Open(loc.File)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for i := 1; scanner.Scan(); i++ {
		if i == loc.Line {
			return strings.TrimSpace(scanner.Text()), nil
		}
	}
	return "", fmt.Errorf("line %d not found in file %s", loc.Line, loc.File)
}

func Run(t *testing.T, name string, loc Loc, f func(*testing.T)) bool {
	return t.Run(name, func(t *testing.T) {
		t.Helper()
		t.Log("using data from", loc)
		f(t)
	})
}