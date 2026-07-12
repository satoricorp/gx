package authoring

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"

	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/version"
)

const (
	defaultDemuxReviewBaseURL     = "https://api.openai.com"
	defaultDemuxReviewModel       = "gpt-4.1-mini"
	defaultDemuxReviewMaxWarnings = 50
	defaultDemuxReviewTimeout     = 5 * time.Minute
	defaultDemuxMaxOutputTokens   = 16000
)

type demuxAIReviewer interface {
	ReviewDemuxProposal(ctx context.Context, req demuxAIReviewRequest) (demuxAIReviewResponse, error)
}

type demuxAIReviewRequest struct {
	Proposal    demuxAIProposalForReview `json:"proposal"`
	Review      demuxAIReviewSummary     `json:"review"`
	MaxWarnings int                      `json:"max_warnings"`
}

type demuxAIReviewResponse struct {
	Proposal      DemuxProposal      `json:"proposal,omitempty"`
	Revisions     []RevisionProposal `json:"revisions,omitempty"`
	LLMConfidence float64            `json:"llm_confidence,omitempty"`
	LLMReasons    []ConfidenceReason `json:"llm_reasons,omitempty"`
	Notes         []string           `json:"notes,omitempty"`
	Model         string             `json:"model,omitempty"`
}

type cloudChatCompletionRequest struct {
	Model               string             `json:"model"`
	Messages            []cloudChatMessage `json:"messages"`
	ResponseFormat      map[string]string  `json:"response_format,omitempty"`
	MaxCompletionTokens int                `json:"max_completion_tokens,omitempty"`
}

type cloudChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type cloudChatCompletionResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type openAIDemuxReviewer struct {
	apiKey  string
	baseURL string
	model   string
	client  openai.Client
}

type cloudDemuxReviewer struct {
	url    string
	token  string
	model  string
	client *http.Client
}

type fallbackDemuxReviewer struct {
	primary  demuxAIReviewer
	fallback demuxAIReviewer
}

type demuxRepairModule struct {
	engine *Engine
}

func (e *Engine) demuxRepair() demuxRepairModule {
	return demuxRepairModule{engine: e}
}

func (e *Engine) ReviewDemuxProposalWithAI(ctx context.Context, selector string, opts DemuxAIReviewOptions) (DemuxAIReviewResult, error) {
	return e.demuxPipeline().repairSavedProposal(ctx, selector, opts)
}

func (e *Engine) RepairDemuxProposal(ctx context.Context, proposal DemuxProposal, opts DemuxAIReviewOptions) (DemuxAIReviewResult, error) {
	return e.demuxPipeline().repairProposal(ctx, proposal, opts)
}

func (e *Engine) repairReviewedDemuxProposal(ctx context.Context, proposal DemuxProposal, review ReviewDemuxResult, opts DemuxAIReviewOptions) (DemuxAIReviewResult, error) {
	return e.demuxRepair().repairReviewed(ctx, proposal, review, opts)
}

func (r demuxRepairModule) repair(ctx context.Context, proposal DemuxProposal, opts DemuxAIReviewOptions) (DemuxAIReviewResult, error) {
	review, err := r.engine.ReviewDemuxPlan(ctx, proposal)
	if err != nil {
		return DemuxAIReviewResult{}, err
	}
	return r.repairReviewed(ctx, proposal, review, opts)
}

func (r demuxRepairModule) repairApplyPreflightFailure(ctx context.Context, proposal DemuxProposal, preflightErr error, opts DemuxAIReviewOptions) (DemuxAIReviewResult, error) {
	review, err := r.engine.ReviewDemuxPlan(ctx, proposal)
	if err != nil {
		return DemuxAIReviewResult{}, err
	}
	message := "apply preflight failed: " + preflightErr.Error()
	review.Valid = false
	review.Errors = append(review.Errors, message)
	review.RepairHints = append(review.RepairHints, RepairHint{
		Kind:       "apply_preflight",
		Suggestion: "revise revision order, grouping, hunk coverage, target_stack, or base_stack so the proposal can be applied in a clean disposable checkout",
	})
	review.Proposal = appendDemuxApplyPreflightWarning(review.Proposal, preflightErr)
	proposal = review.Proposal
	return r.repairReviewed(ctx, proposal, review, opts)
}

