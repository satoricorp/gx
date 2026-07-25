package codereview

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Review modes name what a review actually looked at. They are the honest
// answer to "what did gx read?", and the difference between "reviewed and
// clean" and "never looked" lives here rather than in the findings list.
const (
	// ReviewModeWorkingTree is the interactive case: uncommitted work.
	ReviewModeWorkingTree = "working-tree"
	// ReviewModeRange is a commit range, "<base>...HEAD".
	ReviewModeRange = "range"
	// ReviewModeRepo is a whole-repo review: the caller asked for a scope,
	// prompt, or deep review, so the repo itself is the subject and no diff
	// is required.
	ReviewModeRepo = "repo"
	// ReviewModeNone means nothing was inspected. It is never a pass.
	ReviewModeNone = "none"
)

// baseCandidates mirrors detectDefaultBase in internal/capture/orchestrator so
// review and capture agree on what "the base branch" means.
var baseCandidates = []string{"main", "master", "origin/main", "origin/master"}

// ChangeSet is the resolved answer to "what is this review looking at?". It is
// the single seam the rest of the pipeline reads from: Files drives triage,
// static-tool scoping, and the patch-focused finding filter, while Range picks
// the git command used to diff each file (empty means the working tree).
type ChangeSet struct {
	Mode  string
	Base  string
	Range string
	Files []string
	// Target is a human phrase naming what was inspected, used both in the
	// report and in the "nothing to review" message.
	Target string
}

// Reviewed reports whether the change set actually gave the review something
// to read.
func (c ChangeSet) Reviewed() bool {
	return c.Mode != ReviewModeNone && c.Mode != ""
}

// resolveChangeSet picks the source of changes, in priority order:
//
//	a. an explicit base ref               -> "<base>...HEAD"
//	b. a dirty working tree               -> working tree (the interactive case)
//	c. an auto-detected base branch       -> "<merge-base>...HEAD", or the last
//	   commit when HEAD already is the base branch
//	d. nothing                            -> ReviewModeNone
//
// Case (a) is deliberately absolute: when the caller names a base, an empty
// range is "nothing to review" rather than a silent fallback to another source.
func resolveChangeSet(ctx context.Context, repoRoot string, base string) ChangeSet {
	base = strings.TrimSpace(base)
	if base != "" {
		return rangeChangeSet(ctx, repoRoot, base, base+"...HEAD")
	}
	if files := changedFiles(ctx, repoRoot); len(files) > 0 {
		return ChangeSet{Mode: ReviewModeWorkingTree, Files: files, Target: "the working tree"}
	}
	if committed, ok := committedChangeSet(ctx, repoRoot); ok {
		return committed
	}
	return ChangeSet{Mode: ReviewModeNone, Target: noCommittedChangeTarget(ctx, repoRoot)}
}

// committedChangeSet is case (c): the tree is clean, so review the branch's own
// commits. When HEAD is the base branch itself there is no branch to compare,
// and the last commit is the change worth reviewing.
func committedChangeSet(ctx context.Context, repoRoot string) (ChangeSet, bool) {
	if base, ok := detectReviewBase(ctx, repoRoot); ok {
		if set := rangeChangeSet(ctx, repoRoot, base, base+"...HEAD"); len(set.Files) > 0 {
			return set, true
		}
	}
	if refExists(ctx, repoRoot, "HEAD~1") {
		if set := rangeChangeSet(ctx, repoRoot, "HEAD~1", "HEAD~1...HEAD"); len(set.Files) > 0 {
			return set, true
		}
	}
	return ChangeSet{}, false
}

func rangeChangeSet(ctx context.Context, repoRoot, base, refRange string) ChangeSet {
	files := rangeChangedFiles(ctx, repoRoot, refRange)
	mode := ReviewModeRange
	if len(files) == 0 {
		mode = ReviewModeNone
	}
	return ChangeSet{
		Mode:   mode,
		Base:   base,
		Range:  refRange,
		Files:  files,
		Target: "`" + refRange + "`",
	}
}

// noCommittedChangeTarget names everywhere gx looked, so a "nothing to review"
// result is actionable instead of mysterious.
func noCommittedChangeTarget(ctx context.Context, repoRoot string) string {
	tried := append([]string(nil), baseCandidates...)
	tried = append(tried, "HEAD~1")
	if !refExists(ctx, repoRoot, "HEAD") {
		return "the working tree (no commits in this repository)"
	}
	return fmt.Sprintf("the working tree and a commit range against %s", strings.Join(tried, ", "))
}

// rangeChangedFiles lists the files a ref range touches. Three-dot ranges are
// resolved by git itself, so "<base>...HEAD" is the merge base to HEAD.
func rangeChangedFiles(ctx context.Context, repoRoot, refRange string) []string {
	refRange = strings.TrimSpace(refRange)
	if repoRoot == "" || refRange == "" {
		return nil
	}
	cmd := exec.CommandContext(ctx, "git", "diff", "--name-only", "--no-ext-diff", refRange)
	cmd.Dir = repoRoot
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil
	}
	var files []string
	for _, line := range strings.Split(out.String(), "\n") {
		file := strings.TrimSpace(filepath.ToSlash(line))
		if file != "" {
			files = append(files, file)
		}
	}
	sort.Strings(files)
	return files
}

// detectReviewBase mirrors detectDefaultBase in internal/capture/orchestrator,
// except it reports failure instead of guessing "main": a base that does not
// exist would turn into a review of nothing.
func detectReviewBase(ctx context.Context, repoRoot string) (string, bool) {
	for _, candidate := range baseCandidates {
		if refExists(ctx, repoRoot, candidate) {
			return candidate, true
		}
	}
	return "", false
}

func refExists(ctx context.Context, repoRoot, ref string) bool {
	if strings.TrimSpace(repoRoot) == "" || strings.TrimSpace(ref) == "" {
		return false
	}
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--verify", "--quiet", ref+"^{commit}")
	cmd.Dir = repoRoot
	return cmd.Run() == nil
}
