package codereview

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

const (
	defaultReviewMaxOutputTokens = 6000
	// maxAIContextSnippetBytes is THE per-snippet budget. It used to be 1200
	// here while brief.go read 3200 and internal/publication read 1800, so two
	// thirds of every locally read file was assembled and then thrown away one
	// layer down. One number now, and brief.go references this one.
	//
	// 16000 is measured, not picked: across this repository's 273 non-test
	// source files the median is 3,758 bytes and p90 is 17,380, so 16 KB shows
	// 88% of them end to end (96% in the console repo) against 15% at 1200. The
	// files it still cannot hold whole are the 50–85 KB outliers, and those are
	// the fan-out reader's job — it splits them across snippets rather than
	// truncating, because a review that reads the first 1200 bytes of a file has
	// not read the file.
	maxAIContextSnippetBytes = 16000
	// maxAIDiffSnippetBytes was 6000 while brief.go rendered 20000 — the same
	// defect in the other direction. 32000 holds every per-file diff in this
	// working tree (largest: 29,780 bytes) so no changed file arrives half-read.
	maxAIDiffSnippetBytes      = 32000
	maxAIStaticToolOutputBytes = 4000
	// maxAIContextSnippets was 14, chosen when retrieval returned little more
	// than local docs. Now that the code index, the captured sessions and the
	// review corpus all answer, 14 slots are spent before the first indexed code
	// chunk is reached — the evidence the review goes out of its way to fetch
	// was being retrieved and then discarded. The budget is what it costs to
	// look at the evidence, and looking is the point.
	//
	// The count is no longer the real bound: maxShardContextBytes is, and a
	// brief that exceeds it is split into shards rather than trimmed. These
	// remain as the guard for callers that hand a brief straight to a reviewer
	// without planning shards, internal/publication among them.
	maxAIContextSnippets     = 48
	maxAIDeepContextSnippets = 96
	// maxAIWholeRepoContextSnippets is the whole-repo budget. The repository
	// inventory and its source files are the subject of that review, and the
	// default budget is small enough that they would be cut before the model
	// ever saw them.
	maxAIWholeRepoContextSnippets = 128
	// maxAIRepoInventoryBytes is the file listing's own budget. A map of the
	// repository is worth its bytes: at ~45 bytes a path this holds roughly
	// 4,000 files, which is every file in both repositories gx reviews today.
	maxAIRepoInventoryBytes = 180000
	// Diff snippets are the primary evidence for what changed, so they get their
	// own budget rather than sharing the retrieved-context one. A PR summary
	// builds a wide diff on purpose; capping it at the context limit silently
	// hid most of a large change from the model.
	//
	// 32 was below the file count of an ordinary branch — this working tree has
	// 181 changed files — so the review simply did not see most of the change.
	// Like the context count, this is now a guard rather than the bound: the
	// shard planner packs every changed file into some shard.
	maxAIDiffSnippets     = 240
	maxAIDeepDiffSnippets = 480
	maxAICodeQualityHints = 20
	maxAIModuleSummaries  = 12
	maxAIChangedFiles     = 400
)

// The review composition: two competing Bedrock reviewers plus an independent
// Bedrock judge.
//
// Two flagship reviewers cost max(A,B) in wall clock, not A+B, because
// multiAIReviewer already runs its legs concurrently — so the second opinion is
// close to free in latency and only costs tokens. The pair is deliberately two
// Opus models from DIFFERENT GENERATIONS rather than one Opus plus a smaller
// same-generation sibling: a sibling trained alongside its larger cousin tends
// to miss the same things it does, and correlated misses are exactly what a
// second reviewer is supposed to catch.
//
// The judge is a third model on purpose. It decides which candidate findings
// survive, and a model grading its own output is not a filter. It also sits on
// the sequential path after both reviewers return, which is where latency hurts
// most, so it is the smaller Sonnet rather than a third Opus.
//
// Every ID here MUST be the `us.`-prefixed inference profile form. Bare
// `anthropic.*` model IDs are rejected by bedrock-runtime for on-demand
// invocation ("Invocation of model ID ... with on-demand throughput isn't
// supported"), which the previous default silently was — the Anthropic reviewer
// could never have run. normalizeBedrockModelID enforces this for overrides too,
// and TestDefaultBedrockModelsAreInferenceProfiles is the regression test.
const (
	defaultBedrockReviewModelA = "us.anthropic.claude-opus-4-6-v1"
	defaultBedrockReviewModelB = "us.anthropic.claude-opus-4-5-20251101-v1:0"
	defaultBedrockJudgeModel   = "us.anthropic.claude-sonnet-4-6"
	// defaultBedrockRegion was us-east-1, which is not where these inference
	// profiles are enabled for this deployment; a model that exists in one
	// region reports as "The provided model identifier is invalid" in another,
	// which reads like a bad model name rather than a bad region.
	defaultBedrockRegion = "us-west-2"
)

