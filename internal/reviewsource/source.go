package reviewsource

type ProvenanceSource struct {
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
	Source    string `json:"source"`
}

type TranscriptSource struct {
	SessionID  string  `json:"session_id"`
	RequestID  string  `json:"request_id"`
	ResponseID *string `json:"response_id,omitempty"`
	Provider   string  `json:"provider,omitempty"`
	Model      *string `json:"model,omitempty"`
	CreatedAt  int64   `json:"created_at"`
	Status     string  `json:"status"`
	Source     string  `json:"source"`
}

type Graph struct {
	ProvenanceStatus   string
	LinkedSessionCount int
	ProvenanceSources  []ProvenanceSource
	TranscriptSources  []TranscriptSource
}

func BuildGraph(evidenceStatuses []string, sessionIDs []string, transcriptSources []TranscriptSource) Graph {
	status := provenanceStatus(evidenceStatuses, len(sessionIDs))
	return Graph{
		ProvenanceStatus:   status,
		LinkedSessionCount: len(sessionIDs),
		ProvenanceSources:  provenanceSources(sessionIDs, status),
		TranscriptSources:  normalizeTranscriptSources(transcriptSources, status),
	}
}

func provenanceStatus(evidenceStatuses []string, linkedSessionCount int) string {
	if len(evidenceStatuses) == 0 {
		if linkedSessionCount > 0 {
			return "linked"
		}
		return "absent"
	}
	status := "explicit"
	for _, item := range evidenceStatuses {
		switch item {
		case "absent":
			return "absent"
		case "repo_local":
			status = "repo_local"
		case "":
			if status == "explicit" {
				status = "unknown"
			}
		}
	}
	return status
}

func provenanceSources(sessionIDs []string, status string) []ProvenanceSource {
	if len(sessionIDs) == 0 {
		return nil
	}
	out := make([]ProvenanceSource, 0, len(sessionIDs))
	for _, sessionID := range sessionIDs {
		if sessionID == "" {
			continue
		}
		out = append(out, ProvenanceSource{
			SessionID: sessionID,
			Status:    status,
			Source:    "change_sessions",
		})
	}
	return out
}

func normalizeTranscriptSources(sources []TranscriptSource, status string) []TranscriptSource {
	out := make([]TranscriptSource, 0, len(sources))
	for _, source := range sources {
		if source.SessionID == "" || source.RequestID == "" {
			continue
		}
		source.Status = status
		if source.Source == "" {
			source.Source = "change_sessions"
		}
		out = append(out, source)
	}
	return out
}
