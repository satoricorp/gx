package codereview

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// dirtyRepoWithUncommittedWork reproduces the reported case: real work sitting
// in the working tree. That state makes the diff the review subject, which is
// what used to put a whole-repo review out of reach entirely.
func dirtyRepoWithUncommittedWork(t *testing.T) string {
	t.Helper()
	root := initRepoOnMain(t)
	writeFile(t, root, "internal/app/app.go", "package app\n\nfunc Run() { println(1) }\n")
	return root
}

// repoWideProblemRepo adds a committed, repo-wide problem (a package.json with
// no lockfile) that has nothing to do with the uncommitted edit. It is the
// difference a whole-repo review is supposed to find and a patch-focused
// review is supposed to ignore.
func repoWideProblemRepo(t *testing.T) string {
	t.Helper()
	root := dirtyRepoWithUncommittedWork(t)
	writeFile(t, root, "package.json", "{\n  \"name\": \"example\"\n}\n")
	gitAdd(t, root, "package.json")
	gitCommitMessage(t, root, "add package.json")
	return root
}

// resolveReviewSubject mirrors what the engine does: scan the repository
// first, then resolve the subject, so a test sees the same inputs the review
// itself has. A whole-repo review depends on the repository having something
// in it, and the scan is where that is known.
func resolveReviewSubject(t *testing.T, root string, opts Options) ChangeSet {
	t.Helper()
	ctx := context.Background()
	facts, err := LocalScanner{}.Scan(ctx, root, strings.TrimSpace(opts.Focus))
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	return reviewChangeSet(ctx, root, normalizeOptions(opts), facts.TrackedFileCount)
}

// briefForOptions runs a real review and returns exactly what the model was
// handed. Every claim about what a whole-repo review "reads" has to be made
// against this, not against the report's mode label.
func briefForOptions(t *testing.T, root string, opts Options) (ReviewBrief, Report) {
	t.Helper()
	reviewer := &capturingReviewer{}
	engine := NewEngineWithReviewer(LocalScanner{}, StaticCatalog{}, defaultRules(), LocalContextRetriever{}, reviewer)
	report, err := engine.Review(context.Background(), root, opts)
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	return reviewer.brief, report
}

func contextSnippetsOfKind(brief ReviewBrief, kind string) []ContextSnippet {
	var out []ContextSnippet
	for _, snippet := range brief.Context {
		if snippet.Kind == kind {
			out = append(out, snippet)
		}
	}
	return out
}

func briefMentionsSource(brief ReviewBrief, needle string) bool {
	for _, snippet := range brief.Context {
		if strings.Contains(snippet.Text, needle) {
			return true
		}
	}
	for _, snippet := range brief.Static.DiffSnippets {
		if strings.Contains(snippet.Diff, needle) {
			return true
		}
	}
	return false
}

// WholeRepo names the subject, so it holds in every working-tree state. The
// plain column of this table is the unchanged behavior it must not disturb.
func TestWholeRepoOptionSetsTheSubjectInEveryWorkingTreeState(t *testing.T) {
	cases := []struct {
		name      string
		setup     func(t *testing.T) string
		plainMode string
	}{
		{
			// The reported case.
			name:      "dirty working tree",
			setup:     dirtyRepoWithUncommittedWork,
			plainMode: ReviewModeWorkingTree,
		},
		{
			name: "clean tree with branch commits",
			setup: func(t *testing.T) string {
				root := initRepoOnMain(t)
				commitOnFeatureBranch(t, root)
				return root
			},
			plainMode: ReviewModeRange,
		},
		{
			name:      "fresh repo with a single commit",
			setup:     initRepoOnMain,
			plainMode: ReviewModeNone,
		},
		{
			name: "detached HEAD",
			setup: func(t *testing.T) string {
				root := initRepoOnMain(t)
				commitOnFeatureBranch(t, root)
				runGit(t, root, "checkout", "--detach", "HEAD")
				return root
			},
			plainMode: ReviewModeRange,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := tc.setup(t)

			plain := resolveReviewSubject(t, root, Options{})
			if plain.Mode != tc.plainMode {
				t.Fatalf("plain review Mode = %q, want %q (unchanged behavior)", plain.Mode, tc.plainMode)
			}

			whole := resolveReviewSubject(t, root, Options{WholeRepo: true})
			if whole.Mode != ReviewModeRepo {
				t.Fatalf("WholeRepo Mode = %q, want %q", whole.Mode, ReviewModeRepo)
			}
			if !whole.Reviewed() {
				t.Fatalf("WholeRepo Reviewed() = false, want a whole-repo review to count as reviewed")
			}
			if !strings.Contains(whole.Target, "the repository") {
				t.Fatalf("WholeRepo Target = %q, want it to name the repository", whole.Target)
			}
		})
	}
}

