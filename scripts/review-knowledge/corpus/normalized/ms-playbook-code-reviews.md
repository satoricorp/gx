# Code Reviews

Developers working on projects should conduct peer code reviews on every pull request (or check-in to a shared branch).

## Goals

Code review is a way to have a conversation about the code where participants will:

- **Improve code quality**by identifying and removing defects before they can be introduced into shared code branches.
- **Learn and grow**by having others review the code, we get exposed to unfamiliar design patterns or languages among other topics, and even break some bad habits.
- **Shared understanding**between the developers over the project's code.

## Resources

- Code review tools
- Google's Engineering Practices documentation: How to do a code review
- Best Kept Secrets of Peer Code Review

 Last update:
 August 26, 2024

# FAQ

This is a list of questions / frequently occurring issues when working with code reviews and answers how you can possibly tackle them.

## What Makes a Code Review Different from a PR?

A pull request (PR) is a way to notify a task is finished and ready to be merged into the main working branch (source of truth). A code review is having someone go over the code in a PR and validate it before it is merged, but, in general, code reviews can take place outside PRs too.

| Code Review | Pull Request |
|---|---|
| Source code focused | Intended to enhance and enable code reviews. Includes both source code but can have a broader scope (e.g., docs, integration tests, compiles) |
| Intended for early feedbackbefore submitting a PR | Not intended for early feedback. Created when author is ready to merge |
| Usually a synchronous review with faster feedback cycles (draft PRs as an exception). Examples: scheduled meetings, over-the-shoulder review, pair programming | Usually a tool assisted asynchronous review but can be elevated to a synchronous meeting when needed |

## Why do we Need Code Reviews?

Our peer code reviews are structured around best practices, to find specific kinds of errors. Much like you would still run a linter over mobbed code, you would still ask someone to make the last pass to make sure the code conforms to expected standards and avoids common pitfalls.

## PRs are Too Large, How can we Fix This?

Make sure you size the work items into small clear chunks, so the reviewer will be able to understand the code on their own. The team is instructed to commit early, before the full product backlog item / user story is complete, but rather when an individual item is done. If the work would result in an incomplete feature, make sure it can be turned off, until the full feature is delivered. More information can be found in Pull Requests - Size Guidance.

## How can we Expedite Code Reviews?

Slow code reviews might cause delays in delivering features and cause frustration amongst team members.

### Possible Actions you can Take

- Add a rule for PR turnaround time to your work agreement.
- Set up a slot after the standup to go through pending PRs and assign the ones that are inactive.
- Dedicate a PR review manager who will be responsible to keep things flowing by assigning or notifying people when PR got stale.
- Use tools to better indicate stale reviews - Customize ADO - Task Boards.

## Which Tools can I use to Review a Complex PR?

Checkout the Tools for help on how to perform reviews out of Visual Studio or Visual Studio Code.

## How can we Enforce the Code Review Policies?

By configuring Branch Policies , you can easily enforce code reviews rules.

## We Pair or Mob. How Should This Reflect in our Code Reviews?

There are two ways to perform a code review:

- Pair - Someone outside the pair should perform the code review. One of the other major benefits of code reviews is spreading knowledge about the code base to other members of the team that don't usually work in the part of the codebase under review.
- Mob - A member of the mob who spent less (or no) time at the keyboard should perform the code review.

# Inclusion in Code Review

Below are some points which emphasize why inclusivity in code reviews is important:

- Code reviews are an important part of our job as software professionals.
- In ISE we work with cross cultural teams from across the globe.
- How we communicate affects team morale.
- Inclusive code reviews welcome new developers and make them comfortable with the team.
- Rude or personal attacks doing code reviews alienate - people can unknowingly make rude comments when reviewing pull requests (PRs).

## Types and Examples of Non-Inclusive Code Review Behavior

- Inequitable review assignments.- Example: Assigning most reviews to few people and dismissing some members of the team altogether.

- Negative interpersonal interactions.- Example: Long arguments over subjective topics such as code style.

- Biased decision making.- Example: Comments about the developer and not the code. Assuming code from developer X will always be good and hence not reviewing it properly and vice versa.

## Examples of Inclusive Code Reviews

- Anyone and everyone in the team should be assigned PRs to review.
- Reviewer should be clear about what is an opinion, their personal preference, best practice or a fact. Arguments over personal preferences and opinions are mostly avoidable.
- Using inclusive language and tone in the code review comments. For example, being suggestive rather being prescriptive in the review comments is a good way to get the point across the table.
- It's a good practice for the author of a PR to thank the reviewer for the review, when they have contributed in improving the code or you have learnt something new.
- Using the sandwich method for recommending a code change to a new developer or a new customer: Sandwich the suggestion between 2 compliments. For example: "Great work so far, but I would recommend a few changes here. Btw, I loved the use of XYZ here, nice job!"

## Guidelines for the Author

- Aim to write a code that is easy to read, review and maintain.
- It’s important to ensure that whoever is looking at the code, whether that be the reviewer or a future engineer, can understand the motivations and how your code achieves its goals.
- Proactively asking for targeted help or feedback.
- Respond clearly to questions asked by the reviewers.
- Avoid huge commits by submitting incremental changes. Commits which are large and contain changes to multiple files will lead to unfair review of the code. Biased behavior of reviewers may kick in while reviewing such PRs. For e.g. a huge commit from a senior developer may get approved without thorough review whereas a huge commit from a junior developer may never get reviewed and approved.

## Guidelines for the Reviewer

- Assume positive intent from the author.
- Write clear and elaborate comments.
- Identify subjectivity, choice of coding and best practice. It is good to discuss coding style and subjective coding choices in some other forum and not in the PR. A PR should not become a ground to discuss subjective coding choices and having long arguments over it.
- If you do not understand the code properly, refrain from commenting e.g., "This code is incomprehensible". It is better to have a call with the author and get a basic understanding of their work.
- Be suggestive and not prescriptive. A reviewer should suggest changes and not prescribe changes, let the author decide if they really want to accept the changes proposed.

## Culture and Code Reviews

We in ISE, may come across situations in which code reviews are not ideal and often we are observing non inclusive code review behaviors. Its important to be aware of the fact that culture and communication style of a particular geography also influences how people interact over pull requests. In such cases, assuming positive intent of the author and reviewer is a good start to start analyzing quality of code reviews.

## Dealing with the Impostor Phenomenon

Impostor phenomenon is a psychological pattern in which an individual doubts their skills, talents, or accomplishments and has a persistent internalized fear of being exposed as a "fraud" - Wikipedia.

Someone experiencing impostor phenomenon may find submitting code for a review particularly stressful. It is important to realize that everybody can have meaningful contributions and not to let the perceived weaknesses prevent contributions.

Some tips for overcoming the impostor phenomenon for authors:

- Review the guidelines highlighted above and make sure your code change adhere to them.
- Ask for help from a colleague - pair program with an experienced colleague that you can learn from.

Some tips for overcoming the impostor phenomenon for reviewers:

- Anyone can have valuable insights.
- A fresh new pair of eyes are always welcome.
- Study the review until you have clearly understood it, check the corner cases and look for ways to improve it.
- If something is not clear, a simple specific question should be asked.
- If you have learnt something, you can always compliment the author.
- If possible, pair with someone to review the code so that you can establish a personal connection and have a more profound discussion about the code.

## Tools

Below are some tools which may help in establishing inclusive code review culture within our teams.

# Pull Request Template

