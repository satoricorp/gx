package vcs

import (
	"context"
	"strings"
	"time"

	"github.com/satoricorp/totality/internal/capture/matcher"
	"github.com/satoricorp/totality/internal/storage"
	"github.com/satoricorp/totality/internal/version"
)

// maxSessionsPerChange caps how many sessions one change may link. A single
// Claude conversation can produce dozens of subagent transcripts — each one is
// legitimately its own session row — and the console slices the bundle's
// session list, so an unbounded push would bury the parent conversation behind
// its own subagents.
const maxSessionsPerChange = 25

// AttachSessionsFromHunkLinks writes the sessions and change_sessions rows for
// one push, from the hunk links the capture pipeline already produced.
//
// This is the only writer of `change_sessions` on the git path, and it runs at
// push time rather than commit time on purpose. The rule is that a `sessions`
// row exists because a transcript was observed, never because a commit happened
// to match one: the commit-time matcher this replaced was gated on the sessions
// table already being non-empty while being the only thing that wrote it, so an
// empty table stayed empty forever.
//
// It never fails a push. Callers record the error and carry on.
func (s *Service) AttachSessionsFromHunkLinks(ctx context.Context, repo RepoInfo, links []matcher.HunkLink) (int, error) {
	bySHA, order := agentSessionsByCommit(links)
	if len(order) == 0 {
		return 0, nil
	}
	attached := 0
	err := withBusyRetry(ctx, "attach push sessions", func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()
		repoID, err := upsertRepo(ctx, store, repo)
		if err != nil {
			return err
		}
		attached = 0
		now := time.Now().UnixMilli()
		repoRoot := repo.RootPath
		for _, sha := range order {
			change, err := store.FindChangeByCommitID(ctx, repoID, sha)
			if err != nil {
				return err
			}
			if change == nil {
				// The commit has no recorded revision (no Totality trailer, or the
				// recovery pass could not resolve it). Nothing to link to.
				continue
			}
			candidates := bySHA[sha]
			if len(candidates) > maxSessionsPerChange {
				candidates = candidates[:maxSessionsPerChange]
			}
			sessionIDs := make([]string, 0, len(candidates))
			for _, candidate := range candidates {
				// The sessions row has to exist before change_sessions
				// references it: change_sessions.session_id is a foreign key,
				// and WriteChangeSessions also derives change_session_provenance
				// by reading the sessions row back.
				//
				// UpsertObservedSession, not UpsertSession: a transcript
				// observation only knows the tool name, and these IDs can
				// collide with richer rows written by an ingest pass (Cursor's
				// state.vscdb sessions keep the same "cursor-<composerID>" id
				// the ingest writer mints). Filling blanks is safe; overwriting
				// is not.
				if err := store.UpsertObservedSession(ctx, storage.Session{
					ID:        candidate.sessionID,
					CreatedAt: now,
					Command:   candidate.tool,
					Cwd:       repoRoot,
					TLVersion: version.Current(),
					Source:    &candidate.tool,
					RepoRoot:  &repoRoot,
				}); err != nil {
					return err
				}
				sessionIDs = append(sessionIDs, candidate.sessionID)
			}
			if len(sessionIDs) == 0 {
				continue
			}
			if err := store.WriteChangeSessions(ctx, change.ID, sessionIDs, now); err != nil {
				return err
			}
			attached += len(sessionIDs)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return attached, nil
}

type hunkLinkSession struct {
	sessionID string
	tool      string
}

// agentSessionsByCommit groups agent-authored links by the commit their hunk
// belongs to, preserving first-seen order so the cap keeps the earliest
// evidence rather than an arbitrary slice of a map.
func agentSessionsByCommit(links []matcher.HunkLink) (map[string][]hunkLinkSession, []string) {
	bySHA := map[string][]hunkLinkSession{}
	seen := map[string]struct{}{}
	var order []string
	for _, link := range links {
		if link.Authorship != matcher.AuthorshipAgent {
			continue
		}
		if link.Tier <= 0 || link.Tier > matcher.TierFuzzy {
			continue
		}
		sessionID := strings.TrimSpace(link.SessionID)
		if sessionID == "" {
			continue
		}
		// HunkID is "<commit sha>:<path>:<start>-<end>" and a path may itself
		// contain a colon, so only the first field is ours to read.
		sha, _, ok := strings.Cut(link.HunkID, ":")
		sha = strings.TrimSpace(sha)
		if !ok || sha == "" {
			continue
		}
		key := sha + "\x00" + sessionID
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		if _, known := bySHA[sha]; !known {
			order = append(order, sha)
		}
		bySHA[sha] = append(bySHA[sha], hunkLinkSession{sessionID: sessionID, tool: strings.TrimSpace(link.Tool)})
	}
	return bySHA, order
}
