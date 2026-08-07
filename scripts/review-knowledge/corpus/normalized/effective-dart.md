# Effective Dart

Best practices for building consistent, maintainable, and efficient Dart libraries.

Over the past several years, we've written a ton of Dart code and learned a lot about what works well and what doesn't. We're sharing this with you so you can write consistent, robust, fast code too. There are two overarching themes:

-
 **Be consistent.**When it comes to things like formatting, and casing, arguments about which is better are subjective and impossible to resolve. What we do know is that being*consistent*is objectively helpful.If two pieces of code look different it should be because they *are*different in some meaningful way. When a bit of code stands out and catches your eye, it should do so for a useful reason.
-
 **Be brief.**Dart was designed to be familiar, so it inherits many of the same statements and expressions as C, Java, JavaScript and other languages. But we created Dart because there is a lot of room to improve on what those languages offer. We added a bunch of features, from string interpolation to initializing formals, to help you express your intent more simply and easily.If there are multiple ways to say something, you should generally pick the most concise one. This is not to say you should code golf yourself into cramming a whole program into a single line. The goal is code that is *economical*, not*dense*.

## The guides

#We split the guidelines into a few separate pages for easy digestion:

-
 **Style Guide**– This defines the rules for laying out and organizing code, or at least the parts that`dart format`doesn't handle for you. The style guide also specifies how identifiers are formatted:`camelCase`,`using_underscores`, etc.
-
 **Documentation Guide**– This tells you everything you need to know about what goes inside comments. Both doc comments and regular, run-of-the-mill code comments.
-
 **Usage Guide**– This teaches you how to make the best use of language features to implement behavior. If it's in a statement or expression, it's covered here.
-
 **Design Guide**– This is the softest guide, but the one with the widest scope. It covers what we've learned about designing consistent, usable APIs for libraries. If it's in a type signature or declaration, this goes over it.

For links to all the guidelines, see the summary.

## How to read the guides

#Each guide is broken into a few sections. Sections contain a list of guidelines. Each guideline starts with one of these words:

-
 **DO**guidelines describe practices that should always be followed. There will almost never be a valid reason to stray from them.
-
 **DON'T**guidelines are the converse: things that are almost never a good idea. Hopefully, we don't have as many of these as other languages do because we have less historical baggage.
-
 **PREFER**guidelines are practices that you*should*follow. However, there may be circumstances where it makes sense to do otherwise. Just make sure you understand the full implications of ignoring the guideline when you do.
-
 **AVOID**guidelines are the dual to "prefer": stuff you shouldn't do but where there may be good reasons to on rare occasions.
-
 **CONSIDER**guidelines are practices that you might or might not want to follow, depending on circumstances, precedents, and your own preference.

 Some guidelines describe an **exception** where the rule does *not* apply. When
 listed, the exceptions may not be exhaustive—you might still need to use
 your judgement on other cases.

This sounds like the police are going to beat down your door if you don't have your laces tied correctly. Things aren't that bad. Most of the guidelines here are common sense and we're all reasonable people. The goal, as always, is nice, readable and maintainable code.

The Dart analyzer provides a linter to help you write good, consistent code that follows these and other guidelines. If one or more linter rules exist that can help you follow a guideline then the guideline links to those rules. The links use the following format:

Linter rule: unnecessary_getters_setters

To learn how to use the linter, see Enabling linter rules and the list of linter rules.

## Glossary

#To keep the guidelines brief, we use a few shorthand terms to refer to different Dart constructs.

-
 A **library member**is a top-level field, getter, setter, or function. Basically, anything at the top level that isn't a type.
-
 A **class member**is a constructor, field, getter, setter, function, or operator declared inside a class. Class members can be instance or static, abstract or concrete.
- A - **member**is either a library member or a class member.
-
 A **variable**, when used generally, refers to top-level variables, parameters, and local variables. It doesn't include static or instance fields.
- A - **type**is any named type declaration: a class, typedef, or enum.
-
 A **property**is a top-level variable, getter (inside a class or at the top level, instance or static), setter (same), or field (instance or static). Roughly any "field-like" named construct.

## Summary of all rules

#### Style

#**Identifiers**

-
 DO name types using `UpperCamelCase`.
-
 DO name extensions using `UpperCamelCase`.
-
 DO name packages, directories, and source files using `lowercase_with_underscores`.
-
 DO name import prefixes using `lowercase_with_underscores`.
-
 DO name other identifiers using `lowerCamelCase`.
-
 PREFER using `lowerCamelCase`for constant names.
- DO capitalize acronyms and abbreviations longer than two letters like words.
- PREFER using wildcards for unused callback parameters.
- DON'T use a leading underscore for identifiers that aren't private.
- DON'T use prefix letters.
- DON'T explicitly name libraries.

**Ordering**

-
 DO place `dart:`imports before other imports.
-
 DO place `package:`imports before relative imports.
- DO specify exports in a separate section after all imports.
- DO sort sections alphabetically.

**Formatting**

### Documentation

#**Comments**

**Doc comments**

-
 DO use `///`doc comments to document members and types.
- PREFER writing doc comments for public APIs.
- CONSIDER writing a library-level doc comment.
- CONSIDER writing doc comments for private APIs.
- DO start doc comments with a single-sentence summary.
- DO separate the first sentence of a doc comment into its own paragraph.
- AVOID redundancy with the surrounding context.
- PREFER starting comments of a function or method with third-person verbs if its main purpose is a side effect.
- PREFER starting a non-boolean variable or property comment with a noun phrase.
- PREFER starting a boolean variable or property comment with "Whether" followed by a noun or gerund phrase.
- PREFER a noun phrase or non-imperative verb phrase for a function or method if returning a value is its primary purpose.
- DON'T write documentation for both the getter and setter of a property.
- PREFER starting library or type comments with noun phrases.
- CONSIDER including code samples in doc comments.
- DO use square brackets in doc comments to refer to in-scope identifiers.
- DO use prose to explain parameters, return values, and exceptions.
- DO put doc comments before metadata annotations.

**Markdown**

- AVOID using markdown excessively.
- AVOID using HTML for formatting.
- PREFER backtick fences for code blocks.

**Writing**

### Usage

#**Libraries**

-
 DO use strings in `part of`directives.
-
 DON'T import libraries that are inside the `src`directory of another package.
-
 DON'T allow an import path to reach into or out of `lib`.
- PREFER relative import paths.

**Null**

-
 DON'T explicitly initialize variables to `null`.
-
 DON'T use an explicit default value of `null`.
-
 DON'T use `true`or`false`in equality operations.
-
 AVOID `late`variables if you need to check whether they are initialized.
- CONSIDER type promotion or null-check patterns for using nullable types.

**Strings**

- DO use adjacent strings to concatenate string literals.
- PREFER using interpolation to compose strings and values.
- AVOID using curly braces in interpolation when not needed.

**Collections**

- DO use collection literals when possible.
-
 DON'T use `.length`to see if a collection is empty.
-
 AVOID using `Iterable.forEach()`with a function literal.
-
 DON'T use `List.from()`unless you intend to change the type of the result.
-
 DO use `whereType()`to filter a collection by type.
-
 DON'T use `cast()`when a nearby operation will do.
- AVOID using `cast()`.

**Functions**

- DO use a function declaration to bind a function to a name.
- DON'T create a lambda when a tear-off will do.

**Variables**

-
 DO follow a consistent rule for `var`and`final`on local variables.
- AVOID storing what you can calculate.

**Members**

- DON'T wrap a field in a getter and setter unnecessarily.
-
 PREFER using a `final`field to make a read-only property.
-
 CONSIDER using `=>`for simple members.
-
 DON'T use `this.`except to redirect to a named constructor or to avoid shadowing.
- DO initialize fields at their declaration when possible.

**Constructors**

- DO use initializing formals when possible.
-
 DON'T use `late`when a constructor initializer list will do.
-
 DO use `;`instead of`{}`for empty constructor bodies.
- PREFER using concise constructor syntax.
- DON'T use `new`.
-
 DON'T use `const`redundantly.

**Error handling**

-
 AVOID catches without `on`clauses.
-
 DON'T discard errors from catches without `on`clauses.
-
 DO throw objects that implement `Error`only for programmatic errors.
-
 DON'T explicitly catch `Error`or types that implement it.
-
 DO use `rethrow`to rethrow a caught exception.

**Asynchrony**

### Design

#**Names**

- DO use terms consistently.
- AVOID abbreviations.
- PREFER putting the most descriptive noun last.
- CONSIDER making the code read like a sentence.
- PREFER a noun phrase for a non-boolean property or variable.
- PREFER a non-imperative verb phrase for a boolean property or variable.
-
 CONSIDER omitting the verb for a named boolean *parameter*.
- PREFER the "positive" name for a boolean property or variable.
- PREFER an imperative verb phrase for a function or method whose main purpose is a side effect.
- PREFER a noun phrase or non-imperative verb phrase for a function or method if returning a value is its primary purpose.
- CONSIDER an imperative verb phrase for a function or method if you want to draw attention to the work it performs.
-
 AVOID starting a function or method name with `get`.
-
 PREFER naming a method `to___()`if it copies the object's state to a new object.
-
 PREFER naming a method `as___()`if it returns a different representation backed by the original object.
- AVOID describing the parameters in the function's or method's name.
- DO follow existing mnemonic conventions when naming type parameters.

**Libraries**

**Classes and mixins**

- AVOID defining a one-member abstract class when a simple function will do.
- AVOID defining a class that contains only static members.
- AVOID extending a class that isn't intended to be subclassed.
- DO use class modifiers to control if your class can be extended.
- AVOID implementing a class that isn't intended to be an interface.
- DO use class modifiers to control if your class can be an interface.
-
 PREFER defining a pure `mixin`or pure`class`to a`mixin class`.

**Constructors**

**Members**

-
 PREFER making fields and top-level variables `final`.
- DO use getters for operations that conceptually access properties.
- DO use setters for operations that conceptually change properties.
- DON'T define a setter without a corresponding getter.
- AVOID using runtime type tests to fake overloading.
-
 AVOID public `late final`fields without initializers.
-
 AVOID returning nullable `Future`,`Stream`, and collection types.
-
 AVOID returning `this`from methods just to enable a fluent interface.

**Types**

- DO type annotate variables without initializers.
- DO type annotate fields and top-level variables if the type isn't obvious.
- DON'T redundantly type annotate initialized local variables.
- DO annotate return types on function declarations.
- DO annotate parameter types on function declarations.
- DON'T annotate inferred parameter types on function expressions.
- DON'T type annotate initializing formals.
- DO write type arguments on generic invocations that aren't inferred.
- DON'T write type arguments on generic invocations that are inferred.
- AVOID writing incomplete generic types.
-
 DO annotate with `dynamic`instead of letting inference fail.
- PREFER signatures in function type annotations.
- DON'T specify a return type for a setter.
- DON'T use the legacy typedef syntax.
- PREFER inline function types over typedefs.
- PREFER using function type syntax for parameters.
-
 AVOID using `dynamic`unless you want to disable static checking.
-
 DO use `Future<void>`as the return type of asynchronous members that do not produce values.
-
 AVOID using `FutureOr<T>`as a return type.

**Parameters**

- AVOID positional boolean parameters.
- AVOID optional positional parameters if the user may want to omit earlier parameters.
- AVOID mandatory parameters that accept a special "no argument" value.
- DO use inclusive start and exclusive end parameters to accept a range.

**Equality**

Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Page last updated on 2025-09-04. View source or report an issue.

# Effective Dart: Style

Formatting and naming rules for consistent, readable code.

 A surprisingly important part of good code is good style. Consistent naming,
 ordering, and formatting helps code that *is* the same *look* the same. It takes
 advantage of the powerful pattern-matching hardware most of us have in our
 ocular systems. If we use a consistent style across the entire Dart ecosystem,
 it makes it easier for all of us to learn from and contribute to each others'
 code.

## Identifiers

#Identifiers come in three flavors in Dart.

-
 `UpperCamelCase`names capitalize the first letter of each word, including the first.
-
 `lowerCamelCase`names capitalize the first letter of each word,*except*the first which is always lowercase, even if it's an acronym.
-
 `lowercase_with_underscores`names use only lowercase letters, even for acronyms, and separate words with`_`.

### DO name types using `UpperCamelCase`

 #
 Linter rule: camel_case_types

Classes, enum types, typedefs, and type parameters should capitalize the first letter of each word (including the first word), and use no separators.

```
class SliderMenu {
 ...
}
class HttpRequest {
 ...
}
typedef Predicate<T> = bool Function(T value);
```
This even includes classes intended to be used in metadata annotations.

```
class Foo {
 const Foo([Object? arg]);
}
@Foo(anArg)
class A {
 ...
}
@Foo()
class B {
 ...
}
```
 If the annotation class's constructor takes no parameters, you might want to
 create a separate `lowerCamelCase` constant for it.

```
const foo = Foo();
@foo
class C {
 ...
}
```
### DO name extensions using `UpperCamelCase`

 #
 Linter rule: camel_case_extensions

Like types, extensions should capitalize the first letter of each word (including the first word), and use no separators.

```
extension MyFancyList<T> on List<T> {
 ...
}
extension SmartIterable<T> on Iterable<T> {
 ...
}
```
### DO name packages, directories, and source files using `lowercase_with_underscores`

 #
 Linter rules: file_names, package_names

Some file systems are not case-sensitive, so many projects require filenames to be all lowercase. Using a separating character allows names to still be readable in that form. Using underscores as the separator ensures that the name is still a valid Dart identifier, which may be helpful if the language later supports symbolic imports.

```
my_package
└─ lib
 └─ file_system.dart
 └─ slider_menu.dart
```
```
mypackage
└─ lib
 └─ file-system.dart
 └─ SliderMenu.dart
```
### DO name import prefixes using `lowercase_with_underscores`

 #
 Linter rule: library_prefixes

```
import 'dart:math' as math;
import 'package:angular_components/angular_components.dart' as angular_components;
import 'package:js/js.dart' as js;
```
```
import 'dart:math' as Math;
import 'package:angular_components/angular_components.dart' as angularComponents;
import 'package:js/js.dart' as JS;
```
### DO name other identifiers using `lowerCamelCase`

 #
 Linter rule: non_constant_identifier_names

 Class members, top-level definitions, variables, parameters, and named
 parameters should capitalize the first letter of each word *except* the first
 word, and use no separators.

```
var count = 3;
HttpRequest httpRequest;
void align(bool clearItems) {
 // ...
}
```
### PREFER using `lowerCamelCase` for constant names

 #
 Linter rule: constant_identifier_names

In new code, use `lowerCamelCase` for constant variables, including enum values.

```
const pi = 3.14;
const defaultTimeout = 1000;
final urlScheme = RegExp('^([a-z]+):');
class Dice {
 static final numberGenerator = Random();
}
```
```
const PI = 3.14;
const DefaultTimeout = 1000;
final URL_SCHEME = RegExp('^([a-z]+):');
class Dice {
 static final NUMBER_GENERATOR = Random();
}
```
 You may use `SCREAMING_CAPS` for consistency with existing code,
 as in the following cases:

- When adding code to a file or library that already uses `SCREAMING_CAPS`.
- When generating Dart code that's parallel to Java code—for example, in enumerated types generated from protobufs.

### DO capitalize acronyms and abbreviations longer than two letters like words

#
 Capitalized acronyms can be hard to read,
 and multiple adjacent acronyms can lead to ambiguous names.
 For example, given an identifier `HTTPSFTP`,
 the reader can't tell if it refers to `HTTPS` `FTP` or `HTTP`
 `SFTP`.
 To avoid this,
 capitalize most acronyms and abbreviations like regular words.
 This identifier would be `HttpsFtp` if referring to the former
 or `HttpSftp` for the latter.

Two-letter abbreviations and acronyms are the exception. If both letters are capitalized in English, then they should both stay capitalized when used in an identifier. Otherwise, capitalize it like a word.

