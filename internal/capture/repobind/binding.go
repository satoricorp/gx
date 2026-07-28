package repobind

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// Binding is what a session's working directory resolves to.
//
// Origin is the identity that matters: it is what gets compared against the
// repositories an organization has connected, and it is the same for every
// checkout of a repository. RepoRoot and GitCommonDir are local facts, useful
// for diagnostics and for recognizing that two worktrees are one repository,
// but they must never be used as the gate — paths are a machine's opinion,
// origins are the repository's identity.
//
// A zero Binding means "not identifiable", which is the safe answer: a
// directory that is not a git repository, or a repository with no origin, does
// not belong to any connected repo and its sessions stay local.
type Binding struct {
	Dir          string
	RepoRoot     string
	GitCommonDir string
	OriginURL    string
	Origin       string
}

// Bound reports whether the directory resolved to an identifiable repository.
func (b Binding) Bound() bool { return b.Origin != "" }

// Resolver resolves working directories to repositories, caching per directory.
//
// A session mentions its working directory on nearly every entry — over a
// thousand times in a single sampled transcript — so this is asked far more
// often than it changes.
type Resolver struct {
	mu    sync.Mutex
	cache map[string]Binding
}

// NewResolver returns a Resolver with an empty cache.
func NewResolver() *Resolver { return &Resolver{cache: map[string]Binding{}} }

// ForDirectory resolves one working directory.
func (r *Resolver) ForDirectory(ctx context.Context, dir string) Binding {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return Binding{}
	}
	if abs, err := filepath.Abs(dir); err == nil {
		dir = abs
	}

	r.mu.Lock()
	if cached, ok := r.cache[dir]; ok {
		r.mu.Unlock()
		return cached
	}
	r.mu.Unlock()

	binding := resolveDirectory(ctx, dir)

	r.mu.Lock()
	r.cache[dir] = binding
	r.mu.Unlock()
	return binding
}

func resolveDirectory(ctx context.Context, dir string) Binding {
	binding := Binding{Dir: dir}

	// --show-toplevel from inside a linked worktree reports that worktree, so
	// this is the checkout the session ran in rather than the repository.
	root, ok := git(ctx, dir, "rev-parse", "--show-toplevel")
	if !ok {
		// Not a git repository. Nothing to bind to, and no fallback: guessing
		// from the path is how an unconnected repo gets mistaken for a
		// connected one.
		return binding
	}
	binding.RepoRoot = root

	// The common dir is shared by every worktree of a repository, so two
	// sessions run in different worktrees resolve to one repository.
	if common, ok := git(ctx, dir, "rev-parse", "--git-common-dir"); ok {
		if !filepath.IsAbs(common) {
			common = filepath.Join(root, common)
		}
		binding.GitCommonDir = filepath.Clean(common)
	}

	// A repository with no origin is unidentifiable, not unrestricted.
	raw, ok := git(ctx, dir, "remote", "get-url", "origin")
	if !ok || strings.TrimSpace(raw) == "" {
		return binding
	}
	binding.OriginURL = raw
	binding.Origin = NormalizeOrigin(raw)
	return binding
}

func git(ctx context.Context, dir string, args ...string) (string, bool) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	value := strings.TrimSpace(string(out))
	if value == "" {
		return "", false
	}
	return value, true
}

// BindEvents decides which repository a parsed session belongs to.
//
// The directory the session spent most of its events in wins. A session that
// moves between repositories is real — an agent asked about one project while
// working in another — and the alternative, refusing to bind a session that was
// ever elsewhere, would drop most long sessions.
//
// When the tool recorded a remote itself, that is used as a fallback for a
// directory that no longer resolves. Checkouts get deleted and worktrees get
// removed, and a session should not become unattributable because the folder it
// ran in is gone. Codex records one; most tools do not.
func BindEvents(ctx context.Context, r *Resolver, events []Event) Binding {
	if r == nil {
		r = NewResolver()
	}
	counts := map[string]int{}
	originURL := ""
	for _, e := range events {
		if dir := strings.TrimSpace(e.Directory()); dir != "" {
			counts[dir]++
		}
		if originURL == "" {
			originURL = strings.TrimSpace(e.RecordedOrigin())
		}
	}

	best, bestCount := "", 0
	for dir, n := range counts {
		// Ties resolve by path so the answer does not depend on map order.
		if n > bestCount || (n == bestCount && dir < best) {
			best, bestCount = dir, n
		}
	}

	if best != "" {
		if binding := r.ForDirectory(ctx, best); binding.Bound() {
			return binding
		}
		// The directory is gone or is not a repository. Fall through to what
		// the tool recorded rather than reporting the session as unbindable.
		if origin := NormalizeOrigin(originURL); origin != "" {
			return Binding{Dir: best, OriginURL: originURL, Origin: origin}
		}
		return Binding{Dir: best}
	}
	if origin := NormalizeOrigin(originURL); origin != "" {
		return Binding{OriginURL: originURL, Origin: origin}
	}
	return Binding{}
}

// Event is the part of a parsed session event this package needs. Declared here
// rather than importing the capture package so binding stays a leaf dependency.
type Event interface {
	Directory() string
	RecordedOrigin() string
}
