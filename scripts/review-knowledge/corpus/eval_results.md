# Review Corpus Eval Results

generated_at: 2026-07-08
baseline_namespace: totality-review-knowledge
candidate_namespace: totality-review-knowledge

queries: 50
v1_recall_at_10: 45/50
v2_recall_at_10: 48/50
precedence_violations: 0
status: pass

## Query Results
- when is it acceptable to use panic in go: v1=hit v2=hit expected=google-go-style,uber-go-guide
- go error wrapping with context: v1=hit v2=hit expected=google-go-style,uber-go-guide,go-code-review-comments
- go init function side effects: v1=hit v2=hit expected=google-go-style,uber-go-guide
- go interface naming and small interfaces: v1=hit v2=hit expected=google-go-style,go-code-review-comments
- go context should not be stored in structs: v1=hit v2=hit expected=go-code-review-comments,google-go-style
- go dependency vulnerability reachable functions: v1=hit v2=hit expected=govulncheck,go-security
- go security best practices for untrusted input: v1=hit v2=hit expected=go-security,owasp-asvs
- go table tests for changed behavior: v1=hit v2=hit expected=google-go-style,go-code-review-comments
- go generics guidance should not come from effective go: v1=hit v2=hit expected=google-go-style,uber-go-guide
- go code review comments for goroutine lifetime: v1=hit v2=hit expected=go-code-review-comments,google-go-style
- typescript strict null checks configuration: v1=hit v2=hit expected=tsconfig-reference,typescript-eslint-typed-linting
- typescript unsafe any should be linted: v1=hit v2=hit expected=typescript-eslint-typed-linting,eslint-rules
- typescript default export style guidance: v1=hit v2=hit expected=google-tsguide,typescript-eslint-typed-linting
- typescript module imports and api shape: v1=hit v2=hit expected=google-tsguide
- typescript node security best practices: v1=hit v2=hit expected=nodejs-security,owasp-asvs
- typescript eslint rule pages for mechanical style: v1=hit v2=hit expected=eslint-rules,typescript-eslint-typed-linting
- typescript typed linting performance tradeoff: v1=hit v2=hit expected=typescript-eslint-typed-linting
- javascript security command injection node: v1=hit v2=hit expected=nodejs-security,owasp-cheatsheets
- typescript tsconfig unsafe compiler option: v1=miss v2=hit expected=tsconfig-reference
- typescript design-level style with no lint rule: v1=miss v2=hit expected=google-tsguide
- python line length should follow formatter: v1=hit v2=hit expected=ruff,pep8
- python google docstring convention D211 D212: v1=hit v2=hit expected=ruff,google-pyguide
- python type checking public api boundaries: v1=hit v2=hit expected=mypy,typing-python-spec
- python typing semantics protocol optional: v1=hit v2=hit expected=typing-python-spec,mypy
- python security subprocess shell true: v1=miss v2=miss expected=bandit,owasp-cheatsheets
- python exceptions style google guide: v1=hit v2=hit expected=google-pyguide,pep8
- python imports formatting lint: v1=hit v2=hit expected=ruff,google-pyguide,pep8
- python docstring D203 D213 should not outrank google convention: v1=hit v2=hit expected=ruff,google-pyguide
- python mypy strict optional review: v1=hit v2=hit expected=mypy,typing-python-spec
- python bandit hardcoded password rule: v1=hit v2=hit expected=bandit
- rust unsafe code invariants: v1=hit v2=hit expected=rust-nomicon,rust-unsafe-code-guidelines
- rust public api naming traits errors: v1=hit v2=miss expected=rust-api-guidelines
- rust clippy lint complexity: v1=hit v2=hit expected=clippy
- rust miri undefined behavior tests: v1=hit v2=hit expected=miri,rust-nomicon
- rust dependency advisory database: v1=hit v2=hit expected=rustsec
- rust unsafe aliasing and lifetimes: v1=hit v2=hit expected=rust-nomicon,rust-unsafe-code-guidelines
- rust non normative unsafe code guidelines severity: v1=hit v2=hit expected=rust-unsafe-code-guidelines
- rust api documentation future proofing: v1=hit v2=hit expected=rust-api-guidelines
- rust ffi memory safety: v1=hit v2=hit expected=rust-nomicon,miri
- rust correctness lint should be style tier but retrievable: v1=hit v2=hit expected=clippy
- preventing sql injection in parameterized queries: v1=hit v2=hit expected=owasp-asvs,owasp-sqli-cheatsheet,mysql-prepared-statements
- postgres transaction isolation anomaly: v1=hit v2=hit expected=postgres-txn-isolation
- postgres explain query plan performance: v1=hit v2=hit expected=postgres-explain
- dbt sql style canonical source: v1=hit v2=hit expected=dbt-sql-style,sqlfluff
- sqlfluff rule pages for sql lint: v1=hit v2=hit expected=sqlfluff
- database security least privilege: v1=miss v2=hit expected=owasp-db-cheatsheet,owasp-asvs
- owasp authentication verification requirements: v1=hit v2=hit expected=owasp-asvs
- cwe top 25 cross reference not spine: v1=hit v2=hit expected=cwe-top25,owasp-asvs
- software supply chain scorecard branch protection: v1=hit v2=hit expected=openssf-scorecard,slsa
- how large should a change be for review: v1=miss v2=hit expected=google-eng-practices,smartbear-best-practices
