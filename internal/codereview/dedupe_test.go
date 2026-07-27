package codereview

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
)

// The findings below are verbatim output from live two-leg Bedrock reviews of
// github.com/satoricorp/yeet. They are not written for the test: invented prose
// would let the test pass against a rule that only recognises phrasing a test
// author thought of, which is precisely how the previous 85%-token-overlap rule
// survived, and how the threshold that replaced it survived too.

// egressA/egressB/egressC are one issue in three voices — the case that started
// this work. Two came from one reviewer leg, one from the other.
var egressA = Finding{
	ID:      "bedrock-a.ai.1",
	Title:   "Unrestricted guest egress makes the documented secrets model weaker than it appears",
	Summary: "`docs/secrets.md` documents that `run.env` is injected at 0600 inside a 0700 dir, acknowledges the guest holds the secret, and defers a proxy to 'a separate project'. Meanwhile `docs/networking.md` confirms TSI gives the guest unrestricted outbound. Together this means any agent that exfiltrates secrets cannot be detected or prevented today.",
	Recommendation: "Add a config option `network off|on` (default `on`) that, when off, calls `krun_set_port_map` with an empty map and skips TSI in `guest/launcher/yeet-vm.c` (the API is already mentioned in docs/networking.md line 87). " +
		"Surface a one-line warning in the output when agent secrets are injected AND networking is unrestricted.",
	Corroboration: []string{"Bedrock A (claude-opus-4-6 via AWS)"},
	Anchors:       []FindingAnchor{{File: "guest/launcher/yeet-vm.c", Line: 201}},
}

var egressB = Finding{
	ID:      "bedrock-b.ai.3",
	Title:   "Unrestricted guest egress enables secret exfiltration with no mitigation timeline",
	Summary: "The C launcher (`yeet-vm.c`) explicitly documents that TSI networking gives the guest \"unrestricted egress\" and the docs/secrets.md doc acknowledges \"the guest still holds the secret\" and calls a host-side proxy \"a separate project.\" Meanwhile, `run.env` is sourced inside the VM with the provider API key AND any agent secrets. A malicious or confused agent can trivially `curl` secrets to an external endpoint.",
	Recommendation: "In `guest/launcher/yeet-vm.c`, after the TSI comment at line ~245, add a TODO with a concrete plan: use `krun_set_port_map` to restrict outbound to known model API hosts as a near-term step. " +
		"Longer term, implement the filtering proxy described in docs/networking.md.",
	Corroboration: []string{"Bedrock B (claude-opus-4-5 via AWS)"},
	Anchors:       []FindingAnchor{{File: "guest/launcher/yeet-vm.c", Line: 1}},
}

var egressC = Finding{
	ID:             "bedrock-a.ai.5",
	Title:          "Unrestricted guest egress makes run.env exfiltratable despite 0600 permissions",
	Summary:        "docs/secrets.md acknowledges this explicitly ('the guest still holds the secret… anything bridged in is exfiltratable by an agent that wants to'), and the networking doc confirms TSI gives the guest unrestricted outbound. The current mitigation is file permissions (0600/0700), which protect against other users on the host but do nothing inside the VM where the agent runs as root.",
	Recommendation: "Add a `krun_set_port_map` call or TSI filtering proxy as described in docs/networking.md's security note. As an immediate partial mitigation, rotate `run.env` to a per-iteration nonce that the host validates via loopback.",
	Corroboration:  []string{"Bedrock A (claude-opus-4-6 via AWS)"},
}

// liveEgress1/liveEgress2 are the pair a live run shipped as findings #1 and #2
// AFTER a similarity threshold was in place — the proof that a lexical cut
// cannot do this job. They score below that threshold, and a human reads them as
// one finding at a glance.
var liveEgress1 = Finding{
	ID:             "bedrock-a.ai.1",
	Title:          "Unrestricted network egress from guest VM is a live exfiltration channel for secrets",
	Summary:        "The secrets design doc (docs/secrets.md) explicitly acknowledges that `run.env` is inside the VM and the VM has unrestricted outbound networking via TSI, making any bridged secret exfiltratable. The launcher (`yeet-vm.c`) confirms no egress filtering is configured — `krun_add_vsock_port` is unused and there is no call to `krun_disable_implicit_vsock` or any filtering proxy. Meanwhile `loop.ts` writes the provider API key into `run.env` every iteration. A malicious or confused agent can `curl` the key to any endpoint.",
	Recommendation: "In `guest/launcher/yeet-vm.c`, add an egress policy seam: accept a `--egress-allow` flag listing allowed destination hosts/CIDRs, and when present call `krun_set_net_filter` (or the TSI equivalent) to restrict outbound. As a first step, document the known-required destinations (OpenRouter, GitHub for git fetch) and wire the allowlist from `src/vm.ts` where `runVM` assembles the launcher args.",
	Corroboration:  []string{"Bedrock A (claude-opus-4-6 via AWS)"},
	Anchors:        []FindingAnchor{{File: "guest/launcher/yeet-vm.c", Line: 201}, {File: "docs/secrets.md", Line: 70}},
}

