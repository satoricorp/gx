# Code Review Guidelines

All merge requests for GitLab CE and EE must go through code review to ensure the code is effective, understandable, maintainable, and secure.

## Getting your merge request reviewed, approved, and merged

Before you begin, familiarize yourself with the contribution acceptance criteria.

Have your code **reviewed** by a
reviewer
from your group or a domain expert.

For small, straightforward changes, you can skip the reviewer step and go directly to a maintainer. Examples of small and straightforward changes:

- Fixing a typo or making small copy changes.
- A tiny refactor that doesn’t change any behavior.
- Removing a feature flag that has been default-enabled for more than one month.
- Removing unused methods or classes.
- A well-understood logic change that requires changes to fewer than five lines of code.

Otherwise, have a reviewer in each category the MR touches before passing
to a maintainer. For security assistance, include `@gitlab-com/gl-security/appsec`.

After the reviewer approves, a maintainer reviews and merges. The last required approver merges.

For CODEOWNERS-required approvals, seek domain-specific approvals before generic ones. Domain-specific approvers who are also maintainers should review both aspects and approve once.

### Approval guidelines

As described in the section on the responsibility of the maintainer below, you are recommended to get your merge request approved and merged by maintainers with domain expertise. The optional approval of the first reviewer is not covered here. However, your merge request should be reviewed by a reviewer before passing it to a maintainer as described in the overview section.

| If your merge request includes | It must be approved by a |
|---|---|
| `~backend`changes1 | Backend maintainer. |
| `~database`migrations or changes to expensive queries2 | Database maintainer. Refer to the database review guidelines for more details. |
| `~workhorse`changes | Workhorse maintainer. |
| `~frontend`changes1 | Frontend maintainer. |
| `~UX`user-facing changes3 | Product Designer. Refer to the design and user interface guidelines for details. |
| Adding a new JavaScript library 1 | - Frontend Design System member if the library significantly increases the bundle size. - A legal department member if the license used by the new library hasn’t been approved for use in GitLab. More information about license compatibility can be found in our GitLab Licensing and Compatibility documentation. |
| A new dependency or a file system change | - Distribution team member. See how to work with the Distribution team for more details. - For RubyGems, request an AppSec review. |
| `~documentation`or`~UI text`changes | Technical writer based on assignments in the appropriate DevOps stage group. |
| Changes to development guidelines | Follow the review process and get the approvals accordingly. |
| Changes to AI instruction files under `.ai/` | AI harness DRI. Refer to the AI instruction files review guidelines for more details. |
| End-to-end andnon-end-to-end changes4 | Software Engineer in Test. |
| Only End-to-end changes 4orif the MR author is a Software Engineer in Test | Quality maintainer. |
| A new or updated application limit | Product manager. |
| Analytics Instrumentation (telemetry or analytics) changes | Analytics Instrumentation engineer. |
| A new service to GitLab (Puma, Sidekiq, Gitaly are examples) | Product manager. See the process for adding a service component to GitLab for details. |
| Changes related to authentication | Manage:Authentication. Check the code review section on the group page for more details. Patterns for files known to require review from the team are listed in the in the `Authentication`section of the`CODEOWNERS`file, and the team will be listed in the approvers section of all merge requests that modify these files. |
| Changes related to custom roles or policies | Manage:Authorization Engineer. |

- Specs other than JavaScript specs are considered - `~backend`code. Haml markup is considered- `~frontend`code. However, Ruby code in Haml templates is considered- `~backend`code. When in doubt, request both a frontend and backend review.- For Haml template changes specifically: - **Request backend review**when changes include Ruby logic, method calls, variable assignments, conditionals, loops, data preparation, security checks, or any server-side processing in the template.
- **Request frontend review**when changes affect DOM structure, CSS classes, HTML attributes, accessibility features, user interactions, and responsive design, or visual presentation.
- **Request both reviews**for complex changes that involve both Ruby logic and significant UI modifications, or when backend and frontend are intertwined (such as when backend serves data that is consumed by Vue or JavaScript), to ensure both the backend functionality and frontend user experience are properly evaluated.- **Example:**A Haml template that calls Ruby methods to prepare data attributes for a Vue.js component (for example,- `project_id: @project&.to_global_id`) would benefit from a backend review for Ruby logic correctness and frontend review for the component integration.

