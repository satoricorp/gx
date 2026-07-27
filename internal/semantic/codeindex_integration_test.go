package semantic

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// requireLiveIndexingCredentials skips unless real keys are present. Nothing in
// this file touches a production namespace: every namespace it creates is
// prefixed gx-eval- and deleted before the test returns.
func requireLiveIndexingCredentials(t *testing.T) Config {
	t.Helper()
	if strings.TrimSpace(os.Getenv("TURBOPUFFER_API_KEY")) == "" {
		t.Skip("TURBOPUFFER_API_KEY is not set; skipping live indexing test")
	}
	if strings.TrimSpace(os.Getenv("OPENAI_API_KEY")) == "" {
		t.Skip("OPENAI_API_KEY is not set; skipping live indexing test")
	}
	return CodeIndexConfigFromEnv()
}

func scratchNamespace(t *testing.T, label string) string {
	t.Helper()
	return fmt.Sprintf("gx-eval-%s-%d", label, time.Now().UnixNano())
}

// TestLiveIndexRoundTrip indexes a small fixture repository into a throwaway
// namespace and proves both retrieval legs work end to end: BM25 on the
// identifier field and ANN on the vector.
func TestLiveIndexRoundTrip(t *testing.T) {
	cfg := requireLiveIndexingCredentials(t)
	namespace := scratchNamespace(t, "roundtrip")

	root := t.TempDir()
	writeTestFile(t, root, "internal/semantic/turbopuffer.go", `package semantic

// DeleteStaleCodeDocuments removes code rows that no longer match the indexed
// commit, so a namespace never serves two generations of the same file.
func DeleteStaleCodeDocuments(repoFullName, commitID string) error {
	return nil
}
`)
	writeTestFile(t, root, "internal/billing/invoice.go", `package billing

// ApplyLateFee charges a penalty when an invoice is past due. This file exists
// only so the fixture has a second, unrelated topic to rank against.
func ApplyLateFee(amountCents int) int {
	return amountCents + 500
}
`)
	writeTestFile(t, root, "docs/overview.md", "# Overview\n\nThis repository publishes review artifacts.\n")

	client := NewTurboPufferClientForNamespace(cfg, namespace)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := client.DeleteNamespace(ctx); err != nil {
			t.Logf("cleanup: delete namespace %s: %v", namespace, err)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	statePath := t.TempDir() + "/manifest.json"
	opts := RepoIndexOptions{
		RepoRoot:     root,
		RepoFullName: "acme/fixture",
		OrgID:        "eval-org",
		Namespace:    namespace,
		CommitID:     "fixture-commit",
		BranchName:   "main",
		Reason:       "integration_test",
		StatePath:    statePath,
		Config:       &cfg,
	}
	result, err := IndexRepository(ctx, opts)
	if err != nil {
		t.Fatalf("IndexRepository: %v", err)
	}
	if result.ChunksUpserted == 0 {
		t.Fatalf("indexed nothing: %+v", result)
	}
	t.Logf("indexed %s", result.String())

	// TurboPuffer indexes are read-after-write consistent, but the FTS index
	// is built asynchronously; retry briefly rather than flaking.
	wantFile := "internal/semantic/turbopuffer.go"

	bm25 := retryQuery(t, ctx, client, QueryRequest{
		TextField:         codeFieldSymbol,
		TextQuery:         "DeleteStaleCodeDocuments",
		Limit:             5,
		IncludeAttributes: []string{codeFieldFilePath, codeFieldSymbolName, transcriptFieldSourceKind},
	})
	if got := rowFilePath(bm25[0]); got != wantFile {
		t.Fatalf("BM25 top hit = %q, want %q (rows=%v)", got, wantFile, rowFilePaths(bm25))
	}

	vectors, err := NewOpenAIEmbedder(cfg).Embed(ctx, []string{
		"where do we delete indexed code rows that belong to an old commit",
	})
	if err != nil || len(vectors) != 1 {
		t.Fatalf("embed query: %v", err)
	}
	ann := retryQuery(t, ctx, client, QueryRequest{
		Vector:            vectors[0],
		Limit:             5,
		IncludeAttributes: []string{codeFieldFilePath, codeFieldSymbolName, transcriptFieldSourceKind},
	})
	if got := rowFilePath(ann[0]); got != wantFile {
		t.Fatalf("vector top hit = %q, want %q (rows=%v)", got, wantFile, rowFilePaths(ann))
	}

	hybrid := retryQuery(t, ctx, client, QueryRequest{
		Vector:            vectors[0],
		TextField:         codeFieldSymbol,
		TextQuery:         "DeleteStaleCodeDocuments stale code documents",
		Limit:             5,
		IncludeAttributes: []string{codeFieldFilePath, codeFieldSymbolName, transcriptFieldSourceKind},
	})
	if got := rowFilePath(hybrid[0]); got != wantFile {
		t.Fatalf("hybrid top hit = %q, want %q (rows=%v)", got, wantFile, rowFilePaths(hybrid))
	}

	// A second run over an unchanged checkout must be a no-op.
	second, err := IndexRepository(ctx, opts)
	if err != nil {
		t.Fatalf("second IndexRepository: %v", err)
	}
	if !second.UpToDate() {
		t.Fatalf("second run was not a no-op: %+v", second)
	}
}

func retryQuery(t *testing.T, ctx context.Context, client *TurboPufferClient, req QueryRequest) []QueryRow {
	t.Helper()
	var lastErr error
	for attempt := 0; attempt < 10; attempt++ {
		rows, err := client.Query(ctx, req)
		if err == nil && len(rows) > 0 {
			return rows
		}
		lastErr = err
		time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond)
	}
	t.Fatalf("query returned no rows (last error: %v)", lastErr)
	return nil
}

