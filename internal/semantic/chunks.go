package semantic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/satoricorp/gx/internal/reviewbundle"
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
		if target.Change.ReviewContext == nil {
			continue
		}
		for _, source := range target.Change.ReviewContext.TranscriptSources {
			key := transcriptChunkKey(target.Change.JJChangeID, source)
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
	Change     reviewbundle.ChangePayload
	BranchName string
}

func reviewTargets(bundle reviewbundle.Bundle) []reviewTarget {
	var targets []reviewTarget
	for _, entry := range bundle.Stack {
		targets = append(targets, reviewTarget{
			Change:     entry.Change,
			BranchName: entry.BranchName,
		})
	}
	if len(targets) == 0 && bundle.Change != nil {
		branchName := ""
		if bundle.Repo.BranchName != nil {
			branchName = *bundle.Repo.BranchName
		}
		targets = append(targets, reviewTarget{
			Change:     *bundle.Change,
			BranchName: branchName,
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
	agent := agentProvenanceForSession(target.Change.ReviewContext, source.SessionID, source.Provider, stringValue(source.Model))
	sourceID := transcriptChunkKey(target.Change.JJChangeID, source)
	metadata := TranscriptMetadata{
		SourceID:         sourceID,
		RepoRoot:         bundle.Repo.RootPath,
		BranchName:       target.BranchName,
		JJChangeID:       target.Change.JJChangeID,
		CommitID:         target.Change.CurrentCommitID,
		RevisionTitle:    target.Change.Description,
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
	fmt.Fprintf(&b, "GX session transcript for review.\n")
	fmt.Fprintf(&b, "Repo: %s\n", bundle.Repo.RootPath)
	fmt.Fprintf(&b, "Revision: %s\n", target.Change.Description)
	fmt.Fprintf(&b, "JJ change: %s\n", target.Change.JJChangeID)
	if target.BranchName != "" {
		fmt.Fprintf(&b, "Branch: %s\n", target.BranchName)
	}
	fmt.Fprintf(&b, "Session: %s\n", session.ID)
	fmt.Fprintf(&b, "Command: %s\n", session.Command)
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
		"gx",
		"change", changeID,
		"session", source.SessionID,
		"request", source.RequestID,
		"response", responseID,
	}, ":")
}

func stableChunkID(key string) string {
	sum := sha256.Sum256([]byte(key))
	return "gx-" + hex.EncodeToString(sum[:])[:40]
}

func limitBytes(text string, maxBytes int) string {
	if len(text) <= maxBytes {
		return text
	}
	if maxBytes <= len("\n[truncated]\n") {
		return text[:maxBytes]
	}
	return text[:maxBytes-len("\n[truncated]\n")] + "\n[truncated]\n"
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
