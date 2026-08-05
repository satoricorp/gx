package semantic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/satoricorp/lgtm/internal/reviewbundle"
)

type Chunk struct {
	ID         string
	Text       string
	Attributes map[string]any
}

func BuildSessionChunks(bundle reviewbundle.Bundle, maxBytes int) []Chunk {
	if maxBytes <= 0 {
		maxBytes = defaultMaxChunkBytes
	}
	sessionIndex := indexSessions(bundle.Sessions)
	var chunks []Chunk
	seen := map[string]struct{}{}
	for _, target := range reviewTargets(bundle) {
		if target.ReviewContext == nil {
			continue
		}
		for _, source := range target.ReviewContext.TranscriptSources {
			key := transcriptChunkKey(target.revisionKey(), source)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			chunk, ok := buildTranscriptChunk(bundle, target, source, sessionIndex, maxBytes)
			if ok {
				chunks = append(chunks, chunk)
			}
		}
	}
	return chunks
}

type reviewTarget struct {
	RevisionID    string
	CommitID      string
	Description   string
	BranchName    string
	ReviewContext *reviewbundle.ReviewContextPayload
}

// revisionKey identifies a target even when a commit carries no lgtm trailer.
func (t reviewTarget) revisionKey() string {
	if strings.TrimSpace(t.RevisionID) != "" {
		return t.RevisionID
	}
	return t.CommitID
}

func reviewTargets(bundle reviewbundle.Bundle) []reviewTarget {
	targets := make([]reviewTarget, 0, len(bundle.Revisions))
	for _, revision := range bundle.Revisions {
		targets = append(targets, reviewTarget{
			RevisionID:    revision.RevisionID,
			CommitID:      revision.CommitID,
			Description:   revision.Description,
			BranchName:    revision.BranchName,
			ReviewContext: revision.ReviewContext,
		})
	}
	return targets
}

func indexSessions(sessions []reviewbundle.SessionPayload) map[string]reviewbundle.SessionPayload {
	out := map[string]reviewbundle.SessionPayload{}
	for _, session := range sessions {
		out[session.ID] = session
	}
	return out
}

func buildTranscriptChunk(bundle reviewbundle.Bundle, target reviewTarget, source reviewbundle.ReviewTranscriptSource, sessions map[string]reviewbundle.SessionPayload, maxBytes int) (Chunk, bool) {
	session, ok := sessions[source.SessionID]
	if !ok {
		return Chunk{}, false
	}
	request, response, ok := findRequestResponse(session, source.RequestID, source.ResponseID)
	if !ok {
		return Chunk{}, false
	}
	text := renderTranscriptChunk(bundle, target, session, request, response)
	text = limitBytes(text, maxBytes)
	agent := agentProvenanceForSession(target.ReviewContext, source.SessionID, source.Provider, stringValue(source.Model))
	sourceID := transcriptChunkKey(target.revisionKey(), source)
	metadata := TranscriptMetadata{
		SourceID:         sourceID,
		RepoRoot:         bundle.Repo.RootPath,
		BranchName:       target.BranchName,
		RevisionID:       target.revisionKey(),
		CommitID:         target.CommitID,
		RevisionTitle:    target.Description,
		SessionID:        source.SessionID,
		RequestID:        source.RequestID,
		ResponseID:       stringValue(source.ResponseID),
		AgentTool:        agent.AgentTool,
		Provider:         source.Provider,
		Model:            stringValue(source.Model),
		CreatedAt:        source.CreatedAt,
		ProvenanceStatus: source.Status,
		Text:             text,
	}
	return Chunk{
		ID:         stableChunkID(sourceID),
		Text:       text,
		Attributes: metadata.Attributes(),
	}, true
}

func agentProvenanceForSession(context *reviewbundle.ReviewContextPayload, sessionID, provider, model string) reviewbundle.ReviewAgentProvenance {
	if context == nil {
		return reviewbundle.ReviewAgentProvenance{}
	}
	for _, item := range context.AgentProvenance {
		if item.SessionID != sessionID {
			continue
		}
		if item.Provider == provider && item.ModelID == model {
			return item
		}
	}
	for _, item := range context.AgentProvenance {
		if item.SessionID == sessionID {
			return item
		}
	}
	return reviewbundle.ReviewAgentProvenance{}
}

func findRequestResponse(session reviewbundle.SessionPayload, requestID string, responseID *string) (reviewbundle.RequestPayload, *reviewbundle.ResponsePayload, bool) {
	for _, request := range session.Requests {
		if request.ID != requestID {
			continue
		}
		if responseID == nil || *responseID == "" {
			return request, nil, true
		}
		for _, response := range request.Responses {
			if response.ID == *responseID {
				resp := response
				return request, &resp, true
			}
		}
		return request, nil, true
	}
	return reviewbundle.RequestPayload{}, nil, false
}

func renderTranscriptChunk(bundle reviewbundle.Bundle, target reviewTarget, session reviewbundle.SessionPayload, request reviewbundle.RequestPayload, response *reviewbundle.ResponsePayload) string {
	var b strings.Builder
	fmt.Fprintf(&b, "lgtm session transcript for review.\n")
	fmt.Fprintf(&b, "Repo: %s\n", bundle.Repo.RootPath)
	fmt.Fprintf(&b, "Revision: %s\n", target.Description)
	fmt.Fprintf(&b, "lgtm revision: %s\n", target.revisionKey())
	if target.BranchName != "" {
		fmt.Fprintf(&b, "Branch: %s\n", target.BranchName)
	}
	fmt.Fprintf(&b, "Session: %s\n", session.ID)
	fmt.Fprintf(&b, "Request: %s %s %s\n", request.Provider, request.Method, request.Endpoint)
	if request.Model != nil && *request.Model != "" {
		fmt.Fprintf(&b, "Model: %s\n", *request.Model)
	}
	if len(request.RequestBody) > 0 {
		fmt.Fprintf(&b, "\nRequest body:\n%s\n", string(request.RequestBody))
	}
	if response != nil && len(response.ResponseBody) > 0 {
		fmt.Fprintf(&b, "\nResponse body:\n%s\n", string(response.ResponseBody))
	}
	return b.String()
}

func transcriptChunkKey(changeID string, source reviewbundle.ReviewTranscriptSource) string {
	responseID := ""
	if source.ResponseID != nil {
		responseID = *source.ResponseID
	}
	return strings.Join([]string{
		"lgtm",
		"change", changeID,
		"session", source.SessionID,
		"request", source.RequestID,
		"response", responseID,
	}, ":")
}

func stableChunkID(key string) string {
	sum := sha256.Sum256([]byte(key))
	return "lgtm-" + hex.EncodeToString(sum[:])[:40]
}

// limitBytes truncates on a rune boundary. Source files are UTF-8 and a naive
// byte slice can cut a multi-byte rune in half, which the JSON encoder then
// rewrites as U+FFFD — corrupting the indexed text silently rather than
// failing.
func limitBytes(text string, maxBytes int) string {
	if len(text) <= maxBytes {
		return text
	}
	const marker = "\n[truncated]\n"
	if maxBytes <= len(marker) {
		return truncateAtRuneBoundary(text, maxBytes)
	}
	return truncateAtRuneBoundary(text, maxBytes-len(marker)) + marker
}

func truncateAtRuneBoundary(text string, limit int) string {
	if limit <= 0 {
		return ""
	}
	if limit >= len(text) {
		return text
	}
	for limit > 0 && !utf8.RuneStart(text[limit]) {
		limit--
	}
	return text[:limit]
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
