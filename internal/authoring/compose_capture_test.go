package authoring

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/matcher"
	"github.com/satoricorp/gx/internal/capture/parsers"
)

var errUnexpectedCaptureParse = errors.New("unexpected capture parse")

func TestSessionIDsFromHunkLinks(t *testing.T) {
	links := []matcher.HunkLink{
		{HunkID: "h1", SessionID: "s1", Tier: matcher.TierExact, Authorship: matcher.AuthorshipAgent},
		{HunkID: "h2", SessionID: "s1", Tier: matcher.TierFuzzy, Authorship: matcher.AuthorshipAgent},
		{HunkID: "h3", SessionID: "s2", Tier: matcher.TierTemporal, Authorship: matcher.AuthorshipUnknown},
		{HunkID: "h4", Authorship: matcher.AuthorshipHuman},
	}
	got := sessionIDsFromHunkLinks(links)
	if len(got) != 1 || got[0] != "s1" {
		t.Fatalf("sessionIDsFromHunkLinks() = %#v, want [s1]", got)
	}
}

func TestHunkLinksForRevision(t *testing.T) {
	links := []matcher.HunkLink{
		{HunkID: "working:internal/a.go:1-3"},
		{HunkID: "working:internal/b.go:1-2"},
	}
	got := hunkLinksForRevision(links, []string{"internal/a.go"})
	if len(got) != 1 || got[0].HunkID != links[0].HunkID {
		t.Fatalf("hunkLinksForRevision() = %#v", got)
	}
}

func TestHunkLinksForRevisionPrefersRevisionHunks(t *testing.T) {
	links := []matcher.HunkLink{
		{HunkID: "h1"},
		{HunkID: "h2"},
	}
	got := hunkLinksForRevision(links, []string{"internal/a.go"}, []string{"h2"})
	if len(got) != 1 || got[0].HunkID != "h2" {
		t.Fatalf("hunkLinksForRevision() = %#v, want h2 only", got)
	}
}

func TestApplyCaptureEvidenceToProposalDerivesRevisionSessionIDs(t *testing.T) {
	proposal := DemuxProposal{
		Revisions: []RevisionProposal{
			{ID: "u1", Files: []string{"internal/a.go"}, UseHunks: true, HunkIDs: []string{"h1"}},
			{ID: "u2", Files: []string{"internal/b.go"}, UseHunks: true, HunkIDs: []string{"h2"}},
			{ID: "u3", Files: []string{"internal/c.go"}, UseHunks: true, HunkIDs: []string{"h3"}},
		},
	}
	evidence := ComposeCaptureResult{
		HunkLinks: []matcher.HunkLink{
			{HunkID: "h1", SessionID: "s-a", Tier: matcher.TierExact, Authorship: matcher.AuthorshipAgent, Tool: capture.ToolCodex},
			{HunkID: "h2", SessionID: "s-human", Tier: matcher.TierExact, Authorship: matcher.AuthorshipHuman, Tool: capture.ToolClaude},
		},
		HunkCoverage: 0.5,
		Status:       "matcher",
		Warnings:     []string{"capture: tailed codex session s-a"},
	}

	got := applyCaptureEvidenceToProposal(proposal, evidence, []string{"fallback-session"}, "explicit")

	if got.HunkCoverage != 0.5 || len(got.HunkLinks) != 2 {
		t.Fatalf("proposal capture = coverage %.2f links %#v, want populated", got.HunkCoverage, got.HunkLinks)
	}
	if len(got.Warnings) != 1 || got.Warnings[0] != "capture: tailed codex session s-a" {
		t.Fatalf("proposal warnings = %#v", got.Warnings)
	}
	if len(got.CaptureTools) != 3 || got.CaptureTools[0] != capture.ToolClaude || got.CaptureTools[1] != capture.ToolCodex || got.CaptureTools[2] != capture.ToolCursor {
		t.Fatalf("capture tools = %#v, want normalized defaults", got.CaptureTools)
	}
	if got.Revisions[0].ProvenanceStatus != "matcher" || len(got.Revisions[0].SessionIDs) != 1 || got.Revisions[0].SessionIDs[0] != "s-a" {
		t.Fatalf("revision u1 provenance = %q/%#v, want matched session", got.Revisions[0].ProvenanceStatus, got.Revisions[0].SessionIDs)
	}
	if got.Revisions[1].ProvenanceStatus != "explicit" || len(got.Revisions[1].SessionIDs) != 1 || got.Revisions[1].SessionIDs[0] != "fallback-session" {
		t.Fatalf("revision u2 provenance = %q/%#v, want fallback for human link", got.Revisions[1].ProvenanceStatus, got.Revisions[1].SessionIDs)
	}
	if got.Revisions[2].ProvenanceStatus != "explicit" || len(got.Revisions[2].SessionIDs) != 1 || got.Revisions[2].SessionIDs[0] != "fallback-session" {
		t.Fatalf("revision u3 provenance = %q/%#v, want fallback for unmatched revision", got.Revisions[2].ProvenanceStatus, got.Revisions[2].SessionIDs)
	}
}