var liveEgress2 = Finding{
	ID:             "bedrock-b.ai.2",
	Title:          "Unrestricted VM egress is the primary security gap — add an egress-policy seam now",
	Summary:        "The C launcher (`yeet-vm.c`) explicitly documents unrestricted outbound networking via libkrun's TSI HIJACK_INET. The guest holds the provider API key in `run.env` and has full internet access. A misbehaving or compromised agent can exfiltrate every secret bridged into the VM. The codebase acknowledges this (docs/secrets.md \"the guest still holds the secret\") but has no filtering, no proxy, no per-run nonce, and no monitoring of egress.",
	Recommendation: "In `src/vm.ts` (which assembles the `--env` list and calls the launcher), add a `--env HTTPS_PROXY=…` pointing at a lightweight host-side forward proxy. In the proxy, allowlist only the model-provider domains derived from `KNOWN_PROVIDERS` in `src/secrets.ts`. Gate the feature behind `yeet config egress strict`.",
	Corroboration:  []string{"Bedrock B (claude-opus-4-5-20251101 via AWS)"},
	Anchors:        []FindingAnchor{{File: "guest/launcher/yeet-vm.c", Line: 201}, {File: "src/secrets.ts", Line: 1}},
}

// bridgeStop/bridgePoll/bridgeSync are three DIFFERENT problems in one file, two
// of them on the same line. They are the reason locality cannot decide on its
// own.
var bridgeStop = Finding{
	ID:             "bedrock-a.ai.2",
	Title:          "Race between Bridge tick and VM exit loses final events",
	Summary:        "Bridge.stop() calls drainAgentEvents() once, but tick() is async and guarded by `this.busy`. If the timer fires one last time just before stop() is called, the `busy` flag can be true when stop() runs, causing the final drain in stop() to read from a stale offset. If the VM exits between two ticks and the guest has already deleted `agent.raw.jsonl`, the tail in stop() reads a file that no longer exists and silently returns nothing.",
	Recommendation: "In `src/bridge.ts` Bridge.stop(), await the current tick if busy before draining. Clear the interval first, then `while (this.busy) await new Promise(r => setTimeout(r, 10));` before the final drainAgentEvents().",
	Corroboration:  []string{"Bedrock A (claude-opus-4-6 via AWS)"},
	Anchors:        []FindingAnchor{{File: "src/bridge.ts", Line: 160}},
}

var bridgePoll = Finding{
	ID:             "bedrock-b.ai.4",
	Title:          "Bridge polling loop has no backoff and can spin on malformed request files",
	Summary:        "In `src/bridge.ts`, the Bridge polls the `qa/` directory for `.ask.json` and `.verify.json` files using a tight `setInterval`. If a file is present but malformed, the `JSON.parse` in the polling handler throws, the file remains on disk, and the next poll interval immediately re-attempts, throwing again indefinitely.",
	Recommendation: "In `src/bridge.ts` around line 160, after catching a JSON parse error, rename the offending file to `<name>.bad` and emit a warning event rather than silently retrying every poll tick. Add an exponential backoff counter per-file that resets on success.",
	Corroboration:  []string{"Bedrock B (claude-opus-4-5 via AWS)"},
	Anchors:        []FindingAnchor{{File: "src/bridge.ts", Line: 160}},
}

var bridgeSync = Finding{
	ID:             "bedrock-b.ai.6",
	Title:          "Blocking readFileSync calls inside async functions",
	Summary:        "Multiple async code paths in src/agent.ts, src/cli.ts, src/loop.ts, and src/bridge.ts use synchronous readFileSync. Under Bun's event loop this stalls all in-flight promises, including the bridge's cost-cap ticker.",
	Recommendation: "Change readFileSync calls in hot paths (e.g. bridge.ts:160, loop.ts:262/265/315, cli.ts:132/417) to await readFile from fs/promises. Start with bridge.ts where the loop ticks every second.",
	Corroboration:  []string{"Bedrock B (claude-opus-4-5 via AWS)"},
	Anchors:        []FindingAnchor{{File: "src/bridge.ts", Line: 160}},
}

