package codereview

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/version"
)

const (
	defaultReviewModel           = "gpt-4.1-mini"
	defaultReviewMaxOutputTokens = 6000
	defaultOpenAIBaseURL         = "https://api.openai.com/v1"
	maxAIContextSnippetBytes     = 1200
	maxAIStaticToolOutputBytes   = 4000
	maxAIContextSnippets         = 14
	maxAIDeepContextSnippets     = 40
	maxAICodeQualityHints        = 20
	maxAIModuleSummaries         = 12
	maxAIChangedFiles            = 60
	bedrockReviewModel           = "anthropic.claude-sonnet-4-6"
)

type AIReviewer interface {
	Review(ctx context.Context, brief ReviewBrief) ([]Finding, error)
}

type responsesAIReviewer struct {
	url    string
	token  string
	model  string
	client *http.Client
}

type namedAIReviewer struct {
	name     string
	label    string
	reviewer AIReviewer
}

type multiAIReviewer struct {
	reviewers []namedAIReviewer
}

type fallbackAIReviewer struct {
	primary  AIReviewer
	fallback AIReviewer
}

type bedrockAnthropicReviewer struct {
	region    string
	model     string
	accessKey string
	secretKey string
	client    *http.Client
}

type responseRequest struct {
	Model           string             `json:"model"`
	Instructions    string             `json:"instructions"`
	Input           string             `json:"input"`
	Text            responseTextConfig `json:"text,omitempty"`
	MaxOutputTokens int                `json:"max_output_tokens,omitempty"`
}

type responseTextConfig struct {
	Format map[string]string `json:"format,omitempty"`
}

type responseResult struct {
	Model      string `json:"model"`
	OutputText string `json:"output_text"`
	Output     []struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

type aiReviewResponse struct {
	Recommendations []aiRecommendation `json:"recommendations"`
}

type aiRecommendation struct {
	Title          string   `json:"title"`
	Summary        string   `json:"summary"`
	Benefit        string   `json:"benefit"`
	Recommendation string   `json:"recommendation"`
	Strength       string   `json:"strength"`
	Evidence       []string `json:"evidence"`
}

func reviewerFromEnv() AIReviewer {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_AI")), "0") {
		return nil
	}
	var reviewers []namedAIReviewer
	if reviewer := openAIReviewerFromEnv(); reviewer != nil {
		reviewers = append(reviewers, namedAIReviewer{name: "openai", label: "OpenAI", reviewer: reviewer})
	}
	if reviewer := bedrockAnthropicReviewerFromEnv(); reviewer != nil {
		reviewers = append(reviewers, namedAIReviewer{name: "anthropic", label: "Anthropic", reviewer: reviewer})
	}
	switch len(reviewers) {
	case 0:
		return nil
	case 1:
		return reviewers[0].reviewer
	default:
		return multiAIReviewer{reviewers: reviewers}
	}
}

func openAIReviewerFromEnv() AIReviewer {
	model := strings.TrimSpace(firstNonEmpty(
		os.Getenv("GX_REVIEW_OPENAI_MODEL"),
		os.Getenv("GX_REVIEW_MODEL"),
		os.Getenv("OPENAI_MODEL"),
		defaultReviewModel,
	))
	direct, directErr := directOpenAIReviewerFromEnv(model)
	cloudReviewer := cloudOpenAIReviewerFromEnv(model)
	if direct != nil {
		if cloudReviewer != nil {
			return fallbackAIReviewer{primary: direct, fallback: cloudReviewer}
		}
		return direct
	}
	if directErr != nil {
		return nil
	}
	return cloudReviewer
}

func cloudOpenAIReviewerFromEnv(model string) AIReviewer {
	if url := strings.TrimSpace(os.Getenv("GX_OPENAI_PROXY_URL")); url != "" {
		if token, err := cloud.CloudAPIToken(); err == nil {
			return &responsesAIReviewer{url: url, token: token, model: model, client: &http.Client{Timeout: 120 * time.Second}}
		}
	}
	if baseURL := cloud.CloudBaseURL(); baseURL != "" {
		if token, err := cloud.CloudAPIToken(); err == nil {
			return &responsesAIReviewer{
				url:    strings.TrimRight(baseURL, "/") + "/gx/openai/responses",
				token:  token,
				model:  model,
				client: &http.Client{Timeout: 120 * time.Second},
			}
		}
	}
	return nil
}

