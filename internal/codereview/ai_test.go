package codereview

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseAIReviewOutputIncludesFileLineOverviewAndSources(t *testing.T) {
	brief := ReviewBrief{
		Context: []ContextSnippet{{SourceLabel: "L1", Ref: "REVIEW.md", Kind: "repo_doc", Source: "local", Publisher: "local"}},
		SourceRefs: []SourceRef{{
			ID: "L1", Kind: "local", Title: "REVIEW.md", File: "REVIEW.md", Publisher: "local",
		}},
		SourceCatalog: []SourceBrief{{ID: "google-eng-practices", Publisher: "Google"}},
	}
	output, err := parseAIReviewOutput(`{
		"overview":"This change updates auth handling.",
		"recommendations":[{
			"title":"Check auth rollback",
			"summary":"Session rollback can fail silently.",
			"benefit":"Prevents stale auth state.",
			"recommendation":"Add a test for rollback failure.",
			"strength":"Strong",
			"evidence":["internal/auth/session.go changed"],
			"file":"internal/auth/session.go",
			"line":"42",
			"source_labels":["L1","google-eng-practices","missing"]
		}]
	}`, brief)
	if err != nil {
		t.Fatalf("parseAIReviewOutput() error = %v", err)
	}
	if output.Overview != "This change updates auth handling." {
		t.Fatalf("Overview = %q", output.Overview)
	}
	if len(output.Findings) != 1 {
		t.Fatalf("Findings = %#v", output.Findings)
	}
	finding := output.Findings[0]
	if finding.File != "internal/auth/session.go" || finding.Line != 42 {
		t.Fatalf("File/Line = %q/%d", finding.File, finding.Line)
	}
	if len(finding.ResolvedSources) != 2 {
		t.Fatalf("ResolvedSources = %#v, want local + catalog", finding.ResolvedSources)
	}
	if finding.ResolvedSources[0].Opaque || finding.ResolvedSources[0].File != "REVIEW.md" {
		t.Fatalf("local source = %#v", finding.ResolvedSources[0])
	}
	if !finding.ResolvedSources[1].Opaque || finding.ResolvedSources[1].Publisher != "Google" {
		t.Fatalf("catalog source = %#v", finding.ResolvedSources[1])
	}
}

func TestParseAIReviewOutputEmptyRecommendationsIsValid(t *testing.T) {
	output, err := parseAIReviewOutput(`{"recommendations":[]}`, ReviewBrief{})
	if err != nil {
		t.Fatalf("parseAIReviewOutput() error = %v", err)
	}
	if len(output.Findings) != 0 {
		t.Fatalf("Findings = %#v, want empty", output.Findings)
	}
}

func TestParsePRSummaryReviewIncludesNotableChangesWithStringLine(t *testing.T) {
	summary, err := ParsePRSummaryReview(`{
		"overview":"Adds session validation.",
		"notable_changes":[
			{"file":"internal/auth/session.go","line":"42","note":"Reject expired tokens before refresh."},
			{"file":"","line":1,"note":"skip empty file"},
			{"file":"internal/auth/session.go","line":99,"note":""}
		],
		"recommendations":[]
	}`, ReviewBrief{})
	if err != nil {
		t.Fatalf("ParsePRSummaryReview() error = %v", err)
	}
	if summary.Overview != "Adds session validation." {
		t.Fatalf("Overview = %q", summary.Overview)
	}
	if len(summary.NotableChanges) != 1 {
		t.Fatalf("NotableChanges = %#v, want one anchored entry", summary.NotableChanges)
	}
	change := summary.NotableChanges[0]
	if change.File != "internal/auth/session.go" || change.Line != 42 || change.Note == "" {
		t.Fatalf("NotableChange = %#v", change)
	}
}

func TestParsePRSummaryReviewMissingNotableChangesIsEmpty(t *testing.T) {
	summary, err := ParsePRSummaryReview(`{"overview":"Purpose only.","recommendations":[]}`, ReviewBrief{})
	if err != nil {
		t.Fatalf("ParsePRSummaryReview() error = %v", err)
	}
	if summary.NotableChanges != nil && len(summary.NotableChanges) != 0 {
		t.Fatalf("NotableChanges = %#v, want empty", summary.NotableChanges)
	}
}

