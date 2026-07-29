# Totality MCP Server

TypeScript MCP server (xmcp) that runs over stdio and shells to the local `tl` CLI. Local repository and Git work stay on the user's machine; cloud review context is reached by outbound HTTPS from `tl` when cloud auth or API-key env is configured.

## Primary workflow

1. `git add`, then `git commit` to record work. Totality's `prepare-commit-msg` hook stamps each commit with its Totality revision trailer and `post-commit` records it — no Totality-specific commit verb is required.
2. Publish with plain `git push` when the stack is ready — the Totality pre-push hook captures the agent session, links edits to the changed hunks, and publishes the metadata that becomes the PR summary. Open the PR with `gh pr create`.
3. `tl_review` when codegen needs review context from local facts, previous sessions, PRs, current code changes, and optional prompt guidance.

When a user says "save work", "save using tl", or "save with tl", treat that as
a request to stage with `git add`, commit with `git commit`, and publish ready
stacks with plain `git push` unless the user explicitly asks to keep the work
local.

Totality PR summaries are posted for PRs whose branch was pushed through Totality with `git push`
while the pre-push hook is installed. A PR opened before that push will not get a summary
until the branch is pushed through Totality.

## Agent instructions

Add a block like this to the start or end of `AGENTS.md` or `CLAUDE.md` in repos
where agents should use Totality:

```md
Version control: use plain Git. Totality works through Git hooks, so no Totality-specific save verb is needed.

- `git add`, then `git commit -m "..."` to record work. The Totality `prepare-commit-msg`
  hook stamps the commit with its Totality revision trailer; `post-commit` records it.
- To amend, use `git commit --amend` and preserve the Totality revision trailer.
- `git status` to inspect the working tree.

Publish with plain `git push` (the Totality pre-push hook captures the agent session, links
edits to the changed hunks, and publishes the metadata that becomes the PR summary),
then open the PR with `gh pr create`. Do not run `tl push` or `tl capture push` — they
bypass or suppress the hook.

Use `tl_review` (MCP) or `tl review` (CLI) for review context on the current change.

Totality PR summaries are posted for PRs whose branch was pushed through Totality with `git push`
while the pre-push hook is installed. A PR opened before that push will not get a summary
until the branch is pushed through Totality.
```

## Tools

| Tool | CLI | Purpose |
|------|-----|---------|
| `tl_review` | `tl review --no-publish [prompt]` | Gather local review/context with AI reviewers enabled |

`tl_review` always passes `--no-publish`. `tl review` on its own posts a review
comment on the matching GitHub pull request and records the run to Totality Cloud,
which an agent calling the tool for context mid-codegen should never do — so
publishing stays with the CLI, where a human typed the command.

Saving and publishing are not MCP tools. Record work with plain `git add` + `git commit`
(Totality's hooks stamp and record the revision), push the stack with plain `git push` (the
pre-push hook captures the session and publishes), then open the PR with `gh pr create`.
Inspect state with `git status`.

## Install

The server is published to npm as [`@satoricorp/totality`](https://www.npmjs.com/package/@satoricorp/totality).
The simplest way to register it with any coding agent (Claude Code, Cursor,
Codex, VS Code, Zed, …) is [add-mcp](https://add-mcp.com):

```bash
npx add-mcp @satoricorp/totality --name totality
```

By default that writes project-level config (`.mcp.json`, `.cursor/mcp.json`, …)
in the current directory; add `-g` to register it user-wide for every project
instead.

Any MCP client works with the equivalent stdio config: command `npx`,
args `["-y", "@satoricorp/totality"]`.

The CLI install also bundles the server as a standalone `tl-mcp` binary in
`~/.local/bin`, so no Node is required:

```bash
curl -fsSL https://download.totality.sh/install.sh | sh
```

Building from source:

```bash
cd mcp
bun install
bun run build
```

For cloud auth, log in once with GitHub:

```bash
tl auth login
```

The installer provides the `tl` CLI and `tl-mcp` binary. Repo Git
hooks are installed when a repo is initialized with `tl init`; the `tl_review`
MCP tool is read-only — it never initializes a repo, never writes `~/.totality`,
never touches `.git/index`, and never publishes. The `pre-push` hook runs `tl capture push` for
the pushed ref range, stages captured Claude/Codex/Cursor session context in
`~/.totality/totality.db`, and uploads only when Totality upload credentials are configured.

Development stdio:

```bash
TOTALITY_BINARY="$(command -v tl)" bun run start
```

## Scripts

| Script | Description |
|--------|-------------|
| `bun run dev` | xmcp dev with watch |
| `bun run build` | Production xmcp output plus standalone `dist/tl-mcp` |
| `bun run build:xmcp` | Production xmcp JavaScript output to `dist/` |
| `bun run build:binary` | Standalone stdio binary at `dist/tl-mcp` |
| `bun run start` | stdio MCP (`dist/tl-mcp`) |
| `bun test` | Tool schema and CLI spawn tests |
| `bun run typecheck` | TypeScript check |

## Publishing

`npm publish` from `mcp/` (or `just mcp-publish` from the repo root) builds,
tests, and publishes `@satoricorp/totality`. CI does the same automatically when
a `mcp-v<version>` tag is pushed — the tag must match `package.json`'s version:

```bash
git tag mcp-v0.1.0 && git push origin mcp-v0.1.0
```

CI publishes via npm Trusted Publishing (OIDC) — no token secret. One-time
setup: publish once locally, then on npmjs.com under the package's Settings add
a trusted publisher (GitHub Actions, owner `satoricorp`, repo `totality`,
workflow `mcp.yml`). The npm package ships only the bundled `dist/*.js`; the
compiled `dist/tl-mcp` binary is distributed by the CLI installer instead.

## Environment

| Variable | Description |
|----------|-------------|
| `TOTALITY_BINARY` | Optional path to `tl` executable |
| `TOTALITY_CLOUD_URL` | Totality cloud API base URL for review AI fallback |

Without a `TOTALITY_BINARY` override, MCP uses `~/.local/bin/tl` when present, then falls back to `tl` on `PATH`. Cloud calls use credentials from `tl auth login` when available.

The released `tl-mcp` checks `https://download.totality.sh/cli/manifest.json` and adds an
update notice plus `curl -fsSL https://download.totality.sh/install.sh | sh` to tool responses
when a newer CLI package is available. Set `TOTALITY_MCP_UPDATE_CHECK=0` to disable
that check.
