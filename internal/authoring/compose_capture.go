package authoring

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/exclude"
	"github.com/satoricorp/gx/internal/capture/matcher"
	"github.com/satoricorp/gx/internal/capture/parsers"
	"github.com/satoricorp/gx/internal/capture/parsers/claude"
	"github.com/satoricorp/gx/internal/capture/parsers/codex"
	cursorparser "github.com/satoricorp/gx/internal/capture/parsers/cursor"
	"github.com/satoricorp/gx/internal/capture/redact"
	"github.com/satoricorp/gx/internal/provenance"
)

const workingCopyCommitSHA = "working"

// ComposeCaptureResult holds matcher output for a working-copy compose run.
type ComposeCaptureResult struct {
	HunkLinks    []matcher.HunkLink `json:"hunk_links"`
	HunkCoverage float64            `json:"hunk_coverage"`
	Tools        []string           `json:"tools"`
	SessionIDs   []string           `json:"session_ids,omitempty"`
	Contexts     []SessionContext   `json:"session_contexts,omitempty"`
	Status       string             `json:"status"`
}

// MatchWorkingCopy attaches capture matcher evidence to demux hunks.
func MatchWorkingCopy(ctx context.Context, repoRoot string, hunks []HunkRange, tools []string) (ComposeCaptureResult, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return ComposeCaptureResult{}, fmt.Errorf("repo path: %w", err)
	}
	tools = normalizeCaptureTools(tools)
	if len(hunks) == 0 {
		return ComposeCaptureResult{Tools: tools, Status: provenance.StatusAbsent}, nil
	}

	commitHunks, hunkRefs := hunkRangesToCommitHunks(hunks)
	excluder := exclude.NewMatcher(nil)
	eligibleHunks, _, refs := filterEligibleHunks(commitHunks, hunkRefs, excluder)
	if len(eligibleHunks) == 0 {
		return ComposeCaptureResult{Tools: tools, Status: provenance.StatusAbsent}, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ComposeCaptureResult{}, err
	}
	now := time.Now()
	since := now.Add(-7 * 24 * time.Hour)
	until := now.Add(time.Hour)

	discovered, err := parsers.DiscoverSessions(parsers.DiscoverOptions{
		HomeDir:  home,
		RepoRoot: repoRoot,
		Since:    since,
		Until:    until,
		Tools:    tools,
	})
	if err != nil {
		return ComposeCaptureResult{}, fmt.Errorf("discover sessions: %w", err)
	}

	inventory := capture.NewInventoryCollector()
	claudeParser := &claude.Parser{Inventory: inventory}
	codexParser := &codex.Parser{Inventory: inventory}
	cursorParser := &cursorparser.Parser{Inventory: inventory, Since: since, Until: until}
	activeParsers := buildCaptureParsers(tools, claudeParser, codexParser, cursorParser)

	events, err := parsers.ParseAll(discovered, repoRoot, activeParsers)
	if err != nil {
		return ComposeCaptureResult{}, fmt.Errorf("parse sessions: %w", err)
	}
	eligibleEvents := filterCaptureEvents(events, excluder)

	matchResult := matcher.Match(eligibleEvents, refs, matcher.DefaultConfig())
	links := matcher.BuildHunkLinks(eligibleHunks, eligibleEvents, matchResult)
	links = remapHunkLinkIDs(links, hunks, eligibleHunks)

	t1 := len(matchResult.Tier1HunkIndexes)
	t2 := len(matchResult.Tier2HunkIndexes)
	coverage := 0.0
	if len(eligibleHunks) > 0 {
		coverage = float64(t1+t2) / float64(len(eligibleHunks))
	}

	sessionIDs := sessionIDsFromHunkLinks(links)
	contexts := sessionContextsFromEvents(eligibleEvents, sessionIDs)
	status := provenance.StatusAbsent
	if len(sessionIDs) > 0 {
		status = provenance.StatusMatched
	}

	return ComposeCaptureResult{
		HunkLinks:    links,
		HunkCoverage: coverage,
		Tools:        tools,
		SessionIDs:   sessionIDs,
		Contexts:     contexts,
		Status:       status,
	}, nil
}

func resolveComposeProvenance(ctx context.Context, repoRoot string, capture ComposeCaptureResult) ([]string, []SessionContext, string, error) {
	if len(capture.SessionIDs) > 0 {
		return capture.SessionIDs, capture.Contexts, capture.Status, nil
	}
	store, err := openStore(ctx)
	if err != nil {
		return nil, nil, "", err
	}
	defer store.Close()
	attachment, err := provenance.Resolve(ctx, store, repoRoot)
	if err != nil {
		return nil, nil, "", err
	}
	return attachment.SessionIDs, nil, attachment.Status, nil
}

