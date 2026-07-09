# GX Test Review

This is a design- and architecture-oriented review of the current test suite.
The suite spans CLI, authoring, daemon, cloud, storage, VCS, review bundle,
provider parsing, Cursor ingest, MCP tools, and e2e packages.

## Findings

- **Strongest coverage is on the authoring substrate.** Demux proposal
  normalization, hunk coverage, feasibility warnings, repair hints, provenance
  attachment, JJ-backed revision recording, and e2e demux apply paths are well
  covered.
- **The current tests protect the removed session-to-stack model.** Several
  VCS tests explicitly assert that current stack resolution ignores Cursor/latest
  request routing and does not route by session.
- **The CLI UX is lightly tested.** There are focused tests for summary output,
  demux warnings, help safety, and destructive confirmation, but not a broad
  snapshot suite for the public command output designs.
- **The MCP layer now has focused tool tests, but coverage is still thin.**
  Current tests cover compose, review, and publish command envelopes. Broader
  coverage is still needed for repair-loop output and error formatting.
- **The daemon is covered at primitive level, not product-flow level.** Session
  registry, resolver reuse, redaction, and WebSocket summary extraction are
  tested, but there is no end-to-end "captured session becomes revision
  provenance through daemon lookup" test that exercises the daemon boundary.
- **Multi-channel demux is not covered.** Current tests assume one current
  proposed revision / current stack. There are no tests for demux graph
  components, parallel branches, fan-out, SCCs, textual separability, or
  multi-stack input.
- **Review bundle coverage is valuable but still schema-centric.** The bundle
  tests lock JSON shape, patches, sessions, demux evidence, and basic risk
  computation. They do not yet validate review-surface UX states.

## Package Overview

### `cmd/gx`

- `TestShouldLaunch`: verifies command-line agent launcher detection.
- `TestRootExposesAuthCommand`: ensures `gx auth` is visible from the root.
- `TestRootExposesShortcutCommands`: ensures shortcut aliases such as status/add
  shortcuts are registered.
- `TestRootDoesNotExposeRemovedLegacyCommands`: protects the simplified CLI by
  keeping removed legacy commands out of root help.
- `TestRootDoesNotExposeCompletionCommand`: ensures shell completion is not a
  visible CLI command.
- `TestRootDoesNotExposeDaemonCommand`: ensures internal daemon mode remains
  hidden.
- `TestRootDoesNotExposeGitOrJJPassthroughCommands`: ensures GX does not become
  a generic Git/JJ passthrough CLI.
- `TestRootPrintsInitNoteWhenIdentityMissing`: checks root execution guides an
  unconfigured user toward init/identity setup.
- `TestRootSkipsInitNoteWhenIdentityConfigured`: checks configured users do not
  see unnecessary init guidance.

### `internal/authoring`

- `TestGroupFilesPairsSourceAndTests`: groups source files with likely test
  counterparts.
- `TestProposalIntentUsesPrefix`: creates demux revision intents from an
  explicit prefix.
- `TestCleanFilesSortsAndDedupes`: normalizes file lists deterministically.
- `TestOrderGroupsByStructuralDependencies`: orders file groups using
  structural dependency edges.
- `TestAnnotateHunksWithSymbols`: attaches inferred enclosing symbols to hunks.
- `TestExpandGroupsByChangedSymbolsSplitsSingleFileBySymbol`: splits a single
  file into hunk-level revisions when hunks map cleanly to different symbols.
- `TestExpandGroupsByChangedSymbolsFallsBackWhenAnyHunkHasNoSymbol`: avoids
  unsafe symbol splitting when one hunk cannot be mapped.
- `TestNormalizeDemuxProposalResolvesHunkIDs`: expands `hunk_ids` into full
  hunk payloads and derives files.
- `TestNormalizeDemuxProposalRejectsDuplicateHunkIDs`: rejects plans assigning
  one hunk to multiple revisions.
- `TestValidateDemuxPlanReturnsRepairableErrors`: returns structured repair
  hints for invalid hunk coverage.
- `TestValidateDemuxPlanReturnsNormalizedWarnings`: validates a good plan while
  returning inferred dependency hints.
