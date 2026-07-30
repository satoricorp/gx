package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/satoricorp/totality/internal/version"
)

type CodeReviewFindingRecord struct {
	Fingerprint    string         `json:"fingerprint,omitempty"`
	Outcome        string         `json:"outcome,omitempty"`
	Category       string         `json:"category,omitempty"`
	Language       string         `json:"language,omitempty"`
	FilePath       string         `json:"filePath,omitempty"`
	LineStart      int            `json:"lineStart,omitempty"`
	LineEnd        int            `json:"lineEnd,omitempty"`
	Title          string         `json:"title"`
	Summary        string         `json:"summary"`
	Recommendation string         `json:"recommendation,omitempty"`
	Confidence     int            `json:"confidence,omitempty"`
	Severity       string         `json:"severity,omitempty"`
	Payload        map[string]any `json:"payload,omitempty"`
}

type CodeReviewHistoryRecordRequest struct {
	RepoRootPath string                    `json:"repoRootPath,omitempty"`
	RepoFullName string                    `json:"repoFullName"`
	BranchName   string                    `json:"branchName,omitempty"`
	HeadCommitID string                    `json:"headCommitId,omitempty"`
	SourceKind   string                    `json:"sourceKind,omitempty"`
	SourceRef    string                    `json:"sourceRef,omitempty"`
	Prompt       string                    `json:"prompt,omitempty"`
	Scope        string                    `json:"scope,omitempty"`
	Mode         string                    `json:"mode,omitempty"`
	Reviewer     string                    `json:"reviewer,omitempty"`
	SummaryText  string                    `json:"summaryText,omitempty"`
	SummaryKind  string                    `json:"summaryKind,omitempty"`
	Findings     []CodeReviewFindingRecord `json:"findings,omitempty"`
	Payload      map[string]any            `json:"payload,omitempty"`
}

type CodeReviewHistoryRecordResult struct {
	RunID        string `json:"runId"`
	SummaryID    string `json:"summaryId"`
	FindingCount int    `json:"findingCount"`
	IndexStatus  string `json:"indexStatus"`
	IndexError   string `json:"indexError,omitempty"`
}

type CodeReviewHistorySearchRequest struct {
	RepoFullName string   `json:"repoFullName"`
	Query        string   `json:"query,omitempty"`
	Fingerprints []string `json:"fingerprints,omitempty"`
	FilePaths    []string `json:"filePaths,omitempty"`
	Categories   []string `json:"categories,omitempty"`
	Language     string   `json:"language,omitempty"`
	Limit        int      `json:"limit,omitempty"`
}

type CodeReviewHistoryFindingMatch struct {
	ID             string         `json:"id"`
	RunID          string         `json:"run_id"`
	Fingerprint    string         `json:"fingerprint"`
	Outcome        string         `json:"outcome"`
	Category       string         `json:"category"`
	Language       string         `json:"language"`
	FilePath       string         `json:"file_path"`
	LineStart      int            `json:"line_start"`
	LineEnd        int            `json:"line_end"`
	Title          string         `json:"title"`
	Summary        string         `json:"summary"`
	Recommendation string         `json:"recommendation"`
	Confidence     int            `json:"confidence"`
	Severity       string         `json:"severity"`
	CreatedAtMS    int64          `json:"created_at_ms"`
	Payload        map[string]any `json:"payload"`
}

type CodeReviewHistorySimilarMatch struct {
	ID         string         `json:"id"`
	Score      float64        `json:"score"`
	Text       string         `json:"text"`
	Attributes map[string]any `json:"attributes"`
}

type CodeReviewHistorySummaryMatch struct {
	ID          string         `json:"id"`
	RunID       string         `json:"run_id"`
	SummaryKind string         `json:"summary_kind"`
	SummaryText string         `json:"summary_text"`
	SummaryJSON map[string]any `json:"summary_json"`
	CreatedAtMS int64          `json:"created_at_ms"`
}

type CodeReviewHistorySearchResult struct {
	ExactMatches   []CodeReviewHistoryFindingMatch `json:"exactMatches"`
	SimilarMatches []CodeReviewHistorySimilarMatch `json:"similarMatches"`
	Summaries      []CodeReviewHistorySummaryMatch `json:"summaries"`
}

func (c *Client) RecordCodeReviewHistory(ctx context.Context, reqBody CodeReviewHistoryRecordRequest) (CodeReviewHistoryRecordResult, error) {
	var result CodeReviewHistoryRecordResult
	if c == nil || c.url == "" || strings.TrimSpace(reqBody.RepoFullName) == "" {
		return result, nil
	}
	if err := c.postJSON(ctx, "/v1/code-review-history/record", reqBody, &result); err != nil {
		return CodeReviewHistoryRecordResult{}, err
	}
	return result, nil
}

func (c *Client) SearchCodeReviewHistory(ctx context.Context, reqBody CodeReviewHistorySearchRequest) (CodeReviewHistorySearchResult, error) {
	var result CodeReviewHistorySearchResult
	if c == nil || c.url == "" || strings.TrimSpace(reqBody.RepoFullName) == "" {
		return result, nil
	}
	if err := c.postJSON(ctx, "/v1/code-review-history/search", reqBody, &result); err != nil {
		return CodeReviewHistorySearchResult{}, err
	}
	return result, nil
}

func (c *Client) postJSON(ctx context.Context, path string, payload any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal tx cloud request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cloudURLWithPath(c.url, path), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create tx cloud request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "tx/"+version.Current())
	token, err := CloudAPIToken()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("send tx cloud request: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		detail := strings.TrimSpace(string(raw))
		if detail != "" {
			return fmt.Errorf("tx cloud request %s: status %s: %s", path, resp.Status, detail)
		}
		return fmt.Errorf("tx cloud request %s: status %s", path, resp.Status)
	}
	if out != nil && len(raw) > 0 {
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("decode tx cloud response: %w", err)
		}
	}
	return nil
}
