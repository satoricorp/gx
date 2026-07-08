package codereview

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestParseReviewRiskPaths(t *testing.T) {
	text := strings.Join([]string{
		"# Review",
		"",
		"## high-risk paths",
		"",
		"risk-path: internal/auth/** — auth changes can leak or misuse credentials",
		"risk-path: internal/payments/* - payment flow changes need careful review",
		"",
		"## Other",
		"risk-path: ignored/outside-section — should not parse",
	}, "\n")
	paths := parseReviewRiskPaths(text)
	if len(paths) != 2 {
		t.Fatalf("paths = %#v, want 2 entries", paths)
	}
	if paths[0].Glob != "internal/auth/**" || !strings.Contains(paths[0].Message, "auth changes") {
		t.Fatalf("paths[0] = %#v", paths[0])
	}
	if paths[1].Glob != "internal/payments/*" {
		t.Fatalf("paths[1] = %#v", paths[1])
	}
}

func TestMatchRiskPathGlob(t *testing.T) {
	cases := []struct {
		pattern string
		file    string
		want    bool
	}{
		{"internal/auth/**", "internal/auth/session.go", true},
		{"internal/auth/**", "internal/auth/nested/token.go", true},
		{"internal/auth/**", "internal/storage/auth.go", false},
		{"internal/payments/*", "internal/payments/handler.go", true},
		{"**/Dockerfile", "deploy/Dockerfile", true},
	}
	for _, tc := range cases {
		if got := MatchRiskPathGlob(tc.pattern, tc.file); got != tc.want {
			t.Fatalf("MatchRiskPathGlob(%q, %q) = %v, want %v", tc.pattern, tc.file, got, tc.want)
		}
	}
}

func TestLoadReviewPolicyParsesRiskPaths(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "REVIEW.md", strings.Join([]string{
		"## high-risk paths",
		"risk-path: internal/auth/** — auth changes can leak or misuse credentials",
	}, "\n"))
	policy := LoadReviewPolicy(context.Background(), root)
	if len(policy.RiskPaths) != 1 {
		t.Fatalf("RiskPaths = %#v, want one entry", policy.RiskPaths)
	}
	if policy.RiskPaths[0].Glob != "internal/auth/**" {
		t.Fatalf("RiskPaths[0] = %#v", policy.RiskPaths[0])
	}
}

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

	brief, err := buildReviewBriefForTest(context.Background(), root, normalizeOptions(Options{}), RepoFacts{}, nil, fakeRetriever{})
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

	_, err := retriever.Retrieve(context.Background(), RetrieveInput{
		RepoRoot:     t.TempDir(),
		Options:      normalizeOptions(Options{ReviewPolicy: policy}),
		Facts:        RepoFacts{Files: []string{"internal/webhook/handler.go"}},
		ChangedFiles: []string{"internal/webhook/handler.go"},
		Plan:         ReviewExecutionPlan{RunReviewResources: true},
	})
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

	_, err := retriever.Retrieve(context.Background(), RetrieveInput{
		RepoRoot:     t.TempDir(),
		Options:      normalizeOptions(Options{ReviewPolicy: policy}),
		Facts:        RepoFacts{Files: []string{"internal/session/replay.go"}},
		ChangedFiles: []string{"internal/session/replay.go"},
	})
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(embedder.inputs) != 1 || !strings.Contains(embedder.inputs[0], "Prioritize session replay risks.") {
		t.Fatalf("embedder inputs = %#v, want policy query text", embedder.inputs)
	}
	if len(store.requests) == 0 {
		t.Fatalf("store requests = %#v, want indexed context query", store.requests)
	}
	filters := mustReviewResourceJSON(store.requests[0].Filters)
	if !strings.Contains(filters, "session_context") {
		t.Fatalf("indexed context filters = %s, want session_context", filters)
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

func TestMultiReviewerCallsEveryProviderBeforeLimitingFindings(t *testing.T) {
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
		t.Fatalf("findings = %d, want both reviewers' findings before engine cap", len(findings))
	}
}

func TestMultiReviewerRunsProvidersConcurrently(t *testing.T) {
	reviewer := multiAIReviewer{reviewers: []namedAIReviewer{
		{name: "slow-a", label: "Slow A", reviewer: delayedReviewer{delay: 120 * time.Millisecond, finding: Finding{ID: "one", Title: "one", Summary: "summary one"}}},
		{name: "slow-b", label: "Slow B", reviewer: delayedReviewer{delay: 120 * time.Millisecond, finding: Finding{ID: "two", Title: "two", Summary: "summary two"}}},
	}}

	start := time.Now()
	findings, err := reviewer.Review(context.Background(), ReviewBrief{})
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("findings = %#v, want two findings", findings)
	}
	if elapsed >= 220*time.Millisecond {
		t.Fatalf("Review() took %s, want roughly one child duration", elapsed)
	}
}

func TestMultiReviewerOutputOrderIsDeterministic(t *testing.T) {
	reviewer := multiAIReviewer{reviewers: []namedAIReviewer{
		{name: "first", label: "First", reviewer: delayedReviewer{delay: 80 * time.Millisecond, finding: Finding{ID: "one", Title: "one", Summary: "summary one"}}},
		{name: "second", label: "Second", reviewer: delayedReviewer{delay: 10 * time.Millisecond, finding: Finding{ID: "two", Title: "two", Summary: "summary two"}}},
	}}

	findings, err := reviewer.Review(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if got := []string{findings[0].ID, findings[1].ID}; strings.Join(got, ",") != "first.one,second.two" {
		t.Fatalf("finding order = %#v, want provider order", got)
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

type delayedReviewer struct {
	delay   time.Duration
	finding Finding
}

func (r delayedReviewer) Review(ctx context.Context, _ ReviewBrief) ([]Finding, error) {
	timer := time.NewTimer(r.delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
		return []Finding{r.finding}, nil
	}
}
