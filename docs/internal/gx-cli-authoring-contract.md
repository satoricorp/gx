# GX CLI Authoring Contract

> **Note (2026):** `gx compose`, `gx stacks`, and `gx add` were removed. Normal
> work uses `git add` + `gx commit`, `gx status`, and the hidden `gx generate`.

This document captures the intended behavior for GX authoring, compose, stacks,
and publish workflows.

## Mental Model

GX is a product layer over JJ and Git.

The user-facing objects are:

1. **Session**: captured context from Codex, Claude, Cursor, shell commands,
   files, prompts, responses, and tests.
2. **Revision**: one reviewable logical change. Locally this is a JJ change,
   identified by JJ `change_id`.
3. **Stack**: an ordered line of revisions. Locally this is GX metadata plus a
   JJ/Git-compatible bookmark or branch ref.
4. **Published stack**: the remote GitHub/review representation of that local
   stack.

The important mapping is:

```text
GX revision = JJ change = stable product identity
Git commit = current exported snapshot of that revision
GX stack = ordered container of revisions
Git branch/bookmark = transport/ref for that stack
```

JJ gives GX stable change IDs while commits can be rewritten. Git exists so the
outside world can receive branches and PRs.

## Normal Flow

The normal workflow should be:

```bash
git add <files>
gx commit -m "describe this revision"
gx status
git push
gh pr create
```

`git add` + `gx commit` records staged work as a GX revision. The hidden
`gx generate` command can organize a large working copy into smaller revisions
when needed.

`gx status` shows local features, revisions, and remote state.

Publish with plain `git push` (the GX pre-push hook captures the session and
publishes). Open the PR with `gh pr create`. Do not run `gx push` or
`gx publish` — those paths are disabled or removed; plain `git push` is the
user-facing publish path.

The user should stay on a normal branch checkout after commit and push.
Detached HEAD is not an acceptable steady state.

## Compose

`gx compose` is the core authoring command. It should take messy working-copy
changes and turn them into reviewable revisions grouped into stacks.

The intended compose pipeline is:

1. Collect changed files and hunks.
2. Build a hunk catalog.
3. Infer structural facts: files, symbols, source/test relationships, and
   dependencies.
4. Propose revision boundaries.
5. Assign each revision to a target stack.
6. Review the plan for missing hunks, duplicate hunk ownership, incoherent
   order, invalid routes, and feasibility warnings.
7. Repair locally where deterministic repair is possible.
8. Use AI repair only when deterministic repair cannot produce a valid plan.
9. Preflight the actual apply in a disposable attempt.
10. If preflight fails, feed the concrete failure back into the repair loop.
11. Present only a proposal that should apply successfully.

Compose must not silently degrade to "I applied one valid piece; run me again."
If the user accepts all, GX should create all proposed stacks and revisions. If
the user accepts one stack, that stack should be created and visible
immediately, and the remaining unaccepted changes should stay available for the
next compose run.

A partial accept is valid. A hidden partial fallback is not.

The preflight message:

```text
Checking compose apply in disposable attempt 1/3...
```

means GX is testing the proposal before presenting or applying it. It should not
mutate the real worktree. If it fails, the proposal should be repaired before
the user sees it as ready.

## Stacks

`gx stacks` is the local truth view for accepted work.

It should show:

- stack name
- stack bookmark/branch
- base branch, usually `main`
- draft/published status
- revision count
- each revision's JJ change ID
- each revision's current Git commit ID
- published state
- remote ref if published

It should not require being on that stack. In the normal flow, the user is on
`main`, and `gx stacks` still shows accepted local stacks and published stacks
that are still under review. Merged or closed stacks are removed from the
working stack surface.

A stack record and its bookmark must agree. If GX metadata says a stack head is
change `A`, but the bookmark points at change `B`, that is a bug. Stack
metadata, JJ bookmark state, and Git branch refs must be kept in sync.

## Revisions And Change IDs

GX should store JJ change IDs for every GX revision.

That is the right design because JJ change IDs survive edits better than Git
commit hashes. Git commit hashes are still useful, but they are current
snapshots, not durable identities.

For each revision GX needs:

- JJ `change_id`
- current Git `commit_id`
- description
- stack membership
- stack order
- parent relationship when relevant
- files/hunks covered
- provenance/session links
- demux/compose evidence

Editing should work by change ID. If a user runs `gx edit <revision>`, GX should
find the stack containing that JJ change, enter the correct revision, and
preserve the relationship after the edit rewrites the commit.

## Branches And Bookmarks

GX should hide JJ's awkward checkout state.

Normal users should see real Git branches, not detached HEAD. Internally GX can
use JJ bookmarks and GX-owned refs, but the visible checkout should be attached
to a branch.

The intended visible-branch rules are:

- Normal compose work starts from `main`.
- Compose apply may temporarily move through JJ changes.
- After apply, visible Git checkout returns to `main`.
- Publish may create or update branch refs for stacks.
- After publish, visible Git checkout returns to `main`.
- Edit surgery should attach visible Git to the real stack branch when one
  exists. GX must not create `gx/...` checkout branches for edit/base state.

Branches/bookmarks created by compose should match the stack names shown in
compose. If compose proposes `feature/stack-management`, then accepting it
should create a stack whose visible name/ref in `gx stacks` is
`feature/stack-management`.

## Publish

Publish with plain `git push`. The GX pre-push hook captures the session and
registers publish metadata as the branch goes up. Open the PR separately with
`gh pr create`.

Expected behavior:

```bash
git push
gh pr create
```

means:

- push the stack branch via Git
- the pre-push hook records capture/publish metadata
- upload review context when cloud/GitHub integration is enabled
- open or update the PR with `gh pr create` (do not seed `## Summary`)

Do not run `gx push` or `gx publish`. Those are not the agent/user publish path.

## MCP And Agent Flow

MCP should call the same authoring engine behavior as the CLI.

The ideal agent loop is:

1. `git add` staged files
2. `gx_commit`
3. `gx_status`
4. plain `git push` when ready
5. `gh pr create`
6. `git commit --amend` when updating the latest revision, preserving its GX trailer
7. `gx_review` when review context is needed

There is no `gx_publish` / `gx_push` MCP tool. Agents must publish with plain
`git push` only (never `gx push` or `gx capture push`).

The important invariant is shared: agent commits must create the same revisions
that the human CLI path would create.

## Merge Standard

Merge only changes that support this contract:

- `git add` + `gx commit` (and hidden `gx generate` when needed) record reviewable revisions
- stack/revision metadata stores JJ change IDs and current commit IDs
- agents publish with plain `git push` only (never `gx push` / `gx publish` / MCP push tools)
- the pre-push hook captures session data and records publish metadata
- PRs are opened with `gh pr create` without seeding `## Summary`

Do not merge unrelated draft stacks just because they exist locally. Create a
clean integration stack from `main`, run focused tmp e2e tests plus
`go test ./...`, and only then publish with `git push` or land.

## Acceptance Gate

Before landing authoring changes, run:

```bash
go test ./...
go test -tags e2e ./test/e2e -run 'TestGXCompose|TestGXPublishAllPublishesEveryStack' -count=1
```

Then run a manual dogfood:

```bash
git add <files>
gx commit -m "describe this revision"
gx status
git push
gh pr create
git branch --show-current
```

The result should be boring: staged work becomes a GX revision, `git push`
publishes the branch through the pre-push hook, and the PR opens cleanly.
