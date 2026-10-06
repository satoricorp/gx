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

A plain source build has no gx Cloud endpoints compiled in. It runs, but
sign-in, publishing, and hosted review are off. To install a local-only build
on purpose:

```bash
GX_ALLOW_UNBAKED=1 just install
```

That installs to `~/.local/bin/gx`, which is where gx's git hooks look for the
binary. If you also use the released gx day to day, reinstall it afterwards
with the command in the README.

To run a review from a source build, use your own AWS account. Set
`GX_REVIEW_BEDROCK_DIRECT=1` along with `AWS_ACCESS_KEY_ID`,
`AWS_SECRET_ACCESS_KEY`, and optionally `AWS_SESSION_TOKEN` and `AWS_REGION`.
The account needs access to the Anthropic models on Amazon Bedrock.

Maintainers who build against gx Cloud copy `.env.example` to `.env`, fill it
in, and run `just install`. `just verify-bake` refuses a binary whose cloud
endpoints are missing. AGENTS.md explains why that check exists.

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
