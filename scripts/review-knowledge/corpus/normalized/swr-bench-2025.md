SWR-Bench: Assessing LLM Performance in Real-World Code
Review Comment Generation
ZHENGRAN ZENG∗, Peking University, China
RUIKAI SHI∗, Peking University, China
KEKE HAN, Peking University, China
YIXIN LI, Peking University, China
KAICHENG SUN, Northwestern Polytechnical University, China
YIDONG WANG, Peking University, China
ZHUOHAO YU, Peking University, China
RUI XIE†, Peking University, China
WEI YE†, Peking University, China
SHIKUN ZHANG†, Peking University, China
Automated Code Review (ACR) is crucial for software quality, yet existing benchmarks often fail to re-
flect real-world complexities, hindering the evaluation of modern Large Language Models (LLMs). Current
benchmarks frequently focus on fine-grained code units, lack complete project context, and use inadequate
evaluation metrics. To address these limitations, we introduce SWR-Bench, a new benchmark comprising
1000 manually verified Pull Requests (PRs) from GitHub, offering PR-centric review with full project context.
SWR-Bench employs an objective LLM-based evaluation method that aligns strongly with human judgment
(∼90% agreement) by verifying if issues from a structured ground truth are covered in generated reviews.
Our systematic evaluation of mainstream ACR tools and LLMs on SWR-Bench reveals that current systems
underperform, and ACR tools are more adept at detecting functional errors. Subsequently, we propose and
validate a simple multi-review aggregation strategy that significantly boosts ACR performance, increasing F1
scores by up to 43.67%. Our contributions include the SWR-Bench benchmark, its objective evaluation method,
a comprehensive study of current ACR capabilities, and an effective enhancement approach, offering valuable
insights for advancing ACR research.
CCS Concepts: • Software and its engineering →Maintaining software.
Additional Key Words and Phrases: Automated Code Review, Large Language Models, Benchmark
ACM Reference Format:
Zhengran Zeng, Ruikai Shi, Keke Han, Yixin Li, Kaicheng Sun, Yidong Wang, Zhuohao Yu, Rui Xie, Wei Ye,
and Shikun Zhang. 2026. SWR-Bench: Assessing LLM Performance in Real-World Code Review Comment
Generation. Proc. ACM Softw. Eng. 3, FSE, Article FSE137 (July 2026), 23 pages. https://doi.org/10.1145/3808144
∗Both authors contributed equally to this research.
†Those authors are the corresponding authors.
Authors’ Contact Information: Zhengran Zeng, Peking University, Beijing, China, zhengranzeng@stu.pku.edu.cn; Ruikai
Shi, Peking University, Beijing, China, rkshi25@stu.pku.edu.cn; Keke Han, Peking University, Beijing, China, kkhan25@
stu.pku.edu.cn; Yixin Li, Peking University, Beijing, China, leason_lyx@stu.pku.edu.cn; Kaicheng Sun, Northwestern
Polytechnical University, Xian, China, sunkaicheng@mail.nwpu.edu.cn; Yidong Wang, Peking University, Beijing, China,
2301110730@stu.pku.edu.cn; Zhuohao Yu, Peking University, Beijing, China, zyu@stu.pku.edu.cn; Rui Xie, Peking University,
Beijing, China, ruixie@pku.edu.cn; Wei Ye, Peking University, Beijing, China, wye@pku.edu.cn; Shikun Zhang, Peking
University, Beijing, China, zhangsk@pku.edu.cn.
This work is licensed under a Creative Commons Attribution 4.0 International License.
© 2026 Copyright held by the owner/author(s).
ACM 2994-970X/2026/7-ARTFSE137
https://doi.org/10.1145/3808144
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.
arXiv:2509.01494v2 [cs.SE] 5 Jun 2026

FSE137:2
Zhengran Zeng, Ruikai Shi, Keke Han, Yixin Li, KC Sun, YD Wang, ZH Yu, Rui Xie, Wei Ye and Shikun Zhang
1
Introduction
Code Review (CR) is an indispensable quality assurance practice in software development [2, 3,
6, 24, 49], aimed at identifying and rectifying potential issues in code changes before integration.
While crucial for enhancing code quality, traditional manual code review faces significant hurdles
in modern software development. The increasing scale and complexity of projects mean that
manual reviews are time-consuming, contributing to development costs and potential delays in
feature releases [7, 10, 32]. To mitigate these challenges, Automated Code Review (ACR) Comment
Generation has emerged as a vital field and a prominent research area within software engineering [9,
22, 39, 41, 42], driving the creation of numerous tools and datasets. These ACR works strive to
improve the efficiency and effectiveness of the code review process by shortening feedback cycles,
reducing the manual burden on developers, and improving the consistency and scope of reviews,
thus playing a crucial role in streamlining software development workflows and maintaining high
standards of code quality. More recently, the explosive advancements in Large Language Models
(LLMs) have catalyzed a significant shift in ACR research [15, 21, 26, 35], with a growing emphasis
on LLM-based methodologies.
ACR Tool
Diff Hunk
+ def new_demo_func(item):
+ ... # new logic of func
+ return item
- def old_demo_func(item):
- ... # old logic of func
- return item
 src/demo.py
 Review comment
for this hunk.
import os
...
def old_demo_func(item):
 ... # old logic of func
 return item
...
 src/demo.py
Original Source Code
(a) Traditional hunk-level review with limited con-
text.
ACR Tool
Pull Request
 src
 main.py
 demo.py
 utils
·
+20-12
 Review comment
for this PR.
 src main.py
 utils demo.py
 examples reqs.txt
 README.md
Codebase
(b) SWR-Bench PR-level review with full codebase
context.
Fig. 1. A comparison of review paradigms.
However, a substantial portion of existing ACR benchmarks [22, 39, 41, 42] were established
primarily for deep learning models that predate the widespread capabilities and sophisticated
understanding of modern LLMs. This temporal and methodological gap presents challenges in
evaluating the true potential of contemporary LLM-based ACR tools using these benchmarks. This
discrepancy is particularly evident in the fundamental setup of the review task itself, as illustrated
in Figure 1. Specifically, existing representative benchmarks [22, 39, 41, 42] primarily center their
evaluation on isolated code changes, or diff hunks, with limited contextual information, such as only
a few related source code files (see Figure 1a). This hunk-based approach fails to mirror real-world
developer practices, where an entire Pull Request (PR) is reviewed as a cohesive unit. Consequently,
this narrow scope can lead to missing critical inter-dependency bugs.Meanwhile, these benchmarks
typically do not provide the complete project codebase as necessary context, making it difficult for
ACR tools to understand the global impact of code changes. In contrast, our paradigm provides the
full context (Figure 1b). More critically, the evaluation methodologies for these benchmarks are
also problematic. They generally rely on traditional natural language generation metrics such as
𝐸𝑥𝑎𝑐𝑡-𝑀𝑎𝑡𝑐ℎ[42], 𝐵𝑙𝑒𝑢scores [27] and simple LLM-based ratings [13, 15, 44]. However, metrics like
𝐵𝑙𝑒𝑢primarily measure textual similarity, while the reliability of LLM ratings is often questioned
due to potential inherent biases; both approaches have been shown to be severely inadequate
in assessing whether the crucial issues are truly identified in code review comments [17, 25].
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

SWR-Bench: Assessing LLM Performance in Real-World Code Review Comment Generation
FSE137:3
Furthermore, some studies [9, 35] utilize expensive human annotation, but this approach is difficult
to scale.
Consequently, evaluations on current benchmarks often misrepresent the true capabilities and
practical value of LLMs in complex, real-world code review. To bridge this critical gap, we introduce
SWR-Bench (Software Review Benchmark), a software review benchmark designed for more
realistic evaluation. Specifically, SWR-Bench comprises 1000 manually verified PR instances from
GitHub open-source projects. It directly confronts the issues in existing benchmarks through its
core design features:
• PR-centric Review: Each code review instance is based on a complete PR. This setup is inten-
tionally more challenging, yet more aligned with real-world practices, as we aim to assess the
ability of ACR tools for end-to-end code review. This includes identifying areas requiring review,
pinpointing specific issues, and generating meaningful review comments.
• Comprehensive Context: To simulate a realistic review environment, our benchmark provides
all relevant code changes (commits) and a snapshot of the entire project repository for each
PR, offering comprehensive context that enables tools to accurately assess the global impact of
modifications.
• Objective LLM-based Evaluation: We employ an automated evaluation where an LLM ob-
jectively verifies if issues from a structured ground truth report are covered in the generated
report. This method offers a more objective and semantically relevant assessment than subjective
“LLM-as-judge” ratings [13, 15, 44] or text-similarity metrics [27, 42].
After establishing SWR-Bench, we first validated its evaluation method, confirming strong (∼90%)
agreement with human judgment. We then systematically evaluated mainstream ACR tools and
leading LLMs on SWR-Bench, uncovering key performance aspects in realistic scenarios: 1) current
tools and LLMs do not yet perform sufficiently well; 2) furthermore, we found that LLMs trained with
a reasoning-focused approach exhibit better code review capabilities; and finally, 3) existing tools
are more adept at detecting functional errors, such as bugs, compared to identifying non-functional
issues like outdated documentation. Motivated by these observations, we further introduce a
simple enhancement scheme for LLM-based automated code review. This straightforward strategy
empowers an LLM to synthesize feedback from multiple review sources by analyzing, filtering
valid points, and discarding erroneous suggestions to produce a superior final report. Experiments
demonstrate this strategy substantially boosts ACR performance (increasing issue detection 𝐹1
score by up to 43.67%), suggesting promising avenues for future research. Therefore, our main
contributions in this work include:
1. Benchmark: We introduce SWR-Bench, a benchmark comprising 1,000 real-world PRs with full
project context, providing a more realistic and challenging evaluation platform for ACR systems.
2. Evaluation Method: We developed and validated an objective LLM-based evaluation method
to objectively assess ACR quality by comparing tool outputs against human ground truth.
3. Study: We systematically assessed mainstream LLMs and ACR tools on SWR-Bench, analyzing
their performance, strengths, and limitations in practical code review scenarios.
4. Approach: We developed and validated a simple multiple review strategy that significantly
improves the code review performance of ACR tools on SWR-Bench.
2
Background and Related Work
2.1
Modern Code Review
Modern Code Review (MCR) is a widely adopted quality assurance practice in contemporary
software development [2, 49], crucial for enhancing code quality, fostering knowledge sharing, and
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

