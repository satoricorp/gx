# GX Context

This file defines the domain language GX reviews and agent workflows should use.
It is a glossary, not an implementation spec.

## Core Terms

### GX

The product layer over JJ and Git. GX owns the workflow for capturing coding
context, shaping reviewable revisions, publishing stacks, and preparing review
evidence.

GX is not a generic Git wrapper, a JJ teaching tool, or a generic VCS
abstraction.

### Session

Captured coding context from an agent, editor, or shell. A Session can contain
prompts, responses, tool calls, commands, tests, files touched, and timestamps.

Sessions are provenance only. They do not choose a Stack and do not own a
Revision.

### Revision

One reviewable logical change. A Revision is the GX product identity for a
change.

Locally, a Revision is backed by a JJ change ID. A Git commit is only the
current exported snapshot of that Revision and may change as the Revision is
edited.

### Stack

An ordered line of Revisions. Locally, a Stack is GX metadata plus a JJ/Git
compatible ref. Remotely, it can be published as stacked Git branches and PRs.

### Published Stack

A Stack that GX has exported to the remote Git or review surface.

### Compose

The normal GX authoring flow that turns messy working-copy changes into
reviewable Revisions grouped into Stacks.

Compose should present only proposals that are ready to apply or have explicit
repair guidance.

### Compose Proposal

An addressable proposal for how current working-copy changes should become
ordered Revisions and Stacks. A proposal must account for every selected hunk
exactly once before it is accepted.

### Provenance

Evidence linking a Revision back to Sessions and other captured context.
Provenance can be explicit, repo-local, linked, absent, or unavailable, but it
must not be treated as the source of truth for Stack routing.

### Review Bundle

The versioned publication payload that joins Stack entries, patches, linked
Sessions, provenance, structural facts, and review context.

### Review Context

The per-Revision view of evidence that the review surface consumes. It should
explain what provenance and structural signals are present or missing without
requiring the reader to parse opaque implementation records.

## Product Boundaries

### Authoring Engine

The Go module that owns GX authoring behavior such as compose, edit, status,
stacks, sync, and publish.

### CLI Adapter

The human-facing adapter over the Authoring Engine.

### MCP Adapter

The agent-facing adapter over the Authoring Engine. The TypeScript MCP server
currently shells to the CLI as a temporary adapter.

### VCS Service

The lower-level JJ, Git, and local storage implementation behind the Authoring
Engine.

## Internal Implementation Terms

### Demux

Internal implementation language for splitting messy work into reviewable
Revisions. Public product surfaces should say Compose, Compose Proposal,
Revision, Stack, or route-specific wording instead.

Demux can still appear in internal package names, private helper names, ADRs,
and migration notes when the implementation history matters.

## Invariants

- A Revision is the stable product identity; a Git commit is an exported
  snapshot.
- A Session is provenance only; it must not route work to a Stack.
- Normal authoring starts from the configured base branch and returns there
  after compose or publish.
- Accepted Compose proposals create visible Stacks and Revisions.
- Publishing scans accepted Stacks; it must not depend on the current checkout.
- Review output should use GX product terms before implementation terms.
