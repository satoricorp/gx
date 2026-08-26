package codereview

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/satoricorp/gx/internal/termstyle"
)

const (
	DefaultScope    = "architecture"
	reviewMintANSI  = "\x1b[38;2;61;220;151m"
	reviewResetANSI = "\x1b[0m"
)

var supportedScopes = map[string]struct{}{
	"architecture":    {},
	"security":        {},
	"performance":     {},
	"onboarding":      {},
	"docs":            {},
	"dependencies":    {},
	"testing":         {},
	"maintainability": {},
}

var baselineScopes = []string{"dependencies", "testing", "maintainability"}

type Options struct {
	Scope string
	Deep  bool
	Focus string
	// Base names the ref to review against. When set, the review covers
	// "<Base>...HEAD" (three-dot: what HEAD added since the merge base) and
	// ignores the working tree.
	Base string
	// WholeRepo makes the repository itself the subject of the review, whatever
	// the working tree looks like. It is the explicit answer to "review the
	// codebase, not just my diff": without it a dirty tree always wins and a
	// whole-repo review is unreachable. Any diff that was resolved is still
	// read, so the work in progress informs the review instead of being
	// discarded.
	WholeRepo    bool
	Prompt       string
	Verbose      bool
	PatchFocused bool
	ReviewPolicy *ReviewPolicy
	// MaxFindings caps how many recommendations a review reports. Zero means
	// defaultMaxFindings.
	//
	// It exists because the cap was previously three separate hidden numbers
	// that did not agree: the prompt told the model "at most 5, prefer 2-3",
	// the unjudged fallback path truncated to 3 in code, and neither was
	// reachable by a caller. A repository with six real defects could not
	// report six no matter what, and the output looked like the reviewer's
	// judgement rather than a ceiling nobody could see.
	MaxFindings int

	// Fast trades panel breadth for wall clock: one reviewer leg instead of
	// two, no verification pass, and findings written tightly rather than with
	// worked code examples.
	//
	// It exists because an interactive review and a pull-request summary are
	// different products. Nobody waits on a PR summary, so it should spend
	// whatever it takes to be thorough; a developer who typed a command and is
	// watching a spinner is paying for every second, and measurement says the
	// second leg contributes far less than it costs on a single change (both
	// legs found the same defects on the planted-bug fixture, and the second
	// leg's slowest call sets the shard's wall clock).
	Fast bool

	ProgressWriter io.Writer
	Color          bool
}

