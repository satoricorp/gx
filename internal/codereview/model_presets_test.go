package codereview

import (
	"strings"
	"testing"
)

func TestModelPresetFillsEverySlot(t *testing.T) {
	t.Setenv("GX_REVIEW_MODELS", "budget")
	t.Setenv("GX_REVIEW_BEDROCK_MODEL_A", "")
	t.Setenv("GX_REVIEW_BEDROCK_MODEL_B", "")
	t.Setenv("GX_REVIEW_ANTHROPIC_MODEL", "")
	t.Setenv("GX_REVIEW_JUDGE_MODEL", "")
	t.Setenv("GX_CONSTRAINTS_MODEL", "")
	a, b := resolveBedrockReviewModels()
	if a != defaultBedrockReviewModelA {
		t.Fatalf("budget leg A = %q, want Haiku kept", a)
	}
	if b != "zai.glm-5" {
		t.Fatalf("budget leg B = %q, want zai.glm-5", b)
	}
	if j := resolveBedrockJudgeModel(); j != "nvidia.nemotron-super-3-120b" {
		t.Fatalf("budget judge = %q, want nvidia.nemotron-super-3-120b", j)
	}
	// Non-Anthropic IDs are not given a us. prefix and ride Converse.
	if bedrockActionForModel(b) != "converse" || bedrockActionForModel("nvidia.nemotron-super-3-120b") != "converse" {
		t.Fatalf("open-weight models must route to the Converse API")
	}
}

func TestExplicitEnvBeatsPreset(t *testing.T) {
	t.Setenv("GX_REVIEW_MODELS", "budget")
	t.Setenv("GX_REVIEW_JUDGE_MODEL", "us.anthropic.claude-sonnet-4-6")
	t.Setenv("GX_REVIEW_BEDROCK_MODEL_B", "")
	if j := resolveBedrockJudgeModel(); j != "us.anthropic.claude-sonnet-4-6" {
		t.Fatalf("explicit GX_REVIEW_JUDGE_MODEL must beat the preset, got %q", j)
	}
	if _, b := resolveBedrockReviewModels(); b != "zai.glm-5" {
		t.Fatalf("slots without an explicit env still take the preset, got %q", b)
	}
}

func TestUnknownPresetIsIgnoredByResolversAndCaughtByValidate(t *testing.T) {
	t.Setenv("GX_REVIEW_MODELS", "chaep")
	t.Setenv("GX_REVIEW_BEDROCK_MODEL_B", "")
	if _, b := resolveBedrockReviewModels(); b != defaultBedrockReviewModelB {
		t.Fatalf("unknown preset must leave defaults alone, got %q", b)
	}
	err := ValidateModelPresetEnv()
	if err == nil || !strings.Contains(err.Error(), "budget") || !strings.Contains(err.Error(), "chaep") {
		t.Fatalf("validate should name the bad value and the valid presets, got %v", err)
	}
	t.Setenv("GX_REVIEW_MODELS", "")
	if err := ValidateModelPresetEnv(); err != nil {
		t.Fatalf("unset is fine, got %v", err)
	}
	t.Setenv("GX_REVIEW_MODELS", "Default")
	if err := ValidateModelPresetEnv(); err != nil {
		t.Fatalf("preset names are case-insensitive, got %v", err)
	}
}

func TestFriendlyNamesForOpenWeightModels(t *testing.T) {
	cases := map[string]string{
		"zai.glm-5":                    "GLM 5",
		"zai.glm-4.7":                  "GLM 4.7",
		"nvidia.nemotron-super-3-120b": "Nemotron 3 Super",
		"nvidia.nemotron-nano-3-30b":   "Nemotron 3 Nano",
	}
	for in, want := range cases {
		if got := friendlyModelName(in); got != want {
			t.Errorf("friendlyModelName(%q) = %q, want %q", in, got, want)
		}
	}
}
