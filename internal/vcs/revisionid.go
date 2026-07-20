package vcs

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

const revisionIDByteLength = 16

var revisionIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{16,24}$`)

// GenerateRevisionID returns a new URL-safe 128-bit revision identifier.
func GenerateRevisionID() (string, error) {
	var raw [revisionIDByteLength]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate revision id: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

// ValidRevisionID reports whether id looks like a GX revision identifier.
func ValidRevisionID(id string) bool {
	id = strings.TrimSpace(id)
	if id == "" {
		return false
	}
	return revisionIDPattern.MatchString(id)
}