// secretsRace/secretsMount are two genuinely different defects in one file that
// a similarity threshold merges — different root cause, different threat model,
// different fix. Merging them deleted the second one's recommendation and then
// printed "flagged by both reviewers independently" over the first.
var secretsRace = Finding{
	ID:             "bedrock-a.ai.7",
	Title:          "src/secrets.ts writes the key file with a permissions race",
	Summary:        "`saveKeys` in `src/secrets.ts` writes the key file and only then sets its mode. Between those two calls the file exists with the process umask, so on a shared host another local user can read the provider keys in that window.",
	Recommendation: "In `src/secrets.ts`, create the key file with an exclusive open at mode 0600 so the mode is applied atomically at creation, and write through the returned descriptor. Add a test asserting the mode is never world-readable at any point.",
	Corroboration:  []string{"Bedrock A (claude-opus-4-6 via AWS)"},
	Anchors:        []FindingAnchor{{File: "src/secrets.ts", Line: 66}},
}

var secretsMount = Finding{
	ID:             "bedrock-b.ai.8",
	Title:          "src/secrets.ts keeps the key file inside a directory bind-mounted into the guest",
	Summary:        "The key file path in `src/secrets.ts` resolves under a directory that the launcher passes to the guest as a virtio-fs mount. The agent inside the VM therefore reads the provider keys directly, whatever their mode is on the host, because inside the guest it runs as root.",
	Recommendation: "In `src/secrets.ts`, move the key file out of the mounted tree entirely and bridge only the single provider key a run needs. Add a test that boots a VM and asserts the key file is not visible from the guest.",
	Corroboration:  []string{"Bedrock B (claude-opus-4-5 via AWS)"},
	Anchors:        []FindingAnchor{{File: "src/secrets.ts", Line: 12}},
}

// writerClose/writerSync are the clearest statement of why a lexical score
// cannot decide: they differ by one word in the title, they are two different
// bugs, and they score higher than any true duplicate in this file.
var writerClose = Finding{
	ID:             "bedrock-a.ai.9",
	Title:          "Errors from Close are ignored in the writer",
	Summary:        "`internal/store/writer.go` defers Close and discards its error. For a buffered file that is where a late ENOSPC surfaces, so a write that failed at the very end is reported to the caller as a success.",
	Recommendation: "In `internal/store/writer.go`, capture the deferred Close error into the named return of the write function, and add a test using a stub whose Close fails.",
	Corroboration:  []string{"Bedrock A (claude-opus-4-6 via AWS)"},
	Anchors:        []FindingAnchor{{File: "internal/store/writer.go", Line: 88}},
}

var writerSync = Finding{
	ID:             "bedrock-b.ai.10",
	Title:          "Errors from Sync are ignored in the writer",
	Summary:        "`internal/store/writer.go` defers Sync and discards its error. The data is only in the page cache when the write function reports success, so a power loss loses records the caller was told were durable.",
	Recommendation: "In `internal/store/writer.go`, capture the deferred Sync error into the named return of the write function, and add a test using a stub whose Sync fails.",
	Corroboration:  []string{"Bedrock B (claude-opus-4-5 via AWS)"},
	Anchors:        []FindingAnchor{{File: "internal/store/writer.go", Line: 92}},
}

// pairLabel is the hand-labelled truth for the fixtures above. It drives the
// fake adjudicator and it is also the sample the "no threshold works" test
// measures.
type pairLabel struct {
	left  Finding
	right Finding
	same  bool
	note  string
}

func labelledFindingPairs() []pairLabel {
	return []pairLabel{
		{egressA, egressB, true, "one egress issue, two voices"},
		{egressB, egressC, true, "one egress issue, two voices"},
		{egressA, egressC, true, "one egress issue, two voices"},
		{liveEgress1, liveEgress2, true, "the pair a live run shipped side by side"},
		{bridgeStop, bridgePoll, false, "re-entrancy race vs unbounded retry"},
		{bridgeStop, bridgeSync, false, "re-entrancy race vs blocking IO"},
		{bridgePoll, bridgeSync, false, "unbounded retry vs blocking IO"},
		{secretsRace, secretsMount, false, "creation race vs bind-mount exposure"},
		{writerClose, writerSync, false, "ignored Close error vs missing fsync"},
	}
}

