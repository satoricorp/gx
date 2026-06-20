# gx

`gx` captures terminal agent traffic and exposes a small workflow CLI where gx adds behavior on top of JJ and Git.

This repo also ships an installable skill at `skills/gx` that teaches the intended gx/JJ/Git command boundary.

## License

GX is open source under the GNU Affero General Public License v3.0. Modified
versions offered over a network must provide their corresponding source under
the AGPL.

## Install `just`

Install `just` with Homebrew:

```bash
brew install just
```

## Build

Build a local binary in the repo:

```bash
cd ~/git/gx
go build -o gx ./cmd/gx
```

Or with `just`:

```bash
cd ~/git/gx
just build
```

## Install on PATH

Your shell already includes `~/.local/bin` in `PATH`, so install `gx` there:

```bash
cd ~/git/gx
go build -o ~/.local/bin/gx ./cmd/gx
```

Or with `just`:

```bash
cd ~/git/gx
just install
```

`just install` also installs Bash completion to `~/.local/share/bash-completion/completions/gx`
and Zsh completion to `~/.zfunc/_gx`. `gx` does not expose a `completion` command in the CLI.

After code changes, rebuild with the same command to update the installed binary and completions.

If your Zsh setup does not already load `~/.zfunc`, add this once to `~/.zshrc`:

```bash
fpath=(~/.zfunc $fpath)
autoload -Uz compinit && compinit
```

Verify:

```bash
command -v gx
gx version
```

On first run, `gx` tells you to run `gx init`.

## Run

### Menu-Bar App

Build the native macOS menu-bar app:

```bash
just menubar
```

This creates:

- `apps/menubar/dist/GX.app`
- `apps/menubar/dist/GX-macOS.zip`
- `apps/menubar/dist/GX-macOS.dmg`

Users download `GX-macOS.dmg`, drag `GX.app` to Applications, and launch it.
On launch, the app installs the bundled CLI to `~/.local/bin/gx`, shows
`gx doctor` stats and activity, links to the GX console, and provides MCP
setup instructions.

The `Build macOS App` GitHub Action builds ZIP/DMG artifacts from this repo.
Pushing a `v*` tag also uploads those assets to the matching GitHub Release.
Configure repository variables `GITHUB_CLIENT_ID`, `CONVEX_SITE_URL`, and
`GX_CLOUD_URL` to bake production endpoints into release builds.

### Capture

```bash
gx init
```

`gx init` installs the repository pre-push hook. On push, gx stages capture
context from Claude, Codex, and Cursor transcripts for the pushed commit range.

Check setup and drain queued uploads:

```bash
gx doctor
gx capture sync
```

### Captured Agent Sessions

Captured sessions are recorded in `~/.gx/gx.db`.
When a captured session originates inside a repo, the next `gx add`, `gx edit`,
or `gx demux` can attach unlinked captured sessions from that repo to the GX
revision. Sessions provide provenance only; they do not route work to a body,
bookmark, or branch.

Initialize explicitly:

```bash
gx init
```

The first interactive `gx init` explains that gx stores your name/email in `~/.gx/config.json`,
then prompts for both values.

Or set identity directly during init:

```bash
gx init --name "Some One" --email "someone@example.com"
```

### JJ-Backed Workflow

`gx` now includes workflow commands where gx adds workflow or metadata:

```bash
gx add -m "describe this revision"
gx add --interactive -m "describe selected changes"
gx add --hunk --patch-file /tmp/selected.patch -m "describe selected hunks"
gx demux
gx demux --plan
gx demux --json
gx demux list
gx demux proposals
gx demux review d1
gx demux fix d1
gx demux show d1
gx demux show u1
gx demux apply <proposal-id>
gx edit
gx status
gx stacks
gx sync
gx publish
```

Current behavior:

