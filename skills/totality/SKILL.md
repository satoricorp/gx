---
name: tl
description: >-
  Version-control workflow for repos that use Totality (the tl CLI). Use this skill whenever
  the user asks to commit, save, record, amend, land, publish, ship, or open a PR for
  work in a repo that has Totality initialized (Totality hooks/metadata present, Totality instructions in
  AGENTS.md, or the user mentions "tl") — even if they only say "save my work" or
  "commit this". In a Totality repo you work in plain Git: save with `git add` + `git commit`
  (Totality's hooks record the commit as a reviewable revision), inspect with `git status`, and
  publish with `git push` (never `tl push`), then `gh pr create`. Totality adds AI review via
  `tl review`; use it for review, not for saving.
---

# Totality workflow

Totality is a quality layer for agent-written code. Its Git hooks record your commits as
**revisions** grouped into **stacks** (one stack = one reviewable line of work, which is
just a Git branch and its PR), attach what the agent did to what changed, and publish to
GitHub on `git push`.

You work in plain Git — `git add`, `git commit`, `git push`, `git checkout -b`. There is
no Totality save verb; the hooks do the recording and publishing automatically. Totality MCP exposes
`tl_review` only — save and publish with Git, and use `tl review` for AI review.

## First: make sure the repo is initialized

Totality is initialized when `tl init` has registered the repository (hooks + local metadata).
If Git or tl commands fail with init errors, initialize non-interactively first:

```bash
tl init -y
```

Identity defaults to the repo's Git config; pass `--name` / `--email` only when Git
identity is missing. Run once per repository.

## Default save workflow

When the user says "save", "commit", "save my work", or finishes a task in a Totality repo:

```bash
git add <files>                        # 1. choose scope — one logical change per revision
git commit -m "describe this change"   # 2. a Totality hook records the commit as a revision
git status                             # 3. confirm the tree is clean and the branch is right
git push                               # 4. publish — the Totality pre-push hook does the rest
gh pr create                           # 5. open the PR (only when the user wants one)
```

`git commit` records the change on the current branch and advances `HEAD`; Totality's
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

`tl init` installs a pre-push hook that runs capture/publication on `git push`: it records
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
commits to their new OIDs; keep the Totality revision trailer in the message so metadata stays
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

Publish with plain `git push`. The Totality pre-push hook captures the session and registers
publish/CI status as the branch goes up. Open the pull request with `gh pr create`; do
**not** seed a `## Summary` section in the PR body (leave human notes only). Totality Cloud
appends the rich summary below the existing body once the PR exists.

Do not run `tl push` — the command has been removed. Publish with plain `git push`; the
pre-push hook is the only publish path.

## Review: tl review

```bash
tl review                                    # patch-focused default
tl review "did we break the retry contract?" # steer with a prompt
tl review --repo                             # review the codebase, not just the current change
tl review --scope security                   # architecture, security, performance, etc.
tl review --focus src/auth --deep            # deep pass limited to a path prefix
tl review --verbose                          # include repo facts, docs, changed files
```

Without `--repo`, a dirty working tree is the review subject: tl reviews the diff. Use
`--repo` to review the repository itself — the reviewer is given an inventory of the
repository and its source files, in any language, and findings are no longer filtered
down to changed lines. The uncommitted work stays in focus, so the review still sees
what you just edited, and static checks run over the repository when there is no diff.
`--repo` also wins over `--base`.

The `tlr` alias runs `tl review`. The `tl_review` MCP tool is the same review from an
agent client.

## Do not

- `tl commit` / `tl status` / `tl push` / `tl edit` — removed. Use `git commit`,
  `git status`, `git push`, and `git commit --amend`; Totality's hooks record revisions for you.
- `tl capture push` or any manual capture/attach step — the pre-push hook captures and
  publishes on `git push`, and provenance is inferred automatically.
- Hidden utilities such as `tl base`, `tl generate`, or `tl sync` — removed; use plain
  Git (`git pull --rebase`, `git checkout -b`) instead. Queued uploads retry on the next
  `git push`, and `tl doctor` shows and drains the upload outbox.
- `tl demo` and `tl ops` — removed. There is no walkthrough, and `tl ops diagnose doctor`
  was only ever `tl doctor` under another name; run `tl doctor`.
- `tl report` — folded into `tl doctor --report`, which diagnoses first and then sends
  that diagnosis with recent logs to support.
- Declaring provenance by hand (self-reports, task summaries) — there is no such channel;
  sessions are matched to changed hunks automatically.

## When a command fails

| Symptom | Fix |
| --- | --- |
| `git commit` says "nothing to commit" | Nothing staged — run `git add <files>` first |
| Error mentions `tl auth login` or "not logged in" | Ask the user to run `tl auth login`, then retry |
| Error mentions session expired / `tl auth logout` | Ask the user to run `tl auth logout` then `tl auth login`, then retry |
| `git push` fails: no `origin` / not a GitHub remote | `git remote -v`; point `origin` at GitHub, or skip publishing |
| `git push` rejected (remote moved) | `git pull --rebase`, then retry `git push` |
| No PR summary on the PR | Confirm the branch was pushed with `git push` while hooks are installed; run `tl doctor` |
| Repo not initialized errors | `tl init -y`, then retry |
| Something is broken and the user wants support to see it | `tl doctor --report` sends the diagnosis plus recent tl logs |
