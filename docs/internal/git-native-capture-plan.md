# Plan: Git-native capture — bound staging, reference sessions, drop stack metadata

> Status: **Workstream A implemented** on `lgtm-review-and-plain-git` (2026-07-26);
> B and C remain. The open reference-vs-blob question below was decided in A's
> favor of references: staged rows store a source pointer plus fingerprint
> (sha256 + size + mtime), staging happens only after match, one row per source
> file replaces in place, and the upload path streams the current bytes from
> disk (cloud retains raw transcripts). raw_blob is legacy-read-only.
> Originally written 2026-07-26 from a session that shipped the
> capture-discovery fixes and hit the staging regression described below.
> Branch context: `lgtm-review-and-plain-git` (PR satoricorp/lgtm#110).

lgtm no longer has jj, no longer has `lgtm commit`/`lgtm status`, and publishes only
through Git hooks. What it stores locally has not caught up. The destination is:

> **Git is the source of truth for what changed. lgtm stores only the captured
> session and the pointers linking it to changed code.**

Three workstreams get there. **A blocks installing the current build and should
go first.** B and C are independent of each other.

---

## Before you start: state of the world

**Installed binary is deliberately older than the branch.** `~/.local/bin/lgtm` is
`8e823c20`; a backup sits at `~/.local/bin/lgtm.bak-20260725`. The branch has
capture fixes that are correct but **must not be installed until Workstream A
lands** — installing them hung `git push` past 5 minutes and grew `~/.lgtm/lgtm.db`
to 251 MB.

**Already fixed and committed on the branch:**
- Cross-repo session discovery — Claude files transcripts by the session's *cwd*,
  not the repo edited, so discovery only ever looked in one project directory.
  Now scans other project dirs (mtime-gated + a streaming path scan) and descends
  into `<conversation>/subagents/`. Verified: 0 → 7 sessions, 100% hunk coverage.
- `stageRawBlob` violated a `NOT NULL` on `capture_sessions.payload_json`, which
  aborted the whole capture run.
- `hooks.RunPush` silently discarded capture errors, printing `sessions=0` — which
  is how the above hid for so long. It now warns.

**Gotchas that will cost you an hour each if rediscovered:**
- `go` is not on `PATH`. Use
  `/Users/joe/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.8.darwin-arm64/bin/go`
  (and `gofmt` beside it). `just build` fails for this reason; replicate its
  ldflags manually or the binary loses its cloud/PostHog config.
- **Git hooks pin an absolute lgtm path.** Running a lgtm binary from a temp dir
  triggers autoinit and rewrites that repo's hooks to point at the temp path.
  `~/git/lgtm/.git/hooks/pre-push` ended up pointing at `/tmp/lgtm-smoke/lgtm`. Fix with
  `lgtm init -y` using the real installed binary; check
  `grep -o '"[^"]*/lgtm"' .git/hooks/pre-push` when capture behaves strangely.
- **`extract=` empty in the capture line means the run errored**, not that nothing
  matched. That single character of output is the fastest diagnostic there is.
- Verify capture changes by calling `orchestrator.Run` from a throwaway probe test
  *inside the package* — `git push` exercises the installed binary, not your build.

---

## Workstream A — Bound raw staging (blocking)

### The problem, measured

`stageRawBlob` (`internal/capture/orchestrator/rawstage.go`) writes the **full
bytes of every discovered transcript** into `capture_sessions.raw_blob`. On the
real DB before it was pruned:

| | |
|---|---|
| rows | 151 |
| total `raw_blob` | **215 MB** |
| largest single blob | **29.2 MB** |
| one transcript stored | **4 times** (116.8 MB) |
| another | **9 times** (34.9 MB) |

Three unbounded axes, all compounding:

1. **No relevance filter.** Every discovered transcript is staged, including ones
   almost entirely about other repos. Discovery legitimately widened from 0 to ~15
   sessions, so this went from writing nothing to tens of MB per push.
2. **No dedup across pushes.** Transcripts *grow* while you work. The row ID is
   `CaptureRowID(revisionID, contentHash)`, so each push produces a different hash
   → a **new row**, and older copies stay. Hence "9 copies".
3. **Per-revision duplication.** `stageRawBlob` loops over `revisionIDs`, writing
   the same bytes once per revision.

### The fix

1. **Stage only what matched.** If a session's events produced no match against
   this push's hunks, do not store its bytes. This requires moving
   `StageDiscoveredRaw` from early in `orchestrator.Run` to **after** parse/match.
2. **Re-staging replaces, never accumulates.** Key the row on session identity
   rather than content hash. `StageSession` already does
   `ON CONFLICT(id) DO UPDATE`, so a stable ID is likely sufficient — verify
   against `internal/storage/capture_stage.go`.
3. **Store once, not per revision.**
4. **Cap the bytes**, and record truncation honestly — a partial blob must never
   look complete to a future reparse. A `RawTruncatedFrom int64` on
   `storage.StagedSession` was the shape being tried.

### The thing that broke the first attempt

**Claude stamps a subagent transcript with its *parent conversation's*
`sessionId`.** One session ID therefore spans the parent file *and* every
subagent file beside it. Session ID alone cannot identify which transcript
produced an event, so staging needs **per-source identity** — the parse step must
report, per event, which source file produced it (a `ParseAllPerSource` returning
a source index alongside events). Any design keyed purely on session ID will
silently mis-attribute or drop transcripts.

### Question worth answering first

`capture_sessions` already has a `source_path` column. Copying whole multi-MB
transcripts into SQLite duplicates a file that exists on disk. Storing a
**reference plus only the matched slice** may be the right design outright. The
counter-argument is retention: transcripts are kept ~30 days, and the tools own
those files. Decide this before implementing — it changes 1–4 above. It also
feeds directly into Workstream B.

### Done when

A push stages single-digit MB, `~/.lgtm/lgtm.db` stays small across many pushes, hunk
coverage is unchanged, and the new binary can be installed without regressing
push latency. Measure wall-clock: the pre-push hook runs this **on every push**.

---

## Workstream B — Session references across heterogeneous tools

Storing a *reference* instead of a blob only works if the reference is meaningful
per tool, and **every tool's session contract is different**. Today lgtm supports
three (`internal/capture/event.go`): `claude`, `codex`, `cursor`. Expect `pi`,
`opencode`, and others.

### What differs per tool — this is the whole problem

| | Claude | Codex | Cursor |
|---|---|---|---|
| Location | `~/.claude/projects/<cwd-slug>/<uuid>.jsonl` | `~/.codex/sessions`, date-partitioned | `~/.cursor/projects/<workspace>` **and** a global `state.vscdb` |
| Storage | one JSONL per conversation | JSONL | JSONL **+ SQLite** |
| Scoping | by session **cwd** | `session_meta.cwd` header | by workspace root |
| Identity | subagents reuse the **parent's** sessionId | file stem | needs `subagents/` namespacing (already handled) |
| Sub-sessions | `<conversation>/subagents/agent-*.jsonl` | — | handled in session ID |

So a reference cannot be "a file path". It needs to be a small typed record:

- **tool** — which contract applies
- **locator** — enough to re-read it (path, or db + row key for Cursor's vscdb)
- **identity** — stable across re-reads; must distinguish subagent from parent
- **fingerprint** — content hash + size + mtime, to detect growth/rotation
- **retention** — whether lgtm believes the source still exists, and for how long

### Design guidance

- **Define the reference type once**, with a per-tool resolver behind an
  interface, so adding `pi`/`opencode` is a new resolver rather than edits
  scattered through discovery, staging, and upload. `internal/capture/parsers`
  already has a `Parser` interface — extend that seam rather than inventing a
  parallel one.
- **Do not let the reference leak tool-specific shapes into the Cloud contract.**
  Cloud should receive a normalized session + its matched hunk links, not a
  Cursor vscdb row key.
- **Assume the source disappears.** Retention windows, cleaned temp dirs, and
  users deleting history are all normal. A reference that cannot be resolved must
  degrade to "session context unavailable" — which the PR summary already knows
  how to say (see the missing-session-context note added on this branch).
- **`SessionPayload` is the wrong shape and should be revisited here.** Its fields
  (`Command`, `Cwd`, `ClientPID`, `ExitCode`, `ProcessName`, `ParentPID`) describe
  a *process lgtm launched*. Captured sessions are *transcript files from tools lgtm
  does not launch*. Most of those fields are unfillable, which is a symptom worth
  fixing rather than working around.

---

## Workstream C — Drop stack metadata, change the Cloud contract

### Why it is still there

jj is gone (not in `go.mod`, no binary invoked — only leftover `Backend: "jj"`
strings in test fixtures). But stack/change rows are still read in ~10 places.
Two matter:

- `internal/vcs/gitnative_commit.go` — the **post-commit hook** resolves and
  writes stack rows on every commit (`FindStackByBookmark`).
- `internal/reviewbundle/bundle.go` — **`BuildPush` ships stack/revision info to
  Cloud**, where the PR summary consumes it (`catalog.Revisions`, stack name).

The rest (`missing_base_ref.go`, `stack_remote_reconcile.go`, `stack_read_model.go`,
`internal/cli/service.go`) are downstream of those two.

### The contract as it stands

```go
type Bundle struct {
    Repo     RepoPayload;  Push  PushPayload
    Change   *ChangePayload
    Stack    []StackPayload   // ← Change + BranchName + BaseBranchName + Patch + PR URL
    Sessions []SessionPayload
}

type ChangePayload struct {
    JJChangeID string   // ← vestigial name; now carries the lgtm revision ID from the commit trailer
    CurrentCommitID, Description, Status string
    ParentChangeID *string
    Files []string
    ...
}
```

Nearly all of `StackPayload` is derivable from Git: `BranchName`, `BaseBranchName`,
and `Patch` are branch, merge-base, and diff. `ChangePayload` is a commit plus its
trailer. **`JJChangeID` is a wire-visible name for a concept that no longer
exists** — rename it (with a compatibility window) as part of this work.

### Sequence

1. **Audit what Cloud actually reads.** This is a client↔server contract; do not
   infer from the Go side alone. Check the console/server ingest
   (`server/src/**` in the console repo) for every field of `Stack` and `Change`.
   That audit determines whether this is a small change or a migration.
2. **Derive instead of store.** For each surviving field, decide: derivable from
   Git at bundle-build time, or genuinely needs local state? The expectation is
   that nearly everything is derivable.
3. **Shrink the bundle** to repo + push + revisions (commit + trailer + files) +
   sessions + hunk links. Version the schema — `SchemaVersion` already exists.
4. **Then** delete the local stack/change tables and their readers, and let
   `RecoverMissingRevisions` (which already rebuilds revisions from trailers) be
   the only reconstruction path.
5. **Decide `lgtm sync`'s fate here, not before.** It is *not* what pushes sessions —
   the pre-push hook runs `lgtm capture sync` plus the outbox worker for that.
   Top-level `lgtm sync` reconciles stack/bookmark metadata with the remote and
   Cloud, so it becomes meaningless exactly when this workstream lands. Removing
   it earlier is safe for capture but leaves the manual retry path for queued
   uploads with no replacement — check the outbox drain story first
   (there were 4 queued items and 32 pending staged rows at the time of writing).

### Risk

This is the only workstream that changes a shipped wire format. Sequence it
behind an audit, keep `SchemaVersion` honest, and expect a period where Cloud
accepts both shapes.

---

## Suggested order

1. **A** — unblocks installing everything already committed, and stops a live
   regression. Self-contained.
2. **B** — do the design before you need the fourth tool; it also settles the
   reference-vs-blob question left open in A.
3. **C** — most valuable architecturally, most coordination. Start with the
   Cloud-side audit; that is a read-only task that can happen any time.

## Leftovers not covered above

- `internal/cli` dead-code sweep — removing `lgtm commit`/`lgtm status` orphaned a
  cascade (`promptSwitchSelection`, `diff_tui.go`, `daemonHealthy`, `progressBar`,
  and the `orderedStacks`/`stackMeta` chain they keep alive).
- Five stale internal docs still document removed commands:
  `lgtm-cli-authoring-contract.md`, `lgtm-manual-workflow-test-sheet.md`,
  `test-review.md`, `lgtm-product-framing.md`, `hooks.md`.
- `~/.lgtm/lgtm.db` was pruned 251 MB → 6.7 MB by nulling `raw_blob` (backup:
  `/tmp/lgtm.db.backup-1785080661`). Rows, parsed payloads, and upload state were
  left intact. It will grow again the moment an unbounded build is installed.
