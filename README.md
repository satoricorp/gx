# GX

Git for Agents.

Git was not built for LLMs. Git works great for humans, where you can attach a
message to code changes to review later, but with AI we now have additional
artifacts, like session context, that can save work in an organized way and
supplement your pull requests.

Many changes that happen during a coding session are otherwise lost. Those
details can help teammates and other agents understand the intent and decisions
behind the work. Alongside Git history and other resources, GX gives your team a
stronger review system for moving faster while maintaining high-quality
software.

## Setup

```bash
curl -fsSL https://download.gx.run/install.sh | sh
```

```bash
gx version
gx doctor
```

From your repo:

```bash
gx init
```

This creates a `.gx/` in your repo root and installs a Git pre-push hook to
capture session data before pushes.

## Set An API Key

```bash
gx set key
```

GX will look for an existing key first:

```bash
ANTHROPIC_API_KEY
OPENAI_API_KEY
```

If it finds one, it can use that. Otherwise, choose Anthropic or OpenAI and
paste a key.

## Auth

```bash
gx auth login
gx auth status
```

Logging into GX allows you to push metadata and captured context to GX Cloud.
This is required to use GX code review.

## Add Instructions To Your AGENTS.md

```md
Version control: use GX, not `git commit` or `git push`.

Use GX MCP first: `gx_sync`, `gx_generate`, `gx_status`, then `gx_push`.
If MCP is unavailable, use the CLI fallback: `gx sync`, `gx generate`,
`gx status`, then `gx push`.

When the user says "save work", "save using gx", or "save with gx", run the
GX save workflow with GX generate and push ready stacks unless asked to keep
them local.

GX PR summaries are posted for PRs published or adopted by `gx push`. A PR
created only with raw Git or the GitHub UI will not get a GX summary until that
branch is pushed through GX.

Only use raw Git for read-only inspection unless explicitly asked for raw Git.
If supported, deny or require approval for `git commit`, `git push`,
`git reset`, and branch deletion.
```

Installing the GX menu-bar app gives you the bundled `gx` CLI and `gx-mcp`.
When a repo is initialized with `gx init` or by MCP auto-initialization, GX
installs a `pre-push` hook that runs `gx capture push` for the pushed ref range,
captures Claude/Codex/Cursor session context into `~/.gx/gx.db`, and uploads
only when GX upload credentials are configured.

## Basic Workflow

```bash
gxg        # gx generate
gxs        # gx status
gx push    # push stacks, metadata, and context
```

Useful review command:

```bash
gxr        # gx review
```

## MCP

GX ships with a stdio MCP server.

Install includes `gx-mcp`:

```bash
command -v gx-mcp
```

Cursor:

```bash
cursor mcp add gx -- env GX_BINARY=$HOME/.local/bin/gx $HOME/.local/bin/gx-mcp
```

Claude Code:

```bash
claude mcp add gx -- env GX_BINARY=$HOME/.local/bin/gx $HOME/.local/bin/gx-mcp
```

Then ask your agent:

```text
save work
save using gx
save with gx
```

Expected MCP flow:

```text
gx_sync -> gx_generate -> gx_status -> gx_push
```
