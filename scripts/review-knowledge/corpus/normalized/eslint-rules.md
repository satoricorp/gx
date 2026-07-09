Enforce `return` statements in callbacks of array methods

💡 Suggestions

Results will be shown and updated as you type.

Rules in ESLint are grouped by type to help you understand their purpose. Each rule has emojis denoting:

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

 Some problems reported by this rule are automatically fixable by the `--fix` command line option

Some problems reported by this rule are manually fixable by editor suggestions

This rule is currently frozen and is not accepting feature requests.

These rules relate to possible logic errors in code: Enforce
💡 Suggestions
 Require
✅ Extends
 Enforce
✅ Extends
 Enforce
✅ Extends
 Disallow using an async function as a Promise executor
✅ Extends
 Disallow Disallow reassigning class members
✅ Extends
 Disallow comparing against
✅ Extends
 Disallow assignment operators in conditional expressions
✅ Extends
 Disallow reassigning
✅ Extends
 Disallow expressions where the operation doesn't affect the value
✅ Extends
 Disallow constant expressions in conditions
✅ Extends
 Disallow returning value from constructor Disallow control characters in regular expressions
✅ Extends
 Disallow the use of
✅ Extends
 Disallow duplicate arguments in
✅ Extends
 Disallow duplicate class members
✅ Extends
 Disallow duplicate conditions in if-else-if chains
✅ Extends
 Disallow duplicate keys in object literals
✅ Extends
 Disallow duplicate case labels
✅ Extends
 Disallow duplicate module imports Disallow empty character classes in regular expressions
✅ Extends
 Disallow empty destructuring patterns
✅ Extends
 Disallow reassigning exceptions in
✅ Extends
 Disallow fallthrough of
✅ Extends
 Disallow reassigning
✅ Extends
 Disallow assigning to imported bindings
✅ Extends
 Disallow variable or Disallow invalid regular expression strings in
✅ Extends
 Disallow irregular whitespace
✅ Extends
 Disallow literal numbers that lose precision
✅ Extends
 Disallow characters which are made with multiple code points in character class syntax
✅ Extends

💡 Suggestions
 Disallow
✅ Extends
 Disallow calling global object properties as functions
✅ Extends
 Disallow returning values from Promise executor functions
💡 Suggestions
 Disallow calling some
✅ Extends

💡 Suggestions
 Disallow assignments where both sides are exactly the same
✅ Extends
 Disallow comparisons where both sides are exactly the same Disallow returning values from setters
✅ Extends
 Disallow sparse arrays
✅ Extends
 Disallow template literal placeholder syntax in regular strings Disallow
✅ Extends
 Disallow
✅ Extends
 Disallow the use of undeclared variables unless mentioned in
✅ Extends
 Disallow confusing multiline expressions
✅ Extends
 Disallow unmodified loop conditions Disallow unreachable code after
✅ Extends
 Disallow loops with a body that allows only one iteration Disallow control flow statements in
✅ Extends
 Disallow negating the left operand of relational operators
✅ Extends

💡 Suggestions
 Disallow use of optional chaining in contexts where the
✅ Extends
 Disallow unused private class members
✅ Extends

💡 Suggestions
 Disallow unused variables
✅ Extends

💡 Suggestions
 Disallow the use of variables before they are defined Disallow variable assignments when the value is not used
✅ Extends
 Disallow useless backreferences in regular expressions
✅ Extends
 Disallow assignments that can lead to race conditions due to usage of Require calls to
✅ Extends

💡 Suggestions
 Enforce comparing
✅ Extends

💡 Suggestions
`return` statements in callbacks of array methods`super()` calls in constructors`for` loop update clause moving the counter in the right direction`return` statements in getters`await` inside of loops`-0``const`, `using`, and `await using` variables`debugger``function` definitions`catch` clauses`case` statements`function` declarations`function` declarations in nested blocks`RegExp` constructors`new` operators with global non-constructor functions`Object.prototype` methods directly on objects`this`/`super` before calling `super()` in constructors`let` or `var` variables that are read but never assigned`/*global */` comments`return`, `throw`, `continue`, and `break` statements`finally` blocks`undefined` value is not allowed`await` or `yield``isNaN()` when checking for `NaN``typeof` expressions against valid strings

These rules suggest alternate ways of doing things: Enforce getter and setter pairs in objects and classes ❄️ Frozen Require braces around arrow function bodies
🔧 Fix
 Enforce the use of variables within the scope they are defined ❄️ Frozen Enforce camelcase naming convention ❄️ Frozen Enforce or disallow capitalization of the first letter of a comment
🔧 Fix
 Enforce that class methods utilize Enforce a maximum cyclomatic complexity allowed in a program Require ❄️ Frozen Enforce consistent naming when capturing the current execution context ❄️ Frozen Enforce consistent brace style for all control statements
🔧 Fix
 Require Enforce ❄️ Frozen Enforce default parameters to be last ❄️ Frozen Enforce dot notation whenever possible
🔧 Fix
 Require the use of
🔧 Fix

💡 Suggestions
 ❄️ Frozen Require function names to match the name of the variable or property to which they are assigned Require or disallow named ❄️ Frozen Enforce the consistent use of either Require grouped accessor pairs in object literals and classes Require ❄️ Frozen Disallow specified identifiers ❄️ Frozen Enforce minimum and maximum identifier lengths ❄️ Frozen Require identifiers to match a specified regular expression ❄️ Frozen Require or disallow initialization in variable declarations ❄️ Frozen Require or disallow logical assignment operator shorthand
🔧 Fix

💡 Suggestions
 Enforce a maximum number of classes per file Enforce a maximum depth that blocks can be nested Enforce a maximum number of lines per file Enforce a maximum number of lines of code in a function Enforce a maximum depth that callbacks can be nested Enforce a maximum number of parameters in function definitions Enforce a maximum number of statements allowed in function blocks Require constructor names to begin with a capital letter Disallow the use of Disallow
🔧 Fix

💡 Suggestions
 Disallow bitwise operators Disallow the use of Disallow lexical declarations in case clauses
✅ Extends

💡 Suggestions
 Disallow the use of
💡 Suggestions
 ❄️ Frozen Disallow Disallow deleting variables
✅ Extends
 ❄️ Frozen Disallow equal signs explicitly at the beginning of regular expressions
🔧 Fix
 ❄️ Frozen Disallow
🔧 Fix
 Disallow empty block statements
✅ Extends

💡 Suggestions
 Disallow empty functions
💡 Suggestions
 Disallow empty static blocks
✅ Extends

💡 Suggestions
 Disallow Disallow the use of Disallow extending native types Disallow unnecessary calls to
🔧 Fix
 ❄️ Frozen Disallow unnecessary boolean casts
✅ Extends

🔧 Fix
 ❄️ Frozen Disallow unnecessary labels
🔧 Fix
 Disallow assignments to native objects or read-only global variables
✅ Extends
 ❄️ Frozen Disallow shorthand type conversions
🔧 Fix

💡 Suggestions
 Disallow declarations in the global scope Disallow the use of ❄️ Frozen Disallow inline comments after code Disallow use of Disallow the use of the ❄️ Frozen Disallow labels that share a name with a variable ❄️ Frozen Disallow labeled statements Disallow unnecessary nested blocks ❄️ Frozen Disallow
🔧 Fix
 Disallow function declarations that contain unsafe references inside loop statements ❄️ Frozen Disallow magic numbers Disallow use of chained assignment expressions ❄️ Frozen Disallow multiline strings ❄️ Frozen Disallow negated conditions ❄️ Frozen Disallow nested ternary expressions Disallow Disallow Disallow Disallow
✅ Extends

💡 Suggestions
 Disallow calls to the
💡 Suggestions
 Disallow octal literals
✅ Extends
 Disallow octal escape sequences in string literals Disallow reassigning function parameters ❄️ Frozen Disallow the unary operators Disallow the use of the Disallow variable redeclaration
✅ Extends
 Disallow multiple spaces in regular expressions
✅ Extends

🔧 Fix
 Disallow specified names in exports Disallow specified global variables Disallow specified modules when loaded by Disallow certain properties on certain objects Disallow specified syntax Disallow assignment operators in Disallow Disallow comma operators Disallow variable declarations from shadowing variables declared in the outer scope Disallow identifiers from shadowing restricted names
✅ Extends
 ❄️ Frozen Disallow ternary operators Disallow throwing literals as exceptions ❄️ Frozen Disallow initializing variables to
🔧 Fix
 ❄️ Frozen Disallow the use of ❄️ Frozen Disallow dangling underscores in identifiers ❄️ Frozen Disallow ternary operators when simpler alternatives exist
🔧 Fix
 Disallow unused expressions Disallow unused labels
✅ Extends

🔧 Fix
 Disallow unnecessary calls to Disallow unnecessary
✅ Extends
 ❄️ Frozen Disallow unnecessary computed property keys in objects and classes
🔧 Fix
 ❄️ Frozen Disallow unnecessary concatenation of literals or template literals Disallow unnecessary constructors
💡 Suggestions
 Disallow unnecessary escape characters
✅ Extends

💡 Suggestions
 Disallow renaming import, export, and destructured assignments to the same name
🔧 Fix
 Disallow redundant return statements
🔧 Fix
 Require
🔧 Fix
 ❄️ Frozen Disallow ❄️ Frozen Disallow specified warning terms in comments Disallow
✅ Extends
 ❄️ Frozen Require or disallow method and property shorthand syntax for object literals
🔧 Fix
 ❄️ Frozen Enforce variables to be declared either together or separately in functions
🔧 Fix
 ❄️ Frozen Require or disallow assignment operator shorthand where possible
🔧 Fix
 ❄️ Frozen Require using arrow functions for callbacks
🔧 Fix
 Require
🔧 Fix
 ❄️ Frozen Require destructuring from arrays and/or objects
🔧 Fix
 ❄️ Frozen Disallow the use of
🔧 Fix
 Enforce using named capture group in regular expression
💡 Suggestions
 ❄️ Frozen Disallow
🔧 Fix
 Disallow use of
🔧 Fix
 ❄️ Frozen Disallow using
🔧 Fix
 Require using Error objects as Promise rejection reasons Disallow use of the
💡 Suggestions
 Require rest parameters instead of ❄️ Frozen Require spread operators instead of ❄️ Frozen Require template literals instead of string concatenation
🔧 Fix
 Disallow losing originally caught error when re-throwing custom errors
✅ Extends

💡 Suggestions
 Enforce the use of the radix argument when using
💡 Suggestions
 Disallow async functions which have no
💡 Suggestions
 Enforce the use of
💡 Suggestions
 Require generator functions to contain
✅ Extends
 ❄️ Frozen Enforce sorted
🔧 Fix
 ❄️ Frozen Require object keys to be sorted ❄️ Frozen Require variables within the same declaration block to be sorted
🔧 Fix
 Require or disallow strict mode directives
🔧 Fix
 Require symbol descriptions ❄️ Frozen Require ❄️ Frozen Require or disallow "Yoda" conditions
🔧 Fix
`this``return` statements to either always or never specify values`default` cases in `switch` statements`default` clauses in `switch` statements to be last`===` and `!==``function` expressions`function` declarations or expressions assigned to variables`for-in` loops to include an `if` statement`alert`, `confirm`, and `prompt``Array` constructors`arguments.caller` or `arguments.callee``console``continue` statements`else` blocks after `return` statements in `if` statements`null` comparisons without type-checking operators`eval()``.bind()``eval()`-like methods`this` in contexts where the value of `this` is `undefined``__iterator__` property`if` statements as the only statement in `else` blocks`new` operators outside of assignments or comparisons`new` operators with the `Function` object`new` operators with the `String`, `Number`, and `Boolean` objects`\8` and `\9` escape sequences in string literals`Object` constructor without an argument`++` and `--``__proto__` property`import``return` statements`javascript:` URLs`undefined``undefined` as an identifier`.call()` and `.apply()``catch` clauses`let` or `const` instead of `var``void` operators`with` statements`const` declarations for variables that are never reassigned after declared`Math.pow` in favor of the `**` operator`parseInt()` and `Number.parseInt()` in favor of binary, octal, and hexadecimal literals`Object.prototype.hasOwnProperty.call()` and prefer use of `Object.hasOwn()``Object.assign` with an object literal as the first argument and prefer the use of object spread instead`RegExp` constructor in favor of regular expression literals`arguments``.apply()``parseInt()``await` expression`u` or `v` flag on regular expressions`yield``import` declarations within modules`var` declarations be placed at the top of their containing scope

These rules care about how the code looks rather than how it executes: Require or disallow Unicode byte order mark (BOM)
🔧 Fix

These rules have been deprecated in accordance with the deprecation policy, and replaced by newer rules:
array-bracket-newline
deprecated
 Replaced by
 ❌
🔧 Fix

array-bracket-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

array-element-newline
deprecated
 Replaced by
 ❌
🔧 Fix

arrow-parens
deprecated
 Replaced by
 ❌
🔧 Fix

arrow-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

block-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

brace-style
deprecated
 Replaced by
 ❌
🔧 Fix

callback-return
deprecated
 Replaced by
 ❌
comma-dangle
deprecated
 Replaced by
 ❌
🔧 Fix

comma-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

comma-style
deprecated
 Replaced by
 ❌
🔧 Fix

computed-property-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

dot-location
deprecated
 Replaced by
 ❌
🔧 Fix

eol-last
deprecated
 Replaced by
 ❌
🔧 Fix

func-call-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

function-call-argument-newline
deprecated
 Replaced by
 ❌
🔧 Fix

function-paren-newline
deprecated
 Replaced by
 ❌
🔧 Fix

generator-star-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

global-require
deprecated
 Replaced by
 ❌
handle-callback-err
deprecated
 Replaced by
 ❌
id-blacklist
deprecated
 Replaced by
 ❌
implicit-arrow-linebreak
deprecated
 Replaced by
 ❌
🔧 Fix

indent
deprecated
 Replaced by
 ❌
🔧 Fix

indent-legacy
deprecated
 Replaced by
 ❌
🔧 Fix

jsx-quotes
deprecated
 Replaced by
 ❌
🔧 Fix

key-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

keyword-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

line-comment-position
deprecated
 Replaced by
 ❌
linebreak-style
deprecated
 Replaced by
 ❌
🔧 Fix

lines-around-comment
deprecated
 Replaced by
 ❌
🔧 Fix

lines-around-directive
deprecated
 Replaced by
 ❌
🔧 Fix

lines-between-class-members
deprecated
 Replaced by
 ❌
🔧 Fix

max-len
deprecated
 Replaced by
 ❌
max-statements-per-line
deprecated
 Replaced by
 ❌
multiline-comment-style
deprecated
 Replaced by
 ❌
🔧 Fix

multiline-ternary
deprecated
 Replaced by
 ❌
🔧 Fix

new-parens
deprecated
 Replaced by
 ❌
🔧 Fix

newline-after-var
deprecated
 Replaced by
 ❌
🔧 Fix

newline-before-return
deprecated
 Replaced by
 ❌
🔧 Fix

newline-per-chained-call
deprecated
 Replaced by
 ❌
🔧 Fix

no-buffer-constructor
deprecated
 Replaced by
 ❌
no-catch-shadow
deprecated
 Replaced by
 ❌
no-confusing-arrow
deprecated
 Replaced by
 ❌
🔧 Fix

no-extra-parens
deprecated
 Replaced by
 ❌
🔧 Fix

no-extra-semi
deprecated
 Replaced by
 ❌
🔧 Fix

no-floating-decimal
deprecated
 Replaced by
 ❌
🔧 Fix

no-mixed-operators
deprecated
 Replaced by
 ❌
no-mixed-requires
deprecated
 Replaced by
 ❌
no-mixed-spaces-and-tabs
deprecated
 Replaced by
 ❌
no-multi-spaces
deprecated
 Replaced by
 ❌
🔧 Fix

no-multiple-empty-lines
deprecated
 Replaced by
 ❌
🔧 Fix

no-native-reassign
deprecated
 Replaced by
 ❌
no-negated-in-lhs
deprecated
 Replaced by
 ❌
no-new-object
deprecated
 Replaced by
 ❌
no-new-require
deprecated
 Replaced by
 ❌
no-new-symbol
deprecated
 Replaced by
 ❌
no-path-concat
deprecated
 Replaced by
 ❌
no-process-env
deprecated
 Replaced by
 ❌
no-process-exit
deprecated
 Replaced by
 ❌
no-restricted-modules
deprecated
 Replaced by
 ❌
no-return-await
deprecated
 ❌
💡 Suggestions

no-spaced-func
deprecated
 Replaced by
 ❌
🔧 Fix

no-sync
deprecated
 Replaced by
 ❌
no-tabs
deprecated
 Replaced by
 ❌
no-trailing-spaces
deprecated
 Replaced by
 ❌
🔧 Fix

no-whitespace-before-property
deprecated
 Replaced by
 ❌
🔧 Fix

