package codereview

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// Review modes name what a review actually looked at. They are the honest
// answer to "what did tl read?", and the difference between "reviewed and
// clean" and "never looked" lives here rather than in the findings list.
const (
	// ReviewModeWorkingTree is the interactive case: uncommitted work.
	ReviewModeWorkingTree = "working-tree"
	// ReviewModeRange is a commit range, "<base>...HEAD".
	ReviewModeRange = "range"
	// ReviewModeRepo is a whole-repo review: the repository itself is the
	// subject and no diff is required. It is reached either by explicit
	// request (Options.WholeRepo) or as the fallback for a scope-, prompt-, or
	// deep-directed review that found no diff at all.
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

// wholeRepoChangeSet makes the repository the subject of the review, whatever
// resolveChangeSet found. It deliberately keeps the resolved Files, Base, and
// Range: a whole-repo review that threw away the caller's uncommitted work
// would read everything except the change that prompted the review, so the
// diff stays as the change in focus while the repo is what gets reviewed.
//
// explicitBase reports whether the caller also named a base. Both instructions
// are explicit and they disagree about the subject, so rather than silently
// picking one, WholeRepo sets the subject and the target says so.
//
// repoHasContent is the honesty guard. Asking for a whole-repo review of a
// repository with no files and no diff must not manufacture a review: that is
// still "never looked", and claiming otherwise would let --repo turn the one
// non-passing outcome into a pass.
func wholeRepoChangeSet(changes ChangeSet, explicitBase bool, repoHasContent bool) ChangeSet {
	if !repoHasContent && !changes.Reviewed() {
		return changes
	}
	changes.Mode = ReviewModeRepo
	changes.Target = wholeRepoTarget(changes, explicitBase)
	return changes
}

// wholeRepoTarget names the repository as the subject and, when a diff was
// resolved, names that diff too, so the report never implies tl ignored it.
func wholeRepoTarget(changes ChangeSet, explicitBase bool) string {
	focus := ""
	switch {
	case strings.TrimSpace(changes.Range) != "":
		focus = "`" + strings.TrimSpace(changes.Range) + "`"
	case len(changes.Files) > 0:
		focus = "the working tree"
	}
	if explicitBase {
		if focus == "" {
			focus = "`" + strings.TrimSpace(changes.Base) + "`"
		}
		return fmt.Sprintf("the repository (whole-repo review requested, so it set the subject over --base; %s kept as the change in focus)", focus)
	}
	if focus == "" {
		return "the repository"
	}
	return fmt.Sprintf("the repository (with %s as the change in focus)", focus)
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

// noCommittedChangeTarget names everywhere tl looked, so a "nothing to review"
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
	cmd := gitCommand(ctx, repoRoot, "diff", "--name-only", "--no-ext-diff", refRange)
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
	cmd := gitCommand(ctx, repoRoot, "rev-parse", "--verify", "--quiet", ref+"^{commit}")
	return cmd.Run() == nil
}
