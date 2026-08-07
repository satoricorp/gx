Concerns identiﬁed in code review: A ﬁne-grained, faceted classiﬁcation
Sanuri Gunawardenaa,∗, Ewan Temperoa, Kelly Blincoeb
aSchool of Computer Science, The University of Auckland, Auckland, New Zealand
bDepartment of Electrical, Computer, and Software Engineering, The University of Auckland, Auckland, New Zealand
Abstract
Context: Code review is a valuable software process that helps software practitioners to identify a variety of
defects in code. Even though many code review tools and static analysis tools used to improve the eﬃciency
of the process exist, code review is still costly.
Objective: Understanding the types of defects that code reviews help to identify could reveal other means
of cost improvement. Thus, our goal was to identify defect types detected in real-world code reviews, and
the extent to which code review can be beneﬁted from defect detection tools.
Method: To this end, we classiﬁed 417 comments from code reviews of 7 OSS Java projects using thematic
analysis.
Results: We identiﬁed 116 defect types that we grouped into 15 groups to create a defect classiﬁcation.
Additionally, 38% of these defects could be automatically detected accurately.
Conclusion: We learnt that even though many capable defect detection tools are available today, a sub-
stantial amount of defects that can be detected automatically, reach code review. Also, we identiﬁed several
code review cost reduction opportunities.
Keywords:
code review, code inspection, concerns, types, defects, decisions, manual classiﬁcation,
detection method, detection expertise, non-programmers
1. Introduction
Code review, where code contributions of a soft-
ware practitioner are evaluated by other members
of the team, is a common practice in software
teams [1, 2, 3]. This practice brings signiﬁcant ad-
vantages such as improved code quality and knowl-
edge sharing [4, 5, 6, 7, 8].
During code review,
reviewers may identify functional issues with the
code, suggest changes to the design, check for ad-
herence to coding conventions, suggest alternate
implementations, or simply provide praise [1]. De-
spite its popularity, a common complaint is the cost
of performing code review, in particular the time it
takes [1, 8, 9, 10, 11]. To this end, tools have been
proposed that automatically provide reviewers the
information they need [12], enable eﬀective collab-
oration [13, 14, 15, 16], and detect defects using
∗Corresponding author
Email addresses:
sanuri.gunawardena@auckland.ac.nz (Sanuri
Gunawardena), e.tempero@auckland.ac.nz (Ewan
Tempero), k.blincoe@auckland.ac.nz (Kelly Blincoe)
static analysis [17, 18, 19, 20]. Static analysis tools
would typically be run prior to the code review so
reviewers can focus on identifying more complex po-
tential issues during code review. Despite the wide
availability of such tools, developers have reported
low adoption of these tools due to false positives
and other barriers [21, 22] and complaints are still
being made on the cost of code review [11].
Our goal is to further reduce the cost of code re-
view without compromising its reliability. To this
end, code review cost reduction opportunities could
be identiﬁed through examination of defect types
usually identiﬁed during code reviews. By under-
standing which defects are common, we may tailor
cost-reduction approaches to such defects. Previ-
ous work has investigated defects identiﬁed during
code reviews [23, 24, 25, 8, 26, 18, 27]. Several clas-
siﬁcations of code review ﬁndings and changes have
been proposed (e.g., [27] [24]). However, the current
classiﬁcations are not ﬁne-grained enough to enable
identiﬁcation of all automation opportunities.
For example, the comment “seems like these need
Preprint submitted to Elsevier
August 23, 2022

some javadocs” extracted from CROP [28] refers to
a set of missing Javadocs. This can be automati-
cally detected by performing a text search to see if
all methods include a Javadoc each. The comment
“Javadoc? It looks like there’s just enough here to
pass checkstyle?”, also from CROP [28] refers to
missing Javadoc information.
To determine that
Javadoc lacks information, an understanding of the
what the Javadoc is stating and an understanding
of the related code and its purpose are essential.
Thus, the detection of such concerns cannot be au-
tomated. However, both these defects fall into the
same “Comments” class in the most ﬁne-grained
classiﬁcation existing currently, CRAM [27]. This
means that CRAM does not diﬀerentiate the de-
tection automate-ability of defects. Thus, existing
classiﬁcations do not provide an easy way to iden-
tify all code review cost reduction opportunities.
In this study, we develop a ﬁne-grained classiﬁca-
tion for defects reported in code review comments
with the goal of identifying automation opportu-
nities. In addition, we examine which code review
defects could have been automatically identiﬁed us-
ing popular existing tools to better understand the
adoption of existing automation techniques. We fo-
cus on all potential problems with the code identi-
ﬁed in code review comments that may or may not
be actual defects because our goal is to understand
the entire eﬀort involved in code review, not just
the comments that lead to code changes. Thus, in
the remainder of this paper, we will refer to these
potential problems as concerns.
Our study was guided by the following research
questions:
RQ1:
What types of concerns are identiﬁed in
real-world code reviews?
RQ2:
To what extent could the concerns iden-
tiﬁed during code review be automatically and
accurately identiﬁed?
The “accurately” in RQ2 refers to detecting a
concern without producing false positives and false
negatives.
We use “accurately” with the same
meaning throughout the article.
We extracted and classiﬁed 417 comments from
a subset of code review discussions of 7 open-source
projects obtained from CROP [28], using thematic
analysis [29]. The main contribution of this study
is a detailed code review concern classiﬁcation (see
Figure 5) consisting of 116 concern types that are
grouped into 15 concern groups.
To understand which of the code review concerns
could have been identiﬁed using existing automated
tools, we considered ﬁve popular static analysis
tools and tested whether or not they could have
automatically identiﬁed the concerns described in
the code review comments.
Our results suggest
that many concerns (22%) could have been identi-
ﬁed using existing automated tools, showing there
is signiﬁcant potential to reduce the eﬀort of code
review through better adoption of tools. The de-
tection of another 16% of concerns were not sup-
ported by the tools we studied but also could have
been automated accurately as new heuristics could
be easily deﬁned to enable their detection. In our
discussion, we reﬂect on the types of concerns that
are identiﬁed during code review and suggest future
research avenues to improve the cost of code review
including both automation and other ways.
The remainder of this paper is organized as fol-
lows: We describe the related work in Section 2.
Section 3 describes our methodology and the re-
search questions.
Section 4 presents the results.
Section 5 discusses the implications of our results,
future research directions, and the threats to valid-
ity. Finally, Section 6 oﬀers a brief conclusion.
2. Related Work
2.1. Modern Code Review
Modern code review (MCR) is a light-weight,
tool-assisted,
and asynchronous review process
where code written by a developer is read through
and evaluated by other team members [8, 1, 2, 3].
MCR tools highlight the diﬀerences between two
versions of code and provide collaboration fea-
tures.
Popular MCR tools include Gerrit Code
Review [13], Collaborator [14], Crucible [15], and
Github Code Review [16].
The use of static analysis tools to reduce code
review eﬀort and increase code review quality has
also been explored both in research and prac-
tice [17, 18, 19, 20].
There is a wide range of
static analysis tools available (Eg: SonarQube [30],
Checkstyle [31], FindBugsTM [32], PMD [33], In-
telliJ IDEA Code Inpection [34]) and they have a
wide range of capabilities. Static analysis tools can
either fully or partially automate the detection of
defects, prior to code review. The extent to which
the detection can be automated and the accuracy
2

of detection, however, depends on the nature of the
defect being detected.
A defect that can be well-deﬁned, such as “trail-
ing white-spaces” (unwanted whitespaces at the end
of code statements) can be detected accurately us-
ing static analysis tools. False positives and false
negatives are not produced in such cases. For ex-
ample Checkstyle can detect trailing whitespaces
using its “NoWhitespaceAfter” rule [35] which sim-
ply looks for the presence of a whitespace at the
end of each code line. Thus, the detection of such
defects can be fully automated.
Defects inherently manual to assess [36], such
as “non-self-explanatory identiﬁer names”, do not
have objective deﬁnitions, and cannot be accurately
detected by tools; a tool can only detect a set of
potentially bad identiﬁer names, check an identi-
ﬁer name against a regular expression (class names
should comply with a naming convention (Sonar-
Qube) [30], Type name(Checkstyle) [37], Class
naming convention(IntelliJ) [38]), and recommend
potentially better alternatives [39, 40, 41] but a hu-
man (i.e. the reviewer) must make the decision of
whether the identiﬁer name is bad enough to raise
the concern. The detection of such defects can only
be partially automated as described above.
Thus, it is impossible to completely remove hu-
man intervention in code review. However, human
intervention can be reduced by identifying the con-
cerns that can be fully detected using static analysis
tools and giving the responsibility of such concerns
to static analysis tools rather than reviewers. One
of our attempts is to identify those defects and un-
derstand the role of automation in code reviews to
pave the path to reducing code review cost.
2.2. The Modern Code Review Process
Here, we describe a typical modern code review
process. A review starts when a developer modiﬁes
the original codebase in the repository and submits
a new patch in the form of a commit. This person
is typically called the author of the code. Next, one
or more other developers of the project will review
the submitted source code and provide feedback
in the form of comments.
These developers are
the reviewers.
Then, the author will modify the
current revision of the patch if required, according
to the review feedback and submit the improved
code as a new revision.
There may be multiple
iterations of this step.
If agreed, the proposed
revision is merged into the code base. If rejected,
Repository
Original
code base
Author
makes
changes
Reviewer
comments
Revision
Approved
revision
Review
cycle
Revision
comments
Figure 1: A typical modern code review process (adapted
from Paixao et al. [28])
the revision is abandoned.
2.3. Defect Classiﬁcations
There are several existing code review defect clas-
siﬁcations: Beller et al. [42], Runeson et al. [26],
Bacchelli et al. [8], Lei et al. [25], Siy et al. [23],
M¨antyl¨a and Lassenius [24], and Panichella and Za-
ugg [27].
The two main goals of these classiﬁca-
tions are to identify the value of code review and
tool needs of reviewers.
All of these studies cat-
egorize code review defects into broad categories,
Panichella and Zaug’s Code Review chAnges Model
(CRAM) having the most ﬁne-grained categories
among the classiﬁcations. One of the lowest level
categories of CRAM is “Comments”, deﬁned as
“Explanations of complex code fragments, classes,
methods. Issues include wrongly placed comments,
missing comments, missing or wrong Javadoc, etc.”.
The categories of these classiﬁcations do not pro-
vide a way to identify ﬁne-grained automation op-
portunities, but only high-level tool needs. For ex-
ample, a documentation comment that is unneces-
sarily placed on a private method can be detected
by a simple text search and this opportunity for
using automation cannot be identiﬁed by looking
at the categories of CRAM or any other existing
classiﬁcation due to their high level of abstraction.
Therefore, our classiﬁcation was designed to carry
concern categories that are further ﬁne-grained so
3

