#!/usr/bin/env python3
"""Build and index the GX review corpus.

The v2 corpus is source-manifest driven:
- review guidance indexes into review-corpus-v2;
- research papers index into research-corpus-v1;
- the existing gx-review-knowledge namespace is never deleted or overwritten.
"""

from __future__ import annotations

import argparse
import collections
import datetime as dt
import hashlib
import html.parser
import json
import os
import re
import statistics
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
import urllib.robotparser
from dataclasses import dataclass
from pathlib import Path
from typing import Any

try:
    import yaml  # type: ignore[import-not-found]
except Exception:  # pragma: no cover - fallback is exercised only without PyYAML.
    yaml = None


DEFAULT_MANIFEST = "scripts/review-knowledge/sources.yaml"
DEFAULT_ARTIFACT_DIR = "scripts/review-knowledge/corpus"
DEFAULT_OPENAI_BASE_URL = "https://api.openai.com"
DEFAULT_EMBED_MODEL = "text-embedding-3-small"
DEFAULT_EMBED_DIMS = 512
DEFAULT_TPUF_BASE_URL = "https://gcp-us-central1.turbopuffer.com"
DEFAULT_PRODUCTION_NAMESPACE = "gx-review-knowledge"
DEFAULT_REVIEW_NAMESPACE = "review-corpus-v2"
DEFAULT_RESEARCH_NAMESPACE = "research-corpus-v1"
DEFAULT_BATCH_SIZE = 256
TARGET_MIN_TOKENS = 300
TARGET_MAX_TOKENS = 800


@dataclass(frozen=True)
class Source:
    id: str
    url: str | None
    tier: str
    languages: tuple[str, ...]
    authority: int
    precedence_group: str | None
    superseded_by: str | None
    historical: bool
    license: str
    fetch: str
    crawl_depth: int
    notes: str


@dataclass(frozen=True)
class FetchedDoc:
    source: Source
    text: str
    raw: bytes
    final_url: str | None
    http_status: int | None
    fetched_at: str
    error: str | None = None


@dataclass(frozen=True)
class Chunk:
    id: str
    body: str
    namespace: str
    attributes: dict[str, Any]


