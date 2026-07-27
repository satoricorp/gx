package semantic

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepoIndexStateRoundTrips(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	state := &RepoIndexState{
		Namespace:      "gx-org-acme-widgets-v2",
		RepoFullName:   "acme/widgets",
		EmbeddingModel: "text-embedding-3-large",
		Dimensions:     3072,
		ChunkerVersion: CodeChunkerVersion,
		SchemaVersion:  IndexSchemaVersion,
		Files: map[string]FileIndexState{
			"a.go": {Hash: "hash-a", Size: 12, ChunkHashes: []string{"c1", "c2"}},
			"b.go": {Hash: "hash-b", Size: 4, ChunkHashes: []string{"c3"}},
		},
	}
	if err := SaveRepoIndexState(path, state); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("temp file left behind after an atomic write")
	}

	loaded := LoadRepoIndexState(path)
	if loaded == nil {
		t.Fatal("load returned nil")
	}
	if loaded.UpdatedAt == 0 {
		t.Fatal("UpdatedAt was not stamped")
	}
	if got := loaded.Files["a.go"]; got.Hash != "hash-a" || len(got.ChunkHashes) != 2 {
		t.Fatalf("file state = %#v", got)
	}
	if files := loaded.SortedFiles(); len(files) != 2 || files[0] != "a.go" || files[1] != "b.go" {
		t.Fatalf("SortedFiles() = %v", files)
	}
}

func TestRepoIndexStateCompatibilityGatesReuse(t *testing.T) {
	state := &RepoIndexState{
		Version:        repoIndexStateVersion,
		Namespace:      "ns",
		EmbeddingModel: "text-embedding-3-large",
		Dimensions:     3072,
		ChunkerVersion: 7,
		SchemaVersion:  2,
	}
	if !state.Compatible("ns", "text-embedding-3-large", 3072, 7, 2) {
		t.Fatal("identical parameters must be compatible")
	}
	// Every parameter that can change chunk text or vector shape must
	// invalidate the manifest; reusing it would leave a namespace holding two
	// incompatible generations of rows.
	cases := []struct {
		name       string
		ns, model  string
		dims       int
		chunker    int
		schemaVers int
	}{
		{"namespace", "other", "text-embedding-3-large", 3072, 7, 2},
		{"model", "ns", "text-embedding-3-small", 3072, 7, 2},
		{"dimensions", "ns", "text-embedding-3-large", 1536, 7, 2},
		{"chunker", "ns", "text-embedding-3-large", 3072, 8, 2},
		{"schema", "ns", "text-embedding-3-large", 3072, 7, 3},
	}
	for _, tc := range cases {
		if state.Compatible(tc.ns, tc.model, tc.dims, tc.chunker, tc.schemaVers) {
			t.Fatalf("a changed %s must invalidate the manifest", tc.name)
		}
	}
	var missing *RepoIndexState
	if missing.Compatible("ns", "text-embedding-3-large", 3072, 7, 2) {
		t.Fatal("a nil manifest must not be reusable")
	}
}

func TestLoadRepoIndexStateToleratesGarbage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	if LoadRepoIndexState(path) != nil {
		t.Fatal("a corrupt manifest must be treated as absent, not trusted")
	}
	if LoadRepoIndexState(filepath.Join(t.TempDir(), "missing.json")) != nil {
		t.Fatal("a missing manifest must load as nil")
	}
}

func TestRepoIndexStatePathUsesGXHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	path, err := RepoIndexStatePath("gx-org-acme-widgets-v2")
	if err != nil {
		t.Fatalf("path: %v", err)
	}
	want := filepath.Join(home, "index", "gx-org-acme-widgets-v2.json")
	if path != want {
		t.Fatalf("path = %q, want %q", path, want)
	}
	if _, err := RepoIndexStatePath("  "); err == nil {
		t.Fatal("an empty namespace must be rejected")
	}
}
