# Demux Rebuild Plan

## Purpose

This document preserves the product behavior, data contracts, type contracts,
and executable specification of GX's demux/generate implementation before its
jj-dependent code is removed.

Use it to rebuild the feature only if GX again needs to turn a mixed working
tree into ordered, reviewable Revisions. A future implementation should be
Git-native and should not restore the old jj checkout, bookmark, or operation
machinery.

Historical source baseline:

```text
Repository: satoricorp/gx
Commit: 1abd91d
User command at that commit: gx generate
Internal names: demux, compose
```

Removed source can be inspected without restoring it:

```bash
git show 1abd91d:internal/authoring/proposal.go
git show 1abd91d:internal/authoring/demux.go
git show 1abd91d:internal/authoring/generate_pipeline.go
```

## Product Contract

Given a source snapshot containing mixed changes, generate must:

1. Build an addressable catalog of every changed hunk.
2. Group those hunks into ordered Revision proposals.
3. Pair source and test changes when they form one reviewable concern.
4. Infer dependencies and order prerequisites before dependants.
5. Route Revisions to existing or proposed branches when requested.
6. preserve exact session provenance as evidence, never use sessions as routing
   authority.
7. Validate that every changed hunk is owned exactly once.
8. Return repairable, structured diagnostics instead of partially applying an
   invalid plan.
9. Preflight an accepted plan in disposable isolation before mutating the real
   checkout.
10. Apply the complete accepted plan atomically, or restore the source state.
11. Record planning and provenance evidence against each resulting GX Revision.
12. Leave unaccepted changes available for a later planning run.

The valuable module is the plan graph:

```text
SourceSnapshot
  -> HunkCatalog
  -> Ordered RevisionPlan[]
  -> PlanReview
  -> RepairHints
  -> Preflight
  -> Applied GX Revisions
  -> Review evidence
```

## Behavior Worth Preserving

- Exact-once hunk coverage.
- Source/test and manifest/lockfile pairing.
- Symbol-level splitting only when every affected hunk maps safely.
- Structural dependency ordering.
- Existing-branch matching before proposing new branches.
- Conventional branch kinds such as `feature/`, `bug/`, `docs/`, `test/`, and
  `chore/`.
- Explicit user routes override inferred routes.
- Deterministic repair before model-assisted repair.
- A conservative whole-file or whole-plan fallback when a precise plan cannot
  be proven safe.
- Blocking warnings are distinct from informational diagnostics.
- Session evidence informs grouping but does not choose branches.
- Plan and apply are separate operations, even if one adapter offers a
  plan-and-apply convenience command.

## Behavior Not To Recreate

- jj `@` and `@-` as source state.
- jj change IDs as plan anchors.
- `jj edit`, `jj restore`, `jj new`, bookmarks, or operation rollback.
- Moving the user's visible checkout to an internal base/edit state.
- Separate `compose`, `demux`, and `generate` product surfaces.
- Stringly typed repair logic that parses human error messages when a typed
  diagnostic is available.
- Hidden partial application followed by "run generate again."
- Packet fields naming MCP tools that do not exist.

## Core Data Contracts

The future public names below intentionally replace the old `Demux*` names.
Legacy JSON aliases may be accepted during migration.

### Source Snapshot

The plan must be bound to immutable source state so stale plans cannot apply.

```go
type SourceSnapshot struct {
    RepositoryID string `json:"repository_id"`
    WorktreeRoot string `json:"worktree_root"`
    BaseCommitID string `json:"base_commit_id"`
    TreeID       string `json:"tree_id"`
    IndexTreeID  string `json:"index_tree_id,omitempty"`
    CreatedAt    int64  `json:"created_at"`
}
```

`TreeID` is the hash of the complete planned content. `IndexTreeID` is present
when the plan is intentionally scoped to staged changes. Apply must reject a
plan when its source hash no longer matches.

Legacy equivalents:

```text
proposed_change_id -> no equivalent; remove
proposed_commit_id -> source_snapshot.tree_id or base_commit_id
```

### Hunk Catalog

```go
type Hunk struct {
    ID         string `json:"id"`
    File       string `json:"file"`
    Header     string `json:"header"`
    OldStart   int    `json:"old_start"`
    OldLines   int    `json:"old_lines"`
    NewStart   int    `json:"new_start"`
    NewLines   int    `json:"new_lines"`
    Symbol     string `json:"symbol,omitempty"`
    SymbolKind string `json:"symbol_kind,omitempty"`
    Patch      string `json:"patch,omitempty"`
}
```