FSE137:4
Zhengran Zeng, Ruikai Shi, Keke Han, Yixin Li, KC Sun, YD Wang, ZH Yu, Rui Xie, Wei Ye and Shikun Zhang
Pull Request
Change 1: Doc Improve of ...
Change 2: Bug of ...
Change 3: Bug of ...
① Reviewers review the code in
the Pull Request (PR).
Reviewers
Developer
② Reviewers provide a review report,
identifying the change-actions.
③ The developer verifies the change-
actions and updates the PR code.
Fig. 2. A simplified code review process.
ensuring adherence to project standards. It typically involves developers submitting code changes,
often through PRs, for scrutiny by their peers before integration into the main codebase.
A simplified, yet representative, MCR process is depicted in Figure 2. This process generally
unfolds in three key stages: 1) reviewers scrutinize the code submitted in the pull request; 2)
following their examination, reviewers provide a review report, identifying specific areas requiring
modification, which we term “change-actions”; 3) the developer then verifies these change-actions
and updates the PR code accordingly. The fundamental unit of actionable feedback within this
process is what we term a change-action. We define an effective change-action as a complete directive
that specifies the location(s) of the issue, describes the underlying problem, and proposes a tangible
solution. Therefore, a single change-action precisely reflects a specific issue and its resolution. An
ACR tool capable of generating correct change-actions can thus effectively guide developers in
improving code quality. It is also important to note that a change-action is not strictly tied to a
specific line of code; it can be conceptual and address a high-level concern spanning multiple code
locations, such as an architectural flaw. We use the term ‘change‘ rather than ‘fix‘ or ‘defect‘ because
reviewer intent often extends beyond rectifying bugs to include non-functional improvements in
readability, maintainability, or design, which are not traditional defects but are vital for software
quality [4, 11].
The nature of these change-actions varies significantly, reflecting diverse reviewer considerations.
To classify them, we adopt the detailed taxonomy from established literature shown in Table 1 [4, 11].
We selected this taxonomy because it is a reliable standard that has been consistently used across
multiple studies in the Automated Code Review (ACR) domain. It distinguishes between two
high-level categories: Evolutionary changes (improving future maintainability) and Functional
changes (altering software behavior), which are further decomposed into 11 fine-grained types.
Understanding this classification is fundamental for analyzing the effectiveness of both manual
and automated review processes.
Furthermore, it is important to note that in this work, we focus specifically on the task of
review comment generation, which we consider a critical prerequisite for end-to-end code review
automation. The quality of generated comments directly determines the success of downstream
tasks like automated code refinement, as precise feedback is the necessary input for any subsequent
modifications.
2.2
Automatic Code Review
While modern code review, as described above, is a vital software quality practice, manual code
review faces efficiency and consistency challenges [7, 10, 32], motivating the development of
ACR, aiming to use tools to assist or replace parts of manual review, thereby improving efficiency,
shortening feedback cycles, and enhancing consistency [9].
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

SWR-Bench: Assessing LLM Performance in Real-World Code Review Comment Generation
FSE137:5
Table 1. The taxonomy of change-actions under code review, adapted from [4, 11]. Note: All change-actions
are classified into leaf-node categories only.
ID
Type & Description
E
Evolutionary: Improving code’s future maintainability and structure, not its function.
E.1
Documentation: Modifying in-code information for better human understanding.
E.1.1
Textual: Changing textual code elements such as inline comments and variable/function names.
E.1.2
Language Supported: Using language features for documentation purposes.
E.2
Visual Representation: Modifying code layout for better style and readability.
E.3
Structure: Altering the project’s organization or architecture.
E.3.1
Organization: Reorganizing code by moving or removing parts.
E.3.2
Solution Approach: Altering implementation methods or adding maintainability elements (e.g., tests).
F
Functional: Altering the software’s behavior or interactions.
F.1
Interface: Modifying interactions between codebase parts.
F.2
Logic: Changing the code’s logical operations or algorithms.
F.3
Resource: Altering how variables or resources are managed.
F.4
Check: Modifying checks for unhandled states/conditions.
F.5
Support: Adjusting interactions with external support systems or libraries.
F.6
Larger Defects: PRs that were ultimately not merged (i.e., rejected suggestions), typically involving major,
extensive functional issues or missing features that required fundamental rethinking.
Specifically, ACR technology has progressed through several stages. Initially, development fo-
cused on rule-based and static analysis tools, such as PMD [28] and SonarQube [34], which are effec-
tive but can be rigid and prone to false positives. Subsequently, machine learning and deep learning
approaches [22, 39, 41, 42] emerged for tasks like code change quality assessment [15, 29, 36, 50, 51],
comment generation [22], and automated code repair [1, 19, 47]; while these methods can learn
complex patterns, they may still struggle with deep code intent and complex contexts. More recently,
large language model and AI agent-based methods [15, 21, 26, 30, 35, 36], have demonstrated strong
code understanding and generation capabilities, enabling more natural comments and complex
reasoning. Cutting-edge work in this area, such as CodeAgent [36], explores multi-agent systems
for collaborative review, while commercial tools like PR-Review [29] focus on providing targeted,
actionable feedback that dynamically learns team norms, aiming to address new challenges in
AI-assisted coding like feedback redundancy and unclear prioritization. Despite the potential of
ACR, especially LLM/Agent-based methods, their real-world effectiveness in improving review
quality and efficiency remains a key question. This highlights an urgent need for more effective and
reliable evaluation methods and benchmarks, as current assessment approaches [15, 22, 39, 41, 42]
may not adequately capture the performance of these advanced models in complex, realistic code
review scenarios.
2.3
Code Review Benchmarks
To evaluate various ACR techniques, the research community has constructed several specialized
code review benchmarks. Table 2 provides a comparative overview of representative benchmarks,
assessing them based on their data source, scale, fundamental review unit, the scope of provided
context, and evaluation methodology. This comparison highlights a significant gap in existing
resources, which our work, SWR-Bench, aims to fill.
Tufano et al.’s datasets [39, 41, 42], typically sourced from Gerrit and GitHub projects, center on
method-level triplets (i.e., the submitted method, reviewer comments, and the modified method).
These are used to evaluate tasks like method-level code transformation, implementing changes
based on comments, and generating comments for methods. As indicated in Table 2, these bench-
marks provide no external code context and primarily rely on text-similarity metrics like 𝐵𝑙𝑒𝑢for
evaluation.
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

FSE137:6
Zhengran Zeng, Ruikai Shi, Keke Han, Yixin Li, KC Sun, YD Wang, ZH Yu, Rui Xie, Wei Ye and Shikun Zhang
Table 2. Comparison of Representative Code Review Benchmarks.
Dataset
Source
Size
Unit
Context
Evaluation
Trans-Review-Data [42]
Github, Gerrit
1,719
Method
None
Bleu
AutoTransform-Data [39]
Gerrit
14,750
Method
None
Bleu
T5-Review-Data [41]
Github, Gerrit
17,194
Method
None
Bleu
Code-Reviewer-Data [22]
Github
10,000
Diff Hunk
None
Bleu, Human Eval.
CR-Agent-Data [36]
Github
3,545
Commit
Related Source Code
Human Eval.
Hybrid-Review-Data [15]
Github
1,245
Diff Hunk
Related Source Code
LLM Scoring
SWR-Bench
Github
1,000
Pull Request
Complete Codebase
Objective LLM Eval.
Code-Reviewer-Data [22] shifts the focus to a single diff hunk (i.e., a contiguous block of modified
lines). While closer to a review task, it still operates on fragments of a change rather than the whole.
Its evaluation metrics also include BLEU and Exact-Match, with some later studies [15] adopting
simple LLM-based scoring. More recent benchmarks like CR-Agent-Data [36] and Hybrid-Review-
Data [15] begin to incorporate related source code as context, but still fall short of providing the
complete project environment.
While valuable, these benchmarks have limitations compared to real-world code review. As
summarized in Table 2, they consistently lack a PR-level scope and comprehensive context, instead
focusing on isolated methods or diff hunks. This hinders evaluation of an ACR tool’s ability to
manage real-world complexity. Furthermore, these benchmarks predominantly rely on inadequate
evaluation metrics, such as 𝐸𝑥𝑎𝑐𝑡-𝑀𝑎𝑡𝑐ℎand 𝐵𝑙𝑒𝑢or simple LLM-based ratings metric, which are
known to correlate poorly with human judgment [8, 20], thus leading to potentially misleading
assessments of tool efficacy.
Beyond these public benchmarks, some studies [9, 35] leverage proprietary internal data for ACR
research, validating tools through manual human assessment. While such efforts provide valuable,
deep insights, the private nature of this data and the prohibitive cost of manual validation make
these approaches unscalable and difficult for the broader community to build upon. This leaves a
clear void for a publicly available, context-rich, and scalably evaluable benchmark.
In summary, as Table 2 illustrates, existing public benchmarks consistently lack PR-level scope,
comprehensive project context, and reliable evaluation metrics, three deficiencies that limit their
ability to accurately assess modern LLM-based ACR tools. Meanwhile, proprietary benchmarks with
manual evaluation, though insightful, are neither scalable nor reproducible. These gaps motivate
the design of SWR-Bench, which we detail in the following section.
3
SWR-Bench
This section details our methodology for constructing the automatic code review benchmark.
3.1
Benchmark Construction
The overall construction workflow of SWR-Bench is depicted in Figure 3 and detailed below.
Step 1: Source Data Collection and Initial Filtering. The foundation of SWR-Bench lies in
real-world software development practices. To ensure the benchmark’s realism and quality, we
adopted the 12 open-source Python projects from the well-established SWE-Bench [18], as they
represent the most popular packages on PyPI by download count. This selection guarantees their
code is of high quality and their development practices are highly representative of real-world
scenarios. We utilized the GitHub API to crawl all historical pull requests (PRs) from these projects,
collecting comprehensive metadata for each. This included titles, descriptions, commit histories,
code diffs, review comments, discussion threads, and final PR statuses (e.g., merged, closed).
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

SWR-Bench: Assessing LLM Performance in Real-World Code Review Comment Generation
FSE137:7
Filter & Sampling

 Filtering out Change-PRs missed
by human reviewers using the SZZ
algorithm.
 Resample Clean-PRs to align
their statistical distribution with
Change-PRs.
Scrape PRs
 LLM Verify

 Automatically extracting Change-
Points from PRs with LLMs.
 Enhancing annotation quality
through majority voting.
Human Verify
 Manually verify the correctness of
PR Change-Points.
21,350 PRs For LLM Verify
3,500 PRs For Human Verify
Data Source
12 Popular Repositories
From GitHub
Scrape PR Metadata
Title
Description
Commits
Comments
Status
Refinement Filter
✘PRs without review comments
✘PRs with large-scale changes
✘PRs with rebase operations
SWR-Bench
500
Clean-PRs
500
Change-PRs
PR Metadata
Change-Points
Codebase Checkpoint
Objective LLM as Judge
Fig. 3. The SWR-Bench construction pipeline.
From this raw collection, we performed initial filtering to refine the candidate pool. We excluded
PRs without any review comments, as they lack the necessary interaction for our analysis. We also
removed PRs with large-scale changes (e.g., exceeding 10 commits in one PR) or those with rebase
operations in their commit history. This is because extensive changes make it challenging for the
LLM in the subsequent verification stage to correctly identify and label change-actions, and rebase
operations complicate the construction of reproducible environments. After this initial filtering,
21,350 PRs remained for the next phase.
Step 2: LLM-based change-actions Verification and Classification. In this step, we employed
a LLM to identify and extract “change-actions” from the filtered PRs. Specifically, we provided the
LLM with the complete, chronologically ordered review timeline for each PR, which included its
title, description, all developer-reviewer discussions, and all commit messages. A crucial instruction
was for the LLM to identify each change-action (corresponding to the 11 types in Table 1) and link
it to the specific commit that introduced and fixed the issue, with the full prompt is detailed in [1]
due to space constraints. To ensure high-quality annotations, we employed Gemini-2.5-Pro [12], a
state-of-the-art (SOTA) LLM. Furthermore, we implemented a majority voting technique by making
three independent requests to the LLM for each PR. A PR was considered for further processing
only if all three requests yielded consistent change-actions extraction results. This requirement acts
as a quality filter, excluding PRs where the LLM produced inconsistent results across runs. Lastly,
those consistent PRs were then classified into two categories: Change-PRs”, which contain at least
one identified change-action, and Clean-PRs”, which have no identified change-action.
Step 3: Quality Enhancement through Filtering and Sampling. To further enhance dataset
quality, we performed additional filtering and sampling. First, we applied the SZZ algorithm [31]
to identify issues that might have been missed by human reviewers but were fixed in later commits.
PRs containing such missed changes were removed. This step ensures that our Clean-PRs are
genuinely free of known, non-trivial defects. Consequently, any comment an ACR tool generates
for a Clean-PR is, by definition, considered a false positive. This strict protocol enables a robust
measurement of the false positive rate, a critical metric for assessing a tool’s precision and its
ability to avoid distracting developers with irrelevant suggestions. We acknowledge that SZZ has
known limitations in precision and recall [31], as it may miss a small number of latent bugs or
produce false links. We therefore employed it solely as a supplementary heuristic to improve dataset
completeness, not as a sole ground-truth source, and all results were subject to subsequent manual
verification (Step 4). Moreover, since all tools are evaluated on the same SZZ-augmented dataset,
this still provides a fair basis for comparative analysis.
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