// The exact case the repo owner hit: uncommitted work made the diff the only
// reachable subject, so a review of the codebase was impossible and the report
// said "no material issues in this change" about eight files.
func TestWholeRepoReviewOnADirtyTreeReviewsTheRepository(t *testing.T) {
	root := repoWideProblemRepo(t)

	brief, report := briefForOptions(t, root, Options{WholeRepo: true})

	if !report.Reviewed || report.ReviewMode != ReviewModeRepo {
		t.Fatalf("Reviewed/ReviewMode = %v/%q, want true/%q", report.Reviewed, report.ReviewMode, ReviewModeRepo)
	}
	if !strings.Contains(report.ReviewTarget, "the repository") {
		t.Fatalf("ReviewTarget = %q, want it to name the repository", report.ReviewTarget)
	}
	// The uncommitted work must stay visible: a whole-repo review that dropped
	// the caller's edits would read everything except the reason they asked.
	if !strings.Contains(report.ReviewTarget, "working tree") {
		t.Fatalf("ReviewTarget = %q, want the working tree named as the change in focus", report.ReviewTarget)
	}
	if len(report.ChangedFiles) != 1 || report.ChangedFiles[0] != "internal/app/app.go" {
		t.Fatalf("ChangedFiles = %#v, want the uncommitted file to remain in focus", report.ChangedFiles)
	}
	if len(brief.Static.DiffSnippets) != 1 || !strings.Contains(brief.Static.DiffSnippets[0].Diff, "println(1)") {
		t.Fatalf("brief diff snippets = %#v, want the uncommitted diff", brief.Static.DiffSnippets)
	}
	if brief.ReviewProfile != reviewProfileWholeRepo {
		t.Fatalf("brief review profile = %q, want %q", brief.ReviewProfile, reviewProfileWholeRepo)
	}
	// The substantive difference: a finding that lives outside the diff
	// survives, because the patch-focused filter is off.
	if !hasFinding(report.Findings, "dependencies.javascript-unlocked") {
		t.Fatalf("Findings = %#v, want the repo-wide dependency finding", report.Findings)
	}

	text := RenderMarkdown(report)
	if !strings.Contains(text, "Reviewed the repository") {
		t.Fatalf("RenderMarkdown() does not say the repository was reviewed:\n%s", text)
	}
	if strings.Contains(text, "Nothing to review") {
		t.Fatalf("RenderMarkdown() reported nothing to review for a whole-repo review:\n%s", text)
	}
}

