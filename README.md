# gx

Code review with session context, on top of plain Git.

```bash
curl -fsSL https://download.gx.run/install.sh | sh
```

## Get started

```bash
cd your-repo
gx init
gx auth login
```

Then use Git as usual. gx hooks record each commit and publish on push:

```bash
git add -p
git commit -m "fix the auth timeout"
git push
gx review "fix the auth timeout"
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

Reviews the working tree by default. Optional `intent` is one sentence of what the change was supposed to do.

```bash
gx review
gx review "make checkout survive gateway blips"
gx review --base origin/main
gx review --repo
gx review --fast
gx review --deep
gx review --focus internal/auth
gx review --json
gx review --md
gx review --fail-on blocking   # default; exit 3 if blocking findings remain
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

Model presets: `GX_REVIEW_MODELS=default|budget|glm`. Per-slot overrides: `GX_REVIEW_BEDROCK_MODEL_A`, `GX_REVIEW_BEDROCK_MODEL_B`, `GX_REVIEW_JUDGE_MODEL`, `GX_GATE_MODEL`.

Own-AWS review (skip gx Cloud Bedrock brokerage):

```bash
export GX_REVIEW_BEDROCK_DIRECT=1
export AWS_ACCESS_KEY_ID=… AWS_SECRET_ACCESS_KEY=…   # + AWS_SESSION_TOKEN if needed
```

### `gx enhance [intent]`

Same review, narrowed to the single highest-value fix, written as a prompt for a coding model. Always exits 0.

```bash
gx enhance | pbcopy
gx enhance "fix the auth timeout" | claude -p "apply this"
gx enhance --base origin/main --fast
```

Flags: `--base`, `--repo`, `--fast`, `--deep`, `--focus`, `--client`.

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

Cloud login is required for review unless you set `GX_REVIEW_BEDROCK_DIRECT=1`.

## Daily flow

1. `git add` / `git commit` — hooks record the revision.
2. `git push` — pre-push publishes session + metadata to gx Cloud. Do not run `gx push` / `gx capture push`; they bypass the hook.
3. Open the PR with `gh pr create` (or the UI). Leave human notes only; gx appends the summary after a hooked push.
4. `gx review` before you ship. Blocking finding ⇒ no-ship.

```bash
gx init --global   # optional: hooks in every repo on this machine
```

## Uninstall

```bash
rm -f ~/.local/bin/gx ~/.local/bin/gxr ~/.local/bin/gxe ~/.local/bin/gx-mcp
rm -f ~/.local/share/bash-completion/completions/gx ~/.zfunc/_gx
```

Local data stays in `~/.gx` (or `$GX_HOME`).

## Links

- Product: [gx.run](https://gx.run)
- Installer: [download.gx.run/install.sh](https://download.gx.run/install.sh)
