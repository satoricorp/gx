package review

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/auth"
	"github.com/satoricorp/gx/internal/version"
)

// RemoteContext is the WP-5g review context payload.
type RemoteContext struct {
	Rules         []RemoteRule      `json:"rules"`
	Collisions    []RemoteCollision `json:"collisions"`
	IndexSnippets []string          `json:"index_snippets"`
}

type RemoteRule struct {
	ID       string `json:"id"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	File     string `json:"file,omitempty"`
}

type RemoteCollision struct {
	File      string `json:"file"`
	SessionID string `json:"session_id,omitempty"`
	Detail    string `json:"detail"`
}

type serverReviewContext struct {
	Rules []struct {
		ID        string  `json:"id"`
		RepoScope *string `json:"repoScope"`
		RuleText  string  `json:"ruleText"`
		ScopeExpr *string `json:"scopeExpr"`
		Strength  string  `json:"strength"`
		Status    string  `json:"status"`
	} `json:"rules"`
	Collisions []struct {
		File      string `json:"file"`
		LineStart *int   `json:"lineStart"`
		LineEnd   *int   `json:"lineEnd"`
		Kind      string `json:"kind"`
		Detail    string `json:"detail"`
	} `json:"collisions"`
	IndexSnippets []struct {
		ID   string  `json:"id"`
		Text string  `json:"text"`
		Score *float64 `json:"score"`
	} `json:"indexSnippets"`
}

type remoteClient struct {
	baseURL string
	token   string
	client  *http.Client
}

func newRemoteClient() (*remoteClient, bool) {
	creds, ok := auth.Load()
	if !ok || strings.TrimSpace(creds.Token) == "" {
		return nil, false
	}
	base := strings.TrimRight(strings.TrimSpace(creds.APIURL), "/")
	if base == "" {
		return nil, false
	}
	return &remoteClient{
		baseURL: base,
		token:   creds.Token,
		client:  &http.Client{Timeout: 10 * time.Second},
	}, true
}

func (c *remoteClient) fetchContext(ctx context.Context, repoRoot, base, head string) (RemoteContext, error) {
	query := url.Values{}
	query.Set("repoRoot", repoRoot)
	query.Set("base", base)
	query.Set("head", head)
	endpoint := c.baseURL + "/v1/review/context?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return RemoteContext{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "gx/"+version.Current())
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		return RemoteContext{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return RemoteContext{}, fmt.Errorf("review context status %s", resp.Status)
	}
	var parsed serverReviewContext
	if err := json.Unmarshal(body, &parsed); err != nil {
		return RemoteContext{}, fmt.Errorf("decode review context: %w", err)
	}
	return mapServerReviewContext(parsed), nil
}

func mapServerReviewContext(parsed serverReviewContext) RemoteContext {
	out := RemoteContext{
		Rules:         make([]RemoteRule, 0, len(parsed.Rules)),
		Collisions:    make([]RemoteCollision, 0, len(parsed.Collisions)),
		IndexSnippets: make([]string, 0, len(parsed.IndexSnippets)),
	}
	for _, rule := range parsed.Rules {
		file := ""
		if rule.ScopeExpr != nil {
			file = strings.TrimSpace(*rule.ScopeExpr)
		}
		if file == "" && rule.RepoScope != nil {
			file = strings.TrimSpace(*rule.RepoScope)
		}
		out.Rules = append(out.Rules, RemoteRule{
			ID:       rule.ID,
			Severity: rule.Strength,
			Message:  rule.RuleText,
			File:     file,
		})
	}
	for _, collision := range parsed.Collisions {
		out.Collisions = append(out.Collisions, RemoteCollision{
			File:   collision.File,
			Detail: collision.Detail,
		})
	}
	for _, snippet := range parsed.IndexSnippets {
		text := strings.TrimSpace(snippet.Text)
		if text != "" {
			out.IndexSnippets = append(out.IndexSnippets, text)
		}
	}
	return out
}

func gradeRulesSection(remote RemoteContext, offline bool) Section {
	if offline {
		return Section{
			Name:    SectionRules,
			Grade:   GradeSkip,
			Summary: "Server rules unavailable (offline).",
		}
	}
	violations := 0
	var evidence []string
	for _, rule := range remote.Rules {
		if strings.EqualFold(rule.Severity, "binding") || strings.EqualFold(rule.Severity, "fail") {
			violations++
			item := rule.Message
			if rule.File != "" {
				item = rule.File + ": " + item
			}
			evidence = append(evidence, item)
		}
	}
	section := Section{Name: SectionRules, Grade: GradePass, Summary: "No binding rule violations.", Evidence: evidence}
	if violations > 0 {
		section.Grade = GradeFail
		section.Summary = fmt.Sprintf("%d binding rule violation(s).", violations)
	}
	return section
}

func gradeCollisionsSection(remote RemoteContext, offline bool) Section {
	if offline {
		return Section{
			Name:    SectionCollisions,
			Grade:   GradeSkip,
			Summary: "Server collision data unavailable (offline).",
		}
	}
	if len(remote.Collisions) == 0 {
		return Section{
			Name:    SectionCollisions,
			Grade:   GradePass,
			Summary: "No overlapping agent/human edit collisions reported.",
		}
	}
	section := Section{
		Name:    SectionCollisions,
		Grade:   GradeWarn,
		Summary: fmt.Sprintf("%d collision(s) reported.", len(remote.Collisions)),
	}
	for _, collision := range remote.Collisions {
		item := collision.File
		if collision.Detail != "" {
			item += ": " + collision.Detail
		}
		section.Evidence = append(section.Evidence, item)
	}
	return section
}
