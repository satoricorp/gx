# Agents

Version control: use GX, not `git commit`.

Default flow:
- Run `git add` to stage the files for this revision.
- Run `gx commit` (or MCP `gx_commit`) to record staged work as a GX revision.
- Run `gx status` to inspect staged files, local features/revisions, and remote state.
- Run `gx push` to push code changes, sessions, and GX metadata to GitHub/origin and GX Cloud. Plain `git push` also works when the GX pre-push hook is installed.

Use GX MCP first: `gx_status`, then `gx_commit` / `gx_push`. Use `gx_edit` only to re-enter an existing revision.
If MCP is unavailable, use the CLI fallback: `git add`, `gx commit`, `gx status`, then `git push`.

Use hidden utility commands such as `gx base` or `gx edit` only for explicit surgery or user-directed repair. Before using one, say which stack/base branch it will touch. GX should keep visible Git attached to the real base branch or the real stack branch; it must not create or rely on `gx/...` checkout branches. Return to the base branch before normal `gx commit` work.