```
// Longer than two letters, so always like a word:
Http // "hypertext transfer protocol"
Nasa // "national aeronautics and space administration"
Uri // "uniform resource identifier"
Esq // "esquire"
Ave // "avenue"
// Two letters, capitalized in English, so capitalized in an identifier:
ID // "identifier"
TV // "television"
UI // "user interface"
// Two letters, not capitalized in English, so like a word in an identifier:
Mr // "mister"
St // "street"
Rd // "road"
```
```
HTTP // "hypertext transfer protocol"
NASA // "national aeronautics and space administration"
URI // "uniform resource identifier"
esq // "esquire"
ave // "avenue"
Id // "identifier"
Tv // "television"
Ui // "user interface"
MR // "mister"
ST // "street"
RD // "road"
```
 When any form of abbreviation comes at the beginning
 of a `lowerCamelCase` identifier, the abbreviation should be all lowercase:

```
var httpConnection = connect();
var tvSet = Television();
var mrRogers = 'hello, neighbor';
```

### PREFER using wildcards for unused callback parameters

#
 Sometimes the type signature of a callback function requires a parameter,
 but the callback implementation doesn't *use* the parameter.
 In this case, it's idiomatic to name the unused parameter `_`,
 which declares a wildcard variable
 that is non-binding.

```
futureOfVoid.then((_) {
 print('Operation complete.');
});
```
 Because wildcard variables are non-binding,
 you can name multiple unused parameters `_`.

```
.onError((_, _) {
 print('Operation failed.');
});
```
 This guideline is only for functions that are both *anonymous and local*.
 These functions are usually used immediately in a context where it's
 clear what the unused parameter represents.
 In contrast, top-level functions and method declarations don't have that context,
 so their parameters must be named so that it's clear what each parameter is for,
 even if it isn't used.

### DON'T use a leading underscore for identifiers that aren't private

#Dart uses a leading underscore in an identifier to mark members and top-level declarations as private. This trains users to associate a leading underscore with one of those kinds of declarations. They see "_" and think "private".

There is no concept of "private" for local variables, parameters, local functions, or library prefixes. When one of those has a name that starts with an underscore, it sends a confusing signal to the reader. To avoid that, don't use leading underscores in those names.

### DON'T use prefix letters

#Hungarian notation and other schemes arose in the time of BCPL, when the compiler didn't do much to help you understand your code. Because Dart can tell you the type, scope, mutability, and other properties of your declarations, there's no reason to encode those properties in identifier names.

```
defaultTimeout
```
```
kDefaultTimeout
```
### DON'T explicitly name libraries

#
 Appending a name to the `library` directive is technically possible,
 but is a legacy feature and discouraged.

Dart generates a unique tag for each library based on its path and filename. Naming libraries overrides this generated URI. Without the URI, it can be harder for tools to find the main library file in question.

```
library my_library;
```
```
/// A really great test library.
@TestOn('browser')
library;
```
## Ordering

#To keep the preamble of your file tidy, we have a prescribed order that directives should appear in. Each "section" should be separated by a blank line.

A single linter rule handles all the ordering guidelines: directives_ordering.

### DO place `dart:` imports before other imports

 #
 Linter rule: directives_ordering

```
import 'dart:async';
import 'dart:collection';
import 'package:bar/bar.dart';
import 'package:foo/foo.dart';
```
### DO place `package:` imports before relative imports

 #
 Linter rule: directives_ordering

```
import 'package:bar/bar.dart';
import 'package:foo/foo.dart';
import 'util.dart';
```
### DO specify exports in a separate section after all imports

#Linter rule: directives_ordering

```
import 'src/error.dart';
import 'src/foo_bar.dart';
export 'src/error.dart';
```
```
import 'src/error.dart';
export 'src/error.dart';
import 'src/foo_bar.dart';
```
### DO sort sections alphabetically

#Linter rule: directives_ordering

```
import 'package:bar/bar.dart';
import 'package:foo/foo.dart';
import 'foo.dart';
import 'foo/foo.dart';
```
```
import 'package:foo/foo.dart';
import 'package:bar/bar.dart';
import 'foo/foo.dart';
import 'foo.dart';
```
## Formatting

#
 Like many languages, Dart ignores whitespace. However, *humans* don't. Having a
 consistent whitespace style helps ensure that human readers see code the same
 way the compiler does.

### DO format your code using `dart format`

 #

 Formatting is tedious work and is particularly time-consuming during
 refactoring. Fortunately, you don't have to worry about it. We provide a
 sophisticated automated code formatter called `dart format`
 that does it for
 you. The official whitespace-handling rules for Dart are
 *whatever dart format produces*. The formatter FAQ
 can provide more insight
 into the style choices it enforces.

 The remaining formatting guidelines are for the few things `dart format` cannot
 fix for you.

### CONSIDER changing your code to make it more formatter-friendly

#The formatter does the best it can with whatever code you throw at it, but it can't work miracles. If your code has particularly long identifiers, deeply nested expressions, a mixture of different kinds of operators, etc. the formatted output may still be hard to read.

 When that happens, reorganize or simplify your code. Consider shortening a local
 variable name or hoisting out an expression into a new local variable. In other
 words, make the same kinds of modifications that you'd make if you were
 formatting the code by hand and trying to make it more readable. Think of
 `dart format` as a partnership where you work together, sometimes iteratively,
 to produce beautiful code.

### PREFER lines 80 characters or fewer

#Linter rule: lines_longer_than_80_chars

Readability studies show that long lines of text are harder to read because your eye has to travel farther when moving to the beginning of the next line. This is why newspapers and magazines use multiple columns of text.

 If you really find yourself wanting lines longer than 80 characters, our
 experience is that your code is likely too verbose and could be a little more
 compact. The main offender is usually `VeryLongCamelCaseClassNames`. Ask
 yourself, "Does each word in that type name tell me something critical or
 prevent a name collision?" If not, consider omitting it.

 Note that `dart format` defaults to 80 characters or fewer, though you can
 configure the default.
 It does not split long string literals to fit in 80 columns,
 so you have to do that manually.

 **Exception:** When a URI or file path occurs in a comment or string (usually in
 an import or export), it may remain whole even if it causes the line to go over
 80 characters. This makes it easier to search source files for a path.

 **Exception:** Multi-line strings can contain lines longer than 80 characters
 because newlines are significant inside the string and splitting the lines into
 shorter ones can alter the program.

### DO use curly braces for all flow control statements

#Linter rule: curly_braces_in_flow_control_structures

Doing so avoids the dangling else problem.

```
if (isWeekDay) {
 print('Bike to work!');
} else {
 print('Go dancing or read a book!');
}
```
 **Exception:** When you have an `if` statement with no `else` clause and the
 whole `if` statement fits on one line, you can omit the braces if you prefer:

```
if (arg == null) return defaultValue;
```
If the body wraps to the next line, though, use braces:

```
if (overflowChars != other.overflowChars) {
 return overflowChars < other.overflowChars;
}
```
```
if (overflowChars != other.overflowChars)
 return overflowChars < other.overflowChars;
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Page last updated on 2026-05-12. View source or report an issue.

# Effective Dart: Documentation

Clear, helpful comments and documentation.

It's easy to think your code is obvious today without realizing how much you rely on context already in your head. People new to your code, and even your forgetful future self won't have that context. A concise, accurate comment only takes a few seconds to write but can save one of those people hours of time.

 We all know code should be self-documenting and not all comments are helpful.
 But the reality is that most of us don't write as many comments as we should.
 It's like exercise: you technically *can* do too much, but it's a lot more
 likely that you're doing too little. Try to step it up.

## Comments

#The following tips apply to comments that you don't want included in the generated documentation.

### DO format comments like sentences

#```
// Not if anything comes before it.
if (_chunks.isNotEmpty) return false;
```
Capitalize the first word unless it's a case-sensitive identifier. End it with a period (or "!" or "?", I suppose). This is true for all comments: doc comments, inline stuff, even TODOs. Even if it's a sentence fragment.

### DON'T use block comments for documentation

#```
void greet(String name) {
 // Assume we have a valid name.
 print('Hi, $name!');
}
```
```
void greet(String name) {
 /* Assume we have a valid name. */
 print('Hi, $name!');
}
```
 You can use a block comment (`/* ... */`) to temporarily comment out a section
 of code, but all other comments should use `//`.

## Doc comments

#
 Doc comments are especially handy because `dart doc` parses them
 and generates beautiful doc pages from them.
 A doc comment is any comment that appears before a declaration
 and uses the special `///` syntax that `dart doc` looks for.

### DO use `///` doc comments to document members and types

 #
 Linter rule: slash_for_doc_comments

 Using a doc comment instead of a regular comment enables
 `dart doc` to find it
 and generate documentation for it.

```
/// The number of characters in this chunk when unsplit.
int get length => ...
```
```
// The number of characters in this chunk when unsplit.
int get length => ...
```
 For historical reasons, `dart doc` supports two syntaxes of doc comments: `///`

 ("C# style") and `/** ... */` ("JavaDoc style"). We prefer `///`
 because it's
 more compact. `/**` and `*/` add two content-free lines to a multiline doc
 comment. The `///` syntax is also easier to read in some situations, such as
 when a doc comment contains a bulleted list that uses `*` to mark list items.

If you stumble onto code that still uses the JavaDoc style, consider cleaning it up.

### PREFER writing doc comments for public APIs

#Linter rule: public_member_api_docs

You don't have to document every single library, top-level variable, type, and member, but you should document most of them.

### CONSIDER writing a library-level doc comment

#
 Unlike languages like Java where the class is the only unit of program
 organization, in Dart, a library is itself an entity that users work with
 directly, import, and think about. That makes the `library` directive a great
 place for documentation that introduces the reader to the main concepts and
 functionality provided within. Consider including:

- A single-sentence summary of what the library is for.
- Explanations of terminology used throughout the library.
- A couple of complete code samples that walk through using the API.
- Links to the most important or most commonly used classes and functions.
- Links to external references on the domain the library is concerned with.

 To document a library, place a doc comment before
 the `library` directive and any annotations that might be attached
 at the start of the file.

```
/// A really great test library.
@TestOn('browser')
library;
```
### CONSIDER writing doc comments for private APIs

#Doc comments aren't just for external consumers of your library's public API. They can also be helpful for understanding private members that are called from other parts of the library.

### DO start doc comments with a single-sentence summary

#Start your doc comment with a brief, user-centric description ending with a period. A sentence fragment is often sufficient. Provide just enough context for the reader to orient themselves and decide if they should keep reading or look elsewhere for the solution to their problem.

```
/// Deletes the file at [path] from the file system.
void delete(String path) {
 ...
}
```
```
/// Depending on the state of the file system and the user's permissions,
/// certain operations may or may not be possible. If there is no file at
/// [path] or it can't be accessed, this function throws either [IOError]
/// or [PermissionError], respectively. Otherwise, this deletes the file.
void delete(String path) {
 ...
}
```
### DO separate the first sentence of a doc comment into its own paragraph

#Add a blank line after the first sentence to split it out into its own paragraph. If more than a single sentence of explanation is useful, put the rest in later paragraphs.

 This helps you write a tight first sentence that summarizes the documentation.
 Also, tools like `dart doc` use the first paragraph as a short summary in places
 like lists of classes and members.

```
/// Deletes the file at [path].
///
/// Throws an [IOError] if the file could not be found. Throws a
/// [PermissionError] if the file is present but could not be deleted.
void delete(String path) {
 ...
}
```
```
/// Deletes the file at [path]. Throws an [IOError] if the file could not
/// be found. Throws a [PermissionError] if the file is present but could
/// not be deleted.
void delete(String path) {
 ...
}
```
### AVOID redundancy with the surrounding context

#
 The reader of a class's doc comment can clearly see the name of the class, what
 interfaces it implements, etc. When reading docs for a member, the signature is
 right there, and the enclosing class is obvious. None of that needs to be
 spelled out in the doc comment. Instead, focus on explaining what the reader
 *doesn't* already know.

```
class RadioButtonWidget extends Widget {
 /// Sets the tooltip to [lines].
 ///
 /// The lines should be word wrapped using the current font.
 void tooltip(List<String> lines) {
 ...
 }
}
```
```
class RadioButtonWidget extends Widget {
 /// Sets the tooltip for this radio button widget to the list of strings in
 /// [lines].
 void tooltip(List<String> lines) {
 ...
 }
}
```
Only add doc comments when providing context, caveats, or usage details not immediately obvious from the surrounding context.

### PREFER starting comments of a function or method with third-person verbs if its main purpose is a side effect

#The doc comment should focus on what the code *does*.

```
/// Connects to the server and fetches the query results.
Stream<QueryResult> fetchResults(Query query) => ...
/// Starts the stopwatch if not already running.
void start() => ...
```
### PREFER starting a non-boolean variable or property comment with a noun phrase

#
 The doc comment should stress what the property *is*. This is true even for
 getters which may do calculation or other work. What the caller cares about is
 the *result* of that work, not the work itself.

```
/// The current day of the week, where `0` is Sunday.
int weekday;
/// The number of checked buttons on the page.
int get checkedCount => ...
```
### PREFER starting a boolean variable or property comment with "Whether" followed by a noun or gerund phrase

#
 The doc comment should clarify the states this variable represents.
 This is true even for getters which may do calculation or other work.
 What the caller cares about is the *result* of that work, not the work itself.

```
/// Whether the modal is currently displayed to the user.
bool isVisible;
/// Whether the modal should confirm the user's intent on navigation.
bool get shouldConfirm => ...
/// Whether resizing the current browser window will also resize the modal.
bool get canResize => ...
```
### PREFER a noun phrase or non-imperative verb phrase for a function or method if returning a value is its primary purpose

#
 If a method is *syntactically* a method, but *conceptually* it is a property,
 and is therefore named with a noun phrase or non-imperative verb phrase,
 it should also be documented as such.
 Use a noun-phrase for such non-boolean functions, and
 a phrase starting with "Whether" for such boolean functions,
 just as for a syntactic property or variable.

```
/// The [index]th element of this iterable in iteration order.
E elementAt(int index);
/// Whether this iterable contains an element equal to [element].
bool contains(Object? element);
```
### DON'T write documentation for both the getter and setter of a property

#
 If a property has both a getter and a setter, then create a doc comment for
 only one of them. `dart doc` treats the getter and setter like a single field,
 and if both the getter and the setter have doc comments, then
 `dart doc` discards the setter's doc comment.

```
/// The pH level of the water in the pool.
///
/// Ranges from 0-14, representing acidic to basic, with 7 being neutral.
int get phLevel => ...
set phLevel(int level) => ...
```
```
/// The depth of the water in the pool, in meters.
int get waterDepth => ...
/// Updates the water depth to a total of [meters] in height.
set waterDepth(int meters) => ...
```
### PREFER starting library or type comments with noun phrases

#Doc comments for classes are often the most important documentation in your program. They describe the type's invariants, establish the terminology it uses, and provide context to the other doc comments for the class's members. A little extra effort here can make all of the other members simpler to document.

The documentation should describe an *instance* of the type.

```
/// A chunk of non-breaking output text terminated by a hard or soft newline.
///
/// ...
class Chunk {
 ...
}
```
### CONSIDER including code samples in doc comments

