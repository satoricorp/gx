package cli

import (
	"os"
	"strings"
)

// clientSurfaces are the surfaces a review run can be invoked from. The set is
// closed on purpose: the value is stored on the gx Cloud history row, and an
// arbitrary environment string would let one misconfigured machine mint
// unbounded new categories.
var clientSurfaces = map[string]bool{
	"cli":         true,
	"mcp":         true,
	"skill":       true,
	"slash-gx":    true,
	"slash-gates": true,
}

// resolveClientSurface resolves which surface invoked this run: an explicit
// override (a --client flag) wins, then $GX_CLIENT, then $GX_MCP, then "cli".
// Unknown values fall through rather than being recorded, so the fallback
// chain — not the caller — owns the vocabulary.
func resolveClientSurface(override string) string {
	for _, candidate := range []string{override, os.Getenv("GX_CLIENT")} {
		if value := strings.TrimSpace(candidate); clientSurfaces[value] {
			return value
		}
	}
	if strings.TrimSpace(os.Getenv("GX_MCP")) != "" {
		return "mcp"
	}
	return "cli"
}