- `TestValidateDemuxPlanReturnsReorderRepairHint`: returns a typed reorder hint
  when structural order is wrong.
- `TestDemuxWorkflowStateClassifiesReadyRecommendedAndRequired`: classifies
  demux workflow states for ready, recommended repair, and required repair.
- `TestDemuxWorkflowPacketUsesValidatedProposal`: ensures the machine demux
  packet carries the normalized/validated proposal.
- `TestBlockingFeasibilityWarningsOnlyBlocksWarningSeverity`: gates apply on
  warning-severity feasibility warnings, not info warnings.
- `TestNormalizeDemuxProposalRejectsUnassignedHunk`: rejects plans that leave a
  top-level hunk uncovered.
- `TestNormalizeDemuxProposalAllowsWholeFileCoverage`: allows one whole-file
  revision to cover all hunks in that file.
- `TestNormalizeDemuxProposalAllowsWholeFileRevisionWithContextHunks`: allows
  whole-file plans even when hunk catalog context exists.
- `TestNormalizeDemuxProposalRejectsMixedHunkAndWholeFileCoverage`: rejects
  mixing hunk-level and whole-file coverage for the same file.
- `TestNormalizeDemuxProposalRejectsDuplicateRevisionIDs`: rejects duplicate
  revision identifiers.
- `TestNormalizeDemuxProposalRejectsLaterDependency`: rejects `depends_on`
  references to unknown or later revisions.
- `TestMergeDemuxPlanCarriesHunkCatalog`: preserves proposal catalog fields
  when an LLM plan only supplies revised revisions.
- `TestFeasibilityWarningsReportsUnmappedHunks`: warns when hunks cannot be
  mapped to enclosing symbols.
- `TestFeasibilityWarningsReportsStructuralOrderConflict`: warns when a
  consumer revision is ordered before a producer revision.
- `TestFeasibilityWarningsUsesSymbolOwnerForDependencyTarget`: finds the
  hunk-level producer revision rather than only file-level ownership.
- `TestFeasibilityWarningsUsesReferencingHunkAsDependencySource`: attributes a
  dependency warning to the hunk/revision that references the symbol.
- `TestFeasibilityWarningsReportsMissingInferredDependsOn`: suggests
  `depends_on` for an already-valid order with implicit dependency.
- `TestFeasibilityWarningsSkipsInferredDependencyWhenDependsOnIsExplicit`:
  suppresses inferred dependency hints after explicit `depends_on`.
- `TestFeasibilityWarningsForProposalUsesFinalPlanOrder`: recomputes warnings
  from the final LLM-authored order rather than stale proposal order.
- `TestFeasibilityWarningsReportsSeparatedTestCounterpart`: warns when likely
  source/test counterparts are split.

### `internal/cli`

- `TestRenderStaticLogoPlainWithoutColor`: verifies plain logo rendering when
  color is disabled.
- `TestPromptModifySelection`: verifies the CLI prompt path for selecting a
  revision to modify.
- `TestPrintAddSummaryIncludesHashesSplitAndEditCommands`: checks `gx commit` / `gx add`
  summary includes description, hashes, split status, and follow-up commands.
- `TestPrintDemuxProposalIncludesFeasibilityWarnings`: checks human demux output
  includes feasibility warning details.
- `TestPRHelpDoesNotPublish`: ensures `gx pr --help` is side-effect free.
- `TestConfirmBaseSwitchRequiresYes`: requires exact confirmation before a
  destructive base switch.

### `internal/cloud`

- `TestLoginDeviceFlowAndComplete`: exercises GitHub device login and token
  completion.
- `TestLoginPollsUntilAuthorized`: verifies polling handles pending auth before
  success.
- `TestLoginUsesBakedDefaults`: checks baked cloud/GitHub config is honored.
- `TestLoginNotConfiguredInDevBuild`: checks login fails cleanly without config.
- `TestLogoutRevokesAndClears`: verifies logout revokes and clears credentials.
- `TestRequestGitHubDeviceCodeForm`: verifies GitHub device-code form payload.
- `TestCloudBaseURLStripsGxPrSuffix`: normalizes cloud URLs used for stack
  metadata.