func TestReviewWithOverviewDropsNotableChanges(t *testing.T) {
	reviewer := cannedAIReviewer{payload: `{
		"overview":"Hidden overview.",
		"notable_changes":[{"file":"main.go","line":3,"note":"Notable change."}],
		"recommendations":[{"title":"T","summary":"S","benefit":"B","recommendation":"R"}]
	}`}
	overview, findings, err := reviewer.ReviewWithOverview(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("ReviewWithOverview() error = %v", err)
	}
	if overview != "Hidden overview." || len(findings) != 1 {
		t.Fatalf("overview=%q findings=%#v", overview, findings)
	}
	summary, err := reviewer.ReviewForSummary(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("ReviewForSummary() error = %v", err)
	}
	if len(summary.NotableChanges) != 1 {
		t.Fatalf("ReviewForSummary() NotableChanges = %#v", summary.NotableChanges)
	}
}

func TestPatchFocusedReviewIgnoresNotableChanges(t *testing.T) {
	findings, err := cannedAIReviewer{payload: `{
		"notable_changes":[{"file":"main.go","line":3,"note":"Should not surface in gx review."}],
		"recommendations":[{"title":"T","summary":"S","benefit":"B","recommendation":"R"}]
	}`}.Review(context.Background(), ReviewBrief{ReviewProfile: "patch_focused"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("findings = %#v", findings)
	}
	text := RenderMarkdown(Report{Findings: findings})
	if strings.Contains(text, "Should not surface in gx review") {
		t.Fatalf("RenderMarkdown() leaked notable_changes:\n%s", text)
	}
}

func TestResponsesAIReviewerReviewForSummaryParsesNotableChanges(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"output_text":"{\"overview\":\"Overview.\",\"notable_changes\":[{\"file\":\"main.go\",\"line\":\"3\",\"note\":\"Changed entry point.\"}],\"recommendations\":[]}"}`))
	}))
	defer server.Close()

	reviewer := &responsesAIReviewer{url: server.URL, token: "token", model: "test", client: server.Client()}
	summary, err := reviewer.ReviewForSummary(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("ReviewForSummary() error = %v", err)
	}
	if summary.Overview != "Overview." || len(summary.NotableChanges) != 1 || summary.NotableChanges[0].Line != 3 {
		t.Fatalf("summary = %#v", summary)
	}
}

func TestBedrockAnthropicReviewerReviewForSummaryParsesNotableChanges(t *testing.T) {
	reviewer := &bedrockAnthropicReviewer{
		region: "us-east-1", model: "test-model", accessKey: "key", secretKey: "secret",
		client: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			body := strings.NewReader(`{"content":[{"type":"text","text":"{\"overview\":\"Overview text\",\"notable_changes\":[{\"file\":\"main.go\",\"line\":3,\"note\":\"Changed entry point.\"}],\"recommendations\":[]}"}]}`)
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(body), Header: make(http.Header)}, nil
		})},
	}
	summary, err := reviewer.ReviewForSummary(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("ReviewForSummary() error = %v", err)
	}
	if summary.Overview != "Overview text" || len(summary.NotableChanges) != 1 || summary.NotableChanges[0].Line != 3 {
		t.Fatalf("summary = %#v", summary)
	}
}

func TestResponsesAIReviewerParsesCannedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"output_text":"{\"recommendations\":[{\"title\":\"T\",\"summary\":\"S\",\"benefit\":\"B\",\"recommendation\":\"R\",\"file\":\"main.go\",\"line\":3,\"source_labels\":[\"L1\"]}]}"}`))
	}))
	defer server.Close()

	reviewer := &responsesAIReviewer{url: server.URL, token: "token", model: "test", client: server.Client()}
	brief := ReviewBrief{
		SourceRefs: []SourceRef{{ID: "L1", Kind: "local", File: "main.go", StartLine: 3, Publisher: "local"}},
	}
	findings, err := reviewer.Review(context.Background(), brief)
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(findings) != 1 || findings[0].File != "main.go" || findings[0].Line != 3 {
		t.Fatalf("findings = %#v", findings)
	}
}

