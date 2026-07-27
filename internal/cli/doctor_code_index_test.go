package cli

import (
	"strings"
	"testing"
	"time"
)

// A fresh index must not nag. Age alone is not a fault: an untouched repository
// stays correctly indexed however long it sits, so only drift from HEAD is.
func TestCodeIndexAgeStaysQuietUnderADay(t *testing.T) {
	if got := codeIndexAge(time.Now().UnixMilli()); got != "" {
		t.Fatalf("codeIndexAge(now) = %q, want empty so a fresh index reports nothing", got)
	}
	if got := codeIndexAge(0); got != "" {
		t.Fatalf("codeIndexAge(0) = %q, want empty when the manifest carries no timestamp", got)
	}
}

func TestCodeIndexAgeReportsDaysOnce(t *testing.T) {
	oneDay := time.Now().Add(-30 * time.Hour).UnixMilli()
	if got := codeIndexAge(oneDay); got != " (1 day old)" {
		t.Fatalf("codeIndexAge(30h ago) = %q, want singular day", got)
	}
	// 30 days behind HEAD is the drift measured on this repository that scored
	// 0.000 recall for code written after the indexed commit.
	thirtyDays := time.Now().Add(-30 * 24 * time.Hour).UnixMilli()
	got := codeIndexAge(thirtyDays)
	if !strings.Contains(got, "30 days old") {
		t.Fatalf("codeIndexAge(30d ago) = %q, want it to name the drift", got)
	}
}

func TestShortCommitLeavesShortIDsAlone(t *testing.T) {
	full := "729f7282da046d4de1e0b015161d8b6fb722b837"
	if got := shortCommit(full); got != "729f7282da04" {
		t.Fatalf("shortCommit(full) = %q, want 12 characters", got)
	}
	if got := shortCommit("abc123"); got != "abc123" {
		t.Fatalf("shortCommit(short) = %q, want it returned unchanged", got)
	}
	if got := shortCommit("  " + full + "  "); got != "729f7282da04" {
		t.Fatalf("shortCommit(padded) = %q, want surrounding space trimmed", got)
	}
}
