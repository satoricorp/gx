package semantic

import (
	"context"
	"os/exec"
	"strings"
	"testing"
)

const contextFixture = `package review

import "fmt"

// Modes name what a review actually looked at. They are the honest answer to
// "what did tl read?".
const (
	ModeWorkingTree = "working-tree"
	ModeNone        = "none"
)

// ChangeSet is the resolved answer to "what is this review looking at?".
type ChangeSet struct {
	Mode  string
	Files []string
	Base  string
}

// resolveChangeSet picks the source of changes for one review.
func resolveChangeSet(base string) ChangeSet {
	fmt.Println(base)
	return ChangeSet{Mode: ModeWorkingTree}
}
`

func TestFileContextHarvestsProseAndDeclarationNames(t *testing.T) {
	chunks := ChunkSourceFile("internal/review/changes.go", contextFixture)
	if len(chunks) == 0 {
		t.Fatal("no chunks")
	}
	ctx := chunks[0].Context

	for _, want := range []string{
		"Modes name what a review actually looked at.",
		`ChangeSet is the resolved answer to "what is this review looking at?".`,
		"resolveChangeSet picks the source of changes for one review.",
	} {
		if !strings.Contains(ctx.Summary, want) {
			t.Errorf("summary missing %q\ngot: %s", want, ctx.Summary)
		}
	}

	// The outline names declarations, never the fields inside them. Field names
	// are generic and, repeated on every chunk of the file, they crowd out the
	// prose the context exists to carry.
	outline := strings.Join(ctx.Outline, " ")
	for _, want := range []string{"ChangeSet", "resolveChangeSet", "ModeWorkingTree"} {
		if !strings.Contains(outline, want) {
			t.Errorf("outline missing %q, got %v", want, ctx.Outline)
		}
	}
	for _, unwanted := range []string{"Files", "Base"} {
		for _, symbol := range ctx.Outline {
			if symbol == unwanted {
				t.Errorf("outline contains struct field %q: %v", unwanted, ctx.Outline)
			}
		}
	}
}

func TestChunkTextCarriesFileContextAndWords(t *testing.T) {
	chunks := ChunkSourceFile("internal/review/changes.go", contextFixture)
	text := chunks[0].Text("acme/widgets")
	for _, want := range []string{
		"path words: internal review changes",
		"file declares:",
		"file summary:",
		"symbol words:",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("chunk text missing %q\n%s", want, text)
		}
	}
	// The body budget must apply to the body alone: a richer header may never
	// push source out of the embedded text.
	if !strings.Contains(text, "ModeWorkingTree = \"working-tree\"") {
		t.Errorf("chunk text lost its body\n%s", text)
	}
}

// bigContextFixture is long enough to produce more than one chunk, which is the
// condition for emitting a card.
var bigContextFixture = contextFixture + `
// detectReviewBase walks the candidate branches until one exists.
func detectReviewBase(candidates []string) string {
` + strings.Repeat("\tfmt.Println(\"probing a candidate branch for the review base\")\n", 140) + `	return ""
}
`

func TestFileCardIsTheLastChunkAndDescribesTheWholeFile(t *testing.T) {
	chunks := ChunkSourceFile("internal/review/changes.go", bigContextFixture)
	if len(chunks) < 2 {
		t.Fatalf("fixture produced %d chunks, need at least 2 for a card", len(chunks))
	}
	card := chunks[len(chunks)-1]
	if !card.Card {
		t.Fatalf("last chunk is not the file card: %+v", card)
	}
	if card.StartLine != 1 {
		t.Errorf("card starts at line %d, want 1", card.StartLine)
	}
	if card.SymbolKind != "file" {
		t.Errorf("card kind = %q, want file", card.SymbolKind)
	}
	// Exactly one card per file, and it must be last so its row id is simply
	// the next chunk index.
	cards := 0
	for _, chunk := range chunks {
		if chunk.Card {
			cards++
		}
	}
	if cards != 1 {
		t.Errorf("got %d cards, want 1", cards)
	}
	text := card.Text("acme/widgets")
	if !strings.Contains(text, "This file is internal/review/changes.go.") {
		t.Errorf("card text missing its subject\n%s", text)
	}
	if !strings.Contains(text, "What it is for:") {
		t.Errorf("card text missing prose\n%s", text)
	}
	// The card is a description, not a duplicate of the source.
	if strings.Contains(text, "fmt.Println(base)") {
		t.Errorf("card text contains a source body\n%s", text)
	}
}

func TestFileCardSkippedForSingleChunkFiles(t *testing.T) {
	chunks := ChunkSourceFile("tiny.go", "package tiny\n\n// One does nothing.\nfunc One() {}\n")
	for _, chunk := range chunks {
		if chunk.Card {
			t.Fatalf("emitted a card for a %d-chunk file", len(chunks))
		}
	}
}

