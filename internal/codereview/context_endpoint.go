package codereview

import "context"

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

func (r CompositeContextRetriever) Retrieve(ctx context.Context, repoRoot string, opts Options, facts RepoFacts, hints []ReviewHint) ([]ContextSnippet, error) {
	var out []ContextSnippet
	for _, retriever := range r.Retrievers {
		if retriever == nil {
			continue
		}
		snippets, err := retriever.Retrieve(ctx, repoRoot, opts, facts, hints)
		if err != nil {
			continue
		}
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
