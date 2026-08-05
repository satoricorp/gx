package codereview

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type failingReviewer struct {
	err error
}

func (f failingReviewer) Review(context.Context, ReviewBrief) ([]Finding, error) {
	return nil, f.err
}

func TestRenderMarkdownShowsAIUnavailableWarning(t *testing.T) {
	text := RenderMarkdown(Report{
		DegradedReasons: []string{"AI review status 400: json_object requires json in input"},
	})
	if !strings.Contains(text, "> Warning: AI review unavailable (AI review status 400: json_object requires json in input); results are from deterministic checks only.") {
		t.Fatalf("RenderMarkdown() missing warning:\n%s", text)
	}
}

func TestReviewRecordsAIReviewerFailureWithEngine(t *testing.T) {
	t.Setenv("LGTM_REVIEW_AI", "1")
	t.Setenv("LGTM_REVIEW_JUDGE", "0")
	t.Setenv("LGTM_REVIEW_STATIC_TOOLS", "0")

	root := initRepo(t)
	writeFile(t, root, "go.mod", "module example.com/review\n")
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	gitAdd(t, root, "go.mod", "internal/app/app.go")
	gitCommit(t, root)
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\nfunc Stop() {}\n")

	engine := NewEngineWithReviewer(LocalScanner{}, StaticCatalog{}, defaultRules(), fakeRetriever{}, failingReviewer{err: errors.New("AI reviewers failed: Bedrock A: connection refused")})
	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.Reviewer != "heuristic fallback" {
		t.Fatalf("Reviewer = %q, want heuristic fallback", report.Reviewer)
	}
	if len(report.DegradedReasons) != 1 || !strings.Contains(report.DegradedReasons[0], "AI reviewers failed") {
		t.Fatalf("DegradedReasons = %#v, want AI failure recorded", report.DegradedReasons)
	}
	text := RenderMarkdown(report)
	if !strings.Contains(text, "> Warning: AI review unavailable") {
		t.Fatalf("RenderMarkdown() missing warning:\n%s", text)
	}
}

// TestReviewRecordsMissingAIConfiguration pins requirement 5c: with Bedrock as
// the only review provider, missing AWS credentials means nothing reviewed the
// change. The report used to say "AI reviewer not configured", which named
// neither the provider nor the fix and left a deterministic-only report looking
// like a completed review.
func TestReviewRecordsMissingAIConfiguration(t *testing.T) {
	t.Setenv("LGTM_REVIEW_AI", "1")
	t.Setenv("LGTM_REVIEW_JUDGE", "0")
	t.Setenv("LGTM_REVIEW_STATIC_TOOLS", "0")
	t.Setenv("LGTM_CLOUD_URL", "off")
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")

	root := initRepo(t)
	writeFile(t, root, "go.mod", "module example.com/review\n")
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	gitAdd(t, root, "go.mod", "internal/app/app.go")
	gitCommit(t, root)
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\nfunc Stop() {}\n")

	report, err := NewEngine().Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.DegradedReasons) != 1 {
		t.Fatalf("DegradedReasons = %#v, want one warning", report.DegradedReasons)
	}
	reason := report.DegradedReasons[0]
	for _, want := range []string{"AWS credentials", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"} {
		if !strings.Contains(reason, want) {
			t.Fatalf("DegradedReasons[0] = %q, want it to name %s", reason, want)
		}
	}
	if !strings.Contains(RenderMarkdown(report), "> Warning: AI review unavailable") {
		t.Fatalf("RenderMarkdown() missing warning:\n%s", RenderMarkdown(report))
	}
}
