# Session capture: binding sessions to repos, and indexing without matching

## How capture works today, in five sentences

Every coding tool writes its own record of a session to disk as it goes — Claude
and Codex append JSONL files under their config directories, Cursor writes rows
into a SQLite database. When you `git push`, GX's `pre-push` hook runs `gx
capture push`, which finds the session files overlapping the pushed commits,
parses them into a normalized `SessionEvent` stream, and redacts secrets. It
then matches those events against the commits by **comparing the text the agent
wrote against the lines the commit added** (exact n-gram overlap first, fuzzy
token similarity second), with timestamp proximity as a weaker third signal.
Whatever matches becomes a hunk link — this hunk came from this session — which
is what produces coverage, authorship, and the PR summary. Nothing depends on GX
watching the session happen; it is entirely reconstructed afterward from the
file the tool left behind.

## Two problems, one plan

**Freshness.** The tool has often not finished writing the session when we look.
A push on 2026-07-28 reported `sessions=0 coverage=0.0%`; the identical ref range
replayed minutes later scored `50.0%`. Truncating that transcript reproduces it:
0.0% at 70% of the file, 33.3% at 80%, 50.0% at 90%. Because matching is
content-based, a missing tail is not a degraded match — it is no match.

**Attribution is doing double duty.** Content matching is currently the only
thing keeping other repositories' sessions out of a repo's context. Discovery
casts a wide net — any transcript mentioning the repo's path is a candidate —
and `capture_sessions` today holds sessions sourced from `-Users-joe-git-console`
and `-Users-joe-git-yeet` that were pulled in while capturing for `gx`. Matching
throws them back. If matching stops being a prerequisite for indexing, something
else has to be the gate.

The plan fixes both by binding a session to a **repository** rather than to a
diff, and gating on that binding.

## The binding chain

    session → working directory → git origin → connected repo

Each hop is verifiable and each tool supports it today:

| Tool   | Directory signal                                          | Status                      |
| ------ | --------------------------------------------------------- | --------------------------- |
| Claude | `cwd` on entries (1069 occurrences in a sampled session), plus `gitBranch` | present |
| Codex  | `session_meta.cwd` **and** `git.repository_url` outright   | present; origin needs no inference |
| Cursor | `workspace.json` → `folder`                                | already parsed by the existing Cursor parser |

Claude records `cwd` per entry, not once per session, so a session that moves
between repositories can be attributed per chunk rather than wholesale.

Resolution rules:

1. **Directory → repo root** must be worktree-aware. `internal/capture/repopath`
   already resolves a path against every checkout of a repository; a session run
   in `.claude/worktrees/<name>` belongs to the parent repository.
2. **Repo root → origin** reads `remote.origin.url`, normalized across `ssh` and
   `https` spellings, optional `.git` suffix, and case. Codex supplies the URL
   directly and can skip the inference.
3. **Origin → connected repo** matches the normalized origin against the repos
   connected to the user's GX org. This is the gate.

## What changes

**Indexing no longer requires matching.** A session bound to a connected repo is
pushed as semantic context even when zero hunks matched. Matching remains, and
remains valuable — it is what produces hunk links, authorship, and coverage — but
it becomes a precision layer rather than an admission requirement.

This makes the freshness problem degrade gracefully. A half-written transcript
today yields nothing at all. Under this model it yields "this session worked on
this repo, here is the context," which is worse than a precise hunk link and far
better than silence.

**Cross-repo context is allowed, deliberately.** Sessions bound to any connected
repo may inform any other connected repo in the same org. Related repositories —
`gx` and `console` — genuinely benefit from each other's context. The boundary is
org connection, not repo identity.

**Unconnected repos never leave the machine.** A session bound to a repo with no
connected origin, or to no repo at all, is not indexed and not uploaded. `yeet`
is not connected, so `yeet` sessions stay local. This is the property that makes
the previous point safe.

## The privacy boundary, and what we say about it

The gate is the git origin of the repository a session was working in. A
personal side project has a different origin, or none, so its sessions are never
eligible for a company index.

The residual case is a developer who works on a personal project **inside a
session that is also doing company work**. Chunk-level attribution via per-entry
`cwd` limits this — chunks are attributed to the repo that was current — but a
session deliberately mixing the two will have its company-repo portions indexed.

That is a documented user responsibility, not a mechanism we try to defeat. It
must be stated plainly in user-facing docs:

> GX indexes agent sessions against the repositories your organization has
> connected. A session is bound to a repository by the directory it ran in and
> that repository's git origin. Work done in a repository your organization has
> not connected is never uploaded. If you work on a personal project inside the
> same session as company work, the portions attributed to the connected
> repository will be indexed — keep personal work in a separate session.

This wording is part of the deliverable, not a follow-up.

## Phases

**Phase 1 — Bind sessions to repos.**
Extract the working directory per session (and per chunk where the tool gives it
per entry). Resolve directory → repo root → normalized origin. Persist the
binding on `capture_sessions`. No behavior change yet; this is observability.

**Phase 2 — Gate on connected repos.**
Match the resolved origin against the org's connected repos. Sessions with no
connected origin become ineligible for upload. This narrows what is shared and
should land *before* Phase 3 widens it.

**Phase 3 — Index without matching.**
Push bound, connected sessions as semantic context regardless of hunk links.
Keep matching for hunk links, authorship, and coverage. Normalize and redact on
the same path as today — nothing reaches Postgres unnormalized.

**Phase 4 — Freshness.**
Persist ranges that attributed nothing (`Result.AttributionSuspect` already
detects the shape) along with the source watermarks observed
(`source_bytes`/`content_hash`, already on the table). Re-attempt when a
watermark advances, triggered from invocations that already occur. Lower urgency
after Phase 3, since context no longer depends on matching, but hunk links still
do.

**Phase 5 — Documentation.**
Ship the boundary statement above, plus how bindings are determined and how a
user checks what is connected.

## Risks and open questions

- **Origin normalization is the weak link.** `git@github.com:org/repo.git`,
  `https://github.com/org/repo`, and a fork's origin must resolve consistently, or
  the gate silently fails open or closed. This deserves its own tests, including
  forks and repos with multiple remotes.
- **No origin at all.** Local-only repos have no remote. Treat as unconnected
  (never uploaded) rather than falling back to path heuristics.
- **Re-staging must stay idempotent.** Phases 3 and 4 both re-touch already
  staged rows; the existing `ON CONFLICT ... DO UPDATE` keys must be reused
  deliberately.
- **Redaction on every path.** A re-read or an unmatched-session push must go
  through the same redaction as matched sessions. "It was redacted once" is how a
  secret ships.
- **Do not depend on tool hooks.** Claude offers `PostToolUse`, which would
  remove the lag entirely; Codex, Cursor, and whatever ships next may offer
  nothing comparable. Hooks are an opportunistic accelerator for tools that have
  them, never the mechanism.
- **Do not go fully file-as-truth.** Tools rotate and delete their own
  transcripts and worktrees get removed. Keep normalized events; continue not
  storing raw blobs (`raw_blob` is null in all 22 stored sessions today).
