package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/exclude"
	captureextract "github.com/satoricorp/gx/internal/capture/extract"
	capturegit "github.com/satoricorp/gx/internal/capture/git"
	"github.com/satoricorp/gx/internal/capture/matcher"
	"github.com/satoricorp/gx/internal/capture/parsers"
	"github.com/satoricorp/gx/internal/capture/parsers/claude"
	"github.com/satoricorp/gx/internal/capture/parsers/codex"
	cursorparser "github.com/satoricorp/gx/internal/capture/parsers/cursor"
	"github.com/satoricorp/gx/internal/capture/redact"
	"github.com/satoricorp/gx/internal/capture/report"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/telemetry"
	"github.com/satoricorp/gx/internal/uploadauth"
	"github.com/satoricorp/gx/internal/version"
)

const defaultTimeSlack = 24 * time.Hour

// RunOptions configures one capture pipeline run.
type RunOptions struct {
	RepoRoot    string
	Base        string
	Head        string
	HomeDir     string
	Tools       []string
	StageOnly   bool
	TimeSlack   time.Duration
	ClaudeDir   string
	CodexDir    string
	CursorVSCDB string
	DB          storage.CaptureStager
	Telemetry   telemetry.Client
}

// Result summarizes one orchestrator run.
type Result struct {
	RefRange        string             `json:"refRange"`
	EligibleHunks   int                `json:"eligibleHunks"`
	EligibleEvents  int                `json:"eligibleEvents"`
	HunkCoverage    float64            `json:"hunkCoverage"`
	Tier1Hunks      int                `json:"tier1Hunks"`
	Tier2Hunks      int                `json:"tier2Hunks"`
	HunkLinks       []matcher.HunkLink `json:"hunkLinks"`
	StagedExtractID string             `json:"stagedExtractID,omitempty"`
	StagedSessions  int                `json:"stagedSessions"`
	Tools           []string           `json:"tools"`
	UploadError     string             `json:"uploadError,omitempty"`
}

// StagedExtract is the redacted extract payload written to SQLite.
type StagedExtract struct {
	RefRange       string             `json:"refRange"`
	RepoRoot       string             `json:"repoRoot"`
	Head           string             `json:"head"`
	HunkLinks      []matcher.HunkLink `json:"hunkLinks"`
	FileStats      ExtractFileStats   `json:"fileStats"`
	HunkCoverage   float64            `json:"hunkCoverage"`
	EligibleHunks  int                `json:"eligibleHunks"`
	EligibleEvents int                `json:"eligibleEvents"`
	Tools          []string           `json:"tools"`
}

// ExtractFileStats summarizes excluded files for downstream consumers.
type ExtractFileStats struct {
	ExcludedFiles int `json:"excludedFiles"`
}

// StagedSession is one redacted agent session blob.
type StagedSession struct {
	SessionID string                 `json:"sessionID"`
	Tool      string                 `json:"tool"`
	Events    []capture.SessionEvent `json:"events"`
}

