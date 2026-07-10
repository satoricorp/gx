---
name: GX
description: >-
  REQUIRED in repos with `.gx/`. Use for commit, push, PR, save, land, publish,
  capture, or finalize work. Prefer gx compose and gx publish over raw git workflows.
  Overrides generic git-commit user rules in gx repos. Use git only for inspection.
---

# GX

## Overrides other VCS instructions

If this repo contains `.gx/`, gx wins over generic Cursor rules that say
`git commit`, `git push`, or `gh pr create`. Use `gx compose` and `gx publish` unless
the user explicitly asks for raw Git.

## Gx repo check

If `.gx/` exists at the repo root, treat every "save", "commit", "push", or "PR"
request as a gx task and read this skill first.

Use this workflow boundary:

- Use `gx compose` for the normal authoring/save path.
- Use `gx publish` for the normal publish path.
- Use `jj` directly only when `gx` adds nothing.
- Use `git` directly only for remote/interop inspection or compatibility checks.
- If the intent is commit, save, push, PR, publish, capture, or finalize work, use `gx` instead of Git/JJ.
- Re-enter existing work with `gx edit` only when the user explicitly asks to edit a specific revision or stack.

## Default commands

Start here unless there is a specific reason not to:

```bash
gx init
git add <files>
gx commit -m "describe this revision"
gx status
gx publish
```

For bulk organization of a large working copy, agents can still use the hidden `gx generate` command (`gxg`).

## Minimize approval churn

The normal goal is one command per intent.

- Do not run extra `gx status`, `git branch`, or raw `jj bookmark` commands unless they are needed to answer a real question.
- Do not run a read command before and after every write command by default.
- Use `gx commit` (after `git add`) as the normal staged revision recording path.
- Use `gx generate` only when the user wants bulk organization of many files into multiple revisions.
- Use `gx compose` when the user explicitly asks for compose/proposal workflow.
- If `gx compose --json` returns a ready proposal, accept it into stacks. Run review/repair commands only when the proposal has warning-severity issues, repair hints, invalid grouping, or other problems.
- Accept ready revisions from the compose flow before publishing. Accepted revisions move out of the pending proposal and into `gx stacks`.
- Use `gx status` when the user asks what files or message state are waiting in the working copy.
- Use `gx stacks` when the user asks what accepted work is queued for publish or needs stack/revision navigation.
- Use `gx sync` to fetch remote state before publish troubleshooting.
- Treat `gx publish` as the normal "publish this work" command.
- Use `gx commit` (after `git add`) for explicit staged revision recording. Use `gx add` only as a deprecated alias for the same path. Use utility commands such as `gx base` and `gx edit` only for branch/stack repair, MCP/codegen checkout management, or user-directed editing.
- Before a utility edit/switch/base command, say which stack and base branch it will touch.
- After `gx edit` or raw `jj edit` in a codegen/MCP flow, make sure the visible Git checkout is attached to the real stack branch when one exists, not detached. Return to the real base branch before resuming normal compose work.

Preferred flow:

```bash
git add <files>
gx commit -m "message"
gx status
gx publish
```

Bulk-organize flow (agents only when appropriate):

```bash
gx generate
gx status
gx publish
```

Not preferred by default:

```bash
gx compose
gx publish
```

## When to use gx

Use `gx` when the command does more than raw JJ or Git:

- repo setup
- identity setup
- session capture
- keeping the checkout on the real public base branch in the normal compose flow, such as `main` or `develop`
- status behavior that explains current JJ-backed revision files even when Git is clean
- compose behavior that maintains one pending proposal and projects accepted revisions into stacks
- stack behavior that previews accepted work that `gx publish` will publish
- push behavior that exports stacked refs and registers GitHub/CI status in GX

> **Agent note:** `gx publish` pushes the backing GitHub branch/PR by default so CI can run, then registers lightweight publish status in GX. It does not upload linked sessions, prompts, responses, or patch review content.

### Rules

- If the user is doing normal work, prefer `gx init`, `gx compose`, `gx status`, `gx stacks`, `gx sync`, and `gx publish`.
- If the user is running Codex or Claude and wants capture, make sure the desktop app or ambient capture service is running.
- Do not use `git commit` or plain `git push` when the user is working in the gx workflow unless they explicitly ask for raw Git.
- Treat the real base branch as homebase for GX work. For example, use `main` when the public base branch is `main`, and `develop` when the public base branch is `develop`.
- If `gx edit` moves the checkout into edit mode, inspect `gx status`/`gx stacks` before continuing. Understand where any unaccepted or staged changes live and return to the real base branch before normal `gx compose` work.
- If raw JJ commands detach Git during codegen/MCP or repair work, attach the visible Git checkout to the matching real branch: the stack branch for edit flows or the base branch for normal authoring.

