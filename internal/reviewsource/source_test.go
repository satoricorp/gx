package reviewsource

import "testing"

func TestBuildGraphUsesDemuxProvenanceStatus(t *testing.T) {
	graph := BuildGraph(
		[]string{"explicit"},
		[]string{"session-one"},
		[]TranscriptSource{{SessionID: "session-one", RequestID: "request-one"}},
	)
	if graph.ProvenanceStatus != "explicit" {
		t.Fatalf("status = %q, want explicit", graph.ProvenanceStatus)
	}
	if len(graph.ProvenanceSources) != 1 || graph.ProvenanceSources[0].Status != "explicit" {
		t.Fatalf("provenance sources = %#v", graph.ProvenanceSources)
	}
	if len(graph.TranscriptSources) != 1 || graph.TranscriptSources[0].Status != "explicit" || graph.TranscriptSources[0].Source != "change_sessions" {
		t.Fatalf("transcript sources = %#v", graph.TranscriptSources)
	}
}

func TestBuildGraphUsesLinkedForNonDemuxLinkedSessions(t *testing.T) {
	graph := BuildGraph(nil, []string{"session-one"}, nil)
	if graph.ProvenanceStatus != "linked" {
		t.Fatalf("status = %q, want linked", graph.ProvenanceStatus)
	}
	if graph.LinkedSessionCount != 1 {
		t.Fatalf("linked session count = %d, want 1", graph.LinkedSessionCount)
	}
}

func TestBuildGraphUsesAbsentWithoutEvidenceOrSessions(t *testing.T) {
	graph := BuildGraph(nil, nil, nil)
	if graph.ProvenanceStatus != "absent" {
		t.Fatalf("status = %q, want absent", graph.ProvenanceStatus)
	}
}
