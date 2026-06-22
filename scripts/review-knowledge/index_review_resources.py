#!/usr/bin/env python3
"""Index review guidance URLs into TurboPuffer.

This is intentionally a script instead of GX runtime code:
- the corpus is slow-changing review knowledge, not a deterministic heuristic;
- it can be run by operators when source URLs or seed notes change;
- it keeps the review runtime free to query the corpus later without owning fetch
  and indexing side effects.
"""

from __future__ import annotations

import argparse
import hashlib
import html.parser
import json
import os
import re
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from dataclasses import dataclass
from typing import Any


DEFAULT_MANIFEST = "scripts/review-knowledge/urls.json"
DEFAULT_OPENAI_BASE_URL = "https://api.openai.com"
DEFAULT_EMBED_MODEL = "text-embedding-3-small"
DEFAULT_EMBED_DIMS = 512
DEFAULT_TPUF_BASE_URL = "https://gcp-us-central1.turbopuffer.com"
DEFAULT_NAMESPACE = "gx-review-knowledge"
DEFAULT_CHUNK_CHARS = 2400
DEFAULT_OVERLAP_CHARS = 350
DEFAULT_BATCH_SIZE = 32


@dataclass(frozen=True)
class Document:
    id: str
    title: str
    url: str
    category: str
    authority: str
    evidence_level: str
    languages: tuple[str, ...]
    frameworks: tuple[str, ...]
    risk_tags: tuple[str, ...]
    review_tags: tuple[str, ...]
    seed_notes: str


@dataclass(frozen=True)
class Chunk:
    id: str
    text: str
    attributes: dict[str, Any]


class TextExtractor(html.parser.HTMLParser):
    """Small HTML-to-text extractor good enough for docs pages.

    We avoid BeautifulSoup/readability dependencies so the script works from a
    fresh checkout. This intentionally strips scripts/styles/nav-like markup and
    keeps headings, paragraphs, lists, table cells, and code blocks as text.
    """

    block_tags = {
        "address",
        "article",
        "aside",
        "blockquote",
        "br",
        "dd",
        "div",
        "dl",
        "dt",
        "figcaption",
        "footer",
        "h1",
        "h2",
        "h3",
        "h4",
        "h5",
        "h6",
        "header",
        "li",
        "main",
        "nav",
        "ol",
        "p",
        "pre",
        "section",
        "table",
        "td",
        "th",
        "tr",
        "ul",
    }

    skip_tags = {"script", "style", "svg", "canvas", "noscript"}

    def __init__(self) -> None:
        super().__init__(convert_charrefs=True)
        self.parts: list[str] = []
        self.skip_depth = 0

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        del attrs
        tag = tag.lower()
        if tag in self.skip_tags:
            self.skip_depth += 1
        elif tag in self.block_tags:
            self.parts.append("\n")

    def handle_endtag(self, tag: str) -> None:
        tag = tag.lower()
        if tag in self.skip_tags and self.skip_depth > 0:
            self.skip_depth -= 1
        elif tag in self.block_tags:
            self.parts.append("\n")

    def handle_data(self, data: str) -> None:
        if self.skip_depth == 0:
            self.parts.append(data)

    def text(self) -> str:
        return normalize_text("".join(self.parts))


