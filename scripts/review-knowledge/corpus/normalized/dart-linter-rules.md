# Linter rules

Details about the Dart linter and its style rules you can choose.

 Use the Dart linter to identify possible problems in your Dart code.
 You can use the linter through your IDE
 or with the `dart analyze` command.
 For information on how to enable and disable individual linter rules, see
 individual rules sections of the
 analyzer documentation.

This page lists all the linter rules, with details such as when you might want to use each rule, what code patterns trigger it, and how you might fix your code.

## Sets

#To avoid the need to individually select compatible linter rules, consider starting with a linter rule set, which the following packages provide:

- lints
-
 Contains two rule sets curated by the Dart team. We recommend using at least the `core`rule set, which is used when scoring packages uploaded to pub.dev. Or, better yet, use the`recommended`rule set, a superset of`core`that identifies additional issues and enforces style and format. If you're writing Flutter code, use the rule set in the`flutter_lints`package, which builds on`lints`.

- flutter_lints
-
 Contains the `flutter`rule set, which the Flutter team encourages you to use in Flutter apps, packages, and plugins. This rule set is a superset of the`recommended`set, which is itself a superset of the`core`set that partially determines the score of packages uploaded to pub.dev.

To learn how to use a specific rule set, visit the documentation for enabling and disabling linter rules.

 To find more predefined rule sets,
 check out the `#lints` topic
 on pub.dev.

## Status

#Each rule has a status or maturity level:

- **Stable**
-
 These rules are safe to use and are verified as functional with the latest versions of the Dart language. All rules are considered stable unless they're marked as experimental, deprecated, or removed.
- **Experimental**
-
 These rules are still under evaluation and might never be stabilized. Use these with caution and report any issues you come across.
- **Deprecated**
-
 These rules are no longer suggested for use and might be removed in a future Dart release.
- **Removed**
-
 These rules have already been removed in the latest stable Dart release.

## Quick fixes

#Some rules can be fixed automatically using quick fixes. A quick fix is an automated edit targeted at fixing the issue reported by the linter rule.

 If the rule has a quick fix,
 it can be applied using `dart fix`
 or using your editor with Dart support.
 To learn more, see Quick fixes for analysis issues.

## Rules

#
 The following is an index of all linter rules and
 a short description of their functionality.
 To learn more about a specific rule,
 click the **Learn more** button on its card.

 For an auto-generated list containing all linter rules
 in Dart `3.12.2`,
 check out All linter rules.

Separate the control structure expression from its statement.

Specify `@required` on named parameters without defaults.

Avoid `bool` literals in conditional expressions.

Avoid defining a class that contains only static members.

Avoid overloading operator == and hashCode on classes not marked `@immutable`.

Avoid escaping inner quotes by converting surrounding quotes.

Avoid using `forEach` with a function literal.

Don't declare multiple variables on a single line.

Don't check for `null` in custom `==` operators.

Avoid relative imports for files in `lib/`.

Don't rename parameters of overridden methods.

Avoid returning null from members whose return type is bool, double, int, or num.

Avoid returning this from methods just to enable a fluent interface.

Avoid shadowing type parameters.

Avoid single cascade in expression statements.

Avoid annotating types for function expression parameters.

Avoid overriding a final field to return different values if called multiple times.

Avoid defining unused parameters in constructors.

Avoid using web-only libraries outside Flutter web plugin packages.

Prefer using lowerCamelCase for constant names.

DO use curly braces for all flow control structures.

Attach library doc comments to library directives.

Avoid using deprecated elements from within the package in which they are declared.

DO reference all public properties in debug methods.

 There should be no `Future`-returning calls in synchronous functions unless they are assigned or returned.

Use `;` instead of `{}` for empty constructor bodies.

Define case clauses for all constants in enum-like classes.

Use Flutter TODO format: // TODO(username): message, https://URL-to-issue.

Always override `hashCode` if overriding `==`.

Don't import implementation files from another package.

Explicitly tear-off `call` methods when using an object as a Function.

Initialize the field in the field's initializer.

Avoid runtime type tests with JS interop types where the result may not be platform-consistent.

Conditions should not unconditionally evaluate to `true` or to `false`.

Attach library annotations to library directives.

Use `lowercase_with_underscores` when specifying a library prefix.

Avoid using private types in public APIs.

Don't use more than one case with same value.

Avoid leading underscores for library prefixes.

Avoid leading underscores for local identifiers.

Don't use wildcard parameters or variables.

Name non-constant identifiers using lowerCamelCase.

Don't use `null` check on a potentially nullable type parameter.

Do not pass `null` as an argument where a closure is expected.

Omit obvious type annotations for local variables.

Omit obvious type annotations for top-level and static variables.

Avoid defining a one-member abstract class when a simple function will do.

Prefix library names with the package name and a dot-separated path.

Use adjacent strings to concatenate string literals.

Prefer using `??=` over testing for `null`.

Prefer declaring `const` constructors on `@immutable` classes.

Prefer const literals as parameters of constructors on @immutable classes.

Prefer defining constructors instead of static methods to create instances.

Prefer double quotes where they won't require escape sequences.

Use `=` to separate a named parameter from its default value.

Use => for short members whose body is a single return statement.

Prefer final in for-each loop variable if reference is not reassigned.

Prefer final for variable declarations if they are not reassigned.

Prefer final for parameter declarations if they are not reassigned.

Prefer `for` elements when building maps from iterables.

Use a function declaration to bind a function to a name.

Prefer generic function type aliases.

Prefer if elements to conditional expressions where possible.

Use initializing formals when possible.

Use interpolation to compose strings and values.

Use `isNotEmpty` for `Iterable`s and `Map`s.

Prefer to use `whereType` on iterable.

Prefer typing uninitialized variables and fields.

Don't use the Null type, unless you are positive that you don't want void.

Provide a deprecation message, via `@Deprecated("message")`.

Use trailing commas for all parameter lists and argument lists.

Sort child properties last in widget instance creations.

Specify non-obvious type annotations for local variables.

Specify non-obvious type annotations for top-level and static variables.

Don't use constant patterns with type literals.

 `Future` results in `async` function bodies must be `await`ed or marked
 `unawaited` using `dart:async`.

Use of angle brackets in a doc comment is treated as HTML by Markdown.

Avoid using braces in interpolation when not needed.

Don't use an explicit `const` in a generative enum constructor.

Avoid wrapping fields in getters and setters just to be "safe".

Don't specify the `late` modifier when it is not needed.

Avoid library directives unless they have documentation comments or annotations.

Don't have a library name in a `library` declaration.

Avoid `null` in `null`-aware assignment.

Unnecessary null aware operator on extension on a nullable type.

Avoid using `null` in `??` operators.

Use a non-nullable type for a final variable initialized with a non-nullable value.

Don't override a method to do a super method invocation with the same parameters.

Unnecessary primary constructor bodies can be removed.

Remove unnecessary backslashes in strings.

Unnecessary string interpolation.

Don't access members with `this` unless avoiding shadowing.

Don't use an explicit type name in a constructor.

Do not use `BuildContext` across asynchronous gaps.

Prefer an 8-digit hexadecimal integer (for example, 0xFFFFFFFF) to instantiate a Color.

Use generic function type syntax for parameters.

Use `??` operators to convert `null`s to `bool`s.

Prefer intValue.isOdd/isEven instead of checking the result of % 2.

Use late for private members with a non-nullable type.

If-elements testing for null can be replaced with null-aware elements.

Use rethrow to rethrow a caught exception.

Use a setter for operations that conceptually change a property.

Use string in part of directives.

Use super-initializer parameters where possible.

Start the name of the method with to/_to or as/_as if applicable.

Avoid declaring parameters with `var` and no type annotation.

Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Page last updated on 2026-07-31. View source or report an issue.

# All linter rules

