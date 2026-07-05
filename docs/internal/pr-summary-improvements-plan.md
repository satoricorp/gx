# GX PR Summary Improvements — Implementation Plan

This is a self-contained implementation plan for improving the GitHub PR body summaries that
`gx push` generates. It assumes no prior conversation context. Read the "Current architecture"
section first; every phase references exact files, types, and line anchors in this repo.

A sibling plan exists conceptually for the `gx review` command (not yet committed on this
branch). The two features share `internal/codereview` (the `Finding` type, the AI reviewer,
the system prompt, and the JSON response schema). **Phase 1 of this plan is a single
shared-contract change to that machinery** — it covers the schema needs of both features
(anchors for PR summaries, source attribution for `gx review`) so the two cannot drift while
editing the same code. If a sibling review plan lands schema changes first, reconcile against
whatever shipped; never duplicate the parser or prompt.

## Product requirements (why this plan exists)

1. **The summary must be trust-building, not exhaustive.** A reviewer should read the PR body
   and know how to spend their time: a brief overview of what changed and why, a credible
   blast radius, and only the changes that genuinely need human eyes — each linked to the
   exact hunk on GitHub.
2. **PRs that don't need review must say so explicitly.** Docs-only, tests-only, and other
   minimal/low-impact changes should get an affirmative "no review needed" verdict, not a
   passive "no targets found" sentence. This verdict may only be granted when the AI pass
   actually **succeeded** — never from heuristics alone, and never when model output was
   unavailable or errored.
3. **Blast radius must be grounded and honestly labeled.** It should reflect how much of the
   codebase references the changed code. The v1 mechanism is **lexical reach** (text search
   for changed symbols) — useful, cheap, and capped — and the output must be labeled as
   references, not presented as a true dependency graph. It will miss dynamic dispatch and
   aliased imports and can overcount common names; the rendering must not overclaim.
4. **Every rendered link must be correct.** A finding linked to the wrong hunk is worse than
   no link. Never guess; validate anchors structurally and drop links that don't validate.
5. **Attribution rules:** external/independent knowledge sources are **opaque** (a kind label
   like "review resource" — no titles, no URLs). References to lines of code are **specific
   and linked** — to the PR diff when the code is in the diff, to a GitHub blob permalink
   (pinned to the head SHA) when it is elsewhere in the codebase.
6. **The summary must generalize beyond the gx repo itself.** Today's heuristics hardcode
   gx's own package paths and product vocabulary; they must come from repo config instead.
7. **Degradation must be visible.** If the AI pass did not succeed and the summary is
   heuristics-only, the body should say so.

## Current architecture (as of this plan)

Pipeline entry: `gx push` → `publication.EnqueuePush` (`internal/publication/outbox.go:60`) →
`UpdateGitHubPullRequestBody` (`internal/publication/github_pr.go:22`) →
`GitHubPullRequestBodyFromArtifact` (`internal/publication/pr_body.go:106`) → PATCH to
`/repos/{owner}/{repo}/pulls/{n}` via `internal/github/client.go:257` `UpdatePullRequest`.
The same body generation also runs from `PublishArtifact` (`internal/publication/publication.go:189`).

`GitHubPullRequestBodyFromArtifact` flow (`pr_body.go`):

1. `buildPRBodyCatalog` (`pr_body.go:625`) — parses stack patches into `prHunkSummary` values
   (file, hunk header, line ranges, and a GitHub deep link built from
   `{prURL}/files#diff-{sha256(path)}R{line}`, see `pr_body.go:805-822`), collects
   `prRevisionSummary` values, and computes `prBodyStats` via `bodyStats` (`pr_body.go:884`).
2. `collectDefaultPRSummaryContext` (`pr_body.go:438`) — local docs (`AGENTS.md`, `CONTEXT.md`,
   `REVIEW.md`, `README.md`), changed-file snippets, linked session transcripts
   (`semantic.BuildSessionChunks`), and TurboPuffer queries (indexed code/sessions plus a
   `gx-review-knowledge` namespace) when `semantic.ConfigFromEnv` is enabled.
3. `reviewPRSummaryFindings` (`pr_body.go:383`) — calls `codereview.ReviewerFromEnv()` with
   `prReviewBrief` (`pr_body.go:398`): profile `patch_focused`, up to 32 diff hunks, the
   context snippets, and a PR-specific `ArchitectureRubric`. **Errors are silently swallowed**
   (returns `nil`); there are no retries. Note `ReviewerFromEnv` can return a non-nil
   `multiAIReviewer` whose providers all fail at call time — a non-nil reviewer does not mean
   the AI pass succeeded.
4. `renderGitHubPullRequestBody` (`pr_body.go:116`) — marker comment
   (`<!-- gx:pr-summary:v1 -->`), template opening paragraph (`openingSummary`,
   `pr_body.go:152`), and a `## Needs Review` list from `needsReviewItems` (`pr_body.go:201`):
   AI findings merged with heuristic items, deduped, score-sorted, capped at
   `maxPRReviewItems = 10`.
5. Author-authored PR body content is preserved beneath `<!-- gx:author-notes -->`
   (`appendAuthorNotes` / `splitGeneratedPRBody`, `pr_body.go:1093-1118`).

Key behaviors this plan changes:

- **Opening paragraph is all templates.** `changeSummarySentence` (`pr_body.go:160`) joins
  revision descriptions; `readinessSentence` (`pr_body.go:179`) prints blast radius from
  `bodyStats` — which falls back to raw thresholds (8+ files or 400+ changed lines → "high",
  `pr_body.go:914-921`). No reference analysis; no LLM narrative.
