package codereview

import (
	"encoding/json"
	"fmt"
	"strings"
)

type aiReviewResponse struct {
	Overview         string             `json:"overview"`
	DownstreamImpact string             `json:"downstream_impact"`
	Recommendations  []aiRecommendation `json:"recommendations"`
	NotableChanges   []aiNotableChange  `json:"notable_changes"`
}

type aiReviewOutput struct {
	Overview         string
	DownstreamImpact string
	NotableChanges   []NotableChange
	Findings         []Finding
}

type aiNotableChange struct {
	File string          `json:"file"`
	Line json.RawMessage `json:"line"`
	Note string          `json:"note"`
}

type aiRecommendation struct {
	Title          string          `json:"title"`
	Summary        string          `json:"summary"`
	Benefit        string          `json:"benefit"`
	Recommendation string          `json:"recommendation"`
	Kind           string          `json:"kind"`
	RuleID         string          `json:"rule_id"`
	Strength       string          `json:"strength"`
	Evidence       []string        `json:"evidence"`
	File           string          `json:"file"`
	Line           json.RawMessage `json:"line"`
	SourceLabels   []string        `json:"source_labels"`
	Anchors        []FindingAnchor `json:"anchors"`
	Sources        []string        `json:"sources"`
}

// describeTruncatedCompletion names the output cap as the cause when a reply
// that would not parse was also cut off at it.
//
// It is only consulted after a parse failure, never instead of one: a reply can
// stop at max_tokens with its JSON already complete — the model simply wanted to
// keep talking — and that reply is perfectly usable. Guessing from the stop
// reason alone would throw it away.
func describeTruncatedCompletion(what string, maxOutputTokens int, parseErr error) error {
	return fmt.Errorf("%s response stopped at the %d-token output cap before its JSON was complete, so no verdicts could be read (%w)",
		what, maxOutputTokens, parseErr)
}

func parseAIReviewContent(content string, brief ReviewBrief) ([]Finding, error) {
	output, err := parseAIReviewOutput(content, brief)
	if err != nil {
		return nil, err
	}
	return output.Findings, nil
}

func ParsePRSummaryReview(content string, brief ReviewBrief) (PRSummaryReview, error) {
	output, err := parseAIReviewOutput(content, brief)
	if err != nil {
		return PRSummaryReview{}, err
	}
	return aiReviewOutputToPRSummaryReview(output), nil
}

// parseAIReviewOutput tolerates a model that wraps its JSON in prose or a fenced
// block, exactly as parseJudgeResponse does.
//
// This is the reviewer's half of a guarantee that was lost with the OpenAI
// reviewer. That path could set response_format:json_object and be handed a bare
// object every time; bedrock-runtime has no equivalent, so a leg is free to
// answer with ```json ... ``` whenever it feels like it. It is not a rare shape:
// the first live cloud review run failed on BOTH legs at once with "invalid
// character '`' looking for beginning of value", which is a review that produces
// zero findings and reports only a degraded line. The direct-transport run
// minutes earlier had parsed cleanly, so this is model nondeterminism rather
// than anything about the wire, and either transport can hit it.
//
// The fallback runs only after a strict parse fails, so well-formed responses —
// including every PR-summary golden — take the original path untouched.
func parseAIReviewOutput(content string, brief ReviewBrief) (aiReviewOutput, error) {
	var parsed aiReviewResponse
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		trimmed := extractJSONObject(content, "recommendations", "overview", "notable_changes", "downstream_impact")
		if trimmed == "" {
			return aiReviewOutput{}, fmt.Errorf("decode AI review JSON: %w", err)
		}
		if fallbackErr := json.Unmarshal([]byte(trimmed), &parsed); fallbackErr != nil {
			return aiReviewOutput{}, fmt.Errorf("decode AI review JSON: %w", fallbackErr)
		}
	}
	return aiReviewOutput{
		Overview:         strings.TrimSpace(parsed.Overview),
		DownstreamImpact: strings.TrimSpace(parsed.DownstreamImpact),
		NotableChanges:   aiNotableChangesToNotableChanges(parsed.NotableChanges),
		Findings:         aiRecommendationsToFindings(parsed.Recommendations, brief),
	}, nil
}