def main() -> int:
    parser = argparse.ArgumentParser(description="Index GX review knowledge URLs into TurboPuffer.")
    parser.add_argument("--manifest", default=DEFAULT_MANIFEST, help="Path to urls.json")
    parser.add_argument("--namespace", default=os.getenv("GX_REVIEW_KNOWLEDGE_NAMESPACE", DEFAULT_NAMESPACE))
    parser.add_argument("--limit", type=int, default=0, help="Limit documents for smoke tests")
    parser.add_argument("--chunk-chars", type=int, default=env_int("GX_REVIEW_KNOWLEDGE_CHUNK_CHARS", DEFAULT_CHUNK_CHARS))
    parser.add_argument("--overlap-chars", type=int, default=env_int("GX_REVIEW_KNOWLEDGE_OVERLAP_CHARS", DEFAULT_OVERLAP_CHARS))
    parser.add_argument("--batch-size", type=int, default=env_int("GX_SEMANTIC_BATCH_SIZE", DEFAULT_BATCH_SIZE))
    parser.add_argument("--timeout", type=float, default=12.0, help="Fetch timeout per URL in seconds")
    parser.add_argument("--seed-only", action="store_true", help="Index manifest seed notes without fetching URLs")
    parser.add_argument("--dry-run", action="store_true", help="Validate manifest and print planned work without network/API calls")
    args = parser.parse_args()

    docs = load_manifest(args.manifest)
    if args.limit > 0:
        docs = docs[: args.limit]

    if args.dry_run:
        print_plan(args.manifest, args.namespace, docs)
        return 0

    openai_key = os.getenv("OPENAI_API_KEY", "").strip()
    tpuf_key = os.getenv("TURBOPUFFER_API_KEY", "").strip()
    if not openai_key:
        print("OPENAI_API_KEY is required", file=sys.stderr)
        return 2
    if not tpuf_key:
        print("TURBOPUFFER_API_KEY is required", file=sys.stderr)
        return 2

    chunks = build_chunks(
        docs,
        seed_only=args.seed_only,
        chunk_chars=args.chunk_chars,
        overlap_chars=args.overlap_chars,
        timeout=args.timeout,
    )
    if not chunks:
        print("No chunks produced.")
        return 0

    embedder = OpenAIEmbedder(
        api_key=openai_key,
        base_url=os.getenv("GX_OPENAI_BASE_URL", DEFAULT_OPENAI_BASE_URL),
        model=os.getenv("GX_OPENAI_EMBEDDING_MODEL", DEFAULT_EMBED_MODEL),
        dimensions=env_int("GX_EMBEDDING_DIMENSIONS", DEFAULT_EMBED_DIMS),
    )
    store = TurboPufferStore(
        api_key=tpuf_key,
        base_url=os.getenv("GX_TPUF_BASE_URL", DEFAULT_TPUF_BASE_URL),
        namespace=args.namespace,
        dimensions=embedder.dimensions,
    )

    indexed = 0
    for batch in batches(chunks, args.batch_size):
        vectors = embedder.embed([chunk.text for chunk in batch])
        rows = []
        for chunk, vector in zip(batch, vectors, strict=True):
            row = {"id": chunk.id, "vector": vector}
            row.update(chunk.attributes)
            rows.append(row)
        store.upsert(rows)
        indexed += len(batch)
        print(f"Indexed {indexed}/{len(chunks)} chunks")

    print(f"Done. namespace={args.namespace} chunks={len(chunks)} docs={len(docs)}")
    return 0


def load_manifest(path: str) -> list[Document]:
    with open(path, "r", encoding="utf-8") as handle:
        payload = json.load(handle)
    docs = [parse_document(item) for item in payload.get("documents", [])]
    if not docs:
        raise SystemExit(f"{path}: no documents found")
    ids = [doc.id for doc in docs]
    duplicates = sorted({doc_id for doc_id in ids if ids.count(doc_id) > 1})
    if duplicates:
        raise SystemExit(f"{path}: duplicate document ids: {', '.join(duplicates)}")
    return docs


def parse_document(item: dict[str, Any]) -> Document:
    required = ["id", "title", "url", "category", "authority", "evidence_level", "seed_notes"]
    missing = [name for name in required if not str(item.get(name, "")).strip()]
    if missing:
        raise SystemExit(f"document missing required fields {missing}: {item!r}")
    parsed = urllib.parse.urlparse(item["url"])
    if parsed.scheme not in {"http", "https"} or not parsed.netloc:
        raise SystemExit(f"invalid URL for {item['id']}: {item['url']}")
    return Document(
        id=clean_token(item["id"]),
        title=str(item["title"]).strip(),
        url=str(item["url"]).strip(),
        category=clean_token(item["category"]),
        authority=clean_token(item["authority"]),
        evidence_level=clean_token(item["evidence_level"]),
        languages=tuple(clean_token(value) for value in item.get("languages", [])),
        frameworks=tuple(clean_token(value) for value in item.get("frameworks", [])),
        risk_tags=tuple(clean_token(value) for value in item.get("risk_tags", [])),
        review_tags=tuple(clean_token(value) for value in item.get("review_tags", [])),
        seed_notes=str(item["seed_notes"]).strip(),
    )


def print_plan(manifest_path: str, namespace: str, docs: list[Document]) -> None:
    categories: dict[str, int] = {}
    authorities: dict[str, int] = {}
    for doc in docs:
        categories[doc.category] = categories.get(doc.category, 0) + 1
        authorities[doc.authority] = authorities.get(doc.authority, 0) + 1
    print(f"manifest={manifest_path}")
    print(f"namespace={namespace}")
    print(f"documents={len(docs)}")
    print("categories=" + ", ".join(f"{key}:{value}" for key, value in sorted(categories.items())))
    print("authorities=" + ", ".join(f"{key}:{value}" for key, value in sorted(authorities.items())))
    for doc in docs:
        print(f"- {doc.id} [{doc.category}/{doc.evidence_level}] {doc.url}")


