package codereview

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const maxCodeQualityHints = 40

func collectCodeQualityHints(repoRoot string, facts RepoFacts, opts Options) []CodeQualityHint {
	var hints []CodeQualityHint
	for _, file := range facts.Files {
		if !strings.HasSuffix(file, ".go") || isTestFile(file) || ignorableQualityFile(file) {
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
	var hints []CodeQualityHint
	for i, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}
		lineNo := i + 1
		switch {
		case ignoredNonCleanupResult(trimmed):
			hints = append(hints, CodeQualityHint{
				Kind:   "ignored_result",
				File:   rel,
				Line:   lineNo,
				Text:   trimmed,
				Reason: "A discarded result or error in production code can hide malformed input, failed cleanup, or failed recovery paths.",
			})
		case strings.Contains(trimmed, "time.Sleep("):
			hints = append(hints, CodeQualityHint{
				Kind:   "sleep_polling",
				File:   rel,
				Line:   lineNo,
				Text:   trimmed,
				Reason: "Sleeping in production code is often a brittle readiness or polling seam unless bounded and tested.",
			})
		case strings.Contains(trimmed, "log.Println(") || strings.Contains(trimmed, "log.Printf("):
			hints = append(hints, CodeQualityHint{
				Kind:   "direct_logging",
				File:   rel,
				Line:   lineNo,
				Text:   trimmed,
				Reason: "Direct logging from a Module can make error handling and tests less observable than returning classified errors.",
			})
		case strings.Contains(trimmed, "panic("):
			hints = append(hints, CodeQualityHint{
				Kind:   "panic",
				File:   rel,
				Line:   lineNo,
				Text:   trimmed,
				Reason: "Panic in production code should be justified by an invariant that callers cannot recover from.",
			})
		}
	}
	return hints
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
	case "panic":
		return 0
	case "ignored_result":
		return 1
	case "sleep_polling":
		return 2
	case "direct_logging":
		return 3
	default:
		return 4
	}
}