nonblock-statement-body-position
deprecated
 Replaced by
 ❌
🔧 Fix

object-curly-newline
deprecated
 Replaced by
 ❌
🔧 Fix

object-curly-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

object-property-newline
deprecated
 Replaced by
 ❌
🔧 Fix

one-var-declaration-per-line
deprecated
 Replaced by
 ❌
🔧 Fix

operator-linebreak
deprecated
 Replaced by
 ❌
🔧 Fix

padded-blocks
deprecated
 Replaced by
 ❌
🔧 Fix

padding-line-between-statements
deprecated
 Replaced by
 ❌
🔧 Fix

prefer-reflect
deprecated
 ❌
quote-props
deprecated
 Replaced by
 ❌
🔧 Fix

quotes
deprecated
 Replaced by
 ❌
🔧 Fix

rest-spread-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

semi
deprecated
 Replaced by
 ❌
🔧 Fix

semi-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

semi-style
deprecated
 Replaced by
 ❌
🔧 Fix

space-before-blocks
deprecated
 Replaced by
 ❌
🔧 Fix

space-before-function-paren
deprecated
 Replaced by
 ❌
🔧 Fix

space-in-parens
deprecated
 Replaced by
 ❌
🔧 Fix

space-infix-ops
deprecated
 Replaced by
 ❌
🔧 Fix

space-unary-ops
deprecated
 Replaced by
 ❌
🔧 Fix

spaced-comment
deprecated
 Replaced by
 ❌
🔧 Fix

switch-colon-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

template-curly-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

template-tag-spacing
deprecated
 Replaced by
 ❌
🔧 Fix

wrap-iife
deprecated
 Replaced by
 ❌
🔧 Fix

wrap-regex
deprecated
 Replaced by
 ❌
🔧 Fix

yield-star-spacing
deprecated
 Replaced by
 ❌
🔧 Fix
`array-bracket-newline`
 in `@stylistic/eslint-plugin` `array-bracket-spacing`
 in `@stylistic/eslint-plugin` `array-element-newline`
 in `@stylistic/eslint-plugin` `arrow-parens`
 in `@stylistic/eslint-plugin` `arrow-spacing`
 in `@stylistic/eslint-plugin` `block-spacing`
 in `@stylistic/eslint-plugin` `brace-style`
 in `@stylistic/eslint-plugin` `callback-return`
 in `eslint-plugin-n` `comma-dangle`
 in `@stylistic/eslint-plugin` `comma-spacing`
 in `@stylistic/eslint-plugin` `comma-style`
 in `@stylistic/eslint-plugin` `computed-property-spacing`
 in `@stylistic/eslint-plugin` `dot-location`
 in `@stylistic/eslint-plugin` `eol-last`
 in `@stylistic/eslint-plugin` `function-call-spacing`
 in `@stylistic/eslint-plugin` `function-call-argument-newline`
 in `@stylistic/eslint-plugin` `function-paren-newline`
 in `@stylistic/eslint-plugin` `generator-star-spacing`
 in `@stylistic/eslint-plugin` `global-require`
 in `eslint-plugin-n` `handle-callback-err`
 in `eslint-plugin-n` `id-denylist`
`implicit-arrow-linebreak`
 in `@stylistic/eslint-plugin` `indent`
 in `@stylistic/eslint-plugin` `indent`
 in `@stylistic/eslint-plugin` `jsx-quotes`
 in `@stylistic/eslint-plugin` `key-spacing`
 in `@stylistic/eslint-plugin` `keyword-spacing`
 in `@stylistic/eslint-plugin` `line-comment-position`
 in `@stylistic/eslint-plugin` `linebreak-style`
 in `@stylistic/eslint-plugin` `lines-around-comment`
 in `@stylistic/eslint-plugin` `padding-line-between-statements`
 in `@stylistic/eslint-plugin` `lines-between-class-members`
 in `@stylistic/eslint-plugin` `max-len`
 in `@stylistic/eslint-plugin` `max-statements-per-line`
 in `@stylistic/eslint-plugin` `multiline-comment-style`
 in `@stylistic/eslint-plugin` `multiline-ternary`
 in `@stylistic/eslint-plugin` `new-parens`
 in `@stylistic/eslint-plugin` `padding-line-between-statements`
 in `@stylistic/eslint-plugin` `padding-line-between-statements`
 in `@stylistic/eslint-plugin` `newline-per-chained-call`
 in `@stylistic/eslint-plugin` `no-deprecated-api`
 in `eslint-plugin-n` `no-shadow`
`no-confusing-arrow`
 in `@stylistic/eslint-plugin` `no-extra-parens`
 in `@stylistic/eslint-plugin` `no-extra-semi`
 in `@stylistic/eslint-plugin` `no-floating-decimal`
 in `@stylistic/eslint-plugin` `no-mixed-operators`
 in `@stylistic/eslint-plugin` `no-mixed-requires`
 in `eslint-plugin-n` `no-mixed-spaces-and-tabs`
 in `@stylistic/eslint-plugin` `no-multi-spaces`
 in `@stylistic/eslint-plugin` `no-multiple-empty-lines`
 in `@stylistic/eslint-plugin` `no-global-assign`
`no-unsafe-negation`
`no-object-constructor`
`no-new-require`
 in `eslint-plugin-n` `no-new-native-nonconstructor`
`no-path-concat`
 in `eslint-plugin-n` `no-process-env`
 in `eslint-plugin-n` `no-process-exit`
 in `eslint-plugin-n` `no-restricted-require`
 in `eslint-plugin-n` `function-call-spacing`
 in `@stylistic/eslint-plugin` `no-sync`
 in `eslint-plugin-n` `no-tabs`
 in `@stylistic/eslint-plugin` `no-trailing-spaces`
 in `@stylistic/eslint-plugin` `no-whitespace-before-property`
 in `@stylistic/eslint-plugin` `nonblock-statement-body-position`
 in `@stylistic/eslint-plugin` `object-curly-newline`
 in `@stylistic/eslint-plugin` `object-curly-spacing`
 in `@stylistic/eslint-plugin` `object-property-newline`
 in `@stylistic/eslint-plugin` `one-var-declaration-per-line`
 in `@stylistic/eslint-plugin` `operator-linebreak`
 in `@stylistic/eslint-plugin` `padded-blocks`
 in `@stylistic/eslint-plugin` `padding-line-between-statements`
 in `@stylistic/eslint-plugin` `quote-props`
 in `@stylistic/eslint-plugin` `quotes`
 in `@stylistic/eslint-plugin` `rest-spread-spacing`
 in `@stylistic/eslint-plugin` `semi`
 in `@stylistic/eslint-plugin` `semi-spacing`
 in `@stylistic/eslint-plugin` `semi-style`
 in `@stylistic/eslint-plugin` `space-before-blocks`
 in `@stylistic/eslint-plugin` `space-before-function-paren`
 in `@stylistic/eslint-plugin` `space-in-parens`
 in `@stylistic/eslint-plugin` `space-infix-ops`
 in `@stylistic/eslint-plugin` `space-unary-ops`
 in `@stylistic/eslint-plugin` `spaced-comment`
 in `@stylistic/eslint-plugin` `switch-colon-spacing`
 in `@stylistic/eslint-plugin` `template-curly-spacing`
 in `@stylistic/eslint-plugin` `template-tag-spacing`
 in `@stylistic/eslint-plugin` `wrap-iife`
 in `@stylistic/eslint-plugin` `wrap-regex`
 in `@stylistic/eslint-plugin` `yield-star-spacing`
 in `@stylistic/eslint-plugin`

These rules from older versions of ESLint (before the deprecation policy existed) have been replaced by newer rules:
generator-star
removed
 Replaced by

global-strict
removed
 Replaced by

no-arrow-condition
removed
 Replaced by

no-comma-dangle
removed
 Replaced by

no-empty-class
removed
 Replaced by

no-empty-label
removed
 Replaced by

no-extra-strict
removed
 Replaced by

no-reserved-keys
removed
 Replaced by

no-space-before-semi
removed
 Replaced by

no-wrap-func
removed
 Replaced by

space-after-function-name
removed
 Replaced by

space-after-keywords
removed
 Replaced by

space-before-function-parentheses
removed
 Replaced by

space-before-keywords
removed
 Replaced by

space-in-brackets
removed
 Replaced by

space-return-throw-case
removed
 Replaced by

space-unary-word-ops
removed
 Replaced by

spaced-line-comment
removed
 Replaced by

valid-jsdoc
removed

require-jsdoc
removed
`generator-star-spacing`
`strict`
`no-confusing-arrow`
or

`no-constant-condition`
`comma-dangle`
`no-empty-character-class`
`no-labels`
`strict`
`quote-props`
`semi-spacing`
`no-extra-parens`
`space-before-function-paren`
`keyword-spacing`
`space-before-function-paren`
`keyword-spacing`
`object-curly-spacing`
or

`array-bracket-spacing`
or

`computed-property-spacing`
`keyword-spacing`
`space-unary-ops`
`spaced-comment`

# array-callback-return

Enforce `return` statements in callbacks of array methods

Some problems reported by this rule are manually fixable by editor suggestions

`Array` has several methods for filtering, mapping, and folding.
If we forget to write `return` statement in a callback of those, it’s probably a mistake. If you don’t want to use a return or don’t need the returned results, consider using .forEach instead.

```
// example: convert ['a', 'b', 'c'] --> {a: 0, b: 1, c: 2}
const indexMap = myArray.reduce(function(memo, item, index) {
 memo[item] = index;
}, {}); // Error: cannot set property 'b' of undefined
```
## Rule Details

This rule enforces usage of `return` statement in callbacks of array’s methods.
Additionally, it may also enforce the `forEach` array method callback to **not** return a value by using the `checkForEach` option.

This rule finds callback functions of the following methods, then checks usage of `return` statement.

- `Array.from`
- `Array.fromAsync`
- `Array.prototype.every`
- `Array.prototype.filter`
- `Array.prototype.find`
- `Array.prototype.findIndex`
- `Array.prototype.findLast`
- `Array.prototype.findLastIndex`
- `Array.prototype.flatMap`
- `Array.prototype.forEach`(optional, based on- `checkForEach`parameter)
- `Array.prototype.map`
- `Array.prototype.reduce`
- `Array.prototype.reduceRight`
- `Array.prototype.some`
- `Array.prototype.sort`
- `Array.prototype.toSorted`
- And above of typed arrays if applicable.

Examples of **incorrect** code for this rule:

```
/*eslint array-callback-return: "error"*/
const indexMap = myArray.reduce((memo, item, index) {
 memo[item] = index;
}, {});
const foo = Array.from(nodes, (node) {
 if (node.tagName === "DIV") {
 return true;
 }
});
const bar = foo.filter(function(x) {
 if (x) {
 return true;
 } else {

 }
});
```
Examples of **correct** code for this rule:

```
/*eslint array-callback-return: "error"*/
const indexMap = myArray.reduce(function(memo, item, index) {
 memo[item] = index;
 return memo;
}, {});
const foo = Array.from(nodes, function(node) {
 if (node.tagName === "DIV") {
 return true;
 }
 return false;
});
const bar = foo.map(node => node.getAttribute("id"));
```
## Options

This rule accepts a configuration object with three options:

- `"allowImplicit": false`(default) When set to- `true`, allows callbacks of methods that require a return value to implicitly return- `undefined`with a- `return`statement containing no expression.
- `"checkForEach": false`(default) When set to- `true`, rule will also report- `forEach`callbacks that return a value.
- `"allowVoid": false`(default) When set to- `true`, allows- `void`in- `forEach`callbacks, so rule will not report the return value with a- `void`operator.

**Note:** `{ "allowVoid": true }` works only if `checkForEach` option is set to `true`.

### allowImplicit

Examples of **correct** code for the `{ "allowImplicit": true }` option:

```
/*eslint array-callback-return: ["error", { allowImplicit: true }]*/
const undefAllTheThings = myArray.map(function(item) {
 return;
});
```
### checkForEach

Examples of **incorrect** code for the `{ "checkForEach": true }` option:

```
/*eslint array-callback-return: ["error", { checkForEach: true }]*/
myArray.forEach(function(item) {

});
myArray.forEach(function(item) {
 if (item < 0) {

 }
 handleItem(item);
});
myArray.forEach(function(item) {
 if (item < 0) {

 }
 handleItem(item);
});
myArray.forEach(item handleItem(item));
myArray.forEach(item void handleItem(item));
myArray.forEach(item => {

});
myArray.forEach(item => {

});
```
Examples of **correct** code for the `{ "checkForEach": true }` option:

```
/*eslint array-callback-return: ["error", { checkForEach: true }]*/
myArray.forEach(function(item) {
 handleItem(item)
});
myArray.forEach(function(item) {
 if (item < 0) {
 return;
 }
 handleItem(item);
});
myArray.forEach(function(item) {
 handleItem(item);
 return;
});
myArray.forEach(item => {
 handleItem(item);
});
```
### allowVoid

Examples of **correct** code for the `{ "allowVoid": true }` option:

```
/*eslint array-callback-return: ["error", { checkForEach: true, allowVoid: true }]*/
myArray.forEach(item => void handleItem(item));
myArray.forEach(item => {
 return void handleItem(item);
});
myArray.forEach(item => {
 if (item < 0) {
 return void x;
 }
 handleItem(item);
});
```
## Known Limitations

This rule checks callback functions of methods with the given names, *even if* the object which has the method is *not* an array.

## When Not To Use It

If you don’t want to warn about usage of `return` statement in callbacks of array’s methods, then it’s safe to disable this rule.

## Version

This rule was introduced in ESLint v2.0.0-alpha-1.

# constructor-super

Require `super()` calls in constructors

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Constructors of derived classes must call `super()`.
Constructors of non derived classes must not call `super()`.
If this is not observed, the JavaScript engine will raise a runtime error.

This rule checks whether or not there is a valid `super()` call.

## Rule Details

This rule is aimed to flag invalid/missing `super()` calls.

This is a syntax error because there is no `extends` clause in the class:

```
class A {
 constructor() {
 super();
 }
}
```
Examples of **incorrect** code for this rule:

```
/*eslint constructor-super: "error"*/
class A extends B {
 // Would throw a ReferenceError.
}
// Classes which inherits from a non constructor are always problems.
class C extends null {
 constructor() {
 ; // Would throw a TypeError.
 }
}
class D extends null {
 // Would throw a ReferenceError.
}
```
Examples of **correct** code for this rule:

```
/*eslint constructor-super: "error"*/
class A {
 constructor() { }
}
class B extends C {
 constructor() {
 super();
 }
}
```
## Options

This rule has no options.

## When Not To Use It

If you don’t want to be notified about invalid/missing `super()` callings in constructors, you can safely disable this rule.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

## Version

This rule was introduced in ESLint v0.24.0.

# for-direction

Enforce `for` loop update clause moving the counter in the right direction

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

A `for` loop with a stop condition that can never be reached, such as one with a counter that moves in the wrong direction, will run infinitely. While there are occasions when an infinite loop is intended, the convention is to construct such loops as `while` loops. More typically, an infinite `for` loop is a bug.

## Rule Details

This rule forbids `for` loops where the counter variable changes in such a way that the stop condition will never be met. For example, if the counter variable is increasing (i.e. `i++`) and the stop condition tests that the counter is greater than zero (`i >= 0`) then the loop will never exit.

Note:This rule only checks thedirectionof the counter relative to the stop condition. It does not check the actual values to determine whether the loop will execute at least once. For example, a loop like`for (let i = 0; i < 0; i++) {}`is considered valid by this rule because the counter direction (`i++`) matches the condition (`i <`), even though the condition is false from the start (dead code).

Examples of **incorrect** code for this rule:

```
/*eslint for-direction: "error"*/
 {
}
 {
}
 {
 // counter i is on the left with >, so i++ (increasing) is the wrong direction
}
 {
 // counter i is on the left with >, so i++ (increasing) is the wrong direction
}
 {
 // counter i is on the right with <, so i++ (increasing) is the wrong direction
}
 {
}
const n = -2;
 {
}
```
Examples of **correct** code for this rule:

```
/*eslint for-direction: "error"*/
for (let i = 0; i < 10; i++) {
}
for (let i = 0; 10 > i; i++) { // with counter "i" on the right
}
for (let i = 10; i >= 0; i += this.step) { // direction unknown
}
for (let i = MIN; i <= MAX; i -= 0) { // not increasing or decreasing
}
for (let i = 0; i < 0; i++) {
 // counter i is on the left with <, so i++ (increasing) is the correct direction
 // (loop never executes, but direction is consistent)
}
for (let i = 0; 0 > i; i++) {
 // counter i is on the right with >, so i++ (increasing) is the correct direction
 // (loop never executes, but direction is consistent)
}
```
## Options

This rule has no options.

## Version

This rule was introduced in ESLint v4.0.0-beta.0.

# getter-return

Enforce `return` statements in getters

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

