# gx:recommended v3

The built-in rule pack. Written in the REVIEW.md grammar so the same parser
reads it: every `##` heading below is one rule, its slug is the rule's
permanent name, and an `Advisory:` prefix marks a rule that informs rather
than blocks. A repository's own REVIEW.md may add rules or list exceptions;
it cannot change these.

## No secrets

Credentials, API keys, tokens, private keys, and connection strings with embedded passwords must not appear as literals in source, config committed to the repo, test fixtures that reach production paths, log lines, error messages, panics, or telemetry. This catches a hard-coded key used to "make it work for now", a `.env` value copied into a Go or TS constant, a real token pasted into a test, `log.Printf("auth failed: %v", req)` where the request holds a bearer token, and an error that wraps the raw Authorization header. A clearly fake placeholder (`sk-test-000`, `changeme`, `example.com`) that no live system accepts, a public identifier such as a client ID, or logging that a secret was present (its length, a fixed-width prefix for correlation) without the value itself does not count. Move the value behind the environment, a secrets manager, or the existing config loader; redact it before logging; and if a real secret was committed say so plainly: it needs rotating, not just deleting.

## No injection sinks

Data that came from a user, a request, a file, or another service must not be spliced into a query, a shell command, a URL, or an HTML document without the escaping that sink expects. This catches string-concatenated SQL, `exec.Command("sh", "-c", userInput)`, a template rendered with `template.HTML(input)`, a path built by joining a request parameter without a root check, and a format call whose *template* is attacker-controlled (`fmt.Sprintf(userFormat, x)`). Interpolating data into a template the author wrote is ordinary string building and is not a sink: `f"Quick update on {product}"` is prose, and reporting it as injection costs the reader more than the rule saves. A constant string, a value already validated against an allowlist, a parameter passed through a driver's placeholder mechanism, and an identifier — a table, column, or schema name, which placeholders cannot carry — routed through the library's identifier-quoting helper (`quote_ident`, `psycopg.sql.Identifier`) do not count; quoting the identifier is the correct fix, not a workaround for a missing one. Name the expression that reaches the sink and the call that executes it; a finding that cannot point at both is not this rule. Use the parameterized or escaping form the library provides, and validate at the boundary before the value reaches the sink.

## No swallowed errors

An error returned by a call must be handled, returned, or deliberately discarded with a comment saying why it is safe. This catches `_ = f()` on a call whose failure changes state, `if err != nil { return nil }` that turns a failure into an empty success, an `except: pass`, and a `.catch(() => {})` around a write. Ignoring the error of a best-effort close, a log call, or a cleanup already covered by a later check does not count when the code says so. Return the error, wrap it with context, or handle the specific case and let the rest propagate.

## Tests can fail

A test must be able to fail when the behavior it names is broken. This catches a test with no assertion, an assertion on a value the test itself just set, a mock that returns whatever the code under test asserts, an expected-error test that passes on any error, and a `t.Skip` or `.skip` left in place. A test that only checks the code does not panic counts as passing this rule only when that is what the test's name claims. Assert on the observable result of the behavior, use a real or recorded input where a mock would echo the answer, and delete or fix a test that cannot fail.

## No placeholder code

Code committed as done must not contain a stub standing in for the real thing. This catches `TODO: implement`, a function body that is only `return nil` or `raise NotImplementedError` behind a real call site, a hard-coded value where a computation was asked for, and a `// temporary` shim that ships. A stub that is documented as intentional, gated behind a feature flag, and covered by an issue reference does not count. Implement it, or make the gap explicit at the call site so nothing depends on it silently.

## Advisory: Comments match code

Comments, doc strings, and commit-adjacent prose must describe what the code now does. This catches a doc comment that names a parameter the function no longer takes, a comment explaining a branch that was removed, an example that no longer compiles, and a "returns nil on missing" note above code that now returns an error. Comments that state intent or rationale rather than mechanics do not count when the intent still holds. Update the comment in the same change that moves the code, or delete a comment that only repeated the code.

## Advisory: No unbounded work