- **Finding→hunk linking is lexical.** `hunkLinkForFinding` (`pr_body.go:824`) substring-searches
  the finding prose for a file name and, when nothing matches, **falls back to the first hunk**
  (`pr_body.go:836-838`). `patchRelevantFinding` (`pr_body.go:234`) and `findingScore`
  (`pr_body.go:261`) use keyword lists; `lowValuePRFinding` (`pr_body.go:248`) filters on a
  gx-specific word allowlist.
- **Heuristic items are hardcoded to the gx repo.** `hunkImpact` (`pr_body.go:359`) matches
  `internal/storage/`, `internal/github/`, `internal/publication/`, `internal/vcs/`,
  `internal/auth` and emits fixed gx-specific prose. The PR rubric (`pr_body.go:420-434`)
  says "GX's behavior".
- **Attribution.** `findingAttributions` (`pr_body.go:1022`) classifies evidence by substring
  ("session", "review") and, when a finding has no evidence, **stamps the first two context
  snippets onto it** (`pr_body.go:1048-1053`). `renderAttributions` (`pr_body.go:1070`) renders
  live markdown URLs for TurboPuffer review-knowledge rows (`pr_body.go:1082-1083`) — violating
  requirement 5's opaqueness, and contradicting the system prompt's own "never output source
  titles or URLs" instruction.

Shared `codereview` machinery used here:

- `Finding` (`internal/codereview/findings.go:11`): `ID, Scopes, Title, Summary, Benefit,
  Evidence []Evidence{Label, Value}, Recommendation, Strength, SourceIDs`. Today only
  heuristic rules set `SourceIDs`; the AI JSON schema has no source field, so AI findings
  never carry attribution. Phase 1 fixes the schema for both consumers.
- `ReviewBrief` (`internal/codereview/brief.go:19`) — the JSON payload sent to the model.
  Includes `SourceRefs []SourceRef` (`brief.go:97`) and `SourceCatalog []SourceBrief`
  (`brief.go:92`) with labeled context snippets (`L1`, `R2`, `C3`) produced by
  `labelContextSnippets` (`brief.go:306`) and `sourceRefsFromContextSnippets` (`brief.go:467`).
  **`SourceRef` already carries `Kind`, `Title`, `URL`, `File`, `StartLine`, `EndLine`,
  `SessionID`, etc.** — reuse it; do not invent a parallel type.
- `AIReviewer` (`internal/codereview/ai.go:37`) — exported interface with only
  `Review(ctx, brief) ([]Finding, error)`. Many test fakes and wrappers implement it; **do
  not add methods to this interface.**
- Reviewer wrappers: `multiAIReviewer` (`ai.go:241`), `fallbackAIReviewer` (`ai.go:277`),
  `unavailableAIReviewer` (`ai.go:237`). **Critical current behavior:** `multiAIReviewer`
  returns an **error** when every provider returns zero findings (`ai.go:271-274`,
  "AI reviewers returned no findings"), and `fallbackAIReviewer` treats
  `len(findings) == 0` as failure (`ai.go:279`). Phase 1.6 changes this — a parseable
  empty response must count as success.
- AI callers: `responsesAIReviewer.Review` (`ai.go:298`, OpenAI `/v1/responses`,
  `json_object` format, default model `gpt-4.1-mini`, no retry) and
  `bedrockAnthropicReviewer.Review` (`ai.go:353`). Both parse via `parseAIReviewContent`
  (`ai.go:460`). The output shape is dictated by the last lines of `reviewDeveloperPrompt()`
  (`ai.go:657-688`).
- `renderFindingAttributions` (`internal/codereview/review.go:187`) resolves `finding.SourceIDs`
  only against `report.Sources` (the static catalog), **not** `ReviewBrief.SourceRefs`.
  AI-resolved labels would disappear in `gx review` output unless fixed in Phase 1.
- Test seam in `pr_body.go:37`: `prSummaryReviewerFromEnv` is a package var that tests
  replace (see `pr_body_test.go:14,64` and `github_pr_test.go:19,106`). Phase 1.4 replaces
  it with `prSummaryReviewerFromEnvWithInfo`.

Relevant data already available but unused:

- `ReviewContextPayload.ChangedSymbols` (`internal/reviewbundle/bundle.go:97`) — per-hunk
  changed symbol names. Only counted today (`accumulateContext`, `pr_body.go:932`), never
  traversed. This is the input for lexical reach (Phase 4).
- Revision `CommitID` (`prRevisionSummary`, `pr_body.go:42`) and the PR URL
  (`pullRequestURL`, `pr_body.go:1003`) — enough to build blob permalinks
  `https://github.com/{owner}/{repo}/blob/{sha}/{path}#L{n}`.

Existing tests to keep green:

- `internal/publication/pr_body_test.go` — zero-target fallback, 10-item cap.
- `internal/publication/github_pr_test.go` — end-to-end enqueue → PATCH, blast radius text,
  attribution text, hunk links, author-notes preservation.
- `internal/cli/root_test.go` review command tests and `internal/codereview/review_test.go`
  (shared prompt/schema assertions — `TestAIReviewPromptSeparatesPatchAndDeepReview` asserts
  on prompt content; update it deliberately when the prompt changes).

Constraints for the implementing agent:

- Go codebase. After each phase run:
  `go test ./internal/publication/... ./internal/codereview/... ./internal/cli/... ./internal/github/...`
- Version control uses GX (`gx generate` / `gx status` / `gx push`), not raw `git commit`.
  Keep each phase an independently shippable revision: run `gx generate` after the phase's
  tests pass.
- Env-gate anything new that adds latency or network calls, defaulting to on, following the
  existing pattern (`GX_REVIEW_AI`, `GX_REVIEW_RESOURCES`, ...).
- PR body generation runs inside `gx push`; keep added latency bounded (explicit caps and
  timeouts on any new scanning).
- Do not change the `<!-- gx:pr-summary:v1 -->` / `<!-- gx:author-notes -->` contract.
  Everything above the author-notes marker is GX-owned and may be redesigned freely.
