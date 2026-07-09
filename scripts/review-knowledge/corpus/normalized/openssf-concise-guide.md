# Concise Guide for Developing More Secure Software

*by the Open Source Security Foundation (OpenSSF) Best Practices Working Group, 2023-06-14*

Here is a concise guide for all software developers for secure software development, building, and distribution. All tools or services listed are merely examples.

- **Ensure all privileged developers use**- **multi-factor authentication (MFA) tokens**. This includes those with commit or accept privileges. MFA hinders attackers from “taking over” these accounts.
- **Learn about secure software development.**Take, e.g., the free OpenSSF course or the hands-on Security Knowledge Framework course. SAFECode’s Fundamental Practices for Secure Software Development provides a helpful summary.
- **Use a combination of tools in your CI pipeline to detect vulnerabilities**. See the OpenSSF guide to security tools. Tools shouldn’t be the- *only*mechanism, but they scale.
- **Evaluate software before selecting it as a direct dependency**. Only add it if needed, evaluate it (see Concise Guide for Evaluating Open Source Software), double-check its name (to counter typosquatting), and ensure it’s retrieved from the correct repository.
- **Use package managers**. Use package managers (system, language-level, and/or container-level) to automatically manage dependencies and enable rapid updates.
- **Implement automated tests**. Include negative tests (tests that what shouldn’t happen doesn’t happen) and ensure the test suite is thorough enough to “ship if it passes the tests”.
- **Monitor known vulnerabilities in your software’s direct & indirect dependencies**. E.g., enable basic scanning via GitHub’s dependabot or GitLab dependency scanning. Many other third party Software Composition Analysis (SCA) tools are also available. Quickly update vulnerable dependencies.
- **Keep dependencies reasonably up-to-date**. Otherwise, it’s hard to update for vulnerabilities.
- **Do not push secrets to a repository**. Use tools to detect pushing secrets to a repository.
- **Review before accepting changes**. Enforce this, e.g., using GitHub or GitLab protected branches or an equivalent GitHub ruleset.
- **Prominently document how to report vulnerabilities & prepare for them**.- Use resources like the Guide to coordinated vulnerability disclosure.
- Explicitly disclose security issues affecting vendored dependencies.
- Create a security policy. Provide contacts.

- **Make it easy for your users to update**. Implement stable APIs, e.g., support old names when new ones are added. Use semantic versioning. Have a deprecation process.
- **Sign your project’s important releases**. Use standard tools and signing formats for your distribution. See the cosign tool from the sigstore project to sign containers and other artifacts.
- **Earn an OpenSSF Best Practices badge**for your open source project. At least earn “passing”. Plan and roadmap to eventually earn silver & gold.
- **Improve your**- **OpenSSF Scorecards**- **score (if OSS and on GitHub)**. You can read the Scorecards checks. Use the Allstar monitor.
- **Notify the community of vulnerabilities in your project.**Publish security advisories with accurate & precise information, e.g., what usage & versions are vulnerable, mitigations, and fixed version(s). Get a CVE ID. On GitHub, create your security advisory & request a CVE.
- **Improve your**- **Supply chain Levels for Software Artifacts (SLSA)**- **level**. This hardens the integrity of your build and distribution process against attacks.
- **Publish and consume a software bill of materials (SBOM)**. This lets users verify inventory, id known vulnerabilities, & id potential legal issues. Consider- **SPDX**or- **CycloneDX**. See our SBOM-Everywhere catalog.
- **Onboard your project into**- **LFX Security**- **if you manage a Linux Foundation project**.
- **Apply the**- **CNCF Security TAG Software Supply Chain Best Practices guide**.
- **Implement**- **ASVS**- **and follow relevant**- **cheatsheets**.
- **Apply SAFECode’s**- **Fundamental Practices for Secure Software Development**.
- **Complete a third-party security code review/audit**. Expect this to be USD$50K or more.
- **Continuously improve**. Improve scores, look for tips, & apply as appropriate.
- **Manage succession**. Have clear governance & work to add active, trustworthy maintainer(s).
- **Prefer memory-safe languages**. Many vulnerabilities involve memory safety. Where practical, use memory-safe programming languages (most are) and keep memory safety enabled. Otherwise, use mechanisms like compiler flags, extra tools, and peer review to reduce risk; see Compiler Options Hardening Guide for C and C++.
- **If a source code (unbuilt) package is released, it should only include content from the version control system (VCS), and source package users should rebuild, if needed, to create production (built) package(s)**. E.g., if autotools is used, if a source package is released it should- *not*include a generated- `configure`file, while recipients should ignore pre-generated files like- `configure`and instead rebuild from source (e.g., with- `autoreconf`). This eliminates a malware-hiding mechanism, as illustrated by an attack on xz utils.
- **Ensure production websites only load assets from your own domains**.- *Linking*to other domains is fine, but where practical, don’t directly load assets such as JavaScript, CSS, and media (including images) from domains you do not control. If you do, your site might be subverted if that other domain is subverted, so investigate the risks before doing so. See the subverted polyfill.io revelation in 2024.
- **Apply focused security guidelines**. Consult focused OpenSSF guides as applicable, such as the Compiler Options Hardening Guide for C and C++, npm Best Practices Guide, and the guide to Correctly Using Regular Expressions for Secure Input Validation.

We welcome suggestions and updates! Please open an issue or post a pull request.