Loops, queries, reads, and buffers driven by external input must have a bound the caller can name. This catches reading a whole request body into memory with no limit, a query without a `LIMIT` on a table that grows, a recursive walk without depth or count, and a retry loop with no maximum. Work over a set that is fixed at compile time or bounded by an earlier check does not count. Add the limit at the point the input enters, and reject or page past it rather than growing.

## Advisory: Cleanup on failure

Resources acquired on a code path must be released on every exit from it, including the failure exits. This catches a file opened before an early return, a lock taken without `defer` or `finally`, a temp directory left behind when a later step errors, and a transaction begun but neither committed nor rolled back on the error branch. A resource that is intentionally handed off to a caller who now owns it does not count when the handoff is clear. Use `defer`, `with`, `try/finally`, or the language's equivalent immediately after acquiring, and roll back before returning an error.

## Handles missing values

Code that reads a value which can be absent — a nil pointer, an empty slice, a missing map key, an optional field, an empty string from the environment — must handle the absent case on the path it actually takes. This catches dereferencing a lookup result without checking `ok`, indexing `parts[1]` after a split that may return one element, calling a method on a nullable return, and treating an unset environment variable as configured. A value guaranteed present by construction, a check earlier on the same path, or a type that cannot be absent does not count. Check for the absent case where the value is read, and give it a defined behavior: a default, an error, or a skip.

## Advisory: Dont repeat yourself

A change must use the helper, type, or pattern the codebase already has for a job instead of writing a new one, and must not introduce a block of logic that already appears elsewhere in the same change or the codebase with only names swapped. This catches a second HTTP client wrapper, a re-implemented slugify or path-join, a bespoke retry loop next to the shared one, a new config reader when the loader is three packages over, the same validation copied into three handlers, and a switch that mirrors one two files away. A local function that differs in a way the change explains, a copy made to break an unwanted dependency, or duplication that is cheaper than the coupling a helper would create, does not count. Call or extend the existing helper, or extract the shared block once and call it from each site.

## No invented packages

Every import, dependency, and module reference must name something that exists and is actually installed or declared. This catches an import of a package that does not exist on the registry, a module path that is a plausible-sounding blend of two real ones, a version that was never published, and a sub-package that the real library does not expose. A local package the same change adds does not count when it is present in the diff. Remove the import, or replace it with the real package after confirming its name and version.

## No invented APIs

Every function, method, field, flag, and option a change calls must exist on the type or module it is called on, with the signature the call uses. This catches a method that a library never had, a parameter that sounds right but is not accepted, a struct field renamed in a version the repo does not use, and a CLI flag that the tool does not define. A symbol the same change adds does not count when it is present in the diff. Check the actual declaration and call what is there.

## Advisory: APIs used as documented

A call must respect the documented contract of what it calls: ordering, ownership, thread-safety, error semantics, and what the return value means. This catches reusing a request body after it was consumed, calling a method that is documented as not concurrency-safe from multiple goroutines, treating a "not found" sentinel as a failure, and ignoring a documented requirement to call `Close`. Behavior the documentation leaves open, or a well-established convention the codebase already relies on, does not count. Follow the documented contract, and when the contract is unclear, say so in a comment at the call site.

## Endpoints enforce authz

Every new or changed handler, route, RPC, or job entry point that reads or writes data must check that the caller is allowed to, on the path that reaches the data. This catches a new endpoint added without the auth middleware the neighbors use, a handler that authenticates but never checks tenant or ownership, an admin route reachable by any signed-in user, and an internal endpoint exposed with no check because "only we call it". A path that inherits the check from a router group or middleware the diff visibly attaches, or a public endpoint whose data is public by design, does not count. Attach the existing authorization check, and verify the object-level check (this caller, this record) not just the session check.

## Advisory: No dead code

A change must not leave behind code nothing reaches: functions with no callers, branches that cannot be taken, flags nothing reads, and commented-out blocks. This catches the old implementation kept "just in case" next to the new one, a helper whose only caller was deleted in this change, and an unreachable `else` after a return. Exported API kept for compatibility, or code reached only by reflection or generated callers, does not count when that is stated. Delete it; version control keeps the history.

## Advisory: No speculative layers

