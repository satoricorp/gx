package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/authstore"
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

type PullRequest struct {
	Number int
	URL    string
}

type IssueComment struct {
	ID   int64
	URL  string
	Body string
}

type CreatePullRequestOptions struct {
	Host       string
	Owner      string
	Repo       string
	BaseBranch string
	HeadBranch string
	Title      string
	Body       string
}

type IssueCommentOptions struct {
	Owner  string
	Repo   string
	Number int
	Body   string
	Marker string
}

func NewClient(host string) (*Client, error) {
	token, err := accessToken()
	if err != nil {
		return nil, err
	}
	return NewClientWithToken(host, token, nil), nil
}

func accessToken() (string, error) {
	return authstore.GitHubAccessToken()
}

func tokenError() error {
	if strings.TrimSpace(os.Getenv("GX_MCP")) != "" {
		return fmt.Errorf("github token is not configured for MCP: run `gx auth login` in a terminal, then retry the MCP tool")
	}
	return fmt.Errorf("github token is not configured: run `gx auth login` or set GH_TOKEN/GITHUB_TOKEN")
}

func NewClientWithToken(host, token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{
		baseURL: apiBaseURL(host),
		token:   strings.TrimSpace(token),
		http:    httpClient,
	}
}

func apiBaseURL(host string) string {
	if override := strings.TrimSpace(os.Getenv("GX_GITHUB_API_URL")); override != "" {
		return strings.TrimRight(override, "/")
	}
	host = strings.TrimSpace(host)
	if host == "" || strings.EqualFold(host, "github.com") {
		return "https://api.github.com"
	}
	return "https://" + strings.TrimRight(host, "/") + "/api/v3"
}

func (c *Client) FindPullRequest(ctx context.Context, opts CreatePullRequestOptions) (*PullRequest, error) {
	if c == nil {
		return nil, fmt.Errorf("github client is required")
	}
	q := url.Values{}
	q.Set("state", "open")
	q.Set("head", opts.Owner+":"+opts.HeadBranch)
	endpoint := fmt.Sprintf("/repos/%s/%s/pulls?%s", url.PathEscape(opts.Owner), url.PathEscape(opts.Repo), q.Encode())
	req, err := c.request(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	var payload []struct {
		Number  int    `json:"number"`
		HTMLURL string `json:"html_url"`
	}
	if err := c.do(req, &payload); err != nil {
		return nil, fmt.Errorf("find github pull request: %w", err)
	}
	if len(payload) == 0 || strings.TrimSpace(payload[0].HTMLURL) == "" {
		return nil, nil
	}
	return &PullRequest{Number: payload[0].Number, URL: strings.TrimSpace(payload[0].HTMLURL)}, nil
}

func (c *Client) CreatePullRequest(ctx context.Context, opts CreatePullRequestOptions) (*PullRequest, error) {
	if c == nil {
		return nil, fmt.Errorf("github client is required")
	}
	body, err := json.Marshal(map[string]string{
		"base":  strings.TrimSpace(opts.BaseBranch),
		"head":  strings.TrimSpace(opts.HeadBranch),
		"title": strings.TrimSpace(opts.Title),
		"body":  opts.Body,
	})
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("/repos/%s/%s/pulls", url.PathEscape(opts.Owner), url.PathEscape(opts.Repo))
	req, err := c.request(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var payload struct {
		HTMLURL string `json:"html_url"`
	}
	if err := c.do(req, &payload); err != nil {
		return nil, fmt.Errorf("create github pull request: %w", err)
	}
	if strings.TrimSpace(payload.HTMLURL) == "" {
		return nil, nil
	}
	return &PullRequest{URL: strings.TrimSpace(payload.HTMLURL)}, nil
}

func (c *Client) ListIssueComments(ctx context.Context, opts IssueCommentOptions) ([]IssueComment, error) {
	if c == nil {
		return nil, fmt.Errorf("github client is required")
	}
	endpoint := fmt.Sprintf("/repos/%s/%s/issues/%d/comments?per_page=100", url.PathEscape(opts.Owner), url.PathEscape(opts.Repo), opts.Number)
	req, err := c.request(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	var payload []struct {
		ID      int64  `json:"id"`
		HTMLURL string `json:"html_url"`
		Body    string `json:"body"`
	}
	if err := c.do(req, &payload); err != nil {
		return nil, fmt.Errorf("list github issue comments: %w", err)
	}
	comments := make([]IssueComment, 0, len(payload))
	for _, item := range payload {
		comments = append(comments, IssueComment{
			ID:   item.ID,
			URL:  strings.TrimSpace(item.HTMLURL),
			Body: item.Body,
		})
	}
	return comments, nil
}

func (c *Client) CreateIssueComment(ctx context.Context, opts IssueCommentOptions) (*IssueComment, error) {
	if c == nil {
		return nil, fmt.Errorf("github client is required")
	}
	body, err := json.Marshal(map[string]string{"body": opts.Body})
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("/repos/%s/%s/issues/%d/comments", url.PathEscape(opts.Owner), url.PathEscape(opts.Repo), opts.Number)
	req, err := c.request(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var payload struct {
		ID      int64  `json:"id"`
		HTMLURL string `json:"html_url"`
		Body    string `json:"body"`
	}
	if err := c.do(req, &payload); err != nil {
		return nil, fmt.Errorf("create github issue comment: %w", err)
	}
	return &IssueComment{ID: payload.ID, URL: strings.TrimSpace(payload.HTMLURL), Body: payload.Body}, nil
}

func (c *Client) UpdateIssueComment(ctx context.Context, opts IssueCommentOptions, commentID int64) (*IssueComment, error) {
	if c == nil {
		return nil, fmt.Errorf("github client is required")
	}
	body, err := json.Marshal(map[string]string{"body": opts.Body})
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("/repos/%s/%s/issues/comments/%d", url.PathEscape(opts.Owner), url.PathEscape(opts.Repo), commentID)
	req, err := c.request(ctx, http.MethodPatch, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var payload struct {
		ID      int64  `json:"id"`
		HTMLURL string `json:"html_url"`
		Body    string `json:"body"`
	}
	if err := c.do(req, &payload); err != nil {
		return nil, fmt.Errorf("update github issue comment: %w", err)
	}
	return &IssueComment{ID: payload.ID, URL: strings.TrimSpace(payload.HTMLURL), Body: payload.Body}, nil
}

func (c *Client) UpsertIssueComment(ctx context.Context, opts IssueCommentOptions) (*IssueComment, error) {
	marker := strings.TrimSpace(opts.Marker)
	if marker == "" {
		return c.CreateIssueComment(ctx, opts)
	}
	comments, err := c.ListIssueComments(ctx, opts)
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		if strings.Contains(comment.Body, marker) {
			return c.UpdateIssueComment(ctx, opts, comment.ID)
		}
	}
	return c.CreateIssueComment(ctx, opts)
}

func (c *Client) request(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
	if strings.TrimSpace(c.token) == "" {
		return nil, tokenError()
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("github auth failed: status %s", resp.Status)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		detail := strings.TrimSpace(string(body))
		if detail != "" {
			return fmt.Errorf("status %s: %s", resp.Status, detail)
		}
		return fmt.Errorf("status %s", resp.Status)
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode github response: %w", err)
	}
	return nil
}
