package codereview

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/satoricorp/totality/internal/cloud"
)

// Cloud-mode retrieval: every onboarded user's path. These tests pin that a
// signed-in reviewer with no provider keys still gets indexed code, sessions,
// and shared knowledge — and that when nothing is configured the sources say
// so instead of vanishing.

type fakeCloudSearcher struct {
	requests []cloud.ReviewSearchRequest
	result   cloud.ReviewSearchResult
	err      error
}

func (f *fakeCloudSearcher) SearchReviewIndex(_ context.Context, req cloud.ReviewSearchRequest) (cloud.ReviewSearchResult, error) {
	f.requests = append(f.requests, req)
	if f.err != nil {
		return cloud.ReviewSearchResult{}, f.err
	}
	return f.result, nil
}

// remoteRepoRoot is a real checkout with a GitHub origin, because the cloud
// path resolves the repository the same way production does — from git.
func remoteRepoRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, args := range [][]string{
		{"init", "--quiet"},
		{"remote", "add", "origin", "https://github.com/acme/app.git"},
	} {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	return root
}

func cloudRetrieveInput(t *testing.T) RetrieveInput {
	t.Helper()
	return RetrieveInput{
		RepoRoot:     remoteRepoRoot(t),
		Evidence:     &EvidenceLog{},
		ChangedFiles: []string{"internal/capture/event.go"},
		DiffSnippets: []DiffSnippet{{
			File: "internal/capture/event.go",
			Diff: "@@ -1 +1 @@\n+func RecordCaptureEvent() {}\n",
		}},
	}
}

func evidenceFor(t *testing.T, log *EvidenceLog, source string) EvidenceStatus {
	t.Helper()
	for _, status := range log.Statuses() {
		if status.Source == source {
			return status
		}
	}
	t.Fatalf("no evidence recorded for %q: %#v", source, log.Statuses())
	return EvidenceStatus{}
}

func TestCodeIndexRetrieverUsesCloudWhenSignedIn(t *testing.T) {
	searcher := &fakeCloudSearcher{result: cloud.ReviewSearchResult{
		Available: true,
		Namespace: "totality-org-uuid-acme-app-v2",
		Exists:    true,
		Rows: []cloud.ReviewSearchRow{{
			ID:   "row-1",
			Text: "func RecordCaptureEvent() {}",
			Attributes: map[string]any{
				"file_path":   "internal/capture/event.go",
				"start_line":  float64(10),
				"end_line":    float64(24),
				"symbol_name": "RecordCaptureEvent",
				"source_kind": "code_file",
			},
		}},
	}}
	in := cloudRetrieveInput(t)
	retriever := CodeIndexRetriever{CloudSearcher: searcher}

	snippets, err := retriever.Retrieve(context.Background(), in)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 1 {
		t.Fatalf("snippets = %d, want 1", len(snippets))
	}
	if snippets[0].Kind != "indexed_code" || snippets[0].File != "internal/capture/event.go" {
		t.Fatalf("snippet = %#v, want indexed_code for the changed file", snippets[0])
	}
	if snippets[0].StartLine != 10 || snippets[0].EndLine != 24 {
		t.Fatalf("snippet lines = %d-%d, want 10-24", snippets[0].StartLine, snippets[0].EndLine)
	}

	if len(searcher.requests) != 1 {
		t.Fatalf("cloud requests = %d, want 1", len(searcher.requests))
	}
	req := searcher.requests[0]
	if req.Target != "repo" || req.RepoFullName != "acme/app" {
		t.Fatalf("request = %#v, want repo target for acme/app", req)
	}
	if len(req.SourceKinds) != 1 || req.SourceKinds[0] != "code_file" {
		t.Fatalf("source kinds = %v, want [code_file]", req.SourceKinds)
	}
	if !strings.Contains(req.SymbolQuery, "RecordCaptureEvent") {
		t.Fatalf("symbol query = %q, want the changed identifier", req.SymbolQuery)
	}

	status := evidenceFor(t, in.Evidence, codeIndexEvidenceSource)
	if status.State != EvidenceOK || status.Snippets != 1 {
		t.Fatalf("evidence = %#v, want OK with 1 snippet", status)
	}
	if status.Namespace != "totality-org-uuid-acme-app-v2" {
		t.Fatalf("evidence namespace = %q, want the server-resolved one", status.Namespace)
	}
	if !strings.Contains(status.Detail, "via tl cloud") {
		t.Fatalf("evidence detail = %q, want the cloud attribution", status.Detail)
	}
}

func TestCodeIndexCloudMissingNamespaceOffersConnectRemedy(t *testing.T) {
	searcher := &fakeCloudSearcher{result: cloud.ReviewSearchResult{
		Available: true,
		Namespace: "totality-org-uuid-acme-app-v2",
		Exists:    false,
	}}
	in := cloudRetrieveInput(t)
	retriever := CodeIndexRetriever{CloudSearcher: searcher}

	snippets, err := retriever.Retrieve(context.Background(), in)
	if err != nil || len(snippets) != 0 {
		t.Fatalf("Retrieve() = %d snippets, %v; want none", len(snippets), err)
	}
	status := evidenceFor(t, in.Evidence, codeIndexEvidenceSource)
	if status.State != EvidenceMissing {
		t.Fatalf("evidence state = %q, want missing", status.State)
	}
	if status.Remedy != connectRepositoryRemedy {
		t.Fatalf("remedy = %q, want the connect-repository remedy", status.Remedy)
	}
}

