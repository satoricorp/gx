package codereview

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoadReviewPolicyFetchesReferencesAndParsesReviewerModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><body><h1>Review reference</h1><p>Check authorization before side effects.</p></body></html>")
	}))
	defer server.Close()

	root := t.TempDir()
	writeFile(t, root, "REVIEW.md", strings.Join([]string{
		"# Review policy",
		"",
		"Use gpt-4o as the OpenAI reviewer model.",
		"Use anthropic:claude-sonnet-4-5 for the Anthropic reviewer.",
		"Always check tenant authorization.",
		"Reference: " + server.URL + "/guide",
	}, "\n"))

	policy := LoadReviewPolicy(context.Background(), root)
	if !policy.Present {
		t.Fatalf("policy = %#v, want present", policy)
	}
	if policy.OpenAIModelHint() != "gpt-4o" {
		t.Fatalf("OpenAIModelHint() = %q, want gpt-4o", policy.OpenAIModelHint())
	}
	if policy.AnthropicModelHint() != "anthropic.claude-sonnet-4-5" {
		t.Fatalf("AnthropicModelHint() = %q", policy.AnthropicModelHint())
	}
	if len(policy.References) != 1 || policy.References[0].Status != "ok" {
		t.Fatalf("References = %#v, want one fetched reference", policy.References)
	}
	if !strings.Contains(policy.References[0].Summary, "Check authorization before side effects.") {
		t.Fatalf("reference summary = %q", policy.References[0].Summary)
	}
}

func TestReviewPolicySummarizesOversizedMarkdown(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "REVIEW.md", strings.Repeat("Review this carefully. ", 1000))

	policy := LoadReviewPolicy(context.Background(), root)
	if !policy.Present || !policy.Summarized {
		t.Fatalf("policy = %#v, want summarized present policy", policy)
	}
	if !strings.Contains(policy.Text, "summary generated") {
		t.Fatalf("policy text = %q, want summary marker", policy.Text)
	}
}

func TestBuildReviewBriefUsesPolicyAndReferenceContextWithoutRenderingPolicyText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Prefer tests that exercise rollback behavior.")
	}))
	defer server.Close()

	root := t.TempDir()
	writeFile(t, root, "REVIEW.md", "Always check tenant authorization.\n\n"+server.URL+"/rollback\n")

	brief, err := BuildReviewBrief(context.Background(), root, normalizeOptions(Options{}), RepoFacts{}, nil, fakeRetriever{})
	if err != nil {
		t.Fatalf("BuildReviewBrief() error = %v", err)
	}
	if !hasContextSnippet(brief.Context, "review_policy", "REVIEW.md") {
		t.Fatalf("Context = %#v, want review_policy snippet", brief.Context)
	}
	if !hasContextSnippet(brief.Context, "review_reference", server.URL+"/rollback") {
		t.Fatalf("Context = %#v, want review_reference snippet", brief.Context)
	}

	report := Report{
		Verbose: true,
		Findings: []Finding{{
			ID:             "testing.rollback",
			Title:          "Add rollback coverage",
			Summary:        "Rollback behavior is changed without focused test coverage.",
			Benefit:        "Improves regression safety.",
			Recommendation: "Add a rollback test.",
		}},
		SourceRefs: brief.SourceRefs,
	}
	rendered := RenderMarkdown(report)
	if strings.Contains(rendered, "Always check tenant authorization") || strings.Contains(rendered, "Prefer tests that exercise rollback behavior") {
		t.Fatalf("RenderMarkdown() exposed policy content:\n%s", rendered)
	}
}

func TestReviewPolicyInfluencesReviewResourceQuery(t *testing.T) {
	policy := &ReviewPolicy{
		Present: true,
		Path:    "REVIEW.md",
		Text:    "Prioritize webhook signature verification.",
	}
	embedder := &recordingReviewPolicyEmbedder{vector: []float32{0.1, 0.2}}
	store := &recordingReviewResourceStore{}
	retriever := ReviewResourceRetriever{
		Embedder:  embedder,
		Store:     store,
		Namespace: "gx-review-knowledge",
		Limit:     2,
	}

	_, err := retriever.Retrieve(context.Background(), t.TempDir(), Options{ReviewPolicy: policy}, RepoFacts{Files: []string{"internal/webhook/handler.go"}}, nil)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(embedder.inputs) != 1 || !strings.Contains(embedder.inputs[0], "Prioritize webhook signature verification.") {
		t.Fatalf("embedder inputs = %#v, want policy query text", embedder.inputs)
	}
	if len(store.requests) == 0 || !strings.Contains(store.requests[0].Text, "Prioritize webhook signature verification.") {
		t.Fatalf("store requests = %#v, want policy query text", store.requests)
	}
}

