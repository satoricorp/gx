package codereview

import (
	"os"
	"sort"
	"strings"
)

// Model presets.
//
// A review is four model slots — reviewer A, reviewer B, the judge, and the
// constraints judge — and every one is an env var. That is the right shape
// for a deployment and the wrong shape for an experiment: comparing "the
// default panel" against "a cheaper panel" means exporting four variables,
// remembering which four, and unsetting them afterwards.
//
// GX_REVIEW_MODELS names a preset that fills all four at once. Precedence per
// slot stays: explicit env var > preset > built-in default. So
// GX_REVIEW_MODELS=budget swaps the whole panel, and GX_REVIEW_MODELS=budget
// GX_REVIEW_JUDGE_MODEL=us.anthropic.claude-sonnet-4-6 swaps everything but
// the judge. Unknown preset names are ignored with the defaults left alone,
// so a typo cannot silently pick a model.
//
// Presets carry Bedrock model IDs. Non-Anthropic IDs (zai.*, nvidia.*) ride
// the Converse API automatically — bedrockActionForModel decides per model —
// so a preset can mix vendors freely.

// modelPreset is one named panel.
type modelPreset struct {
	// Why is one line for the operator: what this preset trades for what.
	Why         string
	ReviewerA   string
	ReviewerB   string
	Judge       string
	Constraints string
}

var modelPresets = map[string]modelPreset{
	// The shipped defaults, spelled out so "default" is a valid preset name and
	// so the table documents them next to the alternatives.
	"default": {
		Why:         "two Claude reviewers across tiers, Claude judge — the panel gx ships with",
		ReviewerA:   defaultBedrockReviewModelA,
		ReviewerB:   defaultBedrockReviewModelB,
		Judge:       defaultBedrockJudgeModel,
		Constraints: defaultBedrockReviewModelA,
	},
	// budget keeps Haiku as leg A (fast, cheap, and the model the constraints
	// judge already runs on) and swaps the two expensive slots for open-weight
	// models on Bedrock: GLM 5 as the second reviewer so the panel stays two
	// vendors with uncorrelated misses, and Nemotron 3 Super as the judge —
	// the per-finding verification call, which is where cost concentrates.
	// Both are on-demand in us-west-2 and speak Converse.
	"budget": {
		Why:         "Haiku + GLM 5 review, Nemotron 3 Super judges — open-weight models in the expensive slots",
		ReviewerA:   defaultBedrockReviewModelA,
		ReviewerB:   "zai.glm-5",
		Judge:       "nvidia.nemotron-super-3-120b",
		Constraints: defaultBedrockReviewModelA,
	},
	// glm puts GLM 5 in every graded slot — the "how does one open model do
	// on its own" comparison, with Haiku only where it was already the
	// default.
	"glm": {
		Why:         "GLM 5 in every slot — one open model, end to end",
		ReviewerA:   "zai.glm-5",
		ReviewerB:   "zai.glm-4.7",
		Judge:       "zai.glm-5",
		Constraints: "zai.glm-5",
	},
}

// PresetNames lists the presets in a stable order, for help text and errors.
func PresetNames() []string {
	names := make([]string, 0, len(modelPresets))
	for name := range modelPresets {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// activeModelPreset returns the preset GX_REVIEW_MODELS names, or false when
// it is unset or unknown. Unknown is deliberately not an error here — the
// callers are model resolvers deep in the review path, and the CLI validates
// the name up front where a message can reach the user.
func activeModelPreset() (modelPreset, bool) {
	name := strings.ToLower(strings.TrimSpace(os.Getenv("GX_REVIEW_MODELS")))
	if name == "" {
		return modelPreset{}, false
	}
	preset, ok := modelPresets[name]
	return preset, ok
}

// ValidateModelPresetEnv reports an error naming the valid presets when
// GX_REVIEW_MODELS is set to something that is not one. Called from the CLI
// before a review starts.
func ValidateModelPresetEnv() error {
	name := strings.ToLower(strings.TrimSpace(os.Getenv("GX_REVIEW_MODELS")))
	if name == "" {
		return nil
	}
	if _, ok := modelPresets[name]; ok {
		return nil
	}
	return &unknownPresetError{name: name}
}

type unknownPresetError struct{ name string }

func (e *unknownPresetError) Error() string {
	return "GX_REVIEW_MODELS=" + e.name + " is not a preset; valid: " + strings.Join(PresetNames(), ", ")
}

// presetOr returns the preset's value for a slot when a preset is active and
// that slot is filled, else the fallback. It sits between the env var and the
// built-in default in every resolver.
func presetOr(pick func(modelPreset) string, fallback string) string {
	if preset, ok := activeModelPreset(); ok {
		if v := strings.TrimSpace(pick(preset)); v != "" {
			return v
		}
	}
	return fallback
}
