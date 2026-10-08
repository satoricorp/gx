package cli

import (
	"strings"
	"testing"
)

// noReviewerEnv is a default build with nothing to review with: no gx Cloud
// (DenyNetwork points GX_CLOUD_URL at its tripwire, which would count as one),
// no Anthropic credential, and no request for the caller's own AWS account.
func noReviewerEnv(t *testing.T) {
	t.Helper()
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_CLOUD_URL", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "")
	t.Setenv("GX_REVIEW_AI", "")
	t.Setenv("GX_REVIEW_BEDROCK_DIRECT", "")
	t.Setenv("GX_REVIEW_MODELS", "")
	t.Setenv("GX_REVIEW_BEDROCK_MODEL_A", "")
	t.Setenv("GX_REVIEW_BEDROCK_MODEL_B", "")
}

// A review with no reviewer produced a full report, a verdict, and a findings
// file with no model having read a line of the change. Refusing is not a new
// restriction — there was no reviewer either way — it is the difference
// between saying so and shipping something that reads like a review.
func TestReviewRefusesWithoutAReviewer(t *testing.T) {
	noReviewerEnv(t)

	err := requireReviewer()
	if err == nil {
		t.Fatal("requireReviewer() = nil with no model configured, want a refusal")
	}
	// The message has to name the fixes, or it reads like the DEGRADED line it
	// replaces. None of them may be a service that does not exist.
	for _, want := range []string{"gx review cannot run", "ANTHROPIC_API_KEY", "GX_REVIEW_AI=0"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error = %v, want it to mention %q", err, want)
		}
	}
	if strings.Contains(err.Error(), "gx auth login") {
		t.Fatalf("error = %v sends a build with no gx Cloud to sign in to one", err)
	}
}

// An Anthropic key is the whole of a local review: no gx Cloud, no AWS.
func TestReviewRunsOnAnAnthropicKeyWithoutACloud(t *testing.T) {
	noReviewerEnv(t)
	t.Setenv("ANTHROPIC_API_KEY", "sk-ant-test")

	if err := requireReviewer(); err != nil {
		t.Fatalf("requireReviewer() = %v with ANTHROPIC_API_KEY set, want nil", err)
	}
}

// A build that does configure a cloud keeps its rule: the cloud is the wire, and
// being signed out of it is the problem to report, whatever else is exported.
func TestReviewWithACloudConfiguredStillNeedsItsSession(t *testing.T) {
	noReviewerEnv(t)
	t.Setenv("GX_CLOUD_URL", "https://cloud.example.invalid")
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("ANTHROPIC_API_KEY", "sk-ant-test")

	err := requireReviewer()
	if err == nil || !strings.Contains(err.Error(), "gx auth login") {
		t.Fatalf("requireReviewer() = %v with a cloud configured and no session, want the sign-in remedy", err)
	}
}

// Both opt-outs are still real reviews: one asks for the deterministic rules
// alone, the other reviews on the caller's own AWS account.
func TestReviewAllowsTheExplicitOptOuts(t *testing.T) {
	noReviewerEnv(t)

	t.Setenv("GX_REVIEW_AI", "0")
	if err := requireReviewer(); err != nil {
		t.Fatalf("requireReviewer() = %v with the AI panel off, want nil", err)
	}

	t.Setenv("GX_REVIEW_AI", "")
	t.Setenv("GX_REVIEW_BEDROCK_DIRECT", "1")
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIAEXAMPLE")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	if err := requireReviewer(); err != nil {
		t.Fatalf("requireReviewer() = %v with direct AWS, want nil", err)
	}
}
