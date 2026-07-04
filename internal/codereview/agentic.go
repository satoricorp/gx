package codereview

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	maxAgenticToolCalls       = 6
	maxAgenticToolOutputBytes = 6000
)

type agenticResponsesAIReviewer struct {
	base *responsesAIReviewer
}

func agenticReviewEnabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_AGENTIC")), "1")
}

func maybeAgenticReviewer(reviewer AIReviewer) AIReviewer {
	if !agenticReviewEnabled() || reviewer == nil {
		return reviewer
	}
	switch r := reviewer.(type) {
	case *responsesAIReviewer:
		return &agenticResponsesAIReviewer{base: r}
	case fallbackAIReviewer:
		return fallbackAIReviewer{
			primary:  maybeAgenticReviewer(r.primary),
			fallback: maybeAgenticReviewer(r.fallback),
		}
	default:
		return reviewer
	}
}

func (r *agenticResponsesAIReviewer) Review(ctx context.Context, brief ReviewBrief) ([]Finding, error) {
	if r == nil || r.base == nil {
		return nil, fmt.Errorf("agentic reviewer is not configured")
	}
	brief = compactReviewBriefForAI(brief)
	resolver := newSourceResolver(brief)
	input := any(mustJSON(brief))
	previousID := ""
	for step := 0; step <= maxAgenticToolCalls; step++ {
		result, err := r.createResponse(ctx, input, previousID)
		if err != nil {
			return nil, err
		}
		content := strings.TrimSpace(result.OutputText)
		if content == "" {
			content = strings.TrimSpace(responseOutputText(result))
		}
		if content != "" {
			return parseAIReviewContent(content, resolver)
		}
		calls := agenticToolCalls(result)
		if len(calls) == 0 {
			return nil, fmt.Errorf("agentic AI review response returned no findings or tool calls")
		}
		var outputs []map[string]any
		for _, call := range calls {
			outputs = append(outputs, map[string]any{
				"type":    "function_call_output",
				"call_id": call.CallID,
				"output":  executeAgenticTool(ctx, brief, call.Name, call.Arguments),
			})
		}
		input = outputs
		previousID = result.ID
	}
	return nil, fmt.Errorf("agentic AI review exceeded tool call limit")
}

func (r *agenticResponsesAIReviewer) createResponse(ctx context.Context, input any, previousID string) (responseResult, error) {
	payload := responseRequest{
		Model:        r.base.model,
		Instructions: agenticReviewDeveloperPrompt(),
		Input:        input,
		Text: responseTextConfig{
			Format: map[string]string{"type": "json_object"},
		},
		Tools:           agenticReviewTools(),
		PreviousID:      previousID,
		MaxOutputTokens: defaultReviewMaxOutputTokens,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return responseResult{}, fmt.Errorf("marshal agentic review request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.base.url, bytes.NewReader(body))
	if err != nil {
		return responseResult{}, fmt.Errorf("create agentic review request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.base.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := r.base.client
	if client == nil {
		client = &http.Client{Timeout: 120 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return responseResult{}, fmt.Errorf("request agentic AI review: %w", err)
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		detail := strings.TrimSpace(string(responseBody))
		if detail != "" {
			return responseResult{}, fmt.Errorf("agentic AI review status %s: %s", resp.Status, detail)
		}
		return responseResult{}, fmt.Errorf("agentic AI review status %s", resp.Status)
	}
	var completion responseResult
	if err := json.Unmarshal(responseBody, &completion); err != nil {
		return responseResult{}, fmt.Errorf("decode agentic AI review response: %w", err)
	}
	return completion, nil
}

type agenticToolCall struct {
	CallID    string
	Name      string
	Arguments string
}

func agenticToolCalls(response responseResult) []agenticToolCall {
	var out []agenticToolCall
	for _, item := range response.Output {
		if item.Type != "function_call" {
			continue
		}
		callID := strings.TrimSpace(item.CallID)
		if callID == "" {
			callID = strings.TrimSpace(item.ID)
		}
		if callID == "" || strings.TrimSpace(item.Name) == "" {
			continue
		}
		out = append(out, agenticToolCall{CallID: callID, Name: strings.TrimSpace(item.Name), Arguments: item.Arguments})
	}
	return out
}

func agenticReviewDeveloperPrompt() string {
	return reviewDeveloperPrompt() + "\n" + strings.Join([]string{
		"Agentic mode is enabled. Use tools only when they can verify a specific candidate issue or locate the exact anchor line.",
		"Call read_file for surrounding implementation, grep for symbol usage, list_changed_hunks for changed-line evidence, and search_knowledge for retrieved policy or historical context.",
		"After tool use, emit candidate recommendations in the exact same JSON shape as single-shot review.",
	}, "\n")
}

func agenticReviewTools() []responseTool {
	stringParam := func(name, description string) map[string]any {
		return map[string]any{"type": "string", "description": description, "title": name}
	}
	intParam := func(name, description string) map[string]any {
		return map[string]any{"type": "integer", "description": description, "title": name}
	}
	return []responseTool{
		{
			Type:        "function",
			Name:        "read_file",
			Description: "Read a repository file, optionally bounded by line numbers.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"file":       stringParam("file", "Repository-relative file path."),
					"start_line": intParam("start_line", "Optional 1-based start line."),
					"end_line":   intParam("end_line", "Optional 1-based end line."),
				},
				"required": []string{"file"},
			},
		},
		{
			Type:        "function",
			Name:        "grep",
			Description: "Search repository text files for a literal or regexp pattern.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"pattern": stringParam("pattern", "Pattern to search for."),
					"path":    stringParam("path", "Optional repository-relative path prefix."),
				},
				"required": []string{"pattern"},
			},
		},
		{
			Type:        "function",
			Name:        "list_changed_hunks",
			Description: "Return changed diff snippets from the review brief.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"file": stringParam("file", "Optional repository-relative file path."),
				},
			},
		},
		{
			Type:        "function",
			Name:        "search_knowledge",
			Description: "Search retrieved review knowledge, policies, references, and local context snippets.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": stringParam("query", "Words to search for in retrieved context."),
				},
				"required": []string{"query"},
			},
		},
	}
}

