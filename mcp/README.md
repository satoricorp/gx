# GX MCP Server

TypeScript MCP server (xmcp) that runs over stdio and shells to the local `gx` CLI. Local repo, JJ, and Git work stay on the user's machine; cloud review context is reached by outbound HTTPS from `gx` when cloud auth or API-key env is configured.

## Primary workflow

1. `gx_commit` after `git add` to record staged work as a GX revision (default verb).
2. `gx_status` to inspect staged files, local features, revisions, and remote state.
3. `gx_push` when the stack is ready to publish.
4. `gx_edit` to re-enter an existing revision when amending or continuing work.
5. `gx_review` when codegen needs review context from local facts, previous sessions, PRs, current code changes, and optional prompt guidance.

When a user says "save work", "save using gx", or "save with gx", treat that as
a request to run the GX save workflow: `git add`, then `gx_commit`, then
`gx_status`, then `gx_push` for ready stacks unless the user explicitly asks to keep the work
local.

GX PR summaries are posted for PRs published or adopted by `gx_push`. A PR
created only with raw Git or the GitHub UI will not get a GX summary until that
branch is pushed through GX.

If a repository is not initialized for GX, MCP runs `gx init` non-interactively
before repository tools continue. It uses Git identity when available and falls
back to `GX_MCP_INIT_NAME` / `GX_MCP_INIT_EMAIL`, then safe placeholder values.

## Agent instructions

Add a block like this to the start or end of `AGENTS.md` or `CLAUDE.md` in repos
where agents should use GX:

```md
Version control: use GX, not `git commit` or `git push`.

Use GX MCP first:
- `gx_commit` after `git add` to record staged work (default verb).
- `gx_status` to inspect local and remote stack state.
- `gx_push` to publish ready GX stacks.
- `gx_edit` only to re-enter an existing revision for further edits.

If MCP is unavailable, use the CLI fallback:
- `git add`
- `gx commit -m "..."`
- `gx status`
- `gx push`

When the user says "save work", "save using gx", or "save with gx", run the
GX save workflow: stage with `git add`, record with `gx_commit`, inspect with
`gx_status`, and push ready stacks with `gx_push` unless the user asks to keep them local.

GX PR summaries are posted for PRs published or adopted by `gx_push`. A PR
created only with raw Git or the GitHub UI will not get a GX summary until that
branch is pushed through GX.

Only use raw Git for read-only inspection unless the user explicitly asks for
raw Git. If your agent client supports tool policies, deny or require approval
for `git commit`, `git push`, `git reset`, and branch deletion.
```

## Tools

| Tool | CLI | Purpose |
|------|-----|---------|
| `gx_commit` | `gx commit -m "..."` | Record staged Git changes as a GX revision (default verb) |
| `gx_status` | `gx status --json` | Inspect unstaged files, local features, revisions, and remote state |
| `gx_edit` | `gx edit <rev>` | Re-enter an existing revision to keep working on it |
| `gx_push` | `gx push [stack]` | Push stacks, sessions, metadata, and guarded rewrites |
| `gx_review` | `gx review [prompt]` | Gather local review/context with AI reviewers enabled |

## Install

Released CLI-only install:

```bash
curl -fsSL https://download.gx.run/install.sh | sh
```

That package installs both `gx` and the bundled `gx-mcp` stdio server to
`~/.local/bin`.

```bash
cd mcp
bun install
bun run build
```

Register stdio MCP (example):

```bash
cursor mcp add gx -- env GX_BINARY=$HOME/.local/bin/gx $HOME/.local/bin/gx-mcp
```

For cloud auth, log in once with GitHub:

```bash
gx auth login
```

The GX menu-bar app installs the bundled `gx` CLI and `gx-mcp` binary. Repo Git
hooks are installed when a repo is initialized with `gx init` or by MCP
auto-initialization. The installed `pre-push` hook runs `gx capture push` for
the pushed ref range, stages captured Claude/Codex/Cursor session context in
`~/.gx/gx.db`, and uploads only when GX upload credentials are configured.

Development stdio:

```bash
GX_BINARY="$PWD/../apps/menubar/bin/gx" bun run start
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

## Environment

| Variable | Description |
|----------|-------------|
| `GX_BINARY` | Optional path to `gx` executable |
| `JJ_BINARY` | Path to `jj` executable (default: `jj` on PATH) |
| `GIT_BINARY` | Path to `git` executable (default: `git` on PATH) |
| `GX_CLOUD_URL` | GX cloud API base URL for review AI fallback |
| `GX_MCP_INIT_NAME` / `GX_MCP_INIT_EMAIL` | Optional identity used when MCP auto-runs `gx init` |
| `GX_REVIEW_CONTEXT_URL` / `GX_REVIEW_CONTEXT_TOKEN` | Optional indexed review-context endpoint and token |

Without a `GX_BINARY` override, MCP uses `~/.local/bin/gx` when present, then falls back to `gx` on `PATH`. Cloud calls use credentials from `gx auth login` when available.

The released `gx-mcp` checks `https://download.gx.run/cli/manifest.json` and adds an
update notice plus `curl -fsSL https://download.gx.run/install.sh | sh` to tool responses
when a newer CLI package is available. Set `GX_MCP_UPDATE_CHECK=0` to disable
that check.
