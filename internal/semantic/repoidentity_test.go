package semantic

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// namespaceIndexWrites runs the real indexer over a real checkout — with a
// fake embedder and store, so no credentials and no network — and reports the
// namespace it actually wrote to. Re-deriving the name here instead would make
// the agreement test tautological: it has to observe the writer, not imitate it.
func namespaceIndexWrites(ctx context.Context, t *testing.T, repoRoot, orgID string) string {
	t.Helper()
	cfg := Config{OpenAIEmbeddingModel: "test-model", EmbeddingDimensions: 4}
	result, err := IndexRepository(ctx, RepoIndexOptions{
		RepoRoot:  repoRoot,
		OrgID:     orgID,
		Reason:    "test",
		Config:    &cfg,
		Embedder:  &recordingEmbedder{dims: 4},
		Store:     newRecordingStore(),
		StatePath: filepath.Join(t.TempDir(), "manifest.json"),
	})
	if err != nil {
		t.Fatalf("IndexRepository() error = %v", err)
	}
	return result.Namespace
}

// namespaceReviewReads is the first namespace a review looks in.
//
// internal/codereview builds its targets from RepoIdentity.Candidates(), so
// this is the review's derivation. The end-to-end pairing — the real
// codeIndexTargets against this — is
// TestReviewReadsTheNamespaceIndexWrites in internal/codereview.
func namespaceReviewReads(ctx context.Context, t *testing.T, repoRoot, orgID string) string {
	t.Helper()
	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		t.Fatalf("abs(%q): %v", repoRoot, err)
	}
	candidates := ResolveRepoIdentity(ctx, absRoot, orgID, "").Candidates()
	if len(candidates) == 0 {
		t.Fatal("review resolved no candidate namespaces")
	}
	if candidates[0].Origin != NamespaceOriginPrimary {
		t.Fatalf("first candidate origin = %q, want the primary index namespace", candidates[0].Origin)
	}
	return candidates[0].Namespace
}

// TestIndexAndReviewResolveTheSameNamespace is the regression test for the
// class, not for one instance of it.
//
// A writer and a reader deriving the same identity separately is how three
// separate defects reached production in this codebase. The observed one:
// `gx index` wrote gx-local-yeet-8d862445e7e4-v2 while `gx review` in the same
// checkout searched gx-local-satoricorp-yeet-v2 and reported the repository as
// never indexed. Each case below is a real way the two derivations diverged or
// could diverge; the test fails if anyone re-opens the gap.
func TestIndexAndReviewResolveTheSameNamespace(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name  string
		setup func(t *testing.T, root string)
		orgID string
	}{{
		// (a) A repository with a git remote.
		name: "github remote over ssh",
		setup: func(t *testing.T, root string) {
			gitSet(t, root, "remote.origin.url", "git@github.com:satoricorp/yeet.git")
		},
	}, {
		name: "github remote over https",
		setup: func(t *testing.T, root string) {
			gitSet(t, root, "remote.origin.url", "https://github.com/satoricorp/yeet.git")
		},
	}, {
		// (b) A repository with no remote at all. This is the state the live
		// divergence started from.
		name:  "no remote",
		setup: func(t *testing.T, root string) {},
	}, {
		// Indexing consulted origin only; review consulted origin then
		// upstream. A checkout with an upstream and no origin got a
		// path-derived name from one and a slug from the other.
		name: "upstream but no origin",
		setup: func(t *testing.T, root string) {
			gitSet(t, root, "remote.upstream.url", "git@github.com:satoricorp/yeet.git")
		},
	}, {
		// Indexing ran `git remote get-url`, which expands insteadOf rewrites;
		// review ran `git config --get remote.origin.url`, which does not. With
		// this ordinary setup the first resolved satoricorp/yeet and the second
		// resolved nothing.
		name: "remote written through an insteadOf rewrite",
		setup: func(t *testing.T, root string) {
			gitSet(t, root, "url.https://github.com/.insteadOf", "gh:")
			gitSet(t, root, "remote.origin.url", "gh:satoricorp/yeet.git")
		},
	}, {
		// Signed in: the org prefix has to come from the same place on both
		// sides too.
		name:  "signed in to an org",
		orgID: "2f273110-b6ce-4b3b-95d3-e7c0ca802e83",
		setup: func(t *testing.T, root string) {
			gitSet(t, root, "remote.origin.url", "git@github.com:satoricorp/yeet.git")
		},
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := initTestRepo(t)
			tc.setup(t, root)

			written := namespaceIndexWrites(ctx, t, root, tc.orgID)
			read := namespaceReviewReads(ctx, t, root, tc.orgID)
			if written != read {
				t.Fatalf("`gx index` writes %q but `gx review` reads %q — a review of this repository would report it as never indexed", written, read)
			}
			if strings.TrimSpace(written) == "" {
				t.Fatal("resolved an empty namespace")
			}
		})
	}
}

// TestReviewProbesThePreRemoteNamespace covers the case the shared resolver
// cannot make go away: identity is derived from mutable git config, so a
// repository indexed before it had a remote is indexed under a different name
// than it now answers to. The read side has to look in both.
func TestReviewProbesThePreRemoteNamespace(t *testing.T) {
	ctx := context.Background()
	root := initTestRepo(t)

	// Indexed while the checkout had no remote.
	before := namespaceIndexWrites(ctx, t, root, "")

	gitSet(t, root, "remote.origin.url", "git@github.com:satoricorp/yeet.git")
	after := namespaceIndexWrites(ctx, t, root, "")
	if before == after {
		t.Fatal("adding a remote did not change the namespace; this test no longer covers anything")
	}

	var found bool
	var searched []string
	for _, candidate := range ResolveRepoIdentity(ctx, root, "", "").Candidates() {
		searched = append(searched, candidate.Namespace)
		if candidate.Namespace == before {
			found = true
		}
	}
	if !found {
		t.Fatalf("review searched %v and would miss the pre-remote index at %q", searched, before)
	}
}

// TestRepoFullNameIsNeverGuessed pins the other half of identity. A guessed
// owner is worse than no owner: it silently gives two unrelated checkouts that
// happen to share a directory name the same index.
func TestRepoFullNameIsNeverGuessed(t *testing.T) {
	root := initTestRepo(t)
	if got := RepoFullNameForRoot(context.Background(), root); got != "" {
		t.Fatalf("RepoFullNameForRoot() = %q for a checkout with no remote, want empty", got)
	}
	identity := ResolveRepoIdentity(context.Background(), root, "", "")
	if identity.RepoFullName != "" {
		t.Fatalf("RepoFullName = %q, want empty", identity.RepoFullName)
	}
	if ConsoleNamespaceForRepo(identity.RepoFullName) != "" {
		t.Fatal("a repository with no remote must not be given a console namespace")
	}
	// Two checkouts of the same directory name must not collide.
	other := initTestRepo(t)
	if ResolveRepoIdentity(context.Background(), other, "", "").Namespace == identity.Namespace {
		t.Fatal("two distinct remote-less checkouts resolved to one namespace")
	}
}

func initTestRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	run(t, root, "git", "init")
	return root
}

func gitSet(t *testing.T, root, key, value string) {
	t.Helper()
	run(t, root, "git", "config", key, value)
}

func run(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
	}
}