```
# [Work Item ID](./link-to-the-work-item)
For more information about how to contribute to this repo, visit this [page](https://github.com/microsoft/code-with-engineering-playbook/blob/main/CONTRIBUTING.md)
## Description
---
> Should include a concise description of the changes (bug or feature), it's impact, along with a summary of the solution
## Steps to Reproduce Bug and Validate Solution
---
> Only applicable if the work is to address a bug. Please remove this section if the work is for a feature or story
> Provide details on the environment the bug is found, and detailed steps to recreate the bug.
> This should be detailed enough for a team member to confirm that the bug no longer occurs
## PR Checklist
---
> Use the check-list below to ensure your branch is ready for PR. If the item is not applicable, leave it blank.
- [ ] I have updated the documentation accordingly.
- [ ] I have added tests to cover my changes.
- [ ] All new and existing tests passed.
- [ ] My code follows the code style of this project.
- [ ] I ran the lint checks which produced no new errors nor warnings for my changes.
- [ ] I have checked to ensure there aren't other open Pull Requests for the same update/change.
## Does This Introduce a Breaking Change?
---
- [ ] Yes
- [ ] No
> If this introduces a breaking change, please describe the impact and migration path for existing applications below.
## Testing
---
> - Instructions for testing and validation of your code:
> - What OS was used for testing.
> - Which test sets were used.
> - Description of test scenarios that you have tried.
## Any Relevant Logs or Outputs
---
> - Use this section to attach pictures that demonstrates your changes working / healthy
> - If you are printing something show a screenshot
> - When you want to share long logs upload to:
> `(StorageAccount)/pr-support/attachments/(PR Number)/(yourFiles) using [Azure Storage Explorer](https://azure.microsoft.com/en-us/features/storage-explorer/)` or [portal.azure.com](https://portal.azure.com) and insert the link here.
## Other Information or Known Dependencies
---
> - Any other information or known dependencies that is important to this PR.
> - TODO that are to be done after this PR.
```

 Last update:
 August 22, 2024

# Pull Requests

Changes to any main codebase - main branch in Git repository, for example - must be done using pull requests (PR).

Pull requests enable:

- Code inspection - see Code Reviews
- Running automated qualification of the code- Linters
- Compilation
- Unit tests
- Integration tests etc.

The requirements of pull requests can and should be enforced by policies, which can be set in the most modern version control and work item tracking systems. See Evidence and Measures section for more information.

## General Process

- Implement changes based on the well-defined description and acceptance criteria of the task at hand
- Then, before creating a new pull request: * Make sure the code conforms with the agreed coding conventions * This can be partially automated using linters * Ensure the code compiles and runs without errors or warnings * Write and/or update tests to cover the changes and make sure all new and existing tests pass * Write and/or update the documentation to match the changes
- Once convinced the criteria above are met, create and submit a new pull request adhering to the pull request template
- Follow the code review process to merge the changes to the main codebase

The following diagram illustrates this approach.

```
sequenceDiagram
New branch->>+Pull request: New PR creation
Pull request->>+Code review: Review process
Code review->>+Pull request: Code updates
Pull request->>+New branch: Merge Pull Request
Pull request-->>-New branch: Delete branch
Pull request ->>+ Main branch: Merge after completion
New branch->>+Main branch: Goal of the Pull request
```
## Size Guidance

We should always aim to keep pull requests small. Small PRs have multiple advantages:

- They are easier to review; a clear benefit for the reviewers.
- They are easier to deploy; this is aligned with the strategy of release fast and release often.
- Minimizes possible conflicts and stale PRs.

However, we should keep PRs focused - for example around a functional feature, optimization or code readability and avoid having PRs that include code that is without context or loosely coupled. There is no right size, but keep in mind that a code review is a collaborative process, a big PRs could be difficult and therefore slower to review. We should always strive to have as small PRs as possible that still add value.

## Best Practices

Beyond the size, remember that every PR should:

- be consistent,
- not break the build, and
- include related tests as part of the PR.

Be consistent means that all the changes included on the PR should aim to solve one goal (ex. one user story) and be intrinsically related. Think of this as the Single-responsibility principle in terms of the whole project, the PR should have only one *reason to change* the project.

Start small, it is easier to create a small PR from the start than to break up a bigger one.

These are some strategies to keep PRs small depending on the "cause" of the inevitability, you could break the PR into self-container changes which still add value, release features that are hidden (see feature flag, feature toggling or canary releases) or break the PR into different layers (for example using design patterns like MVC or Observer/Subject). No matter the strategy.

## Pull Request Description

Well written PR descriptions helps maintain a clean, well-structured change history. While every team need not conform to the same specification, it is important that the convention is agreed upon at the start of the project.

One popular specification for open-source projects and others is the Conventional Commits specification, which is structured as:

```
<type>[optional scope]: <description>
[optional body]
[optional footer]
```
The `<type>` in this message can be selected from a list of types defined by the team, but many projects use the list of commit types from the Angular open-source project. It should be clear that `scope`, `body` and `footer` elements are **optional**, but having a required `type` and short description enables the features mentioned above.

See also Pull Request Template

## Resources

- Writing a great pull request description
- Review code-with pull requests (Azure DevOps)
- Collaborating with issues and pull requests (GitHub)
- Google approach to PR size
- Feature Flags
- Facebook approach to hidden features
- Conventional Commits specification
- Angular Commit types

# Code Review Tools

## Customize ADO

### Task Boards

- AzDO: Customize cards
- AzDO: Add columns on task board

### Reviewer Policies

- Setting required reviewer group in AzDO - Automatically include code reviewers

## Configuring Branch Policies

- AzDO: Configure branch policies
- AzDO: Configuring branch policies with the CLI tool:
- GitHub: Configuring protected branches

## VSCode

### GitHub: GitHub Pull Requests

Supports processing GitHub pull requests inside VS Code.

- Open the plugin from the **Activity Bar**
- Select **Assigned To Me**
- Select a PR
- Under **Description**you can choose to**Check Out**the branch and get into**Review Mode**and get a more integrated experience

### Azure DevOps: Azure DevOps Pull Requests

Supports processing Azure DevOps pull requests inside VS Code.

- Open the plugin from the **Activity Bar**
- Select **Assigned To Me**
- Select a PR
- Under **Description**you can choose to**Check Out**the branch and get into**Review Mode**and get a more integrated experience

## Visual Studio

The following extensions can be used to create an integrated code review experience in Visual Studio working with either GitHub or Azure DevOps.

### GitHub: GitHub Extension for Visual Studio

Provides extended functionality for working with pull requests on GitHub directly out of Visual Studio.

- View -> Other Windows -> GitHub
- Click on the **Pull Requests**icon in the task bar
- Double click on a pending pull request

### Azure DevOps: Pull Requests for Visual Studio

Work with pull requests on Azure DevOps directly out of Visual Studio.

- Open Team Explorer
- Click on **Pull Requests**
- Double-click a pull request - the **Pull Request Details**open
- Click on **Checkout**if you want to have the full change locally and have a more integrated experience
- Go through the changes and make comments

## Web

### Reviewable: Seamless multi-round GitHub reviews

Supports multi-round GitHub code reviews, with keyboard shortcuts and more. VS Code extension is in-progress.

- Visit the Review Dashboard to see reviews awaiting your action, that have new comments for you, and more.
- Select a Pull Request from that list.
- Open any file in your browser, in Visual Studio Code, or any editor you've configured by clicking on your profile photo in the top-right
- Select an editor under "External editor link template". VS Code is an option, but so is any editor that supports URI's.
- Review the diff on an overall or per-file basis, leaving comments, code suggestions, and more

 Last update:
 August 22, 2024

# Evidence and Measures

## Evidence

Many of the code quality assurance items can be automated or enforced by policies in modern version control and work item tracking systems. Verification of the policies on the main branch in Azure DevOps (AzDO) or GitHub, for example, may be sufficient evidence that a project team is conducting code reviews.

- The main branches in all repositories have branch policies. - Configure branch policies
- All builds produced out of project repositories include appropriate linters, run unit tests.
- Every bug work item should include a link to the pull request that introduced it, once the error has been diagnosed. This helps with learning.
- Each bug work item should include a note on how the bug might (or might not have) been caught in a code review.
- The project team regularly updates their code review checklists to reflect common issues they have encountered.
- Dev Leads should review a sample of pull requests and/or be co-reviewers with other developers to help everyone improve their skills as code reviewers.

## Measures

The team can collect metrics of code reviews to measure their efficiency. Some useful metrics include:

- Defect Removal Efficiency (DRE) - a measure of the development team's ability to remove defects prior to release
- Time metrics:- Time used preparing for code inspection sessions
- Time used in review sessions

- Lines of code (LOC) inspected per time unit/meeting

It is a perfectly reasonable solution to track these metrics manually e.g. in an Excel sheet. It is also possible to utilize the features of project management platforms - for example, AzDO enables dashboards for metrics including tracking bugs. You may find ready-made plugins for various platforms - see GitHub Marketplace for instance - or you can choose to implement these features yourself.

Remember that since defects removed thanks to reviews is far less costly compared to finding them in production, the cost of doing code reviews is actually negative!

# Process Guidance

## General Guidance

Code reviews should be part of the software engineering team process regardless of the development model. Furthermore, the team should learn to execute reviews in a timely manner. Pull requests (PRs) left hanging can cause additional merge problems and go stale resulting in lost work. Qualified PRs are expected to reflect well-defined, concise tasks, and thus be compact in content. Reviewing a single task should then take relatively little time to complete.

To ensure that the code review process is healthy, inclusive and meets the goals stated above, consider following these guidelines:

- Establish a service-level agreement (SLA) for code reviews and add it to your teams working agreement.
- Although modern DevOps environments incorporate tools for managing PRs, it can be useful to label tasks pending for review or to have a dedicated place for them on the task board - Customize AzDO task boards
- In the daily standup meeting check tasks pending for review and make sure they have reviewers assigned.
- Junior teams and teams new to the process can consider creating separate tasks for reviews together with the tasks themselves.
- Utilize tools to streamline the review process - Code review tools
- Foster inclusive code reviews - Inclusion in Code Review

## Measuring Code Review Process

If the team is finding that code reviews are taking a significant time to merge, and it is becoming a blocker, consider the following additional recommendations:

- Measure the average time it takes to merge a PR per sprint cycle.
- Review during retrospective how the time to merge can be improved and prioritized.
- Assess the time to merge across sprints to see if the process is improving.
- Ping required approvers directly as a reminder.

## Code Reviews Shouldn't Include too Many Lines of Code

It's easy to say a developer can review few hundred lines of code, but when the code surpasses certain amount of lines, the effectiveness of defects discovery will decrease and there is a lesser chance of doing a good review. It's not a matter of setting a code line limit, but rather using common sense. More code there is to review, the higher chances there are letting a bug sneak through. See PR size guidance.

## Automate Whenever Reasonable

Use automation (linting, code analysis etc.) to avoid the need for "nits" and allow the reviewer to focus more on the functional aspects of the PR. By configuring automated builds, tests and checks (something achievable in the CI process), teams can save human reviewers some time and let them focus in areas like design and functionality for proper evaluation. This will ensure higher chances of success as the team is focusing on the things that matter.

# Author Guidance

## Properly Describe Your Pull Request (PR)

- Give the PR a descriptive title, so that other members can easily (in one short sentence) understand what a PR is about.
- Every PR should have a proper description, that shows the reviewer what has been changed and why.

## Add Relevant Reviewers

- Add one or more reviewers (depending on your project's guidelines) to the PR. Ideally, you would add at least someone who has expertise and is familiar with the project, or the language used
- Adding someone less familiar with the project or the language can aid in verifying the changes are understandable, easy to read, and increases the expertise within the team
- In ISE code-with projects with a customer team, it is important to include reviewers from both organizations for knowledge transfer - Customize Reviewers Policy

## Be Open to Receive Feedback

Discuss design/code logic and address all comments as follows:

- Resolve a comment, if the requested change has been made.
- Mark the comment as "won't fix", if you are not going to make the requested changes and provide a clear reasoning- If the requested change is within the scope of the task, "I'll do it later" is not an acceptable reason!
- If the requested change is out of scope, create a new work item (task or bug) for it

- If you don't understand a comment, ask questions in the review itself as opposed to a private chat
- If a thread gets bloated without a conclusion, have a meeting with the reviewer (call them or knock on door)

## Use Checklists

When creating a PR, it is a good idea to add a checklist of objectives of the PR in the description. This helps the reviewers to focus on the key areas of the code changes.

## Link a Task to Your PR

Link the corresponding work items/tasks to the PR. There is no need to duplicate information between the work item and the PR, but if some details are missing in either one, together they provide more context to the reviewer.

## Code Should Have Annotations Before the Review

If you can't avoid large PRs, include explanations of the changes in order to make it easier for the reviewer to review the code, with clear comments the reviewer can identify the goal of every code block.

 Last update:
 August 22, 2024

# Reviewer Guidance

Since parts of reviews can be automated via linters and such, human reviewers can focus on architectural and functional correctness. Human reviewers should focus on:

- The correctness of the business logic embodied in the code.
- The correctness of any new or changed tests.
- The "readability" and maintainability of the overall design decisions reflected in the code.
- The checklist of common errors that the team maintains for each programming language.

Code reviews should use the below guidance and checklists to ensure positive and effective code reviews.

## General Guidance

### Understand the Code You are Reviewing

- Read every line changed.
- If we have a stakeholder review, it’s not necessary to run the PR unless it aids your understanding of the code.
- AzDO orders the files for you, but you should read the code in some logical sequence to aid understanding.
- If you don’t fully understand a change in a file because you don’t have context, click to view the whole file and read through the surrounding code or checkout the changes and view them in IDE.
- Ask the author to clarify.

### Take Your Time and Keep Focus on Scope

You shouldn't review code hastily but neither take too long in one sitting. If you have many pull requests (PRs) to review or if the complexity of code is demanding, the recommendation is to take a break between the reviews to recover and focus on the ones you are most experienced with.

Always remember that a goal of a code review is to verify that the goals of the corresponding task have been achieved. If you have concerns about the related, adjacent code that isn't in the scope of the PR, address those as separate tasks (e.g., bugs, technical debt). Don't block the current PR due to issues that are out of scope.

## Foster a Positive Code Review Culture

Code reviews play a critical role in product quality and it should not represent an arena for long discussions or even worse a battle of egos. What matters is a bug caught, not who made it, not who found it, not who fixed it. The only thing that matters is having the best possible product.

## Be Considerate

- Be positive – encouraging, appreciation for good practices.
- Prefix a “point of polish” with “Nit:”.
- Avoid language that points fingers like “you” but rather use “we” or “this line” -- code reviews are not personal and language matters.
- Prefer asking questions above making statements. There might be a good reason for the author to do something.
- If you make a direct comment, explain why the code needs to be changed, preferably with an example.
- Talking about changes, you can suggest changes to a PR by using the suggestion feature (available in GitHub and Azure DevOps) or by creating a PR to the author branch.
- If a few back-and-forth comments don't resolve a disagreement, have a quick talk with each other (in-person or call) or create a group discussion this can lead to an array of improvements for upcoming PRs. Don't forget to update the PR with what you agreed on and why.

## First Design Pass

### Pull Request Overview

- Does the PR description make sense?
- Do all the changes logically fit in this PR, or are there unrelated changes?
- If necessary, are the changes made reflected in updates to the README or other docs? Especially if the changes affect how the user builds code.

### User Facing Changes

- If the code involves a user-facing change, is there a GIF/photo that explains the functionality? If not, it might be key to validate the PR to ensure the change does what is expected.
- Ensure UI changes look good without unexpected behavior.

### Design

- Do the interactions of the various pieces of code in the PR make sense?
- Does the code recognize and incorporate architectures and coding patterns?

## Code Quality Pass

### Complexity

- Are functions too complex?
- Is the single responsibility principle followed? Function or class should do one ‘thing’.
- Should a function be broken into multiple functions?
- If a method has greater than 3 arguments, it is potentially overly complex.
- Does the code add functionality that isn’t needed?
- Can the code be understood easily by code readers?

### Naming/Readability

- Did the developer pick good names for functions, variables, etc?

### Error Handling

- Are errors handled gracefully and explicitly where necessary?

### Functionality

- Is there parallel programming in this PR that could cause race conditions? Carefully read through this logic.
- Could the code be optimized? For example: are there more calls to the database than need be?
- How does the functionality fit in the bigger picture? Can it have negative effects to the overall system?
- Are there security flaws?
- Does a variable name reveal any customer specific information?
- Is PII and EUII treated correctly? Are we logging any PII information?

### Style

- Are there extraneous comments? If the code isn’t clear enough to explain itself, then the code should be made simpler. Comments may be there to explain why some code exists.
- Does the code adhere to the style guide/conventions that we have agreed upon? We use automated styling like black and prettier.

### Tests

- Tests should always be committed in the same PR as the code itself (‘I’ll add tests next’ is not acceptable).
- Make sure tests are sensible and valid assumptions are made.
- Make sure edge cases are handled as well.
- Tests can be a great source to understand the changes. It can be a strategy to look at tests first to help you understand the changes better.

# YAML(Azure Pipelines) Code Reviews

## Style Guide

Developers should follow the YAML schema reference.

## Code Analysis / Linting

The most popular YAML linter is YAML extension. This extension provides YAML validation, document outlining, auto-completion, hover support and formatter features.

## VS Code Extensions

There is an Azure Pipelines for VS Code extension to add syntax highlighting and autocompletion for Azure Pipelines YAML to VS Code. It also helps you set up continuous build and deployment for Azure WebApps without leaving VS Code.

## YAML in Azure Pipelines Overview

When the pipeline is triggered, before running the pipeline, there are a few phases such as Queue Time, Compile Time and Runtime where variables are interpreted by their runtime expression syntax.

When the pipeline is triggered, all nested YAML files are expanded to run in Azure Pipelines. This checklist contains some tips and tricks for reviewing all nested YAML files.

These documents may be useful when reviewing YAML files:

**Key concepts overview**

- A trigger tells a Pipeline to run.
- A pipeline is made up of one or more stages. A pipeline can deploy to one or more environments.
- A stage is a way of organizing jobs in a pipeline and each stage can have one or more jobs.
- Each job runs on one agent. A job can also be agentless.
- Each agent runs a job that contains one or more steps.
- A step can be a task or script and is the smallest building block of a pipeline.
- A task is a pre-packaged script that performs an action, such as invoking a REST API or publishing a build artifact.
- An artifact is a collection of files or packages published by a run.

## Code Review Checklist

In addition to the Code Review Checklist you should also look for these Azure Pipelines YAML specific code review items.

### Pipeline Structure

- The steps are well understood and components are easily identifiable. Ensure that there is a proper description `displayName:`for every step in the pipeline.
- Steps/stages of the pipeline are checked in Azure Pipelines to have more understanding of components.
- In case you have complex nested YAML files, The pipeline in Azure Pipelines is edited to find trigger root file.
- All the template file references are visited to ensure a small change does not cause breaking changes, changing one file may affect multiple pipelines
- Long inline scripts in YAML file are moved into script files

### YAML Structure

- Re-usable components are split into separate YAML templates.
- Variables are separated per environment stored in templates or variable groups.
- Variable value changes in `Queue Time`,`Compile Time`and`Runtime`are considered.
- Variable syntax values used with `Macro Syntax`,`Template Expression Syntax`and`Runtime Expression Syntax`are considered.
- Variables can change during the pipeline, Parameters cannot.
- Unused variables/parameters are removed in pipeline.
- Does the pipeline meet with stage/job `Conditions`criteria?

### Permission Check & Security

- Secret values shouldn't be printed in pipeline, `issecret`is used for printing secrets for debugging
- If pipeline is using variable groups in Library, ensure pipeline has access to the variable groups created.
- If pipeline has a remote task in other repo/organization, does it have access?
- If pipeline is trying to access a secure file, does it have the permission?
- If pipeline requires approval for environment deployments, Who is the approver?
- Does it need to keep secrets and manage them, did you consider using Azure KeyVault?

### Troubleshooting Tips

- Consider Variable Syntax with Runtime Expressions in the pipeline. Here is a nice sample to understand Expansion of variables.

-
When we assign variable like below it won't set during initialize time, it'll assign during runtime, then we can retrieve some errors based on when template runs. `- task: AzureWebApp@1 displayName: 'Deploy Azure Web App : $(webAppName)' inputs: azureSubscription: '$(azureServiceConnectionId)' appName: '$(webAppName)' package: $(Pipeline.Workspace)/drop/Application$(Build.BuildId).zip startUpCommand: 'gunicorn --bind=0.0.0.0 --workers=4 app:app'`Error: After passing these variables as parameter, it loads values properly. `- template: steps-deployment.yaml parameters: azureServiceConnectionId: ${{ variables.azureServiceConnectionId }} webAppName: ${{ variables.webAppName }}``- task: AzureWebApp@1 displayName: 'Deploy Azure Web App :${{ parameters.webAppName }}' inputs: azureSubscription: '${{ parameters.azureServiceConnectionId }}' appName: '${{ parameters.webAppName }}' package: $(Pipeline.Workspace)/drop/Application$(Build.BuildId).zip startUpCommand: 'gunicorn --bind=0.0.0.0 --workers=4 app:app'`

-
Use `issecret`for printing secrets for debugging`echo "##vso[task.setvariable variable=token;issecret=true]${token}"`

# Bash Code Reviews

## Style Guide

Developers should follow Google's Bash Style Guide.

## Code Analysis / Linting

Projects must check bash code with shellcheck as part of the CI process. Apart from linting, shfmt can be used to automatically format shell scripts. There are few vscode code extensions which are based on shfmt like shell-format which can be used to automatically format shell scripts.

## Project Setup

### vscode-shellcheck

Shellcheck extension should be used in VS Code, it provides static code analysis capabilities and auto fixing linting issues. To use vscode-shellcheck in vscode do the following:

#### Install shellcheck on Your Machine

For macOS

```
brew install shellcheck
```
For Ubuntu:

```
apt-get install shellcheck
```
#### Install shellcheck on VSCode

Find the vscode-shellcheck extension in vscode and install it.

## Automatic Code Formatting

### shell-format

shell-format extension does automatic formatting of your bash scripts, docker files and several configuration files. It is dependent on shfmt which can enforce google style guide checks for bash. To use shell-format in vscode do the following:

#### Install shfmt on Your Machine

Requires Go 1.13 or Later

```
GO111MODULE=on go get mvdan.cc/sh/v3/cmd/shfmt
```
#### Install shell-format on VSCode

Find the shell-format extension in vscode and install it.

## Build Validation

To automate this process in Azure DevOps you can add the following snippet to you `azure-pipelines.yaml` file. This will lint any scripts in the `./scripts/` folder.

```
- bash: |
 echo "This checks for formatting and common bash errors. See wiki for error details and ignore options: https://github.com/koalaman/shellcheck/wiki/SC1000"
 export scversion="stable"
 wget -qO- "https://github.com/koalaman/shellcheck/releases/download/${scversion?}/shellcheck-${scversion?}.linux.x86_64.tar.xz" | tar -xJv
 sudo mv "shellcheck-${scversion}/shellcheck" /usr/bin/
 rm -r "shellcheck-${scversion}"
 shellcheck ./scripts/*.sh
 displayName: "Validate Scripts: Shellcheck"
```
Also, your shell scripts can be formatted in your build pipeline by using the `shfmt` tool. To integrate `shfmt` in your build pipeline do the following:

```
- bash: |
 echo "This step does auto formatting of shell scripts"
 shfmt -l -w ./scripts/*.sh
 displayName: "Format Scripts: shfmt"
```
Unit testing using shunit2 can also be added to the build pipeline, using the following block:

```
- bash: |
 echo "This step unit tests shell scripts by using shunit2"
 ./shunit2
 displayName: "Format Scripts: shfmt"
```
## Pre-Commit Hooks

All developers should run shellcheck and shfmt as pre-commit hooks.

### Step 1- Install pre-commit

Run `pip install pre-commit` to install pre-commit.
Alternatively you can run `brew install pre-commit` if you are using homebrew.

### Step 2- Add shellcheck and shfmt

Add .pre-commit-config.yaml file to root of the go project. Run shfmt on pre-commit by adding it to .pre-commit-config.yaml file like below.

```
- repo: git://github.com/pecigonzalo/pre-commit-fmt
 sha: master
 hooks:
 - id: shell-fmt
 args:
 - --indent=4
```
```
- repo: https://github.com/shellcheck-py/shellcheck-py
 rev: v0.7.1.1
 hooks:
 - id: shellcheck
```
### Step 3

Run `$ pre-commit install` to set up the git hook scripts

## Dependencies

Bash scripts are often used to 'glue together' other systems and tools. As such, Bash scripts can often have numerous and/or complicated dependencies. Consider using Docker containers to ensure that scripts are executed in a portable and reproducible environment that is guaranteed to contain all the correct dependencies. To ensure that dockerized scripts are nevertheless easy to execute, consider making the use of Docker transparent to the script's caller by wrapping the script in a 'bootstrap' which checks whether the script is running in Docker and re-executes itself in Docker if it's not the case. This provides the best of both worlds: easy script execution and consistent environments.

```
if [[ "${DOCKER}" != "true" ]]; then
 docker build -t my_script -f my_script.Dockerfile . > /dev/null
 docker run -e DOCKER=true my_script "$@"
 exit $?
fi
# ... implementation of my_script here can assume that all of its dependencies exist since it's always running in Docker ...
```
## Code Review Checklist

In addition to the Code Review Checklist you should also look for these bash specific code review items

- Does this code use Built-in Shell Options like set -o, set -e, set -u for execution control of shell scripts ?
- Is the code modularized? Shell scripts can be modularized like python modules. Portions of bash scripts should be sourced in complex bash projects.
- Are all exceptions handled correctly? Exceptions should be handled correctly using exit codes or trapping signals.
- Does the code pass all linting checks as per shellcheck and unit tests as per shunit2 ?
- Does the code uses relative paths or absolute paths? Relative paths should be avoided as they are prone to environment attacks. If relative path is needed, check that the `PATH`variable is set.
- Does the code take credentials as user input? Are the credentials masked or encrypted in the script? S

# C# Code Reviews

## Style Guide

Developers should follow Microsoft's C# Coding Conventions and, where applicable, Microsoft's Secure Coding Guidelines.

## Code Analysis / Linting

We strongly believe that consistent style increases readability and maintainability of a code base. Hence, we are recommending analyzers / linters to enforce consistency and style rules.

### Project Setup

We recommend using a common setup for your solution that you can refer to in all the projects that are part of the solution. Create a `common.props` file that contains the defaults for all of your projects:

```
<Project>
...
 <ItemGroup>
 <PackageReference Include="Microsoft.CodeAnalysis.NetAnalyzers" Version="5.0.3">
 <PrivateAssets>all</PrivateAssets>
 <IncludeAssets>runtime; build; native; contentfiles; analyzers; buildtransitive</IncludeAssets>
 </PackageReference>
 <PackageReference Include="StyleCop.Analyzers" Version="1.1.118">
 <PrivateAssets>all</PrivateAssets>
 <IncludeAssets>runtime; build; native; contentfiles; analyzers; buildtransitive</IncludeAssets>
 </PackageReference>
 </ItemGroup>
 <PropertyGroup>
 <TreatWarningsAsErrors>true</TreatWarningsAsErrors>
 </PropertyGroup>
 <ItemGroup Condition="Exists('$(MSBuildThisFileDirectory)../.editorconfig')" >
 <AdditionalFiles Include="$(MSBuildThisFileDirectory)../.editorconfig" />
 </ItemGroup>
...
</Project>
```
You can then reference the `common.props` in your other project files to ensure a consistent setup.

```
<Project Sdk="Microsoft.NET.Sdk.Web">
 <Import Project="..\common.props" />
</Project>
```
The .editorconfig allows for configuration and overrides of rules. You can have an .editorconfig file at project level to customize rules for different projects (test projects for example).

Details about the configuration of different rules.

### .NET analyzers

Microsoft's .NET analyzers has code quality rules and .NET API usage rules implemented as analyzers using the .NET Compiler Platform (Roslyn). This is the replacement for Microsoft's legacy FxCop analyzers.

Enable or install first-party .NET analyzers.

If you are currently using the legacy FxCop analyzers, migrate from FxCop analyzers to .NET analyzers.

### StyleCop Analyzer

The StyleCop analyzer is a nuget package (StyleCop.Analyzers) that can be installed in any of your projects. It's mainly around code style rules and makes sure the team is following the same rules without having subjective discussions about braces and spaces. Detailed information can be found here: StyleCop Analyzers for the .NET Compiler Platform.

The minimum rules set teams should adopt is the Managed Recommended Rules rule set.

## Automatic Code Formatting

Use .editorconfig to configure code formatting rules in your project.

## Build Validation

It's important that you enforce your code style and rules in the CI to avoid any team member merging code that does not comply with your standards into your git repo.

If you are using FxCop analyzers and StyleCop analyzer, it's very simple to enable those in the CI. You have to make sure you are setting up the project using nuget and .editorconfig (see Project setup). Once you have this setup, you will have to configure the pipeline to build your code. That's pretty much it. The FxCop analyzers will run and report the result in your build pipeline. If there are rules that are violated, your build will be red.

```
 - task: DotNetCoreCLI@2
 displayName: 'Style Check & Build'
 inputs:
 command: 'build'
 projects: '**/*.csproj'
```
## Enable Roslyn Support in VSCode

The above steps also work in VS Code provided you enable Roslyn support for Omnisharp. The setting is `omnisharp.enableRoslynAnalyzers` and must be set to `true`. After enabling this setting you must "Restart Omnisharp" (this can be done from the Command Palette in VS Code or by restarting VS Code).

## Code Review Checklist

In addition to the Code Review Checklist you should also look for these C# specific code review items

- Does this code make correct use of asynchronous programming constructs, including proper use of `await`and`Task.WhenAll`including CancellationTokens?
- Is the code subject to concurrency issues? Are shared objects properly protected?
- Is dependency injection (DI) used? Is it setup correctly?
- Are middleware included in this project configured correctly?
- Are resources released deterministically using the IDispose pattern? Are all disposable objects properly disposed (using pattern)?
- Is the code creating a lot of short-lived objects. Could we optimize GC pressure?
- Is the code written in a way that causes boxing operations to happen?
- Does the code handle exceptions correctly?
- Is package management being used (NuGet) instead of committing DLLs?
- Does this code use LINQ appropriately? Pulling LINQ into a project to replace a single short loop or in ways that do not perform well are usually not appropriate.
- Does this code properly validate arguments sanity (i.e. CA1062)? Consider leveraging extensions such as Ensure.That
- Does this code include telemetry (metrics, tracing and logging) instrumentation?
- Does this code leverage the options design pattern by using classes to provide strongly typed access to groups of related settings?
- Instead of using raw strings, are constants used in the main class? Or if these strings are used across files/classes, is there a static class for the constants?
- Are magic numbers explained? There should be no number in the code without at least a comment of why this is here. If the number is repetitive, is there a constant/enum or equivalent?
- Is proper exception handling set up? Catching the exception base class (`catch (Exception)`) is generally not the right pattern. Instead, catch the specific exceptions that can happen e.g.,`IOException`.
- Is the use of #pragma fair?
- Are tests arranged correctly with the **Arrange/Act/Assert**pattern and properly documented in this way?
- If there is an asynchronous method, does the name of the method end with the `Async`suffix?
- If a method is asynchronous, is `Task.Delay`used instead of`Thread.Sleep`?`Task.Delay`is not blocking the current thread and creates a task that will complete without blocking the thread, so in a multi-threaded, multi-task environment, this is the one to prefer.
- Is a cancellation token for asynchronous tasks needed rather than bool patterns?
- Is a minimum level of logging in place? Are the logging levels used sensible?
- Are internal vs private vs public classes and methods used the right way?
- Are auto property set and get used the right way? In a model without constructor and for deserialization, it is ok to have all accessible. For other classes usually a private set or internal set is better.
- Is the `using`pattern for streams and other disposable classes used? If not, better to have the`Dispose`method called explicitly.
- Are the classes that maintain collections in memory, thread safe? When used under concurrency, use lock pattern.

# Go Code Reviews

## Style Guide

Developers should follow the Effective Go Style Guide.

## Code Analysis / Linting

### Project Setup

Below is the project setup that you would like to have in your VS Code.

#### VSCode go Extension

Using the Go extension for Visual Studio Code, you get language features like IntelliSense, code navigation, symbol search, bracket matching, snippets, etc. This extension includes rich language support for go in VS Code.

#### go vet

`go vet` is a static analysis tool that checks for common go errors, such as incorrect use of range loop variables or misaligned printf arguments. Go code should be able to build with no `go vet` errors. This will be part of vscode-go extension.

#### golint

Note:The golint library is deprecated and archived.

The linter revive (below) might be a suitable replacement.

golint can be an effective tool for finding many issues, but it errors on the side of false positives. It is best used by developers when working on code, not as part of an automated build process. This is the default linter which is set up as part of the vscode-go extension.

#### revive

Revive is a linter for go, it provides a framework for development of custom rules, and lets you define a strict preset for enhancing your development & code review processes.

## Automatic Code Formatting

### gofmt

`gofmt` is the automated code format style guide for Go. This is part of the vs-code extension, and it is enabled by default to run on save of every file.

## Aggregator

### golangci-lint

golangci-lint is the replacement for the now deprecated `gometalinter`. It is 2-7x faster than `gometalinter` along with a host of other benefits.

golangci-lint is a powerful, customizable aggregator of linters. By default, several are enabled but not all. A full list of linters and their usages can be found here.

It will allow you to configure each linter and choose which ones you would like to enable in your project.

One awesome feature of `golangci-lint` is that is can be easily introduced to an existing large codebase using the `--new-from-rev COMMITID`. With this setting only newly introduced issues are flagged, allowing a team to improve new code without having to fix all historic issues in a large codebase. This provides a great path to improving code-reviews on existing solutions. golangci-lint can also be setup as the default linter in VS Code.

Installation options for golangci-lint are present at golangci-lint.

To use golangci-lint with VS Code, use the below recommended settings:

```
"go.lintTool":"golangci-lint",
 "go.lintFlags": [
 "--fast"
 ]
```
## Pre-Commit Hooks

All developers should run `gofmt` in a pre-commit hook to ensure standard formatting.

### Step 1- Install pre-commit

Run `pip install pre-commit` to install pre-commit.
Alternatively you can run `brew install pre-commit` if you are using homebrew.

### Step 2- Add go-fmt in pre-commit

Add .pre-commit-config.yaml file to root of the go project. Run go-fmt on pre-commit by adding it to .pre-commit-config.yaml file like below.

```
- repo: git://github.com/dnephin/pre-commit-golang
 rev: master
 hooks:
 - id: go-fmt
```
### Step 3

Run `$ pre-commit install` to set up the git hook scripts

## Build Validation

`gofmt` should be run as a part of every build to enforce the common standard.

To automate this process in Azure DevOps you can add the following snippet to your `azure-pipelines.yaml` file. This will format any scripts in the `./scripts/` folder.

```
- script: go fmt
 workingDirectory: $(System.DefaultWorkingDirectory)/scripts
 displayName: "Run code formatting"
```
`govet` should be run as a part of every build to check code linting.

To automate this process in Azure DevOps you can add the following snippet to your `azure-pipelines.yaml` file. This will check linting of any scripts in the `./scripts/` folder.

```
- script: go vet
 workingDirectory: $(System.DefaultWorkingDirectory)/scripts
 displayName: "Run code linting"
```
Alternatively you can use golangci-lint as a step in the pipeline to do multiple enabled validations(including go vet and go fmt) of golangci-lint.

```
- script: golangci-lint run --enable gofmt --fix
 workingDirectory: $(System.DefaultWorkingDirectory)/scripts
 displayName: "Run code linting"
```
## Sample Build Validation Pipeline in Azure DevOps

```
trigger: master
pool:
 vmImage: 'ubuntu-latest'
steps:
- task: GoTool@0
 inputs:
 version: '1.13.5'
- task: Go@0
 inputs:
 command: 'get'
 arguments: '-d'
 workingDirectory: '$(System.DefaultWorkingDirectory)/scripts'
- script: go fmt
 workingDirectory: $(System.DefaultWorkingDirectory)/scripts
 displayName: "Run code formatting"
- script: go vet
 workingDirectory: $(System.DefaultWorkingDirectory)/scripts
 displayName: 'Run go vet'
- task: Go@0
 inputs:
 command: 'build'
 workingDirectory: '$(System.DefaultWorkingDirectory)'
- task: CopyFiles@2
 inputs:
 TargetFolder: '$(Build.ArtifactStagingDirectory)'
- task: PublishBuildArtifacts@1
 inputs:
 artifactName: drop
```
## Code Review Checklist

The Go language team maintains a list of common Code Review Comments for go that form the basis for a solid checklist for a team working in Go that should be followed in addition to the ISE Code Review Checklist

- Does this code handle errors correctly? This includes not throwing away errors with `_`assignments and returning errors, instead of in-band error values?
- Does this code follow Go standards for method receiver types?
- Does this code pass values when it should?
- Are interfaces in this code defined in the correct packages?
- Do go-routines in this code have clear lifetimes?
- Is parallelism in this code handled via go-routines and channels with synchronous methods?
- Does this code have meaningful Doc Comments?
- Does this code have meaningful Package Comments?
- Does this code use Contexts correctly?
- Do unit tests fail with meaningful messages?

# Java Code Reviews

## Java Style Guide

Developers should follow the Google Java Style Guide.

## Code Analysis / Linting

We strongly believe that consistent style increases readability and maintainability of a code base. Hence, we are recommending analyzers to enforce consistency and style rules.

We make use of Checkstyle using the same configuration used in the Azure Java SDK.

FindBugs and PMD are also commonly used.

## Automatic Code Formatting

Eclipse, and other Java IDEs, support automatic code formatting. If using Maven, some developers also make use of the formatter-maven-plugin.

## Build Validation

It's important to enforce your code style and rules in the CI to avoid any team members merging code that does not comply with standards into your git repo. If building using Azure DevOps, Azure DevOps support Maven and Gradle build tasks using PMD, Checkstyle, and FindBugs code analysis tools as part of every build.

Here is an example yaml for a Maven build task with all three analysis tools enabled:

```
 - task: Maven@3
 displayName: 'Maven pom.xml'
 inputs:
 mavenPomFile: '$(Parameters.mavenPOMFile)'
 checkStyleRunAnalysis: true
 pmdRunAnalysis: true
 findBugsRunAnalysis: true
```
Here is an example yaml for a Gradle build task with all three analysis tools enabled:

```
 - task: Gradle@2
 displayName: 'gradlew build'
 inputs:
 checkStyleRunAnalysis: true
 findBugsRunAnalysis: true
 pmdRunAnalysis: true
```
## Code Review Checklist

In addition to the Code Review Checklist you should also look for these Java specific code review items

- Does the project use Lambda to make code cleaner?
- Is dependency injection (DI) used? Is it setup correctly?
- If the code uses Spring Boot, are you using @Inject instead of @Autowire?
- Does the code handle exceptions correctly?
- Is the Azul Zulu OpenJDK being used?
- Is a build automation and package management tool (Gradle or Maven) being used?

# JavaScript/TypeScript Code Reviews

## Style Guide

Developers should use prettier to do code formatting for JavaScript.

Using an automated code formatting tool like Prettier enforces a well accepted style guide that was collaboratively built by a wide range of companies including Microsoft, Facebook, and AirBnB.

For higher level style guidance not covered by prettier, we follow the AirBnB Style Guide.

## Code Analysis / Linting

### eslint

Per guidance outlined in Palantir's 2019 TSLint road map, TypeScript code should be linted with ESLint. See the typescript-eslint documentation for more information around linting TypeScript code with ESLint.

To install and configure linting with ESLint, install the following packages as dev-dependencies:

```
npm install -D eslint @typescript-eslint/parser @typescript-eslint/eslint-plugin
```
Add a `.eslintrc.js` to the root of your project:

```
module.exports = {
 root: true,
 parser: '@typescript-eslint/parser',
 plugins: [
 '@typescript-eslint',
 ],
 extends: [
 'eslint:recommended',
 'plugin:@typescript-eslint/eslint-recommended',
 'plugin:@typescript-eslint/recommended',
 ],
};
```
Add the following to the `scripts` of your `package.json`:

```
"scripts": {
 "lint": "eslint . --ext .js,.jsx,.ts,.tsx --ignore-path .gitignore"
}
```
This will lint all `.js`, `.jsx`, `.ts`, `.tsx` files in your project and omit any files or
directories specified in your `.gitignore`.

You can run linting with:

```
npm run lint
```
## Setting up Prettier

Prettier is an opinionated code formatter.

Install with `npm` as a dev-dependency:

```
npm install -D prettier eslint-config-prettier eslint-plugin-prettier
```
Add `prettier` to your `.eslintrc.js`:

```
module.exports = {
 root: true,
 parser: '@typescript-eslint/parser',
 plugins: [
 '@typescript-eslint',
 ],
 extends: [
 'eslint:recommended',
 'plugin:@typescript-eslint/eslint-recommended',
 'plugin:@typescript-eslint/recommended',
 'prettier/@typescript-eslint',
 'plugin:prettier/recommended',
 ],
};
```
This will apply the `prettier` rule set when linting with ESLint.

## Auto Formatting with VSCode

VS Code can be configured to automatically perform `eslint --fix` on save.

Create a `.vscode` folder in the root of your project and add the following to your
`.vscode/settings.json`:

```
{
 "editor.codeActionsOnSave": {
 "source.fixAll.eslint": true
 },
}
```
By default, we use the following overrides should be added to the VS Code configuration to standardize on single quotes, a four space drop, and to do ESLinting:

```
{
 "prettier.singleQuote": true,
 "prettier.eslintIntegration": true,
 "prettier.tabWidth": 4
}
```
## Setting Up Testing

Playwright is highly recommended to be set up within a project. its an open source testing suite created by Microsoft.

To install it use this command:

```
npm install playwright
```
Since playwright shows the tests in the browser you have to choose which browser you want it to run if unless using chrome, which is the default. You can do this by

## Build Validation

To automate this process in Azure Devops you can add the following snippet to your pipeline definition yaml file. This will lint any scripts in the `./scripts/` folder.

```
- task: Npm@1
 displayName: 'Lint'
 inputs:
 command: 'custom'
 customCommand: 'run lint'
 workingDir: './scripts/'
```
## Pre-Commit Hooks

All developers should run `eslint` in a pre-commit hook to ensure standard formatting. We highly recommend using an editor integration like vscode-eslint to provide realtime feedback.

- Under `.git/hooks`rename`pre-commit.sample`to`pre-commit`
- Remove the existing sample code in that file
- There are many examples of scripts for this on gist, like pre-commit-eslint
- Modify accordingly to include TypeScript files (include ts extension and make sure typescript-eslint is set up)
- Make the file executable: `chmod +x .git/hooks/pre-commit`

As an alternative husky can be considered to simplify pre-commit hooks.

## Code Review Checklist

In addition to the Code Review Checklist you should also look for these JavaScript and TypeScript specific code review items.

### Javascript / Typescript Checklist

- Does the code stick to our formatting and code standards? Does running prettier and ESLint over the code should yield no warnings or errors respectively?
- Does the change re-implement code that would be better served by pulling in a well known module from the ecosystem?
- Is `"use strict";`used to reduce errors with undeclared variables?
- Are unit tests used where possible, also for APIs?
- Are tests arranged correctly with the **Arrange/Act/Assert**pattern and properly documented in this way?
- Are best practices for error handling followed, as well as `try catch finally`statements?
- Are the `doWork().then(doSomething).then(checkSomething)`properly followed for async calls, including`expect`,`done`?
- Instead of using raw strings, are constants used in the main class? Or if these strings are used across files/classes, is there a static class for the constants?
- Are magic numbers explained? There should be no number in the code without at least a comment of why it is there. If the number is repetitive, is there a constant/enum or equivalent?
- If there is an asynchronous method, does the name of the method end with the `Async`suffix?
- Is a minimum level of logging in place? Are the logging levels used sensible?
- Is document fragment manipulation limited to when you need to manipulate multiple sub elements?
- Does TypeScript code compile without raising linting errors?
- Instead of using raw strings, are constants used in the main class? Or if these strings are used across files/classes, is there a static class for the constants?
- Are magic numbers explained? There should be no number in the code without at least a comment of why it is there. If the number is repetitive, is there a constant/enum or equivalent?
- Is there a proper `/* */`in the various classes and methods?
- Are heavy operations implemented in the backend, leaving the controller as thin as possible?
- Is event handling on the html efficiently done?

# Markdown Code Reviews

## Style Guide

Developers should treat documentation like other source code and follow the same rules and checklists when reviewing documentation as code.

Documentation should both use good Markdown syntax to ensure it's properly parsed, and follow good writing style guidelines to ensure the document is easy to read and understand.

## Markdown

Markdown is a lightweight markup language that you can use to add formatting elements to plaintext text documents. Created by John Gruber in 2004, Markdown is now one of the world’s most popular markup languages.

Using Markdown is different from using a WYSIWYG editor. In an application like Microsoft Word, you click buttons to format words and phrases, and the changes are visible immediately. Markdown isn’t like that. When you create a Markdown-formatted file, you add Markdown syntax to the text to indicate which words and phrases should look different.

You can find more information and full documentation here.

## Linters

Markdown has specific way of being formatted. It is important to respect this formatting, otherwise some interpreters which are strict won't properly display the document. Linters are often used to help developers properly create documents by both verifying proper Markdown syntax, grammar and proper English language.

A good setup includes a markdown linter used during editing and PR build verification, and a grammar linter used while editing the document. The following are a list of linters that could be used in this setup.

### markdownlint

`markdownlint` is a linter for markdown that verifies Markdown syntax, and also enforces rules that make the text more readable. Markdownlint-cli is an easy-to-use CLI based on Markdownlint.

It's available as a ruby gem, an npm package, a Node.js CLI and a VS Code extension. The VS Code extension Prettier also catches all markdownlint errors.

Installing the Node.js CLI

```
npm install -g markdownlint-cli
```
Running markdownlint on a Node.js project

```
markdownlint **/*.md --ignore node_modules
```
Fixing errors automatically

```
markdownlint **/*.md --ignore node_modules --fix
```
A comprehensive list of markdownlint rules is available here.

### write-good

`write-good` is a linter for English text that helps writing better documentation.

```
npm install -g write-good
```
Run write-good

```
write-good *.md
```
Run write-good without installing it

```
npx write-good *.md
```
Write Good is also available as an extension for VS Code

## VSCode Extensions

### Write Good Linter

The `Write Good Linter Extension` integrates with VS Code to give grammar and language advice while editing the document.

### markdownlint Extension

The `markdownlint extension` examines the Markdown documents, showing warnings for rule violations while editing.

## Build Validation

### Linting

To automate linting with `markdownlint` for PR validation in GitHub actions,
you can either use linters aggregator as we do with MegaLinter in this repository or use the following YAML.

```
name: Markdownlint
on:
 push:
 paths:
 - "**/*.md"
 pull_request:
 paths:
 - "**/*.md"
jobs:
 lint:
 runs-on: ubuntu-latest
 steps:
 - uses: actions/checkout@v2
 - name: Use Node.js
 uses: actions/setup-node@v1
 with:
 node-version: 12.x
 - name: Run Markdownlint
 run: |
 npm i -g markdownlint-cli
 markdownlint "**/*.md" --ignore node_modules
```
### Checking Links

To automate link check in your markdown files add `lycheeverse/lychee-action` action to your validation pipeline:

```
 markdown-link-check:
 runs-on: ubuntu-latest
 steps:
 - uses: actions/checkout@v4
 - name: Link Checker
 id: lychee
 uses: lycheeverse/lychee-action@v2
```
More information about this action options can be found at `lychee-action` home page.

## Code Review Checklist

In addition to the Code Review Checklist you should also look for these documentation specific code review items

- Is the document easy to read and understand and does it follow good writing guidelines?
- Is there a single source of truth or is content repeated in more than one document?
- Is the documentation up to date with the code?
- Is the documentation technically, and ethically correct?

## Writing Style Guidelines

The following are some examples of writing style guidelines.

Agree in your team which guidelines you should apply to your project documentation. Save your guidelines together with your documentation, so they are easy to refer back to.

### Wording

- Use inclusive language, and avoid jargon and uncommon words. The docs should be easy to understand
- Be clear and concise, stick to the goal of the document
- Use active voice
- Spell check and grammar check the text
- Always follow chronological order
- Visit Plain English for tips on how to write documentation that is easy to understand.

### Document Organization

- Organize documents by topic rather than type, this makes it easier to find the documentation
- Each folder should have a top-level README.md and any other documents within that folder should link directly or indirectly from that README.md
- Document names with more than one word should use underscores instead of spaces, for example `machine_learning_pipeline_design.md`. The same applies to images

### Headings

- Start with a H1 (single # in markdown) and respect the order H1 > H2 > H3 etc
- Follow each heading with text before proceeding with the next heading
- Avoid putting numbers in headings. Numbers shift, and can create outdated titles
- Avoid using symbols and special characters in headers, this causes problems with anchor links
- Avoid links in headers

### Resources

- Avoid duplication of content, instead link to the `single source of truth`
- Link but don't summarize. Summarizing content on another page leads to the content living in two places
- Use meaningful anchor texts, e.g. instead of writing `Follow the instructions [here](../recipes/markdown.md)`write`Follow the [Markdown guidelines](../recipes/markdown.md)`
- Make sure links to Microsoft docs do not contain the language marker `/en-us/`or`/fr-fr/`, as this is automatically determined by the site itself.

### Lists

- List items should start with capital letters if possible
- Use ordered lists when the items describe a sequence to follow, otherwise use unordered lists
- For ordered lists, prefix each item with `1.`When rendered, the list items will appear with sequential numbering. This avoids number-gaps in list
- Do not add commas `,`or semicolons`;`to the end of list items, and avoid periods`.`unless the list item represents a complete sentence

### Images

- Place images in a separate directory named `img`
- Name images appropriately, avoiding generic names like `screenshot.png`
- Avoid adding large images or videos to source control, link to an external location instead

### Emphasis and Special Sections

- Use **bold**or*italic*to emphasizeFor sections that everyone reading this document needs to be aware of, use blocks
-
Use `backticks`for code, a single backtick for inline code like`pip install flake8`and 3 backticks for code blocks followed by the language for syntax highlighting`def add(num1: int, num2: int): return num1 + num2`

- Use check boxes for task lists- Item 1
- Item 2
- Item 3

- Add a References section to the end of the document with links to external references
-
Prefer tables to lists for comparisons and reports to make research and results more readable Option Pros Cons Option 1 Some pros Some cons Option 2 Some pros Some cons

### General

- Always use Markdown syntax, don't mix with HTML
- Make sure the extension of the files is `.md`- if the extension is missing, a linter might ignore the files

# Python Code Reviews

## Style Guide

Developers should follow the PEP8 style guide with type hints. The use of type hints throughout paired with linting and type hint checking avoids common errors that are tricky to debug.

Projects should check Python code with automated tools.

Linting should be added to build validation, and both linting and code formatting can be added to your pre-commit hooks and VS Code.

## Code Analysis / Linting

The 2 most popular python linters are Pylint and Flake8. Both check adherence to `PEP8` but vary a bit in what other rules they check. In general `Pylint` tends to be a bit more stringent and give more false positives but both are good options for linting python code.

Both `Pylint` and `Flake8` can be configured in VS Code using the VS Code `python extension`.

### Flake8

Flake8 is a simple and fast wrapper around `Pyflakes` (for detecting coding errors) and `pycodestyle` (for pep8).

Install `Flake8`

```
pip install flake8
```
Add an extension for the `pydocstyle` (for doc strings) tool to flake8.

```
pip install flake8-docstrings
```
Add an extension for `pep8-naming` (for naming conventions in pep8) tool to flake8.

```
pip install pep8-naming
```
Run `Flake8`

```
flake8 . # lint the whole project
```
### Pylint

Install `Pylint`

```
pip install pylint
```
Run `Pylint`

```
pylint src # lint the source directory
```
## Automatic Code Formatting

### Black

`Black` is an unapologetic code formatting tool. It removes all need from `pycodestyle` nagging about formatting, so the team can focus on content vs style. It's not possible to configure black for your own style needs.

```
pip install black
```
Format python code

```
black [file/folder]
```
### autopep8

`Autopep8` is more lenient and allows more configuration if you want less stringent formatting.

```
pip install autopep8
```
Format python code

```
autopep8 [file/folder] --in-place
```
### yapf

yapf Yet Another Python Formatter is a python formatter from Google based on ideas from gofmt. This is also more configurable, and a good option for automatic code formatting.

```
pip install yapf
```
Format python code

```
yapf [file/folder] --in-place
```
### Bandit

Bandit is a tool designed by the Python Code Quality Authority (PyCQA) to perform static analysis of Python code, specifically targeting security issues. It scans for common security issues in Python codebase.

- **Installation**: Add Bandit to your development environment with:- `pip install bandit`

## VSCode Extensions

### Python

The `Python language extension` is the base extension you should have installed for python development with VS Code. It enables intellisense, debugging, linting (with the above linters), testing with pytest or unittest, and code formatting with the formatters mentioned above.

### Pyright

The `Pyright extension` augments VS Code with static type checking when you use type hints

```
def add(first_value: int, second_value: int) -> int:
 return first_value + second_value
```
## Build Validation

To automate linting with `flake8` and testing with `pytest` in Azure Devops you can add the following snippet to you `azure-pipelines.yaml` file.

```
trigger:
 branches:
 include:
 - develop
 - master
 paths:
 include:
 - src/*
pool:
 vmImage: 'ubuntu-latest'
jobs:
- job: LintAndTest
 displayName: Lint and Test
 steps:
 - checkout: self
 lfs: true
 - task: UsePythonVersion@0
 displayName: 'Set Python version to 3.6'
 inputs:
 versionSpec: '3.6'
 - script: pip3 install --user -r requirements.txt
 displayName: 'Install dependencies'
 - script: |
 # Install Flake8
 pip3 install --user flake8
 # Install PyTest
 pip3 install --user pytest
 displayName: 'Install Flake8 and PyTest'
 - script: |
 python3 -m flake8
 displayName: 'Run Flake8 linter'
 - script: |
 # Run PyTest tester
 python3 -m pytest --junitxml=./test-results.xml
 displayName: 'Run PyTest Tester'
 - task: PublishTestResults@2
 displayName: 'Publish PyTest results'
 condition: succeededOrFailed()
 inputs:
 testResultsFiles: '**/test-*.xml'
 testRunTitle: 'Publish test results for Python $(python.version)'
```
To perform a PR validation on GitHub you can use a similar YAML configuration with GitHub Actions

## Pre-Commit Hooks

Pre-commit hooks allow you to format and lint code locally before submitting the pull request.

Adding pre-commit hooks for your python repository is easy using the pre-commit package

-
Install pre-commit and add to the requirements.txt `pip install pre-commit`
-
Add a `.pre-commit-config.yaml`file in the root of the repository, with the desired pre-commit actions`repos: - repo: https://github.com/ambv/black rev: stable hooks: - id: black language_version: python3.6 - repo: https://github.com/pre-commit/pre-commit-hooks rev: v1.2.3 hooks: - id: flake8`
-
Each individual developer that wants to set up pre-commit hooks can then run `pre-commit install`

At the next attempted commit any lint failures will block the commit.

Note: Installing pre-commit hooks is voluntary and done by each developer individually. Thus, it's not a replacement for build validation on the server

## Code Review Checklist

In addition to the Code Review Checklist you should also look for these python specific code review items

- Are all new packages used included in requirements.txt
- Does the code pass all lint checks?
- Do functions use type hints, and are there any type hint errors?
- Is the code readable and using pythonic constructs wherever possible.

# Terraform Code Reviews

## Style Guide

Developers should follow the terraform style guide.

Projects should check Terraform scripts with automated tools.

## Code Analysis / Linting

### TFLint

`TFLint` is a Terraform linter focused on possible errors, best practices, etc. Once TFLint installed in the environment, it can be invoked using the VS Code `terraform extension`.

## VSCode Extensions

The following VS Code extensions are widely used.

`Terraform extension`

This extension provides syntax highlighting, linting, formatting and validation capabilities.

`Azure Terraform extension`

This extension provides Terraform command support, resource graph visualization and CloudShell integration inside VS Code.

## Build Validation

Ensure you enforce the style guides during build. The following example script can be used to install terraform, and a linter that then checks for formatting and common errors.

```
#! /bin/bash
set -e
SCRIPT_DIR=$(dirname "$BASH_SOURCE")
cd "$SCRIPT_DIR"
TF_VERSION=0.12.4
TF_LINT_VERSION=0.9.1
echo -e "\n\n>>> Installing Terraform 0.12"
# Install terraform tooling for linting terraform
wget -q https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip -O /tmp/terraform.zip
sudo unzip -q -o -d /usr/local/bin/ /tmp/terraform.zip
echo ""
echo -e "\n\n>>> Install tflint (3rd party)"
wget -q https://github.com/wata727/tflint/releases/download/v${TF_LINT_VERSION}/tflint_linux_amd64.zip -O /tmp/tflint.zip
sudo unzip -q -o -d /usr/local/bin/ /tmp/tflint.zip
echo -e "\n\n>>> Terraform version"
terraform -version
echo -e "\n\n>>> Terraform Format (if this fails use 'terraform fmt -recursive' command to resolve"
terraform fmt -recursive -diff -check
echo -e "\n\n>>> tflint"
tflint
echo -e "\n\n>>> Terraform init"
terraform init
echo -e "\n\n>>> Terraform validate"
terraform validate
```
## Code Review Checklist

In addition to the Code Review Checklist you should also look for these Terraform specific code review items

### Providers

- Are all providers used in the terraform scripts versioned to prevent breaking changes in the future?

### Repository Organization

- The code split into reusable modules?
- Modules are split into separate `.tf`files where appropriate?
- The repository contains a `README.md`describing the architecture provisioned?
- If Terraform code is mixed with application source code, the Terraform code isolated into a dedicated folder?

### Terraform State

- The Terraform project configured using Azure Storage as remote state backend?
- The remote state backend storage account key stored a secure location (e.g. Azure Key Vault)?
- The project is configured to use state files based on the environment, and the deployment pipeline is configured to supply the state file name dynamically?

### Variables

- If the infrastructure will be different depending on the environment (e.g. Dev, UAT, Production), the environment specific parameters are supplied via a `.tfvars`file?
- All variables have `type`information. E.g. a`list(string)`or`string`.
- All variables have a `description`stating the purpose of the variable and its usage.
- `default`values are not supplied for variables which must be supplied by a user.

### Testing

- Unit and integration tests covering the Terraform code exist (e.g. `Terratest`,`terratest-abstraction`)?

### Naming and Code Structure

- Resource definitions and data sources are used correctly in the Terraform scripts?- **resource:**Indicates to Terraform that the current configuration is in charge of managing the life cycle of the object
- **data:**Indicates to Terraform that you only want to get a reference to the existing object, but don’t want to manage it as part of this configuration

- The resource names start with their containing provider's name followed by an underscore? e.g. resource from the provider `postgresql`might be named as`postgresql_database`?
- The `try`function is only used with simple attribute references and type conversion functions? Overuse of the`try`function to suppress errors will lead to a configuration that is hard to understand and maintain.
- Explicit type conversion functions used to normalize types are only returned in module outputs? Explicit type conversions are rarely necessary in Terraform because it will convert types automatically where required.
- The `Sensitive`property on schema set to`true`for the fields that contains sensitive information? This will prevent the field's values from showing up in CLI output.

### General Recommendations

- Try avoiding nesting sub configuration within resources. Create a separate resource section for resources even though they can be declared as sub-element of a resource. For example, declaring subnets within virtual network vs declaring subnets as a separate resources compared to virtual network on Azure.
- Never hard-code any value in configuration. Declare them in `locals`section if a variable is needed multiple times as a static value and are internal to the configuration.
- The `name`s of the resources created on Azure should not be hard-coded or static. These names should be dynamic and user-provided using`variable`block. This is helpful especially in unit testing when multiple tests are running in parallel trying to create resources on Azure but need different names (few resources in Azure need to be named uniquely e.g. storage accounts).
- It is a good practice to `output`the ID of resources created on Azure from configuration. This is especially helpful when adding dynamic blocks for sub-elements/child elements to the parent resource.
- Use the `required_providers`block for establishing the dependency for providers along with pre-determined version.
- Use the `terraform`block to declare the provider dependency with exact version and also the terraform CLI version needed for the configuration.
- Validate the variable values supplied based on usage and type of variable. The validation can be done to variables by adding `validation`block.
- Validate that the component SKUs are the right ones, e.g. standard vs premium.
