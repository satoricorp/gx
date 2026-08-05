package codereview

import (
	"context"
	"fmt"
	"path"
	"sort"
	"strings"
)

// Cross-reviewer finding de-duplication.
//
// The review panel is two flagship models plus a sharded subject, so the same
// problem reaches the judge several times over: once per reviewer leg that saw
// it, and once per shard that contained the relevant file. Nothing merged those
// copies. A live whole-repo review of a small repository produced three separate
// findings — "Unrestricted guest egress makes the documented secrets model
// weaker than it appears", "Unrestricted guest egress enables secret
// exfiltration with no mitigation timeline", and "Unrestricted guest egress
// makes run.env exfiltratable despite 0600 permissions" — which are one issue in
// three voices. The judge cannot fix this: it is delete-only, it rules on each
// candidate in isolation, and it happily confirms all three.
//
// Duplication is also an opportunity, not just noise. Two independently prompted
// flagship models converging on the same problem is the single strongest quality
// signal this design produces, and it was being discarded. Merging records which
// legs agreed (Finding.Corroboration), which then shows in the rendered review,
// the JSON report, and the ranking.
//
// # Why a similarity threshold decides nothing here
//
// The first attempt at this was a lexical score — token overlap on title and on
// title+summary+recommendation — with a threshold calibrated on one live run.
// It does not work, and the reason is worth recording so nobody re-derives it.
//
// Score every pair of those 32 candidates with findingSimilarity and the true
// duplicates and the true distinct pairs interleave, within that one batch:
//
//	0.3526  DISTINCT   "No test coverage for the event log replay migration path"
//	                 / "Event log replay loses touchedFrozenTests after the first iteration"
//	0.3269  DISTINCT   "Secrets file may be written world-readable on creation race"
//	                 / "Secret file store has no file-locking"
//	0.2976  duplicate  "Secret file store has no file-locking"
//	                 / "secrets.ts keys file has a TOCTOU between readAll and writeAll"
//	0.2632  duplicate  "cli.ts conflates parsing, orchestration, and presentation"
//	                 / "cli.ts is a 1300+ line Module with no internal seam"
//
// A missing test for a path and a bug in that path are not one finding; two
// descriptions of one read-modify-write race are. No cut separates those four,
// so no cut separates the population they come from. The fixtures in
// dedupe_test.go make the same point with wider margins and are checked on every
// run by TestNoLexicalThresholdSeparatesDuplicatesFromDistinctFindings: "Errors
// from Close are ignored in the writer" and "Errors from Sync are ignored in the
// writer" score 0.71 on the same file and are two different bugs, while the two
// egress findings a live review shipped side by side — which a human reads as
// one finding at a glance — score 0.34.
//
// Bag-of-words similarity measures shared topic. Duplication is shared *root
// cause*, and the tokens that distinguish one root cause from another ("Close"
// from "Sync") are exactly the ones a set-overlap score cannot weigh.
//
// # What decides instead
//
// Two stages, and the second one is the decision:
//
//  1. A deterministic gate (file overlap, then dedupeGateThreshold on the
//     lexical score) narrows n^2 pairs to the few worth asking about. It is
//     tuned for recall only — a pair it lets through is not a duplicate, it is a
//     question — so shared boilerplate inflating a score costs one adjudication
//     and nothing else.
//  2. A model adjudicates each gated pair: same defect, or two defects? That is
//     a semantic question about root cause, which is what the model is for.
//
// Both stages must agree before anything is merged, and the adjudicator is
// instructed to answer "different" when unsure, because the two errors are not
// symmetric: a missed merge leaves a near-duplicate in the review, which is
// untidy, while a wrong merge deletes a finding a human never sees — a review
// that silently failed to report a problem it had already found. Grouping
// compares each candidate against its group's seed rather than a growing union,
// so A~B and B~C cannot drag in a C that matches neither.
//
// With no adjudicator reachable (no credentials, LGTM_REVIEW_DEDUPE=0, or a failed
// call) nothing semantic is merged. The fallback collapses only near-verbatim
// copies — dedupeVerbatimThreshold on both title and body — which is the old
// behavior and is safe because it is barely a judgement at all. A review that
// could not adjudicate says so in its degraded reasons rather than shipping a
// deduplicated-looking result.
const (
	// dedupeGateThreshold is the lexical score a pair needs before the
	// adjudicator is asked about it. It is a recall floor, not a decision, and
	// it is set by what it must never drop rather than by where the classes
	// separate — they do not separate.
	//
	// Measured over the 32 raw candidates of a live two-leg whole-repo review:
	// 496 pairs, 124 of which share a file, 24 of which clear 0.20. The lowest
	// hand-labelled duplicate in that batch scores 0.2632 ("cli.ts conflates
	// parsing…" against "cli.ts is a 1300+ line Module with no internal seam"),
	// and the lowest across the dedupe_test.go fixtures is 0.3438. So 0.20
	// leaves roughly a quarter of margin under the worst real duplicate seen,
	// and costs 24 yes/no questions in one batched call rather than 496.
	dedupeGateThreshold = 0.20

	// dedupeVerbatimThreshold is the fallback used when no adjudicator answered.
	// Both title and body must clear it, which in practice means one leg emitted
	// the same finding twice. It is deliberately near-verbatim: with no semantic
	// judgement available, only an obvious copy is safe to collapse.
	dedupeVerbatimThreshold = 0.85

	// dedupeMinTextOnlyTokens is how many content tokens the shorter of two
	// file-less findings must have before prose alone is allowed to merge them.
	dedupeMinTextOnlyTokens = 4

	// dedupeAdjudicationBatch is how many pairs one adjudication call carries.
	// Like the judge, the binding constraint is output length, and the answer
	// per pair is small; batches run concurrently.
	dedupeAdjudicationBatch = 25

	// maxAdjudicatedPairs bounds the work when a review produces an unusual
	// number of overlapping findings. Pairs are offered highest-score-first, so
	// the cap drops the least similar questions; anything dropped is reported as
	// unresolved rather than silently answered "different".
	maxAdjudicatedPairs = 200
)

