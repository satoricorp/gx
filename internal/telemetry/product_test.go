package telemetry_test

import (
	"testing"

	"github.com/satoricorp/gx/internal/telemetry"
)

func TestClientSurface(t *testing.T) {
	cases := []struct {
		name     string
		override string
		envVar   string
		mcp      string
		want     string
	}{
		{name: "default", want: "cli"},
		{name: "env", envVar: "skill", want: "skill"},
		{name: "override wins over env", override: "mcp", envVar: "skill", want: "mcp"},
		{name: "slash surface", envVar: "slash-gx", want: "slash-gx"},
		{name: "slash gates surface", envVar: "slash-gates", want: "slash-gates"},
		{name: "whitespace trimmed", envVar: "  mcp  ", want: "mcp"},
		// Unknown values fall through the chain instead of minting a new
		// analytics category.
		{name: "unknown env falls back", envVar: "vscode", want: "cli"},
		{name: "unknown override falls back to env", override: "vscode", envVar: "skill", want: "skill"},
		{name: "legacy GX_MCP still counts as mcp", mcp: "1", want: "mcp"},
		{name: "env wins over GX_MCP", envVar: "slash-gx", mcp: "1", want: "slash-gx"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("GX_CLIENT", tc.envVar)
			t.Setenv("GX_MCP", tc.mcp)
			if got := telemetry.ClientSurface(tc.override); got != tc.want {
				t.Fatalf("ClientSurface(%q) with GX_CLIENT=%q GX_MCP=%q = %q, want %q", tc.override, tc.envVar, tc.mcp, got, tc.want)
			}
		})
	}
}
