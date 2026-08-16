package cloud

import (
	"strings"
	"testing"
)

func TestModelNotAllowedFromBody(t *testing.T) {
	body := []byte(`{"error":"model_not_allowed","message":"Model \"zai.glm-5\" is not available through this endpoint.","allowed_models":["us.anthropic.claude-haiku-4-5-20251001-v1:0","us.anthropic.claude-sonnet-4-6","us.anthropic.claude-opus-4-6-v1"]}`)
	e := modelNotAllowedFromBody(body)
	if e == nil || e.Model != "zai.glm-5" || len(e.Allowed) != 3 {
		t.Fatalf("parse: %+v", e)
	}
	msg := e.Error()
	for _, want := range []string{"zai.glm-5", "claude-haiku-4-5-20251001", "claude-sonnet-4-6", "GX_REVIEW_BEDROCK_DIRECT=1"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message missing %q: %s", want, msg)
		}
	}
	if strings.Contains(msg, "{") || strings.Contains(msg, "allowed_models") {
		t.Errorf("message must not quote the JSON body: %s", msg)
	}
	if modelNotAllowedFromBody([]byte(`{"error":"other"}`)) != nil || modelNotAllowedFromBody([]byte(`not json`)) != nil {
		t.Errorf("only model_not_allowed bodies should be recognized")
	}
}
