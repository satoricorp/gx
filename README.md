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
```

Give `gx review` a model to run, using your own Anthropic API key:
```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

Then use Git as usual. gx hooks record each commit, and the pre-push hook
captures the coding session behind it:

```bash
gx review

git add -p
git commit -m "some commit message"
git push
```

Run `gx update`, or re-run the install command, to upgrade. If `gx review` asks
you to run `gx auth login`, you have an older build that expects a gx Cloud:
upgrade, and it reviews with your key instead.

## Running locally

gx runs entirely on your machine. There is no account and no server: what gx
records stays in `~/.gx` (or `$GX_HOME`), and the only calls it makes are the
ones you give it credentials for.

**Models.** `gx review` and `gx enhance` need one of:

| Set | Reviews with |
| --- | --- |
| `ANTHROPIC_API_KEY` | The Anthropic API, on your account |
| `GX_REVIEW_BEDROCK_DIRECT=1` plus `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` (and optionally `AWS_SESSION_TOKEN`, `AWS_REGION`) | Amazon Bedrock, on your AWS account |
| `GX_REVIEW_AI=0` | No model: only the deterministic checks and your project's own linters and tests |

The default panel is Claude Haiku 4.5 and Claude Sonnet 4.6, verified by Claude
Haiku 4.5. `GX_REVIEW_BEDROCK_MODEL_A`, `GX_REVIEW_BEDROCK_MODEL_B` and
`GX_REVIEW_JUDGE_MODEL` take a Claude model ID such as `claude-opus-5-5`.
A leg named `openai:<model>` runs on the OpenAI API with `OPENAI_API_KEY`, and
`GX_REVIEW_BEDROCK_MODEL_B=off` reviews with one model.

**Retrieval.** With `TURBOPUFFER_API_KEY` and `OPENAI_API_KEY` (for
embeddings), `gx index` indexes the repository and reviews retrieve from it.
Without them, a review reads the change and the files around it.

**GitHub.** With `GH_TOKEN` or `GITHUB_TOKEN` set, `gx review` comments on the
branch's pull request. `--no-comment` skips that.

## Commands

| Command | Alias | Purpose |
| --- | --- | --- |
| `gx init` | | Set up this repo (identity, hooks, MCP, agent `/gx` command) |
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
| `--no-comment` | Skip the PR comment |
| `--no-publish` | Skip the PR comment and any review history |
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

gx update
gx update --check

gx version
gx version --json
```

## Uninstall

```bash
rm -f ~/.local/bin/gx ~/.local/bin/gxr ~/.local/bin/gxe ~/.local/bin/gx-mcp
rm -f ~/.local/share/bash-completion/completions/gx ~/.zfunc/_gx
```

Local data stays in `~/.gx` (or `$GX_HOME`).

## Privacy

gx keeps its data on your machine, in `~/.gx` or `$GX_HOME`. It sends no
analytics or telemetry, and `git push` uploads nothing beyond the push itself:
the pre-push hook only captures the session into `~/.gx`, redacting known
secret formats from session events as it does.

What leaves your machine, and only with credentials you set:

- **Reviews.** `gx review` and `gx enhance` send the change and its review
  context to the model provider you configured: the Anthropic API, Amazon
  Bedrock in your account, or OpenAI for an `openai:` leg.
- **Retrieval.** With TurboPuffer and OpenAI keys set, indexing and review send
  code chunks to OpenAI for embeddings and to TurboPuffer for storage and search.
- **PR comments.** With a GitHub token available, `gx review` posts its comment
  to GitHub. `--no-comment` skips it.
- **Update checks.** `gx update` and the MCP server read the release manifest
  from `download.gx.run`. The MCP server checks at most once an hour.

The gx Cloud client code is still in the repository but dormant: it does
nothing unless a build is configured with a gx Cloud server of its own (see
[CONTRIBUTING.md](CONTRIBUTING.md)).

## Review knowledge

With a gx Cloud configured, `gx review --deep` also draws on a hosted index of
public engineering guidance: standards bodies such as OWASP and NIST, language
and framework documentation, and published code-review research. A local build
has no access to that index and reviews without it. Those documents belong to
their publishers and are not redistributed here.
[docs/review-knowledge-sources.yaml](docs/review-knowledge-sources.yaml) lists
every source and where it lives.

## Build from source

```bash
git clone https://github.com/satoricorp/gx.git
cd gx
go build -o gx ./cmd/gx
```

A source build is a local build, the same as a release.
[CONTRIBUTING.md](CONTRIBUTING.md) covers installing it where gx's git hooks
look for it.

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
