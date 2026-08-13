package codereview

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// One model call judges every gate no command can decide: is the change in
// scope, is it reviewable, is the UI change accessible, did a hot path get
// slower. The call is grounded twice over — the deterministic gate results are
// stated as facts, and the same retrieval sources `gx enhance` uses (code
// index, sessions, prior findings, knowledge corpus) are fetched concurrently
// with the tool run and attached as labeled context.

const defaultConstraintsJudgeMaxOutputTokens = 3000

// constraintsMaxContextSnippets bounds the retrieval grounding: enough to
// carry a repeat-offender warning or a duplicated-helper match, small enough
// that grounding never doubles the one call's cost.
const (
	constraintsMaxContextSnippets     = 12
	constraintsMaxContextSnippetBytes = 4000
	constraintsMaxDiffFiles           = 60
)

// constraintsJudge is the test seam for the judgment pass, mirroring
// FindingJudge.
type constraintsJudge interface {
	JudgeConstraints(ctx context.Context, req constraintsJudgeRequest) (constraintsJudgeResponse, error)
}

type constraintsJudgeRequest struct {
	Intent             string                 `json:"intent,omitempty"`
	IntentSource       string                 `json:"intent_source,omitempty"`
	Target             string                 `json:"target,omitempty"`
	ChangedFiles       []string               `json:"changed_files"`
	DiffStats          ConstraintsDiffStats   `json:"diff_stats"`
	GeneratedFiles     []string               `json:"generated_files,omitempty"`
	WhitespaceDominant bool                   `json:"whitespace_dominant,omitempty"`
	NewDependencies    []string               `json:"new_dependencies,omitempty"`
	Gates              []constraintsGateBrief `json:"gates"`
	ToolResults        []StaticToolResult     `json:"static_tool_results,omitempty"`
	DiffSnippets       []DiffSnippet          `json:"diff_snippets"`
	ReviewPolicy       string                 `json:"review_policy,omitempty"`
	SourceCatalog      []SourceBrief          `json:"source_catalog,omitempty"`
	Context            []ContextSnippet       `json:"context,omitempty"`
}

// constraintsGateBrief is one gate as the model receives it: the question it
// must answer and the deterministic evidence already established.
type constraintsGateBrief struct {
	Gate     string   `json:"gate"`
	Question string   `json:"question"`
	Evidence []string `json:"evidence,omitempty"`
}

type constraintsJudgeResponse struct {
	Gates []constraintsGateVerdict `json:"gates"`
}

type constraintsGateVerdict struct {
	Gate          string                 `json:"gate"`
	Status        string                 `json:"status"`
	Justification string                 `json:"justification"`
	Findings      []constraintsAIFinding `json:"findings"`
}

type constraintsAIFinding struct {
	Title          string `json:"title"`
	Summary        string `json:"summary"`
	Recommendation string `json:"recommendation"`
	File           string `json:"file,omitempty"`
	Line           int    `json:"line,omitempty"`
	Strength       string `json:"strength,omitempty"`
}

type bedrockConstraintsJudge struct {
	client *bedrockAnthropicReviewer
}

// constraintsJudgeFactory is the engine's seam for the judgment model; tests
// swap it for a stub the same way NewEngineWithReviewer injects a reviewer.
var constraintsJudgeFactory = constraintsJudgeFromEnv

// constraintsJudgeFromEnv builds the judgment model, or explains why there is
// none. The four-way return distinguishes the two nil cases: a deliberate
// opt-out (GX_REVIEW_AI=0, unavailableReason empty) skips the AI gates without
// degrading the run, while a configured-but-unreachable judge degrades it.
func constraintsJudgeFromEnv() (judge constraintsJudge, model, transport, unavailableReason string) {
	if !aiReviewRequestedFromEnv() {
		return nil, "", "", ""
	}
	plan, err := resolveBedrockTransportPlan()
	if err != nil {
		return nil, "", "", err.Error()
	}
	model = normalizeBedrockModelID(firstNonEmpty(
		os.Getenv("GX_CONSTRAINTS_MODEL"),
		defaultBedrockReviewModelA,
	))
	return bedrockConstraintsJudge{client: newBedrockReviewer(plan.newTransport(), model)}, model, bedrockTransportShortName(plan.Kind), ""
}