- **Shared-code rule:** any change to `Finding`, the AI JSON schema, `reviewDeveloperPrompt`,
  or the parsers must work for **both** `gx review` and PR summaries, with tests for both
  providers (OpenAI responses + Bedrock) and both consumers (`RenderMarkdown` and
  `renderGitHubPullRequestBody`). Never fork these per feature.

### Naming glossary (avoid four almost-identical nouns)

| Name | What it is |
|------|------------|
| `SourceCatalog` / `Source` | Static catalog entries (OWASP, Fowler, …) with `ID`, `Title`, `URL`, `Publisher` |
| `SourceRef` | Per-snippet metadata in `ReviewBrief` — label (`L1`), kind, file, line range, publisher |
| `source_labels` (JSON) | Raw model output: labels the model says it used (`["R1","L2"]`) |
| `ResolvedSource` | Structured result after resolution — see Phase 1.3 |
| `Finding.SourceIDs` | Resolved catalog IDs for heuristic findings (unchanged meaning) |
| `Finding.ResolvedSources` | Structured attribution for AI findings (new — see Phase 1.1) |

Do **not** add a `Finding.Sources` field; it collides with `report.Sources`, `SourceRefs`, and
`SourceCatalog`.

### Shared contract data flow

```
AI JSON
  ├─ overview?              → PR summary only (pr_summary profile)
  ├─ recommendations[]
  │    ├─ file, line        → PR hunk validation (Phase 2)
  │    └─ source_labels[]   → resolveSourceLabels(brief, labels)
  │                            → []ResolvedSource{ID, Kind, Publisher, Opaque, File, StartLine, ...}
  │                            ├─ gx review: ResolvedSourceLabel + plain ref (Phase 1.5)
  │                            └─ PR summary: ResolvedSourceLabel + hunk/blob links (Phase 2.2)
  └─ parser tolerant path   → old responses without new fields still work
```

---

## Phase 1 — Shared codereview contract (schema, prompt, parser)

**Goal:** one coordinated change to the shared AI contract: structural anchors (`file`,
`line`), source attribution (`source_labels` → structured resolution), an optional PR-only
`overview`, empty-response success semantics, and PR-brief source wiring — all parsed by one
parser, requested by one prompt with a profile guard.

### 1.1 `Finding` type (`internal/codereview/findings.go:11`)

Add fields (all additive; existing consumers ignore them):

```go
File             string           // changed file path the finding is anchored to ("" when not anchored)
Line             int              // changed line within that file's hunk (0 when not anchored)
ResolvedSources  []ResolvedSource // populated for AI findings after label resolution
```

Do **not** store raw model labels on `Finding`. Resolve during parsing and discard labels.

```go
// internal/codereview/sources.go (new)
type ResolvedSource struct {
    ID        string // catalog ID when applicable
    Kind      string // "catalog" | "local_doc" | "session" | "indexed_code" | "review_resource" | ...
    Publisher string // opaque publisher label, e.g. "OWASP", "local", "session"
    Opaque    bool   // true → render publisher/kind only, never title/URL/ref/file
    Ref       string // human-readable ref label (snippet ref, session id, doc path)
    File      string // repo-relative path for codebase links ("" when not applicable)
    StartLine int    // line for path#L{n} links (0 when unknown)
    EndLine   int    // optional range end (0 when unknown)
    URL       string // internal use only (catalog URL, TurboPuffer row URL); never render when Opaque
}
```

**Render rules (strict — no guessing from `ID` alone):**

- `Opaque == true` → show `Publisher` (or `Kind`) only; ignore `Ref`, `File`, `StartLine`, `URL`.
- `Opaque == false` with `File != ""` → link to `File` + `#L{StartLine}` when `StartLine > 0`
  (PR hunk link if file is in the diff; blob permalink if not — Phase 2.2).
- `Opaque == false` with `File == ""` and `Ref != ""` → show `Ref` as plain text (sessions).

Heuristic findings continue to use `SourceIDs` only. AI findings use `ResolvedSources`. Both
render paths understand both (heuristic → lookup catalog; AI → use `ResolvedSources` directly).

### 1.1b Publisher metadata on upstream types

Add `Publisher string` to:

- `Source` (`internal/codereview/catalog.go:3`)
- `SourceBrief` (`internal/codereview/brief.go:92`)
- `SourceRef` (`internal/codereview/brief.go:97`)

Population rules:

- **Static catalog:** set `Publisher` on each `sourceCatalog` entry — e.g. `"OWASP"`,
  `"Google"`, `"Fowler"`, `"NIST"`, `"OpenSSF"`. Copy into `SourceBrief` when building the
  brief's `SourceCatalog`.
- **Retrieved refs (`SourceRef`):** in `sourceRefFromContextSnippet` (`brief.go:484`), set
  `Publisher` from, in order: snippet metadata if present → kind-derived fallback (`"local"`
  for `local_doc`/`changed_file`/`repo_doc`, `"session"` for session kinds, `"indexed"` for
  indexed code) → URL host for external URLs (e.g. `owasp.org` → `"OWASP"`) → `"unknown"`.
- **Resolver:** `resolveSourceLabels` copies `Publisher`, `Ref`, `File`, `StartLine`,
  `EndLine`, `URL` from the matched `SourceRef` or catalog `Source` into `ResolvedSource`.
  Do not infer file/line from `ID` alone.

### 1.2 Response schema and prompt (`internal/codereview/ai.go:657-688`)

- Change the shape line of `reviewDeveloperPrompt()` to:
  `{"overview":string(optional),"recommendations":[{"title":string,"summary":string,"benefit":string,"recommendation":string,"strength":"Strong|Worth exploring|Speculative","evidence":[string],"file":string(optional),"line":number(optional),"source_labels":[string](optional)}]}`