// The claim --repo has to earn: the reviewer is shown the repository. A mode
// label is not a review. Before this, the brief a whole-repo review produced
// was identical in content to the brief for `--scope architecture`, so the
// only thing "whole-repo" about it was the word in the report.
func TestWholeRepoReviewShowsTheReviewerTheRepository(t *testing.T) {
	root := dirtyRepoWithUncommittedWork(t)
	// Committed code that no diff touches and no hint points at: the model has
	// no way to see this except by being given the repository.
	writeFile(t, root, "internal/queue/queue.go", "package queue\n\nfunc Drain() string { return \"unreviewed-repo-code\" }\n")
	gitAdd(t, root, "internal/queue/queue.go")
	gitCommitMessage(t, root, "add queue")

	whole, _ := briefForOptions(t, root, Options{WholeRepo: true})
	if len(contextSnippetsOfKind(whole, "repo_inventory")) != 1 {
		t.Fatalf("repo_inventory snippets = %d, want 1: a whole-repo review has to say what the repository contains", len(contextSnippetsOfKind(whole, "repo_inventory")))
	}
	inventory := contextSnippetsOfKind(whole, "repo_inventory")[0]
	for _, want := range []string{"internal/queue/queue.go", "internal/app/app.go", "tracked files"} {
		if !strings.Contains(inventory.Text, want) {
			t.Fatalf("repo inventory does not list %q:\n%s", want, inventory.Text)
		}
	}
	if !briefMentionsSource(whole, "unreviewed-repo-code") {
		t.Fatalf("brief never shows the model the repository's own code:\n%#v", whole.Context)
	}

	// And the review it used to be indistinguishable from does not do this.
	scoped, _ := briefForOptions(t, root, Options{Scope: "architecture"})
	if len(contextSnippetsOfKind(scoped, "repo_inventory")) != 0 || len(contextSnippetsOfKind(scoped, "repo_source_file")) != 0 {
		t.Fatalf("a scope-directed review gained whole-repo context; --repo must be the thing that reads the repository")
	}
	if scoped.ReviewProfile == whole.ReviewProfile {
		t.Fatalf("review profiles are identical (%q); the two reviews are indistinguishable to the model", scoped.ReviewProfile)
	}
	if whole.Rubric.Goal == scoped.Rubric.Goal {
		t.Fatalf("whole-repo rubric goal is the scoped one:\n%s", whole.Rubric.Goal)
	}
}

// Repository reading has to be language-agnostic. Module summaries, hints, and
// module_file snippets are Go-only, so in a TypeScript repo they contribute
// nothing: without this, an agent asking about the queue code in a TS repo got
// an answer computed without ever seeing the queue code.
func TestWholeRepoReviewReadsANonGoRepository(t *testing.T) {
	root := initRepo(t)
	writeFile(t, root, "package.json", "{\n  \"name\": \"example\"\n}\n")
	writeFile(t, root, "src/queue.ts", "export function drain(): string {\n  return \"typescript-queue-body\";\n}\n")
	writeFile(t, root, "src/auth.ts", "export const token = process.env.TOKEN;\n")
	writeFile(t, root, "src/db.ts", "export const url = process.env.DATABASE_URL;\n")
	gitAdd(t, root, "package.json", "src/queue.ts", "src/auth.ts", "src/db.ts")
	gitCommitMessage(t, root, "initial")
	runGit(t, root, "branch", "-M", "main")
	writeFile(t, root, "src/edit.ts", "export const edited = true;\n")

	brief, report := briefForOptions(t, root, Options{
		WholeRepo: true,
		Prompt:    "is there a race in the queue code?",
	})
	if report.ReviewMode != ReviewModeRepo {
		t.Fatalf("ReviewMode = %q, want %q", report.ReviewMode, ReviewModeRepo)
	}
	if len(brief.Static.Modules) != 0 {
		t.Fatalf("Modules = %#v; this repo has no Go, so module summaries cannot be what carries the repository", brief.Static.Modules)
	}
	if !briefMentionsSource(brief, "typescript-queue-body") {
		t.Fatalf("the queue code the prompt asks about never reached the model:\n%#v", brief.Context)
	}
	// The prompt is the only local signal about what matters, so it decides
	// which files are worth the budget.
	sources := contextSnippetsOfKind(brief, "repo_source_file")
	if len(sources) == 0 || !strings.Contains(sources[0].Ref, "queue") {
		t.Fatalf("first repo source snippet = %#v, want the file the prompt asks about", sources)
	}
}

// The same repo without the flag: unchanged: the diff is the subject and the
// patch-focused filter drops everything that is not in it.
func TestPatchFocusedReviewOnADirtyTreeIsUnchanged(t *testing.T) {
	root := repoWideProblemRepo(t)

	report, err := Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if !report.Reviewed || report.ReviewMode != ReviewModeWorkingTree {
		t.Fatalf("Reviewed/ReviewMode = %v/%q, want true/%q", report.Reviewed, report.ReviewMode, ReviewModeWorkingTree)
	}
	if len(report.ChangedFiles) != 1 || report.ChangedFiles[0] != "internal/app/app.go" {
		t.Fatalf("ChangedFiles = %#v, want the uncommitted file", report.ChangedFiles)
	}
	if hasFinding(report.Findings, "dependencies.javascript-unlocked") {
		t.Fatalf("Findings = %#v, want the patch-focused filter to drop findings outside the diff", report.Findings)
	}
	if text := RenderMarkdown(report); strings.Contains(text, "Reviewed the repository") {
		t.Fatalf("RenderMarkdown() claimed a whole-repo review for a diff review:\n%s", text)
	}
}

