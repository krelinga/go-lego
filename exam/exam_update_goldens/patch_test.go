package main

import (
	"path/filepath"
	"testing"

	"github.com/krelinga/go-lego/exam"
	"github.com/krelinga/go-lego/exam/internal"
)

func TestFindLineStart(t *testing.T) {
	cases := []struct {
		Name    string
		Loc     exam.Loc
		Src     string
		LineNum int
		WantPos int
		WantErr bool
	}{
		{
			Name:    "line 1 of single-line src",
			Loc:     exam.Here(),
			Src:     "hello",
			LineNum: 1,
			WantPos: 0,
		},
		{
			Name:    "line 1 of multi-line src",
			Loc:     exam.Here(),
			Src:     "alpha\nbeta\ngamma",
			LineNum: 1,
			WantPos: 0,
		},
		{
			Name:    "line 2",
			Loc:     exam.Here(),
			Src:     "alpha\nbeta\ngamma",
			LineNum: 2,
			WantPos: 6,
		},
		{
			Name:    "line 3",
			Loc:     exam.Here(),
			Src:     "alpha\nbeta\ngamma",
			LineNum: 3,
			WantPos: 11,
		},
		{
			Name:    "line beyond end",
			Loc:     exam.Here(),
			Src:     "alpha\nbeta",
			LineNum: 3,
			WantErr: true,
		},
		{
			Name:    "empty src line 1",
			Loc:     exam.Here(),
			Src:     "",
			LineNum: 1,
			WantPos: 0,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			pos, err := findLineStart(c.Src, c.LineNum)
			if c.WantErr {
				exam.Try(t, exam.NotNil(err))
			} else {
				exam.Must(t, exam.Nil(err))
				exam.Try(t, exam.Equal(pos, c.WantPos))
			}
		})
	}
}

// srcSingleGolden is test source with one GoldenHere call.
// The opening backtick is on line 4.
const srcSingleGolden = `package test

func TestFoo(t *testing.T) {
	exam.GoldenEqual(s, exam.GoldenHere(` + "`" + `
old content
` + "`" + `))
}`

// srcTwoGoldens is test source with two GoldenHere calls.
// The first opening backtick is on line 4, the second on line 8.
const srcTwoGoldens = `package test

func TestFoo(t *testing.T) {
	exam.GoldenEqual(s, exam.GoldenHere(` + "`" + `
old content 1
` + "`" + `))
	exam.GoldenEqual(s2, exam.GoldenHere(` + "`" + `
old content 2
` + "`" + `))
}`

