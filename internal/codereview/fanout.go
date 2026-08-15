package codereview

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Fan-out: the review reads all of its subject by splitting it across
// concurrent model calls instead of by throwing most of it away.
//
// The limits this replaces were all the same shape — take the first N, where N
// was chosen against a context window. A 181-file branch was reviewed as its
// first ten alphabetical files; a `--repo` review of a 273-file repository read
// twelve of them. Raising N trades one problem for another: a single call
// holding a megabyte of diff answers worse than four calls holding a quarter
// each, because everything in the middle stops being read carefully.
//
// So the budget below is per call, not per review, and the number of calls is
// whatever the subject requires. Coverage is not a knob here: every changed file
// and, in a whole-repo review, every source file lands in exactly one shard.
// What is a knob is concurrency, which is about wall-clock time.
const (
	// maxShardDiffBytes is the diff one model call carries. ~200 KB is roughly
	// 50k tokens of patch — large enough that ordinary branches stay a single
	// call, small enough that no call has a middle to lose.
	maxShardDiffBytes = 200000
	// maxShardRepoSourceBytes is the whole-repo source one call carries. Set
	// alongside the diff budget: a repo shard and a diff shard should cost the
	// model about the same attention.
	maxShardRepoSourceBytes = 200000
	// maxShardContextBytes is the shared context every shard carries — docs,
	// policy, retrieved code, prior findings. It rides along in each call
	// because a shard cannot be reviewed against evidence it was not given.
	maxShardContextBytes = 160000
	// maxChangeInFocusBytes is when a whole-repo shard also carries the diff.
	// Small changes are cheap to repeat and give every repo shard the context
	// of what just happened; large ones get their own shards instead.
	maxChangeInFocusBytes = 40000
	// fastShardDiffBytes is the per-call diff budget for a fast review.
	//
	// The ordinary budget optimizes for one model seeing as much of the change
	// as it can hold, which is right when nobody is waiting. A fast review
	// optimizes for wall clock, and wall clock is set by the slowest single
	// call: a model generates at a roughly fixed rate, so one call writing
	// findings for a 400-line diff takes about four times as long as four calls
	// writing findings for 100 lines each — and those four run concurrently.
	//
	// 50 KB keeps ordinary single-file work in one call (no extra cost for the
	// common case) while splitting the large multi-file changes that were the
	// only ones exceeding the latency budget.
	fastShardDiffBytes = 50000
	// maxReviewShards bounds the fan-out. It is a guard against a repository
	// far larger than anything gx reviews today, not a coverage decision: if it
	// ever binds, the shards it dropped are counted and reported.
	maxReviewShards = 256
	// defaultFanOutConcurrency is how many shards are in flight at once.
	// Overridable with GX_REVIEW_FANOUT_CONCURRENCY. Each shard is two
	// concurrent model calls (one per reviewer leg), so 10 shards is up to 20
	// in-flight Bedrock calls. Throttles this induces are absorbed by the
	// transport's retry (see bedrockRetryBackoffs) rather than failing shards.
	defaultFanOutConcurrency = 10
)

// shardKindDiff and shardKindRepo name what a shard is reading. A review can
// produce both: `--repo` reviews the repository with the current change still
// in focus, and those are different reading jobs.
const (
	shardKindDiff = "change"
	shardKindRepo = "repository"
)

// ReviewShard is one model call's worth of the subject.
type ReviewShard struct {
	Index int
	Total int
	Kind  string
	Brief ReviewBrief
	// Files is what this shard is responsible for reading. It is the unit
	// coverage is counted in.
	Files []string
}

// Label names a shard for progress output and for the finding IDs it produces.
func (s ReviewShard) Label() string {
	if s.Total <= 1 {
		return s.Kind
	}
	return fmt.Sprintf("%s %d/%d", s.Kind, s.Index+1, s.Total)
}

