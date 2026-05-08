package wrap_test

import (
	"strings"
	"testing"

	"github.com/krelinga/go-libs/cmd/comment-reflow/internal/wrap"
	"github.com/krelinga/go-libs/exam"
)

func TestWrap(t *testing.T) {
	cases := []struct {
		Name       string
		Loc        exam.Loc
		Limit      int
		Src        string
		WantChange bool
		Golden     exam.Golden
	}{
		{
			Name:  "NoChangeForShortComment",
			Loc:   exam.Here(),
			Limit: 100,
			Src: `
package p

// This comment is short.
func Foo() {}
`,
			WantChange: false,
		},
		{
			Name:  "NoChangeForCommentAtLimit",
			Loc:   exam.Here(),
			Limit: 11, // "// At Limit" is exactly 11 chars
			Src: `
package p

// At Limit
func Foo() {}
`,
			WantChange: false,
		},
		{
			Name:  "SplitCommentOneCharOverLimit",
			Loc:   exam.Here(),
			Limit: 12, // "// Over Limit" is 13 chars
			Src: `
package p

// Over Limit
func Foo() {}
`,
			WantChange: true,
			Golden: exam.GoldenHere(`
package p

// Over
// Limit
func Foo() {}
`),
		},
		{
			Name:  "NoChangeForInlineComment",
			Loc:   exam.Here(),
			Limit: 10,
			Src: `
package p

func Foo() {
	x := 1 // inline comment
}
`,
			WantChange: false,
		},
		{
			Name:  "NoChangeForTestDirective",
			Loc:   exam.Here(),
			Limit: 10,
			Src: `
package p

//go:generate some very long command that exceeds the line length limit by quite a bit here
//nolint:somelinter,anotherlinter,yetanotherlinter // long nolint directive that is over the limit
func Foo() {}
`,
			WantChange: false,
		},
		{
			Name:  "URLNotSplit",
			Loc:   exam.Here(),
			Limit: 60,
			Src: `
package p

// See https://example.com/some/very/long/path/that/pushes/the/line/over/the/limit for details.
func Foo() {}
`,
			WantChange: true,
			Golden: exam.GoldenHere(`
package p

// See
// https://example.com/some/very/long/path/that/pushes/the/line/over/the/limit
// for details.
func Foo() {}
`),
		},
		{
			Name: "IndentedCommentIsWrappedWithIndent",
			Loc:  exam.Here(),
			Src: `
package p

func Foo() {
	// The quick brown fox jumps over the lazy dog and then some more words here now.
}
`,
			Limit:      60,
			WantChange: true,
			Golden: exam.GoldenHere(`
package p

func Foo() {
	// The quick brown fox jumps over the lazy dog and then
	// some more words here now.
}
`),
		},
		{
			Name: "MultipleCommentGroupsAreReflowed",
			Loc:  exam.Here(),
			Src: `
package p

// This is the first long comment that definitely exceeds the sixty character limit set here.

// This is the second long comment that also definitely exceeds the sixty character limit too.
func Foo() {}
`,
			Limit:      60,
			WantChange: true,
			Golden: exam.GoldenHere(`
package p

// This is the first long comment that definitely exceeds
// the sixty character limit set here.

// This is the second long comment that also definitely
// exceeds the sixty character limit too.
func Foo() {}
`),
		},
		{
			Name: "CannotJoinLinesCloseToLimit",
			Loc:  exam.Here(),
			Src: `
package p

// The quick brown fox jumps over the lazy dog alpha beta.
// The quick brown fox jumps over the lazy dog gamma delta.
func Foo() {}
`,
			Limit:      60,
			WantChange: false,
		},
		{
			Name: "PlainParagraphConsolidation",
			Loc:  exam.Here(),
			Src: `
package p

// Short line one.
// Short line two.
func Foo() {}
`,
			Limit:      60,
			WantChange: true,
			Golden: exam.GoldenHere(`
package p

// Short line one. Short line two.
func Foo() {}
`),
		},
		{
			Name: "ParagraphReflowWhenBothLineAreSlightlyOverLimit",
			Loc:  exam.Here(),
			Src: `
package p

// The quick brown fox jumps over the lazy dog one two three x
// four five six seven eight nine ten eleven twelve thirteen.
func Foo() {}
`,
			Limit:      60,
			WantChange: true,
			Golden: exam.GoldenHere(`
package p

// The quick brown fox jumps over the lazy dog one two three
// x four five six seven eight nine ten eleven twelve
// thirteen.
func Foo() {}
`),
		},
		{
			Name: "ParagraphBreakAtBlankComment",
			Loc:  exam.Here(),
			Src: `
package p

// The quick brown fox jumps over the lazy dog alpha beta gamma delta epsilon.
//
// The quick brown fox jumps over the lazy dog zeta eta theta iota kappa.
func Foo() {}
`,
			Limit:      60,
			WantChange: true,
			Golden: exam.GoldenHere(`
package p

// The quick brown fox jumps over the lazy dog alpha beta
// gamma delta epsilon.
//
// The quick brown fox jumps over the lazy dog zeta eta
// theta iota kappa.
func Foo() {}
`),
		},
		{
			Name: "BulletItemWrapping",
			Loc:  exam.Here(),
			Src: `
package p

// Items:
//
//   - The quick brown fox jumps over the lazy dog and then keeps running far away.
func Foo() {}
`,
			Limit:      60,
			WantChange: true,
			Golden: exam.GoldenHere(`
package p

// Items:
//
//   - The quick brown fox jumps over the lazy dog and then
//     keeps running far away.
func Foo() {}
`),
		},
		{
			Name: "NumberedItemWrapping",
			Loc:  exam.Here(),
			Src: `
package p

// Steps:
//
//  1. The quick brown fox jumps over the lazy dog and then keeps running away.
func Foo() {}
`,
			Limit:      60,
			WantChange: true,
			Golden: exam.GoldenHere(`
package p

// Steps:
//
//  1. The quick brown fox jumps over the lazy dog and then
//     keeps running away.
func Foo() {}
`),
		},
		{
			Name: "ListItemsNotMerged",
			Loc:  exam.Here(),
			Src: `
package p

// Items:
//
//   - Alpha item: the quick brown fox jumps over the lazy dog runs far away.
//   - Beta item: the quick brown fox jumps over the lazy dog runs far away too.
func Foo() {}
`,
			Limit:      60,
			WantChange: true,
			Golden: exam.GoldenHere(`
package p

// Items:
//
//   - Alpha item: the quick brown fox jumps over the lazy
//     dog runs far away.
//   - Beta item: the quick brown fox jumps over the lazy
//     dog runs far away too.
func Foo() {}
`),
		},
	}
	for _, tc := range cases {
		exam.Run(t, tc.Name, tc.Golden.GetLoc(), func(t *testing.T) {
			exam.Must(t, exam.True(strings.HasPrefix(tc.Src, "\n")),
				"test source should start with a newline to make them readable.")
			src := tc.Src[1:] // strip leading newline for processing
			out, err := wrap.File([]byte(src), tc.Limit)
			exam.Must(t, exam.Nil(err))
			if tc.WantChange {
				exam.Must(t, exam.GoldenEqual(string(out), tc.Golden))
			} else {
				exam.Must(t, exam.Equal(src, string(out)))
			}
		})
	}
}