type AIReviewer interface {
	Review(ctx context.Context, brief ReviewBrief) ([]Finding, error)
}

type NotableChange struct {
	File string
	Line int
	Note string
}

type PRSummaryReview struct {
	Overview         string
	DownstreamImpact string
	NotableChanges   []NotableChange
	Findings         []Finding
}

type AIReviewerWithSummary interface {
	AIReviewer
	ReviewForSummary(ctx context.Context, brief ReviewBrief) (PRSummaryReview, error)
}

type ReviewerInfo struct {
	Models []string
	// Transport is how the panel reached bedrock-runtime, phrased for a human
	// ("GX Cloud (https://api.gx.run)" / "direct AWS credentials (us-west-2)").
	// It is reported rather than inferred because the two have different
	// latency and different failure modes, and a review that does not say which
	// one ran leaves both questions unanswerable after the fact.
	Transport string
}

type namedAIReviewer struct {
	name     string
	label    string
	reviewer AIReviewer
}

type multiAIReviewer struct {
	reviewers []namedAIReviewer
	// legFailures records legs that failed while at least one other leg
	// answered. That case used to be silent, and the silence was expensive:
	// the B leg's inference profile ID contains a colon, the SigV4 canonical
	// path was not double-encoded, and every single B call failed with
	// SignatureDoesNotMatch — for as long as the two-model panel has existed.
	// Because leg A answered, ReviewForSummary returned success and the review
	// reported two models in its own "AI models" line while running one.
	//
	// A pointer with its own lock because one multiAIReviewer value is shared
	// across every shard's goroutine.
	legFailures *legFailureLog
	// adjudicator decides which of the two legs' findings are the same finding.
	// Nil is a supported state: nothing semantic is merged, only near-verbatim
	// copies, and the caller reports that it could not do better.
	adjudicator duplicateAdjudicator
}

// legFailureLog collects distinct reviewer-leg failures across concurrent shards.
type legFailureLog struct {
	mu       sync.Mutex
	seen     map[string]struct{}
	messages []string
}

func newLegFailureLog() *legFailureLog {
	return &legFailureLog{seen: map[string]struct{}{}}
}

func (l *legFailureLog) record(label string, err error) {
	if l == nil || err == nil {
		return
	}
	message := label + ": " + err.Error()
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.seen[message]; ok {
		return
	}
	l.seen[message] = struct{}{}
	l.messages = append(l.messages, message)
}

func (l *legFailureLog) reasons() []string {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.messages...)
}

// PartialReviewerFailures is why a review ran with fewer models than it was
// configured for, or nil when every leg answered. The engine reports these as
// degraded reasons: "two flagship models reviewed this" and "one did, twice as
// cheaply as you think" are different reviews and must not print the same.
func PartialReviewerFailures(reviewer AIReviewer) []string {
	multi, ok := reviewer.(multiAIReviewer)
	if !ok {
		return nil
	}
	return multi.legFailures.reasons()
}

type unavailableAIReviewer struct {
	reason string
}

func reviewerInfoFromReviewer(reviewer AIReviewer) ReviewerInfo {
	if reviewer == nil {
		return ReviewerInfo{}
	}
	if multi, ok := reviewer.(multiAIReviewer); ok {
		var models []string
		for _, item := range multi.reviewers {
			if model := reviewerModelName(item.reviewer); model != "" {
				models = append(models, model)
			}
		}
		return ReviewerInfo{Models: models, Transport: ReviewerTransport(multi)}
	}
	if model := reviewerModelName(reviewer); model != "" {
		return ReviewerInfo{Models: []string{model}, Transport: ReviewerTransport(reviewer)}
	}
	return ReviewerInfo{}
}

func reviewerModelName(reviewer AIReviewer) string {
	if r, ok := reviewer.(*bedrockAnthropicReviewer); ok && r != nil {
		return strings.TrimSpace(r.model)
	}
	return ""
}

