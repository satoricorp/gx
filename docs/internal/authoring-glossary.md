# gx Authoring Glossary

gx keeps several names for the same reviewable revision because each layer has its
own wire format. Use these mappings when changing authoring code.

## Core Mapping

- **Revision**: the product term for one reviewable logical change.
- **JJ change**: the durable local VCS object for a revision. This is the
  stable key to use for local provenance because JJ rewrites Git commits.
- **Stack**: the named line of revisions. Locally it is backed by gx metadata
  and a JJ/Git-compatible conventional ref such as `feature/login-flow`.
- **Published stack**: a stack that has been exported to the remote review/Git
  surface.
- **Stack entry**: the review-surface view of a published revision.

In short:

```text
Revision = JJ change = stack entry
Stack = local container/export ref for one line of revisions
Published stack = remote wire format for that stack
```

## Provenance

- **Session**: captured agent or editor context. Sessions are provenance only;
  they do not own or route stacks or revisions.
- **Exact provenance**: `change_sessions(change_id, session_id)` created from an
  explicit handoff such as `GX_SESSION_ID` / `GX_SESSION_IDS`, or from an
  unlinked captured session whose `repo_root` / `cwd` matches the repo being
  recorded. This is still provenance only; it never chooses a stack.
- **Fuzzy provenance**: a retrieval hint based on time/files/session metadata.
  It may scope search, but must not be presented as exact lineage.
- **Provenance attachment module**: the shared module that classifies exact
  provenance as explicit, repo-local, or absent, then attaches resolved sessions
  to recorded revisions. Storage only persists the result.

## Engine Boundary

- **Authoring Engine**: the Go module that owns checkpoint, modify, switch,
  status, publish, and future demux/reorder operations.
- **CLI**: a human adapter over the Authoring Engine.
- **MCP**: an agent adapter over the Authoring Engine. The current TypeScript
  MCP shells out to the CLI as a temporary adapter; a future Go MCP should call
  the Authoring Engine in process.
- **VCS service**: the lower-level JJ/Git/storage implementation behind the
  Authoring Engine.

## Demuxing

- **Checkpoint**: record one already-identified revision.
- **Demux proposal**: an addressable proposal that splits messy working copy
  changes into ordered revisions before applying JJ surgery. The proposal must
  be persisted before adapters expose propose-then-apply flows.
- **Hunk catalog**: the top-level list of candidate hunks in a demux proposal.
  Revision plans should refer to entries by `hunk_ids` instead of copying patch
  payloads into each revision.
- **Demux plan packet**: the MCP-facing planning object returned by
  `gx demux --json` / `gx_demux_changes`. It contains the demux
  proposal, the hunk catalog, and the schema the host model should pass to
  `gx_apply_revision_plan`. The Authoring Engine also runs the first validation
  pass and reports whether the plan is ready to apply, has recommended repairs,
  or requires repair before apply.
- **Review plan**: normalize and check an LLM-authored revision plan without
  applying JJ changes. This is the repair loop gate for MCP: draft, review,
  revise from structured `repair_hints`, review again, then apply.
  Hints can request missing hunk assignment, `depends_on` metadata, or revision
  reordering when structural facts show a bottom-up stack would be incoherent.
  Applying a plan blocks on warning-severity feasibility warnings unless the
  caller explicitly accepts them with `--allow-warnings`.
- **Structural fact**: a local deterministic hint extracted from changed files,
  such as defined symbols, referenced symbols, and file-to-file dependency
  edges. Demux uses these hints for ordering; they are not semantic retrieval.
- **Changed symbol**: the nearest locally inferred symbol enclosing a changed
  hunk. MCP demux planning can use changed symbols to group hunks by
  function/type before applying a revision plan. Deterministic demux may also
  split one file into hunk-level revisions when different hunks map cleanly to
  different changed symbols.
- **Feasibility warning**: a structured warning attached to a demux proposal
  when gx can see that a proposed revision may not be reviewable bottom-up.
  Examples include unmapped hunks, dependency-order conflicts, and separated
  source/test counterparts. Info-level inferred dependency warnings make
  structural dependencies explicit even when the order is already coherent.
- **Apply proposal**: perform the JJ split/reorder operations and attach exact
  provenance to the resulting revisions. Applying validates that each proposed
  hunk is covered exactly once, either by `hunk_ids` or by one whole-file
  revision for that file.
- **Demux route**: the proposed stack target for one proposed revision. An
  omitted route means the revision stays on the current stack,
  preserving today's single-stack demux behavior.
- **Routing evidence**: the heuristic, AI-authored, or user-authored reason a
  demux route points at a specific stack. Routing evidence is a
  planning signal only; exact provenance still comes from session attachment.
- **Route review**: the validation pass that checks demux routes before apply.
  It verifies target stacks exist or are intentionally new, reports ambiguous
  routes, and emits repair hints so AI can resolve invalid or incoherent routing
  before JJ surgery runs.
- **Demux evidence**: durable per-revision metadata written after apply and
  included in review bundles. It links the resulting JJ change to the demux
  proposal revision, covered hunk ids, files, confidence, provenance status,
  feasibility warnings, and relevant structural context.
- **Review bundle**: the versioned push-time payload assembled for the review
  surface. It joins stack entries, patches, linked sessions, and demux evidence.
  Cloud sync uploads this bundle; it does not own the bundle shape.
- **Review publication**: the push-time module that assembles the Review bundle,
  runs optional semantic indexing, and calls an upload adapter. It owns
  non-blocking indexing policy; the cloud module only uploads.
- **Review context**: the per-revision view inside a review bundle. It exposes
  provenance status, per-revision provenance sources, transcript source IDs,
  linked session count, structural signal availability, structural facts,
  changed symbols, feasibility warnings, typed evidence entries, and
  deterministic risk so the console does not need to parse opaque evidence JSON
  before rendering triage.
  Demux evidence fills the rich structural fields when present; revisions
  without demux evidence still get a context shell marked with unavailable
  structural evidence and absent or linked provenance.
- **Semantic transcript index**: an optional retrieval index built from
  revision-scoped transcript source IDs. It stores transcript chunks in
  TurboPuffer with vectors from OpenAI embeddings. It is not provenance truth;
  every indexed chunk points back to session/request/response IDs from the
  review bundle.
- **Local semantic index adapter**: the current opt-in adapter that indexes
  Review bundle transcript chunks from the author's machine during publication.
  It sits behind the publication indexing seam so cloud ingest can later provide
  another adapter after explicit sync consent.
- **Console ingest**: the review-surface module that turns a Review bundle into
  renderable revisions, risk/context fields, transcript source availability, and
  a session index for "why this revision?" retrieval.
