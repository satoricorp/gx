# GX CLI Authoring Contract

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
gx compose
gx stacks
gx publish
```

`gx compose` looks at all unrecorded work, proposes stacks and revisions,
verifies that applying the proposal will actually work, and lets the caller
accept all or selected stacks.

`gx stacks` shows every accepted local stack and its revisions, regardless of
the current Git branch. Once a stack has been merged into its target branch, it
is no longer part of the `gx stacks` surface.

`gx publish` publishes every accepted stack that has unpublished revisions. It
should not depend on a "current stack," because normal operation is from `main`.
If there are no accepted stacks, or all accepted stacks are already published,
publish is a successful no-op.

The user should stay attached to `main` after compose and publish. Detached HEAD
is not an acceptable steady state.

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
- Edit surgery may use `gx/edit` or another internal branch, but only for
  explicit edit workflows.

Branches/bookmarks created by compose should match the stack names shown in
compose. If compose proposes `feature/stack-management`, then accepting it
should create a stack whose visible name/ref in `gx stacks` is
`feature/stack-management`.

## Publish

`gx publish` should publish all accepted stacks with unpublished revisions.

The command should scan GX's stack metadata, find draft stacks with unpublished
revisions, push their branch refs, record publish metadata, and update remote
refs. It should not care which stack the user is "on."

Expected behavior:

```bash
gx publish
```

means:

- publish every unpublished accepted stack
- push each stack branch
- record remote ref
- mark revisions published
- upload review context when cloud/GitHub integration is enabled
- return to `main`

The CLI does not expose a `--no-github` publish path. If GX needs a review-only
dry run, that should be a separate command or explicitly named flag because
`publish` strongly implies branch refs are pushed.

## MCP And Agent Flow

MCP should call the same authoring engine behavior as the CLI.

The ideal agent loop is:

1. `gx_compose_changes`
2. inspect proposal and review state
3. repair with the model if GX reports issues
4. `gx_review_compose_plan`
5. repeat until ready
6. `gx_accept_compose_plan`
7. `gx_publish`

For MCP, `gx compose --json --plan` can avoid spending GX's configured OpenAI
tokens. The caller's model can do the repair. If `use_gx_llm` is enabled, GX can
run its own OpenAI repair loop.

The important invariant is shared: agent acceptance must create the same stacks
and revisions that the human TUI would create.

## Merge Standard

Merge only changes that support this contract:

- compose applies all proposed stacks, not only a valid subset
- selected-stack apply filters hunks correctly
- compose preflight runs in a disposable attempt and feeds failures into repair
- accepted compose stacks show in `gx stacks` under the proposed branch names
- stack/revision metadata stores JJ change IDs and current commit IDs
- publish scans all unpublished accepted stacks
- publish pushes branch refs and records remote refs
- compose and publish return to `main`
- e2e tests cover multi-stack compose, selected-stack accept, `gx stacks`
  visibility, and publish-all

Do not merge unrelated draft stacks just because they exist locally. Create a
clean integration stack from `main`, run focused tmp e2e tests plus
`go test ./...`, and only then publish or land.

## Acceptance Gate

Before landing authoring changes, run:

```bash
go test ./...
go test -tags e2e ./test/e2e -run 'TestGXCompose|TestGXPublishAllPublishesEveryStack' -count=1
```

Then run a manual dogfood:

```bash
gx compose
gx stacks
gx publish
git branch --show-current
```

The result should be boring: proposed stacks become visible stacks, publish
pushes all unpublished stacks, and the final branch is `main`.
