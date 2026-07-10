package authoring

import (
	"errors"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/vcs"
)

func TestGenerateLogicConfidencePenalizesRiskyRouteAndDependencyPlans(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{{ID: "h1", File: "alpha.go"}},
		Revisions: []RevisionProposal{{
			ID:               "u1",
			Intent:           "alpha",
			Files:            []string{"alpha.go"},
			UseHunks:         true,
			HunkIDs:          []string{"h1"},
			TargetStack:      "feature/alpha",
			RouteConfidence:  0.49,
			ProvenanceStatus: "absent",
			Confidence:       0.5,
		}},
		FeasibilityWarnings: []FeasibilityWarning{{
			RevisionID: "u1",
			Severity:   "warning",
			Source:     "structural_dependency",
			DependsOn:  "u0",
			Message:    "revision u1 depends on u0",
		}},
	}

	got := annotateGenerateLogicConfidence(proposal)

	if got.Confidence.LogicConfidence >= generateHighLogicConfidence {
		t.Fatalf("logic confidence = %.2f, want below high gate", got.Confidence.LogicConfidence)
	}
	if got.Confidence.EffectiveConfidence != got.Confidence.LogicConfidence {
		t.Fatalf("effective confidence = %.2f, want logic %.2f", got.Confidence.EffectiveConfidence, got.Confidence.LogicConfidence)
	}
	if !hasConfidenceReason(got.Confidence.LogicReasons, "low_confidence_route") || !hasConfidenceReason(got.Confidence.LogicReasons, "structural_dependency") {
		t.Fatalf("logic reasons = %#v, want route and dependency risks", got.Confidence.LogicReasons)
	}
	for _, reason := range got.Confidence.LogicReasons {
		if reason.Severity == "positive" {
			continue
		}
		if strings.TrimSpace(reason.Suggestion) == "" {
			t.Fatalf("logic reason missing suggestion: %#v", reason)
		}
	}
}

func TestGenerateConfidenceLogIncludesActionAndSuggestions(t *testing.T) {
	proposal := annotateGenerateLogicConfidence(DemuxProposal{
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "first", Files: []string{"alpha.go"}, UseHunks: true, HunkIDs: []string{"h1"}},
			{ID: "u2", Intent: "second", Files: []string{"alpha.go"}, UseHunks: true, HunkIDs: []string{"h2"}},
		},
	})

	got := appendGenerateConfidenceLog(proposal, "skip_llm_repair_medium_confidence", DemuxWorkflowReadyToApply)

	if len(got.Warnings) == 0 {
		t.Fatal("warnings empty, want confidence log")
	}
	warning := got.Warnings[len(got.Warnings)-1]
	for _, want := range []string{
		"generate confidence logic=",
		"state=ready_to_apply",
		"action=skip_llm_repair_medium_confidence",
		"merge shared-file edits",
	} {
		if !strings.Contains(warning, want) {
			t.Fatalf("confidence warning missing %q:\n%s", want, warning)
		}
	}
}

func TestGenerateLLMRepairDefaultsOff(t *testing.T) {
	t.Setenv("GX_GENERATE_LLM_REPAIR", "")
	if generateLLMRepairEnabled() {
		t.Fatal("generate LLM repair default = true, want false")
	}
	t.Setenv("GX_GENERATE_LLM_REPAIR", "1")
	if !generateLLMRepairEnabled() {
		t.Fatal("GX_GENERATE_LLM_REPAIR=1 did not enable LLM repair")
	}
}

func TestConservativeGenerateProposalPreservesValidHunkLevelPlan(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "alpha.go", Patch: "@@\n+const first = true\n"},
			{ID: "h2", File: "alpha.go", Patch: "@@\n+const second = true\n"},
		},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "first", Files: []string{"alpha.go"}, UseHunks: true, HunkIDs: []string{"h1"}, ProvenanceStatus: "explicit", SessionIDs: []string{"s1"}},
			{ID: "u2", Intent: "second", Files: []string{"alpha.go"}, UseHunks: true, HunkIDs: []string{"h2"}, ProvenanceStatus: "explicit", SessionIDs: []string{"s1"}},
		},
	}

	got := conservativeGenerateProposal(proposal, "ship alpha", "test fallback")

	if len(got.Revisions) != 2 {
		t.Fatalf("fallback revisions = %#v, want precise hunk-level revisions preserved", got.Revisions)
	}
	for _, revision := range got.Revisions {
		if !revision.UseHunks || len(revision.HunkIDs) != 1 || len(revision.Hunks) != 1 {
			t.Fatalf("fallback revision = %#v, want hunk-level coverage", revision)
		}
	}
	if got.Revisions[0].ProvenanceStatus != "explicit" || len(got.Revisions[0].SessionIDs) != 1 || got.Revisions[0].SessionIDs[0] != "s1" {
		t.Fatalf("fallback provenance = %q/%#v, want explicit session", got.Revisions[0].ProvenanceStatus, got.Revisions[0].SessionIDs)
	}
	if _, err := normalizeDemuxProposal(got); err != nil {
		t.Fatalf("normalizeDemuxProposal(fallback) error = %v", err)
	}
}

