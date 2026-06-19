package structural

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type FileFact struct {
	File              string   `json:"file"`
	Language          string   `json:"language"`
	DefinedSymbols    []string `json:"defined_symbols,omitempty"`
	ReferencedSymbols []string `json:"referenced_symbols,omitempty"`
	Symbols           []Symbol `json:"symbols,omitempty"`
}

type Symbol struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	File      string `json:"file"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

type DependencyEdge struct {
	FromFile string `json:"from_file"`
	ToFile   string `json:"to_file"`
	Symbol   string `json:"symbol"`
}

type HunkInput struct {
	ID       string
	File     string
	NewStart int
	NewLines int
}

type HunkSymbol struct {
	HunkID    string `json:"hunk_id"`
	File      string `json:"file"`
	Symbol    string `json:"symbol"`
	Kind      string `json:"kind"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

type Facts struct {
	Files []FileFact       `json:"files,omitempty"`
	Edges []DependencyEdge `json:"edges,omitempty"`
}

var identifierPattern = regexp.MustCompile(`[A-Za-z_$][A-Za-z0-9_$]*`)

type definitionPattern struct {
	kind    string
	pattern *regexp.Regexp
}

var definitionPatterns = map[string][]definitionPattern{
	"go": {
		{kind: "function", pattern: regexp.MustCompile(`^\s*func\s+(?:\([^)]*\)\s*)?([A-Za-z_][A-Za-z0-9_]*)\s*\(`)},
		{kind: "type", pattern: regexp.MustCompile(`^\s*type\s+([A-Za-z_][A-Za-z0-9_]*)\b`)},
	},
	"typescript": {
		{kind: "function", pattern: regexp.MustCompile(`^\s*(?:export\s+)?(?:async\s+)?function\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*\(`)},
		{kind: "type", pattern: regexp.MustCompile(`^\s*(?:export\s+)?(?:class|interface|type|enum)\s+([A-Za-z_$][A-Za-z0-9_$]*)\b`)},
	},
	"python": {
		{kind: "definition", pattern: regexp.MustCompile(`^\s*(?:def|class)\s+([A-Za-z_][A-Za-z0-9_]*)\b`)},
	},
}

var ignoredIdentifiers = map[string]struct{}{
	"and": {}, "any": {}, "as": {}, "async": {}, "await": {}, "bool": {}, "break": {},
	"case": {}, "chan": {}, "class": {}, "const": {}, "continue": {}, "def": {},
	"default": {}, "defer": {}, "else": {}, "enum": {}, "export": {}, "false": {},
	"for": {}, "func": {}, "function": {}, "go": {}, "if": {}, "import": {},
	"in": {}, "interface": {}, "let": {}, "map": {}, "nil": {}, "none": {},
	"package": {}, "range": {}, "return": {}, "select": {}, "string": {}, "struct": {},
	"switch": {}, "true": {}, "type": {}, "var": {},
}

func Analyze(root string, files []string) Facts {
	cleaned := cleanFiles(files)
	facts := make([]FileFact, 0, len(cleaned))
	for _, file := range cleaned {
		language := languageForFile(file)
		if language == "" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(file)))
		if err != nil {
			continue
		}
		text := string(data)
		symbols := parseSymbols(file, language, text)
		facts = append(facts, FileFact{
			File:              file,
			Language:          language,
			DefinedSymbols:    symbolNames(symbols),
			ReferencedSymbols: referencedSymbols(text),
			Symbols:           symbols,
		})
	}
	return Facts{
		Files: facts,
		Edges: dependencyEdges(facts),
	}
}

func languageForFile(file string) string {
	lower := strings.ToLower(file)
	switch {
	case strings.HasSuffix(lower, ".go"):
		return "go"
	case strings.HasSuffix(lower, ".ts"), strings.HasSuffix(lower, ".tsx"), strings.HasSuffix(lower, ".js"), strings.HasSuffix(lower, ".jsx"):
		return "typescript"
	case strings.HasSuffix(lower, ".py"):
		return "python"
	default:
		return ""
	}
}

func parseSymbols(file, language, text string) []Symbol {
	seen := map[string]struct{}{}
	var symbols []Symbol
	lines := strings.Split(text, "\n")
	for index, line := range lines {
		for _, def := range definitionPatterns[language] {
			match := def.pattern.FindStringSubmatch(line)
			if len(match) < 2 {
				continue
			}
			name := strings.TrimSpace(match[1])
			if name == "" {
				continue
			}
			key := def.kind + ":" + name
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			symbols = append(symbols, Symbol{
				Name:      name,
				Kind:      def.kind,
				File:      file,
				StartLine: index + 1,
				EndLine:   len(lines),
			})
			break
		}
	}
	for index := range symbols {
		if index < len(symbols)-1 {
			symbols[index].EndLine = symbols[index+1].StartLine - 1
		}
	}
	sort.SliceStable(symbols, func(i, j int) bool {
		if symbols[i].StartLine != symbols[j].StartLine {
			return symbols[i].StartLine < symbols[j].StartLine
		}
		return symbols[i].Name < symbols[j].Name
	})
	return symbols
}