func TestSplitIdentifierRendersIdentifiersAsWords(t *testing.T) {
	if got := wordsOf("resolveChangeSet"); got != "resolve change set" {
		t.Errorf("wordsOf = %q", got)
	}
	if got := wordsOf("HTTPServerConfig"); got != "http server config" {
		t.Errorf("wordsOf acronym = %q", got)
	}
	// A single all-lowercase segment has no case boundary to split on, so it
	// contributes itself and nothing more.
	if got := pathWords("internal/codereview/judge.go"); got != "internal codereview judge" {
		t.Errorf("pathWords = %q", got)
	}
	if got := pathWords("internal/reviewPlan/sessionLinks.go"); got != "internal reviewplan review plan sessionlinks session links" {
		t.Errorf("pathWords camelCase = %q", got)
	}
}

func TestFirstSentenceKeepsEnumerationsAndAbbreviations(t *testing.T) {
	cases := map[string]string{
		"Picks a base, in priority order: a. explicit b. detected. Trailing.": "Picks a base, in priority order: a. explicit b. detected.",
		"Uses a cache, e.g. the manifest. Then uploads.":                      "Uses a cache, e.g. the manifest.",
		"No trailing period": "No trailing period",
	}
	for input, want := range cases {
		if got := firstSentence(input); got != want {
			t.Errorf("firstSentence(%q) = %q, want %q", input, got, want)
		}
	}
}

// TestIndexRepositoryHonoursVendoredGitAttributes proves a repository can
// declare material as not-its-own-source and have the index agree.
func TestIndexRepositoryHonoursVendoredGitAttributes(t *testing.T) {
	root := t.TempDir()
	gitInit(t, root)
	writeTestFile(t, root, "internal/alpha.go", "package alpha\n\n// Alpha does a thing.\nfunc Alpha() string { return \"a\" }\n")
	writeTestFile(t, root, "corpus/guide.md", "# Guide\n\nA vendored third-party document about reviewing code.\n")
	writeTestFile(t, root, ".gitattributes", "corpus/** linguist-vendored\n")

	store := newRecordingStore()
	opts := testIndexOptions(t, root, &recordingEmbedder{}, store)
	if _, err := IndexRepository(context.Background(), opts); err != nil {
		t.Fatalf("IndexRepository: %v", err)
	}
	for _, row := range store.rows {
		if path, _ := row.Attributes[codeFieldFilePath].(string); strings.HasPrefix(path, "corpus/") {
			t.Errorf("indexed vendored file %q", path)
		}
	}
	if store.rowCount() == 0 {
		t.Fatal("indexed nothing at all")
	}
}

// TestIndexRepositoryDeletesRowsForNewlyExcludedFilesOnRebuild is a regression
// test. A rebuild resets the manifest, and the removal sweep used to read that
// freshly emptied manifest — so rows belonging to files the new run no longer
// indexes were never deleted and stayed in the namespace forever, with nothing
// left to record that they were there.
func TestIndexRepositoryDeletesRowsForNewlyExcludedFilesOnRebuild(t *testing.T) {
	root := t.TempDir()
	gitInit(t, root)
	writeTestFile(t, root, "internal/alpha.go", "package alpha\n\n// Alpha does a thing.\nfunc Alpha() string { return \"a\" }\n")
	writeTestFile(t, root, "corpus/guide.md", "# Guide\n\nA vendored third-party document.\n\n## More\n\nStill vendored.\n")

	store := newRecordingStore()
	opts := testIndexOptions(t, root, &recordingEmbedder{}, store)
	first, err := IndexRepository(context.Background(), opts)
	if err != nil {
		t.Fatalf("first IndexRepository: %v", err)
	}
	var vendoredIDs []string
	for _, row := range store.rows {
		if path, _ := row.Attributes[codeFieldFilePath].(string); strings.HasPrefix(path, "corpus/") {
			vendoredIDs = append(vendoredIDs, row.ID)
		}
	}
	if len(vendoredIDs) == 0 {
		t.Fatalf("fixture did not index the corpus at all: %+v", first)
	}

	// Declare the corpus vendored, then force the rebuild a chunker bump causes.
	writeTestFile(t, root, ".gitattributes", "corpus/** linguist-vendored\n")
	store.reset()
	opts.Full = true
	second, err := IndexRepository(context.Background(), opts)
	if err != nil {
		t.Fatalf("second IndexRepository: %v", err)
	}
	if second.FilesRemoved == 0 {
		t.Errorf("rebuild reported no removed files: %+v", second)
	}
	gone := map[string]bool{}
	for _, id := range store.deleted {
		gone[id] = true
	}
	for _, id := range vendoredIDs {
		if !gone[id] {
			t.Errorf("row %s for a newly vendored file was never deleted (orphaned)", id)
		}
	}
}

func gitInit(t *testing.T, root string) {
	t.Helper()
	for _, args := range [][]string{
		{"init", "-q"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "test"},
	} {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Skipf("git unavailable: %v (%s)", err, out)
		}
	}
}
