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
MCP server, installs `/gx` and `/constraints` commands for Claude Code, Codex,
and Cursor, and offers to add gx workflow instructions to `AGENTS.md`.

Most gx commands initialize the repository on first use, so `gx init` is the way
to set your identity and answer the `AGENTS.md` prompt rather than a hard
prerequisite.

To install the hooks once for every repository on the machine instead of
per-repo, use `core.hooksPath`:

```bash
gx init --global
```

To uninstall the CLI:

```bash
rm -f ~/.local/bin/gx ~/.local/bin/gxr ~/.local/bin/gxc ~/.local/bin/gx-mcp
rm -f ~/.local/share/bash-completion/completions/gx ~/.zfunc/_gx
```

This removes the installed binaries and shell completions only. Local gx data
remains in `~/.gx` (or `$GX_HOME`).

## Auth

```bash
gx auth login
gx auth status
gx auth logout
```

Logging into gx allows you to push metadata and captured context to gx Cloud.
This is required to use gx code review: the reviewer runs on Bedrock, and gx
Cloud is what brokers those calls. There is no other provider, and no
`ANTHROPIC_API_KEY`-style local key path — a missing login reads as "no reviewer
ran", never as a clean review.

To review on your own AWS account instead of through gx Cloud:

```bash
export GX_REVIEW_BEDROCK_DIRECT=1
export AWS_ACCESS_KEY_ID=... AWS_SECRET_ACCESS_KEY=...   # plus AWS_SESSION_TOKEN if temporary
```

Override the two review models with `GX_REVIEW_BEDROCK_MODEL_A` and
`GX_REVIEW_BEDROCK_MODEL_B`; set either to `off` to review with a single model.

`OPENAI_API_KEY` is unrelated to the reviewer — it is read only when embedding
this repository for the code index.

## Add Instructions To Your AGENTS.md

`gx init` offers to add this for you. To write it by hand, or to check what init
added, this is the text:

```md
Version control: plain Git. Once `gx init` installs the hooks, gx records and publishes automatically — there is no gx save verb.

Default flow:
- Run `git add` to stage the files for this revision.
- Run `git commit -m "..."` to save. A gx hook records the commit as a reviewable revision.
- Run plain `git push` to publish. The gx pre-push hook captures the session and publishes code changes, sessions, and gx metadata to gx Cloud automatically — do not run `gx push` or `gx capture push` yourself; they bypass the hook.
- Open PRs with `gh pr create` (or the GitHub UI). Do not seed a `## Summary` in the PR body — leave human notes only; gx Cloud appends the rich summary below once the PR exists.
- To amend, use `git commit --amend` and preserve the gx revision trailer in the message.

For AI review, run the `gx_review` MCP tool (or the `gx review` CLI) on the current change.
Before shipping, run the `gx_constraints` MCP tool (or `gx constraints`) to check the change against the pre-ship constraint gates.
```

gx PR summaries are posted for PRs whose branch was pushed through gx with `git
push` while the pre-push hook is installed. A PR opened before that push will not
get a summary until the branch is pushed through gx.

If your agent client supports tool policies, require approval for destructive
reset and branch deletion.

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

`gx review` leaves no gx state on the machine that runs it: it never runs `gx
init`, writes `~/.gx`, or touches `.git/index`, so it is safe in CI and on a
checkout you do not own. It does publish outward by default — posting the PR
review comment and recording review history to gx Cloud. `--no-comment` skips
the comment, and `--no-publish` skips both.

Under `--fail-on` it exits `3` when findings at or above the threshold survive,
`4` when nothing was reviewed at all, and `5` when code was read but the review
that read it ran degraded. Run `gx review --help` for the full flag surface, and
`gx --help` for the rest of the commands.

## MCP

gx ships with a stdio MCP server.

Install includes `gx-mcp`:

```bash
command -v gx-mcp
```

`gx init` registers it automatically with Cursor, Claude Code, and Codex when it
finds them. To register it by hand:

Cursor:

```bash
cursor mcp add gx -- env GX_BINARY=$HOME/.local/bin/gx $HOME/.local/bin/gx-mcp
```

Claude Code:

```bash
claude mcp add gx -- env GX_BINARY=$HOME/.local/bin/gx $HOME/.local/bin/gx-mcp
```

MCP exposes `gx_review` and `gx_constraints`. There is still no save or publish
tool and no gx-specific verb to ask for: the server tells your agent to use
plain Git, and the hooks do the rest.

Expected flow:

```text
git add -> git commit -> git push -> gh pr create
```

`gx init` also installs `/gx` and `/constraints` slash commands for Claude
Code, Codex, and Cursor: `/gx` runs a fast review of the current change, and
`/constraints` runs the pre-ship constraints exit gate.

Publish with plain `git push` only. Do not run `gx push` or `gx capture push`.

## Constraints

`gx constraints` (shortcut `gxc`) is the pre-ship exit gate: it checks the
current change against six constraints — correctness, security, code health,
back-pressure, accessibility, performance — and reports PASS, FAIL, or SKIPPED
for each plus a ship / no-ship verdict. Deterministic checks decide what they
can (the project's tests and linters, a secrets scan over added lines,
dependency audits like govulncheck and npm audit); one AI judgment call covers
the rest, grounded in the same code-index, session, prior-finding, and
knowledge-corpus context `gx review` uses. Run it before opening a PR:

```bash
gx constraints "fix auth timeout"
```

The optional argument states the change's intent, which the back-pressure gate
judges scope against. `--md` prints markdown, `--json` the full report,
`--verbose` lists the files behind each gate, and `--skip-gates` skips named
gates. Exit codes make it a CI gate by default — 0 ship, 3 no-ship, 4 nothing
to check, 5 degraded — and `--report-only` opts out of the coded exits.
