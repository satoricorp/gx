package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/satoricorp/lgtm/internal/capture/orchestrator"
	"github.com/satoricorp/lgtm/internal/hooks"
	"github.com/satoricorp/lgtm/internal/storage"
)

// TestCapturePushSurfacesHookError pins the diagnostic contract of the push
// command: a failed run is reported on stderr, and the staging line still
// prints so its empty `extract=` keeps meaning "the run errored" rather than
// "the run found nothing".
func TestCapturePushSurfacesHookError(t *testing.T) {
	restore := runPushHook
	runPushHook = func(context.Context, hooks.PushOptions) (hooks.PushOutcome, error) {
		return hooks.PushOutcome{}, errors.New("open capture stager: disk is full")
	}
	t.Cleanup(func() { runPushHook = restore })

	var stdout, stderr bytes.Buffer
	cmd := newCapturePushCommand(context.Background())
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(nil)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v, want nil so the hook never blocks the push", err)
	}

	if !strings.Contains(stderr.String(), "disk is full") {
		t.Fatalf("stderr = %q, want the hook error surfaced", stderr.String())
	}
	if !strings.Contains(stdout.String(), "capture staged extract= ") {
		t.Fatalf("stdout = %q, want the staging line with an empty extract id", stdout.String())
	}
}

// TestCapturePushWarnsAboutSkippedTools proves a tool that could not be
// discovered reaches the operator instead of quietly shrinking the capture.
func TestCapturePushWarnsAboutSkippedTools(t *testing.T) {
	restore := runPushHook
	runPushHook = func(context.Context, hooks.PushOptions) (hooks.PushOutcome, error) {
		return hooks.PushOutcome{Result: orchestrator.Result{
			StagedExtractID:   "ext-1",
			DiscoveryProblems: []string{"claude: read /home/u/.claude/projects: permission denied"},
		}}, nil
	}
	t.Cleanup(func() { runPushHook = restore })

	var stdout, stderr bytes.Buffer
	cmd := newCapturePushCommand(context.Background())
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(nil)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr.String(), "permission denied") {
		t.Fatalf("stderr = %q, want the skipped tool reported", stderr.String())
	}
}

// TestCapturePushDistinguishesSkipsFromFailures pins the meaning of the
// staging line. RunPush returns an empty outcome on four separate paths, and
// the command printed the staging line unconditionally, so a paused capture and
// a repo opted out with `git config lgtm.enabled false` emitted the byte-for-byte
// signature of an errored run — `capture staged extract= sessions=0 …` — with
// nothing on stderr. An empty extract id has to mean exactly one thing.
func TestCapturePushDistinguishesSkipsFromFailures(t *testing.T) {
	restore := runPushHook
	runPushHook = func(context.Context, hooks.PushOptions) (hooks.PushOutcome, error) {
		return hooks.PushOutcome{
			SkipReason: "this repository opted out (git config lgtm.enabled false)",
		}, nil
	}
	t.Cleanup(func() { runPushHook = restore })

	var stdout, stderr bytes.Buffer
	cmd := newCapturePushCommand(context.Background())
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(nil)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stdout.String(), "capture staged extract=") {
		t.Fatalf("stdout = %q, want a deliberate skip not rendered as an errored staging line", stdout.String())
	}
	if !strings.Contains(stdout.String(), "capture skipped: this repository opted out") {
		t.Fatalf("stdout = %q, want the skip named", stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want an opt-out to warn about nothing", stderr.String())
	}
}

// TestCapturePushSurfacesFailingUploads is the last link in the only chain by
// which a failed upload can reach a human: `lgtm capture sync` runs detached with
// both streams on os.DevNull, so the push that starts it cannot report its
// result and the next push has to. Without this, 32 rows could 400 on every
// push for weeks behind a clean-looking capture line.
func TestCapturePushSurfacesFailingUploads(t *testing.T) {
	restore := runPushHook
	runPushHook = func(context.Context, hooks.PushOptions) (hooks.PushOutcome, error) {
		return hooks.PushOutcome{
			Result: orchestrator.Result{StagedExtractID: "ext-1"},
			UploadFailures: storage.CaptureUploadFailures{
				Sessions:      32,
				ExhaustedRows: 32,
				LastError:     `status 400: {"error":"sessionId, tool, and content are required"}`,
			},
		}, nil
	}
	t.Cleanup(func() { runPushHook = restore })

	var stdout, stderr bytes.Buffer
	cmd := newCapturePushCommand(context.Background())
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(nil)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr.String(), "32 capture upload(s) failed") {
		t.Fatalf("stderr = %q, want the failing uploads reported", stderr.String())
	}
	if !strings.Contains(stderr.String(), "sessionId, tool, and content are required") {
		t.Fatalf("stderr = %q, want the server's own reason readable", stderr.String())
	}
}

// TestCapturePushSurfacesDroppedWork covers the two errors the hook used to
// discard outright: a failed shareable-marking pass, and a publication that was
// never enqueued. Both leave the push looking successful while the PR silently
// gets no lgtm artifact and the staged rows never upload.
func TestCapturePushSurfacesDroppedWork(t *testing.T) {
	restore := runPushHook
	runPushHook = func(context.Context, hooks.PushOptions) (hooks.PushOutcome, error) {
		return hooks.PushOutcome{
			Result:           orchestrator.Result{StagedExtractID: "ext-1"},
			ShareableError:   "database is locked",
			PublicationError: "resolve head commit: fatal: bad object",
		}, nil
	}
	t.Cleanup(func() { runPushHook = restore })

	var stdout, stderr bytes.Buffer
	cmd := newCapturePushCommand(context.Background())
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(nil)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr.String(), "database is locked") {
		t.Fatalf("stderr = %q, want the marking failure reported", stderr.String())
	}
	if !strings.Contains(stderr.String(), "queued no lgtm review artifact") {
		t.Fatalf("stderr = %q, want the missing publication reported", stderr.String())
	}
}
