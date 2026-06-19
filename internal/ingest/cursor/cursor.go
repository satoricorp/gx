// Package cursor reads Cursor's local chat/composer database and copies new
// sessions and messages into the gx SQLite store.
//
// Cursor stores every chat/composer/agent thread under
// ~/Library/Application Support/Cursor/User/globalStorage/state.vscdb in the
// `cursorDiskKV` table, keyed by:
//
//   - composerData:<composerId>    composer metadata (name, createdAt, lastUpdatedAt)
//   - bubbleId:<composerId>:<id>   individual user/assistant messages
//
// We mirror each composer as a gx `sessions` row (source='cursor') and each
// bubble as a row in `cursor_messages`. Bubble ids are stable Cursor UUIDs, so
// re-ingest is idempotent via INSERT OR IGNORE.
package cursor

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/version"
)

const (
	sessionIDPrefix = "cursor-"
	sessionSource   = "cursor"
)

// Result reports what a single Sync pass added to the gx store.
type Result struct {
	Composers      int
	NewSessions    int
	NewMessages    int
	VSCDBPath      string
	WorkspacesSeen int
}

// SessionWriter is the subset of storage we need; lets us test without a real DB.
type SessionWriter interface {
	UpsertCursorSession(ctx context.Context, session storage.Session) (bool, error)
	UpsertCursorMessage(ctx context.Context, msg storage.CursorMessage) (bool, error)
	TouchSession(ctx context.Context, sessionID string, lastSeenAt int64) error
	UpdateSessionWorkspace(ctx context.Context, sessionID, cwd, repoRoot string) error
}

// Sync reads Cursor's global state.vscdb and writes any new composers + bubbles
// into store. Safe to call repeatedly.
func Sync(ctx context.Context, store SessionWriter) (Result, error) {
	path, err := DefaultVSCDBPath()
	if err != nil {
		return Result{}, err
	}
	return SyncFromPath(ctx, store, path)
}

// SyncFromPath is Sync with an explicit vscdb location (for tests and manual runs).
func SyncFromPath(ctx context.Context, store SessionWriter, vscdbPath string) (Result, error) {
	result := Result{VSCDBPath: vscdbPath}
	if _, err := os.Stat(vscdbPath); err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return result, fmt.Errorf("stat cursor vscdb: %w", err)
	}

	db, err := openReadOnly(vscdbPath)
	if err != nil {
		return result, err
	}
	defer db.Close()

	composers, err := loadComposers(ctx, db)
	if err != nil {
		return result, err
	}
	bubbles, err := loadBubbles(ctx, db)
	if err != nil {
		return result, err
	}
	workspaceFolders := loadWorkspaceFolders()
	result.WorkspacesSeen = len(workspaceFolders)
	result.Composers = len(composers)

	for composerID, bs := range bubbles {
		if len(bs) == 0 {
			continue
		}
		meta := composers[composerID]
		sessionID := sessionIDPrefix + composerID
		createdAt := composerCreatedAt(meta, bs)
		lastSeen := composerLastSeen(meta, bs)
		cwd, repoRoot := pickWorkspace(workspaceFolders, bs)
		command := composerCommand(meta)
		processName := composerID

		inserted, err := store.UpsertCursorSession(ctx, storage.Session{
			ID:          sessionID,
			CreatedAt:   createdAt,
			Command:     command,
			Cwd:         firstNonEmpty(cwd, "."),
			GXVersion:   version.Current(),
			Source:      ptrString(sessionSource),
			ProcessName: ptrString(processName),
			LastSeenAt:  ptrInt64(lastSeen),
			RepoRoot:    ptrStringIfSet(repoRoot),
		})
		if err != nil {
			return result, err
		}
		if inserted {
			result.NewSessions++
		} else {
			if err := store.TouchSession(ctx, sessionID, lastSeen); err != nil {
				return result, err
			}
			if err := store.UpdateSessionWorkspace(ctx, sessionID, cwd, repoRoot); err != nil {
				return result, err
			}
		}

		for _, b := range bs {
			added, err := store.UpsertCursorMessage(ctx, storage.CursorMessage{
				ID:           b.id,
				SessionID:    sessionID,
				CreatedAt:    b.createdAt,
				Role:         b.role,
				Text:         b.text,
				RawJSON:      b.raw,
				InputTokens:  b.inputTokens,
				OutputTokens: b.outputTokens,
			})
			if err != nil {
				return result, err
			}
			if added {
				result.NewMessages++
			}
		}
	}

	return result, nil
}