Auto-generated configuration enabling all linter rules.

 The following is an auto-generated list of all linter rules
 available in the Dart SDK as of version `3.12.2`.
 Add them to your
 `analysis_options.yaml` file
 and adjust as you see fit.

```
linter:
 rules:
 - always_declare_return_types
 - always_put_control_body_on_new_line
 - always_put_required_named_parameters_first
 - always_specify_types
 - always_use_package_imports
 - annotate_overrides
 - annotate_redeclares
 - async_return_with_no_await
 - avoid_annotating_with_dynamic
 - avoid_bool_literals_in_conditional_expressions
 - avoid_catches_without_on_clauses
 - avoid_catching_errors
 - avoid_classes_with_only_static_members
 - avoid_double_and_int_checks
 - avoid_dynamic_calls
 - avoid_empty_else
 - avoid_equals_and_hash_code_on_mutable_classes
 - avoid_escaping_inner_quotes
 - avoid_field_initializers_in_const_classes
 - avoid_final_parameters
 - avoid_function_literals_in_foreach_calls
 - avoid_futureor_void
 - avoid_implementing_value_types
 - avoid_init_to_null
 - avoid_js_rounded_ints
 - avoid_multiple_declarations_per_line
 - avoid_null_checks_in_equality_operators
 - avoid_positional_boolean_parameters
 - avoid_print
 - avoid_private_typedef_functions
 - avoid_redundant_argument_values
 - avoid_relative_lib_imports
 - avoid_renaming_method_parameters
 - avoid_return_types_on_setters
 - avoid_returning_null_for_void
 - avoid_returning_this
 - avoid_setters_without_getters
 - avoid_shadowing_type_parameters
 - avoid_single_cascade_in_expression_statements
 - avoid_slow_async_io
 - avoid_type_to_string
 - avoid_types_as_parameter_names
 - avoid_types_on_closure_parameters
 - avoid_unnecessary_containers
 - avoid_unused_constructor_parameters
 - avoid_void_async
 - avoid_web_libraries_in_flutter
 - await_only_futures
 - camel_case_extensions
 - camel_case_types
 - cancel_subscriptions
 - cascade_invocations
 - cast_nullable_to_non_nullable
 - close_sinks
 - collection_methods_unrelated_type
 - combinators_ordering
 - comment_references
 - conditional_uri_does_not_exist
 - constant_identifier_names
 - control_flow_in_finally
 - curly_braces_in_flow_control_structures
 - dangling_library_doc_comments
 - depend_on_referenced_packages
 - deprecated_consistency
 - deprecated_member_use_from_same_package
 - diagnostic_describe_all_properties
 - directives_ordering
 - discarded_futures
 - do_not_use_environment
 - document_ignores
 - empty_catches
 - empty_constructor_bodies
 - empty_container_bodies
 - empty_statements
 - eol_at_end_of_file
 - exhaustive_cases
 - file_names
 - flutter_style_todos
 - hash_and_equals
 - implementation_imports
 - implicit_call_tearoffs
 - implicit_reopen
 - initialize_in_field_declaration
 - invalid_case_patterns
 - invalid_runtime_check_with_js_interop_types
 - join_return_with_assignment
 - leading_newlines_in_multiline_strings
 - library_annotations
 - library_names
 - library_prefixes
 - library_private_types_in_public_api
 - lines_longer_than_80_chars
 - literal_only_boolean_expressions
 - matching_super_parameters
 - missing_code_block_language_in_doc_comment
 - missing_whitespace_between_adjacent_strings
 - no_adjacent_strings_in_list
 - no_default_cases
 - no_duplicate_case_values
 - no_dynamic_casts
 - no_leading_underscores_for_library_prefixes
 - no_leading_underscores_for_local_identifiers
 - no_literal_bool_comparisons
 - no_logic_in_create_state
 - no_raw_types
 - no_runtimetype_tostring
 - no_self_assignments
 - no_wildcard_variable_uses
 - non_constant_identifier_names
 - noop_primitive_operations
 - null_check_on_nullable_type_parameter
 - null_closures
 - omit_local_variable_types
 - omit_obvious_local_variable_types
 - omit_obvious_property_types
 - one_member_abstracts
 - only_throw_errors
 - overridden_fields
 - package_names
 - package_prefixed_library_names
 - parameter_assignments
 - prefer_adjacent_string_concatenation
 - prefer_asserts_in_initializer_lists
 - prefer_asserts_with_message
 - prefer_collection_literals
 - prefer_conditional_assignment
 - prefer_const_constructors
 - prefer_const_constructors_in_immutables
 - prefer_const_declarations
 - prefer_const_literals_to_create_immutables
 - prefer_constructors_over_static_methods
 - prefer_contains
 - prefer_double_quotes
 - prefer_expression_function_bodies
 - prefer_final_fields
 - prefer_final_in_for_each
 - prefer_final_locals
 - prefer_final_parameters
 - prefer_for_elements_to_map_fromIterable
 - prefer_foreach
 - prefer_function_declarations_over_variables
 - prefer_generic_function_type_aliases
 - prefer_if_elements_to_conditional_expressions
 - prefer_if_null_operators
 - prefer_initializing_formals
 - prefer_inlined_adds
 - prefer_int_literals
 - prefer_interpolation_to_compose_strings
 - prefer_is_empty
 - prefer_is_not_empty
 - prefer_is_not_operator
 - prefer_iterable_whereType
 - prefer_mixin
 - prefer_null_aware_method_calls
 - prefer_null_aware_operators
 - prefer_relative_imports
 - prefer_single_quotes
 - prefer_spread_collections
 - prefer_typing_uninitialized_variables
 - prefer_void_to_null
 - provide_deprecation_message
 - public_member_api_docs
 - recursive_getters
 - remove_deprecations_in_breaking_versions
 - require_trailing_commas
 - secure_pubspec_urls
 - simple_directive_paths
 - simplify_variable_pattern
 - sized_box_for_whitespace
 - sized_box_shrink_expand
 - slash_for_doc_comments
 - sort_child_properties_last
 - sort_constructors_first
 - sort_pub_dependencies
 - sort_unnamed_constructors_first
 - specify_nonobvious_local_variable_types
 - specify_nonobvious_property_types
 - strict_top_level_inference
 - switch_on_type
 - test_types_in_equals
 - throw_in_finally
 - tighten_type_of_initializing_formals
 - type_annotate_public_apis
 - type_init_formals
 - type_literal_in_constant_pattern
 - unawaited_futures
 - unintended_html_in_doc_comment
 - unnecessary_async
 - unnecessary_await_in_return
 - unnecessary_brace_in_string_interps
 - unnecessary_breaks
 - unnecessary_const
 - unnecessary_const_in_enum_constructor
 - unnecessary_constructor_name
 - unnecessary_final
 - unnecessary_getters_setters
 - unnecessary_ignore
 - unnecessary_lambdas
 - unnecessary_late
 - unnecessary_library_directive
 - unnecessary_library_name
 - unnecessary_new
 - unnecessary_null_aware_assignments
 - unnecessary_null_aware_operator_on_extension_on_nullable
 - unnecessary_null_checks
 - unnecessary_null_in_if_null_operators
 - unnecessary_nullable_for_final_variable_declarations
 - unnecessary_overrides
 - unnecessary_parenthesis
 - unnecessary_primary_constructor_body
 - unnecessary_raw_strings
 - unnecessary_statements
 - unnecessary_string_escapes
 - unnecessary_string_interpolations
 - unnecessary_this
 - unnecessary_to_list_in_spreads
 - unnecessary_type_name_in_constructor
 - unnecessary_unawaited
 - unnecessary_underscores
 - unreachable_from_main
 - unrelated_type_equality_checks
 - unsafe_variance
 - use_build_context_synchronously
 - use_colored_box
 - use_declaring_parameters
 - use_decorated_box
 - use_enums
 - use_full_hex_values_for_flutter_colors
 - use_function_type_syntax_for_parameters
 - use_if_null_to_convert_nulls_to_bools
 - use_is_even_rather_than_modulo
 - use_key_in_widget_constructors
 - use_late_for_private_fields_and_variables
 - use_named_constants
 - use_null_aware_elements
 - use_raw_strings
 - use_rethrow_when_possible
 - use_setters_to_change_properties
 - use_string_buffers
 - use_string_in_part_of_directives
 - use_super_parameters
 - use_test_throws_matchers
 - use_to_and_as_if_applicable
 - use_truncating_division
 - valid_regexps
 - var_with_no_type_annotation
 - void_checks
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# always_

 Learn about the always_declare_return_types linter rule.

Declare method return types.

## Details

#**DO** declare method return types.

 When declaring a method or function *always* specify a return type.
 Declaring return types for functions helps improve your codebase by allowing the
 analyzer to more adequately check your code for errors that could occur during
 runtime.

**BAD:**

```
main() { }
_bar() => _Foo();
class _Foo {
 _foo() => 42;
}
```
**GOOD:**

```
void main() { }
_Foo _bar() => _Foo();
class _Foo {
 int _foo() => 42;
}
typedef predicate = bool Function(Object o);
```

## Enable

#
 To enable the `always_declare_return_types` rule, add `always_declare_return_types`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - always_declare_return_types
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `always_declare_return_types: true` under **linter > rules**:

```
linter:
 rules:
 always_declare_return_types: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# always_

 Learn about the always_put_control_body_on_new_line linter rule.