// planReviewShards splits a brief into calls that each fit the per-call budget.
//
// The split is by subject, never by evidence: shared context (docs, policy,
// retrieved code, prior findings) is carried by every shard, and only the
// subject — the diff, and in a whole-repo review the repository's source — is
// partitioned. A shard that saw a quarter of the change but none of the
// project's conventions would produce a quarter of a review four times.
func planReviewShards(brief ReviewBrief, opts Options) ([]ReviewShard, Coverage) {
	shared, sharedDropped := trimSharedContext(brief.Context)
	repoSnippets := partitionRepoSourceSnippets(brief.Context)

	diffSnippets, diffTruncated := truncateDiffSnippets(brief.Static.DiffSnippets)
	// A fast review splits the diff into smaller calls so they run concurrently;
	// see fastShardDiffBytes.
	diffBudget := maxShardDiffBytes
	if opts.Fast {
		diffBudget = fastShardDiffBytes
	}
	diffGroups := packDiffSnippets(diffSnippets, diffBudget)
	repoGroups := groupContextSnippets(repoSnippets, maxShardRepoSourceBytes)

	coverage := Coverage{
		Planned:        true,
		Subject:        "the change",
		Units:          "changed files",
		Total:          len(normalizedChangedFiles(brief.Static.ChangedFiles)),
		Truncated:      diffTruncated,
		ContextDropped: sharedDropped,
	}
	if brief.ReviewProfile == reviewProfileWholeRepo {
		coverage.Subject = "the repository"
		coverage.Units = "reviewable files"
		coverage.Total = len(normalizedChangedFiles(brief.Static.ChangedFiles)) + len(distinctSnippetRefs(repoSnippets))
		if outOfScope := brief.Static.FileCount - coverage.Total; outOfScope > 0 {
			coverage.OutOfScope = outOfScope
		}
	}

	total := len(diffGroups) + len(repoGroups)
	if total > maxReviewShards {
		// Never silently: drop from the tail (lowest-ranked material) and let
		// the coverage numbers say what was left unread.
		keepDiff := len(diffGroups)
		if keepDiff > maxReviewShards {
			keepDiff = maxReviewShards
		}
		diffGroups = diffGroups[:keepDiff]
		keepRepo := maxReviewShards - keepDiff
		if keepRepo < 0 {
			keepRepo = 0
		}
		if keepRepo < len(repoGroups) {
			repoGroups = repoGroups[:keepRepo]
		}
	}

	// Fan out only when the subject does not fit one call. Ordinary changes,
	// docs-only changes, prompt-directed questions and small repositories all
	// stay exactly what they were: one model call, holding everything. Splitting
	// a subject that fits would spend calls to make the review worse, since each
	// shard would then be reasoning about a fragment for no reason.
	if len(diffGroups) <= 1 && len(repoGroups) <= 1 && snippetBytes(diffSnippets)+contextBytes(repoSnippets) <= diffBudget {
		shard := brief
		shard.Static.DiffSnippets = diffSnippets
		shard.Context = append(append([]ContextSnippet(nil), shared...), repoSnippets...)
		files := append(diffSnippetFiles(diffSnippets), distinctSnippetRefs(repoSnippets)...)
		coverage.Read = len(files)
		coverage.Shards = 1
		return []ReviewShard{{Total: 1, Kind: shardKindDiff, Brief: shard, Files: files}}, coverage
	}

	focus := changeInFocusSnippets(diffSnippets)
	shards := make([]ReviewShard, 0, len(diffGroups)+len(repoGroups))
	read := 0
	for _, group := range diffGroups {
		shard := brief
		shard.Static.DiffSnippets = group
		shard.Context = shared
		files := diffSnippetFiles(group)
		read += len(files)
		shards = append(shards, ReviewShard{Kind: shardKindDiff, Brief: shard, Files: files})
	}
	for _, group := range repoGroups {
		shard := brief
		shard.Static.DiffSnippets = focus
		shard.Context = append(append([]ContextSnippet(nil), shared...), group...)
		files := distinctSnippetRefs(group)
		read += len(files)
		shards = append(shards, ReviewShard{Kind: shardKindRepo, Brief: shard, Files: files})
	}
	for i := range shards {
		shards[i].Index = i
		shards[i].Total = len(shards)
	}
	coverage.Read = read
	coverage.Shards = len(shards)
	return shards, coverage
}

// trimSharedContext enforces the per-call context budget in priority order and
// reports how many snippets did not fit.
//
// This is the one place sampling is still allowed, and deliberately: context is
// supporting evidence with a ranking, not the subject. Fanning it out would
// multiply calls without adding coverage of anything anyone asked about.
func trimSharedContext(snippets []ContextSnippet) ([]ContextSnippet, int) {
	var shared []ContextSnippet
	for _, snippet := range snippets {
		if snippet.Kind == "repo_source_file" {
			continue
		}
		shared = append(shared, snippet)
	}
	// Rank a list of indices rather than the snippets themselves, so the kept
	// set can be restored to the brief's original order and duplicate snippets
	// stay distinguishable.
	order := make([]int, len(shared))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(i, j int) bool {
		return contextSnippetPriority(shared[order[i]]) < contextSnippetPriority(shared[order[j]])
	})
	keep := make([]bool, len(shared))
	used := 0
	dropped := 0
	for _, idx := range order {
		size := len(shared[idx].Text)
		// The highest-ranked snippet is always kept, whatever its size: a
		// budget that can reject everything is a budget that can produce an
		// empty context and call it a review.
		if used > 0 && used+size > maxShardContextBytes {
			dropped++
			continue
		}
		used += size
		keep[idx] = true
	}
	out := make([]ContextSnippet, 0, len(shared))
	for i, snippet := range shared {
		if keep[i] {
			out = append(out, snippet)
		}
	}
	return out, dropped
}

