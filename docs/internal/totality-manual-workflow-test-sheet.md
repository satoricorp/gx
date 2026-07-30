# Totality Manual Workflow Test Sheet

> **Note (2026):** Replace `tx compose` / `tx stacks` / `tx add` / `tx publish`
> in this sheet with `git add` + `tx commit`, `tx status`, `tx generate` where
> bulk split is needed, and plain `git push` + `gh pr create` to publish.

Use this sheet against a disposable GitHub repository. The goal is to test the
real flow, not mocked services.

## Run Metadata

- Date:
- Tester:
- Totality version or commit:
- Desktop app version or commit:
- Fixture repo:
- Notes / bug links:

## Fixture Setup

```bash
export TOTALITY_HOME="$(mktemp -d)"
git clone git@github.com:<owner>/<fixture-repo>.git
cd <fixture-repo>
tx auth status || tx auth login
tx init --name "Totality Flow Tester" --email "totality-flow@example.com"
git fetch origin
git switch main
git pull origin main
```

Recommended fixture files:

- `counter.txt`
- `shared.txt`
- `src/message.txt`

## Test 1: Happy Path Publish and Merge

Purpose: prove code change to compose to publish to desktop merge to GitHub
`main`.

### Steps

```bash
printf "counter %s\n" "$(date +%s)" >> counter.txt
printf "message %s\n" "$(date +%s)" >> src/message.txt
tx status
git add counter.txt src/message.txt
tx commit -m "dummy counter and message update"
tx status
git push
gh pr create
```

Then:

- Open the PR in GitHub.
- Open Totality Desktop.
- Find the published stack.
- Verify the PR link/status appears.
- Merge from Totality Desktop.
- Refresh local state.

```bash
git fetch origin
git log --oneline origin/main -5
tx doctor
git status
```

### Pass Criteria

- `git push` succeeds and the pre-push hook runs; `gh pr create` prints a GitHub PR URL.
- GitHub PR contains the dummy changes.
- Totality Desktop shows the published stack.
- Merge from Totality Desktop succeeds.
- `origin/main` contains the dummy changes.
- `tx status` agrees with the merged/published state.

Result: PASS / FAIL

Notes:

## Test 2: Republish Existing PR

Purpose: prove editing a published Totality revision updates the same GitHub PR.

### Steps

Start from an open published PR.

```bash
tx status
tx edit <revision-or-change-id>
printf "republish %s\n" "$(date +%s)" >> counter.txt
git add counter.txt
tx commit -m "republish counter update"
tx status
git push
```

Then verify in GitHub and Totality Desktop.

### Pass Criteria

- Same PR number is reused.
- PR branch/head SHA changes.
- GitHub shows the new dummy change.
- Totality Desktop shows the updated stack/review state.
- Merge still succeeds from Totality Desktop.

Result: PASS / FAIL

Notes:

## Test 3: Conflict Repair by Editing Existing Revision

Purpose: prove a conflicted published PR can be repaired without creating a
disconnected replacement PR.

### Steps

Publish a PR that changes `shared.txt`:

```bash
printf "tx change %s\n" "$(date +%s)" > shared.txt
git add shared.txt
tx commit -m "change shared text"
git push
gh pr create
```

In a separate clone, advance `main` with a conflicting edit:

```bash
git clone git@github.com:<owner>/<fixture-repo>.git /tmp/totality-flow-main
cd /tmp/totality-flow-main
git switch main
printf "main change %s\n" "$(date +%s)" > shared.txt
git add shared.txt
git commit -m "advance main with conflicting shared text"
git push origin main
```

Back in the Totality repo:

```bash
git fetch origin
tx doctor
git status
```

Then:

- Confirm GitHub or Totality Desktop reports the PR as conflicted/not mergeable.
- Use Totality/AI conflict repair.
- Prefer editing the existing revision.
- Run `git push` again.
- Merge from Totality Desktop.

### Pass Criteria

- Conflict is visible before repair.
- Repair updates the same PR number.
- PR becomes mergeable.
- Merge from Totality Desktop succeeds.
- `origin/main` contains the repaired final content.

Result: PASS / FAIL

Notes:

## Test 4: Multi-Revision Stack

Purpose: prove compose/publish/desktop handle multiple revisions in one stack.

### Steps

```bash
printf "alpha %s\n" "$(date +%s)" >> counter.txt
printf "beta %s\n" "$(date +%s)" >> src/message.txt
printf "gamma %s\n" "$(date +%s)" >> shared.txt
git add counter.txt src/message.txt shared.txt
tx generate
tx status
git push
gh pr create
```

Then verify GitHub and Totality Desktop.

### Pass Criteria

- `tx generate` / `tx commit` creates the expected revision structure.
- `tx status` shows all revisions in the intended order.
- GitHub PR contains all changes.
- Totality Desktop displays the stack/revisions coherently.
- Merge from Totality Desktop updates `main`.

Result: PASS / FAIL

Notes:

## Test 5: Publish Returns To Main

Purpose: prove `git push` publishes accepted stacks through the pre-push hook
and leaves a normal branch checkout.

### Steps

```bash
printf "review only %s\n" "$(date +%s)" >> counter.txt
git add counter.txt
tx commit -m "review only counter update"
git push
git branch --show-current
```

### Pass Criteria

- Every accepted unpublished stack is pushed via plain `git push`.
- Totality review/publish metadata is recorded by the pre-push hook.
- The final Git branch is a normal attached checkout (typically `main` or the stack branch per product rules).

Result: PASS / FAIL

Notes:

## Final Sign-Off

- Happy path passed:
- Republish passed:
- Conflict repair passed:
- Multi-revision stack passed:
- publish returns to `main` passed:
- Bugs filed:
- Ship confidence: LOW / MEDIUM / HIGH