// Two explicit instructions disagree about the subject. WholeRepo sets it, the
// base still picks the diff, and the target says so rather than leaving the
// caller to guess which one gx obeyed.
func TestWholeRepoOptionSetsTheSubjectOverAnExplicitBase(t *testing.T) {
	root := initRepoOnMain(t)
	commitOnFeatureBranch(t, root)
	// A dirty file neither instruction selected.
	writeFile(t, root, "internal/app/scratch.go", "package app\n")

	changes := resolveReviewSubject(t, root, Options{Base: "main", WholeRepo: true})
	if changes.Mode != ReviewModeRepo {
		t.Fatalf("Mode = %q, want %q: --repo names the subject", changes.Mode, ReviewModeRepo)
	}
	if changes.Base != "main" || changes.Range != "main...HEAD" {
		t.Fatalf("Base/Range = %q/%q, want main/main...HEAD kept as the change in focus", changes.Base, changes.Range)
	}
	if len(changes.Files) != 1 || changes.Files[0] != "internal/app/feature.go" {
		t.Fatalf("Files = %#v, want the base range's files, not the dirty tree", changes.Files)
	}
	for _, want := range []string{"the repository", "--base", "main...HEAD"} {
		if !strings.Contains(changes.Target, want) {
			t.Fatalf("Target = %q, want it to mention %q so the precedence is visible", changes.Target, want)
		}
	}

	// And --base on its own is untouched.
	baseOnly := resolveReviewSubject(t, root, Options{Base: "main"})
	if baseOnly.Mode != ReviewModeRange {
		t.Fatalf("--base alone Mode = %q, want %q (unchanged behavior)", baseOnly.Mode, ReviewModeRange)
	}
}

// "Never looked" stays its own outcome, distinct from a clean review, and it
// stays impossible to pass a gate with — with or without --repo. Asking for a
// whole-repo review of a repository that has nothing in it must not
// manufacture one: that would make the one non-passing outcome unreachable.
func TestEmptyRepoIsStillNothingToReviewAndNeverAPass(t *testing.T) {
	for _, tc := range []struct {
		name string
		opts Options
	}{
		{name: "plain", opts: Options{}},
		{name: "whole repo", opts: Options{WholeRepo: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := initRepo(t)

			report, err := Review(context.Background(), root, tc.opts)
			if err != nil {
				t.Fatalf("Review() error = %v", err)
			}
			if report.Reviewed || report.ReviewMode != ReviewModeNone {
				t.Fatalf("Reviewed/ReviewMode = %v/%q, want false/%q", report.Reviewed, report.ReviewMode, ReviewModeNone)
			}
			for _, level := range []FailOnLevel{FailOnAny, FailOnSpeculative, FailOnBlocking} {
				if failures := report.GateFailures(level); failures != nil {
					t.Fatalf("GateFailures(%q) = %#v, want nil so callers must handle ReviewModeNone themselves", level, failures)
				}
			}
			text := RenderMarkdown(report)
			if !strings.Contains(text, "Nothing to review") || !strings.Contains(text, "not a clean review") {
				t.Fatalf("RenderMarkdown() missing the nothing-to-review outcome:\n%s", text)
			}
			if strings.Contains(text, "No material issues") {
				t.Fatalf("RenderMarkdown() reported a clean review for a repo it never read:\n%s", text)
			}
		})
	}
}