func partitionRepoSourceSnippets(snippets []ContextSnippet) []ContextSnippet {
	var out []ContextSnippet
	for _, snippet := range snippets {
		if snippet.Kind == "repo_source_file" {
			out = append(out, snippet)
		}
	}
	return out
}

// truncateDiffSnippets applies the per-snippet diff budget and reports how many
// files arrived only in part, so a half-read file is never counted as read.
func truncateDiffSnippets(snippets []DiffSnippet) ([]DiffSnippet, int) {
	truncated := 0
	out := make([]DiffSnippet, 0, len(snippets))
	for _, snippet := range snippets {
		text := truncateAtHunkBoundary(snippet.Diff, maxAIDiffSnippetBytes)
		if len(text) < len(strings.TrimSpace(snippet.Diff)) {
			truncated++
		}
		snippet.Diff = text
		out = append(out, snippet)
	}
	return out, truncated
}

// packDiffSnippets packs diffs into per-call groups, preserving the impact
// ranking so shard 1 carries the highest-impact changes.
func packDiffSnippets(snippets []DiffSnippet, budget int) [][]DiffSnippet {
	if len(snippets) == 0 {
		return nil
	}
	var groups [][]DiffSnippet
	var current []DiffSnippet
	used := 0
	for _, snippet := range snippets {
		size := len(snippet.Diff)
		if len(current) > 0 && used+size > budget {
			groups = append(groups, current)
			current = nil
			used = 0
		}
		current = append(current, snippet)
		used += size
	}
	if len(current) > 0 {
		groups = append(groups, current)
	}
	return groups
}

func snippetBytes(snippets []DiffSnippet) int {
	total := 0
	for _, snippet := range snippets {
		total += len(snippet.Diff)
	}
	return total
}

func contextBytes(snippets []ContextSnippet) int {
	total := 0
	for _, snippet := range snippets {
		total += len(snippet.Text)
	}
	return total
}

func groupContextSnippets(snippets []ContextSnippet, budget int) [][]ContextSnippet {
	if len(snippets) == 0 {
		return nil
	}
	var groups [][]ContextSnippet
	var current []ContextSnippet
	used := 0
	for _, snippet := range snippets {
		size := len(snippet.Text)
		if len(current) > 0 && used+size > budget {
			groups = append(groups, current)
			current = nil
			used = 0
		}
		current = append(current, snippet)
		used += size
	}
	if len(current) > 0 {
		groups = append(groups, current)
	}
	return groups
}

// changeInFocusSnippets is the diff a whole-repo shard also carries, when the
// change is small enough to repeat. A repo shard's job is the repository, but
// knowing what just changed is what turns a generic observation into a finding
// about this branch.
func changeInFocusSnippets(snippets []DiffSnippet) []DiffSnippet {
	total := snippetBytes(snippets)
	if total == 0 || total > maxChangeInFocusBytes {
		return nil
	}
	return snippets
}

func diffSnippetFiles(snippets []DiffSnippet) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, snippet := range snippets {
		file := strings.TrimSpace(snippet.File)
		if file == "" {
			continue
		}
		if _, ok := seen[file]; ok {
			continue
		}
		seen[file] = struct{}{}
		out = append(out, file)
	}
	return out
}

func distinctSnippetRefs(snippets []ContextSnippet) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, snippet := range snippets {
		ref := strings.TrimSpace(snippet.Ref)
		if ref == "" {
			continue
		}
		// Large files are split across several snippets by
		// wholeRepoSourceSnippets; they are one file for coverage purposes.
		if idx := strings.Index(ref, " (part "); idx > 0 {
			ref = ref[:idx]
		}
		if _, ok := seen[ref]; ok {
			continue
		}
		seen[ref] = struct{}{}
		out = append(out, ref)
	}
	return out
}

// shardedReviewResult is what a fan-out review produced.
type shardedReviewResult struct {
	Findings []Finding
	// Story is the reviewer's "what changed that you now own" lane, taken
	// from the first shard that produced one — the summary call carries it
	// beside the findings, so it costs no extra model round-trip.
	Story    []StoryItem
	Coverage Coverage
	// Err is set only when every shard failed, which is a failed review rather
	// than a partial one.
	Err error
}

