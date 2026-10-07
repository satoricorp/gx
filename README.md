<p align="center">
  <img src="docs/gx-chrome-logo-dark.png" alt="gx" width="420" />
</p>

SOTA code review that uses your session context.


With traditional code review, developers held the context for code changes, but now AI sessions hold that context.
This can improve code review drastically, keeping both session history and informing future sessions of best practices.


## Get started

Install GX locally:
```bash
curl -fsSL https://download.gx.run/install.sh | sh
```

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

Sign up at [gx.run](https://gx.run). Re-run the install command to upgrade.

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
- Installer: [download.gx.run/install.sh](https://download.gx.run/install.sh)

## License

Copyright (C) 2026 Joe LaChance.

gx is free software, licensed under the
[GNU Affero General Public License v3.0 only](LICENSE). See [NOTICE](NOTICE).