// labelledAdjudicator answers from the table above and fails the test if it is
// asked about a pair nobody labelled.
//
// The fake stands in for the model's judgement, not for the pipeline. What these
// tests pin is that the gate asks the right questions, that a "same" answer
// merges and records corroboration, and that a "different" answer leaves two
// findings and two fixes standing. Whether the real model answers well is what
// the live review at the end of this work checks; no fake can establish that,
// and none of these tests claims to.
type labelledAdjudicator struct {
	t      *testing.T
	answer map[string]bool
	asked  []string
	err    error
}

func newLabelledAdjudicator(t *testing.T) *labelledAdjudicator {
	t.Helper()
	answer := map[string]bool{}
	for _, pair := range labelledFindingPairs() {
		answer[titlePairKey(pair.left.Title, pair.right.Title)] = pair.same
	}
	return &labelledAdjudicator{t: t, answer: answer}
}

func (a *labelledAdjudicator) AdjudicateDuplicates(_ context.Context, pairs []duplicatePairInput) ([]duplicateVerdict, error) {
	if a.err != nil {
		return nil, a.err
	}
	var out []duplicateVerdict
	for _, pair := range pairs {
		key := titlePairKey(pair.Left.Title, pair.Right.Title)
		a.asked = append(a.asked, key)
		same, ok := a.answer[key]
		if !ok {
			a.t.Fatalf("adjudicator asked about an unlabelled pair:\n  %q\n  %q", pair.Left.Title, pair.Right.Title)
		}
		out = append(out, duplicateVerdict{PairID: pair.PairID, Same: same})
	}
	return out, nil
}

func (a *labelledAdjudicator) askedAbout(left, right Finding) bool {
	key := titlePairKey(left.Title, right.Title)
	for _, asked := range a.asked {
		if asked == key {
			return true
		}
	}
	return false
}

func titlePairKey(left, right string) string {
	titles := []string{left, right}
	sort.Strings(titles)
	return titles[0] + "\x00" + titles[1]
}

func mergeWithLabels(t *testing.T, findings ...Finding) []Finding {
	t.Helper()
	return mergeNearDuplicateFindings(context.Background(), newLabelledAdjudicator(t), findings)
}

func pairScore(left, right Finding) float64 {
	return findingSimilarity(newFindingSignature(left), newFindingSignature(right))
}

// TestNoLexicalThresholdSeparatesDuplicatesFromDistinctFindings is the reason
// this pipeline has a model in it, stated as a test.
//
// If some cut on findingSimilarity separated the labelled duplicates from the
// labelled distinct pairs, the adjudicator would be unjustified complexity and
// this test would say so. It does not: the classes interleave, so every cut
// either merges two real findings into one — deleting one of them — or leaves
// duplicates a reader will call duplicates.
func TestNoLexicalThresholdSeparatesDuplicatesFromDistinctFindings(t *testing.T) {
	lowestDuplicate, highestDistinct := 1.0, 0.0
	var lowestDuplicateNote, highestDistinctNote string
	for _, pair := range labelledFindingPairs() {
		score := pairScore(pair.left, pair.right)
		if pair.same && score < lowestDuplicate {
			lowestDuplicate, lowestDuplicateNote = score, pair.note
		}
		if !pair.same && score > highestDistinct {
			highestDistinct, highestDistinctNote = score, pair.note
		}
	}
	if highestDistinct <= lowestDuplicate {
		t.Fatalf("lexical scores separate the labelled set (distinct max %.4f %q < duplicate min %.4f %q); "+
			"a threshold would do, and the adjudicator is unjustified complexity",
			highestDistinct, highestDistinctNote, lowestDuplicate, lowestDuplicateNote)
	}
	t.Logf("classes interleave: highest distinct pair %.4f (%s) > lowest duplicate pair %.4f (%s)",
		highestDistinct, highestDistinctNote, lowestDuplicate, lowestDuplicateNote)
}

// TestGateAdmitsEveryLabelledDuplicate pins the one thing the deterministic
// stage must get right. It is a recall filter: a duplicate it drops can never be
// merged, whatever the adjudicator would have said.
func TestGateAdmitsEveryLabelledDuplicate(t *testing.T) {
	for _, pair := range labelledFindingPairs() {
		if !pair.same {
			continue
		}
		left, right := newFindingSignature(pair.left), newFindingSignature(pair.right)
		if !sharesAny(left.files, right.files) {
			t.Errorf("file gate drops a real duplicate (%s):\n  %q\n  %q", pair.note, pair.left.Title, pair.right.Title)
			continue
		}
		if score := findingSimilarity(left, right); score < dedupeGateThreshold {
			t.Errorf("gate drops a real duplicate at %.4f < %.2f (%s):\n  %q\n  %q",
				score, dedupeGateThreshold, pair.note, pair.left.Title, pair.right.Title)
		}
	}
}