Selecting
dataset
7 Java
projects
from CROP
Sampling and Data Extraction
Minimum 385 comments,
minimum 55 comments per
project
Random sampling and
concern extraction
417 extracted concerns
Labeling of
concerns and
grouping concern
types
Concern types
classification
Defining
research
questions
Research
questions
Calculating sample size for
unknown or large population
Classifying
concerns
according to the
detection method
Detection method
classification
Figure 2: Methodology
that automation opportunities can be conveniently
identiﬁed. Since CRAM is the most detailed classi-
ﬁcation existing, we compared our classiﬁcation to
CRAM in detail in section 5.4.
There are many other manual defect classiﬁca-
tions available [43, 44, 45, 46, 47, 48, 49, 50, 51,
52, 53, 54, 55, 56, 57, 58].
They extract defects
from various resources such as test reports and cus-
tomer feedback reports [53]. None of these classiﬁ-
cations extract defects from code review data and
thus do not qualify to answer our research ques-
tions. There are also many automatic defect clas-
siﬁcation methods that make defect classiﬁcation
eﬃcient. They either use an existing classiﬁcation
such as ODC [59, 60], classify defects broadly as
bug or other request [61, 62], or create their own
classiﬁcation [63, 64, 65].
Similar to the manual
defects classiﬁcations we discussed before, none of
these classiﬁcations either provide ﬁne-grained de-
fect details or classify defects based on their detec-
tion automate-ability.
Thus, the classiﬁcation we created is the ﬁrst of
its kind. To properly encode data and to not be bi-
ased by the existing classiﬁcations, we designed our
concern classiﬁcation by applying thematic analysis
to code review data instead of building on existing
classiﬁcations.
There are several uses to the concern classiﬁca-
tion we created. It can be used as a reference to
create organizational standards and checklists for
future code reviews.
Also, the classiﬁcation pro-
vides a new way of viewing the code under review if
rigorousness of the review is a goal: While review-
ing code for all aspects of code can be overwhelm-
ing, considering a single concern group at a time
modularizes and simpliﬁes the code review process.
Additionally, the classiﬁcation we produced is the
only existing classiﬁcation of its kind that can po-
tentially help to reduce code review cost eﬀectively.
3. Methodology
This section describes the methodology we fol-
lowed in our study. Fig.2 illustrates this process.
3.1. Research Questions
The focus of this study is to identify the types
of concerns raised during MCRs and the part that
existing static analysis tools can play in detecting
concerns to identify code review cost reduction
opportunities. Our research questions are inspired
by this focus:
RQ1:
What types of concerns are identiﬁed in
real-world code reviews?
RQ2:
To what extent could the concerns iden-
tiﬁed during code review be automatically and
accurately identiﬁed?
3.2. Data Set
A formally published open-source data set of code
review data intended for software engineering re-
searchers and practitioners, “Code Review Open
Platform (CROP)” was used for this study [28].
This data set contains data coming from the Ger-
rit Code Review and thus, review comments are
available in both ﬁle-level and code line-level. This
study analysed only the Java projects from CROP,
which includes 7 of its 11 technical projects. Also,
we analyzed code review comments on Java ﬁles
only.
Fig.3 illustrates the structure of the CROP data
set. At the highest level, it consists of a set of fold-
ers, each folder corresponding to a single project.
Within each project folder, there is a set of sub-
folders, each corresponding to an entire code review
discussion on a unique patch.
4

…
Patch 1 discussion

…
Project 1

…

Revision 1 review
comment 1 Concern 1
comment 2 Concern 2, Concern 3
comment 3
comment 4 Concern 4
comment 5
 …
Patch 2 discussion
Revision 2 review
Project 2
Figure 3: The Structure of CROP Data Set
A unique patch undergoes one or more concep-
tual revisions as a result of the code review process,
before it is accepted or abandoned. Each revision of
the patch produces a revision review ﬁle containing
the author-reviewer communication made on the re-
vision of the patch. Therefore, a single discussion
folder contains one or more revision review ﬁles.
A revision review ﬁle may or may not contain
code review comments based on the outcome of the
code review performed on its patch revision. A code
review comment is a comment made on a code line
or a ﬁle in the patch by the reviewer. It may iden-
tify something that needs to be changed (e.g., a
functional defect), question something that works
as is but could be improved (e.g., a design defect or
non-adherence to code conventions), or suggest al-
ternatives (e.g., a diﬀerent algorithm). A comment
may also acknowledge a point raised by another re-
viewer, or praise the code being reviewed.
Our interest is in comments that may lead to a
change in the code. As noted in the introduction,
we examine these comments to identify concerns.
The CROP data set has a total of 50,959 code
review discussions and 144,906 revision reviews, ac-
cording to the CROP website [66]. The 7 projects
we selected for this study had 22,859 code review
discussions and 63,051 revision reviews. Table 3.2
summarizes the characteristics of the data available
for these 7 projects in the CROP dataset [28] (pop-
ulation) and the sample we extracted as described
below. The project descriptions were obtained from
Paixao and Maia [67].
3.3. Sampling and Data Extraction
We calculated a sample size that would be repre-
sentative of the concerns in the data set while still
being reasonable for manual analysis. As described
above, not all code review comments are describing
concerns. An initial observation of discussions was
conducted to get an idea of the density of comments
that contain concerns. Only 29 out of 100 randomly
selected comments were reporting concerns. Thus,
it was not possible to calculate the exact popula-
tion size of concerns in our data set. We, therefore,
chose to be conservative and calculated a sample
size for a unknown, large population size. Studies
have shown that as the population increases, to ob-
tain a sample with a conﬁdence level of 95% and a
margin of error of 5%, the sample size increases at
a diminishing rate and remains relatively constant
at 384.16 after a population size of 1,000,000 [68].
Thus, we considered 385 as the target sample size.
Line-level Comment
File-level Comment
Figure 4: File-level Comment Vs. Line-level Comment and
Revision Review File Content
Our goal was to extract a similar number of con-
cerns from each project to represent them equally.
Thus, from each project, ﬁrst, a discussion was
randomly selected and concerns were extracted by
manually reading through all revision review ﬁles
5

Population
Sample
Project
Description
Language
Time span
Number of
discussions
Number of
revisions
Number of
developers
Time span
Number of
discussions
Number of
revisions
Number of
developers
1
Couchbase’s driver
implementation
in Java
Java
01-12 to 11-17
917
2638
35
02-12 to 10-17
40
124
14
2
Low-level API
mostly used by
java-client
Java
04-14 to 11-17
841
2301
18
04-14 to 11-17
45
156
8
3
Building blocks
for user interfaces
in Eclipse
Java
02-13 to 11-17
4756
14118
290
03-13 to 11-17
44
147
31
4
Integration of jgit
into the Eclipse IDE
Java
10-09 to 11-17
5337
13231
184
10-09 to 08-17
34
72
22
5
Java
implementation
of Git
Java
09-09 to 11-17
5384
14043
251
04-10 to 11-17
32
54
21
6
C/C++ IDE for
linux developers
Java
06-12 to 11-17
5105
15337
82
08-12 to 03-17
33
82
16
7
Implementation
of a memory
caching system
Java
05-10 to 07-17
519
1383
36
05-10 to 10-11
42
88
11
Table 1: Population and sample summaries
under the selected discussion. Next, another dis-
cussion was randomly selected and the same ex-
traction process was executed on its revision review
ﬁles. This was repeated until at least 55 concerns
(target sample size 385 / Number of projects 7)
were extracted from each project.
Since we ex-
tracted all concerns from each randomly selected
discussion, each project had diﬀerent numbers of
extracted concerns within the range of 56 and 71.
We chose to extract all concerns from a single dis-
cussion, rather than randomly sampling concerns
across each project, to ensure that for each discus-
sion, we captured each type of concern that was
identiﬁed by the reviewers i.e. we extracted both
ﬁle-level and line-level concerns in each discussion.
Fig.4 shows the structure of a revision review ﬁle
and the diﬀerence between line and ﬁle-level com-
ments.
When a single comment was referring to more
than one concern, the comment was split in a way
that each part of the comment was referring to only
one concern. Then each part of the comment was
labeled independently.
There were 22 comments
that we had to split into 2 or more concerns. Simi-
larly, if there were more than one comment referring
to the same concern, they were grouped as a single
comment. There were 9 such concerns.
3.4. Labeling
Concerns
and
Grouping
Concern
Types
To answer RQ1, we built a classiﬁcation of con-
cern types by applying thematic analysis [29] to the
code review data described above. The concerns de-
scribed in code review comments and the associated
code were analyzed to identify the type of concern
and each concern was labeled with a suitable con-
cern type label.
Concern types were then grouped based on com-
mon themes that emerged among them. For exam-
ple, “why dropping public on these? we don’t seem
to do that elsewhere” was labeled “Missing access
modiﬁer”, “Fields that sub-classes might want to
update should be either protected or have protected
setters (protected methods), for example fStatis-
ticsData.” was labeled “Unexpected access modi-
ﬁer”, and both these concerns were grouped under
“Modiﬁers” because they were both concerns re-
lated to modiﬁers used in code.
Some concerns were lacking information on their
root cause.
In such cases, we assigned them la-
bels with higher-levels of abstraction. For example,
“This will not generate a pom ﬁle with the correct
dependencies.” refers to a functional issue in the
code, but does not reveal what causes this issue.
We labeled this concern as “Unexpected logic and
functionality”.
These labels together with other
root cause-level labels constitute the lowest layer
of our classiﬁcation.
We used consistent wording across the labels we
created. “unexpected” was used to mean that the
reviewer observed something diﬀerent compared to
what they were expecting in the code but it may or
may not be incorrect (e.g., unexpected functional-
ity and logic). “Better <ITEM> exists” was used
to mean that the reviewer felt that the current sub-
ject of interest can be further improved (e.g., bet-
ter design exists). “Unnecessary” was used to mean
that something not required is present in the code
(Eg: Unnecessary modiﬁer).
“missing” was used
to mean that a required something is not present
6

in the code (Eg: Missing comment). “unconven-
tional” was used to mean that something in the
code is against the project-conventions (Eg: Uncon-
ventional identiﬁer name, Unconventional license
header pattern).
During the labeling process, the ﬁrst author was
mainly responsible of labeling the concerns. After
100 concerns were labeled and grouped, to ensure
the reliability of the classiﬁcation, the other two
authors categorized 10 concerns each in 5 rounds
(total 50 comments), into the existing labels created
by the ﬁrst author or new labels as required. Each
round was followed by a discussion among the three
authors where the labels were changed, updated,
and moved into diﬀerent groups appropriately.
The concern extraction resulted in 397 concerns.
Once this was completed, additional 20 concerns
(equivalent to 5% of the original sample size) were
extracted by the ﬁrst author from the data set
and categorized to check for theoretical saturation.
All of these concerns could be categorized into the
existing concern types, suggesting that saturation
may have reached.
These concerns were also in-
cluded in the classiﬁcation bringing the total num-
ber of concerns to 417.
Project
1
2
3
4
5
6
7
Extracted number
of concerns
60
56
58
61
71
67
56
Unclear concerns
4
2
4
2
1
4
0
Discussions
including concerns
3
9
12
8
2
5
10
Total evaluated
discussions
40
45
44
34
32
33
42
1 = couchbase-java-client, 2 = couchbase-jvm-core,
3 = eclipse.platform.ui, 4 = egit, 5 = jgit, 6 =
org.eclipse.linuxtools, 7 = spymemcached
Table 2: Concern Extraction Summary
Table 3.4 lists the number of concerns extracted
from each project, the number of concerns that were
diﬃcult to understand i.e.
unclear concerns, the
number of discussions that included concerns, and
the number of total discussions that were evaluated
in each project. A concern was marked as unclear
if it was concluded to be incomprehensible among
the authors. “Document in exec and connect that
exec is called, then connect. It doesn’t make sense
that this is the behavior, but it is.” is an example
of an unclear concern.
As a ﬁnal reliability check, each author other
than the ﬁrst author categorized 15 randomly se-
lected concerns from the 417 categorized concerns,
followed by a discussion. At this stage, there were
no conﬂicts among the authors regarding the label-
ing and grouping. This gave us the conﬁdence that
the classiﬁcation was reliable.
3.5. Detection Method Classiﬁcation
To answer RQ2, we classiﬁed each concern based
on whether it can be automatically detected ac-
curately or it must be manually detected by ap-
plying the criteria given in table 3.4.
The code
containing each concern was analyzed using ﬁve
popular static analysis tools: SonarQube, Check-
style, FindBugsTM, PMD, and IntelliJ IDEA Code
Inspection, and Grammarly to detect grammar-
related issues in API documentation.
The tools
were used to check whether the concern could be
detected without producing any false positives or
false negatives. These tools were selected because
they have been used commonly and therefore im-
proved throughout the past decade based on user
feedback.
If these tools could not detect the concern, the
concern was discussed among the authors to de-
cide if accurate automation was currently possible.
For example, “Do you think we should be address-
ing some of this in the class documentation rather
than repeating it next to each method?
It feels
a bit repetitive.” referring to identical text blocks
in API documentation comments can be detected
by performing a simple text search in the code for
identical Javadocs. Thus, it should be categorized
as an automatically detected concern even though
the standard rules of the studied tools could not
detect it. The complete list of concerns with their
“detection method” classiﬁcation, evidence on the
automatability of concern detection, and the ver-
sions of tools studied are available in our supple-
mentary material [69].
This classiﬁcation of automatically and manually
detected concerns was done in parallel with the con-
Category
Deﬁnition of Category
Automatically
detected
These concerns can be accurately detected
automatically with existing tools or new heuristics
could be easily deﬁned to enable their detection.
”Accurately” implies that false positives and false
negatives will never be produced during
the detection of the concern. (Eg: Missing braces,
Code lines longer than 80 characters)
Manually
detected
The detection of these concerns cannot be automated
accurately. (Eg: Poor design, Low class cohesion)
Table 3: Automatically Vs. Manually Detected Concerns
7

