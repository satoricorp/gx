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
	t.Setenv("GX_GATE_MODEL", "")
	a, b := resolveBedrockReviewModels()
	if a != bedrockHaikuModel {
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

// The shipped panel (2026-08-23): GPT-5.6 Luna on both reviewer legs, Haiku
// judging. The judge is a different model from the reviewers so it never
// grades its own output; both IDs are in their Bedrock-invocable form and take
// the transport path their vendor needs.
func TestShippedDefaultsAreLunaReviewersAndHaikuJudge(t *testing.T) {
	t.Setenv("GX_REVIEW_MODELS", "")
	for _, k := range []string{"GX_REVIEW_BEDROCK_MODEL_A", "GX_REVIEW_BEDROCK_MODEL_B", "GX_REVIEW_ANTHROPIC_MODEL", "GX_REVIEW_JUDGE_MODEL", "GX_GATE_MODEL"} {
		t.Setenv(k, "")
	}
	a, b := resolveBedrockReviewModels()
	judge := resolveBedrockJudgeModel()
	if a != "us.openai.gpt-5.6-luna" || b != "us.openai.gpt-5.6-luna" {
		t.Fatalf("default reviewers = %q, %q; want GPT-5.6 Luna on both legs", a, b)
	}
	if judge != bedrockHaikuModel {
		t.Fatalf("default judge = %q; want Haiku 4.5", judge)
	}
	if judge == a {
		t.Fatalf("the judge must not be a reviewer model")
	}
	if bedrockActionForModel(a) != "converse" || bedrockActionForModel(judge) != "invoke" {
		t.Fatalf("Luna rides Converse and Haiku InvokeModel; got %s / %s", bedrockActionForModel(a), bedrockActionForModel(judge))
	}
	// "default" as an explicit preset name is the same panel.
	t.Setenv("GX_REVIEW_MODELS", "default")
	if pa, pb := resolveBedrockReviewModels(); pa != a || pb != b || resolveBedrockJudgeModel() != judge {
		t.Fatalf("GX_REVIEW_MODELS=default must equal the unset defaults")
	}
}

// The pre-2026-08-23 all-Claude panel stays reachable by name.
func TestClaudePresetIsThePreviousPanel(t *testing.T) {
	t.Setenv("GX_REVIEW_MODELS", "claude")
	for _, k := range []string{"GX_REVIEW_BEDROCK_MODEL_A", "GX_REVIEW_BEDROCK_MODEL_B", "GX_REVIEW_ANTHROPIC_MODEL", "GX_REVIEW_JUDGE_MODEL", "GX_GATE_MODEL"} {
		t.Setenv(k, "")
	}
	a, b := resolveBedrockReviewModels()
	if a != bedrockHaikuModel || b != bedrockSonnetModel || resolveBedrockJudgeModel() != bedrockSonnetModel {
		t.Fatalf("claude preset: a=%q b=%q judge=%q", a, b, resolveBedrockJudgeModel())
	}
}

func TestLunaPresetAndName(t *testing.T) {
	t.Setenv("GX_REVIEW_MODELS", "luna")
	for _, k := range []string{"GX_REVIEW_BEDROCK_MODEL_A", "GX_REVIEW_BEDROCK_MODEL_B", "GX_REVIEW_ANTHROPIC_MODEL", "GX_REVIEW_JUDGE_MODEL", "GX_GATE_MODEL"} {
		t.Setenv(k, "")
	}
	a, b := resolveBedrockReviewModels()
	if a != "us.openai.gpt-5.6-luna" || b != "us.openai.gpt-5.6-luna" || resolveBedrockJudgeModel() != "us.openai.gpt-5.6-luna" {
		t.Fatalf("luna preset: a=%q b=%q judge=%q", a, b, resolveBedrockJudgeModel())
	}
	// The us. profile prefix survives normalization (only anthropic. IDs are
	// rewritten) and the model routes to Converse.
	if normalizeBedrockModelID("us.openai.gpt-5.6-luna") != "us.openai.gpt-5.6-luna" {
		t.Fatalf("normalizer must not touch a us.openai profile ID")
	}
	if bedrockActionForModel("us.openai.gpt-5.6-luna") != "converse" {
		t.Fatalf("openai models on Bedrock speak Converse")
	}
	if got := friendlyModelName("us.openai.gpt-5.6-luna"); got != "GPT-5.6 Luna" {
		t.Fatalf("friendlyModelName(luna) = %q", got)
	}
	if got := friendlyModelName("openai.gpt-oss-120b-1:0"); got != "GPT-oss-120b-1" {
		t.Fatalf("friendlyModelName(gpt-oss) = %q", got)
	}
}
