package vcs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	"github.com/satoricorp/gx/internal/storage"
)

var ErrDetachedHEAD = errors.New("detached HEAD")

func isUnbornHEADError(err error) bool {
	if err == nil {
		return false
	}
	value := strings.ToLower(err.Error())
	return strings.Contains(value, "unknown revision") ||
		strings.Contains(value, "bad revision") ||
		strings.Contains(value, "ambiguous argument 'head'") ||
		strings.Contains(value, "your current branch does not have any commits yet")
}

func (s *Service) currentStagedCommitBranch(ctx context.Context, repoRoot string) (string, error) {
	branch, err := s.runTrimmed(ctx, repoRoot, "git", "branch", "--show-current")
	if err != nil {
		return "", err
	}
	branch = strings.TrimSpace(branch)
	if branch == "" {
		return "", fmt.Errorf("%w; check out a branch before committing", ErrDetachedHEAD)
	}
	return cleanRefName(branch), nil
}

// resolveStagedCommitStack maps a recorded commit onto a stack. The stack is
// the checked-out branch: GX records what git already did and never moves HEAD
// or mints branches of its own.
func (s *Service) resolveStagedCommitStack(ctx context.Context, repo RepoInfo, branch, headCommit string) (StackInfo, error) {
	baseRef := s.publicStackBaseRef(ctx, repo, s.defaultStackBaseRef(repo))
	onBase := branch == baseCheckoutRef(baseRef) || branch == baseRef
	bookmark := branch
	store, err := openStore(ctx)
	if err != nil {
		return StackInfo{}, err
	}
	defer store.Close()
	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return StackInfo{}, err
	}
	if existing, err := store.FindStackByBookmark(ctx, repoID, bookmark); err != nil {
		return StackInfo{}, err
	} else if existing != nil {
		return stackInfoFromStorage(*existing), nil
	}
	baseCommit := s.stackBaseCommitID(ctx, repo.RootPath, baseRef)
	if onBase && strings.TrimSpace(headCommit) != "" {
		baseCommit = strings.TrimSpace(headCommit)
	}
	return StackInfo{
		Name:         stackNameFromBookmark(bookmark),
		BookmarkName: bookmark,
		BaseRef:      baseRef,
		BaseCommitID: baseCommit,
		Status:       "draft",
	}, nil
}

func (s *Service) refExists(ctx context.Context, repoRoot, ref string) bool {
	_, err := s.runTrimmed(ctx, repoRoot, "git", "rev-parse", "--verify", "refs/heads/"+cleanRefName(ref))
	return err == nil
}

func (s *Service) matchStagedSessions(ctx context.Context, repoRoot, headCommit, treeOID string) ([]string, []storage.SessionContext, []storage.SessionEventAttribution) {
	if !s.hasKnownAttachableSessions(ctx, repoRoot) {
		return nil, nil, nil
	}
	// Diff the captured OIDs rather than the live index so matching is based on
	// the exact staged state that was recorded.
	diff, err := s.runTrimmed(ctx, repoRoot, "git", "-c", "core.quotePath=false", "diff", "--unified=0", "--no-ext-diff", headCommit, treeOID)
	if err != nil || strings.TrimSpace(diff) == "" {
		return nil, nil, nil
	}
	hunks := commitHunksFromGitDiff(diff)
	if len(hunks) == 0 {
		return nil, nil, nil
	}
	excluder := exclude.NewMatcher(nil)
	var eligible []capture.CommitHunk
	var refs []matcher.HunkRef
	for _, hunk := range hunks {
		if excluder.IsExcluded(hunk.FilePath) || len(hunk.AddedLines) == 0 {
			continue
		}
		idx := len(eligible)
		eligible = append(eligible, hunk)
		refs = append(refs, matcher.HunkRef{
			Index:      idx,
			CommitSHA:  hunk.CommitSHA,
			FilePath:   hunk.FilePath,
			AddedLines: hunk.AddedLines,
			CommitTime: hunk.CommitTime,
		})
	}
	if len(eligible) == 0 {
		return nil, nil, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, nil
	}
	now := time.Now()
	discovered, err := parsers.DiscoverSessions(parsers.DiscoverOptions{
		HomeDir:  home,
		RepoRoot: repoRoot,
		Since:    now.Add(-7 * 24 * time.Hour),
		Until:    now.Add(time.Hour),
	})
	if err != nil {
		return nil, nil, nil
	}
	inventory := capture.NewInventoryCollector()
	activeParsers := []parsers.Parser{
		&claude.Parser{Inventory: inventory},
		&codex.Parser{Inventory: inventory},
		&cursorparser.Parser{Inventory: inventory, Since: now.Add(-7 * 24 * time.Hour), Until: now.Add(time.Hour)},
	}
	events, err := parsers.ParseAll(discovered, repoRoot, activeParsers)
	if err != nil {
		return nil, nil, nil
	}
	events = filterStagedCaptureEvents(events, excluder)
	events = s.filterUnattributedSessionEvents(ctx, repoRoot, events)
	if len(events) == 0 {
		return nil, nil, nil
	}
	result := matcher.Match(events, refs, matcher.DefaultConfig())
	links := matcher.BuildHunkLinks(eligible, events, result)
	sessionIDs := sessionIDsFromMatchedLinks(links)
	if len(sessionIDs) == 0 {
		return nil, nil, nil
	}
	return sessionIDs, storageSessionContextsFromEvents(events, sessionIDs), sessionEventAttributionsFromOutcomes(events, result.Outcomes, "commit")
}

