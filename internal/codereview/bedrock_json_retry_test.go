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
}

func (s *scriptedTransport) complete(_ context.Context, _, system, _ string, _ int) (bedrockCompletion, error) {
	s.systems = append(s.systems, system)
	i := len(s.systems) - 1
	if i >= len(s.replies) {
		i = len(s.replies) - 1
	}
	return bedrockCompletion{Text: s.replies[i], StopReason: "end_turn"}, nil
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
