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
	"github.com/satoricorp/gx/internal/gxconfig"
	"github.com/satoricorp/gx/internal/storage"
)

const ExitCodeNoStagedChanges = 2

var (
	ErrNoStagedChanges = errors.New("no staged changes")
	ErrDetachedHEAD    = errors.New("detached HEAD")
)

// CodedError carries a process exit code for CLI callers.
type CodedError struct {
	Code int
	Err  error
}

func (e *CodedError) Error() string { return e.Err.Error() }
func (e *CodedError) Unwrap() error { return e.Err }

func codedError(code int, err error) error {
	if err == nil {
		return nil
	}
	return &CodedError{Code: code, Err: err}
}

// stagedCommitAfterImportHook is set by tests to inject failures after jj git import.
var stagedCommitAfterImportHook func() error

type StagedRevisionOptions struct {
	Message             string
	PreferredSessionIDs []string
	SessionContexts     []storage.SessionContext
}

func (s *Service) RecordStagedRevision(ctx context.Context, opts StagedRevisionOptions) (CommitResult, error) {
	if err := ValidateCommitMessage(opts.Message); err != nil {
		return CommitResult{}, err
	}
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return CommitResult{}, err
	}
	var result CommitResult
	err = withRepoLock(repo.RootPath, func() error {
		var commitErr error
		result, commitErr = s.recordStagedRevisionUnlocked(ctx, repo, opts)
		return commitErr
	})
	return result, err
}

