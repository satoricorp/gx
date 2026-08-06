Each code analysis rule belongs to a category of rules. For example, design rules support adherence to the .NET design guidelines, and security rules help prevent security flaws. You can configure the severity level for an entire category of rules. You can also configure additional options on a per-category basis.

The following table shows the different code analysis rule categories and provides a link to the rules in each category. It also lists the configuration value to use in an EditorConfig file to bulk-configure rule severity on a per-category basis. For example, to set the severity of security rule violations to be errors, the EditorConfig entry is `dotnet_analyzer_diagnostic.category-Security.severity = error`.

Tip

Setting the severity for a category of rules using the `dotnet_analyzer_diagnostic.category-<category>.severity` syntax doesn't apply to rules that are disabled by default. However, starting in .NET 6, you can use the AnalysisMode<Category> project property to enable all the rules in a category.

## Design rules

## Documentation rules

| | Value |
| **Link to rules** | Documentation rules |
| **Description** | Documentation rules support writing well-documented libraries through the correct use of XML documentation comments for externally visible APIs. |
| **EditorConfig value** | `dotnet_analyzer_diagnostic.category-Documentation.severity` |
| **MSBuild property value** | `<AnalysisModeDocumentation>` |

## Globalization rules

| | Value |
| **Link to rules** | Globalization rules |
| **Description** | Globalization rules support world-ready libraries and applications. |
| **EditorConfig value** | `dotnet_analyzer_diagnostic.category-Globalization.severity` |
| **MSBuild property value** | `<AnalysisModeGlobalization>` |

## Portability and interoperability rules

| | Value |
| **Link to rules** | Portability and interoperability rules |
| **Description** | Portability rules support portability across different platforms. Interoperability rules support interaction with COM clients. |
| **EditorConfig value** | `dotnet_analyzer_diagnostic.category-Interoperability.severity` |
| **MSBuild property value** | `<AnalysisModeInteroperability>` |

## Maintainability rules

| | Value |
| **Link to rules** | Maintainability rules |
| **Description** | Maintainability rules support library and application maintenance. |
| **EditorConfig value** | `dotnet_analyzer_diagnostic.category-Maintainability.severity` |
| **MSBuild property value** | `<AnalysisModeMaintainability>` |

## Naming rules

| | Value |
| **Link to rules** | Naming rules |
| **Description** | Naming rules support adherence to the naming conventions of the .NET design guidelines. |
| **EditorConfig value** | `dotnet_analyzer_diagnostic.category-Naming.severity` |
| **MSBuild property value** | `<AnalysisModeNaming>` |

| | Value |
| **Link to rules** | Performance rules |
| **Description** | Performance rules support high-performance libraries and applications. |
| **EditorConfig value** | `dotnet_analyzer_diagnostic.category-Performance.severity` |
| **MSBuild property value** | `<AnalysisModePerformance>` |

## SingleFile rules

| | Value |
| **Link to rules** | SingleFile rules |
| **Description** | Single-file rules support single-file applications. |
| **EditorConfig value** | `dotnet_analyzer_diagnostic.category-SingleFile.severity` |
| **MSBuild property value** | `<AnalysisModeSingleFile>` |

## Reliability rules

| | Value |
| **Link to rules** | Reliability rules |
| **Description** | Reliability rules support library and application reliability, such as correct memory and thread usage. |
| **EditorConfig value** | `dotnet_analyzer_diagnostic.category-Reliability.severity` |
| **MSBuild property value** | `<AnalysisModeReliability>` |

## Security rules

| | Value |
| **Link to rules** | Security rules |
| **Description** | Security rules support safer libraries and applications. These rules help prevent security flaws in your program. |
| **EditorConfig value** | `dotnet_analyzer_diagnostic.category-Security.severity` |
| **MSBuild property value** | `<AnalysisModeSecurity>` |

## Style rules

| | Value |
| **Link to rules** | Style rules |
| **Description** | Style rules support consistent code style in your codebase. These rules start with the "IDE" prefix. * |
| **EditorConfig value** | `dotnet_analyzer_diagnostic.category-Style.severity` |
| **MSBuild property value** | `<AnalysisModeStyle>` |

* Use the EditorConfig value `dotnet_analyzer_diagnostic.category-CodeQuality.severity` to enable the following rules: IDE0051, IDE0052, IDE0064, and IDE0076. While these rules start with "IDE", they aren't technically part of the `Style` category.

## Usage rules

| | Value |
| **Link to rules** | Usage rules |
| **Description** | Usage rules support proper usage of .NET. |
| **EditorConfig value** | `dotnet_analyzer_diagnostic.category-Usage.severity` |
| **MSBuild property value** | `<AnalysisModeUsage>` |