- Add instruction lines:
  - *"Set file to the exact changed file path from static.diff_snippets and line to a changed
    line number inside that hunk. If a recommendation cannot be tied to a specific changed
    file, omit file and line."*
  - *"Set source_labels to the labels of context snippets or source_refs you actually relied
    on (e.g. R1, L2). Omit labels you did not use."*
  - **Profile-guarded overview:** *"Only when review_profile is pr_summary: include a
    top-level overview field, 2-3 sentences on what this change does and why, based on the
    revision descriptions and session_transcript context; no file lists, no URLs, no praise.
    For all other profiles, omit overview."*
  - State that `pr_summary` behaves like `patch_focused` for finding selection (current-change
    review, changed-lines evidence, same rejection rules — no quota-filling, no generic
    advice) plus the overview rule.
- `reviewDeveloperPrompt()` remains a single shared prompt. Do **not** create a second prompt
  function for PR summaries; the profile guard is the only branch point.

### 1.3 Parser and label resolver (`internal/codereview/ai.go:460`)

- Rework `parseAIReviewContent` into
  `parseAIReviewOutput(content string) (aiReviewOutput, error)` where
  `aiReviewOutput{Overview string; Findings []Finding}`. Keep
  `parseAIReviewContent(content) ([]Finding, error)` as a thin wrapper that discards the
  overview so existing call sites compile unchanged.
- `overview` is optional everywhere; `line` must tolerate arriving as a JSON string; unknown
  fields are ignored. **A parseable response with an empty `recommendations` array is valid
  output, not an error.**
- After parsing each recommendation, call
  `resolveSourceLabels(brief ReviewBrief, labels []string) []ResolvedSource`:
  - Look up each label in `brief.SourceRefs` (primary) and `brief.SourceCatalog` (fallback
    for static catalog IDs).
  - Map `SourceRef.Kind` → `ResolvedSource.Kind`; set `Opaque: true` for
    `resource` / `reference` / external catalog entries; set `Opaque: false` for `local`,
    `code`, `session`, `policy` kinds (see `sourceRefKind`, `brief.go:509`).
  - Copy link fields from the matched ref: `Ref`, `File`, `StartLine`, `EndLine`, `URL`,
    `Publisher` (§1.1b). Catalog matches populate `Publisher` from `Source.Publisher`;
    `File`/`StartLine` come from `SourceRef` only (catalog entries have no file anchor).
  - Unknown labels are dropped silently.
- Wire resolver into both provider `Review` methods after parsing, before returning findings.

### 1.4 PR-summary test seam and optional overview interface

**Do not add methods to `AIReviewer`.** Keep the exported interface unchanged:

```go
type AIReviewer interface {
    Review(ctx context.Context, brief ReviewBrief) ([]Finding, error)
}
```

Add a separate optional interface in `internal/codereview/ai.go`:

```go
type AIReviewerWithOverview interface {
    AIReviewer
    ReviewWithOverview(ctx context.Context, brief ReviewBrief) (overview string, findings []Finding, err error)
}
```

Implement `ReviewWithOverview` on `responsesAIReviewer`, `bedrockAnthropicReviewer`,
`multiAIReviewer`, and `fallbackAIReviewer`. Each existing `Review` method delegates to
`ReviewWithOverview` and discards the overview — `gx review` call paths unchanged.

Add `ReviewerFromEnvWithInfo() (AIReviewer, ReviewerInfo)` in `codereview` where
`ReviewerInfo{Models []string}`; keep `ReviewerFromEnv()` as a wrapper that discards info.

**Replace the existing test seam** in `pr_body.go:37`. Remove `prSummaryReviewerFromEnv`;
add one seam that covers reviewer + model info:

```go
var prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
    return codereview.ReviewerFromEnvWithInfo()
}
```

`reviewPRSummaryFindings` calls **only** `prSummaryReviewerFromEnvWithInfo()` — never
`codereview.ReviewerFromEnvWithInfo()` directly. Update the existing overrides once:
`pr_body_test.go:14,64` and `github_pr_test.go:19,106`. Fake reviewers should implement
`AIReviewerWithOverview` when tests need overview or `aiSucceeded` behavior.

```go
reviewer, info := prSummaryReviewerFromEnvWithInfo()
if withOverview, ok := reviewer.(codereview.AIReviewerWithOverview); ok {
    overview, findings, err = withOverview.ReviewWithOverview(ctx, brief)
} else {
    findings, err = reviewer.Review(ctx, brief)
}
// info.Models used for provenance footer (Phase 7); nil Models → omit model clause
```

### 1.5 `gx review` attribution fix (`internal/codereview/review.go:187`)

Update `renderFindingAttributions` so AI findings do not lose their sources:

- If `len(finding.ResolvedSources) > 0`: render from `ResolvedSources` using
  `ResolvedSourceLabel(src)` (below) plus plain-text ref/file when non-opaque.
- Else if `len(finding.SourceIDs) > 0`: keep current catalog lookup via `report.Sources`.

**Attribution helper boundary — do not share link rendering across packages:**

- Add exported `codereview.ResolvedSourceLabel(src ResolvedSource) string` in
  `internal/codereview/sources.go`. Returns opaque publisher/kind label, or a plain
  `File`/`Ref` string for non-opaque sources — **no markdown links**. This is the only
  cross-package helper.
- PR-specific link building stays in `internal/publication`: `renderPRResolvedSource(src,
  artifact, catalog)` wraps `ResolvedSourceLabel` and, when `src.File != ""` and not opaque,
  appends a markdown link via hunk catalog or `blobPermalink`. Do **not** put artifact/PR
  URL/head-SHA logic in `codereview`.

Add test: AI finding with `ResolvedSources` containing an opaque `review_resource` and a
non-opaque `local_doc` → both appear in `RenderMarkdown` attribution line (plain text only
for gx review — no blob permalinks in CLI output).

