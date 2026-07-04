package codereview

import (
	"context"
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

	results := collectStaticToolResults(context.Background(), root, RepoFacts{DependencyFiles: []string{"go.mod"}}, Options{})
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

	results := collectStaticToolResults(context.Background(), root, RepoFacts{DependencyFiles: []string{"go.mod"}}, Options{})
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

	results := collectStaticToolResults(context.Background(), root, RepoFacts{DependencyFiles: []string{"go.mod"}}, Options{})
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

func installFakeGo(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "go")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("WriteFile(fake go) error = %v", err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}
