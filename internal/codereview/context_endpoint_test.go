package codereview

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCompositeContextRetrieverIgnoresRetrieverFailures(t *testing.T) {
	snippets, err := (CompositeContextRetriever{Retrievers: []ContextRetriever{
		fakeRetriever{snippets: []ContextSnippet{{Kind: "repo_doc", Ref: "README.md", Source: "local", Text: "readme"}}},
		failingRetriever{},
	}}).Retrieve(context.Background(), "/repo", Options{}, RepoFacts{}, nil)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 1 || snippets[0].Source != "local" {
		t.Fatalf("snippets = %#v", snippets)
	}
}

func TestCompositeContextRetrieverHandlesFailureAndSuccessConcurrently(t *testing.T) {
	snippets, err := (CompositeContextRetriever{Retrievers: []ContextRetriever{
		delayedContextRetriever{delay: 80 * time.Millisecond, snippets: []ContextSnippet{{Kind: "repo_doc", Ref: "README.md", Source: "local", Text: "readme"}}},
		delayedFailingRetriever{delay: 80 * time.Millisecond},
	}}).Retrieve(context.Background(), "/repo", Options{}, RepoFacts{}, nil)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 1 || snippets[0].Ref != "README.md" {
		t.Fatalf("snippets = %#v", snippets)
	}
}

func TestCompositeContextRetrieverOutputOrderIsDeterministic(t *testing.T) {
	snippets, err := (CompositeContextRetriever{Retrievers: []ContextRetriever{
		delayedContextRetriever{delay: 80 * time.Millisecond, snippets: []ContextSnippet{{Kind: "repo_doc", Ref: "first", Source: "local", Text: "first"}}},
		delayedContextRetriever{delay: 10 * time.Millisecond, snippets: []ContextSnippet{{Kind: "repo_doc", Ref: "second", Source: "local", Text: "second"}}},
		delayedContextRetriever{delay: 5 * time.Millisecond, snippets: []ContextSnippet{{Kind: "repo_doc", Ref: "first", Source: "local", Text: "duplicate"}}},
	}}).Retrieve(context.Background(), "/repo", Options{}, RepoFacts{}, nil)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 2 || snippets[0].Ref != "first" || snippets[1].Ref != "second" {
		t.Fatalf("snippets = %#v, want retriever order with duplicate removed", snippets)
	}
}

type failingRetriever struct{}

func (failingRetriever) Retrieve(context.Context, string, Options, RepoFacts, []ReviewHint) ([]ContextSnippet, error) {
	return nil, errors.New("boom")
}

type delayedContextRetriever struct {
	delay    time.Duration
	snippets []ContextSnippet
}

func (r delayedContextRetriever) Retrieve(ctx context.Context, _ string, _ Options, _ RepoFacts, _ []ReviewHint) ([]ContextSnippet, error) {
	timer := time.NewTimer(r.delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
		return r.snippets, nil
	}
}

type delayedFailingRetriever struct {
	delay time.Duration
}

func (r delayedFailingRetriever) Retrieve(ctx context.Context, _ string, _ Options, _ RepoFacts, _ []ReviewHint) ([]ContextSnippet, error) {
	timer := time.NewTimer(r.delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
		return nil, errors.New("boom")
	}
}