Separate the control structure expression from its statement.

## Details

#From the style guide for the flutter repo:

**DO** separate the control structure expression from its statement.

 Don't put the statement part of an `if`, `for`, `while`, `do`
 on the same line
 as the expression, even if it is short. Doing so makes it unclear that there
 is relevant code there. This is especially important for early returns.

**BAD:**

```
if (notReady) return;
if (notReady)
 return;
else print('ok')
while (condition) i += 1;
```
**GOOD:**

```
if (notReady)
 return;
if (notReady)
 return;
else
 print('ok')
while (condition)
 i += 1;
```
Note that this rule can conflict with the Dart formatter, and should not be enabled when the Dart formatter is used.

## Enable

#
 To enable the `always_put_control_body_on_new_line` rule, add `always_put_control_body_on_new_line`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - always_put_control_body_on_new_line
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `always_put_control_body_on_new_line: true` under **linter > rules**:

```
linter:
 rules:
 always_put_control_body_on_new_line: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# always_

 Learn about the always_put_required_named_parameters_first linter rule.

Put required named parameters first.

## Details

#**DO** specify `required` on named parameter before other named parameters.

**BAD:**

```
m({b, c, required a}) ;
```
**GOOD:**

```
m({required a, b, c}) ;
```
**BAD:**

```
m({b, c, @required a}) ;
```
**GOOD:**

```
m({@required a, b, c}) ;
```

## Enable

#
 To enable the `always_put_required_named_parameters_first` rule, add `always_put_required_named_parameters_first`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - always_put_required_named_parameters_first
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `always_put_required_named_parameters_first: true` under **linter > rules**:

```
linter:
 rules:
 always_put_required_named_parameters_first: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# always_

 Learn about the always_require_non_null_named_parameters linter rule.

Specify `@required` on named parameters without defaults.

## Details

#NOTE: This rule is removed in Dart 3.3.0; it is no longer functional.

 **DO** specify `@required` on named parameters without a default value on which
 an `assert(param != null)` is done.

**BAD:**

```
m1({a}) {
 assert(a != null);
}
```
**GOOD:**

```
m1({@required a}) {
 assert(a != null);
}
m2({a: 1}) {
 assert(a != null);
}
```
NOTE: Only asserts at the start of the bodies will be taken into account.

## Enable

#
 To enable the `always_require_non_null_named_parameters` rule, add `always_require_non_null_named_parameters`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - always_require_non_null_named_parameters
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `always_require_non_null_named_parameters: true` under **linter > rules**:

```
linter:
 rules:
 always_require_non_null_named_parameters: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# always_

 Learn about the always_specify_types linter rule.

Specify type annotations.

## Details

#From the style guide for the flutter repo:

**DO** specify type annotations.

 Avoid `var` when specifying that a type is unknown and short-hands that elide
 type annotations. Use `dynamic` if you are being explicit that the type is
 unknown. Use `Object` if you are being explicit that you want an object that
 implements `==` and `hashCode`.

**BAD:**

```
var foo = 10;
final bar = Bar();
const quux = 20;
```
**GOOD:**

```
int foo = 10;
final Bar bar = Bar();
String baz = 'hello';
const int quux = 20;
```
 NOTE: Using the `@optionalTypeArgs` annotation in the `meta` package, API
 authors can special-case type parameters whose type needs to be dynamic but whose
 declaration should be treated as optional. For example, suppose you have a
 `Key` object whose type parameter you'd like to treat as optional. Using the
 `@optionalTypeArgs` would look like this:

```
import 'package:meta/meta.dart';
@optionalTypeArgs
class Key<T> {
 ...
}
void main() {
 Key s = Key(); // OK!
}
```
## Incompatible rules

#The `always_specify_types` lint is incompatible with the following rules:

-
 `avoid_types_on_closure_parameters`
- `omit_local_variable_types`
-
 `omit_obvious_local_variable_types`
-
 `omit_obvious_property_types`

## Enable

#
 To enable the `always_specify_types` rule, add `always_specify_types` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - always_specify_types
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `always_specify_types: true` under **linter > rules**:

```
linter:
 rules:
 always_specify_types: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# always_

 Learn about the always_use_package_imports linter rule.

Avoid relative imports for files in `lib/`.

## Details

#**DO** avoid relative imports for files in `lib/`.

 When mixing relative and absolute imports it's possible to create confusion
 where the same member gets imported in two different ways. One way to avoid
 that is to ensure you consistently use absolute imports for files within the
 `lib/` directory.

This is the opposite of 'prefer_relative_imports'.

 You can also use 'avoid_relative_lib_imports' to disallow relative imports of
 files within `lib/` directory outside of it (for example `test/`).

**BAD:**

```
import 'baz.dart';
import 'src/bag.dart'
import '../lib/baz.dart';
...
```
**GOOD:**

```
import 'package:foo/bar.dart';
import 'package:foo/baz.dart';
import 'package:foo/src/baz.dart';
...
```
## Incompatible rules

#The `always_use_package_imports` lint is incompatible with the following rules:

## Enable

#
 To enable the `always_use_package_imports` rule, add `always_use_package_imports`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - always_use_package_imports
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `always_use_package_imports: true` under **linter > rules**:

```
linter:
 rules:
 always_use_package_imports: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# annotate_

 Learn about the annotate_overrides linter rule.

Annotate overridden members.

## Details

#**DO** annotate overridden methods and fields.

This practice improves code readability and helps protect against unintentionally overriding superclass members.

**BAD:**

```
class Cat {
 int get lives => 9;
}
class Lucky extends Cat {
 final int lives = 14;
}
```
**GOOD:**

```
abstract class Dog {
 String get breed;
 void bark() {}
}
class Husky extends Dog {
 @override
 final String breed = 'Husky';
 @override
 void bark() {}
}
```

## Enable

#
 To enable the `annotate_overrides` rule, add `annotate_overrides` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - annotate_overrides
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `annotate_overrides: true` under **linter > rules**:

```
linter:
 rules:
 annotate_overrides: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# annotate_

 Learn about the annotate_redeclares linter rule.

Annotate redeclared members.

## Details

#**DO** annotate redeclared members.

This practice improves code readability and helps protect against unintentionally redeclaring members or being surprised when a member ceases to redeclare (due for example to a rename refactoring).