// ReviewerTransport is the human phrase for how this reviewer reaches Bedrock,
// or "" when it has no transport (no reviewer, or a test double).
func ReviewerTransport(reviewer AIReviewer) string {
	if multi, ok := reviewer.(multiAIReviewer); ok {
		for _, item := range multi.reviewers {
			if detail := ReviewerTransport(item.reviewer); detail != "" {
				return detail
			}
		}
		return ""
	}
	if r, ok := reviewer.(*bedrockAnthropicReviewer); ok && r != nil && r.transport != nil {
		return r.transport.detail()
	}
	return ""
}

// reviewerFromEnv builds the review panel: two competing Bedrock
// reviewers that run concurrently. There is no other provider.
//
// The OpenAI reviewer and its cloud proxy used to live here as a fallback and
// were deleted rather than left dormant. Bedrock is the only review provider,
// so a second one would only ever be reached when the first was misconfigured —
// and a fallback that answers when the configured reviewer cannot is precisely
// what makes a broken configuration invisible. Missing AWS credentials must read
// as "no reviewer ran", not as a clean review from a model nobody asked for.
//
// Embeddings are unaffected: internal/semantic reads OPENAI_API_KEY directly
// (see defaultEmbedderFactory in turbopuffer_index.go) and never goes through an
// AIReviewer, so the code index still embeds on OpenAI.
func reviewerFromEnv() AIReviewer {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_AI")), "0") {
		return nil
	}
	plan, err := resolveBedrockTransportPlan()
	if err != nil {
		// Both legs carry the same actionable reason, so whichever one a caller
		// inspects says what to fix.
		return multiAIReviewer{
			reviewers: []namedAIReviewer{
				{name: "bedrock-a", label: "Bedrock A", reviewer: unavailableAIReviewer{reason: err.Error()}},
				{name: "bedrock-b", label: "Bedrock B", reviewer: unavailableAIReviewer{reason: err.Error()}},
			},
			legFailures: newLegFailureLog(),
		}
	}
	modelA, modelB := resolveBedrockReviewModels()
	return multiAIReviewer{
		reviewers: []namedAIReviewer{
			{name: "bedrock-a", label: bedrockLegLabel("Bedrock A", modelA, plan.Kind), reviewer: newBedrockReviewer(plan.newTransport(), modelA)},
			{name: "bedrock-b", label: bedrockLegLabel("Bedrock B", modelB, plan.Kind), reviewer: newBedrockReviewer(plan.newTransport(), modelB)},
		},
		legFailures: newLegFailureLog(),
		// Built here rather than passed in because this is the only place that
		// knows a panel of two legs exists at all, and a panel is what produces
		// the duplicates.
		adjudicator: duplicateAdjudicatorFromEnv(),
	}
}

// reviewerDuplicateAdjudicator hands back the panel's own adjudicator so the
// engine's cross-shard de-duplication asks the same model, over the same wire,
// as the per-shard one. A reviewer that is not a panel has no duplicates to
// adjudicate and returns nil.
func reviewerDuplicateAdjudicator(reviewer AIReviewer) duplicateAdjudicator {
	if multi, ok := reviewer.(multiAIReviewer); ok {
		return multi.adjudicator
	}
	return nil
}

func (r unavailableAIReviewer) Review(context.Context, ReviewBrief) ([]Finding, error) {
	return nil, fmt.Errorf("%s", r.reason)
}

func (r unavailableAIReviewer) UnavailableReason() string { return r.reason }

func (r unavailableAIReviewer) Available() bool { return false }

func (m multiAIReviewer) Available() bool {
	for _, item := range m.reviewers {
		if reviewerAvailable(item.reviewer) {
			return true
		}
	}
	return false
}

type availabilityReporter interface {
	Available() bool
}

type unavailabilityReporter interface {
	UnavailableReason() string
}

func reviewerUnavailableReason(item namedAIReviewer) string {
	if reporter, ok := item.reviewer.(unavailabilityReporter); ok {
		if reason := strings.TrimSpace(reporter.UnavailableReason()); reason != "" {
			return reason
		}
	}
	return item.label + " reviewer is not configured"
}