// findingSignature is everything the deterministic gate compares.
type findingSignature struct {
	// files is where the finding says the problem is, gathered from the
	// finding's own File/Anchors and from every path it names in prose.
	files map[string]struct{}
	// title and body are normalized token sets. They are scored separately
	// because a title is a compressed claim and a summary is an argument; a
	// shared title phrase is worth more per token than a shared summary phrase.
	title map[string]struct{}
	body  map[string]struct{}
}

// dedupeOutcome is what one de-duplication pass did, including what it could not
// decide. The engine turns the undecided part into a degraded reason: "we did
// not merge duplicates" and "there were none" must not look the same.
type dedupeOutcome struct {
	Findings []Finding
	// Gated is how many candidate pairs cleared the deterministic gate.
	Gated int
	// Adjudicated is how many of those got an answer.
	Adjudicated int
	// Merged is how many findings were folded into another.
	Merged int
	// Err is the first adjudication failure, if any.
	Err error
}

// Unresolved is how many gated pairs never got an answer — because the
// adjudicator was absent, failed, or the pair fell past maxAdjudicatedPairs.
func (o dedupeOutcome) Unresolved() int {
	if o.Gated <= o.Adjudicated {
		return 0
	}
	return o.Gated - o.Adjudicated
}

// mergeNearDuplicateFindings collapses findings that describe the same problem
// and records, on the survivor, which reviewer legs raised it.
func mergeNearDuplicateFindings(ctx context.Context, adjudicator duplicateAdjudicator, findings []Finding) []Finding {
	return dedupeFindings(ctx, adjudicator, findings).Findings
}