func (s *Service) recordStagedRevisionUnlocked(ctx context.Context, repo RepoInfo, opts StagedRevisionOptions) (result CommitResult, err error) {
	repoRoot := repo.RootPath
	originalBranch, err := s.currentStagedCommitBranch(ctx, repoRoot)
	if err != nil {
		return CommitResult{}, err
	}
	if err := s.preflightStagedIndex(ctx, repoRoot); err != nil {
		return CommitResult{}, err
	}
	headCommit, err := s.runTrimmed(ctx, repoRoot, "git", "rev-parse", "HEAD")
	if err != nil {
		if isUnbornHEADError(err) {
			return CommitResult{}, fmt.Errorf("repository has no commits yet; make an initial commit before gx commit")
		}
		return CommitResult{}, fmt.Errorf("read HEAD: %w", err)
	}
	preWCChangeID, err := s.changeIDForRevIgnoreWorkingCopy(ctx, repoRoot, "@")
	if err != nil {
		return CommitResult{}, fmt.Errorf("read jj working copy: %w", err)
	}
	jjParent, err := s.commitIDForRevIgnoreWorkingCopy(ctx, repoRoot, "@-")
	if err != nil {
		return CommitResult{}, fmt.Errorf("read jj parent: %w", err)
	}
	if strings.TrimSpace(jjParent) != strings.TrimSpace(headCommit) {
		return CommitResult{}, fmt.Errorf("git HEAD and gx working-copy parent diverged; run gx sync or check out the matching branch before gx commit")
	}

	stagedFiles, err := s.stagedFileNames(ctx, repoRoot)
	if err != nil {
		return CommitResult{}, err
	}
	treeOID, err := s.runTrimmed(ctx, repoRoot, "git", "write-tree")
	if err != nil {
		return CommitResult{}, fmt.Errorf("write staged tree: %w", err)
	}
	stagedPatch, err := s.runTrimmed(ctx, repoRoot, "git", "diff", "--cached", "--binary")
	if err != nil {
		return CommitResult{}, fmt.Errorf("snapshot staged diff: %w", err)
	}

	opID, err := s.currentOperationIgnoreWorkingCopy(ctx, repoRoot)
	if err != nil {
		return CommitResult{}, fmt.Errorf("read jj operation: %w", err)
	}
	mutated := false
	tempRef := ""
	var stackBookmark string
	var createdBranch bool
	defer func() {
		if tempRef != "" {
			_, _ = s.runner.Run(ctx, repoRoot, "git", "update-ref", "-d", tempRef)
		}
		if err != nil && mutated {
			_ = s.RestoreOperation(ctx, repoRoot, opID)
			_, _ = s.runner.Run(ctx, repoRoot, "git", "symbolic-ref", "HEAD", "refs/heads/"+originalBranch)
			if createdBranch && stackBookmark != "" && stackBookmark != originalBranch {
				_, _ = s.runner.Run(ctx, repoRoot, "git", "update-ref", "-d", "refs/heads/"+stackBookmark)
			}
			_, _ = s.runner.Run(ctx, repoRoot, "git", "reset", "--mixed", headCommit)
			if treeOID != "" {
				_, _ = s.runner.Run(ctx, repoRoot, "git", "restore", "--source="+treeOID, "--staged", ":/")
			} else if strings.TrimSpace(stagedPatch) != "" {
				patchPath := filepath.Join(repoRoot, ".gx", "rollback-staged.patch")
				_ = os.MkdirAll(filepath.Dir(patchPath), 0o755)
				if writeErr := os.WriteFile(patchPath, []byte(stagedPatch), 0o600); writeErr == nil {
					_, _ = s.runner.Run(ctx, repoRoot, "git", "apply", "--cached", "--whitespace=nowarn", patchPath)
					_ = os.Remove(patchPath)
				}
			}
		}
	}()

	sessionIDs, sessionContexts, eventAttributions := s.matchStagedSessions(ctx, repoRoot)

	name, email, err := s.commitTreeIdentity(ctx, repoRoot)
	if err != nil {
		return CommitResult{}, err
	}
	commitID, err := s.runTrimmed(ctx, repoRoot, "git",
		"-c", "user.name="+name,
		"-c", "user.email="+email,
		"commit-tree", treeOID, "-p", headCommit, "-m", opts.Message)
	if err != nil {
		return CommitResult{}, fmt.Errorf("create staged commit: %w", err)
	}
	commitID = strings.TrimSpace(commitID)
	if commitID == "" {
		return CommitResult{}, fmt.Errorf("create staged commit: git returned an empty commit id")
	}

	stack, createdBranch, err := s.resolveStagedCommitStack(ctx, repo, originalBranch, opts.Message, headCommit)
	if err != nil {
		return CommitResult{}, err
	}
	stackBookmark = stack.BookmarkName
	if err := s.rejectProtectedStackBookmark(ctx, repoRoot, stack.BookmarkName); err != nil {
		return CommitResult{}, err
	}

	tempRef = "refs/heads/gx/import/" + shortID(commitID, 12) + "-" + strconv.FormatInt(time.Now().UnixNano(), 36)
	if _, err := s.runner.Run(ctx, repoRoot, "git", "update-ref", tempRef, commitID); err != nil {
		return CommitResult{}, fmt.Errorf("create temporary gx import ref: %w", err)
	}
	if _, err := s.runner.Run(ctx, repoRoot, "jj", "git", "import"); err != nil {
		return CommitResult{}, fmt.Errorf("import staged commit into jj: %w", err)
	}
	mutated = true
	if stagedCommitAfterImportHook != nil {
		if hookErr := stagedCommitAfterImportHook(); hookErr != nil {
			return CommitResult{}, hookErr
		}
	}
	if err := s.attachGitBranch(ctx, repoRoot, stack.BookmarkName, commitID, false); err != nil {
		return CommitResult{}, err
	}
	repo, err = s.ResolveJJRepoAtPath(ctx, repoRoot)
	if err != nil {
		return CommitResult{}, err
	}
	if err := s.reconcileRepoChanges(ctx, repo, stack.BookmarkName); err != nil {
		return CommitResult{}, err
	}
	if err := s.abandonOrphanWorkingCopyChange(ctx, repoRoot, preWCChangeID); err != nil {
		return CommitResult{}, err
	}
	change, err := s.CurrentChange(ctx, repoRoot, "@-")
	if err != nil {
		return CommitResult{}, err
	}
	if err := validateRecordedChangeDescription(change.Description); err != nil {
		return CommitResult{}, err
	}
	if len(change.Files) == 0 && len(stagedFiles) > 0 {
		change.Files = stagedFiles
	}
	opID, err = s.CurrentOperation(ctx, repoRoot)
	if err != nil {
		return CommitResult{}, err
	}
	result = CommitResult{
		Repo:                     repo,
		Change:                   change,
		Stack:                    &stack,
		OperationID:              opID,
		PreferredSessionIDs:      uniqueStrings(append(opts.PreferredSessionIDs, sessionIDs...)),
		SessionContexts:          append(opts.SessionContexts, sessionContexts...),
		CreatedBranch:            createdBranch,
		ProvenanceStatus:         stagedProvenanceStatus(sessionIDs),
		SkipRepoLocalSessions:    true,
		SessionEventAttributions: eventAttributions,
	}
	if err := recordCommit(ctx, result); err != nil {
		return result, fmt.Errorf("record revision metadata: %w", err)
	}
	if err := s.assertStagedCommitPostcondition(ctx, repoRoot); err != nil {
		return result, err
	}
	return result, nil
}

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
		return "", fmt.Errorf("%w; check out a branch before running gx commit", ErrDetachedHEAD)
	}
	return cleanRefName(branch), nil
}

