package codereview

import "testing"

func TestFriendlyModelName(t *testing.T) {
	cases := map[string]string{
		"us.anthropic.claude-haiku-4-5-20251001-v1:0": "Claude Haiku 4.5",
		"us.anthropic.claude-sonnet-4-6":              "Claude Sonnet 4.6",
		"us.anthropic.claude-opus-4-7-v1:0":           "Claude Opus 4.7",
		"anthropic.claude-sonnet-5":                   "Claude Sonnet 5",
		"claude-haiku-4-5":                            "Claude Haiku 4.5",
		"":                                            "",
		"arn:aws:bedrock:us-west-2::model/x":          "arn:aws:bedrock:us-west-2::model/x",
		"some-other-model":                            "some-other-model",
	}
	for in, want := range cases {
		if got := friendlyModelName(in); got != want {
			t.Errorf("friendlyModelName(%q) = %q, want %q", in, got, want)
		}
	}
}
