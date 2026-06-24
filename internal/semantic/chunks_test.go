package semantic

import (
	"reflect"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/reviewbundle"
)

func TestBuildSessionChunksUsesRevisionTranscriptSources(t *testing.T) {
	responseID := "response-one"
	model := "gpt-5.1-code"
	bundle := reviewbundle.Bundle{
		Repo: reviewbundle.RepoPayload{RootPath: "/repo"},
		Stack: []reviewbundle.StackPayload{{
			BranchName: "feature/alpha",
			Change: reviewbundle.ChangePayload{
				JJChangeID:      "change-alpha",
				CurrentCommitID: "commit-alpha",
				Description:     "feat alpha",
				ReviewContext: &reviewbundle.ReviewContextPayload{
					TranscriptSources: []reviewbundle.ReviewTranscriptSource{{
						SessionID:  "session-one",
						RequestID:  "request-one",
						ResponseID: &responseID,
						Provider:   "openai",
						Model:      &model,
						CreatedAt:  3,
						Status:     "linked",
					}},
					AgentProvenance: []reviewbundle.ReviewAgentProvenance{{
						SessionID: "session-one",
						AgentTool: "codex",
						Provider:  "openai",
						ModelID:   model,
						CreatedAt: 3,
					}},
				},
			},
		}},
		Sessions: []reviewbundle.SessionPayload{{
			ID:      "session-one",
			Command: "codex",
			Requests: []reviewbundle.RequestPayload{{
				ID:             "request-one",
				SessionID:      "session-one",
				Provider:       "openai",
				Endpoint:       "/v1/responses",
				Method:         "POST",
				Model:          &model,
				RequestBody:    []byte(`{"input":"implement alpha"}`),
				RequestHeaders: "{}",
				Responses: []reviewbundle.ResponsePayload{{
					ID:              "response-one",
					RequestID:       "request-one",
					ResponseBody:    []byte(`{"output":"done"}`),
					ResponseHeaders: "{}",
					StatusCode:      200,
				}},
			}},
		}},
	}

	chunks := BuildSessionChunks(bundle, 12000)
	if len(chunks) != 1 {
		t.Fatalf("chunks = %#v, want one", chunks)
	}
	chunk := chunks[0]
	if chunk.ID == "" || chunk.Text == "" {
		t.Fatalf("chunk = %#v, want stable id and text", chunk)
	}
	for _, want := range []string{"feat alpha", "implement alpha", "done"} {
		if !contains(chunk.Text, want) {
			t.Fatalf("chunk text = %q, missing %q", chunk.Text, want)
		}
	}
	wantAttrs := map[string]any{
		"source_kind":       "session_transcript",
		"repo_root":         "/repo",
		"branch_name":       "feature/alpha",
		"jj_change_id":      "change-alpha",
		"session_id":        "session-one",
		"request_id":        "request-one",
		"response_id":       "response-one",
		"agent_tool":        "codex",
		"provider":          "openai",
		"model":             model,
		"provenance_status": "linked",
	}
	for key, want := range wantAttrs {
		if got := chunk.Attributes[key]; !reflect.DeepEqual(got, want) {
			t.Fatalf("attribute %s = %#v, want %#v", key, got, want)
		}
	}
}

func contains(text, sub string) bool {
	return strings.Contains(text, sub)
}
