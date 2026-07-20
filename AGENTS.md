# Agents

Version control: use Git with GX hooks and metadata, not raw `git commit` for normal agent work.

Default flow:
- Run `git add` to stage the files for this revision.
- Run `gx commit` (or MCP `gx_commit`) to record staged work as a GX revision via native `git commit` plus GX hooks.
- Run `gx status` to inspect staged files, local revisions, and remote state.
- Run plain `git push` to ship. The GX pre-push hook captures sessions and publishes GX metadata to GX Cloud automatically — do not run `gx push` or `gx capture push` yourself.
- Open PRs with `gh pr create` (or the GitHub UI). Do not seed a `## Summary` in the PR body — leave human notes only; GX Cloud appends the rich summary below once the PR exists.

Use GX MCP first: `gx_status`, then `gx_commit`; push with plain `git push`.

To amend an existing GX revision, use `git commit --amend` and preserve the GX revision trailer in the message.
