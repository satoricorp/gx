package semantic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// CodeChunkerVersion is bumped whenever chunk boundaries, chunk ids, or the
// embedded text change in a way that invalidates an index built by an older
// build. RepoIndexState records it, so a bump forces a full re-embed instead of
// silently mixing chunk shapes inside one namespace.
const CodeChunkerVersion = 5

const (
	// codeTargetChunkLines is the soft ceiling for a merged unit. Declarations
	// larger than this are windowed; runs of small declarations are fused up to
	// it so a namespace is not flooded with three-line const chunks.
	codeTargetChunkLines = 120
	// codeOverlapLines is how much context a windowed continuation repeats, so
	// a symbol that straddles a window boundary is retrievable from both sides.
	codeOverlapLines = 20
	// codeMergeBelowLines is the size under which a declaration is considered
	// "small" and eligible to fuse with its neighbour.
	codeMergeBelowLines = 30
	// codeMaxChunkBytes caps the embedded text of a single chunk.
	codeMaxChunkBytes = 8000
)

// CodeChunk is one indexable slice of a source file. Line numbers are 1-based
// and inclusive so they can be rendered directly as file:line references.
type CodeChunk struct {
	FilePath    string
	Language    string
	DocType     string
	Package     string
	Symbol      string
	SymbolKind  string
	Symbols     []string
	StartLine   int
	EndLine     int
	Body        string
	ContentHash string
	// Context is the file-level description shared by every chunk of the same
	// file. It is embedded, not stored, and gives a chunk the vocabulary of the
	// file that contains it.
	Context FileContext
	// Card marks the one synthetic chunk per file that describes the file as a
	// whole rather than any span of it.
	Card bool
}

// SymbolText is the value written to the full-text `symbol` attribute.
//
// TurboPuffer's tokenizer does not split camelCase, so indexing the raw
// identifier alone makes BM25 an exact-string lookup: searching "revision" can
// never match "renderRevisionLine". Expanding each identifier into its word
// parts alongside the verbatim spelling keeps exact lookups working and makes
// sub-word lexical search possible. This field is not embedded, so the
// expansion costs nothing in vector quality.
func (c CodeChunk) SymbolText() string {
	var terms []string
	seen := map[string]struct{}{}
	add := func(value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		terms = append(terms, value)
	}
	for _, symbol := range c.Symbols {
		add(symbol)
		for _, part := range SplitIdentifier(symbol) {
			add(part)
		}
	}
	add(c.Package)
	base := path.Base(c.FilePath)
	add(base)
	stem := strings.TrimSuffix(base, path.Ext(base))
	add(stem)
	for _, part := range SplitIdentifier(stem) {
		add(part)
	}
	for _, segment := range strings.Split(path.Dir(c.FilePath), "/") {
		if segment != "" && segment != "." {
			add(segment)
		}
	}
	return strings.Join(terms, " ")
}

// Text renders what actually gets embedded: a header describing where the chunk
// sits and what its file is for, followed by the verbatim source body.
//
// The header carries three kinds of signal that raw source does not:
//
//   - Provenance (path, package, symbol) so a chunk is reachable by name.
//   - The same identifiers rendered as words ("resolve change set"), because
//     questions are asked in prose and an embedding sees a camelCase token as
//     something quite far from the words a reader would use for it.
//   - The file's own doc comments and declaration outline, so every chunk
//     inherits the vocabulary of the file's job. A single declaration rarely
//     restates why its file exists, which is exactly what a conceptual question
//     asks about.
//
// The byte budget applies to the body alone. Letting a rich header eat into the
// body would trade the thing being described for the description of it.
func (c CodeChunk) Text(repoFullName string) string {
	if c.Card {
		return c.cardText(repoFullName)
	}
	var b strings.Builder
	if strings.TrimSpace(repoFullName) != "" {
		fmt.Fprintf(&b, "repo: %s\n", repoFullName)
	}
	fmt.Fprintf(&b, "file: %s\n", c.FilePath)
	if words := pathWords(c.FilePath); words != "" {
		fmt.Fprintf(&b, "path words: %s\n", words)
	}
	fmt.Fprintf(&b, "lines: %d-%d\n", c.StartLine, c.EndLine)
	fmt.Fprintf(&b, "language: %s\n", c.Language)
	if c.Package != "" {
		fmt.Fprintf(&b, "package: %s\n", c.Package)
	}
	if c.SymbolKind != "" {
		fmt.Fprintf(&b, "kind: %s\n", c.SymbolKind)
	}
	if c.Symbol != "" {
		fmt.Fprintf(&b, "symbol: %s\n", c.Symbol)
	}
	if len(c.Symbols) > 1 {
		fmt.Fprintf(&b, "symbols: %s\n", strings.Join(c.Symbols, ", "))
	}
	if words := wordsOf(c.Symbols...); words != "" {
		fmt.Fprintf(&b, "symbol words: %s\n", words)
	}
	if len(c.Context.Outline) > 0 {
		fmt.Fprintf(&b, "file declares: %s\n", strings.Join(c.Context.Outline, ", "))
	}
	if c.Context.Summary != "" {
		fmt.Fprintf(&b, "file summary: %s\n", c.Context.Summary)
	}
	fmt.Fprintf(&b, "doc_type: %s\n---\n", c.DocType)
	b.WriteString(limitBytes(c.Body, codeMaxChunkBytes))
	return b.String()
}

