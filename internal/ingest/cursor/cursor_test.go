package cursor

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/satoricorp/gx/internal/storage"
)

type fakeStore struct {
	sessions map[string]storage.Session
	messages map[string]storage.CursorMessage
	touches  []touch
}

type touch struct {
	sessionID  string
	lastSeenAt int64
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		sessions: map[string]storage.Session{},
		messages: map[string]storage.CursorMessage{},
	}
}

func (f *fakeStore) UpsertCursorSession(_ context.Context, s storage.Session) (bool, error) {
	if _, ok := f.sessions[s.ID]; ok {
		return false, nil
	}
	f.sessions[s.ID] = s
	return true, nil
}

func (f *fakeStore) UpsertCursorMessage(_ context.Context, m storage.CursorMessage) (bool, error) {
	if _, ok := f.messages[m.ID]; ok {
		return false, nil
	}
	f.messages[m.ID] = m
	return true, nil
}

func (f *fakeStore) TouchSession(_ context.Context, sessionID string, lastSeenAt int64) error {
	f.touches = append(f.touches, touch{sessionID, lastSeenAt})
	return nil
}

func (f *fakeStore) UpdateSessionWorkspace(_ context.Context, sessionID, cwd, repoRoot string) error {
	if session, ok := f.sessions[sessionID]; ok {
		if cwd != "" {
			session.Cwd = cwd
		}
		if repoRoot != "" {
			session.RepoRoot = ptrStringIfSet(repoRoot)
		}
		f.sessions[sessionID] = session
	}
	return nil
}

