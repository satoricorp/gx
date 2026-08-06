# Review policy

## high-risk paths

risk-path: internal/storage/** — Storage and schema changes can affect existing repositories and status/read-model correctness across gx.
risk-path: internal/github/** — This code talks to GitHub, so review auth, ownership, retry/error behavior, and whether user-authored PR content is preserved.
risk-path: internal/publication/** — Publication code coordinates gx Cloud, local artifacts, and GitHub updates. Review ordering and partial-failure behavior.
risk-path: internal/vcs/** — VCS service changes affect branch push, PR creation, and local metadata. Review state transitions and failure paths.