- `TestRepoFullNameFromRemoteURL`: parses owner/repo from remote URLs.
- Endpoint/env tests: verify GitHub client id, Convex site URL, cloud URL, and
  API key precedence across env, baked defaults, and unset/dev values.
- `TestCloudURLUsesDefaultWhenUnset`: uses local/default cloud URL.
- `TestCloudURLExplicitDisable`: supports disabling cloud via env.
- `TestCloudURLCustom`: supports custom cloud URL.
- `TestSyncPushBearerAndReviewURL`: verifies sync upload auth and review URL.
- `TestSyncPushRequiresTokenWhenCloudURLSet`: requires auth when cloud upload is
  configured.
- `TestNewClientUsesCloudURL`: creates a cloud client from configured URL.
- `TestNewClientEnvOverridesBakedDefault`: env overrides baked client defaults.
- `TestNewClientNilWhenUnset`: no client is created when cloud is unset.
- `TestDefaultMachineIDCreatesOnce`: persists a stable local machine id.
- `TestSaveLoadClearCloudCredentials`: round-trips stored cloud credentials.
- `TestCloudAPITokenRequiresGitHubToken`: requires GitHub auth for cloud API calls.

### `internal/daemon`

- `TestMemorySessionRegistry`: registers and resolves session ids by process id.
- `TestRedactHeaders`: redacts sensitive headers before persistence/logging.
- `TestAmbientResolverReusesExistingSession`: reuses ambient sessions for the
  same resolved context.
- `TestExtractSummaryFromWebSocketFrames`: extracts useful summary content from
  WebSocket frame traffic.

### `internal/ingest/cursor`

- `TestSyncFromPathIngestsComposersAndBubbles`: ingests Cursor composer/bubble
  records into GX session storage.
- `TestSyncFromPathIsIdempotent`: repeated Cursor ingest does not duplicate
  session data.
- `TestSyncFromPathMissingFile`: handles missing Cursor DB/file input.
- `TestIsMeaningfulFolder`: filters Cursor folders to meaningful repo paths.
- `TestParseBubbleKey`: parses Cursor bubble keys into stable identifiers.

### `internal/launcher`, `service`, `prompt`, `providers`, `postlist`

- `TestCodexCommandArgs`: verifies wrapped Codex launch args.
- `TestClaudeCommandArgsUnchanged`: ensures Claude args are passed through.
- Codex config tests: update/find TOML string values and repair provider config
  in the scoped Codex config format.
- Service manager tests: launchd plist uses ambient daemon, shell block uses
  stable proxy URLs, and shell block insertion is idempotent.
- Prompt tests: legacy/plain prompt behavior for modify selection, identity, and
  required default values.
- Provider tests: assemble Anthropic/OpenAI request/response data and summarize
  OpenAI compaction usage.
- `TestUpsertIdentity`: posts identity data to the postlist endpoint.

### `internal/provenance`

- `TestParseSessionIDsDedupesAndSplitsCommonSeparators`: parses session env
  values robustly.
- `TestResolvePrefersExplicitSessionIDs`: explicit session ids beat fallback.
- `TestResolveReturnsRepoLocalSessionWhenNoExplicitEnv`: repo-local captured
  sessions are used when env is absent.
- `TestAttachWritesResolvedSessions`: writes resolved sessions to
  `change_sessions`.
- `TestWithExplicitSessionEnvRestoresPreviousEnvironment`: temporary explicit
  provenance env restoration is safe.

### `internal/reviewbundle`

- `TestBuildPushIncludesDemuxEvidence`: includes per-revision demux evidence in
  push bundles.
- `TestBuildPushStackIncludesPatchesAndDedupedSessions`: includes patches and
  deduplicated linked sessions.
- `TestBuildPushJSONShapeIsStableForConsoleIngest`: golden-style test for
  console ingest JSON shape.
- `TestBuildReviewContextComputesRiskFromEvidence`: computes basic review risk
  from bundle evidence.

### `internal/storage`

