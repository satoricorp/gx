package codereview

import (
	"context"
	"sync"
)

type CompositeContextRetriever struct {
	Retrievers []ContextRetriever
}

func contextRetrieverFromEnv() ContextRetriever {
	retrievers := []ContextRetriever{LocalContextRetriever{}}
	if retriever := indexedContextRetrieverFromEnv(); retriever != nil {
		retrievers = append(retrievers, retriever)
	}
	if retriever := reviewHistoryRetrieverFromEnv(); retriever != nil {
		retrievers = append(retrievers, retriever)
	}
	if retriever := reviewResourceRetrieverFromEnv(); retriever != nil {
		retrievers = append(retrievers, retriever)
	}
	return CompositeContextRetriever{Retrievers: retrievers}
}

func (r CompositeContextRetriever) Retrieve(ctx context.Context, in RetrieveInput) ([]ContextSnippet, error) {
	results := make([][]ContextSnippet, len(r.Retrievers))
	var wg sync.WaitGroup
	for i, retriever := range r.Retrievers {
		i, retriever := i, retriever
		if retriever == nil {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			snippets, err := retriever.Retrieve(ctx, in)
			if err != nil {
				return
			}
			results[i] = snippets
		}()
	}
	wg.Wait()

	var out []ContextSnippet
	for _, snippets := range results {
		out = append(out, snippets...)
	}
	return dedupeContextSnippets(out), nil
}

func dedupeContextSnippets(snippets []ContextSnippet) []ContextSnippet {
	seen := map[string]struct{}{}
	var out []ContextSnippet
	for _, snippet := range snippets {
		key := snippet.Kind + "\x00" + snippet.Ref + "\x00" + snippet.Source
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, snippet)
	}
	return out
}
