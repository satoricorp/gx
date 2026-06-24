package consoleingest

import (
	"testing"

	"github.com/satoricorp/gx/internal/reviewbundle"
)

func TestIngestStackReviewBundle(t *testing.T) {
	responseID := "response-one"
	bundle := reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push:          reviewbundle.PushPayload{HeadCommitID: "commit-head"},
		Stack: []reviewbundle.StackPayload{{
			BranchName:     "feature/alpha",
			BaseBranchName: "main",
			Patch:          "diff --git a/alpha.txt b/alpha.txt\n",
			Change: reviewbundle.ChangePayload{
				JJChangeID:      "change-alpha",
				CurrentCommitID: "commit-alpha",
				Description:     "feat alpha",
				ReviewContext: &reviewbundle.ReviewContextPayload{
					ProvenanceStatus: "explicit",
					StructuralStatus: "available",
					Risk:             reviewbundle.RiskPayload{Level: "high", Score: 80},
					TranscriptSources: []reviewbundle.ReviewTranscriptSource{{
						SessionID:  "session-one",
						RequestID:  "request-one",
						ResponseID: &responseID,
					}},
				},
			},
		}},
		Sessions: []reviewbundle.SessionPayload{{
			ID: "session-one",
			Requests: []reviewbundle.RequestPayload{{
				ID:          "request-one",
				RequestBody: []byte(`{"input":"why alpha?"}`),
				Responses: []reviewbundle.ResponsePayload{{
					ID:           "response-one",
					ResponseBody: []byte(`{"output":"because alpha"}`),
				}},
			}},
		}},
	}

	artifact := reviewbundle.NewArtifact(bundle)
	artifact.ReviewID = "review-one"
	artifact.ReviewURL = "https://gx.test/reviews/review-one"
	artifact.IndexStatus = "indexed"
	review, err := Ingest(artifact)
	if err != nil {
		t.Fatalf("Ingest() error = %v", err)
	}
	if review.ReviewID != "review-one" || review.ReviewURL != "https://gx.test/reviews/review-one" || review.IndexStatus != "indexed" {
		t.Fatalf("review artifact fields = %#v", review)
	}
	if review.RepoRoot != "/repo" || review.HeadCommitID != "commit-head" {
		t.Fatalf("review = %#v", review)
	}
	if len(review.Revisions) != 1 {
		t.Fatalf("revisions = %#v, want one", review.Revisions)
	}
	revision := review.Revisions[0]
	if revision.BranchName != "feature/alpha" || revision.Description != "feat alpha" || revision.ProvenanceStatus != "explicit" {
		t.Fatalf("revision = %#v", revision)
	}
	if !revision.WhyContextAvailable {
		t.Fatalf("revision = %#v, want why context available", revision)
	}
	why, ok := review.WhyContext("change-alpha")
	if !ok {
		t.Fatal("WhyContext(change-alpha) = false, want true")
	}
	if why.Revision.JJChangeID != "change-alpha" {
		t.Fatalf("why revision = %#v", why.Revision)
	}
	if len(why.Sources) != 1 {
		t.Fatalf("why sources = %#v, want one", why.Sources)
	}
	source := why.Sources[0]
	if source.SessionID != "session-one" || source.RequestID != "request-one" || source.ResponseID == nil || *source.ResponseID != "response-one" {
		t.Fatalf("why source ids = %#v", source)
	}
	if source.RequestText != `{"input":"why alpha?"}` || source.ResponseText != `{"output":"because alpha"}` {
		t.Fatalf("why source text = %#v", source)
	}
}

func TestWhyContextFindsCurrentChangeByCommitID(t *testing.T) {
	responseID := "response-one"
	bundle := reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push:          reviewbundle.PushPayload{HeadCommitID: "commit-alpha"},
		Change: &reviewbundle.ChangePayload{
			JJChangeID:      "change-alpha",
			CurrentCommitID: "commit-alpha",
			Description:     "feat alpha",
			ReviewContext: &reviewbundle.ReviewContextPayload{
				ProvenanceStatus: "linked",
				TranscriptSources: []reviewbundle.ReviewTranscriptSource{{
					SessionID:  "session-one",
					RequestID:  "request-one",
					ResponseID: &responseID,
				}},
			},
		},
		Sessions: []reviewbundle.SessionPayload{{
			ID: "session-one",
			Requests: []reviewbundle.RequestPayload{{
				ID:          "request-one",
				RequestBody: []byte("implement alpha"),
				Responses: []reviewbundle.ResponsePayload{{
					ID:           "response-one",
					ResponseBody: []byte("done"),
				}},
			}},
		}},
	}
	review, err := Ingest(reviewbundle.NewArtifact(bundle))
	if err != nil {
		t.Fatalf("Ingest() error = %v", err)
	}

	why, ok := review.WhyContext("commit-alpha")
	if !ok {
		t.Fatal("WhyContext(commit-alpha) = false, want true")
	}
	if len(why.Sources) != 1 || why.Sources[0].RequestText != "implement alpha" || why.Sources[0].ResponseText != "done" {
		t.Fatalf("why = %#v", why)
	}
}

func TestWhyContextReturnsFalseWithoutMatchingTranscript(t *testing.T) {
	bundle := reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Change: &reviewbundle.ChangePayload{
			JJChangeID:      "change-alpha",
			CurrentCommitID: "commit-alpha",
			Description:     "feat alpha",
			ReviewContext: &reviewbundle.ReviewContextPayload{
				TranscriptSources: []reviewbundle.ReviewTranscriptSource{{
					SessionID: "missing-session",
					RequestID: "request-one",
				}},
			},
		},
	}
	review, err := Ingest(reviewbundle.NewArtifact(bundle))
	if err != nil {
		t.Fatalf("Ingest() error = %v", err)
	}
	if _, ok := review.WhyContext("change-alpha"); ok {
		t.Fatal("WhyContext(change-alpha) = true, want false")
	}
}

func TestIngestRejectsUnknownSchemaVersion(t *testing.T) {
	_, err := Ingest(reviewbundle.Artifact{Bundle: reviewbundle.Bundle{SchemaVersion: reviewbundle.SchemaVersion + 1}})
	if err == nil {
		t.Fatal("expected schema version error")
	}
}

func TestIngestBundleWrapsLegacyBundle(t *testing.T) {
	review, err := IngestBundle(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push:          reviewbundle.PushPayload{HeadCommitID: "commit-alpha"},
		Sessions:      []reviewbundle.SessionPayload{},
	})
	if err != nil {
		t.Fatalf("IngestBundle() error = %v", err)
	}
	if review.IndexStatus != "pending" || review.HeadCommitID != "commit-alpha" {
		t.Fatalf("review = %#v, want artifact defaults", review)
	}
}
