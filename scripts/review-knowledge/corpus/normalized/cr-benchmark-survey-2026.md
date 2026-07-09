A Survey of Code Review Benchmarks and Evaluation
Practices in Pre-LLM and LLM Era
TAUFIQUL ISLAM KHAN, University of Manitoba, Department of Computer Science, Canada
SHAOWEI WANG, University of Manitoba, Department of Computer Science, Canada
HAOXIANG ZHANG, Huawei Canada, Centre for Software Excellence, Canada
TSE-HSUN CHEN, Concordia University, Canada
Code review is a critical practice in modern software engineering, helping developers detect defects early,
improve code quality, and facilitate knowledge sharing. With the rapid advancement of large language models
(LLMs), a growing body of work has explored automated support for code review. However, progress in this
area is hindered by the lack of a systematic understanding of existing benchmarks and evaluation practices.
Current code review datasets are scattered, vary widely in design, and provide limited insight into what
review capabilities are actually being assessed. In this paper, we present a comprehensive survey of code
review benchmarks spanning both the Pre-LLM and LLM eras (2015–2025). We analyze 99 research papers (58
Pre-LLM era and 41 LLM era) and extract key metadata, including datasets, evaluation metrics, data sources,
and target tasks. Based on this analysis, we propose a multi-level taxonomy that organizes code review research
into five domains and 18 fine-grained tasks. Our study reveals a clear shift toward end-to-end generative peer
review, increasing multilingual coverage, and a decline in standalone change understanding tasks. We further
identify limitations of current benchmarks and outline future directions, including broader task coverage,
dynamic runtime evaluation, and taxonomy-guided fine-grained assessment. This survey provides a structured
foundation for developing more realistic and comprehensive benchmarks for LLM-based code review.
Additional Key Words and Phrases: Survey, Benchmark, Evaluation Strategies, Code Review, LLM
ACM Reference Format:
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen. 2026. A Survey of Code Review
Benchmarks and Evaluation Practices in Pre-LLM and LLM Era. 1, 1 (February 2026), 37 pages. https://doi.
org/10.1145/nnnnnnn.nnnnnnn
1
Introduction
Code review is an important practice in software engineering where developers inspect code
changes before they are merged into the main branch. It has become an essential activity in
modern software development as reviewing code changes, especially early on, helps developers find
code defects or design issues before such failures grow into bigger problems. Code review makes
software easier to maintain and safer to release. In addition, code review facilitates collaborative
learning and helps train new team members [3]. When developers read and comment on a code
change, they may suggest clearer or more efficient ways to write the code. These discussions,
explanations, and suggestions help teammates understand the whole project better and learn good
Authors’ Contact Information: Taufiqul Islam khan, University of Manitoba, Department of Computer Science, Canada,
khanti@myumanitoba.ca; Shaowei Wang, University of Manitoba, Department of Computer Science, Canada, shaowei.
wang@umanitoba.ca; Haoxiang Zhang, Huawei Canada, Centre for Software Excellence, Canada, haoxiang.zhang@acm.org;
Tse-Hsun Chen, Concordia University, Montreal, Canada, peterc@encs.concordia.ca.
Permission to make digital or hard copies of all or part of this work for personal or classroom use is granted without fee
provided that copies are not made or distributed for profit or commercial advantage and that copies bear this notice and the
full citation on the first page. Copyrights for components of this work owned by others than the author(s) must be honored.
Abstracting with credit is permitted. To copy otherwise, or republish, to post on servers or to redistribute to lists, requires
prior specific permission and/or a fee. Request permissions from permissions@acm.org.
© 2026 Copyright held by the owner/author(s). Publication rights licensed to ACM.
ACM XXXX-XXXX/2026/2-ART
https://doi.org/10.1145/nnnnnnn.nnnnnnn
, Vol. 1, No. 1, Article . Publication date: February 2026.
arXiv:2602.13377v1 [cs.SE] 13 Feb 2026

2
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
coding practices. Research has shown that early code review can greatly reduce defects before
testing and deployment [6, 64].
However, manual code review is time-consuming and often requires significant effort from
development teams [6]. The recent emergence of Large Language Models (LLMs) has led to the
development of various automation approaches, which show promising results in streamlining
the code review process [87, 102]. Since LLMs are typically trained on vast code repositories,
including commit messages and discussions, they are well-equipped to review pull requests, suggest
improvements, and identify potential issues within code changes.
To properly measure the proposed code review approaches, we need reliable benchmarks
to analyze and understand. Benchmarks demonstrate how well approaches perform on real code
review scenarios. They provide a standard dataset and evaluation metrics so that approaches can
be compared in a fair and consistent way. A few code review benchmarks have been created
recently, such as CodeReviewer [50], CodeReviewQA [51] and CodeFuse-CR-Bench [26], but
these resources are scattered and are not well organized or widely known. In other areas, such
as code generation [100], there are survey papers that summarize available benchmarks and
evaluation strategies, which help researchers understand which datasets to use and how to measure
performance in a particular problem statement. For code review, we still have limited knowledge
about how existing benchmarks are built, what review skills they test and how closely they reflect
real code review practice. This gap presents the need for a more systematic study of code review
benchmarks for LLMs.
To address this gap, we conducted a comprehensive survey that examines code review
benchmarks across both the Pre-LLM and LLM eras from 2015 to present. Our goal is to gather
datasets that have been used in prior code review papers, thereby identifying the code review
tasks that have been measured and analyzed, and understanding how these datasets are used to
measure the code review tasks. We also compare traditional benchmarks in Pre-LLM era with
newer LLM-focused ones to understand how the datasets and evaluation strategy of code review
have changed over time. To conduct a systematic survey of code review datasets and evaluation
metrics, we have collected 61 Pre-LLM code review papers and 45 LLM-based papers published
between Jan. 2015 and Dec. 2025. These papers cover a wide range of areas within code review. We
analyzed key metadata from each study, which includes the datasets, evaluation metrics, and the
specific code review tasks involved in that research.
Our survey provides a comprehensive analysis of the code review landscape, offering the
following key contributions:

A Multi-Level Taxonomy of Code Review Tasks: We categorize existing studies into
five high-level domains—Review Prioritization/Selection, Change Understanding and
Analysis, Peer Review, Review Assessment and Analysis, and Code Refinement—further
subdivided into 18 specific sub-tasks. This framework enables researchers to identify
the most suitable benchmarks for specific research objectives.

Systematic Classification of Benchmarks and Metrics: We classify prior research by
evaluation metrics, data sources, and granularity levels. This systematic mapping assists
researchers in selecting appropriate metrics and facilitates more robust cross-study
comparisons.

Shift from Pre-LLM to LLM: We observe a profound shift in research focus between the
Pre-LLM and LLM eras. Notably, Change Understanding and Analysis has transitioned
from a research cornerstone (14 datasets) to a nearly absent standalone topic (1 dataset),
as tasks are increasingly consolidated into end-to-end generative processes. Our findings
show that Peer Review tasks have come to dominate the field, accounting for nearly
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
3
 5 H Y L H Z  8 Q L W V
  H  J   F R P P L W  3 5 
 5 H S R V L W R U \
 3 U L R U L W L ] D W L R Q
  6 H O H F W L R Q
 8 Q G H U V W D Q
 D Q G 
 $ Q D O \ V L  < H V
 3 D V V
 5 H Y L V L R Q
 / R Z  3 U L R U L W L ] H G
 8 Q L W V
 + L J K  3 U L R U L W L ] H G
 8 Q L W V
 3 H H U
  5 H Y L H Z
 ) H H G E D F N   H  J  
 U H Y L H Z  F R P P H Q W 
 5 H Y L V H G  8 Q L W V
 $ X W R  5 H J U
 H V V L R Q
 & R G H  5 H Y L H Z
 0 H U J H
 1 R
Fig. 1. Code Review Process.
60% of all datasets in the LLM era, compared to the more balanced task distribution
of the Pre-LLM era. We observe a significant evolution from language-specific studies
to cross-language generalization. While 59% of Pre-LLM datasets focused on a single
language, the LLM era is characterized by highly multilingual benchmarks, with 34% of
datasets now covering nine or more programming languages.

Future Directions for Improving LLM Benchmarks: We identify limitations of
current benchmarks for LLM and propose future directions, in terms of improving task
coverage (e.g., developing LLM benchmark for macro-level review responsibilities such
as impact analysis and commit decomposition), transitioning from Static to Dynamic
Evaluation (e.g., developing benchmark include runtime evaluation such as build success
rate, test case verification), and granular benchmarking via Task Taxonomy (e.g., measure
an LLM’s issue resolving capabilities for sub-tasks).
2
Background & Related work
2.1
Code Review process
Figure 1 presents an overview of the code review process. The code review process starts when
a developer submits a review unit, which could be of different granularity, such as chunk-level,
file-level, commit-level, or pull request (PR)-level. The code review process includes the following
steps:
Prioritization/Selection While doing the analysis, some commits appear do not contain much
valuable information or work that would be prioritized in low rank. Those commits demonstrate
that every change is not equally important. Therefore, before formal review, the first task is to
decide which commits worth reviewing and are given a high priority. This priority can come from
information such as the issue description, the part of the code that was changed, or how urgent
, Vol. 1, No. 1, Article . Publication date: February 2026.

4
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
the fix is. Low-impact changes are often deferred, while commits involving major bug fixes, new
features, or critical code paths are often prioritized.
Understanding and Analysis Next, the high-prioritized code changes are moved forward to
the next step. During this step, useful background information is collected to help reviewers
understand the change. This includes details about the affected files, related issues, and any historical
information that may help explain the motivation behind the update. A short summary of the
code modifications is also prepared so the reviewer can quickly see what was changed, why it was
changed, and how large the update is.
Peer Review With this context, the reviewer starts the actual evaluation of the commit. They
check whether a code change is logically correct, follows project standards, and does not introduce
new risks. During this stage, reviewers leave comments and suggestions pointing out unclear logic,
potential bugs, style issues, or areas that need improvement.
Revision Once the review is complete, a decision is made about whether the commit is acceptable.
If the reviewer finds the change correct and complete, it proceeds to merging. If issues are found,
the commit is returned to the developer for revisions. The developer then revisits the units in the
code repository, updates the implementation based on the feedback, and makes any necessary
adjustments to related tests or documentation. Typically, automated regression tests may run to
ensure the revised code does not break existing functionality. The updated code changes then enter
the cycle again, going through review until it satisfies all the requirements and is finally merged
into the main branch of the codebase. This cycle continues until the change becomes stable, correct,
and ready for integration.
2.2
Survey on coding tasks benchmark
Paul et al. [69] critically review evaluation practices for LLM-based code generation, arguing that
benchmarks and metrics remain immature despite rapid model progress. They have surveyed major
models and common benchmarks (e.g., HumanEval, APPS, MBPP) and analyzed how tasks are
sourced, cleaned, decontaminated, and structured across different granularities and languages. They
compare similarity metrics (BLEU, ROUGE, CodeBLEU) with execution-based ones (pass rates,
pass@k), noting their weak correlation with true correctness and developer usefulness. They call for
richer, scenario-driven benchmarks with meaningful metadata. However, unlike code generation,
no systematic benchmark study yet exists for code review tasks, leaving evaluation inconsistent
and incomplete.
Wang et al. [101] survey how benchmarks for code LLMs and agents align with the SDLC
phases, covering tasks such as NL-to-code generation, completion, translation, bug fixing, test gen-
eration, understanding, and security analysis. They catalog representative datasets, inputs/outputs,
languages, and evaluation metrics, and identify gaps in multi-round interactions, cross-phase
workflows, and tasks requiring coordinated capabilities. While the benchmark ecosystem is broad,
it remains unbalanced. A major missing area is code review: despite its central role in quality and
maintenance, there is no systematic SDLC-aware benchmark landscape for review tasks, limiting
fair comparison and realistic evaluation of LLM-based review assistance.
Wang et al. [99] reviewed 112 code-review studies from 2011 to 2019, examining their
methodologies and evaluation methods. They find evaluation-focused work dominates, while
solution proposals remain limited. Most studies rely on Gerrit/GitHub data, yet only about half
release replicable datasets, hindering comparability. The authors identify 457 metrics across sixteen
categories, showing fragmented definitions of review quality and behavior. Dataset choices are
also inconsistent, with most papers building their own. The study concludes that the field lacks a
standardized benchmark for code-review tasks.
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
5
Despite rapid growth in LLM-driven research on code review tasks—such as change un-
derstanding, comment generation, and automated code revision, there remains no comprehensive
benchmark for evaluating these capabilities. Therefore, there is a clear need for a domain-specific
survey of code-review benchmarks to identify existing resources, highlight their limitations, and
outline a roadmap for developing robust, standardized benchmarks in the future.
3
Methodology
We conduct our investigation around three research questions (RQ1–RQ3). Together, these questions
examine how code review approaches have been benchmarked and evaluated in both the Pre-LLM
and LLM eras, helping us reveal established practices, new trends, and gaps that still need to be
addressed.

RQ1: What tasks and evaluation strategies were used for code review studies in
the Pre-LLM era?
This RQ explores what were the common tasks and how such tasks were evaluated
in code review studies before the introduction of LLMs, when traditional machine-
learning or heuristic-based methods dominated. Understanding the tasks, datasets, and
evaluation metrics used in this earlier era helps us establish a baseline for how the
software engineering community historically investigated code review quality, defect
detection, reviewer support, and other core skills. By mapping these foundationaltasks
and their evaluation strategies, we can better identify what remains relevant in the LLM
era, what the limitations are that motivated the shift toward LLM-based techniques, and
where current benchmarks may still lack coverage.

RQ2: What task and evaluation strategies are used for code review studies in
the LLM era?
This question focuses on understanding how the tasks and their evaluation approaches
for code review have evolved in the LLM era. By examining the tasks and datasets
designed specifically for LLM-based systems, we identify what tasks the SE research
community is focusing on and expect these models to demonstrate their abilities. Study-
ing the tasks and their evaluation setups used in these works helps reveal how LLM
performance is currently measured, and where gaps or inconsistencies may still exist in
benchmarking code review.

RQ3: How have tasks and evaluation strategies in code review evolved from the
Pre-LLM to the LLM era?
The emergence of LLMs has fundamentally changed the landscape of code review
research, necessitating a critical examination of how task definitions, dataset character-
istics, and evaluation strategies have evolved. By systematically comparing the Pre-LLM
era with the LLM era, we can understand the trend when integrating LLMs in code
review. Crucially, understanding this transition allows us to identify the inherent limi-
tations of current LLM-era datasets, such as the “vanishing” of code review activities.
Identifying the evolution of code review research is crucial for understanding whether
the current direction of code review research is unintentionally overlooking important
ways of measuring how well a model performs code review tasks.
3.1
Paper Collection
To ensure a comprehensive and systematic collection of research on code review, we adopted a
rigorous multi-stage literature collection process. Our workflow consisted of four primary steps,
which are illustrated and briefly described in Figure 2.
, Vol. 1, No. 1, Article . Publication date: February 2026.

6
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
Fig. 2. Overview of the literature collection process.
(1) Terminology Summary. We began
by conducting preliminary readings and
compiling a set of keyword combinations
related to code review approaches. We
prepared 20 different search terms, in-
cluding: “code-review benchmarks,” “code
review automation,” “machine learning
for code review,” “LLM for code review,”
“AI-assisted review,” “pull-request review
datasets,” “review comment generation,”
“defect detection in code review,” “code
refinement benchmarks,” “review recom-
mendation,” “change-understanding tasks,”
“PR-level analysis,” “review prioritization,” “LLM agents for code review,” “review comment quality
evaluation,” “automated reviewer assignment,” and other related expressions that reflect various
code review tasks.
(2) Literature Retrieval. Using these predefined keywords, we conducted automated searches
across several major academic databases, such as IEEE Xplore, ACM Digital Library, Elsevier
ScienceDirect, ACL Anthology, arXiv, Google Scholar, and SpringerLink. These sources were
selected because they provide broad coverage of core venues in software engineering research
as well as emerging work on AI-assisted benchmarking and code review automation. To reflect
the evolution of modern code review and related benchmark research, we limited our search to
publications from 2015 onward. In total, our search retrieved 160 publications related to code review
from Jan. 2015 to Dec. 2025.
(3) Literature Screening. We then applied a structured filtering process to identify relevant and
high-quality studies. After screening titles, abstracts, and full texts, 77 papers met our inclusion
criteria:

