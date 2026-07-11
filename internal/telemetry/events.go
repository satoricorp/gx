package telemetry

import (
	"context"
)

// Event names for PostHog capture pipeline telemetry.
const (
	EventCaptureCoverage   = "capture.coverage"
	EventMatchRate         = "match.rate"
	EventSessionUploaded   = "session.uploaded"
	EventComposeRun        = "compose.run"
	EventCLIInstall        = "cli.install"
	EventCLIInitRun        = "cli.init.run"
	EventCLICommitRun      = "cli.commit.run"
	EventCLIGenerateRun    = "cli.generate.run"
	EventCLIGeneratePrompt = "cli.generate.prompt"
	EventCLIPushRun        = "cli.push.run"
	EventCLISyncRun        = "cli.sync.run"
	EventCLIReportSent     = "cli.report.sent"
	EventCLIReviewRun      = "cli.review.run"
	EventCLIAuthLogin      = "cli.auth.login"
	EventCLIAuthLogout     = "cli.auth.logout"
	EventSchemaDrift       = "capture.schema_drift"
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

// SchemaDriftProps are properties for capture.schema_drift.
type SchemaDriftProps struct {
	Tool  string
	Paths []string
	Count int
}

// Client emits capture telemetry events.
type Client interface {
	EmitEvent(ctx context.Context, event string, properties map[string]any)
	EmitCaptureCoverage(ctx context.Context, props CaptureCoverageProps)
	EmitMatchRate(ctx context.Context, props MatchRateProps)
	EmitSessionUploaded(ctx context.Context, props SessionUploadedProps)
	EmitComposeRun(ctx context.Context, props ComposeRunProps)
	EmitSchemaDrift(ctx context.Context, props SchemaDriftProps)
}