cern type labelling described above to ensure that
the level of detail of the concern type labels was
suitable to diﬀerentiate between automatically and
manually detected concerns. Our aim was that for
each concern type label, all concerns with that label
should also share the same automatically or manu-
ally detected label. In nearly all cases, by creating
low-level labels, this was possible. There were nine
of the 116 concern types where it was not desirable
to diﬀerentiate between the concerns that could be
automatically detected and those that needed to
be manually detected. For example, “The ID cor-
respond to the package in which this class is embed-
ded.” labeled “Documentation Grammar” could be
automatically detected by Grammarly while “The
list of the traces to add in the tree”, also labeled
“Documentation Grammar” could not be automat-
ically detected using Grammarly [70].
Splitting
such concern types further would have increased
the number of concern types undesirably, and have
made the classiﬁcation unnecessarily complex.
4. Results
4.1. Concern Types Classiﬁcation
We identiﬁed 116 concern types, which are listed
in Fig.5. The concern type that had the largest dis-
tribution was Trailing whitespace (22/417, 5%).
Trailing whitespace is an unnecessary single space
introduced at the end of the code statements during
code changes. This is diﬀerent from Unnecessary
whitespace which refers to a whitespace that is not
required, but present within a code statement and
not at the end of the statement.
The next most common concern type was Miss-
ing Documentation Info (16/417, 4%). Many of
the concern types (43% or 50/116) occurred only
once in the sample. The complete descriptions of
each concern type and how each concern was clas-
siﬁed into these concern types is available in our
online supplementary material [69].
These 116 concern types were grouped into 15
concern groups as shown in Fig.5. Table 4.1 lists the
distribution of concern types and concerns within
each concern group. The column, “Number of Con-
cern Types” provides the count and percentage of
concern types under each group. “Number of Con-
cerns” column provides the count and percentage
of individual concerns under each group. Concern
groups are discussed below in the order of most
number of concerns to least number of concerns in
the sample.
concern group
Number of Concern
Types (Total 116)
Number of Concerns
(Total 417)
Count
%
Count
%
API Documentation
12
10.4
68
16.3
Implementation
19
16.4
62
14.9
Appearance
11
9.5
61
14.6
Logic and functionality
15
12.9
59
14.1
Design
16
13.8
31
7.4
Modiﬁers
6
5.2
25
6.0
Comments
8
6.9
20
4.8
Identiﬁer naming
4
3.4
19
4.6
Errors, warnings,
and logging
7
6
14
3.4
Performance
1
0.9
14
3.4
Header comments
7
6
14
3.4
Annotations
4
3.4
11
2.6
Test cases
3
2.6
10
2.4
Literals
2
1.7
8
1.9
Threads
1
0.9
1
0.2
Table 4: Distribution of Concerns and Concern Types
API Documentation represents concerns re-
lated to a special group of comments, API docu-
mentation comments and contains the most con-
cerns in the sample (68/417 concerns , 16.3%).
“Here the timeout is in milliseconds, but the API
Javadoc suggested seconds. That’s why I wanted
to clarify it.” labeled “Unexpected documentation
info” is an example of a comment that belongs to
this group.
Implementation
(62/417
concerns,
14.9%)
reports
shortcomings
of
the
current
way
of
implementation.
“My
preference
is
to
use
org.eclipse.core.runtime.Platform.getWS()
and
getOS()” labeled “Better library use exists” is
an example which suggests that a better library
method call can be used instead of what is cur-
rently used in the implementation to achieve the
same outcome.
Appearance (61/417 concerns, 14.6%) includes
concerns related to the way code looks, such as un-
necessary braces, trailing whitespaces, and incor-
rect indentation. “Style-nit: In JGit we don’t put
curly braces around single line statements.” labeled
“Unnecessary braces” is an example.
Logic and functionality groups concerns re-
lated to logic that cause the functionality to be dif-
ferent from the expected behaviour (59/417 con-
cerns, 14.1%).
“This will not generate a pom
ﬁle with the correct dependencies.” labeled “Un-
expected logic and functionality” is an example of
such a concern.
Design (31/417 concerns, 7.4%) includes con-
cerns related to the architecture of the software un-
der review. An examples is “I’m pretty sure this
8

- Unexpected comment info (Co1)
- Unexpected TODO comment info (Co2)
- Insufficient comment info (Co3)
- Missing comment (Co4)
- Missing TODO comment (Co5)
- Unnecessary comment (Co6)
- Unnecessary TODO comment (Co7)
- Unclear comment (Co8)
- Uninformative error message (E1)
- Missing error (E2)
- Missing warning (E3)
- Missing logging (E4)
- Unnecessary logging (E5)
- Generic exception throw (E6)
- printStackTrace call (E7)
- Unexpected copyright dates (H1)
- Unexpected copyright author (H2)
- Missing copyright author (H3)
- Incomplete license header description (H4)
- Missing license header (H5)
- Unconventional license header pattern (H6)
- Header comment typo (H7)

Concerns
Annotations
Comments
Performance
Test cases
Implementation
Identifier naming
Design
Literals
Logic and functionality
Modifiers
Errors, warnings, and logging
Appearance
- Unnecessary annotation (A1)
- Missing annotation (A2)
- Unexpected annotation scope (A3)
- Unexpected annotation option (A4)
- Unexpected indentation (Ap1)
- Missing blank line (Ap2)
- Unnecessary blank line (Ap3)
- Trailing whitespace (Ap4)
- Unnecessary whitespace (Excludes Trailing Whitespaces) (Ap5)
- Missing braces (Ap6)
- Unnecessary braces (Ap7)
- Long code line (Ap8)
- Text reflow (Ap9)
- Missing reflow (Ap10)
- Unconventional order of class member placement (Ap11)
- Unnecessary interface (D1) - Unnecessary class (D2) -
Unnecessary method (D3) - Low class cohesion (D4) - Low
package cohesion (D5) - Low method cohesion (D6) - Missing
extend (D7) - Missing interface method (D8) - Reused class
specific constant (D9) - Protected method in non-API class (D10)
- Public method in non-API class (D11) - Candidate for "static"
(D12) - Common logic in child classes (D13) - External reference
of current object (D14) - Enum used instead of inheritance (D15)
- Better design exists (D16)
- Unconventional identifier name (I1)
- Non-self-explanatory identifier name (I2)
- Better identifier name exists (I3)
- Identifier name typo (I4)
- Unnecessary documentation info (J1) - Missing
documentation info (J2) -Unexpected
documentation info (J3) - Repeated
Documentation info (J4) - Documentation
grammar (J5) - Documentation typo (J6) -
Unclear documentation (J7) - Missing
documentation (J8) - Unconventional
documentation pattern (J9) - Unnecessary
documentation on private method (J10) - Missing
import (Jo11)
- Magic number (L1)
- Non-externalized string (L2)
- Unnecessary check (Lo1) - Missing check (Lo2) -
Unexpected check (Lo3) - Unexpected value (Lo4) -
Inconsistent logic and functionality (Lo5) - Missing logic
and functionality (Lo6) - Unexpected logic and
functionality (Lo7) - Unexpected output appearance (Lo8) -
Unexpected output information (Lo9) - Unexpected sentinel
value (Lo10) - Unexpected separator (Lo11) - Observable
closed too early (Lo12) - Better observable type exists
(Lo13) - Duplicate functionality (Lo14) - Unsafe cast (Lo15)
- Unexpected access modifier (M1)
- Unnecessary access modifier (M2)
- Unconventional order of modifiers (M3)
- Inconsistent order of modifiers (M4)
- Missing access modifier (M5)
- Missing modifier (M6)
- Unnecessary modifier (M7)
- Poorly performing logic (P1)
- Unexpected test case (T1)
- Tool-specific test case (T2)
- Missing test case (T3)
- Better library use exists (Im1) - Inconsistent library use (Im2) -
Better API exists (Im3) - Misplaced try (Im4) - Nesting too deep
(Im5) - Method output not reused (Im6) -Simplifiable logic (Im7)
-Unreadable code (Im8) -Extractable code (Im9) -Inflexible logic
(Im10) - Unexpected type (Im11) - Unnecessary parameter
(Im12) - Better implementation exists (Im13) - Class name
prefixed with package path (Im14) - Code duplication (Im15) -
Unexpected order of Enum elements (Im16) - Explicit null
initialization (Im17) - Unused item (Im18) - Dead code (Im19)
API Documentation
Header Comments
- Thread not joined (Th1)
Threads
Figure 5: Code Review Concerns Classiﬁcation
class doesn’t belong here at all. It looks like it’s for
testing??” labeled “Low package cohesion”.
Implementation, Logic and functionality, and
Design-related concerns included both statement-
level concerns (Lo1, D12, Im11) and code block-
level concerns (Lo6, D1, Im13).
The block-level
concerns were the complex and larger concerns that
were labeled with a higher-level of abstraction. It
is up to the author to investigate the root causes of
such concerns and correct them.
Modiﬁers (25/417 concerns, 6.0%) groups con-
cerns related to modiﬁers used in Java.
“Should
this be public?
I’m not quite sure what I’d do
with this.” labeled “Unexpected access modiﬁer”
is an example. Modiﬁer-related concerns were ei-
ther access modiﬁer-related (M1, M2), non-access
modiﬁer-related (M6, M7,) or related to both
(M3,M4).
“ﬁnal” and “static” in Java were ex-
amples of non-access modiﬁers. The concerns re-
lated to both modiﬁer types were regarding the or-
der in which the modiﬁers were placed with respect
to each other.
Comments (20/417 concerns, 4.8%) group con-
tains concerns related to regular inline comments
that describe the code. “I ﬁnd this comment pretty
confusing. What is ‘the changes feed’. I know this
isn’t in this commit, but needs to be somewhere.”
refers to an Unclear comment which falls under this
category. Comment-related concerns were either re-
lated to TODO comments (Co2, Co5, Co7) or nor-
mal code comments (Co1, Co3). TODO comments
were in fact code comments that documented parts
of implementation that was yet to be implemented.
Concerns related to TODO comments were much
less (3 concerns) compared to the normal code com-
ments (17 concerns).
Identiﬁer naming (19/417 concerns, 4.6%) in-
cludes concerns related to identiﬁer names such as
“This is a meaningless variable name. I can’t tell
what it is without digging further.” that was labeled
“Non-self-explanatory identiﬁer name”.
Errors, warnings, and logging (14/417 con-
cerns, 3.4%) includes concerns related to either er-
rors, warnings or logging i.e. concerns related to
9

