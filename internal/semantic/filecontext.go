package semantic

import (
	"path"
	"regexp"
	"strings"
)

const (
	// codeFileSummaryBytes caps the prose harvested from a file's comments.
	// The cap matters: this text is repeated on every chunk of the file, and an
	// unbounded preamble would crowd the chunk's own body out of the embedding.
	codeFileSummaryBytes = 600
	// codeFileOutlineSymbols caps how many declaration names the shared outline
	// carries, for the same reason.
	codeFileOutlineSymbols = 40
	// codeDocSentenceBytes caps one harvested doc comment.
	codeDocSentenceBytes = 160
	// codeFileCardSummaryBytes is the far larger budget for the file card. The
	// card is a single row, not a per-chunk preamble, so there is nothing for a
	// long summary to crowd out — the more of the file's own prose it carries,
	// the better it represents the file.
	codeFileCardSummaryBytes = 3500
)

// FileContext is the file-level material shared by every chunk of one file.
//
// It exists because a chunk embedding only sees its own few dozen lines, while
// the questions asked of the index are about the file's job: "where do we
// decide what a review looks at" is answered by a file whose individual
// declarations never use the words "decide" or "review looks at". Attaching a
// compact, deterministic description of the whole file to each of its chunks
// gives every chunk the vocabulary of its container. This is the same idea as
// LLM-generated contextual retrieval, computed from the source instead, so it
// costs nothing at index time and cannot hallucinate.
type FileContext struct {
	// Summary is prose harvested from the file's own comments, in file order:
	// the leading/package comment first, then the first sentence of each
	// declaration's doc comment. It is kept short because it is repeated on
	// every chunk of the file.
	Summary string
	// FullSummary is the same harvest under the much larger budget the file
	// card can afford.
	FullSummary string
	// Outline is the file's declaration names in file order, deduped.
	Outline []string
}

// Empty reports whether the context carries nothing worth embedding.
func (c FileContext) Empty() bool {
	return strings.TrimSpace(c.Summary) == "" && len(c.Outline) == 0
}

// buildFileContext harvests the shared context from a file's lines and the
// units the chunker already computed, so declaration boundaries are found once.
func buildFileContext(lines []string, units []codeUnit, language string) FileContext {
	var ctx FileContext
	seen := map[string]struct{}{}
	for _, unit := range units {
		// Only the declaration's own name, never the names nested inside it.
		// unit.symbols also carries struct fields and interface methods, which
		// earn their place in the lexical `symbol` column but are poison here:
		// a file of small structs contributes a wall of generic tokens (ID,
		// Text, File, client, reason) that is repeated on every chunk of the
		// file and crowds out the prose summary this context exists to carry.
		symbol := strings.TrimSpace(unit.symbol)
		if symbol == "" {
			continue
		}
		if _, ok := seen[symbol]; ok {
			continue
		}
		seen[symbol] = struct{}{}
		if len(ctx.Outline) < codeFileOutlineSymbols {
			ctx.Outline = append(ctx.Outline, symbol)
		}
	}

	var prose []string
	total := 0
	appendProse := func(text string) bool {
		text = strings.TrimSpace(collapseSpace(text))
		if text == "" {
			return true
		}
		text = firstSentence(text)
		if len(text) > codeDocSentenceBytes {
			text = truncateAtRuneBoundary(text, codeDocSentenceBytes)
		}
		if total+len(text) > codeFileCardSummaryBytes {
			return false
		}
		prose = append(prose, text)
		total += len(text) + 1
		return true
	}

	// The leading comment block of the file: a Go package doc, a TS/JS file
	// banner, a Python module docstring.
	if lead := leadingComment(lines, language); lead != "" {
		appendProse(lead)
	}
	for _, unit := range units {
		if unit.start < 1 || unit.start > len(lines) {
			continue
		}
		doc := unitDocComment(lines, unit, language)
		if doc == "" {
			continue
		}
		if !appendProse(doc) {
			break
		}
	}
	ctx.FullSummary = strings.Join(prose, " ")
	// The per-chunk summary is the same prose truncated to its own budget, so
	// the two never disagree about what the file is for.
	ctx.Summary = joinWithinBudget(prose, codeFileSummaryBytes)
	return ctx
}

func joinWithinBudget(parts []string, maxBytes int) string {
	var kept []string
	total := 0
	for _, part := range parts {
		if total+len(part) > maxBytes {
			break
		}
		kept = append(kept, part)
		total += len(part) + 1
	}
	return strings.Join(kept, " ")
}

// leadingComment returns the comment block at the very top of a file.
func leadingComment(lines []string, language string) string {
	index := 0
	for index < len(lines) && strings.TrimSpace(lines[index]) == "" {
		index++
	}
	text, _ := commentBlockAt(lines, index, language)
	return text
}

// unitDocComment returns the doc comment that introduces a unit. Comment-first
// languages (Go, TS, Rust, Java) put it immediately at the unit start, because
// the chunker folds the doc comment into the declaration's span. Python puts a
// docstring on the line after the `def`, so both positions are tried.
func unitDocComment(lines []string, unit codeUnit, language string) string {
	if text, _ := commentBlockAt(lines, unit.start-1, language); text != "" {
		return text
	}
	if language == "python" && unit.start < len(lines) {
		if text, _ := commentBlockAt(lines, unit.start, language); text != "" {
			return text
		}
	}
	return ""
}

