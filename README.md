<p align="center">
  <img src="docs/gx-chrome-logo-dark.png" alt="gx" width="420" />
</p>

<p align="center">
  <a href="https://github.com/satoricorp/gx/actions/workflows/ci.yml"><img src="https://github.com/satoricorp/gx/actions/workflows/ci.yml/badge.svg?branch=main" alt="CI"></a>
  <a href="https://github.com/satoricorp/gx/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/satoricorp/gx/ci.yml?branch=main&label=unit%20tests" alt="unit tests"></a>
</p>

SOTA code review that uses your session context.
55% less false positives than CodeRabbit.
23% increased accuracy from Greptile.

Code Review tools are expensive and increasingly unnecessary with the improvement in model performance. The main benefit of GX
is to provide patterns that increase review accuracy. Typically, running a review against a single frontier model will catch about
50% of the bugs of a code review tool.

The magic of these tools lie in the pattern used to extract bugs. In general, additional context, such as knowledge bases/documentation,
don't improve bug identification.

Patterns such as using a judge model and multiple model providers provide excellent results out of the box, alongside session context.
This boosts bug identification by 2x while decreasing false positives.

As models continue to improve, it's reasonable to assume these tools will become less important, but for now this will improve your
code reviews and reduce bugs going into production.


## Get started