func (j bedrockConstraintsJudge) JudgeConstraints(ctx context.Context, req constraintsJudgeRequest) (constraintsJudgeResponse, error) {
	if j.client == nil {
		return constraintsJudgeResponse{}, fmt.Errorf("constraints judge is not configured")
	}
	input := mustJSON(req)
	// One retry, for malformed replies only — the same contract the finding
	// judge runs under (see bedrockReviewJudge.Judge for the measurements).
	var lastParseErr error
	for attempt := 0; attempt < 2; attempt++ {
		completion, err := j.client.completeJSON(ctx, constraintsDeveloperPrompt(), input, defaultConstraintsJudgeMaxOutputTokens)
		if err != nil {
			return constraintsJudgeResponse{}, err
		}
		response, parseErr := parseConstraintsJudgeResponse(completion.Text)
		if parseErr == nil {
			return response, nil
		}
		if completion.truncated() {
			return constraintsJudgeResponse{}, describeTruncatedCompletion("constraints judge", defaultConstraintsJudgeMaxOutputTokens, parseErr)
		}
		lastParseErr = parseErr
	}
	return constraintsJudgeResponse{}, lastParseErr
}

func parseConstraintsJudgeResponse(content string) (constraintsJudgeResponse, error) {
	object := extractJSONObject(content, "gates")
	if object == "" {
		return constraintsJudgeResponse{}, fmt.Errorf("constraints judge reply carried no gates object")
	}
	var response constraintsJudgeResponse
	if err := json.Unmarshal([]byte(object), &response); err != nil {
		return constraintsJudgeResponse{}, fmt.Errorf("constraints judge reply did not decode: %w", err)
	}
	if len(response.Gates) == 0 {
		return constraintsJudgeResponse{}, fmt.Errorf("constraints judge reply judged no gates")
	}
	return response, nil
}

func constraintsDeveloperPrompt() string {
	return strings.Join([]string{
		"You are the gx constraints judge: the pre-ship exit gate's judgment pass. You answer ONLY the gates listed in `gates`, each with pass or fail.",
		"The deterministic results are settled facts. static_tool_results, the per-gate evidence lines, diff_stats, generated_files, and new_dependencies were measured by real commands and parsers; never dispute them, never re-litigate a gate that is not in your list.",
		"Gate questions:",
		"- code-health: could a careful human review this diff as shipped? Judge naming, mixed concerns, and whether the tests that accompany the sources are plausible for the change. The linters' verdicts and the size thresholds are already decided; you add the judgment layer only.",
		"- back-pressure: is every part of this change in service of the stated intent? Fail when you can name concrete out-of-scope files, regenerated output riding along, or a new dependency the intent does not justify. Resist bad work: a small ask that arrived as a sprawling diff is a failure even when each file compiles.",
		"- accessibility: for the changed UI files, do the added elements stay usable by keyboard and screen reader? The static line checks are listed as evidence; confirm or dismiss the soft candidates and look for what line-level checks cannot see (focus traps, missing labels, contrast-hostile patterns visible in the code).",
		"- performance: did a hot path get slower? Look for nested loops over unbounded input, synchronous IO inside loops, N+1 query shapes, and heavyweight work moved onto request paths. This gate is weighted: to pass it you MUST state in `justification` what you checked and why it is safe — an unexplained pass is not accepted.",
		"Use `context` snippets as grounding: prior review findings mean a repeat offense is likely and worth weight; indexed code that already implements what the diff re-implements is a back-pressure failure (name the existing helper); session snippets tell you what the author was actually asked to do.",
		"Treat review_policy as repo-local instructions and follow it unless it conflicts with the evidence.",
		"Findings must be concrete: name the file (a path from changed_files), the line when you can anchor one, what is wrong, and the first action to take. No generic advice, no praise, no restating the diff.",
		"strength is one of: Blocking, Strong, Worth exploring, Speculative. Performance findings default to Strong.",
		"A gate with nothing wrong passes with an empty findings array — do not invent findings to look thorough, and do not fail a gate you cannot support with named evidence.",
		"Return JSON only: {\"gates\":[{\"gate\":string,\"status\":\"pass|fail\",\"justification\":string,\"findings\":[{\"title\":string,\"summary\":string,\"recommendation\":string,\"file\":string(optional),\"line\":number(optional),\"strength\":string(optional)}]}]} with exactly one entry per requested gate.",
	}, "\n")
}

