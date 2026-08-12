# Rules and categories#

This page is the reference for every built-in breaking change rule and the four categories that group them.
For end-to-end usage and configuration, see the `buf breaking` usage guide and the `buf.yaml` reference.
Custom rules from Buf check plugins can run alongside or in place of these.

## Pick a category#

The four built-in categories below are designed as a strictness ladder.
Pick one and stick with it; mixing and excluding individual rules is allowed but rarely a good idea, since the rules in each category are chosen to give a coherent compatibility guarantee.
If you don’t know which to pick, start with `FILE` (the default); you can always relax later.

## Categories#

Buf’s breaking rules fit under four categories, from strictest to most lenient:

- `FILE`:- **Default.**Detects changes that move generated code between files, breaking generated source code on a per-file basis. This breaks generated stubs in some languages; for example, it’s safe to move code between files in Go but not in Python.
- `PACKAGE`: Detects changes that break generated source code changes on a per-package basis. It detects changes that break the generated stubs, but only accounting for package-level changes.
- `WIRE_JSON`: Detects changes that break wire (binary) or JSON encoding. Because JSON is common across many transports, this is the recommended minimum level.
- `WIRE`: Detects changes that break wire (binary) encoding.

For details on individual rules and the categories they belong to, see the Rules section below.

`FILE` and `PACKAGE`#

The `FILE` and `PACKAGE` categories protect compatibility in generated code.
For example, deleting an enum or message often removes the corresponding type in generated code.
Any code that refers to that enum or message then fails to compile.

As an example, imagine that you have an `Arena` enum and mark `ARENA_FOO` as deprecated:

Later you remove the field, because it’s no longer supported by the server:

This change is perfectly wire compatible, but all code that referred to `ARENA_FOO` then fails to compile:

In some cases this is desirable, but more commonly you’re sharing your `.proto` files or generated code to clients that you don’t control.
You should choose `FILE` or `PACKAGE` breaking detection if you want to know when you’ll break your client’s code.

Though these rules are code generator specific, you should use `FILE` to protect all generated languages.
`FILE` is absolutely necessary for C++ and Python.

You can use `PACKAGE` to protect languages that are less sensitive to types moving between files within the same package, like Go.

`WIRE` and `WIRE_JSON`#

`WIRE` and `WIRE_JSON` detect breakage of encoded messages.
For example:

- Changing an optional field into a required one. Old messages that don’t have that field encoded fail to read in the new definition.
- Reserving deleted types for which reuse in the future could cause wire incompatibilities.

`WIRE` and `WIRE_JSON` don’t check for breakage in generated source code.
This is advantageous when:

- You control all of your clients for your service. You’re fixing it if it breaks anyway.
- You want your client’s build to break instead of getting errors at runtime. (This assumes your clients are equally happy to immediately stop what they’re doing to fix your service.)
- All of your clients are in a monorepo. You want to determine who’s depending on deprecated features by a broken build instead of at runtime.
- You’re your own client. For example, you’re trying to detect issues reading Protobuf encoded messages from older versions of your program that were persisted to disk or other non-volatile storage.

Using `WIRE_JSON` instead of `WIRE` is safer because Protobuf’s JSON encoding breaks when field names change.

- Use `WIRE_JSON`if you’re using Connect, gRPC-Gateway, or gRPC JSON.
- Use the less strict `WIRE`when you can guarantee that only binary encoded messages are decoded.

## Rules#

The rules are grouped below based on the kind of breaking change they check for: deletions, sameness, and changes to Protobuf file options. Each rule lists the categories that include it.

`ENUM_NO_DELETE`#

**Category:** `FILE`

This checks that no enums are deleted from a given file. Deleting an enum deletes the corresponding generated type, which could be referenced in source code. Instead of deleting an enum, deprecate it:

