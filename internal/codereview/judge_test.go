package codereview

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestReviewerAvailable(t *testing.T) {
	if reviewerAvailable(nil) {
		t.Fatal("reviewerAvailable(nil) = true, want false")
	}
	if reviewerAvailable(multiAIReviewer{reviewers: []namedAIReviewer{{reviewer: unavailableAIReviewer{reason: "missing"}}}}) {
		t.Fatal("reviewerAvailable(unavailable multi) = true, want false")
	}
	if !reviewerAvailable(&responsesAIReviewer{}) {
		t.Fatal("reviewerAvailable(real OpenAI reviewer) = false, want true")
	}
	if !reviewerAvailable(fakeReviewer{}) {
		t.Fatal("reviewerAvailable(fake without Available) = false, want true")
	}
}

func TestJudgeConfirmsFindingAndItSurvives(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "1")
	root := t.TempDir()
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	engine := judgeTestEngine(root, []Finding{judgeTestFinding("ai.review.1", "internal/app/app.go")})
	engine.judge = scriptedJudge{results: []judgeResult{{
		CandidateID:      "ai.review.1",
		Verdict:          "confirmed",
		Impact:           "breaking",
		Severity:         5,
		Confidence:       0.91,
		VerificationNote: "The named file contains the reviewed Module.",
	}}}

	report, err := engine.Review(context.Background(), root, Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Findings) != 1 || report.Findings[0].ID != "ai.review.1" {
		t.Fatalf("Findings = %#v, want confirmed finding", report.Findings)
	}
	if report.Findings[0].Strength != "Strong" {
		t.Fatalf("Strength = %q, want Strong", report.Findings[0].Strength)
	}
}

func TestJudgeWrongFindingDrops(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "1")
	root := t.TempDir()
	writeFile(t, root, "internal/app/app.go", "package app\n")
	engine := judgeTestEngine(root, []Finding{judgeTestFinding("ai.review.1", "internal/app/app.go")})
	engine.judge = scriptedJudge{results: []judgeResult{{CandidateID: "ai.review.1", Verdict: "wrong", Severity: 1, Confidence: 0.9, VerificationNote: "Not supported."}}}

	report, err := engine.Review(context.Background(), root, Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Findings) != 0 {
		t.Fatalf("Findings = %#v, want wrong finding dropped", report.Findings)
	}
}

func TestJudgeUnverifiedFindingDrops(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "1")
	root := t.TempDir()
	writeFile(t, root, "internal/app/app.go", "package app\n")
	engine := judgeTestEngine(root, []Finding{judgeTestFinding("ai.review.1", "internal/app/app.go")})
	engine.judge = scriptedJudge{results: []judgeResult{{CandidateID: "ai.review.1", Verdict: "unverified", Severity: 2, Confidence: 0.4, VerificationNote: "Needs more evidence."}}}

	report, err := engine.Review(context.Background(), root, Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Findings) != 0 {
		t.Fatalf("Findings = %#v, want unverified finding dropped", report.Findings)
	}
}

func TestJudgeErrorStillKeepsDedupedFindings(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "1")
	root := t.TempDir()
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	engine := judgeTestEngine(root, []Finding{
		judgeTestFinding("openai.ai.review.1", "internal/app/app.go"),
		judgeTestFinding("anthropic.ai.review.1", "internal/app/app.go"),
	})
	engine.judge = failingJudge{}

	report, err := engine.Review(context.Background(), root, Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Findings) != 1 {
		t.Fatalf("Findings = %#v, want near-duplicates deduped even though the judge call failed", report.Findings)
	}
}

