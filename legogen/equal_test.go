package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEqualTemplates(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name     string
		typeName string
		strategy string
		pkg      string
		wantFunc string
		wantCall string
	}{
		{
			name:     "comparer",
			typeName: "MyView",
			strategy: "comparer",
			pkg:      "mypkg",
			wantFunc: "func (m MyView) Equal(other MyView) bool {",
			wantCall: "return lego.EqualComparer(m, other)",
		},
		{
			name:     "viewer pointer",
			typeName: "*MyType",
			strategy: "viewer",
			pkg:      "mypkg",
			wantFunc: "func (m *MyType) Equal(other *MyType) bool {",
			wantCall: "return lego.EqualViewer(m, other)",
		},
		{
			name:     "slice",
			typeName: "SliceView",
			strategy: "slice",
			pkg:      "mypkg",
			wantFunc: "func (s SliceView) Equal(other SliceView) bool {",
			wantCall: "return lego.EqualSlice(s, other)",
		},
		{
			name:     "slicego",
			typeName: "IntList",
			strategy: "slicego",
			pkg:      "mypkg",
			wantFunc: "func (i IntList) Equal(other IntList) bool {",
			wantCall: "return lego.EqualSliceGo(i, other)",
		},
		{
			name:     "comparer in lego package",
			typeName: "Foo",
			strategy: "comparer",
			pkg:      "lego",
			wantFunc: "func (f Foo) Equal(other Foo) bool {",
			wantCall: "return EqualComparer(f, other)",
		},
		{
			name:     "viewer in lego package",
			typeName: "*Bar",
			strategy: "viewer",
			pkg:      "lego",
			wantFunc: "func (b *Bar) Equal(other *Bar) bool {",
			wantCall: "return EqualViewer(b, other)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			cmd := exec.Command(binary, "-type", tt.typeName, "-equal", tt.strategy)
			cmd.Dir = dir
			cmd.Env = append(os.Environ(), "GOPACKAGE="+tt.pkg)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("legogen failed: %v\n%s", err, out)
			}

			cleanType := strings.TrimPrefix(tt.typeName, "*")
			outputFile := filepath.Join(dir, strings.ToLower(cleanType)+"_legogen.go")
			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("output file not found: %v", err)
			}

			src := string(content)
			if !strings.Contains(src, tt.wantFunc) {
				t.Errorf("missing function signature %q\n\ngot:\n%s", tt.wantFunc, src)
			}
			if !strings.Contains(src, tt.wantCall) {
				t.Errorf("missing function call %q\n\ngot:\n%s", tt.wantCall, src)
			}
		})
	}
}