func symbolNames(symbols []Symbol) []string {
	seen := map[string]struct{}{}
	var names []string
	for _, symbol := range symbols {
		if _, ok := seen[symbol.Name]; ok {
			continue
		}
		seen[symbol.Name] = struct{}{}
		names = append(names, symbol.Name)
	}
	sort.Strings(names)
	return names
}

func EnclosingSymbols(facts Facts, hunks []HunkInput) []HunkSymbol {
	symbolsByFile := map[string][]Symbol{}
	for _, fact := range facts.Files {
		symbolsByFile[fact.File] = fact.Symbols
	}
	var out []HunkSymbol
	for _, hunk := range hunks {
		symbol, ok := enclosingSymbol(symbolsByFile[hunk.File], hunk.NewStart, hunk.NewLines)
		if !ok {
			continue
		}
		out = append(out, HunkSymbol{
			HunkID:    hunk.ID,
			File:      hunk.File,
			Symbol:    symbol.Name,
			Kind:      symbol.Kind,
			StartLine: symbol.StartLine,
			EndLine:   symbol.EndLine,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].File != out[j].File {
			return out[i].File < out[j].File
		}
		if out[i].StartLine != out[j].StartLine {
			return out[i].StartLine < out[j].StartLine
		}
		return out[i].HunkID < out[j].HunkID
	})
	return out
}

func enclosingSymbol(symbols []Symbol, start, lines int) (Symbol, bool) {
	if len(symbols) == 0 || start <= 0 {
		return Symbol{}, false
	}
	end := start
	if lines > 0 {
		end = start + lines - 1
	}
	for _, symbol := range symbols {
		if symbol.StartLine >= start && symbol.StartLine <= end {
			return symbol, true
		}
	}
	for _, symbol := range symbols {
		if symbol.StartLine <= start && symbol.EndLine >= end {
			return symbol, true
		}
	}
	for _, symbol := range symbols {
		if symbol.StartLine <= end && symbol.EndLine >= start {
			return symbol, true
		}
	}
	return Symbol{}, false
}

func referencedSymbols(text string) []string {
	seen := map[string]struct{}{}
	var symbols []string
	for _, symbol := range identifierPattern.FindAllString(text, -1) {
		lower := strings.ToLower(symbol)
		if _, ok := ignoredIdentifiers[lower]; ok {
			continue
		}
		if _, ok := seen[symbol]; ok {
			continue
		}
		seen[symbol] = struct{}{}
		symbols = append(symbols, symbol)
	}
	sort.Strings(symbols)
	return symbols
}

func dependencyEdges(facts []FileFact) []DependencyEdge {
	type definitionOwner struct {
		file     string
		language string
	}
	definitions := map[string]definitionOwner{}
	for _, fact := range facts {
		for _, symbol := range fact.DefinedSymbols {
			if _, exists := definitions[symbol]; !exists {
				definitions[symbol] = definitionOwner{file: fact.File, language: fact.Language}
			}
		}
	}
	seen := map[DependencyEdge]struct{}{}
	var edges []DependencyEdge
	for _, fact := range facts {
		for _, symbol := range fact.ReferencedSymbols {
			owner, ok := definitions[symbol]
			if !ok || owner.file == fact.File || owner.language != fact.Language {
				continue
			}
			edge := DependencyEdge{FromFile: fact.File, ToFile: owner.file, Symbol: symbol}
			if _, ok := seen[edge]; ok {
				continue
			}
			seen[edge] = struct{}{}
			edges = append(edges, edge)
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].FromFile != edges[j].FromFile {
			return edges[i].FromFile < edges[j].FromFile
		}
		if edges[i].ToFile != edges[j].ToFile {
			return edges[i].ToFile < edges[j].ToFile
		}
		return edges[i].Symbol < edges[j].Symbol
	})
	return edges
}

func cleanFiles(files []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(files))
	for _, file := range files {
		file = filepath.ToSlash(strings.TrimSpace(file))
		if file == "" {
			continue
		}
		if _, ok := seen[file]; ok {
			continue
		}
		seen[file] = struct{}{}
		out = append(out, file)
	}
	sort.Strings(out)
	return out
}
