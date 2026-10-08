# Agents

Version control: plain Git. Once `gx init` installs the hooks, gx records your work automatically — there is no gx save verb.

Default flow:
- Run `git add` to stage the files for this revision.
- Run `git commit -m "..."` to save. A gx hook records the commit as a reviewable revision.
- Run plain `git push` to publish. The gx pre-push hook captures the coding session behind the change — do not run `gx capture push` yourself; it bypasses the hook.
- Open PRs with `gh pr create` (or the GitHub UI).
- To amend, use `git commit --amend` and preserve the gx revision trailer in the message.

For AI review, run the `gx_review` MCP tool (or the `gx review` CLI) on the current change.

## Installing: `just install`, never a bare `go build` onto PATH

Install with `just install`. Never write to `~/.local/bin/gx` (or anywhere on PATH)
with `go build -o` or `go install`: `just install` stamps the version, installs
completions, codesigns on macOS, and runs `just verify-bake` against both the built
and the installed binary.

There is no gx Cloud, so the default is a local build: with no endpoints in `.env`,
gx runs entirely on the machine and its cloud code stays dormant. `verify-bake`
matters only once endpoints are set. `GX_CLOUD_URL`, `GITHUB_CLIENT_ID` and
`CONVEX_SITE_URL` are all or nothing, and it refuses a binary that carries only some
of them, because that failure is silent. A build meant for a cloud but missing its
URL behaves exactly like a local one: pushes look normal, nothing is published, and
nothing says so. While gx Cloud existed, a bare build cost a week of PR summaries
in August 2026 that way. `gx version` printing `dev` instead of a commit sha is the
tell for an unstamped build.

`just` reads `.env` from parent directories too, so every worktree under
`.claude/worktrees` builds with the main checkout's `.env`.

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