func dedupeFindings(ctx context.Context, adjudicator duplicateAdjudicator, findings []Finding) dedupeOutcome {
	outcome := dedupeOutcome{Findings: findings}
	if len(findings) < 2 {
		return outcome
	}
	signatures := make([]findingSignature, len(findings))
	for i, finding := range findings {
		signatures[i] = newFindingSignature(finding)
	}

	pairs := gatedFindingPairs(findings, signatures)
	outcome.Gated = len(pairs)
	if len(pairs) > maxAdjudicatedPairs {
		pairs = pairs[:maxAdjudicatedPairs]
	}
	verdicts, adjudicated, err := adjudicateDuplicatePairs(ctx, adjudicator, findings, pairs)
	outcome.Adjudicated = adjudicated
	outcome.Err = err

	same := func(left, right int) bool {
		if verbatimDuplicate(signatures[left], signatures[right]) {
			return true
		}
		return verdicts[duplicatePairKey(findings[left].ID, findings[right].ID)]
	}

	// Seed-anchored grouping. Each candidate is compared against the member that
	// opened the group, never against a growing union of everything already in
	// it, so similarity cannot chain a third finding in through a second.
	var groups [][]int
	for i := range findings {
		placed := false
		for g := range groups {
			if same(groups[g][0], i) {
				groups[g] = append(groups[g], i)
				placed = true
				break
			}
		}
		if !placed {
			groups = append(groups, []int{i})
		}
	}

	out := make([]Finding, 0, len(groups))
	for _, members := range groups {
		outcome.Merged += len(members) - 1
		out = append(out, collapseFindingGroup(findings, members))
	}
	outcome.Findings = out
	return outcome
}

// pairCandidate is one question for the adjudicator.
type pairCandidate struct {
	left  int
	right int
	score float64
}

// gatedFindingPairs is the deterministic recall filter, highest score first.
func gatedFindingPairs(findings []Finding, signatures []findingSignature) []pairCandidate {
	var pairs []pairCandidate
	for i := 0; i < len(findings); i++ {
		for j := i + 1; j < len(findings); j++ {
			if verbatimDuplicate(signatures[i], signatures[j]) {
				continue // decided without asking
			}
			if !sharesAny(signatures[i].files, signatures[j].files) {
				continue
			}
			score := findingSimilarity(signatures[i], signatures[j])
			if score < dedupeGateThreshold {
				continue
			}
			pairs = append(pairs, pairCandidate{left: i, right: j, score: score})
		}
	}
	sort.SliceStable(pairs, func(a, b int) bool {
		if pairs[a].score != pairs[b].score {
			return pairs[a].score > pairs[b].score
		}
		if pairs[a].left != pairs[b].left {
			return pairs[a].left < pairs[b].left
		}
		return pairs[a].right < pairs[b].right
	})
	return pairs
}

// verbatimDuplicate is the only merge the deterministic side makes on its own.
//
// The file gate is a precondition here, not a weighted term. A code review
// finding is a claim about a place in the code; if two findings name no file in
// common they are claims about different places. Locality alone is nowhere near
// sufficient, which the live run shows plainly: src/bridge.ts attracted three
// unrelated findings (a re-entrancy race, a blocking synchronous read, and an
// unbounded retry on a malformed file), two of which cite the same line.
func verbatimDuplicate(left, right findingSignature) bool {
	if len(left.files) == 0 && len(right.files) == 0 {
		// No locality evidence on either side, so there is nothing to gate on
		// and nothing to send an adjudicator that it could check.
		if smallerSet(left.body, right.body) < dedupeMinTextOnlyTokens {
			// Too little text to have an opinion. An overlap coefficient over
			// two-token sets is 1.0 as soon as one token matches, which is a
			// coin flip dressed up as a measurement.
			return false
		}
		return overlapCoefficient(left.body, right.body) >= dedupeVerbatimThreshold
	}
	if !sharesAny(left.files, right.files) {
		return false
	}
	return overlapCoefficient(left.title, right.title) >= dedupeVerbatimThreshold &&
		overlapCoefficient(left.body, right.body) >= dedupeVerbatimThreshold
}