func TestReviewPolicyInfluencesIndexedContextQuery(t *testing.T) {
	policy := &ReviewPolicy{
		Present: true,
		Path:    "REVIEW.md",
		Text:    "Prioritize session replay risks.",
	}
	embedder := &recordingReviewPolicyEmbedder{vector: []float32{0.3, 0.4}}
	store := &recordingIndexedContextStore{}
	retriever := IndexedContextRetriever{
		Embedder:  embedder,
		Store:     store,
		Namespace: "gx-sessions",
		Limit:     2,
	}

	_, err := retriever.Retrieve(context.Background(), t.TempDir(), Options{ReviewPolicy: policy}, RepoFacts{Files: []string{"internal/session/replay.go"}}, nil)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(embedder.inputs) != 1 || !strings.Contains(embedder.inputs[0], "Prioritize session replay risks.") {
		t.Fatalf("embedder inputs = %#v, want policy query text", embedder.inputs)
	}
	if len(store.requests) == 0 {
		t.Fatalf("store requests = %#v, want indexed context query", store.requests)
	}
}

func TestReviewerFromPolicyUsesOpenAIAndAnthropicModels(t *testing.T) {
	t.Setenv("GX_REVIEW_AI", "1")
	t.Setenv("GX_OPENAI_PROXY_URL", "")
	t.Setenv("GX_CLOUD_URL", "off")
	t.Setenv("OPENAI_API_KEY", "openai-key")
	t.Setenv("OPENAI_BASE_URL", "http://127.0.0.1:43123")
	t.Setenv("AWS_ACCESS_KEY_ID", "aws-key")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "aws-secret")
	t.Setenv("AWS_REGION", "us-west-2")

	reviewer := reviewerFromEnvWithPolicy(&ReviewPolicy{ModelHints: []ReviewModelHint{
		{Provider: "openai", Model: "gpt-4o"},
		{Provider: "anthropic", Model: "anthropic.claude-sonnet-4-5"},
	}})
	multi, ok := reviewer.(multiAIReviewer)
	if !ok {
		t.Fatalf("reviewer = %T, want multiAIReviewer", reviewer)
	}
	if len(multi.reviewers) != 2 {
		t.Fatalf("reviewers = %#v, want OpenAI and Anthropic", multi.reviewers)
	}
	openai, ok := multi.reviewers[0].reviewer.(*responsesAIReviewer)
	if !ok || openai.model != "gpt-4o" {
		t.Fatalf("OpenAI reviewer = %#v, want gpt-4o", multi.reviewers[0].reviewer)
	}
	anthropic, ok := multi.reviewers[1].reviewer.(*bedrockAnthropicReviewer)
	if !ok || anthropic.model != "anthropic.claude-sonnet-4-5" {
		t.Fatalf("Anthropic reviewer = %#v", multi.reviewers[1].reviewer)
	}
}

func TestMultiReviewerCallsEveryProviderWithoutTruncatingFindings(t *testing.T) {
	openai := &countingReviewer{findings: []Finding{
		{ID: "ai.review.1", Title: "one", Summary: "summary one"},
		{ID: "ai.review.2", Title: "two", Summary: "summary two"},
		{ID: "ai.review.3", Title: "three", Summary: "summary three"},
		{ID: "ai.review.4", Title: "four", Summary: "summary four"},
		{ID: "ai.review.5", Title: "five", Summary: "summary five"},
		{ID: "ai.review.6", Title: "six", Summary: "summary six"},
		{ID: "ai.review.7", Title: "seven", Summary: "summary seven"},
	}}
	anthropic := &countingReviewer{findings: []Finding{{
		ID: "ai.review.1", Title: "anthropic", Summary: "summary anthropic",
	}}}
	reviewer := multiAIReviewer{reviewers: []namedAIReviewer{
		{name: "openai", label: "OpenAI", reviewer: openai},
		{name: "anthropic", label: "Anthropic", reviewer: anthropic},
	}}

	findings, err := reviewer.Review(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if openai.calls != 1 || anthropic.calls != 1 {
		t.Fatalf("calls = openai %d anthropic %d, want both called once", openai.calls, anthropic.calls)
	}
	if len(findings) != 8 {
		t.Fatalf("findings = %d, want all provider findings", len(findings))
	}
}

type recordingReviewPolicyEmbedder struct {
	vector []float32
	inputs []string
}

func (e *recordingReviewPolicyEmbedder) Embed(_ context.Context, inputs []string) ([][]float32, error) {
	e.inputs = append(e.inputs, inputs...)
	out := make([][]float32, len(inputs))
	for i := range inputs {
		out[i] = e.vector
	}
	return out, nil
}

type recordingIndexedContextStore struct {
	requests []indexedContextQuery
}

func (s *recordingIndexedContextStore) Query(_ context.Context, req indexedContextQuery) ([]indexedContextRow, error) {
	s.requests = append(s.requests, req)
	return nil, nil
}

type countingReviewer struct {
	calls    int
	findings []Finding
}

func (r *countingReviewer) Review(context.Context, ReviewBrief) ([]Finding, error) {
	r.calls++
	return r.findings, nil
}
