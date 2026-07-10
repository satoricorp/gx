package structural

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestAnalyzeFindsGoDefinitionsAndDependencies(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "internal/core/helper.go", "package core\n\nfunc NewThing() string { return \"x\" }\n")
	writeFile(t, root, "internal/app/app.go", "package app\n\nfunc Run() string { return NewThing() }\n")

	facts := Analyze(root, []string{"internal/app/app.go", "internal/core/helper.go"})
	if len(facts.Files) != 2 {
		t.Fatalf("Analyze() files = %#v", facts.Files)
	}
	if !reflect.DeepEqual(facts.Edges, []DependencyEdge{{
		FromFile: "internal/app/app.go",
		ToFile:   "internal/core/helper.go",
		Symbol:   "NewThing",
	}}) {
		t.Fatalf("Analyze() edges = %#v", facts.Edges)
	}
}

func TestAnalyzeFindsGoImportPackageDependencies(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module example.com/repo\n\ngo 1.22\n")
	writeFile(t, root, "internal/core/helper.go", "package core\n\nfunc Existing() string { return \"x\" }\n")
	writeFile(t, root, "internal/app/app.go", strings.Join([]string{
		"package app",
		"",
		"import \"example.com/repo/internal/core\"",
		"",
		"func Run() string {",
		"  return core.Existing()",
		"}",
	}, "\n"))

	facts := Analyze(root, []string{"internal/app/app.go", "internal/core/helper.go"})
	var found bool
	for _, edge := range facts.Edges {
		if edge.FromFile == "internal/app/app.go" &&
			edge.ToFile == "internal/core/helper.go" &&
			edge.Symbol == "import example.com/repo/internal/core" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Analyze() edges = %#v, want package import dependency", facts.Edges)
	}
}

func TestAnalyzeIgnoresExternalGoImports(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module example.com/repo\n\ngo 1.22\n")
	writeFile(t, root, "internal/app/app.go", strings.Join([]string{
		"package app",
		"",
		"import \"github.com/stretchr/testify/require\"",
		"",
		"func Run() {",
		"  _ = require.New",
		"}",
	}, "\n"))

	facts := Analyze(root, []string{"internal/app/app.go"})
	if len(facts.Edges) != 0 {
		t.Fatalf("Analyze() edges = %#v, want no external import dependency", facts.Edges)
	}
}

func TestAnalyzeIgnoresLocalVariableFalseDependencies(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "internal/cli/root.go", strings.Join([]string{
		"package cli",
		"",
		"func Render() string {",
		"  args := []string{\"demux\"}",
		"  return args[0]",
		"}",
	}, "\n"))
	writeFile(t, root, "mcp/src/tool.ts", strings.Join([]string{
		"export const args = [\"demux\"];",
		"export const metadata = {};",
	}, "\n"))

	facts := Analyze(root, []string{"internal/cli/root.go", "mcp/src/tool.ts"})
	if len(facts.Edges) != 0 {
		t.Fatalf("Analyze() edges = %#v, want no local-variable or cross-language dependencies", facts.Edges)
	}
}

func TestEnclosingSymbolsMapsHunksToNearestSymbol(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "app.go", strings.Join([]string{
		"package app",
		"",
		"func First() string {",
		"  return \"one\"",
		"}",
		"",
		"func Second() string {",
		"  return \"two\"",
		"}",
		"",
	}, "\n"))
	facts := Analyze(root, []string{"app.go"})
	got := EnclosingSymbols(facts, []HunkInput{{ID: "h1", File: "app.go", NewStart: 8, NewLines: 1}})
	want := []HunkSymbol{{
		HunkID:    "h1",
		File:      "app.go",
		Symbol:    "Second",
		Kind:      "function",
		StartLine: 7,
		EndLine:   10,
	}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("EnclosingSymbols() = %#v, want %#v", got, want)
	}
}

func TestEnclosingSymbolsPrefersDefinitionStartingInsideHunkContext(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "app.go", strings.Join([]string{
		"package app",
		"",
		"func First() string {",
		"  return \"one\"",
		"}",
		"",
		"// gap 01",
		"// gap 02",
		"// gap 03",
		"// gap 04",
		"// gap 05",
		"// gap 06",
		"// gap 07",
		"// gap 08",
		"",
		"func Second() string {",
		"  return \"two\"",
		"}",
		"",
	}, "\n"))
	facts := Analyze(root, []string{"app.go"})
	got := EnclosingSymbols(facts, []HunkInput{{ID: "h1", File: "app.go", NewStart: 14, NewLines: 5}})
	if len(got) != 1 || got[0].Symbol != "Second" {
		t.Fatalf("EnclosingSymbols() = %#v, want Second", got)
	}
}

func TestAnalyzeSkipsUnsupportedFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "README.md", "# docs\n")
	facts := Analyze(root, []string{"README.md"})
	if len(facts.Files) != 0 || len(facts.Edges) != 0 {
		t.Fatalf("Analyze() = %#v", facts)
	}
}

func writeFile(t *testing.T, root, name, contents string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
