# Contributing to gx

Thanks for helping. Bug reports, fixes, and docs are all welcome. For anything
larger than a small fix, open a discussion or an issue first so we can agree on
the shape before you write it.

## Setup

You need:

- Go, at the version in `go.mod`
- Git
- [just](https://github.com/casey/just) and zsh, optional, for the build recipes
- [Bun](https://bun.sh), only if you work on the MCP server in `mcp/`

Build and test:

```bash
go build ./...
go vet ./...
go test ./...
```

The tests need no network access, credentials, or gx Cloud account. The few
live tests skip themselves unless their credentials are set.

## Running your build

A source build is a local build: gx runs entirely on your machine, with no
gx Cloud. Install it with:

```bash
just install
```

That installs to `~/.local/bin/gx`, which is where gx's git hooks look for the
binary. If you also use the released gx day to day, reinstall it afterwards
with the command in the README.

To review with a model, set `ANTHROPIC_API_KEY`. To use your own AWS account
instead, set `GX_REVIEW_BEDROCK_DIRECT=1` along with `AWS_ACCESS_KEY_ID`,
`AWS_SECRET_ACCESS_KEY`, and optionally `AWS_SESSION_TOKEN` and `AWS_REGION`;
the account needs access to the Anthropic models on Amazon Bedrock.
`GX_REVIEW_AI=0` runs only the deterministic checks.

The gx Cloud code is dormant, not deleted. `just` reads `.env` from this
directory or any parent, so leave the endpoints in it unset for a local build.
To build against a gx Cloud server of your own, set all of `GITHUB_CLIENT_ID`,
`CONVEX_SITE_URL` and `GX_CLOUD_URL` (see `.env.example`). `just verify-bake`
refuses a binary with only some of them baked in, and AGENTS.md explains why.

## Tests

Two rules come from production bugs that shipped past a green suite. AGENTS.md
has the full story.

- **Seed through the code production uses.** Write test state with the same
  storage functions the product calls. `TestSeedersAreProductionWriters` fails
  on any exported storage method that only tests call.
- **Prefer a weathered database to a pristine one.** Build test state with
  `storagetest.New(t, shapes...)` and repositories with `gxtest.NewWorld(t)`,
  rather than a bare temporary `GX_HOME`. Mark a deliberate cold start with
  `storagetest.NoSessions()`.

## Pull requests

1. Fork the repository and branch from `main`.
2. Keep each PR to one logical change, with tests.
3. Run `gofmt`, `go vet`, and `go test ./...` before pushing. CI runs the same
   checks on Linux and macOS.
4. Describe what changed and why in your own words.

## Sign-off

Sign off every commit to certify the
[Developer Certificate of Origin](https://developercertificate.org/): that you
wrote the change, or otherwise have the right to submit it under this
project's license.

```bash
git commit -s -m "Fix the thing"
```

## License

gx is licensed under the GNU Affero General Public License v3.0 only. By
contributing, you agree that your contribution is licensed under the same terms.
See LICENSE and NOTICE.

## Conduct

Everyone taking part is expected to follow the [code of conduct](CODE_OF_CONDUCT.md).
