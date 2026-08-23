// Package update reads the published CLI manifest and reports whether the
// running binary is behind it.
//
// It exists because nothing in the CLI ever looked. The MCP wrapper checks the
// same manifest and appends a notice to its tool results, so an agent driving
// gx learns about a new build — but someone running `gx review` in a terminal
// had no way to find out, and stayed on whatever they first installed. A user
// spent days on a build that could not log in at all, and the only reason it
// was ever found was that they said so out loud.
package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/version"
)

// DefaultBaseURL is where the release workflow publishes. It is not baked at
// build time because the manifest it serves is public and the same for every
// build; scripts/render-cli-manifest.sh and the MCP wrapper hardcode it too.
const DefaultBaseURL = "https://download.gx.run"

// checkInterval is how long a check is reused. A notice that a new build
// exists does not get more useful for being printed on every command, and the
// request is on the path of an interactive tool.
const checkInterval = 24 * time.Hour

// fetchTimeout bounds the request. The notice is a courtesy: no gx command
// should get measurably slower because download.gx.run is having a bad day.
const fetchTimeout = 2 * time.Second

// Asset is one published archive.
type Asset struct {
	URL          string `json:"url"`
	VersionedURL string `json:"versioned_url"`
	SHA256       string `json:"sha256"`
	Size         int64  `json:"size"`
}

// Manifest is cli/manifest.json.
type Manifest struct {
	SchemaVersion  int              `json:"schema_version"`
	Version        string           `json:"version"`
	GitSHA         string           `json:"git_sha"`
	PublishedAt    string           `json:"published_at"`
	InstallCommand string           `json:"install_command"`
	Assets         map[string]Asset `json:"assets"`
}

// BaseURL is the download host, overridable for testing and for anyone running
// their own mirror. Named to match install.sh's GX_INSTALL_BASE_URL sibling.
func BaseURL() string {
	if custom := strings.TrimSpace(os.Getenv("GX_DOWNLOAD_BASE_URL")); custom != "" {
		return strings.TrimRight(custom, "/")
	}
	return DefaultBaseURL
}

// FetchManifest reads the published manifest.
func FetchManifest(ctx context.Context, client *http.Client) (Manifest, error) {
	if client == nil {
		client = &http.Client{Timeout: fetchTimeout}
	}
	url := BaseURL() + "/cli/manifest.json"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Manifest{}, fmt.Errorf("create manifest request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return Manifest{}, fmt.Errorf("request %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return Manifest{}, fmt.Errorf("request %s: status %s", url, resp.Status)
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Manifest{}, fmt.Errorf("read manifest: %w", err)
	}
	var manifest Manifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return Manifest{}, fmt.Errorf("decode manifest: %w", err)
	}
	if strings.TrimSpace(manifest.GitSHA) == "" {
		return Manifest{}, fmt.Errorf("manifest has no git_sha")
	}
	return manifest, nil
}

// Behind reports whether the running build differs from the published one.
//
// The comparison is on the full commit sha, not the version string, because
// the two builds of one commit are named differently: a push to main publishes
// under a 12-character sha and a tag publishes under its semver, and both
// describe identical binaries. Comparing shas needs no version parsing and
// cannot be confused by that.
//
// It therefore depends on Go's VCS stamping, which release builds carry and a
// build from a linked git worktree does not. A binary with no revision reports
// "not behind" rather than guessing: a wrong nag on every command is worse
// than a missed one, and the builds that lack stamping are development builds
// that Skip already excludes.
func Behind(info version.Info, manifest Manifest) bool {
	local := strings.TrimSpace(info.Revision)
	remote := strings.TrimSpace(manifest.GitSHA)
	if local == "" || remote == "" {
		return false
	}
	return !strings.EqualFold(local, remote)
}

// Skip reports whether this build should not be checked at all, and is
// deliberately generous about it.
//
// A development build is the main case: `just build` bakes version.Version as
// "dev", and its commit is whatever the author has checked out, so a sha
// comparison against the published build says "behind" on work that is
// actually ahead. CI is the other: a notice nobody reads in a log nobody opens
// is pure noise.
func Skip(info version.Info) bool {
	if optedOut() {
		return true
	}
	if strings.TrimSpace(os.Getenv("CI")) != "" {
		return true
	}
	release := strings.TrimSpace(info.Release)
	if release == "" || strings.EqualFold(release, "dev") {
		return true
	}
	return info.Modified
}