#```
/// The lesser of two numbers.
///
/// ```dart
/// min(5, 3) == 3
/// ```
num min(num a, num b) => ...
```
Humans are great at generalizing from examples, so even a single code sample makes an API easier to learn.

### DO use square brackets in doc comments to refer to in-scope identifiers

#Linter rule: comment_references

 If you surround things like variable, method, or type names in square brackets,
 then `dart doc` looks up the name and links to the relevant API docs.
 Parentheses are optional but can
 clarify you're referring to a function or constructor.
 The following partial doc comments illustrate a few cases
 where these comment references can be helpful:

```
/// Throws a [StateError] if ...
///
/// Similar to [anotherMethod()], but ...
```
To link to a member of a specific class, use the class name and member name, separated by a dot:

```
/// Similar to [Duration.inDays], but handles fractional days.
```
 The dot syntax can also be used to refer to named constructors. For the unnamed
 constructor, use `.new` after the class name:

```
/// To create a point, call [Point.new] or use [Point.polar] to ...
```
 To learn more about the references that
 the analyzer and `dart doc` support in doc comments,
 check out Documentation comment references.

### DO use prose to explain parameters, return values, and exceptions

#Other languages use verbose tags and sections to describe what the parameters and returns of a method are.

```
/// Defines a flag with the given name and abbreviation.
///
/// @param name The name of the flag.
/// @param abbr The abbreviation for the flag.
/// @returns The new flag.
/// @throws ArgumentError If there is already an option with
/// the given name or abbreviation.
Flag addFlag(String name, String abbreviation) => ...
```
The convention in Dart is to integrate that into the description of the method and highlight parameters using square brackets.

Consider having sections starting with "The [parameter]" to describe parameters, with "Returns" for the returned value and "Throws" for exceptions. Errors can be documented the same way as exceptions, or just as requirements that must be satisfied, without documenting the precise error which will be thrown.

```
/// Defines a flag with the given [name] and [abbreviation].
///
/// The [name] and [abbreviation] strings must not be empty.
///
/// Returns a new flag.
///
/// Throws a [DuplicateFlagException] if there is already an option named
/// [name] or there is already an option using the [abbreviation].
Flag addFlag(String name, String abbreviation) => ...
```
### DO put doc comments before metadata annotations

#```
/// A button that can be flipped on and off.
@Component(selector: 'toggle')
class ToggleComponent {}
```
```
@Component(selector: 'toggle')
/// A button that can be flipped on and off.
class ToggleComponent {}
```
## Markdown

#
 You are allowed to use most markdown formatting in your doc comments and
 `dart doc` will process it accordingly using the markdown package.

There are tons of guides out there already to introduce you to Markdown. Its universal popularity is why we chose it. Here's just a quick example to give you a flavor of what's supported:

```
/// This is a paragraph of regular text.
///
/// This sentence has *two* _emphasized_ words (italics) and **two**
/// __strong__ ones (bold).
///
/// A blank line creates a separate paragraph. It has some `inline code`
/// delimited using backticks.
///
/// * Unordered lists.
/// * Look like ASCII bullet lists.
/// * You can also use `-` or `+`.
///
/// 1. Numbered lists.
/// 2. Are, well, numbered.
/// 1. But the values don't matter.
///
/// * You can nest lists too.
/// * They must be indented at least 4 spaces.
/// * (Well, 5 including the space after `///`.)
///
/// Code blocks are fenced in triple backticks:
///
/// ```dart
/// this.code
/// .will
/// .retain(its, formatting);
/// ```
///
/// The code language (for syntax highlighting) defaults to Dart. You can
/// specify it by putting the name of the language after the opening backticks:
///
/// ```html
/// <h1>HTML is magical!</h1>
/// ```
///
/// Links can be:
///
/// * https://www.just-a-bare-url.com
/// * [with the URL inline](https://google.com)
/// * [or separated out][ref link]
///
/// [ref link]: https://google.com
///
/// # A Header
///
/// ## A subheader
///
/// ### A subsubheader
///
/// #### If you need this many levels of headers, you're doing it wrong
```
### AVOID using markdown excessively

#When in doubt, format less. Formatting exists to illuminate your content, not replace it. Words are what matter.

### AVOID using HTML for formatting

#
 It *may* be useful to use it in rare cases for things like tables, but in almost
 all cases, if it's too complex to express in Markdown, you're better off not
 expressing it.

### PREFER backtick fences for code blocks

#Markdown has two ways to indicate a block of code: indenting the code four spaces on each line, or surrounding it in a pair of triple-backtick "fence" lines. The former syntax is brittle when used inside things like Markdown lists where indentation is already meaningful or when the code block itself contains indented code.

The backtick syntax avoids those indentation woes, lets you indicate the code's language, and is consistent with using backticks for inline code.

```
/// You can use [CodeBlockExample] like this:
///
/// ```dart
/// var example = CodeBlockExample();
/// print(example.isItGreat); // "Yes."
/// ```
```
```
/// You can use [CodeBlockExample] like this:
///
/// var example = CodeBlockExample();
/// print(example.isItGreat); // "Yes."
```
## Writing

#We think of ourselves as programmers, but most of the characters in a source file are intended primarily for humans to read. English is the language we code in to modify the brains of our coworkers. As for any programming language, it's worth putting effort into improving your proficiency.

This section lists a few guidelines for our docs. You can learn more about best practices for technical writing, in general, from articles such as Technical writing style.

### PREFER brevity

#Be clear and precise, but also terse.

### AVOID abbreviations and acronyms unless they are obvious

#Many people don't know what "i.e.", "e.g." and "et al." mean. That acronym that you're sure everyone in your field knows may not be as widely known as you think.

### PREFER using "this" instead of "the" to refer to a member's instance

#When documenting a member for a class, you often need to refer back to the object the member is being called on. Using "the" can be ambiguous. Prefer having some qualifier after "this", a sole "this" can be ambiguous too.

```
class Box {
 /// The value this box wraps.
 Object? _value;
 /// Whether this box contains a value.
 bool get hasValue => _value != null;
}
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Page last updated on 2026-07-31. View source or report an issue.

# Effective Dart: Usage

Guidelines for using language features to write maintainable code.

 You can use these guidelines every day in the bodies of your Dart code. *Users*
 of your library may not be able to tell that you've internalized the ideas here,
 but *maintainers* of it sure will.

## Libraries

#
 These guidelines help you compose your program out of multiple files in a
 consistent, maintainable way. To keep these guidelines brief, they use "import"
 to cover `import` and `export` directives. The guidelines apply equally to both.

### DO use strings in `part of` directives

 #
 Linter rule: use_string_in_part_of_directives

 Many Dart developers avoid using `part` entirely. They find it easier to reason
 about their code when each library is a single file. If you do choose to use
 `part` to split part of a library out into another file, Dart requires the other
 file to in turn indicate which library it's a part of.

 Dart allows the `part of` directive to use the *name* of a library.
 Naming libraries is a legacy feature that is now discouraged.
 Library names can introduce ambiguity
 when determining which library a part belongs to.

 The preferred syntax is to use a URI string that points
 directly to the library file.
 If you have some library, `my_library.dart`, that contains:

```
library my_library;
part 'some/other/file.dart';
```
Then the part file should use the library file's URI string:

```
part of '../../my_library.dart';
```
Not the library name:

```
part of my_library;
```
### DON'T import libraries that are inside the `src` directory of another package

 #
 Linter rule: implementation_imports

 The `src` directory under `lib` is specified
 to contain
 libraries private to the package's own implementation. The way package
 maintainers version their package takes this convention into account. They are
 free to make sweeping changes to code under `src` without it being a breaking
 change to the package.

That means that if you import some other package's private library, a minor, theoretically non-breaking point release of that package could break your code.

### DON'T allow an import path to reach into or out of `lib`

 #
 Linter rule: avoid_relative_lib_imports

 A `package:` import lets you access
 a library inside a package's `lib` directory
 without having to worry about where the package is stored on your computer.
 For this to work, you cannot have imports that require the `lib`
 to be in some location on disk relative to other files.
 In other words, a relative import path in a file inside `lib`
 can't reach out and access a file outside of the `lib` directory,
 and a library outside of `lib` can't use a relative path
 to reach into the `lib` directory.
 Doing either leads to confusing errors and broken programs.

For example, say your directory structure looks like this:

-
 ## my_package/- ## lib/- api.dart

- ## test/- api_test.dart

And say `api_test.dart` imports `api.dart` in two ways:

```
import 'package:my_package/api.dart';
import '../lib/api.dart';
```
Dart thinks those are imports of two completely unrelated libraries. To avoid confusing Dart and yourself, follow these two rules:

- Don't use `/lib/`in import paths.
- Don't use `../`to escape the`lib`directory.

 Instead, when you need to reach into a package's `lib` directory
 (even from the same package's `test` directory
 or any other top-level directory),
 use a `package:` import.

```
import 'package:my_package/api.dart';
```
 A package should never reach *out* of its `lib` directory and
 import libraries from other places in the package.

### PREFER relative import paths

#Linter rule: prefer_relative_imports

 Whenever the previous rule doesn't come into play, follow this one.
 When an import does *not* reach across `lib`, prefer using relative imports.
 They're shorter.
 For example, say your directory structure looks like this:

-
 ## my_package/- ## lib/- ## src/- stuff.dart
- utils.dart

- api.dart

- ## test/- api_test.dart
- test_utils.dart

Here is how the various libraries should import each other:

```
import 'src/stuff.dart';
import 'src/utils.dart';
```
```
import '../api.dart';
import 'stuff.dart';
```
```
import 'package:my_package/api.dart'; // Don't reach into 'lib'.
import 'test_utils.dart'; // Relative within 'test' is fine.
```
## Null

#### DON'T explicitly initialize variables to `null`

 #
 Linter rule: avoid_init_to_null

 If a variable has a non-nullable type, Dart reports a compile error if you try
 to use it before it has been definitely initialized. If the variable is
 nullable, then it is implicitly initialized to `null` for you. There's no
 concept of "uninitialized memory" in Dart and no need to explicitly initialize a
 variable to `null` to be "safe".

```
Item? bestDeal(List<Item> cart) {
 Item? bestItem;
 for (final item in cart) {
 if (bestItem == null || item.price < bestItem.price) {
 bestItem = item;
 }
 }
 return bestItem;
}
```
```
Item? bestDeal(List<Item> cart) {
 Item? bestItem = null;
 for (final item in cart) {
 if (bestItem == null || item.price < bestItem.price) {
 bestItem = item;
 }
 }
 return bestItem;
}
```
### DON'T use an explicit default value of `null`

 #
 Linter rule: avoid_init_to_null

 If you make a nullable parameter optional but don't give it a default value, the
 language implicitly uses `null` as the default, so there's no need to write it.

```
void error([String? message]) {
 stderr.write(message ?? '\n');
}
```
```
void error([String? message = null]) {
 stderr.write(message ?? '\n');
}
```
### DON'T use `true` or `false` in equality operations

 #

 Using the equality operator to evaluate a *non-nullable* boolean expression
 against a boolean literal is redundant.
 It's always simpler to eliminate the equality operator,
 and use the unary negation operator `!` if necessary:

```
if (nonNullableBool) {
 ...
}
if (!nonNullableBool) {
 ...
}
```
```
if (nonNullableBool == true) {
 ...
}
if (nonNullableBool == false) {
 ...
}
```
 To evaluate a boolean expression that *is nullable*, you should use `??`
 or an explicit `!= null` check.

```
// If you want null to result in false:
if (nullableBool ?? false) {
 ...
}
// If you want null to result in false
// and you want the variable to type promote:
if (nullableBool != null && nullableBool) {
 ...
}
```
```
// Static error if null:
if (nullableBool) {
 ...
}
// If you want null to be false:
if (nullableBool == true) {
 ...
}
```
 `nullableBool == true` is a viable expression,
 but shouldn't be used for several reasons:

- It doesn't indicate the code has anything to do with - `null`.
-
 Because it's not evidently `null`related, it can easily be mistaken for the non-nullable case, where the equality operator is redundant and can be removed. That's only true when the boolean expression on the left has no chance of producing null, but not when it can.
-
 The boolean logic is confusing. If `nullableBool`is null, then`nullableBool == true`means the condition evaluates to`false`.

 The `??` operator makes it clear that something to do with null is happening,
 so it won't be mistaken for a redundant operation.
 The logic is much clearer too;
 the result of the expression being `null` is the same as the boolean literal.

 Using a null-aware operator such as `??` on a variable inside a condition
 doesn't promote the variable to a non-nullable type.
 If you want the variable to be promoted inside the body of the `if` statement,
 it's better to use an explicit `!= null` check instead of `??`.

### AVOID `late` variables if you need to check whether they are initialized

 #

 Dart offers no way to tell if a `late` variable
 has been initialized or assigned to.
 If you access it, it either immediately runs the initializer
 (if it has one) or throws an exception.
 Sometimes you have some state that's lazily initialized
 where `late` might be a good fit,
 but you also need to be able to *tell* if the initialization has happened yet.

 Although you could detect initialization by storing the state in a `late` variable
 and having a separate boolean field
 that tracks whether the variable has been set,
 that's redundant because Dart *internally*
 maintains the initialized status of the `late` variable.
 Instead, it's usually clearer to make the variable non-`late` and nullable.
 Then you can see if the variable has been initialized
 by checking for `null`.

 Of course, if `null` is a valid initialized value for the variable,
 then it probably does make sense to have a separate boolean field.

### CONSIDER type promotion or null-check patterns for using nullable types

#
 Checking that a nullable variable is not equal to `null` promotes the variable
 to a non-nullable type. That lets you access members on the variable and pass it
 to functions expecting a non-nullable type.

Type promotion is only supported, however, for local variables, parameters, and private final fields. Values that are open to manipulation can't be type promoted.

Declaring members private and final, as we generally recommend, is often enough to bypass these limitations. But, that's not always an option.

One pattern to work around type promotion limitations is to use a null-check pattern. This simultaneously confirms the member's value is not null, and binds that value to a new non-nullable variable of the same base type.

```
class UploadException {
 final Response? response;
 UploadException([this.response]);
 @override
 String toString() {
 if (this.response case var response?) {
 return 'Could not complete upload to ${response.url} '
 '(error code ${response.errorCode}): ${response.reason}.';
 }
 return 'Could not upload (no response).';
 }
}
```
Another work around is to assign the field's value to a local variable. Null checks on that variable will promote, so you can safely treat it as non-nullable.

```
class UploadException {
 final Response? response;
 UploadException([this.response]);
 @override
 String toString() {
 final response = this.response;
 if (response != null) {
 return 'Could not complete upload to ${response.url} '
 '(error code ${response.errorCode}): ${response.reason}.';
 }
 return 'Could not upload (no response).';
 }
}
```
 Be careful when using a local variable. If you need to write back to the field,
 make sure that you don't write back to the local variable instead. (Making the
 local variable `final`
 can prevent such mistakes.) Also, if the field might
 change while the local is still in scope, then the local might have a stale
 value.

 Sometimes it's best to simply use `!`
 on the field.
 In some cases, though, using either a local variable or a null-check pattern
 can be cleaner and safer than using `!` every time you need to treat the value
 as non-null:

```
class UploadException {
 final Response? response;
 UploadException([this.response]);
 @override
 String toString() {
 if (response != null) {
 return 'Could not complete upload to ${response!.url} '
 '(error code ${response!.errorCode}): ${response!.reason}.';
 }
 return 'Could not upload (no response).';
 }
}
```
## Strings

#Here are some best practices to keep in mind when composing strings in Dart.

### DO use adjacent strings to concatenate string literals

#Linter rule: prefer_adjacent_string_concatenation

 If you have two string literals—not values, but the actual quoted literal
 form—you do not need to use `+` to concatenate them. Just like in C and
 C++, simply placing them next to each other does it. This is a good way to make
 a single long string that doesn't fit on one line.

```
raiseAlarm(
 'ERROR: Parts of the spaceship are on fire. Other '
 'parts are overrun by martians. Unclear which are which.',
);
```
```
raiseAlarm(
 'ERROR: Parts of the spaceship are on fire. Other ' +
 'parts are overrun by martians. Unclear which are which.',
);
```
### PREFER using interpolation to compose strings and values

#Linter rule: prefer_interpolation_to_compose_strings

 If you're coming from other languages, you're used to using long chains of `+`
 to build a string out of literals and other values. That does work in Dart, but
 it's almost always cleaner and shorter to use interpolation:

```
'Hello, $name! You are ${year - birth} years old.';
```
```
'Hello, ' + name + '! You are ' + (year - birth).toString() + ' y...';
```
 Note that this guideline applies to combining *multiple* literals and values.
 It's fine to use `.toString()` when converting only a single object to a string.

### AVOID using curly braces in interpolation when not needed

#Linter rule: unnecessary_brace_in_string_interps

 If you're interpolating a simple identifier not immediately followed by more
 alphanumeric text, the `{}` should be omitted.