type Report struct {
	RepoRoot          string         `json:"repo_root,omitempty"`
	Scope             string         `json:"scope,omitempty"`
	Deep              bool           `json:"deep"`
	Focus             string         `json:"focus,omitempty"`
	Prompt            string         `json:"prompt,omitempty"`
	BaselineScopes    []string       `json:"baseline_scopes,omitempty"`
	Docs              []FilePresence `json:"docs,omitempty"`
	DependencyFiles   []string       `json:"dependency_files,omitempty"`
	TestFileCount     int            `json:"test_file_count"`
	TrackedFileCount  int            `json:"tracked_file_count"`
	ChangedFiles      []string       `json:"changed_files,omitempty"`
	ObservationLabels []string       `json:"observation_labels,omitempty"`
	Findings          []Finding      `json:"findings"`
	Story             []StoryItem    `json:"story,omitempty"` // "what you now own"; see story.go
	Sources           []Source       `json:"sources,omitempty"`
	SourceRefs        []SourceRef    `json:"source_refs,omitempty"`
	Reviewer          string         `json:"reviewer,omitempty"`
	// ReviewModels and ReviewTransport are which models answered and over which
	// wire. Both questions come up the moment a review is slow or wrong, and
	// neither is recoverable after the fact: the same panel run against local
	// AWS credentials and run through gx Cloud has different latency, different
	// quota behavior, and completely different failure modes.
	ReviewModels      []string     `json:"review_models,omitempty"`
	ReviewTransport   string       `json:"review_transport,omitempty"`
	ContextSnippets   int          `json:"context_snippets"`
	Verbose           bool         `json:"-"`
	Color             bool         `json:"-"`
	Triage            ChangeTriage `json:"triage,omitempty"`
	NoFindingsMessage string       `json:"no_findings_message,omitempty"`
	DegradedReasons   []string     `json:"degraded_reasons,omitempty"`
	// RuleLabeling is why some findings carry no rule name. It is deliberately
	// NOT DegradedReasons: that field gates — gateDegradedReason refuses to pass
	// a review that cannot answer the gate's question — and a finding without a
	// rule name answers it exactly as well as one with. The findings, their
	// severities, and their lanes are identical either way; only the label is
	// missing. Folding this into DegradedReasons turned a failed labeling call
	// into a failed gate, which is the one thing this pass must never do.
	//
	// It is still reported, because "no rule fit this finding" and "the labeler
	// never ran" render identically as a blank and mean opposite things.
	RuleLabeling []string `json:"rule_labeling,omitempty"`
	// DroppedClaims names findings removed because the repository disproved
	// them — a symbol reported absent that the file plainly contains. Like
	// RuleLabeling this is NOT DegradedReasons: a review that caught its own
	// refutable claim did its job, and gating on it would fail a gate for
	// working correctly.
	//
	// It is reported because a filter that silently eats findings cannot be
	// trusted or debugged. If this ever drops something real, this line is how
	// anyone finds out.
	DroppedClaims []string `json:"dropped_claims,omitempty"`
	// Evidence is what each retrieval source contributed, including the ones
	// that contributed nothing because they were missing or unreachable. It is
	// separate from DegradedReasons, which is specifically about the AI
	// reviewer: a review can have a working reviewer and no evidence, or the
	// reverse, and the reader needs to be able to tell which.
	Evidence []EvidenceStatus `json:"evidence,omitempty"`
	// Coverage is how much of the subject reached the model. Evidence is about
	// the supporting material; this is about the subject itself, and it is what
	// decides whether "no material issues found in the repository" is a
	// sentence this review is entitled to.
	Coverage Coverage `json:"coverage,omitempty"`
	// Tools is what the repository's own checkers said — the compiler, the
	// test runner, the vulnerability scanners. They already run on every
	// review and already feed the brief; a failure becomes a finding, but a
	// pass has never been reported anywhere. A reader deciding how much to
	// trust a clean review needs to know a test suite ran and was green.
	Tools []StaticToolResult `json:"tools,omitempty"`

	// Reviewed and the fields below answer "what did this review actually
	// read?". Reviewed is false only when nothing was inspected, which is a
	// distinct outcome from a clean review and must never be reported as one.
	Reviewed    bool   `json:"reviewed"`
	ReviewMode  string `json:"review_mode,omitempty"`
	ReviewBase  string `json:"review_base,omitempty"`
	ReviewRange string `json:"review_range,omitempty"`
	// ReviewTarget is the human phrase naming what was inspected, or where gx
	// looked when it found nothing.
	ReviewTarget string `json:"review_target,omitempty"`
}

type FilePresence struct {
	Path    string `json:"path"`
	Present bool   `json:"present"`
}

func ValidateOptions(opts Options) error {
	opts = normalizeOptions(opts)
	if _, ok := supportedScopes[opts.Scope]; !ok {
		return fmt.Errorf("unsupported review scope %q", opts.Scope)
	}
	return nil
}

// aiReviewRan reports whether a model actually reviewed this change. The engine
// sets Reviewer to "heuristic+ai" exactly when a reviewer ran without error, so
// this reads the one flag that already answers the question rather than
// inferring it from the finding list, which cannot distinguish a model finding
// from a rule finding — nor a healthy review from one that never ran.
func (r Report) aiReviewRan() bool {
	return strings.Contains(r.Reviewer, "ai")
}

