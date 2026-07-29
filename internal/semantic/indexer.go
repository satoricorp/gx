package semantic

import (
	"context"
	"fmt"

	"github.com/satoricorp/totality/internal/reviewbundle"
)

type IndexResult struct {
	Enabled       bool
	Indexed       bool
	Chunks        int
	SessionChunks int
	CodeChunks    int
}

type Indexer struct {
	cfg      Config
	embedder Embedder
	store    VectorStore
}

func NewIndexer(cfg Config, embedder Embedder, store VectorStore) *Indexer {
	return &Indexer{cfg: cfg, embedder: embedder, store: store}
}

func NewLocalIndexerFromEnv() (*Indexer, Config, error) {
	cfg, err := ConfigFromEnv()
	if err != nil || !cfg.Enabled {
		return nil, cfg, err
	}
	return NewIndexer(cfg, NewOpenAIEmbedder(cfg), NewTurboPufferClient(cfg)), cfg, nil
}

func (i *Indexer) IndexBundle(ctx context.Context, bundle reviewbundle.Bundle) (IndexResult, error) {
	if i == nil || !i.cfg.Enabled {
		return IndexResult{}, nil
	}
	sessionChunks := BuildSessionChunks(bundle, i.cfg.MaxChunkBytes)
	codeChunks := BuildRepositoryChunks(bundle)
	chunks := append(sessionChunks, codeChunks...)
	if len(chunks) == 0 {
		return IndexResult{Enabled: true}, nil
	}
	for start := 0; start < len(chunks); start += i.cfg.BatchSize {
		end := start + i.cfg.BatchSize
		if end > len(chunks) {
			end = len(chunks)
		}
		if err := i.indexBatch(ctx, chunks[start:end]); err != nil {
			return IndexResult{
				Enabled:       true,
				Chunks:        start,
				SessionChunks: len(sessionChunks),
				CodeChunks:    len(codeChunks),
			}, err
		}
	}
	if deleter, ok := i.store.(StaleCodeDocumentDeleter); ok && len(codeChunks) > 0 {
		if err := deleter.DeleteStaleCodeDocuments(ctx, repoFullName(bundle), bundle.Push.HeadCommitID); err != nil {
			return IndexResult{
				Enabled:       true,
				Indexed:       true,
				Chunks:        len(chunks),
				SessionChunks: len(sessionChunks),
				CodeChunks:    len(codeChunks),
			}, err
		}
	}
	return IndexResult{
		Enabled:       true,
		Indexed:       true,
		Chunks:        len(chunks),
		SessionChunks: len(sessionChunks),
		CodeChunks:    len(codeChunks),
	}, nil
}

func (i *Indexer) indexBatch(ctx context.Context, chunks []Chunk) error {
	inputs := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		inputs = append(inputs, chunk.Text)
	}
	embeddings, err := i.embedder.Embed(ctx, inputs)
	if err != nil {
		return err
	}
	if len(embeddings) != len(chunks) {
		return fmt.Errorf("embedder returned %d embeddings for %d chunks", len(embeddings), len(chunks))
	}
	rows := make([]VectorRow, 0, len(chunks))
	for index, chunk := range chunks {
		rows = append(rows, VectorRow{
			ID:         chunk.ID,
			Vector:     embeddings[index],
			Attributes: chunk.Attributes,
		})
	}
	return i.store.Upsert(ctx, rows)
}
