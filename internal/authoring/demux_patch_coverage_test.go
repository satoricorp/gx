package authoring

import (
	"errors"
	"testing"
)

func TestValidateDemuxPatchCoverageAllowsValidHunkLevelProposal(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{{ID: "h1", File: "alpha.go", Patch: "diff --git a/alpha.go b/alpha.go\n@@\n+alpha\n"}},
		Revisions: []RevisionProposal{{
			ID:       "u1",
			Intent:   "alpha",
			UseHunks: true,
			HunkIDs:  []string{"h1"},
		}},
	}

	if err := validateDemuxPatchCoverage(proposal); err != nil {
		t.Fatalf("validateDemuxPatchCoverage() error = %v", err)
	}
}

func TestValidateDemuxPatchCoverageAllowsValidWholeFileProposal(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "alpha.go", Patch: "diff --git a/alpha.go b/alpha.go\n@@\n+first\n"},
			{ID: "h2", File: "alpha.go", Patch: "diff --git a/alpha.go b/alpha.go\n@@\n+second\n"},
		},
		Revisions: []RevisionProposal{{
			ID:     "u1",
			Intent: "alpha",
			Files:  []string{"alpha.go"},
		}},
	}

	if err := validateDemuxPatchCoverage(proposal); err != nil {
		t.Fatalf("validateDemuxPatchCoverage() error = %v", err)
	}
}

func TestValidateDemuxPatchCoverageRejectsDuplicateHunkOwnership(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{{ID: "h1", File: "alpha.go", Patch: "patch one"}},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "alpha first", UseHunks: true, HunkIDs: []string{"h1"}},
			{ID: "u2", Intent: "alpha second", UseHunks: true, HunkIDs: []string{"h1"}},
		},
	}

	err := validateDemuxPatchCoverage(proposal)
	var typed DuplicateHunkAssignmentError
	if !errors.As(err, &typed) || typed.HunkID != "h1" || typed.FirstOwnerID != "u1" || typed.SecondOwnerID != "u2" {
		t.Fatalf("validateDemuxPatchCoverage() error = %#v, want duplicate hunk ownership", err)
	}
}

func TestValidateDemuxPatchCoverageRejectsUnknownHunkReference(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{{ID: "h1", File: "alpha.go", Patch: "patch one"}},
		Revisions: []RevisionProposal{{
			ID:       "u1",
			Intent:   "alpha",
			UseHunks: true,
			HunkIDs:  []string{"h9"},
		}},
	}

	err := validateDemuxPatchCoverage(proposal)
	var typed UnknownHunkReferenceError
	if !errors.As(err, &typed) || typed.RevisionID != "u1" || typed.HunkID != "h9" {
		t.Fatalf("validateDemuxPatchCoverage() error = %#v, want unknown hunk reference", err)
	}
}

func TestValidateDemuxPatchCoverageRejectsMissingSourceHunkOwnership(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "alpha.go", Patch: "patch one"},
			{ID: "h2", File: "alpha.go", Patch: "patch two"},
		},
		Revisions: []RevisionProposal{{
			ID:       "u1",
			Intent:   "alpha first",
			UseHunks: true,
			HunkIDs:  []string{"h1"},
		}},
	}

	err := validateDemuxPatchCoverage(proposal)
	var typed MissingHunkOwnerError
	if !errors.As(err, &typed) || typed.HunkID != "h2" || typed.File != "alpha.go" {
		t.Fatalf("validateDemuxPatchCoverage() error = %#v, want missing source hunk ownership", err)
	}
}

func TestValidateDemuxPatchCoverageRejectsMixedWholeFileAndHunkOwnership(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{{ID: "h1", File: "alpha.go", Patch: "patch one"}},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "alpha hunk", UseHunks: true, HunkIDs: []string{"h1"}},
			{ID: "u2", Intent: "alpha file", Files: []string{"alpha.go"}},
		},
	}

	err := validateDemuxPatchCoverage(proposal)
	var typed MixedHunkWholeFileCoverageError
	if !errors.As(err, &typed) || typed.HunkID != "h1" || typed.HunkRevisionID != "u1" || typed.WholeRevisionID != "u2" {
		t.Fatalf("validateDemuxPatchCoverage() error = %#v, want mixed coverage", err)
	}
}

func TestValidateDemuxPatchCoverageRejectsPatchBodyMismatch(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{{ID: "h1", File: "alpha.go", Patch: "diff --git a/alpha.go b/alpha.go\n@@\n+source\n"}},
		Revisions: []RevisionProposal{{
			ID:       "u1",
			Intent:   "alpha",
			UseHunks: true,
			HunkIDs:  []string{"h1"},
			Hunks:    []HunkRange{{ID: "h1", File: "alpha.go", Patch: "diff --git a/alpha.go b/alpha.go\n@@\n+altered\n"}},
		}},
	}

	err := validateDemuxPatchCoverage(proposal)
	var typed HunkPatchBodyMismatchError
	if !errors.As(err, &typed) || typed.RevisionID != "u1" || typed.HunkID != "h1" || typed.File != "alpha.go" {
		t.Fatalf("validateDemuxPatchCoverage() error = %#v, want patch body mismatch", err)
	}
}

func TestValidateDemuxPatchCoverageComparesPerHunkNotNaiveConcatenation(t *testing.T) {
	hunks := parseGitHunks(`diff --git a/app.go b/app.go
--- a/app.go
+++ b/app.go
@@ -1 +1 @@
-alpha
+bravo
@@ -10 +10 @@
-charlie
+delta
`)
	if len(hunks) != 2 {
		t.Fatalf("parseGitHunks() = %#v, want two hunks", hunks)
	}
	proposal := DemuxProposal{
		Hunks: hunks,
		Revisions: []RevisionProposal{{
			ID:       "u1",
			Intent:   "app",
			UseHunks: true,
			HunkIDs:  []string{"h1", "h2"},
			Hunks:    []HunkRange{hunks[0], hunks[1]},
		}},
	}

	if err := validateDemuxPatchCoverage(proposal); err != nil {
		t.Fatalf("validateDemuxPatchCoverage() error = %v", err)
	}
}