func RenderMarkdown(report Report) string {
	var b strings.Builder
	// Degradation is not one condition. "No model reviewed this" and "one of two
	// models reviewed this" are different reviews, and printing "results are
	// from deterministic checks only" over four model-written findings — which
	// is what a partial reviewer failure used to produce — is a false statement
	// about the very thing the reader is deciding how much to trust.
	if len(report.DegradedReasons) > 0 {
		reason := strings.Join(report.DegradedReasons, "; ")
		if report.aiReviewRan() && len(report.Findings) > 0 {
			fmt.Fprintf(&b, "> Warning: the AI review ran degraded (%s); the findings below are real but this review saw less than a healthy one would.\n\n", reason)
		} else if report.aiReviewRan() {
			// A degraded review that found nothing is the case where "no issues"
			// is least trustworthy, and the sentence above presupposes findings
			// that are not there.
			fmt.Fprintf(&b, "> Warning: the AI review ran degraded (%s); it saw less than a healthy one would, so treat \"no issues found\" with less confidence.\n\n", reason)
		} else {
			fmt.Fprintf(&b, "> Warning: AI review unavailable (%s); results are from deterministic checks only.\n\n", reason)
		}
	}
	// Evidence that could not be read is stated unconditionally, not behind
	// --verbose. A review that never reached the code index looks exactly like
	// one that did unless it says so, and "no issues found" means something
	// materially weaker when half the evidence was missing.
	if warnings := EvidenceWarnings(report.Evidence); len(warnings) > 0 {
		fmt.Fprintf(&b, "> Warning: review evidence unavailable — %s. Findings are based on the change and the checkout only.\n\n", strings.Join(warnings, "; "))
	}
	// Coverage is stated whenever it is short, findings or no findings. A
	// review that read a third of a change and found two problems is not a
	// review that found two problems in the change.
	if report.Coverage.Partial() {
		fmt.Fprintf(&b, "> Warning: partial coverage — %s\n\n", report.Coverage.Statement())
	}
	if report.ReviewMode == ReviewModeNone {
		fmt.Fprintln(&b, reviewTitle(report, "## Recommendations"))
		fmt.Fprintf(&b, "- %s\n", NothingToReviewMessage(report))
		return strings.TrimRight(b.String(), "\n") + "\n"
	}
	fmt.Fprintln(&b, reviewTitle(report, "## Recommendations"))
	if report.ReviewMode == ReviewModeRepo {
		// A whole-repo review has a different subject from a diff review, and
		// the reader cannot tell them apart from the findings alone. Naming
		// the subject keeps the reporting invariant honest in both directions.
		//
		// Deliberately keyed on the mode, not on the --repo flag: repo mode is
		// also reached without the flag, when a scope-, prompt-, or
		// deep-directed review finds no diff at all. That review reads the
		// repository too, so it says so too — including in the PR comment and
		// the gx Cloud history entry this same text becomes. Both paths are
		// pinned by tests so the wording cannot drift for one and not the
		// other.
		if target := strings.TrimSpace(report.ReviewTarget); target != "" {
			fmt.Fprintf(&b, "\n_Reviewed %s._\n\n", target)
		}
	}
	if len(report.Findings) == 0 {
		message := strings.TrimSpace(report.NoFindingsMessage)
		if message == "" {
			message = noFindingsMessage(report.ReviewMode, report.Coverage)
		}
		fmt.Fprintf(&b, "- %s\n", message)
	} else {
		for index, finding := range report.Findings {
			fmt.Fprintf(&b, "%s\n", reviewTitle(report, fmt.Sprintf("### %d. %s", index+1, finding.Title)))
			fmt.Fprintln(&b)
			fmt.Fprintf(&b, "%s %s\n", reviewLabel(report, "**Why:**"), finding.Summary)
			if strings.TrimSpace(finding.Benefit) != "" {
				fmt.Fprintln(&b)
				fmt.Fprintf(&b, "%s %s\n", reviewLabel(report, "**Benefit:**"), finding.Benefit)
			}
			fmt.Fprintln(&b)
			fmt.Fprintf(&b, "%s %s\n", reviewLabel(report, "**Do next:**"), finding.Recommendation)
			// Corroboration is stated unconditionally, not behind --verbose.
			// "Both reviewers found this independently" is the strongest
			// confidence signal a two-model panel produces, and a reader who
			// cannot see it has no way to tell a corroborated finding from one
			// model's guess.
			if note := renderCorroboration(finding); note != "" {
				fmt.Fprintln(&b)
				fmt.Fprintf(&b, "%s %s\n", reviewLabel(report, "**Agreement:**"), note)
			}
			// The versions de-duplication folded away, unconditionally. The
			// copy a merge keeps is chosen on strength and specificity, but the
			// other reviewer's fix is often the one a reader wants, and a merge
			// that silently deletes it is a review withholding what it found.
			if alternates := renderMergedFindings(finding); len(alternates) > 0 {
				fmt.Fprintln(&b)
				fmt.Fprintf(&b, "%s\n", reviewLabel(report, "**Also reported as:**"))
				for _, alternate := range alternates {
					fmt.Fprintf(&b, "- %s\n", alternate)
				}
			}
			if attributions := renderFindingAttributions(report, finding); len(attributions) > 0 {
				fmt.Fprintln(&b)
				fmt.Fprintf(&b, "%s %s\n", reviewLabel(report, "**Informed by:**"), strings.Join(attributions, " · "))
			}
			if report.Verbose && len(finding.Evidence) > 0 {
				fmt.Fprintln(&b)
				fmt.Fprintln(&b, reviewLabel(report, "**Evidence:**"))
				for _, evidence := range finding.Evidence {
					fmt.Fprintf(&b, "- %s: %s\n", evidence.Label, evidence.Value)
				}
			}
			if report.Verbose && len(finding.Anchors) > 0 {
				fmt.Fprintln(&b)
				fmt.Fprintln(&b, reviewLabel(report, "**Anchors:**"))
				for _, anchor := range finding.Anchors {
					fmt.Fprintf(&b, "- `%s:%d`\n", anchor.File, anchor.Line)
				}
			}
			fmt.Fprintln(&b)
		}
	}
	fmt.Fprintln(&b)

	if report.Verbose {
		if len(report.Sources) > 0 {
			fmt.Fprintln(&b, reviewTitle(report, "## Sources"))
			for _, source := range report.Sources {
				fmt.Fprintf(&b, "- %s\n", renderSourceCatalogEntry(source))
			}
			fmt.Fprintln(&b)
		}

		if len(report.SourceRefs) > 0 {
			fmt.Fprintln(&b, reviewTitle(report, "## Context Sources"))
			for _, sourceRef := range report.SourceRefs {
				fmt.Fprintf(&b, "- %s\n", renderSourceRef(sourceRef))
			}
			fmt.Fprintln(&b)
		}

		fmt.Fprintln(&b, reviewTitle(report, "## Repo Facts"))
		fmt.Fprintf(&b, "- Tracked/source files scanned: `%d`\n", report.TrackedFileCount)
		fmt.Fprintf(&b, "- Test files: `%d`\n", report.TestFileCount)
		fmt.Fprintf(&b, "- Context snippets: `%d`\n", report.ContextSnippets)
		// Which models, over which wire. Kept together on one line each because
		// they are read together: "slow" is usually a model answer and "denied"
		// is usually a transport answer, and telling them apart starts here.
		if len(report.ReviewModels) > 0 {
			fmt.Fprintf(&b, "- AI models: `%s`\n", strings.Join(report.ReviewModels, "`, `"))
		}
		if transport := strings.TrimSpace(report.ReviewTransport); transport != "" {
			fmt.Fprintf(&b, "- AI transport: %s\n", transport)
		}
		if statement := report.Coverage.Statement(); statement != "" {
			fmt.Fprintf(&b, "- Coverage: %s\n", statement)
		}
		for _, status := range report.Evidence {
			fmt.Fprintf(&b, "- Evidence: %s\n", EvidenceSummaryLine(status))
		}
		if len(report.DependencyFiles) == 0 {
			fmt.Fprintln(&b, "- Dependency manifests: none detected")
		} else {
			fmt.Fprintf(&b, "- Dependency manifests: `%s`\n", strings.Join(report.DependencyFiles, "`, `"))
		}
		fmt.Fprintln(&b)

		fmt.Fprintln(&b, reviewTitle(report, "## Docs"))
		for _, doc := range report.Docs {
			status := "missing"
			if doc.Present {
				status = "present"
			}
			fmt.Fprintf(&b, "- `%s`: %s\n", doc.Path, status)
		}
		fmt.Fprintln(&b)

		fmt.Fprintln(&b, reviewTitle(report, "## Changed Files"))
		if len(report.ChangedFiles) == 0 {
			fmt.Fprintln(&b, "- none detected")
		} else {
			for _, file := range report.ChangedFiles {
				fmt.Fprintf(&b, "- `%s`\n", file)
			}
		}
		fmt.Fprintln(&b)
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}

// NothingToReviewMessage states plainly that no code was inspected. It exists
// so this outcome can never be confused with "reviewed and clean": a gate that
// reads "no material issues" for a diff nobody opened is a silent false pass.
func NothingToReviewMessage(report Report) string {
	target := strings.TrimSpace(report.ReviewTarget)
	if target == "" {
		target = "the working tree"
	}
	return fmt.Sprintf("Nothing to review: no changes found in %s. No code was inspected, so this is not a clean review.", target)
}

// renderCorroboration phrases cross-reviewer agreement, or "" when only one
// reviewer raised the finding. The single-reviewer case is deliberately silent
// rather than saying "flagged by one reviewer": that is the normal case, and
// annotating it would bury the signal in noise.
func renderCorroboration(finding Finding) string {
	if len(finding.Corroboration) < 2 {
		return ""
	}
	if len(finding.Corroboration) == 2 {
		return "Flagged by both reviewers independently — " + strings.Join(finding.Corroboration, " and ") + "."
	}
	return fmt.Sprintf("Flagged by both reviewers independently — %d reviewers: %s.",
		len(finding.Corroboration), strings.Join(finding.Corroboration, ", "))
}

// renderMergedFindings lists the folded-away versions of a finding, each with
// the fix its author proposed. Titles alone would not do: the title is what the
// merge decided was redundant, and the recommendation is what it must not lose.
func renderMergedFindings(finding Finding) []string {
	var out []string
	for _, merged := range finding.MergedFindings {
		title := strings.TrimSpace(merged.Title)
		if title == "" {
			continue
		}
		line := "“" + title + "”"
		if reviewer := strings.TrimSpace(merged.Reviewer); reviewer != "" {
			line += " (" + reviewer + ")"
		}
		if recommendation := strings.TrimSpace(merged.Recommendation); recommendation != "" {
			line += " — do next: " + recommendation
		}
		out = append(out, line)
	}
	return out
}

func renderFindingAttributions(report Report, finding Finding) []string {
	if len(finding.ResolvedSources) > 0 {
		seen := map[string]struct{}{}
		var out []string
		for _, src := range finding.ResolvedSources {
			label := strings.TrimSpace(ResolvedSourceLabel(src))
			if label == "" {
				continue
			}
			if _, ok := seen[label]; ok {
				continue
			}
			seen[label] = struct{}{}
			out = append(out, label)
		}
		return out
	}
	publishers := append([]string(nil), finding.SourcePublishers...)
	sources := sourceMap(report.Sources)
	for _, id := range finding.SourceIDs {
		if source, ok := sources[strings.TrimSpace(id)]; ok {
			publishers = append(publishers, firstNonEmpty(source.Publisher, source.Title, source.ID))
		}
	}
	if len(publishers) == 0 {
		return nil
	}
	seen := map[string]struct{}{}
	var out []string
	for _, publisher := range publishers {
		publisher = strings.TrimSpace(publisher)
		if publisher == "" {
			continue
		}
		if _, ok := seen[publisher]; ok {
			continue
		}
		seen[publisher] = struct{}{}
		out = append(out, publisher)
	}
	return out
}

func sourceMap(sources []Source) map[string]Source {
	out := make(map[string]Source, len(sources))
	for _, source := range sources {
		if strings.TrimSpace(source.ID) == "" {
			continue
		}
		out[source.ID] = source
	}
	return out
}

func renderSourceCatalogEntry(source Source) string {
	title := strings.TrimSpace(source.Title)
	if title == "" {
		title = strings.TrimSpace(source.ID)
	}
	if title == "" {
		return "unknown source"
	}
	if url := strings.TrimSpace(source.URL); url != "" {
		return fmt.Sprintf("[%s](%s)", title, url)
	}
	if id := strings.TrimSpace(source.ID); id != "" {
		return fmt.Sprintf("%s (`%s`)", title, id)
	}
	return title
}

func renderSourceRef(ref SourceRef) string {
	label := strings.TrimSpace(ref.ID)
	if label == "" {
		label = strings.TrimSpace(ref.Kind)
	}
	if label == "" {
		label = "context"
	}
	title := strings.TrimSpace(ref.Title)
	if title == "" {
		title = strings.TrimSpace(ref.File)
	}
	if title == "" {
		title = label
	}
	var details []string
	if url := strings.TrimSpace(ref.URL); url != "" {
		details = append(details, fmt.Sprintf("[%s](%s)", title, url))
	} else if ref.File != "" {
		if ref.StartLine > 0 {
			details = append(details, fmt.Sprintf("`%s:%d`", ref.File, ref.StartLine))
		} else {
			details = append(details, fmt.Sprintf("`%s`", ref.File))
		}
	} else {
		details = append(details, fmt.Sprintf("`%s`", title))
	}
	if ref.Source != "" && ref.Source != "local" {
		details = append(details, fmt.Sprintf("source=%s", ref.Source))
	}
	if ref.Publisher != "" {
		details = append(details, fmt.Sprintf("publisher=%s", ref.Publisher))
	}
	return fmt.Sprintf("`%s` %s", label, strings.Join(details, " · "))
}

func reviewTitle(report Report, text string) string {
	if !report.Color || text == "" || !termstyle.Enabled() {
		return text
	}
	return reviewMintANSI + text + reviewResetANSI
}

func reviewLabel(report Report, text string) string {
	if !report.Color || text == "" || !termstyle.Enabled() {
		return text
	}
	return reviewMintANSI + text + reviewResetANSI
}

func normalizeOptions(opts Options) Options {
	opts.Prompt = strings.TrimSpace(opts.Prompt)
	opts.Focus = strings.TrimSpace(opts.Focus)
	scope := strings.ToLower(strings.TrimSpace(opts.Scope))
	if scope == "" {
		opts.Scope = DefaultScope
		if !opts.Deep && opts.Prompt == "" {
			opts.PatchFocused = true
		}
	} else {
		opts.Scope = scope
	}
	if opts.WholeRepo {
		// A whole-repo review and the patch-focused filter are contradictory:
		// that filter keeps only findings anchored in the diff, which is
		// exactly the reach WholeRepo exists to restore.
		opts.PatchFocused = false
	}
	return opts
}

func baselineFor(scope string) []string {
	var out []string
	for _, baseline := range baselineScopes {
		if baseline != scope {
			out = append(out, baseline)
		}
	}
	return out
}

func depthLabel(deep bool) string {
	if deep {
		return "deep"
	}
	return "shallow"
}

// reviewProfileWholeRepo is the profile token for a whole-repo review. Every
// profile token the brief can carry must be defined in the instructions sent
// with it (see reviewDeveloperPrompt), so this name is shared rather than
// spelled twice.
const reviewProfileWholeRepo = "whole_repo"

// reviewProfile names what kind of review this is for the model. WholeRepo
// comes before Deep because the profile names the subject; depth travels
// separately in the brief's Depth field and is not lost by this ordering.
func reviewProfile(opts Options) string {
	if opts.WholeRepo {
		return reviewProfileWholeRepo
	}
	if opts.Deep {
		return "deep_full_spectrum"
	}
	if opts.PatchFocused {
		return "patch_focused"
	}
	if strings.TrimSpace(opts.Prompt) != "" {
		return "prompt_directed"
	}
	return "scope_focused"
}

type RepoFacts struct {
	Docs             []FilePresence
	ADRFiles         []string
	DependencyFiles  []string
	Files            []string
	GoPackages       []PackageFact
	TestFileCount    int
	TrackedFileCount int
}

type PackageFact struct {
	Path       string
	GoFiles    int
	TestFiles  int
	PublicName bool
}

type LocalScanner struct{}

func (LocalScanner) Scan(ctx context.Context, repoRoot string, focus string) (RepoFacts, error) {
	return scanRepo(ctx, repoRoot, focus)
}

func scanRepo(ctx context.Context, repoRoot, focus string) (RepoFacts, error) {
	docs := []FilePresence{
		{Path: "README.md", Present: exists(repoRoot, "README.md")},
		{Path: "AGENTS.md", Present: exists(repoRoot, "AGENTS.md")},
		{Path: "CONTEXT.md", Present: exists(repoRoot, "CONTEXT.md")},
		{Path: "docs/", Present: exists(repoRoot, "docs")},
	}
	facts := RepoFacts{Docs: docs}
	if files, ok := gitTrackedFiles(ctx, repoRoot); ok {
		for _, rel := range files {
			if focus != "" && !inFocus(rel, focus) {
				continue
			}
			if skipFile(rel) {
				continue
			}
			facts.addFile(rel)
		}
		sort.Strings(facts.DependencyFiles)
		facts.ADRFiles = adrFiles(facts.Files)
		facts.finalize()
		return facts, nil
	}
	err := filepath.WalkDir(repoRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "." {
			return nil
		}
		if entry.IsDir() {
			if skipDir(rel) {
				return filepath.SkipDir
			}
			return nil
		}
		if focus != "" && !inFocus(rel, focus) {
			return nil
		}
		if skipFile(rel) {
			return nil
		}
		facts.addFile(rel)
		return nil
	})
	sort.Strings(facts.DependencyFiles)
	facts.ADRFiles = adrFiles(facts.Files)
	facts.finalize()
	return facts, err
}