func (s *Service) preflightStagedIndex(ctx context.Context, repoRoot string) error {
	if out, err := s.runTrimmed(ctx, repoRoot, "git", "ls-files", "-u"); err != nil {
		return fmt.Errorf("check unmerged paths: %w", err)
	} else if strings.TrimSpace(out) != "" {
		return fmt.Errorf("unmerged paths are staged; resolve conflicts before gx commit")
	}
	raw, err := s.runTrimmed(ctx, repoRoot, "git", "diff", "--cached", "--raw", "--no-abbrev")
	if err != nil {
		return fmt.Errorf("inspect staged diff: %w", err)
	}
	if intentToAdd, submodule := parseCachedRawDiff(raw); intentToAdd {
		return fmt.Errorf("intent-to-add entries are not supported; run git add with real content before gx commit")
	} else if submodule {
		return fmt.Errorf("staged submodule changes are not supported by gx commit")
	}
	if files, err := s.stagedFileNames(ctx, repoRoot); err != nil {
		return err
	} else if len(files) == 0 {
		return codedError(ExitCodeNoStagedChanges, fmt.Errorf("%w; run git add first", ErrNoStagedChanges))
	}
	return nil
}

func parseCachedRawDiff(raw string) (intentToAdd bool, submodule bool) {
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, ":") {
			continue
		}
		fields := strings.Fields(line[1:])
		if len(fields) < 4 {
			continue
		}
		oldMode := fields[0]
		newMode := fields[1]
		newSHA := fields[3]
		if oldMode == "160000" || newMode == "160000" {
			submodule = true
		}
		if newMode != "000000" && strings.HasPrefix(newSHA, "0000000") && newSHA == strings.Repeat("0", len(newSHA)) {
			intentToAdd = true
		}
	}
	return intentToAdd, submodule
}

func (s *Service) stagedFileNames(ctx context.Context, repoRoot string) ([]string, error) {
	out, err := s.runTrimmed(ctx, repoRoot, "git", "diff", "--cached", "--name-only")
	if err != nil {
		return nil, fmt.Errorf("list staged files: %w", err)
	}
	return splitLines(out), nil
}

func (s *Service) commitTreeIdentity(ctx context.Context, repoRoot string) (string, string, error) {
	cfg, _ := gxconfig.Load()
	name := firstNonEmpty(s.gitConfigValue(ctx, repoRoot, "user.name"), strings.TrimSpace(cfg.User.Name))
	email := firstNonEmpty(s.gitConfigValue(ctx, repoRoot, "user.email"), strings.TrimSpace(cfg.User.Email))
	if name == "" || email == "" {
		return "", "", fmt.Errorf("git identity is required for gx commit; run gx init or configure git user.name and user.email")
	}
	return name, email, nil
}

