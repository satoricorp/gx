package semantic

import (
	"context"
	"reflect"
	"testing"

	"github.com/satoricorp/totality/internal/reviewbundle"
)

func TestIndexerIndexesBundleChunks(t *testing.T) {
	responseID := "response-one"
	bundle := reviewbundle.Bundle{
		Repo: reviewbundle.RepoPayload{RootPath: "/repo"},
		Revisions: []reviewbundle.RevisionPayload{{
			RevisionID:  "change-one",
			Description: "feat alpha",
			ReviewContext: &reviewbundle.ReviewContextPayload{
				TranscriptSources: []reviewbundle.ReviewTranscriptSource{{
					SessionID:  "session-one",
					RequestID:  "request-one",
					ResponseID: &responseID,
					Provider:   "openai",
					Status:     "explicit",
				}},
			},
		}},
		Sessions: []reviewbundle.SessionPayload{{
			ID: "session-one",
			Requests: []reviewbundle.RequestPayload{{
				ID:             "request-one",
				Provider:       "openai",
				Method:         "POST",
				Endpoint:       "/v1/responses",
				RequestBody:    []byte(`{"input":"alpha"}`),
				RequestHeaders: "{}",
				Responses: []reviewbundle.ResponsePayload{{
					ID:              "response-one",
					ResponseBody:    []byte(`{"output":"done"}`),
					ResponseHeaders: "{}",
				}},
			}},
		}},
	}
	store := &fakeVectorStore{}
	indexer := NewIndexer(
		Config{Enabled: true, BatchSize: 2, MaxChunkBytes: 12000},
		fakeEmbedder{vectors: [][]float32{{0.1, 0.2}}},
		store,
	)

	result, err := indexer.IndexBundle(context.Background(), bundle)
	if err != nil {
		t.Fatalf("IndexBundle() error = %v", err)
	}
	if !result.Enabled || !result.Indexed || result.Chunks != 1 {
		t.Fatalf("result = %#v", result)
	}
	if len(store.rows) != 1 {
		t.Fatalf("rows = %#v, want one", store.rows)
	}
	if !reflect.DeepEqual(store.rows[0].Vector, []float32{0.1, 0.2}) {
		t.Fatalf("vector = %#v", store.rows[0].Vector)
	}
	if store.rows[0].Attributes["revision_id"] != "change-one" {
		t.Fatalf("attributes = %#v", store.rows[0].Attributes)
	}
	// Dual-written legacy attribute until a full reindex retires it.
	if store.rows[0].Attributes["jj_change_id"] != "change-one" {
		t.Fatalf("attributes = %#v", store.rows[0].Attributes)
	}
}

type fakeEmbedder struct {
	vectors [][]float32
}

func (f fakeEmbedder) Embed(_ context.Context, inputs []string) ([][]float32, error) {
	if len(inputs) != len(f.vectors) {
		return nil, nil
	}
	return f.vectors, nil
}

type fakeVectorStore struct {
	rows []VectorRow
}

func (f *fakeVectorStore) Upsert(_ context.Context, rows []VectorRow) error {
	f.rows = append(f.rows, rows...)
	return nil
}
