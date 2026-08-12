# Review Corpus Eval Results

generated_at: 2026-08-12
baseline_namespace: gx-review-knowledge
candidate_namespace: review-corpus-v2

queries: 112
v1_recall_at_10: 71/112
v2_recall_at_10: 112/112
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
- typescript tsconfig unsafe compiler option: v1=hit v2=hit expected=tsconfig-reference
- typescript design-level style with no lint rule: v1=hit v2=hit expected=google-tsguide
- python line length should follow formatter: v1=hit v2=hit expected=ruff,pep8
- python google docstring convention D211 D212: v1=hit v2=hit expected=ruff,google-pyguide
- python type checking public api boundaries: v1=hit v2=hit expected=mypy,typing-python-spec
- python typing semantics protocol optional: v1=hit v2=hit expected=typing-python-spec,mypy
- python security subprocess shell true: v1=hit v2=hit expected=bandit,owasp-cheatsheets
- python exceptions style google guide: v1=hit v2=hit expected=google-pyguide,pep8
- python imports formatting lint: v1=hit v2=hit expected=ruff,google-pyguide,pep8
- python docstring D203 D213 should not outrank google convention: v1=hit v2=hit expected=ruff,google-pyguide
- python mypy strict optional review: v1=hit v2=hit expected=mypy,typing-python-spec
- python bandit hardcoded password rule: v1=hit v2=hit expected=bandit
- rust unsafe code invariants: v1=hit v2=hit expected=rust-nomicon,rust-unsafe-code-guidelines
- rust public api naming traits errors: v1=hit v2=hit expected=rust-api-guidelines
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
- database security least privilege: v1=hit v2=hit expected=owasp-db-cheatsheet,owasp-asvs
- owasp authentication verification requirements: v1=hit v2=hit expected=owasp-asvs
- cwe top 25 cross reference not spine: v1=hit v2=hit expected=cwe-top25,owasp-asvs
- software supply chain scorecard branch protection: v1=hit v2=hit expected=openssf-scorecard,slsa
- how large should a change be for review: v1=hit v2=hit expected=google-eng-practices,smartbear-best-practices
- label review comments nitpick suggestion blocking non-blocking: v1=hit v2=hit expected=conventional-comments
- semgrep rule metadata category severity confidence cwe references: v1=hit v2=hit expected=semgrep-registry
- nist ssdf secure software development framework practices: v1=hit v2=hit expected=nist-ssdf-800-218
- concise checklist for developing more secure software mfa dependencies: v1=hit v2=hit expected=openssf-concise-guide
- shellcheck double quote variable word splitting globbing: v1=hit v2=hit expected=shellcheck,google-shellguide
- ruff rule code reference flake8 pycodestyle pyflakes: v1=hit v2=hit expected=ruff-rules,ruff
- bandit plugin severity confidence subprocess popen shell injection: v1=hit v2=hit expected=bandit-plugins,bandit
- c# async await exception handling conventions: v1=hit v2=hit expected=csharp-coding-conventions
- dotnet public api naming and design guidelines: v1=hit v2=hit expected=dotnet-framework-design-guidelines,csharp-coding-conventions
- dotnet secure coding untrusted input deserialization: v1=hit v2=hit expected=dotnet-secure-coding,owasp-deserialization-cheatsheet
- roslyn analyzer rule categories quality security: v1=hit v2=hit expected=dotnet-analyzer-categories,dotnet-code-analysis
- enable dotnet code analysis editorconfig severity: v1=hit v2=hit expected=dotnet-code-analysis
- asp.net core authentication authorization review: v1=hit v2=hit expected=aspnet-core-security,owasp-asvs
- effective dart naming and style guidance: v1=hit v2=hit expected=effective-dart
- dart avoid dynamic prefer typed public apis: v1=hit v2=hit expected=effective-dart,dart-linter-rules
- dart linter rule pages for mechanical style: v1=hit v2=hit expected=dart-linter-rules,effective-dart
- flutter widget rebuild performance const constructors: v1=hit v2=hit expected=flutter-perf-best-practices
- flutter app security hardening best practices: v1=hit v2=hit expected=flutter-security,owasp-asvs
- does the transactional annotation apply to a method called from within the same class: v1=miss v2=hit expected=spring-framework-reference
- spring boot configuration property precedence and profile overrides: v1=miss v2=hit expected=spring-boot-reference
- django run code only after the transaction commits: v1=miss v2=hit expected=django-topics
- useeffect dependency array and cleanup contract: v1=miss v2=hit expected=react-reference,react-learn
- you might not need an effect derived state during render: v1=miss v2=hit expected=react-learn
- next.js server component caching and revalidation semantics: v1=miss v2=hit expected=nextjs-docs
- nestjs request scoped provider and interceptor execution order: v1=miss v2=hit expected=nestjs-docs
- vue reactivity caveat mutating props and lifecycle order: v1=miss v2=hit expected=vue-guide
- express async handler error must be passed to next: v1=miss v2=hit expected=express-error-handling
- lazy initialization exception on a detached entity outside the persistence context: v1=miss v2=hit expected=hibernate-user-guide,jakarta-persistence-spec
- does merge cascade to associations for a detached jpa entity: v1=miss v2=hit expected=jakarta-persistence-spec,hibernate-user-guide
- active record callback runs outside the surrounding transaction after_commit: v1=miss v2=hit expected=rails-active-record-callbacks
- n+1 query includes preload eager_load difference: v1=miss v2=hit expected=rails-active-record-querying
- uniqueness validation race condition needs a database unique index: v1=miss v2=hit expected=rails-active-record-validations
- background job must be idempotent because delivery is at least once: v1=miss v2=hit expected=rails-active-job-basics,gcp-pubsub-exactly-once
- sqlalchemy session lifetime expire on commit detached instance: v1=miss v2=hit expected=sqlalchemy-orm
- connection pool size pre ping and recycle for serverless: v1=miss v2=hit expected=sqlalchemy-pooling,prisma-orm-docs
- prisma interactive transaction isolation and connection limit: v1=miss v2=hit expected=prisma-orm-docs
- gorm preload versus joins and hooks inside the transaction: v1=miss v2=hit expected=gorm-docs
- double checked locking requires volatile happens before: v1=miss v2=hit expected=java-memory-model-jls,java-util-concurrent
- do not pool virtual threads and avoid pinning in synchronized blocks: v1=miss v2=hit expected=java-virtual-threads
- unsynchronized map write data race in goroutine: v1=miss v2=hit expected=go-memory-model,go-race-detector
- context cancel must be called and deadlines propagate: v1=miss v2=hit expected=go-context-package,grpc-guides
- blocking call inside the asyncio event loop: v1=miss v2=hit expected=python-asyncio-dev
- fire and forget task garbage collected without a strong reference: v1=miss v2=hit expected=python-asyncio-tasks,python-asyncio-dev
- process.nexttick versus setimmediate microtask ordering: v1=miss v2=hit expected=node-event-loop,mdn-js-execution-model
- synchronous cpu work blocks the node event loop: v1=miss v2=hit expected=node-blocking-event-loop
- sync over async result deadlock and async void: v1=miss v2=hit expected=aspnet-async-guidance,dotnet-async-programming
- reusing a protobuf field number is a breaking change: v1=miss v2=hit expected=protobuf-dos-donts,protobuf-proto3-guide
- which proto changes break wire compatibility versus source compatibility: v1=miss v2=hit expected=buf-breaking-rules,protobuf-proto3-guide
- is adding a method to a public trait a major version bump: v1=miss v2=hit expected=cargo-semver,semver-spec
- adding a field to an exported struct in a go module: v1=miss v2=hit expected=go-module-compatibility
- binary versus source compatibility when changing a public dotnet api: v1=miss v2=hit expected=dotnet-library-change-rules
- retry storm needs exponential backoff with jitter and a budget: v1=hit v2=hit expected=google-sre-book,azure-cloud-design-patterns
- load shedding and cascading failure from overload: v1=miss v2=hit expected=google-sre-book
- rpc call without a deadline and retry policy: v1=miss v2=hit expected=grpc-guides
- idempotency key scope and replayed response for a retried write: v1=miss v2=hit expected=stripe-idempotency
- ack deadline redelivery and consumer side deduplication: v1=miss v2=hit expected=gcp-pubsub-exactly-once
- rubocop cop rationale and default configuration for a lint department: v1=hit v2=hit expected=rubocop-cops,ruby-style-guide
- rails mass assignment and unsafe dynamic render path warning: v1=hit v2=hit expected=brakeman-warning-types,rails-security
- rails migration reversibility and index creation hazards: v1=miss v2=hit expected=rails-active-record-migrations
- dependent destroy versus delete_all callback semantics: v1=miss v2=hit expected=rails-association-basics
- rails cache key based expiry and fragment invalidation: v1=miss v2=hit expected=rails-caching
- java multithreading rule catalog for shared mutable state: v1=miss v2=hit expected=pmd-java-rules,errorprone-bugpatterns