// TestMergesSemanticDuplicatesAcrossReviewers is the original regression test:
// three findings, one issue, no shared title beyond "Unrestricted guest egress",
// raised by two different reviewer legs.
func TestMergesSemanticDuplicatesAcrossReviewers(t *testing.T) {
	merged := mergeWithLabels(t, egressA, egressB, egressC)
	if len(merged) != 1 {
		t.Fatalf("merged %d findings, want 1:\n%s", len(merged), findingTitles(merged))
	}
	if got := len(merged[0].Corroboration); got != 2 {
		t.Fatalf("Corroboration = %v, want both reviewer legs", merged[0].Corroboration)
	}
	if !strings.Contains(merged[0].Corroboration[0], "Bedrock A") || !strings.Contains(merged[0].Corroboration[1], "Bedrock B") {
		t.Errorf("Corroboration = %v, want one entry per leg", merged[0].Corroboration)
	}
	// Both reviewers' locations survive: the reader gets every place either
	// model pointed at, not just the winner's.
	if len(merged[0].Anchors) != 2 {
		t.Errorf("Anchors = %v, want both legs' anchors", merged[0].Anchors)
	}
}

// TestMergesTheDuplicatePairAThresholdShipped is the regression test for the
// defect that survived the first fix: a live review returned these two findings
// as #1 and #2 with a similarity threshold of 0.36 in place.
func TestMergesTheDuplicatePairAThresholdShipped(t *testing.T) {
	// Recorded rather than assumed: the pair scores below the threshold that
	// shipped, which is why no re-tuning of a threshold was the answer.
	if score := pairScore(liveEgress1, liveEgress2); score >= 0.36 {
		t.Fatalf("fixture drift: this pair scores %.4f, but it shipped as two findings under a 0.36 cut", score)
	}
	merged := mergeWithLabels(t, liveEgress1, liveEgress2)
	if len(merged) != 1 {
		t.Fatalf("shipped the same duplicate pair again — merged %d findings, want 1:\n%s", len(merged), findingTitles(merged))
	}
	if len(merged[0].Corroboration) != 2 {
		t.Errorf("Corroboration = %v, want both legs", merged[0].Corroboration)
	}
}

// TestKeepsDistinctFindingsOnTheSameLine is the other half. Merging these would
// delete two real problems, which is strictly worse than printing a near
// duplicate.
func TestKeepsDistinctFindingsOnTheSameLine(t *testing.T) {
	merged := mergeWithLabels(t, bridgeStop, bridgePoll, bridgeSync)
	if len(merged) != 3 {
		t.Fatalf("merged %d findings, want 3 (all distinct):\n%s", len(merged), findingTitles(merged))
	}
	for _, finding := range merged {
		if len(finding.Corroboration) > 1 {
			t.Errorf("%q claims corroboration %v but no reviewer agreed with it", finding.Title, finding.Corroboration)
		}
	}
}

// TestDistinctFindingsAThresholdWouldMergeSurvive covers the pairs a lexical cut
// gets wrong in the expensive direction: both score above any cut that would
// catch the live duplicate above, and both are two bugs, not one.
func TestDistinctFindingsAThresholdWouldMergeSurvive(t *testing.T) {
	live := pairScore(liveEgress1, liveEgress2)
	for _, pair := range []pairLabel{
		{secretsRace, secretsMount, false, "creation race vs bind-mount exposure"},
		{writerClose, writerSync, false, "ignored Close error vs missing fsync"},
	} {
		score := pairScore(pair.left, pair.right)
		if score <= live {
			t.Fatalf("fixture drift: %s scores %.4f, below the live duplicate at %.4f, so it no longer demonstrates the overlap",
				pair.note, score, live)
		}
		adjudicator := newLabelledAdjudicator(t)
		merged := mergeNearDuplicateFindings(context.Background(), adjudicator, []Finding{pair.left, pair.right})
		if len(merged) != 2 {
			t.Fatalf("%s (score %.4f) merged into %d finding(s); a merge here deletes a real defect:\n%s",
				pair.note, score, len(merged), findingTitles(merged))
		}
		if !adjudicator.askedAbout(pair.left, pair.right) {
			t.Errorf("%s survived only because the gate never offered it; the gate must ask and the adjudicator must answer", pair.note)
		}
		for _, finding := range merged {
			if len(finding.Corroboration) > 1 {
				t.Errorf("%q claims corroboration %v after no merge happened", finding.Title, finding.Corroboration)
			}
		}
	}
}

