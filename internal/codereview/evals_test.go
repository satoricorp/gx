// Eval suite for gx review end-to-end behavior.
//
// Each TestEval* scenario builds a synthetic git repository, applies a realistic
// change, and runs the review engine with scripted reviewers/judges (never real
// APIs). Run the suite with:
//
//	go test ./internal/codereview/ -run TestEval -v
package codereview

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func evalDisableAllGates(t *testing.T) {
	t.Helper()
	t.Setenv("GX_REVIEW_AI", "0")
	t.Setenv("GX_REVIEW_JUDGE", "0")
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "0")
	t.Setenv("GX_REVIEW_RESOURCES", "0")
	t.Setenv("GX_REVIEW_INDEXED_CONTEXT", "0")
	t.Setenv("GX_REVIEW_HISTORY_CONTEXT", "0")
}

func evalEngineWithReviewer(t *testing.T, reviewer AIReviewer, retriever ContextRetriever) *Engine {
	t.Helper()
	evalDisableAllGates(t)
	return NewEngineWithReviewer(
		LocalScanner{},
		StaticCatalog{},
		nil,
		retriever,
		reviewer,
	)
}

func evalEngineWithRules(t *testing.T, rules []Rule, reviewer AIReviewer, retriever ContextRetriever) *Engine {
	t.Helper()
	evalDisableAllGates(t)
	return NewEngineWithReviewer(LocalScanner{}, StaticCatalog{}, rules, retriever, reviewer)
}

func evalCommitChange(t *testing.T, root string, initial map[string]string, changes map[string]string) {
	t.Helper()
	for path, content := range initial {
		writeFile(t, root, path, content)
	}
	if len(initial) > 0 {
		paths := make([]string, 0, len(initial))
		for path := range initial {
			paths = append(paths, path)
		}
		gitAdd(t, root, paths...)
		gitCommit(t, root)
	}
	for path, content := range changes {
		writeFile(t, root, path, content)
	}
}

func evalAssertRiskTags(t *testing.T, triage ChangeTriage, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !hasString(triage.RiskTags, want) {
			t.Fatalf("RiskTags = %#v, want %q", triage.RiskTags, want)
		}
	}
}

func evalAssertNoRiskTag(t *testing.T, triage ChangeTriage, tag string) {
	t.Helper()
	if hasString(triage.RiskTags, tag) {
		t.Fatalf("RiskTags = %#v, did not want %q", triage.RiskTags, tag)
	}
}

func evalCountRenderedRecommendations(text string) int {
	count := 0
	for line := range strings.SplitSeq(text, "\n") {
		if strings.HasPrefix(line, "### ") {
			count++
		}
	}
	return count
}

func evalAdvisoryFindings(findings []Finding) []Finding {
	var out []Finding
	for _, finding := range findings {
		if finding.Strength == "Blocking" || strings.HasPrefix(finding.ID, "tools.") {
			continue
		}
		out = append(out, finding)
	}
	return out
}

func TestEvalDocsOnlyReadmeChangelog(t *testing.T) {
	evalDisableAllGates(t)
	root := initRepo(t)
	evalCommitChange(t, root, map[string]string{
		"README.md":    "# Product\n",
		"CHANGELOG.md": "## 1.0.0\nInitial release.\n",
	}, map[string]string{
		"README.md":    "# Product\n\nUpdated overview.\n",
		"CHANGELOG.md": "## 1.0.1\nDocumentation refresh.\n",
	})

	engine := evalEngineWithReviewer(t, panicReviewer{}, panicRetriever{})
	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.Triage.Class != "docs-only" {
		t.Fatalf("Triage.Class = %q, want docs-only", report.Triage.Class)
	}
	if len(report.Findings) != 0 {
		t.Fatalf("Findings = %#v, want zero", report.Findings)
	}
	text := RenderMarkdown(report)
	if !strings.Contains(text, "No material issues in this change (documentation-only).") {
		t.Fatalf("RenderMarkdown() missing docs-only message:\n%s", text)
	}
	for _, unwanted := range []string{"## Sources", "## Context Sources"} {
		if strings.Contains(text, unwanted) {
			t.Fatalf("RenderMarkdown() should not include %q:\n%s", unwanted, text)
		}
	}
}