func optedOut() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("GX_UPDATE_CHECK"))) {
	case "0", "false", "no", "off":
		return true
	}
	return false
}

// AutoEnabled reports whether an out-of-date binary replaces itself.
//
// Default on: the point is that people run the current build without having to
// be told to, and every case where updating is the wrong thing — a development
// build, a dirty tree, CI, a git hook — is already excluded before this is
// consulted. GX_AUTO_UPDATE=0 leaves the notice and nothing else.
func AutoEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("GX_AUTO_UPDATE"))) {
	case "0", "false", "no", "off":
		return false
	}
	return true
}

type cacheFile struct {
	CheckedAt time.Time `json:"checked_at"`
	Manifest  Manifest  `json:"manifest"`
	// AttemptedAt is when a self-update last ran, successfully or not.
	//
	// Separate from CheckedAt because they rate-limit different things. A
	// failed update — no write permission, a broken mirror, a dropped
	// connection — leaves the binary behind, so the next command would find it
	// behind again and re-download tens of megabytes. That repeats on every
	// command until the cause is fixed. This bounds it to one attempt per
	// interval.
	AttemptedAt time.Time `json:"attempted_at,omitempty"`
}

func cachePath() (string, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "update-check.json"), nil
}

func readCache() (cacheFile, bool) {
	path, err := cachePath()
	if err != nil {
		return cacheFile{}, false
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return cacheFile{}, false
	}
	var cached cacheFile
	if err := json.Unmarshal(raw, &cached); err != nil {
		return cacheFile{}, false
	}
	return cached, true
}

func writeCacheFile(cached cacheFile) {
	path, err := cachePath()
	if err != nil {
		return
	}
	raw, err := json.MarshalIndent(cached, "", "  ")
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	// Best effort throughout: a check that cannot cache is still a check, and
	// none of this is worth failing a command over.
	_ = os.WriteFile(path, append(raw, '\n'), 0o644)
}

func writeCache(manifest Manifest, now time.Time) {
	previous, _ := readCache()
	writeCacheFile(cacheFile{
		CheckedAt:   now.UTC(),
		Manifest:    manifest,
		AttemptedAt: previous.AttemptedAt,
	})
}

// AutoUpdateDue reports whether enough time has passed to try again.
func AutoUpdateDue(now time.Time) bool {
	cached, ok := readCache()
	if !ok || cached.AttemptedAt.IsZero() {
		return true
	}
	return now.Sub(cached.AttemptedAt) >= checkInterval
}

// MarkAutoUpdateAttempt records an attempt, so a failing one is not retried on
// every command until it starts working.
func MarkAutoUpdateAttempt(now time.Time) {
	cached, _ := readCache()
	cached.AttemptedAt = now.UTC()
	writeCacheFile(cached)
}

// Latest returns the published manifest, from cache when it was read recently.
//
// Errors are values, not failures to propagate: every caller in the CLI treats
// "could not check" as "say nothing".
func Latest(ctx context.Context, client *http.Client, now time.Time) (Manifest, error) {
	if cached, ok := readCache(); ok && now.Sub(cached.CheckedAt) < checkInterval {
		return cached.Manifest, nil
	}
	manifest, err := FetchManifest(ctx, client)
	if err != nil {
		return Manifest{}, err
	}
	writeCache(manifest, now)
	return manifest, nil
}

// Notice is the one line printed after a command when a newer build exists.
func Notice(manifest Manifest) string {
	name := strings.TrimSpace(manifest.Version)
	if name == "" {
		name = shortSHA(manifest.GitSHA)
	}
	return fmt.Sprintf("gx %s is available (you have %s) — run `gx update`", name, version.Current())
}

func shortSHA(sha string) string {
	sha = strings.TrimSpace(sha)
	if len(sha) <= 8 {
		return sha
	}
	return sha[:8]
}
