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
	// Order matters twice over: it fixes the output order of the composite, and
	// the AI context budget is applied with a stable sort, so a source listed
	// earlier survives a priority tie. The repository's own code and the
	// sessions that wrote it come before external guidance, because a finding
	// about this change is grounded in this repository first.
	if retriever := codeIndexRetrieverFromEnv(); retriever != nil {
		retrievers = append(retrievers, retriever)
	}
	if retriever := sessionContextRetrieverFromEnv(); retriever != nil {
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

// Retrieve runs every sub-retriever concurrently and returns everything they
// found.
//
// A sub-retriever's failure does not fail the review: one unreachable index
// should not stop a review that has a diff and a repository to read. It is,
// however, recorded. This used to be a bare `return` inside the goroutine,
// which meant a review whose evidence sources were all down produced exactly
// the same output as a fully informed one, with nothing in the report — or in
// the brief handed to the model — able to tell the two apart.
//
// Retrievers that know their own namespaces record finer-grained statuses
// themselves; the fallback here covers the ones that only return an error.
func (r CompositeContextRetriever) Retrieve(ctx context.Context, in RetrieveInput) ([]ContextSnippet, error) {
	results := make([][]ContextSnippet, len(r.Retrievers))
	errs := make([]error, len(r.Retrievers))
	reported := make([]bool, len(r.Retrievers))
	var wg sync.WaitGroup
	for i, retriever := range r.Retrievers {
		i, retriever := i, retriever
		if retriever == nil {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Each retriever writes through its own scope of the shared log, so
			// "did this retriever report for itself?" is answered by what it
			// filed rather than by whatever the retrievers running alongside it
			// happened to file first.
			scoped := in
			scoped.Evidence = in.Evidence.scope()
			snippets, err := retriever.Retrieve(ctx, scoped)
			reported[i] = scoped.Evidence.recorded() > 0
			if err != nil {
				errs[i] = err
				return
			}
			results[i] = snippets
		}()
	}
	wg.Wait()

	var out []ContextSnippet
	for i, snippets := range results {
		out = append(out, snippets...)
		retriever := r.Retrievers[i]
		if retriever == nil {
			continue
		}
		if errs[i] != nil {
			if !reported[i] {
				in.Evidence.Record(EvidenceStatus{
					Source: evidenceSourceName(retriever),
					State:  EvidenceUnavailable,
					Detail: errs[i].Error(),
				})
			}
			continue
		}
		if reported[i] {
			continue
		}
		if _, named := retriever.(evidenceNamer); !named {
			// Local retrieval, and anything else that does not name a source,
			// gets no evidence line: reading the checkout cannot be
			// "unavailable" the way an index can, so a line about it would be
			// noise rather than a signal.
			continue
		}
		in.Evidence.Record(EvidenceStatus{
			Source:   evidenceSourceName(retriever),
			State:    evidenceStateForCount(len(snippets)),
			Snippets: len(snippets),
		})
	}
	return dedupeContextSnippets(out), nil
}

func evidenceStateForCount(count int) string {
	if count > 0 {
		return EvidenceOK
	}
	return EvidenceEmpty
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
