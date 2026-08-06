# gx Product Framing

## Core Definition

gx captures coding sessions and attaches them to revisions, then uses that context to publish stacked PRs and generate better review artifacts.

gx is not primarily:

- a Git wrapper
- a JJ teaching tool
- a generic VCS abstraction layer

JJ is the local engine. Git is the remote compatibility layer. gx is the product surface.

## How gx Differs From Graphite

Graphite helps developers manually structure and publish stacked PRs.

gx is intended to:

- capture implementation context
- attach that context to work
- derive reviewable revisions
- publish those revisions as stacked PRs
- generate review summaries from the evidence

Graphite centers the branch or diff as the main unit.

gx centers the revision:

- linked sessions
- linked prompts and responses
- linked files
- linked tests
- linked revisions
- linked PRs

In gx, stacked PRs are an output format, not the whole product.

## Product Thesis

The product value is upstream of stacked PR publication:

- captured context
- inferred work boundaries
- review evidence
- scope and failure-mode summaries

That means gx still has value before perfect stacked PR publishing exists.

## Product Objects

gx should revolve around three core objects:

1. Session
   Prompts, responses, tool calls, commands, tests, files touched.

2. Revision
   A logical, reviewable slice of work with attached session evidence.

3. Stack
   An ordered set of revisions that becomes stacked PRs.

## Why JJ Fits

JJ is useful because gx is about evolving revisions, not raw commits.

JJ provides:

- stable logical change identity
- rewrite-friendly local history operations
- better support for reshaping revisions before publication

But JJ should not be the user-facing product. gx should hide JJ’s awkward parts and own the workflow.

## Product Boundary

Recommended framing:

- gx is the tool the developer uses.
- JJ is the internal local change engine.
- Git is the publication and compatibility format.

The user should not have to think about bookmarks, detached HEAD, or manual ref movement in the normal flow.

## Recommended User Workflow

The intended flow should look like:

```bash
gx init
gx codex
git add .
gx commit -m "scaffold app"
git add .
gx commit -m "add database schema"
git add .
gx commit -m "add secure email ingestion"
git push
gh pr create
```

gx should then:

- attach session context to each revision
- preserve the ordering of the revisions
- publish them as stacked PRs (via `git push` + `gh pr create`)
- generate review summaries for each PR (gx Cloud comment)

## Strategic Implication

gx is not just “stacked PRs made easier.”

gx is a context-native software delivery tool:

- it records how work was done
- attaches that evidence to reviewable revisions
- and turns those revisions into better PRs and better reviews

## Roadmap Implication

The product roadmap should prioritize:

1. session to revision attachment
2. revision creation and modification
3. stack publication
4. PR summary generation
5. cloud lookup by revision, commit, and PR

That is more important than Git parity or JJ parity.
