package codereview

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestEnhanceerAvailable(t *testing.T) {
	if reviewerAvailable(nil) {
		t.Fatal("reviewerAvailable(nil) = true, want false")
	}
	if reviewerAvailable(multiAIReviewer{reviewers: []namedAIReviewer{{reviewer: unavailableAIReviewer{reason: "missing"}}}}) {
		t.Fatal("reviewerAvailable(unavailable multi) = true, want false")
	}
	// A leg is available when it has a wire to Bedrock, not merely when it
	// exists: a reviewer with no transport cannot review, and reporting it as
	// available is how a review with no reviewer reads as a clean review.
	if reviewerAvailable(&bedrockAnthropicReviewer{model: defaultBedrockReviewModelA}) {
		t.Fatal("reviewerAvailable(Bedrock reviewer without a transport) = true, want false")
	}
	withTransport := newBedrockReviewer(
		newDirectBedrockTransport(bedrockCredentials{accessKey: "key", secretKey: "secret", region: "us-west-2"}),
		defaultBedrockReviewModelA,
	)
	if !reviewerAvailable(withTransport) {
		t.Fatal("reviewerAvailable(real Bedrock reviewer) = false, want true")
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

// TestApplyJudgeResultsRecordsRankWithoutReordering pins the measured decision:
// the judge's batch-relative rank is carried onto the finding for offline
// analysis, but the report order stays impact/confidence — on the 30-PR
// benchmark, ordering by the judge's stated rank scored BELOW ordering by its
// own severity+confidence at every top-K (39 vs 42 F1 at top-2).
func TestApplyJudgeResultsRecordsRankWithoutReordering(t *testing.T) {
	findings := []Finding{
		{ID: "f.a"}, {ID: "f.b"}, {ID: "f.c"},
	}
	results := []judgeResult{
		// Rank inverts what impact+confidence would produce; ordering must
		// ignore it and follow impact then confidence anyway.
		{CandidateID: "f.a", Verdict: "confirmed", Impact: "breaking", Severity: 5, Confidence: 0.95, Rank: 3},
		{CandidateID: "f.b", Verdict: "confirmed", Impact: "functional", Severity: 3, Confidence: 0.8, Rank: 1},
		{CandidateID: "f.c", Verdict: "confirmed", Impact: "breaking", Severity: 4, Confidence: 0.9, Rank: 2},
	}
	got := applyJudgeResults(findings, results)
	if gotIDs := findingIDs(got); gotIDs != "f.a,f.c,f.b" {
		t.Fatalf("order = %q, want impact then confidence (a,c,b), ignoring the judge's rank", gotIDs)
	}
	if got[0].JudgeRank != 3 || got[2].JudgeRank != 1 {
		t.Fatalf("JudgeRank not carried: got %d,%d,%d", got[0].JudgeRank, got[1].JudgeRank, got[2].JudgeRank)
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

	report, err := engine.Review(context.Background(), root, Options{Scope: "architecture", MaxFindings: 3})
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

	report, err := engine.Review(context.Background(), root, Options{Scope: "architecture", MaxFindings: 3})
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

	report, err := engine.Review(context.Background(), root, Options{Scope: "architecture", MaxFindings: 3})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Findings) != 3 || hasFinding(report.Findings, "advice.4") {
		t.Fatalf("Findings = %#v, want injected judge ignored and advisory cap applied", report.Findings)
	}
}

// TestNearDupeProviderFindingsMerge covers the one merge the deterministic side
// still makes on its own: two legs that returned the same finding almost word
// for word. It passes a nil adjudicator on purpose — near-verbatim copies must
// collapse with no model reachable at all.
func TestNearDupeProviderFindingsMerge(t *testing.T) {
	findings := mergeNearDuplicateFindings(context.Background(), nil, []Finding{
		{
			ID:             "openai.ai.review.1",
			Title:          "Preserve auth rollback test coverage",
			Summary:        "`internal/auth/session.go` changes rollback behavior without a focused test.",
			Recommendation: "Add a rollback test in `internal/auth/session_test.go`.",
			Corroboration:  []string{"OpenAI"},
		},
		{
			ID:             "anthropic.ai.review.1",
			Title:          "Preserve auth rollback test coverage",
			Summary:        "`internal/auth/session.go` changes rollback behavior without focused coverage.",
			Recommendation: "Add the missing rollback coverage.",
			Corroboration:  []string{"Anthropic"},
		},
	})
	if len(findings) != 1 {
		t.Fatalf("mergeNearDuplicateFindings() len = %d, want 1", len(findings))
	}
	if !strings.Contains(evidenceText(findings[0].Evidence), "Raised independently by 2 reviewers") {
		t.Fatalf("Evidence = %#v, want cross-reviewer corroboration metadata", findings[0].Evidence)
	}
}

func TestNoCredentialEnvironmentDoesNotAttemptJudgeViaUnavailablePlaceholders(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "1")
	t.Setenv("GX_CLOUD_URL", "off")
	// An OPENAI_API_KEY must no longer produce a judge: the OpenAI judge is
	// gone, and embeddings keep that variable set on most developer machines.
	t.Setenv("OPENAI_API_KEY", "sk-embeddings-only")
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")

	judge := judgeFromEnv()
	if judgeAvailable(judge) {
		t.Fatalf("judgeAvailable(%T) = true, want false with no reviewer configured", judge)
	}
	reason, ok := judge.(unavailableReviewJudge)
	if !ok || !strings.Contains(reason.reason, "Cloud") {
		t.Fatalf("judge = %#v, want an unavailable judge naming the Cloud fix", judge)
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

// judgeTestFindingWithStrength builds a distinct advisory finding.
//
// "Distinct" has to mean it, and it has to mean it on the one file these tests
// are about. The fixtures used to differ from each other by a single word in an
// otherwise identical sentence, which no honest de-duplication can be expected
// to keep apart; an earlier attempt at fixing that gave each fixture its own
// file, which kept them apart only by defeating the locality gate — a fixture
// arranged so the merge rule never runs proves nothing about the merge rule.
// They now share `internal/app/app.go`, as a cap test about one changed file
// should, and are five genuinely different problems in it.
func judgeTestFindingWithStrength(id string, strength string) Finding {
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
	title := map[string]string{
		"authorization": "Handler acts on another tenant's records without a permission check",
		"rollback":      "Commit happens before the downstream write is confirmed",
		"observability": "Failures on the Stop path leave no log line or metric",
		"idempotency":   "Retrying a request repeats the side effect instead of replaying it",
		"validation":    "Request body reaches storage with no bounds or type checking",
	}[topic]
	if title == "" {
		title = "Advisory " + topic
	}
	body := map[string]string{
		"authorization": "Any signed-in account can reach Run() and mutate records belonging to a different tenant, because nothing between the router and the store asserts ownership.",
		"rollback":      "The database transaction commits before the downstream write is acknowledged, so a half-applied change survives with no way to undo it.",
		"observability": "Stop() emits nothing at all, so a production hang is invisible until a customer notices and reports it.",
		"idempotency":   "A retried request is treated as a fresh one, repeating the side effect and double-charging whenever the network is flaky.",
		"validation":    "Input is trusted straight from the request body and reaches storage with no bounds or type checking.",
	}[topic]
	if body == "" {
		body = "This path behaves differently from the rest of the package for reasons nothing records."
	}
	fix := map[string]string{
		"authorization": "Assert tenant ownership in Run() before touching the store.",
		"rollback":      "Move the commit after the downstream acknowledgement, or make the downstream write part of the same transaction.",
		"observability": "Emit a structured log line and a counter on the Stop() failure path.",
		"idempotency":   "Key the side effect on a request ID and replay the stored result on retry.",
		"validation":    "Parse and bound-check the request body before it reaches storage.",
	}[topic]
	if fix == "" {
		fix = "Bring this path in line with the rest of the package."
	}
	return Finding{
		ID:             id,
		Scopes:         []string{"architecture", "testing", "maintainability"},
		Title:          title,
		Summary:        body + " See `internal/app/app.go`.",
		Recommendation: fix,
		Strength:       strength,
	}
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

// constantAdjudicator answers every question the same way. It exists to show
// that a fixture reaches the merge machinery at all, which is the thing a
// fixture can quietly stop doing.
type constantAdjudicator bool

func (a constantAdjudicator) AdjudicateDuplicates(_ context.Context, pairs []duplicatePairInput) ([]duplicateVerdict, error) {
	out := make([]duplicateVerdict, 0, len(pairs))
	for _, pair := range pairs {
		out = append(out, duplicateVerdict{PairID: pair.PairID, Same: bool(a)})
	}
	return out, nil
}

// TestAdvisoryCapFixturesAreDistinctOnTheSameFile guards the fixtures, not the
// product.
//
// These five advisories exist to be capped, so they must survive
// de-duplication — but they must survive it for the right reason. An earlier
// revision gave each of them its own file, which made them survive by defeating
// the locality gate: no pair could ever be considered, and the cap test would
// have gone on passing however badly the merge behaved. They share the one file
// the test is about again, and they survive because they describe five
// different problems.
//
// The last assertion is the one that keeps this honest: a near-verbatim sixth
// copy still collapses, so de-duplication is demonstrably live on exactly this
// fixture set rather than inert.
func TestAdvisoryCapFixturesAreDistinctOnTheSameFile(t *testing.T) {
	var findings []Finding
	for _, id := range []string{"advice.1", "advice.2", "advice.3", "advice.4", "advice.5"} {
		findings = append(findings, judgeTestFindingWithStrength(id, "Worth exploring"))
	}
	for _, finding := range findings {
		files := findingFiles(finding)
		if _, ok := files["internal/app/app.go"]; !ok {
			t.Fatalf("%q names %v, not the one file this test is about", finding.Title, sortedSet(files))
		}
	}
	worst := 0.0
	for i := range findings {
		for j := i + 1; j < len(findings); j++ {
			if score := pairScore(findings[i], findings[j]); score > worst {
				worst = score
			}
		}
	}
	t.Logf("most similar pair of the five scores %.4f (gate %.2f)", worst, dedupeGateThreshold)

	for name, adjudicator := range map[string]duplicateAdjudicator{
		"no adjudicator":     nil,
		"honest adjudicator": constantAdjudicator(false),
	} {
		if kept := mergeNearDuplicateFindings(context.Background(), adjudicator, findings); len(kept) != 5 {
			t.Fatalf("%s: kept %d of 5 distinct advisories:\n%s", name, len(kept), findingIDs(kept))
		}
	}

	duplicated := append(append([]Finding(nil), findings...), func() Finding {
		copied := findings[0]
		copied.ID = "advice.6"
		return copied
	}())
	if kept := mergeNearDuplicateFindings(context.Background(), nil, duplicated); len(kept) != 5 {
		t.Fatalf("a verbatim sixth copy did not collapse (%d kept), so de-duplication is inert on this fixture set:\n%s",
			len(kept), findingIDs(kept))
	}
}