`ENUM_SAME_JSON_FORMAT`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`

This checks that an enum doesn’t change from supporting the JSON format to “best effort”.
Enums in `proto2` files are best effort (since the JSON format wasn’t defined when the `proto2` syntax was created).
Enums in `proto3` files support the JSON format, and enums in Editions files can be configured using the `json_format` feature.

This typically doesn’t impact anything except internal compiler validation. When the JSON format is supported, the compiler does more strict checks to prevent name collisions in the JSON names for fields and enum values. However, it could impact the output of code generation plugins and it does impact how the schema can change in the future, possibly allowing changes to a message or an enum that aren’t amenable to using JSON encoding.

`ENUM_SAME_TYPE`#

**Categories:** `FILE`, `PACKAGE`

This checks that an enum doesn’t change from open to closed or vice versa, because whether an enum is open or closed can impact code generation.
Enums in `proto2` files are closed, which means that unrecognized values result in the field being unset (the actual value is stored with other unrecognized fields).
Enums in `proto3` files are open, which means that values not defined in the schema are accepted.
Enums in Editions files default to open but can be configured using the `enum_type` feature.

`ENUM_VALUE_NO_DELETE`#

**Categories:** `FILE`, `PACKAGE`

This checks that no enum value is deleted. Deleting an enum value results in the corresponding value or field being deleted from the generated source code, which could be referenced. Instead of deleting a value, deprecate it:

`ENUM_VALUE_NO_DELETE_UNLESS_NAME_RESERVED`#

**Category:** `WIRE_JSON`

This checks that no enum value is deleted without reserving the name. This is the JSON equivalent of reserving the number, since JSON can use field names rather than numbers for enum values (optional, but allowed). Reserving both the number and the name is preferable in most cases. Here’s an example:

Note that it’s usually better to deprecate enum values than to reserve them in advance.

`ENUM_VALUE_NO_DELETE_UNLESS_NUMBER_RESERVED`#

**Categories:** `WIRE`, `WIRE_JSON`

This checks that no enum value is deleted without reserving the number. Though deleting an enum value isn’t directly a wire-breaking change, reusing these numbers in the future is likely to result in bugs. This is also a JSON breaking change for enum values if they’re serialized as integers (which is an option). Protobuf provides the ability to reserve numbers to prevent them from being reused in the future. For example:

`ENUM_VALUE_SAME_NAME`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`

This checks that a given enum value has the same name for each enum value number.
For example You can’t change `FOO_ONE = 1` to `FOO_TWO = 1`.
Doing so results in potential JSON incompatibilities and broken source code.

Note that for enums with `allow_alias` set, this verifies that the set of names in the current definition covers the set of names in the previous definition.
For example, the new definition `// new` is compatible with `// old`, but `// old` isn’t compatible with `// new`:

```
// old
enum Foo {
 option allow_alias = 1;
 FOO_UNSPECIFIED = 0;
 FOO_BAR = 1;
 FOO_BARR = 1;
}
// new
enum Foo {
 option allow_alias = 1;
 FOO_UNSPECIFIED = 0;
 FOO_BAR = 1;
 FOO_BARR = 1;
 FOO_BARRR = 1;
}
```
`EXTENSION_MESSAGE_NO_DELETE`#

**Categories:** `FILE`, `PACKAGE`

This checks that no extension range is deleted from any message. Though this won’t have any effect on your generated source code, deleting an extension range can result in compile errors for downstream Protobuf schemas, and is generally not recommended.

Note that extension ranges can’t be defined in `proto3` files, so this only impacts sources that use `proto2` syntax or Editions.

`EXTENSION_NO_DELETE`#

**This rule can only be used with v2 configuration files.**

**Category:** `FILE`

This checks that no extensions are deleted from a given file, identified by their fully qualified name. Deleting an extension deletes the corresponding generated extension type, which could be referenced in source code. Instead of deleting an extension, deprecate it:

`FIELD_NO_DELETE`#

**Categories:** `FILE`, `PACKAGE`

This checks that no message field is deleted. Deleting a message field results in the corresponding value or field being deleted from the generated source code, which could be referenced. Instead of deleting a value, deprecate it:

Unlike other field checks, this rule *doesn’t* apply to extensions.

