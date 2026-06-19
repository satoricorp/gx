# Git hooks for gx capture

gx installs a **pre-push** hook that runs the capture pipeline locally before your push completes. Staged extracts and sessions are written to `~/.gx/gx.db` without network access.

## Native git hook

`gx init` writes `.git/hooks/pre-push` that invokes:

```bash
gx capture push --remote "$remote" --ref-range "$range" --repo "$(git rev-parse --show-toplevel)"
```

The ref range is `remote_sha..local_sha` for updates, or the local SHA for new branches.

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