func TestParseCaptureSessionsBoundedTailsOversizedJSONL(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.jsonl")
	prefix := strings.Repeat("x", captureMatchJSONLTailBytes+256)
	if err := os.WriteFile(path, []byte(prefix+"\nrecent\n"), 0o600); err != nil {
		t.Fatalf("write session: %v", err)
	}

	events, warnings, err := parseCaptureSessionsBounded(context.Background(), []parsers.DiscoveredSession{{
		Tool:      "fake",
		Path:      path,
		Kind:      parsers.SessionKindJSONL,
		SessionID: "fake-session",
	}}, "", []parsers.Parser{lineCaptureParser{}}, 10)
	if err != nil {
		t.Fatalf("parseCaptureSessionsBounded: %v", err)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "capture: tailed fake session fake-session") {
		t.Fatalf("warnings = %#v", warnings)
	}
	if len(events) != 1 || events[0].NewText != "recent" {
		t.Fatalf("events = %#v, want recent line only", events)
	}
}

func TestAppendCaptureWarningsDedupes(t *testing.T) {
	got := appendCaptureWarnings([]string{"keep", "capture: one"}, []string{"capture: one", "capture: two", ""})
	want := []string{"keep", "capture: one", "capture: two"}
	if len(got) != len(want) {
		t.Fatalf("warnings = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("warnings = %#v, want %#v", got, want)
		}
	}
}