The get syntax binds an object property to a function that will be called when that property is looked up. It was first introduced in ECMAScript 5:

```
const p = {
 get name(){
 return "nicholas";
 }
};
Object.defineProperty(p, "age", {
 get: function (){
 return 17;
 }
});
```
Note that every `getter` is expected to return a value.

## Rule Details

This rule enforces that a return statement is present in property getters.

Examples of **incorrect** code for this rule:

```
/*eslint getter-return: "error"*/
const p = {
 (){
 // no returns.
 }
};
Object.defineProperty(p, "age", {
 (){
 // no returns.
 }
});
class P{
 (){
 // no returns.
 }
}
```
Examples of **correct** code for this rule:

```
/*eslint getter-return: "error"*/
const p = {
 get name(){
 return "nicholas";
 }
};
Object.defineProperty(p, "age", {
 get: function (){
 return 18;
 }
});
class P{
 get name(){
 return "nicholas";
 }
}
```
## Options

This rule has an object option:

- `"allowImplicit": false`(default) disallows implicitly returning- `undefined`with a- `return`statement.

Examples of **correct** code for the `{ "allowImplicit": true }` option:

```
/*eslint getter-return: ["error", { allowImplicit: true }]*/
const p = {
 get name(){
 return; // return undefined implicitly.
 }
};
```
## When Not To Use It

If your project will not be using ES5 property getters you do not need this rule.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

## Version

This rule was introduced in ESLint v4.2.0.

# no-async-promise-executor

Disallow using an async function as a Promise executor

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

The `new Promise` constructor accepts an *executor* function as an argument, which has `resolve` and `reject` parameters that can be used to control the state of the created Promise. For example:

```
const result = new Promise(function executor(resolve, reject) {
 readFile('foo.txt', function(err, result) {
 if (err) {
 reject(err);
 } else {
 resolve(result);
 }
 });
});
```
The executor function can also be an `async function`. However, this is usually a mistake, for a few reasons:

- If an async executor function throws an error, the error will be lost and won’t cause the newly-constructed `Promise`to reject. This could make it difficult to debug and handle some errors.
- If a Promise executor function is using `await`, this is usually a sign that it is not actually necessary to use the`new Promise`constructor, or the scope of the`new Promise`constructor can be reduced.

## Rule Details

This rule aims to disallow async Promise executor functions.

Examples of **incorrect** code for this rule:

```
/*eslint no-async-promise-executor: "error"*/
const foo = new Promise( (resolve, reject) => {
 readFile('foo.txt', function(err, result) {
 if (err) {
 reject(err);
 } else {
 resolve(result);
 }
 });
});
const result = new Promise( (resolve, reject) => {
 resolve(await foo);
});
```
Examples of **correct** code for this rule:

```
/*eslint no-async-promise-executor: "error"*/
const foo = new Promise((resolve, reject) => {
 readFile('foo.txt', function(err, result) {
 if (err) {
 reject(err);
 } else {
 resolve(result);
 }
 });
});
const result = Promise.resolve(foo);
```
## Options

This rule has no options.

## When Not To Use It

If your codebase doesn’t support `async function` syntax, there’s no need to enable this rule.

## Version

This rule was introduced in ESLint v5.3.0.

# no-await-in-loop

Disallow `await` inside of loops

Performing an operation on each element of an iterable is a common task. However, performing an
`await` as part of each operation may indicate that the program is not taking full advantage of
the parallelization benefits of `async`/`await`.

Often, the code can be refactored to create all the promises at once, then get access to the
results using `Promise.all()` (or one of the other promise concurrency methods). Otherwise, each successive operation will not start until the
previous one has completed.

Concretely, the following function could be refactored as shown:

```
async function foo(things) {
 const results = [];
 for (const thing of things) {
 // Bad: each loop iteration is delayed until the entire asynchronous operation completes
 results.push(await doAsyncWork(thing));
 }
 return results;
}
```
```
async function foo(things) {
 const promises = [];
 for (const thing of things) {
 // Good: all asynchronous operations are immediately started.
 promises.push(doAsyncWork(thing));
 }
 // Now that all the asynchronous operations are running, here we wait until they all complete.
 const results = await Promise.all(promises);
 return results;
}
```
This can be beneficial for subtle error-handling reasons as well. Given an array of promises that might reject, sequential awaiting puts the program at risk of unhandled promise rejections. The exact behavior of unhandled rejections depends on the environment running your code, but they are generally considered harmful regardless. In Node.js, for example, unhandled rejections cause a program to terminate unless configured otherwise.

```
async function foo() {
 const arrayOfPromises = somethingThatCreatesAnArrayOfPromises();
 for (const promise of arrayOfPromises) {
 // Bad: if any of the promises reject, an exception is thrown, and
 // subsequent loop iterations will not run. Therefore, rejections later
 // in the array will become unhandled rejections that cannot be caught
 // by a caller.
 const value = await promise;
 console.log(value);
 }
}
```
```
async function foo() {
 const arrayOfPromises = somethingThatCreatesAnArrayOfPromises();
 // Good: Any rejections will cause a single exception to be thrown here,
 // which may be caught and handled by the caller.
 const arrayOfValues = await Promise.all(arrayOfPromises);
 for (const value of arrayOfValues) {
 console.log(value);
 }
}
```
## Rule Details

This rule disallows the use of `await` within loop bodies.

Examples of **correct** code for this rule:

```
/*eslint no-await-in-loop: "error"*/
async function foo(things) {
 const promises = [];
 for (const thing of things) {
 // Good: all asynchronous operations are immediately started.
 promises.push(doAsyncWork(thing));
 }
 // Now that all the asynchronous operations are running, here we wait until they all complete.
 const results = await Promise.all(promises);
 return results;
}
```
Examples of **incorrect** code for this rule:

```
/*eslint no-await-in-loop: "error"*/
async function foo(things) {
 const results = [];
 for (const thing of things) {
 // Bad: each loop iteration is delayed until the entire asynchronous operation completes
 results.push();
 }
 return results;
}
async function bar(things) {
 for (const thing of things) {

 }
}
```
## Options

This rule has no options.

## When Not To Use It

In many cases the iterations of a loop are not actually independent of each other, and awaiting in the loop is correct. As a few examples:

-
Any code that should run serially, such as implementing a countdown, or processing items sequentially (when each item is already processed in parallel), as in this example: `async function printCountdown() { for (let i = 0; i < 10; i++) { await new Promise(resolve => setTimeout(resolve, 1000)); // sleep 1 second console.log(i); } }`
-
Using standard browser/OS APIs that are inherently serial (as controlled by the underlying operating system), such as `File`or directory contents reading, as in this example:`async function writeNumbersToFile() { for (let i = 0; i < 100; i++) { // Doing this in parallel would interleave (garble) the resulting file contents. await fileWriteStream.write(i + "\n"); } }`
-
Any code where the creation of the Promise allocates a bounded resource (RAM, file descriptors, network bandwidth), as in this example: `async function streamingProcess() { for (let i = 0; i < 10; i++) { await ramIntensiveAction(data[i]); } }``async function checkFiles() { for (let i = 0; i < 10000; i++) { // Here some concurrency may be desirable to reduce I/O bottlenecks, // but not unbounded concurrency as that will fail, // running out of file descriptors. await openFileAndThrowIfContentsMeetCondition(filenames[i]); } }`
-
The output of one iteration might be used as the input to another. `async function loopIterationsDependOnEachOther() { let previousResult = null; for (let i = 0; i < 10; i++) { const result = await doSomething(i, previousResult); if (someCondition(result, previousResult)) { break; } else { previousResult = result; } } }`The previous examples were all similar to this, but they relied on a *side effect*of one iteration being needed for the next, while this example has the dependency explicit via a result variable.
-
Loops may be used to retry asynchronous operations that were unsuccessful. `async function retryUpTo10Times() { for (let i = 0; i < 10; i++) { const wasSuccessful = await tryToDoSomething(); if (wasSuccessful) return 'succeeded!'; // wait to try again. await new Promise(resolve => setTimeout(resolve, 1000)); } return 'failed!'; }`
-
Loops may be used to prevent your code from sending an excessive amount of requests in parallel. `async function makeUpdatesToRateLimitedApi(thingsToUpdate) { // we'll exceed our rate limit if we make all the network calls in parallel. for (const thing of thingsToUpdate) { await updateThingWithRateLimitedApi(thing); } }`

In such cases it makes sense to use `await` within a
loop and it is recommended to disable the rule via a standard ESLint disable comment.

## Version

This rule was introduced in ESLint v3.12.0.

# no-class-assign

Disallow reassigning class members

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

`ClassDeclaration` creates a variable, and we can modify the variable.

```
class A { }
A = 0;
```
But the modification is a mistake in most cases.

## Rule Details

This rule is aimed to flag modifying variables of class declarations.

Examples of **incorrect** code for this rule:

```
/*eslint no-class-assign: "error"*/
class A { }
 = 0;
```
```
/*eslint no-class-assign: "error"*/
 = 0;
class A { }
```
```
/*eslint no-class-assign: "error"*/
class A {
 b() {
 = 0;
 }
}
```
```
/*eslint no-class-assign: "error"*/
let A = class A {
 b() {
 = 0;
 // `let A` is shadowed by the class name.
 }
}
```
Examples of **correct** code for this rule:

```
/*eslint no-class-assign: "error"*/
let A = class A { }
A = 0; // A is a variable.
```
```
/*eslint no-class-assign: "error"*/
let A = class {
 b() {
 A = 0; // A is a variable.
 }
}
```
```
/*eslint no-class-assign: 2*/
class A {
 b(A) {
 A = 0; // A is a parameter.
 }
}
```
## Options

This rule has no options.

## When Not To Use It

If you don’t want to be notified about modifying variables of class declarations, you can safely disable this rule.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

## Version

This rule was introduced in ESLint v1.0.0-rc-1.

# no-compare-neg-zero

Disallow comparing against `-0`

 ✅ Recommended

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

## Rule Details

The rule should warn against code that tries to compare against `-0`, since that will not work as intended. That is, code like `x === -0` will pass for both `+0` and `-0`. The author probably intended `Object.is(x, -0)`.

Examples of **incorrect** code for this rule:

 Open in Playground

```
/* eslint no-compare-neg-zero: "error" */
if () {
 // doSomething()...
}
```
Examples of **correct** code for this rule:

 Open in Playground

```
/* eslint no-compare-neg-zero: "error" */
if (x === 0) {
 // doSomething()...
}
```

 Open in Playground

```
/* eslint no-compare-neg-zero: "error" */
if (Object.is(x, -0)) {
 // doSomething()...
}
```
## Options

This rule has no options.

## Version

This rule was introduced in ESLint v3.17.0.

# no-cond-assign

Disallow assignment operators in conditional expressions

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

In conditional statements, it is very easy to mistype a comparison operator (such as `==`) as an assignment operator (such as `=`). For example:

```
// Check the user's job title
if (user.jobTitle = "manager") {
 // user.jobTitle is now incorrect
}
```
There are valid reasons to use assignment operators in conditional statements. However, it can be difficult to tell whether a specific assignment was intentional.

## Rule Details

This rule disallows ambiguous assignment operators in test conditions of `if`, `for`, `while`, and `do...while` statements.

## Options

This rule has a string option:

- `"except-parens"`(default) allows assignments in test conditions- *only if*they are enclosed in parentheses (for example, to allow reassigning a variable in the test of a- `while`or- `do...while`loop).
- `"always"`disallows all assignments in test conditions.

### except-parens

Examples of **incorrect** code for this rule with the default `"except-parens"` option:

```
/*eslint no-cond-assign: "error"*/
// Unintentional assignment
let x;
if () {
 const b = 1;
}
// Practical example that is similar to an error
const setHeight = function (someNode) {
 do {
 someNode.height = "100px";
 } while ();
}
```
Examples of **correct** code for this rule with the default `"except-parens"` option:

```
/*eslint no-cond-assign: "error"*/
// Assignment replaced by comparison
let x;
if (x === 0) {
 const b = 1;
}
// Practical example that wraps the assignment in parentheses
const setHeight = function (someNode) {
 do {
 someNode.height = "100px";
 } while ((someNode = someNode.parentNode));
}
// Practical example that wraps the assignment and tests for 'null'
const set_height = function (someNode) {
 do {
 someNode.height = "100px";
 } while ((someNode = someNode.parentNode) !== null);
}
```
### always

Examples of **incorrect** code for this rule with the `"always"` option:

```
/*eslint no-cond-assign: ["error", "always"]*/
// Unintentional assignment
let x;
if () {
 const b = 1;
}
// Practical example that is similar to an error
const setHeight = function (someNode) {
 do {
 someNode.height = "100px";
 } while ();
}
// Practical example that wraps the assignment in parentheses
const set_height = function (someNode) {
 do {
 someNode.height = "100px";
 } while (());
}
// Practical example that wraps the assignment and tests for 'null'
const heightSetter = function (someNode) {
 do {
 someNode.height = "100px";
 } while (() !== null);
}
```
Examples of **correct** code for this rule with the `"always"` option:

```
/*eslint no-cond-assign: ["error", "always"]*/
// Assignment replaced by comparison
let x;
if (x === 0) {
 const b = 1;
}
```
## Version

This rule was introduced in ESLint v0.0.9.

# no-const-assign

Disallow reassigning `const`, `using`, and `await using` variables

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Constant bindings cannot be modified. An attempt to modify a constant binding will raise a runtime error.

## Rule Details

This rule is aimed to flag modifying variables that are declared using `const`, `using`, or `await using` keywords.

Examples of **incorrect** code for this rule:

```
/*eslint no-const-assign: "error"*/
const a = 0;
 = 1;
```
```
/*eslint no-const-assign: "error"*/
const a = 0;
 += 1;
```
```
/*eslint no-const-assign: "error"*/
const a = 0;
++;
```
```
/*eslint no-const-assign: "error"*/
if (foo) {
	using a = getSomething();
 = somethingElse;
}
if (bar) {
	await using a = getSomething();
 = somethingElse;
}
```
Examples of **correct** code for this rule:

```
/*eslint no-const-assign: "error"*/
const a = 0;
console.log(a);
```
```
/*eslint no-const-assign: "error"*/
if (foo) {
	using a = getSomething();
	a.execute();
}
if (bar) {
	await using a = getSomething();
	a.execute();
}
```
```
/*eslint no-const-assign: "error"*/
for (const a in [1, 2, 3]) { // `a` is re-defined (not modified) on each loop step.
 console.log(a);
}
```
```
/*eslint no-const-assign: "error"*/
for (const a of [1, 2, 3]) { // `a` is re-defined (not modified) on each loop step.
 console.log(a);
}
```
## Options

This rule has no options.

## When Not To Use It

If you don’t want to be notified about modifying variables that are declared using `const`, `using`, and `await using` keywords, you can safely disable this rule.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

## Version

This rule was introduced in ESLint v1.0.0-rc-1.

# no-constant-binary-expression

Disallow expressions where the operation doesn't affect the value

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Comparisons which will always evaluate to true or false and logical expressions (`||`, `&&`, `??`) which either always short-circuit or never short-circuit are both likely indications of programmer error.

These errors are especially common in complex expressions where operator precedence is easy to misjudge. For example:

```
// One might think this would evaluate as `a + (b ?? c)`:
const x = a + b ?? c;
// But it actually evaluates as `(a + b) ?? c`. Since `a + b` can never be null,
// the `?? c` has no effect.
```
Additionally, this rule detects comparisons to newly constructed objects/arrays/functions/etc. In JavaScript, where objects are compared by reference, a newly constructed object can *never* `===` any other value. This can be surprising for programmers coming from languages where objects are compared by value.

```
// Programmers coming from a language where objects are compared by value might expect this to work:
const isEmpty = x === [];
// However, this will always result in `isEmpty` being `false`.
```
## Rule Details

This rule identifies `==` and `===` comparisons which, based on the semantics of the JavaScript language, will always evaluate to `true` or `false`.

It also identifies `||`, `&&` and `??` logical expressions which will either always or never short-circuit. Additionally, when configured with the `checkRelationalComparisons` option, it can identify constant relational comparisons (such as `<` or `>=`).

Examples of **incorrect** code for this rule:

```
/*eslint no-constant-binary-expression: "error"*/
const value1 = == null;
const value2 = condition ? x : || DEFAULT;
const value3 = == null;
const value4 = === true;
const objIsEmpty = someObj === ;
const arrIsEmpty = someArr === ;
const shortCircuit1 = && condition2;
const shortCircuit2 = || condition2;
const shortCircuit3 = ?? condition2;
```
Examples of **correct** code for this rule:

```
/*eslint no-constant-binary-expression: "error"*/
const value1 = x == null;
const value2 = (condition ? x : {}) || DEFAULT;
const value3 = !(foo == null);
const value4 = Boolean(foo) === true;
const objIsEmpty = Object.keys(someObj).length === 0;
const arrIsEmpty = someArr.length === 0;
```
## Options

