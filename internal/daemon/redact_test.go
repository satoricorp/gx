package daemon

import (
	"net/http"
	"strings"
	"testing"
)

func TestRedactHeaders(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer secret")
	headers.Set("X-Api-Key", "secret-key")
	headers.Set("Content-Type", "application/json")

	got := redactHeaders(headers)

	if strings.Contains(got, "secret") {
		t.Fatalf("expected secrets to be redacted, got %s", got)
	}
	if !strings.Contains(got, "[redacted]") {
		t.Fatalf("expected redaction marker, got %s", got)
	}
	if !strings.Contains(got, "application/json") {
		t.Fatalf("expected safe headers to remain, got %s", got)
	}
}