the error communication of a program. “I realize
this isn’t a new change, but we shouldn’t be calling
e.printStackTrace anywhere.” labeled “printStack-
Trace call” is an example of a concern that belongs
to this group.
Performance
(14/417,
3.4%)
represents
performance-related concerns of program code.
“This one is likely to be a lot more expensive.” is
an example. Additionally, this is one of the groups
that contain the least variety of concern types: just
1 concern type labeled “Poorly-performing logic”.
Header comments (14/417, 3.4%) was another
group that demanded a separate concern group due
to the variety of header comment-related concern
types we found in the sample. “Comma here too
;)” referring to a header comment typo is an ex-
ample. Observing several random ﬁles of the data
set showed that header comments are written us-
ing either Javadoc notation or a regular comment
notation in practice.
Annotations (11/417 concerns, 2.6%) repre-
sents annotation-related concerns. “We should add
@noimplement” labeled “Missing annotation” is an
example of a concern that is annotation-related.
This concern group is language speciﬁc and may
not be applicable to a language other than Java.
Test cases (10/417 concerns, 2.4%) contain con-
cerns related to test cases used to test the soft-
ware.
“also assertFalse(second.hasSubscribers())”
referring to a Missing test case is an example of a
test case-related concern.
Literals (8/417 concerns, 1.9%) contain con-
cerns related to literals used in code. An example
of Literal-related concern is “maybe extract as a
constant?” referring to a Magic number.
Threads (1/417 concerns, 0.2%) included only
1 concern, thus only 1 concern type: Thread not
joined (“Looks good. But there is another issue in
GitRepositoriesViewRemoteHandlingTest. Check-
out job needs to be joined.”).
4.2. Manually Vs.
Automatically Detected Con-
cerns
We deﬁned “automatically detected concerns” as
concerns of which the detection can be automated
accurately.
The remaining concerns were catego-
rized as “manually detected concerns”.
In Fig.5
automatically detected concerns are marked in ital-
ics while manually detected concerns are marked in
bold.
From the 417 concerns, the detection of 22% con-
cerns were supported by the tools we studied and
another 16% (36/417 concerns) were not supported
by the tools we studied but could be detected au-
tomatically and accurately by deﬁning new heuris-
tics easily.
Thus, a total of 38% (157/417) con-
cerns were categorized under “automatically de-
tected concerns” and the remaining 62% (260/417)
of concerns were categorized under “manually de-
tected concerns”.
Out of the 116 concern types,
32% (37/116) of the concern types were categorized
under automatically detected concerns only, 60%
(70/116) of the types under manually detected con-
cerns only, and 8% (9/116) contained both types
of concerns (ﬁg.5).
The concern types that in-
cluded both automatically and manually detected
concerns were Unnecessary annotation, Missing er-
ror, Missing logging, Code duplication, Documenta-
tion grammar, Documentation typo, Missing Doc-
umentation info, Unnecessary Documentation info,
and Missing modiﬁer.
An example of an automatically detected con-
cern is “missing braces” of a control structure which
can occur in Java when it is project’s convention
to use mandatory braces marking control structure
bodies, but the developer decided not to use braces
because the control structure body was only a sin-
gle code line (Java allows this). Many tools pro-
vide checks to enforce control structure braces (Eg:
SonarQube - Control structures should use curly
braces [30], Checkstyle - NeedBraces [71], PMD -
ControlStatementBraces [72], IntelliJ - Control ﬂow
statement without braces [73]) and can automati-
cally detect such cases without producing any false
positives or negatives.
Fig.6.a illustrates the distribution of automati-
cally detected concerns among the concern groups.
From the 157 automatically detected concerns,
the majority of concerns belonged to Appearance
(61/157 concerns, 39%). “Trailing whitespace” was
the most common (22/157 concerns, 14%) concern
type among automatically detected concerns. The
lowest distribution of automatically detected con-
cerns was under Annotations group (1/157, 0.6%).
The only Annotations-related, automatically de-
tected concern was “This is a test plugin, NLS
warnings are not enabled, so these annotations are
not needed.” labeled “Unnecessary annotation”.
An example of a manually detected concern
is “I think we should ﬁnd a diﬀerent name for
this, because an unsubscribe could be because of
a timeout but may also be for diﬀerent reasons.
What about something like isActive()? if unsub-
scribed then its just not active and can be dropped
10

57
43
31
31
20
18
15
14
10
10
9
1
1
0
20
40
60
Logic and functionality
Implementation
Desgin
Javadoc
Comments
Modifiers
Identifier naming
Performance
Annotations
Test cases
Errors, warnings, and logging
Header Comments
Threads
61
37
19
13
8
7
5
4
2
1
0
20
40
60
80
Appearance
Javadoc
Implementation
Header comments
Literals
Modifiers
Errors, warnings, and logging
Identifier naming
Logic and functionality
Annotations

Concerns Count
Concerns Count
b) Manually detected Concerns
a) Automatically detected Concerns
Concerns Group

Total = 157
Total = 260
Figure 6: High-level categorization of concerns - Detection Method
for various reasons.” labeled “Non-self-explanatory
identiﬁer name”.
Even though existing technol-
ogy can be used to identify potentially bad iden-
tiﬁer names [39], whether the identiﬁer name is
bad enough to raise the concern depends on the
reviewer.
Fig.6.b illustrates the distribution of manually
detected concerns among concern groups.
From
the 260 manually detected concerns, the highest
distributed was related to Logic and functionality
(57/260, 22%), and the lowest was Threads and
Header comments-related concerns (1/260, 0.4%).
The only manually detected, Header comment-
related concern was “We generally assign copyright
to either the author or employer, depending on
what is most appropriate.
don’t attribute copy-
right to Eclipse Egit Team, instead use your name
and email or your company” labeled “Unexpected
copyright author”.
5. Discussion
5.1. RQ1: Concerns Identiﬁed in Real-world Code
Reviews
The study revealed 116 concern types identiﬁed
in real-world code reviews and 15 concern groups
that they could be grouped into.
Fig.5 is an
overview of the concern types that were present in
the sample. The API Documentation group con-
tained the most number of concerns (68/417 con-
cerns, 16.3%).
The least concerns were observed
in the Threads group (1/417, 0.2%). Implementa-
tion group contained the highest number of concern
types (19/116 concern types, 16.4%). The least va-
riety (1/116, 0.9%) of concern types was in Per-
formance (“Poorly performing logic”) and Threads
(“Thread not joined”).
The most common con-
cern type in the sample was “Trailing whitespace”,
present 22 times (22/417, 5%). 43% (50/116 con-
cern types) of the concern types were present only
once in the sample.
We did not ﬁnd any security-related concerns in
the sample. Several possible reasons behind this are
explained in the literature. The diversity of the se-
curity issues is overwhelming for general developers
and the lack of eﬀectiveness of the security assur-
ance tools is a challenge and thus it usually leads to
security review avoidance or bringing in expensive
external resources to the organization.
Also, the
lack of expertise and security being a non-functional
requirement adds to its invisibility during code re-
views. [74, 75, 76]. There are many existing static
analysis tools that are designed to support security
defect detection [77, 78, 79, 80, 3]. However, due
to barriers like usability problems, these tools are
not well-adopted [81]. The lack of usable and use-
ful security tools was reﬂected in a survey where
most programmers ideally wanted to see security
warnings on static analysis tools [82].
We found only one concern related to multi-
threaded programming: Thread not joined.
The
presence of this one concern shows that at least
one of the projects is making use of multi-threaded
programming.
The rarity of concurrency con-
cerns in the sample could be due to the nature
of the projects, the nature of the considered sam-
ple, or the complexity of concurrency program-
ming compared to sequential programming.
A
study at Microsoft has shown that Concurrency
bugs take on average several days to detect, repro-
duce, debug and ﬁx [83].
There are many forms
11

of support to help concurrency bugs detection in
code [84, 85, 86, 87, 88]. However, their spurious
results and eﬀectiveness still need to be improved
further [89].
Some of the concern types discovered made it
apparent that coding conventions can be diﬀerent
from project to project. For example, the two in-
verse conventions “unnecessary braces” (in project
“Jgit”) stating that braces should not be used
in control structures with single statement bod-
ies and “missing braces” (in project “Egit”) that
mandates to always use braces in control struc-
tures were discovered in the sample implying dif-
ferences of project-level conventions.
“unneces-
sary braces” also shows that sometimes project
conventions can be contrasting to the generally
accepted standards (Oracle Java Coding Conven-
tions [90]). Header comments is another element
that we observed diﬀerences of. While some pro-
grammers used Javadoc notation for header com-
ments, some used regular comments. There are nu-
merous language-speciﬁc [90] as well as general [91]
coding conventions created by diﬀerent authors, or-
ganizations, and projects [92].
Due to their ob-
jective nature, most coding conventions are well-
supported by existing tools, compared to security
and concurrency bugs [93, 71, 94].
The experience of producing this classiﬁcation
disclosed that diﬀtools (code review tools) partially
automate the detection of two of the simplest con-
cern types, unnecessary trailing white spaces and
blank lines, as a side-eﬀect.
The introduction of
the two concern types could be due to the speciﬁc
development tools or their conﬁgurations used in
the considered projects. Since the diﬀview of code
review tools highlights newly added blank spaces in
red, it is convenient for the reviewers to spot them.
Otherwise, the usual goals of code review tools are
to support review information management and col-
laboration [16, 13].
Additionally, we discovered two forms of subjec-
tivity related to code review decisions.
Concern
types such as “Nesting too deep” and “Long code
line” contain a subjective component that depends
on project conventions i.e. project subjectivity. For
example, the maximum depth of nesting allowed
would be based on the project conventions and may
diﬀer from project to project. The other form is
the reviewer subjectivity that has an eﬀect on the
identiﬁcation of concern types such as “Non-self-
explanatory identiﬁer name”. The meaningfulness
and suitability of an identiﬁer name and the re-
quirement to improve an identiﬁer name quality de-
pends on the reviewer’s perspective and may diﬀer
from reviewer to reviewer. The presence of these
two forms of subjectivity is a major barrier for eﬀec-
tively reducing code review cost using static analy-
sis tools.
5.2. RQ2: Code Review Automation
The results showed that MCR identiﬁes more
than concerns that are inherently manual to as-
sess. A substantial number of concerns (157/417
concerns, 38%) and concern types (46/116 concern
types, 40%) could be detected automatically with-
out producing any false positives or false negatives.
Attempting to capture these concerns by consis-
tently using defect detection tools prior to code
review would reduce the code review cost consid-
erably.
Of the 38% of concerns that can be automati-
cally detected, 22% are supported by the popular
tools we studied, and 16% concerns were not sup-
ported by those tools but could be accurately de-
tected automatically by deﬁning new rules.
The
latter type of concerns are often concerns that are
project-speciﬁc as general rules of the tools cannot
detect them.
“We need to clean up the logging.
Some is JDK logging, some is spy logging.” is an
example.
Many existing tools allow teams to create qual-
ity proﬁles, which could be used to handle project-
speciﬁc rules [95].
Software teams could create
project-speciﬁc custom rules and add them to a
proﬁle at the tool conﬁguration stage. Then, the
quality proﬁle can be assigned to a project, so that
those rules will run only against that particular
project. Software teams can create such a proﬁle
for each project they work on. We believe that hav-
ing custom rules can standardize the project con-
ventions further and ensure consistency throughout
the project. Quality proﬁles may help to manage
the resulting large number of custom rules. Future
work can investigate this further.
The automatically detected concerns in our clas-
siﬁcation do not include concerns that can be par-
tially detected using existing tools because such
concerns cause false positives and false negatives
in static analysis tool results and therefore result
in low adoption of static analysis tools [96, 18, 97].
The low adoption of static analysis tools lead to
manual defect detection which is a barrier to re-
ducing the code review cost. By improving existing
12