func directOpenAIReviewerFromEnv(model string) (AIReviewer, error) {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	baseURLRaw := strings.TrimSpace(os.Getenv("OPENAI_BASE_URL"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("GX_OPENAI_API_KEY"))
		baseURLRaw = strings.TrimSpace(firstNonEmpty(os.Getenv("GX_OPENAI_BASE_URL"), os.Getenv("OPENAI_BASE_URL")))
	}
	if apiKey == "" {
		return nil, nil
	}
	baseURL := strings.TrimRight(strings.TrimSpace(firstNonEmpty(baseURLRaw, defaultOpenAIBaseURL)), "/")
	if !strings.HasSuffix(baseURL, "/v1") {
		baseURL += "/v1"
	}
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("OpenAI base URL is invalid: %w", err)
	}
	return &responsesAIReviewer{
		url:    baseURL + "/responses",
		token:  apiKey,
		model:  model,
		client: &http.Client{Timeout: 120 * time.Second},
	}, nil
}

func bedrockAnthropicReviewerFromEnv() AIReviewer {
	accessKey := strings.TrimSpace(os.Getenv("AWS_ACCESS_KEY_ID"))
	secretKey := strings.TrimSpace(os.Getenv("AWS_SECRET_ACCESS_KEY"))
	if accessKey == "" || secretKey == "" {
		return nil
	}
	region := strings.TrimSpace(firstNonEmpty(
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_DEFAULT_REGION"),
		"us-east-1",
	))
	return &bedrockAnthropicReviewer{
		region:    region,
		model:     bedrockReviewModel,
		accessKey: accessKey,
		secretKey: secretKey,
		client:    &http.Client{Timeout: 120 * time.Second},
	}
}