Download the archive for your platform from the
[latest GitHub release](https://github.com/satoricorp/gx/releases/latest), extract it,
and add the extracted `gx/bin` directory to your `PATH`.

To build from source instead, clone the repository and run `just install`.

Inside your repo:
```bash
gx init
gx auth login
```

Then use Git as usual. gx hooks record each commit and publish on push:

```bash
gx review

git add -p
git commit -m "some commit message"
git push
```

Those commands record and publish on their own. Pull request comments need a GitHub App (you must create this), and the GX API running on a server.

Create the app under GitHub **Settings → Developer settings → GitHub Apps → New GitHub App**:

- **Webhook URL:** `https://<your-api-host>/github/webhook`. Port 3201.
- **Webhook secret:** you pick this
- **Permissions:** Contents (read and write), Pull requests (read and write), Issues (read and write), Checks (read), Actions (read).
- **Subscribe to events:** Installation, Installation repositories, Pull request, Pull request review, Pull request review comment, Issue comment, and Push.

After GitHub creates the app, copy the App ID, generate a private key, and install the app on the account or org you want to run GX against. 

The API is the `server` package in [satoricorp/console](https://github.com/satoricorp/console). From that repository, with Postgres available, put these in `server/.env`:

```bash
DATABASE_URL=postgres://localhost:5432/gx
GITHUB_APP_ID=<app id>
GITHUB_APP_PRIVATE_KEY=<private key pem>
GITHUB_WEBHOOK_SECRET=<webhook secret>
GITHUB_APP_INSTALL_URL=https://github.com/apps/<your-app-slug>/installations/new
```

`GITHUB_APP_PRIVATE_KEY_PATH` can point at the `.pem` file instead of inlining the key. Then:

```bash
cd server
bun install
bun run migrate
bun run dev
```

Point this CLI at that API with `GX_CLOUD_URL=http://localhost:3201`.

## Commands

| Command | Alias | Purpose |
| --- | --- | --- |
| `gx init` | | Set up this repo (identity, hooks, MCP, agent `/gx` command) |
| `gx auth login` | | Log into gx Cloud |
| `gx auth status` | | Show login status |
| `gx auth logout` | | Log out |
| `gx review [intent]` | `gxr` | Review the current change |
| `gx enhance [intent]` | `gxe` | Top fix from that review, as a prompt for your coding model |
| `gx doctor` | | Diagnose (and optionally repair) local gx state |
| `gx update` | | Install the latest published build |
| `gx version` | | Print version |
| `gx help` | | Help |

Most commands auto-init the repo on first use. Run `gx init` once so identity and `AGENTS.md` are set the way you want.

### `gx init`

```bash
gx init                 # this repo
gx init --global        # machine-wide hooks via core.hooksPath
gx init --yes           # defaults, quiet on success
gx init --name "…" --email "…"
```

### `gx review [intent]`

Reviews the working tree by default. Additional flags based on desired direction of review.

```bash
gx review
gx review --base origin/main #Review a specific branch
gx review --json #For agents
```

| Flag | Effect |
| --- | --- |
| `--base <ref>` | Review `<ref>...HEAD` instead of the working tree |
| `--repo` | Review the whole repository (wins over `--base`) |
| `--fast` | One reviewer, no verification pass |
| `--deep` | More local + indexed context |
| `--focus <path>` | Limit to files under a path prefix |
| `--json` / `--md` | Machine / PR-comment output |
| `--no-comment` | Skip PR comment; still record history |
| `--no-publish` | Skip PR comment and history |
| `--fail-on <level>` | Gate: `none`, `any`, `speculative`, `worth-exploring`, `strong`, `blocking` |
| `--scope <lane>` | `architecture` (default), `security`, `performance`, `onboarding`, `docs`, `dependencies`, `testing`, `maintainability` |
| `--verbose` | Include repo facts, docs, and changed files |
| `--max-findings <n>` | Cap reported recommendations |
| `--client <name>` | Invoking surface (`cli`, `mcp`, `skill`, …) |

Findings: **blocking** (exit 3) vs **advisory**. A finding blocks only when both reviewers raised it and verification confirmed it.

### `gx enhance`
Same review, narrowed to the single highest-value fix, written as a prompt for a coding model. Always exits 0.

```bash
gx enhance | pbcopy
gx enhance "fix the auth timeout" | claude -p "apply this"
gx enhance --base origin/main --fast
```

### `gx doctor` / `gx update` / `gx version`

```bash
gx doctor
gx doctor --fix
gx doctor --json
gx doctor --report          # send diagnosis + recent logs to support

gx update
gx update --check

gx version
gx version --json
```

### `gx auth`

```bash
gx auth login
gx auth login --name "studio-mac"
gx auth status
gx auth status --json
gx auth logout
```

## Uninstall

```bash
rm -f ~/.local/bin/gx ~/.local/bin/gxr ~/.local/bin/gxe ~/.local/bin/gx-mcp
rm -f ~/.local/share/bash-completion/completions/gx ~/.zfunc/_gx
```

Local data stays in `~/.gx` (or `$GX_HOME`).

## Privacy

gx keeps its data on your machine, in `~/.gx` or `$GX_HOME`. It sends no
analytics or telemetry.

gx talks to gx Cloud only after you run `gx auth login`, or when you run
`gx doctor --report` yourself. Once you are signed in:

- **On `git push`**, the pre-push hook publishes the pushed commits' patches,
  gx revision metadata, and the AI coding sessions that produced them. gx
  redacts known secret formats from session events before staging them.
- **`gx review` and `gx enhance`** send the change and its review context to gx
  Cloud, which runs the review models. A review history entry is recorded
  unless you pass `--no-publish`.
- **`gx doctor --report`** sends the diagnosis and recent gx logs to support.

A few things happen without gx Cloud:

- **PR comments.** With a GitHub token available, `gx review` posts its comment
  straight to GitHub. `--no-comment` skips it.
- **Your own AWS account.** `GX_REVIEW_BEDROCK_DIRECT=1` sends review requests to
  Amazon Bedrock in your account instead of gx Cloud.
- **Update checks.** `gx update` and the MCP server read the release manifest
  from `download.gx.run`. The MCP server checks at most once an hour.

## Review knowledge

`gx review --deep` also draws on a hosted index of public engineering
guidance: standards bodies such as OWASP and NIST, language and framework
documentation, and published code-review research. Those documents belong to
their publishers and are not redistributed here.
[docs/review-knowledge-sources.yaml](docs/review-knowledge-sources.yaml) lists
every source and where it lives.

## Build from source

```bash
git clone https://github.com/satoricorp/gx.git
cd gx
go build -o gx ./cmd/gx
```

A source build has no gx Cloud endpoints compiled in, so sign-in, publishing,
and hosted review stay off. [CONTRIBUTING.md](CONTRIBUTING.md) covers local
installs and running reviews on your own AWS account.

## Contributing

Issues and pull requests are welcome. Start with
[CONTRIBUTING.md](CONTRIBUTING.md). Questions go in
[Discussions](https://github.com/satoricorp/gx/discussions). Report security
issues privately as described in [SECURITY.md](SECURITY.md).

## Links

- Product: [gx.run](https://gx.run)
- Releases: [github.com/satoricorp/gx/releases](https://github.com/satoricorp/gx/releases)
- Source: [github.com/satoricorp/gx](https://github.com/satoricorp/gx)

## License

Copyright (C) 2026 Joe LaChance.

gx is free software, licensed under the
[GNU Affero General Public License v3.0 only](LICENSE). See [NOTICE](NOTICE).