// Run executes discover → parse → match → redact → stage.
func Run(ctx context.Context, opts RunOptions) (Result, error) {
	if opts.RepoRoot == "" {
		return Result{}, fmt.Errorf("repo root required")
	}
	repoRoot, err := filepath.Abs(opts.RepoRoot)
	if err != nil {
		return Result{}, fmt.Errorf("repo path: %w", err)
	}
	home := opts.HomeDir
	if home == "" {
		home, err = os.UserHomeDir()
		if err != nil {
			return Result{}, fmt.Errorf("home dir: %w", err)
		}
	}
	head := opts.Head
	if head == "" {
		head = "HEAD"
	}
	base := opts.Base
	if base == "" {
		base = detectDefaultBase(repoRoot)
	}
	tools := normalizeToolsList(opts.Tools)
	timeSlack := opts.TimeSlack
	if timeSlack <= 0 {
		timeSlack = defaultTimeSlack
	}

	gitOpts := capturegit.HunkOptions{RepoRoot: repoRoot, Base: base, Head: head}
	commitCount, err := capturegit.CommitCount(gitOpts)
	if err != nil {
		return Result{}, fmt.Errorf("commit count: %w", err)
	}
	refRange := base + ".." + head
	if commitCount == 0 {
		fallbackBase := fallbackHistoryBase(repoRoot, head)
		gitOpts.Base = fallbackBase
		refRange = fallbackBase + ".." + head
		commitCount, err = capturegit.CommitCount(gitOpts)
		if err != nil || commitCount == 0 {
			return Result{}, fmt.Errorf("empty ref range %s", refRange)
		}
	}

	hunks, err := capturegit.CollectHunks(gitOpts)
	if err != nil {
		return Result{}, fmt.Errorf("collect hunks: %w", err)
	}
	firstTS, lastTS, err := capturegit.CommitRange(gitOpts)
	if err != nil {
		return Result{}, fmt.Errorf("commit range: %w", err)
	}

	excluder := exclude.NewMatcher(nil)
	eligibleHunks, excludedFiles, hunkRefs := filterHunks(hunks, excluder)
	since := time.UnixMilli(firstTS).Add(-timeSlack)
	until := time.UnixMilli(lastTS).Add(timeSlack)

	discovered, err := parsers.DiscoverSessions(parsers.DiscoverOptions{
		HomeDir:     home,
		RepoRoot:    repoRoot,
		Since:       since,
		Until:       until,
		ClaudeDir:   opts.ClaudeDir,
		CodexDir:    opts.CodexDir,
		CursorVSCDB: opts.CursorVSCDB,
		Tools:       tools,
	})
	if err != nil {
		return Result{}, fmt.Errorf("discover sessions: %w", err)
	}

	inventory := capture.NewInventoryCollector()
	claudeParser := &claude.Parser{Inventory: inventory}
	codexParser := &codex.Parser{Inventory: inventory}
	cursorParser := &cursorparser.Parser{
		Inventory: inventory,
		Since:     since,
		Until:     until,
	}
	activeParsers := buildToolParsers(tools, claudeParser, codexParser, cursorParser)

	events, err := parsers.ParseAll(discovered, repoRoot, activeParsers)
	if err != nil {
		return Result{}, fmt.Errorf("parse sessions: %w", err)
	}

	eligibleEvents := filterEvents(events, excluder, firstTS-timeSlack.Milliseconds(), lastTS+time.Hour.Milliseconds())
	matchResult := matcher.Match(eligibleEvents, hunkRefs, matcher.DefaultConfig())
	hunkLinks := matcher.BuildHunkLinks(eligibleHunks, eligibleEvents, matchResult)

	t1 := len(matchResult.Tier1HunkIndexes)
	t2 := len(matchResult.Tier2HunkIndexes)
	coverage := 0.0
	if len(eligibleHunks) > 0 {
		coverage = float64(t1+t2) / float64(len(eligibleHunks))
	}
	matchedAgents := len(matchResult.MatchedEvents)

	result := Result{
		RefRange:       refRange,
		EligibleHunks:  len(eligibleHunks),
		EligibleEvents: len(eligibleEvents),
		HunkCoverage:   coverage,
		Tier1Hunks:     t1,
		Tier2Hunks:     t2,
		HunkLinks:      hunkLinks,
		Tools:          tools,
	}

	redactedEvents := redactEvents(eligibleEvents)
	stagedExtract := StagedExtract{
		RefRange:       refRange,
		RepoRoot:       repoRoot,
		Head:           head,
		HunkLinks:      hunkLinks,
		FileStats:      ExtractFileStats{ExcludedFiles: excludedFiles},
		HunkCoverage:   coverage,
		EligibleHunks:  len(eligibleHunks),
		EligibleEvents: len(eligibleEvents),
		Tools:          tools,
	}
	sessions := groupSessions(redactedEvents)

	client := opts.Telemetry
	if client == nil {
		client = telemetry.NewFromEnv()
	}
	client.EmitCaptureCoverage(ctx, telemetry.CaptureCoverageProps{
		Repo:         repoRoot,
		RefRange:     refRange,
		HunkCoverage: coverage,
		Tier1:        t1,
		Tier2:        t2,
		Tools:        tools,
	})
	client.EmitMatchRate(ctx, telemetry.MatchRateProps{
		Repo:           repoRoot,
		RefRange:       refRange,
		AgentPrecision: agentPrecision(len(eligibleEvents), matchedAgents),
		EligibleHunks:  len(eligibleHunks),
		EligibleEvents: len(eligibleEvents),
		Tools:          tools,
	})

	stager := opts.DB
	if stager == nil {
		stager, err = storage.OpenCaptureStager(ctx)
		if err != nil {
			return result, fmt.Errorf("open capture stager: %w", err)
		}
	}
	extractID, err := stageExtract(ctx, stager, repoRoot, refRange, stagedExtract)
	if err != nil {
		return result, fmt.Errorf("stage extract: %w", err)
	}
	result.StagedExtractID = extractID
	for _, session := range sessions {
		if err := stageSession(ctx, stager, session); err != nil {
			return result, fmt.Errorf("stage session %s: %w", session.SessionID, err)
		}
		result.StagedSessions++
	}

	if !opts.StageOnly {
		if creds, ok := uploadauth.Load(); ok {
			headSHA, err := capturegit.RevParseCommit(repoRoot, head)
			if err != nil {
				return result, fmt.Errorf("resolve head commit: %w", err)
			}
			uploader := captureextract.NewClient(creds.APIURL, creds.Token)
			sessionPayloads := make([]captureextract.SessionPayload, 0, len(sessions))
			for _, s := range sessions {
				sessionPayloads = append(sessionPayloads, captureextract.SessionPayload{
					SessionID: s.SessionID,
					Tool:      s.Tool,
					Events:    s.Events,
				})
			}
			if err := uploader.UploadRun(
				ctx,
				repoRoot,
				refRange,
				headSHA,
				version.Current(),
				stagedExtract.HunkLinks,
				stagedExtract.FileStats,
				sessionPayloads,
				client,
			); err != nil {
				result.UploadError = fmt.Sprintf("upload capture: %v", err)
			}
		}
	}

	_ = report.VerdictForCoverage(coverage)
	return result, nil
}