func rowFilePath(row QueryRow) string {
	if value, ok := row[codeFieldFilePath].(string); ok {
		return value
	}
	if attrs, ok := row["attributes"].(map[string]any); ok {
		if value, ok := attrs[codeFieldFilePath].(string); ok {
			return value
		}
	}
	return ""
}

func rowFilePaths(rows []QueryRow) []string {
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, rowFilePath(row))
	}
	return out
}

// TestLiveIndexRealRepository indexes an actual checkout so the reported
// numbers (chunk count, wall clock, bytes) come from a real repository rather
// than a fixture. It is opt-in via GX_INDEX_EVAL_ROOT and writes to a scratch
// namespace it deletes afterwards unless GX_INDEX_EVAL_KEEP is set.
func TestLiveIndexRealRepository(t *testing.T) {
	root := strings.TrimSpace(os.Getenv("GX_INDEX_EVAL_ROOT"))
	if root == "" {
		t.Skip("GX_INDEX_EVAL_ROOT is not set; skipping real-repository indexing")
	}
	cfg := requireLiveIndexingCredentials(t)

	namespace := strings.TrimSpace(os.Getenv("GX_INDEX_EVAL_NAMESPACE"))
	if namespace == "" {
		namespace = scratchNamespace(t, "repo")
	}
	if !strings.HasPrefix(namespace, "gx-eval-") {
		t.Fatalf("refusing to write to %q: evaluation namespaces must be prefixed gx-eval-", namespace)
	}

	client := NewTurboPufferClientForNamespace(cfg, namespace)
	if strings.TrimSpace(os.Getenv("GX_INDEX_EVAL_KEEP")) == "" {
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			if err := client.DeleteNamespace(ctx); err != nil {
				t.Logf("cleanup: delete namespace %s: %v", namespace, err)
			}
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	statePath := t.TempDir() + "/manifest.json"
	if custom := strings.TrimSpace(os.Getenv("GX_INDEX_EVAL_STATE")); custom != "" {
		statePath = custom
	}

	result, err := IndexRepository(ctx, RepoIndexOptions{
		RepoRoot:    root,
		OrgID:       strings.TrimSpace(os.Getenv("GX_INDEX_EVAL_ORG")),
		Namespace:   namespace,
		Reason:      "eval",
		StatePath:   statePath,
		Config:      &cfg,
		Concurrency: 12,
		Logf:        t.Logf,
	})
	if err != nil {
		t.Fatalf("IndexRepository: %v", err)
	}
	t.Logf("REPO INDEX RESULT model=%s dims=%d %s", cfg.OpenAIEmbeddingModel, cfg.EmbeddingDimensions, result.String())
	t.Logf("REPO INDEX BYTES scanned=%d embedded=%d files_scanned=%d files_indexed=%d chunks=%d wall=%s",
		result.BytesScanned, result.BytesEmbedded, result.FilesScanned, result.FilesIndexed, result.ChunksTotal, result.Duration)
	if result.ChunksUpserted == 0 && !result.UpToDate() {
		t.Fatalf("indexed nothing: %+v", result)
	}
}
