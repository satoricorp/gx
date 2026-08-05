# lgtm MCP Server

TypeScript MCP server (xmcp) that runs over stdio and shells to the local `lgtm` CLI. Local repository and Git work stay on the user's machine; cloud review context is reached by outbound HTTPS from `lgtm` when cloud auth or API-key env is configured.

## Primary workflow

1. `git add`, then `git commit` to record work. lgtm's `prepare-commit-msg` hook stamps each commit with its lgtm revision trailer and `post-commit` records it — no lgtm-specific commit verb is required.
2. Publish with plain `git push` when the stack is ready — the lgtm pre-push hook captures the agent session, links edits to the changed hunks, and publishes the metadata that becomes the PR summary. Open the PR with `gh pr create`.
3. `lgtm_review` when codegen needs review context from local facts, previous sessions, PRs, current code changes, and optional prompt guidance.

When a user says "save work", "save using lgtm", or "save with lgtm", treat that as
a request to stage with `git add`, commit with `git commit`, and publish ready
stacks with plain `git push` unless the user explicitly asks to keep the work
local.

lgtm PR summaries are posted for PRs whose branch was pushed through lgtm with `git push`
while the pre-push hook is installed. A PR opened before that push will not get a summary
until the branch is pushed through lgtm.

## Agent instructions

Add a block like this to the start or end of `AGENTS.md` or `CLAUDE.md` in repos
where agents should use lgtm:

```md
Version control: use plain Git. lgtm works through Git hooks, so no lgtm-specific save verb is needed.

- `git add`, then `git commit -m "..."` to record work. The lgtm `prepare-commit-msg`
  hook stamps the commit with its lgtm revision trailer; `post-commit` records it.
- To amend, use `git commit --amend` and preserve the lgtm revision trailer.
- `git status` to inspect the working tree.

Publish with plain `git push` (the lgtm pre-push hook captures the agent session, links
edits to the changed hunks, and publishes the metadata that becomes the PR summary),
then open the PR with `gh pr create`. Do not run `lgtm push` or `lgtm capture push` — they
bypass or suppress the hook.

Use `lgtm_review` (MCP) or `lgtm review` (CLI) for review context on the current change.

lgtm PR summaries are posted for PRs whose branch was pushed through lgtm with `git push`
while the pre-push hook is installed. A PR opened before that push will not get a summary
until the branch is pushed through lgtm.
```

## Tools

| Tool | CLI | Purpose |
|------|-----|---------|
| `lgtm_review` | `lgtm review --no-publish [prompt]` | Gather local review/context with AI reviewers enabled |

`lgtm_review` always passes `--no-publish`. `lgtm review` on its own posts a review
comment on the matching GitHub pull request and records the run to lgtm Cloud,
which an agent calling the tool for context mid-codegen should never do — so
publishing stays with the CLI, where a human typed the command.

Saving and publishing are not MCP tools. Record work with plain `git add` + `git commit`
(lgtm's hooks stamp and record the revision), push the stack with plain `git push` (the
pre-push hook captures the session and publishes), then open the PR with `gh pr create`.
Inspect state with `git status`.

## Install

The server is published to npm as [`@satoricorp/lgtm`](https://www.npmjs.com/package/@satoricorp/lgtm).
The simplest way to register it with any coding agent (Claude Code, Cursor,
Codex, VS Code, Zed, …) is [add-mcp](https://add-mcp.com):

```bash
npx add-mcp @satoricorp/lgtm --name lgtm
```

By default that writes project-level config (`.mcp.json`, `.cursor/mcp.json`, …)
in the current directory; add `-g` to register it user-wide for every project
instead.

Any MCP client works with the equivalent stdio config: command `npx`,
args `["-y", "@satoricorp/lgtm"]`.

The CLI install also bundles the server as a standalone `lgtm-mcp` binary in
`~/.local/bin`, so no Node is required:

```bash
curl -fsSL https://download.lgtm.cx/install.sh | sh
```

Building from source:

```bash
cd mcp
bun install
bun run build
```

For cloud auth, log in once with GitHub:

```bash
lgtm auth login
```

The installer provides the `lgtm` CLI and `lgtm-mcp` binary. Repo Git
hooks are installed when a repo is initialized with `lgtm init`; the `lgtm_review`
MCP tool is read-only — it never initializes a repo, never writes `~/.lgtm`,
never touches `.git/index`, and never publishes. The `pre-push` hook runs `lgtm capture push` for
the pushed ref range, stages captured Claude/Codex/Cursor session context in
`~/.lgtm/lgtm.db`, and uploads only when lgtm upload credentials are configured.

Development stdio:

```bash
LGTM_BINARY="$(command -v lgtm)" bun run start
```

## Scripts

| Script | Description |
|--------|-------------|
| `bun run dev` | xmcp dev with watch |
| `bun run build` | Production xmcp output plus standalone `dist/lgtm-mcp` |
| `bun run build:xmcp` | Production xmcp JavaScript output to `dist/` |
| `bun run build:binary` | Standalone stdio binary at `dist/lgtm-mcp` |
| `bun run start` | stdio MCP (`dist/lgtm-mcp`) |
| `bun test` | Tool schema and CLI spawn tests |
| `bun run typecheck` | TypeScript check |

## Publishing

`npm publish` from `mcp/` (or `just mcp-publish` from the repo root) builds,
tests, and publishes `@satoricorp/lgtm`. CI does the same automatically when
a `mcp-v<version>` tag is pushed — the tag must match `package.json`'s version:

```bash
git tag mcp-v0.1.0 && git push origin mcp-v0.1.0
```

CI publishes via npm Trusted Publishing (OIDC) — no token secret. One-time
setup: publish once locally, then on npmjs.com under the package's Settings add
a trusted publisher (GitHub Actions, owner `satoricorp`, repo `lgtm`,
workflow `mcp.yml`). The npm package ships only the bundled `dist/*.js`; the
compiled `dist/lgtm-mcp` binary is distributed by the CLI installer instead.

## Environment

| Variable | Description |
|----------|-------------|
| `LGTM_BINARY` | Optional path to `lgtm` executable |
| `LGTM_CLOUD_URL` | lgtm cloud API base URL for review AI fallback |

Without a `LGTM_BINARY` override, MCP uses `~/.local/bin/lgtm` when present, then falls back to `lgtm` on `PATH`. Cloud calls use credentials from `lgtm auth login` when available.

The released `lgtm-mcp` checks `https://download.lgtm.cx/cli/manifest.json` and adds an
update notice plus `curl -fsSL https://download.lgtm.cx/install.sh | sh` to tool responses
when a newer CLI package is available. Set `LGTM_MCP_UPDATE_CHECK=0` to disable
that check.
