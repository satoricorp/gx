# Agents

Version control: plain Git. Once `lgtm init` installs the hooks, lgtm records and publishes automatically — there is no lgtm save verb.

Default flow:
- Run `git add` to stage the files for this revision.
- Run `git commit -m "..."` to save. A lgtm hook records the commit as a reviewable revision.
- Run plain `git push` to publish. The lgtm pre-push hook captures the session and publishes code changes, sessions, and lgtm metadata to lgtm Cloud automatically — do not run `lgtm push` or `lgtm capture push` yourself; they bypass the hook.
- Open PRs with `gh pr create` (or the GitHub UI). Do not seed a `## Summary` in the PR body — leave human notes only; lgtm Cloud appends the rich summary below once the PR exists.
- To amend, use `git commit --amend` and preserve the lgtm revision trailer in the message.

For AI review, run the `lgtm_review` MCP tool (or the `lgtm review` CLI) on the current change.

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
  `t.Setenv("LGTM_HOME", t.TempDir())`. Repositories and worktrees come from
  `lgtmtest.NewWorld(t)`. Cold start is also a real case, so keep pristine tests —
  just mark the premise with `storagetest.NoSessions()`.