def build_chunks(
    docs: list[Document],
    *,
    seed_only: bool,
    chunk_chars: int,
    overlap_chars: int,
    timeout: float,
) -> list[Chunk]:
    chunks: list[Chunk] = []
    fetched_at = str(int(time.time()))
    for doc in docs:
        seed_text = render_seed_text(doc)
        chunks.append(make_chunk(doc, "seed", 0, seed_text, fetched_at))
        if seed_only:
            continue
        try:
            fetched_text = fetch_text(doc.url, timeout)
        except Exception as exc:  # noqa: BLE001 - operators need per-URL resilience.
            print(f"warn: {doc.id}: fetch failed: {exc}", file=sys.stderr)
            continue
        for index, text in enumerate(split_text(fetched_text, chunk_chars, overlap_chars)):
            chunks.append(make_chunk(doc, "source", index, text, fetched_at))
    return chunks


def render_seed_text(doc: Document) -> str:
    lines = [
        "GX review knowledge seed",
        f"title: {doc.title}",
        f"url: {doc.url}",
        f"category: {doc.category}",
        f"authority: {doc.authority}",
        f"evidence_level: {doc.evidence_level}",
        f"languages: {join_values(doc.languages)}",
        f"frameworks: {join_values(doc.frameworks)}",
        f"risk_tags: {join_values(doc.risk_tags)}",
        f"review_tags: {join_values(doc.review_tags)}",
        "",
        doc.seed_notes,
    ]
    return normalize_text("\n".join(lines))


def make_chunk(doc: Document, chunk_kind: str, index: int, text: str, fetched_at: str) -> Chunk:
    normalized = normalize_text(text)
    content_hash = sha1(normalized)
    row_id = "gx-review-knowledge-" + sha1(f"{doc.url}:{chunk_kind}:{index}:{content_hash}")[:40]
    attrs = {
        "text": normalized,
        "source_kind": "review_knowledge",
        "source_id": doc.id,
        "url": doc.url,
        "title": doc.title,
        "category": doc.category,
        "authority": doc.authority,
        "evidence_level": doc.evidence_level,
        "languages": join_values(doc.languages),
        "frameworks": join_values(doc.frameworks),
        "risk_tags": join_values(doc.risk_tags),
        "review_tags": join_values(doc.review_tags),
        "chunk_kind": chunk_kind,
        "chunk_index": index,
        "content_hash": content_hash,
        "fetched_at": fetched_at,
    }
    return Chunk(id=row_id, text=normalized, attributes=attrs)


def fetch_text(url: str, timeout: float) -> str:
    request = urllib.request.Request(
        url,
        headers={
            "Accept": "text/html,text/markdown,text/plain;q=0.9,*/*;q=0.1",
            "User-Agent": "gx-review-knowledge-indexer/0.1",
        },
    )
    with urllib.request.urlopen(request, timeout=timeout) as response:
        content_type = response.headers.get("Content-Type", "").lower()
        raw = response.read(1_500_000)
    if not any(kind in content_type for kind in ["text/html", "text/plain", "text/markdown", "application/xhtml+xml", ""]):
        raise ValueError(f"unsupported content type {content_type!r}")
    text = raw.decode("utf-8", errors="replace")
    if "html" in content_type or "<html" in text[:2000].lower():
        extractor = TextExtractor()
        extractor.feed(text)
        text = extractor.text()
    return normalize_text(text)


def split_text(text: str, chunk_chars: int, overlap_chars: int) -> list[str]:
    text = normalize_text(text)
    if not text:
        return []
    if len(text) <= chunk_chars:
        return [text]
    chunks: list[str] = []
    start = 0
    while start < len(text):
        end = min(len(text), start + chunk_chars)
        cut = best_boundary(text, start, end)
        if cut <= start:
            cut = end
        chunk = text[start:cut].strip()
        if chunk:
            chunks.append(chunk)
        if cut >= len(text):
            break
        start = max(cut - overlap_chars, start + 1)
    return chunks


def best_boundary(text: str, start: int, end: int) -> int:
    window = text[start:end]
    for pattern in ["\n\n", "\n", ". ", "? ", "! "]:
        idx = window.rfind(pattern)
        if idx > len(window) * 0.45:
            return start + idx + len(pattern)
    space = window.rfind(" ")
    if space > len(window) * 0.45:
        return start + space
    return end


