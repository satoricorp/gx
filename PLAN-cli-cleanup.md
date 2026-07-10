# Plan: CLI cleanup + new-user demo

Repo: `~/git/gx` (Go, cobra CLI in `internal/cli/`, jj-colocated git repo).
Branch to work on: `feature/merge-generate-surface-improvements-from-feature` (or a new branch off it).

## Thesis (context for every decision below)

`gx commit` is the product. The workflow is git-native:

> stage with `git add` / `git add -p` → `gx commit -m` records a provenance-rich,
> JJ-backed revision (sessions, models, tokens captured ambiently) → review → `git push`.

gx does not replace git's staging or push ergonomics. Anything in the CLI that
predates this — auto-composing the working copy, parallel checkpoint commands,
legacy auth — gets deleted, hidden, or rewritten to serve `gx commit`.

Command inventory ground truth: `internal/cli/root.go` (`root.AddCommand(...)` at ~line 75)
plus `service.go`, `capture.go`, `auth.go`, `login.go`, `set.go`, `demo.go`.

---

## Phase 1 — Delete dead code (zero behavior change)

These commands are **defined but never registered** — no `AddCommand` references
them anywhere. Verify each with grep before deleting; then delete the constructor
and everything that becomes unreferenced.

All in `internal/cli/root.go` unless noted:

| Function | What it was |
|---|---|
| `newComposeCommand` | unregistered wrapper that renames demux → compose |
| `newDemuxCommand` | only constructed by `newComposeCommand` → dead with it |
| `newDemuxListCommand`, `newDemuxReviewCommand`, `newDemuxFixCommand`, `newDemuxShowCommand`, `newDemuxApplyCommand`, `newDemuxApplyPlanCommand`, `newDemuxReviewPlanCommand` | demux subtree, only attached inside `newDemuxCommand` |
| `newComposeDoctorCommand` | second "doctor" ("Check jj provisioning…"), only reachable via demux |
| `configureComposeApplyDefaults`, `renameDemuxCommandSurface` | compose-only helpers |
| `newStackCommand` (`stack --new`) | unregistered |
| `newStacksCommand` (`stacks` browse TUI) + its `list` subcommand | unregistered |

**Careful — shared code that must stay:**
- The authoring/demux **engine** (`internal/authoring/`): used by `gx generate` and `gx commit`. Do not touch.
- Helpers used by `gx generate`: `runDemuxWithLoader`, `runDemuxApplyWithLoader`,
  `printDemuxChangesPacket`, `printDemuxApplySummary`, `demuxAutoAcceptBlockedReason`,
  `useDemuxLoader`, `useStatusInteractive`. Delete only what `go build` + grep prove unreferenced.
- `stacks_tui.go`: `useStatusInteractive` and `runStacksInteractive` are used by
  live commands; only delete pieces that die with `newStacksCommand`.

Method: delete top-level constructors first, then iterate `go build ./...` and
remove newly-unused helpers/types until clean. Run `go test ./internal/cli/ ./internal/authoring/`.
Delete any tests that only exercised the removed commands.

## Phase 2 — Retire reachable legacy commands

1. **`gx login`** (`internal/cli/login.go`, hidden): legacy token-paste upload auth.
   Superseded by `gx auth login`. Delete the command and file; keep the
   `uploadauth` package (used elsewhere). Remove `newLoginCommand` from `root.AddCommand`.
2. **`gx add`** (`newAddCommand`, hidden): old "checkpoint the working copy" flow.
   Competes with `git add` + `gx commit`. Delete the command, `printAddSummary`,
   `addResultJSON`. Keep the engine's checkpoint APIs (used by commit).
3. **`gxg` binary alias**: in `Execute()` (root.go), the `gxg` basename maps to
   `generate`. Remove the mapping and the `gxg` alias on the generate command.
   Keep `gxr` (review) and `gxs` (status).
4. **`gx generate`** itself: keep, keep hidden. Its fate is an open product
   question (see Backlog) — do not delete in this pass.
5. **`gx demo`**: do not delete — replaced in Phase 3.

Update `docs/` and shell completions if they reference removed commands
(`grep -r "gx add\|gx login\|gx compose\|gx stacks\|gxg" docs/ internal/`).

## Phase 3 — New `gx demo`: interactive first-run walkthrough

Rewrite `internal/cli/demo.go`. Unhide it (`Hidden: false`, group `Setup`).

