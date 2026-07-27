package codereview

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

// Path impact: one answer to "which of these files matters most?", shared by
// every place that has to choose.
//
// This logic was written once, inside internal/publication, to rank the hunks a
// PR summary describes. The review command had no ranking at all: it sorted the
// changed files alphabetically and read the first ten, so whether a review saw
// an auth change depended on where its path fell in the alphabet. Two different
// answers to the same question is one answer too many, so the path-based half
// of publication's hunkImpact lives here now and publication calls it. The
// scores, titles, and detail sentences are unchanged, and the size-and-catalog
// clause that needs a PR catalog stays in publication where the catalog is.

// ChangeImpact is how much a path matters and why. Score is 0 when nothing
// matched, which is the common case: most files are ordinary.
type ChangeImpact struct {
	Score  int
	Title  string
	Detail string
}

// Matched reports whether any impact signal fired for this path.
func (i ChangeImpact) Matched() bool { return i.Score > 0 }

// FileImpact scores a path against the project's configured risk paths and then
// against the generic patterns that mean "a human should look at this".
//
// The order of the checks is load-bearing and deliberately preserved from the
// PR-summary implementation: a file matching both `auth` and `schema` is
// reported as the auth change it is. Configured risk paths always win, because
// they are the one signal the project itself supplied.
func FileImpact(file string, policy ReviewPolicy) ChangeImpact {
	file = strings.TrimSpace(file)
	if file == "" {
		return ChangeImpact{}
	}
	base := path.Base(file)
	lowerFile := strings.ToLower(file + " " + base)

	for _, riskPath := range policy.RiskPaths {
		if !MatchRiskPathGlob(riskPath.Glob, file) {
			continue
		}
		title := impactRiskPathTitle(riskPath.Message, base)
		detail := fmt.Sprintf("%s (matched `%s` in %s).", strings.TrimSpace(riskPath.Message), riskPath.Glob, file)
		return ChangeImpact{Score: 700, Title: title, Detail: impactTrimSentence(detail, 240)}
	}

	if impactContainsAny(lowerFile, []string{"auth", "token", "secret", "credential"}) {
		pattern := "auth|token|secret|credential"
		title := fmt.Sprintf("Verify auth handling in %s", base)
		detail := fmt.Sprintf("Path or filename matches generic pattern %q in %s; auth changes can leak, drop, or misuse credentials.", pattern, file)
		return ChangeImpact{Score: 650, Title: title, Detail: impactTrimSentence(detail, 240)}
	}
	if impactContainsAny(lowerFile, []string{"migration", "schema"}) {
		pattern := "migration|schema"
		title := fmt.Sprintf("Verify schema change in %s", base)
		detail := fmt.Sprintf("Path matches generic pattern %q in %s; schema changes can break persistence or migrations.", pattern, file)
		return ChangeImpact{Score: 670, Title: title, Detail: impactTrimSentence(detail, 240)}
	}
	if strings.Contains(lowerFile, "dockerfile") || strings.Contains(lowerFile, ".github/workflows") || strings.Contains(lowerFile, "deploy") {
		pattern := "Dockerfile|.github/workflows|deploy"
		title := fmt.Sprintf("Verify deployment change in %s", base)
		detail := fmt.Sprintf("Path matches generic pattern %q in %s; deployment changes can alter runtime or CI behavior.", pattern, file)
		return ChangeImpact{Score: 660, Title: title, Detail: impactTrimSentence(detail, 240)}
	}
	// Faithful port, quirk included: publication passed `lowerFile` — the
	// composed "<path> <base>" string, not the path — so path.Base never yields
	// a bare lockfile name and this branch has never fired. Preserved rather
	// than fixed, because this function is what ranks PR-summary hunks and the
	// extraction is meant to move that logic, not edit it.
	//
	// It also happens to be the behavior a review wants. Scoring lockfiles at
	// 640 would rank a 500-line go.sum churn above every source file in the
	// change, which is the opposite of the point. See IsLockfilePath for the
	// predicate callers should use when they mean "is this a lockfile".
	if IsLockfilePath(lowerFile) {
		pattern := "lockfiles"
		title := fmt.Sprintf("Verify lockfile change in %s", base)
		detail := fmt.Sprintf("Filename matches generic pattern %q in %s; dependency lockfile changes can alter build or runtime versions.", pattern, file)
		return ChangeImpact{Score: 640, Title: title, Detail: impactTrimSentence(detail, 240)}
	}
	return ChangeImpact{}
}