Requirements:

- IDs are unique within one plan.
- IDs are opaque to consumers.
- The patch is canonical unified-diff text.
- Model requests may omit `patch`, but apply must resolve IDs against the
  server-owned catalog.
- A hunk cannot be copied into a plan and changed by a model.

### Revision Proposal

```go
type RevisionProposal struct {
    ID                string               `json:"id"`
    Intent            string               `json:"intent"`
    Files             []string             `json:"files,omitempty"`
    UseHunks          bool                 `json:"use_hunks,omitempty"`
    HunkIDs           []string             `json:"hunk_ids,omitempty"`
    DependsOn         []string             `json:"depends_on,omitempty"`
    TargetBranch      string               `json:"target_branch,omitempty"`
    BaseBranch        string               `json:"base_branch,omitempty"`
    RouteReason       string               `json:"route_reason,omitempty"`
    RouteConfidence   float64              `json:"route_confidence,omitempty"`
    RouteSource       RouteSource          `json:"route_source,omitempty"`
    ProvenanceStatus  ProvenanceStatus     `json:"provenance_status,omitempty"`
    SessionIDs        []string             `json:"session_ids,omitempty"`
    SessionContexts   []SessionContext     `json:"session_contexts,omitempty"`
    HunkLinks         []HunkLink           `json:"hunk_links,omitempty"`
    Confidence        float64              `json:"confidence,omitempty"`
    EffectiveLOC      int                  `json:"effective_loc,omitempty"`
    ShapeReasons      []string             `json:"shape_reasons,omitempty"`
    SemanticLabels    []SemanticLabelMatch `json:"semantic_labels,omitempty"`
}
```

Rules:

- `ID` is stable only within the plan.
- `Intent` becomes the proposed commit message.
- Whole-file mode uses `Files`.
- Hunk mode sets `UseHunks` and references `HunkIDs`.
- Do not mix whole-file and hunk ownership for the same file.
- `DependsOn` forms an acyclic graph and references only plan Revision IDs.
- The final list is topologically ordered.
- `RouteSource` is `heuristic`, `model`, or `user`.
- User-authored routes cannot be overwritten by heuristics or repair.

### Plan

```go
type RevisionPlan struct {
    ID                   string                `json:"id"`
    Status               PlanStatus            `json:"status"`
    Source               SourceSnapshot        `json:"source"`
    Hunks                []Hunk                `json:"hunks"`
    Revisions            []RevisionProposal    `json:"revisions"`
    StructuralFacts      []StructuralFact      `json:"structural_facts,omitempty"`
    StructuralDeps       []StructuralDependency `json:"structural_dependencies,omitempty"`
    ChangedSymbols       []ChangedSymbol       `json:"changed_symbols,omitempty"`
    FeasibilityWarnings  []FeasibilityWarning  `json:"feasibility_warnings,omitempty"`
    Warnings             []string              `json:"warnings,omitempty"`
    Confidence           PlanConfidence        `json:"confidence,omitempty"`
    AppliedRevisionIDs   []string              `json:"applied_revision_ids,omitempty"`
    CreatedAt            int64                 `json:"created_at"`
    UpdatedAt            int64                 `json:"updated_at"`
}
```

```go
type PlanStatus string

const (
    PlanPending          PlanStatus = "pending"
    PlanPartiallyApplied PlanStatus = "partially_applied"
    PlanApplied          PlanStatus = "applied"
    PlanExpired          PlanStatus = "expired"
)
```

### Review Result And Repair Hints

```go
type PlanReview struct {
    Valid       bool         `json:"valid"`
    State       WorkflowState `json:"state"`
    Plan        RevisionPlan `json:"plan"`
    Diagnostics []Diagnostic `json:"diagnostics,omitempty"`
    RepairHints []RepairHint `json:"repair_hints,omitempty"`
}

type WorkflowState string

const (
    ReadyToApply  WorkflowState = "ready_to_apply"
    RepairNeeded  WorkflowState = "repair_required"
)
```

Do not restore the unused legacy `repair_recommended` state without a concrete
consumer and transition rule.