func (r demuxRepairModule) repairReviewed(ctx context.Context, proposal DemuxProposal, review ReviewDemuxResult, opts DemuxAIReviewOptions) (DemuxAIReviewResult, error) {
	maxWarnings := opts.MaxWarnings
	if maxWarnings <= 0 {
		maxWarnings = defaultDemuxReviewMaxWarnings
	}
	if repaired, changed := r.autoRepairToFixedPoint(ctx, review); changed {
		saved, err := r.engine.SaveDemuxProposal(ctx, repaired.Proposal)
		if err != nil {
			return DemuxAIReviewResult{}, err
		}
		repaired.Proposal = saved
		result := DemuxAIReviewResult{
			Proposal:        saved,
			Review:          repaired,
			State:           demuxWorkflowState(repaired, saved),
			Model:           "deterministic",
			Updated:         true,
			Notes:           []string{"applied deterministic demux repair hints"},
			WarningsTotal:   len(review.Proposal.FeasibilityWarnings),
			RepairHintTotal: len(review.RepairHints),
		}
		proposal = saved
		review = repaired
		if opts.PlanOnly || result.State == DemuxWorkflowReadyToApply {
			return result, nil
		}
	}
	if demuxWorkflowState(review, review.Proposal) == DemuxWorkflowReadyToApply {
		return DemuxAIReviewResult{
			Proposal:        review.Proposal,
			Review:          review,
			State:           DemuxWorkflowReadyToApply,
			Model:           "deterministic",
			Updated:         false,
			Notes:           []string{"deterministic proposal is ready to apply"},
			WarningsTotal:   len(review.Proposal.FeasibilityWarnings),
			RepairHintTotal: len(review.RepairHints),
		}, nil
	}
	if opts.PlanOnly {
		return DemuxAIReviewResult{
			Proposal:        review.Proposal,
			Review:          review,
			State:           demuxWorkflowState(review, review.Proposal),
			Model:           "deterministic",
			Updated:         false,
			Notes:           []string{"no deterministic repair was available; rerun without --plan to ask OpenAI to revise the proposal"},
			WarningsTotal:   len(review.Proposal.FeasibilityWarnings),
			RepairHintTotal: len(review.RepairHints),
		}, nil
	}
	reviewer, model, err := demuxReviewerFromEnv(opts)
	if err != nil {
		return DemuxAIReviewResult{
			Proposal:        review.Proposal,
			Review:          review,
			State:           demuxWorkflowState(review, review.Proposal),
			Model:           "deterministic",
			Updated:         false,
			WarningsTotal:   len(review.Proposal.FeasibilityWarnings),
			RepairHintTotal: len(review.RepairHints),
		}, err
	}
	request, warningsSent, repairHintsSent := aiReviewRequestForProposal(review, maxWarnings)
	response, err := reviewer.ReviewDemuxProposal(ctx, demuxAIReviewRequest{
		Proposal:    request.Proposal,
		Review:      request.Review,
		MaxWarnings: maxWarnings,
	})
	if err != nil {
		return DemuxAIReviewResult{
			Proposal:        review.Proposal,
			Review:          review,
			State:           demuxWorkflowState(review, review.Proposal),
			Model:           model,
			Updated:         false,
			WarningsSent:    warningsSent,
			WarningsTotal:   len(review.Proposal.FeasibilityWarnings),
			RepairHintSent:  repairHintsSent,
			RepairHintTotal: len(review.RepairHints),
		}, err
	}
	revised := proposalFromAIReviewResponse(proposal, response)
	if strings.TrimSpace(revised.ID) == "" {
		return DemuxAIReviewResult{}, fmt.Errorf("AI review did not return proposal.id")
	}
	reviewed, err := r.engine.ReviewDemuxPlan(ctx, revised)
	if err != nil {
		return DemuxAIReviewResult{}, err
	}
	if !reviewed.Valid {
		if repaired, changed := autoRepairDemuxProposal(reviewed.Proposal, reviewed.RepairHints); changed {
			rereviewed, err := r.engine.ReviewDemuxPlan(ctx, repaired)
			if err != nil {
				return DemuxAIReviewResult{}, err
			}
			if rereviewed.Valid || len(rereviewed.Errors) < len(reviewed.Errors) || demuxRepairScore(rereviewed) < demuxRepairScore(reviewed) {
				reviewed = rereviewed
			}
		}
	}
	if !reviewed.Valid {
		detail := strings.Join(reviewed.Errors, "; ")
		if detail == "" && len(reviewed.RepairHints) > 0 {
			detail = reviewed.RepairHints[0].Suggestion
		}
		if detail == "" {
			detail = "no validation details returned"
		}
		return DemuxAIReviewResult{
			Proposal:        reviewed.Proposal,
			Review:          reviewed,
			State:           demuxWorkflowState(reviewed, reviewed.Proposal),
			Model:           firstNonEmpty(response.Model, model),
			Updated:         false,
			Notes:           response.Notes,
			WarningsSent:    warningsSent,
			WarningsTotal:   len(review.Proposal.FeasibilityWarnings),
			RepairHintSent:  repairHintsSent,
			RepairHintTotal: len(review.RepairHints),
		}, fmt.Errorf("AI review returned an invalid proposal: %s", detail)
	}
	for range 4 {
		repaired, changed := autoRepairDemuxProposal(reviewed.Proposal, reviewed.RepairHints)
		if !changed {
			break
		}
		rereviewed, err := r.engine.ReviewDemuxPlan(ctx, repaired)
		if err != nil {
			return DemuxAIReviewResult{}, err
		}
		if !rereviewed.Valid || demuxRepairScore(rereviewed) >= demuxRepairScore(reviewed) {
			break
		}
		reviewed = rereviewed
		if demuxWorkflowState(reviewed, reviewed.Proposal) == DemuxWorkflowReadyToApply {
			break
		}
	}
	saved, err := r.engine.SaveDemuxProposal(ctx, reviewed.Proposal)
	if err != nil {
		return DemuxAIReviewResult{}, err
	}
	reviewed.Proposal = saved
	result := DemuxAIReviewResult{
		Proposal:        saved,
		Review:          reviewed,
		State:           demuxWorkflowState(reviewed, saved),
		Model:           firstNonEmpty(response.Model, model),
		Updated:         true,
		Notes:           response.Notes,
		WarningsSent:    warningsSent,
		WarningsTotal:   len(review.Proposal.FeasibilityWarnings),
		RepairHintSent:  repairHintsSent,
		RepairHintTotal: len(review.RepairHints),
	}
	if result.State != DemuxWorkflowReadyToApply {
		return result, fmt.Errorf("demux fix did not produce an apply-ready proposal: %d blocking warning(s), %d repair hint(s) remain", len(blockingFeasibilityWarnings(result.Proposal.FeasibilityWarnings)), len(result.Review.RepairHints))
	}
	return result, nil
}

