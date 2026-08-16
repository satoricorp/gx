package cloud

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ModelNotAllowedError is gx Cloud's answer when a review asks for a Bedrock
// model that /gx/bedrock/fight does not allow. The server sends the allowlist
// with it, and this carries the list through so the report can say what IS
// allowed in one line instead of quoting a JSON body.
type ModelNotAllowedError struct {
	Model   string
	Allowed []string
}

func (e *ModelNotAllowedError) Error() string {
	if e == nil {
		return "gx Cloud does not allow this model"
	}
	// Show the short, human forms: drop regional prefixes and vendor
	// namespaces the reader does not type.
	names := make([]string, 0, len(e.Allowed))
	for _, m := range e.Allowed {
		names = append(names, shortModelName(m))
	}
	sort.Strings(names)
	list := strings.Join(names, ", ")
	if list == "" {
		return fmt.Sprintf("gx Cloud does not allow model %s through /gx/bedrock/fight", e.Model)
	}
	return fmt.Sprintf("gx Cloud does not allow model %s; it allows %s (to use another model, set GX_REVIEW_BEDROCK_DIRECT=1 with AWS credentials)", shortModelName(e.Model), list)
}

// modelNotAllowedFromBody recognizes the server's model_not_allowed body.
// Returns nil for any other body so callers keep their generic path.
func modelNotAllowedFromBody(body []byte) *ModelNotAllowedError {
	var payload struct {
		Error   string   `json:"error"`
		Message string   `json:"message"`
		Allowed []string `json:"allowed_models"`
	}
	if err := json.Unmarshal(body, &payload); err != nil || payload.Error != "model_not_allowed" {
		return nil
	}
	// The model name is inside the message: Model "x" is not available…
	model := ""
	if i := strings.IndexByte(payload.Message, '"'); i >= 0 {
		if j := strings.IndexByte(payload.Message[i+1:], '"'); j >= 0 {
			model = payload.Message[i+1 : i+1+j]
		}
	}
	return &ModelNotAllowedError{Model: model, Allowed: payload.Allowed}
}

// shortModelName strips a regional profile prefix and, for Anthropic, the
// vendor namespace: us.anthropic.claude-haiku-4-5-20251001-v1:0 →
// claude-haiku-4-5-20251001; zai.glm-5 → zai.glm-5.
func shortModelName(model string) string {
	model = strings.TrimSpace(model)
	for _, p := range []string{"us.", "eu.", "apac.", "global.", "us-gov."} {
		if strings.HasPrefix(model, p) {
			model = strings.TrimPrefix(model, p)
			break
		}
	}
	model = strings.TrimPrefix(model, "anthropic.")
	if i := strings.Index(model, "-v1:"); i > 0 {
		model = model[:i]
	}
	return model
}
