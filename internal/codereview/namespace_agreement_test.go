package codereview

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/satoricorp/gx/internal/semantic"
)

// TestEnhanceReadsTheNamespaceIndexWrites is the end-to-end half of the
// namespace regression: the real `gx index` write path against the real
// `gx enhance` read path, with nothing re-derived in between.
//
// The observed failure: `gx index` in a checkout with no remote wrote
// gx-local-yeet-8d862445e7e4-v2, while `gx review` in that same checkout probed
// gx-local-satoricorp-yeet-v2 and repo-satoricorp-yeet, found neither, and
// reported "gx Cloud has never indexed this repository" — over an index that
// existed, was current, and held every chunk the review needed.
//
// Both sides now go through semantic.ResolveRepoIdentity. This test is what
// stops them going through anything else: it calls IndexRepository (with a fake
// embedder and store, so it needs no credentials) and codeIndexTargets, and
// asserts the namespace the first wrote is one the second searches.
func TestEnhanceReadsTheNamespaceIndexWrites(t *testing.T) {
	cases := []struct {
		name  string
		setup func(t *testing.T, root string)
	}{
		{
			name: "repository with a git remote",
			setup: func(t *testing.T, root string) {
				gitConfig(t, root, "remote.origin.url", "git@github.com:satoricorp/yeet.git")
			},
		},
		{
			name:  "repository with no git remote",
			setup: func(t *testing.T, root string) {},
		},
		{
			// Indexing consulted origin only, review consulted origin then
			// upstream: a fork-style checkout got a path-derived name from one
			// and a slug from the other, at the same instant.
			name: "repository with upstream but no origin",
			setup: func(t *testing.T, root string) {
				gitConfig(t, root, "remote.upstream.url", "git@github.com:satoricorp/yeet.git")
			},
		},
		{
			// Indexing ran `git remote get-url`, which expands insteadOf
			// rewrites; review ran `git config --get`, which returns the raw
			// value and parses "gh:satoricorp/yeet" as no repository at all.
			name: "remote written through an insteadOf rewrite",
			setup: func(t *testing.T, root string) {
				gitConfig(t, root, "url.https://github.com/.insteadOf", "gh:")
				gitConfig(t, root, "remote.origin.url", "gh:satoricorp/yeet.git")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			root := initRepo(t)
			tc.setup(t, root)
			writeFile(t, root, "main.go", "package main\n\nfunc main() {}\n")

			written := runIndexForTest(ctx, t, root)

			targets := codeIndexTargets(ctx, root)
			if len(targets) == 0 {
				t.Fatal("review resolved no namespaces to search")
			}
			var searched []string
			for _, target := range targets {
				searched = append(searched, target.Namespace)
			}
			if targets[0].Namespace != written {
				t.Fatalf("`gx index` wrote %q; review's first namespace is %q (searched %v)", written, targets[0].Namespace, searched)
			}
		})
	}
}

// TestEnhanceStillFindsAnIndexWrittenBeforeTheRemoteExisted covers what the
// shared resolver cannot: repository identity is derived from mutable git
// config, so adding a remote renames the namespace and orphans the index
// already written under the old name. That is exactly how the live divergence
// arose — indexed at 07:34 with no remote, reviewed at 12:05 with one.
func TestEnhanceStillFindsAnIndexWrittenBeforeTheRemoteExisted(t *testing.T) {
	ctx := context.Background()
	root := initRepo(t)
	writeFile(t, root, "main.go", "package main\n\nfunc main() {}\n")

	preRemote := runIndexForTest(ctx, t, root)
	gitConfig(t, root, "remote.origin.url", "git@github.com:satoricorp/yeet.git")

	var searched []string
	found := ""
	for _, target := range codeIndexTargets(ctx, root) {
		searched = append(searched, target.Namespace)
		if target.Namespace == preRemote {
			found = target.Origin
		}
	}
	if found == "" {
		t.Fatalf("review searched %v and would report this repository as never indexed; its index is at %q", searched, preRemote)
	}
	if found != semantic.NamespaceOriginPreRemote {
		t.Errorf("pre-remote namespace reported as %q, want it labelled so the evidence line says what was searched", found)
	}
}

func runIndexForTest(ctx context.Context, t *testing.T, root string) string {
	t.Helper()
	cfg := semantic.Config{OpenAIEmbeddingModel: "test-model", EmbeddingDimensions: 4}
	result, err := semantic.IndexRepository(ctx, semantic.RepoIndexOptions{
		RepoRoot:  root,
		OrgID:     reviewOrgID(),
		Reason:    "test",
		Config:    &cfg,
		Embedder:  namespaceTestEmbedder{},
		Store:     namespaceTestStore{},
		StatePath: filepath.Join(t.TempDir(), "manifest.json"),
	})
	if err != nil {
		t.Fatalf("IndexRepository() error = %v", err)
	}
	return result.Namespace
}

type namespaceTestEmbedder struct{}

func (namespaceTestEmbedder) Embed(_ context.Context, inputs []string) ([][]float32, error) {
	out := make([][]float32, len(inputs))
	for i := range inputs {
		out[i] = make([]float32, 4)
	}
	return out, nil
}

type namespaceTestStore struct{}

func (namespaceTestStore) Upsert(context.Context, []semantic.VectorRow) error { return nil }

func (namespaceTestStore) DeleteRows(context.Context, []string) error { return nil }

func gitConfig(t *testing.T, root, key, value string) {
	t.Helper()
	runGit(t, root, "config", key, value)
}
