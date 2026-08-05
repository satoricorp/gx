package semantic

import (
	"fmt"
	"strings"
	"testing"
)

func TestChunkSourceFileSplitsGoOnDeclarationBoundaries(t *testing.T) {
	source := `package widgets

import "fmt"

// Alpha does alpha things.
func Alpha() string {
	return "alpha"
}

// Beta is a big function that must not share a chunk with Alpha.
func Beta(input string) string {
` + strings.Repeat("\tfmt.Println(input)\n", 60) + `	return input
}

type Gadget struct {
	Name string
	Size int
}

func (g *Gadget) Describe() string {
	return fmt.Sprintf("%s/%d", g.Name, g.Size)
}
`
	chunks := ChunkSourceFile("internal/widgets/gadget.go", source)
	if len(chunks) < 3 {
		t.Fatalf("chunks = %d, want at least 3 declaration chunks", len(chunks))
	}
	for _, chunk := range chunks {
		if chunk.Language != "go" || chunk.Package != "widgets" {
			t.Fatalf("chunk %#v, want go/widgets", chunk)
		}
		if chunk.StartLine < 1 || chunk.EndLine < chunk.StartLine {
			t.Fatalf("chunk lines = %d-%d", chunk.StartLine, chunk.EndLine)
		}
		if chunk.ContentHash == "" {
			t.Fatalf("chunk %s has no content hash", chunk.Symbol)
		}
	}

	bySymbol := map[string]CodeChunk{}
	for _, chunk := range chunks {
		bySymbol[chunk.Symbol] = chunk
	}
	beta, ok := bySymbol["Beta"]
	if !ok {
		t.Fatalf("no chunk for Beta, got %v", chunkSymbols(chunks))
	}
	if beta.SymbolKind != "func" {
		t.Fatalf("Beta kind = %q, want func", beta.SymbolKind)
	}
	if !strings.Contains(beta.Body, "func Beta(") {
		t.Fatalf("Beta chunk body does not contain its declaration: %q", firstLine(beta.Body))
	}

	describe, ok := bySymbol["Gadget.Describe"]
	if !ok {
		t.Fatalf("no chunk for the Gadget.Describe method, got %v", chunkSymbols(chunks))
	}
	if describe.SymbolKind != "method" {
		t.Fatalf("Describe kind = %q, want method", describe.SymbolKind)
	}
	if !containsString(describe.Symbols, "Gadget") {
		t.Fatalf("method chunk symbols = %v, want the receiver type", describe.Symbols)
	}
}

func TestChunkSourceFileCoversEveryLineOfAGoFile(t *testing.T) {
	source := `package cover

// leading comment

const A = 1

func B() {}

// trailing comment after the last declaration
`
	chunks := ChunkSourceFile("cover.go", source)
	if len(chunks) == 0 {
		t.Fatal("no chunks")
	}
	covered := map[int]bool{}
	for _, chunk := range chunks {
		for line := chunk.StartLine; line <= chunk.EndLine; line++ {
			covered[line] = true
		}
	}
	total := len(strings.Split(strings.TrimRight(source, "\n"), "\n"))
	for line := 1; line <= total; line++ {
		if !covered[line] {
			t.Fatalf("line %d of the file is in no chunk (chunks = %v)", line, chunkRanges(chunks))
		}
	}
}

func TestChunkSourceFileExtractsStructAndInterfaceMembers(t *testing.T) {
	source := `package meta

type Store interface {
	Upsert(id string) error
	DeleteRows(ids []string) error
}

type Row struct {
	HeadCommitID string
	BranchName   string
}
`
	chunks := ChunkSourceFile("meta.go", source)
	var all []string
	for _, chunk := range chunks {
		all = append(all, chunk.Symbols...)
	}
	for _, want := range []string{"Store", "Upsert", "DeleteRows", "Row", "HeadCommitID", "BranchName"} {
		if !containsString(all, want) {
			t.Fatalf("symbols %v missing %q", all, want)
		}
	}
}

func TestChunkSourceFileWindowsOversizedDeclarationsWithOverlap(t *testing.T) {
	var b strings.Builder
	b.WriteString("package big\n\nfunc Huge() {\n")
	for i := 0; i < 400; i++ {
		fmt.Fprintf(&b, "\tstep%d()\n", i)
	}
	b.WriteString("}\n")

	chunks := ChunkSourceFile("big.go", b.String())
	if len(chunks) < 3 {
		t.Fatalf("chunks = %d, want the oversized function split into windows", len(chunks))
	}
	var windows []CodeChunk
	for _, chunk := range chunks {
		if chunk.Symbol == "Huge" {
			windows = append(windows, chunk)
		}
	}
	if len(windows) < 3 {
		t.Fatalf("Huge windows = %d, want >= 3", len(windows))
	}
	for i := 1; i < len(windows); i++ {
		previous, current := windows[i-1], windows[i]
		if current.StartLine > previous.EndLine {
			t.Fatalf("window %d starts at %d after previous ended at %d: no overlap", i, current.StartLine, previous.EndLine)
		}
		overlap := previous.EndLine - current.StartLine + 1
		if overlap != codeOverlapLines {
			t.Fatalf("overlap between windows = %d lines, want %d", overlap, codeOverlapLines)
		}
		if current.Symbol != "Huge" {
			t.Fatalf("window %d lost its symbol: %q", i, current.Symbol)
		}
	}
}

func TestChunkSourceFileFallsBackForUnparseableGo(t *testing.T) {
	source := "package broken\n\nfunc Oops( {\n\tthis is not go\n}\n"
	chunks := ChunkSourceFile("broken.go", source)
	if len(chunks) == 0 {
		t.Fatal("unparseable Go produced no chunks; the file would be invisible to review")
	}
}