### 1.6 Empty-response success semantics (BLOCKING correctness fix)

Current wrapper behavior treats zero findings as failure, which would make the Phase 3
verdict impossible:

- `multiAIReviewer.Review` (`ai.go:271-274`) returns `error("AI reviewers returned no
  findings")` when all providers parse but return zero findings.
- `fallbackAIReviewer.Review` (`ai.go:279`) treats `err == nil && len(findings) > 0` as the
  only success, falling through on legitimate empty results.

Change both (in their new `ReviewWithOverview` implementations, with `Review` delegating):

- A provider response that **parses** (valid JSON matching the schema) is a success even
  when `recommendations` is empty. `multiAIReviewer` returns `(overview, []Finding{}, nil)`
  when at least one provider parsed successfully and none produced findings. It returns an
  error only when **every** provider failed to produce a parseable response.
- `fallbackAIReviewer` only invokes the fallback when the primary returns an actual error —
  not when it returns zero findings.
- Preserve the 6-finding cap and dedup behavior for the non-empty path.

Tests (required): canned `{"recommendations":[]}` through `multiAIReviewer` with one healthy
provider → no error, zero findings; both providers erroring → error; primary empty +
fallback configured → fallback **not** called.

### 1.7 PR brief source wiring (BLOCKING correctness fix)

`prReviewBrief` (`pr_body.go:398`) currently sets only `Context:` — no `SourceRefs`, no
`SourceCatalog`, and snippets never pass through `labelContextSnippets`. The model cannot
emit `source_labels` for the PR profile without them, making §1.3's resolver dead code for
PR summaries.

In `prReviewBrief`:

- Run the context snippets through `codereview.labelContextSnippets` (export it or add an
  exported wrapper `LabelContextSnippets`) so each gets an `L*`/`R*` `SourceLabel`.
- Set `SourceRefs: sourceRefsFromContextSnippets(labeled)` (same export treatment).
- Set `SourceCatalog` the same way `BuildReviewBrief` does (`brief.go:162` region) so
  catalog fallback resolution works.
- Change `ReviewProfile` from `patch_focused` to `pr_summary`.

Test: `prReviewBrief` output contains labeled context snippets and non-empty `SourceRefs`
mirroring them.

### 1.8 Tests (phase gate)

- Both providers: canned responses with and without `file`/`line`/`source_labels`/`overview`
  parse correctly (extend the fake-transport pattern in `ai_test.go`).
- Label resolver: labels resolve through `SourceRefs` with correct `Kind`, `Opaque`,
  `Publisher`, `File`, `StartLine`; unknown labels dropped.
- Empty-response semantics per §1.6.
- PR brief wiring per §1.7.
- Prompt assertions: update `TestAIReviewPromptSeparatesPatchAndDeepReview`
  (`review_test.go:597`) — shape line contains `file`, `line`, `source_labels`, `overview`,
  and the `pr_summary` guard.
- `gx review` regression: `RenderMarkdown` shows attribution for AI findings with
  `ResolvedSources`; heuristic findings unchanged.
- Seam migration: all four existing overrides updated; plain `Review` on fakes still works.

## Phase 2 — PR rendering correctness (anchors + attribution)

**Goal:** the renderer never guesses. Links validate against the hunk catalog or don't render;
attribution reflects only actual evidence; external knowledge is opaque; codebase references
are specific links.

All in `internal/publication/pr_body.go` unless noted.

### 2.1 Validated anchors

- Replace `hunkLinkForFinding` (`pr_body.go:824`) with validation logic:
  1. `finding.File` exactly matches a hunk file AND `finding.Line` falls within that hunk's
     new-line range → that hunk's link, with the `R{line}` anchor using `finding.Line`.
  2. `finding.File` matches a hunk file but the line doesn't validate → the file-level anchor
     (`#diff-{sha256(path)}` with no line suffix).
  3. `finding.File` empty or matches nothing → **no link**. Delete the first-hunk fallback
     (`pr_body.go:836-838`). An unlinked finding is honest; a mislinked one is not.
  - Keep the old lexical match only as a tertiary attempt when `finding.File` is empty
    (findings from providers not yet returning the field).
- Replace `patchRelevantFinding` (`pr_body.go:234`) with: `finding.File` is one of
  `catalog.Files` (fall back to the current lexical check only when `File` is empty).
- Replace `findingScore` (`pr_body.go:261`) keyword scoring with strength-based scoring:
  Blocking=300, Strong=200, "Worth exploring"=100, Speculative=0, plus +50 when the anchor
  validated at hunk level. Drop Speculative findings from the PR body entirely.
- Delete `lowValuePRFinding`'s keyword allowlist (`pr_body.go:256-258`); keep only the
  "this is good/correct" praise filter.

### 2.2 Attribution

- **Delete the fallback stamping** of context snippets onto findings with no evidence
  (`findingAttributions`, `pr_body.go:1048-1053`). If a finding has no evidence and no
  resolved sources, it renders with no attribution line.
- **Consume `ResolvedSources` for AI findings.** Primary path:
  `for _, src := range finding.ResolvedSources { renderPRResolvedSource(src, artifact, catalog) }`.
  Uses `codereview.ResolvedSourceLabel(src)` for the label portion, then adds hunk/blob
  markdown links when `src.File != ""` and not opaque. Do **not** treat all `SourceIDs` as
  `review_resource` — that misclassifies local/session refs.
- For heuristic findings, keep the existing `SourceIDs` / evidence paths; classify evidence
  items by resolving against `summaryContext.Snippets` by `Ref`/`SourceLabel` and taking the
  snippet's structured `Kind`/`Source` (as `attributionFromSnippet` does, `pr_body.go:1057`)
  instead of substring sniffing (`pr_body.go:1036-1042`).
- **Opaque external sources:** `ResolvedSource.Opaque == true` → kind/publisher label only,
  never ref/URL/title. Delete the URL-rendering branch (`pr_body.go:1082-1083`) for opaque
  kinds.