func aiReviewOutputToPRSummaryReview(output aiReviewOutput) PRSummaryReview {
	return PRSummaryReview{
		Overview:         output.Overview,
		DownstreamImpact: output.DownstreamImpact,
		NotableChanges:   append([]NotableChange(nil), output.NotableChanges...),
		Findings:         output.Findings,
	}
}

func aiNotableChangesToNotableChanges(raw []aiNotableChange) []NotableChange {
	if len(raw) == 0 {
		return nil
	}
	out := make([]NotableChange, 0, len(raw))
	for _, item := range raw {
		file := strings.TrimSpace(item.File)
		note := strings.TrimSpace(item.Note)
		if file == "" || note == "" {
			continue
		}
		out = append(out, NotableChange{
			File: file,
			Line: parseRecommendationLine(item.Line),
			Note: note,
		})
	}
	return out
}

// normalizeFindingKind folds the model's kind label onto the three the schema
// defines. Empty stays empty rather than defaulting: "the model did not say"
// is information — downstream ranking must not treat an unlabeled finding as a
// labeled one.
func normalizeFindingKind(kind string) string {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "defect", "bug":
		return "defect"
	case "hardening", "robustness":
		return "hardening"
	case "suggestion", "improvement", "style":
		return "suggestion"
	}
	return ""
}

func aiRecommendationsToFindings(recommendations []aiRecommendation, brief ReviewBrief) []Finding {
	var out []Finding
	knownRules := KnownRuleIDs(brief.Rules...)
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
		file, line := recommendationFileLine(rec)
		labels := append([]string(nil), rec.SourceLabels...)
		if len(labels) == 0 {
			labels = append(labels, rec.Sources...)
		}
		var anchors []FindingAnchor
		if file != "" && line > 0 {
			anchors = []FindingAnchor{{File: file, Line: line}}
		} else if len(rec.Anchors) > 0 {
			anchors = rec.Anchors
			if file == "" && len(rec.Anchors) > 0 {
				file = strings.TrimSpace(rec.Anchors[0].File)
				line = rec.Anchors[0].Line
			}
		}
		out = append(out, Finding{
			ID:              fmt.Sprintf("ai.review.%d", i+1),
			Scopes:          []string{"architecture", "dependencies", "testing", "maintainability"},
			Title:           title,
			Summary:         summary,
			Benefit:         benefit,
			Evidence:        evidence,
			File:            file,
			Line:            line,
			Anchors:         anchors,
			Recommendation:  recommendation,
			Kind:            normalizeFindingKind(rec.Kind),
			RuleID:          NormalizeRuleID(rec.RuleID, knownRules),
			Strength:        strength,
			ResolvedSources: resolveSourceLabels(brief, labels),
		})
	}
	return out
}

func recommendationFileLine(rec aiRecommendation) (string, int) {
	file := strings.TrimSpace(rec.File)
	line := parseRecommendationLine(rec.Line)
	if file != "" && line > 0 {
		return file, line
	}
	if len(rec.Anchors) > 0 {
		anchor := rec.Anchors[0]
		return strings.TrimSpace(anchor.File), anchor.Line
	}
	return file, line
}

func parseRecommendationLine(raw json.RawMessage) int {
	if len(raw) == 0 {
		return 0
	}
	var asInt int
	if err := json.Unmarshal(raw, &asInt); err == nil {
		return asInt
	}
	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		asString = strings.TrimSpace(asString)
		if asString == "" {
			return 0
		}
		var parsed int
		_, _ = fmt.Sscanf(asString, "%d", &parsed)
		return parsed
	}
	return 0
}