This rule has an object option:

- `"checkRelationalComparisons": true`(default- `false`) checks for relational comparisons (- `<`,- `<=`,- `>`,- `>=`) where both sides are literal values.

### checkRelationalComparisons

Examples of **incorrect** code for this rule with the `{ "checkRelationalComparisons": true }` option:

```
/*eslint no-constant-binary-expression: ["error", { "checkRelationalComparisons": true } ]*/
const value1 = ;
const value2 = ;
const value3 = ;
const hasStreak = profile.streak ?? ;
```
Examples of **correct** code for this rule with the `{ "checkRelationalComparisons": true }` option:

```
/*eslint no-constant-binary-expression: ["error", { "checkRelationalComparisons": true } ]*/
const value1 = 1 < x;
const value2 = a >= b;
const hasStreak = (profile.streak ?? 0) > 1;
```
## Related Rules

## Version

This rule was introduced in ESLint v8.14.0.

# no-constant-condition

Disallow constant expressions in conditions

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

A constant expression (for example, a literal) as a test condition might be a typo or development trigger for a specific behavior. For example, the following code looks as if it is not ready for production.

```
if (false) {
 doSomethingUnfinished();
}
```
## Rule Details

This rule disallows constant expressions in the test condition of:

- `if`,- `for`,- `while`, or- `do...while`statement
- `?:`ternary expression

Examples of **incorrect** code for this rule:

```
/*eslint no-constant-condition: "error"*/
if () {
 doSomethingUnfinished();
}
if () {
 doSomethingUnfinished();
}
if () {
 doSomethingNever();
}
if () {
 doSomethingAlways();
}
if () {
 doSomethingAlways();
}
if () {
 doSomethingAlways();
}
if () {
 doSomethingUnfinished();
}
if () {
 doSomethingAlways();
}
for (;;) {
 doSomethingForever();
}
while () {
 doSomethingForever();
}
do {
 doSomethingForever();
} while ();
const result = ? a : b;
if(){
 output(input);
}
```
Examples of **correct** code for this rule:

```
/*eslint no-constant-condition: "error"*/
if (x === 0) {
 doSomething();
}
for (;;) {
 doSomethingForever();
}
while (typeof x === "undefined") {
 doSomething();
}
do {
 doSomething();
} while (x);
const result = x !== 0 ? a : b;
if(input === "hello" || input === "bye"){
 output(input);
}
```
## Options

### checkLoops

This is a string option having following values:

- `"all"`- Disallow constant expressions in all loops.
- `"allExceptWhileTrue"`(default) - Disallow constant expressions in all loops except- `while`loops with expression- `true`.
- `"none"`- Allow constant expressions in loops.

Or instead you can set the `checkLoops` value to booleans where `true` is same as `"all"` and `false` is same as `"none"`.

Examples of **incorrect** code for when `checkLoops` is `"all"` or `true`:

```
/*eslint no-constant-condition: ["error", { "checkLoops": "all" }]*/
while () {
 doSomething();
};
for (;;) {
 doSomething();
};
```
```
/*eslint no-constant-condition: ["error", { "checkLoops": true }]*/
while () {
 doSomething();
};
do {
 doSomething();
} while ()
```
Examples of **correct** code for when `checkLoops` is `"all"` or `true`:

```
/*eslint no-constant-condition: ["error", { "checkLoops": "all" }]*/
while (a === b) {
 doSomething();
};
```
```
/*eslint no-constant-condition: ["error", { "checkLoops": true }]*/
for (let x = 0; x <= 10; x++) {
 doSomething();
};
```
Example of **correct** code for when `checkLoops` is `"allExceptWhileTrue"`:

```
/*eslint no-constant-condition: "error"*/
while (true) {
 doSomething();
};
```
Examples of **correct** code for when `checkLoops` is `"none"` or `false`:

```
/*eslint no-constant-condition: ["error", { "checkLoops": "none" }]*/
while (true) {
 doSomething();
 if (condition()) {
 break;
 }
};
do {
 doSomething();
 if (condition()) {
 break;
 }
} while (true)
```
```
/*eslint no-constant-condition: ["error", { "checkLoops": false }]*/
while (true) {
 doSomething();
 if (condition()) {
 break;
 }
};
for (;true;) {
 doSomething();
 if (condition()) {
 break;
 }
};
```
## Related Rules

## Version

This rule was introduced in ESLint v0.4.1.

# no-constructor-return

Disallow returning value from constructor

In JavaScript, returning a value in the constructor of a class may be a mistake. Forbidding this pattern prevents mistakes resulting from unfamiliarity with the language or a copy-paste error.

## Rule Details

This rule disallows return statements in the constructor of a class. Note that returning nothing is allowed.

Examples of **incorrect** code for this rule:

 Open in Playground

```
/*eslint no-constructor-return: "error"*/
class A {
 constructor(a) {
 this.a = a;

 }
}
class B {
 constructor(f) {
 if (!f) {

 }
 }
}
```
Examples of **correct** code for this rule:

 Open in Playground

```
/*eslint no-constructor-return: "error"*/
class C {
 constructor(c) {
 this.c = c;
 }
}
class D {
 constructor(f) {
 if (!f) {
 return; // Flow control.
 }
 f();
 }
}
class E {
 constructor() {
 return;
 }
}
```
## Options

This rule has no options.

## Version

This rule was introduced in ESLint v6.7.0.

# no-control-regex

Disallow control characters in regular expressions

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Control characters are special, invisible characters in the ASCII range 0-31. These characters are rarely used in JavaScript strings so a regular expression containing elements that explicitly match these characters is most likely a mistake.

## Rule Details

This rule disallows control characters and some escape sequences that match control characters in regular expressions.

The following elements of regular expression patterns are considered possible errors in typing and are therefore disallowed by this rule:

- Hexadecimal character escapes from `\x00`to`\x1F`.
- Unicode character escapes from `\u0000`to`\u001F`.
- Unicode code point escapes from `\u{0}`to`\u{1F}`.
- Unescaped raw characters from U+0000 to U+001F.

Control escapes such as `\t` and `\n` are allowed by this rule.

Examples of **incorrect** code for this rule:

```
/*eslint no-control-regex: "error"*/
const pattern1 = ;
const pattern2 = ;
const pattern3 = ;
const pattern4 = ;
const pattern5 = ;
const pattern6 = new RegExp(); // raw U+000C character in the pattern
const pattern7 = new RegExp(); // \x0C pattern
```
Examples of **correct** code for this rule:

```
/*eslint no-control-regex: "error"*/
const pattern1 = /\x20/;
const pattern2 = /\u0020/;
const pattern3 = /\u{20}/u;
const pattern4 = /\t/;
const pattern5 = /\n/;
const pattern6 = new RegExp("\x20");
const pattern7 = new RegExp("\\t");
const pattern8 = new RegExp("\\n");
```
## Options

This rule has no options.

## Known Limitations

When checking `RegExp` constructor calls, this rule examines evaluated regular expression patterns. Therefore, although this rule intends to allow syntax such as `\t`, it doesn’t allow `new RegExp("\t")` since the evaluated pattern (string value of `"\t"`) contains a raw control character (the TAB character).

```
/*eslint no-control-regex: "error"*/
new RegExp("\t"); // disallowed since the pattern is: <TAB>
new RegExp("\\t"); // allowed since the pattern is: \t
```
There is no difference in behavior between `new RegExp("\t")` and `new RegExp("\\t")`, and the intention to match the TAB character is clear in both cases. They are equally valid for the purpose of this rule, but it only allows `new RegExp("\\t")`.

## When Not To Use It

If you need to use control character pattern matching, then you should turn this rule off.

## Related Rules

## Version

This rule was introduced in ESLint v0.1.0.

# no-debugger

Disallow the use of `debugger`

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

The `debugger` statement is used to tell the executing JavaScript environment to stop execution and start up a debugger at the current point in the code. This has fallen out of favor as a good practice with the advent of modern debugging and development tools. Production code should definitely not contain `debugger`, as it will cause the browser to stop executing code and open an appropriate debugger.

## Rule Details

This rule disallows `debugger` statements.

Example of **incorrect** code for this rule:

```
/*eslint no-debugger: "error"*/
function isTruthy(x) {

 return Boolean(x);
}
```
Example of **correct** code for this rule:

```
/*eslint no-debugger: "error"*/
function isTruthy(x) {
 return Boolean(x); // set a breakpoint at this line
}
```
## Options

This rule has no options.

## When Not To Use It

If your code is still very much in development and don’t want to worry about stripping `debugger` statements, then turn this rule off. You’ll generally want to turn it back on when testing code prior to deployment.

## Related Rules

## Version

This rule was introduced in ESLint v0.0.2.

# no-dupe-args

Disallow duplicate arguments in `function` definitions

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

If more than one parameter has the same name in a function definition, the last occurrence “shadows” the preceding occurrences. A duplicated name might be a typing error.

## Rule Details

This rule disallows duplicate parameter names in function declarations or expressions. It does not apply to arrow functions or class methods, because the parser reports the error.

If ESLint parses code in strict mode, the parser (instead of this rule) reports the error.

Examples of **incorrect** code for this rule:

```
/*eslint no-dupe-args: "error"*/
function foo {
 console.log("value of the second a:", a);
}
const bar = function {
 console.log("value of the second a:", a);
};
```
Examples of **correct** code for this rule:

```
/*eslint no-dupe-args: "error"*/
function foo(a, b, c) {
 console.log(a, b, c);
}
const bar = function (a, b, c) {
 console.log(a, b, c);
};
```
## Options

This rule has no options.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

## Version

This rule was introduced in ESLint v0.16.0.

# no-dupe-class-members

Disallow duplicate class members

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

If there are declarations of the same name in class members, the last declaration overwrites other declarations silently. It can cause unexpected behaviors.

```
class Foo {
 bar() { console.log("hello"); }
 bar() { console.log("goodbye"); }
}
const foo = new Foo();
foo.bar(); // goodbye
```
## Rule Details

This rule is aimed to flag the use of duplicate names in class members.

Examples of **incorrect** code for this rule:

```
/*eslint no-dupe-class-members: "error"*/
class A {
 bar() { }
 () { }
}
class B {
 bar() { }
 get () { }
}
class C {
 bar;
 ;
}
class D {
 bar;
 () { }
}
class E {
 static bar() { }
 static () { }
}
```
Examples of **correct** code for this rule:

```
/*eslint no-dupe-class-members: "error"*/
class A {
 bar() { }
 qux() { }
}
class B {
 get bar() { }
 set bar(value) { }
}
class C {
 bar;
 qux;
}
class D {
 bar;
 qux() { }
}
class E {
 static bar() { }
 bar() { }
}
```
This rule additionally supports TypeScript type syntax. It has support for TypeScript’s method overload definitions.

Examples of **correct** TypeScript code for this rule:

```
/* eslint no-dupe-class-members: "error" */
class A {
	foo(value: string): void;
	foo(value: number): void;
	foo(value: string | number) {} // ✅ This is the actual implementation.
}
```
## Options

This rule has no options.

## When Not To Use It

This rule should not be used in ES3/5 environments.

In ES2015 (ES6) or later, if you don’t want to be notified about duplicate names in class members, you can safely disable this rule.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

## Version

This rule was introduced in ESLint v1.2.0.

# no-dupe-else-if

Disallow duplicate conditions in if-else-if chains

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

`if-else-if` chains are commonly used when there is a need to execute only one branch (or at most one branch) out of several possible branches, based on certain conditions.

```
if (a) {
 foo();
} else if (b) {
 bar();
} else if (c) {
 baz();
}
```
Two identical test conditions in the same chain are almost always a mistake in the code. Unless there are side effects in the expressions, a duplicate will evaluate to the same `true` or `false` value as the identical expression earlier in the chain, meaning that its branch can never execute.

```
if (a) {
 foo();
} else if (b) {
 bar();
} else if (b) {
 baz();
}
```
In the above example, `baz()` can never execute. Obviously, `baz()` could be executed only when `b` evaluates to `true`, but in that case `bar()` would be executed instead, since it’s earlier in the chain.

## Rule Details

This rule disallows duplicate conditions in the same `if-else-if` chain.

Examples of **incorrect** code for this rule:

```
/*eslint no-dupe-else-if: "error"*/
if (isSomething(x)) {
 foo();
} else if () {
 bar();
}
if (a) {
 foo();
} else if (b) {
 bar();
} else if (c && d) {
 baz();
} else if () {
 quux();
} else {
 quuux();
}
if (n === 1) {
 foo();
} else if (n === 2) {
 bar();
} else if (n === 3) {
 baz();
} else if () {
 quux();
} else if (n === 5) {
 quuux();
}
```
Examples of **correct** code for this rule:

```
/*eslint no-dupe-else-if: "error"*/
if (isSomething(x)) {
 foo();
} else if (isSomethingElse(x)) {
 bar();
}
if (a) {
 foo();
} else if (b) {
 bar();
} else if (c && d) {
 baz();
} else if (c && e) {
 quux();
} else {
 quuux();
}
if (n === 1) {
 foo();
} else if (n === 2) {
 bar();
} else if (n === 3) {
 baz();
} else if (n === 4) {
 quux();
} else if (n === 5) {
 quuux();
}
```
This rule can also detect some cases where the conditions are not identical, but the branch can never execute due to the logic of `||` and `&&` operators.

Examples of additional **incorrect** code for this rule:

```
/*eslint no-dupe-else-if: "error"*/
if (a || b) {
 foo();
} else if () {
 bar();
}
if (a) {
 foo();
} else if (b) {
 bar();
} else if () {
 baz();
}
if (a) {
 foo();
} else if () {
 bar();
}
if (a && b) {
 foo();
} else if () {
 bar();
}
if (a || b) {
 foo();
} else if () {
 bar();
}
if (a) {
 foo();
} else if (b && c) {
 bar();
} else if () {
 baz();
}
```
Please note that this rule does not compare conditions from the chain with conditions inside statements, and will not warn in the cases such as follows:

```
if (a) {
 if (a) {
 foo();
 }
}
if (a) {
 foo();
} else {
 if (a) {
 bar();
 }
}
```
## Options

This rule has no options.

## When Not To Use It

In rare cases where you really need identical test conditions in the same chain, which necessarily means that the expressions in the chain are causing and relying on side effects, you will have to turn this rule off.

## Related Rules

## Version

This rule was introduced in ESLint v6.7.0.

# no-dupe-keys

Disallow duplicate keys in object literals

 ✅ Recommended

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Multiple properties with the same key in object literals can cause unexpected behavior in your application.

```
const foo = {
 bar: "baz",
 bar: "qux"
};
```
## Rule Details

This rule disallows duplicate keys in object literals.

Examples of **incorrect** code for this rule:

 Open in Playground

```
/*eslint no-dupe-keys: "error"*/
const foo = {
 bar: "baz",
 : "qux"
};
const bar = {
 "bar": "baz",
 : "qux"
};
const baz = {
 0x1: "baz",
 : "qux"
};
```
Examples of **correct** code for this rule:

 Open in Playground

```
/*eslint no-dupe-keys: "error"*/
const foo = {
 bar: "baz",
 quxx: "qux"
};
const obj = {
 "__proto__": baz, // defines object's prototype
 ["__proto__"]: qux // defines a property named "__proto__"
};
```
## Options

This rule has no options.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

## Version

This rule was introduced in ESLint v0.0.9.

# no-duplicate-case

Disallow duplicate case labels

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

If a `switch` statement has duplicate test expressions in `case` clauses, it is likely that a programmer copied a `case` clause but forgot to change the test expression.

## Rule Details

This rule disallows duplicate test expressions in `case` clauses of `switch` statements.

Examples of **incorrect** code for this rule:

```
/*eslint no-duplicate-case: "error"*/
const a = 1,
 one = 1;
switch (a) {
 case 1:
 break;
 case 2:
 break;

 default:
 break;
}
switch (a) {
 case one:
 break;
 case 2:
 break;

 default:
 break;
}
switch (a) {
 case "1":
 break;
 case "2":
 break;

 default:
 break;
}
```
Examples of **correct** code for this rule:

```
/*eslint no-duplicate-case: "error"*/
const a = 1,
 one = 1;
switch (a) {
 case 1:
 break;
 case 2:
 break;
 case 3:
 break;
 default:
 break;
}
switch (a) {
 case one:
 break;
 case 2:
 break;
 case 3:
 break;
 default:
 break;
}
switch (a) {
 case "1":
 break;
 case "2":
 break;
 case "3":
 break;
 default:
 break;
}
```
## Options

This rule has no options.

## When Not To Use It

In rare cases where identical test expressions in `case` clauses produce different values, which necessarily means that the expressions are causing and relying on side effects, you will have to disable this rule.

```
switch (a) {
 case i++:
 foo();
 break;
 case i++: // eslint-disable-line no-duplicate-case
 bar();
 break;
}
```
## Version

This rule was introduced in ESLint v0.17.0.