func hunkLinksForRevision(links []matcher.HunkLink, files []string) []matcher.HunkLink {
	if len(links) == 0 || len(files) == 0 {
		return nil
	}
	fileSet := map[string]struct{}{}
	for _, file := range files {
		fileSet[file] = struct{}{}
	}
	var out []matcher.HunkLink
	for _, link := range links {
		parts := strings.SplitN(link.HunkID, ":", 3)
		if len(parts) < 2 {
			continue
		}
		if _, ok := fileSet[parts[1]]; ok {
			out = append(out, link)
		}
	}
	return out
}

func sessionContextsForRevision(contexts []SessionContext, links []matcher.HunkLink) []SessionContext {
	if len(contexts) == 0 || len(links) == 0 {
		return nil
	}
	sessionIDs := map[string]struct{}{}
	for _, link := range links {
		if link.SessionID != "" && link.Authorship == matcher.AuthorshipAgent && link.Tier <= matcher.TierFuzzy {
			sessionIDs[link.SessionID] = struct{}{}
		}
	}
	if len(sessionIDs) == 0 {
		return nil
	}
	out := make([]SessionContext, 0, len(contexts))
	for _, context := range contexts {
		if _, ok := sessionIDs[context.SessionID]; ok {
			out = append(out, context)
		}
	}
	return out
}

func sessionIDsForRevisionLinks(fallback []string, links []matcher.HunkLink, status string) []string {
	ids := sessionIDsFromHunkLinks(links)
	if len(ids) > 0 {
		return ids
	}
	if status == provenance.StatusMatched {
		return nil
	}
	return append([]string(nil), fallback...)
}

func provenanceStatusForRevision(status string, sessionIDs []string) string {
	if status == provenance.StatusMatched && len(sessionIDs) == 0 {
		return provenance.StatusAbsent
	}
	return status
}

func sessionIDsFromHunkLinks(links []matcher.HunkLink) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, link := range links {
		if link.SessionID == "" || link.Tier > matcher.TierFuzzy {
			continue
		}
		if link.Authorship != matcher.AuthorshipAgent {
			continue
		}
		if _, ok := seen[link.SessionID]; ok {
			continue
		}
		seen[link.SessionID] = struct{}{}
		out = append(out, link.SessionID)
	}
	return out
}

func sessionContextsFromEvents(events []capture.SessionEvent, sessionIDs []string) []SessionContext {
	if len(sessionIDs) == 0 {
		return nil
	}
	wanted := map[string]struct{}{}
	for _, sessionID := range sessionIDs {
		if strings.TrimSpace(sessionID) != "" {
			wanted[sessionID] = struct{}{}
		}
	}
	byID := map[string]*SessionContext{}
	order := make([]string, 0, len(wanted))
	for _, ev := range events {
		if _, ok := wanted[ev.SessionID]; !ok {
			continue
		}
		context := byID[ev.SessionID]
		if context == nil {
			context = &SessionContext{
				SessionID:  ev.SessionID,
				Tool:       ev.Tool,
				Model:      ev.Model,
				Format:     "gx_session_events_v1",
				CapturedAt: ev.TS,
			}
			byID[ev.SessionID] = context
			order = append(order, ev.SessionID)
		}
		if context.Model == "" {
			context.Model = ev.Model
		}
		if context.CapturedAt == 0 || (ev.TS != 0 && ev.TS < context.CapturedAt) {
			context.CapturedAt = ev.TS
		}
		context.ContentRedacted = append(context.ContentRedacted, redactSessionEvent(ev))
	}
	out := make([]SessionContext, 0, len(order))
	for _, sessionID := range order {
		out = append(out, *byID[sessionID])
	}
	return out
}

func redactSessionEvent(ev capture.SessionEvent) capture.SessionEvent {
	ev.OldText = redact.Redact(ev.OldText)
	ev.NewText = redact.Redact(ev.NewText)
	ev.PromptContext = redact.Redact(ev.PromptContext)
	ev.Raw = nil
	return ev
}

func hunkRangesToCommitHunks(hunks []HunkRange) ([]capture.CommitHunk, []matcher.HunkRef) {
	now := time.Now().UnixMilli()
	var commitHunks []capture.CommitHunk
	var refs []matcher.HunkRef
	for _, hunk := range hunks {
		added := addedLinesFromPatch(hunk.Patch)
		if len(added) == 0 {
			continue
		}
		idx := len(commitHunks)
		commitHunks = append(commitHunks, capture.CommitHunk{
			CommitSHA:  workingCopyCommitSHA,
			FilePath:   hunk.File,
			AddedLines: added,
			CommitTime: now,
		})
		refs = append(refs, matcher.HunkRef{
			Index:      idx,
			CommitSHA:  workingCopyCommitSHA,
			FilePath:   hunk.File,
			AddedLines: added,
			CommitTime: now,
		})
	}
	return commitHunks, refs
}