// Decision, pinned so it cannot drift: --deep, --scope, and --prompt choose a
// lens and a depth, not a subject. While a diff exists, that diff stays the
// subject; the repository becomes the subject only on explicit request, or
// when there is no diff anywhere to be found. The doc comment on
// reviewChangeSet used to claim otherwise; the behavior below is the contract.
func TestDeepScopeAndPromptDoNotOverrideADirtyWorkingTree(t *testing.T) {
	cases := []struct {
		name string
		opts Options
	}{
		{name: "deep", opts: Options{Deep: true}},
		{name: "scope", opts: Options{Scope: "security"}},
		{name: "prompt", opts: Options{Prompt: "where does auth read its secret from?"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := dirtyRepoWithUncommittedWork(t)

			changes := resolveReviewSubject(t, root, tc.opts)
			if changes.Mode != ReviewModeWorkingTree {
				t.Fatalf("Mode = %q, want %q: a diff in hand stays the subject", changes.Mode, ReviewModeWorkingTree)
			}

			// ...and each of them reaches the whole repo when asked to.
			withRepo := tc.opts
			withRepo.WholeRepo = true
			whole := resolveReviewSubject(t, root, withRepo)
			if whole.Mode != ReviewModeRepo {
				t.Fatalf("Mode with WholeRepo = %q, want %q", whole.Mode, ReviewModeRepo)
			}
		})
	}
}

// The other half of the same decision: with no diff at all, a scope-, prompt-,
// or deep-directed review still falls back to the repository rather than
// reporting "nothing to review". This is pre-existing behavior and must stay.
func TestDirectedReviewWithoutAnyDiffStillFallsBackToTheRepository(t *testing.T) {
	root := initRepoOnMain(t)

	for _, tc := range []struct {
		name string
		opts Options
	}{
		{name: "deep", opts: Options{Deep: true}},
		{name: "scope", opts: Options{Scope: "security"}},
		{name: "prompt", opts: Options{Prompt: "what runs at startup?"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			changes := resolveReviewSubject(t, root, tc.opts)
			if changes.Mode != ReviewModeRepo {
				t.Fatalf("Mode = %q, want %q", changes.Mode, ReviewModeRepo)
			}
		})
	}

	// The patch-focused default keeps reporting the truth instead.
	if changes := resolveReviewSubject(t, root, Options{}); changes.Mode != ReviewModeNone {
		t.Fatalf("plain review Mode = %q, want %q", changes.Mode, ReviewModeNone)
	}
}

// That fallback path reaches repo mode without anyone passing --repo, and the
// report it produces is published: `gx review` upserts it as the PR comment
// and records it as the gx Cloud history summary. So the wording is pinned
// here too, not only for the flag, and the two must agree.
func TestRepoModeNamesItsSubjectWithoutTheFlagToo(t *testing.T) {
	root := initRepoOnMain(t)

	report, err := Review(context.Background(), root, Options{Deep: true})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.ReviewMode != ReviewModeRepo {
		t.Fatalf("ReviewMode = %q, want %q", report.ReviewMode, ReviewModeRepo)
	}
	text := RenderMarkdown(report)
	if !strings.Contains(text, "Reviewed the repository") {
		t.Fatalf("a repo-mode review does not say what it reviewed:\n%s", text)
	}
	if strings.Contains(text, "No material issues found in this change") {
		t.Fatalf("a repo-mode review reported on \"this change\":\n%s", text)
	}
}

func TestWholeRepoOptionTurnsOffThePatchFocusedFilter(t *testing.T) {
	if opts := normalizeOptions(Options{}); !opts.PatchFocused {
		t.Fatal("normalizeOptions(Options{}).PatchFocused = false, want the patch-focused default intact")
	}
	if opts := normalizeOptions(Options{WholeRepo: true}); opts.PatchFocused {
		t.Fatal("normalizeOptions(WholeRepo).PatchFocused = true, want the diff-only filter off")
	}
	// Even if a caller asks for both, the contradiction resolves toward the
	// broader subject rather than silently filtering the review down again.
	if opts := normalizeOptions(Options{WholeRepo: true, PatchFocused: true}); opts.PatchFocused {
		t.Fatal("normalizeOptions(WholeRepo+PatchFocused).PatchFocused = true, want WholeRepo to win")
	}
	if profile := reviewProfile(normalizeOptions(Options{WholeRepo: true})); profile != reviewProfileWholeRepo {
		t.Fatalf("reviewProfile() = %q, want %q", profile, reviewProfileWholeRepo)
	}
}

