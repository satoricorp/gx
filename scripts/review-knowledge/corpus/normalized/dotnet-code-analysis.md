Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

.NET compiler platform (Roslyn) analyzers inspect your C# or Visual Basic code for code quality and style issues. These analyzers are included with the .NET SDK and you don't need to install them separately. If your project targets .NET 5 or later, code analysis is enabled by default. (For information about enabling code analysis in projects that target .NET Framework, see Enable code analysis in legacy projects.)

If rule violations are found by an analyzer, they're reported as a suggestion, warning, or error, depending on how each rule is configured. Code analysis violations appear with the prefix "CA" or "IDE" to differentiate them from compiler errors.

## Code quality analysis

*Code quality analysis* ("CAxxxx") rules inspect your C# or Visual Basic code for security, performance, design and other issues. Analysis is enabled, by default, for projects that target .NET 5 or later. You can enable code analysis on projects that target earlier .NET versions by setting the EnableNETAnalyzers property to `true`. You can also disable code analysis for your project by setting `EnableNETAnalyzers` to `false`.

Tip

If you're using Visual Studio, many analyzer rules have associated *code fixes* that you can apply to automatically correct the problem. Code fixes are shown in the light bulb icon menu.

### Enabled rules

The following rules are enabled, by default, as errors or warnings in .NET 10. Additional rules are enabled as suggestions.

| Diagnostic ID | Category | Severity | Version added | Description |
|---|---|---|---|---|
| CA1416 | Interoperability | Warning | .NET 5 | Validate platform compatibility |
| CA1417 | Interoperability | Warning | .NET 5 | Do not use `OutAttribute`on string parameters for P/Invokes |
| CA1418 | Interoperability | Warning | .NET 6 | Use valid platform string |
| CA1420 | Interoperability | Warning | .NET 7 | Using features that require runtime marshalling when it's disabled will result in runtime exceptions |
| CA1422 | Interoperability | Warning | .NET 7 | Validate platform compatibility |
| CA1831 | Performance | Warning | .NET 5 | Use `AsSpan`instead of range-based indexers for string when appropriate |
| CA1856 | Performance | Error | .NET 8 | Incorrect usage of `ConstantExpected`attribute |
| CA1857 | Performance | Warning | .NET 8 | A constant is expected for the parameter |
| CA2013 | Reliability | Warning | .NET 5 | Do not use `ReferenceEquals`with value types |
| CA2014 | Reliability | Warning | .NET 5 | Do not use `stackalloc`in loops |
| CA2015 | Reliability | Warning | .NET 5 | Do not define finalizers for types derived from MemoryManager<T> |
| CA2017 | Reliability | Warning | .NET 6 | Parameter count mismatch |
| CA2018 | Reliability | Warning | .NET 6 | The `count`argument to`Buffer.BlockCopy`should specify the number of bytes to copy |
| CA2021 | Reliability | Warning | .NET 8 | Do not call `Enumerable.Cast<T>`or`Enumerable.OfType<T>`with incompatible types |
| CA2022 | Reliability | Warning | .NET 9 | Avoid inexact read with `Stream.Read` |
| CA2023 | Reliability | Warning | .NET 10 | Invalid braces in message template |
| CA2200 | Usage | Warning | .NET 5 | Rethrow to preserve stack details |
| CA2247 | Usage | Warning | .NET 5 | Argument passed to `TaskCompletionSource`constructor should be TaskCreationOptions enum instead of TaskContinuationOptions |
| CA2252 | Usage | Error | .NET 6 | Opt in to preview features |
| CA2255 | Usage | Warning | .NET 6 | The `ModuleInitializer`attribute should not be used in libraries |
| CA2256 | Usage | Warning | .NET 6 | All members declared in parent interfaces must have an implementation in a `DynamicInterfaceCastableImplementation`-attributed interface |
| CA2257 | Usage | Warning | .NET 6 | Members defined on an interface with the `DynamicInterfaceCastableImplementationAttribute`should be`static` |
| CA2258 | Usage | Warning | .NET 6 | Providing a `DynamicInterfaceCastableImplementation`interface in Visual Basic is unsupported |
| CA2259 | Usage | Warning | .NET 7 | `ThreadStatic`only affects static fields |
| CA2260 | Usage | Warning | .NET 7 | Use correct type parameter |
| CA2261 | Usage | Warning | .NET 8 | Do not use `ConfigureAwaitOptions.SuppressThrowing`with`Task<TResult>` |
| CA2264 | Usage | Warning | .NET 9 | Do not pass a non-nullable value to `ArgumentNullException.ThrowIfNull` |
| CA2265 | Usage | Warning | .NET 9 | Do not compare `Span<T>`to`null`or`default` |
| CA2266 | Usage | Warning | .NET 10 | File-based program entry point should start with `#!` |