```go
type RepairHint struct {
    Kind         RepairKind `json:"kind"`
    RevisionID   string     `json:"revision_id,omitempty"`
    HunkID       string     `json:"hunk_id,omitempty"`
    File         string     `json:"file,omitempty"`
    DependsOn    string     `json:"depends_on,omitempty"`
    Symbol       string     `json:"symbol,omitempty"`
    TargetBranch string     `json:"target_branch,omitempty"`
    BaseBranch   string     `json:"base_branch,omitempty"`
    Candidates   []string   `json:"candidates,omitempty"`
    Suggestion   string     `json:"suggestion"`
}
```

Required repair kinds:

```text
unassigned_hunk
duplicate_hunk_assignment
mixed_hunk_and_whole_file
unknown_hunk
patch_mismatch
invalid_dependency
reorder_dependency
invalid_route
cross_branch_dependency
apply_preflight
low_confidence
```

Diagnostics should be typed:

```go
type Diagnostic struct {
    Code       string   `json:"code"`
    Severity   Severity `json:"severity"`
    Message    string   `json:"message"`
    RevisionID string   `json:"revision_id,omitempty"`
    HunkID     string   `json:"hunk_id,omitempty"`
    File       string   `json:"file,omitempty"`
}

type Severity string

const (
    SeverityInfo    Severity = "info"
    SeverityWarning Severity = "warning"
    SeverityError   Severity = "error"
)
```

Errors and warning-severity feasibility diagnostics block apply. Information
diagnostics do not.

### Confidence

```go
type PlanConfidence struct {
    LogicConfidence     float64            `json:"logic_confidence"`
    ModelConfidence     float64            `json:"model_confidence,omitempty"`
    EffectiveConfidence float64            `json:"effective_confidence"`
    LogicReasons        []ConfidenceReason `json:"logic_reasons,omitempty"`
    ModelReasons        []ConfidenceReason `json:"model_reasons,omitempty"`
}

type ConfidenceReason struct {
    Kind       string  `json:"kind"`
    Severity   string  `json:"severity"`
    Message    string  `json:"message"`
    Suggestion string  `json:"suggestion,omitempty"`
    Delta      float64 `json:"delta,omitempty"`
}
```

The old implementation used `0.50` as medium confidence and `0.85` as high
confidence. Treat those as starting calibration values, not permanent product
contracts.

### Apply Result

```go
type ApplyPlanResult struct {
    Plan             RevisionPlan      `json:"plan"`
    Revisions        []AppliedRevision `json:"revisions"`
    RemainingChanges bool              `json:"remaining_changes,omitempty"`
    AcceptedSubset   bool              `json:"accepted_subset,omitempty"`
}

type AppliedRevision struct {
    ProposalID string `json:"proposal_id"`
    RevisionID string `json:"revision_id"`
    CommitID   string `json:"commit_id"`
    Branch     string `json:"branch"`
}
```

`RevisionID` is the durable GX-owned ID described by the Git-native Revision
ADR. `CommitID` is the current Git snapshot.

## Validation Invariants

Validation must be deterministic and side-effect free.

1. Plan Revision IDs are non-empty and unique.
2. Hunk IDs are non-empty and unique.
3. Every referenced hunk exists in the catalog.
4. Every catalog hunk has exactly one owner.
5. A file cannot have both whole-file and hunk-level owners.
6. Embedded patch data, if accepted at an adapter, exactly matches the catalog.
7. Whole-file owners are unique.
8. Dependencies reference known Revisions.
9. Dependencies form a DAG.
10. The Revision list is topologically ordered.
11. Explicit routes reference valid or intentionally new branches.
12. Cross-branch dependencies have a valid base relationship.
13. The source snapshot still matches before apply.
14. Warning-severity diagnostics require explicit acceptance.

Validation errors must produce diagnostics and repair hints. They must never
partially mutate Git state.

## Planning Pipeline

### Phase 1: Collect

- Read the active worktree and Git common-directory identity.
- Select staged, unstaged, or explicitly scoped changes.
- Compute the immutable source snapshot.
- Parse canonical unified diff into the hunk catalog.
- Collect structural facts, symbols, source/test relationships, and session
  evidence.

### Phase 2: Deterministic Plan

Apply these passes in order:

