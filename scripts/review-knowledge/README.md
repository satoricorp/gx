# Review Knowledge Index

This directory seeds a TurboPuffer corpus for review guidance that can later be
queried by `gx review --deep` and GitHub PR review generation.

The corpus is retrieval context, not a replacement for GX heuristics. Keep the
deterministic review rules in code. Use this index for official standards,
language guidance, tool guidance, and repo-independent review knowledge that a
model can cite or use as supporting context.

## Files

- `urls.json`: source URLs plus GX-specific seed notes from the initial review
  rules discussion.
- `index_review_resources.py`: fetches, extracts, chunks, embeds, and upserts
  seed notes plus fetched source text into TurboPuffer.

## Source Policy

Prefer primary or official sources:

- standards: NIST SSDF, OWASP ASVS/Cheat Sheets, OpenSSF Scorecard, SLSA
- large-scale practice: Google, Microsoft, GitLab, Mozilla
- language docs: TypeScript, Python, Go, Rust
- tool docs: typescript-eslint, Ruff, mypy, Bandit, govulncheck, Clippy, Miri,
  RustSec, cargo-deny

Do not add blog posts or generated summaries as high-authority sources unless
there is no primary source. If a seed note is opinionated, keep it in
`seed_notes` and keep the `authority`/`evidence_level` honest.

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
export GX_REVIEW_KNOWLEDGE_NAMESPACE=gx-review-knowledge
export GX_REVIEW_KNOWLEDGE_CHUNK_CHARS=2400
export GX_REVIEW_KNOWLEDGE_OVERLAP_CHARS=350
```

The default namespace is intentionally separate from `GX_TPUF_NAMESPACE`
(`gx-sessions`) so review knowledge does not mix with session transcripts and
repository code chunks.

## Usage

Validate the manifest without network or API calls:

```bash
python3 scripts/review-knowledge/index_review_resources.py --dry-run
```

Index only the curated seed notes from `urls.json`:

```bash
python3 scripts/review-knowledge/index_review_resources.py --seed-only
```

Fetch source pages and index both seed notes and fetched source chunks:

```bash
python3 scripts/review-knowledge/index_review_resources.py
```

Smoke test a small subset:

```bash
python3 scripts/review-knowledge/index_review_resources.py --limit 3 --seed-only
```

## Retrieval Shape

Future review retrieval should query this namespace using signals from the diff:

- language and framework
- changed file category
- risk profile: auth, data migration, dependency, concurrency, unsafe Rust,
  package release, public API
- repo `REVIEW.md` policy text and fetched URL summaries

Suggested priority when constructing model context:

1. repo `REVIEW.md` and fetched references
2. retrieved official standards and language/tool docs
3. prior repo findings or human review comments
4. seed notes from this manifest

This keeps the index useful without making imported guidance more authoritative
than repo-local policy or deterministic GX findings.
