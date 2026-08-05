# lgtm

Git for Agents.

Git was not built for LLMs. Git works great for humans, where you can attach a
message to code changes to review later, but with AI we now have additional
artifacts, like session context, that can save work in an organized way and
supplement your pull requests.

Many changes that happen during a coding session are otherwise lost. Those
details can help teammates and other agents understand the intent and decisions
behind the work. Alongside Git history and other resources, lgtm gives your team a
stronger review system for moving faster while maintaining high-quality
software.

## Setup

```bash
curl -fsSL https://download.lgtm.cx/install.sh | sh
```

Re-run the same command to upgrade or repair an existing installation. The
installer replaces only the lgtm-managed `lgtm` and `lgtm-mcp` binaries and refreshes
the `lgtmr` alias. It does not remove `~/.lgtm`, repository metadata, or
hooks.

```bash
lgtm version
lgtm doctor
```

From your repo:

```bash
lgtm init
```

This configures your lgtm identity, installs Git lifecycle hooks to identify and
record revisions plus a pre-push hook to publish session data, registers the lgtm
MCP server, and offers to add lgtm workflow instructions to `AGENTS.md`.

To uninstall the CLI:

```bash
rm -f ~/.local/bin/lgtm ~/.local/bin/lgtmr ~/.local/bin/lgtm-mcp
```

This removes the installed binaries only. Local lgtm data remains in `~/.lgtm`.

## Set An API Key

lgtm reads a model key from the environment:

```bash
ANTHROPIC_API_KEY
OPENAI_API_KEY
```

## Auth

```bash
lgtm auth login
lgtm auth status
```

Logging into lgtm allows you to push metadata and captured context to lgtm Cloud.
This is required to use lgtm code review.

## Add Instructions To Your AGENTS.md

```md
Version control: plain Git. Once `lgtm init` installs the hooks, lgtm records and
publishes automatically — there is no lgtm save verb.

Save work:
- `git add` to stage.
- `git commit -m "..."` to save. A lgtm hook records the commit as a reviewable revision.
- `git status` to inspect; `lgtm review` for AI review of the current change.

Publish with plain `git push` (the lgtm pre-push hook captures the session and publishes),
then open the PR with `gh pr create`. Do not run `lgtm push` or `lgtm capture push` — they
bypass or suppress the hook.

When the user says "save work", "save using lgtm", or "save with lgtm", stage with
`git add`, save with `git commit`, and publish with plain `git push` unless the user
asks to keep the work local.

To amend, use `git commit --amend` and preserve the lgtm revision trailer.

Use `lgtm_review` (MCP) or `lgtm review` (CLI) for review context on the current change.

lgtm PR summaries are posted for PRs whose branch was pushed through lgtm with `git push`
while the pre-push hook is installed. A PR opened before that push will not get a summary
until the branch is pushed through lgtm.

If your agent client supports tool policies, require approval for destructive reset and
branch deletion.
```

The installer gives you the `lgtm` CLI and `lgtm-mcp`.
When a repo is initialized with `lgtm init`, lgtm
installs `prepare-commit-msg`, `post-commit`, `post-rewrite`, and `pre-push`
hooks. They preserve durable lgtm revision IDs across normal Git commits and
rewrites, capture Claude/Codex/Cursor session context into `~/.lgtm/lgtm.db`, mark
pushed lgtm revisions shareable, and drain uploads in the background when
credentials are configured.

## Basic Workflow

```bash
git add <files>
git commit -m "describe this change"   # a lgtm hook records the revision
git status                             # inspect with plain Git
git push                               # push code; lgtm hook publishes sessions and PR summaries
```

Useful review commands:

```bash
lgtmr                                              # lgtm review, patch-focused
lgtm review --repo "how does capture work?"        # ask about the codebase
lgtm review --base origin/main --fail-on strong --no-publish   # CI gate
```

`lgtm review` is read-only: it never runs `lgtm init`, writes `~/.lgtm`, or touches
`.git/index`, so it is safe in CI and on a checkout you do not own. Under
`--fail-on` it exits `3` when findings survive and `4` when nothing was
reviewed. See [the reference](docs-site/content/docs/cli.mdx) for the full surface.

## MCP

lgtm ships with a stdio MCP server.

Install includes `lgtm-mcp`:

```bash
command -v lgtm-mcp
```

Cursor:

```bash
cursor mcp add lgtm -- env LGTM_BINARY=$HOME/.local/bin/lgtm $HOME/.local/bin/lgtm-mcp
```

Claude Code:

```bash
claude mcp add lgtm -- env LGTM_BINARY=$HOME/.local/bin/lgtm $HOME/.local/bin/lgtm-mcp
```

Then ask your agent:

```text
save work
save using lgtm
save with lgtm
```

Expected flow:

```text
git add -> git commit -> git push -> gh pr create
```

MCP exposes `lgtm_review` only; saving and publishing are plain Git.

Publish with plain `git push` only. Do not run `lgtm push` or `lgtm capture push`.