func TestCodeIndexCloudFailureIsUnavailableNotFatal(t *testing.T) {
	searcher := &fakeCloudSearcher{err: errors.New("status 503")}
	in := cloudRetrieveInput(t)
	retriever := CodeIndexRetriever{CloudSearcher: searcher}

	snippets, err := retriever.Retrieve(context.Background(), in)
	if err != nil {
		t.Fatalf("a cloud outage must not fail the review: %v", err)
	}
	if len(snippets) != 0 {
		t.Fatalf("snippets = %d, want none", len(snippets))
	}
	status := evidenceFor(t, in.Evidence, codeIndexEvidenceSource)
	if status.State != EvidenceUnavailable || !strings.Contains(status.Detail, "503") {
		t.Fatalf("evidence = %#v, want unavailable with the failure detail", status)
	}
}

func TestSessionRetrieverUsesCloudWhenSignedIn(t *testing.T) {
	searcher := &fakeCloudSearcher{result: cloud.ReviewSearchResult{
		Available: true,
		Namespace: "totality-org-uuid-acme-app-v2",
		Exists:    true,
		Rows: []cloud.ReviewSearchRow{{
			ID:   "sess-row",
			Text: "[assistant] chose exponential backoff over jitter",
			Attributes: map[string]any{
				"session_id":  "sess-1",
				"source_kind": "published_session_context",
			},
		}},
	}}
	in := cloudRetrieveInput(t)
	retriever := SessionContextRetriever{CloudSearcher: searcher}

	snippets, err := retriever.Retrieve(context.Background(), in)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 1 || snippets[0].Kind != "indexed_session" || snippets[0].SessionID != "sess-1" {
		t.Fatalf("snippets = %#v, want one indexed_session for sess-1", snippets)
	}
	req := searcher.requests[0]
	if strings.Join(req.SourceKinds, ",") != strings.Join(sessionSourceKinds, ",") {
		t.Fatalf("source kinds = %v, want the session kinds", req.SourceKinds)
	}
	status := evidenceFor(t, in.Evidence, sessionEvidenceSource)
	if status.State != EvidenceOK || status.Snippets != 1 {
		t.Fatalf("evidence = %#v, want OK with 1 snippet", status)
	}
}

func TestKnowledgeRetrieverUsesCloudWhenSignedIn(t *testing.T) {
	searcher := &fakeCloudSearcher{result: cloud.ReviewSearchResult{
		Available: true,
		Namespace: "totality-review-knowledge",
		Exists:    true,
		Rows: []cloud.ReviewSearchRow{{
			ID:   "kn-1",
			Text: "Verify webhook signatures before parsing.",
			Attributes: map[string]any{
				"source_id": "owasp-webhooks",
				"title":     "Webhook verification",
				"publisher": "owasp",
			},
		}},
	}}
	in := cloudRetrieveInput(t)
	retriever := ReviewResourceRetriever{CloudSearcher: searcher}

	snippets, err := retriever.Retrieve(context.Background(), in)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 1 || snippets[0].Kind != "review_resource" {
		t.Fatalf("snippets = %#v, want one review_resource", snippets)
	}
	if searcher.requests[0].Target != "knowledge" {
		t.Fatalf("target = %q, want knowledge", searcher.requests[0].Target)
	}
	status := evidenceFor(t, in.Evidence, reviewKnowledgeEvidenceSource)
	if status.State != EvidenceOK || status.Namespace != "totality-review-knowledge" {
		t.Fatalf("evidence = %#v, want OK against the shared corpus", status)
	}
}

// The regression this pins: with no keys and no login, the knowledge source
// used to be silently absent — no retriever, no evidence line, and no way to
// tell an uninformed review from an informed one.
func TestKnowledgeSourceReportsDisabledWhenNothingConfigured(t *testing.T) {
	t.Setenv("TOTALITY_CLOUD_URL", "off")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("TOTALITY_OPENAI_API_KEY", "")
	t.Setenv("TURBOPUFFER_API_KEY", "")

	retriever := reviewResourceRetrieverFromEnv()
	if retriever == nil {
		t.Fatalf("reviewResourceRetrieverFromEnv() = nil; the source must exist to report itself disabled")
	}
	in := cloudRetrieveInput(t)
	snippets, err := retriever.Retrieve(context.Background(), in)
	if err != nil || len(snippets) != 0 {
		t.Fatalf("Retrieve() = %d snippets, %v; want none", len(snippets), err)
	}
	status := evidenceFor(t, in.Evidence, reviewKnowledgeEvidenceSource)
	if status.State != EvidenceDisabled {
		t.Fatalf("evidence state = %q, want disabled", status.State)
	}
	if !strings.Contains(status.Remedy, "tl auth login") {
		t.Fatalf("remedy = %q, want the sign-in remedy", status.Remedy)
	}
}

func TestCodeIndexSourceReportsSignInRemedyWhenNothingConfigured(t *testing.T) {
	t.Setenv("TOTALITY_CLOUD_URL", "off")
	t.Setenv("TURBOPUFFER_API_KEY", "")

	in := cloudRetrieveInput(t)
	retriever := CodeIndexRetriever{}
	snippets, err := retriever.Retrieve(context.Background(), in)
	if err != nil || len(snippets) != 0 {
		t.Fatalf("Retrieve() = %d snippets, %v; want none", len(snippets), err)
	}
	status := evidenceFor(t, in.Evidence, codeIndexEvidenceSource)
	if status.State != EvidenceDisabled || !strings.Contains(status.Remedy, "tl auth login") {
		t.Fatalf("evidence = %#v, want disabled with the sign-in remedy", status)
	}
}
