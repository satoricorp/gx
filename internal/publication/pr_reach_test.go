package publication

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/reviewbundle"
)

func TestComputeLexicalReachCountsExternalReference(t *testing.T) {
	root, headSHA := initLexicalReachRepo(t)
	prURL := "https://github.com/example/reach/pull/1"
	artifact := reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: root},
		Push: reviewbundle.PushPayload{
			HeadCommitID:         headSHA,
			GitHubPullRequestURL: &prURL,
		},
		Stack: []reviewbundle.StackPayload{{
			Patch: strings.Join([]string{
				"diff --git a/lib.go b/lib.go",
				"--- a/lib.go",
				"+++ b/lib.go",
				"@@ -1,3 +1,3 @@",
				" package main",
				" ",
				"-func HelloWorld() {}",
				"+func HelloWorld() { /* changed */ }",
			}, "\n"),
			Change: reviewbundle.ChangePayload{
				CurrentCommitID: headSHA,
				Description:     "change HelloWorld",
				Files:           []string{"lib.go"},
				ReviewContext: &reviewbundle.ReviewContextPayload{
					ChangedSymbols: []reviewbundle.ReviewChangedSymbol{{
						File: "lib.go", Symbol: "HelloWorld", Kind: "func",
					}},
				},
			},
			GitHubPullRequestURL: &prURL,
		}},
	})

	catalog := buildPRBodyCatalog(artifact)
	reach := computeLexicalReach(context.Background(), root, catalog, artifact)
	if reach.ReferenceCount == 0 {
		t.Fatalf("reach = %#v, want external references", reach)
	}
	if reach.DependentFiles != 3 {
		t.Fatalf("DependentFiles = %d, want 3", reach.DependentFiles)
	}

	stats := catalog.Stats
	applyLexicalReachToStats(&stats, reach)
	if stats.MaxRiskLevel != "medium" {
		t.Fatalf("MaxRiskLevel = %q, want medium", stats.MaxRiskLevel)
	}

	catalog.Stats = stats
	body := renderGitHubPullRequestBody(artifact, catalog, prSummaryContext{}, nil, "", false, reach)
	if !strings.Contains(body, "references to changed symbols (HelloWorld)") {
		t.Fatalf("body missing reach readiness sentence:\n%s", body)
	}
	wantLink := blobPermalink(prURL, headSHA, "consumer.go", 4)
	if !strings.Contains(body, wantLink) {
		t.Fatalf("body missing blob permalink %q:\n%s", wantLink, body)
	}
	if !strings.Contains(body, "## Blast Radius") {
		t.Fatalf("body missing blast radius section:\n%s", body)
	}
}

func TestComputeLexicalReachSkipsShortSymbol(t *testing.T) {
	root, headSHA := initLexicalReachRepo(t)
	artifact := reviewbundle.NewArtifact(reviewbundle.Bundle{
		Repo: reviewbundle.RepoPayload{RootPath: root},
		Push: reviewbundle.PushPayload{HeadCommitID: headSHA},
		Stack: []reviewbundle.StackPayload{{
			Patch: "diff --git a/lib.go b/lib.go\n--- a/lib.go\n+++ b/lib.go\n@@ -1 +1,2 @@\n+func Run() {}\n",
			Change: reviewbundle.ChangePayload{
				Files: []string{"lib.go"},
				ReviewContext: &reviewbundle.ReviewContextPayload{
					ChangedSymbols: []reviewbundle.ReviewChangedSymbol{{Symbol: "Run", File: "lib.go"}},
				},
			},
		}},
	})
	reach := computeLexicalReach(context.Background(), root, buildPRBodyCatalog(artifact), artifact)
	if reach.ReferenceCount != 0 {
		t.Fatalf("reach = %#v, want short symbol skipped", reach)
	}
}

