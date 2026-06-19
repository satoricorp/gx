package consoleingest

import (
	"bytes"
	"fmt"

	"github.com/satoricorp/gx/internal/reviewbundle"
)

type Review struct {
	Event         string
	SchemaVersion int
	ReviewID      string
	ReviewURL     string
	IndexStatus   string
	RepoRoot      string
	HeadCommitID  string
	Revisions     []Revision
	Sessions      map[string]reviewbundle.SessionPayload
}

type Revision struct {
	JJChangeID          string
	CommitID            string
	Description         string
	BranchName          string
	BaseBranchName      string
	Patch               string
	ProvenanceStatus    string
	StructuralStatus    string
	RiskLevel           string
	RiskScore           int
	TranscriptSources   []reviewbundle.ReviewTranscriptSource
	AgentProvenance     []reviewbundle.ReviewAgentProvenance
	WhyContextAvailable bool
}

type WhyContext struct {
	Revision Revision
	Sources  []WhySource
}

type WhySource struct {
	SessionID    string
	RequestID    string
	ResponseID   *string
	Provider     string
	Model        *string
	CreatedAt    int64
	Status       string
	RequestText  string
	ResponseText string
}

func Ingest(artifact reviewbundle.Artifact) (Review, error) {
	if artifact.SchemaVersion != reviewbundle.SchemaVersion {
		return Review{}, fmt.Errorf("unsupported review artifact schema version %d", artifact.SchemaVersion)
	}
	review := Review{
		Event:         artifact.Event,
		SchemaVersion: artifact.SchemaVersion,
		ReviewID:      artifact.ReviewID,
		ReviewURL:     artifact.ReviewURL,
		IndexStatus:   artifact.IndexStatus,
		RepoRoot:      artifact.Repo.RootPath,
		HeadCommitID:  artifact.Push.HeadCommitID,
		Sessions:      sessionIndex(artifact.Sessions),
	}
	if len(artifact.Stack) > 0 {
		for _, entry := range artifact.Stack {
			review.Revisions = append(review.Revisions, revisionFromStackEntry(entry, review.Sessions))
		}
		return review, nil
	}
	if artifact.Change != nil {
		review.Revisions = append(review.Revisions, revisionFromChange(*artifact.Change, "", "", "", review.Sessions))
	}
	return review, nil
}

func IngestBundle(bundle reviewbundle.Bundle) (Review, error) {
	return Ingest(reviewbundle.NewArtifact(bundle))
}

func (r Review) WhyContext(revisionID string) (WhyContext, bool) {
	revision, ok := r.FindRevision(revisionID)
	if !ok {
		return WhyContext{}, false
	}
	sources := make([]WhySource, 0, len(revision.TranscriptSources))
	for _, source := range revision.TranscriptSources {
		session, ok := r.Sessions[source.SessionID]
		if !ok {
			continue
		}
		request, response, ok := findRequestResponse(session, source.RequestID, source.ResponseID)
		if !ok {
			continue
		}
		item := WhySource{
			SessionID:   source.SessionID,
			RequestID:   source.RequestID,
			ResponseID:  source.ResponseID,
			Provider:    source.Provider,
			Model:       source.Model,
			CreatedAt:   source.CreatedAt,
			Status:      source.Status,
			RequestText: printableBytes(request.RequestBody),
		}
		if response != nil {
			item.ResponseText = printableBytes(response.ResponseBody)
		}
		sources = append(sources, item)
	}
	return WhyContext{Revision: revision, Sources: sources}, len(sources) > 0
}

func (r Review) FindRevision(revisionID string) (Revision, bool) {
	for _, revision := range r.Revisions {
		if revisionID == "" ||
			revision.JJChangeID == revisionID ||
			revision.CommitID == revisionID ||
			revision.BranchName == revisionID ||
			revision.Description == revisionID {
			return revision, true
		}
	}
	return Revision{}, false
}

func sessionIndex(sessions []reviewbundle.SessionPayload) map[string]reviewbundle.SessionPayload {
	out := map[string]reviewbundle.SessionPayload{}
	for _, session := range sessions {
		out[session.ID] = session
	}
	return out
}

func revisionFromStackEntry(entry reviewbundle.StackPayload, sessions map[string]reviewbundle.SessionPayload) Revision {
	return revisionFromChange(entry.Change, entry.BranchName, entry.BaseBranchName, entry.Patch, sessions)
}

func revisionFromChange(change reviewbundle.ChangePayload, branchName, baseBranchName, patch string, sessions map[string]reviewbundle.SessionPayload) Revision {
	revision := Revision{
		JJChangeID:     change.JJChangeID,
		CommitID:       change.CurrentCommitID,
		Description:    change.Description,
		BranchName:     branchName,
		BaseBranchName: baseBranchName,
		Patch:          patch,
	}
	if change.ReviewContext != nil {
		revision.ProvenanceStatus = change.ReviewContext.ProvenanceStatus
		revision.StructuralStatus = change.ReviewContext.StructuralStatus
		revision.RiskLevel = change.ReviewContext.Risk.Level
		revision.RiskScore = change.ReviewContext.Risk.Score
		revision.TranscriptSources = change.ReviewContext.TranscriptSources
		revision.AgentProvenance = change.ReviewContext.AgentProvenance
		revision.WhyContextAvailable = whyContextAvailable(change.ReviewContext.TranscriptSources, sessions)
	}
	return revision
}

func whyContextAvailable(sources []reviewbundle.ReviewTranscriptSource, sessions map[string]reviewbundle.SessionPayload) bool {
	for _, source := range sources {
		session, ok := sessions[source.SessionID]
		if !ok {
			continue
		}
		if _, _, ok := findRequestResponse(session, source.RequestID, source.ResponseID); ok {
			return true
		}
	}
	return false
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

func printableBytes(data []byte) string {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return ""
	}
	return string(data)
}