`FIELD_NO_DELETE_UNLESS_NAME_RESERVED`#

**Category:** `WIRE_JSON`

This checks that no message field is deleted without reserving the name. This is the JSON equivalent of reserving the number, since JSON uses field names instead of numbers. Reserving both the number and the name is preferable in most cases:

Note that it’s usually better to deprecate message fields than to reserve them in advance.

Unlike other field checks, this rule *doesn’t* apply to extensions.

`FIELD_NO_DELETE_UNLESS_NUMBER_RESERVED`#

**Categories:** `WIRE`, `WIRE_JSON`

This checks that message field is deleted without reserving the number. Though deleting a message field isn’t directly a wire-breaking change, reusing these numbers in the future is likely to result wire incompatibilities if the type differs. Protobuf provides the ability to reserve numbers to prevent them from being reused in the future. For example:

Note that deprecating a field instead of deleting it has the same effect as reserving the field (as well as reserving the name for JSON).

Unlike other field checks, this rule *doesn’t* apply to extensions.

`FIELD_SAME_CARDINALITY`#

**Categories:** `FILE`, `PACKAGE`

This checks that no field changes its cardinality. The available cardinalities are:

- **optional with implicit presence:**This is the cardinality of fields in- `proto3`files that don’t explicitly specify a label (excluding extension fields and fields in a oneof). It can be enabled in an Editions source files by setting the- `field_presence`feature to- `IMPLICIT`.
- **optional with explicit presence:**This is the cardinality of optional fields in- `proto2`files and fields in- `proto3`files that explicitly use the- `optional`label. It’s also the cardinality for all fields in oneofs and all extensions. This is also the default cardinality for optional fields in Editions source files.
- **required:**This is the cardinality of fields in- `proto2`files that use the- `required`label. It also applies to fields in Editions source files that set the- `field_presence`feature to- `LEGACY_REQUIRED`.
- **repeated:**This is the cardinality of fields that use the- `repeated`label.
- **map:**This is the cardinality of map fields. Under the hood, this is similar to repeated cardinality, except that entries are de-duplicated using the map key.

Though changing an optional field from implicit to explicit presence is typically backwards-compatible, in some runtimes it results in different, incompatible generated code.

For files with `proto3` syntax, message fields already have explicit field presence.
Because of this, adding the `optional` modifier doesn’t change its cardinality and there’s no difference in representation.

`FIELD_SAME_CPP_STRING_TYPE`#

**Categories:** `FILE`, `PACKAGE`

This checks that a given string or bytes field uses the same type in generated C++ code.
For files with `proto2` and `proto3` syntax, this comes from the value for the same value for the `ctype` option.
In edition 2023, this can optionally instead be defined using the C++-specific feature `(pb.cpp).string_type`.

The `ctype` option is a Google-internal field option, so generally you won’t have it set.
The new `(pb.cpp).string_type` feature is intended to be used to migrate generated code to using the more efficient `absl::string_view` type.

`FIELD_SAME_CTYPE`#

**This rule is deprecated.**

It has been replaced with the `FIELD_SAME_CPP_STRING_TYPE` rule.
It isn’t a valid rule name when used with `v2` configuration files.
When using earlier configuration versions, it’s treated as an alias for `FIELD_SAME_CPP_STRING_TYPE`.

`FIELD_SAME_DEFAULT`#

**This rule can only be used with v2 configuration files.**

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`, `WIRE`

This checks that the default value for a field doesn’t change.
In `proto3` files, the default value is always the zero value for the type, but in `proto2` and in Editions, non-repeated, non-message fields can configure an alternate default value.
Changing the default value isn’t backwards-compatible since it means that producers and consumers of the schema interpret serialized data differently.

`FIELD_SAME_JAVA_UTF8_VALIDATION`#

**Categories:** `FILE`, `PACKAGE`

This checks that a given string field uses the same level of UTF8 verification in generated Java code.
With `proto3` syntax, string fields are always validated.
In `proto2` syntax, they aren’t validated by default, but you can opt in to validation by using the `java_string_check_utf8` file option.

In Editions, this can now be controlled on a per-field level.
This allows a `proto2` file to be migrated to Editions and then incrementally updated to using runtime UTF8 verification.
This is done by setting a per-field feature: `(pb.java).utf8_validation`.

Note that there is also a global feature for this named simply `utf8_validation`, which is intended to enable the verification across all languages and runtimes.
Also see `FIELD_SAME_UTF8_VALIDATION`.

`FIELD_SAME_JSON_NAME`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`