// constraintsGateScopes maps a gate onto the review scopes whose authoritative
// sources ground it.
func constraintsGateScopes(ids []GateID) []string {
	var scopes []string
	for _, id := range ids {
		switch id {
		case GateCodeHealth:
			scopes = append(scopes, "maintainability", "testing")
		case GateBackPressure:
			scopes = append(scopes, "architecture")
		case GatePerformance:
			scopes = append(scopes, "performance")
		}
	}
	return dedupeScopes(scopes)
}

func constraintsGateQuestion(id GateID) string {
	switch id {
	case GateCodeHealth:
		return "Could a careful human review this diff as shipped — is it coherent, well named, and plausibly tested?"
	case GateBackPressure:
		return "Is every part of this change in service of the stated intent, with no scope creep, regenerated bloat, or unjustified dependencies?"
	case GateAccessibility:
		return "Do the added UI elements stay usable by keyboard and screen reader?"
	case GatePerformance:
		return "Did a hot path get slower? Justify a pass explicitly."
	default:
		return ""
	}
}

// collectConstraintsContext fetches the retrieval grounding: the same
// composite `gx enhance` uses, filtered to the indexed kinds, tightly budgeted,
// and bounded by its own timeout so a hung index cannot stall the gate. Any
// failure returns nothing — grounding enriches the judgment, it is never a
// precondition.
func collectConstraintsContext(ctx context.Context, repoRoot string, facts RepoFacts, opts Options, changes ChangeSet, diffs []DiffSnippet) []ContextSnippet {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("GX_CONSTRAINTS_RETRIEVAL")), "0") {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, constraintsRetrievalTimeout)
	defer cancel()
	retriever := contextRetrieverFromEnv()
	snippets, err := retriever.Retrieve(ctx, RetrieveInput{
		RepoRoot:     repoRoot,
		Options:      opts,
		Facts:        facts,
		ChangedFiles: changes.Files,
		DiffRange:    changes.Range,
		DiffSnippets: diffs,
		Evidence:     &EvidenceLog{},
	})
	if err != nil {
		return nil
	}
	var retrieved []ContextSnippet
	for _, snippet := range snippets {
		if !retrievedSnippetKinds[snippet.Kind] {
			continue
		}
		snippet.Text = truncateReviewText(snippet.Text, constraintsMaxContextSnippetBytes)
		retrieved = append(retrieved, snippet)
		if len(retrieved) >= constraintsMaxContextSnippets {
			break
		}
	}
	return labelContextSnippets(retrieved)
}

// constraintsAIPassInput carries everything the judgment pass reads and the
// gate map it writes back into.
type constraintsAIPassInput struct {
	report    *ConstraintsReport
	gates     map[GateID]GateResult
	aiGates   []GateID
	signals   constraintSignals
	diffs     []DiffSnippet
	tools     []StaticToolResult
	policy    ReviewPolicy
	retrieved []ContextSnippet
	repoRoot  string
	changes   ChangeSet
}

