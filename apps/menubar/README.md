# GX Menu-Bar App

Native macOS menu-bar app for GX. It bundles the `gx` CLI, installs it plus
the `gxr` shortcut to `~/.local/bin` on launch, shows
`gx doctor --json` status and stats, checks the latest GitHub release, updates
the local CLI when a newer release is available, checks app-bundle updates with
Sparkle, links to gx.run, and provides MCP setup snippets.

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
| **Update CLI to vX.Y.Z** | Downloads the newest matching CLI release asset from GitHub Releases and installs it to `~/.local/bin` |
| **Check for Updates** | Fetches the latest `satoricorp/gx` GitHub release and compares it with `gx version --json` |
| **Install Bundled CLI** | Installs or updates `~/.local/bin/gx` and `gxr` from the bundled CLI |
| **Check for App Updates** | Runs Sparkle against the configured appcast feed and updates `GX.app` in place |
| **Open https://gx.run** | Opens `https://gx.run` |
| **MCP** | Shows setup instructions at `https://docs.gx.run` |
| **Quit** | Exits the app |

## Finding the `gx` binary

Resolution order:

1. `GX_BINARY` env override (dev)
2. Installed CLI: `~/.local/bin/gx`
3. Packaged app: `Contents/Resources/bin/gx`
4. Fallback: `gx` on `PATH`

The app installs or updates `~/.local/bin/gx` and `~/.local/bin/gxr`
from the bundled CLI on launch and through the **Install Bundled CLI**
menu item.

## CLI updates

The tray uses `gx version --json` to read the installed CLI version and fetches
the latest release from:

```text
https://api.github.com/repos/satoricorp/gx/releases/latest
```

When the release tag is newer than the installed CLI `release_version`, the
menu shows **Update CLI to vX.Y.Z**. The updater downloads the matching macOS
archive asset for the current architecture, verifies a sibling `.sha256` asset
when present, extracts it, validates the downloaded `gx version --json`, then
replaces the local binaries in `~/.local/bin`.

For local testing, set `GX_MENUBAR_RELEASE_URL` to a compatible GitHub release
JSON endpoint or fixture server.

## App updates

The menu-bar app uses Sparkle for app-bundle updates. Release builds stamp:

- `GX_APPCAST_URL` into `SUFeedURL` (defaults to `https://download.gx.run/appcast.xml`)
- `GX_SPARKLE_PUBLIC_ED_KEY` into `SUPublicEDKey`
- `GX_APP_BUILD_VERSION` into `CFBundleVersion`

Generate the appcast after packaging:

```bash
GX_APPCAST_DOWNLOAD_URL_PREFIX=https://download.gx.run/menubar \
GX_SPARKLE_PRIVATE_ED_KEY_FILE=/path/to/ed25519.key \
apps/menubar/scripts/generate-appcast.sh
```

For local verification when SwiftPM cannot download Sparkle's binary artifact,
set `GX_LOCAL_SPARKLE_XCFRAMEWORK` to a package-relative `Sparkle.xcframework`
path before running `apps/menubar/scripts/package.sh`.

## Agent setup

After installing the app, add GX instructions to the start or end of your
`AGENTS.md` or `CLAUDE.md` so agents save with GX instead of raw Git:

```md
Version control: use plain Git. GX's hooks record and publish your work automatically —
there is no GX save verb.

- `git add`, then `git commit -m "..."` — a GX hook records the commit as a revision.
- `git status` to inspect the working tree.
- `gx_review` (MCP) or `gx review` (CLI) for AI review of the current change.

To amend an existing GX revision, use `git commit --amend` and preserve the GX revision trailer.

Publish with plain `git push` (the GX pre-push hook captures the session and publishes),
then open the PR with `gh pr create`. Do not run `gx push` or `gx capture push` — they
bypass or suppress the hook.

Save and publish with plain Git:
- `git add`
- `git commit -m "..."`
- `git push`

When the user says "save work", "save using gx", or "save with gx", stage with
`git add`, save with `git commit`, and publish with plain `git push` unless the user
asks to keep the work local. GX's hooks record the revision and publish on push.

GX PR summaries are posted by GX Cloud for PRs whose branch was pushed through GX
with `git push` while the pre-push hook is installed.

Use raw Git for read-only inspection and the `git push` publish step; do not use other raw
Git for the save flow unless the user explicitly asks. If supported, deny or require approval for
`git commit`, `git reset`, and branch deletion.
```

## Git hooks

The menu-bar app gives you the bundled `gx` CLI and `gx-mcp`; repo hooks are
installed when a repo is initialized with `gx init`. GX installs
`prepare-commit-msg` to stamp GX revision
trailers, `post-commit` to record commit metadata, `post-rewrite` to follow
amended or rebased commit OIDs, and `pre-push` to capture and publish each
pushed ref range.

The hook captures Claude, Codex, and Cursor session context into `~/.gx/gx.db`.
It stages locally without network access when upload credentials are absent,
and uploads when GX upload credentials are configured. The tray reads
`gx doctor --json` to report whether registered repo hooks are installed,
missing, or unreachable.

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
notarization, or App Store distribution in this package.

## What's missing (MVP gaps)

- Login flow in tray (use `gx auth login` in terminal for now)
- Launch at login helper
- Developer ID signing and notarization
- Windows/Linux tray
- GitHub App install onboarding (WP-6)
- Real activity API until WP-5h lands (mock/fallback included)

## Reference

Menu bar template icons live in `apps/menubar/assets/trayTemplate*.png`. The Finder/Applications app icon is `apps/menubar/assets/appIcon.icns`, generated from the chrome GX artwork on a dark app tile.