class OpenAIEmbedder:
    def __init__(self, *, api_key: str, base_url: str, model: str, dimensions: int) -> None:
        self.api_key = api_key
        self.base_url = normalize_openai_base_url(base_url)
        self.model = model.strip() or DEFAULT_EMBED_MODEL
        self.dimensions = dimensions

    def embed(self, inputs: list[str]) -> list[list[float]]:
        payload: dict[str, Any] = {
            "model": self.model,
            "input": inputs,
            "encoding_format": "float",
            "dimensions": self.dimensions,
        }
        data = request_json(
            f"{self.base_url}/embeddings",
            payload,
            headers={"Authorization": f"Bearer {self.api_key}"},
        )
        items = data.get("data", [])
        if len(items) != len(inputs):
            raise RuntimeError(f"OpenAI returned {len(items)} embeddings for {len(inputs)} inputs")
        vectors: list[list[float] | None] = [None] * len(inputs)
        for item in items:
            vectors[int(item["index"])] = item["embedding"]
        if any(vector is None for vector in vectors):
            raise RuntimeError("OpenAI response omitted one or more embedding indexes")
        return [vector for vector in vectors if vector is not None]


class TurboPufferStore:
    def __init__(self, *, api_key: str, base_url: str, namespace: str, dimensions: int) -> None:
        self.api_key = api_key
        self.base_url = base_url.rstrip("/")
        self.namespace = namespace.strip("/")
        self.dimensions = dimensions

    def upsert(self, rows: list[dict[str, Any]]) -> None:
        payload = {
            "distance_metric": "cosine_distance",
            "schema": review_knowledge_schema(self.dimensions),
            "upsert_rows": rows,
        }
        request_json(
            f"{self.base_url}/v2/namespaces/{urllib.parse.quote(self.namespace)}",
            payload,
            headers={"Authorization": f"Bearer {self.api_key}"},
        )


def request_json(url: str, payload: dict[str, Any], *, headers: dict[str, str]) -> dict[str, Any]:
    body = json.dumps(payload).encode("utf-8")
    request_headers = {"Content-Type": "application/json", **headers}
    request = urllib.request.Request(url, data=body, headers=request_headers, method="POST")
    try:
        with urllib.request.urlopen(request, timeout=60) as response:
            return json.loads(response.read().decode("utf-8"))
    except urllib.error.HTTPError as exc:
        detail = exc.read().decode("utf-8", errors="replace")[:1000]
        raise RuntimeError(f"{url} failed with status {exc.code}: {detail}") from exc


def review_knowledge_schema(dimensions: int) -> dict[str, Any]:
    string_filter = {"type": "string", "filterable": True}
    string_search = {"type": "string", "full_text_search": True}
    return {
        "vector": {"type": f"[{dimensions}]f32", "ann": True},
        "text": string_search,
        "source_kind": string_filter,
        "source_id": string_filter,
        "url": string_filter,
        "title": string_search,
        "category": string_filter,
        "authority": string_filter,
        "evidence_level": string_filter,
        "languages": string_search,
        "frameworks": string_search,
        "risk_tags": string_search,
        "review_tags": string_search,
        "chunk_kind": string_filter,
        "chunk_index": {"type": "uint", "filterable": True},
        "content_hash": string_filter,
        "fetched_at": string_filter,
    }


def normalize_openai_base_url(raw: str) -> str:
    base = raw.strip().rstrip("/") or DEFAULT_OPENAI_BASE_URL
    return base if base.endswith("/v1") else f"{base}/v1"


def normalize_text(text: str) -> str:
    text = text.replace("\r\n", "\n").replace("\r", "\n")
    text = re.sub(r"[ \t]+\n", "\n", text)
    text = re.sub(r"\n{3,}", "\n\n", text)
    text = re.sub(r"[ \t]{2,}", " ", text)
    return text.strip()


def clean_token(value: Any) -> str:
    return re.sub(r"[^a-zA-Z0-9_.-]+", "-", str(value).strip().lower()).strip("-")


def join_values(values: tuple[str, ...]) -> str:
    return ",".join(value for value in values if value)


def sha1(text: str) -> str:
    return hashlib.sha1(text.encode("utf-8")).hexdigest()


def env_int(name: str, fallback: int) -> int:
    raw = os.getenv(name, "").strip()
    if not raw:
        return fallback
    try:
        return int(raw)
    except ValueError:
        return fallback


def batches(values: list[Chunk], size: int) -> list[list[Chunk]]:
    size = max(1, size)
    return [values[index : index + size] for index in range(0, len(values), size)]


if __name__ == "__main__":
    raise SystemExit(main())