func (f *RepoFacts) addFile(rel string) {
	f.TrackedFileCount++
	f.Files = append(f.Files, rel)
	if isDependencyFile(rel) {
		f.DependencyFiles = append(f.DependencyFiles, rel)
	}
	if isTestFile(rel) {
		f.TestFileCount++
	}
}

func (f *RepoFacts) finalize() {
	sort.Strings(f.Files)
	f.GoPackages = goPackages(f.Files)
}

func goPackages(files []string) []PackageFact {
	byDir := map[string]PackageFact{}
	for _, file := range files {
		if !strings.HasSuffix(file, ".go") {
			continue
		}
		dir := filepath.ToSlash(filepath.Dir(file))
		if dir == "." {
			dir = "."
		}
		fact := byDir[dir]
		fact.Path = dir
		if isTestFile(file) {
			fact.TestFiles++
		} else {
			fact.GoFiles++
		}
		base := filepath.Base(dir)
		fact.PublicName = !strings.HasPrefix(base, "internal") && dir != "."
		byDir[dir] = fact
	}
	var out []PackageFact
	for _, fact := range byDir {
		out = append(out, fact)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Path < out[j].Path
	})
	return out
}

func gitTrackedFiles(ctx context.Context, repoRoot string) ([]string, bool) {
	cmd := gitCommand(ctx, repoRoot, "ls-files")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, false
	}
	var files []string
	for _, line := range strings.Split(out.String(), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, filepath.ToSlash(line))
		}
	}
	return files, true
}

