package semantic

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

type recordingEmbedder struct {
	mu     sync.Mutex
	inputs []string
	dims   int
}

func (e *recordingEmbedder) Embed(_ context.Context, inputs []string) ([][]float32, error) {
	e.mu.Lock()
	e.inputs = append(e.inputs, inputs...)
	e.mu.Unlock()
	dims := e.dims
	if dims == 0 {
		dims = 4
	}
	out := make([][]float32, len(inputs))
	for i := range inputs {
		vector := make([]float32, dims)
		vector[0] = float32(len(inputs[i]))
		out[i] = vector
	}
	return out, nil
}

func (e *recordingEmbedder) count() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.inputs)
}

func (e *recordingEmbedder) reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.inputs = nil
}

type recordingStore struct {
	mu      sync.Mutex
	rows    map[string]VectorRow
	deleted []string
	upserts int
}

func newRecordingStore() *recordingStore {
	return &recordingStore{rows: map[string]VectorRow{}}
}

func (s *recordingStore) Upsert(_ context.Context, rows []VectorRow) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.upserts += len(rows)
	for _, row := range rows {
		s.rows[row.ID] = row
	}
	return nil
}

func (s *recordingStore) DeleteRows(_ context.Context, ids []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deleted = append(s.deleted, ids...)
	for _, id := range ids {
		delete(s.rows, id)
	}
	return nil
}

func (s *recordingStore) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deleted = nil
	s.upserts = 0
}

func (s *recordingStore) rowCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.rows)
}

func testIndexOptions(t *testing.T, root string, embedder Embedder, store CodeVectorStore) RepoIndexOptions {
	t.Helper()
	cfg := Config{
		OpenAIEmbeddingModel: "test-model",
		EmbeddingDimensions:  4,
	}
	return RepoIndexOptions{
		RepoRoot:     root,
		RepoFullName: "acme/widgets",
		OrgID:        "org-1",
		CommitID:     "commit-a",
		BranchName:   "main",
		Reason:       "test",
		Config:       &cfg,
		Embedder:     embedder,
		Store:        store,
		StatePath:    filepath.Join(t.TempDir(), "manifest.json"),
		Concurrency:  2,
	}
}

func TestIndexRepositoryIsIncrementalAcrossRuns(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "internal/alpha.go", "package alpha\n\nfunc Alpha() string {\n\treturn \"a\"\n}\n")
	writeTestFile(t, root, "internal/beta.go", "package beta\n\nfunc Beta() string {\n\treturn \"b\"\n}\n")

	embedder := &recordingEmbedder{}
	store := newRecordingStore()
	opts := testIndexOptions(t, root, embedder, store)

	first, err := IndexRepository(context.Background(), opts)
	if err != nil {
		t.Fatalf("first index: %v", err)
	}
	if first.ChunksUpserted == 0 || first.FilesIndexed != 2 {
		t.Fatalf("first run = %+v, want both files indexed", first)
	}
	if embedder.count() != first.ChunksUpserted {
		t.Fatalf("embedded %d chunks for %d upserts", embedder.count(), first.ChunksUpserted)
	}
	firstRows := store.rowCount()

	// Second run with no edits must embed and upload nothing.
	embedder.reset()
	store.reset()
	second, err := IndexRepository(context.Background(), opts)
	if err != nil {
		t.Fatalf("second index: %v", err)
	}
	if !second.UpToDate() {
		t.Fatalf("second run = %+v, want no work", second)
	}
	if embedder.count() != 0 {
		t.Fatalf("re-embedded %d chunks for an unchanged repository", embedder.count())
	}
	if second.FilesSkipped != 2 || second.ChunksTotal != first.ChunksTotal {
		t.Fatalf("second run = %+v, want 2 skipped files and %d chunks", second, first.ChunksTotal)
	}

	// Editing one file re-embeds only that file.
	embedder.reset()
	store.reset()
	writeTestFile(t, root, "internal/beta.go", "package beta\n\nfunc Beta() string {\n\treturn \"changed\"\n}\n")
	third, err := IndexRepository(context.Background(), opts)
	if err != nil {
		t.Fatalf("third index: %v", err)
	}
	if third.FilesIndexed != 1 || third.FilesSkipped != 1 {
		t.Fatalf("third run = %+v, want exactly one file re-indexed", third)
	}
	if embedder.count() != third.ChunksUpserted || third.ChunksUpserted == 0 {
		t.Fatalf("third run embedded %d chunks for %d upserts", embedder.count(), third.ChunksUpserted)
	}
	if store.rowCount() != firstRows {
		t.Fatalf("row count changed from %d to %d on an in-place edit", firstRows, store.rowCount())
	}

	// Deleting a file removes exactly its rows.
	embedder.reset()
	store.reset()
	if err := os.Remove(filepath.Join(root, "internal", "beta.go")); err != nil {
		t.Fatalf("remove: %v", err)
	}
	fourth, err := IndexRepository(context.Background(), opts)
	if err != nil {
		t.Fatalf("fourth index: %v", err)
	}
	if fourth.FilesRemoved != 1 || fourth.ChunksDeleted == 0 {
		t.Fatalf("fourth run = %+v, want the deleted file's rows removed", fourth)
	}
	for _, id := range store.deleted {
		if !strings.HasPrefix(id, "tlc-") {
			t.Fatalf("deleted id %q is not a code row id", id)
		}
	}
	if store.rowCount() >= firstRows {
		t.Fatalf("row count = %d after a delete, want fewer than %d", store.rowCount(), firstRows)
	}
}