The study must address code review, review automation, or code review related tasks
(e.g., comment generation, review recommendation, defect detection during review).

Priority is given to papers accepted by leading SE and AI venues (e.g., TOSEM, TSE,
EMSE, ICSE, FSE, ASE, ISSTA, ACL, etc.).

For arXiv preprints, we included only papers that had at least one citation to ensure a
minimum level of relevance and scholarly use.

Papers must disclose their dataset details or evaluation approaches. The papers without
such details were filtered from the pool.
(4) Snowball Expansion. Recognizing that keyword searches alone may miss important contri-
butions, we conducted both backward and forward snowballing from our selected papers. This
step allowed us to incorporate an additional 25 studies, capturing influential work that may not
explicitly match our initial keyword set. This process ensured breadth, completeness, and coverage
of relevant benchmark- and review-related tasks.
In total, we collected 99 papers in our final dataset. Of these, 58 belonged to the pre–LLM
era and 41 papers are done in the LLM era. To ensure reliability, each paper underwent a quality
assessment by the authors to verify that the included studies were methodologically sound and
highly relevant to the scope of our survey.
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
7
3.2
Statistics of Papers
Table 1 summarizes the distribution of all papers included in our survey across different publication
venues and lists the corresponding literature for each venue. The table highlights how research
on code review automation and related topics is spread across major conferences and journals,
with arXiv, ICSE, SANER, FSE, MSR, and other flagship venues contributing the largest number of
publications. This distribution helps illustrate where the research community has concentrated
its efforts, which venues are most active in publishing work on code review related technologies,
and how diverse the underlying literature base is. By organizing the papers in this way, the table
provides a clear overview of the breadth of prior work and the venues that have shaped current
developments in this area.
Domain
Venue
Papers No. Literature
Software Engineering (SE)
ICSE
12
Cihan et al. [14]; Tufano et al. [93]; Hu et al. [35]; Tufano et al. [94]; Guo et al. [28];
Tufano et al. [95]; Thongtanunam et al. [91]; Barnett et al. [5]; Egelman et al. [18];
Qiu et al. [72]; Sawant and Sengamedu [78]; Saini and Britto [75]
FSE
6
Sun et al. [88]; Li et al. [50]; Yang et al. [107]; Hong et al. [33]; Chen et al. [10];
Shan et al. [81]
ASE
2
Wang et al. [103]; Ahmed et al. [1]
ICSME
4
Wen et al. [105]; Hanam et al. [31]; Ebert et al. [17]; Shuvo et al. [84]
SANER
6
Sghaier and Sahraoui [79]; Hong et al. [34]; Lin and Thongtanunam [53]; Pornpr-
asit et al. [71]; Chatley and Jones [9]; Siow et al. [85]
MSR
7
Jaoua et al. [40]; Sghaier et al. [80]; Liu et al. [55]; Uchôa et al. [96]; Tao and
Kim [90]; Rahman et al. [73]; Wang et al. [104]
ICPC
2
Brito and Valente [7]; Kanda et al. [42]
ESEM
3
Soltanifar et al. [86]; Sarker et al. [76]; Rahman et al. [74]
APSEC
1
Wang et al. [98]
ISEC
1
Kapur et al. [44]
ICSSP
1
Mukhtarov et al. [65]
SWQD
1
Ochodek et al. [67]
SCAM
1
Khelifi et al. [45]
VISSOFT
1
Balcı et al. [4]
IWSC
1
Ueda et al. [97]
FORGE
1
Işık et al. [39]
ISSRE
2
Lu et al. [60]; Hijazi et al. [32]
FASE
1
Lu et al. [59]
ISCO
1
Lal and Pahwa [48]
EMSE
4
Chouchen and Ouni [13]; Fan et al. [19]; Zhao et al. [112]; Huang et al. [36]
TOSEM
3
Lin et al. [54]; Sarker et al. [77]; Kim et al. [46]
IST
3
Pornprasit and Tantithamthavorn [70]; Liu et al. [57]; Islam et al. [38]
JSS
1
Fregnan et al. [21]
SPE
1
Sharma and Sodhi [82]
AI/ML
ICML
1
Lu et al. [58]
AAAI
1
Shi et al. [83]
ACL
2
Yan et al. [106]; Lin et al. [52]
EMNLP
1
Tang et al. [89]
Data Mining & Data Science
KDD
1
Gupta [29]
HCI
VL/HCC
1
Ge et al. [22]
Others
ARXIV
17
Guo et al. [26]; Maddila et al. [62]; Chervyakov et al. [12]; Zeng et al. [111];
Trakoolgerntong et al. [92]; Icoz and Biricik [37]; Guo et al. [27]; Yu et al. [109];
Goldman et al. [23]; Jiang et al. [41]; Zhao et al. [113]; Yu et al. [110]; Kumar and
Chimalakonda [47]; Naik et al. [66]; Kapadnis et al. [43]; Liu et al. [56]; Haider
et al. [30]
JCST
2
Guo et al. [24]; Guo et al. [25]
COMPSAC
2
Ayinala et al. [2]; Chen et al. [11]
EIT
1
Fish et al. [20]
FIE
1
Crandall et al. [15]
EAAI
1
Cao et al. [8]
IEEE ACCESS
1
Lee and Joe [49]
FEDCSIS
1
Madera and Tomoń [63]
SOFSEM
1
Luna Freire et al. [61]
SENSORS
1
Yin et al. [108]
Table 1. Distribution of papers per venue, grouped by domain and corresponding literature.
, Vol. 1, No. 1, Article . Publication date: February 2026.

8
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
Fig. 3. Venue–year distribution of the analyzed studies.
Figure 3 presents a year-by-year distribution of code-review benchmark publications across a wide
range of venues, illustrating both the growth of the field and how research activities have spread
across different communities. Each bar corresponds to a specific year, and the colored segments
inside the bars represent the number of publications at each venue in that year. The height of a bar
indicates how many benchmark-related papers were published overall during that year, while the
variation in colors shows how many different conferences, journals, and workshops contributed to
this research.
From left to right, the figure shows that early studies (2015–2018) were limited, with only a small
number of papers each year and contributions mainly from traditional software engineering venues
such as MSR, ICPC, ICSME, and ICSSP. As the time progresses, particularly starting around 2020,
the number of papers increases and venues become more diverse, indicating that more papers are
being published and that additional venues are becoming involved. This reflects the gradual shift
toward machine-learning-based techniques and the increasing interest in establishing standardized
datasets and evaluation methods for code-review tasks.
The period from 2022 onward shows a clear acceleration. These years display a broader mix of
venues, including software engineering conferences (ICSE, FSE, and SANER), journals (EMSE and
JSS), AI/NLP venues (EMNLP and ACL), and applied or interdisciplinary outlets. This suggests that
code review benchmarking is no longer confined to a small set of software engineering researchers
but is attracting attention from machine learning and natural language processing communities as
well.
The most dramatic change appears in 2025, where the bar rises sharply compared to all previous
years. A significant portion of this increase comes from arXiv publications, indicating a surge in
preprint activity driven by rapid experimentation with LLM-based methods. At the same time, es-
tablished venues still contribute smaller but steady numbers of papers. This combination highlights
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
9
a shift in publication culture: researchers are disseminating benchmark datasets and evaluation
studies faster and at larger scale, often before formal peer review.
3.3
Constructing the Taxonomy of Code Review-related Tasks
To construct the taxonomy of code review-related tasks, the first two authors perform open coding.
The process follows a three-stage process:

Phase 1: Open/Initial Coding In this phase, the first two authors first conduct the
coding independently. They treat each paper as a primary data source. For every paper,
they perform a manual inspection of the abstract and introduction to identify the specific
problem or task the paper addresses. A “low-level label” that captures the essence of the
paper’s contribution was assigned. For instance, if a paper introduces an algorithm to
generate comments for a commit, a label “Review Comment Generation” was assigned.
After the independent phase, the authors met to compare their labels. Any differences
in the interpretation of a paper’s task were discussed in detail. The inter-agreement
Cohen’s Kappa is 0.81. During the labeling, we also summarize the input/output, and
evaluation setup (e.g., evaluation metrics) for the tasks.

Phase 2: Selective Coding Once all papers are labeled, the first two authors look for
semantic similarities and patterns among the raw codes, and begin to cluster these
individual labels into broader categories following a bottom-up manner. This creates
the High-level Tasks of our taxonomy.
We construct the following high-level tasks:

Review Prioritization / Selection refers to the task of deciding which code changes
should be reviewed and in what order so that limited review effort is used efficiently. It
is commonly formulated as a prediction or ranking problem in both Pre-LLM and LLM-
based approaches. Given information such as the patch or diff content (e.g., changed files,
size, and complexity), developer and reviewer history, review comments, and available
test or CI signals, a model estimates the urgency or importance of each change. The
output may be a binary or multi-class label (e.g., review needed, merged, or abandoned),
a numerical estimate (such as the expected number of review rounds or delay), or a
ranked list of patches indicating which changes should receive attention first.

Change Understanding / Analysis encompasses tasks that support reviewers in
interpreting the intent, impact, and structure of code changes and review requests.
In both Pre-LLM and LLM eras, these tasks are typically formulated as prediction,
classification, or ranking problems, where systems analyze information such as code
diffs or patches, surrounding code context, commit messages, reviewer comments, and
basic project or reviewer metadata. The outputs may be labels or scores (e.g., high-risk or
high-effort), ranked or highlighted change regions, concise structured summaries, or, in
reviewer-intent–focused settings, a change-type label indicating whether the requested
modification involves adding, deleting, or modifying code.

Peer Review refers to a broad class of automated tasks designed to support or emulate
human code review by analyzing code changes and related review information. Across
both Pre-LLM and LLM eras, these tasks are formulated as prediction, ranking, recom-
mendation, classification, localization, or text-generation problems, using inputs such as
code changes (patches or diffs at file, method, commit, or pull-request level), surrounding
code context, review or issue text, review metadata, process information, and reviewer
history. Depending on the formulation, the outputs include defect or rework labels,
relevance scores or ranked recommendations (e.g., similar past examples or useful review
, Vol. 1, No. 1, Article . Publication date: February 2026.

10
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
comments), identified code locations requiring attention, or natural-language review
feedback that reflects a reviewer’s judgment.

Review Assessment / Analysis refers to a set of tasks that examine the quality, tone,
and effectiveness of review-related artifacts and interactions in code review. Across
both Pre-LLM and LLM eras, these tasks are typically formulated as classification or
scoring problems that analyze inputs such as review comment text, code changes or diffs,
commit messages, pull request descriptions, linked issues, and basic reviewer, author,
or process metadata. The outputs are labels or scores that characterize aspects such
as usefulness, clarity, relevance, consistency, toxicity, or overall review adequacy, and
in some settings include indicators of whether feedback was later addressed by code
updates.

Code Refinement refers to tasks that support the code revision stage after review
by helping developers improve an existing change based on review outcomes. In the
surveyed benchmarks, this task is mainly formulated as a code generation or code
transformation problem, where the goal is to generate a revised version of a method or
code fragment that matches the reviewer-approved implementation, and in some cases as
a classification problem that distinguishes refactoring changes from behavior-changing
edits or identifies consistent naming. The typical inputs include the original code before
revision, the revised code used as supervision, code diffs or diff markers, and sometimes
review comments or basic pull-request metadata. The outputs are revised code snippets,
refactored methods, or simple labels indicating the type of change or recommended
names.
Table 2 presents the resultant taxonomy of code review tasks, including high-level task,
sub-task, and their definitions. We end up with five high-level tasks, which align with the code
review process as shown in Figure 1.
4
RQ1: What tasks and evaluation strategies were used for code review studies in the
Pre-LLM era
In this research question, we examine which tasks and datasets have been used to evaluate code
review systems in the Pre-LLM era. We also analyze how these systems have been measured,
focusing on the evaluation metrics applied and the specific code review skills assessed. Figure 4,
presents an overview of our taxonomy from the Pre-LLM era literature, illustrating how the 58
papers map out the landscape of code review research organized across five major high-level tasks
and their sub-tasks.
4.1
Review Prioritization / Selection
4.1.1
Change Qality Prediction: Change Quality Prediction helps decide the code change quality
and predict whether a change would be accepted or not. Studies [13, 19, 38, 83, 112] constructed
dataset to predict if a change would be accepted (merged) at the commit or pull request level,
leveraging various information (e.g., code diffs, patch metadata, file history, developer experience,
collaboration networks, and text), and evaluate performance using classification metrics such as
AUC and F-measure. Table 3 presents the details of each dataset and its evaluation metrics.
4.1.2
Review Prioritization: Review prioritization helps decide which code changes should be
reviewed first and how much effort is needed. Prior studies look at code changes and past history
to find risky or important patches, and help reviewers focus on urgent work. Table 4 summa-
rizes Pre-LLM research on review prioritization. For example, two studies [34, 75] constructed
datasets for review prioritization at two granularities, i.e., commit level and line level, respectively.
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
11
High-level task
Sub-task
Description
Review
Prioritiza-
tion / Selection
Review Prioritization
Determining the urgency and necessity of a review to optimize the
engineering team’s bandwidth.
Change Quality Prediction
Predicting the quality of code changes, whether they require review or
would be accepted or not.
Change Understand-
ing and Analysis
Relevant Change Identifica-
tion
Contextualizing the current changes by tracing their relationship to
historical versions or related files in the repository.
Impact Analysis
Determining the “blast radius” of a change by identifying which down-
stream components or dependencies might impacted.
Change Decomposition
Breaking down a monolithic Pull Request into smaller, logical “atomic”
commits/slicers to reduce reviewer cognitive load.
Visualization
Utilizing graphs or diff-enhancements to provide a mental map of how
the code structure has evolved.
Peer Review
Issue Localization
Pinpointing the specific lines of code or architectural layers where a
reported defect or logical error originates.
Review Comment Generation Generating natural-language feedback that explains code issues to the
author.
Issue Labeling / Classification Categorize review comments into specific issue types (e.g., Bug, Read-
ability, Design, etc.).
Refactoring Identification
Spotting “code smells” or technical debt that should be addressed.
Defective Change Prediction Flagging code changes that contain defects.
Security Detection
Scanning the proposed changes for vulnerabilities such as injection
flaws or credential leaks.
Review Assessment
and Analysis
Review Quality Evaluation
Assessing the utility and constructiveness of a reviewer’s feedback.
Comment–Code Compliance Examining if a review comment accurately matches the code it refers
to.
Review Summarization
Condensing complex discussions into a concise status report.
Sentiment / Toxicity Analysis Monitoring the emotional tone of review interactions.
Code Refinement
Code Revision
Implementing direct fixes to the logic based on the specific issues raised.
Code Refactoring (Rework)
Restructuring the implementation to improve long-term maintainabil-
ity.
Table 2. The definition of each code review-related task.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Shi
et al. [83]
Change acceptance
prediction
Open
Source
6
Java
Chunk level
35,640
changed
hunks
Classification Metrics: F-measure; AUC; t-test;
Mann–Whitney U.
Fan
et al. [19]
Change acceptance
prediction
Open
Source
3
Ruby, Go and 7 oth-
ers
Commit level 166,215 Code
Chnages
Classification Metrics: AUC; ER@K; Precision;
Recall; F-measure
Islam
et al. [38]
Change acceptance
prediction
Open
Source
3
N/A
Commit level 146,612 Code
Changes
Classification Metrics: AUC; ER@K; Precision;
Recall; F-measure
Zhao
et al. [112]
Change acceptance
prediction/ranking
Open
Source
74
Java
PR level
126,180 PRs
Ranking Metrics: NDCG@k
Chouchen
and
Ouni [13]
Change acceptance
prediction
Open
Source
3
N/A
Commit level 146,612
code review
requests
Classification Metrics: AUC
Ranking: Scott–Knott ESD
Table 3. Datasets and Evaluation Metrics for Change Quality Prediction in Pre-LLM era.
Huang et al. [36] constructed a dataset for review round prediction on 19,964 patches, using both
classification metrics and regression metrics (e.g., MAE) to estimate review effort. Finally, two
studies [13, 104] constructed datasets for predicting whether a change would require high review
effort at the chunk level, highlighting the importance of balancing merge likelihood and inspection
effort under limited review capacity. Most of the datasets are sourced from open source projects,
and only one is sourced from private project. In addition, all the datasets are at the commit or PR
level.
, Vol. 1, No. 1, Article . Publication date: February 2026.