func inFocus(rel, focus string) bool {
	focus = strings.Trim(strings.TrimSpace(filepath.ToSlash(focus)), "/")
	return focus == "" || rel == focus || strings.HasPrefix(rel, focus+"/")
}

func (f RepoFacts) observations() []string {
	var out []string
	if !present(f.Docs, "README.md") {
		out = append(out, "`README.md` is missing.")
	}
	if !present(f.Docs, "AGENTS.md") {
		out = append(out, "`AGENTS.md` is missing.")
	}
	if !present(f.Docs, "CONTEXT.md") {
		out = append(out, "`CONTEXT.md` is missing.")
	}
	if len(f.DependencyFiles) == 0 {
		out = append(out, "No dependency manifest was detected.")
	}
	if f.TestFileCount == 0 {
		out = append(out, "No test files were detected.")
	}
	return out
}

func present(docs []FilePresence, path string) bool {
	for _, doc := range docs {
		if doc.Path == path {
			return doc.Present
		}
	}
	return false
}

func adrFiles(files []string) []string {
	var out []string
	for _, file := range files {
		lower := strings.ToLower(file)
		base := filepath.Base(lower)
		if strings.Contains(lower, "/adr/") || strings.HasPrefix(base, "adr-") || strings.HasPrefix(base, "adr_") {
			out = append(out, file)
		}
	}
	sort.Strings(out)
	return out
}