func (m multiAIReviewer) Review(ctx context.Context, brief ReviewBrief) ([]Finding, error) {
	var out []Finding
	var errors []string
	seen := map[string]struct{}{}
	for _, item := range m.reviewers {
		findings, err := item.reviewer.Review(ctx, brief)
		if err != nil {
			errors = append(errors, item.label+": "+err.Error())
			continue
		}
		for _, finding := range findings {
			key := strings.ToLower(strings.TrimSpace(finding.Title + "\x00" + finding.Summary))
			if key == "\x00" {
				continue
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			finding.ID = item.name + "." + finding.ID
			finding.Evidence = append([]Evidence{{Label: "Reviewer", Value: item.label}}, finding.Evidence...)
			out = append(out, finding)
			if len(out) >= 6 {
				return out, nil
			}
		}
	}
	if len(out) > 0 {
		return out, nil
	}
	if len(errors) > 0 {
		return nil, fmt.Errorf("AI reviewers failed: %s", strings.Join(errors, "; "))
	}
	return nil, fmt.Errorf("AI reviewers returned no findings")
}

func (r fallbackAIReviewer) Review(ctx context.Context, brief ReviewBrief) ([]Finding, error) {
	findings, err := r.primary.Review(ctx, brief)
	if err == nil && len(findings) > 0 {
		return findings, nil
	}
	fallbackFindings, fallbackErr := r.fallback.Review(ctx, brief)
	if fallbackErr == nil && len(fallbackFindings) > 0 {
		return fallbackFindings, nil
	}
	if err != nil && fallbackErr != nil {
		return nil, fmt.Errorf("%w; fallback reviewer failed: %v", err, fallbackErr)
	}
	if err != nil {
		return nil, err
	}
	if fallbackErr != nil {
		return nil, fallbackErr
	}
	return nil, fmt.Errorf("AI reviewers returned no findings")
}

func (r *responsesAIReviewer) Review(ctx context.Context, brief ReviewBrief) ([]Finding, error) {
	brief = compactReviewBriefForAI(brief)
	payload := responseRequest{
		Model:        r.model,
		Instructions: reviewDeveloperPrompt(),
		Input:        mustJSON(brief),
		Text: responseTextConfig{
			Format: map[string]string{"type": "json_object"},
		},
		MaxOutputTokens: defaultReviewMaxOutputTokens,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal review request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create review request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "gx/"+version.Current())

	client := r.client
	if client == nil {
		client = &http.Client{Timeout: 120 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request AI review: %w", err)
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		detail := strings.TrimSpace(string(responseBody))
		if detail != "" {
			return nil, fmt.Errorf("AI review status %s: %s", resp.Status, detail)
		}
		return nil, fmt.Errorf("AI review status %s", resp.Status)
	}
	var completion responseResult
	if err := json.Unmarshal(responseBody, &completion); err != nil {
		return nil, fmt.Errorf("decode AI review response: %w", err)
	}
	content := strings.TrimSpace(completion.OutputText)
	if content == "" {
		content = strings.TrimSpace(responseOutputText(completion))
	}
	if content == "" {
		return nil, fmt.Errorf("AI review response returned empty output")
	}
	return parseAIReviewContent(content)
}

func (r *bedrockAnthropicReviewer) Review(ctx context.Context, brief ReviewBrief) ([]Finding, error) {
	brief = compactReviewBriefForAI(brief)
	body, err := json.Marshal(map[string]any{
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens":        defaultReviewMaxOutputTokens,
		"system":            reviewDeveloperPrompt(),
		"messages": []map[string]string{{
			"role":    "user",
			"content": mustJSON(brief),
		}},
	})
	if err != nil {
		return nil, fmt.Errorf("marshal Bedrock review request: %w", err)
	}

	requestPath := "/model/" + url.PathEscape(r.model) + "/invoke"
	endpoint := "https://bedrock-runtime." + r.region + ".amazonaws.com" + requestPath
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create Bedrock review request: %w", err)
	}
	r.signBedrockRequest(req, body, requestPath)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "gx/"+version.Current())

	client := r.client
	if client == nil {
		client = &http.Client{Timeout: 120 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request Bedrock review: %w", err)
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		detail := strings.TrimSpace(string(responseBody))
		if detail != "" {
			return nil, fmt.Errorf("Bedrock review status %s: %s", resp.Status, detail)
		}
		return nil, fmt.Errorf("Bedrock review status %s", resp.Status)
	}

	var completion struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(responseBody, &completion); err != nil {
		return nil, fmt.Errorf("decode Bedrock review response: %w", err)
	}
	var content strings.Builder
	for _, part := range completion.Content {
		if part.Type == "text" {
			content.WriteString(part.Text)
		}
	}
	text := strings.TrimSpace(content.String())
	if text == "" {
		return nil, fmt.Errorf("Bedrock review response returned empty output")
	}
	return parseAIReviewContent(text)
}

func (r *bedrockAnthropicReviewer) signBedrockRequest(req *http.Request, body []byte, requestPath string) {
	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	host := req.URL.Host
	service := "bedrock"
	canonicalPath := strings.ReplaceAll(requestPath, "%", "%25")
	canonicalHeaders := strings.Join([]string{
		"accept:application/json",
		"content-type:application/json",
		"host:" + host,
		"x-amz-date:" + amzDate,
		"",
	}, "\n")
	signedHeaders := "accept;content-type;host;x-amz-date"
	canonicalRequest := strings.Join([]string{
		http.MethodPost,
		canonicalPath,
		"",
		canonicalHeaders,
		signedHeaders,
		sha256Hex(body),
	}, "\n")
	scope := dateStamp + "/" + r.region + "/" + service + "/aws4_request"
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		scope,
		sha256Hex([]byte(canonicalRequest)),
	}, "\n")
	signature := hex.EncodeToString(hmacSHA256(awsSigningKey(r.secretKey, dateStamp, r.region, service), []byte(stringToSign)))
	authorization := strings.Join([]string{
		"AWS4-HMAC-SHA256 Credential=" + r.accessKey + "/" + scope,
		"SignedHeaders=" + signedHeaders,
		"Signature=" + signature,
	}, ", ")

	req.Header.Set("Authorization", authorization)
	req.Header.Set("X-Amz-Date", amzDate)
}

func parseAIReviewContent(content string) ([]Finding, error) {
	var parsed aiReviewResponse
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("decode AI review JSON: %w", err)
	}
	return aiRecommendationsToFindings(parsed.Recommendations), nil
}

func sha256Hex(value []byte) string {
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}

func hmacSHA256(key []byte, value []byte) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(value)
	return mac.Sum(nil)
}

func awsSigningKey(secretKey, dateStamp, region, service string) []byte {
	dateKey := hmacSHA256([]byte("AWS4"+secretKey), []byte(dateStamp))
	regionKey := hmacSHA256(dateKey, []byte(region))
	serviceKey := hmacSHA256(regionKey, []byte(service))
	return hmacSHA256(serviceKey, []byte("aws4_request"))
}

func responseOutputText(response responseResult) string {
	var out strings.Builder
	for _, item := range response.Output {
		for _, content := range item.Content {
			if content.Type == "output_text" {
				out.WriteString(content.Text)
			}
		}
	}
	return out.String()
}

