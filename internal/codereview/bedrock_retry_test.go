package codereview

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// scriptedRoundTripper returns one canned response per call, in order.
type scriptedRoundTripper struct {
	calls     int
	responses []*http.Response
}

func (s *scriptedRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	if s.calls >= len(s.responses) {
		panic("scriptedRoundTripper: more calls than scripted responses")
	}
	resp := s.responses[s.calls]
	s.calls++
	return resp, nil
}

func cannedResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func fastRetryBackoffs(t *testing.T) {
	t.Helper()
	saved := bedrockRetryBackoffs
	bedrockRetryBackoffs = []time.Duration{time.Millisecond, time.Millisecond, time.Millisecond}
	t.Cleanup(func() { bedrockRetryBackoffs = saved })
}

func retryTestTransport(rt http.RoundTripper) *directBedrockTransport {
	return &directBedrockTransport{
		region:    "us-west-2",
		accessKey: "AKIAEXAMPLE",
		secretKey: "secret",
		client:    &http.Client{Transport: rt},
	}
}

const bedrockOKBody = `{"content":[{"type":"text","text":"{\"ok\":true}"}],"stop_reason":"end_turn"}`

func TestBedrockCompleteRetriesThrottleAndServerError(t *testing.T) {
	fastRetryBackoffs(t)
	rt := &scriptedRoundTripper{responses: []*http.Response{
		cannedResponse(http.StatusTooManyRequests, `{"message":"Too many requests, please wait before trying again."}`),
		cannedResponse(http.StatusInternalServerError, `{"message":"Internal Server Error"}`),
		cannedResponse(http.StatusOK, bedrockOKBody),
	}}
	completion, err := retryTestTransport(rt).complete(context.Background(), "us.anthropic.claude-sonnet-4-6", "system", "input", 100)
	if err != nil {
		t.Fatalf("complete() error = %v", err)
	}
	if rt.calls != 3 {
		t.Fatalf("calls = %d, want 3 (two retries then success)", rt.calls)
	}
	if completion.Text != `{"ok":true}` {
		t.Fatalf("Text = %q", completion.Text)
	}
}

func TestBedrockCompleteDoesNotRetryValidationError(t *testing.T) {
	fastRetryBackoffs(t)
	rt := &scriptedRoundTripper{responses: []*http.Response{
		cannedResponse(http.StatusBadRequest, `{"message":"Input is too long for requested model."}`),
	}}
	_, err := retryTestTransport(rt).complete(context.Background(), "us.anthropic.claude-sonnet-4-6", "system", "input", 100)
	if err == nil {
		t.Fatal("expected error")
	}
	if rt.calls != 1 {
		t.Fatalf("calls = %d, want 1 (400s are not retryable)", rt.calls)
	}
}

func TestBedrockCompleteGivesUpAfterBackoffsExhausted(t *testing.T) {
	fastRetryBackoffs(t)
	throttle := func() *http.Response {
		return cannedResponse(http.StatusTooManyRequests, `{"message":"throttled"}`)
	}
	rt := &scriptedRoundTripper{responses: []*http.Response{throttle(), throttle(), throttle(), throttle()}}
	_, err := retryTestTransport(rt).complete(context.Background(), "us.anthropic.claude-sonnet-4-6", "system", "input", 100)
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if want := 1 + len(bedrockRetryBackoffs); rt.calls != want {
		t.Fatalf("calls = %d, want %d", rt.calls, want)
	}
	if !strings.Contains(err.Error(), "throttled") {
		t.Fatalf("err = %v, want the throttle detail preserved", err)
	}
}