func demuxReviewerFromEnv(opts DemuxAIReviewOptions) (demuxAIReviewer, string, error) {
	direct, model, directErr := openAIDemuxReviewerFromEnv(opts)
	cloudReviewer, cloudModel, cloudErr := cloudDemuxReviewerFromEnv(opts)
	if direct != nil {
		if cloudReviewer != nil {
			return fallbackDemuxReviewer{primary: direct, fallback: cloudReviewer}, model, nil
		}
		return direct, model, nil
	}
	if directErr != nil {
		return nil, "", directErr
	}
	if cloudReviewer != nil {
		return cloudReviewer, cloudModel, nil
	}
	if cloudErr != nil {
		return nil, "", fmt.Errorf("gx generate needs a model to repair this plan: sign in with `gx auth login` to use GX Cloud, or set your own key with `gx set key` (cloud unavailable: %w)", cloudErr)
	}
	return nil, "", fmt.Errorf("gx generate needs a model to repair this plan: sign in with `gx auth login` to use GX Cloud, or set your own key with `gx set key`")
}

func openAIDemuxReviewerFromEnv(opts DemuxAIReviewOptions) (demuxAIReviewer, string, error) {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	baseURLRaw := strings.TrimSpace(os.Getenv("OPENAI_BASE_URL"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("GX_OPENAI_API_KEY"))
		baseURLRaw = strings.TrimSpace(firstNonEmpty(os.Getenv("GX_OPENAI_BASE_URL"), os.Getenv("OPENAI_BASE_URL")))
	}
	if apiKey == "" {
		return nil, "", nil
	}
	baseURL := normalizeOpenAIBaseURL(firstNonEmpty(baseURLRaw, defaultDemuxReviewBaseURL))
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, "", fmt.Errorf("OpenAI base URL is invalid: %w", err)
	}
	model := strings.TrimSpace(firstNonEmpty(opts.Model, os.Getenv("GX_DEMUX_REVIEW_MODEL"), defaultDemuxReviewModel))
	reviewer := newOpenAIDemuxReviewer(apiKey, baseURL, model)
	if baseURLRaw != "" && baseURL != normalizeOpenAIBaseURL(defaultDemuxReviewBaseURL) {
		return fallbackDemuxReviewer{
			primary:  reviewer,
			fallback: newOpenAIDemuxReviewer(apiKey, normalizeOpenAIBaseURL(defaultDemuxReviewBaseURL), model),
		}, model, nil
	}
	return reviewer, model, nil
}

