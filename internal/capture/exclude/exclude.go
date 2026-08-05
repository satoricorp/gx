package exclude

import (
	"path/filepath"
	"strings"
)

// DefaultPatterns are sensible generated/vendor exclusions for the spike.
var DefaultPatterns = []string{
	"**/go.sum",
	"**/*_gen.go",
	"**/.lgtm/**",
	"**/_generated/**",
	"**/package-lock.json",
	"**/yarn.lock",
	"**/pnpm-lock.yaml",
	"**/vendor/**",
	"**/node_modules/**",
	"**/*.pb.go",
	"**/dist/**",
	"**/build/**",
}

// Matcher decides whether a repo-relative file path should be excluded.
type Matcher struct {
	patterns []string
}

// NewMatcher builds an exclusion matcher from glob-like patterns.
func NewMatcher(patterns []string) *Matcher {
	if len(patterns) == 0 {
		patterns = DefaultPatterns
	}
	return &Matcher{patterns: patterns}
}

// IsExcluded reports whether filePath matches any configured pattern.
func (m *Matcher) IsExcluded(filePath string) bool {
	filePath = filepath.ToSlash(strings.TrimPrefix(filePath, "./"))
	for _, pattern := range m.patterns {
		if matchGlob(pattern, filePath) {
			return true
		}
	}
	return false
}

func matchGlob(pattern, path string) bool {
	pattern = filepath.ToSlash(pattern)
	if matched, err := filepath.Match(pattern, path); err == nil && matched {
		return true
	}
	if strings.HasSuffix(pattern, "/**") {
		prefix := strings.TrimSuffix(pattern, "/**")
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	if strings.HasPrefix(pattern, "**/") {
		suffix := strings.TrimPrefix(pattern, "**/")
		if strings.HasSuffix(suffix, "/**") {
			mid := strings.TrimSuffix(suffix, "/**")
			if strings.Contains(path, "/"+mid+"/") || strings.HasPrefix(path, mid+"/") {
				return true
			}
			return false
		}
		if strings.Contains(suffix, "*") {
			if matched, err := filepath.Match(suffix, filepath.Base(path)); err == nil && matched {
				return true
			}
			return false
		}
		return strings.HasSuffix(path, "/"+suffix) || path == suffix || strings.Contains(path, "/"+suffix)
	}
	return false
}