// findingSimilarity scores two findings' topical agreement in [0,1]. It is the
// gate's ranking function and nothing more — see the package comment above for
// why it is not, and cannot be, the merge decision.
//
// Both terms are overlap coefficients — shared tokens over the smaller side's
// token count — rather than Jaccard. A model that writes a four-line summary and
// a model that writes a one-line one can still be making the same claim, and
// Jaccard punishes that asymmetry for nothing but a difference in verbosity.
// The smaller-side denominator does mean a terse finding is cheap to match, and
// that shared boilerplate ("…and verify it with a focused test") counts at full
// weight; for a recall gate both are acceptable, and both are reasons this
// number must not be trusted with the decision.
func findingSimilarity(left, right findingSignature) float64 {
	return 0.5*overlapCoefficient(left.title, right.title) + 0.5*overlapCoefficient(left.body, right.body)
}

func newFindingSignature(finding Finding) findingSignature {
	return findingSignature{
		files: findingFiles(finding),
		title: dedupeTokens(finding.Title),
		// The recommendation is in the body because two models that found the
		// same problem tend to propose recognisably the same fix — the egress
		// trio all reach for krun_set_port_map and an allowlist. It also carries
		// per-repo boilerplate that inflates every score equally; that is a real
		// cost to the gate's precision, paid deliberately for its recall.
		body: dedupeTokens(finding.Title + " " + finding.Summary + " " + finding.Recommendation),
	}
}

// collapseFindingGroup folds a group of duplicates into the one finding a reader
// gets, keeping every reviewer's location and every reviewer's prose.
func collapseFindingGroup(findings []Finding, members []int) Finding {
	ordered := append([]int(nil), members...)
	// The survivor is chosen on the finding's own merits, not on arrival order.
	// Ordering the legs by name — as ReviewForSummary does — meant leg A won
	// every cross-leg merge, so for any corroborated finding the reader
	// structurally never saw leg B's version.
	sort.SliceStable(ordered, func(a, b int) bool {
		return betterSurvivor(findings[ordered[a]], findings[ordered[b]])
	})

	survivor := findings[ordered[0]]
	reviewers := map[string]struct{}{}
	providers := map[string]struct{}{}
	for _, reviewer := range findingReviewers(survivor) {
		reviewers[reviewer] = struct{}{}
	}
	providers[findingProvider(survivor)] = struct{}{}

	for _, index := range ordered[1:] {
		loser := findings[index]
		survivor = mergeFindingMetadata(survivor, loser)
		for _, reviewer := range findingReviewers(loser) {
			reviewers[reviewer] = struct{}{}
		}
		providers[findingProvider(loser)] = struct{}{}
	}
	survivor.MergedFindings = dedupeMergedFindings(survivor.MergedFindings)
	survivor = normalizeFindingCitation(survivor)

	// De-duplication runs twice — once per shard inside the panel, once across
	// shards before the judge — so the survivor may already carry a note from
	// the earlier pass. The note computed here is over the whole group and
	// supersedes it; leaving both in produced two identical "Raised
	// independently by 2 reviewers" lines in the verbose evidence listing.
	survivor.Evidence = dropEvidenceLabels(survivor.Evidence, "Corroboration", "Duplicates merged")
	survivor.Corroboration = sortedSet(reviewers)
	switch {
	case len(survivor.Corroboration) > 1:
		survivor.Evidence = append(survivor.Evidence, Evidence{
			Label: "Corroboration",
			Value: fmt.Sprintf("Raised independently by %d reviewers: %s",
				len(survivor.Corroboration), strings.Join(survivor.Corroboration, ", ")),
		})
	case len(members) > 1:
		// Same reviewer, several shards or several phrasings. Worth saying, but
		// it is not corroboration and must not read like it.
		survivor.Evidence = append(survivor.Evidence, Evidence{
			Label: "Duplicates merged",
			Value: fmt.Sprintf("%d near-identical findings from %s", len(members), strings.Join(sortedSet(providers), ", ")),
		})
	}
	return survivor
}

