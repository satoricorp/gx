package cli

import (
	"strings"
	"testing"
)

// A signed-out review produced a full report, a verdict, and a findings file
// with no model having read a line of the change. Refusing is not a new
// restriction — there was no reviewer either way — it is the difference
// between saying so and shipping something that reads like a review.
func TestReviewRefusesWithoutACloudSession(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GX_REVIEW_AI", "")
	t.Setenv("GX_REVIEW_BEDROCK_DIRECT", "")

	err := requireReviewSession()
	if err == nil {
		t.Fatal("requireReviewSession() = nil while signed out, want a refusal")
	}
	// The message has to name the fix and say why, or it reads like the
	// DEGRADED line it replaces.
	if !strings.Contains(err.Error(), "no AI reviewer runs") {
		t.Fatalf("error = %v, want it to say no reviewer runs", err)
	}
	if !strings.Contains(err.Error(), "gx auth login") {
		t.Fatalf("error = %v, want the remedy named", err)
	}
}

// Both opt-outs are still real reviews and must keep working signed out: one
// asks for the deterministic rules alone, the other reviews on the caller's
// own AWS account with gx Cloud nowhere in the path.
func TestReviewAllowsTheExplicitOptOutsWhileSignedOut(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")

	t.Setenv("GX_REVIEW_BEDROCK_DIRECT", "")
	t.Setenv("GX_REVIEW_AI", "0")
	if err := requireReviewSession(); err != nil {
		t.Fatalf("requireReviewSession() = %v with the AI panel off, want nil", err)
	}

	t.Setenv("GX_REVIEW_AI", "")
	t.Setenv("GX_REVIEW_BEDROCK_DIRECT", "1")
	if err := requireReviewSession(); err != nil {
		t.Fatalf("requireReviewSession() = %v with direct AWS, want nil", err)
	}
}