// TestOverMergeCannotFabricateCorroboration states the consequence the previous
// design got wrong twice over: merging two distinct findings from two legs did
// not merely drop one, it printed "flagged by both reviewers independently" over
// the survivor and promoted it above the judge's own confidence.
func TestOverMergeCannotFabricateCorroboration(t *testing.T) {
	merged := mergeWithLabels(t, secretsRace, secretsMount)
	if len(merged) != 2 {
		t.Fatalf("merged %d findings, want 2", len(merged))
	}
	rendered := RenderMarkdown(Report{ReviewMode: ReviewModeRepo, Reviewed: true, Findings: merged})
	if strings.Contains(rendered, "Flagged by both reviewers") {
		t.Errorf("two distinct findings produced a corroboration claim:\n%s", rendered)
	}
	if !strings.Contains(rendered, "move the key file out of the mounted tree") {
		t.Errorf("the second finding's fix is missing from the review:\n%s", rendered)
	}
}

// TestFileGateBlocksUnrelatedFindingsThatShareVocabulary pins the precondition:
// review prose repeats itself, and shared words are not shared subjects.
func TestFileGateBlocksUnrelatedFindingsThatShareVocabulary(t *testing.T) {
	left := Finding{
		ID:            "bedrock-a.ai.1",
		Title:         "No automated test coverage for the secrets module",
		Summary:       "src/secrets.ts has no automated test coverage at all.",
		Corroboration: []string{"Bedrock A"},
	}
	right := Finding{
		ID:            "bedrock-b.ai.1",
		Title:         "No automated test coverage for the search module",
		Summary:       "src/search.ts has no automated test coverage at all.",
		Corroboration: []string{"Bedrock B"},
	}
	adjudicator := newLabelledAdjudicator(t) // fails the test if it is consulted
	if merged := mergeNearDuplicateFindings(context.Background(), adjudicator, []Finding{left, right}); len(merged) != 2 {
		t.Fatalf("merged findings about different files:\n%s", findingTitles(merged))
	}
}

// TestSameLegDuplicatesAreNotCorroboration guards the signal's meaning. One
// model saying a thing in two shards is not two models agreeing, and a review
// that reported it as agreement would be inflating its own confidence.
func TestSameLegDuplicatesAreNotCorroboration(t *testing.T) {
	second := egressC
	second.Corroboration = egressA.Corroboration
	merged := mergeWithLabels(t, egressA, second)
	if len(merged) != 1 {
		t.Fatalf("merged %d findings, want 1", len(merged))
	}
	if len(merged[0].Corroboration) != 1 {
		t.Fatalf("Corroboration = %v, want the single leg that raised it", merged[0].Corroboration)
	}
	for _, evidence := range merged[0].Evidence {
		if evidence.Label == "Corroboration" {
			t.Fatalf("same-leg duplicate reported as corroboration: %q", evidence.Value)
		}
	}
}

