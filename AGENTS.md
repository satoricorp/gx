# Agents

Version control: plain Git. Once `gx init` installs the hooks, GX records and publishes automatically — there is no GX save verb.

Default flow:
- Run `git add` to stage the files for this revision.
- Run `git commit -m "..."` to save. A GX hook records the commit as a reviewable revision.
- Run plain `git push` to publish. The GX pre-push hook captures the session and publishes code changes, sessions, and GX metadata to GX Cloud automatically — do not run `gx push` or `gx capture push` yourself; they bypass the hook.
- Open PRs with `gh pr create` (or the GitHub UI). Do not seed a `## Summary` in the PR body — leave human notes only; GX Cloud appends the rich summary below once the PR exists.
- To amend, use `git commit --amend` and preserve the GX revision trailer in the message.

For AI review, run the `gx_review` MCP tool (or the `gx review` CLI) on the current change.
