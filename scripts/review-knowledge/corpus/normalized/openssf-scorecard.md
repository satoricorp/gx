#
 Build better security habits,

 one test at a time

 Part of the Open Source Security Foundation

## Run the checks

OpenSSF Scorecard can be used in a couple of different ways:

- Run automatically on code you own **using the GitHub Action**
- Run manually on your (or somebody else’s) project **via the Command Line**

### Using the GitHub Action

### Install time: <10 mins

Use the action to automatically scan any code updates for security vulnerabilities. Any time someone commits a change, the action will automatically check the repo and alert you (and other maintainers) if there are problems.

## See it in action

### Installation instructions

- You need to own the repository you are installing the action to, or have admin rights to it.
- Authenticate your access to the repository with a Personal Access Token
- Add Scorecard to your codescanning suite inside GitHub using the link below:

### Using the CLI

### Install time: <10mins

You can use Scorecard on the Command Line. This enables you to:

- Check someone else’s repository
- Select which checks you want to run
- Control how detailed your results are

## See it in action

### Install and run

- Create a GitHub personal access token with 'public_repo' scope. Store the token somewhere safe.
- Choose a language-specific quick start below, or refer to our detailed instructions

Scorecard also has standalone binaries and other platforms troubleshooting and custom configuration available. Learn more here:

## Learn more

We rely on Security Scorecards

[i.e., OpenSSF Scorecard]to ensure we follow secure development best practices.

### The problem

By some estimates* 84% of all codebases have at least one vulnerability, with an average of 158 per codebase. The majority have been in the code for more than 2 years and have documented solutions available.

Even in large tech companies, the tedious process of reviewing code for vulnerabilities falls down the priority list, and there is little insight into known vulnerabilities and solutions that companies can draw on.

That’s where Security Scorecards *[i.e., OpenSSF Scorecard]* is helping. Its focus is to understand the security posture of a project and assess the risks that dependencies introduce.

*Open Source Security and Risk Analysis Report (Synopsys, 2021)

### What is OpenSSF Scorecard?

##### Scorecard assesses open source projects for security risks through a series of automated checks.

It was created by OSS developers to help improve the health of critical projects that the community depends on.

You can use it to proactively assess and make informed decisions about accepting security risks within your codebase. You can also use the tool to evaluate other projects and dependencies, and work with maintainers to improve codebases you might want to integrate.

Scorecard helps you enforce best practices that can guard against:

#### Malicious maintainers

#### Build system compromises

#### Source code compromises

#### Malicious packages

It took less than 5 minutes to install. It quickly analysed the repo and identified easy ways to make the project more secure.

### How it works

Scorecard checks for vulnerabilities affecting different parts of the software supply chain including **source code**, **build**, **dependencies**, **testing**, and project **maintenance**.

Each automated check returns a **score out of 10** and a **risk level**. The risk level **adds a weighting** to the score, and this weighting is compiled into a single, **aggregate score**. This score helps give a sense of the overall security posture of a project.

Alongside the scores, the tool provides remediation prompts to help you **fix problems** and strengthen your development practices.

### The checks

##### The checks collect together security best practises and industry standards

The riskiness of each vulnerability is based on how easy it is to exploit. For example if something can be exploited via a pull request, we consider that a high risk. There are currently 18 checks made across 3 themes: holistic security practises, source code risk assessment and build process risk assessment.

You can learn more about the scoring criteria, risks, and remediation suggestions for each check in the detailed documentation.

#### Holistic security practises

| Code vulnerabilities | Description | Risk |
|---|---|---|
| Vulnerabilities | Does the project have unfixed vulnerabilities? Uses the OSV service. | High |

| Maintenance | Description | Risk |
|---|---|---|
| Dependency Update Tool | Does the project use tools to help update its dependencies e.g. Dependabot, RenovateBot? | High |
| Maintained | Is the project maintained? | High |
| Security Policy | Does the project contain a security policy? | Medium |
| Licence | Does the project declare a licence? | Low |
| CII Best Practices | Does the project have a CII Best Practices Badge? | Low |

| Continuous testing | Description | Risk |
|---|---|---|
| CI Tests | Does the project run tests in CI, e.g. GitHub Actions, Prow? | Low |
| Fuzzing | Does the project use fuzzing tools, e.g. OSS-Fuzz? | Medium |
| SAST | Does the project use static code analysis tools, e.g. CodeQL, LGTM, SonarCloud? | Medium |

#### Source risk assessment

| Name | Description | Risk |
|---|---|---|
| Binary Artifacts | Is the project free of checked-in binaries? | High |
| Branch Protection | Does the project use Branch Protection? | High |
| Dangerous Workflow | Does the project avoid dangerous coding patterns in GitHub Actions? | Critical |
| Code Review | Does the project require code review before code is merged? | High |
| Contributors | Does the project have contributors from multiple organizations? | Low |

#### Build risk assessment

| Name | Description | Risk |
|---|---|---|
| Pinned Dependencies | Does the project declare and pin dependencies? | Medium |
| Token Permissions | Does the project declare GitHub workflow tokens as read only? | High |
| Packaging | Does the project build and publish official packages from CI/CD, e.g. GitHub Publishing? | Medium |
| Signed Releases | Does the project cryptographically sign releases? | High |

Machine checkable properties are an essential part of a sound security process. That’s why we have incorporated Security Scorecards

[i.e., OpenSSF Scorecard]into our dependency acceptance criteria.

### Use cases

##### OpenSSF Scorecard reduces the effort required to continually evaluate changing packages when maintaining a project’s supply chain

#### For individual maintainers

Scorecard is helpful as a pre-launch security checker for a new OSS project or to help to plan improvements to an existing one. If a project is well maintained, it’s more likely to be used by others instead of an alternative. It can also be used to check a new dependency being added to a project, so a maintainer can make an informed decision about the risk of doing so.

#### For an organisation

Scorecard can be included in the continuous integration/continuous deployment processes using the GitHub action and run by default on pull requests.

#### For consumers

Scorecard helps to make informed decisions about security risks and vulnerabilities. Using the public data, it is also possible to evaluate the security posture of over 1 million of the most used OSS projects.

### About the project name

This project was initially called "Security Scorecards" but that form wasn't used consistently. In particular, the repo was named "scorecard" and so was the program. Over time people started referring to either form (singular and plural), with or without "Security", and the inconsitency became prevalent. To end this situation the decision was made to consolidate over the use of the singular form in keeping with the repo and program name, drop the "Security" part and use "OpenSSF" instead to ensure uniqueness. One should therefore refer to this project as "OpenSSF Scorecard" or "Scorecard" for short.

### Part of the OSS community

& many others

OpenSSF Scorecard is being developed and facilitated by contributors from across the OSS ecosystem.

We're part of the Open Source Security Foundation (OpenSSF), a cross-industry collaboration that brings together OSS security initiatives under one foundation and seeks to improve the security of OSS by building a broader community, targeted initiatives, and best practises.

OpenSSF launched Scorecard in November 2020 with the intention of auto-generating a “security score” for open source projects to help users as they decide the trust, risk, and security posture for their use case.

### Get involved

Scorecard is part of the OpenSSF Best Practices Working Group.

If you want to get involved in the OpenSSF Scorecard community or have ideas you'd like to chat about, we'd love to connect.