- We encourage you to seek guidance from a database maintainer if your merge request is potentially introducing expensive queries. It is most efficient to comment on the line of code in question with the SQL queries so they can give their advice.
- User-facing changes include both visual changes (regardless of how minor), and changes to the rendered DOM which impact how a screen reader may announce the content. Groups that do not have dedicated Product Designers do not require a Product Designer to approve feature changes, unless the changes are community contributions.
- End-to-end changes include all files in the - `qa`directory.

### Reviewing a merge request

Understand why the change is necessary (fixes a bug, improves the user experience, refactors the existing code). Then:

- Be thorough to reduce the number of iterations.
- Communicate which ideas you feel strongly about and those you don’t.
- Identify ways to simplify the code while still solving the problem.
- Offer alternative implementations, but assume the author already considered them. (“What do you think about using a custom validator here?”)
- Seek to understand the author’s perspective.
- Check out the branch and test the changes locally. For MRs requiring significant GDK modifications, consider requesting screenshots, videos, or domain-expert verification instead. Your testing might result in opportunities to add automated tests.
- If you don’t understand a piece of code, *say so*.
- Use the Conventional Comment format to convey intent.
Mark non-mandatory suggestions as (`**non-blocking:**`). When only non-blocking suggestions remain, move the MR to the next stage rather than waiting.
- Ensure there are no open dependencies. Check linked issues for blockers. If blocked by open MRs, set an MR dependency.
- After a round of line notes, post a summary note such as “Looks good to me” or “Just a couple things to address.”
- Let the author know if changes are required following your review.

**If the merge request is from a fork, also check the additional guidelines for community contributions.**

### GitLab-specific concerns

GitLab is used in a lot of places, from Omnibus packages, to source installations. GitLab.com itself is a large Enterprise Edition instance. This has some implications:

- **Query changes**should be tested to ensure that they don’t result in worse performance at the scale of GitLab.com. See database review guidelines.
- **Database migrations**must be:- Reversible.
- Performant at the scale of GitLab.com - ask a maintainer to test the migration on the staging environment if you aren’t sure.
- The correct migration type. See the guidance to choose which migration type.

- **Sidekiq workers**cannot change in a backwards-incompatible way.
- **Cached values**may persist across releases. If you are changing the type a cached value returns (say, from a string or nil to an array), change the cache key at the same time.
- **Settings**should be added as a last resort. See Adding a new setting to GitLab Rails.
- **File system access**is not possible in a cloud-native architecture. Ensure that we support object storage for any file storage we need to perform. For more information, see the uploads documentation.

### The responsibility of the merge request author

You are the directly responsible individual (DRI) for finding the best solution. Stay as the assignee throughout the review lifecycle. If you cannot set yourself as an assignee, ask a reviewer to do it.

Do not submit a merge request that is too large, fixes multiple issues, has high complexity, or implements more than one feature.

If the MR touches multiple CODEOWNERS sections, to minimize the approvals needed consider splitting the MR by concerns to parallelize reviews. Before requesting a maintainer review, confirm:

- The MR solves the intended problem in the most appropriate way.
- All requirements are satisfied.
- No remaining bugs, logic problems, uncovered edge cases, or known vulnerabilities exist.

Self-review your MR following the Code Review guidelines. Add inline comments on lines where you made decisions or trade-offs, or where context helps the reviewer understand the code.

Involve domain experts, product managers, UX designers, and database specialists as appropriate. If you are unsure whether your MR needs a domain expert review, it does.

For features spanning 10 or more MRs, work with your EM or Staff Engineer to identify consistent maintainer who share the context.

If your MR touches multiple domains, request a review from an expert in each domain.