**Mechanics**
- Five steps, advanced with Enter. `q` quits. No flags needed for the default run.
- Runs against a **scratch repo** in a temp dir (`gx demo` must never touch the
  user's real repo). Create it, seed one small file, and actually execute
  `gx init`, `git add`, `gx commit` there so output is real. Clean up on exit.
- Step 2 (`gx auth login`): if already signed in, show `already signed in ✓` and
  skip the device flow. Never force a real login inside the demo — show the
  command and what it does, run it only if the user confirms.
- Step 5 (`git push`): the scratch repo has no remote — show the copy and a
  simulated result line rather than failing.
- Style: use the existing `termstyle`/`labelValue`/`muted`/`success` helpers.
  One screen per step. No paragraphs longer than three lines. Generous whitespace.

**Copy (verbatim — this is the spec):**

```
$ gx demo

  gx — version control that captures how code was made

  5 steps · 2 minutes · runs in a scratch repo, your work is untouched


  ── 1/5 · init ─────────────────────────────────────────

  Set up gx in a repository.

    $ gx init

  Installs hooks and starts ambient capture. From here on,
  every AI session that touches this repo is recorded —
  models, tokens, and context ride along with your commits.

  ↵ run it


  ── 2/5 · sign in ──────────────────────────────────────

  Connect to gx cloud.

    $ gx auth login

  GitHub device login. Captures stay on your machine
  until you choose to share them.

  ↵ run it


  ── 3/5 · stage ────────────────────────────────────────

  Choose what goes in — plain git.

    $ git add -p

  gx doesn't replace staging. Pick files and hunks
  exactly like you always have.

  ↵ run it


  ── 4/5 · commit ───────────────────────────────────────

  Record it as a gx revision.

    $ gx commit -m "add rate limiter"

  Looks like a commit. Records more: the sessions, models,
  and tokens behind the change — provenance your reviewer
  can actually use.

  ↵ run it


  ── 5/5 · push ─────────────────────────────────────────

  Ship it.

    $ git push

  Your branch goes up like always. gx attaches the context,
  so reviewers see the why — not just the diff.

  ↵ run it


  ── done ───────────────────────────────────────────────

  That's the loop: stage with git, commit with gx.

    gx status    your stacks at a glance
    gx review    context-aware code review
    gx doctor    check and fix your setup

```

Copy rules if any of it must be adapted: verbs first, no exclamation marks, no
marketing adjectives, nothing the user can't verify on screen. Each step is
title → one-line purpose → the command → ≤3 lines of payoff.

## Phase 4 — Verification

- `go build ./...` and `go vet ./...` clean.
- `go test ./internal/cli/ ./internal/authoring/ ./internal/vcs/` green.
- `gx help` shows: Setup (init, auth, set, demo), Work (commit, review, status),
  Ship (push, sync), Help (doctor, report, version) — nothing else.
- `gx compose`, `gx stacks`, `gx stack`, `gx login`, `gx add` all print
  "unknown command".
- Run `gx demo` end to end in a terminal; every screen matches the copy above.
- Commit in logical phases (1 / 2 / 3), trailer:
  `Co-Authored-By: <model attribution>`.

## Backlog — review for thesis alignment (document, don't build)

Not in scope for this pass; listed so they aren't lost:

- **`gx edit`**: `gx commit --amend` errors with "use `gx edit <rev>`" but `edit`
  is hidden. Bless it (unhide + polish) or give commit a real `--amend`.
- **`gx status`**: its `edit`/`diff` subcommands overlap top-level `edit` — one story.
- **`gx review`**: confirm it consumes commit provenance (sessions, models, tokens
  now captured per-response in `~/.gx/gx.db`).
- **`gx push`/`gx sync`**: confirm they handle commit-created branches
  (`CreatedBranch`) as first-class, not just generate-style stacks.
- **`gx init`**: should offer capture-service install + codex proxy routing
  (`gx ops capture install`, doctor's codex check) — today a fresh install
  captures nothing until the user finds `gx ops`.
- **`gx ops`**: flatten. Service management is load-bearing; `gx ops capture start`
  is three levels deep. Candidate: promote to `gx service` or fold into `doctor --fix`.
- **`gx generate`**: decide fate — retool as an explicitly-invoked escape hatch
  ("split this working copy for me") or retire once commit covers its users.