```
var greeting = 'Hi, $name! I love your ${decade}s costume.';
```
```
var greeting = 'Hi, ${name}! I love your ${decade}s costume.';
```
## Collections

#Out of the box, Dart supports four collection types: lists, maps, queues, and sets. The following best practices apply to collections.

### DO use collection literals when possible

#Linter rule: prefer_collection_literals

Dart has three core collection types: List, Map, and Set. The Map and Set classes have unnamed constructors like most classes do. But because these collections are used so frequently, Dart has nicer built-in syntax for creating them:

```
var points = <Point>[];
var addresses = <String, Address>{};
var counts = <int>{};
```
```
var addresses = Map<String, Address>();
var counts = Set<int>();
```
 Note that this guideline doesn't apply to the *named* constructors for those
 classes. `List.from()`, `Map.fromIterable()`, and friends all have their uses.
 (The List class also has an unnamed constructor, but it is prohibited in null
 safe Dart.)

 Collection literals are particularly powerful in Dart
 because they give you access to the spread operator

 for including the contents of other collections,
 and `if` and `for`
 for performing control flow while
 building the contents:

```
var arguments = [
 ...options,
 command,
 ...?modeFlags,
 for (var path in filePaths)
 if (path.endsWith('.dart')) path.replaceAll('.dart', '.js'),
];
```
```
var arguments = <String>[];
arguments.addAll(options);
arguments.add(command);
if (modeFlags != null) arguments.addAll(modeFlags);
arguments.addAll(
 filePaths
 .where((path) => path.endsWith('.dart'))
 .map((path) => path.replaceAll('.dart', '.js')),
);
```
### DON'T use `.length` to see if a collection is empty

 #
 Linter rules: prefer_is_empty, prefer_is_not_empty

 The Iterable contract does not require that a collection know its length or
 be able to provide it in constant time. Calling `.length` just to see if the
 collection contains *anything* can be painfully slow.

 Instead, there are faster and more readable getters: `.isEmpty` and
 `.isNotEmpty`. Use the one that doesn't require you to negate the result.

```
if (lunchBox.isEmpty) return 'so hungry...';
if (words.isNotEmpty) return words.join(' ');
```
```
if (lunchBox.length == 0) return 'so hungry...';
if (!words.isEmpty) return words.join(' ');
```
### AVOID using `Iterable.forEach()` with a function literal

 #
 Linter rule: avoid_function_literals_in_foreach_calls

 `forEach()` functions are widely used in JavaScript because the built in
 `for-in` loop doesn't do what you usually want. In Dart, if you want to iterate
 over a sequence, the idiomatic way to do that is using a loop.

```
for (final person in people) {
 ...
}
```
```
people.forEach((person) {
 ...
});
```
 Note that this guideline specifically says "function *literal*". If you want to
 invoke some *already existing* function on each element, `forEach()`
 is fine.

```
people.forEach(print);
```
 Also note that it's always OK to use `Map.forEach()`. Maps aren't iterable, so
 this guideline doesn't apply.

### DON'T use `List.from()` unless you intend to change the type of the result

 #
 Given an Iterable, there are two obvious ways to produce a new List that contains the same elements:

```
var copy1 = iterable.toList();
var copy2 = List.from(iterable);
```
 The obvious difference is that the first one is shorter. The *important*
 difference is that the first one preserves the type argument of the original
 object:

```
// Creates a List<int>:
var iterable = [1, 2, 3];
// Prints "List<int>":
print(iterable.toList().runtimeType);
```
```
// Creates a List<int>:
var iterable = [1, 2, 3];
// Prints "List<dynamic>":
print(List.from(iterable).runtimeType);
```
If you *want* to change the type, then calling `List.from()` is useful:

```
var numbers = [1, 2.3, 4]; // List<num>.
numbers.removeAt(1); // Now it only contains integers.
var ints = List<int>.from(numbers);
```
 But if your goal is just to copy the iterable and preserve its original type, or
 you don't care about the type, then use `toList()`.

### DO use `whereType()` to filter a collection by type

 #
 Linter rule: prefer_iterable_whereType

 Let's say you have a list containing a mixture of objects, and you want to get
 just the integers out of it. You could use `where()` like this:

```
var objects = [1, 'a', 2, 'b', 3];
var ints = objects.where((e) => e is int);
```
 This is verbose, but, worse, it returns an iterable whose type probably isn't
 what you want. In the example here, it returns an `Iterable<Object>`
 even though
 you likely want an `Iterable<int>` since that's the type you're filtering it to.

Sometimes you see code that "corrects" the above error by adding `cast()`:

```
var objects = [1, 'a', 2, 'b', 3];
var ints = objects.where((e) => e is int).cast<int>();
```
 That's verbose and causes two wrappers to be created, with two layers of
 indirection and redundant runtime checking. Fortunately, the core library has
 the `whereType()`
 method for this exact use case:

```
var objects = [1, 'a', 2, 'b', 3];
var ints = objects.whereType<int>();
```
 Using `whereType()` is concise, produces an Iterable
 of the desired type,
 and has no unnecessary levels of wrapping.

### DON'T use `cast()` when a nearby operation will do

 #

 Often when you're dealing with an iterable or stream, you perform several
 transformations on it. At the end, you want to produce an object with a certain
 type argument. Instead of tacking on a call to `cast()`, see if one of the
 existing transformations can change the type.

 If you're already calling `toList()`, replace that with a call to
 `List<T>.from()`
 where `T` is the type of resulting list you want.

```
var stuff = <dynamic>[1, 2];
var ints = List<int>.from(stuff);
```
```
var stuff = <dynamic>[1, 2];
var ints = stuff.toList().cast<int>();
```
 If you are calling `map()`, give it an explicit type argument so that it
 produces an iterable of the desired type. Type inference often picks the correct
 type for you based on the function you pass to `map()`, but sometimes you need
 to be explicit.

```
var stuff = <dynamic>[1, 2];
var reciprocals = stuff.map<double>((n) => n * 2);
```
```
var stuff = <dynamic>[1, 2];
var reciprocals = stuff.map((n) => n * 2).cast<double>();
```
### AVOID using `cast()`

 #

 This is the softer generalization of the previous rule. Sometimes there is no
 nearby operation you can use to fix the type of some object. Even then, when
 possible avoid using `cast()` to "change" a collection's type.

Prefer any of these options instead:

-
 **Create it with the right type.**Change the code where the collection is first created so that it has the right type.
-
 **Cast the elements on access.**If you immediately iterate over the collection, cast each element inside the iteration.
-
 **Eagerly cast using**If you'll eventually access most of the elements in the collection, and you don't need the object to be backed by the original live object, convert it using`List.from()`.`List.from()`.The `cast()`method returns a lazy collection that checks the element type on*every operation*. If you perform only a few operations on only a few elements, that laziness can be good. But in many cases, the overhead of lazy validation and of wrapping outweighs the benefits.

Here is an example of **creating it with the right type:**

```
List<int> singletonList(int value) {
 var list = <int>[];
 list.add(value);
 return list;
}
```
```
List<int> singletonList(int value) {
 var list = []; // List<dynamic>.
 list.add(value);
 return list.cast<int>();
}
```
Here is **casting each element on access:**

```
void printEvens(List<Object> objects) {
 // We happen to know the list only contains ints.
 for (final n in objects) {
 if ((n as int).isEven) print(n);
 }
}
```
```
void printEvens(List<Object> objects) {
 // We happen to know the list only contains ints.
 for (final n in objects.cast<int>()) {
 if (n.isEven) print(n);
 }
}
```
Here is **casting eagerly using List.from():**

```
int median(List<Object> objects) {
 // We happen to know the list only contains ints.
 var ints = List<int>.from(objects);
 ints.sort();
 return ints[ints.length ~/ 2];
}
```
```
int median(List<Object> objects) {
 // We happen to know the list only contains ints.
 var ints = objects.cast<int>();
 ints.sort();
 return ints[ints.length ~/ 2];
}
```
 These alternatives don't always work, of course, and sometimes `cast()` is the
 right answer. But consider that method a little risky and undesirable—it
 can be slow and may fail at runtime if you aren't careful.

## Functions

#In Dart, even functions are objects. Here are some best practices involving functions.

### DO use a function declaration to bind a function to a name

#Linter rule: prefer_function_declarations_over_variables

Modern languages have realized how useful local nested functions and closures are. It's common to have a function defined inside another one. In many cases, this function is used as a callback immediately and doesn't need a name. A function expression is great for that.

But, if you do need to give it a name, use a function declaration statement instead of binding a lambda to a variable.

```
void main() {
 void localFunction() {
 ...
 }
}
```
```
void main() {
 var localFunction = () {
 ...
 };
}
```
### DON'T create a lambda when a tear-off will do

#Linter rule: unnecessary_lambdas

 When you refer to a function, method, or named constructor without parentheses,
 Dart creates a *tear-off*. This is a closure that takes the same
 parameters as the function and invokes the underlying function when you call it.
 If your code needs a closure that invokes a named function with the same
 parameters as the closure accepts, don't wrap the call in a lambda.
 Use a tear-off.

```
var charCodes = [68, 97, 114, 116];
var buffer = StringBuffer();
// Function:
charCodes.forEach(print);
// Method:
charCodes.forEach(buffer.write);
// Named constructor:
var strings = charCodes.map(String.fromCharCode);
// Unnamed constructor:
var buffers = charCodes.map(StringBuffer.new);
```
```
var charCodes = [68, 97, 114, 116];
var buffer = StringBuffer();
// Function:
charCodes.forEach((code) {
 print(code);
});
// Method:
charCodes.forEach((code) {
 buffer.write(code);
});
// Named constructor:
var strings = charCodes.map((code) => String.fromCharCode(code));
// Unnamed constructor:
var buffers = charCodes.map((code) => StringBuffer(code));
```
## Variables

#The following best practices describe how to best use variables in Dart.

### DO follow a consistent rule for `var` and `final` on local variables

 #

 Most local variables shouldn't have type annotations and should be declared
 using just `var` or `final`. There are two rules in wide use for when to use one
 or the other:

-
 Use `final`for local variables that are not reassigned and`var`for those that are.
-
 Use `var`for all local variables, even ones that aren't reassigned. Never use`final`for locals. (Using`final`for fields and top-level variables is still encouraged, of course.)

 Either rule is acceptable, but pick *one* and apply it consistently throughout
 your code. That way when a reader sees `var`, they know whether it means that
 the variable is assigned later in the function.

### AVOID storing what you can calculate

#When designing a class, you often want to expose multiple views into the same underlying state. Often you see code that calculates all of those views in the constructor and then stores them:

```
class Circle {
 double radius;
 double area;
 double circumference;
 Circle(double radius)
 : radius = radius,
 area = pi * radius * radius,
 circumference = pi * 2.0 * radius;
}
```
 This code has two things wrong with it. First, it's likely wasting memory. The
 area and circumference, strictly speaking, are *caches*. They are stored
 calculations that we could recalculate from other data we already have. They are
 trading increased memory for reduced CPU usage. Do we know we have a performance
 problem that merits that trade-off?

 Worse, the code is *wrong*. The problem with caches is *invalidation*—how
 do you know when the cache is out of date and needs to be recalculated? Here, we
 never do, even though `radius` is mutable. You can assign a different value and
 the `area` and `circumference` will retain their previous, now incorrect values.

To correctly handle cache invalidation, we would need to do this:

```
class Circle {
 double _radius;
 double get radius => _radius;
 set radius(double value) {
 _radius = value;
 _recalculate();
 }
 double _area = 0.0;
 double get area => _area;
 double _circumference = 0.0;
 double get circumference => _circumference;
 Circle(this._radius) {
 _recalculate();
 }
 void _recalculate() {
 _area = pi * _radius * _radius;
 _circumference = pi * 2.0 * _radius;
 }
}
```
That's an awful lot of code to write, maintain, debug, and read. Instead, your first implementation should be:

```
class Circle {
 double radius;
 Circle(this.radius);
 double get area => pi * radius * radius;
 double get circumference => pi * 2.0 * radius;
}
```
This code is shorter, uses less memory, and is less error-prone. It stores the minimal amount of data needed to represent the circle. There are no fields to get out of sync because there is only a single source of truth.

In some cases, you may need to cache the result of a slow calculation, but only do that after you know you have a performance problem, do it carefully, and leave a comment explaining the optimization.

## Members

#In Dart, objects have members which can be functions (methods) or data (instance variables). The following best practices apply to an object's members.

### DON'T wrap a field in a getter and setter unnecessarily

#Linter rule: unnecessary_getters_setters

In Java and C#, it's common to hide all fields behind getters and setters (or properties in C#), even if the implementation just forwards to the field. That way, if you ever need to do more work in those members, you can without needing to touch the call sites. This is because calling a getter method is different than accessing a field in Java, and accessing a property isn't binary-compatible with accessing a raw field in C#.

Dart doesn't have this limitation. Fields and getters/setters are completely indistinguishable. You can expose a field in a class and later wrap it in a getter and setter without having to touch any code that uses that field.

```
class Box {
 Object? contents;
}
```
```
class Box {
 Object? _contents;
 Object? get contents => _contents;
 set contents(Object? value) {
 _contents = value;
 }
}
```
### PREFER using a `final` field to make a read-only property

 #

 If you have a field that outside code should be able to see but not assign to, a
 simple solution that works in many cases is to simply mark it `final`.

```
class Box {
 final contents = [];
}
```
```
class Box {
 Object? _contents;
 Object? get contents => _contents;
}
```
Of course, if you need to internally assign to the field outside of the constructor, you may need to do the "private field, public getter" pattern, but don't reach for that until you need to.

### CONSIDER using `=>` for simple members

 #
 Linter rule: prefer_expression_function_bodies

 In addition to using `=>` for function expressions, Dart also lets you define
 members with it. That style is a good fit for simple members that just calculate
 and return a value.

```
double get area => (right - left) * (bottom - top);
String capitalize(String name) =>
 '${name[0].toUpperCase()}${name.substring(1)}';
```
 People *writing* code seem to love `=>`, but it's very easy to abuse it and end
 up with code that's hard to *read*. If your declaration is more than a couple of
 lines or contains deeply nested expressions—cascades and conditional
 operators are common offenders—do yourself and everyone who has to read
 your code a favor and use a block body and some statements.

```
Treasure? openChest(Chest chest, Point where) {
 if (_opened.containsKey(chest)) return null;
 var treasure = Treasure(where);
 treasure.addAll(chest.contents);
 _opened[chest] = treasure;
 return treasure;
}
```
```
Treasure? openChest(Chest chest, Point where) => _opened.containsKey(chest)
 ? null
 : _opened[chest] = (Treasure(where)..addAll(chest.contents));
```
 You can also use `=>` on members that don't return a value. This is idiomatic
 when a setter is small and has a corresponding getter that uses `=>`.

```
num get x => center.x;
set x(num value) => center = Point(value, center.y);
```
### DON'T use `this.` except to redirect to a named constructor or to avoid shadowing

 #
 Linter rule: unnecessary_this

 JavaScript requires an explicit `this.` to refer to members on the object whose
 method is currently being executed, but Dart—like C++, Java, and
 C#—doesn't have that limitation.

 There are only two times you need to use `this.`. One is when a local variable
 with the same name shadows the member you want to access:

```
class Box {
 Object? value;
 void clear() {
 this.update(null);
 }
 void update(Object? value) {
 this.value = value;
 }
}
```
```
class Box {
 Object? value;
 void clear() {
 update(null);
 }
 void update(Object? value) {
 this.value = value;
 }
}
```
The other time to use `this.` is when redirecting to a named constructor:

```
class ShadeOfGray {
 final int brightness;
 ShadeOfGray(int val) : brightness = val;
 ShadeOfGray.black() : this(0);
 // This won't parse or compile!
 // ShadeOfGray.alsoBlack() : black();
}
```
```
class ShadeOfGray {
 final int brightness;
 ShadeOfGray(int val) : brightness = val;
 ShadeOfGray.black() : this(0);
 // But now it will!
 ShadeOfGray.alsoBlack() : this.black();
}
```
Note that constructor parameters never shadow fields in constructor initializer lists:

