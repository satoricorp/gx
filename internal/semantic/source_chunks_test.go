package semantic

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/satoricorp/gx/internal/reviewbundle"
)

func TestBuildRepositoryChunksIndexesSourceFiles(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "src/app.go", `package main

func Run() {
	println("hello")
}
`)
	writeTestFile(t, root, "node_modules/pkg/index.js", `export function ignored() {}`)
	writeTestFile(t, root, "assets/logo.png", "not really a png")
	remote := "git@github.com:acme/widgets.git"
	branch := "main"
	bundle := reviewbundle.Bundle{
		Repo: reviewbundle.RepoPayload{
			RootPath:   root,
			RemoteURL:  &remote,
			BranchName: &branch,
		},
		Push: reviewbundle.PushPayload{HeadCommitID: "abc123"},
	}

	chunks := BuildRepositoryChunks(bundle)
	if len(chunks) != 2 {
		t.Fatalf("chunks = %#v, want package preamble and symbol chunks", chunks)
	}
	chunk := chunks[1]
	if chunk.Attributes["source_kind"] != "code_file" ||
		chunk.Attributes["repo_full_name"] != "acme/widgets" ||
		chunk.Attributes["file_path"] != "src/app.go" ||
		chunk.Attributes["symbol"] != "Run" ||
		chunk.Attributes["language"] != "go" {
		t.Fatalf("attributes = %#v", chunk.Attributes)
	}
	if chunk.ID == "" || chunk.Text == "" {
		t.Fatalf("chunk = %#v, want stable id and text", chunk)
	}
}

func writeTestFile(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}