func buildFakeVSCDB(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "state.vscdb")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE cursorDiskKV (key TEXT UNIQUE ON CONFLICT REPLACE, value BLOB)`); err != nil {
		t.Fatalf("create table: %v", err)
	}

	rows := []struct {
		key   string
		value string
	}{
		{
			key:   "composerData:c1",
			value: `{"composerId":"c1","name":"refactor auth","createdAt":1700000000000,"lastUpdatedAt":1700000005000}`,
		},
		{
			key:   "bubbleId:c1:b1",
			value: `{"bubbleId":"b1","type":1,"text":"please fix the bug","timingInfo":{"clientRpcSendTime":1700000001000,"clientEndTime":1700000001500},"tokenCount":{"inputTokens":42,"outputTokens":0}}`,
		},
		{
			key:   "bubbleId:c1:b2",
			value: `{"bubbleId":"b2","type":2,"text":"done","timingInfo":{"clientRpcSendTime":1700000003000,"clientEndTime":1700000004500},"tokenCount":{"inputTokens":0,"outputTokens":17}}`,
		},
		{
			key:   "bubbleId:c2:b3",
			value: `{"bubbleId":"b3","type":1,"text":"orphan composer message","timingInfo":{"clientRpcSendTime":1700000010000}}`,
		},
	}
	for _, r := range rows {
		if _, err := db.Exec(`INSERT INTO cursorDiskKV(key, value) VALUES(?, ?)`, r.key, r.value); err != nil {
			t.Fatalf("insert %s: %v", r.key, err)
		}
	}
	return path
}

func TestSyncFromPathIngestsComposersAndBubbles(t *testing.T) {
	path := buildFakeVSCDB(t)
	store := newFakeStore()

	got, err := SyncFromPath(context.Background(), store, path)
	if err != nil {
		t.Fatalf("SyncFromPath: %v", err)
	}

	if got.NewSessions != 2 {
		t.Fatalf("NewSessions = %d, want 2 (c1 + c2 orphan)", got.NewSessions)
	}
	if got.NewMessages != 3 {
		t.Fatalf("NewMessages = %d, want 3", got.NewMessages)
	}

	c1, ok := store.sessions["cursor-c1"]
	if !ok {
		t.Fatalf("expected cursor-c1 session")
	}
	if c1.Source == nil || *c1.Source != "cursor" {
		t.Fatalf("c1 source = %v, want cursor", c1.Source)
	}
	if c1.Command != "cursor: refactor auth" {
		t.Fatalf("c1 command = %q", c1.Command)
	}
	if c1.CreatedAt != 1700000000000 {
		t.Fatalf("c1 createdAt = %d", c1.CreatedAt)
	}
	if c1.LastSeenAt == nil || *c1.LastSeenAt != 1700000005000 {
		t.Fatalf("c1 lastSeenAt = %v", c1.LastSeenAt)
	}
	if c1.ProcessName == nil || *c1.ProcessName != "c1" {
		t.Fatalf("c1 processName = %v", c1.ProcessName)
	}

	b1, ok := store.messages["b1"]
	if !ok {
		t.Fatalf("expected b1 message")
	}
	if b1.Role != "user" {
		t.Fatalf("b1 role = %q", b1.Role)
	}
	if b1.Text != "please fix the bug" {
		t.Fatalf("b1 text = %q", b1.Text)
	}
	if b1.InputTokens == nil || *b1.InputTokens != 42 {
		t.Fatalf("b1 inputTokens = %v", b1.InputTokens)
	}
	if b1.OutputTokens != nil {
		t.Fatalf("b1 outputTokens = %v, want nil (was zero)", b1.OutputTokens)
	}

	b2 := store.messages["b2"]
	if b2.Role != "assistant" {
		t.Fatalf("b2 role = %q", b2.Role)
	}
	if b2.OutputTokens == nil || *b2.OutputTokens != 17 {
		t.Fatalf("b2 outputTokens = %v", b2.OutputTokens)
	}
}

func TestSyncFromPathIsIdempotent(t *testing.T) {
	path := buildFakeVSCDB(t)
	store := newFakeStore()

	if _, err := SyncFromPath(context.Background(), store, path); err != nil {
		t.Fatalf("first sync: %v", err)
	}
	got, err := SyncFromPath(context.Background(), store, path)
	if err != nil {
		t.Fatalf("second sync: %v", err)
	}
	if got.NewSessions != 0 || got.NewMessages != 0 {
		t.Fatalf("second sync inserted %d sessions / %d messages, want 0/0", got.NewSessions, got.NewMessages)
	}
	if len(store.touches) == 0 {
		t.Fatalf("expected TouchSession calls on second sync, got none")
	}
}

func TestSyncFromPathMissingFile(t *testing.T) {
	store := newFakeStore()
	got, err := SyncFromPath(context.Background(), store, filepath.Join(t.TempDir(), "absent.vscdb"))
	if err != nil {
		t.Fatalf("missing file should not error, got %v", err)
	}
	if got.NewSessions != 0 || got.NewMessages != 0 {
		t.Fatalf("expected zero work on missing file")
	}
}

func TestIsMeaningfulFolder(t *testing.T) {
	cases := []struct {
		path string
		ok   bool
	}{
		{"", false},
		{"/", false},
		{"/Users", false},
		{"/Users/alice", true},
		{"/Users/alice/src/gx", true},
		{"/Users/alice/src/gx/", true},
	}
	for _, c := range cases {
		if got := isMeaningfulFolder(c.path); got != c.ok {
			t.Errorf("isMeaningfulFolder(%q) = %v, want %v", c.path, got, c.ok)
		}
	}
}

func TestParseBubbleKey(t *testing.T) {
	cases := []struct {
		key      string
		composer string
		bubble   string
		ok       bool
	}{
		{"bubbleId:c1:b1", "c1", "b1", true},
		{"bubbleId:c1:b1:nested", "c1", "b1:nested", true},
		{"bubbleId:c1:", "", "", false},
		{"bubbleId::b1", "", "", false},
		{"composerData:c1", "", "", false},
	}
	for _, c := range cases {
		composer, bubble, ok := parseBubbleKey(c.key)
		if composer != c.composer || bubble != c.bubble || ok != c.ok {
			t.Errorf("parseBubbleKey(%q) = (%q,%q,%v); want (%q,%q,%v)", c.key, composer, bubble, ok, c.composer, c.bubble, c.ok)
		}
	}
}
