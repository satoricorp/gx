Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

ASP.NET Core enables developers to configure and manage security. The following list provides links to articles about working with security in ASP.NET Core:

- Authentication
- Authorization
- Data protection
- HTTPS enforcement
- Safe storage of app secrets in development
- XSRF/CSRF prevention
- Cross Origin Resource Sharing (CORS)
- Cross-Site Scripting (XSS) attacks

These security features allow you to build robust and secure ASP.NET Core apps.

For Blazor security coverage, which adds to or supersedes the guidance in this node, see ASP.NET Core Blazor authentication and authorization and the other articles in Blazor's *Security and Identity* node.

## ASP.NET Core security features

ASP.NET Core provides many tools and libraries to secure ASP.NET Core apps, such as built-in identity providers and non-Microsoft identity services like Facebook, Twitter, and LinkedIn. ASP.NET Core provides several approaches to store app secrets.

## Authentication vs. Authorization

Authentication is a process where a user provides credentials that are compared to credentials stored in an operating system, database, app, or resource. When the two sets of credentials match, the user authenticates successfully. They can then perform actions for which they're authorized. The authorization process determines the actions the user is allowed to do.

Another way to think of authentication is to consider it as a way to **enter** a space, where the space is a server, database, app, or resource. Authorization defines **what actions** the user can perform to which objects inside that space (server, database, or app).

## Common vulnerabilities in software

ASP.NET Core and Entity Framework contain features that help you secure your apps and prevent security breaches. The following list of links takes you to documentation detailing techniques to avoid the most common security vulnerabilities in web apps:

- Cross-Site Scripting (XSS) attacks
- SQL queries > SQL injection attacks
- Cross-Site Request Forgery (XSRF/CSRF) attacks
- Open redirect attacks

There are more vulnerabilities that you should be aware of. For more information, see the other articles in the **Security and Identity** section of the table of contents.

## Secure authentication flows

We recommend using the most secure authentication option. For Azure services, the most secure authentication is managed identities.

Avoid using the Resource Owner Password Credentials (ROPG) grant:

- It exposes the user's password to the client.
- It's a significant security risk.
- Use it only when other authentication flows aren't possible.

Managed identities are a secure way to authenticate to services without needing to store credentials in code, environment variables, or configuration files. Managed identities are available for Azure services, and can be used with Azure SQL, Azure Storage, and other Azure services:

- Managed identities in Microsoft Entra for Azure SQL
- Managed identities for App Service and Azure Functions
- Secure authentication flows

When the app is deployed to a test server, an environment variable can be used to set the connection string to a test database server. For more information, see Configuration. Environment variables are commonly stored in plain, unencrypted text. If the machine or process is compromised, environment variables might be accessible to untrusted parties. We recommend against using environment variables to store a production connection string as it's not the most secure approach.

Configuration data guidelines:

- Never store passwords or other sensitive data in configuration provider code or in plain text configuration files. The Secret Manager tool can be used to store secrets in development.
- Don't use production secrets in development or test environments.
- Specify secrets outside of the project so that they can't be accidentally committed to a source code repository.

For more information, see:

- Managed identity best practice recommendations
- Connecting from your application to resources without handling credentials in your code
- Azure services that can use managed identities to access other services
- IETF OAuth 2.0 Security Best Current Practice (Section 2.4. Resource Owner Password Credentials Grant)

For information on other cloud providers, see:

- AWS (Amazon Web Services): AWS Key Management Service (KMS)
- Google Cloud Key Management Service overview

## Enterprise web app patterns

For guidance on creating a reliable, secure, performant, testable, and scalable ASP.NET Core app, see Enterprise web app patterns. A complete production-quality sample web app that implements the patterns is available.

## Related content

ASP.NET Core

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

By Mike Rousos

Authentication is the process of determining a user's identity. Authorization is the process of determining whether a user has access to a resource. In ASP.NET Core, authentication is handled by the authentication service, IAuthenticationService, which is used by authentication middleware. The authentication service uses registered authentication handlers to complete authentication-related actions. Examples of authentication-related actions include:

- Authenticating a user.
- Responding when an unauthenticated user tries to access a restricted resource.

The registered authentication handlers and their configuration options are called "schemes".

Authentication schemes are specified by registering authentication services in `Program.cs`:

- By calling a scheme-specific extension method after a call to AddAuthentication, such as AddJwtBearer or AddCookie. These extension methods use AuthenticationBuilder.AddScheme to register schemes with appropriate settings.
- Less commonly, by calling `AuthenticationBuilder.AddScheme`directly.

For example, the following code registers authentication services and handlers for cookie and JWT bearer authentication schemes:

```
builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
 .AddJwtBearer(JwtBearerDefaults.AuthenticationScheme,
 options => builder.Configuration.Bind("JwtSettings", options))
 .AddCookie(CookieAuthenticationDefaults.AuthenticationScheme,
 options => builder.Configuration.Bind("CookieSettings", options));
```
The `AddAuthentication` parameter JwtBearerDefaults.AuthenticationScheme is the name of the scheme to use by default when a specific scheme isn't requested.

If multiple schemes are used, authorization policies (or authorization attributes) can specify the authentication scheme (or schemes) they depend on to authenticate the user. In the example above, the cookie authentication scheme could be used by specifying its name (CookieAuthenticationDefaults.AuthenticationScheme by default, though a different name could be provided when calling `AddCookie`).

In some cases, the call to `AddAuthentication` is automatically made by other extension methods. For example, when using ASP.NET Core Identity, `AddAuthentication` is called internally.

The Authentication middleware is added in `Program.cs` by calling UseAuthentication. Calling `UseAuthentication` registers the middleware that uses the previously registered authentication schemes. Call `UseAuthentication` before any middleware that depends on users being authenticated.

## Authentication concepts

Authentication is responsible for providing the ClaimsPrincipal for authorization to make permission decisions against. There are multiple authentication scheme approaches to select which authentication handler is responsible for generating the correct set of claims:

- Authentication scheme
- The default authentication scheme, discussed in the next two sections.
- Directly set HttpContext.User.

When there is only a single authentication scheme registered, it becomes the default scheme. If multiple schemes are registered and the default scheme isn't specified, a scheme must be specified in the authorize attribute, otherwise, the following error is thrown:

InvalidOperationException: No authenticationScheme was specified, and there was no DefaultAuthenticateScheme found. The default schemes can be set using either AddAuthentication(string defaultScheme) or AddAuthentication(Action<AuthenticationOptions> configureOptions).

`DefaultScheme`

When there is only a single authentication scheme registered, the single authentication scheme:

- Is automatically used as the DefaultScheme.
- Eliminates the need to specify the `DefaultScheme`in AddAuthentication(IServiceCollection) or AddAuthenticationCore(IServiceCollection).

To disable automatically using the single authentication scheme as the `DefaultScheme`, call `AppContext.SetSwitch("Microsoft.AspNetCore.Authentication.SuppressAutoDefaultScheme")`.

### Authentication scheme

The authentication scheme can select which authentication handler is responsible for generating the correct set of claims. For more information, see Authorize with a specific scheme.

An authentication scheme is a name that corresponds to:

- An authentication handler.
- Options for configuring that specific instance of the handler.

Schemes are useful as a mechanism for referring to the authentication, challenge, and forbid behaviors of the associated handler. For example, an authorization policy can use scheme names to specify which authentication scheme (or schemes) should be used to authenticate the user. When configuring authentication, it's common to specify the default authentication scheme. The default scheme is used unless a resource requests a specific scheme. It's also possible to:

- Specify different default schemes to use for authenticate, challenge, and forbid actions.
- Combine multiple schemes into one using policy schemes.

### Authentication handler

An authentication handler:

- Is a type that implements the behavior of a scheme.
- Is derived from IAuthenticationHandler or AuthenticationHandler<TOptions>.
- Has the primary responsibility to authenticate users.

Based on the authentication scheme's configuration and the incoming request context, authentication handlers:

- Construct AuthenticationTicket objects representing the user's identity if authentication is successful.
- Return 'no result' or 'failure' if authentication is unsuccessful.
- Have methods for challenge and forbid actions for when users attempt to access resources:
- They're unauthorized to access (forbid).
- When they're unauthenticated (challenge).

`RemoteAuthenticationHandler<TOptions>` vs `AuthenticationHandler<TOptions>`

RemoteAuthenticationHandler<TOptions> is the class for authentication that requires a remote authentication step. When the remote authentication step is finished, the handler calls back to the `CallbackPath` set by the handler. The handler finishes the authentication step using the information passed to the HandleRemoteAuthenticateAsync callback path. OAuth 2.0 and OIDC both use this pattern. JWT and cookies don't since they can directly use the bearer header and cookie to authenticate. The remotely hosted provider in this case:

- Is the authentication provider.
- Examples include Facebook, Twitter, Google, Microsoft, and any other OIDC provider that handles authenticating users using the handlers mechanism.

### Authenticate

An authentication scheme's authenticate action is responsible for constructing the user's identity based on request context. It returns an AuthenticateResult indicating whether authentication was successful and, if so, the user's identity in an authentication ticket. See AuthenticateAsync. Authenticate examples include:

- A cookie authentication scheme constructing the user's identity from cookies.
- A JWT bearer scheme deserializing and validating a JWT bearer token to construct the user's identity.

### Challenge

An authentication challenge is invoked by Authorization when an unauthenticated user requests an endpoint that requires authentication. An authentication challenge is issued, for example, when an anonymous user requests a restricted resource or follows a login link. Authorization invokes a challenge using the specified authentication schemes, or the default if none is specified. See ChallengeAsync. Authentication challenge examples include:

- A cookie authentication scheme redirecting the user to a login page.
- A JWT bearer scheme returning a 401 result with a `www-authenticate: bearer`header.

A challenge action should let the user know what authentication mechanism to use to access the requested resource.

### Forbid

An authentication scheme's forbid action is called by Authorization when an authenticated user attempts to access a resource they're not permitted to access. See ForbidAsync. Authentication forbid examples include:

- A cookie authentication scheme redirecting the user to a page indicating access was forbidden.
- A JWT bearer scheme returning a 403 result.
- A custom authentication scheme redirecting to a page where the user can request access to the resource.

A forbid action can let the user know:

- They're authenticated.
- They're not permitted to access the requested resource.

See the following links for differences between challenge and forbid:

## Authentication providers per tenant

ASP.NET Core doesn't have a built-in solution for multi-tenant authentication. While it's possible for customers to write one using the built-in features, we recommend customers consider Orchard Core, ABP Framework, or Finbuckle.MultiTenant for multi-tenant authentication.

Orchard Core is:

- An open-source, modular, and multi-tenant app framework built with ASP.NET Core.
- A content management system (CMS) built on top of that app framework.

See the Orchard Core source for an example of authentication providers per tenant.

ABP Framework supports various architectural patterns including modularity, microservices, domain driven design, and multi-tenancy. See ABP Framework source on GitHub.

Finbuckle.MultiTenant:

- Open source
- Provides tenant resolution
- Lightweight
- Provides data isolation
- Configure app behavior uniquely for each tenant

## Additional resources

By Mike Rousos

Authentication is the process of determining a user's identity. Authorization is the process of determining whether a user has access to a resource. In ASP.NET Core, authentication is handled by the authentication service, IAuthenticationService, which is used by authentication middleware. The authentication service uses registered authentication handlers to complete authentication-related actions. Examples of authentication-related actions include:

- Authenticating a user.
- Responding when an unauthenticated user tries to access a restricted resource.

The registered authentication handlers and their configuration options are called "schemes".

Authentication schemes are specified by registering authentication services in `Program.cs`:

- By calling a scheme-specific extension method after a call to AddAuthentication, such as AddJwtBearer or AddCookie. These extension methods use AuthenticationBuilder.AddScheme to register schemes with appropriate settings.
- Less commonly, by calling `AuthenticationBuilder.AddScheme`directly.

For example, the following code registers authentication services and handlers for cookie and JWT bearer authentication schemes:

```
builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
 .AddJwtBearer(JwtBearerDefaults.AuthenticationScheme,
 options => builder.Configuration.Bind("JwtSettings", options))
 .AddCookie(CookieAuthenticationDefaults.AuthenticationScheme,
 options => builder.Configuration.Bind("CookieSettings", options));
```
The `AddAuthentication` parameter JwtBearerDefaults.AuthenticationScheme is the name of the scheme to use by default when a specific scheme isn't requested.

If multiple schemes are used, authorization policies (or authorization attributes) can specify the authentication scheme (or schemes) they depend on to authenticate the user. In the example above, the cookie authentication scheme could be used by specifying its name (CookieAuthenticationDefaults.AuthenticationScheme by default, though a different name could be provided when calling `AddCookie`).

In some cases, the call to `AddAuthentication` is automatically made by other extension methods. For example, when using ASP.NET Core Identity, `AddAuthentication` is called internally.

The Authentication middleware is added in `Program.cs` by calling UseAuthentication. Calling `UseAuthentication` registers the middleware that uses the previously registered authentication schemes. Call `UseAuthentication` before any middleware that depends on users being authenticated.

## Authentication concepts

Authentication is responsible for providing the ClaimsPrincipal for authorization to make permission decisions against. There are multiple authentication scheme approaches to select which authentication handler is responsible for generating the correct set of claims:

- Authentication scheme
- The default authentication scheme, discussed in the next section.
- Directly set HttpContext.User.

There's no automatic probing of schemes. If the default scheme isn't specified, the scheme must be specified in the authorize attribute, otherwise, the following error is thrown:

InvalidOperationException: No authenticationScheme was specified, and there was no DefaultAuthenticateScheme found. The default schemes can be set using either AddAuthentication(string defaultScheme) or AddAuthentication(Action<AuthenticationOptions> configureOptions).

### Authentication scheme

The authentication scheme can select which authentication handler is responsible for generating the correct set of claims. For more information, see Authorize with a specific scheme.

An authentication scheme is a name that corresponds to:

- An authentication handler.
- Options for configuring that specific instance of the handler.

Schemes are useful as a mechanism for referring to the authentication, challenge, and forbid behaviors of the associated handler. For example, an authorization policy can use scheme names to specify which authentication scheme (or schemes) should be used to authenticate the user. When configuring authentication, it's common to specify the default authentication scheme. The default scheme is used unless a resource requests a specific scheme. It's also possible to:

- Specify different default schemes to use for authenticate, challenge, and forbid actions.
- Combine multiple schemes into one using policy schemes.

### Authentication handler

An authentication handler:

- Is a type that implements the behavior of a scheme.
- Is derived from IAuthenticationHandler or AuthenticationHandler<TOptions>.
- Has the primary responsibility to authenticate users.

Based on the authentication scheme's configuration and the incoming request context, authentication handlers:

- Construct AuthenticationTicket objects representing the user's identity if authentication is successful.
- Return 'no result' or 'failure' if authentication is unsuccessful.
- Have methods for challenge and forbid actions for when users attempt to access resources:
- They're unauthorized to access (forbid).
- When they're unauthenticated (challenge).

`RemoteAuthenticationHandler<TOptions>` vs `AuthenticationHandler<TOptions>`

RemoteAuthenticationHandler<TOptions> is the class for authentication that requires a remote authentication step. When the remote authentication step is finished, the handler calls back to the `CallbackPath` set by the handler. The handler finishes the authentication step using the information passed to the HandleRemoteAuthenticateAsync callback path. OAuth 2.0 and OIDC both use this pattern. JWT and cookies don't since they can directly use the bearer header and cookie to authenticate. The remotely hosted provider in this case:

- Is the authentication provider.
- Examples include Facebook, Twitter, Google, Microsoft, and any other OIDC provider that handles authenticating users using the handlers mechanism.

### Authenticate

An authentication scheme's authenticate action is responsible for constructing the user's identity based on request context. It returns an AuthenticateResult indicating whether authentication was successful and, if so, the user's identity in an authentication ticket. See AuthenticateAsync. Authenticate examples include:

- A cookie authentication scheme constructing the user's identity from cookies.
- A JWT bearer scheme deserializing and validating a JWT bearer token to construct the user's identity.

### Challenge

An authentication challenge is invoked by Authorization when an unauthenticated user requests an endpoint that requires authentication. An authentication challenge is issued, for example, when an anonymous user requests a restricted resource or follows a login link. Authorization invokes a challenge using the specified authentication schemes, or the default if none is specified. See ChallengeAsync. Authentication challenge examples include:

- A cookie authentication scheme redirecting the user to a login page.
- A JWT bearer scheme returning a 401 result with a `www-authenticate: bearer`header.

A challenge action should let the user know what authentication mechanism to use to access the requested resource.

### Forbid

An authentication scheme's forbid action is called by Authorization when an authenticated user attempts to access a resource they're not permitted to access. See ForbidAsync. Authentication forbid examples include:

- A cookie authentication scheme redirecting the user to a page indicating access was forbidden.
- A JWT bearer scheme returning a 403 result.
- A custom authentication scheme redirecting to a page where the user can request access to the resource.

A forbid action can let the user know:

- They're authenticated.
- They're not permitted to access the requested resource.

See the following links for differences between challenge and forbid:

- Challenge and forbid with an operational resource handler.
- Differences between challenge and forbid.

## Authentication providers per tenant

ASP.NET Core doesn't have a built-in solution for multi-tenant authentication. While it's possible for customers to write one using the built-in features, we recommend customers consider Orchard Core or ABP Framework for multi-tenant authentication.

Orchard Core is:

- An open-source, modular, and multi-tenant app framework built with ASP.NET Core.
- A content management system (CMS) built on top of that app framework.

See the Orchard Core source for an example of authentication providers per tenant.

ABP Framework supports various architectural patterns including modularity, microservices, domain driven design, and multi-tenancy. See ABP Framework source on GitHub.

## Additional resources

By Mike Rousos

Authentication is the process of determining a user's identity. Authorization is the process of determining whether a user has access to a resource. In ASP.NET Core, authentication is handled by the authentication service, IAuthenticationService, which is used by authentication middleware. The authentication service uses registered authentication handlers to complete authentication-related actions. Examples of authentication-related actions include:

- Authenticating a user.
- Responding when an unauthenticated user tries to access a restricted resource.

The registered authentication handlers and their configuration options are called "schemes".

Authentication schemes are specified by registering authentication services in `Startup.ConfigureServices`:

- By calling a scheme-specific extension method after a call to AddAuthentication (such as AddJwtBearer or AddCookie, for example). These extension methods use AuthenticationBuilder.AddScheme to register schemes with appropriate settings.
- Less commonly, by calling `AuthenticationBuilder.AddScheme`directly.

For example, the following code registers authentication services and handlers for cookie and JWT bearer authentication schemes:

```
services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
 .AddJwtBearer(JwtBearerDefaults.AuthenticationScheme,
 options => Configuration.Bind("JwtSettings", options))
 .AddCookie(CookieAuthenticationDefaults.AuthenticationScheme,
 options => Configuration.Bind("CookieSettings", options));
```
The `AddAuthentication` parameter JwtBearerDefaults.AuthenticationScheme is the name of the scheme to use by default when a specific scheme isn't requested.

If multiple schemes are used, authorization policies (or authorization attributes) can specify the authentication scheme (or schemes) they depend on to authenticate the user. In the example above, the cookie authentication scheme could be used by specifying its name (CookieAuthenticationDefaults.AuthenticationScheme by default, though a different name could be provided when calling `AddCookie`).

In some cases, the call to `AddAuthentication` is automatically made by other extension methods. For example, when using ASP.NET Core Identity, `AddAuthentication` is called internally.

The Authentication middleware is added in `Startup.Configure` by calling UseAuthentication. Calling `UseAuthentication` registers the middleware that uses the previously registered authentication schemes. Call `UseAuthentication` before any middleware that depends on users being authenticated. When using endpoint routing, the call to `UseAuthentication` must go:

- After UseRouting, so that route information is available for authentication decisions.
- Before UseEndpoints, so that users are authenticated before accessing the endpoints.

## Authentication concepts

Authentication is responsible for providing the ClaimsPrincipal for authorization to make permission decisions against. There are multiple authentication scheme approaches to select which authentication handler is responsible for generating the correct set of claims:

- Authentication scheme
- The default authentication scheme, discussed in the next section.
- Directly set HttpContext.User.

There's no automatic probing of schemes. If the default scheme isn't specified, the scheme must be specified in the authorize attribute, otherwise, the following error is thrown:

InvalidOperationException: No authenticationScheme was specified, and there was no DefaultAuthenticateScheme found. The default schemes can be set using either AddAuthentication(string defaultScheme) or AddAuthentication(Action<AuthenticationOptions> configureOptions).

### Authentication scheme

The authentication scheme can select which authentication handler is responsible for generating the correct set of claims. For more information, see Authorize with a specific scheme.

An authentication scheme is a name that corresponds to:

- An authentication handler.
- Options for configuring that specific instance of the handler.

Schemes are useful as a mechanism for referring to the authentication, challenge, and forbid behaviors of the associated handler. For example, an authorization policy can use scheme names to specify which authentication scheme (or schemes) should be used to authenticate the user. When configuring authentication, it's common to specify the default authentication scheme. The default scheme is used unless a resource requests a specific scheme. It's also possible to:

- Specify different default schemes to use for authenticate, challenge, and forbid actions.
- Combine multiple schemes into one using policy schemes.

### Authentication handler

An authentication handler:

- Is a type that implements the behavior of a scheme.
- Is derived from IAuthenticationHandler or AuthenticationHandler<TOptions>.
- Has the primary responsibility to authenticate users.

Based on the authentication scheme's configuration and the incoming request context, authentication handlers:

- Construct AuthenticationTicket objects representing the user's identity if authentication is successful.
- Return 'no result' or 'failure' if authentication is unsuccessful.
- Have methods for challenge and forbid actions for when users attempt to access resources:
- They're unauthorized to access (forbid).
- When they're unauthenticated (challenge).

`RemoteAuthenticationHandler<TOptions>` vs `AuthenticationHandler<TOptions>`

RemoteAuthenticationHandler<TOptions> is the class for authentication that requires a remote authentication step. When the remote authentication step is finished, the handler calls back to the `CallbackPath` set by the handler. The handler finishes the authentication step using the information passed to the HandleRemoteAuthenticateAsync callback path. OAuth 2.0 and OIDC both use this pattern. JWT and cookies don't since they can directly use the bearer header and cookie to authenticate. The remotely hosted provider in this case:

- Is the authentication provider.
- Examples include Facebook, Twitter, Google, Microsoft, and any other OIDC provider that handles authenticating users using the handlers mechanism.

### Authenticate

An authentication scheme's authenticate action is responsible for constructing the user's identity based on request context. It returns an AuthenticateResult indicating whether authentication was successful and, if so, the user's identity in an authentication ticket. See AuthenticateAsync. Authenticate examples include:

- A cookie authentication scheme constructing the user's identity from cookies.
- A JWT bearer scheme deserializing and validating a JWT bearer token to construct the user's identity.

### Challenge

An authentication challenge is invoked by Authorization when an unauthenticated user requests an endpoint that requires authentication. An authentication challenge is issued, for example, when an anonymous user requests a restricted resource or follows a login link. Authorization invokes a challenge using the specified authentication schemes, or the default if none is specified. See ChallengeAsync. Authentication challenge examples include:

- A cookie authentication scheme redirecting the user to a login page.
- A JWT bearer scheme returning a 401 result with a `www-authenticate: bearer`header.

A challenge action should let the user know what authentication mechanism to use to access the requested resource.

### Forbid

An authentication scheme's forbid action is called by Authorization when an authenticated user attempts to access a resource they're not permitted to access. See ForbidAsync. Authentication forbid examples include:

- A cookie authentication scheme redirecting the user to a page indicating access was forbidden.
- A JWT bearer scheme returning a 403 result.
- A custom authentication scheme redirecting to a page where the user can request access to the resource.

A forbid action can let the user know:

- They're authenticated.
- They're not permitted to access the requested resource.

See the following links for differences between challenge and forbid:

## Authentication providers per tenant

ASP.NET Core framework doesn't have a built-in solution for multi-tenant authentication. While it's possible for customers to write an app with multi-tenant authentication, we recommend using one of the following ASP.NET Core application frameworks that support multi-tenant authentication.

Orchard Core is an open-source, modular, and multi-tenant app framework built with ASP.NET Core that also provides a content management system (CMS). See the Orchard Core source for an example of authentication providers per tenant.

ABP Framework supports various architectural patterns including modularity, microservices, domain-driven design, and multi-tenancy. See ABP Framework source on GitHub.

## Additional resources

ASP.NET Core

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Authorization refers to the process that determines what a user is able to do. For example, an administrative user is allowed to create a document library, add documents, edit documents, and delete them. A nonadministrative user working with the library is only authorized to read the documents.

Authorization is separate and distinct from authentication. However, authorization relies on an authentication mechanism. Authentication is the process of verifying a user's identity, which might result in the creation of one or more identity objects for the user.

For more information about authentication in ASP.NET Core, see Overview of ASP.NET Core Authentication.

## Authorization types

ASP.NET Core authorization provides a simple declarative role and a rich policy-based model. Authorization is expressed in requirements, and handlers evaluate a user's claims against requirements. Imperative checks can be based on simple policies or policies that evaluate both the user identity and properties of the resource that the user is attempting to access.

## Namespaces

Authorization components, including the `AuthorizeAttribute` and `AllowAnonymousAttribute` attributes, are defined in the `Microsoft.AspNetCore.Authorization` namespace.

Consult the documentation on simple authorization.

## Related content

ASP.NET Core

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

ASP.NET Core provides a cryptographic API to protect data, including key management and rotation.

Web apps often need to store sensitive data. The Windows data protection API (DPAPI) isn't intended for use in web apps.

The ASP.NET Core data protection stack is designed to:

- Provide a built-in solution for most web scenarios.
- Address many deficiencies of the previous encryption system.
- Serve as the replacement for the `<machineKey>`element in ASP.NET 1.x - 4.x.

## Problem statement

*I need to persist trusted information for later retrieval, but I don't trust the persistence mechanism.* In web terms, this statement might be written as *I need to round-trip trusted state via an untrusted client.*

Authenticity, integrity, and tamper-proofing are requirements. The canonical example of this scenario is an authentication cookie or bearer token. The server generates an * I am Groot and have xyz permissions* token and sends it to the client. The client presents that token back to the server, but the server needs some kind of assurance that the client didn't forge the token.

Confidentiality is a requirement. Because the persisted state is trusted by the server, this state might contain information that shouldn't be disclosed to an untrusted client. For example:

- A file path
- A permission
- A handle or other indirect reference
- Some server-specific data

Isolation is a requirement. Because modern apps are componentized, individual components want to take advantage of this system without regard to other components in the system. For instance, consider a bearer token component using this stack. It should operate without any interference, for example, from an anti-CSRF mechanism also using the same stack.

Some common assumptions can narrow the scope of requirements:

- All services operating within the cryptosystem are equally trusted.
- The data doesn't need to be generated or consumed outside of the services under our direct control.
- Operations must be fast because each request to the web service might go through the cryptosystem one or more times. The speed requirement makes symmetric cryptography ideal. Asymmetric cryptography isn't used until it's required.

## Design philosophy

ASP.NET Core data protection is an easy to use data protection stack based on the following principles:

- **Configuration should be easy**. The system strives for zero configuration. In situations where developers need to configure a specific aspect, such as the key repository, those specific configurations aren't difficult.
- **Offer basic consumer-facing APIs**. The APIs are straight forward to use correctly and difficult to use incorrectly.
- **Don't require the developer to learn the principles of managing keys**. The system handles algorithm selection and key lifetime on behalf of the developer. Developers don't have access to the raw key material, so they don't need expert knowledge of the principles.
- **Protect keys at rest as much as possible**. The system figures out an appropriate default protection mechanism and applies it automatically.

The data protection APIs aren't primarily intended for indefinite persistence of confidential payloads. Other technologies, such as Windows CNG DPAPI and Azure Rights Management are more suited to the scenario of indefinite storage. They have correspondingly strong key management capabilities. That said, the ASP.NET Core data protection APIs can be used for long-term protection of confidential data.

## Target audience

The data protection system provides APIs that target three main audiences:

- The consumer APIs target application and framework developers. - *I don't want to learn about how the stack operates or about how it's configured. I just want to perform some operation with high probability of using the APIs successfully.*
- The configuration APIs target app developers and system administrators. - *I need to tell the data protection system that my environment requires nondefault paths or settings.*
- The extensibility APIs target developers in charge of implementing custom policy. Usage of these APIs is limited to rare situations and developers with security experience. - *I need to replace an entire component within the system because I have truly unique behavioral requirements. I'm willing to learn uncommonly used parts of the API surface so I can build a plugin that fulfills my requirements.*

## Package layout

The data protection stack consists of five packages:

- Microsoft.AspNetCore.DataProtection.Abstractions contains: - IDataProtectionProvider and IDataProtector interfaces to create data protection services.
- Useful extension methods for working with these types, such as IDataProtector.Protect.
 - If the data protection system is instantiated elsewhere and you're consuming the API, reference - `Microsoft.AspNetCore.DataProtection.Abstractions`.
- Microsoft.AspNetCore.DataProtection contains the core implementation of the data protection system, including: - Core cryptographic operations
- Key management
- Configuration and extensibility
 - To instantiate the data protection system, reference - `Microsoft.AspNetCore.DataProtection`. You might need to reference the data protection system when:- Adding it to an IServiceCollection.
- Modifying or extending its behavior.

- Microsoft.AspNetCore.DataProtection.Extensions contains more APIs which developers might find useful but, which don't belong in the core package. For instance, this package contains: - Factory methods to instantiate the data protection system to store keys at a location on the file system without dependency injection. For more information, see DataProtectionProvider.
- Extension methods for limiting the lifetime of protected payloads. For more information, see ITimeLimitedDataProtector.

- Microsoft.AspNetCore.DataProtection.SystemWeb can be installed into an existing ASP.NET 4.x app to redirect its - `<machineKey>`operations to use the new ASP.NET Core data protection stack. For more information, see Replace the ASP.NET machineKey in ASP.NET Core.
- Microsoft.AspNetCore.Cryptography.KeyDerivation provides an implementation of the PBKDF2 password hashing routine. It's convenient for systems that must handle user passwords securely. For more information, see Hash passwords in ASP.NET Core.

## Related content

ASP.NET Core provides a cryptographic API to protect data, including key management and rotation.

Web apps often need to store sensitive data. The Windows data protection API (DPAPI) isn't intended for use in web apps.

The ASP.NET Core data protection stack was designed to:

- Provide a built in solution for most Web scenarios.
- Address many of the deficiencies of the previous encryption system.
- Serve as the replacement for the `<machineKey>`element in ASP.NET 1.x - 4.x.

## Problem statement

*I need to persist trusted information for later retrieval, but I don't trust the persistence mechanism.* In web terms, this might be written as *I need to round-trip trusted state via an untrusted client.*

Authenticity, integrity, and tamper-proofing is a requirement. The canonical example of this is an authentication cookie or bearer token. The server generates an * I am Groot and have xyz permissions* token and sends it to the client. The client presents that token back to the server, but the server needs some kind of assurance that the client hasn't forged the token.

Confidentiality is a requirement. Since the persisted state is trusted by the server, this state could contain information that shouldn't be disclosed to an untrusted client. For example:

- A file path.
- A permission.
- A handle or other indirect reference.
- Some server-specific data.

Isolation is a requirement. Since modern apps are componentized, individual components want to take advantage of this system without regard to other components in the system. For instance, consider a bearer token component using this stack. It should operate without any interference, for example, from an anti-CSRF mechanism also using the same stack.

Some common assumptions can narrow the scope of requirements:

- All services operating within the cryptosystem are equally trusted.
- The data doesn't need to be generated or consumed outside of the services under our direct control.
- Operations must be fast since each request to the web service might go through the cryptosystem one or more times. The speed requirement makes symmetric cryptography ideal. Asymmetric cryptography isn't used until it's required.

## Design philosophy

ASP.NET Core data protection is an easy to use data protection stack. It's based on the following principles:

- Ease of configuration. The system strives for zero configuration. In situations where developers need to configure a specific aspect, such as the key repository, those specific configurations aren't difficult.
- Offer a basic consumer-facing API. The APIs are straight forward to use correctly and difficult to use incorrectly.
- Developers don't have to learn key management principles. The system handles algorithm selection and key lifetime on behalf of the developer. The developer doesn't have access to the raw key material.
- Keys are protected at rest as much as possible. The system figures out an appropriate default protection mechanism and applies it automatically.

The data protection APIs aren't primarily intended for indefinite persistence of confidential payloads. Other technologies, such as Windows CNG DPAPI and Azure Rights Management are more suited to the scenario of indefinite storage. They have correspondingly strong key management capabilities. That said, the ASP.NET Core data protection APIs can be used for long-term protection of confidential data.

## Audience

The data protection system provides APIs that target three main audiences:

- The consumer APIs target application and framework developers. - *I don't want to learn about how the stack operates or about how it's configured. I just want to perform some operation with high probability of using the APIs successfully.*
- The configuration APIs target app developers and system administrators. - *I need to tell the data protection system that my environment requires non-default paths or settings.*
- The extensibility APIs target developers in charge of implementing custom policy. Usage of these APIs is limited to rare situations and developers with security experience. - *I need to replace an entire component within the system because I have truly unique behavioral requirements. I'm willing to learn uncommonly used parts of the API surface in order to build a plugin that fulfills my requirements.*

## Package layout

The data protection stack consists of five packages:

- Microsoft.AspNetCore.DataProtection.Abstractions contains: - IDataProtectionProvider and IDataProtector interfaces to create data protection services.
- Useful extension methods for working with these types. for example, IDataProtector.Protect
 - If the data protection system is instantiated elsewhere and you're consuming the API, reference - `Microsoft.AspNetCore.DataProtection.Abstractions`.
- Microsoft.AspNetCore.DataProtection contains the core implementation of the data protection system, including: - Core cryptographic operations.
- Key management.
- Configuration and extensibility.
 - To instantiate the data protection system, reference - `Microsoft.AspNetCore.DataProtection`. You might need to reference the data protection system when:- Adding it to an IServiceCollection.
- Modifying or extending its behavior.

- Microsoft.AspNetCore.DataProtection.Extensions contains additional APIs which developers might find useful but which don't belong in the core package. For instance, this package contains: - Factory methods to instantiate the data protection system to store keys at a location on the file system without dependency injection. See DataProtectionProvider.
- Extension methods for limiting the lifetime of protected payloads. See ITimeLimitedDataProtector.

- Microsoft.AspNetCore.DataProtection.SystemWeb can be installed into an existing ASP.NET 4.x app to redirect its - `<machineKey>`operations to use the new ASP.NET Core data protection stack. For more information, see Replace the ASP.NET machineKey in ASP.NET Core.
- Microsoft.AspNetCore.Cryptography.KeyDerivation provides an implementation of the PBKDF2 password hashing routine and can be used by systems that must handle user passwords securely. For more information, see Hash passwords in ASP.NET Core.

## Additional resources

ASP.NET Core

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

By David Galvan and Rick Anderson

To enforce incoming requests to your ASP.NET Core apps to use HTTPS/TLS, you can:

- Require HTTPS for all requests.
- Redirect all HTTP requests to HTTPS.

No API can prevent a client from sending sensitive data on the first request.

This article describes how to configure your ASP.NET Core apps to require HTTPS/TLS or redirect HTTP requests to HTTPS/TLS for secure interaction. Troubleshooting steps are provided for various platforms to resolve untrusted certificate issues.

## API projects

Projects that use Web APIs should either:

- Not listen on HTTP.
- Close the connection with status code 400 (Bad Request) and not serve the request.

To disable HTTP redirection in an API, set the `ASPNETCORE_URLS` environment variable or use the `--urls` command line flag. For more information, see ASP.NET Core runtime environments and 8 ways to set the URLs for an ASP.NET Core app by Andrew Lock.

Warning

Do **not** use RequireHttpsAttribute on Web APIs that receive sensitive information.
`RequireHttpsAttribute` uses HTTP status codes to redirect browsers from HTTP to HTTPS.
API clients might not understand or obey redirects from HTTP to HTTPS, and they might send information over HTTP.

### HSTS and API projects

The secure approach for HTTP Strict Transport Security (HSTS) protocol is to configure API projects to only listen to and respond over HTTPS.

Warning

The default API projects don't include HSTS because it's generally a browser only instruction. Other callers, such as phone or desktop apps, do **not** obey the instruction. Even within browsers, a single authenticated call to an API over HTTP has risks on insecure networks.

### HTTP redirect to HTTPS (ERR_INVALID_REDIRECT on CORS preflight request)

When a request to an endpoint using HTTP is redirected to HTTPS with the UseHttpsRedirection method, the redirection fails with the `ERR_INVALID_REDIRECT` error on the CORS preflight request.

API projects can reject HTTP requests rather than use the `UseHttpsRedirection` method to redirect requests to HTTPS.

## Require HTTPS

For production ASP.NET Core web apps, the following approach is recommended:

- To redirect HTTP requests to HTTPS, use HTTPS redirection middleware (UseHttpsRedirection).
- To send HSTS headers to clients, use HSTS middleware via the UseHsts method.

Note

Apps deployed in a reverse proxy configuration allow the proxy to handle connection security (HTTPS). If the proxy also handles HTTPS redirection, there's no need to use HTTPS redirection middleware. If the proxy server also handles writing HSTS headers (for example, native HSTS support in Internet Information Services (IIS) 10.0 version 1709 or later), then the app doesn't require HSTS middleware. For more information, see Opt-out of HTTPS/HSTS on project creation.

### HTTPS redirection middleware (`UseHttpsRedirection`)

The following code calls the UseHttpsRedirection method in the *Program.cs* file:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
The preceding highlighted code:

- Uses the default HttpsRedirectionOptions.RedirectStatusCode property with the Status307TemporaryRedirect code.
- Uses the default HttpsRedirectionOptions.HttpsPort property (passing null), unless overridden by the `ASPNETCORE_HTTPS_PORT`environment variable or IServerAddressesFeature.

The recommended approach is to use temporary redirects rather than permanent redirects. Link caching can cause unstable behavior in development environments. If you prefer to send a permanent redirect status code when the app is in a non-`Development` environment, see the Configure permanent redirects in production section. Use HSTS to signal to clients that only secure resource requests should be sent to the app (only in production).

Note

Don't confuse the `HTTPS_PORT` configuration key and `ASPNETCORE_HTTPS_PORT` environment variable, which set the port for HTTPS redirection middleware, with the `HTTPS_PORTS` configuration key and `ASPNETCORE_HTTPS_PORTS` environment variable, which set the ports for Kestrel/HTTP.sys endpoint configuration.

### Port configuration

A port must be available for the middleware to redirect an insecure request to HTTPS. If no port is available:

- Redirection to HTTPS doesn't occur.
- The middleware logs the warning *Failed to determine the https port for redirect*.

Specify the HTTPS port by using any of the following approaches:

- Set the - `https_port`host setting:- In host configuration.
- By setting the - `ASPNETCORE_HTTPS_PORT`environment variable.
- By adding a top-level entry in the - *appsettings.json*file:- `{ "https_port": 443, "Logging": { "LogLevel": { "Default": "Information", "Microsoft.AspNetCore": "Warning" } }, "AllowedHosts": "*" }`

- Indicate a port with the secure scheme by using the ASPNETCORE_URLS environment variable. The environment variable configures the server. The middleware indirectly discovers the HTTPS port via IServerAddressesFeature. This approach doesn't work in reverse proxy deployments.
- The ASP.NET Core web templates set an HTTPS URL in the - *Properties/launchsettings.json*file for both Kestrel and IIS Express. The- *launchsettings.json*file is used on the local machine only.
- Configure an HTTPS URL endpoint for a public-facing edge deployment of Kestrel server or HTTP.sys server. Only - **one HTTPS port**is used by the app. The middleware discovers the port via IServerAddressesFeature.

Note

When an app runs in a reverse proxy configuration, IServerAddressesFeature isn't available. Set the port by using one of the other approaches described in this section.

### Edge deployments

When Kestrel or HTTP.sys is used as a public-facing edge server, Kestrel or HTTP.sys must be configured to listen on both:

- The secure port where the client is redirected (typically, 443 in production and 5001 in development).
- The insecure port (typically, 80 in production and 5000 in development).

The insecure port must be accessible by the client for the app to receive an insecure request and redirect the client to the secure port.

For more information, see Kestrel endpoint configuration or HTTP.sys web server implementation in ASP.NET Core.

### Deployment scenarios

Any firewall between the client and server must also have communication ports open for traffic.

If requests are forwarded in a reverse proxy configuration, use forwarded headers middleware before calling HTTPS redirection middleware. Forwarded headers middleware updates the `Request.Scheme` by using the `X-Forwarded-Proto` header. The middleware permits redirect URIs and other security policies to work correctly. When forwarded headers middleware isn't used, the back-end app might not receive the correct scheme and get caught in a redirect loop. A common end user error message is there are too many redirects.

When deploying to Azure App Service, follow the guidance in Enable HTTPS for a custom domain in Azure App Service.

### Options

The following highlighted code calls the AddHttpsRedirection method to configure middleware options:

```
using static Microsoft.AspNetCore.Http.StatusCodes;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
builder.Services.AddHsts(options =>
{
 options.Preload = true;
 options.IncludeSubDomains = true;
 options.MaxAge = TimeSpan.FromDays(60);
 options.ExcludedHosts.Add("example.com");
 options.ExcludedHosts.Add("www.example.com");
});
builder.Services.AddHttpsRedirection(options =>
{
 options.RedirectStatusCode = Status307TemporaryRedirect;
 options.HttpsPort = 5001;
});
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
Calling `AddHttpsRedirection` is only necessary to change the values of `HttpsPort` or `RedirectStatusCode`.

The preceding highlighted code:

- Sets the HttpsRedirectionOptions.RedirectStatusCode property to the Status307TemporaryRedirect code, which is the default value. Use the fields of the StatusCodes class for assignments to `RedirectStatusCode`.
- Sets the HTTPS port to 5001.

#### Configure permanent redirects in production

The middleware defaults to sending a Status307TemporaryRedirect code with all redirects. If you prefer to send a permanent redirect status code when the app is in a non-`Development` environment, wrap the middleware options configuration in a conditional check for a non-`Development` environment.

The following code shows configuration of services in the *Program.cs* file:

```
using static Microsoft.AspNetCore.Http.StatusCodes;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
if (!builder.Environment.IsDevelopment())
{
 builder.Services.AddHttpsRedirection(options =>
 {
 options.RedirectStatusCode = Status308PermanentRedirect;
 options.HttpsPort = 443;
 });
}
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
## HTTPS redirection middleware alternative approach

An alternative to using HTTPS redirection middleware (with the `UseHttpsRedirection` method) is to use URL rewriting middleware (via the `AddRedirectToHttps` method). `AddRedirectToHttps` can also set the status code and port when the redirect is executed. For more information, see URL rewriting middleware.

When the app redirects to HTTPS without the requirement for other redirect rules, the recommendation is to use HTTPS redirection middleware (`UseHttpsRedirection`) as described in this article.

## HTTP Strict Transport Security (HSTS) protocol

Per OWASP, HSTS is an opt-in security enhancement specified by a web app via a response header. When a browser that supports HSTS receives this header:

- The browser stores configuration for the domain that prevents sending any communication over HTTP. The browser forces all communication over HTTPS.
- The browser prevents the user from using untrusted or invalid certificates. The browser disables prompts that allow a user to temporarily trust such a certificate.

Because the client enforces HSTS, there are some limitations:

- The client must support HSTS.
- HSTS requires at least one successful HTTPS request to establish the HSTS policy.
- The application must check every HTTP request and redirect or reject the HTTP request.

ASP.NET Core implements HSTS with the UseHsts extension method. The following code calls `UseHsts` when the app isn't in development mode:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
`UseHsts` isn't recommended in development because the HSTS settings are highly cacheable by browsers. By default, `UseHsts` excludes the local loopback address.

For production environments that are implementing HTTPS for the first time, set the initial HstsOptions.MaxAge property value to a small amount by using one of the TimeSpan methods. Set the value from hours to no more than a single day, in case you need to revert the HTTPS infrastructure to HTTP. After you're confident in the sustainability of the HTTPS configuration, increase the HSTS `max-age` value (commonly, one year).

The following highlighted code:

```
using static Microsoft.AspNetCore.Http.StatusCodes;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
builder.Services.AddHsts(options =>
{
 options.Preload = true;
 options.IncludeSubDomains = true;
 options.MaxAge = TimeSpan.FromDays(60);
 options.ExcludedHosts.Add("example.com");
 options.ExcludedHosts.Add("www.example.com");
});
builder.Services.AddHttpsRedirection(options =>
{
 options.RedirectStatusCode = Status307TemporaryRedirect;
 options.HttpsPort = 5001;
});
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
- Sets the preload parameter of the `Strict-Transport-Security`header. Preload isn't part of the RFC 6797 HSTS specification. Web browsers support preload of HSTS sites on fresh install. For more information, see https://hstspreload.org/.
- Enables the `includeSubDomain`directive, which applies the HSTS policy to host subdomains. For more information, see RFC 6797 HSTS specification (Section 6.1.2).
- Explicitly sets the `max-age`parameter of the`Strict-Transport-Security`header to 60 days. If not set, it defaults to 30 days. For more information, see the`max-age`directive in RFC 6797 HSTS specification (Section 6.1.1).
- Adds `example.com`to the list of hosts to exclude.

`UseHsts` excludes the following loopback hosts:

- `localhost`: The IPv4 loopback address.
- `127.0.0.1`: The IPv4 loopback address.
- `[::1]`: The IPv6 loopback address.

## Opt out of HTTPS/HSTS on project creation

In some back-end service scenarios where connection security is handled at the public-facing edge of the network, configuring connection security at each node isn't required. Web apps that are generated from the templates in Visual Studio or from the dotnet new command enable HTTPS redirection and HSTS. For deployments that don't require these scenarios, you can opt out of HTTPS/HSTS when the app is created from the template.

To opt out of HTTPS/HSTS:

When you create a new ASP.NET Core web app, unselect the **Configure for HTTPS** option:

## Trust the ASP.NET Core HTTPS development certificate

The .NET SDK includes an HTTPS development certificate. The certificate is installed as part of the first-run experience. For example, the `dotnet --info` command produces a variation of the following output:

```
ASP.NET Core
------------
Successfully installed the ASP.NET Core HTTPS Development Certificate.
To trust the certificate run 'dotnet dev-certs https --trust' (Windows and macOS only).
For establishing trust on other platforms refer to the platform specific documentation.
For more information on configuring HTTPS see https://go.microsoft.com/fwlink/?linkid=848054.
```
Installing the .NET SDK installs the ASP.NET Core HTTPS development certificate to the local user certificate store. The certificate is installed, but it isn't trusted. To trust the certificate, perform the one-time step to run the `dotnet dev-certs` tool:

```
dotnet dev-certs https --trust
```
The following command provides help on the `dotnet dev-certs` tool:

```
dotnet dev-certs https --help
```
Warning

Don't create a development certificate in an environment planned for redistribution, such as a container image or virtual machine. This scenario can lead to spoofing and elevation of privilege. To help prevent this situation, set the `DOTNET_GENERATE_ASPNET_CERTIFICATE` environment variable to `false` before calling the .NET CLI for the first time. This approach skips the automatic generation of the ASP.NET Core development certificate during the CLI's first-run experience.

## Set up developer certificate for Docker

To configure the developer certificate for Docker, see GitHub dotnet/aspnetcore.docs issue #6199 - *How to set up the dev certificate when using Docker in development*.

## Linux-specific considerations

Linux distributions differ substantially in how they mark certificates as trusted.

The `dotnet dev-certs` tool is expected to be broadly applicable, but official support is available for Ubuntu and Fedora only. The support specifically aims to ensure trust in Firefox and Chromium-based browsers (Microsoft Edge, Chrome, and Chromium).

### Dependencies

- To establish OpenSSL trust, the `openssl`tool must be on the path.
- To establish browser trust (for example in Microsoft Edge or Firefox), the `certutil`tool must be on the path.

### OpenSSL trust

When an ASP.NET Core development certificate is trusted, the certificate is exported to a folder in the current user's home directory. To have OpenSSL (and clients that consume it) pick up this folder, you need to set the `SSL_CERT_DIR` environment variable. You can set the variable in a single session by running a command like `export SSL_CERT_DIR=$HOME/.aspnet/dev-certs/trust:/usr/lib/ssl/certs` (the exact value is in the output when `--verbose` is passed) or by adding it your (distro- and shell-specific) configuration file (for example *.profile*).

This approach is required to make tools like `curl` trust the development certificate. Alternatively, you can pass `-CAfile` or `-CApath` to each individual `curl` invocation.

Note

Requires 1.1.1h or later or 3.0.0 or later, depending on which major version you're using.

If OpenSSL trust gets into a bad state (for example if `dotnet dev-certs https --clean` fails to remove it), you can frequently resolve the situation by using the c_rehash tool.

### Overrides

If you're using another browser with its own Network Security Services (NSS) store, you can use the `DOTNET_DEV_CERTS_NSSDB_PATHS` environment variable to specify a colon-delimited list of NSS directories (for example, the directory containing `cert9.db`). You can then add the development certificate location to the list in the variable.

If you store the certificates you want OpenSSL to trust in a specific directory, you can use the `DOTNET_DEV_CERTS_OPENSSL_CERTIFICATE_DIRECTORY` environment variable to indicate the certificate location.

Warning

If you set either variable, be sure to set the same values each time trust is updated. If the values change, the tool doesn't know about certificates in the former locations (for example, during certificate cleanup).

### Using sudo

As on other platforms, development certificates are stored and trusted separately for each user.

If you run `dotnet dev-certs` as a different user (for example, by using `sudo`), then *that* specific user (for example `root`) trusts the development certificate.

### Trust HTTPS certificate on Linux with linux-dev-certs

linux-dev-certs is an open-source, community-supported, .NET global tool that provides a convenient way to create and trust a developer certificate on Linux. Microsoft doesn't maintain or support the tool.

The following commands install the tool and create a trusted developer certificate:

```
dotnet tool update -g linux-dev-certs
dotnet linux-dev-certs install
```
For more information or to report issues, see the linux-dev-certs GitHub repository.

#### SUSE Linux Enterprise Server (SLES Linux)

If your configuration includes SUSE Linux Enterprise Server, see GitHub dotnet/aspnetcore.docs issue #28292 - *Trust HTTPS certificate on SLES*.

## Troubleshoot certificate problems (certificate not trusted)

Sometimes when an ASP.NET Core HTTPS development certificate is installed and trusted, the browser warns that the certificate is untrusted. The following sections provide help for troubleshooting this issue.

The ASP.NET Core HTTPS development certificate is used by Kestrel.

To repair the IIS Express certificate, see Stack Overflow issue #20036984 / answer #20048613 - *How do I restore a missing IIS Express SSL Certificate?*

### All platforms - certificate not trusted

For all platforms, try to resolve the untrusted certificate issues with the following steps:

- Run the following commands: - `dotnet dev-certs https --clean dotnet dev-certs https --trust`
- Close any open browser instances, and open the app in a new browser window. - The browser cache stores whether a certificate is trusted. The close/open process helps to refresh the browser cache settings for certificates.

The `dotnet dev-certs https` commands usually solve most browser trust issues. If the `dotnet dev-certs https --clean` command fails and the browser still doesn't trust the certificate, try the platform-specific suggestions in the following sections.

### Docker - certificate not trusted

If you're using Docker, try to resolve the issue with the following steps:

- Delete the - *C:\Users{USER}\AppData\Roaming\ASP.NET\Https*folder.
- Clean the solution. Delete the - *bin*and- *obj*folders.
- Restart the development tool. For example, Visual Studio or Visual Studio Code.

### Windows - certificate not trusted

If you're working in Windows, complete the following troubleshooting steps:

- Check the certificates in the certificate store. Look for a - `localhost`certificate with the- `ASP.NET Core HTTPS development certificate`friendly name in two folders:- *Current User > Personal > Certificates*
- *Current User > Trusted root certification authorities > Certificates*

- Remove all certificates from both Personal and Trusted root certification authorities. - Important - Do - **not**remove the IIS Express localhost certificate.
- Run the following commands: - `dotnet dev-certs https --clean dotnet dev-certs https --trust`
- Close any open browser instances, and open the app in a new browser window.

### OS X - certificate not trusted

If you're working with OS X, try to resolve the issue with the following steps:

- Open KeyChain Access, and then select the System keychain.
- Check for the presence of a localhost certificate.
- Confirm the certificate shows the plus ( - `+`) symbol on the icon, which indicates the certificate is trusted for all users.
- Remove the certificate from the system keychain.
- Run the following commands: - `dotnet dev-certs https --clean dotnet dev-certs https --trust`
- Close any open browser instances, and open the app in a new browser window.

For more information about troubleshooting certificate issues with Visual Studio, see GitHub dotnet/aspnetcore issue #16892) - *HTTPS Error using IIS Express*.

### Linux - certificate not trusted

If you're running Linux, follow these steps to troubleshoot the untrusted certificate:

- Confirm the certificate you're investigating is the user HTTPS developer certificate planned for use by the Kestrel server.
- Check the current user default HTTPS developer Kestrel certificate at the following location: - `ls -la ~/.dotnet/corefx/cryptography/x509stores/my`- The HTTPS developer Kestrel certificate file is the SHA1 thumbprint. When the file is deleted with the - `dotnet dev-certs https --clean`command, the file is regenerated when needed with a different thumbprint.
- Verify the thumbprint of the exported certificate matches by running the following command: - `openssl x509 -noout -fingerprint -sha1 -inform pem -in /usr/local/share/ca-certificates/aspnet/https.crt`- If the certificate thumbprint doesn't match, investigate the following conditions: - Check if the certificate is old.
- Check if the certificate is an exported developer certificate for the root user. - If it is, export the certificate.

- Check the root user certificate in the following folder: - `ls -la /root/.dotnet/corefx/cryptography/x509stores/my`

### IIS Express SSL certificate used with Visual Studio

To fix problems with the IIS Express certificate, select **Repair** in the Visual Studio installer. For more information, see GitHub dotnet/aspnetcore issue #16892) - *HTTPS Error using IIS Express*.

### Group policy prevents trusting self-signed certificates

In some cases, group policy can prevent self-signed certificates from being trusted. For more information, see GitHub dotnet/aspnetcore issue #21173 - *Error trusting HTTPS developer certificate*.

## Related content

Note

If you're using .NET 9 or later SDK, see the updated Linux procedures in the .NET 9 version of this article.

Warning

## API projects

Do **not** use RequireHttpsAttribute on Web APIs that receive sensitive information. `RequireHttpsAttribute` uses HTTP status codes to redirect browsers from HTTP to HTTPS. API clients may not understand or obey redirects from HTTP to HTTPS. Such clients may send information over HTTP. Web APIs should either:

- Not listen on HTTP.
- Close the connection with status code 400 (Bad Request) and not serve the request.

To disable HTTP redirection in an API, set the `ASPNETCORE_URLS` environment variable or use the `--urls` command line flag. For more information, see ASP.NET Core runtime environments and 8 ways to set the URLs for an ASP.NET Core app by Andrew Lock.

## HSTS and API projects

The default API projects don't include HSTS because HSTS is generally a browser only instruction. Other callers, such as phone or desktop apps, do **not** obey the instruction. Even within browsers, a single authenticated call to an API over HTTP has risks on insecure networks. The secure approach is to configure API projects to only listen to and respond over HTTPS.

### HTTP redirection to HTTPS causes ERR_INVALID_REDIRECT on the CORS preflight request

Requests to an endpoint using HTTP that are redirected to HTTPS by UseHttpsRedirection fail with `ERR_INVALID_REDIRECT` on the CORS preflight request.

API projects can reject HTTP requests rather than use `UseHttpsRedirection` to redirect requests to HTTPS.

## Require HTTPS

We recommend that production ASP.NET Core web apps use:

- HTTPS redirection middleware (UseHttpsRedirection) to redirect HTTP requests to HTTPS.
- HSTS middleware (UseHsts) to send HTTP Strict Transport Security (HSTS) protocol headers to clients.

Note

Apps deployed in a reverse proxy configuration allow the proxy to handle connection security (HTTPS). If the proxy also handles HTTPS redirection, there's no need to use HTTPS redirection middleware. If the proxy server also handles writing HSTS headers (for example, native HSTS support in IIS 10.0 (1709) or later), HSTS middleware isn't required by the app. For more information, see Opt-out of HTTPS/HSTS on project creation.

### HTTPS redirection middleware (`UseHttpsRedirection`)

The following code calls UseHttpsRedirection in the `Program.cs` file:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
The preceding highlighted code:

- Uses the default HttpsRedirectionOptions.RedirectStatusCode (Status307TemporaryRedirect).
- Uses the default HttpsRedirectionOptions.HttpsPort (null) unless overridden by the `ASPNETCORE_HTTPS_PORT`environment variable or IServerAddressesFeature.

We recommend using temporary redirects rather than permanent redirects. Link caching can cause unstable behavior in development environments. If you prefer to send a permanent redirect status code when the app is in a non-`Development` environment, see the Configure permanent redirects in production section. We recommend using HSTS to signal to clients that only secure resource requests should be sent to the app (only in production).

### Port configuration

A port must be available for the middleware to redirect an insecure request to HTTPS. If no port is available:

- Redirection to HTTPS doesn't occur.
- The middleware logs the warning "Failed to determine the https port for redirect."

Specify the HTTPS port using any of the following approaches:

- Set the - `https_port`host setting:- In host configuration.
- By setting the - `ASPNETCORE_HTTPS_PORT`environment variable.
- By adding a top-level entry in - `appsettings.json`:- `{ "https_port": 443, "Logging": { "LogLevel": { "Default": "Information", "Microsoft.AspNetCore": "Warning" } }, "AllowedHosts": "*" }`

- Indicate a port with the secure scheme using the ASPNETCORE_URLS environment variable. The environment variable configures the server. The middleware indirectly discovers the HTTPS port via IServerAddressesFeature. This approach doesn't work in reverse proxy deployments.
- The ASP.NET Core web templates set an HTTPS URL in - `Properties/launchsettings.json`for both Kestrel and IIS Express.- `launchsettings.json`is only used on the local machine.
- Configure an HTTPS URL endpoint for a public-facing edge deployment of Kestrel server or HTTP.sys server. Only - **one HTTPS port**is used by the app. The middleware discovers the port via IServerAddressesFeature.

Note

When an app is run in a reverse proxy configuration, IServerAddressesFeature isn't available. Set the port using one of the other approaches described in this section.

### Edge deployments

When Kestrel or HTTP.sys is used as a public-facing edge server, Kestrel or HTTP.sys must be configured to listen on both:

- The secure port where the client is redirected (typically, 443 in production and 5001 in development).
- The insecure port (typically, 80 in production and 5000 in development).

The insecure port must be accessible by the client in order for the app to receive an insecure request and redirect the client to the secure port.

For more information, see Kestrel endpoint configuration or HTTP.sys web server implementation in ASP.NET Core.

### Deployment scenarios

Any firewall between the client and server must also have communication ports open for traffic.

If requests are forwarded in a reverse proxy configuration, use forwarded headers middleware before calling HTTPS redirection middleware. Forwarded headers middleware updates the `Request.Scheme`, using the `X-Forwarded-Proto` header. The middleware permits redirect URIs and other security policies to work correctly. When forwarded headers middleware isn't used, the backend app might not receive the correct scheme and end up in a redirect loop. A common end user error message is that too many redirects have occurred.

When deploying to Azure App Service, follow the guidance in Tutorial: Bind an existing custom SSL certificate to Azure Web Apps.

### Options

The following highlighted code calls AddHttpsRedirection to configure middleware options:

```
using static Microsoft.AspNetCore.Http.StatusCodes;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
builder.Services.AddHsts(options =>
{
 options.Preload = true;
 options.IncludeSubDomains = true;
 options.MaxAge = TimeSpan.FromDays(60);
 options.ExcludedHosts.Add("example.com");
 options.ExcludedHosts.Add("www.example.com");
});
builder.Services.AddHttpsRedirection(options =>
{
 options.RedirectStatusCode = Status307TemporaryRedirect;
 options.HttpsPort = 5001;
});
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
Calling `AddHttpsRedirection` is only necessary to change the values of `HttpsPort` or `RedirectStatusCode`.

The preceding highlighted code:

- Sets HttpsRedirectionOptions.RedirectStatusCode to Status307TemporaryRedirect, which is the default value. Use the fields of the StatusCodes class for assignments to `RedirectStatusCode`.
- Sets the HTTPS port to 5001.

#### Configure permanent redirects in production

The middleware defaults to sending a Status307TemporaryRedirect with all redirects. If you prefer to send a permanent redirect status code when the app is in a non-`Development` environment, wrap the middleware options configuration in a conditional check for a non-`Development` environment.

When configuring services in `Program.cs`:

```
using static Microsoft.AspNetCore.Http.StatusCodes;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
if (!builder.Environment.IsDevelopment())
{
 builder.Services.AddHttpsRedirection(options =>
 {
 options.RedirectStatusCode = Status308PermanentRedirect;
 options.HttpsPort = 443;
 });
}
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
## HTTPS redirection middleware alternative approach

An alternative to using HTTPS redirection middleware (`UseHttpsRedirection`) is to use URL rewriting middleware (`AddRedirectToHttps`). `AddRedirectToHttps` can also set the status code and port when the redirect is executed. For more information, see URL rewriting middleware.

When redirecting to HTTPS without the requirement for additional redirect rules, we recommend using HTTPS redirection middleware (`UseHttpsRedirection`) described in this article.

## HTTP Strict Transport Security (HSTS) protocol

Per OWASP, HTTP Strict Transport Security (HSTS) is an opt-in security enhancement that's specified by a web app through the use of a response header. When a browser that supports HSTS receives this header:

- The browser stores configuration for the domain that prevents sending any communication over HTTP. The browser forces all communication over HTTPS.
- The browser prevents the user from using untrusted or invalid certificates. The browser disables prompts that allow a user to temporarily trust such a certificate.

Because HSTS is enforced by the client, it has some limitations:

- The client must support HSTS.
- HSTS requires at least one successful HTTPS request to establish the HSTS policy.
- The application must check every HTTP request and redirect or reject the HTTP request.

ASP.NET Core implements HSTS with the UseHsts extension method. The following code calls `UseHsts` when the app isn't in development mode:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
`UseHsts` isn't recommended in development because the HSTS settings are highly cacheable by browsers. By default, `UseHsts` excludes the local loopback address.

For production environments that are implementing HTTPS for the first time, set the initial HstsOptions.MaxAge to a small value using one of the TimeSpan methods. Set the value from hours to no more than a single day in case you need to revert the HTTPS infrastructure to HTTP. After you're confident in the sustainability of the HTTPS configuration, increase the HSTS `max-age` value; a commonly used value is one year.

The following highlighted code:

```
using static Microsoft.AspNetCore.Http.StatusCodes;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
builder.Services.AddHsts(options =>
{
 options.Preload = true;
 options.IncludeSubDomains = true;
 options.MaxAge = TimeSpan.FromDays(60);
 options.ExcludedHosts.Add("example.com");
 options.ExcludedHosts.Add("www.example.com");
});
builder.Services.AddHttpsRedirection(options =>
{
 options.RedirectStatusCode = Status307TemporaryRedirect;
 options.HttpsPort = 5001;
});
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
- Sets the preload parameter of the `Strict-Transport-Security`header. Preload isn't part of the RFC HSTS specification, but is supported by web browsers to preload HSTS sites on fresh install. For more information, see https://hstspreload.org/.
- Enables includeSubDomain, which applies the HSTS policy to Host subdomains.
- Explicitly sets the `max-age`parameter of the`Strict-Transport-Security`header to 60 days. If not set, defaults to 30 days. For more information, see the max-age directive.
- Adds `example.com`to the list of hosts to exclude.

`UseHsts` excludes the following loopback hosts:

- `localhost`: The IPv4 loopback address.
- `127.0.0.1`: The IPv4 loopback address.
- `[::1]`: The IPv6 loopback address.

## Opt-out of HTTPS/HSTS on project creation

In some backend service scenarios where connection security is handled at the public-facing edge of the network, configuring connection security at each node isn't required. Web apps that are generated from the templates in Visual Studio or from the dotnet new command enable HTTPS redirection and HSTS. For deployments that don't require these scenarios, you can opt-out of HTTPS/HSTS when the app is created from the template.

To opt-out of HTTPS/HSTS:

Uncheck the **Configure for HTTPS** checkbox.

## Trust the ASP.NET Core HTTPS development certificate on Windows and macOS

For the Firefox browser, see the next section.

The .NET Core SDK includes an HTTPS development certificate. The certificate is installed as part of the first-run experience. For example, `dotnet --info` produces a variation of the following output:

```
ASP.NET Core
------------
Successfully installed the ASP.NET Core HTTPS Development Certificate.
To trust the certificate run 'dotnet dev-certs https --trust' (Windows and macOS only).
For establishing trust on other platforms refer to the platform specific documentation.
For more information on configuring HTTPS see https://go.microsoft.com/fwlink/?linkid=848054.
```
Installing the .NET Core SDK installs the ASP.NET Core HTTPS development certificate to the local user certificate store. The certificate has been installed, but it's not trusted. To trust the certificate, perform the one-time step to run the `dotnet dev-certs` tool:

```
dotnet dev-certs https --trust
```
The following command provides help on the `dotnet dev-certs` tool:

```
dotnet dev-certs https --help
```
Warning

Do not create a development certificate in an environment that will be redistributed, such as a container image or virtual machine. Doing so can lead to spoofing and elevation of privilege. To help prevent this, set the `DOTNET_GENERATE_ASPNET_CERTIFICATE` environment variable to `false` prior to calling the .NET CLI for the first time. This will skip the automatic generation of the ASP.NET Core development certificate during the CLI's first-run experience.

### Trust the HTTPS certificate with Firefox to prevent SEC_ERROR_INADEQUATE_KEY_USAGE error

The Firefox browser uses its own certificate store, and therefore doesn't trust the IIS Express or Kestrel developer certificates.

There are two approaches to trusting the HTTPS certificate with Firefox, create a policy file or configure with the FireFox browser. Configuring with the browser creates the policy file, so the two approaches are equivalent.

#### Create a policy file to trust HTTPS certificate with Firefox

Create a policy file (`policies.json`) at:

- Windows: `%PROGRAMFILES%\Mozilla Firefox\distribution\`
- MacOS: `Firefox.app/Contents/Resources/distribution`
- Linux: See Trust the certificate with Firefox on Linux in this article.

Add the following JSON to the Firefox policy file:

```
{
 "policies": {
 "Certificates": {
 "ImportEnterpriseRoots": true
 }
 }
}
```
The preceding policy file makes Firefox trust certificates from the trusted certificates in the Windows certificate store. The next section provides an alternative approach to create the preceding policy file by using the Firefox browser.

### Configure trust of HTTPS certificate using Firefox browser

Set `security.enterprise_roots.enabled` = `true` using the following instructions:

- Enter `about:config`in the FireFox browser.
- Select **Accept the Risk and Continue**if you accept the risk.
- Select **Show All**
- Set `security.enterprise_roots.enabled`=`true`
- Exit and restart Firefox

For more information, see Setting Up Certificate Authorities (CAs) in Firefox and the mozilla/policy-templates/README file.

## How to set up a developer certificate for Docker

See this GitHub issue.

## Trust HTTPS certificate on Linux

Establishing trust is distribution and browser specific. The following sections provide instructions for some popular distributions and the Chromium browsers (Edge and Chrome) and for Firefox.

### Ubuntu trust the certificate for service-to-service communication

The following instructions don't work for some Ubuntu versions, such as 20.04. For more information, see GitHub issue dotnet/AspNetCore.Docs #23686.

- Install OpenSSL 1.1.1h or later. See your distribution for instructions on how to update OpenSSL.
- Run the following commands: - `dotnet dev-certs https sudo -E dotnet dev-certs https -ep /usr/local/share/ca-certificates/aspnet/https.crt --format PEM sudo update-ca-certificates`

The preceding commands:

- Ensure the current user's developer certificate is created.
- Exports the certificate with elevated permissions needed for the `ca-certificates`folder, using the current user's environment.
- Removing the `-E`flag exports the root user certificate, generating it if necessary. Each newly generated certificate has a different thumbprint. When running as root,`sudo`and`-E`are not needed.

The path in the preceding command is specific for Ubuntu. For other distributions, select an appropriate path or use the path for the Certificate Authorities (CAs).

### Trust HTTPS certificate on Linux using Edge or Chrome

For chromium browsers on Linux:

- Install the - `libnss3-tools`for your distribution.
- Create or verify the - `$HOME/.pki/nssdb`folder exists on the machine.
- Export the certificate with the following command: - `dotnet dev-certs https sudo -E dotnet dev-certs https -ep /usr/local/share/ca-certificates/aspnet/https.crt --format PEM`- The path in the preceding command is specific for Ubuntu. For other distributions, select an appropriate path or use the path for the Certificate Authorities (CAs).
- Run the following commands: - `certutil -d sql:$HOME/.pki/nssdb -A -t "P,," -n localhost -i /usr/local/share/ca-certificates/aspnet/https.crt`
- Exit and restart the browser.

#### Trust the certificate with Firefox on Linux

- Export the certificate with the following command: - `dotnet dev-certs https sudo -E dotnet dev-certs https -ep /usr/local/share/ca-certificates/aspnet/https.crt --format PEM`- The path in the preceding command is specific for Ubuntu. For other distributions, select an appropriate path or use the path for the Certificate Authorities (CAs).
- Create a JSON file at - `/usr/lib/firefox/distribution/policies.json`with the following command:

```
cat <<EOF | sudo tee /usr/lib/firefox/distribution/policies.json
{
 "policies": {
 "Certificates": {
 "Install": [
 "/usr/local/share/ca-certificates/aspnet/https.crt"
 ]
 }
 }
}
EOF
```
Note: Ubuntu 21.10 Firefox comes as a snap package and the installation folder is `/snap/firefox/current/usr/lib/firefox`.

See Configure trust of HTTPS certificate using Firefox browser in this article for an alternative way to configure the policy file using the browser.

### Trust the certificate with Fedora 34

See:

- This GitHub comment
- Fedora: Using Shared System Certificates
- Set up a .NET development environment on Fedora.

### Trust the certificate with other distros

See this GitHub issue.

## Trust HTTPS certificate from Windows Subsystem for Linux

The following instructions don't work for some Linux distributions, such as Ubuntu 20.04. For more information, see GitHub issue dotnet/AspNetCore.Docs #23686.

The Windows Subsystem for Linux (WSL) generates an HTTPS self-signed development certificate, which by default isn't trusted in Windows. The easiest way to have Windows trust the WSL certificate, is to configure WSL to use the same certificate as Windows:

- On - **Windows**- `dotnet dev-certs https -ep https.pfx -p $CREDENTIAL_PLACEHOLDER$ --trust`- Where - `$CREDENTIAL_PLACEHOLDER$`is a password.
- In a WSL window, import the exported certificate on the WSL instance: - `dotnet dev-certs https --clean --import <<path-to-pfx>> --password $CREDENTIAL_PLACEHOLDER$`

The preceding approach is a one time operation per certificate and per WSL distribution. It's easier than exporting the certificate over and over. If you update or regenerate the certificate on windows, you might need to run the preceding commands again.

## Troubleshoot certificate problems such as certificate not trusted

This section provides help when the ASP.NET Core HTTPS development certificate has been installed and trusted, but you still have browser warnings that the certificate is not trusted. The ASP.NET Core HTTPS development certificate is used by Kestrel.

To repair the IIS Express certificate, see this Stackoverflow issue.

### All platforms - certificate not trusted

Run the following commands:

```
dotnet dev-certs https --clean
dotnet dev-certs https --trust
```
Close any browser instances open. Open a new browser window to app. Certificate trust is cached by browsers.

### dotnet dev-certs https --clean Fails

The preceding commands solve most browser trust issues. If the browser is still not trusting the certificate, follow the platform-specific suggestions that follow.

### Docker - certificate not trusted

- Delete the *C:\Users{USER}\AppData\Roaming\ASP.NET\Https*folder.
- Clean the solution. Delete the *bin*and*obj*folders.
- Restart the development tool. For example, Visual Studio or Visual Studio Code.

### Windows - certificate not trusted

- Check the certificates in the certificate store. There should be a `localhost`certificate with the`ASP.NET Core HTTPS development certificate`friendly name both under`Current User > Personal > Certificates`and`Current User > Trusted root certification authorities > Certificates`
- Remove all the found certificates from both Personal and Trusted root certification authorities. Do **not**remove the IIS Express localhost certificate.
- Run the following commands:

```
dotnet dev-certs https --clean
dotnet dev-certs https --trust
```
Close any browser instances open. Open a new browser window to app.

### OS X - certificate not trusted

- Open KeyChain Access.
- Select the System keychain.
- Check for the presence of a localhost certificate.
- Check that it contains a `+`symbol on the icon to indicate it's trusted for all users.
- Remove the certificate from the system keychain.
- Run the following commands:

```
dotnet dev-certs https --clean
dotnet dev-certs https --trust
```
Close any browser instances open. Open a new browser window to app.

See HTTPS Error using IIS Express (dotnet/AspNetCore #16892) for troubleshooting certificate issues with Visual Studio.

### Linux certificate not trusted

Check that the certificate being configured for trust is the user HTTPS developer certificate that will be used by the Kestrel server.

Check the current user default HTTPS developer Kestrel certificate at the following location:

```
ls -la ~/.dotnet/corefx/cryptography/x509stores/my
```
The HTTPS developer Kestrel certificate file is the SHA1 thumbprint. When the file is deleted via `dotnet dev-certs https --clean`, it's regenerated when needed with a different thumbprint.
Check the thumbprint of the exported certificate matches with the following command:

```
openssl x509 -noout -fingerprint -sha1 -inform pem -in /usr/local/share/ca-certificates/aspnet/https.crt
```
If the certificate doesn't match, it could be one of the following:

- An old certificate.
- An exported a developer certificate for the root user. For this case, export the certificate.

The root user certificate can be checked at:

```
ls -la /root/.dotnet/corefx/cryptography/x509stores/my
```
### IIS Express SSL certificate used with Visual Studio

To fix problems with the IIS Express certificate, select **Repair** from the Visual Studio installer. For more information, see this GitHub issue.

### Group policy prevents self-signed certificates from being trusted

In some cases, group policy may prevent self-signed certificates from being trusted. For more information, see this GitHub issue.

## Additional information

Warning

## API projects

Do **not** use RequireHttpsAttribute on Web APIs that receive sensitive information. `RequireHttpsAttribute` uses HTTP status codes to redirect browsers from HTTP to HTTPS. API clients may not understand or obey redirects from HTTP to HTTPS. Such clients may send information over HTTP. Web APIs should either:

- Not listen on HTTP.
- Close the connection with status code 400 (Bad Request) and not serve the request.

To disable HTTP redirection in an API, set the `ASPNETCORE_URLS` environment variable or use the `--urls` command line flag. For more information, see ASP.NET Core runtime environments and 5 ways to set the URLs for an ASP.NET Core app by Andrew Lock.

## HSTS and API projects

The default API projects don't include HSTS because HSTS is generally a browser only instruction. Other callers, such as phone or desktop apps, do **not** obey the instruction. Even within browsers, a single authenticated call to an API over HTTP has risks on insecure networks. The secure approach is to configure API projects to only listen to and respond over HTTPS.

## Require HTTPS

We recommend that production ASP.NET Core web apps use:

- HTTPS redirection middleware (UseHttpsRedirection) to redirect HTTP requests to HTTPS.
- HSTS middleware (UseHsts) to send HTTP Strict Transport Security (HSTS) protocol headers to clients.

Note

Apps deployed in a reverse proxy configuration allow the proxy to handle connection security (HTTPS). If the proxy also handles HTTPS redirection, there's no need to use HTTPS redirection middleware. If the proxy server also handles writing HSTS headers (for example, native HSTS support in IIS 10.0 (1709) or later), HSTS middleware isn't required by the app. For more information, see Opt-out of HTTPS/HSTS on project creation.

### UseHttpsRedirection

The following code calls `UseHttpsRedirection` in the `Startup` class:

```
public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
{
 if (env.IsDevelopment())
 {
 app.UseDeveloperExceptionPage();
 }
 else
 {
 app.UseExceptionHandler("/Error");
 // The default HSTS value is 30 days. You may want to change this for production scenarios, see https://aka.ms/aspnetcore-hsts.
 app.UseHsts();
 }
 app.UseHttpsRedirection();
 app.UseStaticFiles();
 app.UseRouting();
 app.UseAuthorization();
 app.UseEndpoints(endpoints =>
 {
 endpoints.MapRazorPages();
 });
}
```
The preceding highlighted code:

- Uses the default HttpsRedirectionOptions.RedirectStatusCode (Status307TemporaryRedirect).
- Uses the default HttpsRedirectionOptions.HttpsPort (null) unless overridden by the `ASPNETCORE_HTTPS_PORT`environment variable or IServerAddressesFeature.

We recommend using temporary redirects rather than permanent redirects. Link caching can cause unstable behavior in development environments. If you prefer to send a permanent redirect status code when the app is in a non-`Development` environment, see the Configure permanent redirects in production section. We recommend using HSTS to signal to clients that only secure resource requests should be sent to the app (only in production).

### Port configuration

A port must be available for the middleware to redirect an insecure request to HTTPS. If no port is available:

- Redirection to HTTPS doesn't occur.
- The middleware logs the warning "Failed to determine the https port for redirect."

Specify the HTTPS port using any of the following approaches:

- Set the - `https_port`host setting:- In host configuration.
- By setting the - `ASPNETCORE_HTTPS_PORT`environment variable.
- By adding a top-level entry in - `appsettings.json`:- `{ "https_port": 443, "Logging": { "LogLevel": { "Default": "Information", "Microsoft": "Warning", "Microsoft.Hosting.Lifetime": "Information" } }, "AllowedHosts": "*" }`

- Indicate a port with the secure scheme using the ASPNETCORE_URLS environment variable. The environment variable configures the server. The middleware indirectly discovers the HTTPS port via IServerAddressesFeature. This approach doesn't work in reverse proxy deployments.
- In development, set an HTTPS URL in - `launchsettings.json`. Enable HTTPS when IIS Express is used.
- Configure an HTTPS URL endpoint for a public-facing edge deployment of Kestrel server or HTTP.sys server. Only - **one HTTPS port**is used by the app. The middleware discovers the port via IServerAddressesFeature.

Note

When an app is run in a reverse proxy configuration, IServerAddressesFeature isn't available. Set the port using one of the other approaches described in this section.

### Edge deployments

When Kestrel or HTTP.sys is used as a public-facing edge server, Kestrel or HTTP.sys must be configured to listen on both:

- The secure port where the client is redirected (typically, 443 in production and 5001 in development).
- The insecure port (typically, 80 in production and 5000 in development).

The insecure port must be accessible by the client in order for the app to receive an insecure request and redirect the client to the secure port.

For more information, see Kestrel endpoint configuration or HTTP.sys web server implementation in ASP.NET Core.

### Deployment scenarios

Any firewall between the client and server must also have communication ports open for traffic.

If requests are forwarded in a reverse proxy configuration, use forwarded headers middleware before calling HTTPS redirection middleware. Forwarded headers middleware updates the `Request.Scheme`, using the `X-Forwarded-Proto` header. The middleware permits redirect URIs and other security policies to work correctly. When forwarded headers middleware isn't used, the backend app might not receive the correct scheme and end up in a redirect loop. A common end user error message is that too many redirects have occurred.

When deploying to Azure App Service, follow the guidance in Tutorial: Bind an existing custom SSL certificate to Azure Web Apps.

### Options

The following highlighted code calls AddHttpsRedirection to configure middleware options:

```
public void ConfigureServices(IServiceCollection services)
{
 services.AddRazorPages();
 services.AddHsts(options =>
 {
 options.Preload = true;
 options.IncludeSubDomains = true;
 options.MaxAge = TimeSpan.FromDays(60);
 options.ExcludedHosts.Add("example.com");
 options.ExcludedHosts.Add("www.example.com");
 });
 services.AddHttpsRedirection(options =>
 {
 options.RedirectStatusCode = (int) HttpStatusCode.TemporaryRedirect;
 options.HttpsPort = 5001;
 });
}
```
Calling `AddHttpsRedirection` is only necessary to change the values of `HttpsPort` or `RedirectStatusCode`.

The preceding highlighted code:

- Sets HttpsRedirectionOptions.RedirectStatusCode to Status307TemporaryRedirect, which is the default value. Use the fields of the StatusCodes class for assignments to `RedirectStatusCode`.
- Sets the HTTPS port to 5001.

#### Configure permanent redirects in production

The middleware defaults to sending a Status307TemporaryRedirect with all redirects. If you prefer to send a permanent redirect status code when the app is in a non-`Development` environment, wrap the middleware options configuration in a conditional check for a non-`Development` environment.

When configuring services in `Startup.cs`:

```
public void ConfigureServices(IServiceCollection services)
{
 // IWebHostEnvironment (stored in _env) is injected into the Startup class.
 if (!_env.IsDevelopment())
 {
 services.AddHttpsRedirection(options =>
 {
 options.RedirectStatusCode = (int) HttpStatusCode.PermanentRedirect;
 options.HttpsPort = 443;
 });
 }
}
```
## HTTPS redirection middleware alternative approach

An alternative to using HTTPS redirection middleware (`UseHttpsRedirection`) is to use URL rewriting middleware (`AddRedirectToHttps`). `AddRedirectToHttps` can also set the status code and port when the redirect is executed. For more information, see URL rewriting middleware.

When redirecting to HTTPS without the requirement for additional redirect rules, we recommend using HTTPS redirection middleware (`UseHttpsRedirection`) described in this article.

## HTTP Strict Transport Security (HSTS) protocol

Per OWASP, HTTP Strict Transport Security (HSTS) is an opt-in security enhancement that's specified by a web app through the use of a response header. When a browser that supports HSTS receives this header:

- The browser stores configuration for the domain that prevents sending any communication over HTTP. The browser forces all communication over HTTPS.
- The browser prevents the user from using untrusted or invalid certificates. The browser disables prompts that allow a user to temporarily trust such a certificate.

Because HSTS is enforced by the client, it has some limitations:

- The client must support HSTS.
- HSTS requires at least one successful HTTPS request to establish the HSTS policy.
- The application must check every HTTP request and redirect or reject the HTTP request.

ASP.NET Core implements HSTS with the `UseHsts` extension method. The following code calls `UseHsts` when the app isn't in development mode:

```
public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
{
 if (env.IsDevelopment())
 {
 app.UseDeveloperExceptionPage();
 }
 else
 {
 app.UseExceptionHandler("/Error");
 // The default HSTS value is 30 days. You may want to change this for production scenarios, see https://aka.ms/aspnetcore-hsts.
 app.UseHsts();
 }
 app.UseHttpsRedirection();
 app.UseStaticFiles();
 app.UseRouting();
 app.UseAuthorization();
 app.UseEndpoints(endpoints =>
 {
 endpoints.MapRazorPages();
 });
}
```
`UseHsts` isn't recommended in development because the HSTS settings are highly cacheable by browsers. By default, `UseHsts` excludes the local loopback address.

For production environments that are implementing HTTPS for the first time, set the initial HstsOptions.MaxAge to a small value using one of the TimeSpan methods. Set the value from hours to no more than a single day in case you need to revert the HTTPS infrastructure to HTTP. After you're confident in the sustainability of the HTTPS configuration, increase the HSTS `max-age` value; a commonly used value is one year.

The following code:

```
public void ConfigureServices(IServiceCollection services)
{
 services.AddRazorPages();
 services.AddHsts(options =>
 {
 options.Preload = true;
 options.IncludeSubDomains = true;
 options.MaxAge = TimeSpan.FromDays(60);
 options.ExcludedHosts.Add("example.com");
 options.ExcludedHosts.Add("www.example.com");
 });
 services.AddHttpsRedirection(options =>
 {
 options.RedirectStatusCode = (int) HttpStatusCode.TemporaryRedirect;
 options.HttpsPort = 5001;
 });
}
```
- Sets the preload parameter of the `Strict-Transport-Security`header. Preload isn't part of the RFC HSTS specification, but is supported by web browsers to preload HSTS sites on fresh install. For more information, see https://hstspreload.org/.
- Enables includeSubDomain, which applies the HSTS policy to Host subdomains.
- Explicitly sets the `max-age`parameter of the`Strict-Transport-Security`header to 60 days. If not set, defaults to 30 days. For more information, see the max-age directive.
- Adds `example.com`to the list of hosts to exclude.

`UseHsts` excludes the following loopback hosts:

- `localhost`: The IPv4 loopback address.
- `127.0.0.1`: The IPv4 loopback address.
- `[::1]`: The IPv6 loopback address.

## Opt-out of HTTPS/HSTS on project creation

In some backend service scenarios where connection security is handled at the public-facing edge of the network, configuring connection security at each node isn't required. Web apps that are generated from the templates in Visual Studio or from the dotnet new command enable HTTPS redirection and HSTS. For deployments that don't require these scenarios, you can opt-out of HTTPS/HSTS when the app is created from the template.

To opt-out of HTTPS/HSTS:

Uncheck the **Configure for HTTPS** checkbox.

## Trust the ASP.NET Core HTTPS development certificate on Windows and macOS

For the Firefox browser, see the next section.

The .NET Core SDK includes an HTTPS development certificate. The certificate is installed as part of the first-run experience. For example, running `dotnet new webapp` for the first time produces a variation of the following output:

```
Installed an ASP.NET Core HTTPS development certificate.
To trust the certificate, run 'dotnet dev-certs https --trust'
Learn about HTTPS: https://aka.ms/dotnet-https
```
Installing the .NET Core SDK installs the ASP.NET Core HTTPS development certificate to the local user certificate store. The certificate has been installed, but it's not trusted. To trust the certificate, perform the one-time step to run the `dotnet dev-certs` tool:

```
dotnet dev-certs https --trust
```
The following command provides help on the `dotnet dev-certs` tool:

```
dotnet dev-certs https --help
```
Warning

Do not create a development certificate in an environment that will be redistributed, such as a container image or virtual machine. Doing so can lead to spoofing and elevation of privilege. To help prevent this, set the `DOTNET_GENERATE_ASPNET_CERTIFICATE` environment variable to `false` prior to calling the .NET CLI for the first time. This will skip the automatic generation of the ASP.NET Core development certificate during the CLI's first-run experience.

### Trust the HTTPS certificate with Firefox to prevent SEC_ERROR_INADEQUATE_KEY_USAGE error

The Firefox browser uses its own certificate store, and therefore doesn't trust the IIS Express or Kestrel developer certificates.

There are two approaches to trusting the HTTPS certificate with Firefox, create a policy file or configure with the FireFox browser. Configuring with the browser creates the policy file, so the two approaches are equivalent.

#### Create a policy file to trust HTTPS certificate with Firefox

Create a policy file (`policies.json`) at:

- Windows: `%PROGRAMFILES%\Mozilla Firefox\distribution\`
- MacOS: `Firefox.app/Contents/Resources/distribution`
- Linux: See Trust the certificate with Firefox on Linux later in this article.

Add the following JSON to the Firefox policy file:

```
{
 "policies": {
 "Certificates": {
 "ImportEnterpriseRoots": true
 }
 }
}
```
The preceding policy file makes Firefox trust certificates from the trusted certificates in the Windows certificate store. The next section provides an alternative approach to create the preceding policy file by using the Firefox browser.

### Configure trust of HTTPS certificate using Firefox browser

Set `security.enterprise_roots.enabled` = `true` using the following instructions:

- Enter `about:config`in the FireFox browser.
- Select **Accept the Risk and Continue**if you accept the risk.
- Select **Show All**.
- Set `security.enterprise_roots.enabled`=`true`.
- Exit and restart Firefox.

For more information, see Setting Up Certificate Authorities (CAs) in Firefox and the mozilla/policy-templates/README file.

## How to set up a developer certificate for Docker

See this GitHub issue.

## Trust HTTPS certificate on Linux

Establishing trust is distribution and browser specific. The following sections provide instructions for some popular distributions and the Chromium browsers (Edge and Chrome) and for Firefox.

### Ubuntu trust the certificate for service-to-service communication

- Install OpenSSL 1.1.1h or later. See your distribution for instructions on how to update OpenSSL.
- Run the following commands: - `dotnet dev-certs https sudo -E dotnet dev-certs https -ep /usr/local/share/ca-certificates/aspnet/https.crt --format PEM sudo update-ca-certificates`

The preceding commands:

- Ensure the current user's developer certificate is created.
- Export the certificate with elevated permissions needed for the `ca-certificates`folder, using the current user's environment.
- Remove the `-E`flag to export the root user certificate, generating it if necessary. Each newly generated certificate has a different thumbprint. When running as root,`sudo`and`-E`are not needed.

The path in the preceding command is specific for Ubuntu. For other distributions, select an appropriate path or use the path for the Certificate Authorities (CAs).

### Trust HTTPS certificate on Linux using Edge or Chrome

For chromium browsers on Linux:

- Install the - `libnss3-tools`for your distribution.
- Create or verify the - `$HOME/.pki/nssdb`folder exists on the machine.
- Export the certificate with the following command: - `dotnet dev-certs https sudo -E dotnet dev-certs https -ep /usr/local/share/ca-certificates/aspnet/https.crt --format PEM`- The path in the preceding command is specific for Ubuntu. For other distributions, select an appropriate path or use the path for the Certificate Authorities (CAs).
- Run the following commands: - `certutil -d sql:$HOME/.pki/nssdb -A -t "P,," -n localhost -i /usr/local/share/ca-certificates/aspnet/https.crt`
- Exit and restart the browser.

### Trust the certificate with Firefox on Linux

- Export the certificate with the following command: - `dotnet dev-certs https sudo -E dotnet dev-certs https -ep /usr/local/share/ca-certificates/aspnet/https.crt --format PEM`- The path in the preceding command is specific for Ubuntu. For other distributions, select an appropriate path or use the path for the Certificate Authorities (CAs).
- Create a JSON file at - `/usr/lib/firefox/distribution/policies.json`with the following contents:

```
cat <<EOF | sudo tee /usr/lib/firefox/distribution/policies.json
{
 "policies": {
 "Certificates": {
 "Install": [
 "/usr/local/share/ca-certificates/aspnet/https.crt"
 ]
 }
 }
}
EOF
```
See Configure trust of HTTPS certificate using Firefox browser in this article for an alternative way to configure the policy file using the browser.

### Trust the certificate with Fedora 34

#### Firefox on Fedora

```
echo 'pref("general.config.filename", "firefox.cfg");
pref("general.config.obscure_value", 0);' > ./autoconfig.js
echo '//Enable policies.json
lockPref("browser.policies.perUserDir", false);' > firefox.cfg
echo "{
 \"policies\": {
 \"Certificates\": {
 \"Install\": [
 \"aspnetcore-localhost-https.crt\"
 ]
 }
 }
}" > policies.json
dotnet dev-certs https -ep localhost.crt --format PEM
sudo mv autoconfig.js /usr/lib64/firefox/
sudo mv firefox.cfg /usr/lib64/firefox/
sudo mv policies.json /usr/lib64/firefox/distribution/
mkdir -p ~/.mozilla/certificates
cp localhost.crt ~/.mozilla/certificates/aspnetcore-localhost-https.crt
rm localhost.crt
```
#### Trust dotnet-to-dotnet on Fedora

```
sudo cp localhost.crt /etc/pki/tls/certs/localhost.pem
sudo update-ca-trust
rm localhost.crt
```
See this GitHub comment for more information.

### Trust the certificate with other distros

See this GitHub issue.

## Trust HTTPS certificate from Windows Subsystem for Linux

The Windows Subsystem for Linux (WSL) generates an HTTPS self-signed development certificate. To configure the Windows certificate store to trust the WSL certificate:

- Export the developer certificate to a file on - **Windows**- `dotnet dev-certs https -ep C:\<<path-to-folder>>\aspnetcore.pfx -p $CREDENTIAL_PLACEHOLDER$`- Where - `$CREDENTIAL_PLACEHOLDER$`is a password.
- In a WSL window, import the exported certificate on the WSL instance: - `dotnet dev-certs https --clean --import /mnt/c/<<path-to-folder>>/aspnetcore.pfx -p $CREDENTIAL_PLACEHOLDER$`

The preceding approach is a one time operation per certificate and per WSL distribution. It's easier than exporting the certificate over and over. If you update or regenerate the certificate on windows, you might need to run the preceding commands again.

## Troubleshoot certificate problems such as certificate not trusted

This section provides help when the ASP.NET Core HTTPS development certificate has been installed and trusted, but you still have browser warnings that the certificate is not trusted. The ASP.NET Core HTTPS development certificate is used by Kestrel.

To repair the IIS Express certificate, see this Stackoverflow issue.

### All platforms - certificate not trusted

Run the following commands:

```
dotnet dev-certs https --clean
dotnet dev-certs https --trust
```
Close any browser instances that are open. Open a new browser window to the app. Certificate trust is cached by browsers.

### dotnet dev-certs https --clean fails

The preceding commands solve most browser trust issues. If the browser is still not trusting the certificate, follow the platform-specific suggestions that follow.

### Docker - certificate not trusted

- Delete the *C:\Users{USER}\AppData\Roaming\ASP.NET\Https*folder.
- Clean the solution. Delete the *bin*and*obj*folders.
- Restart the development tool. For example, Visual Studio, Visual Studio Code, or Visual Studio for Mac.

### Windows - certificate not trusted

- Check the certificates in the certificate store. There should be a `localhost`certificate with the`ASP.NET Core HTTPS development certificate`friendly name both under`Current User > Personal > Certificates`and`Current User > Trusted root certification authorities > Certificates`
- Remove all the found certificates from both Personal and Trusted root certification authorities. Do **not**remove the IIS Express localhost certificate.
- Run the following commands:

```
dotnet dev-certs https --clean
dotnet dev-certs https --trust
```
Close any browser instances that are open. Open a new browser window to the app. Certificate trust is cached by browsers.

### OS X - certificate not trusted

- Open KeyChain Access.
- Select the System keychain.
- Check for the presence of a localhost certificate.
- Check that it contains a `+`symbol on the icon to indicate it's trusted for all users.
- Remove the certificate from the system keychain.
- Run the following commands:

```
dotnet dev-certs https --clean
dotnet dev-certs https --trust
```
Close any browser instances that are open. Open a new browser window to the app. Certificate trust is cached by browsers.

See HTTPS Error using IIS Express (dotnet/AspNetCore #16892) for troubleshooting certificate issues with Visual Studio.

### Linux certificate not trusted

Check that the certificate being configured for trust is the user HTTPS developer certificate that will be used by the Kestrel server.

Check the current user default HTTPS developer Kestrel certificate at the following location:

```
ls -la ~/.dotnet/corefx/cryptography/x509stores/my
```
The HTTPS developer Kestrel certificate file is the SHA1 thumbprint. When the file is deleted via `dotnet dev-certs https --clean`, it's regenerated when needed with a different thumbprint.
Check the thumbprint of the exported certificate matches with the following command:

```
openssl x509 -noout -fingerprint -sha1 -inform pem -in /usr/local/share/ca-certificates/aspnet/https.crt
```
If the certificate doesn't match, it could be one of the following:

- An old certificate.
- An exported a developer certificate for the root user. For this case, export the certificate.

The root user certificate can be checked at:

```
ls -la /root/.dotnet/corefx/cryptography/x509stores/my
```
### IIS Express SSL certificate used with Visual Studio

To fix problems with the IIS Express certificate, select **Repair** from the Visual Studio installer. For more information, see this GitHub issue.

## Additional information

Note

If you're using .NET 9 or later SDK, see the updated Linux procedures in the .NET 9 version of this article.

Warning

## API projects

Do **not** use RequireHttpsAttribute on Web APIs that receive sensitive information. `RequireHttpsAttribute` uses HTTP status codes to redirect browsers from HTTP to HTTPS. API clients may not understand or obey redirects from HTTP to HTTPS. Such clients may send information over HTTP. Web APIs should either:

- Not listen on HTTP.
- Close the connection with status code 400 (Bad Request) and not serve the request.

To disable HTTP redirection in an API, set the `ASPNETCORE_URLS` environment variable or use the `--urls` command line flag. For more information, see ASP.NET Core runtime environments and 8 ways to set the URLs for an ASP.NET Core app by Andrew Lock.

## HSTS and API projects

The default API projects don't include HSTS because HSTS is generally a browser only instruction. Other callers, such as phone or desktop apps, do **not** obey the instruction. Even within browsers, a single authenticated call to an API over HTTP has risks on insecure networks. The secure approach is to configure API projects to only listen to and respond over HTTPS.

### HTTP redirection to HTTPS causes ERR_INVALID_REDIRECT on the CORS preflight request

Requests to an endpoint using HTTP that are redirected to HTTPS by UseHttpsRedirection fail with `ERR_INVALID_REDIRECT` on the CORS preflight request.

API projects can reject HTTP requests rather than use `UseHttpsRedirection` to redirect requests to HTTPS.

## Require HTTPS

We recommend that production ASP.NET Core web apps use:

- HTTPS redirection middleware (UseHttpsRedirection) to redirect HTTP requests to HTTPS.
- HSTS middleware (UseHsts) to send HTTP Strict Transport Security (HSTS) protocol headers to clients.

Note

Apps deployed in a reverse proxy configuration allow the proxy to handle connection security (HTTPS). If the proxy also handles HTTPS redirection, there's no need to use HTTPS redirection middleware. If the proxy server also handles writing HSTS headers (for example, native HSTS support in IIS 10.0 (1709) or later), HSTS middleware isn't required by the app. For more information, see Opt-out of HTTPS/HSTS on project creation.

### HTTPS redirection middleware (`UseHttpsRedirection`)

The following code calls UseHttpsRedirection in the `Program.cs` file:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
The preceding highlighted code:

- Uses the default HttpsRedirectionOptions.RedirectStatusCode (Status307TemporaryRedirect).
- Uses the default HttpsRedirectionOptions.HttpsPort (null) unless overridden by the `ASPNETCORE_HTTPS_PORT`environment variable or IServerAddressesFeature.

We recommend using temporary redirects rather than permanent redirects. Link caching can cause unstable behavior in development environments. If you prefer to send a permanent redirect status code when the app is in a non-`Development` environment, see the Configure permanent redirects in production section. We recommend using HSTS to signal to clients that only secure resource requests should be sent to the app (only in production).

### Port configuration

A port must be available for the middleware to redirect an insecure request to HTTPS. If no port is available:

- Redirection to HTTPS doesn't occur.
- The middleware logs the warning "Failed to determine the https port for redirect."

Specify the HTTPS port using any of the following approaches:

- Set the - `https_port`host setting:- In host configuration.
- By setting the - `ASPNETCORE_HTTPS_PORT`environment variable.
- By adding a top-level entry in - `appsettings.json`:- `{ "https_port": 443, "Logging": { "LogLevel": { "Default": "Information", "Microsoft.AspNetCore": "Warning" } }, "AllowedHosts": "*" }`

- Indicate a port with the secure scheme using the ASPNETCORE_URLS environment variable. The environment variable configures the server. The middleware indirectly discovers the HTTPS port via IServerAddressesFeature. This approach doesn't work in reverse proxy deployments.
- The ASP.NET Core web templates set an HTTPS URL in - `Properties/launchsettings.json`for both Kestrel and IIS Express.- `launchsettings.json`is only used on the local machine.
- Configure an HTTPS URL endpoint for a public-facing edge deployment of Kestrel server or HTTP.sys server. Only - **one HTTPS port**is used by the app. The middleware discovers the port via IServerAddressesFeature.

Note

When an app is run in a reverse proxy configuration, IServerAddressesFeature isn't available. Set the port using one of the other approaches described in this section.

### Edge deployments

When Kestrel or HTTP.sys is used as a public-facing edge server, Kestrel or HTTP.sys must be configured to listen on both:

- The secure port where the client is redirected (typically, 443 in production and 5001 in development).
- The insecure port (typically, 80 in production and 5000 in development).

The insecure port must be accessible by the client in order for the app to receive an insecure request and redirect the client to the secure port.

For more information, see Kestrel endpoint configuration or HTTP.sys web server implementation in ASP.NET Core.

### Deployment scenarios

Any firewall between the client and server must also have communication ports open for traffic.

If requests are forwarded in a reverse proxy configuration, use forwarded headers middleware before calling HTTPS redirection middleware. Forwarded headers middleware updates the `Request.Scheme`, using the `X-Forwarded-Proto` header. The middleware permits redirect URIs and other security policies to work correctly. When forwarded headers middleware isn't used, the backend app might not receive the correct scheme and end up in a redirect loop. A common end user error message is that too many redirects have occurred.

When deploying to Azure App Service, follow the guidance in Tutorial: Bind an existing custom SSL certificate to Azure Web Apps.

### Options

The following highlighted code calls AddHttpsRedirection to configure middleware options:

```
using static Microsoft.AspNetCore.Http.StatusCodes;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
builder.Services.AddHsts(options =>
{
 options.Preload = true;
 options.IncludeSubDomains = true;
 options.MaxAge = TimeSpan.FromDays(60);
 options.ExcludedHosts.Add("example.com");
 options.ExcludedHosts.Add("www.example.com");
});
builder.Services.AddHttpsRedirection(options =>
{
 options.RedirectStatusCode = Status307TemporaryRedirect;
 options.HttpsPort = 5001;
});
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
Calling `AddHttpsRedirection` is only necessary to change the values of `HttpsPort` or `RedirectStatusCode`.

The preceding highlighted code:

- Sets HttpsRedirectionOptions.RedirectStatusCode to Status307TemporaryRedirect, which is the default value. Use the fields of the StatusCodes class for assignments to `RedirectStatusCode`.
- Sets the HTTPS port to 5001.

#### Configure permanent redirects in production

The middleware defaults to sending a Status307TemporaryRedirect with all redirects. If you prefer to send a permanent redirect status code when the app is in a non-`Development` environment, wrap the middleware options configuration in a conditional check for a non-`Development` environment.

When configuring services in `Program.cs`:

```
using static Microsoft.AspNetCore.Http.StatusCodes;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
if (!builder.Environment.IsDevelopment())
{
 builder.Services.AddHttpsRedirection(options =>
 {
 options.RedirectStatusCode = Status308PermanentRedirect;
 options.HttpsPort = 443;
 });
}
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
## HTTPS redirection middleware alternative approach

An alternative to using HTTPS redirection middleware (`UseHttpsRedirection`) is to use URL rewriting middleware (`AddRedirectToHttps`). `AddRedirectToHttps` can also set the status code and port when the redirect is executed. For more information, see URL rewriting middleware.

When redirecting to HTTPS without the requirement for additional redirect rules, we recommend using HTTPS redirection middleware (`UseHttpsRedirection`) described in this article.

## HTTP Strict Transport Security (HSTS) protocol

Per OWASP, HTTP Strict Transport Security (HSTS) is an opt-in security enhancement that's specified by a web app through the use of a response header. When a browser that supports HSTS receives this header:

- The browser stores configuration for the domain that prevents sending any communication over HTTP. The browser forces all communication over HTTPS.
- The browser prevents the user from using untrusted or invalid certificates. The browser disables prompts that allow a user to temporarily trust such a certificate.

Because HSTS is enforced by the client, it has some limitations:

- The client must support HSTS.
- HSTS requires at least one successful HTTPS request to establish the HSTS policy.
- The application must check every HTTP request and redirect or reject the HTTP request.

ASP.NET Core implements HSTS with the UseHsts extension method. The following code calls `UseHsts` when the app isn't in development mode:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
`UseHsts` isn't recommended in development because the HSTS settings are highly cacheable by browsers. By default, `UseHsts` excludes the local loopback address.

For production environments that are implementing HTTPS for the first time, set the initial HstsOptions.MaxAge to a small value using one of the TimeSpan methods. Set the value from hours to no more than a single day in case you need to revert the HTTPS infrastructure to HTTP. After you're confident in the sustainability of the HTTPS configuration, increase the HSTS `max-age` value; a commonly used value is one year.

The following highlighted code:

```
using static Microsoft.AspNetCore.Http.StatusCodes;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
builder.Services.AddHsts(options =>
{
 options.Preload = true;
 options.IncludeSubDomains = true;
 options.MaxAge = TimeSpan.FromDays(60);
 options.ExcludedHosts.Add("example.com");
 options.ExcludedHosts.Add("www.example.com");
});
builder.Services.AddHttpsRedirection(options =>
{
 options.RedirectStatusCode = Status307TemporaryRedirect;
 options.HttpsPort = 5001;
});
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
- Sets the preload parameter of the `Strict-Transport-Security`header. Preload isn't part of the RFC HSTS specification, but is supported by web browsers to preload HSTS sites on fresh install. For more information, see https://hstspreload.org/.
- Enables includeSubDomain, which applies the HSTS policy to Host subdomains.
- Explicitly sets the `max-age`parameter of the`Strict-Transport-Security`header to 60 days. If not set, defaults to 30 days. For more information, see the max-age directive.
- Adds `example.com`to the list of hosts to exclude.

`UseHsts` excludes the following loopback hosts:

- `localhost`: The IPv4 loopback address.
- `127.0.0.1`: The IPv4 loopback address.
- `[::1]`: The IPv6 loopback address.

## Opt-out of HTTPS/HSTS on project creation

In some backend service scenarios where connection security is handled at the public-facing edge of the network, configuring connection security at each node isn't required. Web apps that are generated from the templates in Visual Studio or from the dotnet new command enable HTTPS redirection and HSTS. For deployments that don't require these scenarios, you can opt-out of HTTPS/HSTS when the app is created from the template.

To opt-out of HTTPS/HSTS:

Uncheck the **Configure for HTTPS** checkbox.

## Trust the ASP.NET Core HTTPS development certificate on Windows and macOS

For the Firefox browser, see the next section.

The .NET Core SDK includes an HTTPS development certificate. The certificate is installed as part of the first-run experience. For example, `dotnet --info` produces a variation of the following output:

```
ASP.NET Core
------------
Successfully installed the ASP.NET Core HTTPS Development Certificate.
To trust the certificate run 'dotnet dev-certs https --trust' (Windows and macOS only).
For establishing trust on other platforms refer to the platform specific documentation.
For more information on configuring HTTPS see https://go.microsoft.com/fwlink/?linkid=848054.
```
Installing the .NET Core SDK installs the ASP.NET Core HTTPS development certificate to the local user certificate store. The certificate has been installed, but it's not trusted. To trust the certificate, perform the one-time step to run the `dotnet dev-certs` tool:

```
dotnet dev-certs https --trust
```
The following command provides help on the `dotnet dev-certs` tool:

```
dotnet dev-certs https --help
```
Warning

Do not create a development certificate in an environment that will be redistributed, such as a container image or virtual machine. Doing so can lead to spoofing and elevation of privilege. To help prevent this, set the `DOTNET_GENERATE_ASPNET_CERTIFICATE` environment variable to `false` prior to calling the .NET CLI for the first time. This will skip the automatic generation of the ASP.NET Core development certificate during the CLI's first-run experience.

### Trust the HTTPS certificate with Firefox to prevent SEC_ERROR_INADEQUATE_KEY_USAGE error

The Firefox browser uses its own certificate store, and therefore doesn't trust the IIS Express or Kestrel developer certificates.

There are two approaches to trusting the HTTPS certificate with Firefox, create a policy file or configure with the FireFox browser. Configuring with the browser creates the policy file, so the two approaches are equivalent.

#### Create a policy file to trust HTTPS certificate with Firefox

Create a policy file (`policies.json`) at:

- Windows: `%PROGRAMFILES%\Mozilla Firefox\distribution\`
- MacOS: `Firefox.app/Contents/Resources/distribution`
- Linux: See Trust the certificate with Firefox on Linux in this article.

Add the following JSON to the Firefox policy file:

```
{
 "policies": {
 "Certificates": {
 "ImportEnterpriseRoots": true
 }
 }
}
```
The preceding policy file makes Firefox trust certificates from the trusted certificates in the Windows certificate store. The next section provides an alternative approach to create the preceding policy file by using the Firefox browser.

### Configure trust of HTTPS certificate using Firefox browser

Set `security.enterprise_roots.enabled` = `true` using the following instructions:

- Enter `about:config`in the FireFox browser.
- Select **Accept the Risk and Continue**if you accept the risk.
- Select **Show All**
- Set `security.enterprise_roots.enabled`=`true`
- Exit and restart Firefox

For more information, see Setting Up Certificate Authorities (CAs) in Firefox and the mozilla/policy-templates/README file.

## How to set up a developer certificate for Docker

See this GitHub issue.

## Trust HTTPS certificate on Linux

Establishing trust is distribution and browser specific. The following sections provide instructions for some popular distributions and the Chromium browsers (Edge and Chrome) and for Firefox.

### Ubuntu trust the certificate for service-to-service communication

The following instructions don't work for some Ubuntu versions, such as 20.04. For more information, see GitHub issue dotnet/AspNetCore.Docs #23686.

- Install OpenSSL 1.1.1h or later. See your distribution for instructions on how to update OpenSSL.
- Run the following commands: - `dotnet dev-certs https sudo -E dotnet dev-certs https -ep /usr/local/share/ca-certificates/aspnet/https.crt --format PEM sudo update-ca-certificates`

The preceding commands:

- Ensure the current user's developer certificate is created.
- Exports the certificate with elevated permissions needed for the `ca-certificates`folder, using the current user's environment.
- Removing the `-E`flag exports the root user certificate, generating it if necessary. Each newly generated certificate has a different thumbprint. When running as root,`sudo`and`-E`are not needed.

The path in the preceding command is specific for Ubuntu. For other distributions, select an appropriate path or use the path for the Certificate Authorities (CAs).

### Trust HTTPS certificate on Linux using Edge or Chrome

For chromium browsers on Linux:

- Install the - `libnss3-tools`for your distribution.
- Create or verify the - `$HOME/.pki/nssdb`folder exists on the machine.
- Export the certificate with the following command: - `dotnet dev-certs https sudo -E dotnet dev-certs https -ep /usr/local/share/ca-certificates/aspnet/https.crt --format PEM`- The path in the preceding command is specific for Ubuntu. For other distributions, select an appropriate path or use the path for the Certificate Authorities (CAs).
- Run the following commands: - `certutil -d sql:$HOME/.pki/nssdb -A -t "P,," -n localhost -i /usr/local/share/ca-certificates/aspnet/https.crt`
- Exit and restart the browser.

#### Trust the certificate with Firefox on Linux

- Export the certificate with the following command: - `dotnet dev-certs https sudo -E dotnet dev-certs https -ep /usr/local/share/ca-certificates/aspnet/https.crt --format PEM`- The path in the preceding command is specific for Ubuntu. For other distributions, select an appropriate path or use the path for the Certificate Authorities (CAs).
- Create a JSON file at - `/usr/lib/firefox/distribution/policies.json`with the following command:

```
cat <<EOF | sudo tee /usr/lib/firefox/distribution/policies.json
{
 "policies": {
 "Certificates": {
 "Install": [
 "/usr/local/share/ca-certificates/aspnet/https.crt"
 ]
 }
 }
}
EOF
```
Note: Ubuntu 21.10 Firefox comes as a snap package and the installation folder is `/snap/firefox/current/usr/lib/firefox`.

See Configure trust of HTTPS certificate using Firefox browser in this article for an alternative way to configure the policy file using the browser.

### Trust the certificate with Fedora 34

See:

- This GitHub comment
- Fedora: Using Shared System Certificates
- Set up a .NET development environment on Fedora.

### Trust the certificate with other distros

See this GitHub issue.

## Trust HTTPS certificate from Windows Subsystem for Linux

The following instructions don't work for some Linux distributions, such as Ubuntu 20.04. For more information, see GitHub issue dotnet/AspNetCore.Docs #23686.

The Windows Subsystem for Linux (WSL) generates an HTTPS self-signed development certificate, which by default isn't trusted in Windows. The easiest way to have Windows trust the WSL certificate, is to configure WSL to use the same certificate as Windows:

- On - **Windows**- `dotnet dev-certs https -ep https.pfx -p $CREDENTIAL_PLACEHOLDER$ --trust`- Where - `$CREDENTIAL_PLACEHOLDER$`is a password.
- In a WSL window, import the exported certificate on the WSL instance: - `dotnet dev-certs https --clean --import <<path-to-pfx>> --password $CREDENTIAL_PLACEHOLDER$`

The preceding approach is a one time operation per certificate and per WSL distribution. It's easier than exporting the certificate over and over. If you update or regenerate the certificate on windows, you might need to run the preceding commands again.

## Troubleshoot certificate problems such as certificate not trusted

This section provides help when the ASP.NET Core HTTPS development certificate has been installed and trusted, but you still have browser warnings that the certificate is not trusted. The ASP.NET Core HTTPS development certificate is used by Kestrel.

To repair the IIS Express certificate, see this Stackoverflow issue.

### All platforms - certificate not trusted

Run the following commands:

```
dotnet dev-certs https --clean
dotnet dev-certs https --trust
```
Close any browser instances open. Open a new browser window to app. Certificate trust is cached by browsers.

### dotnet dev-certs https --clean Fails

The preceding commands solve most browser trust issues. If the browser is still not trusting the certificate, follow the platform-specific suggestions that follow.

### Docker - certificate not trusted

- Delete the *C:\Users{USER}\AppData\Roaming\ASP.NET\Https*folder.
- Clean the solution. Delete the *bin*and*obj*folders.
- Restart the development tool. For example, Visual Studio or Visual Studio Code.

### Windows - certificate not trusted

- Check the certificates in the certificate store. There should be a `localhost`certificate with the`ASP.NET Core HTTPS development certificate`friendly name both under`Current User > Personal > Certificates`and`Current User > Trusted root certification authorities > Certificates`
- Remove all the found certificates from both Personal and Trusted root certification authorities. Do **not**remove the IIS Express localhost certificate.
- Run the following commands:

```
dotnet dev-certs https --clean
dotnet dev-certs https --trust
```
Close any browser instances open. Open a new browser window to app.

### OS X - certificate not trusted

- Open KeyChain Access.
- Select the System keychain.
- Check for the presence of a localhost certificate.
- Check that it contains a `+`symbol on the icon to indicate it's trusted for all users.
- Remove the certificate from the system keychain.
- Run the following commands:

```
dotnet dev-certs https --clean
dotnet dev-certs https --trust
```
Close any browser instances open. Open a new browser window to app.

See HTTPS Error using IIS Express (dotnet/AspNetCore #16892) for troubleshooting certificate issues with Visual Studio.

### Linux certificate not trusted

Check that the certificate being configured for trust is the user HTTPS developer certificate that will be used by the Kestrel server.

Check the current user default HTTPS developer Kestrel certificate at the following location:

```
ls -la ~/.dotnet/corefx/cryptography/x509stores/my
```
The HTTPS developer Kestrel certificate file is the SHA1 thumbprint. When the file is deleted via `dotnet dev-certs https --clean`, it's regenerated when needed with a different thumbprint.
Check the thumbprint of the exported certificate matches with the following command:

```
openssl x509 -noout -fingerprint -sha1 -inform pem -in /usr/local/share/ca-certificates/aspnet/https.crt
```
If the certificate doesn't match, it could be one of the following:

- An old certificate.
- An exported a developer certificate for the root user. For this case, export the certificate.

The root user certificate can be checked at:

```
ls -la /root/.dotnet/corefx/cryptography/x509stores/my
```
### IIS Express SSL certificate used with Visual Studio

To fix problems with the IIS Express certificate, select **Repair** from the Visual Studio installer. For more information, see this GitHub issue.

### Group policy prevents self-signed certificates from being trusted

In some cases, group policy may prevent self-signed certificates from being trusted. For more information, see this GitHub issue.

## Additional information

Note

If you're using .NET 9 or later SDK, see the updated Linux procedures in the .NET 9 version of this article.

Warning

## API projects

Do **not** use RequireHttpsAttribute on Web APIs that receive sensitive information. `RequireHttpsAttribute` uses HTTP status codes to redirect browsers from HTTP to HTTPS. API clients may not understand or obey redirects from HTTP to HTTPS. Such clients may send information over HTTP. Web APIs should either:

- Not listen on HTTP.
- Close the connection with status code 400 (Bad Request) and not serve the request.

To disable HTTP redirection in an API, set the `ASPNETCORE_URLS` environment variable or use the `--urls` command line flag. For more information, see ASP.NET Core runtime environments and 8 ways to set the URLs for an ASP.NET Core app by Andrew Lock.

## HSTS and API projects

The default API projects don't include HSTS because HSTS is generally a browser only instruction. Other callers, such as phone or desktop apps, do **not** obey the instruction. Even within browsers, a single authenticated call to an API over HTTP has risks on insecure networks. The secure approach is to configure API projects to only listen to and respond over HTTPS.

### HTTP redirection to HTTPS causes ERR_INVALID_REDIRECT on the CORS preflight request

Requests to an endpoint using HTTP that are redirected to HTTPS by UseHttpsRedirection fail with `ERR_INVALID_REDIRECT` on the CORS preflight request.

API projects can reject HTTP requests rather than use `UseHttpsRedirection` to redirect requests to HTTPS.

## Require HTTPS

We recommend that production ASP.NET Core web apps use:

- HTTPS redirection middleware (UseHttpsRedirection) to redirect HTTP requests to HTTPS.
- HSTS middleware (UseHsts) to send HTTP Strict Transport Security (HSTS) protocol headers to clients.

Note

Apps deployed in a reverse proxy configuration allow the proxy to handle connection security (HTTPS). If the proxy also handles HTTPS redirection, there's no need to use HTTPS redirection middleware. If the proxy server also handles writing HSTS headers (for example, native HSTS support in IIS 10.0 (1709) or later), HSTS middleware isn't required by the app. For more information, see Opt-out of HTTPS/HSTS on project creation.

### HTTPS redirection middleware (`UseHttpsRedirection`)

The following code calls UseHttpsRedirection in the `Program.cs` file:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
The preceding highlighted code:

- Uses the default HttpsRedirectionOptions.RedirectStatusCode (Status307TemporaryRedirect).
- Uses the default HttpsRedirectionOptions.HttpsPort (null) unless overridden by the `ASPNETCORE_HTTPS_PORT`environment variable or IServerAddressesFeature.

We recommend using temporary redirects rather than permanent redirects. Link caching can cause unstable behavior in development environments. If you prefer to send a permanent redirect status code when the app is in a non-`Development` environment, see the Configure permanent redirects in production section. We recommend using HSTS to signal to clients that only secure resource requests should be sent to the app (only in production).

Note

Don't confuse the `HTTPS_PORT` configuration key and `ASPNETCORE_HTTPS_PORT` environment variable, which set the port for HTTPS redirection middleware, with the `HTTPS_PORTS` configuration key and `ASPNETCORE_HTTPS_PORTS` environment variable, which set the ports for Kestrel/HTTP.sys endpoint configuration.

### Port configuration

A port must be available for the middleware to redirect an insecure request to HTTPS. If no port is available:

- Redirection to HTTPS doesn't occur.
- The middleware logs the warning "Failed to determine the https port for redirect."

Specify the HTTPS port using any of the following approaches:

- Set the - `https_port`host setting:- In host configuration.
- By setting the - `ASPNETCORE_HTTPS_PORT`environment variable.
- By adding a top-level entry in - `appsettings.json`:- `{ "https_port": 443, "Logging": { "LogLevel": { "Default": "Information", "Microsoft.AspNetCore": "Warning" } }, "AllowedHosts": "*" }`

- Indicate a port with the secure scheme using the ASPNETCORE_URLS environment variable. The environment variable configures the server. The middleware indirectly discovers the HTTPS port via IServerAddressesFeature. This approach doesn't work in reverse proxy deployments.
- The ASP.NET Core web templates set an HTTPS URL in - `Properties/launchsettings.json`for both Kestrel and IIS Express.- `launchsettings.json`is only used on the local machine.
- Configure an HTTPS URL endpoint for a public-facing edge deployment of Kestrel server or HTTP.sys server. Only - **one HTTPS port**is used by the app. The middleware discovers the port via IServerAddressesFeature.

Note

When an app is run in a reverse proxy configuration, IServerAddressesFeature isn't available. Set the port using one of the other approaches described in this section.

### Edge deployments

When Kestrel or HTTP.sys is used as a public-facing edge server, Kestrel or HTTP.sys must be configured to listen on both:

- The secure port where the client is redirected (typically, 443 in production and 5001 in development).
- The insecure port (typically, 80 in production and 5000 in development).

The insecure port must be accessible by the client in order for the app to receive an insecure request and redirect the client to the secure port.

For more information, see Kestrel endpoint configuration or HTTP.sys web server implementation in ASP.NET Core.

### Deployment scenarios

Any firewall between the client and server must also have communication ports open for traffic.

If requests are forwarded in a reverse proxy configuration, use forwarded headers middleware before calling HTTPS redirection middleware. Forwarded headers middleware updates the `Request.Scheme`, using the `X-Forwarded-Proto` header. The middleware permits redirect URIs and other security policies to work correctly. When forwarded headers middleware isn't used, the backend app might not receive the correct scheme and end up in a redirect loop. A common end user error message is that too many redirects have occurred.

When deploying to Azure App Service, follow the guidance in Tutorial: Bind an existing custom SSL certificate to Azure Web Apps.

### Options

The following highlighted code calls AddHttpsRedirection to configure middleware options:

```
using static Microsoft.AspNetCore.Http.StatusCodes;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
builder.Services.AddHsts(options =>
{
 options.Preload = true;
 options.IncludeSubDomains = true;
 options.MaxAge = TimeSpan.FromDays(60);
 options.ExcludedHosts.Add("example.com");
 options.ExcludedHosts.Add("www.example.com");
});
builder.Services.AddHttpsRedirection(options =>
{
 options.RedirectStatusCode = Status307TemporaryRedirect;
 options.HttpsPort = 5001;
});
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
Calling `AddHttpsRedirection` is only necessary to change the values of `HttpsPort` or `RedirectStatusCode`.

The preceding highlighted code:

- Sets HttpsRedirectionOptions.RedirectStatusCode to Status307TemporaryRedirect, which is the default value. Use the fields of the StatusCodes class for assignments to `RedirectStatusCode`.
- Sets the HTTPS port to 5001.

#### Configure permanent redirects in production

The middleware defaults to sending a Status307TemporaryRedirect with all redirects. If you prefer to send a permanent redirect status code when the app is in a non-`Development` environment, wrap the middleware options configuration in a conditional check for a non-`Development` environment.

When configuring services in `Program.cs`:

```
using static Microsoft.AspNetCore.Http.StatusCodes;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
if (!builder.Environment.IsDevelopment())
{
 builder.Services.AddHttpsRedirection(options =>
 {
 options.RedirectStatusCode = Status308PermanentRedirect;
 options.HttpsPort = 443;
 });
}
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
## HTTPS redirection middleware alternative approach

An alternative to using HTTPS redirection middleware (`UseHttpsRedirection`) is to use URL rewriting middleware (`AddRedirectToHttps`). `AddRedirectToHttps` can also set the status code and port when the redirect is executed. For more information, see URL rewriting middleware.

When redirecting to HTTPS without the requirement for additional redirect rules, we recommend using HTTPS redirection middleware (`UseHttpsRedirection`) described in this article.

## HTTP Strict Transport Security (HSTS) protocol

Per OWASP, HTTP Strict Transport Security (HSTS) is an opt-in security enhancement that's specified by a web app through the use of a response header. When a browser that supports HSTS receives this header:

- The browser stores configuration for the domain that prevents sending any communication over HTTP. The browser forces all communication over HTTPS.
- The browser prevents the user from using untrusted or invalid certificates. The browser disables prompts that allow a user to temporarily trust such a certificate.

Because HSTS is enforced by the client, it has some limitations:

- The client must support HSTS.
- HSTS requires at least one successful HTTPS request to establish the HSTS policy.
- The application must check every HTTP request and redirect or reject the HTTP request.

ASP.NET Core implements HSTS with the UseHsts extension method. The following code calls `UseHsts` when the app isn't in development mode:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
`UseHsts` isn't recommended in development because the HSTS settings are highly cacheable by browsers. By default, `UseHsts` excludes the local loopback address.

For production environments that are implementing HTTPS for the first time, set the initial HstsOptions.MaxAge to a small value using one of the TimeSpan methods. Set the value from hours to no more than a single day in case you need to revert the HTTPS infrastructure to HTTP. After you're confident in the sustainability of the HTTPS configuration, increase the HSTS `max-age` value; a commonly used value is one year.

The following highlighted code:

```
using static Microsoft.AspNetCore.Http.StatusCodes;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
builder.Services.AddHsts(options =>
{
 options.Preload = true;
 options.IncludeSubDomains = true;
 options.MaxAge = TimeSpan.FromDays(60);
 options.ExcludedHosts.Add("example.com");
 options.ExcludedHosts.Add("www.example.com");
});
builder.Services.AddHttpsRedirection(options =>
{
 options.RedirectStatusCode = Status307TemporaryRedirect;
 options.HttpsPort = 5001;
});
var app = builder.Build();
if (!app.Environment.IsDevelopment())
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
- Sets the preload parameter of the `Strict-Transport-Security`header. Preload isn't part of the RFC HSTS specification, but is supported by web browsers to preload HSTS sites on fresh install. For more information, see https://hstspreload.org/.
- Enables includeSubDomain, which applies the HSTS policy to Host subdomains.
- Explicitly sets the `max-age`parameter of the`Strict-Transport-Security`header to 60 days. If not set, defaults to 30 days. For more information, see the max-age directive.
- Adds `example.com`to the list of hosts to exclude.

`UseHsts` excludes the following loopback hosts:

- `localhost`: The IPv4 loopback address.
- `127.0.0.1`: The IPv4 loopback address.
- `[::1]`: The IPv6 loopback address.

## Opt-out of HTTPS/HSTS on project creation

In some backend service scenarios where connection security is handled at the public-facing edge of the network, configuring connection security at each node isn't required. Web apps that are generated from the templates in Visual Studio or from the dotnet new command enable HTTPS redirection and HSTS. For deployments that don't require these scenarios, you can opt-out of HTTPS/HSTS when the app is created from the template.

To opt-out of HTTPS/HSTS:

Uncheck the **Configure for HTTPS** checkbox.

## Trust the ASP.NET Core HTTPS development certificate on Windows and macOS

For the Firefox browser, see the next section.

The .NET SDK includes an HTTPS development certificate. The certificate is installed as part of the first-run experience. For example, `dotnet --info` produces a variation of the following output:

```
ASP.NET Core
------------
Successfully installed the ASP.NET Core HTTPS Development Certificate.
To trust the certificate run 'dotnet dev-certs https --trust' (Windows and macOS only).
For establishing trust on other platforms refer to the platform specific documentation.
For more information on configuring HTTPS see https://go.microsoft.com/fwlink/?linkid=848054.
```
Installing the .NET SDK installs the ASP.NET Core HTTPS development certificate to the local user certificate store. The certificate has been installed, but it's not trusted. To trust the certificate, perform the one-time step to run the `dotnet dev-certs` tool:

```
dotnet dev-certs https --trust
```
The following command provides help on the `dotnet dev-certs` tool:

```
dotnet dev-certs https --help
```
Warning

Do not create a development certificate in an environment that will be redistributed, such as a container image or virtual machine. Doing so can lead to spoofing and elevation of privilege. To help prevent this, set the `DOTNET_GENERATE_ASPNET_CERTIFICATE` environment variable to `false` prior to calling the .NET CLI for the first time. This will skip the automatic generation of the ASP.NET Core development certificate during the CLI's first-run experience.

### Trust the HTTPS certificate with Firefox to prevent SEC_ERROR_INADEQUATE_KEY_USAGE error

The Firefox browser uses its own certificate store, and therefore doesn't trust the IIS Express or Kestrel developer certificates.

There are two approaches to trusting the HTTPS certificate with Firefox, create a policy file or configure with the FireFox browser. Configuring with the browser creates the policy file, so the two approaches are equivalent.

#### Create a policy file to trust HTTPS certificate with Firefox

Create a policy file (`policies.json`) at:

- Windows: `%PROGRAMFILES%\Mozilla Firefox\distribution\`
- MacOS: `Firefox.app/Contents/Resources/distribution`
- Linux: See Trust the certificate with Firefox on Linux in this article.

Add the following JSON to the Firefox policy file:

```
{
 "policies": {
 "Certificates": {
 "ImportEnterpriseRoots": true
 }
 }
}
```
The preceding policy file makes Firefox trust certificates from the trusted certificates in the Windows certificate store. The next section provides an alternative approach to create the preceding policy file by using the Firefox browser.

### Configure trust of HTTPS certificate using Firefox browser

Set `security.enterprise_roots.enabled` = `true` using the following instructions:

- Enter `about:config`in the FireFox browser.
- Select **Accept the Risk and Continue**if you accept the risk.
- Select **Show All**
- Set `security.enterprise_roots.enabled`=`true`
- Exit and restart Firefox

For more information, see Setting Up Certificate Authorities (CAs) in Firefox and the mozilla/policy-templates/README file.

## How to set up a developer certificate for Docker

See this GitHub issue.

## Trust HTTPS certificate on Linux

Establishing trust is distribution and browser specific. The following sections provide instructions for some popular distributions and the Chromium browsers (Edge and Chrome) and for Firefox.

### Trust HTTPS certificate on Linux with linux-dev-certs

linux-dev-certs is an open-source, community-supported, .NET global tool that provides a convenient way to create and trust a developer certificate on Linux. The tool is not maintained or supported by Microsoft.

The following commands install the tool and create a trusted developer certificate:

```
dotnet tool update -g linux-dev-certs
dotnet linux-dev-certs install
```
For more information or to report issues, see the linux-dev-certs GitHub repository.

### Ubuntu trust the certificate for service-to-service communication

The following instructions don't work for some Ubuntu versions, such as 20.04. For more information, see GitHub issue dotnet/AspNetCore.Docs #23686.

- Install OpenSSL 1.1.1h or later. See your distribution for instructions on how to update OpenSSL.
- Run the following commands: - `dotnet dev-certs https sudo -E dotnet dev-certs https -ep /usr/local/share/ca-certificates/aspnet/https.crt --format PEM sudo update-ca-certificates`

The preceding commands:

- Ensure the current user's developer certificate is created.
- Exports the certificate with elevated permissions needed for the `ca-certificates`folder, using the current user's environment.
- Removing the `-E`flag exports the root user certificate, generating it if necessary. Each newly generated certificate has a different thumbprint. When running as root,`sudo`and`-E`are not needed.

The path in the preceding command is specific for Ubuntu. For other distributions, select an appropriate path or use the path for the Certificate Authorities (CAs).

### Trust HTTPS certificate on Linux using Edge or Chrome

For chromium browsers on Linux:

- Install the - `libnss3-tools`for your distribution.
- Create or verify the - `$HOME/.pki/nssdb`folder exists on the machine.
- Export the certificate with the following command: - `dotnet dev-certs https sudo -E dotnet dev-certs https -ep /usr/local/share/ca-certificates/aspnet/https.crt --format PEM`- The path in the preceding command is specific for Ubuntu. For other distributions, select an appropriate path or use the path for the Certificate Authorities (CAs).
- Run the following commands: - `certutil -d sql:$HOME/.pki/nssdb -A -t "P,," -n localhost -i /usr/local/share/ca-certificates/aspnet/https.crt`
- Exit and restart the browser.

#### Trust the certificate with Firefox on Linux

- Export the certificate with the following command: - `dotnet dev-certs https sudo -E dotnet dev-certs https -ep /usr/local/share/ca-certificates/aspnet/https.crt --format PEM`- The path in the preceding command is specific for Ubuntu. For other distributions, select an appropriate path or use the path for the Certificate Authorities (CAs).
- Create a JSON file at - `/usr/lib/firefox/distribution/policies.json`with the following command:

```
cat <<EOF | sudo tee /usr/lib/firefox/distribution/policies.json
{
 "policies": {
 "Certificates": {
 "Install": [
 "/usr/local/share/ca-certificates/aspnet/https.crt"
 ]
 }
 }
}
EOF
```
Note: Ubuntu 21.10 Firefox comes as a snap package and the installation folder is `/snap/firefox/current/usr/lib/firefox`.

See Configure trust of HTTPS certificate using Firefox browser in this article for an alternative way to configure the policy file using the browser.

### Trust the certificate with Fedora 34

See:

- This GitHub comment
- Fedora: Using Shared System Certificates
- Set up a .NET development environment on Fedora.

### Trust the certificate with other distros

See this GitHub issue.

## Trust HTTPS certificate from Windows Subsystem for Linux

The following instructions don't work for some Linux distributions, such as Ubuntu 20.04. For more information, see GitHub issue dotnet/AspNetCore.Docs #23686.

The Windows Subsystem for Linux (WSL) generates an HTTPS self-signed development certificate, which by default isn't trusted in Windows. The easiest way to have Windows trust the WSL certificate, is to configure WSL to use the same certificate as Windows:

- On - **Windows**- `dotnet dev-certs https -ep https.pfx -p $CREDENTIAL_PLACEHOLDER$ --trust`- Where - `$CREDENTIAL_PLACEHOLDER$`is a password.
- In a WSL window, import the exported certificate on the WSL instance: - `dotnet dev-certs https --clean --import <<path-to-pfx>> --password $CREDENTIAL_PLACEHOLDER$`

The preceding approach is a one time operation per certificate and per WSL distribution. It's easier than exporting the certificate over and over. If you update or regenerate the certificate on windows, you might need to run the preceding commands again.

## Troubleshoot certificate problems such as certificate not trusted

This section provides help when the ASP.NET Core HTTPS development certificate has been installed and trusted, but you still have browser warnings that the certificate is not trusted. The ASP.NET Core HTTPS development certificate is used by Kestrel.

To repair the IIS Express certificate, see this Stackoverflow issue.

### All platforms - certificate not trusted

Run the following commands:

```
dotnet dev-certs https --clean
dotnet dev-certs https --trust
```
Close any browser instances open. Open a new browser window to app. Certificate trust is cached by browsers.

### dotnet dev-certs https --clean Fails

The preceding commands solve most browser trust issues. If the browser is still not trusting the certificate, follow the platform-specific suggestions that follow.

### Docker - certificate not trusted

- Delete the *C:\Users{USER}\AppData\Roaming\ASP.NET\Https*folder.
- Clean the solution. Delete the *bin*and*obj*folders.
- Restart the development tool. For example, Visual Studio or Visual Studio Code.

### Windows - certificate not trusted

- Check the certificates in the certificate store. There should be a `localhost`certificate with the`ASP.NET Core HTTPS development certificate`friendly name both under`Current User > Personal > Certificates`and`Current User > Trusted root certification authorities > Certificates`
- Remove all the found certificates from both Personal and Trusted root certification authorities. Do **not**remove the IIS Express localhost certificate.
- Run the following commands:

```
dotnet dev-certs https --clean
dotnet dev-certs https --trust
```
Close any browser instances open. Open a new browser window to app.

### OS X - certificate not trusted

- Open KeyChain Access.
- Select the System keychain.
- Check for the presence of a localhost certificate.
- Check that it contains a `+`symbol on the icon to indicate it's trusted for all users.
- Remove the certificate from the system keychain.
- Run the following commands:

```
dotnet dev-certs https --clean
dotnet dev-certs https --trust
```
Close any browser instances open. Open a new browser window to app.

See HTTPS Error using IIS Express (dotnet/AspNetCore #16892) for troubleshooting certificate issues with Visual Studio.

### Linux certificate not trusted

Check that the certificate being configured for trust is the user HTTPS developer certificate that will be used by the Kestrel server.

Check the current user default HTTPS developer Kestrel certificate at the following location:

```
ls -la ~/.dotnet/corefx/cryptography/x509stores/my
```
The HTTPS developer Kestrel certificate file is the SHA1 thumbprint. When the file is deleted via `dotnet dev-certs https --clean`, it's regenerated when needed with a different thumbprint.
Check the thumbprint of the exported certificate matches with the following command:

```
openssl x509 -noout -fingerprint -sha1 -inform pem -in /usr/local/share/ca-certificates/aspnet/https.crt
```
If the certificate doesn't match, it could be one of the following:

- An old certificate.
- An exported a developer certificate for the root user. For this case, export the certificate.

The root user certificate can be checked at:

```
ls -la /root/.dotnet/corefx/cryptography/x509stores/my
```
### IIS Express SSL certificate used with Visual Studio

To fix problems with the IIS Express certificate, select **Repair** from the Visual Studio installer. For more information, see this GitHub issue.

### Group policy prevents self-signed certificates from being trusted

In some cases, group policy may prevent self-signed certificates from being trusted. For more information, see this GitHub issue.

## Additional information

ASP.NET Core

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Note

This isn't the latest version of this article. For the current release, see the .NET 10 version of this article.

Warning

This version of ASP.NET Core is no longer supported. For more information, see the .NET and .NET Core Support Policy. For the current release, see the .NET 10 version of this article.

By Rick Anderson and Kirk Larkin

This article explains how to manage sensitive data for an ASP.NET Core app on a development machine. Never store passwords or other sensitive data in source code or configuration files. Production secrets shouldn't be used for development or test. Secrets shouldn't be deployed with the app. Production secrets should be accessed through a controlled means like Azure Key Vault. Azure test and production secrets can be stored and protected with the Azure Key Vault configuration provider.

View or download the sample code (how to download)

For more information on authentication for deployed test and production apps, see Secure authentication flows.

To use user secrets in a .NET console app, see GitHub dotnet/entityframework.docs issue #3939.

## Work with environment variables

Environment variables are used to avoid storage of app secrets in code or in local configuration files. Environment variables override configuration values for all previously specified configuration sources.

Consider an ASP.NET Core web app in which **Individual Accounts** security is enabled. A default database connection string is included in the project *appsettings.json* file with the `DefaultConnection` key. The default connection string is for LocalDB, which runs in user mode and doesn't require a password. During app deployment, you can override the `DefaultConnection` key value with the value from an environment variable. The environment variable might store the complete connection string with sensitive credentials.

Warning

Environment variables are commonly stored as plain, unencrypted text. If the machine or process is compromised, environment variables are accessible to untrusted parties. Extra measures to prevent disclosure of user secrets might be required.

The colon (`:`) separator doesn't work with environment variable hierarchical keys on all platforms. For example, Bash doesn't support colon (`:`) as a separator. All platforms support the double underscore (`__`) syntax and automatically replace it with a colon (`:`).

## Use the Secret Manager tool

Secret Manager is a tool that stores sensitive data during application development. In this context, a piece of sensitive data is an *app secret*.

- App secrets are stored in a separate location from the project tree.
- They're associated with a specific project or shared across several projects.
- They aren't checked into source control.

Warning

Secret Manager doesn't encrypt the stored secrets and shouldn't be treated as a trusted store. It's for development purposes only. The keys and values are stored in a JSON configuration file in the user profile directory.

Secret Manager hides implementation details, such as where and how the values are stored. You can use the tool without knowing these implementation details. The values are stored in a JSON file in the local machine's user profile folder:

File system path:

`%APPDATA%\Microsoft\UserSecrets\<user_secrets_id>\secrets.json`

In the file system path, replace the `<user_secrets_id>` portion with the `UserSecretsId` value specified in your project file.

Don't write code that depends on the location or format of data saved with Secret Manager. These implementation details might change. For example, the secret values aren't encrypted.

## Enable secret storage

Secret Manager operates on project-specific configuration settings stored in your user profile.

### Use the CLI

Secret Manager includes an `init` command. To use user secrets, run the following command in the project directory:

```
dotnet user-secrets init
```
This command adds a `UserSecretsId` element within a `PropertyGroup` of the project file. By default, the inner text of `UserSecretsId` is a GUID. The inner text is arbitrary, but is unique to the project. The following example shows a GUID value of `0000a1a1-b2b2-c3c3-d4d4-eeeeee555555`.

```
<Project Sdk="Microsoft.NET.Sdk.Web">
 <PropertyGroup>
 <TargetFramework>net9.0</TargetFramework>
 <UserSecretsId>0000a1a1-b2b2-c3c3-d4d4-eeeeee555555</UserSecretsId>
 </PropertyGroup>
</Project>
```
### Use Visual Studio

In Visual Studio, right-click the project in Solution Explorer, and select **Manage User Secrets** from the context menu. This gesture adds a `UserSecretsId` element, populated with a GUID, to the project file.

### If 'GenerateAssemblyInfo' is 'false'

If the generation of assembly info attributes (`GenerateAssemblyInfo`) is disabled (set to `false`), manually add the UserSecretsIdAttribute in the *AssemblyInfo.cs* file. For example:

```
[assembly: UserSecretsId("your_user_secrets_id")]
```
When you manually add the `UserSecretsId` attribute to the *AssemblyInfo.cs* file, the `UserSecretsId` value must match the value in the project file.

## Set a secret

Define an app secret consisting of a key and its value. The secret is associated with the project's `UserSecretsId` value. For example, run the following command from the directory in which the project file exists:

```
dotnet user-secrets set "Movies:ServiceApiKey" "12345"
```
In this example, the colon indicates that `Movies` is an object literal with a `ServiceApiKey` property.

You can also use Secret Manager from other directories. Include the `--project` option to supply the file system path at which the project file exists. For example:

```
dotnet user-secrets set "Movies:ServiceApiKey" "12345" --project "C:\apps\WebApp1\src\WebApp1"
```
### JSON structure flattening in Visual Studio

The Visual Studio **Manage User Secrets** gesture opens a *secrets.json* file in the text editor. Replace the contents of the *secrets.json* file with the key-value pairs to store. For example:

```
{
 "Movies": {
 "ConnectionString": "Server=(localdb)\\mssqllocaldb;Database=Movie-1;Trusted_Connection=True;MultipleActiveResultSets=true",
 "ServiceApiKey": "12345"
 }
}
```
The JSON structure is flattened after modifications via the `dotnet user-secrets remove` or `dotnet user-secrets set` command. For example, running `dotnet user-secrets remove "Movies:ConnectionString"` collapses the `Movies` object literal. The modified file resembles the following JSON:

```
{
 "Movies:ServiceApiKey": "12345"
}
```
## Set multiple secrets

A batch of secrets can be set by piping JSON to the `set` command. In the following example, the contents of the *input.json* file is piped to the `set` command.

Run the following command:

```
type .\input.json | dotnet user-secrets set
```
## Access a secret

To access a secret, complete the following steps:

### Register the user secrets configuration source

The user secrets configuration provider registers the appropriate configuration source with the .NET Configuration API.

ASP.NET Core web apps created with the dotnet new command or Visual Studio generate the following code:

```
var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();
app.MapGet("/", () => "Hello World!");
app.Run();
```
The WebApplication.CreateBuilder method initializes a new instance of the WebApplicationBuilder class with preconfigured defaults. The initialized `WebApplicationBuilder` (`builder`) provides default configuration and calls the AddUserSecrets method when the EnvironmentName property is Development.

### Read the secret via the Configuration API

The following examples demonstrate how to read the `Movies:ServiceApiKey` key:

**Program.cs file**

```
var builder = WebApplication.CreateBuilder(args);
var movieApiKey = builder.Configuration["Movies:ServiceApiKey"];
var app = builder.Build();
app.MapGet("/", () => movieApiKey);
app.Run();
```
**Razor Pages page model**

```
public class IndexModel : PageModel
{
 private readonly IConfiguration _config;
 public IndexModel(IConfiguration config)
 {
 _config = config;
 }
 public void OnGet()
 {
 var moviesApiKey = _config["Movies:ServiceApiKey"];
 // call Movies service with the API key
 }
}
```
For more information, see Configuration in ASP.NET Core.

## Map secrets to a POCO

Mapping an entire object literal to a POCO (a simple .NET class with properties) is useful for aggregating related properties.

Assume the application *secrets.json* file contains the following two secrets:

```
{
 "Movies:ConnectionString": "Server=(localdb)\\mssqllocaldb;Database=Movie-1;Trusted_Connection=True;MultipleActiveResultSets=true",
 "Movies:ServiceApiKey": "12345"
}
```
To map the preceding secrets to a POCO, use the .NET Configuration API's object graph binding feature. The following code binds to a custom `MovieSettings` POCO and accesses the `ServiceApiKey` property value:

```
var moviesConfig =
 Configuration.GetSection("Movies").Get<MovieSettings>();
_moviesApiKey = moviesConfig.ServiceApiKey;
```
The `Movies:ConnectionString` and `Movies:ServiceApiKey` secrets are mapped to the respective properties in `MovieSettings`:

```
public class MovieSettings
{
 public string ConnectionString { get; set; }
 public string ServiceApiKey { get; set; }
}
```
## Use string replacement with secrets

Storing passwords in plain text is insecure. Never store secrets in a configuration file such as *appsettings.json*, which might get checked in to a source code repository.

For example, a database connection string stored in an *appsettings.json* file shouldn't include a password. Instead, store the password as a secret, and include the password in the connection string at runtime. For example:

```
dotnet user-secrets set "DbPassword" "`<secret value>`"
```
Replace the `<secret value>` placeholder in the example with the password value. Set the secret's value on a SqlConnectionStringBuilder object's Password property to include it as the password value in the connection string:

```
using System.Data.SqlClient;
var builder = WebApplication.CreateBuilder(args);
var conStrBuilder = new SqlConnectionStringBuilder(
 builder.Configuration.GetConnectionString("Movies"));
conStrBuilder.Password = builder.Configuration["DbPassword"];
var connection = conStrBuilder.ConnectionString;
var app = builder.Build();
app.MapGet("/", () => connection);
app.Run();
```
## List the secrets

Assume the application *secrets.json* file contains the following two secrets:

```
{
 "Movies:ConnectionString": "Server=(localdb)\\mssqllocaldb;Database=Movie-1;Trusted_Connection=True;MultipleActiveResultSets=true",
 "Movies:ServiceApiKey": "12345"
}
```
Run the following command from the directory in which the project file exists:

```
dotnet user-secrets list
```
The following output appears:

```
Movies:ConnectionString = Server=(localdb)\mssqllocaldb;Database=Movie-1;Trusted_Connection=True;MultipleActiveResultSets=true
Movies:ServiceApiKey = 12345
```
In the example, a colon (`:`) in the key names denotes the object hierarchy within the *secrets.json* file.

## Remove a single secret

Assume the application *secrets.json* file contains the following two secrets:

```
{
 "Movies:ConnectionString": "Server=(localdb)\\mssqllocaldb;Database=Movie-1;Trusted_Connection=True;MultipleActiveResultSets=true",
 "Movies:ServiceApiKey": "12345"
}
```
Run the following command from the directory in which the project file exists:

```
dotnet user-secrets remove "Movies:ConnectionString"
```
The application *secrets.json* file is modified to remove the key-value pair associated with the `Movies:ConnectionString` key:

```
{
 "Movies": {
 "ServiceApiKey": "12345"
 }
}
```
The `dotnet user-secrets list` command displays the following message:

```
Movies:ServiceApiKey = 12345
```
## Remove all secrets

Assume the application *secrets.json* file contains the following two secrets:

```
{
 "Movies:ConnectionString": "Server=(localdb)\\mssqllocaldb;Database=Movie-1;Trusted_Connection=True;MultipleActiveResultSets=true",
 "Movies:ServiceApiKey": "12345"
}
```
Run the following command from the directory in which the project file exists:

```
dotnet user-secrets clear
```
All user secrets for the app are deleted from the *secrets.json* file:

```
{}
```
Running the `dotnet user-secrets list` command displays the following message:

```
No secrets configured for this application.
```
## Manage user secrets with Visual Studio

To manage user secrets in Visual Studio, right-click the project in Solution Explorer and select **Manage User Secrets**:

## Migrate user secrets from ASP.NET Framework to ASP.NET Core

You can migrate your stored user secrets from ASP.NET Framework to ASP.NET Core. For more information, see GitHub dotnet/aspnetcore.docs issue #27611 - *User Secrets documentation doesn't mention incompatibility with AssemblyInfo.cs*.

## Work with user secrets in non-web applications

Projects that target `Microsoft.NET.Sdk.Web` automatically include support for user secrets. For projects that target `Microsoft.NET.Sdk`, such as console applications, install the configuration extension and user secrets NuGet packages explicitly.

```
Install-Package Microsoft.Extensions.Configuration
Install-Package Microsoft.Extensions.Configuration.UserSecrets
```
After you install the packages, initialize the project and set secrets the same way as for a web app. The following example shows a console application that retrieves the value of a secret set with the `AppSecret` key:

```
using Microsoft.Extensions.Configuration;
namespace ConsoleApp;
class Program
{
 static void Main(string[] args)
 {
 IConfigurationRoot config = new ConfigurationBuilder()
 .AddUserSecrets<Program>()
 .Build();
 Console.WriteLine(config["AppSecret"]);
 }
}
```
## Related content

By Rick Anderson, Kirk Larkin, Daniel Roth, and Scott Addie

View or download sample code (how to download)

This article explains how to manage sensitive data for an ASP.NET Core app on a development machine. Never store passwords or other sensitive data in source code or configuration files. Production secrets shouldn't be used for development or test. Secrets shouldn't be deployed with the app. Production secrets should be accessed through a controlled means like Azure Key Vault. Azure test and production secrets can be stored and protected with the Azure Key Vault configuration provider.

For more information on authentication for test and production environments, see Secure authentication flows.

## Environment variables

Environment variables are used to avoid storage of app secrets in code or in local configuration files. Environment variables override configuration values for all previously specified configuration sources.

Consider an ASP.NET Core web app in which **Individual User Accounts** security is enabled. A default database connection string is included in the project's `appsettings.json` file with the key `DefaultConnection`. The default connection string is for LocalDB, which runs in user mode and doesn't require a password. During app deployment, the `DefaultConnection` key value can be overridden with an environment variable's value. The environment variable may store the complete connection string with sensitive credentials.

Warning

Environment variables are generally stored in plain, unencrypted text. If the machine or process is compromised, environment variables can be accessed by untrusted parties. Additional measures to prevent disclosure of user secrets may be required.

The colon (`:`) separator doesn't work with environment variable hierarchical keys on all platforms. For example, Bash doesn't support colon (`:`) as a separator. All platforms support the double underscore (`__`) syntax and automatically replace it with a colon (`:`).

## Secret Manager

The Secret Manager tool stores sensitive data during application development. In this context, a piece of sensitive data is an app secret. App secrets are stored in a separate location from the project tree. The app secrets are associated with a specific project or shared across several projects. The app secrets aren't checked into source control.

Warning

The Secret Manager tool doesn't encrypt the stored secrets and shouldn't be treated as a trusted store. It's for development purposes only. The keys and values are stored in a JSON configuration file in the user profile directory.

## How the Secret Manager tool works

The Secret Manager tool hides implementation details, such as where and how the values are stored. You can use the tool without knowing these implementation details. The values are stored in a JSON file in the local machine's user profile folder:

File system path:

`%APPDATA%\Microsoft\UserSecrets\<user_secrets_id>\secrets.json`

In the preceding file paths, replace `<user_secrets_id>` with the `UserSecretsId` value specified in the project file.

Don't write code that depends on the location or format of data saved with the Secret Manager tool. These implementation details may change. For example, the secret values aren't encrypted, but could be in the future.

## Enable secret storage

The Secret Manager tool operates on project-specific configuration settings stored in your user profile.

The Secret Manager tool includes an `init` command in .NET Core SDK 3.0.100 or later. To use user secrets, run the following command in the project directory:

```
dotnet user-secrets init
```
The preceding command adds a `UserSecretsId` element within a `PropertyGroup` of the project file. By default, the inner text of `UserSecretsId` is a GUID. The inner text is arbitrary, but is unique to the project.

```
<PropertyGroup>
 <TargetFramework>netcoreapp3.1</TargetFramework>
 <UserSecretsId>79a3edd0-2092-40a2-a04d-dcb46d5ca9ed</UserSecretsId>
</PropertyGroup>
```
In Visual Studio, right-click the project in Solution Explorer, and select **Manage User Secrets** from the context menu. This gesture adds a `UserSecretsId` element, populated with a GUID, to the project file.

## Set a secret

Define an app secret consisting of a key and its value. The secret is associated with the project's `UserSecretsId` value. For example, run the following command from the directory in which the project file exists:

```
dotnet user-secrets set "Movies:ServiceApiKey" "12345"
```
In the preceding example, the colon denotes that `Movies` is an object literal with a `ServiceApiKey` property.

The Secret Manager tool can be used from other directories too. Use the `--project` option to supply the file system path at which the project file exists. For example:

```
dotnet user-secrets set "Movies:ServiceApiKey" "12345" --project "C:\apps\WebApp1\src\WebApp1"
```
### JSON structure flattening in Visual Studio

Visual Studio's **Manage User Secrets** gesture opens a `secrets.json` file in the text editor. Replace the contents of `secrets.json` with the key-value pairs to be stored. For example:

```
{
 "Movies": {
 "ConnectionString": "Server=(localdb)\\mssqllocaldb;Database=Movie-1;Trusted_Connection=True;MultipleActiveResultSets=true",
 "ServiceApiKey": "12345"
 }
}
```
The JSON structure is flattened after modifications via `dotnet user-secrets remove` or `dotnet user-secrets set`. For example, running `dotnet user-secrets remove "Movies:ConnectionString"` collapses the `Movies` object literal. The modified file resembles the following JSON:

```
{
 "Movies:ServiceApiKey": "12345"
}
```
## Set multiple secrets

A batch of secrets can be set by piping JSON to the `set` command. In the following example, the `input.json` file's contents are piped to the `set` command.

Open a command shell, and execute the following command:

```
type .\input.json | dotnet user-secrets set
```
## Access a secret

To access a secret, complete the following steps:

### Register the user secrets configuration source

The user secrets configuration provider registers the appropriate configuration source with the .NET Configuration API.

The user secrets configuration source is automatically added in Development mode when the project calls CreateDefaultBuilder. `CreateDefaultBuilder` calls AddUserSecrets when the EnvironmentName is Development:

```
public static IHostBuilder CreateHostBuilder(string[] args) =>
 Host.CreateDefaultBuilder(args)
 .ConfigureWebHostDefaults(webBuilder =>
 {
 webBuilder.UseStartup<Startup>();
 });
```
When `CreateDefaultBuilder` isn't called, add the user secrets configuration source explicitly by calling AddUserSecrets in ConfigureAppConfiguration. Call `AddUserSecrets` only when the app runs in the `Development` environment, as shown in the following example:

```
public class Program
{
 public static void Main(string[] args)
 {
 var host = new HostBuilder()
 .ConfigureAppConfiguration((hostContext, builder) =>
 {
 // Add other providers for JSON, etc.
 if (hostContext.HostingEnvironment.IsDevelopment())
 {
 builder.AddUserSecrets<Program>();
 }
 })
 .Build();

 host.Run();
 }
}
```
### Read the secret via the Configuration API

If the user secrets configuration source is registered, the .NET Configuration API can read the secrets. Constructor injection can be used to gain access to the .NET Configuration API. Consider the following examples of reading the `Movies:ServiceApiKey` key:

**Startup class:**

```
public class Startup
{
 private string _moviesApiKey = null;
 public Startup(IConfiguration configuration)
 {
 Configuration = configuration;
 }
 public IConfiguration Configuration { get; }
 public void ConfigureServices(IServiceCollection services)
 {
 _moviesApiKey = Configuration["Movies:ServiceApiKey"];
 }
 public void Configure(IApplicationBuilder app)
 {
 app.Run(async (context) =>
 {
 var result = string.IsNullOrEmpty(_moviesApiKey) ? "Null" : "Not Null";
 await context.Response.WriteAsync($"Secret is {result}");
 });
 }
}
```
**Razor Pages page model:**

```
public class IndexModel : PageModel
{
 private readonly IConfiguration _config;
 public IndexModel(IConfiguration config)
 {
 _config = config;
 }
 public void OnGet()
 {
 var moviesApiKey = _config["Movies:ServiceApiKey"];
 // call Movies service with the API key
 }
}
```
For more information, see Access configuration in Startup and Access configuration in Razor Pages.

## Map secrets to a POCO

Mapping an entire object literal to a POCO (a simple .NET class with properties) is useful for aggregating related properties.

Assume the application *secrets.json* file contains the following two secrets:

```
{
 "Movies:ConnectionString": "Server=(localdb)\\mssqllocaldb;Database=Movie-1;Trusted_Connection=True;MultipleActiveResultSets=true",
 "Movies:ServiceApiKey": "12345"
}
```
To map the preceding secrets to a POCO, use the .NET Configuration API's object graph binding feature. The following code binds to a custom `MovieSettings` POCO and accesses the `ServiceApiKey` property value:

```
var moviesConfig =
 Configuration.GetSection("Movies").Get<MovieSettings>();
_moviesApiKey = moviesConfig.ServiceApiKey;
```
The `Movies:ConnectionString` and `Movies:ServiceApiKey` secrets are mapped to the respective properties in `MovieSettings`:

```
public class MovieSettings
{
 public string ConnectionString { get; set; }
 public string ServiceApiKey { get; set; }
}
```
## String replacement with secrets

Storing passwords in plain text is insecure. Never store secrets in a configuration file such as `appsettings.json`, which might get checked in to a source code repository.

For example, a database connection string stored in `appsettings.json` should not include a password. Instead, store the password as a secret, and include the password in the connection string at runtime. For example:

```
dotnet user-secrets set "DbPassword" "<secret value>"
```
Replace the `<secret value>` placeholder in the preceding example with the password value. Set the secret's value on a SqlConnectionStringBuilder object's Password property to include it as the password value in the connection string:

```
using System.Data.SqlClient;
var builder = WebApplication.CreateBuilder(args);
var conStrBuilder = new SqlConnectionStringBuilder(
 builder.Configuration.GetConnectionString("Movies"));
conStrBuilder.Password = builder.Configuration["DbPassword"];
var connection = conStrBuilder.ConnectionString;
var app = builder.Build();
app.MapGet("/", () => connection);
app.Run();
```
## List the secrets

Assume the application *secrets.json* file contains the following two secrets:

```
{
 "Movies:ConnectionString": "Server=(localdb)\\mssqllocaldb;Database=Movie-1;Trusted_Connection=True;MultipleActiveResultSets=true",
 "Movies:ServiceApiKey": "12345"
}
```
Run the following command from the directory in which the project file exists:

```
dotnet user-secrets list
```
The following output appears:

```
Movies:ConnectionString = Server=(localdb)\mssqllocaldb;Database=Movie-1;Trusted_Connection=True;MultipleActiveResultSets=true
Movies:ServiceApiKey = 12345
```
In the preceding example, a colon in the key names denotes the object hierarchy within `secrets.json`.

## Remove a single secret

Assume the application *secrets.json* file contains the following two secrets:

```
{
 "Movies:ConnectionString": "Server=(localdb)\\mssqllocaldb;Database=Movie-1;Trusted_Connection=True;MultipleActiveResultSets=true",
 "Movies:ServiceApiKey": "12345"
}
```
Run the following command from the directory in which the project file exists:

```
dotnet user-secrets remove "Movies:ConnectionString"
```
The app's `secrets.json` file was modified to remove the key-value pair associated with the `MoviesConnectionString` key:

```
{
 "Movies": {
 "ServiceApiKey": "12345"
 }
}
```
`dotnet user-secrets list` displays the following message:

```
Movies:ServiceApiKey = 12345
```
## Remove all secrets

Assume the application *secrets.json* file contains the following two secrets:

```
{
 "Movies:ConnectionString": "Server=(localdb)\\mssqllocaldb;Database=Movie-1;Trusted_Connection=True;MultipleActiveResultSets=true",
 "Movies:ServiceApiKey": "12345"
}
```
Run the following command from the directory in which the project file exists:

```
dotnet user-secrets clear
```
All user secrets for the app have been deleted from the `secrets.json` file:

```
{}
```
Running `dotnet user-secrets list` displays the following message:

```
No secrets configured for this application.
```
## Manage user secrets with Visual Studio

To manage user secrets in Visual Studio, right click the project in solution explorer and select **Manage User Secrets**:

## Migrating User Secrets from ASP.NET Framework to ASP.NET Core

See this GitHub issue.

## User secrets in non-web applications

Projects that target `Microsoft.NET.Sdk.Web` automatically include support for user secrets. For projects that target `Microsoft.NET.Sdk`, such as console applications, install the configuration extension and user secrets NuGet packages explicitly.

Using PowerShell:

```
Install-Package Microsoft.Extensions.Configuration
Install-Package Microsoft.Extensions.Configuration.UserSecrets
```
Using the .NET CLI:

```
dotnet add package Microsoft.Extensions.Configuration
dotnet add package Microsoft.Extensions.Configuration.UserSecrets
```
Once the packages are installed, initialize the project and set secrets the same way as for a web app. The following example shows a console application that retrieves the value of a secret that was set with the key "AppSecret":

```
using Microsoft.Extensions.Configuration;
namespace ConsoleApp;
class Program
{
 static void Main(string[] args)
 {
 IConfigurationRoot config = new ConfigurationBuilder()
 .AddUserSecrets<Program>()
 .Build();
 Console.WriteLine(config["AppSecret"]);
 }
}
```
## Additional resources

- See this issue and this issue for information on accessing user secrets from IIS.
- Configuration in ASP.NET Core
- Azure Key Vault configuration provider

ASP.NET Core

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

By Fiyaz Hasan and Rick Anderson

Cross-site request forgery is an attack against web-hosted apps whereby a malicious web app can influence the interaction between a client browser and a web app that trusts that browser. These attacks are possible because web browsers send some types of authentication tokens automatically with every request to a website. This form of exploit is also known as a *one-click attack* or *session riding* because the attack takes advantage of the user's previously authenticated session. Cross-site request forgery is also known as XSRF or CSRF.

An example of a CSRF attack:

- A user signs into - `www.good-banking-site.example.com`using forms authentication. The server authenticates the user and issues a response that includes an authentication cookie. The site is vulnerable to attack because it trusts any request that it receives with a valid authentication cookie.
- The user visits a malicious site, - `www.bad-crook-site.example.com`.- The malicious site, - `www.bad-crook-site.example.com`, contains an HTML form similar to the following example:- `<h1>Congratulations! You're a Winner!</h1> <form action="https://www.good-banking-site.example.com/api/account" method="post"> <input type="hidden" name="Transaction" value="withdraw" /> <input type="hidden" name="Amount" value="1000000" /> <input type="submit" value="Click to collect your prize!" /> </form>`- Notice that the form's - `action`posts to the vulnerable site, not to the malicious site. This is the "cross-site" part of CSRF.
- The user selects the submit button. The browser makes the request and automatically includes the authentication cookie for the requested domain, - `www.good-banking-site.example.com`.
- The request runs on the - `www.good-banking-site.example.com`server with the user's authentication context and can perform any action that an authenticated user is allowed to perform.

In addition to the scenario where the user selects the button to submit the form, the malicious site could:

- Run a script that automatically submits the form.
- Send the form submission as an AJAX request.
- Hide the form using CSS.

These alternative scenarios don't require any action or input from the user other than initially visiting the malicious site.

Using HTTPS doesn't prevent a CSRF attack. The malicious site can send `https://www.good-banking-site.example.com/` a request just as easily as it can send an insecure request.

Some attacks target endpoints that respond to GET requests, in which case an image tag can be used to perform the action. This form of attack is common on forum sites that permit images but block JavaScript. Apps that change state on GET requests, where variables or resources are altered, are vulnerable to malicious attacks. **GET requests that change state are insecure. A best practice is to never change state on a GET request.**

CSRF attacks are possible against web apps that use cookies for authentication because:

- Browsers store cookies issued by a web app.
- Stored cookies include session cookies for authenticated users.
- Browsers send all of the cookies associated with a domain to the web app every request regardless of how the request to app was generated within the browser.

However, CSRF attacks aren't limited to exploiting cookies. For example, Basic and Digest authentication are also vulnerable. After a user signs in with Basic or Digest authentication, the browser automatically sends the credentials until the session ends.

In this context, *session* refers to the client-side session during which the user is authenticated. It's unrelated to server-side sessions or ASP.NET Core session middleware.

Users can protect against CSRF vulnerabilities by taking precautions:

- Sign out of web apps when finished using them.
- Clear browser cookies periodically.

However, CSRF vulnerabilities are fundamentally a problem with the web app, not the end user.

## Authentication fundamentals

Cookie-based authentication is a popular form of authentication. Token-based authentication systems are growing in popularity, especially for Single Page Applications (SPAs).

### Cookie-based authentication

When a user authenticates using their username and password they're issued a token containing an authentication ticket. The token can be used for authentication and authorization. The token is stored as a cookie that's sent with every request the client makes. Generating and validating this cookie is performed with the cookie authentication middleware. The middleware serializes a user principal into an encrypted cookie. On subsequent requests, the middleware validates the cookie, recreates the principal, and assigns the principal to the HttpContext.User property.

### Token-based authentication

When a user is authenticated, they're issued a token (not an antiforgery token). The token contains user information in the form of claims or a reference token that points the app to user state maintained in the app. When a user attempts to access a resource that requires authentication, the token is sent to the app with an extra authorization header in the form of a Bearer token. This approach makes the app stateless. In each subsequent request, the token is passed in the request for server-side validation. This token isn't *encrypted*; it's *encoded*. On the server, the token is decoded to access its information. To send the token on subsequent requests, store the token in the browser's local storage. Placing a token in the browser local storage and retrieving it and using it as a bearer token provides protection against CSRF attacks. However, should the app be vulnerable to script injection via XSS or a compromised external JavaScript file, a cyberattacker could retrieve any value from local storage and send it to themselves. ASP.NET Core encodes all server side output from variables by default, reducing the risk of XSS. If you override this behavior by using Html.Raw or custom code with untrusted input then you may increase the risk of XSS.

Don't be concerned about CSRF vulnerability if the token is stored in the browser's local storage. CSRF is a concern when the token is stored in a cookie. For more information, see the GitHub issue SPA code sample adds two cookies.

### Multiple apps hosted at one domain

Shared hosting environments are vulnerable to session hijacking, sign-in CSRF, and other attacks.

Although `example1.contoso.net` and `example2.contoso.net` are different hosts, there's an implicit trust relationship between hosts under the `*.contoso.net` domain. This implicit trust relationship allows potentially untrusted hosts to affect each other's cookies (the same-origin policies that govern AJAX requests don't necessarily apply to HTTP cookies).

Attacks that exploit trusted cookies between apps hosted on the same domain can be prevented by not sharing domains. When each app is hosted on its own domain, there's no implicit cookie trust relationship to exploit.

### Fetch metadata headers

Modern browsers attach Fetch Metadata request headers—most importantly `Sec-Fetch-Site`—to every request. `Sec-Fetch-Site` describes the relationship between the origin that initiated the request and the origin being requested: `same-origin` identifies a request the site made to itself, while `same-site` and `cross-site` identify requests initiated by another origin. The `Origin` header carries the initiating origin and serves as a fallback for browsers that predate Fetch Metadata.

`Sec-Fetch-Site` and `Origin` are forbidden request headers: the browser sets them, and JavaScript running in a page can't override or forge them. That makes them a trustworthy signal for telling a site's own requests apart from cross-site requests without a server-issued token. The automatic CSRF protection built into ASP.NET Core uses this signal to reject cross-site form posts that aren't explicitly trusted.

## Automatic CSRF protection in ASP.NET Core

ASP.NET Core ships an automatic CSRF protection middleware that's enabled by default in apps built with `WebApplication.CreateBuilder`. Unlike the token-based antiforgery system, this middleware doesn't issue or validate tokens. Instead, it inspects the `Sec-Fetch-Site` and `Origin` Fetch metadata headers and records a validation verdict on the request. Components that process submitted form data enforce that verdict, rejecting cross-origin form posts that aren't explicitly trusted.

For most apps, no code changes are required: same-origin browser requests, safe HTTP methods, and non-browser clients (`curl`, server-to-server, mobile apps) all pass through unaffected. The middleware primarily affects apps that accept cross-origin form posts from a browser, such as a site that posts a form to an API on a different origin. Those scenarios need to either configure CORS to declare the trusted origin or opt the endpoint out.

This middleware is *additive* to the token-based antiforgery system. The two protections coexist and can both be active on the same endpoint. For a comparison of when each one applies, see Interaction with token-based antiforgery.

### How it works

For every request, the middleware evaluates a short chain of rules to reach a verdict—*allowed* or *denied*. The checks run in order, and the first match wins:

- **Safe HTTP methods are always allowed.**- `GET`,- `HEAD`,- `OPTIONS`, and- `TRACE`requests pass through. This follows RFC 9110 §9.2.1 and is consistent with the long-standing rule that endpoints shouldn't change state on- `GET`.
- `Sec-Fetch-Site: same-origin`or- `Sec-Fetch-Site: none`is allowed.- `Sec-Fetch-Site`on every request.- `same-origin`covers normal in-app navigation and fetch, and- `none`covers requests initiated directly by the user (typing a URL, using a bookmark). This is the most common code path—most legitimate browser traffic exits here.
- **A trusted origin from CORS is allowed.**If the request carries an- `Origin`header and the endpoint's resolved CORS policy trusts that origin, the request is allowed. The middleware resolves the policy the same way the CORS middleware does: per-endpoint policy from- `[EnableCors("name")]`first, then the default policy registered with- `AddDefaultPolicy`. See Allowing cross-origin clients for important limits on this rule.
- **Any other**When- `Sec-Fetch-Site`value is denied.- `Sec-Fetch-Site`is- `cross-site`or- `same-site`and the origin isn't trusted via CORS, the request is denied.
- **No**the middleware compares the- `Sec-Fetch-Site`, but- `Origin`is present:- `Origin`to- `scheme://host[:port]`built from the request. If they match, the request is allowed; otherwise it's denied. This is the fallback path for browsers older than the Fetch Metadata spec (released ~2020).
- **No**Browsers always send at least one of these on a write request, so a request missing both is almost certainly a non-browser client such as- `Sec-Fetch-Site`and no- `Origin`: the request is allowed.- `curl`, Postman, a mobile app, or a server-to-server caller. CSRF is a browser-only attack vector, so these requests pass through.

The middleware records this verdict on the request rather than ending the request itself. For how and when a denied verdict turns into an HTTP `400 Bad Request` response, see Deferred validation.

### Deferred validation

The middleware doesn't reject a request on its own. Instead, it records its verdict on the request's IAntiforgeryValidationFeature—the same feature the token-based antiforgery system uses—where a denied verdict is recorded as *invalid*. The request continues down the pipeline. An invalid verdict becomes an HTTP `400 Bad Request` only when a component that processes form data observes it. This deferral matches how the token-based system already behaves: the verdict is produced early but enforced at the point where a form is consumed.

The following components read IAntiforgeryValidationFeature and reject a request with `400 - Bad Request` when the recorded verdict is invalid:

- MVC actions protected by antiforgery.
- Minimal API endpoints that bind a form parameter.
- Blazor SSR endpoints.
- Any code that reads the request form directly, which acts as a backstop.

Each consumer first confirms that an antiforgery or CSRF middleware actually ran before it trusts the verdict, so a pipeline without either middleware doesn't produce false rejections.

A consequence of this model is that an endpoint that never reads form data runs even when the verdict is invalid. For example, a JSON API endpoint that binds its body from JSON, or a handler that ignores the request body, isn't rejected automatically on a cross-origin request. The verdict is still recorded on IAntiforgeryValidationFeature for code that wants to inspect it, but nothing enforces it. CSRF is a form-and-cookie attack vector, so endpoints that don't consume a browser-submitted form generally don't need this rejection. Endpoints that do process forms—Razor Pages, MVC views, Blazor SSR, and Minimal API form binding—get the protection automatically.

### Default behavior

The middleware is registered automatically by `WebApplication.CreateBuilder` and runs after authentication and authorization. It validates each request using the registered `ICsrfProtection` implementation, which by default applies the rules described in How it works. The default implementation can be replaced; see Customizing: implement `ICsrfProtection`. To turn the middleware off entirely, see Disabling globally.

The result is that a minimal app with a form-handling endpoint like the following is already protected:

```
var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();
app.MapPost("/widgets", ([FromForm] Widget w) =>
 Results.Created($"/widgets/{w.Id}", w));
app.Run();
```
A browser making a same-origin `POST /widgets` request reaches the endpoint normally. A browser at `https://attacker.example.com` posting the same form is rejected with `400 - Bad Request` when the endpoint binds the form, before the handler body runs. A `curl` request without `Sec-Fetch-Site` or `Origin` is allowed.

Because rejection is deferred to form consumers, an endpoint that doesn't read form data—such as a JSON API that binds its body from JSON—isn't rejected automatically, even on a cross-origin request. The verdict is still recorded on the request for code that wants to inspect it.

The middleware integrates with the existing antiforgery model:

- **Minimal APIs:**Calling- `.DisableAntiforgery()`on an endpoint opts that endpoint out of- *both*the token-based middleware and CSRF protection middleware. The same metadata (- `IAntiforgeryMetadata { RequiresValidation = false }`) is checked by both.
- **MVC controllers and actions:**- `[IgnoreAntiforgeryToken]`also opts the endpoint out of both protections.

### Allowing cross-origin clients

The most common scenario that requires action is a browser-based client that submits a cross-origin form—for example, a site at `https://app.contoso.com` posting a form to an API at `https://api.contoso.com`. Such form posts are denied by default because `Sec-Fetch-Site` is `same-site` or `cross-site` rather than `same-origin`, and the form consumer enforces that verdict with a `400 - Bad Request`.

The CSRF middleware doesn't introduce its own trust list. It reuses the same CORS policy that the CORS middleware resolves for the endpoint: if that policy allows the request's `Origin`, the CSRF middleware records an allowed verdict for the request.

The policy is picked per-endpoint in this order:

- `[EnableCors("api")]`(MVC) or- `.RequireCors("api")`(Minimal API) → the named policy- `"api"`.
- No CORS metadata on the endpoint → the default policy registered with `AddDefaultPolicy`.
- No matching policy (named policy not registered, no default policy, or `services.AddCors()`never called) → no CORS-derived trust. The middleware falls through to the`Sec-Fetch-Site`and Origin-vs-Host rules.

A minimal example using a default policy and a Minimal API endpoint:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddDefaultPolicy(policy =>
 policy.WithOrigins("https://app.contoso.com")
 .AllowAnyHeader()
 .AllowAnyMethod());
});
var app = builder.Build();
app.UseCors();
app.MapPost("/widgets", ([FromForm] Widget w) =>
 Results.Created($"/widgets/{w.Id}", w));
app.Run();
```
For a named policy on a single endpoint:

```
app.MapPost("/widgets", ([FromForm] Widget w) => Results.Created($"/widgets/{w.Id}", w))
 .RequireCors("api");
```
Warning

`AllowAnyOrigin` is intentionally **not** honored as a CSRF trust signal. `AllowAnyOrigin` means "any browser can read this resource," which is a different concern than "any origin may mutate state on the user's behalf." Treating `AllowAnyOrigin` as trusted would turn this middleware into a no-op for cross-origin writes. Apps that need a public-read CORS policy combined with CSRF-protected writes should list trusted write origins explicitly with `WithOrigins` or opt write endpoints out if they don't rely on cookie-based authentication.

`[DisableCors]` on an endpoint isn't a CSRF opt-out. It skips the CORS-derived trust step, and the request still has to satisfy the `Sec-Fetch-Site` and Origin-vs-Host rules. To opt out of CSRF protection, see Opting an endpoint out.

For details on configuring CORS itself—`AddCors`, `AddDefaultPolicy`, `AddPolicy`, `WithOrigins`, and the rest of the policy builder API—see Enable Cross-Origin Requests (CORS) in ASP.NET Core.

### Opting an endpoint out

If an endpoint isn't browser-reachable or is secured by a non-cookie mechanism, such as a bearer token or API key, opt it out individually rather than disabling the middleware globally.

**Minimal APIs** — call `DisableAntiforgery` on the endpoint or group:

```
app.MapPost("/api/webhook", (WebhookPayload p) => Results.Accepted())
 .DisableAntiforgery();
```
**MVC controllers** — apply `[IgnoreAntiforgeryToken]` to the action or controller:

```
[ApiController]
[Route("api/[controller]")]
[IgnoreAntiforgeryToken]
public class WebhookController : ControllerBase
{
 [HttpPost]
 public IActionResult Post([FromBody] WebhookPayload payload) => Accepted();
}
```
Either approach adds `IAntiforgeryMetadata { RequiresValidation = false }` to the endpoint, which the CSRF middleware honors by skipping validation.

Warning

Disabling CSRF protection on an endpoint should only be done when the endpoint isn't vulnerable to CSRF attacks—for example, endpoints that aren't callable from a browser or that are secured with non-cookie authentication such as bearer tokens or API keys. Don't disable CSRF protection on browser-accessible endpoints that rely on cookies for authentication.

### Disabling globally

The middleware can be disabled across the entire app using the `DisableCsrfProtection` configuration key. This is an escape hatch—prefer per-endpoint opt-outs.

In `appsettings.json`:

```
{
 "DisableCsrfProtection": true
}
```
Or as an environment variable:

```
ASPNETCORE_DisableCsrfProtection=true
```
When this key is set to `true`, `WebApplication` skips registering the middleware in the pipeline. The `ICsrfProtection` service remains registered, so anything that resolves it directly continues to work.

Warning

The automatic CSRF middleware also satisfies the antiforgery requirement for endpoints that require validation, even when an app doesn't call `app.UseAntiforgery()`. If an app relies on antiforgery but doesn't call `app.UseAntiforgery()`, disabling the CSRF middleware globally or running on a host that isn't built with `WebApplication`, where the middleware isn't injected, leaves those endpoints with no antiforgery middleware. A request to such an endpoint then throws an exception. Call `app.UseAntiforgery()` in that configuration.

### Browser support

`Sec-Fetch-Site` is supported by all current versions of Chromium-based browsers, Firefox, and Safari. For an authoritative compatibility table, see the MDN reference for `Sec-Fetch-Site`.

Older browsers that predate Fetch Metadata don't send `Sec-Fetch-Site`. For those clients, the middleware falls back to comparing the `Origin` header against the request's scheme and host. Browsers have sent `Origin` on cross-origin write requests for many years, so this fallback covers essentially all legacy browser traffic.

Non-browser clients—`curl`, Postman, mobile apps, server-to-server callers—typically send neither `Sec-Fetch-Site` nor `Origin`. Those requests are allowed because CSRF is a browser-only attack vector that depends on the browser automatically attaching ambient credentials such as cookies. A non-browser client that wants to attack the API has no need for CSRF; it can simply call the API directly with whatever credentials it possesses.

### Customizing: implement `ICsrfProtection`

The decision logic lives behind a one-method interface:

```
namespace Microsoft.AspNetCore.Antiforgery;
public interface ICsrfProtection
{
 ValueTask<CsrfProtectionResult> ValidateAsync(HttpContext context);
}
```
To replace the default implementation, register a singleton in DI. Because the framework uses `TryAddSingleton`, an explicit `AddSingleton` call overrides the default:

```
builder.Services.AddSingleton<ICsrfProtection, AllowlistCsrfProtection>();
```
A custom implementation is useful when the trust model doesn't fit CORS—for example, when a fixed allowlist of partner origins is preferred, or when stricter rules are needed:

```
using Microsoft.AspNetCore.Antiforgery;
using Microsoft.AspNetCore.Http;
public sealed class AllowlistCsrfProtection : ICsrfProtection
{
 private static readonly HashSet<string> SafeMethods =
 new(StringComparer.OrdinalIgnoreCase) { "GET", "HEAD", "OPTIONS", "TRACE" };
 private static readonly HashSet<string> TrustedOrigins =
 new(StringComparer.OrdinalIgnoreCase)
 {
 "https://app.contoso.com",
 "https://admin.contoso.com",
 };
 public ValueTask<CsrfProtectionResult> ValidateAsync(HttpContext context)
 {
 if (SafeMethods.Contains(context.Request.Method))
 {
 return ValueTask.FromResult(CsrfProtectionResult.Allowed());
 }
 var origin = context.Request.Headers.Origin.ToString();
 var allowed = !string.IsNullOrEmpty(origin) && TrustedOrigins.Contains(origin);
 return ValueTask.FromResult(
 allowed ? CsrfProtectionResult.Allowed() : CsrfProtectionResult.Denied());
 }
}
```
The middleware still honors `.DisableAntiforgery()` / `[IgnoreAntiforgeryToken]` regardless of which implementation is registered—opt-out is handled by the middleware itself before `ValidateAsync` is called.

### Interaction with token-based antiforgery

The two CSRF defenses target different layers and are designed to coexist. They also share the same request feature: both record their result on IAntiforgeryValidationFeature, and form consumers enforce whatever verdict is present.

| Aspect | Token-based `AntiforgeryMiddleware` | Automatic CSRF protection middleware |
|---|---|---|
| Introduced | ASP.NET Core 2.0+ | .NET 11 |
| Activation | Opt-in via `app.UseAntiforgery()`(or implicitly by`AddMvc`/`MapRazorPages`/`AddRazorComponents`) | Auto-injected by `WebApplication.CreateBuilder` |
| Validates | Synchronized token (form field + cookie pair) | `Sec-Fetch-Site`/`Origin`headers |
| Requires | ASP.NET Core Data Protection Overview for token encryption | No tokens, no state |
| Browser scope | All browsers that send cookies | All modern browsers; `Origin`fallback for legacy |
| Per-endpoint opt-out | `.DisableAntiforgery()`/`[IgnoreAntiforgeryToken]` | Same — both honor the same metadata |

The token-based middleware specifically protects against the classic CSRF attack pattern in which a malicious site triggers a form `POST` to a vulnerable site using the user's ambient cookies. The automatic CSRF middleware addresses the same threat at the HTTP layer using browser-supplied metadata. Both can be active on the same endpoint, and many apps will benefit from defense in depth:

- **Razor Pages, MVC, and Blazor SSR apps**that already use the token system gain a header-based check that runs before token validation, without changing the token flow.
- **Minimal API apps**that bind forms get a useful default without needing to call- `app.UseAntiforgery()`or thread- `IAntiforgery`through endpoints.
- **APIs called from cross-origin SPAs**can rely on this middleware combined with a CORS allowlist, and skip the token system entirely if the API never serves HTML forms.

The automatic CSRF middleware replaces the token-based system in many scenarios, because both protect the same form-handling endpoints. Keep the token-based system when:

- The app must support browsers that don't send `Sec-Fetch-Site`. See Browser support.
- The app uses IAntiforgeryAdditionalDataProvider to round-trip extra data inside the token.
- A security review or compliance requirement specifies the token defense as an independent layer.

For details on the token-based system, including form integration, AJAX flows, configuration via `AntiforgeryOptions`, and `IAntiforgery` APIs, see Antiforgery in ASP.NET Core.

#### Token validation takes precedence

When an app calls `app.UseAntiforgery()`, the token-based middleware runs after the automatic CSRF middleware. The token middleware clears any verdict the CSRF middleware recorded and replaces it with the result of token validation. The token result is authoritative:

- A request that the CSRF middleware marked invalid becomes valid if it carries a valid token.
- A request that the CSRF middleware allowed is marked invalid if its token is missing or invalid.

This ordering means apps that use the token system see the same end-to-end behavior they had before the automatic middleware existed, while apps that don't use tokens fall back to the CSRF middleware's verdict.

### Blazor static server-side rendering

Blazor static server-side rendering (SSR) endpoints participate in the same deferred model. The Razor Components endpoint trusts the verdict recorded on IAntiforgeryValidationFeature by the upstream middleware and returns `400 - Bad Request` for a form post only when that verdict is invalid. The endpoint no longer validates the request itself.

The behavior depends on which middleware ran:

- Apps that call `app.UseAntiforgery()`are unchanged. The token-based middleware validates each request, and antiforgery tokens are generated for rendered forms as before.
- Apps that don't call `app.UseAntiforgery()`are protected by the automatic CSRF middleware instead. In that configuration, the endpoint skips antiforgery token generation because no token middleware is present to validate a token on a later request.

This is a behavior change for static SSR that previously removed `app.UseAntiforgery()`: they're now protected by the CSRF middleware rather than left unprotected, and they stop emitting antiforgery tokens. For migration guidance, see Migrate from ASP.NET Core in .NET 10 to ASP.NET Core in .NET 11. For the formal breaking-change notice, see Blazor server-side rendering defers antiforgery validation to middleware.

### Troubleshooting

**Symptom:** Same-origin requests from a browser succeed, but cross-origin form posts return `400 - Bad Request` with no body.

**Cause:** The CSRF middleware recorded an invalid verdict for the cross-origin request, and a form-processing component—such as an MVC action, a Minimal API form binding, or a Blazor SSR form post—enforced that verdict with a `400 - Bad Request`. This is the expected default behavior for endpoints that process forms.

**Resolution:** Choose one of the following, depending on the scenario:

- If the calling origin is known and trusted, allow it via CORS.
- If the endpoint isn't browser-reachable or uses non-cookie authentication, opt it out with `.DisableAntiforgery()`or`[IgnoreAntiforgeryToken]`.
- If the entire app needs to opt out (for example, during a migration window), disable globally.

**Diagnosing:** The middleware logs every invalid verdict at `Debug` level under the category `Microsoft.AspNetCore.Antiforgery.CsrfProtectionMiddleware` with event name `CsrfValidationFailed`. Enable `Debug` logging for that category in `appsettings.Development.json`:

```
{
 "Logging": {
 "LogLevel": {
 "Microsoft.AspNetCore.Antiforgery.CsrfProtectionMiddleware": "Debug"
 }
 }
}
```
A recorded verdict then appears in the log as:

```
dbug: Microsoft.AspNetCore.Antiforgery.CsrfProtectionMiddleware[1]
 Cross-origin CSRF protection marked request POST /widgets from origin 'https://attacker.example.com' as invalid.
```
**Reproducing locally:** Use `curl` with an explicit `Origin` header to simulate a cross-origin browser request against a form endpoint:

```
curl -i -X POST https://localhost:{PORT}/widgets \
 -H "Origin: https://attacker.example.com" \
 -H "Content-Type: application/x-www-form-urlencoded" \
 -d "name=test"
```
Replace `{PORT}` with the app's local HTTPS port. The `400 - Bad Request` is observed because the endpoint binds the form, which enforces the recorded verdict. A non-form endpoint returns its normal response because nothing reads the verdict. Without the `Origin` header, the same request is allowed because `curl` doesn't send `Sec-Fetch-Site` either and a request with neither header is treated as a non-browser client.

The token-based antiforgery system described in the rest of this article predates this middleware and remains available. For most apps, the automatic protection is enough on its own. For guidance on when to keep the token-based system and how to migrate, see Migrate from ASP.NET Core in .NET 10 to ASP.NET Core in .NET 11.

## Antiforgery in ASP.NET Core

Warning

ASP.NET Core implements antiforgery using ASP.NET Core Data Protection. The data protection stack must be configured to work in a server farm. For more information, see Configuring data protection.

Antiforgery middleware is added to the Dependency injection container when one of the following APIs is called in `Program.cs`:

For more information, see Antiforgery with Minimal APIs.

The FormTagHelper injects antiforgery tokens into HTML form elements. The following markup in a Razor file automatically generates antiforgery tokens:

```
<form method="post">
 <!-- ... -->
</form>
```
Similarly, IHtmlHelper.BeginForm generates antiforgery tokens by default if the form's method isn't GET.

The automatic generation of antiforgery tokens for HTML form elements happens when the `<form>` tag contains the `method="post"` attribute and either of the following are true:

- The action attribute is empty (`action=""`).
- The action attribute isn't supplied (`<form method="post">`).

Automatic generation of antiforgery tokens for HTML form elements can be disabled:

- Explicitly disable antiforgery tokens with the - `asp-antiforgery`attribute:- `<form method="post" asp-antiforgery="false"> <!-- ... --> </form>`
- The form element is opted-out of Tag Helpers by using the Tag Helper ! opt-out symbol: - `<!form method="post"> <!-- ... --> </!form>`
- Remove the - `FormTagHelper`from the view. The- `FormTagHelper`can be removed from a view by adding the following directive to the Razor view:- `@removeTagHelper Microsoft.AspNetCore.Mvc.TagHelpers.FormTagHelper, Microsoft.AspNetCore.Mvc.TagHelpers`

Note

Razor Pages are automatically protected from XSRF/CSRF. For more information, see XSRF/CSRF and Razor Pages.

The most common approach to protecting against CSRF attacks is to use the *Synchronizer Token Pattern* (STP). STP is used when the user requests a page with form data:

- The server sends a token associated with the current user's identity to the client.
- The client sends back the token to the server for verification.
- If the server receives a token that doesn't match the authenticated user's identity, the request is rejected.

The token is unique and unpredictable. The token can also be used to ensure proper sequencing of a series of requests (for example, ensuring the request sequence of: page 1 > page 2 > page 3). All of the forms in ASP.NET Core MVC and Razor Pages templates generate antiforgery tokens. The following pair of view examples generates antiforgery tokens:

```
<form asp-action="Index" asp-controller="Home" method="post">
 <!-- ... -->
</form>
@using (Html.BeginForm("Index", "Home"))
{
 <!-- ... -->
}
```
Explicitly add an antiforgery token to a `<form>` element without using Tag Helpers with the HTML helper `@Html.AntiForgeryToken`:

```
<form asp-action="Index" asp-controller="Home" method="post">
 @Html.AntiForgeryToken()
 <!-- ... -->
</form>
```
In each of the preceding cases, ASP.NET Core adds a hidden form field similar to the following example:

```
<input name="__RequestVerificationToken" type="hidden" value="CfDJ8NrAkS ... s2-m9Yw">
```
ASP.NET Core includes three filters for working with antiforgery tokens:

## Antiforgery with `AddControllers`

Calling AddControllers does * not* enable antiforgery tokens. AddControllersWithViews must be called to have built-in antiforgery token support.

## Multiple browser tabs and the Synchronizer Token Pattern

Multiple tabs logged in as different users, or one logged in as anonymous, are not supported.

## Configure antiforgery with `AntiforgeryOptions`

Customize AntiforgeryOptions in the app's `Program` file:

```
builder.Services.AddAntiforgery(options =>
{
 // Set Cookie properties using CookieBuilder properties†.
 options.FormFieldName = "AntiforgeryFieldname";
 options.HeaderName = "X-CSRF-TOKEN-HEADERNAME";
 options.SuppressXFrameOptionsHeader = false;
});
```
Set the antiforgery cookie properties using the properties of the CookieBuilder class, as shown in the following table.

| Option | Description |
|---|---|
| Cookie | Determines the settings used to create the antiforgery cookies. |
| FormFieldName | The name of the hidden form field used by the antiforgery system to render antiforgery tokens in views. |
| HeaderName | The name of the header used by the antiforgery system. If `null`, the system considers only form data. |
| SuppressXFrameOptionsHeader | Specifies whether to suppress generation of the `X-Frame-Options`header. By default, the header is generated with a value of "SAMEORIGIN". Defaults to`false`. |

Some browsers don't allow insecure endpoints to set cookies with a 'secure' flag or overwrite cookies whose 'secure' flag is set (for more information, see Deprecate modification of 'secure' cookies from non-secure origins). Since mixing secure and insecure endpoints is a common scenario in apps, ASP.NET Core relaxes the restriction on the secure policy on some cookies, such as the antiforgery cookie, by setting the cookie's SecurePolicy to `CookieSecurePolicy.None`. Even if a malicious user steals an antiforgery cookie, they also must steal the antiforgery token that's typically sent via a form field (more common) or a separate request header (less common) plus the authentication cookie. Cookies related to authentication or authorization use a stronger policy than `CookieSecurePolicy.None`.

Optionally, you can secure the antiforgery cookie in non-`Development` environments using Secure Sockets Layer (SSL), over HTTPS only, with the following AntiforgeryOptions.Cookie property setting in the app's `Program` file:

```
if (!builder.Environment.IsDevelopment())
{
 builder.Services.AddAntiforgery(o =>
 {
 o.Cookie.SecurePolicy = CookieSecurePolicy.Always;
 });
}
```
For more information, see CookieAuthenticationOptions.

## Generate antiforgery tokens with `IAntiforgery`

IAntiforgery provides the API to configure antiforgery features. `IAntiforgery` can be requested in `Program.cs` using WebApplication.Services. The following example uses middleware from the app's home page to generate an antiforgery token and send it in the response as a cookie:

```
app.UseRouting();
app.UseAuthorization();
var antiforgery = app.Services.GetRequiredService<IAntiforgery>();
app.Use((context, next) =>
{
 var requestPath = context.Request.Path.Value;
 if (string.Equals(requestPath, "/", StringComparison.OrdinalIgnoreCase)
 || string.Equals(requestPath, "/index.html", StringComparison.OrdinalIgnoreCase))
 {
 var tokenSet = antiforgery.GetAndStoreTokens(context);
 context.Response.Cookies.Append("XSRF-TOKEN", tokenSet.RequestToken!,
 new CookieOptions { HttpOnly = false });
 }
 return next(context);
});
```
The preceding example sets a cookie named `XSRF-TOKEN`. The client can read this cookie and provide its value as a header attached to AJAX requests. For example, Angular includes built-in XSRF protection that reads a cookie named `XSRF-TOKEN` by default.

### Require antiforgery validation

The ValidateAntiForgeryToken action filter can be applied to an individual action, a controller, or globally. Requests made to actions that have this filter applied are blocked unless the request includes a valid antiforgery token:

```
[HttpPost]
[ValidateAntiForgeryToken]
public IActionResult Index()
{
 // ...
 return RedirectToAction();
}
```
The `ValidateAntiForgeryToken` attribute requires a token for requests to the action methods it marks, including HTTP GET requests. If the `ValidateAntiForgeryToken` attribute is applied across the app's controllers, it can be overridden with the `IgnoreAntiforgeryToken` attribute.

### Automatically validate antiforgery tokens for unsafe HTTP methods only

Instead of broadly applying the `ValidateAntiForgeryToken` attribute and then overriding it with `IgnoreAntiforgeryToken` attributes, the AutoValidateAntiforgeryToken attribute can be used. This attribute works identically to the `ValidateAntiForgeryToken` attribute, except that it doesn't require tokens for requests made using the following HTTP methods:

- GET
- HEAD
- OPTIONS
- TRACE

We recommend use of `AutoValidateAntiforgeryToken` broadly for non-API scenarios. This attribute ensures POST actions are protected by default. The alternative is to ignore antiforgery tokens by default, unless `ValidateAntiForgeryToken` is applied to individual action methods. It's more likely in this scenario for a POST action method to be left unprotected by mistake, leaving the app vulnerable to CSRF attacks. All POSTs should send the antiforgery token.

APIs don't have an automatic mechanism for sending the non-cookie part of the token. The implementation probably depends on the client code implementation. Some examples are shown below:

Class-level example:

```
[AutoValidateAntiforgeryToken]
public class HomeController : Controller
```
Global example:

```
builder.Services.AddControllersWithViews(options =>
{
 options.Filters.Add(new AutoValidateAntiforgeryTokenAttribute());
});
```
### Override global or controller antiforgery attributes

The IgnoreAntiforgeryToken filter is used to eliminate the need for an antiforgery token for a given action (or controller). When applied, this filter overrides `ValidateAntiForgeryToken` and `AutoValidateAntiforgeryToken` filters specified at a higher level (globally or on a controller).

```
[IgnoreAntiforgeryToken]
public IActionResult IndexOverride()
{
 // ...
 return RedirectToAction();
}
```
## Refresh tokens after authentication

Tokens should be refreshed after the user is authenticated by redirecting the user to a view or Razor Pages page.

## JavaScript, AJAX, and SPAs

In traditional HTML-based apps, antiforgery tokens are passed to the server using hidden form fields. In modern JavaScript-based apps and SPAs, many requests are made programmatically. These AJAX requests may use other techniques, such as request headers or cookies, to send the token.

If cookies are used to store authentication tokens and to authenticate API requests on the server, CSRF is a potential problem. If local storage is used to store the token, CSRF vulnerability might be mitigated because values from local storage aren't sent automatically to the server with every request. Using local storage to store the antiforgery token on the client and sending the token as a request header is a recommended approach.

### Blazor

For more information, see ASP.NET Core Blazor authentication and authorization.

### JavaScript

Using JavaScript with views, the token can be created using a service from within the view. Inject the IAntiforgery service into the view and call GetAndStoreTokens:

```
@inject Microsoft.AspNetCore.Antiforgery.IAntiforgery Antiforgery
@{
 ViewData["Title"] = "JavaScript";
 var requestToken = Antiforgery.GetAndStoreTokens(Context).RequestToken;
}
<input id="RequestVerificationToken" type="hidden" value="@requestToken" />
<button id="button" class="btn btn-primary">Submit with Token</button>
<div id="result" class="mt-2"></div>
@section Scripts {
<script>
 document.addEventListener("DOMContentLoaded", () => {
 const resultElement = document.getElementById("result");
 document.getElementById("button").addEventListener("click", async () => {
 const response = await fetch("@Url.Action("FetchEndpoint")", {
 method: "POST",
 headers: {
 RequestVerificationToken:
 document.getElementById("RequestVerificationToken").value
 }
 });
 if (response.ok) {
 resultElement.innerText = await response.text();
 } else {
 resultElement.innerText = `Request Failed: ${response.status}`
 }
 });
 });
</script>
}
```
The preceding example uses JavaScript to read the hidden field value for the AJAX POST header.

This approach eliminates the need to deal directly with setting cookies from the server or reading them from the client. However, when injecting the IAntiforgery service isn't possible, use JavaScript to access tokens in cookies:

- Access tokens in an additional request to the server, typically usually `same-origin`.
- Use the cookie's contents to create a header with the token's value.

Assuming the script sends the token in a request header called `X-XSRF-TOKEN`, configure the antiforgery service to look for the `X-XSRF-TOKEN` header:

```
builder.Services.AddAntiforgery(options => options.HeaderName = "X-XSRF-TOKEN");
```
The following example adds a protected endpoint that writes the request token to a JavaScript-readable cookie:

```
app.UseAuthorization();
app.MapGet("antiforgery/token", (IAntiforgery forgeryService, HttpContext context) =>
{
 var tokens = forgeryService.GetAndStoreTokens(context);
 context.Response.Cookies.Append("XSRF-TOKEN", tokens.RequestToken!,
 new CookieOptions { HttpOnly = false });
 return Results.Ok();
}).RequireAuthorization();
```
The following example uses JavaScript to make an AJAX request to obtain the token and make another request with the appropriate header:

```
var response = await fetch("/antiforgery/token", {
 method: "GET",
 headers: { "Authorization": authorizationToken }
});
if (response.ok) {
 // https://developer.mozilla.org/docs/web/api/document/cookie
 const xsrfToken = document.cookie
 .split("; ")
 .find(row => row.startsWith("XSRF-TOKEN="))
 .split("=")[1];
 response = await fetch("/JavaScript/FetchEndpoint", {
 method: "POST",
 headers: { "X-XSRF-TOKEN": xsrfToken, "Authorization": authorizationToken }
 });
 if (response.ok) {
 resultElement.innerText = await response.text();
 } else {
 resultElement.innerText = `Request Failed: ${response.status}`
 }
} else {
 resultElement.innerText = `Request Failed: ${response.status}`
}
```
Note

When the antiforgery token is provided in both the request header and in the form payload, only the token in the header is validated.

## Antiforgery with Minimal APIs

Call AddAntiforgery and UseAntiforgery(IApplicationBuilder) to register antiforgery services in DI. Antiforgery tokens are used to mitigate cross-site request forgery attacks.

```
var builder = WebApplication.CreateBuilder();
builder.Services.AddAntiforgery();
var app = builder.Build();
app.UseAntiforgery();
app.MapGet("/", () => "Hello World!");
app.Run();
```
The antiforgery middleware:

- Does **not**
- Sets the IAntiforgeryValidationFeature in the HttpContext.Features of the current request.

The antiforgery token is only validated if:

- The endpoint contains metadata implementing IAntiforgeryMetadata where `RequiresValidation=true`.
- The HTTP method associated with the endpoint is a relevant HTTP method of type POST, PUT, or PATCH.
- The request is associated with a valid endpoint.

Antiforgery middleware doesn't short-circuit the request pipeline. Endpoint code always runs, even if token validation fails. To observe the outcome of the token validation, resolve the IAntiforgeryValidationFeature from HttpContext.Features and inspect its IsValid property or the Error property for failure details. This approach is useful when endpoints require custom handling for failed antiforgery validation.

* Note:* When enabled manually, the antiforgery middleware must run after the authentication and authorization middleware to prevent reading form data when the user is unauthenticated.

By default, Minimal APIs that accept form data require antiforgery token validation and fail before running application code if antiforgery validation isn't successful.

Consider the following `GenerateForm` method:

```
public static string GenerateForm(string action,
 AntiforgeryTokenSet token, bool UseToken=true)
{
 string tokenInput = "";
 if (UseToken)
 {
 tokenInput = $@"<input name=""{token.FormFieldName}""
 type=""hidden"" value=""{token.RequestToken}"" />";
 }
 return $@"
 <html><body>
 <form action=""{action}"" method=""POST"" enctype=""multipart/form-data"">
 {tokenInput}
 <input type=""text"" name=""name"" />
 <input type=""date"" name=""dueDate"" />
 <input type=""checkbox"" name=""isCompleted"" />
 <input type=""submit"" />
 </form>
 </body></html>
";
}
```
The preceding code has three arguments, the action, the antiforgery token, and a `bool` indicating whether the token should be used.

Consider the following sample:

```
using Microsoft.AspNetCore.Antiforgery;
using Microsoft.AspNetCore.Mvc;
var builder = WebApplication.CreateBuilder();
builder.Services.AddAntiforgery();
var app = builder.Build();
app.UseAntiforgery();
// Pass token
app.MapGet("/", (HttpContext context, IAntiforgery antiforgery) =>
{
 var token = antiforgery.GetAndStoreTokens(context);
 return Results.Content(MyHtml.GenerateForm("/todo", token), "text/html");
});
// Don't pass a token, fails
app.MapGet("/SkipToken", (HttpContext context, IAntiforgery antiforgery) =>
{
 var token = antiforgery.GetAndStoreTokens(context);
 return Results.Content(MyHtml.GenerateForm("/todo",token, false ), "text/html");
});
// Post to /todo2. DisableAntiforgery on that endpoint so no token needed.
app.MapGet("/DisableAntiforgery", (HttpContext context, IAntiforgery antiforgery) =>
{
 var token = antiforgery.GetAndStoreTokens(context);
 return Results.Content(MyHtml.GenerateForm("/todo2", token, false), "text/html");
});
app.MapPost("/todo", ([FromForm] Todo todo) => Results.Ok(todo));
app.MapPost("/todo2", ([FromForm] Todo todo) => Results.Ok(todo))
 .DisableAntiforgery();
app.Run();
class Todo
{
 public required string Name { get; set; }
 public bool IsCompleted { get; set; }
 public DateTime DueDate { get; set; }
}
public static class MyHtml
{
 public static string GenerateForm(string action,
 AntiforgeryTokenSet token, bool UseToken=true)
 {
 string tokenInput = "";
 if (UseToken)
 {
 tokenInput = $@"<input name=""{token.FormFieldName}""
 type=""hidden"" value=""{token.RequestToken}"" />";
 }
 return $@"
 <html><body>
 <form action=""{action}"" method=""POST"" enctype=""multipart/form-data"">
 {tokenInput}
 <input type=""text"" name=""name"" />
 <input type=""date"" name=""dueDate"" />
 <input type=""checkbox"" name=""isCompleted"" />
 <input type=""submit"" />
 </form>
 </body></html>
 ";
 }
}
```
In the preceding code, posts to:

- `/todo`require a valid antiforgery token.
- `/todo2`do- **not**

```
app.MapPost("/todo", ([FromForm] Todo todo) => Results.Ok(todo));
app.MapPost("/todo2", ([FromForm] Todo todo) => Results.Ok(todo))
 .DisableAntiforgery();
```
Warning

Calling `.DisableAntiforgery()` disables cross-site request forgery (CSRF) protection for the endpoint. This should only be used when an endpoint is not vulnerable to CSRF attacks, such as:

- Endpoints that are not callable from a browser (for example, internal APIs)
- Endpoints secured with non-cookie-based authentication (for example, bearer tokens or API keys)
- Internal or infrastructure endpoints that do not rely on user cookies

Do **not** disable antiforgery validation for browser-accessible endpoints that rely on cookies for authentication or that process user-submitted form data, as this exposes your application to CSRF attacks.

A POST to:

- `/todo`from the form generated by the- `/`endpoint succeeds because the antiforgery token is valid.
- `/todo`from the form generated by the- `/SkipToken`fails because the antiforgery is not included.
- `/todo2`from the form generated by the- `/DisableAntiforgery`endpoint succeeds because the antiforgery is not required.

```
app.MapPost("/todo", ([FromForm] Todo todo) => Results.Ok(todo));
app.MapPost("/todo2", ([FromForm] Todo todo) => Results.Ok(todo))
 .DisableAntiforgery();
```
When a form is submitted without a valid antiforgery token:

- In the `Development`environment, an exception is thrown.
- In the `Production`environment, a message is logged.

### HTTP method limitations and `HttpMethodOverrideMiddleware` interaction

For the middleware-based antiforgery path, `AntiforgeryMiddleware` and `UseAntiforgery()` validate antiforgery tokens only for HTTP **POST**, **PUT**, and **PATCH** requests. Other HTTP methods, such as **DELETE**, aren't validated automatically.

To validate antiforgery tokens for other HTTP methods, resolve IAntiforgery from DI and call ValidateRequestAsync or IsRequestValidAsync explicitly:

```
app.MapDelete("/item/{id}", async (int id, IAntiforgery antiforgery, HttpContext context) =>
{
 await antiforgery.ValidateRequestAsync(context);
 // Process the DELETE request
});
```
Warning

When `HttpMethodOverrideMiddleware` is configured with `FormFieldName` (form-field mode) and placed before `AntiforgeryMiddleware`, a POST request can be overridden to DELETE (or another non-validated method). Because `AntiforgeryMiddleware` validates only POST, PUT, and PATCH, the overridden request bypasses antiforgery validation.

To protect these endpoints:

- Prefer placing `HttpMethodOverrideMiddleware`after antiforgery validation when your pipeline allows it.
- Avoid form-field overrides for endpoints that rely on antiforgery validation.
- If form-field override must run first, validate the antiforgery token explicitly using `IAntiforgery.ValidateRequestAsync`.

For configuration details, see ASP.NET Core middleware.

## Windows authentication and antiforgery cookies

When using Windows Authentication, application endpoints must be protected against CSRF attacks in the same way as done for cookies. The browser implicitly sends the authentication context to the server and endpoints need to be protected against CSRF attacks.

## Extend antiforgery

The IAntiforgeryAdditionalDataProvider type allows developers to extend the behavior of the anti-CSRF system by round-tripping additional data in each token. The GetAdditionalData method is called each time a field token is generated, and the return value is embedded within the generated token. An implementer could return a timestamp, a nonce, or any other value and then call ValidateAdditionalData to validate this data when the token is validated. The client's username is already embedded in the generated tokens, so there's no need to include this information. If a token includes supplemental data but no `IAntiForgeryAdditionalDataProvider` is configured, the supplemental data isn't validated.

## Additional resources

Cross-site request forgery (also known as XSRF or CSRF) is an attack against web-hosted apps whereby a malicious web app can influence the interaction between a client browser and a web app that trusts that browser. These attacks are possible because web browsers send some types of authentication tokens automatically with every request to a website. This form of exploit is also known as a *one-click attack* or *session riding* because the attack takes advantage of the user's previously authenticated session.

An example of a CSRF attack:

- A user signs into - `www.good-banking-site.example.com`using forms authentication. The server authenticates the user and issues a response that includes an authentication cookie. The site is vulnerable to attack because it trusts any request that it receives with a valid authentication cookie.
- The user visits a malicious site, - `www.bad-crook-site.example.com`.- The malicious site, - `www.bad-crook-site.example.com`, contains an HTML form similar to the following example:- `<h1>Congratulations! You're a Winner!</h1> <form action="https://www.good-banking-site.example.com/api/account" method="post"> <input type="hidden" name="Transaction" value="withdraw" /> <input type="hidden" name="Amount" value="1000000" /> <input type="submit" value="Click to collect your prize!" /> </form>`- Notice that the form's - `action`posts to the vulnerable site, not to the malicious site. This is the "cross-site" part of CSRF.
- The user selects the submit button. The browser makes the request and automatically includes the authentication cookie for the requested domain, - `www.good-banking-site.example.com`.
- The request runs on the - `www.good-banking-site.example.com`server with the user's authentication context and can perform any action that an authenticated user is allowed to perform.

In addition to the scenario where the user selects the button to submit the form, the malicious site could:

- Run a script that automatically submits the form.
- Send the form submission as an AJAX request.
- Hide the form using CSS.

These alternative scenarios don't require any action or input from the user other than initially visiting the malicious site.

Using HTTPS doesn't prevent a CSRF attack. The malicious site can send `https://www.good-banking-site.example.com/` a request just as easily as it can send an insecure request.

Some attacks target endpoints that respond to GET requests, in which case an image tag can be used to perform the action. This form of attack is common on forum sites that permit images but block JavaScript. Apps that change state on GET requests, where variables or resources are altered, are vulnerable to malicious attacks. **GET requests that change state are insecure. A best practice is to never change state on a GET request.**

CSRF attacks are possible against web apps that use cookies for authentication because:

- Browsers store cookies issued by a web app.
- Stored cookies include session cookies for authenticated users.
- Browsers send all of the cookies associated with a domain to the web app every request regardless of how the request to app was generated within the browser.

However, CSRF attacks aren't limited to exploiting cookies. For example, Basic and Digest authentication are also vulnerable. After a user signs in with Basic or Digest authentication, the browser automatically sends the credentials until the session ends.

In this context, *session* refers to the client-side session during which the user is authenticated. It's unrelated to server-side sessions or ASP.NET Core session middleware.

Users can protect against CSRF vulnerabilities by taking precautions:

- Sign out of web apps when finished using them.
- Clear browser cookies periodically.

However, CSRF vulnerabilities are fundamentally a problem with the web app, not the end user.

## Authentication fundamentals

Cookie-based authentication is a popular form of authentication. Token-based authentication systems are growing in popularity, especially for Single Page Applications (SPAs).

### Cookie-based authentication

When a user authenticates using their username and password, they're issued a token, containing an authentication ticket that can be used for authentication and authorization. The token is stored as a cookie that's sent with every request the client makes. Generating and validating this cookie is performed by the cookie authentication middleware. The middleware serializes a user principal into an encrypted cookie. On subsequent requests, the middleware validates the cookie, recreates the principal, and assigns the principal to the HttpContext.User property.

### Token-based authentication

When a user is authenticated, they're issued a token (not an antiforgery token). The token contains user information in the form of claims or a reference token that points the app to user state maintained in the app. When a user attempts to access a resource that requires authentication, the token is sent to the app with an extra authorization header in the form of a Bearer token. This approach makes the app stateless. In each subsequent request, the token is passed in the request for server-side validation. This token isn't *encrypted*; it's *encoded*. On the server, the token is decoded to access its information. To send the token on subsequent requests, store the token in the browser's local storage. Placing a token in the browser local storage and retrieving it and using it as a bearer token provides protection against CSRF attacks. However, should the app be vulnerable to script injection via XSS or a compromised external javascript file, a cyberattacker could retrieve any value from local storage and send it to themselves. ASP.NET Core encodes all server side output from variables by default, reducing the risk of XSS. If you override this behavior by using Html.Raw or custom code with untrusted input then you may increase the risk of XSS.

Don't be concerned about CSRF vulnerability if the token is stored in the browser's local storage. CSRF is a concern when the token is stored in a cookie. For more information, see the GitHub issue SPA code sample adds two cookies.

### Multiple apps hosted at one domain

Shared hosting environments are vulnerable to session hijacking, login CSRF, and other attacks.

Although `example1.contoso.net` and `example2.contoso.net` are different hosts, there's an implicit trust relationship between hosts under the `*.contoso.net` domain. This implicit trust relationship allows potentially untrusted hosts to affect each other's cookies (the same-origin policies that govern AJAX requests don't necessarily apply to HTTP cookies).

Attacks that exploit trusted cookies between apps hosted on the same domain can be prevented by not sharing domains. When each app is hosted on its own domain, there's no implicit cookie trust relationship to exploit.

## Antiforgery in ASP.NET Core

Warning

ASP.NET Core implements antiforgery using ASP.NET Core Data Protection. The data protection stack must be configured to work in a server farm. For more information, see Configuring data protection.

Antiforgery middleware is added to the Dependency injection container when one of the following APIs is called in `Program.cs`:

The FormTagHelper injects antiforgery tokens into HTML form elements. The following markup in a Razor file automatically generates antiforgery tokens:

```
<form method="post">
 <!-- ... -->
</form>
```
Similarly, IHtmlHelper.BeginForm generates antiforgery tokens by default if the form's method isn't GET.

The automatic generation of antiforgery tokens for HTML form elements happens when the `<form>` tag contains the `method="post"` attribute and either of the following are true:

- The action attribute is empty (`action=""`).
- The action attribute isn't supplied (`<form method="post">`).

Automatic generation of antiforgery tokens for HTML form elements can be disabled:

- Explicitly disable antiforgery tokens with the - `asp-antiforgery`attribute:- `<form method="post" asp-antiforgery="false"> <!-- ... --> </form>`
- The form element is opted-out of Tag Helpers by using the Tag Helper ! opt-out symbol: - `<!form method="post"> <!-- ... --> </!form>`
- Remove the - `FormTagHelper`from the view. The- `FormTagHelper`can be removed from a view by adding the following directive to the Razor view:- `@removeTagHelper Microsoft.AspNetCore.Mvc.TagHelpers.FormTagHelper, Microsoft.AspNetCore.Mvc.TagHelpers`

Note

Razor Pages are automatically protected from XSRF/CSRF. For more information, see XSRF/CSRF and Razor Pages.

The most common approach to defending against CSRF attacks is to use the *Synchronizer Token Pattern* (STP). STP is used when the user requests a page with form data:

- The server sends a token associated with the current user's identity to the client.
- The client sends back the token to the server for verification.
- If the server receives a token that doesn't match the authenticated user's identity, the request is rejected.

The token is unique and unpredictable. The token can also be used to ensure proper sequencing of a series of requests (for example, ensuring the request sequence of: page 1 > page 2 > page 3). All of the forms in ASP.NET Core MVC and Razor Pages templates generate antiforgery tokens. The following pair of view examples generates antiforgery tokens:

```
<form asp-action="Index" asp-controller="Home" method="post">
 <!-- ... -->
</form>
@using (Html.BeginForm("Index", "Home"))
{
 <!-- ... -->
}
```
Explicitly add an antiforgery token to a `<form>` element without using Tag Helpers with the HTML helper `@Html.AntiForgeryToken`:

```
<form asp-action="Index" asp-controller="Home" method="post">
 @Html.AntiForgeryToken()
 <!-- ... -->
</form>
```
In each of the preceding cases, ASP.NET Core adds a hidden form field similar to the following example:

```
<input name="__RequestVerificationToken" type="hidden" value="CfDJ8NrAkS ... s2-m9Yw">
```
ASP.NET Core includes three filters for working with antiforgery tokens:

## Antiforgery with `AddControllers`

Calling AddControllers does * not* enable antiforgery tokens. AddControllersWithViews must be called to have built-in antiforgery token support.

## Multiple browser tabs and the Synchronizer Token Pattern

With the Synchronizer Token Pattern, only the most recently loaded page contains a valid antiforgery token. Using multiple tabs can be problematic. For example, if a user opens multiple tabs:

- Only the most recently loaded tab contains a valid antiforgery token.
- Requests made from previously loaded tabs fail with an error: `Antiforgery token validation failed. The antiforgery cookie token and request token do not match`

Consider alternative CSRF protection patterns if this poses an issue.

## Configure antiforgery with `AntiforgeryOptions`

Customize AntiforgeryOptions in the app's `Program` file:

```
builder.Services.AddAntiforgery(options =>
{
 // Set Cookie properties using CookieBuilder properties†.
 options.FormFieldName = "AntiforgeryFieldname";
 options.HeaderName = "X-CSRF-TOKEN-HEADERNAME";
 options.SuppressXFrameOptionsHeader = false;
});
```
Set the antiforgery cookie properties using the properties of the CookieBuilder class, as shown in the following table.

| Option | Description |
|---|---|
| Cookie | Determines the settings used to create the antiforgery cookies. |
| FormFieldName | The name of the hidden form field used by the antiforgery system to render antiforgery tokens in views. |
| HeaderName | The name of the header used by the antiforgery system. If `null`, the system considers only form data. |
| SuppressXFrameOptionsHeader | Specifies whether to suppress generation of the `X-Frame-Options`header. By default, the header is generated with a value of "SAMEORIGIN". Defaults to`false`. |

Some browsers don't allow insecure endpoints to set cookies with a 'secure' flag or overwrite cookies whose 'secure' flag is set (for more information, see Deprecate modification of 'secure' cookies from non-secure origins). Since mixing secure and insecure endpoints is a common scenario in apps, ASP.NET Core relaxes the restriction on the secure policy on some cookies, such as the antiforgery cookie, by setting the cookie's SecurePolicy to `CookieSecurePolicy.None`. Even if a malicious user steals an antiforgery cookie, they also must steal the antiforgery token that's typically sent via a form field (more common) or a separate request header (less common) plus the authentication cookie. Cookies related to authentication or authorization use a stronger policy than `CookieSecurePolicy.None`.

Optionally, you can secure the antiforgery cookie in non-`Development` environments using Secure Sockets Layer (SSL), over HTTPS only, with the following AntiforgeryOptions.Cookie property setting in the app's `Program` file:

```
if (!builder.Environment.IsDevelopment())
{
 builder.Services.AddAntiforgery(o =>
 {
 o.Cookie.SecurePolicy = CookieSecurePolicy.Always;
 });
}
```
For more information, see CookieAuthenticationOptions.

## Generate antiforgery tokens with `IAntiforgery`

IAntiforgery provides the API to configure antiforgery features. `IAntiforgery` can be requested in `Program.cs` using WebApplication.Services. The following example uses middleware from the app's home page to generate an antiforgery token and send it in the response as a cookie:

```
app.UseRouting();
app.UseAuthorization();
var antiforgery = app.Services.GetRequiredService<IAntiforgery>();
app.Use((context, next) =>
{
 var requestPath = context.Request.Path.Value;
 if (string.Equals(requestPath, "/", StringComparison.OrdinalIgnoreCase)
 || string.Equals(requestPath, "/index.html", StringComparison.OrdinalIgnoreCase))
 {
 var tokenSet = antiforgery.GetAndStoreTokens(context);
 context.Response.Cookies.Append("XSRF-TOKEN", tokenSet.RequestToken!,
 new CookieOptions { HttpOnly = false });
 }
 return next(context);
});
```
The preceding example sets a cookie named `XSRF-TOKEN`. The client can read this cookie and provide its value as a header attached to AJAX requests. For example, Angular includes built-in XSRF protection that reads a cookie named `XSRF-TOKEN` by default.

### Require antiforgery validation

The ValidateAntiForgeryToken action filter can be applied to an individual action, a controller, or globally. Requests made to actions that have this filter applied are blocked unless the request includes a valid antiforgery token:

```
[HttpPost]
[ValidateAntiForgeryToken]
public IActionResult Index()
{
 // ...
 return RedirectToAction();
}
```
The `ValidateAntiForgeryToken` attribute requires a token for requests to the action methods it marks, including HTTP GET requests. If the `ValidateAntiForgeryToken` attribute is applied across the app's controllers, it can be overridden with the `IgnoreAntiforgeryToken` attribute.

### Automatically validate antiforgery tokens for unsafe HTTP methods only

Instead of broadly applying the `ValidateAntiForgeryToken` attribute and then overriding it with `IgnoreAntiforgeryToken` attributes, the AutoValidateAntiforgeryToken attribute can be used. This attribute works identically to the `ValidateAntiForgeryToken` attribute, except that it doesn't require tokens for requests made using the following HTTP methods:

- GET
- HEAD
- OPTIONS
- TRACE

We recommend use of `AutoValidateAntiforgeryToken` broadly for non-API scenarios. This attribute ensures POST actions are protected by default. The alternative is to ignore antiforgery tokens by default, unless `ValidateAntiForgeryToken` is applied to individual action methods. It's more likely in this scenario for a POST action method to be left unprotected by mistake, leaving the app vulnerable to CSRF attacks. All POSTs should send the antiforgery token.

APIs don't have an automatic mechanism for sending the non-cookie part of the token. The implementation probably depends on the client code implementation. Some examples are shown below:

Class-level example:

```
[AutoValidateAntiforgeryToken]
public class HomeController : Controller
```
Global example:

```
builder.Services.AddControllersWithViews(options =>
{
 options.Filters.Add(new AutoValidateAntiforgeryTokenAttribute());
});
```
### Override global or controller antiforgery attributes

The IgnoreAntiforgeryToken filter is used to eliminate the need for an antiforgery token for a given action (or controller). When applied, this filter overrides `ValidateAntiForgeryToken` and `AutoValidateAntiforgeryToken` filters specified at a higher level (globally or on a controller).

```
[IgnoreAntiforgeryToken]
public IActionResult IndexOverride()
{
 // ...
 return RedirectToAction();
}
```
## Refresh tokens after authentication

Tokens should be refreshed after the user is authenticated by redirecting the user to a view or Razor Pages page.

## JavaScript, AJAX, and SPAs

In traditional HTML-based apps, antiforgery tokens are passed to the server using hidden form fields. In modern JavaScript-based apps and SPAs, many requests are made programmatically. These AJAX requests may use other techniques (such as request headers or cookies) to send the token.

If cookies are used to store authentication tokens and to authenticate API requests on the server, CSRF is a potential problem. If local storage is used to store the token, CSRF vulnerability might be mitigated because values from local storage aren't sent automatically to the server with every request. Using local storage to store the antiforgery token on the client and sending the token as a request header is a recommended approach.

### JavaScript

Using JavaScript with views, the token can be created using a service from within the view. Inject the IAntiforgery service into the view and call GetAndStoreTokens:

```
@inject Microsoft.AspNetCore.Antiforgery.IAntiforgery Antiforgery
@{
 ViewData["Title"] = "JavaScript";
 var requestToken = Antiforgery.GetAndStoreTokens(Context).RequestToken;
}
<input id="RequestVerificationToken" type="hidden" value="@requestToken" />
<button id="button" class="btn btn-primary">Submit with Token</button>
<div id="result" class="mt-2"></div>
@section Scripts {
<script>
 document.addEventListener("DOMContentLoaded", () => {
 const resultElement = document.getElementById("result");
 document.getElementById("button").addEventListener("click", async () => {
 const response = await fetch("@Url.Action("FetchEndpoint")", {
 method: "POST",
 headers: {
 RequestVerificationToken:
 document.getElementById("RequestVerificationToken").value
 }
 });
 if (response.ok) {
 resultElement.innerText = await response.text();
 } else {
 resultElement.innerText = `Request Failed: ${response.status}`
 }
 });
 });
</script>
}
```
The preceding example uses JavaScript to read the hidden field value for the AJAX POST header.

This approach eliminates the need to deal directly with setting cookies from the server or reading them from the client. However, when injecting the IAntiforgery service isn't possible, use JavaScript to access tokens in cookies:

- Access tokens in an additional request to the server, typically usually `same-origin`.
- Use the cookie's contents to create a header with the token's value.

Assuming the script sends the token in a request header called `X-XSRF-TOKEN`, configure the antiforgery service to look for the `X-XSRF-TOKEN` header:

```
builder.Services.AddAntiforgery(options => options.HeaderName = "X-XSRF-TOKEN");
```
The following example adds a protected endpoint that writes the request token to a JavaScript-readable cookie:

```
app.UseAuthorization();
app.MapGet("antiforgery/token", (IAntiforgery forgeryService, HttpContext context) =>
{
 var tokens = forgeryService.GetAndStoreTokens(context);
 context.Response.Cookies.Append("XSRF-TOKEN", tokens.RequestToken!,
 new CookieOptions { HttpOnly = false });
 return Results.Ok();
}).RequireAuthorization();
```
The following example uses JavaScript to make an AJAX request to obtain the token and make another request with the appropriate header:

```
var response = await fetch("/antiforgery/token", {
 method: "GET",
 headers: { "Authorization": authorizationToken }
});
if (response.ok) {
 // https://developer.mozilla.org/docs/web/api/document/cookie
 const xsrfToken = document.cookie
 .split("; ")
 .find(row => row.startsWith("XSRF-TOKEN="))
 .split("=")[1];
 response = await fetch("/JavaScript/FetchEndpoint", {
 method: "POST",
 headers: { "X-XSRF-TOKEN": xsrfToken, "Authorization": authorizationToken }
 });
 if (response.ok) {
 resultElement.innerText = await response.text();
 } else {
 resultElement.innerText = `Request Failed: ${response.status}`
 }
} else {
 resultElement.innerText = `Request Failed: ${response.status}`
}
```
Note

When the antiforgery token is provided in both the request header and in the form payload, only the token in the header is validated.

### Antiforgery with Minimal APIs

`Minimal APIs` do not support the usage of the included filters (`ValidateAntiForgeryToken`, `AutoValidateAntiforgeryToken`, `IgnoreAntiforgeryToken`), however, IAntiforgery provides the required APIs to validate a request.

The following example creates a filter that validates the antiforgery token:

```
internal static class AntiForgeryExtensions
{
 public static TBuilder ValidateAntiforgery<TBuilder>(this TBuilder builder) where TBuilder : IEndpointConventionBuilder
 {
 return builder.AddEndpointFilter(routeHandlerFilter: async (context, next) =>
 {
 try
 {
 var antiForgeryService = context.HttpContext.RequestServices.GetRequiredService<IAntiforgery>();
 await antiForgeryService.ValidateRequestAsync(context.HttpContext);
 }
 catch (AntiforgeryValidationException)
 {
 return Results.BadRequest("Antiforgery token validation failed.");
 }
 return await next(context);
 });
 }
}
```
The filter can then be applied to an endpoint:

```
app.MapPost("api/upload", (IFormFile name) => Results.Accepted())
 .RequireAuthorization()
 .ValidateAntiforgery();
```
## Windows authentication and antiforgery cookies

When using Windows Authentication, application endpoints must be protected against CSRF attacks in the same way as done for cookies. The browser implicitly sends the authentication context to the server and endpoints need to be protected against CSRF attacks.

## Extend antiforgery

The IAntiforgeryAdditionalDataProvider type allows developers to extend the behavior of the anti-CSRF system by round-tripping additional data in each token. The GetAdditionalData method is called each time a field token is generated, and the return value is embedded within the generated token. An implementer could return a timestamp, a nonce, or any other value and then call ValidateAdditionalData to validate this data when the token is validated. The client's username is already embedded in the generated tokens, so there's no need to include this information. If a token includes supplemental data but no `IAntiForgeryAdditionalDataProvider` is configured, the supplemental data isn't validated.

## Additional resources

Cross-site request forgery (also known as XSRF or CSRF) is an attack against web-hosted apps whereby a malicious web app can influence the interaction between a client browser and a web app that trusts that browser. These attacks are possible because web browsers send some types of authentication tokens automatically with every request to a website. This form of exploit is also known as a *one-click attack* or *session riding* because the attack takes advantage of the user's previously authenticated session.

An example of a CSRF attack:

- A user signs into - `www.good-banking-site.example.com`using forms authentication. The server authenticates the user and issues a response that includes an authentication cookie. The site is vulnerable to attack because it trusts any request that it receives with a valid authentication cookie.
- The user visits a malicious site, - `www.bad-crook-site.example.com`.- The malicious site, - `www.bad-crook-site.example.com`, contains an HTML form similar to the following example:- `<h1>Congratulations! You're a Winner!</h1> <form action="https://www.good-banking-site.example.com/api/account" method="post"> <input type="hidden" name="Transaction" value="withdraw" /> <input type="hidden" name="Amount" value="1000000" /> <input type="submit" value="Click to collect your prize!" /> </form>`- Notice that the form's - `action`posts to the vulnerable site, not to the malicious site. This is the "cross-site" part of CSRF.
- The user selects the submit button. The browser makes the request and automatically includes the authentication cookie for the requested domain, - `www.good-banking-site.example.com`.
- The request runs on the - `www.good-banking-site.example.com`server with the user's authentication context and can perform any action that an authenticated user is allowed to perform.

In addition to the scenario where the user selects the button to submit the form, the malicious site could:

- Run a script that automatically submits the form.
- Send the form submission as an AJAX request.
- Hide the form using CSS.

These alternative scenarios don't require any action or input from the user other than initially visiting the malicious site.

Using HTTPS doesn't prevent a CSRF attack. The malicious site can send `https://www.good-banking-site.example.com/` a request just as easily as it can send an insecure request.

Some attacks target endpoints that respond to GET requests, in which case an image tag can be used to perform the action. This form of attack is common on forum sites that permit images but block JavaScript. Apps that change state on GET requests, where variables or resources are altered, are vulnerable to malicious attacks. **GET requests that change state are insecure. A best practice is to never change state on a GET request.**

CSRF attacks are possible against web apps that use cookies for authentication because:

- Browsers store cookies issued by a web app.
- Stored cookies include session cookies for authenticated users.
- Browsers send all of the cookies associated with a domain to the web app every request regardless of how the request to app was generated within the browser.

However, CSRF attacks aren't limited to exploiting cookies. For example, Basic and Digest authentication are also vulnerable. After a user signs in with Basic or Digest authentication, the browser automatically sends the credentials until the session ends.

In this context, *session* refers to the client-side session during which the user is authenticated. It's unrelated to server-side sessions or ASP.NET Core session middleware.

Users can protect against CSRF vulnerabilities by taking precautions:

- Sign out of web apps when finished using them.
- Clear browser cookies periodically.

However, CSRF vulnerabilities are fundamentally a problem with the web app, not the end user.

## Authentication fundamentals

Cookie-based authentication is a popular form of authentication. Token-based authentication systems are growing in popularity, especially for Single Page Applications (SPAs).

### Cookie-based authentication

When a user authenticates using their username and password, they're issued a token, containing an authentication ticket that can be used for authentication and authorization. The token is stored as a cookie that's sent with every request the client makes. Generating and validating this cookie is performed by the cookie authentication middleware. The middleware serializes a user principal into an encrypted cookie. On subsequent requests, the middleware validates the cookie, recreates the principal, and assigns the principal to the HttpContext.User property.

### Token-based authentication

When a user is authenticated, they're issued a token (not an antiforgery token). The token contains user information in the form of claims or a reference token that points the app to user state maintained in the app. When a user attempts to access a resource that requires authentication, the token is sent to the app with an extra authorization header in the form of a Bearer token. This approach makes the app stateless. In each subsequent request, the token is passed in the request for server-side validation. This token isn't *encrypted*; it's *encoded*. On the server, the token is decoded to access its information. To send the token on subsequent requests, store the token in the browser's local storage. Don't be concerned about CSRF vulnerability if the token is stored in the browser's local storage. CSRF is a concern when the token is stored in a cookie. For more information, see the GitHub issue SPA code sample adds two cookies.

### Multiple apps hosted at one domain

Shared hosting environments are vulnerable to session hijacking, login CSRF, and other attacks.

Although `example1.contoso.net` and `example2.contoso.net` are different hosts, there's an implicit trust relationship between hosts under the `*.contoso.net` domain. This implicit trust relationship allows potentially untrusted hosts to affect each other's cookies (the same-origin policies that govern AJAX requests don't necessarily apply to HTTP cookies).

Attacks that exploit trusted cookies between apps hosted on the same domain can be prevented by not sharing domains. When each app is hosted on its own domain, there's no implicit cookie trust relationship to exploit.

## Antiforgery in ASP.NET Core

Warning

ASP.NET Core implements antiforgery using ASP.NET Core Data Protection. The data protection stack must be configured to work in a server farm. For more information, see Configuring data protection.

Antiforgery middleware is added to the Dependency injection container when one of the following APIs is called in `Program.cs`:

The FormTagHelper injects antiforgery tokens into HTML form elements. The following markup in a Razor file automatically generates antiforgery tokens:

```
<form method="post">
 <!-- ... -->
</form>
```
Similarly, IHtmlHelper.BeginForm generates antiforgery tokens by default if the form's method isn't GET.

The automatic generation of antiforgery tokens for HTML form elements happens when the `<form>` tag contains the `method="post"` attribute and either of the following are true:

- The action attribute is empty (`action=""`).
- The action attribute isn't supplied (`<form method="post">`).

Automatic generation of antiforgery tokens for HTML form elements can be disabled:

- Explicitly disable antiforgery tokens with the - `asp-antiforgery`attribute:- `<form method="post" asp-antiforgery="false"> <!-- ... --> </form>`
- The form element is opted-out of Tag Helpers by using the Tag Helper ! opt-out symbol: - `<!form method="post"> <!-- ... --> </!form>`
- Remove the - `FormTagHelper`from the view. The- `FormTagHelper`can be removed from a view by adding the following directive to the Razor view:- `@removeTagHelper Microsoft.AspNetCore.Mvc.TagHelpers.FormTagHelper, Microsoft.AspNetCore.Mvc.TagHelpers`

Note

Razor Pages are automatically protected from XSRF/CSRF. For more information, see XSRF/CSRF and Razor Pages.

The most common approach to defending against CSRF attacks is to use the *Synchronizer Token Pattern* (STP). STP is used when the user requests a page with form data:

- The server sends a token associated with the current user's identity to the client.
- The client sends back the token to the server for verification.
- If the server receives a token that doesn't match the authenticated user's identity, the request is rejected.

The token is unique and unpredictable. The token can also be used to ensure proper sequencing of a series of requests (for example, ensuring the request sequence of: page 1 > page 2 > page 3). All of the forms in ASP.NET Core MVC and Razor Pages templates generate antiforgery tokens. The following pair of view examples generates antiforgery tokens:

```
<form asp-action="Index" asp-controller="Home" method="post">
 <!-- ... -->
</form>
@using (Html.BeginForm("Index", "Home"))
{
 <!-- ... -->
}
```
Explicitly add an antiforgery token to a `<form>` element without using Tag Helpers with the HTML helper `@Html.AntiForgeryToken`:

```
<form asp-action="Index" asp-controller="Home" method="post">
 @Html.AntiForgeryToken()
 <!-- ... -->
</form>
```
In each of the preceding cases, ASP.NET Core adds a hidden form field similar to the following example:

```
<input name="__RequestVerificationToken" type="hidden" value="CfDJ8NrAkS ... s2-m9Yw">
```
ASP.NET Core includes three filters for working with antiforgery tokens:

## Antiforgery with `AddControllers`

Calling AddControllers does * not* enable antiforgery tokens. AddControllersWithViews must be called to have built-in antiforgery token support.

## Multiple browser tabs and the Synchronizer Token Pattern

With the Synchronizer Token Pattern, only the most recently loaded page contains a valid antiforgery token. Using multiple tabs can be problematic. For example, if a user opens multiple tabs:

- Only the most recently loaded tab contains a valid antiforgery token.
- Requests made from previously loaded tabs fail with an error: `Antiforgery token validation failed. The antiforgery cookie token and request token do not match`

Consider alternative CSRF protection patterns if this poses an issue.

## Configure antiforgery with `AntiforgeryOptions`

Customize AntiforgeryOptions in the app's `Program` file:

```
builder.Services.AddAntiforgery(options =>
{
 // Set Cookie properties using CookieBuilder properties†.
 options.FormFieldName = "AntiforgeryFieldname";
 options.HeaderName = "X-CSRF-TOKEN-HEADERNAME";
 options.SuppressXFrameOptionsHeader = false;
});
```
Set the antiforgery cookie properties using the properties of the CookieBuilder class, as shown in the following table.

| Option | Description |
|---|---|
| Cookie | Determines the settings used to create the antiforgery cookies. |
| FormFieldName | The name of the hidden form field used by the antiforgery system to render antiforgery tokens in views. |
| HeaderName | The name of the header used by the antiforgery system. If `null`, the system considers only form data. |
| SuppressXFrameOptionsHeader | Specifies whether to suppress generation of the `X-Frame-Options`header. By default, the header is generated with a value of "SAMEORIGIN". Defaults to`false`. |

Some browsers don't allow insecure endpoints to set cookies with a 'secure' flag or overwrite cookies whose 'secure' flag is set (for more information, see Deprecate modification of 'secure' cookies from non-secure origins). Since mixing secure and insecure endpoints is a common scenario in apps, ASP.NET Core relaxes the restriction on the secure policy on some cookies, such as the antiforgery cookie, by setting the cookie's SecurePolicy to `CookieSecurePolicy.None`. Even if a malicious user steals an antiforgery cookie, they also must steal the antiforgery token that's typically sent via a form field (more common) or a separate request header (less common) plus the authentication cookie. Cookies related to authentication or authorization use a stronger policy than `CookieSecurePolicy.None`.

Optionally, you can secure the antiforgery cookie in non-`Development` environments using Secure Sockets Layer (SSL), over HTTPS only, with the following AntiforgeryOptions.Cookie property setting in the app's `Program` file:

```
if (!builder.Environment.IsDevelopment())
{
 builder.Services.AddAntiforgery(o =>
 {
 o.Cookie.SecurePolicy = CookieSecurePolicy.Always;
 });
}
```
For more information, see CookieAuthenticationOptions.

## Generate antiforgery tokens with `IAntiforgery`

IAntiforgery provides the API to configure antiforgery features. `IAntiforgery` can be requested in `Program.cs` using WebApplication.Services. The following example uses middleware from the app's home page to generate an antiforgery token and send it in the response as a cookie:

```
app.UseRouting();
app.UseAuthorization();
var antiforgery = app.Services.GetRequiredService<IAntiforgery>();
app.Use((context, next) =>
{
 var requestPath = context.Request.Path.Value;
 if (string.Equals(requestPath, "/", StringComparison.OrdinalIgnoreCase)
 || string.Equals(requestPath, "/index.html", StringComparison.OrdinalIgnoreCase))
 {
 var tokenSet = antiforgery.GetAndStoreTokens(context);
 context.Response.Cookies.Append("XSRF-TOKEN", tokenSet.RequestToken!,
 new CookieOptions { HttpOnly = false });
 }
 return next(context);
});
```
The preceding example sets a cookie named `XSRF-TOKEN`. The client can read this cookie and provide its value as a header attached to AJAX requests. For example, Angular includes built-in XSRF protection that reads a cookie named `XSRF-TOKEN` by default.

### Require antiforgery validation

The ValidateAntiForgeryToken action filter can be applied to an individual action, a controller, or globally. Requests made to actions that have this filter applied are blocked unless the request includes a valid antiforgery token:

```
[HttpPost]
[ValidateAntiForgeryToken]
public IActionResult Index()
{
 // ...
 return RedirectToAction();
}
```
The `ValidateAntiForgeryToken` attribute requires a token for requests to the action methods it marks, including HTTP GET requests. If the `ValidateAntiForgeryToken` attribute is applied across the app's controllers, it can be overridden with the `IgnoreAntiforgeryToken` attribute.

### Automatically validate antiforgery tokens for unsafe HTTP methods only

Instead of broadly applying the `ValidateAntiForgeryToken` attribute and then overriding it with `IgnoreAntiforgeryToken` attributes, the AutoValidateAntiforgeryToken attribute can be used. This attribute works identically to the `ValidateAntiForgeryToken` attribute, except that it doesn't require tokens for requests made using the following HTTP methods:

- GET
- HEAD
- OPTIONS
- TRACE

We recommend use of `AutoValidateAntiforgeryToken` broadly for non-API scenarios. This attribute ensures POST actions are protected by default. The alternative is to ignore antiforgery tokens by default, unless `ValidateAntiForgeryToken` is applied to individual action methods. It's more likely in this scenario for a POST action method to be left unprotected by mistake, leaving the app vulnerable to CSRF attacks. All POSTs should send the antiforgery token.

APIs don't have an automatic mechanism for sending the non-cookie part of the token. The implementation probably depends on the client code implementation. Some examples are shown below:

Class-level example:

```
[AutoValidateAntiforgeryToken]
public class HomeController : Controller
```
Global example:

```
builder.Services.AddControllersWithViews(options =>
{
 options.Filters.Add(new AutoValidateAntiforgeryTokenAttribute());
});
```
### Override global or controller antiforgery attributes

The IgnoreAntiforgeryToken filter is used to eliminate the need for an antiforgery token for a given action (or controller). When applied, this filter overrides `ValidateAntiForgeryToken` and `AutoValidateAntiforgeryToken` filters specified at a higher level (globally or on a controller).

```
[IgnoreAntiforgeryToken]
public IActionResult IndexOverride()
{
 // ...
 return RedirectToAction();
}
```
## Refresh tokens after authentication

Tokens should be refreshed after the user is authenticated by redirecting the user to a view or Razor Pages page.

## JavaScript, AJAX, and SPAs

In traditional HTML-based apps, antiforgery tokens are passed to the server using hidden form fields. In modern JavaScript-based apps and SPAs, many requests are made programmatically. These AJAX requests may use other techniques (such as request headers or cookies) to send the token.

If cookies are used to store authentication tokens and to authenticate API requests on the server, CSRF is a potential problem. If local storage is used to store the token, CSRF vulnerability might be mitigated because values from local storage aren't sent automatically to the server with every request. Using local storage to store the antiforgery token on the client and sending the token as a request header is a recommended approach.

### JavaScript

Using JavaScript with views, the token can be created using a service from within the view. Inject the IAntiforgery service into the view and call GetAndStoreTokens:

```
@inject Microsoft.AspNetCore.Antiforgery.IAntiforgery Antiforgery
@{
 ViewData["Title"] = "JavaScript";
 var requestToken = Antiforgery.GetAndStoreTokens(Context).RequestToken;
}
<input id="RequestVerificationToken" type="hidden" value="@requestToken" />
<button id="button" class="btn btn-primary">Submit with Token</button>
<div id="result" class="mt-2"></div>
@section Scripts {
<script>
 document.addEventListener("DOMContentLoaded", () => {
 const resultElement = document.getElementById("result");
 document.getElementById("button").addEventListener("click", async () => {
 const response = await fetch("@Url.Action("FetchEndpoint")", {
 method: "POST",
 headers: {
 RequestVerificationToken:
 document.getElementById("RequestVerificationToken").value
 }
 });
 if (response.ok) {
 resultElement.innerText = await response.text();
 } else {
 resultElement.innerText = `Request Failed: ${response.status}`
 }
 });
 });
</script>
}
```
The preceding example uses JavaScript to read the hidden field value for the AJAX POST header.

This approach eliminates the need to deal directly with setting cookies from the server or reading them from the client. However, when injecting the IAntiforgery service is not possible, JavaScript can also access token in cookies, obtained from an additional request to the server (usually `same-origin`), and use the cookie's contents to create a header with the token's value.

Assuming the script sends the token in a request header called `X-XSRF-TOKEN`, configure the antiforgery service to look for the `X-XSRF-TOKEN` header:

```
builder.Services.AddAntiforgery(options => options.HeaderName = "X-XSRF-TOKEN");
```
The following example adds a protected endpoint that will write the request token to a JavaScript-readable cookie:

```
app.UseAuthorization();
app.MapGet("antiforgery/token", (IAntiforgery forgeryService, HttpContext context) =>
{
 var tokens = forgeryService.GetAndStoreTokens(context);
 context.Response.Cookies.Append("XSRF-TOKEN", tokens.RequestToken!,
 new CookieOptions { HttpOnly = false });
 return Results.Ok();
}).RequireAuthorization();
```
The following example uses JavaScript to make an AJAX request to obtain the token and make another request with the appropriate header:

```
var response = await fetch("/antiforgery/token", {
 method: "GET",
 headers: { "Authorization": authorizationToken }
});
if (response.ok) {
 // https://developer.mozilla.org/docs/web/api/document/cookie
 const xsrfToken = document.cookie
 .split("; ")
 .find(row => row.startsWith("XSRF-TOKEN="))
 .split("=")[1];
 response = await fetch("/JavaScript/FetchEndpoint", {
 method: "POST",
 headers: { "X-XSRF-TOKEN": xsrfToken, "Authorization": authorizationToken }
 });
 if (response.ok) {
 resultElement.innerText = await response.text();
 } else {
 resultElement.innerText = `Request Failed: ${response.status}`
 }
} else {
 resultElement.innerText = `Request Failed: ${response.status}`
}
```
## Windows authentication and antiforgery cookies

When using Windows Authentication, application endpoints must be protected against CSRF attacks in the same way as done for cookies. The browser implicitly sends the authentication context to the server and so endpoints need to be protected against CSRF attacks.

## Extend antiforgery

The IAntiforgeryAdditionalDataProvider type allows developers to extend the behavior of the anti-CSRF system by round-tripping additional data in each token. The GetAdditionalData method is called each time a field token is generated, and the return value is embedded within the generated token. An implementer could return a timestamp, a nonce, or any other value and then call ValidateAdditionalData to validate this data when the token is validated. The client's username is already embedded in the generated tokens, so there's no need to include this information. If a token includes supplemental data but no `IAntiForgeryAdditionalDataProvider` is configured, the supplemental data isn't validated.

## Additional resources

Cross-site request forgery (also known as XSRF or CSRF) is an attack against web-hosted apps whereby a malicious web app can influence the interaction between a client browser and a web app that trusts that browser. These attacks are possible because web browsers send some types of authentication tokens automatically with every request to a website. This form of exploit is also known as a *one-click attack* or *session riding* because the attack takes advantage of the user's previously authenticated session.

An example of a CSRF attack:

- A user signs into - `www.good-banking-site.example.com`using forms authentication. The server authenticates the user and issues a response that includes an authentication cookie. The site is vulnerable to attack because it trusts any request that it receives with a valid authentication cookie.
- The user visits a malicious site, - `www.bad-crook-site.example.com`.- The malicious site, - `www.bad-crook-site.example.com`, contains an HTML form similar to the following example:- `<h1>Congratulations! You're a Winner!</h1> <form action="https://www.good-banking-site.example.com/api/account" method="post"> <input type="hidden" name="Transaction" value="withdraw" /> <input type="hidden" name="Amount" value="1000000" /> <input type="submit" value="Click to collect your prize!" /> </form>`- Notice that the form's - `action`posts to the vulnerable site, not to the malicious site. This is the "cross-site" part of CSRF.
- The user selects the submit button. The browser makes the request and automatically includes the authentication cookie for the requested domain, - `www.good-banking-site.example.com`.
- The request runs on the - `www.good-banking-site.example.com`server with the user's authentication context and can perform any action that an authenticated user is allowed to perform.

In addition to the scenario where the user selects the button to submit the form, the malicious site could:

- Run a script that automatically submits the form.
- Send the form submission as an AJAX request.
- Hide the form using CSS.

These alternative scenarios don't require any action or input from the user other than initially visiting the malicious site.

Using HTTPS doesn't prevent a CSRF attack. The malicious site can send `https://www.good-banking-site.example.com/` a request just as easily as it can send an insecure request.

Some attacks target endpoints that respond to GET requests, in which case an image tag can be used to perform the action. This form of attack is common on forum sites that permit images but block JavaScript. Apps that change state on GET requests, where variables or resources are altered, are vulnerable to malicious attacks. **GET requests that change state are insecure. A best practice is to never change state on a GET request.**

CSRF attacks are possible against web apps that use cookies for authentication because:

- Browsers store cookies issued by a web app.
- Stored cookies include session cookies for authenticated users.
- Browsers send all of the cookies associated with a domain to the web app every request regardless of how the request to app was generated within the browser.

However, CSRF attacks aren't limited to exploiting cookies. For example, Basic and Digest authentication are also vulnerable. After a user signs in with Basic or Digest authentication, the browser automatically sends the credentials until the session ends.

In this context, *session* refers to the client-side session during which the user is authenticated. It's unrelated to server-side sessions or ASP.NET Core session middleware.

Users can protect against CSRF vulnerabilities by taking precautions:

- Sign out of web apps when finished using them.
- Clear browser cookies periodically.

However, CSRF vulnerabilities are fundamentally a problem with the web app, not the end user.

## Authentication fundamentals

Cookie-based authentication is a popular form of authentication. Token-based authentication systems are growing in popularity, especially for Single Page Applications (SPAs).

### Cookie-based authentication

When a user authenticates using their username and password, they're issued a token, containing an authentication ticket that can be used for authentication and authorization. The token is stored as a cookie that's sent with every request the client makes. Generating and validating this cookie is performed by the cookie authentication middleware. The middleware serializes a user principal into an encrypted cookie. On subsequent requests, the middleware validates the cookie, recreates the principal, and assigns the principal to the HttpContext.User property.

### Token-based authentication

When a user is authenticated, they're issued a token (not an antiforgery token). The token contains user information in the form of claims or a reference token that points the app to user state maintained in the app. When a user attempts to access a resource that requires authentication, the token is sent to the app with an extra authorization header in the form of a Bearer token. This approach makes the app stateless. In each subsequent request, the token is passed in the request for server-side validation. This token isn't *encrypted*; it's *encoded*. On the server, the token is decoded to access its information. To send the token on subsequent requests, store the token in the browser's local storage. Don't be concerned about CSRF vulnerability if the token is stored in the browser's local storage. CSRF is a concern when the token is stored in a cookie. For more information, see the GitHub issue SPA code sample adds two cookies.

### Multiple apps hosted at one domain

Shared hosting environments are vulnerable to session hijacking, login CSRF, and other attacks.

Although `example1.contoso.net` and `example2.contoso.net` are different hosts, there's an implicit trust relationship between hosts under the `*.contoso.net` domain. This implicit trust relationship allows potentially untrusted hosts to affect each other's cookies (the same-origin policies that govern AJAX requests don't necessarily apply to HTTP cookies).

Attacks that exploit trusted cookies between apps hosted on the same domain can be prevented by not sharing domains. When each app is hosted on its own domain, there's no implicit cookie trust relationship to exploit.

## ASP.NET Core antiforgery configuration

Warning

ASP.NET Core implements antiforgery using ASP.NET Core Data Protection. The data protection stack must be configured to work in a server farm. For more information, see Configuring data protection.

Antiforgery middleware is added to the Dependency injection container when one of the following APIs is called in `Startup.ConfigureServices`:

In ASP.NET Core 2.0 or later, the FormTagHelper injects antiforgery tokens into HTML form elements. The following markup in a Razor file automatically generates antiforgery tokens:

```
<form method="post">
 ...
</form>
```
Similarly, IHtmlHelper.BeginForm generates antiforgery tokens by default if the form's method isn't GET.

The automatic generation of antiforgery tokens for HTML form elements happens when the `<form>` tag contains the `method="post"` attribute and either of the following are true:

- The action attribute is empty (`action=""`).
- The action attribute isn't supplied (`<form method="post">`).

Automatic generation of antiforgery tokens for HTML form elements can be disabled:

- Explicitly disable antiforgery tokens with the - `asp-antiforgery`attribute:- `<form method="post" asp-antiforgery="false"> ... </form>`
- The form element is opted-out of Tag Helpers by using the Tag Helper ! opt-out symbol: - `<!form method="post"> ... </!form>`
- Remove the - `FormTagHelper`from the view. The- `FormTagHelper`can be removed from a view by adding the following directive to the Razor view:- `@removeTagHelper Microsoft.AspNetCore.Mvc.TagHelpers.FormTagHelper, Microsoft.AspNetCore.Mvc.TagHelpers`

Note

Razor Pages are automatically protected from XSRF/CSRF. For more information, see XSRF/CSRF and Razor Pages.

The most common approach to defending against CSRF attacks is to use the *Synchronizer Token Pattern* (STP). STP is used when the user requests a page with form data:

- The server sends a token associated with the current user's identity to the client.
- The client sends back the token to the server for verification.
- If the server receives a token that doesn't match the authenticated user's identity, the request is rejected.

The token is unique and unpredictable. The token can also be used to ensure proper sequencing of a series of requests (for example, ensuring the request sequence of: page 1 > page 2 > page 3). All of the forms in ASP.NET Core MVC and Razor Pages templates generate antiforgery tokens. The following pair of view examples generates antiforgery tokens:

```
<form asp-controller="Todo" asp-action="Create" method="post">
 ...
</form>
@using (Html.BeginForm("Create", "Todo"))
{
 ...
}
```
Explicitly add an antiforgery token to a `<form>` element without using Tag Helpers with the HTML helper `@Html.AntiForgeryToken`:

```
<form action="/" method="post">
 @Html.AntiForgeryToken()
</form>
```
In each of the preceding cases, ASP.NET Core adds a hidden form field similar to the following example:

```
<input name="__RequestVerificationToken" type="hidden" value="CfDJ8NrAkS ... s2-m9Yw">
```
ASP.NET Core includes three filters for working with antiforgery tokens:

## Antiforgery options

Customize AntiforgeryOptions in `Startup.ConfigureServices`:

```
services.AddAntiforgery(options =>
{
 options.FormFieldName = "AntiforgeryFieldname";
 options.HeaderName = "X-CSRF-TOKEN-HEADERNAME";
 options.SuppressXFrameOptionsHeader = false;
});
```
Set the antiforgery cookie properties using the properties of the CookieBuilder class, as shown in the following table.

| Option | Description |
|---|---|
| Cookie | Determines the settings used to create the antiforgery cookies. |
| FormFieldName | The name of the hidden form field used by the antiforgery system to render antiforgery tokens in views. |
| HeaderName | The name of the header used by the antiforgery system. If `null`, the system considers only form data. |
| SuppressXFrameOptionsHeader | Specifies whether to suppress generation of the `X-Frame-Options`header. By default, the header is generated with a value of "SAMEORIGIN". Defaults to`false`. |

Some browsers don't allow insecure endpoints to set cookies with a 'secure' flag or overwrite cookies whose 'secure' flag is set (for more information, see Deprecate modification of 'secure' cookies from non-secure origins). Since mixing secure and insecure endpoints is a common scenario in apps, ASP.NET Core relaxes the restriction on the secure policy on some cookies, such as the antiforgery cookie, by setting the cookie's SecurePolicy to `CookieSecurePolicy.None`. Even if a malicious user steals an antiforgery cookie, they also must steal the antiforgery token that's typically sent via a form field (more common) or a separate request header (less common) plus the authentication cookie. Cookies related to authentication or authorization use a stronger policy than `CookieSecurePolicy.None`.

Optionally, you can secure the antiforgery cookie in non-`Development` environments using Secure Sockets Layer (SSL), over HTTPS only, with the following AntiforgeryOptions.Cookie property setting in the app's `Startup` class:

```
public class Startup
{
 public Startup(IConfiguration configuration, IHostEnvironment environment)
 {
 Configuration = configuration;
 Environment = environment;
 }
 public IConfiguration Configuration { get; }
 public IHostEnvironment Environment { get; }
 public void ConfigureServices(IServiceCollection services)
 {
 // Other services are registered here
 if (!Environment.IsDevelopment())
 {
 services.AddAntiforgery(o =>
 {
 o.Cookie.SecurePolicy = CookieSecurePolicy.Always;
 });
 }
 }
 public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
 {
 // Request processing pipeline
 }
}
```
For more information, see CookieAuthenticationOptions.

## Configure antiforgery features with IAntiforgery

IAntiforgery provides the API to configure antiforgery features. `IAntiforgery` can be requested in the `Configure` method of the `Startup` class.

In the following example:

- Middleware from the app's home page is used to generate an antiforgery token and send it in the response as a cookie.
- The request token is sent as a JavaScript-readable cookie with the default Angular naming convention described in the AngularJS section.

```
public void Configure(IApplicationBuilder app, IAntiforgery antiforgery)
{
 app.Use(next => context =>
 {
 string path = context.Request.Path.Value;
 if (string.Equals(path, "/", StringComparison.OrdinalIgnoreCase) ||
 string.Equals(path, "/index.html", StringComparison.OrdinalIgnoreCase))
 {
 var tokens = antiforgery.GetAndStoreTokens(context);
 context.Response.Cookies.Append("XSRF-TOKEN", tokens.RequestToken,
 new CookieOptions() { HttpOnly = false });
 }
 return next(context);
 });
}
```
### Require antiforgery validation

ValidateAntiForgeryToken is an action filter that can be applied to an individual action, a controller, or globally. Requests made to actions that have this filter applied are blocked unless the request includes a valid antiforgery token.

```
[HttpPost]
[ValidateAntiForgeryToken]
public async Task<IActionResult> RemoveLogin(RemoveLoginViewModel account)
{
 ManageMessageId? message = ManageMessageId.Error;
 var user = await GetCurrentUserAsync();
 if (user != null)
 {
 var result =
 await _userManager.RemoveLoginAsync(
 user, account.LoginProvider, account.ProviderKey);
 if (result.Succeeded)
 {
 await _signInManager.SignInAsync(user, isPersistent: false);
 message = ManageMessageId.RemoveLoginSuccess;
 }
 }
 return RedirectToAction(nameof(ManageLogins), new { Message = message });
}
```
The `ValidateAntiForgeryToken` attribute requires a token for requests to the action methods it marks, including HTTP GET requests. If the `ValidateAntiForgeryToken` attribute is applied across the app's controllers, it can be overridden with the `IgnoreAntiforgeryToken` attribute.

Note

ASP.NET Core doesn't support adding antiforgery tokens to GET requests automatically.

### Automatically validate antiforgery tokens for unsafe HTTP methods only

ASP.NET Core apps don't generate antiforgery tokens for safe HTTP methods (GET, HEAD, OPTIONS, and TRACE). Instead of broadly applying the `ValidateAntiForgeryToken` attribute and then overriding it with `IgnoreAntiforgeryToken` attributes, the AutoValidateAntiforgeryToken attribute can be used. This attribute works identically to the `ValidateAntiForgeryToken` attribute, except that it doesn't require tokens for requests made using the following HTTP methods:

- GET
- HEAD
- OPTIONS
- TRACE

We recommend use of `AutoValidateAntiforgeryToken` broadly for non-API scenarios. This attribute ensures POST actions are protected by default. The alternative is to ignore antiforgery tokens by default, unless `ValidateAntiForgeryToken` is applied to individual action methods. It's more likely in this scenario for a POST action method to be left unprotected by mistake, leaving the app vulnerable to CSRF attacks. All POSTs should send the antiforgery token.

APIs don't have an automatic mechanism for sending the non-cookie part of the token. The implementation probably depends on the client code implementation. Some examples are shown below:

Class-level example:

```
[Authorize]
[AutoValidateAntiforgeryToken]
public class ManageController : Controller
{
```
Global example:

```
services.AddControllersWithViews(options =>
 options.Filters.Add(new AutoValidateAntiforgeryTokenAttribute()));
```
### Override global or controller antiforgery attributes

The IgnoreAntiforgeryToken filter is used to eliminate the need for an antiforgery token for a given action (or controller). When applied, this filter overrides `ValidateAntiForgeryToken` and `AutoValidateAntiforgeryToken` filters specified at a higher level (globally or on a controller).

```
[Authorize]
[AutoValidateAntiforgeryToken]
public class ManageController : Controller
{
 [HttpPost]
 [IgnoreAntiforgeryToken]
 public async Task<IActionResult> DoSomethingSafe(SomeViewModel model)
 {
 // no antiforgery token required
 }
}
```
## Refresh tokens after authentication

Tokens should be refreshed after the user is authenticated by redirecting the user to a view or Razor Pages page.

## JavaScript, AJAX, and SPAs

In traditional HTML-based apps, antiforgery tokens are passed to the server using hidden form fields. In modern JavaScript-based apps and SPAs, many requests are made programmatically. These AJAX requests may use other techniques (such as request headers or cookies) to send the token.

If cookies are used to store authentication tokens and to authenticate API requests on the server, CSRF is a potential problem. If local storage is used to store the token, CSRF vulnerability might be mitigated because values from local storage aren't sent automatically to the server with every request. Using local storage to store the antiforgery token on the client and sending the token as a request header is a recommended approach.

### JavaScript

Using JavaScript with views, the token can be created using a service from within the view. Inject the IAntiforgery service into the view and call GetAndStoreTokens:

```
@{
 ViewData["Title"] = "AJAX Demo";
}
@inject Microsoft.AspNetCore.Antiforgery.IAntiforgery Xsrf
@functions{
 public string GetAntiXsrfRequestToken()
 {
 return Xsrf.GetAndStoreTokens(Context).RequestToken;
 }
}
<input type="hidden" id="RequestVerificationToken"
 name="RequestVerificationToken" value="@GetAntiXsrfRequestToken()">
<h2>@ViewData["Title"].</h2>
<h3>@ViewData["Message"]</h3>
<div class="row">
 <p><input type="button" id="antiforgery" value="Antiforgery"></p>
 <script>
 var xhttp = new XMLHttpRequest();
 xhttp.onreadystatechange = function() {
 if (xhttp.readyState == XMLHttpRequest.DONE) {
 if (xhttp.status == 200) {
 alert(xhttp.responseText);
 } else {
 alert('There was an error processing the AJAX request.');
 }
 }
 };
 document.addEventListener('DOMContentLoaded', function() {
 document.getElementById("antiforgery").onclick = function () {
 xhttp.open('POST', '@Url.Action("Antiforgery", "Home")', true);
 xhttp.setRequestHeader("RequestVerificationToken",
 document.getElementById('RequestVerificationToken').value);
 xhttp.send();
 }
 });
 </script>
</div>
```
This approach eliminates the need to deal directly with setting cookies from the server or reading them from the client.

The preceding example uses JavaScript to read the hidden field value for the AJAX POST header.

JavaScript can also access tokens in cookies and use the cookie's contents to create a header with the token's value.

```
context.Response.Cookies.Append("CSRF-TOKEN", tokens.RequestToken,
 new Microsoft.AspNetCore.Http.CookieOptions { HttpOnly = false });
```
Assuming the script requests to send the token in a header called `X-CSRF-TOKEN`, configure the antiforgery service to look for the `X-CSRF-TOKEN` header:

```
services.AddAntiforgery(options => options.HeaderName = "X-CSRF-TOKEN");
```
The following example uses JavaScript to make an AJAX request with the appropriate header:

```
function getCookie(cname) {
 var name = cname + "=";
 var decodedCookie = decodeURIComponent(document.cookie);
 var ca = decodedCookie.split(';');
 for (var i = 0; i < ca.length; i++) {
 var c = ca[i];
 while (c.charAt(0) === ' ') {
 c = c.substring(1);
 }
 if (c.indexOf(name) === 0) {
 return c.substring(name.length, c.length);
 }
 }
 return "";
}
var csrfToken = getCookie("CSRF-TOKEN");
var xhttp = new XMLHttpRequest();
xhttp.onreadystatechange = function () {
 if (xhttp.readyState === XMLHttpRequest.DONE) {
 if (xhttp.status === 204) {
 alert('Todo item is created successfully.');
 } else {
 alert('There was an error processing the AJAX request.');
 }
 }
};
xhttp.open('POST', '/api/items', true);
xhttp.setRequestHeader("Content-type", "application/json");
xhttp.setRequestHeader("X-CSRF-TOKEN", csrfToken);
xhttp.send(JSON.stringify({ "name": "Learn C#" }));
```
### AngularJS

AngularJS uses a convention to address CSRF. If the server sends a cookie with the name `XSRF-TOKEN`, the AngularJS `$http` service adds the cookie value to a header when it sends a request to the server. This process is automatic. The client doesn't need to set the header explicitly. The header name is `X-XSRF-TOKEN`. The server should detect this header and validate its contents.

For ASP.NET Core API to work with this convention in your application startup:

- Configure your app to provide a token in a cookie called `XSRF-TOKEN`.
- Configure the antiforgery service to look for a header named `X-XSRF-TOKEN`, which is Angular's default header name for sending the XSRF token.

```
public void Configure(IApplicationBuilder app, IAntiforgery antiforgery)
{
 app.Use(next => context =>
 {
 string path = context.Request.Path.Value;
 if (
 string.Equals(path, "/", StringComparison.OrdinalIgnoreCase) ||
 string.Equals(path, "/index.html", StringComparison.OrdinalIgnoreCase))
 {
 var tokens = antiforgery.GetAndStoreTokens(context);
 context.Response.Cookies.Append("XSRF-TOKEN", tokens.RequestToken,
 new CookieOptions() { HttpOnly = false });
 }
 return next(context);
 });
}
public void ConfigureServices(IServiceCollection services)
{
 services.AddAntiforgery(options => options.HeaderName = "X-XSRF-TOKEN");
}
```
Note

When the antiforgery token is provided in both the request header and in the form payload, only the token in the header is validated.

## Windows authentication and antiforgery cookies

When using Windows Authentication, application endpoints must be protected against CSRF attacks in the same way as done for cookies. The browser implicitly sends the authentication context to the server and so endpoints need to be protected against CSRF attacks.

## Extend antiforgery

The IAntiforgeryAdditionalDataProvider type allows developers to extend the behavior of the anti-CSRF system by round-tripping additional data in each token. The GetAdditionalData method is called each time a field token is generated, and the return value is embedded within the generated token. An implementer could return a timestamp, a nonce, or any other value and then call ValidateAdditionalData to validate this data when the token is validated. The client's username is already embedded in the generated tokens, so there's no need to include this information. If a token includes supplemental data but no `IAntiForgeryAdditionalDataProvider` is configured, the supplemental data isn't validated.

## Additional resources

ASP.NET Core

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Note

This isn't the latest version of this article. For the current release, see the .NET 10 version of this article.

Warning

This version of ASP.NET Core is no longer supported. For more information, see the .NET and .NET Core Support Policy. For the current release, see the .NET 10 version of this article.

By Rick Anderson and Kirk Larkin

This article shows how **C**ross-**O**rigin **R**esource **S**haring (CORS) is enabled in an ASP.NET Core app.

Browser security prevents a web page from making requests to a different domain than the one that served the web page. This restriction is called the *same-origin policy*. The same-origin policy prevents a malicious site from reading sensitive data from another site. Sometimes, you might want to allow other sites to make cross-origin requests to your app. For more information, see the Mozilla CORS article.

Cross Origin Resource Sharing (CORS):

- Is a W3C standard that allows a server to relax the same-origin policy.
- Is **not**a security feature, CORS relaxes security. An API is not safer by allowing CORS. For more information, see How CORS works.
- Allows a server to explicitly allow some cross-origin requests while rejecting others.
- Is safer and more flexible than earlier techniques, such as JSONP.

View or download sample code (how to download)

## Same origin

Two URLs have the same origin if they have identical schemes, hosts, and ports (RFC 6454).

These two URLs have the same origin:

- `https://example.com/foo.html`
- `https://example.com/bar.html`

These URLs have different origins than the previous two URLs:

- `https://example.net`: Different domain
- `https://contoso.example.com/foo.html`: Different subdomain
- `http://example.com/foo.html`: Different scheme
- `https://example.com:9000/foo.html`: Different port

## Enable CORS

There are three ways to enable CORS:

- In middleware using a named policy or default policy.
- Using endpoint routing.
- With the [EnableCors] attribute.

Using the [EnableCors] attribute with a named policy provides the finest control in limiting endpoints that support CORS.

Warning

UseCors must be called in the correct order. For more information, see Middleware order. For example, `UseCors` must be called before UseResponseCaching when using `UseResponseCaching`.

Each approach is detailed in the following sections.

## CORS with named policy and middleware

CORS middleware handles cross-origin requests. The following code applies a CORS policy to all the app's endpoints with the specified origins:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
});
// services.AddResponseCaching();
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors(MyAllowSpecificOrigins);
app.UseAuthorization();
app.MapControllers();
app.Run();
```
The preceding code:

- Sets the policy name to `_myAllowSpecificOrigins`. The policy name is arbitrary.
- Calls the UseCors extension method and specifies the `_myAllowSpecificOrigins`CORS policy.`UseCors`adds the CORS middleware. The call to`UseCors`must be placed after`UseRouting`, but before`UseAuthorization`. For more information, see Middleware order.
- Calls AddCors with a lambda expression. The lambda takes a CorsPolicyBuilder object. Configuration options, such as `WithOrigins`, are described later in this article.
- Enables the `_myAllowSpecificOrigins`CORS policy for all controller endpoints. See endpoint routing to apply a CORS policy to specific endpoints.
- When using response caching middleware, call UseCors before UseResponseCaching.

With endpoint routing, the CORS middleware **must** be configured to execute between the calls to `UseRouting` and `UseEndpoints`.

The AddCors method call adds CORS services to the app's service container:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
});
// services.AddResponseCaching();
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors(MyAllowSpecificOrigins);
app.UseAuthorization();
app.MapControllers();
app.Run();
```
For more information, see CORS policy options in this document.

The CorsPolicyBuilder methods can be chained, as shown in the following code:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com")
 .AllowAnyHeader()
 .AllowAnyMethod();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors(MyAllowSpecificOrigins);
app.UseAuthorization();
app.MapControllers();
app.Run();
```
Note: The specified URL must **not** contain a trailing slash (`/`). If the URL terminates with `/`, the comparison returns `false` and no header is returned.

## UseCors and UseStaticFiles order

Typically, `UseStaticFiles` is called before `UseCors`. Apps that use JavaScript to retrieve static files cross site must call `UseCors` before `UseStaticFiles`.

### CORS with default policy and middleware

The following highlighted code enables the default CORS policy:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddDefaultPolicy(
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
});
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.MapControllers();
app.Run();
```
The preceding code applies the default CORS policy to all controller endpoints.

## Enable Cors with endpoint routing

With endpoint routing, CORS can be enabled on a per-endpoint basis using the RequireCors set of extension methods:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
});
builder.Services.AddControllers();
builder.Services.AddRazorPages();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.UseEndpoints(endpoints =>
{
 endpoints.MapGet("/echo",
 context => context.Response.WriteAsync("echo"))
 .RequireCors(MyAllowSpecificOrigins);
 endpoints.MapControllers()
 .RequireCors(MyAllowSpecificOrigins);
 endpoints.MapGet("/echo2",
 context => context.Response.WriteAsync("echo2"));
 endpoints.MapRazorPages();
});
app.Run();
```
In the preceding code:

- `app.UseCors`enables the CORS middleware. Because a default policy hasn't been configured,- `app.UseCors()`alone doesn't enable CORS.
- The `/echo`and controller endpoints allow cross-origin requests using the specified policy.
- The `/echo2`and Razor Pages endpoints do**not**allow cross-origin requests because no default policy was specified.

The [DisableCors] attribute does **not** disable CORS that has been enabled by endpoint routing with `RequireCors`.

See Test CORS with [EnableCors] attribute and RequireCors method for instructions on testing code similar to the preceding.

## Enable CORS with attributes

Enabling CORS with the [EnableCors] attribute and applying a named policy to only those endpoints that require CORS provides the finest control.

The [EnableCors] attribute provides an alternative to applying CORS globally. The `[EnableCors]` attribute enables CORS for selected endpoints, rather than all endpoints:

- `[EnableCors]`specifies the default policy.
- `[EnableCors("{Policy String}")]`specifies a named policy.

The `[EnableCors]` attribute can be applied to:

- Razor Page `PageModel`
- Controller
- Controller action method

Different policies can be applied to controllers, page models, or action methods with the `[EnableCors]` attribute. When the `[EnableCors]` attribute is applied to a controller, page model, or action method, and CORS is enabled in middleware, **both** policies are applied. **We recommend against combining policies. Use the** `[EnableCors]` **attribute or middleware, not both in the same app.**

The following code applies a different policy to each method:

```
[Route("api/[controller]")]
[ApiController]
public class WidgetController : ControllerBase
{
 // GET api/values
 [EnableCors("AnotherPolicy")]
 [HttpGet]
 public ActionResult<IEnumerable<string>> Get()
 {
 return new string[] { "green widget", "red widget" };
 }
 // GET api/values/5
 [EnableCors("Policy1")]
 [HttpGet("{id}")]
 public ActionResult<string> Get(int id)
 {
 return id switch
 {
 1 => "green widget",
 2 => "red widget",
 _ => NotFound(),
 };
 }
}
```
The following code creates two CORS policies:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("Policy1",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
 options.AddPolicy("AnotherPolicy",
 policy =>
 {
 policy.WithOrigins("http://www.contoso.com")
 .AllowAnyHeader()
 .AllowAnyMethod();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.MapControllers();
app.Run();
```
For the finest control of limiting CORS requests:

- Use `[EnableCors("MyPolicy")]`with a named policy.
- Don't define a default policy.
- Don't use endpoint routing.

The code in the next section meets the preceding list.

### Disable CORS

The [DisableCors] attribute does **not** disable CORS that has been enabled by endpoint routing.

The following code defines the CORS policy `"MyPolicy"`:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: "MyPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com")
 .WithMethods("PUT", "DELETE", "GET");
 });
});
builder.Services.AddControllers();
builder.Services.AddRazorPages();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.UseEndpoints(endpoints => {
 endpoints.MapControllers();
 endpoints.MapRazorPages();
});
app.Run();
```
The following code disables CORS for the `GetValues2` action:

```
[EnableCors("MyPolicy")]
[Route("api/[controller]")]
[ApiController]
public class ValuesController : ControllerBase
{
 // GET api/values
 [HttpGet]
 public IActionResult Get() =>
 ControllerContext.MyDisplayRouteInfo();
 // GET api/values/5
 [HttpGet("{id}")]
 public IActionResult Get(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // PUT api/values/5
 [HttpPut("{id}")]
 public IActionResult Put(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // GET: api/values/GetValues2
 [DisableCors]
 [HttpGet("{action}")]
 public IActionResult GetValues2() =>
 ControllerContext.MyDisplayRouteInfo();
}
```
The preceding code:

- Doesn't enable CORS with endpoint routing.
- Doesn't define a default CORS policy.
- Uses [EnableCors("MyPolicy")] to enable the `"MyPolicy"`CORS policy for the controller.
- Disables CORS for the `GetValues2`method.

See Test CORS for instructions on testing the preceding code.

## CORS policy options

This section describes the various options that can be set in a CORS policy:

- Set the allowed origins
- Set the allowed HTTP methods
- Set the allowed request headers
- Set the exposed response headers
- Credentials in cross-origin requests
- Set the preflight expiration time

AddPolicy is called in `Program.cs`. For some options, it may be helpful to read the How CORS works section first.

## Set the allowed origins

AllowAnyOrigin: Allows CORS requests from all origins with any scheme (`http` or `https`). `AllowAnyOrigin` is insecure because *any website* can make cross-origin requests to the app.

Note

Specifying `AllowAnyOrigin` and `AllowCredentials` is an insecure configuration and can result in cross-site request forgery. The CORS service returns an invalid CORS response when an app is configured with both methods.

`AllowAnyOrigin` affects preflight requests and the `Access-Control-Allow-Origin` header. For more information, see the Preflight requests section.

SetIsOriginAllowedToAllowWildcardSubdomains: Sets the IsOriginAllowed property of the policy to be a function that allows origins to match a configured wildcard domain when evaluating if the origin is allowed.

```
var MyAllowSpecificOrigins = "_MyAllowSubdomainPolicy";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .SetIsOriginAllowedToAllowWildcardSubdomains();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
In the preceding code, `SetIsOriginAllowedToAllowWildcardSubdomains` is called with the wildcard origin `"https://*.example.com"`. This configuration allows CORS requests from any subdomain of `example.com`, such as `https://subdomain.example.com` or `https://api.example.com`. The `*` wildcard character must be included in the origin to enable wildcard subdomain matching.

### Set the allowed HTTP methods

- Allows any HTTP method:
- Affects preflight requests and the `Access-Control-Allow-Methods`header. For more information, see the Preflight requests section.

### Set the allowed request headers

To allow specific headers to be sent in a CORS request, called *author request headers*, call WithHeaders and specify the allowed headers:

```
using Microsoft.Net.Http.Headers;
var MyAllowSpecificOrigins = "_MyAllowSubdomainPolicy";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .WithHeaders(HeaderNames.ContentType, "x-custom-header");
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
To allow all author request headers, call AllowAnyHeader:

```
var MyAllowSpecificOrigins = "_MyAllowSubdomainPolicy";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .AllowAnyHeader();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
`AllowAnyHeader` affects preflight requests and the Access-Control-Request-Headers header. For more information, see the Preflight requests section.

A CORS middleware policy match to specific headers specified by `WithHeaders` is only possible when the headers sent in `Access-Control-Request-Headers` exactly match the headers stated in `WithHeaders`.

For instance, consider an app configured as follows:

```
app.UseCors(policy => policy.WithHeaders(HeaderNames.CacheControl));
```
CORS middleware declines a preflight request with the following request header because `Content-Language` (HeaderNames.ContentLanguage) isn't listed in `WithHeaders`:

```
Access-Control-Request-Headers: Cache-Control, Content-Language
```
The app returns a `204 No Content` response but doesn't send the CORS headers back. Therefore, the browser doesn't attempt the cross-origin request.

### Set the exposed response headers

By default, the browser doesn't expose all of the response headers to the app. For more information, see W3C Cross-Origin Resource Sharing (Terminology): Simple Response Header.

The response headers that are available by default are:

- `Cache-Control`
- `Content-Language`
- `Content-Type`
- `Expires`
- `Last-Modified`
- `Pragma`

The CORS specification calls these headers *simple response headers*. To make other headers available to the app, call WithExposedHeaders:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MyExposeResponseHeadersPolicy",
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .WithExposedHeaders("x-custom-header");
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
### Credentials in cross-origin requests

Credentials require special handling in a CORS request. By default, the browser doesn't send credentials with a cross-origin request. Credentials include cookies and HTTP authentication schemes. To send credentials with a cross-origin request, the client must set `XMLHttpRequest.withCredentials` to `true`.

Using `XMLHttpRequest` directly:

```
var xhr = new XMLHttpRequest();
xhr.open('get', 'https://www.example.com/api/test');
xhr.withCredentials = true;
```
Using jQuery:

```
$.ajax({
 type: 'get',
 url: 'https://www.example.com/api/test',
 xhrFields: {
 withCredentials: true
 }
});
```
Using the Fetch API:

```
fetch('https://www.example.com/api/test', {
 credentials: 'include'
});
```
The server must allow the credentials. To allow cross-origin credentials, call AllowCredentials:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MyMyAllowCredentialsPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .AllowCredentials();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
The HTTP response includes an `Access-Control-Allow-Credentials` header, which tells the browser that the server allows credentials for a cross-origin request.

If the browser sends credentials but the response doesn't include a valid `Access-Control-Allow-Credentials` header, the browser doesn't expose the response to the app, and the cross-origin request fails.

Allowing cross-origin credentials is a security risk. A website at another domain can send a signed-in user's credentials to the app on the user's behalf without the user's knowledge.

The CORS specification also states that setting origins to `"*"` (all origins) is invalid if the `Access-Control-Allow-Credentials` header is present.

## Preflight requests

For some CORS requests, the browser sends an additional OPTIONS request before making the actual request. This request is called a preflight request. The browser can skip the preflight request if **all** the following conditions are true:

- The request method is GET, HEAD, or POST.
- The app doesn't set request headers other than `Accept`,`Accept-Language`,`Content-Language`,`Content-Type`, or`Last-Event-ID`.
- The `Content-Type`header, if set, has one of the following values:- `application/x-www-form-urlencoded`
- `multipart/form-data`
- `text/plain`

The rule on request headers set for the client request applies to headers that the app sets by calling `setRequestHeader` on the `XMLHttpRequest` object. The CORS specification calls these headers *author request headers*. The rule doesn't apply to headers the browser can set, such as `User-Agent`, `Host`, or `Content-Length`.

Note

This article contains URLs created by deploying the sample code to two Azure web sites, `https://cors3.azurewebsites.net` and `https://cors.azurewebsites.net`.

The following is an example response similar to the preflight request made from the **[Put test]** button in the Test CORS section of this document.

```
General:
Request URL: https://cors3.azurewebsites.net/api/values/5
Request Method: OPTIONS
Status Code: 204 No Content
Response Headers:
Access-Control-Allow-Methods: PUT,DELETE,GET
Access-Control-Allow-Origin: https://cors1.azurewebsites.net
Server: Microsoft-IIS/10.0
Set-Cookie: ARRAffinity=8f8...8;Path=/;HttpOnly;Domain=cors1.azurewebsites.net
Vary: Origin
Request Headers:
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9
Access-Control-Request-Method: PUT
Connection: keep-alive
Host: cors3.azurewebsites.net
Origin: https://cors1.azurewebsites.net
Referer: https://cors1.azurewebsites.net/
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
User-Agent: Mozilla/5.0
```
The preflight request uses the HTTP OPTIONS method. It may include the following headers:

- Access-Control-Request-Method: The HTTP method that will be used for the actual request.
- Access-Control-Request-Headers: A list of request headers that the app sets on the actual request. As stated earlier, this doesn't include headers that the browser sets, such as `User-Agent`.

If the preflight request is denied, the app returns a `204 No Content` response but doesn't set the CORS headers. Therefore, the browser doesn't attempt the cross-origin request. For an example of a denied preflight request, see the Test CORS section of this document.

Using the F12 tools, the console app shows an error similar to one of the following, depending on the browser:

- Firefox: Cross-Origin Request Blocked: The Same Origin Policy disallows reading the remote resource at `https://cors1.azurewebsites.net/api/TodoItems1/MyDelete2/5`. (Reason: CORS request did not succeed). Learn More
- Chromium based: Access to fetch at 'https://cors1.azurewebsites.net/api/TodoItems1/MyDelete2/5' from origin 'https://cors3.azurewebsites.net' has been blocked by CORS policy: Response to preflight request doesn't pass access control check: No 'Access-Control-Allow-Origin' header is present on the requested resource. If an opaque response serves your needs, set the request's mode to 'no-cors' to fetch the resource with CORS disabled.

To allow specific headers, call WithHeaders:

```
using Microsoft.Net.Http.Headers;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MyAllowHeadersPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .WithHeaders(HeaderNames.ContentType, "x-custom-header");
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
To allow all author request headers, call AllowAnyHeader:

```
using Microsoft.Net.Http.Headers;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MyAllowAllHeadersPolicy",
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .AllowAnyHeader();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
Browsers aren't consistent in how they set `Access-Control-Request-Headers`. If either:

- Headers are set to anything other than `"*"`
- AllowAnyHeader is called:
Include at least `Accept`,`Content-Type`, and`Origin`, plus any custom headers that you want to support.

### Automatic preflight request code

When the CORS policy is applied either:

- Globally by calling `app.UseCors`in`Program.cs`.
- Using the `[EnableCors]`attribute.

ASP.NET Core responds to the preflight OPTIONS request.

The Test CORS section of this document demonstrates this behavior.

### [HttpOptions] attribute for preflight requests

When CORS is enabled with the appropriate policy, ASP.NET Core generally responds to CORS preflight requests automatically.

The following code uses the [HttpOptions] attribute to create endpoints for OPTIONS requests:

```
[Route("api/[controller]")]
[ApiController]
public class TodoItems2Controller : ControllerBase
{
 // OPTIONS: api/TodoItems2/5
 [HttpOptions("{id}")]
 public IActionResult PreflightRoute(int id)
 {
 return NoContent();
 }
 // OPTIONS: api/TodoItems2
 [HttpOptions]
 public IActionResult PreflightRoute()
 {
 return NoContent();
 }
 [HttpPut("{id}")]
 public IActionResult PutTodoItem(int id)
 {
 if (id < 1)
 {
 return BadRequest();
 }
 return ControllerContext.MyDisplayRouteInfo(id);
 }
```
See Test CORS with [EnableCors] attribute and RequireCors method for instructions on testing the preceding code.

### Set the preflight expiration time

The `Access-Control-Max-Age` header specifies how long the response to the preflight request can be cached. To set this header, call SetPreflightMaxAge:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MySetPreflightExpirationPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .SetPreflightMaxAge(TimeSpan.FromSeconds(2520));
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
## Enable CORS on an endpoint

## How CORS works

This section describes what happens in a CORS request at the level of the HTTP messages.

- CORS is **not**a security feature. CORS is a W3C standard that allows a server to relax the same-origin policy.- For example, a malicious actor could use Cross-Site Scripting (XSS) against your site and execute a cross-site request to their CORS enabled site to steal information.

- An API isn't safer by allowing CORS.
- It's up to the client (browser) to enforce CORS. The server executes the request and returns the response, it's the client that returns an error and blocks the response. For example, any of the following tools will display the server response:
- Fiddler
- .NET HttpClient
- A web browser by entering the URL in the address bar.

- It's up to the client (browser) to enforce CORS. The server executes the request and returns the response, it's the client that returns an error and blocks the response. For example, any of the following tools will display the server response:
- It's a way for a server to allow browsers to execute a cross-origin XHR or Fetch API request that otherwise would be forbidden.
- Browsers without CORS can't do cross-origin requests. Before CORS, JSONP was used to circumvent this restriction. JSONP doesn't use XHR, it uses the `<script>`tag to receive the response. Scripts are allowed to be loaded cross-origin.

- Browsers without CORS can't do cross-origin requests. Before CORS, JSONP was used to circumvent this restriction. JSONP doesn't use XHR, it uses the

The CORS specification introduced several new HTTP headers that enable cross-origin requests. If a browser supports CORS, it sets these headers automatically for cross-origin requests. Custom JavaScript code isn't required to enable CORS.

The following is an example of a cross-origin request from the **Values** test button to `https://cors1.azurewebsites.net/api/values`. The `Origin` header:

- Provides the domain of the site that's making the request.
- Is required and must be different from the host.

**General headers**

```
Request URL: https://cors1.azurewebsites.net/api/values
Request Method: GET
Status Code: 200 OK
```
**Response headers**

```
Content-Encoding: gzip
Content-Type: text/plain; charset=utf-8
Server: Microsoft-IIS/10.0
Set-Cookie: ARRAffinity=8f...;Path=/;HttpOnly;Domain=cors1.azurewebsites.net
Transfer-Encoding: chunked
Vary: Accept-Encoding
X-Powered-By: ASP.NET
```
**Request headers**

```
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9
Connection: keep-alive
Host: cors1.azurewebsites.net
Origin: https://cors3.azurewebsites.net
Referer: https://cors3.azurewebsites.net/
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
User-Agent: Mozilla/5.0 ...
```
In `OPTIONS` requests, the server sets the **Response headers** `Access-Control-Allow-Origin: {allowed origin}` header in the response. For example, in the sample code, the ` Delete [EnableCors]` button `OPTIONS` request contains the following headers:

**General headers**

```
Request URL: https://cors3.azurewebsites.net/api/TodoItems2/MyDelete2/5
Request Method: OPTIONS
Status Code: 204 No Content
```
**Response headers**

```
Access-Control-Allow-Headers: Content-Type,x-custom-header
Access-Control-Allow-Methods: PUT,DELETE,GET,OPTIONS
Access-Control-Allow-Origin: https://cors1.azurewebsites.net
Server: Microsoft-IIS/10.0
Set-Cookie: ARRAffinity=8f...;Path=/;HttpOnly;Domain=cors3.azurewebsites.net
Vary: Origin
X-Powered-By: ASP.NET
```
**Request headers**

```
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9
Access-Control-Request-Headers: content-type
Access-Control-Request-Method: DELETE
Connection: keep-alive
Host: cors3.azurewebsites.net
Origin: https://cors1.azurewebsites.net
Referer: https://cors1.azurewebsites.net/test?number=2
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
User-Agent: Mozilla/5.0
```
In the preceding **Response headers**, the server sets the Access-Control-Allow-Origin header in the response. The `https://cors1.azurewebsites.net` value of this header matches the `Origin` header from the request.

If AllowAnyOrigin is called, the `Access-Control-Allow-Origin: *`, the wildcard value, is returned. `AllowAnyOrigin` allows any origin.

If the response doesn't include the `Access-Control-Allow-Origin` header, the cross-origin request fails. Specifically, the browser disallows the request. Even if the server returns a successful response, the browser doesn't make the response available to the client app.

### HTTP redirection to HTTPS causes ERR_INVALID_REDIRECT on the CORS preflight request

Requests to an endpoint using HTTP that are redirected to HTTPS by UseHttpsRedirection fail with `ERR_INVALID_REDIRECT on the CORS preflight request`.

API projects can reject HTTP requests rather than use `UseHttpsRedirection` to redirect requests to HTTPS.

## CORS in IIS

When deploying to IIS, CORS has to run before Windows Authentication if the server isn't configured to allow anonymous access. To support this scenario, the IIS CORS module needs to be installed and configured for the app.

## Test CORS

The sample download has code to test CORS. See how to download. The sample is an API project with Razor Pages added:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: "MyPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com",
 "https://cors1.azurewebsites.net",
 "https://cors3.azurewebsites.net",
 "https://localhost:44398",
 "https://localhost:5001")
 .WithMethods("PUT", "DELETE", "GET");
 });
});
builder.Services.AddControllers();
builder.Services.AddRazorPages();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.MapControllers();
app.MapRazorPages();
app.Run();
```
Warning

`WithOrigins("https://localhost:<port>");` should only be used for testing a sample app similar to the download sample code.

Note

If you're using `launchSettings.json` in Visual Studio or configuring C# debug settings in VS Code, and using IIS Express to debug locally, ensure that you've configured IIS Express for `"anonymousAuthentication": true`. When `"anonymousAuthentication"` is `false`, the ASP.NET Core web environment host will not see any preflight requests. In particular, if you're using NTLM authentication (`"windowsAuthentication": true`), the first step of the NTLM challenge-response is to send the web browser a 401 challenge, which can make it challenging to verify your preflight route is configured correctly.

The following `ValuesController` provides the endpoints for testing:

```
[EnableCors("MyPolicy")]
[Route("api/[controller]")]
[ApiController]
public class ValuesController : ControllerBase
{
 // GET api/values
 [HttpGet]
 public IActionResult Get() =>
 ControllerContext.MyDisplayRouteInfo();
 // GET api/values/5
 [HttpGet("{id}")]
 public IActionResult Get(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // PUT api/values/5
 [HttpPut("{id}")]
 public IActionResult Put(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // GET: api/values/GetValues2
 [DisableCors]
 [HttpGet("{action}")]
 public IActionResult GetValues2() =>
 ControllerContext.MyDisplayRouteInfo();
}
```
MyDisplayRouteInfo is provided by the Rick.Docs.Samples.RouteInfo NuGet package and displays route information.

Test the preceding sample code by using one of the following approaches:

- Run the sample with `dotnet run`using the default URL of`https://localhost:5001`.
- Run the sample from Visual Studio with the port set to 44398 for a URL of `https://localhost:44398`.

Using a browser with the F12 tools:

- Select the - **Values**button and review the headers in the- **Network**tab.
- Select the - **PUT test**button. See Display OPTIONS requests for instructions on displaying the OPTIONS request. The- **PUT test**creates two requests, an OPTIONS preflight request and the PUT request.
- Select the - `GetValues2 [DisableCors]`- **Console**tab to see the CORS error. Depending on the browser, an error similar to the following is displayed:- Access to fetch at - `'https://cors1.azurewebsites.net/api/values/GetValues2'`from origin- `'https://cors3.azurewebsites.net'`has been blocked by CORS policy: No 'Access-Control-Allow-Origin' header is present on the requested resource. If an opaque response serves your needs, set the request's mode to 'no-cors' to fetch the resource with CORS disabled.

CORS-enabled endpoints can be tested with a tool, such as curl or Fiddler. When using a tool, the origin of the request specified by the `Origin` header must differ from the host receiving the request. If the request isn't *cross-origin* based on the value of the `Origin` header:

- There's no need for CORS middleware to process the request.
- CORS headers aren't returned in the response.

The following command uses `curl` to issue an OPTIONS request with information:

```
curl -X OPTIONS https://cors3.azurewebsites.net/api/TodoItems2/5 -i
```
### Test CORS with [EnableCors] attribute and RequireCors method

Consider the following code which uses endpoint routing to enable CORS on a per-endpoint basis using `RequireCors`:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: "MyPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com",
 "https://cors1.azurewebsites.net",
 "https://cors3.azurewebsites.net",
 "https://localhost:44398",
 "https://localhost:5001")
 .WithMethods("PUT", "DELETE", "GET");
 });
});
builder.Services.AddControllers();
builder.Services.AddRazorPages();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.UseEndpoints(endpoints =>
{
 endpoints.MapGet("/echo",
 context => context.Response.WriteAsync("echo"))
 .RequireCors("MyPolicy");
 endpoints.MapControllers();
 endpoints.MapRazorPages();
});
app.Run();
```
Notice that only the `/echo` endpoint is using the `RequireCors` to allow cross-origin requests using the specified policy. The controllers below enable CORS using [EnableCors] attribute.

The following `TodoItems1Controller` provides endpoints for testing:

```
[Route("api/[controller]")]
[ApiController]
public class TodoItems1Controller : ControllerBase
{
 // PUT: api/TodoItems1/5
 [HttpPut("{id}")]
 public IActionResult PutTodoItem(int id) {
 if (id < 1) {
 return Content($"ID = {id}");
 }
 return ControllerContext.MyDisplayRouteInfo(id);
 }
 // Delete: api/TodoItems1/5
 [HttpDelete("{id}")]
 public IActionResult MyDelete(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // GET: api/TodoItems1
 [HttpGet]
 public IActionResult GetTodoItems() =>
 ControllerContext.MyDisplayRouteInfo();
 [EnableCors("MyPolicy")]
 [HttpGet("{action}")]
 public IActionResult GetTodoItems2() =>
 ControllerContext.MyDisplayRouteInfo();
 // Delete: api/TodoItems1/MyDelete2/5
 [EnableCors("MyPolicy")]
 [HttpDelete("{action}/{id}")]
 public IActionResult MyDelete2(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
}
```
The **Delete [EnableCors]** and **GET [EnableCors]** buttons succeed, because the endpoints have `[EnableCors]` and respond to preflight requests. The other endpoints fails. The **GET** button fails, because the JavaScript sends:

```
 headers: {
 "Content-Type": "x-custom-header"
 },
```
The following `TodoItems2Controller` provides similar endpoints, but includes explicit code to respond to OPTIONS requests:

```
[Route("api/[controller]")]
[ApiController]
public class TodoItems2Controller : ControllerBase
{
 // OPTIONS: api/TodoItems2/5
 [HttpOptions("{id}")]
 public IActionResult PreflightRoute(int id)
 {
 return NoContent();
 }
 // OPTIONS: api/TodoItems2
 [HttpOptions]
 public IActionResult PreflightRoute()
 {
 return NoContent();
 }
 [HttpPut("{id}")]
 public IActionResult PutTodoItem(int id)
 {
 if (id < 1)
 {
 return BadRequest();
 }
 return ControllerContext.MyDisplayRouteInfo(id);
 }
 // [EnableCors] // Not needed as OPTIONS path provided.
 [HttpDelete("{id}")]
 public IActionResult MyDelete(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // [EnableCors] // Warning ASP0023 Route '{id}' conflicts with another action route.
 // An HTTP request that matches multiple routes results in an ambiguous
 // match error.
 [EnableCors("MyPolicy")] // Required for this path.
 [HttpGet]
 public IActionResult GetTodoItems() =>
 ControllerContext.MyDisplayRouteInfo();
 [HttpGet("{action}")]
 public IActionResult GetTodoItems2() =>
 ControllerContext.MyDisplayRouteInfo();
 [EnableCors("MyPolicy")] // Required for this path.
 [HttpDelete("{action}/{id}")]
 public IActionResult MyDelete2(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
}
```
The preceding code can be tested by deploying the sample to Azure. In the **Controller** drop down list, select **Preflight** and then **Set Controller**. All the CORS calls to the `TodoItems2Controller` endpoints succeed.

## Additional resources

By Rick Anderson and Kirk Larkin

This article shows how to enable CORS in an ASP.NET Core app.

Browser security prevents a web page from making requests to a different domain than the one that served the web page. This restriction is called the *same-origin policy*. The same-origin policy prevents a malicious site from reading sensitive data from another site. Sometimes, you might want to allow other sites to make cross-origin requests to your app. For more information, see the Mozilla CORS article.

Cross Origin Resource Sharing (CORS):

- Is a W3C standard that allows a server to relax the same-origin policy.
- Is **not**a security feature, CORS relaxes security. An API is not safer by allowing CORS. For more information, see How CORS works.
- Allows a server to explicitly allow some cross-origin requests while rejecting others.
- Is safer and more flexible than earlier techniques, such as JSONP.

View or download sample code (how to download)

## Same origin

Two URLs have the same origin if they have identical schemes, hosts, and ports (RFC 6454).

These two URLs have the same origin:

- `https://example.com/foo.html`
- `https://example.com/bar.html`

These URLs have different origins than the previous two URLs:

- `https://example.net`: Different domain
- `https://www.example.com/foo.html`: Different subdomain
- `http://example.com/foo.html`: Different scheme
- `https://example.com:9000/foo.html`: Different port

## Enable CORS

There are three ways to enable CORS:

- In middleware using a named policy or default policy.
- Using endpoint routing.
- With the [EnableCors] attribute.

Using the [EnableCors] attribute with a named policy provides the finest control in limiting endpoints that support CORS.

Warning

UseCors must be called in the correct order. For more information, see Middleware order. For example, `UseCors` must be called before UseResponseCaching when using `UseResponseCaching`.

Each approach is detailed in the following sections.

## CORS with named policy and middleware

CORS middleware handles cross-origin requests. The following code applies a CORS policy to all the app's endpoints with the specified origins:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
});
// services.AddResponseCaching();
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors(MyAllowSpecificOrigins);
app.UseAuthorization();
app.MapControllers();
app.Run();
```
The preceding code:

- Sets the policy name to `_myAllowSpecificOrigins`. The policy name is arbitrary.
- Calls the UseCors extension method and specifies the `_myAllowSpecificOrigins`CORS policy.`UseCors`adds the CORS middleware. The call to`UseCors`must be placed after`UseRouting`, but before`UseAuthorization`. For more information, see Middleware order.
- Calls AddCors with a lambda expression. The lambda takes a CorsPolicyBuilder object. Configuration options, such as `WithOrigins`, are described later in this article.
- Enables the `_myAllowSpecificOrigins`CORS policy for all controller endpoints. See endpoint routing to apply a CORS policy to specific endpoints.
- When using response caching middleware, call UseCors before UseResponseCaching.

With endpoint routing, the CORS middleware **must** be configured to execute between the calls to `UseRouting` and `UseEndpoints`.

The AddCors method call adds CORS services to the app's service container:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
});
// services.AddResponseCaching();
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors(MyAllowSpecificOrigins);
app.UseAuthorization();
app.MapControllers();
app.Run();
```
For more information, see CORS policy options in this document.

The CorsPolicyBuilder methods can be chained, as shown in the following code:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com")
 .AllowAnyHeader()
 .AllowAnyMethod();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors(MyAllowSpecificOrigins);
app.UseAuthorization();
app.MapControllers();
app.Run();
```
Note: The specified URL must **not** contain a trailing slash (`/`). If the URL terminates with `/`, the comparison returns `false` and no header is returned.

Warning

`UseCors` must be placed after `UseRouting` and before `UseAuthorization`. This is to ensure that CORS headers are included in the response for both authorized and unauthorized calls.

## UseCors and UseStaticFiles order

Typically, `UseStaticFiles` is called before `UseCors`. Apps that use JavaScript to retrieve static files cross site must call `UseCors` before `UseStaticFiles`.

### CORS with default policy and middleware

The following highlighted code enables the default CORS policy:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddDefaultPolicy(
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
});
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.MapControllers();
app.Run();
```
The preceding code applies the default CORS policy to all controller endpoints.

## Enable Cors with endpoint routing

With endpoint routing, CORS can be enabled on a per-endpoint basis using the RequireCors set of extension methods:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
});
builder.Services.AddControllers();
builder.Services.AddRazorPages();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.UseEndpoints(endpoints =>
{
 endpoints.MapGet("/echo",
 context => context.Response.WriteAsync("echo"))
 .RequireCors(MyAllowSpecificOrigins);
 endpoints.MapControllers()
 .RequireCors(MyAllowSpecificOrigins);
 endpoints.MapGet("/echo2",
 context => context.Response.WriteAsync("echo2"));
 endpoints.MapRazorPages();
});
app.Run();
```
In the preceding code:

- `app.UseCors`enables the CORS middleware. Because a default policy hasn't been configured,- `app.UseCors()`alone doesn't enable CORS.
- The `/echo`and controller endpoints allow cross-origin requests using the specified policy.
- The `/echo2`and Razor Pages endpoints do**not**allow cross-origin requests because no default policy was specified.

The [DisableCors] attribute does **not** disable CORS that has been enabled by endpoint routing with `RequireCors`.

In .NET 7, the `[EnableCors]` attribute must pass a parameter or an ASP0023 Warning is generated from a ambiguous match on the route. .NET 8 or later doesn't generate the `ASP0023` warning.

```
[Route("api/[controller]")]
[ApiController]
public class TodoItems2Controller : ControllerBase
{
 // OPTIONS: api/TodoItems2/5
 [HttpOptions("{id}")]
 public IActionResult PreflightRoute(int id)
 {
 return NoContent();
 }
 // OPTIONS: api/TodoItems2
 [HttpOptions]
 public IActionResult PreflightRoute()
 {
 return NoContent();
 }
 [HttpPut("{id}")]
 public IActionResult PutTodoItem(int id)
 {
 if (id < 1)
 {
 return BadRequest();
 }
 return ControllerContext.MyDisplayRouteInfo(id);
 }
 // [EnableCors] // Not needed as OPTIONS path provided.
 [HttpDelete("{id}")]
 public IActionResult MyDelete(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // [EnableCors] // Warning ASP0023 Route '{id}' conflicts with another action route.
 // An HTTP request that matches multiple routes results in an ambiguous
 // match error.
 [EnableCors("MyPolicy")] // Required for this path.
 [HttpGet]
 public IActionResult GetTodoItems() =>
 ControllerContext.MyDisplayRouteInfo();
 [HttpGet("{action}")]
 public IActionResult GetTodoItems2() =>
 ControllerContext.MyDisplayRouteInfo();
 [EnableCors("MyPolicy")] // Required for this path.
 [HttpDelete("{action}/{id}")]
 public IActionResult MyDelete2(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
}
```
See Test CORS with [EnableCors] attribute and RequireCors method for instructions on testing code similar to the preceding.

## Enable CORS with attributes

Enabling CORS with the [EnableCors] attribute and applying a named policy to only those endpoints that require CORS provides the finest control.

The [EnableCors] attribute provides an alternative to applying CORS globally. The `[EnableCors]` attribute enables CORS for selected endpoints, rather than all endpoints:

- `[EnableCors]`specifies the default policy.
- `[EnableCors("{Policy String}")]`specifies a named policy.

The `[EnableCors]` attribute can be applied to:

- Razor Page `PageModel`
- Controller
- Controller action method

Different policies can be applied to controllers, page models, or action methods with the `[EnableCors]` attribute. When the `[EnableCors]` attribute is applied to a controller, page model, or action method, and CORS is enabled in middleware, **both** policies are applied. **We recommend against combining policies. Use the** `[EnableCors]` **attribute or middleware, not both in the same app.**

The following code applies a different policy to each method:

```
[Route("api/[controller]")]
[ApiController]
public class WidgetController : ControllerBase
{
 // GET api/values
 [EnableCors("AnotherPolicy")]
 [HttpGet]
 public ActionResult<IEnumerable<string>> Get()
 {
 return new string[] { "green widget", "red widget" };
 }
 // GET api/values/5
 [EnableCors("Policy1")]
 [HttpGet("{id}")]
 public ActionResult<string> Get(int id)
 {
 return id switch
 {
 1 => "green widget",
 2 => "red widget",
 _ => NotFound(),
 };
 }
}
```
The following code creates two CORS policies:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("Policy1",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
 options.AddPolicy("AnotherPolicy",
 policy =>
 {
 policy.WithOrigins("http://www.contoso.com")
 .AllowAnyHeader()
 .AllowAnyMethod();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.MapControllers();
app.Run();
```
For the finest control of limiting CORS requests:

- Use `[EnableCors("MyPolicy")]`with a named policy.
- Don't define a default policy.
- Don't use endpoint routing.

The code in the next section meets the preceding list.

### Disable CORS

The [DisableCors] attribute does **not** disable CORS that has been enabled by endpoint routing.

The following code defines the CORS policy `"MyPolicy"`:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: "MyPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com")
 .WithMethods("PUT", "DELETE", "GET");
 });
});
builder.Services.AddControllers();
builder.Services.AddRazorPages();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.UseEndpoints(endpoints => {
 endpoints.MapControllers();
 endpoints.MapRazorPages();
});
app.Run();
```
The following code disables CORS for the `GetValues2` action:

```
[EnableCors("MyPolicy")]
[Route("api/[controller]")]
[ApiController]
public class ValuesController : ControllerBase
{
 // GET api/values
 [HttpGet]
 public IActionResult Get() =>
 ControllerContext.MyDisplayRouteInfo();
 // GET api/values/5
 [HttpGet("{id}")]
 public IActionResult Get(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // PUT api/values/5
 [HttpPut("{id}")]
 public IActionResult Put(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // GET: api/values/GetValues2
 [DisableCors]
 [HttpGet("{action}")]
 public IActionResult GetValues2() =>
 ControllerContext.MyDisplayRouteInfo();
}
```
The preceding code:

- Doesn't enable CORS with endpoint routing.
- Doesn't define a default CORS policy.
- Uses [EnableCors("MyPolicy")] to enable the `"MyPolicy"`CORS policy for the controller.
- Disables CORS for the `GetValues2`method.

See Test CORS for instructions on testing the preceding code.

## CORS policy options

This section describes the various options that can be set in a CORS policy:

- Set the allowed origins
- Set the allowed HTTP methods
- Set the allowed request headers
- Set the exposed response headers
- Credentials in cross-origin requests
- Set the preflight expiration time

AddPolicy is called in `Program.cs`. For some options, it may be helpful to read the How CORS works section first.

## Set the allowed origins

AllowAnyOrigin: Allows CORS requests from all origins with any scheme (`http` or `https`). `AllowAnyOrigin` is insecure because *any website* can make cross-origin requests to the app.

Note

Specifying `AllowAnyOrigin` and `AllowCredentials` is an insecure configuration and can result in cross-site request forgery. The CORS service returns an invalid CORS response when an app is configured with both methods.

`AllowAnyOrigin` affects preflight requests and the `Access-Control-Allow-Origin` header. For more information, see the Preflight requests section.

SetIsOriginAllowedToAllowWildcardSubdomains: Sets the IsOriginAllowed property of the policy to be a function that allows origins to match a configured wildcard domain when evaluating if the origin is allowed.

```
var MyAllowSpecificOrigins = "_MyAllowSubdomainPolicy";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .SetIsOriginAllowedToAllowWildcardSubdomains();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
In the preceding code, `SetIsOriginAllowedToAllowWildcardSubdomains` is called with the wildcard origin `"https://*.example.com"`. This configuration allows CORS requests from any subdomain of `example.com`, such as `https://subdomain.example.com` or `https://api.example.com`. The `*` wildcard character must be included in the origin to enable wildcard subdomain matching.

### Set the allowed HTTP methods

- Allows any HTTP method:
- Affects preflight requests and the `Access-Control-Allow-Methods`header. For more information, see the Preflight requests section.

### Set the allowed request headers

To allow specific headers to be sent in a CORS request, called author request headers, call WithHeaders and specify the allowed headers:

```
using Microsoft.Net.Http.Headers;
var MyAllowSpecificOrigins = "_MyAllowSubdomainPolicy";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .WithHeaders(HeaderNames.ContentType, "x-custom-header");
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
To allow all author request headers, call AllowAnyHeader:

```
var MyAllowSpecificOrigins = "_MyAllowSubdomainPolicy";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .AllowAnyHeader();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
`AllowAnyHeader` affects preflight requests and the Access-Control-Request-Headers header. For more information, see the Preflight requests section.

A CORS middleware policy match to specific headers specified by `WithHeaders` is only possible when the headers sent in `Access-Control-Request-Headers` exactly match the headers stated in `WithHeaders`.

For instance, consider an app configured as follows:

```
app.UseCors(policy => policy.WithHeaders(HeaderNames.CacheControl));
```
CORS middleware declines a preflight request with the following request header because `Content-Language` (HeaderNames.ContentLanguage) isn't listed in `WithHeaders`:

```
Access-Control-Request-Headers: Cache-Control, Content-Language
```
The app returns a *200 OK* response but doesn't send the CORS headers back. Therefore, the browser doesn't attempt the cross-origin request.

### Set the exposed response headers

By default, the browser doesn't expose all of the response headers to the app. For more information, see W3C Cross-Origin Resource Sharing (Terminology): Simple Response Header.

The response headers that are available by default are:

- `Cache-Control`
- `Content-Language`
- `Content-Type`
- `Expires`
- `Last-Modified`
- `Pragma`

The CORS specification calls these headers *simple response headers*. To make other headers available to the app, call WithExposedHeaders:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MyExposeResponseHeadersPolicy",
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .WithExposedHeaders("x-custom-header");
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
### Credentials in cross-origin requests

Credentials require special handling in a CORS request. By default, the browser doesn't send credentials with a cross-origin request. Credentials include cookies and HTTP authentication schemes. To send credentials with a cross-origin request, the client must set `XMLHttpRequest.withCredentials` to `true`.

Using `XMLHttpRequest` directly:

```
var xhr = new XMLHttpRequest();
xhr.open('get', 'https://www.example.com/api/test');
xhr.withCredentials = true;
```
Using jQuery:

```
$.ajax({
 type: 'get',
 url: 'https://www.example.com/api/test',
 xhrFields: {
 withCredentials: true
 }
});
```
Using the Fetch API:

```
fetch('https://www.example.com/api/test', {
 credentials: 'include'
});
```
The server must allow the credentials. To allow cross-origin credentials, call AllowCredentials:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MyMyAllowCredentialsPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .AllowCredentials();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
The HTTP response includes an `Access-Control-Allow-Credentials` header, which tells the browser that the server allows credentials for a cross-origin request.

If the browser sends credentials but the response doesn't include a valid `Access-Control-Allow-Credentials` header, the browser doesn't expose the response to the app, and the cross-origin request fails.

Allowing cross-origin credentials is a security risk. A website at another domain can send a signed-in user's credentials to the app on the user's behalf without the user's knowledge.

The CORS specification also states that setting origins to `"*"` (all origins) is invalid if the `Access-Control-Allow-Credentials` header is present.

## Preflight requests

For some CORS requests, the browser sends an additional OPTIONS request before making the actual request. This request is called a preflight request. The browser can skip the preflight request if **all** the following conditions are true:

- The request method is GET, HEAD, or POST.
- The app doesn't set request headers other than `Accept`,`Accept-Language`,`Content-Language`,`Content-Type`, or`Last-Event-ID`.
- The `Content-Type`header, if set, has one of the following values:- `application/x-www-form-urlencoded`
- `multipart/form-data`
- `text/plain`

The rule on request headers set for the client request applies to headers that the app sets by calling `setRequestHeader` on the `XMLHttpRequest` object. The CORS specification calls these headers author request headers. The rule doesn't apply to headers the browser can set, such as `User-Agent`, `Host`, or `Content-Length`.

The following is an example response similar to the preflight request made from the **[Put test]** button in the Test CORS section of this document.

```
General:
Request URL: https://cors3.azurewebsites.net/api/values/5
Request Method: OPTIONS
Status Code: 204 No Content
Response Headers:
Access-Control-Allow-Methods: PUT,DELETE,GET
Access-Control-Allow-Origin: https://cors1.azurewebsites.net
Server: Microsoft-IIS/10.0
Set-Cookie: ARRAffinity=8f8...8;Path=/;HttpOnly;Domain=cors1.azurewebsites.net
Vary: Origin
Request Headers:
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9
Access-Control-Request-Method: PUT
Connection: keep-alive
Host: cors3.azurewebsites.net
Origin: https://cors1.azurewebsites.net
Referer: https://cors1.azurewebsites.net/
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
User-Agent: Mozilla/5.0
```
The preflight request uses the HTTP OPTIONS method. It may include the following headers:

- Access-Control-Request-Method: The HTTP method that will be used for the actual request.
- Access-Control-Request-Headers: A list of request headers that the app sets on the actual request. As stated earlier, this doesn't include headers that the browser sets, such as `User-Agent`.
- Access-Control-Allow-Methods

If the preflight request is denied, the app returns a `200 OK` response but doesn't set the CORS headers. Therefore, the browser doesn't attempt the cross-origin request. For an example of a denied preflight request, see the Test CORS section of this document.

Using the F12 tools, the console app shows an error similar to one of the following, depending on the browser:

- Firefox: Cross-Origin Request Blocked: The Same Origin Policy disallows reading the remote resource at `https://cors1.azurewebsites.net/api/TodoItems1/MyDelete2/5`. (Reason: CORS request did not succeed). Learn More
- Chromium based: Access to fetch at 'https://cors1.azurewebsites.net/api/TodoItems1/MyDelete2/5' from origin 'https://cors3.azurewebsites.net' has been blocked by CORS policy: Response to preflight request doesn't pass access control check: No 'Access-Control-Allow-Origin' header is present on the requested resource. If an opaque response serves your needs, set the request's mode to 'no-cors' to fetch the resource with CORS disabled.

To allow specific headers, call WithHeaders:

```
using Microsoft.Net.Http.Headers;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MyAllowHeadersPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .WithHeaders(HeaderNames.ContentType, "x-custom-header");
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
To allow all author request headers, call AllowAnyHeader:

```
using Microsoft.Net.Http.Headers;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MyAllowAllHeadersPolicy",
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .AllowAnyHeader();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
Browsers aren't consistent in how they set `Access-Control-Request-Headers`. If either:

- Headers are set to anything other than `"*"`
- AllowAnyHeader is called:
Include at least `Accept`,`Content-Type`, and`Origin`, plus any custom headers that you want to support.

### Automatic preflight request code

When the CORS policy is applied either:

- Globally by calling `app.UseCors`in`Program.cs`.
- Using the `[EnableCors]`attribute.

ASP.NET Core responds to the preflight OPTIONS request.

The Test CORS section of this document demonstrates this behavior.

### [HttpOptions] attribute for preflight requests

When CORS is enabled with the appropriate policy, ASP.NET Core generally responds to CORS preflight requests automatically.

The following code uses the [HttpOptions] attribute to create endpoints for OPTIONS requests:

```
[Route("api/[controller]")]
[ApiController]
public class TodoItems2Controller : ControllerBase
{
 // OPTIONS: api/TodoItems2/5
 [HttpOptions("{id}")]
 public IActionResult PreflightRoute(int id)
 {
 return NoContent();
 }
 // OPTIONS: api/TodoItems2
 [HttpOptions]
 public IActionResult PreflightRoute()
 {
 return NoContent();
 }
 [HttpPut("{id}")]
 public IActionResult PutTodoItem(int id)
 {
 if (id < 1)
 {
 return BadRequest();
 }
 return ControllerContext.MyDisplayRouteInfo(id);
 }
```
See Test CORS with [EnableCors] attribute and RequireCors method for instructions on testing the preceding code.

### Set the preflight expiration time

The `Access-Control-Max-Age` header specifies how long the response to the preflight request can be cached. To set this header, call SetPreflightMaxAge:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MySetPreflightExpirationPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .SetPreflightMaxAge(TimeSpan.FromSeconds(2520));
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
## Enable CORS on an endpoint

## How CORS works

This section describes what happens in a CORS request at the level of the HTTP messages.

- CORS is **not**a security feature. CORS is a W3C standard that allows a server to relax the same-origin policy.- For example, a malicious actor could use Cross-Site Scripting (XSS) against your site and execute a cross-site request to their CORS enabled site to steal information.

- An API isn't safer by allowing CORS.
- It's up to the client (browser) to enforce CORS. The server executes the request and returns the response, it's the client that returns an error and blocks the response. For example, any of the following tools will display the server response:
- Fiddler
- .NET HttpClient
- A web browser by entering the URL in the address bar.

- It's up to the client (browser) to enforce CORS. The server executes the request and returns the response, it's the client that returns an error and blocks the response. For example, any of the following tools will display the server response:
- It's a way for a server to allow browsers to execute a cross-origin XHR or Fetch API request that otherwise would be forbidden.
- Browsers without CORS can't do cross-origin requests. Before CORS, JSONP was used to circumvent this restriction. JSONP doesn't use XHR, it uses the `<script>`tag to receive the response. Scripts are allowed to be loaded cross-origin.

- Browsers without CORS can't do cross-origin requests. Before CORS, JSONP was used to circumvent this restriction. JSONP doesn't use XHR, it uses the

The CORS specification introduced several new HTTP headers that enable cross-origin requests. If a browser supports CORS, it sets these headers automatically for cross-origin requests. Custom JavaScript code isn't required to enable CORS.

Select the **PUT** test button on the deployed sample.
The `Origin` header:

- Provides the domain of the site that's making the request.
- Is required and must be different from the host.

**General headers**

```
Request URL: https://cors1.azurewebsites.net/api/values
Request Method: GET
Status Code: 200 OK
```
**Response headers**

```
Content-Encoding: gzip
Content-Type: text/plain; charset=utf-8
Server: Microsoft-IIS/10.0
Set-Cookie: ARRAffinity=8f...;Path=/;HttpOnly;Domain=cors1.azurewebsites.net
Transfer-Encoding: chunked
Vary: Accept-Encoding
X-Powered-By: ASP.NET
```
**Request headers**

```
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9
Connection: keep-alive
Host: cors1.azurewebsites.net
Origin: https://cors3.azurewebsites.net
Referer: https://cors3.azurewebsites.net/
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
User-Agent: Mozilla/5.0 ...
```
In `OPTIONS` requests, the server sets the **Response headers** `Access-Control-Allow-Origin: {allowed origin}` header in the response. For example, in the sample code, the ` Delete [EnableCors]` button `OPTIONS` request contains the following headers:

**General headers**

```
Request URL: https://cors3.azurewebsites.net/api/TodoItems2/MyDelete2/5
Request Method: OPTIONS
Status Code: 204 No Content
```
**Response headers**

```
Access-Control-Allow-Headers: Content-Type,x-custom-header
Access-Control-Allow-Methods: PUT,DELETE,GET,OPTIONS
Access-Control-Allow-Origin: https://cors1.azurewebsites.net
Server: Microsoft-IIS/10.0
Set-Cookie: ARRAffinity=8f...;Path=/;HttpOnly;Domain=cors3.azurewebsites.net
Vary: Origin
X-Powered-By: ASP.NET
```
**Request headers**

```
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9
Access-Control-Request-Headers: content-type
Access-Control-Request-Method: DELETE
Connection: keep-alive
Host: cors3.azurewebsites.net
Origin: https://cors1.azurewebsites.net
Referer: https://cors1.azurewebsites.net/test?number=2
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
User-Agent: Mozilla/5.0
```
In the preceding **Response headers**, the server sets the Access-Control-Allow-Origin header in the response. The `https://cors1.azurewebsites.net` value of this header matches the `Origin` header from the request.

If AllowAnyOrigin is called, the `Access-Control-Allow-Origin: *`, the wildcard value, is returned. `AllowAnyOrigin` allows any origin.

If the response doesn't include the `Access-Control-Allow-Origin` header, the cross-origin request fails. Specifically, the browser disallows the request. Even if the server returns a successful response, the browser doesn't make the response available to the client app.

### HTTP redirection to HTTPS causes ERR_INVALID_REDIRECT on the CORS preflight request

Requests to an endpoint using HTTP that are redirected to HTTPS by UseHttpsRedirection fail with `ERR_INVALID_REDIRECT on the CORS preflight request`.

API projects can reject HTTP requests rather than use `UseHttpsRedirection` to redirect requests to HTTPS.

## CORS in IIS

When deploying to IIS, CORS has to run before Windows Authentication if the server isn't configured to allow anonymous access. To support this scenario, the IIS CORS module needs to be installed and configured for the app.

## Test CORS

The sample download has code to test CORS. See how to download. The sample is an API project with Razor Pages added:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: "MyPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com",
 "https://cors1.azurewebsites.net",
 "https://cors3.azurewebsites.net",
 "https://localhost:44398",
 "https://localhost:5001")
 .WithMethods("PUT", "DELETE", "GET");
 });
});
builder.Services.AddControllers();
builder.Services.AddRazorPages();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.MapControllers();
app.MapRazorPages();
app.Run();
```
Warning

`WithOrigins("https://localhost:<port>");` should only be used for testing a sample app similar to the download sample code.

The following `ValuesController` provides the endpoints for testing:

```
[EnableCors("MyPolicy")]
[Route("api/[controller]")]
[ApiController]
public class ValuesController : ControllerBase
{
 // GET api/values
 [HttpGet]
 public IActionResult Get() =>
 ControllerContext.MyDisplayRouteInfo();
 // GET api/values/5
 [HttpGet("{id}")]
 public IActionResult Get(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // PUT api/values/5
 [HttpPut("{id}")]
 public IActionResult Put(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // GET: api/values/GetValues2
 [DisableCors]
 [HttpGet("{action}")]
 public IActionResult GetValues2() =>
 ControllerContext.MyDisplayRouteInfo();
}
```
MyDisplayRouteInfo is provided by the Rick.Docs.Samples.RouteInfo NuGet package and displays route information.

Test the preceding sample code by using one of the following approaches:

- Run the sample with `dotnet run`using the default URL of`https://localhost:5001`.
- Run the sample from Visual Studio with the port set to 44398 for a URL of `https://localhost:44398`.

Using a browser with the F12 tools:

- Select the - **Values**button and review the headers in the- **Network**tab.
- Select the - **PUT test**button. See Display OPTIONS requests for instructions on displaying the OPTIONS request. The- **PUT test**creates two requests, an OPTIONS preflight request and the PUT request.
- Select the - `GetValues2 [DisableCors]`- **Console**tab to see the CORS error. Depending on the browser, an error similar to the following is displayed:- Access to fetch at - `'https://cors1.azurewebsites.net/api/values/GetValues2'`from origin- `'https://cors3.azurewebsites.net'`has been blocked by CORS policy: No 'Access-Control-Allow-Origin' header is present on the requested resource. If an opaque response serves your needs, set the request's mode to 'no-cors' to fetch the resource with CORS disabled.

CORS-enabled endpoints can be tested with a tool, such as curl or Fiddler. When using a tool, the origin of the request specified by the `Origin` header must differ from the host receiving the request. If the request isn't *cross-origin* based on the value of the `Origin` header:

- There's no need for CORS middleware to process the request.
- CORS headers aren't returned in the response.

The following command uses `curl` to issue an OPTIONS request with information:

```
curl -X OPTIONS https://cors3.azurewebsites.net/api/TodoItems2/5 -i
```
### Test CORS with [EnableCors] attribute and RequireCors method

Consider the following code which uses endpoint routing to enable CORS on a per-endpoint basis using `RequireCors`:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: "MyPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com",
 "https://cors1.azurewebsites.net",
 "https://cors3.azurewebsites.net",
 "https://localhost:44398",
 "https://localhost:5001")
 .WithMethods("PUT", "DELETE", "GET");
 });
});
builder.Services.AddControllers();
builder.Services.AddRazorPages();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.UseEndpoints(endpoints =>
{
 endpoints.MapGet("/echo",
 context => context.Response.WriteAsync("echo"))
 .RequireCors("MyPolicy");
 endpoints.MapControllers();
 endpoints.MapRazorPages();
});
app.Run();
```
Notice that only the `/echo` endpoint is using the `RequireCors` to allow cross-origin requests using the specified policy. The controllers below enable CORS using [EnableCors] attribute.

The following `TodoItems1Controller` provides endpoints for testing:

```
[Route("api/[controller]")]
[ApiController]
public class TodoItems1Controller : ControllerBase
{
 // PUT: api/TodoItems1/5
 [HttpPut("{id}")]
 public IActionResult PutTodoItem(int id) {
 if (id < 1) {
 return Content($"ID = {id}");
 }
 return ControllerContext.MyDisplayRouteInfo(id);
 }
 // Delete: api/TodoItems1/5
 [HttpDelete("{id}")]
 public IActionResult MyDelete(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // GET: api/TodoItems1
 [HttpGet]
 public IActionResult GetTodoItems() =>
 ControllerContext.MyDisplayRouteInfo();
 [EnableCors("MyPolicy")]
 [HttpGet("{action}")]
 public IActionResult GetTodoItems2() =>
 ControllerContext.MyDisplayRouteInfo();
 // Delete: api/TodoItems1/MyDelete2/5
 [EnableCors("MyPolicy")]
 [HttpDelete("{action}/{id}")]
 public IActionResult MyDelete2(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
}
```
The **Delete [EnableCors]** and **GET [EnableCors]** buttons succeed, because the endpoints have `[EnableCors]` and respond to preflight requests. The other endpoints fails. The **GET** button fails, because the JavaScript sends:

```
 headers: {
 "Content-Type": "x-custom-header"
 },
```
The following `TodoItems2Controller` provides similar endpoints, but includes explicit code to respond to OPTIONS requests:

```
[Route("api/[controller]")]
[ApiController]
public class TodoItems2Controller : ControllerBase
{
 // OPTIONS: api/TodoItems2/5
 [HttpOptions("{id}")]
 public IActionResult PreflightRoute(int id)
 {
 return NoContent();
 }
 // OPTIONS: api/TodoItems2
 [HttpOptions]
 public IActionResult PreflightRoute()
 {
 return NoContent();
 }
 [HttpPut("{id}")]
 public IActionResult PutTodoItem(int id)
 {
 if (id < 1)
 {
 return BadRequest();
 }
 return ControllerContext.MyDisplayRouteInfo(id);
 }
 // [EnableCors] // Not needed as OPTIONS path provided.
 [HttpDelete("{id}")]
 public IActionResult MyDelete(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // [EnableCors] // Warning ASP0023 Route '{id}' conflicts with another action route.
 // An HTTP request that matches multiple routes results in an ambiguous
 // match error.
 [EnableCors("MyPolicy")] // Required for this path.
 [HttpGet]
 public IActionResult GetTodoItems() =>
 ControllerContext.MyDisplayRouteInfo();
 [HttpGet("{action}")]
 public IActionResult GetTodoItems2() =>
 ControllerContext.MyDisplayRouteInfo();
 [EnableCors("MyPolicy")] // Required for this path.
 [HttpDelete("{action}/{id}")]
 public IActionResult MyDelete2(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
}
```
The preceding code can be tested by deploying the sample to Azure.In the **Controller** drop down list, select **Preflight** and then **Set Controller**. All the CORS calls to the `TodoItems2Controller` endpoints succeed.

## Additional resources

By Rick Anderson and Kirk Larkin

This article shows how to enable CORS in an ASP.NET Core app.

Browser security prevents a web page from making requests to a different domain than the one that served the web page. This restriction is called the *same-origin policy*. The same-origin policy prevents a malicious site from reading sensitive data from another site. Sometimes, you might want to allow other sites to make cross-origin requests to your app. For more information, see the Mozilla CORS article.

Cross Origin Resource Sharing (CORS):

- Is a W3C standard that allows a server to relax the same-origin policy.
- Is **not**a security feature, CORS relaxes security. An API is not safer by allowing CORS. For more information, see How CORS works.
- Allows a server to explicitly allow some cross-origin requests while rejecting others.
- Is safer and more flexible than earlier techniques, such as JSONP.

View or download sample code (how to download)

## Same origin

Two URLs have the same origin if they have identical schemes, hosts, and ports (RFC 6454).

These two URLs have the same origin:

- `https://example.com/foo.html`
- `https://example.com/bar.html`

These URLs have different origins than the previous two URLs:

- `https://example.net`: Different domain
- `https://www.example.com/foo.html`: Different subdomain
- `http://example.com/foo.html`: Different scheme
- `https://example.com:9000/foo.html`: Different port

## Enable CORS

There are three ways to enable CORS:

- In middleware using a named policy or default policy.
- Using endpoint routing.
- With the [EnableCors] attribute.

Using the [EnableCors] attribute with a named policy provides the finest control in limiting endpoints that support CORS.

Warning

UseCors must be called in the correct order. For more information, see Middleware order. For example, `UseCors` must be called before UseResponseCaching when using `UseResponseCaching`.

Each approach is detailed in the following sections.

## CORS with named policy and middleware

CORS middleware handles cross-origin requests. The following code applies a CORS policy to all the app's endpoints with the specified origins:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
});
// services.AddResponseCaching();
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors(MyAllowSpecificOrigins);
app.UseAuthorization();
app.MapControllers();
app.Run();
```
The preceding code:

- Sets the policy name to `_myAllowSpecificOrigins`. The policy name is arbitrary.
- Calls the UseCors extension method and specifies the `_myAllowSpecificOrigins`CORS policy.`UseCors`adds the CORS middleware. The call to`UseCors`must be placed after`UseRouting`, but before`UseAuthorization`. For more information, see Middleware order.
- Calls AddCors with a lambda expression. The lambda takes a CorsPolicyBuilder object. Configuration options, such as `WithOrigins`, are described later in this article.
- Enables the `_myAllowSpecificOrigins`CORS policy for all controller endpoints. See endpoint routing to apply a CORS policy to specific endpoints.
- When using response caching middleware, call UseCors before UseResponseCaching.

With endpoint routing, the CORS middleware **must** be configured to execute between the calls to `UseRouting` and `UseEndpoints`.

The AddCors method call adds CORS services to the app's service container:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
});
// services.AddResponseCaching();
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors(MyAllowSpecificOrigins);
app.UseAuthorization();
app.MapControllers();
app.Run();
```
For more information, see CORS policy options in this document.

The CorsPolicyBuilder methods can be chained, as shown in the following code:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com")
 .AllowAnyHeader()
 .AllowAnyMethod();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors(MyAllowSpecificOrigins);
app.UseAuthorization();
app.MapControllers();
app.Run();
```
Note: The specified URL must **not** contain a trailing slash (`/`). If the URL terminates with `/`, the comparison returns `false` and no header is returned.

Warning

`UseCors` must be placed after `UseRouting` and before `UseAuthorization`. This is to ensure that CORS headers are included in the response for both authorized and unauthorized calls.

## UseCors and UseStaticFiles order

Typically, `UseStaticFiles` is called before `UseCors`. Apps that use JavaScript to retrieve static files cross site must call `UseCors` before `UseStaticFiles`.

### CORS with default policy and middleware

The following highlighted code enables the default CORS policy:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddDefaultPolicy(
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
});
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.MapControllers();
app.Run();
```
The preceding code applies the default CORS policy to all controller endpoints.

## Enable Cors with endpoint routing

Enabling CORS on a per-endpoint basis using `RequireCors` * does not support automatic preflight requests.* For more information, see this GitHub issue and Test CORS with endpoint routing and [HttpOptions].

With endpoint routing, CORS can be enabled on a per-endpoint basis using the RequireCors set of extension methods:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
});
builder.Services.AddControllers();
builder.Services.AddRazorPages();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.UseEndpoints(endpoints =>
{
 endpoints.MapGet("/echo",
 context => context.Response.WriteAsync("echo"))
 .RequireCors(MyAllowSpecificOrigins);
 endpoints.MapControllers()
 .RequireCors(MyAllowSpecificOrigins);
 endpoints.MapGet("/echo2",
 context => context.Response.WriteAsync("echo2"));
 endpoints.MapRazorPages();
});
app.Run();
```
In the preceding code:

- `app.UseCors`enables the CORS middleware. Because a default policy hasn't been configured,- `app.UseCors()`alone doesn't enable CORS.
- The `/echo`and controller endpoints allow cross-origin requests using the specified policy.
- The `/echo2`and Razor Pages endpoints do**not**allow cross-origin requests because no default policy was specified.

The [DisableCors] attribute does **not** disable CORS that has been enabled by endpoint routing with `RequireCors`.

See Test CORS with endpoint routing and [HttpOptions] for instructions on testing code similar to the preceding.

## Enable CORS with attributes

Enabling CORS with the [EnableCors] attribute and applying a named policy to only those endpoints that require CORS provides the finest control.

The [EnableCors] attribute provides an alternative to applying CORS globally. The `[EnableCors]` attribute enables CORS for selected endpoints, rather than all endpoints:

- `[EnableCors]`specifies the default policy.
- `[EnableCors("{Policy String}")]`specifies a named policy.

The `[EnableCors]` attribute can be applied to:

- Razor Page `PageModel`
- Controller
- Controller action method

Different policies can be applied to controllers, page models, or action methods with the `[EnableCors]` attribute. When the `[EnableCors]` attribute is applied to a controller, page model, or action method, and CORS is enabled in middleware, **both** policies are applied. **We recommend against combining policies. Use the** `[EnableCors]` **attribute or middleware, not both in the same app.**

The following code applies a different policy to each method:

```
[Route("api/[controller]")]
[ApiController]
public class WidgetController : ControllerBase
{
 // GET api/values
 [EnableCors("AnotherPolicy")]
 [HttpGet]
 public ActionResult<IEnumerable<string>> Get()
 {
 return new string[] { "green widget", "red widget" };
 }
 // GET api/values/5
 [EnableCors("Policy1")]
 [HttpGet("{id}")]
 public ActionResult<string> Get(int id)
 {
 return id switch
 {
 1 => "green widget",
 2 => "red widget",
 _ => NotFound(),
 };
 }
}
```
The following code creates two CORS policies:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("Policy1",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
 options.AddPolicy("AnotherPolicy",
 policy =>
 {
 policy.WithOrigins("http://www.contoso.com")
 .AllowAnyHeader()
 .AllowAnyMethod();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.MapControllers();
app.Run();
```
For the finest control of limiting CORS requests:

- Use `[EnableCors("MyPolicy")]`with a named policy.
- Don't define a default policy.
- Don't use endpoint routing.

The code in the next section meets the preceding list.

### Disable CORS

The [DisableCors] attribute does **not** disable CORS that has been enabled by endpoint routing.

The following code defines the CORS policy `"MyPolicy"`:

```
var MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: "MyPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com")
 .WithMethods("PUT", "DELETE", "GET");
 });
});
builder.Services.AddControllers();
builder.Services.AddRazorPages();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.MapControllers();
app.MapRazorPages();
app.Run();
```
The following code disables CORS for the `GetValues2` action:

```
[EnableCors("MyPolicy")]
[Route("api/[controller]")]
[ApiController]
public class ValuesController : ControllerBase
{
 // GET api/values
 [HttpGet]
 public IActionResult Get() =>
 ControllerContext.MyDisplayRouteInfo();
 // GET api/values/5
 [HttpGet("{id}")]
 public IActionResult Get(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // PUT api/values/5
 [HttpPut("{id}")]
 public IActionResult Put(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // GET: api/values/GetValues2
 [DisableCors]
 [HttpGet("{action}")]
 public IActionResult GetValues2() =>
 ControllerContext.MyDisplayRouteInfo();
}
```
The preceding code:

- Doesn't enable CORS with endpoint routing.
- Doesn't define a default CORS policy.
- Uses [EnableCors("MyPolicy")] to enable the `"MyPolicy"`CORS policy for the controller.
- Disables CORS for the `GetValues2`method.

See Test CORS for instructions on testing the preceding code.

## CORS policy options

This section describes the various options that can be set in a CORS policy:

- Set the allowed origins
- Set the allowed HTTP methods
- Set the allowed request headers
- Set the exposed response headers
- Credentials in cross-origin requests
- Set the preflight expiration time

AddPolicy is called in `Program.cs`. For some options, it may be helpful to read the How CORS works section first.

## Set the allowed origins

AllowAnyOrigin: Allows CORS requests from all origins with any scheme (`http` or `https`). `AllowAnyOrigin` is insecure because *any website* can make cross-origin requests to the app.

Note

Specifying `AllowAnyOrigin` and `AllowCredentials` is an insecure configuration and can result in cross-site request forgery. The CORS service returns an invalid CORS response when an app is configured with both methods.

`AllowAnyOrigin` affects preflight requests and the `Access-Control-Allow-Origin` header. For more information, see the Preflight requests section.

SetIsOriginAllowedToAllowWildcardSubdomains: Sets the IsOriginAllowed property of the policy to be a function that allows origins to match a configured wildcard domain when evaluating if the origin is allowed.

```
var MyAllowSpecificOrigins = "_MyAllowSubdomainPolicy";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .SetIsOriginAllowedToAllowWildcardSubdomains();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
In the preceding code, `SetIsOriginAllowedToAllowWildcardSubdomains` is called with the wildcard origin `"https://*.example.com"`. This configuration allows CORS requests from any subdomain of `example.com`, such as `https://subdomain.example.com` or `https://api.example.com`. The `*` wildcard character must be included in the origin to enable wildcard subdomain matching.

### Set the allowed HTTP methods

- Allows any HTTP method:
- Affects preflight requests and the `Access-Control-Allow-Methods`header. For more information, see the Preflight requests section.

### Set the allowed request headers

To allow specific headers to be sent in a CORS request, called author request headers, call WithHeaders and specify the allowed headers:

```
using Microsoft.Net.Http.Headers;
var MyAllowSpecificOrigins = "_MyAllowSubdomainPolicy";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .WithHeaders(HeaderNames.ContentType, "x-custom-header");
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
To allow all author request headers, call AllowAnyHeader:

```
var MyAllowSpecificOrigins = "_MyAllowSubdomainPolicy";
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .AllowAnyHeader();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
`AllowAnyHeader` affects preflight requests and the Access-Control-Request-Headers header. For more information, see the Preflight requests section.

A CORS middleware policy match to specific headers specified by `WithHeaders` is only possible when the headers sent in `Access-Control-Request-Headers` exactly match the headers stated in `WithHeaders`.

For instance, consider an app configured as follows:

```
app.UseCors(policy => policy.WithHeaders(HeaderNames.CacheControl));
```
CORS middleware declines a preflight request with the following request header because `Content-Language` (HeaderNames.ContentLanguage) isn't listed in `WithHeaders`:

```
Access-Control-Request-Headers: Cache-Control, Content-Language
```
The app returns a *200 OK* response but doesn't send the CORS headers back. Therefore, the browser doesn't attempt the cross-origin request.

### Set the exposed response headers

By default, the browser doesn't expose all of the response headers to the app. For more information, see W3C Cross-Origin Resource Sharing (Terminology): Simple Response Header.

The response headers that are available by default are:

- `Cache-Control`
- `Content-Language`
- `Content-Type`
- `Expires`
- `Last-Modified`
- `Pragma`

The CORS specification calls these headers *simple response headers*. To make other headers available to the app, call WithExposedHeaders:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MyExposeResponseHeadersPolicy",
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .WithExposedHeaders("x-custom-header");
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
### Credentials in cross-origin requests

Credentials require special handling in a CORS request. By default, the browser doesn't send credentials with a cross-origin request. Credentials include cookies and HTTP authentication schemes. To send credentials with a cross-origin request, the client must set `XMLHttpRequest.withCredentials` to `true`.

Using `XMLHttpRequest` directly:

```
var xhr = new XMLHttpRequest();
xhr.open('get', 'https://www.example.com/api/test');
xhr.withCredentials = true;
```
Using jQuery:

```
$.ajax({
 type: 'get',
 url: 'https://www.example.com/api/test',
 xhrFields: {
 withCredentials: true
 }
});
```
Using the Fetch API:

```
fetch('https://www.example.com/api/test', {
 credentials: 'include'
});
```
The server must allow the credentials. To allow cross-origin credentials, call AllowCredentials:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MyMyAllowCredentialsPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .AllowCredentials();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
The HTTP response includes an `Access-Control-Allow-Credentials` header, which tells the browser that the server allows credentials for a cross-origin request.

If the browser sends credentials but the response doesn't include a valid `Access-Control-Allow-Credentials` header, the browser doesn't expose the response to the app, and the cross-origin request fails.

Allowing cross-origin credentials is a security risk. A website at another domain can send a signed-in user's credentials to the app on the user's behalf without the user's knowledge.

The CORS specification also states that setting origins to `"*"` (all origins) is invalid if the `Access-Control-Allow-Credentials` header is present.

## Preflight requests

For some CORS requests, the browser sends an additional OPTIONS request before making the actual request. This request is called a preflight request. The browser can skip the preflight request if **all** the following conditions are true:

- The request method is GET, HEAD, or POST.
- The app doesn't set request headers other than `Accept`,`Accept-Language`,`Content-Language`,`Content-Type`, or`Last-Event-ID`.
- The `Content-Type`header, if set, has one of the following values:- `application/x-www-form-urlencoded`
- `multipart/form-data`
- `text/plain`

The rule on request headers set for the client request applies to headers that the app sets by calling `setRequestHeader` on the `XMLHttpRequest` object. The CORS specification calls these headers author request headers. The rule doesn't apply to headers the browser can set, such as `User-Agent`, `Host`, or `Content-Length`.

The following is an example response similar to the preflight request made from the **[Put test]** button in the Test CORS section of this document.

```
General:
Request URL: https://cors3.azurewebsites.net/api/values/5
Request Method: OPTIONS
Status Code: 204 No Content
Response Headers:
Access-Control-Allow-Methods: PUT,DELETE,GET
Access-Control-Allow-Origin: https://cors1.azurewebsites.net
Server: Microsoft-IIS/10.0
Set-Cookie: ARRAffinity=8f8...8;Path=/;HttpOnly;Domain=cors1.azurewebsites.net
Vary: Origin
Request Headers:
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9
Access-Control-Request-Method: PUT
Connection: keep-alive
Host: cors3.azurewebsites.net
Origin: https://cors1.azurewebsites.net
Referer: https://cors1.azurewebsites.net/
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
User-Agent: Mozilla/5.0
```
The preflight request uses the HTTP OPTIONS method. It may include the following headers:

- Access-Control-Request-Method: The HTTP method that will be used for the actual request.
- Access-Control-Request-Headers: A list of request headers that the app sets on the actual request. As stated earlier, this doesn't include headers that the browser sets, such as `User-Agent`.
- Access-Control-Allow-Methods

If the preflight request is denied, the app returns a `200 OK` response but doesn't set the CORS headers. Therefore, the browser doesn't attempt the cross-origin request. For an example of a denied preflight request, see the Test CORS section of this document.

Using the F12 tools, the console app shows an error similar to one of the following, depending on the browser:

- Firefox: Cross-Origin Request Blocked: The Same Origin Policy disallows reading the remote resource at `https://cors1.azurewebsites.net/api/TodoItems1/MyDelete2/5`. (Reason: CORS request did not succeed). Learn More
- Chromium based: Access to fetch at 'https://cors1.azurewebsites.net/api/TodoItems1/MyDelete2/5' from origin 'https://cors3.azurewebsites.net' has been blocked by CORS policy: Response to preflight request doesn't pass access control check: No 'Access-Control-Allow-Origin' header is present on the requested resource. If an opaque response serves your needs, set the request's mode to 'no-cors' to fetch the resource with CORS disabled.

To allow specific headers, call WithHeaders:

```
using Microsoft.Net.Http.Headers;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MyAllowHeadersPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .WithHeaders(HeaderNames.ContentType, "x-custom-header");
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
To allow all author request headers, call AllowAnyHeader:

```
using Microsoft.Net.Http.Headers;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MyAllowAllHeadersPolicy",
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .AllowAnyHeader();
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
Browsers aren't consistent in how they set `Access-Control-Request-Headers`. If either:

- Headers are set to anything other than `"*"`
- AllowAnyHeader is called:
Include at least `Accept`,`Content-Type`, and`Origin`, plus any custom headers that you want to support.

### Automatic preflight request code

When the CORS policy is applied either:

- Globally by calling `app.UseCors`in`Program.cs`.
- Using the `[EnableCors]`attribute.

ASP.NET Core responds to the preflight OPTIONS request.

Enabling CORS on a per-endpoint basis using `RequireCors` currently does **not** support automatic preflight requests.

The Test CORS section of this document demonstrates this behavior.

### [HttpOptions] attribute for preflight requests

When CORS is enabled with the appropriate policy, ASP.NET Core generally responds to CORS preflight requests automatically. In some scenarios, this may not be the case. For example, using CORS with endpoint routing.

The following code uses the [HttpOptions] attribute to create endpoints for OPTIONS requests:

```
[Route("api/[controller]")]
[ApiController]
public class TodoItems2Controller : ControllerBase
{
 // OPTIONS: api/TodoItems2/5
 [HttpOptions("{id}")]
 public IActionResult PreflightRoute(int id)
 {
 return NoContent();
 }
 // OPTIONS: api/TodoItems2
 [HttpOptions]
 public IActionResult PreflightRoute()
 {
 return NoContent();
 }
 [HttpPut("{id}")]
 public IActionResult PutTodoItem(int id)
 {
 if (id < 1)
 {
 return BadRequest();
 }
 return ControllerContext.MyDisplayRouteInfo(id);
 }
```
See Test CORS with endpoint routing and [HttpOptions] for instructions on testing the preceding code.

### Set the preflight expiration time

The `Access-Control-Max-Age` header specifies how long the response to the preflight request can be cached. To set this header, call SetPreflightMaxAge:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy("MySetPreflightExpirationPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .SetPreflightMaxAge(TimeSpan.FromSeconds(2520));
 });
});
builder.Services.AddControllers();
var app = builder.Build();
```
## How CORS works

This section describes what happens in a CORS request at the level of the HTTP messages.

- CORS is **not**a security feature. CORS is a W3C standard that allows a server to relax the same-origin policy.- For example, a malicious actor could use Cross-Site Scripting (XSS) against your site and execute a cross-site request to their CORS enabled site to steal information.

- An API isn't safer by allowing CORS.
- It's up to the client (browser) to enforce CORS. The server executes the request and returns the response, it's the client that returns an error and blocks the response. For example, any of the following tools will display the server response:
- Fiddler
- .NET HttpClient
- A web browser by entering the URL in the address bar.

- It's up to the client (browser) to enforce CORS. The server executes the request and returns the response, it's the client that returns an error and blocks the response. For example, any of the following tools will display the server response:
- It's a way for a server to allow browsers to execute a cross-origin XHR or Fetch API request that otherwise would be forbidden.
- Browsers without CORS can't do cross-origin requests. Before CORS, JSONP was used to circumvent this restriction. JSONP doesn't use XHR, it uses the `<script>`tag to receive the response. Scripts are allowed to be loaded cross-origin.

- Browsers without CORS can't do cross-origin requests. Before CORS, JSONP was used to circumvent this restriction. JSONP doesn't use XHR, it uses the

The CORS specification introduced several new HTTP headers that enable cross-origin requests. If a browser supports CORS, it sets these headers automatically for cross-origin requests. Custom JavaScript code isn't required to enable CORS.

The following is an example of a cross-origin request from the **Values** test button to `https://cors1.azurewebsites.net/api/values`. The `Origin` header:

- Provides the domain of the site that's making the request.
- Is required and must be different from the host.

**General headers**

```
Request URL: https://cors1.azurewebsites.net/api/values
Request Method: GET
Status Code: 200 OK
```
**Response headers**

```
Content-Encoding: gzip
Content-Type: text/plain; charset=utf-8
Server: Microsoft-IIS/10.0
Set-Cookie: ARRAffinity=8f...;Path=/;HttpOnly;Domain=cors1.azurewebsites.net
Transfer-Encoding: chunked
Vary: Accept-Encoding
X-Powered-By: ASP.NET
```
**Request headers**

```
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9
Connection: keep-alive
Host: cors1.azurewebsites.net
Origin: https://cors3.azurewebsites.net
Referer: https://cors3.azurewebsites.net/
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
User-Agent: Mozilla/5.0 ...
```
In `OPTIONS` requests, the server sets the **Response headers** `Access-Control-Allow-Origin: {allowed origin}` header in the response. For example, the deployed sample, Delete button `OPTIONS` request contains the following headers:

**General headers**

```
Request URL: https://cors3.azurewebsites.net/api/TodoItems2/MyDelete2/5
Request Method: OPTIONS
Status Code: 204 No Content
```
**Response headers**

```
Access-Control-Allow-Headers: Content-Type,x-custom-header
Access-Control-Allow-Methods: PUT,DELETE,GET,OPTIONS
Access-Control-Allow-Origin: https://cors1.azurewebsites.net
Server: Microsoft-IIS/10.0
Set-Cookie: ARRAffinity=8f...;Path=/;HttpOnly;Domain=cors3.azurewebsites.net
Vary: Origin
X-Powered-By: ASP.NET
```
**Request headers**

```
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9
Access-Control-Request-Headers: content-type
Access-Control-Request-Method: DELETE
Connection: keep-alive
Host: cors3.azurewebsites.net
Origin: https://cors1.azurewebsites.net
Referer: https://cors1.azurewebsites.net/test?number=2
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
User-Agent: Mozilla/5.0
```
In the preceding **Response headers**, the server sets the Access-Control-Allow-Origin header in the response. The `https://cors1.azurewebsites.net` value of this header matches the `Origin` header from the request.

If AllowAnyOrigin is called, the `Access-Control-Allow-Origin: *`, the wildcard value, is returned. `AllowAnyOrigin` allows any origin.

If the response doesn't include the `Access-Control-Allow-Origin` header, the cross-origin request fails. Specifically, the browser disallows the request. Even if the server returns a successful response, the browser doesn't make the response available to the client app.

### HTTP redirection to HTTPS causes ERR_INVALID_REDIRECT on the CORS preflight request

Requests to an endpoint using HTTP that are redirected to HTTPS by UseHttpsRedirection fail with `ERR_INVALID_REDIRECT on the CORS preflight request`.

API projects can reject HTTP requests rather than use `UseHttpsRedirection` to redirect requests to HTTPS.

### Display OPTIONS requests

By default, the Chrome and Edge browsers don't show OPTIONS requests on the network tab of the F12 tools. To display OPTIONS requests in these browsers:

- `chrome://flags/#out-of-blink-cors`or- `edge://flags/#out-of-blink-cors`
- disable the flag.
- restart.

Firefox shows OPTIONS requests by default.

## CORS in IIS

When deploying to IIS, CORS has to run before Windows Authentication if the server isn't configured to allow anonymous access. To support this scenario, the IIS CORS module needs to be installed and configured for the app.

## Test CORS

The sample download has code to test CORS. See how to download. The sample is an API project with Razor Pages added:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: "MyPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com",
 "https://cors1.azurewebsites.net",
 "https://cors3.azurewebsites.net",
 "https://localhost:44398",
 "https://localhost:5001")
 .WithMethods("PUT", "DELETE", "GET");
 });
});
builder.Services.AddControllers();
builder.Services.AddRazorPages();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.MapControllers();
app.MapRazorPages();
app.Run();
```
Warning

`WithOrigins("https://localhost:<port>");` should only be used for testing a sample app similar to the download sample code.

The following `ValuesController` provides the endpoints for testing:

```
[EnableCors("MyPolicy")]
[Route("api/[controller]")]
[ApiController]
public class ValuesController : ControllerBase
{
 // GET api/values
 [HttpGet]
 public IActionResult Get() =>
 ControllerContext.MyDisplayRouteInfo();
 // GET api/values/5
 [HttpGet("{id}")]
 public IActionResult Get(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // PUT api/values/5
 [HttpPut("{id}")]
 public IActionResult Put(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // GET: api/values/GetValues2
 [DisableCors]
 [HttpGet("{action}")]
 public IActionResult GetValues2() =>
 ControllerContext.MyDisplayRouteInfo();
}
```
MyDisplayRouteInfo is provided by the Rick.Docs.Samples.RouteInfo NuGet package and displays route information.

Test the preceding sample code by using one of the following approaches:

- Run the sample with `dotnet run`using the default URL of`https://localhost:5001`.
- Run the sample from Visual Studio with the port set to 44398 for a URL of `https://localhost:44398`.

Using a browser with the F12 tools:

- Select the - **Values**button and review the headers in the- **Network**tab.
- Select the - **PUT test**button. See Display OPTIONS requests for instructions on displaying the OPTIONS request. The- **PUT test**creates two requests, an OPTIONS preflight request and the PUT request.
- Select the - `GetValues2 [DisableCors]`- **Console**tab to see the CORS error. Depending on the browser, an error similar to the following is displayed:- Access to fetch at - `'https://cors1.azurewebsites.net/api/values/GetValues2'`from origin- `'https://cors3.azurewebsites.net'`has been blocked by CORS policy: No 'Access-Control-Allow-Origin' header is present on the requested resource. If an opaque response serves your needs, set the request's mode to 'no-cors' to fetch the resource with CORS disabled.

CORS-enabled endpoints can be tested with a tool, such as curl or Fiddler. When using a tool, the origin of the request specified by the `Origin` header must differ from the host receiving the request. If the request isn't *cross-origin* based on the value of the `Origin` header:

- There's no need for CORS middleware to process the request.
- CORS headers aren't returned in the response.

The following command uses `curl` to issue an OPTIONS request with information:

```
curl -X OPTIONS https://cors3.azurewebsites.net/api/TodoItems2/5 -i
```
### Test CORS with endpoint routing and [HttpOptions]

Enabling CORS on a per-endpoint basis using `RequireCors` currently does **not** support automatic preflight requests. Consider the following code which uses endpoint routing to enable CORS:

```
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
 options.AddPolicy(name: "MyPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com",
 "https://cors1.azurewebsites.net",
 "https://cors3.azurewebsites.net",
 "https://localhost:44398",
 "https://localhost:5001")
 .WithMethods("PUT", "DELETE", "GET");
 });
});
builder.Services.AddControllers();
builder.Services.AddRazorPages();
var app = builder.Build();
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseCors();
app.UseAuthorization();
app.MapControllers();
app.MapRazorPages();
app.Run();
```
The following `TodoItems1Controller` provides endpoints for testing:

```
[Route("api/[controller]")]
[ApiController]
public class TodoItems1Controller : ControllerBase
{
 // PUT: api/TodoItems1/5
 [HttpPut("{id}")]
 public IActionResult PutTodoItem(int id)
 {
 if (id < 1)
 {
 return Content($"ID = {id}");
 }
 return ControllerContext.MyDisplayRouteInfo(id);
 }
 // Delete: api/TodoItems1/5
 [HttpDelete("{id}")]
 public IActionResult MyDelete(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // GET: api/TodoItems1
 [HttpGet]
 public IActionResult GetTodoItems() =>
 ControllerContext.MyDisplayRouteInfo();
 [EnableCors]
 [HttpGet("{action}")]
 public IActionResult GetTodoItems2() =>
 ControllerContext.MyDisplayRouteInfo();
 // Delete: api/TodoItems1/MyDelete2/5
 [EnableCors]
 [HttpDelete("{action}/{id}")]
 public IActionResult MyDelete2(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
}
```
Test the preceding code from the test page (`https://cors1.azurewebsites.net/test?number=1`) of the deployed sample.

The **Delete [EnableCors]** and **GET [EnableCors]** buttons succeed, because the endpoints have `[EnableCors]` and respond to preflight requests. The other endpoints fails. The **GET** button fails, because the JavaScript sends:

```
 headers: {
 "Content-Type": "x-custom-header"
 },
```
The following `TodoItems2Controller` provides similar endpoints, but includes explicit code to respond to OPTIONS requests:

```
[Route("api/[controller]")]
[ApiController]
public class TodoItems2Controller : ControllerBase
{
 // OPTIONS: api/TodoItems2/5
 [HttpOptions("{id}")]
 public IActionResult PreflightRoute(int id)
 {
 return NoContent();
 }
 // OPTIONS: api/TodoItems2
 [HttpOptions]
 public IActionResult PreflightRoute()
 {
 return NoContent();
 }
 [HttpPut("{id}")]
 public IActionResult PutTodoItem(int id)
 {
 if (id < 1)
 {
 return BadRequest();
 }
 return ControllerContext.MyDisplayRouteInfo(id);
 }
 // [EnableCors] // Not needed as OPTIONS path provided
 [HttpDelete("{id}")]
 public IActionResult MyDelete(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 [EnableCors] // Rquired for this path
 [HttpGet]
 public IActionResult GetTodoItems() =>
 ControllerContext.MyDisplayRouteInfo();
 [HttpGet("{action}")]
 public IActionResult GetTodoItems2() =>
 ControllerContext.MyDisplayRouteInfo();
 [EnableCors] // Rquired for this path
 [HttpDelete("{action}/{id}")]
 public IActionResult MyDelete2(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
}
```
The preceding code can be tested by deploying the sample to Azure.In the **Controller** drop down list, select **Preflight** and then **Set Controller**. All the CORS calls to the `TodoItems2Controller` endpoints succeed.

## Additional resources

By Rick Anderson and Kirk Larkin

This article shows how to enable CORS in an ASP.NET Core app.

Browser security prevents a web page from making requests to a different domain than the one that served the web page. This restriction is called the *same-origin policy*. The same-origin policy prevents a malicious site from reading sensitive data from another site. Sometimes, you might want to allow other sites to make cross-origin requests to your app. For more information, see the Mozilla CORS article.

Cross Origin Resource Sharing (CORS):

- Is a W3C standard that allows a server to relax the same-origin policy.
- Is **not**a security feature, CORS relaxes security. An API is not safer by allowing CORS. For more information, see How CORS works.
- Allows a server to explicitly allow some cross-origin requests while rejecting others.
- Is safer and more flexible than earlier techniques, such as JSONP.

View or download sample code (how to download)

## Same origin

Two URLs have the same origin if they have identical schemes, hosts, and ports (RFC 6454).

These two URLs have the same origin:

- `https://example.com/foo.html`
- `https://example.com/bar.html`

These URLs have different origins than the previous two URLs:

- `https://example.net`: Different domain
- `https://www.example.com/foo.html`: Different subdomain
- `http://example.com/foo.html`: Different scheme
- `https://example.com:9000/foo.html`: Different port

## Enable CORS

There are three ways to enable CORS:

- In middleware using a named policy or default policy.
- Using endpoint routing.
- With the [EnableCors] attribute.

Using the [EnableCors] attribute with a named policy provides the finest control in limiting endpoints that support CORS.

Warning

UseCors must be called in the correct order. For more information, see Middleware order. For example, `UseCors` must be called before UseResponseCaching when using `UseResponseCaching`.

Each approach is detailed in the following sections.

## CORS with named policy and middleware

CORS middleware handles cross-origin requests. The following code applies a CORS policy to all the app's endpoints with the specified origins:

```
public class Startup
{
 readonly string MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
 public void ConfigureServices(IServiceCollection services)
 {
 services.AddCors(options =>
 {
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
 });
 // services.AddResponseCaching();
 services.AddControllers();
 }
 public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
 {
 if (env.IsDevelopment())
 {
 app.UseDeveloperExceptionPage();
 }
 app.UseHttpsRedirection();
 app.UseStaticFiles();
 app.UseRouting();
 app.UseCors(MyAllowSpecificOrigins);
 // app.UseResponseCaching();
 app.UseAuthorization();
 app.UseEndpoints(endpoints =>
 {
 endpoints.MapControllers();
 });
 }
}
```
The preceding code:

- Sets the policy name to `_myAllowSpecificOrigins`. The policy name is arbitrary.
- Calls the UseCors extension method and specifies the `_myAllowSpecificOrigins`CORS policy.`UseCors`adds the CORS middleware. The call to`UseCors`must be placed after`UseRouting`, but before`UseAuthorization`. For more information, see Middleware order.
- Calls AddCors with a lambda expression. The lambda takes a CorsPolicyBuilder object. Configuration options, such as `WithOrigins`, are described later in this article.
- Enables the `_myAllowSpecificOrigins`CORS policy for all controller endpoints. See endpoint routing to apply a CORS policy to specific endpoints.
- When using response caching middleware, call UseCors before UseResponseCaching.

With endpoint routing, the CORS middleware **must** be configured to execute between the calls to `UseRouting` and `UseEndpoints`.

See Test CORS for instructions on testing code similar to the preceding code.

The AddCors method call adds CORS services to the app's service container:

```
public class Startup
{
 readonly string MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
 public void ConfigureServices(IServiceCollection services)
 {
 services.AddCors(options =>
 {
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
 });
 // services.AddResponseCaching();
 services.AddControllers();
 }
```
For more information, see CORS policy options in this document.

The CorsPolicyBuilder methods can be chained, as shown in the following code:

```
public void ConfigureServices(IServiceCollection services)
{
 services.AddCors(options =>
 {
 options.AddPolicy(MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com")
 .AllowAnyHeader()
 .AllowAnyMethod();
 });
 });
 services.AddControllers();
}
```
Note: The specified URL must **not** contain a trailing slash (`/`). If the URL terminates with `/`, the comparison returns `false` and no header is returned.

### CORS with default policy and middleware

The following highlighted code enables the default CORS policy:

```
public class Startup
{
 public void ConfigureServices(IServiceCollection services)
 {
 services.AddCors(options =>
 {
 options.AddDefaultPolicy(
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
 });
 services.AddControllers();
 }
 public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
 {
 if (env.IsDevelopment())
 {
 app.UseDeveloperExceptionPage();
 }
 app.UseHttpsRedirection();
 app.UseStaticFiles();
 app.UseRouting();
 app.UseCors();
 app.UseAuthorization();
 app.UseEndpoints(endpoints =>
 {
 endpoints.MapControllers();
 });
 }
}
```
The preceding code applies the default CORS policy to all controller endpoints.

## Enable Cors with endpoint routing

Enabling CORS on a per-endpoint basis using `RequireCors` * does not support automatic preflight requests.* For more information, see this GitHub issue and Test CORS with endpoint routing and [HttpOptions].

With endpoint routing, CORS can be enabled on a per-endpoint basis using the RequireCors set of extension methods:

```
public class Startup
{
 readonly string MyAllowSpecificOrigins = "_myAllowSpecificOrigins";
 public void ConfigureServices(IServiceCollection services)
 {
 services.AddCors(options =>
 {
 options.AddPolicy(name: MyAllowSpecificOrigins,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
 });
 services.AddControllers();
 services.AddRazorPages();
 }
 public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
 {
 if (env.IsDevelopment())
 {
 app.UseDeveloperExceptionPage();
 }
 app.UseHttpsRedirection();
 app.UseStaticFiles();
 app.UseRouting();
 app.UseCors();
 app.UseAuthorization();
 app.UseEndpoints(endpoints =>
 {
 endpoints.MapGet("/echo",
 context => context.Response.WriteAsync("echo"))
 .RequireCors(MyAllowSpecificOrigins);
 endpoints.MapControllers()
 .RequireCors(MyAllowSpecificOrigins);
 endpoints.MapGet("/echo2",
 context => context.Response.WriteAsync("echo2"));
 endpoints.MapRazorPages();
 });
 }
}
```
In the preceding code:

- `app.UseCors`enables the CORS middleware. Because a default policy hasn't been configured,- `app.UseCors()`alone doesn't enable CORS.
- The `/echo`and controller endpoints allow cross-origin requests using the specified policy.
- The `/echo2`and Razor Pages endpoints do**not**allow cross-origin requests because no default policy was specified.

The [DisableCors] attribute does **not** disable CORS that has been enabled by endpoint routing with `RequireCors`.

See Test CORS with endpoint routing and [HttpOptions] for instructions on testing code similar to the preceding.

## Enable CORS with attributes

Enabling CORS with the [EnableCors] attribute and applying a named policy to only those endpoints that require CORS provides the finest control.

The [EnableCors] attribute provides an alternative to applying CORS globally. The `[EnableCors]` attribute enables CORS for selected endpoints, rather than all endpoints:

- `[EnableCors]`specifies the default policy.
- `[EnableCors("{Policy String}")]`specifies a named policy.

The `[EnableCors]` attribute can be applied to:

- Razor Page `PageModel`
- Controller
- Controller action method

Different policies can be applied to controllers, page models, or action methods with the `[EnableCors]` attribute. When the `[EnableCors]` attribute is applied to a controller, page model, or action method, and CORS is enabled in middleware, **both** policies are applied. **We recommend against combining policies. Use the** `[EnableCors]` **attribute or middleware, not both in the same app.**

The following code applies a different policy to each method:

```
[Route("api/[controller]")]
[ApiController]
public class WidgetController : ControllerBase
{
 // GET api/values
 [EnableCors("AnotherPolicy")]
 [HttpGet]
 public ActionResult<IEnumerable<string>> Get()
 {
 return new string[] { "green widget", "red widget" };
 }
 // GET api/values/5
 [EnableCors("Policy1")]
 [HttpGet("{id}")]
 public ActionResult<string> Get(int id)
 {
 return id switch
 {
 1 => "green widget",
 2 => "red widget",
 _ => NotFound(),
 };
 }
}
```
The following code creates two CORS policies:

```
public class Startup
{
 public Startup(IConfiguration configuration)
 {
 Configuration = configuration;
 }
 public IConfiguration Configuration { get; }
 public void ConfigureServices(IServiceCollection services)
 {
 services.AddCors(options =>
 {
 options.AddPolicy("Policy1",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com");
 });
 options.AddPolicy("AnotherPolicy",
 policy =>
 {
 policy.WithOrigins("http://www.contoso.com")
 .AllowAnyHeader()
 .AllowAnyMethod();
 });
 });
 services.AddControllers();
 }
 public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
 {
 if (env.IsDevelopment())
 {
 app.UseDeveloperExceptionPage();
 }
 app.UseHttpsRedirection();
 app.UseRouting();
 app.UseCors();
 app.UseAuthorization();
 app.UseEndpoints(endpoints =>
 {
 endpoints.MapControllers();
 });
 }
}
```
For the finest control of limiting CORS requests:

- Use `[EnableCors("MyPolicy")]`with a named policy.
- Don't define a default policy.
- Don't use endpoint routing.

The code in the next section meets the preceding list.

See Test CORS for instructions on testing code similar to the preceding code.

### Disable CORS

The [DisableCors] attribute does **not** disable CORS that has been enabled by endpoint routing.

The following code defines the CORS policy `"MyPolicy"`:

```
public class Startup
{
 public void ConfigureServices(IServiceCollection services)
 {
 services.AddCors(options =>
 {
 options.AddPolicy(name: "MyPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com")
 .WithMethods("PUT", "DELETE", "GET");
 });
 });
 services.AddControllers();
 services.AddRazorPages();
 }
 public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
 {
 if (env.IsDevelopment())
 {
 app.UseDeveloperExceptionPage();
 }
 app.UseHttpsRedirection();
 app.UseStaticFiles();
 app.UseRouting();
 app.UseCors();
 app.UseAuthorization();
 app.UseEndpoints(endpoints =>
 {
 endpoints.MapControllers();
 endpoints.MapRazorPages();
 });
 }
}
```
The following code disables CORS for the `GetValues2` action:

```
[EnableCors("MyPolicy")]
[Route("api/[controller]")]
[ApiController]
public class ValuesController : ControllerBase
{
 // GET api/values
 [HttpGet]
 public IActionResult Get() =>
 ControllerContext.MyDisplayRouteInfo();
 // GET api/values/5
 [HttpGet("{id}")]
 public IActionResult Get(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // PUT api/values/5
 [HttpPut("{id}")]
 public IActionResult Put(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // GET: api/values/GetValues2
 [DisableCors]
 [HttpGet("{action}")]
 public IActionResult GetValues2() =>
 ControllerContext.MyDisplayRouteInfo();
}
```
The preceding code:

- Doesn't enable CORS with endpoint routing.
- Doesn't define a default CORS policy.
- Uses [EnableCors("MyPolicy")] to enable the `"MyPolicy"`CORS policy for the controller.
- Disables CORS for the `GetValues2`method.

See Test CORS for instructions on testing the preceding code.

## CORS policy options

This section describes the various options that can be set in a CORS policy:

- Set the allowed origins
- Set the allowed HTTP methods
- Set the allowed request headers
- Set the exposed response headers
- Credentials in cross-origin requests
- Set the preflight expiration time

AddPolicy is called in `Startup.ConfigureServices`. For some options, it may be helpful to read the How CORS works section first.

## Set the allowed origins

AllowAnyOrigin: Allows CORS requests from all origins with any scheme (`http` or `https`). `AllowAnyOrigin` is insecure because *any website* can make cross-origin requests to the app.

Note

Specifying `AllowAnyOrigin` and `AllowCredentials` is an insecure configuration and can result in cross-site request forgery. The CORS service returns an invalid CORS response when an app is configured with both methods.

`AllowAnyOrigin` affects preflight requests and the `Access-Control-Allow-Origin` header. For more information, see the Preflight requests section.

SetIsOriginAllowedToAllowWildcardSubdomains: Sets the IsOriginAllowed property of the policy to be a function that allows origins to match a configured wildcard domain when evaluating if the origin is allowed.

```
options.AddPolicy("MyAllowSubdomainPolicy",
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .SetIsOriginAllowedToAllowWildcardSubdomains();
 });
```
In the preceding code, `SetIsOriginAllowedToAllowWildcardSubdomains` is called with the wildcard origin `"https://*.example.com"`. This configuration allows CORS requests from any subdomain of `example.com`, such as `https://subdomain.example.com` or `https://api.example.com`. The `*` wildcard character must be included in the origin to enable wildcard subdomain matching.

### Set the allowed HTTP methods

- Allows any HTTP method:
- Affects preflight requests and the `Access-Control-Allow-Methods`header. For more information, see the Preflight requests section.

### Set the allowed request headers

To allow specific headers to be sent in a CORS request, called author request headers, call WithHeaders and specify the allowed headers:

```
options.AddPolicy("MyAllowHeadersPolicy",
 policy =>
 {
 // requires using Microsoft.Net.Http.Headers;
 policy.WithOrigins("http://example.com")
 .WithHeaders(HeaderNames.ContentType, "x-custom-header");
 });
```
To allow all author request headers, call AllowAnyHeader:

```
options.AddPolicy("MyAllowAllHeadersPolicy",
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .AllowAnyHeader();
 });
```
`AllowAnyHeader` affects preflight requests and the Access-Control-Request-Headers header. For more information, see the Preflight requests section.

A CORS middleware policy match to specific headers specified by `WithHeaders` is only possible when the headers sent in `Access-Control-Request-Headers` exactly match the headers stated in `WithHeaders`.

For instance, consider an app configured as follows:

```
app.UseCors(policy => policy.WithHeaders(HeaderNames.CacheControl));
```
CORS middleware declines a preflight request with the following request header because `Content-Language` (HeaderNames.ContentLanguage) isn't listed in `WithHeaders`:

```
Access-Control-Request-Headers: Cache-Control, Content-Language
```
The app returns a *200 OK* response but doesn't send the CORS headers back. Therefore, the browser doesn't attempt the cross-origin request.

### Set the exposed response headers

By default, the browser doesn't expose all of the response headers to the app. For more information, see W3C Cross-Origin Resource Sharing (Terminology): Simple Response Header.

The response headers that are available by default are:

- `Cache-Control`
- `Content-Language`
- `Content-Type`
- `Expires`
- `Last-Modified`
- `Pragma`

The CORS specification calls these headers *simple response headers*. To make other headers available to the app, call WithExposedHeaders:

```
options.AddPolicy("MyExposeResponseHeadersPolicy",
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .WithExposedHeaders("x-custom-header");
 });
```
### Credentials in cross-origin requests

Credentials require special handling in a CORS request. By default, the browser doesn't send credentials with a cross-origin request. Credentials include cookies and HTTP authentication schemes. To send credentials with a cross-origin request, the client must set `XMLHttpRequest.withCredentials` to `true`.

Using `XMLHttpRequest` directly:

```
var xhr = new XMLHttpRequest();
xhr.open('get', 'https://www.example.com/api/test');
xhr.withCredentials = true;
```
Using jQuery:

```
$.ajax({
 type: 'get',
 url: 'https://www.example.com/api/test',
 xhrFields: {
 withCredentials: true
 }
});
```
Using the Fetch API:

```
fetch('https://www.example.com/api/test', {
 credentials: 'include'
});
```
The server must allow the credentials. To allow cross-origin credentials, call AllowCredentials:

```
options.AddPolicy("MyMyAllowCredentialsPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .AllowCredentials();
 });
```
The HTTP response includes an `Access-Control-Allow-Credentials` header, which tells the browser that the server allows credentials for a cross-origin request.

If the browser sends credentials but the response doesn't include a valid `Access-Control-Allow-Credentials` header, the browser doesn't expose the response to the app, and the cross-origin request fails.

Allowing cross-origin credentials is a security risk. A website at another domain can send a signed-in user's credentials to the app on the user's behalf without the user's knowledge.

The CORS specification also states that setting origins to `"*"` (all origins) is invalid if the `Access-Control-Allow-Credentials` header is present.

## Preflight requests

For some CORS requests, the browser sends an additional OPTIONS request before making the actual request. This request is called a preflight request. The browser can skip the preflight request if **all** the following conditions are true:

- The request method is GET, HEAD, or POST.
- The app doesn't set request headers other than `Accept`,`Accept-Language`,`Content-Language`,`Content-Type`, or`Last-Event-ID`.
- The `Content-Type`header, if set, has one of the following values:- `application/x-www-form-urlencoded`
- `multipart/form-data`
- `text/plain`

The rule on request headers set for the client request applies to headers that the app sets by calling `setRequestHeader` on the `XMLHttpRequest` object. The CORS specification calls these headers author request headers. The rule doesn't apply to headers the browser can set, such as `User-Agent`, `Host`, or `Content-Length`.

The following is an example response similar to the preflight request made from the **[Put test]** button in the Test CORS section of this document.

```
General:
Request URL: https://cors3.azurewebsites.net/api/values/5
Request Method: OPTIONS
Status Code: 204 No Content
Response Headers:
Access-Control-Allow-Methods: PUT,DELETE,GET
Access-Control-Allow-Origin: https://cors1.azurewebsites.net
Server: Microsoft-IIS/10.0
Set-Cookie: ARRAffinity=8f8...8;Path=/;HttpOnly;Domain=cors1.azurewebsites.net
Vary: Origin
Request Headers:
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9
Access-Control-Request-Method: PUT
Connection: keep-alive
Host: cors3.azurewebsites.net
Origin: https://cors1.azurewebsites.net
Referer: https://cors1.azurewebsites.net/
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
User-Agent: Mozilla/5.0
```
The preflight request uses the HTTP OPTIONS method. It may include the following headers:

- Access-Control-Request-Method: The HTTP method that will be used for the actual request.
- Access-Control-Request-Headers: A list of request headers that the app sets on the actual request. As stated earlier, this doesn't include headers that the browser sets, such as `User-Agent`.
- Access-Control-Allow-Methods

If the preflight request is denied, the app returns a `200 OK` response but doesn't set the CORS headers. Therefore, the browser doesn't attempt the cross-origin request. For an example of a denied preflight request, see the Test CORS section of this document.

Using the F12 tools, the console app shows an error similar to one of the following, depending on the browser:

- Firefox: Cross-Origin Request Blocked: The Same Origin Policy disallows reading the remote resource at `https://cors1.azurewebsites.net/api/TodoItems1/MyDelete2/5`. (Reason: CORS request did not succeed). Learn More
- Chromium based: Access to fetch at 'https://cors1.azurewebsites.net/api/TodoItems1/MyDelete2/5' from origin 'https://cors3.azurewebsites.net' has been blocked by CORS policy: Response to preflight request doesn't pass access control check: No 'Access-Control-Allow-Origin' header is present on the requested resource. If an opaque response serves your needs, set the request's mode to 'no-cors' to fetch the resource with CORS disabled.

To allow specific headers, call WithHeaders:

```
options.AddPolicy("MyAllowHeadersPolicy",
 policy =>
 {
 // requires using Microsoft.Net.Http.Headers;
 policy.WithOrigins("http://example.com")
 .WithHeaders(HeaderNames.ContentType, "x-custom-header");
 });
```
To allow all author request headers, call AllowAnyHeader:

```
options.AddPolicy("MyAllowAllHeadersPolicy",
 policy =>
 {
 policy.WithOrigins("https://*.example.com")
 .AllowAnyHeader();
 });
```
Browsers aren't consistent in how they set `Access-Control-Request-Headers`. If either:

- Headers are set to anything other than `"*"`
- AllowAnyHeader is called:
Include at least `Accept`,`Content-Type`, and`Origin`, plus any custom headers that you want to support.

### Automatic preflight request code

When the CORS policy is applied either:

- Globally by calling `app.UseCors`in`Startup.Configure`.
- Using the `[EnableCors]`attribute.

ASP.NET Core responds to the preflight OPTIONS request.

Enabling CORS on a per-endpoint basis using `RequireCors` currently does **not** support automatic preflight requests.

The Test CORS section of this document demonstrates this behavior.

### [HttpOptions] attribute for preflight requests

When CORS is enabled with the appropriate policy, ASP.NET Core generally responds to CORS preflight requests automatically. In some scenarios, this may not be the case. For example, using CORS with endpoint routing.

The following code uses the [HttpOptions] attribute to create endpoints for OPTIONS requests:

```
[Route("api/[controller]")]
[ApiController]
public class TodoItems2Controller : ControllerBase
{
 // OPTIONS: api/TodoItems2/5
 [HttpOptions("{id}")]
 public IActionResult PreflightRoute(int id)
 {
 return NoContent();
 }
 // OPTIONS: api/TodoItems2
 [HttpOptions]
 public IActionResult PreflightRoute()
 {
 return NoContent();
 }
 [HttpPut("{id}")]
 public IActionResult PutTodoItem(int id)
 {
 if (id < 1)
 {
 return BadRequest();
 }
 return ControllerContext.MyDisplayRouteInfo(id);
 }
```
See Test CORS with endpoint routing and [HttpOptions] for instructions on testing the preceding code.

### Set the preflight expiration time

The `Access-Control-Max-Age` header specifies how long the response to the preflight request can be cached. To set this header, call SetPreflightMaxAge:

```
options.AddPolicy("MySetPreflightExpirationPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com")
 .SetPreflightMaxAge(TimeSpan.FromSeconds(2520));
 });
```
## How CORS works

This section describes what happens in a CORS request at the level of the HTTP messages.

- CORS is **not**a security feature. CORS is a W3C standard that allows a server to relax the same-origin policy.- For example, a malicious actor could use Cross-Site Scripting (XSS) against your site and execute a cross-site request to their CORS enabled site to steal information.

- An API isn't safer by allowing CORS.
- It's up to the client (browser) to enforce CORS. The server executes the request and returns the response, it's the client that returns an error and blocks the response. For example, any of the following tools will display the server response:
- Fiddler
- .NET HttpClient
- A web browser by entering the URL in the address bar.

- It's up to the client (browser) to enforce CORS. The server executes the request and returns the response, it's the client that returns an error and blocks the response. For example, any of the following tools will display the server response:
- It's a way for a server to allow browsers to execute a cross-origin XHR or Fetch API request that otherwise would be forbidden.
- Browsers without CORS can't do cross-origin requests. Before CORS, JSONP was used to circumvent this restriction. JSONP doesn't use XHR, it uses the `<script>`tag to receive the response. Scripts are allowed to be loaded cross-origin.

- Browsers without CORS can't do cross-origin requests. Before CORS, JSONP was used to circumvent this restriction. JSONP doesn't use XHR, it uses the

The CORS specification introduced several new HTTP headers that enable cross-origin requests. If a browser supports CORS, it sets these headers automatically for cross-origin requests. Custom JavaScript code isn't required to enable CORS.

The following is an example of a cross-origin request from the **Values** test button to `https://cors1.azurewebsites.net/api/values`. The `Origin` header:

- Provides the domain of the site that's making the request.
- Is required and must be different from the host.

**General headers**

```
Request URL: https://cors1.azurewebsites.net/api/values
Request Method: GET
Status Code: 200 OK
```
**Response headers**

```
Content-Encoding: gzip
Content-Type: text/plain; charset=utf-8
Server: Microsoft-IIS/10.0
Set-Cookie: ARRAffinity=8f...;Path=/;HttpOnly;Domain=cors1.azurewebsites.net
Transfer-Encoding: chunked
Vary: Accept-Encoding
X-Powered-By: ASP.NET
```
**Request headers**

```
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9
Connection: keep-alive
Host: cors1.azurewebsites.net
Origin: https://cors3.azurewebsites.net
Referer: https://cors3.azurewebsites.net/
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
User-Agent: Mozilla/5.0 ...
```
In `OPTIONS` requests, the server sets the **Response headers** `Access-Control-Allow-Origin: {allowed origin}` header in the response. For example, the deployed sample, Delete button `OPTIONS` request contains the following headers:

**General headers**

```
Request URL: https://cors3.azurewebsites.net/api/TodoItems2/MyDelete2/5
Request Method: OPTIONS
Status Code: 204 No Content
```
**Response headers**

```
Access-Control-Allow-Headers: Content-Type,x-custom-header
Access-Control-Allow-Methods: PUT,DELETE,GET,OPTIONS
Access-Control-Allow-Origin: https://cors1.azurewebsites.net
Server: Microsoft-IIS/10.0
Set-Cookie: ARRAffinity=8f...;Path=/;HttpOnly;Domain=cors3.azurewebsites.net
Vary: Origin
X-Powered-By: ASP.NET
```
**Request headers**

```
Accept: */*
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9
Access-Control-Request-Headers: content-type
Access-Control-Request-Method: DELETE
Connection: keep-alive
Host: cors3.azurewebsites.net
Origin: https://cors1.azurewebsites.net
Referer: https://cors1.azurewebsites.net/test?number=2
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
User-Agent: Mozilla/5.0
```
In the preceding **Response headers**, the server sets the Access-Control-Allow-Origin header in the response. The `https://cors1.azurewebsites.net` value of this header matches the `Origin` header from the request.

If AllowAnyOrigin is called, the `Access-Control-Allow-Origin: *`, the wildcard value, is returned. `AllowAnyOrigin` allows any origin.

If the response doesn't include the `Access-Control-Allow-Origin` header, the cross-origin request fails. Specifically, the browser disallows the request. Even if the server returns a successful response, the browser doesn't make the response available to the client app.

### Display OPTIONS requests

By default, the Chrome and Edge browsers don't show OPTIONS requests on the network tab of the F12 tools. To display OPTIONS requests in these browsers:

- `chrome://flags/#out-of-blink-cors`or- `edge://flags/#out-of-blink-cors`
- disable the flag.
- restart.

Firefox shows OPTIONS requests by default.

## CORS in IIS

When deploying to IIS, CORS has to run before Windows Authentication if the server isn't configured to allow anonymous access. To support this scenario, the IIS CORS module needs to be installed and configured for the app.

## Test CORS

The sample download has code to test CORS. See how to download. The sample is an API project with Razor Pages added:

```
public class StartupTest2
{
 public void ConfigureServices(IServiceCollection services)
 {
 services.AddCors(options =>
 {
 options.AddPolicy(name: "MyPolicy",
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com",
 "https://cors1.azurewebsites.net",
 "https://cors3.azurewebsites.net",
 "https://localhost:44398",
 "https://localhost:5001")
 .WithMethods("PUT", "DELETE", "GET");
 });
 });
 services.AddControllers();
 services.AddRazorPages();
 }
 public void Configure(IApplicationBuilder app)
 {
 app.UseHttpsRedirection();
 app.UseStaticFiles();
 app.UseRouting();
 app.UseCors();
 app.UseAuthorization();
 app.UseEndpoints(endpoints =>
 {
 endpoints.MapControllers();
 endpoints.MapRazorPages();
 });
 }
}
```
Warning

`WithOrigins("https://localhost:<port>");` should only be used for testing a sample app similar to the download sample code.

The following `ValuesController` provides the endpoints for testing:

```
[EnableCors("MyPolicy")]
[Route("api/[controller]")]
[ApiController]
public class ValuesController : ControllerBase
{
 // GET api/values
 [HttpGet]
 public IActionResult Get() =>
 ControllerContext.MyDisplayRouteInfo();
 // GET api/values/5
 [HttpGet("{id}")]
 public IActionResult Get(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // PUT api/values/5
 [HttpPut("{id}")]
 public IActionResult Put(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // GET: api/values/GetValues2
 [DisableCors]
 [HttpGet("{action}")]
 public IActionResult GetValues2() =>
 ControllerContext.MyDisplayRouteInfo();
}
```
MyDisplayRouteInfo is provided by the Rick.Docs.Samples.RouteInfo NuGet package and displays route information.

Test the preceding sample code by using one of the following approaches:

- Run the sample with `dotnet run`using the default URL of`https://localhost:5001`.
- Run the sample from Visual Studio with the port set to 44398 for a URL of `https://localhost:44398`.

Using a browser with the F12 tools:

- Select the - **Values**button and review the headers in the- **Network**tab.
- Select the - **PUT test**button. See Display OPTIONS requests for instructions on displaying the OPTIONS request. The- **PUT test**creates two requests, an OPTIONS preflight request and the PUT request.
- Select the - `GetValues2 [DisableCors]`- **Console**tab to see the CORS error. Depending on the browser, an error similar to the following is displayed:- Access to fetch at - `'https://cors1.azurewebsites.net/api/values/GetValues2'`from origin- `'https://cors3.azurewebsites.net'`has been blocked by CORS policy: No 'Access-Control-Allow-Origin' header is present on the requested resource. If an opaque response serves your needs, set the request's mode to 'no-cors' to fetch the resource with CORS disabled.

CORS-enabled endpoints can be tested with a tool, such as curl or Fiddler. When using a tool, the origin of the request specified by the `Origin` header must differ from the host receiving the request. If the request isn't *cross-origin* based on the value of the `Origin` header:

- There's no need for CORS middleware to process the request.
- CORS headers aren't returned in the response.

The following command uses `curl` to issue an OPTIONS request with information:

```
curl -X OPTIONS https://cors3.azurewebsites.net/api/TodoItems2/5 -i
```
### Test CORS with endpoint routing and [HttpOptions]

Enabling CORS on a per-endpoint basis using `RequireCors` currently does **not** support automatic preflight requests. Consider the following code which uses endpoint routing to enable CORS:

```
public class StartupEndPointBugTest
{
 readonly string MyPolicy = "_myPolicy";
 // .WithHeaders(HeaderNames.ContentType, "x-custom-header")
 // forces browsers to require a preflight request with GET
 public void ConfigureServices(IServiceCollection services)
 {
 services.AddCors(options =>
 {
 options.AddPolicy(name: MyPolicy,
 policy =>
 {
 policy.WithOrigins("http://example.com",
 "http://www.contoso.com",
 "https://cors1.azurewebsites.net",
 "https://cors3.azurewebsites.net",
 "https://localhost:44398",
 "https://localhost:5001")
 .WithHeaders(HeaderNames.ContentType, "x-custom-header")
 .WithMethods("PUT", "DELETE", "GET", "OPTIONS");
 });
 });
 services.AddControllers();
 services.AddRazorPages();
 }
 public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
 {
 app.UseHttpsRedirection();
 app.UseStaticFiles();
 app.UseRouting();
 app.UseCors();
 app.UseAuthorization();
 app.UseEndpoints(endpoints =>
 {
 endpoints.MapControllers().RequireCors(MyPolicy);
 endpoints.MapRazorPages();
 });
 }
}
```
The following `TodoItems1Controller` provides endpoints for testing:

```
[Route("api/[controller]")]
[ApiController]
public class TodoItems1Controller : ControllerBase
{
 // PUT: api/TodoItems1/5
 [HttpPut("{id}")]
 public IActionResult PutTodoItem(int id)
 {
 if (id < 1)
 {
 return Content($"ID = {id}");
 }
 return ControllerContext.MyDisplayRouteInfo(id);
 }
 // Delete: api/TodoItems1/5
 [HttpDelete("{id}")]
 public IActionResult MyDelete(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 // GET: api/TodoItems1
 [HttpGet]
 public IActionResult GetTodoItems() =>
 ControllerContext.MyDisplayRouteInfo();
 [EnableCors]
 [HttpGet("{action}")]
 public IActionResult GetTodoItems2() =>
 ControllerContext.MyDisplayRouteInfo();
 // Delete: api/TodoItems1/MyDelete2/5
 [EnableCors]
 [HttpDelete("{action}/{id}")]
 public IActionResult MyDelete2(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
}
```
Test the preceding code from the test page (`https://cors1.azurewebsites.net/test?number=1`) of the deployed sample.

The **Delete [EnableCors]** and **GET [EnableCors]** buttons succeed, because the endpoints have `[EnableCors]` and respond to preflight requests. The other endpoints fails. The **GET** button fails, because the JavaScript sends:

```
 headers: {
 "Content-Type": "x-custom-header"
 },
```
The following `TodoItems2Controller` provides similar endpoints, but includes explicit code to respond to OPTIONS requests:

```
[Route("api/[controller]")]
[ApiController]
public class TodoItems2Controller : ControllerBase
{
 // OPTIONS: api/TodoItems2/5
 [HttpOptions("{id}")]
 public IActionResult PreflightRoute(int id)
 {
 return NoContent();
 }
 // OPTIONS: api/TodoItems2
 [HttpOptions]
 public IActionResult PreflightRoute()
 {
 return NoContent();
 }
 [HttpPut("{id}")]
 public IActionResult PutTodoItem(int id)
 {
 if (id < 1)
 {
 return BadRequest();
 }
 return ControllerContext.MyDisplayRouteInfo(id);
 }
 // [EnableCors] // Not needed as OPTIONS path provided
 [HttpDelete("{id}")]
 public IActionResult MyDelete(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
 [EnableCors] // Rquired for this path
 [HttpGet]
 public IActionResult GetTodoItems() =>
 ControllerContext.MyDisplayRouteInfo();
 [HttpGet("{action}")]
 public IActionResult GetTodoItems2() =>
 ControllerContext.MyDisplayRouteInfo();
 [EnableCors] // Rquired for this path
 [HttpDelete("{action}/{id}")]
 public IActionResult MyDelete2(int id) =>
 ControllerContext.MyDisplayRouteInfo(id);
}
```
The preceding code can be tested by deploying the sample to Azure.In the **Controller** drop down list, select **Preflight** and then **Set Controller**. All the CORS calls to the `TodoItems2Controller` endpoints succeed.

## Additional resources

ASP.NET Core

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Cross-Site Scripting (XSS) is a security vulnerability that enables a cyberattacker to place client side scripts (usually JavaScript) into web pages. When other users load affected pages, the cyberattacker's scripts run. The cyberattacker can then steal cookies and session tokens, change the contents of the web page through DOM manipulation, or redirect the browser to another page. XSS vulnerabilities generally occur when an application takes user input and outputs it to a page without validating, encoding, or escaping it.

This article applies primarily to ASP.NET Core MVC with views, Razor Pages, and other apps that return HTML that can be vulnerable to XSS. Web APIs that return data in the form of HTML, XML, or JSON can trigger XSS attacks in their client apps if they don't properly sanitize user input. This behavior depends on how much trust the client app places in the API. If an API accepts user-generated content and returns it in an HTML response, the data is open to attack. A cyberattacker can inject malicious scripts into the content that executes when the response is rendered in the user's browser.

To prevent XSS attacks, web APIs should implement input validation and output encoding. Input validation ensures that user input meets expected criteria and doesn't include malicious code. Output encoding ensures that any data returned by the API is properly sanitized so it can't be executed as code by the user's browser. For more information, see GitHub dotnet/aspnetcore.docs issue #28789.

## Protect your application against XSS

At a basic level, XSS works by tricking your application into inserting a `<script>` tag into your rendered page, or by inserting an `On*` event into an element.

To avoid introducing XSS into the application, developers should implement the following prevention techniques:

- Never put untrusted data into your HTML input, unless you follow the other techniques listed in this section. - Untrusted data is any data controllable by a cyberattacker. Examples include HTML form inputs, query strings, HTTP headers, or even data sourced from a database. A cyberattacker might be able to breach your database even if they can't breach your application.
- Before you put untrusted data into an HTML element, ensure the data is HTML encoded. - HTML encoding takes characters such as the left angle bracket or - *less than*(- `<`) and changes them into a safe form like (- `<`).
- Before you put untrusted data into an HTML attribute, ensure the data is HTML-attribute encoded. - This specialized form of HTML encoding handles double quote ( - `"`), single quote (- `'`), ampersand (- `&`), and less than (- `<`) characters. When dealing with untrusted input, use HTML encoding for general HTML content and HTML attribute encoding for HTML attributes.
- Before you put untrusted data into JavaScript, place the data in an HTML element whose contents you retrieve at runtime. - If you can't follow this technique, ensure the data is JavaScript encoded. JavaScript encoding converts dangerous characters for JavaScript into a hexadecimal equivalent value. For example, JavaScript encoding changes the less than ( - `<`) character into the hex value- `\u003C`.
- Before you put untrusted data into a URL query string, ensure the data is URL encoded.

## Explore HTML encoding with Razor

The Razor engine used in MVC automatically encodes all output sourced from variables, unless you work to prevent this behavior. It uses HTML attribute encoding rules whenever you use the at symbol `@` directive. Because HTML attribute encoding is a superset of HTML encoding, you don't have to consider whether to use HTML encoding or HTML-attribute encoding. You must ensure that you only use the at symbol `@` in an HTML context, and not when attempting to insert untrusted input directly into JavaScript. Razor Tag Helpers also encode input you use in tag parameters.

Consider the following Razor view:

```
@{
 var untrustedInput = "<\"123\">";
}
@untrustedInput
```
This view outputs the contents of the `untrustedInput` variable. The variable includes some characters used in XSS attacks: less than (`<`), double quote (`"`), and right angle bracket or *greater than* (`>`). Examining the source shows the rendered output encoded as:

```
<"123">
```
Warning

ASP.NET Core MVC provides an `HtmlString` class that isn't automatically encoded upon output. This class should never be used in combination with untrusted input because it exposes an XSS vulnerability.

## Explore JavaScript encoding with Razor

In some instances, you might want to insert a value into JavaScript to process in your view. There are two ways to accomplish this task. The safest way to insert values is to place the value in a data attribute of a tag and retrieve it in your JavaScript. For example:

```
@{
 var untrustedInput = "<script>alert(1)</script>";
}
<div id="injectedData"
 data-untrustedinput="@untrustedInput" />
<div id="scriptedWrite" />
<div id="scriptedWrite-html5" />
<script>
 var injectedData = document.getElementById("injectedData");
 // All clients
 var clientSideUntrustedInputOldStyle =
 injectedData.getAttribute("data-untrustedinput");
 // HTML 5 clients only
 var clientSideUntrustedInputHtml5 =
 injectedData.dataset.untrustedinput;
 // Put the injected, untrusted data into the scriptedWrite div tag.
 // Do NOT use document.write() on dynamically generated data as it can lead to XSS.
 document.getElementById("scriptedWrite").innerText += clientSideUntrustedInputOldStyle;
 // Or, you can use createElement() to dynamically create document elements.
 // This instance uses textContent to ensure the data is properly encoded.
 var x = document.createElement("div");
 x.textContent = clientSideUntrustedInputHtml5;
 document.body.appendChild(x);
 // You can also use createTextNode on an element to ensure data is properly encoded.
 var y = document.createElement("div");
 y.appendChild(document.createTextNode(clientSideUntrustedInputHtml5));
 document.body.appendChild(y);
</script>
```
The preceding markup generates the following HTML:

```
<div id="injectedData"
 data-untrustedinput="<script>alert(1)</script>" />
<div id="scriptedWrite" />
<div id="scriptedWrite-html5" />
<script>
 var injectedData = document.getElementById("injectedData");
 // All clients
 var clientSideUntrustedInputOldStyle =
 injectedData.getAttribute("data-untrustedinput");
 // HTML 5 clients only
 var clientSideUntrustedInputHtml5 =
 injectedData.dataset.untrustedinput;
 // Put the injected, untrusted data into the scriptedWrite div tag.
 // Do NOT use document.write() on dynamically generated data as it can lead to XSS.
 document.getElementById("scriptedWrite").innerText += clientSideUntrustedInputOldStyle;
 // Or, you can use createElement() to dynamically create document elements.
 // This instance uses textContent to ensure the data is properly encoded.
 var x = document.createElement("div");
 x.textContent = clientSideUntrustedInputHtml5;
 document.body.appendChild(x);
 // You can also use createTextNode on an element to ensure data is properly encoded.
 var y = document.createElement("div");
 y.appendChild(document.createTextNode(clientSideUntrustedInputHtml5));
 document.body.appendChild(y);
</script>
```
The preceding code generates the following output:

```
<script>alert(1)</script>
<script>alert(1)</script>
<script>alert(1)</script>
```
Warning

Don't concatenate untrusted input in JavaScript to create DOM elements or use `document.write()` on dynamically generated content.

Instead, use one of the following approaches to prevent code from being exposed to DOM-based XSS:

- Call `createElement()`and assign property values with appropriate methods or properties, such as`node.textContent=`or`node.InnerText=`.
- Call the `document.CreateTextNode()`method and append it in the appropriate DOM location.
- Call the `element.SetAttribute()`method.
- Use the `element[attribute]=`assignment.

## Access encoders in code

You can use HTML, JavaScript, and URL encoders in your code in two ways:

- Inject them via dependency injection.
- Use the default encoders contained in the `System.Text.Encodings.Web`namespace.

When you use the default encoders, any customizations applied to character ranges (so they're treated as safe) don't take effect. The default encoders use the safest encoding rules possible.

To use the configurable encoders via dependency injection, your constructors should take an `HtmlEncoder`, `JavaScriptEncoder`, and `UrlEncoder` parameter, as appropriate.

For example:

```
public class HomeController : Controller
{
 HtmlEncoder _htmlEncoder;
 JavaScriptEncoder _javaScriptEncoder;
 UrlEncoder _urlEncoder;
 public HomeController(HtmlEncoder htmlEncoder,
 JavaScriptEncoder javascriptEncoder,
 UrlEncoder urlEncoder)
 {
 _htmlEncoder = htmlEncoder;
 _javaScriptEncoder = javascriptEncoder;
 _urlEncoder = urlEncoder;
 }
}
```
## Encode URL parameters

If you want to build a URL query string with untrusted input as a value, use the `UrlEncoder` parameter to encode the value:

```
var example = "\"Quoted Value with spaces and &\"";
var encodedValue = _urlEncoder.Encode(example);
```
After encoding, the `encodedValue` variable contains the string `%22Quoted%20Value%20with%20spaces%20and%20%26%22`. Spaces, quotes, punctuation, and other unsafe characters are percent encoded to their hexadecimal value. For example, a space character is converted to `%20`.

Warning

Don't use untrusted input as part of a URL path. Always pass untrusted input as a query string value.

## Customize the encoders

By default, encoders use a safe list limited to the Basic Latin Unicode range. All characters outside of the indicated range are encoded as their character code equivalents. This behavior also affects rendering by Razor Tag Helpers and HTML Helpers because they use the encoders to output your strings.

The purpose for this behavior is to protect against unknown or future browser bugs. Previous browser bugs tripped up parsing based on the processing of non-English characters. If your web site makes heavy use of non-Latin characters, such as Chinese, Cyrillic, or others, this behavior is probably not suited for your configuration.

You can customize the encoder safe lists to include Unicode ranges appropriate to the app during startup. Make the customizations in the *Program.cs* file.

For example, you can use the default configuration with a Razor HTML Helper similar to the following HTML:

```
<p>This link text is in Chinese: @Html.ActionLink("汉语/漢語", "Index")</p>
```
The preceding markup is rendered with Chinese text encoded:

```
<p>This link text is in Chinese: <a href="/">汉语/漢語</a></p>
```
To widen the range of characters treated as safe by the encoder, insert the following line into the *Program.cs* file:

```
builder.Services.AddSingleton<HtmlEncoder>(
 HtmlEncoder.Create(allowedRanges: new[] { UnicodeRanges.BasicLatin,
 UnicodeRanges.CjkUnifiedIdeographs }));
```
You can customize the encoder safe lists to include Unicode ranges appropriate to your application during startup, in `ConfigureServices()`.

For example, using the default configuration you might use a Razor HtmlHelper like so;

```
<p>This link text is in Chinese: @Html.ActionLink("汉语/漢語", "Index")</p>
```
When you view the source of the web page you'll see it has been rendered as follows, with the Chinese text encoded;

```
<p>This link text is in Chinese: <a href="/">汉语/漢語</a></p>
```
To widen the characters treated as safe by the encoder you would insert the following line into the `ConfigureServices()` method in `startup.cs`;

```
services.AddSingleton<HtmlEncoder>(
 HtmlEncoder.Create(allowedRanges: new[] { UnicodeRanges.BasicLatin,
 UnicodeRanges.CjkUnifiedIdeographs }));
```
This example widens the safe list to include the Unicode Range CJK Unified Ideographs. The following output shows the rendered view for the wider range of safe characters:

```
<p>This link text is in Chinese: <a href="/">汉语/漢語</a></p>
```
Safe list ranges are specified as Unicode code charts, not languages. The Unicode standard has a list of code charts you can use to find the chart that contains your characters. Each encoder (HTML, JavaScript, URL) must be configured separately.

Note

Customization of the safe list only affects encoders sourced via dependency injection.
If you directly access an encoder via `System.Text.Encodings.Web.*Encoder.Default`, only the default safe list is used, Basic Latin.

## Determine when and where to encode

In general, the accepted practice is that encoding takes place at the point of output and encoded values should never be stored in a database.

Encoding at the point of output allows you to change the use of data. For example, change from HTML to a query string value. This approach enables you to easily search your data without having to encode values before searching. It also allows you to take advantage of any changes or bug fixes made to encoders.

## Use validation as an XSS prevention technique

Validation can be a useful tool in limiting XSS attacks. For example, a numeric string containing only the characters 0-9 doesn't trigger an XSS attack.

Validation is more complicated when HTML is accepted in user input. Parsing HTML input can be difficult, and sometimes, impossible. Markdown, coupled with a parser that strips embedded HTML, is a safer option for accepting rich input.

Never rely on validation alone. Always encode untrusted input before output, no matter what validation or sanitization is performed.

## Related content

ASP.NET Core

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

A web app that redirects to a URL that's specified via the request such as the querystring or form data can potentially be tampered with to redirect users to an external, malicious URL. This tampering is called an open redirection attack.

Whenever your application logic redirects to a specified URL, you must verify that the redirection URL hasn't been tampered with. ASP.NET Core has built-in functionality to help protect apps from open redirect (also known as open redirection) attacks.

## What is an open redirect attack?

Web applications frequently redirect users to a login page when they access resources that require authentication. The redirection typically includes a `returnUrl` querystring parameter so that the user can be returned to the originally requested URL after they have successfully logged in. After the user authenticates, they're redirected to the URL they had originally requested.

Because the destination URL is specified in the querystring of the request, a malicious user could tamper with the querystring. A tampered querystring could allow the site to redirect the user to an external, malicious site. This technique is called an open redirect (or redirection) attack.

### An example attack

A malicious user can develop an attack intended to allow the malicious user access to a user's credentials or sensitive information. To begin the attack, the malicious user convinces the user to click a link to your site's login page with a `returnUrl` querystring value added to the URL. For example, consider an app at `contoso.com` that includes a login page at `http://contoso.com/Account/LogOn?returnUrl=/Home/About`. The attack follows these steps:

- The user clicks a malicious link to `http://contoso.com/Account/LogOn?returnUrl=http://contoso1.com/Account/LogOn`(the second URL is "contoso**1**.com", not "contoso.com").
- The user logs in successfully.
- The user is redirected (by the site) to `http://contoso1.com/Account/LogOn`(a malicious site that looks exactly like real site).
- The user logs in again (giving malicious site their credentials) and is redirected back to the real site.

The user likely believes that their first attempt to log in failed and that their second attempt is successful. The user most likely remains unaware that their credentials are compromised.

In addition to login pages, some sites provide redirect pages or endpoints. Imagine your app has a page with an open redirect, `/Home/Redirect`. A cyberattacker could create, for example, a link in an email that goes to `[yoursite]/Home/Redirect?url=http://phishingsite.com/Home/Login`. A typical user will look at the URL and see it begins with your site name. Trusting that, they will click the link. The open redirect would then send the user to the phishing site, which looks identical to yours, and the user would likely login to what they believe is your site.

## Protecting against open redirect attacks

When developing web applications, treat all user-provided data as untrustworthy. If your application has functionality that redirects the user based on the contents of the URL, ensure that such redirects are only done locally within your app (or to a known URL, not any URL that may be supplied in the querystring).

### LocalRedirect

Use the `LocalRedirect` helper method from the base `Controller` class:

```
public IActionResult SomeAction(string redirectUrl)
{
 return LocalRedirect(redirectUrl);
}
```
`LocalRedirect` will throw an exception if a non-local URL is specified. Otherwise, it behaves just like the `Redirect` method.

### IUrlHelper.IsLocalUrl

Use the IsLocalUrl method to test URLs before redirecting:

The following example shows how to check whether a URL is local before redirecting.

```
private IActionResult RedirectToLocal(string returnUrl)
{
 if (Url.IsLocalUrl(returnUrl))
 {
 return Redirect(returnUrl);
 }
 else
 {
 return RedirectToAction(nameof(HomeController.Index), "Home");
 }
}
```
The `IUrlHelper.IsLocalUrl` method protects users from being inadvertently redirected to a malicious site. You can log the details of the URL that was provided when a non-local URL is supplied in a situation where you expected a local URL. Logging redirect URLs may help in diagnosing redirection attacks.

### Detect if URL is local using `RedirectHttpResult.IsLocalUrl`

The `RedirectHttpResult.IsLocalUrl(url)` helper method detects if a URL is local. A URL is considered local if the following are true:

- It doesn't have the host or authority section.
- It has an absolute path.

URLs using virtual paths `"~/"` are also local.

`IsLocalUrl` is useful for validating URLs before redirecting to them to prevent open redirection attacks.

```
if (RedirectHttpResult.IsLocalUrl(url))
{
 return Results.LocalRedirect(url);
}
```
ASP.NET Core

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

ASP.NET Core Identity:

- Is an API that supports user interface (UI) login functionality.
- Manages users, passwords, profile data, roles, claims, tokens, email confirmation, and more.

Users can create an account with the login information stored in Identity or they can use an external login provider. Supported external login providers include Facebook, Google, Microsoft Account, and Twitter.

For information on how to require authentication for all app users, see Create an ASP.NET Core app with user data protected by authorization.

The Identity source code is available on GitHub. Scaffold Identity and view the generated files to review the template interaction with Identity.

Identity is typically configured using a SQL Server database to store user names, passwords, and profile data. Alternatively, another persistent store can be used, for example, Azure Table Storage.

In this topic, you learn how to use Identity to register, log in, and log out a user. Note: the templates treat username and email as the same for users. For more detailed instructions about creating apps that use Identity, see Next Steps.

For more information on Identity in Blazor apps, see ASP.NET Core Blazor authentication and authorization and the articles that follow it in the Blazor documentation.

ASP.NET Core Identity isn't related to the Microsoft identity platform. Microsoft identity platform is:

- An evolution of the Azure Active Directory (Azure AD) developer platform.
- An alternative identity solution for authentication and authorization in ASP.NET Core apps.

ASP.NET Core Identity adds user interface (UI) login functionality to ASP.NET Core web apps. To secure web APIs and SPAs, use one of the following:

Duende Identity Server is an OpenID Connect and OAuth 2.0 framework for ASP.NET Core. Duende Identity Server enables the following security features:

- Authentication as a Service (AaaS)
- Single sign-on/off (SSO) over multiple application types
- Access control for APIs
- Federation Gateway

Important

Duende Software might require you to pay a license fee for production use of Duende Identity Server. For more information, see Migrate from ASP.NET Core in .NET 5 to .NET 6.

For more information, see the Duende Identity Server documentation (Duende Software website).

View or download the sample code (how to download).

## Create a Blazor Web App with authentication

Create an ASP.NET Core Blazor Web App project with Individual Accounts.

Note

For a Razor Pages experience, see the Create a Razor Pages app with authentication section.

For an MVC experience, see the Create an MVC app with authentication section.

- Select the **Blazor Web App**template. Select**Next**.
- Make the following selections:
- **Authentication type**:- **Individual Accounts**
- **Interactive render mode**:- **Server**
- **Interactivity Location**:- **Global**

- Select **Create**.

The generated project includes Identity Razor components. The components are found in the `Components/Account` folder of the server project. For example:

- `Components/Account/Pages/Register.razor`
- `Components/Account/Pages/Login.razor`
- `Components/Account/Pages/Manage/ChangePassword.razor`

Identity Razor components are described individually in the documentation for specific use cases and are subject to change each release. When you generate a Blazor Web App with Individual Accounts, Identity Razor components are included in the generated project. The Identity Razor components can also be inspected in the `Components/Account` folder of the server project in the Blazor Web App project template (`dotnet/aspnetcore` GitHub repository).

Note

Documentation links to .NET reference source usually load the repository's default branch, which represents the current development for the next release of .NET. To select a tag for a specific release, use the **Switch branches or tags** dropdown list. For more information, see How to select a version tag of ASP.NET Core source code (dotnet/AspNetCore.Docs #26205).

For more information, see ASP.NET Core Blazor authentication and authorization and the articles that follow it in the Blazor documentation. Most of the articles in the *Security and Identity* area of the main ASP.NET Core documentation set apply to Blazor apps. However, the Blazor documentation set contains articles and guidance that supersedes or adds information. We recommend studying the general ASP.NET Core documentation set first, followed by accessing the articles in the Blazor *Security and Identity* documentation.

## Create a Razor Pages app with authentication

Create an ASP.NET Core Web Application (Razor Pages) project with Individual Accounts.

- Select the **ASP.NET Core Web App (Razor Pages)**template. Select**Next**.
- For **Authentication type**, select**Individual Accounts**.
- Select **Create**.

The generated project provides ASP.NET Core Identity as a Razor class library (RCL). The Identity Razor class library exposes endpoints with the `Identity` area. For example:

- `Areas/Identity/Pages/Account/Register`
- `Areas/Identity/Pages/Account/Login`
- `Areas/Identity/Pages/Account/Manage/ChangePassword`

Pages are described individually in the documentation for specific use cases and are subject to change each release. To view all of the pages in the RCL, see the ASP.NET Core reference source (`dotnet/aspnetcore` GitHub repository, `Identity/UI/src/Areas/Identity/Pages` folder). You can *scaffold* individual pages or all of the pages into the app. For more information, see Scaffold Identity in ASP.NET Core projects.

## Create an MVC app with authentication

Create an ASP.NET Core MVC project with Individual Accounts.

- Select the **ASP.NET Core Web App (Model-View-Controller)**template. Select**Next**.
- For **Authentication type**, select**Individual Accounts**.
- Select **Create**.

The generated project provides ASP.NET Core Identity as a Razor class library (RCL). The Identity Razor class library is based on Razor Pages and exposes endpoints with the `Identity` area. For example:

- `Areas/Identity/Pages/Account/Register`
- `Areas/Identity/Pages/Account/Login`
- `Areas/Identity/Pages/Account/Manage/ChangePassword`

Pages are described individually in the documentation for specific use cases and are subject to change each release. To view all of the pages in the RCL, see the ASP.NET Core reference source (`dotnet/aspnetcore` GitHub repository, `Identity/UI/src/Areas/Identity/Pages` folder). You can *scaffold* individual pages or all of the pages into the app. For more information, see Scaffold Identity in ASP.NET Core projects.

### Apply migrations

Apply the migrations to initialize the database.

Run the following command in the Package Manager Console (PMC):

`Update-Database`

### Test Register and Login

Run the app and register a user. Depending on your screen size, you might need to select the navigation toggle button to see the **Register** and **Login** links.

### View the Identity database

- From the **View**menu, select**SQL Server Object Explorer**(SSOX).
- Navigate to **(localdb)MSSQLLocalDB(SQL Server 13)**. Right-click on**dbo.AspNetUsers**>**View Data**:

### Configure Identity services

Services are added in `Program.cs`. The typical pattern is to call methods in the following order:

- `Add{Service}`
- `builder.Services.Configure{Service}`

```
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using WebApp1.Data;
var builder = WebApplication.CreateBuilder(args);
var connectionString = builder.Configuration.GetConnectionString("DefaultConnection");
builder.Services.AddDbContext<ApplicationDbContext>(options =>
 options.UseSqlServer(connectionString));
builder.Services.AddDatabaseDeveloperPageExceptionFilter();
builder.Services.AddDefaultIdentity<IdentityUser>(options => options.SignIn.RequireConfirmedAccount = true)
 .AddEntityFrameworkStores<ApplicationDbContext>();
builder.Services.AddRazorPages();
builder.Services.Configure<IdentityOptions>(options =>
{
 // Password settings.
 options.Password.RequireDigit = true;
 options.Password.RequireLowercase = true;
 options.Password.RequireNonAlphanumeric = true;
 options.Password.RequireUppercase = true;
 options.Password.RequiredLength = 6;
 options.Password.RequiredUniqueChars = 1;
 // Lockout settings.
 options.Lockout.DefaultLockoutTimeSpan = TimeSpan.FromMinutes(5);
 options.Lockout.MaxFailedAccessAttempts = 5;
 options.Lockout.AllowedForNewUsers = true;
 // User settings.
 options.User.AllowedUserNameCharacters =
 "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._@+";
 options.User.RequireUniqueEmail = false;
});
builder.Services.ConfigureApplicationCookie(options =>
{
 // Cookie settings
 options.Cookie.HttpOnly = true;
 options.ExpireTimeSpan = TimeSpan.FromMinutes(5);
 options.LoginPath = "/Identity/Account/Login";
 options.AccessDeniedPath = "/Identity/Account/AccessDenied";
 options.SlidingExpiration = true;
});
var app = builder.Build();
if (app.Environment.IsDevelopment())
{
 app.UseMigrationsEndPoint();
}
else
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthentication();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
The preceding code configures Identity with default option values. Services are made available to the app through dependency injection.

Identity is enabled by calling UseAuthentication. `UseAuthentication` adds authentication middleware to the request pipeline.

The template-generated app doesn't use authorization. `app.UseAuthorization` is included to ensure it's added in the correct order should the app add authorization. `UseRouting`, `UseAuthentication`, and `UseAuthorization` must be called in the order shown in the preceding code.

For more information on `IdentityOptions`, see IdentityOptions and Application Startup.

## ASP.NET Core Identity metrics

ASP.NET Core Identity metrics provide monitoring capabilities for user management and authentication processes. These metrics help you detect unusual sign-in patterns that might indicate security threats, track the performance of identity operations, and understand how users interact with authentication features, such as two-factor authentication. This observability is particularly valuable for apps with strict security requirements or those experiencing high authentication traffic.

For complete details on available metrics and how to use them, see ASP.NET Core metrics.

## Scaffold Register, Login, LogOut, and RegisterConfirmation

Add the `Register`, `Login`, `LogOut`, and `RegisterConfirmation` files. Follow the Scaffold identity into a Razor project with authorization instructions to generate the code shown in this section.

### Examine Register

When a user clicks the **Register** button on the `Register` page, the `RegisterModel.OnPostAsync` action is invoked. The user is created by CreateAsync(TUser) on the `_userManager` object:

```
public async Task<IActionResult> OnPostAsync(string returnUrl = null)
{
 returnUrl = returnUrl ?? Url.Content("~/");
 ExternalLogins = (await _signInManager.GetExternalAuthenticationSchemesAsync())
 .ToList();
 if (ModelState.IsValid)
 {
 var user = new IdentityUser { UserName = Input.Email, Email = Input.Email };
 var result = await _userManager.CreateAsync(user, Input.Password);
 if (result.Succeeded)
 {
 _logger.LogInformation("User created a new account with password.");
 var code = await _userManager.GenerateEmailConfirmationTokenAsync(user);
 code = WebEncoders.Base64UrlEncode(Encoding.UTF8.GetBytes(code));
 var callbackUrl = Url.Page(
 "/Account/ConfirmEmail",
 pageHandler: null,
 values: new { area = "Identity", userId = user.Id, code = code },
 protocol: Request.Scheme);
 await _emailSender.SendEmailAsync(Input.Email, "Confirm your email",
 $"Please confirm your account by <a href='{HtmlEncoder.Default.Encode(callbackUrl)}'>clicking here</a>.");
 if (_userManager.Options.SignIn.RequireConfirmedAccount)
 {
 return RedirectToPage("RegisterConfirmation",
 new { email = Input.Email });
 }
 else
 {
 await _signInManager.SignInAsync(user, isPersistent: false);
 return LocalRedirect(returnUrl);
 }
 }
 foreach (var error in result.Errors)
 {
 ModelState.AddModelError(string.Empty, error.Description);
 }
 }
 // If we got this far, something failed, redisplay form
 return Page();
}
```
### Disable default account verification

With the default templates, the user is redirected to the `Account.RegisterConfirmation` where they can select a link to have the account confirmed. The default `Account.RegisterConfirmation` is used * only* for testing, automatic account verification should be disabled in a production app.

To require a confirmed account and prevent immediate login at registration, set `DisplayConfirmAccountLink = false` in `/Areas/Identity/Pages/Account/RegisterConfirmation.cshtml.cs`:

```
[AllowAnonymous]
public class RegisterConfirmationModel : PageModel
{
 private readonly UserManager<IdentityUser> _userManager;
 private readonly IEmailSender _sender;
 public RegisterConfirmationModel(UserManager<IdentityUser> userManager, IEmailSender sender)
 {
 _userManager = userManager;
 _sender = sender;
 }
 public string Email { get; set; }
 public bool DisplayConfirmAccountLink { get; set; }
 public string EmailConfirmationUrl { get; set; }
 public async Task<IActionResult> OnGetAsync(string email, string returnUrl = null)
 {
 if (email == null)
 {
 return RedirectToPage("/Index");
 }
 var user = await _userManager.FindByEmailAsync(email);
 if (user == null)
 {
 return NotFound($"Unable to load user with email '{email}'.");
 }
 Email = email;
 // Once you add a real email sender, you should remove this code that lets you confirm the account
 DisplayConfirmAccountLink = false;
 if (DisplayConfirmAccountLink)
 {
 var userId = await _userManager.GetUserIdAsync(user);
 var code = await _userManager.GenerateEmailConfirmationTokenAsync(user);
 code = WebEncoders.Base64UrlEncode(Encoding.UTF8.GetBytes(code));
 EmailConfirmationUrl = Url.Page(
 "/Account/ConfirmEmail",
 pageHandler: null,
 values: new { area = "Identity", userId = userId, code = code, returnUrl = returnUrl },
 protocol: Request.Scheme);
 }
 return Page();
 }
}
```
### Log in

The Login form is displayed when:

- The **Log in**link is selected.
- A user attempts to access a restricted page that they aren't authorized to access **or**when they haven't been authenticated by the system.

When the form on the Login page is submitted, the `OnPostAsync` action is called. `PasswordSignInAsync` is called on the `_signInManager` object.

```
public async Task<IActionResult> OnPostAsync(string returnUrl = null)
{
 returnUrl = returnUrl ?? Url.Content("~/");
 if (ModelState.IsValid)
 {
 // This doesn't count login failures towards account lockout
 // To enable password failures to trigger account lockout,
 // set lockoutOnFailure: true
 var result = await _signInManager.PasswordSignInAsync(Input.Email,
 Input.Password, Input.RememberMe, lockoutOnFailure: true);
 if (result.Succeeded)
 {
 _logger.LogInformation("User logged in.");
 return LocalRedirect(returnUrl);
 }
 if (result.RequiresTwoFactor)
 {
 return RedirectToPage("./LoginWith2fa", new
 {
 ReturnUrl = returnUrl,
 RememberMe = Input.RememberMe
 });
 }
 if (result.IsLockedOut)
 {
 _logger.LogWarning("User account locked out.");
 return RedirectToPage("./Lockout");
 }
 else
 {
 ModelState.AddModelError(string.Empty, "Invalid login attempt.");
 return Page();
 }
 }
 // If we got this far, something failed, redisplay form
 return Page();
}
```
For information on how to make authorization decisions, see Introduction to authorization in ASP.NET Core.

### Log out

The **Log out** link invokes the `LogoutModel.OnPost` action.

```
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
using System.Threading.Tasks;
namespace WebApp1.Areas.Identity.Pages.Account
{
 [AllowAnonymous]
 public class LogoutModel : PageModel
 {
 private readonly SignInManager<IdentityUser> _signInManager;
 private readonly ILogger<LogoutModel> _logger;
 public LogoutModel(SignInManager<IdentityUser> signInManager, ILogger<LogoutModel> logger)
 {
 _signInManager = signInManager;
 _logger = logger;
 }
 public void OnGet()
 {
 }
 public async Task<IActionResult> OnPost(string returnUrl = null)
 {
 await _signInManager.SignOutAsync();
 _logger.LogInformation("User logged out.");
 if (returnUrl != null)
 {
 return LocalRedirect(returnUrl);
 }
 else
 {
 return RedirectToPage();
 }
 }
 }
}
```
In the preceding code, the code `return RedirectToPage();` needs to be a redirect so that the browser performs a new request and the identity for the user gets updated.

SignOutAsync clears the user's claims stored in a cookie.

Post is specified in the `Pages/Shared/_LoginPartial.cshtml`:

```
@using Microsoft.AspNetCore.Identity
@inject SignInManager<IdentityUser> SignInManager
@inject UserManager<IdentityUser> UserManager
<ul class="navbar-nav">
@if (SignInManager.IsSignedIn(User))
{
 <li class="nav-item">
 <a class="nav-link text-dark" asp-area="Identity" asp-page="/Account/Manage/Index"
 title="Manage">Hello @User.Identity.Name!</a>
 </li>
 <li class="nav-item">
 <form class="form-inline" asp-area="Identity" asp-page="/Account/Logout"
 asp-route-returnUrl="@Url.Page("/", new { area = "" })"
 method="post" >
 <button type="submit" class="nav-link btn btn-link text-dark">Logout</button>
 </form>
 </li>
}
else
{
 <li class="nav-item">
 <a class="nav-link text-dark" asp-area="Identity" asp-page="/Account/Register">Register</a>
 </li>
 <li class="nav-item">
 <a class="nav-link text-dark" asp-area="Identity" asp-page="/Account/Login">Login</a>
 </li>
}
</ul>
```
## Test Identity

The default web project templates allow anonymous access to the home pages. To test Identity, add `[Authorize]`:

```
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
namespace WebApp1.Pages
{
 [Authorize]
 public class PrivacyModel : PageModel
 {
 private readonly ILogger<PrivacyModel> _logger;
 public PrivacyModel(ILogger<PrivacyModel> logger)
 {
 _logger = logger;
 }
 public void OnGet()
 {
 }
 }
}
```
If you are signed in, sign out. Run the app and select the **Privacy** link. You are redirected to the login page.

### Explore Identity

To explore Identity in more detail:

- Create full identity UI source
- Examine the source of each page and step through the debugger.

## Identity Components

All the Identity-dependent NuGet packages are included in the ASP.NET Core shared framework.

The primary package for Identity is Microsoft.AspNetCore.Identity. This package contains the core set of interfaces for ASP.NET Core Identity, and is included by `Microsoft.AspNetCore.Identity.EntityFrameworkCore`.

## Migrating to ASP.NET Core Identity

For more information and guidance on migrating your existing Identity store, see Migrate Authentication and Identity.

## Setting password strength

See Configuration for a sample that sets the minimum password requirements.

## AddDefaultIdentity and AddIdentity

AddDefaultIdentity was introduced in ASP.NET Core 2.1. Calling `AddDefaultIdentity` is similar to calling the following:

See AddDefaultIdentity source for more information.

## Prevent publish of static Identity assets

To prevent publishing static Identity assets (stylesheets and JavaScript files for Identity UI) to the web root, add the following `ResolveStaticWebAssetsInputsDependsOn` property and `RemoveIdentityAssets` target to the app's project file:

```
<PropertyGroup>
 <ResolveStaticWebAssetsInputsDependsOn>RemoveIdentityAssets</ResolveStaticWebAssetsInputsDependsOn>
</PropertyGroup>
<Target Name="RemoveIdentityAssets">
 <ItemGroup>
 <StaticWebAsset Remove="@(StaticWebAsset)" Condition="%(SourceId) == 'Microsoft.AspNetCore.Identity.UI'" />
 </ItemGroup>
</Target>
```
## Next Steps

- ASP.NET Core Blazor authentication and authorization
- ASP.NET Core Identity source code
- How to work with Roles in ASP.NET Core Identity

- For information on configuring Identity using SQLite, see How to config Identity for SQLite (`dotnet/AspNetCore.Docs`#5131).
- Configure Identity
- Create an ASP.NET Core app with user data protected by authorization
- Add, download, and delete user data to Identity in an ASP.NET Core project
- Enable QR code generation for TOTP authentication
- Migrate Authentication and Identity to ASP.NET Core
- Account confirmation and password recovery
- Multifactor authentication in ASP.NET Core
- Host ASP.NET Core in a web farm

ASP.NET Core Identity:

- Is an API that supports user interface (UI) login functionality.
- Manages users, passwords, profile data, roles, claims, tokens, email confirmation, and more.

Users can create an account with the login information stored in Identity or they can use an external login provider. Supported external login providers include Facebook, Google, Microsoft Account, and Twitter.

For information on how to require authentication for all app users, see Create an ASP.NET Core app with user data protected by authorization.

The Identity source code is available on GitHub. Scaffold Identity and view the generated files to review the template interaction with Identity.

Identity is typically configured using a SQL Server database to store user names, passwords, and profile data. Alternatively, another persistent store can be used, for example, Azure Table Storage.

In this topic, you learn how to use Identity to register, log in, and log out a user. Note: the templates treat username and email as the same for users. For more detailed instructions about creating apps that use Identity, see Next Steps.

ASP.NET Core Identity isn't related to the Microsoft identity platform. Microsoft identity platform is:

- An evolution of the Azure Active Directory (Azure AD) developer platform.
- An alternative identity solution for authentication and authorization in ASP.NET Core apps.

ASP.NET Core Identity adds user interface (UI) login functionality to ASP.NET Core web apps. To secure web APIs and SPAs, use one of the following:

Duende Identity Server is an OpenID Connect and OAuth 2.0 framework for ASP.NET Core. Duende Identity Server enables the following security features:

- Authentication as a Service (AaaS)
- Single sign-on/off (SSO) over multiple application types
- Access control for APIs
- Federation Gateway

Important

Duende Software might require you to pay a license fee for production use of Duende Identity Server. For more information, see Migrate from ASP.NET Core in .NET 5 to .NET 6.

For more information, see the Duende Identity Server documentation (Duende Software website).

View or download the sample code (how to download).

## Create a Web app with authentication

Create an ASP.NET Core Web Application project with Individual User Accounts.

- Select the **ASP.NET Core Web App**template. Name the project**WebApp1**to have the same namespace as the project download. Click**OK**.
- In the **Authentication type**input, select**Individual User Accounts**.

The generated project provides ASP.NET Core Identity as a Razor class library. The Identity Razor class library exposes endpoints with the `Identity` area. For example:

- /Identity/Account/Login
- /Identity/Account/Logout
- /Identity/Account/Manage

### Apply migrations

Apply the migrations to initialize the database.

Run the following command in the Package Manager Console (PMC):

`Update-Database`

### Test Register and Login

Run the app and register a user. Depending on your screen size, you might need to select the navigation toggle button to see the **Register** and **Login** links.

### View the Identity database

- From the **View**menu, select**SQL Server Object Explorer**(SSOX).
- Navigate to **(localdb)MSSQLLocalDB(SQL Server 13)**. Right-click on**dbo.AspNetUsers**>**View Data**:

### Configure Identity services

Services are added in `Program.cs`. The typical pattern is to call methods in the following order:

- `Add{Service}`
- `builder.Services.Configure{Service}`

```
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using WebApp1.Data;
var builder = WebApplication.CreateBuilder(args);
var connectionString = builder.Configuration.GetConnectionString("DefaultConnection");
builder.Services.AddDbContext<ApplicationDbContext>(options =>
 options.UseSqlServer(connectionString));
builder.Services.AddDatabaseDeveloperPageExceptionFilter();
builder.Services.AddDefaultIdentity<IdentityUser>(options => options.SignIn.RequireConfirmedAccount = true)
 .AddEntityFrameworkStores<ApplicationDbContext>();
builder.Services.AddRazorPages();
builder.Services.Configure<IdentityOptions>(options =>
{
 // Password settings.
 options.Password.RequireDigit = true;
 options.Password.RequireLowercase = true;
 options.Password.RequireNonAlphanumeric = true;
 options.Password.RequireUppercase = true;
 options.Password.RequiredLength = 6;
 options.Password.RequiredUniqueChars = 1;
 // Lockout settings.
 options.Lockout.DefaultLockoutTimeSpan = TimeSpan.FromMinutes(5);
 options.Lockout.MaxFailedAccessAttempts = 5;
 options.Lockout.AllowedForNewUsers = true;
 // User settings.
 options.User.AllowedUserNameCharacters =
 "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._@+";
 options.User.RequireUniqueEmail = false;
});
builder.Services.ConfigureApplicationCookie(options =>
{
 // Cookie settings
 options.Cookie.HttpOnly = true;
 options.ExpireTimeSpan = TimeSpan.FromMinutes(5);
 options.LoginPath = "/Identity/Account/Login";
 options.AccessDeniedPath = "/Identity/Account/AccessDenied";
 options.SlidingExpiration = true;
});
var app = builder.Build();
if (app.Environment.IsDevelopment())
{
 app.UseMigrationsEndPoint();
}
else
{
 app.UseExceptionHandler("/Error");
 app.UseHsts();
}
app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseRouting();
app.UseAuthentication();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
```
The preceding code configures Identity with default option values. Services are made available to the app through dependency injection.

Identity is enabled by calling UseAuthentication. `UseAuthentication` adds authentication middleware to the request pipeline.

The template-generated app doesn't use authorization. `app.UseAuthorization` is included to ensure it's added in the correct order should the app add authorization. `UseRouting`, `UseAuthentication`, and `UseAuthorization` must be called in the order shown in the preceding code.

For more information on `IdentityOptions`, see IdentityOptions and Application Startup.

## Scaffold Register, Login, LogOut, and RegisterConfirmation

Add the `Register`, `Login`, `LogOut`, and `RegisterConfirmation` files. Follow the Scaffold identity into a Razor project with authorization instructions to generate the code shown in this section.

### Examine Register

When a user clicks the **Register** button on the `Register` page, the `RegisterModel.OnPostAsync` action is invoked. The user is created by CreateAsync(TUser) on the `_userManager` object:

```
public async Task<IActionResult> OnPostAsync(string returnUrl = null)
{
 returnUrl = returnUrl ?? Url.Content("~/");
 ExternalLogins = (await _signInManager.GetExternalAuthenticationSchemesAsync())
 .ToList();
 if (ModelState.IsValid)
 {
 var user = new IdentityUser { UserName = Input.Email, Email = Input.Email };
 var result = await _userManager.CreateAsync(user, Input.Password);
 if (result.Succeeded)
 {
 _logger.LogInformation("User created a new account with password.");
 var code = await _userManager.GenerateEmailConfirmationTokenAsync(user);
 code = WebEncoders.Base64UrlEncode(Encoding.UTF8.GetBytes(code));
 var callbackUrl = Url.Page(
 "/Account/ConfirmEmail",
 pageHandler: null,
 values: new { area = "Identity", userId = user.Id, code = code },
 protocol: Request.Scheme);
 await _emailSender.SendEmailAsync(Input.Email, "Confirm your email",
 $"Please confirm your account by <a href='{HtmlEncoder.Default.Encode(callbackUrl)}'>clicking here</a>.");
 if (_userManager.Options.SignIn.RequireConfirmedAccount)
 {
 return RedirectToPage("RegisterConfirmation",
 new { email = Input.Email });
 }
 else
 {
 await _signInManager.SignInAsync(user, isPersistent: false);
 return LocalRedirect(returnUrl);
 }
 }
 foreach (var error in result.Errors)
 {
 ModelState.AddModelError(string.Empty, error.Description);
 }
 }
 // If we got this far, something failed, redisplay form
 return Page();
}
```
### Disable default account verification

With the default templates, the user is redirected to the `Account.RegisterConfirmation` where they can select a link to have the account confirmed. The default `Account.RegisterConfirmation` is used * only* for testing, automatic account verification should be disabled in a production app.

To require a confirmed account and prevent immediate login at registration, set `DisplayConfirmAccountLink = false` in `/Areas/Identity/Pages/Account/RegisterConfirmation.cshtml.cs`:

```
[AllowAnonymous]
public class RegisterConfirmationModel : PageModel
{
 private readonly UserManager<IdentityUser> _userManager;
 private readonly IEmailSender _sender;
 public RegisterConfirmationModel(UserManager<IdentityUser> userManager, IEmailSender sender)
 {
 _userManager = userManager;
 _sender = sender;
 }
 public string Email { get; set; }
 public bool DisplayConfirmAccountLink { get; set; }
 public string EmailConfirmationUrl { get; set; }
 public async Task<IActionResult> OnGetAsync(string email, string returnUrl = null)
 {
 if (email == null)
 {
 return RedirectToPage("/Index");
 }
 var user = await _userManager.FindByEmailAsync(email);
 if (user == null)
 {
 return NotFound($"Unable to load user with email '{email}'.");
 }
 Email = email;
 // Once you add a real email sender, you should remove this code that lets you confirm the account
 DisplayConfirmAccountLink = false;
 if (DisplayConfirmAccountLink)
 {
 var userId = await _userManager.GetUserIdAsync(user);
 var code = await _userManager.GenerateEmailConfirmationTokenAsync(user);
 code = WebEncoders.Base64UrlEncode(Encoding.UTF8.GetBytes(code));
 EmailConfirmationUrl = Url.Page(
 "/Account/ConfirmEmail",
 pageHandler: null,
 values: new { area = "Identity", userId = userId, code = code, returnUrl = returnUrl },
 protocol: Request.Scheme);
 }
 return Page();
 }
}
```
### Log in

The Login form is displayed when:

- The **Log in**link is selected.
- A user attempts to access a restricted page that they aren't authorized to access **or**when they haven't been authenticated by the system.

When the form on the Login page is submitted, the `OnPostAsync` action is called. `PasswordSignInAsync` is called on the `_signInManager` object.

```
public async Task<IActionResult> OnPostAsync(string returnUrl = null)
{
 returnUrl = returnUrl ?? Url.Content("~/");
 if (ModelState.IsValid)
 {
 // This doesn't count login failures towards account lockout
 // To enable password failures to trigger account lockout,
 // set lockoutOnFailure: true
 var result = await _signInManager.PasswordSignInAsync(Input.Email,
 Input.Password, Input.RememberMe, lockoutOnFailure: true);
 if (result.Succeeded)
 {
 _logger.LogInformation("User logged in.");
 return LocalRedirect(returnUrl);
 }
 if (result.RequiresTwoFactor)
 {
 return RedirectToPage("./LoginWith2fa", new
 {
 ReturnUrl = returnUrl,
 RememberMe = Input.RememberMe
 });
 }
 if (result.IsLockedOut)
 {
 _logger.LogWarning("User account locked out.");
 return RedirectToPage("./Lockout");
 }
 else
 {
 ModelState.AddModelError(string.Empty, "Invalid login attempt.");
 return Page();
 }
 }
 // If we got this far, something failed, redisplay form
 return Page();
}
```
For information on how to make authorization decisions, see Introduction to authorization in ASP.NET Core.

### Log out

The **Log out** link invokes the `LogoutModel.OnPost` action.

```
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
using System.Threading.Tasks;
namespace WebApp1.Areas.Identity.Pages.Account
{
 [AllowAnonymous]
 public class LogoutModel : PageModel
 {
 private readonly SignInManager<IdentityUser> _signInManager;
 private readonly ILogger<LogoutModel> _logger;
 public LogoutModel(SignInManager<IdentityUser> signInManager, ILogger<LogoutModel> logger)
 {
 _signInManager = signInManager;
 _logger = logger;
 }
 public void OnGet()
 {
 }
 public async Task<IActionResult> OnPost(string returnUrl = null)
 {
 await _signInManager.SignOutAsync();
 _logger.LogInformation("User logged out.");
 if (returnUrl != null)
 {
 return LocalRedirect(returnUrl);
 }
 else
 {
 return RedirectToPage();
 }
 }
 }
}
```
In the preceding code, the code `return RedirectToPage();` needs to be a redirect so that the browser performs a new request and the identity for the user gets updated.

SignOutAsync clears the user's claims stored in a cookie.

Post is specified in the `Pages/Shared/_LoginPartial.cshtml`:

```
@using Microsoft.AspNetCore.Identity
@inject SignInManager<IdentityUser> SignInManager
@inject UserManager<IdentityUser> UserManager
<ul class="navbar-nav">
@if (SignInManager.IsSignedIn(User))
{
 <li class="nav-item">
 <a class="nav-link text-dark" asp-area="Identity" asp-page="/Account/Manage/Index"
 title="Manage">Hello @User.Identity.Name!</a>
 </li>
 <li class="nav-item">
 <form class="form-inline" asp-area="Identity" asp-page="/Account/Logout"
 asp-route-returnUrl="@Url.Page("/", new { area = "" })"
 method="post" >
 <button type="submit" class="nav-link btn btn-link text-dark">Logout</button>
 </form>
 </li>
}
else
{
 <li class="nav-item">
 <a class="nav-link text-dark" asp-area="Identity" asp-page="/Account/Register">Register</a>
 </li>
 <li class="nav-item">
 <a class="nav-link text-dark" asp-area="Identity" asp-page="/Account/Login">Login</a>
 </li>
}
</ul>
```
## Test Identity

The default web project templates allow anonymous access to the home pages. To test Identity, add `[Authorize]`:

```
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
namespace WebApp1.Pages
{
 [Authorize]
 public class PrivacyModel : PageModel
 {
 private readonly ILogger<PrivacyModel> _logger;
 public PrivacyModel(ILogger<PrivacyModel> logger)
 {
 _logger = logger;
 }
 public void OnGet()
 {
 }
 }
}
```
If you are signed in, sign out. Run the app and select the **Privacy** link. You are redirected to the login page.

### Explore Identity

To explore Identity in more detail:

- Create full identity UI source
- Examine the source of each page and step through the debugger.

## Identity Components

All the Identity-dependent NuGet packages are included in the ASP.NET Core shared framework.

The primary package for Identity is Microsoft.AspNetCore.Identity. This package contains the core set of interfaces for ASP.NET Core Identity, and is included by `Microsoft.AspNetCore.Identity.EntityFrameworkCore`.

## Migrating to ASP.NET Core Identity

For more information and guidance on migrating your existing Identity store, see Migrate Authentication and Identity.

## Setting password strength

See Configuration for a sample that sets the minimum password requirements.

## AddDefaultIdentity and AddIdentity

AddDefaultIdentity was introduced in ASP.NET Core 2.1. Calling `AddDefaultIdentity` is similar to calling the following:

See AddDefaultIdentity source for more information.

## Prevent publish of static Identity assets

To prevent publishing static Identity assets (stylesheets and JavaScript files for Identity UI) to the web root, add the following `ResolveStaticWebAssetsInputsDependsOn` property and `RemoveIdentityAssets` target to the app's project file:

```
<PropertyGroup>
 <ResolveStaticWebAssetsInputsDependsOn>RemoveIdentityAssets</ResolveStaticWebAssetsInputsDependsOn>
</PropertyGroup>
<Target Name="RemoveIdentityAssets">
 <ItemGroup>
 <StaticWebAsset Remove="@(StaticWebAsset)" Condition="%(SourceId) == 'Microsoft.AspNetCore.Identity.UI'" />
 </ItemGroup>
</Target>
```
## Next Steps

- See this GitHub issue for information on configuring Identity using SQLite.
- Configure Identity
- Create an ASP.NET Core app with user data protected by authorization
- Add, download, and delete user data to Identity in an ASP.NET Core project
- Enable QR code generation for TOTP authentication
- Migrate Authentication and Identity to ASP.NET Core
- Account confirmation and password recovery
- Two-factor authentication with SMS in ASP.NET Core
- Host ASP.NET Core in a web farm

ASP.NET Core Identity:

- Is an API that supports user interface (UI) login functionality.
- Manages users, passwords, profile data, roles, claims, tokens, email confirmation, and more.

Users can create an account with the login information stored in Identity or they can use an external login provider. Supported external login providers include Facebook, Google, Microsoft Account, and Twitter.

For information on how to require authentication for all app users, see Create an ASP.NET Core app with user data protected by authorization.

The Identity source code is available on GitHub. Scaffold Identity and view the generated files to review the template interaction with Identity.

Identity is typically configured using a SQL Server database to store user names, passwords, and profile data. Alternatively, another persistent store can be used, for example, Azure Table Storage.

In this topic, you learn how to use Identity to register, log in, and log out a user. Note: the templates treat username and email as the same for users. For more detailed instructions about creating apps that use Identity, see Next Steps.

Microsoft identity platform is:

- An evolution of the Azure Active Directory (Azure AD) developer platform.
- An alternative identity solution for authentication and authorization in ASP.NET Core apps.
- Not related to ASP.NET Core Identity.

ASP.NET Core Identity adds user interface (UI) login functionality to ASP.NET Core web apps. To secure web APIs and SPAs, use one of the following:

- Microsoft Entra ID
- Duende IdentityServer. Duende IdentityServer is 3rd party product.

Duende IdentityServer is an OpenID Connect and OAuth 2.0 framework for ASP.NET Core. Duende IdentityServer enables the following security features:

- Authentication as a Service (AaaS)
- Single sign-on/off (SSO) over multiple application types
- Access control for APIs
- Federation Gateway

For more information, see Overview of Duende IdentityServer.

For more information on other authentication providers, see Community OSS authentication options for ASP.NET Core

View or download the sample code (how to download).

## Create a Web app with authentication

Create an ASP.NET Core Web Application project with Individual User Accounts.

- Select **File**>**New**>**Project**.
- Select **ASP.NET Core Web Application**. Name the project**WebApp1**to have the same namespace as the project download. Click**OK**.
- Select an ASP.NET Core **Web Application**, then select**Change Authentication**.
- Select **Individual User Accounts**and click**OK**.

The generated project provides ASP.NET Core Identity as a Razor class library. The Identity Razor class library exposes endpoints with the `Identity` area. For example:

- /Identity/Account/Login
- /Identity/Account/Logout
- /Identity/Account/Manage

### Apply migrations

Apply the migrations to initialize the database.

Run the following command in the Package Manager Console (PMC):

`PM> Update-Database`

### Test Register and Login

Run the app and register a user. Depending on your screen size, you might need to select the navigation toggle button to see the **Register** and **Login** links.

### View the Identity database

- From the **View**menu, select**SQL Server Object Explorer**(SSOX).
- Navigate to **(localdb)MSSQLLocalDB(SQL Server 13)**. Right-click on**dbo.AspNetUsers**>**View Data**:

### Configure Identity services

Services are added in `ConfigureServices`. The typical pattern is to call all the `Add{Service}` methods, and then call all the `services.Configure{Service}` methods.

```
public void ConfigureServices(IServiceCollection services)
{
 services.AddDbContext<ApplicationDbContext>(options =>
 // options.UseSqlite(
 options.UseSqlServer(
 Configuration.GetConnectionString("DefaultConnection")));
 services.AddDefaultIdentity<IdentityUser>(options => options.SignIn.RequireConfirmedAccount = true)
 .AddEntityFrameworkStores<ApplicationDbContext>();
 services.AddRazorPages();
 services.Configure<IdentityOptions>(options =>
 {
 // Password settings.
 options.Password.RequireDigit = true;
 options.Password.RequireLowercase = true;
 options.Password.RequireNonAlphanumeric = true;
 options.Password.RequireUppercase = true;
 options.Password.RequiredLength = 6;
 options.Password.RequiredUniqueChars = 1;
 // Lockout settings.
 options.Lockout.DefaultLockoutTimeSpan = TimeSpan.FromMinutes(5);
 options.Lockout.MaxFailedAccessAttempts = 5;
 options.Lockout.AllowedForNewUsers = true;
 // User settings.
 options.User.AllowedUserNameCharacters =
 "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._@+";
 options.User.RequireUniqueEmail = false;
 });
 services.ConfigureApplicationCookie(options =>
 {
 // Cookie settings
 options.Cookie.HttpOnly = true;
 options.ExpireTimeSpan = TimeSpan.FromMinutes(5);
 options.LoginPath = "/Identity/Account/Login";
 options.AccessDeniedPath = "/Identity/Account/AccessDenied";
 options.SlidingExpiration = true;
 });
}
```
The preceding highlighted code configures Identity with default option values. Services are made available to the app through dependency injection.

Identity is enabled by calling UseAuthentication. `UseAuthentication` adds authentication middleware to the request pipeline.

```
public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
{
 if (env.IsDevelopment())
 {
 app.UseDeveloperExceptionPage();
 app.UseDatabaseErrorPage();
 }
 else
 {
 app.UseExceptionHandler("/Error");
 app.UseHsts();
 }
 app.UseHttpsRedirection();
 app.UseStaticFiles();
 app.UseRouting();
 app.UseAuthentication();
 app.UseAuthorization();
 app.UseEndpoints(endpoints =>
 {
 endpoints.MapRazorPages();
 });
}
```
```
public void ConfigureServices(IServiceCollection services)
{
 services.AddDbContext<ApplicationDbContext>(options =>
 // options.UseSqlite(
 options.UseSqlServer(
 Configuration.GetConnectionString("DefaultConnection")));
 services.AddDatabaseDeveloperPageExceptionFilter();
 services.AddDefaultIdentity<IdentityUser>(options => options.SignIn.RequireConfirmedAccount = true)
 .AddEntityFrameworkStores<ApplicationDbContext>();
 services.AddRazorPages();
 services.Configure<IdentityOptions>(options =>
 {
 // Password settings.
 options.Password.RequireDigit = true;
 options.Password.RequireLowercase = true;
 options.Password.RequireNonAlphanumeric = true;
 options.Password.RequireUppercase = true;
 options.Password.RequiredLength = 6;
 options.Password.RequiredUniqueChars = 1;
 // Lockout settings.
 options.Lockout.DefaultLockoutTimeSpan = TimeSpan.FromMinutes(5);
 options.Lockout.MaxFailedAccessAttempts = 5;
 options.Lockout.AllowedForNewUsers = true;
 // User settings.
 options.User.AllowedUserNameCharacters =
 "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._@+";
 options.User.RequireUniqueEmail = false;
 });
 services.ConfigureApplicationCookie(options =>
 {
 // Cookie settings
 options.Cookie.HttpOnly = true;
 options.ExpireTimeSpan = TimeSpan.FromMinutes(5);
 options.LoginPath = "/Identity/Account/Login";
 options.AccessDeniedPath = "/Identity/Account/AccessDenied";
 options.SlidingExpiration = true;
 });
}
```
The preceding code configures Identity with default option values. Services are made available to the app through dependency injection.

Identity is enabled by calling UseAuthentication. `UseAuthentication` adds authentication middleware to the request pipeline.

```
public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
{
 if (env.IsDevelopment())
 {
 app.UseDeveloperExceptionPage();
 app.UseMigrationsEndPoint();
 }
 else
 {
 app.UseExceptionHandler("/Error");
 app.UseHsts();
 }
 app.UseHttpsRedirection();
 app.UseStaticFiles();
 app.UseRouting();
 app.UseAuthentication();
 app.UseAuthorization();
 app.UseEndpoints(endpoints =>
 {
 endpoints.MapRazorPages();
 });
}
```
The template-generated app doesn't use authorization. `app.UseAuthorization` is included to ensure it's added in the correct order should the app add authorization. `UseRouting`, `UseAuthentication`, `UseAuthorization`, and `UseEndpoints` must be called in the order shown in the preceding code.

For more information on `IdentityOptions` and `Startup`, see IdentityOptions and Application Startup.

## Scaffold Register, Login, LogOut, and RegisterConfirmation

Add the `Register`, `Login`, `LogOut`, and `RegisterConfirmation` files. Follow the Scaffold identity into a Razor project with authorization instructions to generate the code shown in this section.

### Examine Register

When a user clicks the **Register** button on the `Register` page, the `RegisterModel.OnPostAsync` action is invoked. The user is created by CreateAsync(TUser) on the `_userManager` object:

```
public async Task<IActionResult> OnPostAsync(string returnUrl = null)
{
 returnUrl = returnUrl ?? Url.Content("~/");
 ExternalLogins = (await _signInManager.GetExternalAuthenticationSchemesAsync())
 .ToList();
 if (ModelState.IsValid)
 {
 var user = new IdentityUser { UserName = Input.Email, Email = Input.Email };
 var result = await _userManager.CreateAsync(user, Input.Password);
 if (result.Succeeded)
 {
 _logger.LogInformation("User created a new account with password.");
 var code = await _userManager.GenerateEmailConfirmationTokenAsync(user);
 code = WebEncoders.Base64UrlEncode(Encoding.UTF8.GetBytes(code));
 var callbackUrl = Url.Page(
 "/Account/ConfirmEmail",
 pageHandler: null,
 values: new { area = "Identity", userId = user.Id, code = code },
 protocol: Request.Scheme);
 await _emailSender.SendEmailAsync(Input.Email, "Confirm your email",
 $"Please confirm your account by <a href='{HtmlEncoder.Default.Encode(callbackUrl)}'>clicking here</a>.");
 if (_userManager.Options.SignIn.RequireConfirmedAccount)
 {
 return RedirectToPage("RegisterConfirmation",
 new { email = Input.Email });
 }
 else
 {
 await _signInManager.SignInAsync(user, isPersistent: false);
 return LocalRedirect(returnUrl);
 }
 }
 foreach (var error in result.Errors)
 {
 ModelState.AddModelError(string.Empty, error.Description);
 }
 }
 // If we got this far, something failed, redisplay form
 return Page();
}
```
### Disable default account verification

With the default templates, the user is redirected to the `Account.RegisterConfirmation` where they can select a link to have the account confirmed. The default `Account.RegisterConfirmation` is used * only* for testing, automatic account verification should be disabled in a production app.

To require a confirmed account and prevent immediate login at registration, set `DisplayConfirmAccountLink = false` in `/Areas/Identity/Pages/Account/RegisterConfirmation.cshtml.cs`:

```
[AllowAnonymous]
public class RegisterConfirmationModel : PageModel
{
 private readonly UserManager<IdentityUser> _userManager;
 private readonly IEmailSender _sender;
 public RegisterConfirmationModel(UserManager<IdentityUser> userManager, IEmailSender sender)
 {
 _userManager = userManager;
 _sender = sender;
 }
 public string Email { get; set; }
 public bool DisplayConfirmAccountLink { get; set; }
 public string EmailConfirmationUrl { get; set; }
 public async Task<IActionResult> OnGetAsync(string email, string returnUrl = null)
 {
 if (email == null)
 {
 return RedirectToPage("/Index");
 }
 var user = await _userManager.FindByEmailAsync(email);
 if (user == null)
 {
 return NotFound($"Unable to load user with email '{email}'.");
 }
 Email = email;
 // Once you add a real email sender, you should remove this code that lets you confirm the account
 DisplayConfirmAccountLink = false;
 if (DisplayConfirmAccountLink)
 {
 var userId = await _userManager.GetUserIdAsync(user);
 var code = await _userManager.GenerateEmailConfirmationTokenAsync(user);
 code = WebEncoders.Base64UrlEncode(Encoding.UTF8.GetBytes(code));
 EmailConfirmationUrl = Url.Page(
 "/Account/ConfirmEmail",
 pageHandler: null,
 values: new { area = "Identity", userId = userId, code = code, returnUrl = returnUrl },
 protocol: Request.Scheme);
 }
 return Page();
 }
}
```
### Log in

The Login form is displayed when:

- The **Log in**link is selected.
- A user attempts to access a restricted page that they aren't authorized to access **or**when they haven't been authenticated by the system.

When the form on the Login page is submitted, the `OnPostAsync` action is called. `PasswordSignInAsync` is called on the `_signInManager` object.

```
public async Task<IActionResult> OnPostAsync(string returnUrl = null)
{
 returnUrl = returnUrl ?? Url.Content("~/");
 if (ModelState.IsValid)
 {
 // This doesn't count login failures towards account lockout
 // To enable password failures to trigger account lockout,
 // set lockoutOnFailure: true
 var result = await _signInManager.PasswordSignInAsync(Input.Email,
 Input.Password, Input.RememberMe, lockoutOnFailure: true);
 if (result.Succeeded)
 {
 _logger.LogInformation("User logged in.");
 return LocalRedirect(returnUrl);
 }
 if (result.RequiresTwoFactor)
 {
 return RedirectToPage("./LoginWith2fa", new
 {
 ReturnUrl = returnUrl,
 RememberMe = Input.RememberMe
 });
 }
 if (result.IsLockedOut)
 {
 _logger.LogWarning("User account locked out.");
 return RedirectToPage("./Lockout");
 }
 else
 {
 ModelState.AddModelError(string.Empty, "Invalid login attempt.");
 return Page();
 }
 }
 // If we got this far, something failed, redisplay form
 return Page();
}
```
For information on how to make authorization decisions, see Introduction to authorization in ASP.NET Core.

### Log out

The **Log out** link invokes the `LogoutModel.OnPost` action.

```
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
using System.Threading.Tasks;
namespace WebApp1.Areas.Identity.Pages.Account
{
 [AllowAnonymous]
 public class LogoutModel : PageModel
 {
 private readonly SignInManager<IdentityUser> _signInManager;
 private readonly ILogger<LogoutModel> _logger;
 public LogoutModel(SignInManager<IdentityUser> signInManager, ILogger<LogoutModel> logger)
 {
 _signInManager = signInManager;
 _logger = logger;
 }
 public void OnGet()
 {
 }
 public async Task<IActionResult> OnPost(string returnUrl = null)
 {
 await _signInManager.SignOutAsync();
 _logger.LogInformation("User logged out.");
 if (returnUrl != null)
 {
 return LocalRedirect(returnUrl);
 }
 else
 {
 return RedirectToPage();
 }
 }
 }
}
```
In the preceding code, the code `return RedirectToPage();` needs to be a redirect so that the browser performs a new request and the identity for the user gets updated.

SignOutAsync clears the user's claims stored in a cookie.

Post is specified in the `Pages/Shared/_LoginPartial.cshtml`:

```
@using Microsoft.AspNetCore.Identity
@inject SignInManager<IdentityUser> SignInManager
@inject UserManager<IdentityUser> UserManager
<ul class="navbar-nav">
@if (SignInManager.IsSignedIn(User))
{
 <li class="nav-item">
 <a class="nav-link text-dark" asp-area="Identity" asp-page="/Account/Manage/Index"
 title="Manage">Hello @User.Identity.Name!</a>
 </li>
 <li class="nav-item">
 <form class="form-inline" asp-area="Identity" asp-page="/Account/Logout"
 asp-route-returnUrl="@Url.Page("/", new { area = "" })"
 method="post" >
 <button type="submit" class="nav-link btn btn-link text-dark">Logout</button>
 </form>
 </li>
}
else
{
 <li class="nav-item">
 <a class="nav-link text-dark" asp-area="Identity" asp-page="/Account/Register">Register</a>
 </li>
 <li class="nav-item">
 <a class="nav-link text-dark" asp-area="Identity" asp-page="/Account/Login">Login</a>
 </li>
}
</ul>
```
## Test Identity

The default web project templates allow anonymous access to the home pages. To test Identity, add `[Authorize]`:

```
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;
namespace WebApp1.Pages
{
 [Authorize]
 public class PrivacyModel : PageModel
 {
 private readonly ILogger<PrivacyModel> _logger;
 public PrivacyModel(ILogger<PrivacyModel> logger)
 {
 _logger = logger;
 }
 public void OnGet()
 {
 }
 }
}
```
If you are signed in, sign out. Run the app and select the **Privacy** link. You are redirected to the login page.

### Explore Identity

To explore Identity in more detail:

- Create full identity UI source
- Examine the source of each page and step through the debugger.

## Identity Components

All the Identity-dependent NuGet packages are included in the ASP.NET Core shared framework.

The primary package for Identity is Microsoft.AspNetCore.Identity. This package contains the core set of interfaces for ASP.NET Core Identity, and is included by `Microsoft.AspNetCore.Identity.EntityFrameworkCore`.

## Migrating to ASP.NET Core Identity

For more information and guidance on migrating your existing Identity store, see Migrate Authentication and Identity.

## Setting password strength

See Configuration for a sample that sets the minimum password requirements.

## Prevent publish of static Identity assets

To prevent publishing static Identity assets (stylesheets and JavaScript files for Identity UI) to the web root, add the following `ResolveStaticWebAssetsInputsDependsOn` property and `RemoveIdentityAssets` target to the app's project file:

```
<PropertyGroup>
 <ResolveStaticWebAssetsInputsDependsOn>RemoveIdentityAssets</ResolveStaticWebAssetsInputsDependsOn>
</PropertyGroup>
<Target Name="RemoveIdentityAssets">
 <ItemGroup>
 <StaticWebAsset Remove="@(StaticWebAsset)" Condition="%(SourceId) == 'Microsoft.AspNetCore.Identity.UI'" />
 </ItemGroup>
</Target>
```
## Next Steps

- ASP.NET Core Identity source code
- AddDefaultIdentity source
- See this GitHub issue for information on configuring Identity using SQLite.
- Configure Identity
- Create an ASP.NET Core app with user data protected by authorization
- Add, download, and delete user data to Identity in an ASP.NET Core project
- Enable QR code generation for TOTP authentication
- Migrate Authentication and Identity to ASP.NET Core
- Account confirmation and password recovery
- Two-factor authentication with SMS in ASP.NET Core
- Host ASP.NET Core in a web farm

ASP.NET Core

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

ASP.NET Core includes support for authenticator applications for user authentication. Two-factor authentication (2FA) authenticator apps use a Time-based One-time Password Algorithm (TOTP), the industry-recommended approach for 2FA. (TOTP-based 2FA is preferred over SMS 2FA.) Users typically install the authenticator app on a smartphone. The app provides a 6 to 8 digit code that the user enters after they confirm their username and password.

Warning

Keep the ASP.NET Core TOTP code secret. The user can enter the code multiple times and authenticate successfully before it expires.

The ASP.NET Core web app templates support authenticators, but they don't provide support for QR code generation. QR code generators make it easier to set up 2FA. This article provides guidance for Razor Pages and MVC apps on how to add QR code generation to the 2FA configuration page.

The ASP.NET Core web app templates support authenticators but don't provide support for QR code generation. QR code generators make it easier to set up 2FA. This article guides you through adding QR code generation to the 2FA configuration page.

Two-factor authentication doesn't happen by using an external authentication provider, such as Google or Facebook. External sign ins are protected by whatever mechanism the external authentication provider supports. For example, the Microsoft authentication provider requires a hardware key or another 2FA approach. When the default templates require 2FA for both the web app and the external authentication provider, users need to satisfy two 2FA approaches. Requiring two 2FA approaches deviates from established security practices, which typically rely on a single, strong 2FA method for authentication.

If you're working with Blazor in ASP.NET Core 8.0 or later, you can find similar guidance in the following articles:

- Enable QR code generation for TOTP authenticator apps in an ASP.NET Core Blazor Web App
- Enable QR code generation for TOTP authenticator apps in ASP.NET Core Blazor WebAssembly with ASP.NET Core Identity

## Add QR codes to the 2FA configuration page

The following instructions use the *qrcode.js* file from the https://davidshimjs.github.io/qrcodejs/ repo.

- Download the 'qrcode.js' JavaScript library to the - *wwwroot\lib*folder in your project.
- Follow the instructions in Scaffold Identity to generate the - */Areas/Identity/Pages/Account/Manage/EnableAuthenticator.cshtml*file.
- In the - */Areas/Identity/Pages/Account/Manage/EnableAuthenticator.cshtml*file, locate the- `Scripts`section at the end of the file:- `@section Scripts { <partial name="_ValidationScriptsPartial" /> }`
- Create a new JavaScript file named - *qr.js*in the- *wwwroot/js*folder, and add the following code that generates the QR code:- `window.addEventListener("load", () => { const uri = document.getElementById("qrCodeData").getAttribute('data-url'); new QRCode(document.getElementById("qrCode"), { text: uri, width: 150, height: 150 }); });`
- Update the - `Scripts`section to add a reference to the- `qrcode.js`library you previously downloaded.
- Add the - *qr.js*file with the call that generates the QR code:- `@section Scripts { <partial name="_ValidationScriptsPartial" /> <script type="text/javascript" src="~/lib/qrcode.js"></script> <script type="text/javascript" src="~/js/qr.js"></script> }`
- Delete the paragraph that links you to these instructions.
- Run your app. Confirm you can scan the QR code and validate the code the authenticator provides.

## Change the site name in the QR code

The site name in the QR code comes from the project name you select when you create your project. You can change it by looking for the `GenerateQrCodeUri(string email, string unformattedKey)` method in the */Areas/Identity/Pages/Account/Manage/EnableAuthenticator.cshtml.cs* file.

Here's the default code from the template:

```
private string GenerateQrCodeUri(string email, string unformattedKey)
{
 return string.Format(
 AuthenticatorUriFormat,
 _urlEncoder.Encode("Razor Pages"),
 _urlEncoder.Encode(email),
 unformattedKey);
}
```
The second parameter in the call to `string.Format` is your site name, which is obtained from your solution name. You can change it to any value, but it must always be URL encoded.

## Use a different QR code library

You can replace the QR code library with your preferred library. The HTML contains a `qrCode` element into which you can place a QR code by whatever mechanism your library provides.

You can find the correctly formatted URL for the QR code in the following locations:

- `AuthenticatorUri`property of the model
- `data-url`property in the- `qrCodeData`element

## Check TOTP client and server times

TOTP (Time-based One-Time Password) authentication depends on both the server and authenticator device having an accurate time. Tokens only last for 30 seconds. If TOTP 2FA sign-in fails, confirm the server time is accurate, and preferably synchronized to an accurate NTP service.

## Related content

ASP.NET Core

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

By Valeriy Novytskyy and Rick Anderson

This article explains how to build an ASP.NET Core app that enables users to sign in using OAuth 2.0 with credentials from external authentication providers.

Facebook, X (formerly Twitter), Google, and Microsoft providers are covered in the following sections and use the starter project created in this article. Other providers are available in third-party packages such as OpenIddict, AspNet.Security.OAuth.Providers and AspNet.Security.OpenId.Providers.

Enabling users to sign in with their existing credentials is convenient for the users and shifts many of the complexities of managing the sign-in process onto a third party.

## Create a New ASP.NET Core Project

- Select the **ASP.NET Core Web App**template and select**Next**.
- Name the project and select **Next**.
- In the **Authentication type**dropdown, select**Individual Accounts**and select**Create**.

## Apply migrations

- Run the app and select the **Register**link.
- Enter the email and password for the new account, and then select **Register**.
- Follow the instructions to apply migrations.

## Forward request information with a proxy or load balancer

If the app is deployed behind a proxy server or load balancer, some of the original request information might be forwarded to the app in request headers. This information usually includes the secure request scheme (`https`), host, and client IP address. Apps don't automatically read these request headers to discover and use the original request information.

The scheme is used in link generation that affects the authentication flow with external providers. Losing the secure scheme (`https`) results in the app generating incorrect insecure redirect URLs.

Use forwarded headers middleware to make the original request information available to the app for request processing.

For more information, see Configure ASP.NET Core to work with proxy servers and load balancers.

## Use Secret Manager to store tokens assigned by login providers

Social login providers assign **Application Id** and **Application Secret** tokens during the registration process. The exact token names vary by provider. These tokens represent the credentials that the app uses to access the provider's API. The tokens constitute *user secrets* that can be linked to your app configuration with the help of Secret Manager. User secrets are a more secure alternative to storing the tokens in a configuration file, such as `appsettings.json`.

Important

Secret Manager is only for local development and testing. Protect staging and production secrets with the Azure Key Vault configuration provider, which can also be used for local development and testing if you prefer not to use the Secret Manager locally.

For guidance on storing the tokens assigned by each login provider, see Safe storage of app secrets in development.

## Configure login providers

Use the following articles to configure login providers and the app:

- Facebook instructions
- X (formerly Twitter) instructions
- Google instructions
- Microsoft instructions
- Other provider instructions

## Multiple authentication providers

When the app requires multiple providers, chain the provider extension methods on AddAuthentication:

```
builder.Services.AddAuthentication()
 .AddGoogle(options =>
 {
 // Google configuration options
 })
 .AddFacebook(options =>
 {
 // Facebook configuration options
 })
 .AddMicrosoftAccount(options =>
 {
 // Microsoft Account configuration options
 })
 .AddTwitter(options =>
 {
 // X (formerly Twitter) configuration options
 });
```
For detailed configuration guidance on each provider, see their respective articles.

## Optionally set a password

When you register with an external login provider, you don't have a password registered with the app. This alleviates you from creating and remembering a password for the site, but it also makes you completely dependent on the external login provider for site access. If the external login provider is unavailable, you won't be able to sign in to the app.

To create a password and sign in using your email that you set during the sign-in process with external providers:

- Select the **Hello <email alias>**link at the top-right corner to navigate to the**Manage**view:
- Select **Create**:
- Set a valid password, and you can use this credential to sign in with your email address.

## Additional information

- Sign in with Apple Example Integration
- How to customize the login buttons (`dotnet/AspNetCore.Docs`#10563)
- Persist additional data about the user and their access and refresh tokens
- Passkeys (FIDO2/WebAuthn) in ASP.NET Core Identity: A passwordless alternative to social logins introduced in .NET 10.

ASP.NET Core

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

The following table provides a nonexhaustive overview of identity management solutions for ASP.NET Core apps in alphabetical order. These solutions offer features and capabilities to manage user authentication, authorization, and user identity, including options for apps that are:

- Container-based.
- Self-hosted, where you manage the app's installation and infrastructure on your own hardware.
- Managed in a cloud-based service, such as Microsoft Entra.

Depending on your company size and app requirements, many commercial licenses provide "community" or free options.

| Name | Type | License Type | Documentation |
|---|---|---|---|
| ASP.NET Core Identity | Self host | OSS (MIT) | Introduction to Identity on ASP.NET Core |
| Auth0 | Managed | Commercial | Get started |
| Duende IdentityServer | Self host | Commercial | ASP.NET Identity integration |
| Keycloak | Container | OSS (Apache 2.0) | Keycloak securing apps documentation |
| Microsoft Entra ID | Managed | Commercial | Entra documentation |
| Okta | Managed | Commercial | Okta for ASP.NET Core |
| OpenIddict | Self host | OSS (Apache 2.0) | OpenIddict documentation |

ASP.NET Core