// cardText renders the file card: one row per file that describes the file as a
// unit instead of describing a span of it.
//
// It exists because file-level questions were being answered at chunk level. A
// question like "what stops a low value finding from reaching the user" is
// answered jointly by a file's declarations, and no single declaration's
// embedding is close to it — so the file could only be found when one of its
// chunks happened to be lexically lucky. The card gives every file one
// representative vector built from the file's own prose, which is the form the
// question is asked in. It carries no source body: a body would pull the vector
// back toward code and duplicate a real chunk.
func (c CodeChunk) cardText(repoFullName string) string {
	var b strings.Builder
	if strings.TrimSpace(repoFullName) != "" {
		fmt.Fprintf(&b, "repo: %s\n", repoFullName)
	}
	fmt.Fprintf(&b, "file: %s\n", c.FilePath)
	if words := pathWords(c.FilePath); words != "" {
		fmt.Fprintf(&b, "path words: %s\n", words)
	}
	fmt.Fprintf(&b, "language: %s\n", c.Language)
	if c.Package != "" {
		fmt.Fprintf(&b, "package: %s\n", c.Package)
	}
	fmt.Fprintf(&b, "lines: %d-%d\n", c.StartLine, c.EndLine)
	fmt.Fprintf(&b, "doc_type: %s\n---\n", c.DocType)
	fmt.Fprintf(&b, "This file is %s.\n", c.FilePath)
	// Declarations before prose, and the declaration names repeated as words.
	// The reverse layout (prose first, no word expansion) was measured and was
	// worse on every metric — recall@5 0.875 against 0.938, MRR 0.724 against
	// 0.761 — because the word expansion of a file's declaration names is
	// itself a natural-language description of what the file does, and it is
	// usually longer and more specific than the file's doc comments.
	if len(c.Context.Outline) > 0 {
		fmt.Fprintf(&b, "It declares: %s.\n", strings.Join(c.Context.Outline, ", "))
		if words := wordsOf(c.Context.Outline...); words != "" {
			fmt.Fprintf(&b, "In words: %s.\n", words)
		}
	}
	if c.Context.FullSummary != "" {
		fmt.Fprintf(&b, "What it is for: %s\n", c.Context.FullSummary)
	}
	return b.String()
}

// codeUnit is a candidate chunk expressed as a 1-based inclusive line span.
type codeUnit struct {
	start   int
	end     int
	symbol  string
	kind    string
	symbols []string
}

func (u codeUnit) lines() int {
	if u.end < u.start {
		return 0
	}
	return u.end - u.start + 1
}

