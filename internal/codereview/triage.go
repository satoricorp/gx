package codereview

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ChangeTriage struct {
	Class     string   `json:"class,omitempty"`
	RiskTags  []string `json:"risk_tags,omitempty"`
	Languages []string `json:"languages,omitempty"`
	Rationale []string `json:"rationale,omitempty"`
}

type ReviewExecutionPlan struct {
	Triage             ChangeTriage
	ActiveScopes       []string
	RunStaticTools     bool
	RunAI              bool
	RunReviewResources bool
	RiskTags           []string
}

type RetrieveInput struct {
	RepoRoot     string
	Options      Options
	Facts        RepoFacts
	Hints        []ReviewHint
	Plan         ReviewExecutionPlan
	ChangedFiles []string
	DiffSnippets []DiffSnippet
}

func TriageChange(changedFiles []string, diffSnippets []DiffSnippet, opts Options) ChangeTriage {
	files := normalizedChangedFiles(changedFiles)
	triage := ChangeTriage{
		Class:     "code",
		Languages: languageTagsForFiles(files),
	}
	riskTags := stringSet(riskTagsForReview(files, nil, opts))
	var rationale []string

	if len(files) > 0 && allFiles(files, isDocsOnlyReviewFile) {
		triage.Class = "docs-only"
		rationale = append(rationale, "changed files are documentation, text, license, image, or CODEOWNERS files")
	} else if len(files) > 0 && allFiles(files, isConfigOnlyReviewFile) {
		triage.Class = "config-only"
		rationale = append(rationale, "changed files are review policy, agent/security docs, CI, lockfile, or dot/editor configuration files")
	} else if len(files) > 0 && securitySensitiveChange(files, diffSnippets, riskTags) {
		triage.Class = "security-sensitive"
		rationale = append(rationale, "paths or diff content touch security-sensitive review areas")
	} else if len(files) > 0 && mechanicalChange(files, diffSnippets) {
		triage.Class = "mechanical"
		rationale = append(rationale, "diff appears generated, rename-only, or whitespace-dominant")
	} else {
		rationale = append(rationale, "default code review classification")
	}

	triage.RiskTags = sortedKeys(riskTags)
	triage.Rationale = rationale
	return triage
}

func triageEnabled() bool {
	return !strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_TRIAGE")), "0")
}

func reviewPlanFor(opts Options, triage ChangeTriage) ReviewExecutionPlan {
	riskTags := triage.RiskTags
	if len(riskTags) == 0 {
		riskTags = riskTagsForReview(nil, nil, opts)
	}
	plan := ReviewExecutionPlan{
		Triage:             triage,
		ActiveScopes:       activeScopeListForTriage(opts, triage),
		RunStaticTools:     true,
		RunAI:              true,
		RunReviewResources: true,
		RiskTags:           append([]string(nil), riskTags...),
	}
	if triage.Class == "docs-only" && opts.PatchFocused && !opts.Deep && strings.TrimSpace(opts.Prompt) == "" {
		plan.RunStaticTools = false
		plan.RunAI = false
		plan.RunReviewResources = false
	}
	return plan
}

func reviewExecutionPlanConfigured(plan ReviewExecutionPlan) bool {
	return plan.Triage.Class != "" ||
		len(plan.ActiveScopes) > 0 ||
		len(plan.RiskTags) > 0 ||
		plan.RunStaticTools ||
		plan.RunAI ||
		plan.RunReviewResources
}

func activeScopeListForTriage(opts Options, triage ChangeTriage) []string {
	if !opts.PatchFocused {
		return activeScopeList(opts)
	}
	switch triage.Class {
	case "docs-only":
		return []string{"docs"}
	case "config-only":
		return []string{"dependencies", "testing", "maintainability"}
	case "mechanical":
		return []string{"maintainability", "testing"}
	case "security-sensitive":
		return []string{"security", "dependencies", "testing", "maintainability"}
	default:
		return activeScopeList(opts)
	}
}

func shouldShortCircuitDocsOnly(opts Options, explicitScope bool, plan ReviewExecutionPlan) bool {
	return plan.Triage.Class == "docs-only" &&
		opts.PatchFocused &&
		!opts.Deep &&
		!explicitScope &&
		strings.TrimSpace(opts.Prompt) == ""
}

func allFiles(files []string, keep func(string) bool) bool {
	for _, file := range files {
		if !keep(file) {
			return false
		}
	}
	return len(files) > 0
}

func isDocsOnlyReviewFile(rel string) bool {
	lower := strings.ToLower(filepath.ToSlash(strings.TrimSpace(rel)))
	base := filepath.Base(lower)
	switch base {
	case "review.md", "agents.md", "security.md":
		return false
	}
	if base == "codeowners" || strings.HasPrefix(base, "license") || strings.HasPrefix(base, "copying") ||
		strings.HasPrefix(base, "notice") || strings.HasPrefix(base, "readme") || strings.HasPrefix(base, "changelog") ||
		strings.HasPrefix(base, "contributing") || strings.HasPrefix(base, "authors") {
		return true
	}
	if strings.HasPrefix(lower, "docs/") || strings.Contains(lower, "/docs/") {
		return true
	}
	switch filepath.Ext(lower) {
	case ".md", ".markdown", ".mdx", ".rst", ".txt", ".adoc", ".asciidoc", ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".svg", ".pdf":
		return true
	default:
		return false
	}
}