1. Pair source/tests and manifest/lockfile counterparts.
2. Group obvious project scaffolding and infrastructure concerns.
3. Order groups by structural dependencies.
4. Split by changed symbol only when all hunks map safely.
5. Refine groups with session affinity without violating hunk coverage.
6. Shape Revisions to avoid unexplained tiny slices and oversized reviews.
7. Match existing branches before proposing conventional new branches.
8. Infer dependencies and feasibility diagnostics.
9. Persist the plan before exposing it to an adapter.

### Phase 3: Review And Repair

1. Normalize adapter/model output against the server-owned catalog.
2. Validate all invariants.
3. Emit typed diagnostics and repair hints.
4. Apply deterministic repairs to a fixed point.
5. If still blocked and model repair is enabled, send only bounded plan context
   and diagnostics to the model.
6. Validate model output exactly like untrusted user input.
7. Fall back conservatively rather than accepting a low-confidence precise
   split.

Model response contract:

```json
{
  "revisions": [],
  "model_confidence": 0.0,
  "model_reasons": [],
  "notes": []
}
```

The model may propose Revision fields. It cannot modify the source snapshot,
hunk catalog, provenance records, or accepted user routes.

### Phase 4: Preflight

Run the real apply adapter in disposable isolation:

- Create a temporary Git worktree at the source commit.
- Recreate the source tree/index state.
- Apply the complete plan.
- Verify resulting commits contain exactly the planned patch.
- Verify dependency and branch tips.
- Verify no source changes are lost.
- Delete the worktree.

Preflight success is required before automatic apply. A user may explicitly
accept warning diagnostics, but cannot bypass structural validation.

### Phase 5: Apply

- Acquire one repository lock based on Git common directory.
- Recheck the source snapshot.
- Create Revisions in dependency order.
- Generate a durable GX Revision ID for each commit.
- Add the GX trailer before each Git commit is created.
- Attach exact provenance and plan evidence.
- Advance branch refs with compare-and-swap.
- On failure, restore refs and the active worktree/index exactly.
- Mark the plan applied only after every selected Revision is durable.

## Persistence Contracts

Future tables should use product language:

```sql
CREATE TABLE revision_plans (
    id TEXT PRIMARY KEY,
    repo_id INTEGER NOT NULL,
    source_tree_id TEXT NOT NULL,
    status TEXT NOT NULL,
    payload_json TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    applied_at INTEGER
);

CREATE TABLE revision_plan_evidence (
    revision_id TEXT PRIMARY KEY,
    plan_id TEXT NOT NULL,
    proposal_revision_id TEXT NOT NULL,
    intent TEXT NOT NULL,
    files_json TEXT NOT NULL,
    hunk_ids_json TEXT NOT NULL,
    confidence REAL NOT NULL,
    provenance_status TEXT NOT NULL,
    evidence_json TEXT NOT NULL,
    created_at INTEGER NOT NULL
);
```

Legacy mapping:

```text
demux_proposals             -> revision_plans
change_demux_evidence       -> revision_plan_evidence
base_change_id              -> source_tree_id
revision_proposal_id        -> proposal_revision_id
```

Evidence JSON should preserve:

- plan ID and proposal Revision ID;
- intent, files, and hunk IDs;
- confidence and reasons;
- provenance status and linked session IDs;
- structural facts and dependencies;
- changed symbols;
- feasibility diagnostics;
- route evidence;
- model participation and model identifier, when present.

## Review Bundle Contract

The review bundle should expose planning evidence without requiring the console
to understand planner internals:

```go
type RevisionPlanEvidencePayload struct {
    PlanID             string               `json:"plan_id"`
    ProposalRevisionID string               `json:"proposal_revision_id"`
    Intent             string               `json:"intent"`
    Files              []string             `json:"files"`
    HunkIDs            []string             `json:"hunk_ids"`
    Confidence         float64              `json:"confidence"`
    ProvenanceStatus   string               `json:"provenance_status"`
    Diagnostics        []Diagnostic         `json:"diagnostics,omitempty"`
    StructuralFacts    []StructuralFact     `json:"structural_facts,omitempty"`
    StructuralDeps     []StructuralDependency `json:"structural_dependencies,omitempty"`
    ChangedSymbols     []ChangedSymbol      `json:"changed_symbols,omitempty"`
}
```

