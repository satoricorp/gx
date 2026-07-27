// Package gxtest builds the filesystem and git world a gx test runs in: a
// GX_HOME, real git repositories, linked worktrees, trailer-stamped commits and
// agent transcripts.
//
// It imports nothing from gx beyond the standard library, and that is a
// constraint rather than an accident. internal/storage/storagetest layers the
// database state builder on top of this package, and internal/vcs already
// imports internal/storage — so a gx import here would either cycle or drag
// internal/vcs and its git-probing package init into the fast test binaries of
// internal/storage, internal/reviewbundle and internal/provenance.
//
// The one production constant this forces a copy of is the GX revision trailer
// (see RevisionTrailerFormat), and that copy is pinned against the real one by
// TestGXTestRevisionTrailerMatchesProduction in internal/hooks.
package gxtest

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// RevisionTrailerFormat mirrors vcs.RevisionTrailerFormat.
//
// A commit only produces a `changes` row — and therefore only ever gets
// sessions attached to it — if it carries this trailer, so a harness that got
// the format wrong would build repositories in which every session assertion
// silently passes vacuously.
const RevisionTrailerFormat = "GX: https://gx.run/r/%s"

// World is one isolated gx installation: a GX_HOME nothing else writes to and a
// home directory to hang agent transcripts off.
type World struct {
	// GXHome is the value of GX_HOME for the duration of the test. Every
	// storage.Open in this test resolves ~/.gx to this directory.
	GXHome string
	// Home is the fake user home. It is handed to callers explicitly (as
	// hooks.PushOptions.HomeDir) rather than exported through the HOME
	// environment variable, because moving HOME would also move git's global
	// config out from under the repositories this package creates.
	Home string
}

// NewWorld isolates the test from the developer's real ~/.gx.
//
// storage.DefaultDir falls back to ~/.gx whenever GX_HOME is unset, and 32 rows
// in the author's real database are the receipts of tests that forgot to set
// it. Going through this constructor is what makes that impossible.
func NewWorld(t *testing.T) *World {
	t.Helper()
	world := &World{GXHome: t.TempDir(), Home: t.TempDir()}
	t.Setenv("GX_HOME", world.GXHome)
	// Background upload workers would race the test's own database handles and
	// occasionally reach the network.
	t.Setenv("GX_DISABLE_BACKGROUND_WORKERS", "1")
	return world
}

// Repo is one git worktree plus the git metadata directories that identify it.
//
// GitCommonDir is the repository identity gx keys `repos` rows on, and for a
// linked worktree it is the MAIN checkout's common dir while Root and GitDir
// are the worktree's own. Keeping the three separate is the whole point: a
// harness that collapsed them could not express the state Bug B lived in.
type Repo struct {
	// Root is git's own --show-toplevel, i.e. already symlink-resolved. On
	// macOS t.TempDir hands back /var/folders/... while git reports the real
	// /private/var/folders/... underneath, and comparing the two forms is a
	// recurring source of false failures.
	Root         string
	GitDir       string
	GitCommonDir string
}

// IsLinkedWorktree reports whether this worktree shares another checkout's git
// common dir.
func (r *Repo) IsLinkedWorktree() bool {
	return r.GitDir != r.GitCommonDir
}

// NewRepo creates a git repository with one initial commit.
//
// Two commits matter for push tests: with no prior commit the capture
// orchestrator falls back to <head>~20, which does not resolve in a
// single-commit repository and makes the run bail with "empty ref range"
// before any matching happens. Callers add the second commit with Commit.
func (w *World) NewRepo(t *testing.T) *Repo {
	t.Helper()
	return w.NewRepoAt(t, t.TempDir())
}

// NewRepoAt creates the repository at an explicit path.
func (w *World) NewRepoAt(t *testing.T, dir string) *Repo {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create repo dir %s: %v", dir, err)
	}
	for _, args := range [][]string{
		{"init", "-b", "main"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "Test"},
	} {
		Git(t, dir, args...)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# test\n"), 0o644); err != nil {
		t.Fatalf("write README: %v", err)
	}
	Git(t, dir, "add", "README.md")
	Git(t, dir, "commit", "--no-verify", "-m", "initial")
	return describeRepo(t, dir)
}

// AddWorktree creates a linked worktree on a new branch.
//
// The returned Repo has its own Root and GitDir and the MAIN repository's
// GitCommonDir, which is exactly why a push from here writes its rows under a
// `repos` row whose root_path is some other directory entirely.
func (r *Repo) AddWorktree(t *testing.T, dir, branch string) *Repo {
	t.Helper()
	Git(t, r.Root, "worktree", "add", "-b", branch, dir)
	worktree := describeRepo(t, dir)
	if !worktree.IsLinkedWorktree() {
		t.Fatalf("AddWorktree(%s) produced git_dir == git_common_dir (%s); it is not a linked worktree",
			dir, worktree.GitDir)
	}
	if worktree.GitCommonDir != r.GitCommonDir {
		t.Fatalf("AddWorktree(%s) common dir = %s, want the main repo's %s",
			dir, worktree.GitCommonDir, r.GitCommonDir)
	}
	return worktree
}

