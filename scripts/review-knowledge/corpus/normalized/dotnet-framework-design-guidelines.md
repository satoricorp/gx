Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

This section provides guidelines for designing libraries that extend and interact with .NET. The goal is to help library designers ensure API consistency and ease of use by providing a unified programming model that is independent of the programming language used for development. We recommend that you follow these design guidelines when developing classes and components that extend .NET. Inconsistent library design adversely affects developer productivity and discourages adoption.

The guidelines are organized as simple recommendations prefixed with the terms `Do`, `Consider`, `Avoid`, and `Do not`. These guidelines are intended to help class library designers understand the trade-offs between different solutions. There might be situations where good library design requires that you violate these design guidelines. Such cases should be rare, and it is important that you have a clear and compelling reason for your decision.

These guidelines are excerpted from the book *Framework Design Guidelines: Conventions, Idioms, and Patterns for Reusable .NET Libraries, 2nd Edition*, by Krzysztof Cwalina and Brad Abrams, which was published in 2008. The book has since been fully revised in the third edition. Some of the information in these guidelines may be out-of-date.

## In this section

Naming Guidelines

Provides guidelines for naming assemblies, namespaces, types, and members in class libraries.

Type Design Guidelines

Provides guidelines for using static and abstract classes, interfaces, enumerations, structures, and other types.

Member Design Guidelines

Provides guidelines for designing and using properties, methods, constructors, fields, events, operators, and parameters.

Designing for Extensibility

Discusses extensibility mechanisms such as subclassing, using events, virtual members, and callbacks, and explains how to choose the mechanisms that best meet your framework's requirements.

Design Guidelines for Exceptions

Describes design guidelines for designing, throwing, and catching exceptions.

Usage Guidelines

Describes guidelines for using common types such as arrays, attributes, and collections, supporting serialization, and overloading equality operators.

Common Design Patterns

Provides guidelines for choosing and implementing dependency properties and the dispose pattern.

*Portions © 2005, 2009 Microsoft Corporation. All rights reserved.*

*Reprinted by permission of Pearson Education, Inc. from Framework Design Guidelines: Conventions, Idioms, and Patterns for Reusable .NET Libraries, 2nd Edition by Krzysztof Cwalina and Brad Abrams, published Oct 22, 2008 by Addison-Wesley Professional as part of the Microsoft Windows Development Series.*

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Following a consistent set of naming conventions in developing a framework can be a major contribution to the framework’s usability. It allows the framework to be used by many developers on widely separated projects. Beyond consistency of form, the names of framework elements must be easily understood and convey each element's function.

The goal of this chapter is to provide a consistent set of naming conventions that results in names that make immediate sense to developers.

Adopting these naming conventions as general code development guidelines results in more consistent naming throughout your code. However, you're only required to apply them to APIs that are publicly exposed (public or protected types and members, and explicitly implemented interfaces).

## In this section

Capitalization Conventions

General Naming Conventions

Names of Assemblies and DLLs

Names of Namespaces

Names of Classes, Structs, and Interfaces

Names of Type Members

Naming Parameters

Naming Resources

*Portions © 2005, 2009 Microsoft Corporation. All rights reserved.*

*Reprinted by permission of Pearson Education, Inc. from Framework Design Guidelines: Conventions, Idioms, and Patterns for Reusable .NET Libraries, 2nd Edition by Krzysztof Cwalina and Brad Abrams, published Oct 22, 2008 by Addison-Wesley Professional as part of the Microsoft Windows Development Series.*

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

From the CLR perspective, there are only two categories of types—reference types and value types—but for the purpose of a discussion about framework design, we divide types into more logical groups, each with its own specific design rules.

Classes are the general case of reference types. They make up the bulk of types in the majority of frameworks. Classes owe their popularity to the rich set of object-oriented features they support and to their general applicability. Base classes and abstract classes are special logical groups related to extensibility.

Interfaces are types that can be implemented by both reference types and value types. They can thus serve as roots of polymorphic hierarchies of reference types and value types. In addition, interfaces can be used to simulate multiple inheritance, which is not natively supported by the CLR.

Structs are the general case of value types and should be reserved for small, simple types, similar to language primitives.

Enums are a special case of value types used to define short sets of values, such as days of the week, console colors, and so on.

Static classes are types intended to be containers for static members. They are commonly used to provide shortcuts to other operations.

Delegates, exceptions, attributes, arrays, and collections are all special cases of reference types intended for specific uses, and guidelines for their design and usage are discussed elsewhere in this book.

✔️ DO ensure that each type is a well-defined set of related members, not just a random collection of unrelated functionality.

## In this section

Choosing Between Class and Struct

Abstract Class Design

Static Class Design

Interface Design

Struct Design

Enum Design

Nested Types

*Portions © 2005, 2009 Microsoft Corporation. All rights reserved.*

