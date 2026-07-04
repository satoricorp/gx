package codereview

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("GX_REVIEW_AI", "0")
	_ = os.Setenv("GX_REVIEW_JUDGE", "0")
	_ = os.Setenv("GX_REVIEW_STATIC_TOOLS", "0")
	os.Exit(m.Run())
}

func TestReviewUsesDefaultsAndDetectsRepoFacts(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "README.md", "# repo\n")
	writeFile(t, root, "AGENTS.md", "# agents\n")
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "internal/app/app_test.go", "package app\n")

	report, err := Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.Scope != DefaultScope || report.Format != DefaultFormat {
		t.Fatalf("Review() scope/format = %q/%q, want defaults", report.Scope, report.Format)
	}
	if !strings.Contains(strings.Join(report.DependencyFiles, ","), "go.mod") {
		t.Fatalf("dependency files = %#v, want go.mod", report.DependencyFiles)
	}
	if report.TestFileCount != 1 {
		t.Fatalf("TestFileCount = %d, want 1", report.TestFileCount)
	}
	if got := strings.Join(report.BaselineScopes, ","); got != "dependencies,testing,maintainability" {
		t.Fatalf("BaselineScopes = %q", got)
	}
	if hasFinding(report.Findings, "architecture.missing-context") {
		t.Fatalf("Findings = %#v, did not expect missing context finding", report.Findings)
	}
}

func TestUnpromptedDefaultReviewStaysPatchFocused(t *testing.T) {
	opts := normalizeOptions(Options{})
	if !opts.PatchFocused {
		t.Fatalf("PatchFocused = false, want true")
	}
	if profile := reviewProfile(opts); profile != "patch_focused" {
		t.Fatalf("reviewProfile() = %q, want patch_focused", profile)
	}
}

func TestReviewEmitsProgress(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "README.md", "# repo\n")
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	var progress strings.Builder

	_, err := Review(context.Background(), root, Options{ProgressWriter: &progress})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	text := progress.String()
	for _, want := range []string{"Scanning repository", "Building review context", "Checking fallback review rules"} {
		if !strings.Contains(text, want) {
			t.Fatalf("progress missing %q in:\n%s", want, text)
		}
	}
}

func TestReviewUsesTrackedFilesAndSkipsNodeModules(t *testing.T) {
	root := initRepo(t)
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "node_modules/pkg/package.json", "{}\n")
	gitAdd(t, root, "go.mod")

	report, err := Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if strings.Join(report.DependencyFiles, ",") != "go.mod" {
		t.Fatalf("DependencyFiles = %#v, want only go.mod", report.DependencyFiles)
	}
}