// ReviewerUnavailableReason is why no AI reviewer can run, or "" if one can.
//
// Bedrock is now the only review provider, so an unavailable reviewer means no
// review happened at all — there is nothing left to fall back to. The engine
// reported that as "AI reviewer not configured", which named neither the
// provider nor the fix; the reviewer already carries a sentence that does both,
// and this is how the engine reaches it.
func ReviewerUnavailableReason(reviewer AIReviewer) string {
	if reviewer == nil {
		return "AI reviewer not configured"
	}
	if multi, ok := reviewer.(multiAIReviewer); ok {
		var reasons []string
		seen := map[string]struct{}{}
		for _, item := range multi.reviewers {
			if reviewerAvailable(item.reviewer) {
				continue
			}
			reason := reviewerUnavailableReason(item)
			if _, dup := seen[reason]; dup {
				continue
			}
			seen[reason] = struct{}{}
			reasons = append(reasons, reason)
		}
		if len(reasons) == 0 {
			return ""
		}
		return strings.Join(reasons, "; ")
	}
	if reviewerAvailable(reviewer) {
		return ""
	}
	return reviewerUnavailableReason(namedAIReviewer{label: "AI", reviewer: reviewer})
}

func reviewerAvailable(reviewer AIReviewer) bool {
	if reviewer == nil {
		return false
	}
	if reporter, ok := reviewer.(availabilityReporter); ok {
		return reporter.Available()
	}
	return true
}

func (m multiAIReviewer) Review(ctx context.Context, brief ReviewBrief) ([]Finding, error) {
	summary, err := m.ReviewForSummary(ctx, brief)
	return summary.Findings, err
}

func (m multiAIReviewer) ReviewForSummary(ctx context.Context, brief ReviewBrief) (PRSummaryReview, error) {
	type reviewerResult struct {
		item    namedAIReviewer
		summary PRSummaryReview
		err     error
	}
	results := make([]reviewerResult, len(m.reviewers))
	var wg sync.WaitGroup
	for i, item := range m.reviewers {
		i, item := i, item
		results[i] = reviewerResult{item: item}
		if !reviewerAvailable(item.reviewer) {
			// Use the reviewer's own reason. This used to be a fixed "%s
			// reviewer is not configured", which threw away the one string
			// that said whether the credentials were missing, the model was
			// denied, or the region was wrong.
			results[i].err = fmt.Errorf("%s", reviewerUnavailableReason(item))
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if withSummary, ok := item.reviewer.(AIReviewerWithSummary); ok {
				summary, err := withSummary.ReviewForSummary(ctx, brief)
				results[i] = reviewerResult{item: item, summary: summary, err: err}
				return
			}
			findings, err := item.reviewer.Review(ctx, brief)
			results[i] = reviewerResult{item: item, summary: PRSummaryReview{Findings: findings}, err: err}
		}()
	}
	wg.Wait()
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].item.name < results[j].item.name
	})
	var out []Finding
	var overview string
	var downstreamImpact string
	var notableChanges []NotableChange
	var errors []string
	parsed := false
	for _, result := range results {
		item := result.item
		if result.err != nil {
			errors = append(errors, item.label+": "+result.err.Error())
			m.legFailures.record(item.label, result.err)
			continue
		}
		parsed = true
		if overview == "" && strings.TrimSpace(result.summary.Overview) != "" {
			overview = strings.TrimSpace(result.summary.Overview)
		}
		if downstreamImpact == "" && strings.TrimSpace(result.summary.DownstreamImpact) != "" {
			downstreamImpact = strings.TrimSpace(result.summary.DownstreamImpact)
		}
		if len(notableChanges) == 0 && len(result.summary.NotableChanges) > 0 {
			notableChanges = append([]NotableChange(nil), result.summary.NotableChanges...)
		}
		for _, finding := range result.summary.Findings {
			finding.ID = item.name + "." + finding.ID
			finding.Evidence = append([]Evidence{{Label: "Reviewer", Value: item.label}}, finding.Evidence...)
			// The leg that raised it, recorded structurally rather than only in
			// an evidence footnote, so that a merge downstream can tell "two
			// models agree" from "one model said it twice".
			finding.Corroboration = []string{item.label}
			out = append(out, finding)
		}
	}
	out = mergeNearDuplicateFindings(ctx, m.adjudicator, out)
	if parsed {
		return PRSummaryReview{Overview: overview, DownstreamImpact: downstreamImpact, NotableChanges: notableChanges, Findings: out}, nil
	}
	if len(errors) > 0 {
		return PRSummaryReview{}, fmt.Errorf("AI reviewers failed: %s", strings.Join(errors, "; "))
	}
	return PRSummaryReview{}, fmt.Errorf("AI reviewers failed")
}