func TestComputeLexicalReachDisabledByEnv(t *testing.T) {
	t.Setenv("GX_PR_BLAST_RADIUS", "0")
	root, headSHA := initLexicalReachRepo(t)
	artifact := reachTestArtifact(root, headSHA)
	catalog := buildPRBodyCatalog(artifact)
	reach := computeLexicalReach(context.Background(), root, catalog, artifact)
	if reach.ReferenceCount != 0 {
		t.Fatalf("reach = %#v, want disabled", reach)
	}
	body := renderGitHubPullRequestBody(artifact, catalog, prSummaryContext{}, nil, "", false, reach)
	if strings.Contains(body, "references to changed symbols") {
		t.Fatalf("body should fall back without reach wording:\n%s", body)
	}
}

func TestComputeLexicalReachSkippedForDocsOnly(t *testing.T) {
	root := t.TempDir()
	artifact := docsOnlyPRArtifact()
	artifact.Bundle.Repo.RootPath = root
	reach := computeLexicalReach(context.Background(), root, buildPRBodyCatalog(artifact), artifact)
	if reach.ReferenceCount != 0 {
		t.Fatalf("reach = %#v, want skipped for docs-only", reach)
	}
}

func TestApplyLexicalReachToStatsHighFromManyFiles(t *testing.T) {
	stats := prBodyStats{MaxRiskLevel: "low"}
	reach := lexicalReach{
		ReferenceCount: 30,
		DependentFiles: 10,
		Symbols:        []symbolReach{{Symbol: "Foo", References: make([]reachSite, 5)}},
	}
	applyLexicalReachToStats(&stats, reach)
	if stats.MaxRiskLevel != "high" {
		t.Fatalf("MaxRiskLevel = %q, want high", stats.MaxRiskLevel)
	}
	if !strings.Contains(strings.Join(stats.RiskSignals, " "), "lexical_reach:30_refs_in_10_files") {
		t.Fatalf("RiskSignals = %#v", stats.RiskSignals)
	}
}

func TestLexicalReachContextSnippetsInReviewBrief(t *testing.T) {
	reach := lexicalReach{
		ReferenceCount: 2,
		DependentFiles: 1,
		Symbols: []symbolReach{{
			Symbol:     "HelloWorld",
			References: []reachSite{{File: "consumer.go", Line: 3}},
		}},
	}
	brief := prReviewBrief(reviewbundle.NewArtifact(reviewbundle.Bundle{}), prBodyCatalog{}, prSummaryContext{}, reach)
	found := false
	for _, snippet := range brief.Context {
		if snippet.Kind == "lexical_reach" && snippet.Ref == "HelloWorld" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("brief context = %#v, want lexical_reach snippet", brief.Context)
	}
}

func TestBodyStatsSizeAloneCapsAtMedium(t *testing.T) {
	artifact := reviewbundle.NewArtifact(reviewbundle.Bundle{
		Stack: []reviewbundle.StackPayload{{
			Patch: strings.Repeat("+\n", 450),
			Change: reviewbundle.ChangePayload{
				Files: []string{"a.go", "b.go", "c.go", "d.go", "e.go", "f.go", "g.go", "h.go", "i.go"},
			},
		}},
	})
	stats := bodyStats(artifact, revisionSummaries(artifact), nil)
	if stats.MaxRiskLevel == "high" {
		t.Fatalf("MaxRiskLevel = high from size alone, want medium cap")
	}
	if stats.MaxRiskLevel != "medium" {
		t.Fatalf("MaxRiskLevel = %q, want medium", stats.MaxRiskLevel)
	}
}