# no-duplicate-imports

Disallow duplicate module imports

Using a single `import` statement per module will make the code clearer because you can see everything being imported from that module on one line.

In the following example the `module` import on line 1 is repeated on line 3. These can be combined to make the list of imports more succinct.

```
import { merge } from 'module';
import something from 'another-module';
import { find } from 'module';
```
## Rule Details

This rule requires that all imports from a single module that can be merged exist in a single `import` statement.

Example of **incorrect** code for this rule:

```
/*eslint no-duplicate-imports: "error"*/
import { merge } from 'module';
import something from 'another-module';
```
Example of **correct** code for this rule:

```
/*eslint no-duplicate-imports: "error"*/
import { merge, find } from 'module';
import something from 'another-module';
```
Example of **correct** code for this rule:

```
/*eslint no-duplicate-imports: "error"*/
// not mergeable
import { merge } from 'module';
import * as something from 'module';
```
## Options

This rule has an object option:

- `"includeExports"`:- `true`(default- `false`) checks for exports in addition to imports.
- `"allowSeparateTypeImports"`:- `true`(default- `false`) allows a type import alongside a value import from the same module in TypeScript files.

### includeExports

If re-exporting from an imported module, you should add the imports to the `import`-statement, and export that directly, not use `export ... from`.

Example of **incorrect** code for this rule with the `{ "includeExports": true }` option:

```
/*eslint no-duplicate-imports: ["error", { "includeExports": true }]*/
import { merge } from 'module';
```
Example of **correct** code for this rule with the `{ "includeExports": true }` option:

```
/*eslint no-duplicate-imports: ["error", { "includeExports": true }]*/
import { merge, find } from 'module';
export { find };
```
Example of **correct** code for this rule with the `{ "includeExports": true }` option:

```
/*eslint no-duplicate-imports: ["error", { "includeExports": true }]*/
import { merge, find } from 'module';
// cannot be merged with the above import
export * as something from 'module';
// cannot be written differently
export * from 'module';
```
### allowSeparateTypeImports

TypeScript allows importing types using `import type`. By default, this rule flags instances of `import type` that have the same specifier as `import`. The `allowSeparateTypeImports` option allows you to override this behavior.

Example of **incorrect** TypeScript code for this rule with the default `{ "allowSeparateTypeImports": false }` option:

```
/*eslint no-duplicate-imports: ["error", { "allowSeparateTypeImports": false }]*/
import { someValue } from 'module';
```
Example of **correct** TypeScript code for this rule with the default `{ "allowSeparateTypeImports": false }` option:

```
/*eslint no-duplicate-imports: ["error", { "allowSeparateTypeImports": false }]*/
import { someValue, type SomeType } from 'module';
```
Example of **incorrect** TypeScript code for this rule with the `{ "allowSeparateTypeImports": true }` option:

```
/*eslint no-duplicate-imports: ["error", { "allowSeparateTypeImports": true }]*/
import { someValue } from 'module';
import type { SomeType } from 'module';
```
Example of **correct** TypeScript code for this rule with the `{ "allowSeparateTypeImports": true }` option:

```
/*eslint no-duplicate-imports: ["error", { "allowSeparateTypeImports": true }]*/
import { someValue } from 'module';
import type { SomeType, AnotherType } from 'module';
```
## Version

This rule was introduced in ESLint v2.5.0.

# no-empty-character-class

Disallow empty character classes in regular expressions

 ✅ Recommended

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Because empty character classes in regular expressions do not match anything, they might be typing mistakes.

```
const foo = /^abc[]/;
```
## Rule Details

This rule disallows empty character classes in regular expressions.

Examples of **incorrect** code for this rule:

 Open in Playground

```
/*eslint no-empty-character-class: "error"*/
.test("abcdefg"); // false
"abcdefg".match(); // null
.test("abcdefg"); // false
"abcdefg".match(); // null
.test("abcdefg"); // false
"abcdefg".match(); // null
.test("abcdefg"); // false
"abcdefg".match(); // null
const regex = ;
regex.test("abcdefg"); // true, the nested `[]` has no effect
"abcdefg".match(regex); // ["abcd"]
regex.test("abcefg"); // false, the nested `[]` has no effect
"abcefg".match(regex); // null
regex.test("abc"); // false, the nested `[]` has no effect
"abc".match(regex); // null
```
Examples of **correct** code for this rule:

 Open in Playground

```
/*eslint no-empty-character-class: "error"*/
/^abc/.test("abcdefg"); // true
"abcdefg".match(/^abc/); // ["abc"]
/^abc[a-z]/.test("abcdefg"); // true
"abcdefg".match(/^abc[a-z]/); // ["abcd"]
/^abc[^]/.test("abcdefg"); // true
"abcdefg".match(/^abc[^]/); // ["abcd"]
```
## Options

This rule has no options.

## Known Limitations

This rule does not report empty character classes in the string argument of calls to the `RegExp` constructor.

Example of a *false negative* when this rule reports correct code:

```
/*eslint no-empty-character-class: "error"*/
const abcNeverMatches = new RegExp("^abc[]");
```
## Version

This rule was introduced in ESLint v0.22.0.

# no-empty-pattern

Disallow empty destructuring patterns

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

When using destructuring, it’s possible to create a pattern that has no effect. This happens when empty curly braces are used to the right of an embedded object destructuring pattern, such as:

```
// doesn't create any variables
const {a: {}} = foo;
```
In this code, no new variables are created because `a` is just a location helper while the `{}` is expected to contain the variables to create, such as:

```
// creates variable b
const {a: { b }} = foo;
```
In many cases, the empty object pattern is a mistake where the author intended to use a default value instead, such as:

```
// creates variable a
const {a = {}} = foo;
```
The difference between these two patterns is subtle, especially because the problematic empty pattern looks just like an object literal.

## Rule Details

This rule aims to flag any empty patterns in destructured objects and arrays, and as such, will report a problem whenever one is encountered.

Examples of **incorrect** code for this rule:

```
/*eslint no-empty-pattern: "error"*/
const = foo;
const = foo;
const {a: } = foo;
const {a: } = foo;
function foo() {}
function bar() {}
function baz({a: }) {}
function qux({a: }) {}
```
Examples of **correct** code for this rule:

```
/*eslint no-empty-pattern: "error"*/
const {a = {}} = foo;
const {b = []} = foo;
function foo({a = {}}) {}
function bar({a = []}) {}
```
## Options

This rule has an object option for exceptions:

### allowObjectPatternsAsParameters

Set to `false` by default. Setting this option to `true` allows empty object patterns as function parameters.

**Note:** This rule doesn’t allow empty array patterns as function parameters.

Examples of **incorrect** code for this rule with the `{"allowObjectPatternsAsParameters": true}` option:

```
/*eslint no-empty-pattern: ["error", { "allowObjectPatternsAsParameters": true }]*/
function foo({a: }) {}
const bar = function({a: }) {};
const qux = ({a: }) => {};
const quux = ( = bar) => {};
const item = ( = { bar: 1 }) => {};
function baz() {}
```
Examples of **correct** code for this rule with the `{"allowObjectPatternsAsParameters": true}` option:

```
/*eslint no-empty-pattern: ["error", { "allowObjectPatternsAsParameters": true }]*/
function foo({}) {}
const bar = function({}) {};
const qux = ({}) => {};
function baz({} = {}) {}
```
## Version

This rule was introduced in ESLint v1.7.0.

# no-ex-assign

Disallow reassigning exceptions in `catch` clauses

 ✅ Recommended

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

If a `catch` clause in a `try` statement accidentally (or purposely) assigns another value to the exception parameter, it is impossible to refer to the error from that point on.
Since there is no `arguments` object to offer alternative access to this data, assignment of the parameter is absolutely destructive.

## Rule Details

This rule disallows reassigning exceptions in `catch` clauses.

Examples of **incorrect** code for this rule:

 Open in Playground

```
/*eslint no-ex-assign: "error"*/
try {
 // code
} catch (e) {
 = 10;
}
```
Examples of **correct** code for this rule:

 Open in Playground

```
/*eslint no-ex-assign: "error"*/
try {
 // code
} catch (e) {
 const foo = 10;
}
```
## Options

This rule has no options.

## Version

This rule was introduced in ESLint v0.0.9.

# no-fallthrough

Disallow fallthrough of `case` statements

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

The `switch` statement in JavaScript is one of the more error-prone constructs of the language thanks in part to the ability to “fall through” from one `case` to the next. For example:

```
switch(foo) {
 case 1:
 doSomething();
 case 2:
 doSomethingElse();
}
```
In this example, if `foo` is `1`, then execution will flow through both cases, as the first falls through to the second. You can prevent this by using `break`, as in this example:

```
switch(foo) {
 case 1:
 doSomething();
 break;
 case 2:
 doSomethingElse();
}
```
That works fine when you don’t want a fallthrough, but what if the fallthrough is intentional, there is no way to indicate that in the language. It’s considered a best practice to always indicate when a fallthrough is intentional using a comment which matches the `/falls?\s?through/i` regular expression but isn’t a directive:

```
switch(foo) {
 case 1:
 doSomething();
 // falls through
 case 2:
 doSomethingElse();
}
switch(foo) {
 case 1:
 doSomething();
 // fall through
 case 2:
 doSomethingElse();
}
switch(foo) {
 case 1:
 doSomething();
 // fallsthrough
 case 2:
 doSomethingElse();
}
switch(foo) {
 case 1: {
 doSomething();
 // falls through
 }
 case 2: {
 doSomethingElse();
 }
}
```
In this example, there is no confusion as to the expected behavior. It is clear that the first case is meant to fall through to the second case.

## Rule Details

This rule is aimed at eliminating unintentional fallthrough of one case to the other. As such, it flags any fallthrough scenarios that are not marked by a comment.

Examples of **incorrect** code for this rule:

```
/*eslint no-fallthrough: "error"*/
switch(foo) {
 case 1:
 doSomething();

}
```
Examples of **correct** code for this rule:

```
/*eslint no-fallthrough: "error"*/
switch(foo) {
 case 1:
 doSomething();
 break;
 case 2:
 doSomething();
}
function bar(foo) {
 switch(foo) {
 case 1:
 doSomething();
 return;
 case 2:
 doSomething();
 }
}
switch(foo) {
 case 1:
 doSomething();
 throw new Error("Boo!");
 case 2:
 doSomething();
}
switch(foo) {
 case 1:
 case 2:
 doSomething();
}
switch(foo) {
 case 1: case 2:
 doSomething();
}
switch(foo) {
 case 1:
 doSomething();
 // falls through
 case 2:
 doSomething();
}
switch(foo) {
 case 1: {
 doSomething();
 // falls through
 }
 case 2: {
 doSomethingElse();
 }
}
```
Note that the last `case` statement in these examples does not cause a warning because there is nothing to fall through into.

## Options

This rule has an object option:

-
Set the `commentPattern`option to a regular expression string to change the test for intentional fallthrough comment. If the fallthrough comment matches a directive, that takes precedence over`commentPattern`.
-
Set the `allowEmptyCase`option to`true`to allow empty cases regardless of the layout. By default, this rule does not require a fallthrough comment after an empty`case`only if the empty`case`and the next`case`are on the same line or on consecutive lines.
-
Set the `reportUnusedFallthroughComment`option to`true`to prohibit a fallthrough comment from being present if the case cannot fallthrough due to being unreachable. This is mostly intended to help avoid misleading comments occurring as a result of refactoring.

### commentPattern

Examples of **correct** code for the `{ "commentPattern": "break[\\s\\w]*omitted" }` option:

```
/*eslint no-fallthrough: ["error", { "commentPattern": "break[\\s\\w]*omitted" }]*/
switch(foo) {
 case 1:
 doSomething();
 // break omitted
 case 2:
 doSomething();
}
switch(foo) {
 case 1:
 doSomething();
 // caution: break is omitted intentionally
 default:
 doSomething();
}
```
### allowEmptyCase

Examples of **correct** code for the `{ "allowEmptyCase": true }` option:

```
/* eslint no-fallthrough: ["error", { "allowEmptyCase": true }] */
switch(foo){
 case 1:
 case 2: doSomething();
}
switch(foo){
 case 1:
 /*
 Put a message here
 */
 case 2: doSomething();
}
```
### reportUnusedFallthroughComment

Examples of **incorrect** code for the `{ "reportUnusedFallthroughComment": true }` option:

```
/* eslint no-fallthrough: ["error", { "reportUnusedFallthroughComment": true }] */
switch(foo){
 case 1:
 doSomething();
 break;

 case 2: doSomething();
}
function f() {
 switch(foo){
 case 1:
 if (a) {
 throw new Error();
 } else if (b) {
 break;
 } else {
 return;
 }

 case 2:
 break;
 }
}
```
Examples of **correct** code for the `{ "reportUnusedFallthroughComment": true }` option:

```
/* eslint no-fallthrough: ["error", { "reportUnusedFallthroughComment": true }] */
switch(foo){
 case 1:
 doSomething();
 break;
 // just a comment
 case 2: doSomething();
}
```
## When Not To Use It

If you don’t want to enforce that each `case` statement should end with a `throw`, `return`, `break`, or comment, then you can safely turn this rule off.

## Related Rules

## Version

This rule was introduced in ESLint v0.0.7.

# no-func-assign

Disallow reassigning `function` declarations

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

JavaScript functions can be written as a FunctionDeclaration `function foo() { ... }` or as a FunctionExpression `const foo = function() { ... };`. While a JavaScript interpreter might tolerate it, overwriting/reassigning a function written as a FunctionDeclaration is often indicative of a mistake or issue.

```
function foo() {}
foo = bar;
```
## Rule Details

This rule disallows reassigning `function` declarations.

Examples of **incorrect** code for this rule:

```
/*eslint no-func-assign: "error"*/
function foo() {}
 = bar;
function baz() {
 = bar;
}
let a = function hello() {
 = 123;
};
```
Examples of **incorrect** code for this rule, unlike the corresponding rule in JSHint:

```
/*eslint no-func-assign: "error"*/
 = bar;
function foo() {}
```
Examples of **correct** code for this rule:

```
/*eslint no-func-assign: "error"*/
let foo = function () {}
foo = bar;
function baz(baz) { // `baz` is shadowed.
 baz = bar;
}
function qux() {
 const qux = bar; // `qux` is shadowed.
}
```
## Options

This rule has no options.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

## Version

This rule was introduced in ESLint v0.0.9.

# no-import-assign

Disallow assigning to imported bindings

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

The updates of imported bindings by ES Modules cause runtime errors.

## Rule Details

This rule warns the assignments, increments, and decrements of imported bindings.

Examples of **incorrect** code for this rule:

```
/*eslint no-import-assign: "error"*/
import mod, { named } from "./mod.mjs"
import * as mod_ns from "./mod.mjs"
 // ERROR: 'mod' is readonly.
 // ERROR: 'named' is readonly.
 // ERROR: The members of 'mod_ns' are readonly.
 // ERROR: 'mod_ns' is readonly.
// Can't extend 'mod_ns'
 // ERROR: The members of 'mod_ns' are readonly.
```
Examples of **correct** code for this rule:

```
/*eslint no-import-assign: "error"*/
import mod, { named } from "./mod.mjs"
import * as mod_ns from "./mod.mjs"
mod.prop = 1
named.prop = 2
mod_ns.named.prop = 3
// Known Limitation
function test(obj) {
 obj.named = 4 // Not errored because 'obj' is not namespace objects.
}
test(mod_ns) // Not errored because it doesn't know that 'test' updates the member of the argument.
```
## Options

This rule has no options.

## When Not To Use It

If you don’t want to be notified about modifying imported bindings, you can disable this rule.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

Note that the compiler will not catch the `Object.assign()` case. Thus, if you use `Object.assign()` in your codebase, this rule will still provide some value.

## Version

This rule was introduced in ESLint v6.4.0.

# no-inner-declarations

Disallow variable or `function` declarations in nested blocks

In JavaScript, prior to ES6, a function declaration is only allowed in the first level of a program or the body of another function, though parsers sometimes erroneously accept them elsewhere. This only applies to function declarations; named or anonymous function expressions can occur anywhere an expression is permitted.

```
// Good
function doSomething() { }
// Bad
if (test) {
 function doSomethingElse () { }
}
function anotherThing() {
 var fn;
 if (test) {
 // Good
 fn = function expression() { };
 // Bad
 function declaration() { }
 }
}
```
In ES6, block-level functions (functions declared inside a block) are limited to the scope of the block they are declared in and outside of the block scope they can’t be accessed and called, but only when the code is in strict mode (code with `"use strict"` tag or ESM modules). In non-strict mode, they can be accessed and called outside of the block scope.

```
"use strict";
if (test) {
 function doSomething () { }
 doSomething(); // no error
}
doSomething(); // error
```
A variable declaration is permitted anywhere a statement can go, even nested deeply inside other blocks. This is often undesirable due to variable hoisting, and moving declarations to the root of the program or function body can increase clarity. Note that block bindings (`let`, `const`) are not hoisted and therefore they are not affected by this rule.

