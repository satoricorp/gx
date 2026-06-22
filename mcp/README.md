# GX MCP Server

TypeScript MCP server (xmcp) that runs over stdio and shells to the local `gx` CLI. Local repo, JJ, and Git work stay on the user's machine; cloud review context is reached by outbound HTTPS from `gx` when cloud auth or API-key env is configured.

## Primary workflow

1. `gx_sync` before changes, so remote GitHub merges are reflected locally.
2. `gx_compose` to propose changes. Ready proposals auto-accept by default.
3. `gx_publish` to publish accepted stacks.
4. `gx_review` when codegen needs review context from local facts, previous sessions, PRs, and current code changes.

If a repository is not initialized for GX, MCP runs `gx init` non-interactively
before repository tools continue. It uses Git identity when available and falls
back to `GX_MCP_INIT_NAME` / `GX_MCP_INIT_EMAIL`, then safe placeholder values.

## Tools

| Tool | CLI | Purpose |
|------|-----|---------|
| `gx_sync` | `gx sync` | Sync remote Git state before composing or publishing |
| `gx_compose` | `gx compose --json` | Layer working-copy changes into a pending compose proposal from a session-isolated JJ workspace |
| `gx_accept` | `gx compose apply <proposal-id> --json` | Manually accept a held ready compose proposal |
| `gx_publish` | `gx publish [stack]` | Publish accepted stacks |
| `gx_review` | `gx review` | Gather local review/context with AI reviewers enabled |
| `gx_fix` | `gx compose fix <proposal-id> --json` | Repair compose proposal issues with deterministic repair plus LLM repair |
| `gx_set_base` | `gx base --set <default> --json` | Return the GX authoring base to the repo default branch only |

## `gx_compose` actions

| `action` | CLI | Purpose |
|----------|-----|---------|
| `propose` (default) | `gx compose --json` | Layer working-copy changes into the pending compose proposal; auto-accepts ready proposals unless `auto_accept=false` |
| `review-plan` | `gx compose review-plan --json --plan-file …` | Validate an LLM-authored proposal without applying |
| `apply` | `gx compose apply-plan --json --plan-file …` | Accept a reviewed proposal into GX stacks |

Proposal JSON may include `hunk_links` capture evidence (empty when matcher has no sessions).

### Session workspaces

`gx_compose` requires a session id. Pass `session_id` / `session_ids`, or set `GX_SESSION_ID` / `GX_SESSION_IDS`.

Before shelling out to `gx`, the MCP server resolves the requested repo, creates or reuses a JJ workspace at:

```bash
$GX_HOME/workspaces/<repo-hash>/<session-hash>
```

Then it runs `gx compose` from that workspace cwd. This keeps concurrent MCP sessions from composing or applying each other's dirty checkout changes.

Set `GX_MCP_WORKSPACE_ROOT` to override the workspace directory.

## Install

```bash
cd mcp
bun install
bun run build
```

Register stdio MCP (example):

```bash
cursor mcp add gx -- env GX_BINARY=$HOME/.local/bin/gx /Applications/GX.app/Contents/Resources/bin/gx-mcp
```

For cloud auth, log in once with GitHub:

```bash
gx auth login
```

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
| `GX_SESSION_ID` / `GX_SESSION_IDS` | Session ids used to isolate compose workspaces |
| `GX_MCP_WORKSPACE_ROOT` | Override the root directory for MCP session workspaces |
| `GX_CLOUD_URL` | GX cloud API base URL for review and compose AI fallback |
| `GX_MCP_INIT_NAME` / `GX_MCP_INIT_EMAIL` | Optional identity used when MCP auto-runs `gx init` |
| `GX_REVIEW_CONTEXT_URL` / `GX_REVIEW_CONTEXT_TOKEN` | Optional indexed review-context endpoint and token |

Without a `GX_BINARY` override, MCP uses `~/.local/bin/gx` when present, then falls back to `gx` on `PATH`. Cloud calls use credentials from `gx auth login` when available.
