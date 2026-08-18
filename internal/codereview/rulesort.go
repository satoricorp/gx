package codereview

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Rule sorting is a post-verification labeling pass, and where it sits is the
// whole design.
//
// The obvious placement — hand the reviewer the rule pack and ask it to find
// instances — was measured five separate ways: blended into the reviewer
// prompt, weighted, split into sub-rules, given full rule text plus an explicit
// checklist, and sharded across five parallel legs. Every one scored below the
// no-rules control (30, 34, 35.4, 35.8, 27.2 against 44 F1). The mechanism was
// consistent: rules in the reviewer's context compete for attention with
// finding the bug, and a leg assigned rules it cannot satisfy invents findings
// to fill them.
//
// So detection is not told the rules exist. The sorter runs after the judge has
// already confirmed which findings a human will see, and its only output is a
// name for something that was found without it. That placement cannot cost
// recall, because nothing upstream of it changes, and it cannot cost precision,
// because it never adds, removes, or edits a finding — see applyRuleSort, which
// writes exactly one field.
const (
	// ruleSortEnvVar disables the pass, matching GX_REVIEW_JUDGE=0. Labeling is
	// on by default because an unlabeled finding is what shipped before this
	// existed: turning it off degrades to the old output, never to a wrong one.
	ruleSortEnvVar      = "GX_REVIEW_RULE_SORT"
	ruleSortModelEnvVar = "GX_REVIEW_RULE_SORT_MODEL"
	// defaultRuleSortModel is the cheap tier on purpose. Sorting an existing,
	// already-verified finding into a named bucket is a much smaller job than
	// finding it or checking it, and it runs off the critical path of both.
	defaultRuleSortModel = "us.anthropic.claude-haiku-4-5-20251001-v1:0"
	// defaultRuleSortMaxOutputTokens is sized for a short reason plus a slug per
	// finding, not for prose.
	defaultRuleSortMaxOutputTokens = 8000
	defaultRuleSortBatchSize       = 12
	maxConcurrentRuleSortBatches   = 4
)

// RuleSorter names which rule each confirmed finding instances. It is a
// separate interface from FindingJudge because the two must not be merged: the
// judge decides what a human sees, and adding a labeling field to its output
// schema would trade a verification stage that measurably works for a labeling
// convenience.
type RuleSorter interface {
	Available() bool
	Sort(ctx context.Context, req ruleSortRequest) ([]ruleSortResult, error)
}

type bedrockRuleSorter struct {
	client *bedrockAnthropicReviewer
}

type unavailableRuleSorter struct {
	reason string
}

// ruleSortRule is one bucket the sorter may choose. Advisory rides along
// because the pack already distinguishes rules that inform from rules that
// block, and a labeled advisory finding can render differently at no extra
// cost.
type ruleSortRule struct {
	ID string `json:"rule_id"`
	// Text is the rule's full guidance, not policyRuleSummary's one-liner.
	//
	// That helper exists to put a rule beside its ID in a brief, and it caps at
	// 160 bytes mid-word: boolean-polarity reached the sorter as "...negating a
	// condition that should not be negated, o" and variable-misuse as "...not a
	// different, same-typed val", both losing exactly the examples that decide
	// whether a finding is an instance. Measured on 27 findings, no rule that
	// names a specific mechanism was ever chosen while the catalog was built
	// this way — the sorter was picking between truncated fragments, and the
	// general rules survive truncation better because their first clause is
	// already the whole rule.
	//
	// The parser caps rule text at maxReviewRuleBytes, so the catalog is bounded
	// by the pack, and this is the cheapest call in the pipeline.
	Text     string `json:"guidance"`
	Advisory bool   `json:"advisory,omitempty"`
}

