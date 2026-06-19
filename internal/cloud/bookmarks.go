package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CloudBookmark struct {
	ID            string  `json:"id"`
	RepoFullName  string  `json:"repo_full_name"`
	BranchName    string  `json:"branch_name"`
	Title         *string `json:"title"`
	Revision      int     `json:"revision"`
	HeadCommitID  *string `json:"head_commit_id"`
	RemoteHeadSha *string `json:"remote_head_sha"`
	MergeStatus   string  `json:"merge_status"`
	UpdatedAtMs   int64   `json:"updated_at_ms"`
}

type ListBookmarksOptions struct {
	RepoFullName string
	MergeStatus  string
}

func (c *Client) ListBookmarks(ctx context.Context, opts ListBookmarksOptions) ([]CloudBookmark, error) {
	if c == nil || c.url == "" {
		return nil, nil
	}
	base := CloudBaseURL()
	if base == "" {
		return nil, fmt.Errorf("gx cloud base URL is not configured")
	}

	token, err := CloudAPIToken()
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	if strings.TrimSpace(opts.RepoFullName) != "" {
		query.Set("repo_full_name", strings.TrimSpace(opts.RepoFullName))
	}
	if strings.TrimSpace(opts.MergeStatus) != "" {
		query.Set("merge_status", strings.TrimSpace(opts.MergeStatus))
	}

	listURL := base + "/bookmarks"
	if encoded := query.Encode(); encoded != "" {
		listURL += "?" + encoded
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, listURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create bookmarks request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list gx cloud bookmarks: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("list gx cloud bookmarks: unauthorized (run `gx auth login`)")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		detail := strings.TrimSpace(string(body))
		if detail != "" {
			return nil, fmt.Errorf("list gx cloud bookmarks: status %s: %s", resp.Status, detail)
		}
		return nil, fmt.Errorf("list gx cloud bookmarks: status %s", resp.Status)
	}

	var bookmarks []CloudBookmark
	if err := json.Unmarshal(body, &bookmarks); err != nil {
		return nil, fmt.Errorf("decode gx cloud bookmarks: %w", err)
	}
	return bookmarks, nil
}

func RepoFullNameFromRemoteURL(remoteURL string) string {
	remoteURL = strings.TrimSpace(remoteURL)
	if remoteURL == "" {
		return ""
	}
	if match := strings.TrimPrefix(remoteURL, "git@github.com:"); match != remoteURL {
		parts := strings.Split(strings.TrimSuffix(match, ".git"), "/")
		if len(parts) >= 2 {
			return parts[0] + "/" + parts[1]
		}
	}
	if strings.Contains(remoteURL, "github.com") {
		idx := strings.Index(remoteURL, "github.com/")
		if idx >= 0 {
			path := strings.TrimPrefix(remoteURL[idx+len("github.com/"):], "/")
			path = strings.TrimSuffix(path, ".git")
			parts := strings.Split(path, "/")
			if len(parts) >= 2 {
				return parts[0] + "/" + parts[1]
			}
		}
	}
	return ""
}
