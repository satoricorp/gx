package cli

import (
	"strings"
	"testing"
)

func TestReportRedactsSensitiveValues(t *testing.T) {
	raw := "Authorization: Bearer gxcs_secret\nGITHUB_TOKEN=gho_secret\napi_key=sk-secret"

	redacted := redactSensitive(raw)

	if strings.Contains(redacted, "gxcs_secret") || strings.Contains(redacted, "gho_secret") || strings.Contains(redacted, "sk-secret") {
		t.Fatalf("redactSensitive() leaked secret: %q", redacted)
	}
	if got := strings.Count(redacted, "[REDACTED]"); got < 3 {
		t.Fatalf("redactSensitive() redactions = %d, want at least 3 in %q", got, redacted)
	}
}
