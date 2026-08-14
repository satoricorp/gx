# Review Knowledge Index

This directory builds the versioned TurboPuffer corpora used by `gx enhance
--deep` and GitHub PR review generation.

The corpus is retrieval context, not a replacement for gx heuristics. Keep the
deterministic review rules in code. Use this index for official standards,
language guidance, tool guidance, and repo-independent review knowledge that a
model can cite or use as supporting context.

## Files

- `sources.yaml`: approved v2 source manifest. This supersedes `urls.json`.
- `golden_queries.yaml`: 68-query eval gate manifest.
- `index_review_resources.py`: validates, fetches, normalizes, chunks, embeds,
  indexes, refreshes, and scaffolds eval output.
- `urls.json`: retained only as the v1 manifest/reference.

## Source Policy

Prefer primary or official sources:

- standards: NIST SSDF, OWASP ASVS/Cheat Sheets, OpenSSF Scorecard, SLSA
- secure-coding bodies: SEI CERT, MISRA, ISO/IEC secure-coding references
- large-scale practice: Google, Microsoft, GitLab, Mozilla
- language docs: TypeScript, JavaScript, Python, Go, Rust, Java, C, C++,
  SQL, C#, shell, PHP, Kotlin, Swift, Ruby, Dart
- framework and vendor docs: Spring/Spring Boot, Django, React, Next.js, Vue,
  NestJS, Express, Rails, .NET, PostgreSQL, MySQL, SQL Server, dbt, Android,
  Apple platforms
- ORM and data-access semantics: Jakarta Persistence, Hibernate, Active Record,
  SQLAlchemy, Prisma, GORM
- concurrency and async contracts: the Java memory model (JLS 17),
  java.util.concurrent, virtual threads, the Go memory model and `context`,
  asyncio, the Node/JS event loop, .NET async
- API evolution: SemVer, protobuf/Buf breaking-change rules, Cargo SemVer,
  Go module compatibility, .NET library change rules
- distributed-systems correctness: the SRE book, gRPC deadlines/retries,
  idempotency-key and delivery-semantics contracts
- tool docs: typescript-eslint, Ruff, mypy, Bandit, govulncheck, Clippy, Miri,
  RustSec, cargo-deny, Error Prone, SpotBugs, Checkstyle, PMD, clang-tidy,
  Cppcheck, SQLFluff, Roslyn analyzers, ShellCheck, PHPStan, detekt, SwiftLint

Do not add blog posts or generated summaries as high-authority sources unless
there is no primary source. If a seed note is opinionated, keep it in
`seed_notes` and keep the `authority`/`evidence_level` honest.

`fetch: github-md` retrieves a single raw markdown file — link crawling only
runs for `fetch: html`. To ingest a whole directory of docs (the OWASP cheat
sheets, the ASVS chapters, a Sphinx `docs/source`), point the source at a
`github.com/.../tree/` URL with `crawl_depth: 1`: the fetcher expands it via
the GitHub contents API into every `.md`, `.mdx`, and `.rst` file in that
directory. Subdirectories recurse while `crawl_depth` allows (e.g. sqlfluff
uses `crawl_depth: 2` to reach `docs/source/*/`). `.rst` payloads pass through
a best-effort reStructuredText-to-markdown conversion so Sphinx manuals chunk
along their real section structure. Files that another manifest entry ingests
on its own are skipped during expansion so one document never enters the
corpus under two source ids.

Crawls cap their URL count at 50 for `fetch: html` and 200 for
`fetch: github-md`; set `max_urls` on a source to override (the shellcheck
wiki uses 450 to cover every SC check page). A `fetch: manual` source is
skipped at chunk time until its normalized markdown is dropped at
`corpus/normalized/<id>.md`, at which point it chunks like any other source.

robots.txt is fetched with the indexer's own User-Agent. Several sources
(readthedocs sites, smartbear, nvlpubs) used to fail as "robots.txt
disallows" only because the host rejected the default Python-urllib agent on
the robots.txt request itself and robotparser treats that rejection as
disallow-everything; their actual policies permit the fetch.

## Current Coverage

The v2 manifest covers review process, cross-cutting security, TypeScript,
Python, Go, Rust, C#/.NET, Java/Kotlin, Ruby, C/C++, Swift, PHP, Dart/Flutter,
SQL, shell, IaC, style/linter sources, and empirical research. Research sources
are routed to `research-corpus-v1` and are not used for review-time retrieval.

Style and lint sources say how code should *look*. They do not say what a
framework or library *promises*, which is why API-contract misuse measured as
the weakest retrieval category (~36% recall) while TypeScript and Go style were
well covered. The 2026-08-12 contract gap fill added six layers aimed at that:
framework contracts (Spring/Spring Boot, Django, React, Next.js, Vue, NestJS,
Express), ORM and data-access semantics (Jakarta Persistence, Hibernate, Active
Record, SQLAlchemy incl. pooling, Prisma, GORM), per-language concurrency and
async, API evolution and compatibility, distributed-systems correctness, and
Ruby/Java rule-catalog depth. Prefer these over style guides for "is this a
correct use of X" questions; keep them below repo-local `REVIEW.md` policy.

Two known limits on how far a `languages:` tag carries a source. First,
`languageTagsForFiles` in `internal/codereview/review_resources.go` maps file
extensions to language tags and has no case for `.rb`, `.dart`, `.tf`, `.proto`,
or `.html`, so sources tagged with those languages are only reachable through
the unfiltered retrieval leg, never the signal-filtered one. Second, the v2
index has no `frameworks` column (see `review_corpus_schema`), so a
framework-specific source can only be tagged by its host language; framework
targeting rides on the BM25 leg matching framework vocabulary in the body.

