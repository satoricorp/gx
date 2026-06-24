# Agents

Version control: use GX, not `git commit`.

Default flow:
- Run `gx compose` to propose or update the single pending compose proposal.
- Review only if compose reports warnings, repair hints, or proposal issues.
- Accept ready revisions from compose so they move into `gx stacks`.
- Run `gx publish` to publish accepted stacks.

Use hidden utility commands such as `gx add`, `gx switch`, `gx base`, `gx edit`, or `gx modify` only for explicit surgery or user-directed repair. Before using one, say which stack/base branch it will touch. GX should keep visible Git attached to the real base branch or the real stack branch; it must not create or rely on `gx/...` checkout branches. Return to the base branch before normal `gx compose` work.
