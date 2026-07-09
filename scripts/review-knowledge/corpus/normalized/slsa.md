“This is our chance to work with the industry, set a standard which we can all agree to, and work together to raise the collective bar.”

Trishank Karthik Kuppusamy

Staff Security Engineer at Datadog

**Supply-chain Levels for Software Artifacts, or SLSA ("salsa").**

It’s a security framework, a checklist of standards and controls to prevent tampering, improve integrity, and secure packages and infrastructure. It’s how you get from "safe enough" to being as resilient as possible, at any link in the chain.

Any software can introduce vulnerabilities into a supply chain. As a system gets more complex, it’s critical to already have checks and best practices in place to guarantee artifact integrity, that the source code you’re relying on is the code you’re actually using. Without solid foundations and a plan for the system as it grows, it’s difficult to focus your efforts against tomorrow’s next hack, breach or compromise.

More about supply chain attacksSLSA levels are like a common language to talk about how secure software, supply chains and their component parts really are. From source to platform, the levels blend together industry-recognized best practices to create four compliance levels of increasing assurance. These look at the builds, sources and dependencies in open source or commercial software. Starting with easy, basic steps at the lower levels to build up and protect against advanced threats later, bringing SLSA into your work means prioritized, practical measures to prevent unauthorized modifications to software, and a plan to harden that security over time.

Read the level specificationsSLSA is for everyone involved in producing, consuming, and providing infrastructure for software such as build platforms and package ecosystems. SLSA can help create more trust across the entire supply chain. It can be used by producers for protection against tampering and insider threats, by consumers to verify the software they rely on is secure, and by infrastructure providers as a guideline for hardening build platforms and processes.

An industry collaboration

SLSA is led by an initial cross-organization, vendor-neutral steering group committed to improving the security ecosystem for everyone.

Part of the Open Source Security Foundation

Our ethos

Today’s projects, products and services are increasingly complex and open to attack. As that trend continues, we need to scale up our effort to provide more secure, accessible ways to protect the development, distribution and consumption of the software we use, and all the impacted communities behind it.

Get started

Since the release of SLSA v1.0 in 2023, the SLSA community has been hard at work to improve the specification and expand its breadth and depth with updates and new tracks.

Learn how you can get involved!

## Build Environment track

The goal of a Build Environment track is to enable the detection of tampering with core components of the compute environment executing builds.

The current draft version of the Build Environment track includes the following requirements:

- Generation and verification of SLSA Build Provenance for build images.
- Validation of initial build environment system state against known good reference values.
- Deployment of the hosted build platform on a compute system that supports system state measurement and attestation capabilities at the hardware level.

These requirements are **subject to significant change** while this track
is in draft.

## Dependency track

The Dependency track defines requirements aimed at mitigating risks introduced to a project through its dependencies.

There is an ongoing effort to develop this track, building upon the foundation laid by S2C2F. S2C2F provides a guide for the safe consumption of open source dependencies. As part of this ongoing effort, community members from both SLSA and S2C2F collaborated to translate the practices outlined in S2C2F into requirements that align with the structure of existing SLSA specifications.

A first draft of the Dependency track has been published to the SLSA Working Draft site.

Join the discussions and ongoing efforts on the SLSA Dependency Track Slack channel.

# SLSA specification

SLSA is a specification for describing and incrementally improving supply chain security, established by industry consensus. It is organized into a series of levels that describe increasing security guarantees.

This is **Version 1.2** of the SLSA specification. It defines several SLSA
levels and tracks, as well as recommended attestation formats, including
provenance.

## Understanding SLSA

These pages provide an overview of SLSA, how it helps protect against common supply chain attacks, and common use cases. If you’re new to SLSA or supply chain security, start here.

| Page | Description |
|---|---|
| What’s new | The changes brought by this revision of SLSA. |
| About SLSA | An introductory guide to SLSA |
| Supply chain threats | An introduction to supply chain threats |
| Use cases | Use cases |
| Guiding principles | Use cases |
| FAQ | Questions and more information |
| Future directions | Additions and changes being considered for future SLSA versions |
| Tracks | Provides an overview of each track and links to more specific information. |

## Build Track

These pages describe the build track’s security levels and requirements. If you want to achieve a particular level of the SLSA build track these are the requirements you’ll need to meet.

| Page | Description |
|---|---|
| Basics | The SLSA build track is organized into a series of levels that provide increasing supply chain security guarantees. This gives you confidence that software hasn’t been tampered with and can be securely traced back to its source. This page is a descriptive overview of the SLSA build track levels, describing their intent. |
| Terminology | Terminology and model used by SLSA |
| Producing artifacts | Detailed technical requirements for producing software artifacts, intended for platform implementers |
| Distributing provenance | Detailed technical requirements for distributing provenance, intended for platform implementers and software distributors |
| Verifying artifacts | Guidance for verifying software artifacts and their SLSA provenance, intended for platform implementers and software consumers |
| Assessing build platforms | Guidelines for securing SLSA Build L3+ builders, intended for platform implementers |

## Source Track

These pages describe the source track’s security levels and requirements. If you want to achieve a particular level of the SLSA source track these are the requirements you’ll need to meet.