// ruleSortCandidate is a confirmed finding as the sorter sees it. It carries no
// judge verdict, severity, or lane: those are decided, and showing them to a
// model whose only job is naming a rule invites it to relitigate them.
type ruleSortCandidate struct {
	ID             string `json:"candidate_id"`
	Title          string `json:"title"`
	Summary        string `json:"summary,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
	File           string `json:"file,omitempty"`
	Line           int    `json:"line,omitempty"`
}

type ruleSortRequest struct {
	Rules      []ruleSortRule        `json:"rules"`
	Files      []judgeContentSnippet `json:"files,omitempty"`
	Diffs      []judgeDiffSnippet    `json:"diffs,omitempty"`
	Candidates []ruleSortCandidate   `json:"candidates"`
}

type ruleSortResponse struct {
	Results []ruleSortResult `json:"results"`
}

type ruleSortResult struct {
	CandidateID string `json:"candidate_id"`
	// Mechanism and Considered are read into fields nothing renders, for the
	// same reason the judge keeps Analysis: a JSON object is emitted key by key,
	// so what is written before rule_id is what the choice is conditioned on.
	//
	// They are two separate fields because they fix two separate failures
	// measured on 27 real findings. Naming the mechanism in the model's own
	// words, before any rule is in view, stops the choice being driven by which
	// rule name shares vocabulary with the finding's title — that is how an
	// inverted boolean came back as dont-repeat-yourself. Listing the plausible
	// rules before picking one stops the first adequate match from winning:
	// without it, none of the eight rules that name a specific defect mechanism
	// was ever selected across 27 findings, because a general rule always came
	// up first and nothing forced a comparison.
	Mechanism  string   `json:"mechanism"`
	Considered []string `json:"considered"`
	RuleID     string   `json:"rule_id"`
}

func ruleSorterFromEnv() RuleSorter {
	if ruleSortDisabledFromEnv() {
		return nil
	}
	plan, err := resolveBedrockTransportPlan()
	if err != nil {
		return unavailableRuleSorter{reason: err.Error()}
	}
	return bedrockRuleSorter{client: newBedrockReviewer(plan.newTransport(), resolveRuleSortModel())}
}

func ruleSortDisabledFromEnv() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv(ruleSortEnvVar)), "0")
}

// resolveRuleSortModel applies env > default. Like the judge, the reviewed
// repository gets no say in it.
func resolveRuleSortModel() string {
	return normalizeBedrockModelID(firstNonEmpty(
		os.Getenv(ruleSortModelEnvVar),
		defaultRuleSortModel,
	))
}

func resolveRuleSortBatchSize() int {
	if raw := strings.TrimSpace(os.Getenv("GX_REVIEW_RULE_SORT_BATCH")); raw != "" {
		if size, err := strconv.Atoi(raw); err == nil && size > 0 {
			return size
		}
	}
	return defaultRuleSortBatchSize
}

func (s bedrockRuleSorter) Available() bool { return s.client != nil }

func (s unavailableRuleSorter) Available() bool { return false }

func (s unavailableRuleSorter) Sort(context.Context, ruleSortRequest) ([]ruleSortResult, error) {
	return nil, fmt.Errorf("rule sorter is unavailable: %s", s.reason)
}

func (s bedrockRuleSorter) Sort(ctx context.Context, req ruleSortRequest) ([]ruleSortResult, error) {
	if s.client == nil {
		return nil, fmt.Errorf("rule sorter is not configured")
	}
	completion, err := s.client.completeJSON(ctx, ruleSortDeveloperPrompt(), mustJSON(req), defaultRuleSortMaxOutputTokens)
	if err != nil {
		return nil, err
	}
	return parseRuleSortResponse(completion.Text)
}

// parseRuleSortResponse reuses the judge's tolerant object picker: same
// transport, same models, same habit of wrapping JSON in prose.
func parseRuleSortResponse(content string) ([]ruleSortResult, error) {
	object := pickAnsweringJSONObject(content, func(fields map[string]json.RawMessage) int {
		if _, ok := fields["results"]; !ok {
			return -1
		}
		return 1
	})
	if object == "" {
		return nil, fmt.Errorf("decode rule-sort JSON: no complete JSON object with a results field in response")
	}
	var parsed ruleSortResponse
	if err := json.Unmarshal([]byte(object), &parsed); err != nil {
		return nil, fmt.Errorf("decode rule-sort JSON: %w", err)
	}
	return parsed.Results, nil
}

// ruleSortCatalog is the set of buckets in force: the built-in pack first, then
// any rule the repository's REVIEW.md added that the pack does not already
// name.
//
// The `review/*` defect classes are deliberately absent. They are the eight
// open-ended judgment categories the reviewer prompt used to carry, and two of
// them (wrong-logic, data-correctness) are catch-alls that took 24 of 64
// findings in one run and produced 2 true positives. The pack replaces them
// with the specific rules those findings should have landed on, so including
// both would hand the sorter a wide bucket next to the narrow ones and ask it
// to prefer narrow.
func ruleSortCatalog(ctx ReviewContext) ([]ruleSortRule, error) {
	var catalog []ruleSortRule
	seen := map[string]struct{}{}
	// A pack that does not parse is a build-time mistake, and pack.go surfaces
	// it rather than reviewing with fewer rules. Swallowing it here would do
	// exactly what that comment forbids one layer down: the sort would quietly
	// run against only whatever REVIEW.md declared, and a review labeled from a
	// two-rule catalog is indistinguishable in the output from one labeled from
	// the full pack.
	pack, err := RecommendedPack()
	if err != nil {
		return nil, err
	}
	for _, rule := range pack {
		catalog = append(catalog, ruleSortRule{
			ID:       rule.ID,
			Text:     rule.Text,
			Advisory: rule.Advisory,
		})
		seen[strings.ToLower(rule.ID)] = struct{}{}
	}
	for _, def := range ctx.Brief.Rules {
		id := strings.TrimSpace(def.ID)
		if id == "" {
			continue
		}
		if _, dup := seen[strings.ToLower(id)]; dup {
			continue
		}
		seen[strings.ToLower(id)] = struct{}{}
		catalog = append(catalog, ruleSortRule{ID: id, Text: def.Summary})
	}
	return catalog, nil
}

// ruleSortAllowlist is the catalog keyed for NormalizeRuleID: lowercased name
// to canonical spelling. Built from the catalog rather than KnownRuleIDs so
// that a slug the sorter invents, and a `review/*` class it remembers from
// somewhere else, both fold to "" instead of minting a rule.
func ruleSortAllowlist(catalog []ruleSortRule) map[string]string {
	known := make(map[string]string, len(catalog))
	for _, rule := range catalog {
		known[strings.ToLower(rule.ID)] = rule.ID
	}
	for retired, current := range RecommendedPackAliases() {
		if canonical, ok := known[strings.ToLower(current)]; ok {
			known[strings.ToLower(retired)] = canonical
		}
	}
	return known
}

func ruleSortDeveloperPrompt() string {
	return strings.Join([]string{
		"You are labeling code-review findings that have already been found and already been verified. Your only job is to name which rule each finding is an instance of.",
		"You are not reviewing the code. You are not judging whether the finding is correct — that decision is made and is not yours to revisit. Do not rewrite, re-title, re-scope, or re-rank anything.",
		"Answer three fields per candidate, in this order, and do not decide the third before writing the first two.",
		"mechanism: one short clause naming what the code actually does wrong, in your own words, using none of the rule names. Describe the defect, not the finding's phrasing. A finding's title is written to be read by a human, not to match a rule, and matching on its vocabulary is the single most common way to get this wrong.",
		"considered: every rule_id that could plausibly cover that mechanism, narrowest first. Read the entire rules list before answering this — a rule near the end of the list is exactly as eligible as one near the start. Use [] when nothing plausibly applies.",
		"rule_id: the narrowest entry in considered that genuinely names the mechanism, copied exactly. Use \"\" when considered is empty, or when nothing in it survives the test below.",
		"Some rules name a general area (duplication, dead code, comments, scope, silent regressions). Others name one specific defect mechanism: a variable used where another was meant, an inverted boolean, a comparison made without normalizing first, an off-by-one or boundary error, a value that is present but falsy, an operation applied to some cases and not the rest, a test that races what it asserts, an identifier that names the wrong thing. When the mechanism you wrote is one of those, the specific rule wins — always, and even when a general rule also fits. A general rule is correct only when no specific one names the mechanism.",
		"The test, applied to your chosen rule: would an engineer who knows this rule, reading its name printed beside this finding, agree the rule names what went wrong? Not 'is it related' — does it name it. If the honest answer is no, return \"\".",
		"There is no catch-all rule and you must not press any rule into that role. Returning \"\" is a correct, expected, and frequent outcome: roughly one real review comment in ten fits no rule at all, and a review where every finding got a label is a review where some labels are wrong. A wrong rule name is worse than no rule name, because the entire value of a named rule is that it means something specific.",
		"Return JSON only, with shape {\"results\":[{\"candidate_id\":string,\"mechanism\":string,\"considered\":[string],\"rule_id\":string}]}. Include every candidate_id you were given, exactly once.",
	}, "\n")
}

// buildRuleSortRequest gives the sorter the finding, the rules, and the code
// the finding points at. The code matters: several pack rules are separated by
// what the line actually does rather than by how the finding is worded, and a
// sorter working from the finding text alone has to guess between them.
func buildRuleSortRequest(ctx ReviewContext, catalog []ruleSortRule, findings []Finding) ruleSortRequest {
	candidates := make([]ruleSortCandidate, 0, len(findings))
	var ordered []string
	seen := map[string]struct{}{}
	hints := map[string][]int{}
	for _, finding := range findings {
		candidates = append(candidates, ruleSortCandidate{
			ID:             finding.ID,
			Title:          finding.Title,
			Summary:        finding.Summary,
			Recommendation: finding.Recommendation,
			File:           finding.File,
			Line:           finding.Line,
		})
		for _, file := range namedFilesForFinding(ctx, finding) {
			if _, ok := seen[file]; !ok {
				seen[file] = struct{}{}
				ordered = append(ordered, file)
			}
		}
		for file, lines := range findingLineHints(finding) {
			hints[file] = append(hints[file], lines...)
		}
	}
	return ruleSortRequest{
		Rules:      catalog,
		Files:      judgeFileContentSnippets(ctx.Brief.RepoRoot, ordered, hints),
		Diffs:      judgeDiffSnippets(ctx.Brief.Static.DiffSnippets, ordered),
		Candidates: candidates,
	}
}

// ruleSortOutcome reports what the pass managed to do, so a review whose
// labeling failed does not read as a review whose findings fit no rule.
type ruleSortOutcome struct {
	Labeled   int
	Unlabeled int
	Batches   int
	Failed    int
	Err       error
}

// runRuleSort labels findings in place and returns what happened. It is
// deliberately total: a batch that fails leaves its findings unlabeled, which
// is exactly the pre-existing output, so there is no failure mode in which this
// pass makes a review worse than not running it.
func runRuleSort(ctx context.Context, sorter RuleSorter, reviewContext ReviewContext, findings []Finding) ruleSortOutcome {
	if sorter == nil || !sorter.Available() || len(findings) == 0 {
		return ruleSortOutcome{}
	}
	catalog, err := ruleSortCatalog(reviewContext)
	if err != nil {
		return ruleSortOutcome{Batches: 1, Failed: 1, Err: err}
	}
	if len(catalog) == 0 {
		return ruleSortOutcome{}
	}
	known := ruleSortAllowlist(catalog)

	batchSize := resolveRuleSortBatchSize()
	var batches [][]Finding
	for start := 0; start < len(findings); start += batchSize {
		end := start + batchSize
		if end > len(findings) {
			end = len(findings)
		}
		batches = append(batches, findings[start:end])
	}

	type batchResult struct {
		results []ruleSortResult
		err     error
	}
	out := make([]batchResult, len(batches))
	slots := make(chan struct{}, maxConcurrentRuleSortBatches)
	var wg sync.WaitGroup
	for i, batch := range batches {
		i, batch := i, batch
		wg.Add(1)
		go func() {
			defer wg.Done()
			slots <- struct{}{}
			defer func() { <-slots }()
			results, err := sorter.Sort(ctx, buildRuleSortRequest(reviewContext, catalog, batch))
			out[i] = batchResult{results: results, err: err}
		}()
	}
	wg.Wait()

	outcome := ruleSortOutcome{Batches: len(batches)}
	byID := map[string]string{}
	for i := range batches {
		if out[i].err != nil {
			outcome.Failed++
			if outcome.Err == nil {
				outcome.Err = out[i].err
			}
			continue
		}
		for _, result := range out[i].results {
			// NormalizeRuleID against the catalog is the guardrail that makes
			// "never invent a rule" mechanical rather than a request in a
			// prompt: anything not in the pack folds to "" here.
			if id := NormalizeRuleID(result.RuleID, known); id != "" {
				byID[result.CandidateID] = id
			}
		}
	}
	applyRuleSort(findings, byID)
	for i := range findings {
		if findings[i].RuleID != "" {
			outcome.Labeled++
		} else {
			outcome.Unlabeled++
		}
	}
	return outcome
}

// applyRuleSort writes the rule name onto each finding and writes nothing else.
//
// This function is the entire enforcement of "labeling must not change the
// finding". Every other field is untouched by construction rather than by
// review: if a later change needs the sorter to affect severity, ordering, or
// text, it does not belong in this pass at all, because a label that can edit
// what it labels is detection wearing a different hat.
func applyRuleSort(findings []Finding, byID map[string]string) {
	if len(byID) == 0 {
		return
	}
	for i := range findings {
		if id, ok := byID[findings[i].ID]; ok {
			findings[i].RuleID = id
		}
	}
}
