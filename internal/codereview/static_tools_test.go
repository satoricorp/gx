package codereview

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStaticToolsUseAffectedPackagesForNarrowGoChanges(t *testing.T) {
	root := initRepo(t)
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	writeFile(t, root, "internal/other/other.go", "package other\nfunc Other() {}\n")
	gitAdd(t, root, "go.mod", "internal/app/app.go", "internal/other/other.go")
	gitCommit(t, root)
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\nfunc NewBehavior() {}\n")
	installFakeGo(t)
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "1")

	results := collectStaticToolResults(context.Background(), root, RepoFacts{DependencyFiles: []string{"go.mod"}}, Options{}, resolveChangeSet(context.Background(), root, "").Files)
	if len(results) != 2 {
		t.Fatalf("results = %#v, want go test and go vet", results)
	}
	if results[0].Command != "go test ./internal/app" || results[1].Command != "go vet ./internal/app" {
		t.Fatalf("commands = %q / %q, want affected package only", results[0].Command, results[1].Command)
	}
}

func TestStaticToolsFallbackForGoModChanges(t *testing.T) {
	root := initRepo(t)
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	gitAdd(t, root, "go.mod", "internal/app/app.go")
	gitCommit(t, root)
	writeFile(t, root, "go.mod", "module example.com/repo\n\ngo 1.24\n")
	installFakeGo(t)
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "1")

	results := collectStaticToolResults(context.Background(), root, RepoFacts{DependencyFiles: []string{"go.mod"}}, Options{}, resolveChangeSet(context.Background(), root, "").Files)
	if len(results) != 2 {
		t.Fatalf("results = %#v, want go test and go vet", results)
	}
	if results[0].Command != "go test ./..." || results[1].Command != "go vet ./..." {
		t.Fatalf("commands = %q / %q, want repo-wide fallback", results[0].Command, results[1].Command)
	}
}

func TestStaticToolsSkipDocsOnlyChanges(t *testing.T) {
	root := initRepo(t)
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "README.md", "# old\n")
	gitAdd(t, root, "go.mod", "README.md")
	gitCommit(t, root)
	writeFile(t, root, "README.md", "# new\n")
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "1")

	results := collectStaticToolResults(context.Background(), root, RepoFacts{DependencyFiles: []string{"go.mod"}}, Options{}, resolveChangeSet(context.Background(), root, "").Files)
	if len(results) != 0 {
		t.Fatalf("results = %#v, want no static tools for docs-only change", results)
	}
}

func TestGoPackageArgsFallbackAfterTwentyPackages(t *testing.T) {
	root := t.TempDir()
	var files []string
	for i := 0; i < 21; i++ {
		dir := filepath.Join(root, "pkg", "p"+string(rune('a'+i)))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}
		files = append(files, filepath.ToSlash(filepath.Join("pkg", "p"+string(rune('a'+i)), "file.go")))
	}

	args, ok := goPackageArgsForChangedFiles(root, files)
	if !ok || strings.Join(args, " ") != "./..." {
		t.Fatalf("goPackageArgsForChangedFiles() = %#v, %v; want ./... fallback", args, ok)
	}
}

