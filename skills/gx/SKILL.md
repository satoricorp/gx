---
name: gx
description: >-
  Version-control workflow for repos that use GX (the gx CLI). Use this skill whenever
  the user asks to commit, save, record, amend, land, publish, ship, or open a PR for
  work in a repo that has GX initialized (a `.jj` directory at the repo root, GX
  instructions in AGENTS.md, or the user mentions "gx") — even if they only say "save my
  work" or "commit this". In a GX repo this skill overrides generic git habits: record
  changes with `gx commit` (not `git commit`), and publish with plain `git push` (not
  `gx push`), then `gh pr create`. Use git only for staging (`git add`) and read-only
  inspection.
---

# GX CLI workflow

GX is a quality layer for agent-written code. Instead of opaque commits it records
**revisions** grouped into **stacks** (one stack = one reviewable line of work, exported
as one Git branch and PR), attaches what the agent did next to what changed, and
publishes ordinary Git branches and commits to GitHub. Revisions are JJ-backed, so
recorded work stays editable through GX afterward.

The commands below mirror the GX MCP tools (`gx_commit`, `gx_status`, `gx_edit`,
`gx_review`), plus plain `git push` for publishing. Follow the same call shapes so behavior
is identical whether or not MCP is registered.

## First: make sure the repo is initialized

GX is initialized when the repo root has a `.jj` directory. If it does not, initialize
non-interactively before any other gx command (this mirrors what the MCP server does
automatically):

```bash
gx init -y
```

Identity defaults to the repo's Git config; pass `--name` / `--email` only when Git
identity is missing. Run once per repository.

## Default save workflow

When the user says "save", "commit", "save my work", or finishes a task in a GX repo:

```bash
git add <files>                       # 1. choose scope — one logical change per revision
gx commit -m "describe this change"   # 2. record staged changes as a GX revision
gx status --agent                     # 3. verify the stack looks right
git push                              # 4. publish — the GX pre-push hook does the rest
gh pr create                          # 5. open the PR (only when the user wants one)
```

`gx commit` behaves like `git commit`: the revision lands on the current branch and
`HEAD` does not move. To start a new line of work from the base branch, mint the stack's
branch at commit time instead of creating one with git:

```bash
gx commit -m "add exponential backoff helper" -b feature/retry-backoff
```

Exit codes: `0` success, `1` general error, `2` nothing staged (run `git add` first).

### Provenance is captured automatically

You don't need to hand-declare what you did. `gx init` installs a pre-push hook that runs
`gx capture push` on `git push`: it records the coding session — prompts, commands, and
transcript — and attaches that context to the pushed revisions, so reviewers see it next to
the diff with no extra step. Just `git push`; the provenance rides along.

## Inspect state: gx status

```bash
gx status --agent   # stable key=value lines designed for agents to parse
gx status --json    # machine-readable JSON (what MCP gx_status returns)
gx status --all     # include every changed file in the working-tree sections
```

Use status to answer "what is staged, what stacks exist, what will push do" — after a
commit or edit, and before a push. Don't wrap every command in status calls; one check
at a decision point is enough.

## Continue or amend a recorded revision: gx edit

`gx commit --amend` is intentionally unsupported. To keep working on an already
recorded revision:

```bash
gx status --agent        # find the revision id (e.g. tplnrmnorkln)
gx edit <revision-id>    # re-enter it; checkout moves onto that stack's branch
# ...make the follow-up changes...
git add <files>
gx commit -m "describe the follow-up"   # lands on the same stack
gx base --set main       # return to the base branch before unrelated work
```

Accepted ids: the GX change id, Git commit id, or the short id shown by status. After
finishing, always return to the base branch (`gx base --set main`, or the repo's actual
base) — leaving the checkout in edit mode makes the next unrelated commit land on the
wrong stack.

## Publish: git push

```bash
git push             # publish the current stack's branch
gh pr create         # then open the PR (when the user wants one)
```

Publish with plain `git push`. The GX pre-push hook installed by `gx init` captures the
session and registers publish/CI status as the branch goes up — so `git push` is the full
GX publish path. Open the pull request with `gh pr create`; GX posts its PR summary once the
PR exists. Requires a GitHub `origin` remote. Skip publishing when the user asked to keep
work local.

Do not run `gx push`. It takes a separate publish path that suppresses the pre-push hook,
so the session capture the hook performs doesn't happen. Plain `git push` is the path GX is
built around.

## Review: gx review

Run before pushing when the user asks for a review, or when you want findings on your
own changes:

```bash
gx review                                    # patch-focused default
gx review "did we break the retry contract?" # steer with a prompt
gx review --scope security                   # architecture, security, performance,
                                             # onboarding, docs, dependencies,
                                             # testing, maintainability
gx review --focus src/auth --deep            # deep pass limited to a path prefix
gx review --verbose                          # include repo facts, docs, changed files
```

Findings print as markdown and are recorded to review history. AI reviewers resolve
credentials automatically (user key, then GX Cloud); the review always completes even
without them.

## Do not

- `git commit` — the change is never recorded as a GX revision. Use `gx commit` after
  `git add`.
- `gx push` / `gx capture push` — both bypass or suppress the pre-push hook that GX's
  publish path relies on. Publish with plain `git push`.
- `gx commit -a` / `-F` / path arguments — rejected by design; stage with `git add`.
- `gx commit --amend` — use `gx edit <rev>` as above.
- Creating branches with `git switch -c` for new work — use `gx commit -b <branch>`.
- Raw `jj` commands in the normal flow — GX manages the JJ layer; reach for `jj` only
  in explicit repair work, and say why.

## When a command fails

| Symptom | Fix |
| --- | --- |
| `gx commit` exit code 2 | Nothing staged — run `git add <files>` first |
| Error mentions `gx auth login` or "not logged in" | Ask the user to run `gx auth login` in a terminal, then retry |
| Error mentions session expired / `gx auth logout` | Ask the user to run `gx auth logout` then `gx auth login`, then retry |
| `git push` fails: no `origin` / not a GitHub remote | `git remote -v`; point `origin` at GitHub, or skip publishing |
| `git push` rejected (remote moved: branch merged/closed or diverged) | `gx sync`, then retry `git push` |
| Repo not initialized errors | `gx init -y`, then retry the original command |
