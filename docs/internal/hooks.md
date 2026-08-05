# Git lifecycle hooks

lgtm installs `prepare-commit-msg`, `post-commit`, `post-rewrite`, and `pre-push`
hooks. Together they stamp revision trailers, record commits, follow rewritten
commit OIDs, and run capture before a push completes. Staged extracts and
sessions are written to `~/.lgtm/lgtm.db` without network access.

Installing the lgtm CLI package provides `lgtm` and `lgtm-mcp`. It does not rewrite every repository immediately; each repo gets the
Git hook when it is initialized with `lgtm init`, initialized automatically by
the CLI or MCP.

## Native Git hooks

The `pre-push` hook invokes:

```bash
lgtm capture push --remote "$remote" --ref-range "$range" --repo "$(git rev-parse --show-toplevel)"
```

The ref range is `remote_sha..local_sha` for updates, or the local SHA for new branches.

The hook skips deleted refs, captures Claude/Codex/Cursor session context by
default, stages capture data locally when upload credentials are missing, and
uploads when credentials are configured.

## Machine-wide install (`lgtm init --global`)

`lgtm init --global` writes the same lifecycle scripts to `~/.lgtm/hooks`
(`$LGTM_HOME/hooks` when set) and points `git config --global core.hooksPath` at
that directory, so lgtm works in every repo without a per-repo `lgtm init`. It does
not have to run inside a repository.

A global `core.hooksPath` makes Git ignore every repository's own
`.git/hooks/*`, so the global scripts shim more than lgtm's four hooks: each
script does lgtm's work (if any) and then executes the repository's own hook of
the same name, forwarding arguments, stdin, and exit status. `pre-push` and
`post-rewrite` tee stdin to a temp file so both lgtm and the repo hook see every
ref. lgtm's own failures never block a commit or push; a repo hook's non-zero
exit still does. See `globalHookSpecs` in `internal/hooks/global.go` for the
installed set and the deliberate exclusions (`push-to-checkout`,
`proc-receive`, `fsmonitor-watchman`, `p4-*`).

Repo hook lookup deliberately avoids `git rev-parse --git-path hooks/<name>`:
that call honors `core.hooksPath` and would resolve straight back to the global
script. The scripts read the repo-scoped `core.hooksPath` with the global and
system config masked, and otherwise fall back to `--git-common-dir`. Repo hooks
carrying the `# lgtm lifecycle hooks` marker are skipped so a repo that also ran
plain `lgtm init` does not run lgtm twice.

Refusals and opt-outs:

- If `core.hooksPath` is already set globally to a non-lgtm directory, install
  refuses instead of overwriting; re-running against the lgtm directory is
  idempotent.
- `git config lgtm.enabled false` inside a repo excludes it (checked by the Go
  handlers via `hooks.EnabledForRepo`); `~/.lgtm/pause-capture` and
  `LGTM_CAPTURE_PAUSED` still pause capture globally.
- Undo with `git config --global --unset core.hooksPath`.
- `lgtm doctor` reports the global state and flags an overridden or incomplete
  install; it never repairs it silently.

## Manual run

```bash
lgtm capture push --ref-range main..HEAD --repo .
```

## Telemetry

When `LGTM_POSTHOG_KEY` is set, each capture run emits `capture.coverage` and `match.rate` events. Upload telemetry (`session.uploaded`) arrives in WP-1b.

## Agent lifecycle hooks

Claude Code and Codex machine hooks are planned for V1.1. The `hooks.AgentHookConfig` struct reserves the interface; no install yet.

## Git lifecycle recovery

lgtm installs `prepare-commit-msg`, `post-commit`, and `post-rewrite` alongside
`pre-push`. Commit-message stamping is deterministic and runs in-process for
`lgtm commit`, so a missing post hook cannot lose the revision identity.
`post-commit` and `post-rewrite` warn but do not block Git when metadata
recording fails. Before publication, `pre-push` scans revision trailers and
recovers missing rows or stale commit OIDs in the lgtm database.
