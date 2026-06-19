package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/capture"
	capturegit "github.com/satoricorp/gx/internal/capture/git"
	"github.com/satoricorp/gx/internal/capture/exclude"
	"github.com/satoricorp/gx/internal/capture/matcher"
	"github.com/satoricorp/gx/internal/capture/parsers"
	"github.com/satoricorp/gx/internal/capture/parsers/claude"
	"github.com/satoricorp/gx/internal/capture/parsers/codex"
	cursorparser "github.com/satoricorp/gx/internal/capture/parsers/cursor"
	"github.com/satoricorp/gx/internal/capture/report"
)

type funnelStats struct {
	SessionsScanned      int
	EventsTotal          int
	EventsByKind         map[string]int
	EditEvents           int
	EditsWithNewTextPath int
	EligibleHunks        int
	MatchedTier12        int
}

func main() {
	os.Exit(run())
}

func run() int {
	var (
		repoFlag       string
		baseFlag       string
		headFlag       string
		homeFlag       string
		claudeDir      string
		codexDir       string
		cursorVSCDB    string
		toolsFlag      string
		outputPath     string
		inventoryOut   string
		timeSlack      time.Duration
		verboseFunnel  bool
	)
	flag.StringVar(&repoFlag, "repo", ".", "repository root")
	flag.StringVar(&baseFlag, "base", "", "base ref (auto-detect if empty)")
	flag.StringVar(&headFlag, "head", "HEAD", "head ref")
	flag.StringVar(&homeFlag, "home", "", "home directory for session discovery")
	flag.StringVar(&claudeDir, "claude-dir", "", "override Claude projects dir")
	flag.StringVar(&codexDir, "codex-dir", "", "override Codex sessions dir")
	flag.StringVar(&cursorVSCDB, "cursor-vscdb", "", "override Cursor state.vscdb path")
	flag.StringVar(&toolsFlag, "tools", "claude,codex,cursor", "comma-separated tools: claude,codex,cursor")
	flag.StringVar(&outputPath, "o", "", "write JSON report to path")
	flag.StringVar(&inventoryOut, "inventory-out", "docs/work-packages/WP-0-field-inventory.json", "field inventory output path")
	flag.DurationVar(&timeSlack, "time-slack", 24*time.Hour, "session time slack around commit range")
	flag.BoolVar(&verboseFunnel, "verbose-funnel", false, "print discovery funnel metrics")
	flag.Parse()

	repoRoot, err := filepath.Abs(repoFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "repo path: %v\n", err)
		return 1
	}
	home := homeFlag
	if home == "" {
		home, err = os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "home dir: %v\n", err)
			return 1
		}
	}
	tools := splitTools(toolsFlag)

	base := baseFlag
	if base == "" {
		base = detectDefaultBase(repoRoot)
	}

	inventory := capture.NewInventoryCollector()
	claudeParser := &claude.Parser{Inventory: inventory}
	codexParser := &codex.Parser{Inventory: inventory}
	enrichInventory(inventory)

	opts := capturegit.HunkOptions{RepoRoot: repoRoot, Base: base, Head: headFlag}
	commitCount, err := capturegit.CommitCount(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "commit count: %v\n", err)
		return 1
	}

	refRange := base + ".." + headFlag
	var firstTS, lastTS int64
	var hunks []capture.CommitHunk
	if commitCount == 0 {
		fallbackBase := fallbackHistoryBase(repoRoot, headFlag)
		fmt.Fprintf(os.Stderr, "note: %s empty, falling back to %s..%s\n", refRange, fallbackBase, headFlag)
		opts.Base = fallbackBase
		refRange = fallbackBase + ".." + headFlag
		commitCount, err = capturegit.CommitCount(opts)
		if err != nil || commitCount == 0 {
			fmt.Fprintf(os.Stderr, "empty ref range after fallback: %v\n", err)
			return 1
		}
		hunks, err = capturegit.CollectHunks(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "collect hunks: %v\n", err)
			return 1
		}
		firstTS, lastTS, err = capturegit.CommitRange(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "commit range: %v\n", err)
			return 1
		}
	} else {
		hunks, err = capturegit.CollectHunks(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "collect hunks: %v\n", err)
			return 1
		}
		firstTS, lastTS, err = capturegit.CommitRange(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "commit range: %v\n", err)
			return 1
		}
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
		ClaudeDir:   claudeDir,
		CodexDir:    codexDir,
		CursorVSCDB: cursorVSCDB,
		Tools:       tools,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "discover sessions: %v\n", err)
		return 1
	}

	cursorParser := &cursorparser.Parser{
		Inventory: inventory,
		Since:     since,
		Until:     until,
	}
	activeParsers := buildParsers(tools, claudeParser, codexParser, cursorParser)

	events, err := parsers.ParseAll(discovered, repoRoot, activeParsers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse sessions: %v\n", err)
		return 1
	}

	funnel := buildFunnel(discovered, events, eligibleHunks)
	eligibleEvents := filterEvents(events, excluder, firstTS-timeSlack.Milliseconds(), lastTS+time.Hour.Milliseconds())
	matchResult := matcher.Match(eligibleEvents, hunkRefs, matcher.DefaultConfig())
	funnel.EditsWithNewTextPath = len(eligibleEvents)
	funnel.MatchedTier12 = matchResult.AuthorshipHunkCount()

	matchedAgents := countMatchedAgents(matchResult)

	repoStats := report.BuildRepoStats(report.BuildInput{
		Root:                repoRoot,
		RefRange:            refRange,
		CommitCount:         commitCount,
		EligibleHunks:       len(eligibleHunks),
		EligibleAgentEvents: len(eligibleEvents),
		MatchedAgentEvents:  matchedAgents,
		MatchResult:         matchResult,
		ExcludedFiles:       excludedFiles,
	})
	matchReport := report.NewMatchReport([]report.RepoStats{repoStats})

	if verboseFunnel {
		printFunnel(os.Stderr, funnel, tools, matchReport.Aggregate.GateVerdict)
	}
	report.PrintTable(os.Stdout, matchReport)
	fmt.Fprintf(os.Stdout, "gate verdict: %s (hunk_coverage=%.1f%%, tier1+2=%d/%d)\n",
		matchReport.Aggregate.GateVerdict,
		matchReport.Aggregate.HunkCoverage*100,
		matchReport.Aggregate.Tier1Hunks+matchReport.Aggregate.Tier2Hunks,
		matchReport.Aggregate.EligibleHunks,
	)

	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "create report: %v\n", err)
			return 1
		}
		defer f.Close()
		if err := report.WriteJSON(f, matchReport); err != nil {
			fmt.Fprintf(os.Stderr, "write report: %v\n", err)
			return 1
		}
	}

	inv := inventory.Build()
	if err := writeInventory(inventoryOut, inv); err != nil {
		fmt.Fprintf(os.Stderr, "write inventory: %v\n", err)
		return 1
	}
	fmt.Fprintf(os.Stderr, "field inventory: %s\n", inventoryOut)
	return 0
}