func TestIndexRepositoryDeletesTailChunksWhenAFileShrinks(t *testing.T) {
	root := t.TempDir()
	var big strings.Builder
	big.WriteString("package big\n")
	for i := 0; i < 12; i++ {
		big.WriteString("\nfunc F")
		big.WriteString(string(rune('a' + i)))
		big.WriteString("() {\n")
		big.WriteString(strings.Repeat("\tstep()\n", 40))
		big.WriteString("}\n")
	}
	writeTestFile(t, root, "big.go", big.String())

	embedder := &recordingEmbedder{}
	store := newRecordingStore()
	opts := testIndexOptions(t, root, embedder, store)

	first, err := IndexRepository(context.Background(), opts)
	if err != nil {
		t.Fatalf("first index: %v", err)
	}
	if first.ChunksTotal < 5 {
		t.Fatalf("expected several chunks, got %+v", first)
	}

	store.reset()
	embedder.reset()
	writeTestFile(t, root, "big.go", "package big\n\nfunc Fa() {\n\tstep()\n}\n")
	second, err := IndexRepository(context.Background(), opts)
	if err != nil {
		t.Fatalf("second index: %v", err)
	}
	if second.ChunksDeleted != first.ChunksTotal-second.ChunksTotal {
		t.Fatalf("deleted %d rows, want %d", second.ChunksDeleted, first.ChunksTotal-second.ChunksTotal)
	}
	if store.rowCount() != second.ChunksTotal {
		t.Fatalf("store holds %d rows, want %d", store.rowCount(), second.ChunksTotal)
	}
}

func TestIndexRepositoryRebuildsWhenEmbeddingConfigChanges(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "a.go", "package a\n\nfunc A() {}\n")

	embedder := &recordingEmbedder{}
	store := newRecordingStore()
	opts := testIndexOptions(t, root, embedder, store)
	if _, err := IndexRepository(context.Background(), opts); err != nil {
		t.Fatalf("first index: %v", err)
	}

	embedder.reset()
	changed := *opts.Config
	changed.EmbeddingDimensions = 8
	opts.Config = &changed
	embedder.dims = 8
	result, err := IndexRepository(context.Background(), opts)
	if err != nil {
		t.Fatalf("second index: %v", err)
	}
	if !result.FullRebuild || embedder.count() == 0 {
		t.Fatalf("changing embedding dimensions must force a full rebuild, got %+v", result)
	}
}

func TestIndexRepositorySkipsNonSourceFiles(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "keep.go", "package keep\n\nfunc Keep() {}\n")
	writeTestFile(t, root, "node_modules/pkg/index.js", "export function ignored() {}\n")
	writeTestFile(t, root, "assets/logo.png", "not really a png")
	writeTestFile(t, root, "package-lock.json", "{}\n")
	writeTestFile(t, root, "bundle.min.js", "var a=1;\n")
	writeTestFile(t, root, "binary.txt", "abc\x00def")

	embedder := &recordingEmbedder{}
	store := newRecordingStore()
	result, err := IndexRepository(context.Background(), testIndexOptions(t, root, embedder, store))
	if err != nil {
		t.Fatalf("index: %v", err)
	}
	if result.FilesIndexed != 1 {
		t.Fatalf("indexed %d files, want only keep.go (%+v)", result.FilesIndexed, result)
	}
	for _, row := range store.rows {
		if path, _ := row.Attributes[codeFieldFilePath].(string); path != "keep.go" {
			t.Fatalf("indexed unexpected file %q", path)
		}
	}
}