```
class Box extends BaseBox {
 Object? value;
 Box(Object? value) : value = value, super(value);
}
```
This looks surprising, but works like you want. Fortunately, code like this is relatively rare thanks to initializing formals and super initializers.

### DO initialize fields at their declaration when possible

#If a field doesn't depend on any constructor parameters, it can and should be initialized at its declaration. It takes less code and avoids duplication when the class has multiple constructors.

```
class ProfileMark {
 final String name;
 final DateTime start;
 ProfileMark(this.name) : start = DateTime.now();
 ProfileMark.unnamed() : name = '', start = DateTime.now();
}
```
```
class ProfileMark {
 final String name;
 final DateTime start = DateTime.now();
 ProfileMark(this.name);
 ProfileMark.unnamed() : name = '';
}
```
 Some fields can't be initialized at their declarations because they need to reference
 `this`—to use other fields or call methods, for example. However, if the
 field is marked `late`, then the initializer *can* access `this`.

Of course, if a field depends on constructor parameters, or is initialized differently by different constructors, then this guideline does not apply.

## Constructors

#The following best practices apply to declaring constructors for a class.

### DO use initializing formals when possible

#Linter rule: prefer_initializing_formals

Many fields are initialized directly from a constructor parameter, like:

```
class Point {
 double x, y;
 Point(double x, double y) : x = x, y = y;
}
```
We've got to type `x` *four* times here to define a field. We can do better:

```
class Point {
 double x, y;
 Point(this.x, this.y);
}
```
 This `this.` syntax before a constructor parameter is called an "initializing
 formal". You can't always take advantage of it. Sometimes you want to have a
 named parameter whose name doesn't match the name of the field you are
 initializing. But when you *can* use initializing formals, you *should*.

### DON'T use `late` when a constructor initializer list will do

 #
 Dart requires you to initialize non-nullable fields before they can be read. Since fields can be read inside the constructor body, this means you get an error if you don't initialize a non-nullable field before the body runs.

 You can make this error go away by marking the field `late`. That turns the
 compile-time error into a *runtime* error if you access the field before it is
 initialized. That's what you need in some cases, but often the right fix is to
 initialize the field in the constructor initializer list:

```
class Point {
 double x, y;
 Point.polar(double theta, double radius)
 : x = cos(theta) * radius,
 y = sin(theta) * radius;
}
```
```
class Point {
 late double x, y;
 Point.polar(double theta, double radius) {
 x = cos(theta) * radius;
 y = sin(theta) * radius;
 }
}
```
 The initializer list gives you access to constructor parameters and lets you
 initialize fields before they can be read. So, if it's possible to use an initializer list,
 that's better than making the field `late` and losing some static safety and
 performance.

### DO use `;` instead of `{}` for empty constructor bodies

 #
 Linter rule: empty_constructor_bodies

