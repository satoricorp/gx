# GX Menu-Bar App

Native macOS menu-bar app for GX. It bundles the `gx` CLI, installs it plus
the `gxg`, `gxr`, and `gxs` shortcuts to `~/.local/bin` on launch, shows
`gx doctor --json` status and stats, links to gx.run, and provides MCP setup
snippets.

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
| **Status** | Parsed from `gx doctor --json`; green/yellow/red status dot plus hook, auth, Cursor, disk, and backlog checks |
| **Stats** | Capture backlog, disk, and ledger summary from doctor JSON |
| **Update CLI** | Installs or updates `~/.local/bin/gx`, `gxg`, `gxr`, and `gxs` from the bundled CLI |
| **Open https://gx.run** | Opens `https://gx.run` |
| **MCP** | Shows setup instructions at `https://docs.gx.run` |
| **Quit** | Exits the app |

## Finding the `gx` binary

Resolution order:

1. `GX_BINARY` env override (dev)
2. Installed CLI: `~/.local/bin/gx`
3. Packaged app: `Contents/Resources/bin/gx`
4. Fallback: `gx` on `PATH`

The app installs or updates `~/.local/bin/gx`, `~/.local/bin/gxg`,
`~/.local/bin/gxr`, and `~/.local/bin/gxs` from the bundled CLI on launch and
through the **Update CLI** menu item.

## MCP setup

The package includes a standalone `gx-mcp` stdio binary. After
dragging `GX.app` to Applications, use **MCP > Show Instructions** to open
the current setup instructions:

- `https://docs.gx.run`

MCP is stdio-only. Cursor, Codex, or Claude launches the server when it needs a
tool call; the menu-bar app does not run a local HTTP MCP daemon. Cloud review
context uses `gx auth login` credentials:

```bash
gx auth login
```

The local `gx` CLI resolves saved GitHub auth from disk when MCP tools call
cloud endpoints.

For local menu-bar testing from the repo:

```bash
cd mcp && bun install && bun run build && cd ..
go build -o apps/menubar/bin/gx ./cmd/gx
export GX_BINARY="$PWD/apps/menubar/bin/gx"
swift run --package-path apps/menubar
```

Then open **MCP > Show Instructions** from the tray menu. For the installed
app, run `apps/menubar/scripts/package.sh`, move `apps/menubar/dist/GX.app`
to Applications, and use the same **MCP** menu.

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

## What's missing (MVP gaps)

- Login flow in tray (use `gx login` in terminal for now)
- Launch at login helper
- Developer ID signing and notarization
- Windows/Linux tray
- GitHub App install onboarding (WP-6)
- Real activity API until WP-5h lands (mock/fallback included)

## Reference

Menu bar template icons live in `apps/menubar/assets/trayTemplate*.png`. The Finder/Applications app icon is `apps/menubar/assets/appIcon.icns`, generated from the chrome GX artwork on a dark app tile.
