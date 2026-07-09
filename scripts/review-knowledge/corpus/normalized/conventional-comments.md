Comments like this are unhelpful…

By simply prefixing the comment with a label, the intention is clear and the tone dramatically changes.

Labels also prompt the reviewer to give more **actionable** comments.

**suggestion:** This is not worded correctly.

Can we change this to match the wording of the marketing page?

Labeling comments encourages collaboration and saves **hours** of undercommunication and misunderstandings. They are also parseable by machines!

**suggestion:** Let’s avoid using this specific function…

If we reference much of a function marked “Deprecated”, it is almost certain to disagree with us, sooner or later.

Conventional Comments is a standard for formatting comments of any kind of review/feedback process, such as:

Adhering to a consistent format improves reader’s expectations and machine readability. Here’s the format we propose:

```
<label> [decorations]: <subject>
[discussion]
```
For example:

**question (non-blocking):** At this point, does it matter which thread has won?

Maybe to prevent a race condition we should keep looping until they’ve all won?

Can be automatically parsed into:

```
{
 "label": "question",
 "subject": "At this point, does it matter which thread has won?",
 "decorations": ["non-blocking"],
 "discussion": "Maybe to prevent a race condition we should keep looping until they've all won?"
}
```
We strongly suggest using the following labels:

| praise: | Praises highlight something positive. Try to leave at least one of these comments per review. Do notleave false praise (which can actually be damaging).Dolook for something to sincerely praise. |
| nitpick: | Nitpicks are trivial preference-based requests. These should be non-blocking by nature. |
| suggestion: | Suggestions propose improvements to the current subject. It’s important to be explicit and clear on whatis being suggested andwhyit is an improvement. Consider using patches and theblockingornon-blockingdecorations to further communicate your intent. |
| issue: | Issues highlight specific problems with the subject under review. These problems can be user-facing or behind the scenes. It is strongly recommended to pair this comment with a `suggestion`. If you are not sure if a problem exists or not, consider leaving a`question`. |
| todo: | TODO’s are small, trivial, but necessary changes. Distinguishing todo comments from issues: or suggestions: helps direct the reader’s attention to comments requiring more involvement. |
| question: | Questions are appropriate if you have a potential concern but are not quite sure if it’s relevant or not. Asking the author for clarification or investigation can lead to a quick resolution. |
| thought: | Thoughts represent an idea that popped up from reviewing. These comments are non-blocking by nature, but they are extremely valuable and can lead to more focused initiatives and mentoring opportunities. |
| chore: | Chores are simple tasks that must be done before the subject can be “officially” accepted. Usually, these comments reference some common process. Try to leave a link to the process description so that the reader knows how to resolve the chore. |
| note: | Notes are always non-blocking and simply highlight something the reader should take note of. |

If you like to be a bit more expressive with your labels, you may also consider:

| typo: | Typo comments are like todo:, where the main issue is a misspelling. |
| polish: | Polish comments are like a suggestion, where there is nothing necessarily wrong with the relevant content, there’s just some ways to immediately improve the quality. |
| quibble: | Quibbles are very much like nitpick:, except it does not conjure up images of lice and animal hygiene practices. |

Feel free to diverge from this specific list of labels if it seems appropriate.

Decorations give additional context for a comment. They help further classify comments which have the same label (for example, a security suggestion as opposed to a test suggestion).

**suggestion (security):** I’m a bit concerned that we are implementing our own DOM purifying function here…

Could we consider using the framework instead?

Decorations may be specific to each organization. If needed, we recommend establishing a minimal set of decorations (leaving room for discretion) with no ambiguity.

Possible decorations include:

| (non-blocking) | A comment with this decoration should notprevent the subject under review from being accepted. This is helpful for organizations that consider comments blocking by default. |
| (blocking) | A comment with this decoration shouldprevent the subject under review from being accepted, until it is resolved. This is helpful for organizations that consider comments non-blocking by default. |
| (if-minor) | This decoration gives some freedom to the author that they should resolve the comment only if the changes end up being minor or trivial. |

Adding a decoration to a comment should improve understandability and maintain readability. Having a list of many decorations in one comment would conflict with this goal.

**nitpick:** `little star` => `little bat`

Can we update the other references as well?

**chore:** Let’s run the `jabber-walk` CI job to make sure this doesn’t break any known references.

Here are the docs for running this job. Feel free to reach out if you need any help!

If you have a suggestion, comment, idea, or something else altogether, please visit the GitLab project to collaborate. Issues and Merge Requests are welcome!

Check out this page for a showcase of projects and tools built by the wider community.

Check out this page for some general ideas and principles on improving written review communication.

The characters used in the examples are respectfully adapted from Lewis Carroll’s Alice in Wonderland, illustrated by John Tenniel.