func TestEvalMechanicalCommentWhitespaceReformat(t *testing.T) {
	evalDisableAllGates(t)
	root := initRepo(t)
	initial := "package service\n\nfunc Handle() {\n  x := 1\n  y := 2\n}\n\nfunc Other() {\n  a := 3\n  b := 4\n}\n"
	changed := "package service\n\nfunc Handle() {\n    x := 1\n    y := 2\n}\n\nfunc Other() {\n    a := 3\n    b := 4\n}\n"
	evalCommitChange(t, root, map[string]string{"internal/service/handler.go": initial}, map[string]string{
		"internal/service/handler.go": changed,
	})

	engine := evalEngineWithReviewer(t, nil, fakeRetriever{})
	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.Triage.Class != "mechanical" {
		t.Fatalf("Triage.Class = %q, want mechanical", report.Triage.Class)
	}
	scopes := activeScopeListForTriage(normalizeOptions(Options{}), report.Triage)
	if hasString(scopes, "security") {
		t.Fatalf("active scopes = %#v, want no security scope for mechanical change", scopes)
	}
}

func TestEvalRegexUserInputSecuritySensitive(t *testing.T) {
	evalDisableAllGates(t)
	root := initRepo(t)
	evalCommitChange(t, root, map[string]string{
		"internal/search/filter.go": "package search\n\nfunc Filter() {}\n",
	}, map[string]string{
		"internal/search/filter.go": strings.Join([]string{
			"package search",
			"",
			"import (",
			"\"net/http\"",
			"\"regexp\"",
			")",
			"",
			"func Filter(w http.ResponseWriter, r *http.Request) {",
			"  pattern := r.URL.Query().Get(\"q\")",
			"  matcher := regexp.MustCompile(pattern)",
			"  _ = matcher",
			"}",
		}, "\n"),
	})

	engine := evalEngineWithReviewer(t, nil, fakeRetriever{})
	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.Triage.Class != "security-sensitive" {
		t.Fatalf("Triage.Class = %q, want security-sensitive", report.Triage.Class)
	}
	evalAssertRiskTags(t, report.Triage, "regex", "injection")
	scopes := activeScopeListForTriage(normalizeOptions(Options{}), report.Triage)
	if !hasString(scopes, "security") {
		t.Fatalf("active scopes = %#v, want security scope active", scopes)
	}
}

func TestEvalSQLStringConcatenation(t *testing.T) {
	evalDisableAllGates(t)
	root := initRepo(t)
	evalCommitChange(t, root, map[string]string{
		"internal/store/users.go": "package store\n\nfunc Lookup(id string) error { return nil }\n",
	}, map[string]string{
		"internal/store/users.go": strings.Join([]string{
			"package store",
			"",
			"import \"fmt\"",
			"",
			"func Lookup(id string) error {",
			"  query := fmt.Sprintf(\"SELECT id, name FROM users WHERE id=%s\", id)",
			"  _ = query",
			"  return nil",
			"}",
		}, "\n"),
	})

	engine := evalEngineWithReviewer(t, nil, fakeRetriever{})
	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.Triage.Class != "security-sensitive" {
		t.Fatalf("Triage.Class = %q, want security-sensitive", report.Triage.Class)
	}
	evalAssertRiskTags(t, report.Triage, "sql-injection", "injection")
}

func TestEvalSecurityFallbackWithoutAI(t *testing.T) {
	evalDisableAllGates(t)
	root := initRepo(t)
	evalCommitChange(t, root, map[string]string{
		"internal/store/users.go": "package store\n\nfunc Lookup(id string) error { return nil }\n",
	}, map[string]string{
		"internal/store/users.go": strings.Join([]string{
			"package store",
			"",
			"import \"fmt\"",
			"",
			"func Lookup(id string) error {",
			"  query := fmt.Sprintf(\"SELECT id, name FROM users WHERE id=%s\", id)",
			"  _ = query",
			"  return nil",
			"}",
		}, "\n"),
	})

	engine := evalEngineWithReviewer(t, nil, fakeRetriever{})
	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if !hasFinding(report.Findings, "security.quality-hints.1") {
		t.Fatalf("Findings = %#v, want deterministic SQLi fallback finding", report.Findings)
	}
	text := RenderMarkdown(report)
	if strings.Contains(text, "> Warning: AI review unavailable") {
		t.Fatalf("RenderMarkdown() should not warn when AI is intentionally disabled:\n%s", text)
	}
}

