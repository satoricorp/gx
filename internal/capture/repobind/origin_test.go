package repobind

import "testing"

// Every spelling of one repository has to compare equal, because this function
// is the gate: a repository that stops being recognized silently stops having
// its sessions indexed, and two repositories wrongly treated as one put a
// session somewhere it does not belong.
func TestNormalizeOriginCanonicalizesEverySpelling(t *testing.T) {
	same := []string{
		"git@github.com:satoricorp/totality.git",
		"git@github.com:satoricorp/totality",
		"https://github.com/satoricorp/totality.git",
		"https://github.com/satoricorp/totality",
		"https://github.com/satoricorp/totality/",
		"ssh://git@github.com/satoricorp/totality.git",
		"ssh://git@github.com:22/satoricorp/totality.git",
		"git://github.com/satoricorp/totality.git",
		"https://GitHub.com/SatoriCorp/Totality.git",
		"github.com/satoricorp/totality",
	}
	const want = "github.com/satoricorp/totality"
	for _, raw := range same {
		if got := NormalizeOrigin(raw); got != want {
			t.Errorf("NormalizeOrigin(%q) = %q, want %q", raw, got, want)
		}
	}
}

// An origin URL can carry a token. The normalized form is what gets persisted,
// so credentials have to be gone by then.
func TestNormalizeOriginDropsCredentials(t *testing.T) {
	for _, raw := range []string{
		"https://x-access-token:ghs_supersecrettoken@github.com/satoricorp/totality.git",
		"https://joe:hunter2@github.com/satoricorp/totality.git",
		"ssh://git@github.com/satoricorp/totality.git",
	} {
		got := NormalizeOrigin(raw)
		if got != "github.com/satoricorp/totality" {
			t.Errorf("NormalizeOrigin(%q) = %q", raw, got)
		}
		for _, secret := range []string{"ghs_supersecrettoken", "hunter2", "x-access-token", "joe:"} {
			if contains(got, secret) {
				t.Errorf("NormalizeOrigin(%q) leaked %q: %q", raw, secret, got)
			}
		}
	}
}

// Self-hosted forges are repositories too; nothing here may assume GitHub.
func TestNormalizeOriginKeepsTheHost(t *testing.T) {
	cases := map[string]string{
		"git@gitlab.company.com:platform/api.git":     "gitlab.company.com/platform/api",
		"https://bitbucket.org/team/thing.git":        "bitbucket.org/team/thing",
		"git@github.com:satoricorp/totality.git":      "github.com/satoricorp/totality",
		"https://ghe.internal.example/org/repo.git":   "ghe.internal.example/org/repo",
		"git@github.com:satoricorp/nested/deep/x.git": "github.com/satoricorp/nested/deep/x",
	}
	for raw, want := range cases {
		if got := NormalizeOrigin(raw); got != want {
			t.Errorf("NormalizeOrigin(%q) = %q, want %q", raw, got, want)
		}
	}
}

// Different repositories must never collapse together. A gate that returns the
// same identity for two repos is how a session reaches the wrong index.
func TestNormalizeOriginKeepsDifferentReposApart(t *testing.T) {
	distinct := []string{
		"git@github.com:satoricorp/totality.git",
		"git@github.com:satoricorp/console.git",
		"git@github.com:joe/tx.git",
		"git@gitlab.com:satoricorp/totality.git",
		"git@github.com:satoricorp/totality-internal.git",
	}
	seen := map[string]string{}
	for _, raw := range distinct {
		got := NormalizeOrigin(raw)
		if got == "" {
			t.Fatalf("NormalizeOrigin(%q) returned empty", raw)
		}
		if prev, ok := seen[got]; ok {
			t.Fatalf("%q and %q both normalized to %q", prev, raw, got)
		}
		seen[got] = raw
	}
}

// Empty means "not identifiable" and must never behave as a wildcard: a local
// repo with no remote has to stay unbound rather than match something.
func TestNormalizeOriginRejectsUnusableInput(t *testing.T) {
	for _, raw := range []string{"", "   ", "not a url", "https://", "github.com", "/Users/joe/git/tx"} {
		if got := NormalizeOrigin(raw); got != "" {
			t.Errorf("NormalizeOrigin(%q) = %q, want empty", raw, got)
		}
	}
}

func contains(haystack, needle string) bool {
	return len(needle) > 0 && len(haystack) >= len(needle) && indexOf(haystack, needle) >= 0
}

func indexOf(haystack, needle string) int {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}
