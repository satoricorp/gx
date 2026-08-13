# Agents

Version control: plain Git. Once `gx init` installs the hooks, gx records and publishes automatically — there is no gx save verb.

Default flow:
- Run `git add` to stage the files for this revision.
- Run `git commit -m "..."` to save. A gx hook records the commit as a reviewable revision.
- Run plain `git push` to publish. The gx pre-push hook captures the session and publishes code changes, sessions, and gx metadata to gx Cloud automatically — do not run `gx push` or `gx capture push` yourself; they bypass the hook.
- Open PRs with `gh pr create` (or the GitHub UI). Do not seed a `## Summary` in the PR body — leave human notes only; gx Cloud appends the rich summary below once the PR exists.
- To amend, use `git commit --amend` and preserve the gx revision trailer in the message.

To enhance the current change — AI-reported issues and tips to improve it — run the `gx_enhance` MCP tool (or the `gx enhance` CLI).

## Tests: seed through the writer production uses, or do not seed

Two production defects shipped past a fully green suite because every test built
a pristine, freshly migrated database in a temp dir, which never resembles a real
install. Both rules below exist because of that.

- **Never seed through a function production does not call.** `store.WriteSession`
  had zero production callers and fourteen test call sites while the only real
  writer of `sessions` could not bootstrap its first row, so 171 of 174 published
  bundles shipped `sessions: []` with the suite green. `TestSeedersAreProductionWriters`
  (internal/storage/storagetest) fails on any exported storage method with test
  callers and no production ones; its allowlist in
  `internal/storage/storagetest/testdata/test_only_writers.txt` is a ratchet —
  entries may be removed, never added.
- **Prefer a weathered database to a pristine one.** Use
  `storagetest.New(t, shapes...)` / `NewInWorld` with shapes such as
  `DriftedRepoIdentity`, `WeatheredNeighbourRepos`, `LegacyCaptureSessions` and
  `FossilCommitSelfReportSession` instead of hand-rolling
  `t.Setenv("GX_HOME", t.TempDir())`. Repositories and worktrees come from
  `gxtest.NewWorld(t)`. Cold start is also a real case, so keep pristine tests —
  just mark the premise with `storagetest.NoSessions()`.
