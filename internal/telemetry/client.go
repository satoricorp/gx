package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/satoricorp/lgtm/internal/buildconfig"
)

const defaultPostHogHost = "https://f.lgtm.cx"

// ClientImpl sends events to PostHog when configured.
type ClientImpl struct {
	apiKey string
	host   string
	http   *http.Client
}

// NewFromEnv returns a PostHog client or a no-op when LGTM_POSTHOG_KEY is unset.
func NewFromEnv() Client {
	key := postHogKey()
	if key == "" {
		return NopClient{}
	}
	host := postHogHost()
	if host == "" {
		host = defaultPostHogHost
	}
	host = strings.TrimRight(host, "/")
	return &ClientImpl{
		apiKey: key,
		host:   host,
		http:   &http.Client{Timeout: 5 * time.Second},
	}
}

// NopClient discards telemetry events.
type NopClient struct{}

func (NopClient) EmitEvent(context.Context, string, map[string]any)         {}
func (NopClient) EmitCaptureCoverage(context.Context, CaptureCoverageProps) {}
func (NopClient) EmitMatchRate(context.Context, MatchRateProps)             {}
func (NopClient) EmitSessionUploaded(context.Context, SessionUploadedProps) {}
func (NopClient) EmitComposeRun(context.Context, ComposeRunProps)           {}
func (NopClient) EmitSchemaDrift(context.Context, SchemaDriftProps)         {}

func Configured() bool {
	return postHogKey() != ""
}

func (c *ClientImpl) EmitEvent(ctx context.Context, event string, properties map[string]any) {
	c.capture(ctx, event, properties)
}

func (c *ClientImpl) EmitCaptureCoverage(ctx context.Context, props CaptureCoverageProps) {
	c.capture(ctx, EventCaptureCoverage, map[string]any{
		"repo":          props.Repo,
		"ref_range":     props.RefRange,
		"hunk_coverage": props.HunkCoverage,
		"tier1":         props.Tier1,
		"tier2":         props.Tier2,
		"tools":         props.Tools,
	})
}

func (c *ClientImpl) EmitMatchRate(ctx context.Context, props MatchRateProps) {
	c.capture(ctx, EventMatchRate, map[string]any{
		"repo":            props.Repo,
		"ref_range":       props.RefRange,
		"agent_precision": props.AgentPrecision,
		"eligible_hunks":  props.EligibleHunks,
		"eligible_events": props.EligibleEvents,
		"tools":           props.Tools,
	})
}

func (c *ClientImpl) EmitSessionUploaded(ctx context.Context, props SessionUploadedProps) {
	payload := map[string]any{
		"session_id": props.SessionID,
		"tool":       props.Tool,
		"bytes":      props.Bytes,
		"repo":       props.Repo,
		"ref_range":  props.RefRange,
	}
	if props.OrgID != "" {
		payload["org_id"] = props.OrgID
	}
	c.capture(ctx, EventSessionUploaded, payload)
}

func (c *ClientImpl) EmitComposeRun(ctx context.Context, props ComposeRunProps) {
	c.capture(ctx, EventComposeRun, map[string]any{
		"repo":           props.Repo,
		"revision_count": props.RevisionCount,
		"warnings":       props.Warnings,
		"hunk_coverage":  props.HunkCoverage,
		"tools":          props.Tools,
	})
}

func (c *ClientImpl) EmitSchemaDrift(ctx context.Context, props SchemaDriftProps) {
	if props.Count == 0 {
		return
	}
	c.capture(ctx, EventSchemaDrift, map[string]any{
		"tool":  props.Tool,
		"paths": props.Paths,
		"count": props.Count,
	})
}

func (c *ClientImpl) capture(ctx context.Context, event string, properties map[string]any) {
	captureProperties := ProductProperties(ctx, properties)
	body, err := json.Marshal(map[string]any{
		"api_key":    c.apiKey,
		"event":      event,
		"properties": captureProperties,
	})
	if err != nil {
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host+"/capture/", bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return
	}
	_ = resp.Body.Close()
}

func postHogKey() string {
	return strings.TrimSpace(buildconfig.PostHogKeyFromEnvOrEmbedded())
}

func postHogHost() string {
	return strings.TrimSpace(buildconfig.PostHogHostFromEnvOrEmbedded())
}
