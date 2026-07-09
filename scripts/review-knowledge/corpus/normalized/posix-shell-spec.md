This chapter contains the definition of the Shell Command Language.

The shell is a command language interpreter. This chapter describes the syntax of that command language as it is used by the
*sh* utility and the *system*() and *popen*() functions defined in the System Interfaces volume of POSIX.1-2024.

The shell operates according to the following general overview of operations. The specific details are included in the cited sections of this chapter.

The shell reads its input from a file (see *sh*), from the **-c** option or from the
*system*() and *popen*() functions defined
in the System Interfaces volume of POSIX.1-2024. If the first line of a file of shell commands starts with the characters
`"#!"`, the results are unspecified.

The shell breaks the input into tokens: words and operators; see 2.3 Token Recognition.

The shell parses the input into simple commands (see 2.9.1 Simple Commands) and compound commands (see 2.9.4 Compound Commands).

For each word within a command, the shell processes <backslash>-escape sequences inside dollar-single-quotes (see 2.2.4 Dollar-Single-Quotes) and then performs various word expansions (see 2.6 Word Expansions). In the case of a simple command, the results usually include a list of pathnames and fields to be treated as a command name and arguments; see 2.9 Shell Commands.

The shell performs redirection (see 2.7 Redirection) and removes redirection operators and their operands from the parameter list.

The shell executes a function (see 2.9.5 Function Definition Command), built-in (see 2.15 Special Built-In Utilities), executable file, or script, giving the names of the arguments as positional
parameters numbered 1 to *n*, and the name of the command (or in the case of a function within a script, the name of the
script) as special parameter 0 (see 2.9.1.4 Command Search and Execution).

The shell optionally waits for the command to complete and collects the exit status (see 2.8.2 Exit Status for Commands).

Quoting is used to remove the special meaning of certain characters or words to the shell. Quoting can be used to preserve the literal meaning of the special characters in the next paragraph, prevent reserved words from being recognized as such, and prevent parameter expansion and command substitution within here-document processing (see 2.7.4 Here-Document ).

The application shall quote the following characters if they are to represent themselves:

| & ; < > ( ) $ ` \ " ' <space> <tab> <newline>

and the following might need to be quoted under certain circumstances. That is, these characters are sometimes special depending on conditions described elsewhere in this volume of POSIX.1-2024:

* ? [ ] ^ - ! # ~ = % { , }

The various quoting mechanisms are the escape character, single-quotes, double-quotes, and dollar-single-quotes. The here-document represents another form of quoting; see 2.7.4 Here-Document.

A <backslash> that is not quoted shall preserve the literal value of the following character, with the exception of a <newline>. If a <newline> immediately follows the <backslash>, the shell shall interpret this as line continuation. The <backslash> and <newline> shall be removed before splitting the input into tokens. Since the escaped <newline> is removed entirely from the input and is not replaced by any white space, it cannot serve as a token separator.

Enclosing characters in single-quotes (`''`) shall preserve the literal value of each character within the single-quotes.
A single-quote cannot occur within single-quotes.

Enclosing characters in double-quotes (`""`) shall preserve the literal value of all characters within the double-quotes,
with the exception of the characters backquote, <dollar-sign>, and <backslash>, as follows:

The input characters within the quoted string that are also enclosed between `"$("` and the matching `')'` shall
not be affected by the double-quotes, but rather shall define the command(s) whose output replaces the `"$(...)"` when the
word is expanded. The tokenizing rules in 2.3 Token Recognition shall be applied recursively to find the
matching `')'`.

For the four varieties of parameter expansion that provide for substring processing (see 2.6.2 Parameter
Expansion), within the string of characters from an enclosed `"${"` to the matching `'}'`, the double-quotes
within which the expansion occurs shall have no effect on the handling of any special characters.

For parameter expansions other than the four varieties that provide for substring processing, within the string of characters
from an enclosed `"${"` to the matching `'}'`, the double-quotes within which the expansion occurs shall preserve the
literal value of all characters, with the exception of the characters double-quote, backquote, <dollar-sign>, and
<backslash>. If any unescaped double-quote characters occur within the string, other than in embedded command substitutions,
the behavior is unspecified. The backquote and <dollar-sign> characters shall follow the same rules as for characters in
double-quotes described in this section. The <backslash> character shall follow the same rules as for characters in
double-quotes described in this section except that it shall additionally retain its special meaning as an escape character when
followed by `'}'` and this shall prevent the escaped `'}'` from being considered when determining the matching
`'}'` (using the rule in 2.6.2 Parameter Expansion).

A quoted (single-quoted, double-quoted, or dollar-single-quoted) string that begins, but does not end, within the
`"`...`"` sequence

A `"`...`"` sequence that begins, but does not end, within the same double-quoted string

$ ` \ <newline>

or by a double-quote character that would otherwise be considered special (see 2.6.4 Arithmetic Expansion and 2.7.4 Here-Document).

When double-quotes are used to quote a parameter expansion, command substitution, or arithmetic expansion, the literal value of all characters within the result of the expansion shall be preserved.

The application shall ensure that a double-quote that is not within `"$(...)"` nor within `"${...}"` is
immediately preceded by a <backslash> in order to be included within double-quotes. The parameter `'@'` has special
meaning inside double-quotes and is described in 2.5.2 Special Parameters.

A sequence of characters starting with a <dollar-sign> immediately followed by a single-quote (`$'`) shall preserve
the literal value of all characters up to an unescaped terminating single-quote (`'`), with the exception of certain
<backslash>-escape sequences, as follows:

`\"` yields a <quotation-mark> (double-quote) character, but note that <quotation-mark> can be included
unescaped.

`\'` yields an <apostrophe> (single-quote) character.

`\\` yields a <backslash> character.

`\a` yields an <alert> character.

`\b` yields a <backspace> character.

`\e` yields an <ESC> character.

`\f` yields a <form-feed> character.

`\n` yields a <newline> character.

`\r` yields a <carriage-return> character.

`\t` yields a <tab> character.

`\v` yields a <vertical-tab> character.

`\c`*X* yields the control character listed in the **Value** column of *Values for cpio c_mode Field* in the OPERANDS section of the *stty* utility when *X* is one of the characters listed in the **^c** column of the same
table, except that `\c\\` yields the <FS> control character since the <backslash> character has to be
escaped.

`\x`*XX* yields the byte whose value is the hexadecimal value *XX* (one or more hexadecimal digits). If more
than two hexadecimal digits follow `\x`, the results are unspecified.

`\`*ddd* yields the byte whose value is the octal value *ddd* (one to three octal digits).

The behavior of an unescaped <backslash> immediately followed by any other character, including <newline>, is unspecified.

In cases where a variable number of characters can be used to specify an escape sequence (`\x`*XX* and
`\`*ddd*), the escape sequence shall be terminated by the first character that is not of the expected type or, for
`\`*ddd* sequences, when the maximum number of characters specified has been found, whichever occurs first.

These <backslash>-escape sequences shall be processed (replaced with the bytes or characters they yield) immediately prior to word expansion (see 2.6 Word Expansions) of the word in which the dollar-single-quotes sequence occurs.

If a `\x`*XX* or `\`*ddd* escape sequence yields a byte whose value is 0, it is unspecified whether that
null byte is included in the result or if that byte and any following regular characters and escape sequences up to the terminating
unescaped single-quote are evaluated and discarded.

If the octal value specified by `\`*ddd* will not fit in a byte, the results are unspecified.

If a `\e` or `\c`*X* escape sequence specifies a character that does not have an encoding in the locale in
effect when these <backslash>-escape sequences are processed, the result is implementation-defined. However, implementations
shall not replace an unsupported character with bytes that do not form valid characters in that locale's character set.

If a <backslash>-escape sequence represents a single-quote character (for example `\'`), that sequence shall not
terminate the dollar-single-quote sequence.

The shell shall read its input in terms of lines. (For details about how the shell reads its input, see the description of
*sh*.) The input lines can be of unlimited length. These lines shall be parsed using two
major modes: ordinary token recognition and processing of here-documents.

When an **io_here** token has been recognized by the grammar (see 2.10 Shell Grammar), one or more
of the subsequent lines immediately following the next **NEWLINE** token form the body of a here-document and shall be parsed
according to the rules of 2.7.4 Here-Document. Any non-**NEWLINE** tokens (including more
**io_here** tokens) that are recognized while searching for the next **NEWLINE** token shall be saved for processing after
the here-document has been parsed. If a saved token is an **io_here** token, the corresponding here-document shall start on the
line immediately following the line containing the trailing delimiter of the previous here-document. If any saved token includes a
<newline> character, the behavior is unspecified.

When it is not processing an **io_here**, the shell shall break its input into tokens by applying the first applicable rule
below to each character in turn in its input. At the start of input or after a previous token has just been delimited, the first or
next token, respectively, shall start with the first character that has not already been included in a token and is not discarded
according to the rules below. Once a token has started, zero or more characters from the input shall be appended to the token until
the end of the token is delimited according to one of the rules below. When both the start and end of a token have been delimited,
the characters forming the token shall be exactly those in the input between the two delimiters, including any quoting characters.
If a rule below indicates that a token is delimited, and no characters have been included in the token, that empty token shall be
discarded.

If the end of input is recognized, the current token (if any) shall be delimited.

If the previous character was used as part of an operator and the current character is not quoted and can be used with the previous characters to form an operator, it shall be used as part of that (operator) token.

If the previous character was used as part of an operator and the current character cannot be used with the previous characters to form an operator, the operator containing the previous character shall be delimited.

If the current character is an unquoted <backslash>, single-quote, or double-quote or is the first character of an unquoted <dollar-sign> single-quote sequence, it shall affect quoting for subsequent characters up to the end of the quoted text. The rules for quoting are as described in 2.2 Quoting. During token recognition no substitutions shall be actually performed, and the result token shall contain exactly the characters that appear in the input unmodified, including any embedded or enclosing quotes or substitution operators, between the start and the end of the quoted text. The token shall not be delimited by the end of the quoted field.

If the current character is an unquoted `'$'` or `'`'`, the shell shall identify the start of any candidates for
parameter expansion ( 2.6.2 Parameter Expansion), command substitution ( 2.6.3
Command Substitution), or arithmetic expansion ( 2.6.4 Arithmetic Expansion) from their
introductory unquoted character sequences: `'$'` or `"${"`, `"$("` or `'`'`, and `"$(("`,
respectively. The shell shall read sufficient input to determine the end of the unit to be expanded (as explained in the cited
sections). While processing the characters, if instances of expansions or quoting are found nested within the substitution, the
shell shall recursively process them in the manner specified for the construct that is found. For `"$("` and `'`'`
only, if instances of **io_here** tokens are found nested within the substitution, they shall be parsed according to the rules
of 2.7.4 Here-Document; if the terminating `')'` or `'`'` of the substitution occurs
before the **NEWLINE** token marking the start of the here-document, the behavior is unspecified. The characters found from the
beginning of the substitution to its end, allowing for any recursion necessary to recognize embedded constructs, shall be included
unmodified in the result token, including any embedded or enclosing substitution operators or quotes. The token shall not be
delimited by the end of the substitution.

If the current character is not quoted and can be used as the first character of a new operator, the current token (if any) shall be delimited. The current character shall be used as the beginning of the next (operator) token.

If the current character is an unquoted <blank>, any token containing the previous character is delimited and the current character shall be discarded.

If the previous character was part of a word, the current character shall be appended to that word.

If the current character is a `'#'`, it and all subsequent characters up to, but excluding, the next <newline>
shall be discarded as a comment. The <newline> that ends the line is not considered part of the comment.

The current character is used as the start of a new word.

Once a token is delimited, it is categorized as required by the grammar in 2.10 Shell Grammar.

In situations where the shell parses its input as a *program*, once a *complete_command* has been recognized by the
grammar (see 2.10 Shell Grammar), the *complete_command* shall be executed before the next
*complete_command* is tokenized and parsed.

After a token has been categorized as type **TOKEN** (see 2.10.1 Shell Grammar Lexical
Conventions), including (recursively) any token resulting from an alias substitution, the **TOKEN** shall be subject to
alias substitution if all of the following conditions are true:

The **TOKEN** does not contain any quoting characters.

The **TOKEN** is a valid alias name (see XBD *3.10 Alias Name*).

An alias with that name is in effect.

The **TOKEN** did not either fully or, optionally, partially result from an alias substitution of the same alias name at any
earlier recursion level.

Either the **TOKEN** is being considered for alias substitution because it follows an alias substitution whose replacement
value ended with a <blank> (see below) or the **TOKEN** could be parsed as the command name word of a simple command (see
2.10 Shell Grammar), based on this **TOKEN** and the tokens (if any) that preceded it, but ignoring
whether any subsequent characters would allow that.

When a **TOKEN** is subject to alias substitution, the value of the alias shall be processed as if it had been read from the
input instead of the **TOKEN**, with token recognition (see 2.3 Token Recognition) resuming at the
start of the alias value. When the end of the alias value is reached, the shell may behave as if an additional <space>
character had been read from the input after the **TOKEN** that was replaced. If it does not add this <space>, it is
unspecified whether the current token is delimited before token recognition is applied to the character (if any) that followed the
**TOKEN** in the input.

An implementation may defer the effect of a change to an alias but the change shall take effect no later than the completion of
the currently executing *complete_command* (see 2.10 Shell Grammar). Changes to aliases shall not
take effect out of order. Implementations may provide predefined aliases that are in effect when the shell is invoked.

When used as specified by this volume of POSIX.1-2024, alias definitions shall not be inherited by separate invocations of the shell or by the utility execution environments invoked by the shell; see 2.13 Shell Execution Environment .

Reserved words are words that have special meaning to the shell; see 2.9 Shell Commands. The following words shall be recognized as reserved words:

|
 |
 |
 |
 |

This recognition shall only occur when none of the characters is quoted and when the word is used as:

See the grammar in 2.10 Shell Grammar.

When used in circumstances where reserved words are recognized (described above), the following words may be
recognized as reserved words, in which case the results are unspecified except as described below for **time**:

|
 |
 |
 |
 |
 |
 |

When the word **time** is recognized as a reserved word in circumstances where it would, if it were not a
reserved word, be the command name (see 2.9.1.1 Order of Processing) of a simple command that would
execute the *time* utility in a manner other than one for which *time* states that the results are unspecified, the behavior shall be as specified for
the *time* utility.

When used in circumstances where reserved words are recognized (described above), all words whose final character
is a <colon> (`':'`) are reserved; their use in those circumstances produces unspecified results.

A parameter can be denoted by a name, a number, or one of the special characters listed in 2.5.2 Special Parameters. A variable is a parameter denoted by a name.

A parameter is set if it has an assigned value (null is a valid value). Once a variable is set, it can only be
unset by using the *unset* special built-in command.

Parameters can contain arbitrary byte sequences, except for the null byte. The shell shall process their values as characters only when performing operations that are described in this standard in terms of characters.

A positional parameter is a parameter denoted by a decimal representation of a positive integer. The digits denoting the positional parameters shall always be interpreted as a decimal value, even if there is a leading zero. When a positional parameter with more than one digit is specified, the application shall enclose the digits in braces (see 2.6.2 Parameter Expansion).

Examples:

Positional parameters are initially assigned when the shell is invoked (see *sh*), temporarily replaced when a shell function is invoked (see 2.9.5
Function Definition Command), and can be reassigned with the *set* special built-in command.

Listed below are the special parameters and the values to which they shall expand. Only the values of the special parameters are listed; see 2.6 Word Expansions for a detailed summary of all the stages involved in expanding words.

If one of these conditions is true, the initial fields shall be retained as separate fields, except that if the
parameter being expanded was embedded within a word, the first field shall be joined with the beginning part of the original word
and the last field shall be joined with the end part of the original word. In all other contexts the results of the expansion are
unspecified. If there are no positional parameters, the expansion of `'@'` shall generate zero fields, even when
`'@'` is within double-quotes; however, if the expansion is embedded within a word which contains one or more other parts
that expand to a quoted null string, these null string(s) shall still produce an empty field, except that if the other parts are
all within the same double-quotes as the `'@'`, it is unspecified whether the result is zero fields or one empty field.

See the description of the *IFS* variable in 2.5.3 Shell Variables.

Variables shall be initialized from the environment (as defined by XBD *8. Environment Variables* and the *exec* function in the System Interfaces
volume of POSIX.1-2024) and can be given new values with variable assignment commands. Shell variables shall be initialized only
from environment variables that have valid names. If a variable is initialized from the environment, it shall be marked for export
immediately; see the *export* special built-in. New variables can be defined and initialized with
variable assignments, with the *read* or *getopts* utilities, with the *name* parameter in a **for** loop, with the
${*name*=*word*} expansion, or with other mechanisms provided as implementation extensions.

The following variables shall affect the execution of the shell:

This variable, when and only when an interactive shell is invoked, shall be subjected to parameter expansion (see
2.6.2 Parameter Expansion) by the shell and the resulting value shall be used as a pathname of a file.
Before any interactive commands are read, the shell shall tokenize (see 2.3 Token Recognition) the
contents of the file, parse the tokens as a *program* (see 2.10 Shell Grammar), and execute the
resulting commands in the current environment. (In other words, the contents of the *ENV* file are not parsed as a single
*compound_list*. This distinction matters because it influences when aliases take effect.) The file need not be executable. If
the expanded value of *ENV* is not an absolute pathname, the results are unspecified. *ENV* shall be ignored if the
user's real and effective user IDs or real and effective group IDs are different.

