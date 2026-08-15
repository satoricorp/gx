package codereview

import (
	"context"
	"strings"
	"testing"
)

func TestTriageDocsOnly(t *testing.T) {
	triage := TriageChange([]string{"README.md", "docs/overview.md", "CODEOWNERS"}, nil, normalizeOptions(Options{}))
	if triage.Class != "docs-only" {
		t.Fatalf("Class = %q, want docs-only", triage.Class)
	}
}

func TestTriageTestsOnly(t *testing.T) {
	triage := TriageChange([]string{"internal/app/app_test.go", "testdata/fixtures/input.json"}, nil, normalizeOptions(Options{}))
	if triage.Class != "tests-only" {
		t.Fatalf("Class = %q, want tests-only", triage.Class)
	}
}

func TestTriageConfigOnly(t *testing.T) {
	triage := TriageChange([]string{".github/workflows/ci.yml", "go.sum"}, nil, normalizeOptions(Options{}))
	if triage.Class != "config-only" {
		t.Fatalf("Class = %q, want config-only", triage.Class)
	}
}

func TestTriageMechanicalGenerated(t *testing.T) {
	triage := TriageChange([]string{"api/v1/service.pb.go"}, nil, normalizeOptions(Options{}))
	if triage.Class != "mechanical" {
		t.Fatalf("Class = %q, want mechanical", triage.Class)
	}
}

func TestTriagePolicyDocsAreConfigNotDocs(t *testing.T) {
	triage := TriageChange([]string{"REVIEW.md", "AGENTS.md", "SECURITY.md"}, nil, normalizeOptions(Options{}))
	if triage.Class != "config-only" {
		t.Fatalf("Class = %q, want config-only", triage.Class)
	}
}

func TestTriageRegexDiffYieldsRegexAndInjection(t *testing.T) {
	triage := TriageChange([]string{"internal/search/filter.go"}, []DiffSnippet{{
		File: "internal/search/filter.go",
		Diff: "+matcher := regexp.MustCompile(input)\n",
	}}, normalizeOptions(Options{}))
	if triage.Class != "security-sensitive" {
		t.Fatalf("Class = %q, want security-sensitive", triage.Class)
	}
	for _, want := range []string{"regex", "injection"} {
		if !hasString(triage.RiskTags, want) {
			t.Fatalf("RiskTags = %#v, want %q", triage.RiskTags, want)
		}
	}
}

func TestPatchFocusedScopeListFollowsSecuritySensitiveTriage(t *testing.T) {
	opts := normalizeOptions(Options{})
	triage := ChangeTriage{Class: "security-sensitive", RiskTags: []string{"regex"}}
	got := strings.Join(activeScopeListForTriage(opts, triage), ",")
	if got != "security,dependencies,testing,maintainability" {
		t.Fatalf("active scopes = %q", got)
	}
}

func TestDocsOnlyDefaultReviewShortCircuits(t *testing.T) {
	root := initRepo(t)
	writeFile(t, root, "README.md", "# old\n")
	gitAdd(t, root, "README.md")
	gitCommit(t, root)
	writeFile(t, root, "README.md", "# new\n")

	engine := NewEngineWithReviewer(LocalScanner{}, StaticCatalog{}, defaultRules(), panicRetriever{}, panicReviewer{})
	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Findings) != 0 {
		t.Fatalf("Findings = %#v, want zero", report.Findings)
	}
	if report.Triage.Class != "docs-only" {
		t.Fatalf("Triage = %#v, want docs-only", report.Triage)
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

func TestReadmePlusCodeDoesNotShortCircuit(t *testing.T) {
	root := initRepo(t)
	writeFile(t, root, "README.md", "# old\n")
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	gitAdd(t, root, "README.md", "internal/app/app.go")
	gitCommit(t, root)
	writeFile(t, root, "README.md", "# new\n")
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\nfunc More() {}\n")

	retriever := &recordingRetriever{}
	engine := NewEngineWithReviewer(LocalScanner{}, StaticCatalog{}, defaultRules(), retriever, nil)
	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.NoFindingsMessage != "" {
		t.Fatalf("NoFindingsMessage = %q, want empty", report.NoFindingsMessage)
	}
	if retriever.calls == 0 {
		t.Fatal("retriever was not called")
	}
}

func TestPromptedDocsOnlyReviewBypassesShortCircuit(t *testing.T) {
	assertDocsOnlyBypass(t, Options{Prompt: "review docs accuracy"})
}

func TestDeepDocsOnlyReviewBypassesShortCircuit(t *testing.T) {
	assertDocsOnlyBypass(t, Options{Deep: true})
}

func TestExplicitScopeDocsOnlyReviewBypassesShortCircuit(t *testing.T) {
	assertDocsOnlyBypass(t, Options{Scope: "architecture"})
}

func TestEnhanceResourceSignalsUseRetrieveInputChangedFilesAndPlanRiskTags(t *testing.T) {
	signals := reviewResourceSignals(RetrieveInput{
		RepoRoot:     "/missing-repo",
		Options:      normalizeOptions(Options{}),
		Facts:        RepoFacts{Files: []string{"docs/overview.md"}, DependencyFiles: []string{"go.mod"}},
		ChangedFiles: []string{"internal/search/filter.go"},
		Plan:         ReviewExecutionPlan{RiskTags: []string{"regex"}},
	})
	if len(signals.Files) != 1 || signals.Files[0] != "internal/search/filter.go" {
		t.Fatalf("Files = %#v, want changed files from RetrieveInput", signals.Files)
	}
	if len(signals.RiskTags) != 1 || signals.RiskTags[0] != "regex" {
		t.Fatalf("RiskTags = %#v, want plan risk tags", signals.RiskTags)
	}
}

func assertDocsOnlyBypass(t *testing.T, opts Options) {
	t.Helper()
	root := initRepo(t)
	writeFile(t, root, "README.md", "# old\n")
	gitAdd(t, root, "README.md")
	gitCommit(t, root)
	writeFile(t, root, "README.md", "# new\n")

	retriever := &recordingRetriever{}
	engine := NewEngineWithReviewer(LocalScanner{}, StaticCatalog{}, defaultRules(), retriever, nil)
	report, err := engine.Review(context.Background(), root, opts)
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if retriever.calls == 0 {
		t.Fatalf("retriever was not called for opts %#v", opts)
	}
	if report.NoFindingsMessage != "" {
		t.Fatalf("NoFindingsMessage = %q, want empty for bypass", report.NoFindingsMessage)
	}
}

type recordingRetriever struct {
	calls int
	input RetrieveInput
}

func (r *recordingRetriever) Retrieve(_ context.Context, in RetrieveInput) ([]ContextSnippet, error) {
	r.calls++
	r.input = in
	return nil, nil
}

type panicRetriever struct{}

func (panicRetriever) Retrieve(context.Context, RetrieveInput) ([]ContextSnippet, error) {
	panic("retriever should not run")
}

type panicReviewer struct{}

func (panicReviewer) Review(context.Context, ReviewBrief) ([]Finding, error) {
	panic("reviewer should not run")
}

func hasString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
