package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type PublishRegistration struct {
	RepoFullName  string  `json:"repo_full_name"`
	BranchName    string  `json:"branch_name"`
	Title         *string `json:"title,omitempty"`
	HeadCommitID  string  `json:"head_commit_id"`
	RemoteHeadSha *string `json:"remote_head_sha,omitempty"`
	GitHubPRURL   *string `json:"github_pr_url,omitempty"`
}

type PublishRegistrationResult struct {
	ID             string  `json:"id"`
	RepoFullName   string  `json:"repo_full_name"`
	BranchName     string  `json:"branch_name"`
	GitHubPRURL    *string `json:"github_pr_url"`
	GitHubPRNumber *int    `json:"github_pr_number"`
	HeadCommitID   *string `json:"head_commit_id"`
	RemoteHeadSha  *string `json:"remote_head_sha"`
}

func (c *Client) RegisterPublish(ctx context.Context, registration PublishRegistration) (PublishRegistrationResult, error) {
	if c == nil || c.url == "" {
		return PublishRegistrationResult{}, nil
	}
	publishURL := cloudURLWithPath(c.url, "/v1/publish")
	if publishURL == "" {
		return PublishRegistrationResult{}, fmt.Errorf("tx cloud base URL is not configured")
	}
	body, err := json.Marshal(registration)
	if err != nil {
		return PublishRegistrationResult{}, fmt.Errorf("marshal tx publish registration: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, publishURL, bytes.NewReader(body))
	if err != nil {
		return PublishRegistrationResult{}, fmt.Errorf("create tx publish registration request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	token, err := CloudAPIToken()
	if err != nil {
		return PublishRegistrationResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return PublishRegistrationResult{}, fmt.Errorf("register tx publish: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		detail := strings.TrimSpace(string(raw))
		if detail != "" {
			return PublishRegistrationResult{}, fmt.Errorf("register tx publish: status %s: %s", resp.Status, detail)
		}
		return PublishRegistrationResult{}, fmt.Errorf("register tx publish: status %s", resp.Status)
	}
	var result PublishRegistrationResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return PublishRegistrationResult{}, fmt.Errorf("decode tx publish registration: %w", err)
	}
	return result, nil
}
