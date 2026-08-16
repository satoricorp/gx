package codereview

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type aiReviewResponse struct {
	Overview         string             `json:"overview"`
	DownstreamImpact string             `json:"downstream_impact"`
	Recommendations  []aiRecommendation `json:"recommendations"`
	NotableChanges   []aiNotableChange  `json:"notable_changes"`
	Story            []aiStoryItem      `json:"story"`
}

type aiReviewOutput struct {
	Overview         string
	DownstreamImpact string
	NotableChanges   []NotableChange
	Story            []StoryItem
	Findings         []Finding
}

type aiNotableChange struct {
	File string          `json:"file"`
	Line json.RawMessage `json:"line"`
	Note string          `json:"note"`
}

// aiStoryItem is the wire shape of one story entry. Line is raw because the
// model returns it as a number or a quoted string, same as everywhere else.
type aiStoryItem struct {
	Headline    string          `json:"headline"`
	Category    string          `json:"category"`
	Materiality string          `json:"materiality"`
	File        string          `json:"file"`
	Line        json.RawMessage `json:"line"`
	Consequence string          `json:"consequence"`
	Asked       string          `json:"asked"`
	Chose       string          `json:"chose"`
	Watch       []string        `json:"watch"`
	PassedRule  string          `json:"passed_rule"`
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
	Example        string          `json:"example"`
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
		trimmed := extractJSONObject(content, "recommendations", "overview", "notable_changes", "downstream_impact", "story")
		if trimmed == "" {
			// No complete object anywhere. Before giving up, salvage: a reply
			// cut off at the output cap is a well-formed prefix of a JSON
			// object, and every recommendation element that closed before
			// the cut is intact. Measured on a whole-repo review: a Sonnet
			// shard reply of 21 KB stopped at max_tokens with 13 complete
			// findings before the 14th was cut mid-sentence, and all 13 were
			// discarded. Keeping them is strictly better than losing the
			// shard; the caller still reports the truncation as degradation.
			if recs := salvageTruncatedRecommendations(content); len(recs) > 0 {
				parsed.Recommendations = recs
				return aiReviewOutput{
					Findings: aiRecommendationsToFindings(parsed.Recommendations, brief),
				}, errTruncatedButSalvaged
			}
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
		Story:            aiStoryToStoryItems(parsed.Story, brief),
		Findings:         aiRecommendationsToFindings(parsed.Recommendations, brief),
	}, nil
}

func aiReviewOutputToPRSummaryReview(output aiReviewOutput) PRSummaryReview {
	return PRSummaryReview{
		Overview:         output.Overview,
		DownstreamImpact: output.DownstreamImpact,
		NotableChanges:   append([]NotableChange(nil), output.NotableChanges...),
		Story:            append([]StoryItem(nil), output.Story...),
		Findings:         output.Findings,
	}
}

// aiStoryToStoryItems converts the model's story entries. An entry with no
// headline or no consequence is not a story item — those two are what make it
// one — so it is dropped. Category and materiality are folded onto the schema's
// vocabulary; passed_rule is validated against the rules the brief carried, so
// an invented rule id never reaches the reader.
//
// exception_taken is reserved for MandatoryStoryItems: a change exempting
// itself from a rule is reported by the review, not by the model's discretion,
// so a model entry claiming that category keeps its text and loses the label.
func aiStoryToStoryItems(raw []aiStoryItem, brief ReviewBrief) []StoryItem {
	if len(raw) == 0 {
		return nil
	}
	knownRules := KnownRuleIDs(brief.Rules...)
	out := make([]StoryItem, 0, len(raw))
	for _, item := range raw {
		headline := strings.TrimSpace(item.Headline)
		consequence := strings.TrimSpace(item.Consequence)
		if headline == "" || consequence == "" {
			continue
		}
		category := normalizeStoryCategory(item.Category)
		if category == StoryExceptionTaken {
			category = ""
		}
		var watch []string
		for _, w := range item.Watch {
			if w = strings.TrimSpace(w); w != "" {
				watch = append(watch, w)
			}
		}
		out = append(out, StoryItem{
			Headline:    headline,
			Category:    category,
			Materiality: normalizeMateriality(item.Materiality),
			File:        strings.TrimSpace(item.File),
			Line:        parseRecommendationLine(item.Line),
			Consequence: consequence,
			Asked:       strings.TrimSpace(item.Asked),
			Chose:       strings.TrimSpace(item.Chose),
			Watch:       watch,
			PassedRule:  NormalizeRuleID(item.PassedRule, knownRules),
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
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
			Example:         normalizeExampleDiff(rec.Example),
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

// errTruncatedButSalvaged is returned alongside a non-empty output when the
// reply was cut off but complete findings were recovered from its prefix. The
// caller keeps the findings and records the truncation as a degradation
// rather than treating the shard as lost.
var errTruncatedButSalvaged = errors.New("AI review reply was truncated; complete findings before the cut were kept")

// salvageTruncatedRecommendations pulls every complete element out of a
// truncated "recommendations": [ ... array. It walks the array with
// json.Decoder, which tracks strings and escapes, so a "}" inside a summary
// ends nothing, and stops at the first element that does not decode — that is
// the one the cap cut. Anything before it is intact.
func salvageTruncatedRecommendations(content string) []aiRecommendation {
	key := strings.Index(content, "\"recommendations\"")
	if key < 0 {
		return nil
	}
	open := strings.IndexByte(content[key:], '[')
	if open < 0 {
		return nil
	}
	dec := json.NewDecoder(strings.NewReader(content[key+open:]))
	if tok, err := dec.Token(); err != nil || tok != json.Delim('[') {
		return nil
	}
	var out []aiRecommendation
	for dec.More() {
		var rec aiRecommendation
		if err := dec.Decode(&rec); err != nil {
			break // the cut element; everything before it is complete
		}
		out = append(out, rec)
	}
	return out
}

// normalizeExampleDiff keeps an example only when it looks like a small diff
// fragment: some line begins with + or -, and it is short. Anything else — a
// prose paragraph the model put in the wrong field, a whole file — is dropped,
// because the renderer paints it as a hunk and a hunk of prose is worse than
// no example. Fenced ```diff blocks are unwrapped first.
func normalizeExampleDiff(raw string) string {
	text := strings.TrimSpace(raw)
	if text == "" {
		return ""
	}
	if strings.HasPrefix(text, "```") {
		if i := strings.IndexByte(text, '\n'); i >= 0 {
			text = text[i+1:]
		}
		text = strings.TrimSuffix(strings.TrimSpace(text), "```")
		text = strings.TrimSpace(text)
	}
	lines := strings.Split(text, "\n")
	if len(lines) > 12 {
		return ""
	}
	diffLike := false
	for _, l := range lines {
		if strings.HasPrefix(l, "+") || strings.HasPrefix(l, "-") {
			diffLike = true
			break
		}
	}
	if !diffLike {
		return ""
	}
	return strings.Join(lines, "\n")
}
