package postlist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	envURL    = "LGTM_POSTLIST_URL"
	envAPIKey = "LGTM_POSTLIST_API_KEY"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type Identity struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Source    string `json:"source"`
	TLVersion string `json:"lgtm_version,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func NewFromEnv() *Client {
	baseURL := strings.TrimSpace(os.Getenv(envURL))
	if baseURL == "" {
		return nil
	}
	return &Client{
		baseURL:    baseURL,
		apiKey:     strings.TrimSpace(os.Getenv(envAPIKey)),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) UpsertIdentity(ctx context.Context, identity Identity) error {
	if c == nil {
		return nil
	}
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("parse PostList URL: %w", err)
	}
	query := u.Query()
	if query.Get("on_conflict") == "" {
		query.Set("on_conflict", "email")
	}
	u.RawQuery = query.Encode()

	if identity.Source == "" {
		identity.Source = "lgtm"
	}
	if identity.UpdatedAt == "" {
		identity.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	body, err := json.Marshal([]Identity{identity})
	if err != nil {
		return fmt.Errorf("marshal PostList payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("prefer", "resolution=merge-duplicates,return=minimal")
	if c.apiKey != "" {
		req.Header.Set("authorization", "Bearer "+c.apiKey)
		req.Header.Set("apikey", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("PostList request failed: %s", resp.Status)
	}
	return nil
}