func TestBedrockAnthropicReviewerParsesCannedResponse(t *testing.T) {
	reviewer := &bedrockAnthropicReviewer{
		region: "us-east-1", model: "test-model", accessKey: "key", secretKey: "secret",
		client: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			body := strings.NewReader(`{"content":[{"type":"text","text":"{\"overview\":\"Overview text\",\"recommendations\":[]}"}]}`)
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(body), Header: make(http.Header)}, nil
		})},
	}
	overview, findings, err := reviewer.ReviewWithOverview(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("ReviewWithOverview() error = %v", err)
	}
	if overview != "Overview text" || len(findings) != 0 {
		t.Fatalf("overview=%q findings=%#v", overview, findings)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestMultiAIReviewerEmptyParseableResponseSucceeds(t *testing.T) {
	reviewer := multiAIReviewer{reviewers: []namedAIReviewer{{
		name: "openai", label: "OpenAI", reviewer: cannedAIReviewer{payload: `{"recommendations":[]}`},
	}}}
	findings, err := reviewer.Review(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("findings = %#v, want zero", findings)
	}
}

func TestMultiAIReviewerErrorsWhenAllProvidersFail(t *testing.T) {
	reviewer := multiAIReviewer{reviewers: []namedAIReviewer{
		{name: "openai", label: "OpenAI", reviewer: failingAIReviewer{err: errTestParseFail}},
		{name: "anthropic", label: "Anthropic", reviewer: failingAIReviewer{err: errTestParseFail}},
	}}
	if _, err := reviewer.Review(context.Background(), ReviewBrief{}); err == nil {
		t.Fatal("Review() succeeded, want error")
	}
}

func TestFallbackAIReviewerDoesNotFallbackOnEmptyFindings(t *testing.T) {
	var fallbackCalls int
	reviewer := fallbackAIReviewer{
		primary:  cannedAIReviewer{payload: `{"recommendations":[]}`},
		fallback: callbackAIReviewer{fn: func(context.Context, ReviewBrief) ([]Finding, error) {
			fallbackCalls++
			return []Finding{{ID: "fallback", Title: "Fallback", Summary: "S", Benefit: "B", Recommendation: "R"}}, nil
		}},
	}
	findings, err := reviewer.Review(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if fallbackCalls != 0 {
		t.Fatalf("fallbackCalls = %d, want 0", fallbackCalls)
	}
	if len(findings) != 0 {
		t.Fatalf("findings = %#v, want empty", findings)
	}
}

func TestFallbackAIReviewerFallsBackOnPrimaryError(t *testing.T) {
	reviewer := fallbackAIReviewer{
		primary:  failingAIReviewer{err: errTestParseFail},
		fallback: cannedAIReviewer{payload: `{"recommendations":[{"title":"T","summary":"S","benefit":"B","recommendation":"R"}]}`},
	}
	findings, err := reviewer.Review(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("findings = %#v", findings)
	}
}

var errTestParseFail = errTest("parse failed")

type errTest string

func (e errTest) Error() string { return string(e) }

type cannedAIReviewer struct {
	payload string
}

func (c cannedAIReviewer) Review(ctx context.Context, brief ReviewBrief) ([]Finding, error) {
	output, err := parseAIReviewOutput(c.payload, brief)
	if err != nil {
		return nil, err
	}
	return output.Findings, nil
}

func (c cannedAIReviewer) ReviewWithOverview(ctx context.Context, brief ReviewBrief) (string, []Finding, error) {
	output, err := parseAIReviewOutput(c.payload, brief)
	if err != nil {
		return "", nil, err
	}
	return output.Overview, output.Findings, nil
}

func (c cannedAIReviewer) ReviewForSummary(ctx context.Context, brief ReviewBrief) (PRSummaryReview, error) {
	output, err := parseAIReviewOutput(c.payload, brief)
	if err != nil {
		return PRSummaryReview{}, err
	}
	return aiReviewOutputToPRSummaryReview(output), nil
}

type failingAIReviewer struct {
	err error
}

func (f failingAIReviewer) Review(context.Context, ReviewBrief) ([]Finding, error) {
	return nil, f.err
}

type callbackAIReviewer struct {
	fn func(context.Context, ReviewBrief) ([]Finding, error)
}

func (c callbackAIReviewer) Review(ctx context.Context, brief ReviewBrief) ([]Finding, error) {
	return c.fn(ctx, brief)
}
