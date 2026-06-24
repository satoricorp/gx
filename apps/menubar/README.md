# GX Menu-Bar App

Native macOS menu-bar app for GX. It bundles the `gx` CLI, installs it to
`~/.local/bin/gx` on launch, shows `gx doctor --json` stats, links to the GX
site, and provides MCP install snippets.

## Quick start

```bash
just menubar
```

Artifacts:

- `apps/menubar/dist/GX.app`
- `apps/menubar/dist/GX-macOS.zip`
- `apps/menubar/dist/GX-macOS.dmg`

Open the DMG to get the standard drag-to-Applications install window.

## Tray menu

| Item | Behavior |
|------|----------|
| **Doctor** | Parsed from `gx doctor --json`; hook, auth, Cursor, disk, and backlog checks |
| **Stats** | Capture backlog, disk, and ledger summary from doctor JSON |
| **Run Doctor** | Refreshes `gx doctor --json` immediately |
| **Open gx.run** | Opens `GX_CONSOLE_URL`, or `https://gx.run` by default |
| **MCP** | Copies Cursor, Codex, and Claude Desktop install snippets for the bundled MCP server |
| **Quit** | Exits the app |

## Finding the `gx` binary

Resolution order:

1. `GX_BINARY` env override (dev)
2. Installed CLI: `~/.local/bin/gx`
3. Packaged app: `Contents/Resources/bin/gx`
4. Fallback: `gx` on `PATH`

The app installs or updates `~/.local/bin/gx` from the bundled CLI on launch.

## MCP

The package includes a standalone `gx-mcp` stdio binary. After
dragging `GX.app` to Applications, use the menu's MCP items to copy a Cursor
install command, Codex install command, or Claude Desktop JSON. The snippets point MCP at:

- `GX_BINARY=~/.local/bin/gx`
- `/Applications/GX.app/Contents/Resources/bin/gx-mcp`

MCP is stdio-only. Cursor, Codex, or Claude launches the server when it needs a
tool call; the menu-bar app does not run a local HTTP MCP daemon. Cloud review
context uses `gx auth login` credentials:

```bash
gx auth login
```

The MCP snippets do not embed a token. The local `gx` CLI resolves saved
GitHub auth from disk when MCP tools call cloud endpoints.

For local menu-bar testing from the repo:

```bash
cd mcp && bun install && bun run build && cd ..
go build -o apps/menubar/bin/gx ./cmd/gx
export GX_BINARY="$PWD/apps/menubar/bin/gx"
swift run --package-path apps/menubar
```

Then open **MCP** from the tray menu. Copied commands will point at the
local `$PWD/mcp/dist/gx-mcp` build. For the installed app, run
`apps/menubar/scripts/package.sh`, move `apps/menubar/dist/GX.app` to
Applications, and use the same MCP menu.

## Doctor integration

```bash
gx doctor --json
```

The menubar reads `doctor.capture` for registered repo hooks, upload auth, Cursor vscdb, disk space, and staging backlog. Tray states:

- **Green** — capture checks pass
- **Yellow** — backlog > 10 or non-fatal warnings
- **Red** — missing registered repo hook, auth, Cursor, or low disk

## Package

```bash
just menubar
```

The app is ad-hoc signed with `codesign -`. There is no Developer ID,
notarization, Sparkle, or App Store distribution in this package.

## What's Missing

- Login flow in tray (use `gx auth login` in terminal for now)
- Launch at login helper
- Developer ID signing and notarization
- Windows/Linux tray
- GitHub App install onboarding (WP-6)

## Reference

Menu bar template icons live in `apps/menubar/assets/trayTemplate*.png`. The Finder/Applications app icon is `apps/menubar/assets/appIcon.icns`, generated from the chrome GX artwork on a dark app tile.
