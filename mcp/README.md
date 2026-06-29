# GX MCP Server

TypeScript MCP server (xmcp) that runs over stdio and shells to the local `gx` CLI. Local repo, JJ, and Git work stay on the user's machine; cloud review context is reached by outbound HTTPS from `gx` when cloud auth or API-key env is configured.

## Primary workflow

1. `gx_sync` before changes, so remote GitHub merges are reflected locally.
2. `gx_generate` to generate local features and revisions. Safe generated revisions apply automatically and may append to semantically similar stacks.
3. `gx_status` to inspect generated local features, revisions, and remote state.
4. `gx_push` to push generated features to the remote.
5. `gx_review` when codegen needs review context from local facts, previous sessions, PRs, current code changes, and optional prompt guidance.

If `gx_push` reports remote divergence, run `gx_sync`, resolve the divergence, then retry `gx_push` for that stack.

When a user says "save work", "save using gx", or "save with gx", treat that as
a request to run the GX save workflow: `gx_generate`, then `gx_status`, then
`gx_push` for ready stacks unless the user explicitly asks to keep the work
local.

If a repository is not initialized for GX, MCP runs `gx init` non-interactively
before repository tools continue. It uses Git identity when available and falls
back to `GX_MCP_INIT_NAME` / `GX_MCP_INIT_EMAIL`, then safe placeholder values.

## Agent instructions

Add a block like this to the start or end of `AGENTS.md` or `CLAUDE.md` in repos
where agents should use GX:

```md
Version control: use GX, not `git commit` or `git push`.

Use GX MCP first:
- `gx_sync` before generating or pushing when remote changes may have landed.
- `gx_generate` to save work into GX revisions and stacks.
- `gx_status` to inspect local and remote stack state.
- `gx_push` to publish ready GX stacks.

If MCP is unavailable, use the CLI fallback:
- `gx sync`
- `gx generate`
- `gx status`
- `gx push`

When the user says "save work", "save using gx", or "save with gx", run the
GX save workflow: generate the work with GX, inspect status, and push ready
stacks unless the user asks to keep them local.

Only use raw Git for read-only inspection unless the user explicitly asks for
raw Git. If your agent client supports tool policies, deny or require approval
for `git commit`, `git push`, `git reset`, and branch deletion.
```

## Tools

| Tool | CLI | Purpose |
|------|-----|---------|
| `gx_sync` | `gx sync` | Sync remote Git and GX remote state before generating or pushing |
| `gx_generate` | `gx generate --json` | Generate local features and revisions from a session-isolated JJ workspace; may append to similar stacks |
| `gx_status` | `gx status --json` | Inspect unstaged files, local features, revisions, and remote state |
| `gx_push` | `gx push [stack]` | Push generated features, sessions, metadata, and guarded rewrites |
| `gx_review` | `gx review [prompt]` | Gather local review/context with AI reviewers enabled |
| `gx_set_base` | `gx base --set <default> --json` | Return the GX authoring base to the repo default branch only |

### Session workspaces

`gx_generate` requires a session id. Pass `session_id` / `session_ids`, or set `GX_SESSION_ID` / `GX_SESSION_IDS`.

Before shelling out to `gx`, the MCP server resolves the requested repo, creates or reuses a JJ workspace at:

```bash
$GX_HOME/workspaces/<repo-hash>/<session-hash>
```

Then it runs `gx generate` from that workspace cwd. This keeps concurrent MCP sessions from generating each other's dirty checkout changes.

Set `GX_MCP_WORKSPACE_ROOT` to override the workspace directory.

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
| `GX_SESSION_ID` / `GX_SESSION_IDS` | Session ids used to isolate generate workspaces |
| `GX_MCP_WORKSPACE_ROOT` | Override the root directory for MCP session workspaces |
| `GX_CLOUD_URL` | GX cloud API base URL for review and generate AI fallback |
| `GX_MCP_INIT_NAME` / `GX_MCP_INIT_EMAIL` | Optional identity used when MCP auto-runs `gx init` |
| `GX_REVIEW_CONTEXT_URL` / `GX_REVIEW_CONTEXT_TOKEN` | Optional indexed review-context endpoint and token |

Without a `GX_BINARY` override, MCP uses `~/.local/bin/gx` when present, then falls back to `gx` on `PATH`. Cloud calls use credentials from `gx auth login` when available.

The released `gx-mcp` checks `https://download.gx.run/cli/manifest.json` and adds an
update notice plus `curl -fsSL https://download.gx.run/install.sh | sh` to tool responses
when a newer CLI package is available. Set `GX_MCP_UPDATE_CHECK=0` to disable
that check.