// Asking for a broader review must never buy a narrower one. Clearing
// PatchFocused swaps the scope list for `--scope` plus baselines, and that
// list has no `security` in it — the scopes gate rules, findings, and the
// model's source catalog, and the risk tags drive resource retrieval.
func TestWholeRepoReviewKeepsEveryLensThePatchReviewHas(t *testing.T) {
	patch := normalizeOptions(Options{})
	whole := normalizeOptions(Options{WholeRepo: true})

	patchScopes := activeScopeList(patch)
	wholeScopes := activeScopeList(whole)
	for _, scope := range patchScopes {
		if !containsString(wholeScopes, scope) {
			t.Fatalf("whole-repo scopes %v dropped %q, which the patch-focused default (%v) covers", wholeScopes, scope, patchScopes)
		}
	}
	if !containsString(wholeScopes, whole.Scope) {
		t.Fatalf("whole-repo scopes %v do not include the requested scope %q", wholeScopes, whole.Scope)
	}

	files := []string{"internal/app/app.go"}
	patchTags := riskTagsForReview(files, nil, patch)
	wholeTags := riskTagsForReview(files, nil, whole)
	for _, tag := range patchTags {
		if !containsString(wholeTags, tag) {
			t.Fatalf("whole-repo risk tags %v dropped %q, which the patch-focused default (%v) carries", wholeTags, tag, patchTags)
		}
	}

	// The scope list still narrows when a scope is explicitly requested
	// without --repo: that behavior is unchanged.
	if scopes := activeScopeList(normalizeOptions(Options{Scope: "architecture"})); containsString(scopes, "security") {
		t.Fatalf("scope-directed scopes = %v, want the pre-existing narrowing left alone", scopes)
	}
}

// A whole-repo review with no diff still has to run the tools. Every runner
// detects against the change set, and collectStaticToolResults returns before
// detection when that set is empty — so without the repository standing in,
// `gx review --repo --fail-on any` on a clean tree reports the repository
// clean having compiled nothing.
func TestStaticToolsRunOverTheRepositoryWhenAWholeRepoReviewHasNoDiff(t *testing.T) {
	root := initRepoOnMain(t)
	facts, err := LocalScanner{}.Scan(context.Background(), root, "")
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if scope := staticToolScope(facts, normalizeOptions(Options{}), nil); len(scope) != 0 {
		t.Fatalf("patch-focused scope with no changes = %#v, want empty (unchanged behavior)", scope)
	}
	scope := staticToolScope(facts, normalizeOptions(Options{WholeRepo: true}), nil)
	if !containsString(scope, "internal/app/app.go") {
		t.Fatalf("whole-repo tool scope = %#v, want the repository's source files", scope)
	}
	// The change set still wins when there is one: a whole-repo review of a
	// dirty tree scopes its tools to the work in progress, as before.
	changed := []string{"internal/app/app.go"}
	if got := staticToolScope(facts, normalizeOptions(Options{WholeRepo: true}), changed); len(got) != 1 || got[0] != changed[0] {
		t.Fatalf("whole-repo tool scope with a diff = %#v, want the changed files", got)
	}
	if _, ok := goPackageArgsForChangedFiles(root, scope); !ok {
		t.Fatalf("go tools would not run over the repository scope %#v", scope)
	}
}

// Every profile token the brief can carry has to be defined in the
// instructions sent with it. A model told "use review_profile to choose
// behavior", handed an undefined token, and told to return an empty
// recommendations array when nothing meets the active profile is being set up
// to return nothing — which is the symptom this whole change started from.
func TestEveryReviewProfileIsDefinedInThePromptSentWithIt(t *testing.T) {
	profiles := []string{
		reviewProfile(normalizeOptions(Options{})),
		reviewProfile(normalizeOptions(Options{Deep: true})),
		reviewProfile(normalizeOptions(Options{Scope: "security"})),
		reviewProfile(normalizeOptions(Options{Prompt: "why?"})),
		reviewProfile(normalizeOptions(Options{WholeRepo: true})),
		// The PR summary pipeline sets this one by hand.
		"pr_summary",
	}
	for _, profile := range profiles {
		prompt := reviewDeveloperPrompt(ReviewBrief{ReviewProfile: profile})
		if !strings.Contains(prompt, profile) {
			t.Fatalf("review prompt never defines profile %q:\n%s", profile, prompt)
		}
	}

	wholeRepoPrompt := reviewDeveloperPrompt(ReviewBrief{ReviewProfile: reviewProfileWholeRepo})
	for _, want := range []string{"repo_inventory", "repo_source_file", "the repository itself is the subject"} {
		if !strings.Contains(wholeRepoPrompt, want) {
			t.Fatalf("whole-repo prompt missing %q:\n%s", want, wholeRepoPrompt)
		}
	}
	// The prompt also has a standing "prefer findings tied to changed lines
	// over repo-wide advice" instruction. A whole-repo review has to be told
	// that does not apply to it, or it will decline the very findings it was
	// asked for.
	if !strings.Contains(wholeRepoPrompt, "does not apply here") {
		t.Fatalf("whole-repo prompt does not override the changed-lines preference:\n%s", wholeRepoPrompt)
	}
}

