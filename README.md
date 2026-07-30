# Totality

Git for Agents.

Git was not built for LLMs. Git works great for humans, where you can attach a
message to code changes to review later, but with AI we now have additional
artifacts, like session context, that can save work in an organized way and
supplement your pull requests.

Many changes that happen during a coding session are otherwise lost. Those
details can help teammates and other agents understand the intent and decisions
behind the work. Alongside Git history and other resources, Totality gives your team a
stronger review system for moving faster while maintaining high-quality
software.

## Setup

```bash
curl -fsSL https://download.totality.sh/install.sh | sh
```

Re-run the same command to upgrade or repair an existing installation. The
installer replaces only the Totality-managed `tx` and `tx-mcp` binaries and refreshes
the `txr` alias. It does not remove `~/.totality`, repository metadata, or
hooks.

```bash
tx version
tx doctor
```

From your repo:

```bash
tx init
```

This configures your Totality identity, installs Git lifecycle hooks to identify and
record revisions plus a pre-push hook to publish session data, registers the Totality
MCP server, and offers to add Totality workflow instructions to `AGENTS.md`.

To uninstall the CLI:

```bash
rm -f ~/.local/bin/tx ~/.local/bin/txr ~/.local/bin/tx-mcp
```

This removes the installed binaries only. Local Totality data remains in `~/.totality`.

## Set An API Key

Totality reads a model key from the environment:

```bash
ANTHROPIC_API_KEY
OPENAI_API_KEY
```

## Auth

```bash
tx auth login
tx auth status
```

Logging into Totality allows you to push metadata and captured context to Totality Cloud.
This is required to use Totality code review.

## Add Instructions To Your AGENTS.md

```md
Version control: plain Git. Once `tx init` installs the hooks, Totality records and
publishes automatically — there is no Totality save verb.

Save work:
- `git add` to stage.
- `git commit -m "..."` to save. A Totality hook records the commit as a reviewable revision.
- `git status` to inspect; `tx review` for AI review of the current change.

Publish with plain `git push` (the Totality pre-push hook captures the session and publishes),
then open the PR with `gh pr create`. Do not run `tx push` or `tx capture push` — they
bypass or suppress the hook.

When the user says "save work", "save using tx", or "save with tx", stage with
`git add`, save with `git commit`, and publish with plain `git push` unless the user
asks to keep the work local.

To amend, use `git commit --amend` and preserve the Totality revision trailer.

Use `tx_review` (MCP) or `tx review` (CLI) for review context on the current change.

Totality PR summaries are posted for PRs whose branch was pushed through Totality with `git push`
while the pre-push hook is installed. A PR opened before that push will not get a summary
until the branch is pushed through Totality.

If your agent client supports tool policies, require approval for destructive reset and
branch deletion.
```

The installer gives you the `tx` CLI and `tx-mcp`.
When a repo is initialized with `tx init`, Totality
installs `prepare-commit-msg`, `post-commit`, `post-rewrite`, and `pre-push`
hooks. They preserve durable Totality revision IDs across normal Git commits and
rewrites, capture Claude/Codex/Cursor session context into `~/.totality/totality.db`, mark
pushed Totality revisions shareable, and drain uploads in the background when
credentials are configured.

## Basic Workflow

```bash
git add <files>
git commit -m "describe this change"   # a Totality hook records the revision
git status                             # inspect with plain Git
git push                               # push code; Totality hook publishes sessions and PR summaries
```

Useful review commands:

```bash
txr                                              # tx review, patch-focused
tx review --repo "how does capture work?"        # ask about the codebase
tx review --base origin/main --fail-on strong --no-publish   # CI gate
```

`tx review` is read-only: it never runs `tx init`, writes `~/.totality`, or touches
`.git/index`, so it is safe in CI and on a checkout you do not own. Under
`--fail-on` it exits `3` when findings survive and `4` when nothing was
reviewed. See [the reference](docs-site/content/docs/cli.mdx) for the full surface.

## MCP

Totality ships with a stdio MCP server.

Install includes `tx-mcp`:

```bash
command -v tx-mcp
```

Cursor:

```bash
cursor mcp add tx -- env TOTALITY_BINARY=$HOME/.local/bin/tx $HOME/.local/bin/tx-mcp
```

Claude Code:

```bash
claude mcp add tx -- env TOTALITY_BINARY=$HOME/.local/bin/tx $HOME/.local/bin/tx-mcp
```

Then ask your agent:

```text
save work
save using tx
save with tx
```

Expected flow:

```text
git add -> git commit -> git push -> gh pr create
```

MCP exposes `tx_review` only; saving and publishing are plain Git.

Publish with plain `git push` only. Do not run `tx push` or `tx capture push`.