12
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
Fig. 4. Distribution of literature across different code review-related tasks in Pre-LLM era.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Saini
and
Britto [75]
Review
prioritiza-
tion
Private
Projects
1
N/A
Commit level N/A
Classification Metrics: AUC, TPR/FPR
Huang
et al.. [36]
Review round count
prediction
Open
Source
3
Java and Python
Commit level 19,964
Patches
Classification Metrics: F-measure
Regression metrics: MAE, R2
Hong
et al. [34]
Review
prioritiza-
tion
Open
Source
3
Java and Python
Chunk level
19,964
Patches
Classification Metrics: Precision; Recall; Accu-
racy; F-measure
Chouchen
and
Ouni [13]
Review effort predic-
tion
Open
Source
4
N/A
Commit level 146,612
code review
requests
Classification Metrics: ACC; Popt; Recall
Ranking: Scott–Knott ESD
Wang
et al. [104]
Review effort predic-
tion
Private
Projects
4
C#, C, Java, Python Chunk level
More
than
100000
changes
Classification Metrics: Precision; Recall; F-
measure; AUC;
Table 4. Datasets and Evaluation Metrics for Review Prioritization in Pre-LLM era.
4.2
Change Understanding and Analysis
4.2.1
Relevant Change Identification: Change Identification aims to identify what changed in the
code and understand those changes, even across files or versions. It checks for similar, missing,
or inconsistent edits and highlights changes that may need more review. Table 5 presents the
dataset constructed to evaluate change identification in Pre-LLM era. Ayinala et al. [2] studied
inconsistent clone edit detection as a review support task, where the goal is to identify inconsistent
updates among cloned methods; their approach was evaluated at the method level. Chen et al. [11]
focused on detecting incomplete edits in similar code, aiming to help reviewers find missed updates
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
13
across related files. Finally, Fish et al. [20] studied a clone evolution detection task that tracks clone
changes across revisions to support review and maintenance. The dataset in this category usually
focuses on low level, such as file, method, or chunk. All the datasets are in Java and sourced from
open source projects.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Ayinala
et al. [2]
Inconsistent
clone
edit detection
Open
Source
6
Java
Method level 14
 842 clone
instances
Classification Metrics: Precision; Recall; F-
measure.
Chen
et al. [11]
Detect incomplete
edits in similar code
Open
Source
4
Java
File level
2.3
million
Java files
Classification Metrics: Precision; Recall; F-
measure.
Fish
et al. [20]
Clone evolution de-
tection
Open
Source
3
Java
Chunk level 657 Revision
pair
Classification Metrics: Accuracy.
Table 5. Datasets and Evaluation Metrics for Relevant Change Identification in Pre-LLM era.
4.2.2
Change Decomposition. Change Decomposition breaks a large code change into smaller,
meaningful, independent parts. It groups related edits so reviewers can understand complex changes
and prioritize the most impactful changes first that may need deeper review. Table 6 presents
decomposition benchmarks that investigate how complex code changes can be structured to better
support human code review. Most of the studies [5, 24, 61, 90] focus on decomposing a commits
that contain multiple concerns into smaller logically related change slices to improve reviewer
understanding and efficiency. One study [103] looked into tangled commit decomposition at the
commit level, where the non-buggy changes are separated from buggy changes. We observe that
the datasets are all open source projects, focus on Java and C#, typically at the commit or PR level.
Another study [22] aimed to automatically identify and separate behavior preserving refactoring
code changes from other code changes so reviewers can focus on risky modifications.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Barnett
et al. [5]
Independent change
decomposition
Private
Projects
13
C#
Commit level 1,000
code
changes
Classification Metrics: Accuracy.
User Study: Qualitative Interview Feedback.
Tao
and
Kim [90]
Independent change
decomposition
Open
Source
4
Java
Commit level 453 revisions Text-Matching/Generation
Metrics:
Exact-
Match.
Classification Metrics: Accuracy.
User Study: Review Time.
Luna Freire
et al. [61]
Independent change
decomposition
Open
Source
10
Java
PR level
1,000 pull re-
quests
Classification: False Positives and Error Types.
Guo
et al. [24]
Independent change
decomposition
Open
Source
4
Java
Chunk level
28 composite
changes
Classification Metrics: Precision; Recall; F-
measure.
Wang
et al. [103]
Tangled commit de-
composition
Open
Source
7
Java
Commit level 50
tangled
commits
Classification Metrics: Rand Index; MAP; MRR.
Ge
et al. [22]
Refactoring change
decomposition
Open
Source
2
Java
Commit level 3,100
commits
Classification Metrics: Precision; Recall.
Table 6. Datasets and Evaluation Metrics for Change Decomposition in Pre-LLM era.
4.2.3
Impact Analysis: Impact Analysis estimates how much a code change affects the system
by predicting whether changes are risky or impactful and by identifying which parts of the code,
design, or build may require closer review. Table 7 presents datasets for Impact Analysis proposed
in the Pre-LLM era. Uchôa et al. [96] introduced a revision-level design impact classification task in
an open-source setting, where each PR revision is classified based on whether it negatively impacts
software design. The task is formulated at the PR level and evaluated on 57,498 reviewed code
changes from seven Java projects using classification metrics including Precision, Recall, F-measure,
, Vol. 1, No. 1, Article . Publication date: February 2026.

14
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
and AUC, demonstrating that review and change features can effectively flag design-impactful
revisions for reviewer attention. Wen et al. [105] aimed to identify the build deliverables that are
affected by a given patch. Hanam et al. [31] constructed a dataset where impacted code regions for
a given commit are annotated. The datasets for this task are typically at the PR or commit level. The
datasets span various programming languages, although they are sourced from a limited number
of projects.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Uchôa
et al. [96]
Design impact iden-
tification
Open
Source
7
Java
PR level
57,498
re-
viewed code
changes
Classification Metrics: Precision; Recall; F-
measure; AUC.
Wen
et al. [105]
Build
deliverables
impact analysis
Private
Projects
1
C, C++, Java and C# PR level
Over 10 mil-
lion lines of
code
User Study: speed and accuracy of identifying
the set of deliverables.
Hanam
et al. [31]
Semantic change im-
pact annotation
open
Source
3
JavaScript (Node.js) Commit level Mining study
over  2,000
commits
Classification Metrics: Reductions in False Posi-
tives; Reductions in Change-Impact Set Size.
User Study: Search Time; Success Rate
Table 7. Datasets and Evaluation Metrics for Impact Analysis in Pre-LLM era.
4.2.4
Visualization: Visualization helps reviewers understand code changes using clear and easy-
to-understand visuals. It shows risk, effort, structure, or behavior so reviewers can quickly spot
important, risky, or complex changes. Table 8 discusses representative visualization-based bench-
marks proposed in the Pre-LLM era to support change understanding during code review. Wang
et al. [98] studied review risk, effort, and impact analysis at the commit level in open-source Java
systems, where the task is to help reviewers assess how difficult and risky a commit is to review
by analyzing evolutionary coupling information across 10 projects with 41,855 commit counts.
Their work evaluates reviewer-oriented skills related to effort estimation and impact awareness
using behavioral metrics: Effort, Risk, and Impact, complemented by visualization methods such
as Clustering (KMeans), Spider Charts, Coupling Charts, and Density Maps to summarize change
characteristics. Balcı et al. [4] focused on chunk-level risk visualization for code review, formulating
the task as understanding the risk and nature of fine-grained change chunks in Java code, evaluated
on one open-source project using 7 predefined change scenarios; the study emphasizes reviewer
comprehension skills and reports results through a user study using Risk Percentages per Commit
along with Readability Ratings, Understandability Ratings, and Practicality Ratings. Kanda et al. [42]
presented an execution-trace–augmented code diff visualization task at the chunk level, where
reviewers inspect both code differences and execution-trace differences to better understand behav-
ioral changes; the approach is demonstrated on a single Java project using Math-57 from Defects4J,
and no quantitative evaluation metrics are reported, relying instead of illustrative analysis through
the visualization. Finally, Fregnan et al. [21] proposed a graph-based change-set visualization task
at the PR level to help reviewers navigate and understand large PRs in Java projects, evaluated on
138,452 PRs from 208 private projects; their task is framed as supporting reviewer decision-making
and change comprehension, with results reported using classification metrics such as FP and success
rate.
4.3
Peer Review
4.3.1
Review Coment Generation: Review Comment Recommendation or Generation helps re-
viewers by suggesting helpful comments for new code changes, leveraging past reviews to identify
recurring issues and automatically produce clear and useful feedback. Table 9 discusses representa-
tive benchmarks from the Pre-LLM era that study automated review comment recommendation,
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
15
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Wang
et al. [98]
Review risk, effort,
impact analysis and
visualization
Open
Source
10
Java
Commit level 41,855 com-
mit counts
Behavioral Metrics: Effort (C), Risk, Impact
Visualization Methods: Clustering (KMeans), Spi-
der Charts, Coupling Charts, Density Maps
Balcı
et al. [4]
Chunk-level risk vi-
sualization
Open
Source
1
Java
Chunk level
7 predefined
change
scenarios
User study: Risk Percentages per Commit; Read-
ability Ratings; Understandability Ratings; Prac-
ticality Ratings
Kanda
et al. [42]
Execution-
trace–augmented
code diff visualiza-
tion
Open
Source
1
Java
Chunk level
Math-57
from
De-
fects4J
N/A
Fregnan
et al. [21]
Graph-based
change-set
visual-
ization
Private
Projects
208
Java
PR level
138,452 pull
requests
Classification: FP; Success rate.
Table 8. Datasets and Evaluation Metrics for Visualization in Pre-LLM era.
generation, and analysis tasks in code review workflows by modeling reviewer feedback using
historical data from both private and open-source projects. Chatley and Jones [9] presented an
automated code review comment generation task in which the system learns from repository
history and pull request activity to automatically produce natural language review comments for
new changes, with a particular focus on maintainability and architectural issues. The majority of
prior studies in the Pre-LLM era focus on review comment recommendation [29, 33, 84, 85], which
is typically formulated as an information retrieval problem, where similar changes are retrieved
and their associated review comments are reused for a given code change. Ochodek et al. [67]
automatically classify code-review comments by their focus (e.g., design, style, logic, data, API, I/O,
documentation, configuration, build/install, patch review), and constructed a dataset. The datasets
for review recommendation/generation could span board granularity, from chunk level to PR level.
The programming languages primarily focus on popular ones, e.g., Java, Python, C#. The dataset
size ranges from 8 to 11,289.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Chatley and
Jones [9]
NL
Review
com-
ment generation
Private
Projects
21
Ruby
Chunk level
Almost 150k
LOC
Ruby
repo
User study:Comment Resolve Rate.
Siow
et al. [85]
Review
comment
recommendation
Open
Source
19
Java
PR level
57,260 h code
change
Classification Metrics: Recall; MRR.
Gupta [29]
Review
comment
recommendation
Private
Projects
208
C#
Chunk level
22,435 com-
pleted PRs
Classification Metrics: AUC; MRR; Recall
User Study: Acceptance Rate.
Hong
et al. [33]
Review
comment
recommendation
Open
Source
11,289
Java
Method level 151,019
h changed
method
Text-Matching / Generation Metrics: Perfect Pre-
diction; BLEU-4.
Shuvo
et al. [84]
Review
comment
generation
Open
Source
8
Python and Java
Chunk level
56,068
review
comments
Text-Matching / Generation Metrics: BLEU; Exact
Match; Sentence-BERT Cosine Similarity.
Ochodek
et al. [67]
Review
comment
classification
Open
Source
3
C
Chunk level
2,672
comments
Classification Metrics: Accuracy; MCC; AUC.
Table 9. Datasets and Evaluation Metrics for Review Comment Generation in Pre-LLM era.
4.3.2
Defective Change Prediction: Defective Change Prediction estimates whether a code change
is likely to contain bugs. It helps reviewers find risky changes early and decide which ones need
careful review first. Table 10 presents Pre-LLM datasets that focus on predicting defective or risky
changes to support code review decision making. Most of these studies [44, 57, 82, 86] focus on
defect prediction and construct corresponding datasets, where the goal is to predict whether a code
, Vol. 1, No. 1, Article . Publication date: February 2026.

16
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
change will later introduce a defect. For instance, Sharma and Sodhi [82] proposed an approach
to predict defective code changes by retrieving similar code from Stack Overflow and assigning
defective labels to source files that exhibit high similarity. Therefore, they collected millions of Stack
Overflow posts and a curated set of 370 GitHub files. Madera and Tomoń [63] introduced a rework
risk prediction task in an industrial setting that classifies file changes as Correct or Rework in
order to help reviewers and quality assurance teams focus on high-risk modifications. The dataset
in this task spans various programming languages, such as Java, Python, C#. It spans both open
source and private projects. The granularity primarily focuses on file-level prediction, and only
one focuses on commit level.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Soltanifar
et al. [86]
Defect prediction
Open
Source
1
N/A
Commit level 4,750 review
requests
Classification Metrics: Recall; F-measure; AUC.
Statistical / Behavioral Analysis Metrics:
Sharma and
Sodhi [82]
Defect prediction
Open
Source
3
Java, JavaScript and
5 others
File level
370 GitHub
files
Classification
Metrics:
Matching
Ratio;
F-Measure
Kapur
et al. [44]
Defect prediction
Open
Source
2
Java, Python and 3
others
File level
188,200
GitHub files
Classification Metrics: Accuracy; F-Measure; Pre-
cision; Recall.
Liu
et al. [57]
Defect prediction
Private
Projects
14
N/A
File level
2
industrial
systems; 12
programs
Classification Metrics: Accuracy; False positives;
Defect-detection rate
Madera and
Tomoń [63]
Risk prediction
Private
Projects
1
Java and .NET
File level
237,128
file-change
records
Classification Metrics: AUC; Accuracy; Preci-
sion; Recall; PPV (Positive Predictive Value) ;NPV
(Negative Predictive Value)
Table 10. Datasets and Evaluation Metrics for Defective Change Prediction in Pre-LLM era.
4.3.3
Refactoring Identification: Brito and Valente [
7] focused on refactoring detection in real
world pull request review scenarios at PR level using four private Go projects where the dataset
contains 325 pull requests 374 commits approximately 754K LOC and 685 refactoring instances
and the evaluation relies on Classification Metrics Precision and Recall to measure how accurately
refactoring can be detected and surfaced to reviewers during code review. Kim et al. [46] studied
method name consistency checking in large-scale open source projects. The datasets are sourced
from both open source and private projects, and covered various programming languages. The
granularity mainly focus on the chunk level, and size of the datasets from 85 to 17,000, depending
on whether the evaluation is automated or manual.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Kim
et al. [46]
Refactoring
detec-
tion in pull requests
Open
Source
66
Java
Method level 400 reviewed
method
Classification Metrics: Precision; Recall; F-
measure; Exact Match Accuracy (EMAcc).
Brito
and
Valente [7]
Inconsistent method
name detection
Private
Projects
4
Go
PR level
325
pull
requests, 374
commits,
 754K