**BAD:**

```
class C {
 void f() { }
}
extension type E(C c) implements C {
 void f() {
 ...
 }
}
```
**GOOD:**

```
import 'package:meta/meta.dart';
class C {
 void f() { }
}
extension type E(C c) implements C {
 @redeclare
 void f() {
 ...
 }
}
```

## Enable

#
 To enable the `annotate_redeclares` rule, add `annotate_redeclares` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - annotate_redeclares
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `annotate_redeclares: true` under **linter > rules**:

```
linter:
 rules:
 annotate_redeclares: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# async_

 Learn about the async_return_with_no_await linter rule.

Return with no await.

## Details

#
 **DO** use `await` when returning a `Future` from an `async`
 function.

**BAD:**

```
Future<String> futureString(Future<String> value) async {
 return value;
}
Future<int> futureInt(Future<int> value) async => value;
```
**GOOD:**

```
Future<String> futureString(Future<String> value) async {
 return await value;
}
Future<int> futureInt(Future<int> value) => value;
```

## Enable

#
 To enable the `async_return_with_no_await` rule, add `async_return_with_no_await`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - async_return_with_no_await
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `async_return_with_no_await: true` under **linter > rules**:

```
linter:
 rules:
 async_return_with_no_await: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_annotating_with_dynamic linter rule.

Avoid annotating with `dynamic` when not required.

## Details

#**AVOID** annotating with `dynamic` when not required.

 As `dynamic` is the assumed return value of a function or method, it is usually
 not necessary to annotate it.

**BAD:**

```
dynamic lookUpOrDefault(String name, Map map, dynamic defaultValue) {
 var value = map[name];
 if (value != null) return value;
 return defaultValue;
}
```
**GOOD:**

```
lookUpOrDefault(String name, Map map, defaultValue) {
 var value = map[name];
 if (value != null) return value;
 return defaultValue;
}
```

## Enable

#
 To enable the `avoid_annotating_with_dynamic` rule, add `avoid_annotating_with_dynamic`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_annotating_with_dynamic
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_annotating_with_dynamic: true` under **linter > rules**:

```
linter:
 rules:
 avoid_annotating_with_dynamic: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_as linter rule.

Avoid using `as`.

## Details

#NOTE: This rule was removed from the SDK in Dart 3; it is no longer functional. Its advice is compiler-specific and mostly obsolete with null safety.

**AVOID** using `as`.

 If you know the type is correct, use an assertion or assign to a more
 narrowly-typed variable (this avoids the type check in release mode; `as`
 is not
 compiled out in release mode). If you don't know whether the type is
 correct, check using `is` (this avoids the exception that `as`
 raises).

**BAD:**

```
(pm as Person).firstName = 'Seth';
```
**GOOD:**

```
if (pm is Person)
 pm.firstName = 'Seth';
```
but certainly not

**BAD:**

```
try {
 (pm as Person).firstName = 'Seth';
} on CastError { }
```
 Note that an exception is made in the case of `dynamic` since the cast has no
 performance impact.

**OK:**

```
HasScrollDirection scrollable = renderObject as dynamic;
```

## Enable

#
 To enable the `avoid_as` rule, add `avoid_as` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_as
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_as: true` under **linter > rules**:

```
linter:
 rules:
 avoid_as: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_bool_literals_in_conditional_expressions linter rule.

Avoid `bool` literals in conditional expressions.

## Details

#**AVOID** `bool` literals in conditional expressions.

**BAD:**

```
condition ? true : boolExpression
condition ? false : boolExpression
condition ? boolExpression : true
condition ? boolExpression : false
```
**GOOD:**

```
condition || boolExpression
!condition && boolExpression
!condition || boolExpression
condition && boolExpression
```

## Enable

#
 To enable the `avoid_bool_literals_in_conditional_expressions` rule, add `avoid_bool_literals_in_conditional_expressions`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_bool_literals_in_conditional_expressions
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_bool_literals_in_conditional_expressions: true` under **linter > rules**:

```
linter:
 rules:
 avoid_bool_literals_in_conditional_expressions: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_catches_without_on_clauses linter rule.

Avoid catches without on clauses.

## Details

#From Effective Dart:

**AVOID** catches without on clauses.

Using catch clauses without on clauses make your code prone to encountering unexpected errors that won't be thrown (and thus will go unnoticed).

**BAD:**

```
try {
 somethingRisky()
} catch(e) {
 doSomething(e);
}
```
**GOOD:**

```
try {
 somethingRisky()
} on Exception catch(e) {
 doSomething(e);
}
```
A few exceptional cases are allowed:

- If the body of the catch rethrows the exception.
-
 If the caught exception is "directly used" in an argument to `Future.error`,`Completer.completeError`, or`FlutterError.reportError`, or any function with a return type of`Never`.
- If the caught exception is "directly used" in a new throw-expression.

 In these cases, "directly used" means that the exception is referenced within
 the relevant code (like within an argument). If the exception variable is
 referenced *before* the relevant code, for example to instantiate a wrapper
 exception, the variable is not "directly used."

## Enable

#
 To enable the `avoid_catches_without_on_clauses` rule, add `avoid_catches_without_on_clauses`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_catches_without_on_clauses
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_catches_without_on_clauses: true` under **linter > rules**:

```
linter:
 rules:
 avoid_catches_without_on_clauses: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_catching_errors linter rule.

Don't explicitly catch `Error` or types that implement it.

## Details

#**DON'T** explicitly catch `Error` or types that implement it.

Errors differ from Exceptions in that Errors can be analyzed and prevented prior to runtime. It should almost never be necessary to catch an error at runtime.

**BAD:**

```
try {
 somethingRisky();
} on Error catch(e) {
 doSomething(e);
}
```
**GOOD:**

```
try {
 somethingRisky();
} on Exception catch(e) {
 doSomething(e);
}
```

## Enable

#
 To enable the `avoid_catching_errors` rule, add `avoid_catching_errors` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_catching_errors
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_catching_errors: true` under **linter > rules**:

```
linter:
 rules:
 avoid_catching_errors: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_classes_with_only_static_members linter rule.

Avoid defining a class that contains only static members.

## Details

#From Effective Dart:

**AVOID** defining a class that contains only static members.

Creating classes with the sole purpose of providing utility or otherwise static methods is discouraged. Dart allows functions to exist outside of classes for this very reason.

**BAD:**

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
**GOOD:**

```
DateTime mostRecent(List<DateTime> dates) {
 return dates.reduce((a, b) => a.isAfter(b) ? a : b);
}
const _favoriteMammal = 'weasel';
```

## Enable

#
 To enable the `avoid_classes_with_only_static_members` rule, add `avoid_classes_with_only_static_members`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_classes_with_only_static_members
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_classes_with_only_static_members: true` under **linter > rules**:

```
linter:
 rules:
 avoid_classes_with_only_static_members: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_double_and_int_checks linter rule.

Avoid `double` and `int` checks.

## Details

#**AVOID** to check if type is `double` or `int`.

 When compiled to JS, integer values are represented as floats. That can lead to
 some unexpected behavior when using either `is` or `is!` where the type is
 either `int` or `double`.

**BAD:**

```
f(num x) {
 if (x is double) {
 ...
 } else if (x is int) {
 ...
 }
}
```
**GOOD:**

```
f(dynamic x) {
 if (x is num) {
 ...
 } else {
 ...
 }
}
```

## Enable

#
 To enable the `avoid_double_and_int_checks` rule, add `avoid_double_and_int_checks`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_double_and_int_checks
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_double_and_int_checks: true` under **linter > rules**:

```
linter:
 rules:
 avoid_double_and_int_checks: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_dynamic_calls linter rule.

Avoid method calls or property accesses on a `dynamic` target.

## Details

