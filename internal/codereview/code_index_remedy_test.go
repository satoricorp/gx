package codereview

import (
	"strings"
	"testing"

	"github.com/satoricorp/totality/internal/semantic"
)

// TestMissingIndexWarningsCarryTheirRemedy pins the property that makes these
// warnings worth printing: a user told their review ran on partial evidence
// must also be told how to fix it, and told once.
//
// Every route goes to the website. `tl index` is a hidden maintenance command
// that fills one developer's namespace from one developer's checkout; the
// index a review actually reads is the one Totality Cloud maintains from the GitHub
// App, and that is connected on the site.
func TestMissingIndexWarningsCarryTheirRemedy(t *testing.T) {
	// The shape a repository with a git remote and no index anywhere produces:
	// both candidate namespaces absent.
	statuses := []EvidenceStatus{
		{
			Source: codeIndexEvidenceSource, Namespace: "totality-001-acme-api-v2", State: EvidenceMissing,
			Detail: codeIndexMissingDetail(codeIndexTarget{Origin: semantic.NamespaceOriginPrimary}),
			Remedy: codeIndexMissingRemedy(codeIndexTarget{Origin: semantic.NamespaceOriginPrimary}),
		},
		{
			Source: codeIndexEvidenceSource, Namespace: "totality-local-api-abc-v2", State: EvidenceMissing,
			Detail: codeIndexMissingDetail(codeIndexTarget{Origin: semantic.NamespaceOriginPreRemote}),
			Remedy: codeIndexMissingRemedy(codeIndexTarget{Origin: semantic.NamespaceOriginPreRemote}),
		},
	}

	warnings := EvidenceWarnings(statuses)
	if len(warnings) != 1 {
		t.Fatalf("EvidenceWarnings() = %#v, want one warning for the one source", warnings)
	}
	warning := warnings[0]
	if !strings.Contains(warning, "https://totality.sh/repositories") {
		t.Fatalf("warning = %q, want it to say where to get the repository indexed", warning)
	}
	if strings.Contains(warning, "tl index") {
		t.Fatalf("warning = %q, want the website rather than the hidden `tl index` command", warning)
	}
	if count := strings.Count(warning, "https://totality.sh/repositories"); count != 1 {
		t.Fatalf("warning states the remedy %d times, want once:\n%s", count, warning)
	}
}

// TestEveryMissingCodeIndexOriginHasARemedy guards the switch against a new
// namespace origin being added without one, which is how a warning ends up
// telling someone their review is degraded and nothing else.
func TestEveryMissingCodeIndexOriginHasARemedy(t *testing.T) {
	for _, origin := range []string{
		semantic.NamespaceOriginPrimary,
		"some-future-origin",
	} {
		target := codeIndexTarget{Origin: origin}
		if remedy := codeIndexMissingRemedy(target); !strings.Contains(remedy, "https://totality.sh/repositories") {
			t.Errorf("origin %q remedy = %q, want it to point at the console", origin, remedy)
		}
		if detail := codeIndexMissingDetail(target); strings.TrimSpace(detail) == "" {
			t.Errorf("origin %q has no detail", origin)
		}
	}
}
