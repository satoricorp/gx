# Security policy

## Reporting a vulnerability

Please report security issues privately through GitHub:
[open a private vulnerability report](https://github.com/satoricorp/gx/security/advisories/new).
Do not open a public issue, discussion, or pull request for a vulnerability.

Include what you found, how to reproduce it, and the gx version (`gx version`).
You should get an acknowledgement within a few days. Fixes ship in a new release,
and the advisory is published once users have had a chance to update.

## Supported versions

Only the latest release receives security fixes. `gx update` installs it.

## Scope

In scope:

- The `gx` CLI and its git hooks
- The `gx-mcp` MCP server and the `@satoricorp/gx` npm package
- The install script served from `download.gx.run`
- Session capture, including the redaction that runs before a transcript is
  staged for upload

Session capture is the most sensitive surface. Transcripts can contain
anything an AI agent saw, so gx redacts known secret shapes from session events
before it stages them for upload: AWS access keys, GitHub and Slack tokens,
bearer tokens, API-key and password assignments, private keys, and long
high-entropy strings. A secret that survives redaction and reaches an upload is
a vulnerability. Please report it here.

There is no hosted gx service: gx runs on your machine and talks only to the
model, retrieval, and GitHub endpoints you give it credentials for.
