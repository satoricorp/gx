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
	capturegit "github.com/satoricorp/gx/internal/capture/git"
	"github.com/satoricorp/gx/internal/capture/matcher"
	"github.com/satoricorp/gx/internal/capture/parsers"
	"github.com/satoricorp/gx/internal/capture/parsers/claude"
	"github.com/satoricorp/gx/internal/capture/parsers/codex"
	cursorparser "github.com/satoricorp/gx/internal/capture/parsers/cursor"
	"github.com/satoricorp/gx/internal/capture/redact"
	"github.com/satoricorp/gx/internal/capture/repobind"
	"github.com/satoricorp/gx/internal/capture/report"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/telemetry"
	"github.com/satoricorp/gx/internal/vcs"
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
//
// StagedExtractIDs and StagedSessionIDs are the primary keys of the rows this
// run actually wrote. They are how a caller marks exactly what it staged
// shareable, with no dependence on revision resolution returning a non-empty
// set for the pushed range. Both are populated as rows are written, so they
// still describe a run that failed partway.
type Result struct {
	RefRange         string             `json:"refRange"`
	EligibleHunks    int                `json:"eligibleHunks"`
	EligibleEvents   int                `json:"eligibleEvents"`
	HunkCoverage     float64            `json:"hunkCoverage"`
	Tier1Hunks       int                `json:"tier1Hunks"`
	Tier2Hunks       int                `json:"tier2Hunks"`
	HunkLinks        []matcher.HunkLink `json:"hunkLinks"`
	StagedExtractID  string             `json:"stagedExtractID,omitempty"`
	StagedExtractIDs []string           `json:"stagedExtractIDs,omitempty"`
	StagedSessionIDs []string           `json:"stagedSessionIDs,omitempty"`
	StagedSessions   int                `json:"stagedSessions"`
	// AttributionSuspect marks a run that found agent sessions and committed
	// work but linked none of it. That combination is not evidence the change
	// was written by hand — it is the signature of session data that was not
	// on disk yet when capture read it, and it used to print as a clean
	// `sessions=0 coverage=0.0%`, which reads like an answer rather than a
	// gap.
	AttributionSuspect bool     `json:"attributionSuspect,omitempty"`
	Tools              []string `json:"tools"`
	DiscoveryProblems  []string `json:"discoveryProblems,omitempty"`
	UploadError        string   `json:"uploadError,omitempty"`
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

	discovery, err := parsers.Discover(parsers.DiscoverOptions{
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
	discovered := discovery.Sessions

	revisionIDs, err := vcs.RevisionIDsInGitRange(ctx, repoRoot, refRange)
	if err != nil {
		return Result{}, fmt.Errorf("resolve GX revision IDs: %w", err)
	}

	stager := opts.DB
	if stager == nil {
		stager, err = storage.OpenCaptureStager(ctx)
		if err != nil {
			return Result{}, fmt.Errorf("open capture stager: %w", err)
		}
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

	perSource, err := parsers.ParseAllPerSource(discovered, repoRoot, activeParsers)
	if err != nil {
		return Result{}, fmt.Errorf("parse sessions: %w", err)
	}
	var events []capture.SessionEvent
	for _, source := range perSource {
		events = append(events, source.Events...)
	}

	eligibleEvents := filterEvents(events, excluder, firstTS-timeSlack.Milliseconds(), lastTS+time.Hour.Milliseconds())
	matchResult := matcher.Match(eligibleEvents, hunkRefs, matcher.DefaultConfig())
	hunkLinks := matcher.BuildHunkLinks(eligibleHunks, eligibleEvents, matchResult)

	t1 := len(matchResult.Tier1HunkIndexes)
	t2 := len(matchResult.Tier2HunkIndexes)
	coverage := 0.0
	if len(eligibleHunks) > 0 {
		coverage = float64(matchResult.AuthorshipHunkCount()) / float64(len(eligibleHunks))
	}
	matchedAgents := len(matchResult.MatchedEvents)

	result := Result{
		RefRange:           refRange,
		EligibleHunks:      len(eligibleHunks),
		EligibleEvents:     len(eligibleEvents),
		HunkCoverage:       coverage,
		Tier1Hunks:         t1,
		Tier2Hunks:         t2,
		HunkLinks:          hunkLinks,
		Tools:              tools,
		DiscoveryProblems:  discovery.ProblemStrings(),
		AttributionSuspect: attributionSuspect(len(eligibleHunks), len(discovered), len(hunkLinks)),
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

	// Sessions are staged only after matching, and only the sources that
	// produced at least one content-matched event. Discovery deliberately
	// over-collects (it scans every project directory), so staging everything
	// it finds wrote hundreds of megabytes about other repos.
	matchedSources := map[int][]capture.SessionEvent{}
	for idx := range matchResult.RelevantEvents {
		src := eligibleEvents[idx].SourceIndex
		if _, ok := matchedSources[src]; !ok {
			matchedSources[src] = nil
		}
	}
	for _, ev := range redactedEvents {
		if _, ok := matchedSources[ev.SourceIndex]; ok {
			matchedSources[ev.SourceIndex] = append(matchedSources[ev.SourceIndex], ev)
		}
	}

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
	emitSchemaDrift(ctx, client, inventory)

	extractIDs, err := stageExtract(ctx, stager, repoRoot, refRange, stagedExtract, revisionIDs)
	result.StagedExtractIDs = extractIDs
	if len(extractIDs) > 0 {
		result.StagedExtractID = extractIDs[0]
	}
	if err != nil {
		return result, fmt.Errorf("stage extract: %w", err)
	}
	// A session that worked on this repository is context worth keeping even
	// when none of its text reached a commit in this range. Matching answers
	// what a session produced; it cannot answer for work that was explored and
	// abandoned, or for a transcript the agent had not finished writing when
	// capture read it — the case that reports zero coverage and looks exactly
	// like a change nobody used an agent for.
	//
	// Binding is what makes that safe: an unmatched session is staged only when
	// it is bound to this repository, and marking it shareable is gated on the
	// connected set separately.
	repoBinding := repobind.NewResolver().ForDirectory(ctx, repoRoot)
	for i, source := range perSource {
		matched, ok := matchedSources[i]
		if !ok {
			if !boundToRepo(ctx, perSource[i].Events, repoBinding) {
				continue
			}
			// Nothing matched, so there is no matched subset to stage: the
			// session's own events are what make it useful as context.
			matched = perSource[i].Events
		}
		sessionIDs, err := stageMatchedSource(ctx, stager, source.Source, matched, revisionIDs)
		result.StagedSessionIDs = append(result.StagedSessionIDs, sessionIDs...)
		result.StagedSessions += len(sessionIDs)
		if err != nil {
			return result, fmt.Errorf("stage session %s: %w", source.Source.SessionID, err)
		}
	}

	_ = report.VerdictForCoverage(coverage)
	return result, nil
}

// stageExtract writes the extract rows and returns their ids in write order,
// including the ids written before any failure: rows that reached SQLite are
// uploadable regardless of what happened after them.
func stageExtract(ctx context.Context, stager storage.CaptureStager, repoRoot, refRange string, extract StagedExtract, revisionIDs []string) ([]string, error) {
	payload, err := json.Marshal(extract)
	if err != nil {
		return nil, err
	}
	ids := revisionIDs
	if len(ids) == 0 {
		ids = []string{""}
	}
	var staged []string
	for _, revisionID := range ids {
		contentHash := storage.PayloadContentHash(payload)
		row := storage.StagedExtract{
			RepoRoot:    repoRoot,
			RefRange:    refRange,
			PayloadJSON: payload,
			RevisionID:  revisionID,
			ContentHash: contentHash,
		}
		row.ID = storage.CaptureRowID(revisionID, contentHash)
		if row.ID == "" {
			row.ID = uuid.NewString()
		}
		if err := stager.StageExtract(ctx, row); err != nil {
			return staged, err
		}
		staged = append(staged, row.ID)
	}
	return staged, nil
}

func emitSchemaDrift(ctx context.Context, client telemetry.Client, inventory *capture.InventoryCollector) {
	if inventory == nil {
		return
	}
	for tool, paths := range inventory.UnknownPaths() {
		if len(paths) == 0 {
			continue
		}
		client.EmitSchemaDrift(ctx, telemetry.SchemaDriftProps{
			Tool:  tool,
			Paths: paths,
			Count: len(paths),
		})
	}
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
		// A command event carries no file path and no NewText — the text it
		// wrote is in the command itself — so it cannot meet the edit-event
		// conditions and used to be dropped here, taking every sed, heredoc and
		// inline script edit out of attribution before matching ever saw them.
		// The exclusion matcher has nothing to test it against either, which is
		// harmless: it can only claim eligible hunks, and those are already
		// filtered.
		switch {
		case ev.IsCommandEvent():
		case ev.IsEditEvent() && strings.TrimSpace(ev.NewText) != "" && ev.FilePath != "":
			if ex.IsExcluded(ev.FilePath) {
				continue
			}
		default:
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

// attributionSuspect reports a run that had everything it needed and still
// linked nothing.
//
// Capture parses the agent's transcript while the agent is still writing it.
// A push that happens mid-session can read a file that is missing the very
// work being pushed, and the run then reports zero coverage — indistinguishable
// from a change genuinely written by hand. Truncating a real transcript
// reproduces it exactly: at ~70% of its eventual size the same ref range
// scores 0.0%, at 90% it scores 50.0%.
//
// Committed work plus discovered sessions plus no links at all is that shape,
// so callers can say so instead of reporting an absence of agent work.
func attributionSuspect(eligibleHunks, discoveredSessions, hunkLinks int) bool {
	return eligibleHunks > 0 && discoveredSessions > 0 && hunkLinks == 0
}

// boundToRepo reports whether a session was working in this repository.
//
// It is deliberately an identity check, not a relevance one: a session earns a
// place by where it ran, which is knowable for every session, rather than by
// whether its text survived into a commit, which is not.
func boundToRepo(ctx context.Context, events []capture.SessionEvent, repo repobind.Binding) bool {
	if !repo.Bound() {
		return false
	}
	bindable := make([]repobind.Event, 0, len(events))
	for _, e := range events {
		bindable = append(bindable, e)
	}
	return repobind.BindEvents(ctx, nil, bindable).Origin == repo.Origin
}
