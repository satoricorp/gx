package codereview

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

// TestParseAIReviewOutputToleratesFencedAndPrefixedJSON pins the reviewer's
// half of the JSON-shape guarantee that disappeared with the OpenAI reviewer.
//
// That path pinned the shape with response_format:json_object. Bedrock has no
// equivalent, so a leg may answer with a ```json fence or a sentence of preamble
// whenever it likes. This is not hypothetical: the first live cloud review after
// the Bedrock cutover failed on both legs at once with "invalid character '`'
// looking for beginning of value" and returned zero findings.
func TestParseAIReviewOutputToleratesFencedAndPrefixedJSON(t *testing.T) {
	body := `{"overview":"Adds a rate limiter.","recommendations":[{` +
		`"title":"Unsynchronized state",` +
		`"summary":"Allow mutates count without holding mu.",` +
		`"benefit":"Removes a data race under concurrent callers.",` +
		`"recommendation":"Lock mu for the whole read-modify-write.",` +
		`"file":"a.go","line":12}]}`
	for name, content := range map[string]string{
		"fenced with language": "```json\n" + body + "\n```",
		"fenced bare":          "```\n" + body + "\n```",
		"prose preamble":       "Here is the review:\n\n" + body,
		"trailing prose":       body + "\n\nLet me know if you want more detail.",
	} {
		t.Run(name, func(t *testing.T) {
			output, err := parseAIReviewOutput(content, ReviewBrief{})
			if err != nil {
				t.Fatalf("parseAIReviewOutput() error = %v", err)
			}
			if output.Overview != "Adds a rate limiter." {
				t.Fatalf("Overview = %q, want the decoded overview", output.Overview)
			}
			if len(output.Findings) != 1 {
				t.Fatalf("Findings = %d, want 1", len(output.Findings))
			}
		})
	}
}

// TestParseAIReviewOutputStillRejectsNonJSON keeps the fallback from turning a
// genuine failure into a silent empty review: a model that answers with prose
// and no object at all must still be an error, not zero findings.
func TestParseAIReviewOutputStillRejectsNonJSON(t *testing.T) {
	for name, content := range map[string]string{
		"prose only": "I could not review this change.",
		"empty":      "",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := parseAIReviewOutput(content, ReviewBrief{}); err == nil {
				t.Fatalf("parseAIReviewOutput(%q) error = nil, want an error", content)
			}
		})
	}
}

func TestParsePRSummaryReviewIncludesDownstreamImpact(t *testing.T) {
	summary, err := ParsePRSummaryReview(`{
		"overview":"Adjusts cache eviction.",
		"downstream_impact":"Low customer-facing risk; eviction timing may shift but API behavior stays the same.",
		"recommendations":[]
	}`, ReviewBrief{})
	if err != nil {
		t.Fatalf("ParsePRSummaryReview() error = %v", err)
	}
	if summary.DownstreamImpact != "Low customer-facing risk; eviction timing may shift but API behavior stays the same." {
		t.Fatalf("DownstreamImpact = %q", summary.DownstreamImpact)
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

func TestReviewReturnsFindingsWithoutSummaryFields(t *testing.T) {
	reviewer := cannedAIReviewer{payload: `{
		"overview":"Hidden overview.",
		"notable_changes":[{"file":"main.go","line":3,"note":"Notable change."}],
		"recommendations":[{"title":"T","summary":"S","benefit":"B","recommendation":"R"}]
	}`}
	findings, err := reviewer.Review(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("findings=%#v, want 1", findings)
	}
	summary, err := reviewer.ReviewForSummary(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("ReviewForSummary() error = %v", err)
	}
	if summary.Overview != "Hidden overview." || len(summary.NotableChanges) != 1 {
		t.Fatalf("ReviewForSummary() overview=%q notable=%#v", summary.Overview, summary.NotableChanges)
	}
}

func TestPatchFocusedReviewIgnoresNotableChanges(t *testing.T) {
	findings, err := cannedAIReviewer{payload: `{
		"notable_changes":[{"file":"main.go","line":3,"note":"Should not surface in tx review."}],
		"recommendations":[{"title":"T","summary":"S","benefit":"B","recommendation":"R"}]
	}`}.Review(context.Background(), ReviewBrief{ReviewProfile: "patch_focused"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("findings = %#v", findings)
	}
	text := RenderMarkdown(Report{Findings: findings})
	if strings.Contains(text, "Should not surface in tx review") {
		t.Fatalf("RenderMarkdown() leaked notable_changes:\n%s", text)
	}
}

func TestBedrockAnthropicReviewerReviewForSummaryParsesNotableChanges(t *testing.T) {
	reviewer := &bedrockAnthropicReviewer{
		model: "test-model",
		transport: &directBedrockTransport{
			region: "us-east-1", accessKey: "key", secretKey: "secret",
			client: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				body := strings.NewReader(`{"content":[{"type":"text","text":"{\"overview\":\"Overview text\",\"notable_changes\":[{\"file\":\"main.go\",\"line\":3,\"note\":\"Changed entry point.\"}],\"recommendations\":[]}"}]}`)
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(body), Header: make(http.Header)}, nil
			})},
		},
	}
	summary, err := reviewer.ReviewForSummary(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("ReviewForSummary() error = %v", err)
	}
	if summary.Overview != "Overview text" || len(summary.NotableChanges) != 1 || summary.NotableChanges[0].Line != 3 {
		t.Fatalf("summary = %#v", summary)
	}
}