func TestEvalShellExecUserInput(t *testing.T) {
	evalDisableAllGates(t)
	root := initRepo(t)
	evalCommitChange(t, root, map[string]string{
		"internal/tools/runner.go": "package tools\n\nfunc Run(input string) error { return nil }\n",
	}, map[string]string{
		"internal/tools/runner.go": strings.Join([]string{
			"package tools",
			"",
			"import \"os/exec\"",
			"",
			"func Run(input string) error {",
			"  return exec.Command(\"sh\", \"-c\", input).Run()",
			"}",
		}, "\n"),
	})

	engine := evalEngineWithReviewer(t, nil, fakeRetriever{})
	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.Triage.Class != "security-sensitive" {
		t.Fatalf("Triage.Class = %q, want security-sensitive", report.Triage.Class)
	}
	evalAssertRiskTags(t, report.Triage, "command-injection", "injection")
}

func TestEvalHardcodedCredential(t *testing.T) {
	evalDisableAllGates(t)
	root := initRepo(t)
	evalCommitChange(t, root, map[string]string{
		"internal/config/config.go": "package config\n\nconst APIKey = \"\"\n",
	}, map[string]string{
		"internal/config/config.go": strings.Join([]string{
			"package config",
			"",
			"const APIKey = \"sk-live-0123456789abcdef\"",
			"const password = \"super-secret-password\"",
		}, "\n"),
	})

	engine := evalEngineWithReviewer(t, nil, fakeRetriever{})
	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.Triage.Class != "security-sensitive" {
		t.Fatalf("Triage.Class = %q, want security-sensitive", report.Triage.Class)
	}
	evalAssertRiskTags(t, report.Triage, "secure-coding")
}

func TestEvalConfigOnlyLockfileCI(t *testing.T) {
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "1")
	evalDisableAllGates(t)
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "1")

	root := initRepo(t)
	evalCommitChange(t, root, map[string]string{
		"package.json":              "{\n  \"name\": \"example\"\n}\n",
		"package-lock.json":         "{\n  \"name\": \"example\",\n  \"lockfileVersion\": 2\n}\n",
		".github/workflows/ci.yml": "name: CI\non: [push]\njobs:\n  test:\n    runs-on: ubuntu-latest\n",
	}, map[string]string{
		"package-lock.json":         "{\n  \"name\": \"example\",\n  \"lockfileVersion\": 3\n}\n",
		".github/workflows/ci.yml": "name: CI\non: [push]\njobs:\n  test:\n    runs-on: ubuntu-latest\n    steps:\n      - run: npm test\n",
	})

	var progress strings.Builder
	engine := evalEngineWithReviewer(t, nil, fakeRetriever{})
	report, err := engine.Review(context.Background(), root, Options{ProgressWriter: &progress})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.Triage.Class != "config-only" {
		t.Fatalf("Triage.Class = %q, want config-only", report.Triage.Class)
	}
	text := RenderMarkdown(report)
	if strings.Contains(text, "No material issues in this change (documentation-only).") {
		t.Fatalf("RenderMarkdown() should not use docs-only message for config-only change:\n%s", text)
	}
	if strings.Contains(progress.String(), "Running go test") || strings.Contains(progress.String(), "Running go vet") {
		t.Fatalf("static tools should be skipped for non-Go config-only changes, progress:\n%s", progress.String())
	}
}