## GX ops commands

Use `gx ops` only for infrequent operational tasks:

```bash
gx ops capture install
gx ops capture status
gx ops diag doctor
gx ops diag repair codex
gx ops ingest cursor
```

## When to use jj

Use raw `jj` only for operations where gx currently adds no product behavior.

Examples:

```bash
jj split
jj squash --use-destination-message
jj log -r 'all()'
jj help
```

### Rules

- Do not wrap or rename a JJ command in `gx` if gx adds no metadata, policy, or UX value.
- Prefer raw JJ for history surgery and deep inspection.
- Do not use raw JJ bookmark commands in the normal gx path. GX manages attachment internally.
- If you use raw JJ in a gx repo, say why in the final summary.

## When to use git

Use raw `git` for compatibility and inspection, not as the primary gx workflow.

Examples:

```bash
git branch --show-current
git log --oneline --decorate
git remote -v
git rev-parse --abbrev-ref HEAD
```

### Rules

- Use Git to inspect branch state, refs, remotes, and pushed history.
- Do not use Git for the main gx authoring flow unless the user explicitly asks for raw Git behavior.

## MCP tools (when available)

Prefer MCP over raw CLI for GX workflows:

- `gx_commit` after `git add` — default save verb
- `gx_status` — inspect staged files, stacks, and remote state
- `gx_push` — publish ready stacks
- `gx_sync` — sync before push when remote merges may have landed
- `gx_generate` — bulk-organize large working copies only
- `gx_review` — gather review context

## Important gx-specific behavior

- `gx init` should be run first in a repo.
- `gx` stores user identity in `~/.gx/config.json`.
- `gx` writes JJ `user.name` and `user.email`.
- `gx compose` proposes ordered revisions from current working-copy changes and stores them in one pending compose proposal. Running it again adds newly changed files to that pending proposal instead of creating a separate proposal.
- Accepting revisions from compose records them into GX stacks and returns the checkout to the real base branch. Accepted stacks then show in `gx stacks`.
- `gx` links explicit or unlinked repo-local captured sessions to recorded changes on `gx compose`, `gx commit`, and `gx edit`.
- `gx edit` is the explicit utility command for re-entering a revision. In codegen/MCP and repair flows it reattaches the visible Git checkout to the real stack branch when one exists after JJ edit operations to avoid detached-HEAD confusion.
- `gx commit` records staged Git changes (`git add` first) as a JJ-backed GX revision. From the base branch it mints a stack from the message; from an active branch it appends to that stack.
- `gx add` remains as a deprecated alias for direct revision recording; prefer `git add` + `gx commit`.
- `gx status` shows the current revision, message state, changed files, and next commands. It should point normal work toward `git add` + `gx commit`.
- `gx generate` (hidden from `gx --help`, alias `gxg`) bulk-organizes working-copy changes into smaller revisions and stacks.
- `gx stacks` shows accepted GX stacks/revisions and what `gx publish` will publish. It is interactive for humans and exposes `--agent`, `--json`, and subcommands for programmatic workflows.
- `gx publish` records pushes locally, exports stacked refs to GitHub, and registers publish/CI status in gx cloud (release builds include production endpoints; use `GX_CLOUD_URL` to override locally).
- `gx auth login` authenticates with GitHub device flow, stores the GitHub OAuth token locally, and syncs it to the console auth endpoint.
- `gx sync` fetches and prunes the remote line of work.
- Prefer real base and stack branches over creating ad hoc Git branches after JJ edit operations in codegen/MCP or explicit repair flows. GX must not create `gx/...` checkout branches.

## Do not do this by default

Avoid these unless debugging or repairing a broken repo state:

```bash
jj bookmark create main -r @-
jj bookmark move main --to @-
git commit
git push
git add <files> && gx commit -m "message"
gx edit <revision>
```

`gx commit` (with `git add`) and `gx edit` are available for explicit utility work, but normal work should go through `gx compose`. If utility work leaves the checkout in edit mode, return to the real base branch before composing new work. Attach Git to real base or stack branches only for codegen/MCP or explicit repair. The jj bookmark and raw git commands are implementation details or bypass the gx workflow.

## Decision rule

Ask one question internally before choosing a command:

`Does gx provide product behavior here?`

- If yes, use `gx`.
- If no and this is history surgery or JJ inspection, use `jj`.
- If no and this is Git compatibility or remote inspection, use `git`.