FSE137:8
Zhengran Zeng, Ruikai Shi, Keke Han, Yixin Li, KC Sun, YD Wang, ZH Yu, Rui Xie, Wei Ye and Shikun Zhang
Next, to prevent models from relying on superficial heuristics, we addressed the statistical
differences between Change-PRs and Clean-PRs. As highlighted in prior work [51], models can
exploit simple metrics (e.g., lines of code changed, number of commits) to distinguish between PRs
needing changes and those that do not, without truly understanding the code. To mitigate this and
create a more robust evaluation, we performed stratified sampling [46] on the Clean-PRs to align
their distribution of key statistical properties with that of the Change-PRs. While this intentionally
deviates from the natural data distribution, it creates a more challenging and meaningful benchmark
that forces models to engage with code logic.
Step 4: Manual Verification and Refinement. After the preceding steps, approximately 3,500
PRs remained. Given the significant effort required for manual annotation, we randomly sampled
1,000 Change-PRs and 1,000 Clean-PRs for rigorous manual verification. This process was conducted
by five experienced graduate students in computer science, each holding at least a Master’s degree
and possessing over two years of software development experience. Specifically, we assigned 800
PRs to each annotator, organized such that each PR was independently annotated by exactly two
annotators. The primary goal of this manual verification was to check the correctness of each
LLM-identified change-action label, ensuring it conformed to our definition and could be accurately
classified into one of the 11 types in Table 1. Additionally, annotators filtered out Change-PRs
containing only “trivial change-actions”, which we define as changes related solely to formatting
or documentation updates (types E.1 and E.2 in Table 1). This was a deliberate choice to create a
more challenging and realistic benchmark, as a model excelling on SWR-Bench by being forced to
address substantive issues is more likely to be effective in real-world scenarios.
To ensure annotation reliability, we calculated Cohen’s Kappa [45] for each pair of independent
reviews, yielding a score of 66.08. While this indicates “substantial agreement,” it also reflects the
inherent difficulty and ambiguity of identifying specific change-actions from PR discussions. To
mitigate the impact of this ambiguity on dataset quality, all disagreements were resolved through
discussion between the two assigned annotators to ensure the correctness of the final labels. This
rigorous verification process was crucial for ensuring the high quality and reliability of SWR-Bench.
Dataset Statistics. The final SWR-Bench dataset consists of 500 Change-PRs and 500 Clean-PRs,
randomly selected from manually verified and corrected PRs. This balanced composition is vital for
robustly evaluating ACR tools, particularly their false positive rate, a common critique of existing
automated review systems [35, 37].
Each SWR-Bench instance offers a comprehensive dataset for the evaluation of ACR tools, cap-
turing the state of a PR right before its first manual review. This dataset is composed of three
key components: (1) the metadata data of the PR (title, description, commits made before the first
manual review), (2) ground-truth change-action (only for those made prior to the first manual
review) information for Change-PRs, detailing their type, description, and associated commit SHA,
and (3) codebase checkpoint information, allowing for the reconstruction of the full PR codebase at
the relevant commit. This latter feature is particularly important for agent-based ACR tools that
require an interactive environment.
Table 3 presents the overall dataset statistics. We find that, on average, Change-PRs have similar
numbers of modified files and lines of code compared to Clean-PRs. These similarities underscore
the effectiveness of our resampling strategy (Step 3) in ensuring that Clean-PRs are not trivially
distinguishable from Change-PRs based on such coarse-grained metrics.
Furthermore, we analyzed the change-action type distribution (Figure 4). In Raw-PRs (before
manual verification and filtering), functional changes (e.g., logic errors, performance issues) consti-
tuted less than 15% of all identified changes, with a predominance of more evolutionary changes
like documentation or style adjustments. After our manual verification process, which explicitly
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

SWR-Bench: Assessing LLM Performance in Real-World Code Review Comment Generation
FSE137:9
filtered out PRs with only trivial change-actions, the proportion of functional changes in SWR-Bench
increased to 31.8%, while the proportion of simple textual changes was intentionally reduced. Addi-
tionally, earlier quality filtering steps (e.g., majority voting) also contribute to distribution shifts
from real-world data. We argue that such deviations are acceptable, as our goal is to provide a fair
and consistent evaluation scenario for ACR tools. A method that performs well on this intentionally
difficult distribution, where simple cases are underrepresented, is highly likely to generalize well to
practical scenarios.
Table 3. Overall statistics of the SWR-Bench. “Avg.”
denotes average values of all instances.
PR Type
Count
Avg.
Commits
Avg.
Files
Avg.
Lines +
Avg.
Lines -
Avg. Change
Points
Change-PR
500
4.05
6.29
123.64
60.69
1.90
Clean-PR
500
3.25
6.85
112.60
67.22
0.00
All
1000
3.65
6.57
118.12
63.96
0.95
22.3%
2.2%
7.3%
10.0%
26.4%
4.6%
10.3%
2.8%
3.7%
7.5%
2.9%
(a) SWR-Bench Distribution
42.4%
1.7%
10.0%
8.3%
23.7%
2.2%
4.4%
1.1%
1.6%
3.2%
1.2%
(b) Raw-PRs Distribution
E.1.1 Textual Changes
E.1.2 Language Features
E.2 Visual Representation
E.3.1 Organization
E.3.2 Solution
F.1 Interface
F.2 Logic
F.3 Resource
F.4 Check
F.5 Support
F.6 Larger Defects
Fig. 4. Change type distribution of SWR-Bench.
3.2
Evaluation Methodology
Evaluating the output of ACR tools presents unique challenges. The report generated by a model is
typically unstructured natural language text, containing various comments, or identified issues.
Directly comparing this with a structured list of ground truth change-actions is difficult. As discussed
(Section 2), traditional text similarity metrics fail to capture semantic accuracy and practical value in
code review comments [17, 25]. Consequently, we turn to solutions utilizing large language models
as evaluation aids (LLM-as-Judge) [8, 44]. However, we also recognize that traditional LLM-as-Judge
methods often rely on the LLM’s subjective judgment (e.g., scoring comment quality), which can
introduce bias and inconsistency [13, 33, 38, 52]. To mitigate this, we designed an objective LLM
evaluation framework that leverages ground-truth change-actions referencing. Instead of subjective
scoring or ranking, our LLM performs a fact-based matching task: determining if issues in the
model’s report correspond to predefined ground truth change-actions.
Evaluation
Report

Precision
Recall
F1
Accuracy
. . .
Ground-Truth
Change-Actions
Predicted
Review Report
LLM
Evaluation Output

 Pred Change-Action 1
 Change-Type: E.1

 Pred Change-Action 2
 Change-Type: F.2
✘GT Change-Action 1
 Hit By: None
✘GT Change-Action 2
 Hit By: None

✔GT Change-Action 3
 Hit By: Pred Change-Action 2
Fig. 5. The objective LLM-based evaluation pipeline for predicted review report.
Our evaluation process, depicted in Figure 5, unfolds as follows:
The inputs to the evaluation pipeline are the code review report generated by the ACR tool and
the corresponding ground truth change-actions from our SWR-Bench. First, an evaluation LLM
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

FSE137:10 Zhengran Zeng, Ruikai Shi, Keke Han, Yixin Li, KC Sun, YD Wang, ZH Yu, Rui Xie, Wei Ye and Shikun Zhang
parses the ACR tool’s unstructured report to extract distinct predicted change-actions (a single
report generated by the ACT tool contains multiple change-actions). For each predicted change-
action, the LLM is also prompted to determine its change-type (from 11 types in Table 1). These
latter two are subjective LLM assessments but provide useful metadata about the review report.
Next, the core objective evaluation involves the LLM performing a matching task. It compares
each extracted predicted change-action against the PR’s set of ground truth change-actions. The
LLM determines, for each ground truth change-action, whether it has been successfully “hit” (i.e.,
semantically identified) by one or more predicted change-actions.
Based on this explicit matching information, we can define TP, FP, and FN. These definitions then
enable the straightforward calculation of standard performance metrics like 𝑃𝑟𝑒𝑐𝑖𝑠𝑖𝑜𝑛=
𝑇𝑃
𝑇𝑃+𝐹𝑃,
𝑅𝑒𝑐𝑎𝑙𝑙=
𝑇𝑃
𝑇𝑃+𝐹𝑁, and 𝐹1 = 2 × 𝑃𝑟𝑒𝑐𝑖𝑠𝑖𝑜𝑛×𝑅𝑒𝑐𝑎𝑙𝑙
𝑃𝑟𝑒𝑐𝑖𝑠𝑖𝑜𝑛+𝑅𝑒𝑐𝑎𝑙𝑙:
• True Positives (TP): Ground-truth change-actions successfully hit by at least one predicted
change-action.
• False Positives (FP): The number of predicted change-actions that do not hit any ground truth
change-action.
• False Negatives (FN): The number of ground truth change-actions not hit by any predicted
change-actions.
Furthermore, as discussed in Section 3.1, functional changes constitute a critical, but smaller
portion of the change-actions in SWR-Bench. To provide a more nuanced understanding of an ACR
tool’s ability to detect these functional issues, we report performance metrics (𝑃𝑟𝑒𝑐𝑖𝑠𝑖𝑜𝑛, 𝑅𝑒𝑐𝑎𝑙𝑙, and
𝐹1) specifically for functional changes in addition to the overall metrics. This specialized calculation
considers only predicted change-actions that the evaluation LLM has typed as functional change
and ground-truth change-actions with functional type.
This objective, ground-truth-referenced LLM evaluation approach allows for a more accurate
and in-depth assessment of ACR tools, overcoming traditional method limitations. The reliability
of this evaluation approach is further validated through manual verification in later experiments.
4
Evaluation and Study
4.1
Research Questions
Our study is guided by the following research questions:
• RQ1: How reliable is the proposed evaluation methodology?
• RQ2: How do mainstream automated code review tools and LLMs perform on SWR-Bench?
• RQ3: How can the performance of automated code review tools be improved?
4.2
Subjects of Study
4.2.1
Large Language Models (LLMs). To comprehensively evaluate current LLMs on ACR tasks,
we selected a diverse range of models, including closed-source LLMs like the GPT series (o3,
o4-mini, GPT-4o, GPT-5), Gemini series (Gemini-2.5-Pro, Gemini-2.5-Flash), and Claude series
(Claude-3.7-Sonnet, Claude-4-Sonnet, Claude-4-Opus). We also include open-source LLMs such
as DeepSeek series (DeepSeek-R1, DeepSeek-V3) and Qwen2.5 series (Qwen2.5-Chat-7B/14B/32B,
Qwen2.5-R1-7B/14B/32B).
4.2.2
Automated Code Review (ACR) Baselines. For a comprehensive evaluation, we benchmark
against representative ACR tools and methodologies:
• LLM-Review (Prompting Baseline): Our straightforward baseline that directly feeds PR metadata
and diffs to an LLM via a simple prompt, measuring vanilla review capabilities.
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