func TestEvalCapEnforcementWithBlockingToolFinding(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "1")
	t.Setenv("GX_REVIEW_AI", "1")
	evalDisableAllGates(t)
	t.Setenv("GX_REVIEW_JUDGE", "1")
	t.Setenv("GX_REVIEW_AI", "1")

	root := initRepo(t)
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	gitAdd(t, root, "internal/app/app.go")
	gitCommit(t, root)
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\nfunc Stop() {}\n")

	advisories := make([]Finding, 0, 6)
	judgeResults := make([]judgeResult, 0, 5)
	topics := []string{"authorization", "rollback", "observability", "idempotency", "validation", "serialization"}
	for i := 1; i <= 6; i++ {
		id := fmt.Sprintf("ai.review.%d", i)
		topic := topics[i-1]
		advisories = append(advisories, Finding{
			ID:             id,
			Scopes:         []string{"security", "maintainability"},
			Title:          "Advisory " + topic,
			Summary:        fmt.Sprintf("The change in `internal/app/app.go` needs focused %s review.", topic),
			Benefit:        "Improves safety.",
			Recommendation: fmt.Sprintf("Address %s in internal/app/app.go.", topic),
			Strength:       "Worth exploring",
		})
		if i <= 5 {
			judgeResults = append(judgeResults, judgeResult{
				CandidateID:      id,
				Verdict:          "confirmed",
				Severity:         6 - i,
				Confidence:       0.95,
				VerificationNote: "Confirmed by scripted judge.",
			})
		}
	}

	blocking := Finding{
		ID:             "tools.static-failure",
		Scopes:         []string{"testing", "maintainability", "dependencies", "security"},
		Title:          "Fix static tool failures before other advice",
		Summary:        "`go test ./internal/app` failed.",
		Recommendation: "Fix the failing test output.",
		Strength:       "Blocking",
	}

	engine := evalEngineWithRules(t, []Rule{findingRule{findings: []Finding{blocking}}}, fakeReviewer{findings: advisories}, fakeRetriever{})
	engine.judge = scriptedJudge{results: judgeResults}

	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(evalAdvisoryFindings(report.Findings)) != 3 {
		t.Fatalf("advisory Findings = %#v, want exactly 3 after cap", evalAdvisoryFindings(report.Findings))
	}
	if !hasFinding(report.Findings, "tools.static-failure") {
		t.Fatalf("Findings = %#v, want blocking tools.* finding present", report.Findings)
	}
	if report.Findings[0].ID != "tools.static-failure" {
		t.Fatalf("Findings order = %#v, want blocking finding first", report.Findings)
	}
	if got := findingIDs(evalAdvisoryFindings(report.Findings)); got != "ai.review.1,ai.review.2,ai.review.3" {
		t.Fatalf("capped advisory order = %q, want strongest three by judge severity", got)
	}

	text := RenderMarkdown(report)
	if evalCountRenderedRecommendations(text) != 4 {
		t.Fatalf("RenderMarkdown() recommendation count = %d, want 4 (1 blocking + 3 advisory):\n%s", evalCountRenderedRecommendations(text), text)
	}
	if !strings.Contains(text, "Fix static tool failures") {
		t.Fatalf("RenderMarkdown() missing blocking finding:\n%s", text)
	}
}

func TestEvalZeroFindingsHonesty(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "1")
	t.Setenv("GX_REVIEW_AI", "1")
	evalDisableAllGates(t)
	t.Setenv("GX_REVIEW_JUDGE", "1")
	t.Setenv("GX_REVIEW_AI", "1")

	root := initRepo(t)
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	gitAdd(t, root, "internal/app/app.go")
	gitCommit(t, root)
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\nfunc More() {}\n")

	engine := evalEngineWithRules(t, []Rule{}, fakeReviewer{findings: []Finding{{
		ID:             "ai.review.unverified",
		Scopes:         []string{"security"},
		Title:          "Speculative issue",
		Summary:        "Maybe a problem exists.",
		Benefit:        "Would improve safety if true.",
		Recommendation: "Investigate further.",
		Strength:       "Worth exploring",
	}}}, fakeRetriever{})
	engine.judge = scriptedJudge{results: []judgeResult{{
		CandidateID:      "ai.review.unverified",
		Verdict:          "wrong",
		Severity:         1,
		Confidence:       0.99,
		VerificationNote: "Not supported by the diff.",
	}}}

	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Findings) != 0 {
		t.Fatalf("Findings = %#v, want zero after judge rejects candidates", report.Findings)
	}
	text := RenderMarkdown(report)
	if !strings.Contains(text, "No material issues found in this change.") {
		t.Fatalf("RenderMarkdown() missing zero-findings honesty message:\n%s", text)
	}
	if strings.Contains(text, "### 1.") {
		t.Fatalf("RenderMarkdown() should not render numbered recommendations:\n%s", text)
	}
}

