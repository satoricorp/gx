package telemetry

import (
	"context"
)

// Event names for PostHog capture pipeline telemetry.
const (
	EventCaptureCoverage = "capture.coverage"
	EventMatchRate       = "match.rate"
	EventSessionUploaded = "session.uploaded"
	EventComposeRun      = "compose.run"
)

// CaptureCoverageProps are properties for capture.coverage.
type CaptureCoverageProps struct {
	Repo         string
	RefRange     string
	HunkCoverage float64
	Tier1        int
	Tier2        int
	Tools        []string
}

// MatchRateProps are properties for match.rate.
type MatchRateProps struct {
	Repo           string
	RefRange       string
	AgentPrecision float64
	EligibleHunks  int
	EligibleEvents int
	Tools          []string
}

// SessionUploadedProps are properties for session.uploaded (WP-1b).
type SessionUploadedProps struct {
	SessionID string
	Tool      string
	Bytes     int
	OrgID     string
	Repo      string
	RefRange  string
}

// ComposeRunProps are properties for compose.run (WP-2).
type ComposeRunProps struct {
	Repo          string
	RevisionCount int
	Warnings      int
	HunkCoverage  float64
	Tools         []string
}

// Client emits capture telemetry events.
type Client interface {
	EmitCaptureCoverage(ctx context.Context, props CaptureCoverageProps)
	EmitMatchRate(ctx context.Context, props MatchRateProps)
	EmitSessionUploaded(ctx context.Context, props SessionUploadedProps)
	EmitComposeRun(ctx context.Context, props ComposeRunProps)
}