func stageExtract(ctx context.Context, stager storage.CaptureStager, repoRoot, refRange string, extract StagedExtract) (string, error) {
	payload, err := json.Marshal(extract)
	if err != nil {
		return "", err
	}
	id := uuid.NewString()
	err = stager.StageExtract(ctx, storage.StagedExtract{
		ID:          id,
		RepoRoot:    repoRoot,
		RefRange:    refRange,
		PayloadJSON: payload,
	})
	return id, err
}

func stageSession(ctx context.Context, stager storage.CaptureStager, session StagedSession) error {
	payload, err := json.Marshal(session)
	if err != nil {
		return err
	}
	id := uuid.NewString()
	return stager.StageSession(ctx, storage.StagedSession{
		ID:          id,
		SessionID:   session.SessionID,
		Tool:        session.Tool,
		PayloadJSON: payload,
	})
}

func redactEvents(events []capture.SessionEvent) []capture.SessionEvent {
	out := make([]capture.SessionEvent, len(events))
	for i, ev := range events {
		ev.OldText = redact.Redact(ev.OldText)
		ev.NewText = redact.Redact(ev.NewText)
		ev.PromptContext = redact.Redact(ev.PromptContext)
		out[i] = ev
	}
	return out
}

func groupSessions(events []capture.SessionEvent) []StagedSession {
	byKey := map[string]*StagedSession{}
	order := []string{}
	for _, ev := range events {
		key := ev.Tool + ":" + ev.SessionID
		if ev.SessionID == "" {
			continue
		}
		s, ok := byKey[key]
		if !ok {
			s = &StagedSession{SessionID: ev.SessionID, Tool: ev.Tool}
			byKey[key] = s
			order = append(order, key)
		}
		s.Events = append(s.Events, ev)
	}
	out := make([]StagedSession, 0, len(order))
	for _, key := range order {
		out = append(out, *byKey[key])
	}
	return out
}