// ChunkSourceFile splits one source file into semantic chunks.
//
// Go is parsed with go/ast, so boundaries land exactly on declarations and
// symbol names are exact rather than guessed. Other languages use a
// declaration scanner, and anything unrecognised falls back to overlapping
// line windows. Every strategy produces a total, gap-free cover of the file:
// each unit runs up to the start of the next one, so trailing comments and
// stray statements are never dropped.
func ChunkSourceFile(filePath, source string) []CodeChunk {
	if strings.TrimSpace(source) == "" {
		return nil
	}
	source = strings.ReplaceAll(source, "\r\n", "\n")
	lines := strings.Split(source, "\n")
	// A trailing newline yields a final empty element that is not a real line.
	if len(lines) > 1 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return nil
	}
	language := languageFromPath(filePath)
	docType := docTypeFromPath(filePath)

	units, pkg := codeUnitsFor(language, source, len(lines))
	// The file context is harvested from the pre-merge units, so a doc comment
	// is still attributed to the declaration it introduces.
	fileContext := buildFileContext(lines, units, language)
	units = mergeCodeUnits(units, codeTargetChunkLines, codeMergeBelowLines)
	units = splitCodeUnits(units, codeTargetChunkLines, codeOverlapLines)

	chunks := make([]CodeChunk, 0, len(units))
	for _, unit := range units {
		body := strings.Join(lines[unit.start-1:unit.end], "\n")
		if strings.TrimSpace(body) == "" {
			continue
		}
		chunk := CodeChunk{
			FilePath:   filePath,
			Language:   language,
			DocType:    docType,
			Package:    pkg,
			Symbol:     unit.symbol,
			SymbolKind: unit.kind,
			Symbols:    unit.symbols,
			StartLine:  unit.start,
			EndLine:    unit.end,
			Body:       body,
			Context:    fileContext,
		}
		// The file context is part of what gets embedded, so it must be part of
		// the chunk's identity. Without it, editing a doc comment in one
		// declaration would leave every other chunk of the file holding a stale
		// summary that incremental indexing considers up to date.
		chunk.ContentHash = shortHash(strings.Join([]string{
			filePath,
			chunk.Symbol,
			chunk.SymbolKind,
			fmt.Sprintf("%d-%d", chunk.StartLine, chunk.EndLine),
			fileContext.Summary,
			strings.Join(fileContext.Outline, ","),
			body,
		}, "\x00"))
		chunks = append(chunks, chunk)
	}
	if card, ok := fileCardChunk(filePath, language, docType, pkg, len(lines), fileContext, len(chunks)); ok {
		// Appended last on purpose: the card then occupies chunk index N, so
		// the positional manifest, the row-id derivation and the stale-row
		// sweep all cover it with no special case.
		chunks = append(chunks, card)
	}
	return chunks
}

// fileCardChunk builds the one synthetic per-file chunk, or reports that the
// file does not warrant one.
//
// A file with a single chunk is skipped: its one chunk already carries the same
// context header, so a card would be a near-duplicate row competing with the
// chunk that holds the actual source. A file with no prose and no declarations
// is skipped because the card would say nothing its path does not.
func fileCardChunk(filePath, language, docType, pkg string, totalLines int, ctx FileContext, chunkCount int) (CodeChunk, bool) {
	if chunkCount < 2 || ctx.Empty() {
		return CodeChunk{}, false
	}
	card := CodeChunk{
		FilePath:   filePath,
		Language:   language,
		DocType:    docType,
		Package:    pkg,
		Symbol:     path.Base(filePath),
		SymbolKind: "file",
		Symbols:    ctx.Outline,
		StartLine:  1,
		EndLine:    totalLines,
		Context:    ctx,
		Card:       true,
	}
	card.ContentHash = shortHash(strings.Join([]string{
		filePath, "card", ctx.FullSummary, strings.Join(ctx.Outline, ","),
	}, "\x00"))
	return card, true
}

func codeUnitsFor(language, source string, totalLines int) ([]codeUnit, string) {
	if language == "go" {
		if units, pkg, ok := goCodeUnits(source, totalLines); ok {
			return units, pkg
		}
	}
	if language == "markdown" {
		return markdownCodeUnits(source, totalLines), ""
	}
	if units, ok := scannedCodeUnits(language, source, totalLines); ok {
		return units, ""
	}
	return []codeUnit{{start: 1, end: totalLines}}, ""
}

