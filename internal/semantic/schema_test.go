package semantic

import (
	"fmt"
	"strings"
	"testing"
)

func TestTurboPufferSchemaDeclaresHybridSearchFields(t *testing.T) {
	schema := transcriptTurboPufferSchema(3072)

	vector, ok := schema["vector"].(map[string]any)
	if !ok || vector["type"] != "[3072]f32" || vector["ann"] != true {
		t.Fatalf("vector schema = %#v", schema["vector"])
	}

	// `text` is the embedded body: stemming on, so a query for "retry" reaches
	// a chunk that says "retries".
	text, ok := schema[transcriptFieldText].(map[string]any)
	if !ok {
		t.Fatalf("text schema = %#v", schema[transcriptFieldText])
	}
	textFTS, ok := text["full_text_search"].(map[string]any)
	if !ok || textFTS["stemming"] != true {
		t.Fatalf("text full_text_search = %#v, want a stemmed configuration", text["full_text_search"])
	}

	// `symbol` is the identifier field: stemming off, because identifier
	// lookups must be exact. Sub-word recall comes from writing pre-split word
	// parts into the value, not from the tokenizer.
	symbol, ok := schema[codeFieldSymbol].(map[string]any)
	if !ok {
		t.Fatalf("symbol schema = %#v", schema[codeFieldSymbol])
	}
	symbolFTS, ok := symbol["full_text_search"].(map[string]any)
	if !ok || symbolFTS["stemming"] != false {
		t.Fatalf("symbol full_text_search = %#v, want an unstemmed configuration", symbol["full_text_search"])
	}
}

func TestTurboPufferSchemaKeepsEveryAttributeWritersEmit(t *testing.T) {
	schema := transcriptTurboPufferSchema(1536)

	filterable := []string{
		// Code rows.
		codeFieldOrgID, codeFieldRepoFullName, codeFieldFilePath, codeFieldFile,
		codeFieldSymbolName, codeFieldSymbolKind, codeFieldPackage,
		codeFieldLanguage, codeFieldDocType, codeFieldChunkHash,
		codeFieldIndexedReason, codeFieldHeadSha,
		// Session/transcript rows that share the namespace.
		transcriptFieldSourceKind, transcriptFieldSourceID, transcriptFieldRepoRoot,
		transcriptFieldBranchName, transcriptFieldRevisionID, transcriptFieldJJChangeID,
		transcriptFieldCommitID, transcriptFieldSessionID, transcriptFieldRequestID,
		transcriptFieldResponseID, transcriptFieldAgentTool, transcriptFieldProvider,
		transcriptFieldModel, transcriptFieldProvenanceStatus,
	}
	for _, field := range filterable {
		spec, ok := schema[field].(map[string]any)
		if !ok || spec["type"] != "string" || spec["filterable"] != true {
			t.Fatalf("%s schema = %#v, want a filterable string", field, schema[field])
		}
	}

	for _, field := range []string{codeFieldStartLine, codeFieldEndLine, codeFieldIndexedAt, transcriptFieldCreatedAt} {
		spec, ok := schema[field].(map[string]any)
		if !ok || spec["type"] != "uint" {
			t.Fatalf("%s schema = %#v, want uint", field, schema[field])
		}
	}
}

func TestCodeRowAttributesMatchTheDeclaredSchema(t *testing.T) {
	schema := transcriptTurboPufferSchema(1536)
	chunks := ChunkSourceFile("internal/semantic/codeindex.go", "package semantic\n\nfunc IndexRepository() {}\n")
	if len(chunks) == 0 {
		t.Fatal("no chunks")
	}
	row := codeChunkRow(codeRowContext{
		RepoRoot:     "/repo",
		RepoFullName: "acme/widgets",
		OrgID:        "org-1",
		CommitID:     "abc",
		BranchName:   "main",
		Reason:       "review",
	}, chunks[0], 0)

	for key := range row.Attributes {
		if _, ok := schema[key]; !ok {
			t.Fatalf("attribute %q is written but not declared in the schema", key)
		}
	}
	for _, required := range []string{
		transcriptFieldSourceKind, transcriptFieldText, codeFieldSymbol,
		codeFieldFilePath, codeFieldStartLine, codeFieldEndLine,
	} {
		if _, ok := row.Attributes[required]; !ok {
			t.Fatalf("code row is missing required attribute %q", required)
		}
	}
}

func TestMaxEmbeddingDimensionsKnowsTheConfiguredModels(t *testing.T) {
	cases := map[string]int{
		"text-embedding-3-large": 3072,
		"text-embedding-3-small": 1536,
		"some-local-model":       0,
	}
	for model, want := range cases {
		if got := MaxEmbeddingDimensions(model); got != want {
			t.Fatalf("MaxEmbeddingDimensions(%q) = %d, want %d", model, got, want)
		}
	}
}

func TestCodeIndexConfigFromEnvCapsDimensionsToTheModel(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "key")
	t.Setenv("TURBOPUFFER_API_KEY", "key")
	t.Setenv("TOTALITY_OPENAI_EMBEDDING_MODEL", "text-embedding-3-small")
	t.Setenv("TOTALITY_EMBEDDING_DIMENSIONS", "3072")
	cfg := CodeIndexConfigFromEnv()
	if cfg.EmbeddingDimensions != 1536 {
		t.Fatalf("dimensions = %d, want the model's native 1536", cfg.EmbeddingDimensions)
	}
	if cfg.TurboPufferNamespace != "" {
		t.Fatalf("namespace = %q, want empty (resolved per repository)", cfg.TurboPufferNamespace)
	}

	t.Setenv("TOTALITY_OPENAI_EMBEDDING_MODEL", "")
	t.Setenv("TOTALITY_EMBEDDING_DIMENSIONS", "")
	def := CodeIndexConfigFromEnv()
	if def.OpenAIEmbeddingModel != defaultOpenAIEmbedModel || def.EmbeddingDimensions != defaultEmbeddingDims {
		t.Fatalf("default config = %s/%d", def.OpenAIEmbeddingModel, def.EmbeddingDimensions)
	}
	if !def.Enabled {
		t.Fatal("code indexing must be enabled by default when credentials exist")
	}

	t.Setenv("TOTALITY_SEMANTIC_INDEX", "0")
	if CodeIndexConfigFromEnv().Enabled {
		t.Fatal("TOTALITY_SEMANTIC_INDEX=0 must turn code indexing off")
	}
}

func TestVectorRowDimensionGuardMentionsBothWidths(t *testing.T) {
	client := NewTurboPufferClientForNamespace(Config{EmbeddingDimensions: 3072}, "totality-scratch")
	err := client.Upsert(t.Context(), []VectorRow{{ID: "x", Vector: make([]float32, 512)}})
	if err == nil {
		t.Fatal("expected a dimension mismatch error")
	}
	want := fmt.Sprintf("vector dimension %d, want %d", 512, 3072)
	if got := err.Error(); !strings.Contains(got, want) {
		t.Fatalf("error = %q, want it to contain %q", got, want)
	}
}