- **Specific codebase links.** For non-opaque sources with `ResolvedSource.File != ""`:
  - If `File` is in the PR diff → hunk deep link via catalog (2.1 logic), using
    `src.StartLine` when attributing the source snippet itself.
  - If `File` is outside the diff → blob permalink via `blobPermalink(prURL, headSHA,
    src.File, src.StartLine)`. Requires `StartLine > 0` for a line anchor; file-only link
    when line is unknown. Add `blobPermalink(prURL, sha, path, line)` parsing owner/repo
    from `pullRequestURL(artifact)` (`pr_body.go:1003`); head SHA from the last revision's
    `CommitID`. Unit test the helper.
  - **If `pullRequestURL(artifact) == ""`, never emit blob links** — render
    `ResolvedSourceLabel` + plain `File:line` text instead.
- Session attributions (`Kind == "session"`): `ResolvedSourceLabel` + plain `Ref`, no URL.

### 2.3 Tests

- Valid `{file, line}` → link to exactly that hunk and line; unresolvable file → **no** link
  (today it would link to the first hunk — assert the fix).
- AI finding with mixed `ResolvedSources` (opaque review resource + local doc) → correct kinds;
  external resource renders opaque; no-evidence finding renders no attribution; out-of-diff
  codebase ref renders a blob permalink; empty PR URL → plain text, no links.
- Update `TestEnqueueArtifactUpdatesGitHubPullRequestBodyFromReviewBundle`
  (`github_pr_test.go`) expectations.

## Phase 3 — Change classification and an explicit review verdict

**Goal:** the first line of every PR body is a verdict: `No review needed`, `Quick scan`, or
`Full review`, with a one-sentence reason. Deterministic facts propose; a **successful** AI
pass confirms or downgrades — never grants.

### 3.1 Shared classifier: `codereview.TriageChange`

`internal/codereview/triage.go` does not exist yet — this plan owns it. Implement:

```go
type ChangeTriage struct {
    Class     string   // "docs-only" | "tests-only" | "config-only" | "mechanical" | "code" | "security-sensitive"
    RiskTags  []string
    Rationale []string // human-readable, e.g. "all 4 files match docs/** or *.md"
}

func TriageChange(changedFiles []string, diffSnippets []DiffSnippet) ChangeTriage
```

Classification rules (deterministic, no LLM):

- **docs-only**: every file matches `*.md`, `*.mdx`, `*.txt`, `*.rst`, `docs/**`, `LICENSE*`,
  images. Exception: `REVIEW.md` / `AGENTS.md` / `SECURITY.md` are config, not docs.
- **tests-only**: every file matches `*_test.go`, `test/**`, `testdata/**`.
- **config-only**: lockfiles, CI configs, `.gitignore`, editor configs — no source code.
- **mechanical**: pure renames, generated files (`_generated`, `.pb.go`), comment-only hunks.
- **security-sensitive**: any changed path or hunk matching auth/session/crypto/SQL/exec/
  deserialization patterns.
- **code**: everything else.

The publication package adapts its inputs: changed files come from `catalog.Files`, diff
snippets from `catalog.Hunks` (`prHunkSummary.Patch`). Add a small adapter in `pr_body.go`,
not a second classifier.

**Cross-feature rule for `tests-only`:** PR summaries may use `tests-only` for the
`No review needed` verdict. Any future `gx review` short-circuit must **not** treat
`tests-only` like `docs-only` — test changes can hide behavior regressions. Gate any review
short-circuit on an explicit allowlist that omits `tests-only`.

### 3.2 Verdict

```go
func reviewVerdict(triage codereview.ChangeTriage, stats prBodyStats, items []prNeedsReviewItem, aiSucceeded bool) (verdict string, reason string)
```

- `aiSucceeded` means the reviewer call returned a **parseable model response** without error
  — not merely that a reviewer was constructed or attempted, and **not** dependent on finding
  count:

  ```go
  aiSucceeded = (err == nil) && responseParsed
  // responseParsed = JSON unmarshaled and schema accepted; len(findings) may be 0
  ```

  An empty `recommendations` array is a **successful** output (guaranteed by Phase 1.6) and
  is required for the `No review needed` verdict.

  Change `reviewPRSummaryFindings` (`pr_body.go:383-396`) to return
  `(findings []codereview.Finding, overview string, aiSucceeded bool, info codereview.ReviewerInfo)`
  — reviewer + info come from `prSummaryReviewerFromEnvWithInfo()` (Phase 1.4); overview via
  optional `AIReviewerWithOverview` type assertion.
- `No review needed` requires **all** of: triage class ∈ {docs-only, tests-only, config-only,
  mechanical}; `stats.WarningCount == 0`; `stats.MaxRiskLevel == "low"`; zero findings of
  strength Strong or Blocking; and `aiSucceeded == true`.
- `Full review` when: any Strong/Blocking finding, or `stats.MaxRiskLevel == "high"`, or
  triage class is `security-sensitive`, or (after Phase 4) lexical reach is high.
- `Quick scan` otherwise — including every case where `aiSucceeded == false`.
- Reason string comes from `triage.Rationale` plus the deciding condition.

### 3.3 Rendering

In `renderGitHubPullRequestBody` (`pr_body.go:116`), before the opening paragraph:

```markdown
**Review verdict: No review needed** — docs-only change (4 files), no warnings, no findings.
```

When the verdict is `No review needed`, suppress heuristic filler items so the Needs Review
section renders the explicit zero-target line rather than manufactured targets.

Tests: docs-only fixture with successful (empty) AI pass → `No review needed`; same fixture
with reviewer erroring → `Quick scan`; fixture with a Strong finding → `Full review`;
tests-only classification.

## Phase 4 — Lexical reach (grounded blast radius)