func compactReviewBriefForAI(brief ReviewBrief) ReviewBrief {
	deep := strings.EqualFold(brief.Depth, "deep")
	contextLimit := maxAIContextSnippets
	diffLimit := maxAIDiffSnippets
	if deep {
		contextLimit = maxAIDeepContextSnippets
		diffLimit = maxAIDeepDiffSnippets
	}
	if strings.TrimSpace(brief.ReviewProfile) == reviewProfileWholeRepo && contextLimit < maxAIWholeRepoContextSnippets {
		contextLimit = maxAIWholeRepoContextSnippets
	}

	brief.Static.DependencyFiles = limitStrings(brief.Static.DependencyFiles, maxAIChangedFiles)
	brief.Static.ChangedFiles = limitStrings(brief.Static.ChangedFiles, maxAIChangedFiles)
	brief.Static.DiffSnippets = compactDiffSnippets(brief.Static.DiffSnippets, diffLimit)
	brief.Static.Modules = limitModules(brief.Static.Modules, maxAIModuleSummaries)
	brief.Static.ToolResults = compactStaticToolResults(brief.Static.ToolResults)
	brief.Static.CodeQuality = limitCodeQualityHints(brief.Static.CodeQuality, maxAICodeQualityHints)
	brief.Context = compactContextSnippets(brief.Context, contextLimit)
	return brief
}

func compactDiffSnippets(snippets []DiffSnippet, limit int) []DiffSnippet {
	if limit > 0 && len(snippets) > limit {
		snippets = snippets[:limit]
	}
	out := make([]DiffSnippet, 0, len(snippets))
	for _, snippet := range snippets {
		snippet.Diff = truncateAtHunkBoundary(snippet.Diff, maxAIDiffSnippetBytes)
		if strings.TrimSpace(snippet.File) == "" || strings.TrimSpace(snippet.Diff) == "" {
			continue
		}
		out = append(out, snippet)
	}
	return out
}

func compactStaticToolResults(results []StaticToolResult) []StaticToolResult {
	out := make([]StaticToolResult, 0, len(results))
	for _, result := range results {
		if result.ExitCode == 0 && strings.TrimSpace(result.Output) == "" {
			result.Output = ""
			out = append(out, result)
			continue
		}
		result.Output = truncateReviewText(result.Output, maxAIStaticToolOutputBytes)
		out = append(out, result)
	}
	return out
}

func compactContextSnippets(snippets []ContextSnippet, limit int) []ContextSnippet {
	if limit > 0 && len(snippets) > limit {
		snippets = prioritizeContextSnippets(snippets)
		snippets = snippets[:limit]
	}
	out := make([]ContextSnippet, 0, len(snippets))
	for _, snippet := range snippets {
		snippet.Text = truncateReviewText(snippet.Text, contextSnippetByteLimit(snippet))
		out = append(out, snippet)
	}
	return out
}

// contextSnippetByteLimit is the per-snippet budget for the model. The repo
// inventory gets its own: it is a map of the repository, and cutting it to the
// generic budget would show a fraction of the file list and call it the
// repository. That kind only exists in a whole-repo brief, so no other profile
// is affected.
func contextSnippetByteLimit(snippet ContextSnippet) int {
	if snippet.Kind == "repo_inventory" {
		return maxAIRepoInventoryBytes
	}
	return maxAIContextSnippetBytes
}

func prioritizeContextSnippets(snippets []ContextSnippet) []ContextSnippet {
	out := append([]ContextSnippet(nil), snippets...)
	sort.SliceStable(out, func(i, j int) bool {
		return contextSnippetPriority(out[i]) < contextSnippetPriority(out[j])
	})
	return interleaveRetrievedSnippets(out)
}

// retrievedSnippetKinds are the kinds that come from a retrieval index rather
// than from the checkout.
var retrievedSnippetKinds = map[string]bool{
	"indexed_code":        true,
	"indexed_session":     true,
	"code_review_history": true,
	"code_review_summary": true,
	"review_resource":     true,
}