This checks that the `json_name` field option doesn’t change, which breaks JSON compatibility.
Though it’s not usually a generated source code breaking change, some Protobuf plugins may generate code based on this option.
Having this as part of the `FILE` and `PACKAGE` groups also fulfills that the `FILE` and `PACKAGE` categories are supersets of the `WIRE_JSON` category.

`FIELD_SAME_JSTYPE`#

**Categories:** `FILE`, `PACKAGE`

This checks that a given field has the same value for the `jstype` option.
This affects JavaScript generated code.

`FIELD_SAME_LABEL`#

**This rule is deprecated.**

It has been replaced with the following checks:

- `FIELD_SAME_CARDINALITY`(in the- `FILE`and- `PACKAGE`categories)
- `FIELD_WIRE_COMPATIBLE_CARDINALITY`(in the- `WIRE`category)
- `FIELD_WIRE_JSON_COMPATIBLE_CARDINALITY`(in the- `WIRE_JSON`category)

It isn’t a valid rule name when used with `v2` configuration files.

When using `v1` configuration files, it’s treated as an alias for all three checks listed previously.
When using `v1beta1` configuration files, it’s treated as an alias for `FIELD_SAME_CARDINALITY`.

`FIELD_SAME_NAME`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`

This checks that the field name for a given field number doesn’t change.
For example, you can’t change `int64 foo = 1;` to `int64 bar = 1;`.
This affects generated source code, but also affects JSON compatibility because JSON uses field names for serialization.

`FIELD_SAME_ONEOF`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`, `WIRE`

This checks that no field moves into or out of a oneof or changes the oneof it’s a part of. Doing so is almost always a generated source code breaking change. Technically there are exceptions with regard to wire compatibility, but the rules are complex enough that it’s safer to never change a field’s presence inside or outside a given oneof.

`FIELD_SAME_TYPE`#

**Categories:** `FILE`, `PACKAGE`

This checks that a field has the same type. Changing the type of a field can affect the type in the generated source code, wire compatibility, and JSON compatibility. Note that technically, it’s possible to interchange some scalar types. However, most of these result in generated source code changes anyway, and affect JSON compatibility. Instead of worrying about this, just don’t change your field types.

Note that with maps, you may get slightly confusing error messages when changing a field to or from a map and some other type, denoting that the cardinality of the field changed from `repeated` to `map` or the message changed type from `message` to another type.
This is because of the way maps are implemented in Protobuf, where every map is actually just a `repeated` field of an implicit message.
Buf still properly detects this change and outputs an error, so the pass/fail decision remains the same.

`FIELD_SAME_UTF8_VALIDATION`#

**Categories:** `FILE`, `PACKAGE`

This checks that a given string field uses the same level of UTF8 verification at runtime.
With `proto2` syntax, string field contents *aren’t* validated at runtime.
With `proto3` syntax, string fields are always validated.
String fields in Editions files default to being validated at runtime, but that can be changed using the `utf8_validation` feature.

`FIELD_WIRE_COMPATIBLE_CARDINALITY`#

**Categories:** `WIRE`

This rule replaces `FIELD_SAME_CARDINALITY` for the `WIRE` category.
The consequences of this rule are:

- If cardinality changes between “optional with implicit presence” and “optional with explicit presence”, the check passes. Field presence doesn’t impact the binary wire format.
- If cardinality changes between repeated and map, the check passes. Map fields are encoded in the binary wire format as a repeated field of messages.

`FIELD_WIRE_COMPATIBLE_TYPE`#

**Categories:** `WIRE`