// helper calls wrap.File and returns the string output, failing on error.
func mustWrap(t *testing.T, src string, limit int) string {
	t.Helper()
	out, err := wrap.File([]byte(src), limit)
	if err != nil {
		t.Fatalf("wrap.File error: %v", err)
	}
	return string(out)
}

// TestListItemWithContinuation verifies that a bullet item split across a
// marker line and a continuation line is collapsed when the joined body fits
// on a single line.
func TestListItemWithContinuation(t *testing.T) {
	src := `package p

// Items:
//
//   - Short item.
//     And more text.
func Foo() {}
`
	want := `package p

// Items:
//
//   - Short item. And more text.
func Foo() {}
`
	out := mustWrap(t, src, 60)
	if out != want {
		t.Errorf("expected continuation to be joined\ngot:\n%s\nwant:\n%s", out, want)
	}
}

// TestListItemAlreadyOptimal verifies that a multi-line list item which is
// already optimally wrapped (body cannot fit in fewer lines) is left unchanged.
func TestListItemAlreadyOptimal(t *testing.T) {
	src := `package p

// Items:
//
//   - The quick brown fox jumps over the lazy dog alpha.
//     The quick brown fox jumps over the lazy dog beta.
func Foo() {}
`
	out := mustWrap(t, src, 60)
	if out != src {
		t.Errorf("expected no change for already-optimal list item\ngot:\n%s", out)
	}
}