func TestConservativeGenerateProposalRepairsDuplicateHunkByFlatteningImplicatedRevisions(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "alpha.go", Patch: "@@\n+const alpha = true\n"},
			{ID: "h2", File: "beta.go", Patch: "@@\n+const beta = true\n"},
		},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "alpha first", Files: []string{"alpha.go"}, UseHunks: true, HunkIDs: []string{"h1"}},
			{ID: "u2", Intent: "alpha duplicate", Files: []string{"alpha.go"}, UseHunks: true, HunkIDs: []string{"h1"}},
			{ID: "u3", Intent: "beta", Files: []string{"beta.go"}, UseHunks: true, HunkIDs: []string{"h2"}},
		},
	}

	got := conservativeGenerateProposal(proposal, "ship alpha", "test fallback")

	if len(got.Revisions) != 2 {
		t.Fatalf("fallback revisions = %#v, want duplicate owners merged only", got.Revisions)
	}
	if got.Revisions[0].ID != "u1" || got.Revisions[0].UseHunks || len(got.Revisions[0].Files) != 1 || got.Revisions[0].Files[0] != "alpha.go" {
		t.Fatalf("merged duplicate revision = %#v, want whole-file alpha u1", got.Revisions[0])
	}
	if got.Revisions[1].ID != "u3" || !got.Revisions[1].UseHunks || len(got.Revisions[1].HunkIDs) != 1 || got.Revisions[1].HunkIDs[0] != "h2" {
		t.Fatalf("unimplicated revision = %#v, want beta hunk preserved", got.Revisions[1])
	}
	if _, err := normalizeDemuxProposal(got); err != nil {
		t.Fatalf("normalizeDemuxProposal(fallback) error = %v", err)
	}
}

func TestConservativeGenerateProposalRepairsDuplicateWholeFileOwnerByFlatteningImplicatedRevisions(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "alpha.go", Patch: "@@\n+const alpha = true\n"},
			{ID: "h2", File: "beta.go", Patch: "@@\n+const beta = true\n"},
		},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "alpha first", Files: []string{"alpha.go"}},
			{ID: "u2", Intent: "alpha duplicate", Files: []string{"alpha.go"}},
			{ID: "u3", Intent: "beta", Files: []string{"beta.go"}, UseHunks: true, HunkIDs: []string{"h2"}},
		},
	}

	got := conservativeGenerateProposal(proposal, "ship alpha", "test fallback")

	if len(got.Revisions) != 2 || got.Revisions[0].ID != "u1" || got.Revisions[0].UseHunks || got.Revisions[1].ID != "u3" || !got.Revisions[1].UseHunks {
		t.Fatalf("fallback revisions = %#v, want only duplicate whole-file owners merged", got.Revisions)
	}
	if _, err := normalizeDemuxProposal(got); err != nil {
		t.Fatalf("normalizeDemuxProposal(fallback) error = %v", err)
	}
}

func TestConservativeGenerateProposalFallsBackToWholePlanForUnrepairableValidation(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "alpha.go", Patch: "@@\n+const first = true\n"},
			{ID: "h2", File: "alpha.go", Patch: "@@\n+const second = true\n"},
		},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "alpha first", Files: []string{"alpha.go"}, UseHunks: true, HunkIDs: []string{"h1"}},
		},
	}

	got := conservativeGenerateProposal(proposal, "ship alpha", "test fallback")

	if len(got.Revisions) != 1 || got.Revisions[0].UseHunks || len(got.Revisions[0].Files) != 1 || got.Revisions[0].Files[0] != "alpha.go" {
		t.Fatalf("fallback revisions = %#v, want whole-plan conservative fallback", got.Revisions)
	}
	if _, err := normalizeDemuxProposal(got); err != nil {
		t.Fatalf("normalizeDemuxProposal(fallback) error = %v", err)
	}
}

func TestDemuxValidationRepairUnknownErrorFallsBackToWholePlan(t *testing.T) {
	if repair := demuxValidationRepair(errors.New("legacy validation failure")); repair != nil {
		t.Fatalf("demuxValidationRepair() = %#v, want nil", repair)
	}
}

