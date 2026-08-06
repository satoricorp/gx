package redact_test

import (
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/capture/redact"
)

func TestRedact_AWSKey(t *testing.T) {
	secret := "AKIAIOSFODNN7EXAMPLE"
	raw := "export AWS_ACCESS_KEY_ID=" + secret
	out := redact.Redact(raw)
	if strings.Contains(out, secret) {
		t.Fatalf("redacted output still contains secret: %q", out)
	}
	if !strings.Contains(out, "[REDACTED:aws-key:") {
		t.Fatalf("expected aws-key placeholder, got %q", out)
	}
}

func TestRedact_GitHubToken(t *testing.T) {
	secret := "ghp_1234567890abcdefghijklmnopqrstuvwxyz"
	out := redact.Redact("token=" + secret)
	if strings.Contains(out, secret) {
		t.Fatalf("redacted output still contains secret")
	}
}

func TestRedact_SlackToken(t *testing.T) {
	// Build at runtime so push protection does not flag a static Slack token shape.
	prefix := "xox" + "b"
	secret := prefix + "-1234567890-1234567890123-abcdefghijklmnopqrstuvwx"
	out := redact.Redact("slack " + secret)
	if strings.Contains(out, secret) {
		t.Fatalf("redacted output still contains secret")
	}
}

func TestRedact_PasswordKeyword(t *testing.T) {
	secret := "super-secret-password-123"
	out := redact.Redact("password=" + secret)
	if strings.Contains(out, secret) {
		t.Fatalf("redacted output still contains secret")
	}
}

func TestRedact_BearerKeyword(t *testing.T) {
	secret := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.payload"
	out := redact.Redact("Authorization: Bearer " + secret)
	if strings.Contains(out, secret) {
		t.Fatalf("redacted output still contains secret")
	}
}

func TestRedact_HighEntropy(t *testing.T) {
	secret := "aB3dE5fG7hI9jK1lM3nO5pQ7rS9tU1vW3xY5zA7bC9dE"
	out := redact.Redact("key=" + secret)
	if strings.Contains(out, secret) {
		t.Fatalf("redacted output still contains secret")
	}
}

func TestRedact_StablePlaceholder(t *testing.T) {
	secret := "AKIAIOSFODNN7EXAMPLE"
	a := redact.Redact(secret)
	b := redact.Redact(secret)
	if a != b {
		t.Fatalf("placeholders differ: %q vs %q", a, b)
	}
}

func TestRedact_PrivateKeyHeader(t *testing.T) {
	raw := "-----BEGIN RSA PRIVATE KEY-----\nMIIE..."
	out := redact.Redact(raw)
	if strings.Contains(out, "BEGIN RSA PRIVATE KEY") {
		t.Fatalf("private key header not redacted: %q", out)
	}
}