| Page | Description |
|---|---|
| Producing source | Overview of the Source track |
| Verifying source | Guidelines for verifying source provenance |
| Assessing source control systems | Guidelines for assessing source control system security. |
| Example controls | This page provides examples of additional controls that organizations may want to implement as they adopt the SLSA Source track. |

## Cross Track Information

These pages describe information that crosses track boundaries.

| Page | Description |
|---|---|
| Threats & mitigations | Detailed information about specific supply chain attacks and how SLSA helps |
| Verified Properties | SLSA allows a common way to express verified properties that may not fit within a SLSA track. |

## Attestation formats

These pages include the concrete schemas for SLSA attestations. The Provenance and VSA formats are recommended, but not required by the specification.

| Page | Description |
|---|---|
| General model | General attestation mode |
| Provenance | Provides a description of the concept of provenance and links to the various tracks specific definitions. |
| Build Provenance | Suggested build provenance format and explanation |
| Verification Summary | Suggested VSA format and explanation |

# SLSA specification

SLSA is a specification for describing and incrementally improving supply chain security, established by industry consensus. It is organized into a series of levels that describe increasing security guarantees.

This is **Version 1.1** of the SLSA
specification, which defines several SLSA levels and recommended attestation
formats, including provenance.

## Understanding SLSA

These pages provide an overview of SLSA, how it helps protect against common supply chain attacks, and common use cases. If you’re new to SLSA or supply chain security, start here.

| Page | Description |
|---|---|
| What’s new in v1.1 | What’s new in SLSA Version 1.1 |
| About SLSA | An introductory guide to SLSA |
| Supply chain threats | An introduction to supply chain threats |
| Use cases | Use cases |
| Guiding principles | Use cases |
| FAQ | Questions and more information |
| Future directions | Additions and changes being considered for future SLSA versions |

## Core specification

These pages describe SLSA’s security levels and requirements for each track. If you want to achieve SLSA a particular level, these are the requirements you’ll need to meet.

| Page | Description |
|---|---|
| Terminology | Terminology and model used by SLSA |
| Security levels | Overview of SLSA’s tracks and levels, intended for all audiences |
| Producing artifacts | Detailed technical requirements for producing software artifacts, intended for platform implementers |
| Distributing provenance | Detailed technical requirements for distributing provenance, intended for platform implementers and software distributors |
| Verifying artifacts | Guidance for verifying software artifacts and their SLSA provenance, intended for platform implementers and software consumers |
| Verifying build platforms | Guidelines for securing SLSA Build L3+ builders, intended for platform implementers |
| Threats & mitigations | Detailed information about specific supply chain attacks and how SLSA helps |

## Attestation formats

These pages include the concrete schemas for SLSA attestations. The Provenance and VSA formats are recommended, but not required by the specification.

| Page | Description |
|---|---|
| General model | General attestation mode |
| Provenance | Suggested provenance format and explanation |
| Verification Summary | Suggested VSA format and explanation |

# SLSA specification

SLSA is a specification for describing and incrementally improving supply chain security, established by industry consensus. It is organized into a series of levels that describe increasing security guarantees.

This is the Working Draft of what the next version of the SLSA specification might be. It defines several SLSA levels and tracks, as well as recommended attestation formats, including provenance.

This is **Version 1.2** of the SLSA specification, which defines the SLSA Build
and Source tracks.

## Understanding SLSA

These pages provide an overview of SLSA, how it helps protect against common supply chain attacks, and common use cases. If you’re new to SLSA or supply chain security, start here.

| Page | Description |
|---|---|
| What’s new | The changes brought by this Working Draft. |
| About SLSA | An introductory guide to SLSA |
| Supply chain threats | An introduction to supply chain threats |
| Use cases | Use cases |
| Guiding principles | Use cases |
| FAQ | Questions and more information |
| Future directions | Additions and changes being considered for future SLSA versions |
| Tracks | Provides an overview of each track and links to more specific information. |

## Build Track

These pages describe the build track’s security levels and requirements. If you want to achieve a particular level of the SLSA build track these are the requirements you’ll need to meet.

| Page | Description |
|---|---|
| Basics | The SLSA build track is organized into a series of levels that provide increasing supply chain security guarantees. This gives you confidence that software hasn’t been tampered with and can be securely traced back to its source. This page is a descriptive overview of the SLSA build track levels, describing their intent. |
| Terminology | Terminology and model used by SLSA |
| Producing artifacts | Detailed technical requirements for producing software artifacts, intended for platform implementers |
| Distributing provenance | Detailed technical requirements for distributing provenance, intended for platform implementers and software distributors |
| Verifying artifacts | Guidance for verifying software artifacts and their SLSA provenance, intended for platform implementers and software consumers |
| Assessing build platforms | Guidelines for securing SLSA Build L3+ builders, intended for platform implementers |

## Build Environment Track

These pages describe the build environment track’s security levels and requirements. If you want to achieve a particular level of the SLSA build environment track these are the requirements you’ll need to meet.

| Page | Description |
|---|---|
| Attesting build environments | Overview of SLSA’s Attested Build Environment track, intended for all audiences |

## Dependency Track

This pages describes the dependency track’s security levels and requirements. If you want to achieve a particular level of the SLSA dependency track these are the requirements you’ll need to meet.

