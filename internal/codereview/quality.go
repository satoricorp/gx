package codereview

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const maxCodeQualityHints = 40

func collectCodeQualityHints(repoRoot string, facts RepoFacts, opts Options) []CodeQualityHint {
	var hints []CodeQualityHint
	for _, file := range facts.Files {
		if !qualityFileSupported(file) || isTestFile(file) || ignorableQualityFile(file) {
			continue
		}
		if opts.Focus != "" && !inFocus(file, opts.Focus) {
			continue
		}
		hints = append(hints, qualityHintsForFile(repoRoot, file)...)
	}
	sort.SliceStable(hints, func(i, j int) bool {
		if qualityRank(hints[i].Kind) != qualityRank(hints[j].Kind) {
			return qualityRank(hints[i].Kind) < qualityRank(hints[j].Kind)
		}
		if hints[i].File != hints[j].File {
			return hints[i].File < hints[j].File
		}
		return hints[i].Line < hints[j].Line
	})
	if len(hints) > maxCodeQualityHints {
		hints = hints[:maxCodeQualityHints]
	}
	return hints
}

func qualityHintsForFile(repoRoot, rel string) []CodeQualityHint {
	data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(rel)))
	if err != nil {
		return nil
	}
	language := qualityLanguage(rel)
	content := string(data)
	hasAsyncDef := strings.Contains(content, "async def ")
	var hints []CodeQualityHint
	for i, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if qualityLineIgnored(trimmed) {
			continue
		}
		lineNo := i + 1
		switch {
		case shellInjectionPattern(language, trimmed):
			hints = append(hints, qualityHint("shell_injection", rel, lineNo, trimmed, "Command execution with shell parsing or string-built commands can turn input into shell syntax; prefer argv arrays and strict allowlists."))
		case unsafeHTMLPattern(language, trimmed):
			hints = append(hints, qualityHint("unsafe_html", rel, lineNo, trimmed, "Direct HTML injection bypasses framework escaping and needs a trusted sanitization boundary."))
		case sqlInterpolationPattern(language, trimmed):
			hints = append(hints, qualityHint("sql_interpolation", rel, lineNo, trimmed, "SQL assembled with interpolation or concatenation is easy to turn into injection or quoting bugs; prefer parameterized queries."))
		case unsafeRegexCompilePattern(language, trimmed):
			hints = append(hints, qualityHint("unsafe_regex", rel, lineNo, trimmed, "Compiling a regular expression from user or request input can enable ReDoS or unexpected pattern behavior; validate and bound patterns first."))
		case hardcodedSecretPattern(language, trimmed):
			hints = append(hints, qualityHint("hardcoded_secret", rel, lineNo, trimmed, "Hardcoded credentials or private key material in source code can leak through repos, logs, and build artifacts."))
		case language == "go" && ignoredNonCleanupResult(trimmed):
			hints = append(hints, qualityHint("ignored_result", rel, lineNo, trimmed, "A discarded result or error in production code can hide malformed input, failed cleanup, or failed recovery paths."))
		case language == "go" && strings.Contains(trimmed, "time.Sleep("):
			hints = append(hints, qualityHint("sleep_polling", rel, lineNo, trimmed, "Sleeping in production code is often a brittle readiness or polling seam unless bounded and tested."))
		case language == "go" && (strings.Contains(trimmed, "log.Println(") || strings.Contains(trimmed, "log.Printf(")):
			hints = append(hints, qualityHint("direct_logging", rel, lineNo, trimmed, "Direct logging from a Module can make error handling and tests less observable than returning classified errors."))
		case language == "go" && strings.Contains(trimmed, "panic("):
			hints = append(hints, qualityHint("panic", rel, lineNo, trimmed, "Panic in production code should be justified by an invariant that callers cannot recover from."))
		case rustPanicPattern(language, trimmed):
			hints = append(hints, qualityHint("panic", rel, lineNo, trimmed, "Unwrap or expect in production Rust should be justified by an invariant that callers cannot recover from."))
		case unsafeBufferPattern(language, trimmed):
			hints = append(hints, qualityHint("unsafe_buffer", rel, lineNo, trimmed, "Unbounded C/C++ string copy/format APIs are memory-safety hazards; prefer bounded alternatives with explicit lengths."))
		case asyncBlockingPattern(language, trimmed, hasAsyncDef):
			hints = append(hints, qualityHint("async_blocking", rel, lineNo, trimmed, "Blocking calls inside async code can stall unrelated work; use async-aware sleep or IO primitives."))
		}
	}
	return hints
}

