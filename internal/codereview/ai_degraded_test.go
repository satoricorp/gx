package codereview

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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
	t.Setenv("GX_REVIEW_AI", "1")
	t.Setenv("GX_REVIEW_JUDGE", "0")
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "0")

	root := initRepo(t)
	writeFile(t, root, "go.mod", "module example.com/review\n")
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	gitAdd(t, root, "go.mod", "internal/app/app.go")
	gitCommit(t, root)
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\nfunc Stop() {}\n")

	engine := NewEngineWithReviewer(LocalScanner{}, StaticCatalog{}, defaultRules(), fakeRetriever{}, failingReviewer{err: errors.New("AI reviewers failed: OpenAI: connection refused")})
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

func TestReviewRecordsMissingAIConfiguration(t *testing.T) {
	t.Setenv("GX_REVIEW_AI", "1")
	t.Setenv("GX_REVIEW_JUDGE", "0")
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "0")
	t.Setenv("GX_OPENAI_PROXY_URL", "")
	t.Setenv("GX_CLOUD_URL", "off")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GX_OPENAI_API_KEY", "")
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
	if len(report.DegradedReasons) != 1 || report.DegradedReasons[0] != "AI reviewer not configured" {
		t.Fatalf("DegradedReasons = %#v, want missing configuration warning", report.DegradedReasons)
	}
}

func TestCompleteJSONRequestIncludesJSONKeyword(t *testing.T) {
	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		_, _ = w.Write([]byte(`{"output_text":"{\"recommendations\":[]}"}`))
	}))
	defer server.Close()

	reviewer := &responsesAIReviewer{url: server.URL, token: "token", model: "test", client: server.Client()}
	if _, err := reviewer.completeJSON(context.Background(), "Review the patch.", map[string]string{"changed_files": "a.go"}, 100); err != nil {
		t.Fatalf("completeJSON() error = %v", err)
	}
	instructions, _ := body["instructions"].(string)
	input, _ := body["input"].(string)
	if !strings.Contains(strings.ToLower(instructions), "json") {
		t.Fatalf("instructions = %q, want json keyword", instructions)
	}
	if !strings.Contains(strings.ToLower(input), "json") {
		t.Fatalf("input = %q, want json keyword", input)
	}
}