func compactReviewBriefForAI(brief ReviewBrief) ReviewBrief {
	deep := strings.EqualFold(brief.Depth, "deep")
	contextLimit := maxAIContextSnippets
	if deep {
		contextLimit = maxAIDeepContextSnippets
	}

	brief.Static.DependencyFiles = limitStrings(brief.Static.DependencyFiles, maxAIChangedFiles)
	brief.Static.ChangedFiles = limitStrings(brief.Static.ChangedFiles, maxAIChangedFiles)
	brief.Static.DiffSnippets = compactDiffSnippets(brief.Static.DiffSnippets, contextLimit)
	brief.Static.Modules = limitModules(brief.Static.Modules, maxAIModuleSummaries)
	brief.Static.ToolResults = compactStaticToolResults(brief.Static.ToolResults)
	brief.Static.CodeQuality = limitCodeQualityHints(brief.Static.CodeQuality, maxAICodeQualityHints)
	brief.Context = compactContextSnippets(brief.Context, contextLimit)
	return brief
}

func compactDiffSnippets(snippets []DiffSnippet, limit int) []DiffSnippet {
	if limit > 0 && len(snippets) > limit {
		snippets = snippets[:limit]
	}
	out := make([]DiffSnippet, 0, len(snippets))
	for _, snippet := range snippets {
		snippet.Diff = truncateReviewText(snippet.Diff, maxAIContextSnippetBytes)
		if strings.TrimSpace(snippet.File) == "" || strings.TrimSpace(snippet.Diff) == "" {
			continue
		}
		out = append(out, snippet)
	}
	return out
}

func compactStaticToolResults(results []StaticToolResult) []StaticToolResult {
	out := make([]StaticToolResult, 0, len(results))
	for _, result := range results {
		if result.ExitCode == 0 && strings.TrimSpace(result.Output) == "" {
			result.Output = ""
			out = append(out, result)
			continue
		}
		result.Output = truncateReviewText(result.Output, maxAIStaticToolOutputBytes)
		out = append(out, result)
	}
	return out
}

func compactContextSnippets(snippets []ContextSnippet, limit int) []ContextSnippet {
	if limit > 0 && len(snippets) > limit {
		snippets = prioritizeContextSnippets(snippets)
		snippets = snippets[:limit]
	}
	out := make([]ContextSnippet, 0, len(snippets))
	for _, snippet := range snippets {
		snippet.Text = truncateReviewText(snippet.Text, maxAIContextSnippetBytes)
		out = append(out, snippet)
	}
	return out
}

func prioritizeContextSnippets(snippets []ContextSnippet) []ContextSnippet {
	out := append([]ContextSnippet(nil), snippets...)
	sort.SliceStable(out, func(i, j int) bool {
		return contextSnippetPriority(out[i]) < contextSnippetPriority(out[j])
	})
	return out
}

func contextSnippetPriority(snippet ContextSnippet) int {
	switch snippet.Kind {
	case "domain_doc":
		return 0
	case "adr":
		return 1
	case "repo_doc":
		return 2
	case "review_resource":
		return 3
	case "dependency_manifest":
		return 4
	case "code_quality_file":
		return 5
	case "module_file":
		return 6
	default:
		if snippet.Source == "indexed" {
			return 4
		}
		return 7
	}
}

func limitCodeQualityHints(hints []CodeQualityHint, limit int) []CodeQualityHint {
	if limit <= 0 || len(hints) <= limit {
		return hints
	}
	return hints[:limit]
}

func limitModules(modules []ModuleSummary, limit int) []ModuleSummary {
	if limit <= 0 || len(modules) <= limit {
		return modules
	}
	return modules[:limit]
}

func limitStrings(values []string, limit int) []string {
	if limit <= 0 || len(values) <= limit {
		return values
	}
	return values[:limit]
}

func truncateReviewText(text string, limit int) string {
	text = strings.TrimSpace(text)
	if limit <= 0 || len(text) <= limit {
		return text
	}
	return text[:limit] + "\n[truncated]\n"
}

func aiRecommendationsToFindings(recommendations []aiRecommendation) []Finding {
	var out []Finding
	for i, rec := range recommendations {
		title := strings.TrimSpace(rec.Title)
		summary := strings.TrimSpace(rec.Summary)
		benefit := strings.TrimSpace(rec.Benefit)
		recommendation := strings.TrimSpace(rec.Recommendation)
		if title == "" || summary == "" || benefit == "" || recommendation == "" {
			continue
		}
		strength := strings.TrimSpace(rec.Strength)
		if strength == "" {
			strength = "Worth exploring"
		}
		var evidence []Evidence
		for _, item := range rec.Evidence {
			item = strings.TrimSpace(item)
			if item != "" {
				evidence = append(evidence, Evidence{Label: "Evidence", Value: item})
			}
		}
		out = append(out, Finding{
			ID:             fmt.Sprintf("ai.review.%d", i+1),
			Scopes:         []string{"architecture", "dependencies", "testing", "maintainability"},
			Title:          title,
			Summary:        summary,
			Benefit:        benefit,
			Evidence:       evidence,
			Recommendation: recommendation,
			Strength:       strength,
		})
	}
	return out
}