LOC,
685
refactoring
instances.
Classification Metrics: Precision; Recall.
Kim
et al. [46]
Inconsistent method
name detection
Open
Source
66
Java
Method level 13,537 meth-
ods
Classification Metrics: Precision; Recall; F-
measure; Accuracy.
Table 11. Datasets and Evaluation Metrics for Refactoring Identification in Pre-LLM era.
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
17
4.4
Review Assessment and Analysis:
4.4.1
Review Qality Evaluation: Review Quality measures how good a code review is. It checks
whether comments are useful, clear, and relevant to the code, and whether the review properly
addresses important changes. Table 12 discusses the datasets and evaluation setup used in previous
work in this line. Two studies [17, 74] constructed datasets for clarity classification by predicting
whether a review comment is clear or unclear and ranking similar examples to support explanation
and reviewer understanding. Another group of studies [25, 73] constructed a dataset for review
usefulness prediction, where their goal is to predict whether a review comment is useful or not for
developers. We also observe studies to investigate quality from various aspects. Yang et al. [107]
aim to evaluate a comment from four attributes, i.e., emotion, question, evaluation, and suggestion.
For this purpose, they manually curated a dataset, that contains comments and the annotations in
four attributes. Two studies [32, 65] constructs datasets for low quality review classification task
where the code is reviewed without sufficient attention from reviewers. The datasets coverage
various languages such as Java, Python, C, Kotlin, etc. Most of the studies focus on the chunk-level,
and the dataset ranges from small size of 84 code reviews from programming course to large size
of 140,000 code review from open source repositories.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Rahman
et al. [74]
Clarity classification Open
Source
3
N/A
Commit level 3,722 code re-
views
Classification Metrics: Accuracy, Precision, Re-
call, F-measure
Ranking Evaluation: Top-5 accuracy / Top-5 re-
call
Ebert
et al. [17]
Clarity classification Open
Source
N/A
N/A
Chunk level
140,006 code
reviews
Classification Metrics: Precision; Recall; F-
measure.
Rahman
et al. [73]
Usefulness
predic-
tion
Private
Projects
4
Python
Chunk level
1,116 inline
comments
Classification Metrics: Precision; Recall; F-
measure; Accuracy; ROC; Precision; Recall
Guo
et al. [25]
Usefulness
predic-
tion
Private
Projects
7
Java, Kotlin and 4
others
Chunk level
9,477 code re-
views
Classification Metrics: Accuracy; Precision; Re-
call; F-measure
Yang
et al. [107]
Review quality at-
tribute classification
Private
Projects
N/A
N/A
Chunk level
17,000
review
comments
Classification Metrics: Hamming Loss; Subset 0/1
Loss; F-measure; AUC; Precision; Recall
Hijazi
et al. [32]
Low-quality review
detection
Progra-
mming
course /
compe-
tition
4
C
Chunk level
84 code re-
views
Classification Metrics: Accuracy; Precision; Re-
call; F-measure
Mukhtarov
et al. [65]
Low-quality review
detection
Private
Projects
N/A
JavaScript
Chunk level
N/A
User Study: Likert ratings
Table 12. Datasets and Evaluation Metrics for Review Quality Evaluation in Pre-LLM era.
4.4.2
Coment-Code Compliance: Table 13 presents benchmarks for code compliance tasks in the
pre LLM era focusing on how code review systems assess whether code adheres to written develop-
ment or security policies. Sawant and Sengamedu [78] presented policy based code compliance
classification as a supervised learning task where the system receives a natural language policy
description together with an isolated Java code snippet and predicts whether the snippet is not
compliant or irrelevant with respect to the policy.
4.4.3
Sentiment/Toxicity Analysis. Table 14 presents benchmarks that focus on reviewers’ emotions
attitudes and interpersonal signals expressed through code review comments in the Pre-LLM era.
Ahmed et al. [1] presented the task of review comment sentiment classification and constructed
the corresponding dataset. Egelman et al. [18] addressed pushback detection, which aims to detect
, Vol. 1, No. 1, Article . Publication date: February 2026.

18
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Sawant
and
Sen-
gamedu [78]
Policy-based
code
compliance
classification
Open
Source
&
Private
Projects
N/A
Java
Chunk level
231 policies;
28000
code
examples.
Classification Metrics: Accuracy; MRR.
Table 13. Dataset for Comment-Code Compliance in Pre-LLM era.
the reviews that lead to developers’ negative feelings, and eventually cause pushback. Another two
studies [76, 77] proposed approaches to detect toxic review comments and construct corresponding
datasets. Lastly, citeauthorNP39 [72] constructed a dataset for detecting interpersonal conflict across
code review and issue discussion threads. The datasets in this task are typically at a high-level, i.e.,
commit/PR level, and are sourced from both open source and private projects.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Ahmed
et al. [1]
Review
comment
sentiment classifica-
tion
Open
Source
20
N/A
Commit level 2,000 labeled
comments
Classification Metrics: Precision; Recall; F-
measure; Accuracy.
Egelman
et al. [18]
Pushback (negative
feeling) detection
Private
Projects
N/A
N/A
PR level
1,250 change
requests
Classification Metrics: Precision; Recall.
Sarker
et al. [77]
Review
comment
toxicity
classifica-
tion
Open
Source
5
N/A
Commit level 19,651
la-
beled review
comments
Classification Metrics: Accuracy; Precision; Re-
call; F-measure; confusion matrix.
Sarker
et al. [76]
Review
comment
toxicity
classifica-
tion
Open
Source
4
N/A
Commit level 19,651
review
comment
Classification Metrics: Precision; Recall; F-
measure; span-matching metric.
Qiu
et al. [72]
Interpersonal
conflict detection
Private
Projects
4
N/A
Commit level 2,372 labeled
discussions
Classification Metrics: Precision–Recall; AUC;
Precision; Recall; F-measure.
Table 14. Datasets and Evaluation Metrics for Sentiment/Toxicity Analysis in the Pre-LLM era. Note that the
language in this task is natural language.
4.5
Code Refinement
4.5.1
Code Revision: This task uses the original code and review feedback to generate or suggest a
revised version or fixed patterns that match what reviewers expect. Table 15 presents Pre-LLM
benchmarks that study automated code revision and suggestion skills in the context of code review.
The papers related to this task can be grouped into two families: 1) revision pattern suggestion, and
2) revised code generation. In the first family, the study [97] extracted the code revision patterns
from similar patches and suggested relevant patterns when a new review comment is received. In
the second family, the studies [8, 53, 71, 91, 95, 108] typically formulated the problem as a code
generation task, where the model is trained to generate the revised code given a piece of code
and its review comment. For instance, Yin et al. [108] proposed an approach to generate revised
Java methods from original code and optionally review comments. As we can observe, almost all
datasets are at low level, i.e., method- and chunk-level. Only two popular programming languages,
Java and Python, are tested and all datasets are sourced from open source projects.
4.5.2
Code Refactoring (Rework): Table 16 presents Pre-LLM benchmarks related to Code Refac-
toring. Kim et al. [46] investigated refactoring related to method naming recommendation where
models learn reviewer preferred method names from reviewed code and comments, and they
curated 400 reviewed method instances from 66 open source Java projects.
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
19
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Ueda
et al. [97]
Revision
pattern
suggestion
Open
Source
1
Python
Chunk level
616,723 code
changes
Classification Metrics: Accuracy; Confidence Ra-
tio
Tufano
et al. [95]
Revised code gener-
ation
Open
Source
8,954
Java
Method level 17,194 meth-
ods
Text-Matching / Generation Metrics: Exact-
Match; BLEU-4; Levenshtein Distance.
Thongtanunam
et al. [91]
Revised code gener-
ation
Open
Source
3
Java
Method level 630,858
methods
Text-Matching / Generation Metrics: Exact
Matches
Lin
and
Thongta-
nunam [53]
Revised code gener-
ation
Open
Source
20,243
Java
Method level 332,429
methods
Text-Matching / Generation Metrics: Exact
Match; MRR ; Dataflow Match (DFM).
Pornprasit
et al. [71]
Revised code gener-
ation
Open
Source
3
Java
Method level 57,615 meth-
ods
Text-Matching / Generation Metrics: Exact
Match.
Yin
et al. [108]
Revised code gener-
ation
Open
Source
N/A
Java
Method level 15,909 meth-
ods
Text-Matching / Generation Metrics: BLEU
(BLEU-1...BLEU-4 Average); ROUGE-L; Normal-
ized Levenshtein Distance.
Cao
et al. [8]
Revised code gener-
ation
Open
Source
8,954
Java
Method level 17,194 meth-
ods
Text-Matching / Generation Metrics: Perfect Pre-
diction Rate; BLEU-4; Levenshtein Distance.
Table 15. Datasets and Evaluation Metrics for Code Revision in the Pre-LLM era.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Kim
et al. [46]
Method naming rec-
ommendation
Open
Source
66
Java
Method level 400 reviewed
method
Classification Metrics: Precision; Recall; F-
measure; Accuracy
Table 16. Dataset for Code Refactoring (Rework) in Pre-LLM era.
5
RQ2: Code Review Datasets and Evaluation in LLM era
In this research question, we investigate the range of tasks and datasets introduced for evaluating
LLM-driven code-review systems. Figure 5 presents an overview of the literature taxonomy in the
LLM era.
5.1
Review Prioritization/Selection
5.1.1
Change Qality Prediction: In the LLM era, this task primarily involves leveraging language
models to predict whether a code change requires human review and to assess its quality. The
table 17 summarizes representative LLM era datasets for the task. Most of the studies [49, 50, 60]
focus on review necessity prediction, where the goal is to determine whether a code change should
receive human review. These approaches typically operate at the chunk or file level, analyzing
code diffs to filter out changes that do not need attention, allowing reviewers to concentrate on
necessary reviews rather than uniformly inspecting all submissions. Two studies [35, 106] address
quality estimation from a reviewer’s perspective, aiming to judge whether a code change is of
acceptable quality or likely to require comments. These tasks often leverage code diffs enriched
with contextual information such as issue descriptions, pull request metadata and surrounding code
to improve decision-making. Most of the studies use of open-source repositories, multi-language
coverage (e.g. Python, Java, C++, JavaScript), and granularity ranging from chunk-level diffs to
method-level changes. Evaluation is generally based on classification metrics. One study [49]
uniquely incorporates usability-oriented measures (e.g., Blocking Rate and Review-Omission Rate)
and relies on private data from student-written code samples to decide whether a submitted solution
meaningfully attempts a given task.
5.2
Change Understanding and Analysis:
5.2.1
Relevant Change Identification: Table 18 presents the dataset used to evaluate this capability
in the LLM era. Unlike earlier Pre-LLM approaches that concentrated on detecting inconsistencies
, Vol. 1, No. 1, Article . Publication date: February 2026.

20
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
Fig. 5. Distribution of literature across different code review-related tasks in the LLM era.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Li et al. [50] Review
necessity
prediction
Open
Source
1,161
C, C++ and 7 others Chunk level
328,000
diff
hunk
Classification Metrics: Accuracy; Precision; Re-
call; F-measure.
Lu et al. [60] Review
necessity
prediction
Open
Source
N/A
Go, Java and 7 oth-
ers
Chunk level 28
 0 sam-
ples
Classification Metrics: Precision; Recall; F-
measure.
Lee
and
Joe [49]
Review
necessity
prediction
Private
projects
27
Python
File level
93
submis-
sion
of
27
coding
test
questions
Classification Metrics: Blocking Rate; Review-
Omission Rate; Usability Ratings.
Hu
et al. [35]
Change quality pre-
diction
Open
Source
90
Java, JavaScript and
7 others
Method level 67,910
changed
methods.
Classification Metrics: F-measure; Accuracy; Pre-
cision.
Yan
et al. [106]
Change quality pre-
diction
Open
Source
N/A
Python, Java and 7
others
Chunk level
1,800
code
changes.
Classification Metrics: Accuracy; Precision; Re-
call; F-measure.
Table 17. Datasets and Evaluation Metrics for Change Quality Prediction in LLM era.
or missed edits, the dataset of the LLM era treats the task as change type recognition (i.e., delete,
modify, and add), where models classify the intended action using both code fragments and review
context. The dataset introduced by Lin et al. [52] spans a diverse set of open-source projects
across multiple languages that operate at the chunk level. Evaluation relies on classification-based
measures (eg. MCQA accuracy and PPA) to assess correctness and agreement in multi-choice
settings.
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
21
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Lin
et al. [52]
Change Type Recog-
nition
Open
Source
199
PHP, Ruby and 7
others
Chunk level
900
code
review
examples.
Classification Metrics: MCQA accuracy; PPA.
Table 18. Datasets and Evaluation Metrics for Relevant Change Identification in LLM era.
5.3
Peer Review
5.3.1
Defective Change Prediction: Defective Change Prediction in the LLM era leverages language
models to identify buggy or vulnerable code by reasoning about semantic correctness and security
weaknesses which enable automated systems to detect potentially risky code artifacts for further
review. Table 19 summarizes representative datasets used to evaluate LLM-driven defective change
prediction and vulnerability identification. Most studies address defect prediction [37, 49]. These
tasks are explored across both private and real-world open source settings. In educational contexts,
LLMs are used to analyze programming exercise submissions to detect semantic issues (e.g., unmet
requirements, hard-coded solutions, or logical mistakes) that standard test cases may overlook
[49]. Beyond defect prediction, one dataset extends this research toward security-oriented analysis
through CWE prediction [27], identifying standardized vulnerability types in C and C++ source
code.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Lee
and
Joe [49]
Defect prediction
Private
Projects
27
Python
File level
93
answers
for
coding
test
Classification Metrics: Relative Error Detection
Rate (REDR); Fisher’s Exact Test
Icoz
and
Biricik [37]
Defect prediction
Open
Source
N/A
Python
Method level 27,000
labeled code
examples
Classification Metrics: Precision; Recall; F-
measure; Accuracy.
Guo
et al. [27]
CWE prediction
Open
Source
26
C, C++
Method level 251,000 lines
of code
Classification Metrics: True Positives (TP); False
Positives (FP); Precision; Recall
Table 19. Datasets and Evaluation Metrics for Defective Change Prediction in LLM era.
5.3.2
Isue Labeling / Clasification: Issue Labeling and Classification assigns clear problem types
to review comments or code changes so that reviewers can quickly understand what kind of issue
is present (e.g., bugs, readability or maintainability issues) and focus on the most important parts
of the code. Table 20 summarizes datasets constructed by studies in the LLM era. Most studies
focus on issue categorization, where the goal is to classify review comments or code changes into
predefined issue categories [23, 79, 88, 93]. For example, Tufano et al. [93] constructed a dataset
from open-source projects by injecting issues such as code duplication, missing documentation, and
logic bugs into small Java and Python programs at the file level, in order to study how automated
reviews affect human reviewer performance. At an industrial scale, Sun et al. [88] built a dataset
for a production-ready system that categorizes issues such as security vulnerabilities, functional
bugs, and performance problems in private projects using a two-stage LLM pipeline. These studies
are mainly evaluated using standard classification metrics. Beyond simple categorization, Sghaier
and Sahraoui [79] extended the task to multi-perspective labeling, which includes predicting issue
types from review comments, classifying code changes, and locating problematic lines in source
code.
, Vol. 1, No. 1, Article . Publication date: February 2026.