- `TestOpenRepairsCorruptedJJChangeIDs`: repairs older/corrupted JJ change ids
  on open.
- `TestDemuxProposalRoundTrip`: persists and loads demux proposal JSON/status.
- `TestChangeDemuxEvidenceRoundTrip`: persists and reads per-revision demux
  evidence.
- `TestFindAttachableSessionsForRepoReturnsUnlinkedRepoSessions`: finds
  unlinked captured sessions matching the current repo.

### `internal/structural`

- `TestAnalyzeFindsGoDefinitionsAndDependencies`: extracts Go definitions and
  simple dependency edges.
- `TestEnclosingSymbolsMapsHunksToNearestSymbol`: maps changed hunks to nearest
  enclosing symbols.
- `TestEnclosingSymbolsPrefersDefinitionStartingInsideHunkContext`: handles
  hunks whose context includes symbol starts.
- `TestAnalyzeSkipsUnsupportedFiles`: ignores unsupported language/file types.

### `internal/termstyle`

- Color/theme tests verify `NO_COLOR`, default color enablement, plain
  label/value output, status prefixes, empty stack rendering, doctor table
  rendering, and hyperlink fallback without color.

### `internal/vcs`

- Stack tests verify GX stack naming, fallback slugs, JJ stack revsets,
  detached-HEAD recovery with JJ stack refs, git index lock detection, and retry
  behavior.
- Bookmark tests verify edit revision selection for published/draft stacks and
  JJ immutable error detection.
- Description tests validate revision messages and require recorded revisions
  before publish.
- Service tests cover branch push args, modify selection prompts, GitHub remote
  parsing, PR URLs, publish ref selection, stack forking, stack
  resolution from base, ignoring Cursor/session routing for stack selection,
  repo initialization/reuse, identity persistence, revision recording,
  hunk/file split commits, push rejection for empty stacks, metadata attachment,
  repo-local session attachment, prompt defaults, sync defaults, SQLite busy
  retry, and JJ stderr-noise tolerance.

### `test/e2e`

- `TestGXAddMaintainsBookmarkStateAndAttachesExplicitSession`: records a
  revision while preserving stack state and explicit session linkage.
- `TestGXDemuxProposesAndAppliesFileLevelRevisions`: proposes file-level demux
  revisions, applies them, and verifies evidence/provenance.
- `TestGXDemuxChangesWorkflowHappyPathAppliesValidatedProposal`: exercises the
  machine demux workflow packet, validation, and apply flow.
- `TestGXDemuxAttachesRepoLocalSessionWithoutEnv`: proves ambient repo-local
  sessions can attach without explicit env.
- `TestGXDemuxValidatePlanReportsRepairableErrorsWithoutApplying`: validates
  bad plans, returns repairable errors, and avoids JJ side effects.
- `TestGXDemuxApplyPlanResolvesHunkIDs`: applies an LLM-style plan using
  catalog `hunk_ids`.
- `TestGXDemuxProposesSymbolLevelRevisions`: proves deterministic symbol-level
  demux can split one file into multiple revisions.
- `TestGXPRWithMultipleBookmarksPushesAPIAndPreservesState`: pushes stack data,
  uploads API context, and preserves local state across multiple stacks.
- `TestGXSwitchBaseRequiresConfirmationAndReturnsToMain`: protects destructive
  base switch behavior.
- `TestGXSwitchAcceptsImplicitStackAliases`: supports user-friendly stack aliases.

## UX/Test Recommendations

1. Add CLI snapshot tests for the final public design of `gx status`, `gx commit`,
   `gx demux`, `gx demux apply`, `gx pr`, and daemon/capture status.
2. Add MCP tool tests for argument construction and repair-loop behavior.
3. Add a daemon-to-authoring integration test proving active session lookup can
   attach provenance without env variables.
4. Add demux graph tests before multi-channel design: SCC grouping,
   weak-component fan-out, topological ordering, textual overlap edges, and
   separability failure.
5. Add review-surface bundle tests for the UX states: ready, needs repair,
   blocked, accepted-with-warnings, and demux evidence display.