class MarkdownExtractor(html.parser.HTMLParser):
    block_tags = {
        "article",
        "blockquote",
        "br",
        "dd",
        "div",
        "dl",
        "dt",
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
        self.heading_depth: int | None = None

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        del attrs
        tag = tag.lower()
        if tag in self.skip_tags:
            self.skip_depth += 1
            return
        if re.fullmatch(r"h[1-6]", tag):
            self.heading_depth = int(tag[1])
            self.parts.append("\n\n" + ("#" * self.heading_depth) + " ")
        elif tag == "li":
            self.parts.append("\n- ")
        elif tag in self.block_tags:
            self.parts.append("\n")

    def handle_endtag(self, tag: str) -> None:
        tag = tag.lower()
        if tag in self.skip_tags and self.skip_depth > 0:
            self.skip_depth -= 1
            return
        if re.fullmatch(r"h[1-6]", tag):
            self.heading_depth = None
            self.parts.append("\n\n")
        elif tag in self.block_tags:
            self.parts.append("\n")

    def handle_data(self, data: str) -> None:
        if self.skip_depth == 0:
            self.parts.append(data)

    def markdown(self) -> str:
        return normalize_text("".join(self.parts))


class LinkExtractor(html.parser.HTMLParser):
    def __init__(self) -> None:
        super().__init__(convert_charrefs=True)
        self.links: list[str] = []

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        if tag.lower() != "a":
            return
        for key, value in attrs:
            if key.lower() == "href" and value:
                self.links.append(value)


def main() -> int:
    parser = argparse.ArgumentParser(description="Build and index GX review corpus resources.")
    parser.add_argument("command", nargs="?", choices=["dry-run", "fetch", "chunk", "index", "eval", "promote", "refresh"], help="Pipeline command")
    parser.add_argument("--manifest", default=DEFAULT_MANIFEST)
    parser.add_argument("--artifact-dir", default=DEFAULT_ARTIFACT_DIR)
    parser.add_argument("--review-namespace", default=os.getenv("GX_REVIEW_CANDIDATE_NAMESPACE", DEFAULT_REVIEW_NAMESPACE))
    parser.add_argument("--baseline-namespace", default=os.getenv("GX_REVIEW_BASELINE_NAMESPACE", DEFAULT_PRODUCTION_NAMESPACE))
    parser.add_argument("--production-namespace", default=os.getenv("GX_REVIEW_KNOWLEDGE_NAMESPACE", DEFAULT_PRODUCTION_NAMESPACE))
    parser.add_argument("--research-namespace", default=os.getenv("GX_RESEARCH_CORPUS_NAMESPACE", DEFAULT_RESEARCH_NAMESPACE))
    parser.add_argument("--limit", type=int, default=0, help="Limit sources for smoke tests")
    parser.add_argument("--timeout", type=float, default=20.0)
    parser.add_argument("--batch-size", type=int, default=env_int("GX_SEMANTIC_BATCH_SIZE", DEFAULT_BATCH_SIZE))
    parser.add_argument("--seed-only", action="store_true", help="Chunk/index manifest seed notes only")
    parser.add_argument("--dry-run", action="store_true", help="Validate manifest without network/API calls")
    parser.add_argument("--golden", default="scripts/review-knowledge/golden_queries.yaml")
    args = parser.parse_args()

    command = args.command or ("dry-run" if args.dry_run else "index")
    sources = load_sources(args.manifest)
    if args.limit > 0:
        sources = sources[: args.limit]

    if command == "dry-run":
        print_plan(args.manifest, sources, args.review_namespace, args.research_namespace)
        return 0
    if command == "fetch":
        fetch_sources(sources, Path(args.artifact_dir), timeout=args.timeout)
        return 0
    if command == "chunk":
        chunks = build_chunks(sources, Path(args.artifact_dir), seed_only=args.seed_only)
        write_chunks(chunks, Path(args.artifact_dir))
        write_build_log(chunks, sources, Path(args.artifact_dir))
        return 0
    if command == "eval":
        return run_eval(args.golden, args.baseline_namespace, args.review_namespace, Path(args.artifact_dir))
    if command == "promote":
        if args.production_namespace == args.review_namespace:
            print("production namespace already equals candidate namespace; use `index` instead", file=sys.stderr)
            return 2
        if not args.seed_only and not any((Path(args.artifact_dir) / "normalized").glob("*.md")):
            print("No normalized corpus found; run `fetch` first or pass --seed-only.", file=sys.stderr)
            return 2
        args.review_namespace = args.production_namespace
        chunks = [chunk for chunk in build_chunks(sources, Path(args.artifact_dir), seed_only=args.seed_only) if chunk.attributes["tier"] != "research"]
        write_chunks(chunks, Path(args.artifact_dir))
        write_build_log(chunks, sources, Path(args.artifact_dir))
        return index_chunks(chunks, args)
    if command == "refresh":
        fetch_sources(sources, Path(args.artifact_dir), timeout=args.timeout)
        chunks = build_chunks(sources, Path(args.artifact_dir), seed_only=False)
        write_chunks(chunks, Path(args.artifact_dir))
        write_build_log(chunks, sources, Path(args.artifact_dir))
        return index_chunks(chunks, args)

    if not args.seed_only and not any((Path(args.artifact_dir) / "normalized").glob("*.md")):
        print("No normalized corpus found; run `fetch` first or pass --seed-only.", file=sys.stderr)
        return 2
    chunks = build_chunks(sources, Path(args.artifact_dir), seed_only=args.seed_only)
    write_chunks(chunks, Path(args.artifact_dir))
    write_build_log(chunks, sources, Path(args.artifact_dir))
    return index_chunks(chunks, args)


def load_sources(path: str) -> list[Source]:
    manifest_path = Path(path)
    with manifest_path.open("r", encoding="utf-8") as handle:
        if manifest_path.suffix == ".json":
            payload = json.load(handle)
            raw_sources = legacy_json_sources(payload)
        else:
            raw_sources = parse_yaml_manifest(handle.read()).get("sources", [])
    sources = [parse_source(item) for item in raw_sources]
    if not sources:
        raise SystemExit(f"{path}: no sources found")
    duplicates = sorted({source.id for source in sources if [s.id for s in sources].count(source.id) > 1})
    if duplicates:
        raise SystemExit(f"{path}: duplicate source ids: {', '.join(duplicates)}")
    return sources


def parse_yaml_manifest(text: str) -> dict[str, Any]:
    if yaml is not None:
        loaded = yaml.safe_load(text)
        return loaded if isinstance(loaded, dict) else {}
    raise SystemExit("PyYAML is required to parse sources.yaml")


def legacy_json_sources(payload: dict[str, Any]) -> list[dict[str, Any]]:
    out: list[dict[str, Any]] = []
    category_to_tier = {
        "core-process": "process",
        "supply-chain": "security",
        "database": "language",
        "framework": "language",
        "tooling": "style",
        "safety-critical": "security",
    }
    for item in payload.get("documents", []):
        category = clean_token(item.get("category", "language"))
        out.append(
            {
                "id": item.get("id"),
                "url": item.get("url"),
                "tier": category_to_tier.get(category, category),
                "languages": item.get("languages", []),
                "authority": 2,
                "fetch": "html",
                "notes": item.get("seed_notes", ""),
            }
        )
    return out


def parse_source(item: dict[str, Any]) -> Source:
    source_id = clean_token(item.get("id", ""))
    if not source_id:
        raise SystemExit(f"source missing id: {item!r}")
    tier = clean_token(item.get("tier", ""))
    if tier not in {"process", "security", "language", "style", "research"}:
        raise SystemExit(f"{source_id}: invalid tier {tier!r}")
    fetch = clean_token(item.get("fetch", "html"))
    if fetch not in {"html", "pdf", "github-md", "manual"}:
        raise SystemExit(f"{source_id}: invalid fetch {fetch!r}")
    url_value = item.get("url")
    url = str(url_value).strip() if url_value not in {None, ""} else None
    if fetch != "manual":
        if not url:
            raise SystemExit(f"{source_id}: url is required unless fetch=manual")
        parsed = urllib.parse.urlparse(url)
        if parsed.scheme not in {"http", "https"} or not parsed.netloc:
            raise SystemExit(f"{source_id}: invalid URL {url!r}")
    return Source(
        id=source_id,
        url=url,
        tier=tier,
        languages=tuple(clean_token(value) for value in item.get("languages", []) if clean_token(value)),
        authority=int(item.get("authority", 2) or 2),
        precedence_group=nullable_token(item.get("precedence_group")),
        superseded_by=nullable_token(item.get("superseded_by")),
        historical=bool(item.get("historical", False)),
        license=str(item.get("license", "verify") or "verify").strip(),
        fetch=fetch,
        crawl_depth=int(item.get("crawl_depth", 0) or 0),
        notes=str(item.get("notes", "") or "").strip(),
    )


def print_plan(manifest_path: str, sources: list[Source], review_namespace: str, research_namespace: str) -> None:
    by_tier = collections.Counter(source.tier for source in sources)
    by_fetch = collections.Counter(source.fetch for source in sources)
    print(f"manifest={manifest_path}")
    print(f"review_namespace={review_namespace}")
    print(f"research_namespace={research_namespace}")
    print(f"sources={len(sources)}")
    print("tiers=" + ", ".join(f"{key}:{by_tier[key]}" for key in sorted(by_tier)))
    print("fetch=" + ", ".join(f"{key}:{by_fetch[key]}" for key in sorted(by_fetch)))
    for source in sources:
        namespace = research_namespace if source.tier == "research" else review_namespace
        print(f"- {source.id} [{source.tier}; authority={source.authority}; fetch={source.fetch}; namespace={namespace}] {source.url or 'manual'}")


def fetch_sources(sources: list[Source], artifact_dir: Path, *, timeout: float) -> None:
    for dirname in ["raw", "normalized", "metadata"]:
        (artifact_dir / dirname).mkdir(parents=True, exist_ok=True)
    robots: dict[str, urllib.robotparser.RobotFileParser] = {}
    last_by_domain: dict[str, float] = {}
    report: list[str] = ["# Fetch Report", "", f"generated_at: {utc_today()}", ""]
    for source in sources:
        if source.fetch == "manual":
            write_meta(artifact_dir, source, None, None, None, "manual source pending local drop")
            report.append(f"- SKIP {source.id}: manual source pending. license={source.license}")
            continue
        if source.license in {"verify", "unknown"}:
            report.append(f"- LICENSE {source.id}: license={source.license}; verify before redistribution.")
        try:
            fetched = fetch_source(source, robots, last_by_domain, timeout)
        except Exception as exc:  # noqa: BLE001 - report every source failure.
            write_meta(artifact_dir, source, None, None, None, str(exc))
            report.append(f"- FAIL {source.id}: {exc}")
            continue
        raw_suffix = ".pdf" if source.fetch == "pdf" else ".txt"
        (artifact_dir / "raw" / f"{source.id}{raw_suffix}").write_bytes(fetched.raw)
        (artifact_dir / "normalized" / f"{source.id}.md").write_text(fetched.text + "\n", encoding="utf-8")
        write_meta(artifact_dir, source, fetched.final_url, fetched.http_status, sha256_bytes(fetched.raw), fetched.error)
        report.append(f"- OK {source.id}: status={fetched.http_status} final_url={fetched.final_url}")
    (artifact_dir / "fetch_report.md").write_text("\n".join(report) + "\n", encoding="utf-8")
    print(f"Wrote {artifact_dir / 'fetch_report.md'}")


def fetch_source(
    source: Source,
    robots: dict[str, urllib.robotparser.RobotFileParser],
    last_by_domain: dict[str, float],
    timeout: float,
) -> FetchedDoc:
    assert source.url is not None
    queue: list[tuple[str, int]] = [(source.url, 0)]
    seen_urls = {source.url}
    text_parts: list[str] = []
    raw_parts: list[bytes] = []
    final_url = source.url
    status: int | None = None
    while queue:
        url, depth = queue.pop(0)
        if not robots_allowed(url, robots):
            raise RuntimeError(f"robots.txt disallows fetch: {url}")
        throttle(url, last_by_domain)
        raw, final_url, status, content_type = fetch_url(resolve_fetch_url(url, source.fetch), timeout)
        if source.fetch == "pdf" and "application/pdf" not in content_type.lower() and not final_url.lower().endswith(".pdf"):
            pdf_links = [link for link in extract_links(raw, final_url) if link.lower().endswith(".pdf")]
            if pdf_links:
                throttle(pdf_links[0], last_by_domain)
                raw, final_url, status, content_type = fetch_url(pdf_links[0], timeout)
        raw_parts.append(raw)
        text_parts.append(normalize_payload(raw, content_type, source.fetch, final_url))
        if source.fetch == "html" and depth < source.crawl_depth:
            for link in crawl_links(raw, final_url, source.crawl_depth - depth):
                if link in seen_urls:
                    continue
                seen_urls.add(link)
                queue.append((link, depth + 1))
                if len(seen_urls) >= 50:
                    break
    return FetchedDoc(
        source=source,
        text=normalize_text("\n\n".join(text_parts)),
        raw=b"\n\n".join(raw_parts),
        final_url=final_url,
        http_status=status,
        fetched_at=utc_today(),
    )


def robots_allowed(url: str, robots: dict[str, urllib.robotparser.RobotFileParser]) -> bool:
    parsed = urllib.parse.urlparse(url)
    root = f"{parsed.scheme}://{parsed.netloc}"
    parser = robots.get(root)
    if parser is None:
        parser = urllib.robotparser.RobotFileParser()
        parser.set_url(root + "/robots.txt")
        try:
            parser.read()
        except Exception:
            return True
        robots[root] = parser
    return parser.can_fetch("gx-review-corpus-indexer/0.2", url)


def throttle(url: str, last_by_domain: dict[str, float]) -> None:
    domain = urllib.parse.urlparse(url).netloc.lower()
    last = last_by_domain.get(domain, 0)
    delay = 1.0 - (time.time() - last)
    if delay > 0:
        time.sleep(delay)
    last_by_domain[domain] = time.time()


def fetch_url(url: str, timeout: float) -> tuple[bytes, str, int, str]:
    headers = {
        "Accept": "text/html,text/markdown,text/plain,application/pdf;q=0.9,*/*;q=0.1",
        "User-Agent": "gx-review-corpus-indexer/0.2",
    }
    last_error: Exception | None = None
    for attempt in range(3):
        request = urllib.request.Request(url, headers=headers)
        try:
            with urllib.request.urlopen(request, timeout=timeout) as response:
                return response.read(3_000_000), response.geturl(), response.status, response.headers.get("Content-Type", "")
        except urllib.error.HTTPError as exc:
            if exc.code in {401, 402, 403, 404, 451}:
                raise RuntimeError(f"http {exc.code} for {url}") from exc
            last_error = exc
        except (urllib.error.URLError, TimeoutError, ConnectionError) as exc:
            last_error = exc
        if attempt < 2:
            time.sleep(2**attempt)
    raise RuntimeError(f"fetch failed for {url}: {last_error}")


def resolve_fetch_url(url: str, fetch: str) -> str:
    parsed = urllib.parse.urlparse(url)
    if fetch == "github-md" and parsed.netloc == "github.com":
        parts = [part for part in parsed.path.split("/") if part]
        if len(parts) >= 2:
            owner, repo = parts[0], parts[1]
            if len(parts) >= 5 and parts[2] in {"blob", "tree"}:
                branch = parts[3]
                rest = "/".join(parts[4:])
                return f"https://raw.githubusercontent.com/{owner}/{repo}/{branch}/{rest}"
            return f"https://raw.githubusercontent.com/{owner}/{repo}/master/README.md"
    if fetch == "pdf" and parsed.netloc == "arxiv.org" and parsed.path.startswith("/abs/"):
        return "https://arxiv.org/pdf/" + parsed.path.removeprefix("/abs/")
    return url


def normalize_payload(raw: bytes, content_type: str, fetch: str, url: str) -> str:
    if fetch == "pdf" or "application/pdf" in content_type.lower() or url.lower().endswith(".pdf"):
        return pdf_to_text(raw)
    text = raw.decode("utf-8", errors="replace")
    if fetch == "github-md" or "markdown" in content_type.lower():
        return normalize_text(text)
    try:
        import trafilatura  # type: ignore[import-not-found]

        extracted = trafilatura.extract(text, output_format="markdown", include_comments=False, include_tables=True)
        if extracted:
            return normalize_text(extracted)
    except Exception:
        pass
    extractor = MarkdownExtractor()
    extractor.feed(text)
    return extractor.markdown()


def pdf_to_text(raw: bytes) -> str:
    try:
        import fitz  # type: ignore[import-not-found]
    except Exception as exc:
        raise RuntimeError("PyMuPDF is required for PDF normalization") from exc
    parts: list[str] = []
    with fitz.open(stream=raw, filetype="pdf") as doc:
        for page in doc:
            parts.append(page.get_text("text"))
    return normalize_text("\n\n".join(parts))


def crawl_links(raw: bytes, base_url: str, depth: int) -> list[str]:
    if depth <= 0:
        return []
    parsed_base = urllib.parse.urlparse(base_url)
    base_path = parsed_base.path.rstrip("/")
    out: list[str] = []
    seen: set[str] = set()
    for normalized in extract_links(raw, base_url):
        parsed = urllib.parse.urlparse(normalized)
        if parsed.netloc != parsed_base.netloc or not parsed.path.startswith(base_path):
            continue
        if normalized not in seen:
            seen.add(normalized)
            out.append(normalized)
    return out


def extract_links(raw: bytes, base_url: str) -> list[str]:
    extractor = LinkExtractor()
    extractor.feed(raw.decode("utf-8", errors="replace"))
    out: list[str] = []
    seen: set[str] = set()
    for href in extractor.links:
        absolute = urllib.parse.urljoin(base_url, href)
        parsed = urllib.parse.urlparse(absolute)
        normalized = parsed._replace(fragment="", query="").geturl()
        if normalized not in seen:
            seen.add(normalized)
            out.append(normalized)
    return out


def write_meta(artifact_dir: Path, source: Source, final_url: str | None, status: int | None, content_hash: str | None, error: str | None) -> None:
    payload = {
        "source_id": source.id,
        "sha256": content_hash,
        "fetched_at": utc_today(),
        "final_url": final_url,
        "http_status": status,
        "license": source.license,
        "license_status": "needs-verification" if source.license in {"verify", "unknown"} else "recorded",
        "error": error,
    }
    (artifact_dir / "metadata").mkdir(parents=True, exist_ok=True)
    (artifact_dir / "metadata" / f"{source.id}.meta.json").write_text(json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def build_chunks(sources: list[Source], artifact_dir: Path, *, seed_only: bool) -> list[Chunk]:
    chunks: list[Chunk] = []
    for source in sources:
        meta = load_meta(artifact_dir, source.id)
        if source.fetch == "manual" and not seed_only:
            continue
        text = seed_text(source) if seed_only else read_normalized(artifact_dir, source) or seed_text(source)
        for index, (section_path, body) in enumerate(split_markdown(text)):
            chunk = make_chunk(source, section_path, index, body, meta)
            chunks.append(chunk)
    return chunks


def read_normalized(artifact_dir: Path, source: Source) -> str:
    path = artifact_dir / "normalized" / f"{source.id}.md"
    if not path.exists():
        return ""
    return path.read_text(encoding="utf-8")


def load_meta(artifact_dir: Path, source_id: str) -> dict[str, Any]:
    path = artifact_dir / "metadata" / f"{source_id}.meta.json"
    if not path.exists():
        return {}
    return json.loads(path.read_text(encoding="utf-8"))


def seed_text(source: Source) -> str:
    title = title_from_source(source)
    return normalize_text(
        "\n".join(
            [
                f"# {title}",
                f"source_id: {source.id}",
                f"tier: {source.tier}",
                f"url: {source.url or 'manual'}",
                "",
                source.notes,
            ]
        )
    )


def split_markdown(text: str) -> list[tuple[str, str]]:
    blocks: list[tuple[list[str], list[str]]] = []
    current_path: list[str] = []
    current_lines: list[str] = []
    for line in text.splitlines():
        heading = re.match(r"^(#{1,6})\s+(.+?)\s*$", line)
        if heading:
            if current_lines:
                blocks.append((current_path[:], current_lines))
                current_lines = []
            level = len(heading.group(1))
            title = re.sub(r"\s+#*$", "", heading.group(2)).strip()
            current_path = current_path[: level - 1] + [title]
        current_lines.append(line)
    if current_lines:
        blocks.append((current_path[:], current_lines))

    chunks: list[tuple[str, str]] = []
    for path, lines in blocks:
        section = " > ".join(path) if path else "Document"
        body = normalize_text("\n".join(lines))
        if not body:
            continue
        chunks.extend(split_section(section, body))
    return chunks


def split_section(section: str, body: str) -> list[tuple[str, str]]:
    tokens = rough_token_count(body)
    if tokens <= TARGET_MAX_TOKENS:
        return [(section, body)]
    paragraphs = re.split(r"\n(?=\s*(?:[-*]\s+|\d+\.\s+|#{1,6}\s+|\S))", body)
    out: list[tuple[str, str]] = []
    current: list[str] = []
    for paragraph in paragraphs:
        proposed = normalize_text("\n".join(current + [paragraph]))
        if current and rough_token_count(proposed) > TARGET_MAX_TOKENS:
            out.append((section, normalize_text("\n".join(current))))
            current = [paragraph]
        else:
            current.append(paragraph)
    if current:
        out.append((section, normalize_text("\n".join(current))))
    return out


def make_chunk(source: Source, section_path: str, index: int, body: str, meta: dict[str, Any]) -> Chunk:
    title = title_from_source(source)
    breadcrumb = f"Source: {title} > {section_path}"
    chunk_body = normalize_text(breadcrumb + "\n\n" + body)
    section_slug = slugify(section_path) or "document"
    max_slug_bytes = max(8, 64 - len(source.id.encode("utf-8")) - len("::::0000"))
    section_slug = truncate_ascii_token(section_slug, max_slug_bytes)
    chunk_id = f"{source.id}::{section_slug}::{index:04d}"
    superseded_by = chunk_superseded_by(source, section_path, chunk_body)
    namespace = DEFAULT_RESEARCH_NAMESPACE if source.tier == "research" else DEFAULT_REVIEW_NAMESPACE
    attrs = {
        "body": chunk_body,
        "text": chunk_body,
        "source_kind": "review_corpus" if source.tier != "research" else "research_corpus",
        "chunk_id": chunk_id,
        "source_id": source.id,
        "source_url": source.url or "",
        "url": source.url or "",
        "title": title,
        "tier": source.tier,
        "category": source.tier,
        "languages": list(source.languages),
        "language_tags": list(source.languages),
        "authority": source.authority,
        "precedence_group": source.precedence_group or "",
        "superseded_by": superseded_by or source.superseded_by or "",
        "historical": source.historical,
        "section_path": section_path,
        "cwe_ids": extract_cwe_ids(chunk_body),
        "content_sha256": sha256_text(chunk_body),
        "content_hash": sha256_text(chunk_body),
        "fetched_at": str(meta.get("fetched_at") or utc_today()),
        "usefulness_rank": usefulness_rank(source, superseded_by),
        "chunk_index": index,
    }
    return Chunk(id=chunk_id, body=chunk_body, namespace=namespace, attributes=attrs)


def chunk_superseded_by(source: Source, section_path: str, body: str) -> str:
    lower = (section_path + "\n" + body).lower()
    if source.id == "pep8" and any(term in lower for term in ["line length", "formatting", "indentation", "whitespace"]):
        return "ruff"
    if source.id == "ruff" and any(code in body for code in ["D203", "D213"]):
        return "ruff-google-convention"
    if source.id == "effective-go" and any(term in lower for term in ["interface{}", "pre-generics", "generics"]):
        return "google-go-style"
    return ""


def usefulness_rank(source: Source, superseded_by: str) -> int:
    if superseded_by:
        return 90
    if source.historical:
        return 80
    if source.tier == "research":
        return 70
    return max(1, min(50, source.authority * 10))


def write_chunks(chunks: list[Chunk], artifact_dir: Path) -> None:
    (artifact_dir / "chunks").mkdir(parents=True, exist_ok=True)
    grouped: dict[str, list[Chunk]] = collections.defaultdict(list)
    for chunk in chunks:
        grouped[chunk.namespace].append(chunk)
    for namespace, rows in grouped.items():
        path = artifact_dir / "chunks" / f"{namespace}.jsonl"
        with path.open("w", encoding="utf-8") as handle:
            for chunk in rows:
                handle.write(json.dumps({"id": chunk.id, **chunk.attributes}, sort_keys=True) + "\n")
        print(f"Wrote {path} ({len(rows)} chunks)")


def write_build_log(chunks: list[Chunk], sources: list[Source], artifact_dir: Path) -> None:
    counts = collections.Counter(chunk.attributes["source_id"] for chunk in chunks)
    tokens = [rough_token_count(chunk.body) for chunk in chunks]
    lines = ["# Review Corpus Build Log", "", f"generated_at: {utc_today()}", f"chunks: {len(chunks)}", ""]
    lines.append(f"embedding_model: {os.getenv('GX_OPENAI_EMBEDDING_MODEL', DEFAULT_EMBED_MODEL)}")
    lines.append(f"embedding_dimensions: {env_int('GX_EMBEDDING_DIMENSIONS', DEFAULT_EMBED_DIMS)}")
    if tokens:
        lines.append(f"token_histogram: min={min(tokens)} p50={int(statistics.median(tokens))} max={max(tokens)}")
    lines.append("")
    lines.append("## Chunks Per Source")
    for source in sources:
        count = counts[source.id]
        smell = " (investigate: <3 chunks)" if count and count < 3 and source.fetch != "manual" else ""
        lines.append(f"- {source.id}: {count}{smell}")
    (artifact_dir / "build_log.md").write_text("\n".join(lines) + "\n", encoding="utf-8")
    print(f"Wrote {artifact_dir / 'build_log.md'}")


def index_chunks(chunks: list[Chunk], args: argparse.Namespace) -> int:
    openai_key = os.getenv("OPENAI_API_KEY", "").strip()
    tpuf_key = os.getenv("TURBOPUFFER_API_KEY", "").strip()
    if not openai_key:
        print("OPENAI_API_KEY is required for indexing", file=sys.stderr)
        return 2
    if not tpuf_key:
        print("TURBOPUFFER_API_KEY is required for indexing", file=sys.stderr)
        return 2
    embedder = OpenAIEmbedder(
        api_key=openai_key,
        base_url=os.getenv("GX_OPENAI_BASE_URL", DEFAULT_OPENAI_BASE_URL),
        model=os.getenv("GX_OPENAI_EMBEDDING_MODEL", DEFAULT_EMBED_MODEL),
        dimensions=env_int("GX_EMBEDDING_DIMENSIONS", DEFAULT_EMBED_DIMS),
    )
    stores = {
        args.review_namespace: TurboPufferStore(
            api_key=tpuf_key,
            base_url=os.getenv("GX_TPUF_BASE_URL", DEFAULT_TPUF_BASE_URL),
            namespace=args.review_namespace,
            dimensions=embedder.dimensions,
            production_compatible=args.review_namespace == getattr(args, "production_namespace", ""),
        ),
        args.research_namespace: TurboPufferStore(
            api_key=tpuf_key,
            base_url=os.getenv("GX_TPUF_BASE_URL", DEFAULT_TPUF_BASE_URL),
            namespace=args.research_namespace,
            dimensions=embedder.dimensions,
        ),
    }
    grouped: dict[str, list[Chunk]] = collections.defaultdict(list)
    for chunk in chunks:
        namespace = args.research_namespace if chunk.attributes["tier"] == "research" else args.review_namespace
        grouped[namespace].append(chunk)
    for namespace, namespace_chunks in grouped.items():
        indexed = 0
        for batch in batches(namespace_chunks, min(args.batch_size, DEFAULT_BATCH_SIZE)):
            vectors = embedder.embed([chunk.body for chunk in batch])
            rows = []
            for chunk, vector in zip(batch, vectors, strict=True):
                rows.append({"id": chunk.id, "vector": vector, **chunk.attributes})
            stores[namespace].upsert(rows)
            indexed += len(batch)
            print(f"Indexed {indexed}/{len(namespace_chunks)} chunks into {namespace}")
    return 0


def run_eval(golden_path: str, baseline_namespace: str, candidate_namespace: str, artifact_dir: Path) -> int:
    queries = parse_yaml_manifest(Path(golden_path).read_text(encoding="utf-8")).get("queries", [])
    lines = [
        "# Review Corpus Eval Results",
        "",
        f"generated_at: {utc_today()}",
        f"baseline_namespace: {baseline_namespace}",
        f"candidate_namespace: {candidate_namespace}",
        "",
    ]
    if not queries:
        lines.append("No golden queries found.")
        artifact_dir.mkdir(parents=True, exist_ok=True)
        (artifact_dir / "eval_results.md").write_text("\n".join(lines) + "\n", encoding="utf-8")
        return 1
    openai_key = os.getenv("OPENAI_API_KEY", "").strip()
    tpuf_key = os.getenv("TURBOPUFFER_API_KEY", "").strip()
    if not openai_key or not tpuf_key:
        lines.append(f"queries: {len(queries)}")
        lines.append("status: skipped live comparison; OPENAI_API_KEY and TURBOPUFFER_API_KEY are required.")
        artifact_dir.mkdir(parents=True, exist_ok=True)
        (artifact_dir / "eval_results.md").write_text("\n".join(lines) + "\n", encoding="utf-8")
        print(f"Wrote {artifact_dir / 'eval_results.md'}")
        return 0
    embedder = OpenAIEmbedder(
        api_key=openai_key,
        base_url=os.getenv("GX_OPENAI_BASE_URL", DEFAULT_OPENAI_BASE_URL),
        model=os.getenv("GX_OPENAI_EMBEDDING_MODEL", DEFAULT_EMBED_MODEL),
        dimensions=env_int("GX_EMBEDDING_DIMENSIONS", DEFAULT_EMBED_DIMS),
    )
    baseline_store = TurboPufferStore(api_key=tpuf_key, base_url=os.getenv("GX_TPUF_BASE_URL", DEFAULT_TPUF_BASE_URL), namespace=baseline_namespace, dimensions=embedder.dimensions)
    candidate_store = TurboPufferStore(api_key=tpuf_key, base_url=os.getenv("GX_TPUF_BASE_URL", DEFAULT_TPUF_BASE_URL), namespace=candidate_namespace, dimensions=embedder.dimensions)
    baseline_hits = 0
    candidate_hits = 0
    precedence_violations: list[str] = []
    result_rows: list[str] = []
    for item in queries:
        query = str(item.get("query", "")).strip()
        expected = [str(value) for value in item.get("expected_source_ids", [])]
        if not query or not expected:
            continue
        vector = embedder.embed([query])[0]
        baseline_rows = baseline_store.query(vector=vector, text=query, text_attribute="text", limit=10, filters=None, include=eval_baseline_include_attributes())
        candidate_rows = candidate_store.query(vector=vector, text=query, text_attribute="body", limit=10, filters=eval_candidate_filter(), include=eval_candidate_include_attributes())
        baseline_ok = has_expected_source(baseline_rows, expected)
        candidate_ok = has_expected_source(candidate_rows, expected)
        baseline_hits += 1 if baseline_ok else 0
        candidate_hits += 1 if candidate_ok else 0
        bad = precedence_bad_rows(candidate_rows)
        if bad:
            precedence_violations.append(f"{query}: {', '.join(bad)}")
        result_rows.append(f"- {query}: v1={'hit' if baseline_ok else 'miss'} v2={'hit' if candidate_ok else 'miss'} expected={','.join(expected)}")
    total = max(1, len(result_rows))
    lines.append(f"queries: {total}")
    lines.append(f"v1_recall_at_10: {baseline_hits}/{total}")
    lines.append(f"v2_recall_at_10: {candidate_hits}/{total}")
    lines.append(f"precedence_violations: {len(precedence_violations)}")
    passed = candidate_hits >= baseline_hits and not precedence_violations
    lines.append(f"status: {'pass' if passed else 'fail'}")
    lines.append("")
    lines.append("## Query Results")
    lines.extend(result_rows)
    if precedence_violations:
        lines.append("")
        lines.append("## Precedence Violations")
        lines.extend(f"- {item}" for item in precedence_violations)
    artifact_dir.mkdir(parents=True, exist_ok=True)
    (artifact_dir / "eval_results.md").write_text("\n".join(lines) + "\n", encoding="utf-8")
    print(f"Wrote {artifact_dir / 'eval_results.md'}")
    return 0 if passed else 1


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
        data = request_json(f"{self.base_url}/embeddings", payload, headers={"Authorization": f"Bearer {self.api_key}"})
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
    def __init__(self, *, api_key: str, base_url: str, namespace: str, dimensions: int, production_compatible: bool = False) -> None:
        self.api_key = api_key
        self.base_url = base_url.rstrip("/")
        self.namespace = namespace.strip("/")
        self.dimensions = dimensions
        self.production_compatible = production_compatible

    def upsert(self, rows: list[dict[str, Any]]) -> None:
        payload = {
            "distance_metric": "cosine_distance",
            "schema": review_corpus_schema(self.dimensions, production_compatible=self.production_compatible),
            "upsert_rows": [production_compatible_row(row) for row in rows] if self.production_compatible else rows,
        }
        request_json(f"{self.base_url}/v2/namespaces/{urllib.parse.quote(self.namespace)}", payload, headers={"Authorization": f"Bearer {self.api_key}"})

    def query(self, *, vector: list[float], text: str, text_attribute: str, limit: int, filters: Any, include: list[str]) -> list[dict[str, Any]]:
        vector_query: dict[str, Any] = {
            "rank_by": ["vector", "ANN", vector],
            "limit": {"total": limit},
            "include_attributes": include,
        }
        if filters is not None:
            vector_query["filters"] = filters
        payload: dict[str, Any] = vector_query
        if text.strip():
            text_query: dict[str, Any] = {
                "rank_by": [text_attribute, "BM25", text],
                "limit": {"total": limit},
                "include_attributes": include,
            }
            if filters is not None:
                text_query["filters"] = filters
            payload = {"queries": [vector_query, text_query], "rerank_by": ["RRF"]}
        data = request_json(f"{self.base_url}/v2/namespaces/{urllib.parse.quote(self.namespace)}/query", payload, headers={"Authorization": f"Bearer {self.api_key}"})
        rows = data.get("rows") or []
        if not rows and data.get("results"):
            rows = data["results"][0].get("rows") or []
        return rows


def eval_baseline_include_attributes() -> list[str]:
    return [
        "source_id",
        "title",
        "source_kind",
        "text",
    ]


def eval_candidate_include_attributes() -> list[str]:
    return [
        "source_id",
        "title",
        "tier",
        "source_kind",
        "superseded_by",
        "historical",
        "precedence_group",
        "section_path",
        "body",
        "text",
    ]


def eval_candidate_filter() -> list[Any]:
    return [
        "And",
        [
            ["source_kind", "Eq", "review_corpus"],
            ["tier", "NotEq", "research"],
            ["historical", "Eq", False],
            ["superseded_by", "Eq", ""],
        ],
    ]


def has_expected_source(rows: list[dict[str, Any]], expected_source_ids: list[str]) -> bool:
    expected = set(expected_source_ids)
    return any(str(row.get("source_id", "")).strip() in expected for row in rows)


def precedence_bad_rows(rows: list[dict[str, Any]]) -> list[str]:
    bad: list[str] = []
    for row in rows:
        source_id = str(row.get("source_id", "")).strip()
        superseded_by = str(row.get("superseded_by", "")).strip()
        historical = bool(row.get("historical", False))
        if superseded_by or historical:
            bad.append(source_id or str(row.get("title", "unknown")))
    return bad


def request_json(url: str, payload: dict[str, Any], *, headers: dict[str, str]) -> dict[str, Any]:
    body = json.dumps(payload).encode("utf-8")
    request_headers = {"Content-Type": "application/json", **headers}
    retry_statuses = {408, 409, 425, 429, 500, 502, 503, 504}
    last_error: Exception | None = None
    for attempt in range(1, 5):
        request = urllib.request.Request(url, data=body, headers=request_headers, method="POST")
        try:
            with urllib.request.urlopen(request, timeout=60) as response:
                return json.loads(response.read().decode("utf-8"))
        except urllib.error.HTTPError as exc:
            detail = exc.read().decode("utf-8", errors="replace")[:1000]
            error = RuntimeError(f"{url} failed with status {exc.code}: {detail}")
            if exc.code not in retry_statuses:
                raise error from exc
            last_error = error
        except (urllib.error.URLError, TimeoutError, ConnectionError) as exc:
            last_error = RuntimeError(f"{url} request failed: {exc}")
        if attempt < 4:
            time.sleep(min(30, attempt * 5))
    if last_error is not None:
        raise last_error
    raise RuntimeError(f"{url} request failed")


def production_compatible_row(row: dict[str, Any]) -> dict[str, Any]:
    out = dict(row)
    out["authority"] = str(out.get("authority", ""))
    for key in ["languages", "frameworks", "risk_tags", "review_tags"]:
        value = out.get(key)
        if isinstance(value, list):
            out[key] = ",".join(str(item) for item in value)
    return out


def review_corpus_schema(dimensions: int, *, production_compatible: bool = False) -> dict[str, Any]:
    string_filter = {"type": "string", "filterable": True}
    string_search = {"type": "string", "full_text_search": True}
    string_array = {
        "type": "[]string",
        "filterable": True,
        "full_text_search": {"stemming": False, "remove_stopwords": False, "case_sensitive": False},
    }
    return {
        "vector": {"type": f"[{dimensions}]f32", "ann": True},
        "body": string_search,
        "text": string_search,
        "source_kind": string_filter,
        "chunk_id": string_filter,
        "source_id": string_filter,
        "source_url": string_filter,
        "url": string_filter,
        "title": string_search,
        "tier": string_filter,
        "category": string_filter,
        "languages": string_search if production_compatible else string_array,
        "language_tags": string_array,
        "authority": string_filter if production_compatible else {"type": "uint", "filterable": True},
        "precedence_group": string_filter,
        "superseded_by": string_filter,
        "historical": {"type": "bool", "filterable": True},
        "section_path": string_search,
        "cwe_ids": string_array,
        "content_sha256": string_filter,
        "content_hash": string_filter,
        "fetched_at": string_filter,
        "usefulness_rank": {"type": "uint", "filterable": True},
        "chunk_index": {"type": "uint", "filterable": True},
    }


def title_from_source(source: Source) -> str:
    return source.id.replace("-", " ").title()


def extract_cwe_ids(text: str) -> list[str]:
    return sorted(set(re.findall(r"\bCWE-\d+\b", text, flags=re.IGNORECASE)))


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


def nullable_token(value: Any) -> str | None:
    if value in {None, "", "null"}:
        return None
    token = clean_token(value)
    return token or None


def slugify(value: str) -> str:
    return clean_token(value).replace(".", "-")[:80]


def truncate_ascii_token(value: str, max_bytes: int) -> str:
    encoded = value.encode("utf-8")[:max_bytes]
    return encoded.decode("utf-8", errors="ignore").strip("-") or "document"


def sha256_text(text: str) -> str:
    return hashlib.sha256(text.encode("utf-8")).hexdigest()


def sha256_bytes(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def rough_token_count(text: str) -> int:
    return max(1, len(re.findall(r"\S+", text)))


def utc_today() -> str:
    return dt.datetime.now(dt.UTC).date().isoformat()


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