func TestEvalAttributionOpacity(t *testing.T) {
	report := Report{
		Findings: []Finding{{
			ID:               "security.injection",
			Title:            "Avoid dynamic SQL concatenation",
			Summary:          "User input is interpolated into a SQL string.",
			Benefit:          "Reduces injection risk.",
			Recommendation:   "Use parameterized queries.",
			Strength:         "Strong",
			SourcePublishers: []string{"OWASP"},
			SourceIDs:        []string{"owasp-sql-injection"},
		}},
		Sources: []Source{{
			ID:        "owasp-sql-injection",
			Title:     "SQL Injection Prevention Cheat Sheet",
			URL:       "https://owasp.org/www-community/attacks/SQL_Injection",
			Publisher: "OWASP",
		}},
	}

	defaultText := RenderMarkdown(report)
	if !strings.Contains(defaultText, "**Informed by:** OWASP") {
		t.Fatalf("default RenderMarkdown() missing publisher attribution:\n%s", defaultText)
	}
	for _, unwanted := range []string{"owasp.org", "https://", "## Sources", "SQL Injection Prevention Cheat Sheet"} {
		if strings.Contains(defaultText, unwanted) {
			t.Fatalf("default RenderMarkdown() leaked source detail %q:\n%s", unwanted, defaultText)
		}
	}

	report.Verbose = true
	verboseText := RenderMarkdown(report)
	for _, want := range []string{"**Informed by:** OWASP", "## Sources", "owasp.org", "https://", "SQL Injection Prevention Cheat Sheet"} {
		if !strings.Contains(verboseText, want) {
			t.Fatalf("verbose RenderMarkdown() missing %q:\n%s", want, verboseText)
		}
	}
}

func TestEvalPromptDirectedDocsOnlyBypass(t *testing.T) {
	evalDisableAllGates(t)
	root := initRepo(t)
	evalCommitChange(t, root, map[string]string{"README.md": "# old\n"}, map[string]string{"README.md": "# new\n"})

	retriever := &recordingRetriever{}
	engine := evalEngineWithReviewer(t, nil, retriever)
	report, err := engine.Review(context.Background(), root, Options{Prompt: "verify documentation accuracy and completeness"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.Triage.Class != "docs-only" {
		t.Fatalf("Triage.Class = %q, want docs-only", report.Triage.Class)
	}
	if retriever.calls == 0 {
		t.Fatal("retriever was not called; prompt-directed docs-only review should bypass short-circuit")
	}
	if report.NoFindingsMessage != "" {
		t.Fatalf("NoFindingsMessage = %q, want empty when bypassing docs-only short-circuit", report.NoFindingsMessage)
	}
}

func TestEvalNearDuplicateMerging(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "1")
	t.Setenv("GX_REVIEW_AI", "1")
	evalDisableAllGates(t)
	t.Setenv("GX_REVIEW_JUDGE", "1")
	t.Setenv("GX_REVIEW_AI", "1")

	root := initRepo(t)
	writeFile(t, root, "internal/auth/session.go", "package auth\nfunc Save() {}\n")
	gitAdd(t, root, "internal/auth/session.go")
	gitCommit(t, root)
	writeFile(t, root, "internal/auth/session.go", "package auth\nfunc Save() {}\nfunc Rollback() {}\n")

	dupes := []Finding{
		{
			ID:             "openai.ai.review.1",
			Scopes:         []string{"testing", "maintainability"},
			Title:          "Preserve auth rollback test coverage",
			Summary:        "`internal/auth/session.go` changes rollback behavior without a focused test.",
			Benefit:        "Prevents regressions.",
			Recommendation: "Add a rollback test in `internal/auth/session_test.go`.",
			Strength:       "Strong",
		},
		{
			ID:             "anthropic.ai.review.1",
			Scopes:         []string{"testing", "maintainability"},
			Title:          "Preserve auth rollback test coverage",
			Summary:        "`internal/auth/session.go` changes rollback behavior without focused coverage.",
			Benefit:        "Prevents regressions.",
			Recommendation: "Add the missing rollback coverage.",
			Strength:       "Worth exploring",
		},
	}

	engine := evalEngineWithRules(t, nil, fakeReviewer{findings: dupes}, fakeRetriever{})
	engine.judge = scriptedJudge{results: []judgeResult{{
		CandidateID:      "openai.ai.review.1",
		Verdict:          "confirmed",
		Severity:         5,
		Confidence:       0.92,
		VerificationNote: "Rollback path changed without test.",
	}}}

	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Findings) != 1 {
		t.Fatalf("Findings = %#v, want one merged finding", report.Findings)
	}
	if !strings.Contains(evidenceText(report.Findings[0].Evidence), "Similar findings merged") {
		t.Fatalf("Evidence = %#v, want near-duplicate agreement metadata", report.Findings[0].Evidence)
	}
}
