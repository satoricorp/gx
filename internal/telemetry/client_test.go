package telemetry_test

import (
	"context"
	"testing"

	"github.com/satoricorp/gx/internal/telemetry"
)

func TestNopClient(t *testing.T) {
	t.Setenv("GX_POSTHOG_KEY", "")
	client := telemetry.NewFromEnv()
	ctx := context.Background()
	client.EmitCaptureCoverage(ctx, telemetry.CaptureCoverageProps{HunkCoverage: 0.9})
	client.EmitMatchRate(ctx, telemetry.MatchRateProps{AgentPrecision: 0.8})
	client.EmitSessionUploaded(ctx, telemetry.SessionUploadedProps{SessionID: "s1"})
	client.EmitComposeRun(ctx, telemetry.ComposeRunProps{HunkCoverage: 0.5})
}

func TestClientImplImplementsInterface(t *testing.T) {
	var _ telemetry.Client = telemetry.NopClient{}
}