A change must not add interfaces, generics, plugin points, or configuration for needs nothing in the codebase has yet. This catches an interface with one implementation and no test double, a factory that returns one type, an options struct with a single field, and a layer that forwards every call unchanged. An abstraction that a second concrete use in the same change or in the codebase already needs, or a seam the tests actually exercise, does not count. Write the concrete version now and abstract when the second use arrives.

## Advisory: Respects module seams

A change must import and depend in the direction the codebase's layering already establishes. This catches a low-level package importing from the CLI or HTTP layer, a domain package reaching into a storage implementation, an internal package used from outside its module, and a cycle broken only by an init hack. A dependency the module already had, or a boundary the change explicitly and deliberately moves with the rest of the code updated to match, does not count. Move the code to the layer that owns it, or pass the value in instead of importing the layer.

## No race hazards

Shared state read or written from more than one goroutine, thread, request, or process must be protected by a lock, an atomic, a channel, or a design that guarantees single ownership. This catches a map written from a handler with no mutex, a check-then-act on a shared counter, a cached value initialized lazily without synchronization, and a global mutated by concurrent tests. State confined to one goroutine, immutable after construction, or already guarded by a lock the code visibly holds does not count. Name the interleaving that goes wrong, then guard the state or restructure so only one owner touches it.

## Retries are idempotent

Any operation that is retried — by a loop, a queue, a client, or an at-least-once delivery — must be safe to run more than once. This catches a retry around a payment or email send with no idempotency key, a queue consumer that inserts without a dedupe check, a "resend on timeout" where the first attempt may have succeeded, and a migration step that is not re-runnable. A retry around a pure read, or a write keyed so the second attempt is a no-op, does not count. Add an idempotency key or a conditional write, or make the retry decision after checking whether the first attempt landed.

## No vulnerable deps

A dependency added or upgraded by the change must not be a version with a known, applicable vulnerability, and must not be a package that is unmaintained, typosquatted, or unusually new for its role. This catches pinning a version with a published advisory that affects the code path used, adding a package whose name is one letter off a popular one, and pulling in a large transitive tree for one function. An advisory that applies only to a feature the repo does not use, when that is stated, does not count. Move to the fixed version, or pick the maintained package the ecosystem already uses.

## Scope matches intent

Everything in the change must serve what was asked, and everything that was asked must be in the change. This catches a rename or reformat riding along with a fix, a second feature bundled into the diff, regenerated files that dwarf the real change, and a requested part of the task quietly left out. A small adjacent fix that the change explains, or a mechanical update the intent required, does not count. Split the unrelated work into its own change, and finish or explicitly defer the missing part.

## Claims match diff

The description, commit message, PR body, and comments must describe what the diff actually does. This catches "adds tests" when no test file changed, "no behavior change" over a diff that alters a return value, "fixes the race" over a diff that only adds a log line, and a summary that names files the diff does not touch. Wording that is imprecise but not wrong does not count. Rewrite the claim to match the code, or change the code to match the claim.

## No weakened checks

A change must not loosen a check, guard, test, lint rule, or type constraint in order to make something pass. This catches a deleted or skipped test, an assertion made looser, a lint rule disabled inline, an `any` replacing a specific type, a timeout raised without a cause, and a validation branch commented out. A check that was wrong and is fixed with an explanation, or one deliberately retired with the reason in the change, does not count. Fix the code the check was catching, and if the check itself was wrong, say so where it is changed.

## No silent regressions

A change must not remove or alter behavior a caller depends on without saying so. This catches a function that stops handling a case it used to, a default value changed in passing, an error path that now returns success, a field dropped from an output, and a sort order that changed as a side effect. Behavior that was never observable, or a change the description names as intentional, does not count. Preserve the behavior, or make the change explicit in the description and cover it with a test.

## No unguarded destruction

An operation that deletes, drops, truncates, overwrites, force-pushes, or otherwise destroys data or history must be guarded: scoped to what was intended, confirmed or gated, and recoverable where the codebase's norms require it. This catches a `DROP TABLE` or `rm -rf` on a path built from input, a migration that deletes rows without a filter, an overwrite of a user file with no backup, and a cleanup job with an unbounded match. Deleting a temp file the code created, or a destructive step behind an explicit flag the change documents, does not count. Add the guard — a filter, a dry-run, a confirmation, a soft delete — and make the blast radius visible in the code.

