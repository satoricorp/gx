# GX Manual Workflow Test Sheet

> **Note (2026):** Replace `gx compose` / `gx stacks` / `gx add` in this sheet
> with `git add` + `gx commit`, `gx status`, and `gx generate` where bulk split
> is needed.

Use this sheet against a disposable GitHub repository. The goal is to test the
real flow, not mocked services.

## Run Metadata

- Date:
- Tester:
- GX version or commit:
- Desktop app version or commit:
- Fixture repo:
- Notes / bug links:

## Fixture Setup

```bash
export GX_HOME="$(mktemp -d)"
git clone git@github.com:<owner>/<fixture-repo>.git
cd <fixture-repo>
gx auth status || gx auth login
gx init --name "GX Flow Tester" --email "gx-flow@example.com"
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
gx status
gx compose
gx stacks
gx publish
```

Then:

- Open the PR in GitHub.
- Open GX Desktop.
- Find the published stack.
- Verify the PR link/status appears.
- Merge from GX Desktop.
- Refresh local state.

```bash
git fetch origin
git log --oneline origin/main -5
gx sync
gx stacks
```

### Pass Criteria

- `gx publish` prints a branch and GitHub PR URL.
- GitHub PR contains the dummy changes.
- GX Desktop shows the published stack.
- Merge from GX Desktop succeeds.
- `origin/main` contains the dummy changes.
- `gx stacks` agrees with the merged/published state.

Result: PASS / FAIL

Notes:

## Test 2: Republish Existing PR

Purpose: prove editing a published GX revision updates the same GitHub PR.

### Steps

Start from an open published PR.

```bash
gx stacks
gx edit <revision-or-change-id>
printf "republish %s\n" "$(date +%s)" >> counter.txt
gx status
gx publish
```

Then verify in GitHub and GX Desktop.

### Pass Criteria

- Same PR number is reused.
- PR branch/head SHA changes.
- GitHub shows the new dummy change.
- GX Desktop shows the updated stack/review state.
- Merge still succeeds from GX Desktop.

Result: PASS / FAIL

Notes:

## Test 3: Conflict Repair by Editing Existing Revision

Purpose: prove a conflicted published PR can be repaired without creating a
disconnected replacement PR.

### Steps

Publish a PR that changes `shared.txt`:

```bash
printf "gx change %s\n" "$(date +%s)" > shared.txt
gx compose
gx publish
```

In a separate clone, advance `main` with a conflicting edit:

```bash
git clone git@github.com:<owner>/<fixture-repo>.git /tmp/gx-flow-main
cd /tmp/gx-flow-main
git switch main
printf "main change %s\n" "$(date +%s)" > shared.txt
git add shared.txt
git commit -m "advance main with conflicting shared text"
git push origin main
```

Back in the GX repo:

```bash
git fetch origin
gx sync
gx stacks
```

Then:

- Confirm GitHub or GX Desktop reports the PR as conflicted/not mergeable.
- Use GX/AI conflict repair.
- Prefer editing the existing revision.
- Run `gx publish` again.
- Merge from GX Desktop.

### Pass Criteria

- Conflict is visible before repair.
- Repair updates the same PR number.
- PR becomes mergeable.
- Merge from GX Desktop succeeds.
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
gx compose
gx stacks
gx publish
```

Then verify GitHub and GX Desktop.

### Pass Criteria

- `gx compose` creates the expected revision structure.
- `gx stacks` shows all revisions in the intended order.
- GitHub PR contains all changes.
- GX Desktop displays the stack/revisions coherently.
- Merge from GX Desktop updates `main`.

Result: PASS / FAIL

Notes:

## Test 5: Publish Returns To Main

Purpose: prove publish pushes accepted stacks without requiring a current stack
checkout and returns to `main`.

### Steps

```bash
printf "review only %s\n" "$(date +%s)" >> counter.txt
gx compose
gx publish
git branch --show-current
```

### Pass Criteria

- Every accepted unpublished stack is pushed.
- GX review/publish metadata is recorded.
- The final Git branch is `main`.

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