#
 **DO** avoid method calls or accessing properties on an object that is either
 explicitly or implicitly statically typed `dynamic`. Dynamic calls are treated
 slightly different in every runtime environment and compiler, but most
 production modes (and even some development modes) have both compile size and
 runtime performance penalties associated with dynamic calls.

 Additionally, targets typed `dynamic` disables most static analysis, meaning it
 is easier to lead to a runtime `NoSuchMethodError` or `TypeError`
 than properly
 statically typed Dart code.

There is an exception to methods and properties that exist on `Object?`:

- `a.hashCode`
- `a.runtimeType`
- `a.noSuchMethod(someInvocation)`
- `a.toString()`

 ... these members are dynamically dispatched in the web-based runtimes, but not
 in the VM-based ones. Additionally, they are so common that it would be very
 punishing to disallow `any.toString()` or `any == true`, for example.

 Note that despite `Function` being a type, the semantics are close to identical
 to `dynamic`, and calls to an object that is typed `Function`
 will also trigger
 this lint.

Dynamic calls are allowed on cast expressions (`as dynamic` or `as Function`).

**BAD:**

```
void explicitDynamicType(dynamic object) {
 print(object.foo());
}
void implicitDynamicType(object) {
 print(object.foo());
}
abstract class SomeWrapper {
 T doSomething<T>();
}
void inferredDynamicType(SomeWrapper wrapper) {
 var object = wrapper.doSomething();
 print(object.foo());
}
void callDynamic(dynamic function) {
 function();
}
void functionType(Function function) {
 function();
}
```
**GOOD:**

```
void explicitType(Fooable object) {
 object.foo();
}
void castedType(dynamic object) {
 (object as Fooable).foo();
}
abstract class SomeWrapper {
 T doSomething<T>();
}
void inferredType(SomeWrapper wrapper) {
 var object = wrapper.doSomething<Fooable>();
 object.foo();
}
void functionTypeWithParameters(Function() function) {
 function();
}
```

## Enable

#
 To enable the `avoid_dynamic_calls` rule, add `avoid_dynamic_calls` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_dynamic_calls
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_dynamic_calls: true` under **linter > rules**:

```
linter:
 rules:
 avoid_dynamic_calls: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_empty_else linter rule.

Avoid empty statements in else clauses.

## Details

#
 **AVOID** empty statements in the `else` clause of `if` statements.

**BAD:**

```
if (x > y)
 print('1');
else ;
 print('2');
```
 If you want a statement that follows the empty clause to *conditionally* run,
 remove the dangling semicolon to include it in the `else` clause.
 Optionally, also enclose the else's statement in a block.

**GOOD:**

```
if (x > y)
 print('1');
else
 print('2');
```
**GOOD:**

```
if (x > y) {
 print('1');
} else {
 print('2');
}
```
 If you want a statement that follows the empty clause to *unconditionally* run,
 remove the `else` clause.

**GOOD:**

```
if (x > y) print('1');
print('2');
```

## Enable

#
 To enable the `avoid_empty_else` rule, add `avoid_empty_else` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_empty_else
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_empty_else: true` under **linter > rules**:

```
linter:
 rules:
 avoid_empty_else: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_equals_and_hash_code_on_mutable_classes linter rule.

Avoid overloading operator == and hashCode on classes not marked `@immutable`.

## Details

#From Effective Dart:

 **AVOID** overloading operator == and hashCode on classes not marked `@immutable`.

 If a class is not immutable, overloading `operator ==` and `hashCode` can
 lead to unpredictable and undesirable behavior when used in collections.

**BAD:**

```
class B {
 String key;
 const B(this.key);
 @override
 operator ==(other) => other is B && other.key == key;
 @override
 int get hashCode => key.hashCode;
}
```
**GOOD:**

```
@immutable
class A {
 final String key;
 const A(this.key);
 @override
 operator ==(other) => other is A && other.key == key;
 @override
 int get hashCode => key.hashCode;
}
```
 NOTE: The lint checks the use of the `@immutable` annotation, and will trigger
 even if the class is otherwise not mutable. Thus:

**BAD:**

```
class C {
 final String key;
 const C(this.key);
 @override
 operator ==(other) => other is C && other.key == key;
 @override
 int get hashCode => key.hashCode;
}
```

## Enable

#
 To enable the `avoid_equals_and_hash_code_on_mutable_classes` rule, add `avoid_equals_and_hash_code_on_mutable_classes`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_equals_and_hash_code_on_mutable_classes
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_equals_and_hash_code_on_mutable_classes: true` under **linter > rules**:

```
linter:
 rules:
 avoid_equals_and_hash_code_on_mutable_classes: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_escaping_inner_quotes linter rule.

Avoid escaping inner quotes by converting surrounding quotes.

## Details

#Avoid escaping inner quotes by converting surrounding quotes.

**BAD:**

```
var s = 'It\'s not fun';
```
**GOOD:**

```
var s = "It's not fun";
```

## Enable

#
 To enable the `avoid_escaping_inner_quotes` rule, add `avoid_escaping_inner_quotes`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_escaping_inner_quotes
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_escaping_inner_quotes: true` under **linter > rules**:

```
linter:
 rules:
 avoid_escaping_inner_quotes: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_field_initializers_in_const_classes linter rule.

Avoid field initializers in const classes.

## Details

#**AVOID** field initializers in const classes.

 Instead of `final x = const expr;`, you should write `get x => const expr;` and
 not allocate a useless field. As of April 2018 this is true for the VM, but not
 for code that will be compiled to JS.

**BAD:**

```
class A {
 final a = const [];
 const A();
}
```
**GOOD:**

```
class A {
 get a => const [];
 const A();
}
```

## Enable

#
 To enable the `avoid_field_initializers_in_const_classes` rule, add `avoid_field_initializers_in_const_classes`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_field_initializers_in_const_classes
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_field_initializers_in_const_classes: true` under **linter > rules**:

```
linter:
 rules:
 avoid_field_initializers_in_const_classes: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_final_parameters linter rule.

Avoid `final` for parameter declarations.

## Details

#**AVOID** declaring parameters as `final`.

 Declaring parameters as `final` can lead to unnecessarily verbose code,
 especially when using the "parameter_assignments" rule.

**BAD:**

```
void goodParameter(final String label) { // LINT
 print(label);
}
```
**GOOD:**

```
void badParameter(String label) { // OK
 print(label);
}
```
**BAD:**

```
void goodExpression(final int value) => print(value); // LINT
```
**GOOD:**

```
void badExpression(int value) => print(value); // OK
```
**BAD:**

```
[1, 4, 6, 8].forEach((final value) => print(value + 2)); // LINT
```
**GOOD:**

```
[1, 4, 6, 8].forEach((value) => print(value + 2)); // OK
```
## Incompatible rules

#The `avoid_final_parameters` lint is incompatible with the following rules:

## Enable

#
 To enable the `avoid_final_parameters` rule, add `avoid_final_parameters` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_final_parameters
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_final_parameters: true` under **linter > rules**:

```
linter:
 rules:
 avoid_final_parameters: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_function_literals_in_foreach_calls linter rule.

Avoid using `forEach` with a function literal.

## Details

#**AVOID** using `forEach` with a function literal.

 The `for` loop enables a developer to be clear and explicit as to their intent.
 A return in the body of the `for` loop returns from the body of the function,
 where as a return in the body of the `forEach` closure only returns a value
 for that iteration of the `forEach`. The body of a `for` loop can contain
 `await`s, while the closure body of a `forEach` cannot.

**BAD:**

```
people.forEach((person) {
 ...
});
```
**GOOD:**

```
for (var person in people) {
 ...
}
```

## Enable