func addedLinesFromPatch(patch string) []string {
	if strings.TrimSpace(patch) == "" {
		return nil
	}
	var added []string
	for _, line := range strings.Split(patch, "\n") {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			added = append(added, strings.TrimPrefix(line, "+"))
		}
	}
	return added
}

func remapHunkLinkIDs(links []matcher.HunkLink, demuxHunks []HunkRange, eligible []capture.CommitHunk) []matcher.HunkLink {
	if len(links) != len(eligible) {
		return links
	}
	fileToDemuxID := map[string]string{}
	for _, hunk := range demuxHunks {
		fileToDemuxID[hunk.File] = hunk.ID
	}
	out := make([]matcher.HunkLink, len(links))
	for i, link := range links {
		out[i] = link
		if i < len(eligible) {
			if id, ok := fileToDemuxID[eligible[i].FilePath]; ok {
				out[i].HunkID = id
			}
		}
	}
	return out
}

func filterEligibleHunks(hunks []capture.CommitHunk, refs []matcher.HunkRef, ex *exclude.Matcher) ([]capture.CommitHunk, int, []matcher.HunkRef) {
	var eligible []capture.CommitHunk
	var outRefs []matcher.HunkRef
	excluded := 0
	for i, h := range hunks {
		if ex.IsExcluded(h.FilePath) {
			excluded++
			continue
		}
		if len(h.AddedLines) == 0 {
			continue
		}
		idx := len(eligible)
		eligible = append(eligible, h)
		if i < len(refs) {
			ref := refs[i]
			ref.Index = idx
			outRefs = append(outRefs, ref)
		}
	}
	return eligible, excluded, outRefs
}

func filterCaptureEvents(events []capture.SessionEvent, ex *exclude.Matcher) []capture.SessionEvent {
	var out []capture.SessionEvent
	for _, ev := range events {
		if !ev.IsEditEvent() || strings.TrimSpace(ev.NewText) == "" || ev.FilePath == "" {
			continue
		}
		if ex.IsExcluded(ev.FilePath) {
			continue
		}
		ev.FilePath = filepath.ToSlash(ev.FilePath)
		out = append(out, ev)
	}
	return out
}

func normalizeCaptureTools(tools []string) []string {
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

func buildCaptureParsers(tools []string, claudeParser, codexParser, cursorParser parsers.Parser) []parsers.Parser {
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

type composeExportPayload struct {
	ProposalID   string             `json:"proposal_id"`
	RepoRoot     string             `json:"repo_root"`
	HunkLinks    []matcher.HunkLink `json:"hunk_links"`
	HunkCoverage float64            `json:"hunk_coverage"`
	Tools        []string           `json:"tools"`
	Revisions    int                `json:"revision_count"`
	Warnings     int                `json:"warnings"`
	ExportedAt   int64              `json:"exported_at"`
}

func exportComposeProposal(proposal DemuxProposal) error {
	payload := composeExportPayload{
		ProposalID:   proposal.ID,
		RepoRoot:     proposal.RepoRoot,
		HunkLinks:    proposal.HunkLinks,
		HunkCoverage: proposal.HunkCoverage,
		Tools:        proposal.CaptureTools,
		Revisions:    len(proposal.Revisions),
		Warnings:     len(proposal.FeasibilityWarnings),
		ExportedAt:   time.Now().Unix(),
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	if url := strings.TrimSpace(os.Getenv("GX_COMPOSE_EXPORT_URL")); url != "" {
		_ = url // spike: local path only for V1
	}
	dir, err := composeExportDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	name := proposal.ID
	if strings.TrimSpace(name) == "" {
		name = fmt.Sprintf("compose-%d", time.Now().Unix())
	}
	return os.WriteFile(filepath.Join(dir, name+".json"), data, 0o644)
}

func composeExportDir() (string, error) {
	if home := strings.TrimSpace(os.Getenv("GX_HOME")); home != "" {
		return filepath.Join(home, "compose-exports"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gx", "compose-exports"), nil
}

func readLastComposeHunkCoverage(repoRoot string) (float64, bool) {
	dir, err := composeExportDir()
	if err != nil {
		return 0, false
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, false
	}
	var latestPath string
	var latestMod time.Time
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(latestMod) {
			latestMod = info.ModTime()
			latestPath = filepath.Join(dir, entry.Name())
		}
	}
	if latestPath == "" {
		return 0, false
	}
	data, err := os.ReadFile(latestPath)
	if err != nil {
		return 0, false
	}
	var payload composeExportPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return 0, false
	}
	if repoRoot != "" && payload.RepoRoot != "" {
		absRepo, _ := filepath.Abs(repoRoot)
		absPayload, _ := filepath.Abs(payload.RepoRoot)
		if absRepo != "" && absPayload != "" && absRepo != absPayload {
			return 0, false
		}
	}
	return payload.HunkCoverage, true
}