**Goal:** blast radius reflects how much of the codebase **references** the changed code, and
the body links the top external references. This is lexical search, not a dependency graph —
name it and render it accordingly; it must never overclaim.

### 4.1 Reach computation

New file `internal/publication/pr_reach.go`:

```go
type symbolReach struct {
    Symbol     string
    References []reachSite // File, Line — outside the PR's changed files
}

type lexicalReach struct {
    Symbols        []symbolReach
    ReferenceCount int
    DependentFiles int
    Truncated      bool
}

func computeLexicalReach(ctx context.Context, repoRoot string, catalog prBodyCatalog, artifact reviewbundle.Artifact) lexicalReach
```

- Input symbols: `ReviewContextPayload.ChangedSymbols` (`internal/reviewbundle/bundle.go:97`)
  gathered across `artifact.Stack` entries (same walk as `reviewWarnings`, `pr_body.go:871`).
  Fallback when empty: parse top-level `func`/`type` names from hunk patches.
- For each symbol run `git grep -nw --fixed-strings <symbol>` at `repoRoot`, excluding the
  PR's own changed files and test files.
- Guards: skip symbols shorter than 4 chars or with ≥ 200 matches (too generic — set
  `Truncated` and exclude from counts rather than reporting a misleading number); cap at 20
  symbols and 10 references per symbol; whole computation wrapped in a 10-second
  `context.WithTimeout`; env gate `GX_PR_BLAST_RADIUS=0` disables (default enabled). Any
  error → return zero-value `lexicalReach` silently.
- **Skip reach entirely** when Phase 3 triage class ∈ {docs-only, config-only, mechanical} —
  reach on a README-only PR is wasted push latency.

### 4.2 Sequencing change in `GitHubPullRequestBodyFromArtifact`

`buildPRBodyCatalog` currently computes `catalog.Stats` via `bodyStats` before reach exists.
Update the orchestrator (`pr_body.go:106`) to this exact order:

```go
func GitHubPullRequestBodyFromArtifact(ctx context.Context, artifact reviewbundle.Artifact) (string, error) {
    catalog := buildPRBodyCatalog(artifact)          // stats from patches/warnings only
    reach := computeLexicalReach(ctx, artifact.Bundle.Repo.RootPath, catalog, artifact)
    applyLexicalReachToStats(&catalog.Stats, reach) // mutates MaxRiskLevel/MaxRiskScore/RiskSignals
    summaryContext, err := collectPRSummaryContext(ctx, artifact, catalog)
    // ...
}
```

Add `applyLexicalReachToStats(stats *prBodyStats, reach lexicalReach)` in `pr_reach.go`:

- Raise to at least "medium" when `reach.DependentFiles >= 3`.
- Raise to "high" when `>= 10` or any symbol has 25+ external references.
- Append a risk signal like `lexical_reach:N_refs_in_M_files` to `stats.RiskSignals`.
- **Remove the raw-size fallback** in `bodyStats` (`pr_body.go:914-921`) as the sole cause
  of "high": size alone now caps at "medium" (a 500-line docs rewrite must not be "high").

Do not recompute the entire catalog; mutate stats in place after reach.

### 4.3 Rendering (honest labeling)

- `readinessSentence` (`pr_body.go:179`) becomes concrete and honestly worded — references,
  not dependents or call sites:
  `Blast radius is medium: 14 references to changed symbols (PublishArtifact, bodyStats)
  across 6 other files; 4 file(s), +120/-38 lines.`
  When reach is empty or disabled, keep the current file/line phrasing.
- Add a `## Blast Radius` detail block (only when `ReferenceCount > 0`, max 5 lines): each top
  symbol with its reference count and 1–2 blob-permalink links (Phase 2 helper) to the most
  significant external references. If `Truncated`, end the block with
  `Reach is lexical (text search) and truncated; common names were skipped.`
- Add reach facts to the AI brief: append a
  `codereview.ContextSnippet{Kind: "lexical_reach", ...}` per top symbol in `prReviewBrief`
  (`pr_body.go:398`).

Tests: fixture repo (`t.TempDir()` git repo) where a changed symbol is referenced from an
unchanged file → reach counted, permalink rendered, `applyLexicalReachToStats` escalates
risk; generic short symbol skipped; env gate off → sentence falls back; docs-only triage →
reach skipped.

## Phase 5 — LLM-written overview grounded in sessions