func splitTools(raw string) []string {
	var tools []string
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(strings.ToLower(part))
		if part != "" {
			tools = append(tools, part)
		}
	}
	return tools
}

func buildParsers(tools []string, claudeParser, codexParser parsers.Parser, cursorParser parsers.Parser) []parsers.Parser {
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

func buildFunnel(discovered []parsers.DiscoveredSession, events []capture.SessionEvent, eligibleHunks []capture.CommitHunk) funnelStats {
	stats := funnelStats{
		SessionsScanned: parsers.CountSessionsScanned(discovered, events),
		EventsTotal:   len(events),
		EventsByKind:  map[string]int{},
		EligibleHunks: len(eligibleHunks),
	}
	for _, ev := range events {
		kind := ev.Kind
		if kind == "" {
			kind = "unknown"
		}
		stats.EventsByKind[kind]++
		if ev.IsEditEvent() {
			stats.EditEvents++
		}
	}
	return stats
}

func printFunnel(w *os.File, stats funnelStats, tools []string, verdict string) {
	fmt.Fprintf(w, "discovery funnel (tools=%s)\n", strings.Join(tools, ","))
	fmt.Fprintf(w, "  sessions_scanned → %d\n", stats.SessionsScanned)
	fmt.Fprintf(w, "  events_total → %d\n", stats.EventsTotal)
	fmt.Fprintf(w, "  events_by_kind → edit=%d message=%d tool_call=%d tool_result=%d other=%d\n",
		stats.EventsByKind[capture.KindEdit],
		stats.EventsByKind[capture.KindMessage],
		stats.EventsByKind[capture.KindToolCall],
		stats.EventsByKind[capture.KindToolResult],
		stats.EventsByKind["unknown"]+stats.EventsByKind[""],
	)
	fmt.Fprintf(w, "  edit_events → %d\n", stats.EditEvents)
	fmt.Fprintf(w, "  edits_with_newtext+path → %d\n", stats.EditsWithNewTextPath)
	fmt.Fprintf(w, "  |H| → %d\n", stats.EligibleHunks)
	fmt.Fprintf(w, "  |M1∪M2| → %d\n", stats.MatchedTier12)
	fmt.Fprintf(w, "  gate_verdict → %s\n", verdict)
}

func detectDefaultBase(repoRoot string) string {
	for _, candidate := range fallbackBases() {
		if refExists(repoRoot, candidate) {
			return candidate
		}
	}
	return "main"
}

func fallbackBases() []string {
	return []string{
		"main",
		"master",
		"gx/main",
		"origin/main",
		"origin/master",
	}
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

func enrichInventory(inv *capture.InventoryCollector) {
	claudePaths := []struct {
		path, typ, maps, notes string
	}{
		{"type", "string", "kind", "record type discriminator"},
		{"sessionId", "string", "session_id", "Claude session identifier"},
		{"timestamp", "string", "ts", "ISO8601 event time"},
		{"message.role", "string", "", "message role"},
		{"message.model", "string", "model", "model id"},
		{"message.content[].type", "string", "kind", "content block type"},
		{"message.content[].tool_use.name", "string", "kind", "edit tool identifier"},
		{"message.content[].tool_use.input.file_path", "string", "file_path", "target file"},
		{"message.content[].tool_use.input.content", "string", "new_text", "Write tool payload"},
		{"message.content[].tool_use.input.old_string", "string", "old_text", "Edit old text"},
		{"message.content[].tool_use.input.new_string", "string", "new_text", "Edit new text"},
	}
	for _, p := range claudePaths {
		inv.RecordPath(capture.ToolClaude, p.path, p.typ, p.maps, p.notes)
	}
	codexPaths := []struct {
		path, typ, maps, notes string
	}{
		{"type", "string", "kind", "rollout record type"},
		{"timestamp", "string", "ts", "ISO8601 event time"},
		{"payload.session_id", "string", "session_id", "Codex session id"},
		{"payload.model", "string", "model", "model id"},
		{"payload.cwd", "string", "", "working directory"},
		{"payload.name", "string", "kind", "tool name"},
		{"payload.arguments", "string", "new_text", "patch or write payload"},
	}
	for _, p := range codexPaths {
		inv.RecordPath(capture.ToolCodex, p.path, p.typ, p.maps, p.notes)
	}
	cursorPaths := []struct {
		path, typ, maps, notes string
	}{
		{"composerId", "string", "session_id", "Cursor composer id"},
		{"createdAt", "string", "ts", "bubble ISO8601 time"},
		{"timingInfo.clientRpcSendTime", "number", "ts", "bubble send time ms"},
		{"toolFormerData.name", "string", "kind", "tool name (search_replace, write, edit_file, apply_patch)"},
		{"toolFormerData.status", "string", "", "completed tool calls only"},
		{"toolFormerData.params", "string", "new_text", "JSON string of tool params"},
		{"toolFormerData.params.relativeWorkspacePath", "string", "file_path", "repo-relative file path"},
		{"toolFormerData.params.newString", "string", "new_text", "search_replace inserted text"},
		{"toolFormerData.params.oldString", "string", "old_text", "search_replace replaced text"},
		{"toolFormerData.params.contents", "string", "new_text", "edit_file full contents"},
		{"toolFormerData.params.streamingContent", "string", "new_text", "edit_file_v2 streamed contents"},
		{"toolFormerData.result.beforeContentId", "string", "old_text", "content-store before hash"},
		{"toolFormerData.result.afterContentId", "string", "new_text", "content-store after hash"},
		{"composer.content.*", "string", "new_text", "deduplicated file contents keyed by hash"},
		{"toolFormerData.rawArgs", "string", "new_text", "apply_patch / MultiEdit payload"},
		{"codeBlocks[].uri.path", "string", "file_path", "write tool file path"},
		{"codeBlocks[].content", "string", "new_text", "write tool file contents"},
	}
	for _, p := range cursorPaths {
		inv.RecordPath(capture.ToolCursor, p.path, p.typ, p.maps, p.notes)
	}
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

func countMatchedAgents(result matcher.Result) int {
	return len(result.MatchedEvents)
}

func writeInventory(path string, inv capture.FieldInventory) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := inv.MarshalJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