func exists(root, rel string) bool {
	_, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel)))
	return err == nil
}

func skipDir(rel string) bool {
	for _, part := range strings.Split(rel, "/") {
		switch part {
		case ".git", ".jj", ".gx", ".gocache", "node_modules", "dist", "build", ".next", "coverage", ".cache", ".turbo", "vendor":
			return true
		}
	}
	return false
}

func skipFile(rel string) bool {
	lower := strings.ToLower(rel)
	for _, suffix := range []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".svg", ".woff", ".woff2", ".ttf", ".eot", ".mp4", ".mp3", ".zip", ".tar", ".gz", ".pdf", ".exe", ".dll", ".so", ".dylib"} {
		if strings.HasSuffix(lower, suffix) {
			return true
		}
	}
	return false
}

func isDependencyFile(rel string) bool {
	switch filepath.Base(rel) {
	case "go.mod", "go.sum", "package.json", "package-lock.json", "bun.lock", "bun.lockb", "pnpm-lock.yaml", "yarn.lock", "Cargo.toml", "Cargo.lock", "requirements.txt", "pyproject.toml", "poetry.lock",
		"packages.lock.json", "Directory.Packages.props", "pubspec.yaml", "pubspec.lock":
		return true
	}
	// .NET manifests carry the project name, so only the extension is fixed.
	switch strings.ToLower(filepath.Ext(rel)) {
	case ".csproj", ".sln":
		return true
	default:
		return false
	}
}