func executeAgenticTool(ctx context.Context, brief ReviewBrief, name string, rawArgs string) string {
	var args map[string]any
	_ = json.Unmarshal([]byte(rawArgs), &args)
	switch name {
	case "read_file":
		return truncateReviewText(agenticReadFile(brief.RepoRoot, stringArg(args, "file"), intArg(args, "start_line"), intArg(args, "end_line")), maxAgenticToolOutputBytes)
	case "grep":
		return truncateReviewText(agenticGrep(ctx, brief.RepoRoot, stringArg(args, "pattern"), stringArg(args, "path")), maxAgenticToolOutputBytes)
	case "list_changed_hunks":
		return truncateReviewText(agenticListChangedHunks(brief, stringArg(args, "file")), maxAgenticToolOutputBytes)
	case "search_knowledge":
		return truncateReviewText(agenticSearchKnowledge(brief, stringArg(args, "query")), maxAgenticToolOutputBytes)
	default:
		return "unknown tool: " + name
	}
}

func agenticReadFile(repoRoot, file string, startLine, endLine int) string {
	file = normalizeAnchorFile(file)
	if file == "." || strings.Contains(file, "..") {
		return "invalid file path"
	}
	data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(file)))
	if err != nil {
		return "read_file error: " + err.Error()
	}
	lines := strings.Split(string(data), "\n")
	if startLine <= 0 {
		startLine = 1
	}
	if endLine <= 0 || endLine > len(lines) {
		endLine = len(lines)
	}
	if startLine > endLine || startLine > len(lines) {
		return ""
	}
	var b strings.Builder
	for i := startLine; i <= endLine; i++ {
		fmt.Fprintf(&b, "%s:%d: %s\n", file, i, lines[i-1])
	}
	return b.String()
}

func agenticGrep(ctx context.Context, repoRoot, pattern, pathPrefix string) string {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return "grep pattern is required"
	}
	pathPrefix = strings.Trim(strings.TrimSpace(filepath.ToSlash(pathPrefix)), "/")
	re, reErr := regexp.Compile(pattern)
	var b strings.Builder
	matches := 0
	err := filepath.WalkDir(repoRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil || ctx.Err() != nil {
			return err
		}
		rel, relErr := filepath.Rel(repoRoot, path)
		if relErr != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if entry.IsDir() {
			if rel != "." && skipDir(rel) {
				return filepath.SkipDir
			}
			return nil
		}
		if skipFile(rel) || (pathPrefix != "" && rel != pathPrefix && !strings.HasPrefix(rel, pathPrefix+"/")) {
			return nil
		}
		file, openErr := os.Open(path)
		if openErr != nil {
			return nil
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		lineNo := 0
		for scanner.Scan() {
			lineNo++
			line := scanner.Text()
			matched := strings.Contains(line, pattern)
			if reErr == nil {
				matched = re.MatchString(line)
			}
			if !matched {
				continue
			}
			fmt.Fprintf(&b, "%s:%d: %s\n", rel, lineNo, strings.TrimSpace(line))
			matches++
			if matches >= 50 {
				return filepath.SkipAll
			}
		}
		return nil
	})
	if err != nil && err != filepath.SkipAll {
		return "grep error: " + err.Error()
	}
	return b.String()
}

func agenticListChangedHunks(brief ReviewBrief, file string) string {
	file = normalizeAnchorFile(file)
	var b strings.Builder
	for _, snippet := range brief.Static.DiffSnippets {
		if file != "." && file != "" && normalizeAnchorFile(snippet.File) != file {
			continue
		}
		fmt.Fprintf(&b, "file: %s\n%s\n", snippet.File, snippet.Diff)
	}
	return b.String()
}

func agenticSearchKnowledge(brief ReviewBrief, query string) string {
	terms := strings.Fields(strings.ToLower(query))
	if len(terms) == 0 {
		return "search_knowledge query is required"
	}
	var b strings.Builder
	matches := 0
	for _, snippet := range brief.Context {
		text := strings.ToLower(snippet.Text)
		ok := true
		for _, term := range terms {
			if !strings.Contains(text, term) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		label := firstNonEmpty(snippet.SourceLabel, snippet.Ref, snippet.Kind)
		fmt.Fprintf(&b, "[%s] %s\n%s\n", label, snippet.Ref, snippet.Text)
		matches++
		if matches >= 5 {
			break
		}
	}
	return b.String()
}

func stringArg(args map[string]any, key string) string {
	if value, ok := args[key].(string); ok {
		return strings.TrimSpace(value)
	}
	return ""
}

func intArg(args map[string]any, key string) int {
	switch value := args[key].(type) {
	case float64:
		return int(value)
	case int:
		return value
	default:
		return 0
	}
}