func TestBedrockAnthropicReviewerParsesCannedResponse(t *testing.T) {
	reviewer := &bedrockAnthropicReviewer{
		model: "test-model",
		transport: &directBedrockTransport{
			region: "us-east-1", accessKey: "key", secretKey: "secret",
			client: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				body := strings.NewReader(`{"content":[{"type":"text","text":"{\"overview\":\"Overview text\",\"recommendations\":[]}"}]}`)
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(body), Header: make(http.Header)}, nil
			})},
		},
	}
	summary, err := reviewer.ReviewForSummary(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("ReviewForSummary() error = %v", err)
	}
	if summary.Overview != "Overview text" || len(summary.Findings) != 0 {
		t.Fatalf("overview=%q findings=%#v", summary.Overview, summary.Findings)
	}
}

// TestBedrockAnthropicReviewerResolvesFileLineAndSources ports the anchoring
// coverage that used to sit on the deleted OpenAI reviewer: a leg must still
// resolve file, line, and source labels from a canned model response.
func TestBedrockAnthropicReviewerResolvesFileLineAndSources(t *testing.T) {
	reviewer := &bedrockAnthropicReviewer{
		model: "us.anthropic.claude-test",
		transport: &directBedrockTransport{
			region: "us-west-2", accessKey: "key", secretKey: "secret",
			client: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				body := strings.NewReader(`{"content":[{"type":"text","text":"{\"recommendations\":[{\"title\":\"T\",\"summary\":\"S\",\"benefit\":\"B\",\"recommendation\":\"R\",\"file\":\"main.go\",\"line\":3,\"source_labels\":[\"L1\"]}]}"}]}`)
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(body), Header: make(http.Header)}, nil
			})},
		},
	}
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
	if len(findings[0].ResolvedSources) != 1 {
		t.Fatalf("ResolvedSources = %#v, want the local label resolved", findings[0].ResolvedSources)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestMultiAIReviewerEmptyParseableResponseSucceeds(t *testing.T) {
	reviewer := multiAIReviewer{reviewers: []namedAIReviewer{{
		name: "bedrock-a", label: "Bedrock A", reviewer: cannedAIReviewer{payload: `{"recommendations":[]}`},
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
		{name: "bedrock-a", label: "Bedrock A", reviewer: failingAIReviewer{err: errTestParseFail}},
		{name: "bedrock-b", label: "Bedrock B", reviewer: failingAIReviewer{err: errTestParseFail}},
	}}
	if _, err := reviewer.Review(context.Background(), ReviewBrief{}); err == nil {
		t.Fatal("Review() succeeded, want error")
	}
}

// TestMultiAIReviewerSurvivesOneFailedLeg is the property the two-leg panel
// exists for: one model failing must not lose the other model's findings.
func TestMultiAIReviewerSurvivesOneFailedLeg(t *testing.T) {
	reviewer := multiAIReviewer{reviewers: []namedAIReviewer{
		{name: "bedrock-a", label: "Bedrock A", reviewer: failingAIReviewer{err: errTestParseFail}},
		{name: "bedrock-b", label: "Bedrock B", reviewer: cannedAIReviewer{payload: `{"recommendations":[{"title":"T","summary":"S","benefit":"B","recommendation":"R"}]}`}},
	}}
	findings, err := reviewer.Review(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("findings = %#v, want the surviving leg's finding", findings)
	}
	if !strings.Contains(evidenceText(findings[0].Evidence), "Bedrock B") {
		t.Fatalf("Evidence = %#v, want the finding attributed to its leg", findings[0].Evidence)
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

func TestCompactBriefKeepsDiffSnippetsBeyondTheContextBudget(t *testing.T) {
	// Diff snippets are the primary evidence for what changed, so they must not
	// be throttled by the retrieved-context budget. A PR summary runs shallow
	// and deliberately builds a wide diff; capping it at maxAIContextSnippets
	// used to hide most of a large change from the model.
	brief := ReviewBrief{Depth: "shallow"}
	for i := 0; i < maxAIDiffSnippets+5; i++ {
		brief.Static.DiffSnippets = append(brief.Static.DiffSnippets, DiffSnippet{
			File: fmt.Sprintf("file%02d.go", i),
			Diff: "@@ -1 +1,2 @@\n+// change\n",
		})
	}
	for i := 0; i < maxAIContextSnippets+5; i++ {
		brief.Context = append(brief.Context, ContextSnippet{
			Ref:  fmt.Sprintf("doc%02d.md", i),
			Text: "context",
		})
	}

	compacted := compactReviewBriefForAI(brief)

	if got := len(compacted.Static.DiffSnippets); got != maxAIDiffSnippets {
		t.Fatalf("diff snippets = %d, want %d", got, maxAIDiffSnippets)
	}
	// The two budgets must stay independent, and the diff budget must never be
	// the smaller of the two: the diff is the primary evidence for what changed,
	// so retrieved context must not be able to squeeze it out.
	if maxAIDiffSnippets < maxAIContextSnippets {
		t.Fatalf("diff budget %d is below the retrieved-context budget %d; a wide diff would be hidden from the model", maxAIDiffSnippets, maxAIContextSnippets)
	}
	if got := len(compacted.Context); got != maxAIContextSnippets {
		t.Fatalf("context snippets = %d, want %d", got, maxAIContextSnippets)
	}
}
