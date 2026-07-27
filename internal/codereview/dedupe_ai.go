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

// The model that decides whether two candidate findings are one finding.
//
// See dedupe.go for why this cannot be a similarity threshold. In short:
// duplication is agreement about a root cause, and the tokens that separate one
// root cause from another ("Close" from "Sync") are precisely the ones a
// set-overlap score cannot weigh. A model can read both findings and answer the
// question that actually matters — would one fix close both? — so it does.
//
// It runs on the judge's model rather than a reviewer's: this is a small,
// mechanical classification over text that is already in hand, not a code
// review, and it sits on the critical path before verification. Same transport
// precedence as everything else, so a review cannot end up with its reviewers on
// one wire and its de-duplication on another with no way to tell.

const defaultDedupeMaxOutputTokens = 4000

// duplicateAdjudicator answers "same defect or two defects?" for gated pairs.
type duplicateAdjudicator interface {
	AdjudicateDuplicates(ctx context.Context, pairs []duplicatePairInput) ([]duplicateVerdict, error)
}

type duplicatePairInput struct {
	PairID string            `json:"pair_id"`
	Left   duplicateFindingT `json:"finding_a"`
	Right  duplicateFindingT `json:"finding_b"`
}

type duplicateFindingT struct {
	Title          string   `json:"title"`
	Summary        string   `json:"summary,omitempty"`
	Recommendation string   `json:"recommendation,omitempty"`
	Files          []string `json:"files,omitempty"`
}

type duplicateVerdict struct {
	PairID string `json:"pair_id"`
	Same   bool   `json:"same_defect"`
	Reason string `json:"reason,omitempty"`
}

type duplicateAdjudicationResponse struct {
	Results []duplicateVerdict `json:"results"`
}

type bedrockDuplicateAdjudicator struct {
	client *bedrockAnthropicReviewer
}

// duplicateAdjudicatorFromEnvWithPolicy builds the adjudicator, or nil when
// there is none to build. nil is a supported state, not an error: the caller
// falls back to merging near-verbatim copies only and reports that it did.
func duplicateAdjudicatorFromEnvWithPolicy(policy *ReviewPolicy) duplicateAdjudicator {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_DEDUPE")), "0") {
		return nil
	}
	plan, err := resolveBedrockTransportPlan()
	if err != nil {
		return nil
	}
	return bedrockDuplicateAdjudicator{client: newBedrockReviewer(plan.newTransport(), resolveBedrockJudgeModel(policy))}
}

func (a bedrockDuplicateAdjudicator) AdjudicateDuplicates(ctx context.Context, pairs []duplicatePairInput) ([]duplicateVerdict, error) {
	if a.client == nil {
		return nil, fmt.Errorf("duplicate adjudicator is not configured")
	}
	if len(pairs) == 0 {
		return nil, nil
	}
	request := struct {
		Pairs []duplicatePairInput `json:"pairs"`
	}{Pairs: pairs}
	content, err := a.client.completeJSON(ctx, duplicateAdjudicatorPrompt(), mustJSON(request), defaultDedupeMaxOutputTokens)
	if err != nil {
		return nil, err
	}
	return parseDuplicateAdjudication(content)
}

// duplicateAdjudicatorPrompt states the decision, the asymmetry of the two
// errors, and — because the failure mode is specifically "same file, same
// vocabulary, different bug" — worked examples of pairs that must not merge.
func duplicateAdjudicatorPrompt() string {
	return strings.Join([]string{
		"You are GX Review De-duplicator. Two AI reviewers reviewed the same code independently and their findings are about to be shown to an engineer.",
		"For each pair, decide whether finding_a and finding_b report THE SAME DEFECT, so that one fix resolves both.",
		"Return JSON only, with shape {\"results\":[{\"pair_id\":string,\"same_defect\":true|false,\"reason\":string}]}. Return one result per pair, using the pair_id given.",
		"SAME DEFECT means the same root cause in the same code: fixing one fixes the other. Different wording, different emphasis, a different proposed remedy for the same underlying problem, or a different file cited as the place to fix it are all still the same defect.",
		"DIFFERENT DEFECTS means distinct root causes, even when the two findings sit in the same file, share most of their vocabulary, or belong to the same theme. These must NOT be merged:",
		"- an ignored error from Close and an ignored error from Sync in the same writer,",
		"- a permissions race when a secret file is created and that same file being exposed through a bind mount,",
		"- a missing authorization check and a missing request timeout in the same handler,",
		"- a missing test for a code path and a bug in that same code path,",
		"- a retry loop with no jitter and a retry loop that retries non-idempotent requests.",
		"Judge the defect, not the prose. Two findings can share a topic sentence and be two bugs; two findings can share almost no words and be one bug.",
		"When you are not sure, answer false. The errors are not symmetric: a missed merge leaves a near-duplicate in the review, while a wrong merge deletes a real finding that the engineer will then never see.",
		"Keep each reason to one short sentence naming the root cause you compared.",
	}, "\n")
}