func TestStaticToolRunnersDetectEcosystems(t *testing.T) {
	tests := []struct {
		name       string
		files      map[string]string
		deps       []string
		changed    []string
		pathTools  []string
		localTools []string
		want       []string
	}{
		{
			name:      "go stays scoped to the changed package",
			files:     map[string]string{"go.mod": "module example.com/repo\n", "internal/app/app.go": "package app\n"},
			deps:      []string{"go.mod"},
			changed:   []string{"internal/app/app.go"},
			pathTools: []string{"go"},
			want:      []string{"go test ./internal/app", "go vet ./internal/app"},
		},
		{
			name:       "typescript type-checks the whole project",
			files:      map[string]string{"tsconfig.json": "{}\n", "src/a.ts": "export const a = 1\n"},
			changed:    []string{"src/a.ts"},
			localTools: []string{"node_modules/.bin/tsc"},
			want:       []string{"tsc --noEmit"},
		},
		{
			name: "node test script never becomes a static tool",
			files: map[string]string{
				"package.json":  "{\n  \"scripts\": { \"test\": \"jest\" }\n}\n",
				"tsconfig.json": "{}\n",
				"src/a.ts":      "export const a = 1\n",
			},
			changed:    []string{"src/a.ts"},
			localTools: []string{"node_modules/.bin/tsc"},
			want:       []string{"tsc --noEmit"},
		},
		{
			name: "eslint scopes to changed javascript and typescript files",
			files: map[string]string{
				"eslint.config.js": "export default []\n",
				"src/a.ts":         "export const a = 1\n",
				"src/b.js":         "export const b = 1\n",
				"README.md":        "# repo\n",
			},
			changed:    []string{"src/a.ts", "src/b.js", "README.md"},
			localTools: []string{"node_modules/.bin/eslint"},
			want:       []string{"eslint ./src/a.ts ./src/b.js"},
		},
		{
			name:      "eslint is skipped without a config",
			files:     map[string]string{"src/a.js": "export const a = 1\n"},
			changed:   []string{"src/a.js"},
			pathTools: []string{"eslint"},
			want:      nil,
		},
		{
			name:       "deleted typescript file type-checks but is not linted",
			files:      map[string]string{"tsconfig.json": "{}\n", "eslint.config.js": "export default []\n"},
			changed:    []string{"src/gone.ts"},
			localTools: []string{"node_modules/.bin/tsc", "node_modules/.bin/eslint"},
			want:       []string{"tsc --noEmit"},
		},
		{
			name:      "ruff runs without configuration",
			files:     map[string]string{"app/main.py": "x = 1\n"},
			changed:   []string{"app/main.py"},
			pathTools: []string{"ruff", "mypy"},
			want:      []string{"ruff check ./app/main.py"},
		},
		{
			name:      "mypy runs once configured",
			files:     map[string]string{"app/main.py": "x = 1\n", "mypy.ini": "[mypy]\n"},
			changed:   []string{"app/main.py"},
			pathTools: []string{"ruff", "mypy"},
			want:      []string{"ruff check ./app/main.py", "mypy ./app/main.py"},
		},
		{
			name:      "mypy honours pyproject configuration",
			files:     map[string]string{"app/main.py": "x = 1\n", "pyproject.toml": "[tool.mypy]\nstrict = true\n"},
			changed:   []string{"app/main.py"},
			pathTools: []string{"mypy"},
			want:      []string{"mypy ./app/main.py"},
		},
		{
			name:      "cargo checks rust changes",
			files:     map[string]string{"Cargo.toml": "[package]\nname = \"repo\"\n", "src/lib.rs": "pub fn a() {}\n"},
			changed:   []string{"src/lib.rs"},
			pathTools: []string{"cargo"},
			want:      []string{"cargo check"},
		},
		{
			name: "missing binaries are skipped silently",
			files: map[string]string{
				"tsconfig.json":    "{}\n",
				"eslint.config.js": "export default []\n",
				"Cargo.toml":       "[package]\nname = \"repo\"\n",
				"mypy.ini":         "[mypy]\n",
				"src/a.ts":         "export const a = 1\n",
				"src/lib.rs":       "pub fn a() {}\n",
				"app/main.py":      "x = 1\n",
			},
			changed: []string{"src/a.ts", "src/lib.rs", "app/main.py"},
			want:    nil,
		},
		{
			name:      "documentation-only changes run nothing",
			files:     map[string]string{"go.mod": "module example.com/repo\n", "README.md": "# repo\n"},
			deps:      []string{"go.mod"},
			changed:   []string{"README.md"},
			pathTools: []string{"go", "tsc", "eslint", "ruff", "mypy", "cargo"},
			want:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			for rel, content := range test.files {
				writeFile(t, root, rel, content)
			}
			for _, rel := range test.localTools {
				writeExecutable(t, filepath.Join(root, filepath.FromSlash(rel)), "#!/bin/sh\nexit 0\n")
			}
			isolateFakeTools(t, test.pathTools...)

			env := staticToolEnv{
				repoRoot:     root,
				facts:        RepoFacts{DependencyFiles: test.deps},
				changedFiles: test.changed,
			}
			got := detectedStaticToolCommands(env)
			if strings.Join(got, " | ") != strings.Join(test.want, " | ") {
				t.Fatalf("detected commands = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestStaticToolFileArgsRejectsUnsafePaths(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "app/main.py", "x = 1\n")
	writeFile(t, root, "README.md", "# repo\n")
	writeFile(t, root, "outside.py", "x = 1\n")
	writeFile(t, root, "node_modules/pkg/main.py", "x = 1\n")
	if err := os.MkdirAll(filepath.Join(root, "pkg.py"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	tests := []struct {
		name  string
		files []string
		want  []string
	}{
		{name: "tracked file", files: []string{"app/main.py"}, want: []string{"./app/main.py"}},
		{name: "dot-slash prefix is normalized", files: []string{"./app/main.py"}, want: []string{"./app/main.py"}},
		{name: "parent escape", files: []string{"../outside.py"}, want: nil},
		{name: "embedded escape", files: []string{"app/../../outside.py"}, want: nil},
		{name: "absolute path", files: []string{filepath.Join(root, "app", "main.py")}, want: nil},
		{name: "missing file", files: []string{"app/deleted.py"}, want: nil},
		{name: "directory", files: []string{"pkg.py"}, want: nil},
		{name: "other extension", files: []string{"README.md"}, want: nil},
		{name: "vendored tree", files: []string{"node_modules/pkg/main.py"}, want: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := staticToolFileArgs(root, test.files, pythonFileExtensions...)
			if strings.Join(got, " ") != strings.Join(test.want, " ") {
				t.Fatalf("staticToolFileArgs(%#v) = %#v, want %#v", test.files, got, test.want)
			}
		})
	}
}

func TestStaticToolFileArgsCapsFileCount(t *testing.T) {
	root := t.TempDir()
	var files []string
	for i := 0; i < maxStaticToolFileArgs+5; i++ {
		rel := fmt.Sprintf("app/mod%02d.py", i)
		writeFile(t, root, rel, "x = 1\n")
		files = append(files, rel)
	}

	got := staticToolFileArgs(root, files, pythonFileExtensions...)
	if len(got) != maxStaticToolFileArgs {
		t.Fatalf("staticToolFileArgs() returned %d args, want cap of %d", len(got), maxStaticToolFileArgs)
	}
}

func TestLookStaticToolPrefersRepoLocalBinary(t *testing.T) {
	root := t.TempDir()
	local := filepath.Join(root, "node_modules", ".bin", "eslint")
	writeExecutable(t, local, "#!/bin/sh\nexit 0\n")
	isolateFakeTools(t, "eslint", "tsc")

	if got, ok := lookStaticTool(root, "eslint", nodeLocalBinDirs...); !ok || got != local {
		t.Fatalf("lookStaticTool(eslint) = %q, %v; want repo-local %q", got, ok, local)
	}
	if got, ok := lookStaticTool(root, "tsc", nodeLocalBinDirs...); !ok || filepath.Base(got) != "tsc" || strings.HasPrefix(got, root) {
		t.Fatalf("lookStaticTool(tsc) = %q, %v; want PATH fallback", got, ok)
	}
	if got, ok := lookStaticTool(root, "definitely-not-installed", nodeLocalBinDirs...); ok {
		t.Fatalf("lookStaticTool(missing) = %q, %v; want not found", got, ok)
	}
}

func TestStaticToolsRunDetectedNonGoTools(t *testing.T) {
	root := initRepo(t)
	writeFile(t, root, "tsconfig.json", "{}\n")
	writeFile(t, root, "eslint.config.js", "export default []\n")
	writeFile(t, root, "src/a.ts", "export const a = 1\n")
	gitAdd(t, root, "tsconfig.json", "eslint.config.js", "src/a.ts")
	gitCommit(t, root)
	writeFile(t, root, "src/a.ts", "export const a = 2\n")
	writeExecutable(t, filepath.Join(root, "node_modules", ".bin", "tsc"), "#!/bin/sh\nexit 0\n")
	writeExecutable(t, filepath.Join(root, "node_modules", ".bin", "eslint"), "#!/bin/sh\necho \"src/a.ts: error\"\nexit 1\n")
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "1")

	results := collectStaticToolResults(context.Background(), root, RepoFacts{}, Options{}, resolveChangeSet(context.Background(), root, "").Files)
	if len(results) != 2 {
		t.Fatalf("results = %#v, want tsc and eslint", results)
	}
	if results[0].Command != "tsc --noEmit" || results[0].ExitCode != 0 {
		t.Fatalf("tsc result = %#v, want clean tsc --noEmit", results[0])
	}
	if results[1].Command != "eslint ./src/a.ts" || results[1].ExitCode != 1 {
		t.Fatalf("eslint result = %#v, want failing eslint on the changed file", results[1])
	}
	if !strings.Contains(results[1].Output, "src/a.ts: error") {
		t.Fatalf("eslint output = %q, want captured tool output", results[1].Output)
	}
}

func detectedStaticToolCommands(env staticToolEnv) []string {
	var out []string
	for _, runner := range staticToolRunners {
		command, ok := runner.detect(env)
		if !ok {
			continue
		}
		out = append(out, strings.Join(command.argv, " "))
	}
	return out
}

func installFakeGo(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	writeExecutable(t, filepath.Join(dir, "go"), "#!/bin/sh\nexit 0\n")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// isolateFakeTools points PATH at a directory holding only the named no-op
// executables so detection never depends on what the machine has installed.
func isolateFakeTools(t *testing.T, names ...string) {
	t.Helper()
	dir := t.TempDir()
	for _, name := range names {
		writeExecutable(t, filepath.Join(dir, name), "#!/bin/sh\nexit 0\n")
	}
	t.Setenv("PATH", dir)
}

func writeExecutable(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", path, err)
	}
}