#
 To enable the `avoid_function_literals_in_foreach_calls` rule, add `avoid_function_literals_in_foreach_calls`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_function_literals_in_foreach_calls
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_function_literals_in_foreach_calls: true` under **linter > rules**:

```
linter:
 rules:
 avoid_function_literals_in_foreach_calls: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_futureor_void linter rule.

Avoid using 'FutureOr

## Details

#
 **AVOID** using `FutureOr<void>` as the type of a result. This type is
 problematic because it may appear to encode that a result is either a
 `Future<void>`, or the result should be discarded (when it is `void`).
 However, there is no safe way to detect whether we have one or the other
 case (because an expression of type `void` can evaluate to any object
 whatsoever, including a future of any type).

It is also conceptually unsound to have a type whose meaning is something like "ignore this object; also, take a look because it might be a future".

 An exception is made for contravariant occurrences of the type
 `FutureOr<void>` (e.g., for the type of a formal parameter), and no
 warning is emitted for these occurrences. The reason for this exception
 is that the type does not describe a result, it describes a constraint
 on a value provided by others. Similarly, an exception is made for type
 alias declarations, because they may well be used in a contravariant
 position (e.g., as the type of a formal parameter). Hence, in type alias
 declarations, only the type parameter bounds are checked.

 A replacement for the type `FutureOr<void>` which is often useful is
 `Future<void>?`. This type encodes that the result is either a
 `Future<void>` or it is null, and there is no ambiguity at run time
 since no object can have both types.

 It may not always be possible to use the type `Future<void>?` as a
 replacement for the type `FutureOr<void>`, because the latter is a
 supertype of all types, and the former is not. In this case it may be a
 useful remedy to replace `FutureOr<void>` by the type `void`.

**BAD:**

```
FutureOr<void> m() {...}
```
**GOOD:**

```
Future<void>? m() {...}
```
 **This rule is experimental.**
 It is being evaluated, and it might be changed or removed.
 Feedback on its behavior is welcome! The primary relevant issue is
 dart-lang/sdk#59292.

## Enable

#
 To enable the `avoid_futureor_void` rule, add `avoid_futureor_void` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_futureor_void
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_futureor_void: true` under **linter > rules**:

```
linter:
 rules:
 avoid_futureor_void: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_implementing_value_types linter rule.

Don't implement classes that override `==`.

## Details

#**DON'T** implement classes that override `==`.

 The `==` operator is contractually required to be an equivalence relation;
 that is, symmetrically for all objects `o1` and `o2`, `o1 == o2`
 and `o2 == o1`
 must either both be true, or both be false.

NOTE: Dart does not have truevalue types, so instead we consider a class that implements`==`as aproxyfor identifying value types.

 When using `implements`, you do not inherit the method body of `==`, making it
 nearly impossible to follow the contract of `==`. Classes that override
 `==`
 typically are usable directly in tests *without* creating mocks or fakes as
 well. For example, for a given class `Size`:

```
class Size {
 final int inBytes;
 const Size(this.inBytes);
 @override
 bool operator ==(Object other) => other is Size && other.inBytes == inBytes;
 @override
 int get hashCode => inBytes.hashCode;
}
```
**BAD:**

```
class CustomSize implements Size {
 final int inBytes;
 const CustomSize(this.inBytes);
 int get inKilobytes => inBytes ~/ 1000;
}
```
**BAD:**

```
import 'package:test/test.dart';
import 'size.dart';
class FakeSize implements Size {
 int inBytes = 0;
}
void main() {
 test('should not throw on a size >1Kb', () {
 expect(() => someFunction(FakeSize()..inBytes = 1001), returnsNormally);
 });
}
```
**GOOD:**

```
class ExtendedSize extends Size {
 ExtendedSize(int inBytes) : super(inBytes);
 int get inKilobytes => inBytes ~/ 1000;
}
```
**GOOD:**

```
import 'package:test/test.dart';
import 'size.dart';
void main() {
 test('should not throw on a size >1Kb', () {
 expect(() => someFunction(Size(1001)), returnsNormally);
 });
}
```

## Enable

#
 To enable the `avoid_implementing_value_types` rule, add `avoid_implementing_value_types`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_implementing_value_types
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_implementing_value_types: true` under **linter > rules**:

```
linter:
 rules:
 avoid_implementing_value_types: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_init_to_null linter rule.

Don't explicitly initialize variables to `null`.

## Details

#From Effective Dart:

**DON'T** explicitly initialize variables to `null`.

 If a variable has a non-nullable type or is `final`,
 Dart reports a compile error if you try to use it
 before it has been definitely initialized.
 If the variable is nullable and not `const` or `final`,
 then it is implicitly initialized to `null` for you.
 There's no concept of "uninitialized memory" in Dart
 and no need to explicitly initialize a variable to `null` to be "safe".
 Adding `= null` is redundant and unneeded.

**BAD:**

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
**GOOD:**

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

## Enable

#
 To enable the `avoid_init_to_null` rule, add `avoid_init_to_null` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_init_to_null
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_init_to_null: true` under **linter > rules**:

```
linter:
 rules:
 avoid_init_to_null: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_js_rounded_ints linter rule.

Avoid JavaScript rounded ints.

## Details

#
 **AVOID** integer literals that cannot be represented exactly when compiled to
 JavaScript.

 When a program is compiled to JavaScript `int` and `double` become JavaScript
 Numbers. Too large integers (`value < Number.MIN_SAFE_INTEGER` or
 `value > Number.MAX_SAFE_INTEGER`) may be rounded to the closest Number value.

 For instance `1000000000000000001` cannot be represented exactly as a JavaScript
 Number, so `1000000000000000000` will be used instead.

**BAD:**

```
int value = 9007199254740995;
```
**GOOD:**

```
BigInt value = BigInt.parse('9007199254740995');
```

## Enable

#
 To enable the `avoid_js_rounded_ints` rule, add `avoid_js_rounded_ints` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_js_rounded_ints
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_js_rounded_ints: true` under **linter > rules**:

```
linter:
 rules:
 avoid_js_rounded_ints: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_multiple_declarations_per_line linter rule.

Don't declare multiple variables on a single line.

## Details

#**DON'T** declare multiple variables on a single line.

**BAD:**

```
String? foo, bar, baz;
```
**GOOD:**

```
String? foo;
String? bar;
String? baz;
```

## Enable

#
 To enable the `avoid_multiple_declarations_per_line` rule, add `avoid_multiple_declarations_per_line`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_multiple_declarations_per_line
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_multiple_declarations_per_line: true` under **linter > rules**:

```
linter:
 rules:
 avoid_multiple_declarations_per_line: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_null_checks_in_equality_operators linter rule.

Don't check for `null` in custom `==` operators.

## Details

#
 **NOTE:** This lint has been replaced by the
 `non_nullable_equals_parameter` warning and is deprecated.
 Remove all inclusions of this lint from your analysis options.

**DON'T** check for `null` in custom `==` operators.

 As `null` is a special value, no instance of any class (other than `Null`) can
 be equivalent to it. Thus, it is redundant to check whether the other instance
 is `null`.

**BAD:**

```
class Person {
 final String? name;
 @override
 operator ==(Object? other) =>
 other != null && other is Person && name == other.name;
}
```
**GOOD:**

```
class Person {
 final String? name;
 @override
 operator ==(Object? other) => other is Person && name == other.name;
}
```

## Enable

#
 To enable the `avoid_null_checks_in_equality_operators` rule, add `avoid_null_checks_in_equality_operators`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_null_checks_in_equality_operators
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_null_checks_in_equality_operators: true` under **linter > rules**:

```
linter:
 rules:
 avoid_null_checks_in_equality_operators: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_positional_boolean_parameters linter rule.

Avoid positional boolean parameters.

## Details

#From Effective Dart:

**AVOID** positional boolean parameters.

Positional boolean parameters are a bad practice because they are very ambiguous. Using named boolean parameters is much more readable because it inherently describes what the boolean value represents.