// goCodeUnits uses the real Go parser. Doc comments are pulled into their
// declaration's span, and each declaration runs to the start of the next one so
// no line of the file is skipped.
func goCodeUnits(source string, totalLines int) ([]codeUnit, string, bool) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "chunk.go", source, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil || file == nil {
		return nil, "", false
	}
	pkg := ""
	if file.Name != nil {
		pkg = file.Name.Name
	}
	lineOf := func(pos token.Pos) int {
		if !pos.IsValid() {
			return 0
		}
		return fset.Position(pos).Line
	}

	type goDecl struct {
		start   int
		symbol  string
		kind    string
		symbols []string
	}
	decls := make([]goDecl, 0, len(file.Decls))
	for _, node := range file.Decls {
		start := lineOf(node.Pos())
		if start <= 0 {
			continue
		}
		item := goDecl{start: start}
		switch decl := node.(type) {
		case *ast.FuncDecl:
			if decl.Doc != nil {
				if line := lineOf(decl.Doc.Pos()); line > 0 && line < item.start {
					item.start = line
				}
			}
			name := ""
			if decl.Name != nil {
				name = decl.Name.Name
			}
			item.symbol = name
			item.kind = "func"
			item.symbols = appendSymbol(item.symbols, name)
			if decl.Recv != nil && len(decl.Recv.List) > 0 {
				if recv := goTypeName(decl.Recv.List[0].Type); recv != "" {
					item.symbol = recv + "." + name
					item.kind = "method"
					item.symbols = appendSymbol(item.symbols, recv)
				}
			}
		case *ast.GenDecl:
			if decl.Doc != nil {
				if line := lineOf(decl.Doc.Pos()); line > 0 && line < item.start {
					item.start = line
				}
			}
			item.kind = strings.ToLower(decl.Tok.String())
			for _, spec := range decl.Specs {
				switch node := spec.(type) {
				case *ast.TypeSpec:
					if node.Name != nil {
						item.symbols = appendSymbol(item.symbols, node.Name.Name)
					}
					switch typ := node.Type.(type) {
					case *ast.InterfaceType:
						item.kind = "interface"
						item.symbols = appendSymbol(item.symbols, goFieldNames(typ.Methods)...)
					case *ast.StructType:
						item.kind = "type"
						item.symbols = appendSymbol(item.symbols, goFieldNames(typ.Fields)...)
					}
				case *ast.ValueSpec:
					for _, name := range node.Names {
						if name != nil {
							item.symbols = appendSymbol(item.symbols, name.Name)
						}
					}
				case *ast.ImportSpec:
					item.kind = "import"
				}
			}
			if len(item.symbols) > 0 {
				item.symbol = item.symbols[0]
			}
		}
		decls = append(decls, item)
	}
	sort.SliceStable(decls, func(i, j int) bool { return decls[i].start < decls[j].start })

	var units []codeUnit
	if len(decls) == 0 {
		return []codeUnit{{start: 1, end: totalLines, symbol: pkg, kind: "file"}}, pkg, true
	}
	if decls[0].start > 1 {
		units = append(units, codeUnit{start: 1, end: decls[0].start - 1, symbol: pkg, kind: "package"})
	}
	for index, decl := range decls {
		end := totalLines
		if index+1 < len(decls) {
			end = decls[index+1].start - 1
		}
		if end < decl.start {
			continue
		}
		units = append(units, codeUnit{
			start:   decl.start,
			end:     end,
			symbol:  decl.symbol,
			kind:    decl.kind,
			symbols: decl.symbols,
		})
	}
	return units, pkg, true
}

func goTypeName(expr ast.Expr) string {
	switch node := expr.(type) {
	case *ast.Ident:
		return node.Name
	case *ast.StarExpr:
		return goTypeName(node.X)
	case *ast.IndexExpr:
		return goTypeName(node.X)
	case *ast.IndexListExpr:
		return goTypeName(node.X)
	case *ast.SelectorExpr:
		if node.Sel != nil {
			return node.Sel.Name
		}
	}
	return ""
}

func goFieldNames(fields *ast.FieldList) []string {
	if fields == nil {
		return nil
	}
	var out []string
	for _, field := range fields.List {
		for _, name := range field.Names {
			if name != nil {
				out = append(out, name.Name)
			}
		}
	}
	return out
}

func appendSymbol(dst []string, values ...string) []string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || value == "_" {
			continue
		}
		found := false
		for _, existing := range dst {
			if existing == value {
				found = true
				break
			}
		}
		if !found {
			dst = append(dst, value)
		}
	}
	return dst
}

var markdownHeadingPattern = regexp.MustCompile(`^(#{1,6})\s+(.+?)\s*$`)

func markdownCodeUnits(source string, totalLines int) []codeUnit {
	lines := strings.Split(source, "\n")
	var starts []codeUnit
	fenced := false
	for index, line := range lines {
		if index >= totalLines {
			break
		}
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			fenced = !fenced
			continue
		}
		if fenced {
			continue
		}
		match := markdownHeadingPattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		heading := strings.TrimSpace(match[2])
		starts = append(starts, codeUnit{
			start:   index + 1,
			symbol:  heading,
			kind:    "heading",
			symbols: appendSymbol(appendSymbol(nil, heading), strings.Fields(heading)...),
		})
	}
	return unitsFromStarts(starts, totalLines)
}