static analysis tools further (to minimize false pos-
itives and other barriers), we believe that the cost
of code review can be further reduced. However,
for now, we recommend using tools only for those
concerns that can be identiﬁed accurately. Partial
automation and manual code review should be com-
pared in future to determine which is more cost ef-
fective.
To provide tool recommendations, we looked at
the number of concern types that each tool could
support from the concern types we identiﬁed. In-
telliJ IDEA Code Inspection detected the highest
number of concern types (16) accurately.
Gram-
marly detected the 2 concerns that it is designed to
detect: API Documentation Typo and API Docu-
mentation Grammar, when the Javadocs were ex-
tracted and tested on it. SonarQube, Checkstyle,
PMD, and FindBugsTM supported 15, 14, 10, and
one concern types respectively. However, this does
not mean that one tool has more features than the
other.
For example, we observed that FindBugsTM sup-
ported many complex design-related checks but
false positives and false negatives were inevitable
in such cases. Thus, they were not counted as au-
tomatically detected concerns. An example is “Re-
turn value of method without side eﬀect is ignored”
rule supported by FindBugsTM [98]. Such checks
can be called “partial automation” because human
intervention is still required to make the ﬁnal deci-
sion on whether a concern is present or not.
5.3. Opportunities for Reducing Code Review Cost
Conducting this study helped us identify several
opportunities for reducing the code review cost.
5.3.1. Automatically detected concerns
Automatically detected concerns evidently still
reach code review sessions.
Consistent and thor-
ough use of existing defect detection tools prior to
code reviews can help prevent this, reducing the
code review cost by a substantial amount (22%
concerns of the considered sample). Another 16%
concerns were not supported by the popular tools
we studied, but their detection could be automated
accurately. Using new custom rules together with
the standard rules of existing tools, and using these
tools consistently and thoroughly, code review cost
can be reduced. We have identiﬁed the concerns
that can be automatically detected i.e. the oppor-
tunities where human eﬀort and resultant cost can
be prevented (Concern types in italics in ﬁg.5).
5.3.2. Partially Automated Concerns
There are still many concerns identiﬁed in code
reviews that can be partially detected by existing
defect detection tools. These are currently catego-
rized under “manually-detected concerns”. We are
expecting to diﬀerentiate these from fully manually
detected concerns in our future studies so that hu-
man eﬀort required in such cases also can be mini-
mized. Additionally, investigating the eﬀect of par-
tial automation on code review cost is an interesting
future research avenue.
5.3.3. Functional Concerns
From the sample, 26% concerns were functional
concerns.
Some of them may have been de-
tected comparatively inexpensively by having a
high-quality test suite instead of putting code re-
view eﬀort into the task.
An example of such a
concern is “what if the ICommitMessageProvider
returns null? ...” that could be detected by hav-
ing a test case that invokes the related code with
ICommitMessageProvider being null.
5.3.4. Code Review and Expertise
Category
Deﬁnition of Category
No
programming
expertise
These concerns can be detected with no programming
expertise or training of any kind. The review
instructions and the source code representation may
have to be modiﬁed to support non-programmers.
(Eg: “Detect grammatical mistakes in the following
documentation texts extracted from a software system.”)
No
programming
expertise
with training
These concerns can be detected with no programming
expertise, but requires training on other concepts
related to the concern. The review instructions and
the source code representation may have to be modiﬁed
to support non-programmers (Eg: “Detect low-quality
error messages from the following list of error messages
extracted from a software system” will instruct the
reviewer to detect low-quality error messages such as
“Something went wrong!”. The reviewer needs
training on what a low-quality error message is.)
Programming
expertise
These concerns require programming expertise and the
original source code to be detected. (Eg: Logical and
functional concerns, Design-related concerns)
Table 5: Detection Expertise Categories
Code review cost (time and eﬀort) has been a
persistent problem in code review [99, 100, 101,
102, 103, 104, 105, 106, 107, 108, 109, 10, 110,
111, 2, 112]. Many researchers have attempted to
solve this problem by suggesting diﬀerent guide-
lines [113, 114, 115, 116, 117, 118, 119, 120, 121,
122, 123, 97, 124]. Several common measures are
suggested in these guidelines: determining review-
ers based on the available human resources rather
than ﬁxing the number of reviewers by process, pro-
moting tool usage, and not having review meetings
13

or having them only when required. All these mea-
sures achieve cost reductions by ultimately reduc-
ing the number of programmers that have to be
involved in the review process. On the contrary, we
suggest, reducing the amount of money spent on
human resources or in other words, recruiting less
expensive human resources (i.e., non-programmers)
may help to reduce the code review cost.
The domain-speciﬁc expertise and technological-
expertise are two major factors considered in mod-
ern reviewer recommendation systems because of
their eﬀect on the code review quality [125, 126,
127, 128, 129, 130, 122]. Thus, it is common for
domain and technological experts to review code
in practice. However, sometimes, the code review
decision is as simple as spotting a meaningless vari-
able name such as “objRef” and thus might be
able to be identiﬁed by a non-programmer as well.
By delegating such simple code review tasks that
are usually considered as nitpicking by develop-
ers to non-programmers who are less expensive,
the code review cost can be further reduced.
In
fact, at Microsoft, a preliminary “code improve-
ment” review round is conducted prior to the ac-
tual code review session to identify maintainability
issues [8]. This can be an excellent point at which
non-programmers can contribute to code review by
identifying nitpicks that they can so that program-
ming experts can focus on more complex issues in
code.
Organizations with cross-functional teams
can utilize their existing human resources more ef-
fectively in the code review process and improve the
programming productivity of expert programmers
by reducing the workload imposed on them by code
review. Another possibility is to crowd-source these
tasks as micro-jobs.
As the ﬁrst step to study the feasibility of having
non-programmers perform some elements of code
review, we created a non-programmer-centric con-
cerns classiﬁcation, “Detection Expertise”.
The
categorization was done according to the criteria
deﬁned in table 5.3.4.
In creating the categories
for this classiﬁcation we considered 1) whether a
non-programmer can detect the concern or pro-
gramming expertise (i.e.
domain and technologi-
cal expertise acquired by a programmer over time)
is required, and 2) whether providing the non-
programmer instructions on what to detect is suﬃ-
cient or training on concepts related to the concern
is required. Additionally, while verifying this clas-
siﬁcation in the future studies, we will also have to
consider whether the original source code and re-
view instructions can be used for review or modiﬁed
representations of the source code and modiﬁed in-
structions understandable by non-programmers are
more appropriate to get the task done eﬀectively by
the non-programmers.
28
18
18
4
8
7
3
1
3
1
6
1
3
2
0
5
10
15
20
25
30
Appearance
Javadoc
Implementation
Identifier naming
Literals
Modifiers
Errors, warnings, and logging
Comments
Design
Header Comments
Automatically detected concerns
Manually detected concerns
2
1
57
42
29
22
18
17
14
10
10
9
8
1
0
10
20
30
40
50
60
70
Logic and functionality
Implementation
Design
Javadoc
Modifiers
Comments
Performance
Annotations
Test cases
Identifier naming
Errors, warnings, and logging
Threads
Automatically detected concerns
Manually detected concerns

33
19
12
2
1
6
1
0
5
10
15
20
25
30
35
Appearance
Javadoc
Header comments
Errors, warnings, and logging
Annotations
Automatically detected concerns
Manually detected concerns
Concerns Group
Concerns Group
Concerns Group
Concerns Count
Concerns Count
Concerns Count
a) No Expertise (Total 74)
b) No Expertise with Training (Total 103)
c) Programming Expertise (Total 240)
Figure 7: High-level categorization of concerns - Detection
Expertise
To classify, each code review comment and the
associated code were studied by the ﬁrst author to
understand the nature of the concern: whether pro-
gramming expertise is a prerequisite and whether
concept-training is required to detect the presence
of the concern. Once the ﬁrst author had completed
this classiﬁcation, the other two authors categorized
15 randomly selected comments from the sample.
A Fleiss Kappa value of 0.624 implying a substan-
tial agreement was obtained among the 3 authors.
Obtaining a perfect agreement is diﬃcult for this
kind of categorization where the authors had to de-
14

pend on their subjective understanding about the
concerns reported in the classiﬁed comments. How-
ever, following this categorization, the authors had
a discussion where they reached a consensus. Based
on this, the ﬁrst author re-categorized all comments
one last time.
The examples given below explain
the 3 categories in this classiﬁcation further.
The detection of concerns related to the logic of
a program requires the ability to read and under-
stand the code, thus, should be categorized under
“Programming Expertise”.
A missing Javadoc (documentation comment) on
a method can be detected by a person who can
identify method blocks. This does not require the
code to be understood, rather training on identi-
fying method blocks based on the indentation and
the placement of curly brackets. Thus, the concern
“missing documentation” should be categorized un-
der the “No Expertise with Training” category.
To detect the code lines longer than 80 charac-
ters, the reviewer does not need to read and un-
derstand the code i.e. requires no programming ex-
pertise and does not have any complex concepts to
learn i.e. requires no training. The reviewer simply
has to identify the text lines longer than 80 charac-
ters in the code ﬁle. Thus, “long code line” should
be categorized under the “No Expertise” category.
Fig.7 depicts the distribution of concerns in each
expertise level. The complete classiﬁcation is avail-
able online [69].
Programming expertise category with the
most number of concerns in the sample (240/417
concerns,
57%)
contained
12
concern
groups
(Fig.7.c). Logic and Functionality (59/240, 25%)
included the most number of concerns. This cate-
gory is the only instance that Logic and function-
ality, Performance, and Test cases groups can be
observed. This is because they are the most com-
plex aspects of programming and thus they need
programming expertise to detect.
Naturally, the
majority of concerns that require expertise to iden-
tify could not be detected automatically either, due
to their “inherently manual to assess” nature.
No programming expertise with training
category (103/417 concerns, 25%) (Fig.7.b) con-
tained 10 concern groups of which the majority was
Appearance-related (28/103, 27%). The manually
detected concern types that belong to this category
are 3.8% of the sample (16/417 concerns) which is
another promising cost reduction opportunity be-
cause they can be detected with no programming
expertise and with non-extensive concept training.
No programming expertise category (74/417
concerns, 18%) contained 5 concern groups. The
majority
of
concerns
were
Appearance-related
(33/74, 45%)(Fig.7.a). The most eﬀective code re-
view cost reduction opportunity here is in the man-
ually detected concerns because their detection can-
not be automated accurately and therefore human
involvement is required. 1.7% (7/417 concerns) of
the concerns in the sample are manually detected
and require no programming expertise to detect. If
a person with no expertise (i.e a less expensive hu-
man resource) can detect those concerns, the overall
code review cost can be eﬀectively reduced.
There is also another 37% (154/417 concerns)
concerns in the sample that belong to the last two
categories explained above and are also automati-
cally detected concerns. For organizations that pre-
fer manual code review, this is a signiﬁcant oppor-
tunity to reduce the cost of their code review pro-
cess by utilizing less expensive human resources.
Thus, we estimate that in total 42.5% of the con-
cerns in our sample could be identiﬁed by some-
one without programming expertise, which could
enable another avenue for cost savings in software
code review. Programmers will agree that the con-
cerns in this group are mostly “nitpicks”. A study
conducted at Microsoft implies that reviewers tend
to miss more complex issues in code due to these
nitpicks [8]. Another study has shown that ﬁxing
soft maintenance issues lowers the cost of future
changes [23]. Therefore, nitpicks are important to
be discovered. Also, they do carry a certain iden-
tiﬁcation cost. However, this cost is unknown. Fu-
ture studies should examine the cost of identifying
such concerns to see whether they are worth being
delegated to non-programmers. Future studies also
should validate this classiﬁcation.
5.4. A Comparison of Classiﬁcations
Code review defect classiﬁcations existing today
are the products of an evolutionary process of re-
searchers attempting to make classiﬁcations more
informative and more inclusive of the many pos-
sible code review defect types.
Most of these
classiﬁcations either adapt a previous classiﬁcation
or integrate a previous classiﬁcation.
For exam-
ple, the Panichella and Zaugg classiﬁcation [27]
published recently integrates the Beller classiﬁca-
tion [42] which in turn adapts the M¨antyl¨a and
Lassenius classiﬁcation [24]. We compared our clas-
siﬁcation to the most recent and most detailed ex-
15