You can change the severity of these rules to disable them or elevate them to errors. You can also enable more rules.

- For a list of rules that are included with each .NET SDK version, see Analyzer releases.
- For a list of all the code quality rules, see Code quality rules.

### Enable additional rules

*Analysis mode* refers to a predefined code analysis configuration where none, some, or all rules are enabled. In the default analysis mode (`Default`), only a small number of rules are enabled as build warnings. You can change the analysis mode for your project by setting the `<AnalysisMode>` property in the project file. The allowable values are:

| Value | Description |
|---|---|
| `None` | All rules are disabled. You can selectively opt in to individual rules to enable them. |
| `Default` | Default mode, where certain rules are enabled as build warnings, certain rules are enabled as Visual Studio IDE suggestions, and the remainder are disabled. |
| `Minimum` | More aggressive mode than `Default`mode. Certain suggestions that are highly recommended for build enforcement are enabled as build warnings. To see which rules this includes, inspect the%ProgramFiles%/dotnet/sdk/[version]/Sdks/Microsoft.NET.Sdk/analyzers/build/config/analysislevel_[level]_minimum.globalconfigfile. (For .NET 7 and earlier versions, the file extension is.editorconfig.) |
| `Recommended` | More aggressive mode than `Minimum`mode, where more rules are enabled as build warnings. To see which rules this includes, inspect the%ProgramFiles%/dotnet/sdk/[version]/Sdks/Microsoft.NET.Sdk/analyzers/build/config/analysislevel_[level]_recommended.globalconfigfile. (For .NET 7 and earlier versions, the file extension is.editorconfig.) |
| `All` | All rules are enabled as build warnings *. You can selectively opt out of individual rules to disable them.* The following rules are notenabled by setting`AnalysisMode`to`All`or by setting`AnalysisLevel`to`latest-all`: CA1017, CA1045, CA1005, CA1014, CA1060, CA1021, and the code metrics analyzer rules (CA1501, CA1502, CA1505, CA1506, and CA1509). These legacy rules might be deprecated in a future version. However, you can still enable them individually using a`dotnet_diagnostic.CAxxxx.severity = <severity>`entry. |

You can also omit `<AnalysisMode>` in favor of a compound value for the `<AnalysisLevel>` property. For example, the following value enables the recommended set of rules for the latest release: `<AnalysisLevel>latest-Recommended</AnalysisLevel>`. For more information, see `AnalysisLevel`.

To find the default severity for each available rule and whether or not the rule is enabled in `Default` analysis mode, see the full list of rules.

### Treat warnings as errors

If you use the `-warnaserror` flag when you build your projects, all code analysis warnings are also treated as errors. If you do not want code quality warnings (CAxxxx) to be treated as errors in presence of `-warnaserror`, you can set the `CodeAnalysisTreatWarningsAsErrors` MSBuild property to `false` in your project file.

```
<PropertyGroup>
 <CodeAnalysisTreatWarningsAsErrors>false</CodeAnalysisTreatWarningsAsErrors>
</PropertyGroup>
```
You'll still see any code analysis warnings, but they won't break your build.

### Latest updates

By default, you'll get the latest code analysis rules and default rule severities as you upgrade to newer versions of the .NET SDK. If you don't want this behavior, for example, if you want to ensure that no new rules are enabled or disabled, you can override it in one of the following ways:

