package semantic

import (
	"context"
	"os/exec"
	"strings"
)

// Repository identity, resolved in exactly one place.
//
// A TurboPuffer namespace is a function of who the repository is, and the
// writer and the reader have to agree on that function or the reader addresses
// an empty namespace and reports the repository as never indexed. They did not
// agree. `lgtm index` asked git one way and parsed the answer with one parser;
// `lgtm review` asked git a different way and parsed the answer with a different
// copy of the same parser. Verified live: a `lgtm index` run wrote
// lgtm-local-yeet-8d862445e7e4-v2 while `lgtm review` in the same checkout looked in
// lgtm-local-satoricorp-yeet-v2 and reported "lgtm Cloud has never indexed this
// repository" over an index that existed and was current.
//
// This is the third instance of one bug shape in this codebase — a writer and a
// reader each deriving an identity for the same object — after reviewbundle vs
// UpsertRepo and stashed vs live capture rows. So the fix is shaped against the
// class rather than the instance: ResolveRepoIdentity is the single derivation,
// it returns both the write target and every namespace a read should try, and
// there is a test asserting the two paths agree that fails if anyone re-opens
// the gap.
//
// Two ways they could still drift are closed here rather than left to
// convention:
//
//   - Which remote. Indexing consulted only `origin`; review consulted `origin`
//     then `upstream`. A checkout with an upstream and no origin got a
//     path-derived name from one and a slug from the other.
//   - How git is asked. Indexing ran `git remote get-url`, which expands
//     url.<base>.insteadOf rewrites; review ran `git config --get
//     remote.origin.url`, which returns the raw value. With `gh:owner/repo` and
//     an insteadOf rule — an ordinary setup — the first yields owner/repo and
//     the second yields nothing at all.
//
// What is NOT fixed, because it cannot be: identity depends on mutable git
// config. Adding a remote to a repository that already has an index renames its
// namespace and orphans everything written under the old name. That is exactly
// how the live divergence arose — the checkout was indexed at 07:34 with no
// remote and reviewed at 12:05 with one. Candidates() answers it on the read
// side by probing the pre-remote name too, and the review reports which
// namespace actually answered.

// RepoIdentity is who a repository is, for indexing purposes.
type RepoIdentity struct {
	// RepoRoot is the absolute checkout path.
	RepoRoot string
	// RepoFullName is "owner/name" from the checkout's remotes, or "" when it
	// has none. It is never guessed: a repository with no remote has no owner,
	// and inventing one would silently give two different checkouts of the same
	// directory name the same identity.
	RepoFullName string
	// OrgID is the signed-in lgtm org, or "" for a local-only index.
	OrgID string
	// Namespace is where an index WRITE goes.
	Namespace string
}

// NamespaceCandidate is one namespace a read should try, and why.
type NamespaceCandidate struct {
	Namespace string
	// Origin explains the namespace to a human reading the evidence line.
	Origin string
}

// Namespace origins, as they appear in a review's evidence line.
const (
	NamespaceOriginPrimary = "lgtm code index"
	// NamespaceOriginPreRemote is the name this repository's index was written
	// under before the checkout gained a git remote.
	NamespaceOriginPreRemote = "lgtm code index (pre-remote name)"
)

// ResolveRepoIdentity derives a repository's identity from its checkout.
//
// repoFullName may be supplied by a caller that already knows it (the publish
// path does); otherwise it is read from the checkout's remotes. Everything
// downstream — the write namespace and every read candidate — is derived from
// the result, so a caller cannot accidentally use a different one.
func ResolveRepoIdentity(ctx context.Context, repoRoot, orgID, repoFullName string) RepoIdentity {
	repoRoot = strings.TrimSpace(repoRoot)
	repoFullName = strings.TrimSpace(repoFullName)
	if repoFullName == "" {
		repoFullName = RepoFullNameForRoot(ctx, repoRoot)
	}
	orgID = strings.TrimSpace(orgID)
	return RepoIdentity{
		RepoRoot:     repoRoot,
		RepoFullName: repoFullName,
		OrgID:        orgID,
		Namespace:    NamespaceForRepo(orgID, repoFullName, repoRoot),
	}
}

// Candidates lists the namespaces a read should try, most specific first.
//
// The primary is the one a write would go to now — and since the console,
// the server, and this CLI all converged on `lgtm-<orgId>-<slug>-v2` (console
// #58), it is also the cloud indexers' name: there is no separate "console
// namespace" to probe anymore. A `ConsoleNamespaceForRepo` guess at the old
// `repo-<owner>-<repo>` name lived here until it silently drifted — the
// convex function it claimed to mirror had been renamed and re-shaped — which
// is this file's bug class exactly, so the reader now carries no local copy
// of any other system's naming. The pre-remote name remains, because an index
// may predate this checkout's git remote; that costs one metadata probe and
// buys back an index that would otherwise be reported as absent.
func (r RepoIdentity) Candidates() []NamespaceCandidate {
	var out []NamespaceCandidate
	seen := map[string]struct{}{}
	add := func(namespace, origin string) {
		namespace = strings.TrimSpace(namespace)
		if namespace == "" {
			return
		}
		if _, ok := seen[namespace]; ok {
			return
		}
		seen[namespace] = struct{}{}
		out = append(out, NamespaceCandidate{Namespace: namespace, Origin: origin})
	}
	add(r.Namespace, NamespaceOriginPrimary)
	if r.RepoFullName != "" {
		add(NamespaceForRepo(r.OrgID, "", r.RepoRoot), NamespaceOriginPreRemote)
	}
	return out
}

// RepoFullNameForRoot reads "owner/name" from a checkout's remotes.
//
// `git remote get-url` rather than `git config --get remote.<n>.url` because it
// expands url.<base>.insteadOf, so a remote configured as `gh:owner/repo` with a
// rewrite rule resolves to the repository it actually points at. The config
// read stays as a fallback for the case where get-url is unavailable, and both
// are parsed by the same function.
//
// origin then upstream, and nothing else: a fork checkout names its own fork
// `origin`, which is the repository whose index this is.
func RepoFullNameForRoot(ctx context.Context, repoRoot string) string {
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return ""
	}
	for _, remote := range []string{"origin", "upstream"} {
		for _, args := range [][]string{
			{"remote", "get-url", remote},
			{"config", "--get", "remote." + remote + ".url"},
		} {
			if full := repoFullNameFromRemoteURL(gitOutputContext(ctx, repoRoot, args...)); full != "" {
				return full
			}
		}
	}
	return ""
}

func gitOutputContext(ctx context.Context, root string, args ...string) string {
	if ctx == nil {
		ctx = context.Background()
	}
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", root}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