func TestReviewDedupesPrimaryBaselineScope(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "README.md", "# repo\n")

	report, err := Review(context.Background(), root, Options{Scope: "testing"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if got := strings.Join(report.BaselineScopes, ","); got != "dependencies,maintainability" {
		t.Fatalf("BaselineScopes = %q, want dependencies,maintainability", got)
	}
}

func TestReviewFocusLimitsRepoFacts(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "docs/spec_test.go", "package docs\n")
	writeFile(t, root, "internal/app/app_test.go", "package app\n")

	report, err := Review(context.Background(), root, Options{Focus: "internal"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.TestFileCount != 1 {
		t.Fatalf("TestFileCount = %d, want focused test count 1", report.TestFileCount)
	}
	if len(report.DependencyFiles) != 0 {
		t.Fatalf("DependencyFiles = %#v, want none in focus", report.DependencyFiles)
	}
}

func TestValidateOptionsRejectsUnsupportedValues(t *testing.T) {
	if err := ValidateOptions(Options{Scope: "ai-readiness"}); err == nil {
		t.Fatal("ValidateOptions() with unsupported scope succeeded")
	}
	if err := ValidateOptions(Options{Format: "json"}); err == nil {
		t.Fatal("ValidateOptions() with unsupported format succeeded")
	}
}

func TestReviewHTMLErrorIsExplicit(t *testing.T) {
	_, err := Review(context.Background(), t.TempDir(), Options{Format: "html"})
	if err == nil || !strings.Contains(err.Error(), "html review output is not implemented yet") {
		t.Fatalf("Review(html) error = %v", err)
	}
}

func TestRenderMarkdownDefaultsToFindingsOnly(t *testing.T) {
	report := Report{
		RepoRoot:         "/repo",
		Scope:            "security",
		Format:           "markdown",
		Deep:             true,
		Since:            "30d",
		Focus:            "internal",
		BaselineScopes:   []string{"dependencies", "testing", "maintainability"},
		DependencyFiles:  []string{"go.mod"},
		TestFileCount:    2,
		TrackedFileCount: 5,
		Docs: []FilePresence{
			{Path: "README.md", Present: true},
			{Path: "AGENTS.md", Present: false},
		},
		Findings: []Finding{{
			ID:             "testing.no-tests",
			Scopes:         []string{"testing", "maintainability"},
			Title:          "No test files detected",
			Summary:        "No test surface was detected.",
			Benefit:        "Improves regression safety.",
			Evidence:       []Evidence{{Label: "Test files", Value: "0"}},
			Recommendation: "Add tests.",
			Strength:       "Strong",
			SourceIDs:      []string{"fowler-test-pyramid"},
		}},
	}

	text := RenderMarkdown(report)
	for _, want := range []string{
		"## Recommendations",
		"### 1. No test files detected",
		"**Why:** No test surface was detected.",
		"**Benefit:** Improves regression safety.",
		"**Do next:** Add tests.",
		"**Attribution:** `fowler-test-pyramid`",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("RenderMarkdown() missing %q in:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{"# GX Review", "Scope:", "Depth:", "Focus:", "## Repo Facts", "## Changed Files", "Dependency manifests", "Since: `30d`", "Strength:", "Test files: 0", "## Sources", "## Context Sources"} {
		if strings.Contains(text, unwanted) {
			t.Fatalf("RenderMarkdown() should not include %q by default:\n%s", unwanted, text)
		}
	}
}

func TestRenderMarkdownIncludesAttributionSections(t *testing.T) {
	report := Report{
		Findings: []Finding{{
			ID:             "architecture.generic-package-name",
			Title:          "Generic package names reduce Interface clarity",
			Summary:        "Generic package names make it harder to infer what Module Interface callers should rely on.",
			Benefit:        "Improves readability and navigation by making package names describe domain concepts instead of storage buckets for shared code.",
			Recommendation: "Rename or split the package.",
			Strength:       "Worth exploring",
			SourceIDs:      []string{"go-code-review-comments", "custom-source"},
		}},
		Sources: []Source{
			{ID: "go-code-review-comments", Title: "Go Code Review Comments", URL: "https://go.dev/wiki/CodeReviewComments"},
			{ID: "custom-source", Title: "Custom Reference"},
		},
		SourceRefs: []SourceRef{
			{ID: "R1", Kind: "indexed_code", Title: "app.go", URL: "https://example.com/snippet", Source: "turbopuffer:gx"},
			{ID: "L1", Kind: "local", File: "CONTEXT.md", StartLine: 12, Source: "local"},
		},
	}

	text := RenderMarkdown(report)
	for _, want := range []string{
		"**Attribution:** [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) · Custom Reference (`custom-source`)",
		"## Sources",
		"- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)",
		"- Custom Reference (`custom-source`)",
		"## Context Sources",
		"`R1` [app.go](https://example.com/snippet) · source=turbopuffer:gx",
		"`L1` `CONTEXT.md:12`",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("RenderMarkdown() missing %q in:\n%s", want, text)
		}
	}
}

func TestRenderMarkdownVerboseIncludesFacts(t *testing.T) {
	report := Report{
		RepoRoot:         "/repo",
		Scope:            "security",
		Format:           "markdown",
		Since:            "30d",
		BaselineScopes:   []string{"dependencies", "testing", "maintainability"},
		DependencyFiles:  []string{"go.mod"},
		TestFileCount:    2,
		TrackedFileCount: 5,
		Docs:             []FilePresence{{Path: "README.md", Present: true}},
		ChangedFiles:     []string{"main.go"},
		Verbose:          true,
	}
	text := RenderMarkdown(report)
	for _, want := range []string{"## Repo Facts", "Since: `30d`", "Dependency manifests: `go.mod`", "## Changed Files", "`main.go`"} {
		if !strings.Contains(text, want) {
			t.Fatalf("RenderMarkdown(verbose) missing %q in:\n%s", want, text)
		}
	}
}

func TestSourcesAreCollectedAndRendered(t *testing.T) {
	root := t.TempDir()
	report, err := Review(context.Background(), root, Options{Scope: "security"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Sources) == 0 {
		t.Fatal("Review() did not collect internal sources")
	}
	text := RenderMarkdown(report)
	if !strings.Contains(text, "## Sources") {
		t.Fatalf("RenderMarkdown() missing sources section:\n%s", text)
	}
}

func TestStaticToolFailuresRankAboveArchitectureHints(t *testing.T) {
	ctx := ReviewContext{
		Options: Options{Scope: "architecture"},
		Facts: RepoFacts{
			DependencyFiles: []string{"go.mod"},
			GoPackages:      []PackageFact{{Path: "shared", GoFiles: 8, TestFiles: 0}},
		},
		Sources: []Source{
			{ID: "google-eng-practices", Title: "Google Engineering Practices", URL: "https://google.github.io/eng-practices/"},
			{ID: "fowler-architecture", Title: "Martin Fowler Architecture Guide", URL: "https://www.martinfowler.com/architecture/"},
		},
		Brief: ReviewBrief{
			Static: StaticSnapshot{
				ToolResults: []StaticToolResult{{
					Name:     "go test",
					Command:  "go test ./...",
					ExitCode: 1,
					Output:   "FAIL",
				}},
			},
		},
		ActiveScopes: []string{"architecture", "maintainability", "testing", "dependencies", "security"},
	}
	findings := evaluateFindings(ctx, defaultRules())
	if len(findings) == 0 {
		t.Fatal("evaluateFindings() returned no findings")
	}
	if findings[0].ID != "tools.static-failure" {
		t.Fatalf("first finding = %q, want tools.static-failure", findings[0].ID)
	}
}

func TestEngineUsesInjectedScannerCatalogAndRules(t *testing.T) {
	engine := NewEngineWith(
		fakeScanner{facts: RepoFacts{
			Docs:             []FilePresence{{Path: "README.md", Present: true}},
			DependencyFiles:  []string{"go.mod"},
			TestFileCount:    1,
			TrackedFileCount: 2,
		}},
		fakeCatalog{sources: []Source{{ID: "custom-source", Scopes: []string{"architecture"}}}},
		[]Rule{fakeRule{}},
	)

	report, err := engine.Review(context.Background(), "/repo", Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.TrackedFileCount != 2 || len(report.Sources) != 1 || report.Sources[0].ID != "custom-source" {
		t.Fatalf("report = %#v, injected scanner/catalog not used", report)
	}
	if !hasFinding(report.Findings, "custom.rule") {
		t.Fatalf("Findings = %#v, injected rule not used", report.Findings)
	}
}

func TestEngineUsesAIReviewerWhenAvailable(t *testing.T) {
	engine := NewEngineWithReviewer(
		fakeScanner{facts: RepoFacts{
			Docs:             []FilePresence{{Path: "README.md", Present: true}},
			DependencyFiles:  []string{"go.mod"},
			TestFileCount:    1,
			TrackedFileCount: 2,
			GoPackages:       []PackageFact{{Path: "internal/app", GoFiles: 6, TestFiles: 1}},
		}},
		fakeCatalog{sources: []Source{{ID: "custom-source", Scopes: []string{"architecture"}}}},
		[]Rule{fakeRule{}},
		fakeRetriever{snippets: []ContextSnippet{{Kind: "module_file", Ref: "internal/app/app.go", Text: "package app"}}},
		fakeReviewer{findings: []Finding{{
			ID:             "ai.review.1",
			Scopes:         []string{"architecture"},
			Title:          "Deepen the app Module",
			Summary:        "`internal/app` leaks workflow ordering across its Interface.",
			Recommendation: "Move the ordering behind one package seam and test through that Interface.",
			Strength:       "Worth exploring",
		}}},
	)

	report, err := engine.Review(context.Background(), "/repo", Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Findings) != 2 {
		t.Fatalf("Findings = %#v, want merged AI + rule findings", report.Findings)
	}
	if !hasFinding(report.Findings, "ai.review.1") || !hasFinding(report.Findings, "custom.rule") {
		t.Fatalf("Findings = %#v, want both ai.review.1 and custom.rule", report.Findings)
	}
}

func TestEnginePassesReviewPromptToAIReviewer(t *testing.T) {
	reviewer := &capturingReviewer{}
	engine := NewEngineWithReviewer(
		fakeScanner{facts: RepoFacts{
			Docs:             []FilePresence{{Path: "README.md", Present: true}},
			DependencyFiles:  []string{"go.mod"},
			TestFileCount:    1,
			TrackedFileCount: 2,
		}},
		fakeCatalog{sources: []Source{{ID: "custom-source", Scopes: []string{"architecture"}}}},
		[]Rule{fakeRule{}},
		fakeRetriever{},
		reviewer,
	)

	_, err := engine.Review(context.Background(), "/repo", Options{Prompt: "review auth rollback risk"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if reviewer.brief.ReviewPrompt != "review auth rollback risk" {
		t.Fatalf("AI brief ReviewPrompt = %q", reviewer.brief.ReviewPrompt)
	}
	if reviewer.brief.ReviewProfile != "prompt_directed" {
		t.Fatalf("AI brief ReviewProfile = %q, want prompt_directed", reviewer.brief.ReviewProfile)
	}
}

func TestOpenAIReviewerFromEnvPrefersUserOpenAIKey(t *testing.T) {
	t.Setenv("GX_OPENAI_PROXY_URL", "")
	t.Setenv("GX_CLOUD_URL", "off")
	t.Setenv("OPENAI_API_KEY", "openai-key")
	t.Setenv("OPENAI_BASE_URL", "http://127.0.0.1:43123")
	t.Setenv("GX_OPENAI_API_KEY", "gx-key")
	t.Setenv("GX_OPENAI_BASE_URL", "http://127.0.0.1:43124")

	reviewer := openAIReviewerFromEnv()
	got, ok := reviewer.(*responsesAIReviewer)
	if !ok {
		t.Fatalf("reviewer = %T, want *responsesAIReviewer", reviewer)
	}
	if got.url != "http://127.0.0.1:43123/v1/responses" || got.token != "openai-key" {
		t.Fatalf("reviewer config = url %q token %q", got.url, got.token)
	}
}

func TestStaticToolFailureBeatsSpeculativeFindings(t *testing.T) {
	findings := evaluateFindings(ReviewContext{
		ActiveScopes: []string{"architecture", "testing", "maintainability", "dependencies"},
		Sources:      []Source{{ID: "google-eng-practices", Scopes: []string{"maintainability", "testing"}}},
		Brief: ReviewBrief{Static: StaticSnapshot{ToolResults: []StaticToolResult{{
			Name:     "go vet",
			Command:  "go vet ./...",
			ExitCode: 1,
			Output:   "internal/app/app.go:10: unreachable code",
		}}}},
	}, defaultRules())
	if len(findings) == 0 || findings[0].ID != "tools.static-failure" {
		t.Fatalf("Findings = %#v, want static tool failure first", findings)
	}
}

func TestPatchFocusedReviewFiltersGenericArchitectureFindings(t *testing.T) {
	findings := filterPatchFocusedFindings(ReviewContext{
		Options: Options{PatchFocused: true},
		Brief:   ReviewBrief{Static: StaticSnapshot{ChangedFiles: []string{"internal/auth/session.go"}}},
	}, []Finding{
		{
			ID:             "architecture.implementation-heavy-module",
			Scopes:         []string{"architecture", "maintainability"},
			Title:          "Deepen the auth Module",
			Summary:        "`internal/auth` has many files.",
			Benefit:        "Better locality.",
			Recommendation: "Refactor the auth Module.",
		},
		{
			ID:             "ai.review.1",
			Scopes:         []string{"security", "testing"},
			Title:          "Preserve webhook signature verification",
			Summary:        "The changed code in `internal/auth/session.go` can accept an unsigned callback.",
			Benefit:        "Prevents accepting forged auth callbacks.",
			Recommendation: "Update `internal/auth/session.go` to reject missing signatures and add a failing callback test.",
		},
	})
	if len(findings) != 1 || findings[0].Title != "Preserve webhook signature verification" {
		t.Fatalf("Findings = %#v, want only changed-file security finding", findings)
	}
}

func TestPatchFocusedReviewKeepsChangedFileQualityBug(t *testing.T) {
	root := initRepo(t)
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	gitAdd(t, root, "go.mod", "internal/app/app.go")
	gitCommit(t, root)
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() { result, _ := maybe(); _ = result }\nfunc maybe() (int, error) { return 0, nil }\n")

	report, err := Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if !hasFinding(report.Findings, "quality.ignored-results") {
		t.Fatalf("Findings = %#v, want changed-file ignored result finding", report.Findings)
	}
}

func TestBuildReviewBriefUsesArchitectureRubricAndContext(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "README.md", "# repo\n")
	writeFile(t, root, "CONTEXT.md", "# Context\n\n## Core Terms\n\n### Revision\n\nOne reviewable change.\n")
	writeFile(t, root, "docs/adr/0001-revisions.md", "# ADR 0001\n\nUse Revision as product identity.\n")
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	facts := RepoFacts{
		Docs:             []FilePresence{{Path: "README.md", Present: true}, {Path: "CONTEXT.md", Present: true}},
		ADRFiles:         []string{"docs/adr/0001-revisions.md"},
		DependencyFiles:  []string{"go.mod"},
		Files:            []string{"README.md", "CONTEXT.md", "docs/adr/0001-revisions.md", "go.mod", "internal/app/app.go"},
		GoPackages:       []PackageFact{{Path: "internal/app", GoFiles: 6, TestFiles: 0}},
		TestFileCount:    0,
		TrackedFileCount: 5,
	}

	brief, err := BuildReviewBrief(context.Background(), root, normalizeOptions(Options{}), facts, []Source{{ID: "go-package-names", Title: "hidden", URL: "https://example.com", Scopes: []string{"architecture"}}}, LocalContextRetriever{})
	if err != nil {
		t.Fatalf("BuildReviewBrief() error = %v", err)
	}
	if len(brief.Hints) == 0 || brief.Hints[0].Kind != "deepening_candidate" {
		t.Fatalf("Hints = %#v, want deepening candidate", brief.Hints)
	}
	if brief.ReviewProfile != "patch_focused" {
		t.Fatalf("ReviewProfile = %q, want patch_focused", brief.ReviewProfile)
	}
	if !strings.Contains(strings.Join(brief.Rubric.Questions, "\n"), "static.diff_snippets") {
		t.Fatalf("Rubric = %#v, want patch-focused diff question", brief.Rubric)
	}
	if len(brief.Context) == 0 {
		t.Fatalf("Context = %#v, want local snippets", brief.Context)
	}
	if !hasContextSnippet(brief.Context, "domain_doc", "CONTEXT.md") {
		t.Fatalf("Context = %#v, want CONTEXT.md domain snippet", brief.Context)
	}
	if !hasContextSnippet(brief.Context, "adr", "docs/adr/0001-revisions.md") {
		t.Fatalf("Context = %#v, want ADR snippet", brief.Context)
	}
	if len(brief.SourceCatalog) != 1 || brief.SourceCatalog[0].ID != "go-package-names" {
		t.Fatalf("SourceCatalog = %#v", brief.SourceCatalog)
	}
}

func TestPromptedReviewBriefUsesPromptDirectedProfile(t *testing.T) {
	opts := normalizeOptions(Options{Prompt: "review auth rollback risk"})
	if opts.PatchFocused {
		t.Fatalf("PatchFocused = true, want false for prompted review")
	}
	if profile := reviewProfile(opts); profile != "prompt_directed" {
		t.Fatalf("reviewProfile() = %q, want prompt_directed", profile)
	}
	if got := strings.Join(activeScopeList(opts), ","); got != "architecture,dependencies,testing,maintainability" {
		t.Fatalf("activeScopeList() = %q", got)
	}

	brief, err := BuildReviewBrief(context.Background(), t.TempDir(), opts, RepoFacts{}, nil, fakeRetriever{})
	if err != nil {
		t.Fatalf("BuildReviewBrief() error = %v", err)
	}
	if brief.ReviewPrompt != "review auth rollback risk" {
		t.Fatalf("ReviewPrompt = %q", brief.ReviewPrompt)
	}
	if brief.ReviewProfile != "prompt_directed" {
		t.Fatalf("ReviewProfile = %q, want prompt_directed", brief.ReviewProfile)
	}
	if !strings.Contains(brief.Rubric.Goal, "review_prompt") {
		t.Fatalf("Rubric goal = %q, want review_prompt guidance", brief.Rubric.Goal)
	}
}

func TestBuildReviewBriefIncludesDiffSnippetsForCurrentPatch(t *testing.T) {
	root := initRepo(t)
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	gitAdd(t, root, "go.mod", "internal/app/app.go")
	gitCommit(t, root)
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\nfunc NewBehavior() {}\n")
	facts := RepoFacts{
		Docs:             []FilePresence{{Path: "README.md", Present: false}},
		DependencyFiles:  []string{"go.mod"},
		Files:            []string{"go.mod", "internal/app/app.go"},
		GoPackages:       []PackageFact{{Path: "internal/app", GoFiles: 1, TestFiles: 0}},
		TrackedFileCount: 2,
	}

	brief, err := BuildReviewBrief(context.Background(), root, normalizeOptions(Options{}), facts, nil, LocalContextRetriever{})
	if err != nil {
		t.Fatalf("BuildReviewBrief() error = %v", err)
	}
	if len(brief.Static.DiffSnippets) != 1 {
		t.Fatalf("DiffSnippets = %#v, want one changed file diff", brief.Static.DiffSnippets)
	}
	if !strings.Contains(brief.Static.DiffSnippets[0].Diff, "+func NewBehavior()") {
		t.Fatalf("DiffSnippets[0] = %#v, want added function diff", brief.Static.DiffSnippets[0])
	}
}

func TestBuildReviewBriefLabelsRetrievedContext(t *testing.T) {
	brief, err := BuildReviewBrief(
		context.Background(),
		t.TempDir(),
		normalizeOptions(Options{}),
		RepoFacts{},
		nil,
		fakeRetriever{snippets: []ContextSnippet{
			{Kind: "review_resource", Ref: "owasp#1", Source: "turbopuffer:gx-review-knowledge", Text: "Parameterized queries prevent injection."},
			{Kind: "repo_doc", Ref: "REVIEW.md", Source: "local", Text: "Review auth changes carefully."},
		}},
	)
	if err != nil {
		t.Fatalf("BuildReviewBrief() error = %v", err)
	}
	if len(brief.Context) != 2 {
		t.Fatalf("Context = %#v, want two snippets", brief.Context)
	}
	if brief.Context[0].SourceLabel != "R1" || !strings.Contains(brief.Context[0].Text, "[R1] kind=review_resource") {
		t.Fatalf("review resource label = %#v", brief.Context[0])
	}
	if brief.Context[1].SourceLabel != "L1" || !strings.Contains(brief.Context[1].Text, "[L1] kind=repo_doc") {
		t.Fatalf("local label = %#v", brief.Context[1])
	}
}

func TestDeepReviewBriefUsesFullSpectrumRubric(t *testing.T) {
	root := t.TempDir()
	facts := RepoFacts{Docs: []FilePresence{{Path: "README.md", Present: false}}}

	brief, err := BuildReviewBrief(context.Background(), root, normalizeOptions(Options{Deep: true}), facts, nil, LocalContextRetriever{})
	if err != nil {
		t.Fatalf("BuildReviewBrief(deep) error = %v", err)
	}
	if brief.ReviewProfile != "deep_full_spectrum" {
		t.Fatalf("ReviewProfile = %q, want deep_full_spectrum", brief.ReviewProfile)
	}
	questions := strings.Join(brief.Rubric.Questions, "\n")
	for _, want := range []string{"authn/authz", "idempotency", "observability", "performance", "deletion test"} {
		if !strings.Contains(questions, want) {
			t.Fatalf("Deep rubric missing %q in %#v", want, brief.Rubric.Questions)
		}
	}
}

func TestAIReviewPromptSeparatesPatchAndDeepReview(t *testing.T) {
	prompt := reviewDeveloperPrompt()
	for _, want := range []string{
		"patch_focused",
		"prompt_directed",
		"deep_full_spectrum",
		"review_prompt",
		"static.diff_snippets",
		"security/auth",
		"race/idempotency",
		"observability",
		"broad architecture",
		"empty recommendations array",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("review prompt missing %q in:\n%s", want, prompt)
		}
	}
}

func TestBuildReviewBriefIncludesCodeQualityHints(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() { result, _ := maybe(); _ = result }\nfunc maybe() (int, error) { return 0, nil }\n")
	facts := RepoFacts{
		Docs:             []FilePresence{{Path: "README.md", Present: false}, {Path: "CONTEXT.md", Present: false}},
		DependencyFiles:  []string{"go.mod"},
		Files:            []string{"go.mod", "internal/app/app.go"},
		GoPackages:       []PackageFact{{Path: "internal/app", GoFiles: 1, TestFiles: 0}},
		TestFileCount:    0,
		TrackedFileCount: 2,
	}

	brief, err := BuildReviewBrief(context.Background(), root, normalizeOptions(Options{}), facts, nil, LocalContextRetriever{})
	if err != nil {
		t.Fatalf("BuildReviewBrief() error = %v", err)
	}
	if len(brief.Static.CodeQuality) == 0 || brief.Static.CodeQuality[0].Kind != "ignored_result" {
		t.Fatalf("CodeQuality = %#v, want ignored result hint", brief.Static.CodeQuality)
	}
	foundSnippet := false
	for _, snippet := range brief.Context {
		if snippet.Kind == "code_quality_file" && snippet.Ref == "internal/app/app.go" {
			foundSnippet = true
		}
	}
	if !foundSnippet {
		t.Fatalf("Context = %#v, want code quality file snippet", brief.Context)
	}
}

func TestIgnoredResultFindingUsesQualityHints(t *testing.T) {
	findings := evaluateFindings(ReviewContext{
		ActiveScopes: []string{"architecture", "testing", "maintainability"},
		Sources:      []Source{{ID: "google-eng-practices", Scopes: []string{"maintainability", "testing"}}},
		Brief: ReviewBrief{Static: StaticSnapshot{CodeQuality: []CodeQualityHint{{
			Kind: "ignored_result",
			File: "internal/app/app.go",
			Line: 7,
			Text: "proposal, _ = shapeProposal(proposal)",
		}}}},
	}, defaultRules())
	if !hasFinding(findings, "quality.ignored-results") {
		t.Fatalf("Findings = %#v, want ignored result finding", findings)
	}
}

func TestDomainLanguageDriftFindingUsesContextInternalTerms(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "CONTEXT.md", strings.Join([]string{
		"# Context",
		"",
		"## Internal Implementation Terms",
		"",
		"### Demux",
		"",
		"Internal splitting algorithm.",
		"",
	}, "\n"))
	writeFile(t, root, "README.md", "Use demux to split work.\n")

	report, err := Review(context.Background(), root, Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if !hasFinding(report.Findings, "architecture.domain-language-drift") {
		t.Fatalf("Findings = %#v, want domain language drift finding", report.Findings)
	}
}

func TestEvaluateFindingsUsesScopeAndBaselines(t *testing.T) {
	root := t.TempDir()
	report, err := Review(context.Background(), root, Options{Scope: "security"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	for _, want := range []string{
		"testing.no-tests",
		"dependencies.no-manifest",
	} {
		if !hasFinding(report.Findings, want) {
			t.Fatalf("Findings = %#v, want %s from baseline scope", report.Findings, want)
		}
	}
	if hasFinding(report.Findings, "architecture.missing-context") {
		t.Fatalf("Findings = %#v, did not expect missing context finding", report.Findings)
	}
}

type fakeScanner struct {
	facts RepoFacts
}

func (f fakeScanner) Scan(context.Context, string, string) (RepoFacts, error) {
	return f.facts, nil
}

type fakeCatalog struct {
	sources []Source
}

func (f fakeCatalog) SourcesForScopes([]string) []Source {
	return f.sources
}

type fakeRule struct{}

func (fakeRule) ID() string { return "custom.rule" }

func (fakeRule) Scopes() []string { return []string{"architecture"} }

func (fakeRule) Evaluate(ReviewContext) []Finding {
	return []Finding{{
		ID:             "custom.rule",
		Scopes:         []string{"architecture"},
		Title:          "Custom rule",
		Summary:        "Custom summary.",
		Recommendation: "Custom recommendation.",
		Strength:       "Speculative",
		SourceIDs:      []string{"custom-source"},
	}}
}

type fakeRetriever struct {
	snippets []ContextSnippet
}

func (f fakeRetriever) Retrieve(context.Context, string, Options, RepoFacts, []ReviewHint) ([]ContextSnippet, error) {
	return f.snippets, nil
}

type fakeReviewer struct {
	findings []Finding
}

func (f fakeReviewer) Review(context.Context, ReviewBrief) ([]Finding, error) {
	return f.findings, nil
}

type capturingReviewer struct {
	brief ReviewBrief
}

func (c *capturingReviewer) Review(_ context.Context, brief ReviewBrief) ([]Finding, error) {
	c.brief = brief
	return []Finding{{
		ID:             "ai.prompt",
		Title:          "Prompt-directed finding",
		Summary:        "The prompted concern is represented in the review brief.",
		Benefit:        "Keeps prompted review behavior testable.",
		Recommendation: "Keep review_prompt in the AI review brief.",
	}}, nil
}

func TestJavaScriptLockfileFinding(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "package.json", "{}\n")

	report, err := Review(context.Background(), root, Options{Scope: "dependencies"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if !hasFinding(report.Findings, "dependencies.javascript-unlocked") {
		t.Fatalf("Findings = %#v, want javascript lockfile finding", report.Findings)
	}
}

func TestArchitecturePackageFindings(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "README.md", "# repo\n")
	writeFile(t, root, "internal/util/format.go", "package util\n")
	writeFile(t, root, "internal/authoring/engine.go", "package authoring\n")
	writeFile(t, root, "internal/authoring/engine_test.go", "package authoring\n")
	writeFile(t, root, "internal/authoring/proposal.go", "package authoring\n")
	writeFile(t, root, "internal/orchestrator/a.go", "package orchestrator\n")
	writeFile(t, root, "internal/orchestrator/b.go", "package orchestrator\n")
	writeFile(t, root, "internal/orchestrator/c.go", "package orchestrator\n")
	writeFile(t, root, "internal/orchestrator/d.go", "package orchestrator\n")
	writeFile(t, root, "internal/orchestrator/e.go", "package orchestrator\n")
	writeFile(t, root, "internal/orchestrator/f.go", "package orchestrator\n")

	report, err := Review(context.Background(), root, Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	for _, want := range []string{
		"architecture.generic-package-name",
		"architecture.implementation-heavy-module",
		"architecture.package-without-tests",
	} {
		if !hasFinding(report.Findings, want) {
			t.Fatalf("Findings = %#v, want %s", report.Findings, want)
		}
	}
}

func TestDocsWithoutADRsDoesNotCreateFinding(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/overview.md", "# overview\n")

	report, err := Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if hasFinding(report.Findings, "architecture.docs-without-adrs") {
		t.Fatalf("Findings = %#v, did not expect ADR finding", report.Findings)
	}
}

func TestChangedFilesPreservesFirstCharacter(t *testing.T) {
	root := initRepo(t)
	writeFile(t, root, "README.md", "# old\n")
	gitAdd(t, root, "README.md")
	gitCommit(t, root)
	writeFile(t, root, "README.md", "# new\n")

	files := changedFiles(context.Background(), root)
	if len(files) != 1 || files[0] != "README.md" {
		t.Fatalf("changedFiles() = %#v, want README.md", files)
	}
}

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func initRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGit(t, root, "init")
	runGit(t, root, "config", "user.name", "Test User")
	runGit(t, root, "config", "user.email", "test@example.com")
	return root
}

func gitAdd(t *testing.T, root string, files ...string) {
	t.Helper()
	runGit(t, root, append([]string{"add"}, files...)...)
}

func gitCommit(t *testing.T, root string) {
	t.Helper()
	runGit(t, root, "commit", "-m", "initial")
}

func runGit(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s error = %v\n%s", strings.Join(args, " "), err, out)
	}
}

func hasFinding(findings []Finding, id string) bool {
	for _, finding := range findings {
		if finding.ID == id {
			return true
		}
	}
	return false
}

func hasContextSnippet(snippets []ContextSnippet, kind, ref string) bool {
	for _, snippet := range snippets {
		if snippet.Kind == kind && snippet.Ref == ref {
			return true
		}
	}
	return false
}
