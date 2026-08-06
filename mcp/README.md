# gx MCP Server

TypeScript MCP server (xmcp) that runs over stdio and shells to the local `gx` CLI. Local repository and Git work stay on the user's machine; cloud review context is reached by outbound HTTPS from `gx` when cloud auth or API-key env is configured.

## Primary workflow

1. `git add`, then `git commit` to record work. gx's `prepare-commit-msg` hook stamps each commit with its gx revision trailer and `post-commit` records it — no gx-specific commit verb is required.
2. Publish with plain `git push` when the stack is ready — the gx pre-push hook captures the agent session, links edits to the changed hunks, and publishes the metadata that becomes the PR summary. Open the PR with `gh pr create`.
3. `gx_review` when codegen needs review context from local facts, previous sessions, PRs, current code changes, and optional prompt guidance.

When a user says "save work", "save using gx", or "save with gx", treat that as
a request to stage with `git add`, commit with `git commit`, and publish ready
stacks with plain `git push` unless the user explicitly asks to keep the work
local.

gx PR summaries are posted for PRs whose branch was pushed through gx with `git push`
while the pre-push hook is installed. A PR opened before that push will not get a summary
until the branch is pushed through gx.

## Agent instructions

Add a block like this to the start or end of `AGENTS.md` or `CLAUDE.md` in repos
where agents should use gx:

```md
Version control: use plain Git. gx works through Git hooks, so no gx-specific save verb is needed.

- `git add`, then `git commit -m "..."` to record work. The gx `prepare-commit-msg`
  hook stamps the commit with its gx revision trailer; `post-commit` records it.
- To amend, use `git commit --amend` and preserve the gx revision trailer.
- `git status` to inspect the working tree.

Publish with plain `git push` (the gx pre-push hook captures the agent session, links
edits to the changed hunks, and publishes the metadata that becomes the PR summary),
then open the PR with `gh pr create`. Do not run `gx push` or `gx capture push` — they
bypass or suppress the hook.

Use `gx_review` (MCP) or `gx review` (CLI) for review context on the current change.

gx PR summaries are posted for PRs whose branch was pushed through gx with `git push`
while the pre-push hook is installed. A PR opened before that push will not get a summary
until the branch is pushed through gx.
```

## Tools

| Tool | CLI | Purpose |
|------|-----|---------|
| `gx_review` | `gx review --no-publish [prompt]` | Gather local review/context with AI reviewers enabled |

`gx_review` always passes `--no-publish`. `gx review` on its own posts a review
comment on the matching GitHub pull request and records the run to gx Cloud,
which an agent calling the tool for context mid-codegen should never do — so
publishing stays with the CLI, where a human typed the command.

Saving and publishing are not MCP tools. Record work with plain `git add` + `git commit`
(gx's hooks stamp and record the revision), push the stack with plain `git push` (the
pre-push hook captures the session and publishes), then open the PR with `gh pr create`.
Inspect state with `git status`.

## Install

The server is published to npm as [`@satoricorp/gx`](https://www.npmjs.com/package/@satoricorp/gx).
The simplest way to register it with any coding agent (Claude Code, Cursor,
Codex, VS Code, Zed, …) is [add-mcp](https://add-mcp.com):

```bash
npx add-mcp @satoricorp/gx --name gx
```

By default that writes project-level config (`.mcp.json`, `.cursor/mcp.json`, …)
in the current directory; add `-g` to register it user-wide for every project
instead.

Any MCP client works with the equivalent stdio config: command `npx`,
args `["-y", "@satoricorp/gx"]`.

The CLI install also bundles the server as a standalone `gx-mcp` binary in
`~/.local/bin`, so no Node is required:

```bash
curl -fsSL https://download.gx.run/install.sh | sh
```

Building from source:

```bash
cd mcp
bun install
bun run build
```

For cloud auth, log in once with GitHub:

```bash
gx auth login
```

The installer provides the `gx` CLI and `gx-mcp` binary. Repo Git
hooks are installed when a repo is initialized with `gx init`; the `gx_review`
MCP tool is read-only — it never initializes a repo, never writes `~/.gx`,
never touches `.git/index`, and never publishes. The `pre-push` hook runs `gx capture push` for
the pushed ref range, stages captured Claude/Codex/Cursor session context in
`~/.gx/gx.db`, and uploads only when gx upload credentials are configured.

Development stdio:

```bash
GX_BINARY="$(command -v gx)" bun run start
```

## Scripts

| Script | Description |
|--------|-------------|
| `bun run dev` | xmcp dev with watch |
| `bun run build` | Production xmcp output plus standalone `dist/gx-mcp` |
| `bun run build:xmcp` | Production xmcp JavaScript output to `dist/` |
| `bun run build:binary` | Standalone stdio binary at `dist/gx-mcp` |
| `bun run start` | stdio MCP (`dist/gx-mcp`) |
| `bun test` | Tool schema and CLI spawn tests |
| `bun run typecheck` | TypeScript check |

## Publishing

`npm publish` from `mcp/` (or `just mcp-publish` from the repo root) builds,
tests, and publishes `@satoricorp/gx`. CI does the same automatically when
a `mcp-v<version>` tag is pushed — the tag must match `package.json`'s version:

```bash
git tag mcp-v0.1.0 && git push origin mcp-v0.1.0
```

CI publishes via npm Trusted Publishing (OIDC) — no token secret. One-time
setup: publish once locally, then on npmjs.com under the package's Settings add
a trusted publisher (GitHub Actions, owner `satoricorp`, repo `gx`,
workflow `mcp.yml`). The npm package ships only the bundled `dist/*.js`; the
compiled `dist/gx-mcp` binary is distributed by the CLI installer instead.

## Environment

| Variable | Description |
|----------|-------------|
| `GX_BINARY` | Optional path to `gx` executable |
| `GX_CLOUD_URL` | gx cloud API base URL for review AI fallback |

Without a `GX_BINARY` override, MCP uses `~/.local/bin/gx` when present, then falls back to `gx` on `PATH`. Cloud calls use credentials from `gx auth login` when available.

The released `gx-mcp` checks `https://download.gx.run/cli/manifest.json` and adds an
update notice plus `curl -fsSL https://download.gx.run/install.sh | sh` to tool responses
when a newer CLI package is available. Set `GX_MCP_UPDATE_CHECK=0` to disable
that check.
