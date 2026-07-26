---
name: gx
description: >-
  Version-control workflow for repos that use GX (the gx CLI). Use this skill whenever
  the user asks to commit, save, record, amend, land, publish, ship, or open a PR for
  work in a repo that has GX initialized (GX hooks/metadata present, GX instructions in
  AGENTS.md, or the user mentions "gx") — even if they only say "save my work" or
  "commit this". In a GX repo you work in plain Git: save with `git add` + `git commit`
  (GX's hooks record the commit as a reviewable revision), inspect with `git status`, and
  publish with `git push` (never `gx push`), then `gh pr create`. GX adds AI review via
  `gx review`; use it for review, not for saving.
---

# GX workflow

GX is a quality layer for agent-written code. Its Git hooks record your commits as
**revisions** grouped into **stacks** (one stack = one reviewable line of work, which is
just a Git branch and its PR), attach what the agent did to what changed, and publish to
GitHub on `git push`.

You work in plain Git — `git add`, `git commit`, `git push`, `git checkout -b`. There is
no GX save verb; the hooks do the recording and publishing automatically. GX MCP exposes
`gx_review` only — save and publish with Git, and use `gx review` for AI review.

## First: make sure the repo is initialized

GX is initialized when `gx init` has registered the repository (hooks + local metadata).
If Git or gx commands fail with init errors, initialize non-interactively first:

```bash
gx init -y
```

Identity defaults to the repo's Git config; pass `--name` / `--email` only when Git
identity is missing. Run once per repository.

## Default save workflow

When the user says "save", "commit", "save my work", or finishes a task in a GX repo:

```bash
git add <files>                        # 1. choose scope — one logical change per revision
git commit -m "describe this change"   # 2. a GX hook records the commit as a revision
git status                             # 3. confirm the tree is clean and the branch is right
git push                               # 4. publish — the GX pre-push hook does the rest
gh pr create                           # 5. open the PR (only when the user wants one)
```

`git commit` records the change on the current branch and advances `HEAD`; GX's
`prepare-commit-msg` and `post-commit` hooks record it as a **revision** behind the
scenes, leaving Git's output unchanged (`[branch shorthash] message`). To start a new
line of work (a new stack), branch first with plain Git:

```bash
git checkout -b feature/retry-backoff
git add <files>
git commit -m "add exponential backoff helper"
```

If nothing is staged, `git commit` reports "nothing to commit" — run `git add <files>` first.

### Provenance is captured automatically

`gx init` installs a pre-push hook that runs capture/publication on `git push`: it records
the coding session and matches it to the hunks each revision changed, so the review bundle
knows which session produced each change. Just `git push` — the provenance rides along.
There is no manual attach or self-report step; do not try to declare provenance by hand.

## Inspect state: plain Git

A stack is just a Git branch and a revision is just a commit on it — inspect both with Git:

```bash
git status              # working-tree state: staged/unstaged files, current branch
git log --oneline -10   # recent revisions (commits) on this stack (branch)
```

Check state after commits and before a push. One check at a decision point is enough.

## Amend or rewrite recorded work

Use normal Git for history edits. The `post-rewrite` hook follows amended and rebased
commits to their new OIDs; keep the GX revision trailer in the message so metadata stays
linked:

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

Do not run `gx push` — the command has been removed. Publish with plain `git push`; the
pre-push hook is the only publish path.

## Review: gx review

```bash
gx review                                    # patch-focused default
gx review "did we break the retry contract?" # steer with a prompt
gx review --scope security                   # architecture, security, performance, etc.
gx review --focus src/auth --deep            # deep pass limited to a path prefix
gx review --verbose                          # include repo facts, docs, changed files
```

The `gxr` alias runs `gx review`. The `gx_review` MCP tool is the same review from an
agent client.

## Do not

- `gx commit` / `gx status` / `gx push` / `gx edit` — removed. Use `git commit`,
  `git status`, `git push`, and `git commit --amend`; GX's hooks record revisions for you.
- `gx capture push` or any manual capture/attach step — the pre-push hook captures and
  publishes on `git push`, and provenance is inferred automatically.
- Hidden utilities such as `gx base` or `gx generate` — removed; use plain Git
  (`git pull --rebase`, `git checkout -b`) instead. (`gx sync` still exists, but it
  retries queued uploads; it is not part of the save flow.)
- Declaring provenance by hand (self-reports, task summaries) — there is no such channel;
  sessions are matched to changed hunks automatically.

## When a command fails

| Symptom | Fix |
| --- | --- |
| `git commit` says "nothing to commit" | Nothing staged — run `git add <files>` first |
| Error mentions `gx auth login` or "not logged in" | Ask the user to run `gx auth login`, then retry |
| Error mentions session expired / `gx auth logout` | Ask the user to run `gx auth logout` then `gx auth login`, then retry |
| `git push` fails: no `origin` / not a GitHub remote | `git remote -v`; point `origin` at GitHub, or skip publishing |
| `git push` rejected (remote moved) | `git pull --rebase`, then retry `git push` |
| No PR summary on the PR | Confirm the branch was pushed with `git push` while hooks are installed; run `gx doctor` |
| Repo not initialized errors | `gx init -y`, then retry |