// runShardedReview asks the reviewer about each shard concurrently and merges
// what comes back.
//
// Partial failure is reported, not swallowed: a shard that errors leaves a
// region of the subject unreviewed, and Coverage.ShardsFailed is what tells the
// reader their "no issues" covers less than it says.
// reviewShard asks one reviewer about one shard. Every real reviewer can
// answer with a summary — findings plus the story lane in the same reply — so
// that is what is asked for; a bare AIReviewer (a test stub, a future
// findings-only provider) is asked for findings and yields no story.
func reviewShard(ctx context.Context, reviewer AIReviewer, brief ReviewBrief) ([]Finding, []StoryItem, error) {
	if withSummary, ok := reviewer.(AIReviewerWithSummary); ok {
		summary, err := withSummary.ReviewForSummary(ctx, brief)
		if err != nil {
			return nil, nil, err
		}
		return summary.Findings, summary.Story, nil
	}
	findings, err := reviewer.Review(ctx, brief)
	return findings, nil, err
}

func runShardedReview(ctx context.Context, reviewer AIReviewer, shards []ReviewShard, coverage Coverage, opts Options) shardedReviewResult {
	if len(shards) == 0 {
		return shardedReviewResult{Coverage: coverage}
	}
	if len(shards) == 1 {
		findings, story, err := reviewShard(ctx, reviewer, shards[0].Brief)
		if err != nil {
			coverage.ShardsFailed = 1
			coverage.Read = 0
			return shardedReviewResult{Coverage: coverage, Err: err}
		}
		return shardedReviewResult{Findings: findings, Story: story, Coverage: coverage}
	}

	enhanceProgress(opts, fmt.Sprintf("Asking AI reviewer (%d parallel reviews)", len(shards)))
	type result struct {
		index    int
		findings []Finding
		story    []StoryItem
		err      error
	}
	results := make([]result, len(shards))
	sem := make(chan struct{}, fanOutConcurrency(len(shards)))
	var progressMu sync.Mutex
	done := 0
	var wg sync.WaitGroup
	for i, shard := range shards {
		wg.Add(1)
		go func(i int, shard ReviewShard) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			findings, story, err := reviewShard(ctx, reviewer, shard.Brief)
			results[i] = result{index: i, findings: namespaceShardFindings(findings, shard), story: story, err: err}
			// Progress is serialized because the shards are not. A fan-out over
			// a large repository is minutes of silence otherwise.
			progressMu.Lock()
			done++
			state := "done"
			if err != nil {
				state = "failed"
			}
			enhanceProgress(opts, fmt.Sprintf("  %s %s (%d/%d)", shard.Label(), state, done, len(shards)))
			progressMu.Unlock()
		}(i, shard)
	}
	wg.Wait()

	var merged []Finding
	var story []StoryItem
	failed := 0
	unreadFiles := 0
	var firstErr error
	for i, res := range results {
		if res.err != nil {
			failed++
			unreadFiles += len(shards[i].Files)
			if firstErr == nil {
				firstErr = res.err
			}
			continue
		}
		merged = mergeFindings(merged, res.findings)
		// One story per review: shards see disjoint slices of the subject, so
		// their stories are about different files, but the reader wants one
		// list. Take the first shard's and let it lead; a merged story is a
		// follow-up if it proves needed.
		if len(story) == 0 && len(res.story) > 0 {
			story = res.story
		}
	}
	coverage.ShardsFailed = failed
	coverage.Read -= unreadFiles
	if coverage.Read < 0 {
		coverage.Read = 0
	}
	if failed == len(shards) {
		return shardedReviewResult{Coverage: coverage, Err: firstErr}
	}
	return shardedReviewResult{Findings: merged, Story: story, Coverage: coverage}
}

// namespaceShardFindings makes shard findings survive the merge.
//
// The AI reviewer numbers its findings ai.review.1, ai.review.2, … per call,
// and mergeFindings deduplicates on that ID. Without this, shard 2's findings
// would all collide with shard 1's and be dropped: the fan-out would read the
// whole subject and then report a single shard's worth of it. Findings from a
// single-shard review keep their original IDs, so nothing that does not fan out
// changes at all.
func namespaceShardFindings(findings []Finding, shard ReviewShard) []Finding {
	if shard.Total <= 1 {
		return findings
	}
	out := make([]Finding, 0, len(findings))
	for _, finding := range findings {
		if id := strings.TrimSpace(finding.ID); id != "" {
			finding.ID = fmt.Sprintf("%s.s%d", id, shard.Index+1)
		}
		out = append(out, finding)
	}
	return out
}

func fanOutConcurrency(shards int) int {
	limit := defaultFanOutConcurrency
	if raw := strings.TrimSpace(os.Getenv("GX_REVIEW_FANOUT_CONCURRENCY")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > shards {
		limit = shards
	}
	if limit < 1 {
		limit = 1
	}
	return limit
}