- `gx add` finalizes the current JJ-backed revision and starts the next one.
- `gx add -m "message" <filesets...>` records only matching filesets and leaves the remaining changes in the next revision.
- `gx add --interactive` opens JJ's interactive split editor for a human to select hunks or lines, records the selected changes, and leaves the remaining changes in the next revision.
- `gx add --hunk --patch-file <path>` splits using selected hunks from a unified patch file for agent-friendly non-interactive workflows.
- `gx add` always prints the recorded change hash, commit hash, description, split summary, and edit commands.
- `gx demux` proposes ordered revisions from the current working copy, applies deterministic repair passes, then asks OpenAI to refine the proposal when local repair cannot finish it.
- `gx demux --plan` stops after local deterministic planning and repair; it never calls OpenAI.
- `gx demux --json` returns the machine-facing demux workflow packet used by MCP: proposal, review result, workflow state, and review/apply tool guidance.
- `gx demux` accepts filesets plus `--exclude <fileset>` so humans and agents can demux only part of a dirty working copy.
- `gx demux list` lists saved proposals for the current repo; `gx demux proposals` is an alias.
- `gx demux review d1` reviews a saved proposal and prints a compact deterministic state summary without requiring a temporary plan file. Full feasibility warnings, repair hints, and info diagnostics are available with `--raw` or `--json`.
- `gx demux fix d1` runs the same repair pipeline for a saved proposal: deterministic repair first, then OpenAI when local repair cannot finish it. Pass `--plan` for local-only repair.
- `gx demux show d1` shows a saved proposal, where `d1` is the newest pending proposal from `gx demux list`.
- `gx demux show u1` shows one proposed revision from the latest pending proposal; pass `--proposal <id-or-dN>` to inspect a revision from another proposal.
- Agent-only commands such as `gx demux review-plan` and `gx demux apply-plan` are hidden from human help, but remain callable by MCP/agent workflows.
- The demux JSON packet includes a top-level hunk catalog, lightweight local structural facts, and dependency edges for MCP/LLM-authored revision plans.
- Changed hunks are annotated with enclosing symbols when GX can infer them locally.
- When one file has multiple changed symbols, deterministic demux can propose separate hunk-level revisions before the MCP/LLM refines the plan.
- Demux proposals include `feasibility_warnings` for deterministic concerns such as unmapped hunks, dependency-order conflicts, inferred `depends_on` hints, and separated test/source counterparts.
- If no explicit `GX_SESSION_ID(S)` is present, `gx demux` looks for unlinked captured sessions from the same repo and marks those revisions with `provenance_status: "repo_local"`.
- Provenance attachment is centralized: `gx add`, `gx edit`, and `gx demux` all classify sessions as `explicit`, `repo_local`, or `absent` through the same local provenance logic.
- Hidden `gx demux review-plan --plan-file <path>` reviews an LLM-authored revision plan without applying JJ changes; invalid, under-specified, and structurally misordered plans return structured JSON errors or `repair_hints`.
- Hidden `gx demux apply-plan --plan-file <path>` applies a revision plan; hunk-level revisions can refer to `hunk_ids` instead of copying patch payloads. GX checks that every proposed hunk is covered exactly once before applying, and blocks structural warning-severity plans unless `--allow-warnings` is passed after review.
- Applied demux records are persisted per revision and included in push review bundles as `demux_evidence`.
- Push review bundles are assembled by a dedicated local review-bundle module before cloud upload, so the versioned review payload shape is testable without the HTTP adapter.
- Review publication is coordinated by a dedicated publication module; the cloud client is only the HTTP upload adapter, and semantic indexing sits behind an optional indexing seam.
- Review bundles expose first-class `review_context` per revision, including provenance status, per-revision provenance sources, transcript source IDs, linked session count, structural signal availability, structural facts, changed symbols, feasibility warnings, typed evidence entries (`provenance`, `structural`, `risk`), and an initial deterministic risk score derived from local evidence. Revisions without demux evidence still get an honest context shell so the review surface can distinguish missing provenance from unavailable structural evidence.
- A cloud ingest module turns review bundles into stored review artifacts, risk/context fields, transcript source availability, and a session index for the review surface.
- MCP callers should use `gx_compose_changes` for messy diffs. It shells to `gx compose --json --plan` by default, which returns the hunk catalog, an initial review result, and a state without calling GX's configured OpenAI repair model. Apply directly when ready; when repair is recommended or required, revise the proposal with the calling codegen's LLM, review again with `gx_review_compose_plan`, then call `gx_accept_compose_plan`. Pass `use_gx_llm: true` only when the user explicitly wants GX to spend its configured LLM tokens.
- `gx init` initializes a JJ-backed GX repo for the current working tree, or reuses the nearest existing one.
- `gx init` also sets up gx user identity and writes JJ `user.name` / `user.email`.
- if `GX_POSTLIST_URL` is set, `gx init` also sends identity to gx signup
- ambient Codex/Claude sessions are attached to the next `gx add` or `gx edit` for that repo
- `gx edit [rev]` re-enters an existing JJ change.
- `gx edit` with no revision opens an interactive picker of recent mutable changes.
- `gx modify` remains as a compatibility alias for `gx edit`.
- `gx status` shows the current revision, message state, changed files, and next commands. gx stores new changes in revisions, so `git status` may be clean.
- `gx stacks` browses GX stacks and revisions. The interactive view supports stack navigation, revision navigation, edit, and diff. Use `gx edit` to re-enter a revision and continue working on its stack.
- `gx sync` fetches and prunes the remote line of work (default `origin`).
- `gx publish` pushes the recorded GX stack as branch refs, records GX metadata, and prints branch-mapped PR links.
- `gx pr` remains as a hidden compatibility alias for agent and script workflows.
- release builds upload pushed stack context to gx cloud automatically; `gx auth login` works without env setup
- contributors can override cloud endpoints with `CONVEX_SITE_URL`, `GITHUB_CLIENT_ID`, and `GX_CLOUD_URL`
- shell completions are installed with `just install`, not through a CLI subcommand