In Dart, a constructor with an empty body can be terminated with just a semicolon. (In fact, it's required for const constructors.)

```
class Point {
 double x, y;
 Point(this.x, this.y);
}
```
```
class Point {
 double x, y;
 Point(this.x, this.y) {}
}
```
### PREFER using concise constructor syntax

#Linter rule: unnecessary_type_name_in_constructor

 When declaring constructors inside the class body,
 use the concise constructor syntax instead of repeating the class name.
 Use `factory` for factory constructors
 and `new` for all other constructors.

For a complete reference of the concise syntax forms, see the concise constructor syntax table.

```
class Logger {
 final String name;
 factory(String name) => _cache[name] ??= Logger._internal(name);
 new _internal(this.name);
 new fromJson(Map<String, Object?> json) : name = json['name'] as String;
 static final Map<String, Logger> _cache = {};
}
```
```
class Logger {
 final String name;
 factory Logger(String name) => _cache[name] ??= Logger._internal(name);
 Logger._internal(this.name);
 Logger.fromJson(Map<String, Object?> json) : name = json['name'] as String;
 static final Map<String, Logger> _cache = {};
}
```
 Repeating the class name is redundant
 and makes refactoring harder if the class is renamed.
 Using `new` and `factory` keeps constructor declarations concise
 and consistent.

### DON'T use `new`

 #
 Linter rule: unnecessary_new

 The `new` keyword is optional when calling a constructor.
 Its meaning is not clear because factory constructors mean a
 `new` invocation may not actually return a new object.

 The language still permits `new`, but consider
 it deprecated and avoid using it in your code.

```
Widget build(BuildContext context) {
 return Row(
 children: [
 RaisedButton(child: Text('Increment')),
 Text('Click!'),
 ],
 );
}
```
```
Widget build(BuildContext context) {
 return new Row(
 children: [
 new RaisedButton(child: new Text('Increment')),
 new Text('Click!'),
 ],
 );
}
```
### DON'T use `const` redundantly

 #
 Linter rule: unnecessary_const

 In contexts where an expression *must* be constant, the `const` keyword is
 implicit, doesn't need to be written, and shouldn't. Those contexts are any
 expression inside:

- A const collection literal.
- A const constructor call
- A metadata annotation.
- The initializer for a const variable declaration.
-
 A switch case expression—the part right after `case`before the`:`, not the body of the case.

(Default values are not included in this list because future versions of Dart may support non-const default values.)

 Basically, any place where it would be an error to write `new` instead of
 `const`, Dart allows you to omit the `const`.

```
const primaryColors = [
 Color('red', [255, 0, 0]),
 Color('green', [0, 255, 0]),
 Color('blue', [0, 0, 255]),
];
```
```
const primaryColors = const [
 const Color('red', const [255, 0, 0]),
 const Color('green', const [0, 255, 0]),
 const Color('blue', const [0, 0, 255]),
];
```
## Error handling

#Dart uses exceptions when an error occurs in your program. The following best practices apply to catching and throwing exceptions.

### AVOID catches without `on` clauses

 #
 Linter rule: avoid_catches_without_on_clauses

 A catch clause with no `on` qualifier catches *anything* thrown by the code in
 the try block. Pokémon exception handling
 is very likely not what you
 want. Does your code correctly handle StackOverflowError
 or
 OutOfMemoryError? If you incorrectly pass the wrong argument to a method in
 that try block do you want to have your debugger point you to the mistake or
 would you rather that helpful ArgumentError
 get swallowed? Do you want any
 `assert()` statements inside that code to effectively vanish since you're
 catching the thrown AssertionErrors?

 The answer is probably "no", in which case you should filter the types you
 catch. In most cases, you should have an `on` clause that limits you to the
 kinds of runtime failures you are aware of and are correctly handling.

 In rare cases, you may wish to catch any runtime error. This is usually in
 framework or low-level code that tries to insulate arbitrary application code
 from causing problems. Even here, it is usually better to catch Exception

 than to catch all types. Exception is the base class for all *runtime* errors
 and excludes errors that indicate *programmatic* bugs in the code.

### DON'T discard errors from catches without `on` clauses

 #

 If you really do feel you need to catch *everything* that can be thrown from a
 region of code, *do something* with what you catch. Log it, display it to the
 user or rethrow it, but do not silently discard it.

### DO throw objects that implement `Error` only for programmatic errors

 #

 The Error class is the base class for *programmatic*
 errors. When an object
 of that type or one of its subinterfaces like ArgumentError
 is thrown, it
 means there is a *bug* in your code. When your API wants to report to a caller
 that it is being used incorrectly throwing an Error sends that signal clearly.

Conversely, if the exception is some kind of runtime failure that doesn't indicate a bug in the code, then throwing an Error is misleading. Instead, throw one of the core Exception classes or some other type.

### DON'T explicitly catch `Error` or types that implement it

 #
 Linter rule: avoid_catching_errors

This follows from the above. Since an Error indicates a bug in your code, it should unwind the entire callstack, halt the program, and print a stack trace so you can locate and fix the bug.

 Catching errors of these types breaks that process and masks the bug. Instead of
 *adding* error-handling code to deal with this exception after the fact, go back
 and fix the code that is causing it to be thrown in the first place.

### DO use `rethrow` to rethrow a caught exception

 #
 Linter rule: use_rethrow_when_possible

 If you decide to rethrow an exception, prefer using the `rethrow` statement
 instead of throwing the same exception object using `throw`.
 `rethrow` preserves the original stack trace of the exception. `throw`
 on the
 other hand resets the stack trace to the last thrown position.

```
try {
 somethingRisky();
} catch (e) {
 if (!canHandle(e)) throw e;
 handle(e);
}
```
```
try {
 somethingRisky();
} catch (e) {
 if (!canHandle(e)) rethrow;
 handle(e);
}
```
## Asynchrony

#Dart has several language features to support asynchronous programming. The following best practices apply to asynchronous coding.

### PREFER async/await over using raw futures

#
 Asynchronous code is notoriously hard to read and debug, even when using a nice
 abstraction like futures. The `async`/`await` syntax improves readability and
 lets you use all of the Dart control flow structures within your async code.

```
Future<int> countActivePlayers(String teamName) async {
 try {
 var team = await downloadTeam(teamName);
 if (team == null) return 0;
 var players = await team.roster;
 return players.where((player) => player.isActive).length;
 } on DownloadException catch (e) {
 log.error(e);
 return 0;
 }
}
```
```
Future<int> countActivePlayers(String teamName) {
 return downloadTeam(teamName)
 .then((team) {
 if (team == null) return Future.value(0);
 return team.roster.then((players) {
 return players.where((player) => player.isActive).length;
 });
 })
 .onError<DownloadException>((e, _) {
 log.error(e);
 return 0;
 });
}
```
### DON'T use `async` when it has no useful effect

 #

 It's easy to get in the habit of using `async` on any function that does
 anything related to asynchrony. But in some cases, it's extraneous. If you can
 omit the `async` without changing the behavior of the function, do so.

```
Future<int> fastestBranch(Future<int> left, Future<int> right) {
 return Future.any([left, right]);
}
```
```
Future<int> fastestBranch(Future<int> left, Future<int> right) async {
 return Future.any([left, right]);
}
```
Cases where `async` *is* useful include:

- You are using - `await`. (This is the obvious one.)
-
 You are returning an error asynchronously. `async`and then`throw`is shorter than`return Future.error(...)`.
-
 You are returning a value and you want it implicitly wrapped in a future. `async`is shorter than`Future.value(...)`.

```
Future<void> usesAwait(Future<String> later) async {
 print(await later);
}
Future<void> asyncError() async {
 throw 'Error!';
}
Future<String> asyncValue() async => 'value';
```
### CONSIDER using higher-order methods to transform a stream

#This parallels the above suggestion on iterables. Streams support many of the same methods and also handle things like transmitting errors, closing, etc. correctly.

### AVOID using Completer directly

#Many people new to asynchronous programming want to write code that produces a future. The constructors in Future don't seem to fit their need so they eventually find the Completer class and use that.

```
Future<bool> fileContainsBear(String path) {
 var completer = Completer<bool>();
 File(path).readAsString().then((contents) {
 completer.complete(contents.contains('bear'));
 });
 return completer.future;
}
```
 Completer is needed for two kinds of low-level code: new asynchronous
 primitives, and interfacing with asynchronous code that doesn't use futures.
 Most other code should use async/await or `Future.then()`, because
 they're clearer and make error handling easier.

```
Future<bool> fileContainsBear(String path) {
 return File(path).readAsString().then((contents) {
 return contents.contains('bear');
 });
}
```
```
Future<bool> fileContainsBear(String path) async {
 var contents = await File(path).readAsString();
 return contents.contains('bear');
}
```
###
 DO test for `Future<T>` when disambiguating a `FutureOr<T>` whose type argument could be
 `Object`

 #

 Before you can do anything useful with a `FutureOr<T>`, you typically need to do
 an `is` check to see if you have a `Future<T>` or a bare
 `T`. If the type
 argument is some specific type as in `FutureOr<int>`, it doesn't matter which
 test you use, `is int` or `is Future<int>`. Either works because those two types
 are disjoint.

 However, if the value type is `Object` or a type parameter that could possibly
 be instantiated with `Object`, then the two branches overlap. `Future<Object>`

 itself implements `Object`, so `is Object` or `is T`
 where `T` is some type
 parameter that could be instantiated with `Object` returns true even when the
 object is a future. Instead, explicitly test for the `Future` case:

```
Future<T> logValue<T>(FutureOr<T> value) async {
 if (value is Future<T>) {
 var result = await value;
 print(result);
 return result;
 } else {
 print(value);
 return value;
 }
}
```
```
Future<T> logValue<T>(FutureOr<T> value) async {
 if (value is T) {
 print(value);
 return value;
 } else {
 var result = await value;
 print(result);
 return result;
 }
}
```
 In the bad example, if you pass it a `Future<Object>`, it incorrectly treats it
 like a bare, synchronous value.

Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Page last updated on 2026-07-22. View source or report an issue.

# Effective Dart: Design

Design consistent, usable libraries.

Here are some guidelines for writing consistent, usable APIs for libraries.

## Names

#Naming is an important part of writing readable, maintainable code. The following best practices can help you achieve that goal.

### DO use terms consistently

#Use the same name for the same thing, throughout your code. If a precedent already exists outside your API that users are likely to know, follow that precedent.

```
pageCount // A field.
updatePageCount() // Consistent with pageCount.
toSomething() // Consistent with Iterable's toList().
asSomething() // Consistent with List's asMap().
Point // A familiar concept.
```
```
renumberPages() // Confusingly different from pageCount.
convertToSomething() // Inconsistent with toX() precedent.
wrappedAsSomething() // Inconsistent with asX() precedent.
Cartesian // Unfamiliar to most users.
```
The goal is to take advantage of what the user already knows. This includes their knowledge of the problem domain itself, the conventions of the core libraries, and other parts of your own API. By building on top of those, you reduce the amount of new knowledge they have to acquire before they can be productive.

### AVOID abbreviations

#Unless the abbreviation is more common than the unabbreviated term, don't abbreviate. If you do abbreviate, capitalize it correctly.

```
pageCount
buildRectangles
IOStream
HttpRequest
```
```
numPages // "Num" is an abbreviation of "number (of)".
buildRects
InputOutputStream
HypertextTransferProtocolRequest
```
### PREFER putting the most descriptive noun last

#The last word should be the most descriptive of what the thing is. You can prefix it with other words, such as adjectives, to further describe the thing.

```
pageCount // A count (of pages).
ConversionSink // A sink for doing conversions.
ChunkedConversionSink // A ConversionSink that's chunked.
CssFontFaceRule // A rule for font faces in CSS.
```
```
numPages // Not a collection of pages.
CanvasRenderingContext2D // Not a "2D".
RuleFontFaceCss // Not a CSS.
```
### CONSIDER making the code read like a sentence

#When in doubt about naming, write some code that uses your API, and try to read it like a sentence.

```
// "If errors is empty..."
if (errors.isEmpty) {
 // ...
}
// "Hey, subscription, cancel!"
subscription.cancel();
// "Get the monsters where the monster has claws."
monsters.where((monster) => monster.hasClaws);
```
```
// Telling errors to empty itself, or asking if it is?
if (errors.empty) {
 // ...
}
// Toggle what? To what?
subscription.toggle();
// Filter the monsters with claws *out* or include *only* those?
monsters.filter((monster) => monster.hasClaws);
```
 It's helpful to try out your API and see how it "reads" when used in code, but
 you can go too far. It's not helpful to add articles and other parts of speech
 to force your names to *literally* read like a grammatically correct sentence.

```
if (theCollectionOfErrors.isEmpty) {
 // ...
}
monsters.producesANewSequenceWhereEach((monster) => monster.hasClaws);
```
### PREFER a noun phrase for a non-boolean property or variable

#
 The reader's focus is on *what* the property is. If the user cares more about
 *how* a property is determined, then it should probably be a method with a
 verb phrase name.

```
list.length
context.lineWidth
quest.rampagingSwampBeast
```
```
list.deleteItems
```
### PREFER a non-imperative verb phrase for a boolean property or variable

#Boolean names are often used as conditions in control flow, so you want a name that reads well there. Compare:

```
if (window.closeable) ... // Adjective.
if (window.canClose) ... // Verb.
```
Good names tend to start with one of a few kinds of verbs:

-
 a form of "to be": `isEnabled`,`wasShown`,`willFire`. These are, by far, the most common.
-
 an auxiliary verb: `hasElements`,`canClose`,`shouldConsume`,`mustSave`.
-
 an active verb: `ignoresInput`,`wroteFile`. These are rare because they are usually ambiguous.`loggedResult`is a bad name because it could mean "whether or not a result was logged" or "the result that was logged". Likewise,`closingConnection`could be "whether the connection is closing" or "the connection that is closing". Active verbs are allowed when the name can*only*be read as a predicate.

 What separates all these verb phrases from method names is that they are not
 *imperative*. A boolean name should never sound like a command to tell the
 object to do something, because accessing a property doesn't change the object.
 (If the property *does* modify the object in a meaningful way, it should be a
 method.)

```
isEmpty
hasElements
canClose
closesWindow
canShowPopup
hasShownPopup
```
```
empty // Adjective or verb?
withElements // Sounds like it might hold elements.
closeable // Sounds like an interface.
 // "canClose" reads better as a sentence.
closingWindow // Returns a bool or a window?
showPopup // Sounds like it shows the popup.
```
### CONSIDER omitting the verb for a named boolean *parameter*

 #
 This refines the previous rule. For named parameters that are boolean, the name is often just as clear without the verb, and the code reads better at the call site.

```
Isolate.spawn(entryPoint, message, paused: false);
var copy = List.from(elements, growable: true);
var regExp = RegExp(pattern, caseSensitive: false);
```
### PREFER the "positive" name for a boolean property or variable

#
 Most boolean names have conceptually "positive" and "negative" forms where the
 former feels like the fundamental concept and the latter is its
 negation—"open" and "closed", "enabled" and "disabled", etc. Often the
 latter name literally has a prefix that negates the former: "visible" and
 "*in*-visible", "connected" and "*dis*-connected", "zero" and "*non*-zero".

 When choosing which of the two cases that `true` represents—and thus
 which case the property is named for—prefer the positive or more
 fundamental one. Boolean members are often nested inside logical expressions,
 including negation operators. If your property itself reads like a negation,
 it's harder for the reader to mentally perform the double negation and
 understand what the code means.

```
if (socket.isConnected && database.hasData) {
 socket.write(database.read());
}
```
```
if (!socket.isDisconnected && !database.isEmpty) {
 socket.write(database.read());
}
```
 For some properties, there is no obvious positive form. Is a document that has
 been flushed to disk "saved" or "*un*-changed"? Is a document that *hasn't*
 been
 flushed "*un*-saved" or "changed"? In ambiguous cases, lean towards the choice
 that is less likely to be negated by users or has the shorter name.

 **Exception:** With some properties, the negative form is what users
 overwhelmingly need to use. Choosing the positive case would force them to
 negate the property with `!` everywhere. Instead, it may be better to use the
 negative case for that property.

### PREFER an imperative verb phrase for a function or method whose main purpose is a side effect

#Callable members can return a result to the caller and perform other work or side effects. In an imperative language like Dart, members are often called mainly for their side effect: they may change an object's internal state, produce some output, or talk to the outside world.

Those kinds of members should be named using an imperative verb phrase that clarifies the work the member performs.

```
list.add('element');
queue.removeFirst();
window.refresh();
```
This way, an invocation reads like a command to do that work.

### PREFER a noun phrase or non-imperative verb phrase for a function or method if returning a value is its primary purpose

#
 Other callable members have few side effects but return a useful result to the
 caller. If the member needs no parameters to do that, it should generally be a
 getter. But sometimes a logical "property" needs some parameters. For example,
 `elementAt()` returns a piece of data from a collection, but it needs a
 parameter to know *which* piece of data to return.

 This means the member is *syntactically* a method, but *conceptually* it is a
 property, and should be named as such using a phrase that describes *what*
 the
 member returns.

```
var element = list.elementAt(3);
var first = list.firstWhere(test);
var char = string.codeUnitAt(4);
```
 This guideline is deliberately softer than the previous one. Sometimes a method
 has no side effects but is still simpler to name with a verb phrase like
 `list.take()` or `string.split()`.

### CONSIDER an imperative verb phrase for a function or method if you want to draw attention to the work it performs

#When a member produces a result without any side effects, it should usually be a getter or a method with a noun phrase name describing the result it returns. However, sometimes the work required to produce that result is important. It may be prone to runtime failures, or use heavyweight resources like networking or file I/O. In cases like this, where you want the caller to think about the work the member is doing, give the member a verb phrase name that describes that work.

```
var table = database.downloadData();
var packageVersions = packageGraph.solveConstraints();
```
 Note, though, that this guideline is softer than the previous two. The work an
 operation performs is often an implementation detail that isn't relevant to the
 caller, and performance and robustness boundaries change over time. Most of the
 time, name your members based on *what* they do for the caller, not *how*
 they
 do it.

### AVOID starting a function or method name with `get`

 #

 In most cases, the method or function should be
 a getter with `get` removed from the name.
 For example, instead of a method named `getBreakfastOrder()`,
 define a getter named `breakfastOrder`.

 Even if the member does need to be a method because it takes arguments or
 otherwise isn't a good fit for a getter, you should still avoid `get`. Like the
 previous guidelines state, either:

-
 Simply drop `get`and use a noun phrase name like`breakfastOrder()`if the caller mostly cares about the value the method returns.
-
 Use a verb phrase name if the caller cares about the work being done, but pick a verb that more precisely describes the work than `get`, like`create`,`download`,`fetch`,`calculate`,`request`,`aggregate`, etc.

### PREFER naming a method `to___()` if it copies the object's state to a new object

 #
 Linter rule: use_to_and_as_if_applicable

 A *conversion* method is one that returns a new object containing a copy of
 almost all of the state of the receiver but usually in some different form or
 representation. The core libraries have a convention that these methods are
 named starting with `to` followed by the kind of result.

If you define a conversion method, it's helpful to follow that convention.

```
list.toSet();
stackTrace.toString();
dateTime.toLocal();
```
###
 PREFER naming a method `as___()` if it returns a different representation backed by the original object

 #
 Linter rule: use_to_and_as_if_applicable

 Conversion methods are "snapshots". The resulting object has its own copy of the
 original object's state. There are other conversion-like methods that return
 *views*—they provide a new object, but that object refers back to the
 original. Later changes to the original object are reflected in the view.

The core library convention for you to follow is `as___()`.

```
var map = table.asMap();
var list = bytes.asFloat32List();
var future = subscription.asFuture();
```
### AVOID describing the parameters in the function's or method's name

#The user will see the argument at the call site, so it usually doesn't help readability to also refer to it in the name itself.

```
list.add(element);
map.remove(key);
```
```
list.addElement(element)
map.removeKey(key)
```
However, it can be useful to mention a parameter to disambiguate it from other similarly-named methods that take different types:

```
map.containsKey(key);
map.containsValue(value);
```
### DO follow existing mnemonic conventions when naming type parameters

#Single letter names aren't exactly illuminating, but almost all generic types use them. Fortunately, they mostly use them in a consistent, mnemonic way. The conventions are:

-
 `E`for the**element**type in a collection:gooddart`class IterableBase<E> {} class List<E> {} class HashSet<E> {} class RedBlackTree<E> {}`
-
 `K`and`V`for the**key**and**value**types in an associative collection:gooddart`class Map<K, V> {} class Multimap<K, V> {} class MapEntry<K, V> {}`
-
 `R`for a type used as the**return**type of a function or a class's methods. This isn't common, but appears in typedefs sometimes and in classes that implement the visitor pattern:gooddart`abstract class ExpressionVisitor<R> { R visitBinary(BinaryExpression node); R visitLiteral(LiteralExpression node); R visitUnary(UnaryExpression node); }`
-
 Otherwise, use `T`,`S`, and`U`for generics that have a single type parameter and where the surrounding type makes its meaning obvious. There are multiple letters here to allow nesting without shadowing a surrounding name. For example:gooddart`class Future<T> { Future<S> then<S>(FutureOr<S> onValue(T value)) => ... }`Here, the generic method `then<S>()`uses`S`to avoid shadowing the`T`on`Future<T>`.

If none of the above cases are a good fit, then either another single-letter mnemonic name or a descriptive name is fine:

```
class Graph<N, E> {
 final List<N> nodes = [];
 final List<E> edges = [];
}
class Graph<Node, Edge> {
 final List<Node> nodes = [];
 final List<Edge> edges = [];
}
```
In practice, the existing conventions cover most type parameters.

## Libraries

#
 A leading underscore character ( `_` ) indicates that a member is private to its
 library. This is not mere convention, but is built into the language itself.

### PREFER making declarations private

#A public declaration in a library—either top level or in a class—is a signal that other libraries can and should access that member. It is also a commitment on your library's part to support that and behave properly when it happens.

 If that's not what you intend, add the little `_` and be happy. Narrow public
 interfaces are easier for you to maintain and easier for users to learn. As a
 nice bonus, the analyzer will tell you about unused private declarations so you
 can delete dead code. It can't do that if the member is public because it
 doesn't know if any code outside of its view is using it.

### CONSIDER declaring multiple classes in the same library

#Some languages, such as Java, tie the organization of files to the organization of classes—each file may only define a single top level class. Dart does not have that limitation. Libraries are distinct entities separate from classes. It's perfectly fine for a single library to contain multiple classes, top level variables, and functions if they all logically belong together.

Placing multiple classes together in one library can enable some useful patterns. Since privacy in Dart works at the library level, not the class level, this is a way to define "friend" classes like you might in C++. Every class declared in the same library can access each other's private members, but code outside of that library cannot.

 Of course, this guideline doesn't mean you *should* put all of your classes into
 a huge monolithic library, just that you are allowed to place more than one
 class in a single library.

## Classes and mixins

#Dart is a "pure" object-oriented language in that all objects are instances of classes. But Dart does not require all code to be defined inside a class—you can define top-level variables, constants, and functions like you can in a procedural or functional language.

### AVOID defining a one-member abstract class when a simple function will do

#Linter rule: one_member_abstracts

 Unlike Java, Dart has first-class functions, closures, and a nice light syntax
 for using them. If all you need is something like a callback, just use a
 function. If you're defining a class and it only has a single abstract member
 with a meaningless name like `call` or `invoke`, there is a good chance you
 just want a function.

```
typedef Predicate<E> = bool Function(E element);
```
```
abstract class Predicate<E> {
 bool test(E element);
}
```
### AVOID defining a class that contains only static members

#Linter rule: avoid_classes_with_only_static_members

 In Java and C#, every definition *must* be inside a class, so it's common to see
 "classes" that exist only as a place to stuff static members. Other classes are
 used as namespaces—a way to give a shared prefix to a bunch of members to
 relate them to each other or avoid a name collision.

 Dart has top-level functions, variables, and constants, so you don't *need* a
 class just to define something. If what you want is a namespace, a library is a
 better fit. Libraries support import prefixes and show/hide combinators. Those
 are powerful tools that let the consumer of your code handle name collisions in
 the way that works best for *them*.

If a function or variable isn't logically tied to a class, put it at the top level. If you're worried about name collisions, give it a more precise name or move it to a separate library that can be imported with a prefix.

```
DateTime mostRecent(List<DateTime> dates) {
 return dates.reduce((a, b) => a.isAfter(b) ? a : b);
}
const _favoriteMammal = 'weasel';
```
```
class DateUtils {
 static DateTime mostRecent(List<DateTime> dates) {
 return dates.reduce((a, b) => a.isAfter(b) ? a : b);
 }
}
class _Favorites {
 static const mammal = 'weasel';
}
```
 In idiomatic Dart, classes define *kinds of objects*. A type that is never
 instantiated is a code smell.

However, this isn't a hard rule. For example, with constants and enum-like types, it may be natural to group them in a class.

```
class Color {
 static const red = '#f00';
 static const green = '#0f0';
 static const blue = '#00f';
 static const black = '#000';
 static const white = '#fff';
}
```
### AVOID extending a class that isn't intended to be subclassed

#
 If a constructor is changed from a generative constructor to a factory
 constructor, any subclass constructor calling that constructor will break.
 Also, if a class changes which of its own methods it invokes on `this`, that
 may break subclasses that override those methods and expect them to be called
 at certain points.

 Both of these mean that a class needs to be deliberate about whether or not it
 wants to allow subclassing. This can be communicated in a doc comment, or by
 giving the class an obvious name like `IterableBase`. If the author of the class
 doesn't do that, it's best to assume you should *not* extend the class.
 Otherwise, later changes to it may break your code.

### DO use class modifiers to control if your class can be extended

#
 Class modifiers like `final`, `interface`, or `sealed`
 restrict how a class can be extended.
 For example, use `final class A {}` or `interface class B {}`
 to prevent
 extension outside the current library.
 Use these modifiers to communicate your intent, rather than relying on documentation.

### AVOID implementing a class that isn't intended to be an interface

#Implicit interfaces are a powerful tool in Dart to avoid having to repeat the contract of a class when it can be trivially inferred from the signatures of an implementation of that contract.

 But implementing a class's interface is a very tight coupling to that class. It
 means virtually *any* change to the class whose interface you are implementing
 will break your implementation. For example, adding a new member to a class is
 usually a safe, non-breaking change. But if you are implementing that class's
 interface, now your class has a static error because it lacks an implementation
 of that new method.

Library maintainers need the ability to evolve existing classes without breaking users. If you treat every class like it exposes an interface that users are free to implement, then changing those classes becomes very difficult. That difficulty in turn means the libraries you rely on are slower to grow and adapt to new needs.

To give the authors of the classes you use more leeway, avoid implementing implicit interfaces except for classes that are clearly intended to be implemented. Otherwise, you may introduce a coupling that the author doesn't intend, and they may break your code without realizing it.

### DO use class modifiers to control if your class can be an interface

#
 When designing a library, use class modifiers like `final`, `base`, or `sealed`
 to enforce intended
 usage. For example, use `final class C {}` or `base class D {}`
 to prevent
 implementation outside the current library.
 While it's ideal for all libraries to use these modifiers to enforce design intent,
 developers may still encounter cases where they aren't applied. In such cases, be mindful of
 unintended implementation issues.

### PREFER defining a pure `mixin` or pure `class` to a `mixin class`

 #
 Linter rule: prefer_mixin

Dart previously (in language versions 2.12 to 2.19) allowed any class that met certain restrictions (no non-default constructor, no superclass, etc.) to be mixed into other classes. This was confusing because the author of the class might not have intended it to be mixed in.

 Dart 3.0.0 now requires that any type intended to be mixed into other classes,
 as well as treated as a normal class, must be explicitly declared as such with
 the `mixin class` declaration.

 Types that need to be both a mixin and a class should be a rare case, however.
 The `mixin class` declaration is mostly meant to help migrate pre-3.0.0 classes
 being used as mixins to a more explicit declaration. New code should clearly
 define the behavior and intention of its declarations by using only pure `mixin`

 or pure `class` declarations, and avoid the ambiguity of mixin classes.

 Read Migrating classes as mixins

 for more guidance on `mixin` and `mixin class` declarations.

## Constructors

#
 Dart constructors are created by declaring a function with the same name as the
 class and, optionally, an additional identifier. The latter are called *named
 constructors*.

### CONSIDER making your constructor `const` if the class supports it

 #

 If you have a class where all the fields are final, and the constructor does
 nothing but initialize them, you can make that constructor `const`. That lets
 users create instances of your class in places where constants are
 required—inside other larger constants, switch cases, default parameter
 values, etc.

If you don't explicitly make it `const`, they aren't able to do that.

 Note, however, that a `const` constructor is a commitment in your public API. If
 you later change the constructor to non-`const`, it will break users that are
 calling it in constant expressions. If you don't want to commit to that, don't
 make it `const`. In practice, `const` constructors are most useful for simple,
 immutable value-like types.

## Members

#A member belongs to an object and can be either methods or instance variables.

### PREFER making fields and top-level variables `final`

 #
 Linter rule: prefer_final_fields

 State that is not *mutable*—that does not change over time—is
 easier for programmers to reason about. Classes and libraries that minimize the
 amount of mutable state they work with tend to be easier to maintain.
 Of course, it is often useful to have mutable data. But, if you don't need it,
 your default should be to make fields and top-level variables `final`
 when you
 can.

 Sometimes an instance field doesn't change after it has been initialized, but
 can't be initialized until after the instance is constructed. For example, it
 may need to reference `this` or some other field on the instance. In cases like
 that, consider making the field `late final`. When you do, you may also be able
 to initialize the field at its declaration.

### DO use getters for operations that conceptually access properties

#
 Deciding when a member should be a getter versus a method is a subtle but
 important part of good API design, hence this very long guideline.
 Some other language's cultures shy away from getters. They only use them when
 the operation is almost exactly like a field—it does a minuscule amount of
 calculation on state that lives entirely on the object. Anything more complex or
 heavyweight than that gets `()` after the name to signal "computation goin' on
 here!" because a bare name after a `.` means "field".

 Dart is *not* like that. In Dart, *all* dotted names are member invocations that
 may do computation. Fields are special—they're getters whose
 implementation is provided by the language. In other words, getters are not
 "particularly slow fields" in Dart; fields are "particularly fast getters".

 Even so, choosing a getter over a method sends an important signal to the
 caller. The signal, roughly, is that the operation is "field-like". The
 operation, at least in principle, *could* be implemented using a field, as far
 as the caller knows. That implies:

- **The operation does not take any arguments and returns a result.**
-
 **The caller cares mostly about the result.**If you want the caller to worry about*how*the operation produces its result more than they do the result being produced, then give the operation a verb name that describes the work and make it a method.This does *not*mean the operation has to be particularly fast in order to be a getter.`IterableBase.length`is`O(n)`, and that's OK. It's fine for a getter to do significant calculation. But if it does a*surprising*amount of work, you may want to draw their attention to that by making it a method whose name is a verb describing what it does.baddart`connection.nextIncomingMessage; // Does network I/O. expression.normalForm; // Could be exponential to calculate.`
-
 **The operation does not have user-visible side effects.**Accessing a real field does not alter the object or any other state in the program. It doesn't produce output, write files, etc. A getter shouldn't do those things either.The "user-visible" part is important. It's fine for getters to modify hidden state or produce out of band side effects. Getters can lazily calculate and store their result, write to a cache, log stuff, etc. As long as the caller doesn't *care*about the side effect, it's probably fine.baddart`stdout.newline; // Produces output. list.clear; // Modifies object.`
-
 **The operation is**"Idempotent" is an odd word that, in this context, basically means that calling the operation multiple times produces the same result each time, unless some state is explicitly modified between those calls. (Obviously,*idempotent*.`list.length`produces different results if you add an element to the list between calls.)"Same result" here does not mean a getter must literally produce an identical object on successive calls. Requiring that would force many getters to have brittle caching, which negates the whole point of using a getter. It's common, and perfectly fine, for a getter to return a new future or list each time you call it. The important part is that the future completes to the same value, and the list contains the same elements. In other words, the result value should be the same *in the aspects that the caller cares about.*baddart`DateTime.now; // New result each time.`
-
 **The resulting object doesn't expose all of the original object's state.**A field exposes only a piece of an object. If your operation returns a result that exposes the original object's entire state, it's likely better off as a`to___()`or`as___()`method.

If all of the above describe your operation, it should be a getter. It seems like few members would survive that gauntlet, but surprisingly many do. Many operations just do some computation on some state and most of those can and should be getters.

```
rectangle.area;
collection.isEmpty;
button.canShow;
dataSet.minimumValue;
```
### DO use setters for operations that conceptually change properties

#Linter rule: use_setters_to_change_properties

Deciding between a setter versus a method is similar to deciding between a getter versus a method. In both cases, the operation should be "field-like".

For a setter, "field-like" means:

-
 **The operation takes a single argument and does not produce a result value.**
- **The operation changes some state in the object.**
-
 **The operation is idempotent.**Calling the same setter twice with the same value should do nothing the second time as far as the caller is concerned. Internally, maybe you've got some cache invalidation or logging going on. That's fine. But from the caller's perspective, it appears that the second call does nothing.

```
rectangle.width = 3;
button.visible = false;
```
### DON'T define a setter without a corresponding getter

#Linter rule: avoid_setters_without_getters

 Users think of getters and setters as visible properties of an object. A
 "dropbox" property that can be written to but not seen is confusing and
 confounds their intuition about how properties work. For example, a setter
 without a getter means you can use `=` to modify it, but not `+=`.

 This guideline does *not* mean you should add a getter just to permit the setter
 you want to add. Objects shouldn't generally expose more state than they need
 to. If you have some piece of an object's state that can be modified but not
 exposed in the same way, use a method instead.

### AVOID using runtime type tests to fake overloading

#
 It's common for an API to support similar operations
 on different types of parameters.
 To emphasize the similarity, some languages support *overloading*,
 which lets you define multiple methods
 that have the same name but different parameter lists.
 At compile time, the compiler looks at the actual argument types to determine
 which method to call.

 Dart doesn't have overloading.
 You can define an API that looks like overloading
 by defining a single method and then using `is` type tests
 inside the body to look at the runtime types of the arguments and perform the
 appropriate behavior.
 However, faking overloading this way turns a *compile time* method selection
 into a choice that happens at *runtime*.

If callers usually know which type they have and which specific operation they want, it's better to define separate methods with different names to let callers select the right operation. This gives better static type checking and faster performance since it avoids any runtime type tests.

 However, if users might have an object of an unknown type
 and *want* the API to internally use `is` to pick the right operation,
 then a single method where the parameter is a supertype
 of all of the supported types might be reasonable.

### AVOID public `late final` fields without initializers

 #

 Unlike other `final` fields, a `late final` field without an initializer *does*

 define a setter. If that field is public, then the setter is public. This is
 rarely what you want. Fields are usually marked `late` so that they can be
 initialized *internally* at some point in the instance's lifetime, often inside
 the constructor body.

 Unless you *do* want users to call the setter, it's better to pick one of the
 following solutions:

- Don't use `late`.
- Use a factory constructor to compute the `final`field values.
- Use `late`, but initialize the`late`field at its declaration.
-
 Use `late`, but make the`late`field private and define a public getter for it.

### AVOID returning nullable `Future`, `Stream`, and collection types

 #

 When an API returns a container type, it has two ways to indicate the absence of
 data: It can return an empty container or it can return `null`. Users generally
 assume and prefer that you use an empty container to indicate "no data". That
 way, they have a real object that they can call methods on like `isEmpty`.

To indicate that your API has no data to provide, prefer returning an empty collection, a non-nullable future of a nullable type, or a stream that doesn't emit any values.

 **Exception:** If returning `null` *means something different* from yielding an
 empty container, it might make sense to use a nullable type.

### AVOID returning `this` from methods just to enable a fluent interface

 #
 Linter rule: avoid_returning_this

Method cascades are a better solution for chaining method calls.

```
var buffer =
 StringBuffer()
 ..write('one')
 ..write('two')
 ..write('three');
```
```
var buffer =
 StringBuffer()
 .write('one')
 .write('two')
 .write('three');
```
## Types

#
 When you write down a type in your program, you constrain the kinds of values
 that flow into different parts of your code. Types can appear in two kinds of
 places: *type annotations* on declarations and type arguments to *generic
 invocations*.

 Type annotations are what you normally think of when you think of "static
 types". You can type annotate a variable, parameter, field, or return type. In
 the following example, `bool` and `String` are type annotations. They hang off
 the static declarative structure of the code and aren't "executed" at runtime.

```
bool isEmpty(String parameter) {
 bool result = parameter.isEmpty;
 return result;
}
```
 A generic invocation is a collection literal, a call to a generic class's
 constructor, or an invocation of a generic method. In the next example, `num`

 and `int` are type arguments on generic invocations. Even though they are types,
 they are first-class entities that get reified and passed to the invocation at
 runtime.

```
var lists = <num>[1, 2];
lists.addAll(List<num>.filled(3, 4));
lists.cast<int>();
```
 We stress the "generic invocation" part here, because type arguments can *also*
 appear in type annotations:

```
List<int> ints = [1, 2];
```
 Here, `int` is a type argument, but it appears inside a type annotation, not a
 generic invocation. You usually don't need to worry about this distinction, but
 in a couple of places, we have different guidance for when a type is used in a
 generic invocation as opposed to a type annotation.

#### Type inference

#
 Type annotations are optional in Dart.
 If you omit one, Dart tries to infer a type
 based on the nearby context. Sometimes it doesn't have enough information to
 infer a complete type. When that happens, Dart sometimes reports an error, but
 usually silently fills in any missing parts with `dynamic`. The implicit
 `dynamic` leads to code that *looks* inferred and safe, but actually disables
 type checking completely. The rules below avoid that by requiring types when
 inference fails.

 The fact that Dart has both type inference and a `dynamic` type leads to some
 confusion about what it means to say code is "untyped". Does that mean the code
 is dynamically typed, or that you didn't *write* the type? To avoid that
 confusion, we avoid saying "untyped" and instead use the following terminology:

-
 If the code is *type annotated*, the type was explicitly written in the code.
-
 If the code is *inferred*, no type annotation was written, and Dart successfully figured out the type on its own. Inference can fail, in which case the guidelines don't consider that inferred.
-
 If the code is *dynamic*, then its static type is the special`dynamic`type. Code can be explicitly annotated`dynamic`or it can be inferred.

 In other words, whether some code is annotated or inferred is orthogonal to
 whether it is `dynamic` or some other type.

Inference is a powerful tool to spare you the effort of writing and reading types that are obvious or uninteresting. It keeps the reader's attention focused on the behavior of the code itself. Explicit types are also a key part of robust, maintainable code. They define the static shape of an API and create boundaries to document and enforce what kinds of values are allowed to reach different parts of the program.

Of course, inference isn't magic. Sometimes inference succeeds and selects a type, but it's not the type you want. The common case is inferring an overly precise type from a variable's initializer when you intend to assign values of other types to the variable later. In those cases, you have to write the type explicitly.

The guidelines here strike the best balance we've found between brevity and control, flexibility and safety. There are specific guidelines to cover all the various cases, but the rough summary is:

-
 Do annotate when inference doesn't have enough context, even when `dynamic`is the type you want.
- Don't annotate locals and generic invocations unless you need to.
-
 Prefer annotating top-level variables and fields unless the initializer makes the type obvious.

### DO type annotate variables without initializers

#Linter rule: prefer_typing_uninitialized_variables

The type of a variable—top-level, local, static field, or instance field—can often be inferred from its initializer. However, if there is no initializer, inference fails.

```
List<AstNode> parameters;
if (node is Constructor) {
 parameters = node.signature;
} else if (node is Method) {
 parameters = node.parameters;
}
```
```
var parameters;
if (node is Constructor) {
 parameters = node.signature;
} else if (node is Method) {
 parameters = node.parameters;
}
```
### DO type annotate fields and top-level variables if the type isn't obvious

#Linter rule: type_annotate_public_apis

Type annotations are important documentation for how a library should be used. They form boundaries between regions of a program to isolate the source of a type error. Consider:

```
install(id, destination) => ...
```
 Here, it's unclear what `id` is. A string? And what is `destination`? A string
 or a `File` object? Is this method synchronous or asynchronous? This is clearer:

```
Future<bool> install(PackageId id, String destination) => ...
```
In some cases, though, the type is so obvious that writing it is pointless:

```
const screenWidth = 640; // Inferred as int.
```
"Obvious" isn't precisely defined, but these are all good candidates:

- Literals.
- Constructor invocations.
- References to other constants that are explicitly typed.
- Simple expressions on numbers and strings.
-
 Factory methods like `int.parse()`,`Future.wait()`, etc. that readers are expected to be familiar with.

If you think the initializer expression—whatever it is—is sufficiently clear, then you may omit the annotation. But if you think annotating helps make the code clearer, then add one.

 When in doubt, add a type annotation.
 You can also explicitly annotate obvious types.
 If the inferred type relies on values or declarations from other libraries,
 you might want to type annotate *your*
 declaration so that a change to that other library doesn't silently change the
 type of your own API without you realizing.

 This rule applies to both public and private declarations. Just as type
 annotations on APIs help *users* of your code, types on private members help
 *maintainers*.

### DON'T redundantly type annotate initialized local variables

#Linter rule: omit_local_variable_types

 Local variables, especially in modern code where functions tend to be small,
 have very little scope. Omitting the type focuses the reader's attention on the
 more important *name* of the variable and its initialized value.

```
List<List<Ingredient>> possibleDesserts(Set<Ingredient> pantry) {
 var desserts = <List<Ingredient>>[];
 for (final recipe in cookbook) {
 if (pantry.containsAll(recipe)) {
 desserts.add(recipe);
 }
 }
 return desserts;
}
```
```
List<List<Ingredient>> possibleDesserts(Set<Ingredient> pantry) {
 List<List<Ingredient>> desserts = <List<Ingredient>>[];
 for (final List<Ingredient> recipe in cookbook) {
 if (pantry.containsAll(recipe)) {
 desserts.add(recipe);
 }
 }
 return desserts;
}
```
Sometimes the inferred type is not the type you want the variable to have. For example, you may intend to assign values of other types later. In that case, annotate the variable with the type you want.

```
Widget build(BuildContext context) {
 Widget result = Text('You won!');
 if (applyPadding) {
 result = Padding(padding: EdgeInsets.all(8.0), child: result);
 }
 return result;
}
```
### DO annotate return types on function declarations

#Dart doesn't generally infer the return type of a function declaration from its body, unlike some other languages. That means you should write a type annotation for the return type yourself.

```
String makeGreeting(String who) {
 return 'Hello, $who!';
}
```
```
makeGreeting(String who) {
 return 'Hello, $who!';
}
```
 Note that this guideline only applies to *non-local* function declarations:
 top-level, static, and instance methods and getters. Local functions and
 anonymous function expressions infer a return type from their body. In fact, the
 anonymous function syntax doesn't even allow a return type annotation.

### DO annotate parameter types on function declarations

#A function's parameter list determines its boundary to the outside world. Annotating parameter types makes that boundary well defined. Note that even though default parameter values look like variable initializers, Dart doesn't infer an optional parameter's type from its default value.

```
void sayRepeatedly(String message, {int count = 2}) {
 for (var i = 0; i < count; i++) {
 print(message);
 }
}
```
```
void sayRepeatedly(message, {count = 2}) {
 for (var i = 0; i < count; i++) {
 print(message);
 }
}
```
 **Exception:** Function expressions and initializing formals have
 different type annotation conventions, as described in the next two guidelines.

### DON'T annotate inferred parameter types on function expressions

#Linter rule: avoid_types_on_closure_parameters

 Anonymous functions are almost always immediately passed to a method taking a
 callback of some type.
 When a function expression is created in a typed context,
 Dart tries to infer the function's parameter types based on the expected type.
 For example, when you pass a function expression to `Iterable.map()`, your
 function's parameter type is inferred based on the type of callback that `map()`

 expects:

```
var names = people.map((person) => person.name);
```
```
var names = people.map((Person person) => person.name);
```
If the language is able to infer the type you want for a parameter in a function expression, then don't annotate. In rare cases, the surrounding context isn't precise enough to provide a type for one or more of the function's parameters. In those cases, you may need to annotate. (If the function isn't used immediately, it's usually better to make it a named declaration.)

### DON'T type annotate initializing formals

#Linter rule: type_init_formals

 If a constructor parameter is using `this.` to initialize a field,
 or `super.` to forward a super parameter,
 then the type of the parameter
 is inferred to have the same type as
 the field or super-constructor parameter respectively.

```
class Point {
 double x, y;
 Point(this.x, this.y);
}
class MyWidget extends StatelessWidget {
 MyWidget({super.key});
}
```
```
class Point {
 double x, y;
 Point(double this.x, double this.y);
}
class MyWidget extends StatelessWidget {
 MyWidget({Key? super.key});
}
```
### DO write type arguments on generic invocations that aren't inferred

#Dart is pretty smart about inferring type arguments in generic invocations. It looks at the expected type where the expression occurs and the types of values being passed to the invocation. However, sometimes those aren't enough to fully determine a type argument. In that case, write the entire type argument list explicitly.

```
var playerScores = <String, int>{};
final events = StreamController<Event>();
```
```
var playerScores = {};
final events = StreamController();
```
 Sometimes the invocation occurs as the initializer to a variable declaration. If
 the variable is *not* local, then instead of writing the type argument list on the
 invocation itself, you may put a type annotation on the declaration:

```
class Downloader {
 final Completer<String> response = Completer();
}
```
```
class Downloader {
 final response = Completer();
}
```
 Annotating the variable also addresses this guideline because now the type
 arguments *are* inferred.

### DON'T write type arguments on generic invocations that are inferred

#
 This is the converse of the previous rule. If an invocation's type argument list
 *is* correctly inferred with the types you want, then omit the types and let
 Dart do the work for you.

```
class Downloader {
 final Completer<String> response = Completer();
}
```
```
class Downloader {
 final Completer<String> response = Completer<String>();
}
```
Here, the type annotation on the field provides a surrounding context to infer the type argument of constructor call in the initializer.

```
var items = Future.value([1, 2, 3]);
```
```
var items = Future<List<int>>.value(<int>[1, 2, 3]);
```
Here, the types of the collection and instance can be inferred bottom-up from their elements and arguments.

### AVOID writing incomplete generic types

#The goal of writing a type annotation or type argument is to pin down a complete type. However, if you write the name of a generic type but omit its type arguments, you haven't fully specified the type. In Java, these are called "raw types". For example:

```
List numbers = [1, 2, 3];
var completer = Completer<Map>();
```
 Here, `numbers` has a type annotation, but the annotation doesn't provide a type
 argument to the generic `List`. Likewise, the `Map` type argument to
 `Completer`
 isn't fully specified. In cases like this, Dart will *not* try to "fill in" the
 rest of the type for you using the surrounding context. Instead, it silently
 fills in any missing type arguments with `dynamic` (or the bound if the
 class has one). That's rarely what you want.

Instead, if you're writing a generic type either in a type annotation or as a type argument inside some invocation, make sure to write a complete type:

```
List<num> numbers = [1, 2, 3];
var completer = Completer<Map<String, int>>();
```
### DO annotate with `dynamic` instead of letting inference fail

 #

 When inference doesn't fill in a type, it usually defaults to `dynamic`. If
 `dynamic` is the type you want, this is technically the most terse way to get
 it. However, it's not the most *clear* way. A casual reader of your code who
 sees that an annotation is missing has no way of knowing if you intended it to be
 `dynamic`, expected inference to fill in some other type, or simply forgot to
 write the annotation.

 When `dynamic` is the type you want, write that explicitly to make your intent
 clear and highlight that this code has less static safety.

```
dynamic mergeJson(dynamic original, dynamic changes) => ...
```
```
mergeJson(original, changes) => ...
```
Note that it's OK to omit the type when Dart *successfully* infers `dynamic`.

```
Map<String, dynamic> readJson() => ...
void printUsers() {
 var json = readJson();
 var users = json['users'];
 print(users);
}
```
 Here, Dart infers `Map<String, dynamic>` for `json` and then from that infers
 `dynamic` for `users`. It's fine to leave `users`
 without a type annotation. The
 distinction is a little subtle. It's OK to allow inference to *propagate*

 `dynamic` through your code from a `dynamic` type annotation somewhere else, but
 you don't want it to inject a `dynamic` type annotation in a place where your
 code did not specify one.

**Exception**: Type annotations on unused parameters (`_`) can be omitted.

### PREFER signatures in function type annotations

#
 The identifier `Function` by itself without any return type or parameter
 signature refers to the special Function
 type. This type is only
 marginally more useful than using `dynamic`. If you're going to annotate, prefer
 a full function type that includes the parameters and return type of the
 function.

```
bool isValid(String value, bool Function(String) test) => ...
```
```
bool isValid(String value, Function test) => ...
```
 **Exception:** Sometimes, you want a type that represents the union of multiple
 different function types. For example, you may accept a function that takes one
 parameter or a function that takes two. Since we don't have union types, there's
 no way to precisely type that and you'd normally have to use `dynamic`.
 `Function` is at least a little more helpful than that:

```
void handleError(void Function() operation, Function errorHandler) {
 try {
 operation();
 } catch (err, stack) {
 if (errorHandler is Function(Object)) {
 errorHandler(err);
 } else if (errorHandler is Function(Object, StackTrace)) {
 errorHandler(err, stack);
 } else {
 throw ArgumentError('errorHandler has wrong signature.');
 }
 }
}
```
### DON'T specify a return type for a setter

#Linter rule: avoid_return_types_on_setters

Setters always return `void` in Dart. Writing the word is pointless.

```
void set foo(Foo value) {
 ...
}
```
```
set foo(Foo value) {
 ...
}
```
### DON'T use the legacy typedef syntax

#Linter rule: prefer_generic_function_type_aliases

Dart has two notations for defining a named typedef for a function type. The original syntax looks like:

```
typedef int Comparison<T>(T a, T b);
```
That syntax has a couple of problems:

-
 There is no way to assign a name to a *generic*function type. In the above example, the typedef itself is generic. If you reference`Comparison`in your code, without a type argument, you implicitly get the function type`int Function(dynamic, dynamic)`,*not*`int Function<T>(T, T)`. This doesn't come up in practice often, but it matters in certain corner cases.
-
 A single identifier in a parameter is interpreted as the parameter's *name*, not its*type*. Given:baddart`typedef bool TestNumber(num);`Most users expect this to be a function type that takes a `num`and returns`bool`. It is actually a function type that takes*any*object (`dynamic`) and returns`bool`. The parameter's*name*(which isn't used for anything except documentation in the typedef) is "num". This has been a long-standing source of errors in Dart.

The new syntax looks like this:

```
typedef Comparison<T> = int Function(T, T);
```
If you want to include a parameter's name, you can do that too:

```
typedef Comparison<T> = int Function(T a, T b);
```
 The new syntax can express anything the old syntax could express and more, and
 lacks the error-prone misfeature where a single identifier is treated as the
 parameter's name instead of its type. The same function type syntax after the
 `=` in the typedef is also allowed anywhere a type annotation may appear, giving
 us a single consistent way to write function types anywhere in a program.

The old typedef syntax is still supported to avoid breaking existing code, but it's deprecated.

### PREFER inline function types over typedefs

#Linter rule: avoid_private_typedef_functions

In Dart, if you want to use a function type for a field, variable, or generic type argument, you can define a typedef for the function type. However, Dart supports an inline function type syntax that can be used anywhere a type annotation is allowed:

```
class FilteredObservable {
 final bool Function(Event) _predicate;
 final List<void Function(Event)> _observers;
 FilteredObservable(this._predicate, this._observers);
 void Function(Event)? notify(Event event) {
 if (!_predicate(event)) return null;
 void Function(Event)? last;
 for (final observer in _observers) {
 observer(event);
 last = observer;
 }
 return last;
 }
}
```
It may still be worth defining a typedef if the function type is particularly long or frequently used. But in most cases, users want to see what the function type actually is right where it's used, and the function type syntax gives them that clarity.

### PREFER using function type syntax for parameters

#Linter rule: use_function_type_syntax_for_parameters

Dart has a special syntax when defining a parameter whose type is a function. Sort of like in C, you surround the parameter's name with the function's return type and parameter signature:

```
Iterable<T> where(bool predicate(T element)) => ...
```
Before Dart added function type syntax, this was the only way to give a parameter a function type without defining a typedef. Now that Dart has a general notation for function types, you can use it for function-typed parameters as well:

```
Iterable<T> where(bool Function(T) predicate) => ...
```
The new syntax is a little more verbose, but is consistent with other locations where you must use the new syntax.

### AVOID using `dynamic` unless you want to disable static checking

 #

 Some operations work with any possible object. For example, a `log()` method
 could take any object and call `toString()` on it. Two types in Dart permit all
 values: `Object?` and `dynamic`. However, they convey different things. If you
 simply want to state that you allow all objects, use `Object?`. If you want to
 allow all objects *except* `null`, then use `Object`.

 The type `dynamic` not only accepts all objects, but it also permits all
 *operations*. Any member access on a value of type `dynamic` is allowed at
 compile time, but may fail and throw an exception at runtime. If you want
 exactly that risky but flexible dynamic dispatch, then `dynamic` is the right
 type to use.

 Otherwise, prefer using `Object?` or `Object`. Rely on `is` checks and type
 promotion to
 ensure that the value's runtime type supports the member you want to access
 before you access it.

```
/// Returns a Boolean representation for [arg], which must
/// be a String or bool.
bool convertToBool(Object arg) {
 if (arg is bool) return arg;
 if (arg is String) return arg.toLowerCase() == 'true';
 throw ArgumentError('Cannot convert $arg to a bool.');
}
```
 The main exception to this rule is when working with existing APIs that use
 `dynamic`, especially inside a generic type. For example, JSON objects have type
 `Map<String, dynamic>` and your code will need to accept that same type. Even
 so, when using a value from one of these APIs, it's often a good idea to cast it
 to a more precise type before accessing members.

###
 DO use `Future<void>` as the return type of asynchronous members that do not produce values

 #

 When you have a synchronous function that doesn't return a value, you use `void`
 as the return type. The asynchronous equivalent for a method that doesn't
 produce a value, but that the caller might need to await, is `Future<void>`.

 You may see code that uses `Future` or `Future<Null>` instead because older
 versions of Dart didn't allow `void` as a type argument. Now that it does, you
 should use it. Doing so more directly matches how you'd type a similar
 synchronous function, and gives you better error-checking for callers and in the
 body of the function.

 For asynchronous functions that do not return a useful value and where no
 callers need to await the asynchronous work or handle an asynchronous failure,
 use a return type of `void`.

### AVOID using `FutureOr<T>` as a return type

 #

 If a method accepts a `FutureOr<int>`, it is generous in what it
 accepts. Users can call the method with either an `int` or a
 `Future<int>`, so they don't need to wrap an `int` in
 `Future` that you are
 going to unwrap anyway.

 If you *return* a `FutureOr<int>`, users need to check whether get back an `int`

 or a `Future<int>` before they can do anything useful. (Or they'll just
 `await`
 the value, effectively always treating it as a `Future`.) Just return a
 `Future<int>`, it's cleaner. It's easier for users to understand that a function
 is either always asynchronous or always synchronous, but a function that can be
 either is hard to use correctly.

```
Future<int> triple(FutureOr<int> value) async => (await value) * 3;
```
```
FutureOr<int> triple(FutureOr<int> value) {
 if (value is int) return value * 3;
 return value.then((v) => v * 3);
}
```
 The more precise formulation of this guideline is to *only use FutureOr<T> in
 contravariant
 positions.* Parameters are contravariant and return types are
 covariant. In nested function types, this gets flipped—if you have a
 parameter whose type is itself a function, then the callback's return type is
 now in contravariant position and the callback's parameters are covariant. This
 means it's OK for a

*callback's*type to return

`FutureOr<T>`:
 ```
Stream<S> asyncMap<T, S>(
 Iterable<T> iterable,
 FutureOr<S> Function(T) callback,
) async* {
 for (final element in iterable) {
 yield await callback(element);
 }
}
```
## Parameters

#In Dart, optional parameters can be either positional or named, but not both.

### AVOID positional boolean parameters

#Linter rule: avoid_positional_boolean_parameters

 Unlike other types, booleans are usually used in literal form. Values like
 numbers are usually wrapped in named constants, but we typically pass around
 `true` and `false` directly. That can make call sites unreadable if it isn't
 clear what the boolean represents:

```
new Task(true);
new Task(false);
new ListBox(false, true, true);
new Button(false);
```
Instead, prefer using named arguments, named constructors, or named constants to clarify what the call is doing.

```
Task.oneShot();
Task.repeating();
ListBox(scroll: true, showScrollbars: true);
Button(ButtonState.enabled);
```
Note that this doesn't apply to setters, where the name makes it clear what the value represents:

```
listBox.canScroll = true;
button.isEnabled = false;
```
### AVOID optional positional parameters if the user may want to omit earlier parameters

#Optional positional parameters should have a logical progression such that earlier parameters are passed more often than later ones. Users should almost never need to explicitly pass a "hole" to omit an earlier positional argument to pass later one. You're better off using named arguments for that.

```
String.fromCharCodes(Iterable<int> charCodes, [int start = 0, int? end]);
DateTime(
 int year, [
 int month = 1,
 int day = 1,
 int hour = 0,
 int minute = 0,
 int second = 0,
 int millisecond = 0,
 int microsecond = 0,
]);
Duration({
 int days = 0,
 int hours = 0,
 int minutes = 0,
 int seconds = 0,
 int milliseconds = 0,
 int microseconds = 0,
});
```
### AVOID mandatory parameters that accept a special "no argument" value

#
 If the user is logically omitting a parameter, prefer letting them actually omit
 it by making the parameter optional instead of forcing them to pass `null`, an
 empty string, or some other special value that means "did not pass".

 Omitting the parameter is more terse and helps prevent bugs where a sentinel
 value like `null` is accidentally passed when the user thought they were
 providing a real value.

```
var rest = string.substring(start);
```
```
var rest = string.substring(start, null);
```
### DO use inclusive start and exclusive end parameters to accept a range

#If you are defining a method or function that lets a user select a range of elements or items from some integer-indexed sequence, take a start index, which refers to the first item and a (likely optional) end index which is one greater than the index of the last item.

This is consistent with core libraries that do the same thing.

```
[0, 1, 2, 3].sublist(1, 3) // [1, 2]
'abcd'.substring(1, 3) // 'bc'
```
It's particularly important to be consistent here because these parameters are usually unnamed. If your API takes a length instead of an end point, the difference won't be visible at all at the call site.

## Equality

#Implementing custom equality behavior for a class can be tricky. Users have deep intuition about how equality works that your objects need to match, and collection types like hash tables have subtle contracts that they expect elements to follow.

### DO override `hashCode` if you override `==`

 #
 Linter rule: hash_and_equals

 The default hash code implementation provides an *identity* hash—two
 objects generally only have the same hash code if they are the exact same
 object. Likewise, the default behavior for `==` is identity.

 If you are overriding `==`, it implies you may have different objects that are
 considered "equal" by your class. **Any two objects that are equal must have the
 same hash code.** Otherwise, maps and other hash-based collections will fail to
 recognize that the two objects are equivalent.

### DO make your `==` operator obey the mathematical rules of equality

 #
 An equivalence relation should be:

- **Reflexive**:- `a == a`should always return- `true`.
-
 **Symmetric**:`a == b`should return the same thing as`b == a`.
-
 **Transitive**: If`a == b`and`b == c`both return`true`, then`a == c`should too.

 Users and code that uses `==` expect all of these laws to be followed. If your
 class can't obey these rules, then `==` isn't the right name for the operation
 you're trying to express.

### AVOID defining custom equality for mutable classes

#Linter rule: avoid_equals_and_hash_code_on_mutable_classes

 When you define `==`, you also have to define `hashCode`. Both of those should
 take into account the object's fields. If those fields *change* then that
 implies the object's hash code can change.

Most hash-based collections don't anticipate that—they assume an object's hash code will be the same forever and may behave unpredictably if that isn't true.

### DON'T make the parameter to `==` nullable

 #
 Linter rule: avoid_null_checks_in_equality_operators

 The language specifies that `null` is equal only to itself, and that the `==`
 method is called only if the right-hand side is not `null`.

```
class Person {
 final String name;
 // ···
 bool operator ==(Object other) => other is Person && name == other.name;
}
```
```
class Person {
 final String name;
 // ···
 bool operator ==(Object? other) =>
 other != null && other is Person && name == other.name;
}
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Page last updated on 2026-06-04. View source or report an issue.