func qualityHint(kind, file string, line int, text, reason string) CodeQualityHint {
	return CodeQualityHint{Kind: kind, File: file, Line: line, Text: text, Reason: reason}
}

func qualityFileSupported(path string) bool {
	return qualityLanguage(path) != ""
}

func qualityLanguage(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".ts", ".tsx":
		return "typescript"
	case ".js", ".jsx", ".mjs", ".cjs":
		return "javascript"
	case ".rs":
		return "rust"
	case ".sql":
		return "sql"
	case ".java":
		return "java"
	case ".c", ".h":
		return "c"
	case ".cc", ".cpp", ".cxx", ".hh", ".hpp", ".hxx":
		return "cpp"
	default:
		return ""
	}
}

func qualityLineIgnored(line string) bool {
	if line == "" {
		return true
	}
	for _, prefix := range []string{"//", "#", "--", "/*", "*"} {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

func shellInjectionPattern(language, line string) bool {
	lower := strings.ToLower(line)
	switch language {
	case "go":
		if !strings.Contains(line, "exec.Command(") {
			return false
		}
		return strings.Contains(lower, `"bash"`) ||
			strings.Contains(lower, `"sh"`) ||
			strings.Contains(lower, "-c")
	case "python":
		return strings.Contains(lower, "shell=true") || strings.Contains(line, "os.system(") || strings.Contains(line, "os.popen(")
	case "javascript", "typescript":
		return strings.Contains(line, "exec(") || strings.Contains(line, "execSync(") || strings.Contains(lower, "shell: true")
	case "rust", "c", "cpp":
		return strings.Contains(line, "system(")
	case "java":
		return strings.Contains(line, "Runtime.getRuntime().exec(") || strings.Contains(line, "ProcessBuilder(") && strings.Contains(line, "+")
	default:
		return false
	}
}

func unsafeHTMLPattern(language, line string) bool {
	switch language {
	case "javascript", "typescript":
		return strings.Contains(line, "dangerouslySetInnerHTML") || strings.Contains(line, ".innerHTML") || strings.Contains(line, "v-html=")
	case "python":
		return strings.Contains(line, "Markup(") || strings.Contains(line, "|safe")
	case "java":
		return strings.Contains(line, "setEscapeModelStrings(false)")
	default:
		return false
	}
}

func sqlInterpolationPattern(language, line string) bool {
	lower := strings.ToLower(line)
	hasSQL := strings.Contains(lower, "select ") || strings.Contains(lower, "insert ") || strings.Contains(lower, "update ") || strings.Contains(lower, "delete ")
	if !hasSQL {
		return false
	}
	switch language {
	case "python":
		return strings.Contains(line, "f\"") || strings.Contains(line, "f'") || strings.Contains(line, ".format(") || strings.Contains(line, "% ")
	case "javascript", "typescript":
		return strings.Contains(line, "`${") || strings.Contains(line, "${") || strings.Contains(line, " + ")
	case "go":
		return strings.Contains(line, "fmt.Sprintf(") || strings.Contains(line, " + ")
	case "java":
		return strings.Contains(line, " + ") && (strings.Contains(line, "executeQuery(") || strings.Contains(line, "executeUpdate(") || strings.Contains(line, "createStatement("))
	case "rust", "c", "cpp":
		return strings.Contains(line, "format!(") || strings.Contains(line, "sprintf(") || strings.Contains(line, "snprintf(")
	default:
		return false
	}
}

func rustPanicPattern(language, line string) bool {
	return language == "rust" && (strings.Contains(line, ".unwrap()") || strings.Contains(line, ".expect("))
}

func unsafeBufferPattern(language, line string) bool {
	if language != "c" && language != "cpp" {
		return false
	}
	return strings.Contains(line, "strcpy(") || strings.Contains(line, "strcat(") || strings.Contains(line, "sprintf(") || strings.Contains(line, "gets(")
}

func asyncBlockingPattern(language, line string, hasAsyncDef bool) bool {
	switch language {
	case "python":
		return hasAsyncDef && (strings.Contains(line, "time.sleep(") || strings.Contains(line, "requests.get(") || strings.Contains(line, "requests.post("))
	case "javascript", "typescript":
		return strings.Contains(line, "execSync(") || strings.Contains(line, "readFileSync(")
	default:
		return false
	}
}

var (
	awsAccessKeyPattern = regexp.MustCompile(`AKIA[0-9A-Z]{16}`)
	quotedSecretPattern = regexp.MustCompile(`(?i)(password|secret|api[_-]?key|token)\s*=\s*["'][^"']{4,}["']`)
)

func unsafeRegexCompilePattern(language, line string) bool {
	if language != "go" {
		return false
	}
	call := ""
	switch {
	case strings.Contains(line, "regexp.MustCompile("):
		call = "regexp.MustCompile("
	case strings.Contains(line, "regexp.Compile("):
		call = "regexp.Compile("
	default:
		return false
	}
	start := strings.Index(line, call) + len(call)
	rest := strings.TrimSpace(line[start:])
	if rest == "" {
		return false
	}
	if strings.HasPrefix(rest, `"`) || strings.HasPrefix(rest, "`") {
		return false
	}
	return true
}

func hardcodedSecretPattern(language, line string) bool {
	if qualityLineIgnored(line) {
		return false
	}
	if strings.Contains(line, "BEGIN RSA PRIVATE KEY") || strings.Contains(line, "BEGIN PRIVATE KEY") {
		return true
	}
	if awsAccessKeyPattern.MatchString(line) {
		return true
	}
	if quotedSecretPattern.MatchString(line) {
		return true
	}
	if language == "go" && strings.Contains(strings.ToLower(line), "password") && strings.Contains(line, `"`) && !strings.Contains(line, `""`) {
		return strings.Contains(line, "=")
	}
	return false
}

func ignoredNonCleanupResult(line string) bool {
	if !strings.Contains(line, "_ =") && !strings.Contains(line, ", _ :=") && !strings.Contains(line, ", _ =") {
		return false
	}
	for _, allowed := range []string{
		"_ = db.Close()",
		"_ = tx.Rollback()",
		"_ = file.Close()",
		"_ = os.Remove(",
		"_ = ln.Close()",
		"_ = r.Body.Close()",
		"_ = runtime.listener.Close()",
		"_ = runtime.cmd.Process.Signal(",
		"_ = runtime.cmd.Process.Kill()",
		"_ = cmd.Process.Signal(",
		"_ = cmd.Process.Kill()",
		"_ = cmd.Wait()",
		"_ = process.Signal(",
		"_ = w.Write(",
		"_ = json.NewEncoder(",
	} {
		if strings.Contains(line, allowed) {
			return false
		}
	}
	return true
}

func ignorableQualityFile(path string) bool {
	return strings.HasPrefix(path, "internal/providers/") ||
		strings.HasPrefix(path, "internal/codereview/")
}

func qualityRank(kind string) int {
	switch kind {
	case "shell_injection", "sql_interpolation", "unsafe_buffer", "unsafe_html", "unsafe_regex", "hardcoded_secret":
		return 0
	case "panic":
		return 1
	case "ignored_result":
		return 2
	case "async_blocking", "sleep_polling":
		return 3
	case "direct_logging":
		return 4
	default:
		return 5
	}
}