func parseDuplicateAdjudication(content string) ([]duplicateVerdict, error) {
	var parsed duplicateAdjudicationResponse
	if err := json.Unmarshal([]byte(content), &parsed); err == nil {
		return parsed.Results, nil
	}
	trimmed := extractJSONObject(content)
	if trimmed == "" {
		return nil, fmt.Errorf("decode de-duplication JSON: no JSON object in response")
	}
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil, fmt.Errorf("decode de-duplication JSON: %w", err)
	}
	return parsed.Results, nil
}

// adjudicateDuplicatePairs asks about every gated pair and returns the merges it
// was told to make, plus how many pairs actually got an answer.
//
// A pair with no answer is not merged. That is the safe direction, but it is not
// a free one — the review then ships the duplicates it was built to collapse —
// so the count comes back to the caller and reaches the reader as a degraded
// reason instead of disappearing.
func adjudicateDuplicatePairs(ctx context.Context, adjudicator duplicateAdjudicator, findings []Finding, pairs []pairCandidate) (map[string]bool, int, error) {
	if adjudicator == nil || len(pairs) == 0 {
		return nil, 0, nil
	}
	inputs := make([]duplicatePairInput, 0, len(pairs))
	byPairID := map[string]pairCandidate{}
	for _, pair := range pairs {
		id := fmt.Sprintf("p%d-%d", pair.left, pair.right)
		byPairID[id] = pair
		inputs = append(inputs, duplicatePairInput{
			PairID: id,
			Left:   duplicateAdjudicatorFinding(findings[pair.left]),
			Right:  duplicateAdjudicatorFinding(findings[pair.right]),
		})
	}

	var batches [][]duplicatePairInput
	for start := 0; start < len(inputs); start += dedupeAdjudicationBatch {
		end := start + dedupeAdjudicationBatch
		if end > len(inputs) {
			end = len(inputs)
		}
		batches = append(batches, inputs[start:end])
	}
	type batchResult struct {
		verdicts []duplicateVerdict
		err      error
	}
	results := make([]batchResult, len(batches))
	var wg sync.WaitGroup
	for i, batch := range batches {
		i, batch := i, batch
		wg.Add(1)
		go func() {
			defer wg.Done()
			verdicts, err := adjudicator.AdjudicateDuplicates(ctx, batch)
			results[i] = batchResult{verdicts: verdicts, err: err}
		}()
	}
	wg.Wait()

	merges := map[string]bool{}
	answered := map[string]struct{}{}
	var firstErr error
	for _, result := range results {
		if result.err != nil {
			if firstErr == nil {
				firstErr = result.err
			}
			continue
		}
		for _, verdict := range result.verdicts {
			pair, ok := byPairID[strings.TrimSpace(verdict.PairID)]
			if !ok {
				continue // a pair_id nobody asked about decides nothing
			}
			answered[strings.TrimSpace(verdict.PairID)] = struct{}{}
			if verdict.Same {
				merges[duplicatePairKey(findings[pair.left].ID, findings[pair.right].ID)] = true
			}
		}
	}
	return merges, len(answered), firstErr
}

func duplicateAdjudicatorFinding(finding Finding) duplicateFindingT {
	return duplicateFindingT{
		Title:          strings.TrimSpace(finding.Title),
		Summary:        truncateForAdjudication(finding.Summary, 900),
		Recommendation: truncateForAdjudication(finding.Recommendation, 700),
		Files:          sortedSet(findingFiles(finding)),
	}
}

func truncateForAdjudication(text string, limit int) string {
	text = strings.TrimSpace(text)
	if len(text) <= limit {
		return text
	}
	return text[:limit] + "…"
}

// duplicatePairKey is order-independent so a verdict recorded for (a,b) is
// found when the grouping loop asks about (b,a).
func duplicatePairKey(left, right string) string {
	ids := []string{left, right}
	sort.Strings(ids)
	return ids[0] + "\x00" + ids[1]
}
