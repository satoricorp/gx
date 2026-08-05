# ADR 0001: Use lgtm Product Language At Public Seams

## Status

Accepted

## Context

lgtm has three layers that talk about the same work:

- the product layer, where users think in Sessions, Revisions, Stacks, Compose,
  and Publish
- the local implementation layer, where JJ change IDs, Git commits, bookmarks,
  and storage rows make those objects durable
- the adapter layer, where CLI and MCP expose lgtm behavior to humans and agents

Early authoring code used "demux" for the algorithm that splits messy
working-copy changes into reviewable slices. That term is still useful inside
the implementation, but it is not the product concept users should learn.

## Decision

Public lgtm seams use product language:

- Session
- Revision
- Stack
- Published Stack
- Compose
- Compose Proposal
- Publish
- Review Bundle
- Review Context

The stable product identity for a change is the lgtm Revision. Locally, a Revision
is backed by a JJ change ID. A Git commit is the current exported snapshot, not
the durable identity.

The term "demux" is internal implementation language. It can remain in private
helpers, internal package code, migration notes, and ADRs when the implementation
history matters, but user-facing docs, CLI help, MCP tool descriptions, and
review output should prefer Compose, Compose Proposal, Revision, and Stack.

## Consequences

- Architecture review should flag internal implementation terms leaking through
  public seams when `CONTEXT.md` marks them as internal.
- CLI and MCP adapters should converge on compose-named request and response
  shapes, even if the underlying implementation still calls demux helpers.
- Tests and docs should describe product behavior in lgtm terms, not JJ or Git
  implementation details unless the behavior depends on those details.
- ADRs should record future cases where lgtm intentionally exposes implementation
  language because the trade-off is worth it.