// betterSurvivor reports whether left should be the copy the reader gets.
//
// Strongest first, then the finding that cites the most places, then the one
// with the most specific fix (length is a crude proxy, but it is deterministic
// and it favours "call krun_set_port_map with an explicit allowlist and add a
// test asserting an outbound connection fails" over "restrict egress"). ID is
// the final tiebreak so the choice is reproducible run to run.
func betterSurvivor(left, right Finding) bool {
	if l, r := strengthRank(left.Strength), strengthRank(right.Strength); l != r {
		return l < r
	}
	if l, r := len(left.Anchors), len(right.Anchors); l != r {
		return l > r
	}
	if l, r := len(strings.TrimSpace(left.Recommendation)), len(strings.TrimSpace(right.Recommendation)); l != r {
		return l > r
	}
	if l, r := len(strings.TrimSpace(left.Summary)), len(strings.TrimSpace(right.Summary)); l != r {
		return l > r
	}
	return left.ID < right.ID
}

// normalizeFindingCitation keeps File/Line and Anchors telling the same story:
// the primary citation is the survivor's own, and it leads the anchor list that
// merging built from every reviewer's locations.
func normalizeFindingCitation(finding Finding) Finding {
	if strings.TrimSpace(finding.File) == "" {
		if len(finding.Anchors) > 0 {
			finding.File = finding.Anchors[0].File
			finding.Line = finding.Anchors[0].Line
		}
		return finding
	}
	primary := FindingAnchor{File: finding.File, Line: finding.Line}
	out := []FindingAnchor{primary}
	for _, anchor := range finding.Anchors {
		if anchor.File == primary.File && anchor.Line == primary.Line {
			continue
		}
		out = append(out, anchor)
	}
	finding.Anchors = out
	return finding
}

func findingFiles(finding Finding) map[string]struct{} {
	out := map[string]struct{}{}
	add := func(file string) {
		if file = normalizeDedupeFile(file); file != "" {
			out[file] = struct{}{}
		}
	}
	add(finding.File)
	for _, anchor := range finding.Anchors {
		add(anchor.File)
	}
	for _, text := range []string{finding.Title, finding.Summary, finding.Recommendation, evidenceText(finding.Evidence)} {
		// nil known-set and empty repo root: this is a comparison between two
		// findings, not a claim that the path exists. Whether the file is real
		// is the judge's question, not de-duplication's.
		for _, file := range parseReviewFilePaths(text, nil, "") {
			add(file)
		}
	}
	return out
}

// normalizeDedupeFile keeps path-shaped strings and rejects prose that merely
// contains a dot. Without this "e.g." parses as a file and silently satisfies
// the locality gate for any two findings that both use the abbreviation.
func normalizeDedupeFile(file string) string {
	file = strings.ToLower(strings.TrimSpace(file))
	if file == "" {
		return ""
	}
	if strings.Contains(file, "/") {
		return file
	}
	ext := strings.TrimPrefix(path.Ext(file), ".")
	base := strings.TrimSuffix(file, path.Ext(file))
	if len(ext) < 2 || len(base) < 3 {
		return ""
	}
	for _, r := range ext {
		if r < 'a' || r > 'z' {
			return ""
		}
	}
	return file
}

func sharesAny(left, right map[string]struct{}) bool {
	if len(right) < len(left) {
		left, right = right, left
	}
	for key := range left {
		if _, ok := right[key]; ok {
			return true
		}
	}
	return false
}