During compatibility, schema-v1 `demux_evidence` may be populated from the same
record. A new schema should use `revision_plan_evidence`.

## Adapter Contract

A future CLI or MCP adapter should expose three operations:

```text
plan_changes
review_revision_plan
apply_revision_plan
```

The adapter may offer a convenience operation that calls all three, but the
underlying module remains addressable and resumable.

The planning response must include:

```json
{
  "state": "ready_to_apply",
  "plan": {},
  "review": {
    "valid": true,
    "diagnostics": [],
    "repair_hints": []
  }
}
```

Do not include names of follow-up tools in the data contract. Tool discovery
belongs to the adapter, not the plan.

## Reconstruction Module Design

Use a small number of deep modules:

```go
type Planner interface {
    Plan(context.Context, PlanInput) (RevisionPlan, error)
}

type Reviewer interface {
    Review(context.Context, RevisionPlan) (PlanReview, error)
}

type Repairer interface {
    Repair(context.Context, RevisionPlan, PlanReview) (RevisionPlan, error)
}

type Applier interface {
    Preflight(context.Context, RevisionPlan) (PreflightResult, error)
    Apply(context.Context, RevisionPlan, ApplyOptions) (ApplyPlanResult, error)
}
```

Do not create an interface until there are two adapters. Pure validation,
grouping, shaping, and routing should remain ordinary functions.

Recommended module layout:

```text
internal/planning/
  contract.go
  collect.go
  group.go
  route.go
  validate.go
  repair.go
  evidence.go
internal/planning/gitapply/
  preflight.go
  apply.go
```

## Executable Specification To Preserve

Historical tests at commit `1abd91d`:

```text
internal/authoring/demux_test.go
internal/authoring/demux_eval_test.go
internal/authoring/demux_patch_coverage_test.go
internal/authoring/demux_already_applied_test.go
internal/authoring/generate_pipeline_test.go
internal/authoring/compose_capture_test.go
internal/cli/demux_loader_test.go
internal/storage/db_test.go
internal/reviewbundle/bundle_test.go
```

Highest-value scenarios:

- Multi-concern grouping into code/test, infrastructure, tooling, and worker
  Revisions.
- Test-only changes routed together.
- Source/test and package/lock pairing.
- Symbol-level splitting with whole-file fallback for unmapped hunks.
- Duplicate, missing, mixed-mode, and mismatched-patch rejection.
- Dependency ordering and cross-branch dependency diagnostics.
- Explicit user route preservation.
- Session isolation and the invariant that sessions never select branches.
- Already-applied hunk detection.
- Complete multi-branch apply with no hidden remainder.
- Per-Revision evidence and session attachment in the review bundle.

Add missing scenarios before declaring a rebuild complete:

- Partial accept with an intact, re-plannable remainder.
- Invalid plan causes zero Git mutation.
- Successful disposable preflight.
- Explicit warning acceptance.
- Concurrent linked-worktree planning.
- Apply rollback after each mutation point.
- Model timeout, malformed output, and adversarial catalog references.
- Squash/rebase after apply preserves GX Revision identity and evidence.

## Rebuild Sequence

1. Copy the contracts in this document into backend-independent Go types and
   JSON Schema fixtures.
2. Port validation and coverage tests as pure property/table tests.
3. Build a data-driven grouping/routing evaluation corpus from historical
   tests.
4. Implement collection from Git diffs and immutable tree hashes.
5. Implement deterministic planning, review, and fixed-point repair.
6. Add the optional model repair adapter with bounded untrusted input.
7. Implement Git-worktree preflight.
8. Implement transactional Git apply with GX Revision trailers.
9. Persist plan evidence and project it into review bundles.
10. Add CLI/MCP adapters only after the module contract is stable.

## Acceptance Gate

A rebuild is ready only when:

- every source hunk is applied exactly once or explicitly left in the remainder;
- invalid or stale plans make no Git changes;
- a successful apply produces the same patch as the accepted plan;
- each resulting commit has a durable GX Revision ID;
- provenance and plan evidence are attached to the correct Revision;
- linked worktrees do not affect each other's HEAD, index, or files;
- automatic apply requires a successful disposable preflight;
- model output is treated as untrusted and cannot alter source evidence;
- review bundles render useful context without parsing opaque planner JSON;
- no runtime dependency on jj is introduced.
