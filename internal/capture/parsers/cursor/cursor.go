package cursor

import (
	"context"
	"encoding/json"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/satoricorp/totality/internal/capture"
	"github.com/satoricorp/totality/internal/capture/repopath"
)

// Parser normalizes Cursor composer bubbles from state.vscdb into SessionEvents.
type Parser struct {
	Inventory    *capture.InventoryCollector
	Since        time.Time
	Until        time.Time
	contentCache map[string]string
}

func (p *Parser) Tool() string { return capture.ToolCursor }

// ParseFile reads Cursor's global state.vscdb in place, read-only.
//
// The database routinely reaches multiple gigabytes, but most of that is the
// Agents-Window blob store (agentKv:*, unreadable here) and bubbles of
// long-dead sessions. The old approach copied the whole file per push and
// then table-scanned it, which is why oversized databases had to be skipped
// outright. Instead: open the live file read-only (WAL readers do not block
// Cursor's writer), prune composers on their metadata timestamps, and fetch
// bubbles per surviving composer through the key index — a push reads
// megabytes, not gigabytes.
func (p *Parser) ParseFile(path string, repoRoot string) ([]capture.SessionEvent, error) {
	if strings.HasSuffix(path, ".jsonl") {
		return p.parseTranscriptFile(path, repoRoot)
	}
	if p.Inventory != nil {
		p.Inventory.RecordSampleFile(capture.ToolCursor, path)
	}

	db, err := openReadOnly(path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	ctx := context.Background()
	p.contentCache = map[string]string{}
	composers, err := loadComposers(ctx, db)
	if err != nil {
		return nil, err
	}
	composerIDs, err := bubbleComposerIDs(ctx, db)
	if err != nil {
		return nil, err
	}
	workspaces := loadWorkspaceFolders()

	var events []capture.SessionEvent
	for _, composerID := range composerIDs {
		meta, known := composers[composerID]
		if known && !metaMayOverlapWindow(meta, p.Since, p.Until) {
			continue
		}
		bs, err := loadComposerBubbles(ctx, db, composerID)
		if err != nil {
			return nil, err
		}
		if len(bs) == 0 {
			continue
		}
		cwd, detectedRepo := pickWorkspace(workspaces, bs)
		if repoRoot != "" && detectedRepo != "" {
			absRepo, err := filepath.Abs(repoRoot)
			if err == nil {
				absDetected, err := filepath.Abs(detectedRepo)
				if err == nil && absDetected != absRepo {
					continue
				}
			}
		} else if repoRoot != "" && cwd != "" {
			absRepo, err := filepath.Abs(repoRoot)
			if err == nil {
				absCwd, err := filepath.Abs(cwd)
				if err == nil && !strings.HasPrefix(absCwd, absRepo+string(filepath.Separator)) && absCwd != absRepo {
					continue
				}
			}
		}

		sessionID := "cursor-" + composerID
		lastUpdated := composerLastSeen(meta, bs)
		if !inTimeWindow(lastUpdated, p.Since, p.Until) && !bubblesInWindow(bs, p.Since, p.Until) {
			continue
		}

		for _, b := range bs {
			if !inTimeWindow(b.createdAt, p.Since, p.Until) {
				continue
			}
			parsed := p.parseBubble(db, b, sessionID, repoRoot)
			events = append(events, parsed...)
		}
	}
	return events, nil
}

var patchFileRE = regexp.MustCompile(`(?m)^\*\*\* (?:Add|Update) File: (.+)$`)

func parsePatchText(text, sessionID, model string, ts int64, repoRoot string, raw map[string]json.RawMessage) []capture.SessionEvent {
	var events []capture.SessionEvent
	lines := strings.Split(text, "\n")
	var currentFile string
	var added []string
	flush := func() {
		if currentFile == "" || len(added) == 0 {
			return
		}
		events = append(events, capture.SessionEvent{
			SessionID: sessionID,
			Tool:      capture.ToolCursor,
			Model:     model,
			TS:        ts,
			Kind:      capture.KindEdit,
			FilePath:  relPath(currentFile, repoRoot),
			NewText:   strings.Join(added, "\n"),
			Raw:       cloneRaw(raw),
		})
		added = nil
	}
	for _, line := range lines {
		if m := patchFileRE.FindStringSubmatch(line); len(m) == 2 {
			flush()
			currentFile = strings.TrimSpace(m[1])
			continue
		}
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			added = append(added, strings.TrimPrefix(line, "+"))
		}
	}
	flush()
	return events
}

func rawString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	return strings.Trim(string(raw), `"`)
}

func rawInt64(raw json.RawMessage) int64 {
	if len(raw) == 0 {
		return 0
	}
	var n int64
	if err := json.Unmarshal(raw, &n); err == nil {
		return n
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return int64(f)
	}
	return 0
}

// relPath defers to repopath, which knows a repository can have more than one
// checkout. Relativizing against the pushing checkout alone silently mangles
// every edit made in a linked worktree.
func relPath(path, repoRoot string) string {
	return repopath.Rel(path, repoRoot)
}

func cloneRaw(obj map[string]json.RawMessage) map[string]json.RawMessage {
	out := make(map[string]json.RawMessage, len(obj))
	for k, v := range obj {
		out[k] = v
	}
	return out
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func normalizeToolName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.TrimSuffix(name, "_v2")
	return name
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// DiscoverVSCDBPath is deliberately gone. It took a home directory, ignored it
// (`_ = home`) and returned os.UserHomeDir()'s Cursor database instead, so any
// caller sweeping a specific home silently got the current user's real global
// database. Session discovery resolves the path under the home it was given;
// see cursorVSCDBPath in internal/capture/parsers.