22
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Tufano
et al. [93]
Issue categorization
(e.g., code duplica-
tion, documentation
issues,
and
logic
bugs)
Open
Source
6
Java, Python
File level
12 programs
with 48 in-
jected issues
Classification Metrics: Accuracy.
User Study: Time spent; Reviewer confidence
Sun
et al. [88]
Issue categorization
using comment
Private
Projects
N/A
Go, JavaScript and 3
others
Chunk level
120,000
review
comments
Classification Metrics: Precision; Recall.
Sghaier and
Sahraoui [79]
Issue categorization
using comment
Open
Source
19
N/A
Chunk level
262
com-
ments
Classification Metrics: Precision; Recall; F-
measure
Issue categorization
using code change
Open
Source
19
N/A
Chunk level
94,121 code
review pairs
Classification Metrics: Precision; Recall; F-
measure; Accuracy
Goldman
et al. [23]
Issue categorization
(e.g.,
Readability,
Bugs, Maintainabil-
ity, Design, etc)
Open-
source
&
Private
Project
N/A
N/A
Chunk level
94,121
review
comments
Classification Metrics: Accuracy
Table 20. Datasets and Evaluation Metrics for Issue Labeling/Classification in the LLM era.
5.3.3
Isue Localization: Table 21 summarizes datasets that link reviewer feedback or intent to
exact code locations, enabling precise issue localization in real-world code review workflows.
The majority of these studies constructed dataset to address line-level issue or defect localization,
where models are required to predict or rank specific faulty lines within a pull request or code
diff by jointly reasoning over code context and reviewer feedback [35, 79, 109, 111]. While some
approaches concentrate on ranking suspicious lines within modified methods using both textual
context and before-and-after code structure [35], others aim to detect exact faulty regions across
entire diffs by analyzing code logic, recent changes, edge cases, resource handling, and API usage
based on real-world reviewer comments [109]. Change localization focuses on determining where
a requested modification should be applied in response to a natural-language review comment
by reasoning over reviewer intent and pre-change code context [52]. Across these datasets, issue
localization is studied using large-scale open-source code review scenarios, based on pull requests,
diffs, and review comments. The tasks focus on precisely identifying faulty code regions at fine
granularity, such as line, chunk, or method level, rather than whole files. Evaluation in these
datasets mainly uses classification-based accuracy and ranking metrics.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Zeng
et al. [111]
Line-level issue lo-
calization
Open
Source
12
Python
PR level
1,000 pull re-
quests
Classification Metrics: True Positive; False Posi-
tive; False Negative; Precision; Recall; F-measure.
Sghaier and
Sahraoui [79]
Line-level issue lo-
calization
Open
Source
19
N/A
Chunk level
2,349,102
lines of code
Classification Metrics: Precision; Recall; F-
measure
Hu
et al. [35]
Line-level defect lo-
calization
Open
Source
90
Python, Java and 7
others
Method level 13,512
entries.
Classification Metrics: pass@K.
Yu
et al. [109]
Line-level issue lo-
calization
Open
Source
11,324
Java, Python and 3
others
Chunk level
12,881 code
review exam-
ples
Text-Matching / Generation Metrics: IoU (Inter-
section over Union).
Lin
et al. [52]
Change Localization
(where
to
apply
changes)
Open
Source
199
C, C++ and 7 others Chunk level
900
code
review
examples.
Classification Metrics:MCQA accuracy; PPA; Per-
plexity and 5-gram accuracy.
Table 21. Datasets and Evaluation Metrics for Issue Localization in the LLM era.
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
23
5.3.4
Review Coment Generation: Review Comment Generation in the LLM era focuses on
automatically producing helpful, human-like review comments for code changes by leveraging
large language models’ code understanding and natural-language generation capabilities. Table 22
summarizes datasets for this task. Most of the datasets are constructed for natural language (NL)
comment generation [12, 35, 40, 49, 50, 60, 80, 106]. Beyond free-form comments, several studies
introduce structured review generation, requiring models to produce organized feedback across
multiple dimensions such as functionality, complexity, style, documentation, and defects [26, 109],
and some further reframe the task as defect identification and explanation with emphasis on
fault localization and impact analysis in high-risk contexts [40, 58]. Tasks are primarily defined
at chunk and method levels, with specialized line-level settings demanding precise, localized
feedback [35]. The evaluation typically combines text-generation metrics with human or LLM-
as-a-Judge assessments to measure clarity, usefulness, and alignment with reviewer intent [26,
30, 35, 50, 58, 60, 109]. Most datasets are derived from open-source projects, enabling models
to learn realistic reviewer behavior from historical comments and diffs [30, 50, 60, 80], while a
smaller but important subset extends to private or industrial settings, including multilingual review
generation [12], educational code review on student submissions [49], and defect-focused industrial
reviews [58]. [30] relies on a dataset which is reused by several subsequent studies ( [28, 43, 47, 54–
56, 66, 70, 110]). Similarly, [60] belongs to another group of works ( [41, 59, 94, 113]) that employ
closely related or overlapping datasets originating from the same CodeReviewer-style corpus.
5.3.5
Security Detection: Table 23 summarizes datasets that formulate this task as vulnerability
analysis during peer review. The dataset by Tang et al. [89] models security-focused peer review
as automated analysis of code changes, where LLMs inspect commits to detect a broad range of
security issues introduced during development.
5.4
Review Assessment and Analysis
5.4.1
Sentiment / Toxicity Analysis. Table 24 presents an LLM era dataset for analyzing toxicity in
code review comments, focusing on identifying harmful or unprofessional language. The task is
framed at the pull-request level within an educational and competition-oriented framework, where
LLMs automatically evaluate review comments and assign toxicity scores to discourage negative
behavior [45]. Unlike Pre-LLM approaches that relied on static label prediction, this dataset uses
LLMs for context-aware interpretation of reviewer language, prioritizing collaboration quality over
code correctness.
5.4.2
Coment-Code Compliance. Table 25 presents datasets for review comment–code compli-
ance tasks in the LLM era. Most studies focus on review comment–code consistency analysis, where
models check if natural-language artifacts—such as review comments, commit messages, or linked
issues—accurately describe the related code changes (e.g., semantic relevance or missing/tangled
edits) [39, 45, 89]. This task appears in both educational and open-source contexts and requires
reasoning over code diffs and text to detect mismatches or mixed changes. An important extension
is coding format consistency analysis, which verifies whether code modifications follow existing
style conventions, showing that consistency checks now go beyond semantic intent [89]. These
datasets cover pull-request and commit-level analysis, include programming course projects and
large-scale open-source repositories, and span multiple languages, reflecting the diversity of eval-
uation settings. Common evaluation approaches use classification metrics (e.g., precision, recall,
F-measure) or rating-based metrics, combining automated correctness checks with human-like
judgment of semantic alignment.
, Vol. 1, No. 1, Article . Publication date: February 2026.

24
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Li et al. [50] NL
Review
com-
ment generation
Open
Source
1,161
C, C++ and 7 others Chunk level
138,000 hunk
-
comment
pairs
Text-Matching / Generation Metrics: BLEU
Lu
et al. [60];
NL
Review
com-
ment generation
Open
Source
CodeReviewer
+
Tufano
datasets
C, C++ and 7 others Method level Tufano:
168
 0
functions
CRer: 138000
diffs
Text-Matching / Text-Generation Metrics: BLEU,
top-1 results.
Yan
et al. [106]
NL
Review
com-
ment generation
Open
Source
N/A
Python, Java and 7
others
Chunk level
1,800
code
samples.
Text-Matching / Generation Metrics: BLEU;
ROUGE; BERTScore.
Haider
et al. [30];
NL
Review
com-
ment generation
Open
Source
CodeReviewer
dataset
C, C++ and 7 others Chunk level 143
 50 sam-
ples
Text-Matching / Generation Metrics: BLEU;
BERTScore.
Chervyakov
et al. [12]
NL review comment
generation
(Rus-
sian)
Private
projects
N/A
Java, Python, Go,
Scala
Chunk level
689
pairs
of
merge-
request
Text-Matching / Generation Metrics: Judge@k
(LLM-as-a-Judge); BLEU; chrF.
Jaoua
et al. [40]
NL
Review
com-
ment generation
Open
Source
N/A
Java
Chunk level
27,267
Java
examples
Classification Metrics: Accuracy Labels (accurate
/ partially accurate / not accurate)
Lee
and
Joe [49]
NL review comment
generation
Private
projects
27
Python
File level
93 test cases
for 27 ques-
tions written
by 72 student
Text-Matching / Generation Metrics: BERTScore.
Sghaier
et al. [80]
NL review comment
generation
Open
Source
N/A
PHP, Ruby and 7
others
Chunk level
176,613 sam-
ples
Text-Matching / Generation Metrics: BLEU.
Hu
et al. [35]
NL review comment
generation
Open
Source
90
Python, Java and 7
others
Method level 13,512
entries.
Text-Matching / Generation Metrics: ROUGE-1;
ROUGE-L; Edit Similarity.
Lu et al. [58] Defect identification
and explanation
Private
projects
N/A
C++
PR level
45 merge re-
quests
Classification Metrics: False Alarm Rate; LSR
(Line Success Rate); CPI1/CPI2 (Comprehensive
Performance Index 1/2)
Yu
et al. [109]
Structured
review
comment
genera-
tion
Open
Source
11,324
Java, Python and 3
others
Chunk level
2,000 test ex-
amples
Retrieval Metric: Hit Rate
Guo
et al. [26]
Structured
re-
view
generation
(e.g.,
function,
complexity,
style,
documentation, and
defects)
Open
Source
70
Python
Chunk level
601
code
review
examples
Text-Matching / Generation Metrics:BLEU; LLM-
as-a-Judge
Classification Metrics: Precision, Recall, and F-
Measure.
Table 22. Datasets and Evaluation Metrics for Review Comment Generation in LLM era.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Tang
et al. [89]
Vulnerability detec-
tion
Open
Source
180
Python, Java and 7
others
Commit level 3,545
commits
Classification Metrics: Hit Rate
Table 23. Dataset for Security Detection in LLM era.
5.4.3
Review Qality Evaluation: Table 26 summarizes the datasets in the LLM era that primarily
address review comment helpfulness analysis, which can be grouped into two closely related task
settings: educational review quality assessment and real-world development review effectiveness
analysis. In educational contexts, Crandall et al. [15] and Crandall et al. [16] study review helpfulness
as a learning support task in programming courses, where LLM-generated feedback is provided
on student submissions and evaluated at the file level. These studies focus on understanding how
automated reviews support learning, clarity, and defect comprehension, using relatively small-scale
datasets (e.g., AI suggestions, review text pages, and student assignment responses) and relying on
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
25
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Khelifi
et al. [45]
Review
comment
toxicity analysis
Educa-
tional /
Compe-
tition
2
Java, Python
PR level
86
students
in 30 teams
reviewed 40
pull requests
Rating Metrics : LLM toxicity scores (rated on a
negative scale from   1 (least toxic) to   10 (very
toxic))
Table 24. Dataset for Sentiment / Toxicity Analysis in the LLM era.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Khelifi
et al. [45]
Review
comment-
code
consistence
analysis
Programming
course /
compe-
tition
2
Java, Python.
PR-level
40
pull
requests
Rating Metrics: LLM-as-a-Judge relevance scores;
Işık
et al. [39]
Review
comment-
code
consistence
analysis
Open
Source
1
Python
PR level
194
pull-
request
Classification Metrics: Accuracy; Precision; Re-
call; F-measure;
Tang et al. [89]
Review
comment-
code
consistence
analysis
Open
Source
180
Python, Java and 7
others
Commit level 3,545 labeled
commits
Classification Metrics: F-measure; Recall.
Coding format con-
sistence analysis
Open
Source
180
Python, Java and 7
others
Commit level 3,545 labeled
commits
Classification Metrics: F-measure; Recall.
Table 25. Datasets and Evaluation Metrics for Comment–Code Compliance Tasks in the LLM era.
Rating Metrics such as Likert-scale user studies, which require skills in human-centered evaluation
and qualitative analysis. In contrast, studies targeting professional or mixed development settings
emphasize whether review comments lead to concrete code changes. Cihan et al. [14] investigates
review helpfulness in private industrial projects across multiple programming languages at the
PR level, using large-scale pull request data and modeling reviewer impact with Classification
Metrics based on Poisson Regression, which demands statistical modeling and empirical software
engineering expertise. Similarly, Goldman et al. [23] analyzes review helpfulness across open-source
and private projects at a finer chunk-level granularity, leveraging millions of labeled code lines
to measure Resolution Rate with review comments, thereby emphasizing large-scale data mining,
change tracking, and causal reasoning about review effectiveness.
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Crandall
et al. [15]
Review
comment
helpfulness analysis
Programming
course /
compe-
tition
N/A
Java
File level
173 AI sug-
gestions and
58 pages re-
view text.
Rating Metrics: Likert-scale ratings
Cihan
et al. [14]
Review
comment
helpfulness analysis
Private
projects
3
Java, JavaScript and
5 others
PR level
4,335 pull re-
quests
Classification Metrics: Poisson Regression;
Goldman
et al. [23]
Review
comment
helpfulness analysis
Open-
source
&
private
project
N/A
N/A
Chunk level
2,349,102
labeled code
lines
Classification Metrics: Resolution Rate with re-
view comments.
Crandall
et al. [16]
Review helpfulness
analysis
Programming
course /
compe-
tition
N/A
C++
File level
40 responses
of student as-
signments
Rating Metrics: Likert-scale (user study)
Table 26. Datasets and Evaluation Metrics for Review Quality Evaluation in the LLM era.
, Vol. 1, No. 1, Article . Publication date: February 2026.

