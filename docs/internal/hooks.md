# Git lifecycle hooks

GX installs `prepare-commit-msg`, `post-commit`, `post-rewrite`, and `pre-push`
hooks. Together they stamp revision trailers, record commits, follow rewritten
commit OIDs, and run capture before a push completes. Staged extracts and
sessions are written to `~/.gx/gx.db` without network access.

Installing the GX menu-bar app installs the bundled `gx` CLI and `gx-mcp`
binary. It does not rewrite every repository immediately; each repo gets the
Git hook when it is initialized with `gx init`, initialized automatically by
the CLI or MCP.

## Native Git hooks

The `pre-push` hook invokes:

```bash
gx capture push --remote "$remote" --ref-range "$range" --repo "$(git rev-parse --show-toplevel)"
```

The ref range is `remote_sha..local_sha` for updates, or the local SHA for new branches.

The hook skips deleted refs, captures Claude/Codex/Cursor session context by
default, stages capture data locally when upload credentials are missing, and
uploads when credentials are configured.

## Team propagation

### lefthook

Generate a starter template:

```go
hooks.WriteLefthookTemplate(repoRoot, gxBinaryPath)
```

Or add manually to `.lefthook.yml`:

```yaml
pre-push:
  commands:
    gx-capture:
      run: gx capture push --remote "$1" --ref-range "$range"
```

See `internal/hooks/lefthook.go` for the full template.

### husky

Add a `.husky/pre-push` script using `hooks.HuskyPrePushSnippet(gxPath)`.

gx does not install husky or lefthook automatically — copy the snippet into your team's hook manager.

## Manual run

```bash
gx capture push --ref-range main..HEAD --repo .
```

## Telemetry

When `GX_POSTHOG_KEY` is set, each capture run emits `capture.coverage` and `match.rate` events. Upload telemetry (`session.uploaded`) arrives in WP-1b.

## Agent lifecycle hooks

Claude Code and Codex machine hooks are planned for V1.1. The `hooks.AgentHookConfig` struct reserves the interface; no install yet.

## Git lifecycle recovery

GX installs `prepare-commit-msg`, `post-commit`, and `post-rewrite` alongside
`pre-push`. Commit-message stamping is deterministic and runs in-process for
`gx commit`, so a missing post hook cannot lose the revision identity.
`post-commit` and `post-rewrite` warn but do not block Git when metadata
recording fails. Before publication, `pre-push` scans revision trailers and
recovers missing rows or stale commit OIDs in the GX database.