Before requesting review, add MR diff comments for:

- Added linting rules (RuboCop, JS, and so on).
- Added libraries (Ruby gems, JS libs, and so on).
- Links to parent classes or methods, where not obvious.
- Benchmarking results.
- Potentially insecure code.

Ensure reviewers have access to any projects, snippets, or assets needed to validate the solution.

When assigning reviewers, comment to specify which domain each reviewer should focus on. This avoids ambiguity when a team member has expertise in multiple areas. For examples, see MR 75921 and MR 109500.

Only add `TODO` comments to source code if a reviewer requires it.
If you add a `TODO`, include a link to the relevant issue.

Write comments that explain why, not only what the code does.

Request maintainer reviews only when tests pass. If tests are failing, explain why in a comment. Contact maintainers by email or Slack only for immediate requests. In all other cases, adding them as a reviewer is sufficient.

### The responsibility of the reviewer

Reviewers are responsible for reviewing the specifics of the chosen solution.

If unavailable within the Review-response SLO, inform the author, find a replacement using the Review Workload Dashboard, and assign them.

When confident the MR meets all contribution acceptance criteria:

- Select **Approve**.
- `@`mention the author to notify them.
- Request a review from a maintainer with domain expertise, or follow the Reviewer roulette suggestion.

### The responsibility of the maintainer

Maintainers are responsible for the overall health, quality, and consistency of the GitLab codebase. Their reviews focus on architecture, code organization, separation of concerns, tests, DRYness, consistency, and readability.

Maintainers are the DRI for ensuring MRs reasonably meet acceptance criteria.

A maintainer makes sound judgements when evaluating the impact of an MR. If a maintainer feels that an MR is not able to merged, it is their responsibility to say so. The maintainer is also the expert adviser who knows when to pull in others for a second opinion.

When a maintainer approves an MR, they are taking responsibility alongside the author. This means that when there is a production incident, the maintainer may get paged to help resolve issues.

Certain merge requests may target a stable branch. For an overview of how to handle these requests, see the patch release runbook.

## Best practices

### Domain experts

Domain experts are team members who have substantial experience with a specific technology, product feature, or area of the codebase. Team members are encouraged to self-identify as domain experts and add it to their team profiles.

When self-identifying as a domain expert, it is recommended to assign the MR changing the `.yml` file to be merged by an already established Domain Expert or a corresponding Engineering Manager.

We make the following assumption with regards to automatically being considered a domain expert:

- Team members working in a specific stage/group (for example, create: source code) are considered domain experts for that area of the app they work on.
- Team members working on a specific feature (for example, search) are considered domain experts for that feature.

We default to assigning reviews to team members with domain expertise for code reviews. UX reviews default to the recommended reviewer from the Review Roulette. Due to designer capacity limits, areas not supported by a Product Designer will no longer require a UX review unless it is a community contribution. When a suitable domain expert isn’t available, you can choose any team member to review the MR, or follow the Reviewer roulette recommendation (see above for UX reviews). Double check if the person is OOO before assigning them.

To find a domain expert:

- In the Merge Request approvals widget, select View eligible approvers. This widget shows recommended and required approvals per area of the codebase. These rules are defined in Code Owners.
- View the list of team members who work in the stage or group related to the merge request.
- View team members’ domain expertise on the engineering projects page or on the GitLab team page. Domains are self-identified, so use your judgment to map the changes on your merge request to a domain.
- Look for team members who have contributed to the files in the merge request. View the logs by running `git log <file>`.
- Look for team members who have reviewed the files. You can find the relevant merge request by:- Getting the commit SHA by using `git log <file>`.
- Navigating to `https://gitlab.com/gitlab-org/gitlab/-/commit/<SHA>`.
- Selecting the related merge request shown for the commit.

- Getting the commit SHA by using

### Reviewer roulette

Reviewer roulette is an internal tool for GitLab.com, not available on customer installations.

The Danger bot picks a reviewer and maintainer for each codebase area your MR touches. Override the suggestion if you know a better fit.