- Set the - `AnalysisLevel`MSBuild property to a specific value to lock the warnings to that set. When you upgrade to a newer SDK, you'll still get bug fixes for those warnings, but no new warnings will be enabled and no existing warnings will be disabled. For example, to lock the set of rules to those that ship with version 8.0 of the .NET SDK, add the following entry to your project file.- `<PropertyGroup> <AnalysisLevel>8.0</AnalysisLevel> </PropertyGroup>`- Tip - The default value for the - `AnalysisLevel`property is- `latest`, which means you always get the latest code analysis rules as you move to newer versions of the .NET SDK.- For more information, and to see a list of possible values, see AnalysisLevel.
- Install the Microsoft.CodeAnalysis.NetAnalyzers NuGet package to decouple rule updates from .NET SDK updates. For projects that target .NET 5+, installing the package turns off the built-in SDK analyzers. You'll get a build warning if the SDK contains a newer analyzer assembly version than that of the NuGet package. To disable the warning, set the - `_SkipUpgradeNetAnalyzersNuGetWarning`property to- `true`.- Note - If you install the Microsoft.CodeAnalysis.NetAnalyzers NuGet package, you should not add the EnableNETAnalyzers property to either your project file or a - *Directory.Build.props*file. When the NuGet package is installed and the- `EnableNETAnalyzers`property is set to- `true`, a build warning is generated.

## Code-style analysis

*Code-style analysis* ("IDExxxx") rules enable you to define and maintain consistent code style in your codebase. The default enablement settings are:

- Command-line build: Code-style analysis is - *disabled*, by default, for all .NET projects on command-line builds.- You can enable code-style analysis on build, both at the command line and inside Visual Studio. Code style violations appear as warnings or errors with an "IDE" prefix. This enables you to enforce consistent code styles at build time.
- Visual Studio: Code-style analysis is - *enabled*, by default, for all .NET projects inside Visual Studio as code refactoring quick actions.

For a full list of code-style analysis rules, see Code style rules.

### Enable on build

You can enable code-style analysis when building from the command-line and in Visual Studio. (However, for performance reasons, a handful of code-style rules will still apply only in the Visual Studio IDE.)

Follow these steps to enable code-style analysis on build:

- Set the MSBuild property EnforceCodeStyleInBuild to - `true`.
- In an - *.editorconfig*file, configure each "IDE" code style rule that you wish to run on build as a warning or an error. For example:- `[*.{cs,vb}] # IDE0040: Accessibility modifiers required (escalated to a build warning) dotnet_diagnostic.IDE0040.severity = warning`- Tip - Starting in .NET 9, you can also use the option format to specify a severity and have it be respected at build time. For example: - `[*.{cs,vb}] # IDE0040: Accessibility modifiers required (escalated to a build warning) dotnet_style_require_accessibility_modifiers = always:warning`- Alternatively, you can configure an entire category to be a warning or error, by default, and then selectively turn off rules in that category that you don't want to run on build. For example: - `[*.{cs,vb}] # Default severity for analyzer diagnostics with category 'Style' (escalated to build warnings) dotnet_analyzer_diagnostic.category-Style.severity = warning # IDE0040: Accessibility modifiers required (disabled on build) dotnet_diagnostic.IDE0040.severity = silent`

## Suppress a warning

One way to suppress a rule violation is to set the severity option for that rule ID to `none` in an EditorConfig file. For example:

```
dotnet_diagnostic.CA1822.severity = none
```
For more information and other ways to suppress warnings, see How to suppress code analysis warnings.

## Enable code analysis in legacy projects

If your project targets .NET 5 or later, code analysis is enabled by default. If your project targets .NET Standard or .NET Framework, you must manually enable code analysis by setting the EnableNETAnalyzers property to `true`.

If your project uses the legacy project file format, that is, it doesn't reference a project SDK, there are some additional steps you must take to enable code analysis:

- Add a reference to the 📦 Microsoft.CodeAnalysis.NetAnalyzers NuGet package.
- Instead of setting `AnalysisLevel`, which isn't understood by non-SDK-style projects, add the following properties to your project file:

```
 <EffectiveAnalysisLevel>9</EffectiveAnalysisLevel>
 <AnalysisMode>Recommended</AnalysisMode>
```
## Third-party analyzers

In addition to the official .NET analyzers, you can also install third party analyzers, such as StyleCop, Roslynator, Meziantou.Analyzer, XUnit Analyzers, and Sonar Analyzer.
