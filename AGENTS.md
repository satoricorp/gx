# Agents

Version control: use GX, not `git commit`.

Default flow:
- Run `git add` to stage the files for this revision.
- Run `gx commit` (or MCP `gx_commit`) to record staged work as a GX revision.
- Run `gx status` to inspect staged files, local features/revisions, and remote state.
- Run plain `git push` to ship. The GX pre-push hook publishes code changes, sessions, and GX metadata to GX Cloud automatically — do not run `gx push` or `gx capture push` yourself.
- Open PRs with `gh pr create` (or the GitHub UI); the GX summary arrives via the GitHub App once the PR exists.

Use GX MCP first: `gx_status`, then `gx_commit`; push with plain `git push`. Use `gx_edit` only to re-enter an existing revision.

Use hidden utility commands such as `gx base` or `gx edit` only for explicit surgery or user-directed repair. Before using one, say which stack/base branch it will touch. GX should keep visible Git attached to the real base branch or the real stack branch; it must not create or rely on `gx/...` checkout branches. Return to the base branch before normal `gx commit` work.
