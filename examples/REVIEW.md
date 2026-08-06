# Review policy

Payments API. Go service, Postgres, deployed to ECS.
Correctness about money beats everything else here.

## Invariants

Things a reviewer cannot infer from the diff, and that we do not want broken:

- Every write to `ledger_entries` is append-only. Nothing updates or deletes a
  row in that table, ever. A migration that adds an `UPDATE` is a bug.
- Amounts are integer minor units. A `float` anywhere near money is wrong.
- Every external call to the payment processor must be idempotent — pass the
  idempotency key from the request, never generate one inside a retry.

## On purpose

Do not flag these. They look like mistakes and are not:

- `internal/legacy/` is frozen. We are deleting it in Q3. Style, naming, and
  test coverage there are not worth comments.
- We ignore errors from the metrics client on purpose. Metrics must never fail
  a request.
- Handlers take concrete structs rather than interfaces. We tried interfaces and
  reverted; do not suggest them again.

## What blocks a merge

- **Blocking**: anything that can lose, duplicate, or misstate a payment.
  Anything that widens auth scope.
- **Worth fixing, not blocking**: missing tests on new branches, unclear naming,
  slow queries under 100ms.
- **Not worth reporting**: formatting, import order, comment style. The linter
  owns those.

## Known and accepted

Already on the roadmap. Reporting these again is noise:

- `ProcessBatch` is O(n²) over the batch. Batches are capped at 500, so it is
  fine until it isn't.
- The webhook handler has no replay protection. Tracked in PAY-412.

## Always check

Repo-specific checks a general reviewer will not think to make:

- A new endpoint needs an entry in `docs/api.md` and an auth test.
- A schema change needs a migration and a test that runs it against a real
  database, not a mock.
- Anything reading `X-Forwarded-For` must go through `realIP()`.

## high-risk paths

`risk-path: <glob> — <why it is risky>`. gx raises the severity of findings in
these files and quotes your reason in the report.

risk-path: internal/ledger/** — Append-only money records. A bug here is silent and permanent.
risk-path: internal/auth/** — Token issuance and scope checks. Review ownership and expiry.
risk-path: migrations/** — Runs against production data once, unreviewable afterwards.
risk-path: **/Dockerfile — Base image and runtime user affect every deployment.