| Page | Description |
|---|---|
| Consuming dependencies | Overview of the Dependency track |

## Source Track

These pages describe the source track’s security levels and requirements. If you want to achieve a particular level of the SLSA source track these are the requirements you’ll need to meet.

| Page | Description |
|---|---|
| Producing source | Overview of the Source track |
| Verifying source | Guidelines for verifying source provenance |
| Assessing source control systems | Guidelines for assessing source control system security. |
| Example controls | This page provides examples of additional controls that organizations may want to implement as they adopt the SLSA Source track. |

## Cross Track Information

These pages describe information that crosses track boundaries.

| Page | Description |
|---|---|
| Threats & mitigations | Detailed information about specific supply chain attacks and how SLSA helps |
| Verified Properties | SLSA allows a common way to express verified properties that may not fit within a SLSA track. |

## Attestation formats

These pages include the concrete schemas for SLSA attestations. The Provenance and VSA formats are recommended, but not required by the specification.

| Page | Description |
|---|---|
| General model | General attestation mode |
| Provenance | Provides a description of the concept of provenance and links to the various tracks specific definitions. |
| Build Provenance | Suggested build provenance format and explanation |
| Verification Summary | Suggested VSA format and explanation |

# How to SLSA

These instructions tell you how to apply the core SLSA specification to use SLSA in your specific situation.

| Page | Description |
|---|---|
| For developers | How to apply SLSA requirements to your build |
| For organizations | How to apply SLSA to an organization |
| For infrastructure providers | How to implement SLSA in source, build, and package platforms |

# Specification Stages and Versioning

## Specification Stages

Specifications go through various stages of development from which you should have different expectations. This document defines the different stages the SLSA project uses and their meaning for readers and contributors.

Every specification page should prominently display a *Status* section
stating which stage the specification is in with a link to its
definition.

### Draft

This is the first stage of development a specification goes through. At this point, not much should be expected of it. The specification may be very incomplete and may change at any time and even be abandoned. It is therefore not suitable for reference or for implementation beyond experimentation.

A specification may be published several times during this stage as work progresses. The status section of the document may provide additional information as to its development status and whether reviews and feedback are welcome.

See the Governance for other considerations related to a Draft Specification.

### Candidate

At this stage the document is considered to be feature complete and is published as a way to invite final reviews. Editorial changes may still be made but no addition of new features is expected and short of problems being found no significant changes are expected to happen anymore.

See the Governance for other considerations related to a Candidate for Approved Specification.

### Approved

At this stage the document is considered stable. No changes that would constitute a significant departure from the existing specification are expected although changes to address ambiguities and edge cases may still occur.

See the Governance for other considerations related to an Approved Specification.

### Retired

This stage indicates that the specification is no longer maintained. It has been either rendered obsolete by a newer version or abandoned for some reason. The status of the document may provide additional information and point to the new document to use instead if there is one.

## Versioning

SLSA needs revision from time to time, so we version the specification to facilitate conformance efforts and prevent confusion. We assign a single version number to the Core Specification and Attestation Formats, collectively known as the SLSA Specification.

Given a version MAJOR.MINOR, we will increment

- MAJOR version when making backwards incompatible changes to an Attestation Format or adding new requirements to an existing level in the Core Specification.
- MINOR version when adding new tracks or levels to the Core Specification, modifying an existing level without fundamentally changing its meaning, or adding new fields to an Attestation Format in a backwards-compatible way. For more details on Attestation Format versioning, see the in-toto Versioning Rules.

Although we can revise the contents of level, we will never change a level’s
high-level meaning after publication (e.g. SLSA Build Level 2 will retain its
general meaning between major versions). If you require precision when referring
to SLSA levels, then include their version number using the syntax `SLSA [Track] [Level] ([Version])` (e.g. `SLSA Build L3 (v1.0)`). For more
details on SLSA versioning, see the
SLSA v1.0 Proposal.

There’s an active community of members, contributors and collaborators behind the SLSA framework. We’re drawn together by the shared goals of improving software supply chain security and codifying best practices for development, deployment and governance, all collaborating on an objective framework that works for open source projects and organizations, influences policy and regulations, empowers engineers and builds for the future.

## Get involved

The SLSA project is an open source project that strives to make useful and practical standards, tools, and documentation to reduce software supply chain risk in the real world. To succeed, we rely on contributors from a variety organizations to help us improve. Whether that’s reporting successes or challenges, contributing changes to the specification or documentation, or developing tooling, we welcome your contributions.

How to contribute

For questions, suggestions, or status updates, please use one of the following channels.

