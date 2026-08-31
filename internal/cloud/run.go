package cloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// A run is the unit gx Cloud meters: one `gx review`, however many model calls
// it makes. Free accounts get a fixed number of runs; past that, Cloud AI
// answers 402 with a checkout URL until the user subscribes.
//
// The run key is minted once per review and travels in the context, so every
// cloud call the review makes — reviewer legs, judge, dedupe — carries the
// same X-GX-Run header and the server counts them as one.

type runKeyContextKey struct{}

// RunHeader carries the run key on every gx Cloud request made inside a run.
const RunHeader = "X-GX-Run"

// NewRunKey mints a key for one run.
func NewRunKey() string {
	return "run_" + uuid.NewString()
}

// WithRunKey returns a context whose cloud requests are attributed to runKey.
func WithRunKey(ctx context.Context, runKey string) context.Context {
	runKey = strings.TrimSpace(runKey)
	if runKey == "" {
		return ctx
	}
	return context.WithValue(ctx, runKeyContextKey{}, runKey)
}

// RunKeyFrom reports the run key on ctx, or "" outside a run.
func RunKeyFrom(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	key, _ := ctx.Value(runKeyContextKey{}).(string)
	return key
}

// RunReservation is gx Cloud's answer to "may this run go ahead?".
type RunReservation struct {
	Allowed     bool   `json:"allowed"`
	Reason      string `json:"reason"`
	Used        int    `json:"used"`
	Limit       int    `json:"limit"`
	Remaining   int    `json:"remaining"`
	CheckoutURL string `json:"checkout_url"`
}

// FreeRunsNote is a one-line status for the progress log, or "" when there is
// nothing worth saying (subscribers, unknown limits).
func (r RunReservation) FreeRunsNote() string {
	if r.Reason == "subscribed" || r.Limit <= 0 {
		return ""
	}
	if r.Remaining == 1 {
		return fmt.Sprintf("1 free run left after this one (%d/%d used)", r.Used, r.Limit)
	}
	return fmt.Sprintf("%d free runs left after this one (%d/%d used)", r.Remaining, r.Used, r.Limit)
}

// ReserveRun asks gx Cloud to count one run under the key on ctx before any
// model work starts. A 402 comes back as a PaymentRequiredError whose message
// names the checkout URL; other failures are transport errors, and the caller
// may proceed since every model call is gated server-side anyway.
func (c *Client) ReserveRun(ctx context.Context) (RunReservation, error) {
	var result RunReservation
	if c == nil || c.url == "" {
		return result, fmt.Errorf("gx cloud base URL is not configured")
	}
	runKey := RunKeyFrom(ctx)
	if runKey == "" {
		return result, fmt.Errorf("no run key on context")
	}
	payload := struct {
		Kind   string `json:"kind"`
		RunKey string `json:"run_key"`
	}{Kind: "review", RunKey: runKey}
	if err := c.postJSON(ctx, "/v1/runs/reserve", payload, &result); err != nil {
		return RunReservation{}, err
	}
	return result, nil
}
