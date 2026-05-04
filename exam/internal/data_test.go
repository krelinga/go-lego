package internal_test

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/krelinga/go-libs/exam"
	"github.com/krelinga/go-libs/exam/internal"
)

// writeTempFile writes b to a new temp file and returns the path.
func writeTempFile(t *testing.T, b []byte) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "golden_*.bin")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return f.Name()
}

func TestWriteAndReadGoldenEntry(t *testing.T) {
	cases := []struct {
		Name    string
		Loc     exam.Loc
		Entries []internal.GoldenEntry
	}{
		{
			Name: "single entry",
			Loc:  exam.Here(),
			Entries: []internal.GoldenEntry{
				{Path: "/workspace/foo_test.go", Line: 42, Text: "\nhello\n"},
			},
		},
		{
			Name: "multiple entries",
			Loc:  exam.Here(),
			Entries: []internal.GoldenEntry{
				{Path: "/workspace/a_test.go", Line: 10, Text: "\nfirst\n"},
				{Path: "/workspace/b_test.go", Line: 99, Text: "\nsecond\nwith\nmultiple\nlines\n"},
				{Path: "/workspace/a_test.go", Line: 55, Text: "\nthird\n"},
			},
		},
		{
			Name: "entry with empty text",
			Loc:  exam.Here(),
			Entries: []internal.GoldenEntry{
				{Path: "/workspace/empty_test.go", Line: 1, Text: ""},
			},
		},
		{
			Name: "entry with unicode content",
			Loc:  exam.Here(),
			Entries: []internal.GoldenEntry{
				{Path: "/workspace/unicode_test.go", Line: 7, Text: "\n日本語\n"},
			},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "goldens.bin")

			for _, e := range c.Entries {
				exam.Must(t, exam.Nil(internal.WriteGoldenEntry(path, e)))
			}

			got, err := internal.ReadGoldenEntries(path)
			exam.Must(t, exam.Nil(err))
			exam.Must(t, exam.Equal(len(got), len(c.Entries)))
			for i := range c.Entries {
				exam.Try(t, exam.Equal(got[i].Path, c.Entries[i].Path))
				exam.Try(t, exam.Equal(got[i].Line, c.Entries[i].Line))
				exam.Try(t, exam.Equal(got[i].Text, c.Entries[i].Text))
			}
		})
	}
}

func TestWriteGoldenEntry_MagicNumberWrittenOnce(t *testing.T) {
	path := filepath.Join(t.TempDir(), "goldens.bin")
	entry := internal.GoldenEntry{Path: "/workspace/foo_test.go", Line: 1, Text: "\ntext\n"}

	exam.Must(t, exam.Nil(internal.WriteGoldenEntry(path, entry)))
	exam.Must(t, exam.Nil(internal.WriteGoldenEntry(path, entry)))

	data, err := os.ReadFile(path)
	exam.Must(t, exam.Nil(err))

	// Magic number is exactly 8 bytes; it should appear only at the start.
	exam.Must(t, exam.True(len(data) >= 8))
	// The 9th byte onward must not begin with the magic bytes again (i.e. the
	// magic is not repeated for the second entry).
	magic := []byte{'G', 'L', 'D', 'N', 0, 1, 0, 0}
	exam.Try(t, exam.Equal(string(data[:8]), string(magic)))
	// Confirm there is content after the magic (two entries).
	exam.Try(t, exam.Greater(len(data), 8))
}

func TestReadGoldenEntries_WrongMagic(t *testing.T) {
	cases := []struct {
		Name  string
		Loc   exam.Loc
		Bytes []byte
	}{
		{
			Name:  "all zeros",
			Loc:   exam.Here(),
			Bytes: make([]byte, 32),
		},
		{
			Name:  "truncated magic",
			Loc:   exam.Here(),
			Bytes: []byte{'G', 'L', 'D'},
		},
		{
			Name:  "right prefix wrong version",
			Loc:   exam.Here(),
			Bytes: []byte{'G', 'L', 'D', 'N', 0, 2, 0, 0},
		},
		{
			Name:  "empty file",
			Loc:   exam.Here(),
			Bytes: []byte{},
		},
	}
	for _, c := range cases {
		exam.Run(t, c.Name, c.Loc, func(t *testing.T) {
			path := writeTempFile(t, c.Bytes)
			_, err := internal.ReadGoldenEntries(path)
			exam.Try(t, exam.NotNil(err))
		})
	}
}

func TestReadGoldenEntries_TruncatedEntry(t *testing.T) {
	// Write a valid file with one entry, then truncate mid-payload.
	path := filepath.Join(t.TempDir(), "goldens.bin")
	entry := internal.GoldenEntry{Path: "/workspace/foo_test.go", Line: 1, Text: "\ntext\n"}
	exam.Must(t, exam.Nil(internal.WriteGoldenEntry(path, entry)))

	data, err := os.ReadFile(path)
	exam.Must(t, exam.Nil(err))

	// Lop off the last few bytes to simulate a truncated write.
	truncated := data[:len(data)-4]
	path2 := writeTempFile(t, truncated)
	_, err = internal.ReadGoldenEntries(path2)
	exam.Try(t, exam.NotNil(err))
}

func TestReadGoldenEntries_CorruptEntrySize(t *testing.T) {
	// Write a valid magic header then a size field that claims a huge payload.
	magic := []byte{'G', 'L', 'D', 'N', 0, 1, 0, 0}
	sizeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBytes, 1<<32) // 4 GB — clearly not present
	content := append(magic, sizeBytes...)
	path := writeTempFile(t, content)

	_, err := internal.ReadGoldenEntries(path)
	exam.Try(t, exam.NotNil(err))
}

func TestWriteGoldenEntry_AppendsToExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "goldens.bin")
	e1 := internal.GoldenEntry{Path: "/workspace/a_test.go", Line: 1, Text: "\nfirst\n"}
	e2 := internal.GoldenEntry{Path: "/workspace/b_test.go", Line: 2, Text: "\nsecond\n"}

	exam.Must(t, exam.Nil(internal.WriteGoldenEntry(path, e1)))

	size1, err := os.Stat(path)
	exam.Must(t, exam.Nil(err))

	exam.Must(t, exam.Nil(internal.WriteGoldenEntry(path, e2)))

	size2, err := os.Stat(path)
	exam.Must(t, exam.Nil(err))

	// File must have grown.
	exam.Try(t, exam.Greater(size2.Size(), size1.Size()))

	// Both entries must round-trip correctly.
	got, err := internal.ReadGoldenEntries(path)
	exam.Must(t, exam.Nil(err))
	exam.Must(t, exam.Equal(len(got), 2))
	exam.Try(t, exam.Equal(got[0].Path, e1.Path))
	exam.Try(t, exam.Equal(got[1].Path, e2.Path))
}