isting classiﬁcation, Panichella and Zaugg classiﬁ-
cation, also called CRAM [27].
Since our lowest level categories were more ﬁne-
grained than CRAM, we grouped our concern types
under the lowest-level categories of the Panichella
and Zaugg classiﬁcation (the details of this group-
ing is available in our supplementary material [69]).
Based on the deﬁnitions of the low-level categories
of CRAM, we found it diﬃcult to ﬁt 10.34% of our
concern types into their classiﬁcation categories.
From the 15 concern groups in our classiﬁcation 14
overlapped with the categories of CRAM. The re-
maining group, “Annotations” could not be placed
in CRAM because CRAM data set did not contain
any code review comments on code-related anno-
tations. However, according to M¨antyl¨a classiﬁca-
tion that CRAM is indirectly based on, “Annota-
tions” belongs in the “Language supported docu-
mentation” high-level category. The other concern
types for which we could not ﬁnd matching CRAM
low-level categories are better implementation ex-
ists, inﬂexible logic, better observable type exists,
observable closed too early, unexpected logic and
functionality, unexpected sentinel value, and unex-
pected separator.
Conversely, CRAM contained some low-level cat-
egories that we did not include. We did not diﬀer-
entiate, for example, “semantic duplication” and
“duplicate code” whereas in CRAM they were dif-
ferentiated.
“Semantic dead code” and “Dead
code” were another example. CRAM did not dif-
ferentiate API Documentation-related changes and
Comment-related changes whereas we did.
Also,
CRAM does not separate automatically detected
concerns from manually detected concerns.
Our classiﬁcation was an attempt to cleanly sepa-
rate automatically detected concerns and manually
detected concerns.
During the process we learnt
that this was not entirely possible.
However, we
were able to minimize the number of concern types
that contained both automatically and manually
detected concerns to just 9 concern types out of 116
(8%). When concerns in CRAM are considered, 15
“detailed changes” out of 45 (33%) contain both au-
tomatically and manually detected concerns. Thus,
we have been able to improve the separation of
automatically and manually detected concerns fur-
ther.
The existing classiﬁcation studies report a strik-
ing 75:25 ratio of evolvability to functional defects
or changes in industrial and OSS projects [24, 42].
Here, evolvability defects are the defects that af-
fect future development eﬀorts instead of runtime
behavior. Functional defects are the defects that
aﬀect runtime behavior. We also classiﬁed our con-
cern types as evolvability concerns and functional
concerns [69] and found a similar ratio of 74:26.
Thus, our study also supports their implication that
code review is superior to software testing as it not
only ﬁnds the same amount of functional defects as
testing but also identiﬁes a large number of evolv-
ability (non-functional) defects. However, it should
be noted that in contrast to these other studies, we
did not categorize only true positives. Rather, we
considered all reported concerns in the sample that
represent the entire code review eﬀort involved.
5.5. Threats to Validity
5.5.1. Construct Validity
Researcher bias: When a single person (pri-
mary author) is categorizing a large number of con-
cerns, it is diﬃcult to completely eliminate the re-
searcher bias. To minimize the eﬀect of researcher
bias, all categories were deﬁned and the deﬁnitions
were discussed among the authors. During the cat-
egorization process, 10 concerns each in 5 rounds
were categorized by the other two authors and dis-
cussed among the three authors. The labels were
updated and moved during these sessions to better
represent the concern types. Once the categoriza-
tion was completed, 15 randomly selected concerns
from the classiﬁcation were categorized by the other
two authors followed by a discussion. While creat-
ing the detection method categorization, in addi-
tion to deﬁning the categories, 5 well-known defect
detection tools were studied to back up the concerns
that we categorized as “automatically detected”.
The concerns that existing tools could not detect
but obviously could be accurately and automati-
cally detected were thoroughly discussed among the
authors. To minimize the bias during data analysis,
we present the data to backup our discussion and
conclusions.
Also, once the primary author had
completed the data analysis, the resulting impli-
cations and conclusions were discussed among the
authors.
Sampling bias: Sampling bias is a possibility in
any study that works with a sample from a popula-
tion. To minimize this bias, we performed random
sampling at the code review discussion level. Fur-
thermore, once the classiﬁcation was completed, an-
other 5% concerns of the sample size were extracted
and categorized for saturation check which demon-
strated that the saturation may have reached.
16

The sample sizes we used to check the saturation
and the reliability of our classiﬁcation have a pos-
sibility of being insuﬃcient for a thorough check.
However, we discovered the 75:25 evolvability to
functional concerns ratio that have been observed in
the previous code review defect classiﬁcation stud-
ies [27, 42, 24]. This provides us some conﬁdence
that the data we categorized is saturated and reli-
able.
5.5.2. External Validity
The data set selected for this study had been ex-
tracted from a popular code review tool Gerrit Code
Review and consists of 7 OSS Java projects. Thus,
the generalizability of our results towards closed-
source projects, OSS projects that are not consid-
ered in this study, and projects created in other pro-
gramming languages might be limited.
However,
all of our categories overlap with the categories of
other existing classiﬁcations (see section 5.4) and
obtained the ratio of 75:25 evolvability to functional
defects similar to other studies [27, 42, 24]. Due
to these reasons, our results may be applicable to
many other scenarios. Future research can validate
this further.
Our sample size was 385 with 55 or more code re-
view comments extracted from each project. This
sample size maybe too small to identify rare and
highly project-speciﬁc concerns. However, we found
that 16% of our sample concerns required custom
rules to be detected. These may be highly project-
speciﬁc concerns and representative of such con-
cerns in the population. Future work can explore
project-speciﬁc concerns more by classifying larger
samples of code review comments.
6. Conclusion
We have presented a classiﬁcation of concerns
identiﬁed in MCRs of 7 OSS Java projects [28]. We
extracted a sample of 417 code review comments,
and, using thematic analysis [29], we identiﬁed 116
concern types which we grouped into 15 groups.
Additionally, we categorized the concerns based on
the automatability of their detection.
The ﬁrst RQ of this study was to identify
the types of concerns that were detected during
MCRs. We identiﬁed 116 concern types and 15 con-
cern groups. The API Documentation group had
the largest number of concerns (68/417 concerns,
16.3%). The Implementation group had the most
number of concern types (19/116 concern types,
16.4%). The least number of concerns were related
to Threads (1/417, 0.2%). The entire list of concern
types and concern groups are available in ﬁg.5.
The second RQ was aimed at identifying the ex-
tent to which the concerns identiﬁed during MCR
could have been automatically detected accurately
using existing tools. The results suggested that this
was a substantial amount, 22% of concerns in the
sample.
Additionally, 16% of concerns were not
supported by existing, popular tools but automat-
ing the detection of these concerns using custom
rules was possible. This is another opportunity for
further reducing the cost of code review.
Not all concerns require an expert to detect
them (see section 5.3.4). Using less expensive non-
programmers to conduct code review where possi-
ble could also improve the cost of code review. In
our sample, 42.5% concerns were categorized as de-
tectable by non-programmers with or without train-
ing.
There are concerns that can be partially auto-
matically detected using existing tools. This is an-
other cost reduction opportunity where the appli-
cation of defect detection tools could help with to
a certain extent.
As future work, we will explore the last two cost
reduction opportunities discussed above.
We ex-
pect to answer the research question “Can non-
programmers contribute to code review?”. The ver-
iﬁed classiﬁcation may allow organizations to dis-
patch code review tasks to appropriate expertise-
levels and possibly use less expensive resources to
conduct code review tasks. Additionally, the extent
to which partial automation capabilities of existing
tools can help with code review tasks will be ex-
plored so that the use of automation and appropri-
ate expertise-level can work together to reduce the
code review cost eﬀectively.
References
[1] A. Bosu, M. Greiler, C. Bird, Characteristics of useful
code reviews: An empirical study at microsoft, in: in
Proc. IEEE/ACM 12th Working Conference on Min-
ing Software Repositories, 2015, pp. 146–156.
[2] C. Sadowski, E. S¨oderberg, L. Church, M. Sipko,
A. Bacchelli, Modern code review:
a case study at
google, in: in Proc. 40th International Conference on
Software Engineering: Software Engineering in Prac-
tice, 2018, pp. 181–190.
[3] D. Distefano,
M. F¨ahndrich,
F. Logozzo,
P. W.
O’Hearn, Scaling static analyses at facebook, Com-
munications of the ACM 62 (8) (2019) 62–70.
17

