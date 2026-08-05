# Review Knowledge Index

This directory builds the versioned TurboPuffer corpora used by `lgtm review
--deep` and GitHub PR review generation.

The corpus is retrieval context, not a replacement for lgtm heuristics. Keep the
deterministic review rules in code. Use this index for official standards,
language guidance, tool guidance, and repo-independent review knowledge that a
model can cite or use as supporting context.

## Files

- `sources.yaml`: approved v2 source manifest. This supersedes `urls.json`.
- `golden_queries.yaml`: 50-query eval gate manifest.
- `index_review_resources.py`: validates, fetches, normalizes, chunks, embeds,
  indexes, refreshes, and scaffolds eval output.
- `urls.json`: retained only as the v1 manifest/reference.

## Source Policy

Prefer primary or official sources:

- standards: NIST SSDF, OWASP ASVS/Cheat Sheets, OpenSSF Scorecard, SLSA
- secure-coding bodies: SEI CERT, MISRA, ISO/IEC secure-coding references
- large-scale practice: Google, Microsoft, GitLab, Mozilla
- language docs: TypeScript, JavaScript, Python, Go, Rust, Java, C, C++,
  SQL, C#, shell, PHP, Kotlin, Swift
- framework and vendor docs: Spring, .NET, PostgreSQL, MySQL, SQL Server,
  dbt, Android, Apple platforms
- tool docs: typescript-eslint, Ruff, mypy, Bandit, govulncheck, Clippy, Miri,
  RustSec, cargo-deny, Error Prone, SpotBugs, Checkstyle, PMD, clang-tidy,
  Cppcheck, SQLFluff, Roslyn analyzers, ShellCheck, PHPStan, detekt, SwiftLint

Do not add blog posts or generated summaries as high-authority sources unless
there is no primary source. If a seed note is opinionated, keep it in
`seed_notes` and keep the `authority`/`evidence_level` honest.

## Current Coverage

The v2 manifest covers review process, cross-cutting security, TypeScript,
Python, Go, Rust, SQL, shell, style/linter sources, and empirical research.
Research sources are routed to `research-corpus-v1` and are not used for
review-time retrieval.

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
export LGTM_OPENAI_BASE_URL=https://api.openai.com
export LGTM_OPENAI_EMBEDDING_MODEL=text-embedding-3-small
export LGTM_EMBEDDING_DIMENSIONS=512
export LGTM_TPUF_BASE_URL=https://gcp-us-central1.turbopuffer.com
export LGTM_REVIEW_CANDIDATE_NAMESPACE=review-corpus-v2
export LGTM_REVIEW_KNOWLEDGE_NAMESPACE=lgtm-review-knowledge
export LGTM_RESEARCH_CORPUS_NAMESPACE=research-corpus-v1
export LGTM_REVIEW_RESOURCES=1
export LGTM_REVIEW_RESOURCES_TOP_K=8
```

The default namespace is intentionally separate from `LGTM_TPUF_NAMESPACE`
(`lgtm-sessions`) so review knowledge does not mix with session transcripts and
repository code chunks.

`lgtm review` queries this namespace when both `OPENAI_API_KEY` and
`TURBOPUFFER_API_KEY` are available. Set `LGTM_REVIEW_RESOURCES=0` to disable
review-resource retrieval for a run. `LGTM_REVIEW_RESOURCES_TOP_K` controls the
shallow retrieval limit; the default is 8 so repo-local policy files can remain
in the model context alongside review resources. `lgtm review --deep` raises the
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
`OPENAI_API_KEY` and `TURBOPUFFER_API_KEY` it compares `lgtm-review-knowledge`
against `review-corpus-v2` and writes recall/precedence results:

```bash
python3 scripts/review-knowledge/index_review_resources.py eval
```

After the eval gate passes, promote the same v2 review chunks into the
production namespace (`lgtm-review-knowledge`) so production can cut over without
changing namespace names:

```bash
python3 scripts/review-knowledge/index_review_resources.py promote
```

`promote` upserts `source_kind=review_corpus` rows into `lgtm-review-knowledge`.
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
than repo-local policy or deterministic lgtm findings.