func (s *Service) hasKnownAttachableSessions(ctx context.Context, repoRoot string) bool {
	store, err := openStore(ctx)
	if err != nil {
		return false
	}
	defer store.Close()
	sessionIDs, err := store.FindAttachableSessionsForRepo(ctx, repoRoot, 1)
	return err == nil && len(sessionIDs) > 0
}

func commitHunksFromGitDiff(diff string) []capture.CommitHunk {
	now := time.Now().UnixMilli()
	var hunks []capture.CommitHunk
	file := ""
	var added []string
	flush := func() {
		if file == "" || len(added) == 0 {
			added = nil
			return
		}
		hunks = append(hunks, capture.CommitHunk{
			CommitSHA:  "staged",
			FilePath:   filepath.ToSlash(file),
			AddedLines: append([]string(nil), added...),
			CommitTime: now,
		})
		added = nil
	}
	for _, line := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(line, "diff --git "):
			flush()
			file = ""
		case strings.HasPrefix(line, "+++ "):
			file = parseGitDiffPlusPath(line)
		case strings.HasPrefix(line, "@@"):
			flush()
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			added = append(added, strings.TrimPrefix(line, "+"))
		}
	}
	flush()
	return hunks
}

func parseGitDiffPlusPath(line string) string {
	line = strings.TrimPrefix(line, "+++ ")
	if strings.HasPrefix(line, "b/") {
		return strings.TrimPrefix(line, "b/")
	}
	if len(line) >= 2 && line[0] == '"' && line[len(line)-1] == '"' {
		unquoted, err := strconv.Unquote(line)
		if err == nil {
			return strings.TrimPrefix(unquoted, "b/")
		}
	}
	return line
}