*Reprinted by permission of Pearson Education, Inc. from Framework Design Guidelines: Conventions, Idioms, and Patterns for Reusable .NET Libraries, 2nd Edition by Krzysztof Cwalina and Brad Abrams, published Oct 22, 2008 by Addison-Wesley Professional as part of the Microsoft Windows Development Series.*

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Methods, properties, events, constructors, and fields are collectively referred to as members. Members are ultimately the means by which framework functionality is exposed to the end users of a framework.

Members can be virtual or nonvirtual, concrete or abstract, static or instance, and can have several different scopes of accessibility. All this variety provides incredible expressiveness but at the same time requires care on the part of the framework designer.

This chapter offers basic guidelines that should be followed when designing members of any type.

## In This Section

Member Overloading

Property Design

Constructor Design

Event Design

Field Design

Extension Methods

Operator Overloads

Parameter Design

*Portions © 2005, 2009 Microsoft Corporation. All rights reserved.*

*Reprinted by permission of Pearson Education, Inc. from Framework Design Guidelines: Conventions, Idioms, and Patterns for Reusable .NET Libraries, 2nd Edition by Krzysztof Cwalina and Brad Abrams, published Oct 22, 2008 by Addison-Wesley Professional as part of the Microsoft Windows Development Series.*

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

One important aspect of designing a framework is making sure the extensibility of the framework has been carefully considered. This requires that you understand the costs and benefits associated with various extensibility mechanisms. This chapter helps you decide which of the extensibility mechanisms—subclassing, events, virtual members, callbacks, and so on—can best meet the requirements of your framework.

There are many ways to allow extensibility in frameworks. They range from less powerful but less costly to very powerful but expensive. For any given extensibility requirement, you should choose the least costly extensibility mechanism that meets the requirements. Keep in mind that it’s usually possible to add more extensibility later, but you can never take it away without introducing breaking changes.

## In This Section

Unsealed Classes

Protected Members

Events and Callbacks

Virtual Members

Abstractions (Abstract Types and Interfaces)

Base Classes for Implementing Abstractions

Sealing

*Portions © 2005, 2009 Microsoft Corporation. All rights reserved.*

*Reprinted by permission of Pearson Education, Inc. from Framework Design Guidelines: Conventions, Idioms, and Patterns for Reusable .NET Libraries, 2nd Edition by Krzysztof Cwalina and Brad Abrams, published Oct 22, 2008 by Addison-Wesley Professional as part of the Microsoft Windows Development Series.*

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Exception handling has many advantages over return-value-based error reporting. Good framework design helps the application developer realize the benefits of exceptions. This section discusses the benefits of exceptions and presents guidelines for using them effectively.

## In This Section

Exception Throwing

Using Standard Exception Types

Exceptions and Performance

*Portions © 2005, 2009 Microsoft Corporation. All rights reserved.*

*Reprinted by permission of Pearson Education, Inc. from Framework Design Guidelines: Conventions, Idioms, and Patterns for Reusable .NET Libraries, 2nd Edition by Krzysztof Cwalina and Brad Abrams, published Oct 22, 2008 by Addison-Wesley Professional as part of the Microsoft Windows Development Series.*

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

This section contains guidelines for using common types in publicly accessible APIs. It deals with direct usage of built-in Framework types (e.g., serialization attributes) and overloading common operators.

The System.IDisposable interface is not covered in this section, but is discussed in the Dispose Pattern section.

Note

For guidelines and additional information about other common, built-in .NET Framework types, see the reference topics for the following: System.DateTime, System.DateTimeOffset, System.ICloneable, System.IComparable<T>, System.IEquatable<T>, System.Nullable<T>, System.Object, System.Uri.

## In this section

Arrays

Attributes

Collections

Serialization

System.Xml Usage

Equality Operators

*Portions © 2005, 2009 Microsoft Corporation. All rights reserved.*

*Reprinted by permission of Pearson Education, Inc. from Framework Design Guidelines: Conventions, Idioms, and Patterns for Reusable .NET Libraries, 2nd Edition by Krzysztof Cwalina and Brad Abrams, published Oct 22, 2008 by Addison-Wesley Professional as part of the Microsoft Windows Development Series.*

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

There are numerous books on software patterns, pattern languages, and antipatterns that address the very broad subject of patterns. Thus, this chapter provides guidelines and discussion related to a very limited set of patterns that are used frequently in the design of the .NET Framework APIs.

## In This Section

Dependency Properties

Dispose Pattern

*Portions © 2005, 2009 Microsoft Corporation. All rights reserved.*

*Reprinted by permission of Pearson Education, Inc. from Framework Design Guidelines: Conventions, Idioms, and Patterns for Reusable .NET Libraries, 2nd Edition by Krzysztof Cwalina and Brad Abrams, published Oct 22, 2008 by Addison-Wesley Professional as part of the Microsoft Windows Development Series.*