// The brief is compacted before it is sent, and compaction is where a
// whole-repo review can quietly stop being one: the generic per-snippet budget
// would cut the file listing to a fraction of the repository, and the default
// snippet count would drop the repository's code entirely.
func TestTheRepositoryReachesTheModelAfterCompaction(t *testing.T) {
	root := initRepo(t)
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	for i := 0; i < 120; i++ {
		writeFile(t, root, fmt.Sprintf("internal/pkg%02d/file.go", i), fmt.Sprintf("package pkg%02d\n\nfunc F%02d() {}\n", i, i))
	}
	gitAdd(t, root, ".")
	gitCommitMessage(t, root, "initial")
	runGit(t, root, "branch", "-M", "main")

	brief, _ := briefForOptions(t, root, Options{WholeRepo: true})
	compacted := compactReviewBriefForAI(brief)

	inventories := contextSnippetsOfKind(compacted, "repo_inventory")
	if len(inventories) != 1 {
		t.Fatalf("repo_inventory snippets after compaction = %d, want 1", len(inventories))
	}
	// The inventory has its own budget, larger than the generic per-snippet one.
	// This used to be asserted by "the surviving text is longer than the generic
	// budget", which was a proxy that held only while that budget was 1200 bytes
	// and a 121-file listing was 2977. Raising the generic budget to 16000 broke
	// the proxy without touching the property, so the property is asserted
	// directly now: the inventory is exempt, and nothing was cut.
	if contextSnippetByteLimit(inventories[0]) <= maxAIContextSnippetBytes {
		t.Fatalf("repo inventory shares the generic snippet budget (%d bytes); a large repository's listing would be cut", contextSnippetByteLimit(inventories[0]))
	}
	inventory := inventories[0].Text
	if strings.Contains(inventory, "[truncated]") || strings.Contains(inventory, "more files not listed") {
		t.Fatalf("repo inventory was truncated; the model sees a fraction of the repository:\n%s", inventory)
	}
	for _, want := range []string{"internal/pkg00/file.go", "internal/pkg99/file.go", "internal/pkg119/file.go"} {
		if !strings.Contains(inventory, want) {
			t.Fatalf("repo inventory lost %s, so it is not an inventory of this repository:\n%s", want, inventory)
		}
	}
	if len(contextSnippetsOfKind(compacted, "repo_source_file")) == 0 {
		t.Fatalf("compaction dropped every repository source file; %d context snippets survived", len(compacted.Context))
	}
}

// The prompt is shared with the PR summary pipeline, which runs on every push.
// Profile-specific instructions are appended for the profile in hand precisely
// so adding one cannot move what gx writes on a PR.
func TestAddingTheWholeRepoProfileLeavesEveryOtherPromptUnchanged(t *testing.T) {
	base := strings.Join(baseReviewDeveloperPromptLines(), "\n")
	for _, profile := range []string{"", "patch_focused", "pr_summary", "prompt_directed", "scope_focused", "deep_full_spectrum"} {
		prompt := reviewDeveloperPrompt(ReviewBrief{ReviewProfile: profile})
		if prompt != base {
			t.Fatalf("prompt for profile %q is no longer the shared instruction set", profile)
		}
		if strings.Contains(prompt, reviewProfileWholeRepo) {
			t.Fatalf("prompt for profile %q leaked whole-repo instructions:\n%s", profile, prompt)
		}
	}
}