// IsLockfilePath reports whether a path is a generated dependency lockfile.
//
// Lockfiles are the one kind of "dependency file" that is machine-written, and
// the distinction matters twice: a lockfile change is worth naming in a PR
// summary, and a lockfile is the last thing worth spending a review context
// slot on. Both callers need the same list, so there is one.
//
// The list is exactly the one internal/publication shipped, `bun.lock`'s
// absence included. Adding to it would change which hunks a PR summary calls
// out, and this extraction is meant to move that logic, not edit it. A
// `bun.lock` still sorts below source code in the context budget, because
// isDependencyFile catches it either way.
func IsLockfilePath(file string) bool {
	switch strings.ToLower(path.Base(strings.TrimSpace(file))) {
	case "go.sum", "package-lock.json", "yarn.lock", "pnpm-lock.yaml", "bun.lockb", "cargo.lock", "gemfile.lock", "poetry.lock", "composer.lock":
		return true
	default:
		return false
	}
}

// rankFilesByImpact orders paths most-review-worthy first: impact score, then
// non-test before test, then alphabetically so the order is reproducible.
//
// The alphabetical tiebreak is all the review command used to have. It is a
// fine tiebreak and a terrible ranking.
func rankFilesByImpact(files []string, policy ReviewPolicy) []string {
	type ranked struct {
		file  string
		score int
		test  bool
	}
	entries := make([]ranked, 0, len(files))
	for _, file := range files {
		entries = append(entries, ranked{file: file, score: FileImpact(file, policy).Score, test: isTestFile(file)})
	}
	sort.SliceStable(entries, func(i, j int) bool {
		left, right := entries[i], entries[j]
		if left.score != right.score {
			return left.score > right.score
		}
		if left.test != right.test {
			return right.test
		}
		return left.file < right.file
	})
	out := make([]string, 0, len(entries))
	for _, entry := range entries {
		out = append(out, entry.file)
	}
	return out
}

// reviewPolicyOf is the policy on these options, or an empty one. A nil policy
// is ordinary — it means nothing has loaded REVIEW.md yet — and ranking must
// still work without it.
func reviewPolicyOf(opts Options) ReviewPolicy {
	if opts.ReviewPolicy != nil {
		return *opts.ReviewPolicy
	}
	return ReviewPolicy{}
}

func impactRiskPathTitle(message, base string) string {
	message = strings.TrimSpace(message)
	if message == "" {
		return fmt.Sprintf("Verify change in %s", base)
	}
	if idx := strings.IndexAny(message, ".;"); idx > 0 {
		clause := strings.TrimSpace(message[:idx])
		if clause != "" {
			return "Verify " + clause
		}
	}
	return fmt.Sprintf("Verify change in %s", base)
}

func impactContainsAny(text string, needles []string) bool {
	for _, needle := range needles {
		if strings.Contains(text, needle) {
			return true
		}
	}
	return false
}

func impactTrimSentence(value string, limit int) string {
	value = strings.TrimSpace(strings.ReplaceAll(value, "\n", " "))
	if limit <= 0 || len(value) <= limit {
		return value
	}
	candidate := value[:limit]
	if idx := strings.LastIndexAny(candidate, ".!?"); idx >= 80 {
		return strings.TrimSpace(candidate[:idx+1])
	}
	if idx := strings.LastIndex(candidate, ";"); idx >= 80 {
		return strings.TrimSpace(candidate[:idx]) + "."
	}
	return strings.TrimRight(candidate, " .,;:") + "."
}