func reviewDeveloperPrompt() string {
	return strings.Join([]string{
		"You are GX Review. Review the provided patch and context for concrete recommendations, not generic audit facts.",
		"Use review_profile and depth to choose behavior: patch_focused means current-change review; prompt_directed means use review_prompt to guide a broader review of how the current diff affects the surrounding codebase; scope_focused means the requested scope; deep_full_spectrum means full-spectrum review.",
		"When review_prompt is present, answer it directly. Treat static.diff_snippets as evidence for why the prompted concern matters now, but inspect surrounding Modules, Interfaces, tests, docs, local policy, and retrieved context when they explain impact or the correct fix.",
		"For patch_focused reviews, prioritize concrete bugs, security/auth issues, data correctness, race/idempotency, error handling, missing tests, observability, deploy/CI risks, and dependency regressions introduced or exposed by static.diff_snippets.",
		"For patch_focused reviews, broad architecture, naming, docs, cleanup, or Module-depth advice is invalid unless it directly explains a changed-line bug or review risk.",
		"For prompt_directed reviews, prioritize findings where review_prompt, the current diff, and broader repo context intersect. Do not limit yourself to changed lines, but do not emit generic repo-wide advice unrelated to review_prompt.",
		"For deep_full_spectrum reviews, check security, bugs, data integrity, concurrency, idempotency, architecture, testing, observability, performance, dependencies, docs, and operability while still grounding every finding in changed files, tool output, local policy, or retrieved context.",
		"Use the architecture vocabulary exactly when discussing structure: Module, Interface, Implementation, Depth, deep, shallow, seam, adapter, leverage, locality.",
		"Never use component, service, API, boundary, or layer when Module, Interface, seam, or adapter fits.",
		"Treat static facts and hints as clues only. Do not turn file counts, missing docs, or missing tests directly into findings.",
		"Use static_tool_results as hard evidence. Failed tests, vet warnings, compile errors, and linter-like diagnostics should outrank speculative architecture advice. Ignore skipped tool results unless the skipped reason itself is clearly actionable.",
		"Use static.diff_snippets as the primary evidence for what changed. Prefer findings tied to changed lines over repo-wide advice.",
		"Use code_quality_hints as concrete candidates. Confirm whether they matter from the provided snippets before recommending a fix.",
		"Treat source_refs as the attribution set. AI output is synthesis, not evidence; every recommendation must be traceable to static facts, diff snippets, context snippets, or source_refs.",
		"Every recommendation must name at least one changed file path, Module, static tool result, code_quality_hint, or context source label. Do not produce coverage-only or structure-only recommendations without concrete evidence.",
		"The recommendation field must be concrete work: name the specific files or Modules to touch, the first operation to perform, and the verification to run. Avoid vague verbs like assess, consider, clarify, improve, harden, or refactor unless followed by exact code actions.",
		"The benefit field must state the expected payoff in concrete engineering terms: performance, readability, fewer lines of code, better error handling, better testability, lower coupling, faster onboarding, more reproducible dependencies, or better observability.",
		"If you cannot name a concrete payoff, do not emit that recommendation.",
		"In deep_full_spectrum or architecture scope only, explore like the architecture skill: find deepening opportunities, leaked Implementation knowledge, shallow Interfaces, unclear seams, weak locality, test friction, and unjustified adapters.",
		"Apply the deletion test: if deleting a Module removes complexity, call it shallow; if complexity spreads across callers, the Module is earning its keep.",
		"Use dependency categories internally: in-process, local-substitutable, remote but owned, true external. Mention adapters only when the seam needs more than one adapter.",
		"Use source_catalog and labeled context snippets internally. It is okay to mention source labels like R1 or L2 in evidence, but never output source titles or URLs.",
		"Only produce recommendations tied to the provided repo context. Reject generic best-practice advice.",
		"If no concrete issue meets the active profile, return an empty recommendations array.",
		"Return JSON only with shape {\"recommendations\":[{\"title\":string,\"summary\":string,\"benefit\":string,\"recommendation\":string,\"strength\":\"Strong|Worth exploring|Speculative\",\"evidence\":[string]}]}.",
		"Return at most 5 recommendations. Prefer 2-3 high-signal recommendations.",
	}, "\n")
}

func mustJSON(value any) string {
	body, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(body)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
