package cloud

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

	"github.com/satoricorp/totality/internal/storage"
	"github.com/satoricorp/totality/internal/version"
)

// ConnectedRepo is one repository an organization has connected.
type ConnectedRepo struct {
	FullName string `json:"fullName"`
	// Origin is the canonical identity — host, owner, repository, lowercased —
	// rendered by the server so every client compares the same shape.
	Origin string `json:"origin"`
}

type connectedReposResponse struct {
	Repos []ConnectedRepo `json:"repos"`
}

// connectedReposTTL bounds how stale a cached answer may be.
//
// The list changes when somebody installs or removes the GitHub App, which is
// rare, and the cost of being briefly stale is asymmetric: a repository
// connected minutes ago is not yet indexed, while a repository disconnected
// minutes ago is still indexed. An hour keeps the second window short without
// putting the network on the push path every time.
const connectedReposTTL = time.Hour

// ConnectedRepos lists the repositories this organization has connected.
//
// Which repositories are connected is the server's answer, not something the
// client can work out: the local repos table records what tx has seen on this
// machine, which is not the same question and would include repositories nobody
// connected.
func (c *Client) ConnectedRepos(ctx context.Context) ([]ConnectedRepo, error) {
	if c == nil || c.url == "" {
		return nil, fmt.Errorf("tx cloud base URL is not configured")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cloudURLWithPath(c.url, "/v1/connected-repos"), nil)
	if err != nil {
		return nil, fmt.Errorf("create connected repos request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "tx/"+version.Current())
	token, err := CloudAPIToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send connected repos request: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("connected repos: status %s", resp.Status)
	}
	var decoded connectedReposResponse
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return nil, fmt.Errorf("parse connected repos: %w", err)
	}
	return decoded.Repos, nil
}

type connectedReposCache struct {
	FetchedAtMS int64           `json:"fetched_at_ms"`
	Repos       []ConnectedRepo `json:"repos"`
}

// ConnectedOrigins returns the origins of the connected repositories, served
// from a local cache and refreshed when the cache is missing or stale.
//
// The pre-push hook calls this, so the network must not be on its critical
// path. A refresh that fails falls back to whatever was cached: the answer goes
// stale rather than empty, because empty would silently stop indexing every
// repository the moment the machine went offline.
//
// The bool reports whether an answer is available at all. Callers must not read
// an empty list as "nothing is connected" — with no cache and no network there
// is simply no answer, and that is different from a definitive "none".
func ConnectedOrigins(ctx context.Context) ([]string, bool) {
	cached, cachedOK := readConnectedReposCache()
	if cachedOK && time.Since(time.UnixMilli(cached.FetchedAtMS)) < connectedReposTTL {
		return originsOf(cached.Repos), true
	}

	client := NewClient()
	if client == nil {
		return originsOf(cached.Repos), cachedOK
	}
	repos, err := client.ConnectedRepos(ctx)
	if err != nil {
		// Signed out, offline, or the server is unhappy. A stale answer beats
		// no answer; no answer at all is reported as such.
		return originsOf(cached.Repos), cachedOK
	}
	writeConnectedReposCache(connectedReposCache{FetchedAtMS: time.Now().UnixMilli(), Repos: repos})
	return originsOf(repos), true
}

func originsOf(repos []ConnectedRepo) []string {
	out := make([]string, 0, len(repos))
	for _, repo := range repos {
		if origin := strings.TrimSpace(strings.ToLower(repo.Origin)); origin != "" {
			out = append(out, origin)
		}
	}
	return out
}

func connectedReposCachePath() (string, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "connected-repos.json"), nil
}

func readConnectedReposCache() (connectedReposCache, bool) {
	path, err := connectedReposCachePath()
	if err != nil {
		return connectedReposCache{}, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return connectedReposCache{}, false
	}
	var cache connectedReposCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return connectedReposCache{}, false
	}
	return cache, true
}

func writeConnectedReposCache(cache connectedReposCache) {
	path, err := connectedReposCachePath()
	if err != nil {
		return
	}
	// Best effort on purpose: failing to cache costs a refresh next time, and
	// this runs inside a push.
	_ = writeJSONFile(path, cache, 0o600)
}
