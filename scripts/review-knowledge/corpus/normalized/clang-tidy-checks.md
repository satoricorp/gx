```{title} clang-tidy - Clang-Tidy Checks
```

# Clang-Tidy Checks

```{toctree}
:glob: true
:hidden: true

abseil/*
altera/*
android/*
boost/*
bugprone/*
cert/*
clang-analyzer/*
concurrency/*
cppcoreguidelines/*
darwin/*
fuchsia/*
google/*
linuxkernel/*
llvm/*
llvmlibc/*
misc/*
modernize/*
mpi/*
objc/*
openmp/*
performance/*
portability/*
readability/*
```

| Name | Offers fixes |
| --- | --- |
| {doc}`abseil-cleanup-ctad <abseil/cleanup-ctad>` | Yes |
| {doc}`abseil-duration-addition <abseil/duration-addition>` | Yes |
| {doc}`abseil-duration-comparison <abseil/duration-comparison>` | Yes |
| {doc}`abseil-duration-conversion-cast <abseil/duration-conversion-cast>` | Yes |
| {doc}`abseil-duration-division <abseil/duration-division>` | Yes |
| {doc}`abseil-duration-factory-float <abseil/duration-factory-float>` | Yes |
| {doc}`abseil-duration-factory-scale <abseil/duration-factory-scale>` | Yes |
| {doc}`abseil-duration-subtraction <abseil/duration-subtraction>` | Yes |
| {doc}`abseil-duration-unnecessary-conversion <abseil/duration-unnecessary-conversion>` | Yes |
| {doc}`abseil-faster-strsplit-delimiter <abseil/faster-strsplit-delimiter>` | Yes |
| {doc}`abseil-no-internal-dependencies <abseil/no-internal-dependencies>` | |
| {doc}`abseil-no-namespace <abseil/no-namespace>` | |
| {doc}`abseil-redundant-strcat-calls <abseil/redundant-strcat-calls>` | Yes |
| {doc}`abseil-str-cat-append <abseil/str-cat-append>` | Yes |
| {doc}`abseil-string-find-startswith <abseil/string-find-startswith>` | Yes |
| {doc}`abseil-string-find-str-contains <abseil/string-find-str-contains>` | Yes |
| {doc}`abseil-time-comparison <abseil/time-comparison>` | Yes |
| {doc}`abseil-time-subtraction <abseil/time-subtraction>` | Yes |
| {doc}`abseil-unchecked-statusor-access <abseil/unchecked-statusor-access>` | |
| {doc}`abseil-upgrade-duration-conversions <abseil/upgrade-duration-conversions>` | Yes |
| {doc}`altera-id-dependent-backward-branch <altera/id-dependent-backward-branch>` | |
| {doc}`altera-kernel-name-restriction <altera/kernel-name-restriction>` | |
| {doc}`altera-single-work-item-barrier <altera/single-work-item-barrier>` | |
| {doc}`altera-struct-pack-align <altera/struct-pack-align>` | Yes |
| {doc}`altera-unroll-loops <altera/unroll-loops>` | |
| {doc}`android-cloexec-accept <android/cloexec-accept>` | Yes |
| {doc}`android-cloexec-accept4 <android/cloexec-accept4>` | Yes |
| {doc}`android-cloexec-creat <android/cloexec-creat>` | Yes |
| {doc}`android-cloexec-dup <android/cloexec-dup>` | Yes |
| {doc}`android-cloexec-epoll-create <android/cloexec-epoll-create>` | Yes |
| {doc}`android-cloexec-epoll-create1 <android/cloexec-epoll-create1>` | Yes |
| {doc}`android-cloexec-fopen <android/cloexec-fopen>` | Yes |
| {doc}`android-cloexec-inotify-init <android/cloexec-inotify-init>` | Yes |
| {doc}`android-cloexec-inotify-init1 <android/cloexec-inotify-init1>` | Yes |
| {doc}`android-cloexec-memfd-create <android/cloexec-memfd-create>` | Yes |
| {doc}`android-cloexec-open <android/cloexec-open>` | Yes |
| {doc}`android-cloexec-pipe <android/cloexec-pipe>` | Yes |
| {doc}`android-cloexec-pipe2 <android/cloexec-pipe2>` | Yes |
| {doc}`android-cloexec-socket <android/cloexec-socket>` | Yes |
| {doc}`android-comparison-in-temp-failure-retry <android/comparison-in-temp-failure-retry>` | |
| {doc}`boost-use-ranges <boost/use-ranges>` | Yes |
| {doc}`boost-use-to-string <boost/use-to-string>` | Yes |
| {doc}`bugprone-argument-comment <bugprone/argument-comment>` | Yes |
| {doc}`bugprone-assert-side-effect <bugprone/assert-side-effect>` | |
| {doc}`bugprone-assignment-in-if-condition <bugprone/assignment-in-if-condition>` | |
| {doc}`bugprone-assignment-in-selection-statement <bugprone/assignment-in-selection-statement>` | |
| {doc}`bugprone-bad-signal-to-kill-thread <bugprone/bad-signal-to-kill-thread>` | |
| {doc}`bugprone-bitwise-pointer-cast <bugprone/bitwise-pointer-cast>` | |
| {doc}`bugprone-bool-pointer-implicit-conversion <bugprone/bool-pointer-implicit-conversion>` | Yes |
| {doc}`bugprone-branch-clone <bugprone/branch-clone>` | |
| {doc}`bugprone-capturing-this-in-member-variable <bugprone/capturing-this-in-member-variable>` | |
| {doc}`bugprone-casting-through-void <bugprone/casting-through-void>` | |
| {doc}`bugprone-chained-comparison <bugprone/chained-comparison>` | |
| {doc}`bugprone-command-processor <bugprone/command-processor>` | |
| {doc}`bugprone-compare-pointer-to-member-virtual-function <bugprone/compare-pointer-to-member-virtual-function>` | |
| {doc}`bugprone-copy-constructor-init <bugprone/copy-constructor-init>` | Yes |
| {doc}`bugprone-copy-constructor-mutates-argument <bugprone/copy-constructor-mutates-argument>` | |
| {doc}`bugprone-crtp-constructor-accessibility <bugprone/crtp-constructor-accessibility>` | Yes |
| {doc}`bugprone-dangling-handle <bugprone/dangling-handle>` | |
| {doc}`bugprone-default-operator-new-on-overaligned-type <bugprone/default-operator-new-on-overaligned-type>` | |
| {doc}`bugprone-derived-method-shadowing-base-method <bugprone/derived-method-shadowing-base-method>` | |
| {doc}`bugprone-dynamic-static-initializers <bugprone/dynamic-static-initializers>` | |
| {doc}`bugprone-easily-swappable-parameters <bugprone/easily-swappable-parameters>` | |
| {doc}`bugprone-empty-catch <bugprone/empty-catch>` | |
| {doc}`bugprone-exception-copy-constructor-throws <bugprone/exception-copy-constructor-throws>` | |
| {doc}`bugprone-exception-escape <bugprone/exception-escape>` | |
| {doc}`bugprone-float-loop-counter <bugprone/float-loop-counter>` | |
| {doc}`bugprone-fold-init-type <bugprone/fold-init-type>` | |
| {doc}`bugprone-forward-declaration-namespace <bugprone/forward-declaration-namespace>` | |
| {doc}`bugprone-forwarding-reference-overload <bugprone/forwarding-reference-overload>` | |
| {doc}`bugprone-implicit-widening-of-multiplication-result <bugprone/implicit-widening-of-multiplication-result>` | Yes |
| {doc}`bugprone-inaccurate-erase <bugprone/inaccurate-erase>` | Yes |
| {doc}`bugprone-inc-dec-in-conditions <bugprone/inc-dec-in-conditions>` | |
| {doc}`bugprone-incorrect-enable-if <bugprone/incorrect-enable-if>` | Yes |
| {doc}`bugprone-incorrect-enable-shared-from-this <bugprone/incorrect-enable-shared-from-this>` | Yes |
| {doc}`bugprone-incorrect-roundings <bugprone/incorrect-roundings>` | |
| {doc}`bugprone-infinite-loop <bugprone/infinite-loop>` | |
| {doc}`bugprone-integer-division <bugprone/integer-division>` | |
| {doc}`bugprone-invalid-enum-default-initialization <bugprone/invalid-enum-default-initialization>` | |
| {doc}`bugprone-lambda-function-name <bugprone/lambda-function-name>` | |
| {doc}`bugprone-macro-parentheses <bugprone/macro-parentheses>` | Yes |
| {doc}`bugprone-macro-repeated-side-effects <bugprone/macro-repeated-side-effects>` | |
| {doc}`bugprone-misleading-setter-of-reference <bugprone/misleading-setter-of-reference>` | |
| {doc}`bugprone-misplaced-operator-in-strlen-in-alloc <bugprone/misplaced-operator-in-strlen-in-alloc>` | Yes |
| {doc}`bugprone-misplaced-pointer-arithmetic-in-alloc <bugprone/misplaced-pointer-arithmetic-in-alloc>` | Yes |
| {doc}`bugprone-misplaced-widening-cast <bugprone/misplaced-widening-cast>` | |
| {doc}`bugprone-missing-end-comparison <bugprone/missing-end-comparison>` | Yes |
| {doc}`bugprone-move-forwarding-reference <bugprone/move-forwarding-reference>` | Yes |
| {doc}`bugprone-multi-level-implicit-pointer-conversion <bugprone/multi-level-implicit-pointer-conversion>` | |
| {doc}`bugprone-multiple-new-in-one-expression <bugprone/multiple-new-in-one-expression>` | |
| {doc}`bugprone-multiple-statement-macro <bugprone/multiple-statement-macro>` | |
| {doc}`bugprone-narrowing-conversions <bugprone/narrowing-conversions>` | |
| {doc}`bugprone-no-escape <bugprone/no-escape>` | |
| {doc}`bugprone-non-zero-enum-to-bool-conversion <bugprone/non-zero-enum-to-bool-conversion>` | |
| {doc}`bugprone-nondeterministic-pointer-iteration-order <bugprone/nondeterministic-pointer-iteration-order>` | |
| {doc}`bugprone-not-null-terminated-result <bugprone/not-null-terminated-result>` | Yes |
| {doc}`bugprone-optional-value-conversion <bugprone/optional-value-conversion>` | Yes |
| {doc}`bugprone-parent-virtual-call <bugprone/parent-virtual-call>` | Yes |
| {doc}`bugprone-pointer-arithmetic-on-polymorphic-object <bugprone/pointer-arithmetic-on-polymorphic-object>` | |
| {doc}`bugprone-posix-return <bugprone/posix-return>` | Yes |
| {doc}`bugprone-random-generator-seed <bugprone/random-generator-seed>` | |
| {doc}`bugprone-raw-memory-call-on-non-trivial-type <bugprone/raw-memory-call-on-non-trivial-type>` | |
| {doc}`bugprone-redundant-branch-condition <bugprone/redundant-branch-condition>` | Yes |
| {doc}`bugprone-reserved-identifier <bugprone/reserved-identifier>` | Yes |
| {doc}`bugprone-return-const-ref-from-parameter <bugprone/return-const-ref-from-parameter>` | |
| {doc}`bugprone-shared-ptr-array-mismatch <bugprone/shared-ptr-array-mismatch>` | Yes |
| {doc}`bugprone-signal-handler <bugprone/signal-handler>` | |
| {doc}`bugprone-signed-bitwise <bugprone/signed-bitwise>` | |
| {doc}`bugprone-signed-char-misuse <bugprone/signed-char-misuse>` | |
| {doc}`bugprone-sizeof-container <bugprone/sizeof-container>` | |
| {doc}`bugprone-sizeof-expression <bugprone/sizeof-expression>` | |
| {doc}`bugprone-spuriously-wake-up-functions <bugprone/spuriously-wake-up-functions>` | |
| {doc}`bugprone-standalone-empty <bugprone/standalone-empty>` | Yes |
| {doc}`bugprone-std-exception-baseclass <bugprone/std-exception-baseclass>` | |
| {doc}`bugprone-std-namespace-modification <bugprone/std-namespace-modification>` | |
| {doc}`bugprone-string-constructor <bugprone/string-constructor>` | Yes |
| {doc}`bugprone-string-integer-assignment <bugprone/string-integer-assignment>` | Yes |
| {doc}`bugprone-string-literal-with-embedded-nul <bugprone/string-literal-with-embedded-nul>` | |
| {doc}`bugprone-stringview-nullptr <bugprone/stringview-nullptr>` | Yes |
| {doc}`bugprone-suspicious-enum-usage <bugprone/suspicious-enum-usage>` | |
| {doc}`bugprone-suspicious-include <bugprone/suspicious-include>` | |
| {doc}`bugprone-suspicious-memory-comparison <bugprone/suspicious-memory-comparison>` | |
| {doc}`bugprone-suspicious-memset-usage <bugprone/suspicious-memset-usage>` | Yes |
| {doc}`bugprone-suspicious-missing-comma <bugprone/suspicious-missing-comma>` | |
| {doc}`bugprone-suspicious-realloc-usage <bugprone/suspicious-realloc-usage>` | |
| {doc}`bugprone-suspicious-semicolon <bugprone/suspicious-semicolon>` | Yes |
| {doc}`bugprone-suspicious-string-compare <bugprone/suspicious-string-compare>` | Yes |
| {doc}`bugprone-suspicious-stringview-data-usage <bugprone/suspicious-stringview-data-usage>` | |
| {doc}`bugprone-swapped-arguments <bugprone/swapped-arguments>` | Yes |
| {doc}`bugprone-switch-missing-default-case <bugprone/switch-missing-default-case>` | |
| {doc}`bugprone-tagged-union-member-count <bugprone/tagged-union-member-count>` | |
| {doc}`bugprone-terminating-continue <bugprone/terminating-continue>` | Yes |
| {doc}`bugprone-throw-keyword-missing <bugprone/throw-keyword-missing>` | |
| {doc}`bugprone-throwing-static-initialization <bugprone/throwing-static-initialization>` | |
| {doc}`bugprone-too-small-loop-variable <bugprone/too-small-loop-variable>` | |
| {doc}`bugprone-unchecked-optional-access <bugprone/unchecked-optional-access>` | |
| {doc}`bugprone-unchecked-string-to-number-conversion <bugprone/unchecked-string-to-number-conversion>` | |
| {doc}`bugprone-undefined-memory-manipulation <bugprone/undefined-memory-manipulation>` | |
| {doc}`bugprone-undelegated-constructor <bugprone/undelegated-constructor>` | |
| {doc}`bugprone-unhandled-code-paths <bugprone/unhandled-code-paths>` | |
| {doc}`bugprone-unhandled-exception-at-new <bugprone/unhandled-exception-at-new>` | |
| {doc}`bugprone-unhandled-self-assignment <bugprone/unhandled-self-assignment>` | |
| {doc}`bugprone-unintended-char-ostream-output <bugprone/unintended-char-ostream-output>` | Yes |
| {doc}`bugprone-unique-ptr-array-mismatch <bugprone/unique-ptr-array-mismatch>` | Yes |
| {doc}`bugprone-unsafe-functions <bugprone/unsafe-functions>` | |
| {doc}`bugprone-unsafe-to-allow-exceptions <bugprone/unsafe-to-allow-exceptions>` | |
| {doc}`bugprone-unused-local-non-trivial-variable <bugprone/unused-local-non-trivial-variable>` | |
| {doc}`bugprone-unused-raii <bugprone/unused-raii>` | Yes |
| {doc}`bugprone-unused-return-value <bugprone/unused-return-value>` | |
| {doc}`bugprone-use-after-move <bugprone/use-after-move>` | |
| {doc}`bugprone-virtual-near-miss <bugprone/virtual-near-miss>` | Yes |
| {doc}`concurrency-mt-unsafe <concurrency/mt-unsafe>` | |
| {doc}`concurrency-thread-canceltype-asynchronous <concurrency/thread-canceltype-asynchronous>` | |
| {doc}`cppcoreguidelines-avoid-capturing-lambda-coroutines <cppcoreguidelines/avoid-capturing-lambda-coroutines>` | |
| {doc}`cppcoreguidelines-avoid-const-or-ref-data-members <cppcoreguidelines/avoid-const-or-ref-data-members>` | |
| {doc}`cppcoreguidelines-avoid-do-while <cppcoreguidelines/avoid-do-while>` | |
| {doc}`cppcoreguidelines-avoid-goto <cppcoreguidelines/avoid-goto>` | |
| {doc}`cppcoreguidelines-avoid-non-const-global-variables <cppcoreguidelines/avoid-non-const-global-variables>` | |
| {doc}`cppcoreguidelines-avoid-reference-coroutine-parameters <cppcoreguidelines/avoid-reference-coroutine-parameters>` | |
| {doc}`cppcoreguidelines-init-variables <cppcoreguidelines/init-variables>` | Yes |
| {doc}`cppcoreguidelines-interfaces-global-init <cppcoreguidelines/interfaces-global-init>` | |
| {doc}`cppcoreguidelines-macro-usage <cppcoreguidelines/macro-usage>` | |
| {doc}`cppcoreguidelines-misleading-capture-default-by-value <cppcoreguidelines/misleading-capture-default-by-value>` | Yes |
| {doc}`cppcoreguidelines-missing-std-forward <cppcoreguidelines/missing-std-forward>` | |
| {doc}`cppcoreguidelines-no-malloc <cppcoreguidelines/no-malloc>` | |
| {doc}`cppcoreguidelines-no-suspend-with-lock <cppcoreguidelines/no-suspend-with-lock>` | |
| {doc}`cppcoreguidelines-owning-memory <cppcoreguidelines/owning-memory>` | |
| {doc}`cppcoreguidelines-prefer-member-initializer <cppcoreguidelines/prefer-member-initializer>` | Yes |
| {doc}`cppcoreguidelines-pro-bounds-array-to-pointer-decay <cppcoreguidelines/pro-bounds-array-to-pointer-decay>` | |
| {doc}`cppcoreguidelines-pro-bounds-avoid-unchecked-container-access <cppcoreguidelines/pro-bounds-avoid-unchecked-container-access>` | Yes |
| {doc}`cppcoreguidelines-pro-bounds-constant-array-index <cppcoreguidelines/pro-bounds-constant-array-index>` | Yes |
| {doc}`cppcoreguidelines-pro-bounds-pointer-arithmetic <cppcoreguidelines/pro-bounds-pointer-arithmetic>` | |
| {doc}`cppcoreguidelines-pro-type-const-cast <cppcoreguidelines/pro-type-const-cast>` | |
| {doc}`cppcoreguidelines-pro-type-cstyle-cast <cppcoreguidelines/pro-type-cstyle-cast>` | Yes |
| {doc}`cppcoreguidelines-pro-type-member-init <cppcoreguidelines/pro-type-member-init>` | Yes |
| {doc}`cppcoreguidelines-pro-type-reinterpret-cast <cppcoreguidelines/pro-type-reinterpret-cast>` | |
| {doc}`cppcoreguidelines-pro-type-static-cast-downcast <cppcoreguidelines/pro-type-static-cast-downcast>` | Yes |
| {doc}`cppcoreguidelines-pro-type-union-access <cppcoreguidelines/pro-type-union-access>` | |
| {doc}`cppcoreguidelines-pro-type-vararg <cppcoreguidelines/pro-type-vararg>` | |
| {doc}`cppcoreguidelines-rvalue-reference-param-not-moved <cppcoreguidelines/rvalue-reference-param-not-moved>` | |
| {doc}`cppcoreguidelines-slicing <cppcoreguidelines/slicing>` | |
| {doc}`cppcoreguidelines-special-member-functions <cppcoreguidelines/special-member-functions>` | |
| {doc}`cppcoreguidelines-use-enum-class <cppcoreguidelines/use-enum-class>` | |
| {doc}`cppcoreguidelines-virtual-class-destructor <cppcoreguidelines/virtual-class-destructor>` | Yes |
| {doc}`darwin-avoid-spinlock <darwin/avoid-spinlock>` | |
| {doc}`darwin-dispatch-once-nonstatic <darwin/dispatch-once-nonstatic>` | Yes |
| {doc}`fuchsia-default-arguments-calls <fuchsia/default-arguments-calls>` | |
| {doc}`fuchsia-default-arguments-declarations <fuchsia/default-arguments-declarations>` | Yes |
| {doc}`fuchsia-overloaded-operator <fuchsia/overloaded-operator>` | |
| {doc}`fuchsia-statically-constructed-objects <fuchsia/statically-constructed-objects>` | |
| {doc}`fuchsia-temporary-objects <fuchsia/temporary-objects>` | |
| {doc}`fuchsia-trailing-return <fuchsia/trailing-return>` | |
| {doc}`fuchsia-virtual-inheritance <fuchsia/virtual-inheritance>` | |
| {doc}`google-build-explicit-make-pair <google/build-explicit-make-pair>` | |
| {doc}`google-build-using-namespace <google/build-using-namespace>` | |
| {doc}`google-default-arguments <google/default-arguments>` | |
| {doc}`google-global-names-in-headers <google/global-names-in-headers>` | |
| {doc}`google-objc-avoid-nsobject-new <google/objc-avoid-nsobject-new>` | |
| {doc}`google-objc-avoid-throwing-exception <google/objc-avoid-throwing-exception>` | |
| {doc}`google-objc-function-naming <google/objc-function-naming>` | |
| {doc}`google-objc-global-variable-declaration <google/objc-global-variable-declaration>` | |
| {doc}`google-readability-avoid-underscore-in-googletest-name <google/readability-avoid-underscore-in-googletest-name>` | |
| {doc}`google-readability-todo <google/readability-todo>` | |
| {doc}`google-runtime-float <google/runtime-float>` | |
| {doc}`google-runtime-int <google/runtime-int>` | |
| {doc}`google-runtime-operator <google/runtime-operator>` | |
| {doc}`google-upgrade-googletest-case <google/upgrade-googletest-case>` | Yes |
| {doc}`linuxkernel-must-check-errs <linuxkernel/must-check-errs>` | |
| {doc}`llvm-formatv-string <llvm/formatv-string>` | |
| {doc}`llvm-header-guard <llvm/header-guard>` | |
| {doc}`llvm-include-order <llvm/include-order>` | Yes |
| {doc}`llvm-namespace-comment <llvm/namespace-comment>` | |
| {doc}`llvm-prefer-isa-or-dyn-cast-in-conditionals <llvm/prefer-isa-or-dyn-cast-in-conditionals>` | Yes |
| {doc}`llvm-prefer-register-over-unsigned <llvm/prefer-register-over-unsigned>` | Yes |
| {doc}`llvm-prefer-static-over-anonymous-namespace <llvm/prefer-static-over-anonymous-namespace>` | |
| {doc}`llvm-redundant-casting <llvm/redundant-casting>` | Yes |
| {doc}`llvm-twine-local <llvm/twine-local>` | Yes |
| {doc}`llvm-type-switch-case-types <llvm/type-switch-case-types>` | Yes |
| {doc}`llvm-use-new-mlir-op-builder <llvm/use-new-mlir-op-builder>` | Yes |
| {doc}`llvm-use-ranges <llvm/use-ranges>` | Yes |
| {doc}`llvm-use-vector-utils <llvm/use-vector-utils>` | Yes |
| {doc}`llvmlibc-callee-namespace <llvmlibc/callee-namespace>` | |
| {doc}`llvmlibc-implementation-in-namespace <llvmlibc/implementation-in-namespace>` | |
| {doc}`llvmlibc-inline-function-decl <llvmlibc/inline-function-decl>` | Yes |
| {doc}`llvmlibc-restrict-system-libc-headers <llvmlibc/restrict-system-libc-headers>` | Yes |
| {doc}`misc-anonymous-namespace-in-header <misc/anonymous-namespace-in-header>` | |
| {doc}`misc-confusable-identifiers <misc/confusable-identifiers>` | |
| {doc}`misc-const-correctness <misc/const-correctness>` | Yes |
| {doc}`misc-coroutine-hostile-raii <misc/coroutine-hostile-raii>` | |
| {doc}`misc-definitions-in-headers <misc/definitions-in-headers>` | Yes |
| {doc}`misc-explicit-constructor <misc/explicit-constructor>` | Yes |
| {doc}`misc-header-include-cycle <misc/header-include-cycle>` | |
| {doc}`misc-include-cleaner <misc/include-cleaner>` | Yes |
| {doc}`misc-misleading-bidirectional <misc/misleading-bidirectional>` | |
| {doc}`misc-misleading-identifier <misc/misleading-identifier>` | |
| {doc}`misc-misplaced-const <misc/misplaced-const>` | |
| {doc}`misc-multiple-inheritance <misc/multiple-inheritance>` | |
| {doc}`misc-new-delete-overloads <misc/new-delete-overloads>` | |
| {doc}`misc-no-recursion <misc/no-recursion>` | |
| {doc}`misc-non-copyable-objects <misc/non-copyable-objects>` | |
| {doc}`misc-non-private-member-variables-in-classes <misc/non-private-member-variables-in-classes>` | |
| {doc}`misc-override-with-different-visibility <misc/override-with-different-visibility>` | |
| {doc}`misc-predictable-rand <misc/predictable-rand>` | |
| {doc}`misc-redundant-expression <misc/redundant-expression>` | Yes |
| {doc}`misc-static-assert <misc/static-assert>` | Yes |
| {doc}`misc-static-initialization-cycle <misc/static-initialization-cycle>` | |
| {doc}`misc-throw-by-value-catch-by-reference <misc/throw-by-value-catch-by-reference>` | |
| {doc}`misc-unconventional-assign-operator <misc/unconventional-assign-operator>` | |
| {doc}`misc-uniqueptr-reset-release <misc/uniqueptr-reset-release>` | Yes |
| {doc}`misc-unused-alias-decls <misc/unused-alias-decls>` | Yes |
| {doc}`misc-unused-parameters <misc/unused-parameters>` | Yes |
| {doc}`misc-unused-using-decls <misc/unused-using-decls>` | Yes |
| {doc}`misc-use-anonymous-namespace <misc/use-anonymous-namespace>` | |
| {doc}`misc-use-internal-linkage <misc/use-internal-linkage>` | Yes |
| {doc}`modernize-avoid-bind <modernize/avoid-bind>` | Yes |
| {doc}`modernize-avoid-c-arrays <modernize/avoid-c-arrays>` | |
| {doc}`modernize-avoid-c-style-cast <modernize/avoid-c-style-cast>` | Yes |
| {doc}`modernize-avoid-setjmp-longjmp <modernize/avoid-setjmp-longjmp>` | |
| {doc}`modernize-avoid-variadic-functions <modernize/avoid-variadic-functions>` | |
| {doc}`modernize-concat-nested-namespaces <modernize/concat-nested-namespaces>` | Yes |
| {doc}`modernize-deprecated-headers <modernize/deprecated-headers>` | Yes |
| {doc}`modernize-deprecated-ios-base-aliases <modernize/deprecated-ios-base-aliases>` | Yes |
| {doc}`modernize-loop-convert <modernize/loop-convert>` | Yes |
| {doc}`modernize-macro-to-enum <modernize/macro-to-enum>` | Yes |
| {doc}`modernize-make-shared <modernize/make-shared>` | Yes |
| {doc}`modernize-make-unique <modernize/make-unique>` | Yes |
| {doc}`modernize-min-max-use-initializer-list <modernize/min-max-use-initializer-list>` | Yes |
| {doc}`modernize-pass-by-value <modernize/pass-by-value>` | Yes |
| {doc}`modernize-raw-string-literal <modernize/raw-string-literal>` | Yes |
| {doc}`modernize-redundant-void-arg <modernize/redundant-void-arg>` | Yes |
| {doc}`modernize-replace-auto-ptr <modernize/replace-auto-ptr>` | Yes |
| {doc}`modernize-replace-disallow-copy-and-assign-macro <modernize/replace-disallow-copy-and-assign-macro>` | Yes |
| {doc}`modernize-replace-random-shuffle <modernize/replace-random-shuffle>` | Yes |
| {doc}`modernize-return-braced-init-list <modernize/return-braced-init-list>` | Yes |
| {doc}`modernize-shrink-to-fit <modernize/shrink-to-fit>` | Yes |
| {doc}`modernize-type-traits <modernize/type-traits>` | Yes |
| {doc}`modernize-unary-static-assert <modernize/unary-static-assert>` | Yes |
| {doc}`modernize-use-auto <modernize/use-auto>` | Yes |
| {doc}`modernize-use-bool-literals <modernize/use-bool-literals>` | Yes |
| {doc}`modernize-use-constraints <modernize/use-constraints>` | Yes |
| {doc}`modernize-use-default-member-init <modernize/use-default-member-init>` | Yes |
| {doc}`modernize-use-designated-initializers <modernize/use-designated-initializers>` | Yes |
| {doc}`modernize-use-emplace <modernize/use-emplace>` | Yes |
| {doc}`modernize-use-equals-default <modernize/use-equals-default>` | Yes |
| {doc}`modernize-use-equals-delete <modernize/use-equals-delete>` | Yes |
| {doc}`modernize-use-integer-sign-comparison <modernize/use-integer-sign-comparison>` | Yes |
| {doc}`modernize-use-nodiscard <modernize/use-nodiscard>` | Yes |
| {doc}`modernize-use-noexcept <modernize/use-noexcept>` | Yes |
| {doc}`modernize-use-nullptr <modernize/use-nullptr>` | Yes |
| {doc}`modernize-use-override <modernize/use-override>` | Yes |
| {doc}`modernize-use-ranges <modernize/use-ranges>` | Yes |
| {doc}`modernize-use-scoped-lock <modernize/use-scoped-lock>` | Yes |
| {doc}`modernize-use-starts-ends-with <modernize/use-starts-ends-with>` | Yes |
| {doc}`modernize-use-std-bit <modernize/use-std-bit>` | Yes |
| {doc}`modernize-use-std-format <modernize/use-std-format>` | Yes |
| {doc}`modernize-use-std-numbers <modernize/use-std-numbers>` | Yes |
| {doc}`modernize-use-std-print <modernize/use-std-print>` | Yes |
| {doc}`modernize-use-string-view <modernize/use-string-view>` | Yes |
| {doc}`modernize-use-structured-binding <modernize/use-structured-binding>` | Yes |
| {doc}`modernize-use-trailing-return-type <modernize/use-trailing-return-type>` | Yes |
| {doc}`modernize-use-transparent-functors <modernize/use-transparent-functors>` | Yes |
| {doc}`modernize-use-uncaught-exceptions <modernize/use-uncaught-exceptions>` | Yes |
| {doc}`modernize-use-using <modernize/use-using>` | Yes |
| {doc}`mpi-buffer-deref <mpi/buffer-deref>` | Yes |
| {doc}`mpi-type-mismatch <mpi/type-mismatch>` | Yes |
| {doc}`objc-assert-equals <objc/assert-equals>` | Yes |
| {doc}`objc-avoid-nserror-init <objc/avoid-nserror-init>` | |
| {doc}`objc-dealloc-in-category <objc/dealloc-in-category>` | |
| {doc}`objc-forbidden-subclassing <objc/forbidden-subclassing>` | |
| {doc}`objc-missing-hash <objc/missing-hash>` | |
| {doc}`objc-nsdate-formatter <objc/nsdate-formatter>` | |
| {doc}`objc-nsinvocation-argument-lifetime <objc/nsinvocation-argument-lifetime>` | Yes |
| {doc}`objc-property-declaration <objc/property-declaration>` | Yes |
| {doc}`objc-super-self <objc/super-self>` | Yes |
| {doc}`openmp-exception-escape <openmp/exception-escape>` | |
| {doc}`openmp-use-default-none <openmp/use-default-none>` | |
| {doc}`performance-avoid-endl <performance/avoid-endl>` | Yes |
| {doc}`performance-enum-size <performance/enum-size>` | |
| {doc}`performance-expensive-value-or <performance/expensive-value-or>` | Yes |
| {doc}`performance-for-range-copy <performance/for-range-copy>` | Yes |
| {doc}`performance-implicit-conversion-in-loop <performance/implicit-conversion-in-loop>` | |
| {doc}`performance-inefficient-algorithm <performance/inefficient-algorithm>` | Yes |
| {doc}`performance-inefficient-string-concatenation <performance/inefficient-string-concatenation>` | |
| {doc}`performance-inefficient-vector-operation <performance/inefficient-vector-operation>` | Yes |
| {doc}`performance-move-const-arg <performance/move-const-arg>` | Yes |
| {doc}`performance-move-constructor-init <performance/move-constructor-init>` | |
| {doc}`performance-no-automatic-move <performance/no-automatic-move>` | |
| {doc}`performance-no-int-to-ptr <performance/no-int-to-ptr>` | |
| {doc}`performance-noexcept-destructor <performance/noexcept-destructor>` | Yes |
| {doc}`performance-noexcept-move-constructor <performance/noexcept-move-constructor>` | Yes |
| {doc}`performance-noexcept-swap <performance/noexcept-swap>` | Yes |
| {doc}`performance-prefer-single-char-overloads <performance/prefer-single-char-overloads>` | Yes |
| {doc}`performance-string-view-conversions <performance/string-view-conversions>` | Yes |
| {doc}`performance-trivially-destructible <performance/trivially-destructible>` | Yes |
| {doc}`performance-type-promotion-in-math-fn <performance/type-promotion-in-math-fn>` | Yes |
| {doc}`performance-unnecessary-copy-initialization <performance/unnecessary-copy-initialization>` | Yes |
| {doc}`performance-unnecessary-value-param <performance/unnecessary-value-param>` | Yes |
| {doc}`performance-use-std-move <performance/use-std-move>` | Yes |
| {doc}`portability-avoid-pragma-once <portability/avoid-pragma-once>` | |
| {doc}`portability-no-assembler <portability/no-assembler>` | |
| {doc}`portability-restrict-system-includes <portability/restrict-system-includes>` | Yes |
| {doc}`portability-simd-intrinsics <portability/simd-intrinsics>` | |
| {doc}`portability-std-allocator-const <portability/std-allocator-const>` | |
| {doc}`portability-template-virtual-member-function <portability/template-virtual-member-function>` | |
| {doc}`readability-ambiguous-smartptr-reset-call <readability/ambiguous-smartptr-reset-call>` | Yes |
| {doc}`readability-avoid-const-params-in-decls <readability/avoid-const-params-in-decls>` | Yes |
| {doc}`readability-avoid-nested-conditional-operator <readability/avoid-nested-conditional-operator>` | |
| {doc}`readability-avoid-return-with-void-value <readability/avoid-return-with-void-value>` | Yes |
| {doc}`readability-avoid-unconditional-preprocessor-if <readability/avoid-unconditional-preprocessor-if>` | |
| {doc}`readability-braces-around-statements <readability/braces-around-statements>` | Yes |
| {doc}`readability-const-return-type <readability/const-return-type>` | Yes |
| {doc}`readability-container-contains <readability/container-contains>` | Yes |
| {doc}`readability-container-data-pointer <readability/container-data-pointer>` | Yes |
| {doc}`readability-container-size-empty <readability/container-size-empty>` | Yes |
| {doc}`readability-convert-member-functions-to-static <readability/convert-member-functions-to-static>` | Yes |
| {doc}`readability-delete-null-pointer <readability/delete-null-pointer>` | Yes |
| {doc}`readability-duplicate-include <readability/duplicate-include>` | Yes |
| {doc}`readability-else-after-return <readability/else-after-return>` | Yes |
| {doc}`readability-enum-initial-value <readability/enum-initial-value>` | Yes |
| {doc}`readability-function-cognitive-complexity <readability/function-cognitive-complexity>` | |
| {doc}`readability-function-size <readability/function-size>` | |
| {doc}`readability-identifier-length <readability/identifier-length>` | |
| {doc}`readability-identifier-naming <readability/identifier-naming>` | Yes |
| {doc}`readability-implicit-bool-conversion <readability/implicit-bool-conversion>` | Yes |
| {doc}`readability-inconsistent-declaration-parameter-name <readability/inconsistent-declaration-parameter-name>` | Yes |
| {doc}`readability-inconsistent-ifelse-braces <readability/inconsistent-ifelse-braces>` | Yes |
| {doc}`readability-isolate-declaration <readability/isolate-declaration>` | Yes |
| {doc}`readability-magic-numbers <readability/magic-numbers>` | |
| {doc}`readability-make-member-function-const <readability/make-member-function-const>` | Yes |
| {doc}`readability-math-missing-parentheses <readability/math-missing-parentheses>` | Yes |
| {doc}`readability-misleading-indentation <readability/misleading-indentation>` | |
| {doc}`readability-misplaced-array-index <readability/misplaced-array-index>` | Yes |
| {doc}`readability-named-parameter <readability/named-parameter>` | Yes |
| {doc}`readability-non-const-parameter <readability/non-const-parameter>` | Yes |
| {doc}`readability-operators-representation <readability/operators-representation>` | Yes |
| {doc}`readability-qualified-auto <readability/qualified-auto>` | Yes |
| {doc}`readability-redundant-access-specifiers <readability/redundant-access-specifiers>` | Yes |
| {doc}`readability-redundant-casting <readability/redundant-casting>` | Yes |
| {doc}`readability-redundant-control-flow <readability/redundant-control-flow>` | Yes |
| {doc}`readability-redundant-declaration <readability/redundant-declaration>` | Yes |
| {doc}`readability-redundant-function-ptr-dereference <readability/redundant-function-ptr-dereference>` | Yes |
| {doc}`readability-redundant-inline-specifier <readability/redundant-inline-specifier>` | Yes |
| {doc}`readability-redundant-lambda-parameter-list <readability/redundant-lambda-parameter-list>` | Yes |
| {doc}`readability-redundant-member-init <readability/redundant-member-init>` | Yes |
| {doc}`readability-redundant-nested-if <readability/redundant-nested-if>` | Yes |
| {doc}`readability-redundant-parentheses <readability/redundant-parentheses>` | Yes |
| {doc}`readability-redundant-preprocessor <readability/redundant-preprocessor>` | |
| {doc}`readability-redundant-qualified-alias <readability/redundant-qualified-alias>` | Yes |
| {doc}`readability-redundant-smartptr-get <readability/redundant-smartptr-get>` | Yes |
| {doc}`readability-redundant-string-cstr <readability/redundant-string-cstr>` | Yes |
| {doc}`readability-redundant-string-init <readability/redundant-string-init>` | Yes |
| {doc}`readability-redundant-typename <readability/redundant-typename>` | Yes |
| {doc}`readability-reference-to-constructed-temporary <readability/reference-to-constructed-temporary>` | |
| {doc}`readability-simplify-boolean-expr <readability/simplify-boolean-expr>` | Yes |
| {doc}`readability-simplify-subscript-expr <readability/simplify-subscript-expr>` | Yes |
| {doc}`readability-static-accessed-through-instance <readability/static-accessed-through-instance>` | Yes |
| {doc}`readability-static-definition-in-anonymous-namespace <readability/static-definition-in-anonymous-namespace>` | Yes |
| {doc}`readability-string-compare <readability/string-compare>` | Yes |
| {doc}`readability-suspicious-call-argument <readability/suspicious-call-argument>` | |
| {doc}`readability-trailing-comma <readability/trailing-comma>` | Yes |
| {doc}`readability-trivial-switch <readability/trivial-switch>` | |
| {doc}`readability-uniqueptr-delete-release <readability/uniqueptr-delete-release>` | Yes |
| {doc}`readability-uppercase-literal-suffix <readability/uppercase-literal-suffix>` | Yes |
| {doc}`readability-use-anyofallof <readability/use-anyofallof>` | |
| {doc}`readability-use-concise-preprocessor-directives <readability/use-concise-preprocessor-directives>` | Yes |
| {doc}`readability-use-std-min-max <readability/use-std-min-max>` | Yes |

## Check aliases

| Name | Redirect | Offers fixes |
| --- | --- | --- |
| {doc}`cert-arr39-c <cert/arr39-c>` | {doc}`bugprone-sizeof-expression <bugprone/sizeof-expression>` | |
| {doc}`cert-con36-c <cert/con36-c>` | {doc}`bugprone-spuriously-wake-up-functions <bugprone/spuriously-wake-up-functions>` | |
| {doc}`cert-con54-cpp <cert/con54-cpp>` | {doc}`bugprone-spuriously-wake-up-functions <bugprone/spuriously-wake-up-functions>` | |
| {doc}`cert-ctr56-cpp <cert/ctr56-cpp>` | {doc}`bugprone-pointer-arithmetic-on-polymorphic-object <bugprone/pointer-arithmetic-on-polymorphic-object>` | |
| {doc}`cert-dcl03-c <cert/dcl03-c>` | {doc}`misc-static-assert <misc/static-assert>` | Yes |
| {doc}`cert-dcl16-c <cert/dcl16-c>` | {doc}`readability-uppercase-literal-suffix <readability/uppercase-literal-suffix>` | Yes |
| {doc}`cert-dcl37-c <cert/dcl37-c>` | {doc}`bugprone-reserved-identifier <bugprone/reserved-identifier>` | Yes |
| {doc}`cert-dcl50-cpp <cert/dcl50-cpp>` | {doc}`modernize-avoid-variadic-functions <modernize/avoid-variadic-functions>` | |
| {doc}`cert-dcl51-cpp <cert/dcl51-cpp>` | {doc}`bugprone-reserved-identifier <bugprone/reserved-identifier>` | Yes |
| {doc}`cert-dcl54-cpp <cert/dcl54-cpp>` | {doc}`misc-new-delete-overloads <misc/new-delete-overloads>` | |
| {doc}`cert-dcl58-cpp <cert/dcl58-cpp>` | {doc}`bugprone-std-namespace-modification <bugprone/std-namespace-modification>` | |
| {doc}`cert-dcl59-cpp <cert/dcl59-cpp>` | {doc}`misc-anonymous-namespace-in-header <misc/anonymous-namespace-in-header>` | |
| {doc}`cert-env33-c <cert/env33-c>` | {doc}`bugprone-command-processor <bugprone/command-processor>` | |
| {doc}`cert-err09-cpp <cert/err09-cpp>` | {doc}`misc-throw-by-value-catch-by-reference <misc/throw-by-value-catch-by-reference>` | |
| {doc}`cert-err33-c <cert/err33-c>` | {doc}`bugprone-unused-return-value <bugprone/unused-return-value>` | |
| {doc}`cert-err34-c <cert/err34-c>` | {doc}`bugprone-unchecked-string-to-number-conversion <bugprone/unchecked-string-to-number-conversion>` | |
| {doc}`cert-err52-cpp <cert/err52-cpp>` | {doc}`modernize-avoid-setjmp-longjmp <modernize/avoid-setjmp-longjmp>` | |
| {doc}`cert-err58-cpp <cert/err58-cpp>` | {doc}`bugprone-throwing-static-initialization <bugprone/throwing-static-initialization>` | |
| {doc}`cert-err60-cpp <cert/err60-cpp>` | {doc}`bugprone-exception-copy-constructor-throws <bugprone/exception-copy-constructor-throws>` | |
| {doc}`cert-err61-cpp <cert/err61-cpp>` | {doc}`misc-throw-by-value-catch-by-reference <misc/throw-by-value-catch-by-reference>` | |
| {doc}`cert-exp42-c <cert/exp42-c>` | {doc}`bugprone-suspicious-memory-comparison <bugprone/suspicious-memory-comparison>` | |
| {doc}`cert-exp45-c <cert/exp45-c>` | {doc}`bugprone-assignment-in-selection-statement <bugprone/assignment-in-selection-statement>` | |
| {doc}`cert-fio38-c <cert/fio38-c>` | {doc}`misc-non-copyable-objects <misc/non-copyable-objects>` | |
| {doc}`cert-flp30-c <cert/flp30-c>` | {doc}`bugprone-float-loop-counter <bugprone/float-loop-counter>` | |
| {doc}`cert-flp37-c <cert/flp37-c>` | {doc}`bugprone-suspicious-memory-comparison <bugprone/suspicious-memory-comparison>` | |
| {doc}`cert-int09-c <cert/int09-c>` | {doc}`readability-enum-initial-value <readability/enum-initial-value>` | Yes |
| {doc}`cert-mem57-cpp <cert/mem57-cpp>` | {doc}`bugprone-default-operator-new-on-overaligned-type <bugprone/default-operator-new-on-overaligned-type>` | |
| {doc}`cert-msc24-c <cert/msc24-c>` | {doc}`bugprone-unsafe-functions <bugprone/unsafe-functions>` | |
| {doc}`cert-msc30-c <cert/msc30-c>` | {doc}`misc-predictable-rand <misc/predictable-rand>` | |
| {doc}`cert-msc32-c <cert/msc32-c>` | {doc}`bugprone-random-generator-seed <bugprone/random-generator-seed>` | |
| {doc}`cert-msc33-c <cert/msc33-c>` | {doc}`bugprone-unsafe-functions <bugprone/unsafe-functions>` | |
| {doc}`cert-msc50-cpp <cert/msc50-cpp>` | {doc}`misc-predictable-rand <misc/predictable-rand>` | |
| {doc}`cert-msc51-cpp <cert/msc51-cpp>` | {doc}`bugprone-random-generator-seed <bugprone/random-generator-seed>` | |
| {doc}`cert-msc54-cpp <cert/msc54-cpp>` | {doc}`bugprone-signal-handler <bugprone/signal-handler>` | |
| {doc}`cert-oop11-cpp <cert/oop11-cpp>` | {doc}`performance-move-constructor-init <performance/move-constructor-init>` | |
| {doc}`cert-oop54-cpp <cert/oop54-cpp>` | {doc}`bugprone-unhandled-self-assignment <bugprone/unhandled-self-assignment>` | |
| {doc}`cert-oop57-cpp <cert/oop57-cpp>` | {doc}`bugprone-raw-memory-call-on-non-trivial-type <bugprone/raw-memory-call-on-non-trivial-type>` | |
| {doc}`cert-oop58-cpp <cert/oop58-cpp>` | {doc}`bugprone-copy-constructor-mutates-argument <bugprone/copy-constructor-mutates-argument>` | |
| {doc}`cert-pos44-c <cert/pos44-c>` | {doc}`bugprone-bad-signal-to-kill-thread <bugprone/bad-signal-to-kill-thread>` | |
| {doc}`cert-pos47-c <cert/pos47-c>` | {doc}`concurrency-thread-canceltype-asynchronous <concurrency/thread-canceltype-asynchronous>` | |
| {doc}`cert-sig30-c <cert/sig30-c>` | {doc}`bugprone-signal-handler <bugprone/signal-handler>` | |
| {doc}`cert-str34-c <cert/str34-c>` | {doc}`bugprone-signed-char-misuse <bugprone/signed-char-misuse>` | |
| {doc}`clang-analyzer-core.BitwiseShift <clang-analyzer/core.BitwiseShift>` | [Clang Static Analyzer core.BitwiseShift](https://clang.llvm.org/docs/analyzer/checkers.html#core-bitwiseshift) | |
| {doc}`clang-analyzer-core.CallAndMessage <clang-analyzer/core.CallAndMessage>` | [Clang Static Analyzer core.CallAndMessage](https://clang.llvm.org/docs/analyzer/checkers.html#core-callandmessage) | |
| {doc}`clang-analyzer-core.DivideZero <clang-analyzer/core.DivideZero>` | [Clang Static Analyzer core.DivideZero](https://clang.llvm.org/docs/analyzer/checkers.html#core-dividezero) | |
| {doc}`clang-analyzer-core.NonNullParamChecker <clang-analyzer/core.NonNullParamChecker>` | [Clang Static Analyzer core.NonNullParamChecker](https://clang.llvm.org/docs/analyzer/checkers.html#core-nonnullparamchecker) | |
| {doc}`clang-analyzer-core.NullDereference <clang-analyzer/core.NullDereference>` | [Clang Static Analyzer core.NullDereference](https://clang.llvm.org/docs/analyzer/checkers.html#core-nulldereference) | |
| {doc}`clang-analyzer-core.NullPointerArithm <clang-analyzer/core.NullPointerArithm>` | [Clang Static Analyzer core.NullPointerArithm](https://clang.llvm.org/docs/analyzer/checkers.html#core-nullpointerarithm) | |
| {doc}`clang-analyzer-core.StackAddressEscape <clang-analyzer/core.StackAddressEscape>` | [Clang Static Analyzer core.StackAddressEscape](https://clang.llvm.org/docs/analyzer/checkers.html#core-stackaddressescape) | |
| {doc}`clang-analyzer-core.UndefinedBinaryOperatorResult <clang-analyzer/core.UndefinedBinaryOperatorResult>` | [Clang Static Analyzer core.UndefinedBinaryOperatorResult](https://clang.llvm.org/docs/analyzer/checkers.html#core-undefinedbinaryoperatorresult) | |
| {doc}`clang-analyzer-core.VLASize <clang-analyzer/core.VLASize>` | [Clang Static Analyzer core.VLASize](https://clang.llvm.org/docs/analyzer/checkers.html#core-vlasize) | |
| {doc}`clang-analyzer-core.uninitialized.ArraySubscript <clang-analyzer/core.uninitialized.ArraySubscript>` | [Clang Static Analyzer core.uninitialized.ArraySubscript](https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-arraysubscript) | |
| {doc}`clang-analyzer-core.uninitialized.Assign <clang-analyzer/core.uninitialized.Assign>` | [Clang Static Analyzer core.uninitialized.Assign](https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-assign) | |
| {doc}`clang-analyzer-core.uninitialized.Branch <clang-analyzer/core.uninitialized.Branch>` | [Clang Static Analyzer core.uninitialized.Branch](https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-branch) | |
| {doc}`clang-analyzer-core.uninitialized.CapturedBlockVariable <clang-analyzer/core.uninitialized.CapturedBlockVariable>` | [Clang Static Analyzer core.uninitialized.CapturedBlockVariable](https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-capturedblockvariable) | |
| {doc}`clang-analyzer-core.uninitialized.NewArraySize <clang-analyzer/core.uninitialized.NewArraySize>` | [Clang Static Analyzer core.uninitialized.NewArraySize](https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-newarraysize) | |
| {doc}`clang-analyzer-core.uninitialized.UndefReturn <clang-analyzer/core.uninitialized.UndefReturn>` | [Clang Static Analyzer core.uninitialized.UndefReturn](https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-undefreturn) | |
| {doc}`clang-analyzer-cplusplus.ArrayDelete <clang-analyzer/cplusplus.ArrayDelete>` | [Clang Static Analyzer cplusplus.ArrayDelete](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-arraydelete) | |
| {doc}`clang-analyzer-cplusplus.InnerPointer <clang-analyzer/cplusplus.InnerPointer>` | [Clang Static Analyzer cplusplus.InnerPointer](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-innerpointer) | |
| {doc}`clang-analyzer-cplusplus.Move <clang-analyzer/cplusplus.Move>` | [Clang Static Analyzer cplusplus.Move](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-move) | |
| {doc}`clang-analyzer-cplusplus.NewDelete <clang-analyzer/cplusplus.NewDelete>` | [Clang Static Analyzer cplusplus.NewDelete](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-newdelete) | |
| {doc}`clang-analyzer-cplusplus.NewDeleteLeaks <clang-analyzer/cplusplus.NewDeleteLeaks>` | [Clang Static Analyzer cplusplus.NewDeleteLeaks](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-newdeleteleaks) | |
| {doc}`clang-analyzer-cplusplus.PlacementNew <clang-analyzer/cplusplus.PlacementNew>` | [Clang Static Analyzer cplusplus.PlacementNew](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-placementnew) | |
| {doc}`clang-analyzer-cplusplus.PureVirtualCall <clang-analyzer/cplusplus.PureVirtualCall>` | [Clang Static Analyzer cplusplus.PureVirtualCall](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-purevirtualcall) | |
| {doc}`clang-analyzer-cplusplus.StringChecker <clang-analyzer/cplusplus.StringChecker>` | [Clang Static Analyzer cplusplus.StringChecker](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-stringchecker) | |
| {doc}`clang-analyzer-deadcode.DeadStores <clang-analyzer/deadcode.DeadStores>` | [Clang Static Analyzer deadcode.DeadStores](https://clang.llvm.org/docs/analyzer/checkers.html#deadcode-deadstores) | |
| {doc}`clang-analyzer-fuchsia.HandleChecker <clang-analyzer/fuchsia.HandleChecker>` | [Clang Static Analyzer fuchsia.HandleChecker](https://clang.llvm.org/docs/analyzer/checkers.html#fuchsia-handlechecker) | |
| {doc}`clang-analyzer-nullability.NullPassedToNonnull <clang-analyzer/nullability.NullPassedToNonnull>` | [Clang Static Analyzer nullability.NullPassedToNonnull](https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullpassedtononnull) | |
| {doc}`clang-analyzer-nullability.NullReturnedFromNonnull <clang-analyzer/nullability.NullReturnedFromNonnull>` | [Clang Static Analyzer nullability.NullReturnedFromNonnull](https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullreturnedfromnonnull) | |
| {doc}`clang-analyzer-nullability.NullableDereferenced <clang-analyzer/nullability.NullableDereferenced>` | [Clang Static Analyzer nullability.NullableDereferenced](https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullabledereferenced) | |
| {doc}`clang-analyzer-nullability.NullablePassedToNonnull <clang-analyzer/nullability.NullablePassedToNonnull>` | [Clang Static Analyzer nullability.NullablePassedToNonnull](https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullablepassedtononnull) | |
| {doc}`clang-analyzer-nullability.NullableReturnedFromNonnull <clang-analyzer/nullability.NullableReturnedFromNonnull>` | [Clang Static Analyzer nullability.NullableReturnedFromNonnull](https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullablereturnedfromnonnull) | |
| {doc}`clang-analyzer-optin.core.EnumCastOutOfRange <clang-analyzer/optin.core.EnumCastOutOfRange>` | [Clang Static Analyzer optin.core.EnumCastOutOfRange](https://clang.llvm.org/docs/analyzer/checkers.html#optin-core-enumcastoutofrange) | |
| {doc}`clang-analyzer-optin.core.FixedAddressDereference <clang-analyzer/optin.core.FixedAddressDereference>` | [Clang Static Analyzer optin.core.FixedAddressDereference](https://clang.llvm.org/docs/analyzer/checkers.html#optin-core-fixedaddressdereference) | |
| {doc}`clang-analyzer-optin.core.UnconditionalVAArg <clang-analyzer/optin.core.UnconditionalVAArg>` | [Clang Static Analyzer optin.core.UnconditionalVAArg](https://clang.llvm.org/docs/analyzer/checkers.html#optin-core-unconditionalvaarg) | |
| {doc}`clang-analyzer-optin.cplusplus.UninitializedObject <clang-analyzer/optin.cplusplus.UninitializedObject>` | [Clang Static Analyzer optin.cplusplus.UninitializedObject](https://clang.llvm.org/docs/analyzer/checkers.html#optin-cplusplus-uninitializedobject) | |
| {doc}`clang-analyzer-optin.cplusplus.VirtualCall <clang-analyzer/optin.cplusplus.VirtualCall>` | [Clang Static Analyzer optin.cplusplus.VirtualCall](https://clang.llvm.org/docs/analyzer/checkers.html#optin-cplusplus-virtualcall) | |
| {doc}`clang-analyzer-optin.mpi.MPI-Checker <clang-analyzer/optin.mpi.MPI-Checker>` | [Clang Static Analyzer optin.mpi.MPI-Checker](https://clang.llvm.org/docs/analyzer/checkers.html#optin-mpi-mpi-checker) | |
| {doc}`clang-analyzer-optin.osx.OSObjectCStyleCast <clang-analyzer/optin.osx.OSObjectCStyleCast>` | Clang Static Analyzer optin.osx.OSObjectCStyleCast | |
| {doc}`clang-analyzer-optin.osx.cocoa.localizability.EmptyLocalizationContextChecker <clang-analyzer/optin.osx.cocoa.localizability.EmptyLocalizationContextChecker>` | [Clang Static Analyzer optin.osx.cocoa.localizability.EmptyLocalizationContextChecker](https://clang.llvm.org/docs/analyzer/checkers.html#optin-osx-cocoa-localizability-emptylocalizationcontextchecker) | |
| {doc}`clang-analyzer-optin.osx.cocoa.localizability.NonLocalizedStringChecker <clang-analyzer/optin.osx.cocoa.localizability.NonLocalizedStringChecker>` | [Clang Static Analyzer optin.osx.cocoa.localizability.NonLocalizedStringChecker](https://clang.llvm.org/docs/analyzer/checkers.html#optin-osx-cocoa-localizability-nonlocalizedstringchecker) | |
| {doc}`clang-analyzer-optin.performance.GCDAntipattern <clang-analyzer/optin.performance.GCDAntipattern>` | [Clang Static Analyzer optin.performance.GCDAntipattern](https://clang.llvm.org/docs/analyzer/checkers.html#optin-performance-gcdantipattern) | |
| {doc}`clang-analyzer-optin.performance.Padding <clang-analyzer/optin.performance.Padding>` | [Clang Static Analyzer optin.performance.Padding](https://clang.llvm.org/docs/analyzer/checkers.html#optin-performance-padding) | |
| {doc}`clang-analyzer-optin.portability.UnixAPI <clang-analyzer/optin.portability.UnixAPI>` | [Clang Static Analyzer optin.portability.UnixAPI](https://clang.llvm.org/docs/analyzer/checkers.html#optin-portability-unixapi) | |
| {doc}`clang-analyzer-optin.taint.GenericTaint <clang-analyzer/optin.taint.GenericTaint>` | [Clang Static Analyzer optin.taint.GenericTaint](https://clang.llvm.org/docs/analyzer/checkers.html#optin-taint-generictaint) | |
| {doc}`clang-analyzer-optin.taint.TaintedAlloc <clang-analyzer/optin.taint.TaintedAlloc>` | [Clang Static Analyzer optin.taint.TaintedAlloc](https://clang.llvm.org/docs/analyzer/checkers.html#optin-taint-taintedalloc) | |
| {doc}`clang-analyzer-optin.taint.TaintedDiv <clang-analyzer/optin.taint.TaintedDiv>` | [Clang Static Analyzer optin.taint.TaintedDiv](https://clang.llvm.org/docs/analyzer/checkers.html#optin-taint-tainteddiv) | |
| {doc}`clang-analyzer-osx.API <clang-analyzer/osx.API>` | [Clang Static Analyzer osx.API](https://clang.llvm.org/docs/analyzer/checkers.html#osx-api) | |
| {doc}`clang-analyzer-osx.MIG <clang-analyzer/osx.MIG>` | Clang Static Analyzer osx.MIG | |
| {doc}`clang-analyzer-osx.NumberObjectConversion <clang-analyzer/osx.NumberObjectConversion>` | [Clang Static Analyzer osx.NumberObjectConversion](https://clang.llvm.org/docs/analyzer/checkers.html#osx-numberobjectconversion) | |
| {doc}`clang-analyzer-osx.OSObjectRetainCount <clang-analyzer/osx.OSObjectRetainCount>` | Clang Static Analyzer osx.OSObjectRetainCount | |
| {doc}`clang-analyzer-osx.ObjCProperty <clang-analyzer/osx.ObjCProperty>` | [Clang Static Analyzer osx.ObjCProperty](https://clang.llvm.org/docs/analyzer/checkers.html#osx-objcproperty) | |
| {doc}`clang-analyzer-osx.SecKeychainAPI <clang-analyzer/osx.SecKeychainAPI>` | [Clang Static Analyzer osx.SecKeychainAPI](https://clang.llvm.org/docs/analyzer/checkers.html#osx-seckeychainapi) | |
| {doc}`clang-analyzer-osx.cocoa.AtSync <clang-analyzer/osx.cocoa.AtSync>` | [Clang Static Analyzer osx.cocoa.AtSync](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-atsync) | |
| {doc}`clang-analyzer-osx.cocoa.AutoreleaseWrite <clang-analyzer/osx.cocoa.AutoreleaseWrite>` | [Clang Static Analyzer osx.cocoa.AutoreleaseWrite](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-autoreleasewrite) | |
| {doc}`clang-analyzer-osx.cocoa.ClassRelease <clang-analyzer/osx.cocoa.ClassRelease>` | [Clang Static Analyzer osx.cocoa.ClassRelease](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-classrelease) | |
| {doc}`clang-analyzer-osx.cocoa.Dealloc <clang-analyzer/osx.cocoa.Dealloc>` | [Clang Static Analyzer osx.cocoa.Dealloc](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-dealloc) | |
| {doc}`clang-analyzer-osx.cocoa.IncompatibleMethodTypes <clang-analyzer/osx.cocoa.IncompatibleMethodTypes>` | [Clang Static Analyzer osx.cocoa.IncompatibleMethodTypes](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-incompatiblemethodtypes) | |
| {doc}`clang-analyzer-osx.cocoa.Loops <clang-analyzer/osx.cocoa.Loops>` | [Clang Static Analyzer osx.cocoa.Loops](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-loops) | |
| {doc}`clang-analyzer-osx.cocoa.MissingSuperCall <clang-analyzer/osx.cocoa.MissingSuperCall>` | [Clang Static Analyzer osx.cocoa.MissingSuperCall](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-missingsupercall) | |
| {doc}`clang-analyzer-osx.cocoa.NSAutoreleasePool <clang-analyzer/osx.cocoa.NSAutoreleasePool>` | [Clang Static Analyzer osx.cocoa.NSAutoreleasePool](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-nsautoreleasepool) | |
| {doc}`clang-analyzer-osx.cocoa.NSError <clang-analyzer/osx.cocoa.NSError>` | [Clang Static Analyzer osx.cocoa.NSError](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-nserror) | |
| {doc}`clang-analyzer-osx.cocoa.NilArg <clang-analyzer/osx.cocoa.NilArg>` | [Clang Static Analyzer osx.cocoa.NilArg](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-nilarg) | |
| {doc}`clang-analyzer-osx.cocoa.NonNilReturnValue <clang-analyzer/osx.cocoa.NonNilReturnValue>` | [Clang Static Analyzer osx.cocoa.NonNilReturnValue](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-nonnilreturnvalue) | |
| {doc}`clang-analyzer-osx.cocoa.ObjCGenerics <clang-analyzer/osx.cocoa.ObjCGenerics>` | [Clang Static Analyzer osx.cocoa.ObjCGenerics](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-objcgenerics) | |
| {doc}`clang-analyzer-osx.cocoa.RetainCount <clang-analyzer/osx.cocoa.RetainCount>` | [Clang Static Analyzer osx.cocoa.RetainCount](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-retaincount) | |
| {doc}`clang-analyzer-osx.cocoa.RunLoopAutoreleaseLeak <clang-analyzer/osx.cocoa.RunLoopAutoreleaseLeak>` | [Clang Static Analyzer osx.cocoa.RunLoopAutoreleaseLeak](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-runloopautoreleaseleak) | |
| {doc}`clang-analyzer-osx.cocoa.SelfInit <clang-analyzer/osx.cocoa.SelfInit>` | [Clang Static Analyzer osx.cocoa.SelfInit](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-selfinit) | |
| {doc}`clang-analyzer-osx.cocoa.SuperDealloc <clang-analyzer/osx.cocoa.SuperDealloc>` | [Clang Static Analyzer osx.cocoa.SuperDealloc](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-superdealloc) | |
| {doc}`clang-analyzer-osx.cocoa.UnusedIvars <clang-analyzer/osx.cocoa.UnusedIvars>` | [Clang Static Analyzer osx.cocoa.UnusedIvars](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-unusedivars) | |
| {doc}`clang-analyzer-osx.cocoa.VariadicMethodTypes <clang-analyzer/osx.cocoa.VariadicMethodTypes>` | [Clang Static Analyzer osx.cocoa.VariadicMethodTypes](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-variadicmethodtypes) | |
| {doc}`clang-analyzer-osx.coreFoundation.CFError <clang-analyzer/osx.coreFoundation.CFError>` | [Clang Static Analyzer osx.coreFoundation.CFError](https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-cferror) | |
| {doc}`clang-analyzer-osx.coreFoundation.CFNumber <clang-analyzer/osx.coreFoundation.CFNumber>` | [Clang Static Analyzer osx.coreFoundation.CFNumber](https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-cfnumber) | |
| {doc}`clang-analyzer-osx.coreFoundation.CFRetainRelease <clang-analyzer/osx.coreFoundation.CFRetainRelease>` | [Clang Static Analyzer osx.coreFoundation.CFRetainRelease](https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-cfretainrelease) | |
| {doc}`clang-analyzer-osx.coreFoundation.containers.OutOfBounds <clang-analyzer/osx.coreFoundation.containers.OutOfBounds>` | [Clang Static Analyzer osx.coreFoundation.containers.OutOfBounds](https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-containers-outofbounds) | |
| {doc}`clang-analyzer-osx.coreFoundation.containers.PointerSizedValues <clang-analyzer/osx.coreFoundation.containers.PointerSizedValues>` | [Clang Static Analyzer osx.coreFoundation.containers.PointerSizedValues](https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-containers-pointersizedvalues) | |
| {doc}`clang-analyzer-security.ArrayBound <clang-analyzer/security.ArrayBound>` | [Clang Static Analyzer security.ArrayBound](https://clang.llvm.org/docs/analyzer/checkers.html#security-arraybound) | |
| {doc}`clang-analyzer-security.FloatLoopCounter <clang-analyzer/security.FloatLoopCounter>` | [Clang Static Analyzer security.FloatLoopCounter](https://clang.llvm.org/docs/analyzer/checkers.html#security-floatloopcounter) | |
| {doc}`clang-analyzer-security.MmapWriteExec <clang-analyzer/security.MmapWriteExec>` | [Clang Static Analyzer security.MmapWriteExec](https://clang.llvm.org/docs/analyzer/checkers.html#security-mmapwriteexec) | |
| {doc}`clang-analyzer-security.PointerSub <clang-analyzer/security.PointerSub>` | [Clang Static Analyzer security.PointerSub](https://clang.llvm.org/docs/analyzer/checkers.html#security-pointersub) | |
| {doc}`clang-analyzer-security.PutenvStackArray <clang-analyzer/security.PutenvStackArray>` | Clang Static Analyzer security.PutenvStackArray | |
| {doc}`clang-analyzer-security.SetgidSetuidOrder <clang-analyzer/security.SetgidSetuidOrder>` | Clang Static Analyzer security.SetgidSetuidOrder | |
| {doc}`clang-analyzer-security.VAList <clang-analyzer/security.VAList>` | [Clang Static Analyzer security.VAList](https://clang.llvm.org/docs/analyzer/checkers.html#security-valist) | |
| {doc}`clang-analyzer-security.cert.env.InvalidPtr <clang-analyzer/security.cert.env.InvalidPtr>` | [Clang Static Analyzer security.cert.env.InvalidPtr](https://clang.llvm.org/docs/analyzer/checkers.html#security-cert-env-invalidptr) | |
| {doc}`clang-analyzer-security.insecureAPI.DeprecatedOrUnsafeBufferHandling <clang-analyzer/security.insecureAPI.DeprecatedOrUnsafeBufferHandling>` | [Clang Static Analyzer security.insecureAPI.DeprecatedOrUnsafeBufferHandling](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-deprecatedorunsafebufferhandling) | |
| {doc}`clang-analyzer-security.insecureAPI.UncheckedReturn <clang-analyzer/security.insecureAPI.UncheckedReturn>` | [Clang Static Analyzer security.insecureAPI.UncheckedReturn](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-uncheckedreturn) | |
| {doc}`clang-analyzer-security.insecureAPI.bcmp <clang-analyzer/security.insecureAPI.bcmp>` | [Clang Static Analyzer security.insecureAPI.bcmp](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-bcmp) | |
| {doc}`clang-analyzer-security.insecureAPI.bcopy <clang-analyzer/security.insecureAPI.bcopy>` | [Clang Static Analyzer security.insecureAPI.bcopy](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-bcopy) | |
| {doc}`clang-analyzer-security.insecureAPI.bzero <clang-analyzer/security.insecureAPI.bzero>` | [Clang Static Analyzer security.insecureAPI.bzero](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-bzero) | |
| {doc}`clang-analyzer-security.insecureAPI.decodeValueOfObjCType <clang-analyzer/security.insecureAPI.decodeValueOfObjCType>` | [Clang Static Analyzer security.insecureAPI.decodeValueOfObjCType](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-decodevalueofobjctype) | |
| {doc}`clang-analyzer-security.insecureAPI.getpw <clang-analyzer/security.insecureAPI.getpw>` | [Clang Static Analyzer security.insecureAPI.getpw](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-getpw) | |
| {doc}`clang-analyzer-security.insecureAPI.gets <clang-analyzer/security.insecureAPI.gets>` | [Clang Static Analyzer security.insecureAPI.gets](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-gets) | |
| {doc}`clang-analyzer-security.insecureAPI.mkstemp <clang-analyzer/security.insecureAPI.mkstemp>` | [Clang Static Analyzer security.insecureAPI.mkstemp](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-mkstemp) | |
| {doc}`clang-analyzer-security.insecureAPI.mktemp <clang-analyzer/security.insecureAPI.mktemp>` | [Clang Static Analyzer security.insecureAPI.mktemp](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-mktemp) | |
| {doc}`clang-analyzer-security.insecureAPI.rand <clang-analyzer/security.insecureAPI.rand>` | [Clang Static Analyzer security.insecureAPI.rand](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-rand) | |
| {doc}`clang-analyzer-security.insecureAPI.strcpy <clang-analyzer/security.insecureAPI.strcpy>` | [Clang Static Analyzer security.insecureAPI.strcpy](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-strcpy) | |
| {doc}`clang-analyzer-security.insecureAPI.vfork <clang-analyzer/security.insecureAPI.vfork>` | [Clang Static Analyzer security.insecureAPI.vfork](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-vfork) | |
| {doc}`clang-analyzer-unix.API <clang-analyzer/unix.API>` | [Clang Static Analyzer unix.API](https://clang.llvm.org/docs/analyzer/checkers.html#unix-api) | |
| {doc}`clang-analyzer-unix.BlockInCriticalSection <clang-analyzer/unix.BlockInCriticalSection>` | [Clang Static Analyzer unix.BlockInCriticalSection](https://clang.llvm.org/docs/analyzer/checkers.html#unix-blockincriticalsection) | |
| {doc}`clang-analyzer-unix.Chroot <clang-analyzer/unix.Chroot>` | [Clang Static Analyzer unix.Chroot](https://clang.llvm.org/docs/analyzer/checkers.html#unix-chroot) | |
| {doc}`clang-analyzer-unix.Errno <clang-analyzer/unix.Errno>` | [Clang Static Analyzer unix.Errno](https://clang.llvm.org/docs/analyzer/checkers.html#unix-errno) | |
| {doc}`clang-analyzer-unix.Malloc <clang-analyzer/unix.Malloc>` | [Clang Static Analyzer unix.Malloc](https://clang.llvm.org/docs/analyzer/checkers.html#unix-malloc) | |
| {doc}`clang-analyzer-unix.MallocSizeof <clang-analyzer/unix.MallocSizeof>` | [Clang Static Analyzer unix.MallocSizeof](https://clang.llvm.org/docs/analyzer/checkers.html#unix-mallocsizeof) | |
| {doc}`clang-analyzer-unix.MismatchedDeallocator <clang-analyzer/unix.MismatchedDeallocator>` | [Clang Static Analyzer unix.MismatchedDeallocator](https://clang.llvm.org/docs/analyzer/checkers.html#unix-mismatcheddeallocator) | |
| {doc}`clang-analyzer-unix.StdCLibraryFunctions <clang-analyzer/unix.StdCLibraryFunctions>` | [Clang Static Analyzer unix.StdCLibraryFunctions](https://clang.llvm.org/docs/analyzer/checkers.html#unix-stdclibraryfunctions) | |
| {doc}`clang-analyzer-unix.Stream <clang-analyzer/unix.Stream>` | [Clang Static Analyzer unix.Stream](https://clang.llvm.org/docs/analyzer/checkers.html#unix-stream) | |
| {doc}`clang-analyzer-unix.Vfork <clang-analyzer/unix.Vfork>` | [Clang Static Analyzer unix.Vfork](https://clang.llvm.org/docs/analyzer/checkers.html#unix-vfork) | |
| {doc}`clang-analyzer-unix.cstring.BadSizeArg <clang-analyzer/unix.cstring.BadSizeArg>` | [Clang Static Analyzer unix.cstring.BadSizeArg](https://clang.llvm.org/docs/analyzer/checkers.html#unix-cstring-badsizearg) | |
| {doc}`clang-analyzer-unix.cstring.NotNullTerminated <clang-analyzer/unix.cstring.NotNullTerminated>` | [Clang Static Analyzer unix.cstring.NotNullTerminated](https://clang.llvm.org/docs/analyzer/checkers.html#unix-cstring-notnullterminated) | |
| {doc}`clang-analyzer-unix.cstring.NullArg <clang-analyzer/unix.cstring.NullArg>` | [Clang Static Analyzer unix.cstring.NullArg](https://clang.llvm.org/docs/analyzer/checkers.html#unix-cstring-nullarg) | |
| {doc}`clang-analyzer-unix.cstring.UninitializedRead <clang-analyzer/unix.cstring.UninitializedRead>` | [Clang Static Analyzer unix.cstring.UninitializedRead](https://clang.llvm.org/docs/analyzer/checkers.html#unix-cstring-uninitializedread) | |
| {doc}`clang-analyzer-webkit.NoUncountedMemberChecker <clang-analyzer/webkit.NoUncountedMemberChecker>` | [Clang Static Analyzer webkit.NoUncountedMemberChecker](https://clang.llvm.org/docs/analyzer/checkers.html#webkit-nouncountedmemberchecker) | |
| {doc}`clang-analyzer-webkit.RefCntblBaseVirtualDtor <clang-analyzer/webkit.RefCntblBaseVirtualDtor>` | [Clang Static Analyzer webkit.RefCntblBaseVirtualDtor](https://clang.llvm.org/docs/analyzer/checkers.html#webkit-refcntblbasevirtualdtor) | |
| {doc}`clang-analyzer-webkit.UncountedLambdaCapturesChecker <clang-analyzer/webkit.UncountedLambdaCapturesChecker>` | [Clang Static Analyzer webkit.UncountedLambdaCapturesChecker](https://clang.llvm.org/docs/analyzer/checkers.html#webkit-uncountedlambdacaptureschecker) | |
| {doc}`cppcoreguidelines-avoid-c-arrays <cppcoreguidelines/avoid-c-arrays>` | {doc}`modernize-avoid-c-arrays <modernize/avoid-c-arrays>` | |
| {doc}`cppcoreguidelines-avoid-magic-numbers <cppcoreguidelines/avoid-magic-numbers>` | {doc}`readability-magic-numbers <readability/magic-numbers>` | |
| {doc}`cppcoreguidelines-c-copy-assignment-signature <cppcoreguidelines/c-copy-assignment-signature>` | {doc}`misc-unconventional-assign-operator <misc/unconventional-assign-operator>` | |
| {doc}`cppcoreguidelines-explicit-constructor <cppcoreguidelines/explicit-constructor>` | {doc}`misc-explicit-constructor <misc/explicit-constructor>` | Yes |
| {doc}`cppcoreguidelines-explicit-virtual-functions <cppcoreguidelines/explicit-virtual-functions>` | {doc}`modernize-use-override <modernize/use-override>` | Yes |
| {doc}`cppcoreguidelines-macro-to-enum <cppcoreguidelines/macro-to-enum>` | {doc}`modernize-macro-to-enum <modernize/macro-to-enum>` | Yes |
| {doc}`cppcoreguidelines-narrowing-conversions <cppcoreguidelines/narrowing-conversions>` | {doc}`bugprone-narrowing-conversions <bugprone/narrowing-conversions>` | |
| {doc}`cppcoreguidelines-noexcept-destructor <cppcoreguidelines/noexcept-destructor>` | {doc}`performance-noexcept-destructor <performance/noexcept-destructor>` | Yes |
| {doc}`cppcoreguidelines-noexcept-move-operations <cppcoreguidelines/noexcept-move-operations>` | {doc}`performance-noexcept-move-constructor <performance/noexcept-move-constructor>` | Yes |
| {doc}`cppcoreguidelines-noexcept-swap <cppcoreguidelines/noexcept-swap>` | {doc}`performance-noexcept-swap <performance/noexcept-swap>` | Yes |
| {doc}`cppcoreguidelines-non-private-member-variables-in-classes <cppcoreguidelines/non-private-member-variables-in-classes>` | {doc}`misc-non-private-member-variables-in-classes <misc/non-private-member-variables-in-classes>` | |
| {doc}`cppcoreguidelines-use-default-member-init <cppcoreguidelines/use-default-member-init>` | {doc}`modernize-use-default-member-init <modernize/use-default-member-init>` | Yes |
| {doc}`fuchsia-header-anon-namespaces <fuchsia/header-anon-namespaces>` | {doc}`misc-anonymous-namespace-in-header <misc/anonymous-namespace-in-header>` | |
| {doc}`fuchsia-multiple-inheritance <fuchsia/multiple-inheritance>` | {doc}`misc-multiple-inheritance <misc/multiple-inheritance>` | |
| {doc}`google-build-namespaces <google/build-namespaces>` | {doc}`misc-anonymous-namespace-in-header <misc/anonymous-namespace-in-header>` | |
| {doc}`google-explicit-constructor <google/explicit-constructor>` | {doc}`misc-explicit-constructor <misc/explicit-constructor>` | Yes |
| {doc}`google-readability-braces-around-statements <google/readability-braces-around-statements>` | {doc}`readability-braces-around-statements <readability/braces-around-statements>` | Yes |
| {doc}`google-readability-casting <google/readability-casting>` | {doc}`modernize-avoid-c-style-cast <modernize/avoid-c-style-cast>` | Yes |
| {doc}`google-readability-function-size <google/readability-function-size>` | {doc}`readability-function-size <readability/function-size>` | |
| {doc}`google-readability-namespace-comments <google/readability-namespace-comments>` | {doc}`llvm-namespace-comment <llvm/namespace-comment>` | |
| {doc}`llvm-else-after-return <llvm/else-after-return>` | {doc}`readability-else-after-return <readability/else-after-return>` | Yes |
| {doc}`llvm-qualified-auto <llvm/qualified-auto>` | {doc}`readability-qualified-auto <readability/qualified-auto>` | Yes |
| {doc}`performance-faster-string-find <performance/faster-string-find>` | {doc}`performance-prefer-single-char-overloads <performance/prefer-single-char-overloads>` | Yes |

```{title} clang-tidy - abseil-cleanup-ctad
```

# abseil-cleanup-ctad

Suggests switching the initialization pattern of `absl::Cleanup`
instances from the factory function to class template argument
deduction (CTAD), in C++17 and higher.

```cpp
auto c1 = absl::MakeCleanup([] {});

const auto c2 = absl::MakeCleanup(std::function<void()>([] {}));
```

becomes

```cpp
absl::Cleanup c1 = [] {};

const absl::Cleanup c2 = std::function<void()>([] {});
```

```{title} clang-tidy - abseil-duration-addition
```

# abseil-duration-addition

Checks for cases where addition should be performed in the `absl::Time`
domain. When adding two values, and one is known to be an `absl::Time`,
we can infer that the other should be interpreted as an `absl::Duration`
of a similar scale, and make that inference explicit.

Examples:

```cpp
// Original - Addition in the integer domain
int x;
absl::Time t;
int result = absl::ToUnixSeconds(t) + x;

// Suggestion - Addition in the absl::Time domain
int result = absl::ToUnixSeconds(t + absl::Seconds(x));
```

```{title} clang-tidy - abseil-duration-comparison
```

# abseil-duration-comparison

Checks for comparisons which should be in the `absl::Duration` domain instead
of the floating point or integer domains.

N.B.: In cases where a `Duration` was being converted to an integer and then
compared against a floating-point value, truncation during the `Duration`
conversion might yield a different result. In practice this is very rare, and
still indicates a bug which should be fixed.

Examples:

```cpp
// Original - Comparison in the floating point domain
double x;
absl::Duration d;
if (x < absl::ToDoubleSeconds(d)) ...

// Suggested - Compare in the absl::Duration domain instead
if (absl::Seconds(x) < d) ...

// Original - Comparison in the integer domain
int x;
absl::Duration d;
if (x < absl::ToInt64Microseconds(d)) ...

// Suggested - Compare in the absl::Duration domain instead
if (absl::Microseconds(x) < d) ...
```

```{title} clang-tidy - abseil-duration-conversion-cast
```

# abseil-duration-conversion-cast

Checks for casts of `absl::Duration` conversion functions, and recommends
the right conversion function instead.

Examples:

```cpp
// Original - Cast from a double to an integer
absl::Duration d;
int i = static_cast<int>(absl::ToDoubleSeconds(d));

// Suggested - Use the integer conversion function directly.
int i = absl::ToInt64Seconds(d);

// Original - Cast from a double to an integer
absl::Duration d;
double x = static_cast<double>(absl::ToInt64Seconds(d));

// Suggested - Use the integer conversion function directly.
double x = absl::ToDoubleSeconds(d);
```

Note: In the second example, the suggested fix could yield a different result,
as the conversion to integer could truncate. In practice, this is very rare,
and you should use `absl::Trunc` to perform this operation explicitly instead.

```{title} clang-tidy - abseil-duration-division
```

# abseil-duration-division

`absl::Duration` arithmetic works like it does with integers. That means that
division of two `absl::Duration` objects returns an `int64` with any
fractional component truncated toward 0.
See [this link](https://github.com/abseil/abseil-cpp/blob/29ff6d4860070bf8fcbd39c8805d0c32d56628a3/absl/time/time.h#L137)
for more information on arithmetic with `absl::Duration`.

For example:

```cpp
absl::Duration d = absl::Seconds(3.5);
int64 sec1 = d / absl::Seconds(1); // Truncates toward 0.
int64 sec2 = absl::ToInt64Seconds(d); // Equivalent to division.
assert(sec1 == 3 && sec2 == 3);

double dsec = d / absl::Seconds(1); // WRONG: Still truncates toward 0.
assert(dsec == 3.0);
```

If you want floating-point division, you should use either the
`absl::FDivDuration()` function, or one of the unit conversion functions such
as `absl::ToDoubleSeconds()`. For example:

```cpp
absl::Duration d = absl::Seconds(3.5);
double dsec1 = absl::FDivDuration(d, absl::Seconds(1)); // GOOD: No truncation.
double dsec2 = absl::ToDoubleSeconds(d); // GOOD: No truncation.
assert(dsec1 == 3.5 && dsec2 == 3.5);
```

This check looks for uses of `absl::Duration` division that is done in a
floating-point context, and recommends the use of a function that returns a
floating-point value.

```{title} clang-tidy - abseil-duration-factory-float
```

# abseil-duration-factory-float

Checks for cases where the floating-point overloads of various
`absl::Duration` factory functions are called when the more-efficient
integer versions could be used instead.

This check will not suggest fixes for literals which contain fractional
floating point values or non-literals. It will suggest removing
superfluous casts.

Examples:

```cpp
// Original - Providing a floating-point literal.
absl::Duration d = absl::Seconds(10.0);

// Suggested - Use an integer instead.
absl::Duration d = absl::Seconds(10);

// Original - Explicitly casting to a floating-point type.
absl::Duration d = absl::Seconds(static_cast<double>(10));

// Suggested - Remove the explicit cast
absl::Duration d = absl::Seconds(10);
```

```{title} clang-tidy - abseil-duration-factory-scale
```

# abseil-duration-factory-scale

Checks for cases where arguments to `absl::Duration` factory functions are
scaled internally and could be changed to a different factory function. This
check also looks for arguments with a zero value and suggests using
`absl::ZeroDuration()` instead.

Examples:

```cpp
// Original - Internal multiplication.
int x;
absl::Duration d = absl::Seconds(60 * x);

// Suggested - Use absl::Minutes instead.
absl::Duration d = absl::Minutes(x);

// Original - Internal division.
int y;
absl::Duration d = absl::Milliseconds(y / 1000.);

// Suggested - Use absl:::Seconds instead.
absl::Duration d = absl::Seconds(y);

// Original - Zero-value argument.
absl::Duration d = absl::Hours(0);

// Suggested = Use absl::ZeroDuration instead
absl::Duration d = absl::ZeroDuration();
```

```{title} clang-tidy - abseil-duration-subtraction
```

# abseil-duration-subtraction

Checks for cases where subtraction should be performed in the
`absl::Duration` domain. When subtracting two values, and the first one is
known to be a conversion from `absl::Duration`, we can infer that the second
should also be interpreted as an `absl::Duration`, and make that inference
explicit.

Examples:

```cpp
// Original - Subtraction in the double domain
double x;
absl::Duration d;
double result = absl::ToDoubleSeconds(d) - x;

// Suggestion - Subtraction in the absl::Duration domain instead
double result = absl::ToDoubleSeconds(d - absl::Seconds(x));

// Original - Subtraction of two Durations in the double domain
absl::Duration d1, d2;
double result = absl::ToDoubleSeconds(d1) - absl::ToDoubleSeconds(d2);

// Suggestion - Subtraction in the absl::Duration domain instead
double result = absl::ToDoubleSeconds(d1 - d2);
```

Note: As with other `clang-tidy` checks, it is possible that multiple fixes
may overlap (as in the case of nested expressions), so not all occurrences can
be transformed in one run. In particular, this may occur for nested subtraction
expressions. Running `clang-tidy` multiple times will find and fix these
overlaps.

```{title} clang-tidy - abseil-duration-unnecessary-conversion
```

# abseil-duration-unnecessary-conversion

Finds and fixes cases where `absl::Duration` values are being converted to
numeric types and back again.

Floating-point examples:

```cpp
// Original - Conversion to double and back again
absl::Duration d1;
absl::Duration d2 = absl::Seconds(absl::ToDoubleSeconds(d1));

// Suggestion - Remove unnecessary conversions
absl::Duration d2 = d1;

// Original - Division to convert to double and back again
absl::Duration d2 = absl::Seconds(absl::FDivDuration(d1, absl::Seconds(1)));

// Suggestion - Remove division and conversion
absl::Duration d2 = d1;
```

Integer examples:

```cpp
// Original - Conversion to integer and back again
absl::Duration d1;
absl::Duration d2 = absl::Hours(absl::ToInt64Hours(d1));

// Suggestion - Remove unnecessary conversions
absl::Duration d2 = d1;

// Original - Integer division followed by conversion
absl::Duration d2 = absl::Seconds(d1 / absl::Seconds(1));

// Suggestion - Remove division and conversion
absl::Duration d2 = d1;
```

Unwrapping scalar operations:

```cpp
// Original - Multiplication by a scalar
absl::Duration d1;
absl::Duration d2 = absl::Seconds(absl::ToInt64Seconds(d1) * 2);

// Suggestion - Remove unnecessary conversion
absl::Duration d2 = d1 * 2;
```

Note: Converting to an integer and back to an `absl::Duration` might be a
truncating operation if the value is not aligned to the scale of conversion.
In the rare case where this is the intended result, callers should use
`absl::Trunc` to truncate explicitly.

```{title} clang-tidy - abseil-faster-strsplit-delimiter
```

# abseil-faster-strsplit-delimiter

Finds instances of `absl::StrSplit()` or `absl::MaxSplits()` where the
delimiter is a single character string literal and replaces with a character.
The check will offer a suggestion to change the string literal into a
character. It will also catch code using `absl::ByAnyChar()` for just a
single character and will transform that into a single character as well.

These changes will give the same result, but using characters rather than
single character string literals is more efficient and readable.

Examples:

```cpp
// Original - the argument is a string literal.
for (auto piece : absl::StrSplit(str, "B")) {

// Suggested - the argument is a character, which causes the more efficient
// overload of absl::StrSplit() to be used.
for (auto piece : absl::StrSplit(str, 'B')) {

// Original - the argument is a string literal inside absl::ByAnyChar call.
for (auto piece : absl::StrSplit(str, absl::ByAnyChar("B"))) {

// Suggested - the argument is a character, which causes the more efficient
// overload of absl::StrSplit() to be used and we do not need absl::ByAnyChar
// anymore.
for (auto piece : absl::StrSplit(str, 'B')) {

// Original - the argument is a string literal inside absl::MaxSplits call.
for (auto piece : absl::StrSplit(str, absl::MaxSplits("B", 1))) {

// Suggested - the argument is a character, which causes the more efficient
// overload of absl::StrSplit() to be used.
for (auto piece : absl::StrSplit(str, absl::MaxSplits('B', 1))) {
```

```{title} clang-tidy - abseil-no-internal-dependencies
```

# abseil-no-internal-dependencies

Warns if code using Abseil depends on internal details. If something is in a
namespace that includes the word "internal", code is not allowed to depend upon
it because it's an implementation detail. They cannot friend it, include it,
you mention it or refer to it in any way. Doing so violates Abseil's
compatibility guidelines and may result in breakage. See
<https://abseil.io/about/compatibility> for more information.

The following cases will result in warnings:

```cpp
absl::strings_internal::foo();
// warning triggered on this line
class foo {
 friend struct absl::container_internal::faa;
 // warning triggered on this line
};
absl::memory_internal::MakeUniqueResult();
// warning triggered on this line
```

```{title} clang-tidy - abseil-no-namespace
```

# abseil-no-namespace

Ensures code does not open `namespace absl` as that violates Abseil's
compatibility guidelines. Code should not open `namespace absl` as that
conflicts with Abseil's compatibility guidelines and may result in breakage.

Any code that uses:

```cpp
namespace absl {
 ...
}
```

will be prompted with a warning.

See [the full Abseil compatibility guidelines](https://abseil.io/about/compatibility) for more information.

```{title} clang-tidy - abseil-redundant-strcat-calls
```

# abseil-redundant-strcat-calls

Suggests removal of unnecessary calls to `absl::StrCat` when the result is
being passed to another call to `absl::StrCat` or `absl::StrAppend`.

The extra calls cause unnecessary temporary strings to be constructed. Removing
them makes the code smaller and faster.

Examples:

```cpp
std::string s = absl::StrCat("A", absl::StrCat("B", absl::StrCat("C", "D")));
//before

std::string s = absl::StrCat("A", "B", "C", "D");
//after

absl::StrAppend(&s, absl::StrCat("E", "F", "G"));
//before

absl::StrAppend(&s, "E", "F", "G");
//after
```

```{title} clang-tidy - abseil-str-cat-append
```

# abseil-str-cat-append

Flags uses of `absl::StrCat()` to append to a `std::string`. Suggests
`absl::StrAppend()` should be used instead.

The extra calls cause unnecessary temporary strings to be constructed. Removing
them makes the code smaller and faster.

```cpp
a = absl::StrCat(a, b); // Use absl::StrAppend(&a, b) instead.
```

Does not diagnose cases where `absl::StrCat()` is used as a template
argument for a functor.

```{title} clang-tidy - abseil-string-find-startswith
```

# abseil-string-find-startswith

Checks whether a `std::string::find()` or `std::string::rfind()` (and
corresponding `std::string_view` methods) result is compared with 0, and
suggests replacing with `absl::StartsWith()`. This is both a readability and
performance issue.

`starts_with` was added as a built-in function on those types in C++20. If
available, prefer enabling {doc}`modernize-use-starts-ends-with
<../modernize/use-starts-ends-with>` instead of this check.

```cpp
string s = "...";
if (s.find("Hello World") == 0) { /* do something */ }
if (s.rfind("Hello World", 0) == 0) { /* do something */ }
```

becomes

```cpp
string s = "...";
if (absl::StartsWith(s, "Hello World")) { /* do something */ }
if (absl::StartsWith(s, "Hello World")) { /* do something */ }
```

## Options

```{option} StringLikeClasses

Semicolon-separated list of names of string-like classes. By default both
`std::basic_string` and `std::basic_string_view` are considered. The list
of methods to be considered is fixed.
```

```{option} IncludeStyle

A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.
```

```{option} AbseilStringsMatchHeader

The location of Abseil's `strings/match.h`. Defaults to
`absl/strings/match.h`.
```

```{title} clang-tidy - abseil-string-find-str-contains
```

# abseil-string-find-str-contains

Finds `s.find(...) == string::npos` comparisons (for various string-like
types) and suggests replacing with `absl::StrContains()`.

This improves readability and reduces the likelihood of accidentally mixing
`find()` and `npos` from different string-like types.

By default, "string-like types" includes `::std::basic_string`,
`::std::basic_string_view`, and `::absl::string_view`. See the
StringLikeClasses option to change this.

```cpp
std::string s = "...";
if (s.find("Hello World") == std::string::npos) { /* do something */ }

absl::string_view a = "...";
if (absl::string_view::npos != a.find("Hello World")) { /* do something */ }
```

becomes

```cpp
std::string s = "...";
if (!absl::StrContains(s, "Hello World")) { /* do something */ }

absl::string_view a = "...";
if (absl::StrContains(a, "Hello World")) { /* do something */ }
```

## Options

```{option} StringLikeClasses

Semicolon-separated list of names of string-like classes. By default includes
`::std::basic_string`, `::std::basic_string_view`, and
`::absl::string_view`.
```

```{option} IncludeStyle

A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.
```

```{option} AbseilStringsMatchHeader

The location of Abseil's `strings/match.h`. Defaults to
`absl/strings/match.h`.
```

```{title} clang-tidy - abseil-time-comparison
```

# abseil-time-comparison

Prefer comparisons in the `absl::Time` domain instead of the integer domain.

N.B.: In cases where an `absl::Time` is being converted to an integer,
alignment may occur. If the comparison depends on this alignment, doing the
comparison in the `absl::Time` domain may yield a different result. In
practice this is very rare, and still indicates a bug which should be fixed.

Examples:

```cpp
// Original - Comparison in the integer domain
int x;
absl::Time t;
if (x < absl::ToUnixSeconds(t)) ...

// Suggested - Compare in the absl::Time domain instead
if (absl::FromUnixSeconds(x) < t) ...
```

```{title} clang-tidy - abseil-time-subtraction
```

# abseil-time-subtraction

Finds and fixes `absl::Time` subtraction expressions to do subtraction
in the Time domain instead of the numeric domain.

There are two cases of Time subtraction in which deduce additional type
information:

- When the result is an `absl::Duration` and the first argument is an
 `absl::Time`.
- When the second argument is a `absl::Time`.

In the first case, we must know the result of the operation, since without that
the second operand could be either an `absl::Time` or an `absl::Duration`.
In the second case, the first operand *must* be an `absl::Time`, because
subtracting an `absl::Time` from an `absl::Duration` is not defined.

Examples:

```cpp
int x;
absl::Time t;

// Original - absl::Duration result and first operand is an absl::Time.
absl::Duration d = absl::Seconds(absl::ToUnixSeconds(t) - x);

// Suggestion - Perform subtraction in the Time domain instead.
absl::Duration d = t - absl::FromUnixSeconds(x);

// Original - Second operand is an absl::Time.
int i = x - absl::ToUnixSeconds(t);

// Suggestion - Perform subtraction in the Time domain instead.
int i = absl::ToInt64Seconds(absl::FromUnixSeconds(x) - t);
```

```{title} clang-tidy - abseil-unchecked-statusor-access
```

# abseil-unchecked-statusor-access

This check identifies unsafe accesses to values contained in
`absl::StatusOr<T>` objects. Below we will refer to this type as
`StatusOr<T>`.

An access to the value of an `StatusOr<T>` occurs when one of its
`value`, `operator*`, or `operator->` member functions is invoked.
To align with common misconceptions, the check considers these member
functions as equivalent, even though there are subtle differences
related to exceptions vs. undefined behavior.

An access to the value of a `StatusOr<T>` is considered safe if and
only if code in the local scope (e.g. function body) ensures that the
status of the `StatusOr<T>` is ok in all possible execution paths that
can reach the access. That should happen either through an explicit
check, using the `StatusOr<T>::ok` member function, or by constructing
the `StatusOr<T>` in a way that shows that its status is unambiguously
ok (e.g. by passing a value to its constructor).

Below we list some examples of safe and unsafe `StatusOr<T>` access
patterns.

Note: If the check isn't behaving as you would have expected on a code
snippet, please [report it](http://github.com/llvm/llvm-project/issues/new).

## False negatives

This check generally does **not** generate false negatives. That means that if
an access is not marked as unsafe, it is provably safe. If it cannot prove an
access safe, it is assumed to be unsafe. In some cases, the static analysis
cannot prove an access safe even though it is, for a variety of reasons (e.g.
unmodelled invariants of functions called). In these cases, the analysis does
produce false positive reports.

That being said, there are some heuristics used that in very rare cases might
be incorrect:

- [a const method accessor (without arguments) that returns different
 values when called multiple times](#functionstability).

If you think the check generated a false negative, please [report
it](http://github.com/llvm/llvm-project/issues/new).

## Known limitations

This is a non-exhaustive list of constructs that are currently not
modelled in the check and will lead to false positives:

- [Checking a StatusOr and then capturing it in a lambda](#lambdas)
- [Indexing into a container with the same index](#containers)
- [Project specific helper-functions](#uncommonapi),
- [Functions with a stable return value](#functionstability)
- **Any** [cross-function reasoning](#crossfunction). This is by
 design and will not change in the future.

## Checking if the status is ok, then accessing the value

The check recognizes all straightforward ways for checking the status
and accessing the value contained in a `StatusOr<T>` object. For
example:

```cpp
void f(absl::StatusOr<int> x) {
 if (x.ok()) {
 use(*x);
 }
}
```

## Checking if the status is ok, then accessing the value from a copy

The criteria that the check uses is semantic, not syntactic. It
recognizes when a copy of the `StatusOr<T>` object being accessed is
known to have ok status. For example:

```cpp
void f(absl::StatusOr<int> x1) {
 if (x1.ok()) {
 absl::optional<int> x2 = x1;
 use(*x2);
 }
}
```

## Ensuring that the status is ok using common macros

The check is aware of common macros like `ABSL_CHECK` or `ABSL_CHECK_OK`.
Those can be used to ensure that the status of a `StatusOr<T>` object
is ok. For example:

```cpp
void f(absl::StatusOr<int> x) {
 ABSL_CHECK_OK(x);
 use(*x);
}
```

## Ensuring that the status is ok using googletest macros

The check is aware of `googletest` (or `gtest`) macros and matchers.
Accessing the value of a `StatusOr<T>` object is considered safe if it
is preceded by an `ASSERT_` macro that ensures the status is ok.
For example:

```cpp
TEST(MySuite, MyTest) {
 absl::StatusOr<int> x = foo();
 ASSERT_OK(x);
 use(*x);
}

TEST(MySuite, MyOtherTest) {
 absl::StatusOr<int> x = foo();
 ASSERT_THAT(x, absl_testing::IsOk());
 use(*x);
}
```

The following `googletest` macros are supported:

- `ASSERT_OK(...)`
- `ASSERT_TRUE(...)`
- `ASSERT_FALSE(...)`
- `ASSERT_THAT(...)`

The following matchers are supported:

- `IsOk()`
- `StatusIs(...)`
- `IsOkAndHolds(...)`
- `CanonicalStatusIs(...)`

**Note**: `EXPECT_` macros (like `EXPECT_OK` or `EXPECT_TRUE(x.ok())`)
do **not** make subsequent accesses safe because they do not terminate the
test execution.

## Ensuring that the status is ok, then accessing the value in a correlated branch

The check is aware of correlated branches in the code and can figure out
when a `StatusOr<T>` object is ensured to have ok status on all
execution paths that lead to an access. For example:

```cpp
void f(absl::StatusOr<int> x) {
 bool safe = false;
 if (x.ok() && SomeOtherCondition()) {
 safe = true;
 }
 // ... more code...
 if (safe) {
 use(*x);
 }
}
```

## Accessing the value without checking the status

The check flags accesses to the value that are not locally guarded by a
status check:

```cpp
void f1(absl::StatusOr<int> x) {
 use(*x); // unsafe: it is unclear whether the status of `x` is ok.
}

void f2(absl::StatusOr<MyStruct> x) {
 use(x->member); // unsafe: it is unclear whether the status of `x` is ok.
}

void f3(absl::StatusOr<int> x) {
 use(x.value()); // unsafe: it is unclear whether the status of `x` is ok.
}
```

Use `ABSL_CHECK_OK` to signal that you knowingly want to crash on
non-OK values.

NOTE: Even though using `.value()` on a non-`ok()` `StatusOr` is defined
to crash, it is often unintentional. That is why our checker flags those as
well.

## Accessing the value in the wrong branch

The check is aware of the state of a `StatusOr<T>` object in different
branches of the code. For example:

```cpp
void f(absl::StatusOr<int> x) {
 if (x.ok()) {
 } else {
 use(*x); // unsafe: it is clear that the status of `x` is *not* ok.
 }
}
```

(functionstability)=

## Assuming a function result to be stable

The check is aware that function results might not be stable. That is,
consecutive calls to the same function might return different values.
For example:

```cpp
void f(Foo foo) {
 if (foo.x().ok()) {
 use(*foo.x()); // unsafe: it is unclear whether the status of `foo.x()` is ok.
 }
}
```

In such cases it is best to store the result of the function call in a
local variable and use it to access the value. For example:

```cpp
void f(Foo foo) {
 if (const auto& x = foo.x(); x.ok()) {
 use(*x);
 }
}
```

The check **does** assume that `const`-qualified accessor functions
return a stable value if no non-const function was called between the
two calls:

```cpp
class Foo {
 const absl::StatusOr<int>& get() const {
 [...];
 }
}
void f(Foo foo) {
 if (foo.get().ok()) {
 use(*foo.get());
 }
}
```

If there is a call to a non-`const`-qualified function, the check
assumes the return value of the accessor was mutated.

```cpp
class Foo {
 const absl::StatusOr<int>& get() const {
 [...];
 }
 void mutate();
}
void f(Foo foo) {
 if (foo.get().ok()) {
 foo.mutate();
 use(*foo.get()); // unsafe: `mutate()` might have changed the state of the object
 }
}
```

(uncommonapi)=

## Relying on invariants of uncommon APIs

The check is unaware of invariants of uncommon APIs. For example:

```cpp
void f(Foo foo) {
 if (foo.HasProperty("bar")) {
 use(*foo.GetProperty("bar")); // unsafe: it is unclear whether the status of `foo.GetProperty("bar")` is ok.
 }
}
```

In such cases it is best to check explicitly that the status of the
`StatusOr<T>` object is ok. For example:

```cpp
void f(Foo foo) {
 if (const auto& property = foo.GetProperty("bar"); property.ok()) {
 use(*property);
 }
}
```

(crossfunction)=

## Checking if the `StatusOr<T>` is ok, then passing it to another function

The check relies on local reasoning. The check and value access must
both happen in the same function. An access is considered unsafe even if
the caller of the function performing the access ensures that the status
of the `StatusOr<T>` is ok. For example:

```cpp
void g(absl::StatusOr<int> x) {
 use(*x); // unsafe: it is unclear whether the status of `x` is ok.
}

void f(absl::StatusOr<int> x) {
 if (x.ok()) {
 g(x);
 }
}
```

In such cases it is best to either pass the value directly when calling
a function or check that the status of the `StatusOr<T>` is ok in the
local scope of the callee. For example:

```cpp
void g(int val) {
 use(val);
}

void f(absl::StatusOr<int> x) {
 if (x.ok()) {
 g(*x);
 }
}
```

## Aliases created via `using` declarations

The check is aware of aliases of `StatusOr<T>` types that are created
via `using` declarations. For example:

```cpp
using StatusOrInt = absl::StatusOr<int>;

void f(StatusOrInt x) {
 use(*x); // unsafe: it is unclear whether the status of `x` is ok.
}
```

## Containers

The check is more strict than necessary when it comes to containers of
`StatusOr<T>` values. Simply checking that the status of an element of
a container is ok is not sufficient to deem accessing it safe. For
example:

```cpp
void f(std::vector<absl::StatusOr<int>> x) {
 if (x[0].ok()) {
 use(*x[0]); // unsafe: it is unclear whether the status of `x[0]` is ok.
 }
}
```

One needs to grab a reference to a particular object and use that
instead:

```cpp
void f(std::vector<absl::StatusOr<int>> x) {
 absl::StatusOr<int>& x0 = x[0];
 if (x0.ok()) {
 use(*x0);
 }
}
```

A future version could improve the understanding of more safe usage
patterns that involve containers.

## Lambdas

The check is capable of reporting unsafe `StatusOr<T>` accesses in
lambdas, but isn’t smart enough to propagate information from the
surrounding context through the lambda. This means that the following
pattern will be reported as an unsafe access:

```cpp
void f(absl::StatusOr<int> x) {
 if (x.ok()) {
 [&x]() {
 use(*x); // unsafe: it is unclear whether the status of `x` is ok.
 }
 }
}
```

To avoid the issue, you should instead capture the contained object,
either by value or by reference. An init-capture is useful for this,
here capturing by reference:

```cpp
void f(absl::StatusOr<int> x) {
 if (x.ok()) {
 [&x = *x]() {
 use(x);
 }
 }
}
```

Alternatively you could add a check inside the lambda where the value is
accessed:

```cpp
void f(absl::StatusOr<int> x) {
 [&x]() {
 if (x.ok()) {
 use(*x);
 }
 }
}
```

## Reasoning about integers

Because it uses a simple SAT solver, the check cannot reason about integer
inequalities. For instance, the following will result in a false positive:

```cpp
void f(int n, absl::StatusOr<int> x) {
 if (n > 0)
 CHECK_OK(x);
 if (n > 1)
 return *x; // false positive
 return 0;
}
```

In fact, currently this is also the case if the two conditions are identical.

```{title} clang-tidy - abseil-upgrade-duration-conversions
```

# abseil-upgrade-duration-conversions

Finds calls to `absl::Duration` arithmetic operators and factories whose
argument needs an explicit cast to continue compiling after upcoming API
changes.

The operators `*=`, `/=`, `*`, and `/` for `absl::Duration` currently
accept an argument of class type that is convertible to an arithmetic type.
Such a call currently converts the value to an `int64_t`, even in a case such
as `std::atomic<float>` that would result in lossy conversion.

Additionally, the `absl::Duration` factory functions (`absl::Hours`,
`absl::Minutes`, etc) currently accept an `int64_t` or a floating-point
type. Similar to the arithmetic operators, calls with an argument of class type
that is convertible to an arithmetic type go through the `int64_t` path.

These operators and factories will be changed to only accept arithmetic types
to prevent unintended behavior. After these changes are released, passing an
argument of class type will no longer compile, even if the type is implicitly
convertible to an arithmetic type.

Here are example fixes created by this check:

```cpp
std::atomic<int> a;
absl::Duration d = absl::Milliseconds(a);
d *= a;
```

becomes

```cpp
std::atomic<int> a;
absl::Duration d = absl::Milliseconds(static_cast<int64_t>(a));
d *= static_cast<int64_t>(a);
```

Note that this check always adds a cast to `int64_t` in order to preserve the
current behavior of user code. It is possible that this uncovers unintended
behavior due to types implicitly convertible to a floating-point type.

```{title} clang-tidy - altera-id-dependent-backward-branch
```

# altera-id-dependent-backward-branch

Finds ID-dependent variables and fields that are used within loops. This causes
branches to occur inside the loops, and thus leads to performance degradation.

```cpp
// The following code will produce a warning because this ID-dependent
// variable is used in a loop condition statement.
int ThreadID = get_local_id(0);

// The following loop will produce a warning because the loop condition
// statement depends on an ID-dependent variable.
for (int i = 0; i < ThreadID; ++i) {
 std::cout << i << std::endl;
}

// The following loop will not produce a warning, because the ID-dependent
// variable is not used in the loop condition statement.
for (int i = 0; i < 100; ++i) {
 std::cout << ThreadID << std::endl;
}
```

Based on the [Altera SDK for OpenCL: Best Practices Guide](https://www.altera.com/en_US/pdfs/literature/hb/opencl-sdk/aocl_optimization_guide.pdf).

```{title} clang-tidy - altera-kernel-name-restriction
```

# altera-kernel-name-restriction

Finds kernel files and include directives whose filename is `kernel.cl`,
`Verilog.cl`, or `VHDL.cl`. The check is case insensitive.

Such kernel file names cause the offline compiler to generate intermediate
design files that have the same names as certain internal files, which
leads to a compilation error.

Based on the Guidelines for Naming the Kernel section in the
[Intel FPGA SDK for OpenCL Pro Edition: Programming Guide](https://www.intel.com/content/www/us/en/programmable/documentation/mwh1391807965224.html#ewa1412973930963).

**Title:** clang-tidy - altera-single-work-item-barrier

# altera-single-work-item-barrier

Finds OpenCL kernel functions that call a barrier function but do not call
an ID function (`get_local_id`, `get_local_id`, `get_group_id`, or
`get_local_linear_id`).

These kernels may be viable single work-item kernels, but will be forced to
execute as NDRange kernels if using a newer version of the Altera Offline
Compiler (>= v17.01).

If using an older version of the Altera Offline Compiler, these kernel
functions will be treated as single work-item kernels, which could be
inefficient or lead to errors if NDRange semantics were intended.

Based on the `Altera SDK for OpenCL: Best Practices Guide
<https://www.altera.com/en_US/pdfs/literature/hb/opencl-sdk/aocl_optimization_guide.pdf>`_.

Examples:

```c++
// error: function calls barrier but does not call an ID function.
void __kernel barrier_no_id(__global int * foo, int size) {
 for (int i = 0; i < 100; i++) {
 foo[i] += 5;
 }
 barrier(CLK_GLOBAL_MEM_FENCE);
}

// ok: function calls barrier and an ID function.
void __kernel barrier_with_id(__global int * foo, int size) {
 for (int i = 0; i < 100; i++) {
 int tid = get_global_id(0);
 foo[tid] += 5;
 }
 barrier(CLK_GLOBAL_MEM_FENCE);
}

// ok with AOC Version 17.01: the reqd_work_group_size turns this into
// an NDRange.
__attribute__((reqd_work_group_size(2,2,2)))
void __kernel barrier_with_id(__global int * foo, int size) {
 for (int i = 0; i < 100; i++) {
 foo[tid] += 5;
 }
 barrier(CLK_GLOBAL_MEM_FENCE);
}
```

## Options

**Option:** AOCVersion
Defines the version of the Altera Offline Compiler. Defaults to ``1600``
(corresponding to version 16.00).

**Title:** clang-tidy - altera-struct-pack-align

# altera-struct-pack-align

Finds structs that are inefficiently packed or aligned, and recommends
packing and/or aligning of said structs as needed.

Structs that are not packed take up more space than they should, and accessing
structs that are not well aligned is inefficient.

Fix-its are provided to fix both of these issues by inserting and/or amending
relevant struct attributes.

Based on the `Altera SDK for OpenCL: Best Practices Guide
<https://www.altera.com/en_US/pdfs/literature/hb/opencl-sdk/aocl_optimization_guide.pdf>`_.

```c++
// The following struct is originally aligned to 4 bytes, and thus takes up
// 12 bytes of memory instead of 10. Packing the struct will make it use
// only 10 bytes of memory, and aligning it to 16 bytes will make it
// efficient to access.
struct example {
 char a; // 1 byte
 double b; // 8 bytes
 char c; // 1 byte
};

// The following struct is arranged in such a way that packing is not needed.
// However, it is aligned to 4 bytes instead of 8, and thus needs to be
// explicitly aligned.
struct implicitly_packed_example {
 char a; // 1 byte
 char b; // 1 byte
 char c; // 1 byte
 char d; // 1 byte
 int e; // 4 bytes
};

// The following struct is explicitly aligned and packed.
struct good_example {
 char a; // 1 byte
 double b; // 8 bytes
 char c; // 1 byte
} __attribute__((packed)) __attribute__((aligned(16));

// Explicitly aligning a struct to the wrong value will result in a warning.
// The following example should be aligned to 16 bytes, not 32.
struct badly_aligned_example {
 char a; // 1 byte
 double b; // 8 bytes
 char c; // 1 byte
} __attribute__((packed)) __attribute__((aligned(32)));
```

**Title:** clang-tidy - altera-unroll-loops

# altera-unroll-loops

Finds inner loops that have not been unrolled, as well as fully unrolled loops
with unknown loop bounds or a large number of iterations.

Unrolling inner loops could improve the performance of OpenCL kernels. However,
if they have unknown loop bounds or a large number of iterations, they cannot
be fully unrolled, and should be partially unrolled.

Notes:

- This check is unable to determine the number of iterations in a `while` or
 `do..while` loop; hence if such a loop is fully unrolled, a note is emitted
 advising the user to partially unroll instead.

- In `for` loops, our check only works with simple arithmetic increments (
 `+`, `-`, `*`, `/`). For all other increments, partial unrolling is
 advised.

- Depending on the exit condition, the calculations for determining if the
 number of iterations is large may be off by 1. This should not be an issue
 since the cut-off is generally arbitrary.

Based on the `Altera SDK for OpenCL: Best Practices Guide
<https://www.altera.com/en_US/pdfs/literature/hb/opencl-sdk/aocl_optimization_guide.pdf>`_.

```c++
for (int i = 0; i < 10; i++) { // ok: outer loops should not be unrolled
 int j = 0;
 do { // warning: this inner do..while loop should be unrolled
 j++;
 } while (j < 15);

 int k = 0;
 #pragma unroll
 while (k < 20) { // ok: this inner loop is already unrolled
 k++;
 }
}

int A[1000];
#pragma unroll
// warning: this loop is large and should be partially unrolled
for (int a : A) {
 printf("%d", a);
}

#pragma unroll 5
// ok: this loop is large, but is partially unrolled
for (int a : A) {
 printf("%d", a);
}

#pragma unroll
// warning: this loop is large and should be partially unrolled
for (int i = 0; i < 1000; ++i) {
 printf("%d", i);
}

#pragma unroll 5
// ok: this loop is large, but is partially unrolled
for (int i = 0; i < 1000; ++i) {
 printf("%d", i);
}

#pragma unroll
// warning: << operator not supported, recommend partial unrolling
for (int i = 0; i < 1000; i<<1) {
 printf("%d", i);
}

std::vector<int> someVector (100, 0);
int i = 0;
#pragma unroll
// note: loop may be large, recommend partial unrolling
while (i < someVector.size()) {
 someVector[i]++;
}

#pragma unroll
// note: loop may be large, recommend partial unrolling
while (true) {
 printf("In loop");
}

#pragma unroll 5
// ok: loop may be large, but is partially unrolled
while (i < someVector.size()) {
 someVector[i]++;
}
```

## Options

**Option:** MaxLoopIterations
Defines the maximum number of loop iterations that a fully unrolled loop
can have. By default, it is set to `100`.

In practice, this refers to the integer value of the upper bound
within the loop statement's condition expression.

```{title} clang-tidy - android-cloexec-accept
```

# android-cloexec-accept

The usage of `accept()` is not recommended, it's better to use `accept4()`.
Without this flag, an opened sensitive file descriptor would remain open across
a fork+exec to a lower-privileged SELinux domain.

Examples:

```cpp
accept(sockfd, addr, addrlen);

// becomes

accept4(sockfd, addr, addrlen, SOCK_CLOEXEC);
```

```{title} clang-tidy - android-cloexec-accept4
```

# android-cloexec-accept4

`accept4()` should include `SOCK_CLOEXEC` in its type argument to avoid the
file descriptor leakage. Without this flag, an opened sensitive file would
remain open across a fork+exec to a lower-privileged SELinux domain.

Examples:

```cpp
accept4(sockfd, addr, addrlen, SOCK_NONBLOCK);

// becomes

accept4(sockfd, addr, addrlen, SOCK_NONBLOCK | SOCK_CLOEXEC);
```

```{title} clang-tidy - android-cloexec-creat
```

# android-cloexec-creat

The usage of `creat()` is not recommended, it's better to use `open()`.

Examples:

```cpp
int fd = creat(path, mode);

// becomes

int fd = open(path, O_WRONLY | O_CREAT | O_TRUNC | O_CLOEXEC, mode);
```

```{title} clang-tidy - android-cloexec-dup
```

# android-cloexec-dup

The usage of `dup()` is not recommended, it's better to use `fcntl()`,
which can set the close-on-exec flag. Otherwise, an opened sensitive file would
remain open across a fork+exec to a lower-privileged SELinux domain.

Examples:

```cpp
int fd = dup(oldfd);

// becomes

int fd = fcntl(oldfd, F_DUPFD_CLOEXEC);
```

```{title} clang-tidy - android-cloexec-epoll-create
```

# android-cloexec-epoll-create

The usage of `epoll_create()` is not recommended, it's better to use
`epoll_create1()`, which allows close-on-exec.

Examples:

```cpp
epoll_create(size);

// becomes

epoll_create1(EPOLL_CLOEXEC);
```

```{title} clang-tidy - android-cloexec-epoll-create1
```

# android-cloexec-epoll-create1

`epoll_create1()` should include `EPOLL_CLOEXEC` in its type argument to
avoid the file descriptor leakage. Without this flag, an opened sensitive file
would remain open across a fork+exec to a lower-privileged SELinux domain.

Examples:

```cpp
epoll_create1(0);

// becomes

epoll_create1(EPOLL_CLOEXEC);
```

```{title} clang-tidy - android-cloexec-fopen
```

# android-cloexec-fopen

`fopen()` should include `e` in their mode string; so `re` would be
valid. This is equivalent to having set `FD_CLOEXEC` on that descriptor.

Examples:

```cpp
fopen("fn", "r");

// becomes

fopen("fn", "re");
```

```{title} clang-tidy - android-cloexec-inotify-init
```

# android-cloexec-inotify-init

The usage of `inotify_init()` is not recommended, it's better to use
`inotify_init1()`.

Examples:

```cpp
inotify_init();

// becomes

inotify_init1(IN_CLOEXEC);
```

```{title} clang-tidy - android-cloexec-inotify-init1
```

# android-cloexec-inotify-init1

`inotify_init1()` should include `IN_CLOEXEC` in its type argument
to avoid the file descriptor leakage. Without this flag, an opened
sensitive file would remain open across a fork+exec to a
lower-privileged SELinux domain.

Examples:

```cpp
inotify_init1(IN_NONBLOCK);

// becomes

inotify_init1(IN_NONBLOCK | IN_CLOEXEC);
```

```{title} clang-tidy - android-cloexec-memfd-create
```

# android-cloexec-memfd-create

`memfd_create()` should include `MFD_CLOEXEC` in its type argument to avoid
the file descriptor leakage. Without this flag, an opened sensitive file would
remain open across a fork+exec to a lower-privileged SELinux domain.

Examples:

```cpp
memfd_create(name, MFD_ALLOW_SEALING);

// becomes

memfd_create(name, MFD_ALLOW_SEALING | MFD_CLOEXEC);
```

```{title} clang-tidy - android-cloexec-open
```

# android-cloexec-open

A common source of security bugs is code that opens a file without using the
`O_CLOEXEC` flag. Without that flag, an opened sensitive file would remain
open across a fork+exec to a lower-privileged SELinux domain, leaking that
sensitive data. Open-like functions including `open()`, `openat()`, and
`open64()` should include `O_CLOEXEC` in their flags argument.

Examples:

```cpp
open("filename", O_RDWR);
open64("filename", O_RDWR);
openat(0, "filename", O_RDWR);

// becomes

open("filename", O_RDWR | O_CLOEXEC);
open64("filename", O_RDWR | O_CLOEXEC);
openat(0, "filename", O_RDWR | O_CLOEXEC);
```

```{title} clang-tidy - android-cloexec-pipe
```

# android-cloexec-pipe

This check detects usage of `pipe()`. Using `pipe()` is not recommended,
`pipe2()` is the suggested replacement. The check also adds the `O_CLOEXEC`
flag that marks the file descriptor to be closed in child processes.
Without this flag a sensitive file descriptor can be leaked to a
child process, potentially into a lower-privileged SELinux domain.

Examples:

```cpp
pipe(pipefd);
```

Suggested replacement:

```cpp
pipe2(pipefd, O_CLOEXEC);
```

```{title} clang-tidy - android-cloexec-pipe2
```

# android-cloexec-pipe2

This check ensures that `pipe2()` is called with the `O_CLOEXEC` flag.
The check also adds the `O_CLOEXEC` flag that marks the file descriptor
to be closed in child processes.
Without this flag a sensitive file descriptor can be leaked to a child process,
potentially into a lower-privileged SELinux domain.

Examples:

```cpp
pipe2(pipefd, O_NONBLOCK);
```

Suggested replacement:

```cpp
pipe2(pipefd, O_NONBLOCK | O_CLOEXEC);
```

```{title} clang-tidy - android-cloexec-socket
```

# android-cloexec-socket

`socket()` should include `SOCK_CLOEXEC` in its type argument to avoid the
file descriptor leakage. Without this flag, an opened sensitive file would
remain open across a fork+exec to a lower-privileged SELinux domain.

Examples:

```cpp
socket(domain, type, SOCK_STREAM);

// becomes

socket(domain, type, SOCK_STREAM | SOCK_CLOEXEC);
```

**Title:** clang-tidy - android-comparison-in-temp-failure-retry

# android-comparison-in-temp-failure-retry

Diagnoses comparisons that appear to be incorrectly placed in the argument to
the `TEMP_FAILURE_RETRY` macro. Having such a use is incorrect in the vast
majority of cases, and will often silently defeat the purpose of the
`TEMP_FAILURE_RETRY` macro.

For context, `TEMP_FAILURE_RETRY` is `a convenience macro
<https://www.gnu.org/software/libc/manual/html_node/Interrupted-Primitives.html>`_
provided by both glibc and Bionic. Its purpose is to repeatedly run a syscall
until it either succeeds, or fails for reasons other than being interrupted.

Example buggy usage looks like:

```c
char cs[1];
while (TEMP_FAILURE_RETRY(read(STDIN_FILENO, cs, sizeof(cs)) != 0)) {
 // Do something with cs.
}
```

Because `TEMP_FAILURE_RETRY` will check for whether the result
*of the comparison* is `-1`, and retry if so.

If you encounter this, the fix is simple: lift the comparison out of the
`TEMP_FAILURE_RETRY` argument, like so:

```c
char cs[1];
while (TEMP_FAILURE_RETRY(read(STDIN_FILENO, cs, sizeof(cs))) != 0) {
 // Do something with cs.
}
```

## Options

**Option:** RetryMacros
A comma-separated list of the names of retry macros to be checked.
Default is `TEMP_FAILURE_RETRY`.

**Title:** clang-tidy - boost-use-ranges

# boost-use-ranges

Detects calls to standard library iterator algorithms that could be replaced
with a Boost ranges version instead.

## Example

```c++
auto Iter1 = std::find(Items.begin(), Items.end(), 0);
auto AreSame = std::equal(Items1.cbegin(), Items1.cend(), std::begin(Items2),
 std::end(Items2));
```

Transforms to:

```c++
auto Iter1 = boost::range::find(Items, 0);
auto AreSame = boost::range::equal(Items1, Items2);
```

## Supported algorithms

Calls to the following std library algorithms are checked:

`std::accumulate`,
`std::adjacent_difference`,
`std::adjacent_find`,
`std::all_of`,
`std::any_of`,
`std::binary_search`,
`std::copy_backward`,
`std::copy_if`,
`std::copy`,
`std::count_if`,
`std::count`,
`std::equal_range`,
`std::equal`,
`std::fill`,
`std::find_end`,
`std::find_first_of`,
`std::find_if_not`,
`std::find_if`,
`std::find`,
`std::for_each`,
`std::generate`,
`std::includes`,
`std::iota`,
`std::is_partitioned`,
`std::is_permutation`,
`std::is_sorted_until`,
`std::is_sorted`,
`std::lexicographical_compare`,
`std::lower_bound`,
`std::make_heap`,
`std::max_element`,
`std::merge`,
`std::min_element`,
`std::mismatch`,
`std::next_permutation`,
`std::none_of`,
`std::partial_sum`,
`std::partial_sort_copy`,
`std::partition_copy`,
`std::partition_point`,
`std::partition`,
`std::pop_heap`,
`std::prev_permutation`,
`std::push_heap`,
`std::random_shuffle`,
`std::reduce`,
`std::remove_copy_if`,
`std::remove_copy`,
`std::remove_if`,
`std::remove`,
`std::replace_copy_if`,
`std::replace_copy`,
`std::replace_if`,
`std::replace`,
`std::reverse_copy`,
`std::reverse`,
`std::search`,
`std::set_difference`,
`std::set_intersection`,
`std::set_symmetric_difference`,
`std::set_union`,
`std::sort_heap`,
`std::sort`,
`std::stable_partition`,
`std::stable_sort`,
`std::transform`,
`std::unique_copy`,
`std::unique`,
`std::upper_bound`.

The check will also look for the following functions from the
`boost::algorithm` namespace:

`all_of_equal`,
`any_of_equal`,
`any_of`,
`apply_permutation`,
`apply_reverse_permutation`,
`clamp_range`,
`copy_if_until`,
`copy_if_while`,
`copy_if`,
`copy_until`,
`copy_while`,
`find_backward`,
`find_if_backward`,
`find_if_not_backward`,
`find_if_not`,
`find_not_backward`,
`hex_lower`,
`hex`,
`iota`, `all_of`,
`is_decreasing`,
`is_increasing`,
`is_palindrome`,
`is_partitioned_until`,
`is_partitioned`,
`is_permutation`,
`is_sorted_until`,
`is_sorted`,
`is_strictly_decreasing`,
`is_strictly_increasing`,
`none_of_equal`,
`none_of`,
`one_of_equal`,
`one_of`,
`partition_copy`,
`partition_point`,
`reduce`,
`unhex`.

## Reverse Iteration

If calls are made using reverse iterators on containers, The code will be
fixed using the `boost::adaptors::reverse` adaptor.

```c++
auto AreSame = std::equal(Items1.rbegin(), Items1.rend(),
 std::crbegin(Items2), std::crend(Items2));
```

Transforms to:

```c++
auto AreSame = boost::range::equal(boost::adaptors::reverse(Items1),
 boost::adaptors::reverse(Items2));
```

## Options

**Option:** IncludeStyle
A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.
**Option:** IncludeBoostSystem
If `true` (default value) the boost headers are included as system headers
with angle brackets (`#include <boost.hpp>`), otherwise quotes are used
(`#include "boost.hpp"`).
**Option:** UseReversePipe
When `true` (default `false`), fixes which involve reverse ranges will use the
pipe adaptor syntax instead of the function syntax.

.. code-block:: c++

 std::find(Items.rbegin(), Items.rend(), 0);

Transforms to:

.. code-block:: c++

 boost::range::find(Items | boost::adaptors::reversed, 0);

```{title} clang-tidy - boost-use-to-string
```

# boost-use-to-string

This check finds conversion from integer type like `int` to
`std::string` or `std::wstring` using `boost::lexical_cast`,
and replace it with calls to `std::to_string` and `std::to_wstring`.

It doesn't replace conversion from floating points despite the `to_string`
overloads, because it would change the behavior.

```cpp
auto str = boost::lexical_cast<std::string>(42);
auto wstr = boost::lexical_cast<std::wstring>(2137LL);

// Will be changed to
auto str = std::to_string(42);
auto wstr = std::to_wstring(2137LL);
```

```{title} clang-tidy - bugprone-argument-comment
```

# bugprone-argument-comment

Checks that argument comments match parameter names and can optionally add
missing comments for literals, init-lists, and constructed temporaries.

The check understands argument comments in the form `/*parameter_name=*/`
that are placed right before the argument.

```c++
void f(bool foo);

...

f(/*bar=*/true);
// warning: argument name 'bar' in comment does not match parameter name 'foo'
```

The check tries to detect typos and suggest automated fixes for them. It can
also insert missing comments for configured argument kinds.

## Options

```{option} StrictMode
When `false`, the check will ignore leading and trailing
underscores and case when comparing names -- otherwise they are taken into
account. Default is `false`.
```

```{option} IgnoreSingleArgument
When `true`, the check will ignore the single argument. Default is `false`.
```

```{option} CommentAnonymousInitLists
When `true`, the check will add argument comments in the format
`/*ParameterName=*/` right before anonymous braced-init list arguments
such as `{}` and `{1, 2, 3}`. Default is `false`.
```

Before:

```c++
void foo(const std::vector<int> &Dims);

foo({});
```

After:

```c++
void foo(const std::vector<int> &Dims);

foo(/*Dims=*/{});
```

```{option} CommentBoolLiterals
When `true`, the check will add argument comments in the format
`/*ParameterName=*/` right before the boolean literal argument.
Default is `false`.
```

Before:

```c++
void foo(bool TurnKey, bool PressButton);

foo(true, false);
```

After:

```c++
void foo(bool TurnKey, bool PressButton);

foo(/*TurnKey=*/true, /*PressButton=*/false);
```

```{option} CommentCharacterLiterals
When `true`, the check will add argument comments in the format
`/*ParameterName=*/` right before the character literal argument.
Default is `false`.
```

Before:

```c++
void foo(char *Character);

foo('A');
```

After:

```c++
void foo(char *Character);

foo(/*Character=*/'A');
```

```{option} CommentFloatLiterals
When `true`, the check will add argument comments in the format
`/*ParameterName=*/` right before the float/double literal argument.
Default is `false`.
```

Before:

```c++
void foo(float Pi);

foo(3.14159);
```

After:

```c++
void foo(float Pi);

foo(/*Pi=*/3.14159);
```

```{option} CommentIntegerLiterals
When `true`, the check will add argument comments in the format
`/*ParameterName=*/` right before the integer literal argument.
Default is `false`.
```

Before:

```c++
void foo(int MeaningOfLife);

foo(42);
```

After:

```c++
void foo(int MeaningOfLife);

foo(/*MeaningOfLife=*/42);
```

```{option} CommentNullPtrs
When `true`, the check will add argument comments in the format
`/*ParameterName=*/` right before the nullptr literal argument.
Default is `false`.
```

Before:

```c++
void foo(A* Value);

foo(nullptr);
```

After:

```c++
void foo(A* Value);

foo(/*Value=*/nullptr);
```

```{option} CommentParenthesizedTemporaries
When `true`, the check will add argument comments in the format
`/*ParameterName=*/` right before explicit temporary constructions such as
`Type()` and `Type(1, 2, 3)`. Default is `false`.
```

Before:

```c++
struct Dims {
 Dims();
 Dims(int, int, int);
};

void foo(const Dims &DimsValue);

foo(Dims());
foo(Dims(1, 2, 3));
```

After:

```c++
struct Dims {
 Dims();
 Dims(int, int, int);
};

void foo(const Dims &DimsValue);

foo(/*DimsValue=*/Dims());
foo(/*DimsValue=*/Dims(1, 2, 3));
```

```{option} CommentStringLiterals
When `true`, the check will add argument comments in the format
`/*ParameterName=*/` right before the string literal argument.
Default is `false`.
```

Before:

```c++
void foo(const char *String);
void foo(const wchar_t *WideString);

foo("Hello World");
foo(L"Hello World");
```

After:

```c++
void foo(const char *String);
void foo(const wchar_t *WideString);

foo(/*String=*/"Hello World");
foo(/*WideString=*/L"Hello World");
```

```{option} CommentTypedInitLists
When `true`, the check will add argument comments in the format
`/*ParameterName=*/` right before typed braced-init list arguments such
as `Type{}`. Default is `false`.
```

Before:

```c++
void foo(const std::vector<int> &Dims);

foo(std::vector<int>{});
```

After:

```c++
void foo(const std::vector<int> &Dims);

foo(/*Dims=*/std::vector<int>{});
```

```{option} CommentUserDefinedLiterals
When `true`, the check will add argument comments in the format
`/*ParameterName=*/` right before the user defined literal argument.
Default is `false`.
```

Before:

```c++
void foo(double Distance);

double operator"" _km(long double);

foo(402.0_km);
```

After:

```c++
void foo(double Distance);

double operator"" _km(long double);

foo(/*Distance=*/402.0_km);
```

```{title} clang-tidy - bugprone-assert-side-effect
```

# bugprone-assert-side-effect

Finds `assert()` with side effect.

The condition of `assert()` is evaluated only in debug builds so a
condition with side effect can cause different behavior in debug / release
builds.

## Options

```{option} AssertMacros
A comma-separated list of the names of assert macros to be checked.
Default is `assert,NSAssert,NSCAssert`.
```

```{option} CheckFunctionCalls
Whether to treat non-const member and non-member functions as they produce
side effects. Disabled by default because it can increase the number of false
positive warnings.
```

```{option} IgnoredFunctions
A semicolon-separated list of the names of functions or methods to be
considered as not having side-effects. Regular expressions are accepted,
e.g. `[Rr]ef(erence)?$` matches every type with suffix `Ref`, `ref`,
`Reference` and `reference`. The default is empty. If a name in the list
contains the sequence `::` it is matched against the qualified type name
(i.e. `namespace::Type`), otherwise it is matched against only
the type name (i.e. `Type`).
```

```{title} clang-tidy - bugprone-assignment-in-if-condition
```

# bugprone-assignment-in-if-condition

Finds assignments within conditions of `if` statements.
Such assignments are bug-prone because they may have been intended as
equality tests.

This check finds all assignments within `if` conditions, including ones that
are not flagged by `-Wparentheses` due to an extra set of parentheses, and
including assignments that call an overloaded `operator=()`. The identified
assignments violate
[BARR group "Rule 8.2.c"](https://barrgroup.com/embedded-systems/books/embedded-c-coding-standard/statement-rules/if-else-statements).

```cpp
int f = 3;
if(f = 4) { // This is identified by both `Wparentheses` and this check - should it have been: `if (f == 4)` ?
 f = f + 1;
}

if((f == 5) || (f = 6)) { // the assignment here `(f = 6)` is identified by this check, but not by `-Wparentheses`. Should it have been `(f == 6)` ?
 f = f + 2;
}
```

```{title} clang-tidy - bugprone-assignment-in-selection-statement
```

# bugprone-assignment-in-selection-statement

Finds assignments within selection statements.
Such assignments may indicate programmer error because they may have been
intended as equality tests. The selection statements are conditions of `if`
and loop (`for`, `while`, `do`) statements, condition of conditional
operator (`?:`) and any operand of a binary logical operator (`&&`,
`||`). The check finds assignments within these contexts if the single
expression is an assignment or the assignment is contained (recursively) in
last operand of a comma (`,`) operator or true and false expressions in a
conditional operator. The warning is suppressed if the assignment is placed in
extra parentheses, but only if the assignment is the single expression of a
condition (of `if` or a loop statement).

This check corresponds to the CERT rule
[EXP45-C. Do not perform assignments in selection statements](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/expressions-exp/exp45-c/).

# Examples

The check emits a warning in the following cases at the indicated locations:

```c++
int x = 3;

if (x = 4) // should it be `x == 4` instead of 'x = 4' ?
 x = x + 1;

while ((x <= 11) || (x = 22)) // assignment appears as operand of a logical operator
 x += 2;

do {
 x += 5;
} while ((x > 10) ? (x = 11) : (x > 5)); // assignment in loop condition (from `x = 11`)

for (int i = 0; i == 2, x = 5; ++i) // assignment in loop condition (from last operand of comma)
 foo1(i, x);

for (int i = 0; i == 2, (x = 5); ++i) // assignment is not a single expression, parentheses do not prevent the warning
 foo1(i, x);

int a = (x == 2) || (x = 3); // assignment appears in the operand a logical operator
```

The following cases do not produce a warning:

```c++
if ((x = 1)) { // a single assignment between parentheses
 x += 10;

if ((x = 1) != 0) { // assignment appears in a complex expression and without a logical operator
 ++x;

if (foo(x = 9) && array[x = 8]) { // assignment appears in argument of function call or array index
 ++x;

for (int i = 0; i = 2, x == 5; ++i) // assignment does not take part in the condition of the loop
 foo1(i, x);
```

```{title} clang-tidy - bugprone-bad-signal-to-kill-thread
```

# bugprone-bad-signal-to-kill-thread

Finds `pthread_kill` function calls when a thread is terminated by
raising `SIGTERM` signal and the signal kills the entire process, not
just the individual thread. Use any signal except `SIGTERM`.

```cpp
pthread_kill(thread, SIGTERM);
```

This check corresponds to the CERT C Coding Standard rule
[POS44-C. Do not use signals to terminate threads](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/posix-pos/pos44-c/).

`cert-pos44-c` redirects here as an alias of this check.

```{title} clang-tidy - bugprone-bitwise-pointer-cast
```

# bugprone-bitwise-pointer-cast

Warns about code that tries to cast between pointers by means of
`std::bit_cast` or `memcpy`.

The motivation is that `std::bit_cast` is advertised as the safe alternative
to type punning via `reinterpret_cast` in modern C++. However, one should not
blindly replace `reinterpret_cast` with `std::bit_cast`, as follows:

```c++
int x{};
-float y = *reinterpret_cast<float*>(&x);
+float y = *std::bit_cast<float*>(&x);
```

The drop-in replacement behaves exactly the same as `reinterpret_cast`, and
Undefined Behavior is still invoked. `std::bit_cast` is copying the bytes of
the input pointer, not the pointee, into an output pointer of a different type,
which may violate the strict aliasing rules. However, simply looking at the
code, it looks "safe", because it uses `std::bit_cast` which is advertised as
safe.

The solution to safe type punning is to apply `std::bit_cast` on value types,
not on pointer types:

```c++
int x{};
float y = std::bit_cast<float>(x);
```

This way, the bytes of the input object are copied into the output object,
which is much safer. Do note that Undefined Behavior can still occur, if there
is no value of type `To` corresponding to the value representation produced.
Compilers may be able to optimize this copy and generate identical assembly to
the original `reinterpret_cast` version.

Code before C++20 may backport `std::bit_cast` by means of `memcpy`, or
simply call `memcpy` directly, which is equally problematic. This is also
detected by this check:

```c++
int* x{};
float* y{};
std::memcpy(&y, &x, sizeof(x));
```

Alternatively, if a cast between pointers is truly wanted, `reinterpret_cast`
should be used, to clearly convey the intent and enable warnings from compilers
and linters, which should be addressed accordingly.

```{title} clang-tidy - bugprone-bool-pointer-implicit-conversion
```

# bugprone-bool-pointer-implicit-conversion

Checks for conditions based on implicit conversion from a `bool` pointer to
`bool`.

Example:

```cpp
bool *p;
if (p) {
 // Never used in a pointer-specific way.
}
```

```{title} clang-tidy - bugprone-branch-clone
```

# bugprone-branch-clone

Checks for repeated branches in `if/else if/else` chains, consecutive
repeated branches in `switch` statements and identical true and false
branches in conditional operators.

```c++
if (test_value(x)) {
 y++;
 do_something(x, y);
} else {
 y++;
 do_something(x, y);
}
```

In this simple example (which could arise e.g. as a copy-paste error) the
`then` and `else` branches are identical and the code is equivalent the
following shorter and cleaner code:

```c++
test_value(x); // can be omitted unless it has side effects
y++;
do_something(x, y);
```

If this is the intended behavior, then there is no reason to use a conditional
statement; otherwise the issue can be solved by fixing the branch that is
handled incorrectly.

The check detects repeated branches in longer `if/else if/else` chains
where it would be even harder to notice the problem.

The check also detects repeated inner and outer `if` statements that may
be a result of a copy-paste error. This check cannot currently detect
identical inner and outer `if` statements if code is between the `if`
conditions. An example is as follows.

```c++
void test_warn_inner_if_1(int x) {
 if (x == 1) { // warns, if with identical inner if
 if (x == 1) // inner if is here
 ;
 if (x == 1) { // does not warn, cannot detect
 int y = x;
 if (x == 1)
 ;
 }
}
```

In `switch` statements the check only reports repeated branches when they are
consecutive, because it is relatively common that the `case:` labels have
some natural ordering and rearranging them would decrease the readability of
the code. For example:

```c++
switch (ch) {
case 'a':
 return 10;
case 'A':
 return 10;
case 'b':
 return 11;
case 'B':
 return 11;
default:
 return 10;
}
```

Here the check reports that the `'a'` and `'A'` branches are identical
(and that the `'b'` and `'B'` branches are also identical), but does not
report that the `default:` branch is also identical to the first two branches.
If this is indeed the correct behavior, then it could be implemented as:

```c++
switch (ch) {
case 'a':
case 'A':
 return 10;
case 'b':
case 'B':
 return 11;
default:
 return 10;
}
```

Here the check does not warn for the repeated `return 10;`, which is good if
we want to preserve that `'a'` is before `'b'` and `default:` is the last
branch.

Switch cases marked with the `[[fallthrough]]` attribute are ignored.

Finally, the check also examines conditional operators and reports code like:

```c++
return test_value(x) ? x : x;
```

Unlike if statements, the check does not detect chains of conditional
operators.

Note: This check also reports situations where branches become identical only
after preprocessing.

```{title} clang-tidy - bugprone-capturing-this-in-member-variable
```

# bugprone-capturing-this-in-member-variable

Finds lambda captures that capture the `this` pointer and store it as class
members without handle the copy and move constructors and the assignments.

Capture this in a lambda and store it as a class member is dangerous because
the lambda can outlive the object it captures. Especially when the object is
copied or moved, the captured `this` pointer will be implicitly propagated
to the new object. Most of the time, people will believe that the captured
`this` pointer points to the new object, which will lead to bugs.

```c++
struct C {
 C() : Captured([this]() -> C const * { return this; }) {}
 std::function<C const *()> Captured;
};

void foo() {
 C v1{};
 C v2 = v1; // v2.Captured capture v1's 'this' pointer
 assert(v2.Captured() == v1.Captured()); // v2.Captured capture v1's 'this' pointer
 assert(v2.Captured() == &v2); // assertion failed.
}
```

Possible fixes:

- marking copy and move constructors and assignment operators deleted.
- using class member method instead of class member variable with function
 object types.
- passing `this` pointer as parameter.

## Options

```{option} FunctionWrapperTypes
A semicolon-separated list of names of types. Used to specify function
wrapper that can hold lambda expressions.
Default is `::std::function;::std::move_only_function;::boost::function`.
```

```{option} BindFunctions
A semicolon-separated list of fully qualified names of functions that can
capture `this` pointer.
Default is `::std::bind;::boost::bind;::std::bind_front;::std::bind_back;
::boost::compat::bind_front;::boost::compat::bind_back`.
```

```{title} clang-tidy - bugprone-casting-through-void
```

# bugprone-casting-through-void

Detects unsafe or redundant two-step casting operations involving `void*`,
which is equivalent to `reinterpret_cast` as per the
[C++ Standard](https://eel.is/c++draft/expr.reinterpret.cast#7).

Two-step type conversions via `void*` are discouraged for several reasons.

- They obscure code and impede its understandability, complicating maintenance.
- These conversions bypass valuable compiler support, erasing warnings related
 to pointer alignment. It may violate strict aliasing rule and leading to
 undefined behavior.
- In scenarios involving multiple inheritance, ambiguity and unexpected
 outcomes can arise due to the loss of type information, posing runtime
 issues.

In summary, avoiding two-step type conversions through `void*` ensures
clearer code, maintains essential compiler warnings, and prevents ambiguity
and potential runtime errors, particularly in complex inheritance scenarios.
If such a cast is wanted, it shall be done via `reinterpret_cast`,
to express the intent more clearly.

Note: it is expected that, after applying the suggested fix and using
`reinterpret_cast`, the check
{doc}`cppcoreguidelines-pro-type-reinterpret-cast
<../cppcoreguidelines/pro-type-reinterpret-cast>` will emit a warning.
This is intentional: `reinterpret_cast` is a dangerous operation that can
easily break the strict aliasing rules when dereferencing the casted pointer,
invoking Undefined Behavior. The warning is there to prompt users to carefully
analyze whether the usage of `reinterpret_cast` is safe, in which case the
warning may be suppressed.

Examples:

```c++
using IntegerPointer = int *;
double *ptr;

static_cast<IntegerPointer>(static_cast<void *>(ptr)); // WRONG
reinterpret_cast<IntegerPointer>(reinterpret_cast<void *>(ptr)); // WRONG
(IntegerPointer)(void *)ptr; // WRONG
IntegerPointer(static_cast<void *>(ptr)); // WRONG

reinterpret_cast<IntegerPointer>(ptr); // OK, clearly expresses intent.
 // NOTE: dereferencing this pointer violates
 // the strict aliasing rules, invoking
 // Undefined Behavior.
```

```{title} clang-tidy - bugprone-chained-comparison
```

# bugprone-chained-comparison

Check detects chained comparison operators that can lead to unintended
behavior or logical errors.

Chained comparisons are expressions that use multiple comparison operators
to compare three or more values. For example, the expression `a < b < c`
compares the values of `a`, `b`, and `c`. However, this expression does
not evaluate as `(a < b) && (b < c)`, which is probably what the developer
intended. Instead, it evaluates as `(a < b) < c`, which may produce
unintended results, especially when the types of `a`, `b`, and `c` are
different.

To avoid such errors, the check will issue a warning when a chained
comparison operator is detected, suggesting to use parentheses to specify
the order of evaluation or to use a logical operator to separate comparison
expressions.

Consider the following examples:

```c++
int a = 2, b = 6, c = 4;
if (a < b < c) {
 // This block will be executed
}
```

In this example, the developer intended to check if `a` is less than `b`
and `b` is less than `c`. However, the expression `a < b < c` is
equivalent to `(a < b) < c`. Since `a < b` is `true`, the expression
`(a < b) < c` is evaluated as `1 < c`, which is equivalent to `true < c`
and is invalid in this case as `b < c` is `false`.

Even that above issue could be detected as comparison of `int` to `bool`,
there is more dangerous example:

```c++
bool a = false, b = false, c = true;
if (a == b == c) {
 // This block will be executed
}
```

In this example, the developer intended to check if `a`, `b`, and `c` are
all equal. However, the expression `a == b == c` is evaluated as
`(a == b) == c`. Since `a == b` is true, the expression `(a == b) == c`
is evaluated as `true == c`, which is equivalent to `true == true`.
This comparison yields `true`, even though `a` and `b` are `false`, and
are not equal to `c`.

To avoid this issue, the developer can use a logical operator to separate the
comparison expressions, like this:

```c++
if (a == b && b == c) {
 // This block will not be executed
}
```

Alternatively, use of parentheses in the comparison expressions can make the
developer's intention more explicit and help avoid misunderstanding.

```c++
if ((a == b) == c) {
 // This block will be executed
}
```

## Options

```{option} IgnoreMacros
If `true`, the check will not warn on chained comparisons inside macros.
Default is `false`.
```

```{title} clang-tidy - bugprone-command-processor
```

# bugprone-command-processor

Flags calls to `system()`, `popen()`, and `_popen()`, which
execute a command processor. It does not flag calls to `system()` with a null
pointer argument, as such a call checks for the presence of a command processor
but does not actually attempt to execute a command.

## References

This check corresponds to the CERT C Coding Standard rule
[ENV33-C. Do not call system()](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/environment-env/env33-c/).

```{title} clang-tidy - bugprone-compare-pointer-to-member-virtual-function
```

# bugprone-compare-pointer-to-member-virtual-function

Detects unspecified behavior about equality comparison between pointer to
member virtual function and anything other than null-pointer-constant.

```c++
struct A {
 void f1();
 void f2();
 virtual void f3();
 virtual void f4();

 void g1(int);
};

void fn() {
 bool r1 = (&A::f1 == &A::f2); // ok
 bool r2 = (&A::f1 == &A::f3); // bugprone
 bool r3 = (&A::f1 != &A::f3); // bugprone
 bool r4 = (&A::f3 == nullptr); // ok
 bool r5 = (&A::f3 == &A::f4); // bugprone

 void (A::*v1)() = &A::f3;
 bool r6 = (v1 == &A::f1); // bugprone
 bool r6 = (v1 == nullptr); // ok

 void (A::*v2)() = &A::f2;
 bool r7 = (v2 == &A::f1); // false positive, but potential risk if assigning other value to v2.

 void (A::*v3)(int) = &A::g1;
 bool r8 = (v3 == &A::g1); // ok, no virtual function match void(A::*)(int) signature.
}
```

Provide warnings on equality comparisons involve pointers to member virtual
function or variables which is potential pointer to member virtual function and
any entity other than a null-pointer constant.

In certain compilers, virtual function addresses are not conventional pointers
but instead consist of offsets and indexes within a virtual function table
(vtable). Consequently, these pointers may vary between base and derived
classes, leading to unpredictable behavior when compared directly. This issue
becomes particularly challenging when dealing with pointers to pure virtual
functions, as they may not even have a valid address, further complicating
comparisons.

Instead, it is recommended to utilize the `typeid` operator or other
appropriate mechanisms for comparing objects to ensure robust and predictable
behavior in your codebase. By heeding this detection and adopting a more reliable
comparison method, you can mitigate potential issues related to unspecified
behavior, especially when dealing with pointers to member virtual functions or pure
virtual functions, thereby improving the overall stability and maintainability
of your code. In scenarios involving pointers to member virtual functions, it's
only advisable to employ `nullptr` for comparisons.

## Limitations

Does not analyze values stored in a variable. For variable, only analyze all
virtual methods in the same `class` or `struct` and diagnose when assigning
a pointer to member virtual function to this variable is possible.

```{title} clang-tidy - bugprone-copy-constructor-init
```

# bugprone-copy-constructor-init

Finds copy constructors where the constructor doesn't call the copy constructor
of the base class.

```c++
class Copyable {
public:
 Copyable() = default;
 Copyable(const Copyable &) = default;

 int memberToBeCopied = 0;
};

class X2 : public Copyable {
 X2(const X2 &other) {} // Copyable(other) is missing
};
```

Also finds copy constructors where the constructor of
the base class don't have parameter.

```c++
class X3 : public Copyable {
 X3(const X3 &other) : Copyable() {} // other is missing
};
```

Failure to properly initialize base class sub-objects during copy construction
can result in undefined behavior, crashes, data corruption, or other unexpected
outcomes. The check ensures that the copy constructor of a derived class
properly calls the copy constructor of the base class, helping to prevent bugs
and improve code quality.

## Limitations

- It won't generate warnings for empty classes, as there are no class members
 (including base class sub-objects) to worry about.
- It won't generate warnings for base classes that have copy constructor
 private or deleted.
- It won't generate warnings for base classes that are initialized using other
 non-default constructor, as this could be intentional.

The check also suggests a fix-its in some cases.

```{title} clang-tidy - bugprone-copy-constructor-mutates-argument
```

# bugprone-copy-constructor-mutates-argument

Finds assignments to the copied object and its direct or indirect members
in copy constructors and copy assignment operators.

This check corresponds to the CERT C Coding Standard rule
[OOP58-CPP. Copy operations must not mutate the source object](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/object-oriented-programming-oop/oop58-cpp/).

```{title} clang-tidy - bugprone-crtp-constructor-accessibility
```

# bugprone-crtp-constructor-accessibility

Detects error-prone Curiously Recurring Template Pattern usage, when the CRTP
can be constructed outside itself and the derived class.

The CRTP is an idiom, in which a class derives from a template class, where
itself is the template argument. It should be ensured that if a class is
intended to be a base class in this idiom, it can only be instantiated if
the derived class is its template argument.

Example:

```c++
template <typename T> class CRTP {
private:
 CRTP() = default;
 friend T;
};

class Derived : CRTP<Derived> {};
```

Below can be seen some common mistakes that will allow the breaking of the
idiom.

If the constructor of a class intended to be used in a CRTP is public, then
it allows users to construct that class on its own.

Example:

```c++
template <typename T> class CRTP {
public:
 CRTP() = default;
};

class Good : CRTP<Good> {};
Good GoodInstance;

CRTP<int> BadInstance;
```

If the constructor is protected, the possibility of an accidental instantiation
is prevented, however it can fade an error, when a different class is used as
the template parameter instead of the derived one.

Example:

```c++
template <typename T> class CRTP {
protected:
 CRTP() = default;
};

class Good : CRTP<Good> {};
Good GoodInstance;

class Bad : CRTP<Good> {};
Bad BadInstance;
```

To ensure that no accidental instantiation happens, the best practice is to
make the constructor private and declare the derived class as friend. Note
that as a tradeoff, this also gives the derived class access to every other
private members of the CRTP. However, constructors can still be public or
protected if they are deleted.

Example:

```c++
template <typename T> class CRTP {
 CRTP() = default;
 friend T;
};

class Good : CRTP<Good> {};
Good GoodInstance;

class Bad : CRTP<Good> {};
Bad CompileTimeError;

CRTP<int> AlsoCompileTimeError;
```

## Limitations

- The check is not supported below C++11

- The check does not handle when the derived class is passed as a variadic
 template argument

- Accessible functions that can construct the CRTP, like factory functions
 are not checked

The check also suggests a fix-its in some cases.

```{title} clang-tidy - bugprone-dangling-handle
```

# bugprone-dangling-handle

Detect dangling references in value handles like `std::string_view`.
These dangling references can be a result of constructing handles from
temporary values, where the temporary is destroyed soon after the handle
is created.

Examples:

```c++
string_view View = string(); // View will dangle.
string A;
View = A + "A"; // still dangle.

vector<string_view> V;
V.push_back(string()); // V[0] is dangling.
V.resize(3, string()); // V[1] and V[2] will also dangle.

string_view f() {
 // All these return values will dangle.
 return string();
 string S;
 return S;
 char Array[10]{};
 return Array;
}

span<int> g() {
 array<int, 1> V;
 return {V};
 int Array[10]{};
 return {Array};
}
```

## Options

```{option} HandleClasses
A semicolon-separated list of class names that should be treated as handles.
By default only `std::basic_string_view`,
`std::experimental::basic_string_view` and `std::span` are considered.
```

```{title} clang-tidy - bugprone-default-operator-new-on-overaligned-type
```

# bugprone-default-operator-new-on-overaligned-type

Flags uses of default `operator new` where the type has extended
alignment (an alignment greater than the fundamental alignment).

The default `operator new` is guaranteed to provide the correct alignment
if the requested alignment is less or equal to the fundamental alignment.
Only cases are detected (by design) where the `operator new` is not
user-defined and is not a placement new (the reason is that in these cases we
assume that the user provided the correct memory allocation).

## References

This check corresponds to the CERT C++ Coding Standard rule
[MEM57-CPP. Avoid using default operator new for over-aligned types](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/memory-management-mem/mem57-cpp/).

```{title} clang-tidy - bugprone-derived-method-shadowing-base-method
```

# bugprone-derived-method-shadowing-base-method

Finds derived class methods that shadow a (non-virtual) base class method.

In order to be considered "shadowing", methods must have the same signature
(i.e. the same name, same number of parameters, same parameter types, etc).
Only checks public, non-templated methods.

The below example is bugprone because consumers of the `Derived` class will
expect the `reset` method to do the work of `Base::reset()` in addition to
extra work required to reset the `Derived` class. Common fixes include:

- Making the `reset` method polymorphic
- Re-naming `Derived::reset` if it's not meant to intersect with
 `Base::reset`
- Using `using Base::reset` to change the access specifier

This is also a violation of the Liskov Substitution Principle.

```c++
struct Base {
 void reset() {/* reset the base class */};
};

struct Derived : public Base {
 void reset() {/* reset the derived class, but not the base class */};
};
```

```{title} clang-tidy - bugprone-dynamic-static-initializers
```

# bugprone-dynamic-static-initializers

Finds instances of static variables that are dynamically initialized
in header files.

This can pose problems in certain multithreaded contexts. For example,
when disabling compiler generated synchronization instructions for
static variables initialized at runtime (e.g. by `-fno-threadsafe-statics`),
even if a particular project takes the necessary precautions to prevent race
conditions during initialization by providing their own synchronization, header
files included from other projects may not. Therefore, such a check is helpful
for ensuring that disabling compiler generated synchronization for static
variable initialization will not cause problems.

Consider the following code:

```cpp
int foo() {
 static int k = bar();
 return k;
}
```

When synchronization of static initialization is disabled, if two threads both
call `foo` for the first time, there is the possibility that `k` will be double
initialized, creating a race condition.

```{title} clang-tidy - bugprone-easily-swappable-parameters
```

# bugprone-easily-swappable-parameters

Finds function definitions where parameters of convertible types follow each
other directly, making call sites prone to calling the function with
swapped (or badly ordered) arguments.

```c++
void drawPoint(int X, int Y) { /* ... */ }
FILE *open(const char *Dir, const char *Name, Flags Mode) { /* ... */ }
```

A potential call like `drawPoint(-2, 5)` or
`openPath("a.txt", "tmp", Read)` is perfectly legal from the language's
perspective, but might not be what the developer of the function intended.

More elaborate and type-safe constructs, such as opaque typedefs or strong
types should be used instead, to prevent a mistaken order of arguments.

```c++
struct Coord2D { int X; int Y; };
void drawPoint(const Coord2D Pos) { /* ... */ }

FILE *open(const Path &Dir, const Filename &Name, Flags Mode) { /* ... */ }
```

Due to the potentially elaborate refactoring and API-breaking that is necessary
to strengthen the type safety of a project, no automatic fix-its are offered.

## Options

### Extension/relaxation options

Relaxation (or extension) options can be used to broaden the scope of the
analysis and fine-tune the enabling of more mixes between types.
Some mixes may depend on coding style or preference specific to a project,
however, it should be noted that enabling *all* of these relaxations model the
way of mixing at call sites the most.
These options are expected to make the check report for more functions, and
report longer mixable ranges.

````{option} QualifiersMix
Whether to consider parameters of some *cvr-qualified* `T` and a
differently *cvr-qualified* `T` (i.e. `T` and `const T`, `const T`
and `volatile T`, etc.) mixable between one another.
If `false`, the check will consider differently qualified types unmixable.
`True` turns the warnings on.
Defaults to `false`.

The following example produces a diagnostic only if `QualifiersMix` is
enabled:

```c++
void *memcpy(const void *Destination, void *Source, std::size_t N) { /* ... */ }
```
````

`````{option} ModelImplicitConversions
Whether to consider parameters of type `T` and `U` mixable if there
exists an implicit conversion from `T` to `U` and `U` to `T`.
If `false`, the check will not consider implicitly convertible types for
mixability.
`True` turns warnings for implicit conversions on.
Defaults to `true`.

The following examples produce a diagnostic only if
`ModelImplicitConversions` is enabled:

```c++
void fun(int Int, double Double) { /* ... */ }
void compare(const char *CharBuf, std::string String) { /* ... */ }
```

````{note}
Changing the qualifiers of an expression's type (e.g. from `int` to
`const int`) is defined as an *implicit conversion* in the C++
Standard.
However, the check separates this decision-making on the mixability of
differently qualified types based on whether `QualifiersMix` was
enabled.

For example, the following code snippet will only produce a diagnostic
if **both** `QualifiersMix` and `ModelImplicitConversions` are enabled:

```c++
void fun2(int Int, const double Double) { /* ... */ }
```
````
`````

### Filtering options

Filtering options can be used to lessen the size of the diagnostics emitted by
the checker, whether the aim is to ignore certain constructs or dampen the
noisiness.

```{option} MinimumLength
The minimum length required from an adjacent parameter sequence to be
diagnosed.
Defaults to `2`.
Might be any positive integer greater or equal to `2`.
If `0` or `1` is given, the default value `2` will be used instead.

For example, if `3` is specified, the examples above will not be matched.
```

```{option} IgnoredParameterNames
The list of parameter **names** that should never be considered part of a
swappable adjacent parameter sequence.
The value is a `;`-separated list of names.
To ignore unnamed parameters, add `""` to the list verbatim (not the
empty string, but the two quotes, potentially escaped!).
**This option is case-sensitive!**

By default, the following parameter names, and their Uppercase-initial
variants are ignored:
`""` (unnamed parameters), `iterator`, `begin`, `end`, `first`, `last`,
`lhs`, `rhs`.
```

```{option} IgnoredParameterTypeSuffixes
The list of parameter **type name suffixes** that should never be
considered part of a swappable adjacent parameter sequence.
Parameters which type, as written in the source code, end with an element
of this option will be ignored.
The value is a `;`-separated list of names.
**This option is case-sensitive!**

By default, the following, and their lowercase-initial variants are ignored:
`bool`, `It`, `Iterator`, `InputIt`, `ForwardIt`, `BidirIt`, `RandomIt`,
`random_iterator`, `ReverseIt`, `reverse_iterator`,
`reverse_const_iterator`, `RandomIt`, `random_iterator`, `ReverseIt`,
`reverse_iterator`, `reverse_const_iterator`, `Const_Iterator`,
`ConstIterator`, `const_reverse_iterator`, `ConstReverseIterator`.
In addition, `_Bool` (but not `_bool`) is also part of the default value.
```

````{option} SuppressParametersUsedTogether
Suppresses diagnostics about parameters that are used together or in a
similar fashion inside the function's body.
Defaults to `true`.
Specifying `false` will turn off the heuristics.

Currently, the following heuristics are implemented which will suppress the
warning about the parameter pair involved:

- The parameters are used in the same expression, e.g. `f(a, b)` or
 `a < b`.

- The parameters are further passed to the same function to the same
 parameter of that function, of the same overload.
 E.g. `f(a, 1)` and `f(b, 2)` to some `f(T, int)`.

 ```{note}
 The check does not perform path-sensitive analysis, and as such,
 "same function" in this context means the same function declaration.
 If the same member function of a type on two distinct instances are
 called with the parameters, it will still be regarded as
 "same function".
 ```

- The same member field is accessed, or member method is called of the
 two parameters, e.g. `a.foo()` and `b.foo()`.

- Separate `return` statements return either of the parameters on
 different code paths.
````

```{option} NamePrefixSuffixSilenceDissimilarityThreshold
The number of characters two parameter names might be different on *either*
the head or the tail end with the rest of the name the same so that the
warning about the two parameters are silenced.
Defaults to `1`.
Might be any positive integer.
If `0`, the filtering heuristic based on the parameters' names is turned
off.

This option can be used to silence warnings about parameters where the
naming scheme indicates that the order of those parameters do not matter.

For example, the parameters `LHS` and `RHS` are 1-dissimilar suffixes
of each other: `L` and `R` is the different character, while `HS`
is the common suffix.
Similarly, parameters `text1, text2, text3` are 1-dissimilar prefixes
of each other, with the numbers at the end being the dissimilar part.
If the value is at least `1`, such cases will not be reported.
```

## Limitations

**This check is designed to check function signatures!**

The check does not investigate functions that are generated by the compiler
in a context that is only determined from a call site.
These cases include variadic functions, functions in C code that do not have
an argument list, and C++ template instantiations.
Most of these cases, which are otherwise swappable from a caller's standpoint,
have no way of getting "fixed" at the definition point.
In the case of C++ templates, only primary template definitions and explicit
specializations are matched and analyzed.

None of the following cases produce a diagnostic:

```c++
int printf(const char *Format, ...) { /* ... */ }
int someOldCFunction() { /* ... */ }

template <typename T, typename U>
int add(T X, U Y) { return X + Y };

void theseAreNotWarnedAbout() {
 printf("%d %d\n", 1, 2); // Two ints passed, they could be swapped.
 someOldCFunction(1, 2, 3); // Similarly, multiple ints passed.

 add(1, 2); // Instantiates 'add<int, int>', but that's not a user-defined function.
}
```

Due to the limitation above, parameters which type are further dependent upon
template instantiations to *prove* that they mix with another parameter's is
not diagnosed.

```c++
template <typename T>
struct Vector {
 typedef T element_type;
};

// Diagnosed: Explicit instantiation was done by the user, we can prove it
// is the same type.
void instantiated(int A, Vector<int>::element_type B) { /* ... */ }

// Diagnosed: The two parameter types are exactly the same.
template <typename T>
void exact(typename Vector<T>::element_type A,
 typename Vector<T>::element_type B) { /* ... */ }

// Skipped: The two parameters are both 'T' but we cannot prove this
// without actually instantiating.
template <typename T>
void falseNegative(T A, typename Vector<T>::element_type B) { /* ... */ }
```

In the context of *implicit conversions* (when
`ModelImplicitConversions` is
enabled), the modelling performed by the check
warns if the parameters are swappable and the swapped order matches implicit
conversions.
It does not model whether there exists an unrelated third type from which
*both* parameters can be given in a function call.
This means that in the following example, even while `strs()` clearly carries
the possibility to be called with swapped arguments (as long as the arguments
are string literals), will not be warned about.

```c++
struct String {
 String(const char *Buf);
};

struct StringView {
 StringView(const char *Buf);
 operator const char *() const;
};

// Skipped: Directly swapping expressions of the two type cannot mix.
// (Note: StringView -> const char * -> String would be **two**
// user-defined conversions, which is disallowed by the language.)
void strs(String Str, StringView SV) { /* ... */ }

// Diagnosed: StringView implicitly converts to and from a buffer.
void cStr(StringView SV, const char *Buf() { /* ... */ }
```

```{title} clang-tidy - bugprone-empty-catch
```

# bugprone-empty-catch

Detects and suggests addressing issues with empty catch statements.

```c++
try {
 // Some code that can throw an exception
} catch(const std::exception&) {
}
```

Having empty catch statements in a codebase can be a serious problem that
developers should be aware of. Catch statements are used to handle exceptions
that are thrown during program execution. When an exception is thrown, the
program jumps to the nearest catch statement that matches the type of the
exception.

Empty catch statements, also known as "swallowing" exceptions, catch the
exception but do nothing with it. This means that the exception is not handled
properly, and the program continues to run as if nothing happened. This can
lead to several issues, such as:

- *Hidden Bugs*: If an exception is caught and ignored, it can lead to hidden
 bugs that are difficult to diagnose and fix. The root cause of the problem
 may not be apparent, and the program may continue to behave in unexpected
 ways.
- *Security Issues*: Ignoring exceptions can lead to security issues, such as
 buffer overflows or null pointer dereferences. Hackers can exploit these
 vulnerabilities to gain access to sensitive data or execute malicious code.
- *Poor Code Quality*: Empty catch statements can indicate poor code quality
 and a lack of attention to detail. This can make the codebase difficult to
 maintain and update, leading to longer development cycles and increased
 costs.
- *Unreliable Code*: Code that ignores exceptions is often unreliable and can
 lead to unpredictable behavior. This can cause frustration for users and
 erode trust in the software.

To avoid these issues, developers should always handle exceptions properly.
This means either fixing the underlying issue that caused the exception or
propagating the exception up the call stack to a higher-level handler.
If an exception is not important, it should still be logged or reported in
some way so that it can be tracked and addressed later.

If the exception is something that can be handled locally, then it should be
handled within the catch block. This could involve logging the exception or
taking other appropriate action to ensure that the exception is not ignored.

Here is an example:

```c++
try {
 // Some code that can throw an exception
} catch (const std::exception& ex) {
 // Properly handle the exception, e.g.:
 std::cerr << "Exception caught: " << ex.what() << std::endl;
}
```

If the exception cannot be handled locally and needs to be propagated up the
call stack, it should be re-thrown or new exception should be thrown.

Here is an example:

```c++
try {
 // Some code that can throw an exception
} catch (const std::exception& ex) {
 // Re-throw the exception
 throw;
}
```

In some cases, catching the exception at this level may not be necessary, and
it may be appropriate to let the exception propagate up the call stack.
This can be done simply by not using `try/catch` block.

Here is an example:

```c++
void function() {
 // Some code that can throw an exception
}

void callerFunction() {
 try {
 function();
 } catch (const std::exception& ex) {
 // Handling exception on higher level
 std::cerr << "Exception caught: " << ex.what() << std::endl;
 }
}
```

Other potential solution to avoid empty catch statements is to modify the code
to avoid throwing the exception in the first place. This can be achieved by
using a different API, checking for error conditions beforehand, or handling
errors in a different way that does not involve exceptions. By eliminating the
need for try-catch blocks, the code becomes simpler and less error-prone.

Here is an example:

```c++
// Old code:
try {
 mapContainer["Key"].callFunction();
} catch(const std::out_of_range&) {
}

// New code
if (auto it = mapContainer.find("Key"); it != mapContainer.end()) {
 it->second.callFunction();
}
```

In conclusion, empty catch statements are a bad practice that can lead to
hidden bugs, security issues, poor code quality, and unreliable code. By
handling exceptions properly, developers can ensure that their code is
robust, secure, and maintainable.

## Options

```{option} IgnoreCatchWithKeywords
This option can be used to ignore specific catch statements containing
certain keywords. If a `catch` statement body contains (case-insensitive)
any of the keywords listed in this semicolon-separated option, then the
catch will be ignored, and no warning will be raised.
Default value: `@TODO;@FIXME`.
```

```{option} AllowEmptyCatchForExceptions
This option can be used to ignore empty catch statements for specific
exception types. By default, the check will raise a warning if an empty
catch statement is detected, regardless of the type of exception being
caught. However, in certain situations, such as when a developer wants to
intentionally ignore certain exceptions or handle them in a different way,
it may be desirable to allow empty catch statements for specific exception
types.
To configure this option, a semicolon-separated list of exception type names
should be provided. If an exception type name in the list is caught in an
empty catch statement, no warning will be raised.
Default value: empty string.
```

```{title} clang-tidy - bugprone-exception-copy-constructor-throws
```

# bugprone-exception-copy-constructor-throws

Checks whether a thrown object's copy constructor can throw.

Exception objects are required to be copy constructible in C++. However, an
exception's copy constructor should not throw to avoid potential issues when
unwinding the stack. If an exception is thrown during stack unwinding (such
as from a copy constructor of an exception object), the program will
terminate via `std::terminate`.

```c++
class SomeException {
public:
 SomeException() = default;
 SomeException(const SomeException&) { /* may throw */ }
};

void f() {
 throw SomeException(); // warning: thrown exception type's copy constructor can throw
}
```

## References

This check corresponds to the CERT C++ Coding Standard rule
[ERR60-CPP. Exception objects must be nothrow copy constructible](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/exceptions-and-error-handling-err/err60-cpp/).

---
myst:
 enable_extensions:
 - deflist
---

```{title} clang-tidy - bugprone-exception-escape
```

# bugprone-exception-escape

Finds functions which may throw an exception directly or indirectly, but they
should not. The functions which should not throw exceptions are the following:

- Destructors
- Move constructors
- Move assignment operators
- The `main()` functions
- `swap()` functions
- `iter_swap()` functions
- `iter_move()` functions
- Functions marked with `throw()` or `noexcept`
- Other functions given as option

A destructor throwing an exception may result in undefined behavior, resource
leaks or unexpected termination of the program. Throwing move constructor or
move assignment also may result in undefined behavior or resource leak. The
`swap()` operations expected to be non throwing most of the cases and they
are always possible to implement in a non throwing way. Non throwing `swap()`
operations are also used to create move operations. A throwing `main()`
function also results in unexpected termination.

Functions declared explicitly with `noexcept(false)` or `throw(exception)`
will be excluded from the analysis, as even though it is not recommended for
functions like `swap()`, `main()`, move constructors, move assignment
operators and destructors, it is a clear indication of the developer's
intention and should be respected. To check if these special functions are
marked as potentially throwing, the check
{doc}`bugprone-unsafe-to-allow-exceptions <unsafe-to-allow-exceptions>` can be
used.

WARNING! This check may be expensive on large source files.

## Options

```{option} CheckDestructors
When `true`, destructors are analyzed to not throw exceptions.
Default value is `true`.
```

```{option} CheckMoveMemberFunctions
When `true`, move constructors and move assignment operators are analyzed
to not throw exceptions. Default value is `true`.
```

```{option} CheckMain
When `true`, the `main()` function is analyzed to not throw exceptions.
Default value is `true`.
```

```{option} CheckNothrowFunctions
When `true`, functions marked with `noexcept` or `throw()` exception
specifications are analyzed to not throw exceptions. Default value is `true`.
```

```{option} CheckedSwapFunctions
Comma-separated list of swap function names which should not throw exceptions.
Default value is `swap,iter_swap,iter_move`.
```

```{option} FunctionsThatShouldNotThrow
Comma separated list containing function names which should not throw. An
example value for this parameter can be `WinMain` which adds function
`WinMain()` in the Windows API to the list of the functions which should
not throw. Default value is an empty string.
```

```{option} IgnoredExceptions
Comma separated list containing type names which are not counted as thrown
exceptions in the check. Default value is an empty string.
```

```{option} TreatFunctionsWithoutSpecificationAsThrowing
Determines which functions are considered as throwing if they do not have
an explicit exception specification. It can be set to the following values:

- `None`
 : The check will consider functions without an explicit exception
 specification as throwing only if they have a visible definition which
 can be deduced to throw.
- `OnlyUndefined`
 : The check will consider functions with only a declaration available and
 no visible definition as throwing.
- `All`
 : The check will consider all functions without an explicit exception
 specification (such as `noexcept`) as throwing, even if they have a
 visible definition and do not contain any throwing statements.

Default value is `None`.
```

```{title} clang-tidy - bugprone-float-loop-counter
```

# bugprone-float-loop-counter

Flags `for` loops where the induction expression has a floating-point type.

## References

This check corresponds to the CERT C Coding Standard rule
[FLP30-C. Do not use floating-point variables as loop counters](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/floating-point-flp/flp30-c/).

```{title} clang-tidy - bugprone-fold-init-type
```

# bugprone-fold-init-type

The check flags type mismatches in
[folds](<https://en.wikipedia.org/wiki/Fold_(higher-order_function)>)
that might result in loss of precision.

The check supports the following functions:

- `std::accumulate`
- `std::reduce`
- `std::inner_product`

These functions fold an input range into an initial value using the type of the
latter. By default, `std::accumulate` and `std::reduce` use `operator+`
while `std::inner_product` uses `operator+` and `operator*`. This can
cause loss of precision through:

- Truncation: The following code uses a floating point range and an int
 initial value, so truncation will happen at every application of
 `operator+` and the result will be `0`, which might not be what the
 user expected.

```c++
auto a = {0.5f, 0.5f, 0.5f, 0.5f};
return std::accumulate(std::begin(a), std::end(a), 0);
```

- Overflow: The following code also returns `0`.

```c++
auto a = {65536LL * 65536 * 65536};
return std::accumulate(std::begin(a), std::end(a), 0);
```

The check handles overloads with the following transparent standard functors:

- `std::plus`
- `std::minus`
- `std::multiplies`
- `std::divides`
- `std::bit_and`
- `std::bit_or`
- `std::bit_xor`

```{title} clang-tidy - bugprone-forward-declaration-namespace
```

# bugprone-forward-declaration-namespace

Checks if an unused forward declaration is in a wrong namespace.

The check inspects all unused forward declarations and checks if there is any
declaration/definition with the same name existing, which could indicate that
the forward declaration is in a potentially wrong namespace.

```cpp
namespace na { struct A; }
namespace nb { struct A {}; }
nb::A a;
// warning : no definition found for 'A', but a definition with the same name
// 'A' found in another namespace 'nb::'
```

This check can only generate warnings, but it can't suggest a fix at this
point.

```{title} clang-tidy - bugprone-forwarding-reference-overload
```

# bugprone-forwarding-reference-overload

The check looks for perfect forwarding constructors that can hide copy or move
constructors. If a non const lvalue reference is passed to the constructor, the
forwarding reference parameter will be a better match than the const reference
parameter of the copy constructor, so the perfect forwarding constructor will
be called, which can be confusing.
For detailed description of this issue see: Scott Meyers, Effective Modern C++,
Item 26.

Consider the following example:

```c++
class Person {
public:
 // C1: perfect forwarding ctor
 template<typename T>
 explicit Person(T&& n) {}

 // C2: perfect forwarding ctor with parameter default value
 template<typename T>
 explicit Person(T&& n, int x = 1) {}

 // C3: perfect forwarding ctor guarded with enable_if
 template<typename T, typename X = enable_if_t<is_special<T>, void>>
 explicit Person(T&& n) {}

 // C4: variadic perfect forwarding ctor guarded with enable_if
 template<typename... A,
 enable_if_t<is_constructible_v<tuple<string, int>, A&&...>, int> = 0>
 explicit Person(A&&... a) {}

 // C5: perfect forwarding ctor guarded with requires expression
 template<typename T>
 requires requires { is_special<T>; }
 explicit Person(T&& n) {}

 // C6: perfect forwarding ctor guarded with concept requirement
 template<Special T>
 explicit Person(T&& n) {}

 // (possibly compiler generated) copy ctor
 Person(const Person& rhs);
};
```

The check warns for constructors C1 and C2, because those can hide copy and
move constructors. We suppress warnings if the copy and the move constructors
are both disabled (deleted or private), because there is nothing the perfect
forwarding constructor could hide in this case. We also suppress warnings for
constructors like C3-C6 that are guarded with an `enable_if` or a concept,
assuming the programmer was aware of the possible hiding.

## Background

For deciding whether a constructor is guarded with enable_if, we consider the
types of the constructor parameters, the default values of template type parameters
and the types of non-type template parameters with a default literal value. If any
part of these types is `std::enable_if` or `std::enable_if_t`, we assume the
constructor is guarded.

```{title} clang-tidy - bugprone-implicit-widening-of-multiplication-result
```

# bugprone-implicit-widening-of-multiplication-result

The check diagnoses instances where a result of a multiplication is implicitly
widened, and suggests (with fix-it) to either silence the code by making
widening explicit, or to perform the multiplication in a wider type,
to avoid the widening afterwards.

This is mainly useful when operating on very large buffers.
For example, consider:

```c++
void zeroinit(char* base, unsigned width, unsigned height) {
 for(unsigned row = 0; row != height; ++row) {
 for(unsigned col = 0; col != width; ++col) {
 char* ptr = base + row * width + col;
 *ptr = 0;
 }
 }
}
```

This is fine in general, but if `width * height` overflows,
you end up wrapping back to the beginning of `base`
instead of processing the entire requested buffer.

Indeed, this only matters for pretty large buffers (4GB+),
but that can happen very easily for example in image processing,
where for that to happen you "only" need a ~269MPix image.

## Options

```{option} UseCXXStaticCastsInCppSources
When suggesting fix-its for C++ code, should C++-style `static_cast<>()`'s
be suggested, or C-style casts. Defaults to `true`.
```

```{option} UseCXXHeadersInCppSources
When suggesting to include the appropriate header in C++ code,
should `<cstddef>` header be suggested, or `<stddef.h>`.
Defaults to `true`.
```

```{option} IgnoreConstantIntExpr
If the multiplication operands are compile-time constants (like literals or
are `constexpr`) and fit within the source expression type, do not emit a
diagnostic or suggested fix. Only considers expressions where the source
expression is a signed integer type. Defaults to `false`.
```

Examples:

```c++
long mul(int a, int b) {
 return a * b; // warning: performing an implicit widening conversion to type 'long' of a multiplication performed in type 'int'
}

char* ptr_add(char *base, int a, int b) {
 return base + a * b; // warning: result of multiplication in type 'int' is used as a pointer offset after an implicit widening conversion to type 'ssize_t'
}

char ptr_subscript(char *base, int a, int b) {
 return base[a * b]; // warning: result of multiplication in type 'int' is used as a pointer offset after an implicit widening conversion to type 'ssize_t'
}
```

```{title} clang-tidy - bugprone-inaccurate-erase
```

# bugprone-inaccurate-erase

Checks for inaccurate use of the `erase()` method.

Algorithms like `remove()` do not actually remove any element from the
container but return an iterator to the first redundant element at the end
of the container. These redundant elements must be removed using the
`erase()` method. This check warns when not all of the elements will be
removed due to using an inappropriate overload.

For example, the following code erases only one element:

```cpp
std::vector<int> xs;
...
xs.erase(std::remove(xs.begin(), xs.end(), 10));
```

Call the two-argument overload of `erase()` to remove the subrange:

```cpp
std::vector<int> xs;
...
xs.erase(std::remove(xs.begin(), xs.end(), 10), xs.end());
```

```{title} clang-tidy - bugprone-inc-dec-in-conditions
```

# bugprone-inc-dec-in-conditions

Detects when a variable is both incremented/decremented and referenced inside a
complex condition and suggests moving them outside to avoid ambiguity in the
variable's value.

When a variable is modified and also used in a complex condition, it can lead
to unexpected behavior. The side-effect of changing the variable's value within
the condition can make the code difficult to reason about. Additionally, the
developer's intended timing for the modification of the variable may not be
clear, leading to misunderstandings and errors. This can be particularly
problematic when the condition involves logical operators like `&&` and
`||`, where the order of evaluation can further complicate the situation.

Consider the following example:

```c++
int i = 0;
// ...
if (i++ < 5 && i > 0) {
 // do something
}
```

In this example, the result of the expression may not be what the developer
intended. The original intention of the developer could be to increment `i`
after the entire condition is evaluated, but in reality, i will be incremented
before `i > 0` is executed. This can lead to unexpected behavior and bugs in
the code. To fix this issue, the developer should separate the increment
operation from the condition and perform it separately. For example, they can
increment `i` in a separate statement before or after the condition is
evaluated. This ensures that the value of `i` is predictable and consistent
throughout the code.

```c++
int i = 0;
// ...
i++;
if (i <= 5 && i > 0) {
 // do something
}
```

Another common issue occurs when multiple increments or decrements are
performed on the same variable inside a complex condition. For example:

```c++
int i = 4;
// ...
if (i++ < 5 || --i > 2) {
 // do something
}
```

There is a potential issue with this code due to the order of evaluation in
C++. The `||` operator used in the condition statement guarantees that if
the first operand evaluates to `true`, the second operand will not be
evaluated. This means that if `i` were initially `4`, the first operand
`i < 5` would evaluate to `true` and the second operand `i > 2` would
not be evaluated. As a result, the decrement operation `--i` would not be
executed and `i` would hold value `5`, which may not be the intended
behavior for the developer.

To avoid this potential issue, the both increment and decrement operation on
`i` should be moved outside the condition statement.

```{title} clang-tidy - bugprone-incorrect-enable-if
```

# bugprone-incorrect-enable-if

Detects incorrect usages of `std::enable_if` that don't name the nested
`type` type.

In C++11 introduced `std::enable_if` as a convenient way to leverage SFINAE.
One form of using `std::enable_if` is to declare an unnamed template type
parameter with a default type equal to
`typename std::enable_if<condition>::type`. If the author forgets to name
the nested type `type`, then the code will always consider the candidate
template even if the condition is not met.

Below are some examples of code using `std::enable_if` correctly and
incorrect examples that this check flags.

```c++
template <typename T, typename = typename std::enable_if<T::some_trait>::type>
void valid_usage() { ... }

template <typename T, typename = std::enable_if_t<T::some_trait>>
void valid_usage_with_trait_helpers() { ... }

// The below code is not a correct application of SFINAE. Even if
// T::some_trait is not true, the function will still be considered in the
// set of function candidates. It can either incorrectly select the function
// when it should not be a candidates, and/or lead to hard compile errors
// if the body of the template does not compile if the condition is not
// satisfied.
template <typename T, typename = std::enable_if<T::some_trait>>
void invalid_usage() { ... }

// The tool suggests the following replacement for 'invalid_usage':
template <typename T, typename = typename std::enable_if<T::some_trait>::type>
void fixed_invalid_usage() { ... }
```

C++14 introduced the trait helper `std::enable_if_t` which reduces the
likelihood of this error. C++20 introduces constraints, which generally
supersede the use of `std::enable_if`. See
{doc}`modernize-type-traits <../modernize/type-traits>` for another tool
that will replace `std::enable_if` with
`std::enable_if_t`, and see
{doc}`modernize-use-constraints <../modernize/use-constraints>` for another
tool that replaces `std::enable_if` with C++20 constraints. Consider these
newer mechanisms where possible.

```{title} clang-tidy - bugprone-incorrect-enable-shared-from-this
```

# bugprone-incorrect-enable-shared-from-this

Detect classes or structs that do not publicly inherit from
`std::enable_shared_from_this`, because unintended behavior will
otherwise occur when calling `shared_from_this`.

Consider the following code:

````c++
#include <memory>

// private inheritance
class BadExample : std::enable_shared_from_this<BadExample> {

// ``shared_from_this``` unintended behaviour
// `libstdc++` implementation returns uninitialized ``weak_ptr``
 public:
 BadExample* foo() { return shared_from_this().get(); }
 void bar() { return; }
};

void using_not_public() {
 auto bad_example = std::make_shared<BadExample>();
 auto* b_ex = bad_example->foo();
 b_ex->bar();
}
````

Using `libstdc++` implementation, `shared_from_this` will throw
`std::bad_weak_ptr`. When `using_not_public()` is called, this code will
crash without exception handling.

```{title} clang-tidy - bugprone-incorrect-roundings
```

# bugprone-incorrect-roundings

Checks the usage of patterns known to produce incorrect rounding.
Programmers often use:

```cpp
(int)(double_expression + 0.5)
```

to round the double expression to an integer. The problem with this:

1. It is unnecessarily slow.
2. It is incorrect. The number 0.499999975 (smallest representable float
 number below 0.5) rounds to 1.0. Even worse behavior for negative
 numbers where both -0.5f and -1.4f both round to 0.0.

```{title} clang-tidy - bugprone-infinite-loop
```

# bugprone-infinite-loop

Finds obvious infinite loops (loops where the condition variable is not changed
at all).

Finding infinite loops is well-known to be impossible (halting problem).
However, it is possible to detect some obvious infinite loops, for example, if
the loop condition is not changed. This check detects such loops. A loop is
considered infinite if it does not have any loop exit statement (`break`,
`continue`, `goto`, `return`, `throw` or a call to a function called as
`[[noreturn]]`) and all of the following conditions hold for every variable
in the condition:

- It is a local variable.
- It has no reference or pointer aliases.
- It is not a structure or class member.

Furthermore, the condition must not contain a function call to consider the
loop infinite since functions may return different values for different calls.

For example, the following loop is considered infinite `i` is not changed in
the body:

```c++
int i = 0, j = 0;
while (i < 10) {
 ++j;
}
```

```{title} clang-tidy - bugprone-integer-division
```

# bugprone-integer-division

Finds cases where integer division in a floating point context is likely to
cause unintended loss of precision.

No reports are made if divisions are part of the following expressions:

- operands of operators expecting integral or bool types,
- call expressions of integral or bool types, and
- explicit cast expressions to integral or bool types,

as these are interpreted as signs of deliberateness from the programmer.

Examples:

```c++
float floatFunc(float);
int intFunc(int);
double d;
int i = 42;

// Warn, floating-point values expected.
d = 32 * 8 / (2 + i);
d = 8 * floatFunc(1 + 7 / 2);
d = i / (1 << 4);

// OK, no integer division.
d = 32 * 8.0 / (2 + i);
d = 8 * floatFunc(1 + 7.0 / 2);
d = (double)i / (1 << 4);

// OK, there are signs of deliberateness.
d = 1 << (i / 2);
d = 9 + intFunc(6 * i / 32);
d = (int)(i / 32) - 8;
```

```{title} clang-tidy - bugprone-invalid-enum-default-initialization
```

# bugprone-invalid-enum-default-initialization

Detects default initialization (to 0) of variables with `enum` type where
the enum has no enumerator with value of 0.

In C++ a default initialization is performed if a variable is initialized with
initializer list or in other implicit ways, and no value is specified at the
initialization. In such cases the value 0 is used for the initialization.
This also applies to enumerations even if it does not have an enumerator with
value 0. In this way a variable with the `enum` type may contain initially an
invalid value (if the program expects that it contains only the listed
enumerator values).

The check emits a warning only if an `enum` variable is default-initialized
(contrary to not initialized) and the `enum` does not have an enumerator with
value of 0. The type can be a scoped or non-scoped `enum`. Unions are not
handled by the check (if it contains a member of enumeration type).

Note that the `enum` `std::errc` is always ignored because it is expected
to be default initialized, despite not defining an enumerator with the value 0.

```c++
enum class Enum1: int {
 A = 1,
 B
};

enum class Enum0: int {
 A = 0,
 B
};

void f() {
 Enum1 X1{}; // warn: 'X1' is initialized to 0
 Enum1 X2 = Enum1(); // warn: 'X2' is initialized to 0
 Enum1 X3; // no warning: 'X3' is not initialized
 Enum0 X4{}; // no warning: type has an enumerator with value of 0
}

struct S1 {
 Enum1 A;
 S(): A() {} // warn: 'A' is initialized to 0
};

struct S2 {
 int A;
 Enum1 B;
};

S2 VarS2{}; // warn: member 'B' is initialized to 0
```

The check applies to initialization of arrays or structures with initialization
lists in C code too. In these cases elements not specified in the list (and have
enum type) are set to 0.

```c
enum Enum1 {
 Enum1_A = 1,
 Enum1_B
};
struct Struct1 {
 int a;
 enum Enum1 b;
};

enum Enum1 Array1[2] = {Enum1_A}; // warn: omitted elements are initialized to 0
enum Enum1 Array2[2][2] = {{Enum1_A}, {Enum1_A}}; // warn: last element of both nested arrays is initialized to 0
enum Enum1 Array3[2][2] = {{Enum1_A, Enum1_A}}; // warn: elements of second array are initialized to 0

struct Struct1 S1 = {1}; // warn: element 'b' is initialized to 0
```

## Options

```{option} IgnoredEnums
Semicolon-separated list of regexes specifying enums for which this check won't be
enforced. Default is `::std::errc`.
```

```{title} clang-tidy - bugprone-lambda-function-name
```

# bugprone-lambda-function-name

Checks for attempts to get the name of a function from within a lambda
expression. The name of a lambda is always something like `operator()`, which
is almost never what was intended.

Example:

```c++
void FancyFunction() {
 [] { printf("Called from %s\n", __func__); }();
 [] { printf("Now called from %s\n", __FUNCTION__); }();
}
```

Output:

```
Called from operator()
Now called from operator()
```

Likely intended output:

```
Called from FancyFunction
Now called from FancyFunction
```

## Options

```{option} IgnoreMacros
The value `true` specifies that attempting to get the name of a function from
within a macro should not be diagnosed. The default value is `false`.
```

```{title} clang-tidy - bugprone-macro-parentheses
```

# bugprone-macro-parentheses

Finds macros that can have unexpected behavior due to missing parentheses.

Macros are expanded by the preprocessor as-is. As a result, there can be
unexpected behavior; operators may be evaluated in unexpected order and
unary operators may become binary operators, etc.

When the replacement list has an expression, it is recommended to surround
it with parentheses. This ensures that the macro result is evaluated
completely before it is used.

It is also recommended to surround macro arguments in the replacement list
with parentheses. This ensures that the argument value is calculated
properly.

This check corresponds to the CERT C Coding Standard rule
[PRE02-C. Macro replacement lists should be parenthesized.](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/recommendations/preprocessor-pre/pre02-c/)

```{title} clang-tidy - bugprone-macro-repeated-side-effects
```

# bugprone-macro-repeated-side-effects

Checks for repeated argument with side effects in macros.

```{title} clang-tidy - bugprone-misleading-setter-of-reference
```

# bugprone-misleading-setter-of-reference

Finds setter-like member functions that take a pointer parameter and set a
reference member of the same class with the pointed value.

The check detects member functions that take a single pointer parameter,
and contain a single expression statement that dereferences the parameter and
assigns the result to a data member with a reference type.

The fact that a setter function takes a pointer might cause the belief that an
internal reference (if it would be a pointer) is changed instead of the
pointed-to (or referenced) value.

Example:

```c++
class MyClass {
 int &InternalRef; // non-const reference member
public:
 MyClass(int &Value) : InternalRef(Value) {}

 // Warning: This setter could lead to unintended behaviour.
 void setRef(int *Value) {
 InternalRef = *Value; // This assigns to the referenced value, not changing what InternalRef references.
 }
};

int main() {
 int Value1 = 42;
 int Value2 = 100;
 MyClass X(Value1);

 // This might look like it changes what InternalRef references to,
 // but it actually modifies Value1 to be 100.
 X.setRef(&Value2);
}
```

Possible fixes:

- Change the parameter type of the "set" function to non-pointer type (for
 example, a const reference).
- Change the type of the member variable to a pointer and in the "set"
 function assign a value to the pointer (without dereference).

```{title} clang-tidy - bugprone-misplaced-operator-in-strlen-in-alloc
```

# bugprone-misplaced-operator-in-strlen-in-alloc

Finds cases where `1` is added to the string in the argument to `strlen()`,
`strnlen()`, `strnlen_s()`, `wcslen()`, `wcsnlen()`, and
`wcsnlen_s()` instead of the result and the value is used as an argument to a
memory allocation function (`malloc()`, `calloc()`, `realloc()`,
`alloca()`) or the `new[]` operator in `C++`. The check detects error cases
even if one of these functions (except the `new[]` operator) is called by a
constant function pointer. Cases where `1` is added both to the parameter and
the result of the `strlen()`-like function are ignored, as are cases where
the whole addition is surrounded by extra parentheses.

`C` example code:

```c
void bad_malloc(char *str) {
 char *c = (char*) malloc(strlen(str + 1));
}
```

The suggested fix is to add `1` to the return value of `strlen()` and not
to its argument. In the example above the fix would be

```c
char *c = (char*) malloc(strlen(str) + 1);
```

`C++` example code:

```c++
void bad_new(char *str) {
 char *c = new char[strlen(str + 1)];
}
```

As in the `C` code with the `malloc()` function, the suggested fix is to
add `1` to the return value of `strlen()` and not to its argument. In the
example above the fix would be

```c++
char *c = new char[strlen(str) + 1];
```

Example for silencing the diagnostic:

```c
void bad_malloc(char *str) {
 char *c = (char*) malloc(strlen((str + 1)));
}
```

```{title} clang-tidy - bugprone-misplaced-pointer-arithmetic-in-alloc
```

# bugprone-misplaced-pointer-arithmetic-in-alloc

Finds cases where an integer expression is added to or subtracted from the
result of a memory allocation function (`malloc()`, `calloc()`,
`realloc()`, `alloca()`) instead of its argument. The check detects error
cases even if one of these functions is called by a constant function pointer.

Example code:

```c
void bad_malloc(int n) {
 char *p = (char*) malloc(n) + 10;
}
```

The suggested fix is to add the integer expression to the argument of
`malloc` and not to its result. In the example above the fix would be

```c
char *p = (char*) malloc(n + 10);
```

```{title} clang-tidy - bugprone-misplaced-widening-cast
```

# bugprone-misplaced-widening-cast

This check will warn when there is a cast of a calculation result to a bigger
type. If the intention of the cast is to avoid loss of precision then the cast
is misplaced, and there can be loss of precision. Otherwise the cast is
ineffective.

Example code:

```c++
long f(int x) {
 return (long)(x * 1000);
}
```

The result `x * 1000` is first calculated using `int` precision. If the
result exceeds `int` precision there is loss of precision. Then the result is
casted to `long`.

If there is no loss of precision then the cast can be removed or you can
explicitly cast to `int` instead.

If you want to avoid loss of precision then put the cast in a proper location,
for instance:

```c++
long f(int x) {
 return (long)x * 1000;
}
```

## Implicit casts

Forgetting to place the cast at all is at least as dangerous and at least as
common as misplacing it. If {option}`CheckImplicitCasts` is enabled the check
also detects these cases, for instance:

```c++
long f(int x) {
 return x * 1000;
}
```

## Floating point

Currently warnings are only written for integer conversion. No warning is
written for this code:

```c++
double f(float x) {
 return (double)(x * 10.0f);
}
```

## Options

```{option} CheckImplicitCasts
If `true`, enables detection of implicit casts. Default is `false`.
```

```{title} clang-tidy - bugprone-missing-end-comparison
```

# bugprone-missing-end-comparison

Finds instances where the result of a standard algorithm is used in a Boolean
context without being compared to the end iterator.

Standard algorithms such as `std::find`, `std::search`, and
`std::lower_bound` return an iterator to the element if found, or the end
iterator otherwise.

Using the result directly in a Boolean context (like an `if` statement) is
almost always a bug, as it only checks if the iterator itself evaluates to
`true`, which may always be true for many iterator types.

Examples:

```c++
void example() {
 int arr[] = {1, 2, 3};
 int* begin = std::begin(arr);
 int* end = std::end(arr);

 if (std::find(begin, end, 2)) {
 // ...
 }

 // Fixed by the check:
 if ((std::find(begin, end, 2) != end)) {
 // ...
 }

 // C++20 ranges:
 int v[] = {1, 2, 3};
 if (std::ranges::find(v, 2)) {
 // ...
 }

 // Fixed by the check:
 if ((std::ranges::find(v, 2) != std::ranges::end(v))) {
 // ...
 }
}
```

The check also handles range-based algorithms introduced in C++20.

Supported algorithms:

- `std::adjacent_find`
- `std::find`
- `std::find_end`
- `std::find_first_of`
- `std::find_if`
- `std::find_if_not`
- `std::is_sorted_until`
- `std::lower_bound`
- `std::max_element`
- `std::min_element`
- `std::partition_point`
- `std::search`
- `std::search_n`
- `std::upper_bound`
- `std::ranges::adjacent_find`
- `std::ranges::find`
- `std::ranges::find_first_of`
- `std::ranges::find_if`
- `std::ranges::find_if_not`
- `std::ranges::is_sorted_until`
- `std::ranges::lower_bound`
- `std::ranges::max_element`
- `std::ranges::min_element`
- `std::ranges::upper_bound`

## Options

```{option} ExtraAlgorithms
A semicolon-separated list of extra algorithms to check.
The list can contain:

- Iterator-based algorithms. These should follow the standard iterator
 pattern: `func(Iter, Iter, ...)`.
- Range-based algorithms. These are heuristically detected if they take
 exactly two arguments and the first argument is a container or range.
 The fix will insert `std::end(Container)`.

Default is an empty string.
```

```{title} clang-tidy - bugprone-move-forwarding-reference
```

# bugprone-move-forwarding-reference

Warns if `std::move` is called on a forwarding reference, for example:

```c++
template <typename T>
void foo(T&& t) {
 bar(std::move(t));
}
```

[Forwarding references](http://www.open-std.org/jtc1/sc22/wg21/docs/papers/2014/n4164.pdf) should
typically be passed to `std::forward` instead of `std::move`, and this is
the fix that will be suggested.

(A forwarding reference is an rvalue reference of a type that is a deduced
function template argument.)

In this example, the suggested fix would be

```c++
bar(std::forward<T>(t));
```

## Background

Code like the example above is sometimes written with the expectation that
`T&&` will always end up being an rvalue reference, no matter what type is
deduced for `T`, and that it is therefore not possible to pass an lvalue to
`foo()`. However, this is not true. Consider this example:

```c++
std::string s = "Hello, world";
foo(s);
```

This code compiles and, after the call to `foo()`, `s` is left in an
indeterminate state because it has been moved from. This may be surprising to
the caller of `foo()` because no `std::move` was used when calling
`foo()`.

The reason for this behavior lies in the special rule for template argument
deduction on function templates like `foo()` -- i.e. on function templates
that take an rvalue reference argument of a type that is a deduced function
template argument. (See section [temp.deduct.call]/3 in the C++11 standard.)

If `foo()` is called on an lvalue (as in the example above), then `T` is
deduced to be an lvalue reference. In the example, `T` is deduced to be
`std::string &`. The type of the argument `t` therefore becomes
`std::string& &&`; by the reference collapsing rules, this collapses to
`std::string&`.

This means that the `foo(s)` call passes `s` as an lvalue reference, and
`foo()` ends up moving `s` and thereby placing it into an indeterminate
state.

```{title} clang-tidy - bugprone-multi-level-implicit-pointer-conversion
```

# bugprone-multi-level-implicit-pointer-conversion

Detects implicit conversions between pointers of different levels of
indirection.

Conversions between pointer types of different levels of indirection can be
dangerous and may lead to undefined behavior, particularly if the converted
pointer is later cast to a type with a different level of indirection.
For example, converting a pointer to a pointer to an `int` (`int**`) to
a `void*` can result in the loss of information about the original level of
indirection, which can cause problems when attempting to use the converted
pointer. If the converted pointer is later cast to a type with a different
level of indirection and dereferenced, it may lead to access violations,
memory corruption, or other undefined behavior.

Consider the following example:

```c++
void foo(void* ptr);

int main() {
 int x = 42;
 int* ptr = &x;
 int** ptr_ptr = &ptr;
 foo(ptr_ptr); // warning will trigger here
 return 0;
}
```

In this example, `foo()` is called with `ptr_ptr` as its argument. However,
`ptr_ptr` is a `int**` pointer, while `foo()` expects a `void*` pointer.
This results in an implicit pointer level conversion, which could cause issues
if `foo()` dereferences the pointer assuming it's a `int*` pointer.

Using an explicit cast is a recommended solution to prevent issues caused by
implicit pointer level conversion, as it allows the developer to explicitly
state their intention and show their reasoning for the type conversion.
Additionally, it is recommended that developers thoroughly check and verify the
safety of the conversion before using an explicit cast. This extra level of
caution can help catch potential issues early on in the development process,
improving the overall reliability and maintainability of the code.

## Options

```{option} EnableInC
If `true`, enables the check in C code (it is always enabled in C++ code).
Default is `true`.
```

```{title} clang-tidy - bugprone-multiple-new-in-one-expression
```

# bugprone-multiple-new-in-one-expression

Finds multiple `new` operator calls in a single expression, where the
allocated memory by the first `new` may leak if the second allocation fails
and throws exception.

C++ does often not specify the exact order of evaluation of the operands of an
operator or arguments of a function. Therefore if a first allocation succeeds
and a second fails, in an exception handler it is not possible to tell which
allocation has failed and free the memory. Even if the order is fixed the
result of a first `new` may be stored in a temporary location that is not
reachable at the time when a second allocation fails. It is best to avoid any
expression that contains more than one `operator new` call, if exception
handling is used to check for allocation errors.

Different rules apply for are the short-circuit operators `||` and `&&` and
the `,` operator, where evaluation of one side must be completed before the
other starts. Expressions of a list-initialization (initialization or
construction using `{` and `}` characters) are evaluated in fixed order.
Similarly, condition of a `?` operator is evaluated before the branches are
evaluated.

The check reports warning if two `new` calls appear in one expression at
different sides of an operator, or if `new` calls appear in different
arguments of a function call (that can be an object construction with `()`
syntax). These `new` calls can be nested at any level.
For any warning to be emitted the `new` calls should be in a code block where
exception handling is used with catch for `std::bad_alloc` or
`std::exception`. At `||`, `&&`, `,`, `?` (condition and one branch)
operators no warning is emitted. No warning is emitted if both of the memory
allocations are not assigned to a variable or not passed directly to a
function. The reason is that in this case the memory may be intentionally not
freed or the allocated objects can be self-destructing objects.

Examples:

```c++
struct A {
 int Var;
};
struct B {
 B();
 B(A *);
 int Var;
};
struct C {
 int *X1;
 int *X2;
};

void f(A *, B *);
int f1(A *);
int f1(B *);
bool f2(A *);

void foo() {
 A *PtrA;
 B *PtrB;
 try {
 // Allocation of 'B'/'A' may fail after memory for 'A'/'B' was allocated.
 f(new A, new B); // warning: memory allocation may leak if an other allocation is sequenced after it and throws an exception; order of these allocations is undefined

 // List (aggregate) initialization is used.
 C C1{new int, new int}; // no warning

 // Allocation of 'B'/'A' may fail after memory for 'A'/'B' was allocated but not yet passed to function 'f1'.
 int X = f1(new A) + f1(new B); // warning: memory allocation may leak if an other allocation is sequenced after it and throws an exception; order of these allocations is undefined

 // Allocation of 'B' may fail after memory for 'A' was allocated.
 // From C++17 on memory for 'B' is allocated first but still may leak if allocation of 'A' fails.
 PtrB = new B(new A); // warning: memory allocation may leak if an other allocation is sequenced after it and throws an exception

 // 'new A' and 'new B' may be performed in any order.
 // 'new B'/'new A' may fail after memory for 'A'/'B' was allocated but not assigned to 'PtrA'/'PtrB'.
 (PtrA = new A)->Var = (PtrB = new B)->Var; // warning: memory allocation may leak if an other allocation is sequenced after it and throws an exception; order of these allocations is undefined

 // Evaluation of 'f2(new A)' must be finished before 'f1(new B)' starts.
 // If 'new B' fails the allocated memory for 'A' is supposedly handled correctly because function 'f2' could take the ownership.
 bool Z = f2(new A) || f1(new B); // no warning

 X = (f2(new A) ? f1(new A) : f1(new B)); // no warning

 // No warning if the result of both allocations is not passed to a function
 // or stored in a variable.
 (new A)->Var = (new B)->Var; // no warning

 // No warning if at least one non-throwing allocation is used.
 f(new(std::nothrow) A, new B); // no warning
 } catch(std::bad_alloc) {
 }

 // No warning if the allocation is outside a try block (or no catch handler exists for std::bad_alloc).
 // (The fact if exceptions can escape from 'foo' is not taken into account.)
 f(new A, new B); // no warning
}
```

```{title} clang-tidy - bugprone-multiple-statement-macro
```

# bugprone-multiple-statement-macro

Detect multiple statement macros that are used in unbraced conditionals. Only
the first statement of the macro will be inside the conditional and the other
ones will be executed unconditionally.

Example:

```cpp
#define INCREMENT_TWO(x, y) (x)++; (y)++
if (do_increment)
 INCREMENT_TWO(a, b); // (b)++ will be executed unconditionally.
```

```{title} clang-tidy - bugprone-narrowing-conversions
```

# bugprone-narrowing-conversions

`cppcoreguidelines-narrowing-conversions` redirects here as an alias for
this check.

Checks for silent narrowing conversions, e.g: `int i = 0; i += 0.1;`. While
the issue is obvious in this former example, it might not be so in the
following: `void MyClass::f(double d) { int_member_ += d; }`.

We flag narrowing conversions from:

- an integer to a narrower integer (e.g. `char` to `unsigned char`)
 if {option}`WarnOnIntegerNarrowingConversion` is set,
- an integer to a narrower floating-point (e.g. `uint64_t` to `float`)
 if {option}`WarnOnIntegerToFloatingPointNarrowingConversion` is set,
- a floating-point to an integer (e.g. `double` to `int`),
- a floating-point to a narrower floating-point (e.g. `double` to `float`)
 if {option}`WarnOnFloatingPointNarrowingConversion` is set.

This check will flag:

- All narrowing conversions that are not marked by an explicit cast (c-style
 or `static_cast`). For example: `int i = 0; i += 0.1;`,
 `void f(int); f(0.1);`,
- All applications of binary operators with a narrowing conversions.
 For example: `int i; i+= 0.1;`.

Arithmetic with smaller integer types than `int` trigger implicit conversions,
as explained under ["Integral Promotion" on cppreference.com](https://en.cppreference.com/w/cpp/language/implicit_conversion).
This check diagnoses more instances of narrowing than the compiler warning
`-Wconversion` does. The example below demonstrates this behavior.

```c++
// The following function definition demonstrates usage of arithmetic with
// integer types smaller than `int` and how the narrowing conversion happens
// implicitly.
void computation(short argument1, short argument2) {
 // Arithmetic written by humans:
 short result = argument1 + argument2;
 // Arithmetic actually performed by C++:
 short result = static_cast<short>(static_cast<int>(argument1) + static_cast<int>(argument2));
}

void recommended_resolution(short argument1, short argument2) {
 short result = argument1 + argument2;
 // ^ warning: narrowing conversion from 'int' to signed type 'short' is implementation-defined

 // The cppcoreguidelines recommend to resolve this issue by using the GSL
 // in one of two ways. Either by a cast that throws if a loss of precision
 // would occur.
 short result = gsl::narrow<short>(argument1 + argument2);
 // Or it can be resolved without checking the result risking invalid results.
 short result = gsl::narrow_cast<short>(argument1 + argument2);

 // A classical `static_cast` will silence the warning as well if the GSL
 // is not available.
 short result = static_cast<short>(argument1 + argument2);
}
```

## Options

```{option} WarnOnIntegerNarrowingConversion
When `true`, the check will warn on narrowing integer conversion
(e.g. `int` to `size_t`). Default is `true`.
```

```{option} WarnOnIntegerToFloatingPointNarrowingConversion
When `true`, the check will warn on narrowing integer to floating-point
conversion (e.g. `size_t` to `double`). Default is `true`.
```

```{option} WarnOnFloatingPointNarrowingConversion
When `true`, the check will warn on narrowing floating point conversion
(e.g. `double` to `float`). Default is `true`.
```

```{option} WarnWithinTemplateInstantiation
When `true`, the check will warn on narrowing conversions within template
instantiations. Default is `false`.
```

```{option} WarnOnEquivalentBitWidth
When `true`, the check will warn on narrowing conversions that arise from
casting between types of equivalent bit width. (e.g.
`int n = uint(0);` or `long long n = double(0);`) Default is `true`.
```

```{option} IgnoreConversionFromTypes
Narrowing conversions from any type in this semicolon-separated list will be
ignored. This may be useful to weed out commonly occurring, but less commonly
problematic assignments such as `int n = std::vector<char>().size();` or
`int n = std::difference(it1, it2);`. The default list is empty, but one
suggested list for a legacy codebase would be
`size_t;ptrdiff_t;size_type;difference_type`.
```

```{option} PedanticMode
When `true`, the check will warn on assigning a floating point constant
to an integer value even if the floating point value is exactly
representable in the destination type (e.g. `int i = 1.0;`).
Default is `false`.
```

## FAQ

> - What does "narrowing conversion from 'int' to 'float'" mean?

An IEEE754 Floating Point number can represent all integer values in the range
[-2^PrecisionBits, 2^PrecisionBits] where PrecisionBits is the number of bits
in the mantissa.

For `float` this would be [-2^23, 2^23], where `int` can represent values
in the range [-2^31, 2^31-1].

> - What does "implementation-defined" mean?

You may have encountered messages like "narrowing conversion from 'unsigned
int' to signed type 'int' is implementation-defined".
The C/C++ standard does not mandate two's complement for signed integers, and
so the compiler is free to define what the semantics are for converting an
unsigned integer to signed integer. Clang's implementation uses the two's
complement format.

```{title} clang-tidy - bugprone-no-escape
```

# bugprone-no-escape

Finds pointers with the `noescape` attribute that are captured by an
asynchronously-executed block. The block arguments in `dispatch_async()` and
`dispatch_after()` are guaranteed to escape, so it is an error if a pointer
with the `noescape` attribute is captured by one of these blocks.

The following is an example of an invalid use of the `noescape` attribute.

```objc
void foo(__attribute__((noescape)) int *p) {
 dispatch_async(queue, ^{
 *p = 123;
 });
});
```

```{title} clang-tidy - bugprone-non-zero-enum-to-bool-conversion
```

# bugprone-non-zero-enum-to-bool-conversion

Detect implicit and explicit casts of `enum` type into `bool` where
`enum` type doesn't have a zero-value enumerator. If the `enum` is used
only to hold values equal to its enumerators, then conversion to `bool` will
always result in `true` value. This can lead to unnecessary code that reduces
readability and maintainability and can result in bugs.

May produce false positives if the `enum` is used to store other values
(used as a bit-mask or zero-initialized on purpose). To deal with them,
`// NOLINT` or casting first to the underlying type before casting to
`bool` can be used.

It is important to note that this check will not generate warnings if the
definition of the enumeration type is not available.
Additionally, C++11 enumeration classes are supported by this check.

Overall, this check serves to improve code quality and readability by
identifying and flagging instances where implicit or explicit casts from
enumeration types to boolean could cause potential issues.

## Example

```c++
enum EStatus {
 OK = 1,
 NOT_OK,
 UNKNOWN
};

void process(EStatus status) {
 if (!status) {
 // this true-branch won't be executed
 return;
 }
 // proceed with "valid data"
}
```

## Options

```{option} EnumIgnoreList
Option is used to ignore certain enum types when checking for
implicit/explicit casts to bool. It accepts a semicolon-separated list of
(fully qualified) enum type names or regular expressions that match the enum
type names.
The default value is an empty string, which means no enums will be ignored.
```

```{title} clang-tidy - bugprone-nondeterministic-pointer-iteration-order
```

# bugprone-nondeterministic-pointer-iteration-order

Finds nondeterministic usages of pointers in unordered containers.

One canonical example is iteration across a container of pointers.

```c++
{
 int a = 1, b = 2;
 std::unordered_set<int *> UnorderedPtrSet = {&a, &b};
 for (auto i : UnorderedPtrSet)
 f(i);
}
```

Another such example is sorting a container of pointers.

```c++
{
 int a = 1, b = 2;
 std::vector<int *> VectorOfPtr = {&a, &b};
 std::sort(VectorOfPtr.begin(), VectorOfPtr.end());
}
```

Iteration of a containers of pointers may present the order of different
pointers differently across different runs of a program. In some cases this
may be acceptable behavior, in others this may be unexpected behavior. This
check is advisory for this reason.

This check only detects range-based for loops over unordered sets and maps. It
also detects calls sorting-like algorithms on containers holding pointers.
Other similar usages will not be found and are false negatives.

## Limitations

- This check currently does not check if a nondeterministic iteration order is
 likely to be a mistake, and instead marks all such iterations as bugprone.
- std::reference_wrapper is not considered yet.
- Only for loops are considered, other iterators can be included in
 improvements.

```{title} clang-tidy - bugprone-not-null-terminated-result
```

# bugprone-not-null-terminated-result

Finds function calls where it is possible to cause a not null-terminated
result. Usually the proper length of a string is `strlen(src) + 1` or equal
length of this expression, because the null terminator needs an extra space.
Without the null terminator it can result in undefined behavior when the
string is read.

The following and their respective `wchar_t` based functions are checked:

`memcpy`, `memcpy_s`, `memchr`, `memmove`, `memmove_s`,
`strerror_s`, `strncmp`, `strxfrm`

The following is a real-world example where the programmer forgot to increase
the passed third argument, which is `size_t length`. That is why the length
of the allocated memory is not enough to hold the null terminator.

```c
static char *stringCpy(const std::string &str) {
 char *result = reinterpret_cast<char *>(malloc(str.size()));
 memcpy(result, str.data(), str.size());
 return result;
}
```

In addition to issuing warnings, fix-it rewrites all the necessary code.
It also tries to adjust the capacity of the destination array:

```c
static char *stringCpy(const std::string &str) {
 char *result = reinterpret_cast<char *>(malloc(str.size() + 1));
 strcpy(result, str.data());
 return result;
}
```

Note: It cannot guarantee to rewrite every of the path-sensitive memory
allocations.

(memcpytransformation)=

## Transformation rules of 'memcpy()'

It is possible to rewrite the `memcpy()` and `memcpy_s()` calls as the
following four functions: `strcpy()`, `strncpy()`, `strcpy_s()`,
`strncpy_s()`, where the latter two are the safer versions of the former two.
It rewrites the `wchar_t` based memory handler functions respectively.

### Rewrite based on the destination array

- If copy to the destination array cannot overflow [1] the new function should
 be the older copy function (ending with `cpy`), because it is more
 efficient than the safe version.
- If copy to the destination array can overflow [1] and
 {option}`WantToUseSafeFunctions` is set to `true` and it is
 possible to
 obtain the capacity of the destination array then the new function could be
 the safe version (ending with `cpy_s`).
- If the new function is could be safe version and C++ files are analyzed and
 the destination array is plain `char`/`wchar_t` without `un/signed`
 then the length of the destination array can be omitted.
- If the new function is could be safe version and the destination array is
 `un/signed` it needs to be casted to plain `char *`/`wchar_t *`.

[1] It is possible to overflow:

- If the capacity of the destination array is unknown.
- If the given length is equal to the destination array's capacity.

### Rewrite based on the length of the source string

- If the given length is `strlen(source)` or equal length of this expression
 then the new function should be the older copy function (ending with
 `cpy`), as it is more efficient than the safe version (ending with
 `cpy_s`).
- Otherwise we assume that the programmer wanted to copy 'N' characters, so the
 new function is `ncpy`-like which copies 'N' characters.

## Transformations with 'strlen()' or equal length of this expression

It transforms the `wchar_t` based memory and string handler functions
respectively (where only `strerror_s` does not have `wchar_t` based alias).

### Memory handler functions

`memcpy`
Please visit the
{ref}`Transformation rules of 'memcpy()'<MemcpyTransformation>` section.

`memchr`
Usually there is a C-style cast and it is needed to be removed, because the
new function `strchr`'s return type is correct. The given length is going
to be removed.

`memmove`
If safe functions are available the new function is `memmove_s`, which has
a new second argument which is the length of the destination array, it is
adjusted, and the length of the source string is incremented by one.
If safe functions are not available the given length is incremented by one.

`memmove_s`
The given length is incremented by one.

### String handler functions

`strerror_s`
The given length is incremented by one.

`strncmp`
If the third argument is the first or the second argument's `length + 1`
it has to be truncated without the `+ 1` operation.

`strxfrm`
The given length is incremented by one.

## Options

```{option} WantToUseSafeFunctions
The value `true` specifies that the target environment is considered to
implement '\_s' suffixed memory and string handler functions which are safer
than older versions (e.g. 'memcpy_s()'). The default value is `true`.
```

```{title} clang-tidy - bugprone-optional-value-conversion
```

# bugprone-optional-value-conversion

Detects potentially unintentional and redundant conversions where a value is
extracted from an optional-like type and then used to create a new instance of
the same optional-like type.

These conversions might be the result of developer oversight, leftovers from
code refactoring, or other situations that could lead to unintended exceptions
or cases where the resulting optional is always initialized, which might be
unexpected behavior.

To illustrate, consider the following problematic code snippet:

```c++
#include <optional>

void print(std::optional<int>);

int main()
{
 std::optional<int> opt;
 // ...

 // Unintentional conversion from std::optional<int> to int and back to
 // std::optional<int>:
 print(opt.value());

 // ...
}
```

A better approach would be to directly pass `opt` to the `print` function
without extracting its value:

```c++
#include <optional>

void print(std::optional<int>);

int main()
{
 std::optional<int> opt;
 // ...

 // Proposed code: Directly pass the std::optional<int> to the print
 // function.
 print(opt);

 // ...
}
```

By passing `opt` directly to the print function, unnecessary conversions are
avoided, and potential unintended behavior or exceptions are minimized.

Value extraction using `operator *` is matched by default.
The support for non-standard optional types such as `boost::optional` or
`absl::optional` may be limited.

## Options:

```{option} OptionalTypes
Semicolon-separated list of (fully qualified) optional type names or regular
expressions that match the optional types.
Default value is `::std::optional;::absl::optional;::boost::optional`.
```

```{option} ValueMethods
Semicolon-separated list of (fully qualified) method names or regular
expressions that match the methods.
Default value is `::value$;::get$`.
```

```{title} clang-tidy - bugprone-parent-virtual-call
```

# bugprone-parent-virtual-call

Detects and fixes calls to grand-...parent virtual methods instead of calls
to overridden parent's virtual methods.

```cpp
struct A {
 int virtual foo() {...}
};

struct B: public A {
 int foo() override {...}
};

struct C: public B {
 int foo() override { A::foo(); }
// ^^^^^^^^
// warning: qualified name A::foo refers to a member overridden in subclass; did you mean 'B'? [bugprone-parent-virtual-call]
};
```

```{title} clang-tidy - bugprone-pointer-arithmetic-on-polymorphic-object
```

# bugprone-pointer-arithmetic-on-polymorphic-object

Finds pointer arithmetic performed on classes that contain a virtual function.

Pointer arithmetic on polymorphic objects where the pointer's static type is
different from its dynamic type is undefined behavior, as the two types could
have different sizes, and thus the vtable pointer could point to an
invalid address.

Finding pointers where the static type contains a virtual member function is a
good heuristic, as the pointer is likely to point to a different,
derived object.

Example:

```c++
struct Base {
 virtual ~Base();
 int i;
};

struct Derived : public Base {};

void foo(Base* b) {
 b += 1;
 // warning: pointer arithmetic on class that declares a virtual function can
 // result in undefined behavior if the dynamic type differs from the
 // pointer type
}

int bar(const Derived d[]) {
 return d[1].i; // warning due to pointer arithmetic on polymorphic object
}

// Making Derived final suppresses the warning
struct FinalDerived final : public Base {};

int baz(const FinalDerived d[]) {
 return d[1].i; // no warning as FinalDerived is final
}
```

## Options

````{option} IgnoreInheritedVirtualFunctions
When `true`, objects that only inherit a virtual function are not checked.
Classes that do not declare a new virtual function are excluded
by default, as they make up the majority of false positives.
Default is `false`.

```c++
void bar(Base b[], Derived d[]) {
 b += 1; // warning, as Base declares a virtual destructor
 d += 1; // warning only if IgnoreVirtualDeclarationsOnly is set to false
}
```
````

## References

This check corresponds to the SEI Cert rule
[CTR56-CPP. Do not use pointer arithmetic on polymorphic objects](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/containers-ctr/ctr56-cpp/).

```{title} clang-tidy - bugprone-posix-return
```

# bugprone-posix-return

Checks if any calls to `pthread_*` or `posix_*` functions
(except `posix_openpt`) expect negative return values. These functions return
either `0` on success or an `errno` on failure, which is positive only.

Example buggy usage looks like:

```c
if (posix_fadvise(...) < 0) {
```

This will never happen as the return value is always non-negative.
A simple fix could be:

```c
if (posix_fadvise(...) > 0) {
```

```{title} clang-tidy - bugprone-random-generator-seed
```

# bugprone-random-generator-seed

Flags all pseudo-random number engines, engine adaptor
instantiations and `srand()` when initialized or seeded with default
argument, constant expression or any user-configurable type. Pseudo-random
number engines seeded with a predictable value may cause vulnerabilities
e.g. in security protocols.

Examples:

```c++
void foo() {
 std::mt19937 engine1; // Diagnose, always generate the same sequence
 std::mt19937 engine2(1); // Diagnose
 engine1.seed(); // Diagnose
 engine2.seed(1); // Diagnose

 std::time_t t;
 engine1.seed(std::time(&t)); // Diagnose, system time might be controlled by user

 int x = atoi(argv[1]);
 std::mt19937 engine3(x); // Will not warn
}
```

## Options

```{option} DisallowedSeedTypes
A comma-separated list of the type names which are disallowed.
Default is `time_t,std::time_t`.
```

## References

This check corresponds to the CERT C++ Coding Standard rules
[MSC51-CPP. Ensure your random number generator is properly seeded](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/miscellaneous-msc/msc51-cpp/) and
[MSC32-C. Properly seed pseudorandom number generators](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/miscellaneous-msc/msc32-c/).

```{title} clang-tidy - bugprone-raw-memory-call-on-non-trivial-type
```

# bugprone-raw-memory-call-on-non-trivial-type

Flags use of the C standard library functions `memset`, `memcpy` and
`memcmp` and similar derivatives on non-trivial types.

The check will detect the following functions: `memset`, `std::memset`,
`std::memcpy`, `memcpy`, `std::memmove`, `memmove`, `std::strcpy`,
`strcpy`, `memccpy`, `stpncpy`, `strncpy`, `std::memcmp`, `memcmp`,
`std::strcmp`, `strcmp`, `strncmp`.

## Options

```{option} MemSetNames
Specify extra functions to flag that act similarly to `memset`. Specify
names in a semicolon-delimited list. Default is an empty string.
```

```{option} MemCpyNames
Specify extra functions to flag that act similarly to `memcpy`. Specify
names in a semicolon-delimited list. Default is an empty string.
```

```{option} MemCmpNames
Specify extra functions to flag that act similarly to `memcmp`. Specify
names in a semicolon-delimited list. Default is an empty string.
```

This check corresponds to the CERT C++ Coding Standard rule
[OOP57-CPP. Prefer special member functions and overloaded operators to C
Standard Library functions](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/object-oriented-programming-oop/oop57-cpp/).

```{title} clang-tidy - bugprone-redundant-branch-condition
```

# bugprone-redundant-branch-condition

Finds condition variables in nested `if` statements that were also checked in
the outer `if` statement and were not changed.

Simple example:

```c
bool onFire = isBurning();
if (onFire) {
 if (onFire)
 scream();
}
```

Here `onFire` is checked both in the outer `if` and the inner `if`
statement without a possible change between the two checks. The check warns for
this code and suggests removal of the second checking of variable
`onFire`.

The check also detects redundant condition checks if the condition variable
is an operand of a logical "and" (`&&`) or a logical "or" (`||`) operator:

```c
bool onFire = isBurning();
if (onFire) {
 if (onFire && peopleInTheBuilding > 0)
 scream();
}
```

```c
bool onFire = isBurning();
if (onFire) {
 if (onFire || isCollapsing())
 scream();
}
```

In the first case (logical "and") the suggested fix is to remove the redundant
condition variable and keep the other side of the `&&`. In the second case
(logical "or") the whole `if` is removed similarly to the simple case on the
top.

The condition of the outer `if` statement may also be a logical "and"
(`&&`) expression:

```c
bool onFire = isBurning();
if (onFire && fireFighters < 10) {
 if (someOtherCondition()) {
 if (onFire)
 scream();
 }
}
```

The error is also detected if both the outer statement is a logical "and"
(`&&`) and the inner statement is a logical "and" (`&&`) or "or" (`||`).
The inner `if` statement does not have to be a direct descendant of the outer
one.

No error is detected if the condition variable may have been changed between
the two checks:

```c
bool onFire = isBurning();
if (onFire) {
 tryToExtinguish(onFire);
 if (onFire && peopleInTheBuilding > 0)
 scream();
}
```

Every possible change is considered, thus if the condition variable is not
a local variable of the function, it is a volatile or it has an alias (pointer
or reference) then no warning is issued.

## Limitations

The `else` branch is not checked currently for negated condition variable:

```c
bool onFire = isBurning();
if (onFire) {
 scream();
} else {
 if (!onFire) {
 continueWork();
 }
}
```

The check currently only detects redundant checking of single condition
variables. More complex expressions are not checked:

```c
if (peopleInTheBuilding == 1) {
 if (peopleInTheBuilding == 1) {
 doSomething();
 }
}
```

```{title} clang-tidy - bugprone-reserved-identifier
```

# bugprone-reserved-identifier

`cert-dcl37-c` and `cert-dcl51-cpp` redirect
here as an alias for this check.

Checks for usages of identifiers reserved for use by the implementation.

The C and C++ standards both reserve the following names for such use:

- identifiers that begin with an underscore followed by an uppercase letter;
- identifiers in the global namespace that begin with an underscore.

The C standard additionally reserves names beginning with a double underscore,
while the C++ standard strengthens this to reserve names with a double
underscore occurring anywhere.

Violating the naming rules above results in undefined behavior.

```c++
namespace NS {
 void __f(); // name is not allowed in user code
 using _Int = int; // same with this
 #define cool__macro // also this
}
int _g(); // disallowed in global namespace only
```

The check can also be inverted, i.e. it can be configured to flag any
identifier that is *not* a reserved identifier. This mode is for use by e.g.
standard library implementors, to ensure they don't infringe on the user
namespace.

This check does not (yet) check for other reserved names, e.g. macro names
identical to language keywords, and names specifically reserved by language
standards, e.g. C++ 'zombie names' and C future library directions.

This check corresponds to CERT C Coding Standard rule [DCL37-C. Do not declare
or define a reserved identifier](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/declarations-and-initialization-dcl/dcl37-c/)
as well as its C++ counterpart, [DCL51-CPP. Do not declare or define a reserved
identifier](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/declarations-and-initialization-dcl/dcl51-cpp/).

## Options

```{option} Invert
If `true`, inverts the check, i.e. flags names that are not reserved.
Default is `false`.
```

```{option} AllowedIdentifiers
Semicolon-separated list of regular expressions that the check ignores. Default is an
empty string.
```

```{title} clang-tidy - bugprone-return-const-ref-from-parameter
```

# bugprone-return-const-ref-from-parameter

Detects return statements that return a constant reference parameter as
constant reference. This may cause use-after-free errors if the caller
uses xvalues as arguments.

In C++, constant reference parameters can accept xvalues which will be
destructed after the call. When the function returns such a parameter also
as constant reference, then the returned reference can be used after the
object it refers to has been destroyed.

## Example

```c++
struct S {
 int v;
 S(int);
 ~S();
};

const S &fn(const S &a) {
 return a;
}

const S& s = fn(S{1});
s.v; // use after free
```

This issue can be resolved by declaring an overload of the problematic function
where the `const &` parameter is instead declared as `&&`. The developer has
to ensure that the implementation of that function does not produce a
use-after-free, the exact error that this check is warning against.
Marking such an `&&` overload as `deleted`, will silence the warning as
well. In the case of different `const &` parameters being returned depending
on the control flow of the function, an overload where all problematic
`const &` parameters have been declared as `&&` will resolve the issue.

This issue can also be resolved by adding `[[clang::lifetimebound]]`. Clang
enable `-Wdangling` warning by default which can detect mis-uses of the
annotated function. See [lifetimebound attribute](https://clang.llvm.org/docs/AttributeReference.html#lifetimebound)
for details.

```c++
const int &f(const int &a [[clang::lifetimebound]]) { return a; } // no warning
const int &v = f(1); // warning: temporary bound to local reference 'v' will be destroyed at the end of the full-expression [-Wdangling]
```

```{title} clang-tidy - bugprone-shared-ptr-array-mismatch
```

# bugprone-shared-ptr-array-mismatch

Finds initializations of C++ shared pointers to non-array type that are
initialized with an array.

If a shared pointer `std::shared_ptr<T>` is initialized with a new-expression
`new T[]` the memory is not deallocated correctly. The pointer uses plain
`delete` in this case to deallocate the target memory. Instead a `delete[]`
call is needed. A `std::shared_ptr<T[]>` calls the correct delete operator.

The check offers replacement of `shared_ptr<T>` to `shared_ptr<T[]>` if it
is used at a single variable declaration (one variable in one statement).

Example:

```c++
std::shared_ptr<Foo> x(new Foo[10]); // -> std::shared_ptr<Foo[]> x(new Foo[10]);
// ^ warning: shared pointer to non-array is initialized with array [bugprone-shared-ptr-array-mismatch]
std::shared_ptr<Foo> x1(new Foo), x2(new Foo[10]); // no replacement
// ^ warning: shared pointer to non-array is initialized with array [bugprone-shared-ptr-array-mismatch]

std::shared_ptr<Foo> x3(new Foo[10], [](const Foo *ptr) { delete[] ptr; }); // no warning

struct S {
 std::shared_ptr<Foo> x(new Foo[10]); // no replacement in this case
 // ^ warning: shared pointer to non-array is initialized with array [bugprone-shared-ptr-array-mismatch]
};
```

This check partially covers the CERT C++ Coding Standard rule
[MEM51-CPP. Properly deallocate dynamically allocated resources](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/memory-management-mem/mem51-cpp/)
However, only the `std::shared_ptr` case is detected by this check.

```{title} clang-tidy - bugprone-signal-handler
```

# bugprone-signal-handler

Finds specific constructs in signal handler functions that can cause undefined
behavior. The rules for what is allowed differ between C++ language versions.

Checked signal handler rules for C:

- Calls to non-asynchronous-safe functions are not allowed.

Checked signal handler rules for up to and including C++14:

- Calls to non-asynchronous-safe functions are not allowed.
- C++-specific code constructs are not allowed in signal handlers.
 In other words, only the common subset of C and C++ is allowed to be used.
- Calls to functions with non-C linkage are not allowed (including the signal
 handler itself).

The check is disabled on C++17 and later.

Asynchronous-safety is determined by comparing the function's name against a
set of known functions. In addition, the function must come from a system
header include and in a global namespace. The (possible) arguments passed to
the function are not checked. Any function that cannot be determined to be
asynchronous-safe is assumed to be non-asynchronous-safe by the check,
including user functions for which only the declaration is visible.
Calls to user-defined functions with visible definitions are checked
recursively.

This check implements the CERT C Coding Standard rule
[SIG30-C. Call only asynchronous-safe functions within signal handlers](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/signals-sig/sig30-c/)
and the rule
[MSC54-CPP. A signal handler must be a plain old function](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/miscellaneous-msc/msc54-cpp/).
It has the alias names `cert-sig30-c` and `cert-msc54-cpp`.

## Options

```{option} AsyncSafeFunctionSet
Selects which set of functions is considered as asynchronous-safe
(and therefore allowed in signal handlers). It can be set to the following values:

- `minimal`
 : Selects a minimal set that is defined in the CERT SIG30-C rule.
 and includes functions `abort()`, `_Exit()`, `quick_exit()` and
 `signal()`.
- `POSIX`
 : Selects a larger set of functions that is listed in POSIX.1-2017 (see [this
 link](https://pubs.opengroup.org/onlinepubs/9699919799/functions/V2_chap02.html#tag_15_04_03)
 for more information). The following functions are included:
 `_Exit`, `_exit`, `abort`, `accept`, `access`, `aio_error`,
 `aio_return`, `aio_suspend`, `alarm`, `bind`, `cfgetispeed`,
 `cfgetospeed`, `cfsetispeed`, `cfsetospeed`, `chdir`, `chmod`,
 `chown`, `clock_gettime`, `close`, `connect`, `creat`, `dup`,
 `dup2`, `execl`, `execle`, `execv`, `execve`, `faccessat`,
 `fchdir`, `fchmod`, `fchmodat`, `fchown`, `fchownat`, `fcntl`,
 `fdatasync`, `fexecve`, `ffs`, `fork`, `fstat`, `fstatat`,
 `fsync`, `ftruncate`, `futimens`, `getegid`, `geteuid`,
 `getgid`, `getgroups`, `getpeername`, `getpgrp`, `getpid`,
 `getppid`, `getsockname`, `getsockopt`, `getuid`, `htonl`,
 `htons`, `kill`, `link`, `linkat`, `listen`, `longjmp`,
 `lseek`, `lstat`, `memccpy`, `memchr`, `memcmp`, `memcpy`,
 `memmove`, `memset`, `mkdir`, `mkdirat`, `mkfifo`, `mkfifoat`,
 `mknod`, `mknodat`, `ntohl`, `ntohs`, `open`, `openat`,
 `pause`, `pipe`, `poll`, `posix_trace_event`, `pselect`,
 `pthread_kill`, `pthread_self`, `pthread_sigmask`, `quick_exit`,
 `raise`, `read`, `readlink`, `readlinkat`, `recv`, `recvfrom`,
 `recvmsg`, `rename`, `renameat`, `rmdir`, `select`, `sem_post`,
 `send`, `sendmsg`, `sendto`, `setgid`, `setpgid`, `setsid`,
 `setsockopt`, `setuid`, `shutdown`, `sigaction`, `sigaddset`,
 `sigdelset`, `sigemptyset`, `sigfillset`, `sigismember`,
 `siglongjmp`, `signal`, `sigpause`, `sigpending`, `sigprocmask`,
 `sigqueue`, `sigset`, `sigsuspend`, `sleep`, `sockatmark`,
 `socket`, `socketpair`, `stat`, `stpcpy`, `stpncpy`,
 `strcat`, `strchr`, `strcmp`, `strcpy`, `strcspn`, `strlen`,
 `strncat`, `strncmp`, `strncpy`, `strnlen`, `strpbrk`,
 `strrchr`, `strspn`, `strstr`, `strtok_r`, `symlink`,
 `symlinkat`, `tcdrain`, `tcflow`, `tcflush`, `tcgetattr`,
 `tcgetpgrp`, `tcsendbreak`, `tcsetattr`, `tcsetpgrp`,
 `time`, `timer_getoverrun`, `timer_gettime`, `timer_settime`,
 `times`, `umask`, `uname`, `unlink`, `unlinkat`, `utime`,
 `utimensat`, `utimes`, `wait`, `waitpid`, `wcpcpy`,
 `wcpncpy`, `wcscat`, `wcschr`, `wcscmp`, `wcscpy`, `wcscspn`,
 `wcslen`, `wcsncat`, `wcsncmp`, `wcsncpy`, `wcsnlen`, `wcspbrk`,
 `wcsrchr`, `wcsspn`, `wcsstr`, `wcstok`, `wmemchr`, `wmemcmp`,
 `wmemcpy`, `wmemmove`, `wmemset`, `write`

 The function `quick_exit` is not included in the POSIX list but it
 is included here in the set of safe functions.

Default is `POSIX`.
```

```{title} clang-tidy - bugprone-signed-bitwise
```

# bugprone-signed-bitwise

Finds uses of bitwise operations on signed integer types, which may lead to
undefined or implementation defined behavior.

Performing bitwise operations on signed integers can be confusing and may
lead to undefined or implementation dependent behavior. In particular, right
shift a signed integer is implementation dependent, while left shift a signed
integer may result in undefined behavior.

```cpp
int main(){
 int x = -4;
 int y = x >> 1; // y can be -2 or 2147483646
}
```

## Options

```{option} IgnorePositiveIntegerLiterals

If this option is set to `true`, the check will not warn on bitwise
operations with positive integer literals, e.g. `~0`, `2 << 1`, etc.
Default value is `false`.
```

```{title} clang-tidy - bugprone-signed-char-misuse
```

# bugprone-signed-char-misuse

`cert-str34-c` redirects here as an alias for this check. For
the CERT alias, the `DiagnoseSignedUnsignedCharComparisons`
option is set to `false`.

Finds those `signed char` -> integer conversions which might indicate a
programming error. The basic problem with the `signed char`, that it might
store the non-ASCII characters as negative values. This behavior can cause a
misunderstanding of the written code both when an explicit and when an
implicit conversion happens.

When the code contains an explicit `signed char` -> integer conversion, the
human programmer probably expects that the converted value matches with the
character code (a value from [0..255]), however, the actual value is in
[-128..127] interval. To avoid this kind of misinterpretation, the desired way
of converting from a `signed char` to an integer value is converting to
`unsigned char` first, which stores all the characters in the positive
[0..255] interval which matches the known character codes.

In case of implicit conversion, the programmer might not actually be aware
that a conversion happened and char value is used as an integer. There are
some use cases when this unawareness might lead to a functionally imperfect
code. For example, checking the equality of a `signed char` and an
`unsigned char` variable is something we should avoid in C++ code. During
this comparison, the two variables are converted to integers which have
different value ranges. For `signed char`, the non-ASCII characters are
stored as a value in [-128..-1] interval, while the same characters are
stored in the [128..255] interval for an `unsigned char`.

It depends on the actual platform whether plain `char` is handled as
`signed char` by default and so it is caught by this check or not.
To change the default behavior you can use `-funsigned-char` and
`-fsigned-char` compilation options.

Currently, this check warns in the following cases:

- `signed char` is assigned to an integer variable
- `signed char` and `unsigned char` are compared with
 equality/inequality operator
- `signed char` is converted to an integer in the array subscript

See also:
[STR34-C. Cast characters to unsigned char before converting to larger
integer sizes](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/characters-and-strings-str/str34-c/)

A good example from the CERT description when a `char` variable is used to
read from a file that might contain non-ASCII characters. The problem comes
up when the code uses the `-1` integer value as EOF, while the 255 character
code is also stored as `-1` in two's complement form of char type.
See a simple example of this below. This code stops not only when it reaches
the end of the file, but also when it gets a character with the 255 code.

```c++
#define EOF (-1)

int read(void) {
 char CChar;
 int IChar = EOF;

 if (readChar(CChar)) {
 IChar = CChar;
 }
 return IChar;
}
```

A proper way to fix the code above is converting the `char` variable to
an `unsigned char` value first.

```c++
#define EOF (-1)

int read(void) {
 char CChar;
 int IChar = EOF;

 if (readChar(CChar)) {
 IChar = static_cast<unsigned char>(CChar);
 }
 return IChar;
}
```

Another use case is checking the equality of two `char` variables with
different signedness. Inside the non-ASCII value range this comparison between
a `signed char` and an `unsigned char` always returns `false`.

```c++
bool compare(signed char SChar, unsigned char USChar) {
 if (SChar == USChar)
 return true;
 return false;
}
```

The easiest way to fix this kind of comparison is casting one of the arguments,
so both arguments will have the same type.

```c++
bool compare(signed char SChar, unsigned char USChar) {
 if (static_cast<unsigned char>(SChar) == USChar)
 return true;
 return false;
}
```

## Options

```{option} CharTypedefsToIgnore
A semicolon-separated list of typedef names. In this list, we can list
typedefs for `char` or `signed char`, which will be ignored by the
check. This is useful when a typedef introduces an integer alias like
`sal_Int8` or `int8_t`. In this case, human misinterpretation is not
an issue. Default is an empty string.
```

```{option} DiagnoseSignedUnsignedCharComparisons
When `true`, the check will warn on `signed char`/`unsigned char` comparisons,
otherwise these comparisons are ignored. Default is `true`.
```

```{title} clang-tidy - bugprone-sizeof-container
```

# bugprone-sizeof-container

The check finds usages of `sizeof` on expressions of STL container types.
Most likely the user wanted to use `.size()` instead.

All class/struct types declared in namespace `std::` having a const
`size()` method are considered containers, with the exception of
`std::bitset` and `std::array`.

Examples:

```cpp
std::string s;
int a = 47 + sizeof(s); // warning: sizeof() doesn't return the size of the container. Did you mean .size()?

int b = sizeof(std::string); // no warning, probably intended.

std::string array_of_strings[10];
int c = sizeof(array_of_strings) / sizeof(array_of_strings[0]); // no warning, definitely intended.

std::array<int, 3> std_array;
int d = sizeof(std_array); // no warning, probably intended.
```

```{title} clang-tidy - bugprone-sizeof-expression
```

# bugprone-sizeof-expression

The check finds usages of `sizeof` expressions which are most likely errors.

The `sizeof` operator yields the size (in bytes) of its operand, which may be
an expression or the parenthesized name of a type. Misuse of this operator may
be leading to errors and possible software vulnerabilities.

## Suspicious usage of 'sizeof(K)'

A common mistake is to query the `sizeof` of an integer literal. This is
equivalent to query the size of its type (probably `int`). The intent of the
programmer was probably to simply get the integer and not its size.

```c++
#define BUFLEN 42
char buf[BUFLEN];
memset(buf, 0, sizeof(BUFLEN)); // sizeof(42) ==> sizeof(int)
```

## Suspicious usage of 'sizeof(expr)'

In cases, where there is an enum or integer to represent a type, a common
mistake is to query the `sizeof` on the integer or enum that represents the
type that should be used by `sizeof`. This results in the size of the integer
and not of the type the integer represents:

```c++
enum data_type {
 FLOAT_TYPE,
 DOUBLE_TYPE
};

struct data {
 data_type type;
 void* buffer;
 data_type get_type() {
 return type;
 }
};

void f(data d, int numElements) {
 // should be sizeof(float) or sizeof(double), depending on d.get_type()
 int numBytes = numElements * sizeof(d.get_type());
 ...
}
```

## Suspicious usage of 'sizeof(this)'

The `this` keyword is evaluated to a pointer to an object of a given type.
The expression `sizeof(this)` is returning the size of a pointer. The
programmer most likely wanted the size of the object and not the size of the
pointer.

```c++
class Point {
 [...]
 size_t size() { return sizeof(this); } // should probably be sizeof(*this)
 [...]
};
```

## Suspicious usage of 'sizeof(char\*)'

There is a subtle difference between declaring a string literal with
`char* A = ""` and `char A[] = ""`. The first case has the type `char*`
instead of the aggregate type `char[]`. Using `sizeof` on an object
declared with `char*` type is returning the size of a pointer instead of
the number of characters (bytes) in the string literal.

```c++
const char* kMessage = "Hello World!"; // const char kMessage[] = "...";
void getMessage(char* buf) {
 memcpy(buf, kMessage, sizeof(kMessage)); // sizeof(char*)
}
```

## Suspicious usage of 'sizeof(A\*)'

A common mistake is to compute the size of a pointer instead of its pointee.
These cases may occur because of explicit cast or implicit conversion.

```c++
int A[10];
memset(A, 0, sizeof(A + 0));

struct Point point;
memset(point, 0, sizeof(&point));
```

## Suspicious usage of 'sizeof(...)/sizeof(...)'

Dividing `sizeof` expressions is typically used to retrieve the number of
elements of an aggregate. This check warns on incompatible or suspicious cases.

In the following example, the entity has 10-bytes and is incompatible with the
type `int` which has 4 bytes.

```c++
char buf[] = { 0, 1, 2, 3, 4, 5, 6, 7, 8, 9 }; // sizeof(buf) => 10
void getMessage(char* dst) {
 memcpy(dst, buf, sizeof(buf) / sizeof(int)); // sizeof(int) => 4 [incompatible sizes]
}
```

In the following example, the expression `sizeof(Values)` is returning the
size of `char*`. One can easily be fooled by its declaration, but in parameter
declaration the size '10' is ignored and the function is receiving a `char*`.

```c++
char OrderedValues[10] = { 0, 1, 2, 3, 4, 5, 6, 7, 8, 9 };
return CompareArray(char Values[10]) {
 return memcmp(OrderedValues, Values, sizeof(Values)) == 0; // sizeof(Values) ==> sizeof(char*) [implicit cast to char*]
}
```

## Suspicious 'sizeof' by 'sizeof' expression

Multiplying `sizeof` expressions typically makes no sense and is probably a
logic error. In the following example, the programmer used `*` instead of
`/`.

```c++
const char kMessage[] = "Hello World!";
void getMessage(char* buf) {
 memcpy(buf, kMessage, sizeof(kMessage) * sizeof(char)); // sizeof(kMessage) / sizeof(char)
}
```

This check may trigger on code using the arraysize macro. The following code is
working correctly but should be simplified by using only the `sizeof`
operator.

```c++
extern Object objects[100];
void InitializeObjects() {
 memset(objects, 0, arraysize(objects) * sizeof(Object)); // sizeof(objects)
}
```

## Suspicious usage of 'sizeof(sizeof(...))'

Getting the `sizeof` of a `sizeof` makes no sense and is typically an error
hidden through macros.

```c++
#define INT_SZ sizeof(int)
int buf[] = { 42 };
void getInt(int* dst) {
 memcpy(dst, buf, sizeof(INT_SZ)); // sizeof(sizeof(int)) is suspicious.
}
```

## Suspicious usages of 'sizeof(...)' in pointer arithmetic

Arithmetic operators on pointers automatically scale the result with the size
of the pointed typed.
Further use of `sizeof` around pointer arithmetic will typically result in an
unintended result.

### Scaling the result of pointer difference

Subtracting two pointers results in an integer expression (of type
`ptrdiff_t`) which expresses the distance between the two pointed objects in
"number of objects between".
A common mistake is to think that the result is "number of bytes between", and
scale the difference with `sizeof`, such as `P1 - P2 == N * sizeof(T)`
(instead of `P1 - P2 == N`) or `(P1 - P2) / sizeof(T)` instead of
`P1 - P2`.

```c++
void splitFour(const Obj* Objs, size_t N, Obj Delimiter) {
 const Obj *P = Objs;
 while (P < Objs + N) {
 if (*P == Delimiter) {
 break;
 }
 }

 if (P - Objs != 4 * sizeof(Obj)) { // Expecting a distance multiplied by sizeof is suspicious.
 error();
 }
}
```

```c++
void iterateIfEvenLength(int *Begin, int *End) {
 auto N = (Begin - End) / sizeof(int); // Dividing by sizeof() is suspicious.
 if (N % 2)
 return;

 // ...
}
```

### Stepping a pointer with a scaled integer

Conversely, when performing pointer arithmetics to add or subtract from a
pointer, the arithmetic operator implicitly scales the value actually added to
the pointer with the size of the pointee, as `Ptr + N` expects `N` to be
"number of objects to step", and not "number of bytes to step".

Seeing the calculation of a pointer where `sizeof` appears is suspicious,
and the result is typically unintended, often out of bounds.
`Ptr + sizeof(T)` will offset the pointer by `sizeof(T)` elements,
effectively exponentiating the scaling factor to the power of 2.

Similarly, multiplying or dividing a numeric value with the `sizeof` of an
element or the whole buffer is suspicious, because the dimensional connection
between the numeric value and the actual `sizeof` result can not always be
deduced.
While scaling an integer up (multiplying) with `sizeof` is likely **always**
an issue, a scaling down (division) is not always inherently dangerous, in case
the developer is aware that the division happens between an appropriate number
of \_bytes\_ and a `sizeof` value.
Turning {option}`WarnOnOffsetDividedBySizeOf` off will restrict the
warnings to the multiplication case.

This case also checks suspicious `alignof` and `offsetof` usages in
pointer arithmetic, as both return the "size" in bytes and not elements,
potentially resulting in doubly-scaled offsets.

```c++
void printEveryEvenIndexElement(int *Array, size_t N) {
 int *P = Array;
 while (P <= Array + N * sizeof(int)) { // Suspicious pointer arithmetic using sizeof()!
 printf("%d ", *P);

 P += 2 * sizeof(int); // Suspicious pointer arithmetic using sizeof()!
 }
}
```

```c++
struct Message { /* ... */; char Flags[8]; };
void clearFlags(Message *Array, size_t N) {
 const Message *End = Array + N;
 while (Array < End) {
 memset(Array + offsetof(Message, Flags), // Suspicious pointer arithmetic using offsetof()!
 0, sizeof(Message::Flags));
 ++Array;
 }
}
```

For this checked bogus pattern, `cert-arr39-c` redirects here as an alias of
this check.

This check corresponds to the CERT C Coding Standard rule
[ARR39-C. Do not add or subtract a scaled integer to a pointer](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/arrays-arr/arr39-c/).

## Limitations

Cases where the pointee type has a size of `1` byte (such as, and most
importantly, `char`) are excluded.

## Options

```{option} WarnOnSizeOfConstant
When `true`, the check will warn on an expression like
`sizeof(CONSTANT)`. Default is `true`.
```

```{option} WarnOnSizeOfIntegerExpression
When `true`, the check will warn on an expression like `sizeof(expr)`
where the expression results in an integer. Default is `false`.
```

```{option} WarnOnSizeOfThis
When `true`, the check will warn on an expression like `sizeof(this)`.
Default is `true`.
```

```{option} WarnOnSizeOfCompareToConstant
When `true`, the check will warn on an expression like
`sizeof(expr) <= k` for a suspicious constant `k` while `k` is `0` or
greater than `0x8000`. Default is `true`.
```

```{option} WarnOnSizeOfPointerToAggregate
When `true`, the check will warn when the argument of `sizeof` is either a
pointer-to-aggregate type, an expression returning a pointer-to-aggregate
value or an expression that returns a pointer from an array-to-pointer
conversion (that may be implicit or explicit, for example `array + 2` or
`(int *)array`). Default is `true`.
```

```{option} WarnOnSizeOfPointer
When `true`, the check will report all expressions where the argument of
`sizeof` is an expression that produces a pointer (except for a few
idiomatic expressions that are probably intentional and correct).
This detects occurrences of CWE 467. Default is `false`.
```

```{option} WarnOnOffsetDividedBySizeOf
When `true`, the check will warn on pointer arithmetic where the
element count is obtained from a division with `sizeof(...)`,
e.g., `Ptr + Bytes / sizeof(*T)`. Default is `true`.
```

```{option} WarnOnSizeOfInLoopTermination
When `true`, the check will warn about incorrect use of sizeof expression
in loop termination condition. The warning triggers if the `sizeof`
expression appears to be incorrectly used to determine the number of
array/buffer elements.
e.g, `long arr[10]; for(int i = 0; i < sizeof(arr); i++) { ... }`. Default
is `true`.
```

```{title} clang-tidy - bugprone-spuriously-wake-up-functions
```

# bugprone-spuriously-wake-up-functions

Finds `cnd_wait`, `cnd_timedwait`, `wait`, `wait_for`, or
`wait_until` function calls when the function is not invoked from a loop
that checks whether a condition predicate holds or the function has a
condition parameter.

```cpp
if (condition_predicate) {
 condition.wait(lk);
}
```

```c
if (condition_predicate) {
 if (thrd_success != cnd_wait(&condition, &lock)) {
 }
}
```

This check corresponds to the CERT C++ Coding Standard rule
[CON54-CPP. Wrap functions that can spuriously wake up in a loop](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/concurrency-con/con54-cpp/).
and CERT C Coding Standard rule
[CON36-C. Wrap functions that can spuriously wake up in a loop](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/concurrency-con/con36-c/).

```{title} clang-tidy - bugprone-standalone-empty
```

# bugprone-standalone-empty

Warns when `empty()` is used on a range and the result is ignored. Suggests
`clear()` if it is an existing member function.

The `empty()` method on several common ranges returns a Boolean indicating
whether or not the range is empty, but is often mistakenly interpreted as
a way to clear the contents of a range. Some ranges offer a `clear()`
method for this purpose. This check warns when a call to empty returns a
result that is ignored, and suggests replacing it with a call to `clear()`
if it is available as a member function of the range.

For example, the following code could be used to indicate whether a range
is empty or not, but the result is ignored:

```c++
std::vector<int> v;
...
v.empty();
```

A call to `clear()` would appropriately clear the contents of the range:

```c++
std::vector<int> v;
...
v.clear();
```

## Limitations

- Doesn't warn if `empty()` is defined and used with the ignore result in the
 class template definition (for example in the library implementation). These
 error cases can be caught with `[[nodiscard]]` attribute.

```{title} clang-tidy - bugprone-std-exception-baseclass
```

# bugprone-std-exception-baseclass

Ensure that every value that in a `throw` expression is an instance of
`std::exception`.

Deriving all exceptions from `std::exception` allows callers to catch
all exceptions with a single catch block and provides access to the
`what()` method for diagnostics. Throwing arbitrary types creates
hidden contracts, reduces interoperability with the standard library,
and may result in program termination.

```c++
class custom_exception {};

void throwing() noexcept(false) {
 // Problematic throw expressions.
 throw int(42);
 throw custom_exception();
}

class mathematical_error : public std::exception {};

void throwing2() noexcept(false) {
 // These kind of throws are ok.
 throw mathematical_error();
 throw std::runtime_error();
 throw std::exception();
}
```

```{title} clang-tidy - bugprone-std-namespace-modification
```

# bugprone-std-namespace-modification

Warns on modifications of the `std` or `posix` namespaces which can
result in undefined behavior.

The `std` (or `posix`) namespace is allowed to be extended with (class or
function) template specializations that depend on an user-defined type (a type
that is not defined in the standard system headers).

The check detects the following (user provided) declarations in namespace
`std` or `posix`:

- Anything that is not a template specialization.
- Explicit specializations of any standard library function template or class
 template, if it does not have any user-defined type as template argument.
- Explicit specializations of any member function of a standard library class
 template.
- Explicit specializations of any member function template of a standard
 library class or class template.
- Explicit or partial specialization of any member class template of a standard
 library class or class template.

Examples:

```c++
namespace std {
 int x; // warning: modification of 'std' namespace can result in undefined behavior [bugprone-dont-modify-std-namespace]
}

namespace posix::a { // warning: modification of 'posix' namespace can result in undefined behavior
}

template <>
struct ::std::hash<long> { // warning: modification of 'std' namespace can result in undefined behavior
 unsigned long operator()(const long &K) const {
 return K;
 }
};

struct MyData { long data; };

template <>
struct ::std::hash<MyData> { // no warning: specialization with user-defined type
 unsigned long operator()(const MyData &K) const {
 return K.data;
 }
};

namespace std {
 template <>
 void swap<bool>(bool &a, bool &b); // warning: modification of 'std' namespace can result in undefined behavior

 template <>
 bool less<void>::operator()<MyData &&, MyData &&>(MyData &&, MyData &&) const { // warning: modification of 'std' namespace can result in undefined behavior
 return true;
 }
}
```

## References

This check corresponds to the CERT C++ Coding Standard rule
[DCL58-CPP. Do not modify the standard namespaces](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/declarations-and-initialization-dcl/dcl58-cpp/).

```{title} clang-tidy - bugprone-string-constructor
```

# bugprone-string-constructor

Finds string constructors that are suspicious and probably errors.

A common mistake is to swap parameters to the 'fill' string-constructor.

Examples:

```c++
std::string str('x', 50); // should be str(50, 'x')
```

Calling the string-literal constructor with a length bigger than the literal is
suspicious and adds extra random characters to the string.

Examples:

```c++
std::string("test", 200); // Will include random characters after "test".
std::string("test", 2, 5); // Will include random characters after "st".
std::string_view("test", 200);
```

Creating an empty string from constructors with parameters is considered
suspicious. The programmer should use the empty constructor instead.

Examples:

```c++
std::string("test", 0); // Creation of an empty string.
std::string("test", 1, 0);
std::string_view("test", 0);
```

Passing an invalid first character position parameter to constructor will
cause `std::out_of_range` exception at runtime.

Examples:

```c++
std::string("test", -1, 10); // Negative first character position.
std::string("test", 10, 10); // First character position is bigger than string literal character range".
```

## Options

```{option} WarnOnLargeLength
When `true`, the check will warn on a string with a length greater than
{option}`LargeLengthThreshold`. Default is `true`.
```

```{option} LargeLengthThreshold
An integer specifying the large length threshold. Default is `0x800000`.
```

```{option} StringNames
Semicolon-delimited list of class names to apply this check to.
By default `::std::basic_string` applies to `std::string` and
`std::wstring`. Set to e.g. `::std::basic_string;llvm::StringRef;QString`
to perform this check on custom classes.
Default is `::std::basic_string;::std::basic_string_view`.
```

```{title} clang-tidy - bugprone-string-integer-assignment
```

# bugprone-string-integer-assignment

The check finds assignments of an integer to `std::basic_string<CharT>`
(`std::string`, `std::wstring`, etc.). The source of the problem is the
following assignment operator of `std::basic_string<CharT>`:

```c++
basic_string& operator=( CharT ch );
```

Numeric types can be implicitly casted to character types.

```c++
std::string s;
int x = 5965;
s = 6;
s = x;
```

Use the appropriate conversion functions or character literals.

```c++
std::string s;
int x = 5965;
s = '6';
s = std::to_string(x);
```

In order to suppress false positives, use an explicit cast.

```c++
std::string s;
s = static_cast<char>(6);
```

```{title} clang-tidy - bugprone-string-literal-with-embedded-nul
```

# bugprone-string-literal-with-embedded-nul

Finds occurrences of string literal with embedded NUL character and validates
their usage.

## Invalid escaping

Special characters can be escaped within a string literal by using their
hexadecimal encoding like `\x42`. A common mistake is to escape them
like this `\0x42` where the `\0` stands for the NUL character.

```c++
const char* Example[] = "Invalid character: \0x12 should be \x12";
const char* Bytes[] = "\x03\0x02\0x01\0x00\0xFF\0xFF\0xFF";
```

## Truncated literal

String-like classes can manipulate strings with embedded NUL as they are
keeping track of the bytes and the length. This is not the case for a
`char*` (NUL-terminated) string.

A common mistake is to pass a string-literal with embedded NUL to a string
constructor expecting a NUL-terminated string. The bytes after the first NUL
character are truncated.

```c++
std::string str("abc\0def"); // "def" is truncated
str += "\0"; // This statement is doing nothing
if (str == "\0abc") return; // This expression is always true
```

```{title} clang-tidy - bugprone-stringview-nullptr
```

# bugprone-stringview-nullptr

Checks for various ways that the `const CharT*` constructor of
`std::basic_string_view` can be passed a null argument and replaces them
with the default constructor in most cases. For the comparison operators,
braced initializer list does not compile so instead a call to `.empty()`
or the empty string literal are used, where appropriate.

This prevents code from invoking behavior which is unconditionally undefined.
The single-argument `const CharT*` constructor does not check for the null
case before dereferencing its input. The standard is slated to add an
explicitly-deleted overload to catch some of these cases: wg21.link/p2166

To catch the additional cases of `NULL` (which expands to `__null`) and
`0`, first run the `modernize-use-nullptr` check to convert the callers to
`nullptr`.

```c++
std::string_view sv = nullptr;

sv = nullptr;

bool is_empty = sv == nullptr;
bool isnt_empty = sv != nullptr;

accepts_sv(nullptr);

accepts_sv({{}}); // A

accepts_sv({nullptr, 0}); // B
```

is translated into...

```c++
std::string_view sv = {};

sv = {};

bool is_empty = sv.empty();
bool isnt_empty = !sv.empty();

accepts_sv("");

accepts_sv(""); // A

accepts_sv({nullptr, 0}); // B
```

```{note}
The source pattern with trailing comment "A" selects the `(const CharT*)`
constructor overload and then value-initializes the pointer, causing a null
dereference. It happens to not include the `nullptr` literal, but it is
still within the scope of this check.
```

```{note}
The source pattern with trailing comment "B" selects the
`(const CharT*, size_type)` constructor which is perfectly valid, since the
length argument is `0`. It is not changed by this check.
```

```{title} clang-tidy - bugprone-suspicious-enum-usage
```

# bugprone-suspicious-enum-usage

The check detects various cases when an enum is probably misused
(as a bitmask).

1. When "ADD" or "bitwise OR" is used between two enum which come
 from different types and these types value ranges are not disjoint.

The following cases will be investigated only using {option}`StrictMode`. We
regard the enum as a (suspicious)
bitmask if the three conditions below are true at the same time:

- at most half of the elements of the enum are non pow-of-2 numbers (because of
 short enumerations)
- there is another non pow-of-2 number than the enum constant representing all
 choices (the result "bitwise OR" operation of all enum elements)
- enum type variable/enumconstant is used as an argument of a
 `+` or "bitwise
 OR" operator

So whenever the non pow-of-2 element is used as a bitmask element we diagnose a
misuse and give a warning.

2. Investigating the right hand side of `+=` and `|=` operator.
3. Check only the enum value side of a `|` and `+` operator if one of
 them is not enum val.
4. Check both side of `|` or `+` operator where the enum values are from
 the same enum type.

Examples:

```c++
enum { A, B, C };
enum { D, E, F = 5 };
enum { G = 10, H = 11, I = 12 };

unsigned flag;
flag =
 A |
 H; // OK, disjoint value intervals in the enum types ->probably good use.
flag = B | F; // Warning, have common values so they are probably misused.

// Case 2:
enum Bitmask {
 A = 0,
 B = 1,
 C = 2,
 D = 4,
 E = 8,
 F = 16,
 G = 31 // OK, real bitmask.
};

enum Almostbitmask {
 AA = 0,
 BB = 1,
 CC = 2,
 DD = 4,
 EE = 8,
 FF = 16,
 GG // Problem, forgot to initialize.
};

unsigned flag = 0;
flag |= E; // OK.
flag |=
 EE; // Warning at the decl, and note that it was used here as a bitmask.
```

## Options

```{option} StrictMode
When non-null the suspicious bitmask usage will be investigated additionally
to the different enum usage check.
Default is `0`.
```

```{title} clang-tidy - bugprone-suspicious-include
```

# bugprone-suspicious-include

The check detects various cases when an include refers to what appears to be an
implementation file, which often leads to hard-to-track-down ODR violations.

Examples:

```cpp
#include "Dinosaur.hpp" // OK, .hpp files tend not to have definitions.
#include "Pterodactyl.h" // OK, .h files tend not to have definitions.
#include "Velociraptor.cpp" // Warning, filename is suspicious.
#include_next <stdio.c> // Warning, filename is suspicious.
```

## Options

```{option} IgnoredRegex
A regular expression for the file name to be ignored by the check. Default
is empty string.
```

```{title} clang-tidy - bugprone-suspicious-memory-comparison
```

# bugprone-suspicious-memory-comparison

Finds potentially incorrect calls to `memcmp()` based on properties of the
arguments. The following cases are covered:

**Case 1: Non-standard-layout type**

Comparing the object representations of non-standard-layout objects may not
properly compare the value representations.

**Case 2: Types with no unique object representation**

Objects with the same value may not have the same object representation.
This may be caused by padding or floating-point types.

See also:
[EXP42-C. Do not compare padding data](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/expressions-exp/exp42-c/)
and
[FLP37-C. Do not use object representations to compare floating-point values](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/floating-point-flp/flp37-c/)

This check is also related to and partially overlaps the CERT C++ Coding Standard rules
[OOP57-CPP. Prefer special member functions and overloaded operators to
C Standard Library functions](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/object-oriented-programming-oop/oop57-cpp/)
and
[EXP62-CPP. Do not access the bits of an object representation that are not
part of the object's value representation](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/expressions-exp/exp62-cpp/)

`cert-exp42-c` redirects here as an alias of this check.

```{title} clang-tidy - bugprone-suspicious-memset-usage
```

# bugprone-suspicious-memset-usage

This check finds `memset()` calls with potential mistakes in their arguments.
Considering the function as `void* memset(void* destination, int fill_value,
size_t byte_count)`, the following cases are covered:

**Case 1: Fill value is a character `'0'`**

Filling up a memory area with ASCII code 48 characters is not customary,
possibly integer zeroes were intended instead.
The check offers a replacement of `'0'` with `0`. Memsetting character
pointers with `'0'` is allowed.

**Case 2: Fill value is truncated**

Memset converts `fill_value` to `unsigned char` before using it. If
`fill_value` is out of unsigned character range, it gets truncated
and memory will not contain the desired pattern.

**Case 3: Byte count is zero**

Calling memset with a literal zero in its `byte_count` argument is likely
to be unintended and swapped with `fill_value`. The check offers to swap
these two arguments.

Corresponding cpplint.py check name: `runtime/memset`.

Examples:

```c++
void foo() {
 int i[5] = {1, 2, 3, 4, 5};
 int *ip = i;
 char c = '1';
 char *cp = &c;
 int v = 0;

 // Case 1
 memset(ip, '0', 1); // suspicious
 memset(cp, '0', 1); // OK

 // Case 2
 memset(ip, 0xabcd, 1); // fill value gets truncated
 memset(ip, 0x00, 1); // OK

 // Case 3
 memset(ip, sizeof(int), v); // zero length, potentially swapped
 memset(ip, 0, 1); // OK
}
```

```{title} clang-tidy - bugprone-suspicious-missing-comma
```

# bugprone-suspicious-missing-comma

String literals placed side-by-side are concatenated at translation phase 6
(after the preprocessor). This feature is used to represent long string
literal on multiple lines.

For instance, the following declarations are equivalent:

```c++
const char* A[] = "This is a test";
const char* B[] = "This" " is a " "test";
```

A common mistake done by programmers is to forget a comma between two string
literals in an array initializer list.

```c++
const char* Test[] = {
 "line 1",
 "line 2" // Missing comma!
 "line 3",
 "line 4",
 "line 5"
};
```

The array contains the string "line 2line3" at offset 1 (i.e. Test[1]). Clang
won't generate warnings at compile time.

This check may warn incorrectly on cases like:

```c++
const char* SupportedFormat[] = {
 "Error %s",
 "Code " PRIu64, // May warn here.
 "Warning %s",
};
```

## Options

```{option} SizeThreshold
An unsigned integer specifying the minimum size of a string literal to be
considered by the check. Default is `5U`.
```

```{option} RatioThreshold
A string specifying the maximum threshold ratio [0, 1.0] of suspicious string
literals to be considered. Default is `".2"`.
```

```{option} MaxConcatenatedTokens
An unsigned integer specifying the maximum number of concatenated tokens.
Default is `5U`.
```

```{title} clang-tidy - bugprone-suspicious-realloc-usage
```

# bugprone-suspicious-realloc-usage

This check finds usages of `realloc` where the return value is assigned to
the same expression as passed to the first argument:
`p = realloc(p, size);`
The problem with this construct is that if `realloc` fails it returns a
null pointer but does not deallocate the original memory. If no other variable
is pointing to it, the original memory block is not available any more for the
program to use or free. In either case `p = realloc(p, size);` indicates bad
coding style and can be replaced by `q = realloc(p, size);`.

The pointer expression (used at `realloc`) can be a variable or a field
member of a data structure, but can not contain function calls or unresolved
types.

In obvious cases when the pointer used at realloc is assigned to another
variable before the `realloc` call, no warning is emitted. This happens only
if a simple expression in form of `q = p` or `void *q = p` is found in the
same function where `p = realloc(p, ...)` is found. The assignment has to be
before the call to realloc (but otherwise at any place) in the same function.
This suppression works only if `p` is a single variable.

Examples:

```c++
struct A {
 void *p;
};

A &getA();

void foo(void *p, A *a, int new_size) {
 p = realloc(p, new_size); // warning: 'p' may be set to null if 'realloc' fails, which may result in a leak of the original buffer
 a->p = realloc(a->p, new_size); // warning: 'a->p' may be set to null if 'realloc' fails, which may result in a leak of the original buffer
 getA().p = realloc(getA().p, new_size); // no warning
}

void foo1(void *p, int new_size) {
 void *p1 = p;
 p = realloc(p, new_size); // no warning
}
```

```{title} clang-tidy - bugprone-suspicious-semicolon
```

# bugprone-suspicious-semicolon

Finds most instances of stray semicolons that unexpectedly alter the meaning of
the code. More specifically, it looks for `if`, `while`, `for` and
`for-range` statements whose body is a single semicolon, and then analyzes
the context of the code (e.g. indentation) in an attempt to determine whether
that is intentional.

```c++
if (x < y);
{
 x++;
}
```

Here the body of the `if` statement consists of only the semicolon at the end
of the first line, and `x` will be incremented regardless of the condition.

```c++
while ((line = readLine(file)) != NULL);
 processLine(line);
```

As a result of this code, `processLine()` will only be called once, when the
`while` loop with the empty body exits with `line == NULL`. The indentation
of the code indicates the intention of the programmer.

```c++
if (x >= y);
x -= y;
```

While the indentation does not imply any nesting, there is simply no valid
reason to have an `if` statement with an empty body (but it can make sense for
a loop). So this check issues a warning for the code above.

To solve the issue remove the stray semicolon or in case the empty body is
intentional, reflect this using code indentation or put the semicolon in a new
line. For example:

```c++
while (readWhitespace());
 Token t = readNextToken();
```

Here the second line is indented in a way that suggests that it is meant to be
the body of the `while` loop - whose body is in fact empty, because of the
semicolon at the end of the first line.

Either remove the indentation from the second line:

```c++
while (readWhitespace());
Token t = readNextToken();
```

... or move the semicolon from the end of the first line to a new line:

```c++
while (readWhitespace())
 ;

 Token t = readNextToken();
```

In this case the check will assume that you know what you are doing, and will
not raise a warning.

```{title} clang-tidy - bugprone-suspicious-string-compare
```

# bugprone-suspicious-string-compare

Find suspicious usage of runtime string comparison functions.
This check is valid in C and C++.

Checks for calls with implicit comparator and proposed to explicitly add it.

```c++
if (strcmp(...)) // Implicitly compare to zero
if (!strcmp(...)) // Won't warn
if (strcmp(...) != 0) // Won't warn
```

Checks that compare function results (i.e., `strcmp`) are compared to valid
constant. The resulting value is

```
< 0 when lower than,
> 0 when greater than,
== 0 when equals.
```

A common mistake is to compare the result to `1` or `-1`.

```c++
if (strcmp(...) == -1) // Incorrect usage of the returned value.
```

Additionally, the check warns if the results value is implicitly cast to a
*suspicious* non-integer type. It's happening when the returned value is
used in a wrong context.

```c++
if (strcmp(...) < 0.) // Incorrect usage of the returned value.
```

The check will detect the following string comparison functions:
`__builtin_memcmp`, `__builtin_strcasecmp`, `__builtin_strcmp`,
`__builtin_strncasecmp`, `__builtin_strncmp`, `_mbscmp`, `_mbscmp_l`,
`_mbsicmp`, `_mbsicmp_l`, `_mbsnbcmp`, `_mbsnbcmp_l`, `_mbsnbicmp`,
`_mbsnbicmp_l`, `_mbsncmp`, `_mbsncmp_l`, `_mbsnicmp`, `_mbsnicmp_l`,
`_memicmp`, `_memicmp_l`, `_stricmp`, `_stricmp_l`, `_strnicmp`,
`_strnicmp_l`, `_wcsicmp`, `_wcsicmp_l`, `_wcsnicmp`, `_wcsnicmp_l`,
`lstrcmp`, `lstrcmpi`, `memcmp`, `memicmp`, `strcasecmp`, `strcmp`,
`strcmpi`, `stricmp`, `strncasecmp`, `strncmp`, `strnicmp`, `wcscasecmp`,
`wcscmp`, `wcsicmp`, `wcsncmp`, `wcsnicmp`, `wmemcmp`.

## Options

```{option} WarnOnImplicitComparison
When `true`, the check will warn on implicit comparison. Default is `true`.
```

```{option} WarnOnLogicalNotComparison
When `true`, the check will warn on logical not comparison. Default is `false`.
```

```{option} StringCompareLikeFunctions
A string specifying the comma-separated names of the extra string comparison
functions. Default is an empty string.
```

```{title} clang-tidy - bugprone-suspicious-stringview-data-usage
```

# bugprone-suspicious-stringview-data-usage

Identifies suspicious usages of `std::string_view::data()` that could lead to
reading out-of-bounds data due to inadequate or incorrect string null
termination.

It warns when the result of `data()` is passed to a constructor or function
without also passing the corresponding result of `size()` or `length()`
member function. Such usage can lead to unintended behavior, particularly when
assuming the data pointed to by `data()` is null-terminated.

The absence of a `c_str()` method in `std::string_view` often leads
developers to use `data()` as a substitute, especially when interfacing with
C APIs that expect null-terminated strings. However, since `data()` does not
guarantee null termination, this can result in unintended behavior if the API
relies on proper null termination for correct string interpretation.

In today's programming landscape, this scenario can occur when implicitly
converting an `std::string_view` to an `std::string`. Since the constructor
in `std::string` designed for string-view-like objects is `explicit`,
attempting to pass an `std::string_view` to a function expecting an
`std::string` will result in a compilation error. As a workaround, developers
may be tempted to utilize the `.data()` method to achieve compilation,
introducing potential risks.

For instance:

```c++
void printString(const std::string& str) {
 std::cout << "String: " << str << std::endl;
}

void something(std::string_view sv) {
 printString(sv.data());
}
```

In this example, directly passing `sv` to the `printString` function would
lead to a compilation error due to the explicit nature of the `std::string`
constructor. Consequently, developers might opt for `sv.data()` to resolve the
compilation error, albeit introducing potential hazards as discussed.

## Options

```{option} StringViewTypes
Option allows users to specify custom string view-like types for analysis. It
accepts a semicolon-separated list of type names or regular expressions
matching these types. Default value is `::std::basic_string_view;::llvm::StringRef`.
```

```{option} AllowedCallees
Specifies methods, functions, or classes where the result of `.data()` is
passed to. Allows to exclude such calls from the analysis. Accepts a
semicolon-separated list of names or regular expressions matching these
entities. Default value is empty string.
```

```{title} clang-tidy - bugprone-swapped-arguments
```

# bugprone-swapped-arguments

Finds potentially swapped arguments by examining implicit conversions.
It analyzes the types of the arguments being passed to a function and compares
them to the expected types of the corresponding parameters. If there is a
mismatch or an implicit conversion that indicates a potential swap, a warning
is raised.

```c++
void printNumbers(int a, float b);

int main() {
 // Swapped arguments: float passed as int, int as float)
 printNumbers(10.0f, 5);
 return 0;
}
```

Covers a wide range of implicit conversions, including:
- User-defined conversions
- Conversions from floating-point types to boolean or integral types
- Conversions from integral types to boolean or floating-point types
- Conversions from boolean to integer types or floating-point types
- Conversions from (member) pointers to boolean

It is important to note that for most argument swaps, the types need to match
exactly. However, there are exceptions to this rule. Specifically, when the
swapped argument is of integral type, an exact match is not always necessary.
Implicit casts from other integral types are also accepted. Similarly, when
dealing with floating-point arguments, implicit casts between different
floating-point types are considered acceptable.

To avoid confusion, swaps where both swapped arguments are of integral types or
both are of floating-point types do not trigger the warning. In such cases,
it's assumed that the developer intentionally used different integral or
floating-point types and does not raise a warning. This approach prevents false
positives and provides flexibility in handling situations where varying
integral or floating-point types are intentionally utilized.

```{title} clang-tidy - bugprone-switch-missing-default-case
```

# bugprone-switch-missing-default-case

Ensures that switch statements without default cases are flagged, focuses only
on covering cases with non-enums where the compiler may not issue warnings.

Switch statements without a default case can lead to unexpected
behavior and incomplete handling of all possible cases. When a switch statement
lacks a default case, if a value is encountered that does not match any of the
specified cases, the switch statement will do nothing and the program will
continue execution without handling the value.

This check helps identify switch statements that are missing a default case,
allowing developers to ensure that all possible cases are handled properly.
Adding a default case allows for graceful handling of unexpected or unmatched
values, reducing the risk of program errors and unexpected behavior.

Example:

```c++
// Example 1:
// warning: switching on non-enum value without default case may not cover all cases
switch (i) {
case 0:
 break;
}

// Example 2:
enum E { eE1 };
E e = eE1;
switch (e) { // no-warning
case eE1:
 break;
}

// Example 3:
int i = 0;
switch (i) { // no-warning
case 0:
 break;
default:
 break;
}
```

```{note}
Enum types are already covered by compiler warnings (comes under -Wswitch)
when a switch statement does not handle all enum values. This check focuses
on non-enum types where the compiler warnings may not be present.
```

```{seealso}
The [CppCoreGuideline ES.79](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#res-default)
provide guidelines on switch statements, including the recommendation to
always provide a default case.
```

```{title} clang-tidy - bugprone-tagged-union-member-count
```

# bugprone-tagged-union-member-count

Gives warnings for tagged unions, where the number of tags is
different from the number of data members inside the union.

A struct or a class is considered to be a tagged union if it has
exactly one union data member and exactly one enum data member and
any number of other data members that are neither unions or enums.
Furthermore, the types of the union and the enum members must
not come from system header files nor the `std` namespace.

Example:

```c++
enum Tags {
 Tag1,
 Tag2,
};

struct TaggedUnion { // warning: tagged union has more data members (3) than tags (2)
 enum Tags Kind;
 union {
 int I;
 float F;
 char *Str;
 } Data;
};
```

The following example illustrates the exception for unions and enums from
system header files and the `std` namespace.

```c++
#include <pthread.h>

struct NotTaggedUnion {
 enum MyEnum { MyEnumConstant1, MyEnumConstant2 } En;
 pthread_mutex_t Mutex;
};
```

The `pthread_mutex_t` type may be defined as a union behind a `typedef`,
in which case the check could mistake this type as a user-defined tagged union.
After all, it has exactly one enum data member and exactly one union data member.
To avoid false-positive cases originating from this, unions and enums from
system headers and the `std` namespace are ignored when pinpointing the
union part and the enum part of a potential user-defined tagged union.

## How enum constants are counted

The main complicating factor when counting the number of enum constants is that
some of them might be auxiliary values that purposefully don't have a
corresponding union data member and are used for something else. For example
the last enum constant sometimes explicitly "points to" the last declared valid
enum constant or tracks how many enum constants have been declared.

For an illustration:

```c++
enum TagWithLast {
 Tag1 = 0,
 Tag2 = 1,
 Tag3 = 2,
 LastTag = 2
};

enum TagWithCounter {
 Tag1, // is 0
 Tag2, // is 1
 Tag3, // is 2
 TagCount, // is 3
};
```

The check counts the number of distinct values among the enum constants and not
the enum constants themselves. This way the enum constants that are essentially
just aliases of other enum constants are not included in the final count.

Handling of counting enum constants (ones like `TagCount` in the previous
code example) is done by decreasing the number of enum values by one if the name
of the last enum constant starts with a prefix or ends with a suffix specified in
{option}`CountingEnumPrefixes`, {option}`CountingEnumSuffixes` and it's value is
one less than the total number of distinct values in the enum.

When the final count is adjusted based on this heuristic then a diagnostic note
is emitted that shows which enum constant matched the criteria.

The heuristic can be disabled entirely ({option}`EnableCountingEnumHeuristic`)
or configured to follow your naming convention ({option}`CountingEnumPrefixes`,
{option}`CountingEnumSuffixes`).
The strings specified in {option}`CountingEnumPrefixes`,
{option}`CountingEnumSuffixes` are matched case insensitively.

Example counts:

```c++
// Enum count is 3, because the value 2 is counted only once
enum TagWithLast {
 Tag1 = 0,
 Tag2 = 1,
 Tag3 = 2,
 LastTag = 2
};

// Enum count is 3, because TagCount is heuristically excluded
enum TagWithCounter {
 Tag1, // is 0
 Tag2, // is 1
 Tag3, // is 2
 TagCount, // is 3
};
```

## Options

```{option} EnableCountingEnumHeuristic
```

This option enables or disables the counting enum heuristic.
It uses the prefixes and suffixes specified in the options
{option}`CountingEnumPrefixes`, {option}`CountingEnumSuffixes` to find counting enum constants by
using them for prefix and suffix matching.

Default is `true`.

When {option}`EnableCountingEnumHeuristic` is `false`:

```c++
enum TagWithCounter {
 Tag1,
 Tag2,
 Tag3,
 TagCount,
};

struct TaggedUnion {
 TagWithCounter Kind;
 union {
 int A;
 long B;
 char *Str;
 float F;
 } Data;
};
```

When {option}`EnableCountingEnumHeuristic` is `true`:

```c++
enum TagWithCounter {
 Tag1,
 Tag2,
 Tag3,
 TagCount,
};

struct TaggedUnion { // warning: tagged union has more data members (4) than tags (3)
 TagWithCounter Kind;
 union {
 int A;
 long B;
 char *Str;
 float F;
 } Data;
};
```

```{option} CountingEnumPrefixes
```

See {option}`CountingEnumSuffixes` below.

```{option} CountingEnumSuffixes
```

{option}`CountingEnumPrefixes` and {option}`CountingEnumSuffixes` are lists of semicolon
separated strings that are used to search for possible counting enum constants.
These strings are matched case insensitively as prefixes and suffixes
respectively on the names of the enum constants.
If {option}`EnableCountingEnumHeuristic` is `false` then these options do nothing.

The default value of {option}`CountingEnumSuffixes` is
`count` and of
{option}`CountingEnumPrefixes` is the empty string.

When {option}`EnableCountingEnumHeuristic` is `true` and
{option}`CountingEnumSuffixes` is `count;size`:

```c++
enum TagWithCounterCount {
 Tag1,
 Tag2,
 Tag3,
 TagCount,
};

struct TaggedUnionCount { // warning: tagged union has more data members (4) than tags (3)
 TagWithCounterCount Kind;
 union {
 int A;
 long B;
 char *Str;
 float F;
 } Data;
};

enum TagWithCounterSize {
 Tag11,
 Tag22,
 Tag33,
 TagSize,
};

struct TaggedUnionSize { // warning: tagged union has more data members (4) than tags (3)
 TagWithCounterSize Kind;
 union {
 int A;
 long B;
 char *Str;
 float F;
 } Data;
};
```

When {option}`EnableCountingEnumHeuristic` is `true` and
{option}`CountingEnumPrefixes` is `maxsize;last_`

```c++
enum TagWithCounterLast {
 Tag1,
 Tag2,
 Tag3,
 last_tag,
};

struct TaggedUnionLast { // warning: tagged union has more data members (4) than tags (3)
 TagWithCounterLast tag;
 union {
 int I;
 short S;
 char *C;
 float F;
 } Data;
};

enum TagWithCounterMaxSize {
 Tag1,
 Tag2,
 Tag3,
 MaxSizeTag,
};

struct TaggedUnionMaxSize { // warning: tagged union has more data members (4) than tags (3)
 TagWithCounterMaxSize tag;
 union {
 int I;
 short S;
 char *C;
 float F;
 } Data;
};
```

```{option} StrictMode
```

When enabled, the check will also give a warning, when the number of tags
is greater than the number of union data members.

Default is `false`.

When {option}`StrictMode` is `false`:

```c++
struct TaggedUnion {
 enum {
 Tag1,
 Tag2,
 Tag3,
 } Tags;
 union {
 int I;
 float F;
 } Data;
};
```

When {option}`StrictMode` is `true`:

```c++
struct TaggedUnion { // warning: tagged union has fewer data members (2) than tags (3)
 enum {
 Tag1,
 Tag2,
 Tag3,
 } Tags;
 union {
 int I;
 float F;
 } Data;
};
```

```{title} clang-tidy - bugprone-terminating-continue
```

# bugprone-terminating-continue

Detects `do while` loops with a condition always evaluating to false that
have a `continue` statement, as this `continue` terminates the loop
effectively.

```cpp
void f() {
do {
 // some code
 continue; // terminating continue
 // some other code
} while(false);
```

```{title} clang-tidy - bugprone-throw-keyword-missing
```

# bugprone-throw-keyword-missing

Warns about a potentially missing `throw` keyword. If a temporary object
is created, but the object's type derives from (or is the same as) a class
that has 'EXCEPTION', 'Exception' or 'exception' in its name, we can assume
that the programmer's intention was to throw that object.

Example:

```cpp
void f(int i) {
 if (i < 0) {
 // Exception is created but is not thrown.
 std::runtime_error("Unexpected argument");
 }
}
```

```{title} clang-tidy - bugprone-throwing-static-initialization
```

# bugprone-throwing-static-initialization

Finds all `static` or `thread_local` variable declarations where the
initializer for the object may throw an exception.

## Options

```{option} AllowedTypes

 A semicolon-separated list of names of types that will be excluded from
 this check (declarations with matching type will be excluded). Regular
 expressions are accepted, e.g. `[Rr]ef(erence)?$` matches every type with
 suffix `Ref`, `ref`, `Reference` and `reference`. If a name in the
 list contains the sequence `::`, it is matched against the qualified type
 name (i.e. `namespace::Type`), otherwise it is matched against only the
 type name (i.e. `Type`). Default is an empty string.
```

## References

This check corresponds to the CERT C++ Coding Standard rule
[ERR58-CPP. Handle all exceptions thrown before main() begins executing](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/exceptions-and-error-handling-err/err58-cpp/).

```{title} clang-tidy - bugprone-too-small-loop-variable
```

# bugprone-too-small-loop-variable

Detects those `for` loops that have a loop variable with a "too small" type
which means this type can't represent all values which are part of the
iteration range.

```c++
int main() {
 long size = 294967296l;
 for (short i = 0; i < size; ++i) {}
}
```

This `for` loop is an infinite loop because the `short` type can't
represent all values in the `[0..size]` interval.

In a real use case size means a container's size which depends on the
user input.

```c++
int doSomething(const std::vector& items) {
 for (short i = 0; i < items.size(); ++i) {}
}
```

This algorithm works for a small amount of objects, but will lead to freeze for
a larger user input.

It's recommended to enable the compiler warning
`-Wtautological-constant-out-of-range-compare` as well, since check does
not inspect compile-time constant loop boundaries to avoid overlaps with
the warning.

## Options

```{option} MagnitudeBitsUpperLimit
Upper limit for the magnitude bits of the loop variable. If it's set the check
filters out those catches in which the loop variable's type has more magnitude
bits as the specified upper limit.
For example, if the user sets this option to 31 (bits), then a 32-bit `unsigned int`
is ignored by the check, however a 32-bit `int` is not (A 32-bit `signed int`
has 31 magnitude bits).
Default value is `16`.
```

```c++
int main() {
 long size = 294967296l;
 for (unsigned i = 0; i < size; ++i) {} // no warning with MagnitudeBitsUpperLimit = 31 on a system where unsigned is 32-bit
 for (int i = 0; i < size; ++i) {} // warning with MagnitudeBitsUpperLimit = 31 on a system where int is 32-bit
}
```

```{title} clang-tidy - bugprone-unchecked-optional-access
```

# bugprone-unchecked-optional-access

*Note*: This check uses a flow-sensitive static analysis to produce its
results. Therefore, it may be more resource intensive (RAM, CPU) than the
average clang-tidy check.

This check identifies unsafe accesses to values contained in
`std::optional<T>`, `absl::optional<T>`, `base::Optional<T>`,
`folly::Optional<T>`, `bsl::optional`, or
`BloombergLP::bdlb::NullableValue` objects. Below we will refer to all these
types collectively as `optional<T>`.

An access to the value of an `optional<T>` occurs when one of its `value`,
`operator*`, or `operator->` member functions is invoked. To align with
common misconceptions, the check considers these member functions as
equivalent, even though there are subtle differences related to exceptions
versus undefined behavior. See *Additional notes*, below, for more information
on this topic.

An access to the value of an `optional<T>` is considered safe if and only if
code in the local scope (for example, a function body) ensures that the
`optional<T>` has a value in all possible execution paths that can reach the
access. That should happen either through an explicit check, using the
`optional<T>::has_value` member function, or by constructing the
`optional<T>` in a way that shows that it unambiguously holds a value (e.g
using `std::make_optional` which always returns a populated
`std::optional<T>`).

Below we list some examples, starting with unsafe optional access patterns,
followed by safe access patterns.

## Unsafe access patterns

### Access the value without checking if it exists

The check flags accesses to the value that are not locally guarded by
existence check:

```c++
void f(std::optional<int> opt) {
 use(*opt); // unsafe: it is unclear whether `opt` has a value.
}
```

### Access the value in the wrong branch

The check is aware of the state of an optional object in different
branches of the code. For example:

```c++
void f(std::optional<int> opt) {
 if (opt.has_value()) {
 } else {
 use(opt.value()); // unsafe: it is clear that `opt` does *not* have a value.
 }
}
```

### Assume a function result to be stable

The check is aware that function results might not be stable. That is,
consecutive calls to the same function might return different values.
For example:

```c++
void f(Foo foo) {
 if (foo.take().has_value()) {
 use(*foo.take()); // unsafe: it is unclear whether `foo.take()` has a value.
 }
}
```

#### Exception: accessor methods

The check assumes *accessor* methods of a class are stable, with a heuristic to
determine which methods are accessors. Specifically, parameter-free `const`
methods and smart pointer-like APIs (non `const` overloads of `*` when
there is a parallel `const` overload) are treated as accessors. Note that
this is not guaranteed to be safe -- but, it is widely used (safely) in
practice. Calls to non `const` methods are assumed to modify the state of
the object and affect the stability of earlier accessor calls.

### Rely on invariants of uncommon APIs

The check is unaware of invariants of uncommon APIs. For example:

```c++
void f(Foo foo) {
 if (foo.HasProperty("bar")) {
 use(*foo.GetProperty("bar")); // unsafe: it is unclear whether `foo.GetProperty("bar")` has a value.
 }
}
```

### Check if a value exists, then pass the optional to another function

The check relies on local reasoning. The check and value access must
both happen in the same function. An access is considered unsafe even if
the caller of the function performing the access ensures that the
optional has a value. For example:

```c++
void g(std::optional<int> opt) {
 use(*opt); // unsafe: it is unclear whether `opt` has a value.
}

void f(std::optional<int> opt) {
 if (opt.has_value()) {
 g(opt);
 }
}
```

## Safe access patterns

### Check if a value exists, then access the value

The check recognizes all straightforward ways for checking if a value
exists and accessing the value contained in an optional object. For
example:

```c++
void f(std::optional<int> opt) {
 if (opt.has_value()) {
 use(*opt);
 }
}
```

### Check if a value exists, then access the value from a copy

The criteria that the check uses is semantic, not syntactic. It
recognizes when a copy of the optional object being accessed is known to
have a value. For example:

```c++
void f(std::optional<int> opt1) {
 if (opt1.has_value()) {
 std::optional<int> opt2 = opt1;
 use(*opt2);
 }
}
```

### Ensure that a value exists using common macros

The check is aware of common macros like `CHECK` and `DCHECK`. Those can be
used to ensure that an optional object has a value. For example:

```c++
void f(std::optional<int> opt) {
 DCHECK(opt.has_value());
 use(*opt);
}
```

### Ensure that a value exists, then access the value in a correlated branch

The check is aware of correlated branches in the code and can figure out
when an optional object is ensured to have a value on all execution
paths that lead to an access. For example:

```c++
void f(std::optional<int> opt) {
 bool safe = false;
 if (opt.has_value() && SomeOtherCondition()) {
 safe = true;
 }
 // ... more code...
 if (safe) {
 use(*opt);
 }
}
```

## Stabilize function results

Function results are not assumed to be stable across calls, except for
const accessor methods. For more complex accessors (non-const, or depend on
multiple params) it is best to store the result of the function call in a
local variable and use that variable to access the value. For example:

```c++
void f(Foo foo) {
 if (const auto& foo_opt = foo.take(); foo_opt.has_value()) {
 use(*foo_opt);
 }
}
```

## Do not rely on uncommon-API invariants

When uncommon APIs guarantee that an optional has contents, do not rely on it
-- instead, check explicitly that the optional object has a value. For example:

```c++
void f(Foo foo) {
 if (const auto& property = foo.GetProperty("bar")) {
 use(*property);
 }
}
```

instead of the `HasProperty`, `GetProperty`
pairing we saw above.

## Do not rely on caller-performed checks

If you know that all of a function's callers have checked that an optional
argument has a value, either change the function to take the value directly or
check the optional again in the local scope of the callee. For example:

```c++
void g(int val) {
 use(val);
}

void f(std::optional<int> opt) {
 if (opt.has_value()) {
 g(*opt);
 }
}
```

and

```c++
struct S {
 std::optional<int> opt;
 int x;
};

void g(const S &s) {
 if (s.opt.has_value() && s.x > 10) {
 use(*s.opt);
}

void f(S s) {
 if (s.opt.has_value()) {
 g(s);
 }
}
```

## Additional notes

### Aliases created via `using` declarations

The check is aware of aliases of optional types that are created via
`using` declarations. For example:

```c++
using OptionalInt = std::optional<int>;

void f(OptionalInt opt) {
 use(opt.value()); // unsafe: it is unclear whether `opt` has a value.
}
```

### Lambdas

The check does not currently report unsafe optional accesses in lambdas.
A future version will expand the scope to lambdas, following the rules
outlined above. It is best to follow the same principles when using
optionals in lambdas.

### Access with `operator*()` vs. `value()`

Given that `value()` has well-defined behavior (either throwing an exception
or terminating the program), why treat it the same as `operator*()` which
causes undefined behavior (UB)? That is, why is it considered unsafe to access
an optional with `value()`, if it's not provably populated with a value? For
that matter, why is `CHECK()` followed by `operator*()` any better than
`value()`, given that they are semantically equivalent (on configurations
that disable exceptions)?

The answer is that we assume most users do not realize the difference between
`value()` and `operator*()`. Shifting to `operator*()` and some form of
explicit value-presence check or explicit program termination has two
advantages:

- Readability. The check, and any potential side effects like program
 shutdown, are very clear in the code. Separating access from checks can
 actually make the checks more obvious.
- Performance. A single check can cover many or even all accesses within
 scope. This gives the user the best of both worlds -- the safety of a
 dynamic check, but without incurring redundant costs.

### GoogleTest awareness

The check recognizes common macros like `ASSERT_TRUE` and `ASSERT_FALSE`:

```c++
TEST(OptionalTest, CheckValue) {
 std::optional<int> opt;
 EXPECT_TRUE(opt.has_value());
 EXPECT_EQ(opt.value(), 42); // unsafe: EXPECT_TRUE doesn't terminate test.

 ASSERT_TRUE(opt.has_value());
 EXPECT_EQ(opt.value(), 42); // safe: ASSERT_TRUE terminates if no value.
}
```

Less common macros such as `ASSERT_NE(..., nullopt)` and `ASSERT_THAT` are
not currently supported and are ignored, which may result in false positives.

### Options

```{option} IgnoreSmartPointerDereference
If set to `true`, the check ignores optionals that
are reached through overloaded smart-pointer-like dereference (`operator*`,
`operator->`) on classes other than the optional type itself. This helps
avoid false positives where the analysis cannot equate results across such
calls. This does not cover access through `operator[]`. Default is `false`.
```

```{option} IgnoreValueCalls
If set to `true`, the check does not diagnose calls
to `optional::value()`. Diagnostics for `operator*()` and
`operator->()` remain enabled. This is useful for codebases that
intentionally rely on `value()` for defined, guarded access while still
flagging UB-prone operator dereferences. Default is `false`.
```

```{title} clang-tidy - bugprone-unchecked-string-to-number-conversion
```

# bugprone-unchecked-string-to-number-conversion

This check flags calls to string-to-number conversion functions that do not
verify the validity of the conversion, such as `atoi()` or `scanf()`. It
does not flag calls to `strtol()`, or other, related conversion functions
that do perform better error checking.

```c
#include <stdlib.h>

void func(const char *buff) {
 int si;

 if (buff) {
 si = atoi(buff); /* 'atoi' used to convert a string to an integer, but function will
 not report conversion errors; consider using 'strtol' instead. */
 } else {
 /* Handle error */
 }
}
```

## References

This check corresponds to the CERT C Coding Standard rule
[ERR34-C. Detect errors when converting a string to a number](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/error-handling-err/err34-c/).

```{title} clang-tidy - bugprone-undefined-memory-manipulation
```

# bugprone-undefined-memory-manipulation

Finds calls of memory manipulation functions `memset()`, `memcpy()` and
`memmove()` on non-TriviallyCopyable objects resulting in undefined behavior.

Using memory manipulation functions on non-TriviallyCopyable objects can lead
to a range of subtle and challenging issues in C++ code. The most immediate
concern is the potential for undefined behavior, where the state of the object
may become corrupted or invalid. This can manifest as crashes, data corruption,
or unexpected behavior at runtime, making it challenging to identify and
diagnose the root cause. Additionally, misuse of memory manipulation functions
can bypass essential object-specific operations, such as constructors and
destructors, leading to resource leaks or improper initialization.

For example, when using `memcpy` to copy `std::string`, pointer data is
being copied, and it can result in a double free issue.

```c++
#include <cstring>
#include <string>

int main() {
 std::string source = "Hello";
 std::string destination;

 std::memcpy(&destination, &source, sizeof(std::string));

 // Undefined behavior may occur here, during std::string destructor call.
 return 0;
}
```

```{title} clang-tidy - bugprone-undelegated-constructor
```

# bugprone-undelegated-constructor

Finds creation of temporary objects in constructors that look like a
function call to another constructor of the same class.

The user most likely meant to use a delegating constructor or base class
initializer.

```{title} clang-tidy - bugprone-unhandled-code-paths
```

# bugprone-unhandled-code-paths

This check discovers situations where code paths are not fully-covered.

`if-else if` chains that miss a final `else` branch might lead to
unexpected program execution and be the result of a logical error.
If the missing `else` branch is intended you can leave it empty with
a clarifying comment.
This warning can be noisy on some code bases, so it is disabled by default.

```c++
void f1() {
 int i = determineTheNumber();

 if(i > 0) {
 // Some Calculation
 } else if (i < 0) {
 // Precondition violated or something else.
 }
 // ...
}
```

Similar arguments hold for `switch` statements which do not cover all
possible code paths.

```c++
// The missing default branch might be a logical error. It can be kept empty
// if there is nothing to do, making it explicit.
void f2(int i) {
 switch (i) {
 case 0: // something
 break;
 case 1: // something else
 break;
 }
 // All other numbers?
}

// Violates this rule as well, but already emits a compiler warning (-Wswitch).
enum Color { Red, Green, Blue, Yellow };
void f3(enum Color c) {
 switch (c) {
 case Red: // We can't drive for now.
 break;
 case Green: // We are allowed to drive.
 break;
 }
 // Other cases missing
}
```

## Options

```{option} WarnOnMissingElse
Boolean flag that activates a warning for missing `else` branches.
Default is `false`.
```

```{title} clang-tidy - bugprone-unhandled-exception-at-new
```

# bugprone-unhandled-exception-at-new

Finds calls to `new` with missing exception handler for `std::bad_alloc`.

Calls to `new` may throw exceptions of type `std::bad_alloc` that should
be handled. Alternatively, the nonthrowing form of `new` can be
used. The check verifies that the exception is handled in the function
that calls `new`.

If a nonthrowing version is used or the exception is allowed to propagate out
of the function no warning is generated.

The exception handler is checked if it catches a `std::bad_alloc` or
`std::exception` exception type, or all exceptions (catch-all).
The check assumes that any user-defined `operator new` is either
`noexcept` or may throw an exception of type `std::bad_alloc` (or one
derived from it). Other exception class types are not taken into account.

```c++
int *f() noexcept {
 int *p = new int[1000]; // warning: missing exception handler for allocation failure at 'new'
 // ...
 return p;
}
```

```c++
int *f1() { // not 'noexcept'
 int *p = new int[1000]; // no warning: exception can be handled outside
 // of this function
 // ...
 return p;
}

int *f2() noexcept {
 try {
 int *p = new int[1000]; // no warning: exception is handled
 // ...
 return p;
 } catch (std::bad_alloc &) {
 // ...
 }
 // ...
}

int *f3() noexcept {
 int *p = new (std::nothrow) int[1000]; // no warning: "nothrow" is used
 // ...
 return p;
}
```

```{title} clang-tidy - bugprone-unhandled-self-assignment
```

# bugprone-unhandled-self-assignment

`cert-oop54-cpp` redirects here as an alias for this check. For
the CERT alias, the {option}`WarnOnlyIfThisHasSuspiciousField`
is set to `false`.

Finds user-defined copy assignment operators which do not protect the code
against self-assignment either by checking self-assignment explicitly or
using the copy-and-swap or the copy-and-move method.

By default, this check searches only those classes which have any pointer or C
array field to avoid false positives. In case of a pointer or a C array, it's
likely that self-copy assignment breaks the object if the copy assignment
operator was not written with care.

See also:
[OOP54-CPP. Gracefully handle self-copy assignment](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/object-oriented-programming-oop/oop54-cpp/)

A copy assignment operator must prevent that self-copy assignment ruins the
object state. A typical use case is when the class has a pointer field
and the copy assignment operator first releases the pointed object and
then tries to assign it:

```c++
class T {
int* p;

public:
 T(const T &rhs) : p(rhs.p ? new int(*rhs.p) : nullptr) {}
 ~T() { delete p; }

 // ...

 T& operator=(const T &rhs) {
 delete p;
 p = new int(*rhs.p);
 return *this;
 }
};
```

There are two common C++ patterns to avoid this problem. The first is
the self-assignment check:

```c++
class T {
int* p;

public:
 T(const T &rhs) : p(rhs.p ? new int(*rhs.p) : nullptr) {}
 ~T() { delete p; }

 // ...

 T& operator=(const T &rhs) {
 if(this == &rhs)
 return *this;

 delete p;
 p = new int(*rhs.p);
 return *this;
 }
};
```

The second one is the copy-and-swap method when we create a temporary copy
(using the copy constructor) and then swap this temporary object with `this`:

```c++
class T {
int* p;

public:
 T(const T &rhs) : p(rhs.p ? new int(*rhs.p) : nullptr) {}
 ~T() { delete p; }

 // ...

 void swap(T &rhs) {
 using std::swap;
 swap(p, rhs.p);
 }

 T& operator=(const T &rhs) {
 T(rhs).swap(*this);
 return *this;
 }
};
```

There is a third pattern which is less common. Let's call it the copy-and-move
method when we create a temporary copy (using the copy constructor) and then move
this temporary object into `this` (needs a move assignment operator):

```c++
class T {
int* p;

public:
 T(const T &rhs) : p(rhs.p ? new int(*rhs.p) : nullptr) {}
 ~T() { delete p; }

 // ...

 T& operator=(const T &rhs) {
 T t = rhs;
 *this = std::move(t);
 return *this;
 }

 T& operator=(T &&rhs) {
 p = rhs.p;
 rhs.p = nullptr;
 return *this;
 }
};
```

## Options

```{option} WarnOnlyIfThisHasSuspiciousField
When `true`, the check will warn only if the container class of the copy
assignment operator has any suspicious fields (pointer, C array and C++ smart
pointer). Default is `true`.
```

```{title} clang-tidy - bugprone-unintended-char-ostream-output
```

# bugprone-unintended-char-ostream-output

Finds unintended character output from `unsigned char` and `signed char` to
an `ostream`.

Normally, when `unsigned char (uint8_t)` or `signed char (int8_t)` is used,
it is more likely a number than a character. However, when it is passed
directly to `std::ostream`'s `operator<<`, the result is the character
output instead of the numeric value. This often contradicts the developer's
intent to print integer values.

```c++
uint8_t v = 65;
std::cout << v; // output 'A' instead of '65'
```

The check will suggest casting the value to an appropriate type to indicate the
intent, by default, it will cast to `unsigned int` for `unsigned char` and
`int` for `signed char`.

```c++
std::cout << static_cast<unsigned int>(v); // when v is unsigned char
std::cout << static_cast<int>(v); // when v is signed char
```

To avoid lengthy cast statements, add prefix `+` to the variable can
also suppress warnings because unary expression will promote the value
to an `int`.

```c++
std::cout << +v;
```

Or cast to char to explicitly indicate that output should be a character.

```c++
std::cout << static_cast<char>(v);
```

## Options

```{option} AllowedTypes
A semicolon-separated list of type names that will be treated like the `char`
type: the check will not report variables declared with with these types or
explicit cast expressions to these types. Note that this distinguishes type
aliases from the original type, so specifying e.g. `unsigned char` here
will not suppress reports about `uint8_t` even if it is defined as a
`typedef` alias for `unsigned char`.
Default is `unsigned char;signed char`.
```

```{option} CastTypeName
When {option}`CastTypeName` is specified, the fix-it will use
{option}`CastTypeName` as the
cast target type. Otherwise, fix-it will automatically infer the type.
```

```{title} clang-tidy - bugprone-unique-ptr-array-mismatch
```

# bugprone-unique-ptr-array-mismatch

Finds initializations of C++ unique pointers to non-array type that are
initialized with an array.

If a pointer `std::unique_ptr<T>` is initialized with a new-expression
`new T[]` the memory is not deallocated correctly. A plain `delete` is used
in this case to deallocate the target memory. Instead a `delete[]` call is
needed. A `std::unique_ptr<T[]>` uses the correct delete operator. The check
does not emit warning if an `unique_ptr` with user-specified deleter type is
used.

The check offers replacement of `unique_ptr<T>` to `unique_ptr<T[]>` if it
is used at a single variable declaration (one variable in one statement).

Example:

```c++
std::unique_ptr<Foo> x(new Foo[10]); // -> std::unique_ptr<Foo[]> x(new Foo[10]);
// ^ warning: unique pointer to non-array is initialized with array
std::unique_ptr<Foo> x1(new Foo), x2(new Foo[10]); // no replacement
// ^ warning: unique pointer to non-array is initialized with array

D d;
std::unique_ptr<Foo, D> x3(new Foo[10], d); // no warning (custom deleter used)

struct S {
 std::unique_ptr<Foo> x(new Foo[10]); // no replacement in this case
 // ^ warning: unique pointer to non-array is initialized with array
};
```

This check partially covers the CERT C++ Coding Standard rule
[MEM51-CPP. Properly deallocate dynamically allocated resources](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/memory-management-mem/mem51-cpp/)
However, only the `std::unique_ptr` case is detected by this check.

```{title} clang-tidy - bugprone-unsafe-functions
```

# bugprone-unsafe-functions

Checks for functions that have safer, more secure replacements available, or
are considered deprecated due to design flaws.
The check heavily relies on the functions from the
**Annex K.** "Bounds-checking interfaces" of C11.

The check implements the following rules from the CERT C Coding Standard:

- Recommendation [MSC24-C. Do not use deprecated or obsolescent functions](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/recommendations/miscellaneous-msc/msc24-c/).
- Rule [MSC33-C. Do not pass invalid data to the asctime() function](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/miscellaneous-msc/msc33-c/).

`cert-msc24-c` and `cert-msc33-c` redirect
here as aliases of this check.

## Unsafe functions

The following functions are reported if {option}`ReportDefaultFunctions`
is enabled.

If *Annex K.* is available, a replacement from *Annex K.* is suggested for the
following functions:

`asctime`, `asctime_r`, `bsearch`, `ctime`, `fopen`, `fprintf`,
`freopen`, `fscanf`, `fwprintf`, `fwscanf`, `getenv`, `gets`,
`gmtime`, `localtime`, `mbsrtowcs`, `mbstowcs`, `memcpy`,
`memmove`, `memset`, `printf`, `qsort`, `scanf`, `snprintf`,
`sprintf`, `sscanf`, `strcat`, `strcpy`, `strerror`, `strlen`,
`strncat`, `strncpy`, `strtok`, `swprintf`, `swscanf`, `vfprintf`,
`vfscanf`, `vfwprintf`, `vfwscanf`, `vprintf`, `vscanf`,
`vsnprintf`, `vsprintf`, `vsscanf`, `vswprintf`, `vswscanf`,
`vwprintf`, `vwscanf`, `wcrtomb`, `wcscat`, `wcscpy`,
`wcslen`, `wcsncat`, `wcsncpy`, `wcsrtombs`, `wcstok`, `wcstombs`,
`wctomb`, `wmemcpy`, `wmemmove`, `wprintf`, `wscanf`.

If *Annex K.* is not available, replacements are suggested only for the
following functions from the previous list:

- `asctime`, `asctime_r`, suggested replacement: `strftime`
- `gets`, suggested replacement: `fgets`

The following functions are always checked, regardless of *Annex K*
availability:

- `rewind`, suggested replacement: `fseek`
- `setbuf`, suggested replacement: `setvbuf`
- `std::get_temporary_buffer`, suggested replacement: "plain" allocation
 with `operator new[]`

If {option}`ReportMoreUnsafeFunctions` is enabled,
the following functions are also checked:

- `bcmp`, suggested replacement: `memcmp`
- `bcopy`, suggested replacement: `memcpy_s` if *Annex K* is available,
 or `memcpy`
- `bzero`, suggested replacement: `memset_s` if *Annex K* is available,
 or `memset`
- `getpw`, suggested replacement: `getpwuid`
- `vfork`, suggested replacement: `posix_spawn`

Although mentioned in the associated CERT rules, the following functions are
**ignored** by the check:

`atof`, `atoi`, `atol`, `atoll`, `tmpfile`.

The availability of *Annex K* is determined based on the following macros:

- `__STDC_LIB_EXT1__`: feature macro, which indicates the presence of
 *Annex K. "Bounds-checking interfaces"* in the library implementation
- `__STDC_WANT_LIB_EXT1__`: user-defined macro, which indicates that the
 user requests the functions from *Annex K.* to be defined.

Both macros have to be defined to suggest replacement functions from *Annex K.*
`__STDC_LIB_EXT1__` is defined by the library implementation, and
`__STDC_WANT_LIB_EXT1__` must be defined to `1` by the user **before**
including any system headers.

(customfunctions)=

## Custom functions

The option {option}`CustomFunctions` allows the user to define custom functions
to be checked. The format is the following, without newlines:

```
bugprone-unsafe-functions.CustomFunctions="
 functionRegex1[, replacement1[, reason1]];
 functionRegex2[, replacement2[, reason2]];
 ...
"
```

The functions are matched using POSIX extended regular expressions.
*(Note: The regular expressions do not support negative* `(?!)` *matches.)*

The `reason` is optional and is used to provide additional information about
the reasoning behind the replacement. The default reason is
`is marked as unsafe`.

If `replacement` is empty, the default text
`it should not be used` will be
shown instead of the suggestion for a replacement.

If the `reason` starts with the character
`>`, the reason becomes fully
custom. The default suffix is disabled even if a
`replacement` is present,
and only the reason message is shown after the matched function, to allow
better control over the suggestions. (The starting `>` and whitespace
directly after it are trimmed from the message.)

As an example, the following configuration matches only the function
`original` in the default namespace. A similar diagnostic can also be printed
using a fully custom reason.

```c
// bugprone-unsafe-functions.CustomFunctions:
// ^original$, replacement, is deprecated;
// Using the fully custom message syntax:
// ^suspicious$,,> should be avoided if possible.
original(); // warning: function 'original' is deprecated; 'replacement' should be used instead.
suspicious(); // warning: function 'suspicious' should be avoided if possible.
::std::original(); // no-warning
original_function(); // no-warning
```

If the regular expression contains the character `:`, it is matched against
the qualified name (i.e. `std::original`), otherwise the regex is matched
against the unqualified name (`original`). If the regular expression starts
with `::` (or `^::`), it is matched against the fully qualified name
(`::std::original`).

One of the use cases for fully custom messages is suggesting compiler options
and warning flags:

```c
// bugprone-unsafe-functions.CustomFunctions:
// ^memcpy$,,>is recommended to have compiler hardening using '_FORTIFY_SOURCE';
// ^printf$,,>is recommended to have the '-Werror=format-security' compiler warning flag;

memcpy(dest, src, 999'999); // warning: function 'memcpy' is recommended to have compiler hardening using '_FORTIFY_SOURCE'
printf(raw_str); // warning: function 'printf' is recommended to have the '-Werror=format-security' compiler warning flag
```

```{note}
Fully qualified names can contain template parameters on certain C++ classes,
but not on C++ functions. Type aliases are resolved before matching.

As an example, the member function `open` in the class `std::ifstream`
has a fully qualified name of `::std::basic_ifstream<char>::open`.

The example could also be matched with the regex
`::std::basic_ifstream<[^>]*>::open`, which matches all potential template
parameters, but does not match nested template classes.
```

## Options

```{option} ReportMoreUnsafeFunctions
When `true`, additional functions from widely used APIs (such as POSIX) are
added to the list of reported functions.
See the main documentation of the check for the complete list as to what
this option enables.
Default is `true`.
```

```{option} ReportDefaultFunctions
When `true`, the check reports the default set of functions.
Consider changing the setting to false if you only want to see custom
functions matched via {ref}`custom functions<CustomFunctions>`.
Default is `true`.
```

```{option} CustomFunctions
A semicolon-separated list of custom functions to be matched. A matched
function contains a regular expression, an optional name of the replacement
function, and an optional reason, separated by comma. For more information,
see {ref}`Custom functions<CustomFunctions>`.
```

## Examples

```c++
#ifndef __STDC_LIB_EXT1__
#error "Annex K is not supported by the current standard library implementation."
#endif

#define __STDC_WANT_LIB_EXT1__ 1

#include <string.h> // Defines functions from Annex K.
#include <stdio.h>

enum { BUFSIZE = 32 };

void Unsafe(const char *Msg) {
 static const char Prefix[] = "Error: ";
 static const char Suffix[] = "\n";
 char Buf[BUFSIZE] = {0};

 strcpy(Buf, Prefix); // warning: function 'strcpy' is not bounds-checking; 'strcpy_s' should be used instead.
 strcat(Buf, Msg); // warning: function 'strcat' is not bounds-checking; 'strcat_s' should be used instead.
 strcat(Buf, Suffix); // warning: function 'strcat' is not bounds-checking; 'strcat_s' should be used instead.
 if (fputs(buf, stderr) < 0) {
 // error handling
 return;
 }
}

void UsingSafeFunctions(const char *Msg) {
 static const char Prefix[] = "Error: ";
 static const char Suffix[] = "\n";
 char Buf[BUFSIZE] = {0};

 if (strcpy_s(Buf, BUFSIZE, Prefix) != 0) {
 // error handling
 return;
 }

 if (strcat_s(Buf, BUFSIZE, Msg) != 0) {
 // error handling
 return;
 }

 if (strcat_s(Buf, BUFSIZE, Suffix) != 0) {
 // error handling
 return;
 }

 if (fputs(Buf, stderr) < 0) {
 // error handling
 return;
 }
}
```

```{title} clang-tidy - bugprone-unsafe-to-allow-exceptions
```

# bugprone-unsafe-to-allow-exceptions

Finds functions where throwing exceptions is unsafe but the function is still
marked as potentially throwing. Throwing exceptions from the following
functions can be problematic:

- Destructors
- Move constructors
- Move assignment operators
- The `main()` functions
- `swap()` functions
- `iter_swap()` functions
- `iter_move()` functions

A destructor throwing an exception may result in undefined behavior, resource
leaks or unexpected termination of the program. Throwing move constructor or
move assignment also may result in undefined behavior or resource leak. The
`swap()` operations expected to be non throwing most of the cases and they
are always possible to implement in a non throwing way. Non throwing `swap()`
operations are also used to create move operations. A throwing `main()`
function also results in unexpected termination.

The check finds any of these functions if it is marked with `noexcept(false)`
or `throw(exception)`. This would indicate that the function is expected to
throw exceptions. Only the presence of these keywords is checked, not if the
function actually throws any exception. To check if the function actually
throws exception, the check {doc}`bugprone-exception-escape <exception-escape>`
can be used (but it does not warn if a function is explicitly marked as
throwing).

## Options

```{option} CheckedSwapFunctions
Semicolon-separated list of checked swap function names (where throwing
exceptions is unsafe). These functions are checked if the parameter count is
at least 1. Default value is `swap;iter_swap;iter_move`.
```

```{title} clang-tidy - bugprone-unused-local-non-trivial-variable
```

# bugprone-unused-local-non-trivial-variable

Warns when a local non trivial variable is unused within a function.
The following types of variables are excluded from this check:

- trivial and trivially copyable
- references and pointers
- exception variables in catch clauses
- static or thread local
- structured bindings
- variables with `[[maybe_unused]]` attribute
- name-independent variables

This check can be configured to warn on all non-trivial variables by setting
{option}`IncludeTypes` to `.*`, and excluding
specific types using {option}`ExcludeTypes`.

In the this example, `my_lock` would generate a warning that
it is unused.

```c++
std::mutex my_lock;
// my_lock local variable is never used
```

In the next example, `future2` would generate a warning that
it is unused.

```c++
std::future<MyObject> future1;
std::future<MyObject> future2;
// ...
MyObject foo = future1.get();
// future2 is not used.
```

## Options

```{option} IncludeTypes
Semicolon-separated list of regular expressions matching types of variables
to check. By default the following types are checked:

- `::std::.*mutex`
- `::std::future`
- `::std::basic_string`
- `::std::basic_regex`
- `::std::basic_istringstream`
- `::std::basic_stringstream`
- `::std::bitset`
- `::std::filesystem::path`
```

```{option} ExcludeTypes
A semicolon-separated list of regular expressions matching types that are
excluded from the {option}`IncludeTypes` matches. Default is empty string.
```

```{title} clang-tidy - bugprone-unused-raii
```

# bugprone-unused-raii

Finds temporaries that look like RAII objects.

The canonical example for this is a scoped lock.

```c++
{
 scoped_lock(&global_mutex);
 critical_section();
}
```

The destructor of the scoped_lock is called before the `critical_section` is
entered, leaving it unprotected.

We apply a number of heuristics to reduce the false positive count of this
check:

- Ignore code expanded from macros. Testing frameworks make heavy use of this.

- Ignore types with trivial destructors. They are very unlikely to be RAII
 objects and there's no difference when they are deleted.

- Ignore objects at the end of a compound statement (doesn't change behavior).

- Ignore objects returned from a call.

```{title} clang-tidy - bugprone-unused-return-value
```

# bugprone-unused-return-value

Warns on unused function return values. The checked functions can be
configured.

Operator overloading with assignment semantics are ignored.

## Options

```{option} CheckedFunctions
Semicolon-separated list of functions to check.
This parameter supports regexp. The function is checked if the name
and scope matches, with any arguments.
By default the following functions are checked:
`^::std::async$, ^::std::launder$, ^::std::remove$, ^::std::remove_if$,
^::std::unique$, ^::std::unique_ptr::release$, ^::std::basic_string::empty$,
^::std::vector::empty$, ^::std::back_inserter$, ^::std::distance$,
^::std::find$, ^::std::find_if$, ^::std::inserter$, ^::std::lower_bound$,
^::std::make_pair$, ^::std::map::count$, ^::std::map::find$,
^::std::map::lower_bound$, ^::std::multimap::equal_range$,
^::std::multimap::upper_bound$, ^::std::set::count$, ^::std::set::find$,
^::std::setfill$, ^::std::setprecision$, ^::std::setw$, ^::std::upper_bound$,
^::std::vector::at$, ^::bsearch$, ^::ferror$, ^::feof$, ^::isalnum$,
^::isalpha$, ^::isblank$, ^::iscntrl$, ^::isdigit$, ^::isgraph$, ^::islower$,
^::isprint$, ^::ispunct$, ^::isspace$, ^::isupper$, ^::iswalnum$,
^::iswprint$, ^::iswspace$, ^::isxdigit$, ^::memchr$, ^::memcmp$, ^::strcmp$,
^::strcoll$, ^::strncmp$, ^::strpbrk$, ^::strrchr$, ^::strspn$, ^::strstr$,
^::wcscmp$, ^::access$, ^::bind$, ^::connect$, ^::difftime$, ^::dlsym$,
^::fnmatch$, ^::getaddrinfo$, ^::getopt$, ^::htonl$, ^::htons$,
^::iconv_open$, ^::inet_addr$, isascii$, isatty$, ^::mmap$, ^::newlocale$,
^::openat$, ^::pathconf$, ^::pthread_equal$, ^::pthread_getspecific$,
^::pthread_mutex_trylock$, ^::readdir$, ^::readlink$, ^::recvmsg$,
^::regexec$, ^::scandir$, ^::semget$, ^::setjmp$, ^::shm_open$, ^::shmget$,
^::sigismember$, ^::strcasecmp$, ^::strsignal$, ^::ttyname$`

- `std::async()`. Not using the return value makes the call synchronous.
- `std::launder()`. Not using the return value usually means that the
 function interface was misunderstood by the programmer. Only the returned
 pointer is "laundered", not the argument.
- `std::remove()`, `std::remove_if()` and `std::unique()`. The returned
 iterator indicates the boundary between elements to keep and elements to be
 removed. Not using the return value means that the information about which
 elements to remove is lost.
- `std::unique_ptr::release()`. Not using the return value can lead to
 resource leaks if the same pointer isn't stored anywhere else. Often,
 ignoring the `release()` return value indicates that the programmer
 confused the function with `reset()`.
- `std::basic_string::empty()` and `std::vector::empty()`. Not using the
 return value often indicates that the programmer confused the function with
 `clear()`.
```

```{option} CheckedReturnTypes
Semicolon-separated list of function return types to check.
By default the following function return types are checked:
`^::std::error_code$`, `^::std::error_condition$`, `^::std::errc$`,
`^::std::expected$`, `^::boost::system::error_code$`
```

```{option} AllowCastToVoid
Controls whether casting return values to `void` is permitted. Default is `false`.
```

{doc}`cert-err33-c <../cert/err33-c>` is an alias of this check that checks a
fixed and large set of standard library functions.

```{title} clang-tidy - bugprone-use-after-move
```

# bugprone-use-after-move

Warns if an object is used after it has been moved, for example:

```c++
std::string str = "Hello, world!\n";
std::vector<std::string> messages;
messages.emplace_back(std::move(str));
std::cout << str;
```

The last line will trigger a warning that `str` is used after it has been
moved.

The check does not trigger a warning if the object is reinitialized after the
move and before the use. For example, no warning will be output for this code:

```c++
messages.emplace_back(std::move(str));
str = "Greetings, stranger!\n";
std::cout << str;
```

Subsections below explain more precisely what exactly the check considers to be
a move, use, and reinitialization.

The check takes control flow into account. A warning is only emitted if the use
can be reached from the move. This means that the following code does not
produce a warning:

```c++
if (condition) {
 messages.emplace_back(std::move(str));
} else {
 std::cout << str;
}
```

On the other hand, the following code does produce a warning:

```c++
for (int i = 0; i < 10; ++i) {
 std::cout << str;
 messages.emplace_back(std::move(str));
}
```

(The use-after-move happens on the second iteration of the loop.)

In some cases, the check may not be able to detect that two branches are
mutually exclusive. For example (assuming that `i` is an int):

```c++
if (i == 1) {
 messages.emplace_back(std::move(str));
}
if (i == 2) {
 std::cout << str;
}
```

In this case, the check will erroneously produce a warning, even though it is
not possible for both the move and the use to be executed. More formally, the
analysis is [flow-sensitive but not path-sensitive](https://en.wikipedia.org/wiki/Data-flow_analysis#Sensitivities).

## Silencing erroneous warnings

An erroneous warning can be silenced by reinitializing the object after the
move:

```c++
if (i == 1) {
 messages.emplace_back(std::move(str));
 str = "";
}
if (i == 2) {
 std::cout << str;
}
```

If you want to avoid the overhead of actually reinitializing the object,
you can create a dummy function that causes the check to assume the object
was reinitialized:

```c++
template <class T>
void IS_INITIALIZED(T&) {}
```

You can use this as follows:

```c++
if (i == 1) {
 messages.emplace_back(std::move(str));
}
if (i == 2) {
 IS_INITIALIZED(str);
 std::cout << str;
}
```

The check will not output a warning in this case because passing the object
to a function as a non-const pointer or reference counts as a reinitialization
(see section [Reinitialization](#reinitialization) below).

## Unsequenced moves, uses, and reinitializations

In many cases, C++ does not make any guarantees about the order in which
sub-expressions of a statement are evaluated. This means that in code like the
following, it is not guaranteed whether the use will happen before or after the
move:

```c++
void f(int i, std::vector<int> v);
std::vector<int> v = { 1, 2, 3 };
f(v[1], std::move(v));
```

In this kind of situation, the check will note that the use and move are
unsequenced.

The check will also take sequencing rules into account when reinitializations
occur in the same statement as moves or uses. A reinitialization is only
considered to reinitialize a variable if it is guaranteed to be evaluated after
the move and before the use.

## Move

The check currently only considers calls of `std::move` on local variables or
function parameters. It does not check moves of member variables or global
variables.

Any call of `std::move` on a variable is considered to cause a move of that
variable, even if the result of `std::move` is not passed to an rvalue
reference parameter.

This means that the check will flag a use-after-move even on a type that does
not define a move constructor or move assignment operator. This is intentional.
Developers may use `std::move` on such a type in the expectation that the
type will add move semantics in the future. If such a `std::move` has the
potential to cause a use-after-move, we want to warn about it even if the type
does not implement move semantics yet.

Furthermore, if the result of `std::move` *is* passed to an rvalue reference
parameter, this will always be considered to cause a move, even if the function
that consumes this parameter does not move from it, or if it does so only
conditionally. For example, in the following situation, the check will assume
that a move always takes place:

```c++
std::vector<std::string> messages;
void f(std::string &&str) {
 // Only remember the message if it isn't empty.
 if (!str.empty()) {
 messages.emplace_back(std::move(str));
 }
}
std::string str = "";
f(std::move(str));
```

The check will assume that the last line causes a move, even though, in this
particular case, it does not. Again, this is intentional.

There is one special case: A call to `std::move` inside a `try_emplace`
call is conservatively assumed not to move. This is to avoid spurious warnings,
as the check has no way to reason about the `bool` returned by `try_emplace`.

When analyzing the order in which moves, uses and reinitializations happen (see
section [Unsequenced moves, uses, and
reinitializations](#unsequenced-moves-uses-and-reinitializations)), the move is
assumed
to occur in whichever function the result of the `std::move` is passed to.

The check also handles perfect-forwarding with `std::forward` so the
following code will also trigger a use-after-move warning.

```c++
void consume(int);

void f(int&& i) {
 consume(std::forward<int>(i));
 consume(std::forward<int>(i)); // use-after-move
}
```

## Use

Any occurrence of the moved variable that is not a reinitialization (see below)
or an explicit call to the variable destructor is considered to be a use.

An exception to this are objects of type `std::unique_ptr`,
`std::shared_ptr`, `std::weak_ptr`, `std::optional`, and `std::any`,
which can be reinitialized via `reset`. For smart pointers specifically, the
moved-from objects have a well-defined state of being `nullptr`s, and only
`operator*`, `operator->` and `operator[]` are considered bad accesses as
they would be dereferencing a `nullptr`.

User-defined types can be annotated as having the same semantics as standard
smart pointers with `[[clang::annotate("clang-tidy",
"bugprone-use-after-move", "null_after_move")]]`. This expresses that a
moved-from object of this type is a null pointer.

If multiple uses occur after a move, only the first of these is flagged.

## Reinitialization

The check considers a variable to be reinitialized in the following cases:

- The variable occurs on the left-hand side of an assignment.
- The variable is passed to a function as a non-const pointer or non-const
 lvalue reference. (It is assumed that the variable may be an out-parameter
 for the function.)
- `clear()` or `assign()` is called on the variable and the variable is
 of one of the standard container types `basic_string`, `vector`,
 `deque`, `forward_list`, `list`, `set`, `map`, `multiset`,
 `multimap`, `unordered_set`, `unordered_map`, `unordered_multiset`,
 `unordered_multimap`.
- `reset()` is called on the variable and the variable is of type
 `std::unique_ptr`, `std::shared_ptr`, `std::weak_ptr`,
 `std::optional`, or `std::any`.
- A member function marked with the `[[clang::reinitializes]]` attribute is
 called on the variable.
- The variable is passed as an argument to `std::tie` on the left-hand
 side of an assignment (e.g. `std::tie(a, b) = f(...)`). The tuple
 assignment operator writes back through the stored references, which
 reinitializes each named variable.

If the variable in question is a struct and an individual member variable of
that struct is written to, the check does not consider this to be a
reinitialization -- even if, eventually, all member variables of the struct are
written to. For example:

```c++
struct S {
 std::string str;
 int i;
};
S s = { "Hello, world!\n", 42 };
S s_other = std::move(s);
s.str = "Lorem ipsum";
s.i = 99;
```

The check will not consider `s` to be reinitialized after the last line;
instead, the line that assigns to `s.str` will be flagged as a use-after-move.
This is intentional as this pattern of reinitializing a struct is error-prone.
For example, if an additional member variable is added to `S`, it is easy to
forget to add the reinitialization for this additional member. Instead, it is
safer to assign to the entire struct in one go, and this will also avoid the
use-after-move warning.

## Options

```{option} InvalidationFunctions
A semicolon-separated list of regular expressions matching names of functions
that cause their first arguments to be invalidated (e.g., closing a handle).
For member functions, the first argument is considered to be the implicit
object argument (`this`). Default value is an empty string.
```

```{option} ReinitializationFunctions
A semicolon-separated list of regular expressions matching names of functions
that reinitialize the object. For member functions, the implicit object
argument (`*this`) is considered to be reinitialized. For non-member or
static member functions, the first argument is considered to be
reinitialized. Default value is an empty string.
```

```{title} clang-tidy - bugprone-virtual-near-miss
```

# bugprone-virtual-near-miss

Warn if a function is a near miss (i.e. the name is very similar and
the function signature is the same) to a virtual function from a base
class.

Example:

```cpp
struct Base {
 virtual void func();
};

struct Derived : Base {
 virtual void funk();
 // warning: 'Derived::funk' has a similar name and the same signature as virtual method 'Base::func'; did you mean to override it?
};
```

```{title} clang-tidy - cert-arr39-c
```

# cert-arr39-c

The `cert-arr39-c` check is an alias, please see
{doc}`bugprone-sizeof-expression <../bugprone/sizeof-expression>`
for more information.

```{title} clang-tidy - cert-con36-c
```

# cert-con36-c

The `cert-con36-c` check is an alias, please see
{doc}`bugprone-spuriously-wake-up-functions <../bugprone/spuriously-wake-up-functions>`
for more information.

```{title} clang-tidy - cert-con54-cpp
```

# cert-con54-cpp

The `cert-con54-cpp` check is an alias, please see
{doc}`bugprone-spuriously-wake-up-functions <../bugprone/spuriously-wake-up-functions>`
for more information.

```{title} clang-tidy - cert-ctr56-cpp
```

# cert-ctr56-cpp

The `cert-ctr56-cpp` check is an alias, please see
{doc}`bugprone-pointer-arithmetic-on-polymorphic-object <../bugprone/pointer-arithmetic-on-polymorphic-object>`
for more information.

```{title} clang-tidy - cert-dcl03-c
```

# cert-dcl03-c

The `cert-dcl03-c` check is an alias, please see
{doc}`misc-static-assert <../misc/static-assert>`
for more information.

```{title} clang-tidy - cert-dcl16-c
```

# cert-dcl16-c

The `cert-dcl16-c` check is an alias, please see
{doc}`readability-uppercase-literal-suffix <../readability/uppercase-literal-suffix>` for more information.

```{title} clang-tidy - cert-dcl37-c
```

# cert-dcl37-c

The `cert-dcl37-c` check is an alias, please see
{doc}`bugprone-reserved-identifier <../bugprone/reserved-identifier>` for more
information.

```{title} clang-tidy - cert-dcl50-cpp
```

# cert-dcl50-cpp

The `cert-dcl50-cpp` check is an alias, please see
{doc}`modernize-avoid-variadic-functions <../modernize/avoid-variadic-functions>`
for more information.

This check corresponds to the CERT C++ Coding Standard rule
[DCL50-CPP. Do not define a C-style variadic function](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/declarations-and-initialization-dcl/dcl50-cpp/).

```{title} clang-tidy - cert-dcl51-cpp
```

# cert-dcl51-cpp

The `cert-dcl51-cpp` check is an alias, please see
{doc}`bugprone-reserved-identifier <../bugprone/reserved-identifier>` for more information.

```{title} clang-tidy - cert-dcl54-cpp
```

# cert-dcl54-cpp

The `cert-dcl54-cpp` check is an alias, please see
{doc}`misc-new-delete-overloads <../misc/new-delete-overloads>` for more information.

```{title} clang-tidy - cert-dcl58-cpp
```

# cert-dcl58-cpp

The `cert-dcl58-cpp` is an alias, please see
{doc}`bugprone-std-namespace-modification <../bugprone/std-namespace-modification>`
for more information.

This check corresponds to the CERT C++ Coding Standard rule
[DCL58-CPP. Do not modify the standard namespaces](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/declarations-and-initialization-dcl/dcl58-cpp/).

```{title} clang-tidy - cert-dcl59-cpp
```

# cert-dcl59-cpp

The `cert-dcl59-cpp` check is an alias, please see
{doc}`misc-anonymous-namespace-in-header <../misc/anonymous-namespace-in-header>`
for more information.

```{title} clang-tidy - cert-env33-c
```

# cert-env33-c

The `cert-env33-c` check is an alias, please see
{doc}`bugprone-command-processor <../bugprone/command-processor>`
for more information.

This check corresponds to the CERT C Coding Standard rule
[ENV33-C. Do not call system()](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/environment-env/env33-c/).

```{title} clang-tidy - cert-err09-cpp
```

# cert-err09-cpp

The `cert-err09-cpp` check is an alias, please see
{doc}`misc-throw-by-value-catch-by-reference <../misc/throw-by-value-catch-by-reference>`
for more information.

This check corresponds to the CERT C++ Coding Standard recommendation
ERR09-CPP. Throw anonymous temporaries. However, all of the CERT
recommendations have been removed from public view, and so their
justification for the behavior of this check requires an account on
their wiki to view.

**Title:** clang-tidy - cert-err33-c

# cert-err33-c

Warns on unused function return values. Many of the standard library functions
return a value that indicates if the call was successful. Ignoring the returned
value can cause unexpected behavior if an error has occurred. The following
functions are checked:

* aligned_alloc()
* asctime_s()
* at_quick_exit()
* atexit()
* bsearch()
* bsearch_s()
* btowc()
* c16rtomb()
* c32rtomb()
* calloc()
* clock()
* cnd_broadcast()
* cnd_init()
* cnd_signal()
* cnd_timedwait()
* cnd_wait()
* ctime_s()
* fclose()
* fflush()
* fgetc()
* fgetpos()
* fgets()
* fgetwc()
* fopen()
* fopen_s()
* fprintf()
* fprintf_s()
* fputc()
* fputs()
* fputwc()
* fputws()
* fread()
* freopen()
* freopen_s()
* fscanf()
* fscanf_s()
* fseek()
* fsetpos()
* ftell()
* fwprintf()
* fwprintf_s()
* fwrite()
* fwscanf()
* fwscanf_s()
* getc()
* getchar()
* getenv()
* getenv_s()
* gets_s()
* getwc()
* getwchar()
* gmtime()
* gmtime_s()
* localtime()
* localtime_s()
* malloc()
* mbrtoc16()
* mbrtoc32()
* mbsrtowcs()
* mbsrtowcs_s()
* mbstowcs()
* mbstowcs_s()
* memchr()
* mktime()
* mtx_init()
* mtx_lock()
* mtx_timedlock()
* mtx_trylock()
* mtx_unlock()
* printf_s()
* putc()
* putwc()
* raise()
* realloc()
* remove()
* rename()
* setlocale()
* setvbuf()
* scanf()
* scanf_s()
* signal()
* snprintf()
* snprintf_s()
* sprintf()
* sprintf_s()
* sscanf()
* sscanf_s()
* strchr()
* strerror_s()
* strftime()
* strpbrk()
* strrchr()
* strstr()
* strtod()
* strtof()
* strtoimax()
* strtok()
* strtok_s()
* strtol()
* strtold()
* strtoll()
* strtoumax()
* strtoul()
* strtoull()
* strxfrm()
* swprintf()
* swprintf_s()
* swscanf()
* swscanf_s()
* thrd_create()
* thrd_detach()
* thrd_join()
* thrd_sleep()
* time()
* timespec_get()
* tmpfile()
* tmpfile_s()
* tmpnam()
* tmpnam_s()
* tss_create()
* tss_get()
* tss_set()
* ungetc()
* ungetwc()
* vfprintf()
* vfprintf_s()
* vfscanf()
* vfscanf_s()
* vfwprintf()
* vfwprintf_s()
* vfwscanf()
* vfwscanf_s()
* vprintf_s()
* vscanf()
* vscanf_s()
* vsnprintf()
* vsnprintf_s()
* vsprintf()
* vsprintf_s()
* vsscanf()
* vsscanf_s()
* vswprintf()
* vswprintf_s()
* vswscanf()
* vswscanf_s()
* vwprintf_s()
* vwscanf()
* vwscanf_s()
* wcrtomb()
* wcschr()
* wcsftime()
* wcspbrk()
* wcsrchr()
* wcsrtombs()
* wcsrtombs_s()
* wcsstr()
* wcstod()
* wcstof()
* wcstoimax()
* wcstok()
* wcstok_s()
* wcstol()
* wcstold()
* wcstoll()
* wcstombs()
* wcstombs_s()
* wcstoumax()
* wcstoul()
* wcstoull()
* wcsxfrm()
* wctob()
* wctrans()
* wctype()
* wmemchr()
* wprintf_s()
* wscanf()
* wscanf_s()

This check is an alias of check :doc:`bugprone-unused-return-value
<../bugprone/unused-return-value>` with a fixed set of functions.

Suppressing issues by casting to `void` is enabled by default and can be
disabled by setting `AllowCastToVoid` option to `false`.

The check corresponds to a part of CERT C Coding Standard rule `ERR33-C.
Detect and handle standard library errors
<https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/error-handling-err/err33-c/>`_.
The list of checked functions is taken from the rule, with following exception:

* The check can not differentiate if a function is called with `NULL`
 argument. Therefore the following functions are not checked:
 `mblen`, `mbrlen`, `mbrtowc`, `mbtowc`, `wctomb`, `wctomb_s`

```{title} clang-tidy - cert-err34-c
```

# cert-err34-c

The cert-err34-c check is an alias, please see
{doc}`bugprone-unchecked-string-to-number-conversion <../bugprone/unchecked-string-to-number-conversion>`
for more information.

```{title} clang-tidy - cert-err52-cpp
```

# cert-err52-cpp

The `cert-err52-cpp` check is an alias, please see
{doc}`modernize-avoid-setjmp-longjmp <../modernize/avoid-setjmp-longjmp>`
for more information.

This check corresponds to the CERT C++ Coding Standard rule
[ERR52-CPP. Do not use setjmp() or longjmp()](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/exceptions-and-error-handling-err/err52-cpp/).

```{title} clang-tidy - cert-err58-cpp
```

# cert-err58-cpp

The `cert-err58-cpp` check is an alias, please see
{doc}`bugprone-throwing-static-initialization <../bugprone/throwing-static-initialization>`
for more information.

This check corresponds to the CERT C++ Coding Standard rule
[ERR58-CPP. Handle all exceptions thrown before main() begins executing](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/exceptions-and-error-handling-err/err58-cpp/).

```{title} clang-tidy - cert-err60-cpp
```

# cert-err60-cpp

The `cert-err60-cpp` check is an alias, please see
{doc}`bugprone-exception-copy-constructor-throws <../bugprone/exception-copy-constructor-throws>`
for more information.

This check corresponds to the CERT C++ Coding Standard rule
[ERR60-CPP. Exception objects must be nothrow copy constructible](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/exceptions-and-error-handling-err/err60-cpp/).

```{title} clang-tidy - cert-err61-cpp
```

# cert-err61-cpp

The `cert-err61-cpp` check is an alias, please see
{doc}`misc-throw-by-value-catch-by-reference <../misc/throw-by-value-catch-by-reference>`
for more information.

```{title} clang-tidy - cert-exp42-c
```

# cert-exp42-c

The `cert-exp42-c` check is an alias, please see
{doc}`bugprone-suspicious-memory-comparison <../bugprone/suspicious-memory-comparison>` for more information.

```{title} clang-tidy - cert-exp45-c
```

# cert-exp45-c

The `cert-exp45-c` check is an alias, please see
{doc}`bugprone-assignment-in-selection-statement <../bugprone/assignment-in-selection-statement>` for more information.

```{title} clang-tidy - cert-fio38-c
```

# cert-fio38-c

The `cert-fio38-c` check is an alias, please see
{doc}`misc-non-copyable-objects <../misc/non-copyable-objects>` for more information.

This check corresponds to CERT C++ Coding Standard rule [FIO38-C. Do not copy a FILE object](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/input-output-fio/fio38-c/).

```{title} clang-tidy - cert-flp30-c
```

# cert-flp30-c

The `cert-flp30-c` check is an alias, please see
{doc}`bugprone-float-loop-counter <../bugprone/float-loop-counter>`
for more information

This check corresponds to the CERT C Coding Standard rule
[FLP30-C. Do not use floating-point variables as loop counters](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/floating-point-flp/flp30-c/).

```{title} clang-tidy - cert-flp37-c
```

# cert-flp37-c

The `cert-flp37-c` check is an alias, please see
{doc}`bugprone-suspicious-memory-comparison <../bugprone/suspicious-memory-comparison>` for more information.

```{title} clang-tidy - cert-int09-c
```

# cert-int09-c

The `cert-int09-c` check is an alias, please see
{doc}`readability-enum-initial-value <../readability/enum-initial-value>` for more information.

```{title} clang-tidy - cert-mem57-cpp
```

# cert-mem57-cpp

The `cert-mem57-cpp` is an alias, please see
{doc}`bugprone-default-operator-new-on-overaligned-type <../bugprone/default-operator-new-on-overaligned-type>`
for more information.

This check corresponds to the CERT C++ Coding Standard rule
[MEM57-CPP. Avoid using default operator new for over-aligned types](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/memory-management-mem/mem57-cpp/).

```{title} clang-tidy - cert-msc24-c
```

# cert-msc24-c

The `cert-msc24-c` check is an alias, please see
{doc}`bugprone-unsafe-functions <../bugprone/unsafe-functions>` for more information.

```{title} clang-tidy - cert-msc30-c
```

# cert-msc30-c

The `cert-msc30-c` check is an alias, please see
{doc}`misc-predictable-rand <../misc/predictable-rand>` for more information.

This check corresponds to the CERT C Coding Standard rule
[MSC30-C. Do not use the rand() function for generating pseudorandom numbers](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/miscellaneous-msc/msc30-c/).

```{title} clang-tidy - cert-msc32-c
```

# cert-msc32-c

The `cert-msc32-c` check is an alias, please see
{doc}`bugprone-random-generator-seed <../bugprone/random-generator-seed>`
for more information.

This check corresponds to the CERT C Coding Standard rule
[MSC32-C. Properly seed pseudorandom number generators](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/miscellaneous-msc/msc32-c/).

```{title} clang-tidy - cert-msc33-c
```

# cert-msc33-c

The `cert-msc33-c` check is an alias, please see
{doc}`bugprone-unsafe-functions <../bugprone/unsafe-functions>` for more information.

```{title} clang-tidy - cert-msc50-cpp
```

# cert-msc50-cpp

The `cert-msc50-cpp` check is an alias, please see
{doc}`misc-predictable-rand <../misc/predictable-rand>` for more information.

This check corresponds to the CERT C Coding Standard rule
[MSC50-CPP. Do not use std::rand() for generating pseudorandom numbers](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/miscellaneous-msc/msc50-cpp).

```{title} clang-tidy - cert-msc51-cpp
```

# cert-msc51-cpp

The `cert-msc51-cpp` check is an alias, please see
{doc}`bugprone-random-generator-seed <../bugprone/random-generator-seed>`
for more information.

This check corresponds to the CERT C++ Coding Standard rule
[MSC51-CPP. Ensure your random number generator is properly seeded](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/miscellaneous-msc/msc51-cpp/).

```{title} clang-tidy - cert-msc54-cpp
```

# cert-msc54-cpp

The `cert-msc54-cpp` check is an alias, please see
{doc}`bugprone-signal-handler <../bugprone/signal-handler>`
for more information.

```{title} clang-tidy - cert-oop11-cpp
```

# cert-oop11-cpp

The `cert-oop11-cpp` check is an alias, please see
{doc}`performance-move-constructor-init <../performance/move-constructor-init>`
for more information.

This check corresponds to the CERT C++ Coding Standard recommendation
OOP11-CPP. Do not copy-initialize members or base classes from a move
constructor. However, all of the CERT recommendations have been removed from
public view, and so their justification for the behavior of this check requires
an account on their wiki to view.

```{title} clang-tidy - cert-oop54-cpp
```

# cert-oop54-cpp

The `cert-oop54-cpp` check is an alias, please see
{doc}`bugprone-unhandled-self-assignment <../bugprone/unhandled-self-assignment>`
for more information.

```{title} clang-tidy - cert-oop57-cpp
```

# cert-oop57-cpp

The `cert-oop57-cpp` check is an alias, please see
{doc}`bugprone-raw-memory-call-on-non-trivial-type <../bugprone/raw-memory-call-on-non-trivial-type>`
for more information.

This check corresponds to the CERT C++ Coding Standard rule
[OOP57-CPP. Prefer special member functions and overloaded operators to C
Standard Library functions](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/object-oriented-programming-oop/oop57-cpp/).

```{title} clang-tidy - cert-oop58-cpp
```

# cert-oop58-cpp

The `cert-oop58-cpp` check is an alias, please see
{doc}`bugprone-copy-constructor-mutates-argument <../bugprone/copy-constructor-mutates-argument>`
for more information.

```{title} clang-tidy - cert-pos44-c
```

# cert-pos44-c

The `cert-pos44-c` check is an alias, please see
{doc}`bugprone-bad-signal-to-kill-thread <../bugprone/bad-signal-to-kill-thread>` for more information.

```{title} clang-tidy - cert-pos47-c
```

# cert-pos47-c

The `cert-pos47-c` check is an alias, please see
{doc}`concurrency-thread-canceltype-asynchronous <../concurrency/thread-canceltype-asynchronous>` for more information.

```{title} clang-tidy - cert-sig30-c
```

# cert-sig30-c

The `cert-sig30-c` check is an alias, please see {doc}`bugprone-signal-handler <../bugprone/signal-handler>` for more information.

```{title} clang-tidy - cert-str34-c
```

# cert-str34-c

The `cert-str34-c` check is an alias, please see
{doc}`bugprone-signed-char-misuse <../bugprone/signed-char-misuse>`
for more information.

```{title} clang-tidy - clang-analyzer-core.BitwiseShift
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-bitwiseshift
```

# clang-analyzer-core.BitwiseShift

Finds cases where bitwise shift operation causes undefined behaviour.

The `clang-analyzer-core.BitwiseShift` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-bitwiseshift)
for more information.

```{title} clang-tidy - clang-analyzer-core.CallAndMessage
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-callandmessage
```

# clang-analyzer-core.CallAndMessage

Check for logical errors for function calls and Objective-C message expressions
(e.g., uninitialized arguments, null function pointers).

The `clang-analyzer-core.CallAndMessage` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-callandmessage)
for more information.

```{title} clang-tidy - clang-analyzer-core.DivideZero
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-dividezero
```

# clang-analyzer-core.DivideZero

Check for division by zero.

The `clang-analyzer-core.DivideZero` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-dividezero)
for more information.

```{title} clang-tidy - clang-analyzer-core.NonNullParamChecker
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-nonnullparamchecker
```

# clang-analyzer-core.NonNullParamChecker

Check for null pointers passed as arguments to a function whose arguments are
references or marked with the 'nonnull' attribute.

The `clang-analyzer-core.NonNullParamChecker` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-nonnullparamchecker)
for more information.

```{title} clang-tidy - clang-analyzer-core.NullDereference
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-nulldereference
```

# clang-analyzer-core.NullDereference

Check for dereferences of null pointers.

The `clang-analyzer-core.NullDereference` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-nulldereference)
for more information.

```{title} clang-tidy - clang-analyzer-core.NullPointerArithm
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-nullpointerarithm
```

# clang-analyzer-core.NullPointerArithm

Check for undefined arithmetic operations on null pointers.

The `clang-analyzer-core.NullPointerArithm` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-nullpointerarithm)
for more information.

```{title} clang-tidy - clang-analyzer-core.StackAddressEscape
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-stackaddressescape
```

# clang-analyzer-core.StackAddressEscape

Check that addresses to stack memory do not escape the function.

The `clang-analyzer-core.StackAddressEscape` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-stackaddressescape)
for more information.

```{title} clang-tidy - clang-analyzer-core.UndefinedBinaryOperatorResult
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-undefinedbinaryoperatorresult
```

# clang-analyzer-core.UndefinedBinaryOperatorResult

Check for undefined results of binary operators.

The `clang-analyzer-core.UndefinedBinaryOperatorResult` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-undefinedbinaryoperatorresult)
for more information.

```{title} clang-tidy - clang-analyzer-core.VLASize
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-vlasize
```

# clang-analyzer-core.VLASize

Check for declarations of VLA of undefined or zero size.

The `clang-analyzer-core.VLASize` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-vlasize)
for more information.

```{title} clang-tidy - clang-analyzer-core.uninitialized.ArraySubscript
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-arraysubscript
```

# clang-analyzer-core.uninitialized.ArraySubscript

Check for uninitialized values used as array subscripts.

The `clang-analyzer-core.uninitialized.ArraySubscript` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-arraysubscript)
for more information.

```{title} clang-tidy - clang-analyzer-core.uninitialized.Assign
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-assign
```

# clang-analyzer-core.uninitialized.Assign

Check for assigning uninitialized values.

The `clang-analyzer-core.uninitialized.Assign` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-assign)
for more information.

```{title} clang-tidy - clang-analyzer-core.uninitialized.Branch
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-branch
```

# clang-analyzer-core.uninitialized.Branch

Check for uninitialized values used as branch conditions.

The `clang-analyzer-core.uninitialized.Branch` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-branch)
for more information.

```{title} clang-tidy - clang-analyzer-core.uninitialized.CapturedBlockVariable
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-capturedblockvariable
```

# clang-analyzer-core.uninitialized.CapturedBlockVariable

Check for blocks that capture uninitialized values.

The `clang-analyzer-core.uninitialized.CapturedBlockVariable` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-capturedblockvariable)
for more information.

```{title} clang-tidy - clang-analyzer-core.uninitialized.NewArraySize
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-newarraysize
```

# clang-analyzer-core.uninitialized.NewArraySize

Check if the size of the array in a new[] expression is undefined.

The `clang-analyzer-core.uninitialized.NewArraySize` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-newarraysize)
for more information.

```{title} clang-tidy - clang-analyzer-core.uninitialized.UndefReturn
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-undefreturn
```

# clang-analyzer-core.uninitialized.UndefReturn

Check for uninitialized values being returned to the caller.

The `clang-analyzer-core.uninitialized.UndefReturn` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#core-uninitialized-undefreturn)
for more information.

```{title} clang-tidy - clang-analyzer-cplusplus.ArrayDelete
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-arraydelete
```

# clang-analyzer-cplusplus.ArrayDelete

Reports destructions of arrays of polymorphic objects that are destructed as
their base class.

The `clang-analyzer-cplusplus.ArrayDelete` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-arraydelete)
for more information.

```{title} clang-tidy - clang-analyzer-cplusplus.InnerPointer
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-innerpointer
```

# clang-analyzer-cplusplus.InnerPointer

Check for inner pointers of C++ containers used after re/deallocation.

The `clang-analyzer-cplusplus.InnerPointer` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-innerpointer)
for more information.

```{title} clang-tidy - clang-analyzer-cplusplus.Move
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-move
```

# clang-analyzer-cplusplus.Move

Find use-after-move bugs in C++.

The `clang-analyzer-cplusplus.Move` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-move)
for more information.

```{title} clang-tidy - clang-analyzer-cplusplus.NewDelete
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-newdelete
```

# clang-analyzer-cplusplus.NewDelete

Check for double-free and use-after-free problems. Traces memory managed by
new/delete.

The `clang-analyzer-cplusplus.NewDelete` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-newdelete)
for more information.

```{title} clang-tidy - clang-analyzer-cplusplus.NewDeleteLeaks
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-newdeleteleaks
```

# clang-analyzer-cplusplus.NewDeleteLeaks

Check for memory leaks. Traces memory managed by new/delete.

The `clang-analyzer-cplusplus.NewDeleteLeaks` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-newdeleteleaks)
for more information.

```{title} clang-tidy - clang-analyzer-cplusplus.PlacementNew
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-placementnew
```

# clang-analyzer-cplusplus.PlacementNew

Check if default placement new is provided with pointers to sufficient storage
capacity.

The `clang-analyzer-cplusplus.PlacementNew` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-placementnew)
for more information.

```{title} clang-tidy - clang-analyzer-cplusplus.PureVirtualCall
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-purevirtualcall
```

# clang-analyzer-cplusplus.PureVirtualCall

Check pure virtual function calls during construction/destruction.

The `clang-analyzer-cplusplus.PureVirtualCall` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-purevirtualcall)
for more information.

```{title} clang-tidy - clang-analyzer-cplusplus.SelfAssignment
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-selfassignment
```

# clang-analyzer-cplusplus.SelfAssignment

Checks C++ copy and move assignment operators for self assignment.

The `clang-analyzer-cplusplus.SelfAssignment` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-selfassignment)
for more information.

```{title} clang-tidy - clang-analyzer-cplusplus.StringChecker
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-stringchecker
```

# clang-analyzer-cplusplus.StringChecker

Checks C++ std::string bugs.

The `clang-analyzer-cplusplus.StringChecker` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#cplusplus-stringchecker)
for more information.

```{title} clang-tidy - clang-analyzer-deadcode.DeadStores
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#deadcode-deadstores
```

# clang-analyzer-deadcode.DeadStores

Check for values stored to variables that are never read afterwards.

The `clang-analyzer-deadcode.DeadStores` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#deadcode-deadstores)
for more information.

```{title} clang-tidy - clang-analyzer-fuchsia.HandleChecker
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#fuchsia-handlechecker
```

# clang-analyzer-fuchsia.HandleChecker

A Checker that detect leaks related to Fuchsia handles.

The `clang-analyzer-fuchsia.HandleChecker` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#fuchsia-handlechecker)
for more information.

```{title} clang-tidy - clang-analyzer-nullability.NullPassedToNonnull
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullpassedtononnull
```

# clang-analyzer-nullability.NullPassedToNonnull

Warns when a null pointer is passed to a pointer which has a _Nonnull type.

The `clang-analyzer-nullability.NullPassedToNonnull` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullpassedtononnull)
for more information.

```{title} clang-tidy - clang-analyzer-nullability.NullReturnedFromNonnull
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullreturnedfromnonnull
```

# clang-analyzer-nullability.NullReturnedFromNonnull

Warns when a null pointer is returned from a function that has _Nonnull return
type.

The `clang-analyzer-nullability.NullReturnedFromNonnull` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullreturnedfromnonnull)
for more information.

```{title} clang-tidy - clang-analyzer-nullability.NullableDereferenced
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullabledereferenced
```

# clang-analyzer-nullability.NullableDereferenced

Warns when a nullable pointer is dereferenced.

The `clang-analyzer-nullability.NullableDereferenced` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullabledereferenced)
for more information.

```{title} clang-tidy - clang-analyzer-nullability.NullablePassedToNonnull
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullablepassedtononnull
```

# clang-analyzer-nullability.NullablePassedToNonnull

Warns when a nullable pointer is passed to a pointer which has a _Nonnull type.

The `clang-analyzer-nullability.NullablePassedToNonnull` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullablepassedtononnull)
for more information.

```{title} clang-tidy - clang-analyzer-nullability.NullableReturnedFromNonnull
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullablereturnedfromnonnull
```

# clang-analyzer-nullability.NullableReturnedFromNonnull

Warns when a nullable pointer is returned from a function that has _Nonnull
return type.

The `clang-analyzer-nullability.NullableReturnedFromNonnull` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#nullability-nullablereturnedfromnonnull)
for more information.

```{title} clang-tidy - clang-analyzer-optin.core.EnumCastOutOfRange
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-core-enumcastoutofrange
```

# clang-analyzer-optin.core.EnumCastOutOfRange

Check integer to enumeration casts for out of range values.

The `clang-analyzer-optin.core.EnumCastOutOfRange` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-core-enumcastoutofrange)
for more information.

```{title} clang-tidy - clang-analyzer-optin.core.FixedAddressDereference
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-core-fixedaddressdereference
```

# clang-analyzer-optin.core.FixedAddressDereference

Check for dereferences of fixed addresses.

The `clang-analyzer-optin.core.FixedAddressDereference` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-core-fixedaddressdereference)
for more information.

```{title} clang-tidy - clang-analyzer-optin.core.UnconditionalVAArg
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-core-unconditionalvaarg
```

# clang-analyzer-optin.core.UnconditionalVAArg

Check variadic functions unconditionally using va_arg.

The `clang-analyzer-optin.core.UnconditionalVAArg` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-core-unconditionalvaarg)
for more information.

```{title} clang-tidy - clang-analyzer-optin.cplusplus.UninitializedObject
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-cplusplus-uninitializedobject
```

# clang-analyzer-optin.cplusplus.UninitializedObject

Reports uninitialized fields after object construction.

The `clang-analyzer-optin.cplusplus.UninitializedObject` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-cplusplus-uninitializedobject)
for more information.

```{title} clang-tidy - clang-analyzer-optin.cplusplus.VirtualCall
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-cplusplus-virtualcall
```

# clang-analyzer-optin.cplusplus.VirtualCall

Check virtual function calls during construction/destruction.

The `clang-analyzer-optin.cplusplus.VirtualCall` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-cplusplus-virtualcall)
for more information.

```{title} clang-tidy - clang-analyzer-optin.mpi.MPI-Checker
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-mpi-mpi-checker
```

# clang-analyzer-optin.mpi.MPI-Checker

Checks MPI code.

The `clang-analyzer-optin.mpi.MPI-Checker` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-mpi-mpi-checker)
for more information.

```{title} clang-tidy - clang-analyzer-optin.osx.OSObjectCStyleCast
```

# clang-analyzer-optin.osx.OSObjectCStyleCast

Checker for C-style casts of OSObjects.

The clang-analyzer-optin.osx.OSObjectCStyleCast check is an alias of
Clang Static Analyzer optin.osx.OSObjectCStyleCast.

```{title} clang-tidy - clang-analyzer-optin.osx.cocoa.localizability.EmptyLocalizationContextChecker
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-osx-cocoa-localizability-emptylocalizationcontextchecker
```

# clang-analyzer-optin.osx.cocoa.localizability.EmptyLocalizationContextChecker

Check that NSLocalizedString macros include a comment for context.

The `clang-analyzer-optin.osx.cocoa.localizability.EmptyLocalizationContextChecker` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-osx-cocoa-localizability-emptylocalizationcontextchecker)
for more information.

```{title} clang-tidy - clang-analyzer-optin.osx.cocoa.localizability.NonLocalizedStringChecker
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-osx-cocoa-localizability-nonlocalizedstringchecker
```

# clang-analyzer-optin.osx.cocoa.localizability.NonLocalizedStringChecker

Warns about uses of non-localized NSStrings passed to UI methods expecting
localized NSStrings.

The `clang-analyzer-optin.osx.cocoa.localizability.NonLocalizedStringChecker` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-osx-cocoa-localizability-nonlocalizedstringchecker)
for more information.

```{title} clang-tidy - clang-analyzer-optin.performance.GCDAntipattern
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-performance-gcdantipattern
```

# clang-analyzer-optin.performance.GCDAntipattern

Check for performance anti-patterns when using Grand Central Dispatch.

The `clang-analyzer-optin.performance.GCDAntipattern` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-performance-gcdantipattern)
for more information.

```{title} clang-tidy - clang-analyzer-optin.performance.Padding
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-performance-padding
```

# clang-analyzer-optin.performance.Padding

Check for excessively padded structs.

The `clang-analyzer-optin.performance.Padding` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-performance-padding)
for more information.

```{title} clang-tidy - clang-analyzer-optin.portability.UnixAPI
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-portability-unixapi
```

# clang-analyzer-optin.portability.UnixAPI

Finds dynamic memory allocation with size zero.

The `clang-analyzer-optin.portability.UnixAPI` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-portability-unixapi)
for more information.

```{title} clang-tidy - clang-analyzer-optin.taint.GenericTaint
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-taint-generictaint
```

# clang-analyzer-optin.taint.GenericTaint

Reports potential injection vulnerabilities.

The `clang-analyzer-optin.taint.GenericTaint` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-taint-generictaint)
for more information.

```{title} clang-tidy - clang-analyzer-optin.taint.TaintedAlloc
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-taint-taintedalloc
```

# clang-analyzer-optin.taint.TaintedAlloc

Check for memory allocations, where the size parameter might be a tainted
(attacker controlled) value.

The `clang-analyzer-optin.taint.TaintedAlloc` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-taint-taintedalloc)
for more information.

```{title} clang-tidy - clang-analyzer-optin.taint.TaintedDiv
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#optin-taint-tainteddiv
```

# clang-analyzer-optin.taint.TaintedDiv

Check for divisions where the denominator is tainted (attacker controlled) and
might be 0.

The `clang-analyzer-optin.taint.TaintedDiv` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#optin-taint-tainteddiv)
for more information.

```{title} clang-tidy - clang-analyzer-osx.API
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-api
```

# clang-analyzer-osx.API

Check for proper uses of various Apple APIs.

The `clang-analyzer-osx.API` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-api)
for more information.

```{title} clang-tidy - clang-analyzer-osx.MIG
```

# clang-analyzer-osx.MIG

Find violations of the Mach Interface Generator calling convention.

The clang-analyzer-osx.MIG check is an alias of
Clang Static Analyzer osx.MIG.

```{title} clang-tidy - clang-analyzer-osx.NumberObjectConversion
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-numberobjectconversion
```

# clang-analyzer-osx.NumberObjectConversion

Check for erroneous conversions of objects representing numbers into numbers.

The `clang-analyzer-osx.NumberObjectConversion` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-numberobjectconversion)
for more information.

```{title} clang-tidy - clang-analyzer-osx.OSObjectRetainCount
```

# clang-analyzer-osx.OSObjectRetainCount

Check for leaks and improper reference count management for OSObject.

The clang-analyzer-osx.OSObjectRetainCount check is an alias of
Clang Static Analyzer osx.OSObjectRetainCount.

```{title} clang-tidy - clang-analyzer-osx.ObjCProperty
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-objcproperty
```

# clang-analyzer-osx.ObjCProperty

Check for proper uses of Objective-C properties.

The `clang-analyzer-osx.ObjCProperty` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-objcproperty)
for more information.

```{title} clang-tidy - clang-analyzer-osx.SecKeychainAPI
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-seckeychainapi
```

# clang-analyzer-osx.SecKeychainAPI

Check for proper uses of Secure Keychain APIs.

The `clang-analyzer-osx.SecKeychainAPI` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-seckeychainapi)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.AtSync
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-atsync
```

# clang-analyzer-osx.cocoa.AtSync

Check for nil pointers used as mutexes for @synchronized.

The `clang-analyzer-osx.cocoa.AtSync` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-atsync)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.AutoreleaseWrite
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-autoreleasewrite
```

# clang-analyzer-osx.cocoa.AutoreleaseWrite

Warn about potentially crashing writes to autoreleasing objects from different
autoreleasing pools in Objective-C.

The `clang-analyzer-osx.cocoa.AutoreleaseWrite` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-autoreleasewrite)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.ClassRelease
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-classrelease
```

# clang-analyzer-osx.cocoa.ClassRelease

Check for sending 'retain', 'release', or 'autorelease' directly to a Class.

The `clang-analyzer-osx.cocoa.ClassRelease` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-classrelease)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.Dealloc
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-dealloc
```

# clang-analyzer-osx.cocoa.Dealloc

Warn about Objective-C classes that lack a correct implementation of -dealloc.

The `clang-analyzer-osx.cocoa.Dealloc` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-dealloc)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.IncompatibleMethodTypes
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-incompatiblemethodtypes
```

# clang-analyzer-osx.cocoa.IncompatibleMethodTypes

Warn about Objective-C method signatures with type incompatibilities.

The `clang-analyzer-osx.cocoa.IncompatibleMethodTypes` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-incompatiblemethodtypes)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.Loops
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-loops
```

# clang-analyzer-osx.cocoa.Loops

Improved modeling of loops using Cocoa collection types.

The `clang-analyzer-osx.cocoa.Loops` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-loops)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.MissingSuperCall
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-missingsupercall
```

# clang-analyzer-osx.cocoa.MissingSuperCall

Warn about Objective-C methods that lack a necessary call to super.

The `clang-analyzer-osx.cocoa.MissingSuperCall` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-missingsupercall)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.NSAutoreleasePool
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-nsautoreleasepool
```

# clang-analyzer-osx.cocoa.NSAutoreleasePool

Warn for suboptimal uses of NSAutoreleasePool in Objective-C GC mode.

The `clang-analyzer-osx.cocoa.NSAutoreleasePool` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-nsautoreleasepool)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.NSError
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-nserror
```

# clang-analyzer-osx.cocoa.NSError

Check usage of NSError** parameters.

The `clang-analyzer-osx.cocoa.NSError` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-nserror)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.NilArg
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-nilarg
```

# clang-analyzer-osx.cocoa.NilArg

Check for prohibited nil arguments to ObjC method calls.

The `clang-analyzer-osx.cocoa.NilArg` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-nilarg)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.NonNilReturnValue
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-nonnilreturnvalue
```

# clang-analyzer-osx.cocoa.NonNilReturnValue

Model the APIs that are guaranteed to return a non-nil value.

The `clang-analyzer-osx.cocoa.NonNilReturnValue` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-nonnilreturnvalue)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.ObjCGenerics
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-objcgenerics
```

# clang-analyzer-osx.cocoa.ObjCGenerics

Check for type errors when using Objective-C generics.

The `clang-analyzer-osx.cocoa.ObjCGenerics` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-objcgenerics)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.RetainCount
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-retaincount
```

# clang-analyzer-osx.cocoa.RetainCount

Check for leaks and improper reference count management.

The `clang-analyzer-osx.cocoa.RetainCount` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-retaincount)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.RunLoopAutoreleaseLeak
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-runloopautoreleaseleak
```

# clang-analyzer-osx.cocoa.RunLoopAutoreleaseLeak

Check for leaked memory in autorelease pools that will never be drained.

The `clang-analyzer-osx.cocoa.RunLoopAutoreleaseLeak` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-runloopautoreleaseleak)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.SelfInit
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-selfinit
```

# clang-analyzer-osx.cocoa.SelfInit

Check that 'self' is properly initialized inside an initializer method.

The `clang-analyzer-osx.cocoa.SelfInit` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-selfinit)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.SuperDealloc
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-superdealloc
```

# clang-analyzer-osx.cocoa.SuperDealloc

Warn about improper use of '[super dealloc]' in Objective-C.

The `clang-analyzer-osx.cocoa.SuperDealloc` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-superdealloc)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.UnusedIvars
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-unusedivars
```

# clang-analyzer-osx.cocoa.UnusedIvars

Warn about private ivars that are never used.

The `clang-analyzer-osx.cocoa.UnusedIvars` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-unusedivars)
for more information.

```{title} clang-tidy - clang-analyzer-osx.cocoa.VariadicMethodTypes
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-variadicmethodtypes
```

# clang-analyzer-osx.cocoa.VariadicMethodTypes

Check for passing non-Objective-C types to variadic collection initialization
methods that expect only Objective-C types.

The `clang-analyzer-osx.cocoa.VariadicMethodTypes` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-cocoa-variadicmethodtypes)
for more information.

```{title} clang-tidy - clang-analyzer-osx.coreFoundation.CFError
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-cferror
```

# clang-analyzer-osx.coreFoundation.CFError

Check usage of CFErrorRef* parameters.

The `clang-analyzer-osx.coreFoundation.CFError` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-cferror)
for more information.

```{title} clang-tidy - clang-analyzer-osx.coreFoundation.CFNumber
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-cfnumber
```

# clang-analyzer-osx.coreFoundation.CFNumber

Check for proper uses of CFNumber APIs.

The `clang-analyzer-osx.coreFoundation.CFNumber` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-cfnumber)
for more information.

```{title} clang-tidy - clang-analyzer-osx.coreFoundation.CFRetainRelease
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-cfretainrelease
```

# clang-analyzer-osx.coreFoundation.CFRetainRelease

Check for null arguments to CFRetain/CFRelease/CFMakeCollectable.

The `clang-analyzer-osx.coreFoundation.CFRetainRelease` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-cfretainrelease)
for more information.

```{title} clang-tidy - clang-analyzer-osx.coreFoundation.containers.OutOfBounds
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-containers-outofbounds
```

# clang-analyzer-osx.coreFoundation.containers.OutOfBounds

Checks for index out-of-bounds when using 'CFArray' API.

The `clang-analyzer-osx.coreFoundation.containers.OutOfBounds` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-containers-outofbounds)
for more information.

```{title} clang-tidy - clang-analyzer-osx.coreFoundation.containers.PointerSizedValues
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-containers-pointersizedvalues
```

# clang-analyzer-osx.coreFoundation.containers.PointerSizedValues

Warns if 'CFArray', 'CFDictionary', 'CFSet' are created with non-pointer-size
values.

The `clang-analyzer-osx.coreFoundation.containers.PointerSizedValues` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#osx-corefoundation-containers-pointersizedvalues)
for more information.

```{title} clang-tidy - clang-analyzer-security.ArrayBound
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-arraybound
```

# clang-analyzer-security.ArrayBound

Warn about out of bounds access to memory.

The `clang-analyzer-security.ArrayBound` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-arraybound)
for more information.

```{title} clang-tidy - clang-analyzer-security.FloatLoopCounter
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-floatloopcounter
```

# clang-analyzer-security.FloatLoopCounter

Warn on using a floating point value as a loop counter (CERT: FLP30-C,
FLP30-CPP).

The `clang-analyzer-security.FloatLoopCounter` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-floatloopcounter)
for more information.

```{title} clang-tidy - clang-analyzer-security.MmapWriteExec
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-mmapwriteexec
```

# clang-analyzer-security.MmapWriteExec

Warn on mmap() calls with both writable and executable access.

The `clang-analyzer-security.MmapWriteExec` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-mmapwriteexec)
for more information.

```{title} clang-tidy - clang-analyzer-security.PointerSub
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-pointersub
```

# clang-analyzer-security.PointerSub

Check for pointer subtractions on two pointers pointing to different memory
chunks.

The `clang-analyzer-security.PointerSub` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-pointersub)
for more information.

```{title} clang-tidy - clang-analyzer-security.PutenvStackArray
```

# clang-analyzer-security.PutenvStackArray

Finds calls to the function 'putenv' which pass a pointer to an automatic
(stack-allocated) array as the argument.

The clang-analyzer-security.PutenvStackArray check is an alias of
Clang Static Analyzer security.PutenvStackArray.

```{title} clang-tidy - clang-analyzer-security.SetgidSetuidOrder
```

# clang-analyzer-security.SetgidSetuidOrder

Warn on possible reversed order of 'setgid(getgid()))' and 'setuid(getuid())'
(CERT: POS36-C).

The clang-analyzer-security.SetgidSetuidOrder check is an alias of
Clang Static Analyzer security.SetgidSetuidOrder.

```{title} clang-tidy - clang-analyzer-security.VAList
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-valist
```

# clang-analyzer-security.VAList

Warn on misuse of va_list objects.

The `clang-analyzer-security.VAList` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-valist)
for more information.

```{title} clang-tidy - clang-analyzer-security.cert.env.InvalidPtr
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-cert-env-invalidptr
```

# clang-analyzer-security.cert.env.InvalidPtr

Finds usages of possibly invalidated pointers.

The `clang-analyzer-security.cert.env.InvalidPtr` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-cert-env-invalidptr)
for more information.

```{title} clang-tidy - clang-analyzer-security.insecureAPI.DeprecatedOrUnsafeBufferHandling
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-deprecatedorunsafebufferhandling
```

# clang-analyzer-security.insecureAPI.DeprecatedOrUnsafeBufferHandling

Warn on uses of unsecure or deprecated buffer manipulating functions.

The `clang-analyzer-security.insecureAPI.DeprecatedOrUnsafeBufferHandling` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-deprecatedorunsafebufferhandling)
for more information.

```{title} clang-tidy - clang-analyzer-security.insecureAPI.UncheckedReturn
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-uncheckedreturn
```

# clang-analyzer-security.insecureAPI.UncheckedReturn

Warn on uses of functions whose return values must be always checked.

The `clang-analyzer-security.insecureAPI.UncheckedReturn` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-uncheckedreturn)
for more information.

```{title} clang-tidy - clang-analyzer-security.insecureAPI.bcmp
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-bcmp
```

# clang-analyzer-security.insecureAPI.bcmp

Warn on uses of the 'bcmp' function.

The `clang-analyzer-security.insecureAPI.bcmp` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-bcmp)
for more information.

```{title} clang-tidy - clang-analyzer-security.insecureAPI.bcopy
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-bcopy
```

# clang-analyzer-security.insecureAPI.bcopy

Warn on uses of the 'bcopy' function.

The `clang-analyzer-security.insecureAPI.bcopy` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-bcopy)
for more information.

```{title} clang-tidy - clang-analyzer-security.insecureAPI.bzero
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-bzero
```

# clang-analyzer-security.insecureAPI.bzero

Warn on uses of the 'bzero' function.

The `clang-analyzer-security.insecureAPI.bzero` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-bzero)
for more information.

```{title} clang-tidy - clang-analyzer-security.insecureAPI.decodeValueOfObjCType
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-decodevalueofobjctype
```

# clang-analyzer-security.insecureAPI.decodeValueOfObjCType

Warn on uses of the '-decodeValueOfObjCType:at:' method.

The `clang-analyzer-security.insecureAPI.decodeValueOfObjCType` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-decodevalueofobjctype)
for more information.

```{title} clang-tidy - clang-analyzer-security.insecureAPI.getpw
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-getpw
```

# clang-analyzer-security.insecureAPI.getpw

Warn on uses of the 'getpw' function.

The `clang-analyzer-security.insecureAPI.getpw` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-getpw)
for more information.

```{title} clang-tidy - clang-analyzer-security.insecureAPI.gets
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-gets
```

# clang-analyzer-security.insecureAPI.gets

Warn on uses of the 'gets' function.

The `clang-analyzer-security.insecureAPI.gets` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-gets)
for more information.

```{title} clang-tidy - clang-analyzer-security.insecureAPI.mkstemp
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-mkstemp
```

# clang-analyzer-security.insecureAPI.mkstemp

Warn when 'mkstemp' is passed fewer than 6 X's in the format string.

The `clang-analyzer-security.insecureAPI.mkstemp` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-mkstemp)
for more information.

```{title} clang-tidy - clang-analyzer-security.insecureAPI.mktemp
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-mktemp
```

# clang-analyzer-security.insecureAPI.mktemp

Warn on uses of the 'mktemp' function.

The `clang-analyzer-security.insecureAPI.mktemp` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-mktemp)
for more information.

```{title} clang-tidy - clang-analyzer-security.insecureAPI.rand
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-rand
```

# clang-analyzer-security.insecureAPI.rand

Warn on uses of the 'rand', 'random', and related functions.

The `clang-analyzer-security.insecureAPI.rand` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-rand)
for more information.

```{title} clang-tidy - clang-analyzer-security.insecureAPI.strcpy
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-strcpy
```

# clang-analyzer-security.insecureAPI.strcpy

Warn on uses of the 'strcpy' and 'strcat' functions.

The `clang-analyzer-security.insecureAPI.strcpy` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-strcpy)
for more information.

```{title} clang-tidy - clang-analyzer-security.insecureAPI.vfork
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-vfork
```

# clang-analyzer-security.insecureAPI.vfork

Warn on uses of the 'vfork' function.

The `clang-analyzer-security.insecureAPI.vfork` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#security-insecureapi-vfork)
for more information.

```{title} clang-tidy - clang-analyzer-unix.API
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-api
```

# clang-analyzer-unix.API

Check calls to various UNIX/Posix functions.

The `clang-analyzer-unix.API` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-api)
for more information.

```{title} clang-tidy - clang-analyzer-unix.BlockInCriticalSection
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-blockincriticalsection
```

# clang-analyzer-unix.BlockInCriticalSection

Check for calls to blocking functions inside a critical section.

The `clang-analyzer-unix.BlockInCriticalSection` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-blockincriticalsection)
for more information.

```{title} clang-tidy - clang-analyzer-unix.Chroot
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-chroot
```

# clang-analyzer-unix.Chroot

Check improper use of chroot.

The `clang-analyzer-unix.Chroot` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-chroot)
for more information.

```{title} clang-tidy - clang-analyzer-unix.Errno
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-errno
```

# clang-analyzer-unix.Errno

Check for improper use of 'errno'.

The `clang-analyzer-unix.Errno` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-errno)
for more information.

```{title} clang-tidy - clang-analyzer-unix.Malloc
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-malloc
```

# clang-analyzer-unix.Malloc

Check for memory leaks, double free, and use-after-free problems. Traces memory
managed by malloc()/free().

The `clang-analyzer-unix.Malloc` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-malloc)
for more information.

```{title} clang-tidy - clang-analyzer-unix.MallocSizeof
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-mallocsizeof
```

# clang-analyzer-unix.MallocSizeof

Check for dubious malloc arguments involving sizeof.

The `clang-analyzer-unix.MallocSizeof` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-mallocsizeof)
for more information.

```{title} clang-tidy - clang-analyzer-unix.MismatchedDeallocator
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-mismatcheddeallocator
```

# clang-analyzer-unix.MismatchedDeallocator

Check for mismatched deallocators.

The `clang-analyzer-unix.MismatchedDeallocator` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-mismatcheddeallocator)
for more information.

```{title} clang-tidy - clang-analyzer-unix.StdCLibraryFunctions
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-stdclibraryfunctions
```

# clang-analyzer-unix.StdCLibraryFunctions

Check for invalid arguments of C standard library functions, and apply relations
between arguments and return value.

The `clang-analyzer-unix.StdCLibraryFunctions` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-stdclibraryfunctions)
for more information.

```{title} clang-tidy - clang-analyzer-unix.Stream
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-stream
```

# clang-analyzer-unix.Stream

Check stream handling functions.

The `clang-analyzer-unix.Stream` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-stream)
for more information.

```{title} clang-tidy - clang-analyzer-unix.Vfork
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-vfork
```

# clang-analyzer-unix.Vfork

Check for proper usage of vfork.

The `clang-analyzer-unix.Vfork` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-vfork)
for more information.

```{title} clang-tidy - clang-analyzer-unix.cstring.BadSizeArg
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-cstring-badsizearg
```

# clang-analyzer-unix.cstring.BadSizeArg

Check the size argument passed into C string functions for common erroneous
patterns.

The `clang-analyzer-unix.cstring.BadSizeArg` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-cstring-badsizearg)
for more information.

```{title} clang-tidy - clang-analyzer-unix.cstring.NotNullTerminated
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-cstring-notnullterminated
```

# clang-analyzer-unix.cstring.NotNullTerminated

Check for arguments passed to C string functions which are not null-terminated
strings.

The `clang-analyzer-unix.cstring.NotNullTerminated` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-cstring-notnullterminated)
for more information.

```{title} clang-tidy - clang-analyzer-unix.cstring.NullArg
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-cstring-nullarg
```

# clang-analyzer-unix.cstring.NullArg

Check for null pointers being passed as arguments to C string functions.

The `clang-analyzer-unix.cstring.NullArg` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-cstring-nullarg)
for more information.

```{title} clang-tidy - clang-analyzer-unix.cstring.UninitializedRead
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#unix-cstring-uninitializedread
```

# clang-analyzer-unix.cstring.UninitializedRead

Checks if the string manipulation function would read uninitialized bytes.

The `clang-analyzer-unix.cstring.UninitializedRead` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#unix-cstring-uninitializedread)
for more information.

```{title} clang-tidy - clang-analyzer-webkit.NoUncountedMemberChecker
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#webkit-nouncountedmemberchecker
```

# clang-analyzer-webkit.NoUncountedMemberChecker

Check for no uncounted member variables.

The `clang-analyzer-webkit.NoUncountedMemberChecker` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#webkit-nouncountedmemberchecker)
for more information.

```{title} clang-tidy - clang-analyzer-webkit.RefCntblBaseVirtualDtor
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#webkit-refcntblbasevirtualdtor
```

# clang-analyzer-webkit.RefCntblBaseVirtualDtor

Check for any ref-countable base class having virtual destructor.

The `clang-analyzer-webkit.RefCntblBaseVirtualDtor` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#webkit-refcntblbasevirtualdtor)
for more information.

```{title} clang-tidy - clang-analyzer-webkit.UncountedLambdaCapturesChecker
```

```{eval-rst}
.. meta::
 :http-equiv=refresh: 5;URL=https://clang.llvm.org/docs/analyzer/checkers.html#webkit-uncountedlambdacaptureschecker
```

# clang-analyzer-webkit.UncountedLambdaCapturesChecker

Check uncounted lambda captures.

The `clang-analyzer-webkit.UncountedLambdaCapturesChecker` check is an alias, please see
[Clang Static Analyzer Available Checkers](https://clang.llvm.org/docs/analyzer/checkers.html#webkit-uncountedlambdacaptureschecker)
for more information.

**Title:** clang-tidy - concurrency-mt-unsafe

# concurrency-mt-unsafe

Checks for some thread-unsafe functions against a black list of
known-to-be-unsafe functions. Usually they access static variables without
synchronization (e.g. gmtime(3)) or utilize signals in a racy way.
The set of functions to check is specified with the `FunctionSet` option.

Note that using some thread-unsafe functions may be still valid in
concurrent programming if only a single thread is used (e.g. setenv(3)),
however, some functions may track a state in global variables which
would be clobbered by subsequent (non-parallel, but concurrent) calls to
a related function. E.g. the following code suffers from unprotected
accesses to a global state:

```c++
// getnetent(3) maintains global state with DB connection, etc.
// If a concurrent green thread calls getnetent(3), the global state is corrupted.
netent = getnetent();
yield();
netent = getnetent();
```

Examples:

```c++
tm = gmtime(timep); // uses a global buffer

sleep(1); // implementation may use SIGALRM
```

## Options

**Option:** FunctionSet
Specifies which functions in libc should be considered thread-safe,
possible values are `posix`, `glibc`, or `any`.

`posix` means POSIX defined thread-unsafe functions. POSIX.1-2001
in "2.9.1 Thread-Safety" defines that all functions specified in the
standard are thread-safe except a predefined list of thread-unsafe
functions.

Glibc defines some of them as thread-safe (e.g. dirname(3)), but adds
non-POSIX thread-unsafe ones (e.g. getopt_long(3)). Glibc's list is
compiled from GNU web documentation with a search for MT-Safe tag:
https://www.gnu.org/software/libc/manual/html_node/POSIX-Safety-Concepts.html

If you want to identify thread-unsafe API for at least one libc or
unsure which libc will be used, use `any` (default).

```{title} clang-tidy - concurrency-thread-canceltype-asynchronous
```

# concurrency-thread-canceltype-asynchronous

Finds `pthread_setcanceltype` function calls where a thread's cancellation
type is set to asynchronous. Asynchronous cancellation type
(`PTHREAD_CANCEL_ASYNCHRONOUS`) is generally unsafe, use type
`PTHREAD_CANCEL_DEFERRED` instead which is the default. Even with deferred
cancellation, a cancellation point in an asynchronous signal handler may still
be acted upon and the effect is as if it was an asynchronous cancellation.

```cpp
pthread_setcanceltype(PTHREAD_CANCEL_ASYNCHRONOUS, &oldtype);
```

This check corresponds to the CERT C Coding Standard rule
[POS47-C. Do not use threads that can be canceled asynchronously](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/posix-pos/pos47-c/).

`cert-pos47-c` redirects here as an alias of this check.

```{title} clang-tidy - cppcoreguidelines-avoid-c-arrays
```

# cppcoreguidelines-avoid-c-arrays

The `cppcoreguidelines-avoid-c-arrays` check is an alias, please see
{doc}`modernize-avoid-c-arrays <../modernize/avoid-c-arrays>`
for more information.

**Title:** clang-tidy - cppcoreguidelines-avoid-capturing-lambda-coroutines

# cppcoreguidelines-avoid-capturing-lambda-coroutines

Flags C++20 coroutine lambdas with non-empty capture lists that may cause
use-after-free errors and suggests avoiding captures or ensuring the lambda
closure object has a guaranteed lifetime.

This check implements `CP.51
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#rcoro-capture>`_
from the C++ Core Guidelines.

Using coroutine lambdas with non-empty capture lists can be risky, as capturing
variables can lead to accessing freed memory after the first suspension point.
This issue can occur even with refcounted smart pointers and copyable types.
When a lambda expression creates a coroutine, it results in a closure object
with storage, which is often on the stack and will eventually go out of scope.
When the closure object goes out of scope, its captures also go out of scope.
While normal lambdas finish executing before this happens, coroutine lambdas
may resume from suspension after the closure object has been destructed,
resulting in use-after-free memory access for all captures.

Consider the following example:

```c++
int value = get_value();
std::shared_ptr<Foo> sharedFoo = get_foo();
{
 const auto lambda = [value, sharedFoo]() -> std::future<void>
 {
 co_await something();
 // "sharedFoo" and "value" have already been destroyed
 // the "shared" pointer didn't accomplish anything
 };
 lambda();
} // the lambda closure object has now gone out of scope
```

In this example, the lambda object is defined with two captures: value and
`sharedFoo`. When `lambda()` is called, the lambda object is created on the
stack, and the captures are copied into the closure object. When the coroutine
is suspended, the lambda object goes out of scope, and the closure object is
destroyed. When the coroutine is resumed, the captured variables may have been
destroyed, resulting in use-after-free bugs.

In conclusion, the use of coroutine lambdas with non-empty capture lists can
lead to use-after-free errors when resuming the coroutine after the closure
object has been destroyed. This check helps prevent such errors by flagging
C++20 coroutine lambdas with non-empty capture lists and suggesting avoiding
captures or ensuring the lambda closure object has a guaranteed lifetime.

Following these guidelines can help ensure the safe and reliable use of
coroutine lambdas in C++ code.

## Options

**Option:** AllowExplicitObjectParameters
When set to `true`, lambda coroutines that use C++23 "deducing this"
(explicit object parameter, e.g. ``this auto``) are not flagged by this
check, because the captures are moved into the coroutine frame, decoupling
their lifetime from the lambda object.

Default is `false`.

The example from above can be made safe and will pass this check with the
following change:

.. code-block:: c++

 int value = get_value();
 std::shared_ptr<Foo> sharedFoo = get_foo();
 {
 // Pass "this auto" as the first argument to the lambda
 const auto lambda = [value, sharedFoo](this auto) -> std::future<void>
 {
 co_await something();
 };
 lambda();
 } // the lambda closure object has now gone out of scope, but captures are
 // no longer coupled to its lifetime

**Title:** clang-tidy - cppcoreguidelines-avoid-const-or-ref-data-members

# cppcoreguidelines-avoid-const-or-ref-data-members

This check warns when structs or classes that are copyable or movable, and have
const-qualified or reference (lvalue or rvalue) data members. Having such
members is rarely useful, and makes the class only copy-constructible but not
copy-assignable.

Examples:

```c++
// Bad, const-qualified member
struct Const {
 const int x;
}

// Good:
class Foo {
 public:
 int get() const { return x; }
 private:
 int x;
};

// Bad, lvalue reference member
struct Ref {
 int& x;
};

// Good:
struct Foo {
 int* x;
 std::unique_ptr<int> x;
 std::shared_ptr<int> x;
 gsl::not_null<int*> x;
};

// Bad, rvalue reference member
struct RefRef {
 int&& x;
};
```

This check implements `C.12
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#rc-constref>`_
from the C++ Core Guidelines.

Further reading:
[Data members: Never const](https://quuxplusone.github.io/blog/2022/01/23/dont-const-all-the-things/#data-members-never-const).

**Title:** clang-tidy - cppcoreguidelines-avoid-do-while

# cppcoreguidelines-avoid-do-while

Warns when using `do-while` loops. They are less readable than plain
`while` loops, since the termination condition is at the end and the
condition is not checked prior to the first iteration.
This can lead to subtle bugs.

This check implements `ES.75
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#res-do>`_
from the C++ Core Guidelines.

Examples:

```c++
int x;
do {
 std::cin >> x;
 // ...
} while (x < 0);
```

## Options

**Option:** IgnoreMacros
Ignore the check when analyzing macros. This is useful for safely defining function-like macros:

.. code-block:: c++

 #define FOO_BAR(x) \
 do { \
 foo(x); \
 bar(x); \
 } while(0)

Defaults to `false`.

**Title:** clang-tidy - cppcoreguidelines-avoid-goto

# cppcoreguidelines-avoid-goto

The usage of `goto` for control flow is error prone and should be replaced
with looping constructs. Only forward jumps in nested loops are accepted.

This check implements `ES.76
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#es76-avoid-goto>`_
from the C++ Core Guidelines.

For more information on why to avoid programming
with `goto` you can read the famous paper [A Case against the GO TO Statement.](https://www.cs.utexas.edu/users/EWD/ewd02xx/EWD215.PDF).

The check diagnoses `goto` for backward jumps in every language mode. These
should be replaced with `C/C++` looping constructs.

```c++
// Bad, handwritten for loop.
int i = 0;
// Jump label for the loop
loop_start:
do_some_operation();

if (i < 100) {
 ++i;
 goto loop_start;
}

// Better
for(int i = 0; i < 100; ++i)
 do_some_operation();
```

Modern C++ needs `goto` only to jump out of nested loops.

```c++
for(int i = 0; i < 100; ++i) {
 for(int j = 0; j < 100; ++j) {
 if (i * j > 500)
 goto early_exit;
 }
}

early_exit:
some_operation();
```

All other uses of `goto` are diagnosed in `C++`.

## Options

**Option:** IgnoreMacros
If set to `true`, the check will not warn if a ``goto`` statement is
expanded from a macro. Default is `false`.

```{title} clang-tidy - cppcoreguidelines-avoid-magic-numbers
```

# cppcoreguidelines-avoid-magic-numbers

The `cppcoreguidelines-avoid-magic-numbers` check is an alias, please see
{doc}`readability-magic-numbers <../readability/magic-numbers>`
for more information.

**Title:** clang-tidy - cppcoreguidelines-avoid-non-const-global-variables

# cppcoreguidelines-avoid-non-const-global-variables

Finds non-const global variables as described in `I.2
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#i2-avoid-non-const-global-variables>`_
of C++ Core Guidelines.
As [R.6](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#rr-global)
of C++ Core Guidelines is a duplicate of rule `I.2
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#i2-avoid-non-const-global-variables>`_
it also covers that rule.

```c++
char a; // Warns!
const char b = 0;

namespace some_namespace
{
 char c; // Warns!
 const char d = 0;
}

char * c_ptr1 = &some_namespace::c; // Warns!
char *const c_const_ptr = &some_namespace::c; // Warns!
char & c_reference = some_namespace::c; // Warns!

class Foo // No Warnings inside Foo, only namespace scope is covered
{
public:
 char e = 0;
 const char f = 0;
protected:
 char g = 0;
private:
 char h = 0;
};
```

The variables `a`, `c`, `c_ptr1`, `c_const_ptr` and `c_reference`
will all generate warnings since they are either a non-const globally accessible
variable, a pointer or a reference providing global access to non-const data
or both.

## Options

**Option:** AllowInternalLinkage
When set to `true`, static non-const variables and variables in anonymous
namespaces will not generate a warning. The default value is `false`.
**Option:** AllowThreadLocal
When set to `true`, non-const global variables with thread-local storage
duration will not generate a warning. The default value is `false`.
**Option:** IgnoreMacros
When set to `true`, non-const global variables defined in macros will not
generate a warning. The default value is `false`.

```{title} clang-tidy - cppcoreguidelines-avoid-reference-coroutine-parameters
```

# cppcoreguidelines-avoid-reference-coroutine-parameters

Warns when a coroutine accepts reference parameters. After a coroutine suspend
point, references could be dangling and no longer valid. Instead, pass
parameters as values.

Examples:

```cpp
std::future<int> someCoroutine(int& val) {
 co_await ...;
 // When the coroutine is resumed, 'val' might no longer be valid.
 if (val) ...
}
```

This check implements [CP.53](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#rcoro-reference-parameters)
from the C++ Core Guidelines.

```{title} clang-tidy - cppcoreguidelines-c-copy-assignment-signature
```

# cppcoreguidelines-c-copy-assignment-signature

The `cppcoreguidelines-c-copy-assignment-signature` check is an alias,
please see {doc}`misc-unconventional-assign-operator <../misc/unconventional-assign-operator>` for more information.

```{title} clang-tidy - cppcoreguidelines-explicit-constructor
```

# cppcoreguidelines-explicit-constructor

This check is an alias for
{doc}`misc-explicit-constructor <../misc/explicit-constructor>`.

```{title} clang-tidy - cppcoreguidelines-explicit-virtual-functions
```

# cppcoreguidelines-explicit-virtual-functions

The `cppcoreguidelines-explicit-virtual-functions` check is an alias,
please see {doc}`modernize-use-override <../modernize/use-override>`
for more information.

**Title:** clang-tidy - cppcoreguidelines-init-variables

# cppcoreguidelines-init-variables

Checks whether there are local variables that are declared without an initial
value. These may lead to unexpected behavior if there is a code path that reads
the variable before assigning to it.

This rule is part of the `Type safety (Type.5)
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#pro-type-init>`_
profile and [ES.20](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#res-always)
from the C++ Core Guidelines.

Only integers, booleans, floats, doubles and pointers are checked. The fix
option initializes all detected values with the value of zero. An exception is
float and double types, which are initialized to NaN.

As an example a function that looks like this:

```c++
void function() {
 int x;
 char *txt;
 double d;

 // Rest of the function.
}
```

Would be rewritten to look like this:

```c++
#include <math.h>

void function() {
 int x = 0;
 char *txt = nullptr;
 double d = NAN;

 // Rest of the function.
}
```

It warns for the uninitialized enum case, but without a FixIt:

```c++
enum A {A1, A2, A3};
enum A_c : char { A_c1, A_c2, A_c3 };
enum class B { B1, B2, B3 };
enum class B_i : int { B_i1, B_i2, B_i3 };
void function() {
 A a; // Warning: variable 'a' is not initialized
 A_c a_c; // Warning: variable 'a_c' is not initialized
 B b; // Warning: variable 'b' is not initialized
 B_i b_i; // Warning: variable 'b_i' is not initialized
}
```

## Options

**Option:** IncludeStyle
A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.
**Option:** MathHeader
A string specifying the header to include to get the definition of `NAN`.
Default is `<math.h>`.

```{title} clang-tidy - cppcoreguidelines-interfaces-global-init
```

# cppcoreguidelines-interfaces-global-init

This check flags initializers of globals that access extern objects,
and therefore can lead to order-of-initialization problems.

This check implements [I.22](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#ri-global-init)
from the C++ Core Guidelines.

Note that currently this does not flag calls to non-constexpr functions, and
therefore globals could still be accessed from functions themselves.

```{title} clang-tidy - cppcoreguidelines-macro-to-enum
```

# cppcoreguidelines-macro-to-enum

The `cppcoreguidelines-macro-to-enum` check is an alias, please see
{doc}`modernize-macro-to-enum <../modernize/macro-to-enum>`
for more information.

**Title:** clang-tidy - cppcoreguidelines-macro-usage

# cppcoreguidelines-macro-usage

Finds macro usage that is considered problematic because better language
constructs exist for the task.

The relevant sections in the C++ Core Guidelines are
[ES.31](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#es31-dont-use-macros-for-constants-or-functions), and
[ES.32](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#es32-use-all_caps-for-all-macro-names).

Examples:

```c++
#define C 0
#define F1(x, y) ((a) > (b) ? (a) : (b))
#define F2(...) (__VA_ARGS__)
#define F3(x, y) x##y
#define COMMA ,
#define NORETURN [[noreturn]]
#define DEPRECATED attribute((deprecated))
#if LIB_EXPORTS
#define DLLEXPORTS __declspec(dllexport)
#else
#define DLLEXPORTS __declspec(dllimport)
#endif
```

results in the following warnings:

```
4 warnings generated.
test.cpp:1:9: warning: macro 'C' used to declare a constant; consider using a 'constexpr' constant [cppcoreguidelines-macro-usage]
#define C 0
 ^
test.cpp:2:9: warning: function-like macro 'F1' used; consider a 'constexpr' template function [cppcoreguidelines-macro-usage]
#define F1(x, y) ((a) > (b) ? (a) : (b))
 ^
test.cpp:3:9: warning: variadic macro 'F2' used; consider using a 'constexpr' variadic template function [cppcoreguidelines-macro-usage]
#define F2(...) (__VA_ARGS__)
 ^
```

## Options

**Option:** AllowedRegexp
A regular expression to filter allowed macros. For example
`DEBUG*|LIBTORRENT*|TORRENT*|UNI*` could be applied to filter `libtorrent`.
Default value is `^DEBUG_*`.
**Option:** CheckCapsOnly
Boolean flag to warn on all macros except those with CAPS_ONLY names.
This option is intended to ease introduction of this check into older
code bases. Default value is `false`.
**Option:** IgnoreCommandLineMacros
Boolean flag to toggle ignoring command-line-defined macros.
Default value is `true`.

**Title:** clang-tidy - cppcoreguidelines-misleading-capture-default-by-value

# cppcoreguidelines-misleading-capture-default-by-value

Warns when lambda specify a by-value capture default and capture `this`.

By-value capture defaults in member functions can be misleading about whether
data members are captured by value or reference. This occurs because specifying
the capture default `[=]` actually captures the `this` pointer by value,
not the data members themselves. As a result, data members are still indirectly
accessed via the captured `this` pointer, which essentially means they are
being accessed by reference. Therefore, even when using `[=]`, data members
are effectively captured by reference, which might not align with the user's
expectations.

Examples:

```c++
struct AClass {
 int member;
 void misleadingLogic() {
 int local = 0;
 member = 0;
 auto f = [=]() mutable {
 local += 1;
 member += 1;
 };
 f();
 // Here, local is 0 but member is 1
 }

 void clearLogic() {
 int local = 0;
 member = 0;
 auto f = [this, local]() mutable {
 local += 1;
 member += 1;
 };
 f();
 // Here, local is 0 but member is 1
 }
};
```

This check implements `F.54
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#f54-when-writing-a-lambda-that-captures-this-or-any-class-data-member-dont-use--default-capture>`_
from the C++ Core Guidelines.

**Title:** clang-tidy - cppcoreguidelines-missing-std-forward

# cppcoreguidelines-missing-std-forward

Warns when a forwarding reference parameter is not forwarded inside the
function body.

Example:

```c++
template <class T>
void wrapper(T&& t) {
 impl(std::forward<T>(t), 1, 2); // Correct
}

template <class T>
void wrapper2(T&& t) {
 impl(t, 1, 2); // Oops - should use std::forward<T>(t)
}

template <class T>
void wrapper3(T&& t) {
 impl(std::move(t), 1, 2); // Also buggy - should use std::forward<T>(t)
}

template <class F>
void wrapper_function(F&& f) {
 std::forward<F>(f)(1, 2); // Correct
}

template <class F>
void wrapper_function2(F&& f) {
 f(1, 2); // Incorrect - may not invoke the desired qualified function operator
}
```

## Options

**Option:** ForwardFunction
Specify the function used for forwarding. Default is `::std::forward`.
This check implements `F.19
<http://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#rf-forward>`_
from the C++ Core Guidelines.

## Limitations

Explicit object parameters (`this Self&&`) with a type constraint are not
checked to avoid false positives. Such parameters are rarely intended to be
perfectly forwarded.

```{title} clang-tidy - cppcoreguidelines-narrowing-conversions
```

# cppcoreguidelines-narrowing-conversions

This check implements part of [ES.46](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#es46-avoid-lossy-narrowing-truncating-arithmetic-conversions)
from the C++ Core Guidelines.

The `cppcoreguidelines-narrowing-conversions` check is an alias, please see
{doc}`bugprone-narrowing-conversions <../bugprone/narrowing-conversions>`
for more information.

**Title:** clang-tidy - cppcoreguidelines-no-malloc

# cppcoreguidelines-no-malloc

This check handles C-Style memory management using `malloc()`, `realloc()`,
`calloc()` and `free()`. It warns about its use and tries to suggest the
use of an appropriate RAII object.
Furthermore, it can be configured to check against a user-specified list of
functions that are used for memory management (e.g. `posix_memalign()`).

This check implements `R.10
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#rr-mallocfree>`_
from the C++ Core Guidelines.

There is no attempt made to provide fix-it hints, since manual resource
management isn't easily transformed automatically into RAII.

```c++
// Warns each of the following lines.
// Containers like std::vector or std::string should be used.
char* some_string = (char*) malloc(sizeof(char) * 20);
char* some_string = (char*) realloc(sizeof(char) * 30);
free(some_string);

int* int_array = (int*) calloc(30, sizeof(int));

// Rather use a smartpointer or stack variable.
struct some_struct* s = (struct some_struct*) malloc(sizeof(struct some_struct));
```

## Options

**Option:** Allocations
Semicolon-separated list of fully qualified names of memory allocation functions.
Defaults to `::malloc;::calloc`.
**Option:** Deallocations
Semicolon-separated list of fully qualified names of memory allocation functions.
Defaults to `::free`.
**Option:** Reallocations
Semicolon-separated list of fully qualified names of memory allocation functions.
Defaults to `::realloc`.

**Title:** clang-tidy - cppcoreguidelines-no-suspend-with-lock

# cppcoreguidelines-no-suspend-with-lock

Flags coroutines that suspend while a lock guard is in scope at the
suspension point.

When a coroutine suspends, any mutexes held by the coroutine will remain
locked until the coroutine resumes and eventually destructs the lock guard.
This can lead to long periods with a mutex held and runs the risk of deadlock.

Instead, locks should be released before suspending a coroutine.

This check only checks suspending coroutines while a lock_guard is in scope;
it does not consider manual locking or unlocking of mutexes, e.g., through
calls to `std::mutex::lock()`.

Examples:

```c++
future bad_coro() {
 std::lock_guard lock{mtx};
 ++some_counter;
 co_await something(); // Suspending while holding a mutex
}

future good_coro() {
 {
 std::lock_guard lock{mtx};
 ++some_counter;
 }
 // Destroy the lock_guard to release the mutex before suspending the coroutine
 co_await something(); // Suspending while holding a mutex
}
```

This check implements `CP.52
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#rcoro-locks>`_
from the C++ Core Guidelines.

```{title} clang-tidy - cppcoreguidelines-noexcept-destructor
```

# cppcoreguidelines-noexcept-destructor

This check implements [C.37](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#c37-make-destructors-noexcept)
from the C++ Core Guidelines.

The `cppcoreguidelines-noexcept-destructor` check is an alias, please see
{doc}`performance-noexcept-destructor <../performance/noexcept-destructor>`
for more information.

```{title} clang-tidy - cppcoreguidelines-noexcept-move-operations
```

# cppcoreguidelines-noexcept-move-operations

This check implements [C.66](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#c66-make-move-operations-noexcept)
from the C++ Core Guidelines.

The `cppcoreguidelines-noexcept-move-operations` check is an alias, please see
{doc}`performance-noexcept-move-constructor <../performance/noexcept-move-constructor>`
for more information.

```{title} clang-tidy - cppcoreguidelines-noexcept-swap
```

# cppcoreguidelines-noexcept-swap

This check implements [C.83](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#c83-for-value-like-types-consider-providing-a-noexcept-swap-function),
[C.84](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#c84-a-swap-function-must-not-fail)
and [C.85](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#c85-make-swap-noexcept)
from the C++ Core Guidelines.

The `cppcoreguidelines-noexcept-swap` check is an alias, please see
{doc}`performance-noexcept-swap <../performance/noexcept-swap>`
for more information.

```{title} clang-tidy - cppcoreguidelines-non-private-member-variables-in-classes
```

# cppcoreguidelines-non-private-member-variables-in-classes

The `cppcoreguidelines-non-private-member-variables-in-classes` check is
an alias, please see {doc}`misc-non-private-member-variables-in-classes <../misc/non-private-member-variables-in-classes>`
for more information.

**Title:** clang-tidy - cppcoreguidelines-owning-memory

# cppcoreguidelines-owning-memory

This check implements the type-based semantics of `gsl::owner<T*>`, which
allows static analysis on code, that uses raw pointers to handle resources
like dynamic memory, but won't introduce RAII concepts.

This check implements `I.11
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#i11-never-transfer-ownership-by-a-raw-pointer-t-or-reference-t>`_,
`C.33
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#c33-if-a-class-has-an-owning-pointer-member-define-a-destructor>`_,
`R.3
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#r3-a-raw-pointer-a-t-is-non-owning>`_
and `GSL.Views
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#ss-views>`_
from the C++ Core Guidelines.
The definition of a `gsl::owner<T*>` is straight forward

```c++
namespace gsl { template <typename T> owner = T; }
```

It is therefore simple to introduce the owner even without using an implementation of
the [Guideline Support Library](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#s-gsl).

All checks are purely type based and not (yet) flow sensitive.

The following examples will demonstrate the correct and incorrect
initializations of owners, assignment is handled the same way.
Note that both `new` and `malloc()`-like resource functions are
considered to produce resources.

```c++
// Creating an owner with factory functions is checked.
gsl::owner<int*> function_that_returns_owner() { return gsl::owner<int*>(new int(42)); }

// Dynamic memory must be assigned to an owner
int* Something = new int(42); // BAD, will be caught
gsl::owner<int*> Owner = new int(42); // Good
gsl::owner<int*> Owner = new int[42]; // Good as well

// Returned owner must be assigned to an owner
int* Something = function_that_returns_owner(); // Bad, factory function
gsl::owner<int*> Owner = function_that_returns_owner(); // Good, result lands in owner

// Something not a resource or owner should not be assigned to owners
int Stack = 42;
gsl::owner<int*> Owned = &Stack; // Bad, not a resource assigned
```

In the case of dynamic memory as resource, only `gsl::owner<T*>` variables are allowed
to be deleted.

```c++
// Example Bad, non-owner as resource handle, will be caught.
int* NonOwner = new int(42); // First warning here, since new must land in an owner
delete NonOwner; // Second warning here, since only owners are allowed to be deleted

// Example Good, Ownership correctly stated
gsl::owner<int*> Owner = new int(42); // Good
delete Owner; // Good as well, statically enforced, that only owners get deleted
```

The check will furthermore ensure, that functions, that expect a
`gsl::owner<T*>` as argument get called with either a `gsl::owner<T*>` or
a newly created resource.

```c++
void expects_owner(gsl::owner<int*> o) { delete o; }

// Bad Code
int NonOwner = 42;
expects_owner(&NonOwner); // Bad, will get caught

// Good Code
gsl::owner<int*> Owner = new int(42);
expects_owner(Owner); // Good
expects_owner(new int(42)); // Good as well, recognized created resource

// Port legacy code for better resource-safety
gsl::owner<FILE*> File = fopen("my_file.txt", "rw+");
FILE* BadFile = fopen("another_file.txt", "w"); // Bad, warned

// ... use the file

fclose(File); // Ok, File is annotated as 'owner<>'
fclose(BadFile); // BadFile is not an 'owner<>', will be warned
```

## Options

**Option:** LegacyResourceProducers
Semicolon-separated list of fully qualified names of legacy functions that create
resources but cannot introduce ``gsl::owner<>``.
Defaults to `::malloc;::aligned_alloc;::realloc;::calloc;::fopen;::freopen;::tmpfile`.
**Option:** LegacyResourceConsumers
Semicolon-separated list of fully qualified names of legacy functions expecting
resource owners as pointer arguments but cannot introduce ``gsl::owner<>``.
Defaults to `::free;::realloc;::freopen;::fclose`.

## Limitations

Using `gsl::owner<T*>` in a typedef or alias is not handled correctly.

```c++
using heap_int = gsl::owner<int*>;
heap_int allocated = new int(42); // False positive!
```

The `gsl::owner<T*>` is declared as a templated type alias.
In template functions and classes, like in the example below, the information
of the type aliases gets lost. Therefore using `gsl::owner<T*>` in a heavy templated
code base might lead to false positives.

Known code constructs that do not get diagnosed correctly are:

- `std::exchange`
- `std::vector<gsl::owner<T*>>`

```c++
// This template function works as expected. Type information doesn't get lost.
template <typename T>
void delete_owner(gsl::owner<T*> owned_object) {
 delete owned_object; // Everything alright
}

gsl::owner<int*> function_that_returns_owner() { return gsl::owner<int*>(new int(42)); }

// Type deduction does not work for auto variables.
// This is caught by the check and will be noted accordingly.
auto OwnedObject = function_that_returns_owner(); // Type of OwnedObject will be int*

// Problematic function template that looses the typeinformation on owner
template <typename T>
void bad_template_function(T some_object) {
 // This line will trigger the warning, that a non-owner is assigned to an owner
 gsl::owner<T*> new_owner = some_object;
}

// Calling the function with an owner still yields a false positive.
bad_template_function(gsl::owner<int*>(new int(42)));

// The same issue occurs with templated classes like the following.
template <typename T>
class OwnedValue {
public:
 const T getValue() const { return _val; }
private:
 T _val;
};

// Code, that yields a false positive.
OwnedValue<gsl::owner<int*>> Owner(new int(42)); // Type deduction yield T -> int *
// False positive, getValue returns int* and not gsl::owner<int*>
gsl::owner<int*> OwnedInt = Owner.getValue();
```

Another limitation of the current implementation is only the type based
checking. Suppose you have code like the following:

```c++
// Two owners with assigned resources
gsl::owner<int*> Owner1 = new int(42);
gsl::owner<int*> Owner2 = new int(42);

Owner2 = Owner1; // Conceptual Leak of initial resource of Owner2!
Owner1 = nullptr;
```

The semantic of a `gsl::owner<T*>` is mostly like a `std::unique_ptr<T>`,
therefore assignment of two `gsl::owner<T*>` is considered a move, which
requires that the resource `Owner2` must have been released before the
assignment. This kind of condition could be caught in later improvements of
this check with flowsensitive analysis. Currently, the `Clang Static Analyzer`
catches this bug for dynamic memory, but not for general types of resources.

**Title:** clang-tidy - cppcoreguidelines-prefer-member-initializer

# cppcoreguidelines-prefer-member-initializer

Finds member initializations in the constructor body which can be converted
into member initializers of the constructor instead. This not only improves
the readability of the code but also positively affects its performance.
Class-member assignments inside a control statement or following the first
control statement are ignored.

This check implements `C.49
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#c49-prefer-initialization-to-assignment-in-constructors>`_
from the C++ Core Guidelines.

Please note, that this check does not enforce rule `C.48
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#c48-prefer-in-class-initializers-to-member-initializers-in-constructors-for-constant-initializers>`_
from the C++ Core Guidelines. For that purpose
see check :doc:`modernize-use-default-member-init
<../modernize/use-default-member-init>`.

## Example 1

```c++
class C {
 int n;
 int m;
public:
 C() {
 n = 1; // Literal in default constructor
 if (dice())
 return;
 m = 1;
 }
};
```

Here `n` can be initialized in the constructor initializer list, unlike
`m`, as `m`'s initialization follows a control statement (`if`):

```c++
class C {
 int n;
 int m;
public:
 C(): n(1) {
 if (dice())
 return;
 m = 1;
 }
};
```

## Example 2

```c++
class C {
 int n;
 int m;
public:
 C(int nn, int mm) {
 n = nn; // Neither default constructor nor literal
 if (dice())
 return;
 m = mm;
 }
};
```

Here `n` can be initialized in the constructor initializer list, unlike
`m`, as `m`'s initialization follows a control statement (`if`):

```c++
C(int nn, int mm) : n(nn) {
 if (dice())
 return;
 m = mm;
}
```

```{title} clang-tidy - cppcoreguidelines-pro-bounds-array-to-pointer-decay
```

# cppcoreguidelines-pro-bounds-array-to-pointer-decay

This check flags all array to pointer decays.

Pointers should not be used as arrays. `span<T>` is a bounds-checked, safe
alternative to using pointers to access arrays.

This rule is part of the [Bounds safety (Bounds 3)](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#pro-bounds-decay)
profile from the C++ Core Guidelines.

**Title:** clang-tidy - cppcoreguidelines-pro-bounds-avoid-unchecked-container-access

# cppcoreguidelines-pro-bounds-avoid-unchecked-container-access

Finds calls to `operator[]` in STL containers and suggests replacing them
with safe alternatives.
Safe alternatives include STL `at` or GSL `at` functions, `begin()` or
`end()` functions, `range-for` loops, `std::span`, or an appropriate
function from `<algorithms>`.

For example, both

```c++
std::vector<int> a;
int b = a[4];
```

and

```c++
std::unique_ptr<vector> a;
int b = a[0];
```

will generate a warning.

STL containers for which `operator[]` is well-defined for all inputs are
excluded from this check (e.g.: `std::map::operator[]`).

This check enforces part of the `SL.con.3
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#slcon3-avoid-bounds-errors>`_
guideline and is part of the `Bounds Safety (Bounds 4)
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#pro-bounds-arrayindex>`_
profile from the C++ Core Guidelines.

## Options

**Option:** ExcludeClasses
Semicolon-separated list of regular expressions matching class names that
overwrites the default exclusion list. The default is:
`::std::map;::std::unordered_map;::std::flat_map`.
**Option:** FixMode
Determines what fixes are suggested. Either `none`, `at` (use
``a.at(index)`` if a fitting function exists) or `function` (use a
function ``f(a, index)``). The default is `none`.
**Option:** FixFunction
The function to use in the `function` mode. For C++23 and beyond, the
passed function must support the empty subscript operator, i.e., the case
where ``a[]`` becomes ``f(a)``. :option:`FixFunctionEmptyArgs` can be
used to override the suggested function in that case. The default is `gsl::at`.
**Option:** FixFunctionEmptyArgs
The function to use in the `function` mode for the empty subscript operator
case in C++23 and beyond only. If no fixes should be made for empty
subscript operators, pass an empty string. In that case, only the warnings
will be printed. The default is the value of :option:`FixFunction`.

```{title} clang-tidy - cppcoreguidelines-pro-bounds-constant-array-index
```

# cppcoreguidelines-pro-bounds-constant-array-index

This check flags all array subscript expressions on static arrays and
`std::arrays` that either do not have a constant integer expression index or
are out of bounds (for `std::array`). For out-of-bounds checking of static
arrays, see the `-Warray-bounds` Clang diagnostic.

This rule is part of the [Bounds safety (Bounds 2)](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#pro-bounds-arrayindex)
profile from the C++ Core Guidelines.

Optionally, this check can generate fixes using `gsl::at` for indexing.

## Options

```{option} GslHeader

The check can generate fixes after this option has been set to the name of
the include file that contains `gsl::at()`, e.g. `"gsl/gsl.h"`.
Default is an empty string.
```

```{option} IncludeStyle

A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.
```

```{title} clang-tidy - cppcoreguidelines-pro-bounds-pointer-arithmetic
```

# cppcoreguidelines-pro-bounds-pointer-arithmetic

This check flags all usage of pointer arithmetic, because it could lead to an
invalid pointer. Subtraction of two pointers is not flagged by this check.

Pointers should only refer to single objects, and pointer arithmetic is fragile
and easy to get wrong. `span<T>` is a bounds-checked, safe type for accessing
arrays of data.

This rule is part of the [Bounds safety (Bounds 1)](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#pro-bounds-arithmetic)
profile from the C++ Core Guidelines.

## Options

```{option} AllowIncrementDecrementOperators

When enabled, the check will allow using the prefix/postfix increment or
decrement operators on pointers. Default is `false`.
```

**Title:** clang-tidy - cppcoreguidelines-pro-type-const-cast

# cppcoreguidelines-pro-type-const-cast

Imposes limitations on the use of `const_cast` within C++ code. It depends on
the `StrictMode` option setting to determine whether it should flag all
instances of `const_cast` or only those that remove either `const` or
`volatile` qualifier.

Modifying a variable that has been declared as `const` in C++ is generally
considered undefined behavior, and this remains true even when using
`const_cast`. In C++, the `const` qualifier indicates that a variable is
intended to be read-only, and the compiler enforces this by disallowing any
attempts to change the value of that variable.

Removing the `volatile` qualifier in C++ can have serious consequences. This
qualifier indicates that a variable's value can change unpredictably, and
removing it may lead to undefined behavior, optimization problems, and
debugging challenges. It's essential to retain the `volatile` qualifier in
situations where the variable's volatility is a crucial aspect of program
correctness and reliability.

This rule is part of the `Type safety (Type 3)
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#pro-type-constcast>`_
profile and `ES.50: Don’t cast away const
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#es50-dont-cast-away-const>`_
rule from the C++ Core Guidelines.

## Options

**Option:** StrictMode
When this setting is set to `true`, it means that any usage of ``const_cast``
is not allowed. On the other hand, when it's set to `false`, it permits
casting to ``const`` or ``volatile`` types. Default value is `false`.

```{title} clang-tidy - cppcoreguidelines-pro-type-cstyle-cast
```

# cppcoreguidelines-pro-type-cstyle-cast

This check flags all use of C-style casts that perform a `static_cast`
downcast, `const_cast`, or `reinterpret_cast`.

Use of these casts can violate type safety and cause the program to access a
variable that is actually of type X to be accessed as if it were of an
unrelated type Z. Note that a C-style `(T)expression` cast means to perform
the first of the following that is possible: a `const_cast`, a
`static_cast`, a `static_cast` followed by a `const_cast`, a
`reinterpret_cast`, or a `reinterpret_cast` followed by a `const_cast`.
This rule bans `(T)expression` only when used to perform an unsafe cast.

This rule is part of the [Type safety (Type.4)](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#pro-type-cstylecast)
profile from the C++ Core Guidelines.

**Title:** clang-tidy - cppcoreguidelines-pro-type-member-init

# cppcoreguidelines-pro-type-member-init

The check flags user-provided constructor definitions that do not
initialize all fields that would be left in an undefined state by
default construction, e.g. builtins, pointers and record types without
user-provided default constructors containing at least one such
type. If these fields aren't initialized, the constructor will leave
some of the memory in an undefined state.

For C++11 it suggests fixes to add in-class field initializers. For
older versions it inserts the field initializers into the constructor
initializer list. It will also initialize any direct base classes that
need to be zeroed in the constructor initializer list.

The check takes assignment of fields in the constructor body into
account but generates false positives for fields initialized in
methods invoked in the constructor body.

The check also flags variables with automatic storage duration that have record
types without a user-provided constructor and are not initialized. The
suggested fix is to zero initialize the variable via `{}` for C++11 and
beyond or `= {}` for older language versions.

## Options

**Option:** IgnoreArrays
If set to `true`, the check will not warn about array members (including
C arrays and ``std::array``) that are not zero-initialized during
construction. For performance critical code, it may be important to not
initialize fixed-size array members. Default is `false`.
**Option:** UseAssignment
If set to `true`, the check will provide fix-its with literal initializers
\( ``int i = 0;`` \) instead of curly braces \( ``int i{};`` \).
Default is `false`.
This rule is part of the `Type safety (Type.6)
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#pro-type-memberinit>`_
profile from the C++ Core Guidelines.

```{title} clang-tidy - cppcoreguidelines-pro-type-reinterpret-cast
```

# cppcoreguidelines-pro-type-reinterpret-cast

This check flags all uses of `reinterpret_cast` in C++ code.

Use of these casts can violate type safety and cause the program to access a
variable that is actually of type `X` to be accessed as if it were of an
unrelated type `Z`.

This rule is part of the [Type safety (Type.1.1)](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#pro-type-reinterpretcast)
profile from the C++ Core Guidelines.

```{title} clang-tidy - cppcoreguidelines-pro-type-static-cast-downcast
```

# cppcoreguidelines-pro-type-static-cast-downcast

This check flags all usages of `static_cast`, where a base class is casted to
a derived class. In those cases, a fix-it is provided to convert the cast to a
`dynamic_cast`.

Use of these casts can violate type safety and cause the program to access a
variable that is actually of type `X` to be accessed as if it were of an
unrelated type `Z`.

This rule is part of the [Type safety (Type.2)](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#pro-type-downcast)
profile from the C++ Core Guidelines.

## Options

```{option} StrictMode

When set to `false`, no warnings are emitted for casts on non-polymorphic
types. Default is `true`.
```

```{title} clang-tidy - cppcoreguidelines-pro-type-union-access
```

# cppcoreguidelines-pro-type-union-access

This check flags all access to members of unions. Passing unions as a whole is
not flagged.

Reading from a union member assumes that member was the last one written, and
writing to a union member assumes another member with a nontrivial destructor
had its destructor called. This is fragile because it cannot generally be
enforced to be safe in the language and so relies on programmer discipline to
get it right.

This rule is part of the [Type safety (Type.7)](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#pro-type-unions)
profile from the C++ Core Guidelines.

```{title} clang-tidy - cppcoreguidelines-pro-type-vararg
```

# cppcoreguidelines-pro-type-vararg

This check flags all calls to c-style vararg functions and all use of
`va_arg`.

To allow for SFINAE use of vararg functions, a call is not flagged if a literal
0 is passed as the only vararg argument or function is used in unevaluated
context.

Passing to varargs assumes the correct type will be read. This is fragile
because it cannot generally be enforced to be safe in the language and so
relies on programmer discipline to get it right.

This rule is part of the [Type safety (Type.8)](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#pro-type-varargs)
profile from the C++ Core Guidelines.

**Title:** clang-tidy - cppcoreguidelines-rvalue-reference-param-not-moved

# cppcoreguidelines-rvalue-reference-param-not-moved

Warns when an rvalue reference function parameter is never moved within
the function body.

Rvalue reference parameters indicate a parameter that should be moved with
`std::move` from within the function body. Any such parameter that is
never moved is confusing and potentially indicative of a buggy program.

Example:

```c++
void logic(std::string&& Input) {
 std::string Copy(Input); // Oops - forgot to std::move
}
```

Note that parameters that are unused and marked as such will not be diagnosed.

Example:

```c++
void conditional_use([[maybe_unused]] std::string&& Input) {
 // No diagnostic here since Input is unused and marked as such
}
```

## Options

**Option:** AllowPartialMove
 If set to `true`, the check accepts ``std::move`` calls containing any
 subexpression containing the parameter. CppCoreGuideline F.18 officially
 mandates that the parameter itself must be moved. Default is `false`.

.. code-block:: c++

 // 'p' is flagged by this check if and only if AllowPartialMove is false
 void move_members_of(pair<Obj, Obj>&& p) {
 pair<Obj, Obj> other;
 other.first = std::move(p.first);
 other.second = std::move(p.second);
 }

 // 'p' is never flagged by this check
 void move_whole_pair(pair<Obj, Obj>&& p) {
 pair<Obj, Obj> other = std::move(p);
 }
**Option:** AllowImplicitMove
 If set to `true`, the check recognizes implicit move ``return param;``
 where ``param`` is a rvalue reference parameter. Default is `false`.

.. code-block:: c++

 A f(A&& a) {
 return a; // no warning with AllowImplicitMove = true
 }
**Option:** IgnoreUnnamedParams
If set to `true`, the check ignores unnamed rvalue reference parameters.
Default is `false`.
**Option:** IgnoreNonDeducedTemplateTypes
 If set to `true`, the check ignores non-deduced template type rvalue
 reference parameters. Default is `false`.

.. code-block:: c++

 template <class T>
 struct SomeClass {
 // Below, 'T' is not deduced and 'T&&' is an rvalue reference type.
 // This will be flagged if and only if IgnoreNonDeducedTemplateTypes is
 // false. One suggested fix would be to specialize the class for 'T' and
 // 'T&' separately (e.g., see std::future), or allow only one of 'T' or
 // 'T&' instantiations of SomeClass (e.g., see std::optional).
 SomeClass(T&& t) { }
 };

 // Never flagged, since 'T' is a forwarding reference in a deduced context
 template <class T>
 void forwarding_ref(T&& t) {
 T other = std::forward<T>(t);
 }
**Option:** MoveFunction
Specify the function used for moving. Default is `::std::move`.
This check implements `F.18
<http://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#f18-for-will-move-from-parameters-pass-by-x-and-stdmove-the-parameter>`_
from the C++ Core Guidelines.

```{title} clang-tidy - cppcoreguidelines-slicing
```

# cppcoreguidelines-slicing

Flags slicing of member variables or vtable. Slicing happens when copying a
derived object into a base object: the members of the derived object (both
member variables and virtual member functions) will be discarded. This can be
misleading especially for member function slicing, for example:

```cpp
struct B { int a; virtual int f(); };
struct D : B { int b; int f() override; };

void use(B b) { // Missing reference, intended?
 b.f(); // Calls B::f.
}

D d;
use(d); // Slice.
```

This check implements [ES.63](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#es63-dont-slice)
and [C.145](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#c145-access-polymorphic-objects-through-pointers-and-references)
from the C++ Core Guidelines.

**Title:** clang-tidy - cppcoreguidelines-special-member-functions

# cppcoreguidelines-special-member-functions

The check finds classes where some but not all of the special member functions
are defined.

By default the compiler defines a copy constructor, copy assignment operator,
move constructor, move assignment operator and destructor. The default can be
suppressed by explicit user-definitions. The relationship between which
functions will be suppressed by definitions of other functions is complicated
and it is advised that all five are defaulted or explicitly defined.

Note that defining a function with `= delete` is considered to be a
definition.

This check implements `C.21
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#rc-five>`_
from the C++ Core Guidelines.

## Options

**Option:** AllowSoleDefaultDtor
When set to `true` (default is `false`), this check will only trigger on
destructors if they are defined and not defaulted.

.. code-block:: c++

 struct A { // This is fine.
 virtual ~A() = default;
 };

 struct B { // This is not fine.
 ~B() {}
 };

 struct C {
 // This is not checked, because the destructor might be defaulted in
 // another translation unit.
 ~C();
 };
**Option:** AllowMissingMoveFunctions
When set to `true` (default is `false`), this check doesn't flag classes
which define no move operations at all. It still flags classes which define
only one of either move constructor or move assignment operator. With this
option enabled, the following class won't be flagged:

.. code-block:: c++

 struct A {
 A(const A&);
 A& operator=(const A&);
 ~A();
 };
**Option:** AllowMissingMoveFunctionsWhenCopyIsDeleted
When set to `true` (default is `false`), this check doesn't flag classes
which define deleted copy operations but don't define move operations. This
flag is related to Google C++ Style Guide `Copyable and Movable Types
<https://google.github.io/styleguide/cppguide.html#Copyable_Movable_Types>`_.
With this option enabled, the following class won't be flagged:

.. code-block:: c++

 struct A {
 A(const A&) = delete;
 A& operator=(const A&) = delete;
 ~A();
 };
**Option:** AllowImplicitlyDeletedCopyOrMove
When set to `true` (default is `false`), this check doesn't flag classes
which implicitly delete copy or move operations.
With this option enabled, the following class won't be flagged:

.. code-block:: c++

 struct A : boost::noncopyable {
 ~A() { std::cout << "dtor\n"; }
 };
**Option:** IgnoreMacros
If set to `true`, the check will not give warnings for classes defined
inside macros. Default is `true`.

```{title} clang-tidy - cppcoreguidelines-use-default-member-init
```

# cppcoreguidelines-use-default-member-init

This check implements [C.48](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#rc-in-class-initializer)
from the C++ Core Guidelines.

The `cppcoreguidelines-use-default-member-init` check is an alias, please see
{doc}`modernize-use-default-member-init <../modernize/use-default-member-init>`
for more information.

**Title:** clang-tidy - cppcoreguidelines-use-enum-class

# cppcoreguidelines-use-enum-class

Finds unscoped (non-class) `enum` declarations and suggests using
`enum class` instead.

This check implements `Enum.3
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#renum-class>`_
from the C++ Core Guidelines."

Example:

```c++
enum E {}; // use "enum class E {};" instead
enum class E {}; // OK

struct S {
 enum E {}; // use "enum class E {};" instead
 // OK with option IgnoreUnscopedEnumsInClasses
};

namespace N {
 enum E {}; // use "enum class E {};" instead
}
```

## Options

**Option:** IgnoreUnscopedEnumsInClasses
When `true`, ignores unscoped ``enum`` declarations in classes.
Default is `false`.
**Option:** IgnoreMacros
When `true`, ignores unscoped ``enum`` declarations within macros.
Default is `false`.

**Title:** clang-tidy - cppcoreguidelines-virtual-class-destructor

# cppcoreguidelines-virtual-class-destructor

Finds virtual classes whose destructor is neither public and virtual
nor protected and non-virtual. A virtual class's destructor should be specified
in one of these ways to prevent undefined behavior.

This check implements
[C.35](http://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#rc-dtor-virtual)
from the C++ Core Guidelines.

Note that this check will diagnose a class with a virtual method regardless of
whether the class is used as a base class or not.

Fixes are available for user-declared and implicit destructors that are either
public and non-virtual or protected and virtual. No fixes are offered for
private destructors. There, the decision whether to make them private and
virtual or protected and non-virtual depends on the use case and is thus left
to the user.

## Example

For example, the following classes/structs get flagged by the check since they
violate guideline **C.35**:

```c++
struct Foo { // NOK, protected destructor should not be virtual
 virtual void f();
protected:
 virtual ~Foo(){}
};

class Bar { // NOK, public destructor should be virtual
 virtual void f();
public:
 ~Bar(){}
};
```

This would be rewritten to look like this:

```c++
struct Foo { // OK, destructor is not virtual anymore
 virtual void f();
protected:
 ~Foo(){}
};

class Bar { // OK, destructor is now virtual
 virtual void f();
public:
 virtual ~Bar(){}
};
```

```{title} clang-tidy - darwin-avoid-spinlock
```

# darwin-avoid-spinlock

Finds usages of `OSSpinlock`, which is deprecated due to potential livelock
problems.

This check will detect following function invocations:

- `OSSpinlockLock`
- `OSSpinlockTry`
- `OSSpinlockUnlock`

The corresponding information about the problem of `OSSpinlock`: <https://blog.postmates.com/why-spinlocks-are-bad-on-ios-b69fc5221058>

```{title} clang-tidy - darwin-dispatch-once-nonstatic
```

# darwin-dispatch-once-nonstatic

Finds declarations of `dispatch_once_t` variables without static or global
storage. The behavior of using `dispatch_once_t` predicates with automatic or
dynamic storage is undefined by libdispatch, and should be avoided.

It is a common pattern to have functions initialize internal static or global
data once when the function runs, but programmers have been known to miss the
static on the `dispatch_once_t` predicate, leading to an uninitialized flag
value at the mercy of the stack.

Programmers have also been known to make `dispatch_once_t` variables be
members of structs or classes, with the intent to lazily perform some expensive
struct or class member initialization only once; however, this violates the
libdispatch requirements.

See the discussion section of
[Apple's dispatch_once documentation](https://developer.apple.com/documentation/dispatch/1447169-dispatch_once)
for more information.

```{title} clang-tidy - fuchsia-default-arguments-calls
```

# fuchsia-default-arguments-calls

Warns if a function or method is called with default arguments.

For example, given the declaration:

```cpp
int foo(int value = 5) { return value; }
```

A function call expression that uses a default argument will be diagnosed.
Calling it without defaults will not cause a warning:

```cpp
foo(); // warning
foo(0); // no warning
```

See the features disallowed in Fuchsia at <https://fuchsia.dev/fuchsia-src/development/languages/c-cpp/cxx>

```{title} clang-tidy - fuchsia-default-arguments-declarations
```

# fuchsia-default-arguments-declarations

Warns if a function or method is declared with default parameters.

For example, the declaration:

```cpp
int foo(int value = 5) { return value; }
```

will cause a warning.

See the features disallowed in Fuchsia at <https://fuchsia.dev/fuchsia-src/development/languages/c-cpp/cxx>

```{title} clang-tidy - fuchsia-header-anon-namespaces
```

# fuchsia-header-anon-namespaces

The `fuchsia-header-anon-namespaces` check is an alias, please see
{doc}`misc-anonymous-namespace-in-header <../misc/anonymous-namespace-in-header>`
for more information.

```{title} clang-tidy - fuchsia-multiple-inheritance
```

# fuchsia-multiple-inheritance

The `fuchsia-multiple-inheritance` check is an alias, please See
{doc}`misc-multiple-inheritance <../misc/multiple-inheritance>` for details.

See the features disallowed in Fuchsia at <https://fuchsia.dev/fuchsia-src/development/languages/c-cpp/cxx>

```{title} clang-tidy - fuchsia-overloaded-operator
```

# fuchsia-overloaded-operator

Warns if an operator is overloaded, except for the assignment (copy and move)
operators.

For example:

```cpp
int operator+(int); // Warning

B &operator=(const B &Other); // No warning
B &operator=(B &&Other) // No warning
```

See the features disallowed in Fuchsia at <https://fuchsia.dev/fuchsia-src/development/languages/c-cpp/cxx>

**Title:** clang-tidy - fuchsia-statically-constructed-objects

# fuchsia-statically-constructed-objects

Warns if global, non-trivial objects with static storage are constructed,
unless the object is statically initialized with a `constexpr` constructor
or has no explicit constructor.

For example:

```c++
class A {};

class B {
public:
 B(int Val) : Val(Val) {}
private:
 int Val;
};

class C {
public:
 constexpr C(int Val) : Val(Val) {}
 C(int Val1, int Val2) : Val(Val1+Val2) {}

private:
 int Val;
};

static A a; // No warning, as there is no explicit constructor
static C c(0); // No warning, as constructor is constexpr

static B b(0); // Warning, as constructor is not constexpr
static C c2(0, 1); // Warning, as constructor is not constexpr

static int i; // No warning, as it is trivial

extern int get_i();
static C c3(get_i());// Warning, as the constructor is dynamically initialized
```

See the features disallowed in Fuchsia at https://fuchsia.dev/fuchsia-src/development/languages/c-cpp/cxx

**Title:** clang-tidy - fuchsia-temporary-objects

# fuchsia-temporary-objects

Warns on construction of specific temporary objects in the Zircon kernel.
If the object should be flagged, the fully qualified type name must be
explicitly passed to the check.

For example, given the list of classes "Foo" and "NS::Bar", all of the
following will trigger the warning:

```c++
Foo();
Foo F = Foo();
func(Foo());

namespace NS {

Bar();

}
```

With the same list, the following will not trigger the warning:

```c++
Foo F; // Non-temporary construction okay
Foo F(param); // Non-temporary construction okay
Foo *F = new Foo(); // New construction okay

Bar(); // Not NS::Bar, so okay
NS::Bar B; // Non-temporary construction okay
```

Note that objects must be explicitly specified in order to be flagged,
and so objects that inherit a specified object will not be flagged.

This check matches temporary objects without regard for inheritance and so a
prohibited base class type does not similarly prohibit derived class types.

```c++
class Derived : Foo {} // Derived is not explicitly disallowed
Derived(); // and so temporary construction is okay
```

## Options

**Option:** Names
A semi-colon-separated list of fully-qualified names of C++ classes that
should not be constructed as temporaries. Default is empty string.
See the features disallowed in Fuchsia at https://fuchsia.dev/fuchsia-src/development/languages/c-cpp/cxx

**Title:** clang-tidy - fuchsia-trailing-return

# fuchsia-trailing-return

Functions that have trailing returns are disallowed, except for those using
`decltype` specifiers and lambda with otherwise unutterable return types.

For example:

```c++
// No warning
int add_one(const int arg) { return arg; }

// Warning
auto get_add_one() -> int (*)(const int) {
 return add_one;
}
```

Exceptions are made for lambdas and `decltype` specifiers:

```c++
// No warning
auto lambda = [](double x, double y) -> double {return x + y;};

// No warning
template <typename T1, typename T2>
auto fn(const T1 &lhs, const T2 &rhs) -> decltype(lhs + rhs) {
 return lhs + rhs;
}
```

See the features disallowed in Fuchsia at https://fuchsia.dev/fuchsia-src/development/languages/c-cpp/cxx

```{title} clang-tidy - fuchsia-virtual-inheritance
```

# fuchsia-virtual-inheritance

Warns if classes are defined with virtual inheritance.

For example, classes should not be defined with virtual inheritance:

```cpp
class B : public virtual A {}; // warning
```

See the features disallowed in Fuchsia at <https://fuchsia.dev/fuchsia-src/development/languages/c-cpp/cxx>

```{title} clang-tidy - google-build-explicit-make-pair
```

# google-build-explicit-make-pair

Check that `make_pair`'s template arguments are deduced.

G++ 4.6 in C++11 mode fails badly if `make_pair`'s template arguments are
specified explicitly, and such use isn't intended in any case.

Corresponding cpplint.py check name: `build/explicit_make_pair`.

```{title} clang-tidy - google-build-namespaces
```

# google-build-namespaces

The `google-build-namespaces` check is an alias, please see
{doc}`misc-anonymous-namespace-in-header <../misc/anonymous-namespace-in-header>`
for more information.

Finds anonymous namespaces in headers.

<https://google.github.io/styleguide/cppguide.html#Namespaces>

Corresponding cpplint.py check name: `build/namespaces`.

```{title} clang-tidy - google-build-using-namespace
```

# google-build-using-namespace

Finds `using namespace` directives.

The check implements the following rule of the
[Google C++ Style Guide](https://google.github.io/styleguide/cppguide.html#Namespaces):

> You may not use a using-directive to make all names from a namespace
> available.

```cpp
// Forbidden -- This pollutes the namespace.
using namespace foo;
```

Corresponding cpplint.py check name: `build/namespaces`.

```{title} clang-tidy - google-default-arguments
```

# google-default-arguments

Checks that default arguments are not given for virtual methods.

See <https://google.github.io/styleguide/cppguide.html#Default_Arguments>

```{title} clang-tidy - google-explicit-constructor
```

# google-explicit-constructor

This check is an alias for
{doc}`misc-explicit-constructor <../misc/explicit-constructor>`.

See <https://google.github.io/styleguide/cppguide.html#Explicit_Constructors>

```{title} clang-tidy - google-global-names-in-headers
```

# google-global-names-in-headers

Flag global namespace pollution in header files. Right now it only triggers on
`using` declarations and directives.

The relevant style guide section is
<https://google.github.io/styleguide/cppguide.html#Namespaces>.

```{title} clang-tidy - google-objc-avoid-nsobject-new
```

# google-objc-avoid-nsobject-new

Finds calls to `+new` or overrides of it, which are prohibited by the
Google Objective-C style guide.

The Google Objective-C style guide forbids calling `+new` or overriding it in
class implementations, preferring `+alloc` and `-init` methods to
instantiate objects.

An example:

```objc
NSDate *now = [NSDate new];
Foo *bar = [Foo new];
```

Instead, code should use `+alloc`/`-init` or class factory methods.

```objc
NSDate *now = [NSDate date];
Foo *bar = [[Foo alloc] init];
```

This check corresponds to the Google Objective-C Style Guide rule
[Do Not Use +new](https://google.github.io/styleguide/objcguide.html#do-not-use-new).

**Title:** clang-tidy - google-objc-avoid-throwing-exception

# google-objc-avoid-throwing-exception

Finds uses of throwing exceptions usages in Objective-C files.

For the same reason as the Google C++ style guide, we prefer not throwing
exceptions from Objective-C code.

The corresponding C++ style guide rule:
https://google.github.io/styleguide/cppguide.html#Exceptions

Instead, prefer passing in `NSError **` and return `BOOL` to indicate
success or failure.

A counterexample:

```objc
- (void)readFile {
 if ([self isError]) {
 @throw [NSException exceptionWithName:...];
 }
}
```

Instead, returning an error via `NSError **` is preferred:

```objc
- (BOOL)readFileWithError:(NSError **)error {
 if ([self isError]) {
 *error = [NSError errorWithDomain:...];
 return NO;
 }
 return YES;
}
```

The corresponding style guide rule:
https://google.github.io/styleguide/objcguide.html#avoid-throwing-exceptions

```{title} clang-tidy - google-objc-function-naming
```

# google-objc-function-naming

Finds function declarations in Objective-C files that do not follow the pattern
described in the Google Objective-C Style Guide.

The corresponding style guide rule can be found here:
<https://google.github.io/styleguide/objcguide.html#function-names>

All function names should be in Pascal case. Functions whose storage class is
not static should have an appropriate prefix.

The following code sample does not follow this pattern:

```objc
static bool is_positive(int i) { return i > 0; }
bool IsNegative(int i) { return i < 0; }
```

The sample above might be corrected to the following code:

```objc
static bool IsPositive(int i) { return i > 0; }
bool *ABCIsNegative(int i) { return i < 0; }
```

**Title:** clang-tidy - google-objc-global-variable-declaration

# google-objc-global-variable-declaration

Finds global variable declarations in Objective-C files that do not follow the
pattern of variable names in Google's Objective-C Style Guide.

The corresponding style guide rule:
https://google.github.io/styleguide/objcguide.html#variable-names

All the global variables should follow the pattern of `g[A-Z].*` (variables)
or `k[A-Z].*` (constants). The check will suggest a variable name that
follows the pattern if it can be inferred from the original name.

For code:

```objc
static NSString* myString = @"hello";
```

The fix will be:

```objc
static NSString* gMyString = @"hello";
```

Another example of constant:

```objc
static NSString* const myConstString = @"hello";
```

The fix will be:

```objc
static NSString* const kMyConstString = @"hello";
```

However for code that prefixed with non-alphabetical characters like:

```objc
static NSString* __anotherString = @"world";
```

The check will give a warning message but will not be able to suggest
a fix. The user needs to fix it on their own.

**Title:** clang-tidy - google-readability-avoid-underscore-in-googletest-name

# google-readability-avoid-underscore-in-googletest-name

Checks whether there are underscores in googletest test suite names and test
names in test macros:

- `TEST`
- `TEST_F`
- `TEST_P`
- `TYPED_TEST`
- `TYPED_TEST_P`

The `FRIEND_TEST` macro is not included.

For example:

```c++
TEST(TestSuiteName, Illegal_TestName) {}
TEST(Illegal_TestSuiteName, TestName) {}
```

would trigger the check. Underscores are not allowed in test suite name nor
test names.

The `DISABLED_` prefix, which may be used to
disable test suites and individual tests, is removed from the test suite
name and test name before checking for underscores.

This check does not propose any fixes.

```{title} clang-tidy - google-readability-braces-around-statements
```

# google-readability-braces-around-statements

The `google-readability-braces-around-statements` check is an alias, please see
{doc}`readability-braces-around-statements <../readability/braces-around-statements>`
for more information.

```{title} clang-tidy - google-readability-casting
```

# google-readability-casting

The `google-readability-casting` check is an alias, please see
{doc}`modernize-avoid-c-style-cast <../modernize/avoid-c-style-cast>`
for more information.

Finds usages of C-style casts.

<https://google.github.io/styleguide/cppguide.html#Casting>

Corresponding cpplint.py check name: `readability/casting`.

```{title} clang-tidy - google-readability-function-size
```

# google-readability-function-size

The `google-readability-function-size` check is an alias, please see
{doc}`readability-function-size <../readability/function-size>` for more
information.

```{title} clang-tidy - google-readability-namespace-comments
```

# google-readability-namespace-comments

The `google-readability-namespace-comments` check is an alias, please see
{doc}`llvm-namespace-comment <../llvm/namespace-comment>` for more information.

```{title} clang-tidy - google-readability-todo
```

# google-readability-todo

Finds TODO comments without a username or bug number.

The relevant style guide section is
<https://google.github.io/styleguide/cppguide.html#TODO_Comments>.

Corresponding cpplint.py check: `readability/todo`

## Options

```{option} Style

A string specifying the TODO style for fix-it hints. Accepted values are
`Hyphen` and `Parentheses`. Default is `Hyphen`.

* `Hyphen` will format the fix-it as: `// TODO: username - details`.
* `Parentheses` will format the fix-it as: `// TODO(username): details`.
```

```{title} clang-tidy - google-runtime-float
```

# google-runtime-float

Finds uses of `long double` and suggests against their use due to lack of
portability.

The corresponding style guide rule:
<https://google.github.io/styleguide/cppguide.html#Floating-Point_Types>

```{title} clang-tidy - google-runtime-int
```

# google-runtime-int

Finds uses of `short`, `long` and `long long` and suggest replacing them
with `u?intXX(_t)?`.

The corresponding style guide rule:
<https://google.github.io/styleguide/cppguide.html#Integer_Types>.

Corresponding cpplint.py check: `runtime/int`.

## Options

```{option} UnsignedTypePrefix
A string specifying the unsigned type prefix. Default is `uint`.
```

```{option} SignedTypePrefix
A string specifying the signed type prefix. Default is `int`.
```

```{option} TypeSuffix
A string specifying the type suffix. Default is an empty string.
```

```{title} clang-tidy - google-runtime-operator
```

# google-runtime-operator

Finds overloads of unary `operator &`.

<https://google.github.io/styleguide/cppguide.html#Operator_Overloading>

Corresponding cpplint.py check name: `runtime/operator`.

**Title:** clang-tidy - google-upgrade-googletest-case

# google-upgrade-googletest-case

Finds uses of deprecated Google Test version 1.9 APIs with names containing
`case` and replaces them with equivalent APIs with `suite`.

All names containing `case` are being replaced to be consistent with the
meanings of "test case" and "test suite" as used by the International
Software Testing Qualifications Board and ISO 29119.

The new names are a part of Google Test version 1.9 (release pending). It is
recommended that users update their dependency to version 1.9 and then use this
check to remove deprecated names.

The affected APIs are:

- Member functions of `testing::Test`, `testing::TestInfo`,
 `testing::TestEventListener`, `testing::UnitTest`, and any type
 inheriting from these types
- The macros `TYPED_TEST_CASE`, `TYPED_TEST_CASE_P`,
 `REGISTER_TYPED_TEST_CASE_P`, and `INSTANTIATE_TYPED_TEST_CASE_P`
- The type alias `testing::TestCase`

Examples of fixes created by this check:

```c++
class FooTest : public testing::Test {
public:
 static void SetUpTestCase();
 static void TearDownTestCase();
};

TYPED_TEST_CASE(BarTest, BarTypes);
```

becomes

```c++
class FooTest : public testing::Test {
public:
 static void SetUpTestSuite();
 static void TearDownTestSuite();
};

TYPED_TEST_SUITE(BarTest, BarTypes);
```

For better consistency of user code, the check renames both virtual and
non-virtual member functions with matching names in derived types. The check
tries to provide only a warning when a fix cannot be made safely, as is the case
with some template and macro uses.

```{title} clang-tidy - linuxkernel-must-check-errs
```

# linuxkernel-must-check-errs

Checks Linux kernel code to see if it uses the results from the functions in
`linux/err.h`. Also checks to see if code uses the results from functions
that directly return a value from one of these error functions.

This is important in the Linux kernel because `ERR_PTR`, `PTR_ERR`,
`IS_ERR`, `IS_ERR_OR_NULL`, `ERR_CAST`, and `PTR_ERR_OR_ZERO` return
values must be checked, since positive pointers and negative error codes are
being used in the same context. These functions are marked with
`__attribute__((warn_unused_result))`, but some kernel versions do not have
this warning enabled for clang.

Examples:

```c
/* Trivial unused call to an ERR function */
PTR_ERR_OR_ZERO(some_function_call());

/* A function that returns ERR_PTR. */
void *fn() { ERR_PTR(-EINVAL); }

/* An invalid use of fn. */
fn();
```

```{title} clang-tidy - llvm-else-after-return
```

# llvm-else-after-return

The `llvm-else-after-return` check is an alias, please see
{doc}`readability-else-after-return <../readability/else-after-return>`
for more information.

**Title:** clang-tidy - llvm-formatv-string

# llvm-formatv-string

Validates `llvm::formatv` format strings against the provided arguments,
diagnosing mismatched argument counts, unused arguments, and mixed index
styles.

This check diagnoses the following issues:

- The number of replacement indices in the format string does not match the
 number of arguments provided.
- A format string does not use one of the given arguments.
- Mixing of automatic and explicit indices (e.g. `{} {1}`).

```c++
// warning: format string requires 2 arguments, but 1 argument was provided
llvm::formatv("{0} {1}", x);

// warning: format string mixes automatic and explicit indices
llvm::formatv("{} {1}", x, y);

// warning: argument unused in format string
llvm::formatv("{0} {2}", x, y, z);

// OK.
llvm::formatv("{0} {1}", x, y);
llvm::formatv("{} {}", x, y);
llvm::formatv("{0} {0}", x);
```

The check only operates on calls where the format string is a string literal.
Dynamic format strings are not diagnosed.

## Options

**Option:** AdditionalFunctions
A semicolon-separated list of additional fully qualified function names to
check, beyond ``llvm::formatv`` and ``llvm::createStringErrorV``. Each
function must be a variadic template whose last parameter is a parameter
pack. The format string is assumed to be the parameter immediately preceding
the pack.

For example, to check ``::mylib::log(Level, const char *Fmt, Ts&&...)`` set
this option to `::mylib::log`.

Default is the empty string.

```{title} clang-tidy - llvm-header-guard
```

# llvm-header-guard

Finds and fixes header guards that do not adhere to LLVM style.

```{title} clang-tidy - llvm-include-order
```

# llvm-include-order

Checks the correct order of `#includes`.

See <https://llvm.org/docs/CodingStandards.html#include-style>

**Title:** clang-tidy - llvm-namespace-comment

# llvm-namespace-comment

`google-readability-namespace-comments` redirects here as an alias for this
check.

Checks that long namespaces have a closing comment.

https://llvm.org/docs/CodingStandards.html#namespace-indentation

https://google.github.io/styleguide/cppguide.html#Namespaces

```c++
namespace n1 {
void f();
}

// becomes

namespace n1 {
void f();
} // namespace n1
```

## Options

**Option:** ShortNamespaceLines
Requires the closing brace of the namespace definition to be followed by a
closing comment if the body of the namespace has more than
`ShortNamespaceLines` lines of code. The value is an unsigned integer that
defaults to `1U`.
**Option:** SpacesBeforeComments
An unsigned integer specifying the number of spaces before the comment
closing a namespace definition. Default is `1U`.
**Option:** AllowOmittingNamespaceComments
When `true`, the check will accept if no namespace comment is present.
The check will only fail if the specified namespace comment is different
than expected. Default is `false`.

**Title:** clang-tidy - llvm-prefer-isa-or-dyn-cast-in-conditionals

# llvm-prefer-isa-or-dyn-cast-in-conditionals

Looks at conditionals and finds and replaces cases of `cast<>`,
which will assert rather than return a null pointer, and
`dyn_cast<>` where the return value is not captured. Additionally,
finds and replaces cases that match the pattern ``var &&
isa<X>(var)`, where `var`` is evaluated twice.

```c++
// Finds these:
if (auto x = cast<X>(y)) {}
// is replaced by:
if (auto x = dyn_cast<X>(y)) {}

if (cast<X>(y)) {}
// is replaced by:
if (isa<X>(y)) {}

if (dyn_cast<X>(y)) {}
// is replaced by:
if (isa<X>(y)) {}

if (var && isa<T>(var)) {}
// is replaced by:
if (isa_and_nonnull<T>(var.foo())) {}

// Other cases are ignored, e.g.:
if (auto f = cast<Z>(y)->foo()) {}
if (cast<Z>(y)->foo()) {}
if (X.cast(y)) {}
```

```{title} clang-tidy - llvm-prefer-register-over-unsigned
```

# llvm-prefer-register-over-unsigned

Finds historical use of `unsigned` to hold vregs and physregs and rewrites
them to use `Register`.

Currently this works by finding all variables of unsigned integer type whose
initializer begins with an implicit cast from `Register` to `unsigned`.

```cpp
void example(MachineOperand &MO) {
 unsigned Reg = MO.getReg();
 ...
}
```

becomes:

```cpp
void example(MachineOperand &MO) {
 Register Reg = MO.getReg();
 ...
}
```

**Title:** clang-tidy - llvm-prefer-static-over-anonymous-namespace

# llvm-prefer-static-over-anonymous-namespace

Finds function and variable declarations inside anonymous namespace and
suggests replacing them with `static` declarations.

The [LLVM Coding Standards](https://llvm.org/docs/CodingStandards.html#restrict-visibility)
recommend keeping anonymous namespaces as small as possible and only use them
for class declarations. For functions and variables the `static` specifier
should be preferred for restricting visibility.

For example non-compliant code:

```c++
namespace {

class StringSort {
public:
 StringSort(...)
 bool operator<(const char *RHS) const;
};

// warning: place method definition outside of an anonymous namespace
bool StringSort::operator<(const char *RHS) const {}

// warning: prefer using 'static' for restricting visibility
void runHelper() {}

// warning: prefer using 'static' for restricting visibility
int myVariable = 42;

}
```

Should become:

```c++
// Small anonymous namespace for class declaration
namespace {

class StringSort {
public:
 StringSort(...)
 bool operator<(const char *RHS) const;
};

}

// placed method definition outside of the anonymous namespace
bool StringSort::operator<(const char *RHS) const {}

// used 'static' instead of an anonymous namespace
static void runHelper() {}

// used 'static' instead of an anonymous namespace
static int myVariable = 42;
```

## Options

**Option:** AllowVariableDeclarations
When `true`, allow variable declarations to be in anonymous namespace.
Default value is `true`.
**Option:** AllowMemberFunctionsInClass
When `true`, only methods defined in anonymous namespace outside of the
corresponding class will be warned. Default value is `true`.

```{title} clang-tidy - llvm-qualified-auto
```

# llvm-qualified-auto

The `llvm-qualified-auto check` is an alias, please see
{doc}`readability-qualified-auto <../readability/qualified-auto>`
for more information.

**Title:** clang-tidy - llvm-redundant-casting

# llvm-redundant-casting

Points out uses of `cast<>`, `dyn_cast<>`, `isa<>` and their `or_null`
variants that are unnecessary because the argument already is of the target
type, or a derived type thereof.

```c++
struct A {};
A a;
// Finds:
A x = cast<A>(a);
// replaced by:
A x = a;

struct B : public A {};
B b;
// Finds:
A y = cast<A>(b);
// replaced by:
A y = b;

struct C : public A {};
C c;
// Finds:
bool r1 = isa<A>(a) // always true
bool r2 = isa<A>(b) // always true
bool r3 = isa<B, C>(c) // always true
```

Supported functions:
 - `llvm::cast`
 - `llvm::cast_or_null`
 - `llvm::cast_if_present`
 - `llvm::dyn_cast`
 - `llvm::dyn_cast_or_null`
 - `llvm::dyn_cast_if_present`
 - `llvm::isa`
 - `llvm::isa_and_nonnull`
 - `llvm::isa_and_present`

**Title:** clang-tidy - llvm-twine-local

# llvm-twine-local

Looks for local `Twine` variables which are prone to use after frees and
should be generally avoided.

```c++
static Twine Moo = Twine("bark") + "bah";

// becomes

static std::string Moo = (Twine("bark") + "bah").str();
```

The `Twine` does not own the memory of its contents, so it is not
recommended to use `Twine` created from temporary strings or string literals.

```c++
static Twine getModuleIdentifier(StringRef moduleName) {
 return moduleName + "_module";
}
void foo() {
 Twine result = getModuleIdentifier(std::string{"abc"} + "def");
 // temporary std::string is destroyed here, result is dangling
}
```

After applying this fix-it hints, the code will use `std::string` instead of
`Twine` for local variables. However, `Twine` has lots of methods that
are incompatible with `std::string`, so the user may need to adjust the code
manually after applying the fix-it hints.

**Title:** clang-tidy - llvm-type-switch-case-types

# llvm-type-switch-case-types

Finds `llvm::TypeSwitch::Case` calls with redundant explicit template
arguments that can be inferred from the lambda parameter type.

This check identifies two patterns:

1. **Redundant explicit type**: When the lambda parameter type matches the
 `Case` template argument, the explicit type can be removed.

2. **Auto parameter with explicit type**: When a lambda uses `auto` but
 `Case` has an explicit template argument, suggests using an explicit
 type in the lambda instead.

## Example

```c++
llvm::TypeSwitch<Base *, int>(base)
 .Case<DerivedA>([](DerivedA *a) { return 1; }) // Redundant.
 .Case<DerivedB>([](auto b) { return 2; }); // `auto` with explicit type.
```

Transforms to:

```c++
llvm::TypeSwitch<Base *, int>(base)
 .Case([](DerivedA *a) { return 1; }) // Type inferred from lambda.
 .Case<DerivedB>([](auto b) { return 2; }); // Warning only.
```

Note: The second case (`auto` parameter) only emits a warning without a
fix-it, because the deduced type of `auto` depends on `dyn_cast` behavior
which varies between pointer types and MLIR handle types.

```{title} clang-tidy - llvm-use-new-mlir-op-builder
```

# llvm-mlir-op-builder

Checks for uses of MLIR's old/to be deprecated `OpBuilder::create<T>` form
and suggests using `T::create` instead.

## Example

```cpp
builder.create<FooOp>(builder.getUnknownLoc(), "baz");
```

Transforms to:

```cpp
FooOp::create(builder, builder.getUnknownLoc(), "baz");
```

**Title:** clang-tidy - llvm-use-ranges

# llvm-use-ranges

Finds calls to STL library iterator algorithms that could be replaced with
LLVM range-based algorithms from `llvm/ADT/STLExtras.h`.

## Example

```c++
auto it = std::find(vec.begin(), vec.end(), value);
bool all = std::all_of(vec.begin(), vec.end(),
 [](int x) { return x > 0; });
```

Transforms to:

```c++
auto it = llvm::find(vec, value);
bool all = llvm::all_of(vec, [](int x) { return x > 0; });
```

## Supported algorithms

Calls to the following STL algorithms are checked:

`std::accumulate`,
`std::adjacent_find`,
`std::all_of`,
`std::any_of`,
`std::binary_search`,
`std::copy_if`,
`std::copy`,
`std::count_if`,
`std::count`,
`std::equal`,
`std::fill`,
`std::find_if_not`,
`std::find_if`,
`std::find`,
`std::for_each`,
`std::includes`,
`std::is_sorted`,
`std::lower_bound`,
`std::max_element`,
`std::min_element`,
`std::mismatch`,
`std::none_of`,
`std::partition_point`,
`std::partition`,
`std::remove_if`,
`std::replace_copy_if`,
`std::replace_copy`,
`std::replace`,
`std::search`,
`std::stable_sort`,
`std::transform`,
`std::uninitialized_copy`,
`std::unique`,
`std::upper_bound`.

The check will add the necessary `#include "llvm/ADT/STLExtras.h"` directive
when applying fixes.

**Title:** clang-tidy - llvm-use-vector-utils

# llvm-use-vector-utils

Finds calls to `llvm::to_vector` with `llvm::map_range` or
`llvm::make_filter_range` that can be replaced with the more concise
`llvm::map_to_vector` and `llvm::filter_to_vector` utilities from
`llvm/ADT/SmallVectorExtras.h`.

The check will add the necessary `#include "llvm/ADT/SmallVectorExtras.h"`
directive when applying fixes.

## Example

```c++
auto v1 = llvm::to_vector(llvm::map_range(container, func));
auto v2 = llvm::to_vector(llvm::make_filter_range(container, pred));
auto v3 = llvm::to_vector<4>(llvm::map_range(container, func));
auto v4 = llvm::to_vector<4>(llvm::make_filter_range(container, pred));
```

Transforms to:

```c++
auto v1 = llvm::map_to_vector(container, func);
auto v2 = llvm::filter_to_vector(container, pred);
auto v3 = llvm::map_to_vector<4>(container, func);
auto v4 = llvm::filter_to_vector<4>(container, pred);
```

```{title} clang-tidy - llvmlibc-callee-namespace
```

# llvmlibc-callee-namespace

Checks all calls resolve to functions within correct namespace.

```cpp
// Implementation inside the LIBC_NAMESPACE namespace.
// Correct if:
// - LIBC_NAMESPACE is a macro
// - LIBC_NAMESPACE expansion starts with `__llvm_libc`
namespace LIBC_NAMESPACE {

// Allow calls with the fully qualified name.
LIBC_NAMESPACE::strlen("hello");

// Allow calls to compiler provided functions.
(void)__builtin_abs(-1);

// Bare calls are allowed as long as they resolve to the correct namespace.
strlen("world");

// Disallow calling into functions in the global namespace.
::strlen("!");

} // namespace LIBC_NAMESPACE
```

**Title:** clang-tidy - llvmlibc-implementation-in-namespace

# llvmlibc-implementation-in-namespace

Checks that all declarations in the llvm-libc implementation are within the
correct namespace.

```c++
// Implementation inside the LIBC_NAMESPACE_DECL namespace.
// Correct if:
// - LIBC_NAMESPACE_DECL is a macro
// - LIBC_NAMESPACE_DECL expansion starts with `[[gnu::visibility("hidden")]] __llvm_libc`
namespace LIBC_NAMESPACE_DECL {
 LLVM_LIBC_FUNCTION(char *, strcpy, (char *dest, const char *src)) {}
 // Namespaces within LIBC_NAMESPACE_DECL namespace are allowed.
 namespace inner {
 int localVar = 0;
 }
 // Functions with C linkage are allowed.
 extern "C" void str_fuzz() {}
}

// Incorrect: implementation not in the LIBC_NAMESPACE_DECL namespace.
LLVM_LIBC_FUNCTION(char *, strcpy, (char *dest, const char *src)) {}

// Incorrect: outer most namespace is not the LIBC_NAMESPACE_DECL macro.
namespace something_else {
 LLVM_LIBC_FUNCTION(char *, strcpy, (char *dest, const char *src)) {}
}

// Incorrect: outer most namespace expansion does not start with `[[gnu::visibility("hidden")]] __llvm_libc`.
#define LIBC_NAMESPACE_DECL custom_namespace
namespace LIBC_NAMESPACE_DECL {
 LLVM_LIBC_FUNCTION(char *, strcpy, (char *dest, const char *src)) {}
}
```

```{title} clang-tidy - llvmlibc-inline-function-decl
```

# llvmlibc-inline-function-decl

Checks that all implicitly and explicitly inline functions in header files are
tagged with the `LIBC_INLINE` macro, except for functions implicit to classes
or deleted functions.
See the [libc style guide](https://libc.llvm.org/dev/code_style.html) for more
information about this macro.

**Title:** clang-tidy - llvmlibc-restrict-system-libc-headers

# llvmlibc-restrict-system-libc-headers

Finds includes of system libc headers not provided by the compiler within
llvm-libc implementations.

```c++
#include <stdio.h> // Not allowed because it is part of system libc.
#include <stddef.h> // Allowed because it is provided by the compiler.
#include "internal/stdio.h" // Allowed because it is NOT part of system libc.
```

This check is necessary because accidentally including system libc headers can
lead to subtle and hard to detect bugs. For example consider a system libc
whose `dirent` struct has slightly different field ordering than llvm-libc.
While this will compile successfully, this can cause issues during runtime
because they are ABI incompatible.

## Options

**Option:** Includes
A string containing a comma separated glob list of allowed include
filenames. Similar to the -checks glob list for running clang-tidy itself,
the two wildcard characters are `*` and `-`, to include and exclude globs,
respectively. The default is `-*`, which disallows all includes.

This can be used to allow known safe includes such as Linux development
headers. See :doc:`portability-restrict-system-includes
<../portability/restrict-system-includes>` for more
details.

```{title} clang-tidy - misc-anonymous-namespace-in-header
```

# misc-anonymous-namespace-in-header

Finds anonymous namespaces in headers.

Anonymous namespaces in headers can lead to One Definition Rule (ODR)
violations because each translation unit including the header will get its
own unique version of the symbols. This increases binary size and can cause
confusing link-time errors.

## References

This check corresponds to the CERT C++ Coding Standard rule
[DCL59-CPP. Do not define an unnamed namespace in a header file](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/declarations-and-initialization-dcl/dcl59-cpp/).

Corresponding cpplint.py check name: `build/namespaces`.

```{title} clang-tidy - misc-confusable-identifiers
```

# misc-confusable-identifiers

Warn about confusable identifiers, i.e. identifiers that are visually close to
each other, but use different Unicode characters. This detects a potential
attack described in [CVE-2021-42574](https://www.cve.org/CVERecord?id=CVE-2021-42574).

Example:

```text
int fo; // Initial character is U+0066 (LATIN SMALL LETTER F).
int 𝐟o; // Initial character is U+1D41F (MATHEMATICAL BOLD SMALL F) not U+0066 (LATIN SMALL LETTER F).
```

**Title:** clang-tidy - misc-const-correctness

# misc-const-correctness

Finds local variables and function parameters which could be declared as
`const` but are not.

Declaring variables as `const` is required or recommended by many coding
guidelines, such as:
[ES.25](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#es25-declare-an-object-const-or-constexpr-unless-you-want-to-modify-its-value-later-on)
from the C++ Core Guidelines.

Please note that this check's analysis is type-based only. Variables that are
not modified but used to create a non-const handle that might escape the scope
are not diagnosed as potential `const`.

```c++
// Declare a variable, which is not ``const`` ...
int i = 42;
// but use it as read-only. This means that `i` can be declared ``const``.
int result = i * i; // Before transformation
int const result = i * i; // After transformation
```

The check can analyze values, pointers and references and pointees:

```c++
// Normal values like built-ins or objects.
int potential_const_int = 42; // Before transformation
int const potential_const_int = 42; // After transformation
int copy_of_value = potential_const_int;

MyClass could_be_const; // Before transformation
MyClass const could_be_const; // After transformation
could_be_const.const_qualified_method();

// References can be declared const as well.
int &reference_value = potential_const_int; // Before transformation
int const& reference_value = potential_const_int; // After transformation
int another_copy = reference_value;

// The similar semantics of pointers are analyzed.
int *pointer_variable = &potential_const_int; // Before transformation
int const*const pointer_variable = &potential_const_int; // After transformation, both pointer itself and pointee are supported.
int last_copy = *pointer_variable;
```

The automatic code transformation is only applied to variables that are
declared in single declarations. You may want to prepare your code base with
readability-isolate-declaration first.

Note that there is the check
cppcoreguidelines-avoid-non-const-global-variables
to enforce `const` correctness on all globals.

## Limitations

The check does not run on `C` code.

The check will not analyze templated variables, template functions or variables
that are instantiation dependent. Different instantiations can result
in different `const` correctness properties and in general it is not
possible to find all instantiations of a template. The template might
be used differently in an independent translation unit.

## Options

**Option:** AnalyzeValues
Enable or disable the analysis of ordinary value variables, like
``int i = 42;``. Default is `true`.

.. code-block:: c++

 // Warning
 int i = 42;
 // No warning
 int const i = 42;

 // Warning
 int a[] = {42, 42, 42};
 // No warning
 int const a[] = {42, 42, 42};
**Option:** AnalyzeReferences
Enable or disable the analysis of reference variables, like
``int &ref = i;``. Default is `true`.

.. code-block:: c++

 int i = 42;
 // Warning
 int& ref = i;
 // No warning
 int const& ref = i;
**Option:** AnalyzePointers
Enable or disable the analysis of pointers variables, like
``int *ptr = &i;``. For specific checks, see
:option:`WarnPointersAsValues` and :option:`WarnPointersAsPointers`.
Default is `true`.
**Option:** AnalyzeAutoVariables
Enable or disable the analysis of variables declared with ``auto``,
such as ``auto i = 10;`` or ``auto *ptr = &i``. Default is `true`.
**Option:** AnalyzeLambdas
Enable or disable the analysis of lambda variables, like
``auto f = [] { return 10; };``. For this option to have any
effect, `AnalyzeAutoVariables` must be `true` as well.
Default is `true`.
**Option:** AnalyzeParameters
Enable or disable the analysis of function parameters, like
``void foo(int* ptr)``. Only reference and pointer parameters are analyzed.
Unnamed parameters, member functions (including constructors) and lambdas are
excluded from the analysis. Default is `true`.

.. code-block:: c++

 // Warning
 void function(int& param) {}
 // No warning
 void function(const int& param, int&) {}
**Option:** WarnPointersAsValues
This option enables the suggestion for ``const`` of the pointer itself.
Pointer values have two possibilities to be ``const``, the pointer
and the value pointing to. Default is `false`.

.. code-block:: c++

 int value = 42;

 // Warning
 const int * pointer_variable = &value;
 // No warning
 const int *const pointer_variable = &value;
**Option:** WarnPointersAsPointers
This option enables the suggestion for ``const`` of the value pointing to.
Default is `true`.

Requires :option:`AnalyzePointers` to be `true`.

.. code-block:: c++

 int value = 42;

 // No warning
 const int *const pointer_variable = &value;
 // Warning
 int *const pointer_variable = &value;
**Option:** TransformValues
Provides fixit-hints for value types that automatically add ``const`` if
its a single declaration. Default is `true`.

.. code-block:: c++

 // Before
 int value = 42;
 // After
 int const value = 42;

 // Before
 int a[] = {42, 42, 42};
 // After
 int const a[] = {42, 42, 42};

 // Result is modified later in its life-time. No diagnostic and fixit hint will be emitted.
 int result = value * 3;
 result -= 10;
**Option:** TransformReferences
Provides fixit-hints for reference types that automatically add ``const`` if
its a single declaration. Default is `true`.

.. code-block:: c++

 // This variable could still be a constant. But because there is a non-const reference to
 // it, it can not be transformed (yet).
 int value = 42;
 // The reference 'ref_value' is not modified and can be made 'const int &ref_value = value;'
 // Before
 int &ref_value = value;
 // After
 int const &ref_value = value;

 // Result is modified later in its life-time. No diagnostic and fixit hint will be emitted.
 int result = ref_value * 3;
 result -= 10;
**Option:** TransformPointersAsValues
Provides fixit-hints for pointers if their pointee is not changed. This does
not analyze if the value-pointed-to is unchanged! Default is `false`.

Requires 'WarnPointersAsValues' to be 'true'.

.. code-block:: c++

 int value = 42;

 // Before
 const int * pointer_variable = &value;
 // After
 const int *const pointer_variable = &value;

 // Before
 const int * a[] = {&value, &value};
 // After
 const int *const a[] = {&value, &value};

 // Before
 int *ptr_value = &value;
 // After
 int *const ptr_value = &value;

 int result = 100 * (*ptr_value); // Does not modify the pointer itself.
 // This modification of the pointee is still allowed and not diagnosed.
 *ptr_value = 0;

 // The following pointer may not become a 'int *const'.
 int *changing_pointee = &value;
 changing_pointee = &result;
**Option:** TransformPointersAsPointers
Provides fix-it hints for pointers if the value it pointing to is not changed.
Default is `false`.

Requires :option:`WarnPointersAsPointers` to be `true`.

.. code-block:: c++

 int value = 42;

 // Before
 int * pointer_variable = &value;
 // After
 const int * pointer_variable = &value;

 // Before
 int * a[] = {&value, &value};
 // After
 const int * a[] = {&value, &value};
**Option:** AllowedTypes
A semicolon-separated list of names of types that will be excluded from
const-correctness checking. Regular expressions are accepted, e.g.
``[Rr]ef(erence)?$`` matches every type with suffix ``Ref``, ``ref``,
``Reference`` and ``reference``. If a name in the list contains the sequence
`::`, it is matched against the qualified type name
(i.e. ``namespace::Type``), otherwise it is matched against only the type
name (i.e. ``Type``). Default is empty string.

**Title:** clang-tidy - misc-coroutine-hostile-raii

# misc-coroutine-hostile-raii

Detects when objects of certain hostile RAII types persists across suspension
points in a coroutine. Such hostile types include scoped-lockable types and
types belonging to a configurable denylist.

Some objects require that they be destroyed on the same thread that created
them. Traditionally this requirement was often phrased as "must be a local
variable", under the assumption that local variables always work this way.
However this is incorrect with C++20 coroutines, since an intervening
`co_await` may cause the coroutine to suspend and later be resumed on
another thread.

The lifetime of an object that requires being destroyed on the same thread
must not encompass a `co_await` or `co_yield` point. If you create/destroy
an object, you must do so without allowing the coroutine to suspend in the
meantime.

Following types are considered as hostile:

 - Scoped-lockable types: A scoped-lockable object persisting across a
 suspension point is problematic as the lock held by this object could
 be unlocked by a different thread. This would be undefined behaviour.
 This includes all types annotated with the `scoped_lockable` attribute.

 - Types belonging to a configurable denylist.

```c++
// Call some async API while holding a lock.
task coro() {
 const std::lock_guard l(&mu_);

 // Oops! The async Bar function may finish on a different
 // thread from the one that created the lock_guard (and called
 // Mutex::Lock). After suspension, Mutex::Unlock will be called on the wrong thread.
 co_await Bar();
}
```

## Options

**Option:** RAIITypesList
A semicolon-separated list of qualified types which should not be allowed to
persist across suspension points.
Eg: `my::lockable;a::b;::my::other::lockable`
The default value of this option is `std::lock_guard;std::scoped_lock`.
**Option:** AllowedAwaitablesList
A semicolon-separated list of qualified types of awaitables types which can
be safely awaited while having hostile RAII objects in scope.

``co_await``-ing an expression of ``awaitable`` type is considered
safe if the ``awaitable`` type is part of this list.
RAII objects persisting across such a ``co_await`` expression are
considered safe and hence are not flagged.

Example usage:

.. code-block:: c++

 // Consider option AllowedAwaitablesList = "safe_awaitable"
 struct safe_awaitable {
 bool await_ready() noexcept { return false; }
 void await_suspend(std::coroutine_handle<>) noexcept {}
 void await_resume() noexcept {}
 };
 auto wait() { return safe_awaitable{}; }

 task coro() {
 // This persists across both the co_await's but is not flagged
 // because the awaitable is considered safe to await on.
 const std::lock_guard l(&mu_);
 co_await safe_awaitable{};
 co_await wait();
 }

Eg: `my::safe::awaitable;other::awaitable`
Default is an empty string.
**Option:** AllowedCallees
A semicolon-separated list of callee function names which can
be safely awaited while having hostile RAII objects in scope.
Example usage:

.. code-block:: c++

 // Consider option AllowedCallees = "noop"
 task noop() { co_return; }

 task coro() {
 // This persists across the co_await but is not flagged
 // because the awaitable is considered safe to await on.
 const std::lock_guard l(&mu_);
 co_await noop();
 }

Eg: `my::safe::await;other::await`
Default is an empty string.

**Title:** clang-tidy - misc-definitions-in-headers

# misc-definitions-in-headers

Finds non-extern non-inline function and variable definitions in header files,
which can lead to potential ODR violations in case these headers are included
from multiple translation units.

```c++
// Foo.h
int a = 1; // Warning: variable definition.
extern int d; // OK: extern variable.

namespace N {
 int e = 2; // Warning: variable definition.
}

// Warning: variable definition.
const char* str = "foo";

// OK: internal linkage variable definitions are ignored for now.
// Although these might also cause ODR violations, we can be less certain and
// should try to keep the false-positive rate down.
static int b = 1;
const int c = 1;
const char* const str2 = "foo";
constexpr int k = 1;
namespace { int x = 1; }

// Warning: function definition.
int g() {
 return 1;
}

// OK: inline function definition is allowed to be defined multiple times.
inline int e() {
 return 1;
}

class A {
public:
 int f1() { return 1; } // OK: implicitly inline member function definition is allowed.
 int f2();

 static int d;
};

// Warning: not an inline member function definition.
int A::f2() { return 1; }

// OK: class static data member declaration is allowed.
int A::d = 1;

// OK: function template is allowed.
template<typename T>
T f3() {
 T a = 1;
 return a;
}

// Warning: full specialization of a function template is not allowed.
template <>
int f3() {
 int a = 1;
 return a;
}

template <typename T>
struct B {
 void f1();
};

// OK: member function definition of a class template is allowed.
template <typename T>
void B<T>::f1() {}

class CE {
 constexpr static int i = 5; // OK: inline variable definition.
};

inline int i = 5; // OK: inline variable definition.

constexpr int f10() { return 0; } // OK: constexpr function implies inline.

// OK: C++14 variable templates are inline.
template <class T>
constexpr T pi = T(3.1415926L);
```

When `clang-tidy` is invoked with the `--fix-notes` option, this check
provides fixes that automatically add the `inline` keyword to discovered
functions. Please note that the addition of the `inline` keyword to variables
is not currently supported by this check.

**Title:** clang-tidy - misc-explicit-constructor

# misc-explicit-constructor

Checks that constructors callable with a single argument and conversion
operators are marked explicit to avoid the risk of unintentional implicit
conversions.

Consider this example:

```c++
struct S {
 int x;
 operator bool() const { return true; }
};

bool f() {
 S a{1};
 S b{2};
 return a == b;
}
```

The function will return `true`, since the objects are implicitly converted
to `bool` before comparison, which is unlikely to be the intent.

The check will suggest inserting `explicit` before the constructor or
conversion operator declaration. However, copy and move constructors should not
be explicit, as well as constructors taking a single `initializer_list`
argument.

This code:

```c++
struct S {
 S(int a);
 explicit S(const S&);
 operator bool() const;
 ...
```

will become

```c++
struct S {
 explicit S(int a);
 S(const S&);
 explicit operator bool() const;
 ...
```

**Title:** clang-tidy - misc-header-include-cycle

# misc-header-include-cycle

Check detects cyclic `#include` dependencies between user-defined headers.

```c++
// Header A.hpp
#pragma once
#include "B.hpp"

// Header B.hpp
#pragma once
#include "C.hpp"

// Header C.hpp
#pragma once
#include "A.hpp"

// Include chain: A->B->C->A
```

Header files are a crucial part of many C++ programs as they provide a way to
organize declarations and definitions shared across multiple source files.
However, header files can also create problems when they become entangled
in complex dependency cycles. Such cycles can cause issues with compilation
times, unnecessary rebuilds, and make it harder to understand the overall
structure of the code.

To address these issues, a check has been developed to detect cyclic
dependencies between header files, also known as "include cycles".
An include cycle occurs when a header file `A` includes header file `B`,
and `B` (or any subsequent included header file) includes back header file `A`,
resulting in a circular dependency cycle.

This check operates at the preprocessor level and specifically analyzes
user-defined headers and their dependencies. It focuses solely on detecting
include cycles while disregarding other types or function dependencies.
This specialized analysis helps identify and prevent issues related to header
file organization.

By detecting include cycles early in the development process, developers can
identify and resolve these issues before they become more difficult and
time-consuming to fix. This can lead to faster compile times, improved code
quality, and a more maintainable codebase overall. Additionally, by ensuring
that header files are organized in a way that avoids cyclic dependencies,
developers can make their code easier to understand and modify over time.

It's worth noting that only user-defined headers their dependencies are
analyzed, system includes such as standard library headers and third-party
library headers are excluded. System includes are usually well-designed and
free of include cycles, and ignoring them helps to focus on potential issues
within the project's own codebase. This limitation doesn't diminish the
ability to detect `#include` cycles within the analyzed code.

Developers should carefully review any warnings or feedback provided by this
solution. While the analysis aims to identify and prevent include cycles, there
may be situations where exceptions or modifications are necessary. It's
important to exercise judgment and consider the specific context of the
codebase when making adjustments.

## Options

**Option:** IgnoredFilesList
Provides a way to exclude specific files/headers from the warnings raised by
a check. This can be achieved by specifying a semicolon-separated list of
regular expressions or filenames. This option can be used as an alternative
to ``//NOLINT`` when using it is not possible.
The default value of this option is an empty string, indicating that no
files are ignored by default.

**Title:** clang-tidy - misc-include-cleaner

# misc-include-cleaner

Checks for unused and missing includes. Generates findings only for
the main file of a translation unit.
Findings correspond to https://clangd.llvm.org/design/include-cleaner.

Example:

```c++
// foo.h
class Foo{};
// bar.h
#include "baz.h"
class Bar{};
// baz.h
class Baz{};
// main.cc
#include "bar.h" // OK: uses class Bar from bar.h
#include "foo.h" // warning: unused include "foo.h"
Bar bar;
Baz baz; // warning: missing include "baz.h"
```

## Options

**Option:** IgnoreHeaders
A semicolon-separated list of regexes to disable insertion/removal of header
files that match this regex as a suffix. E.g., `foo/.*` disables
insertion/removal for all headers under the directory `foo`. Default is an
empty string, no headers will be ignored.
**Option:** DeduplicateFindings
A boolean that controls whether the check should deduplicate findings for the
same symbol. Defaults to `true`.
**Option:** UnusedIncludes
A boolean that controls whether the check should report unused includes
(includes that are not used directly). Defaults to `true`.
**Option:** MissingIncludes
A boolean that controls whether the check should report missing includes
(header files from which symbols are used but which are not directly included).
Defaults to `true`.

```{title} clang-tidy - misc-misleading-bidirectional
```

# misc-misleading-bidirectional

Warns about unterminated bidirectional unicode sequence, detecting potential attack
as described in the [Trojan Source](https://www.trojansource.codes) attack.

Example:

```cpp
#include <iostream>

int main() {
 bool isAdmin = false;
 /*‮ } ⁦if (isAdmin)⁩ ⁦ begin admins only */
 std::cout << "You are an admin.\n";
 /* end admins only ‮ { ⁦*/
 return 0;
}
```

```{title} clang-tidy - misc-misleading-identifier
```

# misc-misleading-identifier

Finds identifiers that contain Unicode characters with right-to-left direction,
which can be confusing as they may change the understanding of a whole
statement line, as described in [Trojan Source](https://trojansource.codes).

An example of such misleading code follows:

```text
#include <stdio.h>

short int א = (short int)0;
short int ג = (short int)12345;

int main() {
 int א = ג; // a local variable, set to zero?
 printf("ג is %d\n", ג);
 printf("א is %d\n", א);
}
```

```{title} clang-tidy - misc-misplaced-const
```

# misc-misplaced-const

This check diagnoses when a `const` qualifier is applied to a `typedef`/
`using` to a pointer type rather than to the pointee, because such constructs
are often misleading to developers because the `const` applies to the pointer
rather than the pointee.

For instance, in the following code, the resulting type is `int * const`
rather than `const int *`:

```cpp
typedef int *int_ptr;
void f(const int_ptr ptr) {
 *ptr = 0; // potentially quite unexpectedly the int can be modified here
 ptr = 0; // does not compile
}
```

The check does not diagnose when the underlying `typedef`/`using` type is a
pointer to a `const` type or a function pointer type. This is because the
`const` qualifier is less likely to be mistaken because it would be redundant
(or disallowed) on the underlying pointee type.

**Title:** clang-tidy - misc-multiple-inheritance

# misc-multiple-inheritance

Warns if a class inherits from multiple classes that are not pure virtual.

For example, declaring a class that inherits from multiple concrete classes is
disallowed:

```c++
class Base_A {
public:
 virtual int foo() { return 0; }
};

class Base_B {
public:
 virtual int bar() { return 0; }
};

// Warning
class Bad_Child1 : public Base_A, Base_B {};
```

A class that inherits from a pure virtual is allowed:

```c++
class Interface_A {
public:
 virtual int foo() = 0;
};

class Interface_B {
public:
 virtual int bar() = 0;
};

// No warning
class Good_Child1 : public Interface_A, Interface_B {
 virtual int foo() override { return 0; }
 virtual int bar() override { return 0; }
};
```

```{title} clang-tidy - misc-new-delete-overloads
```

# misc-new-delete-overloads

`cert-dcl54-cpp` redirects here as an alias for this check.

The check flags overloaded operator `new()` and operator `delete()`
functions that do not have a corresponding free store function defined within
the same scope.
For instance, the check will flag a class implementation of a non-placement
operator `new()` when the class does not also define a non-placement operator
`delete()` function as well.

The check does not flag implicitly-defined operators, deleted or private
operators, or placement operators.

This check corresponds to CERT C++ Coding Standard rule [DCL54-CPP. Overload allocation and deallocation functions as a pair in the same scope](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/declarations-and-initialization-dcl/dcl54-cpp/).

```{title} clang-tidy - misc-no-recursion
```

# misc-no-recursion

Finds strongly connected functions (by analyzing the call graph for
SCC's (Strongly Connected Components) that are loops),
diagnoses each function in the cycle,
and displays one example of a possible call graph loop (recursion).

References:

- CERT C++ Coding Standard rule [DCL56-CPP. Avoid cycles during initialization of static objects](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/declarations-and-initialization-dcl/dcl56-cpp/).
- JPL Institutional Coding Standard for the C Programming Language
 (JPL DOCID D-60411) rule `2.4 Do not use direct or indirect recursion`.
- OpenCL Specification, Version 1.2 rule [6.9 Restrictions: i. Recursion is not supported.](https://www.khronos.org/registry/OpenCL/specs/opencl-1.2.pdf).

## Limitations

- The check does not handle calls done through function pointers
- The check does not handle C++ destructors

```{title} clang-tidy - misc-non-copyable-objects
```

# misc-non-copyable-objects

`cert-fio38-c` redirects here as an alias for this check.

Flags dereferences and non-pointer declarations of objects that are
not meant to be passed by value, such as C FILE objects or POSIX
`pthread_mutex_t` objects.

## References

This check corresponds to CERT C++ Coding Standard rule [FIO38-C. Do not copy a FILE object](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/input-output-fio/fio38-c/).

```{title} clang-tidy - misc-non-private-member-variables-in-classes
```

# misc-non-private-member-variables-in-classes

`cppcoreguidelines-non-private-member-variables-in-classes` redirects here
as an alias for this check.

Finds classes that contain non-static data members in addition to user-declared
non-static member functions and diagnose all data members declared with a
non-`public` access specifier. The data members should be declared as
`private` and accessed through member functions instead of exposed to derived
classes or class consumers.

## Options

```{option} IgnoreClassesWithAllMemberVariablesBeingPublic

When `true`, allows to completely ignore classes if **all** the member
variables in that class declared with a `public` access specifier.
Default is `false`.
```

```{option} IgnorePublicMemberVariables

When `true`, allows to ignore (not diagnose) **all** the member variables
declared with a `public` access specifier. Default is `false`.
```

**Title:** clang-tidy - misc-override-with-different-visibility

# misc-override-with-different-visibility

Finds virtual function overrides with different visibility than the function
in the base class. This includes for example if a virtual function declared as
`private` is overridden and declared as `public` in a subclass. The
detected change is the modification of visibility resulting from keywords
`public`, `protected`, `private` at overridden virtual functions. The
check applies to any normal virtual function and optionally to destructors or
operators. Use of the `using` keyword is not considered as visibility
change by this check.

```c++
class A {
public:
 virtual void f_pub();
private:
 virtual void f_priv();
};

class B: public A {
public:
 void f_priv(); // warning: changed visibility from private to public
private:
 void f_pub(); // warning: changed visibility from public to private
};

class C: private A {
 // no warning: f_pub becomes private in this case but this is from the
 // private inheritance
};

class D: private A {
public:
 void f_pub(); // warning: changed visibility from private to public
 // 'f_pub' would have private access but is forced to be
 // public
};
```

If the visibility is changed in this way, it can indicate bad design or
programming error.

If a virtual function is private in a subclass but public in the base class, it
can still be accessed from a pointer to the subclass if the pointer is converted
to the base type. Probably private inheritance can be used instead.

A protected virtual function that is made public in a subclass may have valid
use cases but similar (not exactly same) effect can be achieved with the
`using` keyword.

## Options

**Option:** DisallowedVisibilityChange
Controls what kind of change to the visibility will be detected by the check.
Possible values are `any`, `widening`, `narrowing`. For example the
`widening` option will produce warning only if the visibility is changed
from more restrictive (``private``) to less restrictive (``public``).
Default value is `any`.
**Option:** CheckDestructors
If `true`, the check does apply to destructors too. Otherwise destructors
are ignored by the check.
Default value is `false`.
**Option:** CheckOperators
If `true`, the check does apply to overloaded C++ operators (as virtual
member functions) too. This includes other special member functions (like
conversions) too. This option is probably useful only in rare cases because
operators and conversions are not often virtual functions.
Default value is `false`.
**Option:** IgnoredFunctions
This option can be used to ignore the check at specific functions.
To configure this option, a semicolon-separated list of function names
should be provided. The list can contain regular expressions, in this way it
is possible to select all functions of a specific class (like `MyClass::.*`)
or a specific function of any class (like `my_function` or
`::.*::my_function`). The function names are matched at the base class.
Default value is empty string.

```{title} clang-tidy - misc-predictable-rand
```

# misc-predictable-rand

Warns for the usage of `std::rand()`. Pseudorandom number generators use
mathematical algorithms to produce a sequence of numbers with good
statistical properties, but the numbers produced are not genuinely random.
The `std::rand()` function takes a seed (number), runs a mathematical
operation on it and returns the result. By manipulating the seed the result
can be predictable.

## References

This check corresponds to the CERT C Coding Standard rules
[MSC30-C. Do not use the rand() function for generating pseudorandom numbers](https://cmu-sei.github.io/secure-coding-standards/sei-cert-c-coding-standard/rules/miscellaneous-msc/msc30-c/).
[MSC50-CPP. Do not use std::rand() for generating pseudorandom numbers](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/miscellaneous-msc/msc50-cpp/).

**Title:** clang-tidy - misc-redundant-expression

# misc-redundant-expression

Detect redundant expressions which are typically errors due to copy-paste.

Depending on the operator expressions may be

- redundant,

- always `true`,

- always `false`,

- always a constant (zero or one).

Examples:

```c++
((x+1) | (x+1)) // (x+1) is redundant
(p->x == p->x) // always true
(p->x < p->x) // always false
(speed - speed + 1 == 12) // speed - speed is always zero
int b = a | 4 | a // identical expr on both sides
((x=1) | (x=1)) // expression is identical
(DEFINE_1 | DEFINE_1) // same macro on the both sides
((DEF_1 + DEF_2) | (DEF_1+DEF_2)) // expressions differ in spaces only
```

Floats are handled except in the case that NaNs are checked like so:

```c++
int TestFloat(float F) {
 if (F == F) // Identical float values used
 return 1;
 return 0;
}

int TestFloat(float F) {
 // Testing NaN.
 if (F != F && F == F) // does not warn
 return 1;
 return 0;
}
```

```{title} clang-tidy - misc-static-assert
```

# misc-static-assert

`cert-dcl03-c` redirects here as an alias for this check.

Replaces `assert()` with `static_assert()` if the condition is evaluable
at compile time.

The condition of `static_assert()` is evaluated at compile time which is
safer and more efficient.

**Title:** clang-tidy - misc-static-initialization-cycle

# misc-static-initialization-cycle

Finds cyclical initialization of static variables.

The cycle can come from reference to static variables or from (static) function
calls during initialization. Such cycles can cause undefined behavior. In this
context "static" means C++ `static` class members, global variables, global
functions, and `static` variables inside functions.

For the purpose of this check, the initialization of a static variable
*uses* another static variable or function if it appears in the initializer
expression. A function *uses* a static variable or function if the variable
or function appears at any place in the function code (except if the variable
is assigned to). The check can detect cycles in this "usage graph".

The check does not consider conditions in function code and does not follow the
value of static variables (if assigned to another variable). For this reason it
can produce false positives in some cases.

## Examples

```c++
struct S { static int A; };
int B = S::A;
int S::A = B;
```

Cycle in variable initialization.

```c++
int f1(int X, int Y);

struct S { static int A; };

int B = S::A + 1;
int S::A = f1(B, 2);
```

Cyclical initialization: `B` uses value of `S::A`, and `S::A` may use
value of `B` (the check gives always warning regardless of the code of
`f1`).

```c++
struct S { static int A; };
int f1() {
 return S::A;
}
int S::A = f1();
```

This code results in initialization of `S::A` with itself through a function
call. The check would emit a warning in any case when `S::A` appears in
`f1` (even if the return value is not affected by it).

## References

* CERT C++ Coding Standard rule `DCL56-CPP. Avoid cycles during initialization
 of static objects <https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/declarations-and-initialization-dcl/dcl56-cpp/>`_.

**Title:** clang-tidy - misc-throw-by-value-catch-by-reference

# misc-throw-by-value-catch-by-reference

`cert-err09-cpp` and `cert-err61-cpp` redirect here as aliases of this check.

Finds violations of the rule "Throw by value, catch by reference" presented for
example in "C++ Coding Standards" by H. Sutter and A. Alexandrescu, as well as
the CERT C++ Coding Standard rule `ERR61-CPP. Catch exceptions by lvalue
reference <https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/exceptions-and-error-handling-err/err61-cpp/>`_.

Exceptions:
 * Throwing string literals will not be flagged despite being a pointer. They
 are not susceptible to slicing and the usage of string literals is
 idiomatic.
 * Catching character pointers (`char`, `wchar_t`, unicode character
 types) will not be flagged to allow catching string literals.
 * Moved named values will not be flagged as not throwing an anonymous
 temporary. In this case we can be sure that the user knows that the object
 can't be accessed outside catch blocks handling the error.
 * Throwing function parameters will not be flagged as not throwing an
 anonymous temporary. This allows helper functions for throwing.
 * Re-throwing caught exception variables will not be flagged as not throwing
 an anonymous temporary. Although this can usually be done by just writing
 `throw;` it happens often enough in real code.

## Options

**Option:** CheckThrowTemporaries
Triggers detection of violations of the CERT recommendation ERR09-CPP. Throw
anonymous temporaries.
Default is `true`.
**Option:** WarnOnLargeObject
Also warns for any large, trivial object caught by value. Catching a large
object by value is not dangerous but affects the performance negatively. The
maximum size of an object allowed to be caught without warning can be set
using the `MaxSize` option.
Default is `false`.
**Option:** MaxSize
Determines the maximum size of an object allowed to be caught without
warning. Only applicable if :option:`WarnOnLargeObject` is set to `true`. If
the option is set by the user to `std::numeric_limits<uint64_t>::max()` then
it reverts to the default value.
Default is the size of `size_t`.

```{title} clang-tidy - misc-unconventional-assign-operator
```

# misc-unconventional-assign-operator

Finds declarations of assign operators with the wrong return and/or argument
types and definitions with good return type but wrong `return` statements.

- The return type must be `Class&`.
- The assignment may be from the class type by value, const lvalue
 reference, non-const rvalue reference, or from a completely different
 type (e.g. `int`).
- Private and deleted operators are ignored.
- The operator must always return `*this`.

```{title} clang-tidy - misc-uniqueptr-reset-release
```

# misc-uniqueptr-reset-release

Find and replace `unique_ptr::reset(release())` with `std::move()`.

Example:

```cpp
std::unique_ptr<Foo> x, y;
x.reset(y.release()); -> x = std::move(y);
```

If `y` is already rvalue, `std::move()` is not added. `x` and `y` can
also be `std::unique_ptr<Foo>*`.

## Options

```{option} IncludeStyle

A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.
```

```{title} clang-tidy - misc-unused-alias-decls
```

# misc-unused-alias-decls

Finds unused namespace alias declarations.

```cpp
namespace my_namespace {
class C {};
}
namespace unused_alias = ::my_namespace;
```

**Title:** clang-tidy - misc-unused-parameters

# misc-unused-parameters

Finds unused function parameters. Unused parameters may signify a bug in the
code (e.g. when a different parameter is used instead). The suggested fixes
either comment parameter name out or remove the parameter completely, if all
callers of the function are in the same translation unit and can be updated.

The check is similar to the `-Wunused-parameter` compiler diagnostic and
can be used to prepare a codebase to enabling of that diagnostic. By default
the check is more permissive (see `StrictMode`).

```c++
void a(int i) { /*some code that doesn't use `i`*/ }

// becomes

void a(int /*i*/) { /*some code that doesn't use `i`*/ }
```

```c++
static void staticFunctionA(int i);
static void staticFunctionA(int i) { /*some code that doesn't use `i`*/ }

// becomes

static void staticFunctionA()
static void staticFunctionA() { /*some code that doesn't use `i`*/ }
```

## Options

**Option:** StrictMode
When `false` (default value), the check will ignore trivially unused parameters,
i.e. when the corresponding function has an empty body (and in case of
constructors - no constructor initializers). When the function body is empty,
an unused parameter is unlikely to be unnoticed by a human reader, and
there's basically no place for a bug to hide.
**Option:** IgnoreVirtual
Determines whether virtual method parameters should be inspected.
Set to `true` to ignore them. Default is `false`.
**Option:** IgnoreMacroParameters
When `true`, the check will not report unused parameters whose declarations
originate from a macro expansion. This suppresses false positives in code
that uses macros to define function signatures where parameters are
structurally required by the macro but not used in every expansion.
Default is `false`.

```{title} clang-tidy - misc-unused-using-decls
```

# misc-unused-using-decls

Finds unused `using` declarations.

Unused `using` declarations in header files will not be diagnosed
since these using declarations are part of the header's public API.
Allowed header file extensions can be configured via the global
option `HeaderFileExtensions`.

Example:

```cpp
// main.cpp
namespace n { class C; }
using n::C; // Never actually used.
```

**Title:** clang-tidy - misc-use-anonymous-namespace

# misc-use-anonymous-namespace

Finds instances of `static` functions or variables declared at global scope
that could instead be moved into an anonymous namespace.

Anonymous namespaces are the "superior alternative" according to the C++
Standard. `static` was proposed for deprecation, but later un-deprecated to
keep C compatibility [1]. `static` is an overloaded term with different
meanings in different contexts, so it can create confusion.

The following uses of `static` will *not* be diagnosed:

* Functions or variables in header files, since anonymous namespaces in headers
 is considered an antipattern. Allowed header file extensions can be
 configured via the global option `HeaderFileExtensions`.
* `const` or `constexpr` variables, since they already have implicit
 internal linkage in C++.

Examples:

```c++
// Bad
static void foo();
static int x;

// Good
namespace {
 void foo();
 int x;
} // namespace
```

[1] [Undeprecating static](https://www.open-std.org/jtc1/sc22/wg21/docs/cwg_defects.html#1012)

**Title:** clang-tidy - misc-use-internal-linkage

# misc-use-internal-linkage

Detects variables, functions, and classes that can be marked as static or
(in C++) moved into an anonymous namespace to enforce internal linkage.

Any entity that's only used within a single file should be given internal
linkage. Doing so gives the compiler more information, allowing it to better
remove dead code and perform more aggressive optimizations.

Example:

```c++
int v1; // can be marked as static

void fn1() {} // can be marked as static

// already declared as extern
extern int v2;

void fn3(); // without function body in all declaration, maybe external linkage
void fn3();

// === C++-specific ===

struct S1 {}; // can be moved into anonymous namespace

namespace {
 // already in anonymous namespace
 int v2;
 void fn2();
 struct S2 {};
}

// export declarations
export void fn4() {}
export namespace t { void fn5() {} }
export int v2;
export class C {};
```

## Options

**Option:** FixMode
Selects what kind of a fix the check should provide. The default is `UseStatic`.

- `None`
 Don't fix automatically.

- `UseStatic`
 Add ``static`` for internal linkage variable and function.
**Option:** AnalyzeFunctions
Whether to suggest giving functions internal linkage. Default is `true`.
**Option:** AnalyzeVariables
Whether to suggest giving variables internal linkage. Default is `true`.
**Option:** AnalyzeTypes
(C++ only) Whether to suggest giving user-defined types (structs,
classes, unions, and enums) internal linkage. Default is `true`.

**Title:** clang-tidy - modernize-avoid-bind

# modernize-avoid-bind

The check finds uses of `std::bind` and `boost::bind` and replaces them
with lambdas. Lambdas will use value-capture unless reference capture is
explicitly requested with `std::ref` or `boost::ref`.

It supports arbitrary callables including member functions, function objects,
and free functions, and all variations thereof. Anything that you can pass
to the first argument of `bind` should be diagnosable. Currently, the only
known case where a fix-it is unsupported is when the same placeholder is
specified multiple times in the parameter list.

Given:

```c++
int add(int x, int y) { return x + y; }
```

Then:

```c++
void f() {
 int x = 2;
 auto clj = std::bind(add, x, _1);
}
```

is replaced by:

```c++
void f() {
 int x = 2;
 auto clj = [=](auto && arg1) { return add(x, arg1); };
}
```

`std::bind` can be hard to read and can result in larger object files and
binaries due to type information that will not be produced by equivalent
lambdas.

## Options

**Option:** PermissiveParameterList
If the option is set to `true`, the check will append ``auto&&...`` to the end
of every placeholder parameter list. Without this, it is possible for a fix-it
to perform an incorrect transformation in the case where the result of the ``bind``
is used in the context of a type erased functor such as ``std::function`` which
allows mismatched arguments. Default is `false`.
For example:

```c++
int add(int x, int y) { return x + y; }
int foo() {
 std::function<int(int,int)> ignore_args = std::bind(add, 2, 2);
 return ignore_args(3, 3);
}
```

is valid code, and returns `4`. The actual values passed to `ignore_args` are
simply ignored. Without `PermissiveParameterList`, this would be transformed into

```c++
int add(int x, int y) { return x + y; }
int foo() {
 std::function<int(int,int)> ignore_args = [] { return add(2, 2); }
 return ignore_args(3, 3);
}
```

which will *not* compile, since the lambda does not contain an `operator()`
that accepts 2 arguments. With permissive parameter list, it instead generates

```c++
int add(int x, int y) { return x + y; }
int foo() {
 std::function<int(int,int)> ignore_args = [](auto&&...) { return add(2, 2); }
 return ignore_args(3, 3);
}
```

which is correct.

This check requires using C++14 or higher to run.

**Title:** clang-tidy - modernize-avoid-c-arrays

# modernize-avoid-c-arrays

`cppcoreguidelines-avoid-c-arrays` redirects here as an alias for this check.

Finds C-style array types and recommend to use `std::array<>` /
`std::vector<>`. All types of C arrays are diagnosed.

For parameters of incomplete C-style array type, it would be better to
use `std::span` / `gsl::span` as replacement.

However, fix-it are potentially dangerous in header files and are therefore not
emitted right now.

```c++
int a[] = {1, 2}; // warning: do not declare C-style arrays, use 'std::array' instead

int b[1]; // warning: do not declare C-style arrays, use 'std::array' instead

void foo() {
 int c[b[0]]; // warning: do not declare C VLA arrays, use 'std::vector' instead
}

template <typename T, int Size>
class array {
 T d[Size]; // warning: do not declare C-style arrays, use 'std::array' instead

 int e[1]; // warning: do not declare C-style arrays, use 'std::array' instead
};

array<int[4], 2> d; // warning: do not declare C-style arrays, use 'std::array' instead

using k = int[4]; // warning: do not declare C-style arrays, use 'std::array' instead
```

However, the `extern "C"` code is ignored, since it is common to share
such headers between C code, and C++ code.

```c++
// Some header
extern "C" {

int f[] = {1, 2}; // not diagnosed

int j[1]; // not diagnosed

inline void bar() {
 {
 int j[j[0]]; // not diagnosed
 }
}

}
```

Similarly, the `main()` function is ignored. Its second and third parameters
can be either `char* argv[]` or `char** argv`, but cannot be
`std::array<>`.

## Options

**Option:** AllowStringArrays
When set to `true` (default is `false`), variables of character array type
with deduced length, initialized directly from string literals, will be ignored.
This option doesn't affect cases where length can't be deduced, resembling
pointers, as seen in class members and parameters. Example:

.. code:: c++

 const char name[] = "Some name";

**Title:** clang-tidy - modernize-avoid-c-style-cast

# modernize-avoid-c-style-cast

Finds usages of C-style casts.

C-style casts can perform a variety of different conversions (`const_cast`,
`static_cast`, `reinterpret_cast`, or a combination). This makes them
dangerous as the intent is not clear, and they can silently perform unsafe
conversions between incompatible types.

This check is similar to `-Wold-style-cast`, but it suggests automated fixes
in some cases. The reported locations should not be different from the ones
generated by `-Wold-style-cast`.

## Examples

```c++
class A {
 public:
 std::string v;
};

A a;
double *num = (double*)(&a); // Compiles! Hides danger
// num = static_cast<double*>(&a); // Won't compile (good!)
num = reinterpret_cast<double*>(&a); // Compiles, danger is explicit
```

## References

Corresponding cpplint.py check name: `readability/casting`.

```{title} clang-tidy - modernize-avoid-setjmp-longjmp
```

# modernize-avoid-setjmp-longjmp

Flags all call expressions involving `setjmp()` and
`longjmp()` in C++ code.

Exception handling with `throw` and `catch` should be used instead.

## References

This check corresponds to the CERT C++ Coding Standard rule
[ERR52-CPP. Do not use setjmp() or longjmp()](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/exceptions-and-error-handling-err/err52-cpp/).

```{title} clang-tidy - modernize-avoid-variadic-functions
```

# modernize-avoid-variadic-functions

Find all function definitions (but not declarations) of C-style variadic
functions.

Instead of C-style variadic functions, C++ function parameter pack should be
used.

## References

This check corresponds to the CERT C++ Coding Standard rule
[DCL50-CPP. Do not define a C-style variadic function](https://cmu-sei.github.io/secure-coding-standards/sei-cert-cpp-coding-standard/rules/declarations-and-initialization-dcl/dcl50-cpp/).

**Title:** clang-tidy - modernize-concat-nested-namespaces

# modernize-concat-nested-namespaces

Checks for use of nested namespaces such as
`namespace a { namespace b { ... } }`
and suggests changing to the more concise syntax introduced
in C++17: `namespace a::b { ... }`.
Inline namespaces are not modified.

For example:

```c++
namespace n1 {
namespace n2 {
void t();
}
}

namespace n3 {
namespace n4 {
namespace n5 {
void t();
}
}
namespace n6 {
namespace n7 {
void t();
}
}
}

// in c++20
namespace n8 {
inline namespace n9 {
void t();
}
}
```

Will be modified to:

```c++
namespace n1::n2 {
void t();
}

namespace n3 {
namespace n4::n5 {
void t();
}
namespace n6::n7 {
void t();
}
}

// in c++20
namespace n8::inline n9 {
void t();
}
```

**Title:** clang-tidy - modernize-deprecated-headers

# modernize-deprecated-headers

Some headers from C library were deprecated in C++ and are no longer welcome in
C++ codebases. Some have no effect in C++. For more details refer to the C++14
Standard [depr.c.headers] section.

This check replaces C standard library headers with their C++ alternatives and
removes redundant ones.

```c++
// C++ source file...
#include <assert.h>
#include <stdbool.h>

// becomes

#include <cassert>
// No 'stdbool.h' here.
```

Important note: the Standard doesn't guarantee that the C++ headers declare all
the same functions in the global namespace. The check in its current form can
break the code that uses library symbols from the global namespace.

* `<assert.h>`
* `<complex.h>`
* `<ctype.h>`
* `<errno.h>`
* `<fenv.h>` // deprecated since C++11
* `<float.h>`
* `<inttypes.h>`
* `<limits.h>`
* `<locale.h>`
* `<math.h>`
* `<setjmp.h>`
* `<signal.h>`
* `<stdarg.h>`
* `<stddef.h>`
* `<stdint.h>`
* `<stdio.h>`
* `<stdlib.h>`
* `<string.h>`
* `<tgmath.h>` // deprecated since C++11
* `<time.h>`
* `<uchar.h>` // deprecated since C++11
* `<wchar.h>`
* `<wctype.h>`

If the specified standard is older than C++11 the check will only replace
headers deprecated before C++11, otherwise -- every header that appeared in
the previous list.

These headers don't have effect in C++:

* `<iso646.h>`
* `<stdalign.h>`
* `<stdbool.h>`

The checker ignores `include` directives within `extern "C" { ... }` blocks,
since a library might want to expose some API for C and C++ libraries.

```c++
// C++ source file...
extern "C" {
#include <assert.h> // Left intact.
#include <stdbool.h> // Left intact.
}
```

## Options

**Option:** CheckHeaderFile
`clang-tidy` cannot know if the header file included by the currently
analyzed C++ source file is not included by any other C source files.
Hence, to omit false-positives and wrong fixit-hints, we ignore emitting
reports into header files. One can set this option to `true` if they know
that the header files in the project are only used by C++ source files.
Default is `false`.

```{title} clang-tidy - modernize-deprecated-ios-base-aliases
```

# modernize-deprecated-ios-base-aliases

Detects usage of the deprecated member types of `std::ios_base` and replaces
those that have a non-deprecated equivalent.

| Deprecated member type | Replacement |
| -------------------------- | ------------------------- |
| `std::ios_base::io_state` | `std::ios_base::iostate` |
| `std::ios_base::open_mode` | `std::ios_base::openmode` |
| `std::ios_base::seek_dir` | `std::ios_base::seekdir` |
| `std::ios_base::streamoff` | |
| `std::ios_base::streampos` | |

**Title:** clang-tidy - modernize-loop-convert

# modernize-loop-convert

This check converts `for(...; ...; ...)` loops to use the new range-based
loops in C++11.

Three kinds of loops can be converted:

- Loops over statically allocated arrays.
- Loops over containers, using iterators.
- Loops over array-like containers, using `operator[]` and `at()`.

## MinConfidence option

### risky

In loops where the container expression is more complex than just a
reference to a declared expression (a variable, function, enum, etc.),
and some part of it appears elsewhere in the loop, we lower our confidence
in the transformation due to the increased risk of changing semantics.
Transformations for these loops are marked as `risky`, and thus will only
be converted if the minimum required confidence level is set to `risky`.

```c++
int arr[10][20];
int l = 5;

for (int j = 0; j < 20; ++j)
 int k = arr[l][j] + l; // using l outside arr[l] is considered risky

for (int i = 0; i < obj.getVector().size(); ++i)
 obj.foo(10); // using 'obj' is considered risky
```

See
Range-based loops evaluate end() only once
for an example of an incorrect transformation when the minimum required confidence
level is set to `risky`.

### reasonable (Default)

If a loop calls `.end()` or `.size()` after each iteration, the
transformation for that loop is marked as `reasonable`, and thus will
be converted if the required confidence level is set to `reasonable`
(default) or lower.

```c++
// using size() is considered reasonable
for (int i = 0; i < container.size(); ++i)
 cout << container[i];
```

### safe

Any other loops that do not match the above criteria to be marked as
`risky` or `reasonable` are marked `safe`, and thus will be converted
if the required confidence level is set to `safe` or lower.

```c++
int arr[] = {1,2,3};

for (int i = 0; i < 3; ++i)
 cout << arr[i];
```

## Example

Original:

```c++
const int N = 5;
int arr[] = {1,2,3,4,5};
vector<int> v;
v.push_back(1);
v.push_back(2);
v.push_back(3);

// safe conversion
for (int i = 0; i < N; ++i)
 cout << arr[i];

// reasonable conversion
for (vector<int>::iterator it = v.begin(); it != v.end(); ++it)
 cout << *it;

// reasonable conversion
for (vector<int>::iterator it = begin(v); it != end(v); ++it)
 cout << *it;

// reasonable conversion
for (vector<int>::iterator it = std::begin(v); it != std::end(v); ++it)
 cout << *it;

// reasonable conversion
for (int i = 0; i < v.size(); ++i)
 cout << v[i];

// reasonable conversion
for (int i = 0; i < size(v); ++i)
 cout << v[i];
```

After applying the check with minimum confidence level set to
`reasonable` (default):

```c++
const int N = 5;
int arr[] = {1,2,3,4,5};
vector<int> v;
v.push_back(1);
v.push_back(2);
v.push_back(3);

// safe conversion
for (auto & elem : arr)
 cout << elem;

// reasonable conversion
for (auto & elem : v)
 cout << elem;

// reasonable conversion
for (auto & elem : v)
 cout << elem;
```

## Reverse Iterator Support

The converter is also capable of transforming iterator loops which use
`rbegin` and `rend` for looping backwards over a container. Out of the box
this will automatically happen in C++20 mode using the `ranges` library,
however the check can be configured to work without C++20 by specifying a
function to reverse a range and optionally the header file where that function
lives.

## Options

**Option:** UseCxx20ReverseRanges
When set to true convert loops when in C++20 or later mode using
``std::views::reverse``.
Default value is `true`.
**Option:** MakeReverseRangeFunction
Specify the function used to reverse an iterator pair, the function should
accept a class with ``rbegin`` and ``rend`` methods and return a
class with ``begin`` and ``end`` methods that call the ``rbegin`` and
``rend`` methods respectively. Common examples are ``std::views::reverse``
and ``llvm::reverse``.
Default value is an empty string.
**Option:** MakeReverseRangeHeader
Specifies the header file where :option:`MakeReverseRangeFunction` is
declared. For the previous examples this option would be set to
``range/v3/view/reverse.hpp`` and ``llvm/ADT/STLExtras.h`` respectively.
If this is an empty string and :option:`MakeReverseRangeFunction` is set,
the check will proceed on the assumption that the function is already
available in the translation unit.
This can be wrapped in angle brackets to signify to add the include as a
system include.
Default value is an empty string.
**Option:** IncludeStyle
A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.

## Limitations

There are certain situations where the tool may erroneously perform
transformations that remove information and change semantics. Users of the tool
should be aware of the behavior and limitations of the check outlined by
the cases below.

### Comments inside loop headers

Comments inside the original loop header are ignored and deleted when
transformed.

```c++
for (int i = 0; i < N; /* This will be deleted */ ++i) { }
```

### Range-based loops evaluate end() only once

The C++11 range-based for loop calls `.end()` only once during the
initialization of the loop. If in the original loop `.end()` is called after
each iteration the semantics of the transformed loop may differ.

```c++
// The following is semantically equivalent to the C++11 range-based for loop,
// therefore the semantics of the header will not change.
for (iterator it = container.begin(), e = container.end(); it != e; ++it) { }

// Instead of calling .end() after each iteration, this loop will be
// transformed to call .end() only once during the initialization of the loop,
// which may affect semantics.
for (iterator it = container.begin(); it != container.end(); ++it) { }
```

As explained above, calling member functions of the container in the body
of the loop is considered `risky`. If the called member function modifies the
container the semantics of the converted loop will differ due to `.end()`
being called only once.

```c++
bool flag = false;
for (vector<T>::iterator it = vec.begin(); it != vec.end(); ++it) {
 // Add a copy of the first element to the end of the vector.
 if (!flag) {
 // This line makes this transformation 'risky'.
 vec.push_back(*it);
 flag = true;
 }
 cout << *it;
}
```

The original code above prints out the contents of the container including the
newly added element while the converted loop, shown below, will only print the
original contents and not the newly added element.

```c++
bool flag = false;
for (auto & elem : vec) {
 // Add a copy of the first element to the end of the vector.
 if (!flag) {
 // This line makes this transformation 'risky'
 vec.push_back(elem);
 flag = true;
 }
 cout << elem;
}
```

Semantics will also be affected if `.end()` has side effects. For example, in
the case where calls to `.end()` are logged the semantics will change in the
transformed loop if `.end()` was originally called after each iteration.

```c++
iterator end() {
 num_of_end_calls++;
 return container.end();
}
```

### Overloaded operator->() with side effects

Similarly, if `operator->()` was overloaded to have side effects, such as
logging, the semantics will change. If the iterator's `operator->()` was used
in the original loop it will be replaced with `<container element>.<member>`
instead due to the implicit dereference as part of the range-based for loop.
Therefore any side effect of the overloaded `operator->()` will no longer be
performed.

```c++
for (iterator it = c.begin(); it != c.end(); ++it) {
 it->func(); // Using operator->()
}
// Will be transformed to:
for (auto & elem : c) {
 elem.func(); // No longer using operator->()
}
```

### Pointers and references to containers

While most of the check's risk analysis is dedicated to determining whether
the iterator or container was modified within the loop, it is possible to
circumvent the analysis by accessing and modifying the container through a
pointer or reference.

If the container were directly used instead of using the pointer or reference
the following transformation would have only been applied at the `risky`
level since calling a member function of the container is considered `risky`.
The check cannot identify expressions associated with the container that are
different than the one used in the loop header, therefore the transformation
below ends up being performed at the `safe` level.

```c++
vector<int> vec;

vector<int> *ptr = &vec;
vector<int> &ref = vec;

for (vector<int>::iterator it = vec.begin(), e = vec.end(); it != e; ++it) {
 if (!flag) {
 // Accessing and modifying the container is considered risky, but the risk
 // level is not raised here.
 ptr->push_back(*it);
 ref.push_back(*it);
 flag = true;
 }
}
```

### OpenMP

As range-based for loops are only available since OpenMP 5, this check should
not be used on code with a compatibility requirement of OpenMP prior to
version 5. It is **intentional** that this check does not make any attempts to
exclude incorrect diagnostics on OpenMP for loops prior to OpenMP 5.

To prevent this check to be applied (and to break) OpenMP for loops
but still be applied to non-OpenMP for loops the usage of `NOLINT`
(see `clang-tidy-nolint`) on the specific for loops is recommended.

**Title:** clang-tidy - modernize-macro-to-enum

# modernize-macro-to-enum

Replaces groups of adjacent macros with an unscoped anonymous enum.
Using an unscoped anonymous enum ensures that everywhere the macro
token was used previously, the enumerator name may be safely used.

This check can be used to enforce the C++ core guideline `Enum.1:
Prefer enumerations over macros
<https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#enum1-prefer-enumerations-over-macros>`_,
within the constraints outlined below.

Potential macros for replacement must meet the following constraints:

- Macros must expand only to integral literal tokens or expressions
 of literal tokens. The expression may contain any of the unary
 operators `-`, `+`, `~` or `!`, any of the binary operators
 `,`, `-`, `+`, `*`, `/`, `%`, `&`, `|`, `^`, `<`,
 `>`, `<=`, `>=`, `==`, `!=`, `||`, `&&`, `<<`, `>>`
 or `<=>`, the ternary operator `?:` and its
 [GNU extension](https://gcc.gnu.org/onlinedocs/gcc/Conditionals.html).
 Parenthesized expressions are also recognized. This recognizes
 most valid expressions. In particular, expressions with the
 `sizeof` operator are not recognized.
- Macros must be defined on sequential source file lines, or with
 only comment lines in between macro definitions.
- Macros must all be defined in the same source file.
- Macros must not be defined within a conditional compilation block.
 (Conditional include guards are exempt from this constraint.)
- Macros must not be defined adjacent to other preprocessor directives.
- Macros must not be used in any conditional preprocessing directive.
- Macros must not be used as arguments to other macros.
- Macros must not be undefined.
- Macros must be defined at the top-level, not inside any declaration or
 definition.

Each cluster of macros meeting the above constraints is presumed to
be a set of values suitable for replacement by an anonymous enum.
From there, a developer can give the anonymous enum a name and
continue refactoring to a scoped enum if desired. Comments on the
same line as a macro definition or between subsequent macro definitions
are preserved in the output. No formatting is assumed in the provided
replacements, although clang-tidy can optionally format all fixes.

**Warning:**
Initializing expressions are assumed to be valid initializers for
an enum. C requires that enum values fit into an ``int``, but
this may not be the case for some accepted constant expressions.
For instance ``1 << 40`` will not fit into an ``int`` when the size of
an ``int`` is 32 bits.
Examples:

```c++
#define RED 0xFF0000
#define GREEN 0x00FF00
#define BLUE 0x0000FF

#define TM_NONE (-1) // No method selected.
#define TM_ONE 1 // Use tailored method one.
#define TM_TWO 2 // Use tailored method two. Method two
 // is preferable to method one.
#define TM_THREE 3 // Use tailored method three.
```

becomes

```c++
enum {
RED = 0xFF0000,
GREEN = 0x00FF00,
BLUE = 0x0000FF
};

enum {
TM_NONE = (-1), // No method selected.
TM_ONE = 1, // Use tailored method one.
TM_TWO = 2, // Use tailored method two. Method two
 // is preferable to method one.
TM_THREE = 3 // Use tailored method three.
};
```

**Title:** clang-tidy - modernize-make-shared

# modernize-make-shared

This check finds the creation of `std::shared_ptr` objects by explicitly
calling the constructor and a `new` expression, and replaces it with a call
to `std::make_shared`.

```c++
auto my_ptr = std::shared_ptr<MyPair>(new MyPair(1, 2));

// becomes

auto my_ptr = std::make_shared<MyPair>(1, 2);
```

This check also finds calls to `std::shared_ptr::reset()` with a `new`
expression, and replaces it with a call to `std::make_shared`.

```c++
my_ptr.reset(new MyPair(1, 2));

// becomes

my_ptr = std::make_shared<MyPair>(1, 2);
```

## Options

**Option:** MakeSmartPtrFunction
A string specifying the name of make-shared-ptr function. Default is
`std::make_shared`.
**Option:** MakeSmartPtrFunctionHeader
A string specifying the corresponding header of make-shared-ptr function.
Default is `<memory>`.
**Option:** IncludeStyle
A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.
**Option:** IgnoreMacros
If set to `true`, the check will not give warnings inside macros. Default
is `true`.
**Option:** IgnoreDefaultInitialization
If set to `false`, the check does not suggest edits that will transform
default initialization into value initialization, as this can cause
performance regressions. Default is `true`.

**Title:** clang-tidy - modernize-make-unique

# modernize-make-unique

This check finds the creation of `std::unique_ptr` objects by explicitly
calling the constructor and a `new` expression, and replaces it with a call
to `std::make_unique`, introduced in C++14.

```c++
auto my_ptr = std::unique_ptr<MyPair>(new MyPair(1, 2));

// becomes

auto my_ptr = std::make_unique<MyPair>(1, 2);
```

This check also finds calls to `std::unique_ptr::reset()` with a `new`
expression, and replaces it with a call to `std::make_unique`.

```c++
my_ptr.reset(new MyPair(1, 2));

// becomes

my_ptr = std::make_unique<MyPair>(1, 2);
```

## Options

**Option:** MakeSmartPtrFunction
A string specifying the name of make-unique-ptr function. Default is
`std::make_unique`.
**Option:** MakeSmartPtrFunctionHeader
A string specifying the corresponding header of make-unique-ptr function.
Default is `<memory>`.
**Option:** IncludeStyle
A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.
**Option:** IgnoreMacros
If set to `true`, the check will not give warnings inside macros. Default
is `true`.
**Option:** IgnoreDefaultInitialization
If set to `false`, the check does not suggest edits that will transform
default initialization into value initialization, as this can cause
performance regressions. Default is `true`.

**Title:** clang-tidy - modernize-min-max-use-initializer-list

# modernize-min-max-use-initializer-list

Replaces nested `std::min` and `std::max` calls with an initializer list
where applicable.

For instance, consider the following code:

```cpp
int a = std::max(std::max(i, j), k);
```

The check will transform the above code to:

```cpp
int a = std::max({i, j, k});
```

# Performance Considerations

While this check simplifies the code and makes it more readable, it may cause
performance degradation for non-trivial types due to the need to copy objects
into the initializer list.

To avoid this, it is recommended to use `std::ref` or `std::cref` for
non-trivial types:

```cpp
std::string b = std::max({std::ref(i), std::ref(j), std::ref(k)});
```

# Options

**Option:** IncludeStyle
A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.
**Option:** IgnoreNonTrivialTypes
A boolean specifying whether to ignore non-trivial types. Default is `true`.
**Option:** IgnoreTrivialTypesOfSizeAbove
An integer specifying the size (in bytes) above which trivial types are
ignored. Default is `32`.

**Title:** clang-tidy - modernize-pass-by-value

# modernize-pass-by-value

With move semantics added to the language and the standard library updated with
move constructors added for many types it is now interesting to take an
argument directly by value, instead of by const-reference, and then copy. This
check allows the compiler to take care of choosing the best way to construct
the copy.

The transformation is usually beneficial when the calling code passes an
*rvalue* and assumes the move construction is a cheap operation. This short
example illustrates how the construction of the value happens:

```c++
void foo(std::string s);
std::string get_str();

void f(const std::string &str) {
 foo(str); // lvalue -> copy construction
 foo(get_str()); // prvalue -> move construction
}
```

**Note:**
Currently, only constructors are transformed to make use of pass-by-value.
Contributions that handle other situations are welcome!

## Pass-by-value in constructors

Replaces the uses of const-references constructor parameters that are copied
into class fields. The parameter is then moved with `std::move()`.

Since `std::move()` is a library function declared in `<utility>` it may be
necessary to add this include. The check will add the include directive when
necessary.

```c++
 #include <string>

 class Foo {
 public:
- Foo(const std::string &Copied, const std::string &ReadOnly)
- : Copied(Copied), ReadOnly(ReadOnly)
+ Foo(std::string Copied, const std::string &ReadOnly)
+ : Copied(std::move(Copied)), ReadOnly(ReadOnly)
 {}

 private:
 std::string Copied;
 const std::string &ReadOnly;
 };

 std::string get_cwd();

 void f(const std::string &Path) {
 // The parameter corresponding to 'get_cwd()' is move-constructed. By
 // using pass-by-value in the Foo constructor we managed to avoid a
 // copy-construction.
 Foo foo(get_cwd(), Path);
 }
```

If the parameter is used more than once no transformation is performed since
moved objects have an undefined state. It means the following code will be left
untouched:

```c++
#include <string>

void pass(const std::string &S);

struct Foo {
 Foo(const std::string &S) : Str(S) {
 pass(S);
 }

 std::string Str;
};
```

## Limitations

A situation where the generated code can be wrong is when the object referenced
is modified before the assignment in the init-list through a "hidden" reference.

Example:

```c++
 std::string s("foo");

 struct Base {
 Base() {
 s = "bar";
 }
 };

 struct Derived : Base {
- Derived(const std::string &S) : Field(S)
+ Derived(std::string S) : Field(std::move(S))
 { }

 std::string Field;
 };

 void f() {
- Derived d(s); // d.Field holds "bar"
+ Derived d(s); // d.Field holds "foo"
 }
```

### Note about delayed template parsing

When delayed template parsing is enabled, constructors part of templated
contexts; templated constructors, constructors in class templates, constructors
of inner classes of template classes, etc., are not transformed. Delayed
template parsing is enabled by default on Windows as a Microsoft extension:
Clang Compiler User's Manual - Microsoft extensions.

Delayed template parsing can be enabled using the `-fdelayed-template-parsing`
flag and disabled using `-fno-delayed-template-parsing`.

Example:

```c++
 template <typename T> class C {
 std::string S;

 public:
= // using -fdelayed-template-parsing (default on Windows)
= C(const std::string &S) : S(S) {}

+ // using -fno-delayed-template-parsing (default on non-Windows systems)
+ C(std::string S) : S(std::move(S)) {}
 };
```

**Seealso:**
For more information about the pass-by-value idiom, read: `Want Speed? Pass by Value`_.

.. _Want Speed? Pass by Value: https://web.archive.org/web/20140205194657/http://cpp-next.com/archive/2009/08/want-speed-pass-by-value/

## Options

**Option:** IncludeStyle
A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.
**Option:** ValuesOnly
When `true`, the check only warns about copied parameters that are already
passed by value. Default is `false`.
**Option:** IgnoreMacros
When `true`, the check will not give warnings inside macros. Default is
`false`.

**Title:** clang-tidy - modernize-raw-string-literal

# modernize-raw-string-literal

This check selectively replaces string literals containing escaped characters
with raw string literals.

Example:

```c++
const char *const Quotes{"embedded \"quotes\""};
const char *const Paragraph{"Line one.\nLine two.\nLine three.\n"};
const char *const SingleLine{"Single line.\n"};
const char *const TrailingSpace{"Look here -> \n"};
const char *const Tab{"One\tTwo\n"};
const char *const Bell{"Hello!\a And welcome!"};
const char *const Path{"C:\\Program Files\\Vendor\\Application.exe"};
const char *const RegEx{"\\w\\([a-z]\\)"};
```

becomes

```c++
const char *const Quotes{R"(embedded "quotes")"};
const char *const Paragraph{"Line one.\nLine two.\nLine three.\n"};
const char *const SingleLine{"Single line.\n"};
const char *const TrailingSpace{"Look here -> \n"};
const char *const Tab{"One\tTwo\n"};
const char *const Bell{"Hello!\a And welcome!"};
const char *const Path{R"(C:\Program Files\Vendor\Application.exe)"};
const char *const RegEx{R"(\w\([a-z]\))"};
```

The presence of any of the following escapes can cause the string to be
converted to a raw string literal: `\\`, `\'`, `\"`, `\?`,
and octal or hexadecimal escapes for printable ASCII characters.

A string literal containing only escaped newlines is a common way of
writing lines of text output. Introducing physical newlines with raw
string literals in this case is likely to impede readability. These
string literals are left unchanged.

An escaped horizontal tab, form feed, or vertical tab prevents the string
literal from being converted. The presence of a horizontal tab, form feed or
vertical tab in source code is not visually obvious.

## Options

**Option:** DelimiterStem
Custom delimiter to escape characters in raw string literals. It is used in
the following construction: ``R"stem_delimiter(contents)stem_delimiter"``.
The default value is `lit`.
**Option:** ReplaceShorterLiterals
Controls replacing shorter non-raw string literals with longer raw string
literals. Setting this option to `true` enables the replacement.
The default value is `false` (shorter literals are not replaced).

```{title} clang-tidy - modernize-redundant-void-arg
```

# modernize-redundant-void-arg

Finds and removes redundant `void` argument lists.
Works in C++ and in C23 and up.

Examples:

| Initial code | Code with applied fixes |
| --------------------------------- | ------------------------- |
| `int f(void);` | `int f();` |
| `int (*f(void))(void);` | `int (*f())();` |
| `typedef int (*f_t(void))(void);` | `typedef int (*f_t())();` |
| `void (C::*p)(void);` | `void (C::*p)();` |
| `C::C(void) {}` | `C::C() {}` |
| `C::~C(void) {}` | `C::~C() {}` |

**Title:** clang-tidy - modernize-replace-auto-ptr

# modernize-replace-auto-ptr

This check replaces the uses of the deprecated class `std::auto_ptr` by
`std::unique_ptr` (introduced in C++11). The transfer of ownership, done
by the copy-constructor and the assignment operator, is changed to match
`std::unique_ptr` usage by using explicit calls to `std::move()`.

Migration example:

```c++
-void take_ownership_fn(std::auto_ptr<int> int_ptr);
+void take_ownership_fn(std::unique_ptr<int> int_ptr);

 void f(int x) {
- std::auto_ptr<int> a(new int(x));
- std::auto_ptr<int> b;
+ std::unique_ptr<int> a(new int(x));
+ std::unique_ptr<int> b;

- b = a;
- take_ownership_fn(b);
+ b = std::move(a);
+ take_ownership_fn(std::move(b));
 }
```

Since `std::move()` is a library function declared in `<utility>` it may be
necessary to add this include. The check will add the include directive when
necessary.

## Limitations

* If headers modification is not activated or if a header is not allowed to be
 changed this check will produce broken code (compilation error), where the
 headers' code will stay unchanged while the code using them will be changed.

* Client code that declares a reference to an `std::auto_ptr` coming from
 code that can't be migrated (such as a header coming from a 3\ `rd`
 party library) will produce a compilation error after migration. This is
 because the type of the reference will be changed to `std::unique_ptr` but
 the type returned by the library won't change, binding a reference to
 `std::unique_ptr` from an `std::auto_ptr`. This pattern doesn't make much
 sense and usually `std::auto_ptr` are stored by value (otherwise what is
 the point in using them instead of a reference or a pointer?).

```c++
 // <3rd-party header...>
 std::auto_ptr<int> get_value();
 const std::auto_ptr<int> & get_ref();

 // <calling code (with migration)...>
-std::auto_ptr<int> a(get_value());
+std::unique_ptr<int> a(get_value()); // ok, unique_ptr constructed from auto_ptr

-const std::auto_ptr<int> & p = get_ptr();
+const std::unique_ptr<int> & p = get_ptr(); // won't compile
```

* Non-instantiated templates aren't modified.

```c++
template <typename X>
void f() {
 std::auto_ptr<X> p;
}

// only 'f<int>()' (or similar) will trigger the replacement.
```

## Options

**Option:** IncludeStyle
A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.

**Title:** clang-tidy - modernize-replace-disallow-copy-and-assign-macro

# modernize-replace-disallow-copy-and-assign-macro

Finds macro expansions of `DISALLOW_COPY_AND_ASSIGN(Type)` and replaces them
with a deleted copy constructor and a deleted assignment operator.

Before the `delete` keyword was introduced in C++11 it was common practice to
declare a copy constructor and an assignment operator as private members. This
effectively makes them unusable to the public API of a class.

With the advent of the `delete` keyword in C++11 we can abandon the
`private` access of the copy constructor and the assignment operator and
delete the methods entirely.

When running this check on a code like this:

```c++
class Foo {
private:
 DISALLOW_COPY_AND_ASSIGN(Foo);
};
```

It will be transformed to this:

```c++
class Foo {
private:
 Foo(const Foo &) = delete;
 const Foo &operator=(const Foo &) = delete;
};
```

## Limitations

* Notice that the migration example above leaves the `private` access
 specification untouched. You might want to run the check
 modernize-use-equals-delete
 to get warnings for deleted functions in private sections.

## Options

**Option:** MacroName
A string specifying the macro name whose expansion will be replaced.
Default is `DISALLOW_COPY_AND_ASSIGN`.
See: https://en.cppreference.com/w/cpp/language/function#Deleted_functions

**Title:** clang-tidy - modernize-replace-random-shuffle

# modernize-replace-random-shuffle

This check will find occurrences of `std::random_shuffle` and replace it with
`std::shuffle`. In C++17 `std::random_shuffle` will no longer be available
and thus we need to replace it.

Below are two examples of what kind of occurrences will be found and two
examples of what it will be replaced with.

```c++
std::vector<int> v;

// First example
std::random_shuffle(vec.begin(), vec.end());

// Second example
std::random_shuffle(vec.begin(), vec.end(), randomFunc);
```

Both of these examples will be replaced with:

```c++
std::shuffle(vec.begin(), vec.end(), std::mt19937(std::random_device()()));
```

The second example will also receive a warning that `randomFunc` is no longer
supported in the same way as before so if the user wants the same
functionality, the user will need to change the implementation of the
`randomFunc`.

One thing to be aware of here is that `std::random_device` is quite expensive
to initialize. So if you are using the code in a performance critical place,
you probably want to initialize it elsewhere.

Another thing is that the seeding quality of the suggested fix is quite poor:
`std::mt19937` has an internal state of 624 32-bit integers, but is only
seeded with a single integer. So if you require
higher quality randomness, you should consider seeding better, for example:

```c++
std::shuffle(v.begin(), v.end(), []() {
 std::mt19937::result_type seeds[std::mt19937::state_size];
 std::random_device device;
 std::uniform_int_distribution<typename std::mt19937::result_type> dist;
 std::generate(std::begin(seeds), std::end(seeds), [&] { return dist(device); });
 std::seed_seq seq(std::begin(seeds), std::end(seeds));
 return std::mt19937(seq);
}());
```

```{title} clang-tidy - modernize-return-braced-init-list
```

# modernize-return-braced-init-list

Replaces explicit calls to the constructor in a return with a braced
initializer list. This way the return type is not needlessly duplicated in the
function definition and the return statement.

```cpp
Foo bar() {
 Baz baz;
 return Foo(baz);
}

// transforms to:

Foo bar() {
 Baz baz;
 return {baz};
}
```

The check is not applied when the constructed type has a
`std::initializer_list` constructor, since list-initialization would prefer
that constructor and the braced form could therefore select a different
constructor than the original call.

```{title} clang-tidy - modernize-shrink-to-fit
```

# modernize-shrink-to-fit

Replace copy and swap tricks on shrinkable containers with the
`shrink_to_fit()` method call.

The `shrink_to_fit()` method is more readable and more effective than
the copy and swap trick to reduce the capacity of a shrinkable container.
Note that, the `shrink_to_fit()` method is only available in C++11 and up.

**Title:** clang-tidy - modernize-type-traits

# modernize-type-traits

Converts standard library type traits of the form `traits<...>::type` and
`traits<...>::value` into `traits_t<...>` and
`traits_v<...>` respectively.

Also suggests converting `std::remove_cv_t<std::remove_reference_t<...>` into
`std::remove_cvref_t<...>` when targeting C++20 or above.

For example:

```c++
std::is_integral<T>::value
std::is_same<int, float>::value
typename std::add_const<T>::type
std::make_signed<unsigned>::type

std::remove_cv_t<std::remove_reference_t<int>>
```

Would be converted into:

```c++
std::is_integral_v<T>
std::is_same_v<int, float>
std::add_const_t<T>
std::make_signed_t<unsigned>

std::remove_cvref_t<int>
```

## Options

**Option:** IgnoreMacros
If `true` don't diagnose traits defined in macros.

Note: Fixes will never be emitted for code inside of macros.

.. code-block:: c++

 #define IS_SIGNED(T) std::is_signed<T>::value

Defaults to `false`.

## Limitations

Does not currently diagnose uses of type traits with nested name
specifiers (e.g. `std::chrono::is_clock`,
`std::chrono::treat_as_floating_point`).

```{title} clang-tidy - modernize-unary-static-assert
```

# modernize-unary-static-assert

The check diagnoses any `static_assert` declaration with an
empty string literal and provides a fix-it note to replace the
declaration with a single-argument `static_assert` declaration.

The check is only applicable for C++17 and later code.

The following code:

```cpp
void f_textless(int a) {
 static_assert(sizeof(a) <= 10, "");
}
```

is replaced by:

```cpp
void f_textless(int a) {
 static_assert(sizeof(a) <= 10);
}
```

**Title:** clang-tidy - modernize-use-auto

# modernize-use-auto

This check is responsible for using the `auto` type specifier for variable
declarations to *improve code readability and maintainability*. For example:

```c++
std::vector<int>::iterator I = my_container.begin();

// transforms to:

auto I = my_container.begin();
```

The `auto` type specifier will only be introduced in situations where the
variable type matches the type of the initializer expression. In other words
`auto` should deduce the same type that was originally spelled in the source.
However, not every situation should be transformed:

```c++
int val = 42;
InfoStruct &I = SomeObject.getInfo();

// Should not become:

auto val = 42;
auto &I = SomeObject.getInfo();
```

In this example using `auto` for builtins doesn't improve readability. In
other situations it makes the code less self-documenting impairing readability
and maintainability. As a result, `auto` is used only introduced in specific
situations described below.

## Iterators

Iterator type specifiers tend to be long and used frequently, especially in
loop constructs. Since the functions generating iterators have a common format,
the type specifier can be replaced without obscuring the meaning of code while
improving readability and maintainability.

```c++
for (std::vector<int>::iterator I = my_container.begin(),
 E = my_container.end();
 I != E; ++I) {
}

// becomes

for (auto I = my_container.begin(), E = my_container.end(); I != E; ++I) {
}
```

The check will only replace iterator type-specifiers when all of the following
conditions are satisfied:

* The iterator is for one of the standard containers in `std` namespace:

 * `array`
 * `deque`
 * `forward_list`
 * `list`
 * `vector`
 * `map`
 * `multimap`
 * `set`
 * `multiset`
 * `unordered_map`
 * `unordered_multimap`
 * `unordered_set`
 * `unordered_multiset`
 * `queue`
 * `priority_queue`
 * `stack`

* The iterator is one of the possible iterator types for standard containers:

 * `iterator`
 * `reverse_iterator`
 * `const_iterator`
 * `const_reverse_iterator`

* In addition to using iterator types directly, typedefs or other ways of
 referring to those types are also allowed. However, implementation-specific
 types for which a type like `std::vector<int>::iterator` is itself a
 typedef will not be transformed. Consider the following examples:

```c++
// The following direct uses of iterator types will be transformed.
std::vector<int>::iterator I = MyVec.begin();
{
 using namespace std;
 list<int>::iterator I = MyList.begin();
}

// The type specifier for J would transform to auto since it's a typedef
// to a standard iterator type.
typedef std::map<int, std::string>::const_iterator map_iterator;
map_iterator J = MyMap.begin();

// The following implementation-specific iterator type for which
// std::vector<int>::iterator could be a typedef would not be transformed.
__gnu_cxx::__normal_iterator<int*, std::vector> K = MyVec.begin();
```

* The initializer for the variable being declared is not a braced initializer
 list. Otherwise, use of `auto` would cause the type of the variable to be
 deduced as `std::initializer_list`.

## New expressions

Frequently, when a pointer is declared and initialized with `new`, the
pointee type is written twice: in the declaration type and in the
`new` expression. In this case, the declaration type can be replaced with
`auto` improving readability and maintainability.

```c++
TypeName *my_pointer = new TypeName(my_param);

// becomes

auto *my_pointer = new TypeName(my_param);
```

The check will also replace the declaration type in multiple declarations, if
the following conditions are satisfied:

* All declared variables have the same type (i.e. all of them are pointers to
 the same type).
* All declared variables are initialized with a `new` expression.
* The types of all the new expressions are the same than the pointee of the
 declaration type.

```c++
TypeName *my_first_pointer = new TypeName, *my_second_pointer = new TypeName;

// becomes

auto *my_first_pointer = new TypeName, *my_second_pointer = new TypeName;
```

## Cast expressions

Frequently, when a variable is declared and initialized with a cast, the
variable type is written twice: in the declaration type and in the
cast expression. In this case, the declaration type can be replaced with
`auto` improving readability and maintainability.

```c++
TypeName *my_pointer = static_cast<TypeName>(my_param);

// becomes

auto *my_pointer = static_cast<TypeName>(my_param);
```

The check handles `static_cast`, `dynamic_cast`, `const_cast`,
`reinterpret_cast`, functional casts, C-style casts and function templates
that behave as casts, such as `llvm::dyn_cast`, `boost::lexical_cast` and
`gsl::narrow_cast`. Calls to function templates are considered to behave as
casts if the first template argument is explicit and is a type, and the
function returns that type, or a pointer or reference to it.

## Limitations

* If the initializer is an explicit conversion constructor, the check will not
 replace the type specifier even though it would be safe to do so.

* User-defined iterators are not handled at this time.

## Options

**Option:** MinTypeNameLength
If the option is set to non-zero (default `5`), the check will ignore type
names having a length less than the option value. The option affects
expressions only, not iterators.
Spaces between multi-lexeme type names (``long int``) are considered as one.
If the :option:`RemoveStars` option (see below) is set to `true`, then ``*s``
in the type are also counted as a part of the type name.

```c++
// MinTypeNameLength = 0, RemoveStars=0

int a = static_cast<int>(foo()); // ---> auto a = ...
// length(bool *) = 4
bool *b = new bool; // ---> auto *b = ...
unsigned c = static_cast<unsigned>(foo()); // ---> auto c = ...

// MinTypeNameLength = 5, RemoveStars=0

int a = static_cast<int>(foo()); // ---> int a = ...
bool b = static_cast<bool>(foo()); // ---> bool b = ...
bool *pb = static_cast<bool*>(foo()); // ---> bool *pb = ...
unsigned c = static_cast<unsigned>(foo()); // ---> auto c = ...
// length(long <on-or-more-spaces> int) = 8
long int d = static_cast<long int>(foo()); // ---> auto d = ...

// MinTypeNameLength = 5, RemoveStars=1

int a = static_cast<int>(foo()); // ---> int a = ...
// length(int * * ) = 5
int **pa = static_cast<int**>(foo()); // ---> auto pa = ...
bool b = static_cast<bool>(foo()); // ---> bool b = ...
bool *pb = static_cast<bool*>(foo()); // ---> auto pb = ...
unsigned c = static_cast<unsigned>(foo()); // ---> auto c = ...
long int d = static_cast<long int>(foo()); // ---> auto d = ...
```

**Option:** RemoveStars
If the option is set to `true` (default is `false`), the check will remove
stars from the non-typedef pointer types when replacing type names with
``auto``. Otherwise, the check will leave stars. For example:

```c++
TypeName *my_first_pointer = new TypeName, *my_second_pointer = new TypeName;

// RemoveStars = 0

auto *my_first_pointer = new TypeName, *my_second_pointer = new TypeName;

// RemoveStars = 1

auto my_first_pointer = new TypeName, my_second_pointer = new TypeName;
```

```{title} clang-tidy - modernize-use-bool-literals
```

# modernize-use-bool-literals

Finds integer literals which are cast to `bool`.

```cpp
bool p = 1;
bool f = static_cast<bool>(1);
std::ios_base::sync_with_stdio(0);
bool x = p ? 1 : 0;

// transforms to

bool p = true;
bool f = true;
std::ios_base::sync_with_stdio(false);
bool x = p ? true : false;
```

## Options

```{option} IgnoreMacros

If set to `true`, the check will not give warnings inside macros. Default
is `true`.
```

**Title:** clang-tidy - modernize-use-constraints

# modernize-use-constraints

Replace `std::enable_if` with C++20 requires clauses.

`std::enable_if` is a SFINAE mechanism for selecting the desired function or
class template based on type traits or other requirements. `enable_if`
changes the meta-arity of the template, and has other
`adverse side effects
<https://open-std.org/JTC1/SC22/WG21/docs/papers/2016/p0225r0.html>`_
in the code. C++20 introduces concepts and constraints as a cleaner language
provided solution to achieve the same outcome.

This check finds some common `std::enable_if` patterns that can be replaced
by C++20 requires clauses. The tool can replace some of these patterns
automatically, otherwise, the tool will emit a diagnostic without a
replacement. The tool can detect the following `std::enable_if` patterns

1. `std::enable_if` in the return type of a function
2. `std::enable_if` as the trailing template parameter for function templates

Other uses, for example, in class templates for function parameters, are not
currently supported by this tool. Other variants such as `boost::enable_if`
are not currently supported by this tool.

Below are some examples of code using `std::enable_if`.

```c++
// enable_if in function return type
template <typename T>
std::enable_if_t<T::some_trait, int> only_if_t_has_the_trait() { ... }

// enable_if in the trailing template parameter
template <typename T, std::enable_if_t<T::some_trait, int> = 0>
void another_version() { ... }

template <typename T>
typename std::enable_if<T::some_value, Obj>::type existing_constraint() requires (T::another_value) {
 return Obj{};
}

template <typename T, std::enable_if_t<T::some_trait, int> = 0>
struct my_class {};
```

The tool will replace the above code with,

```c++
// warning: use C++20 requires constraints instead of enable_if [modernize-use-constraints]
template <typename T>
int only_if_t_has_the_trait() requires T::some_trait { ... }

// warning: use C++20 requires constraints instead of enable_if [modernize-use-constraints]
template <typename T>
void another_version() requires T::some_trait { ... }

// The tool will emit a diagnostic for the following, but will
// not attempt to replace the code.
// warning: use C++20 requires constraints instead of enable_if [modernize-use-constraints]
template <typename T>
typename std::enable_if<T::some_value, Obj>::type existing_constraint() requires (T::another_value) {
 return Obj{};
}

// The tool will not emit a diagnostic or attempt to replace the code.
template <typename T, std::enable_if_t<T::some_trait, int> = 0>
struct my_class {};
```

**Note:**
System headers are not analyzed by this check.

**Title:** clang-tidy - modernize-use-default-member-init

# modernize-use-default-member-init

This check converts constructors' member initializers into the new
default member initializers in C++11. Other member initializers that match the
default member initializer are removed. This can reduce repeated code or allow
use of '= default'.

```c++
struct A {
 A() : i(5), j(10.0) {}
 A(int i) : i(i), j(10.0) {}
 int i;
 double j;
};

// becomes

struct A {
 A() {}
 A(int i) : i(i) {}
 int i{5};
 double j{10.0};
};
```

**Note:**
Only converts member initializers for built-in types, enums, and pointers.
The `readability-redundant-member-init` check will remove redundant member
initializers for classes.

## Options

**Option:** UseAssignment
If this option is set to `true` (default is `false`), the check will initialize
members with an assignment. For example:

```c++
struct A {
 A() {}
 A(int i) : i(i) {}
 int i = 5;
 double j = 10.0;
};
```

**Option:** IgnoreMacros
If this option is set to `true` (default is `true`), the check will not warn
about members declared inside macros.
**Option:** IgnoreNonVisibleReferences
If this option is set to `true` (default is `true`), the check will not warn
about member initializers that refer to declarations that would not be
visible from the default member initializer created by the fix-it. If set to
`false`, the check will warn about these cases without emitting fix-its.

---
orphan: true
---

```{title} clang-tidy - modernize-use-default
```

# modernize-use-default

This check has been renamed to
{doc}`modernize-use-equals-default <../modernize/use-equals-default>`.

**Title:** clang-tidy - modernize-use-designated-initializers

# modernize-use-designated-initializers

Finds initializer lists for aggregate types which could be written as
designated initializers instead.

With plain initializer lists, it is very easy to introduce bugs when adding new
fields in the middle of a struct or class type. The same confusion might arise
when changing the order of fields.

C++20 supports the designated initializer syntax for aggregate types. By
applying it, we can always be sure that aggregates are constructed correctly,
because every variable being initialized is referenced by its name.

Example:

```
struct S { int i, j; };
```

is an aggregate type that should be initialized as

```
S s{.i = 1, .j = 2};
```

instead of

```
S s{1, 2};
```

which could easily become an issue when `i` and `j` are swapped in the
declaration of `S`.

Even when compiling in a language version older than C++20, depending on your
compiler, designated initializers are potentially supported. Therefore, the
check is by default restricted to C99/C++20 and above. Check out the options
`-Wc99-designator` to get support for mixed designators in initializer list
in C and `-Wc++20-designator` for support of designated initializers in older
C++ language modes.

## Options

**Option:** IgnoreMacros
The value `false` specifies that components of initializer lists expanded from
macros are not checked. The default value is `true`.
**Option:** IgnoreSingleElementAggregates
The value `false` specifies that even initializers for aggregate types with
only a single element should be checked. The default value is `true`.
``std::array`` initializations are always excluded, as the type is a
standard library abstraction and not intended to be initialized with
designated initializers.
**Option:** RestrictToPODTypes
The value `true` specifies that only Plain Old Data (POD) types shall be
checked. This makes the check applicable to even older C++ standards. The
default value is `false`.
**Option:** StrictCStandardCompliance
When set to `false`, the check will not restrict itself to C99 and above.
The default value is `true`.
**Option:** StrictCppStandardCompliance
When set to `false`, the check will not restrict itself to C++20 and above.
The default value is `true`.

**Title:** clang-tidy - modernize-use-emplace

# modernize-use-emplace

The check flags insertions to an STL-style container done by calling the
`push_back`, `push`, or `push_front` methods with an
explicitly-constructed temporary of the container element type. In this case,
the corresponding `emplace` equivalent methods result in less verbose and
potentially more efficient code. Right now the check doesn't support
`insert`. It also doesn't support `insert` functions for associative
containers because replacing `insert` with `emplace` may result in
[speed regression](https://htmlpreview.github.io/?https://github.com/HowardHinnant/papers/blob/master/insert_vs_emplace.html), but it might get support with some addition flag in the future.

The `ContainersWithPushBack`, `ContainersWithPush`, and
`ContainersWithPushFront` options are used to specify the container
types that support the `push_back`, `push`, and `push_front` operations
respectively. The default values for these options are as follows:

* `ContainersWithPushBack`: `std::vector`, `std::deque`,
 and `std::list`.
* `ContainersWithPush`: `std::stack`, `std::queue`,
 and `std::priority_queue`.
* `ContainersWithPushFront`: `std::forward_list`,
 `std::list`, and `std::deque`.

This check also reports when an `emplace`-like method is improperly used,
for example using `emplace_back` while also calling a constructor. This
creates a temporary that requires at best a move and at worst a copy. Almost
all `emplace`-like functions in the STL are covered by this, with
`try_emplace` on `std::map` and `std::unordered_map` being the
exception as it behaves slightly differently than all the others. More
containers can be added with the `EmplacyFunctions` option, so long
as the container defines a `value_type` type, and the `emplace`-like
functions construct a `value_type` object.

Before:

```c++
std::vector<MyClass> v;
v.push_back(MyClass(21, 37));
v.emplace_back(MyClass(21, 37));

std::vector<std::pair<int, int>> w;

w.push_back(std::pair<int, int>(21, 37));
w.push_back(std::make_pair(21L, 37L));
w.emplace_back(std::make_pair(21L, 37L));
```

After:

```c++
std::vector<MyClass> v;
v.emplace_back(21, 37);
v.emplace_back(21, 37);

std::vector<std::pair<int, int>> w;
w.emplace_back(21, 37);
w.emplace_back(21L, 37L);
w.emplace_back(21L, 37L);
```

By default, the check is able to remove unnecessary `std::make_pair` and
`std::make_tuple` calls from `push_back` calls on containers of
`std::pair` and `std::tuple`. Custom tuple-like types can be modified by
the `TupleTypes` option; custom make functions can be modified by the
`TupleMakeFunctions` option.

The other situation is when we pass arguments that will be converted to a type
inside a container.

Before:

```c++
std::vector<boost::optional<std::string> > v;
v.push_back("abc");
```

After:

```c++
std::vector<boost::optional<std::string> > v;
v.emplace_back("abc");
```

In some cases the transformation would be valid, but the code wouldn't be
exception safe. In this case the calls of `push_back` won't be replaced.

```c++
std::vector<std::unique_ptr<int>> v;
v.push_back(std::unique_ptr<int>(new int(0)));
auto *ptr = new int(1);
v.push_back(std::unique_ptr<int>(ptr));
```

This is because replacing it with `emplace_back` could cause a leak of this
pointer if `emplace_back` would throw exception before emplacement (e.g. not
enough memory to add a new element).

For more info read item 42 - "Consider emplacement instead of insertion." of
Scott Meyers "Effective Modern C++".

The default smart pointers that are considered are `std::unique_ptr`,
`std::shared_ptr`, `std::auto_ptr`. To specify other smart pointers or
other classes use the `SmartPointers` option.

Check also doesn't fire if any argument of the constructor call would be:

 - a bit-field (bit-fields can't bind to rvalue/universal reference)

 - a `new` expression (to avoid leak)

 - if the argument would be converted via derived-to-base cast.

This check requires C++11 or higher to run.

## Options

**Option:** ContainersWithPushBack
Semicolon-separated list of class names of custom containers that support
``push_back``.
**Option:** ContainersWithPush
Semicolon-separated list of class names of custom containers that support
``push``.
**Option:** ContainersWithPushFront
Semicolon-separated list of class names of custom containers that support
``push_front``.
**Option:** IgnoreImplicitConstructors
When `true`, the check will ignore implicitly constructed arguments of
``push_back``, e.g.

.. code-block:: c++

 std::vector<std::string> v;
 v.push_back("a"); // Ignored when IgnoreImplicitConstructors is `true`.

Default is `false`.
**Option:** SmartPointers
Semicolon-separated list of class names of custom smart pointers.
**Option:** TupleTypes
Semicolon-separated list of ``std::tuple``-like class names.
**Option:** TupleMakeFunctions
Semicolon-separated list of ``std::make_tuple``-like function names. Those
function calls will be removed from ``push_back`` calls and turned into
``emplace_back``.
**Option:** EmplacyFunctions
Semicolon-separated list of containers without their template parameters
and some ``emplace``-like method of the container. Example:
``vector::emplace_back``. Those methods will be checked for improper use and
the check will report when a temporary is unnecessarily created. All STL
containers with such member functions are supported by default.

### Example

```c++
std::vector<MyTuple<int, bool, char>> x;
x.push_back(MakeMyTuple(1, false, 'x'));
x.emplace_back(MakeMyTuple(1, false, 'x'));
```

transforms to:

```c++
std::vector<MyTuple<int, bool, char>> x;
x.emplace_back(1, false, 'x');
x.emplace_back(1, false, 'x');
```

when `TupleTypes` is set to `MyTuple`, `TupleMakeFunctions`
is set to `MakeMyTuple`, and `EmplacyFunctions` is set to
`vector::emplace_back`.

**Title:** clang-tidy - modernize-use-equals-default

# modernize-use-equals-default

This check replaces default bodies of special member functions with ``=
default;``. The explicitly defaulted function declarations enable more
opportunities in optimization, because the compiler might treat explicitly
defaulted functions as trivial.

```c++
struct A {
 A() {}
 ~A();
};
A::~A() {}

// becomes

struct A {
 A() = default;
 ~A();
};
A::~A() = default;
```

**Note:**
Move-constructor and move-assignment operator are not supported yet.

## Options

**Option:** IgnoreMacros
If set to `true`, the check will not give warnings inside macros and will
ignore special members with bodies contain macros or preprocessor directives.
Default is `true`.

**Title:** clang-tidy - modernize-use-equals-delete

# modernize-use-equals-delete

Identifies unimplemented private special member functions, and recommends using
`= delete` for them. Additionally, it recommends relocating any deleted
member function from the `private` to the `public` section.

Before the introduction of C++11, the primary method to effectively "erase" a
particular function involved declaring it as `private` without providing a
definition. This approach would result in either a compiler error (when
attempting to call a private function) or a linker error (due to an undefined
reference).

However, subsequent to the advent of C++11, a more conventional approach
emerged for achieving this purpose. It involves flagging functions as
`= delete` and keeping them in the `public` section of the class.

To prevent false positives, this check is only active within a translation
unit where all other member functions have been implemented. The check will
generate partial fixes by introducing `= delete`, but the user is responsible
for manually relocating functions to the `public` section.

```c++
// Example: bad
class A {
 private:
 A(const A&);
 A& operator=(const A&);
};

// Example: good
class A {
 public:
 A(const A&) = delete;
 A& operator=(const A&) = delete;
};
```

**Option:** IgnoreMacros
If this option is set to `true` (default is `true`), the check will not warn
about functions declared inside macros.

**Title:** clang-tidy - modernize-use-integer-sign-comparison

# modernize-use-integer-sign-comparison

Replace comparisons between signed and unsigned integers with their safe
C++20 `std::cmp_*` alternative, if available.

The check provides a replacement only for C++20 or later, otherwise
it highlights the problem and expects the user to fix it manually.

Examples of fixes created by the check:

```c++
unsigned int func(int a, unsigned int b) {
 return a == b;
}
```

becomes

```c++
#include <utility>

unsigned int func(int a, unsigned int b) {
 return std::cmp_equal(a, b);
}
```

## Options

**Option:** IncludeStyle
A string specifying which include-style is used, `llvm` or `google`.
Default is `llvm`.
**Option:** EnableQtSupport
Makes C++17 ``q20::cmp_*`` alternative available for Qt-based
applications. Default is `false`.

**Title:** clang-tidy - modernize-use-nodiscard

# modernize-use-nodiscard

Adds `[[nodiscard]]` attributes (introduced in C++17) to member functions in
order to highlight at compile time which return values should not be ignored.

Member functions need to satisfy the following conditions to be considered by
this check:

 - no `[[nodiscard]]`, `[[noreturn]]`,
 `__attribute__((warn_unused_result))`,
 `[[clang::warn_unused_result]]` nor `[[gcc::warn_unused_result]]`
 attribute,
 - non-void return type,
 - non-template return types,
 - const member function,
 - non-variadic functions,
 - no non-const reference parameters,
 - no pointer parameters,
 - no template parameters,
 - no template function parameters,
 - not be a member of a class with mutable member variables,
 - no Lambdas,
 - no conversion functions.

Such functions have no means of altering any state or passing values other than
via the return type. Unless the member functions are altering state via some
external call (e.g. I/O).

## Example

```c++
bool empty() const;
bool empty(int i) const;
```

transforms to:

```c++
[[nodiscard]] bool empty() const;
[[nodiscard]] bool empty(int i) const;
```

## Options

**Option:** ReplacementString
Specifies a macro to use instead of ``[[nodiscard]]``. This is useful when
maintaining source code that needs to compile with a pre-C++17 compiler.

### Example

```c++
bool empty() const;
bool empty(int i) const;
```

transforms to:

```c++
NO_DISCARD bool empty() const;
NO_DISCARD bool empty(int i) const;
```

if the `ReplacementString` option is set to `NO_DISCARD`.

**Note:**
If the :option:`ReplacementString` is not a C++ attribute, but instead a
macro, then that macro must be defined in scope or the fix-it will not be
applied.
**Note:**
For alternative ``__attribute__`` syntax options to mark functions as
``[[nodiscard]]`` in non-c++17 source code.
See https://clang.llvm.org/docs/AttributeReference.html#nodiscard-warn-unused-result

**Title:** clang-tidy - modernize-use-noexcept

# modernize-use-noexcept

This check replaces deprecated dynamic exception specifications with
the appropriate noexcept specification (introduced in C++11). By
default this check will replace `throw()` with `noexcept`,
and `throw(<exception>[,...])` or `throw(...)` with
`noexcept(false)`.

## Example

```c++
void foo() throw();
void bar() throw(int) {}
```

transforms to:

```c++
void foo() noexcept;
void bar() noexcept(false) {}
```

## Options

**Option:** ReplacementString
Users can use :option:`ReplacementString` to specify a macro to use
instead of ``noexcept``. This is useful when maintaining source code
that uses custom exception specification marking other than
``noexcept``. Fix-it hints will only be generated for non-throwing
specifications.

### Example

```c++
void bar() throw(int);
void foo() throw();
```

transforms to:

```c++
void bar() throw(int); // No fix-it generated.
void foo() NOEXCEPT;
```

if the `ReplacementString` option is set to `NOEXCEPT`.

**Option:** UseNoexceptFalse
Enabled by default, disabling will generate fix-it hints that remove
throwing dynamic exception specs, e.g., `throw(<something>)`,
completely without providing a replacement text, except for
destructors and delete operators that are `noexcept(true)` by
default.

### Example

```c++
void foo() throw(int) {}

struct bar {
 void foobar() throw(int);
 void operator delete(void *ptr) throw(int);
 void operator delete[](void *ptr) throw(int);
 ~bar() throw(int);
}
```

transforms to:

```c++
void foo() {}

struct bar {
 void foobar();
 void operator delete(void *ptr) noexcept(false);
 void operator delete[](void *ptr) noexcept(false);
 ~bar() noexcept(false);
}
```

if the `UseNoexceptFalse` option is set to `false`.

**Title:** clang-tidy - modernize-use-nullptr

# modernize-use-nullptr

The check converts the usage of null pointer constants (e.g. `NULL`, `0`)
to use the new C++11 and C23 `nullptr` keyword.

## Example

```c++
void assignment() {
 char *a = NULL;
 char *b = 0;
 char c = 0;
}

int *ret_ptr() {
 return 0;
}
```

transforms to:

```c++
void assignment() {
 char *a = nullptr;
 char *b = nullptr;
 char c = 0;
}

int *ret_ptr() {
 return nullptr;
}
```

## Options

**Option:** IgnoredTypes
Semicolon-separated list of regular expressions to match pointer types for
which implicit casts will be ignored. Default value:
`std::_CmpUnspecifiedParam::;^std::__cmp_cat::__unspec`.
**Option:** NullMacros
Comma-separated list of macro names that will be transformed along with
``NULL``. By default this check will only replace the ``NULL`` macro and will
skip any similar user-defined macros.

### Example

```c++
#define MY_NULL (void*)0
void assignment() {
 void *p = MY_NULL;
}
```

transforms to:

```c++
#define MY_NULL NULL
void assignment() {
 int *p = nullptr;
}
```

if the `NullMacros` option is set to `MY_NULL`.

**Title:** clang-tidy - modernize-use-override

# modernize-use-override

Adds `override` (introduced in C++11) to overridden virtual functions and
removes `virtual` from those functions as it is not required.

`virtual` on non base class implementations was used to help indicate to the
user that a function was virtual. C++ compilers did not use the presence of
this to signify an overridden function.

In C++11 `override` and `final` keywords were introduced to allow
overridden functions to be marked appropriately. Their presence allows
compilers to verify that an overridden function correctly overrides a base
class implementation.

This can be useful as compilers can generate a compile time error when:

 - The base class implementation function signature changes.
 - The user has not created the override with the correct signature.

## Options

**Option:** IgnoreDestructors
If set to `true`, this check will not diagnose destructors. Default is `false`.
**Option:** IgnoreTemplateInstantiations
If set to `true`, instructs this check to ignore virtual function overrides
that are part of template instantiations. Default is `false`.
**Option:** AllowOverrideAndFinal
If set to `true`, this check will not diagnose ``override`` as redundant
with ``final``. This is useful when code will be compiled by a compiler with
warning/error checking flags requiring ``override`` explicitly on overridden
members, such as ``gcc -Wsuggest-override``/``gcc -Werror=suggest-override``.
Default is `false`.
**Option:** AllowVirtualAndOverride
If set to `true`, this check will not diagnose ``virtual`` as redundant
with ``override``. Default is `false`.
**Option:** OverrideSpelling
Specifies a macro to use instead of ``override``. This is useful when
maintaining source code that also needs to compile with a pre-C++11
compiler.
**Option:** FinalSpelling
Specifies a macro to use instead of ``final``. This is useful when
maintaining source code that also needs to compile with a pre-C++11
compiler.
**Note:**
For more information on the use of ``override`` see https://en.cppreference.com/w/cpp/language/override

**Title:** clang-tidy - modernize-use-ranges

# modernize-use-ranges

Detects calls to standard library iterator algorithms that could be replaced
with a ranges version instead.

## Example

```c++
auto Iter1 = std::find(Items.begin(), Items.end(), 0);
auto NewEnd = std::unique(Items.begin(), Items.end());
auto AreSame = std::equal(Items1.cbegin(), Items1.cend(),
 std::begin(Items2), std::end(Items2));
```

Transforms to:

```c++
auto Iter1 = std::ranges::find(Items, 0);
auto NewEnd = std::ranges::unique(Items).begin();
auto AreSame = std::ranges::equal(Items1, Items2);
```

## Supported algorithms

Calls to the following std library algorithms are checked:

`std::adjacent_find`,
`std::all_of`,
`std::any_of`,
`std::binary_search`,
`std::copy_backward`,
`std::copy_if`,
`std::copy`,
`std::destroy`,
`std::equal_range`,
`std::equal`,
`std::fill`,
`std::find_end`,
`std::find_if_not`,
`std::find_if`,
`std::find`,
`std::for_each`,
`std::generate`,
`std::includes`,
`std::inplace_merge`,
`std::iota`,
`std::is_heap_until`,
`std::is_heap`,
`std::is_partitioned`,
`std::is_permutation`,
`std::is_sorted_until`,
`std::is_sorted`,
`std::lexicographical_compare`,
`std::lower_bound`,
`std::make_heap`,
`std::max_element`,
`std::merge`,
`std::min_element`,
`std::minmax_element`,
`std::mismatch`,
`std::move_backward`,
`std::move`,
`std::next_permutation`,
`std::none_of`,
`std::partial_sort_copy`,
`std::partition_copy`,
`std::partition_point`,
`std::partition`,
`std::pop_heap`,
`std::prev_permutation`,
`std::push_heap`,
`std::remove_copy_if`,
`std::remove_copy`,
`std::remove`, `std::remove_if`,
`std::replace_if`,
`std::replace`,
`std::reverse_copy`,
`std::reverse`,
`std::rotate`,
`std::rotate_copy`,
`std::sample`,
`std::search`,
`std::set_difference`,
`std::set_intersection`,
`std::set_symmetric_difference`,
`std::set_union`,
`std::shift_left`,
`std::shift_right`,
`std::sort_heap`,
`std::sort`,
`std::stable_partition`,
`std::stable_sort`,
`std::transform`,
`std::uninitialized_copy`,
`std::uninitialized_default_construct`,
`std::uninitialized_fill`,
`std::uninitialized_move`,
`std::uninitialized_value_construct`,
`std::unique_copy`,
`std::unique`,
`std::upper_bound`.

Note: some range algorithms for `vector<bool>` require C++23 because it uses
proxy iterators.

## Reverse Iteration

If calls are made using reverse iterators on containers, The code will be
fixed using the `std::views::reverse` adaptor.

```c++
auto AreSame = std::equal(Items1.rbegin(), Items1.rend(),
 std::crbegin(Items2), std::crend(Items2));
```

Transforms to:

```c++
auto AreSame = std::ranges::equal(std::views::reverse(Items1),
 std::views::reverse(Items2));
```

## Options

**Option:** IncludeStyle
A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.
**Option:** UseReversePipe
When `true` (default `false`), fixes which involve reverse ranges will use the
pipe adaptor syntax instead of the function syntax.

.. code-block:: c++

 std::find(Items.rbegin(), Items.rend(), 0);

Transforms to:

.. code-block:: c++

 std::ranges::find(Items | std::views::reverse, 0);

**Title:** clang-tidy - modernize-use-scoped-lock

# modernize-use-scoped-lock

Finds uses of `std::lock_guard` and suggests replacing them with C++17's
alternative `std::scoped_lock`.

Fix-its are provided for single declarations of `std::lock_guard` and warning
is emitted for multiple declarations of `std::lock_guard` or `std::scoped_lock`
that can be replaced with a single declaration of `std::scoped_lock`.

## Examples

Single `std::lock_guard` declaration:

```c++
std::mutex M;
std::lock_guard<std::mutex> L(M);
```

Transforms to:

```c++
std::mutex M;
std::scoped_lock L(M);
```

Single `std::lock_guard` declaration with `std::adopt_lock`:

```c++
std::mutex M;
std::lock(M);
std::lock_guard<std::mutex> L(M, std::adopt_lock);
```

Transforms to:

```c++
std::mutex M;
std::lock(M);
std::scoped_lock L(std::adopt_lock, M);
```

Multiple consecutive lock declarations only emit warnings:

```c++
std::mutex M1, M2;
std::lock(M1, M2);
std::lock_guard Lock1(M1, std::adopt_lock); // warning: use single 'std::scoped_lock' instead of multiple locks
std::scoped_lock<std::mutex> Lock2(M2, std::adopt_lock); // note: additional 'std::scoped_lock' declared here
```

## Limitations

The check will not emit warnings if `std::lock_guard` is used implicitly via
`template` parameter:

```c++
template <template <typename> typename Lock>
void TemplatedLock() {
 std::mutex M;
 Lock<std::mutex> L(M); // no warning
}

void instantiate() {
 TemplatedLock<std::lock_guard>();
}
```

## Options

**Option:** WarnOnSingleLocks
When `true`, the check will warn on single ``std::lock_guard`` declarations.
Set this option to `false` if you want to get warnings only on multiple
``std::lock_guard`` declarations that can be replaced with a single
``std::scoped_lock``. Default is `true`.
**Option:** WarnOnUsingAndTypedef
When `true`, the check will emit warnings if ``std::lock_guard`` is used
in ``using`` or ``typedef`` context. Default is `true`.

.. code-block:: c++

 template <typename T>
 using Lock = std::lock_guard<T>; // warning: use 'std::scoped_lock' instead of 'std::lock_guard'

 using LockMutex = std::lock_guard<std::mutex>; // warning: use 'std::scoped_lock' instead of 'std::lock_guard'

 typedef std::lock_guard<std::mutex> LockDef; // warning: use 'std::scoped_lock' instead of 'std::lock_guard'

 using std::lock_guard; // warning: use 'std::scoped_lock' instead of 'std::lock_guard'

```{title} clang-tidy - modernize-use-starts-ends-with
```

# modernize-use-starts-ends-with

Checks for common roundabout ways to express `starts_with` and `ends_with`
and suggests replacing with the simpler method when it is available. Notably,
this will work with `std::string` and `std::string_view`.

Covered scenarios:

| Expression | Replacement |
| ------------------------------------------------------- | ------------------- |
| `u.find(v) == 0` | `u.starts_with(v)` |
| `u.find(v, 0) == 0` | `u.starts_with(v)` |
| `u.find(v, 0, v.size()) == 0` | `u.starts_with(v)` |
| `u.rfind(v, 0) != 0` | `!u.starts_with(v)` |
| `u.rfind(v, 0, v.size()) != 0` | `!u.starts_with(v)` |
| `u.compare(0, v.size(), v) == 0` | `u.starts_with(v)` |
| `u.substr(0, v.size()) == v` | `u.starts_with(v)` |
| `v != u.substr(0, v.size())` | `!u.starts_with(v)` |
| `u.compare(u.size() - v.size(), v.size(), v) == 0` | `u.ends_with(v)` |
| `u.rfind(v) == u.size() - v.size()` | `u.ends_with(v)` |

**Title:** clang-tidy - modernize-use-std-bit

# modernize-use-std-bit

Finds common idioms which can be replaced by standard functions from the
`<bit>` C++20 header.

Covered scenarios:

============================== ==========================
Expression Replacement
------------------------------ --------------------------
`x && !(x & (x - 1))` `std::has_single_bit(x)`
`(x != 0) && !(x & (x - 1))` `std::has_single_bit(x)`
`(x > 0) && !(x & (x - 1))` `std::has_single_bit(x)`
`std::bitset<N>(x).count()` `std::popcount(x)`
`x << 3 | x >> 61` `std::rotl(x, 3)`
`x << 61 | x >> 3` `std::rotr(x, 3)`
============================== ==========================

## Options

**Option:** HonorIntPromotion
When set to ``true`` (default is ``false``), insert explicit cast to make sure the
type of the substituted expression is unchanged. Example:

.. code:: c++

 // Return type is deduced as 'int' (not 'unsigned char') due to implicit conversions.
 auto foo(unsigned char x) {
 return x << 3 | x >> 5;
 }

Becomes:

.. code:: c++

 #include <bit>

 auto foo(unsigned char x) {
 return static_cast<int>(std::rotl(x, 3));
 }

**Title:** clang-tidy - modernize-use-std-format

# modernize-use-std-format

Converts calls to `absl::StrFormat`, or other functions via
configuration options, to C++20's `std::format`, or another function
via a configuration option, modifying the format string appropriately and
removing now-unnecessary calls to `std::string::c_str()` and
`std::string::data()`.

For example, it turns lines like

```c++
return absl::StrFormat("The %s is %3d", description.c_str(), value);
```

into:

```c++
return std::format("The {} is {:3}", description, value);
```

The check uses the same format-string-conversion algorithm as
modernize-use-std-print and its
shortcomings and behaviour in combination with macros are described in the
documentation for that check.

## Options

**Option:** StrictMode
 When `true`, the check will add casts when converting from variadic
 functions and printing signed or unsigned integer types (including
 fixed-width integer types from ``<cstdint>``, ``ptrdiff_t``, ``size_t``
 and ``ssize_t``) as the opposite signedness to ensure that the output
 would matches that of a simple wrapper for ``std::sprintf`` that
 accepted a C-style variable argument list. For example, with
 `StrictMode` enabled,

.. code-block:: c++

 extern std::string strprintf(const char *format, ...);
 int i = -42;
 unsigned int u = 0xffffffff;
 return strprintf("%u %d\n", i, u);

would be converted to

.. code-block:: c++

 return std::format("{} {}\n", static_cast<unsigned int>(i), static_cast<int>(u));

to ensure that the output will continue to be the unsigned representation
of -42 and the signed representation of 0xffffffff (often 4294967254
and -1 respectively). When `false` (which is the default), these casts
will not be added which may cause a change in the output. Note that this
option makes no difference for the default value of
`StrFormatLikeFunctions` since ``absl::StrFormat`` takes a function
parameter pack and is not a variadic function.
**Option:** StrFormatLikeFunctions
A semicolon-separated list of regular expressions matching the
(fully qualified) names of functions to replace, with the requirement that
the first parameter contains the printf-style format string and the
arguments to be formatted follow immediately afterwards. Qualified member
function names are supported, but the replacement function name must be
unqualified. The default value is `absl::StrFormat`.
**Option:** ReplacementFormatFunction
The function that will be used to replace the function set by the
`StrFormatLikeFunctions` option rather than the default
`std::format`. It is expected that the function provides an interface
that is compatible with ``std::format``. A suitable candidate would be
`fmt::format`.
**Option:** FormatHeader
The header that must be included for the declaration of
`ReplacementFormatFunction` so that a ``#include`` directive can be added if
required. If `ReplacementFormatFunction` is `std::format` then this option will
default to ``<format>``, otherwise this option will default to nothing
and no ``#include`` directive will be added.

**Title:** clang-tidy - modernize-use-std-numbers

# modernize-use-std-numbers

Finds constants and function calls to math functions that can be replaced
with C++20's mathematical constants from the `numbers` header and offers
fix-it hints.
Does not match the use of variables with that value, and instead,
offers a replacement for the definition of those variables.
Function calls that match the pattern of how the constant is calculated are
matched and replaced with the `std::numbers` constant.
The use of macros gets replaced with the corresponding `std::numbers`
constant, instead of changing the macro definition.

The following list of constants from the `numbers` header are supported:

* `e`
* `log2e`
* `log10e`
* `pi`
* `inv_pi`
* `inv_sqrtpi`
* `ln2`
* `ln10`
* `sqrt2`
* `sqrt3`
* `inv_sqrt3`
* `egamma`
* `phi`

The list currently includes all constants as of C++20.

The replacements use the type of the matched constant and can remove explicit
casts, i.e., switching between `std::numbers::e`,
`std::numbers::e_v<float>` and `std::numbers::e_v<long double>` where
appropriate.

```c++
double sqrt(double);
double log2(double);
void sink(auto&&) {}
void floatSink(float);

#define MY_PI 3.1415926

void foo() {
 const double Pi = 3.141592653589; // const double Pi = std::numbers::pi
 const auto Use = Pi / 2; // no match for Pi
 static constexpr double Euler = 2.7182818; // static constexpr double Euler = std::numbers::e;

 log2(exp(1)); // std::numbers::log2e;
 log2(Euler); // std::numbers::log2e;
 1 / sqrt(MY_PI); // std::numbers::inv_sqrtpi;
 sink(MY_PI); // sink(std::numbers::pi);
 floatSink(MY_PI); // floatSink(std::numbers::pi);
 floatSink(static_cast<float>(MY_PI)); // floatSink(std::numbers::pi_v<float>);
}
```

## Options

**Option:** DiffThreshold
A floating point value that sets the detection threshold for when literals
match a constant. A literal matches a constant if
``abs(literal - constant) < DiffThreshold`` evaluates to ``true``. Default
is `0.001`.
**Option:** IncludeStyle
A string specifying which include-style is used, `llvm` or `google`. Default
is `llvm`.

**Title:** clang-tidy - modernize-use-std-print

# modernize-use-std-print

Converts calls to `printf`, `fprintf`, `absl::PrintF` and
`absl::FPrintf` to equivalent calls to C++23's `std::print` or
`std::println` as appropriate, modifying the format string appropriately.
The replaced and replacement functions can be customised by configuration
options. Each argument that is the result of a call to
`std::string::c_str()` and `std::string::data()` will have that
now-unnecessary call removed in a similar manner to the
:doc:`readability-redundant-string-cstr
<../readability/redundant-string-cstr>` check.

In other words, it turns lines like:

```c++
fprintf(stderr, "The %s is %3d\n", description.c_str(), value);
```

into:

```c++
std::println(stderr, "The {} is {:3}", description, value);
```

If the `ReplacementPrintFunction` or `ReplacementPrintlnFunction` options
are left at or set to their default values then this check is only enabled
with `-std=c++23` or later.

Macros starting with `PRI` and `__PRI` from `<inttypes.h>` are
expanded, escaping is handled and adjacent strings are concatenated to form
a single `StringLiteral` before the format string is converted. Use of
any other macros in the format string will cause a warning message to be
emitted and no conversion will be performed. The converted format string
will always be a single string literal.

The check doesn't do a bad job, but it's not perfect. In particular:

- It assumes that the format string is correct for the arguments. If you
 get any warnings when compiling with `-Wformat` then misbehaviour is
 possible.

- At the point that the check runs, the AST contains a single
 `StringLiteral` for the format string where escapes have been expanded.
 The check tries to reconstruct escape sequences, they may not be the same
 as they were written (e.g. `"\x41\x0a"` will become `"A\n"` and
 `"ab" "cd"` will become `"abcd"`.)

- It supports field widths, precision, positional arguments, leading zeros,
 leading `+`, alignment and alternative forms.

- Use of any unsupported flags or specifiers will cause the entire
 statement to be left alone and a warning to be emitted. Particular
 unsupported features are:

 - The `%'` flag for thousands separators.

 - The glibc extension `%m`.

- `printf` and similar functions return the number of characters printed.
 `std::print` does not. This means that any invocations that use the
 return value will not be converted. Unfortunately this currently includes
 explicitly-casting to `void`. Deficiencies in this check mean that any
 invocations inside `GCC` compound statements cannot be converted even
 if the resulting value is not used.

If conversion would be incomplete or unsafe then the entire invocation will
be left unchanged.

If the call is deemed suitable for conversion then:

- `printf`, `fprintf`, `absl::PrintF`, `absl::FPrintF` and any
 functions specified by the `PrintfLikeFunctions` option or
 `FprintfLikeFunctions` are replaced with the function specified by the
 `ReplacementPrintlnFunction` option if the format string ends with `\n`
 or `ReplacementPrintFunction` otherwise.
- the format string is rewritten to use the `std::formatter` language. If
 a `\n` is found at the end of the format string not preceded by `r`
 then it is removed and `ReplacementPrintlnFunction` is used rather than
 `ReplacementPrintFunction`.
- any arguments that corresponded to `%p` specifiers that
 `std::formatter` wouldn't accept are wrapped in a `static_cast`
 to `const void *`.
- any arguments that corresponded to `%s` specifiers where the argument
 is of `signed char` or `unsigned char` type are wrapped in a
 `reinterpret_cast<const char *>`.
- any arguments where the format string and the parameter differ in
 signedness will be wrapped in an appropriate `static_cast` if `StrictMode`
 is enabled.
- any arguments that end in a call to `std::string::c_str()` or
 `std::string::data()` will have that call removed.

## Options

**Option:** StrictMode
 When `true`, the check will add casts when converting from variadic
 functions like ``printf`` and printing signed or unsigned integer types
 (including fixed-width integer types from ``<cstdint>``, ``ptrdiff_t``,
 ``size_t`` and ``ssize_t``) as the opposite signedness to ensure that
 the output matches that of ``printf``. This does not apply when
 converting from non-variadic functions such as ``absl::PrintF`` and
 ``fmt::printf``. For example, with `StrictMode` enabled:

.. code-block:: c++

 int i = -42;
 unsigned int u = 0xffffffff;
 printf("%u %d\n", i, u);

would be converted to:

.. code-block:: c++

 std::print("{} {}\n", static_cast<unsigned int>(i), static_cast<int>(u));

to ensure that the output will continue to be the unsigned representation
of `-42` and the signed representation of `0xffffffff` (often
`4294967254` and `-1` respectively.) When `false` (which is the default),
these casts will not be added which may cause a change in the output.
**Option:** PrintfLikeFunctions
A semicolon-separated list of regular expressions matching the
(fully qualified) names of functions to replace, with the requirement
that the first parameter contains the printf-style format string and the
arguments to be formatted follow immediately afterwards. Qualified member
function names are supported, but the replacement function name must be
unqualified. If neither this option nor `FprintfLikeFunctions` are set then
the default value is `printf; absl::PrintF`, otherwise it is the empty
string.
**Option:** FprintfLikeFunctions
A semicolon-separated list of regular expressions matching the
(fully qualified) names of functions to replace, with the requirement
that the first parameter is retained, the second parameter contains the
printf-style format string and the arguments to be formatted follow
immediately afterwards. Qualified member function names are supported,
but the replacement function name must be unqualified. If neither this
option nor `PrintfLikeFunctions` are set then the default value is
`fprintf;absl::FPrintF`, otherwise it is the empty string.
**Option:** ReplacementPrintFunction
The function that will be used to replace ``printf``, ``fprintf`` etc.
during conversion rather than the default ``std::print`` when the
originalformat string does not end with ``\n``. It is expected that the
function provides an interface that is compatible with ``std::print``. A
suitable candidate would be ``fmt::print``.
**Option:** ReplacementPrintlnFunction
The function that will be used to replace ``printf``, ``fprintf`` etc.
during conversion rather than the default ``std::println`` when the
original format string ends with ``\n``. It is expected that the
function provides an interface that is compatible with ``std::println``.
A suitable candidate would be ``fmt::println``.
**Option:** PrintHeader
The header that must be included for the declaration of
`ReplacementPrintFunction` so that a ``#include`` directive can be
added if required. If `ReplacementPrintFunction` is ``std::print``
then this option will default to ``<print>``, otherwise this option will
default to nothing and no ``#include`` directive will be added.

```{title} clang-tidy - mpi-buffer-deref
```

# mpi-buffer-deref

This check verifies if a buffer passed to an MPI (Message Passing Interface)
function is sufficiently dereferenced. Buffers should be passed as a single
pointer or array. As MPI function signatures specify `void *` for their
buffer types, insufficiently dereferenced buffers can be passed, like for
example as double pointers or multidimensional arrays, without a compiler
warning emitted.

Examples:

```cpp
// A double pointer is passed to the MPI function.
char *buf;
MPI_Send(&buf, 1, MPI_CHAR, 0, 0, MPI_COMM_WORLD);

// A multidimensional array is passed to the MPI function.
short buf[1][1];
MPI_Send(buf, 1, MPI_SHORT, 0, 0, MPI_COMM_WORLD);

// A pointer to an array is passed to the MPI function.
short *buf[1];
MPI_Send(buf, 1, MPI_SHORT, 0, 0, MPI_COMM_WORLD);
```

```{title} clang-tidy - objc-assert-equals
```

# objc-assert-equals

Finds improper usages of `XCTAssertEqual` and `XCTAssertNotEqual` and replaces
them with `XCTAssertEqualObjects` or `XCTAssertNotEqualObjects`.

This makes tests less fragile, as many improperly rely on pointer equality for
strings that have equal values. This assumption is not guaranteed by the
language.

```{title} clang-tidy - openmp-exception-escape
```

# openmp-exception-escape

Analyzes OpenMP Structured Blocks and checks that no exception escapes
out of the Structured Block it was thrown in.

As per the OpenMP specification, a structured block is an executable statement,
possibly compound, with a single entry at the top and a single exit at the
bottom. Which means, `throw` may not be used to 'exit' out of the
structured block. If an exception is not caught in the same structured block
it was thrown in, the behavior is undefined.

FIXME: this check does not model SEH, `setjmp`/`longjmp`.

WARNING! This check may be expensive on large source files.

## Options

```{option} IgnoredExceptions

Comma-separated list containing type names which are not counted as thrown
exceptions in the check. Default value is an empty string.
```

**Title:** clang-tidy - performance-avoid-endl

# performance-avoid-endl

Checks for uses of `std::endl` on streams and suggests using the newline
character `'\n'` instead.

Rationale:
Using `std::endl` on streams can be less efficient than using the newline
character `'\n'` because `std::endl` performs two operations: it writes a
newline character to the output stream and then flushes the stream buffer.
Writing a single newline character using `'\n'` does not trigger a flush,
which can improve performance. In addition, flushing the stream buffer can
cause additional overhead when working with streams that are buffered.

Example:

Consider the following code:

```c++
#include <iostream>

int main() {
 std::cout << "Hello" << std::endl;
}
```

Which gets transformed into:

```c++
#include <iostream>

int main() {
 std::cout << "Hello" << '\n';
}
```

This code writes a single newline character to the `std::cout` stream without
flushing the stream buffer.

Additionally, it is important to note that the standard C++ streams (like
`std::cerr`, `std::wcerr`, `std::clog` and `std::wclog`)
always flush after a write operation, unless `std::ios_base::sync_with_stdio`
is set to `false`. regardless of whether `std::endl` or `'\n'` is used.
Therefore, using `'\n'` with these streams will not
result in any performance gain, but it is still recommended to use
`'\n'` for consistency and readability.

If you do need to flush the stream buffer, you can use `std::flush`
explicitly like this:

```c++
#include <iostream>

int main() {
 std::cout << "Hello\n" << std::flush;
}
```

**Title:** clang-tidy - portability-avoid-pragma-once

# portability-avoid-pragma-once

Finds uses of `#pragma once` and suggests replacing them with standard
include guards (`#ifndef`/`#define`/`#endif`) for improved portability.

`#pragma once` is a non-standard extension, despite being widely supported
by modern compilers. Relying on it can lead to portability issues in
some environments.

Some older or specialized C/C++ compilers, particularly in embedded systems,
may not fully support `#pragma once`.

It can also fail in certain file system configurations, like network drives
or complex symbolic links, potentially leading to compilation issues.

Consider the following header file:

```c++
// my_header.h
#pragma once // warning: avoid 'pragma once' directive; use include guards instead
```

The warning suggests using include guards:

```c++
// my_header.h
#ifndef PATH_TO_MY_HEADER_H // Good: use include guards.
#define PATH_TO_MY_HEADER_H

#endif // PATH_TO_MY_HEADER_H
```

**Title:** clang-tidy - readability-ambiguous-smartptr-reset-call

# readability-ambiguous-smartptr-reset-call

Finds potentially erroneous calls to `reset` method on smart pointers when
the pointee type also has a `reset` method. Having a `reset` method in
both classes makes it easy to accidentally make the pointer null when
intending to reset the underlying object.

```c++
struct Resettable {
 void reset() { /* Own reset logic */ }
};

auto ptr = std::make_unique<Resettable>();

ptr->reset(); // Calls underlying reset method
ptr.reset(); // Makes the pointer null
```

Both calls are valid C++ code, but the second one might not be what the
developer intended, as it destroys the pointed-to object rather than resetting
its state. It's easy to make such a typo because the difference between
`.` and `->` is really small.

The recommended approach is to make the intent explicit by using either member
access or direct assignment:

```c++
std::unique_ptr<Resettable> ptr = std::make_unique<Resettable>();

(*ptr).reset(); // Clearly calls underlying reset method
ptr = nullptr; // Clearly makes the pointer null
```

The default smart pointers and classes that are considered are
`std::unique_ptr`, `std::shared_ptr`, `boost::shared_ptr`. To specify
other smart pointers or other classes use the `SmartPointers` option.

**Note:**
The check may emit invalid fix-its and misleading warning messages when
specifying custom smart pointers or other classes in the
:option:`SmartPointers` option. For example, ``boost::scoped_ptr`` does not
have an ``operator=`` which makes fix-its invalid.
**Note:**
Automatic fix-its are enabled only if :program:`clang-tidy` is invoked with
the `--fix-notes` option.

## Options

**Option:** SmartPointers
Semicolon-separated list of fully qualified class names of custom smart
pointers. Default value is `::std::unique_ptr;::std::shared_ptr;
::boost::shared_ptr`.