Contribution guidelines Specification meeting (weekly) GitHub issues (tracks all work) Slack (#slsa) Mailing listArchived meeting notes

The SLSA community no longer holds these meetings regularly. Their meeting notes are archived here for posterity.

Community meeting Tooling Special Interest Group“SLSA’s really the first of its kind, a framework for supply chain and build integrity. What sets it apart is the thriving community behind it, and it’s resonating with different organizations.”

Kim Lewandowski

Founder, Chainguard

## Steering committee

- Adrian Diglio - Microsoft
- Andrew McNamara - Red Hat
- Mike Lieberman - Kusari/CNCF
- Michael Winser - Eclipse Foundation
- Tom Hennen - Google

## Governance

SLSA is a community effort organized within the Open Source Security Foundation (OpenSSF) and released under the Community Specification License 1.0.

For more information about governance, see slsa-framework/governance.

-
 ## Mini Shai-Hulud: Where SLSA's Boundaries Fall
 The Mini Shai-Hulud npm supply chain attack chained GitHub Actions misconfiguration, cache poisoning, and OIDC token theft to publish malicious packages under legitimate identities. The compromised packages carried cryptographically valid SLSA build provenance attestations — but the build platform behind them did not meet SLSA Build L3 isolation requirements, and one that did would have blocked the primary attack vector. This post maps what SLSA's levels actually guarantee, where layering policy closes additional gaps, and what falls entirely outside the framework's scope.**by Andrew McNamara (Red Hat) on 15 May 2026**
-
 ## Supply Chain Robots, Electric Sheep, and SLSA**Guest post by Brett Smith on 18 Dec 2025**By Brett Smith, Distinguished Software Developer at SAS Institute
-
 ## Announcing SLSA v1.2**by SLSA Community on 24 Nov 2025**Today we’re pleased to announce the release of SLSA Version 1.2, the latest version of the SLSA specification.
-
 ## Announcing SLSA v1.2 Release Candidate 2**by SLSA Community on 10 Nov 2025**Today we’re releasing for public review SLSA Version 1.2 RC2, a Release Candidate of SLSA v1.2. We are seeking comments on this specification by November 24th, 2025. This release addresses most of the feedback we received on RC1.
-
 ## SLSA End-to-End With AMPEL & Friends**Guest post by Adolfo García Veytia (puerco) on 21 Oct 2025**This guest post walks through a practical, end-to-end SLSA implementation using 🔴🟡🟢 AMPEL — the Amazing Multipurpose Policy Engine (and L) — along with other tools in the supply chain security ecosystem. You’ll see how each step in a project’s build can be protected through attested data, using VSA receipts to capture and verify each step integrity along the way.
-
 ## SLSA End-to-End: Request for examples**by Andrew McNamara, Tom Hennen on 22 Jul 2025**This is a request for examples (RFE) for an end-to-end implementation of the Supply-chain Levels for Software Artifacts (SLSA) framework. The goal is to create a comprehensive demonstration of how SLSA can be used to secure the software supply chain, from source code to end-user consumption. These implementations will serve as a reference for the community, showcasing best practices and providing a clear adoption path for organizations looking to improve their software supply chain security.
-
 ## Announcing SLSA v1.2 Release Candidate 1**by SLSA Community on 20 Jun 2025**Today we’re releasing for public review SLSA Version 1.2 RC1, a Release Candidate of SLSA v1.2. We are seeking comments on this specification by July 18th, 2025.
-
 ## Source Track Sprint Recap**by Andrew McNamara (Red Hat), Tom Hennen (Google), Zachariah Cox (GitHub) on 29 Apr 2025**Last week the three of us met to try to make more progress on the source track. Async collaboration can work well for some things, but on “squishier” topics a higher-bandwidth engagement can be really helpful. All our work was against the draft version of the spec and before anything becomes official it will go through the approval process. We’d love your feedback on what we accomplished and discussed (summarized below), so please let us know what you think!
-
 ## SLSA v1.1 is now Approved!**by SLSA Community on 21 Apr 2025**Today we’re releasing SLSA Version 1.1 as the latest Approved Specification of SLSA, effectively replacing Version 1.0.
-
 ## Announcing SLSA v1.1 Release Candidate 2**by SLSA Community on 04 Apr 2025**Today we’re releasing the SLSA Version 1.1 RC2 for public review. We are seeking comments on these spec changes by April 18, 2025. This release brings several changes aimed at enhancing the clarity and usability of the original specification. It also introduces backwards-compatible clarifications to the SLSA threat model, attestation model and verification procedure. This includes the addition of verifier metadata to the Verification Summary Attestation (VSA) format. Please, refer to the What’s new section for further details.
-
 ## Defender's Perspective: Dependency Confusion and Typosquatting Attacks**by Meder Kydyraliev (Google) on 13 Aug 2024***Dependency confusion*and*typosquatting*attacks are very similar in their nature. They both exploit the weakness in the way many package managers identify packages using only their names. Successfully exploiting this weakness enables the attacker to run arbitrary code at install time or at application’s run time. These attacks are scalable, portable, and extremely cost-effective to carry out—making them very appealing to malicious actors.
-
 ## Securing software artifacts with Tekton Chains and IBM's DevSecOps**Guest post by Arnaud J Le Hors on 16 Apr 2024**Tekton Chains, and the IBM DevSecOps offering that builds on it, can now be used to secure software artifacts with SLSA.
-
 ## Build your own SLSA 3+ provenance builder on GitHub Actions**by Andres Almiray (JReleaser), Adam Korczynski (Ada Logics), Philip Harrison (GitHub), Laurent Simon (Google) on 28 Aug 2023**It has been an exciting quarter for supply chain security and SLSA, with the release of the SLSA v1.0 specification, SLSA provenance support for npm, and the announcement of new SLSA Level 3 builders for Node.js and containers!
-
 ## Announcing Container-based SLSA 3 Builder on GitHub Actions**by Asra Ali, Razieh Behjati, Tiziano Santoro (Google) on 13 Jun 2023**Following the recent launch of SLSA v1.0, we’re announcing a new, GitHub Actions workflow that achieves SLSA Build Track Level 3 for provenance generation. This lets users generate unforgeable provenance, allowing consumers to trust *and*verify how their software artifacts were built. The container-based SLSA 3 builder is the result of a collaboration between the Google Open Source Security Team (GOSST), the SLSA community, and Project Oak.
-
 ## Bringing Improved Supply Chain Security to the Node.js Ecosystem**by Ian Lewis & Laurent Simon (Google), Fredrik Skogman (GitHub) on 11 May 2023**It has been a big month for supply chain security! GitHub recently announced the public beta for npm package provenance. This adds new functionality to npmjs.com and the npm CLI that allows package maintainers to generate and upload SLSA Build Level 2 provenance along with their packages. Integration with Sigstore enables verification of signature and certificate metadata so users know that the package came from the expected source repository.
-
 ## in-toto and SLSA**Guest post by Aditya Sirish (NYU) and Tom Hennen (Google) representing the in-toto Community on 02 May 2023**As an adopter of SLSA, you have likely encountered the in-toto project. in-toto attestations are part of SLSA’s recommended suite for expressing software supply chain claims. As in-toto maintainers, we’ve interacted with a number of people who know of in-toto through SLSA but don’t fully understand the project. For example, some were surprised to learn that “in-toto isn’t just a format of attestations”, and that the framework also defines verification workflows that make use of attestations that are not SLSA Provenance. So, we decided to author this post as a quick primer on in-toto, how SLSA uses in-toto, and how other attestations can be used to complement SLSA.
-
 ## SLSA v1.0 is now final!**by Mark Lodato on 19 Apr 2023**After almost two years since SLSA’s initial preview release, we are pleased to announce our first official stable version, SLSA v1.0! The full announcement can be found at the OpenSSF press release, and a description of changes can be found at What’s new in v1.0. Thank you to all members of the SLSA community who made this possible through your feedback, suggestions, discussions, and pull requests!
-
 ## Announcing SLSA v1.0 Release Candidate 2**by SLSA Community on 04 Apr 2023**We’re excited to announce SLSA v1.0 Release Candidate 2 (RC2) following the valuable feedback we received on the first release candidate. This is intended to be the final release candidate before marking v1.0 as an Approved Specification.
-
 ## The Breadth and Depth of SLSA**by Mike Lieberman on 03 Apr 2023**Interested in getting involved? Now’s the chance to provide your feedback on the foundational v1 release of the SLSA framework.
-
 ## Announcing SLSA v1.0 Release Candidate**by Mark Lodato, Kris Kooi, Joshua Lock on 24 Feb 2023**Today, we are excited to announce the important milestone of a release candidate (RC) SLSA Specification. This is the first major update to SLSA since its v0.1 release in June 2021, and the RC finalizes multiple revisions to the SLSA specifications and requirements. We’re grateful for the huge community engagement that went into shaping this work.
-
 ## General availability of SLSA 3 Container Generator for GitHub Actions**by Asra Ali, Ian Lewis, Laurent Simon on 01 Feb 2023**Today, we are announcing the general availability of the SLSA 3 Container Generator for GitHub Actions starting with v1.4.0. This free tool allows any GitHub project to produce SLSA level 3 compliant provenance statements so users can verify the origin of container images they use. While previous tools allowed users to generate provenance for file artifacts, the Container Generator is able to support container ecosystems. It does this by allowing provenance statements to be distributed alongside your images in a container registry and integrating directly with Sigstore-compatible tooling for inspection and verification.
-
 ## Safeguarding builds on Google Cloud Build with SLSA**Guest post by Asra Ali, Ian Lewis, Laurent Simon, Stephen Anastos on 05 Dec 2022**Earlier this year, Google Cloud Build (GCB) announced support for Level 3 assurance of Supply-chain Levels for Software Artifacts (SLSA) for container images. Users can now automatically generate verifiable provenance documents (build records) of builds that take place in Cloud Build. Provenance can be used to provide assurance that a trusted builder (in this case, GCB) produced the resulting image through some declared process with trusted source material. To make verification effortless, we are announcing support for verifying the provenance document in the open-source slsa-verifier CLI tool, which previously only had support for GitHub Actions. With the slsa-verifier, everyone — not just the container authors — can verify the SLSA provenance document.
-
 ## Executive Order on Secure Supply Chain — in Plain English**Guest post by Isaac Hepworth on 26 Sep 2022**
-
 ## General availability of SLSA3 Generic Generator for GitHub Actions**by Ian Lewis, Laurent Simon, Asra Ali on 29 Aug 2022**A few months ago Google and GitHub announced the release of a Go builder that would help software developers and consumers more easily verify the origins of software by using verification files known as provenance. Since then, the SLSA community has been working to enable provenance generation for other projects that may use any number of languages or build tools. Today, we’re pleased to announce that we’re adding a new tool to generate similar provenance documents for projects developed in any programming language, while keeping your existing building workflows.
-
 ## All about that Base(line): How Cybersecurity Frameworks are Evolving with Foundational Guidance**Guest post by Jennifer Privette on 25 Jul 2022***In coordination with Aaron Bacchi, Emmy Eide, Melba Lopez, Brandon Lum, and Moshe Zioni*
-
 ## General Availability of SLSA 3 Go native builder for GitHub Actions**by Laurent Simon, Asra Ali, Ian Lewis, Mark Lodato, Jose Palafox, Joshua Lock on 20 Jun 2022**
-
 ## SLSA for Success: Using SLSA to help achieve NIST’s SSDF**Guest post by Isaac Hepworth, Meder Kydyraliev, Brandon Lum on 15 Jun 2022**Since February’s release of the latest version of the Secure Software Development Framework’s (SSDF), software organizations have been poring over the dozens of best practices and tasks laid out by the National Institute of Standards and Technology (NIST) in response to last year’s Executive Order on Cybersecurity. Implementation is tough, though: the guidelines cover organizations of all sizes, cybersecurity sophistication, and operating environment. The descriptive requirements are not prioritized and explicitly not meant to be a checklist to follow. Each organization must find ways to interpret the recommendations for their particular needs.
-
 ## SBOM + SLSA: Accelerating SBOM success with the help of SLSA**Guest post by Brandon Lum, Isaac Hepworth, Meder Kydyraliev on 02 May 2022**
-
 ## SLSA Is No Free Lunch**Guest post by Mike Lieberman on 11 Apr 2022**“What is SLSA?” followed closely by “What does SLSA do for me?” are the two most common questions I get when people learn about SLSA. This has led to a lot of confusion as to how folks apply SLSA, and the benefits they get. You can’t just apply SLSA practices to a pipeline that runs a build, generate a SLSA attestation and magically be protected from supply chain compromise. Contrary to a lot of the hype being thrown around, SLSA is no free lunch, and we must help protect our lunch!
-
 ## Introducing the SLSA Blog**by SLSA Community on 08 Apr 2022**We’re excited to launch our very own blog, from which we will be posting project news, documentation, and other information about SLSA. Stay tuned for more posts coming your way soon.

# Supply chain threats

Attacks can occur at every link in a typical software supply chain, and these kinds of attacks are increasingly public, disruptive, and costly in today’s environment.

This page is an introduction to possible attacks throughout the supply chain and how SLSA could help. For a more technical discussion, see Threats & mitigations.

## Summary

**Note that SLSA does not currently address all of the threats presented here.**
See Threats & mitigations for what is currently addressed and
Terminology for an explanation of the supply chain model.

SLSA’s primary focus is supply chain integrity, with a secondary focus on availability. Integrity means protection against tampering or unauthorized modification at any stage of the software lifecycle. Within SLSA, we divide integrity into source integrity vs build integrity.

**Source integrity:** Ensure that the source revision represents the intent of the producer, that all expected processes were followed and that the revision was not modified after being accepted.

**Build integrity:** Ensure that the package is built from the correct,
unmodified sources and dependencies according to the build recipe defined by the
software producer, and that artifacts are not modified as they pass between
development stages.

**Availability:** Ensure that the package can continue to be built and
maintained in the future, and that all code and change history is available for
investigations and incident response.

### Real-world examples

Many recent high-profile attacks were consequences of supply chain integrity vulnerabilities, and could have been prevented by SLSA’s framework. For example:

| Threats from | Known example | How SLSA could help | |
|---|---|---|---|
| A | Producer | SpySheriff: Software producer purports to offer anti-spyware software, but that software is actually malicious. | SLSA does not directly address this threat but could make it easier to discover malicious behavior in open source software, by forcing it into the publicly available source code. For closed source software SLSA does not provide any solutions for malicious producers. |
| B | Authoring & reviewing | SushiSwap: Contractor with repository access pushed a malicious commit redirecting cryptocurrency to themself. | Two-person review could have caught the unauthorized change. |
| C | Source code management | PHP: Attacker compromised PHP's self-hosted git server and injected two malicious commits. | A better-protected source code system would have been a much harder target for the attackers. |
| D | External build parameters | The Great Suspender: Attacker published software that was not built from the purported sources. | A SLSA-compliant build server would have produced provenance identifying the actual sources used, allowing consumers to detect such tampering. |
| E | Build process | SolarWinds: Attacker compromised the build platform and installed an implant that injected malicious behavior during each build. | Higher SLSA Build levels have stronger security requirements for the build platform, making it more difficult for an attacker to forge the SLSA provenance and gain persistence. |
| F | Artifact publication | CodeCov: Attacker used leaked credentials to upload a malicious artifact to a GCS bucket, from which users download directly. | Provenance of the artifact in the GCS bucket would have shown that the artifact was not built in the expected manner from the expected source repo. |
| G | Distribution channel | Attacks on Package Mirrors: Researcher ran mirrors for several popular package registries, which could have been used to serve malicious packages. | Similar to above (F), provenance of the malicious artifacts would have shown that they were not built as expected or from the expected source repo. |
| H | Package selection | Browserify typosquatting: Attacker uploaded a malicious package with a similar name as the original. | SLSA does not directly address this threat, but provenance linking back to source control can enable and enhance other solutions. |
| I | Usage | Default credentials: Attacker could leverage default credentials to access sensitive data. | SLSA does not address this threat. |
| N/A | Dependency threats (i.e. A-H, recursively) | event-stream: Attacker controls an innocuous dependency and publishes a malicious binary version without a corresponding update to the source code. | Applying SLSA recursively to all dependencies would prevent this particular vector, because the provenance would indicate that it either wasn't built from a proper builder or that the binary did not match the source. |

| Availability threat | Known example | How SLSA could help | |
|---|---|---|---|
| N/A | Dependency becomes unavailable | Mimemagic: Producer intentionally removes package or version of package from repository with no warning. Network errors or service outages may also make packages unavailable temporarily. | SLSA does not directly address this threat. |

A SLSA level helps give consumers confidence that software has not been tampered with and can be securely traced back to source—something that is difficult, if not impossible, to do with most software today.

# Tracks

SLSA is composed of multiple tracks which are each composed of multiple levels. Each track addresses different threats and has its own set of requirements and patterns of use.

## Build Track

The SLSA build track describes increasing levels of trustworthiness and completeness in a package artifact’s provenance. Provenance describes what entity built the artifact, what process they used, and what the inputs were. The lowest level only requires the provenance to exist, while higher levels provide increasing protection against tampering of the build, the provenance, or the artifact.

The primary purpose of the build track is to enable verification that the artifact was built as expected. Consumers have some way of knowing what the expected provenance should look like for a given package and then compare each package artifact’s actual provenance to those expectations. Doing so prevents several classes of supply chain threats.

Each ecosystem (for open source) or organization (for closed source) defines exactly how this is implemented, including: means of defining expectations, what provenance format is accepted, whether reproducible builds are used, how provenance is distributed, when verification happens, and what happens on failure. Guidelines for implementers can be found in the requirements.

## Source Track

The SLSA source track provides producers and consumers with increasing levels of trust in the source code they produce and consume. It describes increasing levels of trustworthiness and completeness of how a source revision was created.

The expected process for creating a new revision is determined solely by that repository’s owner (the organization) who also determines the intent of the software in the repository and administers technical controls to enforce the process.

Consumers can review attestations to verify whether a particular revision meets their standards.

If you’re looking to jump straight in and try SLSA, here’s a quick start guide for the steps to take to reach the different SLSA levels.

## Choosing your SLSA level

For all SLSA levels, you follow the same steps:

- Generate provenance, i.e., document your build process
- Make the provenance available, to allow downstream users to verify it

What differs for each level is the robustness of the build and provenance. For more information about provenance see https://slsa.dev/provenance/.

The tools discussed in this guide are freely available and believed to meet SLSA expectations. The Builder SLSA levels section provides a more complete list. If you think a build option is misclassified or want to add one, please open an issue or submit a PR against this page.

SLSA levels are progressive: SLSA 3 includes all the guarantees of SLSA 2, and SLSA 2 includes all the guarantees of SLSA 1. Currently, though, the work required to achieve lower SLSA levels will not necessarily accrue toward the work needed for higher levels, because achieving a higher level may require migrating to a different build platform altogether. For that reason, **you should start with the highest level that’s possible for your project or organization to avoid wasted work**.

The SLSA level you implement depends on your current build situation:

- If you are using GitHub Actions, jump directly to SLSA 3. You do not need to implement SLSA 1 or SLSA 2.
- If you are using FRSCA, jump to SLSA 2. You do not need to implement SLSA 1.
- If you’re using any other build platform, consult its documentation to find out what SLSA level it supports and how to proceed. If your build platform doesn’t support the SLSA level you are aiming for or you do not use any build platform, you should consider adopting one that does, such as GitHub Actions and FRSCA discussed below.

### Provenance verification

Various build methods require different methods to verify provenance. The SLSA community already has some verification methods and is working on additional solutions for provenance verification. Even if there is currently no simple verification method for a particular build platform, adopting these builders now means you will be “SLSA ready” when more ecosystem tooling is released. Note that when you create provenance, you must supply to users a method for interpreting what you have created.

### Provenance formatting

Your provenance format depends on who will be consuming it. See provenance for an explanation of which format to choose.

### Provenance storage

Containers have a standard place to put the provenance in the OCI container registry. With time, the SLSA community hopes to create standard ecosystem-based repositories for provenance. For now, the convention is to keep the provenance attestation with your artifact. Though Sigstore is becoming more and more popular, the format of the provenance is currently tool-specific.

## SLSA 1

As mentioned before, if you don’t already use a build platform or CI/CD, you should consider adopting a platform that supports SLSA 2 or SLSA 3. This will make the following steps easier and provide for higher SLSA levels later on. Individual developers who wish to put a minimal amount of security on their builds can use SLSA 1.

SLSA 1 requires that the build process is documented. Some tools suggested below also support signed provenance. Though not required for SLSA 1, signing your provenance increases trust in the document by showing that it has not been tampered with.

### Tooling

A build configuration file (i.e., GitHub workflow) qualifies for SLSA 1. It would be considered unsigned, unformatted provenance.

### Build platform plugins or extensions

The following options work with your build platform to produce unsigned, formatted provenance. They do not qualify for SLSA 2 because they are unsigned and not run by the hosted server:

Downstream users may verify the provenance with Cue Policies.

### Build observers with hosted platforms

The following options are user-configured inside a hosted platform. They observe the build process and produce signed, formatted provenance. These options do not qualify for SLSA 2 because they are configured by users, not the hosted platform.

Downstream users may verify the provenance with Cue policies and the signature with Cosign.

**Note:** If you are using one of these options with GitHub Actions, jump to SLSA 3 and use the builder itself to generate provenance.

- If you’re using Tejolote with GitHub Actions, jump to SLSA 3 and generate provenance directly from the builder.
- Tekton Chains – custom resource definition controller that can generate provenance for Kubernetes OCI containers

## SLSA 2

To achieve SLSA 2, the goals are to:

- Run your build on a hosted platform that generates and signs provenance
- Publish the provenance to allow downstream users to verify it

The following is a SLSA 2 builder:

FRSCA is an OpenSSF project that aims at offering a full build pipeline. It is not yet generally available. It qualifies as a SLSA 2 builder because regular users of the platform are not able to inject or alter the contents of the provenance it generates. FRSCA produces signed, formatted provenance that can be verified by the generic SLSA verifier.

## SLSA 3

To achieve SLSA 3, you must:

- Run your build on a hosted platform that generates and signs provenance
- Ensure that build runs cannot influence each other
- Produce signed provenance that can be verified as authentic

### GitHub Actions

If you are building on GitHub Actions, adopt the Language-agnostic GitHub provenance generator / builder to make your build qualify for SLSA 3. Consumers can use the Generic SLSA Verifier for provenance verification.

## Builder SLSA levels

The following table shows known build software packages and the potential SLSA level they are believed to qualify for. Potential levels are reached with proper provenance. Again, if you think a build option is misclassified or want to add one, please open an issue or submit a PR against this page.

Note that this list is provided “as is”. OpenSSF makes no claim as to the reliability of this information. A certification program is under development to provide a more definitive list.

| Builder | Potential SLSA Level |
|---|---|
| FRSCA | 2 |
| GitHub Actions | 3 |
| Google Cloud Build | 3 |
| No Hosted Build Platform | 1 |

# SLSA specification

SLSA is a specification for describing and incrementally improving supply chain security, established by industry consensus. It is organized into a series of levels that describe increasing security guarantees.

This is **Version 1.2** of the SLSA specification. It defines several SLSA
levels and tracks, as well as recommended attestation formats, including
provenance.

## Understanding SLSA

These pages provide an overview of SLSA, how it helps protect against common supply chain attacks, and common use cases. If you’re new to SLSA or supply chain security, start here.

| Page | Description |
|---|---|
| What’s new | The changes brought by this revision of SLSA. |
| About SLSA | An introductory guide to SLSA |
| Supply chain threats | An introduction to supply chain threats |
| Use cases | Use cases |
| Guiding principles | Use cases |
| FAQ | Questions and more information |
| Future directions | Additions and changes being considered for future SLSA versions |
| Tracks | Provides an overview of each track and links to more specific information. |

## Build Track

These pages describe the build track’s security levels and requirements. If you want to achieve a particular level of the SLSA build track these are the requirements you’ll need to meet.

| Page | Description |
|---|---|
| Basics | The SLSA build track is organized into a series of levels that provide increasing supply chain security guarantees. This gives you confidence that software hasn’t been tampered with and can be securely traced back to its source. This page is a descriptive overview of the SLSA build track levels, describing their intent. |
| Terminology | Terminology and model used by SLSA |
| Producing artifacts | Detailed technical requirements for producing software artifacts, intended for platform implementers |
| Distributing provenance | Detailed technical requirements for distributing provenance, intended for platform implementers and software distributors |
| Verifying artifacts | Guidance for verifying software artifacts and their SLSA provenance, intended for platform implementers and software consumers |
| Assessing build platforms | Guidelines for securing SLSA Build L3+ builders, intended for platform implementers |

## Source Track

These pages describe the source track’s security levels and requirements. If you want to achieve a particular level of the SLSA source track these are the requirements you’ll need to meet.

| Page | Description |
|---|---|
| Producing source | Overview of the Source track |
| Verifying source | Guidelines for verifying source provenance |
| Assessing source control systems | Guidelines for assessing source control system security. |
| Example controls | This page provides examples of additional controls that organizations may want to implement as they adopt the SLSA Source track. |

## Cross Track Information

These pages describe information that crosses track boundaries.

| Page | Description |
|---|---|
| Threats & mitigations | Detailed information about specific supply chain attacks and how SLSA helps |
| Verified Properties | SLSA allows a common way to express verified properties that may not fit within a SLSA track. |

## Attestation formats

These pages include the concrete schemas for SLSA attestations. The Provenance and VSA formats are recommended, but not required by the specification.

| Page | Description |
|---|---|
| General model | General attestation mode |
| Provenance | Provides a description of the concept of provenance and links to the various tracks specific definitions. |
| Build Provenance | Suggested build provenance format and explanation |
| Verification Summary | Suggested VSA format and explanation |