// interleaveRetrievedSnippets round-robins the retrieved kinds through the slots
// they already occupy, so the context budget cannot be consumed entirely by
// whichever source happened to return the most rows.
//
// Sorting alone does not give this. Each source has its own top-k, and with a
// strict order the first source fills the budget and the rest are cut whole —
// twelve code chunks and no sessions is not a better-informed review than eight
// code chunks, two sessions and two prior findings, it is a narrower one. Only
// the retrieved block is reordered; local documents keep their positions.
func interleaveRetrievedSnippets(snippets []ContextSnippet) []ContextSnippet {
	var slots []int
	byKind := map[string][]ContextSnippet{}
	var kindOrder []string
	for i, snippet := range snippets {
		if !retrievedSnippetKinds[snippet.Kind] {
			continue
		}
		slots = append(slots, i)
		if _, seen := byKind[snippet.Kind]; !seen {
			kindOrder = append(kindOrder, snippet.Kind)
		}
		byKind[snippet.Kind] = append(byKind[snippet.Kind], snippet)
	}
	if len(kindOrder) < 2 {
		return snippets
	}
	out := append([]ContextSnippet(nil), snippets...)
	next := 0
	for len(slots) > next {
		placed := false
		for _, kind := range kindOrder {
			queue := byKind[kind]
			if len(queue) == 0 {
				continue
			}
			out[slots[next]] = queue[0]
			byKind[kind] = queue[1:]
			next++
			placed = true
			if next >= len(slots) {
				break
			}
		}
		if !placed {
			break
		}
	}
	return out
}

// contextSnippetPriority orders the context budget. Lower is kept first.
//
// The rule the table now obeys, and did not before: code outranks manifests,
// and manifests outrank lockfiles. It used to read dependency_manifest=6,
// module_file=8, repo_source_file=8 — so when the budget bound, the first thing
// evicted was the repository's own source and the thing that survived was
// go.sum. A generated dependency lockfile is the least informative file in any
// repository and it was outranking the code under review.
//
// Prior review findings (code_review_history / code_review_summary) had no case
// at all and fell to the default, which meant every review re-derived what the
// last one already established.
const (
	priorityDomainDoc         = 0
	priorityReviewPolicy      = 1
	priorityADR               = 3
	priorityRepoDoc           = 4
	priorityRepoInventory     = 4
	priorityRetrievedEvidence = 5
	priorityPriorFindings     = 5
	priorityCheckoutSource    = 6
	prioritySessionEvidence   = 7
	priorityQualityFile       = 8
	priorityDependencyFile    = 9
	priorityUnknown           = 10
	priorityLockfile          = 11
)

func contextSnippetPriority(snippet ContextSnippet) int {
	switch snippet.Kind {
	case "domain_doc":
		return priorityDomainDoc
	case "review_policy":
		return priorityReviewPolicy
	case "adr":
		return priorityADR
	case "repo_doc", "codebase_doc":
		return priorityRepoDoc
	// The repository's own map and code outrank retrieved external guidance
	// when the repository is what is being reviewed; they are only present at
	// all in that case.
	case "repo_inventory":
		return priorityRepoInventory

	// Retrieved evidence and curated guidance, unchanged from where the
	// retrieval work put them.
	case "indexed_code", "review_resource":
		return priorityRetrievedEvidence

	// What a previous review already found here. Cheaper to be reminded than to
	// rediscover, and a regression on a known finding is the highest-signal
	// thing a review can report. These had no case at all and fell to the
	// default, below the manifests.
	case "code_review_history", "code_review_summary":
		return priorityPriorFindings

	// Source code read from the checkout. This is the tier that sat at 8, below
	// dependency_manifest at 6 — so when the budget bound, the first snippets
	// evicted were the files under review and the survivor was go.sum.
	case "repo_source_file", "module_file", "changed_file":
		return priorityCheckoutSource

	case "indexed_session", "session_transcript", "lexical_reach":
		return prioritySessionEvidence
	case "code_quality_file":
		return priorityQualityFile

	// Declared dependencies: go.mod and package.json say something a reviewer
	// can act on. Their generated lockfiles do not, and go last of everything —
	// below even an unrecognized snippet kind, which at least might be code.
	case "dependency_manifest":
		return priorityDependencyFile
	case "dependency_lockfile":
		return priorityLockfile

	default:
		if snippet.Source == "indexed" || strings.HasPrefix(snippet.Source, "turbopuffer:") {
			return priorityRetrievedEvidence
		}
		return priorityUnknown
	}
}

func limitCodeQualityHints(hints []CodeQualityHint, limit int) []CodeQualityHint {
	if limit <= 0 || len(hints) <= limit {
		return hints
	}
	return hints[:limit]
}