For pure JJ operations that gx does not extend, use `jj` directly. That includes history surgery and inspection commands such as `jj squash`, `jj log`, and `jj bookmark list`.

Use Git directly only for compatibility and remote inspection, such as checking the current branch, refs, remotes, and raw commit history.

`gx publish` currently expects a normal Git branch when no explicit refspec is provided.
If there is no current branch, pass an explicit push target or create a branch first.

## Data Directory

By default, data is stored here:

```bash
~/.gx/gx.db
```

To isolate GX state for testing or local experiments, override the data directory:

```bash
GX_HOME=/tmp/my-gx-data gx init --name "Some One" --email "someone@example.com"
```

gx user identity is stored in `~/.gx/config.json`.

## gx Signup

If you want to collect identities centrally, set:

```bash
export GX_POSTLIST_URL="https://<your-neon-data-api>/gx_identities"
export GX_POSTLIST_API_KEY="<optional-api-key>"
```

`gx init` will send the current identity to gx signup after local setup.

The expected Neon schema is in:

```bash
docs/postlist-neon.sql
```

## gx Cloud

Shipped `gx` binaries include production cloud endpoints. No env setup is required for
`gx auth login` or `gx publish` cloud upload.

```bash
gx auth login
gx auth status
gx publish
```

After `gx publish`, gx uploads the pushed change, changed files, linked sessions, requests,
responses, token usage, and the GitHub PR URL when available. When the cloud API returns
a review URL, `gx publish` prints it.

### Semantic transcript indexing

Semantic indexing is opt-in because it sends linked session transcript chunks to
OpenAI for embeddings and to TurboPuffer for vector storage. The authoritative
review data remains the relational review bundle; the vector index is only a
retrieval layer for review chat.

```bash
export GX_SEMANTIC_INDEX=1
export OPENAI_API_KEY="<openai-api-key>"
export TURBOPUFFER_API_KEY="<turbopuffer-api-key>"
export GX_TPUF_NAMESPACE="gx-sessions"
gx publish
```

Optional overrides:

```bash
export GX_OPENAI_BASE_URL="https://api.openai.com"
export GX_OPENAI_EMBEDDING_MODEL="text-embedding-3-small"
export GX_EMBEDDING_DIMENSIONS=512
export GX_TPUF_BASE_URL="https://gcp-us-central1.turbopuffer.com"
export GX_SEMANTIC_BATCH_SIZE=64
export GX_SEMANTIC_MAX_CHUNK_BYTES=12000
```

If semantic indexing is misconfigured or unavailable, `gx publish` still uploads the
review bundle and reports the semantic indexing error separately.

### Local development overrides

Copy `.env.example` to `.env` and fill in values. The justfile loads `.env`
automatically for `just build`, `just build-release`, `just run`, and other recipes.

`just build` bakes `GX_*` into the binary (local dev). Runtime `GX_*` exports still
override those baked defaults if set.

```bash
cp .env.example .env
```

Or export overrides manually:

```bash
export CONVEX_SITE_URL="https://<deployment>.convex.site"
export GITHUB_CLIENT_ID="<oauth-app-client-id>"
export GX_CLOUD_URL="http://localhost:3201"
```

Env vars take precedence over baked defaults. A dev build with no env vars and no
`-ldflags` injection reports a clear configuration error on `gx auth login`.
`CONVEX_SITE_URL` points at the Convex HTTP site for auth. `GX_CLOUD_URL`
points at the Postgres-backed server origin used by upload, bookmarks, and
OpenAI proxy routes.

Running `./gx` directly does not load `.env`; use `just run …` or `source .env` first.

### Release builds

Release maintainers inject production values at link time via `.env` or exported vars:

```bash
just build-release   # bakes the current non-release env values from .env
```

Example `.env` release entries:

```bash
GITHUB_CLIENT_ID="<github-oauth-app-client-id>"
CONVEX_SITE_URL="https://<prod-deployment>.convex.site"
GX_CLOUD_URL="https://<api-host>"
VERSION="1.2.3"
```

`GX_CLOUD_URL` is the deployed GX API origin. The CLI appends API paths such as
`/v1/publish`; auth uses `CONVEX_SITE_URL` separately (`/cx/auth/complete`,
`/cx/auth/revoke`). `GX_CLOUD_URL` should not point at the Convex `.site` host.

Only the GitHub OAuth client ID is embedded in the binary — never the client secret.

## Captured Data

Quick inspection:

```bash
sqlite3 ~/.gx/gx.db "select count(*) from sessions; select count(*) from requests; select count(*) from responses;"
sqlite3 ~/.gx/gx.db "select provider, endpoint, method, model from requests order by created_at desc limit 5;"
sqlite3 ~/.gx/gx.db "select status_code, is_streaming, input_tokens, output_tokens from responses order by created_at desc limit 5;"
```

The same database now also stores initial GX workflow metadata for:

- repos
- GX changes
- change revisions
- modify events
- push events