This rule replaces `FIELD_SAME_TYPE` for the `WIRE` category.
The consequences of this rule are:

- If the type changed between int32, uint32, int64, uint64, and bool, the check passes.
- If the type changed between sint32 and sint64, the check passes.
- If the type changed between fixed32 and sfixed32, the check passes.
- If the type changed between fixed64 and sfixed64, the check passes.
- If the type changed from string to bytes, the check passes.
- If the type changed from bytes to string, the check produces an error about string and bytes compatibility. Per the Protobuf docs, you can change between string and bytes IF the data is valid UTF-8, but because we’re only concerned with the API definition and can’t know how a user actually uses the field, the check fails.
- If the previous and current types are both enums, Buf checks them to see if (1) the short names are equal, and (2) the previous enum is a subset of the current enum. A subset is defined as having a subset of the name/number enum values. If the previous enum is a subset, the check passes. This covers the case where someone moves where an enum is defined, but still allows values to be added to this enum in the same change, because adding values to an enum isn’t a breaking change.
- A link to https://developers.google.com/protocol-buffers/docs/proto3#updating is added to failures produced from `FIELD_WIRE_COMPATIBLE_TYPE`.

`FIELD_WIRE_JSON_COMPATIBLE_CARDINALITY`#

**Categories:** `WIRE_JSON`

This rule replaces `FIELD_SAME_CARDINALITY` for the `WIRE_JSON` category.
The consequences of this rule are:

- If cardinality changes between “optional with implicit presence” and “optional with explicit presence”, the check passes. Field presence doesn’t impact the JSON format.

Unlike `FIELD_WIRE_COMPATIBLE_CARDINALITY`, the check still fails if a field changes between “repeated” and “map”, as one uses JSON arrays and the other uses JSON maps in their JSON formats.

`FIELD_WIRE_JSON_COMPATIBLE_TYPE`#

**Categories:** `WIRE_JSON`

This rule replaces `FIELD_SAME_TYPE` for the `WIRE_JSON` category.

JSON allows for some exchanging of types, but due to the way various fields are serialized, the rules are stricter (see the Protocol Buffer docs).
For example, int32, sint32, and uint32 can be exchanged, but 64-bit numbers have a different representation in JSON.
Since sint32 isn’t compatible with int32 or uint32 in `WIRE`, limit this to allow int32 and uint32 to be exchanged in JSON.

The consequences of this rule are:

- If the type changes between int32 and uint32, the check passes.
- If the type changes between int64 and uint64, the check passes.
- If the type changes between fixed32 and sfixed32, the check passes.
- If the type changes between fixed64 and sfixed64, the check passes.
- If the previous and current types are both enums, Buf checks them to see if (1) the short names are equal, and (2) the previous enum is a subset of the current enum. A subset is defined as having a subset of the name/number enum values. If the previous enum is a subset, the check passes. This covers the case where someone moves where an enum is defined, but still allows values to be added to this enum in the same change, because adding values to an enum isn’t a breaking change.
- Links to https://developers.google.com/protocol-buffers/docs/proto3#updating and https://developers.google.com/protocol-buffer/docs/proto3#json are added to failures produced from `FIELD_WIRE_JSON_COMPATIBLE_TYPE`.

`FILE_NO_DELETE`#

**Category:** `FILE`

This checks that no file is deleted. Deleting a file results in its generated header file being deleted as well, which could break source code.

`FILE_SAME_CC_ENABLE_ARENAS`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_CC_GENERIC_SERVICES`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_CSHARP_NAMESPACE`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_GO_PACKAGE`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_JAVA_GENERIC_SERVICES`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_JAVA_MULTIPLE_FILES`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_JAVA_OUTER_CLASSNAME`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_JAVA_PACKAGE`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_JAVA_STRING_CHECK_UTF8`#

**This rule is deprecated.**

It has been replaced with the `FIELD_SAME_JAVA_UTF8_VALIDATION` rule.
It isn’t a valid rule name when used with `v2` configuration files.
When using earlier configuration versions, it’s treated as an alias for `FIELD_SAME_JAVA_UTF8_VALIDATION`.