The roulette skips people whose status contains `OOO`, `PTO`, `Parental Leave`, `Friends and Family`, or `Conference`, or who are at review capacity (set via a number status emoji: 2️⃣–5️⃣).

### Acceptance checklist

This checklist encourages the authors, reviewers, and maintainers of merge requests (MRs) to confirm changes were analyzed for high-impact risks to quality, performance, reliability, security, observability, and maintainability.

Using checklists improves quality in software engineering. This checklist is a straightforward tool to support and bolster the skills of contributors to the GitLab codebase.

#### Quality

For further quality guidelines, see testing.

- You have self-reviewed this MR per code review guidelines.
- The code follows the software design guidelines.
- Ensure automated tests exist following the testing pyramid. Add missing tests or create an issue documenting testing gaps.
- You have considered the technical impacts on GitLab.com, Dedicated and self-managed.
- You have considered the impact of this change on the frontend, backend, and database portions of the system where appropriate and applied the `~ux`,`~frontend`,`~backend`, and`~database`labels accordingly.
- You have tested this MR in all supported browsers, or determined that this testing is not needed.
- You have confirmed that this change is backwards compatible across updates, or you have decided that this does not apply.
- You have properly separated EE content (if any) from FOSS. Consider running the CI pipelines in a FOSS context.
- You have considered that existing data may be surprisingly varied. For example, if adding a new model validation, consider making it optional on existing data.
- You have fixed flaky tests related to this MR, or have explained why they can be ignored. Flaky tests have error `Flaky test '<path/to/test>' was found in the list of files changed by this MR.`but can be in jobs that pass with warnings.

#### Performance, reliability, and availability

- You are confident that this MR does not harm performance, or you have asked a reviewer to help assess the performance impact. (Merge request performance guidelines)
- You have added information for database reviewers in the MR description, or you have decided that it is unnecessary.
- You have considered the availability and reliability risks of this change.
- You have considered the scalability risk based on future predicted growth.
- You have considered the performance, reliability, and availability impacts of this change on large customers who may have significantly more data than the average customer.
- You have considered the performance, reliability, and availability impacts of this change on customers who may run GitLab on the minimum system.
- You are confident that this change is compatible with the Cells architecture. For more information, see Cells development principles.

#### Observability instrumentation

- You have included enough instrumentation to facilitate debugging and proactive performance improvements through observability. See example of adding feature flags, logging, and instrumentation.

#### Documentation

- You have included changelog trailers, or you have decided that they are not needed.
- You have added/updated documentation or decided that documentation changes are unnecessary for this MR.

#### Security

- You have confirmed that if this MR contains changes to processing or storing of credentials or tokens, authorization, and authentication methods, or other items described in the security review guidelines, you have added the `~security`label and you have`@`-mentioned`@gitlab-com/gl-security/appsec`.
- You have reviewed the documentation regarding internal application security reviews for **when**and**how**to request a security review and requested a security review if this is warranted for this change.
- If there are security scan results that are blocking the MR (due to the merge request approval policies):- For true positive findings, they should be corrected before the merge request is merged. This will remove the AppSec approval required by the merge request approval policy.
- For false positive findings, something that should be discussed for risk acceptance, or anything questionable, ping `@gitlab-com/gl-security/appsec`.

#### Deployment

- You have considered using a feature flag for this change because the change may be high risk.
- If you are using a feature flag, you plan to test the change in staging before you test it in production, and you have considered rolling it out to a subset of production customers before rolling it out to all customers.
- You have informed the Infrastructure department of a default setting or new setting change per definition of done, or decided that this is unnecessary.

#### Compliance

- You have confirmed that the correct MR type label has been applied.

### Participating in code review