If *IFS* is not set, it shall behave as normal for an unset variable, except that field splitting by the shell
and line splitting by the *read* utility shall be performed as if the value of *IFS*
is <space><tab><newline>; see 2.6.5 Field Splitting.

The shell shall set *IFS* to <space><tab><newline> when it is invoked.

Set by the shell to a decimal number representing the current sequential line number (numbered starting with 1)
within a script or function before it executes each command. If the user unsets or resets *LINENO ,* the variable may lose its
special meaning for the life of the shell. If the shell is not currently executing a script or function, the value of *LINENO*
is unspecified.

Each time an interactive shell is ready to read a command, the value of this variable shall be subjected to parameter expansion (see 2.6.2 Parameter Expansion) and exclamation-mark expansion (see below). Whether the value is also subjected to command substitution (see 2.6.3 Command Substitution) or arithmetic expansion (see 2.6.4 Arithmetic Expansion) or both is unspecified. After expansion, the value shall be written to standard error.

The expansions shall be performed in two passes, where the result of the first pass is input to the second pass. One of the passes shall perform only the exclamation-mark expansion described below. The other pass shall perform the other expansion(s) according to the rules in 2.6 Word Expansions. Which of the two passes is performed first is unspecified.

The default value shall be `"$ "`. For users who have specific additional implementation-defined
privileges, the default may be another, implementation-defined value.

Exclamation-mark expansion: The shell shall replace each instance of the <exclamation-mark> character
(`'!'`) with the history file number (see *Command History List*)
of the next command to be typed. An <exclamation-mark> character escaped by another <exclamation-mark> character (that
is, `"!!"`) shall expand to a single <exclamation-mark> character.

Each time the user enters a <newline> prior to completing a command line in an interactive shell, the value
of this variable shall be subjected to parameter expansion (see 2.6.2 Parameter Expansion). Whether
the value is also subjected to command substitution (see 2.6.3 Command Substitution) or arithmetic
expansion (see 2.6.4 Arithmetic Expansion) or both is unspecified. After expansion, the value shall be
written to standard error. The default value shall be `"> "`.

When an execution trace (*set* **-x**) is being performed, before each line in the
execution trace, the value of this variable shall be subjected to parameter expansion (see 2.6.2 Parameter
Expansion). Whether the value is also subjected to command substitution (see 2.6.3 Command
Substitution) or arithmetic expansion (see 2.6.4 Arithmetic Expansion) or both is unspecified.
After expansion, the value shall be written to standard error. The default value shall be `"+ "`.

This section describes the various expansions that are performed on words. Not all expansions are performed on every word, as explained in the following sections and elsewhere in this chapter. The expansions that are performed for a given word shall be performed in the following order:

Tilde expansions, parameter expansions, command substitutions, arithmetic expansions, and quote removals that occur within a single word shall expand to a single field, except as described below. The shell shall create multiple fields or no fields from a single word only as a result of field splitting, pathname expansion, or the following cases:

may be subject to an additional implementation-defined form of expansion that can create multiple fields from a single word. This expansion, if supported, shall be applied before all the other word expansions are applied. The other expansions shall then be applied to each field that results from this expansion.

When the expansions in this section are performed other than in the context of preparing a command for execution, they shall be carried out in the current shell execution environment.

When expanding words for a command about to be executed, and the word will be the command name or an argument to the command, the expansions shall be carried out in the current shell execution environment. (The environment for the command to be executed is unknown until the command word is known.)

When expanding the words in a command about to be executed that are used with variable assignments or redirections, it is unspecified whether the expansions are carried out in the current execution environment or in the environment of the command about to be executed.

The `'$'` character is used to introduce parameter expansion, command substitution, or arithmetic
evaluation. If a `'$'` that is neither within single-quotes nor escaped by a <backslash> is immediately followed by a
character that is not a <space>, not a <tab>, not a <newline>, and is not one of the following:

the result is unspecified. If a `'$'` that is neither within single-quotes nor escaped by a
<backslash> is immediately followed by a <space>, <tab>, or a <newline>, or is not followed by any
character, the `'$'` shall be treated as a literal character.

A "tilde-prefix" consists of an unquoted <tilde> character at the beginning of a word, followed by all of
the characters preceding the first unquoted <slash> in the word, or all the characters in the word if there is no
<slash>. In an assignment (see XBD *4.26 Variable Assignment*),
multiple tilde-prefixes can be used: one at the beginning of the word (that is, following the <equals-sign> of the
assignment), or one following any unquoted <colon>, or both. A tilde-prefix in an assignment is terminated by the first
unquoted <colon> or <slash>, or the end of the assignment word.

If the tilde-prefix consists of only the <tilde> character, it shall be replaced by the value of the variable
*HOME .* If *HOME* is unset, the results are unspecified.

Otherwise, the characters in the tilde-prefix following the <tilde> shall be treated as a possible login name
from the user database. If these characters do not form a portable login name (see the description of the *LOGNAME*
environment variable in XBD *8.3 Other Environment Variables*), the
results are unspecified.

~"string" ~'string' ~$var ~\/bin

owing to the presence of `'"'`, `'\''`, `'$'`, `'\\'`, and `'/'` characters in
the login name.

If the characters in the tilde-prefix following the <tilde> form a portable login name, the tilde-prefix
shall be replaced by a pathname of the initial working directory associated with the login name. The pathname shall be obtained as
if by using the *getpwnam*() function as defined in the System Interfaces volume of
POSIX.1-2024. If the system does not recognize the login name, the results are unspecified.

The pathname that replaces the tilde-prefix shall be treated as if quoted to prevent it being altered by field
splitting and pathname expansion; if a <slash> follows the tilde-prefix and the pathname ends with a <slash>, the
trailing <slash> from the pathname should be omitted from the replacement. If the word being expanded consists of only the
<tilde> character and *HOME* is set to the null string, this produces an empty field (as opposed to zero fields) as the
expanded word.

The format for parameter expansion is as follows:

${expression}

where *expression* consists of all characters until the matching `'}'`. Any `'}'` escaped by a
<backslash> or within a quoted string, and characters in embedded arithmetic expansions, command substitutions, and variable
expansions, shall not be examined in determining the matching `'}'`.

The simplest form for parameter expansion is:

${parameter}

The value, if any, of *parameter* shall be substituted.

The parameter name or symbol can be enclosed in braces, which are optional except for positional parameters with
more than one digit or when *parameter* is a name and is followed by a character that could be interpreted as part of the
name.

For a parameter that is not enclosed in braces:

In addition, a parameter expansion can be modified by using one of the following formats. In each case that a value
of *word* is needed (based on the state of *parameter*, as described below), *word* shall be subjected to tilde
expansion, parameter expansion, command substitution, arithmetic expansion, and quote removal. If *word* is not needed, it
shall not be expanded. The `'}'` character that delimits the following parameter expansion modifications shall be determined
as described previously in this section and in 2.2.3 Double-Quotes. If *parameter* is
`'*'` or `'@'`, the result of the expansion is unspecified.

In the parameter expansions shown previously, use of the <colon> in the format shall result in a test for a
parameter that is unset or null; omission of the <colon> shall result in a test for a parameter that is only unset. If
parameter is `'#'` and the colon is omitted, the application shall ensure that *word* is specified (this is necessary
to avoid ambiguity with the string length expansion). The following table summarizes the effect of the <colon>:

|
 |
 |
 |
 |
|---|---|---|---|
|
 | substitute | substitute | substitute |
|
 | substitute | substitute null | substitute |
|
 | substitute | assign | assign |
|
 | substitute | substitute null | assign |
|
 | substitute | error, exit | error, exit |
|
 | substitute | substitute null | error, exit |
|
 | substitute | substitute null | substitute null |
|
 | substitute | substitute | substitute null |

In all cases shown with "substitute", the expression is replaced with the value shown. In all cases shown with
"assign", *parameter* is assigned that value, which also replaces the expression.

The following four varieties of parameter expansion provide for character substring processing. In each case,
pattern matching notation (see 2.14 Pattern Matching Notation), rather than regular expression notation,
shall be used to evaluate the patterns. If *parameter* is `'#'`, `'*'`, or `'@'`, the result of the
expansion is unspecified. If *parameter* is unset and *set* **-u** is in effect, the expansion
shall fail. Enclosing the full parameter expansion string in double-quotes shall not cause the following four varieties of pattern
characters to be quoted, whereas quoting characters within the braces shall have this effect. In each variety, if *word* is
omitted, the empty pattern shall be used.

a=1 set 2 echo ${a}b-$ab-${1}0-${10}-$101b--20--20

foo=asdf echo ${foo-bar}xyz}asdfxyz}foo= echo ${foo-bar}xyz}xyz}unset foo echo ${foo-bar}xyz}barxyz}

${x:-$(ls)}

unset X echo ${X:=abc}abc

unset posix echo ${posix:?}sh: posix: parameter null or not set

set a b c echo ${3:+posix}posix