func overlapCoefficient(left, right map[string]struct{}) float64 {
	smaller := smallerSet(left, right)
	if smaller == 0 {
		return 0
	}
	return float64(intersectionSize(left, right)) / float64(smaller)
}

func smallerSet(left, right map[string]struct{}) int {
	if len(right) < len(left) {
		return len(right)
	}
	return len(left)
}

func intersectionSize(left, right map[string]struct{}) int {
	if len(right) < len(left) {
		left, right = right, left
	}
	count := 0
	for key := range left {
		if _, ok := right[key]; ok {
			count++
		}
	}
	return count
}

// dedupeTokens reduces review prose to the words that carry the claim.
func dedupeTokens(text string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, token := range strings.Fields(normalizeDuplicateText(text)) {
		// Dots and dashes survive normalization so that "session_test.go" stays
		// one token, which means sentence-final punctuation sticks to the last
		// word: "coverage." and "coverage" were two different tokens, and two
		// legs that ended the same sentence differently lost the match. No file
		// path begins or ends with either character, so trimming is free.
		token = strings.Trim(token, ".-")
		token = singularize(token)
		if len(token) < 3 || dedupeStopWords[token] {
			continue
		}
		out[token] = struct{}{}
	}
	return out
}

// singularize collapses the plural/singular split that otherwise costs a real
// token every time one model writes "secrets" and the other writes "secret".
func singularize(token string) string {
	if len(token) > 4 && strings.HasSuffix(token, "s") && !strings.HasSuffix(token, "ss") {
		return strings.TrimSuffix(token, "s")
	}
	return token
}

// dedupeStopWords are function words plus the connective vocabulary every
// review summary uses regardless of subject. Terms that name a subject —
// "race", "test", "secret", "lock" — are deliberately absent: they are the
// signal.
var dedupeStopWords = map[string]bool{
	"the": true, "and": true, "for": true, "with": true, "that": true, "this": true, "from": true, "into": true,
	"can": true, "are": true, "was": true, "were": true, "has": true, "have": true, "but": true, "not": true,
	"which": true, "when": true, "where": true, "what": true, "who": true, "how": true, "why": true,
	"would": true, "could": true, "should": true, "will": true, "may": true, "might": true, "must": true,
	"there": true, "their": true, "them": true, "then": true, "than": true, "they": true, "these": true, "those": true,
	"also": true, "only": true, "just": true, "even": true, "still": true, "yet": true, "both": true, "each": true,
	"any": true, "all": true, "its": true, "does": true, "done": true, "being": true, "been": true,
	"use": true, "used": true, "using": true, "make": true, "made": true, "get": true, "got": true,
	"inside": true, "within": true, "without": true, "before": true, "after": true, "during": true,
	"because": true, "while": true, "however": true, "instead": true, "rather": true, "such": true,
	"here": true, "over": true, "under": true, "between": true, "against": true, "about": true,
	"more": true, "most": true, "less": true, "least": true, "very": true, "same": true, "other": true,
	"code": true, "line": true, "call": true, "calls": true, "case": true, "way": true, "thing": true,
}

func mergeResolvedSources(left, right []ResolvedSource) []ResolvedSource {
	seen := map[string]struct{}{}
	var out []ResolvedSource
	for _, src := range append(append([]ResolvedSource{}, left...), right...) {
		key := strings.TrimSpace(src.ID) + "|" + strings.TrimSpace(ResolvedSourceLabel(src))
		if key == "|" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, src)
	}
	return out
}