func limitModules(modules []ModuleSummary, limit int) []ModuleSummary {
	if limit <= 0 || len(modules) <= limit {
		return modules
	}
	return modules[:limit]
}

func limitStrings(values []string, limit int) []string {
	if limit <= 0 || len(values) <= limit {
		return values
	}
	return values[:limit]
}

func truncateReviewText(text string, limit int) string {
	text = strings.TrimSpace(text)
	if limit <= 0 || len(text) <= limit {
		return text
	}
	return text[:limit] + "\n[truncated]\n"
}

func truncateAtHunkBoundary(text string, limit int) string {
	text = strings.TrimSpace(text)
	if limit <= 0 || len(text) <= limit {
		return text
	}
	cut := limit
	if newline := strings.LastIndex(text[:limit], "\n"); newline > 0 {
		cut = newline
	}
	if hunk := strings.LastIndex(text[:cut], "\n@@"); hunk > limit/2 {
		cut = hunk
	}
	return strings.TrimSpace(text[:cut]) + "\n[truncated]\n"
}

// reviewDeveloperPrompt is the instruction string sent alongside the brief.
//
// It takes the brief because a profile the prompt never defines is worse than
// no profile at all: the model is told "use review_profile to choose behavior",
// is handed a token that appears nowhere in its instructions, and is also told
// to return an empty recommendations array when nothing meets the active
// profile. Profile-specific instructions are therefore appended for the profile
// in hand, which also keeps the pr_summary prompt byte-identical as profiles
// are added.
func reviewDeveloperPrompt(brief ReviewBrief) string {
	lines := baseReviewDeveloperPromptLines()
	if strings.TrimSpace(brief.ReviewProfile) == reviewProfileWholeRepo {
		lines = append(lines, wholeRepoPromptLines()...)
	}
	return strings.Join(lines, "\n")
}