// DefaultVSCDBPath returns the macOS location of Cursor's global state.vscdb.
func DefaultVSCDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if runtime.GOOS != "darwin" {
		return "", fmt.Errorf("cursor ingest is only supported on macOS in v0")
	}
	return filepath.Join(home, "Library", "Application Support", "Cursor", "User", "globalStorage", "state.vscdb"), nil
}

func openReadOnly(path string) (*sql.DB, error) {
	dsn := "file:" + path + "?mode=ro&immutable=0&_pragma=busy_timeout(5000)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open cursor vscdb: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping cursor vscdb: %w", err)
	}
	return db, nil
}

type composerMeta struct {
	name          string
	createdAt     int64
	lastUpdatedAt int64
}

func loadComposers(ctx context.Context, db *sql.DB) (map[string]composerMeta, error) {
	rows, err := db.QueryContext(ctx, `SELECT key, value FROM cursorDiskKV WHERE key LIKE 'composerData:%'`)
	if err != nil {
		return nil, fmt.Errorf("query composers: %w", err)
	}
	defer rows.Close()

	out := map[string]composerMeta{}
	for rows.Next() {
		var key string
		var value []byte
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("scan composer: %w", err)
		}
		composerID := strings.TrimPrefix(key, "composerData:")
		if composerID == "" {
			continue
		}
		var doc struct {
			Name          string `json:"name"`
			CreatedAt     int64  `json:"createdAt"`
			LastUpdatedAt int64  `json:"lastUpdatedAt"`
		}
		_ = json.Unmarshal(value, &doc)
		out[composerID] = composerMeta{
			name:          doc.Name,
			createdAt:     doc.CreatedAt,
			lastUpdatedAt: doc.LastUpdatedAt,
		}
	}
	return out, rows.Err()
}

type bubble struct {
	id           string
	composerID   string
	role         string
	text         string
	raw          []byte
	createdAt    int64
	inputTokens  *int
	outputTokens *int
}

func loadBubbles(ctx context.Context, db *sql.DB) (map[string][]bubble, error) {
	rows, err := db.QueryContext(ctx, `SELECT key, value FROM cursorDiskKV WHERE key LIKE 'bubbleId:%'`)
	if err != nil {
		return nil, fmt.Errorf("query bubbles: %w", err)
	}
	defer rows.Close()

	out := map[string][]bubble{}
	for rows.Next() {
		var key string
		var value []byte
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("scan bubble: %w", err)
		}
		composerID, bubbleID, ok := parseBubbleKey(key)
		if !ok {
			continue
		}
		var doc struct {
			Type       int    `json:"type"`
			Text       string `json:"text"`
			TokenCount struct {
				InputTokens  int `json:"inputTokens"`
				OutputTokens int `json:"outputTokens"`
			} `json:"tokenCount"`
			TimingInfo struct {
				ClientRpcSendTime int64 `json:"clientRpcSendTime"`
				ClientEndTime     int64 `json:"clientEndTime"`
			} `json:"timingInfo"`
		}
		_ = json.Unmarshal(value, &doc)

		b := bubble{
			id:         bubbleID,
			composerID: composerID,
			role:       roleFromType(doc.Type),
			text:       doc.Text,
			raw:        append([]byte(nil), value...),
			createdAt:  firstNonZeroInt64(doc.TimingInfo.ClientRpcSendTime, doc.TimingInfo.ClientEndTime),
		}
		if doc.TokenCount.InputTokens > 0 {
			v := doc.TokenCount.InputTokens
			b.inputTokens = &v
		}
		if doc.TokenCount.OutputTokens > 0 {
			v := doc.TokenCount.OutputTokens
			b.outputTokens = &v
		}
		out[composerID] = append(out[composerID], b)
	}
	return out, rows.Err()
}