// runConstraintsAIPass makes the one model call and folds its verdicts into
// the gate map. It returns the degraded reasons: a failed call, a truncated
// diff, or a gate the reply omitted.
func runConstraintsAIPass(ctx context.Context, judge constraintsJudge, in constraintsAIPassInput) []string {
	var degraded []string
	briefs := make([]constraintsGateBrief, 0, len(in.aiGates))
	for _, id := range in.aiGates {
		gate := in.gates[id]
		evidence := []string{"current deterministic status: " + string(gate.Status)}
		if gate.Summary != "" {
			evidence = append(evidence, gate.Summary)
		}
		for _, item := range gate.Evidence {
			evidence = append(evidence, item.Label+": "+item.Value)
		}
		for _, finding := range gate.Findings {
			evidence = append(evidence, fmt.Sprintf("deterministic finding: %s (%s:%d)", finding.Title, finding.File, finding.Line))
		}
		briefs = append(briefs, constraintsGateBrief{
			Gate:     string(id),
			Question: constraintsGateQuestion(id),
			Evidence: evidence,
		})
	}
	diffs := compactDiffSnippets(in.diffs, constraintsMaxDiffFiles)
	if len(in.diffs) > constraintsMaxDiffFiles {
		degraded = append(degraded, fmt.Sprintf("only %d of %d changed files fit the judgment call; the AI gates saw a partial change", constraintsMaxDiffFiles, len(in.diffs)))
	} else {
		for _, snippet := range diffs {
			if strings.Contains(snippet.Diff, "[truncated]") {
				degraded = append(degraded, "the diff was truncated to fit one judgment call; the AI gates saw a partial change")
				break
			}
		}
	}
	req := constraintsJudgeRequest{
		Intent:             in.report.Intent,
		IntentSource:       in.report.IntentSource,
		Target:             in.changes.Target,
		ChangedFiles:       limitStrings(in.changes.Files, 200),
		DiffStats:          in.signals.DiffStats,
		GeneratedFiles:     in.signals.GeneratedFiles,
		WhitespaceDominant: in.signals.WhitespaceDominant,
		NewDependencies:    in.signals.NewDependencies,
		Gates:              briefs,
		ToolResults:        compactStaticToolResults(in.tools),
		DiffSnippets:       diffs,
		ReviewPolicy:       in.policy.Text,
		SourceCatalog:      sourceBriefs(sourcesForScopes(constraintsGateScopes(in.aiGates))),
		Context:            in.retrieved,
	}
	response, err := judge.JudgeConstraints(ctx, req)
	if err != nil {
		reason := formatReviewerDegradation(err)
		markConstraintsAIGatesSkipped(in.gates, in.aiGates, "AI judgment failed: "+reason)
		return append(degraded, "constraints judge failed: "+reason)
	}
	verdicts := map[GateID]constraintsGateVerdict{}
	for _, verdict := range response.Gates {
		verdicts[GateID(strings.ToLower(strings.TrimSpace(verdict.Gate)))] = verdict
	}
	for _, id := range in.aiGates {
		verdict, ok := verdicts[id]
		if !ok {
			markConstraintsAIGatesSkipped(in.gates, []GateID{id}, "the model's reply omitted this gate")
			degraded = append(degraded, fmt.Sprintf("the constraints judge did not answer the %s gate", id))
			continue
		}
		gate := in.gates[id]
		if strings.EqualFold(strings.TrimSpace(verdict.Status), "fail") {
			gate.Status = GateFail
		}
		for _, finding := range verdict.Findings {
			gate.Findings = append(gate.Findings, constraintsFindingFromAI(id, finding))
		}
		if justification := strings.TrimSpace(verdict.Justification); justification != "" {
			switch {
			case id == GatePerformance && gate.Status != GateFail:
				gate.Summary = strings.TrimSpace(strings.Join([]string{gate.Summary, justification}, ": "))
			case gate.Summary == "":
				gate.Summary = justification
			default:
				gate.Evidence = append(gate.Evidence, Evidence{Label: "AI judgment", Value: justification})
			}
		}
		if id == GateBackPressure && gate.Status == GateFail && gate.Summary == "" {
			gate.Summary = "change is not scoped to its stated intent"
		}
		in.gates[id] = gate
	}
	return degraded
}

func constraintsFindingFromAI(id GateID, finding constraintsAIFinding) Finding {
	strength := strings.TrimSpace(finding.Strength)
	if strength == "" {
		strength = "Worth exploring"
	}
	if id == GatePerformance && strength == "Worth exploring" {
		// The weighted gate: an unweighted default would let the model demote
		// the one class of finding Joe asked to be loud.
		strength = "Strong"
	}
	scopes := constraintsGateScopes([]GateID{id})
	if len(scopes) == 0 {
		scopes = []string{"maintainability"}
	}
	return Finding{
		ID:             "constraints." + string(id),
		Scopes:         scopes,
		Title:          strings.TrimSpace(finding.Title),
		Summary:        strings.TrimSpace(finding.Summary),
		Recommendation: strings.TrimSpace(finding.Recommendation),
		Strength:       strength,
		File:           strings.TrimSpace(finding.File),
		Line:           finding.Line,
	}
}