var (
	lineCommentMarkers = map[string][]string{
		"go":         {"//"},
		"typescript": {"//"},
		"javascript": {"//"},
		"rust":       {"///", "//!", "//"},
		"java":       {"//"},
		"python":     {"#"},
		"ruby":       {"#"},
		"shell":      {"#"},
	}
	docstringFence = regexp.MustCompile(`^\s*("""|''')`)
)

// commentBlockAt reads a contiguous comment block starting at line index
// `start` (0-based) and returns its text with markers stripped, plus the index
// just past the block. An empty result means the line is not a comment.
func commentBlockAt(lines []string, start int, language string) (string, int) {
	if start < 0 || start >= len(lines) {
		return "", start
	}
	markers, ok := lineCommentMarkers[language]
	if !ok {
		markers = []string{"//", "#"}
	}
	trimmed := strings.TrimSpace(lines[start])

	// Triple-quoted docstrings (Python) and block comments (/** ... */).
	if fence := docstringFence.FindString(lines[start]); fence != "" {
		return readFencedBlock(lines, start, strings.TrimSpace(fence))
	}
	if strings.HasPrefix(trimmed, "/*") {
		return readBlockComment(lines, start)
	}

	var parts []string
	index := start
	for index < len(lines) {
		line := strings.TrimSpace(lines[index])
		marker := ""
		for _, candidate := range markers {
			if strings.HasPrefix(line, candidate) {
				marker = candidate
				break
			}
		}
		if marker == "" {
			break
		}
		body := strings.TrimSpace(strings.TrimPrefix(line, marker))
		// A divider line ("// ----") carries no prose.
		if body != "" && strings.Trim(body, "-=*_#/ ") != "" {
			parts = append(parts, body)
		}
		index++
	}
	if len(parts) == 0 {
		return "", start
	}
	return strings.Join(parts, " "), index
}

func readFencedBlock(lines []string, start int, fence string) (string, int) {
	first := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[start]), fence))
	if strings.HasSuffix(first, fence) {
		return strings.TrimSpace(strings.TrimSuffix(first, fence)), start + 1
	}
	parts := []string{}
	if first != "" {
		parts = append(parts, first)
	}
	for index := start + 1; index < len(lines); index++ {
		line := lines[index]
		if strings.Contains(line, fence) {
			tail := strings.TrimSpace(line[:strings.Index(line, fence)])
			if tail != "" {
				parts = append(parts, tail)
			}
			return strings.Join(parts, " "), index + 1
		}
		if text := strings.TrimSpace(line); text != "" {
			parts = append(parts, text)
		}
	}
	return strings.Join(parts, " "), len(lines)
}

func readBlockComment(lines []string, start int) (string, int) {
	var parts []string
	for index := start; index < len(lines); index++ {
		line := strings.TrimSpace(lines[index])
		closed := strings.Contains(line, "*/")
		line = strings.TrimPrefix(line, "/**")
		line = strings.TrimPrefix(line, "/*")
		if idx := strings.Index(line, "*/"); idx >= 0 {
			line = line[:idx]
		}
		line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "*"))
		if line != "" && strings.Trim(line, "-=*_ ") != "" {
			parts = append(parts, line)
		}
		if closed {
			return strings.Join(parts, " "), index + 1
		}
	}
	return strings.Join(parts, " "), len(lines)
}

// firstSentence keeps a doc comment to its opening statement, which is the part
// that describes what the declaration is for. Later sentences are usually
// caveats and would spend the shared budget without adding topic signal.
func firstSentence(text string) string {
	for index := 0; index < len(text)-1; index++ {
		if text[index] != '.' {
			continue
		}
		next := text[index+1]
		if next == ' ' || next == '\t' {
			// Skip abbreviations that are not sentence ends.
			word := lastWord(text[:index])
			// "e.g." and friends, and enumerations ("in priority order: a. ..."),
			// are not sentence ends.
			if word == "e.g" || word == "i.e" || word == "etc" || word == "cf" || len(word) == 1 {
				continue
			}
			return text[:index+1]
		}
	}
	return text
}

func lastWord(text string) string {
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return ""
	}
	return fields[len(fields)-1]
}

var spaceRun = regexp.MustCompile(`\s+`)

func collapseSpace(text string) string {
	return strings.TrimSpace(spaceRun.ReplaceAllString(text, " "))
}

// wordsOf renders identifiers as the natural-language phrase they encode:
// "resolveChangeSet" becomes "resolve change set". Questions are asked in
// prose, and an embedding sees a camelCase identifier as an opaque token far
// from the words a reader would use for it.
func wordsOf(identifiers ...string) string {
	var words []string
	seen := map[string]struct{}{}
	for _, identifier := range identifiers {
		for _, part := range SplitIdentifier(identifier) {
			if _, ok := seen[part]; ok {
				continue
			}
			seen[part] = struct{}{}
			words = append(words, part)
		}
	}
	return strings.Join(words, " ")
}

// pathWords renders a file path as its word parts, so "internal/codereview" is
// reachable from the words "code review".
func pathWords(filePath string) string {
	var words []string
	seen := map[string]struct{}{}
	add := func(value string) {
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" || value == "." {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		words = append(words, value)
	}
	base := path.Base(filePath)
	stem := strings.TrimSuffix(base, path.Ext(base))
	for _, segment := range strings.Split(path.Dir(filePath), "/") {
		add(segment)
		for _, part := range SplitIdentifier(segment) {
			add(part)
		}
	}
	add(stem)
	for _, part := range SplitIdentifier(stem) {
		add(part)
	}
	return strings.Join(words, " ")
}
