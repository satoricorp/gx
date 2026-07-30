# Git lifecycle hooks

Totality installs `prepare-commit-msg`, `post-commit`, `post-rewrite`, and `pre-push`
hooks. Together they stamp revision trailers, record commits, follow rewritten
commit OIDs, and run capture before a push completes. Staged extracts and
sessions are written to `~/.totality/totality.db` without network access.

Installing the Totality CLI package provides `tx` and `tx-mcp`. It does not rewrite every repository immediately; each repo gets the
Git hook when it is initialized with `tx init`, initialized automatically by
the CLI or MCP.

## Native Git hooks

The `pre-push` hook invokes:

```bash
tx capture push --remote "$remote" --ref-range "$range" --repo "$(git rev-parse --show-toplevel)"
```

The ref range is `remote_sha..local_sha` for updates, or the local SHA for new branches.

The hook skips deleted refs, captures Claude/Codex/Cursor session context by
default, stages capture data locally when upload credentials are missing, and
uploads when credentials are configured.

## Machine-wide install (`tx init --global`)

`tx init --global` writes the same lifecycle scripts to `~/.totality/hooks`
(`$TOTALITY_HOME/hooks` when set) and points `git config --global core.hooksPath` at
that directory, so Totality works in every repo without a per-repo `tx init`. It does
not have to run inside a repository.

A global `core.hooksPath` makes Git ignore every repository's own
`.git/hooks/*`, so the global scripts shim more than Totality's four hooks: each
script does Totality's work (if any) and then executes the repository's own hook of
the same name, forwarding arguments, stdin, and exit status. `pre-push` and
`post-rewrite` tee stdin to a temp file so both Totality and the repo hook see every
ref. Totality's own failures never block a commit or push; a repo hook's non-zero
exit still does. See `globalHookSpecs` in `internal/hooks/global.go` for the
installed set and the deliberate exclusions (`push-to-checkout`,
`proc-receive`, `fsmonitor-watchman`, `p4-*`).

Repo hook lookup deliberately avoids `git rev-parse --git-path hooks/<name>`:
that call honors `core.hooksPath` and would resolve straight back to the global
script. The scripts read the repo-scoped `core.hooksPath` with the global and
system config masked, and otherwise fall back to `--git-common-dir`. Repo hooks
carrying the `# tx lifecycle hooks` marker are skipped so a repo that also ran
plain `tx init` does not run Totality twice.

Refusals and opt-outs:

- If `core.hooksPath` is already set globally to a non-Totality directory, install
  refuses instead of overwriting; re-running against the Totality directory is
  idempotent.
- `git config tx.enabled false` inside a repo excludes it (checked by the Go
  handlers via `hooks.EnabledForRepo`); `~/.totality/pause-capture` and
  `TOTALITY_CAPTURE_PAUSED` still pause capture globally.
- Undo with `git config --global --unset core.hooksPath`.
- `tx doctor` reports the global state and flags an overridden or incomplete
  install; it never repairs it silently.

## Manual run

```bash
tx capture push --ref-range main..HEAD --repo .
```

## Telemetry

When `TOTALITY_POSTHOG_KEY` is set, each capture run emits `capture.coverage` and `match.rate` events. Upload telemetry (`session.uploaded`) arrives in WP-1b.

## Agent lifecycle hooks

Claude Code and Codex machine hooks are planned for V1.1. The `hooks.AgentHookConfig` struct reserves the interface; no install yet.

## Git lifecycle recovery

Totality installs `prepare-commit-msg`, `post-commit`, and `post-rewrite` alongside
`pre-push`. Commit-message stamping is deterministic and runs in-process for
`tx commit`, so a missing post hook cannot lose the revision identity.
`post-commit` and `post-rewrite` warn but do not block Git when metadata
recording fails. Before publication, `pre-push` scans revision trailers and
recovers missing rows or stale commit OIDs in the Totality database.