**BAD:**

```
Task(true);
Task(false);
ListBox(false, true, true);
Button(false);
```
**GOOD:**

```
Task.oneShot();
Task.repeating();
ListBox(scroll: true, showScrollbars: true);
Button(ButtonState.enabled);
```

## Enable

#
 To enable the `avoid_positional_boolean_parameters` rule, add `avoid_positional_boolean_parameters`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_positional_boolean_parameters
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_positional_boolean_parameters: true` under **linter > rules**:

```
linter:
 rules:
 avoid_positional_boolean_parameters: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_print linter rule.

Avoid `print` calls in production code.

## Details

#**DO** avoid `print` calls in production code.

 For production code, consider using a logging framework.
 If you are using Flutter, you can use `debugPrint`
 or surround `print` calls with a check for `kDebugMode`

**BAD:**

```
void f(int x) {
 print('debug: $x');
 ...
}
```
**GOOD:**

```
void f(int x) {
 debugPrint('debug: $x');
 ...
}
```
**GOOD:**

```
void f(int x) {
 log('log: $x');
 ...
}
```
**GOOD:**

```
void f(int x) {
 if (kDebugMode) {
 print('debug: $x');
 }
 ...
}
```

## Enable

#
 To enable the `avoid_print` rule, add `avoid_print` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_print
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_print: true` under **linter > rules**:

```
linter:
 rules:
 avoid_print: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_private_typedef_functions linter rule.

Avoid private typedef functions.

## Details

#
 **NOTE:** As a private typedef can make some code more readable,
 this lint rule has been deprecated as of Dart 3.13 and is
 set to be removed in a future release of the Dart SDK.
 Remove all inclusions of this lint from your analysis options.
 If you wish to keep enforcing it, consider implementing the
 rule with an analyzer plugin.

 **AVOID** private typedef functions used only once. Prefer inline function
 syntax.

**BAD:**

```
typedef void _F();
m(_F f);
```
**GOOD:**

```
m(void Function() f);
```

## Enable

#
 To enable the `avoid_private_typedef_functions` rule, add `avoid_private_typedef_functions`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_private_typedef_functions
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_private_typedef_functions: true` under **linter > rules**:

```
linter:
 rules:
 avoid_private_typedef_functions: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_redundant_argument_values linter rule.

Avoid redundant argument values.

## Details

#
 **DON'T** pass an argument that matches the corresponding parameter's default
 value.

 Note that a method override can change the default value of a parameter, so that
 an argument may be equal to one default value, and not the other. Take, for
 example, two classes, `A` and `B` where `B` is a subclass of
 `A`, and `B`
 overrides a method declared on `A`, and that method has a parameter with one
 default value in `A`'s declaration, and a different default value in `B`'s
 declaration. If the static type of the target of the invoked method is `B`, and
 `B`'s default value matches the argument, then the argument can be omitted (and
 if the argument value is different, then a lint is not reported). If, however,
 the static type of the target of the invoked method is `A`, then a lint may be
 reported, but we cannot know statically which method is invoked, so the reported
 lint may be a false positive. Such cases can be ignored inline with a comment
 like `// ignore: avoid_redundant_argument_values`.

**BAD:**

```
void f({bool valWithDefault = true, bool? val}) {
 ...
}
void main() {
 f(valWithDefault: true);
}
```
**GOOD:**

```
void f({bool valWithDefault = true, bool? val}) {
 ...
}
void main() {
 f(valWithDefault: false);
 f();
}
```

## Enable

#
 To enable the `avoid_redundant_argument_values` rule, add `avoid_redundant_argument_values`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_redundant_argument_values
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_redundant_argument_values: true` under **linter > rules**:

```
linter:
 rules:
 avoid_redundant_argument_values: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_relative_lib_imports linter rule.

Avoid relative imports for files in `lib/`.

## Details

#**DO** avoid relative imports for files in `lib/`.

 When mixing relative and absolute imports it's possible to create confusion
 where the same member gets imported in two different ways. An easy way to avoid
 that is to ensure you have no relative imports that include `lib/` in their
 paths.

 You can also use 'always_use_package_imports' to disallow relative imports
 between files within `lib/`.

**BAD:**

```
import 'package:foo/bar.dart';
import '../lib/baz.dart';
...
```
**GOOD:**

```
import 'package:foo/bar.dart';
import 'baz.dart';
...
```

## Enable

#
 To enable the `avoid_relative_lib_imports` rule, add `avoid_relative_lib_imports`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_relative_lib_imports
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_relative_lib_imports: true` under **linter > rules**:

```
linter:
 rules:
 avoid_relative_lib_imports: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_renaming_method_parameters linter rule.

Don't rename parameters of overridden methods.

## Details

#**DON'T** rename parameters of overridden methods.

 Methods that override another method, but do not have their own documentation
 comment, will inherit the overridden method's comment when `dart doc` produces
 documentation. If the inherited method contains the name of the parameter (in
 square brackets), then `dart doc` cannot link it correctly.

**BAD:**

```
abstract class A {
 m(a);
}
abstract class B extends A {
 m(b);
}
```
**GOOD:**

```
abstract class A {
 m(a);
}
abstract class B extends A {
 m(a);
}
```

## Enable

#
 To enable the `avoid_renaming_method_parameters` rule, add `avoid_renaming_method_parameters`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_renaming_method_parameters
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_renaming_method_parameters: true` under **linter > rules**:

```
linter:
 rules:
 avoid_renaming_method_parameters: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_return_types_on_setters linter rule.

Avoid return types on setters.

## Details

#**AVOID** return types on setters.

As setters do not return a value, declaring the return type of one is redundant.

**BAD:**

```
void set speed(int ms);
```
**GOOD:**

```
set speed(int ms);
```

## Enable

#
 To enable the `avoid_return_types_on_setters` rule, add `avoid_return_types_on_setters`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_return_types_on_setters
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_return_types_on_setters: true` under **linter > rules**:

```
linter:
 rules:
 avoid_return_types_on_setters: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_returning_null linter rule.

Avoid returning null from members whose return type is bool, double, int, or num.

## Details

#NOTE: This rule is removed in Dart 3.3.0; it is no longer functional.

 **AVOID** returning null from members whose return type is bool, double, int,
 or num.

Functions that return primitive types such as bool, double, int, and num are generally expected to return non-nullable values. Thus, returning null where a primitive type was expected can lead to runtime exceptions.

**BAD:**

```
bool getBool() => null;
num getNum() => null;
int getInt() => null;
double getDouble() => null;
```
**GOOD:**

```
bool getBool() => false;
num getNum() => -1;
int getInt() => -1;
double getDouble() => -1.0;
```

## Enable

#
 To enable the `avoid_returning_null` rule, add `avoid_returning_null` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_returning_null
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_returning_null: true` under **linter > rules**:

```
linter:
 rules:
 avoid_returning_null: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_returning_null_for_future linter rule.

Avoid returning null for Future.

## Details

#NOTE: This rule is removed in Dart 3.3.0; it is no longer functional.

**AVOID** returning null for Future.

 It is almost always wrong to return `null` for a `Future`. Most of the time the
 developer simply forgot to put an `async` keyword on the function.

## Enable

#
 To enable the `avoid_returning_null_for_future` rule, add `avoid_returning_null_for_future`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_returning_null_for_future
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_returning_null_for_future: true` under **linter > rules**:

```
linter:
 rules:
 avoid_returning_null_for_future: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_returning_null_for_void linter rule.

Avoid returning `null` for `void`.

## Details

#**AVOID** returning `null` for `void`.

 In a large variety of languages `void` as return type is used to indicate that
 a function doesn't return anything. Dart allows returning `null` in functions
 with `void` return type but it also allow using `return;` without specifying any
 value. To have a consistent way you should not return `null` and only use an
 empty return.