func TestApplyJudgeResultsPrecisionFilter(t *testing.T) {
	findings := []Finding{
		{ID: "f.break1"}, {ID: "f.break2"}, {ID: "f.break3"},
		{ID: "f.func_hi"}, {ID: "f.func_lo"},
		{ID: "f.cosmetic"}, {ID: "f.none"}, {ID: "f.wrong"},
	}
	results := []judgeResult{
		{CandidateID: "f.break1", Verdict: "confirmed", Impact: "breaking", Severity: 3, Confidence: 0.7},
		{CandidateID: "f.break2", Verdict: "confirmed", Impact: "breaking", Severity: 5, Confidence: 0.9},
		{CandidateID: "f.break3", Verdict: "confirmed", Impact: "breaking", Severity: 2, Confidence: 0.6},
		{CandidateID: "f.func_hi", Verdict: "confirmed", Impact: "functional", Severity: 4, Confidence: 0.8},
		{CandidateID: "f.func_lo", Verdict: "confirmed", Impact: "functional", Severity: 4, Confidence: 0.3},
		{CandidateID: "f.cosmetic", Verdict: "confirmed", Impact: "cosmetic", Severity: 2, Confidence: 0.9},
		{CandidateID: "f.none", Verdict: "confirmed", Impact: "none", Severity: 1, Confidence: 0.9},
		{CandidateID: "f.wrong", Verdict: "wrong", Impact: "breaking", Severity: 5, Confidence: 0.9},
	}

	got := applyJudgeResults(findings, results)

	// Surfaced: 3 breaking + 1 high-confidence functional (proves no cap-to-3).
	// Dropped: low-confidence functional (noise), cosmetic + none (confidently
	// benign), and the non-confirmed candidate. Order: breaking by confidence,
	// then functional.
	if gotIDs := findingIDs(got); gotIDs != "f.break2,f.break1,f.break3,f.func_hi" {
		t.Fatalf("applyJudgeResults order = %q, want breaking (by confidence) then functional; benign/low-confidence/non-confirmed dropped", gotIDs)
	}
	byID := map[string]Finding{}
	for _, f := range got {
		byID[f.ID] = f
	}
	if byID["f.break2"].Strength != "Strong" {
		t.Fatalf("breaking strength = %q, want Strong", byID["f.break2"].Strength)
	}
	if byID["f.func_hi"].Strength != "Worth exploring" {
		t.Fatalf("functional strength = %q, want Worth exploring", byID["f.func_hi"].Strength)
	}
}

func TestCapKeepsTopThreeAdvisoryFindings(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "0")
	root := t.TempDir()
	findings := []Finding{
		judgeTestFindingWithStrength("advice.5", "Worth exploring"),
		judgeTestFindingWithStrength("advice.4", "Strong"),
		judgeTestFindingWithStrength("advice.3", "Worth exploring"),
		judgeTestFindingWithStrength("advice.2", "Strong"),
		judgeTestFindingWithStrength("advice.1", "Worth exploring"),
	}
	engine := judgeTestEngine(root, findings)

	report, err := engine.Review(context.Background(), root, Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if got := findingIDs(report.Findings); got != "advice.2,advice.4,advice.1" {
		t.Fatalf("Findings = %s, want top three advisories", got)
	}
}

func TestBlockingToolFindingsRenderInAdditionToCap(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "0")
	root := t.TempDir()
	findings := []Finding{{
		ID:             "tools.static-failure",
		Scopes:         []string{"architecture", "testing", "maintainability"},
		Title:          "Fix static tool failures before architecture advice",
		Summary:        "`go test ./...` failed.",
		Recommendation: "Fix the failing test output.",
		Strength:       "Blocking",
	}}
	for i := 1; i <= 4; i++ {
		findings = append(findings, judgeTestFindingWithStrength("advice."+string(rune('0'+i)), "Worth exploring"))
	}
	engine := judgeTestEngine(root, findings)

	report, err := engine.Review(context.Background(), root, Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Findings) != 4 || report.Findings[0].ID != "tools.static-failure" {
		t.Fatalf("Findings = %#v, want blocking tool plus three advisory findings", report.Findings)
	}
	text := RenderMarkdown(report)
	if !strings.Contains(text, "Fix static tool failures") || !strings.Contains(text, "### 4.") {
		t.Fatalf("RenderMarkdown() missing blocking/capped findings:\n%s", text)
	}
}

func TestJudgeDisabledStillCapsAdvisoryFindingsAtThree(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "0")
	root := t.TempDir()
	engine := judgeTestEngine(root, []Finding{
		judgeTestFindingWithStrength("advice.1", "Worth exploring"),
		judgeTestFindingWithStrength("advice.2", "Worth exploring"),
		judgeTestFindingWithStrength("advice.3", "Worth exploring"),
		judgeTestFindingWithStrength("advice.4", "Worth exploring"),
	})
	engine.judge = scriptedJudge{results: []judgeResult{{CandidateID: "advice.4", Verdict: "confirmed", Severity: 5, Confidence: 1}}}

	report, err := engine.Review(context.Background(), root, Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Findings) != 3 || hasFinding(report.Findings, "advice.4") {
		t.Fatalf("Findings = %#v, want injected judge ignored and advisory cap applied", report.Findings)
	}
}

