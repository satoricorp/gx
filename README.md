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

Re-run the same command to upgrade or repair an existing installation. The
installer replaces only the GX-managed `gx` and `gx-mcp` binaries and refreshes
the `gxr` alias. It does not remove `~/.gx`, repository metadata, or
hooks. `--with-menubar` similarly replaces the existing GX app bundle.

```bash
gx version
gx doctor
```

From your repo:

```bash
gx init
```

This configures your GX identity, installs Git lifecycle hooks to identify and
record revisions plus a pre-push hook to publish session data, registers the GX
MCP server, and offers to add GX workflow instructions to `AGENTS.md`.

To uninstall the CLI and menu-bar app:

```bash
rm -f ~/.local/bin/gx ~/.local/bin/gxr ~/.local/bin/gx-mcp
rm -rf /Applications/GX.app ~/Applications/GX.app
```

This removes the installed binaries and app bundle only. Local GX data remains
in `~/.gx`.

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
Version control: plain Git. Once `gx init` installs the hooks, GX records and
publishes automatically — there is no GX save verb.

Save work:
- `git add` to stage.
- `git commit -m "..."` to save. A GX hook records the commit as a reviewable revision.
- `git status` to inspect; `gx review` for AI review of the current change.

Publish with plain `git push` (the GX pre-push hook captures the session and publishes),
then open the PR with `gh pr create`. Do not run `gx push` or `gx capture push` — they
bypass or suppress the hook.

When the user says "save work", "save using gx", or "save with gx", stage with
`git add`, save with `git commit`, and publish with plain `git push` unless the user
asks to keep the work local.

To amend, use `git commit --amend` and preserve the GX revision trailer.

Use `gx_review` (MCP) or `gx review` (CLI) for review context on the current change.

GX PR summaries are posted for PRs whose branch was pushed through GX with `git push`
while the pre-push hook is installed. A PR opened before that push will not get a summary
until the branch is pushed through GX.

If your agent client supports tool policies, require approval for destructive reset and
branch deletion.
```

Installing the GX menu-bar app gives you the bundled `gx` CLI and `gx-mcp`.
When a repo is initialized with `gx init` or by MCP auto-initialization, GX
installs `prepare-commit-msg`, `post-commit`, `post-rewrite`, and `pre-push`
hooks. They preserve durable GX revision IDs across normal Git commits and
rewrites, capture Claude/Codex/Cursor session context into `~/.gx/gx.db`, mark
pushed GX revisions shareable, and drain uploads in the background when
credentials are configured.

## Basic Workflow

```bash
git add <files>
git commit -m "describe this change"   # a GX hook records the revision
git status                             # inspect with plain Git
git push                               # push code; GX hook publishes sessions and PR summaries
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

Expected flow:

```text
git add -> git commit -> git push -> gh pr create
```

MCP exposes `gx_review` only; saving and publishing are plain Git.

Publish with plain `git push` only. Do not run `gx push` or `gx capture push`.