```
// Good
var foo = 42;
// Good
if (foo) {
 let bar1;
}
// Bad
while (test) {
 var bar2;
}
function doSomething() {
 // Good
 var baz = true;
 // Bad
 if (baz) {
 var quux;
 }
}
```
## Rule Details

This rule requires that function declarations and, optionally, variable declarations be in the root of a program, or in the root of the body of a function, or in the root of the body of a class static block.

## Options

This rule has a string and an object option:

- `"functions"`(default) disallows- `function`declarations in nested blocks
- `"both"`disallows- `function`and- `var`declarations in nested blocks
- `{ blockScopedFunctions: "allow" }`(default) this option allows- `function`declarations in nested blocks when code is in strict mode (code with- `"use strict"`tag or ESM modules) and- `languageOptions.ecmaVersion`is set to- `2015`or above. This option can be disabled by setting it to- `"disallow"`.

### functions

Examples of **incorrect** code for this rule with the default `"functions"` option:

```
/*eslint no-inner-declarations: "error"*/
// script, non-strict code
if (test) {

}
function doSomethingElse() {
 if (test) {

 }
}
if (foo)
```
Examples of **correct** code for this rule with the default `"functions"` option:

```
/*eslint no-inner-declarations: "error"*/
function doSomething() { }
function doSomethingElse() {
 function doAnotherThing() { }
}
function doSomethingElse() {
 "use strict";
 if (test) {
 function doAnotherThing() { }
 }
}
class C {
 static {
 function doSomething() { }
 }
}
if (test) {
 asyncCall(id, function (err, data) { });
}
var fn;
if (test) {
 fn = function fnExpression() { };
}
if (foo) var a;
```
### both

Examples of **incorrect** code for this rule with the `"both"` option:

```
/*eslint no-inner-declarations: ["error", "both"]*/
if (test) {

}
function doAnotherThing() {
 if (test) {

 }
}
if (foo)
if (foo)
class C {
 static {
 if (test) {

 }
 }
}
```
Examples of **correct** code for this rule with the `"both"` option:

```
/*eslint no-inner-declarations: ["error", "both"]*/
var bar = 42;
if (test) {
 let baz = 43;
}
function doAnotherThing() {
 var baz = 81;
}
class C {
 static {
 var something;
 }
}
```
### blockScopedFunctions

Example of **incorrect** code for this rule with `{ blockScopedFunctions: "disallow" }` option with `ecmaVersion: 2015`:

```
/*eslint no-inner-declarations: ["error", "functions", { blockScopedFunctions: "disallow" }]*/
// non-strict code
if (test) {

}
function doSomething() {
 if (test) {

 }
}
// strict code
function foo() {
 "use strict";
 if (test) {

 }
}
```
Example of **correct** code for this rule with `{ blockScopedFunctions: "disallow" }` option with `ecmaVersion: 2015`:

```
/*eslint no-inner-declarations: ["error", "functions", { blockScopedFunctions: "disallow" }]*/
function doSomething() { }
function doSomething() {
 function doSomethingElse() { }
}
```
Example of **correct** code for this rule with `{ blockScopedFunctions: "allow" }` option with `ecmaVersion: 2015`:

```
/*eslint no-inner-declarations: ["error", "functions", { blockScopedFunctions: "allow" }]*/
"use strict";
if (test) {
 function doSomething() { }
}
function doSomething() {
 if (test) {
 function doSomethingElse() { }
 }
}
// OR
function foo() {
 "use strict";
 if (test) {
 function bar() { }
 }
}
```
`ESM modules` and both `class` declarations and expressions are always in strict mode.

```
/*eslint no-inner-declarations: ["error", "functions", { blockScopedFunctions: "allow" }]*/
if (test) {
 function doSomething() { }
}
function doSomethingElse() {
 if (test) {
 function doAnotherThing() { }
 }
}
class Some {
 static {
 if (test) {
 function doSomething() { }
 }
 }
}
const C = class {
 static {
 if (test) {
 function doSomething() { }
 }
 }
}
```
## When Not To Use It

By default, this rule disallows inner function declarations only in contexts where their behavior is unspecified and thus inconsistent (pre-ES6 environments) or legacy semantics apply (non-strict mode code). If your code targets pre-ES6 environments or is not in strict mode, you should enable this rule to prevent unexpected behavior.

In ES6+ environments, in strict mode code, the behavior of inner function declarations is well-defined and consistent - they are always block-scoped. If your code targets only ES6+ environments and is in strict mode (ES modules, or code with `"use strict"` directives) then there is no need to enable this rule unless you want to disallow inner functions as a stylistic choice, in which case you should enable this rule with the option `blockScopedFunctions: "disallow"`.

Disable checking variable declarations when using block-scoped-var or if declaring variables in nested blocks is acceptable despite hoisting.

## Version

This rule was introduced in ESLint v0.6.0.

# no-invalid-regexp

Disallow invalid regular expression strings in `RegExp` constructors

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

An invalid pattern in a regular expression literal is a `SyntaxError` when the code is parsed, but an invalid string in `RegExp` constructors throws a `SyntaxError` only when the code is executed.

## Rule Details

This rule disallows invalid regular expression strings in `RegExp` constructors.

Examples of **incorrect** code for this rule:

```
/*eslint no-invalid-regexp: "error"*/
```
Examples of **correct** code for this rule:

```
/*eslint no-invalid-regexp: "error"*/
RegExp('.')
new RegExp
this.RegExp('[')
```
Please note that this rule validates regular expressions per the latest ECMAScript specification, regardless of your parser settings.

If you want to allow additional constructor flags for any reason, you can specify them using the `allowConstructorFlags` option. These flags will then be ignored by the rule.

## Options

This rule has an object option for exceptions:

- `"allowConstructorFlags"`is a case-sensitive array of flags

### allowConstructorFlags

Examples of **correct** code for this rule with the `{ "allowConstructorFlags": ["a", "z"] }` option:

```
/*eslint no-invalid-regexp: ["error", { "allowConstructorFlags": ["a", "z"] }]*/
new RegExp('.', 'a')
new RegExp('.', 'az')
```
## Version

This rule was introduced in ESLint v0.1.4.

# no-irregular-whitespace

Disallow irregular whitespace

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Invalid or irregular whitespace causes issues with ECMAScript 5 parsers and also makes code harder to debug in a similar nature to mixed tabs and spaces.

Various whitespace characters can be inputted by programmers by mistake for example from copying or keyboard shortcuts. Pressing Alt + Space on macOS adds in a non breaking space character for example.

A simple fix for this problem could be to rewrite the offending line from scratch. This might also be a problem introduced by the text editor: if rewriting the line does not fix it, try using a different editor.

Known issues these spaces cause:

- Ogham Space Mark
- Is a valid token separator, but is rendered as a visible glyph in most typefaces, which may be misleading in source code.

- Mongolian Vowel Separator
- Is no longer considered a space separator since Unicode 6.3. It will result in a syntax error in current parsers when used in place of a regular token separator.

- Line Separator and Paragraph Separator
- These have always been valid whitespace characters and line terminators, but were considered illegal in string literals prior to ECMAScript 2019.

- Zero Width Space
- Is NOT considered a separator for tokens and is often parsed as an `Unexpected token ILLEGAL`.
- Is NOT shown in modern browsers making code repository software expected to resolve the visualization.

- Is NOT considered a separator for tokens and is often parsed as an

In JSON, none of the characters listed as irregular whitespace by this rule may appear outside of a string.

## Rule Details

This rule is aimed at catching invalid whitespace that is not a normal tab and space. Some of these characters may cause issues in modern browsers and others will be a debugging issue to spot.

This rule disallows the following characters except where the options allow:

```
\u000B - Line Tabulation (\v) - <VT>
\u000C - Form Feed (\f) - <FF>
\u0085 - Next Line - <NEL>
\u00A0 - No-Break Space - <NBSP>
\u1680 - Ogham Space Mark - <OGSP>
\u180E - Mongolian Vowel Separator - <MVS>
\u2000 - En Quad - <NQSP>
\u2001 - Em Quad - <MQSP>
\u2002 - En Space - <ENSP>
\u2003 - Em Space - <EMSP>
\u2004 - Three-Per-Em - <THPMSP> - <3/MSP>
\u2005 - Four-Per-Em - <FPMSP> - <4/MSP>
\u2006 - Six-Per-Em - <SPMSP> - <6/MSP>
\u2007 - Figure Space - <FSP>
\u2008 - Punctuation Space - <PUNCSP>
\u2009 - Thin Space - <THSP>
\u200A - Hair Space - <HSP>
\u200B - Zero Width Space - <ZWSP>
\u2028 - Line Separator - <LS> - <LSEP>
\u2029 - Paragraph Separator - <PS> - <PSEP>
\u202F - Narrow No-Break Space - <NNBSP>
\u205F - Medium Mathematical Space - <MMSP>
\u3000 - Ideographic Space - <IDSP>
\uFEFF - Zero Width No-Break Space - <BOM>
```
## Options

This rule has an object option for exceptions:

- `"skipStrings": true`(default) allows any whitespace characters in string literals
- `"skipComments": true`allows any whitespace characters in comments
- `"skipRegExps": true`allows any whitespace characters in regular expression literals
- `"skipTemplates": true`allows any whitespace characters in template literals
- `"skipJSXText": true`allows any whitespace characters in JSX text

### skipStrings

Examples of **incorrect** code for this rule with the default `{ "skipStrings": true }` option:

```
/*eslint no-irregular-whitespace: "error"*/
const thing = function()/*<NBSP>*/{
 return 'test';
}
const foo = function(/*<NBSP>*/){
 return 'test';
}
const bar = function/*<NBSP>*/(){
 return 'test';
}
const baz = function/*<Ogham Space Mark>*/(){
 return 'test';
}
const qux = function() {
 return 'test';/*<ENSP>*/
}
const quux = function() {
 return 'test';/*<NBSP>*/
}
const item = function() {
 // Description<NBSP>: some descriptive text
}
/*
Description<NBSP>: some descriptive text
*/
const func = function() {
 return /<NBSP>regexp/;
}
const myFunc = function() {
 return `template<NBSP>string`;
}
```
Examples of **correct** code for this rule with the default `{ "skipStrings": true }` option:

```
/*eslint no-irregular-whitespace: "error"*/
const thing = function() {
 return ' <NBSP>thing';
}
const foo = function() {
 return '<ZWSP>thing';
}
const bar = function() {
 return 'th <NBSP>ing';
}
```
### skipComments

Examples of additional **correct** code for this rule with the `{ "skipComments": true }` option:

```
/*eslint no-irregular-whitespace: ["error", { "skipComments": true }]*/
function thing() {
 // Description <NBSP>: some descriptive text
}
/*
Description <NBSP>: some descriptive text
*/
```
### skipRegExps

Examples of additional **correct** code for this rule with the `{ "skipRegExps": true }` option:

```
/*eslint no-irregular-whitespace: ["error", { "skipRegExps": true }]*/
function thing() {
 return / <NBSP>regexp/;
}
```
### skipTemplates

Examples of additional **correct** code for this rule with the `{ "skipTemplates": true }` option:

```
/*eslint no-irregular-whitespace: ["error", { "skipTemplates": true }]*/
function thing() {
 return `template <NBSP>string`;
}
```
### skipJSXText

Examples of additional **correct** code for this rule with the `{ "skipJSXText": true }` option:

```
/*eslint no-irregular-whitespace: ["error", { "skipJSXText": true }]*/
function Thing() {
 return <div>text in JSX</div>; // <NBSP> before `JSX`
}
```
## When Not To Use It

If you decide that you wish to use whitespace other than tabs and spaces outside of strings in your application.

## Version

This rule was introduced in ESLint v0.9.0.

# no-loss-of-precision

Disallow literal numbers that lose precision

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

This rule would disallow the use of number literals that lose precision at runtime when converted to a JS `Number` due to 64-bit floating-point rounding.

## Rule Details

In JS, `Number`s are stored as double-precision floating-point numbers according to the IEEE 754 standard. Because of this, numbers can only retain accuracy up to a certain amount of digits. If the programmer enters additional digits, those digits will be lost in the conversion to the `Number` type and will result in unexpected behavior.

Examples of **incorrect** code for this rule:

```
/*eslint no-loss-of-precision: "error"*/
const a =
const b =
const c =
const d =
const e =
const f = ;
```
Examples of **correct** code for this rule:

```
/*eslint no-loss-of-precision: "error"*/
const a = 12345
const b = 123.456
const c = 123e34
const d = 12300000000000000000000000
const e = 0x1FFFFFFFFFFFFF
const f = 9007199254740991
const g = 9007_1992547409_91
```
## Options

This rule has no options.

## Version

This rule was introduced in ESLint v7.1.0.

# no-misleading-character-class

Disallow characters which are made with multiple code points in character class syntax

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Some problems reported by this rule are manually fixable by editor suggestions

Unicode includes characters which are made by multiple code points.
RegExp character class syntax (`/[abc]/`) cannot handle characters which are made by multiple code points as a character; those characters will be dissolved to each code point. For example, `❇️` is made by `❇` (`U+2747`) and VARIATION SELECTOR-16 (`U+FE0F`). If this character is in a RegExp character class, it will match either `❇` (`U+2747`) or VARIATION SELECTOR-16 (`U+FE0F`) rather than `❇️`.

This rule reports regular expressions which include multiple code point characters in character class syntax. This rule considers the following characters as multiple code point characters.

**A character with combining characters:**

The combining characters are characters which belong to one of `Mc`, `Me`, and `Mn` Unicode general categories.

```
/^[Á]$/u.test("Á"); //→ false
/^[❇️]$/u.test("❇️"); //→ false
```
**A character with Emoji modifiers:**

```
/^[👶🏻]$/u.test("👶🏻"); //→ false
/^[👶🏽]$/u.test("👶🏽"); //→ false
```
**A pair of regional indicator symbols:**

```
/^[🇯🇵]$/u.test("🇯🇵"); //→ false
```
**Characters that ZWJ joins:**

```
/^[👨👩👦]$/u.test("👨👩👦"); //→ false
```
**A surrogate pair without Unicode flag:**

```
/^[👍]$/.test("👍"); //→ false
// Surrogate pair is OK if with u flag.
/^[👍]$/u.test("👍"); //→ true
```
## Rule Details

This rule reports regular expressions which include multiple code point characters in character class syntax.

Examples of **incorrect** code for this rule:

```
/*eslint no-misleading-character-class: error */
/^[]$/u;
/^[]$/u;
/^[]$/u;
/^[]$/u;
/^[]$/u;
/^[]$/;
new RegExp("[]");
```
Examples of **correct** code for this rule:

```
/*eslint no-misleading-character-class: error */
/^[abc]$/;
/^[👍]$/u;
/^[\q{👶🏻}]$/v;
new RegExp("^[]$");
new RegExp(`[Á-${z}]`, "u"); // variable pattern
```
## Options

This rule has an object option:

- `"allowEscape"`: When set to- `true`, the rule allows any grouping of code points inside a character class as long as they are written using escape sequences. This option only has effect on regular expression literals and on regular expressions created with the- `RegExp`constructor with a literal argument as a pattern.

### allowEscape

Examples of **incorrect** code for this rule with the `{ "allowEscape": true }` option:

```
/* eslint no-misleading-character-class: ["error", { "allowEscape": true }] */
/[]/; // backslash can be omitted
new RegExp();
const pattern = "[\ud83d\udc4d]";
new RegExp();
```
Examples of **correct** code for this rule with the `{ "allowEscape": true }` option:

```
/* eslint no-misleading-character-class: ["error", { "allowEscape": true }] */
/[\ud83d\udc4d]/;
/[\u00B7\u0300-\u036F]/u;
/[👨\u200d👩]/u;
new RegExp("[\x41\u0301]");
new RegExp(`[\u{1F1EF}\u{1F1F5}]`, "u");
new RegExp("[\\u{1F1EF}\\u{1F1F5}]", "u");
```
## When Not To Use It

You can turn this rule off if you don’t want to check RegExp character class syntax for multiple code point characters.

## Version

This rule was introduced in ESLint v5.3.0.

# no-new-native-nonconstructor

Disallow `new` operators with global non-constructor functions

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

It is a convention in JavaScript that global variables beginning with an uppercase letter typically represent classes that can be instantiated using the `new` operator, such as `new Array` and `new Map`. Confusingly, JavaScript also provides some global variables that begin with an uppercase letter that cannot be called using the `new` operator and will throw an error if you attempt to do so. These are typically functions that are related to data types and are easy to mistake for classes. Consider the following example:

```
// throws a TypeError
const foo = new Symbol("foo");
// throws a TypeError
const result = new BigInt(9007199254740991);
```
Both `new Symbol` and `new BigInt` throw a type error because they are functions and not classes. It is easy to make this mistake by assuming the uppercase letters indicate classes.

