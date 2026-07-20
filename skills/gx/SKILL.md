---
name: gx
description: >-
  Version-control workflow for repos that use GX (the gx CLI). Use this skill whenever
  the user asks to commit, save, record, amend, land, publish, ship, or open a PR for
  work in a repo that has GX initialized (GX hooks/metadata present, GX instructions in
  AGENTS.md, or the user mentions "gx") — even if they only say "save my work" or
  "commit this". In a GX repo this skill overrides generic git habits: record changes
  with `gx commit` (or MCP `gx_commit`) after `git add`, inspect with `gx status`, and
  publish with plain `git push` (not `gx push`), then `gh pr create`. Use git for
  staging, amending, rebasing, and read-only inspection.
---

# GX CLI workflow

GX is a quality layer for agent-written code. It records **revisions** grouped into
**stacks** (one stack = one reviewable line of work, exported as one Git branch and PR),
attaches what the agent did to what changed, and publishes ordinary Git branches and
commits to GitHub.

The commands below mirror the GX MCP tools (`gx_commit`, `gx_status`, `gx_review`), plus
plain `git push` for publishing.

## First: make sure the repo is initialized

GX is initialized when `gx init` has registered the repository (hooks + local metadata).
If commands fail with init errors, initialize non-interactively before other gx work:

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
`HEAD` advances. To start a new line of work from the base branch, mint the stack's
branch at commit time:

```bash
gx commit -m "add exponential backoff helper" -b feature/retry-backoff
```

Exit codes: `0` success, `1` general error, `2` nothing staged (run `git add` first).

### Provenance is captured automatically

`gx init` installs a pre-push hook that runs capture/publication on `git push`: it
records the coding session and attaches that context to the pushed revisions. Just
`git push`; the provenance rides along.

## Inspect state: gx status

```bash
gx status --agent   # stable key=value lines designed for agents to parse
gx status --json    # machine-readable JSON (what MCP gx_status returns)
gx status --all     # include every changed file in the working-tree sections
```

Use status after commits and before a push. One check at a decision point is enough.

## Amend or rewrite recorded work

Use normal Git for history edits. Preserve the GX revision trailer in the commit message
when amending or rebasing so metadata stays linked:

```bash
git commit --amend
git rebase -i <base>
```

## Publish: git push

```bash
git push             # publish the current stack's branch
gh pr create         # then open the PR (when the user wants one)
```

Publish with plain `git push`. The GX pre-push hook captures the session and registers
publish/CI status as the branch goes up. Open the pull request with `gh pr create`; do
**not** seed a `## Summary` section in the PR body (leave human notes only). GX Cloud
appends the rich summary below the existing body once the PR exists.

Do not run `gx push`. It bypasses the pre-push hook path GX is built around.

## Review: gx review

```bash
gx review                                    # patch-focused default
gx review "did we break the retry contract?" # steer with a prompt
gx review --scope security                   # architecture, security, performance, etc.
gx review --focus src/auth --deep            # deep pass limited to a path prefix
gx review --verbose                          # include repo facts, docs, changed files
```

## Do not

- `git commit` for normal GX save flows — use `gx commit` after `git add` so revisions are recorded.
- `gx push` / manual `gx capture push` — bypass the pre-push hook publish path.
- `gx commit -a` / `-F` / path arguments — stage with `git add`.
- Hidden utilities such as `gx edit`, `gx base`, `gx generate`, or `gx push` — removed; use Git-native workflows instead.

## When a command fails

| Symptom | Fix |
| --- | --- |
| `gx commit` exit code 2 | Nothing staged — run `git add <files>` first |
| Error mentions `gx auth login` or "not logged in" | Ask the user to run `gx auth login`, then retry |
| Error mentions session expired / `gx auth logout` | Ask the user to run `gx auth logout` then `gx auth login`, then retry |
| `git push` fails: no `origin` / not a GitHub remote | `git remote -v`; point `origin` at GitHub, or skip publishing |
| `git push` rejected (remote moved) | `gx sync`, then retry `git push` |
| Repo not initialized errors | `gx init -y`, then retry the original command |
