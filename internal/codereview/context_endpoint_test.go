package codereview

import (
	"context"
	"errors"
	"testing"
)

func TestCompositeContextRetrieverIgnoresRetrieverFailures(t *testing.T) {
	snippets, err := (CompositeContextRetriever{Retrievers: []ContextRetriever{
		fakeRetriever{snippets: []ContextSnippet{{Kind: "repo_doc", Ref: "README.md", Source: "local", Text: "readme"}}},
		failingRetriever{},
	}}).Retrieve(context.Background(), RetrieveInput{RepoRoot: "/repo"})
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 1 || snippets[0].Source != "local" {
		t.Fatalf("snippets = %#v", snippets)
	}
}

type failingRetriever struct{}

func (failingRetriever) Retrieve(context.Context, RetrieveInput) ([]ContextSnippet, error) {
	return nil, errors.New("boom")
}