## Rule Details

This rule is aimed at preventing the accidental calling of native JavaScript global functions with the `new` operator. These functions are:

- `Symbol`
- `BigInt`

Examples of **incorrect** code for this rule:

```
/*eslint no-new-native-nonconstructor: "error"*/
const foo = new ('foo');
const bar = new (9007199254740991);
```
Examples of **correct** code for this rule:

```
/*eslint no-new-native-nonconstructor: "error"*/
const foo = Symbol('foo');
const bar = BigInt(9007199254740991);
// Ignores shadowed Symbol.
function baz(Symbol) {
 const qux = new Symbol("baz");
}
function quux(BigInt) {
 const corge = new BigInt(9007199254740991);
}
```
## Options

This rule has no options.

## When Not To Use It

This rule should not be used in ES3/5 environments.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

## Related Rules

## Version

This rule was introduced in ESLint v8.27.0.

# no-obj-calls

Disallow calling global object properties as functions

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

ECMAScript provides several global objects that are intended to be used as-is. Some of these objects look as if they could be constructors due to their capitalization (such as `Math`, `JSON`, and `Temporal`) but will throw an error if you try to execute them as functions.

The ECMAScript 5 specification makes it clear that both `Math` and `JSON` cannot be invoked:

The Math object does not have a

`[[Call]]`internal property; it is not possible to invoke the Math object as a function.

The ECMAScript 2015 specification makes it clear that `Reflect` cannot be invoked:

The Reflect object also does not have a

`[[Call]]`internal method; it is not possible to invoke the Reflect object as a function.

The ECMAScript 2017 specification makes it clear that `Atomics` cannot be invoked:

The Atomics object does not have a

`[[Call]]`internal method; it is not possible to invoke the Atomics object as a function.

The ECMAScript Internationalization API Specification makes it clear that `Intl` cannot be invoked:

The Intl object does not have a

`[[Call]]`internal method; it is not possible to invoke the Intl object as a function.

The Temporal proposal specification makes it clear that `Temporal` cannot be invoked:

The Temporal object does not have a

`[[Call]]`internal method; it cannot be invoked as a function.

## Rule Details

This rule disallows calling the `Math`, `JSON`, `Reflect`, `Atomics`, `Intl`, and `Temporal` objects as functions.

This rule also disallows using these objects as constructors with the `new` operator.

Examples of **incorrect** code for this rule:

```
/*eslint no-obj-calls: "error"*/
const math = ;
const newMath = ;
const json = ;
const newJSON = ;
const reflect = ;
const newReflect = ;
const atomics = ;
const newAtomics = ;
const intl = ;
const newIntl = ;
const temporal = ;
const newTemporal = ;
```
Examples of **correct** code for this rule:

```
/*eslint no-obj-calls: "error"*/
function area(r) {
 return Math.PI * r * r;
}
const object = JSON.parse("{}");
const value = Reflect.get({ x: 1, y: 2 }, "x");
const first = Atomics.load(foo, 0);
const segmenterFr = new Intl.Segmenter("fr", { granularity: "word" });
const instant = Temporal.Now.instant();
```
## Options

This rule has no options.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

## Version

This rule was introduced in ESLint v0.0.9.

# no-promise-executor-return

Disallow returning values from Promise executor functions

Some problems reported by this rule are manually fixable by editor suggestions

The `new Promise` constructor accepts a single argument, called an *executor*.

```
const myPromise = new Promise(function executor(resolve, reject) {
 readFile('foo.txt', function(err, result) {
 if (err) {
 reject(err);
 } else {
 resolve(result);
 }
 });
});
```
The executor function usually initiates some asynchronous operation. Once it is finished, the executor should call `resolve` with the result, or `reject` if an error occurred.

The return value of the executor is ignored. Returning a value from an executor function is a possible error because the returned value cannot be used and it doesn’t affect the promise in any way.

## Rule Details

This rule disallows returning values from Promise executor functions.

Only `return` without a value is allowed, as it’s a control flow statement.

Examples of **incorrect** code for this rule:

```
/*eslint no-promise-executor-return: "error"*/
new Promise((resolve, reject) => {
 if (someCondition) {

 }
 getSomething((err, result) => {
 if (err) {
 reject(err);
 } else {
 resolve(result);
 }
 });
});
new Promise((resolve, reject) => );
new Promise(() => {

});
new Promise(r => );
```
Examples of **correct** code for this rule:

```
/*eslint no-promise-executor-return: "error"*/
// Turn return inline into two lines
new Promise((resolve, reject) => {
 if (someCondition) {
 resolve(defaultResult);
 return;
 }
 getSomething((err, result) => {
 if (err) {
 reject(err);
 } else {
 resolve(result);
 }
 });
});
// Add curly braces
new Promise((resolve, reject) => {
 getSomething((err, data) => {
 if (err) {
 reject(err);
 } else {
 resolve(data);
 }
 });
});
new Promise(r => { r(1) });
// or just use Promise.resolve
Promise.resolve(1);
```
## Options

This rule takes one option, an object, with the following properties:

- `allowVoid`: If set to- `true`(- `false`by default), this rule will allow returning void values.

### allowVoid

Examples of **correct** code for this rule with the `{ "allowVoid": true }` option:

```
/*eslint no-promise-executor-return: ["error", { allowVoid: true }]*/
new Promise((resolve, reject) => {
 if (someCondition) {
 return void resolve(defaultResult);
 }
 getSomething((err, result) => {
 if (err) {
 reject(err);
 } else {
 resolve(result);
 }
 });
});
new Promise((resolve, reject) => void getSomething((err, data) => {
 if (err) {
 reject(err);
 } else {
 resolve(data);
 }
}));
new Promise(r => void r(1));
```
## Related Rules

## Version

This rule was introduced in ESLint v7.3.0.

# no-prototype-builtins

Disallow calling some `Object.prototype` methods directly on objects

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Some problems reported by this rule are manually fixable by editor suggestions

In ECMAScript 5.1, `Object.create` was added, which enables the creation of objects with a specified `[[Prototype]]`. `Object.create(null)` is a common pattern used to create objects that will be used as a Map. This can lead to errors when it is assumed that objects will have properties from `Object.prototype`. This rule prevents calling some `Object.prototype` methods directly from an object.

Additionally, objects can have properties that shadow the builtins on `Object.prototype`, potentially causing unintended behavior or denial-of-service security vulnerabilities. For example, it would be unsafe for a webserver to parse JSON input from a client and call `hasOwnProperty` directly on the resulting object, because a malicious client could send a JSON value like `{"hasOwnProperty": 1}` and cause the server to crash.

To avoid subtle bugs like this, it’s better to always call these methods from `Object.prototype`. For example, `foo.hasOwnProperty("bar")` should be replaced with `Object.prototype.hasOwnProperty.call(foo, "bar")`.

## Rule Details

This rule disallows calling some `Object.prototype` methods directly on object instances.

Examples of **incorrect** code for this rule:

```
/*eslint no-prototype-builtins: "error"*/
const hasBarProperty = foo.("bar");
const isPrototypeOfBar = foo.(bar);
const barIsEnumerable = foo.("bar");
```
Examples of **correct** code for this rule:

```
/*eslint no-prototype-builtins: "error"*/
const hasBarProperty = Object.prototype.hasOwnProperty.call(foo, "bar");
const isPrototypeOfBar = Object.prototype.isPrototypeOf.call(foo, bar);
const barIsEnumerable = {}.propertyIsEnumerable.call(foo, "bar");
```
## Options

This rule has no options.

## When Not To Use It

You may want to turn this rule off if your code only touches objects with hardcoded keys, and you will never use an object that shadows an `Object.prototype` method or which does not inherit from `Object.prototype`.

## Version

This rule was introduced in ESLint v2.11.0.

# no-self-assign

Disallow assignments where both sides are exactly the same

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Self assignments have no effect, so probably those are an error due to incomplete refactoring. Those indicate that what you should do is still remaining.

```
foo = foo;
[bar, baz] = [bar, qiz];
```
## Rule Details

This rule is aimed at eliminating self assignments.

Examples of **incorrect** code for this rule:

```
/*eslint no-self-assign: "error"*/
foo = ;
[a, b] = [, ];
[a, ...b] = [x, ...];
({a, b} = {, x});
foo &&= ;
foo ||= ;
foo ??= ;
```
Examples of **correct** code for this rule:

```
/*eslint no-self-assign: "error"*/
foo = bar;
[a, b] = [b, a];
// This pattern is warned by the `no-use-before-define` rule.
let foo = foo;
// The default values have an effect.
[foo = 1] = [foo];
// non-self-assignments with properties.
obj.a = obj.b;
obj.a.b = obj.c.b;
obj.a.b = obj.a.c;
obj[a] = obj["a"];
// This ignores if there is a function call.
obj.a().b = obj.a().b;
a().b = a().b;
// `&=` and `|=` have an effect on non-integers.
foo &= foo;
foo |= foo;
// Known limitation: this does not support computed properties except single literal or single identifier.
obj[a + b] = obj[a + b];
obj["a" + "b"] = obj["a" + "b"];
```
## Options

This rule has the option to check properties as well.

```
{
 "no-self-assign": ["error", {"props": true}]
}
```
- `props`- if this is- `true`,- `no-self-assign`rule warns self-assignments of properties. Default is- `true`.

### props

Examples of **correct** code with the `{ "props": false }` option:

```
/*eslint no-self-assign: ["error", {"props": false}]*/
// self-assignments with properties.
obj.a = obj.a;
obj.a.b = obj.a.b;
obj["a"] = obj["a"];
obj[a] = obj[a];
```
## When Not To Use It

If you don’t want to notify about self assignments, then it’s safe to disable this rule.

## Version

This rule was introduced in ESLint v2.0.0-rc.0.

# no-self-compare

Disallow comparisons where both sides are exactly the same

Comparing a variable against itself is usually an error, either a typo or refactoring error. It is confusing to the reader and may potentially introduce a runtime error.

The only time you would compare a variable against itself is when you are testing for `NaN`. However, it is far more appropriate to use `typeof x === 'number' && isNaN(x)` or the Number.isNaN ES2015 function for that use case rather than leaving the reader of the code to determine the intent of self comparison.

## Rule Details

This error is raised to highlight a potentially confusing and potentially pointless piece of code. There are almost no situations in which you would need to compare something to itself.

Examples of **incorrect** code for this rule:

```
/*eslint no-self-compare: "error"*/
let x = 10;
if () {
 x = 20;
}
```
## Options

This rule has no options.

## Known Limitations

This rule works by directly comparing the tokens on both sides of the operator. It flags them as problematic if they are structurally identical. However, it doesn’t consider possible side effects or that functions may return different objects even when called with the same arguments. As a result, it can produce false positives in some cases, such as:

```
/*eslint no-self-compare: "error"*/
function parseDate(dateStr) {
 return new Date(dateStr);
}
if () {
 // do something
}
let counter = 0;
function incrementUnlessReachedMaximum() {
 return Math.min(counter += 1, 10);
}
if () {
 // ...
}
```
## Version

This rule was introduced in ESLint v0.0.9.

# no-setter-return

Disallow returning values from setters

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Setters cannot return values.

While returning a value from a setter does not produce an error, the returned value is being ignored. Therefore, returning a value from a setter is either unnecessary or a possible error, since the returned value cannot be used.

## Rule Details

This rule disallows returning values from setters and reports `return` statements in setter functions.

Only `return` without a value is allowed, as it’s a control flow statement.

This rule checks setters in:

- Object literals.
- Class declarations and class expressions.
- Property descriptors in `Object.create`,`Object.defineProperty`,`Object.defineProperties`, and`Reflect.defineProperty`methods of the global objects.

Examples of **incorrect** code for this rule:

```
/*eslint no-setter-return: "error"*/
const foo = {
 set a(value) {
 this.val = value;

 }
};
class Foo {
 set a(value) {
 this.val = value * 2;

 }
}
const Bar = class {
 static set a(value) {
 if (value < 0) {
 this.val = 0;

 }
 this.val = value;
 }
};
Object.defineProperty(foo, "bar", {
 set(value) {
 if (value < 0) {

 }
 this.val = value;
 }
});
```
Examples of **correct** code for this rule:

```
/*eslint no-setter-return: "error"*/
const foo = {
 set a(value) {
 this.val = value;
 }
};
class Foo {
 set a(value) {
 this.val = value * 2;
 }
}
const Bar = class {
 static set a(value) {
 if (value < 0) {
 this.val = 0;
 return;
 }
 this.val = value;
 }
};
Object.defineProperty(foo, "bar", {
 set(value) {
 if (value < 0) {
 throw new Error("Negative value.");
 }
 this.val = value;
 }
});
```
## Options

This rule has no options.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

## Related Rules

## Version

This rule was introduced in ESLint v6.7.0.

# no-sparse-arrays

Disallow sparse arrays

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Sparse arrays contain empty slots, most frequently due to multiple commas being used in an array literal, such as:

```
const items = [,,];
```
While the `items` array in this example has a `length` of 2, there are actually no values in `items[0]` or `items[1]`. The fact that the array literal is valid with only commas inside, coupled with the `length` being set and actual item values not being set, make sparse arrays confusing for many developers. Consider the following:

```
const colors = [ "red",, "blue" ];
```
In this example, the `colors` array has a `length` of 3. But did the developer intend for there to be an empty spot in the middle of the array? Or is it a typo?

The confusion around sparse arrays defined in this manner is enough that it’s recommended to avoid using them unless you are certain that they are useful in your code.

## Rule Details

This rule disallows sparse array literals which have “holes” where commas are not preceded by elements. It does not apply to a trailing comma following the last element.

Examples of **incorrect** code for this rule:

```
/*eslint no-sparse-arrays: "error"*/
const items = [];
const colors = [ "red", "blue" ];
```
Examples of **correct** code for this rule:

```
/*eslint no-sparse-arrays: "error"*/
const items = [];
const arr = new Array(23);
// trailing comma (after the last element) is not a problem
const colors = [ "red", "blue", ];
```
## Options

This rule has no options.

## When Not To Use It

If you want to use sparse arrays, then it is safe to disable this rule.

## Version

This rule was introduced in ESLint v0.4.0.

# no-template-curly-in-string

Disallow template literal placeholder syntax in regular strings