- Be kind.
- Accept that many programming decisions are opinions. Discuss trade-offs and resolve quickly.
- Ask questions. Make suggestions, not demands.
- Be explicit. People don’t always understand your intentions online.
- Be humble. Consider a one-on-one call for lengthy misunderstandings and post a follow-up summary.
- Mention the person directly when a comment is addressed specifically to them.
- Read through the entire diff before your first push. Check for unrelated changes and debug code.
- Write a detailed description per the merge request guidelines.
- Don’t take feedback personally. The review is of the code and its impact on production systems.
- Explain why the code exists, not just what it does.
- Try to respond to every comment. Only resolve threads you have fully addressed. If a comment can be addressed in a follow-up issue, work with the maintainer on a path forward.
- Re-request review once you are ready for another round.
- Address all GitLab Duo review comments before requesting a review from human reviewers.

### For authors: getting changes merged faster

- Follow best practices: write clear descriptions, add screenshots and validation steps, address
`dangerbot`comments, and complete the acceptance checklist.
- Follow GitLab patterns. Long discussions delay merging. Consider following the documented approach, then open a separate MR to propose changes to best practices.
- Keep MRs small. Around 200 lines is a good target.- Smaller MRs are reviewed faster and have fewer blocking discussions.
- Use stacked diffs for sequential MRs.
- Split changes so only one maintainer is required per MR (for example, ship database changes before the feature).
- UI with mocked data must be behind a feature flag.

- Minimize the number of reviewers. A database reviewer can also review backend; a fullstack engineer can cover frontend and backend.
- Know your maintainers and assign domain experts. Maintainers prioritize MRs in areas they know well.

### Merging a merge request

Before merging:

- Set the milestone.
- Confirm the correct MR type label is applied.
- Resolve warnings and errors from Danger bot, code quality, and other reports. Post a comment if merging with any failed job.

At least one maintainer must approve before merging. Authors and people who add commits cannot approve their own MR.

If the final approver did not set auto-merge, the MR author may merge their own MR if all required approvals are in place and they have merge rights. This aligns with the GitLab bias for action value.

When ready to merge:

If the merge request is from a fork, also check the guidelines for community contributions.

- Use Squash and merge only if the author set it or the commit history is messy.
- In the **Pipelines**tab, select**Run pipeline**, then enable**Auto-merge**on the**Overview**tab.- Do not merge when the default branch is broken, except in specific cases.
- Start a new pipeline if the latest one was created before approval and the MR has backend changes.
- You may skip a new pipeline if the latest merged results pipeline was created less than 16 hours ago (72 hours for stable branches).

### Community contributions

**Review all changes thoroughly for malicious code before starting a
merged results pipeline.**

When reviewing community MRs:

- Scrutinize new dependencies (e.g. `Gemfile.lock`,`yarn.lock`). They could introduce malicious packages.
- Review links and images, especially in documentation MRs.
- When in doubt, ask `@gitlab-com/gl-security/appsec`to review before starting any pipeline.
- Set the milestone only when the MR is likely to merge in the current milestone.

#### Taking over a community merge request

When an MR needs changes but the author is unresponsive or unable to finish:

- Comment that you (the merge request coach) are taking over.
- Add the `~"coach will finish"`label.
- Create a feature branch from main and merge their branch into it.
- Open a new MR, link the community MR, and add the `~"Community contribution"`label.
- Notify the contributor and follow the regular review process.

### Finding the right balance

Finding the right balance in how deeply to review requires sound judgement. Keep in mind:

- Finding bugs is important, but good design reduces future complexity.
- Enforce code style through automation rather than review comments.
- For non-blocking suggestions, consider approving the MR before passing it back. This reduces time-to-merge.
- Distinguish between doing things right and doing things right now. For example, avoid requiring major refactors in an urgent security fix.
- Doing things well today is usually better than doing something perfectly tomorrow.

## Troubleshooting failing pipelines

- **Unrelated test failures**: Check if the failure also happens on the default branch. If so, wait for the broken master fix, then re-run the pipeline for the MR.
- `danger-review`job failed

For help, comment `@gitlab-bot help` on the MR, or ask in the
Community Discord `contribute` channel.

### Credits

Based on the `thoughtbot` code review guide.
