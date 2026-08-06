# gx

Git for Agents.

Git was not built for LLMs. Git works great for humans, where you can attach a
message to code changes to review later, but with AI we now have additional
artifacts, like session context, that can save work in an organized way and
supplement your pull requests.

Many changes that happen during a coding session are otherwise lost. Those
details can help teammates and other agents understand the intent and decisions
behind the work. Alongside Git history and other resources, gx gives your team a
stronger review system for moving faster while maintaining high-quality
software.

## Setup

```bash
curl -fsSL https://download.gx.run/install.sh | sh
```

Re-run the same command to upgrade or repair an existing installation. The
installer replaces only the gx-managed `gx` and `gx-mcp` binaries and refreshes
the `gxr` alias. It does not remove `~/.gx`, repository metadata, or
hooks.

```bash
gx version
gx doctor
```

From your repo:

```bash
gx init
```

This configures your gx identity, installs Git lifecycle hooks to identify and
record revisions plus a pre-push hook to publish session data, registers the gx
MCP server, and offers to add gx workflow instructions to `AGENTS.md`.

To uninstall the CLI:

```bash
rm -f ~/.local/bin/gx ~/.local/bin/gxr ~/.local/bin/gx-mcp
```

This removes the installed binaries only. Local gx data remains in `~/.gx`.

## Set An API Key

gx reads a model key from the environment:

```bash
ANTHROPIC_API_KEY
OPENAI_API_KEY
```

## Auth

```bash
gx auth login
gx auth status
```

Logging into gx allows you to push metadata and captured context to gx Cloud.
This is required to use gx code review.

## Add Instructions To Your AGENTS.md

```md
Version control: plain Git. Once `gx init` installs the hooks, gx records and
publishes automatically — there is no gx save verb.

Save work:
- `git add` to stage.
- `git commit -m "..."` to save. A gx hook records the commit as a reviewable revision.
- `git status` to inspect; `gx review` for AI review of the current change.

Publish with plain `git push` (the gx pre-push hook captures the session and publishes),
then open the PR with `gh pr create`. Do not run `gx push` or `gx capture push` — they
bypass or suppress the hook.

When the user says "save work", "save using gx", or "save with gx", stage with
`git add`, save with `git commit`, and publish with plain `git push` unless the user
asks to keep the work local.

To amend, use `git commit --amend` and preserve the gx revision trailer.

Use `gx_review` (MCP) or `gx review` (CLI) for review context on the current change.

gx PR summaries are posted for PRs whose branch was pushed through gx with `git push`
while the pre-push hook is installed. A PR opened before that push will not get a summary
until the branch is pushed through gx.

If your agent client supports tool policies, require approval for destructive reset and
branch deletion.
```

The installer gives you the `gx` CLI and `gx-mcp`.
When a repo is initialized with `gx init`, gx
installs `prepare-commit-msg`, `post-commit`, `post-rewrite`, and `pre-push`
hooks. They preserve durable gx revision IDs across normal Git commits and
rewrites, capture Claude/Codex/Cursor session context into `~/.gx/gx.db`, mark
pushed gx revisions shareable, and drain uploads in the background when
credentials are configured.

## Basic Workflow

```bash
git add <files>
git commit -m "describe this change"   # a gx hook records the revision
git status                             # inspect with plain Git
git push                               # push code; gx hook publishes sessions and PR summaries
```

Useful review commands:

```bash
gxr                                              # gx review, patch-focused
gx review --repo "how does capture work?"        # ask about the codebase
gx review --base origin/main --fail-on strong --no-publish   # CI gate
```

`gx review` is read-only: it never runs `gx init`, writes `~/.gx`, or touches
`.git/index`, so it is safe in CI and on a checkout you do not own. Under
`--fail-on` it exits `3` when findings survive and `4` when nothing was
reviewed. See [the reference](docs-site/content/docs/cli.mdx) for the full surface.

## MCP

gx ships with a stdio MCP server.

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