SWR-Bench: Assessing LLM Performance in Real-World Code Review Comment Generation
FSE137:11
• SWR-Agent (Agentic Baseline): Inspired by SWE-Agent [48], we built this baseline to adapt
agentic workflows for code review, allowing the LLM to use tools, explore the codebase, and
gather context.
• CR-Agent [36]: A multi-agent tool using two interacting agents (focusing on formatting and
functional defects, respectively) that discuss and debate to produce a final report.
• Hybrid-Review [16]: Extends the simple prompting baseline by additionally injecting static
analysis reports as supplementary context.
• PR-Review [29]: While similar to LLM-Review, PR-Review uses more sophisticated prompt
engineering. Specifically, it first consolidates code diffs from multiple commits to form the input.
When faced with context length limitations, it prioritizes files based on importance, potentially
excluding less critical files. Finally, it instructs the LLM to output the review report in a specific
structured format, categorizing suggestions into areas such as test-related defects, security-related
defects, and other areas for improvement.
• Code-Reviewer [22]: A traditional hunk-level method using a fine-tuned CodeT5 model [43].
We adapt our PRs into diff hunks for its input and concatenate its generated hunk-level comments
to form a full PR review.
• Llama-Reviewer [26]: Fine-tunes a Llama-Base-7B model [40] for hunk-level review. We evaluate
it using the same hunk-to-PR concatenation protocol as Code-Reviewer.
Lastly, to ensure fair comparison, all evaluated ACR tools used their official default configurations.
RQ1: How reliable is the proposed evaluation methodology?
G-2.5-Pro
G-2.5-Flash
Flash-MV3
Human-1
Human-2
G-2.5-Flash
Flash-MV3
Human-1
Human-2
Human-3
86.7
92.785.7
83.785.182.3
81.874.183.580.2
72.470.677.173.774.6
(a) Hit Agreement
G-2.5-Pro
G-2.5-Flash
Flash-MV3
Human-1
Human-2
G-2.5-Flash
Flash-MV3
Human-1
Human-2
Human-3
76.8
83.881.2
87.676.079.0
63.464.159.368.9
82.374.281.175.960.7
(b) Type Agreement
Human-1
Human-2
Human-3
LLM-Hit
LLM-Score
Human-2
Human-3
LLM-Hit
LLM-Score
BLEU
62.6
61.656.3
62.060.652.8
39.444.931.050.9
-35 -35 -41 -41 -55
(c) Agreement Comparison of
LLM-Hit, LLM-Score and BLEU
0
20
40
60
80
100
Cohen's Kappa Score (%)
Fig. 6. Validation of the proposed objective evaluation methodology. (a) and (b) show the inter-rater agreement
(Cohen’s Kappa) for Hit identification and Type classification, respectively. (c) compares the agreement
(Cohen’s Kappa) of two LLM-as-Judge methodology and BLEU judge methodology with human preferences
in a pairwise comparison task.
To rigorously assess the reliability of our proposed LLM-as-Judge methodology, we conducted
two validation experiments.
First, we measured inter-rater and inter-model agreement for our core evaluation tasks. We
randomly selected 100 code review reports from RQ2 and had them independently annotated
by three human experts and three distinct LLM judges (Gemini-2.5-Pro, Gemini-2.5-Flash and
Gemini-2.5-Flash with majority voting (𝑛=3)). We measured agreement using Cohen’s Kappa for
identifying “Hit” of each ground-truth change-actions and classifying their “Type”.
Hit Agreement (Figure 6a): The agreement for our primary "Hit" metric was exceptionally high,
with Kappa scores ranging from 70.6 to 86.7 across all human-human, human-LLM, and LLM-
LLM pairs. This indicates almost perfect agreement and confirms that identifying "Hits" is an
objective task. The high consistency between Gemini-2.5-Pro and Gemini-2.5-Flash (Kappa = 86.7)
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

FSE137:12 Zhengran Zeng, Ruikai Shi, Keke Han, Yixin Li, KC Sun, YD Wang, ZH Yu, Rui Xie, Wei Ye and Shikun Zhang
demonstrates that the evaluation is stable across different models, implying that the results are
not sensitive to the choice of a specific LLM judge, and therefore are unlikely to be sensitive to
the randomness of multiple runs. Furthermore, the majority voting variant achieved Kappa scores
comparable to the single-call Gemini-2.5-Flash and Gemini-2.5-Pro, indicating that majority voting
provides only marginal improvement for this already highly consistent task.
Type Agreement (Figure 6b): For "Type" classification, human-LLM kappa scores ranged from 63.4
to 87.6. The agreement between LLMs and humans was comparable to the agreement among humans
themselves (human-human Kappa: 60.7 to 75.9). This suggests that the LLM judge’s performance
on Type classification is consistent with that of human experts. Similarly, the majority voting
variant of Gemini-2.5-Flash showed no significant improvement over the single-call version for
Type classification, further confirming that a single LLM call is sufficient for reliable evaluation.
Second, to demonstrate the superiority of our objective methodology over traditional approaches,
we compared it against several alternative evaluation methods. We randomly selected 1,000 pairs
of review reports from RQ2, where each pair containing two distinct reviews for the same PR,
and asked three human experts to choose the better report in each pair. We then tasked the
following automated methods (both using Gemini-2.5-Flash) to perform the same comparison: 1)
LLM-Hit-Judge (ours), our proposed method, which prefers the report that achieves more “Hits” on
ground-truth change-actions with fewer false positives. 2) LLM-Score-Judge, a traditional approach
where the LLM assigns a holistic quality score (1-10) to each report, with the higher-scoring report
being preferred. 3) Bleu [27] and Rouge-L [23], traditional text similarity metrics that compare the
generated review report against the ground-truth review discussion text, with the report achieving
a higher score being preferred.
As shown in Figure 6c, LLM-Hit-Judge achieved substantial agreement with human preferences,
with Cohen’s Kappa scores ranging from 52.8 to 62.0. This level of agreement is on par with the inter-
human agreement (Kappa range: 56.3 to 62.6). In contrast, the baseline LLM-Score-Judge showed
significantly lower agreement with humans (Kappa range: 31.0 to 44.9), highlighting the unreliability
of subjective scoring. Most notably, the traditional text similarity metrics performed drastically
worse. BLEU achieved Kappa scores ranging from only -35 to -41 with human preferences, indicating
agreement even worse than chance. This result validates that our objective, change-action-based
evaluation aligns more closely with expert human judgment.
Ultimately, Gemini-2.5-Pro, Gemini-2.5-Flash and Majority-Voting demonstrated commendable
and similar agreement levels with human evaluators, particularly for the primary “Hit” metric.
Given that Gemini-2.5-Flash demonstrated strong, reliable performance comparable to Gemini-2.5-
Pro at a significantly lower cost (approx. $1.57 for a full SWR-Bench evaluation), we selected it for
all subsequent experiments.
Conclusion 1: Our objective, change-action-based evaluation methodology is highly reliable,
stable across different LLM judges (with or without majority voting), and demonstrates superior
alignment with human expert judgment compared to both traditional subjective scoring methods
and text similarity metrics.
RQ2: How do mainstream automated code review tools and LLMs perform on
SWR-Bench?
In this RQ, we first investigate the performance of studied ACR approaches on SWR-Bench, and
each approaches was evaluated using three powerful LLMs, with the exception of SWR-Agent,
which was not run with DeepSeek-R1 due to the latter’s lack of function call capabilities. For each
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

SWR-Bench: Assessing LLM Performance in Real-World Code Review Comment Generation
FSE137:13
Table 4. Evaluation of studied ACR approaches on SWR-Bench. The table shows hit-based Precision, Recall,
and F1. Avg. Count indicate the average number of predicted change-actions, while Avg. FP Count indicate the
average false positive number.
ACR Tools
LLMs
Overall Change-Actions
Functional Change-Actions
Text Similarity
Precision
Recall
F1
Avg.
Count
Avg. FP
Count
Precision
Recall
F1
Avg.
Count
Avg. FP
Count
Bleu
Rouge-L
LLM-Review
Gemini-2.5-Pro
8.02
21.29
11.65
2.47
2.28
22.60
26.86
24.54
0.35
0.27
0.58
9.64
Claude-3.7-Sonnet
9.85
14.31
11.67
1.35
1.22
24.37
16.02
19.33
0.20
0.15
0.47
10.53
DeepSeek-R1
9.79
25.58
14.16
2.44
2.20
15.53
32.02
20.92
0.61
0.52
0.50
10.19
Mean
9.22
20.39
12.49
2.09
1.90
20.83
24.97
21.60
0.39
0.31
0.52
10.12
SWR-Agent
Gemini-2.5-Pro
9.93
19.14
13.07
1.80
1.62
18.11
25.84
21.30
0.42
0.35
1.95
11.45
Claude-3.7-Sonnet
8.29
22.72
12.15
2.55
2.34
17.92
30.90
22.68
0.51
0.42
1.46
11.77
Mean
9.11
20.93
12.61
2.18
1.98
18.02
28.37
21.99
0.47
0.38
1.31
11.11
CR-Agent
Gemini-2.5-Pro
6.97
17.63
9.98
2.35
2.18
16.26
26.55
20.17
0.48
0.40
0.74
7.43
Claude-3.7-Sonnet
6.30
17.17
9.21
2.54
2.38
18.06
22.78
20.15
0.38
0.31
1.00
10.22
DeepSeek-R1
5.42
19.50
8.48
3.36
3.17
11.45
29.44
16.49
0.77
0.68
1.00
9.87
Mean
6.23
18.10
9.22
2.75
2.58
15.26
26.26
18.94
0.54
0.47
0.91
9.17
Hybrid-Review
Gemini-2.5-Pro
2.83
20.44
4.97
6.64
6.42
11.09
29.28
16.08
0.80
0.71
0.77
7.28
Claude-3.7-Sonnet
2.20
11.23
3.68
4.63
4.52
7.95
15.64
10.55
0.59
0.54
0.39
5.30
DeepSeek-R1
3.33
28.44
5.96
7.95
7.60
10.10
47.31
16.65
1.47
1.31
1.02
8.41
Mean
2.79
20.04
4.87
6.41
6.18
9.71
30.74
14.43
0.95
0.85
0.73
7.00
PR-Review
Gemini-2.5-Pro
16.65
23.18
19.38
1.32
1.10
19.38
40.72
26.26
0.65
0.52
0.40
8.92
Claude-3.7-Sonnet
14.90
23.50
18.23
1.50
1.27
14.72
40.32
21.56
0.86
0.74
0.36
8.38
DeepSeek-R1
14.61
25.50
18.58
1.66
1.41
15.55
50.62
23.79
1.06
0.89
0.31
7.70
Mean
15.39
24.06
18.73
1.49
1.26
16.55
43.89
23.87
0.86
0.72
0.36
8.34
Code-Reviewer
4.19
11.35
6.13
2.55
2.44
9.46
14.67
11.50
0.47
0.42
10.25
22.14
Llama-Reviewer
4.28
23.08
7.22
4.91
4.70
8.62
37.02
13.98
1.32
1.20
7.38
22.89
ACR technique, we also report the mean performance metrics across the LLMs it was paired with.
The detailed results are summarized in Table 4.
The results in Table 4 show that SOTA ACR techniques, when paired with SOTA LLMs, are not
yet ready for real-world code review deployment based on their performance on SWR-Bench. For
instance, the top-performing combination, PR-Review leveraged with Gemini-2.5-Pro, achieved
an 𝐹1 score of only 19.38%. A primary factor limiting higher 𝐹1 scores for all techniques is their
low precision, indicative of a high false positive rate. Specifically, the Avg. FP Count metric reveals
that these approaches frequently generate multiple invalid suggestions per pull request, with some
combinations (e.g., Hybrid-Review with DeepSeek-R1) producing over 7 false positives on average.
This issue is more severe for other four ACR techniques, all of which exhibited precision scores
below 10%. This implies that while these ACR techniques can identify some valuable issues (or
“change-actions”), their overall effectiveness is severely undermined by an excessive number of false
positives. Consequently, developers would need to invest considerable additional effort in verifying
the validity of the generated reports, hindering practical adoption. This observation aligns with
findings from previous empirical studies [37] and further validates the efficacy of SWR-Bench as a
benchmark that closely mirrors realistic code review scenarios.
Conclusion 2: Current ACR approaches, even when augmented with advanced LLMs, demon-
strate limited performance on the SWR-Bench, primarily constrained by high false positive rates.
This significantly hinders their immediate applicability in practical code review workflows.
Performance Analysis by ACR Tools. A granular analysis of 𝐹1 scores reveals distinct performance
tiers among the ACR tools (Table 4). PR-Review, leveraging meticulous prompt engineering, achieved
the highest overall performance (𝑂𝑣𝑒𝑟𝑎𝑙𝑙-𝐹1: 18.73%). This significantly surpassed both SWR-Agent
and LLM-Review (approx. 12%).
Meanwhile, other approaches showed significant limitations. The multi-agent CR-Agent per-
formed poorly (𝑂𝑣𝑒𝑟𝑎𝑙𝑙-𝐹1: 9.22%), likely due to interaction overhead and error propagation—known
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