func TestNearDupeProviderFindingsMerge(t *testing.T) {
	findings := mergeNearDuplicateFindings([]Finding{
		{
			ID:             "openai.ai.review.1",
			Title:          "Preserve auth rollback test coverage",
			Summary:        "`internal/auth/session.go` changes rollback behavior without a focused test.",
			Recommendation: "Add a rollback test in `internal/auth/session_test.go`.",
		},
		{
			ID:             "anthropic.ai.review.1",
			Title:          "Preserve auth rollback test coverage",
			Summary:        "`internal/auth/session.go` changes rollback behavior without focused coverage.",
			Recommendation: "Add the missing rollback coverage.",
		},
	})
	if len(findings) != 1 {
		t.Fatalf("mergeNearDuplicateFindings() len = %d, want 1", len(findings))
	}
	if !strings.Contains(evidenceText(findings[0].Evidence), "Similar findings merged") {
		t.Fatalf("Evidence = %#v, want agreement metadata", findings[0].Evidence)
	}
}

func TestNoCredentialEnvironmentDoesNotAttemptJudgeViaUnavailablePlaceholders(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "1")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GX_OPENAI_API_KEY", "")
	t.Setenv("GX_OPENAI_PROXY_URL", "")
	t.Setenv("GX_CLOUD_URL", "off")

	judge := judgeFromEnvWithPolicy(nil)
	if judgeAvailable(judge) {
		t.Fatalf("judgeAvailable(%T) = true, want false without credentials", judge)
	}
}

func judgeTestEngine(root string, findings []Finding) *Engine {
	return NewEngineWith(
		fakeScanner{facts: RepoFacts{
			Files:            []string{"internal/app/app.go"},
			Docs:             []FilePresence{{Path: "README.md", Present: true}},
			DependencyFiles:  []string{"go.mod"},
			TestFileCount:    1,
			TrackedFileCount: 2,
		}},
		fakeCatalog{},
		[]Rule{findingRule{findings: findings}},
	)
}

func judgeTestFinding(id string, file string) Finding {
	return Finding{
		ID:             id,
		Scopes:         []string{"architecture", "testing", "maintainability"},
		Title:          "Review changed file",
		Summary:        "The change in `" + file + "` needs a focused review.",
		Recommendation: "Update `" + file + "` and verify it with a focused test.",
		Strength:       "Worth exploring",
	}
}

func judgeTestFindingWithStrength(id string, strength string) Finding {
	finding := judgeTestFinding(id, "internal/app/app.go")
	topic := map[string]string{
		"advice.1": "authorization",
		"advice.2": "rollback",
		"advice.3": "observability",
		"advice.4": "idempotency",
		"advice.5": "validation",
	}[id]
	if topic == "" {
		topic = strings.ReplaceAll(id, ".", "-")
	}
	finding.Title = "Advisory " + topic
	finding.Summary = "The change in `internal/app/app.go` needs focused " + topic + " review."
	finding.Recommendation = "Update `internal/app/app.go` for " + topic + " and verify it with a focused test."
	finding.Strength = strength
	return finding
}

func findingIDs(findings []Finding) string {
	var ids []string
	for _, finding := range findings {
		ids = append(ids, finding.ID)
	}
	return strings.Join(ids, ",")
}

type failingJudge struct{}

func (failingJudge) Available() bool {
	return true
}

func (failingJudge) Judge(context.Context, judgeRequest) ([]judgeResult, error) {
	return nil, errors.New("judge request failed")
}

type scriptedJudge struct {
	results []judgeResult
	calls   int
}

func (j scriptedJudge) Available() bool {
	return true
}

func (j scriptedJudge) Judge(context.Context, judgeRequest) ([]judgeResult, error) {
	return j.results, nil
}

type findingRule struct {
	findings []Finding
}

func (r findingRule) ID() string {
	return "finding.rule"
}

func (r findingRule) Scopes() []string {
	return []string{"architecture", "testing", "maintainability"}
}

func (r findingRule) Evaluate(ReviewContext) []Finding {
	return r.findings
}