func TestIndexRepositoryWritesRetrievalAttributes(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "internal/semantic/turbopuffer.go",
		"package semantic\n\n// DeleteStaleCodeDocuments removes rows.\nfunc DeleteStaleCodeDocuments() error {\n\treturn nil\n}\n")

	embedder := &recordingEmbedder{}
	store := newRecordingStore()
	if _, err := IndexRepository(context.Background(), testIndexOptions(t, root, embedder, store)); err != nil {
		t.Fatalf("index: %v", err)
	}
	if len(store.rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(store.rows))
	}
	for _, row := range store.rows {
		attrs := row.Attributes
		if attrs[transcriptFieldSourceKind] != SourceKindCodeFile {
			t.Fatalf("source_kind = %v", attrs[transcriptFieldSourceKind])
		}
		if attrs[codeFieldRepoFullName] != "acme/widgets" || attrs[codeFieldOrgID] != "org-1" {
			t.Fatalf("repo/org attributes = %#v", attrs)
		}
		if attrs[codeFieldFilePath] != "internal/semantic/turbopuffer.go" || attrs[codeFieldFile] != attrs[codeFieldFilePath] {
			t.Fatalf("path attributes = %#v", attrs)
		}
		if attrs[codeFieldSymbolName] != "DeleteStaleCodeDocuments" || attrs[codeFieldSymbolKind] != "func" {
			t.Fatalf("symbol attributes = %#v", attrs)
		}
		symbol, _ := attrs[codeFieldSymbol].(string)
		if !strings.Contains(symbol, "DeleteStaleCodeDocuments") || !strings.Contains(symbol, "stale") {
			t.Fatalf("symbol full-text value = %q, want the identifier and its word parts", symbol)
		}
		if attrs[codeFieldLanguage] != "go" || attrs[codeFieldPackage] != "semantic" {
			t.Fatalf("language attributes = %#v", attrs)
		}
		start, _ := attrs[codeFieldStartLine].(int)
		end, _ := attrs[codeFieldEndLine].(int)
		if start < 1 || end < start {
			t.Fatalf("line attributes = %v-%v", attrs[codeFieldStartLine], attrs[codeFieldEndLine])
		}
		if attrs[transcriptFieldCommitID] != "commit-a" || attrs[codeFieldHeadSha] != "commit-a" {
			t.Fatalf("commit attributes = %#v", attrs)
		}
	}
}

func TestNamespaceForRepoMatchesConsoleShape(t *testing.T) {
	got := NamespaceForRepo("2f273110-b6ce-4b3b-95d3-e7c0ca802e83", "satoricorp/totality", "/Users/joe/git/tl")
	want := "totality-2f273110-b6ce-4b3b-95d3-e7c0ca802e83-satoricorp-totality-v2"
	if got != want {
		t.Fatalf("NamespaceForRepo() = %q, want %q", got, want)
	}
	local := NamespaceForRepo("", "", "/Users/joe/git/yeet")
	if !strings.HasPrefix(local, "totality-local-yeet-") || !strings.HasSuffix(local, "-v2") {
		t.Fatalf("local namespace = %q, want a totality-local-yeet-<hash>-v2 name", local)
	}
	if NamespaceForRepo("", "", "/Users/joe/git/yeet") != local {
		t.Fatal("local namespace is not stable for the same checkout")
	}
	if NamespaceForRepo("", "", "/Users/other/yeet") == local {
		t.Fatal("two different checkouts named yeet share a namespace")
	}
}

func TestBatchChunksRespectsCountAndByteBudgets(t *testing.T) {
	chunks := make([]Chunk, 10)
	for i := range chunks {
		chunks[i] = Chunk{ID: "c", Text: strings.Repeat("x", 100)}
	}
	byCount := batchChunks(chunks, 3, 1_000_000)
	if len(byCount) != 4 {
		t.Fatalf("batches = %d, want 4", len(byCount))
	}
	byBytes := batchChunks(chunks, 100, 250)
	for _, batch := range byBytes {
		if len(batch) > 3 {
			t.Fatalf("batch of %d chunks exceeds the byte budget", len(batch))
		}
	}
	if len(batchChunks(nil, 3, 100)) != 0 {
		t.Fatal("empty input produced batches")
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
