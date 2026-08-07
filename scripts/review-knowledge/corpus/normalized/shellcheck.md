# ShellCheck Wiki Sitemap

The wiki is maintained on GitHub. This page is primarily for the benefit of search engines.
- SC1000 – `$` is not used specially and should therefore be escaped.
- SC1001 – This `\o` will be a regular 'o' in this context.
- SC1003 – Want to escape a single quote? `echo 'This is how it'\''s done'`.
- SC1004 – This backslash+linefeed is literal. Break outside single quotes if you just want to break the line.
- SC1007 – Remove space after `=` if trying to assign a value (or for empty string, use `var=""` ... ).
- SC1008 – This shebang was unrecognized. ShellCheck only supports sh/bash/dash/ksh. Add a 'shell' directive to specify.
- SC1009 – The mentioned parser error was in ...
- SC1010 – Use semicolon or linefeed before `done` (or quote to make it literal).
- SC1011 – This apostrophe terminated the single quoted string!
- SC1012 – `\t` is just literal `t` here. For tab, use `"$(printf '\t')"` instead.
- SC1014 – Use `if cmd; then ..` to check exit code, or `if [ "$(cmd)" = .. ]` to check output.
- SC1015 – This is a Unicode double quote. Delete and retype it.
- SC1016 – This is a Unicode single quote. Delete and retype it.
- SC1017 – Literal carriage return. Run script through `tr -d '\r'` .
- SC1018 – This is a Unicode non-breaking space. Delete it and retype as space.
- SC1019 – Expected this to be an argument to the unary condition.
- SC1020 – You need a space before the `]` or `]]`
- SC1026 – If grouping expressions inside `[[..]]`, use `( .. )`.
- SC1027 – Expected another argument for this operator.
- SC1028 – In `[..]` you have to escape `\( \)` or preferably combine `[..]` expressions.
- SC1029 – In `[[..]]` you shouldn't escape `(` or `)`.
- SC1033 – Test expression was opened with double `[[` but closed with single `]`. Make sure they match.
- SC1034 – Test expression was opened with single `[` but closed with double `]]`. Make sure they match.
- SC1035 – You need a space here
- SC1036 – `(` is invalid here. Did you forget to escape it?
- SC1037 – Braces are required for positionals over 9, e.g. `${10}`.
- SC1038 – Shells are space sensitive. Use `< <(cmd)`, not `<<(cmd)`.
- SC1039 – Remove indentation before end token (or use `<<-` and indent with tabs).
- SC1040 – When using `<<-`, you can only indent with tabs.
- SC1041 – Found `eof` further down, but not on a separate line.
- SC1043 – Found EOF further down, but with wrong casing.
- SC1044 – Couldn't find end token `EOF` in the here document.
- SC1045 – It's not `foo &; bar`, just `foo & bar`.
- SC1046 – Couldn't find `fi` for this `if`.
- SC1047 – Expected `fi` matching previously mentioned `if`.
- SC1048 – Can't have empty then clauses (use `true` as a no-op).
- SC1049 – Did you forget the `then` for this `if`?
- SC1050 – Expected `then`.
- SC1051 – Semicolons directly after `then` are not allowed. Just remove it.
- SC1052 – Semicolons directly after `then` are not allowed. Just remove it.
- SC1053 – Semicolons directly after `else` are not allowed. Just remove it.
- SC1054 – You need a space after the `{`.
- SC1055 – You need at least one command here. Use `true;` as a no-op.
- SC1056 – Expected a `}`. If you have one, try a `;` or `\n` in front of it.
- SC1057 – Did you forget the `do` for this loop?
- SC1058 – Expected `do`.
- SC1059 – Semicolon is not allowed directly after `do`. You can just delete it.
- SC1060 – Can't have empty do clauses (use `true` as a no-op)
- SC1061 – Couldn't find `done` for this `do`.
- SC1062 – Expected `done` matching previously mentioned `do`.
- SC1063 – You need a line feed or semicolon before the `do`.
- SC1064 – Expected a `{` to open the function definition.
- SC1065 – Trying to declare parameters? Don't. Use `()` and refer to params as `$1`, `$2`, …
- SC1066 – Don't use `$` on the left side of assignments.
- SC1067 – For indirection, use arrays, `declare "var$n=value"`, or (for sh) `read`/`eval`
- SC1068 – Don't put spaces around the `=` in assignments.
- SC1069 – You need a space before the `[`.
- SC1070 – Parsing stopped here. Mismatched keywords or invalid parentheses?
- SC1071 – ShellCheck only supports sh/bash/dash/ksh scripts. Sorry!
- SC1072 – Unexpected ..
- SC1073 – Couldn't parse this (thing). Fix to allow more checks.
- SC1074 – Did you forget the `;;` after the previous case item?
- SC1075 – Use `elif` instead of `else if`.
- SC1076 – Trying to do math? Use e.g. `[ $((i/2+7)) -ge 18 ]`.
- SC1077 – For command expansion, the tick should slant left (`` ` `` vs `´`).
- SC1078 – Did you forget to close this double-quoted string?
- SC1079 – This is actually an end quote, but due to next char it looks suspect.
- SC1080 – You need `\` before line feeds to break lines in `[ ]`.
- SC1081 – Scripts are case-sensitive. Use `if`, not `If`.
- SC1082 – This file has a UTF-8 BOM. Remove it with: `LC_CTYPE=C sed '1s/^...//' < yourscript`.
- SC1083 – This `{`/`}` is literal. Check if `;` is missing or quote the expression.
- SC1084 – Use `#!`, not `!#`, for the shebang.
- SC1085 – Did you forget to move the ;; after extending this case item?
- SC1086 – Don't use `$` on the iterator name in for loops.
- SC1087 – Use braces when expanding arrays, e.g. `${array[idx]}` (or `${var}[..` to quiet).
- SC1088 – Parsing stopped here. Invalid use of parentheses?
- SC1089 – Parsing stopped here. Is this keyword correctly matched up?
- SC1090 – Can't follow non-constant source. Use a directive to specify location
- SC1091 – Not following: (error message here)
- SC1092 – Stopping at 100 `source` frames :O
- SC1094 – Parsing of sourced file failed. Ignoring it.
- SC1095 – You need a space or linefeed between the function name and body.
- SC1097 – Unexpected `==`. For assignment, use `=`. For comparison, use `[`/`[[`.
- SC1098 – Quote/escape special characters when using `eval`, e.g. `eval "a=(b)"`.
- SC1099 – You need a space before the `#`.
- SC1100 – This is a Unicode dash. Delete and retype as ASCII minus.
- SC1101 – Delete trailing spaces after `\` to break line (or use quotes for literal space).
- SC1102 – Shells disambiguate `$((` differently or not at all. For `$(command substitution)`, add space after `$(` . For `$((arithmetics))`, fix parsing errors.
- SC1103 – This shell type is unknown. Use e.g. `sh` or `bash`.
- SC1104 – Use `#!`, not just `!`, for the shebang.
- SC1105 – Shells disambiguate `((` differently or not at all. If the first `(` should start a subshell, add a space after it.
- SC1106 – In arithmetic contexts, use `<` instead of `-lt`
- SC1107 – This directive is unknown. It will be ignored.
- SC1108 – You need a space before and after the `=` .
- SC1109 – This is an unquoted HTML entity. Replace with corresponding character.
- SC1110 – This is a Unicode quote. Delete and retype it (or quote to make literal).
- SC1111 – This is a Unicode quote. Delete and retype it (or ignore/singlequote for literal).
- SC1112 – This is a Unicode quote. Delete and retype it (or ignore/doublequote for literal).
- SC1113 – Use `#!`, not just `#`, for the shebang.
- SC1114 – Remove leading spaces before the shebang.
- SC1115 – Remove spaces between `#` and `!` in the shebang.
- SC1116 – Missing `$` on a `$((..))` expression? (or use `( (` for arrays).
- SC1117 – Backslash is literal in `"\n"`. Prefer explicit escaping: `"\\n"`.
- SC1118 – Delete whitespace after the here-doc end token.
- SC1119 – Add a linefeed between end token and terminating `)`.
- SC1120 – No comments allowed after here-doc token. Comment the next line instead.
- SC1121 – Add `;`/`&` terminators (and other syntax) on the line with the `<<`, not here.
- SC1122 – Nothing allowed after end token. To continue a command, put it on the line with the `<<`.
- SC1123 – ShellCheck directives are only valid in front of complete compound commands, like `if`, not e.g. individual `elif` branches.
- SC1124 – ShellCheck directives are only valid in front of complete commands like `case` statements, not individual case branches.
- SC1125 – Invalid `key=value` pair in directive
- SC1126 – Place shellcheck directives before commands, not after.
- SC1127 – Was this intended as a comment? Use `#` in sh.
- SC1128 – The shebang must be on the first line. Delete blanks and move comments.
- SC1129 – You need a space before the `!`.
- SC1130 – You need a space before the :.
- SC1131 – Use `elif` to start another branch.
- SC1132 – This `&` terminates the command. Escape it or add space after `&` to silence.
- SC1133 – Unexpected start of line. If breaking lines, `|`/`||`/`&&` should be at the end of the previous one.
- SC1134 – Error parsing `shellcheckrc`
- SC1135 – Prefer escape over ending quote to make `$` literal. Instead of `"It costs $"5`, use `"It costs \$5"`
- SC1136 – Unexpected characters after terminating `]`. Missing semicolon/linefeed?
- SC1137 – Missing second `(` to start arithmetic for ((;;)) loop
- SC1138 – Remove spaces between (( in arithmetic for loop.
- SC1139 – Use `||` instead of `-o` between test commands.
- SC1140 – Unexpected parameters after condition. Missing `&&`/`||`, or bad expression?
- SC1141 – Unexpected tokens after compound command. Bad redirection or missing `;`/`&&`/`||`/`|`?
- SC1142 – Use `done < <(cmd)` to redirect from process substitution (currently missing one `<`).
- SC1143 – This backslash is part of a comment and does not continue the line.
- SC1144 – `external-sources` can only be enabled in .shellcheckrc, not in individual files.
- SC1145 – Unknown `external-sources` value. Expected `true`/`false`.
- SC2000 – See if you can use `${#variable}` instead
- SC2001 – See if you can use `${variable//search/replace}` instead.
- SC2002 – Useless cat. Consider `cmd < file | ..` or `cmd file | ..` instead.
- SC2003 – expr is antiquated. Consider rewriting this using `$((..))`, `${}` or `[[ ]]`.
- SC2004 – `$`/`${}` is unnecessary on arithmetic variables.
- SC2005 – Useless `echo`? Instead of `echo $(cmd)`, just use `cmd`
- SC2006 – Use `$(...)` notation instead of legacy backticked `` `...` ``.
- SC2007 – Use `$((..))` instead of deprecated `$[..]`.
- SC2008 – `echo` doesn't read from stdin, are you sure you should be piping to it?
- SC2009 – Consider using `pgrep` instead of grepping `ps` output.
- SC2010 – Don't use `ls | grep`. Use a glob or a for loop with a condition to allow non-alphanumeric filenames.
- SC2011 – Use `find -print0` or `find -exec` to better handle non-alphanumeric filenames.
- SC2012 – Use `find` instead of `ls` to better handle non-alphanumeric filenames.
- SC2013 – To read lines rather than words, pipe/redirect to a `while read` loop.
- SC2014 – This will expand once before find runs, not per file found.
- SC2015 – Note that `A && B || C` is not if-then-else. C may run when A is true.
- SC2016 – Expressions don't expand in single quotes, use double quotes for that.
- SC2017 – Increase precision by replacing `a/b*c` with `a*c/b`.
- SC2018 – Use `[:lower:]` to support accents and foreign alphabets.
- SC2019 – Use `[:upper:]` to support accents and foreign alphabets.
- SC2020 – `tr` replaces sets of chars, not words (mentioned due to duplicates).
- SC2021 – Don't use `[]` around ranges in `tr`, it replaces literal square brackets.
- SC2022 – Note that unlike globs, `o*` here matches `ooo` but not `oscar`.
- SC2023 – The shell may override `time` as seen in man time(1). Use `command time ..` for that one.
- SC2024 – `sudo` doesn't affect redirects. Use `..| sudo tee file`
- SC2025 – Make sure all escape sequences are enclosed in `\[..\]` to prevent line wrapping issues.
- SC2026 – This word is outside of quotes. Did you intend to `'nest '"'single quotes'"'` instead?
- SC2027 – The surrounding quotes actually unquote this. Remove or escape them.
- SC2028 – `echo` won't expand escape sequences. Consider `printf`.
- SC2029 – Note that, unescaped, this expands on the client side.
- SC2030 – Modification of var is local (to subshell caused by pipeline).
- SC2031 – var was modified in a subshell. That change might be lost.
- SC2032 – This function can't be invoked via su on line 42.
- SC2033 – Shell functions can't be passed to external commands. Use separate script or sh -c.
- SC2034 – foo appears unused. Verify it or export it.
- SC2035 – Use `./*glob*` or `-- *glob*` so names with dashes won't become options.
- SC2036 – If you wanted to assign the output of the pipeline, use `a=$(b | c)` .
- SC2037 – To assign the output of a command, use `var=$(cmd)` .
- SC2038 – Use `-print0`/`-0` or `find -exec +` to allow for non-alphanumeric filenames.
- SC2039 – In POSIX sh, *something* is undefined.
- SC2040 – `#!/bin/sh` was specified, so ____ is not supported, even when sh is actually bash.
- SC2041 – This is a literal string. To run as a command, use `$(..)` instead of `'..'` .
- SC2042 – Use spaces, not commas, to separate loop elements.
- SC2043 – This loop will only ever run once for a constant value. Did you perhaps mean to loop over `dir/*`, `$var` or `$(cmd)`?
- SC2044 – For loops over find output are fragile. Use `find -exec` or a `while read` loop.
- SC2045 – Iterating over ls output is fragile. Use globs.
- SC2046 – Quote this to prevent word splitting.
- SC2048 – Use `"$@"` (with quotes) to prevent whitespace problems.
- SC2049 – `=~` is for regex, but this looks like a glob. Use `=` instead.
- SC2050 – This expression is constant. Did you forget the `$` on a variable?
- SC2051 – Bash doesn't support variables in brace range expansions.
- SC2053 – Quote the rhs of `=` in `[[ ]]` to prevent glob matching.
- SC2054 – Use spaces, not commas, to separate array elements.
- SC2055 – You probably wanted `&&` here, otherwise it's always true.
- SC2056 – You probably wanted `&&` here
- SC2057 – Unknown binary operator.
- SC2058 – Unknown unary operator.
- SC2059 – Don't use variables in the `printf` format string. Use `printf "..%s.." "$foo"`.
- SC2060 – Quote parameters to `tr` to prevent glob expansion.
- SC2061 – Quote the parameter to `-name` so the shell won't interpret it.
- SC2062 – Quote the grep pattern so the shell won't interpret it.
- SC2063 – Grep uses regex, but this looks like a glob.
- SC2064 – Use single quotes, otherwise this expands now rather than when signalled.
- SC2065 – This is interpreted as a shell file redirection, not a comparison.
- SC2066 – Since you double-quoted this, it will not word split, and the loop will only run once.
- SC2067 – Missing `;` or `+` terminating `-exec`. You can't use `|`/`||`/`&&`, and `;` has to be a separate, quoted argument.
- SC2068 – Double quote array expansions to avoid re-splitting elements.
- SC2069 – To redirect stdout+stderr, `2>&1` must be last (or use `{ cmd > file; } 2>&1` to clarify).
- SC2070 – `-n` doesn't work with unquoted arguments. Quote or use `[[ ]]`.
- SC2071 – `>` is for string comparisons. Use `-gt` instead.
- SC2072 – Decimals are not supported. Either use integers only, or use `bc` or `awk` to compare.
- SC2073 – Escape `\<` to prevent it redirecting (or switch to `[[ .. ]]`).
- SC2074 – Can't use `=~` in `[ ]`. Use `[[..]]` instead.
- SC2075 – Escaping `\<` is required in `[..]`, but invalid in `[[..]]`
- SC2076 – Don't quote rhs of `=~`, it'll match literally rather than as a regex.
- SC2077 – You need spaces around the comparison operator.
- SC2078 – This expression is constant. Did you forget a `$` somewhere?
- SC2079 – `(( ))` doesn't support decimals. Use `bc` or `awk`.
- SC2080 – Numbers with leading 0 are considered octal.
- SC2081 – `[ .. ]` can't match globs. Use `[[ .. ]]` or grep.
- SC2082 – To expand via indirection, use `name="foo$n"; echo "${!name}"`.
- SC2083 – Don't add spaces after the slash in `./file`.
- SC2084 – Remove `$` or use `_=$((expr))` to avoid executing output.
- SC2086 – Double quote to prevent globbing and word splitting.
- SC2087 – Quote `EOF` to make here document expansions happen on the server side rather than on the client.
- SC2088 – Tilde does not expand in quotes. Use `$HOME`.
- SC2089 – Quotes/backslashes will be treated literally. Use an array.
- SC2090 – Quotes/backslashes in this variable will not be respected.
- SC2091 – Remove surrounding `$()` to avoid executing output (or use `eval` if intentional).
- SC2092 – Remove backticks to avoid executing output.
- SC2093 – Remove `exec ` if script should continue after this command.
- SC2094 – Make sure not to read and write the same file in the same pipeline.
- SC2095 – Use `ssh -n` to prevent ssh from swallowing stdin.
- SC2096 – On most OS, shebangs can only specify a single parameter.
- SC2097 – This assignment is only seen by the forked process.
- SC2098 – This expansion will not see the mentioned assignment.
- SC2099 – Use `$((..))` for arithmetics, e.g. `i=$((i + 2))`
- SC2100 – Use `$((..))` for arithmetics, e.g. `i=$((i + 2))`
- SC2101 – Named class needs outer `[]`, e.g. `[[:digit:]]`.
- SC2102 – Ranges can only match single chars (mentioned due to duplicates).
- SC2103 – Use a `( subshell )` to avoid having to `cd` back.
- SC2104 – In functions, use `return` instead of `break`.
- SC2105 – `break` is only valid in loops
- SC2106 – This only exits the subshell caused by the pipeline.
- SC2107 – Instead of `[ a && b ]`, use `[ a ] && [ b ]`.
- SC2108 – In `[[..]]`, use `&&` instead of `-a`.
- SC2109 – Instead of `[ a || b ]`, use `[ a ] || [ b ]`.
- SC2110 – In `[[..]]`, use `||` instead of `-o`.
- SC2111 – ksh does not allow `function` keyword and `()` at the same time.
- SC2112 – `function` keyword is non-standard. Delete it.
- SC2113 – `function` keyword is non-standard. Use `foo()` instead of `function foo`.
- SC2114 – Warning: deletes a system directory.
- SC2115 – Use `"${var:?}"` to ensure this never expands to `/*` .
- SC2116 – Useless echo? Instead of `cmd $(echo foo)`, just use `cmd foo`.
- SC2117 – To run commands as another user, use `su -c` or `sudo`.
- SC2118 – Ksh does not support `|&`. Use `2>&1 |`
- SC2119 – Use `foo "$@"` if function's `$1` should mean script's `$1`.
- SC2120 – foo references arguments, but none are ever passed.
- SC2121 – To assign a variable, use just `var=value`, not `set ..`.
- SC2122 – `>=` is not a valid operator. Use `! a < b` instead.
- SC2123 – `PATH` is the shell search path. Use another name.
- SC2124 – Assigning an array to a string! Assign as array, or use `*` instead of `@` to concatenate.
- SC2125 – Brace expansions and globs are literal in assignments. Quote it or use an array.
- SC2126 – Consider using `grep -c` instead of `grep | wc`
- SC2127 – To use `${ ..; }`, specify `#!/usr/bin/env ksh`.
- SC2128 – Expanding an array without an index only gives the element in the index 0.
- SC2129 – Consider using `{ cmd1; cmd2; } >> file` instead of individual redirects.
- SC2130 – `-eq` is for integer comparisons. Use `=` instead.
- SC2139 – This expands when defined, not when used. Consider escaping.
- SC2140 – Word is of the form `"A"B"C"` (B indicated). Did you mean `"ABC"` or `"A\"B\"C"`?
- SC2141 – Did you mean `IFS=$'\t'` ?
- SC2142 – Aliases can't use positional parameters. Use a function.
- SC2143 – Use `grep -q` instead of comparing output with `[ -n .. ]`.
- SC2144 – `-e` doesn't work with globs. Use a `for` loop.
- SC2145 – Argument mixes string and array. Use `*` or separate argument.
- SC2146 – This action ignores everything before the `-o`. Use `\( \)` to group.
- SC2147 – Literal tilde in PATH works poorly across programs.
- SC2148 – Tips depend on target shell and yours is unknown. Add a shebang.
- SC2149 – Remove `$`/`${}` for numeric index, or escape it for string.
- SC2150 – `-exec` does not automatically invoke a shell. Use `-exec sh -c ..` for that.
- SC2151 – Only one integer 0-255 can be returned. Use stdout for other data.
- SC2152 – Can only return 0-255. Other data should be written to stdout.
- SC2153 – Possible Misspelling: MYVARIABLE may not be assigned. Did you mean MY_VARIABLE?
- SC2154 – var is referenced but not assigned.
- SC2155 – Declare and assign separately to avoid masking return values.
- SC2156 – Injecting filenames is fragile and insecure. Use parameters.
- SC2157 – Argument to implicit `-n` is always true due to literal strings.
- SC2158 – `[ false ]` is true. Remove the brackets
- SC2159 – `[ 0 ]` is true. Use `false` instead.
- SC2160 – Instead of `[ true ]`, just use `true`.
- SC2161 – Instead of `[ 1 ]`, use `true`.
- SC2162 – `read` without `-r` will mangle backslashes.
- SC2163 – This does not export `FOO`. Remove `$`/`${}` for that, or use `${var?}` to quiet.
- SC2164 – Use `cd ... || exit` in case `cd` fails.
- SC2165 – This nested loop overrides the index variable of its parent.
- SC2166 – Prefer `[ p ] && [ q ]` as `[ p -a q ]` is not well-defined.
- SC2167 – This parent loop has its index variable overridden.
- SC2168 – `local` is only valid in functions.
- SC2169 – In dash, *something* is not supported.
- SC2170 – Invalid number for `-eq`. Use `=` to compare as string (or use `$var` to expand as a variable).
- SC2171 – Found trailing `]` outside test. Add missing `[` or quote if intentional.
- SC2172 – Trapping signals by number is not well-defined. Prefer signal names.
- SC2173 – SIGKILL/SIGSTOP can not be trapped.
- SC2174 – When used with `-p`, `-m` only applies to the deepest directory.
- SC2175 – Quote this invalid brace expansion since it should be passed literally to eval
- SC2176 – `time` is undefined for pipelines. time single stage or `bash -c` instead.
- SC2177 – `time` is undefined for compound commands, use `time sh -c` instead.
- SC2178 – Variable was used as an array but is now assigned a string.
- SC2179 – Use `array+=("item")` to append items to an array.
- SC2180 – Bash does not support multidimensional arrays. Use 1D or associative arrays.
- SC2181 – Check exit code directly with e.g. `if mycmd;`, not indirectly with `$?`.
- SC2182 – This printf format string has no variables. Other arguments are ignored.
- SC2183 – This format string has 2 variables, but is passed 1 argument.
- SC2184 – Quote arguments to unset so they're not glob expanded.
- SC2185 – Some finds don't have a default path. Specify `.` explicitly.
- SC2186 – tempfile is deprecated. Use mktemp instead.
- SC2187 – Ash scripts will be checked as Dash. Add `# shellcheck shell=dash` to silence.
- SC2188 – This redirection doesn't have a command. Move to its command (or use `true` as no-op).
- SC2189 – You can't have `|` between this redirection and the command it should apply to.
- SC2190 – Elements in associative arrays need index, e.g. `array=( [index]=value )` .
- SC2191 – The `=` here is literal. To assign by index, use `( [index]=value )` with no spaces. To keep as literal, quote it.
- SC2192 – This array element has no value. Remove spaces after `=` or use `""` for empty string.
- SC2193 – The arguments to this comparison can never be equal. Make sure your syntax is correct.
- SC2194 – This word is constant. Did you forget the `$` on a variable?
- SC2195 – This pattern will never match the case statement's word. Double check them.
- SC2196 – `egrep` is non-standard and deprecated. Use `grep -E` instead.
- SC2197 – `fgrep` is non-standard and deprecated. Use `grep -F` instead.
- SC2198 – Arrays don't work as operands in `[ ]`. Use a loop (or concatenate with `*` instead of `@`).
- SC2199 – Arrays implicitly concatenate in `[[ ]]`. Use a loop (or explicit `*` instead of `@`).
- SC2200 – Brace expansions don't work as operands in `[ ]`. Use a loop.
- SC2201 – Brace expansion doesn't happen in `[[ ]]`. Use a loop.
- SC2202 – Globs don't work as operands in `[ ]`. Use a loop.
- SC2203 – Globs are ignored in `[[ ]]` except right of `=`/`!=`. Use a loop.
- SC2204 – `(..)` is a subshell. Did you mean `[ .. ]`, a test expression?
- SC2205 – `(..)` is a subshell. Did you mean `[ .. ]`, a test expression?
- SC2206 – Quote to prevent word splitting/globbing, or split robustly with mapfile or `read -a`.
- SC2207 – Prefer `mapfile` or `read -a` to split command output (or quote to avoid splitting).
- SC2208 – Use `[[ ]]` or quote arguments to `-v` to avoid glob expansion.
- SC2209 – Use `var=$(command)` to assign output (or quote to assign string).
- SC2210 – This is a file redirection. Was it supposed to be a comparison or fd operation?
- SC2211 – This is a glob used as a command name. Was it supposed to be in `${..}`, array, or is it missing quoting?
- SC2212 – Use `false` instead of empty `[`/`[[` conditionals.
- SC2213 – getopts specified `-n`, but it's not handled by this `case`.
- SC2214 – This case is not specified by getopts.
- SC2215 – This flag is used as a command name. Bad line break or missing `[ .. ]`?
- SC2216 – Piping to `rm`, a command that doesn't read stdin. Wrong command or missing `xargs`?
- SC2217 – Redirecting to `echo`, a command that doesn't read stdin. Bad quoting or missing `xargs`?
- SC2218 – This function is only defined later. Move the definition up.
- SC2219 – Instead of `let expr`, prefer `(( expr ))` .
- SC2220 – Invalid flags are not handled. Add a `*)` case.
- SC2221 – This pattern always overrides a later one.
- SC2222 – This pattern never matches because of a previous pattern.
- SC2223 – This default assignment may cause DoS due to globbing. Quote it.
- SC2224 – This `mv` has no destination. Check the arguments.
- SC2225 – This `cp` has no destination. Check the arguments.
- SC2226 – This `ln` has no destination. Check the arguments, or specify `.` explicitly.
- SC2227 – Redirection applies to the find command itself. Rewrite to work per action (or move to end).
- SC2229 – This does not read `foo`. Remove `$`/`${}` for that, or use `${var?}` to quiet.
- SC2230 – `which` is non-standard. Use builtin `command -v` instead.
- SC2231 – Quote expansions in this `for` loop glob to prevent word splitting, e.g. `"${dir}"/*.txt`.
- SC2232 – Can't use `sudo` with builtins like `cd`. Did you want `sudo sh -c ..` instead?
- SC2233 – Remove superfluous `(..)` around condition to avoid subshell overhead.
- SC2234 – Remove superfluous `(..)` around test command to avoid subshell overhead.
- SC2235 – Use `{ ..; }` instead of `(..)` to avoid subshell overhead.
- SC2236 – Use `-n` instead of `! -z`.
- SC2237 – Use `[ -n .. ]` instead of `! [ -z .. ]`.
- SC2238 – Redirecting to/from command name instead of file. Did you want pipes/xargs (or quote to ignore)?
- SC2239 – Ensure the shebang uses the absolute path to the interpreter.
- SC2240 – The dot command does not support arguments in sh/dash. Set them as variables.
- SC2241 – The exit status can only be one integer 0-255. Use stdout for other data.
- SC2242 – Can only exit with status 0-255. Other data should be written to stdout/stderr.
- SC2243 – Prefer explicit `-n` to check for output (or run command without `[`/`[[` to check for success)
- SC2244 – Prefer explicit `-n` to check non-empty string (or use `=`/`-ne` to check boolean/integer).
- SC2245 – -d only applies to the first expansion of this glob. Use a loop to check any/all.
- SC2246 – This shebang specifies a directory. Ensure the interpreter is a file.
- SC2247 – Flip leading `$` and `"` if this should be a quoted substitution.
- SC2248 – Prefer double quoting even when variables don't contain special characters.
- SC2249 – Consider adding a default `*)` case, even if it just exits with error.
- SC2250 – Prefer putting braces around variable references even when not strictly required.
- SC2251 – This `!` is not on a condition and skips errexit. Add `|| exit 1` or make sure `$?` is checked.
- SC2252 – You probably wanted `&&` here, otherwise it's always true.
- SC2253 – Use `-R` to recurse, or explicitly `a-r` to remove read permissions.
- SC2254 – Quote expansions in case patterns to match literally rather than as a glob.
- SC2255 – `[ ]` does not apply arithmetic evaluation. Evaluate with `$((..))` for numbers, or use string comparator for strings.
- SC2256 – This translated string is the name of a variable. Flip leading `$` and `"` if this should be a quoted substitution.
- SC2257 – Arithmetic modifications in command redirections may be discarded. Do them separately.
- SC2258 – The trailing comma is part of the value, not a separator. Delete or quote it.
- SC2259 – This redirection overrides piped input. To use both, merge or pass filenames.
- SC2260 – This redirection overrides the output pipe. Use `tee` to output to both.
- SC2261 – Multiple redirections compete for stdout. Use `cat`, `tee`, or pass filenames instead.
- SC2262 – This alias can't be defined and used in the same parsing unit. Use a function instead.
- SC2263 – Since they're in the same parsing unit, this command will not refer to the previously mentioned alias.
- SC2264 – This function unconditionally re-invokes itself. Missing `command`?
- SC2265 – Use `&&` for logical AND. Single `&` will background and return true.
- SC2266 – Use `||` for logical OR. Single `|` will pipe.
- SC2267 – GNU `xargs -i` is deprecated in favor of `-I{}`
- SC2268 – Avoid x-prefix in comparisons as it no longer serves a purpose.
- SC2269 – This variable is assigned to itself, so the assignment does nothing.
- SC2270 – To assign positional parameters, use `set -- first second ..` (or use `[ ]` to compare).
- SC2271 – For indirection, use arrays, `declare "var$n=value"`, or (for sh) read/eval
- SC2272 – Command name contains `==`. For comparison, use `[ "$var" = value ]`.
- SC2273 – Sequence of `===`s found. Merge conflict or intended as a commented border?
- SC2274 – Command name starts with `===`. Intended as a commented border?
- SC2275 – Command name starts with `=`. Bad line break?
- SC2276 – This is interpreted as a command name containing `=`. Bad assignment or comparison?
- SC2277 – Use `BASH_ARGV0` to assign to `$0` in bash (or use `[ ]` to compare).
- SC2278 – `$0` can't be assigned in Ksh (but it does reflect the current function).
- SC2279 – `$0` can't be assigned in Dash. This becomes a command name.
- SC2280 – `$0` can't be assigned this way, and there is no portable alternative.
- SC2281 – Don't use `$`/`${}` on the left side of assignments.
- SC2282 – Variable names can't start with numbers, so this is interpreted as a command.
- SC2283 – Use `[ ]` to compare values, or remove spaces around `=` to assign (or quote `'='` if literal).
- SC2284 – Use `[ x = y ]` to compare values (or quote `'=='` if literal).
- SC2285 – Remove spaces around `+=` to assign (or quote `'+='` if literal).
- SC2286 – This empty string is interpreted as a command name. Double check syntax (or use 'true' as a no-op).
- SC2287 – This is interpreted as a command name ending with '/'. Double check syntax.
- SC2288 – This is interpreted as a command name ending with apostrophe. Double check syntax.
- SC2289 – This is interpreted as a command name containing a linefeed. Double check syntax.
- SC2290 – Remove spaces around = to assign.
- SC2291 – Quote repeated spaces to avoid them collapsing into one.
- SC2292 – Prefer `[[ ]]` over `[ ]` for tests in Bash/Ksh.
- SC2293 – When eval'ing @Q-quoted words, use * rather than @ as the index.
- SC2294 – eval negates the benefit of arrays. Drop eval to preserve whitespace/symbols (or eval as string).
- SC2295 – Expansions inside `${..}` need to be quoted separately, otherwise they will match as a pattern.
- SC2296 – Parameter expansions can't start with `{`. Double check syntax.
- SC2297 – Double quotes must be outside `${}`: `${"invalid"}` vs `"${valid}"`.
- SC2298 – `${$x}` is invalid. For expansion, use ${x}. For indirection, use arrays, ${!x} or (for sh) eval.
- SC2299 – Parameter expansions can't be nested. Use temporary variables.
- SC2300 – Parameter expansion can't be applied to command substitutions. Use temporary variables.
- SC2301 – Parameter expansion starts with unexpected quotes. Double check syntax.
- SC2302 – This loops over values. To loop over keys, use `"${!array[@]}"`.
- SC2303 – `i` is an array value, not a key. Use directly or loop over keys instead.
- SC2304 – `*` must be escaped to multiply: `\*`. Modern `$((x * y))` avoids this issue.
- SC2305 – Quote regex argument to expr to avoid it expanding as a glob.
- SC2306 – Escape glob characters in arguments to expr to avoid pathname expansion.
- SC2307 – 'expr' expects 3+ arguments but sees 1. Make sure each operator/operand is a separate argument, and escape <>&|.
- SC2308 – `expr length` has unspecified results. Prefer `${#var}`.
- SC2309 – -eq treats this as a variable. Use = to compare as string (or expand explicitly with $var)
- SC2310 – This function is invoked in an 'if' condition so set -e will be disabled. Invoke separately if failures should cause the script to exit.
- SC2311 – Bash implicitly disabled set -e for this function invocation because it's inside a command substitution. Add set -e; before it or enable inherit_errexit.
- SC2312 – Consider invoking this command separately to avoid masking its return value (or use '|| true' to ignore).
- SC2313 – Quote array indices to avoid them expanding as globs.
- SC2314 – In bats, `!` does not cause a test failure.
- SC2315 – In bats, `!` does not cause a test failure. Fold the `!` into the conditional!
- SC2316 – This applies local to the variable named readonly, which is probably not what you want. Use a separate command or the appropriate `declare` options instead.
- SC2317 – Command appears to be unreachable. Check usage (or ignore if invoked indirectly).
- SC2318 – This assignment is used again in this `declare`, but won't have taken effect. Use two `declare`s.
- SC2319 – This `$?` refers to a condition, not a command. Assign to a variable to avoid it being overwritten.
- SC2320 – This `$?` refers to echo/printf, not a previous command. Assign to variable to avoid it being overwritten.
- SC2321 – Array indices are already arithmetic contexts. Prefer removing the `$((` and `))`.
- SC2322 – In arithmetic contexts, `((x))` is the same as `(x)`. Prefer only one layer of parentheses.
- SC2323 – `a[(x)]` is the same as `a[x]`. Prefer not wrapping in additional parentheses.
- SC2324 – var+=1 will append, not increment. Use (( var += 1 )), declare -i var, or quote number to silence.
- SC2325 – Multiple ! in front of pipelines are a bash/ksh extension. Use only 0 or 1.
- SC2326 – ! is not allowed in the middle of pipelines. Use command group as in `cmd | { ! cmd; }` if necessary.
- SC2327 – This command substitution will be empty because the command's output gets redirected away.
- SC2328 – This redirection takes output away from the command substitution.
- SC2329 – This function is never invoked. Check usage (or ignored if invoked indirectly).
- SC2330 – BusyBox `[[ .. ]]` does not support glob matching. Use a case statement.
- SC2331 – For file existence, prefer standard -e over legacy -a.
- SC2332 – [ ! -o opt ] is always true because -o becomes logical OR. Use [[ ]] or ! [ -o opt ].
- SC2333 – You probably wanted || here, otherwise it's always false.
- SC2334 – You probably wanted || here, otherwise it's always false.
- SC2335 – Use `[ "$var" -ne 1 ]` instead of `[ ! "$var" -eq 1 ]`
- SC3001 – In POSIX sh, process substitution is undefined.
- SC3002 – In POSIX sh, extglob is undefined.
- SC3003 – In POSIX sh, `$'..'` is undefined.
- SC3004 – In POSIX sh, $".." is undefined
- SC3005 – In POSIX sh, arithmetic for loops are undefined.
- SC3006 – In POSIX sh, standalone `((..))` is undefined.
- SC3007 – In POSIX sh, `$[..]` in place of `$((..))` is undefined.
- SC3008 – In POSIX sh, select loops are undefined.
- SC3009 – In POSIX `sh`, brace expansion is undefined.
- SC3010 – In POSIX sh, `[[ ]]` is undefined.
- SC3011 – In POSIX sh, here-strings are undefined.
- SC3012 – In POSIX sh, lexicographical `\<` is undefined.
- SC3013 – In POSIX sh, `-nt` is undefined.
- SC3014 – In POSIX sh, `==` in place of `=` is undefined.
- SC3015 – In POSIX sh, `=~` regex matching is undefined.
- SC3016 – In POSIX sh, unary `-v` (in place of `[ -n "${var+x}" ]`) is undefined.
- SC3017 – In POSIX sh, unary `-a` in place of `-e` is undefined.
- SC3018 – In POSIX sh, `++` is undefined.
- SC3019 – In POSIX sh, exponentials are undefined.
- SC3020 – In POSIX sh, `&>` is undefined.
- SC3021 – In POSIX sh, `>& filename` (as opposed to `>& fd`) is undefined.
- SC3022 – In POSIX sh, named file descriptors is undefined.
- SC3023 – In POSIX sh, FDs outside 0-9 are undefined.
- SC3024 – In POSIX sh, `+=` is undefined.
- SC3025 – In POSIX sh, `/dev/{tcp,udp}` is undefined.
- SC3026 – In POSIX sh, `^` in place of `!` in glob bracket expressions is undefined.
- SC3028 – In POSIX sh, VARIABLE is undefined.
- SC3029 – In POSIX sh, `|&` in place of `2>&1 |` is undefined.
- SC3030 – In POSIX sh, arrays are undefined.
- SC3031 – In POSIX sh, redirecting from/to globs is undefined.
- SC3032 – In POSIX sh, coproc is undefined.
- SC3033 – In POSIX sh, naming functions outside [a-zA-Z_][a-zA-Z0-9_]* is undefined.
- SC3034 – In POSIX sh, `$(<file)` is undefined.
- SC3035 – In POSIX sh, `` `<file` `` is undefined.
- SC3036 – In Dash, echo flags besides -n are not supported.
- SC3037 – In POSIX sh, echo flags are undefined.
- SC3038 – In POSIX sh, exec flags are undefined.
- SC3039 – In POSIX sh, `let` is undefined.
- SC3040 – In POSIX sh, set option *[name]* is undefined.
- SC3041 – In POSIX sh, set flag `-E` is undefined
- SC3042 – In POSIX sh, set flag `--default` is undefined
- SC3043 – In POSIX sh, `local` is undefined.
- SC3044 – In POSIX sh, `declare` is undefined.
- SC3045 – In POSIX sh, some-command-with-flag is undefined.
- SC3046 – In POSIX sh, `source` in place of `.` is undefined.
- SC3047 – In POSIX sh, trapping ERR is undefined.
- SC3048 – In POSIX sh, prefixing signal names with 'SIG' is undefined.
- SC3049 – In POSIX sh, using lower/mixed case for signal names is undefined.
- SC3050 – In POSIX sh, `printf %q` is undefined.
- SC3051 – In POSIX sh, `source` in place of `.` is undefined
- SC3052 – In POSIX sh, arithmetic base conversion is undefined
- SC3053 – In POSIX sh, indirect expansion is undefined.
- SC3054 – In POSIX sh, array references are undefined.
- SC3055 – In POSIX sh, array key expansion is undefined.
- SC3056 – In POSIX sh, name matching prefixes are undefined.
- SC3057 – In POSIX sh, string indexing is undefined.
- SC3059 – In POSIX sh, case modification is undefined.
- SC3060 – In POSIX sh, string replacement is undefined.
- SC3061 – In POSIX sh, `read` without a variable is undefined.
- SC3062 – In POSIX sh, unary -o to check options is undefined.
- SC3063 – In POSIX sh, test -R and namerefs in general are undefined.
- SC3064 – In POSIX sh, test -N is undefined.
- SC3065 – In POSIX sh, test -k is undefined.
- SC3066 – In POSIX sh, test -G is undefined.
- SC3067 – In POSIX sh, test -O is undefined.

`$`
is not used specially and should therefore be escaped.*Note: Removed in v0.3.3
- 2014-05-29*

`echo "$"``echo "\$"``$` is special in double quotes, but there are some cases
where it's interpreted literally:

`echo "\$"``"foo$"`) or before some constructs
(`"$'foo'"`).To avoid relying on strange and shell-specific behavior, any
`$` intended to be literal should be escaped with a
backslash.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`\o`
will be a regular 'o' in this context.```
# Want literal backslash
echo Yay \o/
# Want other characters
bell=\a
```
```
echo 'Yay \o/'
bell="$(printf '\a')"
```
You have escaped something that has no special meaning when escaped. The backslash will be simply be ignored.

If the backslash was supposed to be literal, single quote or escape it.

If you wanted it to expand to something, rewrite the expression to
use `printf` (or in bash, `$'\t'`). If the
sequence in question is `\n`, `\t` or
`\r`, you instead get a SC1012 that
describes this.

None. ShellCheck (as of 2017-07-03,
commit `31bb02d6`)
will not warn when the first letter of a command is unnecessarily
escaped, as this is frequently used to suppress aliases
interactively.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo 'This is how it'\''s done'`.(Note: in v0.4.6, the error message was accidentally missing the backslash)

`echo 'This is not how it\'s done'.``echo 'This is how it'\''s done'.`In POSIX shell, the shell cares about nothing but another single-quote to terminate the quoted segment. Not even backslashes are interpreted.

POSIX.1 Shell Command Language § 2.2.2 Single Quotes:

Enclosing characters in single-quotes (

`''`) shall preserve the literal value of each character within the single-quotes. A single-quote cannot occur within single-quotes.

If you want your single-quoted string to end in a backslash, you can
rewrite as `'string'\\` or ignore this
warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

(This warning was retired after v0.7.2 due to low signal-to-noise ratio)

```
var='This is long \
piece of text'
```
```
var='This is a long '\
'piece of text'
```
You have a single-quoted string containing a backslash followed by a linefeed (newline). Unlike double-quotes or unquoted strings, this has no special meaning. The string will contain a literal backslash and a linefeed.

If you wanted to break the line but not add a linefeed to the string, stop the single quote, break the line, and reopen it. This is demonstrated in the correct code.

If you wanted to break the line and also include the linefeed as a literal, you don't need a backslash:

```
var='This is a multi-line string
with an embedded linefeed'
```
If you do want a string containing a literal backslash+linefeed
combo, such as with `sed`, you can ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`=` if trying to assign a value (or for empty
string, use `var=""` ... ).```
# I want programs to show text in Dutch!
LANGUAGE= nl
```
```
# I want to run the nl command with English error messages!
LANGUAGE= nl
```
```
# I want programs to show text in Dutch!
LANGUAGE=nl
```
```
# I want to run the nl command with English error messages!
LANGUAGE='' nl
```
It's easy to think that `LANGUAGE= nl` would assign
`"nl"` to the variable `LANGUAGE`. It doesn't.

Instead, it runs `nl` (the "number lines" command) and
sets `LANGUAGE` to an empty string in its environment.

Since trying to assign values this way is a common mistake,
ShellCheck warns about it and asks you to be explicit when assigning
empty strings (except for `IFS`, due to the common
`IFS= read ..` idiom).

If you're familiar with this behavior and feel that the explicit version is unnecessary, you can ignore it.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
#!/bin/mywrapper
echo "Hello World"
```
```
#!/bin/mywrapper
# shellcheck shell=bash
echo "Hello World"
```
You have specified a shebang that ShellCheck doesn't recognize. This
can be due to invoking the script via a wrapper, specifying a dummy like
`#!/bin/false` to prevent execution, or trying to check a
script for a non-Bourne shell or tool.

If this really is a sh/bash/dash/ksh script, please add a
`shell` directive after the shebang to tell ShellCheck how to
interpret the script, as in the example. You can also specify the shell
with the `-s` flag.

If this is a script in some other language, like
`#!/bin/sed` for a `sed` script, then sorry --
ShellCheck does not support `sed`, `awk`,
`expect` scripts. It only supports Bourne style shell
scripts.

None.

https://www.gnu.org/software/bash/manual/html_node/Shell-Scripts.html https://mywiki.wooledge.org/BashProgramming?highlight=%28shebang%29#Shebang https://www.gnu.org/software/gawk/manual/html_node/Executable-Scripts.html https://mywiki.wooledge.org/BashPitfalls#On_UTF-8_and_Byte-Order_Marks_.28BOM.29

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

This info warning points to the start of what ShellCheck was parsing when it failed. See [[Parser error]] for example and information.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`done` (or quote to make it
literal).(or `do` `then`, `fi`,
`esac`)

`for f in *; do echo "$f" done`or

`echo $f is done``for f in *; do echo "$f"; done`or

`echo "$f is done"`ShellCheck found a keyword like `done`, `then`,
`fi`, `esac`, etc used as the argument of a
command. This means that the shell will interpret it as a literal string
rather than a shell keyword. To be interpreted as a keyword, it must be
the first word in the line (i.e. after `;`,
`&` or a linefeed).

In the example, `echo "$f" done` is the same as
`echo "$f" "done"`, and the `done` does not
terminate the loop. This is fixed by terminating the `echo`
command with a `;` so that the `done` is the first
word in the next line.

If you're intentionally using `done` as a literal, you can
quote it to make this clear to ShellCheck (and also human readers), e.g.
instead of `echo Task is done`, use
`echo "Task is done"`. This makes no difference to the shell,
but it will silence this warning.

From POSIX-2018, section "C.2.10 Shell Grammar," regarding the
syntax, `if (false) then (echo x) else (echo y) fi`: https://pubs.opengroup.org/onlinepubs/9699919799/xrat/V4_xcu_chap02.html#tag_23_02_10

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo 'Nothing so needs reforming as other peoples' habits.'``echo 'Nothing so needs reforming as other peoples'\'' habits.'`or

`echo "Nothing so needs reforming as other peoples' habits."`When writing a string in single-quotes, you have to make sure that any apostrophes in the text don't accidentally terminate the single-quoted string prematurely.

Escape them properly (see the correct code) or switch quotes to avoid the problem.

```
echo '...peoples\ habits.'
...peoples\ habits.
```
```
$ echo $'...peoples\x27 habits.'
...peoples' habits.
```
None.

https://www.gnu.org/software/bash/manual/html_node/Quoting.html

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`\t`
is just literal `t` here. For tab, use
`"$(printf '\t')"` instead.```
# Want tab
$ var=foo\tbar
$ printf '<%s>\n' "$var"
<footbar>
$ var=foo\\tbar
$ printf '<%s>\n' "$var"
<foo\tbar>
```
or

```
# Want newline
$ var=foo\nbar
$ printf '<%s>\n' "$var"
<foonbar>
$ var=foo\\nbar
$ printf '<%s>\n' "$var"
<foo\nbar>
```
```
$ var="foo$(printf '\t')bar" # As suggested in warning
$ printf '<%s>\n' "$var"
<foo bar>
$ var="$(printf 'foo\tbar')" # Equivalent alternative
$ printf '<%s>\n' "$var"
<foo bar>
```
or

```
$ # Literal, quoted newline
$ line="foo
> bar"
$ printf '<%s>\n' "$line"
<foo
bar>
```
or

```
$ # Newline using ANSI-C quoting
$ line=$'foo\nbar'
$ printf '<%s>\n' "$line"
<foo
bar>
```
ShellCheck has found a `\t`, `\n` or
`\r` in a context where they just become regular letters
`t`, `n` or `r`. Most likely, it was
intended as a tab, newline or carriage return.

To generate such characters (plus other less common ones including
`\a`, `\f` and octal escapes) , use
`printf` as in the example. The exception is for newliness
that would be stripped by command substitution; in these cases, use a
literal quoted newline instead.

Other characters like `\z` generate a SC1001 info message, as the intent is less
certain.

None.

https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-printf https://pubs.opengroup.org/onlinepubs/9799919799/utilities/printf.html https://www.gnu.org/software/bash/manual/html_node/ANSI_002dC-Quoting.html

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`if cmd; then ..` to check exit code, or
`if [ "$(cmd)" = .. ]` to check output.```
# WRONG
if [ grep -q pattern file ]
then
 echo "Found a match"
fi
```
```
# WRONG
if ! [[ logname == $LOGNAME ]]
then
 echo "Possible `su` shell"
fi
```
```
if grep -q pattern file
then
 echo "Found a match"
fi
```
```
if ! [[ $(logname) == $LOGNAME ]]
then
 echo "Possible `su` shell"
fi
```
`[ ... ]` as shell syntax is a simple command that tests
for whether certain conditions are true or false, such as whether the
value assigned to a variable has a non-zero length
(`[ -n "${foo}" ]`) or whether a file system object is a
directory (`[ -d "${dir}" ]`).
`If-then-(elif-then)-else-fi` statements are logical
constructs which themselves contain lists of commands which can include
simple commands.

`[` is just regular command, like `whoami` or
`grep`, but with a funny name (see
`ls -l /bin/[`). It's a shorthand for `test`.
`[[` is similar to both `[` and `test`,
but `[[` offers some additional unary operators, such as '=~'
the regular expression comparison operator. It allows one to use
extglobs such as `@(foo|bar)` (a "bashism"), among some other
less commonly used features.

`[[`, `[` and `test` are often used
within `if...fi` constructs in the conditional commands
position: which is between the 'if' and the 'then.'

There are certain shell syntaxes which can be wrapped
**directly** around simple commands, in particular:

`{ ...;}`, group commands,`$( ... )`, command substitutions,`<( ... )` and `>( ... )`, process
substitutions,`( ... )`, subshells, and`$(( ... ))` and `(( ... ))`, arithmetic
evaluations.Some examples include:

`{ echo {a..z}; echo {0..9};} > ~/f`,`[[ $(logname) == $LOGNAME ]]`,`readarray -t files < <( find ...)`,`(cd /foo || exit 1; tar ...)`, and`dd bs=$((2**12)) count=1 if=/dev/zero of=/tmp/zeroed-block`,
respectively.Note how in example (2) `logname` is enclosed
**directly** within a command substitution, which is itself
enclosed within a `[[` reserved word / conditional expression
/ compound command.

If you want to check the exit status of a certain command, use that command directly as demonstrated in the correct code, above.

If you want to check the output of a command, use
`"$(..)"` to get its output, and then use
`test`/`[` or `[[` to do a string
comparison:

```
# Check output of `whoami` against the string `root`
if [ "$(whoami)" = "root" ]
then
 echo "Running as root"
fi
```
None.

For more information, see this problem in the Bash Pitfall list, or generally Tests and Conditionals in the wooledge.org BashGuide

`if [grep foo myfile]`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo “hello world”``echo "hello world"`Blog software and word processors frequently replaces ASCII quotes
`""` with fancy Unicode quotes, `“”`. To Bash,
Unicode quotes are considered regular literals and not quotes at
all.

Simply delete them and retype them in your editor.

This error was retired after 0.4.5. In this version and earlier, ShellCheck parsed slanted quotes as a valid double quote. This meant that the warning could not simply be ignored. It has since been replaced by SC1110 (outside quotes) and SC1111 (inside double-quotes).

If you really want literal Unicode double quotes, you can put them in single-quotes (or Unicode single-quotes in double-quotes) to make ShellCheck ignore them, e.g.,

`printf 'Warning: “wakeonlan” is not installed.\n'`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo ‘hello world’``echo 'hello world'`Some software, like macOS, Microsoft Word, and WordPress, may automatically replace your regular quotes with slanted Unicode quotes. Try deleting and retyping them, and/or disable “smart quotes” in your editor or OS.

This error was retired after 0.4.5. In this version and earlier, ShellCheck parsed slanted quotes as a valid double-quote. This meant that the warning could not simply be ignored. It has since been replaced by SC1110 (outside quotes) and SC1112 (inside single-quotes).

If you want to use typographic single-quotes, you can put them in double-quotes (or typographic double-quotes in single-quotes) to make ShellCheck ignore them, e.g.,

`printf "Warning: ‘wakeonlan’ is not installed.\n"`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`tr -d '\r'` .```
$ cat -v myscript
#!/bin/sh^M
echo "Hello World"^M
```
```
$ cat -v myscript
#!/bin/sh
echo "Hello World"
```
The script uses Windows/MS-DOS style `\r\n` line
terminators instead of Unix-style `\n` terminators. The
additional `\r` aka `^M` aka carriage return
characters will be treated literally, and results in all sorts strange
bugs and messages.

You can verify this with `cat -v yourfile` and see whether
or not each line ends with a `^M`. To delete them, open the
file in your editor and save the file as "Unix", "Unix/macOS Format",
`:set ff=unix` or similar if it supports it.

If you don't know how to get your editor to save a file with Unix
line terminators, you can use `tr`:

```
tr -d '\r' < badscript > goodscript
or
dos2unix badscript
```
This will read a script `badscript` with possible carriage
returns, and write `goodscript` without them.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

You copy-pasted some code, probably from a blog or web site, which for formatting reasons contained Unicode no-break spaces or Unicode zero-width spaces instead of regular spaces or in words.

To humans, a zero-width space is invisible and a non-breaking space is indistinguishable from a regular space, but the shell does not agree.

If you have just a few, delete the indicated space/word and retype it. If you have tons, do a search-and-replace in your editor (copy-paste an offending space into the search field, and type a regular space into the replace field), or use the following command to remove them:

`sed -e $'s/\xC2\xA0/ /g' -e $'s/\xE2\x80\x8b//g' -i yourfile`On macOS, a non-breaking space can be inserted into most programs by
holding `⌥ Option`+`Space`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ -x ]``[ -x "myfile" ]`ShellCheck has found a unary test operator that does not appear to be followed by a valid shell word.

This could be because of a misplaced `]`, `)`,
or a missing space before the `]`.

Check the syntax, make sure the test operator has an operand, and try again.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`]` or `]]``if [ "$STUFF" = ""]; then``if [ "$STUFF" = "" ]; then`Bourne shells are very whitespace sensitive. Adding or removing spaces can drastically alter the meaning of a script. In these cases, ShellCheck has noticed that you're missing a space at the position indicated.

None.

```
# shellcheck disable=SC1020
if [ "$STUFF" = ""]; then
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[..](..)`, use `( .. )`.```
[[ [ a || b ] && c ]]
[ [ a -o b ] -a c ]]
```
```
[[ ( a || b ) && c ]]
[ \( a -o b \) -a c ]] # or { [ a ] || [ b ]; } && [ c ]
```
`[ .. ]` should not be used to group subexpressions inside
`[[ .. ]]` or `[ .. ]` statements.

For `[[ .. ]]`, use regular parentheses.

For `[ .. ]`, either use escaped parentheses, or
preferably rewrite the expression into multiple `[ .. ]`
joined with `&&`, `||` and
`{ ..; }` groups. The latter is preferred because
`[ .. ]` is undefined for more than 4 arguments in POSIX.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ $a -ne ]``[ $a -ne $b ]`ShellCheck found a `test` operator without an operand.
This could be a copy-paste fail, bad linebreak, or trying to use
`<>` instead of `!=` or
`-ne`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[..]` you have to escape `\( \)` or preferably
combine `[..]` expressions.`[ -e ~/.bashrc -a ( -x /bin/dash -o -x /bin/ash ) ]`In POSIX:

`[ -e ~/.bashrc ] && { [ -x /bin/dash ] || [ -x /bin/ash ]; }`Obsolete XSI syntax:

`[ -e ~/.bashrc -a \( -x /bin/dash -o -x /bin/ash \) ]``[` is implemented as a regular command, so `(`
is not special.

The preferred way is not to group inside `[ .. ]` and
instead compose multiple `[ .. ]` statements using the
shell's `&&`, `||` and
`{ ..; }` syntax, since this is well defined by POSIX.

Some shells, such as Bash, support grouping with
`\( .. \)`, but this is an obsolete XSI-only extension.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[..](..)` you
shouldn't escape `(` or `)`.`[[ -e ~/.bashrc && \( -x /bin/dash || -x /bin/ash \) ]]``[[ -e ~/.bashrc && ( -x /bin/dash || -x /bin/ash ) ]]`You don't have to -- and can't -- escape `(` or
`)` inside a `[[ .. ]]` expression like you do in
`[ .. ]`. Just remove the escaping.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[[` but closed with single
`]`. Make sure they match.(or SC1034 for vice versa)

`[[ -z "$var" ]``[[ -z "$var" ]]`ShellCheck found a test expression `[ ... ]` (POSIX) or
`[[ ... ]]` (ksh/bash), but where the opening and closing
brackets did not match (i.e. `[[ .. ]` or
`[ .. ]]`). The brackets need to match up to work.

Note in particular that `[..]` do *not* work like
parentheses in other languages. You can not do:

```
# Invalid
[[ x ] || [ y ]]
```
You would instead use two separate test expressions joined by
`||`:

```
# Valid basic test expressions (sh/bash/ksh)
[ x ] || [ y ]
# Valid extended test expressions (bash/ksh)
[[ x ]] || [[ y ]]
```
None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

[

]]

See similar error SC1033

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`if ![-z foo ]; then true; fi # if command `[-z' w/ args `foo', `]' fails..``if ! [ -z foo ]; then true; fi # if command `[' w/ args `-z', `foo', `]' fails..`Bourne shells are very whitespace sensitive. Adding or removing spaces can drastically alter the meaning of a script. In these cases, ShellCheck has noticed that you're missing a space at the position indicated.

ShellCheck does not understand Bash
History Expansion, an interactive shell feature also using
`!` (such as `!!` to expand to the previous
command).

These features are disabled by default in shells and very rarely used
in scripts, but may occasionally be found in interactively sourced files
like `.bashrc`. Please ignore the error in these cases.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`(` is
invalid here. Did you forget to escape it?`echo (foo) bar`Depends on your intention:

```
echo "(foo) bar" # Literal parentheses
echo "$(foo) bar" # Command expansion
echo "foo bar" # Tried to use parentheses for grouping or function invocation
```
ShellCheck expected an ordinary shell word but found an opening parenthesis instead.

Determine what you intended the parenthesis to do and rewrite accordingly. Common issues include:

`echo (FAIL) Some tests failed`. In this case, it requires
quoting.`echo Today is (date)`.
Add the missing `$`:
`echo "Today is $(date)"``foo (bar, 42)` to call a function. This
should be `foo bar 42`. Also, shells do not support tuples or
passing arrays as single parameters.Bash allows some parentheses as part of assignment-like tokens to
certain commands, including `export` and `eval`.
This is a workaround in Bash to allow commands that normally would not
be valid:

```
eval foo=(bar) # Valid command
echo foo=(bar) # Invalid syntax
f=foo; eval $f=(bar) # Also invalid
```
In these cases, please quote the command, such as
`eval "foo=(bar)"`. This does not change the behavior, but
stops relying on Bash-specific parsing quirks.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`${10}`.```
echo "Ninth parameter: $9"
echo "Tenth parameter: $10"
```
```
echo "Ninth parameter: $9"
echo "Tenth parameter: ${10}"
```
For legacy reasons, `$10` is interpreted as the variable
`$1` followed by the literal string `0`.

Curly braces are needed to tell the shell that both digits are part of the parameter expansion.

If you wanted the trailing digits to be literal, `${1}0`
will make this clear to both humans and ShellCheck.

In `dash`, `$10` is (wrongly)
interpreted as `${10}`, so some 'reversed' care should also
be taken:

```
bash -c 'set a b c d e f g h i j; echo $10 ${1}0' # POSIX: a0 a0
dash -c 'set a b c d e f g h i j; echo $10 ${1}0' # WRONG: j a0
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`< <(cmd)`, not
`<<(cmd)`.```
while IFS= read -r line
do
 printf "%q\n" "$line"
done <<(curl -s http://example.com)
```
```
while IFS= read -r line
do
 printf "%q\n" "$line"
done < <(curl -s http://example.com)
```
You are using `<<(` which is an invalid
construct.

You probably meant to redirect `<` from process
substitution `<(..)` instead. To do this, a space is
needed between the `<` and `<(..)`, i.e.
`< <(cmd)`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`<<-` and indent
with tabs).```
for f in *.png
do
 cat << HTML
 <img src="$f" /><br/>
 HTML
done > index.html
```
```
for f in *.png
do
 cat << HTML
 <img src="$f" /><br/>
HTML
done > index.html
```
The here document delimiter will not be recognized if it is indented.

You can fix it in one of two ways:

`<<-` instead of `<<`, and
indent the script with tabs only (spaces will not be recognized).Removing the indentation is preferred, since the script won't suddenly break if it's reformatted, copy-pasted, or saved with a different editor.

If the line was supposed to be a literal part of the here document, consider choosing a less ambiguous token.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`<<-`, you can only indent with tabs.Any code using `<<-` that is indented with spaces.
`cat -T script` shows

```
 cat <<- foo
 Hello world
 foo
```
Code using `<<-` must be indented with tabs.
`cat -T script` shows

```
^Icat <<- foo
^I^IHello world
^Ifoo
```
Or simply don't indent the end token:

```
 cat <<- foo
 Hello World
foo
```
`<<-`, by design, only strips tabs. Not spaces.

Your editor may be automatically replacing tabs with spaces, either when you type them or when you save the file or both. If you're unable to make it stop, just don't indent the end token.

None. But note that copy-pasting code to shellcheck.net may also turn correct tabs into spaces on some OS.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`eof` further down, but not on a separate line.`Close matches include '-eof' (!= 'eof').````
cat <<-eof
Hello World
-eof
```
```
cat <<- eof
Hello World
eof
```
Your here document isn't properly terminated.

There is a line containing the terminator you've chosen, but it's not by itself on a separate line.

In the example code, the script uses `<<-eof`, which
is the operator `<<-` followed by `eof`. The
script therefore looks for `eof` and skips right past the
intended terminator because it starts with a dash.

You will get some companion SC1042 errors mentioning lines that contain the string as a substring, though they all point to the start of the here document and not the relevant line:

```
In foo line 4:
Hello
^-- SC1041: Found 'eof' further down, but not on a separate line.
^-- SC1042: Close matches include '-eof' (!= 'eof').
```
Look at your here document and see which line was supposed to terminate it. Then ensure it matches the token exactly, and that it's on its own line with no text before or after.

Under Windows the error might occur due to the standard CRLF line-ending, which is Windows-specific. Try to change the line ending into LF.

None.

Note that SC1041 and SC1042 swapped numbers after v0.4.6 to improve the display order. This rare instance of number reuse was justified by them always occurring together on the same line.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
cat << EOF
Hello World
Eof
```
```
cat << EOF
Hello World
EOF
```
ShellCheck found a here document (`<<`) where the
end token is missing. However, the end token appears with different case
further down. If this was meant to be the end of the here document, make
sure the case matches.

None. This error is only emitted when the here document is incomplete.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`EOF` in the here document.```
cat << EOF
 Hello World
```
```
cat << EOF
 Hello World
EOF
```
The `<<` here document (aka heredoc) was not
properly terminated. The terminating token needs to be on a separate
line without indenting (or indented with tabs only when using
`<<-`).

Note that you can not put here documents in one liners. For such use
cases, use a `<<<` here string:

```
cat << EOF hello world EOF # Wrong: data and terminator can not be on the same line
cat <<< "hello world" # Correct
```
None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

foo &; bar

foo & bar

Both & and ; terminate the command. You should only use one of them.

&

;

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`fi` for
this `if`.```
if true
then
 echo "True"
done
```
```
if true
then
 echo "True"
fi
```
ShellCheck has found an `if` statement that does not
appear to have a matching terminating `fi`.

This could be because it's missing entirely, or because the
`if` statement was incorrectly terminated by a mismatched
`done`, `esac`, `)` or similar. A
companion warning SC1047 is emitted at the point
where ShellCheck expected the `fi`.

Check that the `if` statement is completely and correctly
terminated.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

fi

if

See companion warning SC1046.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`true` as a no-op).```
if [ -e foo ]
then
 # TODO: handle this
fi
```
```
if [ -e foo ]
then
 # TODO: handle this
 true
fi
# Or use the no-op colon operator ":"
if [ -e foo ]
then
 # TODO: handle this
 :
fi
```
Shells do not allow empty `then` clauses. They need at
least one command (and comments are not commands).

If you want a `then` clause that does nothing, use a dummy
command like `true`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`then` for this `if`?```
if true
 echo "foo"
elif true
 echo "bar"
fi
```
```
if true
then
 echo "foo"
elif true
then
 echo "bar"
fi
```
ShellCheck found a parsing error in the script, and determined that
it's most likely due to a missing `then` keyword for the
`if` or `elif` indicated.

Make sure the `then` is there.

Note that the `then` needs a `;` or linefeed
before it. `if true then` is invalid, while
`if true; then` is correct.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`then`.```
if true
 echo "True"
fi
```
```
if true
then
 echo "True"
fi
```
ShellCheck has found an `if` statement that appears to be
missing a `then`.

Make sure the `then` exists, and that it is the first word
of the line (or immediately preceded by a semicolon).

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`then` are not allowed. Just remove it.`if true; then; echo "Hi"; fi``if true; then echo "Hi"; fi``then` keywords should not be followed by semicolons. It's
not valid shell syntax.

You can follow them directly with a line break or another command.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`then` are not allowed. Just remove it.`if true; then; echo "Hi"; fi``if true; then echo "Hi"; fi``then` keywords should not be followed by semicolons. It's
not valid shell syntax.

You can follow them directly with a line break or another command.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`else` are not allowed. Just remove it.`if mycommand; then echo "True"; else; echo "False"; fi``if mycommand; then echo "True"; else echo "False"; fi``else` keywords should not be followed by semicolons. It's
not valid shell syntax.

You can follow them directly with a line break or another command.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`{`.`foo() {echo "hello world"; }``foo() { echo "hello world"; }``{` is only recognized as the start of a command group
when it's a separate token.

If it's not a separate token, like in the problematic example, it
will be considered a literal character, as if writing
`"{echo"` with quotes, and therefore usually cause a syntax
error.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`true;` as a no-op.```
submitbug() {
 # TODO: Implement me
}
```
```
submitbug() {
 # TODO: Implement me
 true
}
```
ShellCheck found an empty code block. This could be an empty function as shown, a loop with an empty body, or similar.

Sh/bash does not allow empty code blocks. Insert at least one
command. If you don't want the block to do anything, `true`
(aka `:`) is a good no-op.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`}`. If you have one, try a `;` or
`\n` in front of it.```
#!/bin/bash
bar() { echo "hello world" }
```
```
#!/bin/bash
bar() { echo "hello world";}
```
`}` is only recognized as the end of a command group when
it's a separate token.

If it's not a separate token, like in the problematic example, it
will be considered a literal character, as if writing
`echo "foo}"` with quotes, and therefore usually cause a
syntax error.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`do` for this loop?```
while read -r line
 echo $line
done
```
```
while read -r line
do
 echo $line
done
```
ShellCheck found a loop that appears to be missing its
`do` statement. Make sure the loop syntax is correct.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`do`.```
for file in *
 echo "$file"
done
```
```
for file in *
do
 echo "$file"
done
```
ShellCheck has found a loop that appears to be missing a
`do` statement. In the problematic code, it was simply
forgotten.

Verify that the `do` exists, and that it's in the correct
place.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`do`. You can just delete
it.```
while true; do; true; done
while true;
do;
 true;
done;
```
```
while true; do true; done
while true;
do
 true;
done;
```
Semicolon `;` is not allowed directly after a
`do` keyword. Follow it directly with either a command or a
linefeed as shown in the example.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`true` as a no-op)```
for i in 1 2 3; do
done
```
```
for i in 1 2 3; do
 true
done
```
An empty `do ... done` block is not valid. Use
`true` or `:` if you need no command at all.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`done`
for this `do`.```
yes() {
 while echo "y"
 do
 true
}
```
```
yes() {
 while echo "y"
 do
 true
 done
}
```
ShellCheck found a `do` without a corresponding
`done`.

Double check that the `done` exists, and that it correctly
matches the indicated `do`. A companion warning SC1062 is emitted where ShellCheck first noticed it
was missing.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

done

do

See companion warning SC1061

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`do`.```
for file in * do
 echo "$file"
done
```
```
for file in *; do
 echo "$file"
done
# or
for file in *
do
 echo "$file"
done
```
ShellCheck found a `do` on the same line as a loop, but
`do` only starts a loop block at the start of a
line/statement. Make the `do` the start of a new
line/statement by inserting a linefeed or semicolon in front of it.

If you wanted to treat `do` as a literal string, you can
quote it to make this clear to ShellCheck and humans:

```
for f in "for" "do" "done"
do
 echo "Shell keywords include: $f"
done
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`{` to open the function definition.```
foo() {
 echo "hello world"
}
foo()
```
```
foo() {
 echo "hello world"
}
foo
```
ShellCheck found what appears to be the start of a function definition, but without a function body.

One common cause is that you are trying to call a function by
appending parentheses, e.g. `foo()` like in C. Bash does not
use or allow parentheses after a function name to call it. The function
`foo` should be called using just `foo` like in
the example.

If you are declaring a function, make sure it looks like the correct
code above, and that it does not try to declare any parameters
(parameters are instead accessed with `$1` and up).

If you are trying to do something else, look up the syntax for what you are trying to do.

POSIX allows the body of a function to be any compound command, e.g.
`foo() for i; do :; done`. Since this usage is rare,
ShellCheck intentionally requires the body to be `{ ..; }`
(or `( ..; )`):

```
foo() {
 for i; do :; done
}
```
This additional structure requirement helps improve error messages and suggestions by not parsing down a path that less advanced users wouldn't expect.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`()` and refer to params as
`$1`, `$2`, …```
foo(input) {
 echo "$input"
}
foo("hello world");
```
```
foo() {
 echo "$1"
}
foo "hello world"
```
Shell script functions behave just like scripts and other commands:

`$1`, `$2` etc. They can not declare parameters by
name.`name arg1 arg2`, and not with
parentheses as C-like languages.None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$` on the left side of assignments.`$greeting="Hello World"``greeting="Hello World"`Alternatively, if the goal was to assign to a variable whose name is
in another variable (indirection), use `declare`:

```
name=foo
declare "$name=hello world"
echo "$foo"
```
Or if you actually wanted to compare the value, use a test expression:

```
if [ "$greeting" = "hello world" ]
then
 echo "Programmer, I presume?"
fi
```
Unlike Perl or PHP, `$` is not used when assigning to a
variable.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`declare "var$n=value"`, or (for sh)
`read`/`eval`Note: Removed in v0.7.2 - 2021-04-20

```
n=1
var$n="hello"
```
For integer indexing in ksh/bash, consider using an indexed array:

```
n=1
var[n]="hello"
echo "${var[n]}"
```
For string indexing in ksh/bash, use an associative array:

```
typeset -A var
n="greeting"
var[$n]="hello"
echo "${var[$n]}"
```
If you actually need a variable with the constructed name in bash,
use `declare`:

```
n="Foo"
declare "var$n=42"
echo "$varFoo"
```
For `sh`, with single line contents, consider
`read`:

```
n="Foo"
read -r "var$n" << EOF
hello
EOF
echo "$varFoo"
```
or with careful escaping, `eval`:

```
n=Foo
eval "var$n='hello'"
echo "$varFoo"
```
`var$n=value` is not a valid way of assigning to a
dynamically created variable name in any shell. Please use one of the
other methods to assign to names via expanded strings. Wooledge BashFaq #6
has significantly more information on the subject.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`=` in assignments.Note: Removed in v0.7.2 - 2021-04-20

`foo = 42``foo=42`Shells are space sensitive. `foo=42` means to assign
`42` to the variable `foo`. `foo = 42`
means to run a command named `foo`, and pass `=`
as `$1` and `42` as `$2`.

If you actually wanted to run a command named foo and provide
`=` as the first argument, simply quote it to make ShellCheck
be quiet: `foo "=" 42`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[`.```
if[ -e file ]
then echo "exists"; fi
```
```
if [ -e file ]
then echo "exists"; fi
```
ShellCheck found a keyword immediately followed by a `[`.
There needs to be a space between them.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

ShellCheck found a syntax error at the indicated location. Barring a bug in ShellCheck itself, your shell will also crash with a syntax error at the same location, so you cannot ignore this check.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
#!/usr/bin/python
print "Hello"
```
You have specified the shebang of an unsupported language or shell dialect.

ShellCheck only supports a limited number of Bourne-based Unix shells: bash, ksh, dash and POSIX sh.

It does not support scripts written for other shells like Zsh, Csh, Tcsh or PowerShell, and it does not support other scripting languages like PHP, Python, JavaScript or SQL.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

*Note: There is a known bug in the
current version when directives appear
within then clauses of if blocks that causes
Shellcheck to report SC1072 on otherwise valid code. Avoid using
directives within then clauses - instead place them at the
top of the if block or another enclosing block. This is
fixed on the online version
and the next release.*

See Parser Error.

This error can also occur with an incomplete
shellcheck directive like `# shellcheck disable` instead
of `# shellcheck disable=all`, or if the corresponding
command has been commented out.

```
# shellcheck disable
echo stuff that shellcheck up to at least v0.10.0 will not even see because of the incorrect directive above
```
```
#shellcheck disable=SC2162
#read VAR
```
```
# shellcheck disable=all
echo stuff that shellcheck will correctly ignore entirely
```
```
##shellcheck disable=SC2162
#read VAR
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

This parsing error points to the structure ShellCheck was trying to parse when a parser error occurred.

Make any necessary fixes and check the script again, since most of ShellCheck's functionality can only be applied to scripts that parse successfully.

See [[Parser error]] for more information.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`;;` after the previous case item?```
while getoptions f option
do
 case "${options}"
 in
 f) FTR="${ARG}"
 \?) exit
 esac
done
```
```
while getoptions f option
do
 case "${options}"
 in
 f) FTR="${ARG}";;
 \?) exit;;
 esac
done
```
Syntax `case` needs `;;` after the previous
case item. If not, syntax error will cause.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`elif` instead of
`else if`.```
if [ "$#" -eq 0 ]
then
 echo "Usage: ..."
else if [ "$#" -lt 2 ]
then
 echo "Missing operand"
fi

```
```
if [ "$#" -eq 0 ]
then
 echo "Usage: ..."
elif [ "$#" -lt 2 ]
then
 echo "Missing operand"
fi
```
Many languages allow alternate branches with `else if`,
but `sh` is not one of them. Use `elif`
instead.

`else if` is a valid (though confusing) way of nesting an
`if` statement in a parent's `else`. If this is
your intention, consider using canonical formatting by putting a
linefeed between `else` and `if`.

This does not change the behavior of the script, but merely makes it
more obvious to ShellCheck (and other humans) that you didn't expect the
`else if` to behave the way it does in C. Alternatively, you
can ignore it with no ill effects.

```
if x
then
 echo "x"
else # line break here
 if y
 then
 echo "y"
 fi
fi
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ $((i/2+7)) -ge 18 ]`.`[ i / 2 + 7 -ge 18 ]``[ $((i / 2 + 7)) -ge 18 ]`ShellCheck found a loose `+*/%` in a test statement. This
usually happens when trying to do arithmetic in a condition, but without
using the arithmetic expansion construct
`$((expression))`.

In C, `if (a+b == c)` is perfectly fine, but in sh this
must be written to first expand the arithmetic operation like
`if [ $((a+b)) = c ]`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

``` vs
`´`).`echo "Your username is ´whoami´"````
echo "Your username is $(whoami)" # Preferred
echo "Your username is `whoami`" # Deprecated, will give [SC2006]
```
In some fonts it's hard to tell ticks apart, but Bash strongly
distinguishes between backticks (grave accent ```), forward
ticks (acute accent `´`) and regular ticks (apostrophe
`'`).

Backticks start command expansions, while forward ticks are literal. To help spot bugs, ShellCheck parses backticks and forward ticks interchangeably.

If you want to write out literal forward ticks, such as fancyful ascii quotation marks:

`echo "``Proprietary software is an injustice.´´ - Richard Stallman"`use single quotes instead:

`echo '``Proprietary software is an injustice.´´ - Richard Stallman'`To nest forward ticks in command expansion, use `$(..)`
instead of ``..``.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
greeting="hello
target="world"
```
```
greeting="hello"
target="world"
```
The first line is missing a quote.

ShellCheck warns when it detects multi-line double quoted, single quoted or backticked strings when the character that follows it looks out of place (and gives a companion warning SC1079 at that spot).

If you do want a multiline variable, just make sure the character after it is a quote, space or line feed.

```
var='multiline
'value
```
can be rewritten for readability and to remove the warning:

```
var='multiline
value'
```
As always ``..`` should be rewritten to
`$(..)`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

# SC1079 – ShellCheck Wiki

 See this page on GitHub
 Sitemap

 ## This
is actually an end quote, but due to next char it looks suspect.

See companion warning SC1078.

 ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`\` before line feeds to break lines in
`[ ]`.```
[ "$filename" =
 "$otherfile" ]
```
```
[ "$filename" = \
 "$otherfile" ]
```
Bash/ksh `[[ ]]]` can include line breaks anywhere, but
`[ ]` requires that you escape them. If you are writing a
multi-line `[ .. ]` statement, make sure to include these
escapes. If the `[ ]` is supposed to be on a single line,
make sure the `]` is there.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`if`, not `If`.```
If true
Then
 echo "hello"
Fi
```
```
if true
then
 echo "hello"
fi
```
Shells are case sensitive and do not accept `If` or
`IF` in place of lowercase `if`.

If you're aware of this and insist on naming a function
`WHILE`, you can quote the name to prevent shellcheck from
thinking you meant `while`. Or if you really want the names,
add things like `alias If=if IF=if` to replace those keywords
and ask shellcheck to ignore them.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`LC_CTYPE=C sed '1s/^...//' < yourscript`.This is an encoding error that can't be seen in the script itself,
but `cat -v` will show three bytes of garbage at the start of
the file:

```
$ cat -v file
M-oM-;M-?#!/bin/bash
echo "hello world"
```
The code is correct when this garbage does not appear.

Some editors may save a file with a Byte Order Mark to mark the file as UTF-8. Shells do not understand this and will give errors on the first line:

```
$ bash myscript
myscript: line 1: #!/bin/sh: No such file or directory
$ dash myscript
myscript: 1: myscript: #!/bin/sh: not found
```
To fix it, remove the byte order mark. One way of doing this is
`LC_CTYPE=C sed '1s/^...//' < yourscript`. Verify that
it's not there with `cat -v`.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`{`/`}` is literal. Check if `;` is
missing or quote the expression.`rmf() { rm -f "$@" }`or

`eval echo \${foo}``rmf() { rm -f "$@"; }`and

`eval "echo \${foo}"`Curly brackets are normally used as syntax in parameter expansion, command grouping and brace expansion.

However, if they don't appear alone at the start of an expression or as part of a parameter or brace expansion, the shell silently treats them as literals. This frequently indicates a bug, so ShellCheck warns about it.

In the example function, the `}` is literal because it's
not at the start of an expression. We fix it by adding a `;`
before it.

In the example eval, the code works fine. However, we can quiet the warning and follow good practice by adding quotes around the literal data.

ShellCheck does not warn about `{}`, since this is
frequently used with `find` and rarely indicates a bug.

This error is harmless when the curly brackets are supposed to be
literal, in e.g. `awk {'print $1'}`. However, it's cleaner
and less error prone to simply include them inside the quotes:
`awk '{print $1}'`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`#!`, not
`!#`, for the shebang.```
!#/bin/sh
echo "Hello World"
```
```
#!/bin/sh
echo "Hello World"
```
The shebang has been accidentally swapped. The `#` should
come first: `#!`, not `!#`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
case $options in
 foobar) echo foo ;; echo bar;;
 *) echo unknown option ;;
esac
```
```
case $options in
 foobar) echo foo ; echo bar;;
 *) echo unknown option ;;
esac
```
There should be no statements between `;;` and the next
case item.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$` on the iterator name in for loops.```
for $var in *
do
 echo "$var"
done
```
```
for var in *
do
 echo "$var"
done
```
The variable is named `var`, and can be expanded to its
value with `$var`.

The `for` loop expects the variable's name, not its value
(and the name can not be specified indirectly).

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`${array[idx]}` (or
`${var}[..` to quiet).`echo "$array[@]"``echo "${array[@]}"`Some languages use the syntax `$array[index]` to access an
index of an arrays, but a shell will interpret this as
`$array` followed by the unrelated literal string (or glob)
`[index]`.

Curly braces are needed to tell the shell that the square brackets are part of the expansion.

If you want the square brackets to be treated literally or as a glob,
use `${var}[idx]` to prevent this warning.

This does not change how the script works, but clarifies your intent to ShellCheck as well as other programmers.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`grep ^(.*)\1$ file`or

`var=myfunction(value)``grep '^(.*)\1$' file`or

`var=$(myfunction value)`Parentheses are shell syntax and must be used correctly.

For commands that expect literal parentheses, such as
`grep` or `find`, the parentheses need to be
quoted or escaped so the shell does not interpret them, but instead
passes them to the command.

For shell syntax, the shell does not use them the way most other languages do, so avoid guessing at syntax based on previous experience. In particular:

Parentheses are NOT used to call functions.

Parentheses are NOT used to group expressions, except in arithmetic contexts.

Parentheses are NOT used in conditional statements or loops.

Parentheses are used differently in different contexts.
`( .. )`, `$( .. )`, `$(( .. ))` and
`var=(..)` are completely separate and independent structures
with different meanings, and can not be broken down into operations on
expressions in parentheses.

In C-like languages, `++` can't be broken down into two
`+` operations, so you can't e.g. use `+ +` or
`+(+)`. In the same way, all of the above are completely
unrelated so that you can't do `$(1+1)` or
`$( (1+1) )` in place of `$(( 1+1 ))`.

If you are trying to use parentheses for shell syntax, look up the actual syntax of the statement you are trying to use.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
if true
then
 echo hello
fi
fi
```
```
if true
then
 echo hello
fi
```
This error is typically seen when there are too many `fi`,
`done` or `esac`s, or when there's a
`do` or `then` without a corresponding
`while`, `for` or `if`. This is often
due to deleting a loop or conditional statement but not its
terminator.

In some cases, it can even be caused by bad quoting:

```
var="foo
if [[ $var = "bar ]
then
 echo true
fi
```
In this case, the `if` ends up inside the double quotes,
leaving the `then` dangling.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`. "${util_path}"````
# shellcheck source=src/util.sh
. "${util_path}"
```
ShellCheck is not able to include sourced files from paths that are determined at runtime. The file will not be read, potentially resulting in warnings about unassigned variables and similar.

Use a Directive to point shellcheck to a fixed location it can read instead.

ShellCheck v0.7.2+ will strip a single expansion followed by a slash,
e.g. `${var}/util.sh` or
`$(dirname "${BASH_SOURCE[0]}")/util.sh`, and treat them as
`./util.sh`. This allows the use of `source-path`
directives or `-P` flags to specify the location.

If you don't care that ShellCheck is unable to account for the file,
specify `# shellcheck source=/dev/null`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Reasons include: file not found, no permissions, not included on the
command line, not allowing `shellcheck` to follow files with
`-x`, etc.

`source somefile`In case you have access to the file, e.g. if it is included in your source code repository:

```
# shellcheck source=somefile
source somefile
```
In case you do not have access to the file:

```
# shellcheck source=/dev/null
source somefile
```
ShellCheck, for whichever reason, is not able to access the source file.

This could be because:

`shellcheck -x` (or specified
`external-sources=true` in the .shellcheckrc) to allow following
other filesFeel free to ignore the error with a directive.

ShellCheck is unable to follow dynamic paths, such as
`source "$somedir/file"`. For these cases, see SC1090: Can't follow non-constant source. Use a directive
to specify location instead. You may be seeing SC1091 because
ShellCheck tried to be helpful and strip a leading dynamic path element
as described on that page.

If you're fine with it, ignore the message with a directive.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`source` frames :OAn initial file sourcing a second file, which in turn sources a third file, which in turn sources a fourth file, ...., which in turn sources a 100th file.

Anything but that.

ShellCheck found a chain of 100+ files sourcing each other. It assumed there must be some internal bug, so it stopped.

If this is intentional, you can cosmetically ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`source mylib````
# shellcheck disable=SC1094
source mylib
```
(or fix `mylib`)

ShellCheck encountered a parsing error in a sourced file,
`mylib` in the example.

Fix parsing error, or just disable it with a directive.

If the file is fine and this is due to a known
`shellcheck` bug, you can ignore it with a directive as in the example.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
function foo{
 echo "hello world"
}
```
Prefer POSIX syntax:

```
foo() {
 echo "hello world"
}
```
Alternatively, add the missing space between function name and
opening `{`:

```
# v-- Here
function foo {
 echo "hello world"
}
```
When using `function` keyword function definitions without
`()`, a space is required between the function name and the
opening `{`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`==`. For assignment, use `=`. For comparison, use
`[`/`[[`.`var==value`Assignment:

`var=value`Comparison:

`[ "$var" = value ]`ShellCheck has noticed that you're using `==` in an
unexpected way. The two most common reasons for this is:

You wanted to assign a value but accidentally used
`==` instead of `=`.

You wanted to compare two values, but neglected to use
`[ .. ]` or `[[ .. ]]`.

If you wanted to assign a literal equals sign, use quotes to make this clear:

`var="=sum(A1:A10)"`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`eval`, e.g.
`eval "a=(b)"`.`eval $var=(a b)``eval "$var=(a b)"`Shells differ widely in how they handle unescaped parentheses in
`eval` expressions.

`eval foo=bar` is allowed by dash, bash and ksh.`eval foo=(bar)` is allowed by bash and ksh, but not
dash.`eval $var=(bar)` is allowed by ksh, but not bash or
dash.`eval foo() ( echo bar; )` is not allowed by any
shell.Since the expression is evaluated as shell script code anyways, it should be passed in as a literal string without relying on special case parsing rules in the target shell. Quote/escape the characters accordingly.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`#`.```
while sleep 1
do# show time
 date
done
```
```
while sleep 1
do # show time
 date
done
```
ShellCheck has noticed that you have a keyword immediately followed
by a `#`. In order for the `#` to start a comment,
it needs to come after a word boundary such as a space.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[[ 3 –gt 2 ]] # Uses unicode en-dash character``[[ 3 -gt 2 ]] # Uses regular ASCII hyphen-minus character`A character that looks similar to `-` has made its way
into your code. This is usually due to copy-pasting from blogs and other
websites that formatted code as text, replacing the ASCII hyphen-minus
with a Unicode dash character.

To fix it, simply delete and retype it.

For a large script, you can use your editor's Search&Replace by copy-pasting the bad dash.

None. If you want a literal Unicode dash character, just quote it.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`\` to break line (or use quotes for
literal space).```
# There are spaces after the backslash:
echo hello \
 world
```
```
# No spaces after the backslash:
echo hello \
 world
```
To break a line you can use `\` before the line break.
However, if there are spaces after the backslash, the escape will apply
to them instead of the line break, and the command will not continue on
the next line.

Delete the trailing spaces to make the line break work correctly.

If you do want a literal escaped space at the end of a line you can ignore this error, but please reconsider and use quotes instead. Trailing whitespace is invisible and frequently stripped on purpose (by editor settings / precommits) or accident (copy-paste), and so should not be relied upon for correctness.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$((` differently or not at all. For
`$(command substitution)`, add space after `$(` .
For `$((arithmetics))`, fix parsing errors.`echo "$((cmd "$@") 2>&1)"``echo "$( (cmd "$@") 2>&1)"`You appear to be using `$((` with two (or more)
parentheses in a row, where the first `$(` should open a
subshell.

This is an ill-defined structure that is parsed differently between different shells and shell versions. Prefer adding spaces to make it unambiguous, both to shells and humans.

Consider the `$(((` in `$(((1)) )`:

Ash, dash and Bash 1 parses it as `$(( (` and subsequently
fail to find the matching `))`. Zsh and Bash 2+ looks ahead
and parses it as `$( ((`. Ksh parses it as
`$( ( (`.

**Alternatively**, you may indeed have correctly spaced
your parentheses, but ShellCheck failed to parse `$((` as an
arithmetic expression while accidentally succeeding in parsing it as
`$(` + `(`.

In these cases, double check the syntax to ensure ShellCheck can
parse the `$((`, or ignore this error and hope that it won't
affect analysis too severely.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`sh` or `bash`.```
# shellcheck shell=zsh
export PAGER=less
```
Any supported shell on the shebang or the `-s` option

```
# shellcheck shell=sh
export PAGER=less
```
Shellcheck only supports a specific range of shell dialects, there are many more applications providing shell like experiences and some of them look and feel like POSIX shell or bash but does not support the same commands.

One notable unsupported shell type is zsh, see issue #809 about supporting zsh - some efforts have been done in the past.

The supported shell types are listed in the help context, at the moment these are

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`#!`, not just
`!`, for the shebang.```
!/bin/sh
echo "Hello"
```
```
#!/bin/sh
echo "Hello"
```
You appear to be specifying an interpreter in a shebang, but it's
missing the hash part. The shebang must always start with
`#!`.

Even the name "shebang" itself comes from "hash" (`#`) +
"bang" (`!`).

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

((

(

See SC1102, the similar warning for ambiguous $((.

$((

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`<` instead of `-lt`Similarly, `>` instead of `-gt`,
`<=` instead of `-le`, `>=`
instead of `-ge`, `==` instead of
`-eq`, `!=` instead of `-ne`.

```
if (( 2 -lt 3 ))
then
 echo "True"
fi
```
```
if (( 2 < 3 ))
then
 echo "True"
fi
```
The comparators `-lt`, `-ge`, `-eq`
and friends are flags for the `test` command aka
`[`. You are instead using it in an arithmetic context, such
as `(( .. ))` or `$(( .. ))`, where you should be
using `<`, `>=`, `==` etc
instead.

In arithmetic contexts, `-lt` is simply interpreted as
"subtract the value of `$lt`", which is clearly not the
intention.

If you do want to subtract `$lt` you can add a space to
make this clear to ShellCheck: `echo $((3 - lt))`

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
# shellcheck foobar=baz
echo "Hello World"
```
Depends on your intention.

ShellCheck doesn't recognize the directive
you're trying to use in a `# shellcheck` comment. See the Directives page for supported directives.

It could be misspelled, or you could be using an older version of shellcheck that doesn't support it yet.

None. If you wish to ignore this warning and continue without it, you need version 0.4.5 (commit 88c56ec) or later and a command grouping:

```
# Ignore an unrecognized directive in 0.4.5 or later:
# shellcheck disable=SC1107
{
 # shellcheck unrecognized=directive
 echo "Hello World"
}
```
Before 0.4.5, unrecognized directives are considered parse errors.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`=` .`[ "$var"= 2 ]``[ "$var" = 2 ]`You appear to be missing the space on the left side of the operator.
Shell in general, and `[` in particular, is space sensitive.
Operators and operands must be separate tokens.

Please ensure that the operator, like the `=` in the
example, has a space both before and after it.

None. If you're comparing values in C style reverse order like
`[ -eq == $1 ]`, use quotes: `[ "-eq" == "$1" ]`.
Also, it's pointless since `[ a = b ]` doesn't assign.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`foo && bar``foo && bar`There is an unquoted HTML entity, such as `&`,
`>` or `<` (instead of
`&`, `>` and `<`) in your
code. This usually happens when copy-pasting from a web site that has
mismanaged its code formatting.

You should go through the entire script and replace HTML entities with their corresponding characters.

Don't rely on ShellCheck to detect all of them. ShellCheck only warns about certain cases in certain contexts, while this issue tends to affect the entire script.

If you want to run a command called `amp` after
backgrounding another command, add a space:
`foo & amp;`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo ‘hello world’``echo 'hello world'`Some software, like OS X, Word and WordPress, may automatically replace your regular quotes with slanted Unicode quotes. Try deleting and retyping them, and/or disable “smart quotes” in your editor or OS.

If you want to use typographic single quotes, you can put them in double quotes (or typographic double quotes in single quotes) to make shellcheck ignore them, e.g.,

`printf "Warning: ‘wakeonlan’ is not installed.\n"`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo "hello world”``echo "hello world"`Some software, like OS X, Word and WordPress, may automatically replace your regular quotes with slanted Unicode quotes. The shell does not recognize these quotes and will not respect them.

In this case, you have slanted double quotes in a double quoted string. Try deleting and retyping them, and/or disable “smart quotes” in your editor or OS.

If you want to use literal slanted double quotes for typographic reasons, you can put them in single quotes to make ShellCheck ignore them:

`printf 'Warning: “wakeonlan” is not installed.\n'`Alternatively, use single slanted Unicode quotes like so:

`printf "Warning: ‘wakeonlan’ is not installed.\n"`You can also just ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo 'hello world’``echo 'hello world'`Some software, like OS X, Word and WordPress, may automatically replace your regular quotes with slanted Unicode quotes. The shell does not recognize these quotes and will not respect them.

In this case, you have slanted single quotes in a single quoted string. Try deleting and retyping them, and/or disable “smart quotes” in your editor or OS.

If you want to use literal slanted single quotes for typographic reasons, you can put them in double quotes to make ShellCheck ignore them:

`printf "Warning: ‘wakeonlan’ is not installed.\\n"`You can also just ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`#!`, not just
`#`, for the shebang.```
# /bin/bash
echo "Hello World"
```
```
#! /bin/bash
echo "Hello World"
```
You appear to be specifying a shebang, but missing the bang (i.e.
`!`). The shebang should always be on the form
`#!/path/shell`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
 #!/bin/sh
echo "Hello world"
```
```
#!/bin/sh
echo "Hello World"
```
The script has leading spaces before the shebang (`#!`).
This is not allowed.

The `#!` should be the first two bytes in the file, as
they're used as a file signature by the OS to determine whether a file
is a script.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`#` and `!` in the shebang.```
# !/bin/sh
echo "Hello World"
```
```
#!/bin/sh
echo "Hello World"
```
The script has spaces between the `#` and `!`
in the shebang. This is not valid.

Remove the spaces so the OS can correctly recognize the file as a script.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$` on a `$((..))` expression? (or use
`( (` for arrays).`var=((foo+1))``var=$((foo+1))`You appear to be missing the `$` on an assignment from an
arithmetic expression `var=$((..))` .

Without the `$`, this is an array expression which is
either nested (ksh) or invalid (bash).

If you are trying to define a multidimensional Ksh array, add spaces
between the `( (` to clarify:

`var=( (1 2 3) (4 5 6) )`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`"\n"`. Prefer explicit escaping:
`"\\n"`.*Note: this warning has been retired due to being too pedantic.
Removed in v0.6.0
- 2018-12-03*

`printf "%s\n" "Hello"``printf "%s\\n" "Hello"`or alternatively, with single quotes:

`printf '%s\n' "Hello"`In a double quoted string, you have escaped a character that has no special behavior when escaped. Instead, it's invoking the fallback behavior of being interpreted literally.

Instead of relying on this implicit fallback, you should escape the backslash explicitly. This makes it clear that it's meant to be passed as a literal backslash in the string parameter.

None. This is a stylistic issue which can be ignored. But can you name the 5 characters that
*are* special when escaped in double quotes?

They are $, `, ", \, or newline. More infos are available in the bash manual.

This warning is no longer emitted as of d8a32da07 (strictly after v0.5).

The number of harmlessly affected `printf`,
`sed` and `grep` statements was significantly
higher than the number of actual unexpanded escape sequences. It may
return some day under a `-pedantic` type flag.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

"▭" below indicates an otherwise invisible space:

```
cat << "eof"
Hello
eof▭
```
```
cat << "eof"
Hello
eof
```
The end token of your here document has trailing whitespace. This is invisible to the naked eye, but shells do not accept it.

Remove the trailing whitespace.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`)`.```
var=$(fmt -s "$COLUMNS" << "eof"
This is a bunch of text
eof)
```
```
var=$(fmt -s "$COLUMNS" << "eof"
This is a bunch of text
eof
)
```
When embedding a here document in `$(..)` or
`(..)`, there needs to be a linefeed (newline) between the
here doc token and the closing `)`. Please insert one.

Failing to do so may cause warnings like this:

```
bash: warning: here-document at line 15 delimited by end-of-file (wanted `eof')`
dash: 5: Syntax error: end of file unexpected (expecting ")")
```
This error may be incorrectly emitted for `ksh`, where
this is allowed. In this case, please either write it in a standard way
or ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
cat << eof # --- Start greeting --
Hello
eof # --- End greeting ---
```
```
cat << eof # --- Start greeting --
Hello
eof
 # --- End greeting ---
```
The terminator token for a here document must be on an entirely separate line. No comments are allowed on this line.

Place the comment somewhere else, such as on the following line.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`;`/`&` terminators (and other syntax) on the
line with the `<<`, not here.```
sudo bash -s << "END"
 cmd1
 cmd2
END &
```
```
sudo bash -s << "END" &
 cmd1
 cmd2
END
```
You are using `&`, `;`,
`&>` or similar after a here document. This is not
allowed.

This should instead be part of the line that initiated the here
document, i.e. the one with the `<<`.

If it helps, look at `<< "END"` as if it was
`< file`, and make sure the resulting command is valid.
This is what the shell does. You can then append here document data
after the command.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`<<`.```
cat << EOF
Hello
EOF | nl
```
```
cat << EOF | nl
Hello
EOF
```
You have a here document, and appear to have added text after the terminating token.

This is not allowed. If it was meant to continue the command, put it
on the line with the `<<`.

If it helps, look at << "END" as if it was < file, and make sure the resulting command is valid. This is what the shell does. You can then append here document data after the command.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`if`, not e.g. individual `elif` branches.```
if [ "$prod" = "true" ]
then
 echo "Prod mode"
# shellcheck disable=2154
elif [ "$debug" = "true" ]
then
 echo "Debug mode"
fi
```
```
# Applies to entire `if...fi` command
# shellcheck disable=2154
if [ "$prod" = "true" ]
then
 echo "Prod mode"
elif [ "$debug" = "true" ]
then
 echo "Debug mode"
fi
```
or

```
if [ "$prod" = "true" ]
then
 echo "Prod mode"
elif # Applies only to this [ .. ] command
 # shellcheck disable=2154
 [ "$debug" = "true" ]
then
 echo "Debug mode"
fi
```
You appear to have put a directive before a non-command keyword, such
as `elif`, `else`, `do`,
`;;` or similar.

Unlike many other linters, ShellCheck comment directives apply to the next shell command, rather than to the next line of text.

This means that you can put a directive in front of a
`while` loop, `if` statement or function
definition, and it will apply to that entire structure.

However, it also means that you can not apply the directive to
non-commands like an individual `elif` or `else`
block since these are not commands by themselves, and rather just parts
of an `if` compound command.

Please move the directive in front of the nearest applicable command
that contains the code you want to apply it to, such as before the
`if`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`case` statements, not individual case branches.```
case $? in
 0) echo "Success" ;;
 # shellcheck disable=2154
 *) echo "$cmd $flag returned failure" ;;
esac
```
```
# Applies to everything in the `case` statement
# shellcheck disable=2154
case $? in
 0) echo "Success" ;;
 *) echo "$cmd $flag returned failure" ;;
esac
```
or

```
case $? in
 0) echo "Success" ;;
 *)
 # Applies to a single command within the `case`
 # shellcheck disable=2154
 echo "$cmd $flag returned failure"
 ;;
esac
```
You appear to have put a directive before a branch in a case statement.

ShellCheck directives can not be scoped to individual branches of
`case` statements, only to the entire `case`, or
to individual commands within it. Please move the directive as
appropriate.

(It is possible to apply directives to all commands within a
`{ ..: }` command group, if you truly wish to apply a
directive to multiple commands but not the full `case`
statement.)

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`key=value` pair in directive`# shellcheck disable=SC2153 (variable not a misspelling)``# shellcheck disable=SC2153 # variable not a misspelling`A comment at the end of a directive must be preceded by a
`#` to avoid it being interpreted as an instruction. The directive page contains
more guidance about commenting style.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`var=1 # shellcheck disable=SC2034````
# shellcheck disable=SC2034
var=1
```
ShellCheck expects directives to come before the relevant command. They are not allowed after.

If this is not a directive and just a comment mentioning ShellCheck, please rewrite or capitalize:

`var=1 # ShellCheck encourages lowercase variable names`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`#` in sh.```
// This is a comment.
/* This too. */
```
```
# This is a comment.
# This too.
```
ShellCheck found what appears to be a C-style comment, a line
starting with `//` or `/*`.

In Bourne based shell scripts, the comment character is
`#`

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
# Copyright 2018 Foobar, All rights reserved
#!/bin/bash
```
```
#!/bin/bash
# Copyright 2018 Foobar, All rights reserved
```
A shebang only has an effect when it appears as the first line in a
script. Specifically, the first two bytes of the file must be
`#!`.

Adding comments, copyright notices or simply an accidental blank line
before it will turn a shebang into an ineffectual comment. This means
that the script is no longer in charge of its own interpreter, and may
fail to run or produce different results depending on the context it's
run (e.g. it may work from `bash` but not from
`zsh` or via `sudo`).

Delete any leading blank lines, and move all comments after the shebang.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`!`.```
while! [ -f file ]
do sleep 1; done
```
```
while ! [ -f file ]
do sleep 1; done
```
ShellCheck found a keyword immediately followed by a `!`.
There needs to be a space between them.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
until make
do:; done
```
```
until make
do :; done
```
ShellCheck found a keyword immediately followed by a `:`.
`:` is a synonym for `true`, the command that
"does nothing, successfully", and as a command name it needs a
space.

`do:` is as invalid as `dotrue`. Use
`do :`, or preferably, `do true` for
readability.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`elif` to start
another branch.```
if false
then
 echo "hi"
elseif true
then
 echo "ho"
fi
```
```
if false
then
 echo "hi"
elif true
then
 echo "ho"
fi
```
ShellCheck noticed that you appear to be using `elseif` or
`elsif` as a keyword. This is not valid.

`sh` instead uses `elif` as its alternative
branch keyword.

If you have made your own function called `elseif` and
intend to call it, you can ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`&` terminates the command. Escape it or add space after
`&` to silence.`curl https://www.google.com/search?q=cats&tbm=isch``curl "https://www.google.com/search?q=cats&tbm=isch"`An unescaped and unquoted `&` terminates the command,
but here it's used in the middle of what would otherwise be a shell
word. This most commonly happens when copying a URL with query string
parameters without escaping the `&`.

Either quote or escape the `&` if you wanted it as a
literal ampersand, or add a space after it to make it easier to see
where the previous command stopped.

If you do want to background one command and run another, e.g.
`sleep 10&wait`, just add a space or linefeed after the
`&` to make this more obvious:
`sleep 10& wait`

This does not change the meaning of the script, it just makes it
clear to ShellCheck (and other humans) that the `&` isn't
supposed to be a part of the shell world.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`|`/`||`/`&&` should be at the
end of the previous one.```
dmesg
 | grep "error"
```
```
dmesg |
 grep "error"
```
ShellCheck has found a line that unexpectedly started with
`|`, `||` or `&&`. This usually
happens when a line is broken incorrectly.

When breaking around a `|`, `||` or
`&&`, there are two options:

`dmesg` is a
complete command by itself, but `dmesg |` is not so the shell
knows to continue on the next line.`\` at the end of the previous line to explicitly
tell the shell to continue on the next.In v0.7.2 and below, this warning triggered incorrectly when starting
a line with `&>`. In these versions, you can either ignore the warning, or move the redirection after the
command name.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

shellcheckrc

SC1134 (error): Failed to process foo, line bar: Fix any mentioned problems and try again.

Need more information.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$` literal. Instead of
`"It costs $"5`, use `"It costs \$5"````
echo "The apples are $""1 each"
eval "var=$"name
```
```
echo "The apples are \$1 each"
eval "var=\$name"
# or better yet: var="${!name}"
```
The script appears to be closing a double quoted string for the sole
purpose of making a dollar sign `$` literal.

While this happens to work, the better solution is instead to escape it with a backslash. This allows the double quoted string to continue uninterrupted, thereby reducing the visual noise of stopping and starting quotes in the middle of a shell word.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`]`. Missing
semicolon/linefeed?*Note: Removed in V0.7.2
2021-04-20*

```
if [ -e "foo.txt" ]: then
 echo "Exists"
fi
```
```
if [ -e "foo.txt" ]; then
 echo "Exists"
fi
```
ShellCheck found unexpected characters after the `]` or
`]]` in a `test` expression. In the example, a
colon was accidentally used instead of a semicolon.

Similarly, a missing space before a comment
(`[ -e foo ]#comment`), an additional square bracket
(`[[ -e foo ]]]`), or a missing semicolon before a
`then` on the same line (`if [ foo ]then`) can
cause this warning.

Make sure the `]` or `]]` is not immediately
followed by another shell word character.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`(` to start arithmetic for ((;;)) loop```
for (i=0; i<10; i++))
do
 echo $i
done
```
```
for ((i=0; i<10; i++))
do
 echo $i
done
```
ShellCheck found an arithmetic `for ((;;))` expression
where either the `((` or the `))` did not come as
a pair. Make sure to use `(( ))` and not
`( )`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`for( (i=0; i<10; i++) ); do echo $i; done``for((i=0; i<10; i++)); do echo $i; done`ShellCheck finds arithmetic for ((;;)) expressions where (( or )) are intervening with spaces

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`||`
instead of `-o` between test commands.And variations, like "Use `&&` instead of
`and`".

```
if [ "$1" = "-v" ] -o [ -n "$VERBOSE" ]
then
 echo "Verbose log"
fi
```
```
if [ "$1" = "-v" ] || [ -n "$VERBOSE" ]
then
 echo "Verbose log"
fi
```
You have a `[ .. ]` or `[[ .. ]]` test
expression followed by `-o`/`-a` (or by
Python-style `or`/`and`).

`-o` and `-a` work *inside*
`[ .. ]`, but they do not work *between* them. The
Python operators `or` and `and` are never
recognized in Bash.

To join two separate test expressions, instead use `||`
for "logical OR", or `&&` for "logical AND".

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`&&`/`||`, or bad expression?`[ "$1" ] input="$1"``[ "$1" ] && input="$1"`ShellCheck found characters (other than redirections) after the
`]` or `]]` in a test expression. This is not
valid.

This sometimes happens when there was an additional expression or
command, but joining `||` or `&&` is
missing. Alternatively, it could happen due to typos (like
`[[ $1 ]]]` with an extra `]`), or generally from
malformed test expressions.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`;`/`&&`/`||`/`|`?```
while echo "$2"; do true; done \
 head -n "$1"
while sleep 1; do date; done > my file
```
```
while echo "$2"; do true; done \
 | head -n "$1"
while sleep 1; do date; done > "my file"
```
ShellCheck found unexpected trailing characters after a compound command.

The only things allowed after compound commands are redirections,
shell keywords, and the various command separators (`;`,
`&`, `|`, `&&`,
`||`).

In the first example, a `|` was missing, causing
`head` to appear as an unexpected trailing word, instead of
being piped to. In the second example, a lack of quoting caused
`file` to appear as an unexpected trailing word, instead of
being part of the redirection.

Examine your statement and correct the problem.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`done < <(cmd)` to redirect from process substitution
(currently missing one `<`).```
sum=0
while IFS="" read -r n
do
 (( sum += n ))
done <(file)
```
```
sum=0
while IFS="" read -r n
do
 (( sum += n ))
done < <(file)
```
ShellCheck found a `done` keyword followed by a process
substitution, e.g. `done <(cmd)`.

The intention was most likely to redirect from this process
substitution, in which case you will need one extra `<`:
`done < <(cmd)`.

This is because `<(cmd)` expands to a filename (e.g.
`/dev/fd/63`), and you need a `<` to redirect
from filenames.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
sed \
 -e "s/HOST/$HOSTNAME/g" \
# -e "s/USER/$USER/g" \
 -e "s/ARCH/$(uname -m)/g" \
 "$buildfile"
```
```
sed \
 -e "s/HOST/$HOSTNAME/g" \
 -e "s/ARCH/$(uname -m)/g" \
 "$buildfile"
# This comment is moved out:
# -e "s/USER/$USER/g" \
```
or using backticked, inlined comments:

```
sed \
 -e "s/HOST/$HOSTNAME/g" \
`# -e "s/USER/$USER/g"` \
 -e "s/ARCH/$(uname -m)/g" \
 "$buildfile"
```
(ShellCheck recognizes this idiom and does not suggest quotes or
`$()`, neither of which would have worked)

ShellCheck found a line continuation followed by a commented line that appears to try to do the same.

Backslash line continuations are not respected in comments, and the line instead simply terminates. This is a problem when commenting out one line in a multi-line command like the example.

Instead, either move the line away from its statement, or use an
``# inline comment`` in an unquoted backtick command
substitution.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`external-sources`
can only be enabled in .shellcheckrc, not in individual files.```
# shellcheck external-sources=true
source /dev/zero
```
Add `external-sources=true` to
`.shellcheckrc`

Due to its origins as an online tool, ShellCheck will by default run in a sandbox where it only reads the files explicitly named on the command line.

The `external-sources` directive
allows disabling this, but must be specified in
`.shellcheckrc`. This is because the sandbox would be useless
if the sandboxed script can disable it for itself.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`external-sources` value. Expected
`true`/`false`.`.shellcheckrc`:

`external-sources=maybe``external-sources=true`The `external-sources` directive
expects a value of `true` or `false`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`${#variable}` instead`${#variable}` will be equal to the number of characters
in `"${variable}"`

This is the same result as
`"$( echo "$variable" | wc -m )"` When "$variable" only
contains single-byte characters, it's also the same as
`"$( echo "$variable" | wc -c )"`

```
#!/usr/bin/env bash
if [ "$( echo "$1" | wc -c )" -gt 1 ]; then
 echo "greater than 1"
fi
if [ "$( echo "$1" | wc -m )" -gt 1 ]; then
 echo "greater than 1"
fi
if [ "${#1}" -gt 1 ]; then
 echo "greater than 1"
fi
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`${variable//search/replace}` instead.`string="bwol" ; echo "$string" | sed -e "s/wo/ow/"``string="bwol" ; echo "${string//wo/ow}"`Here is a demonstration of the different search/replace available in Bash:

```
var="foo foo"
# the following two echo's should be equivalent:
echo "$var" | sed 's/^foo/bar/g'
echo ${var/#foo/bar}
```
| $var | sed expression | bash equivalent | result |
|---|---|---|---|
| foo foo | s/foo/bar/ | ${var/foo/bar} | bar foo |
| foo foo | s/foo/bar/g | ${var//foo/bar} | bar bar |
| foo foo | s/^foo/bar/ | ${var/#foo/bar} | bar foo |
| -foo foo | s/^foo/bar/ | ${var/#foo/bar} | -foo foo |
| -foo foo | s/foo$/bar/ | ${var/%foo/bar} | -foo bar |
| -foo foo- | s/foo$/bar/ | ${var/%foo/bar} | -foo foo- |

Let's assume somewhere earlier in your code, you have put data into a variable (Ex: $string). Now you want to search and replace inside the contents of $string and echo the contents out. You could pass this to sed as done in the example above, but for simple substitutions, a parameter expansion can do it with less overhead.

Occasionally a more complex sed substitution is required. For example, getting the last character of a string.

`string="bwol" ; echo "$string" | sed -e "s/^.*\(.\)$/\1/"`This is a bit simple for the example, and there are alternative ways of doing this in the shell, but this SC2001 flags on several of my crazy complex sed commands beyond this example's scope. Utilizing some of the more complex capabilities of sed is required occasionally, and it is safe to ignore SC2001.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`cmd < file | ..` or
`cmd file | ..` instead.`useless-use-of-cat`This is an optional rule, which means that it has a special "long name" and is not enabled by default. See the optional page for more details. In short, you have to enable it with the long name instead of the "SC" code like you would with a normal rule:

`.shellcheckrc``enable=useless-use-of-cat # SC2002`(Note that this check was not optional prior to version 0.11.0.)

```
cat file | tr ' ' _ | nl
cat file | while IFS= read -r i; do echo "${i%?}"; done
```
```
< file tr ' ' _ | nl
while IFS= read -r i; do echo "${i%?}"; done < file
```
`cat` is a tool for con"cat"enating files. Reading a
single file as input to a program is considered a Useless
Use Of Cat (UUOC).

It's more efficient and less roundabout to simply use redirection.
This is especially true for programs that can benefit from seekable
input, like `tail` or `tar`.

Many tools also accept optional filenames, e.g.
`grep -q foo file` instead of
`cat file | grep -q foo`.

Pointing out UUOC is a long standing shell programming tradition, and removing them from a short-lived pipeline in a loop can speed it up by 2x. However, it's not necessarily a good use of time in practice, and rarely affects correctness. Ignore as you see fit.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$((..))`,
`${}` or `[[ ]]`.```
i=$(expr 1 + 2)
l=$(expr length "$var")
```
```
i=$((1+2))
l=${#var}
```
**WARNING:** constants with a leading 0 are interpreted
as octal numbers by bash, but not by expr. Then you should specify the
base when a leading zero may occur:

```
$ x=08
$ echo $(expr 1 + $x)
9
$ echo $((1 + $x))
-bash: 1 + 08: value too great for base (error token is "08")
$ echo $((1 + 10#$x))
9
```
See issue #1910

The expr utility has a rather difficult syntax [...] In many cases, the arithmetic and string features provided as part of the shell command language are easier to use than their equivalents in expr. Newly written scripts should avoid expr in favor of the new features within the shell.

`sh` doesn't have a great replacement for the
`:` operator (regex match). ShellCheck tries not to warn when
using `expr` with `:`, but e.g.
`op=:; expr string "$op" regex` will still trigger it.

Bash has `[[ string =~ regex ]]`, so not even
`expr .. : ..` is necessary.

Bash doesn't check for integer overflow. When dealing with extremely
large numbers, it may be necessary to use an external tool like
`expr`.

Other than that, all uses of `expr` can be rewritten to
use modern shell features instead.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$`/`${}`
is unnecessary on arithmetic variables.`echo $(($n + ${arr[i]}))``echo $((n + arr[i]))`The `$` or `${..}` on regular variables in
arithmetic contexts is unnecessary, and can even lead to subtle bugs.
This is because the contents of `$((..))` is first expanded
into a string, and then evaluated as an expression:

```
$ a='1+1'
$ echo $(($a * 5)) # becomes 1+1*5
6
$ echo $((a * 5)) # evaluates as (1+1)*5
10
```
The `$` is unavoidable for special variables like
`$1` vs `1`, `$#` vs `#`.
It's also required when adding modifiers to parameters expansions, like
`${#var}` or `${var%-}`. ShellCheck does not warn
about these cases.

The `$` is also required (and not warned about) when you
need to specify the *base* for a variable value:

```
$ a=09
$ echo $((a + 1)) # leading zero forces octal interpretation
bash: 09: value too great for base (error token is "09")
$ echo $((10#a + 1))
bash: 10#a: value too great for base (error token is "10#a")
$ echo $((10#$a + 1))
10
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo`? Instead of `echo $(cmd)`, just use
`cmd``echo "$(whoami)"``whoami`ShellCheck found the unnecessary construct
`echo "$(somecommand here)"`.

This is generally due to a misunderstanding about what
`echo` does. It has no role in "showing on screen" or
similar, but simply writes a string to standard output. This is also how
all other programs output data.

`echo "$(somecommand)"` will capture the output
`somecommand` writes to standard output and write it to
standard output, where it was already going. At best this is a no-op,
but it may have several other negative effects:

`echo "$(find . -name '*.iso')" | xargs sha1sum` which does
not allow iterating files and checksumming at the same time. Similarly,
users don't see incremental updates as programs run.`-n`, stripping NUL bytes and trailing
linefeeds, and expanding escape sequences in some shells but not
others.`echo "$(grep '^user:' /etc/passwd)"` no longer returns with
failure when the user is not found.`ls` vs `echo "$(ls)"` where the former
outputs columns and colors according to user preferences, while the
latter doesn't.To avoid all this, simply replace `echo "$(somecommand)"`
with `somecommand` as in the example. It's shorter, faster,
and more correct.

If you are relying on one of the otherwise detrimental effects for correctness, you can consider one of:

```
# Suppress exit code without the other negative effects
cmd || true
# Disable tty specific output without the other negative effects
cmd | cat
# Buffer up potentially large output without using more memory or modifying the content in any way
cmd > file.tmp
cat file.tmp
# Exactly like `echo "$(cmd)"`, but allows output like `-n` and works the same across shells
printf '%s\n' "$(cmd)"
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$(...)` notation instead of legacy backticked
``...``.`echo "You are running on `uname`"``echo "You are running on $(uname)"`Backtick command substitution ``...`` is legacy syntax
with several issues.

`$(...)` command substitution has none of these problems,
and is therefore strongly encouraged.

Note: The `$(...)` syntax was introduced in the 1989 Korn
Shell (ksh). Finally, in 2011, Solaris 11 was the last operating system
to switch from the Bourne Shell to the Korn Shell. After 2011, all
typical shells have supported the POSIX `$(...)`
notation.

`$(...)` and might require the use of backtick command
substitution. See [mc-devel]
[PATCH] Prefer $() to backticks in sh script and follow-ups.`$(...)` preferred over ``...``
(backticks)?``command`` in
shell programming?ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$((..))` instead of
deprecated `$[..]`.```
n=1
n=$[n+1]
```
```
n=1
n=$((n+1))
```
The `$[..]` syntax was deprecated in Bash 2.0 and replaced
with the standard `$((..))` syntax from Korn shell

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo`
doesn't read from stdin, are you sure you should be piping to it?*Note: Removed in v0.4.7
- 2020-03-07*

`find . | echo``find .`You are piping command output to `echo`, but
`echo` ignores all piped input.

In particular, `echo` is not responsible for putting
output on screen. Commands already output data, and with no further
actions that will end up on screen.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`pgrep` instead of grepping `ps`
output.`ps ax | grep -v grep | grep "$service" > /dev/null``pgrep -f "$service" > /dev/null`If you are just after a pid from a running program, then pgrep is a much safer alternative. Especially if you are also looking for a pid belonging to a certain user or group. All of the parameters are in one command and it can eliminate multiple greps, cuts, seds, awks, etc.

If you want a field that's not the pid, consider doing this through
`ps` + `pgrep` instead of `ps` +
`grep`:

```
for pid in $(pgrep '^python$')
do
 user=$(ps -o user= -p "$pid")
 echo "The process $pid is run by $user"
done
```
This is more robust than `ps .. | grep python | cut ..`
because it does not try to match against unrelated fields, such as if
the user's name was `pythonguru`.

`pgrep` is not POSIX. Please ignore
this warning if you are targeting POSIX userlands.

You can ignore this error if you are trying to
match against something that `pgrep` doesn't support:

```
# pgrep does not support filtering by 'nice' value
# shellcheck disable=SC2009
ps -axo nice=,pid= | grep -v '^ 0'
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`ls | grep`. Use a glob or a for loop with a condition to
allow non-alphanumeric filenames.`ls /directory | grep mystring`or

`rm $(ls | grep -v '\.c$')````
# BASH
shopt -s extglob
rm -- !(*.c)
# POSIX
for f in ./*
do
 case $f in
 *.c) true;;
 *) rm "$f";;
 esac
done
```
Parsing ls is
generally a bad idea because the output is fragile and human readable.
To better handle non-alphanumeric filenames, use a glob. If you need
more advanced matching than a glob can provide, use a `for`
loop.

`ls` has sorting options that are tricky to get right
with other commands. If a specific order of files is needed, ls
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`find -print0` or `find -exec` to better handle
non-alphanumeric filenames.`ls | xargs -n1 wc -w``find . -maxdepth 1 -print0 | xargs -0 -n1 wc -w``find . -maxdepth 1 -exec wc -w {} \;`Using `-print0` separates each output with a NUL
character, rather than a newline, which is safer to pipe into
`xargs`. Alternatively using `-exec` avoids the
problem of piping and parsing filenames in the first place.

See SC2012 for more details on this issue.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`find` instead of `ls` to better handle
non-alphanumeric filenames.`ls -l | grep " $USER " | grep '\.txt$'``NUMGZ="$(ls -l *.gz | wc -l)"``find ./*.txt -user "$USER" # Using the names of the files````
gz_files=(*.gz)
numgz=${#gz_files[@]} # Sometimes, you just need a count
```
`ls` is only intended for human consumption: it has a
loose, non-standard format and may "clean up" filenames to make output
easier to read.

Here's an example:

```
$ ls -l
total 0
-rw-r----- 1 me me 0 Feb 5 20:11 foo?bar
-rw-r----- 1 me me 0 Feb 5 2011 foo?bar
-rw-r----- 1 me me 0 Feb 5 20:11 foo?bar
```
It shows three seemingly identical filenames, and did you spot the
time format change? How it formats and what it redacts can differ
between locale settings, `ls` version, and whether output is
a tty.

`ls` with `find`:(Note that `-maxdepth` is not POSIX, but can be simulated
by having the expression call `-prune` on all directories it
finds, e.g. `find ./* -prune -print`)

`ls` can usually be replaced by `find` if it's
just the filenames, or a count of them, that you're after. Note that if
you are using `ls` to get at the contents of a directory, a
straight substitution of `find` may not yield the same
results as `ls`. Here is an example:

```
$ ls -c1 .snapshot
rnapdev1-svm_4_05am_6every4hours.2019-04-01_1605
rnapdev1-svm_4_05am_6every4hours.2019-04-01_2005
rnapdev1-svm_4_05am_6every4hours.2019-04-02_0005
rnapdev1-svm_4_05am_6every4hours.2019-04-02_0405
rnapdev1-svm_4_05am_6every4hours.2019-04-02_0805
rnapdev1-svm_4_05am_6every4hours.2019-04-02_1205
snapmirror.1501b4aa-3f82-11e8-9c31-00a098cef13d_2147868328.2019-04-01_190000
```
versus

```
$ find .snapshot -maxdepth 1
.snapshot
.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-02_0005
.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-02_0405
.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-02_0805
.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-01_1605
.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-01_2005
.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-02_1205
.snapshot/snapmirror.1501b4aa-3f82-11e8-9c31-00a098cef13d_2147868328.2019-04-01_190000
```
You can see two differences here. The first is that the
`find` output has the full paths to the found files, relative
to the current working directory from which `find` was run
whereas `ls` only has the filenames. You may have to adjust
your code to not add the directory to the filenames as you process them
when moving from `ls` to `find`, or (with GNU
find) use `-printf '%P\n'` to print just the filename.

The second difference in the two outputs is that the
`find` command includes the searched directory as an entry.
This can be eliminated by also using `-mindepth 1` to skip
printing the root path, or using a negative name option for the searched
directory:

```
$ find .snapshot -maxdepth 1 ! -name .snapshot
.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-02_0005
.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-02_0405
.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-02_0805
.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-01_1605
.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-01_2005
.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-02_1205
.snapshot/snapmirror.1501b4aa-3f82-11e8-9c31-00a098cef13d_2147868328.2019-04-01_190000
```
**Note:** If the directory argument to `find`
is a fully expressed path (`/home/somedir/.snapshot`), then
you should use `basename` on the `-name`
filter:

```
$ theDir="$HOME/.snapshot"
$ find "$theDir" -maxdepth 1 ! -name "$(basename $theDir)"
/home/matt/.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-02_0005
/home/matt/.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-02_0405
/home/matt/.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-02_0805
/home/matt/.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-01_1605
/home/matt/.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-01_2005
/home/matt/.snapshot/rnapdev1-svm_4_05am_6every4hours.2019-04-02_1205
/home/matt/.snapshot/snapmirror.1501b4aa-3f82-11e8-9c31-00a098cef13d_2147868328.2019-04-01_190000
```
If trying to parse out any other fields, first see whether
`stat` (GNU, OS X, FreeBSD) or `find -printf`
(GNU) can give you the data you want directly. When trying to determine
file size, try: `wc -c`. This is more portable as
`wc` is a mandatory unix command, unlike `stat`
and `find -printf`. It may be slower as unoptimized
`wc -c` may read the entire file rather than just checking
its properties. On some systems, `wc -c` adds whitespace to
the file size which can be trimmed by double expansion:
`$(( $(wc -c < "filename") ))`

If the information is intended for the user and not for processing
(`ls -l ~/dir | nl; echo "Ok to delete these files?"`) you
can ignore this error with a directive.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`while read`
loop.```
for line in $(cat file | grep -v '^ *#')
do
 echo "Line: $line"
done
```
```
grep -v '^ *#' file | while IFS= read -r line
do
 echo "Line: $line"
done
```
or without a subshell (bash, zsh, ksh):

```
while IFS= read -r line
do
 echo "Line: $line"
done < <(grep -v '^ *#' file)
```
or without a subshell, with a pipe (more portable, but write a file on the filesystem):

```
mkfifo mypipe
grep -v '^ *#' file > mypipe &
while IFS= read -r line
do
 echo "Line: $line"
done < mypipe
rm mypipe
```
NOTE: `grep -v '^ *#'` is a placeholder example and not
needed. To just loop through a file:

```
while IFS= read -r line
do
 echo "Line: $line"
done < file
# or: done <<< "$variable"
```
For loops by default (subject to `$IFS`) read word by
word. Additionally, glob expansion will occur.

Given this text file:

```
foo *
bar
```
The for loop will print:

```
Line: foo
Line: aardwark.jpg
Line: bullfrog.jpg
...
```
The while loop will print:

```
Line: foo *
Line: bar
```
If you do want to read word by word, you can set `$IFS`
appropriately and disable globbing with `set -f`, and then ignore this warning. Alternatively, you can pipe
through `tr ' ' '\n'` to turn words into lines, and then use
`while read`. In Bash/Ksh, you can also use a
`while read -a` loop to get an array of words per line.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`find . -name '*.tar' -exec tar xf {} -C "$(dirname {})" \;``find . -name '*.tar' -exec sh -c 'tar xf "$1" -C "$(dirname "$1")"' _ {} \;`Bash evaluates any command substitutions before the command they
feature in is executed. In this case, the command is `find`.
This means that `$(dirname {})` will run
**before** `find` runs, and not
**while** `find` runs.

To run shell code for each file, we can write a tiny script and
inline it with `sh -c`. We add `_` as a dummy
argument that becomes `$0`, and a filename argument that
becomes `$1` in the inlined script:

```
$ sh -c 'echo "$1 is in $(dirname "$1")"' _ "mydir/myfile"
mydir/myfile is in mydir
```
This command can be executed by `find -exec`, with
`{}` as the filename argument. It executes shell which
interprets the inlined script once for each file. Note that the inlined
script is single quoted, again to ensure that the expansion does not
happen prematurely .

If you don't care (or if you prefer) that it's only expanded once, like when dynamically selecting the executable to be used by all invocations, you can ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`A && B || C` is not if-then-else. C may run
when A is true.`[[ $dryrun ]] && echo "Would delete file" || rm file````
if [[ $dryrun ]]
then
 echo "Would delete file"
else
 rm file
fi
```
It's common to use `A && B` to run `B`
when `A` is true, and `A || C` to run
`C` when `A` is false.

However, combining them into `A && B || C` is not
the same as `if A then B else C`.

In this case, if `A` is true but `B` is false,
`C` will run.

For the code sample above, if the script was run with stdout closed
for any reason (such as explicitly running
`script --dryrun >&-`), echo would fail and the file
would be deleted, even though `$dryrun` was set!

If an `if` clause is used instead, this problem is
avoided.

We can think of the example above as

`((([[ $dryrun ]]) && echo "Would delete file") || rm file)`expressing the left-associativity of the `&&`
`||` operators.

Whenever a command (strictly, a pipeline) succeeds or fails, the
execution proceeds following the next `&&` (for
success) or `||` (for failure). (More strictly, the
parentheses should be replaced with `{ command; }` to avoid
making a subshell, but that's ugly and boring.)

Ignore this warning when you actually do intend to run C when either A or B fails.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
name=World
echo 'Hello $name' # Outputs Hello $name
```
```
name=World
echo "Hello $name" # Outputs Hello World
```
ShellCheck found an expansion like `$var`,
`$(cmd)`, or ``cmd`` in single quotes.

Single quotes express all such expansions. If you want the expression to expand, use double quotes instead.

If switching to double quotes would require excessive escaping of other metacharacters, note that you can mix and match quotes in the same shell word:

`dialog --msgbox "Filename $file may not contain any of: "'`&;"\#%$' 10 70`If you know that you want the expression literally without expansion, you can ignore this message:

```
# We want this to output $PATH without expansion
# shellcheck disable=SC2016
echo 'PATH=$PATH:/usr/local/bin' >> ~/.bashrc
```
```
# We also want this variable to expand "$BASH_SOURCE:$LINE..." during an execution trace.
# shellcheck disable=SC2016
PS4='+$BASH_SOURCE:$LINENO:$FUNCNAME: '
```
```
# We want to control which environment variables envsubst replaces
# shellcheck disable=SC2016
envsubst '${SERVICE_HOST}:${SERVICE_PORT}' config.template > config
```
ShellCheck also does not warn about escaped expansions in double quotes:

`echo "PATH=\$PATH:/usr/local/bin" >> ~/.bashrc`This suggestion is primarily meant to help newbies who assume single
and double quotes are basically the same, like in Python and JavaScript.
It's not at all meant to discourage experienced users from using single
quotes in general. If you are well aware of the difference, please do
not hesitate to permanently disable this suggestion with
`disable=SC2016` in your `.shellcheckrc`.

ShellCheck tries to increase the signal-to-noise ratio of this
warning by ignoring certain well known commands that frequently expect
literal dollar signs, such as `sh` and `perl`.
However, there's a long tail of less common commands and flags that also
frequently expect `$`s, and it's not in ShellCheck's scope to
try to keep track of them all. When you come across such a command,
please ignore the suggestion, either permanently or
for that one instance.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`a/b*c` with `a*c/b`.`percent=$((count/total*100))``percent=$((count*100/total))`If integer division is performed before multiplication, the intermediate result will be truncated causing a loss of precision.

In this case, if `count=1` and `total=2`, then
the problematic code results in `percent=0`, while the
correct code gives `percent=50`.

If you want and expect truncation you can ignore this message.

ShellCheck doesn't warn when `b` and `c` are
identical expressions, e.g. `a/10*10`, under the assumption
that the intent is to rounded to the nearest 10 rather than the no-op of
multiply by `1`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[:lower:]` to support accents and foreign alphabets.`PLATFORM="$(uname -s | tr 'A-Z' 'a-z')"``PLATFORM="$(uname -s | tr '[:upper:]' '[:lower:]')"``A-Z` and `a-z` are commonly intended to mean
"all uppercase" and "all lowercase letters" respectively. This ignores
accented characters in English, and foreign characters in other
languages:

```
$ tr 'a-z' 'A-Z' <<< "My fiancée ordered a piña colada."
MY FIANCéE ORDERED A PIñA COLADA.
```
Instead, you can use `[:lower:]` and
`[:upper:]` to explicitly specify case:

```
$ tr '[:lower:]' '[:upper:]' <<< "My fiancée ordered a piña colada."
MY FIANCÉE ORDERED A PIÑA COLADA.
```
If you don't want `a-z` to match `é` or
`A-Z` to match `Ñ`, you can ignore this
message.

As of 2019-09-08, BusyBox `tr` does not support character
classes, so you would have to ignore this message.

Note that the examples used here are multibyte characters in UTF-8. Many implementations (including GNU) fails to deal with them.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

[:upper:]

See the equivalent warning for lowercase matching: SC2018

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`tr`
replaces sets of chars, not words (mentioned due to duplicates).`echo 'hello world' | tr 'hello' 'goodbye'``echo 'hello world' | sed -e 's/hello/goodbye/g'``tr` is for `tr`ansliteration, turning some
characters into other characters. It doesn't match strings or words,
only individual characters.

In this case, it transliterates h->g, e->o, l->d, o->y, resulting in the string "goddb wbrdd" instead of "goodbye world".

The solution is to use a tool that does string search and replace, such as sed.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[]` around ranges in `tr`, it replaces
literal square brackets.`tr -cd '[a-z]'``tr -cd 'a-z'`Ancient System V `tr` required brackets around operands,
but modern implementations including POSIX, GNU, OS X and *BSD instead
treat them as literals.

Unless you want to operate on literal square brackets, don't include them.

If you do want to replace literal square brackets, reorder the
expression (e.g. `a-z[]` to make it clear that the brackets
are not special).

ShellCheck does not warn about correct usage of `[..]` in
character and equivalence classes like `[:lower:]` and
`[=e=]`.

Busybox requires `CONFIG_FEATURE_TR_CLASSES` for this to
work. If busybox was compiled without this feature, character classes
will obviously not be available, and the correct usage is to use
`a-z`, not `[a-z]`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`o*` here matches `ooo` but not
`oscar`.`grep 'foo*'`when wanting to match `food` and `foosball`,
but not `mofo` or `keyfob`.

`grep '^foo'`As a glob, `foo*` means "Any string starting with foo",
e.g. `food` and `foosball`.

As a regular expression, "foo*" means "f followed by 1 or more o's, anywhere", e.g. "mofo" or "keyfob".

This construct is way more common as a glob than as a regex, so ShellCheck notifies you about it.

If you're aware of the above, you can ignore this message. If you'd
like shellcheck to be quiet, use a directive or
`'fo[o]*'`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`time` as seen in man time(1). Use
`command time ..` for that one.`time -some some``command time -some some``time` is a built-in command. If you would like to use
`time` from `$PATH`, you need to use
`command` to execute it as a regular command.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`sudo`
doesn't affect redirects. Use `..| sudo tee file`or "Use `..| sudo tee -a file`" instead of
`>>` to append.

or "Use `sudo cat file | ..`" instead of `<`
to read.

```
# Write to a file
sudo echo 3 > /proc/sys/vm/drop_caches
# Append to a file
sudo echo 'export FOO=bar' >> /etc/profile
# Read from a file
sudo wc -l < /etc/shadow
```
```
# Write to a file
echo 3 | sudo tee /proc/sys/vm/drop_caches > /dev/null
# Append to a file
echo 'export FOO=bar' | sudo tee -a /etc/profile > /dev/null
# Read from a file
sudo cat /etc/shadow | wc -l
```
Redirections are performed by the current shell before
`sudo` is started. This means that it will use the current
shell's user and permissions to open and read from or write to the
file.

`sudo command < file` with
`sudo cat file | command`.`sudo command > file` with
`command | sudo tee file > /dev/null``tee -a`.The substitutions work by having a command open the file for reading or writing, instead of relying on the current shell. Since the command is run with elevated privileges, it will have access to files that the current user does not.

Note: there is nothing special about `tee`. It's just the
simplest command that can both truncate and append to files without help
from the shell. Here are equivalent alternatives:

Truncating:

```
echo 'data' | sudo dd of=file
echo 'data' | sudo sed 'w file'
```
Appending:

```
echo 'data' | sudo awk '{ print $0 >> "file" }'
echo 'data' | sudo sh -c 'cat >> file'
```
If you want to run a command as root but redirect as the normal user, you can ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`\[..\]` to prevent
line wrapping issues.`PS1='\e[36m\$ \e(B\e[m'``PS1='\[\e[36m\]\$ \[\e(B\e[m\]'`Bash is unable to determine exactly which parts of your prompt are
text and which are terminal codes. You have to help it by wrapping
invisible control codes in `\[..\]` (and ensuring that
visible characters are not wrapped in `\[..\]`).

Note: ShellCheck offers this as a helpful hint and not a robust check. Don't rely on ShellCheck to verify that your prompt is correct.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`'nest '"'single quotes'"'` instead?`alias server_uptime='ssh $host 'uptime -p''``alias server_uptime='ssh $host '"'uptime -p'"`In the first case, the user has four single quotes on a line, wishfully hoping that the shell will match them up as outer quotes around a string with literal single quotes:

```
# v--------match--------v
alias server_uptime='ssh $host 'uptime -p''
# ^--match--^
```
The shell, meanwhile, always terminates single quoted strings at the first possible single quote:

```
# v---match--v
alias server_uptime='ssh $host 'uptime -p''
# ^^
```
Which is the same thing as
`alias server_uptime='ssh $host uptime' -p`.

There is no way to nest single quotes. However, single quotes can be placed literally in double quotes, so we can instead concatenate a single quoted string and a double quoted string:

```
# v--match---v
alias server_uptime='ssh $host '"'uptime -p'"
# ^---match---^
```
This results in an alias with embedded single quotes.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo "You enter "$HOSTNAME". You can smell the wumpus." >> /etc/issue``echo "You enter $HOSTNAME. You can smell the wumpus." >> /etc/issue`Always quoting variables and command expansions is good practice, but blindly putting quotes left and right of them is not.

In this case, ShellCheck has noticed that the quotes around the expansion are unquoting it, because the left quote is terminating an existing double quoted string, while the right quote starts a new one:

```
echo "You enter "$HOSTNAME". You can smell the wumpus."
 |----------| |---------------------------|
 Quoted No quotes Quoted
```
If the quotes were supposed to be literal, they should be escaped. If the quotes were supposed to quote an expansion (as in the example), they should be removed because this is already a double quoted string.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo`
won't expand escape sequences. Consider `printf`.`echo "Name:\t$value"``printf 'Name:\t%s\n' "$value"`Backslash escapes like `\t` and `\n` are not
expanded by echo, and become literal backslash-t, backslash-n.

`printf` does expand these sequences, and should be used
instead.

Other, non-portable methods include `echo -e '\t'` and
`echo $'\t'`. ShellCheck will warn if this is used in a
script with shebang `#!/bin/sh`.

If you actually wanted a literal backslash-t, use

`echo "\\t"`None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`ssh host "echo $HOSTNAME"``ssh host "echo \$HOSTNAME"`or

`ssh host 'echo $HOSTNAME'`Bash expands all arguments that are not escaped/singlequoted. This means that the problematic code is identical to

`ssh host "echo clienthostname"`and will print out the client's hostname, not the server's hostname.

By escaping the `$` in `$HOSTNAME`, it will be
transmitted literally and evaluated on the server instead.

If you do want your string expanded on the client side, you can safely ignore this message.

Keep in mind that the expanded string will be evaluated again on the
server side, so for arbitrary variables and command output, you may need
to add a layer of escaping with e.g. `printf %q`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

See companion warning SC2031.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

There are many ways of accidentally creating subshells, but a common one is piping to a loop:

```
n=0
printf "%s\n" {1..10} | while read i; do (( n+=i )); done
echo $n
```
```
# Bash specific: process substitution. Also try shopts like lastpipe.
n=0
while read i; do (( n+=i )); done < <(printf "%s\n" {1..10})
echo $n
```
In `sh`, temporary files, FIFOs or file descriptors can be
used instead. When the output of the command can be stored to a variable
before entering the loop, here documents are a preferable
alternative.

```
n=0
SUMMANDS="1
2
3
4
5
6
7
8
9
10
"
while read i; do n=$(( n + i )); done <<SUMMANDS_HEREDOC_INPUT
$SUMMANDS
SUMMANDS_HEREDOC_INPUT
echo $n
```
With Bash 4.2+ you can also use `shopt -s lastpipe` which
will change the pipe behaviour to be similar to Ksh and Zsh (see
Rationale below) as
long as job control is not active (e.g. inside a script):

```
#!/usr/bin/env bash
shopt -s lastpipe
n=0
printf "%s\n" {1..10} | while read i; do (( n+=i )); done
echo $n
```
Variables set in subshells are not available outside the subshell. This is a wide topic, and better described on the Wooledge Bash Wiki.

Here are some constructs that cause subshells (shellcheck may not
warn about all of them). In each case, you can replace
`subshell1` by a command or function that sets a variable,
e.g. simply `var=foo`, and the variable will appear to be
unset after the command is run. Similarly, you can replace
`regular` with `var=foo`, and it will be set
afterwards:

Pipelines:

```
subshell1 | subshell2 | subshell3 # Dash, Ash, Bash (default)
subshell1 | subshell2 | regular # Ksh, Zsh, Bash (with lastpipe=on and no job control)
```
Command substitution:

`regular "$(subshell1)" "`subshell2`"`Process substitution:

`regular <(subshell1) >(subshell2)`Some forms of grouping:

```
( subshell )
{ regular; }
```
Backgrounding:

```
subshell1 &
subshell2 &
```
Anything executed by external processes:

```
find . -exec subshell1 {} \;
find . -print0 | xargs -0 subshell2
sudo subshell3
su -c subshell4
```
This applies not only to setting variables, but also setting shell options and changing directories.

You can ignore this error if you don't care that the changes aren't reflected, because work on the value branches and shouldn't be recombined.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

(or from xargs, sudo, chroot, or other commands)

See companion warning SC2033.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
foo() { bar --baz "$@"; frob --baz "$@"; };
find . -exec foo {} +
```
`find . -exec sh -c 'bar --baz "$@"; frob --baz "$@";' -- {} +`Shell functions are only known to the shell. External commands like
`find`, `xargs`, `su` and
`sudo` do not recognize shell functions.

Instead, the function contents can be executed in a shell, either
through `sh -c` or by creating a separate shell script as an
executable file.

If you're intentionally passing a word that happens to have the same name as a declared function, you can quote it to make shellcheck ignore it, e.g.

```
nobody() {
 sudo -u "nobody" "$@"
}
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
foo=42
echo "$FOO"
```
```
foo=42
echo "$foo"
```
Variables not used for anything are often associated with bugs, so ShellCheck warns about them.

Also note that something like `local let foo=42` does not
make a `let` statement local -- it instead declares an
additional local variable named `let`.

This warning may be falsely emitted when a variable is only referenced indirectly, and for variables that are intentionally unused.

It's ShellCheck's intended behavior to emit this warning for any variable that is only referenced though indirection:

```
# foo generates a warning, even though it has five indirect references
foo=42
name=foo
echo "${!name} $((name))"
export "$name"; eval "echo $name"
declare -n name; echo "$name"
```
This is an intentional design decision and not a bug. If you have variables that will not have direct references, consider using an associative array in bash, or just ignore the warning.

Tracking indirect references is a common problem for compilers and static analysis tools, and it is known to be unsolvable in the most general case. There are two ways to handle unresolved indirections (which in a realistic program will essentially be all of them):

Compilers are forced to do the former, but static analysis tools
generally do the latter. This includes `pylint`,
`jshint` and `shellcheck` itself. This is a design
decision meant to make the tools more helpful at the expense of some
noise. For consistency and to avoid giving the impression that it should
work more generally, ShellCheck does not attempt to resolve even trivial
indirections.

For throwaway variables, consider using `_` as a
dummy:

```
read _ last _ zip _ _ <<< "$str"
echo "$last, $zip"
```
Or optionally as a prefix for dummy variables (ShellCheck >0.7.2).

`sh read _first last _email zip _lat _lng <<< "$str" echo "$last, $zip"`

For versions <= 0.7.2, the message can optionally be ignored with a directive:

```
# shellcheck disable=SC2034 # Unused variables left for readability
read first last email zip lat lng <<< "$str"
echo "$last, $zip"
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`./*glob*` or `-- *glob*` so names with dashes
won't become options.`rm *``rm ./*`or

`rm -- *`Since files and arguments are strings passed the same way, programs can't properly determine which is which, and rely on dashes to determine what's what.

A file named `-f` (`touch -- -f`) will not be
deleted by the problematic code. It will instead be interpreted as a
command line option, and `rm` will even report success.

Using `./*` will instead cause the glob to be expanded
into `./-f`, which no program will treat as an option.

Similarly, `--` by convention indicates the end of
options, and nothing after it will be treated like flags (except for
some programs possibly still special casing `-` as e.g.
stdin).

Note that changing `*` to `./*` in GNU Tar
parameters will add `./` prefix to path names in the created
archive. This may cause subtle problems (eg. to search for a specific
file in archive, the `./` prefix must be specified as well).
So using `-- *` is a safer fix for GNU Tar commands.

`echo` and `printf` do not have issues unless
the glob is the first word in the command. ShellCheck 0.7.2+ does not
warn for these commands.

For more information, see "Filenames and Pathnames in Shell: How to do it Correctly".

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`a=$(b | c)` .`sum=find | wc -l``sum=$(find | wc -l)`The intention in this code was that `sum` would in some
way get the value of the command `find | wc -l`.

However, `|` has precedence over the `=`, so
the command is a two stage pipeline consisting of `sum=find`
and `wc -l`.

`sum=find` is a plain string assignment. Since it happens
by itself in an independent pipeline stage, it has no effect: it
produces no output, and the variable disappears when the pipeline stage
finishes. Because the assignment produces no output, `wc -l`
will count 0 lines.

To instead actually assign a variable with the output of a command,
command substitution `$(..)` can be used.

None. This warning is triggered whenever the first stage of a pipeline is a single assignment, which is never correct.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`var=$(cmd)` .`var=grep -c pattern file``var=$(grep -c pattern file)`To assign the output of a command to a variable, use
`$(command substitution)`. Just typing a command after the
`=` sign does not work.

None.

This warning triggers generally for `var=value -flag` and
`var=value *glob*`. See related warning SC2209 which matches
`var=commonCommand`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-print0`/`-0` or `find -exec +` to
allow for non-alphanumeric filenames.`find . -type f | xargs md5sum````
find . -type f -print0 | xargs -0 md5sum
find . -type f -exec md5sum {} +
```
By default, `xargs` interprets spaces and quotes in an
unsafe and unexpected way. Whenever it's used, it should be used with
`-0` or `--null` to split on `\0`
bytes, and `find` should be made to output `\0`
separated filenames.

POSIX does not require find or xargs to support null terminators, so
you can also use `find -exec +`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

*Note: This warning has been retired in favor of individual SC3xxx
warnings for each individual issue. Removed in V0.7.2
- 2021-04-19*

You have declared that your script works with `/bin/sh`,
but you are using features that have undefined behavior according to the
POSIX specification.

It may currently work for you, but it can or will fail on other OS, the same OS with different configurations, from different contexts (like initramfs/chroot), or in different versions of the same OS, including future updates to your current system.

Either declare that your script requires a specific shell like
`#!/bin/bash` or `#!/bin/dash`, or rewrite the
script in a portable way.

For help with rewrites, the Ubuntu wiki has a list of portability
issues that broke people's `#!/bin/sh` scripts when
Ubuntu switched from Bash to Dash. See also Bashism on wooledge's
wiki. ShellCheck may not warn about all these issues.

`$'c-style-escapes'`bash, ksh:

`a=$' \t\n'`POSIX:

`a="$(printf '%b_' ' \t\n')"; a="${a%_}" # protect trailing \n`Want some good news? See http://austingroupbugs.net/view.php?id=249#c590.

`==` operator is not supported in POSIX
`sh`

Bash:

```
if [ "$a" == "$b" ]; then
 echo "equal"
fi
```
POSIX:

```
if [ "$a" = "$b" ]; then
 echo "equal"
fi
```
`$"msgid"`Bash:

`echo $"foo $(bar) baz"`POSIX:

```
. gettext.sh # GNU Gettext sh library
# ...
barout=$(bar)
eval_gettext 'foo $barout baz' # See GNU Gettext doc for more info.
```
Or you can change them to normal double quotes so you go without
`gettext`.

`${var:1}` (substring
expansion)https://wiki.ubuntu.com/DashAsBinSh#A.24.7Bfoo:3.5B:1.5D.7D

`for` loopsBash:

`for ((init; test; next)); do foo; done`POSIX:

```
: $((init))
while [ $((test)) -ne 0 ]; do foo; : $((next)); done
```
Bash:

`printf "%s\n" "$(( 2**63 ))"`POSIX:

The POSIX standard does not allow for exponents. However, you can
replicate them completely built-in using a POSIX compatible function. As
an example, the `pow` function from here.

```
pow() {
 set -- "$1" "$2" 1
 while [ "$2" -gt 0 ]; do
 set -- "$1" $(($2-1)) $(($1*$3))
 done
 # %d = signed decimal, %u = unsigned decimal
 # Either should overflow to 0
 printf "%d\n" "$3"
}
```
To compare:

```
$ echo "$(( 2**62 ))"
4611686018427387904
$ pow 2 62
4611686018427387904
```
Alternatively, if you don't mind using an external program, you can
use `bc`. Be aware though: `bash` and other
programs may abide by a certain maximum integer that `bc`
does not (for `bash` that's: 64-bit signed long int, failing
back to 32-bit signed long int).

Example:

```
# Note the overflow that gives a negative number
$ echo "$(( 2**63 ))"
-9223372036854775808
# No such problem
$ echo 2^63 | bc
9223372036854775808
# 'bc' just keeps on going
$ echo 2^1280 | bc
20815864389328798163850480654728171077230524494533409610638224700807\
21611934672059602447888346464836968484322790856201558276713249664692\
98162798132113546415258482590187784406915463666993231671009459188410\
95379622423387354295096957733925002768876520583464697770622321657076\
83317005651120933244966378183760369413644440628104205339687097746591\
6057756101739472373801429441421111406337458176
```
`((..))`Bash:

```
((a=c+d))
((d)) && echo d is true.
```
POSIX:

```
: $((a=c+d)) # discard the output of the arith expn with `:` command
[ $((d)) -ne 0 ] && echo d is true. # manually check non-zero => true
```
`select` loopsIt takes extra care over terminal columns to make select loop look
like bash's, which generates a list with multiple items on one line, or
like `ls`.

It is, however, still possible to make a naive translation for
`select foo in bar baz; do eat; done`:

```
while
 _i=0 _foo= foo=
 for _name in bar baz; do echo "$((_i+=1))) $_name"; done
 printf '$# '; read _foo
do
 case _foo in 1) foo=bar;; 2) foo=baz;; *) continue;; esac
 eat
done
```
Bash, ksh:

`read aaa bbb <<< $(grep foo bar)`POSIX:

```
read aaa bbb << EOF
$(grep foo bar)
EOF
```
See https://unix.stackexchange.com/tags/echo/info.

`${var/pat/replacement}`Bash:

`echo "${TERM/%-256*}"`POSIX:

```
echo "$TERM" | sed -e 's/-256.*$//g'
# Special case for this since we are matching the end (the start [#] also works):
echo "${TERM%-256*}"
```
`printf %q`Bash:

`printf '%q ' "$@"`POSIX:

```
# TODO: Interpret it back to printf escapes for hard-to-copy chars like \t?
# See also: http://git.savannah.gnu.org/cgit/libtool.git/tree/gl/build-aux/funclib.sh?id=c60e054#n1029
reuse_quote()(
 for i; do
 __i_quote=$(printf '%s\n' "$i" | sed -e "s/'/'\\\\''/g"; echo x)
 printf "'%s'" "${__i_quote%x}"
 done
)
reuse_quote "$@"
```
`jobs` flagsThe only acceptable flags under POSIX sh for `jobs` are
`-l` and `-p` (see
spec). Common flags supported by other shells are `-s`
and `-r`, to check for stopped/suspended jobs and running
jobs. A portable alternative is using `grep` or
`awk`:

```
"$(jobs | awk '/(S|s)(topped|uspended)/')" # instead of jobs -s
"$(jobs | awk '/(R|r)(unning)/')" # instead of jobs -r
```
Although the state of stopped jobs is `Stopped` in Bash
and dash, and it's the one specified by POSIX, `Suspended` is
also a valid alternative (but Zsh happens to not respect the
capitalization, that's why we try to match `suspended`).
Similarly, the state of running jobs is `Running` according
to POSIX. Bash and dash respect this, but Zsh uses
`running`.

Change:

`>& and &>`To:

`command > file 2>&1 or command 2>&1 | othercommand`No Comments / Exceptions

`SIG`Instead of e.g.:

`trap my_handler SIGTERM`use:

```
trap my_handler TERM
# or (`trap -l` for a list of signal numbers; not every one is portable!)
trap my_handler 15
```
Bash:

```
<command>
disown %<command>
```
POSIX:

`nohup <command>`Note that while `nohup` can be used to achieve the same
result, their semantics is different. Also note that `nohup`
will, by default, redirect input and output.

Some errors, currently having their own pages, were linked to this page in older ShellCheck versions.

Depends on what your expected POSIX shell providers would use.

Some features have POSIX proposals:

`local`: https://www.austingroupbugs.net/bug_view_page.php?bug_id=767`$'c-style-escape'`: see aboveShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`#!/bin/sh`
was specified, so ____ is not supported, even when sh is actually
bash.*Note: Removed in v0.4.2
- 2016-01-10*

The shebang indicates that the script works with
`/bin/sh`, but you are using non-standard features that may
not work with `/bin/sh`, **even if /bin/sh is actually
bash**. Bash behaves differently when invoked as `sh`,
and disabling support for the highlighted feature is one part of
that.

Specify `#!/usr/bin/env bash` to ensure that bash (or your
shell of choice) will be used, or rewrite the script to be more
portable.

The Ubuntu wiki has a list of portability issues and suggestions on how to rewrite them.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$(..)` instead
of `'..'` .```
for i in 'seq 1 10'
do
 echo "$i"
done
```
```
for i in $(seq 1 10)
do
 echo "$i"
done
```
The intent was to run the code in the single quotes. This would have
worked with slanted backticks, ``..``, but here the very
similar looking single quotes `'..'` were used, resulting in
a string literal instead of command output.

This is one of the many problems with backticks, so it's better to
use `$(..)` to expand commands.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
for f in foo,bar,baz
do
 echo "$f"
done
```
```
for f in foo bar baz
do
 echo "$f"
done
```
ShellCheck found a `for` loop where the items appeared to
be delimited by commas. These will be treated as literal commas. Use
spaces instead.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`dir/*`, `$var` or
`$(cmd)`?```
for var in value
do
 echo "$var"
done
```
Correct code depends on what you want to do.

To iterate over files in a directory, instead of
`for var in /my/dir` use:

`for var in /my/dir/* ; do echo "$var"; done`To iterate over lines in a file or command output, use a while read loop instead:

`mycommand | while IFS= read -r line; do echo "$line"; done`To iterate over *words* written to a command or function's
stdout, instead of `for var in myfunction`, use

`for var in $(myfunction); do echo "$var"; done`To iterate over *words* in a variable, instead of
`for var in myvariable`, use

`for var in $myvariable; do echo "$var"; done`ShellCheck has detected that your for loop iterates over a single, constant value. This is most likely a bug in your code, caused by you not expanding the value in the way you want.

You should make sure that whatever you loop over will expand into multiple words.

If you stylistically choose to use `for` loops with a
single element, e.g. to align with other code or to set up for future
extensibility (`for target in x86_64; do ..`), you can ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`find -exec` or a
`while read` loop.```
for file in $(find mydir -mtime -7 -name '*.mp3')
do
 (( count++ ))
 echo "Playing file no. $count"
 play "$file"
done
echo "Played $count files"
```
This will fail for filenames containing spaces and similar, such as
`My File.mp3`, and has a series of potential globbing issues
depending on other filenames in the directory like (if you have
`MyFile2.mp3` and `MyFile[2014].mp3`, the former
file will play twice and the latter will not play at all).

There are many possible fixes, each with its pros and cons.

The most general fix (that requires the least amount of thinking to
apply) is having `find` output a `\0` separated
list of files and consuming them in a `while read` loop:

```
while IFS= read -r -d '' file
do
 (( count++ ))
 echo "Playing file no. $count"
 play "$file"
done < <(find mydir -mtime -7 -name '*.mp3' -print0)
echo "Played $count files"
```
In usage it's very similar to the `for` loop: it gets its
output from a `find` statement, it executes a shell script
body, it allows updating/aggregating variables, and the variables are
available when the loop ends.

It requires Bash, and works with GNU, Busybox, OS X, FreeBSD and OpenBSD find, but not POSIX find.

`find`
is just matching globs recursivelyIf you don't need `find` logic like `-mtime -7`
and just use it to match globs recursively (all `*.mp3` files
under a directory), you can instead use `globstar` and
`nullglob` instead of `find`, and still use a
`for` loop:

```
shopt -s globstar nullglob
for file in mydir/**/*.mp3
do
 (( count++ ))
 echo "Playing file no. $count"
 play "$file"
done
echo "Played $count files"
```
This is bash 4 specific.

If you need POSIX compliance, this is a fair approach:

```
find mydir ! -name "$(printf "*\n*")" -name '*.mp3' > tmp
while IFS= read -r file
do
 count=$((count + 1))
 echo "Playing file #$count"
 play "$file"
done < tmp
rm tmp
echo "Played $count files"
```
The only problem is for filenames containing line feeds. A
`! -name "$(printf "*\n*")"` has been added to simply skip
these files, just in case there are any.

If you don't need variables to be available after the loop (here, if
you don't need to print the final play count at the end), you can skip
the `tmp` file and just pipe from `find` to
`while`.

If you don't need a shell script loop body or any form of variable like if we only wanted to play the file, we can dramatically simplify while maintaining POSIX compatibility:

```
# Simple and POSIX
find mydir -name '*.mp3' -exec play {} \;
```
This does not allow things like `count=$((count + 1))`
because -exec does not invoke a shell to interpret that variable
assignment, it simply runs an external command. That external command
can be a shell, however; see below.

If we do need a shell script body but no aggregation, you can do the
above but invoking `sh` (this is still POSIX):

```
find mydir -name '*.mp3' -exec sh -c '
 echo "Playing ${1%.mp3}"
 play "$1"
 ' sh {} \;
```
This would not be possible without `sh`, because
`${1%.mp3}` is a shell construct that `find` can't
evaluate by itself. If we had tried to `count=$((count + 1))`
in this loop, we would have found that the value never changes.

Note that using `+` instead of `\;`, and using
an embedded `for file in "$@"` loop rather than
`"$1"`, will not allow aggregating variables. This is because
for large lists, `find` will invoke the command multiple
times, each time with some chunk of the input.

`for var in $(find ...)` loops rely on word splitting and
will evaluate globs, which will wreck havoc with filenames containing
whitespace or glob characters.

`find -exec` `for i in glob` and
`find`+`while` do not rely on word splitting, so
they avoid this problem.

If you know about and carefully apply `IFS=$'\n'` and
`set -f`, you could choose to ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
for f in $(ls *.wav)
do
 echo "$f"
done
```
```
for f in *.wav
do
 [[ -e "$f" ]] || break # handle the case of no *.wav files
 echo "$f"
done
```
Also note that in Bash, `shopt -s nullglob` will allow the
loop to run 0 times instead of 1 if there are no matches. There are also
several
other conditions to be aware of.

When looping over a set of files, it's always better to use globs when possible. Using command expansion causes word splitting and glob expansion, which will cause problems for certain filenames (typically first seen when trying to process a file with spaces in the name).

The following files can or will break the first loop:

```
touch 'filename with spaces.wav'
touch 'filename with * globs.wav'
touch 'More_Globs[2003].wav'
touch 'files_with_fønny_chæracters_in_certain_locales.wav'
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`ls -l $(getfilename)````
# getfilename outputs 1 file
ls -l "$(getfilename)"
# getfilename outputs multiple files, linefeed separated
getfilename | while IFS='' read -r line
do
 ls -l "$line"
done
```
When command expansions are unquoted, word splitting and globbing will occur. This often manifests itself by breaking when filenames contain spaces.

Trying to fix it by adding quotes or escapes to the data will not work. Instead, quote the command substitution itself.

If the command substitution outputs multiple pieces of data, use a loop instead.

In rare cases you actually want word splitting, such as in

```
# shellcheck disable=SC2046
gcc $(pkg-config --libs openssl) client.c
```
This is because `pkg-config` outputs
`-lssl -lcrypto`, which you want to break up by spaces into
`-lssl` and `-lcrypto`.

A bash alternative in these cases is to use `read -a` for
words or `mapfile` for lines. ksh can also use
`read -a`, or a `while read` loop for lines. In
this case, since `pkg-config` outputs words, you could
use:

```
# Read words from one line into an array in bash and ksh, then expand each as args
read -ra args < <(pkg-config --libs openssl)
gcc "${args[@]}" client.c # expand as args
# Read lines into an array, then expand each as args
readarray -t file_args < <(find /etc/ -type f | grep some-check)
your-linter.sh "${file_args[@]}" # expand as args
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`"$@"` (with quotes) to prevent whitespace problems.Or: Use "${array[@]}" (with quotes) to prevent whitespace problems.

```
cp $* ~/dir
cp ${array[*]} ~/dir
```
```
cp "$@" ~/dir
cp "${array[@]}" ~/dir
```
`$*` and `${array[*]}`, unquoted, is subject to
word splitting and globbing.

Let's say you have three arguments or array elements:
`baz`, `foo bar` and `*`

`"$@"` and `"${array[@]}"`will expand into
exactly that: `baz`, `foo bar` and
`*`

`$*` and `${array[*]}` will expand into
multiple other arguments: `baz`, `foo`,
`bar`, `file.txt` and
`otherfile.jpg`

Since the latter is rarely expected or desired, ShellCheck warns about it.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`=~`
is for regex, but this looks like a glob. Use `=`
instead.`[[ $file =~ *.txt ]]``[[ $file = *.txt ]]`You are using `=~` to match against a regex --
specifically a Extended Regular Expression (ERE) -- but the right-hand
side looks more like a glob:

`*`, like in `*.txt`
`.txt`, like
`readme.txt` but not `foo.sh``txt`, such as `*itxt` but not
`test.txt``*`, like in
`s*`.
`s`, such
as `shell` and `set`.`s`s, such as
`dog` (because it does in fact contain zero or more
`s`'s)Please ensure that the pattern is correct as an ERE, or switch to glob matching if that's what you intended.

This is similar to SC2063, where
`grep "*foo*"` produces an equivalent warning.

If you are aware of the difference, you can ignore this message, but this warning is not emitted
for the more probable EREs `\*.txt`, `\.txt$`,
`^s` or `s+`, so it should rarely be
necessary.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$` on a
variable?```
if [ myvar = "test" ]
then
 echo "Test mode"
fi
```
```
if [ "$myvar" = "test" ]
then
 echo "Test mode"
fi
```
ShellCheck has found a `[ .. ]` or `[[ .. ]]`
comparison that only involves literal strings. The intention was
probably to check a variable or command output instead.

This is usually due to missing `$` or bad quoting:

```
if [[ "myvar" = "test" ]] # always false because myvar is a literal string
if [[ "$myvar" = "test" ]] # correctly compares a variable
if [ 'grep -c foo bar' -ge 10 ] # always false because grep doesn't run
if [ "$(grep -c foo bar)" -ge 10 ] # correctly checks grep output
```
None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
for i in {1..$n}
do
 echo "$i"
done
```
```
for ((i=1; i <= n; i++))
do
 echo "$i"
done
```
In Bash, brace expansion happens before variable expansion. This means that brace expansion will not account for variables.

For integers, use an arithmetic for loop instead. For zero-padded numbers or letters, use of eval may be warranted:

```
from="a" to="m"
for c in $(eval "echo {$from..$to}"); do echo "$c"; done
```
or more carefully (if `from`/`to` could be user
input, or if the brace expansion could have spaces):

```
from="a" to="m"
while IFS= read -d '' -r c
do
 echo "Read $c"
done < <(eval "printf '%s\0' $(printf "{%q..%q}.jpg" "$from" "$to")")
```
None (if you're writing for e.g. zsh, you can use a directive to disable this check)

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`=` in `[[ ]]` to prevent glob matching.`[[ $a = $b ]]``[[ $a = "$b" ]]`When the right-hand side of `=`, `==` or
`!=` is unquoted in `[[ .. ]]`, it will be treated
like a glob.

This has some unexpected consequences like
`[[ $var = $var ]]` being false (for `var='[a]'`),
or `[[ $foo = $bar ]]` giving a different result from
`[[ $bar = $foo ]]`.

The most common intention is to compare one variable to another as strings, in which case the right-hand side must be quoted.

If you explicitly want to match against a pattern, you can ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
flags=("-l", "-d", "--sort=size")
ls "${flags[@]}"
```
```
flags=("-l" "-d" "--sort=size")
ls "${flags[@]}"
```
You appear to have used commas to separate array elements in an array assignment. Other languages require this, but bash instead treats the commas as literal strings.

In the problematic code, the first element is `-l,` with
the trailing comma, and the executed command ends up being
`ls -l, -d, --sort=size`.

In the correct code, the trailing commas have been removed, and the
command will be `ls -l -d --sort=size` as expected.

None (if you actually want a trailing comma in your strings, move them inside the quotes).

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`&&` here, otherwise it's always
true.```
if [[ $1 != foo || $1 != bar ]]
then
 echo "$1 is not foo or bar"
fi
```
```
if [[ $1 != foo && $1 != bar ]]
then
 echo "$1 is not foo or bar"
fi
```
This is not a bash issue, but a simple, common logical mistake applicable to all languages.

`[[ $1 != foo || $1 != bar ]]` is always true (when
`foo != bar`):

`$1 = foo` then `$1 != bar` is true, so the
statement is true.`$1 = bar` then `$1 != foo` is true, so the
statement is true.`$1 = cow` then `$1 != foo` is true, so the
statement is true.`[[ $1 != foo && $1 != bar ]]` matches when
`$1` is neither `foo` nor `bar`:

`$1 = foo`, then `$1 != foo` is false, so
the statement is false.`$1 = bar`, then `$1 != bar` is false, so
the statement is false.`$1 = cow`, then both `$1 != foo` and
`$1 != bar` is true, so the statement is true.This statement is identical to
`! [[ $1 = foo || $1 = bar ]]`, which also works
correctly.

Rare.

```
if [[ $FOO != $BAR || $FOO != $COW ]]
then
 echo "$FOO and $BAR and $COW are not all equal"
fi
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`&&` here```
if (( $1 != 0 || $1 != 3 ))
then
 echo "$1 is not 0 or 3"
fi
```
```
if (( $1 != 0 && $1 != 3 ))
then
 echo "$1 is not 0 or 3"
fi
```
This is not a bash issue, but a simple, common logical mistake applicable to all languages.

`(( $1 != 0 || $1 != 3 ))` is always true:

`$1 = 0` then `$1 != 3` is true, so the
statement is true.`$1 = 3` then `$1 != 0` is true, so the
statement is true.`$1 = 42` then `$1 != 0` is true, so the
statement is true.`(( $1 != 0 && $1 != 3 ))` is true only when
`$1` is not `0` and not `3`:

`$1 = 0`, then `$1 != 3` is false, so the
statement is false.`$1 = 3`, then `$1 != 0` is false, so the
statement is false.`$1 = 42`, then both `$1 != 0` and
`$1 != 3` is true, so the statement is true.This statement is identical to
`! (( $1 == 0 || $1 == 3 ))`, which also works correctly.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ "$var" -leq 42 ]``[ "$var" -le 42 ]`You are using an unknown binary operator in a `test`
expression. Choose one that exists.

In bash, use `help test` to see a list of supported
operators:

```
 FILE1 -nt FILE2 True if file1 is newer than file2 (according to
 modification date).
 FILE1 -ot FILE2 True if file1 is older than file2.
 FILE1 -ef FILE2 True if file1 is a hard link to file2.
 STRING1 = STRING2
 True if the strings are equal.
 STRING1 != STRING2
 True if the strings are not equal.
 STRING1 < STRING2
 True if STRING1 sorts before STRING2 lexicographically.
 STRING1 > STRING2
 True if STRING1 sorts after STRING2 lexicographically.
 arg1 OP arg2 Arithmetic tests. OP is one of -eq, -ne,
 -lt, -le, -gt, or -ge.
Arithmetic binary operators return true if ARG1 is equal, not-equal,
less-than, less-than-or-equal, greater-than, or greater-than-or-equal
than ARG2.
```
None. If you've tested and verified that the operator works but the latest version of ShellCheck says it's unknown, please submit a bug report.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ -E 42 ]``[ -e 42 ]`You are using an unknown unary operator in a `test`
expression. Perhaps it's a typo?

In bash, you can use `help test` to see a list of
supported operators:

```
 -a FILE True if file exists.
 -b FILE True if file is block special.
 -c FILE True if file is character special.
 -d FILE True if file is a directory.
 -e FILE True if file exists.
 -f FILE True if file exists and is a regular file.
 -g FILE True if file is set-group-id.
 -h FILE True if file is a symbolic link.
 -L FILE True if file is a symbolic link.
 -k FILE True if file has its `sticky' bit set.
 -p FILE True if file is a named pipe.
 -r FILE True if file is readable by you.
 -s FILE True if file exists and is not empty.
 -S FILE True if file is a socket.
 -t FD True if FD is opened on a terminal.
 -u FILE True if the file is set-user-id.
 -w FILE True if the file is writable by you.
 -x FILE True if the file is executable by you.
 -O FILE True if the file is effectively owned by you.
 -G FILE True if the file is effectively owned by your group.
 -N FILE True if the file has been modified since it was last read.
```
None. If you've tested and verified that the operator works but the latest version of ShellCheck says it's unknown, please submit a bug report.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`printf` format string. Use
`printf "..%s.." "$foo"`.`printf "Hello, $NAME\n"``printf "Hello, %s\n" "$NAME"``printf` interprets escape sequences and format specifiers
in the format string. If variables are included, any escape sequences or
format specifiers in the data will be interpreted too, when you most
likely wanted to treat it as data. Example:

```
coverage='96%'
printf "Unit test coverage: %s\n" "$coverage"
printf "Unit test coverage: $coverage\n"
```
The first printf writes `Unit test coverage: 96%`.

The second writes
`bash: printf: `\': invalid format character`

Sometimes you may actually want to interpret data as a format string, like in:

```
octToAscii() { printf "\\$1"; }
octToAscii 130
```
In this case, use the `%b` format specifier that expands
escape sequences without interpreting other format specifiers:

```
octToAscii() { printf '%b' "\0$1"; }
octToAscii 130
```
Sometimes you might have a pattern in a variable:

```
filepattern="file-%d.jpg"
printf -v filename "$filepattern" "$number"
```
This has no good rewrite. Please ignore the warning with a directive.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`tr` to prevent glob expansion.`tr -cd [:digit:]``tr -cd '[:digit:]'`From the shell's point of view, unquoted `[:digit:]` is a
glob equivalent to `[digit:]` that matches any single
character filename from the group, such as `d` or
`t`, in the current directory.

If someone starts learning D and creates a directory named
`d` to hold the source code, the glob will be expanded and
the script will end up executing `tr -cd d` instead, which is
clearly unintended.

Quoting the argument prevents this, and will pass it correctly as the
literal string `[:digit:]` no matter which files exist in the
current directory.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-name` so the shell won't interpret
it.`find . -name *.txt``find . -name '*.txt'`Several find options take patterns to match against, including
`-ilname`, `-iname`, `-ipath`,
`-iregex`, `-iwholename`, `-lname`,
`-name`, `-path`, `-regex` and
`-wholename`.

These compete with the shell's pattern expansion, and must therefore
be quoted so that they are passed literally to `find`.

The example command may end up executing as
`find . -name README.txt` after the shell has replaced the
`*.txt` with a matching file `README.txt` from the
current directory.

This may happen today or suddenly in the future.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`grep foo* file``grep "foo*" file`The regex passed to grep frequently contains characters that collide
with globs. The code above is supposed to match "f followed by 1 or more
o's", but if the directory contains a file called "foo.txt", an unquoted
pattern will cause it to become `grep foo.txt file`.

To prevent this, always quote the regex passed to grep, especially when it contains one or more glob character.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`grep '*foo*'``grep 'foo' # or more explicitly, grep '.*foo.*'`In globs, `*` matches any number of any character.

In regex, `*` matches any number of the preceding
character.

`grep` uses regex, not globs, so this means that
`grep '*foo'` is nonsensical because there's no preceding
character for `*`.

If the intention was to match "any number of characters followed by
foo", use `'.*foo'`. Also note that since grep matches
substrings, this will match "fishfood". Use anchors to prevent this,
e.g. `foo$`.

This also means that `f*` will match "hello", because
`f*` matches 0 (or more) "f"s and there are indeed 0 "f"
characters in "hello". Again, use `grep 'f'` to find strings
containing "f", or `grep '^f'` to find strings starting with
"f".

If you're aware of the differences between globs and regex, you can ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
trap "echo \"Finished on $(date)\"" EXIT
trap "rm -fr '$testdir1' '$testdir2'" $TRAP_SIGNALS
trap "rm $tmp" $TRAP_SIGNALS
trap "${remove_aar_temp}" SIGTERM SIGQUIT EXIT
trap "$(shopt -p extglob)" RETURN
```
`trap 'echo "Finished on $(date)"' EXIT`With double quotes, all parameter and command expansions will expand when the trap is defined rather than when it's executed.

In the example with the Problematic code, the message will contain the date on which the trap was declared, and not the date on which the script exits.

Using single quotes will prevent expansion at declaration time, and save it for execution time.

If you don't care that the trap code is expanded early because the commands/variables won't change during execution of the script, or because you want to use the current and not the future values, then you can ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ 1 >2 ] || [ 3>'aaa bb' ] # Simple example of problematic code``[ 1 -gt 2 ] || [ 3 \> 'aaa bb' ] # arithmetical, lexicographical`A word that looks like a redirection in simple shell commands causes it to be interpreted as a redirection. ShellCheck would guess that you don't want it in tests.

When it's among a continuous list of redirections at the end of a
simple `test` command, it's more likely that the user really
meant to do a redirection. Or any other case that you mean to do
that.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`for s in "$(mycommand)"; do echo "$s"; done`The correct code depends on your intention. Let's say you're in a
directory with the files `file.png` and
`My cat.png`, and you want to loop over a command that
outputs (or variable that contains):

```
hello world
My *.png
```
`hello world`,
`My *.png`)`mycommand | while IFS= read -r s; do echo "$s"; done``hello`, `world`,
`My`, `file.png`, `My cat.png`):```
# relies on the fact that IFS by default contains space-tab-linefeed
for s in $(mycommand); do echo "$s"; done
```
`hello world`,
`My cat.png`)```
# explicitly set IFS to contain only a line feed
IFS='
'
for s in $(mycommand); do echo "$s"; done
```
You get this warning because you have a loop that will only ever run exactly one iteration. Since you have a loop, you clearly expect it to run more than once. You just have to decide how it should be split up.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`;` or `+` terminating `-exec`. You
can't use `|`/`||`/`&&`, and
`;` has to be a separate, quoted argument.```
find . -type f -exec shellcheck {} | wc -l \;
find . -exec echo {} ;
```
```
find . -type f -exec sh -c 'shellcheck "$1" | wc -l' -- {} \;
find . -exec echo {} \;
```
`find -exec` is still subject to all normal shell rules,
so all shell features like `|`, `||`,
`&` and `&&` will apply to the
`find` command itself, and not to the command you are trying
to construct with `-exec`.

`find . -exec foo {} && bar {} \;` means run the
command `find . -exec foo {}`, and if find is successful, run
the command `bar "{}" ";"`.

To instead go through each file and run
`foo file && bar file` on it, invoke a shell that can
interpret `&&`:

`find . -exec sh 'foo "$1" && bar "$1"' -- {} \;`You can also use find `-a` instead of shell
`&&`:

`find . -exec foo {} \; -a -exec bar {} \;`This will have the same effect (`-a` is also the default
when two commands are specified, and can therefore be omitted).

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`cp $@ ~/dir``cp "$@" ~/dir`Double quotes around `$@` (and similarly,
`${array[@]}`) prevents globbing and word splitting of
individual elements, while still expanding to multiple separate
arguments.

Let's say you have four arguments: `baz`,
`foo bar`, `*` and `/*/*/*/*`

`"$@"` will expand into exactly that: `baz`,
`foo bar`, `*` and `/*/*/*/*`

`$@` will expand into multiple other arguments:
`baz`, `foo`, `bar`,
`file.txt`, `otherfile.jpg`, and (eventually) a
list of most files on the system

Since the latter is rarely expected or desired, ShellCheck warns about it.

When you want globbing of individual elements.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`2>&1` must be last (or use
`{ cmd > file; } 2>&1` to clarify).`firefox 2>&1 > /dev/null``firefox > /dev/null 2>&1`When it comes to redirection, order matters.

The problematic code means "Point stderr to where stdout is currently pointing (the terminal). Then point stdout to /dev/null".

The correct code means "Point stdout to /dev/null. Then point stderr to where stdout is currently pointing (/dev/null)".

In other words, the problematic code hides stdout and shows stderr. The correct code hides both stderr and stdout, which is usually the intention.

If you actually do want to redirect stdout to a file, and then turn stderr into the new stdout, you can make this more explicit with braces:

`{ firefox > /dev/null; } 2>&1`Also note that this warning does not trigger when output is piped or captured.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-n`
doesn't work with unquoted arguments. Quote or use
`[[ ]]`.```
if [ -n $var ]
then
 echo "var has a value"
else
 echo "var is empty"
fi
```
In POSIX:

```
if [ -n "$var" ]
then
 echo "var has a value"
else
 echo "var is empty"
fi
```
In bash/ksh:

```
if [[ -n $var ]]
then
 echo "var has a value"
else
 echo "var is empty"
fi
```
When `$var` is unquoted, a blank value will cause it to
wordsplit and disappear. If `$var` is empty, these two
statements are identical:

```
[ -n $var ]
[ -n ]
```
`[ string ]` is shorthand for testing if a string is
empty. This is still true if `string` happens to be
`-n`. `[ -n ]` is therefore true, and by extension
so is `[ -n $var ]`.

To fix this, either quote the variable, or (if your shell supports
it) use `[[ -n $var ]]` which generally has fewer caveats
than `[`.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`>` is
for string comparisons. Use `-gt` instead.```
if [[ $var > 10 ]]
then
 echo "Incorrectly triggers when var=5"
fi
```
```
if [[ $var -gt 10 ]]
then
 echo "Correct numerical comparison"
fi
```
`<` and `>`, in both `[[` and
`[` (when escaped) will do a lexicographical comparison, not
a numerical comparison.

This means that `[[ 5 > 10 ]]` is true because 5 comes
after 10 alphabetically. Meanwhile `[[ 5 -gt 10 ]]` is false
because 5 does not come after 10 numerically.

If you want to compare numbers by value, use the numerical comparison
operators `-gt`, `-ge`, `-lt` and
`-le`.

If the strings happen to be version numbers and you're using <, or > to compare them as strings, and you consider this an acceptable thing to do, then you can ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`bc` or
`awk` to compare.`[[ 2 -lt 3.14 ]]````
[[ 200 -lt 314 ]] # Use fixed point math
[[ $(echo "2 < 3.14" | bc) == 1 ]] # Use bc
```
Bash and Posix sh does not support decimals in numbers. Decimals should either be avoided, or compared using a tool that does support them.

If the strings happen to be version numbers and you're using
`<`, or `>` to compare them as strings, and
you consider this an acceptable thing to do, then you can ignore this
warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`\<` to prevent it redirecting (or switch to
`[[ .. ]]`).```
if [ "aardvark" < "zebra" ]
then
 echo "Alphabetical!"
fi
```
```
if [ "aardvark" \< "zebra" ]
then
 echo "Alphabetical!"
fi
```
or optionally in Bash/Ksh:

```
if [[ "aardvark" < "zebra" ]]
then
 echo "Alphabetical!"
fi
```
You are using the operator `<` or `>` in
a `[` test expression.

In this context, it will be considered a file redirection operator
instead, so `[ "aardvark" < "zebra" ]` is equivalent to
`[ "aardvark" ] < ./zebra`, which is true if there exists
a readable file `zebra` in the current directory.

If you wanted to compare two strings lexicographically
(alphabetically), escape the `<` or `>` with
a backslash as in the correct example.

If you want to compare two numbers numerically, use `-lt`
or `-ge` instead.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`=~` in
`[ ]`. Use `[..](..)` instead.`[ "$input" =~ DOC[0-9]*\.txt ] && echo "match"``[[ "$input" =~ DOC[0-9]*\.txt ]] && echo "match"``=~` only works in `[[ .. ]]` tests. It does
not work with `test` or `[` in any shell.

If you're targeting POSIX `sh`, rewrite in terms of
`case` or `grep` instead.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`\<` is required in `[..]`, but invalid in
`[..](..)``[[ aardvark \< zebra ]]``[[ aardvark < zebra ]]`Grammatically speaking, `[` is considered a normal command
name, so `<` and `>` are interpreted as
redirections. When using the lexicographical string operators
`<` and `>` in `[ .. ]`, they
must be escaped (e.g. `\<` or `"<"`).

`[[` is considered its own grammatical construct, and
therefore it does not require (nor does it allow) escaping
`<` or `>`.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`=~`, it'll match literally rather than as a
regex.`[[ $foo =~ "^fo+ bar$" ]]``[[ $foo =~ ^fo+\ bar$ ]]`Quotes on the right hand side of `=~` can be used to match
literally, so that `[[ $1 =~ ^"$2".* ]]` works even if
`$2` contains regex metacharacters. This mirrors the behavior
of globs, `[[ $1 = "$2"* ]]`.

This also means that the problematic code tries to match literal carets and plus signs instead of interpreting them as regular expression matchers. To match as a regex, the regex metacharacters it must be unquoted. Literal parts of the expression can be quoted with double or single quotes, or escaped.

If you do want to match literally just to do a plain substring
search, e.g. `[[ $foo =~ "bar" ]]`, you could ignore this
message, but consider using a more canonical glob match instead:
`[[ $foo == *"bar"* ]]`.

`compat31` `compat31`
See http://stackoverflow.com/questions/218156/bash-regex-with-quotes

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[[ 0=1 ]]``[[ 0 = 1 ]]``[[ 0 = 1 ]]` means "check if 0 and 1 are equal".

`[[ str ]]` is short form for `[[ -n str ]]`,
and means "check if `str` is non-empty". It doesn't matter if
`str` happens to contain `0=1`.

Always use spaces around the comparison operator in `[..]`
and `[..](..)`, otherwise it won't be recognized as an
operator.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$` somewhere?```
if [ "myvar" ]
then
 echo "myvar is set"
fi
```
```
if [ "$myvar" ]
then
 echo "myvar is set"
fi
```
ShellCheck has found a `[ .. ]` or `[[ .. ]]`
statement that just contains a literal string. Such a check does not do
anything useful, and will always be true (or always false, for empty
strings).

This is usually due to missing `$` or bad quoting:

```
if [[ STY ] # always true
if [[ $STY ]] # checks variable $STY
if [[ 'grep foo bar' ]] # always true
if [[ `grep foo bar` ]] # checks grep output (poorly)
if grep -q foo bar # checks for grep match (preferred)
```
None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`(( ))`
doesn't support decimals. Use `bc` or `awk`.Bash arithmetic conditional evaluation can only be performed on integers. More detail: Bash has limited data types which include integer, but everything is effectively untyped.

Suggested workarounds to this constraint use bc or awk, here are some examples.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo $(( 16 - 08 ))``echo $(( 16 - 8 ))`ShellCheck found an integer literal with a leading zero, but containing the digits 8 or 9.

This is invalid, as the integer will be interpreted as an octal value (e.g. 0777 == 0x1FF == 511).

To have the value parsed in base 10, either remove the leading zeros as in the example, or specify the radix explicitly:

`echo $((10#08)) `None

The BASH manual "Shell Arithmetic" chapter

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ .. ]` can't
match globs. Use `[[ .. ]]` or grep.```
if [ $var == *[^0-9]* ]
then
 echo "$var is not numeric"
fi
```
```
if [[ $var == *[^0-9]* ]]
then
 echo "$var is not numeric"
fi
```
`[ .. ]` aka `test` can not match against
globs.

In bash/ksh, you can instead use `[[ .. ]]` which supports
this behavior.

In sh, you can rewrite to use `grep`.

```
if echo $var | grep -q '^[0-9]*$'; then
 echo "$var is numeric"
fi
```
None. If you are not trying to match a glob, quote the argument (e.g.
`[ $var == '*' ]` to match literal asterisk.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`name="foo$n"; echo "${!name}"`.```
var_1="hello world"
n=1
echo "${var_$n}"
```
Bash/ksh:

```
# Use arrays instead of dynamic names
declare -a var
var[1]="hello world"
n=1
echo "${var[n]}"
```
or

```
# Expand variable names dynamically
var_1="hello world"
n=1
name="var_$n"
echo "${!name}"
```
POSIX sh:

```
# Expand dynamically with eval
var_1="hello world"
n=1
eval "tmp=\$var_$n"
echo "${tmp}"
```
You can expand a variable `var_1` with
`${var_1}`, but you can not generate the string
`var_1` with an embedded expansion, like
`${var_$n}`.

Instead, if at all possible, you should use an array. Bash and ksh support both numerical and associative arrays, and an example is shown above.

If you can't use arrays, you can indirectly reference variables by
creating a temporary variable with its name, e.g.
`myvar="var_$n"` and then expanding it indirectly with
`${!myvar}`. This will give the contents of the variable
`var_1`.

If using POSIX sh, where neither arrays nor `${!var}` is
available, `eval` can be used. You must be careful in
sanitizing the data used to construct the variable name to avoid
arbitrary code execution.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`./file`.```
gcc -o myfile file.c
./ myfile
```
```
gcc -o myfile file.c
./myfile
```
Contrary to popular belief, there is no command or syntax
`./` that runs a file.

`./myfile` is simply the shortest path equivalent to
`myfile` that specifies a directory and therefore causes a
shell to run it as-is, instead of trying to find its directory using
$PATH.

Therefore, to run a file in the current directory, use
`./myfile` and not `./ myfile`.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$` or use `_=$((expr))` to avoid executing
output.```
i=4
$(( i++ ))
```
Bash, Ksh:

```
i=4
(( i++ ))
```
POSIX (assuming `++` is supported):

```
i=4
_=$(( i++ ))
```
Alternative POSIX version that does not preserve the exit code:

`: $(( i++ ))``$((..))` expands to a number. If it's the only word on
the line, the shell will try to execute this number as a command
name:

```
$ i=4
$ $(( i++ ))
4: command not found
$ echo $i
5
```
To avoid trying to execute the number as a command name, use one of the methods mentioned:

```
$ i=4
$ _=$(( i++ ))
$ echo $i
5
```
None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
echo $1
for i in $*; do :; done # this one and the next one also apply to expanding arrays.
for i in $@; do :; done
```
```
echo "$1"
for i in "$@"; do :; done # or, 'for i; do'
```
The first code looks like "print the first argument". It's actually "Split the first argument by IFS (spaces, tabs and line feeds). Expand each of them as if it was a glob. Join all the resulting strings and filenames with spaces. Print the result."

The second one looks like "iterate through all arguments". It's actually "join all the arguments by the first character of IFS (space), split them by IFS and expand each of them as globs, and iterate on the resulting list". The third one skips the joining part.

Quoting variables prevents word splitting and glob expansion, and prevents the script from breaking when input contains spaces, line feeds, glob characters and such.

Strictly speaking, only expansions themselves need to be quoted, but for stylistic reasons, entire arguments with multiple variable and literal parts are often quoted as one:

```
$HOME/$dir/dist/bin/$file # Unquoted (bad)
"$HOME"/"$dir"/dist/bin/"$file" # Minimal quoting (good)
"$HOME/$dir/dist/bin/$file" # Canonical quoting (good)
```
When quoting composite arguments, make sure to exclude globs and
brace expansions, which lose their special meaning in double quotes:
`"$HOME/$dir/src/*.c"` will not expand, but
`"$HOME/$dir/src"/*.c` will.

Note that `$( )` starts a new context, and variables in it
have to be quoted independently:

```
echo "This $variable is quoted $(but this $variable is not)"
echo "This $variable is quoted $(and now this "$variable" is too)"
```
Sometimes you want to split on spaces, like when building a command line:

```
options="-j 5 -B"
[[ $debug == "yes" ]] && options="$options -d"
make $options file
```
Just quoting this doesn't work. Instead, you should have used an array (bash, ksh, zsh):

```
options=(-j 5 -B) # ksh88: set -A options -- -j 5 -B
[[ $debug == "yes" ]] && options=("${options[@]}" -d)
make "${options[@]}" file
```
or a function (POSIX):

```
make_with_flags() {
 [ "$debug" = "yes" ] && set -- -d "$@"
 make -j 5 -B "$@"
}
make_with_flags file
```
To split on spaces but not perform glob expansion, POSIX has a
`set -f` to disable globbing. You can disable word splitting
by setting `IFS=''`.

Similarly, you might want an optional argument:

```
debug=""
[[ $1 == "--trace-commands" ]] && debug="-x"
bash $debug script
```
Quoting this doesn't work, since in the default case,
`"$debug"` would expand to one empty argument while
`$debug` would expand into zero arguments. In this case, you
can use an array with zero or one elements as outlined above, or you can
use an unquoted expansion with an alternate value:

```
debug=""
[[ $1 == "--trace-commands" ]] && debug="yes"
bash ${debug:+"-x"} script
```
This is better than an unquoted value because the alternative value
can be properly quoted, e.g.
`wget ${output:+ -o "$output"}`.

Here are two common cases where this warning seems unnecessary but may still be beneficial:

```
cmd <<< $var # Requires quoting on Bash 3 (but not 4+)
: ${var=default} # Should be quoted to avoid DoS when var='*/*/*/*/*/*'
```
As always, this warning can be ignored on a case-by-case basis.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`EOF` to make here document expansions happen on the server
side rather than on the client.```
ssh host.example.com << EOF
 echo "Logged in on $HOSTNAME"
EOF
```
```
ssh host.example.com << "EOF"
 echo "Logged in on $HOSTNAME"
EOF
```
When the end token of a here document is unquoted, parameter expansion and command substitution will happen on in contents of the here doc.

This means that before sending the commands to the server, the client
replaces `$HOSTNAME` with localhost, thereby sending
`echo "Logged in on localhost"` to the server. This has the
effect of printing the client's hostname instead of the server's.

Scripts with any kind of variable use are especially problematic because all references will be expanded before the script run. For example,

```
ssh host << EOF
 x="$(uname -a)"
 echo "$x"
EOF
```
will never print anything, neither client nor server details, since before evaluation, it will be expanded to:

```
 x="Linux localhost ... x86_64 GNU/Linux"
 echo ""
```
By quoting the here token, local expansion will not take place, so
the server sees `echo "Logged in on $HOSTNAME"` which is
expanded and printed with the server's hostname, which is usually the
intention.

If the client should expand some or all variables, this message can and should be ignored.

To expand a mix of local and remote variables, the here doc end token should be unquoted, and the remote variables should be escaped, e.g.

```
ssh host.example.com << EOF
 echo "Logged in on \$HOSTNAME from $HOSTNAME"
EOF
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$HOME`.`rm "~/Desktop/$filename"``rm "$HOME/Desktop/$filename"`Tilde does not expand to the user's home directory when it's single
or double quoted. Use double quotes and `$HOME` instead.

Alternatively, the `~/` can be left unquoted, as in
`rm ~/"Desktop/$filename"`.

If you don't want the tilde to be expanded, you can ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
args='-lh "My File.txt"'
ls $args
```
In Bash/Ksh with arrays:

```
args=(-lh "My File.txt")
ls "${args[@]}"
```
or in POSIX overwriting `"$@"`:

```
set -- -lh "My File.txt"
ls "$@"
```
or in POSIX via functions:

```
myls() { ls "-lh" "My File.txt"; }
myls
```
Bash does not interpret data as code. Consider almost any other languages, such as Python:

```
#!/usr/bin/env python3
print(1 + 1) # prints 2
a = "1 + 1"
print(a) # prints 1 + 1, not 2
```
Here, `1 + 1` is Python syntax for adding numbers.
However, passing a literal string containing this expression does not
cause Python to interpret it, see the `+`, and produce the
calculated result.

Similarly, `"My File.txt"` is Bash syntax for a single
word with a space in it. However, passing a literal string containing
this expression does not cause Bash to interpret it, see the quotes, and
produce the tokenized result.

The solution is to use an array instead, whenever possible.

If due to `sh` compatibility you can't use arrays, you can
sometimes use functions instead. Instead of trying to create a set of
arguments that has to be passed to a command, create a function that
calls the function with arguments plus some more:

```
ffmpeg_with_args() {
 ffmpeg -filter_complex '[#0x2ef] setpts=PTS+1/TB [sub] ; [#0x2d0] [sub] overlay' "$@"
}
ffmpeg_with_args -i "My File.avi" "Output.avi"
```
In other cases, you may have to use `eval` instead, though
this is often fragile and insecure. If you get it wrong, it'll appear to
work great in all test cases, and may still lead to various forms of
security vulnerabilities and breakage:

```
quote() { local q=${1//\'/\'\\\'\'}; echo "'$q'"; }
args="-lh $(quote "My File.txt")"
eval ls "$args" # Do not use unless you understand implications
```
If you ever accidentally forget to use proper quotes, such as with:

```
for f in *.txt; do
 args="-lh '$1'" # Example security exploit
 eval ls "$args" # Do not copy and use
done
```
Then you can use `touch "'; rm -rf \$'\x2F'; '.txt"` (or
someone can trick you into downloading a file with this name, or create
a zip file or git repo containing it, or changing their nick and have
your chat client create the file for a chat log, or...), and running the
script to list your files will run the command
`rm -rf /`.

Few and far between, such as, prompt variables. This from
`man bash` "PROMPTING":

`After the string is decoded, it is expanded via parameter expansion, command substitution, arithmetic expansion, and quote removal, subject to the value of the promptvars shell option (see the description of the shopt command under SHELL BUILTIN COMMANDS below). This can have unwanted side effects if escaped portions of the string appear within command substitution or contain characters special to word expansion.`

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

See companion warning, SC2089.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$()` to avoid executing output (or use
`eval` if intentional).```
if $(which epstopdf)
then
 echo "Found epstopdf"
fi
```
or

```
make_command() {
 printf 'cat header %q footer > %q\n' "$1" "$2" | tee log
}
$(make_command foo.html output/foo.html)
```
```
if which epstopdf
then
 echo "Found epstopdf"
fi
```
or

```
make_command() {
 printf 'cat header %q footer > %q\n' "$1" "$2" | tee log
}
eval "$(make_command foo.html output/foo.html)"
```
ShellCheck has detected that you have a command that just consists of
a command substitution. This often happens when you want to run a
command (possibly from a variable name), without realizing that
`$(..)` is for capturing and not for executing.

For example, if you have this shell function:

`sayhello() { echo "hello world"; }`Then `$(sayhello)` will:

`sayhello`, capturing "hello world"`hello world`, resulting in
`bash: hello: command not found`Meanwhile, just `sayhello` will:

`sayhello`, outputting "hello world" to screenNote that this is equally true if the command is in a variable, e.g.
`x=sayhello; $($x)`.

If you *do* have a command that outputs a second command,
similar to how `ssh-agent` outputs `export`
commands to run, then you should do this via `eval`. This
way, quotes, pipes, redirections, semicolons, and other shell constructs
will work as expected. Note that this kind of design is best avoided
when possible, since correctly escaping all values can be difficult and
error prone.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

Backticks does the same thing as $(..). See SC2091 for a description of the same problem with this syntax.

$(..)

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`exec` if script should continue after this command.```
echo "Starting compilation"
exec ./compile
echo "Starting deployment"
exec ./deploy
```
```
echo "Starting compilation"
./compile
echo "Starting deployment"
./deploy
```
The script contains an `exec` command followed by other
commands in the same block. This is likely an error because the script
will not resume after an `exec` command.

Instead, "exec" refers to the Unix process model's idea of execution
(see `execve(2)`),
in which the current process stops its current program and replaces it
with a new one. This is mainly used in wrapper scripts.

To execute another script or program and then continue, simply drop
the `exec` as in the example.

If the code after the `exec` is only there to handle a
failure in executing the command you can ignore this warning. For this
reason, ShellCheck suppresses the warning if `exec` is only
followed by `echo`/`exit` commands.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`grep foo file.txt | sed -e 's/foo/bar/g' > file.txt``grep foo file.txt | sed -e 's/foo/bar/g' > tmpfile && mv tmpfile file.txt`Each step in a pipeline runs in parallel.

In this case, `grep foo file.txt` will immediately try to
read `file.txt` while `sed .. > file.txt` will
immediately try to truncate it.

This is a race condition, and results in the file being partially or (far more likely) entirely truncated.

Note that this can also be a problem when you write to a file and
read from it later in the pipe. The second command (which reads the
file) may not see all the output of the first. An exception in this case
is a non-greedy file reader like `less`, for example
`python foo.py 2> errfile.txt | less - errfile.txt` will
successfully allow you to see stdout and stderr separately in less.

You can ignore this error if:

`echo log.txt > log.txt`.`cat file | sed s/foo/bar/ > file`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`ssh -n` to prevent ssh from swallowing stdin.The same error applies to multiple commands, like
`ffmpeg -nostdin` and
`mplayer -noconsolecontrols`.

```
while read -r host
do
 ssh "$host" "uptime"
done < hosts.txt
```
```
while read -r host
do
 ssh -n "$host" "uptime"
done < hosts.txt
```
or

```
while read -r host
do
 ssh "$host" <<'EOF'
uptime
EOF
done < hosts.txt
```
or

By using a pipe and avoiding the use of the stdin file descriptor, this ensures that commands in the loop are not interfered with.

```
exec 3< hosts.txt
while read -r host
do
 ssh "$host" "uptime"
done <&3
# Close the file descriptor
exec 3<&-
```
Commands that process stdin will compete with the `read`
statement for input. This is especially tricky for commands you wouldn't
expect reads from stdin, like `ssh .. uptime`,
`ffmpeg` and `mplayer`.

The most common symptom of this is a `while read` loop
only running once, even though the input contains many lines. This is
because the rest of the lines are swallowed by the offending
command.

To refuse such commands input, you can use a command specific option
like `ssh -n` or `ffmpeg -nostdin`.

More generally, you can also redirect their stdin with
`< /dev/null`. This works for all commands with this
behavior.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`#!/usr/bin/env bash -x````
#!/usr/bin/env bash
set -x
```
Most operating systems, including POSIX, Linux and FreeBSD, allow
only a single parameter in the shebang. The example is equivalent to
calling `env 'bash -x'` instead of
`env 'bash' '-x'`, and it will therefore fail.

The shebang should be rewritten to use at most one parameter. Shell options can instead be set in the body of the script.

macOS X currently allows multiple words in the shebang. Scripts running on OSX exclusively can ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`name=World cmd -m "Hello $name"````
name=World
cmd -m "Hello $name"
```
If the original goal was to limit the scope of the variable, this can also be done in a subshell:

```
(
 name=World
 cmd -m "Hello $name"
) # 'name' does not leave this subshell
```
In `name=World cmd "$name"`, `name=World` is
passed in as part of the environment to `cmd` (i.e., in the
`envp` parameter to execve(2)). This means that
`cmd` and its children will see the parameter, but no other
processes will.

However, `"$name"` is not expanded by `cmd`.
`"$name"` is expanded by the shell before `cmd` is
ever executed, and thus it will not use the new value.

The solution is to set the variable first, then use it as a
parameter. If limited scope is desired, a `( subshell )` can
be used.

In the strange and fabricated scenarios where the script and a program uses a variable name for two different purposes, you can ignore this message. This is hard to conceive, since scripts should use lowercase variable names specifically to avoid collisions with the environment.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

See companion warning SC2097.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$((..))` for
arithmetics, e.g. `i=$((i + 2))````
i=3
i=i + 2
```
```
i=3
i=$((i + 2))
```
Unlike most languages, variable assignments in shell scripts are space sensitive and (almost) always assign strings.

To evaluate a mathematical expressions, use `$((..))` as
in the correct code:

`i=$((i + 2)) # Spaces are fine inside $((...))`In the problematic code, `i=i + 2` will give an error
`+: command not found` because the expression is interpreted
similar to something like `LC_ALL=C wc -c` instead of
numerical addition:

```
 Prefix assignment Command Argument
 LC_ALL=C wc -c
 i=i + 2
```
If you wanted to assign a literal string, quote it:

`game_score="0 - 2"`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$((..))` for
arithmetics, e.g. `i=$((i + 2))````
i=3
i=i+2
```
```
i=3
i=$((i + 2))
```
Unlike most languages, variable assignments (almost) always assigns
strings and not expressions. In the example code, `i` will
become the string `i+2` instead of the intended
`5`.

To instead evaluate a mathematical expressions, use
`$((..))` as in the correct code.

If you wanted to assign a literal string, quote it:

`description="friendly-looking"`ShellCheck (as of v0.5) doesn't recognize Bash/Ksh numeric variables
created with `declare -i` where this syntax is valid. Using
`$((..))` still works, but you can also ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[]`, e.g. `[:digit:](:digit:)`.`gzip file[:digit:]*.txt``gzip file[:digit:](:digit:)*.txt`Predefined character groups are supposed to be used inside character
ranges. `[:digit:]` matches one of "digt:" just like
`[abc]` matches one of "abc". `[:digit:](:digit:)`
matches a digit.

When passing an argument to `tr` which parses these by
itself without relying on globbing, you should quote it instead, e.g.
`tr -d '[:digit:]'`

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo [100-999].txt``echo [1-9][0-9][0-9].txt`ShellCheck found a glob range expression (such as `[a-z]`)
that contains multiple of the same character.

Range expressions can only be used to match a single character in a
given set, so `[ab]` and `[abba]` will both match
the same thing: either one `a` or one `b`.

Having multiple of the same character often means you're trying to
match more than one character, such as in the problematic example where
someone tried to match any number from 100 to 999. Instead, it matches a
single digit just like `[0-9].txt`, and specifies 0, 1 and 9
multiple times.

In Bash, most uses can be rewritten using extglob and/or brace expansion. For example:

```
cat *.[dev,prod,test].conf # Doesn't work
cat *.{dev,prod,test}.conf # Works in bash
cat *.@(dev|prod|test).conf # Works in bash with `shopt -s extglob`
```
In POSIX sh, you may have to write multiple globs, one after the other:

`cat *.dev.conf *.prod.conf *.test.conf`There is currently a bug in which a range expression whose contents
is a variable gets parsed verbatim, e.g. `[$foo]`. In this
case, either ignore the warning or make the square brackets part of the
variable contents instead.

v0.7.2 and below would unintentionally show this warning for
subscripts in arrays in `[[ -v array[xx] ]]` and other
dereferencing operators. In these versions, you can either ignore the message or quote the word (as in
`[[ -v 'array[xx]' ]]`)

Note that IPv6 URLs trigger this warning, but the correct solution in this case is to quote them:

`curl 'http://[2607:f8b0:4002:c0c::65]/'`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`( subshell )` to avoid having to `cd` back.```
for dir in */
do
 cd "$dir"
 convert index.png index.jpg
 cd ..
done
```
```
for dir in */
do
 (
 cd "$dir" || exit
 convert index.png index.jpg
 )
done
```
or

```
for dir in */
do
 cd "$dir" || exit
 convert index.png index.jpg
 cd ..
done
```
When doing `cd dir; somestuff; cd ..`, `cd dir`
can fail when permissions are lacking, if the dir was deleted, or if
`dir` is actually a file.

In this case, `somestuff` will run in the wrong directory
and `cd ..` will take you to an even more wrong directory. In
a loop, this will likely cause the next `cd` to fail as well,
propagating this error and running these commands far away from the
intended directories.

Check `cd`s exit status and/or use subshells to limit the
effects of `cd`.

If you set variables you can't use a subshell. In that case, you
should definitely check the exit status of `cd`, which will
also silence this suggestion.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`return` instead of `break`.```
foo() {
 if [[ -z $1 ]]
 then
 break
 fi
 echo "Hello $1"
}
```
```
foo() {
 if [[ -z $1 ]]
 then
 return 1
 fi
 echo "Hello $1"
}
```
`break` or `continue` are used to abort or
continue a loop, and are not the right way to exit a function. Use
`return` instead.

The `break` or `continue` may be intended for a
loop that calls the function:

```
# Rarely valid
foo() { break; echo $?; }
while true; do foo; done
```
This is undefined behavior in POSIX sh. Different shells do different things.

When the function is called from a loop:

`ksh` keeps going and `$?` is 0.`bash` version 4.4+ prints an error "break: only
meaningful in a `for', `while', or `until' loop", the function keeps
going, and `$?` is 0.`bash` versions before 4.4, will return from the
function, break the loop calling the function, or exit a subshell if
there's one in between.`dash`, BusyBox `ash`: like above.When the function is not called from a loop:

`bash` versions print an error "break: only
meaningful in a `for', `while', or `until' loop", the function keeps
going, and `$?` is 0.`ksh`, `dash` and `ash` silently
keep going and `$?` is 0.Due to the many different implementations, many of which are not
helpful, it's recommended to use proper flow control. A typical solution
is making sure the function `return`s success/failure, and
calling `myfunction || break` in the loop.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`break` is only valid
in loops```
case "$1" in
 -v)
 verbose=1
 break
 ;;
 -d)
 debug=1
esac
```
```
case "$1" in
 -v)
 verbose=1
 ;;
 -d)
 debug=1
esac
```
`break` or `continue` was found outside a loop.
These statements are valid only in loops. In particular,
`break` is not required in `case` statements as
there is no implicit fall-through.

To return from a function or sourced script, use `return`.
To exit a script, use `exit`.

It's possible to `break`/`continue` in a
function without a loop. The call will then affect the loop – if any –
that the function is invoked from, but this is obviously not good coding
practice.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
for i in a b c; do
 echo hi | grep -q bye | break
done
```
```
for i in a b c; do
 echo hi | grep -q bye || break
done
```
The most common cause of this issue is probably using a single
`|` when `||` was intended. The reason this
message appears, though, is that a construction like this, intended to
surface a failure inside of a loop:

`for i in a b c; do false | break; done; echo ${PIPESTATUS[@]}`may appear to work:

```
$ for i in a b c; do false | break; done; echo ${PIPESTATUS[@]}
1 0
```
What's actually happening, though, becomes clear if we add some
`echo`s; the entire loop completes, and the
`break` has no effect.

```
$ for i in a b c; do echo $i; false | break; done; echo ${PIPESTATUS[@]}
a
b
c
1 0
$ for i in a b c; do false | break; echo $i; done; echo ${PIPESTATUS[@]}
a
b
c
0
```
Because bash processes pipelines by creating subshells, control
statements like `break` only take effect in the subshell.

Contrast with the related, but different, problem in this link.

Bash Reference Manual: Pipelines, esp.:

Each command in a pipeline is executed in its own subshell.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ a && b ]`, use
`[ a ] && [ b ]`.`[ "$1" = "-v" && -z "$2" ]``[ "$1" = "-v" ] && [ -z "$2" ]``&&` can not be used in a `[ .. ]` test
expression. Instead, make two `[ .. ]` expressions and put
the `&&` between them.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[..](..)`, use
`&&` instead of `-a`.`[[ "$1" = "-v" -a -z "$2" ]]``[[ "$1" = "-v" && -z "$2" ]]``-a` for logical AND is not supported in a
`[[ .. ]]` expression. Use `&&`
instead.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ a || b ]`, use `[ a ] || [ b ]`.`[ "$1" = "-v" || "$1" = "-help" ]``[ "$1" = "-v" ] || [ "$1" = "-help" ]``||` cannot be used in a `[ .. ]` test
expression. Instead, make two `[ .. ]` expressions and put
the `||` between them.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[..](..)`, use
`||` instead of `-o`.`[[ "$1" = "-v" -o "$1" = "-help" ]]``[[ "$1" = "-v" || "$1" = "-help" ]]``-o` for logical OR is not supported in a
`[[ .. ]]` expression. Use `||` instead.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`function` keyword and `()` at the
same time.```
#!/bin/ksh
function foo() {
 echo "Hello World"
}
```
```
# POSIX sh function
foo() {
 echo "Hello World"
}
```
or

```
# ksh extended function
function foo {
 echo "Hello World"
}
```
Ksh allows two ways of defining functions: POSIX sh style
`foo() { ..; }` and Ksh specific
`function foo { ..; }`.

ShellCheck found a function definition that uses both at the same
time, `function foo() { ..; }` which is not allowed. Use one
or the other.

Note that the two are not identical, for example:

`typeset` in a `function foo` will
create a local variable, while in `foo()` it will create a
global variable.`function foo` has its own trap context, while
`foo()` shares them with the current process.`function foo` will set `$0` to foo, while
`foo()` will inherit `$0` from the current
process.In Bash, `function foo() { ..; }` is allowed, and
`function foo` and `foo()` are identical. This
warning does not trigger when the shebang is e.g.
`#!/bin/bash`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`function`
keyword is non-standard. Delete it.```
#!/bin/sh
function hello() {
 echo "Hello World"
}
```
```
#!/bin/sh
hello() {
 echo "Hello World"
}
```
The `function` keyword is a feature of Bash and Ksh, so
code that uses it may be intended to be a Bash or Ksh script
instead:

```
#!/bin/bash
function hello() {
 echo "Hello World"
}
```
`function` is a non-standard keyword that can be used to
declare functions in Bash and Ksh.

In POSIX `sh` and `dash`, a function is instead
declared without the `function` keyword as in the correct
example.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`function`
keyword is non-standard. Use `foo()` instead of
`function foo`.```
#!/bin/sh
function hello {
 echo "Hello World"
}
```
```
#!/bin/sh
hello() {
 echo "Hello World"
}
```
`function` is a non-standard keyword that can be used to
declare functions in Bash and Ksh.

In POSIX `sh` and `dash`, a function is defined
without a `function` keyword. Instead, the function name is
followed by `()` as in the correct example.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`rm -rf /usr /lib/nvidia-current/xorg/xorg``rm -rf /usr/lib/nvidia-current/xorg/xorg`The example line of code was an actual bug in the Bumblebee NVIDIA driver.

Due to an accidental space, it deleted `/usr` instead of
just the particular directory.

If you do intend to delete a system directory, such as when working in a chroot or initramfs, you can disable this message with a directive:

```
# shellcheck disable=SC2114
rm -rf /usr
```
Previous versions of shellcheck, up to and including 0.4.6, would
ignore `rm` statements containing a `--` (an
arbitrary convention). This is no longer the case.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`"${var:?}"` to ensure this never expands to `/*`
.`rm -rf "$STEAMROOT/"*``rm -rf "${STEAMROOT:?}/"*`If `STEAMROOT` is empty, this will end
up deleting everything in the system's root directory.

Using `:?` will cause the command to fail if the variable
is null or unset. Similarly, you can use `:-` to set a
default value if applicable.

In the case of command substitution, assign to a variable first and
then use `:?`. This is relevant even if the command seems
simple and obviously correct, since forks and execs can fail due to
external system limits and conditions, resulting in a blank
substitution.

For more details about `:?` see the "Parameter Expansion"
section of the Bash man page.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`cmd $(echo foo)`, just use
`cmd foo`.```
greeting=$(echo "Hello, $name")
# or
tar czf "$(echo "$(date +%F).tar.gz")" *
```
```
greeting="Hello, $name"
# or
tar czf "$(date +%F).tar.gz" *
```
You appear to be using `echo` to write a value to stdout,
and then using `$(..)` or ``..`` to capture the
value again. This is as pointless as mailing yourself a postcard: you
already have what you want, so there's no need to send it on a round
trip.

You can just replace `$(echo myvalue)` with
`myvalue`.

Sometimes this pattern is used because of side effect of
`echo` or expansions. For example, here
`$(echo ..)` is used to expand a glob.

```
glob="*.png"
files="$(echo $var)"
```
The `echo` is not useless, but this code is problematic
because it concatenates filenames by spaces. This will break filenames
containing spaces and other characters later when the list is split
again. Better options are:

`files=( $glob ); echo "The first file is ${files[0]}"``set -- $glob; echo "The first file is $1"``for file in $glob; do ...`All three methods will let you avoid issues with special characters in filenames.

As another example, here `$(echo ..)` is used to expand
escape sequences:

```
unexpanded='var\tvalue'
expanded="$(echo "$var")"
```
In this case, use `printf` instead. It's well defined with
regard to escape sequences.

Finally, if you really do want to concatenate a series of elements by
a character like space, consider doing it explicitly with
`for` or `printf` (e.g.
`printf '%s\n' $glob`).

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`su -c` or
`sudo`.```
whoami
su
whoami
```
```
whoami
sudo whoami
```
It's commonly believed that `su` makes a session run as
another user. In reality, it starts an entirely new shell, independent
of the one currently running your script.

`su; whoami` will start a root shell and wait for it to
exit before running `whoami`. It will not start a root shell
and then proceed to run `whoami` in it.

To run commands as another user, use `sudo some command`
or `su -c 'some command'`. `sudo` is preferred
when available, as it doesn't require additional quoting and can be
configured to run passwordless if desired.

If you're aware of the above and want to e.g. start an interactive shell for a user, feel free to ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`|&`. Use `2>&1 |````
#!/usr/bin/ksh
make |& tee ~/log
```
```
#!/usr/bin/ksh
make 2>&1 | tee ~/log
```
You are using the Bash specific shorthand `|&`, but
your script is running with Ksh. Rewrite it to its full,
POSIX-compatible form as shown in the example.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

foo "$@"

$1

See companion warning SC2120.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
sayhello() {
 echo "Hello $1"
}
sayhello
```
`./myscript World` just prints "Hello " instead of "Hello
World".

```
sayhello() {
 echo "Hello $1"
}
sayhello "$@"
```
`./myscript World` now prints "Hello World".

In a function, `$1` and up refers to the function's
parameters, not the script's parameters.

If you want to process your script's parameters in a function, you
have to explicitly pass them. You can do this with
`myfunction "$@"`.

Note that `"$@"` refers to the current context's
positional parameters, so if you call a function from a function, you
have to pass in `"$@"` to both of them:

```
first() { second "$@"; }
second() { echo "The first script parameter is: $1"; }
first "$@"
```
If the parameters are optional and you currently just don't want to use them, you can ignore this message. In versions strictly greater than v0.6.0, ignoring SC2120 on a function will also disable SC2119 on each of the call sites.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`var=value`, not `set ..`.```
set var=42
set var 42
```
`var=42``set` is not used to set or assign variables in Bourne
shells. It's used to set shell options and positional parameters.

To assign variables, use `var=value` with no
`set` or other qualifiers.

If you actually do want to set positional parameters, simply quoting
them or using `--` will make shellcheck stop warning, e.g.
`set -- var1 var2` or `set "foo=bar"`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`>=`
is not a valid operator. Use `! a < b` instead.`[[ a >= b ]]``[[ ! a < b ]]`The operators `<=` and `>=` are not
supported by Bourne shells. Instead of "less than or equal", rewrite as
"not greater than".

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`PATH` is
the shell search path. Use another name.```
PATH=/my/dir
cat "$PATH/myfile"
```
Good practice: always use lowercase for unexported variables.

```
path=/my/dir
cat "$path/myfile"
```
Bad practice: use another uppercase name.

```
MYPATH=/my/dir
cat "$MYPATH/myfile"
```
`PATH` is where the shell looks for the commands it
executes. By inadvertently overwriting it, the shell will be unable to
find commands (like `cat` in this case).

You get this warning when ShellCheck suspects that you didn't meant to overwrite it (because it's a single path with no path separators).

Best shell scripting practice is to always use lowercase variable names to avoid accidentally overwriting exported and internal variables.

If you're aware of the above and really do want to set your shell
search path to `/my/dir`, you can ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`*` instead of
`@` to concatenate.```
# Want to store multiple elements in var
var=$@
for i in $var; do ..; done
```
or

```
set -- Hello World
# Want to concatenate multiple elements into a single string
msg=$@
echo "You said $msg"
```
```
# Bash: use an array variable
var=( "$@" )
for i in "${var[@]}"; do ..; done
# POSIX sh: without array support, one possible workaround
# is to store elements concatenated with a delimiter (here linefeed/newline)
var=$(printf '%s\n' "$@")
printf '%s\n' "$var" | while IFS='' read -r line; do ..; done
```
or

```
#!/bin/sh
set -- Hello World
# Explicitly concatenates all the array elements into a single string
msg=$*
echo "You said $msg"
```
Arrays and `$@` can contain multiple elements. Simple
variables contain only one. When assigning multiple elements to one
element, the default behavior depends on the shell (bash concatenates
with spaces, zsh concatenates with first char of `IFS`).

Since doing this usually indicates a bug, ShellCheck warns and asks you to be explicit about what you want.

If you want to assign N elements as N elements in Bash or Ksh, use an
array, e.g. `myArray=( "$@" )`.

Dash and POSIX sh do not support arrays. In this case, either
concatenate the values with some delimiter that you can split on later
(the example uses linefeeds and splits them back up with a
`while read` loop), or keep the values as positional
parameters without putting them in an intermediate variable.

If you want to assign N elements as 1 element by concatenating them,
use `*` instead of `@`, e.g.
`myVar=${myArray[*]}` (this separates elements with the first
character of `IFS`, usually space).

The same is true for `${@: -1}`, which results in 0 or 1
elements: `var=${*: -1}` assigns the last element or an empty
string.

None.

`filelist="${filelist[@]}" "$filename"`What was meant is:

`filelist=("${filelist[@]}" "$filename")`Note: This syntax is also compatible with older shells, as opposed to
`filelist+=("$filename")` which works in later shells (bash
3.1+ and zsh 4.2+), only.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
foo={1..9}
echo "$foo"
```
```
foo=/some/path/*
echo "$foo"
```
```
foo=( {1..9} )
echo "${foo[@]}"
```
```
foo=(/some/path/*)
echo "${foo[@]}"
```
Note that either of these will trigger SC3030
("In POSIX sh, array references are undefined") if you are using
`sh` and not e.g. `bash`.

`echo *.png {1..9}` expands to all png files and numbers
from 1 to 9, but `var=*.png` or `var={1..9}` will
just assign the literal strings `'*.png'` and
`'{1..9}'`.

To make the variable contain all png files or 1 through 9, use an array as demonstrated.

If you intended to assign these values as literals, quote them (e.g.
`var="*.png"`).

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`grep -c` instead of `grep | wc``grep foo | wc -l``grep -c foo`Instead of:

`grep foo *.log | wc -l`You can pipe all the file contents into `grep` (passing
the files directly to `grep` causes `-c` to print
each file's count separately, rather than the total):

`cat *.log | grep foo -c`This is purely a stylistic issue. `grep` can count lines
without piping to `wc`.

Often this number is only used to see whether there are matches (i.e.
`== 0`). In these cases it's clearer and more efficient to
use `grep -q` and check its exit status:

```
if grep -q pattern file; then
 echo "The file contains the pattern"
fi
```
Also note that in `foo | grep bar | wc -l`,
`wc` will mask the exit code of `grep` by default
(i.e. without `set -o pipefail`), and always return success.
If replacing with `foo | grep -c bar`, `grep` will
exit non-zero when there are no matches. This is generally desirable
(see above), but may require handling when used with
`set -e`.

If you find piping to `wc` is clearer in a given situation
it's fine to ignore this error.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`${ ..; }`,
specify `#!/usr/bin/env ksh`.(Or "To use cases with `;;&`, specify
`#!/usr/bin/env bash`)

```
#!/usr/bin/env bash
var=${ mycmd; };
```
or

```
#!/usr/bin/env ksh
case "$1" in
 foo) echo "Foo!" ;;&
 f*) echo "F-something at least" ;;
esac
```
```
#!/usr/bin/env ksh
var=${ mycmd; };
```
or

```
#!/usr/bin/env bash
case "$1" in
 foo) echo "Foo!" ;;&
 f*) echo "F-something at least" ;;
esac
```
You are using a shell syntax feature not supported by the script's shell. Either rewrite the construct, or switch to a different shell interpreter.

ShellCheck 0.10.0 and below warns about `${ ..; }` command
expansions when using Bash. However, Bash 5.3 added support for this
construct. If you are using this construct in Bash 5.3, either ignore
the warning or upgrade ShellCheck.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
myarray=(foo bar)
for f in $myarray
do
 cat "$f"
done
```
```
myarray=(foo bar)
for f in "${myarray[@]}"
do
 cat "$f"
done
```
When referencing arrays, `$myarray` is equivalent to
`${myarray[0]}` -- which is usually the first of multiple
elements. This is also true for associative arrays. Therefore, if 0
(zero) is not a valid key, `$myarray` expands to an empty
string.

To get all elements as separate parameters, use the index
`@` (and make sure to double quote). In the example,
`echo "${myarray[@]}"` is equivalent to
`echo "foo" "bar"`.

To get all elements as a single parameter, concatenated by the first
character in `IFS`, use the index `*`. In the
example, `echo "${myarray[*]}"` is equivalent to
`echo "foo bar"`.

There is a known
issue with this check's handling of `local` variables,
causing ShellCheck to flag variables that were previously declared as
arrays, even if they are in different scopes.

The easiest workaround is to simply use different variable names. Alternatively, you can ignore the check.

It is also possible to satisfy ShellCheck by declaring the
`local` variable separately from assigning to it, e.g.:

```
foo () {
 local -a baz
 baz+=("foo" "bar")
 echo "${baz[@]}"
}
bar () {
 local baz # ShellCheck gets confused if these lines are merged as local baz="qux"
 baz="qux"
 echo "$baz"
}
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`{ cmd1; cmd2; } >> file` instead of individual
redirects.```
echo foo >> file
date >> file
cat stuff >> file
```
```
{
 echo foo
 date
 cat stuff
} >> file
```
Rather than adding `>> something` after every single
line, you can simply group the relevant commands and redirect the group.
So the file has to be opened and closed only once and it means a
performance gain.

This is mainly a stylistic issue, and can freely be ignored.

Note: shell traps which would ordinarily emit output to stdout or stderr on catching their condition will have output swallowed into the redirect when the trap is triggered from within the grouping.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-eq` is
for integer comparisons. Use `=` instead.*Note: This warning seems replace by SC2170 or SC2309. Removed in
V0.4.2
- 2016-01-10*

`[[ $foo -eq "Y" ]]``[[ $foo = "Y" ]]`Shells have two sets of comparison operators: for integers
(`-eq`, `-gt`, ...) and strings (`=`,
`>`, ...). ShellCheck has noticed that you're using an
integer comparison with string data.

If you are in fact comparing integers, double check your parameters.
Certain mistakes like `$$foo` or `${bar}}` can
introduce non-numeric characters into otherwise numeric arguments.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`alias whereami="echo $PWD"``alias whereami='echo $PWD'`With double quotes, this particular alias will be defined as
`echo /home/me`, so it will always print the same path. This
is rarely intended.

By using single quotes or escaping any expansions, we define the
alias as `echo $PWD`, which will be expanded when we use the
alias. This is the far more common use case.

Note that even if you expect that the variable will never change, it may still be better to quote it. This prevents a second round of evaluation later:

```
default="Can't handle failure, aborting"
trap "echo $default; exit 1" err
false
```
The trap now has a syntax error, because instead of running
`echo $default`, it runs `echo Can't handle ..`
which has an unmatched single quote. Avoid early expansion unless you're
equally comfortable putting `eval` in there.

If you don't mind that your alias definition is expanded at define time (and its result expanded again at evaluation time), you can ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`"A"B"C"` (B indicated). Did you mean
`"ABC"` or `"A\"B\"C"`?`echo "<img src="foo.png" />" > file.html`or

`export "var"="42"``echo "<img src=\"foo.png\" />" > file.html`or

`export "var=42"`This warning triggers when an unquoted literal string is found suspiciously sandwiched between two double quoted strings.

This usually indicates one of:

`<img>` example)`export`
example).Without escaping, the inner two quotes of the sandwich (the end quote of the first section and the start quote of the second section) are no-ops. The following two statements are identical, so the quotes that were intended to be part of the html output are instead removed:

```
echo "<img src="foo.png" />" > file.html
echo "<img src=foo.png />" > file.html
```
Similarly, these statements are identical, but work as intended:

```
export "var"="42"
export "var=42"
```
If you know that the quotes are ineffectual but you prefer it stylistically, you can ignore this message.

It's common not to realize that double quotes can span multiple elements, or to stylistically prefer to quote individual variables. For example, these statements are identical, but the first is laboriously and redundantly quoted:

```
http://"$user":"$password"@"$host"/"$path"
"http://$user:$password@$host/$path"
```
When ShellCheck detects the first style (i.e. the double quotes include only a single element each), it will suppress the warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`IFS=$'\t'` ?`IFS="\t"``IFS=$'\t'`or POSIX:

`IFS="$(printf '\t')"``IFS="\t"` splits on backslash and the letter "t".
`IFS=$'\t'` splits on tab.

It's extremely rare to want to split on the letter "n" or "t", rather than linefeed or tab.

See https://github.com/koalaman/shellcheck/wiki/SC1012

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`alias archive='mv "$@" /backup'``archive() { mv "$@" /backup; }`Aliases just substitute the start of a command with something else.
They therefore can't use positional parameters, such as `$1`.
Rewrite your alias as a function.

If your alias ends up quoting the value, e.g.
`alias cut_first="awk '{print \$1}'"`, you can technically ignore this error. However, you should consider
turning this alias into a more readable function instead:
`cut_first() { awk '{print $1}' "$@"; }`

You can also ignore this warning if you
intentionally referenced the positional parameters of its relevant
context, knowing that it won't refer to the parameters of the alias
itself. For example,
`alias whatisthis='echo "This is $0 -$-" #'` will show the
shell name with flags, i.e. `This is dash -smi` or
`This is bash -himBs`, and is correct usage because it does
not intend for `$0` to reflect anything related to the
`whatisthis` alias or its invocation.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`grep -q` instead of comparing output with
`[ -n .. ]`.```
if [ -n "$(find . | grep 'IMG[0-9]')" ]
then
 echo "Images found"
fi
```
```
if find . | grep -q 'IMG[0-9]'
then
 echo "Images found"
fi
```
The problematic code has to iterate the entire directory and read all matching lines into memory before making a decision.

The correct code is cleaner and stops at the first matching line, avoiding both iterating the rest of the directory and reading data into memory.

The `pipefail` bash option may interfere with this
rewrite, since the `if` will now in effect be evaluating the
statuses of all commands instead of just the last one. Be careful using
them together.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-e`
doesn't work with globs. Use a `for` loop.```
if [ -e dir/*.mp3 ]
then
 echo "There are mp3 files."
fi
```
```
for file in dir/*.mp3
do
 if [ -e "$file" ]
 then
 echo "There are mp3 files"
 break
 fi
done
```
`[ -e file* ]` only works if there's 0 or 1 matches. If
there are multiple, it becomes `[ -e file1 file2 ]`, and the
test fails.

`[[ -e file* ]]` doesn't work at all.

Instead, use a for loop to expand the glob and check each result individually.

If you are looking for the existence of a directory, do:

```
for f in /path/to/your/files*; do
 ## Check if the glob gets expanded to existing files.
 ## If not, f here will be exactly the pattern above
 ## and the exists test will evaluate to false.
 [ -e "$f" ] && echo "files do exist" || echo "files do not exist"
 ## This is all we needed to know, so we can break after the first iteration
 break
done
```
If you are sure there will only ever be exactly 0 or 1 matches -- and
`nullglob` is not enabled -- then the test happens to
work.

You may still want to consider making this assumption explicit and failing fast if it's ever violated:

```
files=( dir/file* )
[ "${#files[@]}" -ge 2 ] && exit 1
if [ -e "${files[0]}" ]
then
 echo "The file exists"
else
 echo "No such file"
fi
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`*` or separate argument.`printf "Error: %s\n" "Bad parameters: $@"``printf "Error: %s\n" "Bad parameters: $*"``printf "Error: %s\n" "Bad parameters: ${ARRAY_VAR[@]}"``printf "Error: %s\n" "Bad parameters: " "${ARRAY_VAR[@]}"`The behavior when concatenating a string and array is rarely intended. The preceding string is prefixed to the first array element, while the succeeding string is appended to the last one. The middle array elements are unaffected.

E.g., with the parameters
`foo`,`bar`,`baz`,
`"--flag=$@"` is equivalent to the three arguments
`"--flag=foo" "bar" "baz"`.

If the intention is to concatenate all the array elements into one
argument, use `$*`. This concatenates based on
`IFS`.

If the intention is to provide each array element as a separate argument, put the array expansion in its own argument.

The POSIX specified behavior of `$@` (and by extension
arrays) as part of other strings is often unexpected:

if the parameter being expanded was embedded within a word, the first field shall be joined with the beginning part of the original word and the last field shall be joined with the end part of the original word. In all other contexts the results of the expansion are unspecified. If there are no positional parameters, the expansion of '@' shall generate zero fields, even when '@' is within double-quotes; however, if the expansion is embedded within a word which contains one or more other parts that expand to a quoted null string, these null string(s) shall still produce an empty field, except that if the other parts are all within the same double-quotes as the '@', it is unspecified whether the result is zero fields or one empty field.

If you're aware of this and intend to take advantage of it, you can ignore this warning. However, you can usually also rewrite it into a less surprising form. For example, here's a wrapper script that uses this behavior to substitute certain commands by defining a function for them:

```
#!/bin/sh
fixed_fgrep() { grep -F "$@"; }
fixed_echo() { printf '%s\n' "$*"; }
fixed_seq() { echo "seq is not portable" >&2; return 1; }
if command -v "fixed_$1" > /dev/null 2>&1
then
 # shellcheck disable=SC2145 # I know how fixed_"$@" behaves and it's correct!
 fixed_"$@"
else
 "$@"
fi
```
Here's the same script without relying on this behavior:

```
#!/bin/sh
fixed_fgrep() { grep -F "$@"; }
fixed_echo() { printf '%s\n' "$*"; }
fixed_seq() { echo "seq is not portable" >&2; return 1; }
cmd="$1"
shift
if command -v "fixed_$cmd" > /dev/null 2>&1
then
 # Perhaps more straight forward with fewer surprises:
 "fixed_$cmd" "$@"
else
 "$cmd" "$@"
fi
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-o`. Use
`\( \)` to group.`find . -name '*.avi' -o -name '*.mkv' -exec cp {} /media \;``find . \( -name '*.avi' -o -name '*.mkv' \) -exec cp {} /media \;`Note that the space between `\(` and `-name` in
this command is significant.

In `find`, two predicates with no operator between them is
considered a logical, short-circuiting AND (as if using
`-a`). E.g., `-name '*.mkv' -exec ..` is the same
as `-name '*.mkv' -a -exec ..`.

`-a` has higher precedence than `-o`, so
`-name '*.avi' -o -name '*.mkv' -a -exec ..` is equivalent to
`-name '*.avi' -o \( -name '*.mkv' -a -exec .. \)`.

In other words, the problematic code means "if name matches
`*.avi`, do nothing. Otherwise, if it matches
`*.mkv`, execute a command.".

In the correct code, we use `\( \)` to group to get the
evaluation order we want. The correct code means "if name matches
`*.avi` or `*.mkv`, then execute a command", which
was what was intended.

If you're aware of this, you can either ignore this error or group to make it explicit. For example, to decompress all gz files except tar.gz, you can use:

`find . -name '*.tar.gz' -o \( -name '*.gz' -exec gzip -d {} + \)`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`PATH="$PATH:~/bin"``PATH="$PATH:$HOME/bin"`Having literal `~` in PATH is a bad idea. Bash handles it,
but nothing else does.

This means that even if you're always using Bash, you should avoid it because any invoked program that relies on PATH will effectively ignore those entries.

For example, `make` may say
`foo: Command not found` even though `foo` works
fine from the shell and Make and Bash both use the same PATH. You'll get
similar messages from any non-bash scripts invoked, and
`whereis` will come up empty.

Use `$HOME` or full path instead.

If your directory name actually contains a literal tilde, you can ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo "$RANDOM" # Does this work?````
#!/bin/sh
echo "$RANDOM" # Unsupported in sh. Produces warning.
```
or

```
#!/bin/bash
echo "$RANDOM" # Supported in bash. No warnings.
```
Different shells support different features. To give effective advice, ShellCheck needs to know which shell your script is going to run on. You will get a different numbers of warnings about different things depending on your target shell.

If you add a shebang (e.g. `#!/bin/bash` as the first
line), the OS will use this interpreter when the script is executed, and
ShellCheck will use this shell when offering advice.

If you for any reason can't or won't add a shebang, there are multiple other ways to let shellcheck know which shell you're coding for:

`-s` or `--shell`
flag, e.g. `shellcheck -s bash myfile``# shellcheck shell=ksh` before the first command in the
file.`.bash`, `.ksh` or
`.dash` extension (`.sh` will not assume
`--shell=sh` since it's so generic)None. Please either add a shebang, directive, extension or use
`-s` to maximize ShellCheck's usefulness.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$`/`${}` for numeric index, or escape it for
string.*Note: Removed in v0.5.0
- 2018-06-01*

```
# Regular array
index=42
echo $((array[$index]))
```
or

```
# Associative array
index=banana
echo $((array[$index]))
```
```
# Regular array
index=42
echo $((array[index]))
```
or

```
# Associative array
index=banana
echo $((array[\$index]))
```
For a numerically indexed array, the `$` is mostly
pointless and can be removed like in SC2004.

For associative arrays, the `$` should be escaped to avoid
accidental dereferencing:

```
declare -A array
index='$1'
array[$index]=42
echo "$(( array[$index] ))" # bash: array: bad array subscript
echo "$(( array[\$index] ))" # 42
```
None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-exec`
does not automatically invoke a shell. Use `-exec sh -c ..`
for that.`find . -type f -exec 'cat {} | wc -l' \;````
find . -type f -exec sh -c 'cat {} | wc -l' \; # Insecure
find . -type f -exec sh -c 'cat "$1" | wc -l' _ {} \; # Secure
```
Sometimes the command can also be rewritten to not require
`find` to invoke a shell:

`find . -type f -exec wc -l {} \; | cut -d ' ' -f 1`find `-exec` and `-execdir` uses
`execve(2)` style semantics, meaning it expects an executable
and zero or more arguments that should be passed to it.

It does not use `system(3)` style semantics, meaning it
does not accept a shell command as a string, to be parsed and evaluated
by the system's command interpreter.

If you want `find` to execute a shell command, you have to
specify `sh` (or `bash`) as the executable,
`-c` as first argument and your shell command as the
second.

To prevent command injection, the filename can be passed as a separate argument to sh and referenced as a positional parameter.

This warning would trigger falsely if executing a program with spaces in the path, if no other arguments were specified.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
myfunc() {
 return foo bar
}
```
```
myfunc() {
 echo foo
 echo bar
 return 0
}
```
In bash, `return` can only be used to signal success or
failure (0 = success, 1-255 = failure).

To return textual or multiple values from a function, write them to stdout and capture them with command substitution instead.

See SC2152 for more information.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
myfunc() {
 return "Hello $USER"
}
```
```
myfunc() {
 echo "Hello $USER"
 return 0
}
```
In many languages, `return` is used to return from the
function with a final result.

In sh/bash, `return` can only be used to signal success or
failure (0 = success, 1-255 = failure), more akin to
`throw/raise` in other languages.

Results should instead be written to stdout and captured:

```
message=$(myfunc)
echo "The function wrote: $message"
```
In functions that return small integers, such as getting the cpu
temperature, the value should still be written to stdout.
`return` should be reserved for error conditions, such as
"can't determine CPU temperature". Error or failure messages should be
written to stderr.

Note in particular that `return -1` is equivalent to
`return 255`, but that `return 1` is the more
canonical way of expressing the first possible error code.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
MY_VARIABLE="hello world"
echo "$MYVARIABLE"
```
```
MY_VARIABLE="hello world"
echo "$MY_VARIABLE"
```
ShellCheck has noticed that you reference a variable that is not assigned in the script, but which has a name similar to another known variable. You should verify that the variable name is spelled correctly.

Note: This error only triggers for environment variables (all uppercase variables), and only when they have names similar to another known variable in the script. If the variable is script-local, it should by convention have a lowercase name, and will in that case be caught by SC2154 whether or not it resembles another name.

If you've double checked and ensured that you did not intend to
reference the specified variable, you can disable this message with a directive. The message will also not appear for
guarded references like `${ENVVAR:-default}` or
`${ENVVAR:?Unset error message here}`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`check-unassigned-uppercase`This is an optional rule, which means that it has a special "long name" and is not enabled by default. See the optional page for more details. In short, you have to enable it with the long name instead of the "SC" code like you would with a normal rule:

`.shellcheckrc``enable=check-unassigned-uppercase # SC2154````
var=name
n=42
echo "$var_$n.jpg" # overextended
```
or

```
target="world"
echo "hello $tagret" # misspelled
```
or

`echo "Result: ${mycmd -a myfile}" # trying to execute commands````
var=name
n=42
echo "${var}_${n}.jpg"
```
or

```
target="world"
echo "hello $target"
```
or

`echo "Result: $(mycmd -a myfile)"`ShellCheck has noticed that you reference a variable that is not assigned. Double check that the variable is indeed assigned, and that the name is not misspelled.

Note: This message only triggers for variables with lowercase
characters in their name (`foo` and `kFOO` but not
`FOO`) due to the standard convention of using lowercase
variable names for unexported, local variables.

ShellCheck intentionally does not attempt to figure out runtime or
dynamic assignments like with `source "$(date +%F).sh"` or
`eval var=value`. See SC2034 for an
extended discussion of why this is the case.

If you know for a fact that the variable is set, you can use
`${var:?}` to fail if the variable is unset (or empty),
initialize it to a default value if uninitialized with
`: "${var:=}"`, or explicitly initialize/declare it with
`var=""` or `declare var`. You can also disable
the message with a directive.

POSIX - Parameter expansion:

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`export`:`export foo="$(mycmd)"````
foo="$(mycmd)"
export foo
```
In the original code, the return value of `mycmd` is
ignored, and `export` will instead always return true. This
may prevent conditionals, `set -e` and traps from working
correctly.

When first marked for export and assigned separately, the return
value of the assignment will be that of `mycmd`. This avoids
the problem.

Note that ShellCheck does not warn about masking of local read-only
variables, such as `local -r foo=$(cmd)`, even though this
also masks the return value. This is because the alternative
`local foo; foo=$(cmd); local -r foo` is repetitive and
cumbersome. To see warnings for this and many other additional cases of
suppressed exit codes, enable
`check-extra-masked-returns`.

If you intend to ignore the return value of an assignment, you can either ignore this warning or use

```
foo=$(mycmd) || true
export foo
```
Shellcheck does not warn about `export foo=bar` because
`bar` is a literal and not a command substitution with an
independent return value.

`local`:`local foo="$(mycmd)"````
local foo
foo=$(mycmd)
```
The exit status of the command is overridden by the exit status of the creation of the local variable. For example:

```
$ f() { local foo=$(false) && echo "error was hidden"; }; f
error was hidden
$ f() { local foo; foo=$(false) && echo "error was hidden"; }; f
```
`readonly`:`readonly foo="$(mycmd)"````
foo="$(mycmd)"
readonly foo
```
A serious quoting problem with dash is another reason to declare and
assign separately. Dash is the default, `/bin/sh`
shell on Ubuntu. More specifically, dash version 0.5.8-2.10 and
others cannot run these two examples:

```
f(){ local e=$1; }
f "1 2"
export g=$(printf '%s' "foo 2")
```
While this runs fine in other shells, dash doesn't treat any of these as assignments and fails both like this:

```
local: 2: bad variable name
export: 2: bad variable name
```
The direct workaround to this bug is to quote the right-hand-side of the assignment. Separating declaraction and assignment also makes this runs fine in any shell.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`find . -name '*.mp3' -exec sh -c 'i="{}"; sox "$i" "${i%.mp3}.wav"' \;``find . -name '*.mp3' -exec sh -c 'i="$1"; sox "$i" "${i%.mp3}.wav"' shell {} \;`In the problematic example, the filename is passed by injecting it into a shell string. Any shell metacharacters in the filename will be interpreted as part of the script, and not as part of the filename. This can break the script and allow arbitrary code execution exploits.

In the correct example, the filename is passed as a parameter. It
will be safely treated as literal text. Note that when using the shell
command with `-c`, the first parameter to the shell command
(in the example "shell") becomes `$0` in the shell command's
environment, where it is used e.g. in shell error messages (you can set
it to an arbitrary value, but it makes sense to set it to the shell's
name). You should not use the first parameter to the shell command as a
data processing parameter because you cannot, for example, access
`$0` via `$*` in the shell command (because
`$*` starts with `$1`), and as previously
mentioned, `$0` is used in the shell command's error
messages, which would be confusing.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-n` is always true due to literal strings.(Or: Argument to -z is always false due to literal strings. )

```
if [ "$foo " ]
then
 echo "this is always true because of the trailing space"
fi
```
```
if [ "$foo" ]
then
 echo "correctly checks value"
fi
```
Since `[ str ]` and `[ -n str ]` check that the
string is non-empty, any literal characters in the string -- including a
space character like in the example -- will cause the test to always be
true.

Equivalently, since `[ -z str ]` checks that the string is
empty, any literal character in the string will cause the test to always
be false.

Double check the string: you may have added trailing characters, or bad quotes or syntax. Some examples include:

`[ "$foo " ]` like in the example, where the space
becomes part of the string`[ "{$foo}" ]` instead of `[ "${foo}" ]`,
where the `{` becomes part of the string`[ "$foo -gt 0" ]` instead of
`[ "$foo" -gt "0" ]`, where the `-gt` becomes part
of the stringNone.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ false ]` is
true. Remove the brackets```
if [ false ]
then
 echo "triggers anyways"
fi
```
```
if false
then
 echo "never triggers"
fi
```
`[ str ]` checks whether `str` is non-empty. It
doesn't matter if `str` is `false`, it will still
be evaluated for non-emptyness.

Instead, use the command `false` which -- as the manual
puts it -- does nothing, unsuccessfully.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ 0 ]` is true. Use
`false` instead.```
if [ 0 ]
then
 echo "always triggers"
fi
```
```
if false
then
 echo "never triggers"
fi
```
`[ str ]` checks whether `str` is non-empty. It
doesn't matter if `str` is `0`, it will still be
evaluated for non-emptyness.

Instead, use the command `false` which -- as the manual
puts it -- does nothing, unsuccessfully.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ true ]`, just use `true`.```
if [ true ]
then
 echo "always triggers"
fi
```
```
if true
then
 echo "always triggers"
fi
```
This is a stylistic suggestion to use `true` instead of
`[ true ]`.

`[ true ]` seems to suggest that the value "true" is
somehow relevant to the statement. This is not the case, it doesn't
matter. You can replace it with `[ false ]` or
`[ wombat ]`, and it will still always be true:

| String | In brackets | Outside brackets |
|---|---|---|
| true | true | true |
| false | true | false |
| wombat | true | unknown command |

It's therefore better to use it without brackets, so that the "true" actually matters.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ 1 ]`, use
`true`.```
while [ 1 ]
do
 echo "infinite loop"
done
```
```
while true
do
 echo "infinite loop"
done
```
This is a stylistic suggestion to use `true` instead of
`[ 1 ]`.

`[ 1 ]` seems to suggest that the value "1" is somehow
relevant to the statement. This is not the case: it doesn't matter. You
can replace it with `[ 0 ]` or `[ wombat ]`, and
it will still always be true.

If you instead use `true`, the value is actually
considered and can be inverted by replacing with `false`.

On bash, you can also use `(( 1 ))`, which evaluates to
true much like in C. `(( 0 ))` is similarly false.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`read`
without `-r` will mangle backslashes.```
echo "Enter name:"
read name
```
```
echo "Enter name:"
read -r name
```
or

```
echo "Enter name:"
IFS= read -r name
```
By default, `read` will interpret backslashes before
spaces and line feeds (i.e. you can use backslashes in your string as an
escape character). This is rarely expected or desired.

Normally you just want to read data *including backslashes*
which are part of the input string and have no special escape meaning,
which is what `read -r` does. You should always use
`-r` unless you have a good reason not to:

-r

If this option is given, backslash does not act as an escape character.

Even with `read -r`, leading and trailing whitespace will
be stripped from the input. Although this may sometimes be desirable or
harmless it is often surprising and difficult to catch. Clearing the
`IFS` disables this behavior, so `IFS= read -r` is
generally safest.

If you want backslashes to affect field splitting and line terminators instead of being read, you can disable this message with a directive.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`FOO`. Remove `$`/`${}`
for that, or use `${var?}` to quiet.```
MYVAR=foo
export $MYVAR
```
```
MYVAR=foo
export MYVAR
```
`export` takes a variable name, but shellcheck has noticed
that you give it an expanded variable instead. The problematic code does
not export `MYVAR` but a variable called `foo` if
any.

If this is intentional and you do want to export `foo`
instead of `MYVAR`, you can either use a directive:

```
# shellcheck disable=SC2163
export "$MYVAR"
```
Or after (but not including) version 0.4.7, take advantage of the fact that ShellCheck only warns when no parameter expansion modifiers are applied:

```
export "${MYVAR}" # ShellCheck warns
export "${MYVAR?}" # No warning
```
`${MYVAR?}` fails when `MYVAR` is unset, which
is fine since `export` would have failed too. The main side
effect is an improved runtime error message in that case.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`cd ... || exit`
in case `cd` fails.```
cd generated_files
rm -r *.c
```
```
func(){
 cd foo
 do_something
}
```
```
cd generated_files || exit
rm -r *.c
```
```
# For functions, you may want to use return:
func(){
 cd foo || return
 do_something
}
```
`cd` can fail for a variety of reasons: misspelled paths,
missing directories, missing permissions, broken symlinks and more.

If/when it does, the script will keep going and do all its operations in the wrong directory. This can be messy, especially if the operations involve creating or deleting a lot of files.

To avoid this, make sure you handle the cases when `cd`
fails. Ways to do this include

`cd foo || exit` as suggested to abort immediately, using
exit code from failed `cd` command`cd foo || { echo "Failure"; exit 1; }` abort with custom
message`cd foo || ! echo "Failure"` omitting "abort with custom
message"`if cd foo; then echo "Ok"; else echo "Fail"; fi` for
custom handling`<(cd foo && cmd)` as an alternative to
`<(cd foo || exit; cmd)` in `<(..)`,
`$(..)` or `( )`ShellCheck does not give this warning when `cd` is on the
left of a `||` or `&&`, or the condition
of a `if`, `while` or `until` loop.
Having a `set -e` command anywhere in the script will disable
this message, even though it won't necessarily prevent the issue.

If you are accounting for `cd` failures in a way
shellcheck doesn't realize, you can disable this message with a directive.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

And companion warning "This parent loop has its index variable overridden." in SC2167.

```
for((i=0; i<10; i++))
do
 for i in *
 do
 echo "$i"
 done
done
```
```
for((i=0; i<10; i++))
do
 for j in *
 do
 echo "$j"
 done
done
```
When nesting loops, especially arithmetic for loops, using the same loop variable can cause unexpected results.

In the problematic code, `i` will contain the last
filename from the inner loop, which will be interpreted as a value in
the next iteration out the outer loop. This results in either an
infinite loop or a syntax error, depending on whether the last filename
is a valid shell variable name.

In nested for-in loops, variable merely shadow each other and won't cause infinite loops or syntax errors, but reusing the variable name is rarely intentional.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ p ] && [ q ]` as `[ p -a q ]` is not
well-defined.And likewise, prefer `[ p ] || [ q ]` over
`[ p -o q ]`.

`[ "$1" = "test" -a -z "$2" ]``[ "$1" = "test" ] && [ -z "$2" ]``-a` and `-o` to mean AND and OR in a
`[ .. ]` test expression is not well defined, and can cause
incorrect results when arguments start with dashes or contain
`!`. From
POSIX:

The -a and -o binary primaries and the '(' and ')' operators have been removed. (Many expressions using them were ambiguously defined by the grammar depending on the specific expressions being evaluated.) Scripts using these expressions should be converted to the forms given below. Even though many implementations will continue to support these forms, scripts should be extremely careful when dealing with user-supplied input that could be confused with these and other primaries and operators. Unless the application developer knows all the cases that produce input to the script, invocations like:

`test "$1" -a "$2"`should be written as:

`test "$1" && test "$2"`

Using multiple `[ .. ]` expressions with shell AND/OR
operators `&&` and `||` is well defined
and therefore preferred (but note that they have equal precedence, while
`-a`/`-o` is unspecified but usually implemented
as `-a` having higher precedence).

If the shell variant being used is ksh derived (such as the bash
shell) it will have the shell builtin command `[[ ... ]]`.
This has the operators `&&`, `||`,
`(`, `)`, `!` which safely avoid the
ambiguity by noting which arguments were quoted and requiring the
operators to be unquoted (except by the `[[ ... ]]` construct
itself).

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

See companion warning SC2165.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`local` is only
valid in functions.```
local foo=bar
echo "$foo"
```
```
foo=bar
echo "$foo"
```
In Bash, `local` can only be used in functions. In other
contexts, it's an error.

It's possible to source files containing `local` from a
function context but not from any other context. This is not good
practice, but in these cases you can ignore this
error.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

*Note: This warning has been retired in favor of individual SC3xxx
warnings for each individual issue. Removed in V0.7.2
- 2021-04-19*

You are writing a script for `dash`, but you're using a
feature that `dash` doesn't support. See SC2039, the equivalent warning for POSIX sh, for
possible workarounds.

See also Ubuntu's DashAsBinSh migration guide for how to make bash-specific scripts dash-compatible.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-eq`. Use `=` to compare as string (or
use `$var` to expand as a variable).```
read -r n
if [ n -lt 0 ]
then
 echo "bad input"
fi
if [ "$USER" -eq root ]
then
 echo "You are root"
fi
```
```
read -r n
if [ "$n" -lt 0 ] # Numerical comparison
then
 echo "bad input"
fi
if [ "$USER" = root ] # String comparison
then
 echo "You are root"
fi
```
You are comparing a string value with a numerical operator, such as
`-eq`, `-ne`, `-lt` or
`-gt`. These only work for numbers.

If you want to compare the value as a string, switch to the
equivalent string operator: `=`, `!=`
`\<` or `\>`.

If you want to compare it as a number, such as
`n=42; while [ n -gt 1024/8 ]; ..`, then keep the operator
and expand the operands yourself with `$var` or
`$((expr))`: `while [ "$n" -gt $((1024/8)) ]`

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`]` outside test. Add missing `[` or
quote if intentional.`if foo -eq bar ]; then true; fi`or

`tr -d ]``if [ foo -eq bar ]; then true; fi`or

`tr -d ']'`ShellCheck found a non-test command that ends with `]` or
`]]`.

If this was intended to be a test expression like in the first
example, add the missing `[` or `[[`.

If the `]` was intended to be literal, like in
`tr -d ]`, you can quote to make this obvious.

`tr -d ]` is valid and not different from
`tr -d ']'`, so in these cases you can ignore the error
instead.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`trap myfunc 28``trap myfunc WINCH`Signal numbers can vary between platforms. Prefer signal names, which are fixed.

Signal numbers 1, 2, 3, 6, 9, 14 and 15 are specified as parts of the optional POSIX XSI and ShellCheck will not warn about these.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`trap 'echo "unkillable"' KILL`Not applicable. This is not possible.

SIGKILL and SIGSTOP can not be caught/ignored (according to POSIX and as implemented on platforms including Linux and FreeBSD). Trying to trap this signal has undefined results.

None. If you come across one, please file an issue about it.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-p`, `-m` only applies to the deepest
directory.`mkdir -p -m 0755 foo/bar/baz````
mkdir -p foo/bar/baz
chmod 0755 foo/bar/baz foo/bar foo
```
When using `-m 0755`, the mode of the directory created
will be set to 0755. When using `-p`, parent directories
which do not exist will be created, but the mode specified by
`-m` will only be used on the last directory. The parent
directories will get their access mode the default way, via umask(2).

ShellCheck does not warn if the path only has one component, as in
`mkdir -p -m 0755 mydir`, but will not attempt to determine
whether this applies for a variable as in
`mkdir -p -m 0755 "$mydir"`. You can mkdir/chmod separately
or ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`eval echo {1..$n}``eval "echo {1..$n}"`Using `eval somecommand {1..$n}` depends both on bash
silently failing to interpret the brace expansion, and on it passing
failing brace expansions literally.

Rather than depending on these questionable features (which already behave differently in other shells), use the explicit, predictable way of passing values literally: quoting.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`time`
is undefined for pipelines. time single stage or `bash -c`
instead.`time foo | bar`To time the most relevant stage:

`foo | { time bar; }`To time everything in a pipeline:

`time bash -c 'foo | bar'`Note that you can not use `time sh -c` to time an entire
pipeline, because POSIX does not guarantee that anything other than the
last stage is waited upon by the shell.

This behavior is explicitly left undefined in POSIX.

None. This warning is not emitted in `ksh` or
`bash` where `time` is defined for pipelines.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`time`
is undefined for compound commands, use `time sh -c`
instead.`time for i in *.bmp; do convert "$i" "$i.png"; done``time sh -c 'for i in *.bmp; do convert "$i" "$i.png"; done'``time` is only defined for Simple Commands by
POSIX. Timing loops, command groups and similar is not.

None. If you use a shell that supports this (e.g. bash, ksh), specify this shell in the shebang.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
flags[0]="-r"
flags[1]="--delete-after"
if [ "$dryrun" ]
then
 flags="--dry-run"
fi
```
```
flags[0]="-r"
flags[1]="--delete-after"
if [ "$dryrun" ]
then
 flags=( "--dry-run" )
fi
```
ShellCheck noticed that you have used a variable as an array, but
then assign it a string. `array=foo` is equivalent to
`array[0]=foo`, and leaves the rest of the elements
unaffected.

In the incorrect code, `"${flags[@]}"` would contain
`--dry-run` `--delete-after`.

To set an array to only a single, given element, you should use
`array=( foo )`.

In the correct code, `"${flags[@]}"` will contain
`--dry-run` only.

Another possible cause is accidentally missing the `$` on
a previous assignment: `var=(my command); var=bar` instead of
`var=$(my command); var=bar`. If the variable is not intended
to be an array, ensure that it's never assigned as one.

There is a known
issue with this check's handling of `local` variables,
causing ShellCheck to flag variables that were previously declared as
arrays, even if they are in different scopes.

The easiest workaround is to simply use different variable names. Alternatively, you can ignore the check.

It is also possible to satisfy ShellCheck by declaring the
`local` variable separately from assigning to it, e.g.:

```
foo () {
 local -a baz
 baz+=("foo" "bar")
 echo "${baz[@]}"
}
bar () {
 local baz # ShellCheck gets confused if these lines are merged as local baz="qux"
 baz="qux"
 echo "$baz"
}
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`array+=("item")` to append items to an array.```
var=(one two)
var+=three
```
```
var=(one two)
var+=( three )
```
It looks like you are trying to append a string to an array with
`var+=string`. This instead appends to the first element of
the array (equivalent to `var[0]+=three`).

In the problematic code, the array will therefore contain
`onethree` `two`.

Instead, append an array to the array with
`var+=( elements )`. This will append the new items to the
array.

In the correct code, it will contain `one`
`two` `three` as expected.

If ShellCheck mistakenly thinks the variable is an array when it's not (e.g. because the same name was used in a different context), you can ignore this error.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
foo[1][2]=bar
echo "${foo[1][2]}"
```
In bash4, consider using associative arrays:

```
declare -A foo
foo[1,2]=bar
echo "${foo[1,2]}"
```
Otherwise, do your own index arithmetic:

```
size=10
foo[1*size+2]=bar
echo "${foo[1*size+2]}"
```
Bash does not support multidimensional arrays. Rewrite it to use 1D
arrays. Associative arrays map arbitrary strings to values, and are
therefore useful since you can construct keys like `"1,2,3"`
or `"val1;val2;val3"` to index them.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`if mycmd;`, not indirectly with
`$?`.```
make mytarget
if [ $? -ne 0 ]
then
 echo "Build failed"
fi
```
```
if ! make mytarget
then
 echo "Build failed"
fi
```
For the Solaris 10 Bourne shell:

```
if make mytarget
then
 :
else
 echo "Build failed"
fi
```
Running a command and then checking its exit status `$?`
against 0 is redundant.

Instead of just checking the exit code of a command, it checks the
exit code of a command (e.g. `[`) that checks the exit code
of a command.

Apart from the redundancy, there are other reasons to avoid this pattern:

`echo "make finished"` after
`make` will cause the `if` statement to silently
start comparing `echo`'s status instead.`set -e` aka
`errexit` will exit immediately if the command fails, even
though they're followed by a clause that handles failure.`$?` is overwritten by
`[`/`[[`, so you can't get the original value in
the relevant then/else block (e.g.
`if mycmd; then echo "Success"; else echo "Failed with $?"; fi`).To check that a command returns success, use
`if mycommand; then ...`.

To check that a command returns failure, use
`if ! mycommand; then ...`. Notice that `!` will
overwrite `$?` value.

To additionally capture output with command substitution:
`if ! output=$(mycommand); then ...`

This also applies to `while`/`until` loops.

The default Solaris 10 Bourne shell does not support negating exit
statuses with `!`, so `! mycommand` tries to
invoke a utility named "!" instead. To test for failure, use
`if mycommand; then :; else ...; fi` and
`until mycommand; do ...; done`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
place="world"
printf hello $place
```
```
place="world"
printf "hello %s\n" "$place"
```
ShellCheck has noticed that you're using a `printf` with
multiple arguments, but where the first argument has no `%s`
or equivalent variable placeholders.

`echo` accepts zero or more strings to write, e.g.
`echo hello world`.

`printf` instead accepts one pattern/template with zero or
more `%s`-style placeholders, and one argument for each
placeholder.

Rewrite your command using the right semantics, otherwise all arguments after the first one will be ignored:

```
$ printf hello world\\n
hello
$ printf "hello world\n"
hello world
$ printf "hello %s\n" "world"
hello world
```
If you wanted a no-op, use `:` instead.

`: ${place=world}`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`printf "Hello %s, welcome to %s.\n" "$USER"``printf "Hello %s, welcome to %s.\n" "$USER" "$HOSTNAME"`ShellCheck has noticed that you're using a `printf` format
string with more `%s` variables than arguments to fill
them.

In the problematic example case, the last `%s` will just
become an empty string every time.

Either remove the unused variables from the format string, or add enough arguments to fill them.

When using the Ksh/Bash `%T` timestamp extension, such as
`printf 'The time is %(%H:%M)T\n'`, an argument of
`-1` and no argument are both taken to mean the current time.
In these cases, consider specifying `-1` explicitly.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`unset foo[index]``unset 'foo[index]'`Arguments to `unset` are subject to regular glob
expansion. This is especially relevant when unsetting indices in arrays,
where `[..]` is considered a glob character group.

In the problematic code, having a file called `food` in
the current directory will result in `unset foo[index]`
expanding to `unset food`, which will silently succeed
without unsetting the element.

Quoting so that the `[..]` is passed literally to
`unset` solves the issue.

Note that you can unset element using variable for index name like this:

`unset 'foo[$var]'`If you know that pathname expansion is disabled you can ignore this
message. `set -o noglob` (and variations like invoking the
script with `#!/bin/bash -f`) will prevent glob expansion of
arguments to `unset`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`.` explicitly.`find -type f``find . -type f`When not provided a search path, GNU and Busybox `find`
will use a default path of `.`, the current directory.

On POSIX, macOS/OSX, FreeBSD, OpenBSD and NetBSD, it will instead result in an error.

Explicitly specifying a path works across all implementations, and is therefore preferred.

You will get a false positive if you concatenate a series of pre-path flags:

`find -XLE .`In such cases, please either use `find -X -L -E .` or ignore the message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`tmp=$(tempfile)``tmp=$(mktemp)``tempfile` is a Debian specific utility for creating
temporary files. Its man page notes:

tempfile is deprecated; you should use mktemp(1) instead.

Neither `tempfile` nor `mktemp` are POSIX, but
`tempfile` is Debian specific while `mktemp` works
on GNU, OSX, BusyBox, *BSD and Solaris.

ShellCheck will not recognize when a function overrides this name.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`# shellcheck shell=dash` to silence.```
#!/bin/ash
echo "Hello World"
```
```
#!/bin/ash
# shellcheck shell=dash
echo "Hello World"
```
ShellCheck has no first class support for `ash`, but it
does support its Debian fork `dash` and defaults to this
whenever `ash` is specified.

Unfortunately, while the two are similar, they are not completely
compatible. For example, `ash` supports `echo -e`
but `dash` does not, so ShellCheck will incorrectly warn
about it.

You can use a directive to let ShellCheck know you're aware of this problem.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`true` as no-op).```
{
 echo "Report for $(date +%F)"
 uptime
 df -h
}
 > report.txt
```
```
{
 echo "Report for $(date +%F)"
 uptime
 df -h
} > report.txt
```
ShellCheck found a redirection that doesn't actually redirect from/to anything.

This could indicate a bug, such as in the problematic code where an
additional linefeed causes `report.txt` to be truncated
instead of containing report output, or in
`foo & > bar`, where either
`foo &> bar` or `foo > bar &` was
intended.

However, it could also be intentionally used to truncate a file or
check that it's readable. You can make this more explicit for both
ShellCheck and human readers by using `true` or
`:` as a dummy command, e.g. `true > file` or
`: > file`.

There are no semantic problems with using `> foo` over
`true > foo`, so if you don't see this as a potential
source of bugs or confusion, you can ignore it.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`|` between this redirection and the command it
should apply to.`< file.txt | grep foo``< file.txt grep foo # or more canonically: grep foo < file.txt`ShellCheck has found a stage in a pipeline that consists of a redirection but no command. This doesn't make sense because a redirection without a command will not read or write any data.

This is most likely to occur when deleting a command that had a
redirection, but leaving a `|` behind instead of moving the
redirection to a different command:

```
# Match lines with line numbers
nl < foo.txt | grep bar
# Incorrect attempt at removing line numbers. grep now has no input:
< foo.txt | grep bar
# Line numbers correctly removed. grep now reads foo.txt as intended.
grep bar < foo.txt
```
It's technically valid to do e.g.
`echo foo | > "$(cat)"` to truncate a file called "foo",
but please consider rewriting such code.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`array=( [index]=value )` .```
declare -A foo
foo=( myvalue )
```
```
declare -A foo
foo=( [key]=myvalue )
```
You appear to be initializing or appending an array element to an associative array without giving it an index. In an indexed array, elements will be auto-indexed by incremented characters. In associative arrays, the index must be given explicitly.

This could happen because of invalid spaces or otherwise malformed
index assignment, such as `array=( [key] = value )`. This
should instead be `array=( [key]=value )`.

ShellCheck may be confused when a variable name is reused in different contexts. If ShellCheck mistakenly believes the array is associative, please ignore this error.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`=` here is literal. To assign by index, use
`( [index]=value )` with no spaces. To keep as literal, quote
it.`array=( [index] = value )``array=( [index]=value )`The shell doesn't care about the `=` sign in your array
assignment because it's not part of a recognized index assignment.
Instead, it's considered a literal character and becomes part of an
array element's value.

In the example problematic code, this is because the `=`
was intended to set the index, but the shell will not recognize it when
it is surrounded by spaces.

Make sure to remove any spaces around the `=` when
assigning by index, such as in the correct code.

If you wanted the `=` to be a literal part of the array
element, add quotes around it, such as `env=( "LC_CTYPE=C" )`
or `specialChars=( "=" "%" ";" )` .

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`=` or use
`""` for empty string.`array=([1]=one [2]= two)``array=([1]=one [2]=two)`You have an array element on the form `[index]=`. The
shell will interpret this as an independent element with index
`index` and value `<empty string>`.

This may happen as part of the expression
`[index]= value`, where the space is not allowed and causes
the shell to interpret it as
`[index]="" [index+1]=value`.

If you wanted the element to have a value, remove the spaces after
`=`, e.g. `[index]=value`.

If you wanted to assign an empty string, explicitly use empty quotes:
`[index]=""`. This makes no difference to the shell, but will
make your intention clear to shellcheck and other humans.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
[ $var+1 == 5 ] # Unevaluated math
[ "{$var}" == "value" ] # Swapped around $ and {
[ "$(cmd1) | cmd2" == "42" ] # Ended with ) too soon
[[ "$var " == *.png ]] # Trailing space
```
```
[ $((var+1)) == 5 ] # Evaluated math
[ "${var}" == "value" ] # Correct variable expansion
[ "$(cmd1 | cmd2)" == "42" ] # Correct command substitution
[[ "$var" == *.png ]] # No trailing space
```
ShellCheck has determined that the two values you're comparing can never be equal.

Most of the time, this happens because of a syntax issue that introduced unintended literal characters into one of the arguments.

The left-hand side in the problematic examples will always contain (respectively) curly braces, pipe and trailing space. The right-hand sides are literal values and a pattern without trailing spaces, so they will never be equal. The statement is therefore useless, strongly indicating a bug.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$` on a variable?```
case foo in
 bar) echo "Match"
esac
```
```
case $foo in
 bar) echo "Match"
esac
```
You are using a `case` statement to compare a literal
word.

You most likely wanted to treat this word as a `$variable`
or `$(command)` instead.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
case "$var " in # Trailing space
 value) echo "Match"
esac
```
```
case "${var}" in # No trailing space
 value) echo "Match"
esac
```
ShellCheck has detected that one of the patterns in a
`case` statement will never match.

Often, this is due to mistakes in the case statement word that results in unintended literal characters. In the problematic code, there's a trailing space that will prevent the match from ever succeeding.

For more examples of when this could happen, see SC2193 for
the equivalent warning for `[[ .. ]]` statements.

Note that ShellCheck warns about individual patterns in a branch, and
will flag `*.png` in this example even though the branch is
not dead:

```
case "${img}.jpg" in
 *.png | *.jpg) echo "It's an image"
esac
```
None. If you encounter a bug and wish to ignore
this warning, make sure the directive goes in front of the
`case` and not the individual branch.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`egrep`
is non-standard and deprecated. Use `grep -E` instead.`egrep 'foo|bar' file``grep -E 'foo|bar' file``egrep` is a non-standard command. Its functionality is
provided in POSIX by `grep -E`. POSIX
grep says:

This grep has been enhanced in an upwards-compatible way to provide the exact functionality of the historical egrep and fgrep commands as well. It was the clear intention of the standard developers to consolidate the three greps into a single command.

man grep for GNU says:

Direct invocation as either egrep or fgrep is deprecated

ShellCheck will fail to recognize when functions override
`egrep`. Consider giving it a different name or ignore this error.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`fgrep`
is non-standard and deprecated. Use `grep -F` instead.`fgrep '*.*' file``grep -F '*.*' file``fgrep` is a non-standard command. Its functionality is
provided in POSIX by `grep -F`. POSIX
grep says:

This grep has been enhanced in an upwards-compatible way to provide the exact functionality of the historical egrep and fgrep commands as well. It was the clear intention of the standard developers to consolidate the three greps into a single command.

man grep for GNU says:

Direct invocation as either egrep or fgrep is deprecated

ShellCheck will fail to recognize when functions override
`fgrep`. Consider giving it a different name or ignore this error.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ ]`. Use a loop (or concatenate
with `*` instead of `@`).```
ext=png
allowedExt=(jpg bmp png)
[ "$ext" = "${allowedExt[@]}" ] && echo "Extension is valid"
```
```
ext=png
allowedExt=(jpg bmp png)
for value in "${allowedExt[@]}"
do
 [ "$ext" = "$value" ] && echo "Extension is valid"
done
```
Array expansions become a series of words in `[ .. ]`.
Operators expect single words only.

The problematic code is equivalent to
`[ "$ext" = jpg bmp png ]`, which is invalid syntax. A
typical error message is `bash: [: too many arguments` or
`dash: somefile: unexpected operator`.

Instead, use a `for` loop to iterate over values, and
apply your condition to each.

Alternatively, if you want to concatenate all the values in the array
into a single string for your test, use `"$*"` or
`"${array[*]}"`.

If you are dynamically building an a test expression, make your array
the only thing in the test expression. ShellCheck will not emit a
warning for: `set -- 1 -lt 2; [ "$@" ]`

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[[ ]]`. Use a loop (or explicit
`*` instead of `@`).```
ext=png
allowedExt=(jpg bmp png)
[[ "$ext" = "${allowedExt[@]}" ]] && echo "Extension is valid"
```
```
ext=png
allowedExt=(jpg bmp png)
for value in "${allowedExt[@]}"
do
 [[ "$ext" = "$value" ]] && echo "Extension is valid"
done
```
Array expansions in `[[ .. ]]` will implicitly concatenate
into a single string, much like in assignments. The problematic code is
equivalent to `[ "$ext" = "jpg bmp png" ]`.

Instead, use a `for` loop to iterate over values, and
apply your condition to each.

Alternatively, if you do want to concatenate all the values in the
array into a single string for your test, use `"$*"` or
`"${array[*]}"` to make this explicit.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ ]`. Use a loop.`[ "$file" = index.{htm,html,php} ] && echo "This is the main file"````
for main in index.{htm,html,php}
do
 [ "$file" = "$main" ] && echo "This is the main file"
done
```
Brace expansions in `[ ]` will expand to a sequence of
words. Operators work on single words.

The problematic code is equivalent to
`[ "$file" = index.htm index.html index.php ]`, which is
invalid syntax. A typical error message is
`bash: [: too many arguments` or
`dash: somefile: unexpected operator`.

Instead, use a `for` loop to iterate over values, and
apply your condition to each.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[[ ]]`. Use a loop.`[[ "$file" = index.{htm,html,php} ]] && echo "This is the main file"````
for main in index.{htm,html,php}
do
 [[ "$file" = "$main" ]] && echo "This is the main file"
done
```
Brace expansions doesn't happen in `[[ ]]`. They will just
be interpreted literally.

Instead, use a `for` loop to iterate over values, and
apply your condition to each.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ ]`. Use a loop.`[ current.log -nt backup/*.log ] && echo "This is the latest file"````
newerThanAll=true
for log in backup/*.log
do
 [ current.log -nt "$log" ] || newerThanAll=false
done
[ "$newerThanAll" = "true" ] && echo "This is the latest file"
```
Globs in `[ ]` will expand to a sequence of words, one per
matching filename. Meanwhile, operators work on single words.

The problematic code is equivalent to
`[ current.log -nt backup/file1.log backup/file2.log backup/file3.log ]`,
which is invalid syntax. A typical error message is
`bash: [: too many arguments` or
`dash: somefile: unexpected operator`.

Instead, use a `for` loop to iterate over matching
filenames, and apply your condition to each.

If you know your glob will only ever match one file, you can check this explicitly and use the first file:

```
set -- backup/*.log
[ $# -eq 1 ] || { echo "There are too many matches."; exit 1; }
[ file.log -nt "$1" ] && echo "This is the latest file"
```
Alternatively, ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[[ ]]` except right of
`=`/`!=`. Use a loop.`[[ current.log -nt backup/*.log ]] && echo "This is the latest file"````
newerThanAll=true
for log in backup/*.log
do
 [[ current.log -nt "$log" ]] || newerThanAll=false
done
[[ "$newerThanAll" = "true" ]] && echo "This is the latest file"
```
Globs in `[[ ]]` will not filename expand, and will be
treated literally (or as patterns on the right-hand side of
`=`, `==` and `!=`).

The problematic code is equivalent to
`[[ current.log -nt 'backup/*.png' ]`, and will look for a
file with a literal asterisk in the name.

Instead, you can iterate over the filenames you want with a loop, and apply your condition to each filename.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`(..)`
is a subshell. Did you mean `[ .. ]`, a test expression?```
if ( -d mydir )
then
 echo "It's a directory"
fi
```
```
if [ -d mydir ]
then
 echo "It's a directory"
fi
```
Tests like `-d` to see if something is a directory or
`-z` to see if it's non-empty are actually flags to the
`test` command, and only work as tests in that context.
`[` is an alias for `test`, so you'll frequently
see them written as `[ -d mydir ]`.

`( .. )` is completely unrelated, and is a subshell mostly
used to scope shell modifications. They should not be used in
`if` or `while` statements in shell scripts.

If you wanted to test a condition, rewrite the `( .. )` to
`[ .. ]`.

None.

This error is triggered by having a unary test operator as the first command name in a subshell, which won't normally happen. Note that there's a similar warning SC2205 with a higher false positive rate.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`(..)`
is a subshell. Did you mean `[ .. ]`, a test expression?```
if ( 1 -lt 2 )
then
 echo "1 is less than 2"
fi
```
```
if [ 1 -lt 2 ]
then
 echo "1 is less than 2"
fi
```
Tests like `-eq` to check numeric equality or
`\<` for string comparison only work are actually
parameters to the `test` command, and only work as tests in
that context. `[` is an alias for `test`, so
you'll frequently see them written as `[ 1 -eq 2 ]`.

`( .. )` is completely unrelated, and is a subshell mostly
used to scope shell modifications. They should not be used in
`if` or `while` statements in shell scripts.

If you wanted to test a condition, rewrite the `( .. )` to
`[ .. ]`.

This error is triggered by having a binary operator as the first
parameter in a subshell, and could falsely trigger on e.g.
`if ( grep -eq "foo|bar" file )`. In these cases, check
whether the subshell is actually needed.

Note that there's a similar looking error SC2204 with a low false positive rate.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`read -a`.`array=( $var )`If the variable should be a single element:

`array=( "$var" )`If it's multiple lines, each of which should be an element:

```
# For bash
mapfile -t array <<< "$var"
# For ksh
printf '%s\n' "$var" | while IFS="" read -r line; do array+=("$line"); done
```
If it's a line with multiple words (separated by spaces, other delimiters can be chosen with IFS), each of which should be an element:

```
# For bash
IFS=" " read -r -a array <<< "$var"
# For ksh
IFS=" " read -r -A array <<< "$var"
```
You are expanding a variable unquoted in an array. This will invoke the shell's sloppy word splitting and glob expansion.

Instead, prefer explicitly splitting (or not splitting):

`mapfile`,
`read -ra` and/or `while` loops as
appropriate.This prevents the shell from doing unwanted splitting and glob expansion, and therefore avoiding problems with data containing spaces or special characters.

If you have already taken care (through setting IFS and
`set -f`) to have word splitting work the way you intend, you
can ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`mapfile` or `read -a` to split command output (or
quote to avoid splitting).`array=( $(mycommand) )`If it outputs multiple lines, each of which should be an element:

```
# For bash 4.4+, must not be in posix mode, may use temporary files
mapfile -t array < <(mycommand)
# For bash 3.x+, must not be in posix mode, may use temporary files
array=()
while IFS='' read -r line; do array+=("$line"); done < <(mycommand)
# For ksh, and bash 4.2+ with the lastpipe option enabled (may require disabling monitor mode)
array=()
mycommand | while IFS="" read -r line; do array+=("$line"); done
```
If it outputs a line with multiple words (separated by spaces, other delimiters can be chosen with IFS), each of which should be an element:

```
# For bash, uses temporary files
IFS=" " read -r -a array <<< "$(mycommand)"
# For bash 4.2+ with the lastpipe option enabled (may require disabling monitor mode)
array=()
mycommand | IFS=" " read -r -a array
# For ksh
IFS=" " read -r -A array <<< "$(mycommand)"
```
If the output should be a single element:

`array=( "$(mycommand)" )`You are doing unquoted command expansion in an array. This will invoke the shell's sloppy word splitting and glob expansion.

Instead, prefer explicitly splitting (or not splitting):

`mapfile`, `read -ra` and/or `while`
loops as appropriate.This prevents the shell from doing unwanted splitting and glob expansion, and therefore avoiding problems with output containing spaces or special characters.

If you have already taken care (through setting IFS and
`set -f`) to have word splitting work the way you intend, you
can ignore this warning.

Another exception is the wish for error handling:
`array=( $(mycommand) ) || die-with-error` works the way it
looks while a similar `mapfile` construct like
`mapfile -t array < <(mycommand)` **doesn't
fail** and you will have to write more code for error
handling.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[[ ]]` or quote arguments to `-v` to avoid glob
expansion.`[ -v foo[0] ] ``[ -v 'foo[0]' ]`With `[`, arguments will undergo glob expansion. If a file
`foo0` exists when the problematic code is run, it will check
for the variable `foo0` instead of the array entry
`foo[0]`. If there additionally exists a `foo1`,
it will simply fail with an error.

Use `[[ ]]` or quote the argument.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`var=$(command)` to assign output (or quote to assign
string).```
user=whoami # Want to run whoami and assign output
PAGER=cat git log # Want to assign the string "cat"
```
```
user=$(whoami)
PAGER="cat" git log
```
Putting `var=` in front of a command will not assign its
output. Use `var=$(my command here)` to execute the command
and capture its output.

If you do want to assign a literal string, use quotes to make this clear to shellcheck and humans alike.

None.

Quoting a single command (as in `PAGER="cat"` above)
doesn't change how the script works. It's purely to show shellcheck (and
humans) that a literal assignment of a command name is intentional.

This warning triggers generally when a variable is assigned an
unquoted command name (from a list of hard coded names). See related
warning SC2037 which detects the same kind of error
through the patterns `var=value -flag` and
`var=value *glob*`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`if x > 5; then echo "true"; fi`or

`foo > /dev/null 2>1``if (( x > 5 )); then echo "true"; fi`or

`foo > /dev/null 2>&1`You are redirecting to or from a filename that is an integer. For
example, `ls > file` where `file` happens to be
`3`.

This is not likely to be intentional. The most common causes are:

`x > 5`. This
should instead be `[ "$x" -gt 5 ]` or
`(( x > 5 ))`.`grep -c foo file > 100` instead of
`[ "$(grep -c foo file)" -gt 100 ]``1>2` instead
of `1>&2`.If you do want to create a file named `4`, you can quote
it to silence shellcheck and make it more clear to humans that it's not
supposed to be taken numerically.

If you use the `&>` form of redirection, as in
`foo > /dev/null 2&>1`, it will trigger this
warning. You can safely ignore this warning if that is what triggered
it, or change your redirection operator to the semantically preferable
`>&`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`${..}`, array, or is it missing quoting?```
for f in $(*.png); do echo "$f"; done # Trying to loop over a glob
array=$(*.txt) # Trying to assign an array
echo "$(array[1])" # Trying to expand an array
```
```
for f in *.png; do echo "$f"; done
array=(*.txt)
echo "${array[1]}"
```
You are using a glob as a command name. This is usually a mistake caused by one of the following:

``*foo*`` or `$(*foo*)` to
expand a glob.`var=$(*.txt)` instead of `var=(*.txt)`
to assign an array.`$(..)` instead of `${..}` when
expanding an array element.Look up and double check the syntax of what you're trying to do.

None. If you want to specify a command name via glob, e.g. to not
hard code version in `./myprogram-*/foo`, expand to array or
parameters first to allow handling the cases of 0 or 2+ matches.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`false`
instead of empty `[`/`[[` conditionals.```
if [ ]
then
 echo "Temporarily disabled"
fi
```
```
if false
then
 echo "Temporarily disabled"
fi
```
`[ ]` is a somewhat obscure way of expressing falsehood,
and the behavior is likely intended to allow the incorrectly quoted
command `[ $var ]` to still work when the variable is
unset.

POSIX has a more descriptive command `false` for this.

None. This is a stylistic suggestion, and has no effect on how the script works.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-n`, but it's not handled by this
`case`.```
while getopts "vrn" n
do
 case "$n" in
 v) echo "Verbose" ;;
 r) echo "Recursive" ;;
 \?) usage;;
 esac
done
```
```
while getopts "vrn" n
do
 case "$n" in
 v) echo "Verbose" ;;
 r) echo "Recursive" ;;
 n) echo "Dry-run" ;; # -n handled here
 \?) usage;;
 esac
done
```
You have a `while getopts` loop where the corresponding
`case` statement fails to handle one of the flags.

Either add a case to handle the flag, or remove it from the
`getopts` option string.

ShellCheck may not correctly recognize less canonical uses of
`while getopts ..; do case ..;`, such as when modifying the
variable before using it:

```
while getopts "rf-:" OPT; do
 if [ "$OPT" = "-" ]; then # long option: reformulate OPT and OPTARG
 OPT="${OPTARG%%=*}" # extract long option name
 OPTARG="${OPTARG#$OPT}" # extract long option argument (may be empty)
 OPTARG="${OPTARG#=}" # if long option argument, remove assigning `=`
 fi
 case "$OPT" in
 r) ... ;;
 f) ... ;;
 my-long-option) ... ;;
 esac
done
```
In such cases you can do one of:

`getopt` (no "s") which supports
long options natively.`-)` branch.ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
while getopts "vr" n
do
 case "$n" in
 v) echo "Verbose" ;;
 r) echo "Recursive" ;;
 n) echo "Dry-run" ;;
 *) usage;;
 esac
done
```
```
while getopts "vrn" n # 'n' added here
do
 case "$n" in
 v) echo "Verbose" ;;
 r) echo "Recursive" ;;
 n) echo "Dry-run" ;;
 *) usage;;
 esac
done
```
You have a `case` statement in a
`while getopts` loop that matches a flag that hasn't been
provided in the `getopts` option string.

Either add the flag to the options list, or delete the case statement.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ .. ]`?```
if -e .bashrc
then
 echo ".bashrc already exists"
fi
```
or

```
find . -name '*.mkv'
 -exec mplayer {} \;
```
```
if [ -e .bashrc ]
then
 echo ".bashrc already exists"
fi
```
or

```
find . -name '*.mkv' \
 -exec mplayer {} \;
```
You are using a name that starts with a dash as a command name. This is almost always a bug.

There are two typical ways in which this happens:

`[ .. ]` or `[[ .. ]]` around a test
expression, like in the first example example.If you actually have a command that starts with a dash—which you should really reconsider—you can quote the name (or at least the leading dash). This makes no difference to the shell, but makes it clear to ShellCheck and humans that this is not intended as a flag.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`rm`, a command that doesn't read stdin. Wrong command or
missing `xargs`?```
ls | echo # Want to print result
cat files | rm # Want to delete items from a file
find . -type f | cp dir # Want to process 'find' output
rm file | true # Want to ignore errors
```
```
ls
cat files | while IFS= read -r file; do rm -- "$file"; done
find . -type f -exec cp {} dir \;
rm file || true
```
You are piping to one of several commands that don't read from stdin.

This may happen when:

`echo`
where `cat` was intended.`|` on the previous
line.`xargs`, because stdin should be passed as
positional parameters instead (use `xargs -0` if at all
possible).`||` instead of `|`Check your logic, and rewrite the command so data is passed correctly.

If you've overridden a command to return output, you can either rename it to make this obvious, or ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo`, a command that doesn't read stdin. Bad quoting or
missing `xargs`?```
echo << eof
 Hello World
eof
```
```
cat << eof
 Hello World
eof
```
You are redirecting to one of several commands that don't read from stdin.

This may happen when:

`echo`
where `cat` was intended.`echo <p>Hello` which tries to read from a file
`p`.`xargs`, e.g. `mv -t dir < files`
instead of `xargs mv -t dir < files` (or more safely,
`tr '\n' '\0' < files | xargs -0 mv -t dir`), because
stdin should be passed as parameters.Check your logic, and rewrite the command so data is passed correctly.

If you've overridden a command to return output, you can either rename it to make this obvious, or ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
#!/bin/sh
myfunction
myfunction() {
 echo "Hello World"
}
```
```
#!/bin/sh
myfunction() {
 echo "Hello World"
}
myfunction
```
You are calling a function that you are defining later in the file. The function definition must come first.

Function definitions are much like variable assignments, and define a name at the point the definition is "executed". This is why they must happen before their first use.

This is especially apparent when defining functions conditionally:

```
case "$(uname -s)" in
 Linux) hi() { echo "Hello from Linux"; } ;;
 Darwin) hi() { echo "Hello from macOS"; } ;;
 *) hi() { echo "Hello from something else"; } ;;
esac
hi
```
None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`let expr`, prefer `(( expr ))` .`let a++``(( a++ )) || true`Note,

`|| true`bits ignore error status code when incrementing from`0`to`1`

The `(( .. ))` arithmetic compound command evaluates
expressions in the same way as `let`, except it's not subject
to glob expansion and therefore requires no additional quoting or
escaping.

This warning only triggers in Bash/Ksh scripts. In Sh/Dash, neither
`let` nor `(( .. ))` are defined, but can be
simulated with `[ $(( expr )) -ne 0 ]` to retain exit code,
or `: $(( expr ))` to ignore it. For portability, the
`$(( expr ))` syntax is defined in POSIX standard.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`*)` case.```
#!/bin/sh
while getopts "vr" f
do
 case "$f" in
 v) echo "verbose" ;;
 r) echo "recursive" ;;
 esac
done
```
```
#!/bin/sh
while getopts "vr" f
do
 case "$f" in
 v) echo "verbose" ;;
 r) echo "recursive" ;;
 *) echo "usage: $0 [-v] [-r]" >&2
 exit 1 ;;
 esac
done
```
The `case` statement handling `getopts`
arguments does not have a default branch to handle unknown flags.

When a flag is not recognized, such as if passing `-Z` to
the example code, `getopts` will set the variable to a
literal question mark `?`. This should be handled along with
all the valid flags, usually by printing a usage message and exiting
with failure.

Using a `\?)` or `?)` case will also match
invalid flags, but`*)` would additionally match things like
the empty string if the variable name was misspelled.

If your script's logic handles unrecognized flags in another way,
e.g. after the `case` statement, you can ignore this
warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
case "$1" in
 -?) echo "Usage: $0 [-n]";;
 -n) echo "Hello World";;
 *) exit 1;;
esac
```
```
case "$1" in
 -\?) echo "Usage: $0 [-n]";;
# '-?') echo "Usage: $0 [-n]";; # Also valid
 -n) echo "Hello World";;
 *) exit 1;;
esac
```
You have specified multiple patterns in a `case`
statement, where one will always override the other. The pattern being
overridden is indicated with a SC2222 warning.

In the example, `-?` actually matches a dash followed by
any character, such as `-n`. This means that the later
`-n` branch will never trigger. In this case, the correct
solution is to escape the `-\?` so that it doesn't match
`-n`.

Another common reason for this is accidentally duplicating a branch. In this case, fix or delete the duplicate branch.

None. One could argue that having
`-*|--*) echo "Invalid flag";` is a readability issue, even
though the second pattern follows from the first. In this case, you can
either rearrange the pattern from most to least specific, i.e.
`--*|-*)` or ignore the error.

When ignoring this error, remember that ShellCheck directives have to
go in front of the `case` statement, and not in front of the
branch:

```
# shellcheck disable=SC2221,SC2222
case "$1" in
 -n) ...;;
 # no directive here
 -*|--*) echo "Unknown flag" ;;
esac
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

See companion warning SC2221.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`quote-safe-variables`This is an optional rule, which means that it has a special "long name" and is not enabled by default. See the optional page for more details. In short, you have to enable it with the long name instead of the "SC" code like you would with a normal rule:

`.shellcheckrc``enable=quote-safe-variables # SC2223``: ${COLUMNS:=80}``: "${COLUMNS:=80}"`This statement is an idiomatic way of assigning a default value to an
environment variable. However, even though it's passed to `:`
which ignores arguments, it's better to quote it.

If `COLUMNS='/*/*/*/*/*/*'`, the unquoted, problematic
code may spend 30+ minutes trashing the disk as it unnecessarily tries
to glob expand the value.

The correct code uses double quotes to avoid glob expansion, and therefore does not have this problem.

When quoting, make sure to update any inner quotes:

```
: ${var:='foo'} # Assigns foo without quotes
: "${var:='foo'}" # Assigns 'foo' with quotes
```
None, though this issue is largely theoretical.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`mv` has no destination. Check the arguments.`mv "$file $dir"``mv "$file" "$dir"`ShellCheck found an `mv` command with a single parameter.
This may be because the source and destination was accidentally merged
into a single argument, or because the line was broken in an invalid
way.

Fix the `mv` statement by correctly specifying both source
and destination.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`cp` has no destination. Check the arguments.`cp "$file $dir"``cp "$file" "$dir"`ShellCheck found a `cp` command with a single parameter.
This may be because the source and destination was accidentally merged
into a single argument, or because the line was broken in an invalid
way.

Fix the `cp` statement by correctly specifying both source
and destination.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`ln` has no destination. Check the arguments, or specify
`.` explicitly.`ln "$file $dir"`or

`ln /foo/bar/baz``ln "$file" "$dir"`or

`ln /foo/bar/baz .`ShellCheck found a `ln` command with a single parameter.
This may be because the source and destination was accidentally merged
into a single argument, because the line was broken in an invalid way,
or because you're using a non-standard invocation of `ln`
that defaults to linking the argument into the current directory.

If you wanted to specify both source and destination, fix the
`ln` statement.

If you wanted to link a file into the current directory, prefer using
the more explicit and POSIX standard invocation
`ln /your/file .`

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`find . -name '*.ppm' -exec pnmtopng {} > {}.png \;``find . -name '*.ppm' -exec sh -c 'pnmtopng "$1" > "$1.png"' _ {} \;`ShellCheck detected a `find` command with a redirection in
the middle.

This redirection may have been intended to apply only to a specific
action like `-exec` or `-print`, but it does in
fact apply to the entire `find` command:

```
# This command
find . -name '*.ppm' -exec pnmtopng {} > {}.png \;
# Is the same as this
{
 find . -name '*.ppm' -exec pnmtopng {} \;
} > {}.png
```
To perform a redirection per action, rewrite it with e.g.
`-exec sh -c '...' _ {} \;`

If the redirection is something like `> /dev/null`
where you don't mind it applying to the whole `find` and not
individual results, you can move the redirection to the end of command
to make it clear to ShellCheck (and humans) that it's not meant per
command:

```
find . -exec foo {} > /dev/null \; # Ambiguous syntax. Is it per -exec or not?
find . -exec foo {} \; > /dev/null # Identical command with clear intent.
```
There is no difference in behavior between the two.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`foo`. Remove `$`/`${}`
for that, or use `${var?}` to quiet.`read $foo``read foo``read` takes a variable name, but ShellCheck has noticed
that you give it an expanded variable instead. This will populate
whatever the variable expands to instead of the variable itself. For
example:

```
foo=bar
read $foo # Reads data into 'bar', not into 'foo'
read foo # Reads data into 'foo'
```
If this is intentional and you do want to read a variable through an indirect reference, you can silence this warning with a directive:

```
# shellcheck disable=SC2229
read "$foo"
```
Or take advantage of the fact that ShellCheck only warns when no parameter expansion modifiers are applied:

```
read "${foo}" # ShellCheck warns
read "${foo?}" # No warning
```
`${foo?}` fails when `foo` is unset, which is
fine since `read` would have failed too. The main side effect
is an improved runtime error message in that case.

`read $foo`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`which`
is non-standard. Use builtin `command -v` instead.`deprecate-which`This is an optional rule, which means that it has a special "long name" and is not enabled by default. See the optional page for more details. In short, you have to enable it with the long name instead of the "SC" code like you would with a normal rule:

`.shellcheckrc``enable=deprecate-which # SC2230``which grep````
# For the path of a single, unaliased, external command,
# or to check whether this will just "run" in this shell:
command -v grep
# To check whether commands exist, without obtaining a reusable path:
hash grep
```
`which` is a non-standard, external tool that locates an
executable in PATH. `command -v` is a POSIX standard builtin,
which uses the same lookup mechanism that the shell itself would.

This check is opt-in only in 0.7.1+, and you may choose to ignore it in earlier versions. `which` is
very common, and some prefer its executable-or-nothing behavior over
`command -v`'s handling of builtins, functions and
aliases.

`command -v`
does not check ALL parameters`command -v` succeeds (with exit code 0) if *any*
command exists:

```
# grep is in /usr/bin/grep
# foobar is not in path
#
$ command -v -- grep foobar; echo $?
0
```
In the above example, it should have failed and exited with 1 unless
*all* commands exist, if it were to be a replacement for
`which`. Other problems associated with `command`
include its inclusion of builtins, aliases, and functions.

An alternative is:

`$ hash <file1> <file2>`Which observes the standard behaviour of failures.

To obtain a path, `type -p` can be used instead. Like
`command -v`, it has a similarly quirky behavior with
builtins, aliases, and functions, although this is arguably milder since
it would print nothing for these cases. The failure condition is similar
to `hash`.

`shellcheck` issue: #1162 command
-v is not a direct replacement for which (Discussion)ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`for` loop glob to prevent word splitting,
e.g. `"${dir}"/*.txt`.```
for file in ${dir}/*.txt
do
 echo "Found ${file}"
done
```
```
for file in "${dir}"/*.txt
do
 echo "Found ${file}"
done
```
When iterating over globs containing expansions, you can still quote all expansions in the path to better handle whitespace and special characters.

Just make sure glob characters are outside quotes.
`"${dir}/*.txt"` will not glob expand, but
`"${dir}"/*.txt` or `"${dir}"/*."${ext}"`
will.

Exceptions similar to SC2086 apply. If the
variable is expected to contain globs, such as if
`dir="tmp/**"` in the example, you can ignore this
message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`sudo` with builtins like `cd`. Did you want
`sudo sh -c ..` instead?```
sudo cd /root
pwd
```
`sudo sh -c 'cd /root && pwd'`Due to the Unix process model, `sudo` can only change the
privileges of a new, external process. It can not grant privileges to a
currently running process.

This means that shell builtins -- commands that are interpreted by
the current shell rather than through program invocation -- cannot be
run with `sudo`. This includes `cd`,
`source`, `read`, and others.

Instead you can run a shell with `sudo`, and have that
shell run the builtins you want. Just be aware that what happens in that
shell stays in that shell:

```
sudo sh -c 'cd /root && pwd' # This shows /root
pwd # This shows the original directory
```
None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`(..)` around condition to avoid subshell
overhead.```
if ([ "$x" -gt 0 ])
then true; fi
```
```
if [ "$x" -gt 0 ]
then true; fi
```
The shell syntax is `if cmd`, `elif cmd`,
`while cmd` and `until cmd` without any
parentheses. Instead, parentheses are an independent construct used to
create subshells.

ShellCheck has noticed that you're wrapping `(..)` around
one or more test commands. This is unnecessary, and the resulting fork
adds quite a lot of overhead:

```
$ i=0; time while ( [ "$i" -lt 10000 ] ); do i=$((i+1)); done
real 0m6.998s
user 0m3.453s
sys 0m3.464s
$ i=0; time while [ "$i" -lt 10000 ]; do i=$((i+1)); done
real 0m0.055s
user 0m0.054s
sys 0m0.001s
```
Just delete the surrounding `(..)` since they serve no
purpose and only slows the script down.

This issue only affects performance, not correctness, so it can be safely ignored.

If you are considering doing it to stylistically match C-like
languages, please note that this is not conventional and that you'd
probably recommend someone use `if (1 == 2)` over
`if (system("[ 1 = 2 ]"))` in C no matter which language
they're used to.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`(..)` around test command to avoid subshell
overhead.`([ "$x" -gt 0 ]) && foo``[ "$x" -gt 0 ] && foo`You are wrapping a single test command in `(..)`, creating
an unnecessary subshell. This serves no purpose, but is significantly
slower:

```
$ i=0; time while ( [ "$i" -lt 10000 ] ); do i=$((i+1)); done
real 0m6.998s
user 0m3.453s
sys 0m3.464s
$ i=0; time while [ "$i" -lt 10000 ]; do i=$((i+1)); done
real 0m0.055s
user 0m0.054s
sys 0m0.001s
```
Just delete the surrounding `(..)` since they serve no
purpose and only slows the script down.

This issue only affects performance, not correctness, and can be ignored for stylistic reasons.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`{ ..; }` instead of `(..)` to avoid subshell
overhead.`([ "$x" ] || [ "$y" ]) && [ "$z" ]``{ [ "$x" ] || [ "$y" ]; } && [ "$z" ]`You appear to be using `(..)` to group test commands. This
creates a subshell, making it unnecessarily slow. Avoid this by using
`{ ..; }` to group.

Be careful to note that unlike `(..)`, this requires both
a space after the `{` and a semicolon before the
`}`.

For example, `(cmd)`, `(cmd;)` and
`( cmd )` are all valid, but `{cmd}`,
`{cmd;}` and `{ cmd }` are all syntax errors
because they lack either or both of the spaces and semicolon. The
correct form is `{ cmd; }`

Here's a small benchmark showing that the subshell version is more than 100x slower:

```
$ i=0; time for i in {1..10000}; do ([ "$x" ] || [ "$y" ]) && [ "$z" ]; done
real 0m7.122s
user 0m4.204s
sys 0m2.825s
$ i=0; time for i in {1..10000}; do { [ "$x" ] || [ "$y" ]; } && [ "$z" ]; done
real 0m0.055s
user 0m0.055s
sys 0m0.000s
```
None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-n` instead of
`! -z`.(or "Use `-z` instead of `! -n`")

```
if [ ! -n "$JAVA_HOME" ]; then echo "JAVA_HOME not specified"; fi
if [ ! -z "$STY" ]; then echo "You are already running screen"; fi
```
```
if [ -z "$JAVA_HOME" ]; then echo "JAVA_HOME not specified"; fi
if [ -n "$STY" ]; then echo "You are already running screen"; fi
```
You have negated `test -z` or `test -n`,
resulting in a needless double-negative. You can just use the other
operator instead:

```
# Identical tests to verify that a value is assigned
[ ! -z foo ] # Not has no value
[ -n foo ] # Has value
# Identical tests to verify that a value is empty
[ ! -n foo ] # Not is non-empty
[ -z foo ] # Is empty
```
This is a stylistic issue that does not affect correctness. If you prefer the original expression, you can Ignore it with a directive or flag.

`[ ! -z $var ]` might work, but
`[ -n $var]` will not. `[ -n "$var" ]` will do
what you expect.ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ -n .. ]` instead
of `! [ -z .. ]`.(or "Use `[ -z .. ]` instead of
`! [ -n .. ]`.)

```
if ! [ -n "$JAVA_HOME" ]; then echo "JAVA_HOME not specified"; fi
if ! [ -z "$STY" ]; then echo "You are already running screen"; fi
```
```
if [ -z "$JAVA_HOME" ]; then echo "JAVA_HOME not specified"; fi
if [ -n "$STY" ]; then echo "You are already running screen"; fi
```
You have negated `test -z` or `test -n`,
resulting in a needless double-negative. You can just use the other
operator instead:

```
# Identical tests to verify that a value is assigned
! [ -z foo ] # Not has no value
[ -n foo ] # Has value
# Identical tests to verify that a value is empty
! [ -n foo ] # Not is non-empty
[ -z foo ] # Is empty
```
This is a stylistic issue that does not affect correctness. If you prefer the original expression, you can Ignore it with a directive or flag.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
cat file > tr -d '\r'
cat file > rm
```
```
cat file | tr -d '\r' # tr reads stdin
cat file | xargs -d '\n' rm # rm reads arguments
```
You are using file redirection, but the filename is an unquoted command name. Instead of running the command and feeding data to it, this just writes to a file with the same name.

To run the command and feed data to it, determine how it gets its data:

`xargs` as in the second exampleNote that `xargs` has many pitfalls when it comes to
spaces and quotes. `cat file | xargs rm` will appear to work
during testing, but fails for filenames like `My File.txt` or
`Can't_Fight_This_Feeling.mp3`. The example uses the GNU
extension `-d '\n'` to more safely handle these names.

If you actually did want to write a file named after a command, simply quote the filename to let ShellCheck know you meant it literally and not as a command name. This does not change anything about how the script works:

```
# Write to a file literally named 'rm', does not try to delete anything
echo "A potentially dangerous command" > "rm"
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
#!bin/sh
echo "Hello World"
```
```
#!/bin/sh
echo "Hello World"
```
The script's interpreter, as specified in the shebang, does not start
with a `/`.

The interpreter should always be specified by absolute path to ensure that the script can be executed from any directory. When it's not, it's generally a typo like in the problematic example.

If you don't know where the interpreter is and you hoped to use
`#! bash`, this is not an option. Use
`/usr/bin/env` instead:

```
#!/usr/bin/env bash
echo "Hello World"
```
While not required by POSIX, `env` can essentially always
be found in `/usr/bin` and will search the PATH for the
specified executable.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
#!/bin/sh
. include/myscript example.com 80
```
```
#!/bin/sh
host=example.com port=80 . include/myscript
```
In Bash and Ksh, you can use `. myscript arg1 arg2..` to
set `$1` and `$2` in the sourced script.

This is not the case in Dash, where any additional arguments are ignored, or in POSIX sh where the behavior is unspecified.

Instead, assign arguments to variables and rewrite the sourced script to read from them.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`exit foo bar````
echo foo
echo bar
exit
```
In bash, `exit` can only be used to signal success or
failure (0 = success, 1-255 = failure).

To exit with textual or multiple values from a function, write them to stdout and capture them with command substitution instead.

See SC2242 for more information.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`exit "Bad filename"````
echo "Bad filename" >&2
exit 1
```
`exit` can only be used to signal success or failure (0 =
success, 1-255 = failure). It can not be used to return string data, and
it can not be used to print error messages.

String data should be written stdout, before an `exit 0`
to exit with success.

Errors should instead be written to stderr, with an
`exit 1` (or higher) to exit with failure:

```
if [ ! -f "$1" ]
then
 echo "$1 is not a regular file" >&2
 exit 1
fi
```
Note in particular that `exit -1` is equivalent to
`exit 255`, but that `exit 1` is the more
canonical way of expressing the first possible error code.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-n` to check for output (or run command without
`[`/`[[` to check for success)`avoid-nullary-conditions`This is an optional rule, which means that it has a special "long name" and is not enabled by default. See the optional page for more details. In short, you have to enable it with the long name instead of the "SC" code like you would with a normal rule:

`.shellcheckrc``enable=avoid-nullary-conditions # SC2243````
if [ "$(mycommand --myflags)" ]
then
 echo "True"
fi
```
```
# Check that the command outputs something on stdout
if [ -n "$(mycommand --myflags)" ]
then
 echo "The command had output on stdout"
fi
# Check instead that the command succeeded (exit code = 0)
if mycommand --myflags
then
 echo "The command reported success"
fi
```
(if the command instead outputs "0" or "false", see SC2244 for integer and "boolean" comparisons)

`[ "$(mycommand)" ]` is equivalent to
`[ -n "$(mycommand)" ]` and checks whether the command's
output on stdout was non-empty.

Users more familiar with other languages are often surprised to learn
that it is nothing like e.g. `if (myfunction())`, since it
does not care about what the command/function `return`s.

Using an explicit `-n` helps clarify that this is purely a
string operation. And of course, if the intention was to check whether
the command ran successfully, now would be a good time to fix it as in
the alternate example.

If you are familiar with the semantics of `[`, you can ignore this suggestion with no ill effects.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-n` to check non-empty string (or use
`=`/`-ne` to check boolean/integer).```
if [ "$1" ]
then
 echo "True"
fi
```
```
# Check if $1 is empty or non-empty
if [ -n "$1" ]
then
 echo "True, $1 is a non-empty value"
fi
# Check instead if $1 is true or false, as in Java
[ "$1" = "true" ]
# Check instead if $1 is non-zero or zero, as in C
[ "$1" -ne 0 ]
# Check instead if $1 is defined (even if just assigned the empty string) or undefined
[ "${1+x}" = "x" ]
```
`[ "$var" ]` is equivalent to `[ -n "$var" ]`
and checks that a string is non-empty.

Users more familiar with other languages are often surprised to learn
that `[ "$var" ]` is true when:

`var=false``var=0``var=null``var=" "`Adding the explicit `-n` helps clarify that this is a
string comparison, and not related to any concept of boolean values or
"truthiness" as it is in most languages.

If you are familiar with the semantics of `[`, you can ignore this stylistic suggestion with no ill
effects.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
#!/bin/ksh
if [ -f ksh* ]
then
 echo "The file exists"
fi
```
```
#!/bin/ksh
for f in ksh*
do
 if [ -f "$f" ]
 then
 echo "Found a matching file: $f"
 fi
done
```
Ksh has the curious behavior of ignoring anything after an
unrecognized flag to `test`/`[`, which means that
file checking operators against globs will effectively apply the
operator to the first expansion:

```
[ -f ksh* ] # This
[ -f ksh93u ksh93u.tar ksh93u.tar.gz ] # Becomes this
[ -f ksh93u ] # And is interpreted like this
```
This is an issue when you have multiple matches for a glob. Instead
of checking some or all, it only checks the first result and ignores the
rest. To ensure that all results are considered (either to check that
*any* or *all* results match the operator), use a loop
explicitly.

If you really only want to match the first result of the glob expansion as sorted alphabetically in the current locale, you can make this intention explicit:

```
matches=( ksh* )
if [ -f "${matches[0]}" ]
then
 echo "The first result is a file"
fi
```
If you only care that entries exists, use `-e`. ShellCheck
does not warn in this case, since all files resulting from glob
expansion necessarily exist.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
#!/bin/sh/
echo "Hello World"
```
```
#!/bin/sh
echo "Hello World"
```
The shebang specifies which file to use as an interpreter, but
probably due to some kind of typo, your script's interpreter ends in a
`/`, indicating a directory.

Ensure it points to a valid executable filename.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$` and `"` if this should be a quoted
substitution.`var=$"(whoami)"``var="$(whoami)"`ShellCheck has found a `$"(` or `$"{` . This is
most likely due to flipping the dollar-sign and double quote:

```
echo $"(cmd)" # Supposed to be "$(cmd)"
echo $"{var}" # Supposed to be "${var}"
```
Instead of quoted substitutions, these will be interpreted as
localized string resources (`$".."`) containing literal
parentheses or curly braces. If this was not intentional, you should
flip the `"` and `$` like in the example.

If you intentionally wanted a localized string literal
`$".."` that starts with `(` or `{`,
either ignore this error or start it with a
different character.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

This is an optional suggestion. It must be
explicitly enabled with a directive
`enable=quote-safe-variables` in a `# shellcheck`
comment or `.shellcheckrc`

```
subdir='example'
cd ${subdir}
```
```
subdir='example'
cd "${subdir}"
```
ShellCheck normally warns about unquoted variable use due to potential globbing or word splitting issues. See SC2086 for details. However if it is determined that a variable does not have have spaces or special characters it will omit that warning. This optional warning exists to suggest that quotes be used even in this scenario. If the code is later changed such that special characters can appear in the variable, having its use already quoted will prevent issues.

This optional warning is also helpful if ShellCheck's analysis of the variable contents is wrong because of indirect modification of the variable or because unknown commands implemented as shell functions have modified the variable.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`*)` case, even if it just exits with
error.`add-default-case`This is an optional rule, which means that it has a special "long name" and is not enabled by default. See the optional page for more details. In short, you have to enable it with the long name instead of the "SC" code like you would with a normal rule:

`.shellcheckrc``enable=add-default-case # SC2249````
case "$1" in
 start) start_service ;;
 stop) stop_service ;;
 restart|reload|force-reload)
 stop_service;
 start_service;;
esac
```
```
case "$1" in
 start) start_service ;;
 stop) stop_service ;;
 restart|reload|force-reload)
 stop_service;
 start_service;;
 *)
 echo >&2 "Invalid choice: $1"
 exit 1
esac
```
ShellCheck found a `case` statement that may not be
considering all possible cases. This may mean that only the happy paths are
accounted for.

Consider adding a default case to handle other values. If you don't know what to do or don't believe it'll ever happen, exiting with an error is good, fail-fast practice.

The example is adapted from a real world Debian init script, which due to a missing default case reports success on any misspelled command (here with underscore instead of dash):

```
$ /etc/init.d/screen-cleanup force_reload && echo success
success
```
This suggestion only triggers in verbose mode
(`-S verbose`).

If you don't have a default case because the default should be to take no action, consider adding a comment to other humans:

```
case "$(uname)" in
 CYGWIN*) cygwin=1;;
 MINGW*) mingw=1;;
 *) ;; # No special workarounds identified
esac
```
If you believe that it's impossible for the expression to have any
other value, it's considered good practice to add the equivalent of an
`assert(0)` to fail fast if this assumption should turn out
to be incorrect in the current or future versions:

```
case "$result" in
 true) proceed;;
 false) cancel;;
 *) echo >&2 "Submit bug report: '$result' should be true or false."
 exit 127
esac
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`require-variable-braces`This is an optional rule, which means that it has a special "long name" and is not enabled by default. See the optional page for more details. In short, you have to enable it with the long name instead of the "SC" code like you would with a normal rule:

`.shellcheckrc``enable=require-variable-braces # SC2250````
partial_path='example'
curl "http://example.com/$partial_path_version/explain.html"
```
```
partial_path='example'
curl "http://example.com/${partial_path}_version/explain.html"
```
If a variable gets called, and there is a string that gets appended to the variable that could get misinterpreted as possibly part of the name of the variable. Then it will not call the right variable.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`!` is not on a condition and skips errexit. Add
`|| exit 1` or make sure `$?` is checked.```
set -e
! false
rest
```
```
set -e
! false || exit 1
```
ShellCheck has found a command inverted with `!` that may
have no effect. In particular, it does not appear as a condition in an
`if` statement or `while` loop, or as the final
command in a script or function.

The most common reason for this is thinking that it'll trigger
`set -e` aka `errexit` if a command succeeds, as
in the example. This is not the case: `!` will inhibit
errexit both on success and failure of the inverted command.

Adding `|| exit` will instead exit with failure when the
command succeeds.

ShellCheck will not detect cases where `$?` is implicitly
or explicitly used to check the value afterwards:

```
set -e
check_success() { [ $? -eq 0 ] || exit 1; }
! false; check_success
! true; check_success
```
In this case, you can ignore the warning.

`set -e` and negated return
codeShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`&&` here, otherwise it's always
true.```
if [ "$1" != foo ] || [ "$1" != bar ]
then
 echo "$1 is not foo or bar"
fi
```
```
if [ "$1" != foo ] && [ "$1" != bar ]
then
 echo "$1 is not foo or bar"
fi
```
This is not a bash issue, but a simple, common logical mistake applicable to all languages.

`[ "$1" != foo ] || [ "$1" != bar ]` is always true (when
`foo != bar`):

`$1 = foo` then `$1 != bar` is true, so the
statement is true.`$1 = bar` then `$1 != foo` is true, so the
statement is true.`$1 = cow` then `$1 != foo` is true, so the
statement is true.`[ $1 != foo ] && [ $1 != bar ]` matches when
`$1` is neither `foo` nor `bar`:

`$1 = foo`, then `$1 != foo` is false, so
the statement is false.`$1 = bar`, then `$1 != bar` is false, so
the statement is false.`$1 = cow`, then both `$1 != foo` and
`$1 != bar` is true, so the statement is true.This statement is identical to
`! [ "$1" = foo ] || [ "$1" = bar ]`, which also works
correctly (by De Morgan's
law)

This warning is equivalent to SC2055 and SC2056, which trigger for intra-`test`
expressions and arithmetic contexts respectively.

Rare.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`-R` to recurse, or explicitly `a-r` to remove
read permissions.```
chmod -r 0700 dir
chmod -r file
```
```
chmod -R 0700 dir
chmod a-r file
```
Many tools use `-r` for recursive operation, but in
`chmod` this removes read permissions.

If you wanted to change permissions recursively, change the flag to
`-R`. If you wanted to remove read permissions, consider
using `a-r` explicitly to make this more obvious.

If you're using it correctly and don't mind the potential for confusion, you can save a single character by ignoring this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
case $input in
 - ) echo "Reading from stdin..." ;;
 $output ) echo "Input should be different from output" ;;
esac
```
```
case $input in
 - ) echo "Reading from stdin..." ;;
 "$output" ) echo "Input should be different from output" ;;
esac
```
When unquoted variables and command expansions are used in case branch patterns, they will be interpreted as globs.

This can lead to some surprising behavior, such as
`case $x in $x) trigger;; esac` not triggering in some cases,
such as when `x='Pride and Prejudice [1813].epub'`.

To match the literal content of the variable or expansion, make sure to double quote the expansion.

If you intended to match a dynamically generated pattern, you can ignore this suggestion with a directive.

`[[ $x = $x ]]`.ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ ]`
does not apply arithmetic evaluation. Evaluate with `$((..))`
for numbers, or use string comparator for strings.`[ 2*3 -eq array[i] ]``[ $((2*3)) -eq $((array[i])) ]`When using `[[ .. ]]` with numerical comparators
(`-eq`, `-lt`, etc), the value on either side will
be evaluated as an arithmetic expression. This means that
`2*3` will be evaluated to `6`, and `x`
will be evaluated to the contents of the variable `$x`.

When using `[ .. ]`, this does not happen.
`2*3` and `x` will both be considered invalid
numbers. Instead, use e.g. `$((2*3))` to evaluate the
expression before passing it to `[ .. ]`.

Alternatively, if the expression should be considered a string, quote
the expression and use a string comparison operator like `=`
and `!=`.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$`
and `"` if this should be a quoted substitution.```
var="foo"
echo $"var"
```
```
var="foo"
echo "$var"
```
`$".."` is a localized string, for example,
`echo $"Hello $USER"` along with the proper translation files
can be used to have the script say "Bonjour, youruser" in French
locales.

In this case, ShellCheck found a localized string whose contents is
also the name of a variable. This could have happened because the user
wanted a far more common quoted substitution, e.g. `"$var"`,
but accidentally switched the leading `$` and
`"`.

If you do want a localized string whose contents is also an active variable, you can ignore this warning or rename the variable.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`curl "$URL" > "image$((i++)).jpg"````
i=$((i+1))
curl "$URL" > "image$i.jpg"
```
You are using an arithmetic expression that modifies a variable, e.g.
`$((x+=1))` or `$((x++))`, in the name of a file
to redirect from/to, in a here document, or in a here string.

The scope of these modifications depends on whether the command itself will fork:

```
echo foo > $((var++)).txt # Updates in BusyBox and Bash
cat foo > $((var++)).txt # Updates in Busybox, not in Bash
gcc foo > $((var++)).txt # Does not update in either
gcc() { /opt/usr/bin/gcc "$@"; }
gcc foo > $((var++)).txt # Now suddenly updates in both
```
Rather than rely on knowing which commands do and don't fork, or are and aren't overridden, simply do the updates in a separate command as in the correct code.

If you know your variable is scoped the way you want it, you can ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
for f in foo, bar, baz
do
 echo "$f"
done
```
```
for f in foo bar baz
do
 echo "$f"
done
```
or

```
for f in "foo," "bar," "baz,"
do
 echo "$f"
done
```
ShellCheck found a `for` loop where the items appear to be
separated by commas. These will be treated as literal commas. If the
commas are part of the value, enclose them in quotes, or remove
them.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`gzcat yesterday.log.gz | grep "$USER" < today.log````
# Specify non-piped inputs as filenames
gzcat yesterday.log.gz | grep "$USER" - today.log
# Or merge multiple inputs into a single stream
{ gzcat yesterday.log.gz; cat today.log; } | grep "$USER"
```
A process only has a single standard input stream. Pipes and input redirections both overwrite it, so you can't use both at the same time. If you try, the redirection takes precedence and the input pipe is closed.

Many commands support specifying multiple filenames, where one can be
stdin (canonically by specifying `-` as a filename, or
alternatively by using `/dev/stdin`). In these cases, you can
rewrite the command to use one piped input, and as many extra files (or
process substitutions) as you want.

For commands that only process a single input stream (like
`tr`), you can also concatenate multiple commands or files
into a single stream using a `{ command group; }` as in the
example.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`tee` to output to
both.`env > environment.txt | grep '^ANDROID'``env | tee environment.txt | grep '^ANDROID'`A process only has a single standard output stream. Pipes and output redirections both overwrite it, so you can't use both at the same time. If you try, the redirection takes precedence and the output pipe is closed.

If you want to dump output to a file while also piping it, use
`tee` as in the example.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`cat`, `tee`,
or pass filenames instead.(or `stdin`, or `stderr`, or
`FD 3`)

`grep foo < input1 < input2 > output1 > output2 > output3````
# Merge inputs into a single stream, write outputs individually
cat input1 input2 | grep foo | tee output1 output2 > output3
# Pass inputs as filenames, write outputs individually
grep foo input1 input2 | tee output1 output2 > output3
```
A file descriptor, whether stdin, stdout, stderr, or non-standard ones, can only point to a single file/pipe.

For input, many commands support processing multiple filenames. In
these cases you can just specify the filenames instead of redirecting.
Alternatively, you can use `cat` to merge multiple filenames
into a single stream.

For output, you can use `tee` to write to multiple output
sinks in parallel.

Zsh will automatically `cat` inputs and `tee`
outputs, but none of the shells supported by ShellCheck do.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
function checksum() {
 type md5 && alias md5sum=md5
 md5sum "$@" # This calls `md5sum`, not `md5`
}
```
```
function checksum() {
 type md5 && md5sum() { md5 "$@"; }
 md5sum "$@" # Now this would call `md5` when applicable
}
```
Alias expansion happens at parse time, which means to have an effect,
the `alias` command must be executed not just before the
alias is invoked, but before the invocation is parsed.

A shell will parse commands until it has a complete set of commands
followed by a linefeed. This includes compound commands like
`{ brace; groups; }` and
`while loops; do true; done`. Here are some examples:

```
# A single command followed by a linefeed is one unit
unit 1
# These commands are in the same parsing unit because
# there is no line feed between them
unit 2; unit 2;
# These commands are in the same parsing unit because
# they are part of the same top level brace group
{
 unit 3
 unit 3
}
# These commands are in the same parsing unit because
# there is no linefeed between the groups.
{
 unit 4
}; {
 unit 4
}
```
Any alias defined in a command in `unit 1` would not take
effect until `unit 2` and beyond. Similarly, an alias defined
in unit 2 will only take effect in unit 3 and 4.

In the problematic example, the alias is defined and used in a function. Since a function definition is a single compound command, it's considered a single parsing unit. The alias would therefore not have an effect (this is true even if the function is invoked twice, because it's only parsed once).

Does this sound confusing and counter-intuitive? It is. Save yourself the trouble and always use functions instead of aliases.

If the flagged commands are not expected to use the alias, you can
ignore this error. ShellCheck may incorrectly flag this if the alias
definition and usage were in different branches of an `if`
statement.

You can ignore this warning with a directive. All warnings may always be disabled either before the relevant command or before any outer compound commands, but in this case it's especially useful:

```
# shellcheck disable=SC2262 # Option A, before compound command
if true
then
 # shellcheck disable=SC2262 # Option B, before alias command
 alias foo=bar
 # With either Option A or B, this SC2263 message is auto-suppressed
 foo
fi
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

# SC2263 – ShellCheck Wiki

 See this page on GitHub
 Sitemap

 ## Since
they're in the same parsing unit, this command will not refer to the
previously mentioned alias.

See companion warning SC2262

 ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`command`?```
ls() {
 ls --color=always "$@"
}
cd() {
 cd "$@" && ls
}
```
Note that `command` is the literal name of a shell
builtin. You should not replace it:

```
ls() {
 command ls --color=always "$@"
}
cd() {
 command cd "$@" && ls
}
```
ShellCheck found a function that immediately and unconditionally re-invokes itself, causing infinite recursion.

This generally happens when writing a wrapper function with the same name as an existing command, but forgetting to make sure it invokes the existing command and not itself. This is what happened in both of the problematic examples.

To invoke a command when a function by the same name is defined, i.e.
to suppress function lookup during execution, use the command
confusingly named `command`. For example, to run the system's
`ls` instead of the shell function `ls`, use
`command ls`.

ShellCheck does not intend to warn about infinite recursion or fork bombs in general. This warning is purely meant for unintentional bugs in well meaning wrapper functions.

If ShellCheck is triggering on an intentionally malicious fork bomb, either ignore the issue, or simply add a leading command or condition:

`:() { true && :|: & }`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`&&` for logical AND. Single `&` will
background and return true.```
if [ "$1" = "install" ] & [ "$USER" != "root" ]
then
 echo "Must be root to install"
fi
```
```
if [ "$1" = "install" ] && [ "$USER" != "root" ]
then
 echo "Must be root to install"
fi
```
ShellCheck found a `test` command followed by a
`&`. This runs the test in the background, effectively
ignoring it. To specify "logical AND" between two commands, use
`&&`.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`||` for
logical OR. Single `|` will pipe.```
if [ "$1" = "--verbose" ] | [ "$1" = "-v" ]
then
 verbose=1
fi
```
```
if [ "$1" = "--verbose" ] || [ "$1" = "-v" ]
then
 verbose=1
fi
```
ShellCheck found a `test` command followed by a
`|`. This was undoubtedly intended as a logical OR
(`||`).

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`xargs -i` is deprecated in favor of `-I{}````
# Implicit replacement string
xargs -i ls {}
# Explicit replacement string
xargs -imyfilename ls myfilename
```
```
xargs -I {} ls {}
xargs -I filename ls filename
```
`xargs -i` is a GNU specific option. It has been
deprecated in favor of the POSIX standard option `-I`.

Note that `-i` will implicitly use `{}` as a
token if nothing is specified, while `-I` requires it to be
explicit.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
[ "x$pass" = "xswordfish" ]
test x"$var" = x
```
```
[ "$pass" = "swordfish" ]
test "$var" = ""
```
Some older shells would get confused if the first argument started
with a dash, or consisted of `!` or `(`. As a
workaround, people would prefix variables and values to be compared with
`x` to ensure the left-hand side always started with an
alphanumeric character.

POSIX ensures this is not necessary, and all modern shells now follow suit.

Bash 1.14 from 1992 incorrectly fails this test. This was fixed for Bash 2.0 in 1996:

```
var='!'
[ "$var" = "!" ]
```
Dash 0.5.4 from 2007 incorrectly passes this test. This was fixed for Dash 0.5.5 in 2008:

```
x='(' y=')'
[ "$x" = "$y" ]
```
Zsh (while not supported by ShellCheck) fixed the same problem in 2015.

If you are targeting especially old shells, you can ignore this warning (or use a different letter).

`[ "x$var" = "xval" ]`?ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`var="$var"````
# If the goal is to do nothing
true
```
ShellCheck found a variable that is assigned to itself, e.g.
`x=$x`. This obviously has no effect.

Double check what the assignment was supposed to do.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`set -- first second ..`
(or use `[ ]` to compare).```
if [ -z "$1" ]
then
 $1="help"
fi
```
or

```
if $1="help"
then
 echo "Usage: $0 filename"
fi
```
```
if [ -z "$1" ]
then
 set -- "help"
fi
```
or

```
if [ $1 = "help" ]
then
 echo "Usage: $0 filename"
fi
```
You have a command on the form `$2=value`.

If the goal is to assign a new value to the positional parameters,
use the `set` builtin: `set -- one two ..` will
cause `$1` to be "one" and `$2` to be "two".

If you instead want to compare the value, use `[ ]` and
add spaces: `[ "$1" = "foo" ]`

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`declare "var$n=value"`, or (for sh)
read/eval```
n=1
var$n="hello"
```
For integer indexing in ksh/bash, consider using an indexed array:

```
n=1
var[n]="hello"
echo "${var[n]}"
```
For string indexing in ksh/bash, use an associative array:

```
typeset -A var
n="greeting"
var[$n]="hello"
echo "${var[$n]}"
```
If you actually need a variable with the constructed name in bash,
use `declare`:

```
n="Foo"
declare "var$n=42"
echo "$varFoo"
```
For `sh`, with single line contents, consider
`read`:

```
n="Foo"
read -r "var$n" << EOF
hello
EOF
echo "$varFoo"
```
or with careful escaping, `eval`:

```
n=Foo
eval "var$n='hello'"
echo "$varFoo"
```
`var$n=value` is not a valid way of assigning to a
dynamically created variable name in any shell. Please use one of the
other methods to assign to names via expanded strings. Wooledge BashFaq #6
has significantly more information on the subject.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`==`. For comparison, use
`[ "$var" = value ]`.`$a/$b==foo/bar``[ "$a/$b" = "foo/bar" ] `ShellCheck found a command name that contains a `==`. Most
likely, this was intended as a kind of comparison.

To compare two values, use `[ value1 = value2 ]`. Both the
brackets and the spaces around the `=` are relevant.

None, though you can quote the `==` to suppress the
warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`===`s found. Merge conflict or intended as a commented
border?`=======`Either resolve the merge conflict, or use `# =======` for
a border

ShellCheck found a series of `=======`s. If this was
supposed to be a border or separator, use a comment.

However, it could also be left behind from a source control merge conflict:

```
<<<<<<< HEAD
echo "Goodbye World"
=======
echo "Hello World!"
>>>>>>> mybranch
```
In this case, make sure the merge conflict is correctly resolved, and all the markers removed.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`===`. Intended as a commented border?`===================== MAIN SECTION =======================``# ===================== MAIN SECTION =======================`ShellCheck found a command that starts with a series of
`===`s. This may have been intended as a border, but is
missing the `#` to turn it into a harmless comment.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`=`. Bad line break?```
my_variable
 =value
```
`myvariable=value`ShellCheck found a command name starting with a `=`. This
was likely not meant as a new command, but instead a continuation from a
previous line.

Make sure the `=` is used correctly.

None, though you can quote the value to make ShellCheck ignore it,
e.g. `"=foo"`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`=`. Bad
assignment or comparison?```
"$var"=42
if "$var"=42
then
 true
fi
```
```
var=42
if [ "$var" = 42 ]
then
 true
fi
```
ShellCheck found a command name containing an unquoted equals sign
`=`. This was likely intended as either a comparison or an
assignment.

To compare two values, use e.g. `[ "$var" = "42" ]`

To assign a value, use e.g. `var="42"`

None, though you can quote the `=` to make ShellCheck
ignore it: `"$var=42"`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`BASH_ARGV0` to assign to `$0` in bash (or use
`[ ]` to compare).```
#!/bin/bash
$0=myscriptname
```
```
#!/bin/bash
BASH_ARGV0=myscriptname
```
You appear to be trying to assign a new value to `$0` in a
Bash script. To do this, instead assign to the special variable
`BASH_ARGV0`.

If you instead wanted to compare the value of `$0`, use a
comparison like `[ "$0" = "myname" ]`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$0`
can't be assigned in Ksh (but it does reflect the current
function).```
#!/bin/ksh
$0=myname
echo "Usage: $0 --help"
```
```
#!/bin/ksh
myname() {
 echo "Usage: $0 --help"
}
myname
```
You appear to be trying to assign a new value to `$0` in
Ksh.

This is not possible. However, `$0` will reflect the
current function name, so if you wrap your code in a function with your
chosen name, you can have `$0` expand to it.

If you instead wanted to compare the value of `$0`, use a
comparison like `[ "$0" = "myname" ]`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$0`
can't be assigned in Dash. This becomes a command name.```
#!/bin/dash
$0=myname
```
`$0` can not be changed in Dash.

You appear to be trying to assign a new value to `$0` in
Dash.

Dash does not support this. Write around it, or switch to Bash.

If you instead wanted to compare the value of `$0`, use a
comparison like `[ "$0" = "myname" ]`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$0`
can't be assigned this way, and there is no portable alternative.```
#!/bin/sh
$0=myname
```
`$0` can not be changed in a portable way.

You appear to be trying to assign a new value to `$0` in a
`sh` script.

There is no portable way to do this. Write around it, or switch to Bash.

If you instead wanted to compare the value of `$0`, use a
comparison like `[ "$0" = "myname" ]`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$`/`${}` on the left side of assignments.```
$greeting="Hello World"
${greeting}="Hello World"
```
`greeting="Hello World"`Alternatively, if the goal was to assign to a variable whose name is
in another variable (indirection), use `declare`:

```
name=foo
declare "$name=hello world"
echo "$foo"
```
Or if you actually wanted to compare the value, use a test expression:

```
if [ "$greeting" = "hello world" ]
then
 echo "Programmer, I presume?"
fi
```
Unlike Perl or PHP, `$` is not used on the left-hand side
of `=` when assigning to a variable.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`411toppm=true``_411toppm=true`You appear to be assigning to a variable name that starts with a digit. This is not allowed: variables must start with A-Z, a-z or _.

Switch to a variable name that does not start with a digit.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ ]` to compare values, or remove spaces around
`=` to assign (or quote `'='` if literal).```
# Assignment
var = value
# Comparison
if $var = value
then
 echo "Match"
fi
```
```
# Assignment
var=value
# Comparison
if [ "$var" = value ]
then
 echo "Match"
fi
```
ShellCheck found an unquoted `=` after a word.

If this was supposed to be a comparison, use square brackets:
`[ "$var" = value ]`

If this was supposed to be an assignment, remove spaces around
`=`: `var=value`

If the `=` was meant literally, quote it:

`grep '=true' file.cfg`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[ x = y ]` to compare values (or quote `'=='` if
literal).```
if $var == value
then
 echo "Match"
fi
```
```
if [ "$var" = value ]
then
 echo "Match"
fi
```
ShellCheck found an unquoted `==` after a word.

This was most likely supposed to be a comparison, so use square brackets as in the correct code.

If the `==` was supposed to be literal, you can quote it
to make ShellCheck ignore it:

`grep '===' file.js`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`+=` to assign (or quote `'+='` if
literal).`var += "my text"``var+="my text"`ShellCheck found an unquoted `+=` after a word. To append
text to a variable, remove spaces around `+=` as in the
example.

If the `+=` was supposed to be literal, you can quote it
to make ShellCheck ignore it:

`grep '+=' files..`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
jq
 ''
 file.json
```
```
jq \
 '' \
 file.json
```
ShellCheck found an empty string used as a command name. This is never valid.

If the command is intended to do nothing, use `true` aka
`:` instead. Otherwise, determine why an empty string ended
up as a command name and fix it accordingly. In the example, each line
was interpreted as a separate command due to missing line
continuations.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`df/``df /`ShellCheck found a command name that ends with `/`. Since
directories are not valid commands, this is always wrong.

The most common reason is bad quoting or escaping, such as in the example where a space was missing between a command and its argument.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

(also applies to multiple other characters like
`.,([{<>}])#"'`)

```
var=$("wget 'http://www.shellcheck.net/'")
echo Usage: $0 {start|stop|restart}
array=val1, val2, val3
```
```
var="$(wget 'http://www.shellcheck.net/')"
echo "Usage: $0 {start|stop|restart}"
array=(val1 val2 val3)
```
ShellCheck found a command name ending with a symbol, such as a comma, parenthesis, quote, or similar. This is almost always due to a syntax issue in the script.

In the examples, bad quoting and invalid array syntax caused the shell to try to run commands ending in apostrophe, curly brace, and comma, respectively.

If you have a command that *does* end in a symbol, you can
ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

(or tab)

```
'''
This script greets the planet
'''
echo "Hello World"
```
```
# This script greets the planet
echo "Hello World"
```
ShellCheck found a command name containing an especially unusual character like a tab or linefeed. This is most likely due to a syntax issue.

In the example, this was due to a Python style documentation string, which a shell will merely interpreted as a multi-line command name sandwiched between two empty strings.

If you have a command name that *does* contain a tab or
linefeed you can ignore this message, but... wow.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`export LC_ALL = "POSIX"``export LC_ALL="POSIX"`Parameters to `export`, `declare`,
`local`, `typeset` and `readonly` may
not have spaces around the `=` or `+=` operator.
This is the same as for regular variable assignments:

```
export var = value # Invalid: spaces around =
export var =value # Invalid: space before =
export var= value # Invalid: space after =
export var=value # Valid
```
This is because each individual argument to these commands is
interpreted as a string in the format `name=value`. By adding
spaces, you are instead passing the three strings `var`,
`=`, `value`, none of which follow this
format.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo Hello World``echo "Hello World"`ShellCheck found multiple unquoted spaces between words passed to
`echo`. Due to the way arguments are interpreted and passed
in the shell, these will collapse into a single space:

```
$ echo Hello World
Hello World
```
If you want to output multiple spaces, such as when creating a notice or banner, make sure the spaces are quoted, e.g. by adding (or extending) double quotes to include them:

```
$ echo "Hello World"
Hello World
```
If you're aware of this behavior and didn't want multiple spaces to show up in the output, you can either remove the unnecessary spaces or ignore this issue.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`[[ ]]`
over `[ ]` for tests in Bash/Ksh.`require-double-brackets`This is an optional rule, which means that it has a special "long name" and is not enabled by default. See the optional page for more details. In short, you have to enable it with the long name instead of the "SC" code like you would with a normal rule:

`.shellcheckrc``enable=require-double-brackets # SC2292``[ -e /etc/issue ] ``[[ -e /etc/issue ]]`ShellCheck has been explicitly asked to warn about uses of
`[ .. ]` in favor of the extended Bash/Ksh test
`[[ .. ]]`.

`[[ .. ]]` suppresses word splitting and globbing,
supports a wider variety of tests, and is generally safer and better
defined than `[ .. ]`. For an in-depth list of differences,
see the Related Resources.

This check is not enabled by default, and may have been turned on for your current project by someone who wants it enforced. You can still ignore it with a directive.

This suggestion does not trigger for Sh or Dash scripts, even when
explicitly enabled, as these shells don't support
`[[ .. ]]`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`eval "$MYCOMMAND ${@@Q}"``eval "$MYCOMMAND ${*@Q}"`ShellCheck noticed that you are calling `eval` and
including an escaped array. However, the array is passed as multiple
arguments and relies on being implicitly joined together to form a
single shell string, which `eval` can then evaluate.

Instead, prefer building your shell string with explicit string
concatenation by using `*` instead of `@` for the
index, such as `${*@Q}` or `${array[*]@Q}`.

This suggestion is equivalent to SC2124, but for
`eval` arguments rather than string variables.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
check() {
 eval "$@" || exit
}
```
```
check() {
 "$@" || exit
}
```
ShellCheck found `eval` used on an array (or equivalently,
`"$@"`). This is problematic because it effectively throws
away all boundary information and rebuilds it from shell words.

Let's say you invoke
`check sed -i '$d' "my file.txt"`:

`eval "$@"` will:

`sed -i $d my file.txt``sed`,
`-i`, `$d`, `my`
`file.txt``$d` is unset):
`sed`, `-i`, `my`,
`file.txt``sed -i 'my' 'file.txt'``"$@"` will

`sed -i '$d' 'my file.txt'`Note that while `"$@"` is essentially always better than
`eval "$@"`, it's easy to unintentionally introduce a
dependency on bad behavior through the shell debugging anti-strategy of
"adding quotes until it works":

```
# Works with problematic example because of double-escaping, fails with correct example
check ls -l "'My File.txt'"
# Works with correct example the way it was always intended:
check ls -l "My File.txt"
```
The correct example is still better, but the function invocation has to be tweaked as well.

If each of the array elements is a carefully escaped shell command or
word, use `*` instead of `@` to explicitly join
the elements on spaces which is what would happen anyways:

```
on_exit=(
 'rm /tmp/myfile; '
 'echo "Finished on $(date)" > log.txt; '
)
# Equivalent to `eval "${on_exit[@]}"`, but more explicit
eval "${on_exit[*]}"
# Even better in this case, as it does not require
# semicolons and commands don't interfere:
for cmd in "${on_exit[@]}"
do
 eval "$cmd"
done
```
If you require `eval` for another part of the command,
explicitly transform the array into a series of escaped shell words.
This ensures that the array elements will `eval` back to
themselves:

```
# Assumed to be outside of our control,
# otherwise we would output this in an array as well:
COMMAND='dialog --menu "Choose file:" 15 40 4'
# Our array:
array=(
 1 "My File.txt"
 2 "My Other File.txt"
)
eval "$COMMAND ${array[*]@Q}" # Bash 4+
eval "$COMMAND $(printf "%q " "${array[@]}")" # Bash 1+
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`${..}` need to be quoted separately, otherwise they
will match as a pattern.```
relative_path() {
 printf '%s\n' "${2#$1}"
}
# Results in "/tmp/King_Kong_[1933]/extras/trailer.mkv" because the prefix fails to match
relative_path "/tmp/King_Kong_[1933]/" "/tmp/King_Kong_[1933]/extras/trailer.mkv"
# Results in "cover.jpg" even though the prefix is different
relative_path "/tmp/King_Kong_[1933]/" "/tmp/King_Kong_3/cover.jpg"
```
```
relative_path() {
 printf '%s\n' "${2#"$1"}"
}
# Results in "extras/trailer.mkv" as expected
relative_path "/tmp/King_Kong_[1933]/" "/tmp/King_Kong_[1933]/extras/trailer.mkv"
# Results in "/tmp/King_Kong_3/cover.jpg" as expected
relative_path "/tmp/King_Kong_[1933]/" "/tmp/King_Kong_3/cover.jpg"
```
When using expansions in a parameter expansion prefix/suffix expression, the expansion needs to be quoted separately or it will match as a pattern. The quotes around the outer parameter expansion do not protect against this.

This means that any variable that contains e.g. brackets, asterisks
or question marks may not match as expected. In the example,
`[1933]` was interpreted as a pattern character range and
would therefore match `/tmp/King_Kong_3/` but not
`/tmp/King_Kong_[1933]/` as was the intention.

If you wanted to treat the string as a pattern, such as
`suffix=".*"; file="${var%$suffix}";` then you can ignore this suggestion.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`{`. Double check syntax.(or any other character)

`echo "Hello ${{name}"``echo "Hello ${name}"`ShellCheck found a parameter expansion `${something}` that
starts with an invalid character. In the example, this was caused by
accidentally duplicating the `{` in
`${{name}`.

Double check the syntax of what you're trying to do.

Some Zsh specific parameter expansions like `${(q)value}`
trigger this warning, but ShellCheck does not support Zsh.

If this warning triggers in code that works on Bash, Ksh, Dash or Sh, please submit a bug.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`${}`: `${"invalid"}` vs
`"${valid}"`.`echo ${"USER"}``echo "${USER}"`ShellCheck found a parameter expansion containing what appears to be a quoted variable name.

While the parameter expansion itself must be quoted, as in
`"${valid}"`, the quotes may not appear inside the
`{}` as in `${"invalid"}`.

Also note that translated strings like `$"Hello"` may not
use curly braces.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`${$x}`
is invalid. For expansion, use ${x}. For indirection, use arrays, ${!x}
or (for sh) eval.(or `${${x}}` is invalid)

```
# Expecting $RETRIES or 3 if unset
retries=${$RETRIES:-3}
```
or

```
mypath="/tmp/foo.txt"
var=mypath
result=${$var##*/} # Expecting ${mypath##*/}, i.e. 'foo.txt'
```
`retries=${RETRIES:-3}`or

```
mypath="/tmp/foo.txt"
var=mypath
result=${!var}
result=${result##*/}
```
ShellCheck found a parameter expansion `${..}` where the
first element was a second parameter expansion, either
`${$x..}` or `${${x}..}`. This is not valid.

In the first example, the extra `$` was unintentional and
should simply be deleted.

In the second example, `${$var##*/}` was used in the hopes
that it would expand to `${myvar##*/}` and subsequently strip
the path. This is not possible, and `var` must instead be
expanded indirectly in a separate step, before the path can be stripped
as usual. More information and other approaches can be found in the
description of SC2082.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
path="/path/to/MyFile.mp3"
echo "Playing ${${path##*/}%.*}" # Expect: Playing MyFile
```
```
path="/path/to/MyFile.mp3"
tmp=${path##*/}
echo "Playing ${tmp%.*}"
```
ShellCheck found what appears to be a nested parameter expansion. In
the example, it was hoping to combine `${var##*/}` to strip
the directory and `${var%.*}` to strip the extension.

Parameter expansions can't be nested. Use temporary variables instead, so that each parameter expansion only does a single operation.

Alternatively, if the goal is to dynamically generate and expand a variable name, see SC2082.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo "Building ${$(git rev-parse --show-toplevel)##*/}"````
tmp=$(git rev-parse --show-toplevel)
echo "Building ${tmp##*/}"
```
ShellCheck found a parameter expansion that begins with a command
substitution, such as `$(..)` or ``..``. This is
not valid. Parameter expansion only works on variables (normal or
special).

In the example, the user hoped to apply the construct
`${var##*/}`, stripping the path, to the current git root
directory as output by `git rev-parse --show-toplevel`. Since
parameter expansion only works on variable, the command substitution
must be assigned to a variable first like in the correct example.

If the goal was instead to dynamically generate a variable name to expand, see SC2082.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`echo ${"Hello World"}``echo $"Hello World"`ShellCheck found a parameter expansion `${...}` that
contains an unexpected syntax element, such as single or double
quotes.

In the example, this was due to wrapping a translated string
`$".."` with curly braces, which is not valid.

It is unclear what the intention is with invalid expansions like
`${'foo'}`, `${$"foo"}`, so please look up how to
do what you were trying to do.

If this warning triggers for working code, please submit a bug.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`"${!array[@]}"`.Plus companion warning SC2303:
`i is an array value, not a key. Use directly or loop over keys instead.`

```
array=(foo bar)
for v in "${array[@]}"
do
 echo "Value is ${array[$v]}"
done
```
Either loop over values

```
for v in "${array[@]}"
do
 echo "Value is $v"
done
```
or loop over keys:

```
for k in "${!array[@]}" # Note `!`
do
 echo "Key is $k"
 echo "Value is ${array[$k]}"
done
```
ShellCheck found a `for` loop over array *values*,
where the variable is used as an array *key*.

In the problematic example, the loop will print
`Value is foo` twice. On the second iteration,
`v=bar`, and `bar` is unset and considered zero,
so `${array[$v]}` becomes `${array[bar]}` becomes
`${array[0]}` becomes `foo`.

If you don't care about the key, simply loop over array values and
use `$v` to refer to the array value, like in the first
correct example.

If you do want the key, loop over array keys with
`"${!array[@]}"`, use `$k` to refer to the array
key, and `${array[$k]}` to refer to the array value.

If you do want to use values from the arrays as keys in the same array, you can ignore these messages with a directive:

```
declare -A fatherOf=(
 ["Eric Bloodaxe"]="Harald Fairhair"
 ["Harald Fairhair"]="Halfdan the Black"
 ["Halfdan the Black"]="Gudrød the Hunter"
 ["Gudrød the Hunter"]="Halfdan the Mild"
)
# shellcheck disable=SC2302,SC2303
for i in "${fatherOf[@]}"
do
 echo "${fatherOf[$i]:-(missing)} begat $i"
done
```
ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

Sitemap

i

See companion warning SC2302

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`*`
must be escaped to multiply: `\*`. Modern
`$((x * y))` avoids this issue.`result=$(expr 2 * 3)````
# Modern, efficient, POSIX standard approach
result=$(( 2 * 3 ))
# Older, slower approach
result=$(expr 2 \* 3)
```
ShellCheck found an `expr` command whose operator is an
unescaped asterisk `*`.

When using `expr`, each argument is expanded the same way
as for any other command. This means that `expr 2 * 3` will
turn into
`expr 2 Desktop Documents Downloads Music Pictures 3`
depending on the files in the current directory, causing an error like
`expr: syntax error: unexpected argument ‘Desktop’`

The best way to avoid this is to avoid `expr` and instead
use `$((..))` instead. If you for any reason prefer the 200x
slower, heavyweight process of forking a new process, you can escape the
`*`. Both ways are demonstrated in the correct example.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`expr "$input" : [0-9]*``expr "$input" : "[0-9]*"`ShellCheck found an `expr` command using `:` to
match a regex, but the regex is not quoted and therefore being treated
as a glob.

This means that if the problematic code is ever executed in a
directory containing a file matching `[0-9]*`, such as
`2021-reports` or `12 Angry Men [1957].mkv`, it
will be replaced be replaced and cause the command to error or
incorrectly match.

The regex should be quoted to avoid this, like in the correct example.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`f=$(expr "$c" * 9 / 5 + 32)`Prefer rewriting to a modern style (see SC2003):

`f=$((c * 9 / 5 + 32))`If you do not wish to do so, at least escape the glob characters when
passing them to `expr`:

`f=$(expr "$c" \* 9 / 5 + 32)``expr` is a command so `expr 2 * 2` will
consider `*` to mean "all files in the current directory".
This causes the expression to fail to evaluate unless you are in an
empty directory with the `failglob` and `nullglob`
options turned off.

Prefer rewriting it using the modern, POSIX standard arithmetic
expansion `$((..))`. If you do not wish to do so, you can
escape any characters like `*` to avoid the shell performing
pathname expansion on them.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
# | not escaped
expr 1 | 2
# > not escaped
expr "$foo" >= "$bar"
# Missing spaces around +
expr 1+2
# Unexpected quoting around an expression
expr "1 + 2"
```
```
expr 16 \| 7
expr "$foo" \>= "$bar"
expr 1 + 2
```
ShellCheck found an `expr` command with 1 or 2 arguments.
`expr` normally expects 3 or more.

Generally, this happens for one of two reasons:

`|`, `&`,
`>`, `>=`, `<`,
`<=`, which needs to be escaped to avoid the shell
interpreting it as a pipe, backgrounded command, or redirection.Make sure each operator or operand to `expr` is a separate
argument, and that anything containing shell metacharacters is escaped.
The correct code shows examples of each.

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`expr length`
has unspecified results. Prefer `${#var}`.or
`'expr match' has unspecified results. Prefer 'expr str : regex'.`

or
`'expr substr' has unspecified results. Prefer 'cut' or ${var#???}.`

or
`'expr index' has unspecified results. Prefer x=${var%%[chars]*}; $((${#x}+1)).`

```
# Find length of string
length=$(expr length "$var")
# Match string against regex
expr match "$input" "[0-9]*"
# Find character index in string
pos=$(expr index "$input" ":")
# Get substring by index
col2=$(expr substr "foo bar baz" 8 5)
```
```
# Find length of string
length=${#var}
# Match string against regex
expr "$input" : "[0-9]*"
# Find character index in string
pos=${input%%:*} pos=$((${#pos}+1))
# Get substring by index (bash)
str="foo bar baz"
col2="${str:7:5}"
# Get substring by index (POSIX)
col2="$(printf 'foo bar baz\n' | cut -c 8-12)"
```
You are using a `expr` with `length`,
`match`, `index`, or `substr`. These
forms did not make it into POSIX, and fail on platforms like MacOS and
FreeBSD. Consider replacing them with portable equivalents:

`length`can be trivially replaced with `${#var}`

`match`can be trivially replaced with the POSIX form
`expr str : regex`

`index`if you only need a numerical index as part of trying to extract a piece of the string, consider replacing it with parameter expansion:

```
str="mykey=myvalue"
key="${str%%=*}" # Remove everything after first =, no index required
value="${str#*=}" # Remove everything before first =, no index required
```
otherwise, you can find the index of the first `=` using
parameter expansion and string length:

```
str="mykey=myvalue"
x=${str%%=*} # Assign x="mystr"
index=$((${#x}+1)) # Add 1 to length of x
```
`substr`Extract a substring via character index is generally fragile. For example, in this example, any minor changes to the format, including just the version increasing from 8.9 to 8.10, will cause the following snippet to fail:

```
str="VIM - Vi IMproved 8.2 (2019 Dec 12, compiled Feb 15 2021 12:29:39)"
version=$(expr substr "$str" 19 3)
```
Instead, consider a different approach:

```
x="${str%% (*}" # Delete ` (` and everything after, giving "VIM - Vi IMproved 8.2"
version="${x##* }" # Delete everything before last space, giving "8.2"
# Get the fifth word separated by spaces
IFS=" " read -r _ _ _ _ version _ << EOF
$str
EOF
```
If you still want to use character index, this is trivially done in
Bash/Ksh with `${var:offset:length}` (0-based).

In POSIX, you can generally use `cut`, though be careful
if the value can contain multiple lines.

If you know your script will only run on platforms where these forms are supported, like GNU or BusyBox, you can ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
read -p "Continue? [y/n] " var
[ "$var" -eq "n" ] && exit 1
```
```
#read -p "Continue? [y/n] " var
[ "$var" = "n" ] && exit 1
```
ShellCheck found a string used as an argument to a numerical operator
like `-eq`, `-ne`, `-lt`,
`-ge`. Such strings will be treated as arithmetic
expressions, meaning `n` will refer to a variable
`$n`, and `24/12` will be evaluated into
`2`.

In the problematic example, the intention was instead to compare
`"n"` as a string, so it should use the equivalent string
operator instead, in this case `=`.

It is perfectly valid to use variables as operands. ShellCheck will not flag any value that is an unquoted variable name assigned in the script:

```
a=42; [[ "a" -eq 0 ]] # Flagged due to quotes
 [[ b -eq 0 ]] # Flagged due to not being assigned
c=42; [[ c -eq 0 ]] # Not flagged
```
However, ShellCheck does not know whether you intended
`foo/bar` to be division or a file path.

If you intended to divide `$foo` and `$bar`,
you can either make it explicit with
`[[ $((foo/bar)) -ge 0 ]]`, or simply ignore the warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`check-set-e-suppressed`This is an optional rule, which means that it has a special "long name" and is not enabled by default. See the optional page for more details. In short, you have to enable it with the long name instead of the "SC" code like you would with a normal rule:

`.shellcheckrc``enable=check-set-e-suppressed # SC2310````
#!/bin/sh
#shellcheck enable=check-set-e-suppressed
set -e
backup() {
 cp *.txt /backup
 rm *.txt # Runs even if copy fails!
}
if backup
then
 echo "Backup successful"
fi
```
```
#!/bin/sh
#shellcheck enable=check-set-e-suppressed
set -e
backup() {
 cp *.txt /backup
 rm *.txt
}
backup
echo "Backup successful"
```
ShellCheck found a function used as a condition in a script where
`set -e` is enabled. This means that the function will run
without `set -e`, and will power through any errors.

This applies to `if`, `while`, and
`until` statements, commands negated with `!`, as
well as the left-hand side of `||` and
`&&`. It does not matter how deeply the command is
nested in such a structure.

In the problematic example, the intent was that an error like
`cp: error writing '/backup/important.txt': No space left on device`
would cause the script to abort. Instead, since the function is invoked
in an `if` statement, the script will proceed to delete all
the files even though it failed to back them up.

The fix is to call it outside of an `if` statement. There
is no point in checking whether the command succeeded, since the script
would abort if it didn't. You may also want to consider replacing
`set -e` with explicit `|| exit` after every
relevant command to avoid such surprises.

If you don't care that the function runs without `set -e`,
you can disable this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

(This warning is optional and must be explicitly enabled)

```
#!/bin/bash
#shellcheck enable=check-set-e-suppressed
set -e
deploy() {
 make -j "$(nproc)" foo test
 cp ./foo /var/www/example.com/cgi-bin
 ./foo --version 2>&1
}
version=$(deploy)
echo "Successfully deployed $version"
```
```
#!/bin/bash
#shellcheck enable=check-set-e-suppressed
set -e
shopt -s inherit_errexit
deploy() {
 make -j "$(nproc)" foo test
 cp ./foo /var/www/example.com/cgi-bin
 ./foo --version 2>&1
}
version=$(deploy)
echo "Successfully deployed $version"
```
ShellCheck found a Bash function invoked in a command substitution
with `set -e` enabled.

Unlike other shells, Bash disables `set -e` in command
substitution by default. This means that the function will not exit on
error.

In the problematic code, the hope was that the function (and therefore the script) would fail if the build and test suite failed. Instead, the deployment continues even if the tests fail.

This can be fixed by either using
`version=$(set -e; deploy)` or by enabling
`inherit_errexit` as in the correct example.

If you don't care that the function runs without `set -e`,
you can ignore this warning.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`check-extra-masked-returns`This is an optional rule, which means that it has a special "long name" and is not enabled by default. See the optional page for more details. In short, you have to enable it with the long name instead of the "SC" code like you would with a normal rule:

`.shellcheckrc``enable=check-extra-masked-returns # SC2312````
set -e
cd "$(get_chroot_dir)/etc"
tar xf "${config}"
```
```
set -e
dir="$(get_chroot_dir)"
cd "${dir}/etc"
tar xf "${config}"
```
```
set -e
dir="$(get_chroot_dir)"
[[ -d "${dir}" ]] || exit 1
cd "${dir}/etc"
tar xf "${config}"
```
In the problematic example, the exit code for
`get_chroot_dir` is ignored because it is used in a command
substitution in the argument of another command.

If the command shows `error: Can't determine chroot` and
exits with failure without outputting a directory, then the command
being run will be `cd "/etc"` and the script will proceed to
overwrite the host system's configuration.

By assigning it to a variable first, the exit code of the command
will propagate into the exit code of the assignment, so that it can be
checked explicitly with `if` or implicitly with
`set -e`.

If you don't care about the command's exit status, already handle it
through a side channel like `<(cmd; echo $? > status)`,
or (in the case of background processes and process substitution) wait
on the result like `<(cmd) ; wait $!`, then you can either
ignore the suggestion with a directive, or use
`|| true` (or `|| :`) to suppress it.

Note that you can combine file
descriptor duplication with `wait` to reference process
substitution output (or input) while retaining the exit codes of those
processes. For example:

```
generate_data() {
 declare i
 for (( i = 0 ; i < 5 ; ++i ))
 do
 date -d "$RANDOM hours"
 done
}
consume_data() {
 declare line
 while IFS= read -r line
 do
 echo Consuming line: "$line"
 done
}
declare \
 input_file_descriptor \
 process
# The following statement
#
# - uses process substitution to allow us to read the output of `generate_data`
# via a filename and
# - duplicates the file descriptor for that file so that it is not
# immediately closed
#
# Note that process substitution uses either `pipe(2)` or named pipes (FIFOs)
# with `O_RDONLY` or `O_WRONLY`, and so the file descriptor that is duplicated
# via `[N]<&WORD` is only opened for reads
exec {input_file_descriptor}< <(
 generate_data
)
process=$!
# Returns non-zero if `consume_data` does
consume_data <&"$input_file_descriptor"
# Returns non-zero if `generate_data` does
wait "$process"
```
This can be particularly helpful with `readarray` for
robust array handling

https://mywiki.wooledge.org/BashPitfalls#cmd1_.26.26_cmd2_.7C.7C_cmd3

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`read -r foo[index]``read -r "foo[index]"`ShellCheck found an array element passed to read, where the
`[]` was not quoted. This means the array index
`[index]` will be treated as a glob range, and the word may
be replaced or trigger `failglob`.

In the problematic example, having a directory named
`food` will cause the command to become
`read -r food` instead, since `food` matches the
glob `foo[index]`. The result is assigning a value to the
wrong variable.

Quote or escape the pattern as shown to ensure it always reads into
the array `foo` at index `index`.

None.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`!`
does not cause a test failure.```
#!/usr/bin/env bats
@test "test" {
 # ... code
 ! test_file_exists
 # ... more code
}
```
```
#!/usr/bin/env bats
@test "test" {
 # ... code
 run ! test_file_exists
 # ... more code
}
```
Bats uses `set -e` and `trap ERR` to catch test
failures as early as possible. Although the return code of a
`!` negated command is inverted, they will never trigger
`errexit`, due to a bash design decision (see Related Resources). This means that tests
which use `!` can never fail.

Starting with bats 1.5.0 you can use `!` inside
`run`. If you are still using an older bats version, you can
rewrite `! <command>` to
`<command> && exit 1`.

The return code of the last command in the test will be the exit code
of the test function. This means that you can use
`! <command>` on the last line of the test and it will
still fail appropriately. However, you are encouraged to still use
`run !` in this case for consistency.

Stackoverflow: Why do
I need parenthesis In bash `set -e` and negated return
code

bash manpage (look
at `trap [-lp] [[arg] sigspec ...]`):

The ERR trap is not executed [...] if the command's return value is being inverted via !

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`!` does not cause a test failure. Fold the
`!` into the conditional!```
#!/usr/bin/env bats
@test "test" {
 # ... code
 ! [ $status == 0 ]
 # ... more code
}
```
```
#!/usr/bin/env bats
@test "test" {
 # ... code
 [ $status != 0 ]
 # ... more code
}
```
Bats uses `set -e` and `trap ERR` to catch test
failures as early as possible. Although the return code of a
`!` negated command is inverted, they will never trigger
`errexit`, due to a bash design decision (see Related Resources). This means that tests
which use `!` can never fail.

The return code of the last command in the test will be the exit code
of the test function. This means that you can use
`! <command>` on the last line of the test and it will
still fail appropriately. However, you are encouraged to still transform
the code in this case for consistency.

SC2314: In bats, ! does not cause a test failure (for non-conditionals)

Stackoverflow: Why do
I need parenthesis In bash `set -e` and negated return
code

bash manpage (look
at `trap [-lp] [[arg] sigspec ...]`):

The ERR trap is not executed [...] if the command's return value is being inverted via !

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`declare`
options instead.```
local readonly foo=3
readonly export bar=4
```
```
local foo=3
readonly foo
readonly bar=4
export bar
```
or

```
declare -r foo=3
declare -rx bar=4
```
In most languages, declaration modifiers like
`public`/`static`/`const` are keywords
and you can apply multiple to any declaration.

In shell scripting they are instead command names, and anything after
them is an argument. This means that `readonly local foo`
will create two readonly variables: `local`, and
`foo`. Neither will be local.

Instead, either use multiple commands, or use a single
`declare` command with appropriate flags
(`declare` will automatically make a variable local when
invoked in a function, unless `-g` is passed to explicitly
make it global).

If you want to name your variable `local`, you can quote
it as in `readonly "local"` to make your intention clear to
ShellCheck and other humans.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
usage() {
 echo >&2 "Usage: $0 -i input"
 exit 1
}
if [ "$1" = "--help" ]
then
 usage
 exit 0 # Unreachable
fi
```
```
usage() {
 echo >&2 "Usage: $0 -i input"
}
if [ "$1" = "--help" ]
then
 usage
 exit 0
fi
```
The problematic code wanted to exit with success if the user
explicitly asked for `--help`. However, since the
`usage` function already had an `exit 1`, this
statement could never run.

One possible solution is to change `usage()` to only echo,
and let callers be responsible for exiting.

ShellCheck may incorrectly believe that code is unreachable if it's invoked by variable name or in a trap. In such a case, please Ignore the message.

Note in particular that since unreachable commands may come in
clusters, it's useful to use ShellCheck's filewide or functionwide
ignore directives. A `disable` directive before a function
ignores the entire function:

```
#!/bin/bash
...
# shellcheck disable=SC2317 # Don't warn about unreachable commands in this function
start() {
 echo Starting
 /etc/init.d/foo start
}
"$1"
exit 0
```
A disable directive after the shebang, before any commands, will ignore the entire file:

```
#!/bin/bash
# Test script #1
# shellcheck disable=SC2317 # Don't warn about unreachable commands in this file
echo "Temporarily disabled"
exit 0
run-test1
run-test2
run-test3
```
Defined functions are assumed to be reachable when the script ends (not exits) since another file may source and invoke them.

You have defined two functions in the same file you are sourcing
whose names are the same but defined differently within their bodies.
Then shellcheck will state that every line of the body of the earlier
seen function definition will be unreachable which is how bash would
operate when sourcing the file. It **unclear** what
shellcheck would output if the earlier definition appeared in a
difference file that was seen first. Apparently doing a quick test. It
does **NOT** notice.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`declare`, but won't have
taken effect. Use two `declare`s.(or `local`, `typeset`, `readonly`,
`export`)

`declare -i first=$1 current=$first````
declare -i first=$1
declare -i current=$first
```
When assigning variables via a command, such as `declare`,
`typeset`, `local` etc, the expansion of all
arguments happen before all assignments. This means that you can't have
a variable assigned and then referenced in the same command.

In the example, if `$1` is 42, the arguments will first be
expanded in the current environment into
`-i first=42 current=`. They will then be passed to
`declare` which will perform the assignments.

To correctly set `current=$first` so that it uses the new
value of `first`, use two separate commands as shown.

Note that this only applies when assigning via commands, because
arguments are always expanded before commands are invoked. If assigning
without a command, as in `first=$1 current=$first`, it will
work as expected.

If you want to reference the value as it existed before the command,
e.g. if swapping variables with `declare x=$y y=$x`, you can
ignore this message. However, consider rewriting it anyways for the
benefit of any humans reading the code.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$?` refers to a condition, not a command. Assign to a
variable to avoid it being overwritten.```
mycommand
if [ $? -ne 0 ] && [ $? -ne 14 ]
then
 echo "Command failed"
fi
```
or

```
mycommand
[ $? -gt 0 ] && exit $?
```
```
mycommand
ret=$?
if [ $ret -ne 0 ] && [ $ret -ne 14 ]
then
 echo "Command failed"
fi
```
or

`mycommand || exit $?`ShellCheck found a `$?` that always refers to a condition
like `[ .. ]`, `[[ .. ]]`, or
`test`.

This most commonly happens when trying to inspect `$?`
before doing something with it, such as inspecting it again or exiting
with it, without realizing that any such inspection will also overwrite
`$?`.

In the first problematic example, `[ $? -ne 14 ]` will
never be true because it only runs after `[ $? -ne 0 ]` has
modified $? to be 0. The solution is to assign `$?` from
`mycommand` to a variable so that the variable can be
inspected repeatedly.

In the second problematic example, `exit $?` will always
`exit 0`, because it only runs if `[ $? -gt 0 ]`
returns success and sets `$?` to 0. The solution could again
be to assign `$?` to a variable first, or (as shown) use
`mycommand || exit $?` as there is no condition to overwrite
`$?`.

None. Note that ShellCheck does not warn if the usage of
`$?` after `[ .. ]` is unconditional, as in
`[ -d dir ]; return $?`.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$?` refers to echo/printf, not a previous command. Assign to
variable to avoid it being overwritten.```
mycommand
echo "Command exited with $?"
if [ $? -ne 0 ]
then
 echo "Failed"
fi
```
```
mycommand
ret=$?
echo "Command exited with $ret"
if [ $ret -ne 0 ]
then
 echo "Failed"
fi
```
ShellCheck found a `$?` that always refers to
`echo` or `printf`.

This most commonly happens when trying to show `$?` before
doing something with it, without realizing that any such action will
also overwrite `$?`.

In the problematic example,
`echo "Command exited with $?"` was intended to show the exit
code before acting on it, but the act of showing `$?` also
overwrote it, so the condition is always false. The solution is to
assign `$?` to a variable first, so that it can be used
repeatedly.

If you intentionally refer to `echo` to get the result of
a write, you can ignore this message. Alternatively, write it out as in
`if echo $$ > "$pidfile"; then status=0; else status=1; fi`

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`$((` and `))`.`values[$((i+1))]=1``values[i+1]=1`In indexed arrays (but not associative ones), the array index is
already an arithmetic context. There is no point or value in wrapping it
in an additional, explicit `$((..))`.

If ShellCheck has failed to realize that your array is associative, or if you for stylistic reasons prefer the redundancy, you can ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`((x))` is the same as `(x)`.
Prefer only one layer of parentheses.`value=$(( ((offset + index)) * 512 )) ``value=$(( (offset + index) * 512 ))`ShellCheck found doubly nested parentheses in an arithmetic
expression. While the syntax for an arithmetic expansion is
`$((..))`, this does not imply that parentheses in the
expression should also be doubled. `(x)`, `((x))`,
and `(((x)))` are all identical, so you might as well use
only one layer of parentheses.

This is a stylistic suggestion. If you prefer keeping both parentheses, you can ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`a[(x)]`
is the same as `a[x]`. Prefer not wrapping in additional
parentheses.Similarly `$(( (x) ))` is the same as
`$(( x ))`, `(( (x) ))` is the same as
`(( x ))`.

```
array[(x+1)]=val
echo $(( (x+1) ))
```
```
array[x+1]=val
echo $((x+1))
```
ShellCheck found an entire arithmetic expression wrapped in
parentheses. This does not serve a purpose since the expression is
already clearly delimited by the construct it's in, such as
`array[..]=` or `$((..))` in the example.

Note: ShellCheck does *not* warn about redundant parentheses
in subexpressions, such as `(a*b)+c`. Feel free to use
parentheses to clarify the order of operations any way you'd like.
ShellCheck only emits this suggestion when the *entire*
expression is wrapped *twice*.

If you prefer having an extra layer of parentheses for stylistic reasons, you can ignore this message.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
var=2 n=3
var+=$n
```
In bash/ksh, use an `(( arithmetic context ))`

`(( var += n ))`or declare the variable as an integer type:

```
declare -i var=2
n=4
var+=$n
```
For POSIX sh, use an `$((arithmetic expansion))`:

`var=$((var+n))`The problematic code attempts to add 2 and 3 to get 5.

Instead, `+=` on a string variable will concatenate, so
the result is 23.

If you *do* want to concatenate a number, for example to
append trailing zeroes, you can silence the warning by quoting the
number:

`var+="000"`ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
#!/bin/sh
! ! true
```
```
#!/bin/sh
true
```
POSIX (and Dash) does not allow multiple `!` pipeline
negations in a row. It's also logically unnecessary.

Use either zero or one `!`.

Scripts whose shebang declares it will run with Ksh and Bash will not trigger this warning.

If you really want to negate multiple times on POSIX or Dash, e.g. to
normalize exit codes to 0 or 1, use `cmd || false` or a
command group:

`! { ! true; } `ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

`cmd | { ! cmd; }` if necessary.`cat | ! tee /dev/full`Either negate the entire pipeline (this is equivalent unless
`pipefail` is set):

`! cat | tee /dev/full`Or use a command group to negate a single stage:

`cat | { ! tee /dev/full; }`POSIX
specifies that a status negation operator `!` is only
used to negate the status of an entire pipeline, not individual
stages.

By default the status of a pipeline is that of the last command, so
use `!` in front of the pipeline to negate as necessary.

If you have set the option `pipefail` to OR the status of
each stage together, and want to negate the status of only a single
stage, you can use negate inside a
`{ ! command group; }`.

Ksh supports `!` in front of individual pipeline stages.
ShellCheck does not warn when the shebang declares that the script will
run with Ksh.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

# SC2327 – ShellCheck Wiki

 See this page on GitHub
 Sitemap

 ## This
command substitution will be empty because the command's output gets
redirected away.

See companion warning SC2328.

 ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

(and companion warning "This command substitution will be empty because the command's output gets redirected away" which points to the relevant command substitution)

`var=$(tr -d ':' < input.txt > output.txt)````
# If the output should be captured INSTEAD OF being written to file
var=$(tr -d ':' < input.txt)
# If the output should be captured IN ADDITION to being written to file
var=$(tr -d ':' < input.txt | tee output.txt)
# If the output should NOT BE captured at all
tr -d ':' < input.txt > output.txt
```
ShellCheck has found a command substitution (`$(..)`,
``..``) that appears to never capture any output because the
command's output is being redirected.

Decide whether you want the output to be captured (by removing the
redirection), to go into wherever it's redirected (by not running the
command in a command substitution), or both (by using `tee`
to copy the output to both file and stdout).

None

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.

```
#!/bin/sh
f() {
 echo "Hello World"
}
exit
```
```
#!/bin/sh
f() {
 echo "Hello World"
}
f
exit
```
ShellCheck found a function that goes out of scope before it's ever
invoked. Verify that the function is called. It could be misspelled, or
its invocation could also be unreachable (as in
`f() { ..; }; exit; f`).

Note that if the example script did not end in `exit`,
this warning would not be emitted. This is because the function could be
invoked by another script that `source`s it.

ShellCheck is currently bad at figuring out functions that are
invoked via `trap`. In such cases, please ignore the message with a directive.

ShellCheck is a static analysis tool for shell scripts. This page is part of its documentation.