func TestCoverageCheckedGenerateProposalDowngradesMismatchWithWarning(t *testing.T) {
	sourceHunk := HunkRange{ID: "h1", File: "alpha.go", Patch: "diff --git a/alpha.go b/alpha.go\n@@\n+source\n"}
	coverageInput := DemuxProposal{
		Hunks: []HunkRange{sourceHunk},
		Revisions: []RevisionProposal{{
			ID:       "u1",
			Intent:   "alpha",
			UseHunks: true,
			HunkIDs:  []string{"h1"},
			Hunks:    []HunkRange{{ID: "h1", File: "alpha.go", Patch: "diff --git a/alpha.go b/alpha.go\n@@\n+altered\n"}},
		}},
	}
	readyProposal := DemuxProposal{
		Hunks: []HunkRange{sourceHunk},
		Revisions: []RevisionProposal{{
			ID:       "u1",
			Intent:   "alpha",
			UseHunks: true,
			HunkIDs:  []string{"h1"},
			Hunks:    []HunkRange{sourceHunk},
		}},
	}

	got, downgraded := coverageCheckedGenerateProposal(coverageInput, readyProposal, "ship alpha")

	if !downgraded {
		t.Fatal("coverageCheckedGenerateProposal() downgraded = false, want true")
	}
	if !hasGenerateWarning(got.Warnings, "patch coverage mismatch:") {
		t.Fatalf("warnings = %#v, want patch coverage mismatch warning", got.Warnings)
	}
	if err := validateDemuxPatchCoverage(got); err != nil {
		t.Fatalf("validateDemuxPatchCoverage(fallback) error = %v", err)
	}
	if got.Revisions[0].Hunks[0].Patch != sourceHunk.Patch {
		t.Fatalf("fallback hunk patch = %q, want source patch", got.Revisions[0].Hunks[0].Patch)
	}
}

func TestCoverageCheckedGenerateProposalKeepsValidPrecisePlan(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "alpha.go", Patch: "patch one"},
			{ID: "h2", File: "alpha.go", Patch: "patch two"},
		},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "alpha first", UseHunks: true, HunkIDs: []string{"h1"}, Hunks: []HunkRange{{ID: "h1", File: "alpha.go", Patch: "patch one"}}},
			{ID: "u2", Intent: "alpha second", UseHunks: true, HunkIDs: []string{"h2"}, Hunks: []HunkRange{{ID: "h2", File: "alpha.go", Patch: "patch two"}}},
		},
	}

	got, downgraded := coverageCheckedGenerateProposal(proposal, proposal, "ship alpha")

	if downgraded {
		t.Fatalf("coverageCheckedGenerateProposal() downgraded = true, proposal = %#v", got)
	}
	if len(got.Revisions) != 2 || !got.Revisions[0].UseHunks || !got.Revisions[1].UseHunks {
		t.Fatalf("proposal revisions = %#v, want precise hunk-level plan", got.Revisions)
	}
}

func TestConservativeGenerateProposalPreservesSeparateFileRevisions(t *testing.T) {
	proposal := DemuxProposal{
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "update docs", Files: []string{"README.md", "docs/index.mdx"}, ProvenanceStatus: "absent", Confidence: 0.55},
			{ID: "u2", Intent: "update generate pipeline", Files: []string{"internal/authoring/generate_pipeline.go", "internal/authoring/generate_pipeline_test.go"}, ProvenanceStatus: "absent", Confidence: 0.55},
			{ID: "u3", Intent: "update CLI flag", Files: []string{"internal/cli/root.go", "internal/cli/root_test.go"}, ProvenanceStatus: "absent", Confidence: 0.55},
		},
	}

	got := conservativeGenerateProposal(proposal, "update README", "test fallback")

	if len(got.Revisions) != 3 {
		t.Fatalf("fallback revisions = %#v, want separate heuristic revisions preserved", got.Revisions)
	}
	for _, revision := range got.Revisions {
		if revision.UseHunks || len(revision.HunkIDs) != 0 || len(revision.Hunks) != 0 {
			t.Fatalf("fallback revision = %#v, want whole-file revision", revision)
		}
	}
	if got.Revisions[0].Intent != "update docs" || got.Revisions[1].Intent != "update generate pipeline" || got.Revisions[2].Intent != "update CLI flag" {
		t.Fatalf("fallback intents = %#v, want original grouping names", got.Revisions)
	}
}