ECMAScript 6 allows programmers to create strings containing variable or expressions using template literals, instead of string concatenation, by writing expressions like `${variable}` between two backtick quotes (`). It can be easy to use the wrong quotes when wanting to use template literals, by writing `"${variable}"`, and end up with the literal value `"${variable}"` instead of a string containing the value of the injected expressions.

## Rule Details

This rule aims to warn when a regular string contains what looks like a template literal placeholder. It will warn when it finds a string containing the template literal placeholder (`${something}`) that uses either `"` or `'` for the quotes.

Examples of **incorrect** code for this rule:

```
/*eslint no-template-curly-in-string: "error"*/
;
;
;
```
Examples of **correct** code for this rule:

```
/*eslint no-template-curly-in-string: "error"*/
`Hello ${name}!`;
`Time: ${12 * 60 * 60 * 1000}`;
templateFunction`Hello ${name}`;
```
## Options

This rule has no options.

## When Not To Use It

This rule should not be used in ES3/5 environments.

## Version

This rule was introduced in ESLint v3.3.0.

# no-this-before-super

Disallow `this`/`super` before calling `super()` in constructors

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

In the constructor of derived classes, if `this`/`super` are used before `super()` calls, it raises a reference error.

This rule checks `this`/`super` keywords in constructors, then reports those that are before `super()`.

## Rule Details

This rule is aimed to flag `this`/`super` keywords before `super()` callings.

Examples of **incorrect** code for this rule:

```
/*eslint no-this-before-super: "error"*/
class A1 extends B {
 constructor() {
 .a = 0;
 super();
 }
}
class A2 extends B {
 constructor() {
 .foo();
 super();
 }
}
class A3 extends B {
 constructor() {
 .foo();
 super();
 }
}
class A4 extends B {
 constructor() {
 super(.foo());
 }
}
```
Examples of **correct** code for this rule:

```
/*eslint no-this-before-super: "error"*/
class A1 {
 constructor() {
 this.a = 0; // OK, this class doesn't have an `extends` clause.
 }
}
class A2 extends B {
 constructor() {
 super();
 this.a = 0; // OK, this is after `super()`.
 }
}
class A3 extends B {
 foo() {
 this.a = 0; // OK. this is not in a constructor.
 }
}
```
## Options

This rule has no options.

## When Not To Use It

If you don’t want to be notified about using `this`/`super` before `super()` in constructors, you can safely disable this rule.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

## Version

This rule was introduced in ESLint v0.24.0.

# no-unassigned-vars

Disallow `let` or `var` variables that are read but never assigned

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

This rule flags `let` or `var` declarations that are never assigned a value but are still read or used in the code. Since these variables will always be `undefined`, their usage is likely a programming mistake.

For example, if you check the value of a `status` variable, but it was never given a value, it will always be `undefined`:

```
let status;
// ...forgot to assign a value to status...
if (status === 'ready') {
 console.log('Ready!');
}
```
## Rule Details

Examples of **incorrect** code for this rule:

```
/*eslint no-unassigned-vars: "error"*/
let ;
if (status === 'ready') {
 console.log('Ready!');
}
let ;
greet(user);
function test() {
 let ;
 return error || "Unknown error";
}
let ;
const { debug } = options || {};
let ;
while (!flag) {
 // Do something...
}
let ;
function init() {
 return config?.enabled;
}
```
In TypeScript:

```
/*eslint no-unassigned-vars: "error"*/
let ;
console.log(value);
```
Examples of **correct** code for this rule:

```
/*eslint no-unassigned-vars: "error"*/
let message = "hello";
console.log(message);
let user;
user = getUser();
console.log(user.name);
let count;
count = 1;
count++;
// Variable is unused (should be reported by `no-unused-vars` only)
let temp;
let error;
if (somethingWentWrong) {
 error = "Something went wrong";
}
console.log(error);
let item;
for (item of items) {
 process(item);
}
let config;
function setup() {
 config = { debug: true };
}
setup();
console.log(config);
let one = undefined;
if (one === two) {
 // Noop
}
```
In TypeScript:

```
/*eslint no-unassigned-vars: "error"*/
declare let value: number | undefined;
console.log(value);
declare module "my-module" {
 let value: string;
 export = value;
}
```
## Options

This rule has no options.

## When Not To Use It

You can disable this rule if your code intentionally uses variables that are declared and used, but are never assigned a value. This might be the case in:

- Legacy codebases where uninitialized variables are used as placeholders.
- Certain TypeScript use cases where variables are declared with a type and intentionally left unassigned (though using `declare`is preferred).

## Related Rules

## Version

This rule was introduced in ESLint v9.27.0.

# no-undef

Disallow the use of undeclared variables unless mentioned in `/*global */` comments

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

This rule can help you locate potential ReferenceErrors resulting from misspellings of variable and parameter names, or accidental implicit globals (for example, from forgetting the `var` keyword in a `for` loop initializer).

## Rule Details

Any reference to an undeclared variable causes a warning, unless the variable is explicitly mentioned in a `/*global ...*/` comment, or specified in the `globals` key in the configuration file. A common use case for these is if you intentionally use globals that are defined elsewhere (e.g. in a script sourced from HTML).

Examples of **incorrect** code for this rule:

```
/*eslint no-undef: "error"*/
const foo = ();
const bar = + 1;
```
Examples of **correct** code for this rule with `global` declaration:

```
/*global someFunction, a*/
/*eslint no-undef: "error"*/
const foo = someFunction();
const bar = a + 1;
```
Note that this rule does not disallow assignments to read-only global variables. See no-global-assign if you also want to disallow those assignments.

This rule also does not disallow redeclarations of global variables. See no-redeclare if you also want to disallow those redeclarations.

## Options

- `typeof`set to true will warn for variables used inside typeof check (Default false).

### typeof

Examples of **correct** code for the default `{ "typeof": false }` option:

```
/*eslint no-undef: "error"*/
if (typeof UndefinedIdentifier === "undefined") {
 // do something ...
}
```
You can use this option if you want to prevent `typeof` check on a variable which has not been declared.

Examples of **incorrect** code for the `{ "typeof": true }` option:

```
/*eslint no-undef: ["error", { "typeof": true }] */
if(typeof === "string"){}
```
Examples of **correct** code for the `{ "typeof": true }` option with `global` declaration:

```
/*global a*/
/*eslint no-undef: ["error", { "typeof": true }] */
if(typeof a === "string"){}
```
## When Not To Use It

If explicit declaration of global variables is not to your taste.

## Compatibility

This rule provides compatibility with treatment of global variables in JSHint and JSLint.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

## Related Rules

## Version

This rule was introduced in ESLint v0.0.9.

# no-unexpected-multiline

Disallow confusing multiline expressions

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Semicolons are usually optional in JavaScript, because of automatic semicolon insertion (ASI). You can require or disallow semicolons with the semi rule.

The rules for ASI are relatively straightforward: As once described by Isaac Schlueter, a newline character always ends a statement, just like a semicolon, **except** where one of the following is true:

- The statement has an unclosed paren, array literal, or object literal or ends in some other way that is not a valid way to end a statement. (For instance, ending with `.`or`,`.)
- The line is `--`or`++`(in which case it will decrement/increment the next token.)
- It is a `for()`,`while()`,`do`,`if()`, or`else`, and there is no`{`
- The next line starts with `[`,`(`,`+`,`*`,`/`,`-`,`,`,`.`, or some other binary operator that can only be found between two tokens in a single expression.

In the exceptions where a newline does **not** end a statement, a typing mistake to omit a semicolon causes two unrelated consecutive lines to be interpreted as one expression. Especially for a coding style without semicolons, readers might overlook the mistake. Although syntactically correct, the code might throw exceptions when it is executed.

## Rule Details

This rule disallows confusing multiline expressions where a newline looks like it is ending a statement, but is not.

Examples of **incorrect** code for this rule:

```
/*eslint no-unexpected-multiline: "error"*/
const foo = bar
1 || 2).baz();
const hello = 'world'
1, 2, 3].forEach(addNumber);
const x = function() {}
hello`
const y = function() {}
y
hello`
const z = foo
regex/g.test(bar)
```
Examples of **correct** code for this rule:

```
/*eslint no-unexpected-multiline: "error"*/
const foo = bar;
(1 || 2).baz();
const baz = bar
;(1 || 2).baz()
const hello = 'world';
[1, 2, 3].forEach(addNumber);
const hi = 'world'
void [1, 2, 3].forEach(addNumber);
const x = function() {};
`hello`
const tag = function() {}
tag `hello`
```
## Options

This rule has no options.

## When Not To Use It

You can turn this rule off if you are confident that you will not accidentally introduce code like this.

Note that the patterns considered problems are **not** flagged by the semi rule.

## Version

This rule was introduced in ESLint v0.24.0.

# no-unmodified-loop-condition

Disallow unmodified loop conditions

Variables in a loop condition often are modified in the loop. If not, it’s possibly a mistake.

```
while (node) {
 doSomething(node);
}
```
```
while (node) {
 doSomething(node);
 node = node.parent;
}
```
## Rule Details

This rule finds references which are inside of loop conditions, then checks the variables of those references are modified in the loop.

If a reference is inside of a binary expression or a ternary expression, this rule checks the result of
the expression instead.
If a reference is inside of a dynamic expression (e.g. `CallExpression`,
`YieldExpression`, …), this rule ignores it.

Examples of **incorrect** code for this rule:

```
/*eslint no-unmodified-loop-condition: "error"*/
let node = something;
while () {
 doSomething(node);
}
node = other;
for (let j = 0; < 5;) {
 doSomething(j);
}
while ( !== root) {
 doSomething(node);
}
```
Examples of **correct** code for this rule:

```
/*eslint no-unmodified-loop-condition: "error"*/
while (node) {
 doSomething(node);
 node = node.parent;
}
for (let j = 0; j < items.length; ++j) {
 doSomething(items[j]);
}
// OK, the result of this binary expression is changed in this loop.
while (node !== root) {
 doSomething(node);
 node = node.parent;
}
// OK, the result of this ternary expression is changed in this loop.
while (node ? A : B) {
 doSomething(node);
 node = node.parent;
}
// A property might be a getter which has side effect...
// Or "doSomething" can modify "obj.foo".
while (obj.foo) {
 doSomething(obj);
}
// A function call can return various values.
while (check(obj)) {
 doSomething(obj);
}
```
## Options

This rule has no options.

## When Not To Use It

If you don’t want to notified about references inside of loop conditions, then it’s safe to disable this rule.

## Version

This rule was introduced in ESLint v2.0.0-alpha-2.

# no-unreachable

Disallow unreachable code after `return`, `throw`, `continue`, and `break` statements

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

Because the `return`, `throw`, `continue`, and `break` statements unconditionally exit a block of code, any statements after them cannot be executed. Unreachable statements are usually a mistake.

```
function fn() {
 x = 1;
 return x;
 x = 3; // this will never execute
}
```
Another kind of mistake is defining instance fields in a subclass whose constructor doesn’t call `super()`. Instance fields of a subclass are only added to the instance after `super()`. If there are no `super()` calls, their definitions are never applied and therefore are unreachable code.

```
class C extends B {
 #x; // this will never be added to instances
 constructor() {
 return {};
 }
}
```
## Rule Details

This rule disallows unreachable code after `return`, `throw`, `continue`, and `break` statements. This rule also flags definitions of instance fields in subclasses whose constructors don’t have `super()` calls.

Examples of **incorrect** code for this rule:

```
/*eslint no-unreachable: "error"*/
function foo() {
 return true;

}
function bar() {
 throw new Error("Oops!");

}
while(value) {
 break;

}
throw new Error("Oops!");
function baz() {
 if (Math.random() < 0.5) {
 return;
 } else {
 throw new Error();
 }

}
```
Examples of **correct** code for this rule, because of JavaScript function and variable hoisting:

```
/*eslint no-unreachable: "error"*/
function foo() {
 return bar();
 function bar() {
 return 1;
 }
}
function bar() {
 return x;
 var x;
}
switch (foo) {
 case 1:
 break;
 var x;
}
```
Examples of additional **incorrect** code for this rule:

```
/*eslint no-unreachable: "error"*/
class C extends B {
 // unreachable
 constructor() {
 return {};
 }
}
```
Examples of additional **correct** code for this rule:

```
/*eslint no-unreachable: "error"*/
class D extends B {
 #x;
 #y = 1;
 a;
 b = 1;
 constructor() {
 super();
 }
}
class E extends B {
 #x;
 #y = 1;
 a;
 b = 1;
 // implicit constructor always calls `super()`
}
class F extends B {
 static #x;
 static #y = 1;
 static a;
 static b = 1;
 constructor() {
 return {};
 }
}
```
## Options

This rule has no options.

## Handled by TypeScript

It is safe to disable this rule when using TypeScript because TypeScript's compiler enforces this check.

TypeScript must be configured with `allowUnreachableCode: false` for it to consider unreachable code an error.

## Version

This rule was introduced in ESLint v0.0.6.

# no-unreachable-loop

Disallow loops with a body that allows only one iteration

A loop that can never reach the second iteration is a possible error in the code.

```
for (let i = 0; i < arr.length; i++) {
 if (arr[i].name === myName) {
 doSomething(arr[i]);
 // break was supposed to be here
 }
 break;
}
```
In rare cases where only one iteration (or at most one iteration) is intended behavior, the code should be refactored to use `if` conditionals instead of `while`, `do-while` and `for` loops. It’s considered a best practice to avoid using loop constructs for such cases.

## Rule Details

This rule aims to detect and disallow loops that can have at most one iteration, by performing static code path analysis on loop bodies.

In particular, this rule will disallow a loop with a body that exits the loop in all code paths. If all code paths in the loop’s body will end with either a `break`, `return` or a `throw` statement, the second iteration of such loop is certainly unreachable, regardless of the loop’s condition.

This rule checks `while`, `do-while`, `for`, `for-in` and `for-of` loops. You can optionally disable checks for each of these constructs.

Examples of **incorrect** code for this rule:

```
/*eslint no-unreachable-loop: "error"*/
function verifyList(head) {
 let item = head;

}
function findSomething(arr) {

}
```
Examples of **correct** code for this rule:

```
/*eslint no-unreachable-loop: "error"*/
while (foo) {
 doSomething(foo);
 foo = foo.parent;
}
function verifyList(head) {
 let item = head;
 do {
 if (verify(item)) {
 item = item.next;
 } else {
 return false;
 }
 } while (item);
 return true;
}
function findSomething(arr) {
 for (let i = 0; i < arr.length; i++) {
 if (isSomething(arr[i])) {
 return arr[i];
 }
 }
 throw new Error("Doesn't exist.");
}
for (key in obj) {
 if (key.startsWith("_")) {
 continue;
 }
 firstKey = key;
 firstValue = obj[key];
 break;
}
for (foo of bar) {
 if (foo.id === id) {
 doSomething(foo);
 break;
 }
}
```
Please note that this rule is not designed to check loop conditions, and will not warn in cases such as the following examples.

Examples of additional **correct** code for this rule:

```
/*eslint no-unreachable-loop: "error"*/
do {
 doSomething();
} while (false)
for (let i = 0; i < 1; i++) {
 doSomething(i);
}
for (const a of [1]) {
 doSomething(a);
}
```
## Options

This rule has an object option, with one option:

- `"ignore"`- an optional array of loop types that will be ignored by this rule.

### ignore

You can specify up to 5 different elements in the `"ignore"` array:

- `"WhileStatement"`- to ignore all- `while`loops.
- `"DoWhileStatement"`- to ignore all- `do-while`loops.
- `"ForStatement"`- to ignore all- `for`loops (does not apply to- `for-in`and- `for-of`loops).
- `"ForInStatement"`- to ignore all- `for-in`loops.
- `"ForOfStatement"`- to ignore all- `for-of`loops.

Examples of **correct** code for this rule with the `"ignore"` option:

```
/*eslint no-unreachable-loop: ["error", { "ignore": ["ForInStatement", "ForOfStatement"] }]*/
for (let key in obj) {
 hasEnumerableProperties = true;
 break;
}
for (const a of b) break;
```
## Known Limitations

Static code path analysis, in general, does not evaluate conditions. Due to this fact, this rule might miss reporting cases such as the following:

```
for (let i = 0; i < 10; i++) {
 doSomething(i);
 if (true) {
 break;
 }
}
```
## Related Rules

## Version

This rule was introduced in ESLint v7.3.0.

# no-unsafe-finally

Disallow control flow statements in `finally` blocks

 Using the `recommended` config from `@eslint/js` in a configuration file
 enables this rule

JavaScript suspends the control flow statements of `try` and `catch` blocks until the execution of `finally` block finishes. So, when `return`, `throw`, `break`, or `continue` is used in `finally`, control flow statements inside `try` and `catch` are overwritten, which is considered as unexpected behavior. Such as:

```
// We expect this function to return 1;
(() => {
 try {
 return 1; // 1 is returned but suspended until finally block ends
 } catch(err) {
 return 2;
 } finally {
 return 3; // 3 is returned before 1, which we did not expect
 }
})();
// > 3
```
```
// We expect this function to throw an error, then return
(() => {
 try {
 throw new Error("Try"); // error is thrown but suspended until finally block ends
 } finally {
 return 3; // 3 is returned before the error is thrown, which we did not expect
 }
})();
// > 3
```
```
// We expect this function to throw Try(...) error from the catch block
(() => {
 try {
 throw new Error("Try")
 } catch(err) {
 throw err; // The error thrown from try block is caught and rethrown
 } finally {
 throw new Error("Finally"); // Finally(...) is thrown, which we did not expect
 }
})();
// > Uncaught Error: Finally(...)
```
```
// We expect this function to return 0 from try block.
(() => {
 label: try {
 return 0; // 0 is returned but suspended until finally block ends
 } finally {
 break label; // It breaks out the try-finally block, before 0 is returned.
 }
 return 1;
})();
// > 1
```
## Rule Details

This rule disallows `return`, `throw`, `break`, and `continue` statements inside `finally` blocks. It allows indirect usages, such as in `function` or `class` definitions.

Examples of **incorrect** code for this rule:

```
/*eslint no-unsafe-finally: "error"*/
let foo = function() {
 try {
 return 1;
 } catch(err) {
 return 2;
 } finally {

 }
};
```
```
/*eslint no-unsafe-finally: "error"*/
let foo = function() {
 try {
 return 1;
 } catch(err) {
 return 2;
 } finally {

 }
};
```
Examples of **correct** code for this rule:

```
/*eslint no-unsafe-finally: "error"*/
let foo = function() {
 try {
 return 1;
 } catch(err) {
 return 2;
 } finally {
 console.log("hola!");
 }
};
```
```
/*eslint no-unsafe-finally: "error"*/
let foo = function() {
 try {
 return 1;
 } catch(err) {
 return 2;
 } finally {
 let a = function() {
 return "hola!";
 }
 }
};
```
```
/*eslint no-unsafe-finally: "error"*/
let foo = function(a) {
 try {
 return 1;
 } catch(err) {
 return 2;
 } finally {
 switch(a) {
 case 1: {
 console.log("hola!")
 break;
 }
 }
 }
};
```
## Options

This rule has no options.

## When Not To Use It

If you want to allow control flow operations in `finally` blocks, you can turn this rule off.

Note that using `return` with `try`/`finally` in generator functions can be used to override the value returned.
If you want to use this feature of generator functions, you can disable this rule using an `eslint-disable` comment.

## Version

This rule was introduced in ESLint v2.9.0.