func TestChunkSourceFileTypeScriptAndMarkdown(t *testing.T) {
	ts := `import { thing } from "./thing";

export function runIncrementalIndex(input: string) {
  return input;
}

export type IndexChunk = {
  id: string;
};
`
	chunks := ChunkSourceFile("server/src/indexing/turbopuffer.ts", ts)
	var symbols []string
	for _, chunk := range chunks {
		symbols = append(symbols, chunk.Symbols...)
	}
	if !containsString(symbols, "runIncrementalIndex") || !containsString(symbols, "IndexChunk") {
		t.Fatalf("typescript symbols = %v", symbols)
	}
	if chunks[0].Language != "typescript" {
		t.Fatalf("language = %q", chunks[0].Language)
	}

	filler := strings.Repeat("prose line\n", 40)
	md := "# Title\n\n" + filler + "## Section One\n\n" + filler +
		"```\n## not a heading\n```\n\n" + filler + "## Section Two\n\n" + filler
	mdChunks := ChunkSourceFile("docs/guide.md", md)
	var headings []string
	for _, chunk := range mdChunks {
		headings = append(headings, chunk.Symbol)
	}
	if !containsString(headings, "Section One") || !containsString(headings, "Section Two") {
		t.Fatalf("markdown headings = %v", headings)
	}
	if containsString(headings, "not a heading") {
		t.Fatalf("a heading inside a fenced code block became a chunk boundary: %v", headings)
	}
	if mdChunks[0].DocType != "architecture_doc" {
		t.Fatalf("markdown doc_type = %q", mdChunks[0].DocType)
	}
}

func TestSplitIdentifierCoversCodeCasing(t *testing.T) {
	cases := map[string][]string{
		"AttachSessionsFromHunkLinks": {"attach", "sessions", "from", "hunk", "links"},
		"renderRevisionLine":          {"render", "revision", "line"},
		"LGTM_REVIEW_JUDGE":       {"lgtm", "review", "judge"},
		"max_upload_attempts":         {"max", "upload", "attempts"},
		"HTTPServerConfig":            {"http", "server", "config"},
		"turbopuffer":                 nil,
	}
	for input, want := range cases {
		got := SplitIdentifier(input)
		if len(got) != len(want) {
			t.Fatalf("SplitIdentifier(%q) = %v, want %v", input, got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("SplitIdentifier(%q) = %v, want %v", input, got, want)
			}
		}
	}
}

func TestSymbolTextExpandsIdentifiersForLexicalSearch(t *testing.T) {
	chunk := CodeChunk{
		FilePath: "internal/semantic/turbopuffer.go",
		Package:  "semantic",
		Symbol:   "TurboPufferClient.DeleteStaleCodeDocuments",
		Symbols:  []string{"DeleteStaleCodeDocuments", "TurboPufferClient"},
	}
	got := chunk.SymbolText()
	// The verbatim identifier must survive: exact-symbol lookup is the highest
	// precision query a reviewer issues.
	for _, want := range []string{
		"DeleteStaleCodeDocuments", "TurboPufferClient",
		"delete", "stale", "code", "documents", "turbo", "puffer", "client",
		"semantic", "turbopuffer.go", "internal",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("SymbolText() = %q, missing %q", got, want)
		}
	}
}

func TestCodeChunkTextCarriesProvenanceHeader(t *testing.T) {
	chunk := CodeChunk{
		FilePath:   "internal/semantic/codeindex.go",
		Language:   "go",
		DocType:    "code_chunk",
		Package:    "semantic",
		Symbol:     "IndexRepository",
		SymbolKind: "func",
		Symbols:    []string{"IndexRepository", "RepoIndexOptions"},
		StartLine:  10,
		EndLine:    42,
		Body:       "func IndexRepository() {}",
	}
	text := chunk.Text("satoricorp/lgtm")
	for _, want := range []string{
		"repo: satoricorp/lgtm",
		"file: internal/semantic/codeindex.go",
		"lines: 10-42",
		"language: go",
		"package: semantic",
		"kind: func",
		"symbol: IndexRepository",
		"symbols: IndexRepository, RepoIndexOptions",
		"func IndexRepository() {}",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("Text() = %q, missing %q", text, want)
		}
	}
}

func TestCodeRowIDIsStableAcrossCommits(t *testing.T) {
	first := CodeRowID("satoricorp/lgtm", "internal/semantic/codeindex.go", 3)
	second := CodeRowID("satoricorp/lgtm", "internal/semantic/codeindex.go", 3)
	if first != second || first == "" {
		t.Fatalf("CodeRowID is not stable: %q vs %q", first, second)
	}
	if CodeRowID("satoricorp/lgtm", "internal/semantic/codeindex.go", 4) == first {
		t.Fatal("CodeRowID collides across chunk indexes")
	}
	if CodeRowID("satoricorp/console", "internal/semantic/codeindex.go", 3) == first {
		t.Fatal("CodeRowID collides across repositories")
	}
}

func chunkSymbols(chunks []CodeChunk) []string {
	var out []string
	for _, chunk := range chunks {
		out = append(out, chunk.Symbol)
	}
	return out
}

func chunkRanges(chunks []CodeChunk) []string {
	var out []string
	for _, chunk := range chunks {
		out = append(out, fmt.Sprintf("%d-%d", chunk.StartLine, chunk.EndLine))
	}
	return out
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func firstLine(body string) string {
	if index := strings.IndexByte(body, '\n'); index >= 0 {
		return body[:index]
	}
	return body
}