The six .NET documents from the v1 `urls.json` manifest are re-ingested here as
v2 sources (same ids where possible) rather than served from their legacy
production rows: runtime retrieval filters on `source_kind = "review_corpus"`,
which the v1 rows (`source_kind = "review_knowledge"`) never match, and the
legacy rows lack the v2 filter attributes. (As of 2026-08-06 the legacy rows
turned out to be gone entirely — the earlier `totality-review-knowledge`
production namespace was deleted in a namespace cleanup, so the v2 re-ingest is
the only copy. The 2026-08-06 promote created `gx-review-knowledge` fresh, the
namespace the shipped retriever queries.)

Treat SQL and shell as cross-cutting review profiles as much as languages. A
TypeScript, Python, Java, Go, Rust, or C# change can still need SQL review when
it changes queries, migrations, indexes, data retention, permissions, or dbt
models. Any language can need shell review when it changes CI, package scripts,
deployment scripts, release automation, destructive commands, or operational
runbooks.

## Environment

Required for indexing:

```bash
export OPENAI_API_KEY=...
export TURBOPUFFER_API_KEY=...
```

Optional:

```bash
export GX_OPENAI_BASE_URL=https://api.openai.com
export GX_OPENAI_EMBEDDING_MODEL=text-embedding-3-small
export GX_EMBEDDING_DIMENSIONS=512
export GX_TPUF_BASE_URL=https://gcp-us-central1.turbopuffer.com
export GX_REVIEW_CANDIDATE_NAMESPACE=review-corpus-v2
export GX_REVIEW_KNOWLEDGE_NAMESPACE=gx-review-knowledge
export GX_RESEARCH_CORPUS_NAMESPACE=research-corpus-v1
export GX_REVIEW_RESOURCES=1
export GX_REVIEW_RESOURCES_TOP_K=8
```

The default namespace is intentionally separate from `GX_TPUF_NAMESPACE`
(`gx-sessions`) so review knowledge does not mix with session transcripts and
repository code chunks.

`gx enhance` queries this namespace when both `OPENAI_API_KEY` and
`TURBOPUFFER_API_KEY` are available. Set `GX_REVIEW_RESOURCES=0` to disable
review-resource retrieval for a run. `GX_REVIEW_RESOURCES_TOP_K` controls the
shallow retrieval limit; the default is 8 so repo-local policy files can remain
in the model context alongside review resources. `gx enhance --deep` raises the
minimum resource limit to 24.

The v2 index stores filterable metadata including `tier`, `languages`,
`authority`, `precedence_group`, `superseded_by`, `historical`, `section_path`,
`cwe_ids`, `content_sha256`, `fetched_at`, and `usefulness_rank`. The `body`
field is full-text indexed for BM25; runtime retrieval uses vector + BM25
multi-query RRF.

## Usage

Validate the manifest without network or API calls:

```bash
python3 scripts/review-knowledge/index_review_resources.py dry-run
```

Fetch and normalize sources:

```bash
python3 scripts/review-knowledge/index_review_resources.py fetch
```

Chunk normalized sources and write build artifacts:

```bash
python3 scripts/review-knowledge/index_review_resources.py chunk
```

`corpus/raw/` and `corpus/chunks/*.jsonl` are gitignored, so a fresh clone has
neither. Neither is an input: `fetch` writes `raw/` and nothing reads it back,
and `chunk` rebuilds the chunk files from `normalized/` offline, with no network
or API calls. Run `chunk` after cloning if you need them locally. They stopped
being committed once both had passed GitHub's 100MB per-file hard limit, which
rejects the push outright.

`corpus/normalized/` and `corpus/metadata/` stay committed. `normalized/` is the
real input to `chunk` and `index`, and for sources that refuse the crawler its
markdown is placed there by hand, so it does not survive a re-`fetch`.

Index normalized chunks into candidate `review-corpus-v2` and
`research-corpus-v1`:

```bash
python3 scripts/review-knowledge/index_review_resources.py index
```

Smoke test a small subset without network/API calls:

```bash
python3 scripts/review-knowledge/index_review_resources.py dry-run --limit 3
python3 scripts/review-knowledge/index_review_resources.py chunk --limit 3 --seed-only
```

Run the eval gate. Without API credentials this writes a skipped report; with
`OPENAI_API_KEY` and `TURBOPUFFER_API_KEY` it compares `gx-review-knowledge`
against `review-corpus-v2` and writes recall/precedence results:

```bash
python3 scripts/review-knowledge/index_review_resources.py eval
```

After the eval gate passes, promote the same v2 review chunks into the
production namespace (`gx-review-knowledge`) so production can cut over without
changing namespace names:

```bash
python3 scripts/review-knowledge/index_review_resources.py promote
```

`promote` upserts `source_kind=review_corpus` rows into `gx-review-knowledge`.
It does not delete legacy `source_kind=review_knowledge` rows, so old readers
continue working during the application rollout.

## Retrieval Shape

Future review retrieval should query this namespace using signals from the diff:

- language and framework
- changed file category
- risk profile: auth, data migration, dependency, concurrency, unsafe Rust,
  package release, public API
- repo `REVIEW.md` policy text

Suggested priority when constructing model context:

1. repo `REVIEW.md`
2. retrieved official standards and language/tool docs
3. prior repo findings or human review comments
4. seed notes from this manifest

Do not delete legacy `source_kind=review_knowledge` rows until v2 passes the
eval gate and completes the 7-day soak.

This keeps the index useful without making imported guidance more authoritative
than repo-local policy or deterministic gx findings.
