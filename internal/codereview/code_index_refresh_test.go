package codereview

import (
	"context"
	"testing"
)

// TestRefreshCodeIndexSkipsWithoutCredentials pins that the refresh seam is
// inert when there is nothing to refresh into. It must not fabricate evidence
// and must not attempt a network call.
func TestRefreshCodeIndexSkipsWithoutCredentials(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GX_OPENAI_API_KEY", "")
	t.Setenv("TURBOPUFFER_API_KEY", "")

	log := &EvidenceLog{}
	refreshCodeIndex(context.Background(), t.TempDir(), log)

	if statuses := log.Statuses(); len(statuses) != 0 {
		t.Fatalf("Statuses() = %#v, want no evidence recorded when unconfigured", statuses)
	}
}

// TestRefreshCodeIndexHonoursKillSwitch pins the escape hatch: a review must be
// able to run without touching the index at all.
func TestRefreshCodeIndexHonoursKillSwitch(t *testing.T) {
	t.Setenv("GX_REVIEW_INDEX_REFRESH", "0")
	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("TURBOPUFFER_API_KEY", "test-key")

	log := &EvidenceLog{}
	refreshCodeIndex(context.Background(), t.TempDir(), log)

	if statuses := log.Statuses(); len(statuses) != 0 {
		t.Fatalf("Statuses() = %#v, want the refresh skipped entirely", statuses)
	}
}

// TestRefreshCodeIndexSkipsEmptyRepoRoot guards the seam against being called
// before the repository root is known.
func TestRefreshCodeIndexSkipsEmptyRepoRoot(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("TURBOPUFFER_API_KEY", "test-key")

	log := &EvidenceLog{}
	refreshCodeIndex(context.Background(), "  ", log)

	if statuses := log.Statuses(); len(statuses) != 0 {
		t.Fatalf("Statuses() = %#v, want no refresh without a repository root", statuses)
	}
}
