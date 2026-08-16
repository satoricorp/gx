package codereview

import (
	"context"
	"strings"
	"testing"
)

// scriptedTransport replays a fixed sequence of replies and records the
// system prompt each call received, so a test can assert both that a retry
// happened and that it carried the JSON-only reminder.
type scriptedTransport struct {
	replies []string
	systems []string
	stop    string // stop reason for every reply; "" means end_turn
}

func (s *scriptedTransport) complete(_ context.Context, _, system, _ string, _ int) (bedrockCompletion, error) {
	s.systems = append(s.systems, system)
	i := len(s.systems) - 1
	if i >= len(s.replies) {
		i = len(s.replies) - 1
	}
	stop := s.stop
	if stop == "" {
		stop = "end_turn"
	}
	return bedrockCompletion{Text: s.replies[i], StopReason: stop}, nil
}

func (s *scriptedTransport) detail() string { return "scripted" }

func TestReviewForSummaryRetriesOnceWhenReplyHasNoJSON(t *testing.T) {
	transport := &scriptedTransport{replies: []string{
		// What a leg actually sent on a binary-only diff: prose, no JSON.
		"I reviewed the change. The file is a compiled Mach-O binary, so there is nothing to review at the source level.",
		`{"recommendations":[]}`,
	}}
	r := newBedrockReviewer(transport, "us.anthropic.claude-haiku-4-5-20251001-v1:0")
	out, err := r.ReviewForSummary(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("second reply was valid JSON; want success, got %v", err)
	}
	if len(out.Findings) != 0 {
		t.Fatalf("expected empty findings, got %d", len(out.Findings))
	}
	if len(transport.systems) != 2 {
		t.Fatalf("expected exactly one retry (2 calls), got %d", len(transport.systems))
	}
	if strings.HasPrefix(transport.systems[0], reviewJSONOnlyReminder) {
		t.Fatalf("first call must use the plain system prompt")
	}
	if !strings.HasPrefix(transport.systems[1], reviewJSONOnlyReminder) {
		t.Fatalf("retry must lead with the JSON-only reminder, got:\n%s", transport.systems[1][:120])
	}
}

func TestReviewForSummaryGivesUpAfterOneRetry(t *testing.T) {
	transport := &scriptedTransport{replies: []string{"prose only", "still prose"}}
	r := newBedrockReviewer(transport, "us.anthropic.claude-haiku-4-5-20251001-v1:0")
	if _, err := r.ReviewForSummary(context.Background(), ReviewBrief{}); err == nil {
		t.Fatalf("two unparseable replies must surface an error")
	}
	if len(transport.systems) != 2 {
		t.Fatalf("must stop after one retry, made %d calls", len(transport.systems))
	}
}

func TestReviewForSummaryDoesNotRetryAValidFirstReply(t *testing.T) {
	transport := &scriptedTransport{replies: []string{"```json\n{\"recommendations\":[]}\n```\n\nSome trailing summary prose."}}
	r := newBedrockReviewer(transport, "us.anthropic.claude-haiku-4-5-20251001-v1:0")
	if _, err := r.ReviewForSummary(context.Background(), ReviewBrief{}); err != nil {
		t.Fatalf("fenced JSON with trailing prose parses today and must keep doing so: %v", err)
	}
	if len(transport.systems) != 1 {
		t.Fatalf("a parseable reply must not be retried, made %d calls", len(transport.systems))
	}
}

func TestReviewForSummarySalvagesCompleteFindingsFromATruncatedReply(t *testing.T) {
	// A reply the output cap cut mid-way through the third recommendation:
	// two complete elements, then a dangling string. What Sonnet's whole-repo
	// shards actually looked like — 13 complete findings thrown away because
	// the 14th was cut.
	truncated := `{"recommendations":[` +
		`{"title":"t1","summary":"s1","benefit":"b1","recommendation":"r1","kind":"defect","strength":"Strong","file":"a.go","line":3},` +
		`{"title":"t2","summary":"s2","benefit":"b2","recommendation":"r2","kind":"defect","strength":"Strong","file":"b.go","line":9},` +
		`{"title":"t3","summary":"cut off mid sen`
	transport := &scriptedTransport{replies: []string{truncated}, stop: "max_tokens"}
	r := newBedrockReviewer(transport, "us.anthropic.claude-sonnet-4-6")
	out, err := r.ReviewForSummary(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("salvageable truncation must not be an error, got %v", err)
	}
	if len(out.Findings) != 2 {
		t.Fatalf("expected the 2 complete findings, got %d", len(out.Findings))
	}
	if out.Findings[0].Title != "t1" || out.Findings[1].Title != "t2" {
		t.Fatalf("wrong findings salvaged: %q %q", out.Findings[0].Title, out.Findings[1].Title)
	}
	if out.Truncated == "" || !strings.Contains(out.Truncated, "output cap") {
		t.Fatalf("truncation must be reported on the summary, got %q", out.Truncated)
	}
	if len(transport.systems) != 1 {
		t.Fatalf("a truncated reply must not be retried at the same cap, made %d calls", len(transport.systems))
	}
}

func TestSalvageTruncatedRecommendationsEdgeCases(t *testing.T) {
	if got := salvageTruncatedRecommendations(`no json here`); got != nil {
		t.Fatalf("no recommendations key → nil, got %v", got)
	}
	if got := salvageTruncatedRecommendations(`{"recommendations":[`); got != nil {
		t.Fatalf("empty truncated array → nil, got %v", got)
	}
	// A "}" inside a string must not end an element early.
	one := `{"recommendations":[{"title":"has } brace","summary":"s","benefit":"b","recommendation":"r"},{"title":"cut`
	got := salvageTruncatedRecommendations(one)
	if len(got) != 1 || got[0].Title != "has } brace" {
		t.Fatalf("brace inside string mishandled: %+v", got)
	}
}
