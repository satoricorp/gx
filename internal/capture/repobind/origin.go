// Package repobind answers which repository an agent session was working in.
//
// Attribution used to answer that question by matching the text an agent wrote
// against the lines a commit added. That works, but it conflates two jobs:
// deciding what a session produced, and deciding what a session belongs to.
// Only the second one can gate what leaves the machine, and it has to hold for
// sessions whose work was never committed, or whose transcript had not finished
// being written when capture read it.
//
// The binding runs: session -> working directory -> git origin -> connected
// repo. This package owns the middle two hops.
package repobind

import (
	"fmt"
	"net/url"
	"strings"
)

// NormalizeOrigin reduces a git remote URL to a canonical "host/owner/repo",
// lowercased, so the spellings of one repository compare equal:
//
//	git@github.com:Org/Repo.git
//	https://github.com/org/repo
//	ssh://git@github.com:22/org/repo.git
//	https://x-access-token:ghs_abc@github.com/org/repo.git
//
// all become "github.com/org/repo".
//
// This is the gate's identity function. Fail to normalize two spellings of the
// same repo and a connected repository silently stops being recognized; treat
// two different repos as equal and a session reaches an index it does not
// belong in. Both failures are quiet, which is why this is a pure function with
// its own tests rather than something assembled inline at the call site.
//
// Empty is returned for anything unrecognized. Callers must treat empty as
// "not identifiable", never as a wildcard.
func NormalizeOrigin(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	host, path := splitOrigin(raw)
	if host == "" || path == "" {
		return ""
	}

	// Credentials belong to nobody's identity, and an origin URL can carry a
	// token. Dropping the userinfo here is also what keeps one out of the
	// database: the normalized form is what gets persisted.
	if at := strings.LastIndex(host, "@"); at >= 0 {
		host = host[at+1:]
	}
	// A port says how to reach the host, not which host it is.
	if colon := strings.LastIndex(host, ":"); colon >= 0 {
		if _, err := fmt.Sscanf(host[colon+1:], "%d", new(int)); err == nil {
			host = host[:colon]
		}
	}

	path = strings.Trim(path, "/")
	path = strings.TrimSuffix(path, ".git")
	path = strings.Trim(path, "/")
	if host == "" || path == "" {
		return ""
	}
	// Host and owner are case-insensitive; repository names are compared
	// case-insensitively by every forge this targets, so the whole thing
	// lowercases rather than guessing which segment is safe.
	return strings.ToLower(host + "/" + path)
}

// splitOrigin separates the host from the repository path across the shapes git
// accepts: scp-like (git@host:owner/repo.git), a real URL, and bare host/path.
func splitOrigin(raw string) (host, path string) {
	if strings.Contains(raw, "://") {
		parsed, err := url.Parse(raw)
		if err != nil || parsed.Host == "" {
			return "", ""
		}
		return parsed.Host, parsed.Path
	}
	// scp-like syntax has no scheme and separates host from path with a colon.
	// It is only scp-like when the colon is not a port, which is what
	// distinguishes git@host:owner/repo from host:1234/owner/repo.
	if at := strings.Index(raw, "@"); at >= 0 {
		rest := raw[at+1:]
		if colon := strings.Index(rest, ":"); colon >= 0 {
			return rest[:colon], rest[colon+1:]
		}
		return "", ""
	}
	if colon := strings.Index(raw, ":"); colon >= 0 && !strings.Contains(raw[:colon], "/") {
		return raw[:colon], raw[colon+1:]
	}
	// A bare host/owner/repo with no scheme and no colon.
	if slash := strings.Index(raw, "/"); slash > 0 {
		return raw[:slash], raw[slash:]
	}
	return "", ""
}

// SameRepo reports whether two remote URLs name the same repository.
func SameRepo(a, b string) bool {
	na, nb := NormalizeOrigin(a), NormalizeOrigin(b)
	return na != "" && na == nb
}