func cloudDemuxReviewerFromEnv(opts DemuxAIReviewOptions) (demuxAIReviewer, string, error) {
	url := strings.TrimSpace(os.Getenv("GX_OPENAI_PROXY_URL"))
	if url == "" {
		baseURL := cloud.CloudBaseURL()
		if baseURL == "" {
			return nil, "", fmt.Errorf("gx cloud base URL is not configured")
		}
		url = strings.TrimRight(baseURL, "/") + "/gx/openai/chat-completions"
	}
	token, err := cloud.CloudAPIToken()
	if err != nil {
		return nil, "", err
	}
	model := strings.TrimSpace(firstNonEmpty(opts.Model, os.Getenv("GX_DEMUX_REVIEW_MODEL"), defaultDemuxReviewModel))
	return &cloudDemuxReviewer{
		url:    url,
		token:  token,
		model:  model,
		client: &http.Client{Timeout: defaultDemuxReviewTimeout},
	}, model, nil
}

func newOpenAIDemuxReviewer(apiKey, baseURL, model string) *openAIDemuxReviewer {
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(apiKey),
		option.WithRequestTimeout(defaultDemuxReviewTimeout),
	)
	return &openAIDemuxReviewer{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client:  client,
	}
}

func (r fallbackDemuxReviewer) ReviewDemuxProposal(ctx context.Context, req demuxAIReviewRequest) (demuxAIReviewResponse, error) {
	response, err := r.primary.ReviewDemuxProposal(ctx, req)
	if err == nil {
		return response, nil
	}
	fallbackResponse, fallbackErr := r.fallback.ReviewDemuxProposal(ctx, req)
	if fallbackErr == nil {
		return fallbackResponse, nil
	}
	return demuxAIReviewResponse{}, fmt.Errorf("%w; fallback reviewer failed: %v", err, fallbackErr)
}