func agentPrecision(eligibleEvents, matched int) float64 {
	if eligibleEvents == 0 {
		return 0
	}
	return float64(matched) / float64(eligibleEvents)
}

func normalizeToolsList(tools []string) []string {
	if len(tools) == 0 {
		return []string{capture.ToolClaude, capture.ToolCodex, capture.ToolCursor}
	}
	out := make([]string, 0, len(tools))
	seen := map[string]struct{}{}
	for _, tool := range tools {
		tool = strings.TrimSpace(strings.ToLower(tool))
		if tool == "" {
			continue
		}
		if _, ok := seen[tool]; ok {
			continue
		}
		seen[tool] = struct{}{}
		out = append(out, tool)
	}
	return out
}

func buildToolParsers(tools []string, claudeParser, codexParser, cursorParser parsers.Parser) []parsers.Parser {
	var out []parsers.Parser
	for _, tool := range tools {
		switch tool {
		case capture.ToolClaude:
			out = append(out, claudeParser)
		case capture.ToolCodex:
			out = append(out, codexParser)
		case capture.ToolCursor:
			out = append(out, cursorParser)
		}
	}
	return out
}

func detectDefaultBase(repoRoot string) string {
	for _, candidate := range []string{"main", "master", "origin/main", "origin/master"} {
		if refExists(repoRoot, candidate) {
			return candidate
		}
	}
	return "main"
}

func fallbackHistoryBase(repoRoot, head string) string {
	count, err := capturegit.RevListCount(repoRoot, head)
	if err != nil || count <= 1 {
		return head + "~20"
	}
	back := 100
	if count-1 < back {
		back = count - 1
	}
	return fmt.Sprintf("%s~%d", head, back)
}

func refExists(repoRoot, ref string) bool {
	_, err := capturegit.CommitCount(capturegit.HunkOptions{RepoRoot: repoRoot, Base: ref, Head: "HEAD"})
	return err == nil
}

func filterHunks(hunks []capture.CommitHunk, ex *exclude.Matcher) ([]capture.CommitHunk, int, []matcher.HunkRef) {
	var eligible []capture.CommitHunk
	var refs []matcher.HunkRef
	excluded := 0
	seenExcluded := map[string]struct{}{}
	for _, h := range hunks {
		if ex.IsExcluded(h.FilePath) {
			if _, ok := seenExcluded[h.FilePath]; !ok {
				seenExcluded[h.FilePath] = struct{}{}
				excluded++
			}
			continue
		}
		if len(h.AddedLines) == 0 {
			continue
		}
		if strings.Contains(h.FilePath, "..") {
			continue
		}
		idx := len(eligible)
		eligible = append(eligible, h)
		refs = append(refs, matcher.HunkRef{
			Index:      idx,
			CommitSHA:  h.CommitSHA,
			FilePath:   filepath.ToSlash(h.FilePath),
			AddedLines: h.AddedLines,
			CommitTime: h.CommitTime,
		})
	}
	return eligible, excluded, refs
}

func filterEvents(events []capture.SessionEvent, ex *exclude.Matcher, sinceMS, untilMS int64) []capture.SessionEvent {
	var out []capture.SessionEvent
	for _, ev := range events {
		if !ev.IsEditEvent() || strings.TrimSpace(ev.NewText) == "" || ev.FilePath == "" {
			continue
		}
		if ex.IsExcluded(ev.FilePath) {
			continue
		}
		if ev.TS > 0 && (ev.TS < sinceMS || ev.TS > untilMS) {
			continue
		}
		ev.FilePath = filepath.ToSlash(ev.FilePath)
		out = append(out, ev)
	}
	return out
}