`FILE_SAME_OBJC_CLASS_PREFIX`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_OPTIMIZE_FOR`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_PACKAGE`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`, `WIRE`

This checks that a given file has the same `package` value.
Changing the package value results in a ton of issues downstream in various languages, and for the `FILE` category, this effectively results in any types declared within that file being considered deleted.

`FILE_SAME_PHP_CLASS_PREFIX`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_PHP_GENERIC_SERVICES`#

**As of v1.32.0 of buf (and v26.0 of Protobuf), there is no longer a php_generic_services file option.
This rule is deprecated and has no replacement.**

It isn’t a valid rule name when used with `v2` configuration files.
When using earlier configuration versions, it’s effectively ignored.

`FILE_SAME_PHP_METADATA_NAMESPACE`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_PHP_NAMESPACE`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_PY_GENERIC_SERVICES`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_RUBY_PACKAGE`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_SWIFT_PREFIX`#

**Categories:** `FILE`, `PACKAGE`

This checks that the value of this file option doesn’t change values between versions of your Protobuf schema. Changing this value results in differences in your generated source code.

`FILE_SAME_SYNTAX`#

**Categories:** `FILE`, `PACKAGE`

This checks that the file syntax doesn’t change between `proto2`, `proto3`, and Editions.
Changing the syntax may result in differences in generated code for some languages.

For many plugins, including the core Google-provided code generators, the syntax doesn’t necessarily impact code generation but instead the *syntax-specific semantics* do. These syntax-specific semantics are also validated in other rules:

So it’s often okay to ignore this `FILE_SAME_SYNTAX` rule and leave the others enabled, which allows you to migrate your sources from `proto2` or `proto3` to Editions without getting errors from breaking change detection (as long as your migration isn’t actually changing semantics).

However, it depends on what code generation plugins are used.

`MESSAGE_NO_DELETE`#

**Category:** `FILE`

This checks that no messages are deleted from a given file. Deleting a message deletes the corresponding generated type, which could be referenced in source code. Instead of deleting a message, deprecate it:

`MESSAGE_NO_REMOVE_STANDARD_DESCRIPTOR_ACCESSOR`#

**Categories:** `FILE`, `PACKAGE`

This checks that the `no_standard_descriptor_accessor` message option isn’t changed from `false` or unset to `true`.
Changing this option to `true` results in the `descriptor()` accessor not being generated in certain languages, which is a generated source code breaking change.
Protobuf has issues with fields that are named “descriptor”, with any capitalization and with any number of underscores
before and after “descriptor”.
Don’t name fields this.

`MESSAGE_SAME_JSON_FORMAT`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`

This checks that a message doesn’t change from supporting the JSON format to “best effort”.
Messages in `proto2` files are best effort (since the JSON format wasn’t defined when the `proto2` syntax was created).
Messages in `proto3` files support the JSON format, and messages in Editions files can be configured using the `json_format` feature.

This typically doesn’t impact anything except internal compiler validation: when the JSON format is supported, the compiler does more strict checks to prevent name collisions in the JSON names for fields and enum values. However, it could impact the output of code generation plugins, and it does impact how the schema can change in the future, possibly allowing changes to a message or an enum that aren’t amenable to using JSON encoding.

`MESSAGE_SAME_MESSAGE_SET_WIRE_FORMAT`#

**This rule is deprecated and has no replacement.**

This rule is no longer valid and is effectively ignored since `protobuf-go` has removed support for `MessageSets`.
Formerly, this checked that the `message_set_wire_format` message option is the same.
Since this is a `proto1` construct that’s long deprecated, we congratulate you if you are using this for any current Protobuf schema, as you are a champion of maintaining backwards compatible APIs over many years.
Instead of failing breaking change detection, perhaps you should get an award. 🏆

`MESSAGE_SAME_REQUIRED_FIELDS`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`, `WIRE`