func baseReviewDeveloperPromptLines() []string {
	return []string{
		"You are GX Review. Review the provided patch and context for concrete recommendations, not generic audit facts.",
		"Use review_profile and depth to choose behavior: patch_focused means current-change review; prompt_directed means use review_prompt to guide a broader review of how the current diff affects the surrounding codebase; scope_focused means the requested scope; deep_full_spectrum means full-spectrum review.",
		"Use triage.class and triage.risk_tags to weight your review: for security-sensitive changes prioritize the tagged risks; for mechanical changes only report real breakage.",
		"When review_prompt is present, answer it directly. Treat static.diff_snippets as evidence for why the prompted concern matters now, but inspect surrounding Modules, Interfaces, tests, docs, local policy, and retrieved context when they explain impact or the correct fix.",
		"For patch_focused and pr_summary reviews, prioritize concrete bugs, security/auth issues, data correctness, race/idempotency, error handling, missing tests, observability, deploy/CI risks, and dependency regressions introduced or exposed by static.diff_snippets.",
		"For patch_focused and pr_summary reviews, broad architecture, naming, docs, cleanup, or Module-depth advice is invalid unless it directly explains a changed-line bug or review risk.",
		"For prompt_directed reviews, prioritize findings where review_prompt, the current diff, and broader repo context intersect. Do not limit yourself to changed lines, but do not emit generic repo-wide advice unrelated to review_prompt.",
		"For deep_full_spectrum reviews, check security, bugs, data integrity, concurrency, idempotency, architecture, testing, observability, performance, dependencies, docs, and operability while still grounding every finding in changed files, tool output, local policy, or retrieved context.",
		"Treat REVIEW.md review_policy snippets as repo-local review instructions. Follow them unless they conflict with the explicit review_prompt, hard evidence in the changed patch, or safety/security requirements.",
		"Use the architecture vocabulary exactly when discussing structure: Module, Interface, Implementation, Depth, deep, shallow, seam, adapter, leverage, locality.",
		"Never use component, service, API, boundary, or layer when Module, Interface, seam, or adapter fits.",
		"Treat static facts and hints as clues only. Do not turn file counts, missing docs, or missing tests directly into findings.",
		"Use static_tool_results as hard evidence. Failed tests, vet warnings, compile errors, and linter-like diagnostics should outrank speculative architecture advice. Ignore skipped tool results unless the skipped reason itself is clearly actionable.",
		"Use static.diff_snippets as the primary evidence for what changed. Prefer findings tied to changed lines over repo-wide advice.",
		"Use code_quality_hints as concrete candidates. Confirm whether they matter from the provided snippets before recommending a fix.",
		"Treat source_refs as the attribution set. AI output is synthesis, not evidence; every recommendation must be traceable to static facts, diff snippets, context snippets, or source_refs.",
		"Every recommendation must name at least one changed file path, Module, static tool result, code_quality_hint, or context source label. Do not produce coverage-only or structure-only recommendations without concrete evidence.",
		"The recommendation field must be concrete work: name the specific files or Modules to touch, the first operation to perform, and the verification to run. Avoid vague verbs like assess, consider, clarify, improve, harden, or refactor unless followed by exact code actions.",
		"The benefit field must state the expected payoff in concrete engineering terms: performance, readability, fewer lines of code, better error handling, better testability, lower coupling, faster onboarding, more reproducible dependencies, or better observability.",
		"If you cannot name a concrete payoff, do not emit that recommendation.",
		"In deep_full_spectrum or architecture scope only, explore like the architecture skill: find deepening opportunities, leaked Implementation knowledge, shallow Interfaces, unclear seams, weak locality, test friction, and unjustified adapters.",
		"Apply the deletion test: if deleting a Module removes complexity, call it shallow; if complexity spreads across callers, the Module is earning its keep.",
		"Use dependency categories internally: in-process, local-substitutable, remote but owned, true external. Mention adapters only when the seam needs more than one adapter.",
		"Use source_catalog and labeled context snippets internally. It is okay to mention source labels like R1 or L2 in evidence, but never output source titles or URLs.",
		"Only produce recommendations tied to the provided repo context. Reject generic best-practice advice.",
		"If no concrete issue meets the active profile, return an empty recommendations array.",
		"Set file to the exact changed file path from static.diff_snippets and line to a changed line number inside that hunk. If a recommendation cannot be tied to a specific changed file, omit file and line.",
		"Set source_labels to the labels of context snippets or source_refs you actually relied on (e.g. R1, L2). Omit labels you did not use.",
		"Only when review_profile is pr_summary: include a top-level overview field, 2-3 sentences on what this change does and why, based on the revision descriptions, session_transcript context, and session_context intent/edit trail; no file lists, no URLs, no praise. For all other profiles, omit overview.",
		"Only when review_profile is pr_summary: include notable_changes — 3 to 6 entries, each the single most important changed line of one logical change. file must be an exact changed file path from static.diff_snippets and line a changed line inside that hunk. note is one sentence describing what changed and why it matters, no file paths, no URLs. Omit entries you cannot anchor. For all other profiles, omit notable_changes.",
		"Only when review_profile is pr_summary: include downstream_impact — 1 to 3 sentences on customer-facing risk (could this introduce bugs or issues for customers?) and how the change shifts the status quo of the codebase or application, including potential downstream effects. Calibrate depth to diff size: tiny localized changes get one brief sentence (e.g. low risk to existing behavior); large multi-area changes get a broader assessment. No file lists, no URLs, no praise. For all other profiles, omit downstream_impact.",
		"pr_summary behaves like patch_focused for finding selection (current-change review, changed-lines evidence, same rejection rules — no quota-filling, no generic advice) plus the overview and downstream_impact rules.",
		"Return JSON only with shape {\"overview\":string(optional),\"downstream_impact\":string(optional),\"notable_changes\":[{\"file\":string,\"line\":number,\"note\":string}](optional),\"recommendations\":[{\"title\":string,\"summary\":string,\"benefit\":string,\"recommendation\":string,\"strength\":\"Strong|Worth exploring|Speculative\",\"evidence\":[string],\"file\":string(optional),\"line\":number(optional),\"source_labels\":[string](optional)}]}.",
		"Return at most 5 recommendations. Prefer 2-3 high-signal recommendations.",
	}
}

func mustJSON(value any) string {
	body, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(body)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func aiReviewRequestedFromEnv() bool {
	value := strings.TrimSpace(os.Getenv("GX_REVIEW_AI"))
	if value == "" {
		return true
	}
	if strings.EqualFold(value, "0") || strings.EqualFold(value, "false") || strings.EqualFold(value, "no") {
		return false
	}
	return true
}

func formatReviewerDegradation(err error) string {
	if err == nil {
		return "unknown error"
	}
	msg := strings.TrimSpace(err.Error())
	if msg == "" {
		return "unknown error"
	}
	if len(msg) > 160 {
		return msg[:157] + "..."
	}
	return msg
}
