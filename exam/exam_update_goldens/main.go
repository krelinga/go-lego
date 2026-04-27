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

	workspaceRoot, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: getting working directory: %v\n", err)
		os.Exit(1)
	}

	// Validate all paths before touching any file.
	for _, e := range entries {
		if err := validateEntryPath(e.Path, workspaceRoot); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
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

// applyDiffs reads path, delegates the patch logic to applyDiffsToSrc, and
// writes the result back.  entries must be sorted by Line ascending.
func applyDiffs(path string, entries []internal.GoldenEntry) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}
	result, err := applyDiffsToSrc(string(data), entries)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(result), 0644)
}