HOME=/usr/posix echo ${#HOME}10

x=file.c echo ${x%.c}.ofile.o

x=posix/src/std echo ${x%%/*}posix

x=$HOME/src/cmd echo ${x#$HOME}/src/cmd

x=/one/two/three echo ${x##*/}three

The double-quoting of patterns is different depending on where the double-quotes are placed:

Command substitution allows the output of one or more commands to be substituted in place of the commands themselves. Command substitution shall occur when command(s) are enclosed as follows:

$(commands)

or (backquoted version):

`commands`

The shell shall expand the command substitution by executing *commands* in a subshell environment (see
2.13 Shell Execution Environment) and replacing the command substitution (the text of the *commands*
string plus the enclosing `"$()"` or backquotes) with the standard output of the command(s); if the output ends with one or
more bytes that have the encoded value of a <newline> character, they shall not be included in the replacement. Any such
bytes that occur elsewhere shall be included in the replacement; however, they might be treated as field delimiters and eliminated
during field splitting, depending on the value of *IFS* and quoting that is in effect. If the output contains any null bytes,
the behavior is unspecified.

Within the backquoted style of command substitution, if the command substitution is not within double-quotes,
<backslash> shall retain its literal meaning, except when followed by: `'$'`, `'`'`, or <backslash>. See
2.2.3 Double-Quotes for the handling of <backslash> when the command substitution is within
double-quotes. The search for the matching backquote shall be satisfied by the first unquoted non-escaped backquote; during this
search, if a non-escaped backquote is encountered within a shell comment, a here-document, an embedded command substitution of the
$(*commands*) form, or a quoted string, undefined results occur. A quoted string that begins, but does not end, within the
`"`...`"` sequence produces undefined results.

With the $(*commands*) form, all characters following the open parenthesis to the matching closing parenthesis
constitute the *commands* string.

With both the backquoted and $(*commands*) forms, the *commands* string shall be tokenized (see 2.3 Token Recognition) and parsed (see 2.10 Shell Grammar). It is unspecified whether
the *commands* string is parsed and executed incrementally as a *program* (as for a shell script), or is parsed as a
single *compound_list* that is executed after the string has been completely parsed. In addition, it is unspecified whether
the terminating `')'` of the $(*commands*) form can result from alias substitution. With the $(*commands*) form
any syntactically correct *program* can be used for *commands*, except that:

The results of command substitution shall not be processed for further tilde expansion, parameter expansion, command substitution, or arithmetic expansion.

Command substitution can be nested. To specify nesting within the backquoted version, the application shall precede the inner backquotes with <backslash> characters; for example:

\`commands\`

The syntax of the shell command language has an ambiguity for expansions beginning with `"$(("`, which can
introduce an arithmetic expansion or a command substitution that starts with a subshell. Arithmetic expansion has precedence; that
is, the shell shall first determine whether it can parse the expansion as an arithmetic expansion and shall only parse the
expansion as a command substitution if it determines that it cannot parse the expansion as an arithmetic expansion. The shell need
not evaluate nested expansions when performing this determination. If it encounters the end of input without already having
determined that it cannot parse the expansion as an arithmetic expansion, the shell shall treat the expansion as an incomplete
arithmetic expansion and report a syntax error. A conforming application shall ensure that it separates the `"$("` and
`'('` into two tokens (that is, separate them with white space) in a command substitution that starts with a subshell. For
example, a command substitution containing a single subshell could be written as:

$( (commands) )

Arithmetic expansion provides a mechanism for evaluating an arithmetic expression and substituting its value. The format for arithmetic expansion shall be as follows:

$((expression))

The expression shall be treated as if it were in double-quotes, except that a double-quote inside the expression is not treated specially. The shell shall expand all tokens in the expression for parameter expansion, command substitution, and quote removal.

Next, the shell shall treat this as an arithmetic expression and substitute the value of the expression. The
arithmetic expression shall be processed according to the rules given in *1.1.2.1 Arithmetic Precision and Operations*, with the following
exceptions:

All changes to variables in an arithmetic expression shall be in effect after the arithmetic expansion, as in the
parameter expansion `"${x=value}"`.

If the shell variable *x* contains a value that forms a valid integer constant, optionally including a leading
<plus-sign> or <hyphen-minus>, then the arithmetic expansions `"$((x))"` and `"$(($x))"` shall return the
same value.

As an extension, the shell may recognize arithmetic expressions beyond those listed. The shell may use a signed
integer type with a rank larger than the rank of **signed long**. The shell may use a real-floating type instead of **signed
long** as long as it does not affect the results in cases where there is no overflow. If the expression is invalid, or the
contents of a shell variable used in the expression are not recognized by the shell, the expansion fails and the shell shall write
a diagnostic message to standard error indicating the failure.

A simple example using arithmetic expansion:

# repeat a command 100 times x=100 while [ $x -gt 0 ] docommandx=$(($x-1)) done

After parameter expansion ( 2.6.2 Parameter Expansion), command substitution (
2.6.3 Command Substitution), and arithmetic expansion ( 2.6.4 Arithmetic
Expansion), if the shell variable *IFS* (see 2.5.3 Shell Variables) is set and its value is
not empty, or if *IFS* is unset, the shell shall scan each field containing results of expansions and substitutions that did
not occur in double-quotes for field splitting; zero, one or multiple fields can result.

For the remainder of this section, any reference to the results of an expansion, or results of expansions, shall be interpreted to mean the results from one or more unquoted variable or arithmetic expansions, or unquoted command substitutions.

If the *IFS* variable is set and has an empty string as its value, no field splitting shall occur. However, if
an input field which contained the results of an expansion is entirely empty, it shall be removed. Note that this occurs before
quote removal; any input field that contains any quoting characters can never be empty at this point. After the removal of any such
fields from the input, the possibly modified input field list shall become the output.

Each input field shall be considered in sequence, first to last, with the results of the algorithm described in this section causing output fields to be generated, which shall remain in the same order as the input fields from which they originated.

Fields which contain no results from expansions shall not be affected by field splitting, and shall remain unaltered, simply moving from the list of input fields to be next in the list of output fields.

In the remainder of this description, it is assumed that there is present in the field at least one expansion result; this assumption will not be restated. Field splitting only ever alters those parts of the field.

For the purposes of this section, the term "*IFS* white space" is used to mean any of the white-space bytes
(see XBD *3.413 White Space*, *3.414 White-Space Byte*, and *3.415 White-Space Character*) <space>, <tab>, or <newline>
from the portable character set (see XBD *6.1 Portable Character Set*)
which are present in the value of the *IFS* variable, and perhaps other white-space characters. It is implementation-defined
whether other white-space characters which appear in the value of *IFS* are also considered as "*IFS* white space". The
three characters above specified as *IFS* white-space bytes are always *IFS* white space, when they occur in the value of
*IFS ,* regardless of whether they are white-space characters in any relevant locale. For other locale-specific white-space
characters allowed by the implementation it is unspecified whether the character is considered as *IFS* white space if it is
white space at the time it is assigned to the *IFS* variable, or if it is white space at the time field splitting occurs. (The
locale might have changed between those events.)

If the *IFS* variable is unset, then for the purposes of this section, but without altering the value of the
variable, its value shall be considered to contain the three single-byte characters <space>, <tab>, and <newline>
from the portable character set, all of which are *IFS* white-space characters.

The shell shall use the byte sequences that form the characters in the value of the *IFS* variable as
delimiters. Each of the characters <space>, <tab>, and <newline> which appears in the value of *IFS* shall
be a single-byte delimiter. The shell shall use these delimiters as field terminators to split the results of expansions, along
with other adjacent bytes, into separate fields, as described below. Note that these delimiters terminate a field; they do not, of
themselves, cause a new field to start—subsequent bytes that are not from the results of an expansion, or that do not form
*IFS* white-space characters are required for a new field to begin.

Note that the shell processes arbitrary bytes from the input fields; there is no requirement that those bytes form valid characters.

If the results of the algorithm are that no fields are delimited; that is, if the input field is wholly empty or
consists entirely of *IFS* white space, the result shall be zero fields (rather than an empty field).

For the purposes of this section, when a field is said to be delimited, then the candidate field, as generated below shall become an output field. When the algorithm transforms a candidate into an output field it shall be appended to the current list of output fields.

Each field containing the results from an expansion shall be processed in order, intermixed with fields not containing the results of expansions, processed as described above, as if by using the following algorithm, examining bytes in the input field, from beginning to end:

Append this byte (or sequence) to the candidate, and remove it from the input. Start the next iteration of the loop.

Append the first byte of the sequence to the candidate, and remove that byte from the input. Start the next iteration of the loop.

Remove that byte sequence from the input, consider the new leading input byte sequence, and repeat this step.

Remove that byte sequence from the input, but note it was observed.

At this point, if the candidate is not empty, or if a sequence of bytes representing an *IFS* character that
is not *IFS* white space was seen at step 4, then a field is said to have been delimited, and the candidate shall become an
output field.

The ordered list of output fields so produced, which might be empty, shall replace the list of input fields.

After field splitting, if *set* **-f** is not in effect, each field in the resulting
command line shall be expanded using the algorithm described in 2.14 Pattern Matching Notation, qualified
by the rules in 2.14.3 Patterns Used for Filename Expansion.

The quote character sequence <dollar-sign> single-quote and the single-character quote characters
(<backslash>, single-quote, and double-quote) that were present in the original word shall be removed unless they have
themselves been quoted. Note that the single-quote character that terminates a <dollar-sign> single-quote sequence is itself
a single-character quote character.

Redirection is used to open and close files for the current shell execution environment (see 2.13 Shell Execution Environment) or for any command. Redirection operators can be used with numbers representing
file descriptors (see XBD *3.141 File Descriptor*) as described
below.

The overall format used for redirection is:

[n]redir-op word

The number *n* is an optional one or more digit decimal number designating the file descriptor number; the
application shall ensure it is delimited from any preceding text and immediately precedes the redirection operator *redir-op*
(with no intervening <blank> characters allowed). If *n* is quoted, the number shall not be recognized as part of the
redirection expression. For example:

echo \2>a

writes the character 2 into file **a**. If any part of *redir-op* is quoted, no redirection expression is
recognized. For example:

echo 2\>a

writes the characters 2>*a* to standard output. The optional number, redirection operator, and *word*
shall not appear in the arguments provided to the command to be executed (if any).

The shell may support an additional format used for redirection:

{location}redir-op word

where *location* is non-empty and indicates a location where an integer value can be stored, such as the name
of a shell variable. If this format is supported its behavior is implementation-defined.

The largest file descriptor number supported in shell redirections is implementation-defined; however, all implementations shall support at least 0 to 9, inclusive, for use by the application.

If the redirection operator is `"<<"` or `"<<-"`, the word that follows the redirection
operator shall be subjected to quote removal; it is unspecified whether any of the other expansions occur. For the other
redirection operators, the word that follows the redirection operator shall be subjected to tilde expansion, parameter expansion,
command substitution, arithmetic expansion, and quote removal. Pathname expansion shall not be performed on the word by a
non-interactive shell; an interactive shell may perform it, but if the expansion would result in more than one word it is
unspecified whether the redirection proceeds without pathname expansion being performed or the redirection fails.

If more than one redirection operator is specified with a command, the order of evaluation is from beginning to end.

A failure to open or create a file shall cause a redirection to fail.

Input redirection shall cause the file whose name results from the expansion of *word* to be opened for
reading on the designated file descriptor, or standard input if the file descriptor is not specified.

The general format for redirecting input is:

[n]<word

where the optional *n* represents the file descriptor number. If the number is omitted, the redirection shall
refer to standard input (file descriptor 0).

The two general formats for redirecting output are:

[n]>word[n]>|word

where the optional *n* represents the file descriptor number. If the number is omitted, the redirection shall
refer to standard output (file descriptor 1).

Output redirection using the `'>'` format shall fail if the *noclobber* option is set (see the
description of *set* **-C**) and the file named by the expansion of *word* exists and is either a
regular file or a symbolic link that resolves to a regular file; it may also fail if the file is a symbolic link that does not
resolve to an existing file. The check for existence, file creation, and open operations shall be performed atomically as is done
by the *open*() function as defined in System Interfaces volume of POSIX.1-2024 when the
O_CREAT and O_EXCL flags are set, except that if the file exists and is a symbolic link, the open operation need not fail with
[EEXIST] unless the symbolic link resolves to an existing regular file. Performing these operations atomically ensures that the
creation of lock files and unique (often temporary) files is reliable, with important caveats detailed in *C.2.7.2 Redirecting Output*. The check for the type of the file need not be
performed atomically with the check for existence, file creation, and open operations. If not, there is a potential race condition
that may result in a misleading shell diagnostic message when redirection fails. See XRAT *C.2.7.2 Redirecting Output* for more details.

In all other cases (*noclobber* not set, redirection using `'>'` does not fail for the reasons
stated above, or redirection using the `">|"` format), output redirection shall cause the file whose name results from
the expansion of *word* to be opened for output on the designated file descriptor, or standard output if none is specified. If
the file does not exist, it shall be created as an empty file; otherwise, it shall be opened as if the *open*() function was called with the O_TRUNC flag set.

Appended output redirection shall cause the file whose name results from the expansion of word to be opened for
output on the designated file descriptor. The file shall be opened as if the *open*()
function as defined in the System Interfaces volume of POSIX.1-2024 was called with the O_APPEND flag set. If the file does not
exist, it shall be created.

The general format for appending redirected output is as follows:

[n]>>word

where the optional *n* represents the file descriptor number. If the number is omitted, the redirection refers
to standard output (file descriptor 1).

The redirection operators `"<<"` and `"<<-"` both allow redirection of subsequent lines
read by the shell to the input of a command. The redirected lines are known as a "here-document".

The here-document shall be treated as a single word that begins after the next **NEWLINE** token and continues
until there is a line containing only the delimiter and a <newline>, with no <blank> characters in between. Then the
next here-document starts, if there is one. For the purposes of locating this terminating line, the end of a *command_string*
operand (see *sh*) shall be treated as a <newline> character, and the end of the
*commands* string in `$( commands)` and

[n]<<word here-document delimiter

where the optional *n* represents the file descriptor number. If the number is omitted, the here-document
refers to standard input (file descriptor 0). It is unspecified whether the file descriptor is opened as a regular file or some
other type of file. Portable applications cannot rely on the file descriptor being seekable (see XSH *lseek*()).

If any part of *word* is quoted, not counting double-quotes outside a command substitution if the
here-document is inside one, the delimiter shall be formed by performing quote removal on *word*, and the here-document lines
shall not be expanded. Otherwise:

If the redirection operator is `"<<-"`, all leading <tab> characters shall be stripped from
input lines after <backslash><newline> line continuation (when it applies) has been performed, and from the line
containing the trailing delimiter. Stripping of leading <tab> characters shall occur as the here-document is read from the
shell input (and consequently does not affect any <tab> characters that result from expansions).

If more than one `"<<"` or `"<<-"` operator is specified on a line, the here-document
associated with the first operator shall be supplied first by the application and shall be read first by the shell.

When a here-document is read from a terminal device and the shell is interactive, it shall write the contents of
the variable *PS2*, processed as described in 2.5.3 Shell Variables, to standard error before
reading each line of input until the delimiter has been recognized.

An example of a here-document follows:

cat <<eof1; cat <<eof2 Hi, eof1 Helene. eof2

The redirection operator:

[n]<&word

shall duplicate one input file descriptor from another, or shall close one. If *word* evaluates to one or more
digits, the file descriptor denoted by *n*, or standard input if *n* is not specified, shall be made to be a copy of the
file descriptor denoted by *word*; if the digits in *word* do not represent an already open file descriptor, a
redirection error shall result (see 2.8.1 Consequences of Shell Errors); if the file descriptor
denoted by *word* represents an open file descriptor that is not open for input, a redirection error may result. If
*word* evaluates to `'-'`, file descriptor *n*, or standard input if *n* is not specified, shall be closed.
Attempts to close a file descriptor that is not open shall not constitute an error. If *word* evaluates to something else, the
behavior is unspecified.

The redirection operator:

[n]>&word

shall duplicate one output file descriptor from another, or shall close one. If *word* evaluates to one or
more digits, the file descriptor denoted by *n*, or standard output if *n* is not specified, shall be made to be a copy
of the file descriptor denoted by *word*; if the digits in *word* do not represent an already open file descriptor, a
redirection error shall result (see 2.8.1 Consequences of Shell Errors); if the file descriptor
denoted by *word* represents an open file descriptor that is not open for output, a redirection error may result. If
*word* evaluates to `'-'`, file descriptor *n*, or standard output if *n* is not specified, is closed.
Attempts to close a file descriptor that is not open shall not constitute an error. If *word* evaluates to something else, the
behavior is unspecified.

The redirection operator:

[n]<>word

shall cause the file whose name is the expansion of *word* to be opened for both reading and writing on the
file descriptor denoted by *n*, or standard input if *n* is not specified. If the file does not exist, it shall be
created.

Certain errors shall cause the shell to write a diagnostic message to standard error and exit as shown in the following table:

|
 |
 |
 |
 |
|---|---|---|---|
| Shell language syntax error | shall exit | shall not exit | yes |
| Special built-in utility error | shall exit | shall not exit | no |
| Other utility (not a special | shall not exit | shall not exit | no |
| Redirection error with | shall exit | shall not exit | yes |
| Redirection error with | shall not exit | shall not exit | yes |
| Redirection error with | shall not exit | shall not exit | yes |
| Redirection error with other | shall not exit | shall not exit | yes |
| Variable assignment error | shall exit | shall not exit | yes |
| Expansion error | shall exit | shall not exit | yes |
| Command not found | may exit | shall not exit | yes |
| Unrecoverable read error | shall exit | shall exit | yes |

An expansion error is one that occurs when the shell expansions defined in 2.6 Word
Expansions are carried out (for example, `"${x!y}"`, because `'!'` is not a valid operator); an implementation
may treat these as syntax errors if it is able to detect them during tokenization, rather than during expansion.

If any of the errors shown as "shall exit" or "may exit" occur in a subshell environment, the shell shall (respectively, may) exit from the subshell environment with a non-zero status and continue in the environment from which that subshell environment was invoked.

In all of the cases shown in the table where an interactive shell is required not to exit and a non-interactive shell is required to exit, an interactive shell shall not perform any further processing of the command in which the error occurred.

Each command has an exit status that can influence the behavior of other shell commands. The exit status of commands that are not utilities is documented in this section. The exit status of the standard utilities is documented in their respective sections.

The exit status of a command shall be determined as follows:

This section describes the basic structure of shell commands. The following command descriptions each describe a format of the command that is only used to aid the reader in recognizing the command type, and does not formally represent the syntax. In particular, the representations include spacing between tokens in some places where <blank>s would not be necessary (when one of the tokens is an operator). Each description discusses the semantics of the command; for a formal definition of the command language, consult 2.10 Shell Grammar.

A *command* is one of the following:

Unless otherwise stated, the exit status of a command shall be that of the last simple command executed by the command. There shall be no limit on the size of any shell command other than that imposed by the underlying system (memory constraints, {ARG_MAX}, and so on).

A "simple command" is a sequence of optional variable assignments and redirections, in any sequence, optionally followed by words and redirections.

When a given simple command is required to be executed (that is, when any conditional construct such as an AND-OR
list or a **case** statement has not bypassed the simple command), the following expansions, assignments, and redirections shall
all be performed from the beginning of the command text to the end:

In the preceding list, the order of steps 3 and 4 may be reversed if no command name results from step 2 or if the command name matches the name of a special built-in utility; see 2.15 Special Built-In Utilities.

When determining whether a command name is a declaration utility, an implementation may use only lexical analysis. It is unspecified whether assignment context will be used if the command name would only become recognized as a declaration utility after word expansions.

Variable assignments shall be performed as follows:

If any of the variable assignments attempt to assign a value to a variable for which the *readonly* attribute
is set in the current shell environment (regardless of whether the assignment is made in that environment), a variable assignment
error shall occur. See 2.8.1 Consequences of Shell Errors for the consequences of these errors.

If a simple command has no command name after word expansion (see 2.9.1.1 Order of Processing), any redirections shall be performed in a subshell environment; it is unspecified whether this subshell environment is the same one as that used for a command substitution within the command. (To affect the current execution environment, see the exec special built-in.) If any of the redirections performed in the current shell execution environment fail, the command shall immediately fail with an exit status greater than zero, and the shell shall write an error message indicating the failure. See 2.8.1 Consequences of Shell Errors for the consequences of these failures on interactive and non-interactive shells.

Additionally, if there is no command name but the command contains a command substitution, the command shall complete with the exit status of the command substitution whose exit status was the last to be obtained. Otherwise, the command shall complete with a zero exit status.

If a simple command has a command name and an optional list of arguments after word expansion (see 2.9.1.1 Order of Processing), the following actions shall be performed:

|
 |
 |
 |
 |
 |

Once a utility has been searched for and found (either as a result of this specific search or as part of an
unspecified shell start-up activity), an implementation may remember its location and need not search for the utility again unless
the *PATH* variable has been the subject of an assignment. If the remembered location fails for a subsequent invocation, the
shell shall repeat the search to find the new location for the utility, if any.

If the utility would be executed with file descriptor 0, 1, or 2 closed, implementations may execute the utility with the file descriptor open to an unspecified file. If a standard utility or a conforming application is executed with file descriptor 0 not open for reading or with file descriptor 1 or 2 not open for writing, the environment in which the utility or application is executed shall be deemed non-conforming, and consequently the utility or application might not behave as described in this standard.

If the execution is being made via the *exec* special built-in utility, the shell shall
not create a separate utility environment for this execution; the new process image shall replace the current shell execution
environment. If the current shell environment is a subshell environment, the new process image shall replace the subshell
environment and the shell shall continue in the environment from which that subshell environment was invoked.

In either case, execution of the utility in the specified environment shall be performed as follows:

If the *execl*() function fails due to an error equivalent to the
[ENOEXEC] error defined in the System Interfaces volume of POSIX.1-2024, the shell shall execute a command equivalent to having a
shell invoked with the pathname resulting from the search as its first operand, with any remaining arguments passed to the new
shell, except that the value of `"$0"` in the new shell may be set to the command name. The shell may apply a heuristic
check to determine if the file to be executed could be a script and may bypass this command execution if it determines that the
file cannot be a script. In this case, it shall write an error message, and the command shall fail with an exit status of 126.

It is unspecified whether environment variables that were passed to the shell when it was invoked, but were not
used to initialize shell variables (see 2.5.3 Shell Variables) because they had invalid names, are
included in the environment passed to *execl*() and (if *execl*() fails as described above) to the new shell.

If the *execl*() function fails due to an error equivalent to the
[ENOEXEC] error, the shell shall execute a command equivalent to having a shell invoked with the command name as its first operand,
with any remaining arguments passed to the new shell. The shell may apply a heuristic check to determine if the file to be executed
could be a script and may bypass this command execution if it determines that the file cannot be a script. In this case, it shall
write an error message, and the command shall fail with an exit status of 126.

It is unspecified whether environment variables that were passed to the shell when it was invoked, but were not
used to initialize shell variables (see 2.5.3 Shell Variables) because they had invalid names, are
included in the environment passed to *execl*() and (if *execl*() fails as described above) to the new shell.

A *pipeline* is a sequence of one or more commands separated by the control operator `'|'`. For each
command but the last, the shell shall connect the standard output of the command to the standard input of the next command as if by
creating a pipe and passing the write end of the pipe as the standard output of the command and the read end of the pipe as the
standard input of the next command.

The format for a pipeline is:

[!]command1[|command2...]

If the pipeline begins with the reserved word **!** and *command1* is a subshell command, the application
shall ensure that the **(** operator at the beginning of *command1* is separated from the **!** by one or more
<blank> characters. The behavior of the reserved word **!** immediately followed by the **(** operator is
unspecified.

The standard output of *command1* shall be connected to the standard input of *command2*. The standard
input, standard output, or both of a command shall be considered to be assigned by the pipeline before any redirection specified by
redirection operators that are part of the command (see 2.7 Redirection).

If the pipeline is not in the background (see 2.9.3.1 Asynchronous AND-OR Lists and 2.11 Job Control), the shell shall wait for the last command specified in the pipeline to complete, and may also wait for all commands to complete.

The exit status of a pipeline shall depend on whether or not the *pipefail* option (see set) is enabled and whether or not the pipeline begins with the **!** reserved word, as described in the
following table. The *pipefail* option determines which command in the pipeline the exit status is derived from; the **!**
reserved word causes the exit status to be the logical NOT of the exit status of that command. The shell shall use the
*pipefail* setting at the time it begins execution of the pipeline, not the setting at the time it sets the exit status of the
pipeline. (For example, in `command1 | set -o pipefail` the exit status of `command1` has no effect on the exit
status of the pipeline, even if the shell executes `set -o pipefail` in the current shell environment.)

|
 |
 |
 |
|---|---|---|
| no | no | The exit status of the last (rightmost) command specified in the pipeline. |
| no | yes | Zero, if the last (rightmost) command in the pipeline returned a non-zero exit status; otherwise, 1. |
| yes | no | Zero, if all commands in the pipeline returned an exit status of 0; otherwise, the exit status of the last (rightmost) command specified in the pipeline that returned a non-zero exit status. |
| yes | yes | Zero, if any command in the pipeline returned a non-zero exit status; otherwise, 1. |

An *AND-OR list* is a sequence of one or more pipelines separated by the operators `"&&"` and
`"||"`.

A *list* is a sequence of one or more AND-OR lists separated by the operators `';'` and
`'&'`.

The operators `"&&"` and `"||"` shall have equal precedence and shall be evaluated with left
associativity. For example, both of the following commands write solely **bar** to standard output:

false && echo foo || echo bar true || echo foo && echo bar

A `';'` separator or a `';'` or <newline> terminator shall cause the preceding AND-OR list to
be executed sequentially; an `'&'` separator or terminator shall cause asynchronous execution of the preceding AND-OR
list.

The term "compound-list" is derived from the grammar in 2.10 Shell Grammar; it is
equivalent to a sequence of *lists*, separated by <newline> characters, that can be preceded or followed by an arbitrary
number of <newline> characters.

The following is an example that illustrates <newline> characters in compound-lists:

while # a couple of <newline>s

# a list date && who || ls; cat file # a couple of <newline>s

# another list wc file > output & true

do # 2 lists ls cat file done

If an AND-OR list is terminated by the control operator <ampersand> (`'&'`), the shell shall
execute the AND-OR list asynchronously in a subshell environment. This subshell shall execute in the background; that is, the shell
shall not wait for the subshell to terminate before executing the next command (if any); if there are no further commands to
execute, the shell shall not wait for the subshell to terminate before exiting.

If job control is enabled (see set, **-m**), the AND-OR list shall become a
job-control background job and a job number shall be assigned to it. If job control is disabled, the AND-OR list may become a
non-job-control background job, in which case a job number shall be assigned to it; if no job number is assigned it shall become a
background command but not a background job.

A job-control background job can be controlled as described in 2.11 Job Control.

The process ID associated with the asynchronous AND-OR list shall become known in the current shell execution environment; see 2.13 Shell Execution Environment. This process ID shall remain known until any one of the following occurs (and, unless otherwise specified, may continue to remain known after it occurs).

The implementation need not retain more than the {CHILD_MAX} most recent entries in its list of known process IDs in the current shell execution environment.

If, and only if, job control is disabled, the standard input for the subshell in which an asynchronous AND-OR list
is executed shall initially be assigned to an open file description that behaves as if **/dev/null** had been opened for reading
only. This initial assignment shall be overridden by any explicit redirection of standard input within the AND-OR list.

If the shell is interactive and the asynchronous AND-OR list became a background job, the job number and the process ID associated with the job shall be written to standard error using the format:

"[%d] %d\n", <job-number>, <process-id>

If the shell is interactive and the asynchronous AND-OR list did not become a background job, the process ID associated with the asynchronous AND-OR list shall be written to standard error in an unspecified format.

The exit status of an asynchronous AND-OR list shall be zero.

The exit status of the subshell in which the AND-OR list is asynchronously executed can be obtained using the
*wait* utility.

AND-OR lists that are separated by a <semicolon> (`';'`) shall be executed sequentially. The format
for executing AND-OR lists sequentially shall be:

aolist1[;aolist2]...

Each AND-OR list shall be expanded and executed in the order specified.

If job control is enabled, the AND-OR lists shall form all or part of a foreground job that can be controlled as described in 2.11 Job Control.

The exit status of a sequential AND-OR list shall be the exit status of the last pipeline in the AND-OR list that is executed.

The control operator `"&&"` denotes an AND list. The format shall be:

command1[&&command2]...

First *command1* shall be executed. If its exit status is zero, *command2* shall be executed, and so on,
until a command has a non-zero exit status or there are no more commands left to execute. The commands are expanded only if they
are executed.

The exit status of an AND list shall be the exit status of the last command that is executed in the list.

The control operator `"||"` denotes an OR List. The format shall be:

command1[||command2]...

First, *command1* shall be executed. If its exit status is non-zero, *command2* shall be executed, and so
on, until a command has a zero exit status or there are no more commands left to execute.

The exit status of an OR list shall be the exit status of the last command that is executed in the list.

The shell has several programming constructs that are "compound commands", which provide control flow for commands. Each of these compound commands has a reserved word or control operator at the beginning, and a corresponding terminator reserved word or operator at the end. In addition, each can be followed by redirections on the same line as the terminator. Each redirection shall apply to all the commands within the compound command that do not explicitly override that redirection.

In the descriptions below, the exit status of some compound commands is stated in terms of the exit status of a
*compound-list*. The exit status of a *compound-list* shall be the value that the special parameter `'?'` (see
2.5.2 Special Parameters) would have immediately after execution of the *compound-list*.

The format for grouping commands is as follows:

If a character sequence beginning with `"(("` would be parsed by the shell as an arithmetic expansion if
preceded by a `'$'`, shells which implement an extension whereby `"(( expression))"` is evaluated as an
arithmetic expression may treat the

The exit status of a grouping command shall be the exit status of *compound-list*.

The **for** loop shall execute a sequence of commands for each member in a list of *items*. The **for**
loop requires that the reserved words **do** and **done** be used to delimit the sequence of commands.

The format for the **for** loop is as follows:

forname[in[word...]]docompound-listdone

First, the list of words following **in** shall be expanded to generate a list of items. Then, the variable
*name* shall be set to each item, in turn, and the *compound-list* executed each time. If no items result from the
expansion, the *compound-list* shall not be executed. Omitting:

inword...

shall be equivalent to:

in "$@"

If there is at least one item in the list of items, the exit status of a **for** command shall be the exit
status of the last *compound-list* executed. If there are no items, the exit status shall be zero.

The conditional construct **case** shall execute the *compound-list* corresponding to the first
*pattern* (see 2.14 Pattern Matching Notation), if any are present, that is matched by the string
resulting from the tilde expansion, parameter expansion, command substitution, arithmetic expansion, and quote removal of the given
word. The reserved word **in** shall denote the beginning of the patterns to be matched. Multiple patterns with the same
*compound-list* shall be delimited by the `'|'` symbol. The control operator `')'` terminates a list of patterns
corresponding to a given action. The terminated pattern list and the following *compound-list* is called a **case**
statement *clause*. Each **case** statement clause, with the possible exception of the last, shall be terminated with
either `";;"` or `";&"`. The **case** construct terminates with the reserved word **esac** (**case**
reversed).

The format for the **case** construct is as follows:

casewordin[[(]pattern[|pattern]... )compound-list terminator]...[[(]pattern[|pattern]... )compound-list]esac

Where *terminator* is either `";;"` or `";&"` and is optional for the last
*compound-list*.

In order from the beginning to the end of the **case** statement, each *pattern* that labels a
*compound-list* shall be subjected to tilde expansion, parameter expansion, command substitution, and arithmetic expansion,
and the result of these expansions shall be compared against the expansion of *word*, according to the rules described in
2.14 Pattern Matching Notation (which also describes the effect of quoting parts of the pattern). After
the first match, no more patterns in the **case** statement shall be expanded, and the *compound-list* of the matching
clause shall be executed. If the **case** statement clause is terminated by `";;"`, no further clauses shall be examined.
If the **case** statement clause is terminated by `";&"`, then the *compound-list* (if any) of each subsequent
clause shall be executed, in order, until either a clause terminated by `";;"` is reached and its *compound-list* (if
any) executed or there are no further clauses in the **case** statement. The order of expansion and comparison of multiple
*pattern*s that label a *compound-list* statement is unspecified.

The exit status of **case** shall be zero if no patterns are matched. Otherwise, the exit status shall be the
exit status of the *compound-list* of the last clause to be executed.

The **if** command shall execute a *compound-list* and use its exit status to determine whether to execute
another *compound-list*.

The format for the **if** construct is as follows:

ifcompound-listthencompound-list[elifcompound-listthencompound-list]...[elsecompound-list]fi

The **if** *compound-list* shall be executed; if its exit status is zero, the **then**
*compound-list* shall be executed and the command shall complete. Otherwise, each **elif** *compound-list* shall be
executed, in turn, and if its exit status is zero, the **then** *compound-list* shall be executed and the command shall
complete. Otherwise, the **else** *compound-list* shall be executed.

The exit status of the **if** command shall be the exit status of the **then** or **else**
*compound-list* that was executed, or zero, if none was executed.

The **while** loop shall continuously execute one *compound-list* as long as another *compound-list*
has a zero exit status.

The format of the **while** loop is as follows:

whilecompound-list-1docompound-list-2done

The *compound-list-1* shall be executed, and if it has a non-zero exit status, the **while** command shall
complete. Otherwise, the *compound-list-2* shall be executed, and the process shall repeat.

The exit status of the **while** loop shall be the exit status of the last *compound-list-2* executed, or
zero if none was executed.

while some_command; st=$?; false; do ...

The exit status of *compound-list-1* is available through the special parameter `'?'` during execution
of *compound-list-2*, but is known to be zero at that point anyway.

The **until** loop shall continuously execute one *compound-list* as long as another *compound-list*
has a non-zero exit status.

The format of the **until** loop is as follows:

untilcompound-list-1docompound-list-2done

The *compound-list-1* shall be executed, and if it has a zero exit status, the **until** command completes.
Otherwise, the *compound-list-2* shall be executed, and the process repeats.

The exit status of the **until** loop shall be the exit status of the last *compound-list-2* executed, or
zero if none was executed.

A function is a user-defined name that is used as a simple command to call a compound command with new positional parameters. A function is defined with a "function definition command".

The format of a function definition command is as follows:

fname( )compound-command[io-redirect...]

The function is named *fname*; the application shall ensure that it is a name (see XBD *3.216 Name*) and that it is not the name of a special built-in utility. An
implementation may allow other characters in a function name as an extension. The implementation shall maintain separate name
spaces for functions and variables.

The argument *compound-command* represents a compound command, as described in 2.9.4
Compound Commands.

When the function is declared, none of the expansions in 2.6 Word Expansions shall be
performed on the text in *compound-command* or *io-redirect*; all expansions shall be performed as normal each time the
function is called. Similarly, the optional *io-redirect* redirections and any variable assignments within
*compound-command* shall be performed during the execution of the function itself, not the function definition. See 2.8.1 Consequences of Shell Errors for the consequences of failures of these operations on interactive and
non-interactive shells.

When a function is executed, it shall have the syntax-error properties described for special built-in utilities in the first item in the enumerated list at the beginning of 2.15 Special Built-In Utilities.

The *compound-command* shall be executed whenever the function name is specified as the name of a simple
command (see 2.9.1.4 Command Search and Execution). The operands to the command temporarily shall
become the positional parameters during the execution of the *compound-command*; the special parameter `'#'` also shall
be changed to reflect the number of operands. The special parameter 0 shall be unchanged. When the function completes, the values
of the positional parameters and the special parameter `'#'` shall be restored to the values they had before the function
was executed. If the special built-in *return* (see return) is executed in the
*compound-command*, the function completes and execution shall resume with the next command after the function call.

The exit status of a function definition shall be zero if the function was declared successfully; otherwise, it shall be greater than zero. The exit status of a function invocation shall be the exit status of the last command executed by the function.

The following grammar defines the Shell Command Language. This formal syntax shall take precedence over the preceding text syntax description.

The input language to the shell shall be first recognized at the character level. The resulting tokens shall be classified by their immediate context according to the following rules (applied in order). These rules shall be used to determine what a "token" is that is subject to parsing at the token level. The rules for token recognition in 2.3 Token Recognition shall apply.

Further distinction on **TOKEN** is context-dependent. It may be that the same **TOKEN** yields **WORD**,
a **NAME**, an **ASSIGNMENT_WORD**, or one of the reserved words below, dependent upon the context. Some of the productions
in the grammar below are annotated with a rule number from the following list. When a **TOKEN** is seen where one of those
annotated productions could be used to reduce the symbol, the applicable rule shall be applied to convert the token identifier type
of the **TOKEN** to:

The reduction shall then proceed based upon the token identifier type yielded by the rule applied. When more than
one rule applies, the highest numbered rule shall apply (which in turn may refer to another rule). (Note that except in rule 7, the
presence of an `'='` in the token has no effect.)

The **WORD** tokens shall have the word expansion rules applied to them immediately before the associated
command is executed, not at the time the command is parsed.

When the **TOKEN** is exactly a reserved word, the token identifier for that reserved word shall result.
Otherwise, the token **WORD** shall be returned. Also, if the parser is in any state where only a reserved word could be the
next correct token, proceed as above.

Rule 1 is not directly referenced in the grammar, but is referred to by other rules, or applies globally.

The expansions specified in 2.7 Redirection shall occur. As specified there, exactly one field can result (or the result is unspecified), and there are additional requirements on pathname expansion.

Quote removal shall be applied to the word to determine the delimiter that is used to find the end of the here-document that begins after the next <newline>.

When the **TOKEN** is exactly the reserved word **esac**, the token identifier for **esac** shall result.
Otherwise, the token **WORD** shall be returned.

When the **TOKEN** meets the requirements for a name (see XBD *3.216 Name*), the token identifier **NAME** shall result. Otherwise, the
token **WORD** shall be returned.

When the **TOKEN** is exactly the reserved word **in**, the token identifier for **in** shall result.
Otherwise, the token **WORD** shall be returned.

When the **TOKEN** is exactly the reserved word **in** or **do**, the token identifier for **in** or
**do** shall result, respectively. Otherwise, the token **WORD** shall be returned.

(For a. and b.: As indicated in the grammar, a *linebreak* precedes the tokens **in** and **do**. If
<newline> characters are present at the indicated location, it is the token after them that is treated in this fashion.)

If the **TOKEN** is exactly a reserved word, the token identifier for that reserved word shall result.
Otherwise, 7b shall be applied.

If the **TOKEN** contains an unquoted (as determined while applying rule 4 from 2.3 Token
Recognition) <equals-sign> character that is not part of an embedded parameter expansion, command substitution, or
arithmetic expansion construct (as determined while applying rule 5 from 2.3 Token Recognition):

Otherwise, the token **WORD** shall be returned.

If a returned **ASSIGNMENT_WORD** token begins with a valid name, assignment of the value after the first
<equals-sign> to the name shall occur as specified in 2.9.1 Simple Commands. If a returned
**ASSIGNMENT_WORD** token does not begin with a valid name, the way in which the token is processed is unspecified.

When the **TOKEN** is exactly a reserved word, the token identifier for that reserved word shall result.
Otherwise, when the **TOKEN** meets the requirements for a name, the token identifier **NAME** shall result. Otherwise, rule
7 applies.

Word expansion and assignment shall never occur, even when required by the rules above, when this rule is being
parsed. Each **TOKEN** that might either be expanded or have assignment applied to it shall instead be returned as a single
**WORD** consisting only of characters that are exactly the token described in 2.3 Token Recognition
.

/* ------------------------------------------------------- The grammar symbols ------------------------------------------------------- */ %token WORD %token ASSIGNMENT_WORD %token NAME %token NEWLINE %token IO_NUMBER %token IO_LOCATION

%token AND_IF OR_IF DSEMI SEMI_AND /* '&&' '||' ';;' ';&' */

%token DLESS DGREAT LESSAND GREATAND LESSGREAT DLESSDASH /* '<<' '>>' '<&' '>&' '<>' '<<-' */

%token CLOBBER /* '>|' */

/* The following are the reserved words. */

%token If Then Else Elif Fi Do Done /* 'if' 'then' 'else' 'elif' 'fi' 'do' 'done' */

%token Case Esac While Until For /* 'case' 'esac' 'while' 'until' 'for' */

/* These are reserved words, not operator tokens, and are recognized when reserved words are recognized. */

%token Lbrace Rbrace Bang /* '{' '}' '!' */

%token In /* 'in' */

/* ------------------------------------------------------- The Grammar ------------------------------------------------------- */ %start program %% program : linebreak complete_commands linebreak | linebreak ; complete_commands: complete_commands newline_list complete_command | complete_command ; complete_command : list separator_op | list ; list : list separator_op and_or | and_or ; and_or : pipeline | and_or AND_IF linebreak pipeline | and_or OR_IF linebreak pipeline ; pipeline : pipe_sequence | Bang pipe_sequence ; pipe_sequence : command | pipe_sequence '|' linebreak command ; command : simple_command | compound_command | compound_command redirect_list | function_definition ; compound_command : brace_group | subshell | for_clause | case_clause | if_clause | while_clause | until_clause ; subshell : '(' compound_list ')' ; compound_list : linebreak term | linebreak term separator ; term : term separator and_or | and_or ; for_clause : For name do_group | For name sequential_sep do_group | For name linebreak in sequential_sep do_group | For name linebreak in wordlist sequential_sep do_group ; name : NAME /* Apply rule 5 */ ; in : In /* Apply rule 6 */ ; wordlist : wordlist WORD | WORD ; case_clause : Case WORD linebreak in linebreak case_list Esac | Case WORD linebreak in linebreak case_list_ns Esac | Case WORD linebreak in linebreak Esac ; case_list_ns : case_list case_item_ns | case_item_ns ; case_list : case_list case_item | case_item ; case_item_ns : pattern_list ')' linebreak | pattern_list ')' compound_list ; case_item : pattern_list ')' linebreak DSEMI linebreak | pattern_list ')' compound_list DSEMI linebreak | pattern_list ')' linebreak SEMI_AND linebreak | pattern_list ')' compound_list SEMI_AND linebreak ; pattern_list : WORD /* Apply rule 4 */ | '(' WORD /* Do not apply rule 4 */ | pattern_list '|' WORD /* Do not apply rule 4 */ ; if_clause : If compound_list Then compound_list else_part Fi | If compound_list Then compound_list Fi ; else_part : Elif compound_list Then compound_list | Elif compound_list Then compound_list else_part | Else compound_list ; while_clause : While compound_list do_group ; until_clause : Until compound_list do_group ; function_definition : fname '(' ')' linebreak function_body ; function_body : compound_command /* Apply rule 9 */ | compound_command redirect_list /* Apply rule 9 */ ; fname : NAME /* Apply rule 8 */ ; brace_group : Lbrace compound_list Rbrace ; do_group : Do compound_list Done /* Apply rule 6 */ ; simple_command : cmd_prefix cmd_word cmd_suffix | cmd_prefix cmd_word | cmd_prefix | cmd_name cmd_suffix | cmd_name ; cmd_name : WORD /* Apply rule 7a */ ; cmd_word : WORD /* Apply rule 7b */ ; cmd_prefix : io_redirect | cmd_prefix io_redirect | ASSIGNMENT_WORD | cmd_prefix ASSIGNMENT_WORD ; cmd_suffix : io_redirect | cmd_suffix io_redirect | WORD | cmd_suffix WORD ; redirect_list : io_redirect | redirect_list io_redirect ; io_redirect : io_file | IO_NUMBER io_file | IO_LOCATION io_file /* Optionally supported */ | io_here | IO_NUMBER io_here | IO_LOCATION io_here /* Optionally supported */ ; io_file : '<' filename | LESSAND filename | '>' filename | GREATAND filename | DGREAT filename | LESSGREAT filename | CLOBBER filename ; filename : WORD /* Apply rule 2 */ ; io_here : DLESS here_end | DLESSDASH here_end ; here_end : WORD /* Apply rule 3 */ ; newline_list : NEWLINE | newline_list NEWLINE ; linebreak : newline_list | /* empty */ ; separator_op : '&' | ';' ; separator : separator_op linebreak | newline_list ; sequential_sep : ';' linebreak | newline_list ;

Job control is defined (see XBD *3.181 Job Control*) as a
facility that allows users selectively to stop (suspend) the execution of processes and continue (resume) their execution at a
later point. It is jointly supplied by the terminal I/O driver and a command interpreter. The shell is one such command interpreter
and job control in the shell is enabled by set **-m** (which is enabled by default in interactive
shells). The remainder of this section describes the job control facility provided by the shell. Requirements relating to
background jobs stated in this section only apply to job-control background jobs.

If the shell has a controlling terminal and it is the controlling process for the terminal session, it shall initially set the foreground process group ID associated with the terminal to its own process group ID. Otherwise, if it has a controlling terminal, it shall initially perform the following steps if interactive and may perform them if non-interactive:

Subsequently, the shell shall change the foreground process group associated with its controlling terminal when a foreground job is running as noted in the description below.

When job control is enabled, the shell shall create one or more jobs when it executes a list (see 2.9.3 Lists) that has one of the following forms:

For the purposes of job control, a list that includes more than one asynchronous AND-OR list shall be treated as if it were split into multiple separate lists, each ending with an asynchronous AND-OR list.

When a job consisting of a single asynchronous AND-OR list is created, it shall form a *background job* and
the associated process ID shall be that of a child process that is made a process group leader, with all other processes (if any)
that the shell creates to execute the AND-OR list initially having this process ID as their process group ID.

For a list consisting of one or more sequentially executed AND-OR lists followed by at most one asynchronous AND-OR
list, the whole list shall form a single *foreground job* up until the sequentially executed AND-OR lists have all completed
execution, at which point the asynchronous AND-OR list (if any) shall form a background job as described above.

For each pipeline in a foreground job, if the pipeline is executed while the list is still a foreground job, the set of processes comprising the pipeline, and any processes descended from it, shall all be in the same process group, unless the shell executes some of the commands in the pipeline in the current shell execution environment and others in a subshell environment; in this case the process group ID of the current shell need not change (or cannot change if it is the session leader), and consequently the process group ID that the other processes all share may differ from the process group ID of the current shell (which means that a SIGSTOP, SIGTSTP, SIGTTIN, or SIGTTOU signal sent to one of those process groups does not cause the whole pipeline to stop).

A background job that was created on execution of an asynchronous AND-OR list can be brought into the foreground by
means of the *fg* utility (if supported); in this case the entire job shall become a single
foreground job. If a process that the shell subsequently waits for is part of this foreground job and is stopped by a signal, the
entire job shall become a suspended job and the behavior shall be as if the process had been stopped while the job was running in
the background.

When a foreground job is created, or a background job is brought into the foreground by the *fg* utility, if the shell has a controlling terminal it shall set the foreground process group ID
associated with the terminal as follows:

When a foreground job terminates, or becomes a suspended job (see below), if the shell has a controlling terminal it shall set the foreground process group ID associated with the terminal to the process group ID of the shell.

Each background job (whether suspended or not) shall have associated with it a job number and a process ID that is
known in the current shell execution environment. When a background job is brought into the foreground by means of the *fg* utility, the associated job number shall be removed from the shell's background jobs list and
the associated process ID shall be removed from the list of process IDs known in the current shell execution environment.

If a process that the shell is waiting for is part of a foreground job that was started as a foreground job and is stopped by a catchable signal (SIGTSTP, SIGTTIN, or SIGTTOU):

command1 && command2 | { command3 || command4; } & fg

If a process that the shell is waiting for is part of a foreground job that was started as a foreground job and is stopped by a SIGSTOP signal, the behavior shall be as described above for a catchable signal unless the shell was executing a built-in utility in the current shell execution environment when the SIGSTOP was delivered, resulting in the shell itself being stopped by the signal, in which case if the shell subsequently receives a SIGCONT signal and has one or more child processes that remain stopped, the shell shall create a suspended job as if only those child processes had been stopped.

When a suspended job is created as a result of a foreground job being stopped, it shall be assigned a job number,
and an interactive shell shall write, and a non-interactive shell may write, a message to standard error, formatted as described by
the *jobs* utility (without the **-l** option) for a suspended job. The message may
indicate that the commands comprising the job include commands that have already completed; in this case the completed commands
shall not be repeated if execution of the job is subsequently continued. If the shell is interactive, it shall save the terminal
settings before changing them to the settings it needs to read further commands.

When a process associated with a background job is stopped by a SIGSTOP, SIGTSTP, SIGTTIN, or SIGTTOU signal, the
shell shall convert the (non-suspended) background job into a suspended job and an interactive shell shall write a message to
standard error, formatted as described by the *jobs* utility (without the **-l**
option) for a suspended job, at the following time:

Execution of a suspended job can be continued as a foreground job by means of the *fg* utility (if supported), or as a (non-suspended) background job either by means of the *bg* utility (if supported) or by sending the stopped processes a SIGCONT signal. The *fg* and *bg* utilities shall send a SIGCONT signal to the
process group of the process(es) whose stopped wait status caused the shell to suspend the job. If the shell has a controlling
terminal, the *fg* utility shall send the SIGCONT signal after it has set the foreground
process group ID associated with the terminal (see above). If the *fg* utility is used from
an interactive shell to bring into the foreground a suspended job that was created from a foreground job, before it sends the
SIGCONT signal the *fg* utility shall restore the terminal settings to the ones that the
shell saved when the job was suspended.

When a background job completes or is terminated by a signal, an interactive shell shall write a message to
standard error, formatted as described by the *jobs* utility (without the **-l**
option) for a job that completed or was terminated by a signal, respectively, at the following time:

In each case above where an interactive shell writes a message immediately prior to writing the next prompt for input, the same message may also be written by a non-interactive shell, at any of the following times:

If job control is disabled (see the description of *set* **-m**) when the shell executes
an asynchronous AND-OR list, the commands in the list shall inherit from the shell a signal action of ignored (SIG_IGN) for the
SIGINT and SIGQUIT signals. In all other cases, commands executed by the shell shall inherit the same signal actions as those
inherited by the shell from its parent unless a signal action is modified by the *trap* special built-in
(see trap)

When a signal for which a trap has been set is received while the shell is waiting for the completion of a utility
executing a foreground command, the trap associated with that signal shall not be executed until after the foreground command has
completed. When the shell is waiting, by means of the *wait* utility, for asynchronous
commands to complete, the reception of a signal for which a trap has been set shall cause the *wait* utility to return immediately with an exit status >128, immediately after which the
trap associated with that signal shall be taken.

If multiple signals are pending for the shell for which there are associated trap actions, the order of execution of trap actions is unspecified.

A shell execution environment consists of the following:

Utilities other than the special built-ins (see 2.15 Special Built-In Utilities) shall be invoked in a separate environment that consists of the following. The initial value of these objects shall be the same as that for the parent shell, except as noted below.

The environment of the shell process shall not be changed by the utility unless explicitly specified by the utility
description (for example, *cd* and *umask*).

A subshell environment shall be created as a duplicate of the shell environment, except that:

Changes made to the subshell environment shall not affect the shell environment. Command substitution, commands that are grouped with parentheses, and asynchronous AND-OR lists shall be executed in a subshell environment. Additionally, each command of a multi-command pipeline is in a subshell environment; as an extension, however, any or all commands in a pipeline may be executed in the current environment. Except where otherwise stated, all other commands shall be executed in the current shell environment.

The pattern matching notation described in this section is used to specify patterns for matching character strings
in the shell. This notation is also used by some other utilities (*find*, *pax*, and optionally *make*) and by some system
interfaces (*fnmatch*(), *glob*(), and
*wordexp*()). Historically, pattern matching notation is related to, but slightly
different from, the regular expression notation described in XBD *9. Regular
Expressions*. For this reason, the description of the rules for this pattern matching notation are based on the description
of regular expression notation, modified to account for the differences.

If an attempt is made to use pattern matching notation to match a string that contains one or more bytes that do not form part of a valid character, the behavior is unspecified. Since pathnames can contain such bytes, portable applications need to ensure that the current locale is the C or POSIX locale when performing pattern matching (or expansion) on arbitrary pathnames.

The following patterns shall match a single character: ordinary characters, special pattern characters, and pattern bracket expressions. The pattern bracket expression also shall match a single collating element.

In a pattern, or part of one, where a shell-quoting <backslash> can be used, a <backslash> character
shall escape the following character as described in 2.2.1 Escape Character (Backslash), regardless of
whether or not the <backslash> is inside a bracket expression. (The sequence `"\\"` represents one literal
<backslash>.)

In a pattern, or part of one, where a shell-quoting <backslash> cannot be used to preserve the literal value of a character that would otherwise be treated as special:

All of the requirements and effects of quoting on ordinary, shell special, and special pattern characters shall
apply to escaping in this context, except where specified otherwise. (Situations where this applies include word expansions when a
pattern used in pathname expansion is not present in the original word but results from an earlier expansion, or the argument to
the *find* -*name* or -*path* primary as passed to *find*, or the *pattern* argument to the *fnmatch*() and *glob*() functions when FNM_NOESCAPE
or GLOB_NOESCAPE is not set in *flags*, respectively.)

If a pattern ends with an unescaped <backslash>, the behavior is unspecified.

An ordinary character is a pattern that shall match itself. In a pattern, or part of one, where a shell-quoting <backslash> can be used, an ordinary character can be any character in the supported character set except for NUL, those special shell characters in 2.2 Quoting that require quoting, and the three special pattern characters described below. In a pattern, or part of one, where a shell-quoting <backslash> cannot be used to preserve the literal value of a character that would otherwise be treated as special, an ordinary character can be any character in the supported character set except for NUL and the three special pattern characters described below. Matching shall be based on the bit pattern used for encoding the character, not on the graphic representation of the character. If any character (ordinary, shell special, or pattern special) is quoted, or escaped with a <backslash>, that pattern shall match the character itself. The application shall ensure that it quotes or escapes any character that would otherwise be treated as special, in order for it to be matched as an ordinary character.

When unquoted, unescaped, and not inside a bracket expression, the following three characters shall have special meaning in the specification of patterns:

The following rules are used to construct patterns matching multiple characters from patterns matching a single character:

The rules described so far in 2.14.1 Patterns Matching a Single Character and 2.14.2 Patterns Matching Multiple Characters are qualified by the following rules that apply when pattern matching notation is used for filename expansion:

It is unspecified whether an explicit <period> in a bracket expression matching list, such as
`"[.abc]"`, can match a leading <period> in a filename.

For example, given the pattern:

/foo/bar/x*/bam

search permission is needed for directories **/** and **foo**, search and read permissions are needed for
directory **bar**, and search permission is needed for each **x*** directory.

If the pattern matches any existing filenames or pathnames, the pattern shall be replaced with those filenames and
pathnames, sorted according to the collating sequence in effect in the current locale. If this collating sequence does not have a
total ordering of all characters (see XBD *7.3.2 LC_COLLATE*), any
filenames or pathnames that collate equally shall be further compared byte-by-byte using the collating sequence for the POSIX
locale.

If the pattern contains an open bracket (`'['`) that does not introduce a bracket expression as in XBD
*9.3.5 RE Bracket Expression*, it is unspecified whether other unquoted
`'*'`, `'?'`, `'['` or <backslash> characters within the same slash-delimited component of the pattern
retain their special meanings or are treated as ordinary characters. For example, the pattern `"a*[/b*"` may match all
filenames beginning with `'b'` in the directory `"a*["` or it may match all filenames beginning with `'b'` in
all directories with names beginning with `'a'` and ending with `'['`.

If the pattern does not match any existing filenames or pathnames, the pattern string shall be left unchanged.

The following "special built-in" utilities shall be supported in the shell command language. The output of each command, if any, shall be written to standard output, subject to the normal redirection and piping possible with all commands.

The term "built-in" implies that there is no need to execute a separate executable file because the utility is implemented in the shell itself. An implementation may choose to make any utility a built-in; however, the special built-in utilities described here differ from regular built-in utilities in two respects:

The special built-in utilities in this section need not be provided in a manner accessible via the *exec*
family of functions defined in the System Interfaces volume of POSIX.1-2024.

Some of the special built-ins are described as conforming to XBD *12.2 Utility Syntax Guidelines*. For those that are not, the requirement in
*1.4 Utility Description Defaults* that `"--"` be recognized as a
first argument to be discarded does not apply and a conforming application shall not use that argument.

return to top of page

break — exit from for, while, or until loop

break[n]

If

nis specified, thebreakutility shall exit from thenth enclosingfor,while, oruntilloop. Ifnis not specified,breakshall behave as ifnwas specified as 1. Execution shall continue with the command immediately following the exited loop. The application shall ensure that the value ofnis a positive decimal integer. Ifnis greater than the number of enclosing loops, the outermost enclosing loop shall be exited. If there is no enclosing loop, the behavior is unspecified.A loop shall enclose a

breakorcontinuecommand if the loop lexically encloses the command. A loop lexically encloses abreakorcontinuecommand if the command is:

- Executing in the same execution environment (see 2.13 Shell Execution Environment) as the compound-list of the loop's do-group (see 2.10.2 Shell Grammar Rules), and
- Contained in a compound-list associated with the loop (either in the compound-list of the loop's do-group or, if the loop is a
whileoruntilloop, in the compound-list following thewhileoruntilreserved word), and- Not in the body of a function whose function definition command (see 2.9.5 Function Definition Command) is contained in a compound-list associated with the loop.
If

nis greater than the number of lexically enclosing loops and there is a non-lexically enclosing loop in progress in the same execution environment as thebreakorcontinuecommand, it is unspecified whether that loop encloses the command.

None.

See the DESCRIPTION.

Not used.

None.

None.

Default.

Not used.

The standard error shall be used only for diagnostic messages.

None.

None.

- 0
- Successful completion.
- >0
- The
nvalue was not an unsigned decimal integer greater than or equal to 1.

Default.

None.

for i in * do if test -d "$i" then break fi doneThe results of running the following example are unspecified: there are two loops in progress when the

breakcommand is executed, and they are in the same execution environment, but neither loop is lexically enclosing thebreakcommand. (There are no loops lexically enclosing thecontinuecommands, either.)foo() { for j in 1 2; do echo 'break 2' >/tmp/do_break echo " sourcing /tmp/do_break ($j)..." # the behavior of the break from running the following command # results in unspecified behavior: . /tmp/do_break

do_continue() { continue 2; } echo " running do_continue ($j)..." # the behavior of the continue in the following function call # results in unspecified behavior (if execution reaches this # point): do_continue

trap 'continue 2' USR1 echo " sending SIGUSR1 to self ($j)..." # the behavior of the continue in the trap invoked from the # following signal results in unspecified behavior (if # execution reaches this point): kill -s USR1 $$ sleep 1 done } for i in 1 2; do echo "running foo ($i)..." foo done

In early proposals, consideration was given to expanding the syntax of

breakandcontinueto refer to a label associated with the appropriate loop as a preferable alternative to thenmethod. However, this volume of POSIX.1-2024 does reserve the name space of command names ending with a <colon>. It is anticipated that a future implementation could take advantage of this and provide something like:outofloop: for i in a b c d e do for j in 0 1 2 3 4 5 6 7 8 9 do if test -r "${i}${j}" then break outofloop fi done doneand that this might be standardized after implementation experience is achieved.

None.

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.

POSIX.1-2008, Technical Corrigendum 2, XCU/TC2-2008/0046 [842] is applied.

Austin Group Defect 1058 is applied, clarifying that the requirement for

nto be a positive decimal integer is a requirement on the application.

colon — null utility

:[argument...]

This utility shall do nothing except return a 0 exit status. It is used when a command is needed, as in the

thencondition of anifcommand, but nothing is to be done by the command.

This utility shall not recognize the

"--"argument in the manner specified by Guideline 10 of XBD12.2 Utility Syntax Guidelines.Implementations shall not support any options.

See the DESCRIPTION.

Not used.

None.

None.

Default.

Not used.

Not used.

None.

None.

Zero.

None.

See the APPLICATION USAGE for

true.

: "${X=abc}" if false then : else printf '%s\n' "$X" fiabcAs with any of the special built-ins, the null utility can also have variable assignments and redirections associated with it, such as:

x=y : > zwhich sets variable

xto the valuey(so that it persists after the null utility completes) and creates or truncates filez; if the file cannot be created or truncated, a non-interactive shell exits (see 2.8.1 Consequences of Shell Errors).

None.

None.

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.

SD5-XCU-ERN-97 is applied, updating the SYNOPSIS.

Austin Group Defect 1272 is applied, clarifying that the null utility does not process its arguments, does not recognize the

"--"end-of-options delimiter, does not support any options, and does not write to standard error.Austin Group Defect 1640 is applied, changing the APPLICATION USAGE section.

continue — continue for, while, or until loop

continue[n]

If

nis specified, thecontinueutility shall return to the top of thenth enclosingfor,while, oruntilloop. Ifnis not specified,continueshall behave as ifnwas specified as 1. Returning to the top of the loop involves repeating the condition list of awhileoruntilloop or performing the next assignment of aforloop, and re-executing the loop if appropriate.The application shall ensure that the value of

nis a positive decimal integer. Ifnis greater than the number of enclosing loops, the outermost enclosing loop shall be used. If there is no enclosing loop, the behavior is unspecified.The meaning of "enclosing" shall be as specified in the description of the

breakutility.

None.

See the DESCRIPTION.

Not used.

None.

None.

Default.

Not used.

The standard error shall be used only for diagnostic messages.

None.

None.

- 0
- Successful completion.
- >0
- The
nvalue was not an unsigned decimal integer greater than or equal to 1.

Default.

None.

for i in * do if test -d "$i" then continue fi printf '"%s" is not a directory.\n' "$i" done

None.

None.

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.

The example is changed to use the

printfutility rather thanecho.POSIX.1-2008, Technical Corrigendum 2, XCU/TC2-2008/0046 [842] is applied.

Austin Group Defect 1058 is applied, clarifying that the requirement for

nto be a positive decimal integer is a requirement on the application.

dot — execute commands in the current environment

.file

The shell shall tokenize (see 2.3 Token Recognition) the contents of the

file, parse the tokens (see 2.10 Shell Grammar), and execute the resulting commands in the current environment. It is unspecified whether the commands are parsed and executed as aprogram(as for a shell script) or are parsed as a singlecompound_listthat is executed after the entire file has been parsed.If

filedoes not contain a <slash>, the shell shall use the search path specified byPATHto find the directory containingfile. Unlike normal command search, however, the file searched for by thedotutility need not be executable. If no readable file is found, a non-interactive shell shall abort; an interactive shell shall write a diagnostic message to standard error.The

dotspecial built-in shall support XBD12.2 Utility Syntax Guidelines, except for Guidelines 1 and 2.

None.

See the DESCRIPTION.

Not used.

See the DESCRIPTION.

See the DESCRIPTION.

Default.

Not used.

The standard error shall be used only for diagnostic messages.

None.

None.

If no readable file was found or if the commands in the file could not be parsed, and the shell is interactive (and therefore does not abort; see 2.8.1 Consequences of Shell Errors), the exit status shall be non-zero. Otherwise, return the value of the last command executed, or a zero exit status if no command is executed.

Default.

None.

cat foobarfoo=hello bar=world. ./foobar echo $foo $barhello world

Some older implementations searched the current directory for the

file, even if the value ofPATHdisallowed it. This behavior was omitted from this volume of POSIX.1-2024 due to concerns about introducing the susceptibility to trojan horses that the user might be trying to avoid by leavingdotout ofPATH .The KornShell version of

dottakes optional arguments that are set to the positional parameters. This is a valid extension that allows adotscript to behave identically to a function.

None.

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.

SD5-XCU-ERN-164 is applied.

POSIX.1-2008, Technical Corrigendum 1, XCU/TC1-2008/0038 [114] and XCU/TC1-2008/0039 [214] are applied.

Austin Group Defect 252 is applied, adding a requirement for

dotto support XBD12.2 Utility Syntax Guidelines(except for Guidelines 1 and 2, since the utility's name is'.').Austin Group Defect 953 is applied, clarifying how the commands in the

fileare parsed.Austin Group Defect 1265 is applied, updating the DESCRIPTION to align with the changes made to 2.8.1 Consequences of Shell Errors between Issue 6 and Issue 7.

eval — construct command by concatenating arguments

eval[argument...]

The

evalutility shall construct a command string by concatenatingarguments together, separating each with a <space> character. The constructed command string shall be tokenized (see 2.3 Token Recognition), parsed (see 2.10 Shell Grammar), and executed by the shell in the current environment. It is unspecified whether the commands are parsed and executed as aprogram(as for a shell script) or are parsed as a singlecompound_listthat is executed after the entire constructed command string has been parsed.

None.

See the DESCRIPTION.

Not used.

None.

None.

Default.

Not used.

The standard error shall be used only for diagnostic messages.

None.

None.

If there are no

arguments, or only nullarguments,evalshall return a zero exit status; otherwise, it shall return the exit status of the command defined by the string of concatenatedarguments separated by <space> characters, or a non-zero exit status if the concatenation could not be parsed as a command and the shell is interactive (and therefore did not abort).

Default.

Since

evalis not required to recognize the"--"end of options delimiter, in cases where the argument(s) toevalmight begin with'-'it is recommended that the first argument is prefixed by a string that will not alter the commands to be executed, such as a <space> character:eval " $commands"or:

eval " $(some_command)"

foo=10 x=foo y='$'$x echo $y$fooeval y='$'$x echo $y10

This standard allows, but does not require,

evalto recognize"--". Although this means applications cannot use"--"to protect against options supported as an extension (or errors reported for unsupported options), the nature of theevalutility is such that other means can be used to provide this protection (see APPLICATION USAGE above).

None.

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.

SD5-XCU-ERN-97 is applied, updating the SYNOPSIS.

POSIX.1-2008, Technical Corrigendum 1, XCU/TC1-2008/0040 [114], XCU/TC1-2008/0041 [163], and XCU/TC1-2008/0042 [163] are applied.

Austin Group Defect 953 is applied, clarifying how the commands in the constructed command string are parsed.

exec — perform redirections in the current shell or execute a utility

exec[utility[argument...]]

If

execis specified with no operands, any redirections associated with theexeccommand shall be made in the current shell execution environment. If any file descriptors with numbers greater than 2 are opened by those redirections, it is unspecified whether those file descriptors remain open when the shell invokes another utility. Scripts concerned that child shells could misuse open file descriptors can always close them explicitly, as shown in one of the following examples. If the result of the redirections would be that file descriptor 0, 1, or 2 is closed, implementations may open the file descriptor to an unspecified file.If

execis specified with autilityoperand, the shell shall execute a non-built-in utility as described in 2.9.1.6 Non-built-in Utility Execution withutilityas the command name and theargumentoperands (if any) as the command arguments.If the

execcommand fails, a non-interactive shell shall exit from the current shell execution environment;[UP]an interactive shell may exit from a subshell environment but shall not exit if the current shell environment is not a subshell environment.If the

execcommand fails and the shell does not exit, any redirections associated with theexeccommand that were successfully made shall take effect in the current shell execution environment.The

execspecial built-in shall support XBD12.2 Utility Syntax Guidelines.

None.

See the DESCRIPTION.

Not used.

None.

The following environment variable shall affect the execution of

exec:

PATH- Determine the search path when looking for the utility given as the
utilityoperand; see XBD8.3 Other Environment Variables.

Default.

Not used.

The standard error shall be used only for diagnostic messages.

None.

None.

If

utilityis specified and is executed,execshall not return to the shell; rather, the exit status of the current shell execution environment shall be the exit status ofutility. Ifutilityis specified and an attempt to execute it as a non-built-in utility fails, the exit status shall be as described in 2.9.1.6 Non-built-in Utility Execution. If a redirection error occurs (see 2.8.1 Consequences of Shell Errors), the exit status shall be a value in the range 1-125. Otherwise,execshall return a zero exit status.

Default.

None.

Open

readfileas file descriptor 3 for reading:exec 3< readfileOpen

writefileas file descriptor 4 for writing:exec 4> writefileMake file descriptor 5 a copy of file descriptor 0:

exec 5<&0Close file descriptor 3:

exec 3<&-Cat the file

maggieby replacing the current shell with thecatutility:exec cat maggieAn application that is not concerned with strict conformance can make use of optional

%gsupport known to be present in the implementation'sprintfutility by ensuring that any shell built-in version is not executed instead, and using a subshell so that the shell continues afterwards:(exec printf '%g\n' "$float_value")

Most historical implementations were not conformant in that:

foo=bar exec cmddid not pass

footocmd.

None.

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.

SD5-XCU-ERN-97 is applied, updating the SYNOPSIS.

Austin Group Defect 252 is applied, adding a requirement for

execto support XBD12.2 Utility Syntax Guidelines.Austin Group Defect 1157 is applied, clarifying the execution of non-built-in utilities.

Austin Group Defect 1587 is applied, changing the ENVIRONMENT VARIABLES section.

exit — cause the shell to exit

exit[n]

The

exitutility shall cause the shell to exit from its current execution environment. If the current execution environment is a subshell environment, the shell shall exit from the subshell environment and continue in the environment from which that subshell environment was invoked; otherwise, the shell utility shall terminate. The wait status of the shell or subshell shall be determined by the unsigned decimal integern, if specified.If

nis specified and has a value between 0 and 255 inclusive, the wait status of the shell or subshell shall indicate that it exited with exit statusn. Ifnis specified and has a value greater than 256 that corresponds to an exit status the shell assigns to commands terminated by a valid signal (see 2.8.2 Exit Status for Commands), the wait status of the shell or subshell shall indicate that it was terminated by that signal. No other actions associated with the signal, such as execution oftrapactions or creation of a core image, shall be performed by the shell.If

nis specified and is not an unsigned decimal integer, or has a value of 256, or has a value greater than 256 but not corresponding to an exit status the shell assigns to commands terminated by a valid signal, the wait status of the shell or subshell is unspecified.If

nis not specified, the result shall be as ifnwere specified with the current value of the special parameter'?'(see 2.5.2 Special Parameters), except that if theexitcommand would cause the end of execution of atrapaction, the value for the special parameter'?'that is considered "current" shall be the value it had immediately preceding thetrapaction.A

trapaction onEXITshall be executed before the shell terminates, except when theexitutility is invoked in thattrapaction itself, in which case the shell shall exit immediately. It is unspecified whether setting a newtrapaction onEXITduring execution of atrapaction onEXITwill cause the newtrapaction to be executed before the shell terminates.

None.

See the DESCRIPTION.

Not used.

None.

None.

Default.

Not used.

The standard error shall be used only for diagnostic messages.

None.

None.

The

exitutility causes the shell to exit from its current execution environment, and therefore does not itself return an exit status.

Default.

As explained in other sections, certain exit status values have been reserved for special uses and should be used by applications only for those purposes:

- 126
- A file to be executed was found, but it was not an executable utility.
- 127
- A utility to be executed was not found.
- 128
- An unrecoverable read error was detected by the shell while reading commands, except from the
fileoperand of thedotspecial built-in.- >128
- A command was interrupted by a signal.

Exit with a

truevalue:exit 0Exit with a

falsevalue:exit 1Propagate error handling from within a subshell:

( command1 || exit 1 command2 || exit 1 exec command3 ) > outputfile || exit 1 echo "outputfile created successfully"

The behavior of

exitwhen given an invalid argument or unknown option is unspecified, because of differing practices in the various historical implementations. A value larger than 255 might be truncated by the shell, and be unavailable even to a parent process that useswaitid() to get the full exit value. It is recommended that implementations that detect any usage error should cause a non-zero exit status (or, if the shell is interactive and the error does not cause the shell to abort, store a non-zero value in"$?"), but even this was not done historically in all shells.See also

C.2.8.2 Exit Status for Commands.

None.

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.

POSIX.1-2008, Technical Corrigendum 2, XCU/TC2-2008/0047 [717], XCU/TC2-2008/0048 [960], XCU/TC2-2008/0049 [717], and XCU/TC2-2008/0050 [960] are applied.

Austin Group Defect 51 is applied, specifying the behavior when

nhas a value greater than 256 that corresponds to an exit status the shell assigns to commands terminated by a valid signal.Austin Group Defect 1029 is applied, changing "

trap" to "trapaction" in the DESCRIPTION section.Austin Group Defect 1309 is applied, changing the EXIT STATUS section.

Austin Group Defect 1425 is applied, clarifying the requirements for a

trapaction onEXIT.Austin Group Defect 1602 is applied, clarifying the behavior of

exitin atrapaction.Austin Group Defect 1629 is applied, adding exit status 128 to the APPLICATION USAGE section.

export — set the export attribute for variables

export name[=word]...

export -p

The shell shall give the

exportattribute to the variables corresponding to the specifiednames, which shall cause them to be in the environment of subsequently executed commands. If the name of a variable is followed by =word, then the value of that variable shall be set toword.The

exportspecial built-in shall be a declaration utility. Therefore, ifexportis recognized as the command name of a simple command, then subsequent words of the formname=wordshall be expanded in an assignment context. See 2.9.1.1 Order of Processing.The

exportspecial built-in shall support XBD12.2 Utility Syntax Guidelines.When

-pis specified,exportshall write to the standard output the names and values of all exported variables, in the following format:"export %s=%s\n", <name>, <value>if

nameis set, and:"export %s\n", <name>if

nameis unset.The shell shall format the output, including the proper use of quoting, so that it is suitable for reinput to the shell as commands that achieve the same exporting results, except:

- Read-only variables with values cannot be reset.
- Variables that were unset at the time they were output need not be reset to the unset state if a value is assigned to the variable between the time the state was saved and the time at which the saved output is reinput to the shell.
When no arguments are given, the results are unspecified.

See the DESCRIPTION.

See the DESCRIPTION.

Not used.

None.

None.

Default.

See the DESCRIPTION.

The standard error shall be used only for diagnostic messages.

None.

None.

- 0
- Successful completion.
- >0
- At least one operand could not be processed as requested, such as a
nameoperand that could not be exported or an attempt to modify areadonlyvariable using aname=wordoperand, or the-poption was specified and a write error occurred.

Default.

Note that, unless

Xwas previously marked readonly, the value of"$?"after:export X=$(false)will be 0 (because

exportsuccessfully setXto the empty string) and that execution continues, even ifset-eis in effect. In order to detect command substitution failures, a user must separate the assignment from the export, as in:X=$(false) export XIn shells that support extended assignment syntax, for example to allow an array to be populated with a single assignment, such extensions can typically only be used in assignments specified as arguments to

exportif the command word is literallyexport, and not if it is some other word that expands toexport. For example:# Shells that support array assignment as an extension generally # support this: export x=(1 2 3); echo ${x[0]} # outputs 1 # But generally do not support this: e=export; $e x=(1 2 3); echo ${x[0]} # syntax error

Export

PWDandHOMEvariables:export PWD HOMESet and export the

PATHvariable:export PATH="/local/bin:$PATH"Save and restore all exported variables:

export -p > temp-file unseta lot of variables...processing. ./temp-file

Note:- If LANG, LC_CTYPE or LC_ALL are left altered or unset in the above example prior to sourcing
temp-file, the results may be undefined.

Some historical shells use the no-argument case as the functional equivalent of what is required here with

-p. This feature was left unspecified because it is not historical practice in all shells, and some scripts may rely on the now-unspecified results on their implementations. Attempts to specify the-poutput as the default case were unsuccessful in achieving consensus. The-poption was added to allow portable access to the values that can be saved and then later restored using; for example, adotscript.Some implementations extend the shell's assignment syntax, for example to allow an array to be populated with a single assignment, and in order for such an extension to be usable in assignments specified as arguments to

exportthese shells haveexportas a separate token in their grammar. This standard only permits an extension of this nature when the input to the shell would contain a syntax error according to the standard grammar. Note that althoughexportcan be a separate token in the shell's grammar, it cannot be a reserved word sinceexportis a candidate for alias substitution whereas reserved words are not (see 2.3.1 Alias Substitution).

None.

2.9.1.1 Order of Processing, 2.15 Special Built-In Utilities

IEEE PASC Interpretation 1003.2 #203 is applied, clarifying the format when a variable is unset.

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/6 is applied, adding the following text to the end of the first paragraph of the DESCRIPTION: "If the name of a variable is followed by =

word, then the value of that variable shall be set toword.". The reason for this change is that the SYNOPSIS forexportincludes:export name[=word]...but the meaning of the optional "=

word" is never explained in the text.

POSIX.1-2008, Technical Corrigendum 1, XCU/TC1-2008/0043 [352] is applied.

POSIX.1-2008, Technical Corrigendum 2, XCU/TC2-2008/0051 [654] and XCU/TC2-2008/0052 [960] are applied.

Austin Group Defect 351 is applied, requiring

exportto be a declaration utility.Austin Group Defect 367 is applied, changing the EXIT STATUS section.

Austin Group Defect 1258 is applied, changing the EXAMPLES section.

Austin Group Defect 1393 is applied, changing the APPLICATION USAGE and RATIONALE sections.

readonly — set the readonly attribute for variables

readonly name[=word]...

readonly -p

The variables whose

names are specified shall be given thereadonlyattribute. The values of variables with thereadonlyattribute cannot be changed by subsequent assignment or use of theexport,getopts,readonly, orreadutilities, nor can those variables be unset by theunsetutility. As described in XBD8.1 Environment Variable Definition, conforming applications shall not request to mark a variable asreadonlyif it is documented as being manipulated by a shell built-in utility, as it may render those utilities unable to complete successfully. If the name of a variable is followed by =word, then the value of that variable shall be set toword.The

readonlyspecial built-in shall be a declaration utility. Therefore, ifreadonlyis recognized as the command name of a simple command, then subsequent words of the formname=wordshall be expanded in an assignment context. See 2.9.1.1 Order of Processing.The

readonlyspecial built-in shall support XBD12.2 Utility Syntax Guidelines.When

-pis specified,readonlywrites to the standard output the names and values of all read-only variables, in the following format:"readonly %s=%s\n", <name>, <value>if

nameis set, and"readonly %s\n", <name>if

nameis unset.The shell shall format the output, including the proper use of quoting, so that it is suitable for reinput to the shell as commands that achieve the same value and

readonlyattribute-setting results in a shell execution environment in which:

- Variables with values at the time they were output do not have the
readonlyattribute set.- Variables that were unset at the time they were output do not have a value at the time at which the saved output is reinput to the shell.
When no arguments are given, the results are unspecified.

See the DESCRIPTION.

See the DESCRIPTION.

Not used.

None.

None.

Default.

See the DESCRIPTION.

The standard error shall be used only for diagnostic messages.

None.

None.

- 0
- Successful completion.
- >0
- At least one operand could not be processed as requested, such as a
nameoperand that could not be markedreadonlyor an attempt to modify an alreadyreadonlyvariable using aname=wordoperand, or the-poption was specified and a write error occurred.

Default.

In shells that support extended assignment syntax, for example to allow an array to be populated with a single assignment, such extensions can typically only be used in assignments specified as arguments to

readonlyif the command word is literallyreadonly, and not if it is some other word that expands toreadonly. For example:# Shells that support array assignment as an extension generally # support this: readonly x=(1 2 3); echo ${x[0]} # outputs 1 # But generally do not support this: r=readonly; $r x=(1 2 3); echo ${x[0]} # syntax error

readonly HOME

Some historical shells preserve the

readonlyattribute across separate invocations. This volume of POSIX.1-2024 allows this behavior, but does not require it.The

-poption allows portable access to the values that can be saved and then later restored using, for example, adotscript. Also see the RATIONALE for export for a description of the no-argument and-poutput cases and a related example.Read-only functions were considered, but they were omitted as not being historical practice or particularly useful. Furthermore, functions must not be read-only across invocations to preclude "spoofing" (spoofing is the term for the practice of creating a program that acts like a well-known utility with the intent of subverting the real intent of the user) of administrative or security-relevant (or security-conscious) shell scripts.

Attempts to set the

readonlyattribute on certain variables, such asPWD ,may have surprising results. Eitherreadonlywill reject the attempt, or the attempt will succeed but the shell will continue to alter the contents ofPWDduring thecdutility, or the attempt will succeed and render thecdutility inoperative (since it must not change directories if it cannot also updatePWD ).Some implementations extend the shell's assignment syntax, for example to allow an array to be populated with a single assignment, and in order for such an extension to be usable in assignments specified as arguments to

readonlythese shells havereadonlyas a separate token in their grammar. This standard only permits an extension of this nature when the input to the shell would contain a syntax error according to the standard grammar. Note that althoughreadonlycan be a separate token in the shell's grammar, it cannot be a reserved word sincereadonlyis a candidate for alias substitution whereas reserved words are not (see 2.3.1 Alias Substitution).

None.

2.9.1.1 Order of Processing, 2.15 Special Built-In Utilities

IEEE PASC Interpretation 1003.2 #203 is applied, clarifying the format when a variable is unset.

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/7 is applied, adding the following text to the end of the first paragraph of the DESCRIPTION: "If the name of a variable is followed by =

word, then the value of that variable shall be set toword.". The reason for this change is that the SYNOPSIS forreadonlyincludes:readonly name[=word]...but the meaning of the optional "=

word" is never explained in the text.

POSIX.1-2008, Technical Corrigendum 2, XCU/TC2-2008/0052 [960] is applied.

Austin Group Defect 351 is applied, requiring

readonlyto be a declaration utility.Austin Group Defect 367 is applied, clarifying that the values of

readonlyvariables cannot be changed by subsequent use of theexport,getopts,readonly, orreadutilities, and changing the EXIT STATUS, EXAMPLES and RATIONALE sections.Austin Group Defect 1393 is applied, changing the APPLICATION USAGE and RATIONALE sections.

return — return from a function or dot script

return[n]

The

returnutility shall cause the shell to stop executing the current function ordotscript. If the shell is not currently executing a function ordotscript, the results are unspecified.

None.

See the DESCRIPTION.

Not used.

None.

None.

Default.

Not used.

The standard error shall be used only for diagnostic messages.

None.

None.

The exit status shall be

n, if specified, except that the behavior is unspecified ifnis not an unsigned decimal integer or is greater than 255. Ifnis not specified, the result shall be as ifnwere specified with the current value of the special parameter'?'(see 2.5.2 Special Parameters), except that if thereturncommand would cause the end of execution of atrapaction, the value for the special parameter'?'that is considered "current" shall be the value it had immediately preceding thetrapaction.

Default.

None.

None.

The behavior of

returnwhen not in a function ordotscript differs between the System V shell and the KornShell. In the System V shell this is an error, whereas in the KornShell, the effect is the same asexit.The results of returning a number greater than 255 are undefined because of differing practices in the various historical implementations. Some shells AND out all but the low-order 8 bits; others allow larger values, but not of unlimited size.

See the discussion of appropriate exit status values under exit.

None.

2.9.5 Function Definition Command, 2.15 Special Built-In Utilities, dot

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.

POSIX.1-2008, Technical Corrigendum 1, XCU/TC1-2008/0044 [214] and XCU/TC1-2008/0045 [214] are applied.

POSIX.1-2008, Technical Corrigendum 2, XCU/TC2-2008/0052 [960] is applied.

Austin Group Defect 1309 is applied, changing the EXIT STATUS section.

Austin Group Defect 1602 is applied, clarifying the behavior of

returnin atrapaction.

set — set or unset options and positional parameters

set[-abCefhmnuvx] [-ooption] [argument...]

set[+abCefhmnuvx] [+ooption] [argument...]

set --[argument...]

set -o

set +o

If no

options orarguments are specified,setshall write the names and values of all shell variables in the collation sequence of the current locale. Eachnameshall start on a separate line, using the format:"%s=%s\n", <name>, <value>The

valuestring shall be written with appropriate quoting; see the description of shell quoting in 2.2 Quoting. The output shall be suitable for reinput to the shell, setting or resetting, as far as possible, the variables that are currently set; read-only variables cannot be reset.When options are specified, they shall set or unset attributes of the shell, as described below. When

arguments are specified, they cause positional parameters to be set or unset, as described below. Setting or unsetting attributes and positional parameters are not necessarily related actions, but they can be combined in a single invocation ofset.The

setspecial built-in shall support XBD12.2 Utility Syntax Guidelinesexcept that options can be specified with either a leading <hyphen-minus> (meaning enable the option) or <plus-sign> (meaning disable it) unless otherwise specified.Implementations shall support the options in the following list in both their <hyphen-minus> and <plus-sign> forms. These options can also be specified as options to

sh.

-a- Set the
exportattribute for all variable assignments. When this option is on, whenever a value is assigned to a variable in the current shell execution environment, theexportattribute shall be set for the variable. This applies to all forms of assignment, including those made as a side-effect of variable expansions or arithmetic expansions, and those made as a result of the operation of thecd,getopts, orreadutilities.

Note:- As discussed in 2.9.1 Simple Commands, not all variable assignments happen in the current execution environment. When an assignment happens in a separate execution environment the
exportattribute is still set for the variable, but that does not affect the current execution environment.
-b- This option shall be supported if the implementation supports the User Portability Utilities option. When job control and
-bare both enabled, the shell shall write asynchronous notifications of background job completions (including termination by a signal), and may write asynchronous notifications of background job suspensions. See 2.11 Job Control for details. When job control is disabled, the-boption shall have no effect. Asynchronous notification shall not be enabled by default.
-C- (Uppercase C.) Prevent existing regular files from being overwritten by the shell's
'>'redirection operator (see 2.7.2 Redirecting Output); the">|"redirection operator shall override thisnoclobberoption for an individual file.
-e- When this option is on, when any command fails (for any of the reasons listed in 2.8.1 Consequences of Shell Errors or by returning an exit status greater than zero), the shell immediately shall exit, as if by executing the
exitspecial built-in utility with no arguments, with the following exceptions:

- The failure of any individual command in a multi-command pipeline, or of any subshell environments in which command substitution was performed during word expansion, shall not cause the shell to exit. Only the failure of the pipeline itself shall be considered.
- The
-esetting shall be ignored when executing the compound list following thewhile,until,if, orelifreserved word, a pipeline beginning with the!reserved word, or any command of an AND-OR list other than the last.- If the exit status of a compound command other than a subshell command was the result of a failure while
-ewas being ignored, then-eshall not apply to this command.This requirement applies to the shell environment and each subshell environment separately. For example, in:

set -e; (false; echo one) | cat; echo twothe

falsecommand causes the subshell to exit without executingecho one; however,echo twois executed because the exit status of the pipeline(false; echo one) | catis zero.In

set -e; echo $(false; echo one) twothe

falsecommand causes the subshell in which the command substitution is performed to exit without executingecho one; the exit status of the subshell is ignored and the shell then executes the word-expanded commandecho two.
-f- The shell shall disable pathname expansion.
-h
[OB]Setting this option may speed upPATHsearches (see XBD8. Environment Variables). This option may be enabled by default.
-m- This option shall be supported if the implementation supports the User Portability Utilities option. When this option is enabled, the shell shall perform job control actions as described in 2.11 Job Control. This option shall be enabled by default for interactive shells.
-n- The shell shall read commands but does not execute them; this can be used to check for shell script syntax errors. Interactive shells and subshells of interactive shells, recursively, may ignore this option.
-o- Write the current settings of the options to standard output in an unspecified format.
+o- Write the current option settings to standard output in a format that is suitable for reinput to the shell as commands that achieve the same options settings.
-ooption

Set various options, many of which shall be equivalent to the single option letters. The following values ofoptionshall be supported:

allexport- Equivalent to
-a.
errexit- Equivalent to
-e.
ignoreeof- Prevent an interactive shell from exiting on end-of-file. This setting prevents accidental logouts when <control>-D is entered. A user shall explicitly
exitto leave the interactive shell. This option shall be supported if the system supports the User Portability Utilities option.
monitor- Equivalent to
-m. This option shall be supported if the system supports the User Portability Utilities option.
noclobber- Equivalent to
-C(uppercase C).
noglob- Equivalent to
-f.
noexec- Equivalent to
-n.
nolog
[OB]Prevent the entry of function definitions into the command history; seeCommand History List. This option may have no effect; it is kept for compatibility with previous versions of the standard. This option shall be supported if the system supports the User Portability Utilities option.
notify- Equivalent to
-b.
nounset- Equivalent to
-u.
pipefail- Derive the exit status of a pipeline from the exit statuses of all of the commands in the pipeline, not just the last (rightmost) command, as described in 2.9.2 Pipelines.
verbose- Equivalent to
-v.
vi- Allow shell command line editing using the built-in
vieditor. Enablingvimode shall disable any other command line editing mode provided as an implementation extension. This option shall be supported if the system supports the User Portability Utilities option.It need not be possible to set

vimode on for certain block-mode terminals.
xtrace- Equivalent to
-x.
-u- When the shell tries to expand, in a parameter expansion or an arithmetic expansion, an unset parameter other than the
'@'and'*'special parameters, it shall write a message to standard error and the expansion shall fail with the consequences specified in 2.8.1 Consequences of Shell Errors.
-v- The shell shall write its input to standard error as it is read.
-x- The shell shall write to standard error a trace for each command after it expands the command and before it executes it. It is unspecified whether the command that turns tracing off is traced.
The default for all these options shall be off (unset) unless stated otherwise in the description of the option or unless the shell was invoked with them on; see

sh.The remaining arguments shall be assigned in order to the positional parameters. The special parameter

'#'shall be set to reflect the number of positional parameters. All positional parameters shall be unset before any new values are assigned.If the first argument is

'-', the results are unspecified.The special argument

"--"immediately following thesetcommand name can be used to delimit the arguments if the first argument begins with'+'or'-', or to prevent inadvertent listing of all shell variables when there are no arguments. The commandset--withoutargumentshall unset all positional parameters and set the special parameter'#'to zero.

See the DESCRIPTION.

See the DESCRIPTION.

Not used.

None.

None.

Default.

See the DESCRIPTION.

The standard error shall be used only for diagnostic messages.

None.

None.

- 0
- Successful completion.
- >0
- An invalid option was specified, or an error occurred.

Default.

Application writers should avoid relying on

set-ewithin functions. For example, in the following script:set -e start() { some_server echo some_server started successfully } start || echo >&2 some_server failedthe

-esetting is ignored within the function body (because the function is a command in an AND-OR list other than the last). Therefore, ifsome_serverfails, the function carries on to echo"some_server started successfully", and the exit status of the function is zero (which means"some_server failed"is not output).Use of

set-ncauses the shell to parse the rest of the script without executing any commands, meaning thatset+ncannot be used to undo the effect. Syntax checking is more commonly done viash-nscript_name.

Write out all variables and their values:

setSet $1, $2, and $3 and set

"$#"to 3:set c a bTurn on the

-xand-voptions:set -xvUnset all positional parameters:

set --Set $1 to the value of

x, even if it begins with'-'or'+':set -- "$x"Set the positional parameters to the expansion of

x, even ifxexpands with a leading'-'or'+':set -- $x

The

set-- form is listed specifically in the SYNOPSIS even though this usage is implied by the Utility Syntax Guidelines. The explanation of this feature removes any ambiguity about whether theset-- form might be misinterpreted as being equivalent tosetwithout any options or arguments. The functionality of this form has been adopted from the KornShell. In System V,set-- only unsets parameters if there is at least one argument; the only way to unset all parameters is to useshift. Using the KornShell version should not affect System V scripts because there should be no reason to issue it without arguments deliberately; if it were issued as, for example:set -- "$@"and there were in fact no arguments resulting from

"$@", unsetting the parameters would have no result.The

set+ form in early proposals was omitted as being an unnecessary duplication ofsetalone and not widespread historical practice.The

noclobberoption was changed to allowset-Cas well as theset-onoclobberoption. The single-letter version was added so that the historical"$-"paradigm would not be broken; see 2.5.2 Special Parameters.The description of the

-eoption is intended to match the behavior of the 1988 version of the KornShell.The

-hoption is related to command name hashing. Seehash. The normative description is deliberately vague because the way this option works varies between shell implementations.Earlier versions of this standard specified

-has a way to locate and remember utilities to be invoked by functions as those functions are defined (the utilities are normally located when the function is executed). However, this did not match existing practice in most shells.The following

setoptions were omitted intentionally with the following rationale:

-k- The
-koption was originally added by the author of the Bourne shell to make it easier for users of pre-release versions of the shell. In early versions of the Bourne shell the constructsetname=valuehad to be used to assign values to shell variables. The problem with-kis that the behavior affects parsing, virtually precluding writing any compilers. To explain the behavior of-k, it is necessary to describe the parsing algorithm, which is implementation-defined. For example:set -k; echoname=valueand:

set -k echoname=valuebehave differently. The interaction with functions is even more complex. What is more, the

-koption is never needed, since the command line could have been reordered.
-t- The
-toption is hard to specify and almost never used. The only known use could be done with here-documents. Moreover, the behavior withkshandshdiffers. The reference page says that it exits after reading and executing one command. What is one command? If the input isdate;date,shexecutes bothdatecommands whilekshdoes only the first.Consideration was given to rewriting

setto simplify its confusing syntax. A specific suggestion was that theunsetutility should be used to unset options instead of using the non-getopt()-able +optionsyntax. However, the conclusion was reached that the historical practice of using +optionwas satisfactory and that there was no compelling reason to modify such widespread historical practice.The

-ooption was adopted from the KornShell to address user needs. In addition to its generally friendly interface,-ois needed to provide thevicommand line editing mode, for which historical practice yields no single-letter option name. (Although it might have been possible to invent such a letter, it was recognized that other editing modes would be developed and-oprovides ample name space for describing such extensions.)Historical implementations are inconsistent in the format used for

-ooption status reporting. The+oformat without an option-argument was added to allow portable access to the options that can be saved and then later restored using, for instance, a dot script.Historically,

shdid trace the commandset+x, butkshdid not.The

ignoreeofsetting prevents accidental logouts when the end-of-file character (typically <control>-D) is entered. A user shall explicitlyexitto leave the interactive shell.The

set-moption was added to apply only to the UPE because it applies primarily to interactive use, not shell script applications.The ability to do asynchronous notification became available in the 1988 version of the KornShell. To have it occur, the user had to issue the command:

trap "jobs -n" CLDThe C shell provides two different levels of an asynchronous notification capability. The environment variable

notifyis analogous to what is done inset-borset-onotify. When set, it notifies the user immediately of background job completions. When unset, this capability is turned off.The other notification ability comes through the built-in utility

notify. The syntax is:notify[%job ...]By issuing

notifywith no operands, it causes the C shell to notify the user asynchronously when the state of the current job changes. If given operands,notifyasynchronously informs the user of changes in the states of the specified jobs.To add asynchronous notification to the POSIX shell, neither the KornShell extensions to

trap, nor the C shellnotifyenvironment variable seemed appropriate (notifyis not a proper POSIX environment variable name).The

set-boption was selected as a compromise.The

notifybuilt-in was considered to have more functionality than was required for simple asynchronous notification.Historically, some shells applied the

-uoption to all parameters including$@and$*. The standard developers felt that this was a misfeature since it is normal and common for$@and$*to be used in shell scripts regardless of whether they were passed any arguments. Treating these uses as an error when no arguments are passed reduces the value of-ufor its intended purpose of finding spelling mistakes in variable names and uses of unset positional parameters.

A future version of this standard may remove the

-onologoption.

2.15 Special Built-In Utilities,

hashXBD

4.26 Variable Assignment,12.2 Utility Syntax Guidelines

The obsolescent

setcommand name followed by'-'has been removed.The following new requirements on POSIX implementations derive from alignment with the Single UNIX Specification:

- The
nologoption is added toset-o.IEEE PASC Interpretation 1003.2 #167 is applied, clarifying that the options default also takes into account the description of the option.

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/8 is applied, changing the square brackets in the example in RATIONALE to be in bold, which is the typeface used for optional items.

Austin Group Interpretation 1003.1-2001 #027 is applied, clarifying the behavior if the first argument is

'-'.SD5-XCU-ERN-97 is applied, updating the SYNOPSIS.

XSI shading is removed from the

-hfunctionality.POSIX.1-2008, Technical Corrigendum 1, XCU/TC1-2008/0046 [52], XCU/TC1-2008/0047 [155,280], XCU/TC1-2008/0048 [52], XCU/TC1-2008/0049 [52], and XCU/TC1-2008/0050 [155,430] are applied.

POSIX.1-2008, Technical Corrigendum 2, XCU/TC2-2008/0053 [584], XCU/TC2-2008/0054 [717], XCU/TC2-2008/0055 [717], and XCU/TC2-2008/0056 [960] are applied.

Austin Group Defect 559 is applied, changing the description of the

-uoption.Austin Group Defect 789 is applied, adding

-opipefail.Austin Group Defect 981 is applied, changing the description of the

-onologoption and the FUTURE DIRECTIONS section.Austin Group Defects 1009 and 1555 are applied, changing the description of the

-aoption.Austin Group Defect 1016 is applied, changing the description of the

-Coption.Austin Group Defect 1055 is applied, adding a paragraph about the

-noption to the APPLICATION USAGE section.Austin Group Defect 1063 is applied, changing the description of the

-hoption.Austin Group Defect 1150 is applied, changing the description of the

-eoption.Austin Group Defect 1207 is applied, clarifying which option-arguments of the

-ooption are related to the User Portability Utilities option.Austin Group Defect 1254 is applied, changing the descriptions of the

-band-moptions.Austin Group Defect 1384 is applied, allowing subshells of interactive shells to ignore the

-noption.

shift — shift positional parameters

shift[n]

The positional parameters shall be shifted. Positional parameter 1 shall be assigned the value of parameter (1+

n), parameter 2 shall be assigned the value of parameter (2+n), and so on. The parameters represented by the numbers"$#"down to"$#-n+1"shall be unset, and the parameter'#'is updated to reflect the new number of positional parameters.The value

nshall be an unsigned decimal integer less than or equal to the value of the special parameter'#'. Ifnis not given, it shall be assumed to be 1. Ifnis 0, the positional and special parameters are not changed.

None.

See the DESCRIPTION.

Not used.

None.

None.

Default.

Not used.

The standard error shall be used only for diagnostic messages and the warning message specified in EXIT STATUS.

None.

None.

If the

noperand is invalid or is greater than"$#", this may be treated as an error and a non-interactive shell may exit; if the shell does not exit in this case, a non-zero exit status shall be returned and a warning message shall be written to standard error. Otherwise, zero shall be returned.

Default.

None.

$set a b c d e$shift 2$echo $*c d e

None.

None.

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.

POSIX.1-2008, Technical Corrigendum 1, XCU/TC1-2008/0051 [459] is applied.

Austin Group Defect 1265 is applied, updating the EXIT STATUS and STDERR sections to align with the changes made to 2.8.1 Consequences of Shell Errors between Issue 6 and Issue 7.

times — write process times

times

The

timesutility shall write the accumulated user and system times for the shell and for all of its child processes, in the following POSIX locale format:"%dm%fs %dm%fs\n%dm%fs %dm%fs\n", <shell user minutes>, <shell user seconds>, <shell system minutes>, <shell system seconds>, <children user minutes>, <children user seconds>, <children system minutes>, <children system seconds>The four pairs of times shall correspond to the members of the

<sys/times.h>tmsstructure (defined in XBD14. Headers) as returned bytimes():tms_utime,tms_stime,tms_cutime, andtms_cstime, respectively.

None.

None.

Not used.

None.

None.

Default.

See the DESCRIPTION.

The standard error shall be used only for diagnostic messages.

None.

None.

- 0
- Successful completion.
- >0
- An error occurred.

Default.

None.

$times0m0.43s 0m1.11s 8m44.18s 1m43.23s

The

timesspecial built-in from the Single UNIX Specification is now required for all conforming shells.

None.

2.15 Special Built-In Utilities

XBD

<sys/times.h>

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/9 is applied, changing text in the DESCRIPTION from: "Write the accumulated user and system times for the shell and for all of its child processes ..." to: "The

timesutility shall write the accumulated user and system times for the shell and for all of its child processes ...".

POSIX.1-2008, Technical Corrigendum 2, XCU/TC2-2008/0056 [960] is applied.

trap — trap signals

trapn[condition...]

trap -p[condition...]

trap[action condition...]

If the

-poption is not specified and the first operand is an unsigned decimal integer, the shell shall treat all operands as conditions, and shall reset each condition to the default value. Otherwise, if the-poption is not specified and there are operands, the first operand shall be treated as an action and the remaining as conditions.If

actionis'-', the shell shall reset eachconditionto the default value. Ifactionis null (""), the shell shall ignore each specifiedconditionif it arises. Otherwise, the argumentactionshall be read and executed by the shell when one of the corresponding conditions arises. The action oftrapshall override a previous action (either default action or one explicitly set). The value of"$?"after thetrapaction completes shall be the value it had before thetrapaction was executed.The condition can be EXIT, 0 (equivalent to EXIT), or a signal specified using a symbolic name, without the SIG prefix, as listed in the tables of signal names in the

<signal.h>header defined in XBD14. Headers; for example, HUP, INT, QUIT, TERM. Implementations may permit names with the SIG prefix or ignore case in signal names as an extension. Setting a trap for SIGKILL or SIGSTOP produces undefined results.The EXIT condition shall occur when the shell terminates normally (exits), and may occur when the shell terminates abnormally as a result of delivery of a signal (other than SIGKILL) whose

trapaction is the default.The environment in which the shell executes a

trapaction on EXIT shall be identical to the environment immediately after the last command executed before thetrapaction on EXIT was executed.If

actionis neither'-'nor the empty string, then each time a matchingconditionarises, theactionshall be executed in a manner equivalent to:evalactionSignals that were ignored on entry to a non-interactive shell cannot be trapped or reset, although no error need be reported when attempting to do so. An interactive shell may reset or catch signals ignored on entry. Traps shall remain in place for a given shell until explicitly changed with another

trapcommand.When a subshell is entered, traps that are not being ignored shall be set to the default actions, except in the case of a command substitution containing only a single

trapcommand, when the traps need not be altered. Implementations may check for this case using only lexical analysis; for example, if`trap`and$( trap -- )do not alter the traps in the subshell, cases such as assigningvar=trapand then using$($var)may still alter them. This does not imply that thetrapcommand cannot be used within the subshell to set new traps.The

trapcommand with no operands shall write to standard output a list of commands associated with each of a set of conditions; if the-poption is not specified, this set shall contain only the conditions that are not in the default state (including signals that were ignored on entry to a non-interactive shell); if the-poption is specified, the set shall contain all conditions, except that it is unspecified whether conditions corresponding to the SIGKILL and SIGSTOP signals are included in the set. If the command is executed in a subshell, the implementation does not perform the optional check described above for a command substitution containing only a singletrapcommand, and notrapcommands with operands have been executed since entry to the subshell, the list shall contain the commands that were associated with each condition immediately before the subshell environment was entered. Otherwise, the list shall contain the commands currently associated with each condition. The format shall be:"trap -- %s %s ...\n", <action>, <condition> ...The shell shall format the output, including the proper use of quoting, so that it is suitable for reinput to the shell as commands that achieve the same trapping results for the set of conditions included in the output, except for signals that were ignored on entry to the shell as described above. If this set includes conditions corresponding to the SIGKILL and SIGSTOP signals, the shell shall accept them when the output is reinput to the shell (where accepting them means they do not cause a non-zero exit status, a diagnostic message, or undefined behavior). For example:

save_traps=$(trap -p)

... eval "$save_traps"or:

save_traps=$(trap -p INT QUIT) trap "some command" INT QUIT

...

eval "$save_traps"

[XSI]XSI-conformant systems also allow numeric signal numbers for the conditions corresponding to the following signal names:

- 1
- SIGHUP
- 2
- SIGINT
- 3
- SIGQUIT
- 6
- SIGABRT
- 9
- SIGKILL
- 14
- SIGALRM
- 15
- SIGTERM
If an invalid signal name

[XSI]or number is specified, thetraputility shall write a warning message to standard error.The

trapspecial built-in shall conform to XBD12.2 Utility Syntax Guidelines.

The following option shall be supported:

-p- Write to standard output a list of commands associated with each
conditionoperand. The behavior when there are no operands is specified in the DESCRIPTION section.The shell shall format the output, including the proper use of quoting, so that it is suitable for reinput to the shell as commands that achieve the same trapping results for the specified set of conditions. If a

conditionoperand is a condition corresponding to the SIGKILL or SIGSTOP signal, andtrap-pwithout any operands would not include it in the set of conditions for which it writes output, the behavior is undefined if the output is reinput to the shell.

See the DESCRIPTION.

Not used.

None.

None.

Default.

See the DESCRIPTION.

The standard error shall be used only for diagnostic messages and warning messages about invalid signal names

[XSI]or numbers.

None.

None.

If the trap name

[XSI]or number is invalid, a non-zero exit status shall be returned; otherwise, zero shall be returned. For both interactive and non-interactive shells, invalid signal names[XSI]or numbers shall not be considered an error and shall not cause the shell to abort.

Default.

When the

-poption is not used, sincetrapwith no operands does not output commands to restore traps that are currently set to default, these need to be restored separately. The RATIONALE section shows examples and describes their drawbacks.

Write out a list of all traps and actions:

trapSet a trap so the

logoututility in the directory referred to by theHOMEenvironment variable executes when the shell terminates:trap '"$HOME"/logout' EXITor:

trap '"$HOME"/logout' 0Unset traps on INT, QUIT, TERM, and EXIT:

trap - INT QUIT TERM EXIT

Implementations may permit lowercase signal names as an extension. Implementations may also accept the names with the SIG prefix; no known historical shell does so. The

trapandkillutilities in this volume of POSIX.1-2024 are now consistent in their omission of the SIG prefix for signal names. Somekillimplementations do not allow the prefix, andkill-llists the signals without prefixes.Trapping SIGKILL or SIGSTOP is syntactically accepted by some historical implementations, but it has no effect. Portable POSIX applications cannot attempt to trap these signals.

The output format is not historical practice. Since the output of historical

trapcommands is not portable (because numeric signal values are not portable) and had to change to become so, an opportunity was taken to format the output in a way that a shell script could use to save and then later reuse a trap if it wanted.The KornShell uses an

ERRtrap that is triggered wheneverset-ewould cause an exit. This is allowable as an extension, but was not mandated, as other shells have not used it.The text about the environment for the EXIT trap invalidates the behavior of some historical versions of interactive shells which, for example, close the standard input before executing a trap on 0. For example, in some historical interactive shell sessions the following trap on 0 would always print

"--":trap 'read foo; echo "-$foo-"' 0The command:

trap 'eval " $cmd"' 0causes the contents of the shell variable

cmdto be executed as a command when the shell exits. Using:trap '$cmd' 0does not work correctly if

cmdcontains any special characters such as quoting or redirections. Using:trap " $cmd" 0also works (the leading <space> character protects against unlikely cases where

cmdis a decimal integer or begins with'-'), but it expands thecmdvariable when thetrapcommand is executed, not when the exit action is executed.The

-poption was added because without it the method used to restore traps needs to include special handling of traps that are set to default whentrapwith no operands is used to save the current traps. One example is:save_traps=$(trap) trap "some command" INT QUIT save_traps="trap - INT QUIT; $save_traps"

...

eval "$save_traps"but this method relies on hard-coding the commands to reset the traps that are being set. It also has a race condition if INT or QUIT was not set to default when saved, since it first sets them to default and then restores the saved traps. A more general approach would be:

save_traps=$(trap) ... for sig in EXIT $( kill -l ) do case "$sig" in SIGKILL | KILL | sigkill | kill | SIGSTOP | STOP | sigstop | stop) ;; *) trap - $sig ;; esac done eval "$save_traps"This has the same race condition since it first sets all traps (that can be set) to default and then restores those that were not previously set to default.

Historically, some shells behaved the same with and without

-pwhen there are no operands. This standard requires that the set of conditions differs between the two cases: with-pit is all conditions (except possibly SIGKILL and SIGSTOP); without-pit is only the conditions that are not in the default state.

None.

XSI-conforming implementations provide the mapping of signal names to numbers given above (previously this had been marked obsolescent). Other implementations need not provide this optional mapping.

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.

SD5-XCU-ERN-97 is applied, updating the SYNOPSIS.

Austin Group Interpretation 1003.1-2001 #116 is applied.

POSIX.1-2008, Technical Corrigendum 1, XCU/TC1-2008/0052 [53,268,440], XCU/TC1-2008/0053 [53,268,440], XCU/TC1-2008/0054 [163], XCU/TC1-2008/0055 [163], and XCU/TC1-2008/0056 [163] are applied.

Austin Group Defect 621 is applied, clarifying when the EXIT condition occurs.

Austin Group Defect 1029 is applied, clarifying the execution of

trapactions.Austin Group Defects 1211 and 1212 are applied, adding the

-poption and clarifying that, when-pis not specified, the output oftrapwith no operands does not list conditions that are in the default state.Austin Group Defect 1265 is applied, updating the DESCRIPTION, STDERR and EXIT STATUS sections to align with the changes made to 2.8.1 Consequences of Shell Errors between Issue 6 and Issue 7.

Austin Group Defect 1285 is applied, inserting a blank line between the two SYNOPSIS lines.

unset — unset values and attributes of variables and functions

unset[-fv]name...

The

unsetutility shall unset each variable or function definition specified bynamethat does not have thereadonlyattribute and remove any attributes other thanreadonlythat have been given toname(see 2.15 Special Built-In Utilitiesexportandreadonly).If

-vis specified,namerefers to a variable name and the shell shall unset it and remove it from the environment. Read-only variables cannot be unset.If

-fis specified,namerefers to a function and the shell shall unset the function definition.If neither

-fnor-vis specified,namerefers to a variable; if a variable by that name does not exist, it is unspecified whether a function by that name, if any, shall be unset.Unsetting a variable or function that was not previously set shall not be considered an error and does not cause the shell to abort.

The

unsetspecial built-in shall support XBD12.2 Utility Syntax Guidelines.Note that:

VARIABLE=is not equivalent to an

unsetofVARIABLE; in the example,VARIABLEis set to"". Also, the variables that can beunsetshould not be misinterpreted to include the special parameters (see 2.5.2 Special Parameters).

See the DESCRIPTION.

See the DESCRIPTION.

Not used.

None.

None.

Default.

Not used.

The standard error shall be used only for diagnostic messages.

None.

None.

- 0
- All
nameoperands were successfully unset.- >0
- At least one
namecould not be unset.

Default.

None.

Unset

VISUALvariable:unset -v VISUALUnset the functions

fooandbar:unset -f foo bar

Consideration was given to omitting the

-foption in favor of anunfunctionutility, but the standard developers decided to retain historical practice.The

-voption was introduced because System V historically used one name space for both variables and functions. Whenunsetis used without options, System V historically unset either a function or a variable, and there was no confusion about which one was intended. A portable POSIX application can useunsetwithout an option to unset a variable, but not a function; the-foption must be used.

None.

IEEE Std 1003.1-2001/Cor 1-2002, item XCU/TC1/D6/5 is applied so that the reference page sections use terms as described in the Utility Description Defaults (

1.4 Utility Description Defaults). No change in behavior is intended.

SD5-XCU-ERN-97 is applied, updating the SYNOPSIS.

Austin Group Defect 1075 is applied, clarifying that

unsetremoves attributes, other thanreadonly, from the variables it unsets.