// TestCorroborationSurvivesToOutput checks the two places a reader actually
// meets a finding. Recording corroboration and then not showing it would leave
// the defect fixed in theory only.
func TestCorroborationSurvivesToOutput(t *testing.T) {
	merged := mergeWithLabels(t, egressA, egressB)
	report := Report{
		ReviewMode: ReviewModeRepo,
		Reviewed:   true,
		Findings:   merged,
	}

	rendered := RenderMarkdown(report)
	if !strings.Contains(rendered, "Flagged by both reviewers") {
		t.Errorf("rendered review does not surface corroboration:\n%s", rendered)
	}

	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal report: %v", err)
	}
	var decoded struct {
		Findings []struct {
			Corroboration []string `json:"corroboration"`
		} `json:"findings"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal report: %v", err)
	}
	if len(decoded.Findings) != 1 || len(decoded.Findings[0].Corroboration) != 2 {
		t.Errorf("JSON report lost corroboration: %s", raw)
	}
}

// TestSingleReviewerFindingReadsAsSingleReviewer is the negative of the above:
// the two must be distinguishable in the output, or the signal says nothing.
func TestSingleReviewerFindingReadsAsSingleReviewer(t *testing.T) {
	merged := mergeNearDuplicateFindings(context.Background(), nil, []Finding{egressA})
	rendered := RenderMarkdown(Report{ReviewMode: ReviewModeRepo, Reviewed: true, Findings: merged})
	if strings.Contains(rendered, "Flagged by both reviewers") {
		t.Errorf("single-reviewer finding claims corroboration:\n%s", rendered)
	}
}

// TestMergeKeepsTheOtherReviewersFix pins what a merge may not cost. The copy
// the reader meets is one of two, and the one folded away frequently carries the
// more concrete remedy; dropping it silently made the review less useful on the
// findings it was most confident about.
func TestMergeKeepsTheOtherReviewersFix(t *testing.T) {
	merged := mergeWithLabels(t, liveEgress1, liveEgress2)
	if len(merged) != 1 {
		t.Fatalf("merged %d findings, want 1", len(merged))
	}
	if len(merged[0].MergedFindings) != 1 {
		t.Fatalf("MergedFindings = %#v, want the folded-away version", merged[0].MergedFindings)
	}

	dropped := liveEgress2.Recommendation
	if merged[0].Recommendation == dropped {
		dropped = liveEgress1.Recommendation
	}
	rendered := RenderMarkdown(Report{ReviewMode: ReviewModeRepo, Reviewed: true, Findings: merged})
	if !strings.Contains(rendered, "Also reported as:") {
		t.Fatalf("merged review does not show the other reviewer's version:\n%s", rendered)
	}
	if !strings.Contains(rendered, firstSentence(dropped)) {
		t.Errorf("merging deleted the other reviewer's fix %q from the review:\n%s", firstSentence(dropped), rendered)
	}

	raw, err := json.Marshal(Report{Findings: merged})
	if err != nil {
		t.Fatalf("marshal report: %v", err)
	}
	if !strings.Contains(string(raw), "merged_findings") {
		t.Errorf("JSON report lost the folded-away finding: %s", raw)
	}
}

// TestSurvivorIsChosenOnMeritNotArrivalOrder pins the other half of the same
// defect. ReviewForSummary sorts legs by name, so an arrival-order survivor made
// leg A win every cross-leg merge by construction, and leg B's fix was
// structurally unreachable for every corroborated finding.
func TestSurvivorIsChosenOnMeritNotArrivalOrder(t *testing.T) {
	forward := mergeWithLabels(t, egressA, egressB)
	backward := mergeWithLabels(t, egressB, egressA)
	if len(forward) != 1 || len(backward) != 1 {
		t.Fatalf("merged %d and %d findings, want 1 each", len(forward), len(backward))
	}
	if forward[0].ID != backward[0].ID {
		t.Errorf("survivor depends on arrival order: %q forward, %q reversed", forward[0].ID, backward[0].ID)
	}
	if forward[0].Recommendation != backward[0].Recommendation {
		t.Errorf("surviving recommendation depends on arrival order")
	}
	// Whichever wins, the other leg's fix still reaches the reader.
	rendered := RenderMarkdown(Report{ReviewMode: ReviewModeRepo, Reviewed: true, Findings: backward})
	for _, recommendation := range []string{egressA.Recommendation, egressB.Recommendation} {
		if !strings.Contains(rendered, firstSentence(recommendation)) {
			t.Errorf("a leg's fix did not survive the merge: %q\n%s", firstSentence(recommendation), rendered)
		}
	}
}

// TestWithoutAnAdjudicatorNothingSemanticIsMerged pins the fallback. No
// credentials must mean "did not merge", never "merged by guessing".
func TestWithoutAnAdjudicatorNothingSemanticIsMerged(t *testing.T) {
	outcome := dedupeFindings(context.Background(), nil, []Finding{egressA, egressB, egressC})
	if len(outcome.Findings) != 3 {
		t.Fatalf("merged %d findings with no adjudicator; the fallback must not guess:\n%s",
			len(outcome.Findings), findingTitles(outcome.Findings))
	}
	if outcome.Unresolved() == 0 {
		t.Fatalf("gated %d pairs, adjudicated %d, reported nothing unresolved", outcome.Gated, outcome.Adjudicated)
	}

	// A near-verbatim copy is still collapsed: that is not a judgement call.
	copied := egressA
	copied.ID = "bedrock-a.ai.9"
	verbatim := dedupeFindings(context.Background(), nil, []Finding{egressA, copied})
	if len(verbatim.Findings) != 1 {
		t.Fatalf("near-verbatim copies survived the fallback:\n%s", findingTitles(verbatim.Findings))
	}
}

// TestAdjudicationFailureIsReportedNotSwallowed keeps a failed call from looking
// like a clean de-duplication.
func TestAdjudicationFailureIsReportedNotSwallowed(t *testing.T) {
	adjudicator := newLabelledAdjudicator(t)
	adjudicator.err = fmt.Errorf("bedrock: throttled")
	outcome := dedupeFindings(context.Background(), adjudicator, []Finding{egressA, egressB, egressC})
	if len(outcome.Findings) != 3 {
		t.Fatalf("merged findings on a failed adjudication:\n%s", findingTitles(outcome.Findings))
	}
	if outcome.Err == nil || outcome.Unresolved() == 0 {
		t.Fatalf("failed adjudication reported Err=%v Unresolved=%d", outcome.Err, outcome.Unresolved())
	}
}

// TestCorroboratedFindingOutranksSingleReviewerFinding pins the promotion. Two
// models agreeing is evidence, and evidence should move a finding up the page.
func TestCorroboratedFindingOutranksSingleReviewerFinding(t *testing.T) {
	solo := Finding{ID: "bedrock-a.ai.1", Title: "Solo", Corroboration: []string{"Bedrock A"}}
	both := Finding{ID: "bedrock-b.ai.2", Title: "Both", Corroboration: []string{"Bedrock A", "Bedrock B"}}
	results := []judgeResult{
		{CandidateID: "bedrock-a.ai.1", Verdict: "confirmed", Impact: impactFunctional, Confidence: 0.9},
		{CandidateID: "bedrock-b.ai.2", Verdict: "confirmed", Impact: impactFunctional, Confidence: 0.9},
	}
	kept := applyJudgeResults([]Finding{solo, both}, results)
	if len(kept) != 2 {
		t.Fatalf("kept %d findings, want 2", len(kept))
	}
	if kept[0].Title != "Both" {
		t.Errorf("ranking put %q first; corroborated finding should lead", kept[0].Title)
	}
}

// TestPartialReviewerFailureDoesNotClaimNoAIReview pins the wording. A live run
// with one leg failing printed "AI review unavailable …; results are from
// deterministic checks only" above four findings the surviving model had
// written, one of them corroborated by the leg the same banner called absent.
func TestPartialReviewerFailureDoesNotClaimNoAIReview(t *testing.T) {
	report := Report{
		ReviewMode:      ReviewModeRepo,
		Reviewed:        true,
		Reviewer:        "heuristic+ai",
		DegradedReasons: []string{"one reviewer did not run — Bedrock B: decode AI review JSON"},
		Findings:        mergeNearDuplicateFindings(context.Background(), nil, []Finding{egressA}),
	}
	rendered := RenderMarkdown(report)
	if strings.Contains(rendered, "deterministic checks only") {
		t.Errorf("partial reviewer failure reported as no AI review at all:\n%s", rendered)
	}
	if !strings.Contains(rendered, "ran degraded") {
		t.Errorf("partial reviewer failure not reported:\n%s", rendered)
	}

	// The genuinely-no-reviewer case must keep saying so.
	report.Reviewer = "heuristic fallback"
	report.Findings = nil
	if !strings.Contains(RenderMarkdown(report), "deterministic checks only") {
		t.Errorf("a review with no AI findings must still say so:\n%s", RenderMarkdown(report))
	}
}

func firstSentence(text string) string {
	text = strings.TrimSpace(text)
	if index := strings.Index(text, ". "); index > 0 {
		return text[:index]
	}
	return text
}

func findingTitles(findings []Finding) string {
	var b strings.Builder
	for _, finding := range findings {
		b.WriteString("  - " + finding.Title + "\n")
	}
	return b.String()
}

// TestSecondDeduplicationPassDoesNotRestateCorroboration pins a small honesty
// bug in the verbose evidence listing. De-duplication runs twice — once per
// shard inside the panel, once across shards before the judge — and a live run
// printed "Raised independently by 2 reviewers: …" twice under one finding,
// which reads like two separate corroborations of the same claim.
func TestSecondDeduplicationPassDoesNotRestateCorroboration(t *testing.T) {
	first := mergeWithLabels(t, egressA, egressB)
	second := mergeWithLabels(t, append(first, egressC)...)
	if len(second) != 1 {
		t.Fatalf("merged %d findings, want 1", len(second))
	}
	count := 0
	for _, evidence := range second[0].Evidence {
		if evidence.Label == "Corroboration" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("Corroboration evidence appears %d times after two passes: %#v", count, second[0].Evidence)
	}
}