## Variable misuse

Code that reads, returns, or passes a value must use the specific variable the surrounding logic actually computed or was given, not a different, same-typed value that happens to be in scope and compiles. This catches modifying a local copy but returning the original, comparing a value against the wrong parameter, passing an outer-scope value where the inner one was intended, and using a start-time variable where the end-time variable was meant. A deliberate use of an outer, default, or fallback value, explained by the surrounding code or its comments, does not count. Trace each identifier back to where it was last assigned or received, and confirm it is the one the current line's logic actually needs, not merely a plausible name in scope.

## Boolean polarity

A change must not invert the sense of a guard, gate, or permission check: using AND where the logic needs OR, negating a condition that should not be negated, or swapping which branch a flag or flag default takes. This catches a permission check that requires two roles when either should qualify, a feature flag guarded by the wrong sign, and a boolean whose default silently flipped. A condition correctly negated to match a genuinely different requirement, explained by the change, does not count. Write out in plain language which inputs the condition is meant to admit and which it is meant to reject, then verify the operator and every negation against that sentence.

## Normalized comparisons

Values compared for equality that can legitimately vary in case, whitespace, or encoding — hostnames, tokens, codes, anything a person might type — must be normalized the same way on both sides before comparing, unless the comparison is deliberately case-sensitive by contract. This catches one side of a comparison being lower-cased while the other is a raw user-supplied value, and validation done with indexOf or substring matching that silently fails on a differently-cased but equivalent input. A comparison against a value already guaranteed to be in canonical form upstream does not count. Normalize both operands the same way immediately before the comparison, and say in the code where that normalization happens.

## Boundary arithmetic

Code that computes an index, offset, or window boundary from external input or mutable state must be checked against the edge cases: the first element, the last element, an empty collection, and an offset that lands exactly on a limit. This catches a negative-offset branch that slices a collection with a negative start index, a paginator that produces an empty page when offset equals total, and a time window whose lower and upper bounds are computed from two different reference points. A boundary already validated earlier on the same path does not count. Name the boundary case explicitly in a comment or test, and verify the arithmetic against it, not only against the common case.

## Falsy-but-present values

A value that evaluates as falsy in the language but is semantically present and valid — zero, an empty string, an empty collection, `false` itself — must not be treated the same as absent, unset, or not-provided. This catches a rate, count, or flag of zero being skipped by a truthiness check meant to test "was this given", and an empty string being treated as no answer when it is a valid one. A value where the language's falsy state truly does mean nothing was provided, stated by the code's own contract, does not count. Check for absence explicitly — a null, undefined, or key-exists check — instead of relying on truthiness.

## Operation completeness

An operation described or named as a single action but implemented as several related writes, deletes, or external calls must complete every one of them on every path, or explicitly roll back the ones that already succeeded. This catches canceling an external side effect, such as a scheduled email or a charge, without deleting the local record that was tracking it, and a multi-table write that updates one table but not a table that depends on it. A step intentionally deferred to a background job or a later change, stated in the code, does not count. Enumerate every step the operation implies and verify each one runs on every return path, including early returns and error branches.

## Test synchronization

A test that waits for asynchronous, scheduled, or background work must synchronize on an observable condition — a callback, a flag, a polled state — not a fixed-duration sleep or delay that can pass or fail independent of whether the work actually happened. This catches a fixed sleep standing in for a wait-until-condition, and a sleep whose target function was mocked or monkeypatched elsewhere in the test so it no longer waits at all. A sleep that bounds genuinely time-based behavior under test, not a proxy for another event, does not count. Replace the fixed wait with a wait on the actual signal the test cares about.

## Advisory: Accurate identifiers

Names — of functions, methods, tests, variables, and properties — must be spelled correctly and must describe what they actually name or test. This catches a misspelled method name that ships as part of a public API, and a test name that describes a different input, behavior, or expectation than what the test body actually exercises. A name that is imprecise but not misleading, or a pre-existing name outside the change's scope, does not count. Fix the spelling, or rename so the identifier matches what the code or test actually does.