[4] S. Nazir, N. Fatima, S. Chuprat, Modern code review
beneﬁts-primary ﬁndings of a systematic literature re-
view, in:
in Proc. 3rd International Conference on
Software Engineering and Information Management,
2020, p. 210–215.
[5] A. Bosu, J. C. Carver, C. Bird, J. Orbeck, C. Chockley,
Process aspects and social dynamics of contemporary
code review: Insights from open source development
and industrial practice at microsoft, IEEE Transac-
tions on Software Engineering 43 (1) (2017) 56–75.
[6] T. Baum, O. Liskin, K. Niklas, K. Schneider, Fac-
tors inﬂuencing code review processes in industry, in:
in Proc.24th ACM SIGSOFT International Sympo-
sium on Foundations of Software Engineering, 2016,
p. 85–96.
[7] J. Wang, P. C. Shih, J. M. Carroll, Revisiting Li-
nus’s law: Beneﬁts and challenges of open source soft-
ware peer review, International Journal of Human-
Computer Studies 77 (2015) 52–65.
[8] A. Bacchelli, C. Bird, Expectations, outcomes, and
challenges of modern code review, in: in Proc. 35th In-
ternational Conference on Software Engineering, 2013,
pp. 712–721.
[9] J. Czerwonka, M. Greiler, J. Tilford, Code reviews
do not ﬁnd bugs. how the current code review best
practice slows us down, in: in Proc. 37th IEEE Inter-
national Conference on Software Engineering, Vol. 2,
2015, pp. 27–28.
[10] T. Baum, O. Liskin, K. Niklas, K. Schneider, Fac-
tors inﬂuencing code review processes in industry, in:
in Proc. 24th acm sigsoft international symposium on
foundations of software engineering, 2016, pp. 85–96.
[11] C. Staﬀ, Codeﬂow: Improving the code review pro-
cess at microsoft, Communications of the ACM 62 (2)
(2019) 36–44.
[12] L. Pascarella, D. Spadini, F. Palomba, M. Bruntink,
A. Bacchelli, Information needs in contemporary code
review, Proceedings of the ACM on Human-Computer
Interaction 2 (CSCW) (2018) 1–27.
[13] Gerrit code review.
URL https://www.gerritcodereview.com/
[14] Collaborator.
URL
https://smartbear.com/product/
collaborator/overview/
[15] Crucible.
URL
https://www.atlassian.com/software/
crucible
[16] Github code review.
URL https://github.com/features/code-review/
[17] V. Balachandran, Reducing human eﬀort and improv-
ing quality in peer code reviews using automatic static
analysis and reviewer recommendation, in: in Proc.
35th International Conference on Software Engineer-
ing, 2013, pp. 931–940.
[18] S. Panichella, V. Arnaoudova, M. Di Penta, G. Anto-
niol, Would static analysis tools help developers with
code reviews?, in: in Proc. 22nd International Confer-
ence on Software Analysis, Evolution, and Reengineer-
ing, 2015, pp. 161–170.
[19] M. Fadhel, Towards automating code reviews (2020).
[20] D. Singh, V. R. Sekar, K. T. Stolee, B. Johnson, Eval-
uating how static analysis tools can reduce code re-
view eﬀort, in: in Proc. IEEE Symposium on Visual
Languages and Human-Centric Computing, 2017, pp.
101–105.
[21] B. Johnson, Y. Song, E. Murphy-Hill, R. Bowdidge,
Why don’t software developers use static analysis tools
to ﬁnd bugs?, in: in Proc. of the International Confer-
ence on Software Engineering, 2013, p. 672–681.
[22] C. Vassallo, S. Panichella, F. Palomba, S. Proksch,
H. C. Gall, A. Zaidman, How developers engage with
static analysis tools in diﬀerent contexts, Empirical
Software Engineering (2019) 1–39.
[23] H. Siy, L. Votta, Does the modern code inspection
have value?, in: in Proc. International Conference on
Software Maintenance, 2001, pp. 281–289.
[24] M. V. M¨antyl¨a, C. Lassenius, What types of defects are
really discovered in code reviews?, IEEE Transactions
on Software Engineering 35 (3) (2009) 430–448.
[25] Q. Lei, Z. He, H. Fuqun, L. Bin, Classiﬁcation of
air on-board software code defects and investigations,
Procedia Engineering 15 (2011) 3577–3583.
[26] P. Runeson, A. Steﬁk, A. Andrews, Variation factors in
the design and analysis of replicated controlled exper-
iments, Empirical Software Engineering 19 (6) (2014)
1781–1808.
[27] S. Panichella, N. Zaugg, An empirical investigation
of relevant changes and automation needs in modern
code review, Empirical Software Engineering (2020)
1–40.
[28] M. Paixao, J. Krinke, D. Han, M. Harman, Crop:
Linking code reviews to source code changes, in: Inter-
national Conference on Mining Software Repositories,
MSR, 2018.
[29] V. Braun, V. Clarke, Using thematic analysis in
psychology, Qualitative research in psychology 3 (2)
(2006) 77–101.
[30] Sonarqube rules.
URL
https://docs.sonarqube.org/latest/user-
guide/rules/
[31] Checkstyle standard checks.
URL
https://checkstyle.sourceforge.io/checks.
html
[32] Findbugs - bug descriptions.
URL
http://findbugs.sourceforge.net/
bugDescriptions.html
[33] Pmd.
URL https://pmd.github.io/
[34] Intellij idea code inspection.
URL
https://www.jetbrains.com/help/idea/code-
inspection.html
[35] Checkstyle standard checks - nowhitespaceafter.
URL
https://checkstyle.sourceforge.io/config_
whitespace.html\#NoWhitespaceAfter
[36] N. Cassee, B. Vasilescu, A. Serebrenik, The silent
helper: the impact of continuous integration on code
reviews, in: in Proc. 27th IEEE International Confer-
ence on Software Analysis, Evolution and Reengineer-
ing, 2020, pp. 423–434.
[37] Checkstyle - type name.
URL
https://checkstyle.sourceforge.io/config_
naming.html\#TypeName
[38] Intellij idea - identiﬁer naming.
URL
https://www.jetbrains.com/help/clion/
naming-conventions.html
[39] J. Lacomis,
P. Yin,
E. Schwartz,
M. Allamanis,
C. Le Goues, G. Neubig, B. Vasilescu, Dire: A neu-
ral approach to decompiled identiﬁer naming, in: in
Proc. 34th IEEE/ACM International Conference on
Automated Software Engineering, 2019, pp. 628–639.
18

[40] B. Lin, S. Scalabrino, A. Mocci, R. Oliveto, G. Bavota,
M. Lanza, Investigating the use of code analysis and
nlp to promote a consistent usage of identiﬁers, in:
in Proc. IEEE 17th International Working Conference
on Source Code Analysis and Manipulation, 2017, pp.
81–90.
[41] M. Allamanis, E. T. Barr, C. Bird, C. Sutton, Learning
natural coding conventions, in: in Proc. 22nd ACM
SIGSOFT International Symposium on Foundations
of Software Engineering, 2014, pp. 281–293.
[42] M. Beller, A. Bacchelli, A. Zaidman, E. Juergens,
Modern code reviews in open-source projects: Which
problems do they ﬁx?, in: in Proc. 11th working con-
ference on mining software repositories, 2014, pp. 202–
211.
[43] J. Agnelo, N. Laranjeiro, J. Bernardino, Using orthog-
onal defect classiﬁcation to characterize nosql database
defects, Journal of Systems and Software 159 (2020)
110451.
[44] X. Xia, X. Zhou, D. Lo, X. Zhao, Y. Wang, An em-
pirical study of bugs in software build system, IEICE
TRANSACTIONS on Information and Systems 97 (7)
(2014) 1769–1780.
[45] U. Hunny, Orthogonal security defect classiﬁcation for
secure software development, Ph.D. thesis (2012).
[46] F. Thung, S. Wang, D. Lo, L. Jiang, An empirical
study of bugs in machine learning systems, in: in Proc.
23rd International Symposium on Software Reliability
Engineering, 2012, pp. 271–280.
[47] N. Li, Z. Li, X. Sun, Classiﬁcation of software defect
detected by black-box testing: An empirical study, in:
2010 Second World Congress on Software Engineering,
Vol. 2, 2010, pp. 234–240.
[48] M. Grottke, A. P. Nikora, K. S. Trivedi, An empiri-
cal investigation of fault types in space mission sys-
tem software, in: in Proc. International Conference
on Dependable Systems & Networks (DSN), 2010, pp.
447–456.
[49] Ieee standard classiﬁcation for software anomalies,
IEEE Std 1044-2009 (Revision of IEEE Std 1044-1993)
(2010) 1–23.
[50] R. S. Pressman, Software engineering: a practitioner’s
approach, Palgrave macmillan, 2005.
[51] R. C. Michael Inies, Fault links: identifying module
and fault types and their relationship, Master’s thesis,
University of Kentucky (2004).
[52] B. Beizer, Software testing techniques, Dreamtech
Press, 2003.
[53] X. Huang, Software reliability, safety and quality as-
surance, Beijing: Publishing House of Electronics In-
dustry 112 (2002).
[54] C. Kaner, J. Falk, H. Q. Nguyen, Testing computer
software, John Wiley & Sons, 1999.
[55] W. S. Humphrey, A discipline for software engineering,
Addison-Wesley Longman Publishing Co., Inc., 1995.
[56] R. B. Grady, Practical software metrics for project
management and process improvement, Prentice-Hall,
Inc., 1992.
[57] L. H. Putnam, W. Myers, Measures for excellence: re-
liable software on time, within budget, Prentice Hall
Professional Technical Reference, 1991.
[58] V. R. Basili, R. W. Selby, Comparing the eﬀectiveness
of software testing strategies, IEEE transactions on
software engineering (12) (1987) 1278–1296.
[59] F. Lopes, J. Agnelo, C. A. Teixeira, N. Laranjeiro,
J. Bernardino, Automating orthogonal defect classiﬁ-
cation using machine learning algorithms, Future Gen-
eration Computer Systems 102 (2020) 932–947.
[60] J. Mabrey, Automated defect classiﬁcation using ma-
chine learning, Ph.D. thesis, North Carolina Agricul-
tural and Technical State University (2020).
[61] I. Chawla, S. K. Singh, An automated approach for
bug categorization using fuzzy logic, in: in Proc. of
the 8th India Software Engineering Conference, 2015,
pp. 90–99.
[62] N. Pingclasai, H. Hata, K.-i. Matsumoto, Classifying
bug reports to bugs and other requests using topic
modeling, in: in Proc. 20th Asia-Paciﬁc Software En-
gineering Conference (APSEC), Vol. 2, 2013, pp. 13–
18.
[63] C. Liu, Y. Zhao, Y. Yang, H. Lu, Y. Zhou, B. Xu,
An ast-based approach to classifying defects, in: in
Proc.International Conference on Software Quality,
Reliability and Security-Companion, 2015, pp. 14–21.
[64] X. Xia, D. Lo, X. Wang, B. Zhou, Automatic defect
categorization based on fault triggering conditions, in:
in Proc. 19th International Conference on Engineering
of Complex Computer Systems, 2014, pp. 39–48.
[65] L. Yu, C. Kong, L. Xu, J. Zhao, H. Zhang, Mining bug
classiﬁer and debug strategy association rules for web-
based applications, in:
in Proc. International Con-
ference on Advanced Data Mining and Applications,
2008, pp. 427–434.
[66] Code review open platform (crop).
URL https://crop-repo.github.io/\#structure
[67] M. Paixao, P. H. Maia, Rebasing in code review con-
sidered harmful:
A large-scale empirical investiga-
tion, in: 2019 19th International Working Conference
on Source Code Analysis and Manipulation (SCAM),
2019, pp. 45–55.
[68] R. V. Krejcie, D. W. Morgan, Determining sample size
for research activities, Educational and psychological
measurement 30 (3) (1970) 607–610.
[69] S. Gunawardena, Supporting documentation - code
review cost reduction opportunities.
URL
https://github.com/sgun571/Code-Review-
Cost-Reduction-Opportunities
[70] grammarly.
URL https://www.grammarly.com/
[71] Checkstyle - needbraces.
URL
https://checkstyle.sourceforge.io/config_
blocks.html\#NeedBraces
[72] Pmd - controlstatementbraces.
URL
https://pmd.github.io/latest/pmd_rules_
java_codestyle.html\#controlstatementbraces
[73] Intellij idea - control ﬂow statement without braces.
URL
https://www.jetbrains.com/help/idea/list-
of-java-inspections.html
[74] A. M. Jamil, L. b. Othmane, A. Valani, M. Ab-
delkhalek, A. Tek, The current practices of chang-
ing secure software: an empirical study, in: in Proc.
35th Annual ACM Symposium on Applied Comput-
ing, 2020, pp. 1566–1575.
[75] M. Tahaei, K. Vaniea, A survey on developer-centred
security, in: in Proc. IEEE European Symposium on
Security and Privacy Workshops, 2019, pp. 129–138.
[76] T. Thomas, Exploring the usability and eﬀectiveness
of interactive annotation and code review for the de-
tection of security vulnerabilities, in: in Proc. IEEE
Symposium on Visual Languages and Human-Centric
19