This checks that messages have no added or deleted required fields. Note that required fields are considered strongly deprecated since it’s impossible to change or remove a required field in a safe way. All existing required fields should be treated as permanent and immutable.

`ONEOF_NO_DELETE`#

**Categories:** `FILE`, `PACKAGE`

This checks that no oneof is deleted from a message. Various languages generate types for oneofs, which would no longer be present if the oneof is deleted.

`PACKAGE_ENUM_NO_DELETE`#

**Category:** `PACKAGE`

This has the same effect as `ENUM_NO_DELETE`, except that it verifies that types aren’t deleted from a given package, while letting them move between files in the same package.

`PACKAGE_EXTENSION_NO_DELETE`#

**This rule can only be used with v2 configuration files.**

**Category:** `PACKAGE`

This has the same effect as `EXTENSION_NO_DELETE`, except that it verifies that types aren’t deleted from a given package, while letting them move between files in the same package.

`PACKAGE_MESSAGE_NO_DELETE`#

**Category:** `PACKAGE`

This has the same effect as `MESSAGE_NO_DELETE`, except that it verifies that types aren’t deleted from a given package, while letting them move between files in the same package.

`PACKAGE_NO_DELETE`#

**Category:** `PACKAGE`

This checks that every package that existed in your previous version still exists in the current schemas. Deleting a package usually deletes other types that break generated code.

`PACKAGE_SERVICE_NO_DELETE`#

**Category:** `PACKAGE`

This has the same effect as `SERVICE_NO_DELETE`, except that it verifies that types aren’t deleted
from a given package, while letting them move between files in the same package.

`RESERVED_ENUM_NO_DELETE`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`, `WIRE`

This checks that no reserved number range or reserved name is deleted from any enum. Deleting a reserved value means that future versions of your Protobuf schema could use names or numbers in those ranges, and if the ranges were reserved, it was probably because an enum value was deleted.

`RESERVED_MESSAGE_NO_DELETE`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`, `WIRE`

This checks that no reserved number range or reserved name is deleted from any message. Deleting a reserved value means that future versions of your Protobuf schema could use names or numbers in those ranges, and if the ranges were reserved, it was probably because a field was deleted.

`RPC_NO_DELETE`#

**Categories:** `FILE`, `PACKAGE`

This checks that no RPC is deleted from a service. Doing so isn’t a wire-breaking change (although client calls fail if a server doesn’t implement a given RPC); however, existing source code may reference a given RPC. Instead of deleting an RPC, deprecate it.

`RPC_SAME_CLIENT_STREAMING`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`, `WIRE`

This checks that RPC signatures don’t change. Doing so breaks both generated source code and over-the-wire RPC calls.

`RPC_SAME_IDEMPOTENCY_LEVEL`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`, `WIRE`

This checks that the `idempotency_level`
RPC option doesn’t change.
Doing so can result in different HTTP verbs being used.

`RPC_SAME_REQUEST_TYPE`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`, `WIRE`

This checks that RPC signatures don’t change. Doing so breaks both generated source code and over-the-wire RPC calls.

`RPC_SAME_RESPONSE_TYPE`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`, `WIRE`

This checks that RPC signatures don’t change. Doing so breaks both generated source code and over-the-wire RPC calls.

`RPC_SAME_SERVER_STREAMING`#

**Categories:** `FILE`, `PACKAGE`, `WIRE_JSON`, `WIRE`

This checks that RPC signatures don’t change. Doing so breaks both generated source code and over-the-wire RPC calls.

`SERVICE_NO_DELETE`#

**Category:** `FILE`

This checks that no services are deleted from a given file. Deleting a service deletes the corresponding generated type, which could be referenced in source code. Instead of deleting a service, deprecate it:

## What’s not covered#

The rules above cover every field on a standard `FileDescriptorSet`, tracked against current Protobuf releases as new fields are added.
If something is missing, let us know.

Custom options are intentionally out of scope: their effects are user-defined, so a generic compatibility rule can’t reason about them. For project-specific compatibility checks (including custom-option semantics), write a Buf check plugin.