**Goal:** the opening paragraph explains what changed and *why* (from linked agent-session
transcripts — context competitors don't have), instead of concatenated revision descriptions.

The contract work already happened in Phase 1 (profile-guarded `overview` field, tolerant
parser, optional `AIReviewerWithOverview` interface). This phase is consumption only:

- `reviewPRSummaryFindings` already returns the overview as of Phase 3.2; wire it through
  `GitHubPullRequestBodyFromArtifact` into rendering.
- Rendering: `openingSummary` (`pr_body.go:152`) becomes: sanitized overview paragraph
  (`sanitizePRVisibleText`, trimmed to ~500 chars, strip any URLs defensively) followed by
  the deterministic verdict/blast-radius sentences. When the AI pass failed or the overview
  is empty, fall back to the current `changeSummarySentence` — behavior with AI off is
  unchanged.
- Session transcripts are already in the context (`sessionPRContextSnippets`,
  `pr_body.go:493`, capped at 4 snippets); no new plumbing needed.
- `gx review` must be unaffected: its briefs never use the `pr_summary` profile, so the model
  omits `overview` for them, and its call sites use plain `Review`.

Tests: fake reviewer implementing `AIReviewerWithOverview` returns an overview → rendered
first; returns no overview → template fallback; overview containing a URL → URL stripped; a
`patch_focused` brief's prompt output path never renders an overview.

## Phase 6 — Generalize the hardcoded gx heuristics

**Goal:** path-risk heuristics and rubric wording work for any repo using gx.

- **REVIEW.md config.** Extend `ReviewPolicy` parsing (`internal/codereview/review_policy.go:27`)
  to recognize a `high-risk paths` section, e.g. lines of the form
  `risk-path: internal/auth/** — auth changes can leak or misuse credentials`. Parse into
  `RiskPaths []RiskPath{Glob, Message string}` on `ReviewPolicy`. Follow the existing
  parsing style used for model hints (`review_policy.go:143-158`).
- **`hunkImpact` rewrite** (`pr_body.go:359`): match configured `RiskPaths` first (score 700,
  title/detail derived from the configured message). Ship a small generic fallback set —
  path/name patterns for `auth|token|secret|credential`, `migration|schema`,
  `Dockerfile|.github/workflows|deploy`, lockfiles — replacing the gx-specific
  `internal/storage/`, `internal/github/`, `internal/publication/`, `internal/vcs/` branches.
  Move today's gx-specific entries into this repo's own `REVIEW.md` so gx's PRs keep their
  current quality (dogfooding the config).
- The publication package needs the policy: call `codereview.LoadReviewPolicy(ctx, root)` in
  `GitHubPullRequestBodyFromArtifact` and thread it to `hunkImpact` and `prReviewBrief`.
- **De-gx-ify the rubric** (`pr_body.go:420-434`): "GX's behavior, state model, remote side
  effects" → "the project's behavior, persistence/state, remote side effects". Keep the
  reject list.
- **Kill canned repetition:** heuristic item details must mention the matched file/path
  (e.g. include `path.Base(hunk.File)` and the matched glob) so two different PRs don't
  render byte-identical warnings.

Tests: repo with `REVIEW.md` risk-paths → configured item rendered with configured message;
repo without → generic fallback only; gx's own `REVIEW.md` gains the migrated entries.

## Phase 7 — Provenance footer and retry

**Goal:** readers can tell how the summary was produced; transient AI failures don't silently
degrade quality.

- **Retry:** in `reviewPRSummaryFindings` (`pr_body.go:383`), retry the reviewer call once
  after a 2-second sleep on error (respect `ctx` cancellation). Keep the final
  swallow-and-continue behavior; `aiSucceeded` is true only if some attempt returned a
  parseable response (per Phase 3.2 — zero findings still counts as success).
- **Provenance footer:** last line of the GX-owned section, italic, e.g.
  `*Generated by GX — model gpt-4.1-mini; context: codebase, session, indexed code/session.*`
  or, when the AI pass failed: `*Generated by GX — heuristics only (AI unavailable).*`
  Model names come from `ReviewerInfo` returned by `prSummaryReviewerFromEnvWithInfo()`
  (Phase 1.4). When `info.Models` is nil/empty, omit the model clause. Context source names
  already exist (`summaryContext.Sources`).

**Deferred (explicitly out of scope for v1): model escalation.** Routing high-risk PRs to a
stronger model adds constructor surface area and model-policy questions that don't need
answering yet. Provenance + retry is enough for v1; revisit after the Phase 8 corpus can
measure whether the default model's findings are the weak point.

Tests: fake reviewer (via `prSummaryReviewerFromEnvWithInfo`) failing once then succeeding →
findings present (possibly zero) and `aiSucceeded` true; footer reflects heuristics-only
when all attempts fail; footer lists models when fake returns
`ReviewerInfo{Models: []string{"test-model"}}`.

## Phase 8 — Eval harness (quality regression net)

**Goal:** a golden corpus that locks in verdicts, items, links, and attribution so Phases 1–7
(and future prompt changes) can be validated mechanically.

- New: `internal/publication/testdata/prsummary/<case>/` fixtures. Each case: an artifact
  JSON (build `reviewbundle.Artifact` values the way `github_pr_test.go` does), an optional
  canned AI response JSON, and `expected.md` (the golden body).
- New test `TestPRSummaryGolden` in `internal/publication/pr_summary_golden_test.go`:
  for each case, override `prSummaryReviewerFromEnvWithInfo` (`pr_body.go:37`) to return a
  fake reviewer that replays the canned response (or errors for heuristics-only cases) plus
  optional `ReviewerInfo`. Fake should implement `AIReviewerWithOverview` when the case
  includes an overview. Run `GitHubPullRequestBodyFromArtifact`, diff against
  `expected.md`. Support `-update` flag to regenerate goldens.
- Minimum corpus (≈8 cases): docs-only (`No review needed`), docs-only with AI failing
  (`Quick scan`), tests-only, single risky code change with valid anchors, finding with an
  unresolvable file (no link), high-reach change (blast radius section with the lexical
  disclosure), REVIEW.md risk-path repo, external review-resource attribution (renders
  opaque).
- Deterministic only — no network, no real LLM. Live-model quality tracking is out of scope;
  leave a `TODO` in the test file pointing at GX Cloud review history as the future data
  source.

---

## Phase order and dependencies

```
Phase 1 (shared codereview contract: file/line/source_labels→ResolvedSource,
         optional overview interface, empty-response success, PR brief source wiring)
   │
   ├──► Phase 2 (PR rendering: validated anchors + honest attribution via ResolvedSources)
   │        │
   │        └──► Phase 3 (verdict via TriageChange; aiSucceeded; tests-only split)
   │                 │
   │                 ├──► Phase 4 (lexical reach — applyLexicalReachToStats after catalog)
   │                 ├──► Phase 5 (overview consumption — optional interface from Phase 1)
   │                 └──► Phase 7 (provenance/retry — needs aiSucceeded)
   ├──► Phase 6 (configurable risk paths — independent)
   └──► Phase 8 (golden evals — last; locks everything in)
```

Ship order: 1, 2, 3, 4, 5, 6, 7, 8. Each phase is one GX revision (`gx generate` after the
tests pass, then `gx push`).
