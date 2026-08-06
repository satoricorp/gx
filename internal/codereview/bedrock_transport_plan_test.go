package codereview

import (
	"strings"
	"testing"
)

// Cloud is the default wire. These tests pin that ambient AWS credentials — the
// kind a developer has exported for unrelated work — never silently choose the
// direct-to-AWS path.

func TestResolveTransportIgnoresAmbientAWSCredentials(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIAEXAMPLE")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	t.Setenv(bedrockDirectEnvVar, "")
	// No cloud URL configured, so planning must fail rather than fall back to
	// the credentials above.
	t.Setenv("GX_CLOUD_URL", "")

	_, err := resolveBedrockTransportPlan()
	if err == nil {
		t.Fatal("expected an error: ambient AWS credentials must not select the direct wire")
	}
	if !strings.Contains(err.Error(), bedrockDirectEnvVar) {
		t.Errorf("error should name the opt-in var, got: %v", err)
	}
}

func TestResolveTransportUsesDirectOnlyWhenAsked(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIAEXAMPLE")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	t.Setenv("AWS_REGION", "us-west-2")
	t.Setenv(bedrockDirectEnvVar, "1")

	plan, err := resolveBedrockTransportPlan()
	if err != nil {
		t.Fatalf("resolveBedrockTransportPlan() error = %v", err)
	}
	if plan.Kind != bedrockTransportKindDirect {
		t.Fatalf("Kind = %q, want %q", plan.Kind, bedrockTransportKindDirect)
	}
}

func TestResolveTransportDirectWithoutCredentialsDoesNotFallBackToCloud(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	t.Setenv(bedrockDirectEnvVar, "1")
	t.Setenv("GX_CLOUD_URL", "https://api.example.com")

	plan, err := resolveBedrockTransportPlan()
	if err == nil {
		t.Fatalf("expected an error, got plan kind %q", plan.Kind)
	}
	// Rerouting an explicit direct request to Cloud would hide the
	// misconfiguration the caller needs to see.
	if plan.Kind == bedrockTransportKindCloud {
		t.Error("explicit direct request must not silently reroute to Cloud")
	}
}

func TestBedrockDirectRequestedAcceptsCommonTruthyValues(t *testing.T) {
	for _, on := range []string{"1", "true", "TRUE", "yes", "on"} {
		t.Setenv(bedrockDirectEnvVar, on)
		if !bedrockDirectRequested() {
			t.Errorf("bedrockDirectRequested() = false for %q", on)
		}
	}
	for _, off := range []string{"", "0", "false", "no", "off", "maybe"} {
		t.Setenv(bedrockDirectEnvVar, off)
		if bedrockDirectRequested() {
			t.Errorf("bedrockDirectRequested() = true for %q", off)
		}
	}
}