Computing, 2015, pp. 295–296.
[77] B. Chess, J. West, Secure programming with static
analysis, 2007.
[78] V. B. Livshits, M. S. Lam, Finding security vulner-
abilities in java applications with static analysis., in:
USENIX Security Symposium, Vol. 14, 2005, pp. 18–
18.
[79] R. K. McLean, Comparing static security analysis
tools using open source software, in: in Proc. Sixth
International Conference on Software Security and Re-
liability Companion, 2012, pp. 68–74.
[80] A. Masood, J. Java, Static analysis for web service
security-tools & techniques for a secure development
life cycle, in: 2015 IEEE International Symposium on
Technologies for Homeland Security (HST), 2015, pp.
1–6.
[81] J. Smith, L. N. Q. Do, E. Murphy-Hill, Why can’t
johnny ﬁx vulnerabilities:
A usability evaluation of
static analysis tools for security, in: Sixteenth Sympo-
sium on Usable Privacy and Security ({SOUPS} 2020),
2020, pp. 221–238.
[82] L. N. Q. Do, J. Wright, K. Ali, Why do software devel-
opers use static analysis tools? a user-centered study
of developer needs and motivations, IEEE Transac-
tions on Software Engineering (Jun 2020).
[83] P. Godefroid, N. Nagappan, Concurrency at microsoft:
An exploratory survey, in: CAV workshop on exploit-
ing concurrency eﬃciently and correctly, 2008.
[84] H. Chen, S. Guo, Y. Xue, Y. Sui, C. Zhang, Y. Li,
H. Wang, Y. Liu, {MUZZ}: Thread-aware grey-box
fuzzing for eﬀective bug hunting in multithreaded
programs, in: 29th {USENIX} Security Symposium
({USENIX} Security 20), 2020, pp. 2325–2342.
[85] F. Eichinger, V. Pankratius, P. W. Große, K. B¨ohm,
Localizing defects in multithreaded programs by min-
ing dynamic call graphs, in: International Academic
and Industrial Conference on Practice and Research
Techniques, 2010, pp. 56–71.
[86] O. Edelstein, E. Farchi, E. Goldin, Y. Nir, G. Ratsaby,
S. Ur, Framework for testing multi-threaded java pro-
grams, Concurrency and Computation: Practice and
Experience 15 (3-5) (2003) 485–499.
[87] M. A. Al Mamun, A. Khanam, H. Grahn, R. Feldt,
Comparing four static analysis tools for java concur-
rency bugs, in: Third Swedish Workshop on Multi-
Core Computing (MCC-10), 2010, pp. 18–19.
[88] H. Liu, G. Li, J. F. Lukman, J. Li, S. Lu, H. S. Gunawi,
C. Tian, Dcatch: Automatically detecting distributed
concurrency bugs in cloud systems, ACM SIGARCH
Computer Architecture News 45 (1) (2017) 677–691.
[89] D. Kester, M. Mwebesa, J. S. Bradbury, How good is
static analysis at ﬁnding concurrency bugs?, in: 2010
10th IEEE Working Conference on Source Code Anal-
ysis and Manipulation, 2010, pp. 115–124.
[90] Code conventions for the java tm programming lan-
guage - if.
URL
https://www.oracle.com/java/technologies/
javase/codeconventions-statements.html#449
[91] R. S. Laramee,
Bob’s concise coding conventions
(c3), Advances in Computer Science and Engineering
(ACSE) 4 (1) (2010) 23–26.
[92] Coding conventions.
URL
https://en.wikipedia.org/wiki/Coding_
conventions
[93] Checkstyle standard checks - linelength.
URL
https://checkstyle.sourceforge.io/config_
sizes.html\#LineLength
[94] Intellij idea - line is longer than allowed by code style.
URL
https://www.jetbrains.com/help/phpstorm/
general-line-is-longer-than-allowed-by-code-
style.html
[95] Quality proﬁles - sonarqube.
URL
https://docs.sonarqube.org/latest/
instance-administration/quality-profiles/
[96] B. Johnson, Y. Song, E. Murphy-Hill, R. Bowdidge,
Why don’t software developers use static analysis tools
to ﬁnd bugs?, in: in Proc. 35th International Confer-
ence on Software Engineering (ICSE), 2013, pp. 672–
681.
[97] D. Singh, V. R. Sekar, K. T. Stolee, B. Johnson, Eval-
uating how static analysis tools can reduce code review
eﬀort, in: in Proc. IEEE Symposium on Visual Lan-
guages and Human-Centric Computing (VL/HCC),
2017, pp. 101–105.
[98] Findbugs - return value of method without side eﬀect
is ignored.
URL
http://findbugs.sourceforge.net/
bugDescriptions.html#RV_RETURN_VALUE_IGNORED_
NO_SIDE_EFFECT
[99] T. Gilb, D. Graham, Software inspections, Addison-
Wesley Reading, Masachusetts, 1993.
[100] D. O’Neill, Issues in software inspection, IEEE Soft-
ware 14 (1) (1997) 18–19.
[101] T. Hall, D. Wilson, N. Baddoo, Towards implementing
successful software inspections, in: Proceedings Inter-
national Conference on Software Methods and Tools.
SMT 2000, 2000, pp. 127–136.
[102] M. Ciolkowski, O. Laitenberger, S. Biﬄ, Software re-
views, the state of the practice, IEEE software 20 (6)
(2003) 46–51.
[103] Y.-K. Wong, An exploratory study of software re-
view in practice, in:
PICMET’03:
Portland Inter-
national Conference on Management of Engineering
and Technology Technology Management for Reshap-
ing the World, 2003., 2003, pp. 301–308.
[104] J.-S. Oh, H.-J. Choi, A reﬂective practice of auto-
mated and manual code reviews for a studio project,
in: Fourth Annual ACIS International Conference on
Computer and Information Science (ICIS’05), 2005,
pp. 37–42.
[105] L. Harjumaa, I. Tervonen, A. Huttunen, Peer re-
views in real life-motivators and demotivators, in:
Fifth International Conference on Quality Software
(QSIC’05), 2005, pp. 29–36.
[106] P. Sliz, A. Morin, Optimizing peer review of software
code, Science 341 (6143) (2013) 236–237.
[107] S. Jayatilake, S. De Silva, U. Settinayake, S. Yapa,
J. Jayamanne, A. Ruwanthika, C. Manawadu, Role of
software inspections in the sri lankan software develop-
ment industry, in: 2013 8th International Conference
on Computer Science & Education, 2013, pp. 697–702.
[108] G.
Gousios,
A.
Zaidman,
M.-A.
Storey,
A.
Van
Deursen,
Work
practices
and
challenges
in pull-based development: The integrator’s perspec-
tive, in: 2015 IEEE/ACM 37th IEEE International
Conference on Software Engineering, Vol. 1, 2015, pp.
358–368.
[109] O. Kononenko, O. Baysal, M. W. Godfrey, Code re-
view quality: How developers see it, in: Proceedings
of the 38th international conference on software engi-
20

neering, 2016, pp. 1028–1038.
[110] L. MacLeod, M. Greiler, M.-A. Storey, C. Bird, J. Cz-
erwonka, Code reviewing in the trenches: Challenges
and best practices, IEEE Software 35 (4) (2017) 34–42.
[111] T. Baum, H. Leßmann, K. Schneider, The choice of
code review process: A survey on the state of the prac-
tice, in: International Conference on Product-Focused
Software Process Improvement, 2017, pp. 111–127.
[112] F. Ebert, F. Castor, N. Novielli, A. Serebrenik, An
exploratory study on confusion in code reviews, Em-
pirical Software Engineering 26 (1) (2021) 1–48.
[113] P. M. Johnson, An instrumented approach to improv-
ing software quality through formal technical review,
in: Proceedings of 16th International Conference on
Software Engineering, 1994, pp. 113–122.
[114] F. Belli, R. Crisan, Towards automation of checklist-
based code-reviews, in: Proceedings of ISSRE’96: 7th
International Symposium on Software Reliability En-
gineering, 1996, pp. 24–33.
[115] A. A. Porter, H. P. Siy, C. A. Toman, L. G. Votta, An
experiment to assess the cost-beneﬁts of code inspec-
tions in large scale software development, IEEE trans-
actions on software engineering 23 (6) (1997) 329–346.
[116] F. Belli, R. Crisan, Empirical performance analysis
of computer-supported code-reviews, in: Proceedings
The Eighth International Symposium on Software Re-
liability Engineering, 1997, pp. 245–255.
[117] K. Chan, An agent-based approach to computer as-
sisted code inspections, in:
Proceedings 2001 Aus-
tralian Software Engineering Conference, 2001, pp.
147–152.
[118] D. Kelly, T. Shepard, Task-directed software inspec-
tion, Journal of Systems and Software 73 (2) (2004)
361–368.
[119] E. Farchi, S. Ur, Selective homeworkless reviews, in:
2008 1st International Conference on Software Testing,
Veriﬁcation, and Validation, 2008, pp. 404–413.
[120] J. Ratcliﬀe, Moving software quality upstream: The
positive impact of lightweight peer code review, in:
Paciﬁc NW software quality conference, 2009, pp. 1–
10.
[121] B. Xu, Cost eﬃcient software review in an e-business
software development project, in: 2010 International
Conference on E-Business and E-Government, 2010,
pp. 2680–2683.
[122] V. Balachandran, Reducing human eﬀort and improv-
ing quality in peer code reviews using automatic static
analysis and reviewer recommendation, in: in Proc.
35th International Conference on Software Engineer-
ing, 2013, pp. 931–940.
[123] S. Misra, L. Fern´andez, R. Colomo-Palacios, A simpli-
ﬁed model for software inspection, Journal of software:
evolution and process 26 (12) (2014) 1297–1315.
[124] C. Staﬀ, Codeﬂow: improving the code review pro-
cess at microsoft, Communications of the ACM 62 (2)
(2019) 36–44.
[125] S. Rebai, A. Amich, S. Molaei, M. Kessentini, R. Kaz-
man, Multi-objective code reviewer recommendations:
balancing expertise, availability and collaborations,
Automated Software Engineering (2020) 1–28.
[126] Z. Xia, H. Sun, J. Jiang, X. Wang, X. Liu, A hybrid
approach to code reviewer recommendation with col-
laborative ﬁltering, in: 2017 6th International Work-
shop on Software Mining (SoftwareMining), 2017, pp.
24–31.
[127] S. McIntosh, Y. Kamei, B. Adams, A. E. Hassan, An
empirical study of the impact of modern code review
practices on software quality, Empirical Software En-
gineering 21 (5) (2016) 2146–2189.
[128] M. M. Rahman, C. K. Roy, J. A. Collins, Correct: code
reviewer recommendation in github based on cross-
project and technology experience, in: Proceedings of
the 38th International Conference on Software Engi-
neering Companion, 2016, pp. 222–231.
[129] Y. Wang, X. Wang, Y. Jiang, Y. Liang, Y. Liu, A code
reviewer assignment model incorporating the compe-
tence diﬀerences and participant preferences, Foun-
dations of Computing and Decision Sciences 41 (1)
(2016) 77–91.
[130] P. Thongtanunam, C. Tantithamthavorn, R. G. Kula,
N. Yoshida, H. Iida, K.-i. Matsumoto, Who should re-
view my code? a ﬁle location-based code-reviewer rec-
ommendation approach for modern code review, in:
2015 IEEE 22nd International Conference on Software
Analysis, Evolution, and Reengineering (SANER),
2015, pp. 141–150.
21