func TestComputeLexicalReachTruncatedForCommonSymbol(t *testing.T) {
	root := t.TempDir()
	runReachGit(t, root, "init")
	runReachGit(t, root, "config", "user.email", "t@example.com")
	runReachGit(t, root, "config", "user.name", "test")
	for i := 0; i < 210; i++ {
		name := filepath.Join("pkg", "f"+strconv.Itoa(i)+".go")
		writeReachFile(t, root, name, "package pkg\n\nfunc CommonName() {}\n")
	}
	writeReachFile(t, root, "changed.go", "package main\n\nfunc CommonName() { /* edit */ }\n")
	runReachGit(t, root, "add", ".")
	runReachGit(t, root, "commit", "-m", "bulk")
	artifact := reviewbundle.NewArtifact(reviewbundle.Bundle{
		Repo: reviewbundle.RepoPayload{RootPath: root},
		Stack: []reviewbundle.StackPayload{{
			Patch: "diff --git a/changed.go b/changed.go\n--- a/changed.go\n+++ b/changed.go\n@@ -1 +1,2 @@\n package main\n+func CommonName() {}\n",
			Change: reviewbundle.ChangePayload{
				Files: []string{"changed.go"},
				ReviewContext: &reviewbundle.ReviewContextPayload{
					ChangedSymbols: []reviewbundle.ReviewChangedSymbol{{Symbol: "CommonName", File: "changed.go"}},
				},
			},
		}},
	})
	reach := computeLexicalReach(context.Background(), root, buildPRBodyCatalog(artifact), artifact)
	if !reach.Truncated {
		t.Fatalf("reach = %#v, want truncated for common symbol", reach)
	}
	if reach.ReferenceCount != 0 {
		t.Fatalf("reach = %#v, want common symbol excluded from counts", reach)
	}
}

func initLexicalReachRepo(t *testing.T) (root, headSHA string) {
	t.Helper()
	root = t.TempDir()
	runReachGit(t, root, "init")
	runReachGit(t, root, "config", "user.email", "t@example.com")
	runReachGit(t, root, "config", "user.name", "test")
	writeReachFile(t, root, "go.mod", "module example.com/reach\n\ngo 1.22\n")
	writeReachFile(t, root, "consumer.go", "package main\n\nfunc use() {\n\tHelloWorld()\n}\n")
	writeReachFile(t, root, "consumer2.go", "package main\n\nfunc also() { HelloWorld() }\n")
	writeReachFile(t, root, "consumer3.go", "package main\n\nfunc more() { HelloWorld() }\n")
	writeReachFile(t, root, "lib.go", "package main\n\nfunc HelloWorld() {}\n")
	runReachGit(t, root, "add", ".")
	runReachGit(t, root, "commit", "-m", "initial")
	writeReachFile(t, root, "lib.go", "package main\n\nfunc HelloWorld() { /* changed */ }\n")
	runReachGit(t, root, "add", "lib.go")
	runReachGit(t, root, "commit", "-m", "change HelloWorld")
	headSHA = strings.TrimSpace(runReachGitOutput(t, root, "rev-parse", "HEAD"))
	return root, headSHA
}

func reachTestArtifact(root, headSHA string) reviewbundle.Artifact {
	prURL := "https://github.com/example/reach/pull/1"
	return reviewbundle.NewArtifact(reviewbundle.Bundle{
		Repo: reviewbundle.RepoPayload{RootPath: root},
		Push: reviewbundle.PushPayload{
			HeadCommitID:         headSHA,
			GitHubPullRequestURL: &prURL,
		},
		Stack: []reviewbundle.StackPayload{{
			Patch: "diff --git a/lib.go b/lib.go\n--- a/lib.go\n+++ b/lib.go\n@@ -1,3 +1,3 @@\n package main\n \n-func HelloWorld() {}\n+func HelloWorld() { /* changed */ }\n",
			Change: reviewbundle.ChangePayload{
				CurrentCommitID: headSHA,
				Files:           []string{"lib.go"},
				ReviewContext: &reviewbundle.ReviewContextPayload{
					ChangedSymbols: []reviewbundle.ReviewChangedSymbol{{Symbol: "HelloWorld", File: "lib.go"}},
				},
			},
			GitHubPullRequestURL: &prURL,
		}},
	})
}

func runReachGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func runReachGitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

func writeReachFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
