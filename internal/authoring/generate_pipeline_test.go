package authoring

import (
	"strings"
	"testing"
)

func TestGenerateLogicConfidencePenalizesRiskyHunkAndRoutePlans(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{{ID: "h1", File: "alpha.go"}},
		Revisions: []RevisionProposal{{
			ID:               "u1",
			Intent:           "alpha",
			Files:            []string{"alpha.go"},
			UseHunks:         true,
			HunkIDs:          []string{"h1"},
			TargetStack:      "feature/alpha",
			RouteConfidence:  0.5,
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
	if len(got.Confidence.LogicReasons) < 2 {
		t.Fatalf("logic reasons = %#v, want hunk and dependency risks", got.Confidence.LogicReasons)
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

func TestConservativeGenerateProposalProducesHighConfidenceWholeFilePlan(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "alpha.go"},
			{ID: "h2", File: "alpha.go"},
		},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "first", Files: []string{"alpha.go"}, UseHunks: true, HunkIDs: []string{"h1"}, ProvenanceStatus: "explicit", SessionIDs: []string{"s1"}},
			{ID: "u2", Intent: "second", Files: []string{"alpha.go"}, UseHunks: true, HunkIDs: []string{"h2"}, ProvenanceStatus: "explicit", SessionIDs: []string{"s1"}},
		},
	}

	got := conservativeGenerateProposal(proposal, "ship alpha", "test fallback")

	if len(got.Revisions) != 1 {
		t.Fatalf("fallback revisions = %#v, want one whole-file revision", got.Revisions)
	}
	revision := got.Revisions[0]
	if revision.UseHunks || len(revision.HunkIDs) != 0 || len(revision.Hunks) != 0 {
		t.Fatalf("fallback revision = %#v, want whole-file coverage", revision)
	}
	if len(revision.Files) != 1 || revision.Files[0] != "alpha.go" {
		t.Fatalf("fallback files = %#v, want alpha.go", revision.Files)
	}
	if revision.ProvenanceStatus != "explicit" || len(revision.SessionIDs) != 1 || revision.SessionIDs[0] != "s1" {
		t.Fatalf("fallback provenance = %q/%#v, want explicit session", revision.ProvenanceStatus, revision.SessionIDs)
	}
	if got.Confidence.EffectiveConfidence < generateHighLogicConfidence {
		t.Fatalf("fallback confidence = %.2f, want high", got.Confidence.EffectiveConfidence)
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
	if docs.Intent != "update GX documentation" {
		t.Fatalf("docs intent = %q, want deterministic docs name", docs.Intent)
	}
	if len(docs.Files) != 2 || docs.Files[0] != "README.md" || docs.Files[1] != "docs/index.mdx" {
		t.Fatalf("docs files = %#v, want README and docs index", docs.Files)
	}
	if docs.TargetStack != "docs/documentation" {
		t.Fatalf("docs target = %q, want deterministic docs route", docs.TargetStack)
	}
	if docs.RouteConfidence != 0.96 {
		t.Fatalf("route confidence = %.2f, want deterministic route confidence", docs.RouteConfidence)
	}
	if len(docs.SessionIDs) != 2 || docs.SessionIDs[0] != "s1" || docs.SessionIDs[1] != "s2" {
		t.Fatalf("session ids = %#v, want de-duplicated sessions", docs.SessionIDs)
	}
	if got.Revisions[1].Intent != "add gx generate --legacy flag" {
		t.Fatalf("cli intent = %q, want deterministic legacy flag name", got.Revisions[1].Intent)
	}
	if got.Revisions[1].TargetStack != "feature/cli" {
		t.Fatalf("cli target = %q, want deterministic cli route", got.Revisions[1].TargetStack)
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