func (s *Service) resolveStagedCommitStack(ctx context.Context, repo RepoInfo, branch, message, headCommit string) (StackInfo, bool, error) {
	baseRef := s.publicStackBaseRef(ctx, repo, s.defaultStackBaseRef(repo))
	onBase := branch == baseCheckoutRef(baseRef) || branch == baseRef
	bookmark := branch
	created := false
	if onBase {
		var err error
		bookmark, err = s.uniqueStackBookmark(ctx, repo.RootPath, stackBookmarkName(message, ""))
		if err != nil {
			return StackInfo{}, false, err
		}
		created = true
	}
	if s.isProtectedRef(ctx, repo.RootPath, bookmark) {
		return StackInfo{}, false, fmt.Errorf("branch %s is protected; check out a feature branch or commit from the base branch to mint a GX stack", bookmark)
	}
	store, err := openStore(ctx)
	if err != nil {
		return StackInfo{}, false, err
	}
	defer store.Close()
	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return StackInfo{}, false, err
	}
	if existing, err := store.FindStackByBookmark(ctx, repoID, bookmark); err != nil {
		return StackInfo{}, false, err
	} else if existing != nil {
		return stackInfoFromStorage(*existing), created, nil
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
	}, created, nil
}

func (s *Service) uniqueStackBookmark(ctx context.Context, repoRoot, base string) (string, error) {
	base = cleanRefName(base)
	if base == "" {
		base = "feature/commit"
	}
	candidate := base
	for i := 2; i < 100; i++ {
		if !s.refExists(ctx, repoRoot, candidate) {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	return "", fmt.Errorf("could not create a unique stack branch for %s", base)
}

func (s *Service) refExists(ctx context.Context, repoRoot, ref string) bool {
	_, err := s.runTrimmed(ctx, repoRoot, "git", "rev-parse", "--verify", "refs/heads/"+cleanRefName(ref))
	return err == nil
}

func (s *Service) assertStagedCommitPostcondition(ctx context.Context, repoRoot string) error {
	if files, err := s.stagedFileNames(ctx, repoRoot); err != nil {
		return err
	} else if len(files) > 0 {
		return fmt.Errorf("gx commit recorded the revision but staged changes remain; run git status")
	}
	return nil
}

func (s *Service) commitIDForRevIgnoreWorkingCopy(ctx context.Context, repoRoot, rev string) (string, error) {
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "--ignore-working-copy", "-r", rev, "--no-graph", "-T", "commit_id")
	if err != nil {
		return "", err
	}
	return lastNonEmptyLine(out), nil
}

func (s *Service) changeIDForRevIgnoreWorkingCopy(ctx context.Context, repoRoot, rev string) (string, error) {
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "--ignore-working-copy", "-r", rev, "--no-graph", "-T", "change_id")
	if err != nil {
		return "", err
	}
	return lastNonEmptyLine(out), nil
}

func (s *Service) currentOperationIgnoreWorkingCopy(ctx context.Context, repoRoot string) (string, error) {
	return s.runStdoutTrimmed(ctx, repoRoot, "jj", "op", "log", "--ignore-working-copy", "-n", "1", "--no-graph", "-T", "id")
}

func (s *Service) abandonOrphanWorkingCopyChange(ctx context.Context, repoRoot, changeID string) error {
	changeID = strings.TrimSpace(changeID)
	if changeID == "" {
		return nil
	}
	desc, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "--ignore-working-copy", "-r", changeID, "--no-graph", "-T", "description")
	if err != nil || strings.TrimSpace(desc) != "" {
		return nil
	}
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl)
	if err != nil {
		return nil
	}
	for _, line := range splitLines(out) {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[1]) == changeID {
			return nil
		}
	}
	_, err = s.runner.Run(ctx, repoRoot, "jj", "abandon", changeID)
	return err
}

func (s *Service) matchStagedSessions(ctx context.Context, repoRoot string) ([]string, []storage.SessionContext, []storage.SessionEventAttribution) {
	if !s.hasKnownAttachableSessions(ctx, repoRoot) {
		return nil, nil, nil
	}
	diff, err := s.runTrimmed(ctx, repoRoot, "git", "-c", "core.quotePath=false", "diff", "--cached", "--unified=0", "--no-ext-diff")
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
	repo, err := store.FindRepoByRoot(ctx, repoRoot)
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
