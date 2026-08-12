It’s common for even the best programmers to make simple mistakes. And sometimes a refactoring which seems safe can leave behind code which will never do what’s intended.

We’re used to getting help from the compiler, but it doesn’t do much beyond static type checking. Using Error Prone to augment the compiler’s type analysis, you can catch more mistakes before they cost you time, or end up as bugs in production. We use Error Prone in Google’s Java build system to eliminate classes of serious bugs from entering our code, and we’ve open-sourced it, so you can too!

**Error Prone …**

```
import java.util.Set;
import java.util.HashSet;
public class ShortSet {
 public static void main (String[] args) {
 Set<Short> s = new HashSet<>();
 for (short i = 0; i < 100; i++) {
 s.add(i);
 s.remove(i - 1);
 }
 System.out.println(s.size());
 }
}
```
```
$ bazel build :hello
ERROR: example/myproject/BUILD:29:1: Java compilation in rule '//example/myproject:hello'
ShortSet.java:6: error: [CollectionIncompatibleType] Argument 'i - 1' should not be passed to this method;
its type int is not compatible with its collection's type argument Short
 s.remove(i - 1);
 ^
 (see https://errorprone.info/bugpattern/CollectionIncompatibleType)
1 error
```
Doug Lea, on learning of a bug we discovered in ConcurrentHashMap:

Definitely embarrassing. I guess I’m back to liking Error Prone even though it sometimes annoys me :-)

It’s common for even the best programmers to make simple mistakes. And sometimes a refactoring which seems safe can leave behind code which will never do what’s intended.

We’re used to getting help from the compiler, but it doesn’t do much beyond static type checking. Using Error Prone to augment the compiler’s type analysis, you can catch more mistakes before they cost you time, or end up as bugs in production. We use Error Prone in Google’s Java build system to eliminate classes of serious bugs from entering our code, and we’ve open-sourced it, so you can too!

**Error Prone …**

```
import java.util.Set;
import java.util.HashSet;
public class ShortSet {
 public static void main (String[] args) {
 Set<Short> s = new HashSet<>();
 for (short i = 0; i < 100; i++) {
 s.add(i);
 s.remove(i - 1);
 }
 System.out.println(s.size());
 }
}
```
```
$ bazel build :hello
ERROR: example/myproject/BUILD:29:1: Java compilation in rule '//example/myproject:hello'
ShortSet.java:6: error: [CollectionIncompatibleType] Argument 'i - 1' should not be passed to this method;
its type int is not compatible with its collection's type argument Short
 s.remove(i - 1);
 ^
 (see https://errorprone.info/bugpattern/CollectionIncompatibleType)
1 error
```
Doug Lea, on learning of a bug we discovered in ConcurrentHashMap:

Definitely embarrassing. I guess I’m back to liking Error Prone even though it sometimes annoys me :-)

This list is auto-generated from our sources.

Patterns which are marked **Experimental** will not be evaluated against your
code, unless you specifically configure Error Prone. The default checks are
marked **On by default**, and each release promotes some experimental checks
after we’ve vetted them against Google’s codebase.

**AlwaysThrows**

Statically detects calls that will always fail at runtime.

**AndroidInjectionBeforeSuper**

AndroidInjection.inject() should always be invoked before calling super.lifecycleMethod()

**ArrayEquals**

Reference equality used to compare arrays

**ArrayFillIncompatibleType**

Arrays.fill(Object[], Object) called with incompatible types.

**ArrayHashCode**

hashcode method on array does not hash array contents

**ArrayToString**

Calling toString on an array does not provide useful information

**ArraysAsListPrimitiveArray**

Arrays.asList does not autobox primitive arrays, as one might expect.

**AsyncCallableReturnsNull**

AsyncCallable should not return a null Future, only a Future whose result is null.

**AsyncFunctionReturnsNull**

AsyncFunction should not return a null Future, only a Future whose result is null.

**AutoValueBuilderDefaultsInConstructor**

Defaults for AutoValue Builders should be set in the factory method returning Builder instances, not the constructor

**AutoValueConstructorOrderChecker**

Arguments to AutoValue constructor are in the wrong order

**BadAnnotationImplementation**

Classes that implement Annotation must override equals and hashCode. Consider using AutoAnnotation instead of implementing Annotation by hand.

**BadShiftAmount**

Shift by an amount that is out of range

**BanJNDI**

Using JNDI may deserialize user input via the `Serializable` API which is extremely dangerous

**BoxedPrimitiveEquality**

Comparison using reference equality instead of value equality. Reference equality of boxed primitive types is usually not useful, as they are value objects, and it is bug-prone, as instances are cached for some values but not others.

**BundleDeserializationCast**

Object serialized in Bundle may have been flattened to base type.

**ChainingConstructorIgnoresParameter**

The called constructor accepts a parameter with the same name and type as one of its caller's parameters, but its caller doesn't pass that parameter to it. It's likely that it was intended to.

**CheckNotNullMultipleTimes**

A variable was checkNotNulled multiple times. Did you mean to check something else?

**CheckReturnValue**

The result of this call must be used

**CollectionIncompatibleType**

Incompatible type as argument to Object-accepting Java collections method

**CollectionToArraySafeParameter**

The type of the array parameter of Collection.toArray needs to be compatible with the array type

**ComparableType**

Implementing 'Comparable<T>' where T is not the same as the implementing class is incorrect, since it violates the symmetry contract of compareTo.

**ComparingThisWithNull**

this == null is always false, this != null is always true

**ComparisonOutOfRange**

Comparison to value that is out of range for the compared type

**CompatibleWithAnnotationMisuse**

@CompatibleWith's value is not a type argument.

**CompileTimeConstant**

Non-compile-time constant expression passed to parameter with @CompileTimeConstant type annotation.

**ComputeIfAbsentAmbiguousReference**

computeIfAbsent passes the map key to the provided class's constructor

**ConditionalExpressionNumericPromotion**

A conditional expression with numeric operands of differing types will perform binary numeric promotion of the operands; when these operands are of reference types, the expression's result may not be of the expected type.

**ConstantOverflow**

Compile-time constant expression overflows

**DaggerProvidesNull**

Dagger @Provides methods may not return null unless annotated with @Nullable

**DangerousLiteralNull**

This method is null-hostile: passing a null literal to it is always wrong

**DeadException**

Exception created but not thrown

**DeadThread**

Thread created but not started

**DereferenceWithNullBranch**

Dereference of an expression with a null branch

**DiscardedPostfixExpression**

The result of this unary operation on a lambda parameter is discarded

**DoNotCall**

This method should not be called.

**DoNotMock**

Identifies undesirable mocks.

**DoubleBraceInitialization**

Prefer collection factory methods or builders to the double-brace initialization pattern.

**DuplicateMapKeys**

Map#ofEntries will throw an IllegalArgumentException if there are any duplicate keys

**DurationFrom**

Duration.from(Duration) returns itself; from(Period) throws a runtime exception.

**DurationGetTemporalUnit**

Duration.get() only works with SECONDS or NANOS.

**DurationTemporalUnit**

Duration APIs only work for DAYS or exact durations.

**DurationToLongTimeUnit**

Unit mismatch when decomposing a Duration or Instant to call a <long, TimeUnit> API

**EqualsHashCode**

Classes that override equals should also override hashCode.

**EqualsNaN**

== NaN always returns false; use the isNaN methods instead

**EqualsNull**

The contract of Object.equals() states that for any non-null reference value x, x.equals(null) should return false. If x is null, a NullPointerException is thrown. Consider replacing equals() with the == operator.

**EqualsReference**

== must be used in equals method to check equality to itself or an infinite loop will occur.

**EqualsWrongThing**

Comparing different pairs of fields/getters in an equals implementation is probably a mistake.

**FloggerFormatString**

Invalid printf-style format string

**FloggerLogString**

Arguments to log(String) must be compile-time constants or parameters annotated with @CompileTimeConstant. If possible, use Flogger's formatting log methods instead.

**FloggerLogVarargs**

logVarargs should be used to pass through format strings and arguments.

**FloggerSplitLogStatement**

Splitting log statements and using Api instances directly breaks logging.

**ForOverride**

Method annotated @ForOverride must be protected or package-private and only invoked from declaring class, or from an override of the method

**FormatString**

Invalid printf-style format string

**FormatStringAnnotation**

Invalid format string passed to formatting method.

**FromTemporalAccessor**

Certain combinations of javaTimeType.from(TemporalAccessor) will always throw a DateTimeException or return the parameter directly.

**FunctionalInterfaceMethodChanged**

Casting a lambda to this @FunctionalInterface can cause a behavior change from casting to a functional superinterface, which is surprising to users. Prefer decorator methods to this surprising behavior.

**FuturesGetCheckedIllegalExceptionType**

Futures.getChecked requires a checked exception type with a standard constructor.

**FuzzyEqualsShouldNotBeUsedInEqualsMethod**

DoubleMath.fuzzyEquals should never be used in an Object.equals() method

**GetClassOnAnnotation**

Calling getClass() on an annotation may return a proxy class

**GetClassOnClass**

Calling getClass() on an object of type Class returns the Class object for java.lang.Class; you probably meant to operate on the object directly

**GuardedBy**

Checks for unguarded accesses to fields and methods with @GuardedBy annotations

**GuiceAssistedInjectScoping**

Scope annotation on implementation class of AssistedInject factory is not allowed

**GuiceAssistedParameters**

A constructor cannot have two @Assisted parameters of the same type unless they are disambiguated with named @Assisted annotations.

**GuiceInjectOnFinalField**

Although Guice allows injecting final fields, doing so is disallowed because the injected value may not be visible to other threads.

**HashtableContains**

contains() is a legacy method that is equivalent to containsValue()

**IdentityBinaryExpression**

A binary expression where both operands are the same is usually incorrect.

**IdentityHashMapBoxing**

Using IdentityHashMap with a boxed type as the key is risky since boxing may produce distinct instances

**Immutable**

Type declaration annotated with @Immutable is not immutable

**ImpossibleNullComparison**

This value cannot be null, and comparing it to null may be misleading.

**Incomparable**

Types contained in sorted collections must implement Comparable.

**IncompatibleArgumentType**

Passing argument to a generic method with an incompatible type.

**IncompatibleModifiers**

This annotation has incompatible modifiers as specified by its @IncompatibleModifiers annotation

**IndexOfChar**

The first argument to indexOf is a Unicode code point, and the second is the index to start the search from

**InexactVarargsConditional**

Conditional expression in varargs call contains array and non-array arguments

**InfiniteRecursion**

This method calls itself unconditionally; it will throw StackOverflowError

**InjectMoreThanOneScopeAnnotationOnClass**

A class can be annotated with at most one scope annotation.

**InjectOnMemberAndConstructor**

Members shouldn't be annotated with @Inject if constructor is already annotated @Inject

**InlineMeValidator**

Ensures that the @InlineMe annotation is used correctly.

**InstantTemporalUnit**

Instant APIs only work for NANOS, MICROS, MILLIS, SECONDS, MINUTES, HOURS, HALF_DAYS and DAYS.

**InvalidJavaTimeConstant**

This checker errors on calls to java.time methods using values that are guaranteed to throw a DateTimeException.

**InvalidPatternSyntax**

Invalid syntax used for a regular expression

**InvalidTimeZoneID**

Invalid time zone identifier. TimeZone.getTimeZone(String) will silently return GMT instead of the time zone you intended.

**InvalidZoneId**

Invalid zone identifier. ZoneId.of(String) will throw exception at runtime.

**IsInstanceIncompatibleType**

This use of isInstance will always evaluate to false.

**IsInstanceOfClass**

The argument to Class#isInstance(Object) should not be a Class

**IsLoggableTagLength**

Log tag too long, cannot exceed 23 characters.

**JUnit3TestNotRun**

Test method will not be run; please correct method signature (Should be public, non-static, and method name should begin with "test").

**JUnit4ClassAnnotationNonStatic**

This method should be static

**JUnit4SetUpNotRun**

setUp() method will not be run; please add JUnit's @Before annotation

**JUnit4TearDownNotRun**

tearDown() method will not be run; please add JUnit's @After annotation

**JUnit4TestNotRun**

This looks like a test method but is not run; please add @Test and @Ignore, or, if this is a helper method, reduce its visibility.

**JUnit4TestsNotRunWithinEnclosed**

This test is annotated @Test, but given it's within a class using the Enclosed runner, will not run.

**JUnitAssertSameCheck**

An object is tested for reference equality to itself using JUnit library.

**JUnitParameterMethodNotFound**

The method for providing parameters was not found.

**JavaxInjectOnAbstractMethod**

Abstract and default methods are not injectable with javax.inject.Inject

**JodaToSelf**

Use of Joda-Time's DateTime.toDateTime(), Duration.toDuration(), Instant.toInstant(), Interval.toInterval(), and Period.toPeriod() are not allowed.

**LabelledBreakTarget**

Labels should only be used on loops.

**LenientFormatStringValidation**

The number of arguments provided to lenient format methods should match the positional specifiers.

**LiteByteStringUtf8**

This pattern will silently corrupt certain byte sequences from the serialized protocol message. Use ByteString or byte[] directly

**LocalDateTemporalAmount**

LocalDate.plus() and minus() does not work with Durations. LocalDate represents civil time (years/months/days), so java.time.Period is the appropriate thing to add or subtract instead.

**LockOnBoxedPrimitive**

It is dangerous to use a boxed primitive as a lock as it can unintentionally lead to sharing a lock with another piece of code.

**LoopConditionChecker**

Loop condition is never modified in loop body.

**LossyPrimitiveCompare**

Using an unnecessarily-wide comparison method can lead to lossy comparison

**MathRoundIntLong**

Math.round(Integer) results in truncation

**MemorySegmentReferenceEquality**

Do not compare MemorySegments using reference equality.

**MislabeledAndroidString**

Certain resources in `android.R.string` have names that do not match their content

**MisleadingEmptyVarargs**

`thenThrow` with no arguments is a no-op, despite reading like it makes the mock throw.

**MisleadingEscapedSpace**

Using \s anywhere except at the end of a line in a text block is potentially misleading.

**MisplacedScopeAnnotations**

Scope annotations used as qualifier annotations don't have any effect. Move the scope annotation to the binding location or delete it.

**MissingSuperCall**

Overriding method is missing a call to overridden super method

**MissingTestCall**

A terminating method call is required for a test helper to have any effect.

**MisusedDayOfYear**

Use of 'DD' (day of year) in a date pattern with 'MM' (month of year) is not likely to be intentional, as it would lead to dates like 'March 73rd'.

**MisusedWeekYear**

Use of "YYYY" (week year) in a date pattern without "ww" (week in year). You probably meant to use "yyyy" (year) instead.

**MixedDescriptors**

The field number passed into #findFieldByNumber belongs to a different proto to the Descriptor.

**MockitoUsage**

Missing method call for verify(mock) here

**ModifyingCollectionWithItself**

Using a collection function with itself as the argument.

**MoreThanOneInjectableConstructor**

This class has more than one @Inject-annotated constructor. Please remove the @Inject annotation from all but one of them.

**MustBeClosedChecker**

This method returns a resource which must be managed carefully, not just left for garbage collection. If it is a constant that will persist for the lifetime of your program, move it to a private static final field. Otherwise, you should use it in a try-with-resources.

**NCopiesOfChar**

The first argument to nCopies is the number of copies, and the second is the item to copy

**NoCanIgnoreReturnValueOnClasses**

@CanIgnoreReturnValue should not be applied to classes as it almost always overmatches (as it applies to constructors and all methods), and the CIRVness isn't conferred to its subclasses.

**NonCanonicalStaticImport**

Static import of type uses non-canonical name

**NonFinalCompileTimeConstant**

@CompileTimeConstant parameters should be final or effectively final

**NonRuntimeAnnotation**

Calling getAnnotation on an annotation that is not retained at runtime

**NullArgumentForNonNullParameter**

Null is not permitted for this parameter.

**NullNeedsCastForVarargs**

This call passes a null *array*, so it always produces NullPointerException. To pass a null *element*, cast to the element type.

**NullTernary**

This conditional expression may evaluate to null, which will result in an NPE when the result is unboxed.

**NullableOnContainingClass**

Type-use nullability annotations should annotate the inner class, not the outer class (e.g., write `A.@Nullable B` instead of `@Nullable A.B`).

**OptionalEquality**

Comparison using reference equality instead of value equality

**OptionalMapUnusedValue**

Optional.ifPresent is preferred over Optional.map when the return value is unused

**OptionalOfRedundantMethod**

Optional.of() always returns a non-empty optional. Using ifPresent/isPresent/orElse/orElseGet/orElseThrow/isPresent/or/orNull method on it is unnecessary and most probably a bug.

**OverlappingQualifierAndScopeAnnotation**

Annotations cannot be both Scope annotations and Qualifier annotations: this causes confusion when trying to use them.

**OverridesJavaxInjectableMethod**

This method is not annotated with @Inject, but it overrides a method that is annotated with @javax.inject.Inject. The method will not be Injected.

**PackageInfo**

Declaring types inside package-info.java files is very bad form

**ParametersButNotParameterized**

This test has @Parameters but is using the default JUnit4 runner. The parameters will have no effect.

**ParcelableCreator**

Detects classes which implement Parcelable but don't have CREATOR

**PeriodFrom**

Period.from(Period) returns itself; from(Duration) throws a runtime exception.

**PeriodGetTemporalUnit**

Period.get() only works with YEARS, MONTHS, or DAYS.

**PeriodTimeMath**

When adding or subtracting from a Period, Duration is incompatible.

**PreconditionsInvalidPlaceholder**

Preconditions only accepts the %s placeholder in error message strings

**PrivateSecurityContractProtoAccess**

Access to a private protocol buffer field is forbidden. This protocol buffer carries a security contract, and can only be created using an approved library. Direct access to the fields is forbidden.

**ProtoBuilderReturnValueIgnored**

Unnecessary call to proto's #build() method. If you don't consume the return value of #build(), the result is discarded and the only effect is to verify that all required fields are set, which can be expressed more directly with #isInitialized().

**ProtoStringFieldReferenceEquality**

Comparing protobuf fields of type String using reference equality

**ProtoTruthMixedDescriptors**

The arguments passed to `ignoringFields` are inconsistent with the proto which is the subject of the assertion.

**ProtocolBufferOrdinal**

To get the tag number of a protocol buffer enum, use getNumber() instead.

**ProvidesMethodOutsideOfModule**

@Provides methods need to be declared in a Module to have any effect.

**RandomCast**

Casting a random number in the range [0.0, 1.0) to an integer or long always results in 0.

**RandomModInteger**

Use Random.nextInt(int). Random.nextInt() % n can have negative results

**RecordAccessorInCompactConstructor**

Record accessors read uninitialized fields inside a compact constructor. Use the component parameter instead.

**RectIntersectReturnValueIgnored**

Return value of android.graphics.Rect.intersect() must be checked

**RedundantSetterCall**

A field was set twice in the same chained expression.

**RequiredModifiers**

This annotation is missing required modifiers as specified by its @RequiredModifiers annotation

**RestrictedApi**

Check for non-allowlisted callers to RestrictedApiChecker.

**ReturnValueIgnored**

Return value of this method must be used

**SelfAssertion**

This assertion will always fail or succeed.

**SelfAssignment**

Variable assigned to itself

**SelfComparison**

An object is compared to itself

**SelfEquals**

Testing an object for equality with itself will always be true.

**SetUnrecognized**

Setting a proto field to an UNRECOGNIZED value will result in an exception at runtime when building.

**ShouldHaveEvenArgs**

This method must be called with an even number of arguments.

**SizeGreaterThanOrEqualsZero**

Comparison of a size >= 0 is always true, did you intend to check for non-emptiness?

**StreamToString**

Calling toString on a Stream does not provide useful information

**StringBuilderInitWithChar**

StringBuilder does not have a char constructor; this invokes the int constructor.

**StringJoin**

String.join(CharSequence) performs no joining (it always returns the empty string); String.join(CharSequence, CharSequence) performs no joining (it just returns the 2nd parameter).

**SubstringOfZero**

String.substring(0) returns the original String

**SuppressWarningsDeprecated**

Suppressing "deprecated" is probably a typo for "deprecation"

**TemporalAccessorGetChronoField**

TemporalAccessor.get() only works for certain values of ChronoField.

**TestParametersNotInitialized**

This test has @TestParameter fields but is using the default JUnit4 runner. The parameters will not be initialised beyond their default value.

**TheoryButNoTheories**

This test has members annotated with @Theory, @DataPoint, or @DataPoints but is using the default JUnit4 runner.

**ThreadBuilderNameWithPlaceholder**

Thread.Builder.name() does not accept placeholders (e.g., %d or %s). threadBuilder.name(String) accepts a constant name and threadBuilder.name(String, int) accepts a constant name *prefix* and an initial counter value.

**ThrowIfUncheckedKnownChecked**

throwIfUnchecked(knownCheckedException) is a no-op.

**ThrowNull**

Throwing 'null' always results in a NullPointerException being thrown.

**TreeToString**

Tree#toString shouldn't be used for Trees deriving from the code being compiled, as it discards whitespace and comments.

**TryFailThrowable**

Catching Throwable/Error masks failures from fail() or assert*() in the try block

**TypeParameterQualifier**

Type parameter used as type qualifier

**UnicodeDirectionalityCharacters**

Unicode directionality modifiers can be used to conceal code in many editors.

**UnicodeInCode**

Avoid using non-ASCII Unicode characters outside of comments and literals, as they can be confusing.

**UnnecessaryCheckNotNull**

This null check is unnecessary; the expression can never be null

**UnnecessaryTypeArgument**

Non-generic methods should not be invoked with type arguments

**UnsafeWildcard**

Certain wildcard types can confuse the compiler.

**UnusedAnonymousClass**

Instance created but never used

**UnusedCollectionModifiedInPlace**

Collection is modified in place, but the result is not used

**VarTypeName**

`var` should not be used as a type name.

**WrongOneof**

This field is guaranteed not to be set given it's within a switch over a one_of.

**XorPower**

The `^` operator is binary XOR, not a power operator.

**ZoneIdOfZ**

Use ZoneOffset.UTC instead of ZoneId.of("Z").

**ASTHelpersSuggestions**

Prefer ASTHelpers instead of calling this API directly

**AddressSelection**

Prefer InetAddress.getAllByName to APIs that convert a hostname to a single IP address

**AlmostJavadoc**

This comment contains Javadoc or HTML tags, but isn't started with a double asterisk (/**); is it meant to be Javadoc?

**AlreadyChecked**

This condition has already been checked.

**AmbiguousMethodReference**

Method reference is ambiguous

**AnnotateFormatMethod**

This method uses a pair of parameters as a format string and its arguments, but the enclosing method wasn't annotated. Doing so gives compile-time rather than run-time protection against malformed format strings.

**ArgumentSelectionDefectChecker**

Arguments are in the wrong order or could be commented for clarity.

**ArrayAsKeyOfSetOrMap**

Arrays do not override equals() or hashCode, so comparisons will be done on reference equality only. If neither deduplication nor lookup are needed, consider using a List instead. Otherwise, use IdentityHashMap/Set, a Map from a library that handles object arrays, or an Iterable/List of pairs.

**ArrayRecordComponent**

Record components should not be arrays.

**AssertEqualsArgumentOrderChecker**

Arguments are swapped in assertEquals-like call

**AssertSameIncompatible**

The types passed to this assertion are incompatible.

**AssertThrowsBlockToExpression**

assertThrows calls with lambdas containing a single statement can be expressed more concisely

**AssertThrowsMinimizer**

Minimize the amount of logic in assertThrows

**AssertThrowsMultipleStatements**

The lambda passed to assertThrows should contain exactly one statement

**AssertionFailureIgnored**

This assertion throws an AssertionError if it fails, which will be caught by an enclosing try block.

**AssignmentExpression**

The use of an assignment expression can be surprising and hard to read; consider factoring out the assignment to a separate statement.

**AssistedInjectAndInjectOnSameConstructor**

@AssistedInject and @Inject cannot be used on the same constructor.

**AttemptedNegativeZero**

-0 is the same as 0. For the floating-point negative zero, use -0.0.

**AutoValueBoxedValues**

AutoValue instances should not usually contain boxed types that are not Nullable. We recommend removing the unnecessary boxing.

**AutoValueFinalMethods**

Make toString(), hashCode() and equals() final in AutoValue classes, so it is clear to readers that AutoValue is not overriding them

**AutoValueImmutableFields**

AutoValue recommends using immutable collections

**AutoValueSubclassLeaked**

Do not refer to the autogenerated AutoValue_ class outside the file containing the corresponding @AutoValue base class.

**AvoidCommonTypeNames**

Never reuse class names from java.lang

**AvoidValueSetter**

Prefer using the enum-accepting rather than the int-accepting setter for enum fields.

**BadComparable**

Possible sign flip from narrowing conversion

**BadImport**

Importing nested classes/static methods/static fields with commonly-used names can make code harder to read, because it may not be clear from the context exactly which type is being referred to. Qualifying the name with that of the containing class can make the code clearer.

**BadInstanceof**

instanceof used in a way that is equivalent to a null check.

**BareDotMetacharacter**

"." is rarely useful as a regex, as it matches any character. To match a literal '.' character, instead write "\.".

**BigDecimalEquals**

BigDecimal#equals has surprising behavior: it also compares scale.

**BigDecimalLiteralDouble**

new BigDecimal(double) loses precision in this case.

**BooleanLiteral**

This expression can be written more clearly with a boolean literal.

**BoxedPrimitiveConstructor**

valueOf or autoboxing provides better time and space performance

**BoxingComparator**

Comparator.comparing() unnecessarily boxes numerical primitives; please use the primitive-specific method instead (e.g., comparingInt()).

**BugPatternNaming**

Giving BugPatterns a name different to the enclosing class can be confusing

**ByteBufferBackingArray**

ByteBuffer.array() shouldn't be called unless ByteBuffer.arrayOffset() is used or if the ByteBuffer was initialized using ByteBuffer.wrap() or ByteBuffer.allocate().

**CacheLoaderNull**

The result of CacheLoader#load must be non-null.

**CanonicalDuration**

Duration can be expressed more clearly with different units

**CatchAndPrintStackTrace**

Logging or rethrowing exceptions should usually be preferred to catching and calling printStackTrace

**CatchFail**

Ignoring exceptions and calling fail() is unnecessary, and makes test output less useful

**ChainedAssertionLosesContext**

Inside a Subject, use check(…) instead of assert*() to preserve user-supplied messages and other settings.

**CharacterGetNumericValue**

getNumericValue has unexpected behaviour: it interprets A-Z as base-36 digits with values 10-35, but also supports non-arabic numerals and miscellaneous numeric unicode characters like ㊷; consider using Character.digit or UCharacter.getUnicodeNumericValue instead

**ClassCanBeStatic**

Inner class is non-static but does not reference enclosing class

**ClassInitializationDeadlock**

Possible class initialization deadlock

**ClassNewInstance**

Class.newInstance() bypasses exception checking; prefer getDeclaredConstructor().newInstance()

**CloseableProvides**

Providing Closeable resources makes their lifecycle unclear

**ClosingStandardOutputStreams**

Don't use try-with-resources to manage standard output streams, closing the stream will cause subsequent output to standard output or standard error to be lost

**CollectionUndefinedEquality**

This type does not have well-defined equals behavior.

**CollectorShouldNotUseState**

Collector.of() should not use state

**ComparableAndComparator**

Class should not implement both `Comparable` and `Comparator`

**CompareToZero**

The result of #compareTo or #compare should only be compared to 0. It is an implementation detail whether a given type returns strictly the values {-1, 0, +1} or others.

**ComplexBooleanConstant**

Non-trivial compile time constant boolean expressions shouldn't be used.

**DateChecker**

Warns against suspect looking calls to java.util.Date APIs

**DateFormatConstant**

DateFormat is not thread-safe, and should not be used as a constant field.

**DeeplyNested**

Very deeply nested code may lead to StackOverflowErrors during compilation

**DefaultCharset**

Implicit use of the platform default charset, which can result in differing behaviour between JVM executions or incorrect behavior if the encoding of the data source doesn't match expectations.

**DefaultPackage**

Java classes shouldn't use default package

**DeprecatedVariable**

Applying the @Deprecated annotation to local variables or parameters has no effect

**DirectInvocationOnMock**

Methods should not be directly invoked on mocks. Should this be part of a verify(..) call?

**DistinctVarargsChecker**

Method expects distinct arguments at some/all positions

**DoNotCallSuggester**

Consider annotating methods that always throw with @DoNotCall. Read more at https://errorprone.info/bugpattern/DoNotCall

**DoNotClaimAnnotations**

Don't 'claim' annotations in annotation processors; Processor#process should unconditionally return `false`

**DoNotMockAutoValue**

AutoValue classes represent pure data classes, so mocking them should not be necessary. Construct a real instance of the class instead.

**DoubleCheckedLocking**

Double-checked locking on non-volatile fields is unsafe

**DuplicateAssertion**

This assertion is duplicate.

**DuplicateBranches**

Both branches contain identical code

**DuplicateDateFormatField**

Reuse of DateFormat fields is most likely unintentional

**EffectivelyPrivate**

This declaration has public or protected modifiers, but is effectively private.

**EmptyBlockTag**

A block tag (@param, @return, @throws, @deprecated) has an empty description. Block tags without descriptions don't add much value for future readers of the code; consider removing the tag entirely or adding a description.

**EmptyCatch**

Caught exceptions should not be ignored

**EmptySetMultibindingContributions**

@Multibinds is a more efficient and declarative mechanism for ensuring that a set multibinding is present in the graph.

**EmptyTopLevelDeclaration**

Empty top-level type declarations should be omitted

**EnumOrdinal**

You should almost never invoke the Enum.ordinal() method or depend on the enum values by index.

**EqualsGetClass**

Prefer instanceof to getClass when implementing Object#equals. Note that this may be a behaviour change.

**EqualsIncompatibleType**

An equality test between objects with incompatible types always returns false

**EqualsUnsafeCast**

The contract of #equals states that it should return false for incompatible types, while this implementation may throw ClassCastException.

**EqualsUsingHashCode**

Implementing #equals by just comparing hashCodes is fragile. Hashes collide frequently, and this will lead to false positives in #equals.

**ErroneousBitwiseExpression**

This expression evaluates to 0. If this isn't an error, consider expressing it as a literal 0.

**ErroneousThreadPoolConstructorChecker**

Thread pool size will never go beyond corePoolSize if an unbounded queue is used

**EscapedEntity**

HTML entities in @code/@literal tags will appear literally in the rendered javadoc.

**ExpensiveLenientFormatString**

String.format is passed to a lenient formatting method, which can be unwrapped to improve efficiency.

**ExposedPrivateType**

Private member classes should not be referenced in signatures of non-private members.

**ExtendingJUnitAssert**

When only using JUnit Assert's static methods, you should import statically instead of extending.

**ExtendsObject**

`T extends Object` is redundant (unless you are using the Checker Framework).

**FallThrough**

Switch case may fall through

**Finalize**

Do not override finalize

**Finally**

If you return or throw from a finally, then values returned or thrown from the try-catch block will be ignored. Consider using try-with-resources instead.

**FloatCast**

Use parentheses to make the precedence explicit

**FloatingPointAssertionWithinEpsilon**

This fuzzy equality check is using a tolerance less than the gap to the next number. You may want a less restrictive tolerance, or to assert equality.

**FloatingPointLiteralPrecision**

Floating point literal loses precision

**FloggerArgumentToString**

Use Flogger's printf-style formatting instead of explicitly converting arguments to strings. Note that Flogger does more than just call toString; for instance, it formats arrays sensibly.

**FloggerPerWithoutRateLimit**

per() methods are no-ops unless combined with atMostEvery(), every(), or onAverageEvery()

**FloggerStringConcatenation**

Prefer string formatting using printf placeholders (e.g. %s) instead of string concatenation

**FormatStringShouldUsePlaceholders**

Using a format string avoids string concatenation in the common case.

**FragmentInjection**

Classes extending PreferenceActivity must implement isValidFragment such that it does not unconditionally return true to prevent vulnerability to fragment injection attacks.

**FragmentNotInstantiable**

Subclasses of Fragment must be instantiable via Class#newInstance(): the class must be public, static and have a public nullary constructor

**FutureReturnValueIgnored**

Return value of methods returning Future must be checked. Ignoring returned Futures suppresses exceptions thrown from the code that completes the Future.

**FutureTransformAsync**

Use transform instead of transformAsync when all returns are an immediate future.

**GetClassOnEnum**

Calling getClass() on an enum may return a subclass of the enum type

**GuiceNestedCombine**

Nesting Modules.combine() here is unnecessary.

**HidingField**

Hiding fields of superclasses may cause confusion and errors

**ICCProfileGetInstance**

This method searches the class path for the given file, prefer to read the file and call getInstance(byte[]) or getInstance(InputStream)

**IdentityHashMapUsage**

IdentityHashMap usage shouldn't be intermingled with Map

**IfChainToSwitch**

This if-chain may be converted into a switch

**IgnoredPureGetter**

Getters on AutoValues, AutoBuilders, and Protobuf Messages are side-effect free, so there is no point in calling them if the return value is ignored. While there are no side effects from the getter, the receiver may have side effects.

**ImmutableAnnotationChecker**

Annotations should always be immutable

**ImmutableEnumChecker**

Enums should always be immutable

**InconsistentCapitalization**

It is confusing to have a field and a parameter under the same scope that differ only in capitalization.

**InconsistentHashCode**

Including fields in hashCode which are not compared in equals violates the contract of hashCode.

**IncorrectMainMethod**

'main' methods must be public, static, and void

**IncrementInForLoopAndHeader**

This for loop increments the same variable in the header and in the body

**InheritDoc**

Invalid use of @inheritDoc.

**InjectInvalidTargetingOnScopingAnnotation**

A scoping annotation's Target should include TYPE and METHOD.

**InjectOnBugCheckers**

BugChecker constructors should be marked @Inject.

**InjectOnConstructorOfAbstractClass**

Constructors on abstract classes are never directly @Inject'ed, only the constructors of their subclasses can be @Inject'ed.

**InjectScopeAnnotationOnInterfaceOrAbstractClass**

Scope annotation on an interface or abstract class is not allowed

**InjectedConstructorAnnotations**

Injected constructors cannot be optional nor have binding annotations

**InlineFormatString**

Prefer to create format strings inline, instead of extracting them to a single-use constant

**InlineMeInliner**

Callers of this API should be inlined.

**InlineMeSuggester**

This deprecated API looks inlineable. If you'd like the body of the API to be automatically inlined to its callers, please annotate it with @InlineMe. NOTE: the suggested fix makes the method final if it was not already.

**InlineTrivialConstant**

Consider inlining this constant

**InputStreamSlowMultibyteRead**

Please also override int read(byte[], int, int), otherwise multi-byte reads from this input stream are likely to be slow.

**InstanceOfAndCastMatchWrongType**

Casting inside an if block should be plausibly consistent with the instanceof type

**IntFloatConversion**

Conversion from int to float may lose precision; use an explicit cast to float if this was intentional

**IntLiteralCast**

Consider using a literal of the desired type instead of casting an int literal

**IntLongMath**

Expression of type int may overflow before being assigned to a long

**InterruptedInCatchBlock**

Did you mean to call Thread.currentThread().interrupt() instead of Thread.interrupted()?

**InvalidBlockTag**

This tag is invalid.

**InvalidInlineTag**

This tag is invalid.

**InvalidLink**

This @link tag looks wrong.

**InvalidParam**

This @param tag doesn't refer to a parameter of the method.

**InvalidSnippet**

This tag is invalid.

**InvalidThrows**

The documented method doesn't actually throw this checked exception.

**InvalidThrowsLink**

Don't use {@link} or {@code} in @throws tags; mention the exception name directly (e.g., @throws IOException, not @throws {@link IOException}).

**IterableAndIterator**

Class should not implement both `Iterable` and `Iterator`

**JUnit3FloatingPointComparisonWithoutDelta**

Floating-point comparison without error tolerance

**JUnit4ClassUsedInJUnit3**

Some JUnit4 construct cannot be used in a JUnit3 context. Convert your class to JUnit4 style to use them.

**JUnit4EmptyMethods**

Empty JUnit4 @Before, @After, @BeforeClass, and @AfterClass methods are unnecessary and should be deleted.

**JUnitAmbiguousTestClass**

Test class inherits from JUnit 3's TestCase but has JUnit 4 @Test or @RunWith annotations.

**JUnitIncompatibleType**

The types passed to this assertion are incompatible.

**JUnitMethodInvoked**

Directly invoking a JUnit test method is discouraged; only the JUnit test runner should call these methods. If you need to share logic between tests, extract a helper method or class.

**JavaDurationGetSecondsGetNano**

duration.getNano() only accesses the underlying nanosecond adjustment from the whole second.

**JavaDurationGetSecondsToToSeconds**

Prefer duration.toSeconds() over duration.getSeconds()

**JavaDurationWithNanos**

Use of java.time.Duration.withNanos(int) is not allowed.

**JavaDurationWithSeconds**

Use of java.time.Duration.withSeconds(long) is not allowed.

**JavaInstantGetSecondsGetNano**

instant.getNano() only accesses the underlying nanosecond adjustment from the whole second.

**JavaLocalDateTimeGetNano**

localDateTime.getNano() only access the nanos-of-second field. It's rare to only use getNano() without a nearby getSecond() call.

**JavaLocalTimeGetNano**

localTime.getNano() only accesses the nanos-of-second field. It's rare to only use getNano() without a nearby getSecond() call.

**JavaPeriodGetDays**

period.getDays() only accesses the "days" portion of the Period, and doesn't represent the total span of time of the period. Consider using org.threeten.extra.Days to extract the difference between two civil dates if you want the whole time.

**JavaTimeDefaultTimeZone**

java.time APIs that silently use the default system time-zone are not allowed.

**JavaUtilDate**

Date has a bad API that leads to bugs; prefer java.time.Instant or LocalDate.

**JavaxInjectOnFinalField**

@javax.inject.Inject cannot be put on a final field.

**JdkObsolete**

Suggests alternatives to obsolete JDK classes.

**JodaConstructors**

Use of certain JodaTime constructors are not allowed.

**JodaDateTimeConstants**

Using the `*PER*` constants in `DateTimeConstants` is problematic because they encourage manual date/time math.

**JodaDurationWithMillis**

Use of duration.withMillis(long) is not allowed. Please use Duration.millis(long) instead.

**JodaInstantWithMillis**

Use of instant.withMillis(long) is not allowed. Use Instant.ofEpochMilli(long) instead.

**JodaNewPeriod**

This may have surprising semantics, e.g. new Period(LocalDate.parse("1970-01-01"), LocalDate.parse("1970-02-02")).getDays() == 1, not 32.

**JodaPlusMinusLong**

Use of JodaTime's type.plus(long) or type.minus(long) is not allowed (where <type> = {Duration,Instant,DateTime,DateMidnight}). Please use type.plus(Duration.millis(long)) or type.minus(Duration.millis(long)) instead.

**JodaTimeConverterManager**

Joda-Time's ConverterManager makes the semantics of DateTime/Instant/etc construction subject to global static state. If you need to define your own converters, use a helper.

**JodaWithDurationAddedLong**

Use of JodaTime's type.withDurationAdded(long, int) (where <type> = {Duration,Instant,DateTime}). Please use type.withDurationAdded(Duration.millis(long), int) instead.

**ListRemoveAmbiguous**

Ambiguous call to List.remove; clarify if index-based or value-based removal was intended by adding a comment

**LiteEnumValueOf**

Instead of looking up a lite enum by name, use its numeric value since that is the stable part of the protocol defined by the enum.

**LiteProtoToString**

toString() on lite protos will not generate a useful representation of the proto from optimized builds. Consider whether using some subset of fields instead would provide useful information.

**LockNotBeforeTry**

Calls to Lock#lock should be immediately followed by a try block which releases the lock.

**LockOnNonEnclosingClassLiteral**

Lock on the class other than the enclosing class of the code block can unintentionally prevent the locked class being used properly.

**LogicalAssignment**

Assignment where a boolean expression was expected; use == if this assignment wasn't expected or add parentheses for clarity.

**LongDoubleConversion**

Conversion from long to double may lose precision; use an explicit cast to double if this was intentional

**LongFloatConversion**

Conversion from long to float may lose precision; use an explicit cast to float if this was intentional

**LoopOverCharArray**

toCharArray allocates a new array, using charAt is more efficient

**LoopToTestParameter**

Migrate loops in tests to use github.com/google/TestParameterInjector. Test parameterization executes each input case in strict isolation, ensuring that a single failure doesn't halt the rest of your test case while providing clear, per-case reporting without the need for manual loops.

**MalformedInlineTag**

This Javadoc tag is malformed. The correct syntax is {@tag and not @{tag.

**MathAbsoluteNegative**

Math.abs() does not always give a non-negative result. Please consider other methods for positive numbers, such as IntMath.saturatedAbs() or Math.floorMod().

**MemoizeConstantVisitorStateLookups**

Anytime you need to look up a constant value from VisitorState, improve performance by creating a cache for it with VisitorState.memoize

**MisformattedTestData**

This test data will be more readable if correctly formatted.

**MissingCasesInEnumSwitch**

Switches on enum types should either handle all values, or have a default case.

**MissingFail**

Not calling fail() when expecting an exception masks bugs

**MissingImplementsComparable**

Classes implementing valid compareTo function should implement Comparable interface

**MissingOverride**

method overrides method in supertype; expected @Override

**MissingRefasterAnnotation**

The Refaster template contains a method without any Refaster annotations

**MissingSummary**

A summary line is required on public/protected Javadocs.

**MixedMutabilityReturnType**

This method returns both mutable and immutable collections or maps from different paths. This may be confusing for users of the method.

**MockIllegalThrows**

This exception can't be thrown by the mocked method.

**MockNotUsedInProduction**

This mock is instantiated and configured, but is never passed to production code. It should be either removed or used.

**ModifiedButNotUsed**

A collection or proto builder was created, but its values were never accessed.

**ModifyCollectionInEnhancedForLoop**

Modifying a collection while iterating over it in a loop may cause a ConcurrentModificationException to be thrown or lead to undefined behavior.

**ModifySourceCollectionInStream**

Modifying the backing source during stream operations may cause unintended results.

**MultimapKeys**

Iterating over `Multimap.keys()` does not collapse duplicates. Did you mean `keySet()`?

**MultipleNullnessAnnotations**

This type use has conflicting nullness annotations

**MultipleParallelOrSequentialCalls**

Multiple calls to either parallel or sequential are unnecessary and cause confusion.

**MultipleUnaryOperatorsInMethodCall**

Avoid having multiple unary operators acting on the same variable in a method call

**MutablePublicArray**

Non-empty arrays are mutable, so this `public static final` array is not a constant and can be modified by clients of this class. Prefer an ImmutableList, or provide an accessor method that returns a defensive copy.

**NamedLikeContextualKeyword**

Avoid naming of classes and methods that is similar to contextual keywords. When invoking such a method, qualify it.

**NarrowCalculation**

This calculation may lose precision compared to its target type.

**NarrowingCompoundAssignment**

Compound assignments may hide dangerous casts

**NegativeCharLiteral**

Casting a negative signed literal to an (unsigned) char might be misleading.

**NestedInstanceOfConditions**

Nested instanceOf conditions of disjoint types create blocks of code that never execute

**NewFileSystem**

Starting in JDK 13, this call is ambiguous with FileSystem.newFileSystem(Path,Map)

**NonApiType**

Certain types should not be passed across API boundaries.

**NonAtomicVolatileUpdate**

This update of a volatile variable is non-atomic

**NonCanonicalType**

This type is referred to by a non-canonical name, which may be misleading.

**NonOverridingEquals**

equals method doesn't override Object.equals

**NotJavadoc**

Avoid using `/**` for comments which aren't actually Javadoc.

**NullOptional**

Passing a literal null to an Optional parameter is almost certainly a mistake. Did you mean to provide an empty Optional?

**NullableConstructor**

Constructors should not be annotated with @Nullable since they cannot return null

**NullableOptional**

Using an Optional variable which is expected to possibly be null is discouraged. It is best to indicate the absence of the value by assigning it an empty optional.

**NullablePrimitive**

Nullness annotations should not be used for primitive types since they cannot be null

**NullablePrimitiveArray**

@Nullable type annotations should not be used for primitive types since they cannot be null

**NullableTypeParameter**

Nullness annotations directly on type parameters are interpreted differently by different tools

**NullableVoid**

void-returning methods should not be annotated with nullness annotations, since they cannot return null

**NullableWildcard**

Nullness annotations directly on wildcard types are interpreted differently by different tools

**ObjectEqualsForPrimitives**

Avoid unnecessary boxing by using plain == for primitive types.

**ObjectToString**

Calling toString on Objects that don't override toString() doesn't provide useful information

**ObjectsHashCodePrimitive**

Objects.hashCode(Object o) should not be passed a primitive value

**OperatorPrecedence**

Use grouping parenthesis to make the operator precedence explicit

**OptionalMapToOptional**

Mapping to another Optional will yield a nested Optional. Did you mean flatMap?

**OptionalNotPresent**

This Optional has been confirmed to be empty at this point, so the call to `get()` or `orElseThrow()` will always throw.

**OrphanedFormatString**

String literal contains format specifiers, but is not passed to a format method

**OutlineNone**

Setting CSS outline style to none or 0 (while not otherwise providing visual focus indicators) is inaccessible for users navigating a web page without a mouse.

**OverrideThrowableToString**

To return a custom message with a Throwable class, one should override getMessage() instead of toString().

**Overrides**

Varargs doesn't agree for overridden method

**OverridesGuiceInjectableMethod**

This method is not annotated with @Inject, but it overrides a method that is annotated with @com.google.inject.Inject. Guice will inject this method, and it is recommended to annotate it explicitly.

**OverridingMethodInconsistentArgumentNamesChecker**

Arguments of overriding method are inconsistent with overridden method.

**ParameterName**

Detects `/* name= */`-style comments on actual parameters where the name doesn't match the formal parameter

**PatternMatchingInstanceof**

This code can be simplified to use a pattern-matching instanceof.

**PreconditionsCheckNotNullRepeated**

Including the first argument of checkNotNull in the failure message is not useful, as it will always be `null`.

**PreferCharsetOverload**

Prefer calling overloads that accept a Charset over those that accept a String encoding name.

**PreferInstanceofOverGetKind**

Prefer instanceof over getKind() checks where possible, as these work well with pattern matching instanceofs

**PreferTestParameter**

When exhaustively testing all values of a single enum or boolean parameter, prefer @TestParameter over @TestParameters.

**PreferThrowsTag**

Prefer the @throws javadoc tag instead of @exception.

**PrimitiveAtomicReference**

Using compareAndSet with boxed primitives is dangerous, as reference rather than value equality is used. Consider using AtomicInteger, AtomicLong, AtomicBoolean from JDK or AtomicDouble from Guava instead.

**ProtectedMembersInFinalClass**

Protected members in final classes can be package-private

**ProtoDurationGetSecondsGetNano**

getNanos() only accesses the underlying nanosecond-adjustment of the duration.

**ProtoTimestampGetSecondsGetNano**

getNanos() only accesses the underlying nanosecond-adjustment of the instant.

**QualifierOrScopeOnInjectMethod**

Qualifiers/Scope annotations on @Inject methods don't have any effect. Move the qualifier annotation to the binding location.

**ReachabilityFenceUsage**

reachabilityFence should always be called inside a finally block

**RecordComponentOverride**

@Override annotations on record components don't do anything.

**RedundantControlFlow**

This continue statement is redundant and can be removed. It may be misleading.

**RefactorSwitch**

This switch can be refactored to be more readable

**ReferenceEquality**

Comparison using reference equality instead of value equality

**RethrowReflectiveOperationExceptionAsLinkageError**

Prefer LinkageError for rethrowing ReflectiveOperationException as unchecked

**ReturnAtTheEndOfVoidFunction**

`return;` is unnecessary at the end of void methods and constructors.

**ReturnFromVoid**

Void methods should not have a @return tag.

**RobolectricShadowDirectlyOn**

Migrate off a deprecated overload of org.robolectric.shadow.api.Shadow#directlyOn

**RuleNotRun**

This TestRule isn't annotated with @Rule, so won't be run.

**RxReturnValueIgnored**

Returned Rx objects must be checked. Ignoring a returned Rx value means it is never scheduled for execution

**SameNameButDifferent**

This type name shadows another in a way that may be confusing.

**ScannerUseDelimiter**

Scanner.useDelimiter is not an efficient way to read an entire InputStream

**SelfAlwaysReturnsThis**

Non-abstract instance methods named 'self()' or 'getThis()' that return the enclosing class must always 'return this'

**SelfSet**

This setter seems to be invoked with a value from its own getter. Is it redundant?

**ShortCircuitBoolean**

Prefer the short-circuiting boolean operators && and || to & and |.

**StatementSwitchToExpressionSwitch**

This statement switch can be converted to a new-style arrow switch

**StaticAssignmentInConstructor**

This assignment is to a static field. Mutating static state from a constructor is highly error-prone.

**StaticAssignmentOfThrowable**

Saving instances of Throwable in static fields is discouraged, prefer to create them on-demand when an exception is thrown

**StaticGuardedByInstance**

Writes to static fields should not be guarded by instance locks

**StaticMockMember**

@Mock members of test classes shouldn't share state between tests and preferably be non-static

**StreamResourceLeak**

Streams that encapsulate a closeable resource should be closed using try-with-resources

**StreamToIterable**

Using stream::iterator creates a one-shot Iterable, which may cause surprising failures.

**StringCaseLocaleUsage**

Specify a `Locale` when calling `String#to{Lower,Upper}Case`. (Note: there are multiple suggested fixes; the third may be most appropriate if you're dealing with ASCII Strings.)

**StringCharset**

Prefer StandardCharsets over using string names for charsets

**StringConcatToTextBlock**

This string literal can be written more clearly as a text block

**StringSplitter**

String.split(String) has surprising behavior

**SuperCallToObjectMethod**

`super.equals(obj)` and `super.hashCode()` are often bugs when they call the methods defined in `java.lang.Object`

**SwigMemoryLeak**

SWIG generated code that can't call a C++ destructor will leak memory

**SynchronizeOnNonFinalField**

Synchronizing on non-final fields is not safe: if the field is ever updated, different threads may end up locking on different objects.

**SystemConsoleNull**

System.console() no longer returns null in JDK 22 and newer versions

**ThreadJoinLoop**

Thread.join needs to be immediately surrounded by a loop until it succeeds. Consider using Uninterruptibles.joinUninterruptibly.

**ThreadLocalUsage**

ThreadLocals should be stored in static fields

**ThreadPriorityCheck**

Relying on the thread scheduler is discouraged.

**ThreeLetterTimeZoneID**

Three-letter time zone identifiers are deprecated, may be ambiguous, and might not do what you intend; the full IANA time zone ID should be used instead.

**ThrowIfUncheckedKnownUnchecked**

`throwIfUnchecked(knownUnchecked)` is equivalent to `throw knownUnchecked`.

**ThrowableEqualsHashCode**

Overriding Throwable.equals() or hashCode() is discouraged.

**TimeInStaticInitializer**

Accessing the current time in a static initialiser captures the time at class loading, which is rarely desirable.

**TimeUnitConversionChecker**

This TimeUnit conversion looks buggy: converting from a smaller unit to a larger unit (and passing a constant), converting to/from the same TimeUnit, or converting TimeUnits where the result is statically known to be 0 or 1 are all buggy patterns.

**ToStringReturnsNull**

An implementation of Object.toString() should never return null.

**TraditionalSwitchExpression**

Prefer -> switches for switch expressions

**TruthAssertExpected**

The actual and expected values appear to be swapped, which results in poor assertion failure messages. The actual value should come first.

**TruthConstantAsserts**

Truth Library assert is called on a constant.

**TruthGetOrDefault**

Asserting on getOrDefault is unclear; prefer containsEntry or doesNotContainKey

**TruthIncompatibleType**

Argument is not compatible with the subject's type.

**TypeEquals**

TypeMirror should be compared using Types#isSameType, not equality operators or equals().

**TypeNameShadowing**

Type parameter declaration shadows another named type

**TypeParameterShadowing**

Type parameter declaration overrides another type parameter already declared

**TypeParameterUnusedInFormals**

Declaring a type parameter that is only used in the return type is a misuse of generics: operations on the type parameter are unchecked, it hides unsafe casts at invocations of the method, and it interacts badly with method overload resolution.

**URLEqualsHashCode**

Avoid hash-based containers of java.net.URL–the containers rely on equals() and hashCode(), which cause java.net.URL to make blocking internet connections.

**UndefinedEquals**

This type is not guaranteed to implement a useful equals() method.

**UnicodeEscape**

Using unicode escape sequences for printable ASCII characters is obfuscated, and potentially dangerous.

**UnnamedVariable**

Consider renaming unused variables and lambda parameters to `_`

**UnnecessaryAssignment**

Fields annotated with @Inject/@Mock/@TestParameter should not be manually assigned to, as they should be initialized by a framework. Remove the assignment if a framework is being used, or the annotation if one isn't.

**UnnecessaryAsync**

Variables which are initialized and do not escape the current scope do not need to worry about concurrency. Using the non-concurrent type will reduce overhead and verbosity.

**UnnecessaryBreakInSwitch**

This break is unnecessary, fallthrough does not occur in -> switches

**UnnecessaryCopy**

This collection is already immutable (just not ImmutableList/ImmutableMap); copying it is unnecessary.

**UnnecessaryLambda**

Returning a lambda from a helper method or saving it in a constant is unnecessary; prefer to implement the functional interface method directly and use a method reference instead.

**UnnecessaryLongToIntConversion**

Converting a long or Long to an int to pass as a long parameter is usually not necessary. If this conversion is intentional, consider `Longs.constrainToRange()` instead.

**UnnecessaryMethodInvocationMatcher**

It is not necessary to wrap a MethodMatcher with methodInvocation().

**UnnecessaryMethodReference**

This method reference is unnecessary, and can be replaced with the variable itself.

**UnnecessaryParentheses**

These parentheses are unnecessary; it is unlikely the code will be misinterpreted without them

**UnnecessaryQualifier**

A qualifier annotation has no effect here.

**UnnecessaryStringBuilder**

Prefer string concatenation over explicitly using `StringBuilder#append`, since `+` reads better and has equivalent or better performance.

**UnrecognisedJavadocTag**

This Javadoc tag wasn't recognised by the parser. Is it malformed somehow, perhaps with mismatched braces?

**UnsafeFinalization**

Finalizer may run before native code finishes execution

**UnsafeReflectiveConstructionCast**

Prefer `asSubclass` instead of casting the result of `newInstance`, to detect classes of incorrect type before invoking their constructors. This way, if the class is of the incorrect type, it will throw an exception before invoking its constructor.

**UnsynchronizedOverridesSynchronized**

Unsynchronized method overrides a synchronized method.

**UnusedLabel**

This label is unused.

**UnusedMethod**

Unused.

**UnusedNestedClass**

This nested class is unused, and can be removed.

**UnusedTypeParameter**

This type parameter is unused and can be removed.

**UnusedVariable**

Unused.

**UseBinds**

@Binds is a more efficient and declarative mechanism for delegating a binding.

**VariableNameSameAsType**

variableName and type with the same name would refer to the static field instead of the class

**VoidUsed**

Using a Void-typed variable is potentially confusing, and can be replaced with a literal `null`.

**WaitNotInLoop**

Because of spurious wakeups, Object.wait() and Condition.await() must always be called in a loop

**WakelockReleasedDangerously**

On Android versions < P, a wakelock acquired with a timeout may be released by the system before calling `release`, even after checking `isHeld()`. If so, it will throw a RuntimeException. Please wrap in a try/catch block.

**AutoFactoryAtInject**

@AutoFactory and @Inject should not be used in the same type.

**BanClassLoader**

Using dangerous ClassLoader APIs may deserialize untrusted user input into bytecode, leading to remote code execution vulnerabilities

**BanSerializableRead**

Deserializing user input via the `Serializable` API is extremely dangerous

**ClassName**

The source file name should match the name of the top-level class it contains

**ComparisonContractViolated**

This comparison method violates the contract

**DeduplicateConstants**

This expression was previously declared as a constant; consider replacing this occurrence.

**DepAnn**

Item documented with a @deprecated javadoc note is not annotated with @Deprecated

**EmptyIf**

Empty statement after if

**ExtendsAutoValue**

Do not extend an @AutoValue-like classes in non-generated code.

**InsecureCryptoUsage**

A standard cryptographic operation is used in a mode that is prone to vulnerabilities

**IterablePathParameter**

Path implements Iterable<Path>; prefer Collection<Path> for clarity

**Java8ApiChecker**

Use of class, field, or method that is not compatible with JDK 8

**LongLiteralLowerCaseSuffix**

Prefer 'L' to 'l' for the suffix to long literals

**MissingRuntimeRetention**

Scoping and qualifier annotations must have runtime retention.

**MoreThanOneQualifier**

Using more than one qualifier annotation on the same element is not allowed.

**NoAllocation**

@NoAllocation was specified on this method, but something was found that would trigger an allocation

**RefersToDaggerCodegen**

Don't refer to Dagger's internal or generated code

**StaticOrDefaultInterfaceMethod**

Static and default interface methods are not natively supported on older Android devices.

**StaticQualifiedUsingExpression**

A static variable or method should be qualified with a class name, not expression

**SystemExitOutsideMain**

Code that contains System.exit() is untestable.

**TestExceptionChecker**

Using @Test(expected=…) is discouraged, since the test will pass if *any* statement in the test method throws the expected exception

**ThreadSafe**

Type declaration annotated with @ThreadSafe is not thread safe

**UseCorrectAssertInTests**

Java assert is used in testing code. For testing purposes, prefer using Truth-based assertions.

**AnnotationPosition**

Annotations should be positioned after Javadocs, but before modifiers.

**AssertFalse**

Assertions may be disabled at runtime and do not guarantee that execution will halt here; consider throwing an exception instead

**AssistedInjectAndInjectOnConstructors**

@AssistedInject and @Inject should not be used on different constructors in the same class.

**AvoidObjectArrays**

Object arrays are inferior to collections in almost every way. Prefer immutable collections (e.g., ImmutableSet, ImmutableList, etc.) over an object array whenever possible.

**BinderIdentityRestoredDangerously**

A call to Binder.clearCallingIdentity() should be followed by Binder.restoreCallingIdentity() in a finally block. Otherwise the wrong Binder identity may be used by subsequent code.

**BindingToUnqualifiedCommonType**

This code declares a binding for a common value type without a Qualifier annotation.

**BuilderReturnThis**

Builder instance method does not return 'this'

**CanIgnoreReturnValueSuggester**

Methods that always return 'this' (or return an input parameter) should be annotated with @com.google.errorprone.annotations.CanIgnoreReturnValue

**CannotMockFinalClass**

Mockito cannot mock final classes

**CannotMockMethod**

Mockito cannot mock final or static methods, and can't detect this at runtime

**CatchingUnchecked**

This catch block catches `Exception`, but can only catch unchecked exceptions. Consider catching RuntimeException (or something more specific) instead so it is more apparent that no checked exceptions are being handled.

**CheckedExceptionNotThrown**

This method cannot throw a checked exception that it claims to. This may cause consumers of the API to incorrectly attempt to handle, or propagate, this exception.

**ConstantPatternCompile**

Variables initialized with Pattern#compile calls on constants can be constants

**DefaultLocale**

Implicit use of the JVM default locale, which can result in differing behaviour between JVM executions.

**DifferentNameButSame**

This type is referred to in different ways within this file, which may be confusing.

**EqualsBrokenForNull**

equals() implementation may throw NullPointerException when given null

**ExpectedExceptionChecker**

Prefer assertThrows to ExpectedException

**ExplicitArrayForVarargs**

Avoid explicit array creation for varargs

**FloggerLogWithCause**

Setting the caught exception as the cause of the log message may provide more context for anyone debugging errors.

**FloggerMessageFormat**

Invalid message format-style format specifier ({0}), expected printf-style (%s)

**FloggerRedundantIsEnabled**

Logger level check is already implied in the log() call. An explicit atLEVEL().isEnabled() check is redundant.

**FloggerRequiredModifiers**

FluentLogger.forEnclosingClass should always be saved to a private static final field.

**FloggerWithCause**

Calling withCause(Throwable) with an inline allocated Throwable is discouraged. Consider using withStackTrace(StackSize) instead, and specifying a reduced stack size (e.g. SMALL, MEDIUM or LARGE) instead of FULL, to improve performance.

**FloggerWithoutCause**

Use withCause to associate Exceptions with log statements

**FunctionalInterfaceClash**

Overloads will be ambiguous when passing lambda arguments.

**HardCodedSdCardPath**

Hardcoded reference to /sdcard

**IdentifierName**

Methods and non-static variables should be named in lowerCamelCase

**InconsistentOverloads**

The ordering of parameters in overloaded methods should be as consistent as possible (when viewed from left to right)

**InitializeInline**

Initializing variables in their declaring statement is clearer, where possible.

**InterfaceWithOnlyStatics**

This interface only contains static fields and methods; consider making it a final class instead to prevent subclassing.

**InterruptedExceptionSwallowed**

This catch block appears to be catching an explicitly declared InterruptedException as an Exception/Throwable and not handling the interruption separately.

**Interruption**

Always pass 'false' to 'Future.cancel()', unless you are propagating a cancellation-with-interrupt from another caller

**MissingDefault**

The Google Java Style Guide requires that each switch statement includes a default statement group (even if it contains no code) unless the switch statement covers all values of an enum.

**MissingJavadoc**

Public types must have Javadoc comments.

**MockitoDoSetup**

Prefer using when/thenReturn over doReturn/when for additional type safety.

**MutableGuiceModule**

Fields in Guice modules should be final

**NegativeBoolean**

Prefer positive boolean names

**NonCanonicalStaticMemberImport**

Static import of member uses non-canonical name

**NonFinalStaticField**

Static fields should almost always be final.

**PreferJavaTimeOverload**

Prefer using java.time-based APIs when available. Note that this checker does not and cannot guarantee that the overloads have equivalent semantics, but that is generally the case with overloaded methods.

**PreferPreconditions**

Consider using Preconditions instead of explicit if-throw for parameter validation.

**PreferredInterfaceType**

This type can be more specific.

**PrimitiveArrayPassedToVarargsMethod**

Passing a primitive array to a varargs method is usually wrong

**QualifierWithTypeUse**

Injection frameworks currently don't understand Qualifiers in TYPE_PARAMETER or TYPE_USE contexts.

**RecordComponentAccessorAnnotationConflict**

Annotation on record component is ignored.

**RedundantNullCheck**

Null check on an expression that is statically determined to be non-null according to language semantics or nullness annotations.

**RedundantOverride**

This overriding method is redundant, and can be removed.

**RedundantThrows**

Thrown exception is a subtype of another

**StringFormatWithLiteral**

There is no need to use String.format() when all the arguments are literals.

**StronglyTypeByteString**

This primitive byte array is only used to construct ByteStrings. It would be clearer to strongly type the field instead.

**StronglyTypeTime**

This primitive integral type is only used to construct time types. It would be clearer to strongly type the field instead.

**SunApi**

Usage of internal proprietary API which may be removed in a future release

**SuppressWarningsWithoutExplanation**

Use of @SuppressWarnings should be accompanied by a comment describing why the warning is safe to ignore.

**SystemOut**

Production code should not print to standard out or standard error. Standard out and standard error should only be used for debugging.

**ThrowSpecificExceptions**

Base exception classes should be treated as abstract. If the exception is intended to be caught, throw a domain-specific exception. Otherwise, prefer a more specific exception for clarity. Common alternatives include: AssertionError, IllegalArgumentException, IllegalStateException, and (Guava's) VerifyException.

**TimeUnitMismatch**

An value that appears to be represented in one unit is used where another appears to be required (e.g., seconds where nanos are needed)

**TooManyParameters**

A large number of parameters on public APIs should be avoided.

**TransientMisuse**

Static fields are implicitly transient, so the explicit modifier is unnecessary

**TruthContainsExactlyElementsInUsage**

containsExactly is preferred over containsExactlyElementsIn when creating new iterables.

**TryWithResourcesVariable**

This variable is unnecessary, the try-with-resources resource can be a reference to a final or effectively final variable

**UnescapedEntity**

Javadoc is interpreted as HTML, so HTML entities such as &, <, > must be escaped. If this finding seems wrong (e.g. is within a @code or @literal tag), check whether the tag could be malformed and not recognised by the compiler.

**UnnecessarilyFullyQualified**

This fully qualified name is unambiguous to the compiler if imported.

**UnnecessarilyUsedValue**

The result of this API is ignorable, so it does not need to be captured / assigned into an `unused` variable.

**UnnecessarilyVisible**

Some methods (such as those annotated with @Inject or @Provides) are only intended to be called by a framework, and so should have default visibility.

**UnnecessaryAnonymousClass**

Implementing a functional interface is unnecessary; prefer to implement the functional interface method directly and use a method reference instead.

**UnnecessaryDefaultInEnumSwitch**

Switch handles all enum values: an explicit default case is unnecessary and defeats error checking for non-exhaustive switches.

**UnnecessaryFinal**

Since Java 8, it's been unnecessary to make local variables and parameters `final` for use in lambdas or anonymous classes. Marking them as `final` is weakly discouraged, as it adds a fair amount of noise for minimal benefit.

**UnnecessaryOptionalGet**

This code can be simplified by directly using the lambda parameters instead of calling get..() on optional.

**UnnecessarySemicolon**

Unnecessary semicolons should be omitted. For empty block statements, prefer {}.

**UnnecessaryTestMethodPrefix**

A `test` prefix for a JUnit4 test is redundant, and a holdover from JUnit3. The `@Test` annotation makes it clear it's a test.

**UnsafeLocaleUsage**

Possible unsafe operation related to the java.util.Locale class.

**UnusedException**

This catch block catches an exception and re-throws another, but swallows the caught exception rather than setting it as a cause. This can make debugging harder.

**UrlInSee**

URLs should not be used in @see tags; they are designed for Java elements which could be used with @link.

**UsingJsr305CheckReturnValue**

Prefer ErrorProne's @CheckReturnValue over JSR305's version.

**Var**

Non-constant variable missing @Var annotation

**Varifier**

Consider using `var` here to avoid boilerplate.

**YodaCondition**

The non-constant portion of a comparison generally comes first. For equality, prefer e.equals(CONSTANT) if e is non-null or Objects.equals(e, CONSTANT) if e may be null. For standard operators, prefer e <OPERATION> CONSTANT.

**AddNullMarkedToClass**

Apply @NullMarked to this class

**AddNullMarkedToPackageInfo**

Apply @NullMarked to this package

**AnnotationMirrorToString**

AnnotationMirror#toString doesn't use fully qualified type names, prefer auto-common's AnnotationMirrors#toString

**AnnotationValueToString**

AnnotationValue#toString doesn't use fully qualified type names, prefer auto-common's AnnotationValues#toString

**BooleanParameter**

Use parameter comments to document ambiguous literals

**ClassNamedLikeTypeParameter**

This class's name looks like a Type Parameter.

**ConstantField**

Fields with CONSTANT_CASE names should be both static and final

**EqualsMissingNullable**

Method overrides Object.equals but does not have @Nullable on its parameter

**FieldCanBeFinal**

This field is only assigned during initialization; consider making it final

**FieldCanBeLocal**

This field can be replaced with a local variable in the methods that use it.

**FieldCanBeStatic**

A final field initialized at compile-time with an instance of an immutable type can be static.

**FieldMissingNullable**

Field is assigned (or compared against) a definitely null value but is not annotated @Nullable

**ForEachIterable**

This loop can be replaced with an enhanced for loop.

**ImmutableMemberCollection**

If you don't intend to mutate a member collection prefer using Immutable types.

**ImmutableRefactoring**

Refactors uses of the JSR 305 @Immutable to Error Prone's annotation

**ImmutableSetForContains**

This private static ImmutableList is only used for contains, containsAll or isEmpty checks; prefer ImmutableSet.

**ImplementAssertionWithChaining**

Prefer check(…), which usually generates more readable failure messages.

**LambdaFunctionalInterface**

Use Java's utility functional interfaces instead of Function<A, B> for primitive types.

**MethodCanBeStatic**

This method does not reference the enclosing instance and can be static

**MissingBraces**

The Google Java Style Guide requires braces to be used with if, else, for, do and while statements, even when the body is empty or contains only a single statement.

**MixedArrayDimensions**

C-style array declarations should not be used

**MultiVariableDeclaration**

Variable declarations should declare only one variable

**MultipleTopLevelClasses**

Source files should not contain multiple top-level class declarations

**PackageLocation**

Package names should match the directory they are declared in

**ParameterComment**

Non-standard parameter comment; prefer `/* paramName= */ arg`

**ParameterMissingNullable**

Parameter has handling for null but is not annotated @Nullable

**PrivateConstructorForNoninstantiableModule**

Add a private constructor to modules that will not be instantiated by Dagger.

**PrivateConstructorForUtilityClass**

Classes which are not intended to be instantiated should be made non-instantiable with a private constructor. This includes utility classes (classes with only static members), and the main class.

**PublicApiNamedStreamShouldReturnStream**

Public methods named stream() are generally expected to return a type whose name ends with Stream. Consider choosing a different method name instead.

**RemoveUnusedImports**

Unused imports

**RequireNonNullRefactoring**

Refactor explicit null checks to Objects.requireNonNull

**ReturnMissingNullable**

Method returns a definitely null value but is not annotated @Nullable

**ReturnsNullCollection**

Method has a collection return type and returns {@code null} in some cases but does not annotate the method as @Nullable. See Effective Java 3rd Edition Item 54.

**ScopeOnModule**

Scopes on modules have no function and will soon be an error.

**SwitchDefault**

The default case of a switch should appear at the end of the last statement group

**SymbolToString**

Element#toString shouldn't be used for comparison as it is expensive and fragile.

**ThrowsUncheckedException**

Unchecked exceptions do not need to be declared in the method signature.

**TraditionalJavadocToMarkdown**

Converts traditional Javadoc comments into Markdown Javadoc comments

**TryFailRefactoring**

Prefer assertThrows to try/fail

**TypeParameterNaming**

Type parameters must be a single letter with an optional numeric suffix, or an UpperCamelCase name followed by the letter 'T'.

**TypeToString**

TypeMirror#toString shouldn't be used for comparison as it is expensive and fragile.

**UngroupedOverloads**

Constructors and methods with the same name should appear sequentially with no other code in between, even when modifiers such as static or private differ between the methods. Please re-order or re-name methods.

**UnnecessaryBoxedAssignment**

This expression can be implicitly boxed.

**UnnecessaryBoxedVariable**

It is unnecessary for this variable to be boxed. Use the primitive instead.

**UnnecessarySetDefault**

Unnecessary call to NullPointerTester#setDefault

**UnnecessaryStaticImport**

Using static imports for types is unnecessary

**UseEnumSwitch**

Prefer using a switch instead of a chained if-else for enums

**VarWithPrimitive**

Avoid using `var` with primitive types. Explicit primitive type names are short and clear, and `var` provides no benefit in readability while potentially hiding the type.

**VoidMissingNullable**

The type Void is not annotated @Nullable

**WildcardImport**

Wildcard imports, static or otherwise, should not be used

Our goal is to make it simple to add Error Prone checks to your existing Java
compilation. Please note that Error Prone must be run on JDK 21 or newer. (It
can still be used to build Java 8 code by setting the appropriate `-source` /
`-target` / `-bootclasspath` or `--release` flags.) If you still have to build
with an older JDK version you can use
older unmaintained Error Prone versions.

Please join our mailing list to know when a new version is released!

Error Prone works out of the box with Bazel.

```
java_library(
 name = "hello",
 srcs = ["Hello.java"],
)
```
```
$ bazel build :hello
ERROR: example/myproject/BUILD:29:1: Java compilation in rule '//example/myproject:hello'
examples/maven/error_prone_should_flag/src/main/java/Main.java:20: error: [DeadException] Exception created but not thrown
 new Exception();
 ^
 (see https://errorprone.info/bugpattern/DeadException)
 Did you mean 'throw new Exception();'?
1 error
BazelJavaBuilder threw exception: java compilation returned status ERROR
INFO: Elapsed time: 1.989s, Critical Path: 1.69s
```
Edit your `pom.xml` file to add settings to the maven-compiler-plugin:

```
 <build>
 <plugins>
 <plugin>
 <groupId>org.apache.maven.plugins</groupId>
 <artifactId>maven-compiler-plugin</artifactId>
 <version>3.11.0</version>
 <configuration>
 <source>17</source>
 <target>8</target>
 <encoding>UTF-8</encoding>
 <compilerArgs>
 <arg>-XDcompilePolicy=simple</arg>
 <arg>--should-stop=ifError=FLOW</arg>
 <arg>-Xplugin:ErrorProne</arg>
 <!-- necessary only if building with JDK21: https://github.com/google/error-prone/issues/5426 -->
 <arg>-XDaddTypeAnnotationsToSymbol=true</arg>
 </compilerArgs>
 <annotationProcessorPaths>
 <path>
 <groupId>com.google.errorprone</groupId>
 <artifactId>error_prone_core</artifactId>
 <version>${error-prone.version}</version>
 </path>
 <!-- Other annotation processors go here.
 If 'annotationProcessorPaths' is set, processors will no longer be
 discovered on the regular -classpath; see also 'Using Error Prone
 together with other annotation processors' below. -->
 </annotationProcessorPaths>
 </configuration>
 </plugin>
 </plugins>
 </build>
```
Additional flags are required due to JEP 396: Strongly Encapsulate JDK Internals by Default.

If your `maven-compiler-plugin` uses an external executable, e.g. because
`<fork>` is `true` or because the `maven-toolchains-plugin` is enabled, add the
following under `compilerArgs` in the configuration above:

```
 <arg>-J--add-exports=jdk.compiler/com.sun.tools.javac.api=ALL-UNNAMED</arg>
 <arg>-J--add-exports=jdk.compiler/com.sun.tools.javac.file=ALL-UNNAMED</arg>
 <arg>-J--add-exports=jdk.compiler/com.sun.tools.javac.main=ALL-UNNAMED</arg>
 <arg>-J--add-exports=jdk.compiler/com.sun.tools.javac.model=ALL-UNNAMED</arg>
 <arg>-J--add-exports=jdk.compiler/com.sun.tools.javac.parser=ALL-UNNAMED</arg>
 <arg>-J--add-exports=jdk.compiler/com.sun.tools.javac.processing=ALL-UNNAMED</arg>
 <arg>-J--add-exports=jdk.compiler/com.sun.tools.javac.tree=ALL-UNNAMED</arg>
 <arg>-J--add-exports=jdk.compiler/com.sun.tools.javac.util=ALL-UNNAMED</arg>
 <arg>-J--add-opens=jdk.compiler/com.sun.tools.javac.code=ALL-UNNAMED</arg>
 <arg>-J--add-opens=jdk.compiler/com.sun.tools.javac.comp=ALL-UNNAMED</arg>
```
Otherwise, add the following to the .mvn/jvm.config file:

```
--add-exports jdk.compiler/com.sun.tools.javac.api=ALL-UNNAMED
--add-exports jdk.compiler/com.sun.tools.javac.file=ALL-UNNAMED
--add-exports jdk.compiler/com.sun.tools.javac.main=ALL-UNNAMED
--add-exports jdk.compiler/com.sun.tools.javac.model=ALL-UNNAMED
--add-exports jdk.compiler/com.sun.tools.javac.parser=ALL-UNNAMED
--add-exports jdk.compiler/com.sun.tools.javac.processing=ALL-UNNAMED
--add-exports jdk.compiler/com.sun.tools.javac.tree=ALL-UNNAMED
--add-exports jdk.compiler/com.sun.tools.javac.util=ALL-UNNAMED
--add-opens jdk.compiler/com.sun.tools.javac.code=ALL-UNNAMED
--add-opens jdk.compiler/com.sun.tools.javac.comp=ALL-UNNAMED
```
See the flags documentation for details on how to customize the plugin’s behavior.

The gradle plugin is an external contribution. The documentation and code is at tbroyer/gradle-errorprone-plugin.

Download the following artifacts from maven:

`error_prone_core-${EP_VERSION?}-with-dependencies.jar` from
https://repo1.maven.org/maven2/com/google/errorprone/error_prone_core/`dataflow-errorprone-${DATAFLOW_VERSION?}.jar` from
https://repo1.maven.org/maven2/io/github/eisop/dataflow-errorprone/and add the following javac task to your project’s `build.xml` file:

```
 <path id="processorpath.ref">
 <pathelement location="${user.home}/.m2/repository/com/google/errorprone/error_prone_core/${error-prone.version}/error_prone_core-${error-prone.version}-with-dependencies.jar"/>
 <pathelement location="${user.home}/.m2/repository/io/github/eisop/dataflow-errorprone/${dataflow.version}/dataflow-errorprone-${dataflow.version}.jar"/>
 </path>
 <javac srcdir="src" destdir="build" fork="yes" includeantruntime="no">
 <compilerarg value="-XDcompilePolicy=simple"/>
 <compilerarg value="--should-stop=ifError=FLOW"/>
 <compilerarg value="-processorpath"/>
 <compilerarg pathref="processorpath.ref"/>
 <compilerarg value="-Xplugin:ErrorProne -Xep:DeadException:ERROR" />
 <compilerarg value="-J--add-exports=jdk.compiler/com.sun.tools.javac.api=ALL-UNNAMED" />
 <compilerarg value="-J--add-exports=jdk.compiler/com.sun.tools.javac.file=ALL-UNNAMED" />
 <compilerarg value="-J--add-exports=jdk.compiler/com.sun.tools.javac.main=ALL-UNNAMED" />
 <compilerarg value="-J--add-exports=jdk.compiler/com.sun.tools.javac.model=ALL-UNNAMED" />
 <compilerarg value="-J--add-exports=jdk.compiler/com.sun.tools.javac.parser=ALL-UNNAMED" />
 <compilerarg value="-J--add-exports=jdk.compiler/com.sun.tools.javac.processing=ALL-UNNAMED" />
 <compilerarg value="-J--add-exports=jdk.compiler/com.sun.tools.javac.tree=ALL-UNNAMED" />
 <compilerarg value="-J--add-exports=jdk.compiler/com.sun.tools.javac.util=ALL-UNNAMED" />
 <compilerarg value="-J--add-opens=jdk.compiler/com.sun.tools.javac.code=ALL-UNNAMED" />
 <compilerarg value="-J--add-opens=jdk.compiler/com.sun.tools.javac.comp=ALL-UNNAMED" />
 </javac>
```
Setting the following `--add-exports=` flags is required on JDK 16 and newer due
to
JEP 396: Strongly Encapsulate JDK Internals by Default:

To add the plugin, start the IDE and find the Plugins dialog. Browse Repositories, choose Category: Build, and find the Error-prone plugin. Right-click and choose “Download and install”. The IDE will restart after you’ve exited these dialogs.

To enable Error Prone, choose ```
Settings | Compiler | Java Compiler | Use
compiler: Javac with error-prone
```
 and also make sure ```
Settings | Compiler | Use
external build
```
 is NOT selected.

Ideally, you should find out about failed Error Prone checks as you code in
eclipse, thanks to the continuous compilation by ECJ (eclipse compiler for
Java). But this is an architectural challenge, as Error Prone currently relies
heavily on the `com.sun.*` APIs for accessing the AST and symbol table.

For now, Eclipse users should use the Findbugs eclipse plugin instead, as it catches many of the same issues.

Error Prone supports the
`com.sun.source.util.Plugin`
API, which can be used by adding Error Prone to the `-processorpath` and setting
the `-Xplugin` flag.

Example:

```
wget https://repo1.maven.org/maven2/com/google/errorprone/error_prone_core/${EP_VERSION?}/error_prone_core-${EP_VERSION?}-with-dependencies.jar
wget https://repo1.maven.org/maven2/io/github/eisop/dataflow-errorprone/${DATAFLOW_VERSION?}/dataflow-errorprone-${DATAFLOW_VERSION?}.jar
javac \
 -J--add-exports=jdk.compiler/com.sun.tools.javac.api=ALL-UNNAMED \
 -J--add-exports=jdk.compiler/com.sun.tools.javac.file=ALL-UNNAMED \
 -J--add-exports=jdk.compiler/com.sun.tools.javac.main=ALL-UNNAMED \
 -J--add-exports=jdk.compiler/com.sun.tools.javac.model=ALL-UNNAMED \
 -J--add-exports=jdk.compiler/com.sun.tools.javac.parser=ALL-UNNAMED \
 -J--add-exports=jdk.compiler/com.sun.tools.javac.processing=ALL-UNNAMED \
 -J--add-exports=jdk.compiler/com.sun.tools.javac.tree=ALL-UNNAMED \
 -J--add-exports=jdk.compiler/com.sun.tools.javac.util=ALL-UNNAMED \
 -J--add-opens=jdk.compiler/com.sun.tools.javac.code=ALL-UNNAMED \
 -J--add-opens=jdk.compiler/com.sun.tools.javac.comp=ALL-UNNAMED \
 -XDcompilePolicy=simple \
 --should-stop=ifError=FLOW \
 -processorpath error_prone_core-${EP_VERSION?}-with-dependencies.jar:dataflow-errorprone-${DATAFLOW_VERSION?}.jar \
 '-Xplugin:ErrorProne -XepDisableAllChecks -Xep:CollectionIncompatibleType:ERROR' \
 ShortSet.java
```
```
ShortSet.java:8: error: [CollectionIncompatibleType] Argument 'i - 1' should not be passed to this method; its type int is not compatible with its collection's type argument Short
 s.remove(i - 1);
 ^
 (see https://errorprone.info/bugpattern/CollectionIncompatibleType)
1 error
```
The `--add-exports` and `--add-opens` flags are required due to
JEP 396: Strongly Encapsulate JDK Internals by Default:

If you’re an end-user of the build system, you can file a bug to request integration.

If you develop a build system, you should create an integration for your users! Here are some basics to get you started:

Error-prone is implemented as a compiler hook, using an internal mechanism in
javac. To install our hook, we override the `main()` method in
`com.sun.tools.javac.main.Main`.

Find the spot in your build system where javac’s main method is called. This is assuming you call javac in-process, rather than shell’ing out to the javac executable on the machine (which would be pretty lame since it’s hard to know where that’s located).

First, add Error Prone’s core library to the right classpath. It will need to be
visible to the classloader which currently locates the javac Main class. Then
replace the call of `javac.main.Main.main()` with the Error Prone compiler:

`return new ErrorProneCompiler.Builder().build().compile(args) == 0`

All of the above instructions use the `javac` option `-processorpath` which has
side-effect of causing `javac` to no longer scan the compile classpath for
annotation processors. If you are using other annotation processors in addition
to Error Prone, such as
AutoValue, then you will
need to add their JARs to your `-processorpath` argument. The mechanics of this
will vary according to the build tool you are using.

(Compiling the Java 8 language level is still supported by using a javac from a
newer JDK, and setting the appropriate `-source`/`-target`/`-bootclasspath` or
`--release` flags).

For instructions on using Error Prone 2.10.0 with JDK 8, see this older version of the installation instructions.

Suppress false positives by adding the suppression annotation @SuppressWarnings("AlwaysThrows") to the enclosing element.

@SuppressWarnings("AlwaysThrows")

Members injection should always be called as early as possible to avoid uninitialized @Inject members. This is also crucial to protect against bugs during configuration changes and reattached Fragments to make sure that each framework type is injected in the appropriate order.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AndroidInjectionBeforeSuper")` to the enclosing element.

Generally when comparing arrays for equality, the programmer intends to check
that the contents of the arrays are equal rather than that they are actually the
same object. But many commonly used equals methods compare arrays for reference
equality rather than content equality. These include the instance `.equals()`
method, Guava’s `com.google.common.base.Objects#equal()`, JDK’s
`java.util.Objects#equals()`, and Android’s
`androidx.core.ObjectsCompat#equals()`.

If reference equality is needed, `==` should be used instead for clarity.
Otherwise, use `java.util.Arrays#equals()` to compare the contents of the
arrays.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ArrayEquals")` to the enclosing element.

`Arrays.fill(Object[], Object)` is used to copy a reference into every slot of
an array.

For example:

```
String[] foo = new String[42];
Arrays.fill(foo, "life");
// 42 references to the same String instance of "life" in the foo array
```
However, because of Array covariance (e.g.: `String[]` is assignable to
`Object[]`), and the signature of Arrays.fill is ```
Arrays.fill(Object[],
Object)
```
, this also allows you to do the following:

```
String[] foo = new String[42];
Arrays.fill(foo, 42); // ArrayStoreException! Integer can't be put into a String[]
```
This check detects the above circumstances, and won’t let you attempt to put
`Integer`s into a `String[]`.

`List<T>` doesn’t have the same issue, since generic types are *not* covariant.

```
List<String> foo = new ArrayList<>();
foo.add(42); // Compile time error: Integer is not assignable to String
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ArrayFillIncompatibleType")` to the enclosing element.

Computing a hashcode for an array is tricky. Typically you want a hashcode that
depends on the value of each element in the array, but many of the common ways
to do this actually return a hashcode based on the *identity* of the array
rather than its contents.

This check flags attempts to compute a hashcode from an array that do not take the contents of the array into account. There are several ways to mess this up:

Call the instance `.hashCode()` method on an array.

Call the JDK method `java.util.Objects#hashCode()` with an argument of array
type.

Call the JDK method `java.util.Objects#hash()` or the Guava method
`com.google.common.base.Objects#hashCode()` with multiple arguments, at
least one of which is an array.

Call the JDK method `java.util.Objects#hash()` or the Guava method
`com.google.common.base.Objects#hashCode()` with a single argument of
*primitive* array type. Because these are varags methods that take
`Object...`, the primitive array is autoboxed into a single-element Object
array, and these methods use the identity hashcode of the primitive array
rather than examining its contents. Note that calling these methods on an
argument of *Object* array type actually does the right thing because no
boxing is needed.

Please use either `java.util.Arrays#hashCode()` (for single-dimensional arrays)
or `java.util.Arrays#deepHashCode()` (for multidimensional arrays) to compute a
hash value that depends on the contents of the array. If you really intended to
compute the identity hash code, consider using
`java.lang.System#identityHashCode()` instead for clarity.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ArrayHashCode")` to the enclosing element.

The `toString` method on an array will print its identity, such as
`[I@4488aabb`. This is almost never needed. Use `Arrays.toString` to print a
human-readable summary.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ArrayToString")` to the enclosing element.

Arrays.asList does not autobox primitive arrays, as one might expect. If you intended to autobox the primitive array, use an asList method from Guava that does autobox. If you intended to create a singleton list containing the primitive array, use Collections.singletonList to make your intent clearer.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ArraysAsListPrimitiveArray")` to the enclosing element.

Methods like Futures.whenAllComplete(…).callAsync(…) will throw a NullPointerException if the provided AsyncCallable returns a null Future. To produce a Future with an output of null, instead return immediateFuture(null).

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AsyncCallableReturnsNull")` to the enclosing element.

Methods like Futures.transformAsync and Futures.catchingAsync will throw a NullPointerException if the provided AsyncFunction returns a null Future. To produce a Future with an output of null, instead return immediateFuture(null).

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AsyncFunctionReturnsNull")` to the enclosing element.

#
 AutoValueBuilderDefaultsInConstructor

 Defaults for AutoValue Builders should be set in the factory method returning Builder instances, not the constructor

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AutoValueBuilderDefaultsInConstructor")` to the enclosing element.

AutoValue constructors are synthesized with their parameters in the same order as the abstract accessor methods. Calls to the constructor need to match this ordering.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AutoValueConstructorOrderChecker")` to the enclosing element.

Implementations of `Annotation` must override `equals` and `hashCode` to match
the expectations defined in that interface, to ensure compatibility with the
annotation instances that source annotations create. Without this, operations
that care about equality of annotations (such as qualified dependency injection
bindings) will fail in mysterious ways.

```
class Foo {
 @SomeAnnotation("hello") public void annotatedMethod() {}
 private static class HelloAnnotationImpl implements SomeAnnotation {
 @Override
 public Class<? extends Annotation> annotationType() {
 return SomeAnnotation.class;
 }
 @Override
 public String value() {
 return "hello";
 }
 }
 static void test() {
 Annotation manual = new HelloAnnotationImpl();
 Annotation fromMethod = Foo.class.getMethod("annotatedMethod").getDeclaredAnnotations()[0];
 manual.equals(fromMethod); // false, violating equality expectations of Annotation!
 }
}
```
It is very difficult to write these methods correctly, so consider using
`@AutoAnnotation` to generate a properly-functioning implementation of
`Annotation`:

```
class Foo {
 @SomeAnnotation("hello") public void annotatedMethod() {}
 @AutoAnnotation
 private static SomeAnnotation someAnnotationInstance(String value) {
 return new AutoAnnotation_Foo_someAnnotationInstance(value);
 }
 static void test() {
 Annotation manual = someAnnotationInstance("hello");
 Annotation fromMethod = Foo.class.getMethod("annotatedMethod").getDeclaredAnnotations()[0];
 manual.equals(fromMethod); // true, hooray!
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("BadAnnotationImplementation")` to the enclosing element.

For shift operations on int types, only the five lowest-order bits of the shift amount are used as the shift distance. This means that shift amounts that are not in the range 0 to 31, inclusive, are silently mapped to values in that range. For example, a shift of an int by 32 is equivalent to shifting by 0, i.e., a no-op.

See JLS 15.19, “Shift Operators”, for more details.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("BadShiftAmount")` to the enclosing element.

JNDI (“Java Naming and Directory Interface”) is a Java JDK API representing an abstract directory service such as DNS, a file system or LDAP. Critically, JNDI allows Java objects to be serialized and deserialized on the wire in implementing systems. This means that if a Java application is allowed to perform a JNDI lookup over some transport protocol then the server it connects to can execute arbitrary attacker-defined code. See this Black Hat talk for more information.

This checker bans usage of every API in the Java JDK that can result in
deserialising an unsafe object via JNDI. The list of APIs is generated from
static callgraph analysis of the JDK, rooted at `javax.naming.Context.lookup`
and is as follows:

`javax.naming.Context.lookup``javax.jdo.JDOHelper.getPersistenceManagerFactory``javax.jdo.JDOHelperTest.testGetPMFBadJNDI``javax.jdo.JDOHelperTest.testGetPMFBadJNDIGoodClassLoader``javax.jdo.JDOHelperTest.testGetPMFNullJNDI``javax.jdo.JDOHelperTest.testGetPMFNullJNDIGoodClassLoader``javax.management.remote.JMXConnectorFactory.connect``javax.management.remote.rmi.RMIConnector.connect``javax.management.remote.rmi.RMIConnector.findRMIServer``javax.management.remote.rmi.RMIConnector.findRMIServerJNDI``javax.management.remote.rmi.RMIConnector.RMIClientCommunicatorAdmin.doStart``javax.naming.directory.InitialDirContext.bind``javax.naming.directory.InitialDirContext.createSubcontext``javax.naming.directory.InitialDirContext.getAttributes``javax.naming.directory.InitialDirContext.getSchema``javax.naming.directory.InitialDirContext.getSchemaClassDefinition``javax.naming.directory.InitialDirContext.modifyAttributes``javax.naming.directory.InitialDirContext.rebind``javax.naming.directory.InitialDirContext.search``javax.naming.InitialContext.doLookup``javax.naming.InitialContext.lookup``javax.naming.spi.ContinuationContext.lookup``javax.naming.spi.ContinuationDirContext.bind``javax.naming.spi.ContinuationDirContext.createSubcontext``javax.naming.spi.ContinuationDirContext.getAttributes``javax.naming.spi.ContinuationDirContext.getSchema``javax.naming.spi.ContinuationDirContext.getSchemaClassDefinition``javax.naming.spi.ContinuationDirContext.getTargetContext``javax.naming.spi.ContinuationDirContext.modifyAttributes``javax.naming.spi.ContinuationDirContext.rebind``javax.naming.spi.ContinuationDirContext.search``javax.sql.rowset.spi.ProviderImpl.getDataSourceLock``javax.sql.rowset.spi.ProviderImpl.getProviderGrade``javax.sql.rowset.spi.ProviderImpl.getRowSetReader``javax.sql.rowset.spi.ProviderImpl.getRowSetWriter``javax.sql.rowset.spi.ProviderImpl.setDataSourceLock``javax.sql.rowset.spi.ProviderImpl.supportsUpdatableView``javax.sql.rowset.spi.SyncFactory.enumerateBindings``javax.sql.rowset.spi.SyncFactory.getInstance``javax.sql.rowset.spi.SyncFactory.initJNDIContext``javax.sql.rowset.spi.SyncFactory.parseJNDIContext`A small subset of these are banned directly. The rest are banned indirectly by
banning the `lookup()`, `bind()`, `rebind()`, `getAttributes()`,
`modifyAttriutes()`, `createSubcontext()`, `getSchema()`,
`getSchemaClassDefinition()` and `search()` methods on any subclass
(implementer) of `javax.naming.Context`. The indirect ban is necessary due to
these methods being vulnerable in previously noted subclasses in the JDK. If
they were not banned at the Context level, a cast to Context would make the
vulnerable call invisible to static analysis.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("BanJNDI")` to the enclosing element.

*Alternate names: NumericEquality*

Comparison using reference equality instead of value equality. The inputs to this comparison are boxed primitive types, where reference equality is particularly bug-prone: Primitive wrapper classes cache instances for some (but usually not all) values, so == may be equivalent to equals() for some values but not others. Additionally, not all versions of the runtime and other libraries use the cache in the same cases, so upgrades may change behavior. Furthermore, reference identity is usually not useful for primitive wrappers, as they are immutable types whose equals() method fully compares their values.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("BoxedPrimitiveEquality")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("BundleDeserializationCast") to the enclosing element.

@SuppressWarnings("BundleDeserializationCast")

When a class exposes multiple constructors, they’re generally used as a means of initializing default parameters. If a chaining constructor ignores a parameter, it’s likely the parameter needed to be plumbed to the chained constructor.

```
MissileLauncher(Location target) {
 this(target, false);
}
MissileLauncher(boolean askForConfirmation) {
 this(TEST_TARGET, false); // should be askForConfirmation
}
MissileLauncher(Location target, boolean askForConfirmation) {
 ...
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ChainingConstructorIgnoresParameter")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("CheckNotNullMultipleTimes") to the enclosing element.

@SuppressWarnings("CheckNotNullMultipleTimes")

*Alternate names: ResultOfMethodCallIgnored, ReturnValueIgnored*

When code calls a non-`void` method, it should usually use the value that the
method returns.

Consider the following code, which ignores the return value of `concat`:

```
string.concat("\n");
```
That code is a no-op because `concat` doesn’t modify `string`; it returns a new
string for the caller to use, as in:

```
string = string.concat("\n");
```
To avoid this bug, Error Prone requires callers to use the return value of
`concat` and some other well-known methods.

Additionally, Error Prone can be configured to require callers to use the return value of any methods that you choose.

Most methods are like `concat`: Calls to those methods must use their return
values.

However, there are exceptions. For example, `set.add(element)` returns a
`boolean`: The return value is `false` if `element` was *already* contained in
`set`. Typically, callers don’t need to know this, so they don’t need to use the
return value.

For Error Prone’s `CheckReturnValue` check to be useful, it needs to know which
methods are like `concat` and which are like `add`.

`@CheckReturnValue` and `@CanIgnoreReturnValue`The `@CheckReturnValue` annotation (available in JSR-3051 or in
Error Prone) marks methods whose return values must be used. This error
is triggered when one of these methods is called but the result is not used.

`@CheckReturnValue` may be applied to a class or package 2 to
indicate that all methods in that class or package must have their return values
used.

For convenience, we provide an annotation, `@CanIgnoreReturnValue`, to
exempt specific methods or classes from this behavior. `@CanIgnoreReturnValue`
is available from the Error Prone annotations package,
`com.google.errorprone.annotations`.

If you really want to ignore the return value of a method annotated with
`@CheckReturnValue`, a cleaner alternative to `@SuppressWarnings` is to assign
the result to a variable that starts with `unused`:

```
public void setNameFormat(String nameFormat) {
 String unused = format(nameFormat, 0); // fail fast if the format is bad or null
 this.nameFormat = nameFormat;
}
```
`@CheckReturnValue` is ignored under the following conditions (which saves users
from having to use either an `unused` variable or `@SuppressWarnings`):

Calls from `Mockito.verify()` or `Stubber.when()`; e.g.,
`Mockito.verify(t).foo()` or `doReturn(val).when(t).foo()` (where `foo()` is
annotated with `@CheckReturnValue`). Here, the method calls are just used to
program the mock object, not to be consumed directly.

Code that does exception testing with JUnit, where the intent is that the method call should throw an exception:

Uses of JUnit 4.13 or JUnit5’s `assertThrows` methods:

```
assertThrows(IndexOutOfBoundsException.class, () -> list.get(-1));
```
The `try/execute/fail/catch` pattern

```
 try {
 list.get(-1);
 fail("Expected a IndexOutOfBoundsException to be thrown on a negative index");
 } catch (IndexOutOfBoundsException expected) {
 }
```
JUnit’s `ExpectedException`

```
expectedException.expect(IndexOutOfBoundsException.class);
list.get(-1); // If this throws IOOBE, the test passes.
```
Of note, the JSR-305 project was never fully approved, so the JSR-305 version of the annotation is not actually official and causes issues with Java 9 and the Module System. Prefer to use the Error Prone version. ↩

To annotate a package, create a
`package-info.java` file in the package directory, add a package statement,
and annotate the package statement. ↩

Querying a collection for an element it cannot possibly contain is almost certainly a bug.

In a generic collection type, query methods such as `Map.get(Object)` and
`Collection.remove(Object)` accept a parameter that identifies a potential
element to *look for* in that collection. This check reports cases where this
element *cannot* be present because its type and the collection’s generic
element type are “incompatible.” A typical example:

```
Set<Long> values = ...
if (values.contains(42)) { ... }
```
This code looks reasonable, but there’s a problem: The `Set` contains `Long`
instances, but the argument to `contains` is an `Integer`. Because no instance
can be of type `Integer` and of type `Long` at the same time, the `contains`
check always fails. This is clearly not what the developer intended.

Why does the collection API permit this kind of mistake? Why not declare the
method as `contains(E)`? After all, that is what the collections API does for
methods that *store* an element in the collection: They require that passed type
be strictly *assignable to* the collection’s element type. For example:

```
void addIntegerOne(Set<? extends Number> numbers) {
 numbers.add(42); // won't compile
}
```
The code above rightly won’t compile, because `numbers` *might* be a (for
example) `Set<Long>`, and adding an `Integer` value would corrupt it.

But this restriction is necessary only for methods that insert elements. Methods that only query or remove elements cannot corrupt the collection:

```
void removeIntegerOne(Set<? extends Number> numbers) {
 numbers.remove(42); // should compile (and does)
}
```
In this case, the `Integer` `42` might be contained in `numbers`, and should be
removed if it is, but if `numbers` is a `Set<Long>`, no harm is done.

We’d like to define `contains` in a way that rejects the bad call but permits
the good one. But Java’s type system is not powerful enough. Our solution is
static analysis.

The specific restriction we would like to express for the two types is not
assignability, but “compatibility”. Informally, we mean that it must at least be
*possible* for some instance to be of both types. Formally, we require that a
“casting conversion” exist between the types as defined by
JLS 5.5.1.

The result is that the method can be defined as `contains(Object)`, permitting
the “good” call above, but that Error Prone will give errors for incompatible
arguments, preventing the “bad.”

`E` have been better?We might say: Sure, a buggy `remove` call can’t corrupt a collection. And sure,
someone might want to pass an `Object` reference that happens to contain an `E`.
But isn’t that a low standard for an API? We don’t normally write code that way:

```
void throwIfUnchecked(Object throwable) { // no "need" to require Throwable
 if (throwable instanceof RuntimeException) {
 throw (RuntimeException) throwable;
 }
 if (throwable instanceof Error) {
 throw (Error) throwable;
 }
}
```
Such code would invite bugs. To avoid that, we require a `Throwable`. Users who
have an `Object` reference that might be a `Throwable` can test `instanceof` and
cast. So why not require the same thing in the collections API?

Of course, we can’t really change the API of `Collection`. But if we were
designing a similar API, what would we do – require `E` or accept any `Object`?

The burden of proof falls on accepting `Object`, since doing so permits buggy
code. And we’re not going to settle for “it occasionally saves users a cast.”

*The main reason to accept Object is to permit a fast, type-safe Set.equals
implementation.*

(`equals` is actually just one example of the general problem, which arises with
many uses of wildcards. Once you’ve read the following, consider the problem of
implementing `Collection.removeAll(Collection<?>)` without `contains(Object)`.
Then consider how the problem would exist even if the signature were
`removeAll(Collection<? extends E>)`. The `removeAll` problem is at least
“solvable” by changing the signature to `removeAll(Collection<E>)`, but that
signature may reject useful calls.)

Here’s how: `equals` necessarily accepts a plain `Object`. It can test whether
that `Object` is a `Set`, but it can’t know the element type it was originally
declared with. In short, `equals` has to operate on a `Set<?>`.

If `contains` were to require an `E`, `equals` would be in trouble because it
doesn’t know what `E` is. In particular, it wouldn’t be able to call
`otherSet.contains(myElement)` for any of its elements.

It would have only two options: It could copy the entire other `Set` into a
`Set<Object>`, or it could perform an unchecked cast. Copying is wasteful, so in
practice, `equals` would need an unchecked cast. This is probably acceptable,
but we might feel strange for defining an API that can be implemented only by
performing unchecked casts.

Does a cleaner implementation (and occasional convenience to callers) outweigh
the bugs that accepting `Object` enables? That’s a tough question. The good news
is that this Error Prone check gives you some of the best of both worlds.

It is technically possible for a `Set<Integer>` to contain a `String` element,
but only if an `unchecked` warning was earlier ignored or improperly suppressed.
Such practice should never be treated as acceptable, so it makes no practical
difference to our arguments above.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("CollectionIncompatibleType")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("CollectionToArraySafeParameter") to the enclosing element.

@SuppressWarnings("CollectionToArraySafeParameter")

The type argument of `Comparable` should always be the type of the current
class.

For example, do this:

```
class Foo implements Comparable<Foo> {
 public int compareTo(Foo other) { ... }
}
```
not this:

```
class Foo implements Comparable<Bar> {
 public int compareTo(Bar other) { ... }
}
```
Implementing `Comparable` for a different type breaks the API contract, which
requires `x.compareTo(y) == -y.compareTo(x)` for all `x` and `y`. If `x` and `y`
are different types, this behaviour can’t be guaranteed.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ComparableType")` to the enclosing element.

The boolean expression this != null always returns true and similarly this == null always returns false.

Suppress false positives by adding the suppression annotation @SuppressWarnings("ComparingThisWithNull") to the enclosing element.

@SuppressWarnings("ComparingThisWithNull")

This checker looks for comparisons to values that are too high or too low for the compared type. For example, bytes may have a value in the range -128 to 127. Comparing a byte for equality with a value outside that range will always evaluate to false and usually indicates an error in the code.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ComparisonOutOfRange")` to the enclosing element.

The `@CompatibleWith` annotation is used to mark parameters that need extra type
checking on arguments passed to the method. The annotation was not appropriately
placed on a parameter with a valid type argument. See the javadoc for more
details.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("CompatibleWithAnnotationMisuse")` to the enclosing element.

A method or constructor with one or more parameters whose declaration is
annotated with the `@CompileTimeConstant` type annotation must only be invoked
with corresponding actual parameters that are computed as compile-time constant
expressions, specifically expressions that:

`@CompileTimeConstant` annotation, or`String`s,
orFor example, the following are valid compile-time constants:

`"some literal string"``"literal string" + compileTimeConstantParameter``debug ? compileTimeConstantParameter : "foo"``ImmutableList.of("a", "b", "c")`When applied to fields, this check enforces that the field is `final` and has an
initializer which satisfies the above conditions.

Getting Java 8 references to methods with `@CompileTimeConstant` parameters is
disallowed because we couldn’t check if the method reference is later applied to
a compile-time constant. Use the methods directly instead.

For the same reason, it’s also disallowed to create lambda expressions with
`@CompileTimeConstant` parameters.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("CompileTimeConstant")` to the enclosing element.

With ambiguous constructor references used in `java.util.Map#computeIfAbsent`
function parameter, it becomes unclear what the code intends to do.

```
map.computeIfAbsent(someLong, AtomicLong::new).incrementAndGet()
```
Code of this form can seemingly look like it’s trying to make a counter,
creating the key if absent. Unfortunately the code is surprising, because it
will call the wrong `AtomicLong` constructor. Instead, it will create the
counter initialized with the key value, which is probably not desired.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ComputeIfAbsentAmbiguousReference")` to the enclosing element.

A conditional expression with numeric second and third operands of differing types may give surprising results.

For example:

```
Object t = true ? Double.valueOf(0) : Integer.valueOf(0);
System.out.println(t.getClass()); // class java.lang.Double
Object f = false ? Double.valueOf(0) : Integer.valueOf(0);
System.out.println(f.getClass()); // class java.lang.Double !!
```
Despite the apparent intent to get a `Double` in one case, and an `Integer` in
the other, the result is a `Double` in both cases.

This is because the rules in JLS § 15.25.2 state that differing numeric types will undergo binary numeric promotion. As such, the latter case is evaluated as:

```
Object f =
 Double.valueOf(
 false
 ? Double.valueOf(0).doubleValue()
 : (double) Integer.valueOf(0).intValue());
```
To get a different type in the two cases, one can either explicitly cast the operands to a non-boxable type:

```
Object f = false ? ((Object) Double.valueOf(0)) : ((Object) Integer.valueOf(0));
System.out.println(t.getClass()); // class java.lang.Integer
```
Or use if/else:

```
Object f;
if (false) {
 f = Double.valueOf(0);
} else {
 f = Integer.valueOf(0);
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ConditionalExpressionNumericPromotion")` to the enclosing element.

Compile-time constant expressions that overflow are a potential source of bugs.

Literals without an explicit `L` suffix have type `int`, so the following
multiplication expression is evaluated as an integer before being widened to
`long`. The value is greater than `Integer.MAX_VALUE`, so it wraps around to
`-1857093632`.

```
static final long NANOS_PER_DAY = 24 * 60 * 60 * 1000 * 1000 * 1000;
```
The intent was probably for the multiplication expression to be evaluated as a
`long` instead of an `int`.

```
static final long NANOS_PER_DAY = 24L * 60 * 60 * 1000 * 1000 * 1000;
```
If you find yourself doing this kind of time-based math, consider using an API
that provides a safer, more readable and strongly-typed solution like the
`java.time.Duration`
API.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ConstantOverflow")` to the enclosing element.

Dagger `@Provides` methods may not return null unless annotated with
`@Nullable`. Such a method will cause a `NullPointerException` at runtime if the
`return null` path is ever taken.

If you believe the `return null` path can never be taken, please throw a
`RuntimeException` instead. Otherwise, please annotate the method with
`@Nullable`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DaggerProvidesNull")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("DangerousLiteralNull") to the enclosing element.

@SuppressWarnings("DangerousLiteralNull")

*Alternate names: ThrowableInstanceNeverThrown*

The exception is created with `new`, but is not thrown, and the reference is
lost.

Creating an exception without using it is unlikely to be correct, so we assume that you wanted to throw the exception.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DeadException")` to the enclosing element.

The Thread is created with `new`, but is never started and is not otherwise
captured.

Threads must be started with `start()` to actually execute.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DeadThread")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("DereferenceWithNullBranch") to the enclosing element.

@SuppressWarnings("DereferenceWithNullBranch")

Suppress false positives by adding the suppression annotation @SuppressWarnings("DiscardedPostfixExpression") to the enclosing element.

@SuppressWarnings("DiscardedPostfixExpression")

This check prevents calls to methods annotated with Error Prone’s `@DoNotCall`
annotation (`com.google.errorprone.annotations.DoNotCall`).

The check disallows invocations and method references of the annotated method.

There are a few situations where this can be useful, including methods that are required to satisfy the contract of an interface, but that are not supported.

A method annotated with `@DoNotCall` should always be `final` or `abstract`. If
an `abstract` method is annotated `@DoNotCall` Error Prone will ensure all
implementations of that method also have the annotation. Methods annotated with
`@DoNotCall` should *not* be private, since a private method that should not be
called can simply be removed.

TIP: Marking methods annotated with `@DoNotCall` as `@Deprecated` is
recommended, since it provides IDE users with more immediate feedback.

Example:

`java.util.Collection#add` should never be called on an immutable collection
implementation:

```
package com.google.common.collect.ImmutableList;
class ImmutableList<E> implements List<E> {
 // ...
 /**
 * Guaranteed to throw an exception and leave the list unmodified.
 *
 * @deprecated Unsupported operation.
 */
 @Deprecated
 @DoNotCall("guaranteed to throw an exception and leave the list unmodified")
 @Override
 public final void add(E e) {
 throw new UnsupportedOperationException();
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("DoNotCall")` to the enclosing element.

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
DoNotMock
Identifies undesirable mocks.
Severity
ERROR

The double-brace initialization pattern should be avoided—especially in non-static contexts.

The double-brace pattern uses an instance-initializer in an anonymous inner class to express the initialization of a class (often a collection) in a single step.

Inner classes in a non-static context are terrific sources of memory leaks! If you pass the collection somewhere that retains it, the entire instance you created it from can no longer be garbage collected. Even if it is completely unreachable. And if someone serializes the map? Yep, the entire creating instance goes along for the ride (or if that fails, serializing the map fails, which is also awfully strange). All this is completely nonobvious.

Luckily, there are more readable and more performant alternatives in the factory
methods and builders for `ImmutableList`, `ImmutableSet`, and `ImmutableMap`.

The `List.of`, `Set.of`, and `Map.of` static factories
added in Java 9 are also a good choice.

That is, prefer this:

```
ImmutableList.of("Denmark", "Norway", "Sweden");
```
Not this:

```
new ArrayList<>() {
 {
 add("Denmark");
 add("Norway");
 add("Sweden");
 }
};
```
TIP: Neither the guava immutable collections nor the static factory methods
added in a JDK 9 support `null` elements. The double-brace pattern is still best
avoided for collections that contain null. Consider using `Arrays.asList` to
initialize `List`s and `Set`s with `null` values, and refactoring `Map`
initializers into a helper method.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DoubleBraceInitialization")` to the enclosing element.

JDK 9 has
`Map#ofEntries`
factory which throws runtime error when provided multiple entries with the same
key.

For eg, the following code is erroneously adding two entries with `Foo` as key.

```
Map<String, String> map = Map.ofEntries(
 Map.entry("Foo", "Bar"),
 Map.entry("Ping", "Pong"),
 Map.entry("Kit", "Kat"),
 Map.entry("Foo", "Bar"));
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("DuplicateMapKeys")` to the enclosing element.

Duration.from(TemporalAmount) will always throw a UnsupportedTemporalTypeException when passed a Period and return itself when passed a Duration.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DurationFrom")` to the enclosing element.

`Duration.get(TemporalUnit)` only works when passed `ChronoUnit.SECONDS` or `ChronoUnit.NANOS`. All other values are guaranteed to throw a `UnsupportedTemporalTypeException`. In general, you should avoid `duration.get(ChronoUnit)`. Instead, please use `duration.toNanos()`, `Durations.toMicros(duration)`, `duration.toMillis()`, `duration.getSeconds()`, `duration.toMinutes()`, `duration.toHours()`, or `duration.toDays()`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DurationGetTemporalUnit")` to the enclosing element.

Duration APIs only work for TemporalUnits with exact durations or ChronoUnit.DAYS. E.g., Duration.of(1, ChronoUnit.YEARS) is guaranteed to throw a DateTimeException.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DurationTemporalUnit")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("DurationToLongTimeUnit") to the enclosing element.

@SuppressWarnings("DurationToLongTimeUnit")

The contract for `Object.hashCode` states that if two objects are equal, then
calling the `hashCode()` method on each of the two objects must produce the same
result. Implementing `equals()` but not `hashCode()` causes broken behaviour
when trying to store the object in a collection.

See Effective Java 3rd Edition §11 for more information and a
discussion of how to correctly implement `hashCode()`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("EqualsHashCode")` to the enclosing element.

As per JLS 15.21.1, == NaN comparisons always return false, even NaN == NaN. Instead, use the isNaN methods to check for NaN.

Suppress false positives by adding the suppression annotation @SuppressWarnings("EqualsNaN") to the enclosing element.

@SuppressWarnings("EqualsNaN")

The contract of `Object.equals()` states that for any non-null reference value
`x`, `x.equals(null)` should return `false`. Thus code such as

```
if (x.equals(null)) {
 ...
}
```
either returns `false`, or throws a `NullPointerException` if `x` is `null`. The
nested block may never execute.

This check replaces `x.equals(null)` with `x == null`, and `!x.equals(null)`
with `x != null`. If the author intended for `x.equals(null)` to return `true`,
consider this as fragile code as it breaks the contract of `Object.equals()`.

See Effective Java 3rd Edition §10: Objey the general contract when overriding equals for more details.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("EqualsNull")` to the enclosing element.

.equals() to the same object will result in infinite recursion

Suppress false positives by adding the suppression annotation @SuppressWarnings("EqualsReference") to the enclosing element.

@SuppressWarnings("EqualsReference")

An equals method compares non-corresponding fields from itself and the other instance:

```
class Frobnicator {
 private int a;
 private int b;
 @Override
 public boolean equals(@Nullable Object other) {
 if (!(other instanceof Frobnicator)) {
 return false;
 }
 Frobnicator that = (Frobnicator) other;
 return a == that.a && b == that.a; // BUG: should be b == that.b
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("EqualsWrongThing")` to the enclosing element.

Alternate names: FormatString

Suppress false positives by adding the suppression annotation @SuppressWarnings("FloggerFormatString") to the enclosing element.

@SuppressWarnings("FloggerFormatString")

#
 FloggerLogString

 Arguments to log(String) must be compile-time constants or parameters annotated with @CompileTimeConstant. If possible, use Flogger's formatting log methods instead.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("FloggerLogString")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("FloggerLogVarargs") to the enclosing element.

@SuppressWarnings("FloggerLogVarargs")

Suppress false positives by adding the suppression annotation @SuppressWarnings("FloggerSplitLogStatement") to the enclosing element.

@SuppressWarnings("FloggerSplitLogStatement")

A method that overrides a @ForOverride method should not be invoked directly. Instead, it should be invoked only from the class in which it was declared. For example, if overriding Converter.doForward, you should invoke it through Converter.convert. For testing, factor out the code you want to run to a separate method.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ForOverride")` to the enclosing element.

Format strings for the printf family of functions must follow the specification in the documentation for java.util.Formatter.

The syntax for format specifiers is:

```
%[argument_index$][flags][width][.precision]conversion
```
Format strings can have the following errors:

Duplicate flags are provided in the format specifier:

```
String.format("e = %++10.4f", Math.E);
```
A conversion and flag are incompatible:

```
String.format("%#b", Math.E);
```
The argument is a character with an invalid Unicode code point.

```
String.format("%c", 0x110000);
```
The argument corresponding to the format specifier is of an incompatible type:

```
String.format("%f", "abcd");
```
An illegal combination of flags is given:

```
String.format("%-010d", 5);
```
The conversion does not support a precision:

```
String.format("%.c", 'c');
```
The conversion does not support a width:

```
String.format("%1n");
```
There is a format specifier which does not have a corresponding argument or if an argument index refers to an argument that does not exist:

```
String.format("%<s", "test");
```
The format width is required:

```
String.format("e = %-f", Math.E);
```
An unknown conversion is given:

```
String.format("%r", "hello");
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("FormatString")` to the enclosing element.

Methods can be annotated with Error Prone’s `@FormatMethod` annotation to
indicate that calls to this function should be treated similarly to
`String.format`: One of the parameters is a ‘format string’ (the first String
parameter or the only parameter annotated with `@FormatString`), and the
subsequent parameters are used as format arguments to that format string.

For example:

```
@FormatMethod
void myLogMethod(@FormatString String fmt, Object... args) {}
// ERROR: 2nd format argument isn't a number
myLogMessage("My log message: %d and %d", 3, "has a message");
```
In order to avoid complex runtime issues when the format string part is
dynamically constructed, leading to a mismatch between the arguments and format
strings, we require that the ‘format string’ argument in calls to
`@FormatMethod`-annotated methods be one of:

`@FormatString`-annotated variableWe will then check that the format string and format arguments match.

For more information on possible format string errors, see the documentation on the FormatString check.

The import for `@FormatMethod` is:

```
import com.google.errorprone.annotations.FormatMethod;
```
Suppress false positives by adding the suppression annotation @SuppressWarnings(“FormatStringAnnotation”) to the enclosing element.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("FormatStringAnnotation")` to the enclosing element.

Not all java.time types can be created via from(TemporalAccessor). For example, you can create a Month from a LocalDate (Month.from(localDate)) because a LocalDate consists of a year, month, and day. However, you cannot create a LocalDate from a Month (since it doesn’t have the year or day information). Instead of throwing a DateTimeException at runtime, this checker validates the type transformations at compile time using static type information.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("FromTemporalAccessor")` to the enclosing element.

There are only two correct ways to have one `@FunctionalInterface` extend
another `@FunctionalInterface`—the way where you leave the abstract method
alone, and the way where you make a different abstract method but a default
method for the original abstract method simply delegates to the new name.

If you do anything else, you create a situation where what looks like a “cast” actually changes behavior. That is really quite bad for understandability.

For example, if the method `bar()` changes the behaviour of `qux()`, then the
same lambda cast to `A` or `B` could have completely different behaviour.

```
@FunctionalInterface
interface A {
 Foo bar();
}
@FunctionalInterface
interface B extends A {
 Foo qux();
 @Override
 default Foo bar() {
 // anything here but exactly `return qux();` or perhaps `return (SomeType) qux();`
 }
}
```
If you need to change the behaviour of an existing lambda, prefer an explicit decorator method, e.g.:

```
static Runnable crashTerminating(Runnable r) {
 return () -> { ...wrapping behavior goes here... }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("FunctionalInterfaceMethodChanged")` to the enclosing element.

The passed exception type must not be a RuntimeException, and it must expose a public constructor whose only parameters are of type String or Throwable. getChecked will reject any other type with an IllegalArgumentException.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("FuturesGetCheckedIllegalExceptionType")` to the enclosing element.

From documentation: DoubleMath.fuzzyEquals is not transitive, so it is not suitable for use in Object#equals implementations.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("FuzzyEqualsShouldNotBeUsedInEqualsMethod")` to the enclosing element.

Instances of an annotation interface generally return a random proxy class when
`getClass()` is called on them; to get the actual annotation type use
`annotationType()`.

In the following example, calling `getClass()` on the annotation instance
returns a proxy class like `com.sun.proxy.$Proxy1`, while `annotationType()`
returns `Deprecated`.

```
@Deprecated
public class Test {
 static void printAnnotationClass(Annotation annotation) {
 System.err.println(annotation.getClass());
 System.err.println(annotation.annotationType());
 }
 public static void main(String[] args) {
 printAnnotationClass(Test.class.getAnnotation(Deprecated.class));
 }
}
```
Prints:

```
class com.sun.proxy.$Proxy1
interface java.lang.Deprecated
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("GetClassOnAnnotation")` to the enclosing element.

Calling `getClass()` on an object of type Class returns the Class object for
java.lang.Class. Usually this is a mistake, and people intend to operate on the
object itself (for example, to print an error message). If you really did intend
to operate on the Class object for java.lang.Class, please use `Class.class`
instead for clarity.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("GetClassOnClass")` to the enclosing element.

*Alternate names: GuardedByChecker*

The GuardedBy analysis checks that fields or methods annotated with
`@GuardedBy(lock)` are only accessed when the specified lock is held.

Example:

```
import com.google.errorprone.annotations.concurrent.GuardedBy;
class Account {
 @GuardedBy("this")
 private int balance;
 public synchronized int getBalance() {
 return balance; // OK: implicit 'this' lock is held.
 }
 public synchronized void withdraw(int amount) {
 setBalance(balance - amount); // OK: implicit 'this' lock is held.
 }
 public void deposit(int amount) {
 setBalance(balance + amount); // ERROR: access to 'balance' not guarded by 'this'.
 }
 @GuardedBy("this")
 private void setBalance(int newBalance) {
 checkState(newBalance >= 0, "Balance cannot be negative.");
 balance = newBalance; // OK: 'this' must be held by caller of 'setBalance'.
 }
}
```
This above example uses implicit locks (via the ‘synchronized’ modifier). The analysis also supports synchronized statements and java.util.concurrent locks.

The analysis provides a way of associating members with locks. A member is a field or a method. A lock can be the implicit lock of an object, or a java.util.concurrent Lock.

An implicit lock is acquired using the built in synchronization features of the language. Adding the ‘synchronized’ modifier to an instance method causes the implicit lock of the enclosing instance to be acquired for the duration of the method. Adding the ‘synchronized’ modifier to a static method is similar, except the implicit lock of the Class object is acquired instead.

The Locks defined in java.util.concurrent are acquired with explicit lock()/unlock() methods. The use of these methods in Java should always correspond to a try/finally block, to ensure that the locks are released on all execution paths.

The following syntax can be used to describe a lock:

| ```
this
```
 | The implicit object lock of the enclosing class. |
| ```
ClassName.this
```
 | The implicit object lock of the enclosing class specified by ClassName. (For inner classes, the ClassName.this designation allows you to specify which 'this' reference is intended.) |
| ```
fieldName
```
 | The final instance field specified by fieldName. |
| ```
methodName()
```
 | The instance method specified by methodName(). Methods called to return locks should be deterministic. |
| ```
ClassName.class
```
 | The implicit lock of specified Class object. |
| ```
ClassName.fieldName
```
 | The static final field specified by fieldName. |
| ```
ClassName.methodName()
```
 | The static method specified by methodName(). Methods called to return locks should be deterministic. |
| ```
itself
```
 | The annotated field. |

com.google.errorprone.annotations.concurrent.GuardedBy

The @GuardedBy annotation is used to document that a member (a field or a method) can only be accessed when the specified lock is held.

@GuardedBy can be used with both implicit locks and java.util.concurrent Locks.

```
final Lock lock = new ReentrantLock();
@GuardedBy("lock")
int x;
void m() {
 x++; // error: access of 'x' not guarded by 'lock'
 lock.lock();
 try {
 x++; // OK: guarded by 'lock'
 } finally {
 lock.unlock();
 }
}
```
Note: there are a couple more annotations called `@GuardedBy`, including
`javax.annotation.concurrent.GuardedBy` and
`org.checkerframework.checker.lock.qual.GuardedBy`. The check recognizes those
versions of the annotation, but we recommend using
`com.google.errorprone.annotations.concurrent.GuardedBy`.

Anonymous classes and lambdas need to re-acquire locks that may be held by an enclosing block. For example, consider:

```
class Transaction {
 @GuardedBy("this")
 int x;
 public synchronized void handle() {
 doSomething(() -> {
 x++; // Error: access of 'x' not guarded by 'Transaction.this'
 });
 }
}
```
The analysis is intra-procedural, meaning it doesn’t consider the implementation
of `doSomething`.

In general, the checker doesn’t know if `doSomething` immediately calls the
provided lambda while the lock is still held by the enclosing method `handle`,
for example:

```
private void doSomething(Runnable r) {
 r.run();
}
```
… or whether the lambda could be called later after `handle` has released the
lock, for example:

```
private void doSomething(Runnable r) {
 // runs `r` at some point in the future
 someExecutor.execute(r);
}
```
However, the check does special-case some method calls which are known to immediately call the provided lambda or method reference.

```
class Names {
 @GuardedBy("this")
 List<String> names = new ArrayList<>();
 public void addName(String name) {
 List<String> copyOfNames;
 synchronized (this) {
 copyOfNames = names; // OK: access of 'names' guarded by 'this'
 }
 copyOfNames.add(name); // should be an error: this access is not thread-safe!
 }
}
```
The analysis does not track aliasing, so it’s possible to circumvent the safety it provides by copying references to guarded members.

In the example, the guarded field ‘names’ can be accessed via a copy even if the required lock is not held.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("GuardedBy")` to the enclosing element.

Classes that AssistedInject factories create may not be annotated with scope annotations, such as @Singleton. This will cause a Guice error at runtime.

See this bug report for details.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("GuiceAssistedInjectScoping")` to the enclosing element.

From the javadoc of FactoryModuleBuilder:

The types of the factory method’s parameters must be distinct. To use multiple parameters of the same type, use a named

`@Assisted`annotation to disambiguate the parameters. The names must be applied to the factory method’s parameters:

```
public interface PaymentFactory {
 Payment create(
 @Assisted("startDate") Date startDate,
 @Assisted("dueDate") Date dueDate,
 Money amount);
 }
```
…and to the concrete type’s constructor parameters:

```
public class RealPayment implements Payment {
 @Inject
 public RealPayment(
 CreditService creditService,
 AuthService authService,
 @Assisted("startDate") Date startDate,
 @Assisted("dueDate") Date dueDate,
 @Assisted Money amount) {
 ...
 }
 }
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("GuiceAssistedParameters")` to the enclosing element.

From the Guice wiki:

Injecting

`final`fields is not recommended because the injected value may not be visible to other threads.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("GuiceInjectOnFinalField")` to the enclosing element.

`Hashtable.contains(Object)` and `ConcurrentHashMap.contains(Object)` are legacy
methods for testing if the given object is a value in the hash table. They are
often mistaken for `containsKey`, which checks whether the given object is a
*key* in the hash table.

If you intended to check whether the given object is a key in the hash table,
use `containsKey` instead. If you really intended to check whether the given
object is a value in the hash table, use `containsValue` for clarity.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("HashtableContains")` to the enclosing element.

*Alternate names: SelfEquality*

Using the same expressions as both arguments to the following binary expressions is usually a mistake:

`a && a`, `a || a`, `a & a`, or `a | a` is equivalent to `a``a <= a`, `a >= a`, or `a == a` is always `true``a < a`, `a > a`, `a != a`, or `a ^ a` is always `false``a / a` is always `1``a % a` or `a - a` is always `0`If the expression has side-effects, consider refactoring one of the expressions with side effects into a local. For example, prefer this:

```
// check twice, just to be sure
boolean isTrue = foo.isTrue();
if (isTrue && foo.isTrue()) {
 // ...
}
```
to this:

```
if (foo.isTrue() && foo.isTrue()) {
 // ...
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("IdentityBinaryExpression")` to the enclosing element.

Usage of `java.util.IdentityHashMap` with a boxed primitive type as a key is
risky and can yield unexpected results because `java.util.IdentityHashMap` uses
reference-equality when comparing keys and reference equality for primitive
wrappers is particularly bug-prone: Primitive wrapper classes cache instances
for some (but usually not all) values, so == may be equivalent to equals() for
some values but not others. Additionally, not all versions of the runtime and
other libraries use the cache in the same cases, so upgrades may change
behavior.

Thus:

```
 Map<Integer, Foo> map = new IdentityHashMap<>();
 int n = randomInt();
 map.put(n, x);
 map.get(n); // This could be null since boxing happens twice and could produce distinct values.
```
But:

```
 Map<Integer, Foo> map = new IdentityHashMap<>();
 Integer n = randomInt();
 map.put(n, x);
 map.get(n); // This cannot be null because it's the same instance.
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("IdentityHashMapBoxing")` to the enclosing element.

This check validates that all classes annotated with Error Prone’s `@Immutable`
annotation (`com.google.errorprone.annotations.Immutable`) are deeply immutable.
It also checks that any class extending an `@Immutable`-annotated class or
implementing an `@Immutable`-annotated interface are also immutable.

NOTE: Other versions of the annotation, such as
`javax.annotation.concurrent.Immutable`, are currently *not* enforced.

An object is immutable if its state cannot be observed to change after construction. Immutable objects are inherently thread-safe.

A class is immutable if all instances of that class are immutable. The immutability of a class can only be fully guaranteed if the class is final, otherwise one must ensure all subclasses are also immutable.

A conservative definition of object immutability is:

`this` reference does not escape the
constructor).The requirement that all reference fields be immutable ensures *deep*
immutability, meaning all contained state is also immutable. A weaker property,
common with container classes, is *shallow* immutability, which allows some of
the object’s fields to point to mutable objects. One example of shallow
immutability is guava’s ImmutableList, which may contain mutable elements.

It is possible to implement immutable classes with some internal mutable state, as long as callers can never observe changes to that state. For example, some state may be lazily initialized to improve performance.

It is also technically possible to have an immutable object with non-final
fields (see the implementation of `String#hashCode()` for an example), but doing
this correctly requires subtle reasoning about safe data races and deep
knowledge of the Java Memory Model.

If you have an immutable class with mutable fields as described above, you can
mark it as such by suppressing the Immutable check on it. This
will allow your class to be included in other `@Immutable` classes.

For more information about immutability, see:

When an `@Immutable` class has type parameters that are used in the type of that
class’s fields, that class is called an immutable generic container. Usages of
immutable generic container classes, such as `ImmutableList`, are only actually
deemed immutable if the arguments to all such type parameters are also deemed
immutable. For example, an `ImmutableList<String>` is deemed immutable since
`String`s are immutable. However, an `ImmutableList<Object>` is not deemed
immutable since `Object`s are not provably immutable.

When creating generic container classes, Error Prone requires that you declare whether that container is allowed to be used with mutable, or only with immutable type parameters.

`@Immutable(containerOf = ...)`If you want to allow your immutable generic container to possibly contain
mutable types, use `@Immutable`’s `containerOf` method:

```
@Immutable(containerOf = "T")
class ImmutableHolder<T> {
 final T ref;
 ...
}
```
Error Prone will allow you to instantiate an `ImmutableHolder<String>` and use
it as a field in another `@Immutable` class. You may instantiate an
`ImmutableHolder<Object>`, but since it is mutable, Error Prone would report an
error if that was a field of another `@Immutable` class.

`@ImmutableTypeParameter`If you want to allow your `@Immutable` generic container to only contain
immutable types, use `@ImmutableTypeParameter`:

```
@Immutable
class ImmutableContainer<@ImmutableTypeParameter T> {
 final T ref;
 ...
}
```
Error Prone will allow you to instantiate a `ImmutableContainer<String>` and use
it as a field in another `@Immutable` class. However, it is a compiler error to
instantiate an `ImmutableContainer<Object>`.

`@ImmutableTypeParameter` can restrict generic parameters to immutable types
only but the generic itself may not be immutable.

```
class MutableContainer<@ImmutableTypeParameter T> {
 ...
}
```
You can also use `@ImmutableTypeParameter` to annotate a method’s type
parameters:

```
class SomeMutableClass {
 <@ImmutableTypeParameter T> ImmutableList<T> putInImmutableList(T t) {
 return ImmutableList.of(t);
 }
}
```
If your `@Immutable` class has a type parameter that is not used in the type of
your class’s fields, then there is no need to use `containerOf` or
`@ImmutableTypeParameter`:

```
@Immutable
class NonContainer<T> {
 ... // No fields whose type contains T
 void process(T element) {
 // process 'element', which won't violate NonContainer's immutability.
 }
}
```
Suppress false positives by adding an `@SuppressWarnings("Immutable")`
annotation to the enclosing element, or the offending field.

To suppress warnings in AutoValue classes, add `@AutoValue.CopyAnnotations` to
ensure the suppression is also applied to the generated sub-class:

```
@AutoValue
@AutoValue.CopyAnnotations
@Immutable
@SuppressWarnings("Immutable")
class MyAutoValue {
 ...
}
```

*Alternate names: ProtoFieldNullComparison*

This checker looks for comparisons of protocol buffer fields with null. If a proto field is not specified, its field accessor will return a non-null default value. Thus, the result of calling one of these accessors can never be null, and comparisons like these often indicate a nearby error.

If you need to distinguish between an unset optional value and a default value,
you have two options. In most cases, you can simply use the `hasField()` method.
proto3 however does not generate `hasField()` methods for primitive types
(including `string` and `bytes`). In those cases you will need to wrap your
field in `google.protobuf.StringValue` or similar.

NOTE: This check applies to normal (server) protos and Lite protos. The
deprecated nano runtime does produce objects which use `null` values to indicate
field absence.

```
void test(MyProto proto) {
 if (proto.getField() == null) {
 ...
 }
 if (proto.getRepeatedFieldList() != null) {
 ...
 }
 if (proto.getRepeatedField(1) != null) {
 ...
 }
}
```
```
void test(MyProto proto) {
 if (!proto.hasField()) {
 ...
 }
 if (!proto.getRepeatedFieldList().isEmpty()) {
 ...
 }
 if (proto.getRepeatedFieldCount() > 1) {
 ...
 }
}
```
If the presence of a field is required information in proto3, the field can be wrapped. For example,

```
message MyMessage {
 google.protobuf.StringValue my_string = 1;
}
```
Presence can then be tested using `myMessage.hasMyString()`, and the value
retrieved using `myMessage.getMyString().getValue()`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ImpossibleNullComparison")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("Incomparable") to the enclosing element.

@SuppressWarnings("Incomparable")

The called method is annotated with `@CompatibleWith`, which enforces that the
parameter passed to the method can potentially be cast to the appropriate
generic type. However, the type of the parameter passed can’t be cast to the
appropriate generic type.

This is useful when a method can’t just take a parameter of the generic type to
allow developers to safely operate with instances held with a wildcard type when
using an instance as both a *consumer* and *producer* of values. This should
*not* be the default, as most interfaces are either one or the other. Containers
and container-like class are the most likely places to use this tool.

TIP: More explanation can be found on the page for CollectionIncompatibleType

```
interface Container<T> {
 void add(T thing);
 boolean contains(@CompatibleWith("T") Object thing);
 boolean containsAsT(T thing);
}
void containmentCheck(Container<? extends Number> container) {
 container.contains(2); // OK, int can be cast to Number
 container.contains(2.0); // OK, double can be cast to Number
 container.contains("a"); // Not OK, String can't be cast to number
 // Does not compile, since, for example, container might be Container<Double>, and Integer
 // can't be cast to Double.
 container.containsAsT(2);
 Container<String> stringContainer = ...;
 stringContainer.contains("a"); // OK
 // OK, since Object *could* be cast to String
 stringContainer.contains(new Object() {});
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("IncompatibleArgumentType")` to the enclosing element.

The @IncompatibleModifiers annotation declares that the target annotation is incompatible with a set of provided modifiers. This check ensures that all annotations respect their @IncompatibleModifiers specifications.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("IncompatibleModifiers")` to the enclosing element.

#
 IndexOfChar

 The first argument to indexOf is a Unicode code point, and the second is the index to start the search from

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("IndexOfChar")` to the enclosing element.

When calling a varargs method, you can either pass an explicit array of arguments, or individual arguments:

```
void f(Object... xs) {
 System.err.println(Arrays.deepToString(xs));
}
```
Both of the following print `[1, 2]`:

```
f(new Object[] {1, 2}) // prints "[1, 2]"
f(1, 2) // prints "[1, 2]"
```
If the argument to the varargs method is a conditional expression, and either branch is not an array, the result of the expression will be implicitly wrapped in an array.

```
f(flag ? 1 : 2) // prints [1] or [2]
```
This means that if one branch is an array and the other branch is not, the array branch will become a multi-dimensional array:

```
f(flag ? new Object[] {1, 2} : 3); // prints [[1, 2]] or [3]
```
To avoid the implicit array creation, the other argument can be explicitly wrapped in an array:

```
f(flag ? new Object[] {1, 2} : new Object[] {3}); // prints [1, 2] or [3]
```
Or, if the multi-dimensional array was intentional, it can be written explicitly as:

```
f(flag ? new Object[][] 1 : new Object[] {3}); // prints [[1, 2]] or [3]
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("InexactVarargsConditional")` to the enclosing element.

A method that always calls itself will cause a StackOverflowError.

```
int oops() {
 return oops();
}
```
```
Exception in thread "main" java.lang.StackOverflowError
 at Test.oops(X.java:3)
 at Test.oops(X.java:3)
 ...
```
The fix may be to call another method with the same name:

```
void process(String name, int id) {
 process(name, id); // error
 process(name, id, /*verbose=*/ true); // ok
}
void process(String name, int id, boolean verbose) {
 // ...
}
```
or to call the method on a different instance:

```
class Delegate implements Processor {
 Processor delegate;
 void process(String name, int id) {
 process(name, id); // error
 delegate.process(name, id); // ok
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("InfiniteRecursion")` to the enclosing element.

*Alternate names: MoreThanOneScopeAnnotationOnClass*

Annotating a class with more than one scope annotation is invalid according to the JSR-330 specification.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InjectMoreThanOneScopeAnnotationOnClass")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("InjectOnMemberAndConstructor") to the enclosing element.

@SuppressWarnings("InjectOnMemberAndConstructor")

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
InlineMeValidator
Ensures that the @InlineMe annotation is used correctly.
Severity
ERROR

Suppress false positives by adding the suppression annotation @SuppressWarnings("InstantTemporalUnit") to the enclosing element.

@SuppressWarnings("InstantTemporalUnit")

#
 InvalidJavaTimeConstant

 This checker errors on calls to java.time methods using values that are guaranteed to throw a DateTimeException.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InvalidJavaTimeConstant")` to the enclosing element.

This error is triggered by calls to regex-accepting methods with invalid string literals. These calls would cause a PatternSyntaxException at runtime.

We deliberately do not check `java.util.regex.Pattern#compile` as many of its
users are deliberately testing the regex compiler or using a vacuously true
regex.

`"."` is also discouraged, as it is a valid regex but is easy to mistake for
`"\\."`. Instead of e.g. `str.replaceAll(".", "x")`, prefer ```
Strings.repeat("x",
str.length())
```
 or `CharMatcher.ANY.replaceFrom(str, "x")`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InvalidPatternSyntax")` to the enclosing element.

TimeZone.getTimeZone(String) silently returns GMT when an invalid time zone identifier is passed in.

Suppress false positives by adding the suppression annotation @SuppressWarnings("InvalidTimeZoneID") to the enclosing element.

@SuppressWarnings("InvalidTimeZoneID")

Suppress false positives by adding the suppression annotation @SuppressWarnings("InvalidZoneId") to the enclosing element.

@SuppressWarnings("InvalidZoneId")

Suppress false positives by adding the suppression annotation @SuppressWarnings("IsInstanceIncompatibleType") to the enclosing element.

@SuppressWarnings("IsInstanceIncompatibleType")

Passing an argument of type `Class` to `Class#instanceOf(Class)` is usually a
mistake.

Calling `clazz.instanceOf(obj)` for some `clazz` with type `Class<T>` is
equivalent to `obj instanceof T`. The `instanceOf` method exists for cases where
the type `T` is not known statically.

When a class literal is passed as the argument of `instanceOf`, the result will
only be true if the class literal on left hand side is equal to `Class.class`.

For example, the following code returns true if and only if the type `A` is
equal to `Class` (i.e. lhs is equal to `Class.class`).

```
<A, B> boolean f(Class<A> lhs, Class<B> rhs) {
 return lhs.instanceOf(rhs); // equivalent to 'lhs == Class.class'
}
```
To test if the type represented by a class literal is a subtype of the type
represented by some other class literal, `isAssignableFrom` should be used
instead:

```
<A, B> boolean f(Class<A> lhs, Class<B> rhs) {
 return lhs.isAssignableFrom(rhs); // equivalent to 'B instanceof A'
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("IsInstanceOfClass")` to the enclosing element.

Log.isLoggable(tag, level) throws an IllegalArgumentException if its tag argument is more than 23 characters long.

Log.isLoggable(tag, level)

IllegalArgumentException

Suppress false positives by adding the suppression annotation @SuppressWarnings("IsLoggableTagLength") to the enclosing element.

@SuppressWarnings("IsLoggableTagLength")

JUnit 3 requires that test method names start with “`test`”. The method that
triggered this error looks like it is supposed to be a test, but misspells the
required prefix; has `@Test` annotation, but no prefix; or has the wrong method
signature. As a consequence, JUnit 3 will ignore it.

If you meant to disable this test on purpose, or this is a helper method, change
the name to something more descriptive, like “`disabledTestSomething()`”. You
don’t need an `@Test` annotation, but if you want to keep it, add `@Ignore` too.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JUnit3TestNotRun")` to the enclosing element.

JUnit4 provides two annotations (`@BeforeClass` and
`@AfterClass`) that are applied to methods that are run once per
**test class**. These complement the more-often used `@Before` and `@After`
which are applied to methods that are run once per **test method**.

JUnit4 runs `@BeforeClass` and `@AfterClass` methods without making an instance
of the test class, meaning that the methods must be `static`. JUnit4 will fail
to run any `@BeforeClass` or `@AfterClass` method that isn’t also `static`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JUnit4ClassAnnotationNonStatic")` to the enclosing element.

JUnit 3 provides the method setUp(), to be overridden by subclasses when the test needs to perform some pre-test initialization. In JUnit 4, this is accomplished by annotating such a method with @Before.

The method that triggered this error matches the definition of setUp() from JUnit3, but was not annotated with @Before and thus won’t be run by the JUnit4 runner.

If you intend for this setUp() method not to run by the JUnit4 runner, but perhaps manually be invoked in certain test methods, please rename the method or mark it private.

If the method is part of an abstract test class hierarchy where this class’s setUp() is invoked by a superclass method that is annotated with @Before, then please rename the abstract method or add @Before to the superclass’s definition of setUp()

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JUnit4SetUpNotRun")` to the enclosing element.

JUnit 3 provides the overridable method tearDown(), to be overridden by subclasses when the test needs to perform some post-test de-initialization. In JUnit 4, this is accomplished by annotating such a method with @After. The method that triggered this error matches the definition of tearDown() from JUnit3, but was not annotated with @After and thus won’t be run by the JUnit4 runner.

If you intend for this tearDown() method not to be run by the JUnit4 runner, but perhaps be manually invoked after certain test methods, please rename the method or mark it private.

If the method is part of an abstract test class hierarchy where this class’s tearDown() is invoked by a superclass method that is annotated with @After, then please rename the abstract method or add @After to the superclass’s definition of tearDown().

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JUnit4TearDownNotRun")` to the enclosing element.

Unlike in JUnit 3, JUnit 4 tests will not be run unless annotated with @Test. The test method that triggered this error looks like it was meant to be a test, but was not so annotated, so it will not be run. If you intend for this test method not to run, please add both an @Test and an @Ignore annotation to make it clear that you are purposely disabling it. If this is a helper method and not a test, consider reducing its visibility to non-public, if possible.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JUnit4TestNotRun")` to the enclosing element.

#
 JUnit4TestsNotRunWithinEnclosed

 This test is annotated @Test, but given it's within a class using the Enclosed runner, will not run.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JUnit4TestsNotRunWithinEnclosed")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("JUnitAssertSameCheck") to the enclosing element.

@SuppressWarnings("JUnitAssertSameCheck")

Suppress false positives by adding the suppression annotation @SuppressWarnings("JUnitParameterMethodNotFound") to the enclosing element.

@SuppressWarnings("JUnitParameterMethodNotFound")

The `Inject` annotation cannot be applied to abstract methods, per the JSR-330
spec, since injectors will only inject those methods if the concrete implementer
of the abstract method has the `Inject` annotation as well. See
OverridesJavaxInjectableMethod for more examples of this interaction.

Currently, default methods in interfaces are not injected if they have
`Inject` for similar reasons, although future updates to dependency injection
frameworks may allow this, since the default methods are not abstract.

See the Guice wiki page on JSR-330 for more.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JavaxInjectOnAbstractMethod")` to the enclosing element.

Joda-Time’s DateTime.toDateTime(), Duration.toDuration(), Instant.toInstant(), Interval.toInterval(), and Period.toPeriod() are always unnecessary, since they simply ‘return this’. There is no reason to ever call them.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JodaToSelf")` to the enclosing element.

Labels should only be used on loops. For other statements, consider refactoring to express the control flow without labelled breaks.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("LabelledBreakTarget")` to the enclosing element.

This check ensures that the number of arguments passed to ‘lenient’ formatting
methods like `Preconditions.checkArgument` match the number of format
specifiers.

WARNING: Only the exact two-character placeholder sequence `%s` is recognized by
these methods. Any others will be ignored, and not used for argument
substitution.

The APIs checked by this bugpattern include:

`com.google.common.base.Strings#lenientFormat``com.google.common.base.Preconditions#check*``com.google.common.base.Verify#verify*``com.google.common.truth.Truth#assertWithMessage``com.google.common.truth.Subject#check``com.google.common.truth.StandardSubjectBuilder#withMessage`Suppress false positives by adding the suppression annotation `@SuppressWarnings("LenientFormatStringValidation")` to the enclosing element.

When serializing bytes from a `MessageLite`, one can use `toByteString` to get a
`ByteString`, effectively an immutable wrapper over a `byte[]`. This ByteString
can be passed around and deserialized into a message using
`MyMessage.Builder.mergeFrom(ByteString)`.

`ByteString#toStringUtf8` copies UTF-8 encoded byte data living inside the
`ByteString` to a `java.lang.String`, replacing any
invalid UTF-8 byte sequences with � (the Unicode
replacement character).

In this circumstance, a protocol message is being serialized to a `ByteString`,
then immediately turned into a Java `String` using the `toStringUtf8` method.
However, serialized protocol buffers are arbitrary binary data and not
UTF-8-encoded data. Thus, the resulting `String` may not match the actual
serialized bytes from the protocol message.

Instead of holding the serialized protocol message in a Java `String`, carry
around the actual bytes in a `ByteString`, `byte[]`, or some other equivalent
container for arbitrary binary data.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("LiteByteStringUtf8")` to the enclosing element.

#
 LocalDateTemporalAmount

 LocalDate.plus() and minus() does not work with Durations. LocalDate represents civil time (years/months/days), so java.time.Period is the appropriate thing to add or subtract instead.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("LocalDateTemporalAmount")` to the enclosing element.

Instances of boxed primitive types may be cached by the standard library
`valueOf` method. This method is used for autoboxing. This means that using a
boxed primitive as a lock can result in unintentionally sharing a lock with
another piece of code.

Consider using an explicit lock `Object` instead of locking on a boxed
primitive. That is, prefer this:

```
private final Object lock = new Object();
void doSomething() {
 synchronized (lock) {
 // ...
 }
}
```
instead of this:

```
private final Integer lock = 42;
void doSomething() {
 synchronized (lock) {
 // ...
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("LockOnBoxedPrimitive")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("LoopConditionChecker") to the enclosing element.

@SuppressWarnings("LoopConditionChecker")

Implicit widening conversions when comparing two primitives with methods like Float.compare can lead to lossy comparison. For example, `Float.compare(Integer.MAX_VALUE, Integer.MAX_VALUE - 1) == 0`. Use a compare method with non-lossy conversion, or ideally no conversion if possible.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("LossyPrimitiveCompare")` to the enclosing element.

Math.round() called with an integer or long type results in truncation because Math.round only accepts floats or doubles and some integers and longs can’t be represented with float.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MathRoundIntLong")` to the enclosing element.

`MemorySegment` is a value-based class. Comparing `MemorySegment` instances
using `==` or `!=` is bug-prone because `MemorySegment` implementations may not
be unique for the same underlying memory. Use `Objects.equals()` (or `.equals()`
if the receiver is known to be non-null) instead.

For example:

```
MemorySegment seg = ...;
if (seg == MemorySegment.NULL) { // reference equality
 ...
}
```
should be:

```
MemorySegment seg = ...;
if (Objects.equals(seg, MemorySegment.NULL)) { // value equality
 ...
}
```
Reference equality between any two `MemorySegment` instances (e.g., `a == b`) is
flagged by this check, and `Objects.equals(a, b)` is the preferred alternative.

See also https://bugs.openjdk.org/browse/JDK-8381012

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MemorySegmentReferenceEquality")` to the enclosing element.

Certain resources in `android.R.string` have names that do not match their
content: `android.R.string.yes` is actually “OK” and `android.R.string.no` is
“Cancel”. Avoid these string resources and prefer ones whose names *do* match
their content. If you need “Yes” or “No” you must create your own string
resources.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MislabeledAndroidString")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("MisleadingEmptyVarargs") to the enclosing element.

@SuppressWarnings("MisleadingEmptyVarargs")

When Java introduced text blocks as a feature, it also introduced a new string
escape sequence `\s`. This escape sequence is another way to write a normal
space, but it has the advantage that it can be used at the end of a line in a
text block, where a normal space would be stripped.

This new escape sequence can easily be confused with the regex `\s`, which is a
metacharacter that matches any kind of whitespace character. To write that
metacharacter in a Java string, you must still write `\\s`: an escaped backslash
followed by an `s`.

There is little reason to ever write the Java escape `\s` except at the end of a
line. Either use a normal space, or switch to `\\s` if you are trying to write
the regex metacharacter.

```
// Each line here is five characters long.
String colors = """
 one \s
 two \s
 three
 """;
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("MisleadingEscapedSpace")` to the enclosing element.

#
 MisplacedScopeAnnotations

 Scope annotations used as qualifier annotations don't have any effect. Move the scope annotation to the binding location or delete it.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MisplacedScopeAnnotations")` to the enclosing element.

API providers may annotate a method with an annotation like
`androidx.annotation.CallSuper` or
`javax.annotation.OverridingMethodsMustInvokeSuper` to require that overriding
methods invoke the super method. This check enforces those annotations.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MissingSuperCall")` to the enclosing element.

Some test helpers such as `EqualsTester` require a terminating method call to be
of any use.

```
 @Test
 public void string() {
 new EqualsTester()
 .addEqualityGroup("hello", new String("hello"))
 .addEqualityGroup("world", new String("world"))
 .addEqualityGroup(2, Integer.valueOf(2));
 // Oops: forgot to call `testEquals()`
 }
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("MissingTestCall")` to the enclosing element.

#
 MisusedDayOfYear

 Use of 'DD' (day of year) in a date pattern with 'MM' (month of year) is not likely to be intentional, as it would lead to dates like 'March 73rd'.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MisusedDayOfYear")` to the enclosing element.

“YYYY” in a date pattern means “week year”. The week year is defined to begin at the beginning of the week that contains the year’s first Thursday. For example, the week year 2015 began on Monday, December 29, 2014, since January 1, 2015, was on a Thursday.

“Week year” is intended to be used for week dates, e.g. “2015-W01-1”, but is often mistakenly used for calendar dates, e.g. 2014-12-29, in which case the year may be incorrect during the last week of the year. If you are formatting anything other than a week date, you should use the year specifier “yyyy” instead.

This isn’t an idle risk; Twitter had a
significant outage
in ~~2015~~ 2014
due to this bug.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MisusedWeekYear")` to the enclosing element.

A field `Descriptor` was created by mixing the message `Descriptor` from one
proto message with the field number from another. For example:

```
Foo.getDescriptor().findFieldByNumber(Bar.ID_FIELD_NUMBER)
```
This accesses the `Descriptor` of a field in `Foo` with a field number from
`Bar`. One of these was probably intended:

```
Foo.getDescriptor().findFieldByNumber(Foo.ID_FIELD_NUMBER)
Bar.getDescriptor().findFieldByNumber(Bar.ID_FIELD_NUMBER)
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("MixedDescriptors")` to the enclosing element.

Calls to `Mockito.when` should always be accompanied by a call to a method like
`thenReturn`.

```
when(mock.get()).thenReturn(answer); // correct
when(mock.get()) // oops!
```
Similarly, calls to `Mockito.verify` should call the verified method *outside*
the call to `verify`.

```
verify(mock).execute(); // correct
verify(mock.execute()); // oops!
```
For more information, see the Mockito documentation.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MockitoUsage")` to the enclosing element.

Invoking a collection method with the same collection as the argument is likely incorrect.

`collection.addAll(collection)` may cause an infinite loop, duplicate the
elements, or do nothing, depending on the type of Collection and
implementation class.`collection.retainAll(collection)` is a no-op.`collection.removeAll(collection)` is the same as `collection.clear()`.`collection.containsAll(collection)` is always true.Suppress false positives by adding the suppression annotation `@SuppressWarnings("ModifyingCollectionWithItself")` to the enclosing element.

*Alternate names: inject-constructors, InjectMultipleAtInjectConstructors*

Injection frameworks may use `@Inject` to determine how to construct an object
in the absence of other instructions. Annotating `@Inject` on a constructor
tells the injection framework to use that constructor. However, if multiple
`@Inject` constructors exist, injection frameworks can’t reliably choose between
them.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MoreThanOneInjectableConstructor")` to the enclosing element.

*Alternate names: MustBeClosed*

Methods or constructors annotated with `@MustBeClosed` require that the returned
resource is closed. This is enforced by checking that invocations occur within
the resource variable initializer of a try-with-resources statement:

```
try (AutoCloseable resource = createTheResource()) {
 doSomething(resource);
}
```
or the `return` statement of another method annotated with `@MustBeClosed`:

```
@MustBeClosed
AutoCloseable createMyResource() {
 return createTheResource();
}
```
To support legacy code, the following pattern is also supported:

```
AutoCloseable resource = createTheResource();
try {
 doSomething(resource);
} finally {
 resource.close();
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("MustBeClosedChecker")` to the enclosing element.

Calling nCopies('a', 10) returns a list with 97 copies of 10, not a list with 10 copies of ‘a’.

nCopies('a', 10)

Suppress false positives by adding the suppression annotation @SuppressWarnings("NCopiesOfChar") to the enclosing element.

@SuppressWarnings("NCopiesOfChar")

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
NoCanIgnoreReturnValueOnClasses
@CanIgnoreReturnValue should not be applied to classes as it almost always overmatches (as it applies to constructors and all methods), and the CIRVness isn't conferred to its subclasses.
Severity
ERROR

Types should always be imported by their canonical name. The canonical name of a top-level class is the fully-qualified name of the package, followed by a ‘.’, followed by the name of the class. The canonical name of a member class is the canonical name of its declaring class, followed by a ‘.’, followed by the name of the member class.

Fully-qualified member class names are not guaranteed to be canonical. Consider some member class M declared in a class C. There may be another class D that extends C and inherits M. Therefore M can be accessed using the fully-qualified name of D, followed by a ‘.’, followed by ‘M’. Since M is not declared in D, this name is not canonical.

The JLS §7.5.3 requires all single static imports to *start* with a canonical
type name, but the fully-qualified name of the imported member is not required
to be canonical.

Importing types using non-canonical names is unnecessary and unclear, and should be avoided.

Example:

```
package a;
class One {
 static class Inner {}
}
```
```
package a;
class Two extends One {}
```
An import of `Inner` should always refer to it using the canonical name
`a.One.Inner`, not `a.Two.Inner`.

If a method’s formal parameter is annotated with `@CompileTimeConstant`, the
method will always be invoked with an argument that is a static constant. If the
parameter itself is non-final, then it is a mutable reference to immutable data.
This is rarely useful, and can be confusing when trying to use the parameter in
a context that requires an compile-time constant. For example:

```
void f(@CompileTimeConstant y) {}
void g(@CompileTimeConstant x) {
 x = f(x); // x is not a constant
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("NonFinalCompileTimeConstant")` to the enclosing element.

Calling getAnnotation on an annotation that does not have its Retention set to RetentionPolicy.RUNTIME will always return null.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NonRuntimeAnnotation")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("NullArgumentForNonNullParameter") to the enclosing element.

@SuppressWarnings("NullArgumentForNonNullParameter")

#
 NullNeedsCastForVarargs

 This call passes a null *array*, so it always produces NullPointerException. To pass a null *element*, cast to the element type.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NullNeedsCastForVarargs")` to the enclosing element.

If a conditional expression evaluates to `null`, unboxing it will result in a
`NullPointerException`.

For example:

```
int x = flag ? foo : null;
```
If `flag` is false, `null` will be auto-unboxed from an `Integer` to `int`,
resulting in a NullPointerException.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NullTernary")` to the enclosing element.

The correct syntax to apply a `TYPE_USE` annotation to an inner class is
`A.@Nullable B`.

For a `TYPE_USE` `@Nullable` annotation, `@Nullable A.B` is legal Java if `B` is
a non-static inner class:

```
class A {
 @Target(TYPE_USE)
 @interface Nullable {}
 class B {}
 static class C {}
 void test(A.@Nullable B x) {} // B is annotated ('A' is the enclosing instance type)
 void test(A.@Nullable C x) {} // C is annotated ('A' is a 'scoping construct' here)
}
```
```
 void test(@Nullable A.B x) {} // compiles, but likely incorrect: annotates the enclosing instance type 'A', which can never be null
 void test(@Nullable A.C x) {} // compile error: 'A' cannot be annotated
```
However, for `@Nullable` (and `@NonNull`, and friends), annotating the outer
class is meaningless. The reference to the outer class (`A.this`) can never be
`null`, so any nullability annotations are redundant.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NullableOnContainingClass")` to the enclosing element.

Optionals should be compared for value equality using .equals(), and not for reference equality using == and !=.

.equals()

==

!=

Suppress false positives by adding the suppression annotation @SuppressWarnings("OptionalEquality") to the enclosing element.

@SuppressWarnings("OptionalEquality")

Suppress false positives by adding the suppression annotation @SuppressWarnings("OptionalMapUnusedValue") to the enclosing element.

@SuppressWarnings("OptionalMapUnusedValue")

#
 OptionalOfRedundantMethod

 Optional.of() always returns a non-empty optional. Using ifPresent/isPresent/orElse/orElseGet/orElseThrow/isPresent/or/orNull method on it is unnecessary and most probably a bug.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("OptionalOfRedundantMethod")` to the enclosing element.

Qualifiers and Scoping annotations have different semantic meanings and a single annotation should not be both a qualifier and a scoping annotation.

If an annotation is both a scoping annotation and a qualifier, unless great care is taken with its application and usage, the semantics of objects annotated with the annotation are unclear.

Take a look at this example:

```
@Retention(RetentionPolicy.RUNTIME)
@Scope
@Qualifier
@interface DayScoped {}
static class Allowance {}
static class DailyAllowance extends Allowance {}
static class Spender {
 @Inject
 Spender(Allowance allowance) {}
}
static class BindingModule extends AbstractModule {
 ...
 @Provides
 @DayScoped
 Allowance providesAllowance() {
 return new DailyAllowance();
 }
}
```
Here, the `Allowance` instance used by Spender isn’t actually scoped to a single
day, as the `@Provides` method applies the `DayScoped` scoping only to the
`@DayScoped Allowance`. Instead, the default constructor of `Allowance` is used
to create a new instance every time a `Spender` is created.

If `@DayScope` wasn’t a `Qualifier`, the provider method would do the right
thing: the un-annotated `Allowance` binding would be scoped to `DayScope`,
implemented by a single `DailyAllowance` instance per day.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("OverlappingQualifierAndScopeAnnotation")` to the enclosing element.

When classes declare that they have an `@javax.inject.Inject`ed method,
dependency injection tools must call those methods after first calling any
`@javax.inject.Inject` constructor, and performing any field injection. These
methods are part of the initialization contract for the object.

When subclasses override methods annotated with `@javax.inject.Inject` and
*don’t* also annotate themselves with `@javax.inject.Inject`, the injector will
not call those methods as part of the subclass’s initialization. This may
unexpectedly cause assumptions taken in the superclass (e.g.: this
post-initialization routine is finished, meaning that I can safely use this
field) to no longer hold.

This compile error is intended to prevent this unintentional breaking of assumptions. Possible resolutions to this error include:

`@Inject` the overridden method, calling the `super` method to maintain the
initialization contract.`final` to avoid subclasses unintentionally
masking the injected method.`protected` method for
this subclass to use.Suppress false positives by adding the suppression annotation `@SuppressWarnings("OverridesJavaxInjectableMethod")` to the enclosing element.

Classes should not be declared inside `package-info.java` files.

Typically package-info.java contains only a package declaration, preceded immediately by the annotations on the package. While the file could technically contain the source code for one or more classes with package access, it would be very bad form.

– JLS 7.4

Suppress false positives by adding the suppression annotation `@SuppressWarnings("PackageInfo")` to the enclosing element.

#
 ParametersButNotParameterized

 This test has @Parameters but is using the default JUnit4 runner. The parameters will have no effect.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ParametersButNotParameterized")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("ParcelableCreator") to the enclosing element.

@SuppressWarnings("ParcelableCreator")

Period.from(TemporalAmount) will always throw a DateTimeException when passed a Duration and return itself when passed a Period.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("PeriodFrom")` to the enclosing element.

`Period.get(TemporalUnit)` only works when passed `ChronoUnit.YEARS`, `ChronoUnit.MONTHS`, or `ChronoUnit.DAYS`. All other values are guaranteed to throw an `UnsupportedTemporalTypeException`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("PeriodGetTemporalUnit")` to the enclosing element.

Period.(plus|minus)(TemporalAmount) will always throw a DateTimeException when passed a Duration.

Suppress false positives by adding the suppression annotation @SuppressWarnings("PeriodTimeMath") to the enclosing element.

@SuppressWarnings("PeriodTimeMath")

The Guava Preconditions checks take error message template strings that look similar to format strings, but only accept the %s format (not %d, %f, etc.). This check points out places where a Preconditions error message template string has a non-%s format, or where the number of arguments does not match the number of %s formats in the string.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("PreconditionsInvalidPlaceholder")` to the enclosing element.

#
 PrivateSecurityContractProtoAccess

 Access to a private protocol buffer field is forbidden. This protocol buffer carries a security contract, and can only be created using an approved library. Direct access to the fields is forbidden.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("PrivateSecurityContractProtoAccess")` to the enclosing element.

#
 ProtoBuilderReturnValueIgnored

 Unnecessary call to proto's #build() method. If you don't consume the return value of #build(), the result is discarded and the only effect is to verify that all required fields are set, which can be expressed more directly with #isInitialized().

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ProtoBuilderReturnValueIgnored")` to the enclosing element.

Comparing strings with == is almost always an error, but it is an error 100% of the time when one of the strings is a protobuf field. Additionally, protobuf fields cannot be null, so Object.equals(Object) is always more correct.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ProtoStringFieldReferenceEquality")` to the enclosing element.

ProtoTruth’s `#ignoringFields` method accepts integer field numbers, so
supplying field numbers from the wrong protocol buffers is possible. For
example:

```
message Bar {
 optional string name = 1;
}
message Foo {
 optional string name = 1;
 optional Bar bar = 2;
}
```
```
void assertOnFoo(Foo foo) {
 assertThat(foo).ignoringFields(Bar.NAME_FIELD_NUMBER).isEqualTo(...);
}
```
This will ignore the `Foo#name` field rather than `Bar#name`. The field number
can be turned into a `Descriptor` object to resolve the correct nested field to
ignore:

```
void assertOnFoo(Foo foo) {
 assertThat(foo)
 .ignoringFieldDescriptors(
 Bar.getDescriptor().findFieldByNumber(Bar.NAME_FIELD_NUMBER))
 .isEqualTo(...);
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ProtoTruthMixedDescriptors")` to the enclosing element.

The generated Java source files for Protocol Buffer enums have `getNumber()` as
accessors for the tag number in the protobuf file.

In addition, since it’s a java enum, it also has the `ordinal()` method,
returning its positional index within the generated java enum.

The `ordinal()` order of the generated Java enums isn’t guaranteed, and can
change when a new enum value is inserted into a proto enum. The `getNumber()`
value won’t change for an enum value (since making that change is a
backwards-incompatible change for the protocol buffer).

You should very likely use `getNumber()` in preference to `ordinal()` in all
circumstances since it’s a more stable value.

Note: If you’re changing code that was already using ordinal(), it’s likely that getNumber() will return a different real value. Tread carefully to avoid mismatches if the ordinal was persisted elsewhere.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ProtocolBufferOrdinal")` to the enclosing element.

Guice `@Provides` methods annotate methods that are used as a means of declaring
bindings. However, this is only helpful inside of a module. Methods outside of
these modules are not used for binding declaration.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ProvidesMethodOutsideOfModule")` to the enclosing element.

`Math.random()`, `Random#nextFloat`, and `Random#nextDouble` return results in
the range `[0.0, 1.0)`. Therefore, casting the result to `(int)` or `(long)`
*always* results in the value of `0`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("RandomCast")` to the enclosing element.

`Random.nextInt() % n` has

`1` to `n-1` inclusive`-1` to `-(n-1)` inclusiveMany users expect a uniformly distributed random integer between `0` and `n-1`
inclusive, but you must use random.nextInt(n) to get that behavior. If the
original behavior is truly desired, use ```
(random.nextBoolean() ? 1 : -1) *
random.nextInt(n)
```
.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("RandomModInteger")` to the enclosing element.

In a Java `record`, using the accessor method (like `d()`) inside a **compact
constructor** (the one without arguments, `Foo { ... }`) reads the record’s
underlying field before it has been set. This means the method will always
return `null` (for objects) or `0`/`false` (for primitives), regardless of what
arguments were passed to the constructor.

```
record User(String name) {
 User {
 // BUG: name() reads the uninitialized field 'this.name', which is currently
 // null. This throws a NullPointerException immediately.
 if (name().isEmpty()) {
 throw new IllegalArgumentException("Name cannot be empty");
 }
 }
}
```
```
record User(String name) {
 User {
 // CORRECT: Reads the constructor parameter 'name'.
 if (name.isEmpty()) {
 throw new IllegalArgumentException("Name cannot be empty");
 }
 }
}
```
The **Compact Constructor** in Java is a special initialization block that runs
*before* the record’s fields are automatically assigned.

`User { ... }`
runs first.`this.name = name;`) only When you call `name()`, it attempts to read `this.name`. Since the assignment
hasn’t happened yet, `this.name` still holds its default value, which is `null`
(JLS §4.12.5).

To fix this, refer to the component by its name (e.g., `name`). This accesses
the **parameter** passed to the constructor, which holds the correct value.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("RecordAccessorInCompactConstructor")` to the enclosing element.

`android.graphics.Rect.intersect(Rect r)` and
`android.graphics.Rect.intersect(int, int, int, int)` do not always modify the
rectangle to the intersected result. If the rectangles do not intersect, no
change is made and the original rectangle is not modified. These methods return
false to indicate that this has happened.

If you don’t check the return value of these methods, you may end up drawing the wrong rectangle.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("RectIntersectReturnValueIgnored")` to the enclosing element.

*Alternate names: ProtoRedundantSet*

Proto and AutoValue builders provide a fluent interface for constructing instances. Unlike argument lists, however, they do not prevent the user from providing multiple values for the same field.

Setting the same field multiple times in the same chained expression is
pointless (as the intermediate value will be overwritten), and can easily mask a
bug, especially if the setter is called with *different* arguments.

```
return MyProto.newBuilder()
 .setFoo(copy.getFoo())
 .setFoo(copy.getBar())
 .build();
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("RedundantSetterCall")` to the enclosing element.

This annotation is itself annotated with @RequiredModifiers and can only be used when the specified modifiers are present. You are attempting to use it on an element that is missing one or more required modifiers.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("RequiredModifiers")` to the enclosing element.

Calls to APIs marked `@RestrictedApi` are prohibited without a corresponding
allowlist annotation.

The intended use-case for `@RestrictedApi` is to restrict calls to annotated
methods so that each usage of those APIs must be reviewed separately. For
example, an API might lead to security bugs unless the programmer uses it
correctly.

See the
javadoc for `@RestrictedApi`
for more details.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("RestrictedApi")` to the enclosing element.

*Alternate names: ResultOfMethodCallIgnored, CheckReturnValue*

Certain library methods do nothing useful if their return value is ignored. For example, String.trim() has no side effects, and you must store the return value of String.intern() to access the interned string. This check encodes a list of methods in the JDK whose return value must be used and issues an error if they are not.

`Optional.orElseThrow`Don’t call `orElseThrow` just for its side-effects. When the result of a call to
`orElseThrow` is discarded, it may be unclear to future readers whether the
result is being discarded deliberately or accidentally. That is, avoid:

```
// return value of orElseThrow() is silently ignored here
optional.orElseThrow(() -> new AssertionError("something has gone terribly wrong"));
```
Instead of calling `orElseThrow` for its side-effects, prefer an explicit call
to `isPresent()`, or use one of `checkState` or `checkArgument` from
Guava’s `Preconditions` class.

```
if (!optional.isPresent()) {
 throw new AssertionError("something has gone terribly wrong");
}
```
```
checkState(optional.isPresent(), "something has gone terribly wrong");
```
```
checkArgument(optional.isPresent(), "something has gone terribly wrong");
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ReturnValueIgnored")` to the enclosing element.

*Alternate names: TruthSelfEquals*

When using Truth, if a test subject and the argument to `isEqualTo` are the same
instance (for example `assertThat(x).isEqualTo(x)`), then the assertion will
always pass. Truth implements `isEqualTo` using [`Objects#equal`] , which tests
its arguments for reference equality and returns true without calling `equals()`
if both arguments are the same instance.

JUnit’s `assertEquals` (and similar) methods are implemented in terms of
`Object#equals`. However, this is not explicitly documented, so isn’t a
contractual guarantee of the assertion methods.

To test the implementation of an `equals` method, use
Guava’s EqualsTester, or explicitly call `equals` as part of the
test.

In our experience, `assertThat(x).isEqualTo(x)` and similar are *more likely to
be typos* than assertions about an `equals` method. This alone is sufficient
motivation to choose a dedicated approach for testing `equals` implementations.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("SelfAssertion")` to the enclosing element.

The left-hand side and right-hand side of this assignment are the same. It has no effect.

This also handles assignments in which the right-hand side is a call to Preconditions.checkNotNull(), which returns the variable that was checked for non-nullity. If you just intended to check that the variable is non-null, please don’t assign the result to the checked variable; just call Preconditions.checkNotNull() as a bare statement.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("SelfAssignment")` to the enclosing element.

The arguments to compareTo method are the same object, so it always returns 0. Either change the arguments to point to different objects or substitute 0.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("SelfComparison")` to the enclosing element.

The arguments to equals method are the same object, so it always returns true. Either change the arguments to point to different objects or substitute true.

For test cases, instead of explicitly testing equals, use EqualsTester from Guava.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("SelfEquals")` to the enclosing element.

#
 SetUnrecognized

 Setting a proto field to an UNRECOGNIZED value will result in an exception at runtime when building.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("SetUnrecognized")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("ShouldHaveEvenArgs") to the enclosing element.

@SuppressWarnings("ShouldHaveEvenArgs")

A standard means of checking non-emptiness of an array or collection is to test if the size of that collection is greater than 0. However, one may accidentally check if the size is greater than or equal to 0, which is always true.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("SizeGreaterThanOrEqualsZero")` to the enclosing element.

The `toString` method on a `Stream` will print its identity, such as
`java.util.stream.ReferencePipeline$Head@6d06d69c`. This is rarely what was
intended.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("StreamToString")` to the enclosing element.

StringBuilder does not have a char constructor, so instead this code creates a StringBuilder with initial size equal to the code point of the specified char.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("StringBuilderInitWithChar")` to the enclosing element.

#
 StringJoin

 String.join(CharSequence) performs no joining (it always returns the empty string); String.join(CharSequence, CharSequence) performs no joining (it just returns the 2nd parameter).

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("StringJoin")` to the enclosing element.

String.substring(int) gives you the substring from the index to the end, inclusive. Calling that method with an index of 0 will return the same String.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("SubstringOfZero")` to the enclosing element.

To suppress warnings to deprecated methods, you should add the annotation
`@SuppressWarnings("deprecation")` and not `@SuppressWarnings("deprecated")`

Suppress false positives by adding the suppression annotation `@SuppressWarnings("SuppressWarningsDeprecated")` to the enclosing element.

TemporalAccessor.get(ChronoField) only works for certain values of ChronoField. E.g., DayOfWeek only supports DAY_OF_WEEK. All other values are guaranteed to throw an UnsupportedTemporalTypeException.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("TemporalAccessorGetChronoField")` to the enclosing element.

#
 TestParametersNotInitialized

 This test has @TestParameter fields but is using the default JUnit4 runner. The parameters will not be initialised beyond their default value.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("TestParametersNotInitialized")` to the enclosing element.

#
 TheoryButNoTheories

 This test has members annotated with @Theory, @DataPoint, or @DataPoints but is using the default JUnit4 runner.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("TheoryButNoTheories")` to the enclosing element.

#
 ThreadBuilderNameWithPlaceholder

 Thread.Builder.name() does not accept placeholders (e.g., %d or %s). threadBuilder.name(String) accepts a constant name and threadBuilder.name(String, int) accepts a constant name _prefix_ and an initial counter value.

 - Severity
- ERROR
- Tags
- LikelyError

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ThreadBuilderNameWithPlaceholder")` to the enclosing element.

`throwIfUnchecked(knownCheckedException)` is a no-op (aside from performing a
null check). `propagateIfPossible(knownCheckedException)` is a complete no-op.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ThrowIfUncheckedKnownChecked")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("ThrowNull") to the enclosing element.

@SuppressWarnings("ThrowNull")

`Tree#toString` shouldn’t be used for Trees deriving from the code being
compiled, as it discards whitespace and comments.

This check only runs inside Error Prone code. Suggested replacements include:

`VisitorState#getConstantExpression` for escaping constants in
generated code.`VisitorState#getSourceForNode` : it will give you the original source text.
Note that for synthetic trees (e.g.: implicit constructors), that source may
be `null`.`this` and `super`, try `tree.getName().contentEquals("this")``ASTHelpers.getSymbol(tree).getSimpleName().toString()`Suppress false positives by adding the suppression annotation `@SuppressWarnings("TreeToString")` to the enclosing element.

When testing that a line of code throws an expected exception, it is typical to
execute that line in a try block with a `fail()` or `assert*()` on the line
following. The expectation is that the expected exception will be thrown, and
execution will continue in the catch block, and the `fail()` or `assert*()` will
not be executed.

`fail()` and `assert*()` throw AssertionErrors, which are a subtype of
Throwable. That means that if the catch block catches Throwable, then execution
will always jump to the catch block, and the test will always pass.

To fix this, you usually want to catch Exception rather than Throwable. If you
need to catch throwable (e.g., the expected exception is an AssertionError),
then add logic in your catch block to ensure that the AssertionError that was
caught is not the same one thrown by the call to `fail()` or `assert*()`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("TryFailThrowable")` to the enclosing element.

Using a type parameter as a qualifier in the name of a type or expression is equivalent to referencing the type parameter’s upper bound directly.

For example, this signature:

```
static <T extends Message> T populate(T.Builder builder) {}
```
Is identical to the following:

```
static <T extends Message> T populate(Message.Builder builder) {}
```
The use of `T.Builder` is unnecessary and misleading. Always refer to the type
by its canonical name `Message.Builder` instead.

This check may not be suppressed.

Suppress false positives by adding the suppression annotation @SuppressWarnings("UnicodeDirectionalityCharacters") to the enclosing element.

@SuppressWarnings("UnicodeDirectionalityCharacters")

Using non-ASCII Unicode characters in code can be confusing, and potentially unsafe.

For example, homoglyphs can result in a different method to the one that was expected being invoked.

```
import static com.google.common.base.Objects.equal;
public void isAuthenticated(String password) {
 // The "l" here is not what it seems.
 return equaⅼ(password, this.password());
}
// ...
private boolean equaⅼ(String a, String b) {
 return true;
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnicodeInCode")` to the enclosing element.

Alternate names: PreconditionsCheckNotNull, PreconditionsCheckNotNullPrimitive

Suppress false positives by adding the suppression annotation @SuppressWarnings("UnnecessaryCheckNotNull") to the enclosing element.

@SuppressWarnings("UnnecessaryCheckNotNull")

JLS §15.12.2.1 allows non-generic methods to be invoked with type arguments:

a non-generic method may be potentially applicable to an invocation that supplies explicit type arguments. Indeed, it may turn out to be applicable. In such a case, the type arguments will simply be ignored.

This rule stems from issues of compatibility and principles of substitutability. Since interfaces or superclasses may be generified independently of their subtypes, we may override a generic method with a non-generic one. However, the overriding (non-generic) method must be applicable to calls to the generic method, including calls that explicitly pass type arguments. Otherwise the subtype would not be substitutable for its generified supertype.

There is no reason to do this in cases where the method being called is not an override of a generic method.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnnecessaryTypeArgument")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("UnsafeWildcard") to the enclosing element.

@SuppressWarnings("UnsafeWildcard")

Creating a side-effect-free anonymous class and never using it is usually a mistake.

For example:

```
public static void main(String[] args) {
 new Thread(new Runnable() {
 @Override public void run() {
 preventMissionCriticalDisasters();
 }
 }); // did you mean to call Thread#start()?
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnusedAnonymousClass")` to the enclosing element.

Several of the methods in `java.util.Collections`, such as `sort` and `shuffle`,
modify collections in place. If you call one of these methods on a
newly-allocated collection and don’t use it later, you are doing unnecessary
work. You probably meant to keep a reference to the newly-allocated copy of your
collection and use that in the rest of your code.

For example, this code sorts a new `ArrayList` and then throws away the result,
returning the unsorted original collection:

```
public Collection<String> sort(Collection<String> foos) {
 Collections.sort(new ArrayList<>(foos));
 return foos;
}
```
The author probably meant:

```
public Collection<String> sort(Collection<String> foos) {
 List<String> sortedFoos = new ArrayList<>(foos);
 Collections.sort(sortedFoos);
 return sortedFoos;
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnusedCollectionModifiedInPlace")` to the enclosing element.

As of JDK 10 var is a restricted local variable type and cannot be used for type declarations (see JEP 286).

var

Suppress false positives by adding the suppression annotation @SuppressWarnings("VarTypeName") to the enclosing element.

@SuppressWarnings("VarTypeName")

When switching over a proto `one_of`, getters that don’t match the current case
are guaranteed to be return a default instance:

```
switch (foo.getBlahCase()) {
 case FOO:
 return foo.getFoo();
 case BAR:
 return foo.getFoo(); // should be foo.getBar()
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("WrongOneof")` to the enclosing element.

The `^` binary XOR operator is sometimes mistaken for a power operator, but e.g.
`2 ^ 2` evaluates to `0`, not `4`.

Consider expressing powers of `2` using a bit shift instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("XorPower")` to the enclosing element.

Avoid the magic constant (ZoneId.of(“Z”)) in favor of a more descriptive API: ZoneOffset.UTC

Suppress false positives by adding the suppression annotation @SuppressWarnings("ZoneIdOfZ") to the enclosing element.

@SuppressWarnings("ZoneIdOfZ")

Suppress false positives by adding the suppression annotation @SuppressWarnings("ASTHelpersSuggestions") to the enclosing element.

@SuppressWarnings("ASTHelpersSuggestions")

Avoid APIs that convert a hostname to a single IP address:

`java.net.Socket(String,int)``java.net.InetSocketAddress(String,int)``java.net.InetAddress.html#getByName(String)`Depending on the value of the
`-Djava.net.preferIPv6Addresses=true`
system property, those APIs will return an IPv4 or IPv6 address. If a client
only has IPv4 connectivity, it will fail to connect with
`-Djava.net.preferIPv6Addresses=true`. If a client only has IPv6 connectivity,
it will fail to connect with `-Djava.net.preferIPv6Addresses=false`.

The preferred alternative is for clients to consider all addresses returned by
`java.net.InetAddress.html#getAllByName(String)`,
and try to connect to each one until a successful connection is made.

TIP: To resolve a loopback address, prefer `InetAddress.getLoopbackAddress()`
over hard-coding an IPv4 or IPv6 loopback address with
`InetAddress.getByName("127.0.0.1")` or `InetAddress.getByName("::1")`.

This is, prefer this:

```
 Socket doConnect(String hostname, int port) throws IOException {
 IOException exception = null;
 for (InetAddress address : InetAddress.getAllByName(hostname)) {
 try {
 return new Socket(address, port);
 } catch (IOException e) {
 if (exception == null) {
 exception = e;
 } else {
 exception.addSuppressed(e);
 }
 }
 }
 throw exception;
 }
```
```
 Socket doConnect(String hostname, int port) throws IOException {
 IOException exception = null;
 for (InetAddress address : InetAddress.getAllByName(hostname)) {
 try {
 Socket s = new Socket();
 s.connect(new InetSocketAddress(address, port));
 return s;
 } catch (IOException e) {
 if (exception == null) {
 exception = e;
 } else {
 exception.addSuppressed(e);
 }
 }
 }
 throw exception;
 }
```
instead of this:

```
 Socket doConnect(String hostname, int port) throws IOException {
 return new Socket(hostname, port);
 }
```
```
 void doConnect(String hostname, int port) throws IOException {
 Socket s = new Socket();
 s.connect(new InetSocketAddress(hostname, port));
 }
```
```
 void doConnect(String hostname, int port) throws IOException {
 Socket s = new Socket();
 s.connect(new InetSocketAddress(InetAddress.getByName(hostname), port));
 }
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("AddressSelection")` to the enclosing element.

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
AlmostJavadoc
This comment contains Javadoc or HTML tags, but isn't started with a double asterisk (/**); is it meant to be Javadoc?
Severity
WARNING
Tags
Style

Checking a provably-identical condition twice can be a sign of a logic error.

For example:

```
public Optional<T> first(Optional<T> a, Optional<T> b) {
 if (a.isPresent()) {
 return a;
 } else if (a.isPresent()) { // Oops--should be checking `b`.
 return b;
 } else {
 return Optional.empty();
 }
}
```
It can also be a sign of redundancy, which can just be removed.

```
public void act() {
 if (enabled) {
 frobnicate();
 } else if (!enabled) { // !enabled is guaranteed to be true here, so the check can be removed
 doSomethingElse();
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("AlreadyChecked")` to the enclosing element.

The following code is fine in Java 7, but poses a problem in Java 8.

```
static class A {
 B c(D d) { return null; }
 static B c(A a, D d) { return null; }
}
```
The method reference `A::c` of the instance method `c` has an implicit first
parameter for the `this` pointer. So both methods that `A::c` could resolve to
are compatible with `BiFunction<A, D, B>`, and the method reference is
ambiguous.

```
void f(BiFunction<A, D, B> f) { ... }
```
```
error: incompatible types: invalid method reference
 f(A::c);
 ^
 reference to c is ambiguous
 both method c(A,D) in A and method c(D) in A match
```
Consider renaming one of the methods to avoid the ambiguity.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AmbiguousMethodReference")` to the enclosing element.

This method passes a pair of parameters through to `String#format`, but the
enclosing method wasn’t annotated `@FormatMethod`. Doing so gives compile-time
rather than run-time protection against malformed format strings. Consider
annotating the format string with
`@com.google.errorprone.annotations.FormatString` and the method with
`@FormatMethod` to allow compile-time checking for well-formed format strings.

```
static void log(String format, String... args) {
 Log.w(format, args);
}
void frobnicate(int a, int b) {
 if (a < b) {
 // Whoops: didn't provide enough format args.
 log("%s < %s", a);
 }
}
```
```
@FormatMethod
static void log(@FormatString String format, String... args) {
 Log.w(format, args);
}
```
WARNING: There’s a very high chance that manual intervention will be required
after applying this fix, either due to existing uses of the method not passing
in valid format strings, or methods which delegate to this one requiring the
`@FormatMethod` annotation as well. Please ensure that everything depending on
this code still compiles after applying.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AnnotateFormatMethod")` to the enclosing element.

If permuting the arguments of a method call means that the argument names are a better match for the parameter names than the original ordering then this might indicate that they have been accidentally swapped. There are also legitimate reasons for the names not to match such as when rotating an image (swap width and height). In this case we suggest annotating the names with a comment to make the deliberate swap clear to future readers of the code. Argument names annotated with a comment containing the parameter name will not generate a warning.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ArgumentSelectionDefectChecker")` to the enclosing element.

#
 ArrayAsKeyOfSetOrMap

 Arrays do not override equals() or hashCode, so comparisons will be done on reference equality only. If neither deduplication nor lookup are needed, consider using a List instead. Otherwise, use IdentityHashMap/Set, a Map from a library that handles object arrays, or an Iterable/List of pairs.

There are two main problems with having a component of a record be an array.

By default, the generated `equals` and `hashCode` will just call `equals` or
`hashCode` on the array. Two distinct arrays are never considered equal by
`equals` even if their contents are the same. The generated `toString` is
similarly not useful, since it will be something like `[B@723279cf`.

Arrays are mutable, but records should not be mutable. A client of a record with an array component can change the contents of the array.

Instead of an array component, consider something like `ImmutableList<String>`,
or, for primitive arrays, something like `ByteString` or `ImmutableIntArray`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ArrayRecordComponent")` to the enclosing element.

JUnit’s assertEquals (and similar) are defined to take the expected value first and the actual value second. Getting these the wrong way round will cause a confusing error message if the assertion fails.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AssertEqualsArgumentOrderChecker")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("AssertSameIncompatible") to the enclosing element.

@SuppressWarnings("AssertSameIncompatible")

Suppress false positives by adding the suppression annotation @SuppressWarnings("AssertThrowsBlockToExpression") to the enclosing element.

@SuppressWarnings("AssertThrowsBlockToExpression")

Prefer to minimize the amount of logic in a lambda passed to `assertThrows`.

An `assertThrows` assertion will pass if any code in the provided lambda throws
the given exception. If the lambda contains multiple expressions that may throw,
the test may incorrectly pass if an earlier expression unexpectedly throws. For
example, consider:

```
assertThrows(
 IllegalArgumentException.class,
 () ->
 EncodingUtil.convertAudioFormatToMp3Configuration(
 AudioEncoding.newBuilder()
 .setAudioFormat(AudioFormat.AUDIO_FORMAT_MP3)
 .setAudioQuality(AudioQuality.UNRECOGNIZED)
 .build()));
```
This assertion is intended to check that
`convertAudioFormatToMp3Configuration()` throws `IllegalArgumentException`, but
the assertion will always pass because the setup logic in
`setAudioQuality(AudioQuality.UNRECOGNIZED)` is incorrect and `build()` will
throw `IllegalArgumentException`.

Instead, you should minimize the amount of logic inside of your `assertThrows`:

```
AudioEncoding audioEncoding =
 AudioEncoding.newBuilder()
 .setAudioFormat(AudioFormat.AUDIO_FORMAT_MP3)
 .setAudioQuality(AudioQuality.UNRECOGNIZED) // BOOM goes the dynamite!
 .build();
assertThrows(
 IllegalArgumentException.class,
 () -> EncodingUtil.convertAudioFormatToMp3Configuration(audioEncoding));
```
The test above now (correctly) fails and exposes the setup issue.

This check tries to pick a reasonable name for the extracted variable, but there may be a better name. If the extracted variable seems verbose, consider renaming it to improve readability.

In some cases a judicious use of `var` could improve readability, following the
guidelines for usage of `var`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AssertThrowsMinimizer")` to the enclosing element.

If the body of the lambda passed to `assertThrows` contains multiple statements,
execution of the lambda will stop at the first statement that throws an
exception and all subsequent statements will be ignored.

This means that:

Don’t do this:

```
assertThrows(
 UnsupportedOperationException.class,
 () -> {
 AppendOnlyList list = new AppendOnlyList();
 list.add(0, "a");
 list.remove(0);
 assertThat(list).containsExactly("a");
 });
```
Do this instead:

```
AppendOnlyList list = new AppendOnlyList();
list.add(0, "a");
assertThrows(
 UnsupportedOperationException.class,
 () -> list.remove(0));
assertThat(list).containsExactly("a");
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("AssertThrowsMultipleStatements")` to the enclosing element.

JUnit’s `fail()` and `assert*` methods throw an `AssertionError`, so using the
try/fail/catch pattern to test for `AssertionError` (or any of its super-types)
is incorrect. The following example will never fail:

```
try {
 doSomething();
 fail("expected doSomething to throw AssertionError");
} catch (AssertionError expected) {
 // expected exception
}
```
To avoid this issue, prefer JUnit’s `assertThrows()` API:

```
import static com.google.common.truth.Truth.assertThat;
import static org.junit.Assert.assertThrows;
@Test
public void testFailsWithAssertionError() {
 AssertionError thrown = assertThrows(
 AssertionError.class,
 () -> {
 doSomething();
 });
 assertThat(thrown).hasMessageThat().contains("something went terribly wrong");
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("AssertionFailureIgnored")` to the enclosing element.

Using the result of an assignment expression can be quite unclear, except for
simple cases where the result is used to initialise another variable, such as
`x = y = 0`.

Consider a common pattern of lazily initialising a field:

```
class Lazy<T> {
 private T t = null;
 abstract T create();
 public T get() {
 if (t != null) {
 return t;
 }
 return t = create();
 }
}
```
```
class Lazy<T> {
 private T t = null;
 abstract T create();
 public T get() {
 if (t != null) {
 return t;
 }
 t = create();
 return t;
 }
}
```
At the cost of another line, it’s now clearer what’s happening. (Note that
neither the before nor the after are thread-safe; this particular example would
be better served with `Suppliers.memoizing`.)

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AssignmentExpression")` to the enclosing element.

Using @AssistedInject and @Inject on the same constructor is a runtimeerror in Guice.

Suppress false positives by adding the suppression annotation @SuppressWarnings("AssistedInjectAndInjectOnSameConstructor") to the enclosing element.

@SuppressWarnings("AssistedInjectAndInjectOnSameConstructor")

Because `0` is an integer constant, `-0` is an integer, too. Integers have no
concept of “negative zero,” so it is the same as plain `0`.

The value is then widened to a floating-point number. And while floating-point
numbers have a concept of “negative zero,” the integral `0` is widened to the
floating-point “positive” zero.

To write a negative zero, you have to write a constant that is a floating-point
number. One simple way to do that is to write `-0.0`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AttemptedNegativeZero")` to the enclosing element.

AutoValue classes reject `null` values, unless the property is annotated with
`@Nullable`. For this reason, the usage of boxed primitives (e.g. `Long`) is
discouraged, except when annotated as `@Nullable`. Otherwise they can be
replaced with the corresponding primitive. There could be some cases where the
usage of a boxed primitive might be intentional to avoid boxing the value again
after invoking the getter.

Suppress violations by using `@SuppressWarnings("AutoValueBoxedValues")` on the
relevant `abstract` getter and/or setter.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AutoValueBoxedValues")` to the enclosing element.

Consider that other developers will try to read and understand your value class while looking only at your hand-written class, not the actual (generated) implementation class. If you mark your concrete methods final, they won’t have to wonder whether the generated subclass might be overriding them. This is especially helpful if you are underriding equals, hashCode or toString!

Reference: https://github.com/google/auto/blob/master/value/userguide/practices.md#mark-all-concrete-methods-final

NOTE:
Since `@Memoized` methods can’t be final,
the check doesn’t flag them.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AutoValueFinalMethods")` to the enclosing element.

*Alternate names: mutable*

AutoValue instances should
be deeply immutable. Therefore, we recommend using immutable types for fields.
E.g., use `ImmutableMap` instead of `Map`, `ImmutableSet` instead of `Set`, etc.

Read more at:

Suppress violations by using `@SuppressWarnings("AutoValueImmutableFields")` on
the relevant `abstract` getter.

@AutoValue-annotated classes may form part of your API, but the AutoValue_ generated classes should not. The fact that the generated classes are visible to other classes within the same package is an implementation detail, and is best avoided. Ideally, any reference to the AutoValue_-prefixed class should be confined to a single factory method, with other factories delegating to it if necessary.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AutoValueSubclassLeaked")` to the enclosing element.

*Alternate names: JavaLangClash*

Class names from `java.lang` should never be reused. From
Java Puzzlers:

Avoid reusing the names of platform classes, and never reuse class names from

`java.lang`, because these names are automatically imported everywhere. Programmers are used to seeing these names in their unqualified form and naturally assume that these names refer to the familiar classes from`java.lang`. If you reuse one of these names, the unqualified name will refer to the new definition any time it is used inside its own package.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AvoidCommonTypeNames")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("AvoidValueSetter") to the enclosing element.

@SuppressWarnings("AvoidValueSetter")

A narrowing integral conversion can cause a sign flip, since it simply discards all but the n lowest order bits, where n is the number of bits used to represent the target type (JLS 5.1.3). In a compare or compareTo method, this can cause incorrect and unstable sort orders.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("BadComparable")` to the enclosing element.

#
 BadImport

 Importing nested classes/static methods/static fields with commonly-used names can make code harder to read, because it may not be clear from the context exactly which type is being referred to. Qualifying the name with that of the containing class can make the code clearer.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("BadImport")` to the enclosing element.

Flags `instanceof` checks where the expression can be determined to be a
supertype of the type it is compared to.

JLS 15.28
specifically calls `instanceof` out as *not* being a compile-time constant
expression, so the usage of this pattern can lead to unreachable code that won’t
be flagged by the compiler:

```
class Foo {
 void doSomething() {
 if (this instanceof Foo) { // BAD: always true
 return;
 }
 interestingProcessing();
 }
}
```
In general, an `instanceof` comparison against a superclass is equivalent to a
null check:

```
foo instanceof Foo
```
```
foo != null
```
Pattern-matching `instanceof`s introduce some extra complexity into this. It may
be tempting to use an `instanceof` check to define a narrowly-scoped local
variable which gets reused within an expression, for example,

```
return proto.getSubMessage() instanceof SubMessage sm
 && sm.getForename().equals("John")
 && sm.getSurname().equals("Smith");
```
We feel this urge should be resisted. While this is a clever trick to avoid an
extra line, it is not a true `instanceof` check, and declaring a variable
normally is clearer:

```
SubMessage sm = proto.getSubMessage();
return sm.getForename().equals("John") && sm.getSurname().equals("Smith");
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("BadInstanceof")` to the enclosing element.

Alternate names: InvalidPatternSyntax

Suppress false positives by adding the suppression annotation @SuppressWarnings("BareDotMetacharacter") to the enclosing element.

@SuppressWarnings("BareDotMetacharacter")

`BigDecimal`’s equals method compares the scale of the representation as well as
the numeric value, which may not be expected.

```
BigDecimal a = new BigDecimal("1.0");
BigDecimal b = new BigDecimal("1.00");
a.equals(b); // false!
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("BigDecimalEquals")` to the enclosing element.

BigDecimal’s `double` can lose precision in surprising ways.

```
 // these are the same:
 new BigDecimal(0.1)
 new BigDecimal("0.1000000000000000055511151231257827021181583404541015625")
```
Prefer the `BigDecimal.valueOf(double)` method or the `new BigDecimal(String)`
constructor.

NOTE `BigDecimal.valueOf(double)` does not suffer from the same problem; it is
equivalent to `new BigDecimal(Double.valueOf(double))`, and while `0.1` is not
exactly representable, `Double.valueOf(0.1)` yields `"0.1"`. As long as
FloatingPointLiteralPrecision doesn’t
generate a warning, `BigDecimal.valueOf` is safe.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("BigDecimalLiteralDouble")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("BooleanLiteral") to the enclosing element.

@SuppressWarnings("BooleanLiteral")

Constructors of primitive wrapper objects (e.g. `new Boolean(true)` will be
deprecated in Java 9. The `valueOf` factory methods (e.g.
`Boolean.valueOf(true)`) should always be preferred. Those methods are called
implicitly by autoboxing, which is often more convenient than an explicit call.
`Integer x = Integer.valueOf(23);` and `Integer x = 23;` are equivalent.

The explicit constructors always return a fresh instance, resulting in
unnecessary allocations. The `valueOf` methods return cached instances for
frequently requested values, offering significantly better space and time
performance.

Relying on the unique reference identity of the instances returned by the explicit constructors is extremely bad practice. Primitives should always be treated as identity-less value types, even in their boxed representations.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("BoxedPrimitiveConstructor")` to the enclosing element.

#
 BoxingComparator

 Comparator.comparing() unnecessarily boxes numerical primitives; please use the primitive-specific method instead (e.g., comparingInt()).

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("BoxingComparator")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("BugPatternNaming") to the enclosing element.

@SuppressWarnings("BugPatternNaming")

`ByteBuffer` provides a view of an underlying bytes storage. The two most common
implementations are the non-direct byte buffers, which are backed by a bytes
array, and direct byte buffers, which usually are off-heap directly mapped to
memory. Thus, not all `ByteBuffer` implementations are backed by a bytes array,
and when they are, the beginning of the array returned by `.array()` may not
necessarily correspond to the beginning of the `ByteBuffer`.

Since it’s so finicky, **use of .array() is discouraged. Use .get(...) to
copy the underlying data instead.**

But, if you *absolutely must* use `.array()` to look behind the `ByteBuffer`
curtain, check all of the following:

`.hasArray()``.arrayOffset()``.remaining()`If you know that the `ByteBuffer` was created locally with
`ByteBuffer.wrap(...)` or `ByteBuffer.allocate(...)`, it’s safe, you can use the
backing array normally, without checking the above.

Do this:

```
// Use `.get(...)` to copy the byte[] without changing the current position.
public void foo(ByteBuffer buffer) throws Exception {
 byte[] bytes = new byte[buffer.remaining()];
 buffer.get(bytes);
 buffer.position(buffer.position() - bytes.length); // Restores the buffer position
 // ...
}
```
or this:

```
// Use `.array()` only if you also check `.hasArray()`, `.arrayOffset()`, and `.remaining()`.
public void foo(ByteBuffer buffer) throws Exception {
 if (buffer.hasArray()) {
 int startIndex = buffer.arrayOffset();
 int curIndex = buffer.arrayOffset() + buffer.position();
 int endIndex = curIndex + buffer.remaining();
 // Access elements of `.array()` with the above indices ...
 }
}
```
or this:

```
// No checking necessary when the buffer was constructed locally with `allocate(...)`.
public void foo() throws Exception {
 ByteBuffer buffer = ByteBuffer.allocate(Long.SIZE / Byte.SIZE);
 buffer.putLong(1L);
 // ...
 buffer.array();
 // ...
}
```
or this:

```
// No checking necessary when the buffer was constructed locally with `wrap(...)`.
public void foo(byte[] bytes) throws Exception {
 ByteBuffer buffer = ByteBuffer.wrap(bytes);
 // ...
 buffer.array();
 // ...
}
```
Not this:

```
public void foo(ByteBuffer buffer) {
 byte[] dataAsBytesArray = buffer.array();
 // ...
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ByteBufferBackingArray")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("CacheLoaderNull") to the enclosing element.

@SuppressWarnings("CacheLoaderNull")

Prefer to express durations using the largest possible unit, e.g.
`Duration.ofDays(1)` instead of `Duration.ofSeconds(86400)`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("CanonicalDuration")` to the enclosing element.

Discarding an exception after calling `printStackTrace` should usually be
avoided.

```
try {
 // ...
} catch (IOException e) {
 logger.log(INFO, "something has gone terribly wrong", e);
}
```
```
try {
 // ...
} catch (IOException e) {
 throw new UncheckedIOException(e); // New in Java 8
}
```
```
try {
 // ...
} catch (IOException e) {
 e.printStackTrace();
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("CatchAndPrintStackTrace")` to the enclosing element.

Ignoring an exception and calling `fail()` is unnecessary, since an uncaught
exception will already cause a test to fail. It also makes the test output less
useful, since the exception’s message and stack trace is lost.

Do this:

```
@Test
public void testFoo() throws Exception {
 int x = foos(); // the test fails if this throws
 assertThat(x).isEqualTo(42);
}
```
or this:

```
@Test
public void testFoo() throws Exception {
 int x;
 try {
 x = foos();
 } catch (Exception e) {
 throw new AssertionError("the test failed", e); // wraps the exception with additional context
 }
 assertThat(x).isEqualTo(42);
}
```
Not this:

```
@Test
public void testFoo() {
 int x;
 try {
 x = foos();
 } catch (Exception e) {
 fail("the test failed"); // the exception message and stack trace is lost
 }
 assertThat(x).isEqualTo(42);
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("CatchFail")` to the enclosing element.

Assertions made *inside the implementation of another Truth assertion* should
use `check`, not `assertThat`.

Before:

```
class MyProtoSubject extends Subject {
 ...
 public void hasFoo(Foo expected) {
 assertThat(actual.foo()).isEqualTo(expected);
 }
}
```
After:

```
class MyProtoSubject extends Subject {
 ...
 public void hasFoo(Foo expected) {
 check("foo()").that(actual.foo()).isEqualTo(expected);
 }
}
```
Benefits of `check` include:

`myProto.foo()`, and
it includes the full value of `myProto` for reference.`assertWithMessage`, that message, which
is lost in the `assertThat` version, is shown by the `check` version.`check` makes it possible to test the assertion with `ExpectFailure` and
to use `Expect` or `assume`. `assertThat`, by contrast, overrides any
user-specified failure behavior.Suppress false positives by adding the suppression annotation `@SuppressWarnings("ChainedAssertionLosesContext")` to the enclosing element.

`Character.getNumericValue` has unexpected behaviour: it interprets A-Z as
base-36 digits with values 10-35, but also supports non-arabic numerals and
miscellaneous numeric unicode characters like ㊷. For example:

`Character.getNumericValue('V' /* ASCII V */) == 31``Character.getNumericValue('Ⅴ' /* U+2164, Roman numeral 5 */) == 5``Character.getNumericValue('௧' /* U+0BF2, Tamil Digit One */) == 1````
Character.getNumericValue('௲' /* U+0BF2, Tamil Number One Thousand */) ==
1000
```
```
Character.getNumericValue('㊷' /* U+32B7, Circled Number Forty Two */) ==
42
```
`UCharacter.getNumericValue`
has the same behavior.

Consider using:

`UCharacter.getUnicodeNumericValue`:
Handles all unicode codepoints with numeric values including fractions,
roman numerals, and other miscellaneous numeric characters. Returns the
value as a double. Does not assign a value to A-Z.`Character.digit`:
Handles unicode codepoints in the “decimal digit” category and the letters
A-Z which are interpreted as base-36 digits with values 10-35. Does not
handle characters like roman numerals and ㊷. You can use a `radix` value of
10 or less to avoid interpreting A-Z as digits.Suppress false positives by adding the suppression annotation `@SuppressWarnings("CharacterGetNumericValue")` to the enclosing element.

*Alternate names: InnerClassMayBeStatic*

An inner class should be static unless it references members of its enclosing class. An inner class that is made non-static unnecessarily uses more memory and does not make the intent of the class clear.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ClassCanBeStatic")` to the enclosing element.

To avoid deadlocks, class initializers should not reference subtypes of the current class.

For example:

```
class Foo {
 public static final Bar INSTANCE = new Bar();
 public static class Bar extends Foo {}
}
```
There is a circular reference between the class initializers for `Foo` and
`Bar`: `Foo` depends on `Bar` in the initializer for a `static` field, and
initializing `Bar` requires initializing its supertype `Foo`. If one thread
starts initializing `Foo` and another thread simultaneously starts initializing
`Bar`, it will result in a deadlock.

The best solution is to refactor to break the cycle, by defining the constant field in a separate class from the supertype of the field.

For example, using this approach to fix the sample above would result in something like:

```
class Foos {
 public static final Bar INSTANCE = new Bar();
 public static class Foo {}
 public static class Bar extends Foo {}
}
```
That refactoring may be too invasive (say the code is part of an API, and there are many references to the current structure).

If the subclass is never referenced outside the current file (i.e. `Bar` is
never used outside of `Foo`, it is only referenced via `Foo.INSTANCE`), making
`Bar` `private` makes deadlocks less likely (see caveats below in the discussion
about `private` classes):

```
class Foo {
 public static final Foo INSTANCE = new Bar();
 private static class Bar extends Foo {}
}
```
If the subclass *is* referenced outside the current field, deadlocks can be
avoided by ensuring that the subclass has only private constructors (or `static`
factory methods), so that the only way to initialize the subclass is to first
initialize the containing class:

```
class Foo {
 public static final Foo INSTANCE = new Bar();
 private static class Bar extends Foo {
 private Bar() {}
 }
}
```
If the subclass needs to be directly created by code outside the current file, a static factory can be added as a member of the outer class, for example:

```
class Foo {
 public static final Foo INSTANCE = new Bar();
 private static class Bar extends Foo {
 private Bar() {}
 }
 public static Bar createBar() {
 return new Bar();
 }
}
```
AutoValue implementation classes are necessarily non-`private`, since they are
generated into separate files. However examples like the following can’t
deadlock as long as the only reference to the `AutoValue_Base` class is inside
`Base`, since there is no way for a thread to cause `AutoValue_Base` to be
initialized without first having initialized `Base`:

```
@AutoValue
abstract class Base {
 abstract String bar();
 static final Object DEFAULT = new AutoValue_Base("bar");
 static Base of(String bar) {
 return new AutoValue_Base(bar);
 }
}
```
There is a separate Error Prone check, https://errorprone.info/bugpattern/AutoValueSubclassLeaked, to
prevent `AutoValue_` classes from being accessed outside the file containing the
corresponding `@AutoValue` base class.

The check ignores references that cross from a `private` inner class (or any
class inside it) to its immediately enclosing class, e.g.

```
public class A {
 private static Object benignCycle = new B.C();
 private static class B {
 public static class C extends A { }
 }
}
```
There is a cycle `A` -> `A.B.C` -> `A`, but (without reflection) it’s not
possible to access `A.B.C` in a way that causes initialization until after A is
initialized.

There are situations where deadlocks involving `private` classes can still
occur, but the heuristic of ignoring paths cycles from `private` members is good
enough for most real-world examples of deadlocks that have been observed.

In the following, `A.C` can trigger initialization of `B`, despite `B` being
private.

```
 public class A {
 private static Object bad_cycle = new B();
 private static class B extends A { }
 public static class C extends B { }
 }
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ClassInitializationDeadlock")` to the enclosing element.

The documentation for `Class#newInstance` includes the following
warning:

Note that this method propagates any exception thrown by the nullary constructor, including a checked exception. Use of this method effectively bypasses the compile-time exception checking that would otherwise be performed by the compiler. The

`Constructor.newInstance`method avoids this problem by wrapping any exception thrown by the constructor in a (checked)`InvocationTargetException`.

Always prefer `myClass.getConstructor().newInstance()` to calling
`myClass.newInstance()` directly. The `Class#newInstance` method is slated for
deprecation in JDK 9.

Note that migrating to `Class#getConstructor()` and `Constructor#newInstance`
requires handling three new exceptions: `IllegalArgumentException`,
`NoSuchMethodException`, and `InvocationTargetException`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ClassNewInstance")` to the enclosing element.

If you provide
`Closeable`
resources through dependency injection, it can be difficult to effectively
manage the lifecycle of the Closable:

```
class MyModule extends AbstractModule() {
 @Provides
 FileOutputStream provideFileStream() {
 return new FileOutputStream("/tmp/outfile");
 }
}
...
class Client {
 private final FileOutputStream fos;
 @Inject Client(FileOutputStream fos, OtherDependency other, ...) {
 this.fos = fos;
 }
 void doSomething() throws IOException {
 fos.write("hello!");
 }
}
```
There are a number of issues with this approach as it relates to resource management:

It’s not clear which class has the responsibility of closing the
`FileOutputStream` resource:

`Client`, then it
makes sense for the client to close it.`@Singleton` to the
`@Provides` method), then suddenly it’s `Client` to close the resource. If one `Client` closes the stream, then
all of the other users of that stream will be dealing with a closed
resource. There needs to be some other ‘resource manager’ object that
also gets that stream, and its closing functions are call in the right
place. That can be tricky to do correctly.`Client`
fails (resulting in a `ProvisionException`), then the `FileOutputStream`
that may have been constructed leaks and isn’t properly closed, even if
`Client` normally closes its resources correctly.The preferred solution is to not inject closable resources, but instead, objects that can expose short-lived closable resources that are used as necessary. The following example uses Guava’s CharSink as the resource manager object:

```
class MyModule extends AbstractModule() {
 @Provides
 CharSink provideCharSink() {
 return Files.asCharSink(new File("/tmp/outfile"), StandardCharsets.UTF_8);
 }
}
...
class Client {
 private final CharSink sink;
 @Inject Client(CharSink sink, OtherDependency other, ...) {
 this.sink = sink;
 }
 void doSomething() throws IOException {
 sink.write("hello!"); // Opens the file at this point, and closes once its done.
 }
}
```
If there’s not a similar non-closable resource, you can write a simple wrapper:

```
class ResourceManager {
 @Inject ResourceManager(@Config String configs, ...) {}
 /**
 * Returns a new thing for you to use and dispose of
 */
 OutputStream provideInstance() { return new...(); }
}
...
class Client {
 private final ResourceManager resource;
 @Inject Client(ResourceManager resource, OtherDependency other, ...) {
 this.resource = resource;
 }
 void doSomething() {
 try (OutputStream actualStream = resource.provideInstance()) {
 // write to actualStream, closing with try-with-resources
 }
 }
}
```
This pattern can be extended to other resources: as opposed to injecting database connection handles directly, inject connection pool objects that require your object to ask for those connection objects when they’re needed and close them safely.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("CloseableProvides")` to the enclosing element.

Closing the standard output streams `System.out` or `System.err` will cause all
subsequent standard output to be dropped, including stack traces from exceptions
that propagate to the top level.

Avoid using try-with-resources to manage `PrintWriter`s or `OutputStream`s that
wrap `System.out` or `System.err`, since the try-with-resource statement will
close the underlying streams.

That is, prefer this:

```
PrintWriter pw = new PrintWriter(new OutputStreamWriter(System.err));
pw.println("hello");
pw.flush();
```
Instead of this:

```
try (PrintWriter pw = new PrintWriter(new OutputStreamWriter(System.err))) {
 pw.println("hello");
}
```
Consider the following example:

```
import java.io.OutputStreamWriter;
import java.io.PrintWriter;
import static java.nio.charset.StandardCharsets.UTF_8;
public class X {
 public static void main(String[] args) {
 System.err.println("one");
 try (PrintWriter err = new PrintWriter(new OutputStreamWriter(System.err, UTF_8))) {
 err.print("two");
 }
 // System.err has been closed, no more output will be printed!
 System.err.println("three");
 throw new AssertionError();
 }
}
```
The program will print the following, and return with exit code 1. Note that the
last `println` doesn’t produce any output, and the exception’s stack trace is
not printed:

```
one
two
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ClosingStandardOutputStreams")` to the enclosing element.

Using `Collection`s (and other types which rely on equality) to contain elements
with undefined equality is error-prone. For example, `Collection` itself does
not have well-defined equality: the `List` and `Set` subinterfaces are not
necessarily comparable.

```
ImmutableList<Collection<Integer>> collectionsOfIntegers =
 ImmutableList.of(ImmutableSet.of(1, 2), ImmutableSet.of(3, 4));
collectionsOfIntegers.contains(ImmutableSet.of(1, 2)); // true
collectionsOfIntegers.contains(ImmutableList.of(1, 2)); // false
```
```
boolean containsTest(Collection<CharSequence>> charSequences) {
 // True if `charSequences` actually contains Strings, but otherwise not necessarily.
 return charSequences.contains("test");
}
```
In this case, an appropriate fix may be,

```
boolean containsTest(Collection<CharSequence>> charSequences) {
 return charSequences.stream().anyMatch("test"::contentEquals);
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("CollectionUndefinedEquality")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("CollectorShouldNotUseState") to the enclosing element.

@SuppressWarnings("CollectorShouldNotUseState")

A `Comparator` is an object that knows how to compare other objects, whereas an
object implementing `Comparable` knows how to compare itself to other objects of
the same type.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ComparableAndComparator")` to the enclosing element.

The contract for `Comparator#compare` and `Comparable#compareTo` states that the
result is an integer which is `< 0` for less than, `== 0` for equality and `> 0`
for greater than. While most implementations return `-1`, `0` and `+1` for those
cases respectively, this is not guaranteed. Always comparing to `0` is the
safest use of the return value.

```
 boolean <T> isLessThan(Comparator<T> comparator, T a, T b) {
 // Fragile: it's not guaranteed that `comparator` returns -1 to mean
 // "less than".
 return comparator.compare(a, b) == -1;
 }
```
```
 boolean <T> isLessThan(Comparator<T> comparator, T a, T b) {
 return comparator.compare(a, b) < 0;
 }
```
Even comparisons which are otherwise correct are clearer to other readers of the
code if turned into a comparison to `0`, e.g.:

```
 boolean <T> greaterThan(Comparator<T> comparator, T a, T b) {
 return comparator.compare(a, b) >= 1;
 }
```
```
 boolean <T> greaterThan(Comparator<T> comparator, T a, T b) {
 return comparator.compare(a, b) > 0;
 }
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("CompareToZero")` to the enclosing element.

When a boolean expression is a compile-time constant (e.g.: `2 < 1`, `1 == 1`,
`'a' < 'A'`), these expressions can be directly replaced with `true` or `false`,
as appropriate. In any context where these expressions are used, `true` or
`false` is a more readable alternative:

```
if (2 < 1) {
 // Some code I don't want to run right now
}
while (1 == 1) {
 // Some loop that I will manually break out of
}
assert 1 != 2; // I want to force an AssertionFailure if assertions are enabled
```
```
if (false) {
 // Some code I don't want to run right now
}
while (true) {
 // Some loop that I will manually break out of
}
assert false; // I want to force an AssertionFailure if assertions are enabled
```
When some boolean expression is a compile-time constant unexpectedly, it generally represents a bug in the code:

```
for (int i = 0; i < 100; i++) {
 System.out.println("Is " + i + " greater than 50?: " + (1 > 50));
}
// Prints "... false" 100 times, since i > 50 is mistyped as 1 > 50
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ComplexBooleanConstant")` to the enclosing element.

java.util.Date uses 1900-based years, 0-based months, 1-based days, and 0-based hours/minutes/seconds. Additionally, it allows for negative values or very large values (which rollover).

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DateChecker")` to the enclosing element.

`DateFormat` is not thread-safe. The documentation recommends creating
separate format instances for each thread. If multiple threads access a format
concurrently, it must be synchronized externally.

The Google Java Style Guide §5.2.4 requires `CONSTANT_CASE` to only be
used for static final fields whose contents are deeply immutable and whose
methods have no detectable side effects, so fields of type `DateFormat` should
not use `CONSTANT_CASE`.

TIP: Consider using the `java.time` API added in Java8, in particular
`DateTimeFormatter`. One its many advantages over `DateFormat` is that it is
immutable and thread-safe.

If the date formatter is accessed by multiple threads, consider using
`ThreadLocal`:

```
private static final ThreadLocal<DateFormat> DATE_FORMAT =
 ThreadLocal.withInitial(() -> new SimpleDateFormat("yyyy-MM-dd HH:mm"));
```
If the field is never accessed by multiple threads, rename it to use
`lowerCamelCase`.

```
// not thread safe
private static final DateFormat dateFormat =
 new SimpleDateFormat("yyyy-MM-dd HH:mm");
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("DateFormatConstant")` to the enclosing element.

Having an extremely long Java statement with many chained method calls can cause
compilation to fail with a `StackOverflowError` when the compiler tries to
recursively process it.

This is a common problem in generated code.

As an alternative to extremely long chained method calls, e.g. for builders, consider something like the following for collections with hundreds or thousands of entries:

```
private static final ImmutableList<String> FEATURES = createFeatures();
private static final ImmutableList<String> createFeatures() {
 ImmutableList.Builder<String> builder = ImmutableList.<String>builder();
 builder.add("foo");
 builder.add("bar");
 ...
 return builder.build();
}
```
over code like this:

```
private static final ImmutableList<String> FEATURES =
 ImmutableList.<String>builder()
 .add("foo")
 .add("bar")
 ...
 .build();
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("DeeplyNested")` to the enclosing element.

A `Charset` is a mapping between sequences of
16-bit Unicode code units and sequences of bytes. Charsets are used
when encoding characters into bytes and decoding bytes into characters.

Using APIs that rely on the JVM’s default Charset under the hood is dangerous.
The default charset can vary from machine to machine or JVM to JVM. This can
lead to unstable character encoding/decoding between runs of your program, even
for ASCII characters (e.g.: `A` is `0100 0001` in `UTF-8`, but is ```
0000 0000
0100 0001
```
 in `UTF-16`).

If you need stable encoding/decoding, you must specify an explicit charset. The
`StandardCharsets` class provides these constants for you.

When in doubt, use UTF-8.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DefaultCharset")` to the enclosing element.

Declaring classes in the default package is discouraged.

Suppress false positives by adding the suppression annotation @SuppressWarnings("DefaultPackage") to the enclosing element.

@SuppressWarnings("DefaultPackage")

`@Deprecated` annotations should not be applied to local variables and
parameters, since they have no effect there.

The
javadoc for `@Deprecated`
says

Use of the

`@Deprecated`annotation on a local variable declaration or on a parameter declaration or a package declaration has no effect on the warnings issued by a compiler.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DeprecatedVariable")` to the enclosing element.

Direct invocations on mocks should be avoided in tests.

When you call a method on a mock, the call normally does only what you have
configured it to do (through calls to `when(...).thenReturn/thenAnswer`, etc.)
and makes a record of the call that can be read by later `verify(...)` calls.
Both of these are rarely what you want:

`verify(foo).bar()`, then the
reader will expect the test to succeed only because the code under test
called `bar()`, not because the test itself did.Sometimes, test authors, especially those familiar with other mocking frameworks (like EasyMock), will call a method on a mock for one of two reasons:

`verify(foo).bar()`, and it must do so ```
@Test
public void balanceIsChecked() {
 Account account = mock(Account.class);
 LoanChecker loanChecker = new LoanChecker(account);
 assertThat(loanChecker.checkEligibility()).isFalse();
 // Should be verify(account).checkBalance();, or be removed if the call to
 // `checkEligibility` is sufficient proof the code is behaving as intended.
 account.checkBalance();
}
```
There is at least one edge case in which a call to a mock has different effects:
A call to a `final` method will normally *not* be intercepted by Mockito, so it
will run the implementation of that method in the code under test. Sometimes,
that method will call other methods on the mock object, producing effects
similar to if the test had called those methods directly. Sometimes, the method
will have other effects. Both kinds of effects can be confusing, so prefer to
avoid such calls when possible.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DirectInvocationOnMock")` to the enclosing element.

Various methods which take variable-length arguments throw runtime exceptions
like `IllegalArgumentException` when the arguments are not distinct.

This checker warns on using non-distinct parameters in various varargs methods when the usage is either redundant or will result in a runtime exception.

Bad:

```
ImmutableSet.of(first, second, second, third);
```
Good:

```
ImmutableSet.of(first, second, third);
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("DistinctVarargsChecker")` to the enclosing element.

#
 DoNotCallSuggester

 Consider annotating methods that always throw with @DoNotCall. Read more at https://errorprone.info/bugpattern/DoNotCall

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DoNotCallSuggester")` to the enclosing element.

Do not ‘claim’ annotations in annotation processors; `Processor#process`
should unconditionally `return false`. Claiming annotations prevents other
processors from seeing the annotations, which there’s usually no reason to do.
It’s also fragile, since it relies on the order the processors run in, and
there’s no robust way in most build systems to ensure a particular processor
sees the annotations first.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DoNotClaimAnnotations")` to the enclosing element.

`@AutoValue` is used to represent pure data classes. Mocking these should not be
necessary: prefer constructing them in the same way production code would.

To make the argument another way: the fact that `AutoValue` classes are not
`final` is an implementation detail of the way they’re generated. They should be
regarded as logically final insofar as they must not be extended by
non-generated code. If they were final, they also would not be mockable.

Instead of mocking:

```
@Test
public void test() {
 MyAutoValue myAutoValue = mock(MyAutoValue.class);
 when(myAutoValue.getFoo()).thenReturn("foo");
}
```
Prefer simply constructing an instance:

```
@Test
public void test() {
 MyAutoValue myAutoValue = MyAutoValue.create("foo");
}
```
If your `AutoValue` has multiple required fields, and only one is relevant for a
test, consider using a builder or `with`-style methods to create
test instances with just the fields you care about. Consider using

```
private MyAutoValue.Builder myAutoValueBuilder() {
 return MyAutoValue.builder().bar(42).baz(false);
}
@Test
public void test() {
 MyAutoValue myAutoValue = myAutoValueBuilder.foo("foo").build();
}
```
or:

```
private static final MyAutoValue MY_AUTO_VALUE =
 MyAutoValue.create(/* foo= */ "", /* bar= */ 42, /* baz= */ false);
@Test
public void test() {
 MyAutoValue myAutoValue = MY_AUTO_VALUE.withFoo("foo");
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("DoNotMockAutoValue")` to the enclosing element.

Using double-checked locking on mutable objects in non-volatile fields is not thread-safe.

If the field is not volatile, the compiler may re-order the code in the accessor. For more information, see:

The canonical example of *correct* double-checked locking for lazy
initialization is:

```
class Foo {
 /** This foo's bar. Lazily initialized via double-checked locking. */
 private volatile Bar bar;
 public Bar getBar() {
 Bar value = bar;
 if (value == null) {
 synchronized (this) {
 value = bar;
 if (value == null) {
 bar = value = computeBar();
 }
 }
 }
 return value;
 }
 private Bar computeBar() { ... }
}
```
Double-checked locking should only be used in performance critical classes. For code that is less performance sensitive, there are simpler, more readable approaches. Effective Java recommends two alternatives:

For lazily initializing instance fields, consider a *synchronized accessor*. In
modern JVMs with efficient uncontended synchronization the performance
difference is often negligible.

```
// Lazy initialization of instance field - synchronized accessor
private Object field;
synchronized Object get() {
 if (field == null) {
 field = computeValue();
 }
 return field;
}
```
If the field being initialized is static, consider using the *lazy
initialization holder class* idiom:

```
// Lazy initialization holder class idiom for static fields
private static class Holder {
 static final Object field = computeValue();
}
static Object get() {
 return Holder.field;
}
```
If the object being initialized with double-checked locking is
immutable,
then it is safe for the field to be non-volatile. *However*, the use of volatile
is still encouraged because it is almost free on x86 and makes the code more
obviously correct.

Note that immutable has a very specific meaning in this context:

[An immutable object] is transitively reachable from a final field, has not changed since the final field was set, and a reference to the object containing the final field did not escape the constructor.

Double-checked locking on non-volatile fields is in general unsafe because the
compiler and JVM can re-order code from the object’s constructor to occur
*after* the object is written to the field.

The final modifier prevents that re-ordering from occurring, and guarantees that all of the object’s final fields have been written to before a reference to that object is published.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DoubleCheckedLocking")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("DuplicateAssertion") to the enclosing element.

@SuppressWarnings("DuplicateAssertion")

Branching constructs (`if` statements, `conditional` expressions) should contain
difference code in the two branches. Repeating identical code in both branches
is usually a bug.

For example:

```
condition ? same : same
```
```
if (condition) {
 same();
} else {
 same();
}
```
this usually indicates a typo where one of the branches was supposed to contain different logic:

```
condition ? something : somethingElse
```
```
if (condition) {
 doSomething();
} else {
 doSomethingElse();
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("DuplicateBranches")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("DuplicateDateFormatField") to the enclosing element.

@SuppressWarnings("DuplicateDateFormatField")

Suppress false positives by adding the suppression annotation @SuppressWarnings("EffectivelyPrivate") to the enclosing element.

@SuppressWarnings("EffectivelyPrivate")

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
EmptyBlockTag
A block tag (@param, @return, @throws, @deprecated) has an empty description. Block tags without descriptions don't add much value for future readers of the code; consider removing the tag entirely or adding a description.
Severity
WARNING

The Google Java Style Guide §6.2 states:

It is very rarely correct to do nothing in response to a caught exception. (Typical responses are to log it, or if it is considered “impossible”, rethrow it as an AssertionError.)

When it truly is appropriate to take no action whatsoever in a catch block, the reason this is justified is explained in a comment.

When writing tests that expect an exception to be thrown, prefer using
`Assert.assertThrows` instead of writing a try-catch. That is,
prefer this:

```
assertThrows(NoSuchElementException.class, () -> emptyStack.pop());
```
instead of this:

```
try {
 emptyStack.pop();
 fail();
} catch (NoSuchElementException expected) {
}
```

When using Dagger Multibinding, you can use methods like the below to make sure that there’s a (potentially empty) Set binding for your type:

```
@Provides @ElementsIntoSet Set<MyType> provideEmptySetOfMyType() {
 return new HashSet<>();
}
```
However, there’s a slightly easier way to express this:

```
@Multibinds abstract Set<MyType> provideEmptySetOfMyType();
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("EmptySetMultibindingContributions")` to the enclosing element.

A semi-colon at the top level of a Java file is treated as an empty type declaration in the grammar, but it’s confusing and unnecessary.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("EmptyTopLevelDeclaration")` to the enclosing element.

You should almost never invoke the `Enum.ordinal()` method, nor depend on some
enum constant being at a particular index of the `values()` array. The ordinal
of a given enum value is not guaranteed to be stable across builds because of
the potential for enum values to be added, removed, or reordered. The ordinal
exists only to support low-level utilities like `EnumSet`.

Prefer using enum value directly:

```
ImmutableMap<MyEnum, String> MAPPING =
 ImmutableMap.<MyEnum, String>builder()
 .put(MyEnum.FOO, "Foo")
 .put(MyEnum.BAR, "Bar")
 .buildOrThrow();
```
instead of relying on the ordinal:

```
ImmutableMap<Integer, String> MAPPING =
 ImmutableMap.<Integer, String>builder()
 .put(MyEnum.FOO.ordinal(), "Foo")
 .put(MyEnum.BAR.ordinal(), "Bar")
 .buildOrThrow();
```
If you need a stable number for serialisation, consider defining an explicit field on the enum:

```
enum MyStableEnum {
 FOO(1),
 BAR(2),
 ;
 private final int wireCode;
 MyStableEnum(int wireCode) {
 this.wireCode = wireCode;
 }
}
```
rather than relying on the ordinal values:

```
enum MyUnstableEnum {
 FOO,
 BAR,
}
MyUnstableEnum fromWire(int wireCode) {
 return MyUnstableEnum.values()[wireCode];
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("EnumOrdinal")` to the enclosing element.

Requiring the argument to `equals` to be of a specific concrete type is
incorrect, because it precludes any subtype of this class from obeying the
substitution principle.
Such code should be modified to use an `instanceof` test instead of `getClass`.

TL;DR: use composition rather than inheritance to add fields to value types.

The most common objection to this rule arises from a scenario like the following:

```
class Point {
 double x() {...}
 double y() {...}
}
class ColoredPoint {
 double x() {...}
 double y() {...}
 Color color() {...}
}
```
It’s reasonable to think first of modeling `ColoredPoint` as a subclass of
`Point`. First, because on a conceptual level, there is a clear “is-a”
relationship between the two, which we have been taught to model as inheritance.
But also because of a few concrete advantages:

`x` and `y` fields/accessors for free.`Point` for free.`ColoredPoint` to anything that expects a `Point`. (This is
by far the primary advantage, since the first two are just one-time-only
implementation helpers for the `Point` authors themselves.)Although these same advantages *can* be achieved via composition, it’s no longer
quite “for free”; the frequent need to call `myColoredPoint.asPoint()` is a pain
that feels unjustified, and thus subclassing is often chosen.

Unfortunately, `equals` now creates a big problem. Two `Points` should be seen
as interchangeable whenever they have the same `x` and `y` coordinates. But two
`ColoredPoints` are only equivalent if their coordinates *and* color are the
same. This, plus the general contract of `equals` (e.g. symmetry), ends up
forcing `Point(1, 2)` to respond `false` if asked whether it is equal to
`ColoredPoint(1, 2, BLUE)`.

This can be achieved by changing `equals` to be based on `getClass` instead of
`instanceof`. But do we even want to do that? Recall the main advantage we cited
for using subtyping: so that we “can pass `ColoredPoint` to anything that
expects a `Point`”, we said. Yet we are in fact *not* achieving that after all.
Put simply, `ColoredPoint` is incapable of functioning properly *as a Point*
because it cannot participate in equality checks with other

`Point`s.The composition approach may be annoying (the frequent need for `.asPoint()`),
but aside from that annoyance, it at least achieves the three advantages
correctly.

In summation, while it is true *conceptually* that a `ColoredPoint` “is-a”
`Point`, a `ColoredPoint` is nevertheless unable to properly *function* as a
`Point`, and modeling it with subtyping is not a good idea for that reason.

In this case there is no disadvantage to the `getClass` trick - but there’s no
great advantage to it either. Most unsafe idioms have circumstances in which
they are safe, but this doesn’t change the fact that they are *generally* unsafe
and not worth propagating and legitimizing.

See Effective Java 3rd Edition §10 (“Obey the general contract when overriding equals”).

Suppress false positives by adding the suppression annotation `@SuppressWarnings("EqualsGetClass")` to the enclosing element.

Consider the following code:

```
String x = "42";
Integer y = 42;
if (x.equals(y)) {
 System.out.println("What is this, Javascript?");
} else {
 System.out.println("Types have meaning here.");
}
```
We understand that no `Integer` will be equal to any `String`. However, the
signature of the `equals` method accepts any Object, so the compiler will
happily allow us to pass an Integer to the equals method. That method will
always return false, which is probably not what we intended.

This check detects circumstances where the equals method is called when the two
objects in question can *never* be equal to each other. We check the following
equality methods:

`java.lang.Object.equals(Object)``java.util.Objects.equals(Object, Object)``com.google.common.base.Objects.equal(Object, Object)`Good! Many tests of equals methods neglect to test that equals on an unrelated object return false.

We recommend using Guava’s EqualsTester to perform tests of your equals method. Simply give it a collection of objects of your class, broken into groups that should be equal to each other, and EqualsTester will ensure that:

`hashCode` of each object in a group is the same as the hash code of
each other member of the groupWhich should exhaustively check all of the properties of `equals` and
`hashCode`.

The javadoc of `Object.equals(Object)` defines object equality very
precisely:

The equals method implements an equivalence relation on non-null object references:

It is reflexive: for any non-null reference value x, x.equals(x) should return true.

It is symmetric: for any non-null reference values x and y, x.equals(y) should return true if and only if y.equals(x) returns true.

It is transitive: for any non-null reference values x, y, and z, if x.equals(y) returns true and y.equals(z) returns true, then x.equals(z) should return true.

It is consistent: for any non-null reference values x and y, multiple invocations of x.equals(y) consistently return true or consistently return false, provided no information used in equals comparisons on the objects is modified.

For any non-null reference value x, x.equals(null) should return false.

TIP: EqualsTester validates each of these properties.

For most simple value objects (e.g.: a `Point` containing `x` and `y`
coordinates), this generally means that the equals method will only return true
if the other object has the exact same class, and each of the components is
equal to the corresponding component in the other object. Here, there are
numerous tools in the Java ecosystem to generate the appropriate `equals` and
`hashCode` method implementations, including AutoValue.

Another pattern often seen is to declare a common supertype with a defined
`equals` method (like `List`, which defines equality by having equal elements in
the same order). Then, different subclasses of that supertype (`LinkedList` and
`ArrayList`) can be equal to other classes with that supertype, since the
concrete class of the `List` is irrelevant. This checker will allow these types
of equality, as we detect when two objects share a common supertype with an
`equals` implementation and allow that to succeed.

Outside of these two general groups of equals methods, however, it’s very
difficult to produce correctly-behaving equals methods. Most of the time, when
`equals` is implemented in a non-obvious manner, one or more of the properties
above isn’t satisfied (generally the symmetric property). This can result in
subtle bugs, explained below.

`equals()````
class Foo {
 private String foo; // Some property
 public boolean equals(Object other) {
 if (other instanceof String) {
 return other.equals(foo); // We want to be able to call equals with a String
 }
 if (other instanceof Foo) {
 return ((Foo) other).foo.equals(foo); // Simplified, avoid null checks
 }
 return false;
 }
 public int hashCode() {
 return foo.hashCode();
 }
}
```
Here, `Foo`’s equals method is defined to accept a `String` value in addition to
other `Foo`’s. This may appear to work at first, but you end up with some
complex situations:

```
Foo a = new Foo("hello");
Foo b = new Foo("hello");
String hi = "hello";
if (a.equals(b)) {
 System.out.println("yes"); // Is printed, expected
}
if (b.equals(hi)) {
 System.out.println("yes"); // Is printed, abusing equals
}
if (hi.equals(b)) {
 System.out.println("no"); // Isn't printed, since String doesn't equals() Foo
}
Set<Foo> set = new HashSet<Foo>();
set.add(a);
set.add(b);
if (set.contains(hi)) {
 // Maybe? Depends on which way HashSet decides to call .equals()
 System.out.println("contained");
 // Is it removed? It's not guaranteed to be, since the .equals() method could
 // be called the other way in the remove path. Object.equals documentation
 // specifies it's supposed to be symmetric, so this could work.
 boolean removed = set.remove(hi);
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("EqualsIncompatibleType")` to the enclosing element.

Implementations of `#equals` should return `false` for different types, not
throw.

```
class Data {
 private int a;
 @Override
 public boolean equals(Object other) {
 Data that = (Data) other; // BAD: This may throw ClassCastException.
 return a == that.a;
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("EqualsUnsafeCast")` to the enclosing element.

Don’t implement `#equals` using just a `hashCode` comparison:

```
class MyClass {
 private final int a;
 private final int b;
 private final String c;
 ...
 @Override
 public boolean equals(@Nullable Object o) {
 return o.hashCode() == hashCode();
 }
 @Override
 public int hashCode() {
 return Objects.hashCode(a, b, c);
 }
```
The number of `Object`s with randomly distributed `hashCode` required to give a
50% chance of collision (and therefore, with this pattern, erroneously correct
equality) is only ~77k.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("EqualsUsingHashCode")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("ErroneousBitwiseExpression") to the enclosing element.

@SuppressWarnings("ErroneousBitwiseExpression")

Whenever a `ThreadPoolExecutor` is constructed with an unbounded `workQueue`,
the pool size will never go beyond `corePoolSize`. Using `maximumPoolSize`
greater than `corePoolSize` in such case will not have any impact on the maximum
bound of pool size.

Bad:

```
new ThreadPoolExecutor(
 /* corePoolSize= */ 1,
 /* maximumPoolSize= */ 10,
 /* keepAliveTime= */ 60,
 TimeUnit.SECONDS,
 new LinkedBlockingQueue<>());
```
Good:

```
new ThreadPoolExecutor(
 /* corePoolSize= */ 10,
 /* maximumPoolSize= */ 10,
 /* keepAliveTime= */ 60,
 TimeUnit.SECONDS,
 new LinkedBlockingQueue<>());
```
```
new ThreadPoolExecutor(
 /* corePoolSize= */ 1,
 /* maximumPoolSize= */ 10,
 /* keepAliveTime= */ 60,
 TimeUnit.SECONDS,
 new LinkedBlockingQueue<>(QUEUE_CAPACITY));
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ErroneousThreadPoolConstructorChecker")` to the enclosing element.

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
EscapedEntity
HTML entities in @code/@literal tags will appear literally in the rendered javadoc.
Severity
WARNING

*Alternate names: PreconditionsExpensiveString*

Lenient format strings, such as those accepted by `Preconditions`, are often
constructed lazily. The message is rarely needed, so it should either be cheap
to construct or constructed only when needed. This check ensures that these
messages are not constructed using expensive methods that are evaluated eagerly.

Prefer this:

```
checkNotNull(foo, "hello %s", name);
```
instead of this:

```
checkNotNull(foo, String.format("hello %s", name));
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ExpensiveLenientFormatString")` to the enclosing element.

Classes referenced in signatures of non-private APIs should usually also be non-private.

For example, consider:

```
class Foo {
 private record Person(String name) {}
 public static void printMessage(Person person) {
 System.out.printf("hello %s\n", person.name());
 }
}
```
Because `Person` is `private`, clients outside the current compilation unit will
be unable to create an instance of it to pass to `printMessage`.

Prefer to make classes referenced in APIs at least as visible as the API, or
else if the API is only intended to be used in the current compilation unit it
should also be `private`.

That is, prefer:

```
class Foo {
 public record Person(String name) {}
 public static void printMessage(Person person) {
 System.out.printf("hello %s\n", person.name());
 }
}
```
or:

```
class Foo {
 private record Person(String name) {}
 private static void printMessage(Person person) {
 System.out.printf("hello %s\n", person.name());
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ExposedPrivateType")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("ExtendingJUnitAssert") to the enclosing element.

@SuppressWarnings("ExtendingJUnitAssert")

`T extends Object` is redundant; both `<T>` and `<T extends Object>` compile to
identical class files.— unless you are using the Checker Framework.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ExtendsObject")` to the enclosing element.

*Alternate names: fallthrough*

The Google Java Style Guide §4.8.4.2 requires that within a switch
block, each statement group either terminates abruptly (with a `break`,
`continue`, `return` or `throw` statement), or is marked with a comment to
indicate that execution will or might continue into the next statement group.
This special comment is not required in the last statement group of the switch
block.

Example:

```
switch (input) {
 case 1:
 case 2:
 prepareOneOrTwo();
 // fall through
 case 3:
 handleOneTwoOrThree();
 break;
 default:
 handleLargeNumber(input);
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("FallThrough")` to the enclosing element.

Do not override `Object.finalize`.

Starting in JDK 18, the method is deprecated for removal, see JEP 421: Deprecate Finalization for Removal.

The Google Java Style Guide §6.4 states:

It is extremely rare to override

`Object.finalize`.Tip: Don’t do it. If you absolutely must, first read and understand Effective Java Item 8, “Avoid finalizers and cleaners” very carefully, and then don’t do it.

Suppress false positives by adding the suppression annotation to the enclosing element:

```
@SuppressWarnings("Finalize") // TODO(user): remove overrides of finalize
```

*Alternate names: finally, ThrowFromFinallyBlock*

Terminating a finally block abruptly preempts the outcome of the try and catch blocks, and will cause the result of any previously executed return or throw statements to be ignored. Finally blocks should be written so they always complete normally.

Consider the following code. In the case where `doWork` throws `SomeException`,
the finally block will still be executed. If closing the input stream *also*
fails, then the exception that was thrown in the catch block will be preempted
by the exception thrown by `close()`, and the first exception will be lost.

```
InputStream in = openInputStream();
try {
 doWork(in);
} catch (SomeException e) {
 throw new SomeError(e);
} finally {
 in.close(); // exception could be thrown here
}
```
This code is easily fixed using try-with-resources. Below, the input stream will
always be closed, and if `doWork` fails and an `IOException` is thrown, then the
try-with-resources uses the `Throwable.addSuppressed()` method added in Java 7
to propagate both exceptions back to the caller.

```
try (InputStream in = openInputStream()) {
 doWork(in);
} catch (SomeException e) {
 throw new SomeError(e);
}
```
If Java 7 is not available, we recommend Guava’s Closer API.

```
Closer closer = Closer.create();
try {
 InputStream in = closer.register(openInputStream());
 doWork(in);
} catch (Throwable e) {
 throw closer.rethrow(e);
} finally {
 closer.close();
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("Finally")` to the enclosing element.

Casts have higher precedence than binary expressions, so `(int) 0.5f * 100` is
equivalent to `((int) 0.5f) * 100` = `0 * 100` = `0`, not `(int) (0.5f * 100)` =
`50`.

To avoid this common source of error, add explicit parentheses to make the precedence explicit. For example, instead of this:

```
long SIZE_IN_GB = (long) 1.5 * 1024 * 1024 * 1024; // this is 1GB, not 1.5GB!
// this is 0, not a long in the range [0, 1000000000]
long rand = (long) new Random().nextDouble() * 1000000000
```
Prefer:

```
long SIZE_IN_GB = (long) (1.5 * 1024 * 1024 * 1024);
long rand = (long) (new Random().nextDouble() * 1000000000);
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("FloatCast")` to the enclosing element.

Both JUnit and Truth allow for asserting equality of floating point numbers with an absolute tolerance. For example, the following statements are equivalent,

```
double EPSILON = 1e-20;
assertThat(actualValue).isWithin(EPSILON).of(Math.PI);
assertEquals(Math.PI, actualValue, EPSILON);
```
What’s not immediately obvious is that both of these assertions are checking
exact equality between `Math.PI` and `actualValue`, because the next `double`
after `Math.PI` is `Math.PI + 4.44e-16`.

This means that using the same tolerance to compare several floating point values with different magnitude can be prone to error,

```
float TOLERANCE = 1e-5f;
assertThat(pressure).isWithin(TOLERANCE).of(1f); // GOOD
assertThat(pressure).isWithin(TOLERANCE).of(10f); // GOOD
assertThat(pressure).isWithin(TOLERANCE).of(100f); // BAD -- misleading equals check
```
A larger tolerance should be used if the goal of the test is to allow for some
floating point errors, or, if not, `isEqualTo` makes the intention more clear.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("FloatingPointAssertionWithinEpsilon")` to the enclosing element.

`double` and `float` literals that can’t be precisely represented should be
avoided.

Example:

```
double d = 1.9999999999999999999999999999999;
System.err.println(d); // prints 2.0
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("FloatingPointLiteralPrecision")` to the enclosing element.

Prefer to let Flogger transform your arguments to strings, instead of calling
`toString()` explicitly.

For example, prefer the following:

```
logger.atInfo().log("hello '%s'", world);
```
instead of this, which eagerly calls `world.toString()` even if `INFO` level
logging is disabled.

```
logger.atInfo().log("hello '%s'", world.toString());
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("FloggerArgumentToString")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("FloggerPerWithoutRateLimit") to the enclosing element.

@SuppressWarnings("FloggerPerWithoutRateLimit")

Prefer string formatting to concatenating format arguments together, to avoid work at the log site.

That is, prefer this:

```
logger.atInfo().log("processing: %s", request);
```
to this, which calls `request.toString()` even if `INFO` logging is disabled:

```
logger.atInfo().log("processing: " + request);
```
More information: https://google.github.io/flogger/formatting

Suppress false positives by adding the suppression annotation `@SuppressWarnings("FloggerStringConcatenation")` to the enclosing element.

It usually hurts performance to eagerly generate error messages with +, as you pay the cost of the string conversion whether or not the condition fails. It’s usually more efficient to use %s as a placeholder and to pass the additional variables as further arguments.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("FormatStringShouldUsePlaceholders")` to the enclosing element.

#
 FragmentInjection

 Classes extending PreferenceActivity must implement isValidFragment such that it does not unconditionally return true to prevent vulnerability to fragment injection attacks.

 - Severity
- WARNING
- Tags
- LikelyError

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("FragmentInjection")` to the enclosing element.

#
 FragmentNotInstantiable

 Subclasses of Fragment must be instantiable via Class#newInstance(): the class must be public, static and have a public nullary constructor

 - Severity
- WARNING
- Tags
- LikelyError

*Alternate names: ValidFragment*

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("FragmentNotInstantiable")` to the enclosing element.

Methods that return `java.util.concurrent.Future` and its subclasses generally
indicate errors by returning a future that eventually fails.

If you don’t check the return value of these methods, you will never find out if they threw an exception.

Nested futures can also result in missed cancellation signals or suppressed exceptions - see Avoiding Nested Futures for details.

In certain scenarios like tests, there might be a need of not using the future
values. One can suppress such false positives by either suppressing the check
directly or by saving the future in variables named with prefix `unused`. For
example:

```
@SuppressWarnings("FutureReturnValueIgnored")
@Test
public void futureInvocation_noMemoryLeak() {
 functionReturningFuture();
}
```
```
@Test
public void futureInvocation_noMemoryLeak() {
 Future<?> unusedFuture = functionReturningFuture();
}
```

The usage of `transformAsync`, `callAsync` and `submitAsync` is not necessary
when all the return values of the transformation function are immediate futures.
In this case, the usage of `transform`, `call` and `submit` is preferred.

Note that `transform` cannot be used if the body of the transformation function
throws checked exceptions.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("FutureTransformAsync")` to the enclosing element.

Enum values that declare methods are a subclass of the actual enum type, so
calling `getClass()` returns a synthetic subclass of the enum. To retrieve the
type of the enum, use `getDeclaringClass()`.

In the following example, `Binop.MULT.getClass()` returns the anonymous class
`Binop$2`, while `Binop.MULT.getDeclaringClass()` returns the class `Binop`.

```
enum Binop {
 MULT {
 @Override
 int apply(int x) {
 return x * x;
 }
 },
 ADD {
 @Override
 int apply(int x) {
 return x + x;
 }
 };
 abstract int apply(int x);
}
```
```
public class Test {
 static void printEnumClass(Enum theEnum) {
 System.err.println(theEnum.getClass());
 System.err.println(theEnum.getDeclaringClass());
 }
 public static void main(String[] args) {
 printEnumClass(Binop.ADD);
 }
}
```
Prints:

```
class Binop$2
class Binop
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("GetClassOnEnum")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("GuiceNestedCombine") to the enclosing element.

@SuppressWarnings("GuiceNestedCombine")

*Alternate names: hiding, OvershadowingSubclassFields*

If a class has a field of the same name as any field visible to it on any of its superclasses or superinterfaces, the subclass’ field is said to “hide” the superclass’ field.

When this circumstance occurs, users of the class declaring the hiding field can’t interact with the fields from the superclass.

Let’s take a look at how field hiding might cause problems:

```
class Super {
 public String foo = "bar";
}
class Sub extends Super {
 private int foo = 0; // the same name, so this hides `Super`'s `foo`
}
class Main {
 void stringFn(String s) { /*...*/ }
 public static void main(String... args) {
 // Looking at the API of `Super`, I should be able to access a string `foo`
 // on any object of type `Super` or its subclasses, right?
 stringFn(new Sub().foo); // Oops! `foo` is not visible, and the wrong type!
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("HidingField")` to the enclosing element.

`ICC_Profile.getInstance(String)` searches the entire classpath, which is often
unnecessary and can result in slow performance for applications with long
classpaths. Prefer `getInstance(byte[])` or `getInstance(InputStream)` instead.

See also https://bugs.openjdk.org/browse/JDK-8191622.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ICCProfileGetInstance")` to the enclosing element.

`java.util.IdentityHashMap` uses reference equality to compare keys. This is
in violation of the contract of `java.util.Map`,
which states that object equality (the keys’ `equals` methods) should be used
for key comparison. This peculiarity can lead to confusion and subtle bugs,
especially when the two types of maps are used together. This check attempts to
reduce confusion between these two very different kinds of maps in a few ways:

`IdentityHashMap`’s reference-equality behavior would be surprising if
its declared type was `java.util.Map`. The type of an `IdentityHashMap` is
required to be explicit, i.e., `IdentityHashMap` should not be upcast to
`Map`.`IdentityHashMap`’s `equals` method will only compare equal to a `Map` if
the keys of both are the same instances. This makes little sense for an
object-equality map, so `IdentityHashMap.equals()` should only be used to
compare to other `IdentityHashMap`s.`Map` and `IdentityHashMap` are different enough that they
should be considered different, incompatible types. Converting one to the
other is discouraged, as it is often a mistake.```
Map<String, String> bad(Map<String, String> aMap, IdentityHashMap<String, String> identityMap) {
 // Don't assign an `IdentityHashMap` to a plain `Map`-typed variable
 Map<String, String> myMap = identityMap;
 // Don't use `IdentityHashMap`'s `equals` method with a `Map`
 identityMap.equals(aMap);
 // Don't convert between `IdentityHashMap` and `Map`
 identityMap.putAll(aMap);
 identityMap = new IdentityHashMap<>(aMap);
 return identityMap;
}
```
```
// Keep `IdentityHashMap`'s type information around so maintainers know when
// reference equality is being used.
IdentityHashMap<String, String> good() {
 IdentityHashMap<String, String> identityMap = new IdentityHashMap<>();
 // ...
 return identityMap;
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("IdentityHashMapUsage")` to the enclosing element.

We’re trying to make long chains of `if` statements clearer (and potentially
faster) by converting them into `switch`es.

`if` statements`if ... else if ... else if ...` statements has `K` total
branches, at runtime one needs to check `O(K)` conditions on average
(assuming equal likelihood of each branch)`if (...)`. Besides redundancy, this introduces a potential bug vector:
an `if` in the chain could unintentionally have a slightly different
condition than others, an ordering bug (see below), `switch`es:`1`, `2`, …), `enum` values, `null`, and pattern
matching, including mixtures of these```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private void foo(Suit suit) {
 if (suit == Suit.SPADE) {
 System.out.println("spade");
 } else if (suit == Suit.HEART || suit == Suit.DIAMOND) {
 System.out.println("red suit");
 } else if (suit == Suit.CLUB) {
 System.out.println("club");
 }
}
```
Which can be converted into:

```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private void foo(Suit suit) {
 switch (suit) {
 case Suit.SPADE -> System.out.println("spade");
 case Suit.HEART, Suit.DIAMOND -> System.out.println("red suit");
 case Suit.CLUB -> System.out.println("club");
 }
}
```
If the flag `-XepOpt:IfChainToSwitch:EnableSafe=true` is set, the output will
include an empty `case null`; this more closely matches the behavior of the
original if-chain when `suit` is `null`, although is more verbose and may not
match the intent.

```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private void foo(Suit suit) {
 switch (suit) {
 case Suit.SPADE -> System.out.println("spade");
 case Suit.HEART, Suit.DIAMOND -> System.out.println("red suit");
 case Suit.CLUB -> System.out.println("club");
 case null -> {}
 }
}
```
Note that with the new `switch` style (`->`), exhaustiveness checking is
provided by the compiler for ‘enhanced’ switch statements (see JLS §14.11.2),
and for traditional switch statements by Error Prone
(https://errorprone.info/bugpattern/MissingCasesInEnumSwitch). That is, if a new `Suit` value were to
be added to the `enum`, then the `switch` would raise a compile-time error,
whereas the original chain of `if` statements would need to be manually detected
and edited.

This conversion works for `instanceof`s too:

```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private void describeObject(Object obj) {
 if (obj instanceof String) {
 System.out.println("It's a string!");
 } else if (obj instanceof Number n) {
 System.out.println("It's a number!");
 } else if (obj instanceof Object) {
 System.out.println("It's an object!");
 }
}
```
This can be converted as follows (if the `instanceof` does not originally have a
pattern variable, then `unused` will be inserted):

```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private void describeObject(Object obj) {
 switch(obj) {
 case String unused -> System.out.println("It's a string!");
 case Number n -> System.out.println("It's a number!");
 case Object unused -> System.out.println("It's an object!");
 }
}
```
In later Java versions, an unnamed variable (`_`) can be used in place of
`unused`.

With `if` chains, it’s possible to write code such as:

```
private void describeObject(Object obj) {
 if (obj instanceof Object) {
 System.out.println("It's an object!");
 } else if (obj instanceof Number n) {
 System.out.println("It's a number!");
 } else if (obj instanceof String) {
 System.out.println("It's a string!");
 }
}
```
When calling `describeObject("hello")`, one might expect to have ```
It's a
string!
```
 printed, but this is not what happens. Because the `Object` check
happens first in code, it matches, resulting in `It's an object!`. This behavior
is most likely a bug, and can sometimes be hard to spot. This checker will
automatically reorder `case`s if needed to correct the issue, like this:

```
private void describeObject(Object obj) {
 switch(obj) {
 case Number n -> System.out.println("It's a number!");
 case String unused -> System.out.println("It's a string!");
 case Object unused -> System.out.println("It's an object!");
 }
}
```
In this way, the behavior of the switch (`It's a string!`) and the original
if-chain (`It's an object!`) are different. To prevent the checker from changing
behavior in this way, set the flag `-XepOpt:IfChainToSwitch:EnableSafe=true`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("IfChainToSwitch")` to the enclosing element.

IgnoredPureGetter
Getters on AutoValues, AutoBuilders, and Protobuf Messages are side-effect free, so there is no point in calling them if the return value is ignored. While there are no side effects from the getter, the receiver may have side effects.

Severity

WARNING

Suppression

Suppress false positives by adding the suppression annotation @SuppressWarnings("IgnoredPureGetter") to the enclosing element.

*Alternate names: Immutable*

Annotations should always be immutable.

Static state is dangerous to begin with, but much worse for annotations. Annotation instances are usually constants, and it is very surprising if their state ever changes, or if they are not thread-safe.

TIP: prefer [`@AutoAnnotation`] to writing annotation implementations by hand.

To make annotation implementations immutable, ensure:

`ImmutableList` and `ImmutableSet` instead of `List` and `Set`.`com.google.errorprone.annotations.Immutable` to its declaration.TIP: annotating the declaration of an annotation with `@Immutable` is
unnecessary – Error Prone assumes annotations are immutable by default.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ImmutableAnnotationChecker")` to the enclosing element.

*Alternate names: Immutable*

All fields in your enum class should be final and either be primitive or refer to deeply immutable objects.

Static state is dangerous to begin with, but much worse for enums. We all think of enum values as constants – and even refer to them as “enum constants” – and would be very surprised if any of their state ever changed, or was not thread-safe.

To make enums immutable, ensure:

The types of all fields of the enum are deeply immutable. For example, use
`ImmutableList` and `ImmutableSet` instead of `List` and `Set`.

Types are considered immutable if they are primitives, in a set of types
that are built in to Error Prone (e.g. `java.lang.String`,
`java.util.UUID`), or are annotated with
`com.google.errorprone.annotations.Immutable`.

Other versions of the annotation, such as
`javax.annotation.concurrent.Immutable`, are currently *not* recognized.
See https://errorprone.info/bugpattern/Immutable for additional discussion.

If the type you’re using inside the enum can be annotated with
`@Immutable`, you should do that:

```
// WARNING: E is not immutable, since MyValueObject isn't Immutable
enum E {
 private final MyValueObject mvo = new MyValueObject();
}
// Add @Immutable here
final class MyValueObject {}
```
Note that MyValueObject must actually be immutable. The
Immutable checker will raise a compile-time error if the
`MyValueObject` class isn’t actually immutable.

If the type you’re using inside the enum is not considered immutable, and you can’t annotate the type because it’s outside the project, consider using an immutable replacement of the type, or suppress this check on the enum with a comment about why the fields in question are immutable.

TIP: annotating the declaration of the enum class with `@Immutable` is
unnecessary – Error Prone assumes enums are immutable by default.

Example:

```
import com.google.errorprone.annotations.Immutable;
@Immutable
class Foo {
 final int id;
 public Foo(int id) {
 this.id = id;
 }
}
```
```
// The declaration doesn't need to be annotated with @Immutable.
enum E {
 A("A", ImmutableList.of(new Foo(1), new Foo(2))),
 B("B", ImmutableList.of(new Foo(3)));
 // All fields are final, and deeply immutable:
 // Error Prone knows String is immutable.
 private final String name;
 // Error Prone knows ImmutableList<T> is an immutable collection of some
 // objects, and it recognizes the @Immutable annotation on the declaration of
 // Foo, so it can safely determine that this ImmutableList is deeply
 // immutable.
 private final ImmutableList<Foo> foos;
 private E(String name, ImmutableList<Foo> foos) {
 this.name = name;
 this.foos = foos;
 }
 public ImmutableList<Foo> foos() {
 return foos;
 }
 public String name() {
 return foos;
 }
}
```
TIP: Instead of creating an enum with functional interface fields (`Predicate`,
`Function`, etc.), declare abstract methods that are overridden by each
constant. For example, do this:

```
enum Types {
 STRING {
 @Override public boolean hasCompatibleType(Object o) {
 return o instanceof String;
 }
 },
 NUMBER {
 @Override public boolean hasCompatibleType(Object o) {
 return o instanceof Number;
 }
 },
 // ...
 public abstract boolean hasCompatibleType(Object o);
}
```
… not this:

```
enum Types {
 STRING(o -> o instanceof String),
 NUMBER(o -> o instanceof Number),
 // ...
 final Predicate<Object> hasCompatibleType;
 Types(Predicate<Object> hasCompatibleType) {
 this.hasCompatibleType = hasCompatibleType;
 }
}
```
This has several advantages on top of sidestepping this checker, e.g. not tying
you to a particular functional interface type – your callers should e.g. use
`STRING::hasCompatibleType` instead of `STRING.hasCompatibleType` which only
works for one interface type.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ImmutableEnumChecker")` to the enclosing element.

It is confusing to have two or more variables under the same scope that differ only in capitalization. Make sure that both of these follow the casing guide (Google Java Style Guide §5.3) and to be consistent if more than one option is possible.

This checker will only find parameters that differ in capitalization with fields that can be accessed from the parameter’s scope.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InconsistentCapitalization")` to the enclosing element.

Implementations of `Object#hashCode` should not incorporate fields which the
implementation of `Object#equals` does not. This violates the contract of
`hashCode`: specifically, equal objects must have equal hashCodes.

```
class Foo {
 private final int a;
 private final int b;
 Foo(int a, int b) {
 this.a = a;
 this.b = b;
 }
 @Override
 public boolean equals(@Nullable Object o) {
 return o instanceof Foo && ((Foo) o).a == a;
 }
 @Override
 public int hashCode() {
 return a + 31 * b;
 }
}
Foo first = new Foo(10, 20);
Foo second = new Foo(10, 40);
first.equals(second) // true
first.hashCode() == second.hashCode() // false
```
The fix for this class is either to include a comparison of `b` in the `#equals`
method, or remove `b` from `#hashCode`. The former is more likely to be correct.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InconsistentHashCode")` to the enclosing element.

Prior to Java 25, a `main` method must be `public`, `static`, and return `void`
(see JLS §12.1.4).

For example, the following method is confusing, because it is an overload of a
valid `main` method (it has the same name and signature), but is not a valid
`main` method:

```
class Test {
 static void main(String[] args) {
 System.err.println("hello world");
 }
}
```
```
$ java T.java
error: 'main' method is not declared 'public static'
```
For Java 25 and later, a `main` method must return `void` (see JLS §12.1.4).
The `public` and `static` requirements have been dropped.

For example, the following method is confusing, because it is an overload of a
valid `main` method (it has the same name and arguments), but does not return
`void`:

```
class Test {
 public static int main(String[] args) {
 System.err.println("hello world");
 return 0;
 }
}
```
TIP: If you’re declaring a method that isn’t intended to be used as the main
method of your program, prefer to use a name other than `main`. It’s confusing
to humans and static analysis to see methods like ```
private int main(String[]
args)
```
.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("IncorrectMainMethod")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("IncrementInForLoopAndHeader") to the enclosing element.

@SuppressWarnings("IncrementInForLoopAndHeader")

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
InheritDoc
Invalid use of @inheritDoc.
Severity
WARNING
Tags
Style

`@Scope` annotations should be applicable to TYPE (annotating classes that
should be scoped) and to METHOD (annotating `@Provides` methods to apply scoping
to the returned object.

If an annotation’s use is restricted by `@Target` and it doesn’t include those
two element types, the annotation can’t be used where it should be able to be
used.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InjectInvalidTargetingOnScopingAnnotation")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("InjectOnBugCheckers") to the enclosing element.

@SuppressWarnings("InjectOnBugCheckers")

When dependency injection frameworks call constructors, they can only do so on
constructors of concrete classes, which can delegate to superclass constructors.
In the case of abstract classes, their constructors are only called by their
concrete subclasses, not directly by injection frameworks, so the `@Inject`
annotation has no effect.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InjectOnConstructorOfAbstractClass")` to the enclosing element.

Scoping annotations are not inherited. They will have no effect if added to an abstract type.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InjectScopeAnnotationOnInterfaceOrAbstractClass")` to the enclosing element.

The constructor is annotated with @Inject(optional=true), or it is annotated with @Inject and a binding annotation. This will cause a Guice runtime error.

For more information, see https://github.com/google/guice/wiki/InjectionPoints.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InjectedConstructorAnnotations")` to the enclosing element.

Prefer to use literal format strings directly in the call to a formatting method over extracting them to a constant. That is, prefer this:

```
throw new IllegalArgumentException(
 String.format("Uh oh, can't use %s with %s", badArgA, badArgB));
```
to this:

```
private static final String ERROR_MESSAGE = "Uh oh, can't use %s with %s";
...
throw new IllegalArgumentException(String.format(ERROR_MESSAGE, badArgA, badArgB));
```
Extracting the format string to a constant makes it harder to read the
`String.format` call, and to see that the correct format arguments are being
passed.

If a single format string is used by multiple calls to `String.format`, consider
extracting a helper method instead of making the string a constant:

```
String errorMessage(String badArgA, String badArgB) {
 return String.format("Uh oh, can't use %s with %s", badArgA, badArgB);
}
...
throw new IllegalArgumentException(errorMessage(badArgA, badArgB));
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("InlineFormatString")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("InlineMeInliner") to the enclosing element.

@SuppressWarnings("InlineMeInliner")

#
 InlineMeSuggester

 This deprecated API looks inlineable. If you'd like the body of the API to be automatically inlined to its callers, please annotate it with @InlineMe. NOTE: the suggested fix makes the method final if it was not already.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InlineMeSuggester")` to the enclosing element.

Constants should be given names that emphasize the semantic meaning of the value. If the name of the constant doesn’t convey any information that isn’t clear from the value, consider inlining it.

For example, prefer this:

```
System.err.println(1);
System.err.println("");
```
to this:

```
private static final int ONE = 1;
private static final String EMPTY_STRING = "";
...
System.err.println(ONE);
System.err.println(EMPTY_STRING);
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("InlineTrivialConstant")` to the enclosing element.

`java.io.InputStream` defines a single abstract method: `int read()`, which
subclasses implement to return bytes from the logical input stream.

However, in most circumstances, readers from `InputStreams` use higher-level
methods like `read(byte[], int offset, int length)` to read multiple bytes at a
time into a buffer. The default implementation of this method is to repeatedly
call `read()`. However, most InputStream implementations could do much better if
they can read multiple bytes at once (at the very least, avoiding unneeded
`byte` -> `int` -> `byte` casts that are needed when implementing the read()
method over an underlying `byte` source).

The class in question implements `int read()` without also overriding ```
int
read(byte[], int, int)
```
 and will thus be subject to the costs associated with
the default behavior of the multibyte read method.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InputStreamSlowMultibyteRead")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("InstanceOfAndCastMatchWrongType") to the enclosing element.

@SuppressWarnings("InstanceOfAndCastMatchWrongType")

Implicit conversions from `int` to `float` may lose precision when when calling
methods with overloads that accept both`float` and `double`.

For example, `Math.scalb` has overloads
`Math.scalb(float, int)`
and
`Math.scalb(double, int)`.

When passing an `int` as the first argument, `Math.scalb(float, int)` will be
selected. If the result of `Math.scalb(float, int)` is then used as a `double`,
this may result in a loss of precision.

To avoid this, an explicit cast to `double` can be used to call the
`Match.scalb(double, int)` overload:

```
int x = ...
int y = ...
double f = Math.scalb((double) x, 2);
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("IntFloatConversion")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("IntLiteralCast") to the enclosing element.

@SuppressWarnings("IntLiteralCast")

Performing an arithmetic expression on arguments of type int and then assigning the result to a long is error-prone. The result is widened to a long as the final step, and the intermediate results may overflow.

For example, the following expression exceeds `Integer.MAX_VALUE` and overflows
to `-1857093632`:

```
long nanosPerDay = 24 * 60 * 60 * 1000 * 1000 * 1000;
```
The corrected code (which has a value of `86400000000000`) is:

```
long nanosPerDay = 24L * 60 * 60 * 1000 * 1000 * 1000;
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("IntLongMath")` to the enclosing element.

When attempting to fail-fast when interrupted, you should preserve the
interrupted status of the thread by calling
`Thread.currentThread().interrupt()`:

```
try {
 mightTimeOutOrBeCancelled(); // for example myFuture.get()
} catch (InterruptedException e) {
 Thread.currentThread().interrupt(); // Restore the interrupted status
 throw new MyCheckedException("[describe what task was interrupted]", e);
}
```
The current code is likely accidentally calling `Thread.interrupted()`, which
clears the interrupt bit. (That is typically a no-op in this situation because
the interrupt bit is generally already unset when `InterruptedException` has
been thrown.) `thread.interrupt()` correctly restores the interrupt bit.

More information:

https://web.archive.org/web/20201217182342/https://www.ibm.com/developerworks/java/library/j-jtp05236/index.html

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InterruptedInCatchBlock")` to the enclosing element.

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
InvalidBlockTag
This tag is invalid.
Severity
WARNING

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
InvalidInlineTag
This tag is invalid.
Severity
WARNING

This error is triggered by a Javadoc `@link` tag that either is syntactically
invalid or can’t be resolved. See javadoc documentation for an
explanation of how to correctly format the contents of this tag.

Use the erased type of method parameters in `@link` tags. For example, write
`{@link #foo(List)}` instead of `{@link #foo(List<Bah>)}`. Javadoc does yet not
support generics in `@link` tags, due to a bug:
JDK-5096551.

This check is very limited in terms of which unresolved links it can be *sure*
are unresolvable. Code within Google is often compiled on a per-package or even
per-file basis, and Error Prone only has visibility into the current
compilation.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InvalidLink")` to the enclosing element.

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
InvalidParam
This @param tag doesn't refer to a parameter of the method.
Severity
WARNING

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
InvalidSnippet
This tag is invalid.
Severity
WARNING

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
InvalidThrows
The documented method doesn't actually throw this checked exception.
Severity
WARNING

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
InvalidThrowsLink
Don't use {@link} or {@code} in @throws tags; mention the exception name directly (e.g., @throws IOException, not @throws {@link IOException}).
Severity
WARNING
Tags
Style

An `Iterator` is a *state-ful* instance that enables you to check whether it has
more elements (via `hasNext()`) and moves to the next one if any (via `next()`),
while an `Iterable` is a representation of literally iterable elements. An
`Iterable` can generate multiple valid `Iterator`s, though.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("IterableAndIterator")` to the enclosing element.

JUnit 4’s floating-point overloads of `assertEquals(expected, actual)` always
throw an exception, and some floating-point calls to JUnit 4’s `assertEquals` do
not even compile.

To continue comparing floating-point numbers using `Double.equals` semantics,
you may be able to cast one argument to `Object` or use Truth’s
`assertThat(actual).isEqualTo(expected)` /
`assertWithMessage(message).that(actual).isEqualTo(expected)`.

Alternatively, you can switch to tolerance-based equality testing, which changes
your code’s behavior for negative zero (in JUnit and Truth) and for infinities
and NaN (in Truth). If you want that, use JUnit’s ```
assertEquals(expected,
actual, delta)
```
 or Truth’s `isWithin(...).of(...)`, possibly with a tolerance of
zero.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JUnit3FloatingPointComparisonWithoutDelta")` to the enclosing element.

#
 JUnit4ClassUsedInJUnit3

 Some JUnit4 construct cannot be used in a JUnit3 context. Convert your class to JUnit4 style to use them.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JUnit4ClassUsedInJUnit3")` to the enclosing element.

#
 JUnit4EmptyMethods

 Empty JUnit4 @Before, @After, @BeforeClass, and @AfterClass methods are unnecessary and should be deleted.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JUnit4EmptyMethods")` to the enclosing element.

For JUnit3-style tests, behavior is defined in `junit.framework.TestCase` and
tests add behavior by overriding methods. For JUnit4-style tests, special
behavior happens with fields and methods annotated with JUnit 4 annotations.
Having JUnit4-style tests extend from `junit.framework.TestCase` (directly or
indirectly) historically has been a source of test bugs and unexpected behavior
(e.g.: teardown logic and/or verification does not run because JUnit doesn’t
call the inherited code).

Error Prone also cannot infer whether the test class runs with JUnit 3 or JUnit 4.

Thus, even if the test class runs with JUnit 4, Error Prone will not run
additional checks which can catch common errors with JUnit 4 test classes.
Either use only JUnit4 classes and annotations and remove the inheritance from
TestCase, or use only JUnit 3 and remove the `@Test` annotations. When looking
for replacements for base test classes, consider using Rules (see the `@Rule`
annotation and implementations of `TestRule` and `MethodRule`).

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JUnitAmbiguousTestClass")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("JUnitIncompatibleType") to the enclosing element.

@SuppressWarnings("JUnitIncompatibleType")

#
 JUnitMethodInvoked

 Directly invoking a JUnit test method is discouraged; only the JUnit test runner should call these methods. If you need to share logic between tests, extract a helper method or class.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JUnitMethodInvoked")` to the enclosing element.

If you call duration.getNano(), you must also call duration.getSeconds() in ‘nearby’ code. If you are trying to convert this duration to nanoseconds, you probably meant to use duration.toNanos() instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JavaDurationGetSecondsGetNano")` to the enclosing element.

duration.getSeconds() is a decomposition API which should always be used alongside duration.getNano(). duration.toSeconds() is a conversion API, and the preferred way to convert to seconds.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JavaDurationGetSecondsToToSeconds")` to the enclosing element.

Duration’s withNanos(int) method is often a source of bugs because it returns a copy of the current Duration instance, but *only* the nano field is mutated (the seconds field is copied directly). Use Duration.ofSeconds(duration.getSeconds(), nanos) instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JavaDurationWithNanos")` to the enclosing element.

Duration’s withSeconds(long) method is often a source of bugs because it returns a copy of the current Duration instance, but *only* the seconds field is mutated (the nanos field is copied directly). Use Duration.ofSeconds(seconds, duration.getNano()) instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JavaDurationWithSeconds")` to the enclosing element.

If you call instant.getNano(), you must also call instant.getEpochSecond() in ‘nearby’ code. If you are trying to convert this instant to nanoseconds, you probably meant to use Instants.toEpochNanos(instant) instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JavaInstantGetSecondsGetNano")` to the enclosing element.

#
 JavaLocalDateTimeGetNano

 localDateTime.getNano() only access the nanos-of-second field. It's rare to only use getNano() without a nearby getSecond() call.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JavaLocalDateTimeGetNano")` to the enclosing element.

#
 JavaLocalTimeGetNano

 localTime.getNano() only accesses the nanos-of-second field. It's rare to only use getNano() without a nearby getSecond() call.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JavaLocalTimeGetNano")` to the enclosing element.

#
 JavaPeriodGetDays

 period.getDays() only accesses the "days" portion of the Period, and doesn't represent the total span of time of the period. Consider using org.threeten.extra.Days to extract the difference between two civil dates if you want the whole time.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JavaPeriodGetDays")` to the enclosing element.

Using APIs that silently use the default system time-zone is dangerous. The default system time-zone can vary from machine to machine or JVM to JVM. You must choose an explicit ZoneId.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JavaTimeDefaultTimeZone")` to the enclosing element.

The `Date` API is full of
major design flaws and pitfalls
and should be avoided at all costs. Prefer the `java.time` APIs, specifically,
`java.time.Instant` (for physical time) and `java.time.LocalDate[Time]` (for
civil time).

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JavaUtilDate")` to the enclosing element.

According to the JSR-330 spec, the @javax.inject.Inject annotation cannot go on final fields.

Suppress false positives by adding the suppression annotation @SuppressWarnings("JavaxInjectOnFinalField") to the enclosing element.

@SuppressWarnings("JavaxInjectOnFinalField")

Some JDK APIs are obsolete and have preferred alternatives.

`LinkedList``LinkedList` almost never out-performs `ArrayList` or `ArrayDeque`1. If you
are using `LinkedList` as a list, prefer `ArrayList`. If you are using
`LinkedList` as a stack or queue/deque, prefer `ArrayDeque`.

Migration gotcha: `LinkedList` permits `null` elements; `ArrayDeque` rejects
them. The documentation for `Deque` strongly discourages users from inserting
`null`, even into implementations that permit it. So, if you are using a
`LinkedList` for this purpose, you should likely stop, and you will *need* to
stop in order to migrate to `ArrayDeque`.

`Vector``Vector` performs synchronization that is usually unnecessary; prefer
`ArrayList`.

If a synchronized collection is necessary, use `Collections.synchronizedList` or
a data structure from `java.util.concurrent`.

`Hashtable` and `Dictionary`This is a nonstandard class that predates the Java Collections Framework; prefer
`LinkedHashMap` or `HashMap`.

If synchronization is necessary, `java.util.concurrent.ConcurrentHashMap` is
usually a good choice.

`java.util.Stack``Stack` is a nonstandard class that predates the Java Collections Framework;
prefer `ArrayDeque`.

If a synchronized collection is necessary, use a data structure from
`java.util.concurrent`.

When migrating from `Stack` to `Deque`, note that the `Stack` methods
`push`/`pop`/`peek`/`add`/`iterator` correspond to the `Deque` methods
`addFirst`/`removeFirst`/`peekFirst`/`addFirst`/`descendingIterator`.

`StringBuffer``StringBuffer` performs synchronization that is rarely necessary and has
significant performance overhead. Prefer `StringBuilder`, which does not do
synchronization.

If synchronization is necessary, consider creating an explicit lock object and
using `synchronized` blocks.

`Enumeration`An ancient precursor to `Iterator`.

`SortedSet` and `SortedMap`Replaced by `NavigableSet` and `NavigableMap` in Java 6.

`Matcher#hitEnd()` and `Matcher#requireEnd()`The methods [`Matcher#hitEnd()`][hitEnd] and
[`Matcher#requireEnd()`][requireEnd] are mainly intended for implementing
scanning functionalities (like `java.util.Scanner`), where you need to decide
whether to read more input from a stream to find a longer match.

If you are not implementing your own scanning or streaming parser, and the input being matched against is already fully loaded in memory, using these methods is likely a mistake.

If you want to check if a match extends to the end of the input, you can simply compare the match end index with the length of the input string:

```
// Instead of:
if (matcher.find() && matcher.hitEnd()) { ... }
// Use:
if (matcher.find() && matcher.end() == input.length()) { ... }
```
If you are implementing a scanner or streaming parser, consider using
`java.util.Scanner` directly, which already encapsulates this logic correctly.

```
Scanner scanner = new Scanner(input);
while (scanner.hasNext(pattern)) {
 String match = scanner.next(pattern);
 ...
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("JdkObsolete")` to the enclosing element.

People generally choose `LinkedList` because they want fast insertion and
removal. However, `LinkedList` has slow traversal, and typically you need
to traverse the list to find the place to insert/remove. It turns out that
the cost of traversing to a location in a `LinkedList` is approximately 4x
the cost of copying an element in an `ArrayList`. Thus, `LinkedList`’s
traversal cost dominates and results in poorer performance. More info:
https://stuartmarks.wordpress.com/2015/12/18/some-java-list-benchmarks/
[hitEnd]: https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/util/regex/Matcher.html#hitEnd()
[requireEnd]: https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/util/regex/Matcher.html#requireEnd() ↩

Use JodaTime’s static factories instead of the ambiguous constructors.

Suppress false positives by adding the suppression annotation @SuppressWarnings("JodaConstructors") to the enclosing element.

@SuppressWarnings("JodaConstructors")

Manual date/time math leads to overflows, unit mismatches, and weak typing. Prefer to use strong types (e.g., `java.time.Duration` or `java.time.Instant`) and their APIs to perform date/time math.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JodaDateTimeConstants")` to the enclosing element.

Joda-Time’s ‘duration.withMillis(long)’ method is often a source of bugs because it doesn’t mutate the current instance but rather returns a new immutable Duration instance. Please use Duration.millis(long) instead. If your Duration is better expressed in terms of other units, use standardSeconds(long), standardMinutes(long), standardHours(long), or standardDays(long) instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JodaDurationWithMillis")` to the enclosing element.

Joda-Time’s ‘instant.withMillis(long)’ method is often a source of bugs because it doesn’t mutate the current instance but rather returns a new immutable Instant instance. Please use Instant.ofEpochMilli(long) instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JodaInstantWithMillis")` to the enclosing element.

JodaNewPeriod
This may have surprising semantics, e.g. new Period(LocalDate.parse("1970-01-01"), LocalDate.parse("1970-02-02")).getDays() == 1, not 32.

Severity

WARNING

Suppression

Suppress false positives by adding the suppression annotation @SuppressWarnings("JodaNewPeriod") to the enclosing element.

JodaTime’s type.plus(long) and type.minus(long) methods are often a source of bugs because the units of the parameters are ambiguous. Please use type.plus(Duration.millis(long)) or type.minus(Duration.millis(long)) instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JodaPlusMinusLong")` to the enclosing element.

#
 JodaTimeConverterManager

 Joda-Time's ConverterManager makes the semantics of DateTime/Instant/etc construction subject to global static state. If you need to define your own converters, use a helper.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JodaTimeConverterManager")` to the enclosing element.

JodaTime’s type.withDurationAdded(long, int) is often a source of bugs because the units of the parameters are ambiguous. Please use type.withDurationAdded(Duration.millis(long), int) instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("JodaWithDurationAddedLong")` to the enclosing element.

There are two overloads of the method `List.remove`:
`remove(int index)` removes an element at the specified index, and
`remove(Object element)` removes the specified element. When used
with a list of integers (`List<Integer>`), the overload resolution can be
confusing and may lead to bugs if the wrong overload is selected or if the
code’s intent is not clear to readers.

Consider the following code:

```
List<Integer> list = new ArrayList<>();
// ...
list.remove(1);
```
In this case, `list.remove(1)` calls `remove(int index)`, removing the element
at index 1. If the intention was to remove the element with value 1, this code
is incorrect.

To make the intent explicit and avoid ambiguity, add a comment before the argument:

If you meant to remove by value (element):

```
list.remove(/* element */ ii);
```
If you meant to remove by index:

```
list.remove(/* index */ 1);
```
If you need to change the type to select the correct overload, you can use explicit boxing/unboxing:

```
list.remove(Integer.valueOf(i)); // remove by element
list.remove(ii.intValue()); // remove by index
```
JDK-8384074 discusses adding a
new `List::removeAtIndex` method to the JDK to avoid this potential overload
confusion.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ListRemoveAmbiguous")` to the enclosing element.

Byte code optimizers can change the implementation of `toString()` in lite
runtime and thus using `valueOf(String)` is discouraged. Instead of converting
enums to string and back, its numeric value should be used instead as it is the
stable part of the protocol defined by the enum.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("LiteEnumValueOf")` to the enclosing element.

After bytecode optimization, Lite protos do not generate a meaningful
`#toString()`. This can be acceptable for debugging, as development builds will
typically not be optimized, but can pose a problem if the default lite
`#toString` finds its way into production code as part of an important error
message.

```
 public void validate(MyProto myProto) {
 if (!myProto.hasFrobnicator()) {
 // MyProto missing frobnicator: asf@6531767b
 throw new IllegalArgumentException("MyProto missing frobnicator: " + myProto);
 }
 }
```
The fix will be highly dependent on the case at hand. There may be an identifier associated with the message that would be useful to log, or even some serialized version of the entire proto.

NOTE: Logging fields of a proto will force those fields to be retained after
optimization. This is not an issue if they’re already being used, but writing a
comprehensive `#toString` method for a proto which is otherwise largely unused
would increase code size.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("LiteProtoToString")` to the enclosing element.

Wherever possible, calls to Lock#lock should be immediately followed by a `try`
block with a `finally` clause which releases the lock,

```
lock.lock();
try {
 frobnicate();
} finally {
 lock.unlock();
}
```
Placing the call to `lock` *inside* the `try` block is suboptimal. The
documentation for `Lock` allows for `Lock#lock` to throw an unchecked exception
if, for example, it determines that the program will deadlock. In this
situation, `#unlock` will be called even if `#lock` throws,

```
try {
 lock.lock();
 frobnicate();
} finally {
 lock.unlock();
}
```
Doing work between the `lock` invocation and the start of the `try` block is
potentially very bad, as the lock will go unreleased if the intermediate work
throws,

```
lock.lock();
checkState(frobnicator.ready());
try {
 frobnicator.frobnicate();
} finally {
 lock.unlock();
}
```
`frobnicator.ready()` should either be moved before the `#lock` call or inside
the `try` block, depending on whether the lock must be acquired before calling
it.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("LockNotBeforeTry")` to the enclosing element.

Having a lock on a class, other than the enclosing class of the code block, can
unintentionally prevent the locked class from being used properly when other
classes effectively lock on the same resource. From a maintainability
perspective, it can be time-consuming to ensure the `synchronized` blocks are
working as expected. Hence, locking on a class other than the enclosing class of
the `synchronized` code block is discouraged by Error Prone. Locking on the
enclosing class or an instance is a preferred practice.

For example, a lock on `Other.class` rather than `Example.class` will trigger an
Error Prone warning:

```
class Example {
 void method() {
 synchronized (Other.class) {
 }
 }
}
```
A lock on the instance or the enclosing class of the `synchronized` code block
will not trigger the warning:

```
class Example {
 void method() {
 synchronized (Example.class) {
 }
 synchronized (this) {
 }
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("LockOnNonEnclosingClassLiteral")` to the enclosing element.

When an assignment expression is used as the condition of a loop, it isn’t clear to the reader whether the assignment was deliberate or it was intended to be an equality test. Parenthesis should be used around assignments in loop conditions to make it clear to the reader that the assignment is deliberate.

That is, instead of this:

```
void f(boolean x) {
 while (x = checkSomething()) {
 // ...
 }
}
```
Prefer `while ((x = checkSomething())) {` or `while (x == checkSomething()) {`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("LogicalAssignment")` to the enclosing element.

A cast from `long` to `double` may lose precision. Prefer an explicit cast to an
implicit conversion if this was intentional.

Consider
`com.google.protobuf.util.Values`
which has a method `of(double value)`, and the following example:

```
// Values.of receives a long, which is implicitly converted to double:
long value = 123L;
Values.of(value);
```
Prefer this (to make existing behavior explicit):

```
long value = 123L;
Values.of((double) value);
```
From JLS §5.1.2:

A widening primitive conversion from

`int`to`float`, or from`long`to`float`, or from`long`to`double`, may result in loss of precision - that is, the result may lose some of the least significant bits of the value. In this case, the resulting floating-point value will be a correctly rounded version of the integer value, using IEEE 754 round-to-nearest mode

Suppress false positives by adding the suppression annotation `@SuppressWarnings("LongDoubleConversion")` to the enclosing element.

A cast from `long` to `float` may lose precision. Prefer an explicit cast to an
implicit conversion if this was intentional.

Consider
`java.awt.Color`
which has constructors `Color(float r, float g, float b)` and ```
Color(int r, int
g, int b)
```
, and the following example:

```
// Math.round returns a double, which is implicitly converted to float:
new Color(Math.round(18.0), Math.round(0.0), Math.round(18.0));
```
Prefer this (to make existing behavior explicit):

```
new Color((float) Math.round(18.0), (float) Math.round(0.0), (float) Math.round(18.0));
```
or this (if this implicit conversion to `float` was unintentional):

```
new Color((int) Math.round(18.0), (int) Math.round(0.0), (int) Math.round(18.0));
```
From JLS §5.1.2:

A widening primitive conversion from

`int`to`float`, or from`long`to`float`, or from`long`to`double`, may result in loss of precision - that is, the result may lose some of the least significant bits of the value. In this case, the resulting floating-point value will be a correctly rounded version of the integer value, using IEEE 754 round-to-nearest mode

Suppress false positives by adding the suppression annotation `@SuppressWarnings("LongFloatConversion")` to the enclosing element.

`String#toCharArray`
allocates a new array. Calling `charAt` is more efficient, because it avoids
creating a new array with a copy of the character data.

That is, prefer this:

```
boolean isDigits(String string) {
 for (int i = 0; i < string.length(); i++) {
 char c = string.charAt(i);
 if (!Character.isDigit(c)) {
 return false;
 }
 }
 return true;
}
```
to this:

```
boolean isDigits(String string) {
 // this allocates a new char[]
 for (char c : string.toCharArray()) {
 if (!Character.isDigit(c)) {
 return false;
 }
 }
 return true;
}
```
Note that many loops over characters can be expressed using streams with
`String#chars` or `String#codePoints`, for example:

```
boolean isDigits(String string) {
 string.codePoints().allMatch(Character::isDigit);
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("LoopOverCharArray")` to the enclosing element.

#
 LoopToTestParameter

 Migrate loops in tests to use github.com/google/TestParameterInjector. Test parameterization executes each input case in strict isolation, ensuring that a single failure doesn't halt the rest of your test case while providing clear, per-case reporting without the need for manual loops.

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
MalformedInlineTag
This Javadoc tag is malformed. The correct syntax is {@tag and not @{tag.
Severity
WARNING

*Alternate names: MathAbsoluteRandom*

`Math.abs`
returns a negative number when called with the largest negative number.

Example:

```
int veryNegative = Math.abs(Integer.MIN_VALUE);
long veryNegativeLong = Math.abs(Long.MIN_VALUE);
```
When trying to generate positive random numbers or fingerprints by using
`Math.abs` around a random positive-or-negative integer (or long), there will a
rare edge case where the returned value will be negative.

This is because there is no positive integer with the same magnitude as
`Integer.MIN_VALUE`, which is equal to `-Integer.MAX_VALUE - 1`. Floating point
numbers don’t suffer from this problem, as the sign is stored in a separate bit.

Instead, one should use random number generation functions that are guaranteed to generate positive numbers:

```
Random r = new Random();
int positiveNumber = r.nextInt(Integer.MAX_VALUE);
```
or map negative numbers onto the non-negative range:

```
long lng = r.nextLong();
long bestForHashCodes = lng & Long.MAX_VALUE;
long bestForMath = LongMath.saturatedAbs(lng);
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("MathAbsoluteNegative")` to the enclosing element.

#
 MemoizeConstantVisitorStateLookups

 Anytime you need to look up a constant value from VisitorState, improve performance by creating a cache for it with VisitorState.memoize

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MemoizeConstantVisitorStateLookups")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("MisformattedTestData") to the enclosing element.

@SuppressWarnings("MisformattedTestData")

Consider a switch statement that doesn’t handle all possible values and doesn’t have a default:

```
enum Colors { RED, GREEN, BLUE }
switch (color) {
 case RED:
 case GREEN:
 paint(color);
 break;
}
```
The author’s intent isn’t clear. There are three possibilities:

The default case is known to be impossible. This could be made clear by
adding:

`default: throw new AssertionError(color);`

The code intentionally ‘falls out’ of the switch on the default case, and
execution continues below. This could be made clear by adding:

`default: // fall out`

The code has a bug, and the missing cases should have been handled.

To avoid this ambiguity, the Google Java Style Guide requires each switch statement on an enum type to either handle all values of the enum, or have a default statement group.

If libraries are compiled against different versions of the same enum it’s possible for the switch statement to encounter an enum value despite it otherwise being thought to be exhaustive. If there is no default branch code execution will simply fall out of the switch statement.

Since developers may have assumed this to be impossible, it may be helpful to add a default branch when library skew is a concern, however, you may not want to give up checking to ensure that all cases are handled. Therefore if a default branch exists with a comment containing “skew”, the default will not be considered for exhaustiveness. For example:

```
enum TrafficLightColour { RED, GREEN, YELLOW }
void approachIntersection(TrafficLightColour state) {
 switch (state) {
 case GREEN:
 proceed();
 break;
 case YELLOW:
 case RED:
 stop();
 break;
 default: // In case of skew we may get an unknown value, always stop.
 stop();
 break;
 }
}
```
In this case the default branch is providing runtime safety for unknown enum values while also still enforcing that all known enum values are handled.

Note: The UnnecessaryDefaultInEnumSwitch check will not classify the default as unnecessary if it has the “skew” comment.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MissingCasesInEnumSwitch")` to the enclosing element.

*Alternate names: missing-fail*

When testing for exceptions in JUnit, it is easy to forget the call to `fail()`:

```
try {
 someOperationThatShouldThrow();
 // forget to call Assert.fail()
} catch (SomeException expected) {
 assertThat(expected).hasMessage("Operation failed");
}
```
This is better:

```
import static org.junit.Assert.fail;
try {
 someOperationThatShouldThrow();
 fail();
} catch (SomeException expected) {
 assertThat(expected).hasMessage("Operation failed");
}
```
But using `assertThrows` is preferable and the least error prone:

```
import static org.junit.Assert.assertThrows;
SomeException expected =
 assertThrows(SomeException.class, () -> someOperationThatShouldThrow());
assertThat(expected).hasMessage("Operation failed");
```
Without the call to `fail()`, the test is broken: it will pass if the exception
is never thrown *or* if the exception is thrown with the expected message.

If the try/catch block is defensive and the exception may not always be thrown, then the exception should be named ‘tolerated’.

This checker uses heuristics that identify as many occurrences of the problem as possible while keeping the false positive rate low (low single-digit percentages).

Five methods were explored to detect missing `fail()` calls, triggering if no
`fail()` is used in a `try/catch` statement within a JUnit test class:

`assert*()` method in the catch block.`catch` block.`catch` block is empty.`try` block contains only a single statement.Only the first three yield useful results and also required some more refinement
to reduce false positives. In addition, the checker does not trigger on comments
in the `catch` block due to implementation complexity.

To reduce false positives, no match is found if any of the following is true:

`fail` in its name is present in either catch or try block.`throw` statement or synonym (`assertTrue(false)`, etc.) is present in
either `catch` or `try` block.`setUp`, `tearDown`, `@Before`, `@After`,
`suite` or`main` method.`try` or `catch` block or immediately after.`InterruptedException`, `AssertionError`,
`junit.framework.AssertionFailedError` or `Throwable`.`while(true)` loop.`try` or `catch` block contains a `continue;` statement.`try/catch` statement also contains a `finally` statement.`catch` block.In addition, for occurrences which matched because they have a call to an
`assert*()` method in the catch block, no match is found if any of the following
characteristics are present:

`assertTrue/False(boolean variable or field)` in the catch block.`try` block is an `assert*()` (that is not a
noop): `assertFalse(false)`, `assertTrue(true))` or `Mockito.verify()` call.Suppress false positives by adding the suppression annotation `@SuppressWarnings("MissingFail")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("MissingImplementsComparable") to the enclosing element.

@SuppressWarnings("MissingImplementsComparable")

The Google Java Style Guide §6.1 requires that a method is marked with
the `@Override` annotation whenever it is legal. This includes a class method
overriding a superclass method, a class method implementing an interface method,
and an interface method respecifying a superinterface method.

Exception: `@Override` may be omitted when the parent method is `@Deprecated`.
If the flag `-XepOpt:MissingOverride:IgnoreInterfaceOverrides=true` is used,
`@Override` can be omitted for an interface method respecifying a superinterface
method.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MissingOverride")` to the enclosing element.

A Refaster template consists of multiple methods. Typically, each method in the class has an annotation. If a method has no annotation, this is likely an oversight.

```
static final class MethodLacksBeforeTemplateAnnotation {
 @BeforeTemplate
 boolean before1(String string) {
 return string.equals("");
 }
 // @BeforeTemplate is missing
 boolean before2(String string) {
 return string.length() == 0;
 }
 @AfterTemplate
 @AlsoNegation
 boolean after(String string) {
 return string.isEmpty();
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("MissingRefasterAnnotation")` to the enclosing element.

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
MissingSummary
A summary line is required on public/protected Javadocs.
Severity
WARNING

It is dangerous for a method to return a mutable instance in some circumstances, but immutable in others. Doing so may lead users of your API to make incorrect assumptions about the mutability of the return type. For example, consider this method:

```
List<Integer> primeFactors(int n) {
 if (isPrime(n)) {
 return Collections.singletonList(n);
 }
 List<Integer> factors = new ArrayList<>();
 for (...) {
 factors.add(i);
 }
 return factors;
}
```
If someone were to add another method to include the trivial factor `1`, a bug
will be introduced.

```
List<Integer> primeFactorsAndOne(int n) {
 List<Integer> primeFactors = primeFactors(n);
 primeFactors.add(1);
 return primeFactors;
}
```
`primeFactorsAndOne` will behave as intended for composite numbers, but throw an
exception for primes.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MixedMutabilityReturnType")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("MockIllegalThrows") to the enclosing element.

@SuppressWarnings("MockIllegalThrows")

#
 MockNotUsedInProduction

 This mock is instantiated and configured, but is never passed to production code. It should be either removed or used.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MockNotUsedInProduction")` to the enclosing element.

Collections and proto builders which are created and mutated but never used may be a sign of a bug, for example:

```
 MyProto.Builder builder = MyProto.newBuilder();
 if (field != null) {
 MyProto.NestedField.Builder nestedBuilder = MyProto.NestedField.newBuilder();
 nestedBuilder.setValue(field);
 // Oops--forgot to do anything with nestedBuilder.
 }
 return builder.build();
```
Likewise, converting a proto to a builder and modifying it is a no-op unless something is done with the return value:

```
 void setFoo(MyProto proto, String foo) {
 proto.toBuilder().setFoo(foo).build();
 }
```
As protos are immutable, either the return value must be used:

```
 @CheckReturnValue
 MyProto withFoo(MyProto proto, String foo) {
 return proto.toBuilder().setFoo(foo).build();
 }
```
or the Builder modified in place:

```
 void setFoo(MyProto.Builder protoBuilder, String foo) {
 protoBuilder.setFoo(foo);
 }
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ModifiedButNotUsed")` to the enclosing element.

From the javadoc for
`Iterator.remove`:

The behavior of an iterator is unspecified if the underlying collection is modified while the iteration is in progress in any way other than by calling this method, unless an overriding class has specified a concurrent modification policy.

That is, prefer this:

```
Iterator<String> it = ids.iterator();
while (it.hasNext()) {
 if (shouldRemove(it.next())) {
 it.remove();
 }
}
```
to this:

```
for (String id : ids) {
 if (shouldRemove(id)) {
 ids.remove(id); // will cause a ConcurrentModificationException!
 }
}
```
TIP: This pattern is simpler with Java 8’s
`Collection.removeIf`:

```
ids.removeIf(id -> shouldRemove(id));
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ModifyCollectionInEnhancedForLoop")` to the enclosing element.

From the javadoc for
`java.util.stream: Non-interference`:

Accordingly, behavioral parameters in stream pipelines whose source might not be concurrent should never modify the stream’s data source. A behavioral parameter is said to interfere with a non-concurrent data source if it modifies, or causes to be modified, the stream’s data source. The need for non-interference applies to all pipelines, not just parallel ones. Unless the stream source is concurrent, modifying a stream’s data source during execution hg of a stream pipeline can cause exceptions, incorrect answers, or nonconformant behavior.

That is, prefer this:

```
mutableValues.stream()
 .filter(x -> x < 5)
 .collect(Collectors.toList()) // Terminate stream before source modification.
 .forEach(mutableValues::remove);
```
to this:

```
mutableValues.stream()
 .filter(x -> x < 5)
 .forEach(mutableValues::remove);
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ModifySourceCollectionInStream")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("MultimapKeys") to the enclosing element.

@SuppressWarnings("MultimapKeys")

Suppress false positives by adding the suppression annotation @SuppressWarnings("MultipleNullnessAnnotations") to the enclosing element.

@SuppressWarnings("MultipleNullnessAnnotations")

Suppress false positives by adding the suppression annotation @SuppressWarnings("MultipleParallelOrSequentialCalls") to the enclosing element.

@SuppressWarnings("MultipleParallelOrSequentialCalls")

Increment operators in method calls are dubious and while argument lists are evaluated left-to-right, documentation suggests that code not rely on this specification. In addition, code is clearer when each expression contains at most one side effect.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MultipleUnaryOperatorsInMethodCall")` to the enclosing element.

Nonzero-length arrays are mutable. Declaring one `public static final` indicates
that the developer expects it to be a constant, which is not the case. Making it
`public` is especially dangerous since clients of this code can modify the
contents of the array.

There are two ways to fix this problem:

`ImmutableList`.`private` and add a `public` method that returns a copy of
the `private` array.See Effective Java 3rd Edition §15 for more details.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MutablePublicArray")` to the enclosing element.

The problem we’re trying to prevent is clashes between the names of classes/methods and contextual keywords. Clashes can occur when (1.) naming a class/method, or (2.) when invoking.

Change problematic names for classes and methods.

```
class Foo {
 ...
 // This can clash with the contextual keyword "yield"
 void yield() {
 ...
 }
}
```
Another example:

```
// This can clash with Java modules (JPMS)
static class module {
 ...
}
```
In recent versions of Java, `yield` is a restricted identifier:

```
class T {
 void yield() {}
 {
 yield();
 }
}
```
```
$ javac --release 20 T.java
T.java:3: error: invalid use of a restricted identifier 'yield'
 yield();
 ^
 (to invoke a method called yield, qualify the yield with a receiver or type name)
1 error
```
To invoke existing methods called `yield`, use qualified names:

```
class T {
 void yield() {}
 {
 this.yield();
 }
}
```
```
class T {
 static void yield() {}
 {
 T.yield();
 }
}
```
```
class T {
 void yield() {}
 class I {
 {
 T.this.yield();
 }
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("NamedLikeContextualKeyword")` to the enclosing element.

Integer division is suspicious when the target type of the expression is a float. For example:

```
private static final double ONE_HALF = 1 / 2; // Actually 0.
```
If you specifically want the integer result, consider pulling out variable to make it clear what’s happening. For example, instead of:

```
// Deduct 10% from the grade for every week the assignment is late:
float adjustedGrade = grade - (days / 7) * .1;
```
Prefer:

```
// Deduct 10% from the grade for every week the assignment is late:
int fullWeeks = days / 7;
float adjustedGrade = grade - fullWeeks * .1;
```
Similarly, multiplication of two `int` values which are then cast to a `long` is
problematic, as the `int` multiplication could overflow. It’s better to perform
the multiplication using `long` arithmetic.

```
long secondsToNanos(int seconds) {
 return seconds * 1_000_000_000; // Oops; starts overflowing around 2.15 seconds.
}
```
Instead, prefer:

```
long secondsToNanos(int seconds) {
 return seconds * 1_000_000_000L; // Or ((long) seconds) * 1_000_000_000.
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("NarrowCalculation")` to the enclosing element.

The compound assignment `E1 op= E2` could be mistaken for being equivalent to
`E1 = E1 op E2`. However, this is not the case: compound assignment operators
automatically cast the result of the computation to the type on the left hand
side. So `E1 op= E2` is actually equivalent to `E1 = (T) (E1 op E2)`, where `T`
is the type of `E1`.

If the type of the expression is wider than the type of the variable (i.e. the variable is a byte, char, short, or float), then the compound assignment will perform a narrowing primitive conversion. Attempting to perform the equivalent simple assignment would generate a compilation error.

For example, the following does not compile:

```
byte b = 0;
b = b << 1;
// ^
// error: incompatible types: possible lossy conversion from int to byte
```
However, the compound assignment form is allowed:

```
byte b = 0;
b <<= 1;
```
Similarly, if the expression is a floating point type (float or double), and the variable is an integral type (long, int, short, byte, or char), then an implicit conversion will be performed.

Example:

```
long l = 180;
l = l * 2.0f;
// ^
// error: incompatible types: possible lossy conversion from float to long
```
Again, the compound assignment form is permitted:

```
long l = 180;
l *= 2.0f;
```
See Puzzle #9 in ‘Java Puzzlers: Traps, Pitfalls, and Corner Cases’ for more information.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NarrowingCompoundAssignment")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("NegativeCharLiteral") to the enclosing element.

@SuppressWarnings("NegativeCharLiteral")

Suppress false positives by adding the suppression annotation @SuppressWarnings("NestedInstanceOfConditions") to the enclosing element.

@SuppressWarnings("NestedInstanceOfConditions")

Starting in JDK 13, calls to `FileSystem.newFileSystem(path, null)` are
ambiguous.

The calls match both:

To disambiguate, add a cast to the desired type, to preserve the pre-JDK 13 behaviour.

That is, prefer this:

```
FileSystem.newFileSystem(path, (ClassLoader) null);
```
Instead of this:

```
FileSystem.newFileSystem(path, null);
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("NewFileSystem")` to the enclosing element.

Flags instances of non-API types from being accepted or returned in public APIs.

`int[]`, `Integer[]`, `double[]`, `Double[]`,
`long[]`, `Long[]`). Prefer `ImmutableIntArray`, `ImmutableDoubleArray`, or
`ImmutableLongArray` instead.
 `int... rest`) are `ArrayList`, `LinkedList`, `HashSet`, `LinkedHashSet`, `TreeSet`,
`HashMap`, `LinkedHashMap`, `TreeMap`). Prefer interface types (`List`,
`Set`, `Map`).`ImmutableCollection`,
`ImmutableList`, `ImmutableSet`, or `ImmutableMap` as method parameters.
Prefer accepting `Collection`, `List`, `Set`, `Map`, or `Iterable` for
parameter generality.`java.util.Optional` or
`com.google.common.base.Optional` as method parameters. Prefer method
overloading: creating one signature with the parameter and one without (or
use `@Nullable` parameters).`com.google.common.base.Pair`:`Pair` across API boundaries.
Define a well-named class or record instead.`Iterator` (prefer `Stream` or collecting
to an `ImmutableList`/`ImmutableSet`) or accepting `Stream` as a parameter
(prefer `Iterable` or `Collection`).
 `Iterator`s limits caller options`com.google.protobuf.Duration`, `Timestamp`, or
`com.google.type.*` types across public APIs instead of standard
`java.time.*` types (`Duration`, `Instant`, `LocalDate`, etc.).`FluentLogger` or `GoogleLogger` instances
across method boundaries; this can break standard per-class logger
initialization patterns.`List` rather than `ArrayList`) to give callers
flexibility in implementation details.`ImmutableIntArray`,
`ImmutableDoubleArray`, and `ImmutableLongArray` provide immutable, safe,
and efficient alternatives.Suppress false positives by adding the suppression annotation `@SuppressWarnings("NonApiType")` to the enclosing element.

The volatile modifier ensures that updates to a variable are propagated predictably to other threads. A read of a volatile variable always returns the most recent write by any thread.

However, this does not mean that all updates to a volatile variable are atomic. For example, if you increment or decrement a volatile variable, you are actually doing (1) a read of the variable, (2) an increment or decrement of a local copy, and (3) a write back to the variable. Each step is atomic individually, but the whole sequence is not, and it will cause a race condition if two threads try to increment or decrement a volatile variable at the same time. The same is true for compound assignment, e.g. foo += bar.

If you intended for this update to be atomic, you should wrap all update operations on this variable in a synchronized block. If the variable is an integer, you could use an AtomicInteger instead of a volatile int.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NonAtomicVolatileUpdate")` to the enclosing element.

Types being referred to by non-canonical names can be confusing. For example,

```
public final class Entries {
 private final ImmutableList<ImmutableMap.Entry<String, Long>> entries;
 public Entries(Map<String, Long> map) {
 this.entries = ImmutableList.copyOf(map.entrySet());
 }
}
```
There is nothing special about `ImmutableMap.Entry`; it is precisely the same
type as `Map.Entry`. This example makes it look deceptively as though
`ImmutableList<ImmutableMap.Entry<?, ?>>` is an immutable type and therefore
safe to store indefinitely, when really it offers no more safety than
`ImmutableList<Map.Entry<?, ?>>`. You should use ```
ImmutableList<Map.Entry<?,
?>>
```
 instead, so it’s obvious what type you’re referring to.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NonCanonicalType")` to the enclosing element.

Defining a method that looks like `Object#equals` but doesn’t actually override
`Object#equals` is dangerous. The result of the comparison could differ
depending on the declared type of the argument passed into the `equals` call.

For example, consider this code:

```
public class Example {
 private int value;
 public Example(int value) {
 this.value = value;
 }
 public boolean equals(Example other) {
 return this.value == other.value;
 }
 public static void main(String[] args) {
 Example exampleA = new Example(1);
 Example exampleB = new Example(1);
 System.out.println(exampleA.equals(exampleB));
 }
}
```
This will print `true`. Suppose you refactor it so that the variable `exampleB`
is declared as an `Object` instead. Now this code will print `false`, because
Java’s overload resolution will choose the default `equals(Object)`
implementation instead of the `equals(Example)` method defined in this class.

If this equals method is intended to be a type-specific helper for an `equals`
method that *does* override `Object#equals`, either inline it into the
overriding `equals` method or rename it to something other than `equals` to
avoid ambiguity in overload resolution.

If you don’t want to write and maintain `equals` and `hashCode` methods by hand,
consider rewriting this class to use
AutoValue.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NonOverridingEquals")` to the enclosing element.

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
NotJavadoc
Avoid using `/**` for comments which aren't actually Javadoc.
Severity
WARNING

Passing a literal `null` to an `Optional` accepting parameter is likely a bug.
`Optional` is already designed to encode missing values through a non-`null`
instance.

```
Optional<Integer> double(Optional<Integer> i) {
 return i.map(i -> i * 2);
}
Optional<Integer> doubled = double(null);
```
```
Optional<Integer> doubled = double(Optional.empty());
```
This is a scenario that can easily happen when refactoring code from accepting
`@Nullable` parameters to accept `Optional`s. Note that the check will not match
if the parameter is explicitly annotated `@Nullable`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NullOptional")` to the enclosing element.

Constructors never return null.

Suppress false positives by adding the suppression annotation @SuppressWarnings("NullableConstructor") to the enclosing element.

@SuppressWarnings("NullableConstructor")

`Optional` is a container object which may or may not contain a value. The
presence or absence of the contained value should be demonstrated by the
`Optional` object itself.

Using an Optional variable which is expected to possibly be null is discouraged.
An nullable Optional which uses `null` to indicate the absence of the value will
lead to extra work for `null` checking when using the object and even cause
exceptions such as `NullPointerException`. It is best to indicate the absence of
the value by assigning it an empty optional.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NullableOptional")` to the enclosing element.

Primitives can never be null.

Suppress false positives by adding the suppression annotation @SuppressWarnings("NullablePrimitive") to the enclosing element.

@SuppressWarnings("NullablePrimitive")

For `@Nullable` type annotations (such as
`org.checkerframework.checker.nullness.qual.Nullable`), `@Nullable byte[]` means
a ‘non-null array of nullable bytes’, and `byte @Nullable []` means a ‘nullable
array of non-null bytes’. Since primitive types cannot be null, the former is
incorrect.

Some other nullness annotations (such as `javax.annotation.Nullable`) are
*declaration* annotations rather than *type* annotations. Their meaning is
different: For such annotations, `@Nullable byte[]` refers to ‘a nullable array
of non-null bytes,’ and `byte @Nullable []` is rejected by javac. Thus, this
check never reports errors for usages of declaration annotations.

See also: https://checkerframework.org/manual/#faq-array-syntax-meaning

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NullablePrimitiveArray")` to the enclosing element.

Nullness annotations directly on type parameters are interpreted differently by different tools. — unless you are using the Checker Framework.

`class Foo<@Nullable T>`To the Checker Framework, this means that a given type argument *must* be
nullable. For example, in the Checker Framework’s default JDK stubs,
the type argument to `ThreadLocal` must be nullable.

To Kotlin, this means that the type argument *can* be nullable but need not be
so.

If you want the Checker Framework interpretation, then keep your code as it is,
and suppress this warning. If you want the “*can* be nullable” interpretation,
change to `class Foo<T extends @Nullable Object>`.

The effects of that change would be:

It is a behavior change for the Checker Framework, one that could even
produce local or non-local the Checker Framework failures. As discussed
above, it may be a desirable change, and it’s likely to be safe unless you
are using it with `ThreadLocal`.

This is probably not a behavior change for Kotlin.

`class Foo<@NonNull T>`To the Checker Framework, this means to allow *any* nullness for the type
argument(!). (It
sets the *lower* bound
to `@NonNull`, which is already the default there.)

To Kotlin, this means that the type argument must be non-nullable.

Users probably want `class Foo<T extends @NonNull Object>`. Or, if they’re
within the scope of `@NullMarked` and they are using tools that recognize it
(such as Kotlin), they may prefer `class Foo<T>`, which is equivalent but
shorter.

The effects of that change would be:

The JSpecify spec says that usages of their annotations on type parameters are unrecognized (Javadoc, spec). This specification choice is motivated by the disagreement in tool behavior discussed above.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NullableTypeParameter")` to the enclosing element.

void-returning methods cannot return null.

Suppress false positives by adding the suppression annotation @SuppressWarnings("NullableVoid") to the enclosing element.

@SuppressWarnings("NullableVoid")

Nullness annotations directly on wildcard types are interpreted differently by different tools. — unless you are using the Checker Framework.

`Foo<@Nullable ?>`To the Checker Framework, this means that the type argument *must* be nullable.
They use this
in `ExecutorService`.

To Kotlin, this has no effect. That means that the type argument
*can* be nullable but need not be so.

While Checker Framework users do sometimes want `Foo<@Nullable ?>`, we commonly
see them use it in places where `Foo<?>` would also be correct and would be more
flexible. To fully preserve Kotlin behavior, Kotlin users may wish to write
“`Foo<? extends @Nullable Object>`” (unless they are within the scope of
`@NullMarked`, in which case `Foo<?>` is equivalent).

The effects of that change would be:

`ExecutorService` and the `Future`
objects that it produces.`Foo<?>` is probably not a behavior change for Kotlin.`Foo<@NonNull ?>`To the Checker Framework, this means that the type argument must be non-nullable. (See the Checker Framework docs.)

To Kotlin, this has no effect. That means that the type argument can still be nullable.

We recommend a change to `Foo<? extends @NonNull Object>` (or, within the scope
of `@NullMarked`, `Foo<? extends Object>`).

The effects of that change would be:

The JSpecify spec says that usages of their annotations on wildcard types are unrecognized (Javadoc, spec). This specification choice is motivated by the disagreement in tool behavior discussed above.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NullableWildcard")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("ObjectEqualsForPrimitives") to the enclosing element.

@SuppressWarnings("ObjectEqualsForPrimitives")

Calling `toString` on objects that don’t override `toString()` doesn’t provide
useful information (just the class name and the `hashCode()`).

Consider overriding toString() function to return a meaningful String describing the object.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ObjectToString")` to the enclosing element.

`Objects.hashCode` takes an `Object` parameter, and will either return `0` when
the parameter is `null`, or call the underlying `hashCode` function of the
Object reference.

Passing a primitive value to `Objects.hashCode` function results in boxing the
primitive, then calling the boxed object’s `hashCode`. You can get the same
result by using, e.g.: `Long.hashCode(long)` to get the effective hash code of a
primitive `long`. If you’re calling this method outside of your own `hashCode()`
implementation, prefer to use the `BoxedClass.hashCode(primitive)` functions to
avoid unwanted boxed.

If you’re implementing a `hashCode` function for your **own** class that
consists of a single primitive value, you may want to consider some of these
alternatives:

```
@Override
public int hashCode() {
 // This function will box intValue into an Integer, and wrap *that* in an
 // array, but will generate a hashCode which is likely to be different than
 // the hashCode of the boxed version of the intValue. This makes it easier
 // to add more fields to the class and hashCode method (just by adding more
 // fields to the hash call), but comes at a potential performance penalty.
 return Objects.hash(intValue);
}
```
```
@Override
public int hashCode() {
 // This function will avoid boxing the int to an Integer, and is an explicit
 // acknowledgement that the hashCode() of *this* class is the same as the
 // hash code of the underlying intValue.
 return Integer.hashCode(intValue);
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ObjectsHashCodePrimitive")` to the enclosing element.

The Google Java Style Guide §4.7 states:

Optional grouping parentheses are omitted only when author and reviewer agree that there is no reasonable chance the code will be misinterpreted without them, nor would they have made the code easier to read. It is not reasonable to assume that every reader has the entire Java operator precedence table memorized.

Use grouping parentheses to disambiguate expressions that could be misinterpreted.

For example, consider this:

```
boolean d = (a && b) || c;
boolean e = (a || b) ? c : d;
int z = (x + y) << 2;
```
Instead of this:

```
boolean r = a && b || c;
boolean e = a || b ? c : d;
int z = x + y << 2;
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("OperatorPrecedence")` to the enclosing element.

Using `Optional#map` (or `Optional#transform`) to map to another `Optional`
might indicate an error, or at least quite hard-to-reason-about code that may
turn into an error.

For example,

```
class AccountManager {
 /** Retrieves the current user, or absent if not logged in. */
 Optional<Account> getUser() {
 ...
 }
 /**
 * Returns an administrative token for the user, or absent if they do not have
 * such privileges.
 */
 Optional<Token> getAdminToken() {
 }
 void doScaryAdminThings() {
 if (getUser().map(AccountManager::getAdminToken).isPresent()) {
 // do privileged things.
 }
 }
}
```
In this case, assuming `getAdminToken` does not throw, the conditional is
equivalent to `getUser().isPresent()`, which is not what was intended.
`getUser().flatMap(AccountManager::getAdminToken).isPresent()` would be correct.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("OptionalMapToOptional")` to the enclosing element.

Calling `get()` on an `Optional` that is not present will result in a
`NoSuchElementException`.

This check detects cases where `get()` is called when the optional is definitely
not present, e.g.:

```
if (!o.isPresent()) {
 return o.get(); // this will throw a NoSuchElementException
}
```
```
if (o.isEmpty()) {
 return o.get(); // this will throw a NoSuchElementException
}
```
These cases are almost definitely bugs; the intent may have been to invert the test:

```
if (o.isPresent()) {
 return o.get();
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("OptionalNotPresent")` to the enclosing element.

Passing a string that contains format specifiers to a method that does not perform string formatting is usually a mistake.

Do this:

```
if (!isValid(arg)) {
 throw new IllegalArgumentException(String.format("invalid arg: %s", arg));
}
```
or this:

```
logger.atWarning().log("invalid arg: %s", arg);
```
Not this:

```
if (!isValid(arg)) {
 throw new IllegalArgumentException("invalid arg: %s");
}
```
or this:

```
logger.atWarning().log("invalid arg: %s");
```
If the method you’re calling actually accepts a format string, you can annotate
that method with `@FormatMethod` to ensure that callers correctly pass
format strings (and to inform Error Prone that the method call you’re making
doesn’t orphan a format string).

Suppress false positives by adding the suppression annotation `@SuppressWarnings("OrphanedFormatString")` to the enclosing element.

The `outline` CSS property provides visual indicators as to which element is
currently selected within a web page. These are the dotted lines you see
surrounding links, etc. when you tab to them using the keyboard.

These indicators are important for users navigating without a mouse (such as
those with visual or mobility impairments). Setting `outline` style to `"none"`
or `0` removes these indicators, leaving these users without any way to tell
where they are within the page, therefore making the page inaccessible.

Caveat: `outline` is not the *only* way to emphasize selected elements. You may
instead choose to change the background color, add an underline, or otherwise
make them visually distinct. Learn more & get alternative suggestions at
OutlineNone.com.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("OutlineNone")` to the enclosing element.

Many logging tools build a string representation out of getMessage() and ignores toString() completely.

Suppress false positives by adding the suppression annotation @SuppressWarnings("OverrideThrowableToString") to the enclosing element.

@SuppressWarnings("OverrideThrowableToString")

*Alternate names: overrides*

Even though varargs methods are different than methods with an array parameter
as the last parameter, varargs methods are compiled into bytecode as methods
with an array as the last parameter. When a varargs method is *called*, the Java
compiler will insert instructions to automatically box the varargs arguments
into an array.

This detail means that, for example, you can’t declare two methods in the same class where the final parameter is an array in one method, and a varargs of the same type in the other:

```
class Foo {
 void bah(double a, double... others) {}
 void bah(double baz, double[] myArray) {} // ERROR: bah(double, double[]) already defined
}
```
This also means that one method with varargs can override another method with an array as the final parameter:

```
class A {
 void something(int... ints) {}
}
class B extends A {
 @Override
 void something(int[] ints) {}
}
```
This overriding may be unintentional (since the signatures ‘look’ different, the programmer may be unaware that an overriding has occurred).

Even if this overriding is intentional, it causes inconsistencies at call-sites, as the code required to invoke the overridden method depends on the static type of the variable being operated on.

Given the example classes above, observe the result on the client side:

```
class Client {
 public static void main(String[] args) {
 B b = new B();
 A a = b;
 a.something(new int[]{1}); // OK, array invocation of varargs method
 b.something(new int[]{1}); // OK, direct array invocation
 a.something(2); // OK, varargs invocation with 1 element
 // Very strange compile-time error:
 // error: A.something(int...) is defined in an inaccessible class or interface
 b.something(1);
 }
}
```
To avoid these ambiguities, use the same parameter style (varargs or explicit arrays) when overriding methods.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("Overrides")` to the enclosing element.

Unlike with `@javax.inject.Inject`, if a method overrides a method annotated
with `@com.google.inject.Inject`, Guice will inject it even if it itself is not
annotated. This differs from the behavior of methods that override
`@javax.inject.Inject` methods since according to the JSR-330 spec, a method
that overrides a method annotated with `@javax.inject.Inject` will not be
injected unless it itself is annotated with `@Inject`. Because of this
difference, it is recommended that you annotate this method explicitly.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("OverridesGuiceInjectableMethod")` to the enclosing element.

Inconsistently ordered parameters in method overrides mostly indicate an accidental bug in the overriding method. An example for an overriding method with inconsistent parameter names:

```
class A {
 public void foo(int foo, int baz) { ... }
}
class B extends A {
 @Override
 public void foo(int baz, int foo) { ... }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("OverridingMethodInconsistentArgumentNamesChecker")` to the enclosing element.

In certain contexts literal arguments - such as `0`, `""`, `true` and `false`,
or `null` - can make it difficult for readers to know what a method will do.
Defining methods that take boolean parameters or otherwise expect users to pass
in ambiguous literals is generally discouraged. However, when you must call such
a method you’re encouraged to use the parameter name as an inline comment at the
call site, so that readers don’t need to look at the method declaration to
understand the parameter’s purpose.

Error Prone recognizes such comments that use the following formatting, and emits an error if the comment doesn’t match the name of the corresponding formal parameter:

```
booleanMethod(/* enableFoo= */ true);
```
Varargs methods are also supported using `...` syntax: ```
void varargsMethod(/*
states...= */ true, true, false);
```

The check also recognizes comments with whitespace variations, e.g. `/*foo =*/`,
but the form `/* foo= */` is preferred.

If the comment deliberately does not match the formal parameter name, using a
regular block comment without the `=` is recommended: `/* enableFoo */`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ParameterName")` to the enclosing element.

Pattern matching with `instanceof` allows writing this:

```
void handle(Object o) {
 if (o instanceof Point(int x, int y)) {
 handlePoint(x, y);
 } else if (o instanceof String s) {
 handleString(s);
 }
}
```
which is more concise than an instanceof and a separate cast:

```
void handle(Object o) {
 if (o instanceof Point) {
 Point point = (Point) o;
 handlePoint(point.x(), point.y());
 } else if (o instanceof String) {
 String s = (String) o;
 handleString(s);
 }
}
```
For more information on pattern matching and `instanceof`, see
Pattern Matching for the instanceof Operator

Suppress false positives by adding the suppression annotation `@SuppressWarnings("PatternMatchingInstanceof")` to the enclosing element.

#
 PreconditionsCheckNotNullRepeated

 Including the first argument of checkNotNull in the failure message is not useful, as it will always be `null`.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("PreconditionsCheckNotNullRepeated")` to the enclosing element.

APIs that accept a `String` charset name often have an overload that accepts a
`java.nio.charset.Charset`. Prefer using the `Charset` overload, as it provides
stronger typing.

If a `Charset` instance is being converted to a `String` (e.g. via `.name()`,
`.displayName()`, or `.toString()`) just to call the `String` overload, the
conversion can simply be removed.

**Note:** This check relies on method parameter name information being preserved
in class files at runtime. Therefore, compiling the targeted libraries or JDK
stubs with the `-parameters` flag (introduced in
JEP 118) is required for the check to discover
candidate parameters reliably.

```
OutputStreamWriter writer = new OutputStreamWriter(out, "UTF-8");
String s = new String(bytes, "ISO-8859-1");
log("hello", charset.name());
```
```
OutputStreamWriter writer = new OutputStreamWriter(out, UTF_8);
String s = new String(bytes, ISO_8859_1);
log("hello", charset);
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("PreferCharsetOverload")` to the enclosing element.

#
 PreferInstanceofOverGetKind

 Prefer instanceof over getKind() checks where possible, as these work well with pattern matching instanceofs

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("PreferInstanceofOverGetKind")` to the enclosing element.

#
 PreferTestParameter

 When exhaustively testing all values of a single enum or boolean parameter, prefer @TestParameter over @TestParameters.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("PreferTestParameter")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("PreferThrowsTag") to the enclosing element.

@SuppressWarnings("PreferThrowsTag")

PrimitiveAtomicReference
Using compareAndSet with boxed primitives is dangerous, as reference rather than value equality is used. Consider using AtomicInteger, AtomicLong, AtomicBoolean from JDK or AtomicDouble from Guava instead.

Severity

WARNING

Suppression

Suppress false positives by adding the suppression annotation @SuppressWarnings("PrimitiveAtomicReference") to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("ProtectedMembersInFinalClass") to the enclosing element.

@SuppressWarnings("ProtectedMembersInFinalClass")

If you call duration.getNanos(), you must also call duration.getSeconds() in ‘nearby’ code. If you are trying to convert this duration to nanoseconds, you probably meant to use Durations.toNanos(duration) instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ProtoDurationGetSecondsGetNano")` to the enclosing element.

If you call timestamp.getNanos(), you must also call timestamp.getSeconds() in ‘nearby’ code. If you are trying to convert this timestamp to nanoseconds, you probably meant to use Timestamps.toNanos(timestamp) instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ProtoTimestampGetSecondsGetNano")` to the enclosing element.

The `@Inject` annotation is applied to **injection points** - methods,
constructors, or fields that a dependency injection tool like Guice or Dagger
will invoke or populate to construct an object.

Scope annotations (annotations that are themselves annotated with
`@ScopeAnnotation`) are used to denote the ‘effective scope’ of an object.
Examples of scope annotations include `@Singleton`, `@RequestScope`,
`@SessionScope`. These annotations can be used either on the declaration of a
type (so that every instantiation of that object will be handled with the
appropriate scope), or on a **provider method** (to declare that the particular
binding declared by that provider method is subject to that annotation’s scoping
rules):

```
@Singleton
class ExpensiveGlobalObject {
 @Inject ExpensiveGlobalObject(SomeDependencies stuff) {...}
}
class SomeModule extends AbstractModule {
 ...
 @Provides
 @Singleton
 OtherExpensiveObject provideOtherExpensiveObject() { return new Gold(); }
}
```
Qualifier annotations (annotations that are themselves annotated with
`@Qualifier` or `@BindingAnnotation`) are used to distinguish different
instances of the same type of object (the `@Red Robot` instead of the ```
@Blue
Robot
```
). These annotations can be used *inside* injection points, or on those
provider methods:

```
class RobotArena {
 @Inject RobotArena(@Red Robot redBot, @Blue Robot blueBot) {...}
}
class SomeModule extends AbstractModule {
 ...
 @Provides
 @Red
 Robot provideRedRobot() { return new Robot("red"); }
}
```
Notably: there is no method declaration where `@Inject` and either a scope or a
qualifier annotation can reasonably coexist. Either the method is an **injection
point**, and the dependency injection system will ignore the scope or qualifier
annotation, or the method is a **provider** method, and the `@Inject` method
won’t have any effect, as modules are constructed without a dependency injection
container.

```
class Example {
 // Here, @Singleton is ignored. Perhaps the scope should go onto Example, to make dependency
 // injection systems treat the Example type as a singleton.
 @Inject @Singleton
 Example(Robot something) {}
}
```
```
class MyModule extends AbstractModule {
 ...
 // The `@Inject` doesn't do anything: Guice will use this method to define a Singleton binding
 // for Database. When a Database is constructed by Guice, a DatabaseCredentials will be
 // constructed and this method will be invoked, but it isn't invoked until then.
 @Provides @Singleton @Inject
 Database providesSingletonDatabase(DatabaseCredentials db) { ... }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("QualifierOrScopeOnInjectMethod")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("ReachabilityFenceUsage") to the enclosing element.

@SuppressWarnings("ReachabilityFenceUsage")

Because `@Override` is a `@Target(METHOD)` annotation, it’s automatically
permitted on record component declarations. It would normally cause the
annotation to be copied to the generated accessor method, but in this case it’s
a SOURCE-retention annotation so there’s nothing to do, and the annotation is
ignored.

Note that the annotation *does not* mean that the generated accessor method for
the record is overriding something.

Also note that a hand-written accessor method in a record class can also always
use `@Override` regardless of supertype methods (see JLS §9.6.4.4):

If a method declaration in class or interface Q is annotated with

`@Override`, then one of the following three conditions must be true, or a compile-time error occurs:…

Q is a record class (§8.10), and the method is an accessor method for a record component of Q (§8.10.3)

For additional discussion, see this compile-dev@ thread and Error Prone issue #5174.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("RecordComponentOverride")` to the enclosing element.

Unnecessary control flow statements may be misleading in some contexts.

For instance, this method to find groups of prime order has a bug:

```
static ImmutableList<Group> filterGroupsOfPrimeOrder(Iterable<Group> groups) {
 ImmutableList.Builder<Group> filtered = ImmutableList.builder();
 for (Group group : groups) {
 for (int i = 2; i < group.order(); i++) {
 if (group.order() % i == 0) {
 continue;
 }
 }
 filtered.add(group);
 }
 return filtered.build();
}
```
The `continue` statement only breaks out of the innermost loop, so the input is
returned unchanged.

The most readable alternative is probably to avoid a nested loop entirely:

```
static ImmutableList<Group> filterGroupsOfPrimeOrder(Iterable<Group> groups) {
 ImmutableList.Builder<Group> filtered = ImmutableList.builder();
 for (Group group : groups) {
 if (!isPrime(group.order())) {
 continue;
 }
 filtered.add(group);
 }
 return filtered.build();
}
```
A labelled break statement would also be correct, but is quite uncommon:

```
static ImmutableList<Group> filterGroupsOfPrimeOrder(Iterable<Group> groups) {
 ImmutableList.Builder<Group> filtered = ImmutableList.builder();
 outer:
 for (Group group : groups) {
 for (int i = 2; i < group.order(); i++) {
 if (group.order() % i == 0) {
 continue outer;
 }
 }
 filtered.add(group);
 }
 return filtered.build();
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("RedundantControlFlow")` to the enclosing element.

This checker aims to improve the readability of new-style arrow `switch`es by
simplifying and refactoring them.

`case` consists only of a `yield`ing some value, it can be re-written
with just the value. For example, `case FOO -> { yield "bar"; }` can be
shortened to `case FOO -> "bar";````
case FOO -> { System.out.println("bar");
}
```
 can shortened to `case FOO -> System.out.println("bar");``return switch`:`enum` is covered by a `case` which `return`s, the
`return` can be factored out```
enum SideOfCoin {OBVERSE, REVERSE};
private String renderName(SideOfCoin sideOfCoin) {
 switch(sideOfCoin) {
 case OBVERSE -> {
 return "Heads";
 }
 case REVERSE -> {
 return "Tails";
 }
 }
 // This should never happen, but removing this will cause a compile-time error
 throw new RuntimeException("Unknown side of coin");
}
```
The transformed code is simpler and elides the “should never happen” handler.

```
enum SideOfCoin {OBVERSE, REVERSE};
private String renderName(SideOfCoin sideOfCoin) {
 return switch(sideOfCoin) {
 case OBVERSE -> "Heads";
 case REVERSE -> "Tails";
 };
}
```
`enum`, if the `case`s are
exhaustive, then a similar refactoring can be performed.`switch`:`case` just assigns a value to the same variable, the assignment
can potentially be factored out`switch`, the definition and assignment can potentially be combined```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private void updateScore(Suit suit) {
 int score = 0;
 switch(suit) {
 case HEARTS, DIAMONDS -> {
 score = -1;
 }
 case SPADES -> {
 score = 2;
 }
 case CLUBS -> {
 score = 3;
 }
 }
}
```
Which can be consolidated:

```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private void updateScore(Suit suit) {
 int score = switch(suit) {
 case HEARTS, DIAMONDS -> -1;
 case SPADES -> 2;
 case CLUBS -> 3;
 };
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("RefactorSwitch")` to the enclosing element.

Reference types should normally be compared for value equality with `equals()`,
not for object identity with `==` or `!=`.

`equals()` method uses object identity anyway?This check allows you to use `==` and `!=` in cases in which it can reliably
determine that they would behave identically to `equals()`, such as for enums or
for final classes that inherit the default `equals()` implementation from
`Object`. You might prefer `==` and `!=` in such cases, since it is less
verbose, especially compared to using `Objects.equals(a, b)` in cases in which
you need to tolerate null values.

It’s dangerous to rely on your instances being interned. We have no tooling to check or enforce that, and it’s easy to get wrong.

`Boolean` values? We `TRUE` and `FALSE` (and `null`). Surely Well, no, because some tricky client can always generate a new instance with
`new Boolean(true)`. Comparing with `equals` always works; comparing with `==`
doesn’t.

The check allows implementations of `Object#equals()` to perform reference
equality tests on the type equality is being implemented for. For example:

```
abstract class Foo {
 abstract String bar();
 @Override
 public boolean equals(Object other) {
 if (this == other) {
 return true; // fast path, reference equality is allowed here
 }
 if (!(other instanceof Foo)) {
 return false;
 }
 Foo that = (Foo) other;
 // value equality should still be used for types other than `Foo`
 return bar().equals(that.bar());
 }
}
```
In other cases, calling `Type#equals()` should be just as fast, because that
method will likely be inlined, and the first thing it will likely do is that
same instance comparison.

Alternatively, if you’re okay with accepting `null`, you could call
`java.util.Objects.equals()`, which first does a reference equality comparison
and then falls back to content equality for non-null arguments.

Assertion libraries provide clearer ways to assert this, with the bonus of providing better failure messages:

Truth:

```
assertThat(a).isSameInstanceAs(b);
assertThat(a).isNotSameInstanceAs(b);
```
AssertJ:

```
assertThat(a).isSameAs(b);
assertThat(a).isNotSameAs(b);
```
Classes override `equals` to express when two instances should be treated as
interchangeable with each other. Predominant Java libraries and practices are
built on that assumption. Defining a “magic instance” for such a type goes
against this whole practice, leaving you vulnerable to unexpected bugs.

Consider choosing a sentinel value within the domain of the type (the moral
equivalent of `-1` for indexOf function calls) that you could compare against
using the normal `equals` method.

Use `Optional<V>` as the value type of your map instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ReferenceEquality")` to the enclosing element.

Consider using `LinkageError` instead of `AssertionError` when rethrowing
reflective exceptions as unchecked exceptions, since it conveys more information
when reflection fails due to an incompatible change in the classpath.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("RethrowReflectiveOperationExceptionAsLinkageError")` to the enclosing element.

Detects no-op `return` statements in `void` functions when they occur at the end
of the method.

Instead of:

```
public void stuff() {
 int x = 5;
 return;
}
```
do:

```
public void stuff() {
 int x = 5;
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ReturnAtTheEndOfVoidFunction")` to the enclosing element.

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
ReturnFromVoid
Void methods should not have a @return tag.
Severity
WARNING
Tags
Style

Suppress false positives by adding the suppression annotation @SuppressWarnings("RobolectricShadowDirectlyOn") to the enclosing element.

@SuppressWarnings("RobolectricShadowDirectlyOn")

Suppress false positives by adding the suppression annotation @SuppressWarnings("RuleNotRun") to the enclosing element.

@SuppressWarnings("RuleNotRun")

Methods that return an ignored [Observable | Single | Flowable | Maybe ] generally indicate errors.

If you don’t check the return value of these methods, the observables may never execute. It also means the error case is not being handled

Suppress false positives by adding the suppression annotation `@SuppressWarnings("RxReturnValueIgnored")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("SameNameButDifferent") to the enclosing element.

@SuppressWarnings("SameNameButDifferent")

`Scanner.useDelimiter("\\A")` is not an efficient way to read an entire
`InputStream`.

```
Scanner scanner = new Scanner(inputStream, UTF_8).useDelimiter("\\A");
String s = scanner.hasNext() ? scanner.next() : "";
```
`Scanner` separates its input into “tokens” based on a delimiter that is a
regular expression. The regular expression `\A` matches the beginning of the
input, only, so there is no later delimiter and the single token consists of
every character read from the `InputStream`.

This works, but has multiple drawbacks:

`InputStream`. In that case there is no
token after `\A`. That’s why the extract above checks `hasNext()`. If you
forget to do that, you get `NoSuchElementException` in the empty case.It swallows `IOException`. Quoting the `Scanner` specification:

A scanner can read text from any object which implements the

`Readable`interface. If an invocation of the underlying readable’s`read()`method throws an`IOException`then the scanner assumes that the end of the input has been reached. The most recent`IOException`thrown by the underlying readable can be retrieved via the`ioException()`method.

`inputStream.readAllBytes()`.Instead, prefer one of the following alternatives:

Since Java 9, it has been possible to write this:

```
String s = new String(inputStream.readAllBytes(), UTF_8);
```
On Android, that does require API level 33, though. Guava’s
`ByteStreams.toByteArray(inputStream)`
is equivalent to `inputStream.readAllBytes()`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ScannerUseDelimiter")` to the enclosing element.

A common pattern for abstract `Builders` is to declare an instance method named
`self()`, which subtypes override and implement as `return this` (see Effective
Java 3rd Edition, Item 2).

Returning anything other than `this` from an instance method named `self()` with
a return type that matches the enclosing class will be confusing for readers and
callers.

If an unchecked cast is required, prefer a single-statement cast, with the suppression on the method (rather than the statement). For example:

```
 @SuppressWarnings("unchecked")
 default U self() {
 return (U) this;
 }
```
Instead of:

```
 default U self() {
 @SuppressWarnings("unchecked")
 U self = (U) this;
 return self;
 }
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("SelfAlwaysReturnsThis")` to the enclosing element.

A setter invoked with a value from the corresponding getter is often a mistake, for example:

```
if (from.hasFrobnicator()) {
 to.setFrobnicator(to.getFrobnicator());
}
```
This is easy to accidentally write, but is clearly meant to be,

```
if (from.hasFrobnicator()) {
 to.setFrobnicator(from.getFrobnicator());
}
```
The Java proto API is tolerant enough that the former code will compile and
execute fine, but it will set `frobnicator` to the default value for that field.

This pattern is occasionally used to ensure that a field is always present, even if it takes the default value, for example,

```
// ensure "always_present" is present
builder.setAlwaysPresent(builder.getAlwaysPresent());
```
This is not a no-op, but we’d encourage being more explicit about the condition,

```
if (!builder.hasAlwaysPresent()) {
 builder.setAlwaysPresent(false);
}
```
Or if `builder` is otherwise untouched, `builder.setAlwaysPresent(false)`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("SelfSet")` to the enclosing element.

The boolean operators `&&` and `||` should almost always be used instead of `&`
and `|`.

If the right hand side is an expression that has side effects or is expensive to
compute, `&&` and `||` will short-circuit but `&` and `|` will not, which may be
surprising or cause slowness.

If evaluating both operands is necessary for side effects, consider refactoring to make that explicit. For example, prefer this:

```
boolean rhs = hasSideEffects();
if (lhs && rhs) {
 // ...
}
```
to this:

```
if (lhs & hasSideEffects()) {
 // ...
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ShortCircuitBoolean")` to the enclosing element.

We’re trying to make `switch`es simpler to understand at a glance.
Misunderstanding the control flow of a `switch` is a common source of bugs.

As part of this simplification, new-style arrow (`->`) switches are encouraged
instead of old-style colon (`:`) switches. And where possible, neighboring cases
are grouped together (e.g. `case A, B, C`).

`:`) `switch`es:`case` and the `case`’s code. For example, ```
case
HEARTS:
```
`switch` block is large, just
skimming each `case` can be toilsome. Fall-through can also be conditional
(see example 5. below). In this scenario, one would potentially need to
reason about all possible flows for each `case`. When conditionally
falling-through multiple `case`s, the number of potential control flows can
grow rapidly`case`s are propagated down to later
`case`s, however the `->`) `switch`es:`case` and the `case`’s code. For example, ```
case
HEARTS ->
```
`case`s fall through; no control flow analysis needed`case`s (within a `switch`)`case`s; if you define a local
variable within a `case`, it can only be used within that specific `case`.```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private void foo(Suit suit) {
 switch(suit) {
 case HEARTS:
 System.out.println("Red hearts");
 break;
 case DIAMONDS:
 System.out.println("Red diamonds");
 break;
 case SPADES:
 // Fall through
 case CLUBS:
 bar();
 System.out.println("Black suit");
 }
}
```
Which can be simplified by grouping and using a new-style switch:

```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private void foo(Suit suit) {
 switch(suit) {
 case HEARTS -> System.out.println("Red hearts");
 case DIAMONDS -> System.out.println("Red diamonds");
 case SPADES, CLUBS -> {
 bar();
 System.out.println("Black suit");
 }
 }
}
```
`return switch ...`Sometimes `switch` is used with a `return` for each `case`, like this:

```
enum SideOfCoin {OBVERSE, REVERSE};
private String renderName(SideOfCoin sideOfCoin) {
 switch(sideOfCoin) {
 case OBVERSE:
 return "Heads";
 case REVERSE:
 return "Tails";
 }
 // This should never happen, but removing this will cause a compile-time error
 throw new RuntimeException("Unknown side of coin");
}
```
Note that even though a `case` is present for each possible value of the `enum`,
a boilerplate “should never happen” clause is still needed. The transformed code
is simpler and doesn’t need a “should never happen” clause.

```
enum SideOfCoin {OBVERSE, REVERSE};
private String renderName(SideOfCoin sideOfCoin) {
 return switch(sideOfCoin) {
 case OBVERSE -> "Heads";
 case REVERSE -> "Tails";
 };
}
```
If you nevertheless wish to define an explicit “should never happen” clause,
this can be accomplished by placing the logic inside a `default` case. For
example:

```
enum SideOfCoin {OBVERSE, REVERSE};
private String foo(SideOfCoin sideOfCoin) {
 return switch(sideOfCoin) {
 case OBVERSE -> "Heads";
 case REVERSE -> "Tails";
 default -> throw new RuntimeException("Unknown side of coin"); // should never happen
 };
}
```
When the checker detects an existing `default` that appears to be redundant, it
may suggest a secondary auto-fix which removes the redundant `default` and its
code (if any).

`switch`If every branch of a `switch` is making an assignment to the same variable, the
code can be simplified into a combined assignment and `switch`:

```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
int score = 0;
private void updateScore(Suit suit) {
 switch(suit) {
 case HEARTS:
 // Fall thru
 case DIAMONDS:
 score += -1;
 break;
 case SPADES:
 score += 2;
 break;
 case CLUBS:
 score += 3;
 }
}
```
This can be simplified as follows:

```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
int score = 0;
private void updateScore(Suit suit) {
 score += switch(suit) {
 case HEARTS, DIAMONDS -> -1;
 case SPADES -> 2;
 case CLUBS -> 3;
 };
}
```
Taking this one step further: if a local variable is defined, and then
immediately followed by a `switch` in which every `case` assigns to that same
variable, then all three (the `switch`, the variable declaration, and the
assignment) can be merged:

```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private void updateStatus(Suit suit) {
 int score;
 switch(suit) {
 case HEARTS:
 // Fall thru
 case DIAMONDS:
 score = 1;
 break;
 case SPADES:
 score = 2;
 break;
 case CLUBS:
 score = 3;
 }
 ...
}
```
Becomes:

```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private void updateStatus(Suit suit) {
 int score = switch(suit) {
 case HEARTS, DIAMONDS -> 1;
 case SPADES -> 2;
 case CLUBS -> 3;
 };
 ...
}
```
`switch`Even when the simplifications discussed above are not applicable, conversion to
new arrow `switch`es can be automated by this checker:

```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private void processEvent(Suit suit) {
 switch (suit) {
 case CLUBS:
 String message = "hello";
 var anotherMessage = "salut";
 processMessages(message, anotherMessage);
 break;
 case DIAMONDS:
 anotherMessage = "bonjour";
 processMessage(anotherMessage);
 }
}
```
Note that the local variables referenced in multiple cases are hoisted up out of
the `switch` statement, and `var` declarations are converted to explicit types,
resulting in:

```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private void processEvent(Suit suit) {
 String anotherMessage;
 switch (suit) {
 case CLUBS -> {
 String message = "hello";
 anotherMessage = "salut";
 processMessages(message, anotherMessage);
 }
 case DIAMONDS -> {
 anotherMessage = "bonjour";
 processMessage(anotherMessage);
 }
 }
}
```
Here’s an example of a complex statement `switch` with conditional fall-through
and various control flows. Unfortunately, the checker does not currently have
the ability to automatically convert such code to new-style arrow `switch`es.
Manually converting the code could be a good opportunity to improve its
readability.

How many potential execution paths can you spot?

```
enum Suit {HEARTS, CLUBS, SPADES, DIAMONDS};
private int foo(Suit suit){
 switch(suit) {
 case HEARTS:
 if (bar()) {
 break;
 }
 // Fall through
 case CLUBS:
 if (baz()) {
 return 1;
 } else if (baz2()) {
 throw new AssertionError(...);
 }
 // Fall through
 case SPADES:
 // Fall through
 case DIAMONDS:
 return 0;
 }
 return -1;
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("StatementSwitchToExpressionSwitch")` to the enclosing element.

Assigning to a static variable from a constructor is highly indicative of a bug, or error-prone design.

Common reasons are:

The field simply should be an instance field, and there’s a bug.

An attempt is being made to lazily initialize a static field. In this case,
first consider whether lazy initialization is necessary: it often isn’t. If
it is, doing it from a constructor is very hairy: the static field could be
accessed from a static method before the class is even initialized. Consider
using a memoized `Supplier`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("StaticAssignmentInConstructor")` to the enclosing element.

The problem we’re trying to prevent is unhelpful stack traces that don’t contain information about where the Exception was thrown from. This problem can sometimes arise when an attempt is being made to cache or reuse a Throwable (often, a particular Exception). In this case, consider whether this is really is necessary: it often isn’t. Could a Throwable simply be instantiated when needed?

```
// this always has the same stack trace
static final MyException MY_EXCEPTION = new MyException("something terrible has happened!");
```
```
throw new MyException("something terrible has happened!");
```
```
static MyException myException() {
 return new MyException("something terrible has happened!");
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("StaticAssignmentOfThrowable")` to the enclosing element.

Protecting writes to a static field by synchronizing on an instance lock is not thread-safe.

In the following example, two difference instances of `Test` can each acquire
their own instance lock and call `initialize()` at the same time.

```
class Test {
 private final Object lock = new Object();
 static initialized = false;
 static void initialize() { /* ... */ }
 Test() {
 synchronized (lock) {
 if (!initialized) {
 initialize();
 // error: modification of static variable guarded by instance variable 'lock'
 initialized = true;
 }
 // ...
 }
 }
}
```
Static fields should generally be guarded by static locks, and instance fields guarded by instance locks.

The example above could be made thread-safe by locking on the enclosing `Class`:

```
synchronized (Test.class) {
 if (!initialized) {
 initialize();
 initialized = true;
 }
}
```
To update a static counter from an instance method, consider using
`AtomicInteger`
instead of incrementing a static `int` field.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("StaticGuardedByInstance")` to the enclosing element.

Making a `@Mock` instance `static` will share the state across tests and can
make them order dependent and/or unclear to reader/maintainer. Removing
`static`, will ensure fresh mock instances are used for each tests.

Additionally, if the `@Mock` instance is marked as `static` to make it
serializable, there is a cleaner way to do it :
https://javadoc.io/doc/org.mockito/mockito-core/latest/org/mockito/Mockito.html#serilization_across_classloader
and
https://javadoc.io/static/org.mockito/mockito-core/3.3.3/org/mockito/Mock.html#serializable–

Suppress false positives by adding the suppression annotation `@SuppressWarnings("StaticMockMember")` to the enclosing element.

*Alternate names: FilesLinesLeak*

The problem is described in the javadoc for `Files`.

`Files.newDirectoryStream`When not using the try-with-resources construct, then directory stream’s close method should be invoked after iteration is completed so as to free any resources held for the open directory.

`Files.list`The returned stream encapsulates a

`DirectoryStream`. If timely disposal of file system resources is required, the try-with-resources construct should be used to ensure that the stream’s close method is invoked after the stream operations are completed.

`Files.walk`The returned stream encapsulates one or more

`DirectoryStreams`. If timely disposal of file system resources is required, the try-with-resources construct should be used to ensure that the stream’s close method is invoked after the stream operations are completed. Operating on a closed stream will result in an`IllegalStateException`.

`Files.find`The returned stream encapsulates one or more

`DirectoryStreams`. If timely disposal of file system resources is required, the try-with-resources construct should be used to ensure that the stream’s close method is invoked after the stream operations are completed. Operating on a closed stream will result in an`IllegalStateException`.

`Files.lines`The returned stream encapsulates a

`Reader`. If timely disposal of file system resources is required, the try-with-resources construct should be used to ensure that the stream’s close method is invoked after the stream operations are completed.

To ensure the stream is closed, always use try-with-resources. For example, when
using `Files.lines`, do this:

```
String input;
try (Stream<String> stream = Files.lines(path)) {
 input = stream.collect(Collectors.joining(", "));
}
```
Not this:

```
// the Reader is never closed!
String input = Files.lines(path).collect(Collectors.joining(", ");
```
Methods that return `Stream`s that encapsulate a closeable resource can be
annotated with `com.google.errorprone.annotations.MustBeClosed` to ensure their
callers remember to close the stream.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("StreamResourceLeak")` to the enclosing element.

Using `stream::iterator` to create an `Iterable` results in an `Iterable` which
can only be iterated once, and will throw an `IllegalArgumentException` on
subsequent attempts.

The contract of `Iterable` is poorly defined, but many APIs will assume that
`Iterable`s allow multiple iteration.

To give a concrete example using protocol buffers,

```
MyProto construct(Stream<Integer> ints) {
 return MyProto.newBuilder()
 .addAllSubMessage(
 ints.map(i -> SubMessage.newBuilder().setVal(i).build())::iterator)
 .build();
}
```
With the current implementation of the protocol buffer generated code, this will work, but the following minor change will lead to re-iteration and failure,

```
MyProto construct(Stream<Integer> ints) {
 MyProto.Builder builder = MyProto.newBuilder();
 builder.addSubMessageBuilder().setVal(0);
 return builder
 .addAllSubMessage(
 ints.map(i -> SubMessage.newBuilder().setVal(i).build())::iterator)
 .build(); // build iterates twice, and throws.
}
```
To avoid such pitfalls, the `Stream` can either be terminated with
`forEachOrdered`

```
MyProto construct(Stream<Integer> ints) {
 MyProto.Builder builder = MyProto.newBuilder();
 ints.map(i -> SubMessage.newBuilder().setVal(i).build())
 .forEachOrdered(builder::addSubMessage);
 return builder.build();
}
```
or collected (with the caveat that the contents of the `Stream` will now be
materialized into memory at once),

```
MyProto construct(Stream<Integer> ints) {
 return MyProto.newBuilder()
 .addAllSubMessage(
 ints.map(i -> SubMessage.newBuilder().setVal(i).build())
 .collect(toImmutableList()))
 .build();
}
```
or suppressed using `@SuppressWarnings("StreamToIterable")`. Only suppress if
you’re sure the API you’re using will only iterate the resulting one-shot
`Iterable` once, and can’t accept a `Stream` or an `Iterator`.

`String.toLowerCase()` (and `toUpperCase`) without specifying a `Locale` can
have surprising results.

For example, if this is used on a device and the user’s `Locale` is Türkiye,
then `"I".toLowerCase()` will yield a lowercase dotless I (“ı”). This could be
extremely dangerous if you were expecting to operate on ASCII text to generate
machine-readable identifiers.

If this kind of regionalisation is desired, use
`.toLowerCase(Locale.getDefault())` to make that explicit. If not,
`.toLowerCase(Locale.ROOT)` or `.toLowerCase(Locale.ENGLISH)` will give you
casing independent of the user’s current `Locale`. If you know that you’re
operating on ASCII, prefer `Ascii.toLower/UpperCase` to make that explicit.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("StringCaseLocaleUsage")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("StringCharset") to the enclosing element.

@SuppressWarnings("StringCharset")

Using text blocks for strings that span multiple lines can make code easier to read.

For example, prefer this:

```
String message =
 """
 'The time has come,' the Walrus said,
 'To talk of many things:
 Of shoes -- and ships -- and sealing-wax --
 Of cabbages -- and kings --
 And why the sea is boiling hot --
 And whether pigs have wings.'
 """;
```
instead of this:

```
String message =
 "'The time has come,' the Walrus said,\n"
 + "'To talk of many things:\n"
 + "Of shoes -- and ships -- and sealing-wax --\n"
 + "Of cabbages -- and kings --\n"
 + "And why the sea is boiling hot --\n"
 + "And whether pigs have wings.'\n";
```
If the string should not contain a trailing newline, use a `\ ` to escape the
final newline in the text block. That is, these two strings are equivalent:

```
String s = "hello\n" + "world";
```
```
String s =
 """
 hello
 world\
 """;
```
The suggested fixes for this check preserve the exact contents of the original
string, so if the original string doesn’t include a trailing newline the fix
will use a `\ ` to escape the last newline.

If the whitespace in the string isn’t significant, for example because the
string value will be parsed by a parser that doesn’t care about the trailing
newlines, consider removing the final `\ ` to improve the readability of the
string.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("StringConcatToTextBlock")` to the enclosing element.

`String.split(String)` and `Pattern.split(CharSequence)` have surprising
behaviour. For example, consider the following puzzler from
https://konigsberg.blogspot.com/2009/11/final-thoughts-java-puzzler-splitting.html:

```
String[] nothing = "".split(":");
String[] bunchOfNothing = ":".split(":");
```
The result is `[""]` and `[]`!

More examples:

| input | `input.split(":")` | `Pattern.compile(":").split(input)` | `Splitter.on(':').split(input)` |
|---|---|---|---|
| `""` | `[""]` | `[""]` | `[""]` |
| `":"` | `[]` | `[]` | `["", ""]` |
| `":::"` | `[]` | `[]` | `["", "", "", ""]` |
| `"a:::"` | `["a"]` | `["a"]` | `["a", "", "", ""]` |
| `":::b"` | `["", "", "", "b"]` | `["", "", "", "b"]` | `["", "", "", "b"]` |

Prefer either:

Guava’s
`Splitter`,
which has less surprising behaviour and provides explicit control over the
handling of empty strings and the trimming of whitespace with `trimResults`
and `omitEmptyStrings`.

`String.split(String, int)`
or
`Pattern.split(CharSequence, int)`
and setting an explicit ‘limit’ to `-1` to match the behaviour of
`Splitter`.

TIP: if you use `Splitter`, consider extracting the instance to a `static`
`final` field.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("StringSplitter")` to the enclosing element.

*Alternate names: SuperEqualsIsObjectEquals*

Implementations of `equals` and `hashCode` should usually not delegate to
`Object.equals` and `Object.hashCode`.

Those two methods implement equality based on object identity. That
implementation *is* sometimes what the author intended. (This check attempts to
identify those cases and *not* report a warning for them. But in some cases, it
still produces a warning when it shouldn’t.)

But when `super.equals` and `super.hashCode` call the methods defined in
`Object`, the developer often did *not* intend to use object identity. Often,
developers write something like:

```
private final int id;
@Override
public boolean equals(Object obj) {
 if (obj instanceof Foo) {
 return super.equals(obj) && id == ((Foo) obj).id;
 }
 return false;
}
@Override
public int hashCode() {
 return super.hashCode() ^ id;
}
```
This appears to be an attempt to define equality in terms of the `id` field in
this class and any fields in the superclass. However, when the superclass that
defines `equals` or `hashCode` is `Object`, the code instead defines equality in
terms of a mix of object identity and field values. The result is equivalent to
defining it in terms of identity alone—which is equivalent to not overriding
`equals` and `hashCode` at all!

Typically, the code should be rewritten to remove the `super` calls entirely:

```
private final int id;
@Override
public boolean equals(Object obj) {
 if (obj instanceof Foo) {
 return id == ((Foo) obj).id;
 }
 return false;
}
@Override
public int hashCode() {
 return id;
}
```
Note that the suggested edits for this check instead preserve behavior, which
likely means preserving bugs! However, in cases in which object identity *is*
intended, we recommend applying the suggested edit to make that behavior
explicit in the code:

```
// This class's definition of equality is unusual and perhaps not ideal.
// But it is at least explicit.
private final Integer id;
@Override
public boolean equals(Object obj) {
 if (obj instanceof Foo) {
 if (id == null) {
 return this == obj;
 }
 return id.equals(((Foo) obj).id);
 }
 return false;
}
@Override
public int hashCode() {
 return id != null ? id : System.identityHashCode(this);
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("SuperCallToObjectMethod")` to the enclosing element.

SWIG is a tool that will automatically generate Java bindings to C++ code. It is possible to %ignore in SWIG a C++ object’s destructor, this is trivially achieved when using %ignoreall and then selectively %unignore-ing an API. SWIG cleans up C++ objects using a delete method, which is most commonly called by a finalizer. When a SWIG generated delete method can’t call a destructor, as it is hidden, the delete method throws an exception. However, in the case of a hidden C++ destructor SWIG also doesn’t generate a finalizer, and so the most common call to the delete method is removed. The consequence of this is that the SWIG objects leak their C++ counterpart and no warnings or exceptions are thrown.

This check looks for the pattern of a memory leaking SWIG generated object and warns about the potential memory leak. The most straightforward fix is to the SWIG input code to tell it not to %ignore the C++ code’s destructor.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("SwigMemoryLeak")` to the enclosing element.

Possible fixes:

If the field is never reassigned, add the missing `final` modifier.

If the field needs to be mutable, create a separate lock by adding a private final field and synchronizing on it to guard all accesses.

If the field is lazily initialized, annotation it with
`com.google.errorprone.annotations.concurrent.LazyInit`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("SynchronizeOnNonFinalField")` to the enclosing element.

Null-checking `System.console()` is not a reliable way to detect if the console
is connected to a terminal.

See JDK 22 Release Note: JLine As The Default Console Provider:

`System.console()`now returns a`Console`object when the standard streams are redirected or connected to a virtual terminal. In prior releases,`System.console()`returned`null`for these cases. This change may impact code that uses the return from`System.console()`to test if the VM is connected to a terminal. If needed, running with`-Djdk.console=java.base`will restore older behavior where the console is only returned when it is connected to a terminal.

A new method

`Console.isTerminal()`has been added to test if console is connected to a terminal.

and JDK 25 release note Release Note: Default Console Implementation No Longer Based On JLine:

The default Console obtained via

`System.console()`is no longer based on JLine. Since JDK 20, the JDK has included a JLine-based Console implementation, offering a richer user experience and better support for virtual terminal environments, such as IDEs. This implementation was initially opt-in via a system property in JDK 20 and JDK 21 and became the default in JDK 22. However, maintaining the JLine-based Console proved challenging. As a result, in JDK 25, it has reverted to being opt-in, as it was in JDK 20 and JDK 21.

To prepare for this change while remaining compatible with JDK versions prior to
JDK 22, consider using reflection to call `Console#isTerminal` on JDK versions
that support it:

```
 @SuppressWarnings("SystemConsoleNull") // https://errorprone.info/bugpattern/SystemConsoleNull
 private static boolean systemConsoleIsTerminal() {
 Console systemConsole = System.console();
 if (Runtime.version().feature() < 22 || systemConsole == null) {
 return systemConsole != null;
 }
 try {
 return (Boolean) Console.class.getMethod("isTerminal").invoke(systemConsole);
 } catch (ReflectiveOperationException e) {
 throw new LinkageError(e.getMessage(), e);
 }
 }
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("SystemConsoleNull")` to the enclosing element.

Thread.join() can be interrupted, and so requires users to catch InterruptedException. Most users should be looping until the join() actually succeeds.

Instead of writing your own try-catch and loop to handle it properly, you may
use **Uninterruptibles.joinUninterruptibly** which does the same for you.

Example:

```
Thread thread = new Thread(new Runnable() {...});
Uninterruptibles.joinUninterruptibly(thread);
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ThreadJoinLoop")` to the enclosing element.

`ThreadLocal`s should be stored in `static` variables to avoid memory leaks. If
a `ThreadLocal` is stored in an instance (non-static) variable, there will be
`M \* N` instances of the `ThreadLocal` value where `M` is the number
of threads, and `N` is the number of instances of the containing class. Each
instance may remain live as long the thread that stored it stays live.

Example:

```
class C {
 private final ThreadLocal<D> local = new ThreadLocal<D>();
 public f() {
 D d = local.get();
 if (d == null) {
 d = new D(this);
 local.set(d);
 }
 d.doSomething();
 }
}
```
The fix is often to make the field `static`:

```
private static final ThreadLocal<D> local = new ThreadLocal<D>();
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ThreadLocalUsage")` to the enclosing element.

Don’t rely on the thread scheduler for correctness or performance. Instead, ensure that the average number of runnable threads is not significantly greater than the number of processors, i.e. by using the executor framework and an appropriately sized thread pool.

For more information, see Effective Java 3rd Edition §84.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ThreadPriorityCheck")` to the enclosing element.

According to the Javadoc of `java.util.TimeZone`:

For compatibility with JDK 1.1.x, some other three-letter time zone IDs (such as “PST”, “CTT”, “AST”) are also supported. However, their use is deprecated because the same abbreviation is often used for multiple time zones (for example, “CST” could be U.S. “Central Standard Time” and “China Standard Time”), and the Java platform can then only recognize one of them.

Aside from the ambiguity between timezones, there is inconsistency in the
observance of Daylight Savings Time for the returned time zone, meaning the
`TimeZone` obtained may not be what you expect. Examples include:

`DateTime.getTimeZone("PST")` does observe daylight savings time; however,
the identifier implies that it is Pacific `DateTime.getTimeZone("EST")` (and `"MST"` and `"HST"`) do not observe
daylight savings time. However, this is inconsistent with PST (and others),
so you may believe that daylight savings time will be observed.This check will only suggest replacements which yield the same rules as the existing three-letter ID for at least part of the year (e.g. it will suggest “America/Chicago” and “Etc/GMT+6” but not “Asia/Shanghai” as a replacement for “CST”).

Certain 3-letter time zone IDs are not flagged by this check, specifically if
the ID appears in `ZoneId.getAvailableZoneIds()`, e.g. “UTC”, “GMT”, “PRC”.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ThreeLetterTimeZoneID")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("ThrowIfUncheckedKnownUnchecked") to the enclosing element.

@SuppressWarnings("ThrowIfUncheckedKnownUnchecked")

Throwables should not override `equals()` or `hashCode()`.

**Exceptions are Events, Not Value Objects** Philosophically, an exception
represents a unique historical event: something went wrong at a specific
time, in a specific thread, at a specific line of code. It is not a data
container or a value type (like a `String` or a `Money` object). Even if two
`IllegalArgumentException`s are thrown with the exact same message (```
"ID
cannot be null"
```
), and perhaps even identical stack traces, they represent
two distinct failures that happened independently. Treating them as “equal”
conceptually conflates two different events.

**The Stack Trace Problem** When an exception is instantiated (or thrown),
Java populates its stack trace via `fillInStackTrace()`.

`equals()`
comparison? If you do, comparing arrays of `StackTraceElement` is
computationally expensive. Exceptions can also have causes and
suppressed exceptions. This adds expense, too. Also, will all the
transitive causes and suppressed exceptions themselves implement
`equals()` the way you want? Plus, causes and suppressed exceptions make
exceptions mutable, and mutable objects generally shouldn’t implement
`equals()`.`ServiceA` is equal to an exception thrown in
`ServiceB` just because they share a message or an error code. This
masks critical debugging context.**It Hides Bad Architecture** The primary reason you would need to override
these methods is if you are placing Exceptions into a `HashSet`, or using
them as keys in a `HashMap`. If you are doing this, you are likely using
Exceptions for normal business logic or control flow, which is a known
anti-pattern. Exceptions are for exceptional circumstances; they are heavy
(because of the stack trace) and slow to generate.

**Java Still Sometimes Compares Exceptions by Object Identity** Even if two
exceptions are “the same” according to `equals`, Java will still print both
if it encounters them during a call to `printStackTrace`, and it will allow
one to be a suppressed exception of the other. In both ways, two exceptions
that are “the same” will continue to be treated as different by Java.

If you find yourself needing to compare exceptions, **extract the state into a
separate value object.**

Instead of this:

```
public class UserNotFoundException extends RuntimeException {
 private final String userId;
 private final String groupId;
 // Override equals and hashCode based on userId and groupId...
}
```
Do this: Create a custom Exception that *contains* a value object (like a Record
or a POJO), and compare the value objects instead.

```
public class UserNotFoundException extends RuntimeException {
 private final ErrorDetails details;
 public UserNotFoundException(ErrorDetails details) {
 super("%s is not a member of %s:".formatted(details.userId(), details.groupId()));
 this.details = details;
 }
 public ErrorDetails getDetails() { return details; }
}
// Compare the details, not the exceptions!
if (e1.getDetails().equals(e2.getDetails())) { ... }
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ThrowableEqualsHashCode")` to the enclosing element.

#
 TimeInStaticInitializer

 Accessing the current time in a static initialiser captures the time at class loading, which is rarely desirable.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("TimeInStaticInitializer")` to the enclosing element.

This checker flags potential problems with TimeUnit conversions: 1) conversions that are statically known to be equal to 0 or 1; 2) conversions that are converting from a given unit back to the same unit; 3) conversions that are converting from a smaller unit to a larger unit and passing a constant value

Suppress false positives by adding the suppression annotation `@SuppressWarnings("TimeUnitConversionChecker")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("ToStringReturnsNull") to the enclosing element.

@SuppressWarnings("ToStringReturnsNull")

The newer arrow (`->`) syntax for switches is preferred to the older colon (`:`)
syntax. The main reason for continuing to use the colon syntax in switch
*statements* is that it allows fall-through from one statement group to the
next. But in a switch *expression*, fall-through would only be useful if the
code that falls through has side effects. Burying side effects inside a switch
expression makes code hard to understand.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("TraditionalSwitchExpression")` to the enclosing element.

Arguments to a fluent Truth assertion appear to be reversed based on the argument names.

```
 int expected = 1;
 assertThat(expected).isEqualTo(codeUnderTest());
```
This is problematic as the quality of Truth’s error message depends on the
argument order. If `codeUnderTest()` returns `2`, this code will output:

```
expected: 2
but was : 1
```
Which will likely make debugging the problem harder. Truth assertions should follow the opposite order to JUnit assertions. Compare:

```
 assertThat(actual).isEqualTo(expected);
 assertEquals(expected, actual);
```
See https://truth.dev/faq#order for more details.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("TruthAssertExpected")` to the enclosing element.

The arguments to assertThat method is a constant. It should be a variable or a method invocation. For eg. switch assertThat(1).isEqualTo(methodCall()) to assertThat(methodCall()).isEqualTo(1).

Suppress false positives by adding the suppression annotation `@SuppressWarnings("TruthConstantAsserts")` to the enclosing element.

Expectation of ```
Truth.assertThat(map.getOrDefault(key,
defaultValue)).isEqualTo(expectedValue)
```
 is unclear if the `defaultValue` is
same as `expectedValue`. If the test passes, its hard to say if `map` contained
`key, expectedValue` as an entry. Most likely, developer intended to verify that
`map` doesn’t contain `'key` or perhaps map `key` isn’t mapped to
`defaultValue`.

Additionally, same assertion can be simplified if `defaultValue` and
`expectedValue` are different constants to
`Truth.assertThat(map.get(key)).isEqualTo(expectedValue)`.

That is, prefer this:

```
public static void doSomething(Map<String, String> map, String key, String expectedValue) {
 assertThat(map.get(key)).isEqualTo(expectedValue);
 assertThat(map).doesNotContainKey(key);
 assertThat(map).containsEntry(key, expectedValue);
}
```
to this:

```
public static void doSomething(Map<String, String> map, String key, String defaultValue, String expectedValue) {
 assertThat(map.getOrDefault(key, defaultValue)).isEqualTo(expectedValue);
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("TruthGetOrDefault")` to the enclosing element.

This assertion using Truth is guaranteed to be true or false based on the types being passed to it. For example,

```
Optional<Integer> x = Optional.of(1);
assertThat(x).isNotEqualTo(Optional.of(2L));
```
```
ImmutableList<Long> xs = ImmutableList.of(1L, 2L);
assertThat(xs).doesNotContain(1);
```
These will always be true, given `Integer` is not comparable to `Long`. This
isn’t such a big issue for `isEqualTo` assertions, given the test will fail.
However, it can be insidious for `isNotEqualTo` or `doesNotContain`, given the
assertion will be vacuously true without providing any test coverage.

One false positive for this is where the goal is to test the `equals`
implementation of the type under test, i.e. that comparison to a different type
is `false` and does not throw:

```
assertThat(myCustomType).isNotEqualTo("");
```
For such cases, consider whether your type can be implemented using `AutoValue`
to remove the need to implement `equals` by hand. If it can’t, consider
`EqualsTester`.

```
new EqualsTester()
 .addEqualityGroup(MyCustomType.create(1), MyCustomType.create(1))
 .addEqualityGroup("")
 .testEquals();
```
Although consider omitting an explicit comparison with a different type, as
`EqualsTester` does this already by default.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("TruthIncompatibleType")` to the enclosing element.

`TypeMirror`
doesn’t override `Object.equals` and instances are not interned by javac, so
testing types for equality should be done with
`Types#isSameType`
instead.

If you’re implementing an Error Prone `BugChecker`, you can get a `Types`
instance from `VisitorState`.

If you’re implementing `AnnotationProcessor`, you can get the `Types` instance
from `javax.annotation.processing.ProcessingEnvironment`.

For more discussion of preferred APIs for comparing types, see https://errorprone.info/bugpattern/TypeToString

Suppress false positives by adding the suppression annotation `@SuppressWarnings("TypeEquals")` to the enclosing element.

When declaring type parameters, it’s possible to declare a type parameter with the same name as another type in scope, “shadowing” that type and potentially causing confusing or unintended behavior.

```
class Bar {
 ...
 public void doSomething(T object) {
 // Here, object is the static class T in this file
 }
 public <T> void doSomethingElse(T object) {
 // Here, object is a generic T
 }
 ...
 public static class T {...}
}
```
This checker warns when a type parameter shadows another type and suggests a possible renaming for the type parameter.

Note, however, that in some cases it may be preferable to rename or delete the shadowed type rather than the type parameter shadowing it, such as in cases where the type parameter is always instantiated with the same type.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("TypeNameShadowing")` to the enclosing element.

When declaring type parameters, it’s possible to declare a type parameter with the same name as another type parameter in scope causing unintended incompatibilities when trying to use them together.

```
public class Foo<T> {
 void instanceMethod(T t) {}
 <T> void genericMethod(T t) {
 instanceMethod(t); // FAIL: T declared in this method doesn't correspond to Foo<T>'s T
 }
}
```
In some cases, the type variable being declared has no relation to the type variable being shadowed. It may be appropriate to rename the shadowing type variable:

```
class Logger<T> {
 void log(T t) {}
 <T> logOther(T t) { ... } // Really should be <O> void logOther(O o), since this T is unrelated.
}
```
Depending on the nature of the surrounding code, you might be able to remove the generic declaration on a method, or convert the generic method into a static method that doesn’t inherit the surrounding type parameter:

```
class Holder<T> {
 T held;
 public Foo<T> fooIt() {
 return fooify(held);
 }
 private <T> Foo<T> fooify(T t) { ... } // Could be static, or non-generic
}
```
If an inner class declaration shadows a type variable, you may be able to remove the type variable, make it a static inner class, or rename the type variable:

```
class BoxingBox<T> {
 // Works if you make the class static.
 // If you remove the type parameter, you'll need to update stuff to List<Container>
 class Container<T> {
 T held;
 }
 List<Container<T>> stuff = new ArrayList<>();
 T get(int index) {
 return stuff.get(index).held;
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("TypeParameterShadowing")` to the enclosing element.

A method’s type parameters should always be referenced in the declaration of one or more formal parameters. Type parameters that are only used in the return type are a source of type-unsafety. First, operations on the type will be unchecked after the type parameter is erased. For example:

```
static <T> T doCast(Object o) {
 return (T) o; // this will always succeed, since T is erased
}
```
The ‘doCast’ method would be better implemented as:

```
static <T> T doCast(Class<T> clazz, Object o) {
 return clazz.cast(o); // has the expected behaviour
}
```
Second, this pattern causes unsafe casts to occur at invocations of the method. Consider the following snippet, which uses the first (incorrect) implementation of ‘doCast’:

```
this.<String>doCast(42); // succeeds
String s = doCast(42); // fails at runtime
```
Finally, relying on the type parameter to be inferred can have surprising results, and interacts badly with overloaded methods. Consider:

```
<T> T getThing()
void assertThat(int a, int b)
void assertThat(Object a, Object b)
```
This invocation will be ambiguous:

```
// both method assertThat(int,int) and method assertThat(Object,Object) match
assertThat(42, getThing());
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("TypeParameterUnusedInFormals")` to the enclosing element.

The `equals` and `hashCode` methods of `java.net.URL` make blocking network
calls. When you place a `URL` into a hash-based container, the container invokes
those methods.

Prefer `java.net.URI`. Or, if you must use `URL` in a
collection, prefer to use a non-hash-based container like a `List<URL>`, and
avoid calling methods like `contains` (which calls `equals`) on it.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("URLEqualsHashCode")` to the enclosing element.

This code uses `Object.equals` (or similar method) with a type that does not
have well-defined `equals` behavior: `Collection`, `Iterable`, `Multimap`,
`Queue`, or `CharSequence`. Such a call to `equals` may return `false` in
cases where equality was expected. `SparseArray` and `LongSparseArray` do
not implement `equals`, so will fall back to reference equality.

`Collection` or `Iterable`Your code might be working correctly, but only if there is some *subtype* of
`Iterable` which *does* have well-defined equals behavior, and you are certain
that both operands are definitely of that type at runtime. (The common examples
of such types are [`List`], [`Set`], and `Multiset`, or any subtypes of
those.)

If this describes your situation, congratulations: you don’t have a bug. To make
this warning go away, change the references in your code to be of that more
specific static type, not `Collection` or `Iterable`. This lets the bug
checker (and human readers) *know* that there is no risk of a false negative.

The minimal solution is to cast or copy “just in time” before calling `equals`,
but ideally you can make broader changes, to adopt the proper interface more
widely in your code.

On the other hand, if you might be mixing a `List` and a non-`List`, etc., you
are at risk. If you can’t correct that situation, one alternative solution is to
use `Iterables.elementsEqual` (which checks for *order-dependent* equality,
like `List.equals`).

`Multimap`The discussion above generally applies in this case as well; the well-behaved
subtypes to choose from are `ListMultimap` and `SetMultimap`.

There is no equivalent to `Iterables.elementsEqual` in this case, however.

`Queue`All known `Queue` implementations besides `LinkedList` use reference equality
instead of value-based equality. *Some* of the workarounds discussed above may
apply.

`CharSequence`When comparing a `String` to a `CharSequence`, prefer `String#contentEquals`.
When comparing the content of two `CharSequence`s, you may want to compare the
string representation: `lhs.toString().contentEquals(rhs)`.

`SparseArray` and `LongSparseArray`These must be iterated over and compared manually, element by element.

`java.util.Date`Subtypes of `Date` (like `java.sql.Timestamp`) break substitutability, so
comparing `Date`s with `equals` is unreliable.

TIP: `java.util.Date` is a legacy, bug-prone API. Prefer `java.time.Instant` or
`java.time.LocalDateTime`.

`ImmutableCollection`Prefer subtypes such as `ImmutableSet` or `ImmutableList`, which have
well-defined `equals`.

`Future`Prefer calling `#get` and checking the result, or using
`assertThat(...).isSameInstanceAs()` or `assertThat(...).isNotSameInstanceAs()`
for cases where the same/different instance is expected.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("UndefinedEquals")` to the enclosing element.

Using unicode escapes in Java for printable characters is obfuscated. Worse,
given the compiler allows unicode literals outside of `String` literals, it can
be potentially unsafe.

Prefer using literal characters for printable characters.

For an example of malicious code, consider:

```
class Evil {
 public static void main(String... args) {
 // Don't run this, it would be really unsafe!
 // \u000d Runtime.exec("rm -rf /");
 }
}
```
`\u000d` encodes a newline character, so `Runtime.exec` appears on its own line
and will execute.

NOTE: Unicode escapes are defined as a preprocessing step in the Java compiler
(see JLS §3.3). After compilation, there is no runtime difference whatsoever
between a Unicode escape and using the equivalent character in source. That is,
writing `"hello \u0077\u006f\u0072\u006c\u0064"` is equivalent to ```
"hello
world"
```
 in the compiled `.class` file and at runtime.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnicodeEscape")` to the enclosing element.

Prefer using an unnamed variable (_) to denote variables and patterns that are intentionally unused.

_

Suppress false positives by adding the suppression annotation @SuppressWarnings("UnnamedVariable") to the enclosing element.

@SuppressWarnings("UnnamedVariable")

The `@Mock` annotation is used to automatically initialize mocks using
`MockitoAnnotations.initMocks`, or `MockitoJUnitRunner`.

Variables annotated this way should not be explicitly initialized, as this will be overwritten by automatic initialization.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnnecessaryAssignment")` to the enclosing element.

#
 UnnecessaryAsync

 Variables which are initialized and do not escape the current scope do not need to worry about concurrency. Using the non-concurrent type will reduce overhead and verbosity.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnnecessaryAsync")` to the enclosing element.

The newer arrow (`->`) syntax for switches does not permit fallthrough between
cases. A `break` statement is allowed to break out of the switch, but including
a `break` as the last statement in a case body is unnecessary.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnnecessaryBreakInSwitch")` to the enclosing element.

The collections returned by the protobuf API are unmodifiable: attempts to
modify them will result in an `UnsupportedOperationException`.

The protobuf collections are not implemented using Guava’s immutable collections, so copying them as a Guava collection isn’t free.

This check suggests omitting copies of unmodifiable protobuf collections as Guava immutable collections in places where a Guava immutable collection is not required, for example to replace this:

```
List<String> foos(MyProtoMessage message) {
 return ImmutableList.copyOf(message.getFoos());
}
```
with this:

```
List<String> foos(MyProtoMessage message) {
 return message.getFoos();
}
```
Note that the check will not report diagnostics if the result of the copy is being used somewhere that requires a Guava Immutable collection, for example it is returned from the current method, or passed to an API that expects a Guava immutable collection:

```
ImmutableList<String> foos(MyProtoMessage message) {
 return ImmutableList.copyOf(message.getFoos());
}
```
The check will suggest refactoring local variables, if the local variable is never used in a context that requires a Guava immutable collection. Using Guava’s collections has some benefit in this case, since it makes the immutability of the collection clear and can enable other compile-time static analysis for accidental attempts to modify the collection. But this value is limited if it is only used as a local variable that never escapes the current method, and there is some runtime cost from making the copy.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnnecessaryCopy")` to the enclosing element.

Prefer method references to constant lambda expressions, or to a helper method that does nothing but return a lambda.

That is, prefer this:

```
private static Bar getBar(Foo foo) {
 return BarService.lookupBar(foo, defaultCredentials());
}
...
return someStream().map(MyClass::getBar)....;
```
to this:

```
private static final Function<Foo, Bar> GET_BAR_FUNCTION =
 foo -> BarService.lookupBar(foo, defaultCredentials());
...
return someStream().map(GET_BAR_FUNCTION)....;
```
Advantages of using a method include:

`com.google.common.base.Function` or a
`java.util.function.Function`?)Be aware that this change is not purely syntactic: it affects the semantics of your program in some small ways. In particular, evaluating the same method reference twice is not guaranteed to return an identical object.

This means that, first, inlining the reference instead of storing a lambda may cause additional memory allocations - usually this very slight performance cost is worth the improved readability, but use your judgment if the performance matters to you.

Secondly, if the correctness of your program depends on reference equality of
your lambda, inlining it may break you. Ideally, you should *not* depend on
reference equality for a lambda, but if you are doing so, consider not making
this change.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnnecessaryLambda")` to the enclosing element.

#
 UnnecessaryLongToIntConversion

 Converting a long or Long to an int to pass as a long parameter is usually not necessary. If this conversion is intentional, consider `Longs.constrainToRange()` instead.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnnecessaryLongToIntConversion")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("UnnecessaryMethodInvocationMatcher") to the enclosing element.

@SuppressWarnings("UnnecessaryMethodInvocationMatcher")

Using a method reference to refer to the abstract method of the target type is unnecessary. For example,

```
Stream<Integer> filter(Stream<Integer> xs, Predicate<Integer> predicate) {
 return xs.filter(predicate::test);
}
```
```
Stream<Integer> filter(Stream<Integer> xs, Predicate<Integer> predicate) {
 return xs.filter(predicate);
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnnecessaryMethodReference")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("UnnecessaryParentheses") to the enclosing element.

@SuppressWarnings("UnnecessaryParentheses")

A `@Qualifier` or a `@BindingAnnotation` has no effect here, and can be removed.
Its presence may be misleading.

For example:

```
final class MyInjectableClass {
 @Username private final String username;
 @Inject
 MyInjectableClass(@Username String username) {
 this.username = username;
 }
}
```
The annotation on the constructor parameter is important, but the field annotation is redundant.

```
final class MyInjectableClass {
 private final String username;
 @Inject
 MyInjectableClass(@Username String username) {
 this.username = username;
 }
}
```
There are a couple of ways this check can lead to false positives:

You’re using a custom framework we don’t know about which makes the location of the finding an injection point. File a bug, and we’ll happily incorporate it.

Your annotation is annotated with `@Qualifier` or `@BindingAnnotation` but
isn’t actually used as a qualifier (perhaps you have a framework that just
uses it reflectively). Try removing those annotations from the *annotation*.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnnecessaryQualifier")` to the enclosing element.

#
 UnnecessaryStringBuilder

 Prefer string concatenation over explicitly using `StringBuilder#append`, since `+` reads better and has equivalent or better performance.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnnecessaryStringBuilder")` to the enclosing element.

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
UnrecognisedJavadocTag
This Javadoc tag wasn't recognised by the parser. Is it malformed somehow, perhaps with mismatched braces?
Severity
WARNING

When invoking `static native` methods from objects that override
`finalize`,
the finalizer may run before the native method finishes execution.

Consider the following example:

```
public class GameRunner {
 // Pointer to some native resource, stored as a long
 private long nativeResourcePtr;
 // Allocates the native resource that this object "owns"
 private static native long doNativeInit();
 // Releases the native resource
 private static native void cleanUpNativeResources(long nativeResourcePtr);
 public GameRunner() {
 nativeResourcePtr = doNativeInit();
 }
 @Override
 protected void finalize() {
 cleanUpNativeResources(nativeResourcePtr);
 nativeResourcePtr = 0;
 super.finalize();
 }
 public void run() {
 GameLibrary.playGame(nativeResourcePtr); // Bug!
 }
}
public class GameLibrary {
 // Plays the game using the native resource
 public static native void playGame(long nativeResourcePtr);
}
```
During the execution of `GameRunner.run`, the call to `playGame` may not hold
the `this` reference live, and its finalizer may run, cleaning up the native
resources while the native code is still executing.

You can fix this by making the `static native` method not `static`, or by
changing the `static native` method so that it is passed the enclosing instance
and not the `nativeResourcePtr` value directly.

If you can use Java 9, the new method
`Reference.reachabilityFence`
will keep the reference passed in live. Be sure to call it in a finally block:

```
 public void run() {
 try {
 playAwesomeGame(nativeResourcePtr);
 } finally {
 Reference.reachabilityFence(this);
 }
 }
}
```
Note: This check doesn’t currently detect this problem when passing fields from objects other than “this”. That is equally unsafe. Consider passing the Java object instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnsafeFinalization")` to the enclosing element.

Prefer `asSubclass` instead of casting the result of `newInstance` to detect
classes of incorrect type before invoking their constructors. This way, if the
class is of the incorrect type, it will throw an exception before invoking its
constructor.

```
(Foo) Class.forName(someString).getDeclaredConstructor(...).newInstance(args);
```
Should be written as

```
Class.forName(someString).asSubclass(Foo.class).getDeclaredConstructor(...).newInstance();
```
This has caused issues in the past:

CVE-2014-7911 - https://seclists.org/fulldisclosure/2014/Nov/51

Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnsafeReflectiveConstructionCast")` to the enclosing element.

Thread-safe methods should never be overridden by methods that are not thread-safe. Doing so violates behavioural subtyping, and can result in bugs if the subtype is used in contexts that rely on the thread-safety of the supertype.

Overriding a `synchronized` method with a method that is not `synchronized` can
be a sign that the thread-safety of the supertype is not being preserved.

```
class Counter {
 private int count = 0;
 synchronized void increment() {
 count++;
 }
}
```
```
// MyCounter is not thread safe!
class MyCounter extends Counter {
 private int count = 0;
 void increment() {
 count++;
 }
}
```
Note that there are many ways to implement a thread-safe method without using
the `synchronized` modifier (e.g. `synchronized` statements using explicit
locks, or other locking constructs). When overriding a `synchronized` method
with a method that is thread-safe but does not have the `synchronized` modifier,
consider adding `@SuppressWarnings("UnsynchronizedOverridesSynchronized")` and
an explanation.

```
class MyCounter extends Counter {
 private AtomicInteger count = AtomicInteger();
 @SuppressWarnings("UnsynchronizedOverridesSynchronized") // AtomicInteger is thread-safe
 void increment() {
 count.getAndIncrement();
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("UnsynchronizedOverridesSynchronized")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("UnusedLabel") to the enclosing element.

@SuppressWarnings("UnusedLabel")

*Alternate names: Unused, unused, UnusedParameters*

The presence of an unused method may indicate a bug. This check highlights
*private* methods which are unused and can be safely removed without considering
the impact on other source files.

Methods and fields which are used by reflection can be annotated with `@Keep` to
suppress the warning.

This annotation can also be applied to annotations, to suppress the warning for any member annotated with that annotation:

```
import com.google.errorprone.annotations.Keep;
@Keep
@Retention(RetentionPolicy.RUNTIME)
@interface Field {}
...
public class Data {
 @Field private int a; // no warning.
 ...
}
```
All false positives can be suppressed by annotating the method with
`@SuppressWarnings("unused")` or prefixing its name with `unused`.

Toggle navigation
Error Prone
Bug Patterns
Docs
GitHub
UnusedNestedClass
This nested class is unused, and can be removed.
Severity
WARNING
Alternate names: unused

Suppress false positives by adding the suppression annotation @SuppressWarnings("UnusedTypeParameter") to the enclosing element.

@SuppressWarnings("UnusedTypeParameter")

*Alternate names: unused, UnusedParameters*

The presence of an unused variable may indicate a bug. This check highlights private fields, and parameters of private methods, which are unused and can be safely removed without considering the impact on other source files. “Private” in this context also includes effectively-private members, like public members of private classes.

False positives on fields and parameters can be suppressed by prefixing the
variable name with `unused`, e.g.:

```
private static void authenticate(User user, Application unusedApplication) {
 checkState(user.isAuthenticated());
}
```
Fields which are used by reflection can be annotated with `@Keep` to suppress
the warning.

This annotation can also be applied to annotations, to suppress the warning for any member annotated with that annotation:

```
import com.google.errorprone.annotations.Keep;
@Keep
@Retention(RetentionPolicy.RUNTIME)
@interface Field {}
...
public class Data {
 @Field private int a; // no warning.
 ...
}
```
All false positives can be suppressed by annotating the variable with
`@SuppressWarnings("unused")`.

A @Provides or @Produces method that returns its single parameter has long been Dagger’s only mechanism for delegating a binding. Since the delegation is implemented via a user-defined method there is a disproportionate amount of overhead for such a conceptually simple operation. @Binds was introduced to provide a declarative way of delegating from one binding to another in a way that allows for minimal overhead in the implementation. @Binds should always be preferred over @Provides or @Produces for delegation.

For instance, the following `@Provides` method

```
@Provides static Heater provideHeater(ElectricHeater heater) {
 return heater;
}
```
is equivalent to the following preferred `@Binds` method.

```
@Binds abstract Heater bindHeater(ElectricHeater impl);
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("UseBinds")` to the enclosing element.

When a field/variable name is the same as the field/variable type, it is difficult to determine which to use at which time.

For example,

```
private static String String;
```
This would cause future use of String.something within this class to refer to the static field String, instead of the class String.

This is worth calling out to avoid confusion and is a violation of Google Java style naming conventions

Instead of this naming style, the correct way would be:

```
private static String string;
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("VariableNameSameAsType")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("VoidUsed") to the enclosing element.

@SuppressWarnings("VoidUsed")

`Object.wait()` is supposed to block until either another thread invokes the
`Object.notify()` or `Object.notifyAll()` method, or a specified amount of time
has elapsed. The various `Condition.await()` methods have similar behavior.
However, it is possible for a thread to wake up without either of those
occurring; these are called *spurious wakeups*.

Because of spurious wakeups, `Object.wait()` and `Condition.await()` must always
be called in a loop. The correct fix for this varies depending on what you are
trying to do.

The incorrect code for this typically looks like:

Thread 1:

```
synchronized (this) {
 if (!condition) {
 wait();
 }
 doStuffAssumingConditionIsTrue();
}
```
Thread 2:

```
synchronized (this) {
 condition = true;
 notify();
}
```
If the call to `wait()` unblocks because of a spurious wakeup, then
`doStuffAssumingConditionIsTrue()` will be called even though `condition` is
still false. Instead of the `if`, you should use a `while`:

Thread 1:

```
synchronized (this) {
 while (!condition) {
 wait();
 }
 doStuffAssumingConditionIsTrue();
}
```
This ensures that you only proceed to `doStuffAssumingConditionIsTrue()` if
`condition` is true. Note that the check of the condition variable must be
inside the synchronized block; otherwise you will have a race condition between
checking and setting the condition variable.

The incorrect code for this typically looks like:

Thread 1:

```
synchronized (this) {
 wait();
 doStuffAfterEvent();
}
```
Thread 2:

```
// when event occurs
synchronized (this) {
 notify();
}
```
If the call to `wait()` unblocks because of a spurious wakeup, then
`doStuffAfterEvent()` will be called even though the event has not yet occurred.
You should rewrite this code so that the occurrence of the event sets a
condition variable as well as calls `notify()`, and the `wait()` is wrapped in a
while loop checking the condition variable. That is, it should look just like
the previous example.

The incorrect code for this typically looks like:

```
synchronized (this) {
 if (!condition) {
 wait(timeout);
 }
 doStuffAssumingConditionIsTrueOrTimeoutHasOccurred();
}
```
A spurious wakeup could cause this to proceed to
`doStuffAssumingConditionIsTrueOrTimeoutHasOccurred()` even if the condition is
still false and time less than the timeout has elapsed. Instead, you should
write:

```
synchronized (this) {
 long now = System.currentTimeMillis();
 long deadline = now + timeout;
 while (!condition && now < deadline) {
 wait(deadline - now);
 now = System.currentTimeMillis();
 }
 doStuffAssumingConditionIsTrueOrTimeoutHasOccurred();
}
```
First, a warning: This type of waiting/sleeping is often done when the real
intent is to wait until some operation completes, and then proceed. If that’s
what you’re trying to do, please consider rewriting your code to use one of the
patterns above. Otherwise you are depending on system-specific timing that
*will* change when you run on different machines.

The incorrect code for this typically looks like:

```
synchronized (this) {
 // Give some time for the foos to bar
 wait(1000);
}
```
A spurious wakeup could cause this not to wait for a full 1000 ms. Instead, you
should use `Thread.sleep()`, which is not subject to spurious wakeups:

```
Thread.sleep(1000);
```
The incorrect code for this typically looks like:

```
synchronized (this) {
 // wait forever
 wait();
}
```
A spurious wakeup could cause this not to wait forever. You should wrap the call
to `wait()` in a `while (true)` loop:

```
synchronized (this) {
 // wait forever
 while (true) {
 wait();
 }
}
```
See Java Concurrency in Practice section 14.2.2, “Waking up too soon,”
the Javadoc for `Object.wait()`,
and the “Implementation Considerations” section in
the Javadoc for `Condition`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("WaitNotInLoop")` to the enclosing element.

#
 WakelockReleasedDangerously

 On Android versions < P, a wakelock acquired with a timeout may be released by the system before calling `release`, even after checking `isHeld()`. If so, it will throw a RuntimeException. Please wrap in a try/catch block.

 - Severity
- WARNING
- Tags
- FragileCode

@AutoFactory classes should not be @Inject-ed, inject the generated factory instead. Classes that are annotated with @AutoFactory are intended to be constructed by invoking the factory method on the generated factory. Typically this is because some of the necessary constructor arguments are not part of the binding graph. Generated @AutoFactory classes are automatically marked @Inject - prefer to inject that instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AutoFactoryAtInject")` to the enclosing element.

The Java class loading APIs can lead to remote code execution vulnerabilities if not used carefully. Interpreting potentially untrusted input as bytecode can give an attacker control of the application.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("BanClassLoader")` to the enclosing element.

The Java `Serializable` API is very powerful, and very dangerous. Any
consumption of a serialized object that cannot be explicitly trusted will likely
result in a critical remote code execution bug that will give an attacker
control of the application. (See
Effective Java 3rd Edition §85)

Consider using less powerful serialization methods, such as JSON or XML.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("BanSerializableRead")` to the enclosing element.

*Alternate names: TopLevelName*

While the Java programming language requires that `public` top-level classes are
declared in a source file matching their name (e.g.: `Foo.java` for ```
public
class Foo {}
```
, it is possible to declare a non-public top-level class in a file
with a different name (e.g.: `Bar.java` for `class Foo {}`).

The Google Java Style Guide §2.1 states, “The source file name consists of the case-sensitive name of the top-level class it contains, plus the .java extension.”

Since `@SuppressWarnings` cannot be applied to package declarations, this
warning can be suppressed by annotating any top-level class in the compilation
unit with `@SuppressWarnings("ClassName")`.

The comparison contract states that `sgn(compare(x, y)) == -sgn(compare(y, x))`.
(An immediate corollary is that `compare(x, x) == 0`.) This comparison
implementation either a) cannot return 0, b) cannot return a negative value but
may return a positive value, or c) cannot return a positive value but may return
a negative value.

The results of violating this contract can include `TreeSet.contains` never
returning true or `Collections.sort` failing with an IllegalArgumentException
arbitrarily.

In the long term, essentially all Comparators should be rewritten to use the Java 8 Comparator factory methods, but our automated migration tools will, of course, only work for correctly implemented Comparators.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ComparisonContractViolated")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("DeduplicateConstants") to the enclosing element.

@SuppressWarnings("DeduplicateConstants")

*Alternate names: dep-ann*

A declaration has the `@deprecated` Javadoc tag but no `@Deprecated` annotation.
Please add an `@Deprecated` annotation to this declaration in addition to the
`@deprecated` tag in the Javadoc.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DepAnn")` to the enclosing element.

*Alternate names: empty*

An if statement contains an empty statement as the then clause. A semicolon may have been inserted by accident.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("EmptyIf")` to the enclosing element.

`@AutoValue` classes are intended to be closed, with a single implementation
with known semantics. Implementing them by hand is extremely dangerous.

Here are some common cases where we have seen code that extends `@AutoValue`
classes and recommendations for what to do instead:

**Overriding the getters to return given values.** This should usually be
replaced by creating an instance of the `@AutoValue` class with those
values.

**Having the @AutoValue.Builder class extend the @AutoValue class so it
inherits the abstract getters.** This is wrong since the Builder doesn’t
satisfy the contract of its superclass and is better implemented either by
repeating the methods or by having both the

`@AutoValue` class and the
Builder implement a common interface with these getters.In other cases of extending the `@AutoValue` class and implementing its
abstract methods, it would be more correct to have the `@AutoValue` class
and the other class have a common supertype.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("ExtendsAutoValue")` to the enclosing element.

*Alternate names: InsecureCipherMode*

This checker looks for usages of standard cryptographic algorithms in configurations that are prone to vulnerabilities. There are currently three classes of problems that are covered by this checker:

Creating an instance of `javax.crypto.Cipher` using either the default
settings or the notoriously insecure ECB mode. In particular, Java’s default
`Cipher.getInstance(AES)` returns a cipher object that operates in ECB mode.
Dynamically constructed transformation strings are also flagged, as they may
conceal an instance of ECB mode. The problem with ECB mode is that
encrypting the same block of plaintext always yields the same block of
ciphertext. Hence, repetitions in the plaintext translate into repetitions
in the ciphertext, which can be readily used to conduct cryptanalysis. The
use of IES-based cipher algorithms also raises an error, as all currently
available implementations use ECB mode under the hood.

Using the Diffie-Hellmann protocol on prime fields. Most library implementations of Diffie-Hellman on prime fields have serious issues that can be exploited by an attacker. Any operation that may involve this protocol will be flagged by the checker. Implementations of the protocol based on elliptic curves (ECDH) are secure and should be used instead.

Using DSA for digital signatures. Some widely used crypto libraries accept invalid DSA signatures in specific configurations. The checker will flag all cryptographic operations that may involve DSA.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("InsecureCryptoUsage")` to the enclosing element.

`java.nio.file.Path` implements `Iterable<Path>`, and provides an iterator
over the name elements of the path. Declaring a parameter of type
`Iterable<Path>` is not recommended, since it allows clients to pass either an
`Iterable` of `Path`s, or a single `Path`. Using `Collection<Path>` prevents
clients from accidentally passing a single `Path`.

Example:

```
void printPaths(Iterable<Path> paths) {
 for (Path path : paths) System.err.println(path);
}
```
```
printPaths(Paths.get("/tmp/hello"));
tmp
hello
```
```
printPaths(ImmutableList.of(Paths.get("/tmp/hello")));
/tmp/hello
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("IterablePathParameter")` to the enclosing element.

Code that needs to be compatible with Java 8 cannot use types or members that are only present in newer class libraries

Suppress false positives by adding the suppression annotation @SuppressWarnings("Java8ApiChecker") to the enclosing element.

@SuppressWarnings("Java8ApiChecker")

A long literal can have a suffix of ‘L’ or ‘l’, but the former is less likely to be confused with a ‘1’ in most fonts.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("LongLiteralLowerCaseSuffix")` to the enclosing element.

*Alternate names: InjectScopeOrQualifierAnnotationRetention*

Qualifier and Scope annotations are used by dependency injection frameworks to adjust their behavior. Not having runtime retention on scoping or qualifier annotations will cause unexpected behavior in frameworks that use reflection:

```
class CreditCardProcessor { @Inject CreditCardProcessor(...) }
@Qualifier
@interface ForTests
@Provides
@ForTests
CreditCardProcessor providesTestProcessor() { return new TestCreditCardProcessor(...) }
...
@Inject
MyApp(CreditCardProcessor processor) {
 processor.issueCharge(...); // Issues a charge against a fake!
}
```
Since the Qualifier doesn’t have runtime retention, the Guice provider method doesn’t see the annotation, and will use the TestCreditCardProcessor for the normal CreditCardProcessor injection point.

NOTE: Even for dependency injection frameworks traditionally considered to be
compile-time dependent, the JSR-330 specification still requires runtime
retention for both `Qualifier` and `Scope`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("MissingRuntimeRetention")` to the enclosing element.

Alternate names: InjectMoreThanOneQualifier

An element can be qualified by at most one qualifier.

Suppress false positives by adding the suppression annotation @SuppressWarnings("MoreThanOneQualifier") to the enclosing element.

@SuppressWarnings("MoreThanOneQualifier")

Like many other languages, Java provides automatic memory management. In Java, this feature incurs an runtime cost, and can also lead to unpredictable execution pauses. In most cases, this is a reasonable tradeoff, but sometimes the loss of performance or predictability is unacceptable. Examples include pause-sensitive user interface handlers, high query rate server response handlers, or other soft-realtime applications.

In these situations, you can annotate a few carefully written methods with @NoAllocation. Methods with this annotation will avoid allocations in most cases, reducing pressure on the garbage collector. Note that allocations may still occur in methods with @NoAllocation if the compiler or runtime system inserts them.

To ease the use of exceptions, allocations are allowed if they occur within a throw statement. But if the throw statement contains a nested class with methods annotated with @NoAllocation, those methods will be disallowed from allocating.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("NoAllocation")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("RefersToDaggerCodegen") to the enclosing element.

@SuppressWarnings("RefersToDaggerCodegen")

Static and default interface methods are not natively supported on Android versions earlier than 7.0. Enable this check for compatibility with older devices. See Android Java 8 Documentation.

To declare default or static methods in interfaces, add a
`@SuppressWarnings("StaticOrDefaultInterfaceMethod")` annotation to the
enclosing element.

*Alternate names: static, static-access, StaticAccessedFromInstance*

To refer to a static member of another class, we typically *qualify* that member
name by prepending the name of the class it’s in, and a dot:
`TheClass.theMethod()`.

But the Java language also permits you to qualify this call using any
*expression* (typically, a variable) whose static type is the class that
contains the method: `instanceOfTheClass.theMethod()`.

Doing this creates the appearance of an ordinary polymorphic method call, but it behaves very differently. For example:

```
public class Main {
 static class TheClass {
 public static int theMethod() {
 return 1;
 }
 }
 static class TheSubclass extends TheClass {
 public static int theMethod() {
 return 2;
 }
 }
 public static void main(String[] args) {
 TheClass instanceOfTheClass = new TheSubclass();
 System.out.println(instanceOfTheClass.theMethod());
 }
}
```
`TheSubclass` appears to “override” `theMethod`, so we might expect this code to
print the number `2`.

The code, however, prints the number `1`. The runtime type of
`instanceOfTheClass`, `TheSubclass`, is ignored; only the static type of the
reference, as seen by `javac`, matters.

In fact, the instance that `instanceOfTheClass` points to is *entirely*
irrelevant. To prove this, set the variable to `null` and run again. The program
will *still* print `1`, not throw a `NullPointerException`!

Qualifying a static reference in this way creates an unnecessarily confusing
situation. To prevent it, only qualify static method calls using a class name,
never an expression (that is, `TheClass.theMethod()` or `theMethod()`, but not
`anInstanceOfTheClass.theMethod()`).

Suppress false positives by adding the suppression annotation `@SuppressWarnings("StaticQualifiedUsingExpression")` to the enclosing element.

Calling `System.exit` terminates the java process and returns a status code.
Since it is disruptive to shut down the process within library code,
`System.exit` should not be called outside of a main method.

Instead of calling `System.exit` consider throwing an unchecked exception to
signal failure.

For example, prefer this:

```
public static void main(String[] args) {
 try {
 doSomething(args[0]);
 } catch (MyUncheckedException e) {
 System.err.println(e.getMessage());
 System.exit(1);
 }
}
// In library code
public static void doSomething(String s) {
 try {
 doSomethingElse(s);
 } catch (MyCheckedException e) {
 throw new MyUncheckedException(e);
 }
}
```
to this:

```
public static void main(String[] args) {
 doSomething(args[0]);
}
// In library code
public static void doSomething(String s) {
 try {
 doSomethingElse(s)
 } catch (MyCheckedException e) {
 System.err.println(e.getMessage());
 System.exit(1);
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("SystemExitOutsideMain")` to the enclosing element.

The use of `@Test(expected = FooException.class)` is strongly discouraged, since
the test passes if *any* statement throws an exception of the expected type.

For example, if `add(0, "a")` throws an `UnsupportedOperationException` below,
the test will pass without even executing `remove(0)`, much less testing whether
it throws the right kind of exception. Such false negatives are particularly
likely when testing for common unchecked exceptions like `NullPointerException`.

```
@Test(expected = UnsupportedOperationException.class)
public void testRemoveFails() {
 AppendOnlyList list = new AppendOnlyList();
 list.add(0, "a");
 list.remove(0);
}
```
To avoid this issue, prefer JUnit’s `assertThrows()` API:

```
import static org.junit.Assert.assertThrows;
@Test
public void testRemoveFails() {
 AppendOnlyList list = new AppendOnlyList();
 list.add(0, "a");
 assertThrows(UnsupportedOperationException.class, () -> {
 list.remove(0);
 });
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("TestExceptionChecker")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("ThreadSafe") to the enclosing element.

@SuppressWarnings("ThreadSafe")

Java assert statements are not run unless explicitly enabled via runtime flags to the JVM invocation.

If asserts are not enabled, then a test using assert would continue to pass even
if a bug is introduced since these statements will not be executed. To avoid
this, use one of the assertion libraries that are always enabled, such as
JUnit’s `org.junit.Assert` or Google’s Truth library. These will also produce
richer contextual failure diagnostics to aid and accelerate debugging.

Don’t do this:

```
@Test
public void testArray() {
 String[] arr = getArray();
 assert arr != null;
 assert arr.length == 1;
 assert arr[0].equals("hello");
}
```
Do this instead:

```
import static com.google.common.truth.Truth.assertThat;
@Test
public void testArray() {
 String[] arr = getArray();
 assertThat(arr).isNotNull();
 assertThat(arr).hasLength(1);
 assertThat(arr[0]).isEqualTo("hello");
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("UseCorrectAssertInTests")` to the enclosing element.

Per the style guide, `TYPE_USE` annotations should appear
immediately before the type being annotated, and after any modifiers:

```
public <K, V> @Nullable V getOrNull(final Map<K, V> map, final @Nullable K key) {
 return map.get(key);
}
```
Non-`TYPE_USE` annotations should appear before modifiers, as they annotate the
entire element (method, variable, class):

```
@VisibleForTesting
public void reset() {
 // ...
}
```
Javadoc must appear before any annotations, or the compiler will fail to recognise it as Javadoc:

```
@Nullable
/** Might return a frobnicator. */
Frobnicator getFrobnicator();
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("AnnotationPosition")` to the enclosing element.

Java assertions do not necessarily execute at runtime; they may be enabled and disabled depending on which options are passed to the JVM invocation. An assert false statement may be intended to ensure that the program never proceeds beyond that statement. If the correct execution of the program depends on that being the case, consider throwing an exception instead, so that execution is halted regardless of runtime configuration.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AssertFalse")` to the enclosing element.

Mixing @Inject and @AssistedInject leads to confusing code and the documentation specifies not to do it. See https://google.github.io/guice/api-docs/latest/javadoc/com/google/inject/assistedinject/AssistedInject.html

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AssistedInjectAndInjectOnConstructors")` to the enclosing element.

Object arrays are inferior to collections in almost every way. Prefer `Set`,
`List`, or `Multiset` over an object array whenever possible.

A few of these issues are covered, in much greater detail, in Effective Java Item 28: Prefer lists to arrays.

Don’t use object arrays as method parameters:

```
public void createUsers(User[] users) { ... }
```
Use an `Iterable` instead:

```
public void createUsers(Iterable<User> users) { ... }
```
Don’t use object arrays as method return values:

```
public User[] loadUsers() { ... }
```
Use an `ImmutableList` (or `ImmutableSet`) instead:

```
public ImmutableList<User> loadUsers() { ... }
```
If you have a 2-dimensional array (e.g., `Foo[][]`), consider using an
`ImmutableTable<Integer, Integer, Foo>` instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("AvoidObjectArrays")` to the enclosing element.

#
 BinderIdentityRestoredDangerously

 A call to Binder.clearCallingIdentity() should be followed by Binder.restoreCallingIdentity() in a finally block. Otherwise the wrong Binder identity may be used by subsequent code.

 - Severity
- WARNING
- Tags
- FragileCode

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("BinderIdentityRestoredDangerously")` to the enclosing element.

Guice bindings are keyed by a pair of (optional Annotation, Type).

In most circumstances, one doesn’t need the annotation, as there’s really just one active implementation:

```
bind(CoffeeMaker.class).to(RealCoffeeMaker.class);
...
@Inject Office(CoffeeMaker coffeeMaker) {}
```
However, in other circumstances, you want to bind a simple value (an integer,
String, double, etc.). You should use a Qualifier annotation to allow you to get
the *right* Integer back:

```
bindConstant().annotatedWith(HttpPort.class).to(80);
...
@Inject MyWebServer(@HttpPort Integer httpPort) {}
```
NOTE: Make sure that your annotation has the `@Qualifier` meta-annotation on
it, otherwise injection systems can’t see them. Guice users can optionally use
`@BindingAnnotation`, but Guice also understands `@Qualifier`.

This works great, but if your integer binding *doesn’t* include a Qualifier, it
just means that you can ask Guice for “the Integer”, and it will give you a
value back:

```
bind(Integer.class).toInstance(80);
...
@Inject MyWebServer(Integer httpsPort) {}
```
To avoid confusion in these circumstances, please use a Qualifier annotation when binding simple value types.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("BindingToUnqualifiedCommonType")` to the enclosing element.

This check identifies instance methods in builder classes, and requires that
they either `return this;`, or are explicitly annotated with
`@CheckReturnValue`.

Instance methods in builders typically return `this`, to allow chaining.
Ignoring this result does not indicate a https://errorprone.info/bugpattern/CheckReturnValue bug. For
example, both of the following are fine:

```
Foo.Builder builder = Foo.builder();
builder.setBar("bar"); // return value is deliberately unused
return builder.build();
```
```
Foo.Builder builder =
 Foo.builder().setBar("bar").build();
```
Rarely, a builder method may return a new instance, which should not be ignored.
This check requires these methods to be annotated with `@CheckReturnValue`:

```
class Builder {
 @CheckReturnValue
 Builder setFoo(String foo) {
 return new Builder(foo); // returns a new builder instead of this!
 }
}
```
This check allows the https://errorprone.info/bugpattern/CheckReturnValue enforcement to assume the
return value of instance methods in builders can safely be ignored, unless the
method is explicitly annotated with `@CheckReturnValue`.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("BuilderReturnThis")` to the enclosing element.

#
 CanIgnoreReturnValueSuggester

 Methods that always return 'this' (or return an input parameter) should be annotated with @com.google.errorprone.annotations.CanIgnoreReturnValue

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("CanIgnoreReturnValueSuggester")` to the enclosing element.

Mockito cannot mock final classes. See https://github.com/mockito/mockito/wiki/FAQ for details.

Suppress false positives by adding the suppression annotation @SuppressWarnings("CannotMockFinalClass") to the enclosing element.

@SuppressWarnings("CannotMockFinalClass")

*Alternate names: MockitoBadFinalMethod, CannotMockFinalMethod*

Mockito cannot mock `final` or `static` methods, and cannot tell at runtime that
this is attempted and fail with an error (as mocking `final` classes does).

`when(mock.finalMethod())` will invoke the real implementation of `finalMethod`.
In some cases, this may wind up accidentally doing what’s intended:

```
when(converter.convert(a)).thenReturn(b);
```
`convert` is final, but under the hood, calls `doForward`, so we wind up mocking
that method instead.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("CannotMockMethod")` to the enclosing element.

#
 CatchingUnchecked

 This catch block catches `Exception`, but can only catch unchecked exceptions. Consider catching RuntimeException (or something more specific) instead so it is more apparent that no checked exceptions are being handled.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("CatchingUnchecked")` to the enclosing element.

Java allows methods to declare that they throw checked exceptions even when they don’t. This can lead to call sites being forced to explicitly handle or propagate exceptions which provably can never occur. It may also lead readers of the code to search for possibly throwing paths where none exist.

```
private static void validateRequest(Request request) throws IOException {
 checkArgument(request.hasFoo(), "foo must be specified");
}
Response handle(Request request) {
 try {
 validateRequest(request);
 } catch (IOException e) { // Required, but unreachable.
 return failedResponse();
 }
 // ...
}
```
Including unthrown exceptions can be reasonable where the method is overridable, as overriding methods will not be able to declare that they throw any exceptions not included.

Suppress false positives by adding the suppression annotation `@SuppressWarnings("CheckedExceptionNotThrown")` to the enclosing element.

Method bodies should generally not call
`java.util.regex.Pattern#compile(String)` with constant arguments. Instead,
define a constant to store that Pattern. This can avoid recompilation of the
regex every time the method is invoked.

That is, prefer this:

```
private static final Pattern REGEX_PATTERN = Pattern.compile("a+");
public static boolean doSomething(String input) {
 Matcher matcher = REGEX_PATTERN.matcher(input);
 if (matcher.matches()) {
 ...
 }
}
```
to this:

```
public static boolean doSomething(String input) {
 Matcher matcher = Pattern.compile("a+").matcher(input);
 if (matcher.matches()) {
 ...
 }
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ConstantPatternCompile")` to the enclosing element.

#
 DefaultLocale

 Implicit use of the JVM default locale, which can result in differing behaviour between JVM executions.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DefaultLocale")` to the enclosing element.

#
 DifferentNameButSame

 This type is referred to in different ways within this file, which may be confusing.

 - Severity
- WARNING
- Tags
- Style

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("DifferentNameButSame")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("EqualsBrokenForNull") to the enclosing element.

@SuppressWarnings("EqualsBrokenForNull")

Any additional statements after the statement that is expected to throw will never be executed in a passing test. This can lead to inappropriately passing tests where later incorrect assertions are skipped by the thrown exception. For instance, the final assertion in the following example will never be executed if the call throws as expected.

```
@Test
public void testRemoveFails() {
 AppendOnlyList list = new AppendOnlyList();
 list.add(0, "a");
 thrown.expect(UnsupportedOperationException.class);
 thrown.expectMessage("hello");
 list.remove(0); // throws
 assertThat(list).hasSize(1); // never executed
}
```
To avoid this issue, prefer JUnit’s `assertThrows()` API:

```
import static org.junit.Assert.assertThrows;
@Test
public void testRemoveFails() {
 AppendOnlyList list = new AppendOnlyList();
 list.add(0, "a");
 UnsupportedOperationException thrown = assertThrows(
 UnsupportedOperationException.class,
 () -> {
 list.remove(0);
 });
 assertThat(thrown).hasMessageThat().contains("hello");
 assertThat(list).hasSize(1);
}
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("ExpectedExceptionChecker")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("ExplicitArrayForVarargs") to the enclosing element.

@SuppressWarnings("ExplicitArrayForVarargs")

#
 FloggerLogWithCause

 Setting the caught exception as the cause of the log message may provide more context for anyone debugging errors.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("FloggerLogWithCause")` to the enclosing element.

Flogger uses printf-style format specifiers, such as %s and %d. Message format-style specifiers like {0} don’t work.

Suppress false positives by adding the suppression annotation @SuppressWarnings("FloggerMessageFormat") to the enclosing element.

@SuppressWarnings("FloggerMessageFormat")

#
 FloggerRedundantIsEnabled

 Logger level check is already implied in the log() call. An explicit atLEVEL().isEnabled() check is redundant.

## Suppression

Suppress false positives by adding the suppression annotation `@SuppressWarnings("FloggerRedundantIsEnabled")` to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("FloggerRequiredModifiers") to the enclosing element.

@SuppressWarnings("FloggerRequiredModifiers")

FloggerWithCause
Calling withCause(Throwable) with an inline allocated Throwable is discouraged. Consider using withStackTrace(StackSize) instead, and specifying a reduced stack size (e.g. SMALL, MEDIUM or LARGE) instead of FULL, to improve performance.

Severity

WARNING

Suppression

Suppress false positives by adding the suppression annotation @SuppressWarnings("FloggerWithCause") to the enclosing element.

Suppress false positives by adding the suppression annotation @SuppressWarnings("FloggerWithoutCause") to the enclosing element.

@SuppressWarnings("FloggerWithoutCause")

Passing lambdas to an overloaded method may be ambiguous if two overloads have parameters that are functional interfaces with equivalent methods.

Prefer to avoid ambiguous overloads, and consider renaming one of the methods.

For example `Function<String, Integer>` and `IntFunction<String>` are both
compatible with the lambda `x -> x.hashCode()`.

```
void f(Function<String, Integer> x) {}
void f(IntFunction<String> x) {}
```
```
error: reference to f is ambiguous
 f(x -> x.hashCode());
 ^
 both method f(Function<String,Integer>) in Test and method f(IntFunction<String>) in Test match
```
To avoid the ambiguity, callers will have to use an explicit cast:

```
f((IntFunction<String>) x -> x.hashCode());
f((Function<String, Integer>) x -> x.hashCode());
```
The situation is more complicated with expression-bodied lambdas. Consider:

```
void doIt(Function<String, String> f);
void doIt(Consumer<String> c);
```
JLS 15.12.2.1
says that lambdas whose body is a statement expression are compatible with
functional interfaces whose function type is void-returning *or* value
returning:

A lambda expression (§15.27) is potentially compatible with a functional interface type (§9.8) if all of the following are true:

The arity of the target type’s function type is the same as the arity of the lambda expression.

If the target type’s function type has a void return, then the lambda body is either a statement expression (§14.8) or a void-compatible block (§15.27.2).

If the target type’s function type has a (non-void) return type, then the lambda body is either an expression or a value-compatible block (§15.27.2).

So, if you have:

```
doIt(x -> System.gc());
```
it’s an implicitly typed statement-expression-bodied lambda that’s compatible with both overloads.

Any of the following disambiguate the overloads:

```
doIt((String x) -> x.toString()); // explicitly typed, calls f(Function)
doIt(x -> (x.toString())); // non-statement expression body, calls f(Function)
doIt(x -> {x.toString();}); // statement body, calls f(Consumer)
```
Suppress false positives by adding the suppression annotation `@SuppressWarnings("FunctionalInterfaceClash")` to the enclosing element.

Error Prone lets the user enable and disable specific checks as well as override their built-in severity levels (warning vs. error) by passing options to the Error Prone compiler invocation.

A valid Error Prone command-line option looks like:

```
-Xep:<checkName>[:severity]
```
`checkName` is required and is the canonical name of the check, e.g.
“ReferenceEquality”. `severity` is one of {“OFF”, “WARN”, “ERROR”}. Multiple
flags must be passed to enable or disable multiple checks. The last flag for a
specific check wins.

Examples of usage follow:

```
-Xep:ReferenceEquality [turns on ReferenceEquality check with the severity level from its BugPattern annotation]
-Xep:ReferenceEquality:OFF [turns off ReferenceEquality check]
-Xep:ReferenceEquality:WARN [turns on ReferenceEquality check as a warning]
-Xep:ReferenceEquality:ERROR [turns on ReferenceEquality check as an error]
-Xep:ReferenceEquality:OFF -Xep:ReferenceEquality [turns on ReferenceEquality check]
```
There are also a few blanket severity-changing flags:

`-XepAllErrorsAsWarnings``-XepAllSuggestionsAsWarnings``-XepAllDisabledChecksAsWarnings``-XepDisableAllChecks``-XepDisableAllWarnings``-XepDisableWarningsInGeneratedCode` : Disables warnings in classes
annotated with `@Generated`With any of the blanket flags, you can pass additional flags afterward to tweak the level of individual checks. E.g., this flag combination disables all checks except for ReferenceEquality:

```
-XepDisableAllChecks -Xep:ReferenceEquality:ERROR
```
Additionally, you can completely exclude certain paths from any Error Prone
checking via the `-XepExcludedPaths` flag. The flag takes as an argument a
regular expression that is matched against a source file’s path to determine
whether it should be excluded. So, to exclude files in any sub-directory of a
path containing `build/generated`, use the option:

```
-XepExcludedPaths:.*/build/generated/.*
```
If you pass a flag that refers to an unknown check name, by default Error Prone
will throw an error. You can allow the use of unknown check names by passing the
`-XepIgnoreUnknownCheckNames` flag.

We no longer support the old-style Error Prone disabling flags that used the
`-Xepdisable:<checkName>` syntax.

There are a couple of flags for configuration of patching in suggested fixes,
e.g. `-XepPatchChecks:VALUE` and `-XepPatchLocation:VALUE`. See the
patching docs for more info.

To configure checks, you can use custom flags to pass info directly to BugCheckers. A valid custom flag looks like this:

```
-XepOpt:[Namespace:]FlagName[=Value]
```
By convention, if a flag is only relevant to one check or a group of checks, to
prevent name collision, you should prefix your flag’s name with an optional
namespace and a colon, e.g. `-XepOpt:JUnit4TestNotRun:ExpandedHeuristic=true`.

If a flag is set with no value provided, that flag is set to `true`, e.g.
`-XepOpt:MakeAwesome` is equivalent to `-XepOpt:MakeAwesome=true`.

Some examples:

```
-XepOpt:FlagName=SomeValue (flags["FlagName"] = "SomeValue")
-XepOpt:BooleanFlag (flags["BooleanFlag"] = "true")
-XepOpt:ListFlag=1,2,3 (flags["ListFlag"] = "1,2,3")
-XepOpt:Namespace:SomeFlag=AValue (flags["Namespace:SomeFlag"] = "AValue")
```
These flags can be accessed in a BugChecker just by adding a one-argument
constructor that takes an `ErrorProneFlags` object, like so:

```
public class MyChecker extends BugChecker implements SomeTreeMatcher {
 private final boolean coolness;
 public MyChecker(ErrorProneFlags flags) {
 // The ErrorProneFlags get* methods return an Optional<*>, use
 // Optional.orElse(...) and related methods to get with default, etc.
 this.coolness = flags.getBoolean("ErrorProne:IsCool").orElse(true);
 }
 public Description matchSomething(...) {...}
}
```
All arguments that are passed as `-Xep*` after `-Xplugin:ErrorProne` can also be
passed through a configuration file using `@`.

This allows easier sharing of errorprone configurations between various build systems (cli, maven, gradle, bazel, etc) and works around platform-dependent line wrapping rules.

The configuration file allows `#` as comment, both full line and inline.

Example:

```
-XepAllErrorsAsWarnings
-XepAllSuggestionsAsWarnings
-XepDisableWarningsInGeneratedCode
# This library back-ports `java.time`, so we need the old APIs
-Xep:JavaUtilDate:OFF
# Bug #123
-Xep:JUnit4TestNotRun:OFF
# Indends and trailing spaces are OK:
 -Xep:MisusedDayOfYear:WARN
 -Xep:MisusedWeekYear:WARN
# Several flags in one single line are OK:
-Xep:EffectivelyPrivate:OFF -Xep:ReturnValueIgnored:OFF -Xep:EmptyBlockTag:OFF
# Empty lines are OK
-Xep:ReturnValueIgnored:WARN # Inline comments are OK
```
Using it in command line:

```
@~/project/errorprone.cfg
```
NOTE: The `~` in the example above will not work in Windows, or Gradle, or
Maven, as it is expanded by the shell. You will need to use the methods specific
to your platform and build system.

NOTE: It is supported to pass a mixture of flags and several arguments files. The final value of a flag will be the one set that last time, regardless of whether that was done directly or in an arguments file.

To pass Error Prone flags to Maven, use the `compilerArgs` parameter in the
plugin’s configuration. The flags must be appended to the `arg` entry containing
`-Xplugin:ErrorProne`. To enable warnings, the `showWarnings` parameter must
also be set:

```
<project>
 <build>
 <plugins>
 <plugin>
 <artifactId>maven-compiler-plugin</artifactId>
 <configuration>
 <showWarnings>true</showWarnings>
 <compilerArgs>
 <arg>-XDcompilePolicy=simple</arg>
 <arg>--should-stop=ifError=FLOW</arg>
 <arg>-Xplugin:ErrorProne -Xep:DeadException:WARN -Xep:GuardedBy:OFF</arg>
 </compilerArgs>
 </configuration>
 </plugin>
 </plugins>
 </build>
</project>
```
Be aware that when running on JDK 8 the flags cannot be wrapped across multiple
lines. JDK 9 and above do allow the flags to be separated by newlines. That is,
the second `<arg>` element above can also be formatted as follows on JDK 9+, but
*not* on JDK 8:

```
<arg>
 -Xplugin:ErrorProne \
 -Xep:DeadException:WARN \
 -Xep:GuardedBy:OFF
</arg>
```
NOTE: using multi-line `<arg>`s does not work on Windows when `<fork>` is
enabled, see https://github.com/google/error-prone/issues/4256

NOTE: using an argument file (with `@`) allows bypassing this line wrapping bug:

```
<compilerArgs>
 ...
 <arg>-Xplugin:ErrorProne @${project.basedir}/errorprone.cfg</arg>
</compilerArgs>
```
And move the flags to `errorprone.cfg`:

```
-Xep:DeadException:WARN
-Xep:GuardedBy:OFF
```
