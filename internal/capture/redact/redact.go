package redact

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"
)

// Kind constants for redaction placeholders.
const (
	KindAWS         = "aws-key"
	KindGitHub      = "github-token"
	KindSlack       = "slack-token"
	KindAPIKey      = "api-key"
	KindBearer      = "bearer"
	KindPassword    = "password"
	KindPrivateKey  = "private-key"
	KindHighEntropy = "high-entropy"
)

type rule struct {
	kind    string
	pattern *regexp.Regexp
}

var builtinRules = []rule{
	{KindAWS, regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
	{KindGitHub, regexp.MustCompile(`gh[pousr]_[A-Za-z0-9_]{20,}`)},
	{KindSlack, regexp.MustCompile(`xox[baprs]-[0-9A-Za-z-]{10,}`)},
	{KindAPIKey, regexp.MustCompile(`(?i)(?:api[_-]?key|secret[_-]?key)\s*[:=]\s*['"]?([A-Za-z0-9_\-./+=]{12,})`)},
	{KindBearer, regexp.MustCompile(`(?i)Bearer\s+([A-Za-z0-9_\-.=+/]{12,})`)},
	{KindPassword, regexp.MustCompile(`(?i)(?:password|passwd|pwd)\s*[:=]\s*['"]?([^\s'"]{4,})`)},
	{KindPrivateKey, regexp.MustCompile(`-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----`)},
}

var keywordContexts = []*regexp.Regexp{
	regexp.MustCompile(`(?i)Authorization:\s*[^\s]+`),
	regexp.MustCompile(`(?i)Bearer\s+[A-Za-z0-9_\-.=+/]{12,}`),
}

var entropyToken = regexp.MustCompile(`[A-Za-z0-9+/=_\-]{20,}`)

// Redact applies Tier-1 secret redaction to text.
func Redact(text string) string {
	if text == "" {
		return text
	}
	out := text
	for _, r := range builtinRules {
		out = r.pattern.ReplaceAllStringFunc(out, func(match string) string {
			secret := match
			if sub := r.pattern.FindStringSubmatch(match); len(sub) > 1 && sub[1] != "" {
				secret = sub[1]
			}
			return placeholder(r.kind, secret)
		})
	}
	for _, pattern := range keywordContexts {
		out = pattern.ReplaceAllStringFunc(out, func(match string) string {
			return placeholder(KindBearer, match)
		})
	}
	out = entropyToken.ReplaceAllStringFunc(out, func(token string) string {
		if strings.HasPrefix(token, "[REDACTED:") {
			return token
		}
		if shannonEntropy(token) >= 4.2 && len(token) >= 24 {
			return placeholder(KindHighEntropy, token)
		}
		return token
	})
	return out
}

func placeholder(kind, secret string) string {
	hash := shortHash(secret)
	return fmt.Sprintf("[REDACTED:%s:%s]", kind, hash)
}

func shortHash(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:2])
}

func shannonEntropy(s string) float64 {
	if s == "" {
		return 0
	}
	freq := map[rune]int{}
	for _, r := range s {
		if unicode.IsPrint(r) || r == '\t' {
			freq[r]++
		}
	}
	n := float64(len(s))
	var entropy float64
	for _, count := range freq {
		p := float64(count) / n
		entropy -= p * math.Log2(p)
	}
	return entropy
}
