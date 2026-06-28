# Agents

Version control: use GX, not `git commit`.

Default flow:
- Run `gx generate` to create local features/stacks and revisions from current work.
- Run `gx status` to inspect unstaged files, local features/revisions, and Remote state.
- Run `gx push` to push code changes, sessions, and GX metadata to GitHub/origin and GX Cloud.

Use hidden utility commands such as `gx add`, `gx switch`, `gx base`, `gx edit`, or `gx modify` only for explicit surgery or user-directed repair. Before using one, say which stack/base branch it will touch. GX should keep visible Git attached to the real base branch or the real stack branch; it must not create or rely on `gx/...` checkout branches. Return to the base branch before normal `gx generate` work.
