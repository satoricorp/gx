package codereview

import (
	"context"
	"strings"
	"testing"
)

// symbolAskingJudge abstains on the first round asking where a name is used,
// then decides once the search results arrive.
type symbolAskingJudge struct {
	symbol string
	// secondRound is the verdict returned once the search results arrive.
	secondRound string
	requests    []judgeRequest
}

func (*symbolAskingJudge) Available() bool { return true }

func (j *symbolAskingJudge) Judge(_ context.Context, req judgeRequest) ([]judgeResult, error) {
	j.requests = append(j.requests, req)
	var results []judgeResult
	for _, candidate := range req.Candidates {
		if len(j.requests) == 1 {
			results = append(results, judgeResult{
				CandidateID:   candidate.ID,
				Verdict:       "insufficient_evidence",
				NeededSymbols: []string{j.symbol},
			})
			continue
		}
		results = append(results, judgeResult{
			CandidateID: candidate.ID, Verdict: j.secondRound,
			Impact: impactBreaking, Severity: 3, Confidence: 0.9,
			VerificationNote: "decided from the search results",
		})
	}
	return results, nil
}

// The question that decided four of the thirteen wrong findings — "does every
// caller pass onPaid" — is a search, not a path, so needed_files could never
// express it and the judge guessed. This is that loop closed: it asks, the
// search runs, and the files holding the references come back with the
// matching lines in view.
func TestJudgeCanAskWhereASymbolIsUsed(t *testing.T) {
	root := gitRepoWith(t, map[string]string{
		"page.tsx":               "return <PreviewGateBanner onPaid={onPaid} />;\n",
		"use-lifecycle-state.ts": "export function useLifecycle({ onPaid }) {}\n",
	})
	judge := &symbolAskingJudge{symbol: "onPaid", secondRound: "wrong"}
	reviewContext := ReviewContext{Brief: ReviewBrief{RepoRoot: root}}
	candidates := []Finding{{
		ID: "f1", Title: "`onPaid` is dropped by some caller", File: "page.tsx", Line: 1,
	}}

	outcome := runJudge(context.Background(), judge, reviewContext, candidates)

	if len(j2Requests(judge)) != 2 {
		t.Fatalf("judge called %d times, want a second round after the symbol ask", len(judge.requests))
	}
	second := judge.requests[1]
	var files []string
	for _, snippet := range second.Files {
		files = append(files, snippet.File)
	}
	for _, want := range []string{"page.tsx", "use-lifecycle-state.ts"} {
		if !hasFileNamed(files, want) {
			t.Fatalf("second round files = %v, want the searched references included", files)
		}
	}
	// The candidate came back decided, so it must not still count as an
	// abstention: an answered question is not a missing one.
	if outcome.Abstained != 0 {
		t.Fatalf("Abstained = %d, want the resolved candidate uncounted", outcome.Abstained)
	}
	// Only confirmed verdicts survive, so a refuted candidate is gone. That is
	// the whole point: the search let the judge disprove a claim it had
	// previously been unable to check, and the reader never sees it.
	if len(outcome.Judged) != 0 {
		t.Fatalf("Judged = %d, want the refuted candidate dropped", len(outcome.Judged))
	}
}

// The mirror case. A loop that resolves every ask into a refutation would be
// worse than no loop at all, so a candidate the search supports must survive.
func TestSymbolSearchStillLetsARealFindingThrough(t *testing.T) {
	root := gitRepoWith(t, map[string]string{
		"page.tsx": "return <PreviewGateBanner onPaid={onPaid} />;\n",
	})
	judge := &symbolAskingJudge{symbol: "onPaid", secondRound: "confirmed"}
	reviewContext := ReviewContext{Brief: ReviewBrief{RepoRoot: root}}
	candidates := []Finding{{
		ID: "f1", Title: "`onPaid` is dropped by some caller", File: "page.tsx", Line: 1,
	}}

	outcome := runJudge(context.Background(), judge, reviewContext, candidates)

	if len(outcome.Judged) != 1 {
		t.Fatalf("Judged = %d, want the confirmed candidate kept", len(outcome.Judged))
	}
	if outcome.Abstained != 0 {
		t.Fatalf("Abstained = %d, want the resolved candidate uncounted", outcome.Abstained)
	}
}

func j2Requests(j *symbolAskingJudge) []judgeRequest { return j.requests }

func hasFileNamed(haystack []string, needle string) bool {
	for _, item := range haystack {
		if strings.EqualFold(item, needle) {
			return true
		}
	}
	return false
}