func isTestFile(rel string) bool {
	lower := strings.ToLower(rel)
	base := filepath.Base(lower)
	return strings.HasSuffix(base, "_test.go") ||
		strings.Contains(lower, "/test/") ||
		strings.Contains(lower, "/tests/") ||
		strings.Contains(base, ".test.") ||
		strings.Contains(base, ".spec.")
}

func changedFiles(ctx context.Context, repoRoot string) []string {
	// --untracked-files=all, because the default collapses a new directory into
	// a single `internal/gxtest/` entry. That entry is not a file: it produces
	// no diff, reads as nothing, and takes an entire directory of brand-new
	// unreviewed code out of the review without anything saying so. Listing the
	// files individually is what makes them reviewable and what makes the
	// coverage count true.
	cmd := gitCommand(ctx, repoRoot, "status", "--short", "--untracked-files=all")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil
	}
	var files []string
	for _, line := range strings.Split(out.String(), "\n") {
		if len(line) < 4 {
			continue
		}
		file := strings.TrimSpace(line[3:])
		if strings.Contains(file, " -> ") {
			parts := strings.Split(file, " -> ")
			file = parts[len(parts)-1]
		}
		if file != "" && !isVendoredOrGeneratedPath(file) {
			files = append(files, file)
		}
	}
	sort.Strings(files)
	return files
}