FSE137:14 Zhengran Zeng, Ruikai Shi, Keke Han, Yixin Li, KC Sun, YD Wang, ZH Yu, Rui Xie, Wei Ye and Shikun Zhang
challenges in multi-agent systems [5]. The worst performer, Hybrid-Review (𝑂𝑣𝑒𝑟𝑎𝑙𝑙-𝐹1: 4.87%),
suffered from extremely low precision (with the highest mean Avg. FP Count of 6.18), indicating
that simply integrating raw static analysis outputs is ineffective without careful processing.
Furthermore, the poor performance of traditional fine-tuned models (Code-Reviewer: 6.13% 𝐹1,
Llama-Reviewer: 7.22% 𝐹1) is also notable. It highlights a critical flaw in prior approaches: models
optimized for isolated code hunks fail at the holistic, contextual task of reviewing a full pull request.
This underscores the value of SWR-Bench in providing a more realistic evaluation and guiding
research toward end-to-end solutions.
Table 4 also reports Bleu [27] and Rouge-L [23] scores. Interestingly, text similarity inversely
correlates with 𝐹1: hunk-level models achieve high Bleu but low 𝐹1, whereas PR-level tools (e.g.,
PR-Review) show the opposite. Our manual analysis reveals this stems from formatting differences
rather than review quality. PR-Review generates structured, detailed reports that differ from
brief, colloquial human ground truths, yielding low Bleu despite capturing actual issues (high 𝐹1).
Conversely, hunk-level models mimic human chat, often generating short, ambiguous comments
that inflate Bleu but fail to identify actionable defects (low 𝐹1). Consistent with RQ1 (Figure 6c),
this format bias confirms that Bleu and Rouge-L are unreliable proxies for review quality, validating
our change-action hit-based evaluation.
Lastly, comparing PR-Review and LLM-Review, both of which operate via a single-turn interaction
with the LLM, PR-Review’s superior performance can be attributed to its more sophisticated prompt
engineering, these refinements appear to significantly reduce false positives, thereby enhancing
precision and 𝐹1 score.
Collectively, these results underscore that refined prompt engineering, as demonstrated by
PR-Review, is currently the most effective strategy for improving ACR capabilities. Agent-based
approaches require further research to overcome their architectural challenges and unlock their
potential. The failure of models trained on isolated code hunks validates the necessity of holistic,
end-to-end benchmarks like SWR-Bench to guide the field towards more practical solutions.
Conclusion 3: PR-Review, leveraging meticulous prompt engineering, achieves the best perfor-
mance on SWR-Bench, underscoring the promise of prompt engineering for improving efficacy.
Table 5. PR-Review (Gemini-2.5-Pro) perfor-
mance vs. number of ground-truth changes
(𝑁).
𝑁
Count
Precision
Recall
F1
1
266
30.18
38.35
33.77
2
139
35.05
24.46
28.81
3
56
34.18
16.07
21.86
4
17
29.63
11.76
16.84
5+
22
44.12
8.88
14.78
Table 6. Performance across different change types on
SWR-Bench.
Change Type
Precision
Recall
F1
Avg.
Count
E.1.1 Textual Changes
15.85
13.03
14.30
0.17
E.1.2 Language Features
8.28
7.59
7.85
0.02
E.2 Visual Representation
12.21
4.39
6.05
0.02
E.3.1 Organization
24.94
12.56
16.45
0.05
E.3.2 Solution Approach
12.66
19.15
15.21
0.37
F.1 Interface
16.97
40.83
23.55
0.12
F.2 Logic
17.28
54.60
26.20
0.35
F.3 Resource
15.89
53.45
24.26
0.10
F.4 Check
13.36
39.70
19.60
0.12
F.5 Support
15.95
35.80
21.74
0.16
F.6 Larger Defects
53.31
18.88
27.65
0.01
Performance in Multi-Change PR Scenarios. A distinctive feature of SWR-Bench is that each
Change-PR contains multiple ground-truth change-actions (on average 1.90 per PR, as shown in
Table 3). To understand how baselines perform in such scenarios, we analyzed the relationship
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

SWR-Bench: Assessing LLM Performance in Real-World Code Review Comment Generation
FSE137:15
between the number of ground-truth change-actions (𝑁) in a PR and the ACR tool’s performance. As
detailed in Table 5 (using PR-Review with Gemini-2.5-Pro as a representative example), we observed
that as the number of change-actions per PR increases, the Recall tends to decrease sharply (from
38.35% for 𝑁= 1 to 8.88% for 𝑁≥5), indicating that tools struggle to identify all issues in more
complex PRs. Conversely, Precision remains relatively stable (fluctuating roughly between 29.63%
and 44.12%), suggesting that the quality of individual predictions does not degrade significantly
with PR complexity. This analysis highlights that multi-change PRs remain a significant challenge.
Conclusion 4: ACR tools struggle to comprehensively review complex PRs. As the number of
issues in a PR increases, ACR tools experience a sharp drop in Recall while Precision remains
stable, highlighting their inability to comprehensively review multi-change PRs.
Performance Analysis by Change Type. To gain a deeper understanding of how ACR tools handle
different types of change-actions, we further analyzed the average performance of PR-Review across
various change types defined in Table 1. The results, disaggregated by change type, are shown
in Table 6. A key observation is that PR-Review demonstrates significantly stronger detection
capabilities for functional change-actions compared to evolutionary change-actions. Specifically,
for the F.2 Logic change, PR-Review achieved an 𝐹1 score of 26.20%, with most functional change
types also yielding 𝐹1 scores above 21%. In stark contrast, the highest 𝐹1 score for an evolutionary
change type, E.3.1 Organization, was merely 16.45%.
We hypothesize that this disparity arises from the inherent nature of these change categories.
Evolutionary changes often represent optional or stylistic improvements, where the criteria for
what constitutes a necessary change can vary significantly between human reviewers [4] and
LLMs. For instance, with E.3.2 Solution Approach, which was the most prevalent evolutionary
change type, one human reviewer might suggest an alternative implementation they deem superior,
while another might find the current approach acceptable. Consequently, accurately detecting such
optional evolutionary changes presents a considerable challenge.
Based on these findings, we recommend that future development of ACR tools prioritize en-
hancing the detection capabilities for functional changes to improve the precision for functional
change-actions. Furthermore, ACR tools could consider presenting reports for functional and evolu-
tionary changes separately.
Conclusion 5: ACR tools demonstrate better performance in detecting functional changes,
which likely due to the subjective nature of many evolutionary changes. Accordingly, future
ACR tools could prioritize robust detection of functional changes to enhance practical utility.
False Positive Analysis. While our previous analysis highlighted the practical utility of ACR tools
is critically undermined by the generation of false positives, we therefore performed a qualitative
analysis to diagnose the root causes of these inaccuracies. We manually inspected 100 functional
false positives, randomly sampled from a total of 1,101 false positives generated by PR-Review
(based on Gemini-2.5-Pro), categorizing them as summarized in Table 7.
The taxonomy was derived through an iterative process based on Grounded Theory [5]. Three
authors independently examined and labeled each of the 100 sampled false positives, assigning
initial descriptive codes to characterize the root cause. The authors then collaboratively discussed
and reconciled their labels through multiple rounds of discussion, iteratively merging, splitting,
and refining categories until a stable set of themes emerged. This process resulted in the five major
categories and one “Other” category presented in Table 7.
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

FSE137:16 Zhengran Zeng, Ruikai Shi, Keke Han, Yixin Li, KC Sun, YD Wang, ZH Yu, Rui Xie, Wei Ye and Shikun Zhang
Table 7. Categories and Proportions of Common False Positives
Categories
Description
Proportion
Lack of Contextual
Understanding
The tool applies programming rules rigidly and in isolation, ignoring
crucial logical consistency established by the surrounding code within the
same project or commit.
48%
Over-sensitivity to
Modification
The tool lacks an understanding of the developer’s intent and tends to
treat any large-scale code modification or refactoring as a potential risk.
17%
Vague and
Unactionable Feedback
The review feedback provided by the tool is often too broad and generic,
lacking specific, actionable steps.
16%
Lack of Domain
Knowledge
The tool fails to recognize specialized coding conventions specific to a
project or technical domain (i.e., domain knowledge), causing it to flag
correct and idiomatic code as anomalous.
13%
Misjudgment of
"Anti-best-practices"
The tool relies on superficial heuristics and makes incorrect judgments
when developers intentionally deviate from conventional best practices to
achieve higher-level goals like test effectiveness or readability.
3%
Other
A catch-all for other rare causes.
3%
As the table shows, the most prevalent issue is a Lack of Contextual Understanding (48%). These
errors typically arise from the LLM’s insufficient contextual reasoning, leading to a misunder-
standing of the code. For example, in astropy/astropy-1199, the tool warned that parameter
unpacking (*result) could raise a TypeError, but it failed to recognize that the variable had been
explicitly constructed as a tuple within the same commit, thus ensuring the operation’s safety.
Other significant sources of error are Over-sensitivity to Modification (17%) and Vague and
Unactionable Feedback (16%). These stem from the model’s cautious bias, leading to low-value
suggestions like "verifying" every new logic in scikit-learn/scikit-learn-25025, offering no
specific flaw or actionable advice. Such issues could be mitigated through prompt engineering.
Additionally, a Lack of Domain Knowledge (13%) contributes to false positives when the LLM
is unaware of project-specific conventions. In sympy/sympy-17801, the idiomatic expression is
S.true, used for checking symbolic boolean values, was incorrectly flagged, even though it is
the standard practice in the library to avoid a TypeError. This type of error could potentially be
addressed by providing additional domain knowledge in the prompt.
In summary, these findings demonstrate that false positives arise from a fundamental inability to
comprehend developer intent, code context, and domain-specific knowledge. Future work must
therefore enhance the reasoning abilities of LLMs and optimize prompts to generate more accurate,
actionable feedback.
Conclusion 6: Current ACR tools struggle to grasp three critical elements: developer intent,
code context, and domain knowledge. Addressing this requires advancing LLMs’ core reasoning
capabilities and employing better prompting strategies to improve the quality of review reports.
Impact of LLMs and Reasoning Enhancement. Following PR-Review’s superior performance, we
evaluated its efficacy on SWR-Bench with diverse LLMs, including variants with and without
reasoning-enhancement training (Table 8).
Overall, LLM performance with PR-Review remains suboptimal for practical code review (highest
𝑂𝑣𝑒𝑟𝑎𝑙𝑙-𝐹1 19%). Within this context, we first noted a significant divergence: the LLM performance
ranking on SWR-Bench does not consistently mirror trends observed in other SE benchmarks (e.g.,
SWE-Bench [18] and LiveCodeBench [14]). This discrepancy underscores SWR-Bench’s value in
capturing unique code review challenges and speculatively points to a potential misalignment be-
tween current LLM training/optimization and code review’s specific demands, where overemphasis
on other SE tasks might inadvertently reduce code review proficiency.
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