func filterStagedCaptureEvents(events []capture.SessionEvent, ex *exclude.Matcher) []capture.SessionEvent {
	out := make([]capture.SessionEvent, 0, len(events))
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

func (s *Service) filterUnattributedSessionEvents(ctx context.Context, repoRoot string, events []capture.SessionEvent) []capture.SessionEvent {
	store, err := openStore(ctx)
	if err != nil {
		return events
	}
	defer store.Close()
	repoInfo, err := s.ResolveGXRepoAtPath(ctx, repoRoot)
	var repo *storage.Repo
	if err == nil {
		repo, err = store.FindRepoByIdentity(ctx, repoInfo.GitCommonDir, repoInfo.RootPath)
	} else {
		repo, err = store.FindRepoByRoot(ctx, repoRoot)
	}
	if err != nil || repo == nil {
		return events
	}
	attributed, err := store.AttributedSessionEventKeys(ctx, repo.ID)
	if err != nil || len(attributed) == 0 {
		return events
	}
	out := make([]capture.SessionEvent, 0, len(events))
	for _, ev := range events {
		fingerprint := matcher.EventFingerprint(ev)
		if fingerprint == "" {
			out = append(out, ev)
			continue
		}
		if _, ok := attributed[storage.SessionEventAttributionKey(ev.Tool, ev.SessionID, fingerprint)]; ok {
			continue
		}
		out = append(out, ev)
	}
	return out
}

func sessionIDsFromMatchedLinks(links []matcher.HunkLink) []string {
	var ids []string
	for _, link := range links {
		if link.Authorship != matcher.AuthorshipAgent || link.SessionID == "" || link.Tier > matcher.TierFuzzy {
			continue
		}
		ids = append(ids, link.SessionID)
	}
	return uniqueStrings(ids)
}

func storageSessionContextsFromEvents(events []capture.SessionEvent, sessionIDs []string) []storage.SessionContext {
	wanted := map[string]struct{}{}
	for _, sessionID := range sessionIDs {
		wanted[sessionID] = struct{}{}
	}
	byID := map[string]*storage.SessionContext{}
	var order []string
	for _, ev := range events {
		if _, ok := wanted[ev.SessionID]; !ok {
			continue
		}
		context := byID[ev.SessionID]
		if context == nil {
			context = &storage.SessionContext{
				SessionID:  ev.SessionID,
				Tool:       ev.Tool,
				Format:     "gx_session_events_v1",
				CapturedAt: ev.TS,
			}
			if strings.TrimSpace(ev.Model) != "" {
				model := ev.Model
				context.Model = &model
			}
			byID[ev.SessionID] = context
			order = append(order, ev.SessionID)
		}
		if context.CapturedAt == 0 || (ev.TS != 0 && ev.TS < context.CapturedAt) {
			context.CapturedAt = ev.TS
		}
		redacted := redactStagedSessionEvent(ev)
		encoded, err := appendSessionContextEvent(context.ContentJSON, redacted)
		if err == nil {
			context.ContentJSON = encoded
		}
	}
	out := make([]storage.SessionContext, 0, len(order))
	for _, sessionID := range order {
		out = append(out, *byID[sessionID])
	}
	return out
}

func sessionEventAttributionsFromOutcomes(events []capture.SessionEvent, outcomes []matcher.MatchOutcome, via string) []storage.SessionEventAttribution {
	var out []storage.SessionEventAttribution
	seen := map[string]struct{}{}
	for _, outcome := range outcomes {
		if outcome.Tier != matcher.TierExact && outcome.Tier != matcher.TierFuzzy {
			continue
		}
		if outcome.EventIndex < 0 || outcome.EventIndex >= len(events) {
			continue
		}
		ev := events[outcome.EventIndex]
		fingerprint := matcher.EventFingerprint(ev)
		if fingerprint == "" || ev.SessionID == "" || ev.Tool == "" {
			continue
		}
		key := storage.SessionEventAttributionKey(ev.Tool, ev.SessionID, fingerprint)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		confidence := outcome.Score
		out = append(out, storage.SessionEventAttribution{
			Tool:             ev.Tool,
			SessionID:        ev.SessionID,
			EventFingerprint: fingerprint,
			AttributedVia:    via,
			Confidence:       &confidence,
		})
	}
	return out
}

func redactStagedSessionEvent(ev capture.SessionEvent) capture.SessionEvent {
	ev.OldText = redact.Redact(ev.OldText)
	ev.NewText = redact.Redact(ev.NewText)
	ev.PromptContext = redact.Redact(ev.PromptContext)
	ev.Raw = nil
	return ev
}

func appendSessionContextEvent(existing []byte, ev capture.SessionEvent) ([]byte, error) {
	var events []capture.SessionEvent
	if len(existing) > 0 {
		if err := json.Unmarshal(existing, &events); err != nil {
			return nil, err
		}
	}
	events = append(events, ev)
	return json.Marshal(events)
}

func stagedProvenanceStatus(sessionIDs []string) string {
	if len(sessionIDs) > 0 {
		return "matched"
	}
	return "absent"
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