// scannedCodeUnits finds declaration boundaries for languages without a
// bundled parser. It is deliberately conservative: only lines that look like a
// top-level or one-level-nested declaration open a new unit, so a stray match
// inside a function body cannot shred the file into noise.
func scannedCodeUnits(language, source string, totalLines int) ([]codeUnit, bool) {
	matchers, ok := codeSymbolMatchers[language]
	if !ok {
		return nil, false
	}
	lines := strings.Split(source, "\n")
	var starts []codeUnit
	for index, line := range lines {
		if index >= totalLines {
			break
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if indent > 4 {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || isCommentLine(trimmed) {
			continue
		}
		for _, matcher := range matchers {
			match := matcher.pattern.FindStringSubmatch(trimmed)
			if len(match) < 2 {
				continue
			}
			starts = append(starts, codeUnit{
				start:   index + 1,
				symbol:  match[1],
				kind:    matcher.kind,
				symbols: appendSymbol(nil, match[1]),
			})
			break
		}
	}
	if len(starts) == 0 {
		return []codeUnit{{start: 1, end: totalLines}}, true
	}
	return unitsFromStarts(starts, totalLines), true
}

func unitsFromStarts(starts []codeUnit, totalLines int) []codeUnit {
	if len(starts) == 0 {
		return []codeUnit{{start: 1, end: totalLines}}
	}
	var units []codeUnit
	if starts[0].start > 1 {
		units = append(units, codeUnit{start: 1, end: starts[0].start - 1, kind: "preamble"})
	}
	for index, start := range starts {
		end := totalLines
		if index+1 < len(starts) {
			end = starts[index+1].start - 1
		}
		if end < start.start {
			continue
		}
		start.end = end
		units = append(units, start)
	}
	return units
}

func isCommentLine(trimmed string) bool {
	for _, prefix := range []string{"//", "#", "*", "/*", "--", "<!--"} {
		if strings.HasPrefix(trimmed, prefix) {
			return true
		}
	}
	return false
}

type codeSymbolMatcher struct {
	kind    string
	pattern *regexp.Regexp
}

var codeSymbolMatchers = map[string][]codeSymbolMatcher{
	"typescript": tsSymbolMatchers(),
	"javascript": tsSymbolMatchers(),
	"python": {
		{kind: "class", pattern: regexp.MustCompile(`^class\s+([A-Za-z_]\w*)`)},
		{kind: "func", pattern: regexp.MustCompile(`^(?:async\s+)?def\s+([A-Za-z_]\w*)`)},
	},
	"rust": {
		{kind: "impl", pattern: regexp.MustCompile(`^(?:pub(?:\([^)]*\))?\s+)?impl(?:<[^>]*>)?\s+([A-Za-z_]\w*)`)},
		{kind: "type", pattern: regexp.MustCompile(`^(?:pub(?:\([^)]*\))?\s+)?(?:struct|enum|trait|type)\s+([A-Za-z_]\w*)`)},
		{kind: "func", pattern: regexp.MustCompile(`^(?:pub(?:\([^)]*\))?\s+)?(?:async\s+)?(?:unsafe\s+)?fn\s+([A-Za-z_]\w*)`)},
	},
	"java": {
		{kind: "type", pattern: regexp.MustCompile(`^(?:public\s+|private\s+|protected\s+|abstract\s+|final\s+|static\s+)*(?:class|interface|enum|record)\s+([A-Za-z_]\w*)`)},
		{kind: "method", pattern: regexp.MustCompile(`^(?:public\s+|private\s+|protected\s+|abstract\s+|final\s+|static\s+|synchronized\s+)+[\w<>\[\],.?\s]+\s+([A-Za-z_]\w*)\s*\(`)},
	},
	"ruby": {
		{kind: "class", pattern: regexp.MustCompile(`^(?:class|module)\s+([A-Za-z_][\w:]*)`)},
		{kind: "func", pattern: regexp.MustCompile(`^def\s+([A-Za-z_][\w.?!]*)`)},
	},
}

func tsSymbolMatchers() []codeSymbolMatcher {
	return []codeSymbolMatcher{
		{kind: "func", pattern: regexp.MustCompile(`^(?:export\s+)?(?:default\s+)?(?:async\s+)?function\s*\*?\s*([A-Za-z_$][\w$]*)`)},
		{kind: "class", pattern: regexp.MustCompile(`^(?:export\s+)?(?:default\s+)?(?:abstract\s+)?class\s+([A-Za-z_$][\w$]*)`)},
		{kind: "type", pattern: regexp.MustCompile(`^(?:export\s+)?(?:declare\s+)?(?:interface|type|enum)\s+([A-Za-z_$][\w$]*)`)},
		{kind: "const", pattern: regexp.MustCompile(`^(?:export\s+)?(?:declare\s+)?(?:const|let|var)\s+([A-Za-z_$][\w$]*)\s*[:=]`)},
	}
}

// mergeCodeUnits fuses runs of small declarations so the index is not dominated
// by three-line chunks that carry no context, while leaving units that are
// already substantial on their own untouched.
func mergeCodeUnits(units []codeUnit, targetLines, mergeBelow int) []codeUnit {
	if len(units) <= 1 {
		return units
	}
	out := make([]codeUnit, 0, len(units))
	current := units[0]
	for _, next := range units[1:] {
		combined := next.end - current.start + 1
		// Both sides must be small. Fusing a small declaration into a large one
		// would bury the large declaration's boundary and cost the chunk its
		// file:line precision, which review renders directly.
		smallEnough := current.lines() < mergeBelow && next.lines() < mergeBelow
		if smallEnough && combined <= targetLines && next.start == current.end+1 {
			current.end = next.end
			current.symbols = appendSymbol(current.symbols, next.symbols...)
			switch {
			case symbolKindRank(next.kind) > symbolKindRank(current.kind):
				// A real declaration outranks a package clause or import
				// block, so the merged chunk is named after the declaration a
				// reader would search for.
				current.symbol = next.symbol
				current.kind = next.kind
			case current.symbol == "":
				current.symbol = next.symbol
				current.kind = next.kind
			case next.symbol != "" && next.kind != current.kind &&
				symbolKindRank(next.kind) == symbolKindRank(current.kind):
				current.kind = "group"
			}
			continue
		}
		out = append(out, current)
		current = next
	}
	return append(out, current)
}

// symbolKindRank orders declaration kinds by how well they name a chunk. A
// package clause or import block is scaffolding; a function or type is what a
// reader searches for.
func symbolKindRank(kind string) int {
	switch kind {
	case "func", "method", "class", "impl":
		return 3
	case "type", "interface", "struct":
		return 2
	case "const", "var", "heading":
		return 1
	default:
		return 0
	}
}

// splitCodeUnits windows any unit that is still oversized, repeating overlap
// lines so a declaration spanning a boundary stays retrievable from either
// window. Every window keeps the parent unit's symbols.
func splitCodeUnits(units []codeUnit, maxLines, overlap int) []codeUnit {
	if maxLines <= 0 {
		return units
	}
	if overlap < 0 || overlap >= maxLines {
		overlap = maxLines / 4
	}
	out := make([]codeUnit, 0, len(units))
	for _, unit := range units {
		if unit.lines() <= maxLines {
			out = append(out, unit)
			continue
		}
		for start := unit.start; start <= unit.end; {
			end := start + maxLines - 1
			if end > unit.end {
				end = unit.end
			}
			window := unit
			window.start = start
			window.end = end
			out = append(out, window)
			if end >= unit.end {
				break
			}
			start = end - overlap + 1
		}
	}
	return out
}

// SplitIdentifier breaks camelCase, PascalCase, snake_case, kebab-case and
// SCREAMING_CASE identifiers into lowercase word parts. Acronym runs are kept
// whole ("HTTPServer" -> http, server).
func SplitIdentifier(identifier string) []string {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return nil
	}
	var parts []string
	var current strings.Builder
	flush := func() {
		if current.Len() == 0 {
			return
		}
		part := strings.ToLower(current.String())
		current.Reset()
		if len(part) < 2 {
			return
		}
		for _, existing := range parts {
			if existing == part {
				return
			}
		}
		parts = append(parts, part)
	}
	runes := []rune(identifier)
	for index, r := range runes {
		switch {
		case r == '_' || r == '-' || r == '.' || r == '/' || r == ' ':
			flush()
		case unicode.IsUpper(r):
			// Start a new word on lower->upper, and on the last capital of an
			// acronym run that is followed by a lowercase letter (HTTPServer).
			if index > 0 && (unicode.IsLower(runes[index-1]) || unicode.IsDigit(runes[index-1])) {
				flush()
			} else if index > 0 && index+1 < len(runes) && unicode.IsUpper(runes[index-1]) && unicode.IsLower(runes[index+1]) {
				flush()
			}
			current.WriteRune(r)
		default:
			current.WriteRune(r)
		}
	}
	flush()
	if len(parts) == 1 && strings.EqualFold(parts[0], identifier) {
		return nil
	}
	return parts
}

func hashBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