SWR-Bench: Assessing LLM Performance in Real-World Code Review Comment Generation
FSE137:17
Table 8. LLMs performance evaluation for code review with PR-Review on SWR-Bench.
LLMs
Overall Change-Actions
Functional Change-Actions
Precision
Recall
F1
Avg.
Count
Precision
Recall
F1
Avg.
Count
Reasoning LLM
Gemini-2.5-Pro
16.65
23.18
19.38
1.32
19.38
40.72
26.26
0.65
Gemini-2.5-Flash
16.88
13.91
15.25
0.78
17.92
24.42
20.67
0.41
DeepSeek-R1
14.61
25.5
18.58
1.66
15.55
50.62
23.79
1.06
GPT-o3
14.05
25.58
18.13
1.73
14.41
50.78
22.45
1.13
GPT-5
14.69
35.93
20.85
2.32
14.87
65.03
24.2
1.43
Qwen-2.5-R1-32B
11.69
20.86
14.98
1.69
12.65
38.19
19
0.93
Qwen-2.5-R1-14B
13.21
20.13
15.95
1.45
15.01
40.94
21.96
0.87
Qwen-2.5-R1-7B
6.83
8.33
7.51
1.16
6.8
11.18
8.46
0.5
Standard LLM
Claude-4-Opus
14.94
19.68
16.99
1.25
16.02
37.54
22.45
0.74
Claude-4-Sonnet
13.84
20.76
16.61
1.42
14.81
39.45
21.54
0.87
Claude-3.7-Sonnet
14.9
23.5
18.23
1.5
14.72
40.32
21.56
0.86
GPT-4o
14.13
27.78
18.73
1.86
15.61
44.48
23.11
0.85
DeepSeek-V3
15.49
20.17
17.52
1.22
16.15
33.56
21.81
0.61
Qwen-2.5-32B
11.43
17.28
13.76
1.44
13.68
29.71
18.73
0.68
Qwen-2.5-14B
8.67
9.38
9.01
1.03
13.65
17.68
15.41
0.4
Qwen-2.5-7B
8.87
16.86
11.63
1.8
10
24.6
14.22
0.77
Turning to the impact of reasoning enhancement training, we found further nuances. Models
like Gemini-2.5-Pro and DeepSeek-R1, known for their reasoning capabilities, generally performed
better. The Qwen-2.5 series provided a clear illustration: its standard versions (without specific
reasoning enhancement) yielded lower 𝐹1 scores (max 13.76%). However, their reasoning-enhanced
counterparts showed marked improvement. Specifically, Qwen-2.5-R1-14B achieved an 𝐹1 of 15.95%
(Qwen-2.5-R1-7B was an exception due to output formatting issues). This strongly suggests that
reasoning-enhancement training is a crucial factor for improving LLM effectiveness in code review.
Conclusion 7: Current LLMs underperform for practical review on SWR-Bench; however,
reasoning-enhanced models perform better, highlighting reasoning-enhanced as a crucial ad-
vancement path.
RQ3: How can the performance of automated code review tools be improved?
73
27
16
38
28
21
10
34
20
15
6
16
13
14
29
20
15
4
6
1
1
4
9
5
6
6
9
5
4
15
36
Gemini-2.5-Flash
Claude-3.7-Sonnet
DeepSeek-Reasoner
Gemini-2.5-Pro
GPT-4o
(a) Overlap across different models.
34
27
8
22
4
9
4
20
4
6
5
3
4
8
3
21
8
5
2
3
2
7
11
7
1
4
3
4
7
7
27
Gemini-2.5-Flash Run-1
Run-2
Run-3
Run-4
Run-5
(b) Overlap across multiple runs of the same model.
Fig. 7. Venn diagrams illustrating the overlap of identified change-actions.
Conclusion 7 indicates reasoning-enhanced LLMs perform better, the overall practical utility
of current ACR approaches remains limited. This prompted a closer examination of the nature of
LLM-generated reviews, particularly their reliability and consistency. Specifically, we questioned
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

FSE137:18 Zhengran Zeng, Ruikai Shi, Keke Han, Yixin Li, KC Sun, YD Wang, ZH Yu, Rui Xie, Wei Ye and Shikun Zhang
whether current LLMs provide stable and consistent feedback across different invocations or when
compared to other models.
To investigate this, we conducted an analysis into the overlap of specific ground truth change-
actions identified by different LLMs and by multiple runs of the same LLM using Venn diagrams
(Figures 7a and 7b). This analysis revealed considerable variability and instability in the sets of
successfully identified defects across these models and runs. Notably, for different LLMs only 36 suc-
cessfully identified change-actions overlapped (Figure 7a), for the same LLM over five independent
runs, only 27 successfully identified change-actions overlapped (Figure 7b). This low consistency
suggests unstable detection for most change-actions, implying significant randomness in LLM iden-
tification. While this variability might partly stem from the subjective nature of many evolutionary
changes, it strongly indicates that LLMs’ current grasp of code review nuances may be superficial,
with detections sometimes attributable to stochastic factors rather than deep comprehension.
LLM
Pull Request
 src
 main.py
 demo.py
 utils
·
+20-12
 src main.py
 utils demo.py
 examples reqs.txt
 README.md
Codebase
Final Report
n Pred Action 2
n Pred Action 4
n Pred Action 6
LLM
PROMPT: Review
this PR based on the
diff and codebase.
PROMPT: Aggregate
n draft reviews into a
refined review report.
Self-Agg: same LLM generates n reviews.
Multi-Agg: different LLMs generate n reviews.
 Draft Reports
n Pred Action 1
n Pred Action 6
 Draft Reports
n Pred Action 1
n Pred Action 4
 Draft Reports
n Pred Action 1
n Pred Action 2
Fig. 8. The overall workflow of the proposed Multi-Review approach.
This observed variability and the implied limitations of single-pass reviews directly motivated
our exploration into whether integrating diverse review outputs could yield a more comprehensive
and reliable final report. To this end, building upon PR-Review, we propose an enhanced approach
named Multi-Review (Figure 8). The core idea is to execute PR-Review (or another ACR tool)
multiple times on the same code change to generate 𝑛independent review reports, which are then
aggregated into a final review report using an additional LLM call. Specifically, all 𝑛reports are
concatenated into a single input, and an LLM is prompted with a meta-instruction to synthesize
them into a single, superior review. To validate Multi-Review’s efficacy, we conducted experiments
on SWR-Bench with two aggregation strategies: 1) Self-Agg, which employs the same LLM to
both generate and aggregate multiple independent review reports, testing whether a single model
can overcome its own stochasticity by consolidating multiple internal opinions; and 2) Multi-
Agg (Cross-Model Aggregation), which aggregates reports generated by multiple different LLMs,
leveraging the complementary strengths of different models to produce a more complete and
accurate review.
We selected Gemini-2.5-Flash as a representative LLM to investigate the impact of aggregat-
ing varying numbers of review reports (𝑛∈{0, 1, 3, 5, 10}, where 𝑛= 0 represents the baseline
PR-Review without aggregation). The results, presented in Figure 9a, demonstrate a significant
performance uplift with Multi-Review. Regardless of whether aggregating results from different
models or multiple runs of a single model, the 𝐹1 scores showed marked improvement. Specifically,
Gemini-2.5-Flash with Self-Agg (𝑛= 10) achieved an 𝐹1 of 21.91% (a 43.67% increase) and a 𝑅𝑒𝑐𝑎𝑙𝑙
of 30.44% (a 118.83% increase). This larger gain in 𝑅𝑒𝑐𝑎𝑙𝑙indicates Multi-Review’s effectiveness
in identifying more true defects. However, to further boost the 𝐹1 score, enhancing 𝑃𝑟𝑒𝑐𝑖𝑠𝑖𝑜𝑛by
reducing false positives remains a key priority. To further assess the generalizability of Multi-
Review, we experimented with smaller, open-source LLMs from the Qwen-2.5 series (7B, 14B, 32B),
investigating whether these models could achieve significant self-enhancement via Multi-Review
in resource-constrained scenarios. Specifically, Multi-Review consistently boosted 𝐹1 scores for the
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