// vendoredPathSegments are directory names whose contents are never a review
// subject: installed dependencies, build output, and tool caches.
//
// `git status` only hides these when they are gitignored, and the moment a
// sub-project is added without its .gitignore they all arrive as untracked
// files. That is not hypothetical: a nine-line change in a repository with an
// un-ignored `dotcom/` reported 21,122 changed files and fanned out to 18
// shards — 36 model calls, two and a half minutes, to review nine lines. The
// model spent that budget reading `node_modules`, and the findings it returned
// were about the dependencies rather than the change.
//
// Filtering here rather than at the shard planner is deliberate: these files
// are not "too many to read", they are not the change. Coverage counts should
// never have included them, so a review that skips them is not a review that
// dropped anything.
var vendoredPathSegments = map[string]bool{
	"node_modules":     true,
	".next":            true,
	".nuxt":            true,
	".svelte-kit":      true,
	".turbo":           true,
	".parcel-cache":    true,
	"bower_components": true,
	"vendor":           true,
	"dist":             true,
	"build":            true,
	"out":              true,
	"target":           true,
	".venv":            true,
	"venv":             true,
	"__pycache__":      true,
	".mypy_cache":      true,
	".pytest_cache":    true,
	".tox":             true,
	".gradle":          true,
	".terraform":       true,
	"coverage":         true,
	".nyc_output":      true,
	".cache":           true,
	".output":          true,
	"Pods":             true,
	"DerivedData":      true,
	// dotnet build writes bin/ and obj/ beside every project. Only obj is
	// listed: a bare `bin` segment is also where repos keep hand-written
	// entrypoint scripts (Rails binstubs, bin/setup), which are review
	// subjects, and obj alone already catches the bulk of MSBuild output.
	"obj":        true,
	".dart_tool": true,
}

// isVendoredOrGeneratedPath reports whether a repo-relative path lives inside a
// dependency, build-output, or cache directory at any depth.
func isVendoredOrGeneratedPath(rel string) bool {
	rel = strings.TrimSpace(rel)
	if rel == "" {
		return false
	}
	for _, segment := range strings.Split(filepath.ToSlash(rel), "/") {
		if vendoredPathSegments[segment] {
			return true
		}
	}
	return false
}