// Commit is one trailer-stamped commit.
type Commit struct {
	SHA        string
	RevisionID string
}

// Commit writes files and commits them carrying a GX revision trailer, so the
// push path's RecoverMissingRevisions can create the `changes` row that
// sessions attach to.
func (r *Repo) Commit(t *testing.T, files map[string]string, subject string) Commit {
	t.Helper()
	revisionID := NewRevisionID(t)
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	for _, name := range names {
		path := filepath.Join(r.Root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("create dir for %s: %v", name, err)
		}
		if err := os.WriteFile(path, []byte(files[name]), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
		Git(t, r.Root, "add", name)
	}
	Git(t, r.Root, "commit", "--no-verify", "-m", StampRevisionTrailer(subject, revisionID))
	return Commit{SHA: r.Rev(t, "HEAD"), RevisionID: revisionID}
}

// Rev resolves a git revision to its full SHA.
func (r *Repo) Rev(t *testing.T, ref string) string {
	t.Helper()
	return GitOutput(t, r.Root, "rev-parse", ref)
}

// RefRange is the base..head string a pre-push hook receives.
func (r *Repo) RefRange(t *testing.T, base, head string) string {
	t.Helper()
	return base + ".." + head
}

// Git runs a git command in dir and fails the test if it errors.
func Git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

// GitOutput runs a git command in dir and returns its trimmed stdout.
func GitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out))
}

// NewRevisionID mints a GX revision identifier in the same shape
// vcs.GenerateRevisionID produces: a URL-safe base64 encoding of 128 random
// bits, which vcs.ValidRevisionID accepts.
func NewRevisionID(t *testing.T) string {
	t.Helper()
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		t.Fatalf("generate revision id: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw[:])
}

// RevisionTrailerLine renders the GX identity trailer for a revision.
func RevisionTrailerLine(revisionID string) string {
	revisionID = strings.TrimSpace(revisionID)
	if revisionID == "" {
		return ""
	}
	return fmt.Sprintf(RevisionTrailerFormat, revisionID)
}

// StampRevisionTrailer appends the GX identity trailer to a commit subject.
func StampRevisionTrailer(message, revisionID string) string {
	trailer := RevisionTrailerLine(revisionID)
	if trailer == "" {
		return message
	}
	body := strings.TrimRight(message, "\n")
	if body == "" {
		return trailer
	}
	return body + "\n\n" + trailer
}

// ClaudeProjectSlug is the directory name Claude Code derives from a repo root.
func ClaudeProjectSlug(repoRoot string) string {
	abs, err := filepath.Abs(repoRoot)
	if err != nil {
		abs = repoRoot
	}
	slug := strings.ReplaceAll(abs, string(filepath.Separator), "-")
	if !strings.HasPrefix(slug, "-") {
		slug = "-" + slug
	}
	return slug
}

// Transcript describes one Claude transcript to write.
//
// InFileSessionID is separate from the file path on purpose: Claude stamps
// every subagent transcript with the PARENT conversation's sessionId, so the
// pipeline must take identity from the path and ignore the in-file value.
// Leaving it empty defaults it to ConversationID.
type Transcript struct {
	ConversationID  string
	Subagent        string
	InFileSessionID string
	File            string
	Content         string
}

// WriteClaudeTranscript writes a one-line Claude transcript into the world's
// home directory whose Write tool call produced exactly the committed file, so
// the matcher links it at tier 1. It returns the transcript path.
func (w *World) WriteClaudeTranscript(t *testing.T, repo *Repo, transcript Transcript) string {
	t.Helper()
	projects := filepath.Join(w.Home, ".claude", "projects", ClaudeProjectSlug(repo.Root))
	path := filepath.Join(projects, transcript.ConversationID+".jsonl")
	if transcript.Subagent != "" {
		path = filepath.Join(projects, transcript.ConversationID, "subagents", transcript.Subagent+".jsonl")
	}
	inFileSessionID := transcript.InFileSessionID
	if inFileSessionID == "" {
		inFileSessionID = transcript.ConversationID
	}
	line, err := json.Marshal(map[string]any{
		"type":      "assistant",
		"cwd":       repo.Root,
		"sessionId": inFileSessionID,
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"message": map[string]any{
			"role":  "assistant",
			"model": "claude-probe",
			"content": []map[string]any{{
				"type": "tool_use",
				"id":   "t1",
				"name": "Write",
				"input": map[string]any{
					"file_path": filepath.ToSlash(filepath.Join(repo.Root, transcript.File)),
					"content":   transcript.Content,
				},
			}},
		},
	})
	if err != nil {
		t.Fatalf("marshal transcript: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create transcript dir: %v", err)
	}
	if err := os.WriteFile(path, append(line, '\n'), 0o644); err != nil {
		t.Fatalf("write transcript: %v", err)
	}
	return path
}

func describeRepo(t *testing.T, dir string) *Repo {
	t.Helper()
	root := GitOutput(t, dir, "rev-parse", "--show-toplevel")
	return &Repo{
		Root:         root,
		GitDir:       filepath.Clean(GitOutput(t, root, "rev-parse", "--path-format=absolute", "--git-dir")),
		GitCommonDir: filepath.Clean(GitOutput(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir")),
	}
}