func (r *cloudDemuxReviewer) ReviewDemuxProposal(ctx context.Context, req demuxAIReviewRequest) (demuxAIReviewResponse, error) {
	payload := cloudChatCompletionRequest{
		Model: r.model,
		Messages: []cloudChatMessage{
			{Role: "developer", Content: demuxAIReviewDeveloperPrompt()},
			{Role: "user", Content: mustJSONForPrompt(req)},
		},
		ResponseFormat:      map[string]string{"type": "json_object"},
		MaxCompletionTokens: defaultDemuxMaxOutputTokens,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return demuxAIReviewResponse{}, fmt.Errorf("marshal gx cloud OpenAI request: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, r.url, bytes.NewReader(body))
	if err != nil {
		return demuxAIReviewResponse{}, fmt.Errorf("create gx cloud OpenAI request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+r.token)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "gx/"+version.Current())

	resp, err := r.client.Do(httpReq)
	if err != nil {
		return demuxAIReviewResponse{}, fmt.Errorf("request gx cloud OpenAI demux review: %w", err)
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode == http.StatusPaymentRequired {
		return demuxAIReviewResponse{}, cloud.NewPaymentRequiredError(responseBody)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		detail := strings.TrimSpace(string(responseBody))
		if detail != "" {
			return demuxAIReviewResponse{}, fmt.Errorf("gx cloud OpenAI demux review: status %s: %s", resp.Status, detail)
		}
		return demuxAIReviewResponse{}, fmt.Errorf("gx cloud OpenAI demux review: status %s", resp.Status)
	}

	var completion cloudChatCompletionResponse
	if err := json.Unmarshal(responseBody, &completion); err != nil {
		return demuxAIReviewResponse{}, fmt.Errorf("decode gx cloud OpenAI response: %w", err)
	}
	if len(completion.Choices) == 0 {
		return demuxAIReviewResponse{}, fmt.Errorf("gx cloud OpenAI demux review response returned no choices")
	}
	var out demuxAIReviewResponse
	if err := json.Unmarshal([]byte(completion.Choices[0].Message.Content), &out); err != nil {
		return demuxAIReviewResponse{}, fmt.Errorf("decode gx cloud AI demux review JSON: %w", err)
	}
	if out.Model == "" {
		out.Model = completion.Model
	}
	return out, nil
}

func (r *openAIDemuxReviewer) ReviewDemuxProposal(ctx context.Context, req demuxAIReviewRequest) (demuxAIReviewResponse, error) {
	response, err := r.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: shared.ChatModel(r.model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.DeveloperMessage(demuxAIReviewDeveloperPrompt()),
			openai.UserMessage(mustJSONForPrompt(req)),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: demuxPtr(shared.NewResponseFormatJSONObjectParam()),
		},
		MaxCompletionTokens: param.NewOpt(int64(defaultDemuxMaxOutputTokens)),
	})
	if err != nil {
		return demuxAIReviewResponse{}, fmt.Errorf("request OpenAI demux review: %w", err)
	}
	if len(response.Choices) == 0 {
		return demuxAIReviewResponse{}, fmt.Errorf("OpenAI demux review response returned no choices")
	}
	var out demuxAIReviewResponse
	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &out); err != nil {
		return demuxAIReviewResponse{}, fmt.Errorf("decode AI demux review JSON: %w", err)
	}
	if out.Model == "" {
		out.Model = response.Model
	}
	return out, nil
}

func normalizeOpenAIBaseURL(raw string) string {
	baseURL := strings.TrimRight(strings.TrimSpace(raw), "/")
	if baseURL == "" {
		baseURL = strings.TrimRight(defaultDemuxReviewBaseURL, "/")
	}
	if strings.HasSuffix(baseURL, "/v1") {
		return baseURL
	}
	return baseURL + "/v1"
}

func demuxPtr[T any](value T) *T {
	return &value
}

func demuxAIReviewDeveloperPrompt() string {
	return strings.Join([]string{
		"You repair GX demux proposals.",
		"Return only JSON with shape {\"revisions\": <full ordered revision array>, \"llm_confidence\": number, \"llm_reasons\": [confidence reason], \"notes\": [string]}.",
		"Return every revision, in final order. Do not return patches or the top-level proposal object.",
		"Preserve revision ids, provenance_status, session_ids, and valid hunk_ids.",
		"Use hunk_ids from the provided hunk catalog; do not invent hunk ids.",
		"Every hunk must be assigned exactly once.",
		"For unassigned_hunk hints, assign the named hunk_id to an existing revision or create a new revision for that hunk.",
		"Fix ordering and grouping to reduce repair_hints and warning-severity feasibility warnings.",
		"Preserve target_stack and base_stack when they are coherent.",
		"For invalid_demux_route hints, set the revision target_stack to the hinted target_stack.",
		"For cross_stack_dependency hints, either keep dependent revisions on one target_stack or set the dependent revision base_stack to the hinted base_stack.",
		"Your output must be apply-ready: no warning-severity feasibility_warnings should remain after GX reviews the returned revisions.",
		"If ordering alone cannot fix a warning, merge the dependent revisions instead of leaving a blocking dependency.",
		"Use the provided confidence_summary as the deterministic logic score for the heuristic draft.",
		"Return llm_confidence for semantic grouping quality only; GX will combine it with logic_confidence using the lower value.",
	}, "\n")
}

func aiReviewRequestForProposal(review ReviewDemuxResult, maxWarnings int) (demuxAIReviewRequest, int, int) {
	warnings := limitedWarningSeverity(review.Proposal.FeasibilityWarnings, maxWarnings)
	repairHints := limitRepairHints(review.RepairHints, maxWarnings)
	return demuxAIReviewRequest{
		Proposal: proposalForAIReview(review.Proposal, maxWarnings),
		Review: demuxAIReviewSummary{
			Valid:               review.Valid,
			Errors:              append([]string(nil), review.Errors...),
			FeasibilityWarnings: warnings,
			RepairHints:         repairHints,
		},
		MaxWarnings: maxWarnings,
	}, len(warnings), len(repairHints)
}

func proposalForAIReview(proposal DemuxProposal, maxWarnings int) demuxAIProposalForReview {
	hunks := make([]HunkRange, 0, len(proposal.Hunks))
	for _, hunk := range proposal.Hunks {
		hunk.Patch = ""
		hunks = append(hunks, hunk)
	}
	revisions := append([]RevisionProposal(nil), proposal.Revisions...)
	for index := range revisions {
		revisions[index].Hunks = nil
	}
	return demuxAIProposalForReview{
		ID:               proposal.ID,
		RepoRoot:         proposal.RepoRoot,
		ProposedChangeID: proposal.ProposedChangeID,
		ProposedCommitID: proposal.ProposedCommitID,
		Status:           proposal.Status,
		PlanInstructions: append([]string(nil), proposal.PlanInstructions...),
		Confidence:       proposal.Confidence,
		Hunks:            hunks,
		StructuralDeps:   limitStructuralDependencies(proposal.StructuralDeps, maxWarnings),
		ChangedSymbols:   append([]ChangedSymbol(nil), proposal.ChangedSymbols...),
		Revisions:        revisions,
	}
}

func proposalFromAIReviewResponse(base DemuxProposal, response demuxAIReviewResponse) DemuxProposal {
	if len(response.Revisions) > 0 {
		base.Revisions = response.Revisions
		base.FeasibilityWarnings = nil
		base.Confidence.LLMConfidence = response.LLMConfidence
		base.Confidence.LLMReasons = append([]ConfidenceReason(nil), response.LLMReasons...)
		return base
	}
	if strings.TrimSpace(response.Proposal.ID) == "" {
		response.Proposal.ID = base.ID
	}
	if response.LLMConfidence > 0 {
		response.Proposal.Confidence.LLMConfidence = response.LLMConfidence
		response.Proposal.Confidence.LLMReasons = append([]ConfidenceReason(nil), response.LLMReasons...)
	}
	return response.Proposal
}

func (r demuxRepairModule) autoRepairToFixedPoint(ctx context.Context, review ReviewDemuxResult) (ReviewDemuxResult, bool) {
	best := review
	bestScore := demuxRepairScore(review)
	current := review
	changedAny := false
	for range 10 {
		repaired, changed := autoRepairDemuxProposal(current.Proposal, current.RepairHints)
		if !changed {
			break
		}
		reviewed, err := r.engine.ReviewDemuxPlan(ctx, repaired)
		if err != nil {
			break
		}
		changedAny = true
		score := demuxRepairScore(reviewed)
		if score <= bestScore {
			best = reviewed
			bestScore = score
		}
		current = reviewed
		if len(blockingFeasibilityWarnings(reviewed.Proposal.FeasibilityWarnings)) == 0 {
			break
		}
	}
	return best, changedAny
}

func demuxRepairScore(review ReviewDemuxResult) int {
	if !review.Valid {
		return 1_000_000 + len(review.Errors)*1_000 + len(review.RepairHints)
	}
	return len(blockingFeasibilityWarnings(review.Proposal.FeasibilityWarnings))*10 + len(review.RepairHints)
}

func limitedWarningSeverity(warnings []FeasibilityWarning, limit int) []FeasibilityWarning {
	if limit <= 0 || len(warnings) == 0 {
		return nil
	}
	out := make([]FeasibilityWarning, 0, minInt(limit, len(warnings)))
	for _, warning := range warnings {
		if !strings.EqualFold(strings.TrimSpace(warning.Severity), "warning") {
			continue
		}
		out = append(out, warning)
		if len(out) >= limit {
			break
		}
	}
	return out
}

func limitRepairHints(hints []RepairHint, limit int) []RepairHint {
	if limit <= 0 || len(hints) == 0 {
		return nil
	}
	if len(hints) <= limit {
		return append([]RepairHint(nil), hints...)
	}
	return append([]RepairHint(nil), hints[:limit]...)
}

func limitStructuralDependencies(deps []StructuralDependency, limit int) []StructuralDependency {
	if limit <= 0 || len(deps) == 0 {
		return nil
	}
	if len(deps) <= limit {
		return append([]StructuralDependency(nil), deps...)
	}
	return append([]StructuralDependency(nil), deps[:limit]...)
}

func autoRepairDemuxProposal(proposal DemuxProposal, hints []RepairHint) (DemuxProposal, bool) {
	if len(proposal.Revisions) == 0 {
		return proposal, false
	}
	revisions, cleanedDeps := removeInvalidDependsOn(proposal.Revisions)
	proposal.Revisions = revisions
	proposal, repairedRoutes := repairDemuxRoutes(proposal, hints)
	proposal, shaped := shapeDemuxProposalForReview(proposal)
	if len(proposal.Revisions) < 2 || len(hints) == 0 {
		return proposal, cleanedDeps || repairedRoutes || shaped
	}

	deps := map[string]map[string]struct{}{}
	for _, hint := range hints {
		revisionID := strings.TrimSpace(hint.RevisionID)
		dependsOn := strings.TrimSpace(hint.DependsOn)
		if revisionID == "" || dependsOn == "" || revisionID == dependsOn {
			continue
		}
		switch hint.Kind {
		case "reorder_dependency", "inferred_dependency":
		default:
			continue
		}
		if deps[revisionID] == nil {
			deps[revisionID] = map[string]struct{}{}
		}
		deps[revisionID][dependsOn] = struct{}{}
	}
	if len(deps) == 0 {
		return proposal, cleanedDeps || repairedRoutes || shaped
	}

	ordered, changedOrder := orderRevisionsByDependencies(proposal.Revisions, deps)
	ordered, cleanedAfterOrder := removeInvalidDependsOn(ordered)
	ordered, changedDeps := addEarlierDependsOn(ordered, deps)
	ordered, prunedTransitiveDeps := pruneTransitiveDependsOn(ordered)
	proposal.Revisions = ordered
	proposal, shapedAfterOrder := shapeDemuxProposalForReview(proposal)
	return proposal, cleanedDeps || repairedRoutes || shaped || changedOrder || cleanedAfterOrder || changedDeps || prunedTransitiveDeps || shapedAfterOrder
}

func repairDemuxRoutes(proposal DemuxProposal, hints []RepairHint) (DemuxProposal, bool) {
	if len(hints) == 0 {
		return proposal, false
	}
	byID := map[string]int{}
	for index, revision := range proposal.Revisions {
		byID[revision.ID] = index
	}
	changed := false
	for _, hint := range hints {
		index, ok := byID[strings.TrimSpace(hint.RevisionID)]
		if !ok {
			continue
		}
		revision := &proposal.Revisions[index]
		switch hint.Kind {
		case "invalid_demux_route":
			targetStack := strings.TrimSpace(hint.TargetStack)
			if targetStack != "" && revision.TargetStack != targetStack {
				revision.TargetStack = targetStack
				changed = true
			}
		case "cross_stack_dependency":
			baseStack := strings.TrimSpace(hint.BaseStack)
			if baseStack != "" && revision.BaseStack != baseStack {
				revision.BaseStack = baseStack
				changed = true
			}
		}
	}
	return proposal, changed
}

func removeInvalidDependsOn(revisions []RevisionProposal) ([]RevisionProposal, bool) {
	out := append([]RevisionProposal(nil), revisions...)
	seen := map[string]struct{}{}
	changed := false
	for index := range out {
		var kept []string
		for _, dep := range out[index].DependsOn {
			dep = strings.TrimSpace(dep)
			if dep == "" {
				continue
			}
			if _, ok := seen[dep]; ok {
				kept = append(kept, dep)
				continue
			}
			changed = true
		}
		out[index].DependsOn = kept
		seen[out[index].ID] = struct{}{}
	}
	return out, changed
}

func addEarlierDependsOn(revisions []RevisionProposal, deps map[string]map[string]struct{}) ([]RevisionProposal, bool) {
	out := append([]RevisionProposal(nil), revisions...)
	indexByID := map[string]int{}
	for index, revision := range out {
		indexByID[revision.ID] = index
	}
	changed := false
	for index := range out {
		for dep := range deps[out[index].ID] {
			depIndex, ok := indexByID[dep]
			if !ok || depIndex >= index || revisionDependsOn(out[index], dep) {
				continue
			}
			out[index].DependsOn = append(out[index].DependsOn, dep)
			changed = true
		}
	}
	return out, changed
}

func pruneTransitiveDependsOn(revisions []RevisionProposal) ([]RevisionProposal, bool) {
	out := append([]RevisionProposal(nil), revisions...)
	byID := map[string]RevisionProposal{}
	changed := false
	for index := range out {
		deps := appendUniqueStrings(nil, out[index].DependsOn...)
		if len(deps) < 2 {
			out[index].DependsOn = deps
			byID[out[index].ID] = out[index]
			continue
		}
		var kept []string
		for _, dep := range deps {
			redundant := false
			for _, other := range deps {
				if other == dep {
					continue
				}
				if dependencyReachable(other, dep, byID, map[string]bool{}) {
					redundant = true
					break
				}
			}
			if redundant {
				changed = true
				continue
			}
			kept = append(kept, dep)
		}
		out[index].DependsOn = kept
		byID[out[index].ID] = out[index]
	}
	return out, changed
}

func dependencyReachable(from, target string, byID map[string]RevisionProposal, seen map[string]bool) bool {
	if from == target {
		return true
	}
	if seen[from] {
		return false
	}
	seen[from] = true
	revision, ok := byID[from]
	if !ok {
		return false
	}
	for _, dep := range revision.DependsOn {
		if dependencyReachable(dep, target, byID, seen) {
			return true
		}
	}
	return false
}

func orderRevisionsByDependencies(revisions []RevisionProposal, deps map[string]map[string]struct{}) ([]RevisionProposal, bool) {
	byID := map[string]RevisionProposal{}
	originalIndex := map[string]int{}
	for index, revision := range revisions {
		byID[revision.ID] = revision
		originalIndex[revision.ID] = index
	}

	indegree := map[string]int{}
	outgoing := map[string][]string{}
	for _, revision := range revisions {
		indegree[revision.ID] = 0
	}
	for revisionID, revisionDeps := range deps {
		if _, ok := byID[revisionID]; !ok {
			continue
		}
		for dep := range revisionDeps {
			if _, ok := byID[dep]; !ok {
				continue
			}
			outgoing[dep] = append(outgoing[dep], revisionID)
			indegree[revisionID]++
		}
	}

	var ready []string
	for _, revision := range revisions {
		if indegree[revision.ID] == 0 {
			ready = append(ready, revision.ID)
		}
	}
	sort.SliceStable(ready, func(i, j int) bool {
		return originalIndex[ready[i]] < originalIndex[ready[j]]
	})

	var ordered []RevisionProposal
	for len(ready) > 0 {
		current := ready[0]
		ready = ready[1:]
		ordered = append(ordered, byID[current])
		next := outgoing[current]
		sort.SliceStable(next, func(i, j int) bool {
			return originalIndex[next[i]] < originalIndex[next[j]]
		})
		for _, dependent := range next {
			indegree[dependent]--
			if indegree[dependent] == 0 {
				ready = append(ready, dependent)
			}
		}
		sort.SliceStable(ready, func(i, j int) bool {
			return originalIndex[ready[i]] < originalIndex[ready[j]]
		})
	}
	if len(ordered) != len(revisions) {
		return revisions, false
	}
	changed := false
	for index := range ordered {
		if ordered[index].ID != revisions[index].ID {
			changed = true
			break
		}
	}
	return ordered, changed
}

func mustJSONForPrompt(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(data)
}

func envInt(name string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