func parseBubbleKey(key string) (composerID, bubbleID string, ok bool) {
	rest := strings.TrimPrefix(key, "bubbleId:")
	if rest == key {
		return "", "", false
	}
	sep := strings.IndexByte(rest, ':')
	if sep <= 0 || sep == len(rest)-1 {
		return "", "", false
	}
	return rest[:sep], rest[sep+1:], true
}

func roleFromType(t int) string {
	switch t {
	case 1:
		return "user"
	case 2:
		return "assistant"
	default:
		return "unknown"
	}
}

func composerCreatedAt(meta composerMeta, bs []bubble) int64 {
	if meta.createdAt > 0 {
		return meta.createdAt
	}
	var earliest int64
	for _, b := range bs {
		if b.createdAt > 0 && (earliest == 0 || b.createdAt < earliest) {
			earliest = b.createdAt
		}
	}
	return earliest
}

func composerLastSeen(meta composerMeta, bs []bubble) int64 {
	latest := meta.lastUpdatedAt
	for _, b := range bs {
		if b.createdAt > latest {
			latest = b.createdAt
		}
	}
	return latest
}

func composerCommand(meta composerMeta) string {
	if meta.name != "" {
		return "cursor: " + meta.name
	}
	return "cursor"
}

// loadWorkspaceFolders walks Cursor's per-workspace storage and returns the
// `folder` URIs declared in each workspace.json. We use these to attribute a
// composer to a repo when the bubble payload references a workspace.
func loadWorkspaceFolders() map[string]string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	root := filepath.Join(home, "Library", "Application Support", "Cursor", "User", "workspaceStorage")
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	out := map[string]string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, entry.Name(), "workspace.json"))
		if err != nil {
			continue
		}
		var doc struct {
			Folder string `json:"folder"`
		}
		if err := json.Unmarshal(data, &doc); err != nil || doc.Folder == "" {
			continue
		}
		path := uriToPath(doc.Folder)
		if !isMeaningfulFolder(path) {
			continue
		}
		out[entry.Name()] = path
	}
	return out
}

// isMeaningfulFolder rejects paths that would trivially match every bubble
// payload via substring search (root, single-segment system paths, empty).
func isMeaningfulFolder(path string) bool {
	trimmed := strings.TrimRight(path, "/")
	if trimmed == "" {
		return false
	}
	if strings.Count(trimmed, "/") < 2 {
		return false
	}
	return true
}

func uriToPath(folderURI string) string {
	u, err := url.Parse(folderURI)
	if err != nil {
		return folderURI
	}
	if u.Scheme == "file" {
		return u.Path
	}
	return folderURI
}

// pickWorkspace picks the most-frequent workspace folder referenced inside the
// bubble payloads. v0 keeps it simple: scan the raw JSON for "file://" URIs and
// match against known workspace folders. If nothing matches, return ("","").
func pickWorkspace(workspaces map[string]string, bs []bubble) (cwd, repoRoot string) {
	if len(workspaces) == 0 {
		return "", ""
	}
	counts := map[string]int{}
	for _, b := range bs {
		for _, folder := range workspaces {
			if folder != "" && strings.Contains(string(b.raw), folder) {
				counts[folder]++
			}
		}
	}
	var best string
	bestN := 0
	for folder, n := range counts {
		if n > bestN {
			best = folder
			bestN = n
		}
	}
	if best == "" {
		return "", ""
	}
	return best, detectRepoRoot(best)
}

func detectRepoRoot(cwd string) string {
	for dir := cwd; dir != "" && dir != string(filepath.Separator); dir = filepath.Dir(dir) {
		if exists(filepath.Join(dir, ".git")) || exists(filepath.Join(dir, ".jj")) {
			return dir
		}
	}
	return ""
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ptrString(v string) *string {
	return &v
}

func ptrStringIfSet(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func ptrInt64(v int64) *int64 {
	return &v
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func firstNonZeroInt64(values ...int64) int64 {
	for _, v := range values {
		if v != 0 {
			return v
		}
	}
	return 0
}