func TestApplyDiffsToSrc(t *testing.T) {
	cases := []struct {
		Name    string
		Loc     exam.Loc
		Src     string
		Entries []internal.GoldenEntry
		Want    string
		WantErr bool
	}{
		{
			// Replace the content of a single golden with text of the same line count.
			Name:    "single patch same line count",
			Loc:     exam.Here(),
			Src:     srcSingleGolden,
			Entries: []internal.GoldenEntry{{Line: 4, Text: "\nnew content\n"}},
			Want: `package test

func TestFoo(t *testing.T) {
	exam.GoldenEqual(s, exam.GoldenHere(` + "`" + `
new content
` + "`" + `))
}`,
		},
		{
			// Replace content with more lines; the closing backtick shifts down.
			Name: "single patch grows",
			Loc:  exam.Here(),
			Src:  srcSingleGolden,
			Entries: []internal.GoldenEntry{
				{Line: 4, Text: "\nline A\nline B\nline C\n"},
			},
			Want: `package test

func TestFoo(t *testing.T) {
	exam.GoldenEqual(s, exam.GoldenHere(` + "`" + `
line A
line B
line C
` + "`" + `))
}`,
		},
		{
			// Replace content with fewer lines; the closing backtick shifts up.
			Name:    "single patch shrinks",
			Loc:     exam.Here(),
			Src:     srcSingleGolden,
			Entries: []internal.GoldenEntry{{Line: 4, Text: "\n"}},
			Want: `package test

func TestFoo(t *testing.T) {
	exam.GoldenEqual(s, exam.GoldenHere(` + "`" + `
` + "`" + `))
}`,
		},
		{
			// Two patches with the same line count: no lineOffset adjustment needed.
			Name: "two patches same line count",
			Loc:  exam.Here(),
			Src:  srcTwoGoldens,
			Entries: []internal.GoldenEntry{
				{Line: 4, Text: "\nnew 1\n"},
				{Line: 7, Text: "\nnew 2\n"},
			},
			Want: `package test

func TestFoo(t *testing.T) {
	exam.GoldenEqual(s, exam.GoldenHere(` + "`" + `
new 1
` + "`" + `))
	exam.GoldenEqual(s2, exam.GoldenHere(` + "`" + `
new 2
` + "`" + `))
}`,
		},
		{
			// First patch adds 2 lines; the second entry's line must be shifted by +2.
			Name: "two patches first grows",
			Loc:  exam.Here(),
			Src:  srcTwoGoldens,
			Entries: []internal.GoldenEntry{
				{Line: 4, Text: "\nline A\nline B\nline C\n"},
				// In the original source this GoldenHere is on line 7.
				// After the first patch inserts 2 extra lines it becomes line 9.
				{Line: 7, Text: "\nnew 2\n"},
			},
			Want: `package test

func TestFoo(t *testing.T) {
	exam.GoldenEqual(s, exam.GoldenHere(` + "`" + `
line A
line B
line C
` + "`" + `))
	exam.GoldenEqual(s2, exam.GoldenHere(` + "`" + `
new 2
` + "`" + `))
}`,
		},
		{
			// First patch removes 1 line; the second entry's line is shifted by -1.
			Name: "two patches first shrinks",
			Loc:  exam.Here(),
			Src:  srcTwoGoldens,
			Entries: []internal.GoldenEntry{
				{Line: 4, Text: "\n"},
				{Line: 7, Text: "\nnew 2\n"},
			},
			Want: `package test

func TestFoo(t *testing.T) {
	exam.GoldenEqual(s, exam.GoldenHere(` + "`" + `
` + "`" + `))
	exam.GoldenEqual(s2, exam.GoldenHere(` + "`" + `
new 2
` + "`" + `))
}`,
		},
		{
			Name:    "no entries returns src unchanged",
			Loc:     exam.Here(),
			Src:     srcSingleGolden,
			Entries: nil,
			Want:    srcSingleGolden,
		},
		{
			Name:    "line beyond end of src",
			Loc:     exam.Here(),
			Src:     "one line",
			Entries: []internal.GoldenEntry{{Line: 5, Text: "\nnew\n"}},
			WantErr: true,
		},
		{
			Name:    "no opening backtick on target line",
			Loc:     exam.Here(),
			Src:     "no backtick here\n",
			Entries: []internal.GoldenEntry{{Line: 1, Text: "\nnew\n"}},
			WantErr: true,
		},
		{
			Name:    "no closing backtick",
			Loc:     exam.Here(),
			Src:     "f(`unclosed",
			Entries: []internal.GoldenEntry{{Line: 1, Text: "\nnew\n"}},
			WantErr: true,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			got, err := applyDiffsToSrc(c.Src, c.Entries)
			if c.WantErr {
				exam.Try(t, exam.NotNil(err))
			} else {
				exam.Must(t, exam.Nil(err))
				exam.Try(t, exam.Equal(got, c.Want))
			}
		})
	}
}

func TestValidateEntryPath(t *testing.T) {
	root := "/workspace/myproject"
	cases := []struct {
		Name    string
		Loc     exam.Loc
		Path    string
		WantErr bool
	}{
		{
			Name: "valid go file in workspace",
			Loc:  exam.Here(),
			Path: filepath.Join(root, "pkg", "foo_test.go"),
		},
		{
			Name: "valid go file at workspace root",
			Loc:  exam.Here(),
			Path: filepath.Join(root, "main.go"),
		},
		{
			Name:    "not a go file - txt extension",
			Loc:     exam.Here(),
			Path:    filepath.Join(root, "pkg", "data.txt"),
			WantErr: true,
		},
		{
			Name:    "not a go file - no extension",
			Loc:     exam.Here(),
			Path:    filepath.Join(root, "pkg", "somefile"),
			WantErr: true,
		},
		{
			Name:    "outside workspace - sibling directory",
			Loc:     exam.Here(),
			Path:    "/workspace/otherproject/foo.go",
			WantErr: true,
		},
		{
			Name:    "outside workspace - system go files",
			Loc:     exam.Here(),
			Path:    "/usr/local/go/src/fmt/print.go",
			WantErr: true,
		},
		{
			Name:    "path traversal attempt",
			Loc:     exam.Here(),
			Path:    filepath.Clean(filepath.Join(root, "..", "..", "etc", "passwd.go")),
			WantErr: true,
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			err := validateEntryPath(c.Path, root)
			if c.WantErr {
				exam.Try(t, exam.NotNil(err))
			} else {
				exam.Try(t, exam.Nil(err))
			}
		})
	}
}