func TestParseCaptureSessionsBoundedSkipsOversizedCursorVSCDBWithWarning(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.vscdb")
	if err := os.WriteFile(path, []byte(strings.Repeat("x", captureMatchMaxVSCDBBytes+1)), 0o600); err != nil {
		t.Fatalf("write vscdb: %v", err)
	}

	events, warnings, err := parseCaptureSessionsBounded(context.Background(), []parsers.DiscoveredSession{{
		Tool: "fake",
		Path: path,
		Kind: parsers.SessionKindCursorVSCDB,
	}}, "", []parsers.Parser{failingCaptureParser{}}, 10)
	if err != nil {
		t.Fatalf("parseCaptureSessionsBounded: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("events = %#v, want none", events)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "capture: skipped cursor state.vscdb") {
		t.Fatalf("warnings = %#v", warnings)
	}
}

func TestCapCaptureSessionsKeepsMostRecent(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old.jsonl")
	newPath := filepath.Join(dir, "new.jsonl")
	for _, path := range []string{oldPath, newPath} {
		if err := os.WriteFile(path, []byte("{}\n"), 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
	if err := os.Chtimes(oldPath, time.Unix(100, 0), time.Unix(100, 0)); err != nil {
		t.Fatalf("chtimes old: %v", err)
	}
	if err := os.Chtimes(newPath, time.Unix(200, 0), time.Unix(200, 0)); err != nil {
		t.Fatalf("chtimes new: %v", err)
	}
	got := capCaptureSessions([]parsers.DiscoveredSession{
		{Tool: "fake", Path: oldPath},
		{Tool: "fake", Path: newPath},
	}, 1)
	if len(got) != 1 || got[0].Path != newPath {
		t.Fatalf("capCaptureSessions() = %#v, want newest path", got)
	}
}

type lineCaptureParser struct{}

func (lineCaptureParser) Tool() string { return "fake" }

func (lineCaptureParser) ParseFile(path string, repoRoot string) ([]capture.SessionEvent, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var events []capture.SessionEvent
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		events = append(events, capture.SessionEvent{
			SessionID: "fake-session",
			Tool:      "fake",
			Kind:      capture.KindEdit,
			NewText:   line,
		})
	}
	return events, scanner.Err()
}

type failingCaptureParser struct{}

func (failingCaptureParser) Tool() string { return "fake" }

func (failingCaptureParser) ParseFile(path string, repoRoot string) ([]capture.SessionEvent, error) {
	return nil, errUnexpectedCaptureParse
}

func TestMatchWorkingCopyForGeneratePlanningCancellationWarning(t *testing.T) {
	previous := matchWorkingCopyForGenerate
	matchWorkingCopyForGenerate = func(ctx context.Context, repoRoot string, hunks []HunkRange, tools []string) (ComposeCaptureResult, error) {
		<-ctx.Done()
		return ComposeCaptureResult{}, ctx.Err()
	}
	defer func() { matchWorkingCopyForGenerate = previous }()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	capture, warning := matchWorkingCopyForGeneratePlanning(ctx, "/repo", []HunkRange{{ID: "h1", File: "a.go", Patch: "+x\n"}})

	if len(capture.HunkLinks) != 0 || capture.HunkCoverage != 0 {
		t.Fatalf("capture = %#v, want empty on cancellation", capture)
	}
	if !strings.Contains(warning, "capture matching skipped") || !strings.Contains(warning, "context canceled") {
		t.Fatalf("warning = %q, want cancellation diagnostic", warning)
	}
}

func TestMatchWorkingCopyForGeneratePlanningErrorWarning(t *testing.T) {
	previous := matchWorkingCopyForGenerate
	matchWorkingCopyForGenerate = func(ctx context.Context, repoRoot string, hunks []HunkRange, tools []string) (ComposeCaptureResult, error) {
		return ComposeCaptureResult{}, context.Canceled
	}
	defer func() { matchWorkingCopyForGenerate = previous }()

	_, warning := matchWorkingCopyForGeneratePlanning(context.Background(), "/repo", []HunkRange{{ID: "h1", File: "a.go", Patch: "+x\n"}})

	if !strings.Contains(warning, "capture matching skipped") {
		t.Fatalf("warning = %q, want capture diagnostic", warning)
	}
}

func TestAddedLinesFromPatch(t *testing.T) {
	patch := "@@ -1 +1,2 @@\n-old\n+new\n+line\n"
	got := addedLinesFromPatch(patch)
	want := []string{"new", "line"}
	if len(got) != len(want) {
		t.Fatalf("addedLinesFromPatch() = %#v, want %#v", got, want)
	}
}

func TestRemapHunkLinkIDs(t *testing.T) {
	demuxHunks := []HunkRange{{ID: "h1", File: "a.go", Patch: "+x\n"}}
	eligible := []capture.CommitHunk{{FilePath: "a.go"}}
	links := []matcher.HunkLink{{HunkID: "working:a.go:1-1"}}
	got := remapHunkLinkIDs(links, demuxHunks, eligible)
	if got[0].HunkID != "h1" {
		t.Fatalf("remapHunkLinkIDs() hunk id = %q, want h1", got[0].HunkID)
	}
}