// mergeFindingMetadata folds the loser into the survivor.
//
// The loser's prose is kept, not dropped. Concatenating two descriptions of one
// problem reads worse than either alone, so the survivor's Title/Summary/
// Recommendation stay the ones the reader meets — but the other reviewer's
// framing and, more importantly, its proposed fix are recorded on
// MergedFindings, rendered under the survivor, and serialised in the JSON
// report. Merging must never be the reason a concrete recommendation vanishes.
func mergeFindingMetadata(left Finding, right Finding) Finding {
	left.MergedFindings = append(left.MergedFindings, mergedFindingRecord(right))
	left.MergedFindings = append(left.MergedFindings, right.MergedFindings...)
	left.Evidence = append(left.Evidence, right.Evidence...)
	left.Anchors = mergeFindingAnchors(left.Anchors, right.Anchors)
	left.SourceIDs = uniqueStrings(append(left.SourceIDs, right.SourceIDs...))
	left.SourcePublishers = uniqueStrings(append(left.SourcePublishers, right.SourcePublishers...))
	left.ResolvedSources = mergeResolvedSources(left.ResolvedSources, right.ResolvedSources)
	left.Scopes = uniqueStrings(append(left.Scopes, right.Scopes...))
	if strengthRank(right.Strength) < strengthRank(left.Strength) {
		left.Strength = right.Strength
	}
	if strings.TrimSpace(left.Benefit) == "" {
		left.Benefit = right.Benefit
	}
	return left
}

func dropEvidenceLabels(evidence []Evidence, labels ...string) []Evidence {
	drop := map[string]struct{}{}
	for _, label := range labels {
		drop[label] = struct{}{}
	}
	// A fresh slice, not evidence[:0]: findings share backing arrays after
	// mergeFindingMetadata's appends, and filtering in place would rewrite
	// another finding's evidence.
	out := make([]Evidence, 0, len(evidence))
	for _, item := range evidence {
		if _, ok := drop[item.Label]; ok {
			continue
		}
		out = append(out, item)
	}
	return out
}

func mergedFindingRecord(finding Finding) MergedFinding {
	reviewer := ""
	if len(finding.Corroboration) > 0 {
		reviewer = finding.Corroboration[0]
	}
	return MergedFinding{
		ID:             finding.ID,
		Reviewer:       reviewer,
		Title:          strings.TrimSpace(finding.Title),
		Summary:        strings.TrimSpace(finding.Summary),
		Recommendation: strings.TrimSpace(finding.Recommendation),
	}
}

func dedupeMergedFindings(merged []MergedFinding) []MergedFinding {
	seen := map[string]struct{}{}
	var out []MergedFinding
	for _, item := range merged {
		key := item.ID + "|" + item.Title
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, item)
	}
	return out
}

// mergeFindingAnchors keeps both reviewers' locations. Two models pointing at
// two lines of the same problem give the reader more to go on than one does.
func mergeFindingAnchors(left, right []FindingAnchor) []FindingAnchor {
	seen := map[string]struct{}{}
	var out []FindingAnchor
	for _, anchor := range append(append([]FindingAnchor{}, left...), right...) {
		if strings.TrimSpace(anchor.File) == "" {
			continue
		}
		key := fmt.Sprintf("%s:%d", anchor.File, anchor.Line)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, anchor)
	}
	return out
}

// findingReviewers is the set of reviewer legs a finding is already credited to.
// Findings that never went through a reviewer leg (deterministic rules) report
// nothing, so corroboration counting never silently treats two rule findings as
// two agreeing models.
func findingReviewers(finding Finding) []string {
	if len(finding.Corroboration) > 0 {
		return finding.Corroboration
	}
	return nil
}

func normalizeDuplicateText(text string) string {
	text = strings.ToLower(text)
	var b strings.Builder
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '/' || r == '_' || r == '-' || r == '.' {
			b.WriteRune(r)
			continue
		}
		b.WriteByte(' ')
	}
	return b.String()
}

func findingProvider(finding Finding) string {
	id := strings.TrimSpace(finding.ID)
	if id == "" {
		return "unknown"
	}
	if index := strings.Index(id, "."); index > 0 {
		return id[:index]
	}
	return id
}

func sortedSet(values map[string]struct{}) []string {
	out := make([]string, 0, len(values))
	for value := range values {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
