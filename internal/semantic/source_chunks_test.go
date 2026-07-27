package semantic

import (
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
	if len(chunks) != 1 {
		t.Fatalf("chunks = %d, want one chunk for the single indexable file", len(chunks))
	}
	chunk := chunks[0]
	if chunk.Attributes[transcriptFieldSourceKind] != SourceKindCodeFile ||
		chunk.Attributes[codeFieldRepoFullName] != "acme/widgets" ||
		chunk.Attributes[codeFieldFilePath] != "src/app.go" ||
		chunk.Attributes[codeFieldSymbolName] != "Run" ||
		chunk.Attributes[codeFieldLanguage] != "go" {
		t.Fatalf("attributes = %#v", chunk.Attributes)
	}
	// The publish path and `gx index` must mint the same id for the same
	// chunk, otherwise one namespace ends up holding two copies of every file.
	if chunk.ID != CodeRowID("acme/widgets", "src/app.go", 0) {
		t.Fatalf("chunk id = %q, want the shared content-addressed code row id", chunk.ID)
	}
	if chunk.Text == "" {
		t.Fatalf("chunk = %#v, want embedded text", chunk)
	}
}

func TestDocTypeFromPath(t *testing.T) {
	cases := map[string]string{
		"internal/semantic/codeindex.go":      "code_chunk",
		"internal/semantic/codeindex_test.go": "test_file",
		"server/test/review.test.ts":          "test_file",
		"docs/how-it-works.mdx":               "architecture_doc",
		"README.md":                           "architecture_doc",
		"tsconfig.json":                       "config_file",
	}
	for path, want := range cases {
		if got := docTypeFromPath(path); got != want {
			t.Fatalf("docTypeFromPath(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestRepoFullNameFromRemoteURL(t *testing.T) {
	cases := map[string]string{
		"git@github.com:satoricorp/gx.git":     "satoricorp/gx",
		"https://github.com/satoricorp/gx.git": "satoricorp/gx",
		"https://github.com/satoricorp/gx":     "satoricorp/gx",
		"https://example.com/not/github":       "",
		"":                                     "",
	}
	for input, want := range cases {
		if got := repoFullNameFromRemoteURL(input); got != want {
			t.Fatalf("repoFullNameFromRemoteURL(%q) = %q, want %q", input, got, want)
		}
	}
}