26
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
5.4.4
Review Sumarization: Table 27 summarizes LLM era datasets that explore pull-request
link summarization as a supporting technique for code review. Rather than revising code directly,
this task aims to assist reviewers by generating concise summaries of external resources (e.g.,
documentation, issue pages, or web links) referenced within pull requests, thereby reducing context
switching and improving review efficiency. The existing studies address both educational and open-
source datasets where models are expected to interpret the surrounding pull-request context and
produce informative link summaries that help reviewers quickly grasp the relevance of referenced
materials. While educational settings emphasize reviewer engagement and perceived usefulness
through human-centered evaluations [45], open-source datasets focus more on the quality and
semantic faithfulness of generated summaries using automated text similarity measures [92].
Despite differences in evaluation emphasis, both settings treat link summarization as a PR-level
auxiliary task designed to support, rather than replace human review activities [45, 92].
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Khelifi
et al. [45]
Pull-request
link
summarization
Programming
course/competition
2
Java , Python
PR-level
86 students,
working
in
30
teams,
reviewed 40
pull requests
User studies
Trakoolgerntong
et al. [92]
Pull-request
link
summarization
Open
Source
50
N/A
PR level
365
sample
links
Text-Matching / Generation Metrics: BLEU; ME-
TEOR; ROUGE-1; ROUGE-2; Sentence Similarity;
BERTScore (Precision, Recall, F-measure
Table 27. Datasets and Evaluation Metrics for Review Summarization in the LLM era.
5.5
Code Refinement
5.5.1
Code Revision: Table 28 summarizes datasets for code revision. The datasets of existing
studies can be categorized into two primary tasks, e.g., revised code generation ( where models
update code to satisfy reviewer feedback) and best candidate selection (where models identify
the correct revision among alternatives). The majority of LLM era datasets focus on revised code
generation, where models must understand natural-language review comments and apply the
requested changes to produce a corrected, review-compliant version of the code [50, 52, 60, 62,
80, 89]. Those datasets emphasize end-to-end code refinement, requiring joint reasoning over
reviewer intent and code context. Most studies rely on open-source repositories, while a smaller
number extend to private industrial settings to evaluate whether generated patches for real-world
deployment [62]. In contrast, best candidate selection is explored in a smaller set of datasets, where
the goal is to identify the correct revision from multiple candidate changes instead of generating
code directly, emphasizing discriminative reasoning over change intent and location [52]. Overall,
LLM era datasets reflect a clear shift toward reviewer-aware code revision, with broader language
coverage, more diverse data sources, and an increased focus on producing or selecting revisions
that align with human reviewer expectations in realistic review scenarios.
6
RQ3: Comparison between Pre-LLM and LLM eras
6.1
Distribution of tasks
Figure 6 provides a complementary view of how research efforts are distributed across code
review–related sub-tasks in the Pre-LLM and LLM eras. In the Pre-LLM era, the datasets are evenly
distributed across all sub-tasks. In contrast, the LLM era exhibits a more uneven distribution of
datasets across sub-tasks. The Change Understanding and Analysis category represents the
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
27
Literature Task
Source Project No. Language
Granularity Data Points Evaluation Metrics
Li et al. [50] Revised code gener-
ation
Open
Source
1,161
C, C++ and 7 others Chunk level
176,000 code-
revisions
Text-Matching / Generation Metrics: BLEU; Exact
Match.
Lu et al. [60] Revised code gener-
ation
Open
Source
N/A
Go, Java and 7 oth-
ers
Method level Tufano:
168
 0
functions
CRer: 138000
diffs
Text-Matching / Text-Generation Metrics: BLEU
Tang
et al. [89]
Revised code gener-
ation
Open
Source
N/A
N/A
Commit level 3,545
Classification Metrics: Edit Progress.
Maddila
et al. [62]
Revised code gener-
ation
Private
projects
N/A
Hack; PHP
Chunk level
2,900
SFT
pairs
Text-Matching / Generation Metrics: Exact
Match; Successful Patch Generation (SPG);
Lin
et al. [52]
Revised code gener-
ation
Open
Source
199
C, C++ and 7 others Chunk level
900
code
review
examples
Text-Matching / Generation Metrics: Exact
Match; Perplexity; 5-gram accuracy.
Best candidate selec-
tion
Open
Source
199
C, C++ and 7 others Chunk level
900
code
review
examples
Classification Metrics: Invariant MCQA accuracy;
PPA; Perplexity; 5-gram accuracy.
Sghaier
et al. [80]
Revised code gener-
ation
Open
Source
N/A
PHP, Ruby and 7
others
Chunk level
20,000 pairs
review com-
ments
Text-Matching / Generation Metrics: CodeBLEU;
Exact Match.
Table 28. Datasets and Evaluation Metrics for Code Revision in LLM era.
Fig. 6. Comparative distribution of datasets across code review sub-tasks in the Pre-LLM vs. LLM eras. Values
atop bars indicate the absolute number of datasets.
most drastic shift between the two eras. In the Pre-LLM era, this was a cornerstone of research
(14 datasets), while in the LLM era, it has nearly vanished as a standalone topic (1 dataset). For
instance, the datasets for Change Decomposition and Impact analysis have completely vanished in
the LLM era. This suggests a fundamental change in research philosophy: moving from helping
humans understand code to having machines perform the task directly end-to-end.
, Vol. 1, No. 1, Article . Publication date: February 2026.

28
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
Fig. 7. Distribution of programming language coverage in datasets across Pre-LLM and LLM eras. Values
indicate the percentage of datasets within each era.
The most obvious trend is the sheer dominance of Peer Review, which accounts for
nearly 60% of datasets in LLM era. In the Pre-LLM era, Peer Review was balanced with other tasks
like Change Understanding and Refinement. In the LLM era, it accounts for nearly 60% of all datasets.
If we look at the sub-tasks, the focus has a radical transformation. LLM era benchmarks reflect a
shift from retrieval-based recommendation toward end-to-end review generation, with growing
emphasis on reasoning quality, contextual understanding, and adaptability across review scenarios,
positioning LLMs as capable virtual reviewers that combine natural-language fluency with fine-
grained code comprehension and review judgment. Also, behind the text-matching/generation
metrics (e.g., BLEU scores and ROUGE), more and more studies use LLM-as-Judge as the evaluation
method. The reason Peer Review grew while Change Understanding shrank is consolidation. In the
LLM era, researchers probably consider that understanding is no longer a separate research goal,
instead it is a prerequisite that is now bundled into Comment Generation.
6.2
Program language
The transition to the LLM era marks a significant shift from language-specific research
to cross-language generalization. Figure 7 illustrates this evolution in the number of program-
ming languages covered per dataset. In the Pre-LLM era, research was highly concentrated on
single-language studies, with nearly three-fifths (59%) of all datasets restricted to a single program-
ming language. In this period, datasets employing more than one language were the minority,
accounting for only 41% of the total. In contrast, the LLM era exhibits a much more dispersed and
multilingual distribution. While single-language datasets remain present, their relative proportion
has plummeted to 24%, a decrease of more than half compared to the Pre-LLM era. There is a
pronounced surge in multi-language studies, with datasets covering five or more languages now
representing 43% of the landscape. Most notably, datasets covering nine or more languages have
become a dominant category in the LLM era (34%), whereas they were virtually non-existent (2%)
in the earlier period. This trend suggests that modern code review benchmarks are increasingly
designed to evaluate the zero-shot transfer capabilities and broad linguistic versatility inherent in
LLMs.
Figure 8 presents the specific programming languages used in datasets of both eras. In the
Pre-LLM era, Java dominates the datasets, accounting for approximately 61% of all datasets, making
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
29
Fig. 8. A specific programming language is covered by the datasets in the Pre-LLM and LLM eras. The values
above each bar indicate the absolute number of datasets for each language.
it by far the most frequently used language. C represents the second most common language at
roughly 13%, followed by Python (about 12%). Most other languages, including C++, PHP, Ruby,
Go, and HTML, individually account for only a small fraction of studies. In contrast, the LLM era
exhibits a substantially more diverse language landscape. Although Java remains prominent, its
relative share decreases to approximately 34%, indicating reduced reliance on a single dominant
language. Python becomes the most widely used language in this era, representing about 41% of
datasets. Additionally, other languages such as C++ (around 28%), C (about 26%), JavaScript (25%),
Go (15%), Ruby (20%), and PHP (21%) gain noticeably greater representation compared to the Pre-
LLM era. The LLM era also introduces several languages that were absent or negligible previously,
including Rust, TypeScript, Scala, SQL, and Hack, reflecting an expansion in experimental scope.
The proportion of datasets that do not specify a programming language decreases markedly from
roughly 15% in the Pre-LLM era to about 4% in the LLM era. This reduction suggests improved
reporting practices in more recent studies.
7
Discussion
7.1
Limitation of existing benchmark in LLM era and future direction
7.1.1
Improve Task Coverage. While current benchmarks in the LLM era primarily evaluate com-
ment generation and code revision, they largely overlook macro-level review responsibilities such
as impact analysis (predicting how changes affect downstream dependencies) and commit decom-
position (identifying “tangential” changes that should be split into separate PRs). These tasks are
foundational to code understanding; mastering them is essential for identifying defects during
peer reviews and generating accurate fixes. Future research should bridge this gap by expanding
task coverage to include these overlooked areas. Furthermore, while some datasets for these tasks
exist, they are often limited to some specific languages like Java, C, and C#. We encourage the
development of new benchmarks to expand to other languages such as Python and Go to align
with the current LLM landscape.
7.1.2
Transitioning from Static to Dynamic Evaluation. Current benchmarks in all code review
tasks (e.g., comment generation, code revision) rely on static metrics (e.g., text match metrics,
, Vol. 1, No. 1, Article . Publication date: February 2026.

30
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
classification metrics) to verify the alignment between the generated comment or code and the
ground truth. These are "shallow" matches that fail to capture the functional correctness of a review.
An LLM might suggest a fix that is grammatically perfect but introduces a deadlock or fails to
compile, yet still receives a high score from static text-matching metrics. More comprehensive
Software Engineering metrics that involve runtime information, such as compilation/build success
rate, regression testing, are strongly encouraged to be integrated in the evaluation pipeline. For
instance, future research could introduce sandboxed verification that integrates automated execution
environments using Docker containers. When an LLM proposes a code revision, the benchmark
should automatically attempt to build the project and examine if the revised code is executable and
breaks existing functionality.
7.1.3
Granular Benchmarking via Task Taxonomy. Existing code revision benchmarks only broadly
measure the LLMs’ overall performance on code revision, rarely breaking it down into different
tasks during the code review process (e.g., refactoring, enhancing documentation and readability,
ensuring correctness and functionality, and optimizing performance). We encourage future research
to build code revision benchmark with task taxonomy that can concretely measure an LLM’s
capabilities from different dimensions. For instance, in our survey, we observe certain works focus
on categorizing the issues during peer review; future research could reuse their taxonomy and
curate a code revision dataset.
7.2
Threats to Validity
Internal Validity: One threat to our data collection is the potential omission of relevant studies.
While we employed a snowballing approach starting from major SE (Software Engineering) and AI
venues [68], there is a risk that certain papers might have been missed. To minimize this, we utilized
multiple iterations of both forward and backward snowballing. We also cross-referenced our initial
seed set with well-known repositories and previous surveys to ensure the most influential datasets
were captured. To extend our search space, we also collected papers from ArXiv. However, we
acknowledge that the rapid growth of AI-driven code review means new datasets are published
weekly, making absolute exhaustiveness a moving target.
Construct Validity: The categorization of tasks is subject to human interpretation. Human
bias or fatigue could lead to inconsistent labeling across the survey. We addressed this by imple-
menting a dual-reviewer protocol. Each paper was labeled independently by at least two authors.
In cases of disagreement, a third senior researcher acted as an arbitrator to reach a consensus.
8
Conclusion
Our survey reveals a transformative shift in the code review landscape, characterized by a transition
from human-centric assistance to autonomous, end-to-end generation. While the Pre-LLM era
featured a balanced distribution of efforts, the LLM era has nearly abandoned these as standalone
tasks, with only one dataset represented as such. Research now favors comprehensive Peer Review
tasks, which account for nearly 60% of LLM-era datasets.
Furthermore, benchmarks have evolved from language-specific studies toward broad cross-
language generalization. The prevalence of single-language datasets decreased from 59% in the
Pre-LLM era to 24% in the LLM era, while highly multilingual benchmarks covering nine or
more languages increased to 34%. Despite these shifts, future research should address current
benchmark limitations by incorporating macro-level responsibilities, such as impact analysis,
and transitioning toward dynamic evaluation methods. By providing a multi-level taxonomy and
systematic classification of 18 sub-tasks, our work establishes a foundation for more rigorous and
context-aware code review automation.
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
31
References
[1] Toufique Ahmed, Amiangshu Bosu, Anindya Iqbal, and Shahram Rahimi. 2017. SentiCR: A customized sentiment
analysis tool for code review interactions. In 2017 32nd IEEE/ACM International Conference on Automated Software
Engineering (ASE). 106–111. doi:10.1109/ASE.2017.8115623
[2] Krishna Teja Ayinala, Kwok Sun Cheng, Kwangsung Oh, Teukseob Song, and Myoungkyu Song. 2020. Code Inspection
Support for Recurring Changes with Deep Learning in Evolving Software. In 2020 IEEE 44th Annual Computers,
Software, and Applications Conference (COMPSAC). 931–942. doi:10.1109/COMPSAC48688.2020.0-149
[3] Alberto Bacchelli and Christian Bird. 2013. Expectations, outcomes, and challenges of modern code review. In 2013
35th International Conference on Software Engineering (ICSE). IEEE, 712–721.
[4] Faruk Balcı, Dilruba Sultan Haliloğlu, Onur Şahin, Cankat Tilki, Mehmet Ata Yurtsever, and Eray Tüzün. 2021.
Augmenting Code Review Experience Through Visualization. In 2021 Working Conference on Software Visualization
(VISSOFT). 110–114. doi:10.1109/VISSOFT52517.2021.00021
[5] Mike Barnett, Christian Bird, João Brunet, and Shuvendu K. Lahiri. 2015. Helping Developers Help Themselves:
Automatic Decomposition of Code Review Changesets. In 2015 IEEE/ACM 37th IEEE International Conference on
Software Engineering, Vol. 1. 134–144. doi:10.1109/ICSE.2015.35
[6] Amiangshu Bosu, Jeffrey C Carver, Christian Bird, Jonathan Orbeck, and Christopher Chockley. 2016. Process aspects
and social dynamics of contemporary code review: Insights from open source development and industrial practice at
microsoft. IEEE Transactions on Software Engineering 43, 1 (2016), 56–75.
[7] Rodrigo Brito and Marco Túlio Valente. 2021. RAID: Tool Support for Refactoring-Aware Code Reviews. 2021 IEEE/ACM
29th International Conference on Program Comprehension (ICPC) (2021), 265–275. https://api.semanticscholar.org/
CorpusID:232307313
[8] Zhenzhen Cao, Sijia Lv, Xinlong Zhang, Hui Li, Qian Ma, Tingting Li, Cheng Guo, and Shikai Guo. 2024. Structuring
Meaningful Code Review Automation in Developer Community. Engineering Applications of Articial Intelligence 127
(2024), 106970. doi:10.1016/j.engappai.2023.106970
[9] Robert Chatley and Lawrence Jones. 2018. Diggit: Automated code review via software repository mining. In 2018
IEEE 25th International Conference on Software Analysis, Evolution and Reengineering (SANER). 567–571. doi:10.1109/
SANER.2018.8330261
[10] Lawrence Chen, Rui Abreu, Tobi Akomolede, Peter C. Rigby, Satish Chandra, and Nachiappan Nagappan. 2022.
Leveraging test plan quality to improve code review efficacy. In Proceedings of the 30th ACM Joint European Software
Engineering Conference and Symposium on the Foundations of Software Engineering (Singapore, Singapore) (ESEC/FSE
2022). Association for Computing Machinery, New York, NY, USA, 1320–1330. doi:10.1145/3540250.3558952
[11] Zhiyuan Chen, Maneesh Mohanavilasam, Young-Woo Kwon, and Myoungkyu Song. 2017. Tool Support for Managing
Clone Refactorings to Facilitate Code Review in Evolving Software. In 2017 IEEE 41st Annual Computer Software and
Applications Conference (COMPSAC) (2017-07), Vol. 1. 288–297. doi:10.1109/COMPSAC.2017.242
[12] Artem Chervyakov, Alexander Kharitonov, Pavel Zadorozhny, Adamenko Pavel, Rodion Levichev, Dmitrii Vorobev,
Dmitrii Salikhov, Aidar Valeev, Alena Pestova, Maria Dziuba, Ilseyar Alimova, Artem Zavgorodnev, Aleksandr
Medvedev, Stanislav Moiseev, Elena Bruches, Daniil Grebenkin, Roman Derunets, Vikulov Vladimir, Anton Emelyanov,
Dmitrii Babaev, Vladimir V. Ivanov, Valentin Malykh, and Alena Fenogenova. 2025. MERA Code: A Unified Framework
for Evaluating Code Generation Across Tasks. arXiv:2507.12284 [cs.SE] https://arxiv.org/abs/2507.12284
[13] Moataz Chouchen and Ali Ouni. 2023. A multi-objective effort-aware approach for early code review prediction and
prioritization. Empirical Software Engineering 29, 1 (2023), 29. doi:10.1007/s10664-023-10431-7
[14] Umut Cihan, Vahid Haratian, Arda İçöz, Mert Kaan Gül, Ömercan Devran, Emircan Furkan Bayendur, Baykal Mehmet
Uçar, and Eray Tüzün. 2025. Automated Code Review in Practice. In 2025 IEEE/ACM 47th International Conference on
Software Engineering: Software Engineering in Practice (ICSE-SEIP). 425–436. doi:10.1109/ICSE-SEIP66354.2025.00043
[15] Aaron S. Crandall, Bryan J. Fischer, and Johannah L. Crandall. 2024. WIP: ARTful Insights from a Pilot Study
on GPT-Based Automatic Code Reviews in Undergraduate Computer Science Programs. In 2024 IEEE Frontiers in
Education Conference (FIE). 1–5. doi:10.1109/FIE61694.2024.10893407
[16] Aaron S. Crandall, Gina Sprint, and Bryan Fischer. 2023. Generative Pre-Trained Transformer (GPT) Models as a
Code Review Feedback Tool in Computer Science Programs. J. Comput. Sci. Coll. 39, 1 (Oct. 2023), 38–47.
[17] Felipe Ebert, Fernando Castor, Nicole Novielli, and Alexander Serebrenik. 2017. Confusion Detection in Code Reviews.
In 2017 IEEE International Conference on Software Maintenance and Evolution (ICSME). 549–553. doi:10.1109/ICSME.
2017.40
[18] Carolyn D. Egelman, Emerson Murphy-Hill, Elizabeth Kammer, Margaret Morrow Hodges, Collin Green, Ciera Jaspan,
and James Lin. 2020. Predicting developers’ negative feelings about code review. In Proceedings of the ACM/IEEE
42nd International Conference on Software Engineering (Seoul, South Korea) (ICSE ’20). Association for Computing
Machinery, New York, NY, USA, 174–185. doi:10.1145/3377811.3380414
, Vol. 1, No. 1, Article . Publication date: February 2026.

32
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
[19] Yuanrui Fan, Xin Xia, David Lo, and Shanping Li. 2018. Early prediction of merged code changes to prioritize
reviewing tasks. Empirical Softw. Engg. 23, 6 (Dec. 2018), 3346–3393. doi:10.1007/s10664-018-9602-0
[20] Alex Fish, Thuy Linh Nguyen, and Myoungkyu Song. 2018. CloneMap: A Clone-Aware Code Inspection Tool
in Evolving Software. In 2018 IEEE International Conference on Electro/Information Technology (EIT). 0368–0372.
doi:10.1109/EIT.2018.8500143
[21] Enrico Fregnan, Josua Fr¨"ohlich, Davide Spadini, and Alberto Bacchelli. 2023. Graph-based visualization of merge
requests for code review. J. Syst. Softw. 195, C (Jan. 2023), 20 pages. doi:10.1016/j.jss.2022.111506
[22] Xi Ge, Saurabh Sarkar, Jim Witschey, and Emerson Murphy-Hill. 2017. Refactoring-aware code review. In 2017 IEEE
Symposium on Visual Languages and Human-Centric Computing (VL/HCC). 71–79. doi:10.1109/VLHCC.2017.8103453
[23] Saul Goldman, Hong Yi Lin, Jirat Pasuksmit, Patanamon Thongtanunam, Kla Tantithamthavorn, Zhe Wang, Ray Zhang,
Ali Behnaz, Fan Jiang, Michael Siers, Ryan Jiang, Mike Buller, Minwoo Jeong, and Ming Wu. 2025. What Types of Code
Review Comments Do Developers Most Frequently Resolve? arXiv:2510.05450 [cs.SE] https://arxiv.org/abs/2510.05450
[24] Bo Guo, Young-Woo Kwon, and Myoungkyu Song. 2019. Decomposing Composite Changes for Code Review and
Regression Test Selection in Evolving Software. Journal of Computer Science and Technology 34, 2 (2019), 416–436.
doi:10.1007/s11390-019-1917-9
[25] Bo Guo, Young-Woo Kwon, and Myoungkyu Song. 2019. Decomposing Composite Changes for Code Review and
Regression Test Selection in Evolving Software. Journal of Computer Science and Technology 34, 2 (2019), 416–436.
doi:10.1007/s11390-019-1917-9
[26] Hanyang Guo, Xunjin Zheng, Zihan Liao, Hang Yu, Peng DI, Ziyin Zhang, and Hong-Ning Dai. 2025. CodeFuse-
CR-Bench: A Comprehensiveness-aware Benchmark for End-to-End Code Review Evaluation in Python Projects.
arXiv:2509.14856 [cs.SE] https://arxiv.org/abs/2509.14856
[27] Jinyao Guo, Chengpeng Wang, Xiangzhe Xu, Zian Su, and Xiangyu Zhang. 2025. RepoAudit: An Autonomous
LLM-Agent for Repository-Level Code Auditing. arXiv:2501.18160 [cs.SE] https://arxiv.org/abs/2501.18160
[28] Qi Guo, Junming Cao, Xiaofei Xie, Shangqing Liu, Xiaohong Li, Bihuan Chen, and Xin Peng. 2024. Exploring the
Potential of ChatGPT in Automated Code Refinement: An Empirical Study. In Proceedings of the IEEE/ACM 46th
International Conference on Software Engineering (Lisbon, Portugal) (ICSE ’24). Association for Computing Machinery,
New York, NY, USA, Article 34, 13 pages. doi:10.1145/3597503.3623306
[29] Anshul Gupta. 2018. Intelligent code reviews using deep learning. https://api.semanticscholar.org/CorpusID:52219239
[30] Md. Asif Haider, Ayesha Binte Mostofa, Sk. Sabit Bin Mosaddek, Anindya Iqbal, and Toufique Ahmed. 2024. Prompting
and Fine-tuning Large Language Models for Automated Code Review Comment Generation. arXiv:2411.10129 [cs.SE]
https://arxiv.org/abs/2411.10129
[31] Quinn Hanam, Ali Mesbah, and Reid Holmes. 2019. Aiding Code Change Understanding with Semantic Change
Impact Analysis. In 2019 IEEE International Conference on Software Maintenance and Evolution (ICSME). 202–212.
doi:10.1109/ICSME.2019.00031
[32] Haytham Hijazi, José Cruz, João Castelhano, Ricardo Couceiro, Miguel Castelo-Branco, Paulo de Carvalho, and
Henrique Madeira. 2021. iReview: an Intelligent Code Review Evaluation Tool using Biofeedback. In 2021 IEEE 32nd
International Symposium on Software Reliability Engineering (ISSRE). 476–485. doi:10.1109/ISSRE52982.2021.00056
[33] Yang Hong, Chakkrit Tantithamthavorn, Patanamon Thongtanunam, and Aldeida Aleti. 2022. CommentFinder: a
simpler, faster, more accurate code review comments recommendation. In Proceedings of the 30th ACM Joint European
Software Engineering Conference and Symposium on the Foundations of Software Engineering (Singapore, Singapore)
(ESEC/FSE 2022). Association for Computing Machinery, New York, NY, USA, 507–519. doi:10.1145/3540250.3549119
[34] Yang Hong, Chakkrit Kla Tantithamthavorn, and Patanamon Pick Thongtanunam. 2022. Where Should I Look at?
Recommending Lines that Reviewers Should Pay Attention To. In 2022 IEEE International Conference on Software
Analysis, Evolution and Reengineering (SANER). 1034–1045. doi:10.1109/SANER53432.2022.00121
[35] Ruida Hu, Xinchen Wang, Xin-Cheng Wen, Zhao Zhang, Bo Jiang, Pengfei Gao, Chao Peng, and Cuiyun Gao. 2025.
Benchmarking LLMs for Fine-Grained Code Review with Enriched Context in Practice. arXiv:2511.07017 [cs.SE]
https://arxiv.org/abs/2511.07017
[36] Yuan Huang, Xingjian Liang, Zhihao Chen, Nan Jia, Xiapu Luo, Xiangping Chen, Zibin Zheng, and Xiaocong Zhou.
2021. Reviewing rounds prediction for code patches. Empirical Software Engineering 27, 1 (2021), 7. doi:10.1007/s10664-
021-10035-z
[37] Busra Icoz and Goksel Biricik. 2025. Automated Code Review Using Large Language Models with Symbolic Reasoning.
arXiv:2507.18476 [cs.SE] https://arxiv.org/abs/2507.18476
[38] Khairul Islam, Toufique Ahmed, Rifat Shahriyar, Anindya Iqbal, and Gias Uddin. 2022. Early prediction for merged
vs abandoned code changes in modern code reviews. Information and Software Technology 142 (2022), 106756.
doi:10.1016/j.infsof.2021.106756
[39] Ali Tunahan Işık, Hatice Kübra Çağlar, and Eray Tüzün. 2025. Enhancing Pull Request Reviews: Leveraging Large
Language Models to Detect Inconsistencies Between Issues and Pull Requests. In 2025 IEEE/ACM Second International
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
33
Conference on AI Foundation Models and Software Engineering (Forge). 168–178. doi:10.1109/Forge66646.2025.00027
[40] Imen Jaoua, Oussama Ben Sghaier, and Houari Sahraoui. 2025. Combining Large Language Models with Static Ana-
lyzers for Code Review Generation . In 2025 IEEE/ACM 22nd International Conference on Mining Software Repositories
(MSR). IEEE Computer Society, Los Alamitos, CA, USA, 174–186. doi:10.1109/MSR66628.2025.00038
[41] Yanjie Jiang, Hui Liu, Tianyi Chen, Fu Fan, Chunhao Dong, Kui Liu, and Lu Zhang. 2025. Deep Assessment of Code
Review Generation Approaches: Beyond Lexical Similarity. arXiv:2501.05176 [cs.SE] https://arxiv.org/abs/2501.05176
[42] Tetsuya Kanda, Kazumasa Shimari, and Katsuro Inoue. 2022. didiffff: A Viewer for Comparing Changes in both Code
and Execution Traces. In 2022 IEEE/ACM 30th International Conference on Program Comprehension (ICPC). 528–532.
doi:10.1145/3524610.3527877
[43] Manav Nitin Kapadnis, Atharva Naik, and Carolyn Rose. 2025. CRScore++: Reinforcement Learning with Verifiable
Tool and AI Feedback for Code Review. arXiv:2506.00296 [cs.SE] https://arxiv.org/abs/2506.00296
[44] Ritu Kapur, Balwinder Sodhi, Poojith U Rao, and Shipra Sharma. 2021. Using Paragraph Vectors to improve our
existing code review assisting tool-CRUSO. In Proceedings of the 14th Innovations in Software Engineering Conference
(Formerly Known as India Software Engineering Conference) (Bhubaneswar, Odisha, India) (ISEC ’21). Association for
Computing Machinery, New York, NY, USA, Article 10, 11 pages. doi:10.1145/3452383.3452393
[45] Jasem Khelifi, Moataz Chouchen, Ali Ouni, Dong Wang, Raula Gaikovina Kula, Salma Hamza, and Mohamed Wiem
Mkaouer. 2024. GitRev: An LLM-Based Gamification Framework for Modern Code Review Activities. In 2024 IEEE
International Conference on Source Code Analysis and Manipulation (SCAM). 235–241. doi:10.1109/SCAM63643.2024.
00031
[46] Kisub Kim, Xin Zhou, Dongsun Kim, Julia Lawall, Kui Liu, Tegawende F. Bissyande, Jacques Klein, Jaekwon Lee,
and David Lo. 2025. How Are We Detecting Inconsistent Method Names? An Empirical Study from Code Review
Perspective. ACM Trans. Softw. Eng. Methodol. 34, 6, Article 178 (July 2025), 27 pages. doi:10.1145/3711901
[47] Jahnavi Kumar and Sridhar Chimalakonda. 2024. Code Review Automation Via Multi-task Federated LLM – An
Empirical Study. arXiv:2412.15676 [cs.SE] https://arxiv.org/abs/2412.15676
[48] Harsh Lal and Gaurav Pahwa. 2017. Code review analysis of software system using machine learning techniques. In
2017 11th International Conference on Intelligent Systems and Control (ISCO). 8–13. doi:10.1109/ISCO.2017.7855962
[49] Dong-Kyu Lee and Inwhee Joe. 2025. A GPT-Based Code Review System With Accurate Feedback for Programming
Education. IEEE Access 13 (2025), 105724–105737. doi:10.1109/ACCESS.2025.3581139
[50] Zhiyu Li, Shuai Lu, Daya Guo, Nan Duan, Shailesh Jannu, Grant Jenks, Deep Majumder, Jared Green, Alexey
Svyatkovskiy, Shengyu Fu, and Neel Sundaresan. 2022. Automating code review activities by large-scale pre-training.
In Proceedings of the 30th ACM Joint European Software Engineering Conference and Symposium on the Foundations of
Software Engineering (Singapore, Singapore) (ESEC/FSE 2022). Association for Computing Machinery, New York, NY,
USA, 1035–1047. doi:10.1145/3540250.3549081
[51] Hong-Yi Lin, Chunhua Liu, Haoyu Gao, Patanamon Thongtanunam, and Christoph Treude. 2025. Codereviewqa: The
code review comprehension assessment for large language models. In Findings of the Association for Computational
Linguistics: ACL 2025. 9138–9166.
[52] Hong Yi Lin, Chunhua Liu, Haoyu Gao, Patanamon Thongtanunam, and Christoph Treude. 2025. CodeReviewQA: The
Code Review Comprehension Assessment for Large Language Models. In Findings of the Association for Computational
Linguistics: ACL 2025, Wanxiang Che, Joyce Nabende, Ekaterina Shutova, and Mohammad Taher Pilehvar (Eds.).
Association for Computational Linguistics, Vienna, Austria, 9138–9166. doi:10.18653/v1/2025.findings-acl.476
[53] Hong Yi Lin and Patanamon Thongtanunam. 2023. Towards Automated Code Reviews: Does Learning Code Structure
Help?. In 2023 IEEE International Conference on Software Analysis, Evolution and Reengineering (SANER). 703–707.
doi:10.1109/SANER56733.2023.00075
[54] Hong Yi Lin, Patanamon Thongtanunam, Christoph Treude, Michael W. Godfrey, Chunhua Liu, and Wachiraphan
Charoenwet. 2025. Leveraging Reviewer Experience in Code Review Comment Generation. ACM Trans. Softw. Eng.
Methodol. (Aug. 2025). doi:10.1145/3762183 Just Accepted.
[55] Chunhua Liu, Hong Yi Lin, and Patanamon Thongtanunam. 2025. Too Noisy To Learn: Enhancing Data Quality for
Code Review Comment Generation . In 2025 IEEE/ACM 22nd International Conference on Mining Software Repositories
(MSR). IEEE Computer Society, Los Alamitos, CA, USA, 236–248. doi:10.1109/MSR66628.2025.00043
[56] Fang Liu, Simiao Liu, Yinghao Zhu, Xiaoli Lian, and Li Zhang. 2025. SecureReviewer: Enhancing Large Language
Models for Secure Code Review through Secure-aware Fine-tuning. arXiv:2510.26457 [cs.SE] https://arxiv.org/abs/
2510.26457
[57] Shaoying Liu, Honghui Li, Zhouxian Jiang, Xiuru Li, Feng Liu, and Yan Zhong. 2021. Rigorous code review by reverse
engineering. Information and Software Technology 133 (2021), 106503. doi:10.1016/j.infsof.2020.106503
[58] Junyi Lu, Lili Jiang, Xiaojia Li, Jianbing Fang, Fengjun Zhang, Li Yang, and Chun Zuo. 2025. Towards Practical
Defect-Focused Automated Code Review. arXiv:2505.17928 [cs.SE] https://arxiv.org/abs/2505.17928
, Vol. 1, No. 1, Article . Publication date: February 2026.

34
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
[59] Junyi Lu, Xiaojia Li, Zihan Hua, Lei Yu, Shiqi Cheng, Li Yang, Fengjun Zhang, and Chun Zuo. 2025. DeepCRCEval:
Revisiting the Evaluation of Code Review Comment Generation. In Fundamental Approaches to Software Engineering:
28th International Conference, FASE 2025, Held as Part of the International Joint Conferences on Theory and Practice of
Software, ETAPS 2025, Hamilton, ON, Canada, May 38, 2025, Proceedings (Hamilton, ON, Canada). Springer-Verlag,
Berlin, Heidelberg, 43–64. doi:10.1007/978-3-031-90900-9_3
[60] Junyi Lu, Lei Yu, Xiaojia Li, Li Yang, and Chun Zuo. 2023. LLaMA-Reviewer: Advancing Code Review Automation
with Large Language Models through Parameter-Efficient Fine-Tuning . In 2023 IEEE 34th International Symposium
on Software Reliability Engineering (ISSRE). IEEE Computer Society, Los Alamitos, CA, USA, 647–658. doi:10.1109/
ISSRE59848.2023.00026
[61] Victor da C. Luna Freire, João Brunet, and Jorge C. A. de Figueiredo. 2018. Automatic Decomposition of Java Open
Source Pull Requests: A Replication Study. In SOFSEM 2018: Theory and Practice of Computer Science, A Min Tjoa,
Ladjel Bellatreche, Stefan Biffl, Jan van Leeuwen, and Jiří Wiedermann (Eds.). Springer International Publishing,
Cham, 255–268.
[62] Chandra Maddila, Negar Ghorbani, James Saindon, Parth Thakkar, Vijayaraghavan Murali, Rui Abreu, Jingyue Shen,
Brian Zhou, Nachiappan Nagappan, and Peter C. Rigby. 2025. AI-Assisted Fixes to Code Review Comments at Scale.
arXiv:2507.13499 [cs.SE] https://arxiv.org/abs/2507.13499
[63] Michał Madera and Rafał Tomoń. 2017. A case study on machine learning model for code review expert system in
software engineering. In 2017 Federated Conference on Computer Science and Information Systems (FedCSIS). 1357–1363.
doi:10.15439/2017F536
[64] Shane McIntosh, Yasutaka Kamei, Bram Adams, and Ahmed E Hassan. 2016. An empirical study of the impact of
modern code review practices on software quality. Empirical Software Engineering 21, 5 (2016), 2146–2189.
[65] Ziya Mukhtarov, Mannan Abdul, Mokhlaroyim Raupova, Javid Baghirov, Osama Tanveer, Haluk Altunel, and Eray
Tüzün. 2023. Towards Better Code Reviews: Using Mutation Testing to Improve Reviewer Attention. In 2023 IEEE/ACM
International Conference on Software and System Processes (ICSSP). 92–96. doi:10.1109/ICSSP59042.2023.00020
[66] Atharva Naik, Marcus Alenius, Daniel Fried, and Carolyn Rose. 2025. CRScore: Grounding Automated Evaluation of
Code Review Comments in Code Claims and Smells. arXiv:2409.19801 [cs.SE] https://arxiv.org/abs/2409.19801
[67] Miroslaw Ochodek, Miroslaw Staron, Wilhelm Meding, and Ola Söder. 2022. Automated Code Review Comment
Classification to Improve Modern Code Reviews. In Software Quality: The Next Big Thing in Software Engineering and
Quality, Daniel Mendez, Manuel Wimmer, Dietmar Winkler, Stefan Biffl, and Johannes Bergsmann (Eds.). Springer
International Publishing, Cham, 23–40.
[68] Charlie Parker, Sam Scott, and Alistair Geddes. 2019. Snowball sampling. SAGE research methods foundations (2019).
[69] Debalina Ghosh Paul, Hong Zhu, and Ian Bayley. 2024. Benchmarks and Metrics for Evaluations of Code Generation:
A Critical Review. In 2024 IEEE International Conference on Articial Intelligence Testing (AITest). 87–94. doi:10.1109/
AITest62860.2024.00019
[70] Chanathip Pornprasit and Chakkrit Tantithamthavorn. 2024. Fine-tuning and prompt engineering for large language
models-based code review automation. Inf. Softw. Technol. 175, C (Nov. 2024), 12 pages. doi:10.1016/j.infsof.2024.107523
[71] Chanathip Pornprasit, Chakkrit Tantithamthavorn, Patanamon Thongtanunam, and Chunyang Chen. 2023. D-ACT:
Towards Diff-Aware Code Transformation for Code Review Under a Time-Wise Evaluation. In 2023 IEEE International
Conference on Software Analysis, Evolution and Reengineering (SANER). 296–307. doi:10.1109/SANER56733.2023.00036
[72] Huilian Sophie Qiu, Bogdan Vasilescu, Christian Kästner, Carolyn Egelman, Ciera Jaspan, and Emerson Murphy-Hill.
2022. Detecting Interpersonal Conflict in Issues and Code Review: Cross Pollinating Open- and Closed-Source
Approaches. In 2022 IEEE/ACM 44th International Conference on Software Engineering: Software Engineering in Society
(ICSE-SEIS). 41–55. doi:10.1145/3510458.3513019
[73] Mohammad Masudur Rahman, Chanchal K. Roy, and Raula G. Kula. 2017. Predicting Usefulness of Code Review
Comments Using Textual Features and Developer Experience. In 2017 IEEE/ACM 14th International Conference on
Mining Software Repositories (MSR). 215–226. doi:10.1109/MSR.2017.17
[74] Shadikur Rahman, Umme Ayman Koana, and Maleknaz Nayebi. 2022. Example Driven Code Review Explanation.
In Proceedings of the 16th ACM / IEEE International Symposium on Empirical Software Engineering and Measurement
(Helsinki, Finland) (ESEM ’22). Association for Computing Machinery, New York, NY, USA, 307–312. doi:10.1145/
3544902.3546639
[75] Nishrith Saini and Ricardo Britto. 2021. Using Machine Intelligence to Prioritise Code Review Requests. In 2021
IEEE/ACM 43rd International Conference on Software Engineering: Software Engineering in Practice (ICSE-SEIP). 11–20.
doi:10.1109/ICSE-SEIP52600.2021.00010
[76] Jaydeb Sarker, Sayma Sultana, Steven R. Wilson, and Amiangshu Bosu. 2023. ToxiSpanSE: An Explainable Toxicity
Detection in Code Review Comments . In 2023 ACM/IEEE International Symposium on Empirical Software Engineering
and Measurement (ESEM). IEEE Computer Society, Los Alamitos, CA, USA, 1–12. doi:10.1109/ESEM56168.2023.
10304855
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
35
[77] Jaydeb Sarker, Asif Kamal Turzo, Ming Dong, and Amiangshu Bosu. 2023. Automated Identification of Toxic Code
Reviews Using ToxiCR. ACM Trans. Softw. Eng. Methodol. 32, 5, Article 118 (July 2023), 32 pages. doi:10.1145/3583562
[78] Neela Sawant and Srinivasan H. Sengamedu. 2023. Code Compliance Assessment as a Learning Problem. In Proceedings
of the 45th International Conference on Software Engineering: Software Engineering in Practice (Melbourne, Australia)
(ICSE-SEIP ’23). IEEE Press, 445–454. doi:10.1109/ICSE-SEIP58684.2023.00046
[79] Oussama Ben Sghaier and Houari Sahraoui. 2023. A Multi-Step Learning Approach to Assist Code Review. In 2023
IEEE International Conference on Software Analysis, Evolution and Reengineering (SANER). 450–460. doi:10.1109/
SANER56733.2023.00049
[80] Oussama Ben Sghaier, Martin Weyssow, and Houari Sahraoui. 2025. Harnessing Large Language Models for Curated
Code Reviews . In 2025 IEEE/ACM 22nd International Conference on Mining Software Repositories (MSR). IEEE Computer
Society, Los Alamitos, CA, USA, 187–198. doi:10.1109/MSR66628.2025.00039
[81] Qianhua Shan, David Sukhdeo, Qianying Huang, Seth Rogers, Lawrence Chen, Elise Paradis, Peter C. Rigby, and
Nachiappan Nagappan. 2022. Using nudges to accelerate code reviews at scale. In Proceedings of the 30th ACM Joint
European Software Engineering Conference and Symposium on the Foundations of Software Engineering (Singapore,
Singapore) (ESEC/FSE 2022). Association for Computing Machinery, New York, NY, USA, 472–482. doi:10.1145/
3540250.3549104
[82] Shipra Sharma and Balwinder Sodhi. [n. d.]. Using Stack Overflow content to assist in code review. Software:
Practice and Experience 49, 8 ([n. d.]), 1255–1277. arXiv:https://onlinelibrary.wiley.com/doi/pdf/10.1002/spe.2706
doi:10.1002/spe.2706
[83] Shu-Ting Shi, Ming Li, David Lo, Ferdian Thung, and Xuan Huo. 2019. Automatic code review by learning the
revision of source code. In Proceedings of the Thirty-Third AAAI Conference on Articial Intelligence and Thirty-First
Innovative Applications of Articial Intelligence Conference and Ninth AAAI Symposium on Educational Advances
in Articial Intelligence (Honolulu, Hawaii, USA) (AAAI’19/IAAI’19/EAAI’19). AAAI Press, Article 603, 8 pages.
doi:10.1609/aaai.v33i01.33014910
[84] Ohiduzzaman Shuvo, Parvez Mahbub, and Mohammad Masudur Rahman. 2023. Recommending Code Reviews
Leveraging Code Changes with Structured Information Retrieval. In 2023 IEEE International Conference on Software
Maintenance and Evolution (ICSME). 194–206. doi:10.1109/ICSME58846.2023.00029
[85] Jing Kai Siow, Cuiyun Gao, Lingling Fan, Sen Chen, and Yang Liu. 2020. CORE: Automating Review Recommendation
for Code Changes. In 2020 IEEE 27th International Conference on Software Analysis, Evolution and Reengineering
(SANER). 284–295. doi:10.1109/SANER48275.2020.9054794
[86] Behjat Soltanifar, Atakan Erdem, and Ayse Bener. 2016. Predicting Defectiveness of Software Patches. In Proceedings
of the 10th ACM/IEEE International Symposium on Empirical Software Engineering and Measurement (Ciudad Real,
Spain) (ESEM ’16). Association for Computing Machinery, New York, NY, USA, Article 22, 10 pages. doi:10.1145/
2961111.2962601
[87] Tao Sun, Jian Xu, Yuanpeng Li, Zhao Yan, Ge Zhang, Lintao Xie, Lu Geng, Zheng Wang, Yueyan Chen, Qin Lin, et al.
2025. Bitsai-cr: Automated code review via llm in practice. In Proceedings of the 33rd ACM International Conference on
the Foundations of Software Engineering. 274–285.
[88] Tao Sun, Jian Xu, Yuanpeng Li, Zhao Yan, Ge Zhang, Lintao Xie, Lu Geng, Zheng Wang, Yueyan Chen, Qin Lin,
Wenbo Duan, Kaixin Sui, and Yuanshuo Zhu. 2025. BitsAI-CR: Automated Code Review via LLM in Practice. In
Proceedings of the 33rd ACM International Conference on the Foundations of Software Engineering (Clarion Hotel
Trondheim, Trondheim, Norway) (FSE Companion ’25). Association for Computing Machinery, New York, NY, USA,
274–285. doi:10.1145/3696630.3728552
[89] Xunzhu Tang, Kisub Kim, Yewei Song, Cedric Lothritz, Bei Li, Saad Ezzini, Haoye Tian, Jacques Klein, and Tégawendé F.
Bissyandé. 2024. CodeAgent: Autonomous Communicative Agents for Code Review. In Conference on Empirical
Methods in Natural Language Processing. https://api.semanticscholar.org/CorpusID:267412469
[90] Yida Tao and Sunghun Kim. 2015. Partitioning Composite Code Changes to Facilitate Code Review. In 2015 IEEE/ACM
12th Working Conference on Mining Software Repositories. 180–190. doi:10.1109/MSR.2015.24
[91] Patanamon Thongtanunam, Chanathip Pornprasit, and Chakkrit Tantithamthavorn. 2022. AutoTransform: automated
code transformation to support modern code review process. In Proceedings of the 44th International Conference on
Software Engineering (Pittsburgh, Pennsylvania) (ICSE ’22). Association for Computing Machinery, New York, NY,
USA, 237–248. doi:10.1145/3510003.3510067
[92] Panya Trakoolgerntong, Tao Xiao, Masanari Kondo, Chaiyong Ragkhitwetsagul, Morakot Choetkiertikul, Pattaraporn
Sangaroonsilp, and Yasutaka Kamei. 2025. AILINKPREVIEWER: Enhancing Code Reviews with LLM-Powered Link
Previews. arXiv:2511.09223 [cs.SE] https://arxiv.org/abs/2511.09223
[93] Rosalia Tufano, Alberto Martin-Lopez, Ahmad Tayeb, Ozren Dabić, Sonia Haiduc, and Gabriele Bavota. 2025. Deep
Learning-Based Code Reviews: A Paradigm Shift or a Double-Edged Sword?. In Proceedings of the IEEE/ACM 47th
International Conference on Software Engineering (Ottawa, Ontario, Canada) (ICSE ’25). IEEE Press, 1640–1652. doi:10.
, Vol. 1, No. 1, Article . Publication date: February 2026.

36
Taufiqul Islam khan, Shaowei Wang, Haoxiang Zhang, and Tse-Hsun Chen
1109/ICSE55347.2025.00060
[94] Rosalia Tufano, Simone Masiero, Antonio Mastropaolo, Luca Pascarella, Denys Poshyvanyk, and Gabriele Bavota.
2022. Using pre-trained models to boost code review automation. In Proceedings of the 44th International Conference
on Software Engineering (Pittsburgh, Pennsylvania) (ICSE ’22). Association for Computing Machinery, New York, NY,
USA, 2291–2302. doi:10.1145/3510003.3510621
[95] Rosalia Tufano, Luca Pascarella, Michele Tufano, Denys Poshyvanyk, and Gabriele Bavota. 2021. Towards Automating
Code Review Activities. In Proceedings of the 43rd International Conference on Software Engineering (Madrid, Spain)
(ICSE ’21). IEEE Press, 163–174. doi:10.1109/ICSE43902.2021.00027
[96] Anderson Uchôa, Caio Barbosa, Daniel Coutinho, Willian Oizumi, Wesley K. G. Assunção, Silvia Regina Vergilio,
Juliana Alves Pereira, Anderson Oliveira, and Alessandro Garcia. 2021. Predicting Design Impactful Changes in
Modern Code Review: A Large-Scale Empirical Study. In 2021 IEEE/ACM 18th International Conference on Mining
Software Repositories (MSR). 471–482. doi:10.1109/MSR52588.2021.00059
[97] Yuki Ueda, Takashi Ishio, Akinori Ihara, and Kenichi Matsumoto. 2019. Mining Source Code Improvement Patterns
from Similar Code Review Works. In 2019 IEEE 13th International Workshop on Software Clones (IWSC). 13–19.
doi:10.1109/IWSC.2019.8665852
[98] Chen Wang, Xiaoyuan Xie, Peng Liang, and Jifeng Xuan. 2017. Multi-Perspective Visualization to Assist Code Change
Review. In 2017 24th Asia-Pacic Software Engineering Conference (APSEC). 564–569. doi:10.1109/APSEC.2017.66
[99] Dong Wang, Yuki Ueda, Raula Gaikovina Kula, Takashi Ishio, and Kenichi Matsumoto. 2021. Can we benchmark Code
Review studies? A systematic mapping study of methodology, dataset, and metric. Journal of Systems and Software
180 (2021), 111009. doi:10.1016/j.jss.2021.111009
[100] Kaixin Wang, Tianlin Li, Xiaoyu Zhang, Chong Wang, Weisong Sun, Yang Liu, and Bin Shi. 2025. Software Development
Life Cycle Perspective: A Survey of Benchmarks for Code Large Language Models and Agents. arXiv preprint
arXiv:2505.05283 (2025).
[101] Kaixin Wang, Tianlin Li, Xiaoyu Zhang, Chong Wang, Weisong Sun, Yang Liu, and Bin Shi. 2025. Software Development
Life Cycle Perspective: A Survey of Benchmarks for Code Large Language Models and Agents. ArXiv abs/2505.05283
(2025). https://api.semanticscholar.org/CorpusID:278394402
[102] Luqiao Wang, Yangtao Zhou, Huiying Zhuang, Qingshan Li, Di Cui, Yutong Zhao, and Lu Wang. 2024. Unity is
strength: Collaborative llm-based agents for code reviewer recommendation. In Proceedings of the 39th IEEE/ACM
International Conference on Automated Software Engineering. 2235–2239.
[103] Min Wang, Zeqi Lin, Yanzhen Zou, and Bing Xie. 2019. CoRA: Decomposing and Describing Tangled Code Changes
for Reviewer. In 2019 34th IEEE/ACM International Conference on Automated Software Engineering (ASE). 1050–1061.
doi:10.1109/ASE.2019.00101
[104] Song Wang, Chetan Bansal, Nachiappan Nagappan, and Adithya Abraham Philip. 2019. Leveraging Change Intents for
Characterizing and Identifying Large-Review-Effort Changes. In Proceedings of the Fifteenth International Conference on
Predictive Models and Data Analytics in Software Engineering (Recife, Brazil) (PROMISE’19). Association for Computing
Machinery, New York, NY, USA, 46–55. doi:10.1145/3345629.3345635
[105] Ruiyin Wen, David Gilbert, Michael G. Roche, and Shane McIntosh. 2018. BLIMP Tracer: Integrating Build Impact
Analysis with Code Review. In 2018 IEEE International Conference on Software Maintenance and Evolution (ICSME).
685–694. doi:10.1109/ICSME.2018.00078
[106] Weixiang Yan, Haitian Liu, Yunkun Wang, Yunzhe Li, Qian Chen, Wen Wang, Tingyu Lin, Weishan Zhao, Li Zhu, Hari
Sundaram, and Shuiguang Deng. 2024. CodeScope: An Execution-based Multilingual Multitask Multidimensional
Benchmark for Evaluating LLMs on Code Understanding and Generation. In Proceedings of the 62nd Annual Meeting of
the Association for Computational Linguistics (Volume 1: Long Papers), Lun-Wei Ku, Andre Martins, and Vivek Srikumar
(Eds.). Association for Computational Linguistics, Bangkok, Thailand, 5511–5558. doi:10.18653/v1/2024.acl-long.301
[107] Lanxin Yang, Jinwei Xu, Yifan Zhang, He Zhang, and Alberto Bacchelli. 2023. EvaCRC: Evaluating Code Review
Comments. In Proceedings of the 31st ACM Joint European Software Engineering Conference and Symposium on the
Foundations of Software Engineering (San Francisco, CA, USA) (ESEC/FSE 2023). Association for Computing Machinery,
New York, NY, USA, 275–287. doi:10.1145/3611643.3616245
[108] Ying Yin, Yuhai Zhao, Yiming Sun, and Chen Chen. 2023. Automatic Code Review by Learning the Structure
Information of Code Graph. Sensors 23, 5 (2023). doi:10.3390/s23052551
[109] Yongda Yu, Guohao Shi, Xianwei Wu, Haochuan He, XueMing Gu, Qianqian Zhao, Kui Liu, Qiushi Wang, Zhao Tian,
Haifeng Shen, and Guoping Rong. 2025. Fine-Tuning LLMs to Analyze Multiple Dimensions of Code Review: A
Maximum Entropy Regulated Long Chain-of-Thought Approach. arXiv:2509.21170 [cs.SE] https://arxiv.org/abs/2509.
21170
[110] Yongda Yu, Lei Zhang, Guoping Rong, Haifeng Shen, Jiahao Zhang, Haoxiang Yan, Guohao Shi, Dong Shao, Ruiqi
Pan, Yuan Li, Qiushi Wang, and Zhao Tian. 2025. Distilling Desired Comments for Enhanced Code Review with
Large Language Models. arXiv:2412.20340 [cs.SE] https://arxiv.org/abs/2412.20340
, Vol. 1, No. 1, Article . Publication date: February 2026.

A Survey of Code Review Benchmarks and Evaluation Practices in Pre-LLM and LLM Era
37
[111] Zhengran Zeng, Ruikai Shi, Keke Han, Yixin Li, Kaicheng Sun, Yidong Wang, Zhuohao Yu, Rui Xie, Wei Ye, and
Shikun Zhang. 2025. Benchmarking and Studying the LLM-based Code Review. arXiv:2509.01494 [cs.SE] https:
//arxiv.org/abs/2509.01494
[112] Guoliang Zhao, Daniel Alencar da Costa, and Ying Zou. 2019. Improving the pull requests review process using
learning-to-rank algorithms. Empirical Software Engineering 24, 4 (2019), 2140–2170. doi:10.1007/s10664-019-09696-8
[113] Zelin Zhao, Zhaogui Xu, Jialong Zhu, Peng Di, Yuan Yao, and Xiaoxing Ma. 2023. The Right Prompts for the Job:
Repair Code-Review Defects with Large Language Model. arXiv:2312.17485 [cs.SE] https://arxiv.org/abs/2312.17485
, Vol. 1, No. 1, Article . Publication date: February 2026.