SWR-Bench: Assessing LLM Performance in Real-World Code Review Comment Generation
FSE137:19
0
1
3
5
10
Aggregate Report Count (#)
10.0
12.5
15.0
17.5
20.0
22.5
Overall F1 (%)
0
1
3
5
10
Aggregate Report Count (#)
10
20
30
40
Overall Recall (%)
Gemini-2.5-Flash Multi-Agg
Gemini-2.5-Flash Self-Agg
Qwen-Chat-7B Self-Agg
Qwen-Chat-14B Self-Agg
Qwen-Chat-32B Self-Agg
(a) Performance of Multi-Review for varying numbers of aggre-
gated reports (𝑛).
0.00
0.01
0.02
0.03
0.04
0.05
0.06
0.07
Avg. Cost ($)
14
16
18
20
22
24
Overall F1 (%)
PR-Review
(Flash)
PR-Review
(Pro)
n=1
n=3
n=5
n=10
n=1
n=3 n=5
n=10
n=1
n=3
n=5
n=10
PR-Review + Gemini-2.5-Flash
PR-Review + Gemini-2.5-Pro
Gemini-2.5-Flash Self-Agg
Gemini-2.5-Flash Multi-Agg
Gemini-2.5-Pro Self-Agg
(b) Cost-benefit analysis across different
models and aggregation sizes (𝑛).
Fig. 9. Evaluation of the Multi-Review strategy. (a) shows the performance impact of aggregating varying
numbers of reports, while (b) presents the trade-off between Overall F1 and average API cost per PR (in US
dollars).
Qwen series. Notably, Qwen-Chat-7B Self-Agg (𝑛= 10) improved its 𝐹1 by 26.13% (to 14.67%), and
Qwen-Chat-32B Self-Agg (𝑛= 10) enhanced its 𝐹1 by 19.25% (to 16.41%), allowing their aggregated
results to approach the baseline PR-Review performance of some larger commercial LLMs.
Furthermore, Figure 9b presents a cost-benefit analysis of the Multi-Review, including results with
the larger Gemini-2.5-Pro model. A key finding is that Gemini-2.5-Flash Self-Agg (𝑛=5) achieves an
F1 of 20.48%, surpassing the single-pass Gemini-2.5-Pro baseline (F1: 19.38%) at a lower API cost
($3.68E-03 vs. $5.86E-03). This demonstrates that multiple runs of a smaller, cheaper model can be
more cost-effective than a single run of a larger, more expensive one. Meanwhile, Gemini-2.5-Pro
Self-Agg (𝑛=5) achieves the highest F1 of 23.84% (a 23.01% improvement over its single-pass baseline),
confirming that the aggregation strategy benefits models across the entire performance spectrum.
However, performance gains exhibit diminishing returns beyond 𝑛=5, while costs continue to
scale linearly, suggesting that 𝑛=5 represents a practical sweet spot. Regarding time efficiency,
generating initial drafts takes 38.75 seconds per PR. This step is fully parallelizable, meaning its
duration remains constant regardless of 𝑛. The final aggregation step requires minimal additional
time, taking only 4.67, 8.76, 12.80, and 24.35 seconds for 𝑛=1, 3, 5, and 10 reports respectively. This
low time overhead confirms the approach is practical for real-world deployment.
Conclusion 8: Multi-Review improves code review performance across difference model sizes
by aggregating multiple reports. Notably, aggregating several runs of a smaller model can
cost-effectively match or exceed the accuracy of a single run from a larger model. This highlights
the substantial potential of improving current tools through systematic output integration.
5
Discussion and Implications
The findings from our evaluation on SWR-Bench provide concrete insights into automated code
review. We summarize these implications as follow:
Objective Evaluation through Fact Matching: As demonstrated in Conclusion 1, traditional
metrics like text similarity fail to reflect true semantic quality. Furthermore, relying on subjective
large language model scoring to evaluate generated reviews shows poor agreement with human
expert preferences. Researchers should therefore evaluate tools using objective fact matching
against predefined ground truth to verify semantic issues reliably.
From Code Generation to Code Critique. Code review is fundamentally different from code
generation. Generation relies on pattern completion, whereas review demands logical deduction,
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

FSE137:20 Zhengran Zeng, Ruikai Shi, Keke Han, Yixin Li, KC Sun, YD Wang, ZH Yu, Rui Xie, Wei Ye and Shikun Zhang
counterfactual reasoning, and cross-file dependency analysis. This distinction explains why models
optimized primarily for code generation underperform on SWR-Bench, and why reasoning-enhanced
models consistently achieve better results (Conclusion 7). Future work should move beyond prompt
engineering and explore dedicated training objectives for code critique, such as reinforcement
learning with process reward models that evaluate the logical rigor of identified issues rather than
only the correctness of generated code.
Mitigating False Positives in Practice. The high false positive rates of current tools (Conclusion 2)
severely limit their practical adoption by overwhelming developers with invalid alerts. Since
most false positives arise from insufficient contextual understanding and ignorance of project-
specific conventions (Table 7), future systems must integrate historical repository data to learn
implicit coding norms. Furthermore, tool designers should decouple feedback based on severity.
Because models detect functional changes more reliably than evolutionary ones (Conclusion 5),
presenting functional defects as mandatory requirements and evolutionary suggestions as optional
recommendations will significantly enhance the practical value of automated reviews.
Enhancing Review Performance through Report Aggregation: Single passes of large language
models exhibit significant randomness and miss different defects across runs. Aggregating multiple
independent review reports into a single final output effectively mitigates this stochasticity and
significantly improves overall defect detection. Conclusion 8 demonstrates that this aggregation
strategy also allows multiple runs of smaller models to match or exceed the performance of a single
run from a larger and more expensive model.
6
Threats to Validity
Internal Validity. Threats to internal validity mainly arise from potential bugs in our implemen-
tation and the accuracy of manual verification. To mitigate these threats, we have conducted a
detailed review of the code, and our code has been made publicly available at [1] for independent
verification.
External Validity. Threats to external validity mainly arise from the LLMs, ACR tools, and the
quality of the selected projects. To mitigate these threats, we chose representative tools based on an
extensive literature review and selected 12 popular, well-maintained open-source GitHub projects.
Construct Validity. The primary threat to construct validity centers on the accuracy and reliability
of our objective LLM-based evaluation method. We mitigated this through human validation, which
confirmed high consistency between our method’s results and human expert evaluations.
7
Conclusions
This paper addressed critical limitations in existing automated code review benchmarks by intro-
ducing SWR-Bench, a novel benchmark featuring 1,000 real-world pull requests with full project
context and an objective LLM-based evaluation. Our study on SWR-Bench revealed that while
current LLM-based ACR systems generally underperform, they show better aptitude for functional
error detection. To improve performance, we proposed and validated a multi-review aggregation
strategy that significantly boosts 𝐹1 scores. SWR-Bench, along with our findings and proposed
enhancement, provides a more realistic platform and valuable insights for advancing practical ACR
research and development.
Data Availability
The replication package for our study, containing the necessary source code and scripts to reproduce
our experiments, is available at the anonymous repository [1].
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

SWR-Bench: Assessing LLM Performance in Real-World Code Review Comment Generation
FSE137:21
References
[1] [n. d.]. SWR-Bench: A Benchmarking Suite for Serverless Cold Start Time Reduction. https://github.com/ZZR0/
SWRench. Accessed: 2025-05-23..
[2] Deepika Badampudi, Michael Unterkalmsteiner, and Ricardo Britto. 2023. Modern Code Reviews - Survey of Literature
and Practice. ACM Trans. Softw. Eng. Methodol. 32, 4 (2023), 107:1–107:61. doi:10.1145/3585004
[3] Gabriele Bavota and Barbara Russo. 2015. Four eyes are better than two: On the impact of code reviews on software
quality. In 2015 IEEE International Conference on Software Maintenance and Evolution, ICSME 2015, Bremen, Germany,
September 29 - October 1, 2015, Rainer Koschke, Jens Krinke, and Martin P. Robillard (Eds.). IEEE Computer Society,
81–90. doi:10.1109/ICSM.2015.7332454
[4] Moritz Beller, Alberto Bacchelli, Andy Zaidman, and Elmar Jürgens. 2014. Modern code reviews in open-source projects:
which problems do they fix?. In 11th Working Conference on Mining Software Repositories, MSR 2014, Proceedings, May
31 - June 1, 2014, Hyderabad, India, Premkumar T. Devanbu, Sung Kim, and Martin Pinzger (Eds.). ACM, 202–211.
doi:10.1145/2597073.2597082
[5] Mert Cemri, Melissa Z. Pan, Shuyi Yang, Lakshya A. Agrawal, Bhavya Chopra, Rishabh Tiwari, Kurt Keutzer, Aditya G.
Parameswaran, Dan Klein, Kannan Ramchandran, Matei Zaharia, Joseph E. Gonzalez, and Ion Stoica. 2025. Why Do
Multi-Agent LLM Systems Fail? CoRR abs/2503.13657 (2025). arXiv:2503.13657 doi:10.48550/ARXIV.2503.13657
[6] H. Alperen Çetin, Emre Dogan, and Eray Tüzün. 2021. A review of code reviewer recommendation studies: Challenges
and future directions. Sci. Comput. Program. 208 (2021), 102652. doi:10.1016/J.SCICO.2021.102652
[7] H. Alperen Çetin, Emre Dogan, and Eray Tüzün. 2021. A review of code reviewer recommendation studies: Challenges
and future directions. Sci. Comput. Program. 208 (2021), 102652. doi:10.1016/J.SCICO.2021.102652
[8] Jialin Chen, Aosong Feng, Ziyu Zhao, Juan Garza, Gaukhar Nurbek, Cheng Qin, Ali Maatouk, Leandros Tassiulas,
Yifeng Gao, and Rex Ying. 2025. MTBench: A Multimodal Time Series Benchmark for Temporal Reasoning and
Question Answering. CoRR abs/2503.16858 (2025). arXiv:2503.16858 doi:10.48550/ARXIV.2503.16858
[9] Umut Cihan, Vahid Haratian, Arda Içöz, Mert Kaan Gül, Ömercan Devran, Emircan Furkan Bayendur, Baykal Mehmet
Uçar, and Eray Tüzün. 2024. Automated Code Review In Practice. CoRR abs/2412.18531 (2024). arXiv:2412.18531
doi:10.48550/ARXIV.2412.18531
[10] Yegor Denisov-Blanch, Igor Ciobanu, Simon Obstbaum, and Michal Kosinski. 2024. Predicting Expert Evaluations in
Software Code Reviews. CoRR abs/2409.15152 (2024). arXiv:2409.15152 doi:10.48550/ARXIV.2409.15152
[11] Enrico Fregnan, Fernando Petrulio, and Alberto Bacchelli. 2022. The evolution of the code during review: an investiga-
tion on review changes. Empir. Softw. Eng. 27, 7 (2022), 177. doi:10.1007/S10664-022-10205-7
[12] Google. 2025. Gemini 2.5: Our most intelligent ai model. https://blog.google/technology/google-deepmind/gemini-
model-thinking-updates-march-2025/#gemini-2-5-thinking.
[13] Jiawei Gu, Xuhui Jiang, Zhichao Shi, Hexiang Tan, Xuehao Zhai, Chengjin Xu, Wei Li, Yinghan Shen, Shengjie Ma,
Honghao Liu, Yuanzhuo Wang, and Jian Guo. 2024. A Survey on LLM-as-a-Judge. CoRR abs/2411.15594 (2024).
arXiv:2411.15594 doi:10.48550/ARXIV.2411.15594
[14] Naman Jain, King Han, Alex Gu, Wen-Ding Li, Fanjia Yan, Tianjun Zhang, Sida Wang, Armando Solar-Lezama, Koushik
Sen, and Ion Stoica. 2025. LiveCodeBench: Holistic and Contamination Free Evaluation of Large Language Models for
Code. In The Thirteenth International Conference on Learning Representations, ICLR 2025, Singapore, April 24-28, 2025.
OpenReview.net. https://openreview.net/forum?id=chfJJYC3iL
[15] Imen Jaoua, Oussama Ben Sghaier, and Houari A. Sahraoui. 2025. Combining Large Language Models with Static
Analyzers for Code Review Generation. CoRR abs/2502.06633 (2025). arXiv:2502.06633 doi:10.48550/ARXIV.2502.06633
[16] Imen Jaoua, Oussama Ben Sghaier, and Houari A. Sahraoui. 2025. Combining Large Language Models with Static
Analyzers for Code Review Generation. CoRR abs/2502.06633 (2025). arXiv:2502.06633 doi:10.48550/ARXIV.2502.06633
[17] Yanjie Jiang, Hui Liu, Tianyi Chen, Fu Fan, Chunhao Dong, Kui Liu, and Lu Zhang. 2025. Deep Assessment of
Code Review Generation Approaches: Beyond Lexical Similarity. CoRR abs/2501.05176 (2025). arXiv:2501.05176
doi:10.48550/ARXIV.2501.05176
[18] Carlos E. Jimenez, John Yang, Alexander Wettig, Shunyu Yao, Kexin Pei, Ofir Press, and Karthik R. Narasimhan. 2024.
SWE-bench: Can Language Models Resolve Real-world Github Issues?. In The Twelfth International Conference on
Learning Representations, ICLR 2024, Vienna, Austria, May 7-11, 2024. OpenReview.net. https://openreview.net/forum?
id=VTF8yNQM66
[19] René Just, Darioush Jalali, and Michael D. Ernst. 2014. Defects4J: a database of existing faults to enable controlled
testing studies for Java programs. In International Symposium on Software Testing and Analysis, ISSTA ’14, San Jose, CA,
USA - July 21 - 26, 2014, Corina S. Pasareanu and Darko Marinov (Eds.). ACM, 437–440. doi:10.1145/2610384.2628055
[20] Sumith Kulal, Panupong Pasupat, Kartik Chandra, Mina Lee, Oded Padon, Alex Aiken, and Percy Liang. 2019. SPoC:
Search-based Pseudocode to Code. In Advances in Neural Information Processing Systems 32: Annual Conference on Neural
Information Processing Systems 2019, NeurIPS 2019, December 8-14, 2019, Vancouver, BC, Canada, Hanna M. Wallach,
Hugo Larochelle, Alina Beygelzimer, Florence d’Alché-Buc, Emily B. Fox, and Roman Garnett (Eds.). 11883–11894.
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

FSE137:22 Zhengran Zeng, Ruikai Shi, Keke Han, Yixin Li, KC Sun, YD Wang, ZH Yu, Rui Xie, Wei Ye and Shikun Zhang
https://proceedings.neurips.cc/paper/2019/hash/7298332f04ac004a0ca44cc69ecf6f6b-Abstract.html
[21] Jahnavi Kumar and Sridhar Chimalakonda. 2024. Code Review Automation Via Multi-task Federated LLM - An
Empirical Study. CoRR abs/2412.15676 (2024). arXiv:2412.15676 doi:10.48550/ARXIV.2412.15676
[22] Zhiyu Li, Shuai Lu, Daya Guo, Nan Duan, Shailesh Jannu, Grant Jenks, Deep Majumder, Jared Green, Alexey Svy-
atkovskiy, Shengyu Fu, and Neel Sundaresan. 2022. Automating code review activities by large-scale pre-training. In
Proceedings of the 30th ACM Joint European Software Engineering Conference and Symposium on the Foundations of
Software Engineering, ESEC/FSE 2022, Singapore, Singapore, November 14-18, 2022, Abhik Roychoudhury, Cristian Cadar,
and Miryung Kim (Eds.). ACM, 1035–1047. doi:10.1145/3540250.3549081
[23] Chin-Yew Lin. 2004. ROUGE: A Package for Automatic Evaluation of Summaries. In Text Summarization Branches Out.
Association for Computational Linguistics, Barcelona, Spain, 74–81. https://aclanthology.org/W04-1013/
[24] Junwei Liu, Kaixin Wang, Yixuan Chen, Xin Peng, Zhenpeng Chen, Lingming Zhang, and Yiling Lou. 2024. Large
Language Model-Based Agents for Software Engineering: A Survey. CoRR abs/2409.02977 (2024). arXiv:2409.02977
doi:10.48550/ARXIV.2409.02977
[25] Junyi Lu, Xiaojia Li, Zihan Hua, Lei Yu, Shiqi Cheng, Li Yang, Fengjun Zhang, and Chun Zuo. 2025. DeepCRCEval:
Revisiting the Evaluation of Code Review Comment Generation. In Fundamental Approaches to Software Engineering -
28th International Conference, FASE 2025, Held as Part of the International Joint Conferences on Theory and Practice of
Software, ETAPS 2025, Hamilton, ON, Canada, May 3-8, 2025, Proceedings (Lecture Notes in Computer Science, Vol. 15693),
Artur Boronat and Gordon Fraser (Eds.). Springer, 43–64. doi:10.1007/978-3-031-90900-9_3
[26] Junyi Lu, Lei Yu, Xiaojia Li, Li Yang, and Chun Zuo. 2023. LLaMA-Reviewer: Advancing Code Review Automation with
Large Language Models through Parameter-Efficient Fine-Tuning. In 34th IEEE International Symposium on Software
Reliability Engineering, ISSRE 2023, Florence, Italy, October 9-12, 2023. IEEE, 647–658. doi:10.1109/ISSRE59848.2023.00026
[27] Kishore Papineni, Salim Roukos, Todd Ward, and Wei-Jing Zhu. 2002. Bleu: a Method for Automatic Evaluation of
Machine Translation. In Proceedings of the 40th Annual Meeting of the Association for Computational Linguistics, July
6-12, 2002, Philadelphia, PA, USA. ACL, 311–318. doi:10.3115/1073083.1073135
[28] PMD. 2000. https://pmd.github.io/
[29] Qodana AI. 2025. PR-Agent: AI-Powered Pull Request Agent. https://github.com/qodo-ai/pr-agent
[30] Zeeshan Rasheed, Malik Abdul Sami, Muhammad Waseem, Kai-Kristian Kemell, Xiaofeng Wang, Anh Nguyen, Kari
Systä, and Pekka Abrahamsson. 2024. AI-powered Code Review with LLMs: Early Results. CoRR abs/2404.18496 (2024).
arXiv:2404.18496 doi:10.48550/ARXIV.2404.18496
[31] Giovanni Rosa, Luca Pascarella, Simone Scalabrino, Rosalia Tufano, Gabriele Bavota, Michele Lanza, and Rocco Oliveto.
2023. A comprehensive evaluation of SZZ Variants through a developer-informed oracle. J. Syst. Softw. 202 (2023),
111729. doi:10.1016/J.JSS.2023.111729
[32] Oussama Ben Sghaier and Houari A. Sahraoui. 2024. Improving the Learning of Code Review Successive Tasks with
Cross-Task Knowledge Distillation. Proc. ACM Softw. Eng. 1, FSE (2024), 1086–1106. doi:10.1145/3643775
[33] Lin Shi, Weicheng Ma, and Soroush Vosoughi. 2024. Judging the Judges: A Systematic Investigation of Position Bias in
Pairwise Comparative Assessments by LLMs. CoRR abs/2406.07791 (2024). arXiv:2406.07791 doi:10.48550/ARXIV.2406.
07791
[34] SonarSource. 2006. SonarQube. https://www.sonarsource.com/products/sonarqube/
[35] Tao Sun, Jian Xu, Yuanpeng Li, Zhao Yan, Ge Zhang, Lintao Xie, Lu Geng, Zheng Wang, Yueyan Chen, Qin Lin, Wenbo
Duan, and Kaixin Sui. 2025. BitsAI-CR: Automated Code Review via LLM in Practice. CoRR abs/2501.15134 (2025).
arXiv:2501.15134 doi:10.48550/ARXIV.2501.15134
[36] Xunzhu Tang, Kisub Kim, Yewei Song, Cedric Lothritz, Bei Li, Saad Ezzini, Haoye Tian, Jacques Klein, and Tegawendé F.
Bissyandé. 2024. CodeAgent: Autonomous Communicative Agents for Code Review. In Proceedings of the 2024
Conference on Empirical Methods in Natural Language Processing, EMNLP 2024, Miami, FL, USA, November 12-16, 2024,
Yaser Al-Onaizan, Mohit Bansal, and Yun-Nung Chen (Eds.). Association for Computational Linguistics, 11279–11313.
https://aclanthology.org/2024.emnlp-main.632
[37] Minaoar Hossain Tanzil, Junaed Younus Khan, and Gias Uddin. 2024. ChatGPT Incorrectness Detection in Software
Reviews. In Proceedings of the 46th IEEE/ACM International Conference on Software Engineering, ICSE 2024, Lisbon,
Portugal, April 14-20, 2024. ACM, 180:1–180:12. doi:10.1145/3597503.3639194
[38] Aman Singh Thakur, Kartik Choudhary, Venkat Srinik Ramayapally, Sankaran Vaidyanathan, and Dieuwke Hupkes.
2024. Judging the Judges: Evaluating Alignment and Vulnerabilities in LLMs-as-Judges. CoRR abs/2406.12624 (2024).
arXiv:2406.12624 doi:10.48550/ARXIV.2406.12624
[39] Patanamon Thongtanunam, Chanathip Pornprasit, and Chakkrit Tantithamthavorn. 2022. AutoTransform: Automated
Code Transformation to Support Modern Code Review Process. In 44th IEEE/ACM 44th International Conference on
Software Engineering, ICSE 2022, Pittsburgh, PA, USA, May 25-27, 2022. ACM, 237–248. doi:10.1145/3510003.3510067
[40] Hugo Touvron, Thibaut Lavril, Gautier Izacard, Xavier Martinet, Marie-Anne Lachaux, Timothée Lacroix, Baptiste
Rozière, Naman Goyal, Eric Hambro, Faisal Azhar, Aurélien Rodriguez, Armand Joulin, Edouard Grave, and Guillaume
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.

SWR-Bench: Assessing LLM Performance in Real-World Code Review Comment Generation
FSE137:23
Lample. 2023. LLaMA: Open and Efficient Foundation Language Models. CoRR abs/2302.13971 (2023). arXiv:2302.13971
doi:10.48550/ARXIV.2302.13971
[41] Rosalia Tufano, Simone Masiero, Antonio Mastropaolo, Luca Pascarella, Denys Poshyvanyk, and Gabriele Bavota.
2022. Using Pre-Trained Models to Boost Code Review Automation. In 44th IEEE/ACM 44th International Conference on
Software Engineering, ICSE 2022, Pittsburgh, PA, USA, May 25-27, 2022. ACM, 2291–2302. doi:10.1145/3510003.3510621
[42] Rosalia Tufano, Luca Pascarella, Michele Tufano, Denys Poshyvanyk, and Gabriele Bavota. 2021. Towards Automating
Code Review Activities. In 43rd IEEE/ACM International Conference on Software Engineering, ICSE 2021, Madrid, Spain,
22-30 May 2021. IEEE, 163–174. doi:10.1109/ICSE43902.2021.00027
[43] Yue Wang, Weishi Wang, Shafiq R. Joty, and Steven C. H. Hoi. 2021. CodeT5: Identifier-aware Unified Pre-trained
Encoder-Decoder Models for Code Understanding and Generation. In Proceedings of the 2021 Conference on Em-
pirical Methods in Natural Language Processing, EMNLP 2021, Virtual Event / Punta Cana, Dominican Republic, 7-11
November, 2021, Marie-Francine Moens, Xuanjing Huang, Lucia Specia, and Scott Wen-tau Yih (Eds.). Association for
Computational Linguistics, 8696–8708. doi:10.18653/V1/2021.EMNLP-MAIN.685
[44] Yidong Wang, Zhuohao Yu, Wenjin Yao, Zhengran Zeng, Linyi Yang, Cunxiang Wang, Hao Chen, Chaoya Jiang, Rui Xie,
Jindong Wang, Xing Xie, Wei Ye, Shikun Zhang, and Yue Zhang. 2024. PandaLM: An Automatic Evaluation Benchmark
for LLM Instruction Tuning Optimization. In The Twelfth International Conference on Learning Representations, ICLR
2024, Vienna, Austria, May 7-11, 2024. OpenReview.net. https://openreview.net/forum?id=5Nn2BLV7SB
[45] Wikipedia. 2025. Cohen’s kappa. https://en.wikipedia.org/wiki/Cohen’s_kappa. Accessed: 2025-09-11.
[46] Wikipedia. 2025. Stratified sampling. https://en.wikipedia.org/wiki/Stratified_sampling. Accessed: 2025-09-11.
[47] Chunqiu Steven Xia, Yinlin Deng, Soren Dunn, and Lingming Zhang. 2024. Agentless: Demystifying LLM-based
Software Engineering Agents. CoRR abs/2407.01489 (2024). arXiv:2407.01489 doi:10.48550/ARXIV.2407.01489
[48] John Yang, Carlos E. Jimenez, Alexander Wettig, Kilian Lieret, Shunyu Yao, Karthik Narasimhan, and Ofir Press.
2024. SWE-agent: Agent-Computer Interfaces Enable Automated Software Engineering. In Advances in Neural
Information Processing Systems 38: Annual Conference on Neural Information Processing Systems 2024, NeurIPS 2024,
Vancouver, BC, Canada, December 10 - 15, 2024, Amir Globersons, Lester Mackey, Danielle Belgrave, Angela Fan,
Ulrich Paquet, Jakub M. Tomczak, and Cheng Zhang (Eds.).
http://papers.nips.cc/paper_files/paper/2024/hash/
5a7c947568c1b1328ccc5230172e1e7c-Abstract-Conference.html
[49] Zezhou Yang, Cuiyun Gao, Zhaoqiang Guo, Zhenhao Li, Kui Liu, Xin Xia, and Yuming Zhou. 2024. A Survey on
Modern Code Review: Progresses, Challenges and Opportunities. CoRR abs/2405.18216 (2024). arXiv:2405.18216
doi:10.48550/ARXIV.2405.18216
[50] Zhengran Zeng, Yidong Wang, Rui Xie, Wei Ye, and Shikun Zhang. 2024. CoderUJB: An Executable and Unified Java
Benchmark for Practical Programming Scenarios. In Proceedings of the 33rd ACM SIGSOFT International Symposium on
Software Testing and Analysis, ISSTA 2024, Vienna, Austria, September 16-20, 2024, Maria Christakis and Michael Pradel
(Eds.). ACM, 124–136. doi:10.1145/3650212.3652115
[51] Zhengran Zeng, Yuqun Zhang, Haotian Zhang, and Lingming Zhang. 2021. Deep just-in-time defect prediction: how
far are we?. In ISSTA ’21: 30th ACM SIGSOFT International Symposium on Software Testing and Analysis, Virtual Event,
Denmark, July 11-17, 2021, Cristian Cadar and Xiangyu Zhang (Eds.). ACM, 427–438. doi:10.1145/3460319.3464819
[52] Lianghui Zhu, Xinggang Wang, and Xinlong Wang. 2025. JudgeLM: Fine-tuned Large Language Models are Scalable
Judges. In The Thirteenth International Conference on Learning Representations, ICLR 2025, Singapore, April 24-28, 2025.
OpenReview.net. https://openreview.net/forum?id=xsELpEPn4A
Received 2026-02-25; accepted 2026-03-24
Proc. ACM Softw. Eng., Vol. 3, No. FSE, Article FSE137. Publication date: July 2026.