func TestPolishGenerateProposalCoalescesDocsAndNamesRevisions(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "README.md"},
			{ID: "h2", File: "docs/index.mdx"},
		},
		Revisions: []RevisionProposal{
			{
				ID:              "u1",
				Intent:          "update README.md",
				Files:           []string{"README.md"},
				HunkIDs:         []string{"h1"},
				TargetStack:     "docs/documentation",
				RouteConfidence: 0.72,
				Confidence:      0.55,
				SessionIDs:      []string{"s1"},
			},
			{
				ID:              "u2",
				Intent:          "update docs/index.mdx",
				Files:           []string{"docs/index.mdx"},
				HunkIDs:         []string{"h2"},
				TargetStack:     "docs/documentation",
				RouteConfidence: 0.66,
				Confidence:      0.70,
				SessionIDs:      []string{"s1", "s2"},
			},
			{
				ID:     "u3",
				Intent: "update internal/cli/root",
				Files:  []string{"internal/cli/root.go", "internal/cli/root_test.go"},
			},
		},
	}

	got := polishGenerateProposal(proposal)

	if len(got.Revisions) != 2 {
		t.Fatalf("revisions = %#v, want docs coalesced into two revisions", got.Revisions)
	}
	docs := got.Revisions[0]
	if docs.Intent != "update documentation" {
		t.Fatalf("docs intent = %q, want deterministic docs name", docs.Intent)
	}
	if len(docs.Files) != 2 || docs.Files[0] != "README.md" || docs.Files[1] != "docs/index.mdx" {
		t.Fatalf("docs files = %#v, want README and docs index", docs.Files)
	}
	if docs.TargetStack != "docs/documentation" {
		t.Fatalf("docs target = %q, want deterministic docs route", docs.TargetStack)
	}
	if docs.RouteConfidence != 0.66 {
		t.Fatalf("route confidence = %.2f, want lowest merged route confidence", docs.RouteConfidence)
	}
	if len(docs.SessionIDs) != 2 || docs.SessionIDs[0] != "s1" || docs.SessionIDs[1] != "s2" {
		t.Fatalf("session ids = %#v, want de-duplicated sessions", docs.SessionIDs)
	}
	if got.Revisions[1].Intent != "update internal/cli/root" {
		t.Fatalf("cli intent = %q, want original intent preserved", got.Revisions[1].Intent)
	}
	if got.Revisions[1].TargetStack != "" {
		t.Fatalf("cli target = %q, want no deterministic route", got.Revisions[1].TargetStack)
	}
}

func TestComposeApplyPreflightDefaultsOff(t *testing.T) {
	t.Setenv("GX_GENERATE_VERIFY", "")
	t.Setenv("GX_COMPOSE_SKIP_APPLY_PREFLIGHT", "")
	t.Setenv("GX_GENERATE_APPLY_PREFLIGHT_ATTEMPTS", "")
	t.Setenv("GX_COMPOSE_APPLY_PREFLIGHT_ATTEMPTS", "")

	if defaultDemuxApplyPreflightAttempts != 0 {
		t.Fatalf("defaultDemuxApplyPreflightAttempts = %d, want 0", defaultDemuxApplyPreflightAttempts)
	}

	repoRoot := t.TempDir()
	engine := NewEngineWithVCS(vcs.NewServiceWithRunner(&demuxAlreadyAppliedFakeRunner{repoRoot: repoRoot}))
	pipeline := engine.demuxPipeline()
	proposal := DemuxProposal{
		ID:       "p1",
		RepoRoot: repoRoot,
		Revisions: []RevisionProposal{{
			ID:          "u1",
			Intent:      "test",
			Files:       []string{"a.txt"},
			TargetStack: "feature/test",
		}},
	}
	got, err := pipeline.preflightApplyReadyProposal(t.Context(), proposal, ProposeDemuxOptions{})
	if err != nil {
		t.Fatalf("preflightApplyReadyProposal: %v", err)
	}
	foundSkip := false
	for _, warning := range got.Warnings {
		if strings.Contains(warning, "skipped disposable apply preflight") {
			foundSkip = true
			break
		}
	}
	if !foundSkip {
		t.Fatalf("warnings = %#v, want skip preflight warning", got.Warnings)
	}
}

func TestGenerateApplyPreflightDefaultsOff(t *testing.T) {
	t.Setenv("GX_GENERATE_VERIFY", "")
	t.Setenv("GX_COMPOSE_SKIP_APPLY_PREFLIGHT", "")

	if generateApplyPreflightEnabled() {
		t.Fatal("generate preflight default = true, want fast path default false")
	}

	t.Setenv("GX_GENERATE_VERIFY", "1")
	if !generateApplyPreflightEnabled() {
		t.Fatal("GX_GENERATE_VERIFY=1 did not enable preflight")
	}

	t.Setenv("GX_COMPOSE_SKIP_APPLY_PREFLIGHT", "1")
	if generateApplyPreflightEnabled() {
		t.Fatal("GX_COMPOSE_SKIP_APPLY_PREFLIGHT=1 should override generate verify")
	}
}

func hasConfidenceReason(reasons []ConfidenceReason, kind string) bool {
	_, ok := confidenceReason(reasons, kind)
	return ok
}

func confidenceReason(reasons []ConfidenceReason, kind string) (ConfidenceReason, bool) {
	for _, reason := range reasons {
		if reason.Kind == kind {
			return reason, true
		}
	}
	return ConfidenceReason{}, false
}

func hasGenerateWarning(warnings []string, prefix string) bool {
	for _, warning := range warnings {
		if strings.HasPrefix(warning, prefix) {
			return true
		}
	}
	return false
}