func isConfigOnlyReviewFile(rel string) bool {
	lower := strings.ToLower(filepath.ToSlash(strings.TrimSpace(rel)))
	base := filepath.Base(lower)
	switch base {
	case "review.md", "agents.md", "security.md", "dependabot.yml", "dependabot.yaml",
		"codeql.yml", "codeql.yaml", ".editorconfig", ".gitattributes", ".gitignore":
		return true
	}
	if isDependencyFile(rel) || strings.HasSuffix(lower, ".lock") {
		return true
	}
	if strings.HasPrefix(lower, ".github/workflows/") || strings.HasPrefix(lower, ".vscode/") || strings.HasPrefix(lower, ".idea/") {
		return true
	}
	if strings.HasPrefix(base, ".") {
		return true
	}
	if strings.Contains(base, "config") && (strings.HasSuffix(base, ".json") || strings.HasSuffix(base, ".yaml") || strings.HasSuffix(base, ".yml") || strings.HasSuffix(base, ".toml")) {
		return true
	}
	return false
}

func securitySensitiveChange(files []string, snippets []DiffSnippet, tags map[string]struct{}) bool {
	for _, file := range files {
		lower := strings.ToLower(filepath.ToSlash(file))
		switch {
		case isDependencyFile(file):
			tags["dependencies"] = struct{}{}
			tags["supply-chain"] = struct{}{}
			return true
		case strings.Contains(lower, "auth"), strings.Contains(lower, "session"), strings.Contains(lower, "permission"), strings.Contains(lower, "crypto"), strings.Contains(lower, "secret"), strings.Contains(lower, "token"):
			tags["secure-coding"] = struct{}{}
			return true
		case strings.HasSuffix(lower, ".sql"), strings.Contains(lower, "migration"), strings.Contains(lower, "schema"), strings.Contains(lower, "db/"):
			tags["database-security"] = struct{}{}
			tags["sql-injection"] = struct{}{}
			tags["injection"] = struct{}{}
			return true
		case strings.Contains(lower, "parser"), strings.Contains(lower, "deserialize"), strings.Contains(lower, "unmarshal"), strings.Contains(lower, "path"):
			return true
		case strings.Contains(lower, "scripts/"), strings.HasSuffix(lower, ".sh"):
			tags["command-injection"] = struct{}{}
			tags["injection"] = struct{}{}
			return true
		}
	}
	for _, snippet := range snippets {
		if securitySensitiveText(snippet.File+"\n"+snippet.Diff, tags) {
			return true
		}
	}
	return false
}

func securitySensitiveText(text string, tags map[string]struct{}) bool {
	lower := strings.ToLower(text)
	matched := false
	add := func(tag string) {
		tags[tag] = struct{}{}
		matched = true
	}
	for _, needle := range []string{"regexp.", "regex", "mustcompile", "compile("} {
		if strings.Contains(lower, needle) {
			add("regex")
			add("injection")
			break
		}
	}
	for _, needle := range []string{"select ", "insert ", "update ", "delete ", "sql.", "query(", "execcontext(", "exec("} {
		if strings.Contains(lower, needle) {
			add("sql-injection")
			add("injection")
			break
		}
	}
	for _, needle := range []string{"exec.command", "os/exec", "sh -c", "bash -c", "shell"} {
		if strings.Contains(lower, needle) {
			add("command-injection")
			add("injection")
			break
		}
	}
	for _, needle := range []string{"json.unmarshal", "yaml.unmarshal", "deserialize", "deserializ", "pickle", "gob.decode"} {
		if strings.Contains(lower, needle) {
			add("deserialization")
			break
		}
	}
	for _, needle := range []string{"../", "..\\\\", "filepath.clean", "filepath.join", "path traversal", "zip slip"} {
		if strings.Contains(lower, needle) {
			add("path-traversal")
			break
		}
	}
	for _, needle := range []string{"auth", "session", "crypto", "permission", "csrf", "jwt", "token", "password", "secret", "input parsing"} {
		if strings.Contains(lower, needle) {
			add("secure-coding")
			break
		}
	}
	return matched
}

func mechanicalChange(files []string, snippets []DiffSnippet) bool {
	return allFiles(files, isGeneratedReviewFile) || renameOnlyDiff(snippets) || whitespaceDominantDiff(snippets)
}

func isGeneratedReviewFile(rel string) bool {
	lower := strings.ToLower(filepath.ToSlash(rel))
	base := filepath.Base(lower)
	return strings.Contains(lower, "/generated/") ||
		strings.Contains(base, "generated") ||
		strings.HasSuffix(base, ".pb.go") ||
		strings.HasSuffix(base, "_generated.go") ||
		strings.HasSuffix(base, "_gen.go")
}

func renameOnlyDiff(snippets []DiffSnippet) bool {
	if len(snippets) == 0 {
		return false
	}
	for _, snippet := range snippets {
		for _, line := range strings.Split(snippet.Diff, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "diff --git ") || strings.HasPrefix(line, "similarity index ") ||
				strings.HasPrefix(line, "rename from ") || strings.HasPrefix(line, "rename to ") {
				continue
			}
			return false
		}
	}
	return true
}

func whitespaceDominantDiff(snippets []DiffSnippet) bool {
	total := 0
	whitespaceOnly := 0
	for _, snippet := range snippets {
		var removed []string
		for _, line := range strings.Split(snippet.Diff, "\n") {
			if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
				continue
			}
			if strings.HasPrefix(line, "-") {
				removed = append(removed, strings.TrimSpace(strings.TrimPrefix(line, "-")))
				total++
			}
			if strings.HasPrefix(line, "+") {
				added := strings.TrimSpace(strings.TrimPrefix(line, "+"))
				total++
				for len(removed) > 0 {
					candidate := removed[0]
					removed = removed[1:]
					if candidate == added {
						whitespaceOnly += 2
						break
					}
				}
			}
		}
	}
	return total >= 4 && whitespaceOnly*100/total >= 80
}

func stringSet(values []string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			out[value] = struct{}{}
		}
	}
	return out
}

func sortedStrings(values []string) []string {
	out := append([]string(nil), values...)
	sort.Strings(out)
	return out
}