**BAD:**

```
void f1() {
 return null;
}
Future<void> f2() async {
 return null;
}
```
**GOOD:**

```
void f1() {
 return;
}
Future<void> f2() async {
 return;
}
```

## Enable

#
 To enable the `avoid_returning_null_for_void` rule, add `avoid_returning_null_for_void`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_returning_null_for_void
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_returning_null_for_void: true` under **linter > rules**:

```
linter:
 rules:
 avoid_returning_null_for_void: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_returning_this linter rule.

Avoid returning this from methods just to enable a fluent interface.

## Details

#From Effective Dart:

**AVOID** returning this from methods just to enable a fluent interface.

 Returning `this` from a method is redundant; Dart has a cascade operator which
 allows method chaining universally.

Returning `this` is allowed for:

- operators
- methods with a return type different of the current class
- methods defined in parent classes / mixins or interfaces
- methods defined in extensions

**BAD:**

```
var buffer = StringBuffer()
 .write('one')
 .write('two')
 .write('three');
```
**GOOD:**

```
var buffer = StringBuffer()
 ..write('one')
 ..write('two')
 ..write('three');
```

## Enable

#
 To enable the `avoid_returning_this` rule, add `avoid_returning_this` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_returning_this
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_returning_this: true` under **linter > rules**:

```
linter:
 rules:
 avoid_returning_this: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_setters_without_getters linter rule.

Avoid setters without getters.

## Details

#**DON'T** define a setter without a corresponding getter.

Defining a setter without defining a corresponding getter can lead to logical inconsistencies. Doing this could allow you to set a property to some value, but then upon observing the property's value, it could easily be different.

**BAD:**

```
class Bad {
 int l, r;
 set length(int newLength) {
 r = l + newLength;
 }
}
```
**GOOD:**

```
class Good {
 int l, r;
 int get length => r - l;
 set length(int newLength) {
 r = l + newLength;
 }
}
```

## Enable

#
 To enable the `avoid_setters_without_getters` rule, add `avoid_setters_without_getters`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_setters_without_getters
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_setters_without_getters: true` under **linter > rules**:

```
linter:
 rules:
 avoid_setters_without_getters: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_shadowing_type_parameters linter rule.

Avoid shadowing type parameters.

## Details

#**AVOID** shadowing type parameters.

**BAD:**

```
class A<T> {
 void fn<T>() {}
}
```
**GOOD:**

```
class A<T> {
 void fn<U>() {}
}
```

## Enable

#
 To enable the `avoid_shadowing_type_parameters` rule, add `avoid_shadowing_type_parameters`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_shadowing_type_parameters
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_shadowing_type_parameters: true` under **linter > rules**:

```
linter:
 rules:
 avoid_shadowing_type_parameters: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_single_cascade_in_expression_statements linter rule.

Avoid single cascade in expression statements.

## Details

#**AVOID** single cascade in expression statements.

**BAD:**

```
o..m();
```
**GOOD:**

```
o.m();
```

## Enable

#
 To enable the `avoid_single_cascade_in_expression_statements` rule, add `avoid_single_cascade_in_expression_statements`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_single_cascade_in_expression_statements
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_single_cascade_in_expression_statements: true` under **linter > rules**:

```
linter:
 rules:
 avoid_single_cascade_in_expression_statements: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_slow_async_io linter rule.

Avoid slow asynchronous `dart:io` methods.

## Details

#
 **AVOID** using the following asynchronous file I/O methods because they are
 much slower than their synchronous counterparts.

- `Directory.exists`
- `Directory.stat`
- `File.lastModified`
- `File.exists`
- `File.stat`
- `FileSystemEntity.isDirectory`
- `FileSystemEntity.isFile`
- `FileSystemEntity.isLink`
- `FileSystemEntity.type`

**BAD:**

```
import 'dart:io';
Future<Null> someFunction() async {
 var file = File('/path/to/my/file');
 var now = DateTime.now();
 if ((await file.lastModified()).isBefore(now)) print('before'); // LINT
}
```
**GOOD:**

```
import 'dart:io';
Future<Null> someFunction() async {
 var file = File('/path/to/my/file');
 var now = DateTime.now();
 if (file.lastModifiedSync().isBefore(now)) print('before'); // OK
}
```

## Enable

#
 To enable the `avoid_slow_async_io` rule, add `avoid_slow_async_io` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_slow_async_io
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_slow_async_io: true` under **linter > rules**:

```
linter:
 rules:
 avoid_slow_async_io: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_type_to_string linter rule.

Avoid

## Details

#
 **DO** avoid calls to

**BAD:**

```
void bar(Object other) {
 if (other.runtimeType.toString() == 'Bar') {
 doThing();
 }
}
Object baz(Thing myThing) {
 return getThingFromDatabase(key: myThing.runtimeType.toString());
}
```
**GOOD:**

```
void bar(Object other) {
 if (other is Bar) {
 doThing();
 }
}
class Thing {
 String get thingTypeKey => ...
}
Object baz(Thing myThing) {
 return getThingFromDatabase(key: myThing.thingTypeKey);
}
```

## Enable

#
 To enable the `avoid_type_to_string` rule, add `avoid_type_to_string` under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_type_to_string
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_type_to_string: true` under **linter > rules**:

```
linter:
 rules:
 avoid_type_to_string: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_types_as_parameter_names linter rule.

Avoid types as parameter names.

## Details

#**AVOID** using a parameter name that is the same as an existing type.

**BAD:**

```
m(f(int));
```
**GOOD:**

```
m(f(int v));
```

## Enable

#
 To enable the `avoid_types_as_parameter_names` rule, add `avoid_types_as_parameter_names`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_types_as_parameter_names
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_types_as_parameter_names: true` under **linter > rules**:

```
linter:
 rules:
 avoid_types_as_parameter_names: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_types_on_closure_parameters linter rule.

Avoid annotating types for function expression parameters.

## Details

#**AVOID** annotating types for function expression parameters.

Annotating types for function expression parameters is usually unnecessary because the parameter types can almost always be inferred from the context, thus making the practice redundant.

**BAD:**

```
var names = people.map((Person person) => person.name);
```
**GOOD:**

```
var names = people.map((person) => person.name);
```
## Incompatible rules

#The `avoid_types_on_closure_parameters` lint is incompatible with the following rules:

## Enable

#
 To enable the `avoid_types_on_closure_parameters` rule, add `avoid_types_on_closure_parameters`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_types_on_closure_parameters
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_types_on_closure_parameters: true` under **linter > rules**:

```
linter:
 rules:
 avoid_types_on_closure_parameters: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.

# avoid_

 Learn about the avoid_unnecessary_containers linter rule.

Avoid unnecessary containers.

## Details

#**AVOID** wrapping widgets in unnecessary containers.

 Wrapping a widget in `Container` with no other parameters set has no effect
 and makes code needlessly more complex.

**BAD:**

```
Widget buildRow() {
 return Container(
 child: Row(
 children: <Widget>[
 const MyLogo(),
 const Expanded(
 child: Text('...'),
 ),
 ],
 )
 );
}
```
**GOOD:**

```
Widget buildRow() {
 return Row(
 children: <Widget>[
 const MyLogo(),
 const Expanded(
 child: Text('...'),
 ),
 ],
 );
}
```

## Enable

#
 To enable the `avoid_unnecessary_containers` rule, add `avoid_unnecessary_containers`
 under
 **linter > rules** in your `analysis_options.yaml`
 file:

```
linter:
 rules:
 - avoid_unnecessary_containers
```
 If you're instead using the YAML map syntax to configure linter rules,
 add `avoid_unnecessary_containers: true` under **linter > rules**:

```
linter:
 rules:
 avoid_unnecessary_containers: true
```
Unless stated otherwise, the documentation on this site reflects Dart 3.12.2. Report an issue.
