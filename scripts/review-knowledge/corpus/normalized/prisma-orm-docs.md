---
title: Prisma Client
description: Prisma Client is Prisma ORM's generated, type-safe query builder for Node.js, Bun, and Deno applications.
url: /orm/prisma-client
metaTitle: Prisma Client overview
metaDescription: Learn what Prisma Client is, how to generate it, and where to go next for querying, relations, and transactions.
---

Prisma Client is Prisma ORM's generated query builder. It is tailored to your schema, fully typed, and designed to make common database work feel like ordinary application code.

## What Prisma Client gives you

- Typed query methods based on your models
- Autocomplete for filters, relations, ordering, and nested writes
- Predictable plain JavaScript objects as query results
- A single client API that works across PostgreSQL, MySQL, SQLite, MongoDB, and more

## Quick start

### 1. Define a generator in your schema

```prisma title="schema.prisma"
generator client {
 provider = "prisma-client"
 output = "./generated"
}
```

### 2. Install Prisma Client

```npm
npm install @prisma/client
```

### 3. Generate the client

```npm
npx prisma generate
```

If you want more detail on this step, see [Generating Prisma Client](/orm/prisma-client/setup-and-configuration/generating-prisma-client).

### 4. Import and use the generated client

```ts
import { PrismaClient } from "./generated/client";

const prisma = new PrismaClient();

const users = await prisma.user.findMany();
```

## Common tasks

- [Set up and configure Prisma Client](/orm/prisma-client/setup-and-configuration/introduction)
- [Generate Prisma Client](/orm/prisma-client/setup-and-configuration/generating-prisma-client)
- [Run CRUD queries](/orm/prisma-client/queries/crud)
- [Work with relations](/orm/prisma-client/queries/relation-queries)
- [Use transactions](/orm/prisma-client/queries/transactions)
- [Use raw SQL when you need it](/orm/prisma-client/using-raw-sql)

## Related reference docs

- [Prisma Client API reference](/orm/reference/prisma-client-reference)
- [Prisma schema generators](/orm/prisma-schema/overview/generators)
- [Prisma CLI generate command](/cli/generate)

---
title: Add methods to Prisma Client
description: 'Extend the functionality of Prisma Client, client component'
url: /orm/prisma-client/client-extensions/client
metaTitle: 'Prisma Client extensions: client component'
metaDescription: 'Extend the functionality of Prisma Client, client component'
---

You can use the `client` [Prisma Client extensions](/orm/prisma-client/client-extensions) component to add top-level methods to Prisma Client.

## Extend Prisma Client

Use the `$extends` [client-level method](/orm/reference/prisma-client-reference#client-methods) to create an _extended client_. An extended client is a variant of the standard Prisma Client that is wrapped by one or more extensions. Use the `client` extension component to add top-level methods to Prisma Client.

To add a top-level method to Prisma Client, use the following structure:

```ts
const prisma = new PrismaClient().$extends({
 client?: { ... }
})
```

### Example

The following example uses the `client` component to add two methods to Prisma Client:

- `$log` outputs a message.
- `$totalQueries` returns the number of queries executed by the current client instance.

```ts
let total = 0;
const prisma = new PrismaClient().$extends({
 client: {
 $log: (s: string) => console.log(s),
 async $totalQueries() {
 return total;
 },
 },
 query: {
 $allModels: {
 async $allOperations({ query, args }) {
 total += 1;
 return query(args);
 },
 },
 },
});

async function main() {
 prisma.$log("Hello world");
 const totalQueries = await prisma.$totalQueries();
 console.log(totalQueries);
}
```

---
title: Shared packages & examples
description: Explore the Prisma Client extensions that have been built by Prisma and its community
url: /orm/prisma-client/client-extensions/extension-examples
metaTitle: Prisma Client extensions | Shared packages & examples
metaDescription: Explore the Prisma Client extensions that have been built by Prisma and its community
---

## Extensions made by Prisma

The following is a list of extensions we've built at Prisma:

| Extension | Description |
| :------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------- |
| [`@prisma/extension-read-replicas`](https://github.com/prisma/extension-read-replicas) | Adds read replica support to Prisma Client |

## Extensions made by Prisma's community

The following is a list of extensions created by the community. If you want to create your own package, refer to the [Shared Prisma Client extensions](/orm/prisma-client/client-extensions/shared-extensions) documentation.

| Extension | Description |
| :--------------------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------------------------- |
| [`prisma-extension-supabase-rls`](https://github.com/dthyresson/prisma-extension-supabase-rls) | Adds support for Supabase Row Level Security with Prisma |
| [`prisma-extension-bark`](https://github.com/adamjkb/bark) | Implements the Materialized Path pattern that allows you to easily create and interact with tree structures in Prisma |
| [`prisma-cursorstream`](https://github.com/etabits/prisma-cursorstream) | Adds cursor-based streaming |
| [`prisma-gpt`](https://github.com/aliyeysides/prisma-gpt) | Lets you query your database using natural language |
| [`prisma-extension-caching`](https://github.com/isaev-the-poetry/prisma-extension-caching) | Adds the ability to cache complex queries |
| [`prisma-extension-cache-manager`](https://github.com/random42/prisma-extension-cache-manager) | Caches model queries with any [cache-manager](https://www.npmjs.com/package/cache-manager) compatible cache |
| [`prisma-extension-random`](https://github.com/nkeil/prisma-extension-random) | Lets you query for random rows in your database |
| [`prisma-paginate`](https://github.com/sandrewTx08/prisma-paginate) | Adds support for paginating read queries |
| [`prisma-extension-streamdal`](https://github.com/streamdal/prisma-extension-streamdal) | Adds support for Code-Native data pipelines using Streamdal |
| [`prisma-rbac`](https://github.com/multipliedtwice/prisma-rbac) | Adds customizable role-based access control |
| [`prisma-extension-redis`](https://github.com/yxx4c/prisma-extension-redis) | Extensive Prisma extension designed for efficient caching and cache invalidation using Redis and Dragonfly Databases |
| [`prisma-cache-extension`](https://github.com/Shikhar97/prisma-cache) | Prisma extension for caching and invalidating cache with Redis(other Storage options to be supported) |
| [`prisma-extension-casl`](https://github.com/dennemark/prisma-extension-casl) | Prisma client extension that utilizes CASL to enforce authorization logic on most simple and nested queries. |
| [`prisma-emitter-extension`](https://github.com/feggaa/prisma-emitter-extension) | Prisma extension for emit events on CRUD operations based on configurable listeners. |

If you have built an extension and would like to see it featured, feel free to add it to the list by opening a pull request.

## Examples

:::info

The following example extensions are provided as examples only, and without warranty. They are supposed to show how Prisma Client extensions can be created using approaches documented here. We recommend using these examples as a source of inspiration for building your own extensions.

:::

| Example | Description |
| :------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------ |
| [`audit-log-context`](https://github.com/prisma/prisma-client-extensions/tree/main/audit-log-context) | Provides the current user's ID as context to Postgres audit log triggers |
| [`callback-free-itx`](https://github.com/prisma/prisma-client-extensions/tree/main/callback-free-itx) | Adds a method to start interactive transactions without callbacks |
| [`computed-fields`](https://github.com/prisma/prisma-client-extensions/tree/main/computed-fields) | Adds virtual / computed fields to result objects |
| [`input-transformation`](https://github.com/prisma/prisma-client-extensions/tree/main/input-transformation) | Transforms the input arguments passed to Prisma Client queries to filter the result set |
| [`input-validation`](https://github.com/prisma/prisma-client-extensions/tree/main/input-validation) | Runs custom validation logic on input arguments passed to mutation methods |
| [`instance-methods`](https://github.com/prisma/prisma-client-extensions/tree/main/instance-methods) | Adds Active Record-like methods like `save()` and `delete()` to result objects |
| [`json-field-types`](https://github.com/prisma/prisma-client-extensions/tree/main/json-field-types) | Uses strongly-typed runtime parsing for data stored in JSON columns |
| [`model-filters`](https://github.com/prisma/prisma-client-extensions/tree/main/model-filters) | Adds reusable filters that can composed into complex `where` conditions for a model |
| [`obfuscated-fields`](https://github.com/prisma/prisma-client-extensions/tree/main/obfuscated-fields) | Prevents sensitive data (e.g. `password` fields) from being included in results |
| [`query-logging`](https://github.com/prisma/prisma-client-extensions/tree/main/query-logging) | Wraps Prisma Client queries with simple query timing and logging |
| [`readonly-client`](https://github.com/prisma/prisma-client-extensions/tree/main/readonly-client) | Creates a client that only allows read operations |
| [`retry-transactions`](https://github.com/prisma/prisma-client-extensions/tree/main/retry-transactions) | Adds a retry mechanism to transactions with exponential backoff and jitter |
| [`row-level-security`](https://github.com/prisma/prisma-client-extensions/tree/main/row-level-security) | Uses Postgres row-level security policies to isolate data a multi-tenant application |
| [`static-methods`](https://github.com/prisma/prisma-client-extensions/tree/main/static-methods) | Adds custom query methods to Prisma Client models |
| [`transformed-fields`](https://github.com/prisma/prisma-client-extensions/tree/main/transformed-fields) | Demonstrates how to use result extensions to transform query results and add i18n to an app |
| [`exists-method`](https://github.com/prisma/prisma-client-extensions/tree/main/exists-fn) | Demonstrates how to add an `exists` method to all your models |
| [`update-delete-ignore-not-found `](https://github.com/prisma/prisma-client-extensions/tree/main/update-delete-ignore-not-found) | Demonstrates how to add the `updateIgnoreOnNotFound` and `deleteIgnoreOnNotFound` methods to all your models. |

## Going further

- Learn more about [Prisma Client extensions](/orm/prisma-client/client-extensions).

---
title: What are Client Extensions
description: Extend the functionality of Prisma Client
url: /orm/prisma-client/client-extensions
metaTitle: Prisma Client extensions
metaDescription: Extend the functionality of Prisma Client
---

You can use Prisma Client extensions to add functionality to your models, result objects, and queries, or to add client-level methods.

You can create an extension with one or more of the following component types:

- `model`: [add custom methods or fields to your models](/orm/prisma-client/client-extensions/model)
- `client`: [add client-level methods to Prisma Client](/orm/prisma-client/client-extensions/client)
- `query`: [create custom Prisma Client queries](/orm/prisma-client/client-extensions/query)
- `result`: [add custom fields to your query results](/orm/prisma-client/client-extensions/result)

For example, you might create an extension that uses the `model` and `client` component types.

## About Prisma Client extensions

When you use a Prisma Client extension, you create an _extended client_. An extended client is a lightweight variant of the standard Prisma Client that is wrapped by one or more extensions. The standard client is not mutated. You can add as many extended clients as you want to your project. [Learn more about extended clients](#extended-clients).

You can associate a single extension, or multiple extensions, with an extended client. [Learn more about multiple extensions](#multiple-extensions).

You can [share your Prisma Client extensions](/orm/prisma-client/client-extensions/shared-extensions) with other Prisma ORM users, and [import Prisma Client extensions developed by other users](/orm/prisma-client/client-extensions/shared-extensions#install-a-shared-packaged-extension) into your Prisma ORM project.

### Extended clients

Extended clients interact with each other, and with the standard client, as follows:

- Each extended client operates independently in an isolated instance.
- Extended clients cannot conflict with each other, or with the standard client.
- All extended clients and the standard client share the same connection pool.

> **Note**: The author of an extension can modify this behavior since they're able to run arbitrary code as part of an extension. For example, an extension might actually create an entirely new `PrismaClient` instance (including its own query engine and connection pool). Be sure to check the documentation of the extension you're using to learn about any specific behavior it might implement.

### Example use cases for extended clients

Because extended clients operate in isolated instances, they can be a good way to do the following, for example:

- Implement row-level security (RLS), where each HTTP request has its own client with its own RLS extension, customized with session data. This can keep each user entirely separate, each in a separate client.
- Add a `user.current()` method for the `User` model to get the currently logged-in user.
- Enable more verbose logging for requests if a debug cookie is set.
- Attach a unique request id to all logs so that you can correlate them later, for example to help you analyze the operations that Prisma Client carries out.
- Remove a `delete` method from models unless the application calls the admin endpoint and the user has the necessary privileges.

## Add an extension to Prisma Client

You can create an extension using two primary ways:

- Use the client-level [`$extends`](/orm/reference/prisma-client-reference#client-methods) method

 ```ts
 const prisma = new PrismaClient().$extends({
 name: 'signUp', // Optional: name appears in error logs
 model: { // This is a `model` component
 user: { ... } // The extension logic for the `user` model goes inside the curly braces
 },
 })
 ```

- Use the `Prisma.defineExtension` method to define an extension and assign it to a variable, and then pass the extension to the client-level `$extends` method

 ```ts
 import { Prisma } from '@prisma/client'

 // Define the extension
 const myExtension = Prisma.defineExtension({
 name: 'signUp', // Optional: name appears in error logs
 model: { // This is a `model` component
 user: { ... } // The extension logic for the `user` model goes inside the curly braces
 },
 })

 // Pass the extension to a Prisma Client instance
 const prisma = new PrismaClient().$extends(myExtension)
 ```

 :::tip

 This pattern is useful for when you would like to separate extensions into multiple files or directories within a project.

 :::

The above examples use the [`model` extension component](/orm/prisma-client/client-extensions/model) to extend the `User` model.

In your `$extends` method, use the appropriate extension component or components ([`model`](/orm/prisma-client/client-extensions/model), [`client`](/orm/prisma-client/client-extensions/client), [`result`](/orm/prisma-client/client-extensions/result) or [`query`](/orm/prisma-client/client-extensions/query)).

## Name an extension for error logs

You can name your extensions to help identify them in error logs. To do so, use the optional field `name`. For example:

```ts
const prisma = new PrismaClient().$extends({
 name: `signUp`, // (Optional) Extension name
 model: {
 user: { ... }
 },
})
```

## Multiple extensions

You can associate an extension with an [extended client](#about-prisma-client-extensions) in one of two ways:

- You can associate it with an extended client on its own, or
- You can combine the extension with other extensions and associate all of these extensions with an extended client. The functionality from these combined extensions applies to the same extended client.
 Note: [Combined extensions can conflict](#conflicts-in-combined-extensions).

You can combine the two approaches above. For example, you might associate one extension with its own extended client and associate two other extensions with another extended client. [Learn more about how client instances interact](#extended-clients).

### Apply multiple extensions to an extended client

In the following example, suppose that you have two extensions, `extensionA` and `extensionB`. There are two ways to combine these.

#### Option 1: Declare the new client in one line

With this option, you apply both extensions to a new client in one line of code.

```ts
// First of all, store your original Prisma Client in a variable as usual
const prisma = new PrismaClient();

// Declare an extended client that has an extensionA and extensionB
const prismaAB = prisma.$extends(extensionA).$extends(extensionB);
```

You can then refer to `prismaAB` in your code, for example `prismaAB.myExtensionMethod()`.

#### Option 2: Declare multiple extended clients

The advantage of this option is that you can call any of the extended clients separately.

```ts
// First of all, store your original Prisma Client in a variable as usual
const prisma = new PrismaClient();

// Declare an extended client that has extensionA applied
const prismaA = prisma.$extends(extensionA);

// Declare an extended client that has extensionB applied
const prismaB = prisma.$extends(extensionB);

// Declare an extended client that is a combination of clientA and clientB
const prismaAB = prismaA.$extends(extensionB);
```

In your code, you can call any of these clients separately, for example `prismaA.myExtensionMethod()`, `prismaB.myExtensionMethod()`, or `prismaAB.myExtensionMethod()`.

### Conflicts in combined extensions

When you combine two or more extensions into a single extended client, then the _last_ extension that you declare takes precedence in any conflict. In the example in option 1 above, suppose there is a method called `myExtensionMethod()` defined in `extensionA` and a method called `myExtensionMethod()` in `extensionB`. When you call `prismaAB.myExtensionMethod()`, then Prisma Client uses `myExtensionMethod()` as defined in `extensionB`.

### Middleware chaining with query extensions

Chain [`query`](/orm/prisma-client/client-extensions/query) extensions to compose middleware. Extensions execute in order—first in, first out:

```ts
import { PrismaClient } from "./generated/prisma";

const prisma = new PrismaClient()
 // Extension 1: Logging - measures query execution time
 .$extends({
 query: {
 $allModels: {
 async $allOperations({ model, operation, args, query }) {
 const start = Date.now();
 const result = await query(args);
 console.log(`[LOGGING] ${model}.${operation}: ${Date.now() - start}ms`);
 return result;
 },
 },
 },
 })
 // Extension 2: Audit - logs write operations
 .$extends({
 query: {
 $allModels: {
 async $allOperations({ model, operation, args, query }) {
 if (["create", "update", "delete"].includes(operation)) {
 console.log(`[AUDIT] ${operation} on ${model}:`, JSON.stringify(args));
 }
 return query(args);
 },
 },
 },
 });
```

## Type of an extended client

You can infer the type of an extended Prisma Client instance using the [`typeof`](https://www.typescriptlang.org/docs/handbook/2/typeof-types.html) utility as follows:

```ts
const extendedPrismaClient = new PrismaClient().$extends({
 /** extension */
});

type ExtendedPrismaClient = typeof extendedPrismaClient;
```

If you're using Prisma Client as a singleton, you can get the type of the extended Prisma Client instance using the `typeof` and [`ReturnType`](https://www.typescriptlang.org/docs/handbook/utility-types.html#returntypetype) utilities as follows:

```ts
function getExtendedClient() {
 return new PrismaClient().$extends({
 /* extension */
 });
}

type ExtendedPrismaClient = ReturnType<typeof getExtendedClient>;
```

## Extending model types with `Prisma.Result`

You can use the `Prisma.Result` type utility to extend model types to include properties added via client extensions. This allows you to infer the type of the extended model, including the extended properties.

### Example

The following example demonstrates how to use `Prisma.Result` to extend the `User` model type to include a `__typename` property added via a client extension.

```ts
import { PrismaClient, Prisma } from "@prisma/client";

const prisma = new PrismaClient().$extends({
 result: {
 user: {
 __typename: {
 needs: {},
 compute() {
 return "User";
 },
 },
 },
 },
});

type ExtendedUser = Prisma.Result<typeof prisma.user, { select: { id: true } }, "findFirstOrThrow">;

async function main() {
 const user: ExtendedUser = await prisma.user.findFirstOrThrow({
 select: {
 id: true,
 __typename: true,
 },
 });

 console.log(user.__typename); // Output: 'User'
}

main();
```

The `Prisma.Result` type utility is used to infer the type of the extended `User` model, including the `__typename` property added via the client extension.

## Limitations

### Usage of client-level methods in extended clients

[Client-level methods](/orm/reference/prisma-client-reference#client-methods) do not necessarily exist on extended clients. For these clients you will need to first check for existence before using.

```ts
const xPrisma = new PrismaClient().$extends(...);

if (xPrisma.$connect) {
 xPrisma.$connect()
}
```

### Usage with nested operations

The `query` extension type does not support nested read and write operations.

---
title: Add custom methods to your models
description: 'Extend the functionality of Prisma Client, model component'
url: /orm/prisma-client/client-extensions/model
metaTitle: 'Prisma Client extensions: model component'
metaDescription: 'Extend the functionality of Prisma Client, model component'
---

You can use the `model` [Prisma Client extensions](/orm/prisma-client/client-extensions) component type to add custom methods to your models.

Possible uses for the `model` component include the following:

- New operations to operate alongside existing Prisma Client operations, such as `findMany`
- Encapsulated business logic
- Repetitive operations
- Model-specific utilities

## Add a custom method

Use the `$extends` [client-level method](/orm/reference/prisma-client-reference#client-methods) to create an _extended client_. An extended client is a variant of the standard Prisma Client that is wrapped by one or more extensions. Use the `model` extension component to add methods to models in your schema.

### Add a custom method to a specific model

To extend a specific model in your schema, use the following structure. This example adds a method to the `user` model.

```ts
const prisma = new PrismaClient().$extends({
 name?: '<name>', // (optional) names the extension for error logs
 model?: {
 user: { ... } // in this case, we extend the `user` model
 },
});
```

#### Example

The following example adds a method called `signUp` to the `user` model. This method creates a new user with the specified email address:

```ts
const prisma = new PrismaClient().$extends({
 model: {
 user: {
 async signUp(email: string) {
 await prisma.user.create({ data: { email } });
 },
 },
 },
});
```

You would call `signUp` in your application as follows:

```ts
const user = await prisma.user.signUp("john@prisma.io");
```

### Add a custom method to all models in your schema

To extend _all_ models in your schema, use the following structure:

```ts
const prisma = new PrismaClient().$extends({
 name?: '<name>', // `name` is an optional field that you can use to name the extension for error logs
 model?: {
 $allModels: { ... }
 },
})
```

#### Example

The following example adds an `exists` method to all models.

```ts
const prisma = new PrismaClient().$extends({
 model: {
 $allModels: {
 async exists<T>(this: T, where: Prisma.Args<T, "findFirst">["where"]): Promise<boolean> {
 // Get the current model at runtime
 const context = Prisma.getExtensionContext(this);

 const result = await (context as any).findFirst({ where });
 return result !== null;
 },
 },
 },
});
```

You would call `exists` in your application as follows:

```ts
// `exists` method available on all models
await prisma.user.exists({ name: "Alice" });
await prisma.post.exists({
 OR: [{ title: { contains: "Prisma" } }, { content: { contains: "Prisma" } }],
});
```

## Call a custom method from another custom method

You can call a custom method from another custom method, if the two methods are declared on the same model. For example, you can call a custom method on the `user` model from another custom method on the `user` model. It does not matter if the two methods are declared in the same extension or in different extensions.

To do so, use `Prisma.getExtensionContext(this).methodName`. Note that you cannot use `prisma.user.methodName`. This is because `prisma` is not extended yet, and therefore does not contain the new method.

For example:

```ts
const prisma = new PrismaClient().$extends({
 model: {
 user: {
 firstMethod() {
 ...
 },
 secondMethod() {
 Prisma.getExtensionContext(this).firstMethod()
 }
 }
 }
})
```

## Get the current model name at runtime

You can get the name of the current model at runtime with `Prisma.getExtensionContext(this).$name`. You might use this to write out the model name to a log, to send the name to another service, or to branch your code based on the model.

For example:

```ts
// `context` refers to the current model
const context = Prisma.getExtensionContext(this);

// `context.$name` returns the name of the current model
console.log(context.$name);

// Usage
await (context as any).findFirst({ args });
```

Refer to [Add a custom method to all models in your schema](#example-1) for a concrete example for retrieving the current model name at runtime.

## Advanced type safety: type utilities for defining generic extensions

You can improve the type-safety of `model` components in your shared extensions with [type utilities](/orm/prisma-client/client-extensions/type-utilities).

---
title: Create custom Prisma Client queries
description: 'Extend the functionality of Prisma Client, query component'
url: /orm/prisma-client/client-extensions/query
metaTitle: 'Prisma Client extensions: query component'
metaDescription: 'Extend the functionality of Prisma Client, query component'
---

You can use the `query` [Prisma Client extensions](/orm/prisma-client/client-extensions) component type to hook into the query life-cycle and modify an incoming query or its result.

You can use Prisma Client extensions `query` component to create independent clients with customized behavior. You can bind one client to a specific filter or user, and another client to another filter or user. For example, you might do this to get [user isolation](/orm/prisma-client/client-extensions#extended-clients) in a row-level security (RLS) extension. The `query` extension component provides end-to-end type safety for all your custom queries.

## Extend Prisma Client query operations

Use the `$extends` [client-level method](/orm/reference/prisma-client-reference#client-methods) to create an [extended client](/orm/prisma-client/client-extensions#about-prisma-client-extensions). An extended client is a variant of the standard Prisma Client that is wrapped by one or more extensions.

Use the `query` extension component to modify queries. You can modify a custom query in the following:

- [A specific operation in a specific model](#modify-a-specific-operation-in-a-specific-model)
- [A specific operation in all models of your schema](#modify-a-specific-operation-in-all-models-of-your-schema)
- [All Prisma Client operations](#modify-all-prisma-client-operations)
- [All operations in a specific model](#modify-all-operations-in-a-specific-model)
- [All operations in all models of your schema](#modify-all-operations-in-all-models-of-your-schema)
- [A specific top-level raw query operation](#modify-a-top-level-raw-query-operation)

To create a custom query, use the following structure:

```ts
const prisma = new PrismaClient().$extends({
 name?: 'name',
 query?: {
 user: { ... } // in this case, we add a query to the `user` model
 },
});
```

The properties are as follows:

- `name`: (optional) specifies a name for the extension that appears in error logs.
- `query`: defines a custom query.

### Modify a specific operation in a specific model

The `query` object can contain functions that map to the names of the [Prisma Client operations](/orm/reference/prisma-client-reference#model-queries), such as `findUnique()`, `findFirst`, `findMany`, `count`, and `create`. The following example modifies `user.findMany` to a use a customized query that finds only users who are older than 18 years:

```ts
const prisma = new PrismaClient().$extends({
 query: {
 user: {
 async findMany({ model, operation, args, query }) {
 // take incoming `where` and set `age`
 args.where = { ...args.where, age: { gt: 18 } };

 return query(args);
 },
 },
 },
});

await prisma.user.findMany(); // returns users whose age is greater than 18
```

In the above example, a call to `prisma.user.findMany` triggers `query.user.findMany`. Each callback receives a type-safe `{ model, operation, args, query }` object that describes the query. This object has the following properties:

- `model`: the name of the containing model for the query that we want to extend.

 In the above example, the `model` is a string of type `"User"`.

- `operation`: the name of the operation being extended and executed.

 In the above example, the `operation` is a string of type `"findMany"`.

- `args`: the specific query input information to be extended.

 This is a type-safe object that you can mutate before the query happens. You can mutate any of the properties in `args`. Exception: you cannot mutate `include` or `select` because that would change the expected output type and break type safety.

- `query`: a promise for the result of the query.
 - You can use `await` and then mutate the result of this promise, because its value is type-safe. TypeScript catches any unsafe mutations on the object.

### Modify a specific operation in all models of your schema

To extend the queries in all the models of your schema, use `$allModels` instead of a specific model name. For example:

```ts
const prisma = new PrismaClient().$extends({
 query: {
 $allModels: {
 async findMany({ model, operation, args, query }) {
 // set `take` and fill with the rest of `args`
 args = { ...args, take: 100 };

 return query(args);
 },
 },
 },
});
```

### Modify all operations in a specific model

Use `$allOperations` to extend all operations in a specific model.

For example, the following code applies a custom query to all operations on the `user` model:

```ts
const prisma = new PrismaClient().$extends({
 query: {
 user: {
 $allOperations({ model, operation, args, query }) {
 /* your custom logic here */
 return query(args);
 },
 },
 },
});
```

### Modify all Prisma Client operations

Use the `$allOperations` method to modify all query methods present in Prisma Client. The `$allOperations` can be used on both model operations and raw queries.

You can modify all methods as follows:

```ts
const prisma = new PrismaClient().$extends({
 query: {
 $allOperations({ model, operation, args, query }) {
 /* your custom logic for modifying all Prisma Client operations here */
 return query(args);
 },
 },
});
```

In the event a [raw query](/orm/prisma-client/using-raw-sql/raw-queries) is invoked, the `model` argument passed to the callback will be `undefined`.

For example, you can use the `$allOperations` method to log queries as follows:

```ts
const prisma = new PrismaClient().$extends({
 query: {
 async $allOperations({ operation, model, args, query }) {
 const start = performance.now();
 const result = await query(args);
 const end = performance.now();
 const time = end - start;
 console.log(
 util.inspect(
 { model, operation, args, time },
 { showHidden: false, depth: null, colors: true },
 ),
 );
 return result;
 },
 },
});
```

### Modify all operations in all models of your schema

Use `$allModels` and `$allOperations` to extend all operations in all models of your schema.

To apply a custom query to all operations on all models of your schema:

```ts
const prisma = new PrismaClient().$extends({
 query: {
 $allModels: {
 $allOperations({ model, operation, args, query }) {
 /* your custom logic for modifying all operations on all models here */
 return query(args);
 },
 },
 },
});
```

### Modify a top-level raw query operation

To apply custom behavior to a specific top-level raw query operation, use the name of a top-level raw query function instead of a model name:

```ts copy tab="Relational databases"
const prisma = new PrismaClient().$extends({
 query: {
 $queryRaw({ args, query, operation }) {
 // handle $queryRaw operation
 return query(args);
 },
 $executeRaw({ args, query, operation }) {
 // handle $executeRaw operation
 return query(args);
 },
 $queryRawUnsafe({ args, query, operation }) {
 // handle $queryRawUnsafe operation
 return query(args);
 },
 $executeRawUnsafe({ args, query, operation }) {
 // handle $executeRawUnsafe operation
 return query(args);
 },
 },
});
```

```ts copy tab="MongoDB"
const prisma = new PrismaClient().$extends({
 query: {
 $runCommandRaw({ args, query, operation }) {
 // handle $runCommandRaw operation
 return query(args);
 },
 },
});
```

### Mutate the result of a query

You can use `await` and then mutate the result of the `query` promise.

```ts
const prisma = new PrismaClient().$extends({
 query: {
 user: {
 async findFirst({ model, operation, args, query }) {
 const user = await query(args);

 if (user.password !== undefined) {
 user.password = "******";
 }

 return user;
 },
 },
 },
});
```

:::info

We include the above example to show that this is possible. However, for performance reasons we recommend that you use the [`result` component type](/orm/prisma-client/client-extensions/result) to override existing fields. The `result` component type usually gives better performance in this situation because it computes only on access. The `query` component type computes after query execution.

:::

## Wrap a query into a batch transaction

You can wrap your extended queries into a [batch transaction](/orm/prisma-client/queries/transactions). For example, you can use this to enact row-level security (RLS).

The following example extends `findFirst` so that it runs in a batch transaction.

```ts
const transactionExtension = Prisma.defineExtension((prisma) =>
 prisma.$extends({
 query: {
 user: {
 // Get the input `args` and a callback to `query`
 async findFirst({ args, query, operation }) {
 const [result] = await prisma.$transaction([query(args)]); // wrap the query in a batch transaction, and destructure the result to return an array
 return result; // return the first result found in the array
 },
 },
 },
 }),
);
const prisma = new PrismaClient().$extends(transactionExtension);
```

---
title: Add custom fields and methods to query results
description: 'Extend the functionality of Prisma Client, result component'
url: /orm/prisma-client/client-extensions/result
metaTitle: 'Prisma Client extensions: result component'
metaDescription: 'Extend the functionality of Prisma Client, result component'
---

You can use the `result` [Prisma Client extensions](/orm/prisma-client/client-extensions) component type to add custom fields and methods to query results.

Use the `$extends` [client-level method](/orm/reference/prisma-client-reference#client-methods) to create an _extended client_. An extended client is a variant of the standard Prisma Client that is wrapped by one or more extensions.

To add a custom [field](#add-a-custom-field-to-query-results) or [method](#add-a-custom-method-to-the-result-object) to query results, use the following structure. In this example, we add the custom field `myComputedField` to the result of a `user` model query.

```ts
const prisma = new PrismaClient().$extends({
 name?: 'name',
 result?: {
 user: { // in this case, we extend the `user` model
 myComputedField: { // the name of the new computed field
 needs: { ... },
 compute() { ... }
 },
 },
 },
});
```

The parameters are as follows:

- `name`: (optional) specifies a name for the extension that appears in error logs.
- `result`: defines new fields and methods to the query results.
- `needs`: an object which describes the dependencies of the result field.
- `compute`: a method that defines how the virtual field is computed when it is accessed.

## Add a custom field to query results

You can use the `result` extension component to add fields to query results. These fields are computed at runtime and are type-safe.

In the following example, we add a new virtual field called `fullName` to the `user` model.

```ts
const prisma = new PrismaClient().$extends({
 result: {
 user: {
 fullName: {
 // the dependencies
 needs: { firstName: true, lastName: true },
 compute(user) {
 // the computation logic
 return `${user.firstName} ${user.lastName}`;
 },
 },
 },
 },
});

const user = await prisma.user.findFirst();

// return the user's full name, such as "John Doe"
console.log(user.fullName);
```

In above example, the input `user` of `compute` is automatically typed according to the object defined in `needs`. `firstName` and `lastName` are of type `string`, because they are specified in `needs`. If they are not specified in `needs`, then they cannot be accessed.

## Re-use a computed field in another computed field

The following example computes a user's title and full name in a type-safe way. `titleFullName` is a computed field that reuses the `fullName` computed field.

```ts
const prisma = new PrismaClient()
 .$extends({
 result: {
 user: {
 fullName: {
 needs: { firstName: true, lastName: true },
 compute(user) {
 return `${user.firstName} ${user.lastName}`;
 },
 },
 },
 },
 })
 .$extends({
 result: {
 user: {
 titleFullName: {
 needs: { title: true, fullName: true },
 compute(user) {
 return `${user.title} (${user.fullName})`;
 },
 },
 },
 },
 });
```

### Considerations for fields

- For performance reasons, Prisma Client computes results on access, not on retrieval.
- You can only create computed fields that are based on scalar fields.
- You can only use computed fields with `select` and you cannot aggregate them. For example:

 ```ts
 const user = await prisma.user.findFirst({
 select: { email: true },
 });
 console.log(user.fullName); // undefined
 ```

## Add a custom method to the result object

You can use the `result` component to add methods to query results. The following example adds a new method, `save` to the result object.

```ts
const prisma = new PrismaClient().$extends({
 result: {
 user: {
 save: {
 needs: { id: true },
 compute(user) {
 return () => prisma.user.update({ where: { id: user.id }, data: user });
 },
 },
 },
 },
});

const user = await prisma.user.findUniqueOrThrow({ where: { id: someId } });
user.email = "mynewmail@mailservice.com";
await user.save();
```

## Using `omit` query option with `result` extension component

You can use the [`omit` (Preview) option](/orm/reference/prisma-client-reference#omit) with [custom fields](#add-a-custom-field-to-query-results) and fields needed by custom fields.

### `omit` fields needed by custom fields from query result

If you `omit` a field that is a dependency of a custom field, it will still be read from the database even though it will not be included in the query result.

The following example omits the `password` field, which is a dependency of the custom field `sanitizedPassword`:

```ts
const xprisma = prisma.$extends({
 result: {
 user: {
 sanitizedPassword: {
 needs: { password: true },
 compute(user) {
 return sanitize(user.password);
 },
 },
 },
 },
});

const user = await xprisma.user.findFirstOrThrow({
 omit: {
 password: true,
 },
});
```

In this case, although `password` is omitted from the result, it will still be queried from the database because it is a dependency of the `sanitizedPassword` custom field.

### `omit` custom field and dependencies from query result

To ensure omitted fields are not queried from the database at all, you must omit both the custom field and its dependencies.

The following example omits both the custom field `sanitizedPassword` and the dependent `password` field:

```ts
const xprisma = prisma.$extends({
 result: {
 user: {
 sanitizedPassword: {
 needs: { password: true },
 compute(user) {
 return sanitize(user.password);
 },
 },
 },
 },
});

const user = await xprisma.user.findFirstOrThrow({
 omit: {
 sanitizedPassword: true,
 password: true,
 },
});
```

In this case, omitting both `password` and `sanitizedPassword` will exclude both from the result as well as prevent the `password` field from being read from the database.

## Limitation

As of now, Prisma Client's result extension component does not support relation fields. This means that you cannot create custom fields or methods based on related models or fields in a relational relationship (e.g., user.posts, post.author). The needs parameter can only reference scalar fields within the same model. Follow [issue #20091 on GitHub](https://github.com/prisma/prisma/issues/20091).

```ts
const prisma = new PrismaClient().$extends({
 result: {
 user: {
 postsCount: {
 needs: { posts: true }, // This will not work because posts is a relation field
 compute(user) {
 return user.posts.length; // Accessing a relation is not allowed
 },
 },
 },
 },
});
```

---
title: Type utilities
description: 'Advanced type safety: improve type safety in your custom model methods'
url: /orm/prisma-client/client-extensions/type-utilities
metaTitle: 'Prisma Client Extensions: Type utilities'
metaDescription: 'Advanced type safety: improve type safety in your custom model methods'
---

Several type utilities exist within Prisma Client that can assist in the creation of highly type-safe extensions.

## Type Utilities

[Prisma Client type utilities](/orm/prisma-client/type-safety) are utilities available within your application and Prisma Client extensions and provide useful ways of constructing safe and extendable types for your extension.

The type utilities available are:

- `Exact<Input, Shape>`: Enforces strict type safety on `Input`. `Exact` makes sure that a generic type `Input` strictly complies with the type that you specify in `Shape`. It [narrows](https://www.typescriptlang.org/docs/handbook/2/narrowing.html) `Input` down to the most precise types.
- `Args<Type, Operation>`: Retrieves the input arguments for any given model and operation. This is particularly useful for extension authors who want to do the following:
 - Re-use existing types to extend or modify them.
 - Benefit from the same auto-completion experience as on existing operations.
- `Result<Type, Arguments, Operation>`: Takes the input arguments and provides the result for a given model and operation. You would usually use this in conjunction with `Args`. As with `Args`, `Result` helps you to re-use existing types to extend or modify them.
- `Payload<Type, Operation>`: Retrieves the entire structure of the result, as scalars and relations objects for a given model and operation. For example, you can use this to determine which keys are scalars or objects at a type level.

The following example creates a new operation, `exists`, based on `findFirst`. It has all of the arguments that `findFirst`.

```ts
const prisma = new PrismaClient().$extends({
 model: {
 $allModels: {
 // Define a new `exists` operation on all models
 // T is a generic type that corresponds to the current model
 async exists<T>(
 // `this` refers to the current type, e.g. `prisma.user` at runtime
 this: T,

 // The `exists` function will use the `where` arguments from the current model, `T`, and the `findFirst` operation
 where: Prisma.Args<T, "findFirst">["where"],
 ): Promise<boolean> {
 // Retrieve the current model at runtime
 const context = Prisma.getExtensionContext(this);

 // Prisma Client query that retrieves data based
 const result = await (context as any).findFirst({ where });
 return result !== null;
 },
 },
 },
});

async function main() {
 const user = await prisma.user.exists({ name: "Alice" });
 const post = await prisma.post.exists({
 OR: [{ title: { contains: "Prisma" } }, { content: { contains: "Prisma" } }],
 });
}
```

## Add a custom property to a method

The following example illustrates how you can add custom arguments, to a method in an extension:

```ts highlight=16
type CacheStrategy = {
 swr: number;
 ttl: number;
};

const prisma = new PrismaClient().$extends({
 model: {
 $allModels: {
 findMany<T, A>(
 this: T,
 args: Prisma.Exact<
 A,
 // For the `findMany` method, use the arguments from model `T` and the `findMany` method
 // and intersect it with `CacheStrategy` as part of `findMany` arguments
 Prisma.Args<T, "findMany"> & CacheStrategy
 >,
 ): Prisma.Result<T, A, "findMany"> {
 // method implementation with the cache strategy
 },
 },
 },
});

async function main() {
 await prisma.post.findMany({
 cacheStrategy: {
 ttl: 360,
 swr: 60,
 },
 });
}
```

The example here is only conceptual. For the actual caching to work, you will have to implement the logic. If you're interested in a caching extension/ service, we recommend taking a look at [Prisma Accelerate](https://www.prisma.io/accelerate).

---
title: Debugging
description: This page explains how to enable debugging output for Prisma Client by setting the `DEBUG` environment variable
url: /orm/prisma-client/debugging-and-troubleshooting/debugging
metaTitle: Debugging (Reference)
metaDescription: This page explains how to enable debugging output for Prisma Client by setting the `DEBUG` environment variable.
---

You can enable debugging output in Prisma Client and Prisma CLI via the [`DEBUG`](/orm/reference/environment-variables-reference#debug) environment variable. It accepts two namespaces to print debugging output:

- `prisma:engine`: Prints relevant debug messages happening in a Prisma ORM [engine](https://github.com/prisma/prisma-engines/)
- `prisma:client`: Prints relevant debug messages happening in the Prisma Client runtime
- `prisma*`: Prints all debug messages from Prisma Client or CLI
- `*`: Prints all debug messages

:::info

Prisma Client can be configured to log warnings, errors and information related to queries sent to the database. See [Configuring logging](/orm/prisma-client/observability-and-logging/logging) for more information.

:::

## Setting the `DEBUG` environment variable

Here are examples for setting these debugging options in bash:

```bash
# enable only `prisma:engine`-level debugging output
export DEBUG="prisma:engine"

# enable only `prisma:client`-level debugging output
export DEBUG="prisma:client"

# enable both `prisma-client`- and `engine`-level debugging output
export DEBUG="prisma:client,prisma:engine"
```

To enable all `prisma` debugging options, set `DEBUG` to `prisma*`:

```bash
export DEBUG="prisma*"
```

On Windows, use `set` instead of `export`:

```bash
set DEBUG="prisma*"
```

To enable _all_ debugging options, set `DEBUG` to `*`:

```bash
export DEBUG="*"
```

---
title: Handling exceptions and errors
description: This page covers how to handle exceptions and errors
url: /orm/prisma-client/debugging-and-troubleshooting/handling-exceptions-and-errors
metaTitle: Handling exceptions and errors (Reference)
metaDescription: This page covers how to handle exceptions and errors
---

In order to handle different types of errors you can use `instanceof` to check what the error is and handle it accordingly.

The following example tries to create a user with an already existing email record. This will throw an error because the `email` field has the `@unique` attribute applied to it.

```prisma title="schema.prisma"
model User {
 id Int @id @default(autoincrement())
 email String @unique
 name String?
}
```

Use the `Prisma` namespace to access the error type. The [error code](/orm/reference/error-reference#error-codes) can then be checked and a message can be printed.

```ts
import { PrismaPg } from "@prisma/adapter-pg";
import { PrismaClient, Prisma } from "../generated/prisma/client";

const connectionString = `${process.env.DATABASE_URL}`;
const adapter = new PrismaPg({ connectionString });
const prisma = new PrismaClient({ adapter });

try {
 await client.user.create({ data: { email: "alreadyexisting@mail.com" } });
} catch (e) {
 if (e instanceof Prisma.PrismaClientKnownRequestError) {
 // The .code property can be accessed in a type-safe manner
 if (e.code === "P2002") {
 console.log(
 "There is a unique constraint violation, a new user cannot be created with this email",
 );
 }
 }
 throw e;
}
```

See [Errors reference](/orm/reference/error-reference) for a detailed breakdown of the different error types and their codes.

---
title: Caveats when deploying to AWS platforms
description: Known caveats when deploying to an AWS platform
url: /orm/prisma-client/deployment/caveats-when-deploying-to-aws-platforms
metaTitle: Caveats when deploying to AWS platforms
metaDescription: Known caveats when deploying to an AWS platform
---

The following describes some caveats you might face when deploying to different AWS platforms.

## AWS RDS Proxy

Prisma ORM is compatible with AWS RDS Proxy. However, there is no benefit in using it for connection pooling with Prisma ORM due to the way RDS Proxy pins connections:

> "Your connections to the proxy can enter a state known as pinning. When a connection is pinned, each later transaction uses the same underlying database connection until the session ends. Other client connections also can't reuse that database connection until the session ends. The session ends when Prisma Client's connection is dropped." - [AWS RDS Proxy Docs](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy-pinning.html)

[Prepared statements (of any size) or query statements greater than 16 KB cause RDS Proxy to pin the session.](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/rds-proxy-pinning.html) Because Prisma ORM uses prepared statements for all queries, you won't see any benefit when using RDS Proxy with Prisma ORM.

## AWS Elastic Beanstalk

AWS Elastic Beanstalk is a PaaS-like deployment service that abstracts away infrastructure and allows you to deploy applications to AWS quickly.

When deploying an app using Prisma Client to AWS Elastic Beanstalk, Prisma ORM generates the Prisma Client code into `node_modules`. This is typically done in a `postinstall` hook defined in a `package.json`.

Because Beanstalk limits the ability to write to the filesystem in the `postinstall` hook, you need to create an [`.npmrc`](https://docs.npmjs.com/cli/v6/configuring-npm/npmrc/) file in the root of your project and add the following configuration:

```yaml title=".npmrc" showLineNumbers
unsafe-perm=true
```

Enabling `unsafe-perm` forces _npm_ to run as _root_, avoiding the filesystem access problem, thereby allowing the `prisma generate` command in the `postinstall` hook to generate your code.

### Error: @prisma/client did not initialize yet

This error happens because AWS Elastic Beanstalk doesn't install `devDependencies`, which means that it doesn't pick up the Prisma CLI. To remedy this you can either:

1. Add the `prisma` CLI package to your `dependencies` instead of the `devDependencies`. (Making sure to run `npm install` afterward to update the `package-lock.json`).
2. Or install your `devDependencies` on AWS Elastic Beanstalk instances. To do this you must set the AWS Elastic Beanstalk `NPM_USE_PRODUCTION` environment property to false.

## AWS RDS Postgres

When using Prisma ORM with AWS RDS Postgres, you may encounter connection issues or the following error during migration or runtime:

```bash
Error: P1010: User <username> was denied access on the database <database>
```

### Cause

AWS RDS enforces SSL connections by default, and Prisma parses the database connection string with `rejectUnauthorized: true`, which requires a valid SSL certificate. If the certificate is not configured properly, Prisma cannot connect to the database.

### Solution

To resolve this issue, update the `DATABASE_URL` environment variable to include the `sslmode=no-verify` option. This bypasses strict SSL certificate verification and allows Prisma to connect to the database. Update your `.env` file as follows:

```bash
DATABASE_URL=postgresql://<username>:<password>@<host>/<database>?sslmode=no-verify&schema=public
```

### Why This Works

The `sslmode=no-verify` setting passes `rejectUnauthorized: false` to the SSL configuration via the [pg-connection-string](https://github.com/brianc/node-postgres/blob/95d7e620ef8b51743b4cbca05dd3c3ce858ecea7/packages/pg-connection-string/README.md?plain=1#L71) package. This disables strict certificate validation, allowing Prisma to establish a connection with the RDS database.

### Note

While using `sslmode=no-verify` can be a quick fix, it bypasses SSL verification and might not meet security requirements for production environments. In such cases, ensure that a valid SSL certificate is properly configured.

---
title: Deploying database changes with Prisma Migrate
description: Learn how to deploy database changes with Prisma Migrate
url: /orm/prisma-client/deployment/deploy-database-changes-with-prisma-migrate
metaDescription: Learn how to deploy database changes with Prisma Migrate.
metaTitle: Deploying database changes with Prisma Migrate
---

To apply pending migrations to staging, testing, or production environments, run the `migrate deploy` command as part of your CI/CD pipeline:

```npm
npx prisma migrate deploy
```

:::info

This guide **does not apply for MongoDB**.<br />
Instead of `migrate deploy`, [`db push`](/orm/prisma-migrate/workflows/prototyping-your-schema) is used for [MongoDB](/orm/core-concepts/supported-databases/mongodb).

:::

Exactly when to run `prisma migrate deploy` depends on your platform. For example, a simplified [Heroku](/orm/prisma-client/deployment/traditional/deploy-to-heroku) workflow includes:

1. Ensuring the `./prisma/migration` folder is in source control
2. Running `prisma migrate deploy` during the [release phase](https://devcenter.heroku.com/articles/release-phase)

Ideally, `migrate deploy` should be part of an automated CI/CD pipeline, and we do not generally recommend running this command locally to deploy changes to a production database (for example, by temporarily changing the `DATABASE_URL` environment variable). It is not generally considered good practice to store the production database URL locally.

Beware that in order to run the `prisma migrate deploy` command, you need access to the `prisma` dependency that is typically added to the `devDependencies`. Some platforms like Vercel, prune development dependencies during the build, thereby preventing you from calling the command. This can be worked around by making the `prisma` a production dependency, by moving it to `dependencies` in your `package.json`.
For more information about the `migrate deploy` command, see:

- [`migrate deploy` reference](/orm/reference/prisma-cli-reference#migrate-deploy)
- [How `migrate deploy` works](/orm/prisma-migrate/workflows/development-and-production#production-and-testing-environments)
- [Production troubleshooting](/orm/prisma-migrate/workflows/patching-and-hotfixing)

## Deploying database changes using GitHub Actions

As part of your CI/CD, you can run `prisma migrate deploy` as part of your pipeline to apply pending migrations to your production database.

Here is an example action that will run your migrations against your database:

```yaml title="deploy.yml" highlight=17-20 showLineNumbers
name: Deploy
on:
 push:
 paths:
 - prisma/migrations/** # [!code highlight]
 branches:
 - main

jobs:
 deploy:
 runs-on: ubuntu-latest
 steps:
 - name: Checkout repo
 uses: actions/checkout@v3
 - name: Setup Node
 uses: actions/setup-node@v3
 - name: Install dependencies
 run: npm install
 - name: Apply all pending migrations to the database
 run: npx prisma migrate deploy
 env:
 DATABASE_URL: ${{ secrets.DATABASE_URL }}
```

The highlighted line shows that this action will only run if there is a change in the `prisma/migrations` directory, so `npx prisma migrate deploy` will only run when migrations are updated.

Ensure you have the `DATABASE_URL` variable [set as a secret in your repository](https://docs.github.com/en/actions/security-for-github-actions/security-guides/using-secrets-in-github-actions), without quotes around the connection string.

## Pre-deploy migration safety checks

Before running `prisma migrate deploy`, you can analyze your migration SQL files for potentially dangerous patterns using a migration safety tool like [pgfence](/guides/integrations/pgfence). pgfence detects operations that acquire heavy locks (such as `CREATE INDEX` without `CONCURRENTLY` or `ALTER COLUMN TYPE`), reports risk levels, and provides safe rewrite recipes.

To add pgfence as a pre-deploy step in your GitHub Actions workflow:

```yaml
- name: Run migration safety check
 run: npx @flvmnt/pgfence analyze --ci --max-risk medium prisma/migrations/**/migration.sql

- name: Apply all pending migrations to the database
 run: npx prisma migrate deploy
 env:
 DATABASE_URL: ${{ secrets.DATABASE_URL }}
```

For a full setup guide, see the [pgfence integration guide](/guides/integrations/pgfence).

---
title: Deploy migrations from a local environment
description: Learn how to deploy Node.js and TypeScript applications that are using Prisma Client locally
url: /orm/prisma-client/deployment/deploy-migrations-from-a-local-environment
metaTitle: Deploy migrations from a local environment
metaDescription: Learn how to deploy Node.js and TypeScript applications that are using Prisma Client locally.
---

There are two scenarios where you might consider deploying migrations directly from a local environment to a production environment.

- You have a local CI/CD pipeline
- You are [baselining](/orm/prisma-migrate/workflows/baselining) a production environment

This page outlines some examples of how you can do that and **why we would generally not recommend it**.

## Local CI/CD pipeline

If you do not have an automated CI/CD process, you can technically deploy new migrations from your local environment to production in the following ways:

1. Make sure your migration history is up to date. You can do this through running `prisma migrate dev`, which will generate a migration history from the latest changes made.
2. Swap your local connection URL for your production connection URL

```bash title=".env"
DATABASE_URL="postgresql://johndoe:randompassword@localhost:5432/my_local_database" # [!code --]

DATABASE_URL="postgresql://johndoe:randompassword@prod-db.example.com:5432/my_production_database" # [!code ++]
```

3. Run `prisma migrate deploy`

**⛔ We strongly discourage this solution due to the following reasons**

- You risk exposing your production database connection URL to version control.
- You may accidentally use your production connection URL instead and in turn **override or delete your production database**.

**✅ We recommend setting up an automated CI/CD pipeline**

The pipeline should handle deployment to staging and production environments, and use `migrate deploy` in a pipeline step. See the [deployment guides](/orm/prisma-client/deployment/deploy-database-changes-with-prisma-migrate) for examples.

## Baselining a production database

When you add Prisma Migrate to an **existing database**, you must [baseline](/orm/prisma-migrate/workflows/baselining) the production database. Baselining is performed **once**, and can be done from a local instance.

![Baselining production from local with Prisma ORM](/img/orm/baseline-production-from-local.png)

---
title: Deploy Prisma ORM
description: Learn more about the different deployment paradigms for Node.js applications and how they affect deploying an application using Prisma Client
url: /orm/prisma-client/deployment/deploy-prisma
metaTitle: Deploying Prisma ORM-based projects
metaDescription: Learn more about the different deployment paradigms for Node.js applications and how they affect deploying an application using Prisma Client.
---

Projects using Prisma Client can be deployed to many different cloud platforms. Given the variety of cloud platforms and different names, it's noteworthy to mention the different deployment paradigms, as they affect the way you deploy an application using Prisma Client.

## Deployment paradigms

Each paradigm has different tradeoffs that affect the performance, scalability, and operational costs of your application.

Moreover, the user traffic pattern of your application is also an important factor to consider. For example, any application with consistent user traffic may be better suited for a [continuously running paradigm](#traditional-servers), whereas an application with sudden spikes may be better suited to [serverless](#serverless-functions).

### Traditional servers

Your application is [traditionally deployed](/orm/prisma-client/deployment/traditional/deploy-to-heroku) if a Node.js process is continuously running and handles multiple requests at the same time. Your application could be deployed to a Platform-as-a-Service (PaaS) like [Heroku](/orm/prisma-client/deployment/traditional/deploy-to-heroku), [Koyeb](/orm/prisma-client/deployment/traditional/deploy-to-koyeb), or [Render](/orm/prisma-client/deployment/traditional/deploy-to-render); as a Docker container to Kubernetes; or as a Node.js process on a virtual machine or bare metal server.

See also: [Connection management in long-running processes](/orm/prisma-client/setup-and-configuration/databases-connections#long-running-processes)

### Serverless Functions

Your application is [serverless](/orm/prisma-client/deployment/serverless/deploy-to-vercel) if the Node.js processes of your application (or subsets of it broken into functions) are started as requests come in, and each function only handles one request at a time. Your application would most likely be deployed to a Function-as-a-Service (FaaS) offering, such as [AWS Lambda](/orm/prisma-client/deployment/serverless/deploy-to-aws-lambda) or [Azure Functions](/orm/prisma-client/deployment/serverless/deploy-to-azure-functions)

Serverless environments have the concept of warm starts, which means that for subsequent invocations of the same function, it may use an already existing container that has the allocated processes, memory, file system (`/tmp` is writable on AWS Lambda), and even DB connection still available.

Typically, any piece of code [outside the handler](https://docs.aws.amazon.com/lambda/latest/dg/welcome.html) remains initialized.

See also: [Connection management in serverless environments](/orm/prisma-client/setup-and-configuration/databases-connections#serverless-environments-faas)

### Edge Functions

Your application is [edge deployed](/orm/prisma-client/deployment/edge/overview) if your application is [serverless](#serverless-functions) and the functions are distributed across one or more regions close to the user.

Typically, edge environments also have a different runtime than a traditional or serverless environment, leading to common APIs being unavailable.

---
title: Logging
description: Learn how to configure Prisma Client to log the raw SQL queries it sends to the database and other information
url: /orm/prisma-client/observability-and-logging/logging
metaTitle: Logging
metaDescription: Learn how to configure Prisma Client to log the raw SQL queries it sends to the database and other information.
---

Use the `PrismaClient` [`log`](/orm/reference/prisma-client-reference#log) parameter to configure [log levels](/orm/reference/prisma-client-reference#log-levels) , including warnings, errors, and information about the queries sent to the database.

Prisma Client supports two types of logging:

- Logging to [stdout](https://en.wikipedia.org/wiki/Standard_streams) (default)
- Event-based logging (use [`$on()`](/orm/reference/prisma-client-reference#on) method to [subscribe to events](#event-based-logging))

:::info

You can also use the `DEBUG` environment variable to enable debugging output in Prisma Client. See [Debugging](/orm/prisma-client/debugging-and-troubleshooting/debugging) for more information.

:::

:::info

If you want a detailed insight into your Prisma Client's performance at the level of individual operations, see [Tracing](/orm/prisma-client/observability-and-logging/opentelemetry-tracing).

:::

## Log to stdout

The simplest way to print _all_ log levels to stdout is to pass in an array `LogLevel` objects:

```ts
const prisma = new PrismaClient({
 log: ["query", "info", "warn", "error"],
});
```

This is the short form of passing in an array of `LogDefinition` objects where the value of `emit` is always `stdout`:

```ts
const prisma = new PrismaClient({
 log: [
 {
 emit: "stdout",
 level: "query",
 },
 {
 emit: "stdout",
 level: "error",
 },
 {
 emit: "stdout",
 level: "info",
 },
 {
 emit: "stdout",
 level: "warn",
 },
 ],
});
```

## Event-based logging

To use event-based logging:

1. Set `emit` to `event` for a specific log level, such as query
2. Use the `$on()` method to subscribe to the event

The following example subscribes to all `query` events and write the `duration` and `query` to console:

```ts highlight=4,5,22-26;normal tab="Relational databases"
const prisma = new PrismaClient({
 log: [
 {
 emit: "event",
 level: "query",
 },
 {
 emit: "stdout",
 level: "error",
 },
 {
 emit: "stdout",
 level: "info",
 },
 {
 emit: "stdout",
 level: "warn",
 },
 ],
});

prisma.$on("query", (e) => {
 console.log("Query: " + e.query);
 console.log("Params: " + e.params);
 console.log("Duration: " + e.duration + "ms");
});
```

```sql
Query: SELECT "public"."User"."id", "public"."User"."email", "public"."User"."name" FROM "public"."User" WHERE 1=1 OFFSET $1
Params: [0]
Duration: 3ms
Query: SELECT "public"."Post"."id", "public"."Post"."title", "public"."Post"."authorId" FROM "public"."Post" WHERE "public"."Post"."authorId" IN ($1,$2,$3,$4) OFFSET $5
Params: [2, 7, 18, 29]
Duration: 2ms
```

```ts highlight=4,5,22-25;normal tab="MongoDB"
const prisma = new PrismaClient({
 log: [
 {
 emit: "event",
 level: "query",
 },
 {
 emit: "stdout",
 level: "error",
 },
 {
 emit: "stdout",
 level: "info",
 },
 {
 emit: "stdout",
 level: "warn",
 },
 ],
});

prisma.$on("query", (e) => {
 console.log("Query: " + e.query);
});
```

```bash
Query: db.User.aggregate([ { $project: { _id: 1, email: 1, name: 1, }, }, ])
Query: db.Post.aggregate([ { $match: { userId: { $in: [ "622f0bbbdf635a42016ee325", ], }, }, }, { $project: { _id: 1, slug: 1, title: 1, body: 1, userId: 1, }, }, ])
```

The exact [event (`e`) type and the properties available](/orm/reference/prisma-client-reference#event-types) depends on the log level.

---
title: OpenTelemetry tracing
description: Diagnose application performance with detailed traces of each query
url: /orm/prisma-client/observability-and-logging/opentelemetry-tracing
metaTitle: OpenTelemetry tracing
metaDescription: Diagnose application performance with detailed traces of each query.
---

Tracing provides a detailed log of the activity that Prisma Client carries out, at an operation level, including the time taken to execute each query. It helps you analyze your application's performance and identify bottlenecks. Tracing is fully compliant with [OpenTelemetry](https://opentelemetry.io/), so you can use it as part of your end-to-end application tracing system.

:::info

Tracing gives you a highly detailed, operation-level insight into your Prisma ORM project.

:::

:::tip[Correlate database queries with traces]

You can add the `traceparent` header to your SQL queries as comments using the [`@prisma/sqlcommenter-trace-context`](/orm/prisma-client/observability-and-logging/sql-comments#trace-context) plugin. This enables correlation between distributed traces and database queries in your monitoring tools.

:::

## About tracing

When you enable tracing, Prisma Client outputs the following:

- One trace for each operation (e.g. findMany) that Prisma Client makes.
- In each trace, one or more [spans](https://opentelemetry.io/docs/specs/otel/trace/api/#span). Each span represents the length of time that one stage of the operation takes, such as serialization, or a database query. Spans are represented in a tree structure, where child spans indicate that execution is happening within a larger parent span.

The number and type of spans in a trace depends on the type of operation the trace covers, but an example is as follows:

![Example Prisma Client trace structure showing parent and child spans for a database operation (serialization, query engine, database query).](/img/orm/prisma-client/observability-and-logging/trace-diagram.png)

You can [send tracing output to the console](#send-tracing-output-to-the-console), or analyze it in any OpenTelemetry-compatible tracing system, such as [Jaeger](https://www.jaegertracing.io/), [Honeycomb](https://www.honeycomb.io/distributed-tracing) and [Datadog](https://www.datadoghq.com/). On this page, we give an example of how to send tracing output to Jaeger, which you can [run locally](#visualize-traces-with-jaeger).

## Trace output

For each trace, Prisma Client outputs a series of spans. The number and type of these spans depends on the Prisma Client operation. A typical Prisma trace has the following spans:

- `prisma:client:operation`: Represents the entire Prisma Client operation, from Prisma Client to the database and back. It contains details such as the model and method called by Prisma Client. Depending on the Prisma operation, it contains one or more of the following spans:
 - `prisma:client:connect`: Represents how long it takes for Prisma Client to connect to the database.
 - `prisma:client:serialize`: Represents how long it takes to validate and transform a Prisma Client operation into a query for the query engine.
 - `prisma:engine:query`: Represents how long a query takes in the query engine.
 - `prisma:engine:connection`: Represents how long it takes for Prisma Client to get a database connection.
 - `prisma:engine:db_query`: Represents the database query that was executed against the database. It includes the query in the tags, and how long the query took to run.
 - `prisma:engine:serialize`: Represents how long it takes to transform a raw response from the database into a typed result.
 - `prisma:engine:response_json_serialization`: Represents how long it takes to serialize the database query result into a JSON response to the Prisma Client.

For example, given the following Prisma Client code:

```ts
prisma.user.findMany({
 where: {
 email: email,
 },
 include: {
 posts: true,
 },
});
```

The trace is structured as follows:

- `prisma:client:operation`
 - `prisma:client:serialize`
 - `prisma:engine:query`
 - `prisma:engine:connection`
 - `prisma:engine:db_query`: details of the first SQL query or command...
 - `prisma:engine:db_query`: ...details of the next SQL query or command...
 - `prisma:engine:serialize`
 - `prisma:engine:response_json_serialization`

## Considerations and prerequisites

If your application sends a large number of spans to a [collector](https://opentelemetry.io/docs/collector/), this can have a significant performance impact. For information on how to minimize this impact, see [Reducing performance impact](#reduce-performance-impact).

To use tracing, you must do the following:

1. [Install the appropriate dependencies](#step-1-install-up-to-date-prisma-orm-dependencies).
1. [Install OpenTelemetry packages](#step-2-install-opentelemetry-packages).
1. [Register tracing in your application](#step-3-register-tracing-in-your-application).

## Get started with tracing in Prisma ORM

This section explains how to install and register tracing in your application.

### Step 1. Install Prisma ORM dependencies

Install the `prisma`, `@prisma/client`, and `@prisma/instrumentation` npm packages. You will also need to install the `@opentelemetry/api` package as it's a peer dependency.

```npm
npm install prisma@latest --save-dev
npm install @prisma/client@latest --save
npm install @prisma/instrumentation@latest --save
npm install @opentelemetry/api@latest --save
```

<details>
<summary>Tracing on previous versions of Prisma ORM</summary>

Tracing was added in version `4.2.0` of Prisma ORM as a Preview feature. For versions of Prisma ORM between `4.2.0` and `6.1.0`, you need to enable the `tracing` Preview feature in your Prisma schema file.

```prisma
generator client {
 provider = "prisma-client"
 output = "./generated"
 previewFeatures = ["tracing"]
}
```

</details>

### Step 2: Install OpenTelemetry packages

Now install the appropriate OpenTelemetry packages, as follows:

```npm
npm install @opentelemetry/semantic-conventions \
 @opentelemetry/exporter-trace-otlp-http \
 @opentelemetry/sdk-trace-base \
 @opentelemetry/sdk-trace-node \
 @opentelemetry/resources
```

### Step 3: Register tracing in your application

The following code provides two examples of configuring OpenTelemetry tracing in Prisma:

1. Using `@opentelemetry/sdk-trace-node` (existing example), which gives fine-grained control over tracing setup.
2. Using `@opentelemetry/sdk-node`, which offers a simpler configuration and aligns with OpenTelemetry's JavaScript getting started guide.

---

#### Option 1: Using `@opentelemetry/sdk-trace-node`

This setup gives you fine-grained control over instrumentation and tracing. You need to customize this configuration for your specific application. This approach is concise and easier for users who need a quick setup for sending traces to OTLP-compatible backends, such as Honeycomb, Jaeger, or Datadog.

```ts
// Imports
import { ATTR_SERVICE_NAME, ATTR_SERVICE_VERSION } from "@opentelemetry/semantic-conventions";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-http";
import { SimpleSpanProcessor } from "@opentelemetry/sdk-trace-base";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { PrismaInstrumentation, registerInstrumentations } from "@prisma/instrumentation";
import { resourceFromAttributes } from "@opentelemetry/resources";

// Configure the trace provider
const provider = new NodeTracerProvider({
 resource: resourceFromAttributes({
 [ATTR_SERVICE_NAME]: "example application", // Replace with your service name
 [ATTR_SERVICE_VERSION]: "0.0.1", // Replace with your service version
 }),
 spanProcessors: [
 // Configure how spans are processed and exported. In this case, we're sending spans
 // as we receive them to an OTLP-compatible collector (e.g., Jaeger).
 new SimpleSpanProcessor(new OTLPTraceExporter()),
 ],
});

// Register your auto-instrumentors
registerInstrumentations({
 tracerProvider: provider,
 instrumentations: [new PrismaInstrumentation()],
});

// Register the provider globally
provider.register();
```

This approach provides maximum flexibility but may involve additional configuration steps.

#### Option 2: Using `@opentelemetry/sdk-node`

For many users, especially beginners, the `NodeSDK` class simplifies OpenTelemetry setup by bundling common defaults into a single, unified configuration.

```ts
// Imports
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-proto";
import { NodeSDK } from "@opentelemetry/sdk-node";
import { PrismaInstrumentation } from "@prisma/instrumentation";

// Configure the OTLP trace exporter
const traceExporter = new OTLPTraceExporter({
 url: "https://api.honeycomb.io/v1/traces", // Replace with your collector's endpoint
 headers: {
 "x-honeycomb-team": "HONEYCOMB_API_KEY", // Replace with your Honeycomb API key or collector auth header
 },
});

// Initialize the NodeSDK
const sdk = new NodeSDK({
 serviceName: "my-service-name", // Replace with your service name
 traceExporter,
 instrumentations: [new PrismaInstrumentation()],
});

// Start the SDK
sdk.start();

// Handle graceful shutdown
process.on("SIGTERM", async () => {
 try {
 await sdk.shutdown();
 console.log("Tracing shut down successfully");
 } catch (err) {
 console.error("Error shutting down tracing", err);
 } finally {
 process.exit(0);
 }
});
```

Choose the `NodeSDK` approach if:

- You are starting with OpenTelemetry and want a simplified setup.
- You need to quickly integrate tracing with minimal boilerplate.
- You are using an OTLP-compatible tracing backend like Honeycomb, Jaeger, or Datadog.

Choose the `NodeTracerProvider` approach if:

- You need detailed control over how spans are created, processed, and exported.
- You are using custom span processors or exporters.
- Your application requires specific instrumentation or sampling strategies.

OpenTelemetry is highly configurable. You can customize the resource attributes, what components gets instrumented, how spans are processed, and where spans are sent.

You can find a complete example that includes metrics in [this sample application](https://github.com/garrensmith/prisma-metrics-sample).

## Tracing how-tos

### Visualize traces with Jaeger

[Jaeger](https://www.jaegertracing.io/) is a free and open source OpenTelemetry collector and dashboard that you can use to visualize your traces.

The following screenshot shows an example trace visualization:

![Jaeger UI](/img/orm/prisma-client/observability-and-logging/jaeger.png)

To run Jaeger locally, use the following [Docker](https://www.docker.com/) command:

```console
docker run --rm --name jaeger -d -e COLLECTOR_OTLP_ENABLED=true -p 16686:16686 -p 4318:4318 jaegertracing/all-in-one:latest
```

You'll now find the tracing dashboard available at `http://localhost:16686/`. When you use your application with tracing enabled, you'll start to see traces in this dashboard.

### Send tracing output to the console

The following example sends output tracing to the console with `ConsoleSpanExporter` from `@opentelemetry/sdk-trace-base`.

```ts
// Imports
import { ATTR_SERVICE_NAME, ATTR_SERVICE_VERSION } from "@opentelemetry/semantic-conventions";
import { ConsoleSpanExporter, SimpleSpanProcessor } from "@opentelemetry/sdk-trace-base";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { AsyncHooksContextManager } from "@opentelemetry/context-async-hooks";
import * as api from "@opentelemetry/api";
import { PrismaInstrumentation, registerInstrumentations } from "@prisma/instrumentation";
import { resourceFromAttributes } from "@opentelemetry/resources";

// Export the tracing
export function otelSetup() {
 const contextManager = new AsyncHooksContextManager().enable();

 api.context.setGlobalContextManager(contextManager);

 //Configure the console exporter
 const consoleExporter = new ConsoleSpanExporter();

 // Configure the trace provider
 const provider = new NodeTracerProvider({
 resource: resourceFromAttributes({
 [ATTR_SERVICE_NAME]: "example application", // Replace with your service name
 [ATTR_SERVICE_VERSION]: "0.0.1", // Replace with your service version
 }),
 spanProcessors: [
 // Configure how spans are processed and exported. In this case, we're sending spans
 // as we receive them to the console
 new SimpleSpanProcessor(consoleExporter),
 ],
 });

 // Register your auto-instrumentors
 registerInstrumentations({
 tracerProvider: provider,
 instrumentations: [new PrismaInstrumentation()],
 });

 // Register the provider
 provider.register();
}
```

### Trace interactive transactions

When you perform an interactive transaction, you'll see the following span in addition to the [standard spans](#trace-output):

- `prisma:client:transaction`: A [root span](https://opentelemetry.io/docs/concepts/observability-primer/#distributed-traces) that wraps the `prisma` span.

As an example, take the following Prisma schema:

```prisma title="schema.prisma" showLineNumbers
generator client {
 provider = "prisma-client"
 output = "./generated"
}

datasource db {
 provider = "postgresql"
}

model User {
 id Int @id @default(autoincrement())
 email String @unique
}

model Audit {
 id Int @id
 table String
 action String
}
```

Given the following interactive transaction:

```ts
await prisma.$transaction(async (tx) => {
 const user = await tx.user.create({
 data: {
 email: email,
 },
 });

 await tx.audit.create({
 data: {
 table: "user",
 action: "create",
 id: user.id,
 },
 });

 return user;
});
```

The trace is structured as follows:

- `prisma:client:transaction`
- `prisma:client:connect`
- `prisma:engine:itx_runner`
 - `prisma:engine:connection`
 - `prisma:engine:db_query`
 - `prisma:engine:itx_query_builder`
 - `prisma:engine:db_query`
 - `prisma:engine:db_query`
 - `prisma:engine:serialize`
 - `prisma:engine:itx_query_builder`
 - `prisma:engine:db_query`
 - `prisma:engine:db_query`
 - `prisma:engine:serialize`
- `prisma:client:operation`
 - `prisma:client:serialize`
- `prisma:client:operation`
 - `prisma:client:serialize`

### Add more instrumentation

A nice benefit of OpenTelemetry is the ability to add more instrumentation with only minimal changes to your application code.

For example, to add HTTP and [ExpressJS](https://expressjs.com/) tracing, add the following instrumentations to your OpenTelemetry configuration. These instrumentations add spans for the full request-response lifecycle. These spans show you how long your HTTP requests take.

```js
// Imports
import { ExpressInstrumentation } from "@opentelemetry/instrumentation-express";
import { HttpInstrumentation } from "@opentelemetry/instrumentation-http";

// Register your auto-instrumentors
registerInstrumentations({
 tracerProvider: provider,
 instrumentations: [
 new HttpInstrumentation(),
 new ExpressInstrumentation(),
 new PrismaInstrumentation(),
 ],
});
```

For a full list of available instrumentation, take a look at the [OpenTelemetry Registry](https://opentelemetry.io/ecosystem/registry/?language=js&component=instrumentation).

### Customize resource attributes

You can adjust how your application's traces are grouped by changing the resource attributes to be more specific to your application:

```js
const provider = new NodeTracerProvider({
 resource: new Resource({
 [ATTR_SERVICE_NAME]: "weblog",
 [ATTR_SERVICE_VERSION]: "1.0.0",
 }),
});
```

There is an ongoing effort to standardize common resource attributes. Whenever possible, it's a good idea to follow the [standard attribute names](https://github.com/open-telemetry/semantic-conventions/blob/main/docs/general/trace.md).

### Reduce performance impact

If your application sends a large number of spans to a collector, this can have a significant performance impact. You can use the following approaches to reduce this impact:

- [Use the BatchSpanProcessor](#send-traces-in-batches-using-the-batchspanprocessor)
- [Send fewer spans to the collector](#send-fewer-spans-to-the-collector-with-sampling)

#### Send traces in batches using the `BatchSpanProcessor`

In a production environment, you can use OpenTelemetry's `BatchSpanProcessor` to send the spans to a collector in batches rather than one at a time. However, during development and testing, you might not want to send spans in batches. In this situation, you might prefer to use the `SimpleSpanProcessor`.

You can configure your tracing configuration to use the appropriate span processor, depending on the environment, as follows:

```ts
import { SimpleSpanProcessor, BatchSpanProcessor } from "@opentelemetry/sdk-trace-base";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";

const spanProcessors = [];
if (process.env.NODE_ENV === "production") {
 spanProcessors.push(new BatchSpanProcessor(otlpTraceExporter));
} else {
 spanProcessors.push(new SimpleSpanProcessor(otlpTraceExporter));
}

const provider = new NodeTracerProvider({
 spanProcessors,
 // ...other configurations
});
```

#### Send fewer spans to the collector with sampling

Another way to reduce the performance impact is to [use probability sampling](https://opentelemetry.io/docs/specs/otel/trace/tracestate-probability-sampling/) to send fewer spans to the collector. This reduces the collection cost of tracing but still gives a good representation of what is happening in your application.

An example implementation looks like this:

```ts highlight=3,7;add
import { ATTR_SERVICE_NAME, ATTR_SERVICE_VERSION } from "@opentelemetry/semantic-conventions";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { TraceIdRatioBasedSampler } from "@opentelemetry/core";
import { resourceFromAttributes } from "@opentelemetry/resources";

const provider = new NodeTracerProvider({
 sampler: new TraceIdRatioBasedSampler(0.1), // [!code ++]
 resource: resourceFromAttributes({
 // we can define some metadata about the trace resource
 [ATTR_SERVICE_NAME]: "test-tracing-service",
 [ATTR_SERVICE_VERSION]: "1.0.0",
 }),
});
```

## Troubleshoot tracing

### My traces aren't showing up

The order in which you set up tracing matters. In your application, ensure that you register tracing and instrumentation before you import any instrumented dependencies. For example:

```ts
import { registerTracing } from "./tracing";

registerTracing({
 name: "tracing-example",
 version: "0.0.1",
});

// You must import any dependencies after you register tracing.
import { PrismaClient } from "../prisma/generated/client";
import async from "express-async-handler";
import express from "express";
```

---
title: SQL comments
description: 'Add metadata to your SQL queries as comments for improved observability, debugging, and tracing'
url: /orm/prisma-client/observability-and-logging/sql-comments
metaTitle: SQL comments
metaDescription: 'Add metadata to your SQL queries as comments for improved observability, debugging, and tracing.'
---

SQL comments allow you to append metadata to your database queries, making it easier to correlate queries with application context. Prisma ORM supports the [sqlcommenter format](https://google.github.io/sqlcommenter/) developed by Google, which is widely supported by database monitoring tools.

SQL comments are useful for:

- **Observability**: Correlate database queries with application traces using `traceparent`
- **Query insights**: Tag queries with metadata for analysis in database monitoring tools
- **Debugging**: Add custom context to queries for easier troubleshooting

## Installation

Install one or more first-party plugins depending on your use case:

```npm
npm install @prisma/sqlcommenter-query-tags
npm install @prisma/sqlcommenter-trace-context
```

Install the core SQL commenter types package to create your own plugin:

```npm
npm install @prisma/sqlcommenter
```

## Basic usage

Pass an array of SQL commenter plugins to the `comments` option when creating a `PrismaClient` instance:

```ts
import { PrismaClient } from "../prisma/generated/client";
import { PrismaPg } from "@prisma/adapter-pg";
import { queryTags } from "@prisma/sqlcommenter-query-tags";
import { traceContext } from "@prisma/sqlcommenter-trace-context";

const adapter = new PrismaPg({ connectionString: process.env.DATABASE_URL });

const prisma = new PrismaClient({
 adapter,
 comments: [queryTags(), traceContext()],
});
```

With this configuration, your SQL queries will include metadata as comments:

```sql
SELECT "id", "name" FROM "User" /*application='my-app',traceparent='00-abc123...-01'*/
```

## First-party plugins

Prisma provides two official SQL commenter plugins:

### Query tags

The `@prisma/sqlcommenter-query-tags` package allows you to add arbitrary tags to queries within an async context using `AsyncLocalStorage`.

```ts
import { queryTags, withQueryTags } from "@prisma/sqlcommenter-query-tags";
import { PrismaClient } from "../prisma/generated/client";

const prisma = new PrismaClient({
 adapter,
 comments: [queryTags()],
});

// Wrap your queries to add tags
const users = await withQueryTags({ route: "/api/users", requestId: "abc-123" }, () =>
 prisma.user.findMany(),
);
```

The resulting SQL includes the tags as comments:

```sql
SELECT ... FROM "User" /*requestId='abc-123',route='/api/users'*/
```

#### Multiple queries in one scope

All queries within the callback share the same tags:

```ts
const result = await withQueryTags({ traceId: "trace-456" }, async () => {
 const users = await prisma.user.findMany();
 const posts = await prisma.post.findMany();
 return { users, posts };
});
```

#### Nested scopes with tag replacement

By default, nested `withQueryTags` calls replace the outer tags entirely:

```ts
await withQueryTags({ requestId: "req-123" }, async () => {
 // Queries here have: requestId='req-123'

 await withQueryTags({ userId: "user-456" }, async () => {
 // Queries here only have: userId='user-456'
 // requestId is NOT included
 await prisma.user.findMany();
 });
});
```

#### Nested scopes with tag merging

Use `withMergedQueryTags` to merge tags with the outer scope:

```ts
import { withQueryTags, withMergedQueryTags } from "@prisma/sqlcommenter-query-tags";

await withQueryTags({ requestId: "req-123", source: "api" }, async () => {
 await withMergedQueryTags({ userId: "user-456", source: "handler" }, async () => {
 // Queries here have: requestId='req-123', userId='user-456', source='handler'
 await prisma.user.findMany();
 });
});
```

You can also remove tags in nested scopes by setting them to `undefined`:

```ts
await withQueryTags({ requestId: "req-123", debug: "true" }, async () => {
 await withMergedQueryTags({ userId: "user-456", debug: undefined }, async () => {
 // Queries here have: requestId='req-123', userId='user-456'
 // debug is removed
 await prisma.user.findMany();
 });
});
```

### Trace context

The `@prisma/sqlcommenter-trace-context` package adds W3C Trace Context (`traceparent`) headers to your queries, enabling correlation between distributed traces and database queries.

```ts
import { traceContext } from "@prisma/sqlcommenter-trace-context";
import { PrismaClient } from "../prisma/generated/client";

const prisma = new PrismaClient({
 adapter,
 comments: [traceContext()],
});
```

When tracing is enabled and the current span is sampled, queries include the `traceparent`:

```sql
SELECT * FROM "User" /*traceparent='00-0af7651916cd43dd8448eb211c80319c-b9c7c989f97918e1-01'*/
```

:::info

The trace context plugin requires [`@prisma/instrumentation`](/orm/prisma-client/observability-and-logging/opentelemetry-tracing) to be configured. The `traceparent` is only added when tracing is active and the span is sampled.

:::

The `traceparent` header follows the [W3C Trace Context](https://www.w3.org/TR/trace-context/) specification:

```
{version}-{trace-id}-{parent-id}-{trace-flags}
```

Where:

- `version`: Always `00` for the current spec
- `trace-id`: 32 hexadecimal characters representing the trace ID
- `parent-id`: 16 hexadecimal characters representing the parent span ID
- `trace-flags`: 2 hexadecimal characters; `01` indicates sampled

## Creating custom plugins

You can create your own SQL commenter plugins to add custom metadata to queries.

### Plugin structure

A SQL commenter plugin is a function that receives query context and returns key-value pairs:

```ts
import type { SqlCommenterPlugin, SqlCommenterContext } from "@prisma/sqlcommenter";

const myPlugin: SqlCommenterPlugin = (context: SqlCommenterContext) => {
 return {
 application: "my-app",
 version: "1.0.0",
 };
};
```

### Using custom plugins

Pass your custom plugins to the `comments` option:

```ts
const prisma = new PrismaClient({
 adapter,
 comments: [myPlugin],
});
```

### Conditional keys

Return `undefined` for keys you want to exclude from the comment. Keys with `undefined` values are automatically filtered out:

```ts
const conditionalPlugin: SqlCommenterPlugin = (context) => ({
 model: context.query.modelName, // undefined for raw queries, automatically omitted
 action: context.query.action,
});
```

### Query context

Plugins receive a `SqlCommenterContext` object containing information about the query:

```ts
interface SqlCommenterContext {
 query: SqlCommenterQueryInfo;
 sql?: string;
}
```

The `query` property provides information about the Prisma operation:

| Property | Type | Description |
| ----------- | ------------------------------------------------------ | ---------------------------------------------------------------------- |
| `type` | `'single'` \| `'compacted'` | Whether this is a single query or a batched query |
| `modelName` | `string` \| `undefined` | The model being queried (e.g., `"User"`). Undefined for raw queries. |
| `action` | `string` | The Prisma operation (e.g., `"findMany"`, `"createOne"`, `"queryRaw"`) |
| `query` | `unknown` (single) or `queries: unknown[]` (compacted) | The full query object(s). Structure is not part of the public API. |

The `sql` property is the raw SQL query generated from this Prisma query. It is always available when `PrismaClient` connects to the database and renders SQL queries directly. When using Prisma Accelerate, SQL rendering happens on Accelerate side and the raw SQL strings are not available when SQL commenter plugins are executed on the `PrismaClient` side.

#### Single vs. compacted queries

- **Single queries** (`type: 'single'`): A single Prisma query is being executed
- **Compacted queries** (`type: 'compacted'`): Multiple queries have been batched into a single SQL statement (e.g., automatic `findUnique` batching)

### Example: Application metadata

```ts
import type { SqlCommenterPlugin } from "@prisma/sqlcommenter";

const applicationTags: SqlCommenterPlugin = (context) => ({
 application: "my-service",
 environment: process.env.NODE_ENV ?? "development",
 operation: context.query.action,
 model: context.query.modelName,
});
```

### Example: Async context propagation

Use `AsyncLocalStorage` to propagate context through your application:

```ts
import { AsyncLocalStorage } from "node:async_hooks";
import type { SqlCommenterPlugin } from "@prisma/sqlcommenter";

interface RequestContext {
 route: string;
 userId?: string;
}

const requestStorage = new AsyncLocalStorage<RequestContext>();

const requestContextPlugin: SqlCommenterPlugin = () => {
 const context = requestStorage.getStore();
 return {
 route: context?.route,
 userId: context?.userId,
 };
};

// Usage in a request handler
requestStorage.run({ route: "/api/users", userId: "user-123" }, async () => {
 await prisma.user.findMany();
});
```

### Combining multiple plugins

Plugins are called in array order, and their outputs are merged. Later plugins can override keys from earlier plugins:

```ts
import type { SqlCommenterPlugin } from "@prisma/sqlcommenter";
import { queryTags } from "@prisma/sqlcommenter-query-tags";
import { traceContext } from "@prisma/sqlcommenter-trace-context";

const appPlugin: SqlCommenterPlugin = () => ({
 application: "my-app",
 version: "1.0.0",
});

const prisma = new PrismaClient({
 adapter,
 comments: [appPlugin, queryTags(), traceContext()],
});
```

## Framework integration

### Hono

Hono's middleware properly awaits downstream handlers:

```ts
import { createMiddleware } from "hono/factory";
import { withQueryTags } from "@prisma/sqlcommenter-query-tags";

app.use(
 createMiddleware(async (c, next) => {
 await withQueryTags(
 {
 route: c.req.path,
 method: c.req.method,
 requestId: c.req.header("x-request-id") ?? crypto.randomUUID(),
 },
 () => next(),
 );
 }),
);
```

### Koa

Koa's middleware properly awaits downstream handlers:

```ts
import { withQueryTags } from "@prisma/sqlcommenter-query-tags";

app.use(async (ctx, next) => {
 await withQueryTags(
 {
 route: ctx.path,
 method: ctx.method,
 requestId: ctx.get("x-request-id") || crypto.randomUUID(),
 },
 () => next(),
 );
});
```

### Fastify

Wrap individual route handlers:

```ts
import { withQueryTags } from "@prisma/sqlcommenter-query-tags";

fastify.get("/users", (request, reply) => {
 return withQueryTags(
 {
 route: "/users",
 method: "GET",
 requestId: request.id,
 },
 () => prisma.user.findMany(),
 );
});
```

### Express

Express middleware uses callbacks, so wrap route handlers directly:

```ts
import { withQueryTags } from "@prisma/sqlcommenter-query-tags";

app.get("/users", (req, res, next) => {
 withQueryTags(
 {
 route: req.path,
 method: req.method,
 requestId: req.header("x-request-id") ?? crypto.randomUUID(),
 },
 () => prisma.user.findMany(),
 )
 .then((users) => res.json(users))
 .catch(next);
});
```

### NestJS

Use an interceptor to wrap handler execution:

```ts
import { Injectable, NestInterceptor, ExecutionContext, CallHandler } from "@nestjs/common";
import { Observable, from, lastValueFrom } from "rxjs";
import { withQueryTags } from "@prisma/sqlcommenter-query-tags";

@Injectable()
export class QueryTagsInterceptor implements NestInterceptor {
 intercept(context: ExecutionContext, next: CallHandler): Observable<unknown> {
 const request = context.switchToHttp().getRequest<Request>();
 return from(
 withQueryTags(
 {
 route: request.url,
 method: request.method,
 requestId: request.headers.get("x-request-id") ?? crypto.randomUUID(),
 },
 () => lastValueFrom(next.handle()),
 ),
 );
 }
}

// Apply globally in main.ts
app.useGlobalInterceptors(new QueryTagsInterceptor());
```

## Output format

Plugin outputs are merged, sorted alphabetically by key, URL-encoded, and formatted according to the [sqlcommenter specification](https://google.github.io/sqlcommenter/spec/):

```sql
SELECT "id", "name" FROM "User" /*application='my-app',environment='production',model='User'*/
```

Key behaviors:

- Plugins are called synchronously in array order
- Later plugins override earlier ones if they return the same key
- Keys with `undefined` values are filtered out (they do not remove keys set by earlier plugins)
- Keys and values are URL-encoded per the sqlcommenter spec
- Single quotes in values are escaped as `\'`
- Comments are appended to the end of SQL queries

## API reference

### `SqlCommenterTags`

```ts
type SqlCommenterTags = { readonly [key: string]: string | undefined };
```

Key-value pairs to add as SQL comments. Keys with `undefined` values are automatically filtered out.

### `SqlCommenterPlugin`

```ts
interface SqlCommenterPlugin {
 (context: SqlCommenterContext): SqlCommenterTags;
}
```

A function that receives query context and returns key-value pairs. Return an empty object to add no comments for a particular query.

### `SqlCommenterContext`

```ts
interface SqlCommenterContext {
 query: SqlCommenterQueryInfo;
 sql?: string;
}
```

Context provided to plugins containing information about the query.

- **`query`**: Information about the Prisma query being executed. See [`SqlCommenterQueryInfo`](#sqlcommenterqueryinfo).
- **`sql`**: The SQL query being executed. It is only available when using driver adapters but not when using Accelerate.

### `SqlCommenterQueryInfo`

```ts
type SqlCommenterQueryInfo =
 | ({ type: "single" } & SqlCommenterSingleQueryInfo)
 | ({ type: "compacted" } & SqlCommenterCompactedQueryInfo);
```

Information about the query or queries being executed.

### `SqlCommenterSingleQueryInfo`

```ts
interface SqlCommenterSingleQueryInfo {
 modelName?: string;
 action: string;
 query: unknown;
}
```

Information about a single Prisma query.

### `SqlCommenterCompactedQueryInfo`

```ts
interface SqlCommenterCompactedQueryInfo {
 modelName?: string;
 action: string;
 queries: unknown[];
}
```

Information about a compacted batch query.

---
title: 'Aggregation, grouping, and summarizing'
description: 'Use Prisma Client to aggregate, group by, count, and select distinct.'
url: /orm/prisma-client/queries/aggregation-grouping-summarizing
metaTitle: 'Aggregation, grouping, and summarizing (Concepts)'
metaDescription: 'Use Prisma Client to aggregate, group by, count, and select distinct.'
---

Prisma Client allows you to count records, aggregate number fields, and select distinct field values.

## Aggregate

Prisma Client allows you to [`aggregate`](/orm/reference/prisma-client-reference#aggregate) on the **number** fields (such as `Int` and `Float`) of a model. The following query returns the average age of all users:

```ts
const aggregations = await prisma.user.aggregate({
 _avg: { age: true },
});

console.log('Average age:' + aggregations._avg.age);
```

You can combine aggregation with filtering and ordering. For example, the following query returns the average age of users:

- Ordered by `age` ascending
- Where `email` contains `prisma.io`
- Limited to the 10 users

```ts
const aggregations = await prisma.user.aggregate({
 _avg: { age: true },
 where: {
 email: {
 contains: 'prisma.io',
 },
 },
 orderBy: { age: 'asc' },
 take: 10,
});

console.log('Average age:' + aggregations._avg.age);
```

### Aggregate values are nullable

Aggregations on **nullable fields** can return a `number` or `null`. This excludes `count`, which always returns 0 if no records are found.

Consider the following query, where `age` is nullable in the schema:

```ts
const aggregations = await prisma.user.aggregate({
 _avg: { age: true },
 _count: { age: true },
});
```

```json
{
 "_avg": { "age": null },
 "_count": { "age": 9 }
}
```

The query returns `{ _avg: { age: null } }` in either of the following scenarios:

- There are no users
- The value of every user's `age` field is `null`

This allows you to differentiate between the true aggregate value (which could be zero) and no data.

## Group by

Prisma Client's [`groupBy()`](/orm/reference/prisma-client-reference#groupby) allows you to **group records** by one or more field values - such as `country`, or `country` and `city` and **perform aggregations** on each group, such as finding the average age of people living in a particular city.

The following example groups all users by the `country` field and returns the total number of profile views for each country:

```ts
const groupUsers = await prisma.user.groupBy({
 by: ['country'],
 _sum: { profileViews: true },
});
```

```json
[
 { country: 'Germany', _sum: { profileViews: 126 } },
 { country: 'Sweden', _sum: { profileViews: 0 } },
];
```

If you have a single element in the `by` option, you can use the following shorthand syntax to express your query:

```ts
const groupUsers = await prisma.user.groupBy({
 by: 'country',
});
```

### `groupBy()` and filtering

`groupBy()` supports two levels of filtering: `where` and `having`.

#### Filter records with `where`

Use `where` to filter all records **before grouping**. The following example groups users by country and sums profile views, but only includes users where the email address contains `prisma.io`:

```ts
const groupUsers = await prisma.user.groupBy({
 by: ['country'],
 where: {
 // [!code highlight]
 email: {
 // [!code highlight]
 contains: 'prisma.io', // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 _sum: {
 profileViews: true,
 },
});
```

#### Filter groups with `having`

Use `having` to filter **entire groups** by an aggregate value such as the sum or average of a field, not individual records - for example, only return groups where the _average_ `profileViews` is greater than 100:

```ts
const groupUsers = await prisma.user.groupBy({
 by: ['country'],
 where: {
 email: {
 contains: 'prisma.io',
 },
 },
 _sum: { profileViews: true, },
 having: {
 // [!code highlight]
 profileViews: {
 // [!code highlight]
 _avg: {
 // [!code highlight]
 gt: 100, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
});
```

##### Use case for `having`

The primary use case for `having` is to filter on aggregations. We recommend that you use `where` to reduce the size of your data set as far as possible _before_ grouping, because doing so ✔ reduces the number of records the database has to return and ✔ makes use of indices.

For example, the following query groups all users that are _not_ from Sweden or Ghana:

```ts
const fd = await prisma.user.groupBy({
 by: ['country'],
 where: {
 country: {
 // [!code highlight]
 notIn: ['Sweden', 'Ghana'], // [!code highlight]
 }, // [!code highlight]
 },
 _sum: {
 profileViews: true,
 },
 having: {
 profileViews: {
 _min: {
 gte: 10,
 },
 },
 },
});
```

The following query technically achieves the same result, but excludes users from Ghana _after_ grouping. This does not confer any benefit and is not recommended practice.

```ts
const groupUsers = await prisma.user.groupBy({
 by: ['country'],
 where: {
 country: {
 // [!code highlight]
 not: 'Sweden', // [!code highlight]
 }, // [!code highlight]
 },
 _sum: {
 profileViews: true,
 },
 having: {
 country: {
 // [!code highlight]
 not: 'Ghana', // [!code highlight]
 }, // [!code highlight]
 profileViews: {
 _min: {
 gte: 10,
 },
 },
 },
});
```

> **Note**: Within `having`, you can only filter on aggregate values _or_ fields available in `by`.

### `groupBy()` and ordering

The following constraints apply when you combine `groupBy()` and `orderBy`:

- You can `orderBy` fields that are present in `by`
- You can `orderBy` aggregate (Preview in 2.21.0 and later)
- If you use `skip` and/or `take` with `groupBy()`, you must also include `orderBy` in the query

#### Order by aggregate group

You can **order by aggregate group**. The following example sorts each `city` group by the number of users in that group (largest group first):

```ts
const groupBy = await prisma.user.groupBy({
 by: ['city'],
 _count: {
 city: true,
 },
 orderBy: {
 _count: {
 city: 'desc',
 },
 },
});
```

```json
[
 { city: 'Berlin', count: { city: 3 } },
 { city: 'Paris', count: { city: 2 } },
 { city: 'Amsterdam', count: { city: 1 } },
];
```

#### Order by field

The following query orders groups by country, skips the first two groups, and returns the 3rd and 4th group:

```ts
const groupBy = await prisma.user.groupBy({
 by: ['country'],
 _sum: {
 profileViews: true,
 },
 orderBy: {
 country: 'desc',
 },
 skip: 2,
 take: 2,
});
```

### `groupBy()` FAQ

#### Can I use `select` with `groupBy()`?

You cannot use `select` with `groupBy()`. However, all fields included in `by` are automatically returned.

#### What is the difference between using `where` and `having` with `groupBy()`?

`where` filters all records before grouping, and `having` filters entire groups and supports filtering on an aggregate field value, such as the average or sum of a particular field in that group.

#### What is the difference between `groupBy()` and `distinct`?

Both `distinct` and `groupBy()` group records by one or more unique field values. `groupBy()` allows you to aggregate data within each group - for example, return the average number of views on posts from Denmark - whereas distinct does not.

## Count

### Count records

Use [`count()`](/orm/reference/prisma-client-reference#count) to count the number of records or non-`null` field values. The following example query counts all users:

```ts
const userCount = await prisma.user.count();
```

### Count relations

To return a count of relations (for example, a user's post count), use the `_count` parameter with a nested `select` as shown:

```ts
const usersWithCount = await prisma.user.findMany({
 include: {
 _count: {
 select: { posts: true },
 },
 },
});
```

```json
{ id: 1, _count: { posts: 3 } },
{ id: 2, _count: { posts: 2 } },
{ id: 3, _count: { posts: 2 } },
{ id: 4, _count: { posts: 0 } },
{ id: 5, _count: { posts: 0 } }
```

The `_count` parameter:

- Can be used inside a top-level `include` _or_ `select`
- Can be used with any query that returns records (including `delete`, `update`, and `findFirst`)
- Can return [multiple relation counts](#return-multiple-relation-counts)
- Can [filter relation counts](#filter-the-relation-count) (from version 4.3.0)

#### Return a relations count with `include`

The following query includes each user's post count in the results:

```ts
const usersWithCount = await prisma.user.findMany({
 include: {
 _count: {
 select: { posts: true },
 },
 },
});
```

```json
{ id: 1, _count: { posts: 3 } },
{ id: 2, _count: { posts: 2 } },
{ id: 3, _count: { posts: 2 } },
{ id: 4, _count: { posts: 0 } },
{ id: 5, _count: { posts: 0 } }
```

#### Return a relations count with `select`

The following query uses `select` to return each user's post count _and no other fields_:

```ts
const usersWithCount = await prisma.user.findMany({
 select: {
 _count: {
 select: { posts: true },
 },
 },
});
```

```json
{
 _count: {
 posts: 3;
 }
}
```

#### Return multiple relation counts

The following query returns a count of each user's `posts` and `recipes` and no other fields:

```ts
const usersWithCount = await prisma.user.findMany({
 select: {
 _count: {
 select: {
 posts: true,
 recipes: true,
 },
 },
 },
});
```

```json
{
 "_count": {
 "posts": 3,
 "recipes": 9
 }
}
```

#### Filter the relation count

Use `where` to filter the fields returned by the `_count` output type. You can do this on [scalar fields](/orm/prisma-schema/data-model/models#scalar-fields) and [relation fields](/orm/prisma-schema/data-model/models#relation-fields).

For example, the following query returns all user posts with the title "Hello!":

```ts
// Count all user posts with the title "Hello!"
await prisma.user.findMany({
 select: {
 _count: {
 select: {
 posts: { where: { title: 'Hello!' } },
 },
 },
 },
});
```

The following query finds all user posts with comments from an author named "Alice":

```ts
// Count all user posts that have comments
// whose author is named "Alice"
await prisma.user.findMany({
 select: {
 _count: {
 select: {
 posts: {
 where: { comments: { some: { author: { is: { name: 'Alice' } } } } },
 },
 },
 },
 },
});
```

### Count non-`null` field values

In [2.15.0](https://github.com/prisma/prisma/releases/2.15.0) and later, you can count all records as well as all instances of non-`null` field values. The following query returns a count of:

- All `User` records (`_all`)
- All non-`null` `name` values (not distinct values, just values that are not `null`)

```ts
const userCount = await prisma.user.count({
 select: {
 _all: true, // Count all records
 name: true, // Count all non-null field values
 },
});
```

```json
{ "_all": 30, "name": 10 }
```

### Filtered count

`count` supports filtering. The following example query counts all users with more than 100 profile views:

```ts
const userCount = await prisma.user.count({
 where: {
 profileViews: {
 gte: 100,
 },
 },
});
```

The following example query counts a particular user's posts:

```ts
const postCount = await prisma.post.count({
 where: {
 authorId: 29,
 },
});
```

## Select distinct

Prisma Client allows you to filter duplicate rows from a Prisma Query response to a [`findMany`](/orm/reference/prisma-client-reference#findmany) query using [`distinct`](/orm/reference/prisma-client-reference#distinct) . `distinct` is often used in combination with [`select`](/orm/reference/prisma-client-reference#select) to identify certain unique combinations of values in the rows of your table.

The following example returns all fields for all `User` records with distinct `name` field values:

```ts
const result = await prisma.user.findMany({
 where: {},
 distinct: ['name'],
});
```

The following example returns distinct `role` field values (for example, `ADMIN` and `USER`):

```ts
const distinctRoles = await prisma.user.findMany({
 distinct: ['role'],
 select: {
 role: true,
 },
});
```

```json
[
 { role: 'USER', },
 { role: 'ADMIN', },
];
```

### `distinct` under the hood

Prisma Client's `distinct` option does not use SQL `SELECT DISTINCT`. Instead, `distinct` uses:

- A `SELECT` query
- In-memory post-processing to select distinct

It was designed in this way in order to **support `select` and `include`** as part of `distinct` queries.

The following example selects distinct on `gameId` and `playerId`, ordered by `score`, in order to return **each player's highest score per game**. The query uses `include` and `select` to include additional data:

- Select `score` (field on `Play`)
- Select related player name (relation between `Play` and `User`)
- Select related game name (relation between `Play` and `Game`)

<details>

<summary>Expand for sample schema</summary>

```prisma title="schema.prisma"
model User {
 id Int @id @default(autoincrement())
 name String?
 play Play[]
}

model Game {
 id Int @id @default(autoincrement())
 name String?
 play Play[]
}

model Play {
 id Int @id @default(autoincrement())
 score Int? @default(0)
 playerId Int?
 player User? @relation(fields: [playerId], references: [id])
 gameId Int?
 game Game? @relation(fields: [gameId], references: [id])
}
```

</details>

```ts
const distinctScores = await prisma.play.findMany({
 distinct: ['playerId', 'gameId'],
 orderBy: {
 score: 'desc',
 },
 select: {
 score: true,
 game: {
 select: {
 name: true,
 },
 },
 player: {
 select: {
 name: true,
 },
 },
 },
});
```

```json
[
 {
 "score": 900,
 "game": { "name": "Pacman" },
 "player": { "name": "Bert Bobberton" }
 },
 {
 "score": 400,
 "game": { "name": "Pacman" },
 "player": { "name": "Nellie Bobberton" }
 }
]
```

Without `select` and `distinct`, the query would return:

```json
[
 {
 "gameId": 2,
 "playerId": 5
 },
 {
 "gameId": 2,
 "playerId": 10
 }
]
```

---
title: CRUD
description: "Learn how to perform create, read, update, and delete operations"
url: /orm/prisma-client/queries/crud
metaTitle: CRUD (Reference)
metaDescription: How to perform CRUD with Prisma Client.
---

This page describes how to perform CRUD operations with Prisma Client:

- [Create](#create) - Insert records
- [Read](#read) - Query records
- [Update](#update) - Modify records
- [Delete](#delete) - Remove records

See the [Prisma Client API reference](/orm/reference/prisma-client-reference) for detailed method documentation.

## Create

### Create a single record

```ts
const user = await prisma.user.create({
 data: {
 email: "elsa@prisma.io",
 name: "Elsa Prisma",
 },
});
```

## Where to go next

- [Prisma Client API reference](/orm/reference/prisma-client-reference) if you need the full method and option surface
- [Relation queries](/orm/prisma-client/queries/relation-queries) if your CRUD workflow needs nested writes or relational reads
- [Prisma Migrate getting started](/orm/prisma-migrate/getting-started) if you're still evolving the database schema behind these queries
- [Quickstart with Prisma Postgres](/prisma-orm/quickstart/prisma-postgres) if you want a managed database to try these examples against quickly

The `id` is auto-generated. Your schema determines which fields are mandatory.

### Create multiple records

```ts
const createMany = await prisma.user.createMany({
 data: [
 { name: "Bob", email: "bob@prisma.io" },
 { name: "Yewande", email: "yewande@prisma.io" },
 ],
 skipDuplicates: true, // Skip records with duplicate unique fields
});
// Returns: { count: 2 }
```

:::note
`skipDuplicates` is not supported on MongoDB, SQLServer, or SQLite.
:::

### Create and return multiple records

Supported by PostgreSQL, CockroachDB, and SQLite.

```ts
const users = await prisma.user.createManyAndReturn({
 data: [
 { name: "Alice", email: "alice@prisma.io" },
 { name: "Bob", email: "bob@prisma.io" },
 ],
});
```

See [Nested writes](/orm/prisma-client/queries/relation-queries#nested-writes) for creating records with relations.

## Read

### Get record by ID or unique field

```ts
// By unique field
const user = await prisma.user.findUnique({
 where: { email: "elsa@prisma.io" },
});

// By ID
const user = await prisma.user.findUnique({
 where: { id: 99 },
});
```

### Get all records

```ts
const users = await prisma.user.findMany();
```

### Get first matching record

```ts
const user = await prisma.user.findFirst({
 where: { posts: { some: { likes: { gt: 100 } } } },
 orderBy: { id: "desc" },
});
```

### Filter records

```ts
// Single field filter
const users = await prisma.user.findMany({
 where: { email: { endsWith: "prisma.io" } },
});

// Multiple conditions with OR/AND
const users = await prisma.user.findMany({
 where: {
 OR: [{ name: { startsWith: "E" } }, { AND: { profileViews: { gt: 0 }, role: "ADMIN" } }],
 },
});

// Filter by related records
const users = await prisma.user.findMany({
 where: {
 email: { endsWith: "prisma.io" },
 posts: { some: { published: false } },
 },
});
```

See [Filtering and sorting](/orm/prisma-client/queries/filtering-and-sorting) for more examples.

### Select fields

```ts
const user = await prisma.user.findUnique({
 where: { email: "emma@prisma.io" },
 select: { email: true, name: true },
});
// Returns: { email: 'emma@prisma.io', name: "Emma" }
```

### Include related records

```ts
const users = await prisma.user.findMany({
 where: { role: "ADMIN" },
 include: { posts: true },
});
```

See [Select fields](/orm/prisma-client/queries/select-fields) and [Relation queries](/orm/prisma-client/queries/relation-queries) for more.

## Update

### Update a single record

```ts
const updateUser = await prisma.user.update({
 where: { email: "viola@prisma.io" },
 data: { name: "Viola the Magnificent" },
});
```

### Update multiple records

```ts
const updateUsers = await prisma.user.updateMany({
 where: { email: { contains: "prisma.io" } },
 data: { role: "ADMIN" },
});
// Returns: { count: 19 }
```

### Update and return multiple records

Supported by PostgreSQL, CockroachDB, and SQLite.

```ts
const users = await prisma.user.updateManyAndReturn({
 where: { email: { contains: "prisma.io" } },
 data: { role: "ADMIN" },
});
```

### Upsert (update or create)

```ts
const upsertUser = await prisma.user.upsert({
 where: { email: "viola@prisma.io" },
 update: { name: "Viola the Magnificent" },
 create: { email: "viola@prisma.io", name: "Viola the Magnificent" },
});
```

:::tip
To emulate `findOrCreate()`, use `upsert()` with an empty `update` parameter.
:::

### Atomic number operations

```ts
await prisma.post.updateMany({
 data: {
 views: { increment: 1 },
 likes: { increment: 1 },
 },
});
```

See [Relation queries](/orm/prisma-client/queries/relation-queries) for connecting and disconnecting related records.

## Delete

### Delete a single record

The following query uses [`delete()`](/orm/reference/prisma-client-reference#delete) to delete a single `User` record:

```ts
const deleteUser = await prisma.user.delete({
 where: {
 email: "bert@prisma.io",
 },
});
```

Attempting to delete a user with one or more posts result in an error, as every `Post` requires an author - see [cascading deletes](#cascading-deletes-deleting-related-records).

### Delete multiple records

The following query uses [`deleteMany()`](/orm/reference/prisma-client-reference#deletemany) to delete all `User` records where `email` contains `prisma.io`:

```ts
const deleteUsers = await prisma.user.deleteMany({
 where: {
 email: {
 contains: "prisma.io",
 },
 },
});
```

Attempting to delete a user with one or more posts result in an error, as every `Post` requires an author - see [cascading deletes](#cascading-deletes-deleting-related-records).

### Delete all records

The following query uses [`deleteMany()`](/orm/reference/prisma-client-reference#deletemany) to delete all `User` records:

```ts
const deleteUsers = await prisma.user.deleteMany({});
```

Be aware that this query will fail if the user has any related records (such as posts). In this case, you need to [delete the related records first](#cascading-deletes-deleting-related-records).

### Cascading deletes (deleting related records)

:::tip

You can configure cascading deletes using [referential actions](/orm/prisma-schema/data-model/relations/referential-actions).

:::

The following query uses [`delete()`](/orm/reference/prisma-client-reference#delete) to delete a single `User` record:

```ts
const deleteUser = await prisma.user.delete({
 where: {
 email: "bert@prisma.io",
 },
});
```

However, the example schema includes a **required relation** between `Post` and `User`, which means that you cannot delete a user with posts:

```
The change you are trying to make would violate the required relation 'PostToUser' between the `Post` and `User` models.
```

To resolve this error, you can:

- Make the relation optional:

 ```prisma highlight=3,4;add|5,6;delete
 model Post {
 id Int @id @default(autoincrement())
 author User? @relation(fields: [authorId], references: [id]) // [!code ++]
 authorId Int? // [!code ++]
 author User @relation(fields: [authorId], references: [id]) // [!code --]
 authorId Int // [!code --]
 }
 ```

- Change the author of the posts to another user before deleting the user.

- Delete a user and all their posts with two separate queries in a transaction (all queries must succeed):

 ```ts
 const deletePosts = prisma.post.deleteMany({
 where: {
 authorId: 7,
 },
 });

 const deleteUser = prisma.user.delete({
 where: {
 id: 7,
 },
 });

 const transaction = await prisma.$transaction([deletePosts, deleteUser]);
 ```

### Delete all records from all tables

Sometimes you want to remove all data from all tables but keep the actual tables. This can be particularly useful in a development environment and whilst testing.

The following shows how to delete all records from all tables with Prisma Client and with Prisma Migrate.

#### Deleting all data with `deleteMany()`

When you know the order in which your tables should be deleted, you can use the [`deleteMany`](/orm/reference/prisma-client-reference#deletemany) function. This is executed synchronously in a [`$transaction`](/orm/prisma-client/queries/transactions) and can be used with all types of databases.

```ts
const deletePosts = prisma.post.deleteMany();
const deleteProfile = prisma.profile.deleteMany();
const deleteUsers = prisma.user.deleteMany();

// The transaction runs synchronously so deleteUsers must run last.
await prisma.$transaction([deleteProfile, deletePosts, deleteUsers]);
```

✅ **Pros**:

- Works well when you know the structure of your schema ahead of time
- Synchronously deletes each tables data

❌ **Cons**:

- When working with relational databases, this function doesn't scale as well as having a more generic solution which looks up and `TRUNCATE`s your tables regardless of their relational constraints. Note that this scaling issue does not apply when using the MongoDB connector.

> **Note**: The `$transaction` performs a cascading delete on each models table so they have to be called in order.

#### Deleting all data with raw SQL / `TRUNCATE`

If you are comfortable working with raw SQL, you can perform a `TRUNCATE` query on a table using [`$executeRawUnsafe`](/orm/prisma-client/using-raw-sql/raw-queries#executerawunsafe).

In the following examples, the first tab shows how to perform a `TRUNCATE` on a Postgres database by using a `$queryRaw` look up that maps over the table and `TRUNCATES` all tables in a single query.

The second tab shows performing the same function but with a MySQL database. In this instance the constraints must be removed before the `TRUNCATE` can be executed, before being reinstated once finished. The whole process is run as a `$transaction`

```ts tab="PostgreSQL"
const tablenames = await prisma.$queryRaw<
 Array<{ tablename: string }>
>`SELECT tablename FROM pg_tables WHERE schemaname='public'`;

const tables = tablenames
 .map(({ tablename }) => tablename)
 .filter((name) => name !== "_prisma_migrations")
 .map((name) => `"public"."${name}"`)
 .join(", ");

try {
 await prisma.$executeRawUnsafe(`TRUNCATE TABLE ${tables} CASCADE;`);
} catch (error) {
 console.log({ error });
}
```

```ts tab="MySQL"
const transactions: PrismaPromise<any>[] = [];
transactions.push(prisma.$executeRaw`SET FOREIGN_KEY_CHECKS = 0;`);

const tablenames = await prisma.$queryRaw<
 Array<{ TABLE_NAME: string }>
>`SELECT TABLE_NAME from information_schema.TABLES WHERE TABLE_SCHEMA = 'tests';`;

for (const { TABLE_NAME } of tablenames) {
 if (TABLE_NAME !== "_prisma_migrations") {
 try {
 transactions.push(prisma.$executeRawUnsafe(`TRUNCATE ${TABLE_NAME};`));
 } catch (error) {
 console.log({ error });
 }
 }
}

transactions.push(prisma.$executeRaw`SET FOREIGN_KEY_CHECKS = 1;`);

try {
 await prisma.$transaction(transactions);
} catch (error) {
 console.log({ error });
}
```

✅ **Pros**:

- Scalable
- Very fast

❌ **Cons**:

- Can't undo the operation
- Using reserved SQL key words as tables names can cause issues when trying to run a raw query

#### Deleting all records with Prisma Migrate

If you use Prisma Migrate, you can use `migrate reset`, this will:

1. Drop the database
2. Create a new database
3. Apply migrations
4. Seed the database with data

## Advanced query examples

### Create a deeply nested tree of records

- A single `User`
- Two new, related `Post` records
- Connect or create `Category` per post

```ts
const u = await prisma.user.create({
 include: {
 posts: {
 include: {
 categories: true,
 },
 },
 },
 data: {
 email: "emma@prisma.io",
 posts: {
 create: [
 {
 title: "My first post",
 categories: {
 connectOrCreate: [
 {
 create: { name: "Introductions" },
 where: {
 name: "Introductions",
 },
 },
 {
 create: { name: "Social" },
 where: {
 name: "Social",
 },
 },
 ],
 },
 },
 {
 title: "How to make cookies",
 categories: {
 connectOrCreate: [
 {
 create: { name: "Social" },
 where: {
 name: "Social",
 },
 },
 {
 create: { name: "Cooking" },
 where: {
 name: "Cooking",
 },
 },
 ],
 },
 },
 ],
 },
 },
});
```

---
title: Excluding fields
description: Learn how to exclude fields from Prisma Client results with the omit option.
url: /orm/prisma-client/queries/excluding-fields
metaTitle: Excluding fields
metaDescription: Learn how to use omit in Prisma Client to exclude fields globally or per query.
---

Use `omit` when you want Prisma Client to return the normal result shape except for a few specific fields.

## Omit a field for one query

```ts
const user = await prisma.user.findUnique({
 where: { id: 1 },
 omit: {
 password: true,
 },
});
```

## Omit fields globally

You can also configure `omit` on the client itself:

```ts
const prisma = new PrismaClient({
 omit: {
 user: {
 password: true,
 },
 },
});
```

## When to use omit vs select

- Use [`select`](/orm/prisma-client/queries/select-fields) when you want to return only a small subset of fields.
- Use `omit` when the default result is mostly correct and you only want to remove a few sensitive or noisy fields.

## Related pages

- [Select fields](/orm/prisma-client/queries/select-fields)
- [Prisma Client API reference](/orm/reference/prisma-client-reference#omit)

---
title: Filtering and sorting
description: Learn how to filter Prisma Client queries with where and sort results with orderBy.
url: /orm/prisma-client/queries/filtering-and-sorting
metaTitle: Filtering and sorting
metaDescription: Learn how to filter Prisma Client results with where, combine operators, use relation filters, and sort with orderBy.
---

Prisma Client lets you narrow results with `where` and order them with `orderBy`.

## Filtering with where

Use `where` to match records by field values:

```ts
const users = await prisma.user.findMany({
 where: {
 email: {
 endsWith: "prisma.io",
 },
 },
});
```

## Combining operators

You can compose filters with operators such as `OR`, `AND`, and `NOT`:

```ts
const users = await prisma.user.findMany({
 where: {
 OR: [
 { email: { endsWith: "gmail.com" } },
 { email: { endsWith: "company.com" } },
 ],
 NOT: {
 email: {
 endsWith: "admin.company.com",
 },
 },
 },
});
```

## Filter on related records

Relation filters let you match records based on related data:

```ts
const users = await prisma.user.findMany({
 where: {
 posts: {
 some: {
 published: true,
 },
 },
 },
});
```

For more relation-specific patterns, see [Relation queries](/orm/prisma-client/queries/relation-queries).

## Sort results with orderBy

Use `orderBy` to control result ordering:

```ts
const posts = await prisma.post.findMany({
 orderBy: {
 title: "asc",
 },
});
```

You can also combine filtering and sorting:

```ts
const posts = await prisma.post.findMany({
 where: {
 published: true,
 },
 orderBy: {
 createdAt: "desc",
 },
});
```

## Case-insensitive filtering

Case sensitivity depends on your database provider and collation settings. For PostgreSQL, Prisma Client also supports specific case-insensitive filter modes on supported operators. See the [Prisma Client API reference](/orm/reference/prisma-client-reference#mode) for details.

## Sort by relation

You can sort by properties on related records when the query shape supports it. For example, you might order posts by their author's name or users by related aggregates.

## Sort by relevance (PostgreSQL and MySQL)

On supported databases, Prisma Client can sort search results by relevance using `_relevance`. This is especially useful when combined with [full-text search](/orm/prisma-client/queries/full-text-search).

## Sort with null records first or last

Prisma Client supports explicit null ordering on supported databases so you can keep incomplete values grouped at the beginning or end of a result set.

## Related pages

- [Pagination](/orm/prisma-client/queries/pagination)
- [Select fields](/orm/prisma-client/queries/select-fields)
- [Prisma Client API reference](/orm/reference/prisma-client-reference#filter-conditions-and-operators)

---
title: Full-text search
description: Learn how to search text fields with Prisma Client using your database's native full-text search support.
badge: preview
url: /orm/prisma-client/queries/full-text-search
metaTitle: Full-text search
metaDescription: Learn how to use full-text search with Prisma Client for PostgreSQL and MySQL.
---

Prisma Client supports full-text search for MySQL and for PostgreSQL with the `fullTextSearchPostgres` preview feature.

## Enabling full-text search for PostgreSQL

Add the preview flag to your generator and re-generate Prisma Client:

```prisma title="schema.prisma"
generator client {
 provider = "prisma-client"
 output = "./generated"
 previewFeatures = ["fullTextSearchPostgres"]
}
```

```npm
npx prisma generate
```

## Search within a text field

```ts
const posts = await prisma.post.findMany({
 where: {
 body: {
 search: "cat | dog",
 },
 },
});
```

## Sort by relevance

Prisma Client also supports relevance-based ordering on supported databases:

```ts
const posts = await prisma.post.findMany({
 orderBy: {
 _relevance: {
 fields: ["title"],
 search: "database",
 sort: "desc",
 },
 },
});
```

## Related pages

- [Filtering and sorting](/orm/prisma-client/queries/filtering-and-sorting)
- [Prisma Client API reference](/orm/reference/prisma-client-reference#search)
- [Preview features](/orm/reference/preview-features/client-preview-features)

---
title: Pagination
description: Learn how to paginate Prisma Client query results with offset pagination and cursor-based pagination.
url: /orm/prisma-client/queries/pagination
metaTitle: Pagination
metaDescription: Learn how to use offset pagination and cursor-based pagination with Prisma Client.
---

Prisma Client supports both offset pagination and cursor-based pagination.

## Offset pagination

Use `skip` and `take` when you need page numbers or shallow navigation through a result set:

```ts
const posts = await prisma.post.findMany({
 skip: 20,
 take: 10,
});
```

Offset pagination is straightforward, but it becomes more expensive as the offset grows.

## Cursor-based pagination

Use `cursor` and `take` when you want stable, scalable pagination for feeds, timelines, or large datasets:

```ts
const firstPage = await prisma.post.findMany({
 take: 10,
 orderBy: {
 id: "asc",
 },
});

const lastPost = firstPage[firstPage.length - 1];

const nextPage = lastPost
 ? await prisma.post.findMany({
 take: 10,
 skip: 1,
 cursor: {
 id: lastPost.id,
 },
 orderBy: {
 id: "asc",
 },
 })
 : [];
```

## Which approach to choose

- Use offset pagination when users need to jump directly to a numbered page.
- Use cursor-based pagination when you care more about performance and consistency as the dataset grows.

## Related pages

- [Filtering and sorting](/orm/prisma-client/queries/filtering-and-sorting)
- [Query Insights](/query-insights)
- [Prisma Client API reference](/orm/reference/prisma-client-reference)

---
title: Relation queries
description: 'Prisma Client provides convenient queries for working with relations, such as a fluent API, nested writes (transactions), nested reads and relation filters'
url: /orm/prisma-client/queries/relation-queries
metaTitle: Relation queries (Concepts)
metaDescription: 'Prisma Client provides convenient queries for working with relations, such as a fluent API, nested writes (transactions), nested reads and relation filters.'
---

A key feature of Prisma Client is the ability to query [relations](/orm/prisma-schema/data-model/relations) between two or more models. Relation queries include:

- [Nested reads](#nested-reads) (sometimes referred to as _eager loading_) via [`select`](/orm/reference/prisma-client-reference#select) and [`include`](/orm/reference/prisma-client-reference#include)
- [Nested writes](#nested-writes) with [transactional](/orm/prisma-client/queries/transactions) guarantees
- [Filtering on related records](#relation-filters)

Prisma Client also has a [fluent API for traversing relations](#fluent-api).

## Nested reads

Nested reads allow you to read related data from multiple tables in your database - such as a user and that user's posts. You can:

- Use [`include`](/orm/reference/prisma-client-reference#include) to include related records, such as a user's posts or profile, in the query response.
- Use a nested [`select`](/orm/reference/prisma-client-reference#select) to include specific fields from a related record. You can also nest `select` inside an `include`.

### Relation load strategies (Preview)

You can decide on a per-query-level _how_ you want Prisma Client to execute a relation query (i.e. what _load strategy_ should be applied) via the `relationLoadStrategy` option for PostgreSQL databases.

Because the `relationLoadStrategy` option is currently in Preview, you need to enable it via the `relationJoins` preview feature flag in your Prisma schema file:

```prisma title="schema.prisma" showLineNumbers
generator client {
 provider = "prisma-client"
 output = "./generated"
 previewFeatures = ["relationJoins"]
}
```

After adding this flag, you need to run `prisma generate` again to re-generate Prisma Client. The `relationJoins` feature is currently available on PostgreSQL, CockroachDB and MySQL.

Prisma Client supports two load strategies for relations:

- `join` (default): Uses a database-level `LATERAL JOIN` (PostgreSQL) or correlated subqueries (MySQL) and fetches all data with a single query to the database.
- `query`: Sends multiple queries to the database (one per table) and joins them on the application level.

Another important difference between these two options is that the `join` strategy uses JSON aggregation on the database level. That means that it creates the JSON structures returned by Prisma Client already in the database which saves computation resources on the application level.

#### Examples

You can use the `relationLoadStrategy` option on the top-level in any query that supports `include` or `select`.

Here is an example with `include`:

```ts
const users = await prisma.user.findMany({
 relationLoadStrategy: "join", // or 'query'
 include: {
 posts: true,
 },
});
```

And here is another example with `select`:

```ts
const users = await prisma.user.findMany({
 relationLoadStrategy: "join", // or 'query'
 select: {
 posts: true,
 },
});
```

#### When to use which load strategy?

- The `join` strategy (default) will be more effective in most scenarios. On PostgreSQL, it uses a combination of `LATERAL JOINs` and JSON aggregation to reduce redundancy in result sets and delegate the work of transforming the query results into the expected JSON structures on the database server. On MySQL, it uses correlated subqueries to fetch the results with a single query.
- There may be edge cases where `query` could be more performant depending on the characteristics of the dataset and query. We recommend that you profile your database queries to identify these situations.
- Use `query` if you want to save resources on the database server and do heavy-lifting of merging and transforming data in the application server which might be easier to scale.

### Include a relation

The following example returns a single user and that user's posts:

```ts
const user = await prisma.user.findFirst({
 include: {
 posts: true,
 },
});
```

```json
{
 id: 19,
 name: null,
 email: 'emma@prisma.io',
 profileViews: 0,
 role: 'USER',
 coinflips: [],
 posts: [
 {
 id: 20,
 title: 'My first post',
 published: true,
 authorId: 19,
 comments: null,
 views: 0,
 likes: 0
 },
 {
 id: 21,
 title: 'How to make cookies',
 published: true,
 authorId: 19,
 comments: null,
 views: 0,
 likes: 0
 }
 ]
}
```

### Include all fields for a specific relation

The following example returns a post and its author:

```ts
const post = await prisma.post.findFirst({
 include: {
 author: true,
 },
});
```

```json
{
 id: 17,
 title: 'How to make cookies',
 published: true,
 authorId: 16,
 comments: null,
 views: 0,
 likes: 0,
 author: {
 id: 16,
 name: null,
 email: 'orla@prisma.io',
 profileViews: 0,
 role: 'USER',
 coinflips: [],
 },
}
```

### Include deeply nested relations

You can nest `include` options to include relations of relations. The following example returns a user's posts, and each post's categories:

```ts
const user = await prisma.user.findFirst({
 include: {
 posts: {
 include: {
 categories: true,
 },
 },
 },
});
```

```json
{
 "id": 40,
 "name": "Yvette",
 "email": "yvette@prisma.io",
 "profileViews": 0,
 "role": "USER",
 "coinflips": [],
 "testing": [],
 "city": null,
 "country": "Sweden",
 "posts": [
 {
 "id": 66,
 "title": "How to make an omelette",
 "published": true,
 "authorId": 40,
 "comments": null,
 "views": 0,
 "likes": 0,
 "categories": [
 {
 "id": 3,
 "name": "Easy cooking"
 }
 ]
 },
 {
 "id": 67,
 "title": "How to eat an omelette",
 "published": true,
 "authorId": 40,
 "comments": null,
 "views": 0,
 "likes": 0,
 "categories": []
 }
 ]
}
```

### Select specific fields of included relations

You can use a nested `select` to choose a subset of fields of relations to return. For example, the following query returns the user's `name` and the `title` of each related post:

```ts
const user = await prisma.user.findFirst({
 select: {
 name: true,
 posts: {
 select: {
 title: true,
 },
 },
 },
});
```

```json
{
 name: "Elsa",
 posts: [ { title: 'My first post' }, { title: 'How to make cookies' } ]
}
```

You can also nest a `select` inside an `include` - the following example returns _all_ `User` fields and the `title` field of each post:

```ts
const user = await prisma.user.findFirst({
 include: {
 posts: {
 select: {
 title: true,
 },
 },
 },
});
```

```json
{
 "id": 1,
 "name": null,
 "email": "martina@prisma.io",
 "profileViews": 0,
 "role": "USER",
 "coinflips": [],
 "posts": [
 { "title": "How to grow salad" },
 { "title": "How to ride a horse" }
 ]
}
```

Note that you **cannot** use `select` and `include` _on the same level_. This means that if you choose to `include` a user's post and `select` each post's title, you cannot `select` only the users' `email`:

```ts
// The following query returns an exception
const user = await prisma.user.findFirst({
 select: { // This won't work! // [!code --]
 email: true
 }
 include: { // This won't work! // [!code --]
 posts: {
 select: {
 title: true
 }
 }
 },
})
```

```text no-copy
Invalid `prisma.user.findUnique()` invocation:

{
 where: {
 id: 19
 },
 select: {
 ~~~~~~
 email: true
 },
 include: {
 ~~~~~~~
 posts: {
 select: {
 title: true
 }
 }
 }
}

Please either use `include` or `select`, but not both at the same time.
```

Instead, use nested `select` options:

```ts
const user = await prisma.user.findFirst({
 select: {
 // This will work!
 email: true,
 posts: {
 select: {
 title: true,
 },
 },
 },
});
```

## Relation count

In [3.0.1](https://github.com/prisma/prisma/releases/3.0.1) and later, you can [`include` or `select` a count of relations](/orm/prisma-client/queries/aggregation-grouping-summarizing#count-relations) alongside fields - for example, a user's post count.

```ts
const relationCount = await prisma.user.findMany({
 include: {
 _count: {
 select: { posts: true },
 },
 },
});
```

```text no-copy
{ id: 1, _count: { posts: 3 } },
{ id: 2, _count: { posts: 2 } },
{ id: 3, _count: { posts: 2 } },
{ id: 4, _count: { posts: 0 } },
{ id: 5, _count: { posts: 0 } }
```

## Filter a list of relations

When you use `select` or `include` to return a subset of the related data, you can **filter and sort the list of relations** inside the `select` or `include`.

For example, the following query returns list of titles of the unpublished posts associated with the user:

```ts
const result = await prisma.user.findFirst({
 select: {
 posts: {
 where: {
 published: false,
 },
 orderBy: {
 title: "asc",
 },
 select: {
 title: true,
 },
 },
 },
});
```

You can also write the same query using `include` as follows:

```ts
const result = await prisma.user.findFirst({
 include: {
 posts: {
 where: {
 published: false,
 },
 orderBy: {
 title: "asc",
 },
 },
 },
});
```

## Nested writes

A nested write allows you to write **relational data** to your database in **a single transaction**.

Nested writes:

- Provide **transactional guarantees** for creating, updating or deleting data across multiple tables in a single Prisma Client query. If any part of the query fails (for example, creating a user succeeds but creating posts fails), Prisma Client rolls back all changes.
- Support any level of nesting supported by the data model.
- Are available for [relation fields](/orm/prisma-schema/data-model/relations#relation-fields) when using the model's create or update query. The following section shows the nested write options that are available per query.

### Create a related record

You can create a record and one or more related records at the same time. The following query creates a `User` record and two related `Post` records:

```ts
const result = await prisma.user.create({
 data: {
 email: "elsa@prisma.io",
 name: "Elsa Prisma",
 posts: {
 // [!code highlight]
 create: [{ title: "How to make an omelette" }, { title: "How to eat an omelette" }], // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true, // Include all posts in the returned object
 },
});
```

```json
{
 id: 29,
 name: 'Elsa',
 email: 'elsa@prisma.io',
 profileViews: 0,
 role: 'USER',
 coinflips: [],
 posts: [
 {
 id: 22,
 title: 'How to make an omelette',
 published: true,
 authorId: 29,
 comments: null,
 views: 0,
 likes: 0
 },
 {
 id: 23,
 title: 'How to eat an omelette',
 published: true,
 authorId: 29,
 comments: null,
 views: 0,
 likes: 0
 }
 ]
}
```

### Create a single record and multiple related records

There are two ways to create or update a single record and multiple related records - for example, a user with multiple posts:

- Use a nested [`create`](/orm/reference/prisma-client-reference#create) query
- Use a nested [`createMany`](/orm/reference/prisma-client-reference#nested-createmany-options) query

In most cases, a nested `create` will be preferable unless the [`skipDuplicates` query option](/orm/reference/prisma-client-reference#nested-createmany-options) is required. Here's a quick table describing the differences between the two options:

| Feature | `create` | `createMany` | Notes |
| :------------------------------------ | :------- | :----------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Supports nesting additional relations | ✔ | ✘ \* | For example, you can create a user, several posts, and several comments per post in one query.<br />\* You can manually set a foreign key in a has-one relation - for example: `{ authorId: 9}` |
| Supports 1-n relations | ✔ | ✔ | For example, you can create a user and multiple posts (one user has many posts) |
| Supports m-n relations | ✔ | ✘ | For example, you can create a post and several categories (one post can have many categories, and one category can have many posts) |
| Supports skipping duplicate records | ✘ | ✔ | Use `skipDuplicates` query option. |

#### Using nested `create`

The following query uses nested [`create`](/orm/reference/prisma-client-reference#create) to create:

- One user
- Two posts
- One post category

The example also uses a nested `include` to include all posts and post categories in the returned data.

```ts
const result = await prisma.user.create({
 data: {
 email: "yvette@prisma.io",
 name: "Yvette",
 posts: {
 // [!code highlight]
 create: [
 // [!code highlight]
 {
 // [!code highlight]
 title: "How to make an omelette", // [!code highlight]
 categories: {
 // [!code highlight]
 create: {
 // [!code highlight]
 name: "Easy cooking", // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 { title: "How to eat an omelette" }, // [!code highlight]
 ], // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 // Include posts
 posts: {
 include: {
 categories: true, // Include post categories
 },
 },
 },
});
```

```json
{
 "id": 40,
 "name": "Yvette",
 "email": "yvette@prisma.io",
 "profileViews": 0,
 "role": "USER",
 "coinflips": [],
 "testing": [],
 "city": null,
 "country": "Sweden",
 "posts": [
 {
 "id": 66,
 "title": "How to make an omelette",
 "published": true,
 "authorId": 40,
 "comments": null,
 "views": 0,
 "likes": 0,
 "categories": [
 {
 "id": 3,
 "name": "Easy cooking"
 }
 ]
 },
 {
 "id": 67,
 "title": "How to eat an omelette",
 "published": true,
 "authorId": 40,
 "comments": null,
 "views": 0,
 "likes": 0,
 "categories": []
 }
 ]
}
```

Here's a visual representation of how a nested create operation can write to several tables in the database as once:

![Diagram showing how a nested create operation writes to multiple database tables (User, Post, Category) in a single transaction.](/img/orm/nested-create.png)

#### Using nested `createMany`

The following query uses a nested [`createMany`](/orm/reference/prisma-client-reference#createmany) to create:

- One user
- Two posts

The example also uses a nested `include` to include all posts in the returned data.

```ts
const result = await prisma.user.create({
 data: {
 email: "saanvi@prisma.io",
 posts: {
 // [!code highlight]
 createMany: {
 // [!code highlight]
 data: [{ title: "My first post" }, { title: "My second post" }], // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true,
 },
});
```

```json
{
 "id": 43,
 "name": null,
 "email": "saanvi@prisma.io",
 "profileViews": 0,
 "role": "USER",
 "coinflips": [],
 "testing": [],
 "city": null,
 "country": "India",
 "posts": [
 {
 "id": 70,
 "title": "My first post",
 "published": true,
 "authorId": 43,
 "comments": null,
 "views": 0,
 "likes": 0
 },
 {
 "id": 71,
 "title": "My second post",
 "published": true,
 "authorId": 43,
 "comments": null,
 "views": 0,
 "likes": 0
 }
 ]
}
```

Note that it is **not possible** to nest an additional `create` or `createMany` inside the highlighted query, which means that you cannot create a user, posts, and post categories at the same time.

As a workaround, you can send a query to create the records that will be connected first, and then create the actual records. For example:

```ts
const categories = await prisma.category.createManyAndReturn({
 data: [{ name: "Fun" }, { name: "Technology" }, { name: "Sports" }],
 select: {
 id: true,
 },
});

const posts = await prisma.post.createManyAndReturn({
 data: [
 {
 title: "Funniest moments in 2024",
 categoryId: categories.find((category) => category.name === "Fun")!.id,
 },
 {
 title: "Linux or macOS — what's better?",
 categoryId: categories.find((category) => category.name === "Technology")!.id,
 },
 {
 title: "Who will win the next soccer championship?",
 categoryId: categories.find((category) => category.name === "Sports")!.id,
 },
 ],
});
```

If you want to create _all_ records in a single database query, consider using a [`$transaction`](/orm/prisma-client/queries/transactions#the-transaction-api) or [type-safe, raw SQL](/orm/prisma-client/using-raw-sql/typedsql).

### Create multiple records and multiple related records

You cannot access relations in a `createMany()` or `createManyAndReturn()` query, which means that you cannot create multiple users and multiple posts in a single nested write. The following is **not** possible:

```ts
const createMany = await prisma.user.createMany({
 data: [
 {
 name: "Yewande",
 email: "yewande@prisma.io",
 posts: {
 // [!code --]
 // Not possible to create posts! // [!code --]
 }, // [!code --]
 },
 {
 name: "Noor",
 email: "noor@prisma.io",
 posts: {
 // [!code --]
 // Not possible to create posts! // [!code --]
 }, // [!code --]
 },
 ],
});
```

### Connect multiple records

The following query creates ([`create`](/orm/reference/prisma-client-reference#create) ) a new `User` record and connects that record ([`connect`](/orm/reference/prisma-client-reference#connect) ) to three existing posts:

```ts
const result = await prisma.user.create({
 data: {
 email: "vlad@prisma.io",
 posts: {
 // [!code highlight]
 connect: [{ id: 8 }, { id: 9 }, { id: 10 }], // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true, // Include all posts in the returned object
 },
});
```

```json
{
 id: 27,
 name: null,
 email: 'vlad@prisma.io',
 profileViews: 0,
 role: 'USER',
 coinflips: [],
 posts: [
 {
 id: 10,
 title: 'An existing post',
 published: true,
 authorId: 27,
 comments: {},
 views: 0,
 likes: 0
 }
 ]
}
```

:::info[Note]
Prisma Client throws an exception if any of the post records cannot be found: `connect: [{ id: 8 }, { id: 9 }, { id: 10 }]`
:::

### Connect a single record

You can [`connect`](/orm/reference/prisma-client-reference#connect) an existing record to a new or existing user. The following query connects an existing post (`id: 11`) to an existing user (`id: 9`)

```ts
const result = await prisma.user.update({
 where: {
 id: 9,
 },
 data: {
 posts: {
 // [!code highlight]
 connect: {
 // [!code highlight]
 id: 11, // [!code highlight]
 }, // [!code highlight]
 },
 },
 include: {
 posts: true,
 },
});
```

### Connect _or_ create a record

If a related record may or may not already exist, use [`connectOrCreate`](/orm/reference/prisma-client-reference#connectorcreate) to connect the related record:

- Connect a `User` with the email address `viola@prisma.io` _or_
- Create a new `User` with the email address `viola@prisma.io` if the user does not already exist

```ts
const result = await prisma.post.create({
 data: {
 title: "How to make croissants",
 author: {
 // [!code highlight]
 connectOrCreate: {
 // [!code highlight]
 where: {
 // [!code highlight]
 email: "viola@prisma.io", // [!code highlight]
 }, // [!code highlight]
 create: {
 // [!code highlight]
 email: "viola@prisma.io", // [!code highlight]
 name: "Viola", // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 author: true,
 },
});
```

```json
{
 id: 26,
 title: 'How to make croissants',
 published: true,
 authorId: 43,
 views: 0,
 likes: 0,
 author: {
 id: 43,
 name: 'Viola',
 email: 'viola@prisma.io',
 profileViews: 0,
 role: 'USER',
 coinflips: []
 }
}
```

### Disconnect a related record

To `disconnect` one out of a list of records (for example, a specific blog post) provide the ID or unique identifier of the record(s) to disconnect:

```ts
const result = await prisma.user.update({
 where: {
 id: 16,
 },
 data: {
 posts: {
 // [!code highlight]
 disconnect: [{ id: 12 }, { id: 19 }], // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true,
 },
});
```

```json
{
 id: 16,
 name: null,
 email: 'orla@prisma.io',
 profileViews: 0,
 role: 'USER',
 coinflips: [],
 posts: []
}
```

To `disconnect` _one_ record (for example, a post's author), use `disconnect: true`:

```ts
const result = await prisma.post.update({
 where: {
 id: 23,
 },
 data: {
 author: {
 // [!code highlight]
 disconnect: true, // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 author: true,
 },
});
```

```json
{
 id: 23,
 title: 'How to eat an omelette',
 published: true,
 authorId: null,
 comments: null,
 views: 0,
 likes: 0,
 author: null
}
```

### Disconnect all related records

To [`disconnect`](/orm/reference/prisma-client-reference#disconnect) _all_ related records in a one-to-many relation (a user has many posts), `set` the relation to an empty list as shown:

```ts
const result = await prisma.user.update({
 where: {
 id: 16,
 },
 data: {
 posts: {
 // [!code highlight]
 set: [], // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true,
 },
});
```

```json
{
 id: 16,
 name: null,
 email: 'orla@prisma.io',
 profileViews: 0,
 role: 'USER',
 coinflips: [],
 posts: []
}
```

### Delete all related records

Delete all related `Post` records:

```ts
const result = await prisma.user.update({
 where: {
 id: 11,
 },
 data: {
 posts: {
 // [!code highlight]
 deleteMany: {}, // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true,
 },
});
```

### Delete specific related records

Update a user by deleting all unpublished posts:

```ts
const result = await prisma.user.update({
 where: {
 id: 11,
 },
 data: {
 posts: {
 // [!code highlight]
 deleteMany: {
 // [!code highlight]
 published: false, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true,
 },
});
```

Update a user by deleting specific posts:

```ts
const result = await prisma.user.update({
 where: {
 id: 6,
 },
 data: {
 posts: {
 // [!code highlight]
 deleteMany: [{ id: 7 }], // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true,
 },
});
```

### Update all related records (or filter)

You can use a nested `updateMany` to update _all_ related records for a particular user. The following query unpublishes all posts for a specific user:

```ts
const result = await prisma.user.update({
 where: {
 id: 6,
 },
 data: {
 posts: {
 // [!code highlight]
 updateMany: {
 // [!code highlight]
 where: {
 // [!code highlight]
 published: true, // [!code highlight]
 }, // [!code highlight]
 data: {
 // [!code highlight]
 published: false, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true,
 },
});
```

### Update a specific related record

```ts
const result = await prisma.user.update({
 where: {
 id: 6,
 },
 data: {
 posts: {
 // [!code highlight]
 update: {
 // [!code highlight]
 where: {
 // [!code highlight]
 id: 9, // [!code highlight]
 }, // [!code highlight]
 data: {
 // [!code highlight]
 title: "My updated title", // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true,
 },
});
```

### Update _or_ create a related record

The following query uses a nested `upsert` to update `"bob@prisma.io"` if that user exists, or create the user if they do not exist:

```ts
const result = await prisma.post.update({
 where: {
 id: 6,
 },
 data: {
 author: {
 // [!code highlight]
 upsert: {
 // [!code highlight]
 create: {
 // [!code highlight]
 email: "bob@prisma.io", // [!code highlight]
 name: "Bob the New User", // [!code highlight]
 }, // [!code highlight]
 update: {
 // [!code highlight]
 email: "bob@prisma.io", // [!code highlight]
 name: "Bob the existing user", // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 author: true,
 },
});
```

### Add new related records to an existing record

You can nest `create` or `createMany` inside an `update` to add new related records to an existing record. The following query adds two posts to a user with an `id` of 9:

```ts
const result = await prisma.user.update({
 where: {
 id: 9,
 },
 data: {
 posts: {
 // [!code highlight]
 createMany: {
 // [!code highlight]
 data: [{ title: "My first post" }, { title: "My second post" }], // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true,
 },
});
```

## Relation filters

### Filter on "-to-many" relations

Prisma Client provides the [`some`](/orm/reference/prisma-client-reference#some), [`every`](/orm/reference/prisma-client-reference#every), and [`none`](/orm/reference/prisma-client-reference#none) options to filter records by the properties of related records on the "-to-many" side of the relation. For example, filtering users based on properties of their posts.

For example:

| Requirement | Query option to use |
| --------------------------------------------------------------------------------- | ----------------------------------- |
| "I want a list of every `User` that has _at least one_ unpublished `Post` record" | `some` posts are unpublished |
| "I want a list of every `User` that has _no_ unpublished `Post` records" | `none` of the posts are unpublished |
| "I want a list of every `User` that has _only_ unpublished `Post` records" | `every` post is unpublished |

For example, the following query returns `User` that meet the following criteria:

- No posts with more than 100 views
- All posts have less than, or equal to 50 likes

```ts
const users = await prisma.user.findMany({
 where: {
 posts: {
 // [!code highlight]
 none: {
 // [!code highlight]
 views: {
 // [!code highlight]
 gt: 100, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 every: {
 // [!code highlight]
 likes: {
 // [!code highlight]
 lte: 50, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true,
 },
});
```

### Filter on "-to-one" relations

Prisma Client provides the [`is`](/orm/reference/prisma-client-reference#is) and [`isNot`](/orm/reference/prisma-client-reference#isnot) options to filter records by the properties of related records on the "-to-one" side of the relation. For example, filtering posts based on properties of their author.

For example, the following query returns `Post` records that meet the following criteria:

- Author's name is not Bob
- Author is older than 40

```ts
const users = await prisma.post.findMany({
 where: {
 author: {
 // [!code highlight]
 isNot: {
 // [!code highlight]
 name: "Bob", // [!code highlight]
 }, // [!code highlight]
 is: {
 // [!code highlight]
 age: {
 // [!code highlight]
 gt: 40, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 include: {
 author: true,
 },
});
```

### Filter on absence of "-to-many" records

For example, the following query uses `none` to return all users that have zero posts:

```ts
const usersWithZeroPosts = await prisma.user.findMany({
 where: {
 posts: {
 // [!code highlight]
 none: {}, // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true,
 },
});
```

### Filter on absence of "-to-one" relations

The following query returns all posts that don't have an author relation:

```js
const postsWithNoAuthor = await prisma.post.findMany({
 where: {
 author: null, // or author: { } // [!code highlight]
 },
 include: {
 author: true,
 },
});
```

### Filter on presence of related records

The following query returns all users with at least one post:

```ts
const usersWithSomePosts = await prisma.user.findMany({
 where: {
 posts: {
 // [!code highlight]
 some: {}, // [!code highlight]
 }, // [!code highlight]
 },
 include: {
 posts: true,
 },
});
```

## Fluent API

The fluent API lets you _fluently_ traverse the [relations](/orm/prisma-schema/data-model/relations) of your models via function calls. Note that the _last_ function call determines the return type of the entire query (the respective type annotations are added in the code snippets below to make that explicit).

This query returns all `Post` records by a specific `User`:

```ts
const postsByUser: Post[] = await prisma.user
 .findUnique({ where: { email: "alice@prisma.io" } })
 .posts();
```

This is equivalent to the following `findMany` query:

```ts
const postsByUser = await prisma.post.findMany({
 where: {
 author: {
 email: "alice@prisma.io",
 },
 },
});
```

The main difference between the queries is that the fluent API call is translated into two separate database queries while the other one only generates a single query (see this [GitHub issue](https://github.com/prisma/prisma/issues/1984))

This request returns all categories by a specific post:

```ts
const categoriesOfPost: Category[] = await prisma.post
 .findUnique({ where: { id: 1 } })
 .categories();
```

Note that you can chain as many queries as you like. In this example, the chaining starts at `Profile` and goes over `User` to `Post`:

```ts
const posts: Post[] = await prisma.profile
 .findUnique({ where: { id: 1 } })
 .user()
 .posts();
```

The only requirement for chaining is that the previous function call must return only a _single object_ (e.g. as returned by a `findUnique` query or a "to-one relation" like `profile.user()`).

The following query is **not possible** because `findMany` does not return a single object but a _list_:

```ts
// This query is illegal
const posts = await prisma.user.findMany().posts();
```

---
title: Select fields
description: Learn how to return only the fields and relations you need with select and include in Prisma Client.
url: /orm/prisma-client/queries/select-fields
metaTitle: Select fields
metaDescription: Learn how to use select and include in Prisma Client to return only the fields and relations you need.
---

By default, Prisma Client returns all scalar fields for a model and no relations. Use `select` and `include` to make the result smaller, clearer, and more intentional.

## Return the default fields

If you do not pass `select`, `include`, or `omit`, Prisma Client returns all scalar fields for the model and excludes relations from the result.

## Select specific fields

Use `select` when you only need a few scalar fields:

```ts
const user = await prisma.user.findFirst({
 select: {
 email: true,
 name: true,
 },
});
```

## Return nested objects by selecting relation fields

Use `include` when you want related records alongside the main result:

```ts
const user = await prisma.user.findFirst({
 include: {
 posts: true,
 },
});
```

## Nest selections

You can combine both patterns to keep relation payloads focused:

```ts
const user = await prisma.user.findFirst({
 select: {
 email: true,
 posts: {
 select: {
 title: true,
 published: true,
 },
 },
 },
});
```

## Omit fields instead of selecting everything manually

If you mostly want the default result but need to exclude a few fields, see [Excluding fields](/orm/prisma-client/queries/excluding-fields).

## Related pages

- [Relation queries](/orm/prisma-client/queries/relation-queries)
- [Filtering and sorting](/orm/prisma-client/queries/filtering-and-sorting)
- [Prisma Client API reference](/orm/reference/prisma-client-reference#select)

---
title: Transactions and batch queries
description: This page explains the transactions API of Prisma Client
url: /orm/prisma-client/queries/transactions
metaTitle: Transactions and batch queries (Reference)
metaDescription: This page explains the transactions API of Prisma Client.
---

A database transaction is a sequence of read/write operations guaranteed to succeed or fail as a whole (ACID properties: Atomic, Consistent, Isolated, Durable).

Prisma Client supports transactions in several ways:

| Scenario | Technique |
| :------------------ | :--------------------------------------- |
| Dependent writes | Nested writes |
| Independent writes | `$transaction([])` API, Batch operations |
| Read, modify, write | Interactive transactions |

## Nested writes

A [nested write](/orm/prisma-client/queries/relation-queries#nested-writes) performs multiple operations on related records in a single transaction:

```ts
// Create user with posts in a single transaction
const user = await prisma.user.create({
 data: {
 email: "alice@prisma.io",
 posts: {
 create: [{ title: "Post 1" }, { title: "Post 2" }],
 },
 },
});
```

## Batch operations

These bulk operations run as transactions:

- `createMany()` / `createManyAndReturn()`
- `updateMany()` / `updateManyAndReturn()`
- `deleteMany()`

## The `$transaction` API

### Sequential operations

Pass an array of queries to execute sequentially in a transaction:

```ts
const [posts, totalPosts] = await prisma.$transaction([
 prisma.post.findMany({ where: { title: { contains: "prisma" } } }),
 prisma.post.count(),
]);
```

With options:

```ts
await prisma.$transaction(
 [prisma.resource.deleteMany({ where: { name: "name" } }), prisma.resource.createMany({ data })],
 { isolationLevel: Prisma.TransactionIsolationLevel.Serializable },
);
```

### Interactive transactions

For complex logic between queries, use interactive transactions:

```ts
const result = await prisma.$transaction(async (tx) => {
 const sender = await tx.account.update({
 data: { balance: { decrement: 100 } },
 where: { email: "alice@prisma.io" },
 });

 if (sender.balance < 0) {
 throw new Error("Insufficient funds");
 }

 return await tx.account.update({
 data: { balance: { increment: 100 } },
 where: { email: "bob@prisma.io" },
 });
});
```

:::warning
Keep transactions short. Long-running transactions hurt performance and can cause deadlocks.
:::

**Options:**

```ts
await prisma.$transaction(
 async (tx) => {
 /* ... */
 },
 {
 maxWait: 5000, // Max wait to acquire transaction (default: 2000ms)
 timeout: 10000, // Max transaction run time (default: 5000ms)
 isolationLevel: Prisma.TransactionIsolationLevel.Serializable,
 },
);
```

### Transaction isolation level

:::info

This feature is not available on MongoDB, because MongoDB does not support isolation levels.

:::

You can set the transaction [isolation level](https://www.prisma.io/dataguide/intro/database-glossary#isolation-levels) for transactions.

#### Set the isolation level

To set the transaction isolation level, use the `isolationLevel` option in the second parameter of the API.

For sequential operations:

```ts
await prisma.$transaction(
 [
 // Prisma Client operations running in a transaction...
 ],
 {
 isolationLevel: Prisma.TransactionIsolationLevel.Serializable, // optional, default defined by database configuration
 },
);
```

For an interactive transaction:

```jsx
await prisma.$transaction(
 async (prisma) => {
 // Code running in a transaction...
 },
 {
 isolationLevel: Prisma.TransactionIsolationLevel.Serializable, // optional, default defined by database configuration
 maxWait: 5000, // default: 2000
 timeout: 10000, // default: 5000
 },
);
```

#### Supported isolation levels

Prisma Client supports the following isolation levels if they are available in the underlying database:

- `ReadUncommitted`
- `ReadCommitted`
- `RepeatableRead`
- `Snapshot`
- `Serializable`

The isolation levels available for each database connector are as follows:

| Database | `ReadUncommitted` | `ReadCommitted` | `RepeatableRead` | `Snapshot` | `Serializable` |
| ----------- | ----------------- | --------------- | ---------------- | ---------- | -------------- |
| PostgreSQL | ✔️ | ✔️ | ✔️ | No | ✔️ |
| MySQL | ✔️ | ✔️ | ✔️ | No | ✔️ |
| SQL Server | ✔️ | ✔️ | ✔️ | ✔️ | ✔️ |
| CockroachDB | No | No | No | No | ✔️ |
| SQLite | No | No | No | No | ✔️ |

By default, Prisma Client sets the isolation level to the value currently configured in your database.

The isolation levels configured by default in each database are as follows:

| Database | Default |
| ----------- | ---------------- |
| PostgreSQL | `ReadCommitted` |
| MySQL | `RepeatableRead` |
| SQL Server | `ReadCommitted` |
| CockroachDB | `Serializable` |
| SQLite | `Serializable` |

#### Database-specific information on isolation levels

See the following resources:

- [Transaction isolation levels in PostgreSQL](https://www.postgresql.org/docs/9.3/runtime-config-client.html#GUC-DEFAULT-TRANSACTION-ISOLATION)
- [Transaction isolation levels in Microsoft SQL Server](https://learn.microsoft.com/en-us/sql/t-sql/statements/set-transaction-isolation-level-transact-sql?view=sql-server-ver15)
- [Transaction isolation levels in MySQL](https://dev.mysql.com/doc/refman/8.0/en/innodb-transaction-isolation-levels.html)

CockroachDB and SQLite only support the `Serializable` isolation level.

### Transaction timing issues

:::info

- The solution in this section does not apply to MongoDB, because MongoDB does not support [isolation levels](https://www.prisma.io/dataguide/intro/database-glossary#isolation-levels).
- The timing issues discussed in this section do not apply to CockroachDB and SQLite, because these databases only support the highest `Serializable` isolation level.

:::

When two or more transactions run concurrently in certain [isolation levels](https://www.prisma.io/dataguide/intro/database-glossary#isolation-levels), timing issues can cause write conflicts or deadlocks, such as the violation of unique constraints. For example, consider the following sequence of events where Transaction A and Transaction B both attempt to execute a `deleteMany` and a `createMany` operation:

1. Transaction B: `createMany` operation creates a new set of rows.
1. Transaction B: The application commits transaction B.
1. Transaction A: `createMany` operation.
1. Transaction A: The application commits transaction A. The new rows conflict with the rows that transaction B added at step 2.

This conflict can occur at the isolation level `ReadCommitted`, which is the default isolation level in PostgreSQL and Microsoft SQL Server. To avoid this problem, you can set a higher isolation level (`RepeatableRead` or `Serializable`). You can set the isolation level on a transaction. This overrides your database isolation level for that transaction.

To avoid transaction write conflicts and deadlocks on a transaction:

1. On your transaction, use the `isolationLevel` parameter to `Prisma.TransactionIsolationLevel.Serializable`.

 This ensures that your application commits multiple concurrent or parallel transactions as if they were run serially. When a transaction fails due to a write conflict or deadlock, Prisma Client returns a [P2034 error](/orm/reference/error-reference#p2034).

2. In your application code, add a retry around your transaction to handle any P2034 errors, as shown in this example:

 ```ts
 import { Prisma, PrismaClient } from "../prisma/generated/client";

 const prisma = new PrismaClient();
 async function main() {
 const MAX_RETRIES = 5;
 let retries = 0;

 let result;
 while (retries < MAX_RETRIES) {
 try {
 result = await prisma.$transaction(
 [
 prisma.user.deleteMany({
 where: {
 /** args */
 },
 }),
 prisma.post.createMany({
 data: {
 /** args */
 },
 }),
 ],
 {
 isolationLevel: Prisma.TransactionIsolationLevel.Serializable,
 },
 );
 break;
 } catch (error) {
 if (error.code === "P2034") {
 retries++;
 continue;
 }
 throw error;
 }
 }
 }
 ```

### Using `$transaction` within `Promise.all()`

If you wrap a `$transaction` inside a call to `Promise.all()`, the queries inside the transaction will be executed _serially_ (i.e. one after another):

```ts
await prisma.$transaction(async (prisma) => {
 await Promise.all([
 prisma.user.findMany(),
 prisma.user.findMany(),
 prisma.user.findMany(),
 prisma.user.findMany(),
 prisma.user.findMany(),
 prisma.user.findMany(),
 prisma.user.findMany(),
 prisma.user.findMany(),
 prisma.user.findMany(),
 prisma.user.findMany(),
 ]);
});
```

This may be counterintuitive because `Promise.all()` usually _parallelizes_ the calls passed into it.

The reason for this behaviour is that:

- One transaction means that all queries inside it have to be run on the same connection.
- A database connection can only ever execute one query at a time.
- As one query blocks the connection while it is doing its work, putting a transaction into `Promise.all` effectively means that queries should be ran one after another.

## Dependent writes

Writes are **dependent** when operations rely on the result of a preceding operation (e.g., using a database-generated ID).

### Nested writes for dependent operations

Use nested writes when you need to create related records atomically:

```ts
const team = await prisma.team.create({
 data: {
 name: "Aurora Adventures",
 members: {
 create: { email: "alice@prisma.io" },
 },
 },
});
```

If any operation fails, Prisma Client rolls back the entire transaction.

:::note
The `$transaction([])` API cannot pass IDs between operations - use nested writes when you need the generated ID from one record to create another.
:::

## Independent writes

Writes are **independent** if they don't rely on the result of a previous operation. Use these for:

- Updating the status of multiple orders to "Dispatched"
- Marking a list of emails as "Read"

### Bulk operations

```ts
const updateUsers = await prisma.user.updateMany({
 where: { email: { contains: "prisma.io" } },
 data: { role: "ADMIN" },
});
```

### Using `$transaction([])` for independent writes

```ts
const [deleteResult, createResult] = await prisma.$transaction([
 prisma.post.deleteMany({ where: { authorId: 7 } }),
 prisma.user.delete({ where: { id: 7 } }),
]);
```

### Scenario: Pre-computed IDs and the `$transaction([])` API

If you pre-compute IDs (e.g., using UUIDs), you can use either nested writes or `$transaction([])` since both operations know the ID upfront.

#### When to use bulk operations

Consider bulk operations as a solution if:

- ✔ You want to update a batch of the _same type_ of record, like a batch of emails

#### Scenario: Marking emails as read

You are building a service like gmail.com, and your customer wants a **"Mark as read"** feature that allows users to mark all emails as read. Each update to the status of an email is an independent write because the emails do not depend on one another - for example, the "Happy Birthday! 🍰" email from your aunt is unrelated to the promotional email from IKEA.

In the following schema, a `User` can have many received emails (a one-to-many relationship):

```ts
model User {
 id Int @id @default(autoincrement())
 email String @unique
 receivedEmails Email[] // Many emails
}

model Email {
 id Int @id @default(autoincrement())
 user User @relation(fields: [userId], references: [id])
 userId Int
 subject String
 body String
 unread Boolean
}
```

Based on this schema, you can use `updateMany` to mark all unread emails as read:

```ts
await prisma.email.updateMany({
 where: {
 user: {
 id: 10,
 },
 unread: true,
 },
 data: {
 unread: false,
 },
});
```

#### Can I use nested writes with bulk operations?

No - neither `updateMany` nor `deleteMany` currently supports nested writes. For example, you cannot delete multiple teams and all of their members (a cascading delete):

```ts highlight=8;delete
await prisma.team.deleteMany({
 where: {
 id: {
 in: [2, 99, 2, 11],
 },
 },
 data: {
 members: {}, // Cannot access members here // [!code --]
 },
});
```

#### Can I use bulk operations with the `$transaction([])` API?

Yes — for example, you can include multiple `deleteMany` operations inside a `$transaction([])`.

### `$transaction([])` API

The `$transaction([])` API is generic solution to independent writes that allows you to run multiple operations as a single, atomic operation - if any operation fails, Prisma Client rolls back the entire transaction.

Its also worth noting that operations are executed according to the order they are placed in the transaction.

```ts
await prisma.$transaction([iRunFirst, iRunSecond, iRunThird]);
```

> **Note**: Using a query in a transaction does not influence the order of operations in the query itself.

As Prisma Client evolves, use cases for the `$transaction([])` API will increasingly be replaced by more specialized bulk operations (such as `createMany`) and nested writes.

#### When to use the `$transaction([])` API

Consider the `$transaction([])` API if:

- ✔ You want to update a batch that includes different types of records, such as emails and users. The records do not need to be related in any way.
- ✔ You want to batch raw SQL queries (`$executeRaw`) - for example, for features that Prisma Client does not yet support.

#### Scenario: Privacy legislation

GDPR and other privacy legislation give users the right to request that an organization deletes all of their personal data. In the following example schema, a `User` can have many posts and private messages:

```prisma
model User {
 id Int @id @default(autoincrement())
 posts Post[]
 privateMessages PrivateMessage[]
}

model Post {
 id Int @id @default(autoincrement())
 user User @relation(fields: [userId], references: [id])
 userId Int
 title String
 content String
}

model PrivateMessage {
 id Int @id @default(autoincrement())
 user User @relation(fields: [userId], references: [id])
 userId Int
 message String
}
```

If a user invokes the right to be forgotten, we must delete three records: the user record, private messages, and posts. It is critical that _all_ delete operations succeed together or not at all, which makes this a use case for a transaction. However, using a single bulk operation like `deleteMany` is not possible in this scenario because we need to delete across three models. Instead, we can use the `$transaction([])` API to run three operations together - two `deleteMany` and one `delete`:

```ts
const id = 9; // User to be deleted

const deletePosts = prisma.post.deleteMany({
 where: {
 userId: id,
 },
});

const deleteMessages = prisma.privateMessage.deleteMany({
 where: {
 userId: id,
 },
});

const deleteUser = prisma.user.delete({
 where: {
 id: id,
 },
});

await prisma.$transaction([deletePosts, deleteMessages, deleteUser]); // Operations succeed or fail together
```

#### Scenario: Pre-computed IDs and the `$transaction([])` API

Dependent writes are not supported by the `$transaction([])` API - if operation A relies on the ID generated by operation B, use [nested writes](#nested-writes). However, if you _pre-computed_ IDs (for example, by generating GUIDs), your writes become independent. Consider the sign-up flow from the nested writes example:

```ts
await prisma.team.create({
 data: {
 name: "Aurora Adventures",
 members: {
 create: {
 email: "alice@prisma.io",
 },
 },
 },
});
```

Instead of auto-generating IDs, change the `id` fields of `Team` and `User` to a `String` (if you do not provide a value, a UUID is generated automatically). This example uses UUIDs:

```prisma highlight=2,9;delete|3,10;add
model Team {
 id Int @id @default(autoincrement()) // [!code --]
 id String @id @default(uuid()) // [!code ++]
 name String
 members User[]
}

model User {
 id Int @id @default(autoincrement()) // [!code --]
 id String @id @default(uuid()) // [!code ++]
 email String @unique
 teams Team[]
}
```

Refactor the sign-up flow example to use the `$transaction([])` API instead of nested writes:

```ts
import { v4 } from "uuid";

const teamID = v4();
const userID = v4();

await prisma.$transaction([
 prisma.user.create({
 data: {
 id: userID,
 email: "alice@prisma.io",
 team: {
 id: teamID,
 },
 },
 }),
 prisma.team.create({
 data: {
 id: teamID,
 name: "Aurora Adventures",
 },
 }),
]);
```

Technically you can still use nested writes with pre-computed APIs if you prefer that syntax:

```ts
import { v4 } from "uuid";

const teamID = v4();
const userID = v4();

await prisma.team.create({
 data: {
 id: teamID,
 name: "Aurora Adventures",
 members: {
 create: {
 id: userID,
 email: "alice@prisma.io",
 team: {
 id: teamID,
 },
 },
 },
 },
});
```

There's no compelling reason to switch to manually generated IDs and the `$transaction([])` API if you are already using auto-generated IDs and nested writes.

## Read, modify, write

In some cases you may need to perform custom logic as part of an atomic operation - also known as the [read-modify-write pattern](https://en.wikipedia.org/wiki/Read%E2%80%93modify%E2%80%93write). The following is an example of the read-modify-write pattern:

- Read a value from the database
- Run some logic to manipulate that value (for example, contacting an external API)
- Write the value back to the database

All operations should **succeed or fail together** without making unwanted changes to the database, but you do not necessarily need to use an actual database transaction. This section of the guide describes two ways to work with Prisma Client and the read-modify-write pattern:

- Designing idempotent APIs
- Optimistic concurrency control

### Idempotent APIs

Idempotency is the ability to run the same logic with the same parameters multiple times with the same result: the **effect on the database** is the same whether you run the logic once or one thousand times. For example:

- **NOT IDEMPOTENT**: Upsert (update-or-insert) a user in the database with email address `"letoya@prisma.io"`. The `User` table **does not** enforce unique email addresses. The effect on the database is different if you run the logic once (one user created) or ten times (ten users created).
- **IDEMPOTENT**: Upsert (update-or-insert) a user in the database with the email address `"letoya@prisma.io"`. The `User` table **does** enforce unique email addresses. The effect on the database is the same if you run the logic once (one user created) or ten times (existing user is updated with the same input).

Idempotency is something you can and should actively design into your application wherever possible.

#### When to design an idempotent API

- ✔ You need to be able to retry the same logic without creating unwanted side-effects in the databases

#### Scenario: Upgrading a Slack team

You are creating an upgrade flow for Slack that allows teams to unlock paid features. Teams can choose between different plans and pay per user, per month. You use Stripe as your payment gateway, and extend your `Team` model to store a `stripeCustomerId`. Subscriptions are managed in Stripe.

```prisma highlight=5;normal
model Team {
 id Int @id @default(autoincrement())
 name String
 User User[]
 stripeCustomerId String? // [!code highlight]
}
```

The upgrade flow looks like this:

1. Count the number of users
2. Create a subscription in Stripe that includes the number of users
3. Associate the team with the Stripe customer ID to unlock paid features

```ts
const teamId = 9;
const planId = "plan_id";

// Count team members
const numTeammates = await prisma.user.count({
 where: {
 teams: {
 some: {
 id: teamId,
 },
 },
 },
});

// Create a customer in Stripe for plan-9454549
const customer = await stripe.customers.create({
 externalId: teamId,
 plan: planId,
 quantity: numTeammates,
});

// Update the team with the customer id to indicate that they are a customer
// and support querying this customer in Stripe from our application code.
await prisma.team.update({
 data: {
 customerId: customer.id,
 },
 where: {
 id: teamId,
 },
});
```

This example has a problem: you can only run the logic _once_. Consider the following scenario:

1. Stripe creates a new customer and subscription, and returns a customer ID
2. Updating the team **fails** - the team is not marked as a customer in the Slack database
3. The customer is charged by Stripe, but paid features are not unlocked in Slack because the team lacks a valid `customerId`
4. Running the same code again either:
 - Results in an error because the team (defined by `externalId`) already exists - Stripe never returns a customer ID
 - If `externalId` is not subject to a unique constraint, Stripe creates yet another subscription (**not idempotent**)

You cannot re-run this code in case of an error and you cannot change to another plan without being charged twice.

The following refactor (highlighted) introduces a mechanism that checks if a subscription already exists, and either creates the description or updates the existing subscription (which will remain unchanged if the input is identical):

```ts highlight=12-27;normal
// Calculate the number of users times the cost per user
const numTeammates = await prisma.user.count({
 where: {
 teams: {
 some: {
 id: teamId,
 },
 },
 },
});

// Find customer in Stripe // [!code highlight]
let customer = await stripe.customers.get({ externalId: teamID }); // [!code highlight]

if (customer) {
 // [!code highlight]
 // If team already exists, update // [!code highlight]
 customer = await stripe.customers.update({
 // [!code highlight]
 externalId: teamId, // [!code highlight]
 plan: "plan_id", // [!code highlight]
 quantity: numTeammates, // [!code highlight]
 });
} else {
 customer = await stripe.customers.create({
 // If team does not exist, create customer
 externalId: teamId,
 plan: "plan_id",
 quantity: numTeammates,
 });
}

// Update the team with the customer id to indicate that they are a customer
// and support querying this customer in Stripe from our application code.
await prisma.team.update({
 data: {
 customerId: customer.id,
 },
 where: {
 id: teamId,
 },
});
```

You can now retry the same logic multiple times with the same input without adverse effect. To further enhance this example, you can introduce a mechanism whereby the subscription is cancelled or temporarily deactivated if the update does not succeed after a set number of attempts.

### Optimistic concurrency control

Optimistic concurrency control (OCC) is a model for handling concurrent operations on a single entity that does not rely on 🔒 locking. Instead, we **optimistically** assume that a record will remain unchanged in between reading and writing, and use a concurrency token (a timestamp or version field) to detect changes to a record.

If a ❌ conflict occurs (someone else has changed the record since you read it), you cancel the transaction. Depending on your scenario, you can then:

- Re-try the transaction (book another cinema seat)
- Throw an error (alert the user that they are about to overwrite changes made by someone else)

This section describes how to build your own optimistic concurrency control. See also: Plans for [application-level optimistic concurrency control on GitHub](https://github.com/prisma/prisma/issues/4988)

#### When to use optimistic concurrency control

- ✔ You anticipate a high number of concurrent requests (multiple people booking cinema seats)
- ✔ You anticipate that conflicts between those concurrent requests will be rare

Avoiding locks in an application with a high number of concurrent requests makes the application more resilient to load and more scalable overall. Although locking is not inherently bad, locking in a high concurrency environment can lead to unintended consequences - even if you are locking individual rows, and only for a short amount of time. For more information, see:

- [Why ROWLOCK Hints Can Make Queries Slower and Blocking Worse in SQL Server](https://kendralittle.com/2016/02/04/why-rowlock-hints-can-make-queries-slower-and-blocking-worse-in-sql-server/)

#### Scenario: Reserving a seat at the cinema

You are creating a booking system for a cinema. Each movie has a set number of seats. The following schema models movies and seats:

```ts
model Seat {
 id Int @id @default(autoincrement())
 userId Int?
 claimedBy User? @relation(fields: [userId], references: [id])
 movieId Int
 movie Movie @relation(fields: [movieId], references: [id])
}

model Movie {
 id Int @id @default(autoincrement())
 name String @unique
 seats Seat[]
}
```

The following sample code finds the first available seat and assigns that seat to a user:

```ts
const movieName = "Hidden Figures";

// Find first available seat
const availableSeat = await prisma.seat.findFirst({
 where: {
 movie: {
 name: movieName,
 },
 claimedBy: null,
 },
});

// Throw an error if no seats are available
if (!availableSeat) {
 throw new Error(`Oh no! ${movieName} is all booked.`);
}

// Claim the seat
await prisma.seat.update({
 data: {
 claimedBy: userId,
 },
 where: {
 id: availableSeat.id,
 },
});
```

However, this code suffers from the "double-booking problem" - it is possible for two people to book the same seats:

1. Seat 3A returned to Sorcha (`findFirst`)
2. Seat 3A returned to Ellen (`findFirst`)
3. Seat 3A claimed by Sorcha (`update`)
4. Seat 3A claimed by Ellen (`update` - overwrites Sorcha's claim)

Even though Sorcha has successfully booked the seat, the system ultimately stores Ellen's claim. To solve this problem with optimistic concurrency control, add a `version` field to the seat:

```prisma highlight=7;normal
model Seat {
 id Int @id @default(autoincrement())
 userId Int?
 claimedBy User? @relation(fields: [userId], references: [id])
 movieId Int
 movie Movie @relation(fields: [movieId], references: [id])
 version Int // [!code highlight]
}
```

Next, adjust the code to check the `version` field before updating:

```ts highlight=19-38;normal
const userEmail = "alice@prisma.io";
const movieName = "Hidden Figures";

// Find the first available seat
// availableSeat.version might be 0
const availableSeat = await client.seat.findFirst({
 where: {
 Movie: {
 name: movieName,
 },
 claimedBy: null,
 },
});

if (!availableSeat) {
 throw new Error(`Oh no! ${movieName} is all booked.`);
}

// Only mark the seat as claimed if the availableSeat.version // [!code highlight]
// matches the version we're updating. Additionally, increment the // [!code highlight]
// version when we perform this update so all other clients trying // [!code highlight]
// to book this same seat will have an outdated version. // [!code highlight]
const seats = await client.seat.updateMany({
 // [!code highlight]
 data: {
 // [!code highlight]
 claimedBy: userEmail, // [!code highlight]
 version: {
 // [!code highlight]
 increment: 1, // [!code highlight]
 }, // [!code highlight]
 }, // [!code highlight]
 where: {
 // [!code highlight]
 id: availableSeat.id, // [!code highlight]
 version: availableSeat.version, // This version field is the key; only claim seat if in-memory version matches database version, indicating that the field has not been updated // [!code highlight]
 }, // [!code highlight]
}); // [!code highlight]

if (seats.count === 0) {
 // [!code highlight]
 throw new Error(`That seat is already booked! Please try again.`); // [!code highlight]
} // [!code highlight]
```

It is now impossible for two people to book the same seat:

1. Seat 3A returned to Sorcha (`version` is 0)
2. Seat 3A returned to Ellen (`version` is 0)
3. Seat 3A claimed by Sorcha (`version` is incremented to 1, booking succeeds)
4. Seat 3A claimed by Ellen (in-memory `version` (0) does not match database `version` (1) - booking does not succeed)

### Interactive transactions

If you have an existing application, it can be a significant undertaking to refactor your application to use optimistic concurrency control. Interactive Transactions offers a useful escape hatch for cases like this.

To create an interactive transaction, pass an async function into [$transaction](#transaction-api).

The first argument passed into this async function is an instance of the Prisma Client. Below, we will call this instance `tx`. Any Prisma Client call invoked on this `tx` instance is encapsulated into the transaction.

In the example below, Alice and Bob each have $100 in their account. If they try to send more money than they have, the transfer is rejected.

The expected outcome would be for Alice to make 1 transfer for $100 and the other transfer would be rejected. This would result in Alice having $0 and Bob having $200.

```ts
import { PrismaClient } from "../prisma/generated/client";
const prisma = new PrismaClient();

async function transfer(from: string, to: string, amount: number) {
 return await prisma.$transaction(async (tx) => {
 // 1. Decrement amount from the sender.
 const sender = await tx.account.update({
 data: {
 balance: {
 decrement: amount,
 },
 },
 where: {
 email: from,
 },
 });

 // 2. Verify that the sender's balance didn't go below zero.
 if (sender.balance < 0) {
 throw new Error(`${from} doesn't have enough to send ${amount}`);
 }

 // 3. Increment the recipient's balance by amount
 const recipient = tx.account.update({
 data: {
 balance: {
 increment: amount,
 },
 },
 where: {
 email: to,
 },
 });

 return recipient;
 });
}

async function main() {
 // This transfer is successful
 await transfer("alice@prisma.io", "bob@prisma.io", 100);
 // This transfer fails because Alice doesn't have enough funds in her account
 await transfer("alice@prisma.io", "bob@prisma.io", 100);
}

main();
```

In the example above, both `update` queries run within a database transaction. When the application reaches the end of the function, the transaction is **committed** to the database.

If the application encounters an error along the way, the async function will throw an exception and automatically **rollback** the transaction.

You can learn more about interactive transactions in this [section](#interactive-transactions).

:::warning

**Use interactive transactions with caution**. Keeping transactions
open for a long time hurts database performance and can even cause deadlocks.
Try to avoid performing network requests and executing slow queries inside your
transaction functions. We recommend you get in and out as quick as possible!

:::

## Conclusion

Prisma Client supports multiple ways of handling transactions, either directly through the API or by supporting your ability to introduce optimistic concurrency control and idempotency into your application. If you feel like you have use cases in your application that are not covered by any of the suggested options, please open a [GitHub issue](https://github.com/prisma/prisma/issues/new/choose) to start a discussion.

---
title: Custom model and field names
description: Learn how you can decouple the naming of Prisma models from database tables to improve the ergonomics of the generated Prisma Client API
url: /orm/prisma-client/setup-and-configuration/custom-model-and-field-names
metaTitle: Custom model and field names
metaDescription: Learn how you can decouple the naming of Prisma models from database tables to improve the ergonomics of the generated Prisma Client API.
---

The Prisma Client API is generated based on the models in your [Prisma schema](/orm/prisma-schema/overview). Models are _typically_ 1:1 mappings of your database tables.

In some cases, especially when using [introspection](/orm/prisma-schema/introspection), it might be useful to _decouple_ the naming of database tables and columns from the names that are used in your Prisma Client API. This can be done via the [`@map` and `@@map`](/orm/prisma-schema/data-model/models#mapping-model-names-to-tables-or-collections) attributes in your Prisma schema.

You can use `@map` and `@@map` to rename MongoDB fields and collections respectively. This page uses a relational database example.

## Example: Relational database

Assume you have a PostgreSQL relational database schema looking similar to this:

```sql
CREATE TABLE users (
	user_id SERIAL PRIMARY KEY NOT NULL,
	name VARCHAR(256),
	email VARCHAR(256) UNIQUE NOT NULL
);
CREATE TABLE posts (
	post_id SERIAL PRIMARY KEY NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	title VARCHAR(256) NOT NULL,
	content TEXT,
	author_id INTEGER REFERENCES users(user_id)
);
CREATE TABLE profiles (
	profile_id SERIAL PRIMARY KEY NOT NULL,
	bio TEXT,
	user_id INTEGER NOT NULL UNIQUE REFERENCES users(user_id)
);
CREATE TABLE categories (
	category_id SERIAL PRIMARY KEY NOT NULL,
	name VARCHAR(256)
);
CREATE TABLE post_in_categories (
	post_id INTEGER NOT NULL REFERENCES posts(post_id),
	category_id INTEGER NOT NULL REFERENCES categories(category_id)
);
CREATE UNIQUE INDEX post_id_category_id_unique ON post_in_categories(post_id int4_ops,category_id int4_ops);
```

When introspecting a database with that schema, you'll get a Prisma schema looking similar to this:

```prisma
model categories {
 category_id Int @id @default(autoincrement())
 name String? @db.VarChar(256)
 post_in_categories post_in_categories[]
}

model post_in_categories {
 post_id Int
 category_id Int
 categories categories @relation(fields: [category_id], references: [category_id], onDelete: NoAction, onUpdate: NoAction)
 posts posts @relation(fields: [post_id], references: [post_id], onDelete: NoAction, onUpdate: NoAction)

 @@unique([post_id, category_id], map: "post_id_category_id_unique")
}

model posts {
 post_id Int @id @default(autoincrement())
 created_at DateTime? @default(now()) @db.Timestamptz(6)
 title String @db.VarChar(256)
 content String?
 author_id Int?
 users users? @relation(fields: [author_id], references: [user_id], onDelete: NoAction, onUpdate: NoAction)
 post_in_categories post_in_categories[]
}

model profiles {
 profile_id Int @id @default(autoincrement())
 bio String?
 user_id Int @unique
 users users @relation(fields: [user_id], references: [user_id], onDelete: NoAction, onUpdate: NoAction)
}

model users {
 user_id Int @id @default(autoincrement())
 name String? @db.VarChar(256)
 email String @unique @db.VarChar(256)
 posts posts[]
 profiles profiles?
}
```

There are a few "issues" with this Prisma schema when the Prisma Client API is generated:

**Adhering to Prisma ORM's naming conventions**

Prisma ORM has a [naming convention](/orm/reference/prisma-schema-reference#naming-conventions) of **camelCasing** and using the **singular form** for Prisma models. If these naming conventions are not met, the Prisma schema can become harder to interpret and the generated Prisma Client API will feel less natural. Consider the following, generated model:

```prisma
model users {
 user_id Int @id @default(autoincrement())
 name String? @db.VarChar(256)
 email String @unique @db.VarChar(256)
 posts posts[]
 profiles profiles?
}
```

Although `profiles` refers to a 1:1 relation, its type is currently called `profiles` in plural, suggesting that there might be many `profiles` in this relation. With Prisma ORM conventions, the models and fields were _ideally_ named as follows:

```prisma
model User {
 user_id Int @id @default(autoincrement())
 name String? @db.VarChar(256)
 email String @unique @db.VarChar(256)
 posts Post[]
 profile Profile?
}
```

Because these fields are "Prisma ORM-level" [relation fields](/orm/prisma-schema/data-model/relations#relation-fields) that do not manifest you can manually rename them in your Prisma schema.

**Naming of annotated relation fields**

Foreign keys are represented as a combination of a [annotated relation fields](/orm/prisma-schema/data-model/relations#relation-fields) and its corresponding relation scalar field in the Prisma schema. Here's how all the relations from the SQL schema are currently represented:

```prisma
model categories {
 category_id Int @id @default(autoincrement())
 name String? @db.VarChar(256)
 post_in_categories post_in_categories[] // virtual relation field
}

model post_in_categories {
 post_id Int // relation scalar field
 category_id Int // relation scalar field
 categories categories @relation(fields: [category_id], references: [category_id], onDelete: NoAction, onUpdate: NoAction) // virtual relation field
 posts posts @relation(fields: [post_id], references: [post_id], onDelete: NoAction, onUpdate: NoAction)

 @@unique([post_id, category_id], map: "post_id_category_id_unique")
}

model posts {
 post_id Int @id @default(autoincrement())
 created_at DateTime? @default(now()) @db.Timestamptz(6)
 title String @db.VarChar(256)
 content String?
 author_id Int?
 users users? @relation(fields: [author_id], references: [user_id], onDelete: NoAction, onUpdate: NoAction)
 post_in_categories post_in_categories[]
}

model profiles {
 profile_id Int @id @default(autoincrement())
 bio String?
 user_id Int @unique
 users users @relation(fields: [user_id], references: [user_id], onDelete: NoAction, onUpdate: NoAction)
}

model users {
 user_id Int @id @default(autoincrement())
 name String? @db.VarChar(256)
 email String @unique @db.VarChar(256)
 posts posts[]
 profiles profiles?
}
```

## Using @map and @@map to rename fields and models in the Prisma Client API

You can "rename" fields and models that are used in Prisma Client by mapping them to the "original" names in the database using the `@map` and `@@map` attributes. For the example above, you could e.g. annotate your models as follows.

_After_ you introspected your database with `prisma db pull`, you can manually adjust the resulting Prisma schema as follows:

```prisma
model Category {
 id Int @id @default(autoincrement()) @map("category_id")
 name String? @db.VarChar(256)
 post_in_categories PostInCategories[]

 @@map("categories")
}

model PostInCategories {
 post_id Int
 category_id Int
 categories Category @relation(fields: [category_id], references: [id], onDelete: NoAction, onUpdate: NoAction)
 posts Post @relation(fields: [post_id], references: [id], onDelete: NoAction, onUpdate: NoAction)

 @@unique([post_id, category_id], map: "post_id_category_id_unique")
 @@map("post_in_categories")
}

model Post {
 id Int @id @default(autoincrement()) @map("post_id")
 created_at DateTime? @default(now()) @db.Timestamptz(6)
 title String @db.VarChar(256)
 content String?
 author_id Int?
 users User? @relation(fields: [author_id], references: [id], onDelete: NoAction, onUpdate: NoAction)
 post_in_categories PostInCategories[]

 @@map("posts")
}

model Profile {
 id Int @id @default(autoincrement()) @map("profile_id")
 bio String?
 user_id Int @unique
 users User @relation(fields: [user_id], references: [id], onDelete: NoAction, onUpdate: NoAction)

 @@map("profiles")
}

model User {
 id Int @id @default(autoincrement()) @map("user_id")
 name String? @db.VarChar(256)
 email String @unique @db.VarChar(256)
 posts Post[]
 profiles Profile?

 @@map("users")
}
```

With these changes, you're now adhering to Prisma ORM's naming conventions and the generated Prisma Client API feels more "natural":

```ts
// Nested writes
const profile = await prisma.profile.create({
 data: {
 bio: "Hello World",
 users: {
 create: {
 name: "Alice",
 email: "alice@prisma.io",
 },
 },
 },
});

// Fluent API
const userByProfile = await prisma.profile
 .findUnique({
 where: { id: 1 },
 })
 .users();
```

:::info

`prisma db pull` preserves the custom names you defined via `@map` and `@@map` in your Prisma schema on re-introspecting your database.

:::

## Renaming relation fields

Prisma ORM-level [relation fields](/orm/prisma-schema/data-model/relations#relation-fields) (sometimes referred to as "virtual relation fields") only exist in the Prisma schema, but do not actually manifest in the underlying database. You can therefore name these fields whatever you want.

Consider the following example of an ambiguous relation in a SQL database:

```sql
CREATE TABLE "User" (
 id SERIAL PRIMARY KEY
);
CREATE TABLE "Post" (
 id SERIAL PRIMARY KEY,
 "author" integer NOT NULL,
 "favoritedBy" INTEGER,
 FOREIGN KEY ("author") REFERENCES "User"(id),
 FOREIGN KEY ("favoritedBy") REFERENCES "User"(id)
);
```

Prisma ORM's introspection will output the following Prisma schema:

```prisma
model Post {
 id Int @id @default(autoincrement())
 author Int
 favoritedBy Int?
 User_Post_authorToUser User @relation("Post_authorToUser", fields: [author], references: [id], onDelete: NoAction, onUpdate: NoAction)
 User_Post_favoritedByToUser User? @relation("Post_favoritedByToUser", fields: [favoritedBy], references: [id], onDelete: NoAction, onUpdate: NoAction)
}

model User {
 id Int @id @default(autoincrement())
 Post_Post_authorToUser Post[] @relation("Post_authorToUser")
 Post_Post_favoritedByToUser Post[] @relation("Post_favoritedByToUser")
}
```

Because the names of the virtual relation fields `Post_Post_authorToUser` and `Post_Post_favoritedByToUser` are based on the generated relation names, they don't look very friendly in the Prisma Client API. In that case, you can rename the relation fields. For example:

```prisma highlight=11-12;edit
model Post {
 id Int @id @default(autoincrement())
 author Int
 favoritedBy Int?
 User_Post_authorToUser User @relation("Post_authorToUser", fields: [author], references: [id], onDelete: NoAction, onUpdate: NoAction)
 User_Post_favoritedByToUser User? @relation("Post_favoritedByToUser", fields: [favoritedBy], references: [id], onDelete: NoAction, onUpdate: NoAction)
}

model User {
 id Int @id @default(autoincrement())
 //edit-start
 writtenPosts Post[] @relation("Post_authorToUser")
 favoritedPosts Post[] @relation("Post_favoritedByToUser")
 //edit-end
}
```

:::info

`prisma db pull` preserves custom relation fields defined in your Prisma schema on re-introspecting your database.

:::

---
title: Database polyfills
description: Prisma Client provides features that are not achievable with relational databases. These features are referred to as "polyfills" and explained on this page.
url: /orm/prisma-client/setup-and-configuration/database-polyfills
metaTitle: Database polyfills (Concepts)
metaDescription: Prisma Client provides features that are not achievable with relational databases. These features are referred to as "polyfills" and explained on this page.
---

Prisma Client provides features that are typically either not achievable with particular databases or require extensions. These features are referred to as _polyfills_. For all databases, this includes:

- Initializing [ID](/orm/prisma-schema/data-model/models#defining-an-id-field) values with `cuid` and `uuid` values
- Using [`@updatedAt`](/orm/prisma-schema/data-model/models#defining-attributes) to store the time when a record was last updated

For relational databases, this includes:

- [Implicit many-to-many relations](/orm/prisma-schema/data-model/relations/many-to-many-relations#implicit-many-to-many-relations)

For MongoDB, this includes:

- [Relations in general](/orm/prisma-schema/data-model/relations) - foreign key relations between documents are not enforced in MongoDB

---
title: Configuring error formatting
description: This page explains how to configure the formatting of errors when using Prisma Client
url: /orm/prisma-client/setup-and-configuration/error-formatting
metaTitle: Configuring error formatting (Concepts)
metaDescription: This page explains how to configure the formatting of errors when using Prisma Client.
---

By default, Prisma Client uses [ANSI escape characters](https://en.wikipedia.org/wiki/ANSI_escape_code) to pretty print the error stack and give recommendations on how to fix a problem. While this is very useful when using Prisma Client from the terminal, in contexts like a GraphQL API, you only want the minimal error without any additional formatting.

This page explains how error formatting can be configured with Prisma Client.

## Formatting levels

There are 3 error formatting levels:

1. **Pretty Error** (default): Includes a full stack trace with colors, syntax highlighting of the code and extended error message with a possible solution for the problem.
2. **Colorless Error**: Same as pretty errors, just without colors.
3. **Minimal Error**: The raw error message.

In order to configure these different error formatting levels, there are two options:

- Setting the config options via environment variables
- Providing the config options to the `PrismaClient` constructor

## Formatting via environment variables

- [`NO_COLOR`](/orm/reference/environment-variables-reference#no_color): If this env var is provided, colors are stripped from the error messages. Therefore you end up with a **colorless error**. The `NO_COLOR` environment variable is a standard described [here](https://no-color.org/).
- `NODE_ENV=production`: If the env var `NODE_ENV` is set to `production`, only the **minimal error** will be printed. This allows for easier digestion of logs in production environments.

### Formatting via the `PrismaClient` constructor

Alternatively, use the `PrismaClient` [`errorFormat`](/orm/reference/prisma-client-reference#errorformat) parameter to set the error format:

```ts
const prisma = new PrismaClient({
 errorFormat: "pretty",
});
```

---
title: Generating Prisma Client
description: Learn when and how to run prisma generate, configure the generator output, and import the generated client in your app.
url: /orm/prisma-client/setup-and-configuration/generating-prisma-client
metaTitle: Generating Prisma Client
metaDescription: Learn how prisma generate works, how to configure the generator output, and how to import the generated Prisma Client.
---

`prisma generate` creates Prisma Client from the models and generator configuration in your `schema.prisma` file.

## Define a generator

In Prisma ORM v7, the `output` field is required:

```prisma title="schema.prisma"
generator client {
 provider = "prisma-client"
 output = "./generated"
}
```

## Generate the client

Run the following command whenever you add models, change fields, or update generator settings:

```npm
npx prisma generate
```

If you want the CLI-specific options such as `--watch` or `--generator`, see the [`prisma generate` command reference](/cli/generate).

## Import the generated client

Import Prisma Client from the output path you configured:

```ts
import { PrismaClient } from "./generated/client";

const prisma = new PrismaClient();
```

## When to run generate

You should run `prisma generate` after:

- changing your Prisma schema
- updating generator configuration
- enabling features that affect the client API
- pulling schema changes from another branch or teammate

In many projects it also makes sense to run `prisma generate` in `postinstall` or before your production build so deployments always use a current client.

## Related pages

- [Introduction to Prisma Client](/orm/prisma-client/setup-and-configuration/introduction)
- [Generators in the Prisma schema](/orm/prisma-schema/overview/generators)
- [Prisma CLI reference for generate](/cli/generate)

---
title: Introduction to Prisma Client
description: Learn how to set up and configure Prisma Client in your project
url: /orm/prisma-client/setup-and-configuration/introduction
metaTitle: Introduction to Prisma Client
metaDescription: Learn how to set up Prisma Client.
---

Prisma Client is an auto-generated and type-safe query builder that's _tailored_ to your data. The easiest way to get started with Prisma Client is by following the **[Quickstart](/prisma-orm/quickstart/sqlite)**.

[Quickstart (5 min)](/prisma-orm/quickstart/sqlite)

## Prerequisites

In order to set up Prisma Client, you need a Prisma Config and a [Prisma schema file](/orm/prisma-schema/overview):

```ts title="prisma.config.ts" tab="Prisma Config"
import 'dotenv/config';
import { defineConfig, env } from 'prisma/config';

export default defineConfig({
 schema: './prisma/schema.prisma',
 datasource: {
 url: env('DATABASE_URL'),
 },
});
```

```prisma title="schema.prisma" tab="Prisma Schema"
datasource db {
 provider = "postgresql"
}

generator client {
 provider = "prisma-client"
 output = "../src/generated/prisma"
}

model User {
 id Int @id @default(autoincrement())
 createdAt DateTime @default(now())
 email String @unique
 name String?
}
```

## Installation

[Install the Prisma CLI](/orm/reference/prisma-cli-reference), the Prisma Client library, and the [driver adapter](/orm/core-concepts/supported-databases/database-drivers) for your database:

```npm tab="PostgreSQL"
npm install prisma --save-dev
npm install @prisma/client @prisma/adapter-pg pg
```

```npm tab="MySQL / MariaDB"
npm install prisma --save-dev
npm install @prisma/client @prisma/adapter-mariadb mariadb
```

```npm tab="SQLite"
npm install prisma --save-dev
npm install @prisma/client @prisma/adapter-better-sqlite3 better-sqlite3
```

:::note

Prisma 7 requires a [driver adapter](/orm/core-concepts/supported-databases/database-drivers) to connect to your database. Make sure your `package.json` includes `"type": "module"` for ESM support. See the [upgrade guide](/guides/upgrade-prisma-orm/v7) for details.

:::

## Generate the Client API

Prisma Client is based on the models in Prisma Schema. To provide the correct types, you need generate the client code:

```npm
npx prisma generate
```

This will create a `generated` directory based on where you set the `output` to in the Prisma Schema. Any time your import Prisma Client, it will need to come from this generated client API.

## Importing Prisma Client

With the client generated, import it along with your [driver adapter](/orm/core-concepts/supported-databases/database-drivers) and create a new instance:

```ts tab="PostgreSQL"
import { PrismaClient } from "./path/to/generated/prisma";
import { PrismaPg } from "@prisma/adapter-pg";

const adapter = new PrismaPg({
 connectionString: process.env.DATABASE_URL!,
});

export const prisma = new PrismaClient({ adapter });
```

```ts tab="MySQL / MariaDB"
import { PrismaClient } from "./path/to/generated/prisma";
import { PrismaMariaDb } from "@prisma/adapter-mariadb";

const adapter = new PrismaMariaDb({
 host: "localhost",
 user: "root",
 database: "mydb",
});

export const prisma = new PrismaClient({ adapter });
```

```ts tab="SQLite"
import { PrismaClient } from "./path/to/generated/prisma";
import { PrismaBetterSqlite3 } from "@prisma/adapter-better-sqlite3";

const adapter = new PrismaBetterSqlite3({
 url: "file:./dev.db",
});

export const prisma = new PrismaClient({ adapter });
```

```ts tab="PostgreSQL (Edge)"
import { PrismaClient } from "./path/to/generated/prisma/edge";
import { PrismaPostgresAdapter } from "@prisma/adapter-ppg";

const adapter = new PrismaPostgresAdapter({
 connectionString: process.env.DATABASE_URL!,
});

export const prisma = new PrismaClient({ adapter });
```

:::warning

`PrismaClient` requires a driver adapter in Prisma 7. Calling `new PrismaClient()` without an `adapter` will result in an error.

:::

Find out what [driver adapter](/orm/core-concepts/supported-databases/database-drivers) is needed for your database.

Your application should generally only create **one instance** of `PrismaClient`. How to achieve this depends on whether you are using Prisma ORM in a [long-running application](/orm/prisma-client/setup-and-configuration/databases-connections#prismaclient-in-long-running-applications) or in a [serverless environment](/orm/prisma-client/setup-and-configuration/databases-connections#prismaclient-in-serverless-environments).

Creating multiple instances of `PrismaClient` will create multiple connection pools and can hit the connection limit for your database. Too many connections may start to **slow down your database** and eventually lead to errors such as:

```bash
Error in connector: Error querying the database: db error: FATAL: sorry, too many clients already
 at PrismaClientFetcher.request
```

## Use Prisma Client to send queries to your database

Once you have instantiated `PrismaClient`, you can start sending queries in your code:

```ts
// run inside `async` function
const newUser = await prisma.user.create({
 data: {
 name: "Alice",
 email: "alice@prisma.io",
 },
});

const users = await prisma.user.findMany();
```

## Evolving your application

Whenever you make changes to your database that are reflected in the Prisma schema, you need to manually re-generate Prisma Client to update the generated code in your output directory:

```npm
npx prisma generate
```

---
title: Read replicas
description: Learn how to set up and use read replicas with Prisma Client
url: /orm/prisma-client/setup-and-configuration/read-replicas
metaTitle: Read replicas
metaDescription: Learn how to set up and use read replicas with Prisma Client
---

Read replicas enable you to distribute workloads across database replicas for high-traffic workloads. The [read replicas extension](https://github.com/prisma/extension-read-replicas), `@prisma/extension-read-replicas`, adds support for read-only database replicas to Prisma Client.

If you run into a bug or have feedback, create a GitHub issue [here](https://github.com/prisma/extension-read-replicas/issues/new).

## Setup the read replicas extension

Install the extension:

```npm
npm install @prisma/extension-read-replicas@latest
```

Initialize the extension by extending your Prisma Client instance and provide the extension with full `PrismaClient` instances for your read replicas. The default approach is to use driver adapters:

```ts
import { readReplicas } from "@prisma/extension-read-replicas";
import { PrismaPg } from "@prisma/adapter-pg";
import { PrismaClient } from "./generated/prisma/client";

// Create main client with adapter
const mainAdapter = new PrismaPg({
 connectionString: process.env.DATABASE_URL!,
});

const mainClient = new PrismaClient({ adapter: mainAdapter });

// Create replica client with adapter
const replicaAdapter = new PrismaPg({
 connectionString: process.env.REPLICA_URL!,
});

const replicaClient = new PrismaClient({ adapter: replicaAdapter });

// Extend main client with read replicas
const prisma = mainClient.$extends(readReplicas({ replicas: [replicaClient] }));

// Query is run against the database replica
await prisma.post.findMany();

// Query is run against the primary database
await prisma.post.create({
 data: {
 /** */
 },
});
```

All read operations (e.g. `findMany`) are executed against the database replica. All write operations (e.g. `create`, `update`) and `$transaction` queries are executed against your primary database.

## Configure multiple database replicas

The `replicas` property accepts an array of `PrismaClient` instances for all your database replicas:

```ts
import { readReplicas } from "@prisma/extension-read-replicas";
import { PrismaPg } from "@prisma/adapter-pg";
import { PrismaClient } from "./generated/prisma/client";

// Create main client
const mainAdapter = new PrismaPg({
 connectionString: process.env.DATABASE_URL!,
});
const mainClient = new PrismaClient({ adapter: mainAdapter });

// Create multiple replica clients
const replicaAdapter1 = new PrismaPg({
 connectionString: process.env.DATABASE_URL_REPLICA_1!,
});
const replicaClient1 = new PrismaClient({ adapter: replicaAdapter1 });

const replicaAdapter2 = new PrismaPg({
 connectionString: process.env.DATABASE_URL_REPLICA_2!,
});
const replicaClient2 = new PrismaClient({ adapter: replicaAdapter2 });

// Configure multiple replicas
const prisma = mainClient.$extends(
 readReplicas({
 replicas: [replicaClient1, replicaClient2],
 }),
);
```

If you have more than one read replica configured, a database replica will be randomly selected to execute your query.

## Executing read operations against your primary database

You can use the `$primary()` method to explicitly execute a read operation against your primary database:

```ts
const posts = await prisma.$primary().post.findMany();
```

## Executing operations against a database replica

You can use the `$replica()` method to explicitly execute your query against a replica instead of your primary database:

```ts
const result = await prisma.$replica().user.findFirst(...)
```

---
title: Composite types
description: Work with composite types and embedded documents in MongoDB
url: /orm/prisma-client/special-fields-and-types/composite-types
metaTitle: Composite types
metaDescription: Composite types
---

:::warning

Composite types are only available with MongoDB.

:::

[Composite types](/orm/prisma-schema/data-model/models), known as [embedded documents](https://www.mongodb.com/docs/manual/data-modeling/#embedded-data) in MongoDB, allow you to embed records within other records.

This page explains how to:

- [find](#finding-records-that-contain-composite-types-with-find-and-findmany) records that contain composite types using `findFirst` and `findMany`
- [create](#creating-records-with-composite-types-using-create-and-createmany) new records with composite types using `create` and `createMany`
- [update](#changing-composite-types-within-update-and-updatemany) composite types within existing records using `update` and `updateMany`
- [delete](#deleting-records-that-contain-composite-types-with-delete-and-deletemany) records with composite types using `delete` and `deleteMany`

## Example schema

We’ll use this schema for the examples that follow:

```prisma title="schema.prisma" showLineNumbers
generator client {
 provider = "prisma-client-js"
}

datasource db {
 provider = "mongodb"
 url = env("DATABASE_URL")
}

model Product {
 id String @id @default(auto()) @map("_id") @db.ObjectId
 name String @unique
 price Float
 colors Color[]
 sizes Size[]
 photos Photo[]
 orders Order[]
}

model Order {
 id String @id @default(auto()) @map("_id") @db.ObjectId
 product Product @relation(fields: [productId], references: [id])
 color Color
 size Size
 shippingAddress Address
 billingAddress Address?
 productId String @db.ObjectId
}

enum Color {
 Red
 Green
 Blue
}

enum Size {
 Small
 Medium
 Large
 XLarge
}

type Photo {
 height Int @default(200)
 width Int @default(100)
 url String
}

type Address {
 street String
 city String
 zip String
}
```

In this schema, the `Product` model has a `Photo[]` composite type, and the `Order` model has two composite `Address` types. The `shippingAddress` is required, but the `billingAddress` is optional.

## Considerations when using composite types

There are currently some limitations when using composite types in Prisma Client:

- [`findUnique()`](/orm/reference/prisma-client-reference#findunique) can't filter on composite types
- [`aggregate`](/orm/prisma-client/queries/aggregation-grouping-summarizing#aggregate), [`groupBy()`](/orm/prisma-client/queries/aggregation-grouping-summarizing#group-by), [`count`](/orm/prisma-client/queries/aggregation-grouping-summarizing#count) don’t support composite operations

## Default values for required fields on composite types

When you carry out a database read on a composite type and all of the following conditions are true, Prisma Client inserts the default value into the result:

- A field on the composite type is [required](/orm/prisma-schema/data-model/models#optional-and-mandatory-fields), and
- this field has a [default value](/orm/prisma-schema/data-model/models#defining-a-default-value), and
- this field is not present in the returned document or documents.

Note:

- This is the same behavior as with [model fields](/orm/reference/prisma-schema-reference#model-field-scalar-types).
- On read operations, Prisma Client inserts the default value into the result, but does not insert the default value into the database.

In our example schema, suppose that you add a required field to `photo`. This field, `bitDepth`, has a default value:

```prisma title="schema.prisma" highlight=4;add
...
type Photo {
 ...
 bitDepth Int @default(8) // [!code ++]
}

...
```

Suppose that you then run `npx prisma db push` to [update your database](/orm/reference/prisma-cli-reference#db-push) and regenerate your Prisma Client with `npx prisma generate`. Then, you run the following application code:

```ts
console.dir(await prisma.product.findMany({}), { depth: Infinity });
```

The `bitDepth` field has no content because you have only just added this field, so the query returns the default value of `8`.

## Finding records that contain composite types with `find` and `findMany`

Records can be filtered by a composite type within the `where` operation.

The following section describes the operations available for filtering by a single type or multiple types, and gives examples of each.

### Filtering for one composite type

Use the `is`, `equals`, `isNot` and `isSet` operations to change a single composite type:

- `is`: Filter results by matching composite types. Requires one or more fields to be present _(e.g. Filter orders by the street name on the shipping address)_
- `equals`: Filter results by matching composite types. Requires all fields to be present. _(e.g. Filter orders by the full shipping address)_
- `isNot`: Filter results by non-matching composite types
- `isSet` : Filter optional fields to include only results that have been set (either set to a value, or explicitly set to `null`). Setting this filter to `true` will exclude `undefined` results that are not set at all.

For example, use `is` to filter for orders with a street name of `'555 Candy Cane Lane'`:

```ts
const orders = await prisma.order.findMany({
 where: {
 shippingAddress: {
 is: {
 street: "555 Candy Cane Lane",
 },
 },
 },
});
```

Use `equals` to filter for orders which match on all fields in the shipping address:

```ts
const orders = await prisma.order.findMany({
 where: {
 shippingAddress: {
 equals: {
 street: "555 Candy Cane Lane",
 city: "Wonderland",
 zip: "52337",
 },
 },
 },
});
```

You can also use a shorthand notation for this query, where you leave out the `equals`:

```ts
const orders = await prisma.order.findMany({
 where: {
 shippingAddress: {
 street: "555 Candy Cane Lane",
 city: "Wonderland",
 zip: "52337",
 },
 },
});
```

Use `isNot` to filter for orders that do not have a `zip` code of `'52337'`:

```ts
const orders = await prisma.order.findMany({
 where: {
 shippingAddress: {
 isNot: {
 zip: "52337",
 },
 },
 },
});
```

Use `isSet` to filter for orders where the optional `billingAddress` has been set (either to a value or to `null`):

```ts
const orders = await prisma.order.findMany({
 where: {
 billingAddress: {
 isSet: true,
 },
 },
});
```

### Filtering for many composite types

Use the `equals`, `isEmpty`, `every`, `some` and `none` operations to filter for multiple composite types:

- `equals`: Checks exact equality of the list
- `isEmpty`: Checks if the list is empty
- `every`: Every item in the list must match the condition
- `some`: One or more of the items in the list must match the condition
- `none`: None of the items in the list can match the condition
- `isSet` : Filter optional fields to include only results that have been set (either set to a value, or explicitly set to `null`). Setting this filter to `true` will exclude `undefined` results that are not set at all.

For example, you can use `equals` to find products with a specific list of photos (all `url`, `height` and `width` fields must match):

```ts
const product = prisma.product.findMany({
 where: {
 photos: {
 equals: [
 {
 url: "1.jpg",
 height: 200,
 width: 100,
 },
 {
 url: "2.jpg",
 height: 200,
 width: 100,
 },
 ],
 },
 },
});
```

You can also use a shorthand notation for this query, where you leave out the `equals` and specify just the fields that you want to filter for:

```ts
const product = prisma.product.findMany({
 where: {
 photos: [
 {
 url: "1.jpg",
 height: 200,
 width: 100,
 },
 {
 url: "2.jpg",
 height: 200,
 width: 100,
 },
 ],
 },
});
```

Use `isEmpty` to filter for products with no photos:

```ts
const product = prisma.product.findMany({
 where: {
 photos: {
 isEmpty: true,
 },
 },
});
```

Use `some` to filter for products where one or more photos has a `url` of `"2.jpg"`:

```ts
const product = prisma.product.findFirst({
 where: {
 photos: {
 some: {
 url: "2.jpg",
 },
 },
 },
});
```

Use `none` to filter for products where no photos have a `url` of `"2.jpg"`:

```ts
const product = prisma.product.findFirst({
 where: {
 photos: {
 none: {
 url: "2.jpg",
 },
 },
 },
});
```

## Creating records with composite types using `create` and `createMany`

:::info

When you create a record with a composite type that has a unique constraint, note that MongoDB does not enforce unique values inside a record. [Learn more](#duplicate-values-in-unique-fields-of-composite-types).

:::

Composite types can be created within a `create` or `createMany` method using the `set` operation. For example, you can use `set` within `create` to create an `Address` composite type inside an `Order`:

```ts
const order = await prisma.order.create({
 data: {
 // Normal relation
 product: { connect: { id: "some-object-id" } },
 color: "Red",
 size: "Large",
 // Composite type
 shippingAddress: {
 set: {
 street: "1084 Candycane Lane",
 city: "Silverlake",
 zip: "84323",
 },
 },
 },
});
```

You can also use a shorthand notation where you leave out the `set` and specify just the fields that you want to create:

```ts
const order = await prisma.order.create({
 data: {
 // Normal relation
 product: { connect: { id: "some-object-id" } },
 color: "Red",
 size: "Large",
 // Composite type
 shippingAddress: {
 street: "1084 Candycane Lane",
 city: "Silverlake",
 zip: "84323",
 },
 },
});
```

For an optional type, like the `billingAddress`, you can also set the value to `null`:

```ts
const order = await prisma.order.create({
 data: {
 // Normal relation
 product: { connect: { id: "some-object-id" } },
 color: "Red",
 size: "Large",
 // Composite type
 shippingAddress: {
 street: "1084 Candycane Lane",
 city: "Silverlake",
 zip: "84323",
 },
 // Embedded optional type, set to null
 billingAddress: {
 set: null,
 },
 },
});
```

To model the case where an `product` contains a list of multiple `photos`, you can `set` multiple composite types at once:

```ts
const product = await prisma.product.create({
 data: {
 name: "Forest Runners",
 price: 59.99,
 colors: ["Red", "Green"],
 sizes: ["Small", "Medium", "Large"],
 // New composite type
 photos: {
 set: [
 { height: 100, width: 200, url: "1.jpg" },
 { height: 100, width: 200, url: "2.jpg" },
 ],
 },
 },
});
```

You can also use a shorthand notation where you leave out the `set` and specify just the fields that you want to create:

```ts
const product = await prisma.product.create({
 data: {
 name: "Forest Runners",
 price: 59.99,
 // Scalar lists that we already support
 colors: ["Red", "Green"],
 sizes: ["Small", "Medium", "Large"],
 // New composite type
 photos: [
 { height: 100, width: 200, url: "1.jpg" },
 { height: 100, width: 200, url: "2.jpg" },
 ],
 },
});
```

These operations also work within the `createMany` method. For example, you can create multiple `product`s which each contain a list of `photos`:

```ts
const product = await prisma.product.createMany({
 data: [
 {
 name: "Forest Runners",
 price: 59.99,
 colors: ["Red", "Green"],
 sizes: ["Small", "Medium", "Large"],
 photos: [
 { height: 100, width: 200, url: "1.jpg" },
 { height: 100, width: 200, url: "2.jpg" },
 ],
 },
 {
 name: "Alpine Blazers",
 price: 85.99,
 colors: ["Blue", "Red"],
 sizes: ["Large", "XLarge"],
 photos: [
 { height: 100, width: 200, url: "1.jpg" },
 { height: 150, width: 200, url: "4.jpg" },
 { height: 200, width: 200, url: "5.jpg" },
 ],
 },
 ],
});
```

## Changing composite types within `update` and `updateMany`

:::info

When you update a record with a composite type that has a unique constraint, note that MongoDB does not enforce unique values inside a record. [Learn more](#duplicate-values-in-unique-fields-of-composite-types).

:::

Composite types can be set, updated or removed within an `update` or `updateMany` method. The following section describes the operations available for updating a single type or multiple types at once, and gives examples of each.

### Changing a single composite type

Use the `set`, `unset` `update` and `upsert` operations to change a single composite type:

- Use `set` to set a composite type, overriding any existing value
- Use `unset` to unset a composite type. Unlike `set: null`, `unset` removes the field entirely
- Use `update` to update a composite type
- Use `upsert` to `update` an existing composite type if it exists, and otherwise `set` the composite type

For example, use `update` to update a required `shippingAddress` with an `Address` composite type inside an `Order`:

```ts
const order = await prisma.order.update({
 where: {
 id: "some-object-id",
 },
 data: {
 shippingAddress: {
 // Update just the zip field
 update: {
 zip: "41232",
 },
 },
 },
});
```

For an optional embedded type, like the `billingAddress`, use `upsert` to create a new record if it does not exist, and update the record if it does:

```ts
const order = await prisma.order.update({
 where: {
 id: "some-object-id",
 },
 data: {
 billingAddress: {
 // Create the address if it doesn't exist,
 // otherwise update it
 upsert: {
 set: {
 street: "1084 Candycane Lane",
 city: "Silverlake",
 zip: "84323",
 },
 update: {
 zip: "84323",
 },
 },
 },
 },
});
```

You can also use the `unset` operation to remove an optional embedded type. The following example uses `unset` to remove the `billingAddress` from an `Order`:

```ts
const order = await prisma.order.update({
 where: {
 id: "some-object-id",
 },
 data: {
 billingAddress: {
 // Unset the billing address
 // Removes "billingAddress" field from order
 unset: true,
 },
 },
});
```

You can use [filters](/orm/prisma-client/special-fields-and-types/composite-types#finding-records-that-contain-composite-types-with-find-and-findmany) within `updateMany` to update all records that match a composite type. The following example uses the `is` filter to match the street name from a shipping address on a list of orders:

```ts
const orders = await prisma.order.updateMany({
 where: {
 shippingAddress: {
 is: {
 street: "555 Candy Cane Lane",
 },
 },
 },
 data: {
 shippingAddress: {
 update: {
 street: "111 Candy Cane Drive",
 },
 },
 },
});
```

### Changing multiple composite types

Use the `set`, `push`, `updateMany` and `deleteMany` operations to change a list of composite types:

- `set`: Set an embedded list of composite types, overriding any existing list
- `push`: Push values to the end of an embedded list of composite types
- `updateMany`: Update many composite types at once
- `deleteMany`: Delete many composite types at once

For example, use `push` to add a new photo to the `photos` list:

```ts
const product = prisma.product.update({
 where: {
 id: "62de6d328a65d8fffdae2c18",
 },
 data: {
 photos: {
 // Push a photo to the end of the photos list
 push: [{ height: 100, width: 200, url: "1.jpg" }],
 },
 },
});
```

Use `updateMany` to update photos with a `url` of `1.jpg` or `2.png`:

```ts
const product = prisma.product.update({
 where: {
 id: "62de6d328a65d8fffdae2c18",
 },
 data: {
 photos: {
 updateMany: {
 where: {
 url: "1.jpg",
 },
 data: {
 url: "2.png",
 },
 },
 },
 },
});
```

The following example uses `deleteMany` to delete all photos with a `height` of 100:

```ts
const product = prisma.product.update({
 where: {
 id: "62de6d328a65d8fffdae2c18",
 },
 data: {
 photos: {
 deleteMany: {
 where: {
 height: 100,
 },
 },
 },
 },
});
```

## Upserting composite types with `upsert`

:::info

When you create or update the values in a composite type that has a unique constraint, note that MongoDB does not enforce unique values inside a record. [Learn more](#duplicate-values-in-unique-fields-of-composite-types).

:::

To create or update a composite type, use the `upsert` method. You can use the same composite operations as the `create` and `update` methods above.

For example, use `upsert` to either create a new product or add a photo to an existing product:

```ts
const product = await prisma.product.upsert({
 where: {
 name: "Forest Runners",
 },
 create: {
 name: "Forest Runners",
 price: 59.99,
 colors: ["Red", "Green"],
 sizes: ["Small", "Medium", "Large"],
 photos: [
 { height: 100, width: 200, url: "1.jpg" },
 { height: 100, width: 200, url: "2.jpg" },
 ],
 },
 update: {
 photos: {
 push: { height: 300, width: 400, url: "3.jpg" },
 },
 },
});
```

## Deleting records that contain composite types with `delete` and `deleteMany`

To remove records which embed a composite type, use the `delete` or `deleteMany` methods. This will also remove the embedded composite type.

For example, use `deleteMany` to delete all products with a `size` of `"Small"`. This will also delete any embedded `photos`.

```ts
const deleteProduct = await prisma.product.deleteMany({
 where: {
 sizes: {
 equals: "Small",
 },
 },
});
```

You can also use [filters](/orm/prisma-client/special-fields-and-types/composite-types#finding-records-that-contain-composite-types-with-find-and-findmany) to delete records that match a composite type. The example below uses the `some` filter to delete products that contain a certain photo:

```ts
const product = await prisma.product.deleteMany({
 where: {
 photos: {
 some: {
 url: "2.jpg",
 },
 },
 },
});
```

## Ordering composite types

You can use the `orderBy` operation to sort results in ascending or descending order.

For example, the following command finds all orders and orders them by the city name in the shipping address, in ascending order:

```ts
const orders = await prisma.order.findMany({
 orderBy: {
 shippingAddress: {
 city: "asc",
 },
 },
});
```

## Duplicate values in unique fields of composite types

Be careful when you carry out any of the following operations on a record with a composite type that has a unique constraint. In this situation, MongoDB does not enforce unique values inside a record.

- When you create the record
- When you add data to the record
- When you update data in the record

If your schema has a composite type with a `@@unique` constraint, MongoDB prevents you from storing the same value for the constrained value in two or more of the records that contain this composite type. However, MongoDB does does not prevent you from storing multiple copies of the same field value in a single record.

Note that you can [use Prisma ORM relations to work around this issue](#use-prisma-orm-relations-to-enforce-unique-values-in-a-record).

For example, in the following schema, `MailBox` has a composite type, `addresses`, which has a `@@unique` constraint on the `email` field.

```prisma
type Address {
 email String
}

model MailBox {
 name String
 addresses Address[]

 @@unique([addresses.email])
}
```

The following code creates a record with two identical values in `address`. MongoDB does not throw an error in this situation, and it stores `alice@prisma.io` in `addresses` twice.

```ts
await prisma.MailBox.createMany({
 data: [
 {
 name: "Alice",
 addresses: {
 set: [
 {
 address: "alice@prisma.io", // Not unique
 },
 {
 address: "alice@prisma.io", // Not unique
 },
 ],
 },
 },
 ],
});
```

Note: MongoDB throws an error if you try to store the same value in two separate records. In our example above, if you try to store the email address `alice@prisma.io` for the user Alice and for the user Bob, MongoDB does not store the data and throws an error.

### Use Prisma ORM relations to enforce unique values in a record

In the example above, MongoDB did not enforce the unique constraint on a nested address name. However, you can model your data differently to enforce unique values in a record. To do so, use Prisma ORM [relations](/orm/prisma-schema/data-model/relations) to turn the composite type into a collection. Set a relationship to this collection and place a unique constraint on the field that you want to be unique.

In the following example, MongoDB enforces unique values in a record. There is a relation between `Mailbox` and the `Address` model. Also, the `name` field in the `Address` model has a unique constraint.

```prisma
model Address {
 id String @id @default(auto()) @map("_id") @db.ObjectId
 name String
 mailbox Mailbox? @relation(fields: [mailboxId], references: [id])
 mailboxId String? @db.ObjectId

 @@unique([name])
}

model Mailbox {
 id String @id @default(auto()) @map("_id") @db.ObjectId
 name String
 addresses Address[] @relation
}
```

```ts
await prisma.MailBox.create({
 data: {
 name: "Alice",
 addresses: {
 create: [
 { name: "alice@prisma.io" }, // Not unique
 { name: "alice@prisma.io" }, // Not unique
 ],
 },
 },
});
```

If you run the above code, MongoDB enforces the unique constraint. It does not allow your application to add two addresses with the name `alice@prisma.io`.

---
title: Fields & types
description: Learn how to use about special fields and types with Prisma Client
url: /orm/prisma-client/special-fields-and-types
metaTitle: Fields & types
metaDescription: Learn how to use about special fields and types with Prisma Client.
---

This section covers various special fields and types you can use with Prisma Client.

## Working with `Decimal`

`Decimal` fields are represented by the [`Decimal.js` library](https://mikemcl.github.io/decimal.js/). The following example demonstrates how to import and use `Prisma.Decimal`:

```ts
import { PrismaClient, Prisma } from "@prisma/client";

const newTypes = await prisma.sample.create({
 data: {
 cost: new Prisma.Decimal(24.454545),
 },
});
```

<br />

You can also perform arithmetic operations:

```ts
import { PrismaClient, Prisma } from "@prisma/client";

const newTypes = await prisma.sample.create({
 data: {
 cost: new Prisma.Decimal(24.454545).plus(1),
 },
});
```

`Prisma.Decimal` uses Decimal.js, see [Decimal.js docs](https://mikemcl.github.io/decimal.js) to learn more.

:::warning

The use of the `Decimal` field [is not currently supported in MongoDB](https://github.com/prisma/prisma/issues/12637).

:::

## Working with `BigInt`

### Overview

`BigInt` fields are represented by the [`BigInt` type](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/BigInt) (Node.js 10.4.0+ required). The following example demonstrates how to use the `BigInt` type:

```ts
import { PrismaClient, Prisma } from "@prisma/client";

const newTypes = await prisma.sample.create({
 data: {
 revenue: BigInt(534543543534),
 },
});
```

### Serializing `BigInt`

Prisma Client returns records as plain JavaScript objects. If you attempt to use `JSON.stringify` on an object that includes a `BigInt` field, you will see the following error:

```
Do not know how to serialize a BigInt
```

To work around this issue, use a customized implementation of `JSON.stringify`:

```js
JSON.stringify(
 this,
 (key, value) => (typeof value === "bigint" ? value.toString() : value), // return everything else unchanged
);
```

## Working with `Bytes`

`Bytes` fields are represented by the [`Uint8Array`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Uint8Array) type. The following example demonstrates how to use the `Uint8Array` type:

```ts
import { PrismaClient, Prisma } from "@prisma/client";

const newTypes = await prisma.sample.create({
 data: {
 myField: new Uint8Array([1, 2, 3, 4]),
 },
});
```

## Working with `DateTime`

:::note

There currently is a [bug](https://github.com/prisma/prisma/issues/9516) that doesn't allow you to pass in `DateTime` values as strings and produces a runtime error when you do. `DateTime` values need to be passed as [`Date`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date) objects (i.e. `new Date('2024-12-04')` instead of `'2024-12-04'`).

:::

When creating records that have fields of type [`DateTime`](/orm/reference/prisma-schema-reference#datetime), Prisma Client accepts values as [`Date`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date) objects adhering to the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) standard.

Consider the following schema:

```prisma
model User {
 id Int @id @default(autoincrement())
 birthDate DateTime?
}
```

Here are some examples for creating new records:

##### Jan 01, 1998; 00 h 00 min and 000 ms

```ts
await prisma.user.create({
 data: {
 birthDate: new Date("1998"),
 },
});
```

##### Dec 01, 1998; 00 h 00 min and 000 ms

```ts
await prisma.user.create({
 data: {
 birthDate: new Date("1998-12"),
 },
});
```

##### Dec 24, 1998; 00 h 00 min and 000 ms

```ts
await prisma.user.create({
 data: {
 birthDate: new Date("1998-12-24"),
 },
});
```

##### Dec 24, 1998; 06 h 22 min 33s and 444 ms

```ts
await prisma.user.create({
 data: {
 birthDate: new Date("1998-12-24T06:22:33.444Z"),
 },
});
```

---
title: Null and undefined
description: How Prisma Client handles null and undefined
preview: false
url: /orm/prisma-client/special-fields-and-types/null-and-undefined
metaTitle: Null and undefined in Prisma Client (Reference)
metaDescription: How Prisma Client handles null and undefined
---

:::warning

In Prisma ORM, if `undefined` is passed as a value, it is not included in the generated query. This behavior can lead to unexpected results and data loss. We strongly recommend enabling the `strictUndefinedChecks` preview feature described below.

For documentation on the current behavior (without the `strictUndefinedChecks` Preview feature) see [current behavior](#current-behavior).

:::

## Strict undefined checks (Preview feature)

The `strictUndefinedChecks` preview feature changes how Prisma Client handles `undefined` values, offering better protection against accidental data loss or unintended query behavior.

### Enabling strict undefined checks

To enable this feature, add the following to your Prisma schema:

```prisma
generator client {
 provider = "prisma-client"
 output = "./generated"
 previewFeatures = ["strictUndefinedChecks"]
}
```

### Using strict undefined checks

When this feature is enabled:

1. Explicitly setting a field to `undefined` in a query will cause a runtime error.
2. To skip a field in a query, use the new `Prisma.skip` symbol instead of `undefined`.

Example usage:

```typescript
// This will throw an error
prisma.user.create({
 data: {
 name: "Alice",
 email: undefined, // Error: Cannot explicitly use undefined here
 },
});

// Use `Prisma.skip` (a symbol provided by Prisma) to omit a field
prisma.user.create({
 data: {
 name: "Alice",
 email: Prisma.skip, // This field will be omitted from the query
 },
});
```

This change helps prevent accidental deletions or updates, such as:

```typescript
// Before: This would delete all users
prisma.user.deleteMany({
 where: {
 id: undefined
 }
})

// After: This will throw an error
Invalid \`prisma.user.deleteMany()\` invocation in
/client/tests/functional/strictUndefinedChecks/test.ts:0:0
 XX })
 XX
 XX test('throws on undefined input field', async () => {
→ XX const result = prisma.user.deleteMany({
 where: {
 id: undefined
 ~~~~~~~~~
 }
 })
Invalid value for argument \`where\`: explicitly \`undefined\` values are not allowed."
```

### Migration path

To migrate existing code:

```typescript
// Before
let optionalEmail: string | undefined;

prisma.user.create({
 data: {
 name: "Alice",
 email: optionalEmail,
 },
});

// After
prisma.user.create({
 data: {
 name: "Alice",
 email: optionalEmail ?? Prisma.skip, // [!code highlight]
 },
});
```

### `exactOptionalPropertyTypes`

In addition to `strictUndefinedChecks`, we also recommend enabling the TypeScript compiler option `exactOptionalPropertyTypes`. This option enforces that optional properties must match exactly, which can help catch potential issues with `undefined` values in your code. While `strictUndefinedChecks` will raise runtime errors for invalid `undefined` usage, `exactOptionalPropertyTypes` will catch these issues during the build process.

Learn more about `exactOptionalPropertyTypes` in the [TypeScript documentation](https://www.typescriptlang.org/tsconfig/#exactOptionalPropertyTypes).

### Feedback

As always, we welcome your feedback on this feature. Please share your thoughts and suggestions in the [GitHub discussion for this Preview feature](https://github.com/prisma/prisma/discussions/25271).

## current behavior

Prisma Client differentiates between `null` and `undefined`:

- `null` is a **value**
- `undefined` means **do nothing**

:::info

This is particularly important to account for in [a **Prisma ORM with GraphQL context**, where `null` and `undefined` are interchangeable](#null-and-undefined-in-a-graphql-resolver).

:::

The data below represents a `User` table. This set of data will be used in all of the examples below:

| id | name | email |
| --- | ------- | ----------------- |
| 1 | Nikolas | nikolas@gmail.com |
| 2 | Martin | martin@gmail.com |
| 3 | _empty_ | sabin@gmail.com |
| 4 | Tyler | tyler@gmail.com |

### `null` and `undefined` in queries that affect _many_ records

This section will cover how `undefined` and `null` values affect the behavior of queries that interact with or create multiple records in a database.

#### Null

Consider the following Prisma Client query which searches for all users whose `name` value matches the provided `null` value:

```ts
const users = await prisma.user.findMany({
 where: {
 name: null,
 },
});
```

```json
[
 {
 "id": 3,
 "name": null,
 "email": "sabin@gmail.com"
 }
]
```

Because `null` was provided as the filter for the `name` column, Prisma Client will generate a query that searches for all records in the `User` table whose `name` column is _empty_.

#### Undefined

Now consider the scenario where you run the same query with `undefined` as the filter value on the `name` column:

```ts
const users = await prisma.user.findMany({
 where: {
 name: undefined,
 },
});
```

```json
[
 {
 "id": 1,
 "name": "Nikolas",
 "email": "nikolas@gmail.com"
 },
 {
 "id": 2,
 "name": "Martin",
 "email": "martin@gmail.com"
 },
 {
 "id": 3,
 "name": null,
 "email": "sabin@gmail.com"
 },
 {
 "id": 4,
 "name": "Tyler",
 "email": "tyler@gmail.com"
 }
]
```

Using `undefined` as a value in a filter essentially tells Prisma Client you have decided _not to define a filter_ for that column.

An equivalent way to write the above query would be:

```ts
const users = await prisma.user.findMany();
```

This query will select every row from the `User` table.

:::info

Using `undefined` as the value of any key in a Prisma Client query's parameter object will cause Prisma ORM to act as if that key was not provided at all.

:::

Although this section's examples focused on the `findMany` function, the same concepts apply to any function that can affect multiple records, such as `updateMany` and `deleteMany`.

### `null` and `undefined` in queries that affect _one_ record

This section will cover how `undefined` and `null` values affect the behavior of queries that interact with or create a single record in a database.

:::warning

`null` is not a valid filter value in a `findUnique()` query.

:::

The query behavior when using `null` and `undefined` in the filter criteria of a query that affects a single record is very similar to the behaviors described in the previous section.

#### Null

Consider the following query where `null` is used to filter the `name` column:

```ts
const user = await prisma.user.findFirst({
 where: {
 name: null,
 },
});
```

```json
[
 {
 "id": 3,
 "name": null,
 "email": "sabin@gmail.com"
 }
]
```

Because `null` was used as the filter on the `name` column, Prisma Client will generate a query that searches for the first record in the `User` table whose `name` value is _empty_.

#### Undefined

If `undefined` is used as the filter value on the `name` column instead, _the query will act as if no filter criteria was passed to that column at all_.

Consider the query below:

```ts
const user = await prisma.user.findFirst({
 where: {
 name: undefined,
 },
});
```

```json
[
 {
 "id": 1,
 "name": "Nikolas",
 "email": "nikolas@gmail.com"
 }
]
```

In this scenario, the query will return the very first record in the database.

Another way to represent the above query is:

```ts
const user = await prisma.user.findFirst();
```

Although this section's examples focused on the `findFirst` function, the same concepts apply to any function that affects a single record.

### `null` and `undefined` in a GraphQL resolver

For this example, consider a database based on the following Prisma schema:

```prisma
model User {
 id Int @id @default(autoincrement())
 email String @unique
 name String?
}
```

In the following GraphQL mutation that updates a user, both `authorEmail` and `name` accept `null`. From a GraphQL perspective, this means that fields are **optional**:

```ts
type Mutation {
 // Update author's email or name, or both - or neither!
 updateUser(id: Int!, authorEmail: String, authorName: String): User!
}
```

However, if you pass `null` values for `authorEmail` or `authorName` on to Prisma Client, the following will happen:

- If `args.authorEmail` is `null`, the query will **fail**. `email` does not accept `null`.
- If `args.authorName` is `null`, Prisma Client changes the value of `name` to `null`. This is probably not how you want an update to work.

```ts
updateUser: (parent, args, ctx: Context) => {
 return ctx.prisma.user.update({
 where: { id: Number(args.id) },
 data: {
 email: args.authorEmail, // email cannot be null // [!code highlight]
 name: args.authorName // name set to null - potentially unwanted behavior // [!code highlight]
 },
 })
},
```

Instead, set the value of `email` and `name` to `undefined` if the input value is `null`. Doing this is the same as not updating the field at all:

```ts
updateUser: (parent, args, ctx: Context) => {
 return ctx.prisma.user.update({
 where: { id: Number(args.id) },
 data: {
 email: args.authorEmail != null ? args.authorEmail : undefined, // If null, do nothing // [!code highlight]
 name: args.authorName != null ? args.authorName : undefined // If null, do nothing // [!code highlight]
 },
 })
},
```

### The effect of `null` and `undefined` on conditionals

There are some caveats to filtering with conditionals which might produce unexpected results. When filtering with conditionals you might expect one result but receive another given how Prisma Client treats nullable values.

The following table provides a high-level overview of how the different operators handle 0, 1 and `n` filters.

| Operator | 0 filters | 1 filter | n filters |
| -------- | ----------------- | ---------------------- | -------------------- |
| `OR` | return empty list | validate single filter | validate all filters |
| `AND` | return all items | validate single filter | validate all filters |
| `NOT` | return all items | validate single filter | validate all filters |

This example shows how an `undefined` parameter impacts the results returned by a query that uses the [`OR`](/orm/reference/prisma-client-reference#or) operator.

```ts
interface FormData {
 name: string;
 email?: string;
}

const formData: FormData = {
 name: "Emelie",
};

const users = await prisma.user.findMany({
 where: {
 OR: [
 {
 email: {
 contains: formData.email,
 },
 },
 ],
 },
});

// returns: []
```

The query receives filters from a formData object, which includes an optional email property. In this instance, the value of the email property is `undefined`. When this query is run no data is returned.

This is in contrast to the [`AND`](/orm/reference/prisma-client-reference#and) and [`NOT`](/orm/reference/prisma-client-reference) operators, which will both return all the users
if you pass in an `undefined` value.

> This is because passing an `undefined` value to an `AND` or `NOT` operator is the same
> as passing nothing at all, meaning the `findMany` query in the example will run without any filters and return all the users.

```ts
interface FormData {
 name: string;
 email?: string;
}

const formData: FormData = {
 name: "Emelie",
};

const users = await prisma.user.findMany({
 where: {
 AND: [
 {
 email: {
 contains: formData.email,
 },
 },
 ],
 },
});

// returns: { id: 1, email: 'ems@boop.com', name: 'Emelie' }

const users = await prisma.user.findMany({
 where: {
 NOT: [
 {
 email: {
 contains: formData.email,
 },
 },
 ],
 },
});

// returns: { id: 1, email: 'ems@boop.com', name: 'Emelie' }
```

---
title: Working with compound IDs and unique constraints
description: 'How to read, write, and filter by compound IDs and unique constraints'
url: /orm/prisma-client/special-fields-and-types/working-with-composite-ids-and-constraints
metaTitle: Working with compound IDs and unique constraints (Concepts)
metaDescription: 'How to read, write, and filter by compound IDs and unique constraints.'
---

Composite IDs and compound unique constraints can be defined in your Prisma schema using the [`@@id`](/orm/reference/prisma-schema-reference) and [`@@unique`](/orm/reference/prisma-schema-reference) attributes.

:::warning

**MongoDB does not support `@@id`**<br />
MongoDB does not support composite IDs, which means you cannot identify a model with a `@@id` attribute.

:::

A composite ID or compound unique constraint uses the combined values of two fields as a primary key or identifier in your database table. In the following example, the `postId` field and `userId` field are used as a composite ID for a `Like` table:

```prisma highlight=22;normal
model User {
 id Int @id @default(autoincrement())
 name String
 post Post[]
 likes Like[]
}

model Post {
 id Int @id @default(autoincrement())
 content String
 User User? @relation(fields: [userId], references: [id])
 userId Int?
 likes Like[]
}

model Like {
 postId Int
 userId Int
 User User @relation(fields: [userId], references: [id])
 Post Post @relation(fields: [postId], references: [id])

 @@id([postId, userId]) // [!code highlight]
}
```

Querying for records from the `Like` table (e.g. using `prisma.like.findMany()`) would return objects that look as follows:

```json
{
 "postId": 1,
 "userId": 1
}
```

Although there are only two fields in the response, those two fields make up a compound ID named `postId_userId`.

You can also create a named compound ID or compound unique constraint by using the `@@id` or `@@unique` attributes' `name` field. For example:

```prisma highlight=7;normal
model Like {
 postId Int
 userId Int
 User User @relation(fields: [userId], references: [id])
 Post Post @relation(fields: [postId], references: [id])

 @@id(name: "likeId", [postId, userId]) // [!code highlight]
}
```

## Where you can use compound IDs and unique constraints

Compound IDs and compound unique constraints can be used when working with _unique_ data.

Below is a list of Prisma Client functions that accept a compound ID or compound unique constraint in the `where` filter of the query:

- `findUnique()`
- `findUniqueOrThrow`
- `delete`
- `update`
- `upsert`

A composite ID and a composite unique constraint is also usable when creating relational data with `connect` and `connectOrCreate`.

## Filtering records by a compound ID or unique constraint

Although your query results will not display a compound ID or unique constraint as a field, you can use these compound values to filter your queries for unique records:

```ts highlight=3-6;normal
const like = await prisma.like.findUnique({
 where: {
 likeId: {
 userId: 1,
 postId: 1,
 },
 },
});
```

:::info

Note composite ID and compound unique constraint keys are only available as filter options for _unique_ queries such as `findUnique()` and `findUniqueOrThrow`. See the [section](/orm/prisma-client/special-fields-and-types/working-with-composite-ids-and-constraints#where-you-can-use-compound-ids-and-unique-constraints) above for a list of places these fields may be used.

:::

## Deleting records by a compound ID or unique constraint

A compound ID or compound unique constraint may be used in the `where` filter of a `delete` query:

```ts highlight=3-6;normal
const like = await prisma.like.delete({
 where: {
 likeId: {
 userId: 1,
 postId: 1,
 },
 },
});
```

## Updating and upserting records by a compound ID or unique constraint

A compound ID or compound unique constraint may be used in the `where` filter of an `update` query:

```ts highlight=3-6;normal
const like = await prisma.like.update({
 where: {
 likeId: {
 userId: 1,
 postId: 1,
 },
 },
 data: {
 postId: 2,
 },
});
```

They may also be used in the `where` filter of an `upsert` query:

```ts highlight=3-6;normal
await prisma.like.upsert({
 where: {
 likeId: {
 userId: 1,
 postId: 1,
 },
 },
 update: {
 userId: 2,
 },
 create: {
 userId: 2,
 postId: 1,
 },
});
```

## Filtering relation queries by a compound ID or unique constraint

Compound IDs and compound unique constraint can also be used in the `connect` and `connectOrCreate` keys used when connecting records to create a relationship.

For example, consider this query:

```ts highlight=6-9;normal
await prisma.user.create({
 data: {
 name: "Alice",
 likes: {
 connect: {
 likeId: {
 postId: 1,
 userId: 2,
 },
 },
 },
 },
});
```

The `likeId` compound ID is used as the identifier in the `connect` object that is used to locate the `Like` table's record that will be linked to the new user: `"Alice"`.

Similarly, the `likeId` can be used in `connectOrCreate`'s `where` filter to attempt to locate an existing record in the `Like` table:

```ts highlight=10-13;normal
await prisma.user.create({
 data: {
 name: "Alice",
 likes: {
 connectOrCreate: {
 create: {
 postId: 1,
 },
 where: {
 likeId: {
 postId: 1,
 userId: 1,
 },
 },
 },
 },
 },
});
```

---
title: Working with Json fields
description: 'How to read, write, and filter by Json fields'
url: /orm/prisma-client/special-fields-and-types/working-with-json-fields
metaTitle: Working with Json fields (Concepts)
metaDescription: 'How to read, write, and filter by Json fields.'
---

Use the [`Json`](/orm/reference/prisma-schema-reference#json) Prisma ORM field type to read, write, and perform basic filtering on JSON types in the underlying database. In the following example, the `User` model has an optional `Json` field named `extendedPetsData`:

```prisma highlight=6;normal
model User {
 id Int @id @default(autoincrement())
 email String @unique
 name String?
 posts Post[]
 extendedPetsData Json? // [!code highlight]
}
```

Example field value:

```json
{
 "pet1": {
 "petName": "Claudine",
 "petType": "House cat"
 },
 "pet2": {
 "petName": "Sunny",
 "petType": "Gerbil"
 }
}
```

The `Json` field supports a few additional types, such as `string` and `boolean`. These additional types exist to match the types supported by [`JSON.parse()`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/JSON/parse):

```ts
export type JsonValue = string | number | boolean | null | JsonObject | JsonArray;
```

## Use cases for JSON fields

Reasons to store data as JSON rather than representing data as related models include:

- You need to store data that does not have a consistent structure
- You are importing data from another system and do not want to map that data to Prisma models

## Reading a `Json` field

You can use the `Prisma.JsonArray` and `Prisma.JsonObject` utility classes to work with the contents of a `Json` field:

```ts
const { PrismaClient, Prisma } = require("@prisma/client");

const user = await prisma.user.findFirst({
 where: {
 id: 9,
 },
});

// Example extendedPetsData data:
// [{ name: 'Bob the dog' }, { name: 'Claudine the cat' }]

if (
 user?.extendedPetsData &&
 typeof user?.extendedPetsData === "object" &&
 Array.isArray(user?.extendedPetsData)
) {
 const petsObject = user?.extendedPetsData as Prisma.JsonArray;

 const firstPet = petsObject[0];
}
```

See also: [Advanced example: Update a nested JSON key value](#advanced-example-update-a-nested-json-key-value)

## Writing to a `Json` field

The following example writes a JSON object to the `extendedPetsData` field:

```ts
var json = [{ name: "Bob the dog" }, { name: "Claudine the cat" }] as Prisma.JsonArray;

const createUser = await prisma.user.create({
 data: {
 email: "birgitte@prisma.io",
 extendedPetsData: json,
 },
});
```

> **Note**: JavaScript objects (for example, `{ extendedPetsData: "none"}`) are automatically converted to JSON.

See also: [Advanced example: Update a nested JSON key value](#advanced-example-update-a-nested-json-key-value)

## Filter on a `Json` field (simple)

You can filter rows of `Json` type.

### Filter on exact field value

The following query returns all users where the value of `extendedPetsData` matches the `json` variable exactly:

```ts
var json = [{ name: "Bob the dog" }, { name: "Claudine the cat" }];

const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 equals: json,
 },
 },
});
```

The following query returns all users where the value of `extendedPetsData` does **not** match the `json` variable exactly:

```ts
var json = [{ name: "Bob the dog" }, { name: "Claudine the cat" }];

const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 not: json,
 },
 },
});
```

## Filter on a `Json` field (advanced)

You can also filter rows by the data inside a `Json` field. We call this **advanced `Json` filtering**. This functionality is supported by [PostgreSQL](/orm/core-concepts/supported-databases/postgresql) and [MySQL](/orm/core-concepts/supported-databases/mysql) only with [different syntaxes for the `path` option](#path-syntax-depending-on-database).

:::warning

PostgreSQL does not support [filtering on object key values in arrays](#filtering-on-object-key-value-inside-array).

:::

### `path` syntax depending on database

The filters below use a `path` option to select specific parts of the `Json` value to filter on. The implementation of that filtering differs between connectors:

- The [MySQL connector](/orm/core-concepts/supported-databases/mysql) uses [MySQL's implementation of JSON path](https://dev.mysql.com/doc/refman/8.0/en/json.html#json-path-syntax)
- The [PostgreSQL connector](/orm/core-concepts/supported-databases/postgresql) uses the custom JSON functions and operators [supported in version 12 _and earlier_](https://www.postgresql.org/docs/11/functions-json.html)

For example, the following is a valid MySQL `path` value:

```
$petFeatures.petName
```

The following is a valid PostgreSQL `path` value:

```
["petFeatures", "petName"]
```

### Filter on object property

You can filter on a specific property inside a block of JSON. In the following examples, the value of `extendedPetsData` is a one-dimensional, unnested JSON object:

```json highlight=11;normal
{
 "petName": "Claudine",
 "petType": "House cat"
}
```

The following query returns all users where the value of `petName` is `"Claudine"`:

```ts tab="PostgreSQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: ["petName"],
 equals: "Claudine",
 },
 },
});
```

```ts tab="MySQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: "$.petName",
 equals: "Claudine",
 },
 },
});
```

The following query returns all users where the value of `petType` _contains_ `"cat"`:

```ts tab="PostgreSQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: ["petType"],
 string_contains: "cat",
 },
 },
});
```

```ts tab="MySQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: "$.petType",
 string_contains: "cat",
 },
 },
});
```

The following string filters are available:

- [`string_contains`](/orm/reference/prisma-client-reference#string_contains)
- [`string_starts_with`](/orm/reference/prisma-client-reference#string_starts_with)
- [`string_ends_with`](/orm/reference/prisma-client-reference#string_ends_with) .

To use case insensitive filter with these, you can use the [`mode`](/orm/reference/prisma-client-reference#mode) option:

```ts tab="PostgreSQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: ["petType"],
 string_contains: "cat",
 mode: "insensitive",
 },
 },
});
```

```ts tab="MySQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: "$.petType",
 string_contains: "cat",
 mode: "insensitive",
 },
 },
});
```

### Filter on nested object property

You can filter on nested JSON properties. In the following examples, the value of `extendedPetsData` is a JSON object with several levels of nesting.

```json
{
 "pet1": {
 "petName": "Claudine",
 "petType": "House cat"
 },
 "pet2": {
 "petName": "Sunny",
 "petType": "Gerbil",
 "features": {
 "eyeColor": "Brown",
 "furColor": "White and black"
 }
 }
}
```

The following query returns all users where `"pet2"` &rarr; `"petName"` is `"Sunny"`:

```ts tab="PostgreSQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: ["pet2", "petName"],
 equals: "Sunny",
 },
 },
});
```

```ts tab="MySQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: "$.pet2.petName",
 equals: "Sunny",
 },
 },
});
```

The following query returns all users where:

- `"pet2"` &rarr; `"petName"` is `"Sunny"`
- `"pet2"` &rarr; `"features"` &rarr; `"furColor"` contains `"black"`

```ts tab="PostgreSQL"
const getUsers = await prisma.user.findMany({
 where: {
 AND: [
 {
 extendedPetsData: {
 path: ["pet2", "petName"],
 equals: "Sunny",
 },
 },
 {
 extendedPetsData: {
 path: ["pet2", "features", "furColor"],
 string_contains: "black",
 },
 },
 ],
 },
});
```

```ts tab="MySQL"
const getUsers = await prisma.user.findMany({
 where: {
 AND: [
 {
 extendedPetsData: {
 path: "$.pet2.petName",
 equals: "Sunny",
 },
 },
 {
 extendedPetsData: {
 path: "$.pet2.features.furColor",
 string_contains: "black",
 },
 },
 ],
 },
});
```

### Filtering on an array value

You can filter on the presence of a specific value in a scalar array (strings, integers). In the following example, the value of `extendedPetsData` is an array of strings:

```json
["Claudine", "Sunny"]
```

The following query returns all users with a pet named `"Claudine"`:

```ts tab="PostgreSQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 array_contains: ["Claudine"],
 },
 },
});
```

```ts tab="MySQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 array_contains: "Claudine",
 },
 },
});
```

:::info

In PostgreSQL, the value of `array_contains` must be an array and not a string, even if the array only contains a single value.

:::

The following array filters are available:

- [`array_contains`](/orm/reference/prisma-client-reference#array_contains)
- [`array_starts_with`](/orm/reference/prisma-client-reference#array_starts_with)
- [`array_ends_with`](/orm/reference/prisma-client-reference#array_ends_with)

### Filtering on nested array value

You can filter on the presence of a specific value in a scalar array (strings, integers). In the following examples, the value of `extendedPetsData` includes nested scalar arrays of names:

```json
{
 "cats": { "owned": ["Bob", "Sunny"], "fostering": ["Fido"] },
 "dogs": { "owned": ["Ella"], "fostering": ["Prince", "Empress"] }
}
```

#### Scalar value arrays

The following query returns all users that foster a cat named `"Fido"`:

```ts tab="PostgreSQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: ["cats", "fostering"],
 array_contains: ["Fido"],
 },
 },
});
```

```ts tab="MySQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: "$.cats.fostering",
 array_contains: "Fido",
 },
 },
});
```

:::info

In PostgreSQL, the value of `array_contains` must be an array and not a string, even if the array only contains a single value.

:::

The following query returns all users that foster cats named `"Fido"` _and_ `"Bob"`:

```ts tab="PostgreSQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: ["cats", "fostering"],
 array_contains: ["Fido", "Bob"],
 },
 },
});
```

```ts tab="MySQL"
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: "$.cats.fostering",
 array_contains: ["Fido", "Bob"],
 },
 },
});
```

#### JSON object arrays

```ts tab="PostgreSQL"
const json = [{ status: "expired", insuranceID: 92 }];

const checkJson = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: ["insurances"],
 array_contains: json,
 },
 },
});
```

```ts tab="MySQL"
const json = { status: "expired", insuranceID: 92 };

const checkJson = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: "$.insurances",
 array_contains: json,
 },
 },
});
```

- If you are using PostgreSQL, you must pass in an array of objects to match, even if that array only contains one object:

 ```json5
 [{ status: "expired", insuranceID: 92 }]
 // PostgreSQL
 ```

 If you are using MySQL, you must pass in a single object to match:

 ```json5
 { status: "expired", insuranceID: 92 }
 // MySQL
 ```

- If your filter array contains multiple objects, PostgreSQL will only return results if _all_ objects are present - not if at least one object is present.

- You must set `array_contains` to a JSON object, not a string. If you use a string, Prisma Client escapes the quotation marks and the query will not return results. For example:

 ```ts
 array_contains: '[{"status": "expired", "insuranceID": 92}]';
 ```

 is sent to the database as:

 ```
 [{\"status\": \"expired\", \"insuranceID\": 92}]
 ```

### Targeting an array element by index

You can filter on the value of an element in a specific position.

```json
{ "owned": ["Bob", "Sunny"], "fostering": ["Fido"] }
```

```ts tab="PostgreSQL"
const getUsers = await prisma.user.findMany({
 where: {
 comments: {
 path: ["owned", "1"],
 string_contains: "Bob",
 },
 },
});
```

```ts tab="MySQL"
const getUsers = await prisma.user.findMany({
 where: {
 comments: {
 path: "$.owned[1]",
 string_contains: "Bob",
 },
 },
});
```

### Filtering on object key value inside array

Depending on your provider, you can filter on the key value of an object inside an array.

:::warning

Filtering on object key values within an array is **only** supported by the [MySQL database connector](/orm/core-concepts/supported-databases/mysql). However, you can still [filter on the presence of entire JSON objects](#json-object-arrays).

:::

In the following example, the value of `extendedPetsData` is an array of objects with a nested `insurances` array, which contains two objects:

```json
[
 {
 "petName": "Claudine",
 "petType": "House cat",
 "insurances": [
 { "insuranceID": 92, "status": "expired" },
 { "insuranceID": 12, "status": "active" }
 ]
 },
 {
 "petName": "Sunny",
 "petType": "Gerbil"
 },
 {
 "petName": "Gerald",
 "petType": "Corn snake"
 },
 {
 "petName": "Nanna",
 "petType": "Moose"
 }
]
```

The following query returns all users where at least one pet is a moose:

```ts
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: "$[*].petType",
 array_contains: "Moose",
 },
 },
});
```

- `$[*]` is the root array of pet objects
- `petType` matches the `petType` key in any pet object

The following query returns all users where at least one pet has an expired insurance:

```ts
const getUsers = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: "$[*].insurances[*].status",
 array_contains: "expired",
 },
 },
});
```

- `$[*]` is the root array of pet objects
- `insurances[*]` matches any `insurances` array inside any pet object
- `status` matches any `status` key in any insurance object

## Advanced example: Update a nested JSON key value

The following example assumes that the value of `extendedPetsData` is some variation of the following:

```json
{
 "petName": "Claudine",
 "petType": "House cat",
 "insurances": [
 { "insuranceID": 92, "status": "expired" },
 { "insuranceID": 12, "status": "active" }
 ]
}
```

The following example:

1. Gets all users
1. Change the `"status"` of each insurance object to `"expired"`
1. Get all users that have an expired insurance where the ID is `92`

```ts tab="PostgreSQL"
const userQueries: string | any[] = [];

getUsers.forEach((user) => {
 if (
 user.extendedPetsData &&
 typeof user.extendedPetsData === "object" &&
 !Array.isArray(user.extendedPetsData)
 ) {
 const petsObject = user.extendedPetsData as Prisma.JsonObject;

 const i = petsObject["insurances"];

 if (i && typeof i === "object" && Array.isArray(i)) {
 const insurancesArray = i as Prisma.JsonArray;

 insurancesArray.forEach((i) => {
 if (i && typeof i === "object" && !Array.isArray(i)) {
 const insuranceObject = i as Prisma.JsonObject;

 insuranceObject["status"] = "expired";
 }
 });

 const whereClause = Prisma.validator<Prisma.UserWhereInput>()({
 id: user.id,
 });

 const dataClause = Prisma.validator<Prisma.UserUpdateInput>()({
 extendedPetsData: petsObject,
 });

 userQueries.push(
 prisma.user.update({
 where: whereClause,
 data: dataClause,
 }),
 );
 }
 }
});

if (userQueries.length > 0) {
 console.log(userQueries.length + " queries to run!");
 await prisma.$transaction(userQueries);
}

const json = [{ status: "expired", insuranceID: 92 }];

const checkJson = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: ["insurances"],
 array_contains: json,
 },
 },
});

console.log(checkJson.length);
```

```ts tab="MySQL"
const userQueries: string | any[] = [];

getUsers.forEach((user) => {
 if (
 user.extendedPetsData &&
 typeof user.extendedPetsData === "object" &&
 !Array.isArray(user.extendedPetsData)
 ) {
 const petsObject = user.extendedPetsData as Prisma.JsonObject;

 const insuranceList = petsObject["insurances"]; // is a Prisma.JsonArray

 if (Array.isArray(insuranceList)) {
 insuranceList.forEach((insuranceItem) => {
 if (insuranceItem && typeof insuranceItem === "object" && !Array.isArray(insuranceItem)) {
 insuranceItem["status"] = "expired"; // is a Prisma.JsonObject
 }
 });

 const whereClause = Prisma.validator<Prisma.UserWhereInput>()({
 id: user.id,
 });

 const dataClause = Prisma.validator<Prisma.UserUpdateInput>()({
 extendedPetsData: petsObject,
 });

 userQueries.push(
 prisma.user.update({
 where: whereClause,
 data: dataClause,
 }),
 );
 }
 }
});

if (userQueries.length > 0) {
 console.log(userQueries.length + " queries to run!");
 await prisma.$transaction(userQueries);
}

const json = { status: "expired", insuranceID: 92 };

const checkJson = await prisma.user.findMany({
 where: {
 extendedPetsData: {
 path: "$.insurances",
 array_contains: json,
 },
 },
});

console.log(checkJson.length);
```

## Using `null` Values

There are two types of `null` values possible for a `JSON` field in an SQL database.

- Database `NULL`: The value in the database is a `NULL`.
- JSON `null`: The value in the database contains a JSON value that is `null`.

To differentiate between these possibilities, we've introduced three _null enums_ you can use:

- `JsonNull`: Represents the `null` value in JSON.
- `DbNull`: Represents the `NULL` value in the database.
- `AnyNull`: Represents both `null` JSON values and `NULL` database values. (Only when filtering)

:::info

- When filtering using any of the _null enums_ you can not use a shorthand and leave the `equals` operator off.
- These _null enums_ do not apply to MongoDB because there the difference between a JSON `null` and a database `NULL` does not exist.
- The _null enums_ do not apply to the `array_contains` operator in all databases because there can only be a JSON `null` within a JSON array. Since there cannot be a database `NULL` within a JSON array, `{ array_contains: null }` is not ambiguous.

:::

For example:

```prisma
model Log {
 id Int @id
 meta Json
}
```

Here is an example of using `AnyNull`:

```ts highlight=7;normal
import { Prisma } from "@prisma/client";

prisma.log.findMany({
 where: {
 data: {
 meta: {
 equals: Prisma.AnyNull,
 },
 },
 },
});
```

### Inserting `null` Values

This also applies to `create`, `update` and `upsert`. To insert a `null` value
into a `Json` field, you would write:

```ts highlight=5;normal
import { Prisma } from "@prisma/client";

prisma.log.create({
 data: {
 meta: Prisma.JsonNull,
 },
});
```

And to insert a database `NULL` into a `Json` field, you would write:

```ts highlight=5;normal
import { Prisma } from "@prisma/client";

prisma.log.create({
 data: {
 meta: Prisma.DbNull,
 },
});
```

### Filtering by `null` Values

To filter by `JsonNull` or `DbNull`, you would write:

```ts highlight=6;normal
import { Prisma } from "@prisma/client";

prisma.log.findMany({
 where: {
 meta: {
 equals: Prisma.AnyNull,
 },
 },
});
```

:::info

These _null enums_ do not apply to MongoDB because MongoDB does not differentiate between a JSON `null` and a database `NULL`. They also do not apply to the `array_contains` operator in all databases because there can only be a JSON `null` within a JSON array. Since there cannot be a database `NULL` within a JSON array, `{ array_contains: null }` is not ambiguous.

:::

## Typed `Json` Fields

Prisma's `Json` fields are untyped by default. To add strong typing, you can use the external package [prisma-json-types-generator](https://www.npmjs.com/package/prisma-json-types-generator).

1. First, install the package and add the generator to your `schema.prisma`:

 ```npm
 npm install -D prisma-json-types-generator
 ```

 ```prisma title="schema.prisma"
 generator client {
 provider = "prisma-client"
 output = "./generated"
 }

 generator json {
 provider = "prisma-json-types-generator"
 }
 ```

2. Next, link a field to a TypeScript type using an [AST comment](/orm/prisma-schema/overview#comments).

 ```prisma highlight=4;normal title="schema.prisma" showLineNumbers
 model Log {
 id Int @id

 /// [LogMetaType] // [!code highlight]
 meta Json
 }
 ```

3. Then, define `LogMetaType` in a type declaration file (e.g., `types.ts`) that is included in your `tsconfig.json`.

 ```ts title="types.ts" showLineNumbers
 declare global {
 namespace PrismaJson {
 type LogMetaType = { timestamp: number; host: string };
 }
 }

 // This file must be a module.
 export {};
 ```

Now, `Log.meta` will be strongly typed as `{ timestamp: number; host: string }`.

### Typing `String` Fields and Advanced Features

You can also apply these techniques to `String` fields. This is especially useful for creating string-based enums directly in your schema when your database does not support enum types.

```prisma
model Post {
 id Int @id

 /// !['draft' | 'published']
 status String

 /// [LogMetaType]
 meta Json[]
}
```

This results in `post.status` being strongly typed as `'draft' | 'published'` and `post.meta` as `LogMetaType[]`.

For a complete guide on configuration, monorepo setup, and other advanced features, please refer to the [official `prisma-json-types-generator` documentation](https://github.com/arthurfiorette/prisma-json-types-generator#readme).

## `Json` FAQs

### Can you select a subset of JSON key/values to return?

No - it is not yet possible to [select which JSON elements to return](https://github.com/prisma/prisma/issues/2431). Prisma Client returns the entire JSON object.

### Can you filter on the presence of a specific key?

No - it is not yet possible to filter on the presence of a specific key.

### Is case insensitive filtering supported?

Yes - you can use the `mode: 'insensitive'` option with string filters like `string_contains`, `string_starts_with`, and `string_ends_with`. See [Filter on object property](#filter-on-object-property) for examples.

### Can you sort an object property within a JSON value?

No, [sorting object properties within a JSON value](https://github.com/prisma/prisma/issues/10346) (order-by-prop) is not currently supported.

### How to set a default value for JSON fields?

When you want to set a `@default` value the `Json` type, you need to enclose it with double-quotes inside the `@default` attribute (and potentially escape any "inner" double-quotes using a backslash), for example:

```prisma
model User {
 id Int @id @default(autoincrement())
 json1 Json @default("[]")
 json2 Json @default("{ \"hello\": \"world\" }")
}
```

---
title: Working with scalar lists
description: 'How to read, write, and filter by scalar lists / arrays'
url: /orm/prisma-client/special-fields-and-types/working-with-scalar-lists-arrays
metaTitle: Working with scalar lists/arrays (Concepts)
metaDescription: 'How to read, write, and filter by scalar lists / arrays.'
---

[Scalar lists](/orm/reference/prisma-schema-reference#-modifier) are represented by the `[]` modifier and are only available if the underlying database supports scalar lists. The following example has one scalar `String` list named `pets`:

```prisma highlight=4;normal tab="Relational databases"
model User {
 id Int @id @default(autoincrement())
 name String
 pets String[] // [!code highlight]
}
```

```prisma highlight=4;normal tab="MongoDB"
model User {
 id String @id @default(auto()) @map("_id") @db.ObjectId
 name String
 pets String[] // [!code highlight]
}
```

Example field value:

```json5
["Fido", "Snoopy", "Brian"]
```

## Setting the value of a scalar list

The following example demonstrates how to [`set`](/orm/reference/prisma-client-reference) the value of a scalar list (`coinflips`) when you create a model:

```ts
const createdUser = await prisma.user.create({
 data: {
 email: "eloise@prisma.io",
 coinflips: [true, true, true, false, true],
 },
});
```

## Unsetting the value of a scalar list

:::warning

This method is available on MongoDB only.

:::

The following example demonstrates how to [`unset`](/orm/reference/prisma-client-reference#unset) the value of a scalar list (`coinflips`):

```ts
const createdUser = await prisma.user.create({
 data: {
 email: "eloise@prisma.io",
 coinflips: {
 unset: true,
 },
 },
});
```

Unlike `set: null`, `unset` removes the list entirely.

## Adding items to a scalar list

:::note

Available for PostgreSQL, CockroachDB, and MongoDB.

:::

Use the [`push`](/orm/reference/prisma-client-reference#push) method to add a single value to a scalar list:

```ts
const userUpdate = await prisma.user.update({
 where: {
 id: 9,
 },
 data: {
 coinflips: {
 push: true,
 },
 },
});
```

In earlier versions, you have to overwrite the entire value. The following example retrieves user, uses `push()` to add three new coin flips, and overwrites the `coinflips` field in an `update`:

```ts
const user = await prisma.user.findUnique({
 where: {
 email: "eloise@prisma.io",
 },
});

if (user) {
 console.log(user.coinflips);

 user.coinflips.push(true, true, false);

 const updatedUser = await prisma.user.update({
 where: {
 email: "eloise@prisma.io",
 },
 data: {
 coinflips: user.coinflips,
 },
 });

 console.log(updatedUser.coinflips);
}
```

## Filtering scalar lists

:::note

Available for PostgreSQL, CockroachDB, and MongoDB.

:::

Use [scalar list filters](/orm/reference/prisma-client-reference#scalar-list-filters) to filter for records with scalar lists that match a specific condition. The following example returns all posts where the tags list includes `databases` _and_ `typescript`:

```ts
const posts = await prisma.post.findMany({
 where: {
 tags: {
 hasEvery: ["databases", "typescript"],
 },
 },
});
```

### `NULL` values in arrays

:::note

This section applies to PostgreSQL and CockroachDB.

:::

When using scalar list filters with a relational database connector, array fields with a `NULL` value are not considered by the following conditions:

- `NOT` (array does not contain X)
- `isEmpty` (array is empty)

This means that records you might expect to see are not returned. Consider the following examples:

- The following query returns all posts where the `tags` **do not** include `databases`:

 ```ts
 const posts = await prisma.post.findMany({
 where: {
 NOT: {
 tags: {
 has: "databases",
 },
 },
 },
 });
 ```

 - ✔ Arrays that do not contain `"databases"`, such as `{"typescript", "graphql"}`
 - ✔ Empty arrays, such as `[]`

 The query does not return:
 - ✘ `NULL` arrays, even though they do not contain `"databases"`

The following query returns all posts where `tags` is empty:

```ts
const posts = await prisma.post.findMany({
 where: {
 tags: {
 isEmpty: true,
 },
 },
});
```

The query returns:

- ✔ Empty arrays, such as `[]`

The query does not return:

- ✘ `NULL` arrays, even though they could be considered empty

To work around this issue, you can set the default value of array fields to `[]`.

---
title: Integration testing
description: Learn how to setup and run integration tests with Prisma and Docker
url: /orm/prisma-client/testing/integration-testing
metaTitle: Integration testing with Prisma
metaDescription: Learn how to setup and run integration tests with Prisma and Docker
---

Integration tests focus on testing how separate parts of the program work together. In the context of applications using a database, integration tests usually require a database to be available and contain data that is convenient to the scenarios intended to be tested.

One way to simulate a real world environment is to use [Docker](https://www.docker.com/get-started/) to encapsulate a database and some test data. This can be spun up and torn down with the tests and so operate as an isolated environment away from your production databases.

> **Note:** Prisma's testing series offers a more detailed walkthrough of setting up integration tests against a real database if you want a deeper companion resource.

## Prerequisites

This guide assumes you have [Docker](https://docs.docker.com/get-started/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed on your machine as well as `Jest` setup in your project.

See our [system requirements](/orm/reference/system-requirements) for all minimum version requirements.

The following e-commerce schema will be used throughout the guide. This varies from the traditional `User` and `Post` models used in other parts of the docs, mainly because it is unlikely you will be running integration tests against your blog.

<details>

<summary>Ecommerce schema</summary>

```prisma title="schema.prisma"
// Can have 1 customer
// Can have many order details
model CustomerOrder {
 id Int @id @default(autoincrement())
 createdAt DateTime @default(now())
 customer Customer @relation(fields: [customerId], references: [id])
 customerId Int
 orderDetails OrderDetails[]
}

// Can have 1 order
// Can have many products
model OrderDetails {
 id Int @id @default(autoincrement())
 products Product @relation(fields: [productId], references: [id])
 productId Int
 order CustomerOrder @relation(fields: [orderId], references: [id])
 orderId Int
 total Decimal
 quantity Int
}

// Can have many order details
// Can have 1 category
model Product {
 id Int @id @default(autoincrement())
 name String
 description String
 price Decimal
 sku Int
 orderDetails OrderDetails[]
 category Category @relation(fields: [categoryId], references: [id])
 categoryId Int
}

// Can have many products
model Category {
 id Int @id @default(autoincrement())
 name String
 products Product[]
}

// Can have many orders
model Customer {
 id Int @id @default(autoincrement())
 email String @unique
 address String?
 name String?
 orders CustomerOrder[]
}
```

</details>

The guide uses a singleton pattern for Prisma Client setup. Refer to the [singleton](/orm/prisma-client/testing/unit-testing#singleton) docs for a walk through of how to set that up.

## Add Docker to your project

![Docker compose code pointing towards image of container holding a Postgres database](/img/orm/prisma-client/testing/Docker_Diagram_V1.png)

With Docker and Docker compose both installed on your machine you can use them in your project.

1. Begin by creating a `docker-compose.yml` file at your projects root. Here you will add a Postgres image and specify the environments credentials.

```yml title="docker-compose.yml"
# Set the version of docker compose to use
version: "3.9"

# The containers that compose the project
services:
 db:
 image: postgres:13
 restart: always
 container_name: integration-tests-prisma
 ports:
 - "5433:5432"
 environment:
 POSTGRES_USER: prisma
 POSTGRES_PASSWORD: prisma
 POSTGRES_DB: tests
```

> **Note**: The compose version used here (`3.9`) is the latest at the time of writing, if you are following along be sure to use the same version for consistency.

The `docker-compose.yml` file defines the following:

- The Postgres image (`postgres`) and version tag (`:13`). This will be downloaded if you do not have it locally available.
- The port `5433` is mapped to the internal (Postgres default) port `5432`. This will be the port number the database is exposed on externally.
- The database user credentials are set and the database given a name.

2. To connect to the database in the container, create a new connection string with the credentials defined in the `docker-compose.yml` file. For example:

```bash title=".env.test"
DATABASE_URL="postgresql://prisma:prisma@localhost:5433/tests"
```

:::info

The above `.env.test` file is used as part of a multiple `.env` file setup. Checkout the [using multiple .env files.](/orm/more/dev-environment/environment-variables) section to learn more about setting up your project with multiple `.env` files

:::

3. To create the container in a detached state so that you can continue to use the terminal tab, run the following command:

```bash
docker compose up -d
```

4. Next you can check that the database has been created by executing a `psql` command inside the container. Make a note of the container id.

 ```
 docker ps
 ```

 ```bash
 CONTAINER ID IMAGE COMMAND CREATED STATUS PORTS NAMES
 1322e42d833f postgres:13 "docker-entrypoint.s…" 2 seconds ago Up 1 second 0.0.0.0:5433->5432/tcp integration-tests-prisma
 ```

> **Note**: The container id is unique to each container, you will see a different id displayed.

5. Using the container id from the previous step, run `psql` in the container, login with the created user and check the database is created:

 ```
 docker exec -it 1322e42d833f psql -U prisma tests
 ```

 ```bash
 tests=# \l
 List of databases
 Name | Owner | Encoding | Collate | Ctype | Access privileges

 postgres | prisma | UTF8 | en_US.utf8 | en_US.utf8 |
 template0 | prisma | UTF8 | en_US.utf8 | en_US.utf8 | =c/prisma +
 | | | | | prisma=CTc/prisma
 template1 | prisma | UTF8 | en_US.utf8 | en_US.utf8 | =c/prisma +
 | | | | | prisma=CTc/prisma
 tests | prisma | UTF8 | en_US.utf8 | en_US.utf8 |
 (4 rows)
 ```

## Integration testing

Integration tests will be run against a database in a **dedicated test environment** instead of the production or development environments.

### The flow of operations

The flow for running said tests goes as follows:

1. Start the container and create the database
1. Migrate the schema
1. Run the tests
1. Destroy the container

Each test suite will seed the database before all the test are run. After all the tests in the suite have finished, the data from all the tables will be dropped and the connection terminated.

### The function to test

The ecommerce application you are testing has a function which creates an order. This function does the following:

- Accepts input about the customer making the order
- Accepts input about the product being ordered
- Checks if the customer has an existing account
- Checks if the product is in stock
- Returns an "Out of stock" message if the product doesn't exist
- Creates an account if the customer doesn't exist in the database
- Create the order

An example of how such a function might look can be seen below:

```ts title="create-order.ts"
import prisma from "../client";

export interface Customer {
 id?: number;
 name?: string;
 email: string;
 address?: string;
}

export interface OrderInput {
 customer: Customer;
 productId: number;
 quantity: number;
}

/**
 * Creates an order with customer.
 * @param input The order parameters
 */
export async function createOrder(input: OrderInput) {
 const { productId, quantity, customer } = input;
 const { name, email, address } = customer;

 // Get the product
 const product = await prisma.product.findUnique({
 where: {
 id: productId,
 },
 });

 // If the product is null its out of stock, return error.
 if (!product) return new Error("Out of stock");

 // If the customer is new then create the record, otherwise connect via their unique email
 await prisma.customerOrder.create({
 data: {
 customer: {
 connectOrCreate: {
 create: {
 name,
 email,
 address,
 },
 where: {
 email,
 },
 },
 },
 orderDetails: {
 create: {
 total: product.price,
 quantity,
 products: {
 connect: {
 id: product.id,
 },
 },
 },
 },
 },
 });
}
```

### The test suite

The following tests will check if the `createOrder` function works as it should do. They will test:

- Creating a new order with a new customer
- Creating an order with an existing customer
- Show an "Out of stock" error message if a product doesn't exist

Before the test suite is run the database is seeded with data. After the test suite has finished a [`deleteMany`](/orm/reference/prisma-client-reference#deletemany) is used to clear the database of its data.

:::tip

Using `deleteMany` may suffice in situations where you know ahead of time how your schema is structured. This is because the operations need to be executed in the correct order according to how the model relations are setup.

However, this doesn't scale as well as having a more generic solution that maps over your models and performs a truncate on them. For those scenarios and examples of using raw SQL queries see [Deleting all data with raw SQL / `TRUNCATE`](/orm/prisma-client/queries/crud#deleting-all-data-with-raw-sql--truncate)

:::

```ts title="__tests__/create-order.ts"
import prisma from "../src/client";
import { createOrder, Customer, OrderInput } from "../src/functions/index";

beforeAll(async () => {
 // create product categories
 await prisma.category.createMany({
 data: [{ name: "Wand" }, { name: "Broomstick" }],
 });

 console.log("✨ 2 categories successfully created!");

 // create products
 await prisma.product.createMany({
 data: [
 {
 name: 'Holly, 11", phoenix feather',
 description: "Harry Potters wand",
 price: 100,
 sku: 1,
 categoryId: 1,
 },
 {
 name: "Nimbus 2000",
 description: "Harry Potters broom",
 price: 500,
 sku: 2,
 categoryId: 2,
 },
 ],
 });

 console.log("✨ 2 products successfully created!");

 // create the customer
 await prisma.customer.create({
 data: {
 name: "Harry Potter",
 email: "harry@hogwarts.io",
 address: "4 Privet Drive",
 },
 });

 console.log("✨ 1 customer successfully created!");
});

afterAll(async () => {
 const deleteOrderDetails = prisma.orderDetails.deleteMany();
 const deleteProduct = prisma.product.deleteMany();
 const deleteCategory = prisma.category.deleteMany();
 const deleteCustomerOrder = prisma.customerOrder.deleteMany();
 const deleteCustomer = prisma.customer.deleteMany();

 await prisma.$transaction([
 deleteOrderDetails,
 deleteProduct,
 deleteCategory,
 deleteCustomerOrder,
 deleteCustomer,
 ]);

 await prisma.$disconnect();
});

it("should create 1 new customer with 1 order", async () => {
 // The new customers details
 const customer: Customer = {
 id: 2,
 name: "Hermione Granger",
 email: "hermione@hogwarts.io",
 address: "2 Hampstead Heath",
 };
 // The new orders details
 const order: OrderInput = {
 customer,
 productId: 1,
 quantity: 1,
 };

 // Create the order and customer
 await createOrder(order);

 // Check if the new customer was created by filtering on unique email field
 const newCustomer = await prisma.customer.findUnique({
 where: {
 email: customer.email,
 },
 });

 // Check if the new order was created by filtering on unique email field of the customer
 const newOrder = await prisma.customerOrder.findFirst({
 where: {
 customer: {
 email: customer.email,
 },
 },
 });

 // Expect the new customer to have been created and match the input
 expect(newCustomer).toEqual(customer);
 // Expect the new order to have been created and contain the new customer
 expect(newOrder).toHaveProperty("customerId", 2);
});

it("should create 1 order with an existing customer", async () => {
 // The existing customers email
 const customer: Customer = {
 email: "harry@hogwarts.io",
 };
 // The new orders details
 const order: OrderInput = {
 customer,
 productId: 1,
 quantity: 1,
 };

 // Create the order and connect the existing customer
 await createOrder(order);

 // Check if the new order was created by filtering on unique email field of the customer
 const newOrder = await prisma.customerOrder.findFirst({
 where: {
 customer: {
 email: customer.email,
 },
 },
 });

 // Expect the new order to have been created and contain the existing customer with an id of 1 (Harry Potter from the seed script)
 expect(newOrder).toHaveProperty("customerId", 1);
});

it("should show 'Out of stock' message if productId doesn't exit", async () => {
 // The existing customers email
 const customer: Customer = {
 email: "harry@hogwarts.io",
 };
 // The new orders details
 const order: OrderInput = {
 customer,
 productId: 3,
 quantity: 1,
 };

 // The productId supplied doesn't exit so the function should return an "Out of stock" message
 await expect(createOrder(order)).resolves.toEqual(new Error("Out of stock"));
});
```

## Running the tests

This setup isolates a real world scenario so that you can test your applications functionality against real data in a controlled environment.

You can add some scripts to your projects `package.json` file which will setup the database and run the tests, then afterwards manually destroy the container.

:::warning

If the test doesn't work for you, you'll need to ensure the test database is properly set up and ready before the suite runs.

:::

```json title="package.json"
 "scripts": {
 "docker:up": "docker compose up -d",
 "docker:down": "docker compose down",
 "test": "yarn docker:up && yarn prisma migrate deploy && jest -i"
 },
```

The `test` script does the following:

1. Runs `docker compose up -d` to create the container with the Postgres image and database.
1. Applies the migrations found in `./prisma/migrations/` directory to the database, this creates the tables in the container's database.
1. Executes the tests.

Once you are satisfied you can run `yarn docker:down` to destroy the container, its database and any test data.

---
title: Unit testing
description: Learn how to setup and run unit tests with Prisma Client
url: /orm/prisma-client/testing/unit-testing
metaTitle: Unit testing with Prisma ORM
metaDescription: Learn how to setup and run unit tests with Prisma Client
---

Unit testing aims to isolate a small portion (unit) of code and test it for logically predictable behaviors. It generally involves mocking objects or server responses to simulate real world behaviors. Some benefits to unit testing include:

- Quickly find and isolate bugs in code.
- Provides documentation for each module of code by way of indicating what certain code blocks should be doing.
- A helpful gauge that a refactor has gone well. The tests should still pass after code has been refactored.

In the context of Prisma ORM, this generally means testing a function which makes database calls using Prisma Client.

A single test should focus on how your function logic handles different inputs (such as a null value or an empty list).

This means that you should aim to remove as many dependencies as possible, such as external services and databases, to keep the tests and their environments as lightweight as possible.

> **Note**: Prisma's testing series provides additional background on unit testing patterns with Prisma ORM if you want to go deeper after this guide.

## Prerequisites

This guide assumes you have the JavaScript testing library [`Jest`](https://jestjs.io/) and [`ts-jest`](https://github.com/kulshekhar/ts-jest) already setup in your project.

## Mocking Prisma Client

To ensure your unit tests are isolated from external factors you can mock Prisma Client, this means you get the benefits of being able to use your schema (**_type-safety_**), without having to make actual calls to your database when your tests are run.

This guide will cover two approaches to mocking Prisma Client, a singleton instance and dependency injection. Both have their merits depending on your use cases. To help with mocking Prisma Client the [`jest-mock-extended`](https://github.com/marchaos/jest-mock-extended) package will be used.

```npm
npm install jest-mock-extended@2.0.4 --save-dev
```

:::danger

At the time of writing, this guide uses `jest-mock-extended` version `^2.0.4`.

:::

### Singleton

The following steps guide you through mocking Prisma Client using a singleton pattern.

1. Create a file at your projects root called `client.ts` and add the following code. This will instantiate a Prisma Client instance.

 ```ts title="client.ts"
 import "dotenv/config";
 import { PrismaPg } from "@prisma/adapter-pg";
 import { PrismaClient } from "../generated/prisma/client";

 const connectionString = `${process.env.DATABASE_URL}`;

 const adapter = new PrismaPg({ connectionString });
 const prisma = new PrismaClient({ adapter });

 export { prisma };
 ```

2. Next create a file named `singleton.ts` at your projects root and add the following:

 ```ts title="singleton.ts"
 import { PrismaClient } from "../generated/prisma/client";
 import { mockDeep, mockReset, DeepMockProxy } from "jest-mock-extended";

 import prisma from "./client";

 jest.mock("./client", () => ({
 __esModule: true,
 default: mockDeep<PrismaClient>(),
 }));

 beforeEach(() => {
 mockReset(prismaMock);
 });

 export const prismaMock = prisma as unknown as DeepMockProxy<PrismaClient>;
 ```

The singleton file tells Jest to mock a default export (the Prisma Client instance in `./client.ts`), and uses the `mockDeep` method from `jest-mock-extended` to enable access to the objects and methods available on Prisma Client. It then resets the mocked instance before each test is run.

Next, add the `setupFilesAfterEnv` property to your `jest.config.js` file with the path to your `singleton.ts` file.

```js title="jest.config.js" highlight=5;add showLineNumbers
module.exports = {
 clearMocks: true,
 preset: "ts-jest",
 testEnvironment: "node",
 setupFilesAfterEnv: ["<rootDir>/singleton.ts"], {/* [!code ++] */}
};
```

### Dependency injection

Another popular pattern that can be used is dependency injection.

1. Create a `context.ts` file and add the following:

 ```ts title="context.ts"
 import { PrismaClient } from "../generated/prisma/client";
 import { mockDeep, DeepMockProxy } from "jest-mock-extended";

 export type Context = {
 prisma: PrismaClient;
 };

 export type MockContext = {
 prisma: DeepMockProxy<PrismaClient>;
 };

 export const createMockContext = (): MockContext => {
 return {
 prisma: mockDeep<PrismaClient>(),
 };
 };
 ```

:::tip

If you find that you're seeing a circular dependency error highlighted through mocking Prisma Client, try adding `"strictNullChecks": true`
to your `tsconfig.json`.

:::

2. To use the context, you would do the following in your test file:

 ```ts
 import { MockContext, Context, createMockContext } from "../context";

 let mockCtx: MockContext;
 let ctx: Context;

 beforeEach(() => {
 mockCtx = createMockContext();
 ctx = mockCtx as unknown as Context;
 });
 ```

This will create a new context before each test is run via the `createMockContext` function. This (`mockCtx`) context will be used to make a mock call to Prisma Client and run a query to test. The `ctx` context will be used to run a scenario query that is tested against.

## Example unit tests

A real world use case for unit testing Prisma ORM might be a signup form. Your user fills in a form which calls a function, which in turn uses Prisma Client to make a call to your database.

All of the examples that follow use the following schema model:

```prisma title="schema.prisma" showLineNumbers
model User {
 id Int @id @default(autoincrement())
 email String @unique
 name String?
 acceptTermsAndConditions Boolean
}
```

The following unit tests will mock the process of

- Creating a new user
- Updating a users name
- Failing to create a user if terms are not accepted

The functions that use the dependency injection pattern will have the context injected (passed in as a parameter) into them, whereas the functions that use the singleton pattern will use the singleton instance of Prisma Client.

```ts title="functions-with-context.ts"
import { Context } from "./context";

interface CreateUser {
 name: string;
 email: string;
 acceptTermsAndConditions: boolean;
}

export async function createUser(user: CreateUser, ctx: Context) {
 if (user.acceptTermsAndConditions) {
 return await ctx.prisma.user.create({
 data: user,
 });
 } else {
 return new Error("User must accept terms!");
 }
}

interface UpdateUser {
 id: number;
 name: string;
 email: string;
}

export async function updateUsername(user: UpdateUser, ctx: Context) {
 return await ctx.prisma.user.update({
 where: { id: user.id },
 data: user,
 });
}
```

```ts title="functions-without-context.ts"
import prisma from "./client";

interface CreateUser {
 name: string;
 email: string;
 acceptTermsAndConditions: boolean;
}

export async function createUser(user: CreateUser) {
 if (user.acceptTermsAndConditions) {
 return await prisma.user.create({
 data: user,
 });
 } else {
 return new Error("User must accept terms!");
 }
}

interface UpdateUser {
 id: number;
 name: string;
 email: string;
}

export async function updateUsername(user: UpdateUser) {
 return await prisma.user.update({
 where: { id: user.id },
 data: user,
 });
}
```

The tests for each methodology are fairly similar, the difference is how the mocked Prisma Client is used.

The **_dependency injection_** example passes the context through to the function that is being tested as well as using it to call the mock implementation.

The **_singleton_** example uses the singleton client instance to call the mock implementation.

```ts title="__tests__/with-singleton.ts"
import { createUser, updateUsername } from "../functions-without-context";
import { prismaMock } from "../singleton";

test("should create new user ", async () => {
 const user = {
 id: 1,
 name: "Rich",
 email: "hello@prisma.io",
 acceptTermsAndConditions: true,
 };

 prismaMock.user.create.mockResolvedValue(user);

 await expect(createUser(user)).resolves.toEqual({
 id: 1,
 name: "Rich",
 email: "hello@prisma.io",
 acceptTermsAndConditions: true,
 });
});

test("should update a users name ", async () => {
 const user = {
 id: 1,
 name: "Rich Haines",
 email: "hello@prisma.io",
 acceptTermsAndConditions: true,
 };

 prismaMock.user.update.mockResolvedValue(user);

 await expect(updateUsername(user)).resolves.toEqual({
 id: 1,
 name: "Rich Haines",
 email: "hello@prisma.io",
 acceptTermsAndConditions: true,
 });
});

test("should fail if user does not accept terms", async () => {
 const user = {
 id: 1,
 name: "Rich Haines",
 email: "hello@prisma.io",
 acceptTermsAndConditions: false,
 };

 prismaMock.user.create.mockImplementation();

 await expect(createUser(user)).resolves.toEqual(new Error("User must accept terms!"));
});
```

```ts title="__tests__/with-dependency-injection.ts"
import { MockContext, Context, createMockContext } from "../context";
import { createUser, updateUsername } from "../functions-with-context";

let mockCtx: MockContext;
let ctx: Context;

beforeEach(() => {
 mockCtx = createMockContext();
 ctx = mockCtx as unknown as Context;
});

test("should create new user ", async () => {
 const user = {
 id: 1,
 name: "Rich",
 email: "hello@prisma.io",
 acceptTermsAndConditions: true,
 };
 mockCtx.prisma.user.create.mockResolvedValue(user);

 await expect(createUser(user, ctx)).resolves.toEqual({
 id: 1,
 name: "Rich",
 email: "hello@prisma.io",
 acceptTermsAndConditions: true,
 });
});

test("should update a users name ", async () => {
 const user = {
 id: 1,
 name: "Rich Haines",
 email: "hello@prisma.io",
 acceptTermsAndConditions: true,
 };
 mockCtx.prisma.user.update.mockResolvedValue(user);

 await expect(updateUsername(user, ctx)).resolves.toEqual({
 id: 1,
 name: "Rich Haines",
 email: "hello@prisma.io",
 acceptTermsAndConditions: true,
 });
});

test("should fail if user does not accept terms", async () => {
 const user = {
 id: 1,
 name: "Rich Haines",
 email: "hello@prisma.io",
 acceptTermsAndConditions: false,
 };

 mockCtx.prisma.user.create.mockImplementation();

 await expect(createUser(user, ctx)).resolves.toEqual(new Error("User must accept terms!"));
});
```

---
title: Type safety Overview
description: 'Prisma Client provides full type safety for queries, even for partial queries or included relations. This page explains how to leverage the generated types and utilities'
url: /orm/prisma-client/type-safety
metaTitle: 'Type safety'
metaDescription: 'Prisma Client provides full type safety for queries, even for partial queries or included relations. This page explains how to leverage the generated types and utilities.'

---

The generated code for Prisma Client contains several helpful types and utilities that you can use to make your application more type-safe. This page describes patterns for leveraging them.

> **Note**: If you're interested in advanced type safety topics with Prisma ORM, be sure to check out this [blog post](https://www.prisma.io/blog/satisfies-operator-ur8ys8ccq7zb) about improving your Prisma Client workflows with the new TypeScript `satisfies` keyword.

## Importing generated types

You can import the `Prisma` namespace and use dot notation to access types and utilities. The following example shows how to import the `Prisma` namespace and use it to access and use the `Prisma.UserSelect` [generated type](#what-are-generated-types):

```ts
import { Prisma } from "../path/to/generated/prisma/client";

// Build 'select' object
const userEmail: Prisma.UserSelect = {
 email: true,
};

// Use select object
const createUser = await prisma.user.create({
 data: {
 email: "bob@prisma.io",
 },
 select: userEmail,
});
```

See also: [Using the `Prisma.UserCreateInput` generated type](/orm/prisma-client/queries/crud#create-a-single-record)

## What are generated types?

Generated types are TypeScript types that are derived from your models. You can use them to create typed objects that you pass into top-level methods like `prisma.user.create(...)` or `prisma.user.update(...)`, or options such as `select` or `include`.

For example, `select` accepts an object of type `UserSelect`. Its object properties match those that are supported by `select` statements according to the model.

The first tab below shows the `UserSelect` generated type and how each property on the object has a type annotation. The second tab shows the original schema from which the type was generated.

```ts tab="Generated type"
type Prisma.UserSelect = {
 id?: boolean | undefined;
 email?: boolean | undefined;
 name?: boolean | undefined;
 posts?: boolean | Prisma.PostFindManyArgs | undefined;
 profile?: boolean | Prisma.ProfileArgs | undefined;
}
```

```prisma tab="Model"
model User {
 id Int @id @default(autoincrement())
 email String @unique
 name String?
 posts Post[]
 profile Profile?
}
```

In TypeScript the concept of [type annotations](https://www.typescriptlang.org/docs/handbook/2/everyday-types.html#type-annotations-on-variables) is when you declare a variable and add a type annotation to describe the type of the variable. See the below example.

```ts
const myAge: number = 37;
const myName: string = "Rich";
```

Both of these variable declarations have been given a type annotation to specify what primitive type they are, `number` and `string` respectively. Most of the time this kind of annotation is not needed as TypeScript will infer the type of the variable based on how its initialized. In the above example `myAge` was initialized with a number so TypeScript guesses that it should be typed as a number.

Going back to the `UserSelect` type, if you were to use dot notation on the created object `userEmail`, you would have access to all of the fields on the `User` model that can be interacted with using a `select` statement.

```prisma
model User {
 id Int @id @default(autoincrement())
 email String @unique
 name String?
 posts Post[]
 profile Profile?
}
```

```ts
import { Prisma } from "../path/to/generated/prisma/client";

const userEmail: Prisma.UserSelect = {
 email: true,
};

// properties available on the typed object
userEmail.id;
userEmail.email;
userEmail.name;
userEmail.posts;
userEmail.profile;
```

In the same mould, you can type an object with an `include` generated type then your object would have access to those properties on which you can use an `include` statement.

```ts
import { Prisma } from "../path/to/generated/prisma/client";

const userPosts: Prisma.UserInclude = {
 posts: true,
};

// properties available on the typed object
userPosts.posts;
userPosts.profile;
```

> See the [model query options](/orm/reference/prisma-client-reference#model-query-options) reference for more information about the different types available.

### Generated `UncheckedInput` types

The `UncheckedInput` types are a special set of generated types that allow you to perform some operations that Prisma Client considers "unsafe", like directly writing [relation scalar fields](/orm/prisma-schema/data-model/relations). You can choose either the "safe" `Input` types or the "unsafe" `UncheckedInput` type when doing operations like `create`, `update`, or `upsert`.

For example, this Prisma schema has a one-to-many relation between `User` and `Post`:

```prisma
model Post {
 id Int @id @default(autoincrement())
 title String @db.VarChar(255)
 content String?
 author User @relation(fields: [authorId], references: [id])
 authorId Int
}

model User {
 id Int @id @default(autoincrement())
 email String @unique
 name String?
 posts Post[]
}
```

The first tab shows the `PostUncheckedCreateInput` generated type. It contains the `authorId` property, which is a relation scalar field. The second tab shows an example query that uses the `PostUncheckedCreateInput` type. This query will result in an error if a user with an `id` of `1` does not exist.

```ts tab="Generated type"
type PostUncheckedCreateInput = {
 id?: number;
 title: string;
 content?: string | null;
 authorId: number;
};
```

```ts tab="Example query"
prisma.post.create({
 data: {
 title: "First post",
 content: "Welcome to the first post in my blog...",
 authorId: 1,
 },
});
```

The same query can be rewritten using the "safer" `PostCreateInput` type. This type does not contain the `authorId` field but instead contains the `author` relation field.

```ts tab="Generated type"
type PostCreateInput = {
 title: string;
 content?: string | null;
 author: UserCreateNestedOneWithoutPostsInput;
};

type UserCreateNestedOneWithoutPostsInput = {
 create?: XOR<UserCreateWithoutPostsInput, UserUncheckedCreateWithoutPostsInput>;
 connectOrCreate?: UserCreateOrConnectWithoutPostsInput;
 connect?: UserWhereUniqueInput;
};
```

```ts tab="Example query"
prisma.post.create({
 data: {
 title: "First post",
 content: "Welcome to the first post in my blog...",
 author: {
 connect: {
 id: 1,
 },
 },
 },
});
```

This query will also result in an error if an author with an `id` of `1` does not exist. In this case, Prisma Client will give a more descriptive error message. You can also use the [`connectOrCreate`](/orm/reference/prisma-client-reference#connectorcreate) API to safely create a new user if one does not already exist with the given `id`.

We recommend using the "safe" `Input` types whenever possible.

## Type utilities

To help you create highly type-safe applications, Prisma Client provides a set of type utilities that tap into input and output types. These types are fully dynamic, which means that they adapt to any given model and schema. You can use them to improve the auto-completion and developer experience of your projects.

This is especially useful in [shared Prisma Client extensions](/orm/prisma-client/client-extensions/shared-extensions).

The following type utilities are available in Prisma Client:

- `Exact<Input, Shape>`: Enforces strict type safety on `Input`. `Exact` makes sure that a generic type `Input` strictly complies with the type that you specify in `Shape`. It [narrows](https://www.typescriptlang.org/docs/handbook/2/narrowing.html) `Input` down to the most precise types.
- `Args<Type, Operation>`: Retrieves the input arguments for any given model and operation. This is particularly useful for extension authors who want to do the following:
 - Re-use existing types to extend or modify them.
 - Benefit from the same auto-completion experience as on existing operations.
- `Result<Type, Arguments, Operation>`: Takes the input arguments and provides the result for a given model and operation. You would usually use this in conjunction with `Args`. As with `Args`, `Result` helps you to re-use existing types to extend or modify them.
- `Payload<Type, Operation>`: Retrieves the entire structure of the result, as scalars and relations objects for a given model and operation. For example, you can use this to determine which keys are scalars or objects at a type level.

As an example, here's a quick way you can enforce that the arguments to a function matches what you will pass to a `post.create`:

```ts
type PostCreateBody = Prisma.Args<typeof prisma.post, "create">["data"];

const addPost = async (postBody: PostCreateBody) => {
 const post = await prisma.post.create({ data: postBody });
 return post;
};

await addPost(myData);
// ^ guaranteed to match the input of `post.create`
```

---
title: Operating against partial structures of your model types
description: This page documents various scenarios for using the generated types from the Prisma namespace
url: /orm/prisma-client/type-safety/operating-against-partial-structures-of-model-types
metaTitle: Operating against partial structures of your model types
metaDescription: This page documents various scenarios for using the generated types from the Prisma namespace
---

When using Prisma Client, every model from your [Prisma schema](/orm/prisma-schema/overview) is translated into a dedicated TypeScript type. For example, assume you have the following `User` and `Post` models:

```prisma
model User {
 id Int @id
 email String @unique
 name String?
 posts Post[]
}

model Post {
 id Int @id
 author User @relation(fields: [userId], references: [id])
 title String
 published Boolean @default(false)
 userId Int
}
```

The Prisma Client code that's generated from this schema contains this representation of the `User` type:

```ts
export type User = {
 id: string;
 email: string;
 name: string | null;
};
```

## Problem: Using variations of the generated model type

### Description

In some scenarios, you may need a _variation_ of the generated `User` type. For example, when you have a function that expects an instance of the `User` model that carries the `posts` relation. Or when you need a type to pass only the `User` model's `email` and `name` fields around in your application code.

### Solution

As a solution, you can customize the generated model type using Prisma Client's helper types.

The `User` type only contains the model's [scalar](/orm/prisma-schema/data-model/models#scalar-fields) fields, but doesn't account for any relations. That's because [relations are not included by default](/orm/prisma-client/queries/select-fields#return-the-default-fields) in Prisma Client queries.

However, sometimes it's useful to have a type available that **includes a relation** (i.e. a type that you'd get from an API call that uses [`include`](/orm/prisma-client/queries/select-fields#return-nested-objects-by-selecting-relation-fields)). Similarly, another useful scenario could be to have a type available that **includes only a subset of the model's scalar fields** (i.e. a type that you'd get from an API call that uses [`select`](/orm/prisma-client/queries/select-fields#select-specific-fields)).

One way of achieving this would be to define these types manually in your application code:

```ts
// 1: Define a type that includes the relation to `Post`
type UserWithPosts = {
 id: string;
 email: string;
 name: string | null;
 posts: Post[];
};

// 2: Define a type that only contains a subset of the scalar fields
type UserPersonalData = {
 email: string;
 name: string | null;
};
```

While this is certainly feasible, this approach increases the maintenance burden upon changes to the Prisma schema as you need to manually maintain the types. A cleaner solution to this is to use the `UserGetPayload` type that is generated and exposed by Prisma Client under the `Prisma` namespace in combination with TypeScript's `satisfies` operator.

The following example uses the `satisfies` operator to create two type-safe objects and then uses the `Prisma.UserGetPayload` utility function to create a type that can be used to return all users and their posts.

```ts
import { Prisma } from "@prisma/client";

// 1: Define a type that includes the relation to `Post`
const userWithPosts = { include: { posts: true } } satisfies Prisma.UserDefaultArgs;

// 2: Define a type that only contains a subset of the scalar fields
const userPersonalData = { select: { email: true, name: true } } satisfies Prisma.UserDefaultArgs;

// 3: This type will include a user and all their posts
type UserWithPosts = Prisma.UserGetPayload<typeof userWithPosts>;
```

The main benefits of the latter approach are:

- Cleaner approach as it leverages Prisma Client's generated types
- Reduced maintenance burden and improved type safety when the schema changes

## Problem: Getting access to the return type of a function

### Description

When doing [`select`](/orm/reference/prisma-client-reference#select) or [`include`](/orm/reference/prisma-client-reference#include) operations on your models and returning these variants from a function, it can be difficult to gain access to the return type, e.g:

```ts
// Function definition that returns a partial structure
async function getUsersWithPosts() {
 const users = await prisma.user.findMany({ include: { posts: true } });
 return users;
}
```

Extracting the type that represents "users with posts" from the above code snippet requires some advanced TypeScript usage:

```ts
// Function definition that returns a partial structure
async function getUsersWithPosts() {
 const users = await prisma.user.findMany({ include: { posts: true } });
 return users;
}

// Extract `UsersWithPosts` type with
type ThenArg<T> = T extends PromiseLike<infer U> ? U : T;
type UsersWithPosts = ThenArg<ReturnType<typeof getUsersWithPosts>>;

// run inside `async` function
const usersWithPosts: UsersWithPosts = await getUsersWithPosts();
```

### Solution

You can use native the TypeScript utility type [`Awaited`](https://www.typescriptlang.org/docs/handbook/utility-types.html#awaitedtype) and [`ReturnType`](https://www.typescriptlang.org/docs/handbook/utility-types.html#returntypetype) to solve the problem elegantly:

```ts
type UsersWithPosts = Awaited<ReturnType<typeof getUsersWithPosts>>;
```

---
title: How to use Prisma ORM's type system
description: How to use Prisma ORM's type system
url: /orm/prisma-client/type-safety/prisma-type-system
metaDescription: How to use Prisma ORM's type system
metaTitle: How to use Prisma ORM's type system
---

This guide introduces Prisma ORM's type system and explains how to introspect existing native types in your database, and how to use types when you apply schema changes to your database with Prisma Migrate or `db push`.

## How does Prisma ORM's type system work?

Prisma ORM uses _types_ to define the kind of data that a field can hold. To make it easy to get started, Prisma ORM provides a small number of core [scalar types](/orm/reference/prisma-schema-reference#model-field-scalar-types) that should cover most default use cases. For example, take the following blog post model:

```prisma title="schema.prisma" showLineNumbers
datasource db {
 provider = "postgresql"
}

model Post {
 id Int @id
 title String
 createdAt DateTime
}
```

The `title` field of the `Post` model uses the `String` scalar type, while the `createdAt` field uses the `DateTime` scalar type.

Databases also have their own type system, which defines the type of value that a column can hold. Most databases provide a large number of data types to allow fine-grained control over exactly what a column can store. For example, a database might provide inbuilt support for multiple sizes of integers, or for XML data. The names of these types vary between databases. For example, in PostgreSQL the column type for booleans is `boolean`, whereas in MySQL the `tinyint(1)` type is typically used.

In the blog post example above, we are using the PostgreSQL connector. This is specified in the `datasource` block of the Prisma schema.

### Default type mappings

To allow you to get started with our core scalar types, Prisma ORM provides _default type mappings_ that map each scalar type to a default type in the underlying database. For example:

- by default Prisma ORM's `String` type gets mapped to PostgreSQL's `text` type and MySQL's `varchar` type
- by default Prisma ORM's `DateTime` type gets mapped to PostgreSQL's `timestamp(3)` type and SQL Server's `datetime2` type

See Prisma ORM's [database connector pages](/orm/core-concepts/supported-databases) for the default type mappings for a given database. For example, [this table](/orm/core-concepts/supported-databases/postgresql#prisma-to-postgresql) gives the default type mappings for PostgreSQL.
To see the default type mappings for all databases for a specific given Prisma ORM type, see the [model field scalar types section](/orm/reference/prisma-schema-reference#model-field-scalar-types) of the Prisma schema reference. For example, [this table](/orm/reference/prisma-schema-reference#float) gives the default type mappings for the `Float` scalar type.

### Native type mappings

Sometimes you may need to use a more specific database type that is not one of the default type mappings for your Prisma ORM type. For this purpose, Prisma ORM provides [native type attributes](/orm/prisma-schema/data-model/models#native-types-mapping) to refine the core scalar types. For example, in the `createdAt` field of your `Post` model above you may want to use a date-only column in your underlying PostgreSQL database, by using the `date` type instead of the default type mapping of `timestamp(3)`. To do this, add a `@db.Date` native type attribute to the `createdAt` field:

```prisma title="schema.prisma" showLineNumbers
model Post {
 id Int @id
 title String
 createdAt DateTime @db.Date
}
```

Native type mappings allow you to express all the types in your database. However, you do not need to use them if the Prisma ORM defaults satisfy your needs. This leads to a shorter, more readable Prisma schema for common use cases.

## How to introspect database types

When you [introspect](/orm/prisma-schema/introspection) an existing database, Prisma ORM will take the database type of each table column and represent it in your Prisma schema using the correct Prisma ORM type for the corresponding model field. If the database type is not the default database type for that Prisma ORM scalar type, Prisma ORM will also add a native type attribute.

As an example, take a `User` table in a PostgreSQL database, with:

- an `id` column with a data type of `serial`
- a `name` column with a data type of `text`
- an `isActive` column with a data type of `boolean`

You can create this with the following SQL command:

```sql
CREATE TABLE "public"."User" (
 id serial PRIMARY KEY NOT NULL,
 name text NOT NULL,
 "isActive" boolean NOT NULL
);
```

Introspect your database with the following command run from the root directory of your project:

```npm
npx prisma db pull
```

You will get the following Prisma schema:

```prisma title="schema.prisma" showLineNumbers
model User {
 id Int @id @default(autoincrement())
 name String
 isActive Boolean
}
```

The `id`, `name` and `isActive` columns in the database are mapped respectively to the `Int`, `String` and `Boolean` Prisma ORM types. The database types are the _default_ database types for these Prisma ORM types, so Prisma ORM does not add any native type attributes.

Now add a `createdAt` column to your database with a data type of `date` by running the following SQL command:

```sql
ALTER TABLE "public"."User"
ADD COLUMN "createdAt" date NOT NULL;
```

Introspect your database again:

```npm
npx prisma db pull
```

Your Prisma schema now includes the new `createdAt` field with a Prisma ORM type of `DateTime`. The `createdAt` field also has a `@db.Date` native type attribute, because PostgreSQL's `date` is not the default type for the `DateTime` type:

```prisma title="schema.prisma" highlight=5;add showLineNumbers
model User {
 id Int @id @default(autoincrement())
 name String
 isActive Boolean
 createdAt DateTime @db.Date // [!code ++]
}
```

## How to use types when you apply schema changes to your database

When you apply schema changes to your database using Prisma Migrate or `db push`, Prisma ORM will use both the Prisma ORM scalar type of each field and any native attribute it has to determine the correct database type for the corresponding column in the database.

As an example, create a Prisma schema with the following `Post` model:

```prisma title="schema.prisma" showLineNumbers
model Post {
 id Int @id
 title String
 createdAt DateTime
 updatedAt DateTime @db.Date
}
```

This `Post` model has:

- an `id` field with a Prisma ORM type of `Int`
- a `title` field with a Prisma ORM type of `String`
- a `createdAt` field with a Prisma ORM type of `DateTime`
- an `updatedAt` field with a Prisma ORM type of `DateTime` and a `@db.Date` native type attribute

Now apply these changes to an empty PostgreSQL database with the following command, run from the root directory of your project:

```npm
npx prisma db push
```

You will see that the database has a newly created `Post` table, with:

- an `id` column with a database type of `integer`
- a `title` column with a database type of `text`
- a `createdAt` column with a database type of `timestamp(3)`
- an `updatedAt` column with a database type of `date`

Notice that the `@db.Date` native type attribute modifies the database type of the `updatedAt` column to `date`, rather than the default of `timestamp(3)`.

## More on using Prisma ORM's type system

For further reference information on using Prisma ORM's type system, see the following resources:

- The [database connector](/orm/core-concepts/supported-databases) page for each database provider has a type mapping section with a table of default type mappings between Prisma ORM types and database types, and a table of database types with their corresponding native type attribute in Prisma ORM. For example, the type mapping section for PostgreSQL is [here](/orm/core-concepts/supported-databases/postgresql#prisma-to-postgresql).
- The [model field scalar types](/orm/reference/prisma-schema-reference#model-field-scalar-types) section of the Prisma schema reference has a subsection for each Prisma ORM scalar type. This includes a table of default mappings for that Prisma ORM type in each database, and a table for each database listing the corresponding database types and their native type attributes in Prisma ORM. For example, the entry for the `String` Prisma ORM type is [here](/orm/reference/prisma-schema-reference#string).

---
title: Write your own SQL
description: Learn how to use raw SQL queries in Prisma Client
url: /orm/prisma-client/using-raw-sql
metaTitle: Write Your Own SQL in Prisma Client
metaDescription: Learn how to use raw SQL queries in Prisma Client.
---

While the Prisma Client API aims to make all your database queries intuitive, type-safe, and convenient, there may still be situations where raw SQL is the best tool for the job.

This can happen for various reasons, such as the need to optimize the performance of a specific query or because your data requirements can't be fully expressed by Prisma Client's query API.

In most cases, [TypedSQL](#writing-type-safe-queries-with-prisma-client-and-typedsql) allows you to express your query in SQL while still benefiting from Prisma Client's excellent user experience. However, since TypedSQL is statically typed, it may not handle certain scenarios, such as dynamically generated `WHERE` clauses. In these cases, you will need to use [`$queryRaw`](/orm/prisma-client/using-raw-sql/raw-queries#queryraw) or [`$executeRaw`](/orm/prisma-client/using-raw-sql/raw-queries#executeraw), or their unsafe counterparts.

## Writing type-safe queries with Prisma Client and TypedSQL

### What is TypedSQL?

TypedSQL is a new feature of Prisma ORM that allows you to write your queries in `.sql` files while still enjoying the great developer experience of Prisma Client. You can write the code you're comfortable with and benefit from fully-typed inputs and outputs.

With TypedSQL, you can:

1. Write complex SQL queries using familiar syntax
2. Benefit from full IDE support and syntax highlighting for SQL
3. Import your SQL queries as fully typed functions in your TypeScript code
4. Maintain the flexibility of raw SQL with the safety of Prisma's type system

TypedSQL is particularly useful for:

- Complex reporting queries that are difficult to express using Prisma's query API
- Performance-critical operations that require fine-tuned SQL
- Leveraging database-specific features not yet supported in Prisma's API

By using TypedSQL, you can write efficient, type-safe database queries without sacrificing the power and flexibility of raw SQL. This feature allows you to seamlessly integrate custom SQL queries into your Prisma-powered applications, ensuring type safety and improving developer productivity.

For a detailed guide on how to get started with TypedSQL, including setup instructions and usage examples, please refer to our [TypedSQL documentation](/orm/prisma-client/using-raw-sql/typedsql).

## Raw queries

While not as ergonomic as [TypedSQL](#writing-type-safe-queries-with-prisma-client-and-typedsql), raw queries are still supported and useful when TypedSQL queries are not possible due to features not yet supported in TypedSQL or when the query is dynamically generated.

### Alternative approaches to raw SQL queries in relational databases

Prisma ORM supports four methods to execute raw SQL queries in relational databases:

- [`$queryRaw`](/orm/prisma-client/using-raw-sql/raw-queries#queryraw)
- [`$executeRaw`](/orm/prisma-client/using-raw-sql/raw-queries#executeraw)
- [`$queryRawUnsafe`](/orm/prisma-client/using-raw-sql/raw-queries#queryrawunsafe)
- [`$executeRawUnsafe`](/orm/prisma-client/using-raw-sql/raw-queries#executerawunsafe)

These commands are similar to using TypedSQL, but they are not type-safe and are written as strings in your code rather than in dedicated `.sql` files.

### Alternative approaches to raw queries in document databases

For MongoDB, Prisma ORM supports three methods to execute raw queries:

- [`$runCommandRaw`](/orm/prisma-client/using-raw-sql/raw-queries#runcommandraw)
- [`<model>.findRaw`](/orm/prisma-client/using-raw-sql/raw-queries#findraw)
- [`<model>.aggregateRaw`](/orm/prisma-client/using-raw-sql/raw-queries#aggregateraw)

These methods allow you to execute raw MongoDB commands and queries, providing flexibility when you need to use MongoDB-specific features or optimizations.

`$runCommandRaw` is used to execute database commands, `<model>.findRaw` is used to find documents that match a filter, and `<model>.aggregateRaw` is used for aggregation operations.

Similar to raw queries in relational databases, these methods are not type-safe and require manual handling of the query results.

---
title: Raw queries
description: Learn how you can send raw SQL and MongoDB queries to your database using the raw() methods from the Prisma Client API
url: /orm/prisma-client/using-raw-sql/raw-queries
metaTitle: Raw queries
metaDescription: Learn how you can send raw SQL and MongoDB queries to your database using the raw() methods from the Prisma Client API.
---

:::warning

We recommend using [TypedSQL](/orm/prisma-client/using-raw-sql) for type-safe SQL queries instead of the raw queries described below.

:::

Prisma Client supports sending raw queries to your database. You may wish to use raw queries if:

- you want to run a heavily optimized query
- you require a feature that Prisma Client does not yet support (please [consider raising an issue](https://github.com/prisma/prisma/issues/new/choose))

Raw queries are available for all relational databases Prisma ORM supports, as well as MongoDB. For more details, see the relevant sections:

- [Raw queries with relational databases](#raw-queries-with-relational-databases)
- [Raw queries with MongoDB](#raw-queries-with-mongodb)

## Raw queries with relational databases

For relational databases, Prisma Client exposes four methods that allow you to send raw queries. You can use:

- `$queryRaw` to return actual records (for example, using `SELECT`).
- `$executeRaw` to return a count of affected rows (for example, after an `UPDATE` or `DELETE`).
- `$queryRawUnsafe` to return actual records (for example, using `SELECT`) using a raw string.
- `$executeRawUnsafe` to return a count of affected rows (for example, after an `UPDATE` or `DELETE`) using a raw string.

The methods with "Unsafe" in the name are a lot more flexible but are at **significant risk of making your code vulnerable to SQL injection**.

The other two methods are safe to use with a simple template tag, no string building, and no concatenation. **However**, caution is required for more complex use cases as it is still possible to introduce SQL injection if these methods are used in certain ways. For more details, see the [SQL injection prevention](#sql-injection-prevention) section below.

> **Note**: All methods in the above list can only run **one** query at a time. You cannot append a second query - for example, calling any of them with `select 1; select 2;` will not work.

### `$queryRaw`

`$queryRaw` returns actual database records. For example, the following `SELECT` query returns all fields for each record in the `User` table:

```ts no-lines
const result = await prisma.$queryRaw`SELECT * FROM User`;
```

The method is implemented as a [tagged template](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Template_literals#tagged_templates), which allows you to pass a template literal where you can easily insert your [variables](#using-variables). In turn, Prisma Client creates prepared statements that are safe from SQL injections:

```ts no-lines
const email = "emelie@prisma.io";
const result = await prisma.$queryRaw`SELECT * FROM User WHERE email = ${email}`;
```

You can also use the [`Prisma.sql`](#tagged-template-helpers) helper, in fact, the `$queryRaw` method will **only accept** a template string or the `Prisma.sql` helper:

```ts no-lines
const email = "emelie@prisma.io";
const result = await prisma.$queryRaw(Prisma.sql`SELECT * FROM User WHERE email = ${email}`);
```

:::warning

If you use string building to incorporate untrusted input into queries passed to this method, then you open up the possibility for SQL injection attacks. SQL injection attacks can expose your data to modification or deletion. The preferred mechanism would be to include the text of the query at the point that you run this method. For more information on this risk and also examples of how to prevent it, see the [SQL injection prevention](#sql-injection-prevention) section below.

:::

#### Considerations

Be aware that:

- Template variables cannot be used inside SQL string literals. For example, the following query would **not** work:

 ```ts no-lines
 const name = "Bob";
 await prisma.$queryRaw`SELECT 'My name is ${name}';`;
 ```

 Instead, you can either pass the whole string as a variable, or use string concatenation:

 ```ts no-lines
 const name = "My name is Bob";
 await prisma.$queryRaw`SELECT ${name};`;
 ```

 ```ts no-lines
 const name = "Bob";
 await prisma.$queryRaw`SELECT 'My name is ' || ${name};`;
 ```

- Template variables can only be used for data values (such as `email` in the example above). Variables cannot be used for identifiers such as column names, table names or database names, or for SQL keywords. For example, the following two queries would **not** work:

 ```ts no-lines
 const myTable = "user";
 await prisma.$queryRaw`SELECT * FROM ${myTable};`;
 ```

 ```ts no-lines
 const ordering = "desc";
 await prisma.$queryRaw`SELECT * FROM Table ORDER BY ${ordering};`;
 ```

- Prisma maps any database values returned by `$queryRaw` and `$queryRawUnsafe` to their corresponding JavaScript types. [Learn more](#raw-query-type-mapping).

- `$queryRaw` does not support dynamic table names in PostgreSQL databases. [Learn more](#dynamic-table-names-in-postgresql)

#### Return type

`$queryRaw` returns an array. Each object corresponds to a database record:

```json5
[
 { id: 1, email: "emelie@prisma.io", name: "Emelie" },
 { id: 2, email: "yin@prisma.io", name: "Yin" },
]
```

You can also [type the results of `$queryRaw`](#typing-queryraw-results).

#### Signature

```ts no-lines
$queryRaw<T = unknown>(query: TemplateStringsArray | Prisma.Sql, ...values: any[]): PrismaPromise<T>;
```

#### Typing `$queryRaw` results

`PrismaPromise<T>` uses a [generic type parameter `T`](https://www.typescriptlang.org/docs/handbook/generics.html). You can determine the type of `T` when you invoke the `$queryRaw` method. In the following example, `$queryRaw` returns `User[]`:

```ts
// import the generated `User` type from the `@prisma/client` module
import { User } from "@prisma/client";

const result = await prisma.$queryRaw<User[]>`SELECT * FROM User`;
// result is of type: `User[]`
```

> **Note**: If you do not provide a type, `$queryRaw` defaults to `unknown`.

If you are selecting **specific fields** of the model or want to include relations, refer to the documentation about [leveraging Prisma Client's generated types](/orm/prisma-client/type-safety/operating-against-partial-structures-of-model-types#problem-using-variations-of-the-generated-model-type) if you want to make sure that the results are properly typed.

#### Type caveats when using raw SQL

When you type the results of `$queryRaw`, the raw data might not always match the suggested TypeScript type. For example, the following Prisma model includes a `Boolean` field named `published`:

```prisma highlight=3;normal
model Post {
 id Int @id @default(autoincrement())
 published Boolean @default(false) // [!code highlight]
 title String
 content String?
}
```

The following query returns all posts. It then prints out the value of the `published` field for each `Post`:

```ts
const result = await prisma.$queryRaw<Post[]>`SELECT * FROM Post`;

result.forEach((x) => {
 console.log(x.published);
});
```

For regular CRUD queries, the Prisma Client query engine standardizes the return type for all databases. **Using the raw queries does not**. If the database provider is MySQL, the returned values are `1` or `0`. However, if the database provider is PostgreSQL, the values are `true` or `false`.

> **Note**: Prisma sends JavaScript integers to PostgreSQL as `INT8`. This might conflict with your user-defined functions that accept only `INT4` as input. If you use `$queryRaw` in conjunction with a PostgreSQL database, update the input types to `INT8`, or cast your query parameters to `INT4`.

#### Dynamic table names in PostgreSQL

[It is not possible to interpolate table names](#considerations). This means that you cannot use dynamic table names with `$queryRaw`. Instead, you must use [`$queryRawUnsafe`](#queryrawunsafe), as follows:

```ts
let userTable = "User";
let result = await prisma.$queryRawUnsafe(`SELECT * FROM ${userTable}`);
```

Note that if you use `$queryRawUnsafe` in conjunction with user inputs, you risk SQL injection attacks. [Learn more](#queryrawunsafe).

### `$queryRawUnsafe()`

The `$queryRawUnsafe()` method allows you to pass a raw string (or template string) to the database.

:::warning

If you use this method with user inputs (in other words, `SELECT * FROM table WHERE columnName = ${userInput}`), then you open up the possibility for SQL injection attacks. SQL injection attacks can expose your data to modification or deletion.<br /><br />

Wherever possible you should use the `$queryRaw` method instead. When used correctly `$queryRaw` method is significantly safer but note that the `$queryRaw` method can also be made vulnerable in certain circumstances. For more information, see the [SQL injection prevention](#sql-injection-prevention) section below.

:::

The following query returns all fields for each record in the `User` table:

```ts
// import the generated `User` type from the `@prisma/client` module
import { User } from "@prisma/client";

const result = await prisma.$queryRawUnsafe("SELECT * FROM User");
```

You can also run a parameterized query. The following example returns all users whose email contains the string `emelie@prisma.io`:

```ts
prisma.$queryRawUnsafe("SELECT * FROM users WHERE email = $1", "emelie@prisma.io");
```

> **Note**: Prisma sends JavaScript integers to PostgreSQL as `INT8`. This might conflict with your user-defined functions that accept only `INT4` as input. If you use a parameterized `$queryRawUnsafe` query in conjunction with a PostgreSQL database, update the input types to `INT8`, or cast your query parameters to `INT4`.

For more details on using parameterized queries, see the [parameterized queries](#parameterized-queries) section below.

#### Signature

```ts no-lines
$queryRawUnsafe<T = unknown>(query: string, ...values: any[]): PrismaPromise<T>;
```

### `$executeRaw`

`$executeRaw` returns the _number of rows affected by a database operation_, such as `UPDATE` or `DELETE`. This function does **not** return database records. The following query updates records in the database and returns a count of the number of records that were updated:

```ts
const result: number =
 await prisma.$executeRaw`UPDATE User SET active = true WHERE emailValidated = true`;
```

The method is implemented as a [tagged template](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Template_literals#tagged_templates), which allows you to pass a template literal where you can easily insert your [variables](#using-variables). In turn, Prisma Client creates prepared statements that are safe from SQL injections:

```ts
const emailValidated = true;
const active = true;

const result: number =
 await prisma.$executeRaw`UPDATE User SET active = ${active} WHERE emailValidated = ${emailValidated};`;
```

:::warning

If you use string building to incorporate untrusted input into queries passed to this method, then you open up the possibility for SQL injection attacks. SQL injection attacks can expose your data to modification or deletion. The preferred mechanism would be to include the text of the query at the point that you run this method. For more information on this risk and also examples of how to prevent it, see the [SQL injection prevention](#sql-injection-prevention) section below.

:::

#### Considerations

Be aware that:

- `$executeRaw` does not support multiple queries in a single string (for example, `ALTER TABLE` and `CREATE TABLE` together).
- Prisma Client submits prepared statements, and prepared statements only allow a subset of SQL statements. For example, `START TRANSACTION` is not permitted. You can learn more about [the syntax that MySQL allows in Prepared Statements here](https://dev.mysql.com/doc/refman/8.0/en/sql-prepared-statements.html).
- [`PREPARE` does not support `ALTER`](https://www.postgresql.org/docs/current/sql-prepare.html) - see the [workaround](#alter-limitation-postgresql).
- Template variables cannot be used inside SQL string literals. For example, the following query would **not** work:

 ```ts no-lines
 const name = "Bob";
 await prisma.$executeRaw`UPDATE user SET greeting = 'My name is ${name}';`;
 ```

 Instead, you can either pass the whole string as a variable, or use string concatenation:

 ```ts no-lines
 const name = "My name is Bob";
 await prisma.$executeRaw`UPDATE user SET greeting = ${name};`;
 ```

 ```ts no-lines
 const name = "Bob";
 await prisma.$executeRaw`UPDATE user SET greeting = 'My name is ' || ${name};`;
 ```

- Template variables can only be used for data values (such as `email` in the example above). Variables cannot be used for identifiers such as column names, table names or database names, or for SQL keywords. For example, the following two queries would **not** work:

 ```ts no-lines
 const myTable = "user";
 await prisma.$executeRaw`UPDATE ${myTable} SET active = true;`;
 ```

 ```ts no-lines
 const ordering = "desc";
 await prisma.$executeRaw`UPDATE User SET active = true ORDER BY ${desc};`;
 ```

#### Return type

`$executeRaw` returns a `number`.

#### Signature

```ts
$executeRaw<T = unknown>(query: TemplateStringsArray | Prisma.Sql, ...values: any[]): PrismaPromise<number>;
```

### `$executeRawUnsafe()`

The `$executeRawUnsafe()` method allows you to pass a raw string (or template string) to the database. Like `$executeRaw`, it does **not** return database records, but returns the number of rows affected.

:::warning

If you use this method with user inputs (in other words, `SELECT * FROM table WHERE columnName = ${userInput}`), then you open up the possibility for SQL injection attacks. SQL injection attacks can expose your data to modification or deletion.<br /><br />

Wherever possible you should use the `$executeRaw` method instead. When used correctly `$executeRaw` method is significantly safer but note that the `$executeRaw` method can also be made vulnerable in certain circumstances. For more information, see the [SQL injection prevention](#sql-injection-prevention) section below.

:::

The following example uses a template string to update records in the database. It then returns a count of the number of records that were updated:

```ts
const emailValidated = true;
const active = true;

const result = await prisma.$executeRawUnsafe(
 `UPDATE User SET active = ${active} WHERE emailValidated = ${emailValidated}`,
);
```

The same can be written as a parameterized query:

```ts
const result = prisma.$executeRawUnsafe(
 "UPDATE User SET active = $1 WHERE emailValidated = $2",
 "yin@prisma.io",
 true,
);
```

For more details on using parameterized queries, see the [parameterized queries](#parameterized-queries) section below.

#### Signature

```ts no-lines
$executeRawUnsafe<T = unknown>(query: string, ...values: any[]): PrismaPromise<number>;
```

### Raw query type mapping

Prisma maps any database values returned by `$queryRaw` and `$queryRawUnsafe`to their corresponding [JavaScript types](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Data_structures). This behavior is the same as for regular Prisma query methods like `findMany()`.

As an example, take a raw query that selects columns with `BigInt`, `Bytes`, `Decimal` and `Date` types from a table:

```ts
const result = await prisma.$queryRaw`SELECT bigint, bytes, decimal, date FROM "Table";`;

console.log(result);
```

```bash
{ bigint: BigInt("123"), bytes: <Buffer 01 02>), decimal: Decimal("12.34"), date: Date("<some_date>") }
```

In the `result` object, the database values have been mapped to the corresponding JavaScript types.

The following table shows the conversion between types used in the database and the JavaScript type returned by the raw query:

| Database type | JavaScript type |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| Text | `String` |
| 32-bit integer | `Number` |
| 32-bit unsigned integer | `BigInt` |
| Floating point number | `Number` |
| Double precision number | `Number` |
| 64-bit integer | `BigInt` |
| Decimal / numeric | `Decimal` |
| Bytes | `Uint8Array` |
| Json | `Object` |
| DateTime | `Date` |
| Date | `Date` |
| Time | `Date` |
| Uuid | `String` |
| Xml | `String` |

Note that the exact name for each database type will vary between databases – for example, the boolean type is known as `boolean` in PostgreSQL and `STRING` in CockroachDB. See the [Scalar types reference](/orm/reference/prisma-schema-reference#model-field-scalar-types) for full details of type names for each database.

### Raw query typecasting behavior

Raw queries with Prisma Client might require parameters to be in the expected types of the SQL function or query. Prisma Client does not do subtle, implicit casts.

As an example, take the following query using PostgreSQL's `LENGTH` function, which only accepts the `text` type as an input:

```ts
await prisma.$queryRaw`SELECT LENGTH(${42});`;
```

This query returns an error:

```bash wrap
// ERROR: function length(integer) does not exist
// HINT: No function matches the given name and argument types. You might need to add explicit type casts.
```

The solution in this case is to explicitly cast `42` to the `text` type:

```ts
await prisma.$queryRaw`SELECT LENGTH(${42}::text);`;
```

### Transactions

You can use `.$executeRaw()` and `.$queryRaw()` inside a [transaction](/orm/prisma-client/queries/transactions).

### Using variables

`$executeRaw` and `$queryRaw` are implemented as [**tagged templates**](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Template_literals#tagged_templates). Tagged templates are the recommended way to use variables with raw SQL in the Prisma Client.

The following example includes a placeholder named `${userId}`:

```ts
const userId = 42;
const result = await prisma.$queryRaw`SELECT * FROM User WHERE id = ${userId};`;
```

✔ Benefits of using the tagged template versions of `$queryRaw` and `$executeRaw` include:

- Prisma Client escapes all variables.
- Tagged templates are database-agnostic - you do not need to remember if variables should be written as `$1` (PostgreSQL) or `?` (MySQL).
- [SQL Template Tag](https://github.com/blakeembrey/sql-template-tag) give you access to [useful helpers](#tagged-template-helpers).
- Embedded, named variables are easier to read.

> **Note**: You cannot pass a table or column name into a tagged template placeholder. For example, you cannot `SELECT ?` and pass in `*` or `id, name` based on some condition.

#### Tagged template helpers

Prisma Client specifically uses [SQL Template Tag](https://github.com/blakeembrey/sql-template-tag), which exposes a number of helpers. For example, the following query uses `join()` to pass in a list of IDs:

```ts
import { Prisma } from "@prisma/client";

const ids = [1, 3, 5, 10, 20];
const result = await prisma.$queryRaw`SELECT * FROM User WHERE id IN (${Prisma.join(ids)})`;
```

The following example uses the `empty` and `sql` helpers to change the query depending on whether `userName` is empty:

```ts
import { Prisma } from "@prisma/client";

const userName = "";
const result = await prisma.$queryRaw`SELECT * FROM User ${
 userName ? Prisma.sql`WHERE name = ${userName}` : Prisma.empty // Cannot use "" or NULL here!
}`;
```

#### `ALTER` limitation (PostgreSQL)

PostgreSQL [does not support using `ALTER` in a prepared statement](https://www.postgresql.org/docs/current/sql-prepare.html), which means that the following queries **will not work**:

```ts
await prisma.$executeRaw`ALTER USER prisma WITH PASSWORD "${password}"`;
await prisma.$executeRaw(Prisma.sql`ALTER USER prisma WITH PASSWORD "${password}"`);
```

You can use the following query, but be aware that this is potentially **unsafe** as `${password}` is not escaped:

```ts
await prisma.$executeRawUnsafe('ALTER USER prisma WITH PASSWORD "$1"', password})
```

### Unsupported types

[`Unsupported` types](/orm/reference/prisma-schema-reference#unsupported) need to be cast to Prisma Client supported types before using them in `$queryRaw` or `$queryRawUnsafe`. For example, take the following model, which has a `location` field with an `Unsupported` type:

```tsx
model Country {
 location Unsupported("point")?
}
```

The following query on the unsupported field will **not** work:

```tsx
await prisma.$queryRaw`SELECT location FROM Country;`;
```

Instead, cast `Unsupported` fields to any supported Prisma Client type, **if your `Unsupported` column supports the cast**.

The most common type you may want to cast your `Unsupported` column to is `String`. For example, on PostgreSQL, this would map to the `text` type:

```tsx
await prisma.$queryRaw`SELECT location::text FROM Country;`;
```

The database will thus provide a `String` representation of your data which Prisma Client supports.

For details of supported Prisma types, see the [Prisma connector overview](/orm/core-concepts/supported-databases) for the relevant database.

## SQL injection prevention

The ideal way to avoid SQL injection in Prisma Client is to use the ORM models to perform queries wherever possible.

Where this is not possible and raw queries are required, Prisma Client provides various raw methods, but it is important to use these methods safely.

This section will provide various examples of using these methods safely and unsafely. You can test these examples in the [Prisma Playground](https://playground.prisma.io/examples).

### In `$queryRaw` and `$executeRaw`

#### Simple, safe use of `$queryRaw` and `$executeRaw`

These methods can mitigate the risk of SQL injection by escaping all variables when you use tagged templates and sends all queries as prepared statements.

```ts
$queryRaw`...`; // Tagged template
$executeRaw`...`; // Tagged template
```

The following example is safe ✅ from SQL Injection:

```ts
const inputString = `'Sarah' UNION SELECT id, title FROM "Post"`;
const result = await prisma.$queryRaw`SELECT id, name FROM "User" WHERE name = ${inputString}`;

console.log(result);
```

#### Unsafe use of `$queryRaw` and `$executeRaw`

However, it is also possible to use these methods in unsafe ways.

One way is by artificially generating a tagged template that unsafely concatenates user input.

The following example is vulnerable ❌ to SQL Injection:

```ts
// Unsafely generate query text
const inputString = `'Sarah' UNION SELECT id, title FROM "Post"`; // SQL Injection
const query = `SELECT id, name FROM "User" WHERE name = ${inputString}`;

// Version for Typescript
const stringsArray: any = [...[query]];

// Version for Javascript
const stringsArray = [...[query]];

// Use the `raw` property to impersonate a tagged template
stringsArray.raw = [query];

// Use queryRaw
const result = await prisma.$queryRaw(stringsArray);
console.log(result);
```

Another way to make these methods vulnerable is misuse of the `Prisma.raw` function.

The following examples are all vulnerable ❌ to SQL Injection:

```ts
const inputString = `'Sarah' UNION SELECT id, title FROM "Post"`;
const result = await prisma.$queryRaw`SELECT id, name FROM "User" WHERE name = ${Prisma.raw(
 inputString,
)}`;
console.log(result);
```

```ts
const inputString = `'Sarah' UNION SELECT id, title FROM "Post"`;
const result = await prisma.$queryRaw(
 Prisma.raw(`SELECT id, name FROM "User" WHERE name = ${inputString}`),
);
console.log(result);
```

```ts
const inputString = `'Sarah' UNION SELECT id, title FROM "Post"`;
const query = Prisma.raw(`SELECT id, name FROM "User" WHERE name = ${inputString}`);
const result = await prisma.$queryRaw(query);
console.log(result);
```

#### Safely using `$queryRaw` and `$executeRaw` in more complex scenarios

##### Building raw queries separate to query execution

If you want to build your raw queries elsewhere or separate to your parameters you will need to use one of the following methods.

In this example, the `sql` helper method is used to build the query text by safely including the variable. It is safe ✅ from SQL Injection:

```ts
// inputString can be untrusted input
const inputString = `'Sarah' UNION SELECT id, title FROM "Post"`;

// Safe if the text query below is completely trusted content
const query = Prisma.sql`SELECT id, name FROM "User" WHERE name = ${inputString}`;

const result = await prisma.$queryRaw(query);
console.log(result);
```

In this example which is safe ✅ from SQL Injection, the `sql` helper method is used to build the query text including a parameter marker for the input value. Each variable is represented by a marker symbol (`?` for MySQL, `$1`, `$2`, and so on for PostgreSQL). Note that the examples just show PostgreSQL queries.

```ts
// Version for Typescript
const query: any;

// Version for Javascript
const query;

// Safe if the text query below is completely trusted content
query = Prisma.sql`SELECT id, name FROM "User" WHERE name = $1`;

// inputString can be untrusted input
const inputString = `'Sarah' UNION SELECT id, title FROM "Post"`;
query.values = [inputString];

const result = await prisma.$queryRaw(query);
console.log(result);
```

> **Note**: PostgreSQL variables are represented by `$1`, etc

##### Building raw queries elsewhere or in stages

If you want to build your raw queries somewhere other than where the query is executed, the ideal way to do this is to create an `Sql` object from the segments of your query and pass it the parameter value.

In the following example we have two variables to parameterize. The example is safe ✅ from SQL Injection as long as the query strings being passed to `Prisma.sql` only contain trusted content:

```ts
// Example is safe if the text query below is completely trusted content
const query1 = `SELECT id, name FROM "User" WHERE name = `; // The first parameter would be inserted after this string
const query2 = ` OR name = `; // The second parameter would be inserted after this string

const inputString1 = "Fred";
const inputString2 = `'Sarah' UNION SELECT id, title FROM "Post"`;

const query = Prisma.sql([query1, query2, ""], inputString1, inputString2);
const result = await prisma.$queryRaw(query);
console.log(result);
```

> Note: Notice that the string array being passed as the first parameter `Prisma.sql` needs to have an empty string at the end as the `sql` function expects one more query segment than the number of parameters.

If you want to build your raw queries into one large string, this is still possible but requires some care as it is uses the potentially dangerous `Prisma.raw` method. You also need to build your query using the correct parameter markers for your database as Prisma won't be able to provide markers for the relevant database as it usually is.

The following example is safe ✅ from SQL Injection as long as the query strings being passed to `Prisma.raw` only contain trusted content:

```ts
// Version for Typescript
const query: any;

// Version for Javascript
const query;

// Example is safe if the text query below is completely trusted content
const query1 = `SELECT id, name FROM "User" `;
const query2 = `WHERE name = $1 `;

query = Prisma.raw(`${query1}${query2}`);

// inputString can be untrusted input
const inputString = `'Sarah' UNION SELECT id, title FROM "Post"`;
query.values = [inputString];

const result = await prisma.$queryRaw(query);
console.log(result);
```

### In `$queryRawUnsafe` and `$executeRawUnsafe`

#### Using `$queryRawUnsafe` and `$executeRawUnsafe` unsafely

If you cannot use tagged templates, you can instead use [`$queryRawUnsafe`](/orm/prisma-client/using-raw-sql/raw-queries#queryrawunsafe) or [`$executeRawUnsafe`](/orm/prisma-client/using-raw-sql/raw-queries#executerawunsafe). However, **be aware that these functions significantly increase the risk of SQL injection vulnerabilities in your code**.

The following example concatenates `query` and `inputString`. Prisma Client ❌ **cannot** escape `inputString` in this example, which makes it vulnerable to SQL injection:

```ts
const inputString = '"Sarah" UNION SELECT id, title, content FROM Post'; // SQL Injection
const query = "SELECT id, name, email FROM User WHERE name = " + inputString;
const result = await prisma.$queryRawUnsafe(query);

console.log(result);
```

#### Parameterized queries

As an alternative to tagged templates, `$queryRawUnsafe` supports standard parameterized queries where each variable is represented by a symbol (`?` for MySQL, `$1`, `$2`, and so on for PostgreSQL). Note that the examples just show PostgreSQL queries.

The following example is safe ✅ from SQL Injection:

```ts
const userName = "Sarah";
const email = "sarah@prisma.io";
const result = await prisma.$queryRawUnsafe(
 "SELECT * FROM User WHERE (name = $1 OR email = $2)",
 userName,
 email,
);
```

> **Note**: PostgreSQL variables are represented by `$1` and `$2`

As with tagged templates, Prisma Client escapes all variables when they are provided in this way.

> **Note**: You cannot pass a table or column name as a variable into a parameterized query. For example, you cannot `SELECT ?` and pass in `*` or `id, name` based on some condition.

##### Parameterized PostgreSQL `ILIKE` query

When you use `ILIKE`, the `%` wildcard character(s) should be included in the variable itself, not the query (`string`). This example is safe ✅ from SQL Injection.

```ts
const userName = "Sarah";
const emailFragment = "prisma.io";
const result = await prisma.$queryRawUnsafe(
 'SELECT * FROM "User" WHERE (name = $1 OR email ILIKE $2)',
 userName,
 `%${emailFragment}`,
);
```

> **Note**: Using `%$2` as an argument would not work

## Raw queries with MongoDB

For MongoDB, Prisma Client exposes three methods that allow you to send raw queries. You can use:

- `$runCommandRaw` to run a command against the database
- `<model>.findRaw` to find zero or more documents that match the filter.
- `<model>.aggregateRaw` to perform aggregation operations on a collection.

### `$runCommandRaw()`

`$runCommandRaw()` runs a raw MongoDB command against the database. As input, it accepts all [MongoDB database commands](https://www.mongodb.com/docs/manual/reference/command/), with the following exceptions:

- `find` (use [`findRaw()`](#findraw) instead)
- `aggregate` (use [`aggregateRaw()`](#aggregateraw) instead)

When you use `$runCommandRaw()` to run a MongoDB database command, note the following:

- The object that you pass when you invoke `$runCommandRaw()` must follow the syntax of the MongoDB database command.
- You must connect to the database with an appropriate role for the MongoDB database command.

In the following example, a query inserts two records with the same `_id`. This bypasses normal document validation.

```ts no-lines
prisma.$runCommandRaw({
 insert: "Pets",
 bypassDocumentValidation: true,
 documents: [
 {
 _id: 1,
 name: "Felinecitas",
 type: "Cat",
 breed: "Russian Blue",
 age: 12,
 },
 {
 _id: 1,
 name: "Nao Nao",
 type: "Dog",
 breed: "Chow Chow",
 age: 2,
 },
 ],
});
```

:::warning

Do not use `$runCommandRaw()` for queries which contain the `"find"` or `"aggregate"` commands, because you might be unable to fetch all data. This is because MongoDB returns a [cursor](https://www.mongodb.com/docs/manual/tutorial/iterate-a-cursor/) that is attached to your MongoDB session, and you might not hit the same MongoDB session every time. For these queries, you should use the specialised [`findRaw()`](#findraw) and [`aggregateRaw()`](#aggregateraw) methods instead.

:::

#### Return type

`$runCommandRaw()` returns a `JSON` object whose shape depends on the inputs.

#### Signature

```ts no-lines
$runCommandRaw(command: InputJsonObject): PrismaPromise<JsonObject>;
```

### `findRaw()`

`<model>.findRaw()` returns actual database records. It will find zero or more documents that match the filter on the `User` collection:

```ts no-lines
const result = await prisma.user.findRaw({
 filter: { age: { $gt: 25 } },
 options: { projection: { _id: false } },
});
```

#### Return type

`<model>.findRaw()` returns a `JSON` object whose shape depends on the inputs.

#### Signature

```ts no-lines
<model>.findRaw(args?: {filter?: InputJsonObject, options?: InputJsonObject}): PrismaPromise<JsonObject>;
```

- `filter`: The query predicate filter. If unspecified, then all documents in the collection will match the [predicate](https://www.mongodb.com/docs/manual/reference/mql/query-predicates/).
- `options`: Additional options to pass to the [`find` command](https://www.mongodb.com/docs/manual/reference/command/find/#command-fields).

### `aggregateRaw()`

`<model>.aggregateRaw()` returns aggregated database records. It will perform aggregation operations on the `User` collection:

```ts no-lines
const result = await prisma.user.aggregateRaw({
 pipeline: [
 { $match: { status: "registered" } },
 { $group: { _id: "$country", total: { $sum: 1 } } },
 ],
});
```

#### Return type

`<model>.aggregateRaw()` returns a `JSON` object whose shape depends on the inputs.

#### Signature

```ts no-lines
<model>.aggregateRaw(args?: {pipeline?: InputJsonObject[], options?: InputJsonObject}): PrismaPromise<JsonObject>;
```

- `pipeline`: An array of aggregation stages to process and transform the document stream via the [aggregation pipeline](https://www.mongodb.com/docs/atlas/data-federation/supported-unsupported/supported-aggregation/).
- `options`: Additional options to pass to the [`aggregate` command](https://www.mongodb.com/docs/manual/reference/command/aggregate/#command-fields).

#### Caveats

When working with custom objects like `ObjectId` or `Date,` you will have to pass them according to the [MongoDB extended JSON Spec](https://www.mongodb.com/docs/manual/reference/mongodb-extended-json/#type-representations).
Example:

```ts no-lines
const result = await prisma.user.aggregateRaw({
 pipeline: [
 { $match: { _id: { $oid: id } } },
 // ^ notice the $oid convention here
 ],
});
```

---
title: SafeQL & Prisma Client
description: 'Learn how to use SafeQL and Prisma Client extensions to work around features not natively supported by Prisma, such as PostGIS'
url: /orm/prisma-client/using-raw-sql/safeql
metaTitle: Integrate SafeQL with Prisma Client
metaDescription: 'Learn how to use SafeQL and Prisma Client extensions to work around features not natively supported by Prisma, such as PostGIS.'
---

## Overview

This page explains how to improve the experience of writing raw SQL in Prisma ORM. It uses [Prisma Client extensions](/orm/prisma-client/client-extensions) and [SafeQL](https://safeql.dev) to create custom, type-safe Prisma Client queries which abstract custom SQL that your app might need (using `$queryRaw`).

The example will be using [PostGIS](https://postgis.net/) and PostgreSQL, but is applicable to any raw SQL queries that you might need in your application.

:::note

This page builds on the [legacy raw query methods](/orm/prisma-client/using-raw-sql/raw-queries) available in Prisma Client. While many use cases for raw SQL in Prisma Client are covered by [TypedSQL](/orm/prisma-client/using-raw-sql/typedsql), using these legacy methods is still the recommended approach for working with `Unsupported` fields.

:::

## What is SafeQL?

[SafeQL](https://safeql.dev/) allows for advanced linting and type safety within raw SQL queries. After setup, SafeQL works with Prisma Client `$queryRaw` and `$executeRaw` to provide type safety when raw queries are required.

SafeQL runs as an [ESLint](https://eslint.org/) plugin and is configured using ESLint rules. This guide doesn't cover setting up ESLint and we will assume that you already having it running in your project.

## Prerequisites

To follow along, you will be expected to have:

- A [PostgreSQL](https://www.postgresql.org/) database with PostGIS installed
- Prisma ORM set up in your project
- ESLint set up in your project

## Geographic data support in Prisma ORM

At the time of writing, Prisma ORM does not support working with geographic data, specifically using [PostGIS](https://github.com/prisma/prisma/issues/2789).

A model that has geographic data columns will be stored using the [`Unsupported`](/orm/reference/prisma-schema-reference#unsupported) data type. Fields with `Unsupported` types are present in the generated Prisma Client and will be typed as `any`. A model with a required `Unsupported` type does not expose write operations such as `create`, and `update`.

Prisma Client supports write operations on models with a required `Unsupported` field using `$queryRaw` and `$executeRaw`. You can use Prisma Client extensions and SafeQL to improve the type-safety when working with geographical data in raw queries.

## 1. Set up Prisma ORM for use with PostGIS

If you haven't already, enable the `postgresqlExtensions` Preview feature and add the `postgis` PostgreSQL extension in your Prisma schema:

```prisma
generator client {
 provider = "prisma-client"
 output = "./generated"
 previewFeatures = ["postgresqlExtensions"] // [!code ++]
}

datasource db {
 provider = "postgresql"
 extensions = [postgis] // [!code ++]
}
```

:::warning

If you are not using a hosted database provider, you will likely need to install the `postgis` extension. Refer to [PostGIS's docs](http://postgis.net/documentation/getting_started/#installing-postgis) to learn more about how to get started with PostGIS. If you're using Docker Compose, you can use the following snippet to set up a PostgreSQL database that has PostGIS installed:

```yaml
version: "3.6"
services:
 pgDB:
 image: postgis/postgis:13-3.1-alpine
 restart: always
 ports:
 - "5432:5432"
 volumes:
 - db_data:/var/lib/postgresql/data
 environment:
 POSTGRES_PASSWORD: password
 POSTGRES_DB: geoexample
volumes:
 db_data:
```

:::

Next, create a migration and execute a migration to enable the extension:

```npm
npx prisma migrate dev --name add-postgis
```

For reference, the output of the migration file should look like the following:

```sql title="migrations/TIMESTAMP_add_postgis/migration.sql"
-- CreateExtension
CREATE EXTENSION IF NOT EXISTS "postgis";
```

You can double-check that the migration has been applied by running `prisma migrate status`.

## 2. Create a new model that uses a geographic data column

Add a new model with a column with a `geography` data type once the migration is applied. For this guide, we'll use a model called `PointOfInterest`.

```prisma
model PointOfInterest {
 id Int @id @default(autoincrement())
 name String
 location Unsupported("geography(Point, 4326)")
}
```

You'll notice that the `location` field uses an [`Unsupported`](/orm/reference/prisma-schema-reference#unsupported) type. This means that we lose a lot of the benefits of Prisma ORM when working with `PointOfInterest`. We'll be using [SafeQL](https://safeql.dev/) to fix this.

Like before, create and execute a migration using the `prisma migrate dev` command to create the `PointOfInterest` table in your database:

```npm
npx prisma migrate dev --name add-poi
```

For reference, here is the output of the SQL migration file generated by Prisma Migrate:

```sql title="migrations/TIMESTAMP_add_poi/migration.sql"
-- CreateTable
CREATE TABLE "PointOfInterest" (
 "id" SERIAL NOT NULL,
 "name" TEXT NOT NULL,
 "location" geography(Point, 4326) NOT NULL,

 CONSTRAINT "PointOfInterest_pkey" PRIMARY KEY ("id")
);
```

## 3. Integrate SafeQL

SafeQL is easily integrated with Prisma ORM in order to lint `$queryRaw` and `$executeRaw` Prisma operations. You can reference [SafeQL's integration guide](https://safeql.dev/compatibility/prisma.html) or follow the steps below.

### 3.1. Install the `@ts-safeql/eslint-plugin` npm package

```npm
npm install -D @ts-safeql/eslint-plugin libpg-query
```

This ESLint plugin is what will allow for queries to be linted.

### 3.2. Add `@ts-safeql/eslint-plugin` to your ESLint plugins

Next, add `@ts-safeql/eslint-plugin` to your list of ESLint plugins. In our example we are using an `.eslintrc.js` file, but this can be applied to any way that you [configure ESLint](https://eslint.org/docs/latest/use/configure/).

```js title=".eslintrc.js" highlight=3
/** @type {import('eslint').Linter.Config} */
module.exports = {
 "plugins": [..., "@ts-safeql/eslint-plugin"],
 ...
}
```

### 3.3 Add `@ts-safeql/check-sql` rules

Now, setup the rules that will enable SafeQL to mark invalid SQL queries as ESLint errors.

```js title=".eslintrc.js" highlight=4-22;add
/** @type {import('eslint').Linter.Config} */
module.exports = {
 plugins: [..., '@ts-safeql/eslint-plugin'],
 rules: { // [!code ++]
 '@ts-safeql/check-sql': [ // [!code ++]
 'error', // [!code ++]
 { // [!code ++]
 connections: [ // [!code ++]
 { // [!code ++]
 // The migrations path: // [!code ++]
 migrationsDir: './prisma/migrations', // [!code ++]
 targets: [ // [!code ++]
 // This makes `prisma.$queryRaw` and `prisma.$executeRaw` commands linted // [!code ++]
 { tag: 'prisma.+($queryRaw|$executeRaw)', transform: '{type}[]' }, // [!code ++]
 ], // [!code ++]
 }, // [!code ++]
 ], // [!code ++]
 }, // [!code ++]
 ], // [!code ++]
 }, // [!code ++]
} // [!code ++]
```

> **Note**: If your `PrismaClient` instance is called something different than `prisma`, you need to adjust the value for `tag` accordingly. For example, if it is called `db`, the value for `tag` should be `'db.+($queryRaw|$executeRaw)'`.

### 3.4. Connect to your database

Finally, set up a `connectionUrl` for SafeQL so that it can introspect your database and retrieve the table and column names you use in your schema. SafeQL then uses this information for linting and highlighting problems in your raw SQL statements.

Our example relies on the [`dotenv`](https://github.com/motdotla/dotenv) package to get the same connection string that is used by Prisma ORM. We recommend this in order to keep your database URL out of version control.

If you haven't installed `dotenv` yet, you can install it as follows:

```npm
npm install dotenv
```

Then update your ESLint config as follows:

```js title=".eslintrc.js" highlight=1,6-9,16;add
require("dotenv").config(); // [!code ++]

/** @type {import('eslint').Linter.Config} */
module.exports = {
 plugins: ["@ts-safeql/eslint-plugin"],
 // exclude `parserOptions` if you are not using TypeScript // [!code ++]
 parserOptions: {
 // [!code ++]
 project: "./tsconfig.json", // [!code ++]
 }, // [!code ++]
 rules: {
 "@ts-safeql/check-sql": [
 "error",
 {
 connections: [
 {
 connectionUrl: process.env.DATABASE_URL, // [!code ++]
 // The migrations path:
 migrationsDir: "./prisma/migrations",
 targets: [
 // what you would like SafeQL to lint. This makes `prisma.$queryRaw` and `prisma.$executeRaw`
 // commands linted
 { tag: "prisma.+($queryRaw|$executeRaw)", transform: "{type}[]" },
 ],
 },
 ],
 },
 ],
 },
};
```

SafeQL is now fully configured to help you write better raw SQL using Prisma Client.

## 4. Creating extensions to make raw SQL queries type-safe

In this section, we'll create two [`model`](/orm/prisma-client/client-extensions/model) extensions with custom queries to be able to work conveniently with the `PointOfInterest` model:

1. A `create` query that allows us to create new `PointOfInterest` records in the database
1. A `findClosestPoints` query that returns the `PointOfInterest` records that are closest to a given coordinate

### 4.1. Adding an extension to create `PointOfInterest` records

The `PointOfInterest` model in the Prisma schema uses an `Unsupported` type. As a consequence, the generated `PointOfInterest` type in Prisma Client can't be used to carry values for latitude and longitude.

We will resolve this by defining two custom types that better represent our model in TypeScript:

```ts
type MyPoint = {
 latitude: number;
 longitude: number;
};

type MyPointOfInterest = {
 name: string;
 location: MyPoint;
};
```

Next, you can add a `create` query to the `pointOfInterest` property of your Prisma Client:

```ts highlight=19;normal
const prisma = new PrismaClient().$extends({
 model: {
 pointOfInterest: {
 async create(data: { name: string; latitude: number; longitude: number }) {
 // Create an object using the custom types from above
 const poi: MyPointOfInterest = {
 name: data.name,
 location: {
 latitude: data.latitude,
 longitude: data.longitude,
 },
 };

 // Insert the object into the database
 const point = `POINT(${poi.location.longitude} ${poi.location.latitude})`;
 await prisma.$queryRaw`
 INSERT INTO "PointOfInterest" (name, location) VALUES (${poi.name}, ST_GeomFromText(${point}, 4326));
 `;

 // Return the object
 return poi;
 },
 },
 },
});
```

Notice that the SQL in the line that's highlighted in the code snippet gets checked by SafeQL! For example, if you change the name of the table from `"PointOfInterest"` to `"PointOfInterest2"`, the following error appears:

```
error Invalid Query: relation "PointOfInterest2" does not exist @ts-safeql/check-sql
```

This also works with the column names `name` and `location`.

You can now create new `PointOfInterest` records in your code as follows:

```ts
const poi = await prisma.pointOfInterest.create({
 name: "Berlin",
 latitude: 52.52,
 longitude: 13.405,
});
```

### 4.2. Adding an extension to query for closest to `PointOfInterest` records

Now let's make a Prisma Client extension in order to query this model. We will be making an extension that finds the closest points of interest to a given longitude and latitude.

```ts
const prisma = new PrismaClient().$extends({
 model: {
 pointOfInterest: {
 async create(data: { name: string; latitude: number; longitude: number }) {
 // ... same code as before
 },

 async findClosestPoints(latitude: number, longitude: number) {
 // Query for closest points of interest
 const result = await prisma.$queryRaw<
 {
 id: number | null;
 name: string | null;
 st_x: number | null;
 st_y: number | null;
 }[]
 >`SELECT id, name, ST_X(location::geometry), ST_Y(location::geometry)
 FROM "PointOfInterest"
 ORDER BY ST_DistanceSphere(location::geometry, ST_MakePoint(${longitude}, ${latitude})) DESC`;

 // Transform to our custom type
 const pois: MyPointOfInterest[] = result.map((data) => {
 return {
 name: data.name,
 location: {
 latitude: data.st_x || 0,
 longitude: data.st_y || 0,
 },
 };
 });

 // Return data
 return pois;
 },
 },
 },
});
```

Now, you can use our Prisma Client as normal to find close points of interest to a given longitude and latitude using the custom method created on the `PointOfInterest` model.

```ts
const closestPointOfInterest = await prisma.pointOfInterest.findClosestPoints(53.5488, 9.9872);
```

Similar to before, we again have the benefit of SafeQL to add extra type safety to our raw queries. For example, if we removed the cast to `geometry` for `location` by changing `location::geometry` to just `location`, we would get linting errors in the `ST_X`, `ST_Y` or `ST_DistanceSphere` functions respectively.

```bash
error Invalid Query: function st_distancesphere(geography, geometry) does not exist @ts-safeql/check-sql
```

## Conclusion

While you may sometimes need to drop down to raw SQL when using Prisma ORM, you can use various techniques to make the experience of writing raw SQL queries with Prisma ORM better.

In this article, you have used SafeQL and Prisma Client extensions to create custom, type-safe Prisma Client queries to abstract PostGIS operations which are currently not natively supported in Prisma ORM.

---
title: TypedSQL
description: Learn how to use TypedSQL to write fully type-safe SQL queries that are compatible with any SQL console and Prisma Client
badge: preview
url: /orm/prisma-client/using-raw-sql/typedsql
metaTitle: Writing Type-safe SQL with TypedSQL and Prisma Client
metaDescription: Learn how to use TypedSQL to write fully type-safe SQL queries that are compatible with any SQL console and Prisma Client.
---

## Getting started with TypedSQL

To start using TypedSQL in your Prisma project, follow these steps:

1. Ensure you have `@prisma/client` and `prisma` installed:

 ```npm
 npm install @prisma/client@latest
 npm install -D prisma@latest
 ```

1. Add the `typedSql` preview feature flag to your `schema.prisma` file:

 ```prisma
 generator client {
 provider = "prisma-client"
 previewFeatures = ["typedSql"]
 output = "../src/generated/prisma"
 }
 ```

 :::tip[Using driver adapters with TypedSQL]

 If you are deploying Prisma in serverless or edge environments, you can use [driver adapters](/orm/core-concepts/supported-databases/database-drivers#driver-adapters) to connect through JavaScript database drivers. Driver adapters are compatible with TypedSQL, with the exception of `@prisma/adapter-better-sqlite3`. For SQLite support, use [`@prisma/adapter-libsql`](https://www.npmjs.com/package/@prisma/adapter-libsql) instead. All other driver adapters are supported.

 :::

1. Create a `sql` directory inside your `prisma` directory. This is where you'll write your SQL queries.

 ```bash
 mkdir -p prisma/sql
 ```

 :::note[Custom SQL folder location]

 Starting with Prisma 6.12.0, you can configure a custom location for your SQL files using the Prisma config file. Create a `prisma.config.ts` file in your project root and specify the `typedSql.path` option:

 ```typescript title="prisma.config.ts"
 import "dotenv/config";
 import { defineConfig } from "prisma/config";

 export default defineConfig({
 schema: "./prisma/schema.prisma",
 typedSql: {
 path: "./prisma/sql",
 },
 });
 ```

 :::

1. Create a new `.sql` file in your `prisma/sql` directory. For example, `getUsersWithPosts.sql`. Note that the file name must be a valid JS identifier and cannot start with a `$`.

1. Write your SQL queries in your new `.sql` file. For example:

 ```sql title="prisma/sql/getUsersWithPosts.sql"
 SELECT u.id, u.name, COUNT(p.id) as "postCount"
 FROM "User" u
 LEFT JOIN "Post" p ON u.id = p."authorId"
 GROUP BY u.id, u.name
 ```

1. Generate Prisma Client with the `sql` flag to ensure TypeScript functions and types for your SQL queries are created:

 :::warning

 Make sure that any pending migrations are applied before generating the client with the `sql` flag.

 :::

 ```bash
 prisma generate --sql
 ```

 If you don't want to regenerate the client after every change, this command also works with the existing `--watch` flag:

 ```bash
 prisma generate --sql --watch
 ```

1. Now you can import and use your SQL queries in your TypeScript code:

```typescript title="/src/index.ts"
import { PrismaClient } from "./generated/prisma/client";
import { getUsersWithPosts } from "./generated/prisma/sql";

const prisma = new PrismaClient();

const usersWithPostCounts = await prisma.$queryRawTyped(getUsersWithPosts());
console.log(usersWithPostCounts);
```

:::note

If you do not customize the generator `output`, you can import from `@prisma/client` and `@prisma/client/sql` instead.

:::

## Passing Arguments to TypedSQL Queries

To pass arguments to your TypedSQL queries, you can use parameterized queries. This allows you to write flexible and reusable SQL statements while maintaining type safety. Here's how to do it:

1. In your SQL file, use placeholders for the parameters you want to pass. The syntax for placeholders depends on your database engine:

For PostgreSQL, use the positional placeholders `$1`, `$2`, etc. For MySQL, use `?`. In SQLite, you can use positional (`$1`, `$2`), general (`?`), or named placeholders (`:minAge`, `:maxAge`):

```sql title="prisma/sql/getUsersByAge.sql" tab="PostgreSQL"
SELECT id, name, age
FROM users
WHERE age > $1 AND age < $2
```

```sql title="prisma/sql/getUsersByAge.sql" tab="MySQL"
SELECT id, name, age
FROM users
WHERE age > ? AND age < ?
```

```sql title="prisma/sql/getUsersByAge.sql" tab="SQLite"
SELECT id, name, age
FROM users
WHERE age > :minAge AND age < :maxAge
```

:::note

See below for information on how to [define argument types in your SQL files](#defining-argument-types-in-your-sql-files).

:::

1. When using the generated function in your TypeScript code, pass the arguments as additional parameters to `$queryRawTyped`:

```typescript title="/src/index.ts"
import { PrismaClient } from "./generated/prisma/client";
import { getUsersByAge } from "./generated/prisma/sql";

const prisma = new PrismaClient();

const minAge = 18;
const maxAge = 30;
const users = await prisma.$queryRawTyped(getUsersByAge(minAge, maxAge));
console.log(users);
```

By using parameterized queries, you ensure type safety and protect against SQL injection vulnerabilities. The TypedSQL generator will create the appropriate TypeScript types for the parameters based on your SQL query, providing full type checking for both the query results and the input parameters.

### Passing array arguments to TypedSQL

TypedSQL supports passing arrays as arguments for PostgreSQL. Use PostgreSQL's `ANY` operator with an array parameter.

```sql title="prisma/sql/getUsersByIds.sql"
SELECT id, name, email
FROM users
WHERE id = ANY($1)
```

```typescript title="/src/index.ts"
import { PrismaClient } from "./generated/prisma/client";
import { getUsersByIds } from "./generated/prisma/sql";

const prisma = new PrismaClient();

const userIds = [1, 2, 3];
const users = await prisma.$queryRawTyped(getUsersByIds(userIds));
console.log(users);
```

TypedSQL will generate the appropriate TypeScript types for the array parameter, ensuring type safety for both the input and the query results.

:::note

When passing array arguments, be mindful of the maximum number of placeholders your database supports in a single query. For very large arrays, you may need to split the query into multiple smaller queries.

:::

### Defining argument types in your SQL files

Argument typing in TypedSQL is accomplished via specific comments in your SQL files. These comments are of the form:

```sql
-- @param {Type} $N:alias optional description
```

Where `Type` is a valid database type, `N` is the position of the argument in the query, and `alias` is an optional alias for the argument that is used in the TypeScript type.

As an example, if you needed to type a single string argument with the alias `name` and the description "The name of the user", you would add the following comment to your SQL file:

```sql
-- @param {String} $1:name The name of the user
```

To indicate that a parameter is nullable, add a question mark after the alias:

```sql
-- @param {String} $1:name? The name of the user (optional)
```

Currently accepted types are `Int`, `BigInt`, `Float`, `Boolean`, `String`, `DateTime`, `Json`, `Bytes`, `null`, and `Decimal`.

Taking the [example from above](#passing-arguments-to-typedsql-queries), the SQL file would look like this:

```sql
-- @param {Int} $1:minAge
-- @param {Int} $2:maxAge
SELECT id, name, age
FROM users
WHERE age > $1 AND age < $2
```

The format of argument type definitions is the same regardless of the database engine.

:::note

Manual argument type definitions are not supported for array arguments. For these arguments, you will need to rely on the type inference provided by TypedSQL.

:::

## Examples

For practical examples of how to use TypedSQL, please refer to the [TypedSQL example in the Prisma Examples repo](https://github.com/prisma/prisma-examples/tree/latest/generator-prisma-client/basic-typedsql).

## Limitations of TypedSQL

### Supported Databases

TypedSQL supports modern versions of MySQL and PostgreSQL without any further configuration. For MySQL versions older than 8.0 and all SQLite versions, you will need to manually [describe argument types](#defining-argument-types-in-your-sql-files) in your SQL files. The types of inputs are inferred in all supported versions of PostgreSQL and MySQL 8.0 and later.

TypedSQL does not work with MongoDB, as it is specifically designed for SQL databases.

### Active Database Connection Required

TypedSQL requires an active database connection to function properly. This means you need to have a running database instance that Prisma can connect to when generating the client with the `--sql` flag. TypedSQL uses the connection string defined in `prisma.config.ts` (`datasource.url`) to establish this connection.

### Dynamic SQL Queries with Dynamic Columns

TypedSQL does not natively support constructing SQL queries with dynamically added columns. When you need to create a query where the columns are determined at runtime, you must use the `$queryRawUnsafe` and `$executeRawUnsafe` methods. These methods allow for the execution of raw SQL, which can include dynamic column selections.

**Example of a query using dynamic column selection:**

```typescript
const columns = "name, email, age"; // Columns determined at runtime
const result = await prisma.$queryRawUnsafe(`SELECT ${columns} FROM Users WHERE active = true`);
```

In this example, the columns to be selected are defined dynamically and included in the SQL query. While this approach provides flexibility, it requires careful attention to security, particularly to [avoid SQL injection vulnerabilities](/orm/prisma-client/using-raw-sql/raw-queries#sql-injection-prevention). Additionally, using raw SQL queries means foregoing the type-safety and DX of TypedSQL.

## Acknowledgements

This feature was heavily inspired by [PgTyped](https://github.com/adelsz/pgtyped) and [SQLx](https://github.com/launchbadge/sqlx). Additionally, SQLite parsing is handled by SQLx.

---
title: Shared Prisma Client extensions
description: Share extensions or import shared extensions into your Prisma project
url: /orm/prisma-client/client-extensions/shared-extensions
metaTitle: 'Shared Prisma Client extensions'
metaDescription: 'Share extensions or import shared extensions into your Prisma project'

---

You can share your [Prisma Client extensions](/orm/prisma-client/client-extensions) with other users, either as packages or as modules, and import extensions that other users create into your project.

If you would like to build a shareable extension, we also recommend using the [`prisma-client-extension-starter`](https://github.com/prisma/prisma-client-extension-starter) template.

To explore examples of Prisma's official Client extensions and those made by the community, visit [this](/orm/prisma-client/client-extensions/extension-examples) page.

## Install a shared, packaged extension

In your project, you can install any Prisma Client extension that another user has published to `npm`. To do so, run the following command:

```npm
npm install prisma-extension-<package-name>
```

For example, if the package name for an available extension is `prisma-extension-find-or-create`, you could install it as follows:

```npm
npm install prisma-extension-find-or-create
```

To import the `find-or-create` extension from the example above, and wrap your client instance with it, you could use the following code. This example assumes that the extension name is `findOrCreate`.

```ts
import findOrCreate from "prisma-extension-find-or-create";
import { PrismaClient } from "../generated/prisma/client";
const prisma = new PrismaClient();
const xprisma = prisma.$extends(findOrCreate);
const user = await xprisma.user.findOrCreate();
```

When you call a method in an extension, use the constant name from your `$extends` statement, not `prisma`. In the above example, `xprisma.user.findOrCreate` works, but `prisma.user.findOrCreate` does not, because the original `prisma` is not modified.

## Create a shareable extension

When you want to create extensions other users can use, and that are not tailored just for your schema, Prisma ORM provides utilities to allow you to create shareable extensions.

To create a shareable extension:

1. Define the extension as a module using `Prisma.defineExtension`
2. Use one of the methods that begin with the `$all` prefix such as [`$allModels`](/orm/prisma-client/client-extensions/model#add-a-custom-method-to-all-models-in-your-schema) or [`$allOperations`](/orm/prisma-client/client-extensions/query#modify-all-prisma-client-operations)

### Define an extension

Use the `Prisma.defineExtension` method to make your extension shareable. You can use it to package the extension to either separate your extensions into a separate file or share it with other users as an npm package.

The benefit of `Prisma.defineExtension` is that it provides strict type checks and auto completion for authors of extension in development and users of shared extensions.

### Use a generic method

Extensions that contain methods under `$allModels` apply to every model instead of a specific one. Similarly, methods under `$allOperations` apply to a client instance as a whole and not to a named component, e.g. `result` or `query`.

You do not need to use the `$all` prefix with the [`client`](/orm/prisma-client/client-extensions/client) component, because the `client` component always applies to the client instance.

For example, a generic extension might take the following form:

```ts
export default Prisma.defineExtension({
 name: "prisma-extension-find-or-create", //Extension name
 model: {
 $allModels: {
 // new method
 findOrCreate(/* args */) {
 /* code for the new method */
 return query(args);
 },
 },
 },
});
```

Refer to the following pages to learn the different ways you can modify Prisma Client operations:

- [Modify all Prisma Client operations](/orm/prisma-client/client-extensions/query#modify-all-prisma-client-operations)
- [Modify a specific operation in all models of your schema](/orm/prisma-client/client-extensions/query#modify-a-specific-operation-in-all-models-of-your-schema)
- [Modify all operations in all models of your schema](/orm/prisma-client/client-extensions/query#modify-all-operations-in-all-models-of-your-schema)

### Publishing the shareable extension to npm

You can then share the extension on `npm`. When you choose a package name, we recommend that you use the `prisma-extension-<package-name>` convention, to make it easier to find and install.

### Call a client-level method from your packaged extension

:::warning

There's currently a limitation for extensions that reference a `PrismaClient` and call a client-level method, like the example below.

If you trigger the extension from inside a [transaction](/orm/prisma-client/queries/transactions) (interactive or batched), the extension code will issue the queries in a new connection and ignore the current transaction context.

Learn more in this issue on GitHub: [Client extensions that require use of a client-level method silently ignore transactions](https://github.com/prisma/prisma/issues/20678).

:::

In the following situations, you need to refer to a Prisma Client instance that your extension wraps:

- When you want to use a [client-level method](/orm/reference/prisma-client-reference#client-methods), such as `$queryRaw`, in your packaged extension.
- When you want to chain multiple `$extends` calls in your packaged extension.

However, when someone includes your packaged extension in their project, your code cannot know the details of the Prisma Client instance.

You can refer to this client instance as follows:

```ts
Prisma.defineExtension((client) => {
 // The Prisma Client instance that the extension user applies the extension to
 return client.$extends({
 name: "prisma-extension-<extension-name>",
 });
});
```

For example:

```ts
export default Prisma.defineExtension((client) => {
 return client.$extends({
 name: "prisma-extension-find-or-create",
 query: {
 $allModels: {
 async findOrCreate({ args, query, operation }) {
 return (await client.$transaction([query(args)]))[0];
 },
 },
 },
 });
});
```

### Advanced type safety: type utilities for defining generic extensions

You can improve the type-safety of your shared extensions using [type utilities](/orm/prisma-client/client-extensions/type-utilities).

---
title: Fine-Grained Authorization (Permit)
description: 'Learn how to implement RBAC, ABAC, and ReBAC authorization in your Prisma applications'
url: /orm/prisma-client/client-extensions/shared-extensions/permit-rbac
metaDescription: 'Learn how to implement RBAC, ABAC, and ReBAC authorization in your Prisma applications'
metaTitle: Fine-Grained Authorization (Permit)
---

:::info[Quick summary]
This page explains how to implement fine-grained authorization (FGA) in Prisma ORM applications using the `@permitio/permit-prisma` extension. It introduces different access control models—RBAC, ABAC, and ReBAC—supported by Permit.io, and guides you on choosing the right model to protect your database operations with precise, programmable permissions.

:::

Database operations often require careful control over who can access or modify which data. While Prisma ORM excels at data modeling and database access, it doesn't include built-in authorization capabilities. This guide shows how to implement fine-grained authorization in your Prisma applications using the `@permitio/permit-prisma` extension.

Fine-grained authorization (FGA) provides detailed and precise control over what data users can access or modify at a granular level. Without proper authorization, your application might expose sensitive data or allow unauthorized modifications, creating security vulnerabilities.

## Access control models

This extension supports three access control models from Permit.io:

### Role-based Access Control (RBAC)

**What it is**: Users are assigned roles (Admin, Editor, Viewer) with predefined permissions to perform actions on resource types.

**Example**: An "Editor" role can update any document in the system.

**Best for**: Simple permission structures where access is determined by job function or user level.

### Attribute-Based Access Control (ABAC)

**What it is**: Access decisions based on attributes of users, resources, or environment.

**Examples**:

- Allow access if `user.department == document.department`
- Allow updates if `document.status == "DRAFT"`

**How it works with the extension**: When `enableAttributeSync` is on, resource attributes are automatically synced to Permit.io for policy evaluation.

**Best for**: Dynamic rules that depend on context or data properties.

### Relationship-Based Access Control (ReBAC)

**What it is**: Permissions based on relationships between users and specific resource instances.

**Example**: A user is an "Owner" of document-123 but just a "Viewer" of document-456.

**How it works with the extension**:

- Resource instances are synced to Permit.io (with `enableResourceSync: true`)
- Permission checks include the specific resource instance ID

**Best for**: Collaborative applications where users need different permissions on different instances of the same resource type.

### Choosing the right model

- **RBAC**: When you need simple, role-based access control
- **ABAC**: When decisions depend on data properties or contextual information
- **ReBAC**: When users need different permissions on different instances

## Usage

### Prerequisites

Before implementing fine-grained authorization with Prisma, make sure you have:

- A Prisma application with existing models and queries
- Basic understanding of authorization concepts
- Node.js and npm installed

### Installation

Install the extension alongside Prisma Client:

```npm
npm install @permitio/permit-prisma @prisma/client
```

You'll also need to sign up for a [Permit account](https://app.permit.io) to define your authorization policies.

> **Note:**
> Ensure that the Permit PDP container is running. It is recommended to run it using Docker for better performance, security, and availability. For instructions, refer to the Permit documentation: [Deploy Permit to Production](https://docs.permit.io/how-to/deploy/deploy-to-production/) and [PDP Overview](https://docs.permit.io/concepts/pdp/overview/).

## Basic setup

First, extend your Prisma Client with the Permit extension:

```typescript
import { PrismaClient } from "@prisma/client";
import { createPermitClientExtension } from "@permitio/permit-prisma";

const prisma = new PrismaClient().$extends(
 createPermitClientExtension({
 permitConfig: {
 token: process.env.PERMIT_API_KEY, // Your Permit API key
 pdp: "http://localhost:7766", // PDP address (local or cloud)
 },
 enableAutomaticChecks: true, // Automatically enforce permissions
 }),
);
```

## Implementing RBAC (Role-Based Access Control)

RBAC uses roles to determine access permissions. For example, "Admin" roles can perform all actions while "Viewer" roles can only read data.

1. **Define resources and actions in Permit.io dashboard**:
 - Create resources matching your Prisma models (e.g., "document")
 - Define actions (e.g., "create", "read", "update", "delete")
 - Create roles with permission sets (e.g., "admin", "editor", "viewer")
2. **Set the active user in your code**:

```typescript
// Set the current user context before performing operations
prisma.$permit.setUser("john@example.com");

// All subsequent operations will be checked against this user's permissions
const documents = await prisma.document.findMany();
```

## Implementing ABAC (Attribute-Based Access Control)

ABAC extends access control by considering user attributes, resource attributes, and context.

1. **Configure the extension for ABAC**:

```typescript
const prisma = new PrismaClient().$extends(
 createPermitClientExtension({
 permitConfig: { token: process.env.PERMIT_API_KEY, pdp: "http://localhost:7766" },
 enableAutomaticChecks: true,
 }),
);
```

2. **Set user with attributes:**

```typescript
prisma.$permit.setUser({
 key: "doctor@hospital.com",
 attributes: { department: "cardiology" },
});

// Will succeed only if user department matches record department (per policy)
const records = await prisma.medicalRecord.findMany({
 where: { department: "cardiology" },
});
```

## Implementing ReBAC (Relationship-Based Access Control)

ReBAC models permissions based on relationships between users and specific resource instances.

1. **Configure the extension for ReBAC**:

```typescript
const prisma = new PrismaClient().$extends(
 createPermitClientExtension({
 permitConfig: { token: process.env.PERMIT_API_KEY, pdp: "http://localhost:7766" },
 accessControlModel: "rebac",
 enableAutomaticChecks: true,
 enableResourceSync: true, // Sync resource instances with Permit.io
 enableDataFiltering: true, // Filter queries by permissions
 }),
);
```

2. ** Access instance-specific resources:**

```typescript
prisma.$permit.setUser("owner@example.com");

// Will only succeed if the user has permission on this specific file
const file = await prisma.file.findUnique({
 where: { id: "file-123" },
});
```

## Manual permission checks

For more control, you can perform explicit permission checks:

```typescript
// Check if user can update a document
const canUpdate = await prisma.$permit.check(
 "john@example.com", // user
 "update", // action
 "document", // resource
);

if (canUpdate) {
 await prisma.document.update({
 where: { id: "doc-123" },
 data: { title: "Updated Title" },
 });
}

// Or enforce permissions (throws if denied)
await prisma.$permit.enforceCheck("john@example.com", "delete", {
 type: "document",
 key: "doc-123",
});
```

## Common use cases

Here are some common scenarios where fine-grained authorization is valuable:

- **Multi-tenant applications**: Isolate data between different customers
- **Healthcare applications**: Ensure patient data is only accessible to authorized staff
- **Collaborative platforms**: Grant different permissions on shared resources
- **Content management systems**: Control who can publish, edit, or view content

## Summary

By integrating the `@permitio/permit-prisma` extension with your Prisma ORM application, you can implement sophisticated authorization policies that protect your data and ensure users only access what they're permitted to see. The extension supports all major authorization models (RBAC, ABAC, ReBAC) and provides both automatic and manual permission enforcement.

## Next steps

- [Create a free Permit.io account](https://app.permit.io)
- [View the full extension documentation](https://github.com/permitio/permit-prisma)

---
title: Deploy to Cloudflare Workers & Pages
description: Learn the things you need to know in order to deploy an app that uses Prisma Client for talking to a database to a Cloudflare Worker or to Cloudflare Pages
url: /orm/prisma-client/deployment/edge/deploy-to-cloudflare
metaTitle: Deploy to Cloudflare Workers & Pages
metaDescription: Learn the things you need to know in order to deploy an app that uses Prisma Client for talking to a database to a Cloudflare Worker or to Cloudflare Pages.
---

:::info[Quick summary]
This page covers everything you need to know to deploy an app with Prisma ORM to a [Cloudflare Worker](https://developers.cloudflare.com/workers/) or to [Cloudflare Pages](https://developers.cloudflare.com/pages).
:::

<details>
<summary>Questions answered in this page</summary>

- How to deploy Prisma to Cloudflare Workers?
- Which drivers work on Workers/Pages?
- How to configure DATABASE_URL and envs?

</details>

## General considerations when deploying to Cloudflare Workers

This section covers _general_ things you need to be aware of when deploying to Cloudflare Workers or Pages and are using Prisma ORM, regardless of the database provider you use.

### Using Prisma Postgres

You can use Prisma Postgres and deploy to Cloudflare Workers.

After you create a Worker, run:

```npm
npx prisma@latest init --db
```

Enter a name for your project and choose a database region.

This command:

- Connects your CLI to your [Prisma Data Platform](https://console.prisma.io) account. If you're not logged in or don't have an account, your browser will open to guide you through creating a new account or signing into your existing one.
- Creates a `prisma` directory containing a `schema.prisma` file for your database models.
- Creates a `.env` file with your `DATABASE_URL`.

### Using an edge-compatible driver

When deploying a Cloudflare Worker that uses Prisma ORM, you need to use an [edge-compatible driver](/orm/prisma-client/deployment/edge/overview#edge-compatibility-of-database-drivers) and its respective [driver adapter](/orm/core-concepts/supported-databases/database-drivers#driver-adapters) for Prisma ORM.

The edge-compatible drivers for Cloudflare Workers and Pages are:

- [Neon Serverless](https://neon.tech/docs/serverless/serverless-driver) uses HTTP to access the database
- [PlanetScale Serverless](https://planetscale.com/docs/tutorials/planetscale-serverless-driver) uses HTTP to access the database
- [`node-postgres`](https://node-postgres.com/) (`pg`) uses Cloudflare's `connect()` (TCP) to access the database
- [`@libsql/client`](https://github.com/tursodatabase/libsql-client-ts) is used to access Turso databases via HTTP
- [Cloudflare D1](/orm/prisma-client/deployment/edge/deploy-to-cloudflare) is used to access D1 databases

There's [also work being done](https://github.com/sidorares/node-mysql2/pull/2289) on the `node-mysql2` driver which will enable access to traditional MySQL databases from Cloudflare Workers and Pages in the future as well.

If your application uses PostgreSQL, we recommend using [Prisma Postgres](/postgres). It is fully supported on edge runtimes and does not require a specialized edge-compatible driver. Review the [Prisma Postgres serverless driver limitations](/postgres/database/serverless-driver#limitations) to understand current constraints.

### Setting your database connection URL as an environment variable

First, ensure that your `datasource` block in your Prisma schema is configured correctly. Database connection URLs are configured in `prisma.config.ts`:

```prisma
datasource db {
 provider = "postgresql" // this might also be `mysql` or another value depending on your database
}
```

```ts title="prisma.config.ts"
import "dotenv/config";
import { defineConfig, env } from "prisma/config";

export default defineConfig({
 schema: "prisma/schema.prisma",
 datasource: {
 url: env("DATABASE_URL"),
 },
});
```

#### Development

When using your Worker in **development**, you can configure your database connection via the [`.dev.vars` file](https://developers.cloudflare.com/workers/configuration/secrets/#local-development-with-secrets) locally.

Assuming you use the `DATABASE_URL` environment variable from above, you can set it inside `.dev.vars` as follows:

```bash title=".dev.vars"
DATABASE_URL="your-database-connection-string"
```

In the above snippet, `your-database-connection-string` is a placeholder that you need to replace with the value of your own connection string, for example:

```bash title=".dev.vars"
DATABASE_URL="postgresql://admin:mypassword42@somehost.aws.com:5432/mydb"
```

Note that the `.dev.vars` file is not compatible with `.env` files which are typically used by Prisma ORM.

This means that you need to make sure that Prisma ORM gets access to the environment variable when needed, e.g. when running a Prisma CLI command like `prisma migrate dev`.

There are several options for achieving this:

- Run your Prisma CLI commands using [`dotenv`](https://www.npmjs.com/package/dotenv-cli) to specify from where the CLI should read the environment variable, for example:
 ```bash
 dotenv -e .dev.vars -- npx prisma migrate dev
 ```
- Create a script in `package.json` that reads `.dev.vars` via [`dotenv`](https://www.npmjs.com/package/dotenv-cli). You can then execute `prisma` commands as follows: `npm run env -- npx prisma migrate dev`. Here's a reference for the script:
 ```js title="package.json"
 "scripts": { "env": "dotenv -e .dev.vars" }
 ```
- Duplicate the `DATABASE_URL` and any other relevant env vars into a new file called `.env` which can then be used by Prisma ORM.

:::note

If you're using an approach that requires `dotenv`, you need to have the [`dotenv-cli`](https://www.npmjs.com/package/dotenv-cli) package installed. You can do this e.g. by using this command to install the package locally in your project: `npm install -D dotenv-cli`.

:::

#### Production

When deploying your Worker to **production**, you'll need to set the database connection using the `wrangler` CLI:

```npm
npx wrangler secret put DATABASE_URL
```

The command is interactive and will ask you to enter the value for the `DATABASE_URL` env var as the next step in the terminal.

:::note

This command requires you to be authenticated, and will ask you to log in to your Cloudflare account in case you are not.

:::

### Size limits on free accounts

Cloudflare has a [size limit of 3 MB for Workers on the free plan](https://developers.cloudflare.com/workers/platform/limits/). If your application bundle with Prisma ORM exceeds that size, we recommend upgrading to a paid Worker plan.

### Deploying a Next.js app to Cloudflare Pages with `@cloudflare/next-on-pages`

Cloudflare offers an option to run Next.js apps on Cloudflare Pages with [`@cloudflare/next-on-pages`](https://github.com/cloudflare/next-on-pages), see the [docs](https://developers.cloudflare.com/pages/framework-guides/nextjs/ssr/get-started/) for instructions.

Based on some testing, we found the following:

- You can deploy using the PlanetScale or Neon Serverless Driver.
- Traditional PostgreSQL deployments using `pg` don't work because `pg` itself currently does not work with `@cloudflare/next-on-pages` (see [here](https://github.com/cloudflare/next-on-pages/issues/605)).

Feel free to reach out to us on [Discord](https://pris.ly/discord?utm_source=docs&utm_medium=inline_text) if you find that anything has changed about this.

## Database-specific considerations & examples

This section provides database-specific instructions for deploying a Cloudflare Worker with Prisma ORM.

### Prerequisites

As a prerequisite for the following section, you need to have a Cloudflare Worker running locally and the Prisma CLI installed.

If you don't have that yet, you can run these commands:

```npm
npm create cloudflare@latest prisma-cloudflare-worker-example -- --type hello-world
cd prisma-cloudflare-worker-example
npm install prisma --save-dev && npm install @prisma/client
npx prisma init --output ../generated/prisma
```

You'll further need a database instance of your database provider of choice available. Refer to the respective documentation of the provider for setting up that instance.

We'll use the default `User` model for the example below:

```prisma
model User {
 id Int @id @default(autoincrement())
 email String @unique
 name String?
}
```

### PostgreSQL (traditional)

If you are using a traditional PostgreSQL database that's accessed via TCP and the `pg` driver, you need to:

- use the `@prisma/adapter-pg` database adapter (learn more [here](/orm/core-concepts/supported-databases/postgresql#using-driver-adapters))
- set `node_compat = true` in `wrangler.toml` (see the [Cloudflare docs](https://developers.cloudflare.com/workers/runtime-apis/nodejs/))

#### 1. Configure Prisma schema & database connection

:::note

If you don't have a project to deploy, follow the instructions in the [Prerequisites](#prerequisites) to bootstrap a basic Cloudflare Worker with Prisma ORM in it.

:::

First, ensure that the database connection is configured properly. Database connection URLs are configured in `prisma.config.ts`:

```prisma title="schema.prisma"
generator client {
 provider = "prisma-client"
 output = "./generated"
}

datasource db {
 provider = "postgresql"
}
```

```ts title="prisma.config.ts"
import "dotenv/config";
import { defineConfig, env } from "prisma/config";

export default defineConfig({
 schema: "prisma/schema.prisma",
 datasource: {
 url: env("DATABASE_URL"),
 },
});
```

Next, you need to set the `DATABASE_URL` environment variable to the value of your database connection string. You'll do this in a file called `.dev.vars` used by Cloudflare:

```bash title=".dev.vars"
DATABASE_URL="postgresql://admin:mypassword42@somehost.aws.com:5432/mydb"
```

Because the Prisma CLI by default is only compatible with `.env` files, you can adjust your `package.json` with the following script that loads the env vars from `.dev.vars`. You can then use this script to load the env vars before executing a `prisma` command.

Add this script to your `package.json`:

```js title="package.json" highlight=5;add
{
 // ...
 "scripts": {
 // ....
 "env": "dotenv -e .dev.vars"
 },
 // ...
}
```

Now you can execute Prisma CLI commands as follows while ensuring that the command has access to the env vars in `.dev.vars`:

```npm
npm run env -- npx prisma
```

#### 2. Install dependencies

Next, install the required packages:

```npm
npm install @prisma/adapter-pg
```

#### 3. Set `node_compat = true` in `wrangler.toml`

In your `wrangler.toml` file, add the following line:

```toml title="wrangler.toml"
node_compat = true
```

:::note

For Cloudflare Pages, using `node_compat` is not officially supported. If you want to use `pg` in Cloudflare Pages, you can find a workaround [here](https://github.com/cloudflare/workers-sdk/pull/2541#issuecomment-1954209855).

:::

#### 4. Migrate your database schema (if applicable)

If you ran `npx prisma init` above, you need to migrate your database schema to create the `User` table that's defined in your Prisma schema (if you already have all the tables you need in your database, you can skip this step):

```npm
npm run env -- npx prisma migrate dev --name init
```

#### 5. Use Prisma Client in your Worker to send a query to the database

Here is a sample code snippet that you can use to instantiate `PrismaClient` and send a query to your database:

```ts
import { PrismaClient } from "./generated/client";
import { PrismaPg } from "@prisma/adapter-pg";

export default {
 async fetch(request, env, ctx) {
 const adapter = new PrismaPg({ connectionString: env.DATABASE_URL });
 const prisma = new PrismaClient({ adapter });

 const users = await prisma.user.findMany();
 const result = JSON.stringify(users);
 ctx.waitUntil(prisma.$disconnect());
 return new Response(result);
 },
};
```

#### 6. Run the Worker locally

To run the Worker locally, you can run the `wrangler dev` command:

```npm
npx wrangler dev
```

#### 7. Set the `DATABASE_URL` environment variable and deploy the Worker

To deploy the Worker, you first need to the `DATABASE_URL` environment variable [via the `wrangler` CLI](https://developers.cloudflare.com/workers/configuration/secrets/#secrets-on-deployed-workers):

```npm
npx wrangler secret put DATABASE_URL
```

The command is interactive and will ask you to enter the value for the `DATABASE_URL` env var as the next step in the terminal.

:::note

This command requires you to be authenticated, and will ask you to log in to your Cloudflare account in case you are not.

:::

Then you can go ahead then deploy the Worker:

```npm
npx wrangler deploy
```

The command will output the URL where you can access the deployed Worker.

### PlanetScale

If you are using a PlanetScale database, you need to:

- use the `@prisma/adapter-planetscale` database adapter (learn more [here](/orm/core-concepts/supported-databases/mysql#planetscale))
- manually remove the conflicting `cache` field:

 ```ts
 export default {
 async fetch(request, env, ctx) {
 const adapter = new PrismaPlanetScale({
 url: env.DATABASE_URL,
 // see https://github.com/cloudflare/workerd/issues/698
 fetch(url, init) {
 delete init["cache"];
 return fetch(url, init);
 },
 });
 const prisma = new PrismaClient({ adapter });

 // ...
 },
 };
 ```

#### 1. Configure Prisma schema & database connection

:::note

If you don't have a project to deploy, follow the instructions in the [Prerequisites](#prerequisites) to bootstrap a basic Cloudflare Worker with Prisma ORM in it.

:::

First, ensure that the database connection is configured properly. Database connection URLs are configured in `prisma.config.ts`:

```prisma title="schema.prisma"
generator client {
 provider = "prisma-client"
 output = "./generated"
}

datasource db {
 provider = "mysql"
 relationMode = "prisma" // required for PlanetScale (as by default foreign keys are disabled)
}
```

```ts title="prisma.config.ts"
import "dotenv/config";
import { defineConfig, env } from "prisma/config";

export default defineConfig({
 schema: "prisma/schema.prisma",
 datasource: {
 url: env("DATABASE_URL"),
 },
});
```

Next, you need to set the `DATABASE_URL` environment variable to the value of your database connection string. You'll do this in a file called `.dev.vars` used by Cloudflare:

```bash title=".dev.vars"
DATABASE_URL="mysql://32qxa2r7hfl3102wrccj:password@us-east.connect.psdb.cloud/demo-cf-worker-ps?sslaccept=strict"
```

Because the Prisma CLI by default is only compatible with `.env` files, you can adjust your `package.json` with the following script that loads the env vars from `.dev.vars`. You can then use this script to load the env vars before executing a `prisma` command.

Add this script to your `package.json`:

```js title="package.json" highlight=5;add
{
 // ...
 "scripts": {
 // ....
 "env": "dotenv -e .dev.vars"
 },
 // ...
}
```

Now you can execute Prisma CLI commands as follows while ensuring that the command has access to the env vars in `.dev.vars`:

```npm
npm run env -- npx prisma
```

#### 2. Install dependencies

Next, install the required packages:

```npm
npm install @prisma/adapter-planetscale
```

#### 3. Migrate your database schema (if applicable)

If you ran `npx prisma init` above, you need to migrate your database schema to create the `User` table that's defined in your Prisma schema (if you already have all the tables you need in your database, you can skip this step):

```npm
npm run env -- npx prisma db push
```

#### 4. Use Prisma Client in your Worker to send a query to the database

Here is a sample code snippet that you can use to instantiate `PrismaClient` and send a query to your database:

```ts
import { PrismaClient } from "./generated/client";
import { PrismaPlanetScale } from "@prisma/adapter-planetscale";

export default {
 async fetch(request, env, ctx) {
 const adapter = new PrismaPlanetScale({
 url: env.DATABASE_URL,
 // see https://github.com/cloudflare/workerd/issues/698
 fetch(url, init) {
 delete init["cache"];
 return fetch(url, init);
 },
 });
 const prisma = new PrismaClient({ adapter });

 const users = await prisma.user.findMany();
 const result = JSON.stringify(users);
 ctx.waitUntil(prisma.$disconnect());
 return new Response(result);
 },
};
```

#### 6. Run the Worker locally

To run the Worker locally, you can run the `wrangler dev` command:

```npm
npx wrangler dev
```

#### 7. Set the `DATABASE_URL` environment variable and deploy the Worker

To deploy the Worker, you first need to the `DATABASE_URL` environment variable [via the `wrangler` CLI](https://developers.cloudflare.com/workers/configuration/secrets/#secrets-on-deployed-workers):

```npm
npx wrangler secret put DATABASE_URL
```

The command is interactive and will ask you to enter the value for the `DATABASE_URL` env var as the next step in the terminal.

:::note

This command requires you to be authenticated, and will ask you to log in to your Cloudflare account in case you are not.

:::

Then you can go ahead then deploy the Worker:

```npm
npx wrangler deploy
```

The command will output the URL where you can access the deployed Worker.

### Neon

If you are using a Neon database, you need to:

- use the `@prisma/adapter-neon` database adapter (learn more [here](/orm/core-concepts/supported-databases/postgresql#using-driver-adapters))

#### 1. Configure Prisma schema & database connection

:::note

If you don't have a project to deploy, follow the instructions in the [Prerequisites](#prerequisites) to bootstrap a basic Cloudflare Worker with Prisma ORM in it.

:::

First, ensure that the database connection is configured properly. Database connection URLs are configured in `prisma.config.ts`:

```prisma title="schema.prisma"
generator client {
 provider = "prisma-client"
 output = "./generated"
}

datasource db {
 provider = "postgresql"
}
```

```ts title="prisma.config.ts"
import "dotenv/config";
import { defineConfig, env } from "prisma/config";

export default defineConfig({
 schema: "prisma/schema.prisma",
 datasource: {
 url: env("DATABASE_URL"),
 },
});
```

Next, you need to set the `DATABASE_URL` environment variable to the value of your database connection string. You'll do this in a file called `.dev.vars` used by Cloudflare:

```bash title=".dev.vars"
DATABASE_URL="postgresql://janedoe:password@ep-nameless-pond-a23b1mdz.eu-central-1.aws.neon.tech/neondb?sslmode=require"
```

Because the Prisma CLI by default is only compatible with `.env` files, you can adjust your `package.json` with the following script that loads the env vars from `.dev.vars`. You can then use this script to load the env vars before executing a `prisma` command.

Add this script to your `package.json`:

```js title="package.json" highlight=5;add
{
 // ...
 "scripts": {
 // ....
 "env": "dotenv -e .dev.vars"
 },
 // ...
}
```

Now you can execute Prisma CLI commands as follows while ensuring that the command has access to the env vars in `.dev.vars`:

```npm
npm run env -- npx prisma
```

#### 2. Install dependencies

Next, install the required packages:

```npm
npm install @prisma/adapter-neon
```

#### 3. Migrate your database schema (if applicable)

If you ran `npx prisma init` above, you need to migrate your database schema to create the `User` table that's defined in your Prisma schema (if you already have all the tables you need in your database, you can skip this step):

```npm
npm run env -- npx prisma migrate dev --name init
```

#### 5. Use Prisma Client in your Worker to send a query to the database

Here is a sample code snippet that you can use to instantiate `PrismaClient` and send a query to your database:

```ts
import { PrismaClient } from "./generated/client";
import { PrismaNeon } from "@prisma/adapter-neon";

export default {
 async fetch(request, env, ctx) {
 const adapter = new PrismaNeon({ connectionString: env.DATABASE_URL });
 const prisma = new PrismaClient({ adapter });

 const users = await prisma.user.findMany();
 const result = JSON.stringify(users);
 ctx.waitUntil(prisma.$disconnect());
 return new Response(result);
 },
};
```

#### 6. Run the Worker locally

To run the Worker locally, you can run the `wrangler dev` command:

```npm
npx wrangler dev
```

#### 7. Set the `DATABASE_URL` environment variable and deploy the Worker

To deploy the Worker, you first need to the `DATABASE_URL` environment variable [via the `wrangler` CLI](https://developers.cloudflare.com/workers/configuration/secrets/#secrets-on-deployed-workers):

```npm
npx wrangler secret put DATABASE_URL
```

The command is interactive and will ask you to enter the value for the `DATABASE_URL` env var as the next step in the terminal.

:::note

This command requires you to be authenticated, and will ask you to log in to your Cloudflare account in case you are not.

:::

Then you can go ahead then deploy the Worker:

```npm
npx wrangler deploy
```

The command will output the URL where you can access the deployed Worker.

### Cloudflare D1

:::info[Using Cloudflare D1]
For step-by-step instructions on using Prisma ORM with [Cloudflare D1](https://developers.cloudflare.com/d1/) (schema setup, migrations, and deploying your Worker), see the dedicated [Cloudflare D1 deployment guide](/guides/deployment/cloudflare-d1).
:::

---
title: Deploy to Deno Deploy
metaTitle: Deploy to Deno Deploy
metaDescription: Learn how to deploy a TypeScript application using Prisma ORM to Deno Deploy.
url: /orm/prisma-client/deployment/edge/deploy-to-deno-deploy
---

With this guide, you can learn how to build and deploy a REST API to [Deno Deploy](https://deno.com/deploy). The application uses Prisma ORM to manage tasks in a [Prisma Postgres](/postgres) database.

This guide covers Deno CLI, Deno Deploy, Prisma Client with the Postgres adapter, and Prisma Postgres.

## Prerequisites

- A free [Prisma Data Platform](https://console.prisma.io/login) account
- A free [Deno Deploy](https://deno.com/deploy) account
- [Deno](https://docs.deno.com/runtime/#install-deno) v2.0 or later installed
- (Recommended) [Deno extension for VS Code](https://docs.deno.com/runtime/reference/vscode/)

## 1. Set up your application

Create a new directory and initialize your Prisma project:

```npm
mkdir prisma-deno-deploy
cd prisma-deno-deploy
deno run -A npm:prisma@latest init --db
```

Enter a name for your project and choose a database region.

This command:
- Connects to the [Prisma Data Platform](https://console.prisma.io) (opens browser for authentication)
- Creates a `prisma/schema.prisma` file for your database models
- Creates a `.env` file with your `DATABASE_URL`
- Creates a `prisma.config.ts` configuration file

## 2. Configure Deno

Create a `deno.json` file with the following configuration:

```json title="deno.json"
{
 "nodeModulesDir": "auto",
 "compilerOptions": {
 "lib": ["deno.window"],
 "types": ["node"]
 },
 "imports": {
 "@prisma/adapter-pg": "npm:@prisma/adapter-pg@^7.0.0",
 "@prisma/client": "npm:@prisma/client@^7.0.0",
 "prisma": "npm:prisma@^7.0.0"
 },
 "tasks": {
 "dev": "deno run -A --env=.env --watch index.ts",
 "db:generate": "deno run -A --env=.env npm:prisma generate",
 "db:push": "deno run -A --env=.env npm:prisma db push",
 "db:migrate": "deno run -A --env=.env npm:prisma migrate dev",
 "db:studio": "deno run -A --env=.env npm:prisma studio"
 }
}
```

:::note
The `nodeModulesDir: "auto"` setting is required for Prisma to work correctly with Deno. The `compilerOptions` ensure TypeScript understands Deno globals and npm packages. The import map allows you to use clean import paths like `@prisma/adapter-pg` instead of `npm:@prisma/adapter-pg`.
:::

Install the dependencies:

```bash
deno install
deno install --allow-scripts
```

## 3. Define your data model

Edit `prisma/schema.prisma` to add the Deno runtime and a `Task` model:

```prisma title="prisma/schema.prisma" highlight=4;add
generator client {
 provider = "prisma-client"
 output = "../generated/prisma"
 runtime = "deno"
}

datasource db {
 provider = "postgresql"
}

model Task {
 id Int @id @default(autoincrement())
 title String
 description String?
 completed Boolean @default(false)
 createdAt DateTime @default(now())
 updatedAt DateTime @updatedAt
}
```

## 4. Push the schema to your database

Apply the schema to your database and generate Prisma Client:

```bash
deno task db:push
```

This command:
1. Creates the `Task` table in your Prisma Postgres database
2. Generates the Prisma Client with full type safety

:::info
If you see TypeScript errors in your IDE after generating, restart the Deno language server (`Cmd/Ctrl + Shift + P` → "Deno: Restart Language Server") to refresh the types.
:::

## 5. Create your application

Create `index.ts` with a REST API for managing tasks:

```ts title="index.ts"
import { PrismaPg } from "@prisma/adapter-pg";
import { PrismaClient } from "./generated/prisma/client.ts";

// Initialize Prisma Client with the Postgres adapter
const connectionString = Deno.env.get("DATABASE_URL")!;
const adapter = new PrismaPg({ connectionString });
const prisma = new PrismaClient({ adapter });

// Helper to create JSON responses
function json(data: unknown, status = 200): Response {
 return new Response(JSON.stringify(data, null, 2), {
 status,
 headers: { "Content-Type": "application/json" },
 });
}

// Request handler
async function handler(request: Request): Promise<Response> {
 const url = new URL(request.url);
 const path = url.pathname;
 const method = request.method;

 try {
 // GET /tasks - List all tasks
 if (method === "GET" && path === "/tasks") {
 const tasks = await prisma.task.findMany({
 orderBy: { createdAt: "desc" },
 });
 return json(tasks);
 }

 // POST /tasks - Create a new task
 if (method === "POST" && path === "/tasks") {
 const body = await request.json();
 const task = await prisma.task.create({
 data: {
 title: body.title,
 description: body.description,
 },
 });
 return json(task, 201);
 }

 // GET /tasks/:id - Get a specific task
 const taskMatch = path.match(/^\/tasks\/(\d+)$/);
 if (taskMatch) {
 const id = parseInt(taskMatch[1]);

 if (method === "GET") {
 const task = await prisma.task.findUnique({ where: { id } });
 if (!task) return json({ error: "Task not found" }, 404);
 return json(task);
 }

 // PATCH /tasks/:id - Update a task
 if (method === "PATCH") {
 const body = await request.json();
 const task = await prisma.task.update({
 where: { id },
 data: body,
 });
 return json(task);
 }

 // DELETE /tasks/:id - Delete a task
 if (method === "DELETE") {
 await prisma.task.delete({ where: { id } });
 return json({ message: "Task deleted" });
 }
 }

 // GET / - API info
 if (method === "GET" && path === "/") {
 return json({
 name: "Prisma + Deno Task API",
 version: "1.0.0",
 endpoints: {
 "GET /tasks": "List all tasks",
 "POST /tasks": "Create a task",
 "GET /tasks/:id": "Get a task",
 "PATCH /tasks/:id": "Update a task",
 "DELETE /tasks/:id": "Delete a task",
 },
 });
 }

 return json({ error: "Not found" }, 404);
 } catch (error) {
 console.error(error);
 return json({ error: "Internal server error" }, 500);
 }
}

// Start the server
Deno.serve({ port: 8000 }, handler);
```

This creates a full CRUD API with the following endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | API info |
| GET | `/tasks` | List all tasks |
| POST | `/tasks` | Create a new task |
| GET | `/tasks/:id` | Get a specific task |
| PATCH | `/tasks/:id` | Update a task |
| DELETE | `/tasks/:id` | Delete a task |

## 6. Test your application locally

Start the development server:

```bash
deno task dev
```

Test the API with curl:

```bash
# Get API info
curl http://localhost:8000/

# Create a task
curl -X POST http://localhost:8000/tasks \
 -H "Content-Type: application/json" \
 -d '{"title": "Learn Prisma", "description": "Complete the Deno guide"}'

# List all tasks
curl http://localhost:8000/tasks

# Update a task (mark as completed)
curl -X PATCH http://localhost:8000/tasks/1 \
 -H "Content-Type: application/json" \
 -d '{"completed": true}'

# Delete a task
curl -X DELETE http://localhost:8000/tasks/1
```

You should see JSON responses for each request. The task ID will increment with each new task created.

## 7. Create a GitHub repository

You need a GitHub repository to deploy to Deno Deploy.

Create a `.gitignore` file:

```text title=".gitignore"
.env
node_modules/
generated/
deno.lock
```

Initialize and push your repository:

```bash
git init -b main
git remote add origin https://github.com/<username>/prisma-deno-deploy
git add .
git commit -m "Initial commit"
git push -u origin main
```

## 8. Deploy to Deno Deploy

1. Go to [https://dash.deno.com/](https://dash.deno.com/)
2. Click **New Project** and select your GitHub repository
3. Configure the deployment:
 - **Framework preset**: No Preset
 - **Install command**: `deno install`
 - **Build command**: `deno run -A npm:prisma generate`
 - **Entrypoint**: `index.ts`
4. Click **Create & Deploy**

The first deployment will fail because you need to add the database connection string.

### Add environment variables

1. Go to your project's **Settings** > **Environment Variables**
2. Add a new variable:
 - **Key**: `DATABASE_URL`
 - **Value**: Your Prisma Postgres connection string (copy from your `.env` file)
3. Click **Save**

Trigger a new deployment by clicking **Redeploy** or pushing a new commit.

## 9. Test your deployed API

Once deployed, test your API at your Deno Deploy URL:

```bash
# Replace with your actual Deno Deploy URL
curl https://your-project.deno.dev/

# Create a task
curl -X POST https://your-project.deno.dev/tasks \
 -H "Content-Type: application/json" \
 -d '{"title": "Deploy to production"}'

# List tasks
curl https://your-project.deno.dev/tasks
```

## Summary

You successfully deployed a REST API to Deno Deploy using:

- **Deno** as the runtime with native TypeScript support
- **Prisma ORM** with the Postgres adapter for type-safe database access
- **Prisma Postgres** as the managed database

Your project structure should look like this:

```
prisma-deno-deploy/
├── deno.json
├── index.ts
├── prisma/
│ └── schema.prisma
├── prisma.config.ts
├── generated/
│ └── prisma/
│ └── ...
└── .env
```

### Next steps

- Add authentication using [Deno KV](https://docs.deno.com/deploy/reference/deno_kv/) for sessions
- Add request validation with [Zod](https://zod.dev/)
- Explore [Prisma Client extensions](/orm/prisma-client/client-extensions) for custom functionality
- Set up [Prisma Migrate](/orm/prisma-migrate) for schema versioning in production

---
title: Deploy to Vercel Edge Functions & Middleware
description: Learn the things you need to know in order to deploy an Edge function that uses Prisma Client for talking to a database
url: /orm/prisma-client/deployment/edge/deploy-to-vercel
metaTitle: Deploy to Vercel Edge Functions & Middleware
metaDescription: Learn the things you need to know in order to deploy an Edge function that uses Prisma Client for talking to a database.
---

This page covers everything you need to know to deploy an app that uses Prisma Client for talking to a database in [Vercel Edge Middleware](https://vercel.com/docs/functions/edge-middleware) or a [Vercel Function](https://vercel.com/docs/functions) deployed to the [Vercel Edge Runtime](https://vercel.com/docs/functions/runtimes/edge-runtime).

<details>
<summary>Questions answered in this page</summary>

- How to deploy Prisma on Vercel Edge?
- Which database drivers are supported?
- How to configure env and postinstall?

</details>

Vercel supports both Node.js and edge runtimes for Vercel Functions. The Node.js runtime is the default and recommended for most use cases.

:::note

By default, Vercel Functions use the Node.js runtime. You can explicitly set the runtime if needed:

```typescript
export const runtime = "edge"; // 'nodejs' is the default

export function GET(request: Request) {
 return new Response(`I am a Vercel Function!`, {
 status: 200,
 });
}
```

:::

## General considerations when deploying to Vercel Edge Functions & Edge Middleware

### Using Prisma Postgres

You can use Prisma Postgres in Vercel's edge runtime. Follow this guide for an end-to-end tutorial on [deploying an application to Vercel using Prisma Postgres](/guides/frameworks/nextjs).

### Using an edge-compatible driver

Vercel's Edge Runtime currently only supports a limited set of database drivers:

- [Neon Serverless](https://neon.tech/docs/serverless/serverless-driver) uses HTTP to access the database (also compatible with [Vercel Postgres](https://vercel.com/docs/storage/vercel-postgres))
- [PlanetScale Serverless](https://planetscale.com/docs/tutorials/planetscale-serverless-driver) uses HTTP to access the database
- [`@libsql/client`](https://github.com/tursodatabase/libsql-client-ts) is used to access Turso databases

Note that [`node-postgres`](https://node-postgres.com/) (`pg`) is currently _not_ supported on Vercel Edge Functions.

When deploying a Vercel Edge Function that uses Prisma ORM, you need to use one of these [edge-compatible drivers](/orm/prisma-client/deployment/edge/overview#edge-compatibility-of-database-drivers) and its respective [driver adapter](/orm/core-concepts/supported-databases/database-drivers#driver-adapters) for Prisma ORM.

:::note

If your application uses PostgreSQL, we recommend using [Prisma Postgres](/postgres). It is fully supported on edge runtimes and does not require a specialized edge-compatible driver. For other databases, [Prisma Accelerate](/accelerate) extends edge compatibility so you can connect to _any_ database from _any_ edge function provider.

:::

### Setting your database connection URL as an environment variable

First, ensure that your `datasource` block in your Prisma schema is configured correctly. Database connection URLs are configured in `prisma.config.ts`:

```prisma
datasource db {
 provider = "postgresql" // this might also be `mysql` or another value depending on your database
}
```

```ts title="prisma.config.ts"
import "dotenv/config";
import { defineConfig, env } from "prisma/config";

export default defineConfig({
 schema: "prisma/schema.prisma",
 datasource: {
 url: env("DATABASE_URL"),
 },
});
```

#### Development

When in **development**, you can configure your database connection via the `DATABASE_URL` environment variable (e.g. [using `.env` files](/orm/more/dev-environment/environment-variables)).

#### Production

When deploying your Edge Function to **production**, you'll need to set the database connection using the `vercel` CLI:

```npm
npx vercel env add DATABASE_URL
```

This command is interactive and will ask you to select environments and provide the value for the `DATABASE_URL` in subsequent steps.

Alternatively, you can configure the environment variable [via the UI](https://vercel.com/docs/projects/environment-variables#creating-environment-variables) of your project in the Vercel Dashboard.

### Generate Prisma Client in `postinstall` hook

In your `package.json`, you should add a `"postinstall"` section as follows:

```js title="package.json" showLineNumbers
{
 // ...,
 "postinstall": "prisma generate"
}
```

### Size limits on free accounts

Vercel has a [size limit of 1 MB on free accounts](https://vercel.com/docs/functions/limitations). If your application bundle with Prisma ORM exceeds that size, we recommend upgrading to a paid account or using Prisma Accelerate to deploy your application.

## Database-specific considerations & examples

This section provides database-specific instructions for deploying a Vercel Edge Functions with Prisma ORM.

### Prerequisites

As a prerequisite for the following section, you need to have a Vercel Edge Function (which typically comes in the form of a Next.js API route) running locally and the Prisma and Vercel CLIs installed.

If you don't have that yet, you can run these commands to set up a Next.js app from scratch (following the instructions of the [Vercel Functions Quickstart](https://vercel.com/docs/functions/quickstart)):

```npm
npm install -g vercel
```

```npm
npx create-next-app@latest
```

```npm
npm install prisma --save-dev && npm install @prisma/client
```

```npm
npx prisma init --output ../app/generated/prisma
```

We'll use the default `User` model for the example below:

```prisma
model User {
 id Int @id @default(autoincrement())
 email String @unique
 name String?
}
```

### Vercel Postgres

If you are using Vercel Postgres, you need to:

- use the `@prisma/adapter-neon` database adapter because Vercel Postgres uses [Neon](https://neon.tech/) under the hood
- be aware that Vercel by default calls the pooled connection string `POSTGRES_PRISMA_URL` while the direct, non-pooled connection string is exposed as `POSTGRES_URL_NON_POOLING`. Configure `prisma.config.ts` so that the Prisma CLI uses the direct connection string:

 ```ts title="prisma.config.ts"
 import { defineConfig, env } from "prisma/config";

 export default defineConfig({
 schema: "prisma/schema.prisma",
 datasource: {
 url: env("POSTGRES_URL_NON_POOLING"), // direct connection for Prisma CLI
 },
 });
 ```

#### 1. Configure Prisma schema & database connection

:::note

If you don't have a project to deploy, follow the instructions in the [Prerequisites](#prerequisites) to bootstrap a basic Next.js app with Prisma ORM in it.

:::

First, ensure that the database connection is configured properly. Database connection URLs are configured in `prisma.config.ts`:

```prisma title="schema.prisma" showLineNumbers
generator client {
 provider = "prisma-client"
 output = "./generated"
}

datasource db {
 provider = "postgresql"
}
```

```ts title="prisma.config.ts"
import "dotenv/config";
import { defineConfig, env } from "prisma/config";

export default defineConfig({
 schema: "prisma/schema.prisma",
 datasource: {
 url: env("POSTGRES_URL_NON_POOLING"), // direct connection for Prisma CLI
 },
});
```

Next, you need to set the `POSTGRES_PRISMA_URL` and `POSTGRES_URL_NON_POOLING` environment variable to the values of your database connection.

If you ran `npx prisma init`, you can use the `.env` file that was created by this command to set these:

```bash title=".env"
POSTGRES_PRISMA_URL="postgres://user:password@host-pooler.region.postgres.vercel-storage.com:5432/name?pgbouncer=true&connect_timeout=15"
POSTGRES_URL_NON_POOLING="postgres://user:password@host.region.postgres.vercel-storage.com:5432/name"
```

#### 2. Install dependencies

Next, install the required packages:

```npm
npm install @prisma/adapter-neon
```

#### 3. Configure `postinstall` hook

Next, add a new key to the `scripts` section in your `package.json`:

```js title="package.json"
{
 // ...
 "scripts": {
 // ...
 "postinstall": "prisma generate" // [!code ++]
 }
}
```

#### 4. Migrate your database schema (if applicable)

If you ran `npx prisma init` above, you need to migrate your database schema to create the `User` table that's defined in your Prisma schema (if you already have all the tables you need in your database, you can skip this step):

```npm
npx prisma migrate dev --name init
```

#### 5. Use Prisma Client in your Vercel Edge Function to send a query to the database

If you created the project from scratch, you can create a new edge function as follows.

First, create a new API route, e.g. by using these commands:

```bash
mkdir src/app/api
mkdir src/app/api/edge
touch src/app/api/edge/route.ts
```

Here is a sample code snippet that you can use to instantiate `PrismaClient` and send a query to your database in the new `app/api/edge/route.ts` file you just created:

```ts title="app/api/edge/route.ts" showLineNumbers
import { NextResponse } from "next/server";
import { PrismaClient } from "./generated/client";
import { PrismaNeon } from "@prisma/adapter-neon";

export const runtime = "nodejs"; // can also be set to 'edge'

export async function GET(request: Request) {
 const adapter = new PrismaNeon({ connectionString: process.env.POSTGRES_PRISMA_URL });
 const prisma = new PrismaClient({ adapter });

 const users = await prisma.user.findMany();

 return NextResponse.json(users, { status: 200 });
}
```

#### 6. Run the Edge Function locally

Run the app with the following command:

```npm
npm run dev
```

You can now access the Edge Function via `http://localhost:3000/api/edge`.

#### 7. Set the `POSTGRES_PRISMA_URL` environment variable and deploy the Edge Function

Run the following command to deploy your project with Vercel:

```npm
npx vercel deploy
```

Note that once the project was created on Vercel, you will need to set the `POSTGRES_PRISMA_URL` environment variable (and if this was your first deploy, it likely failed). You can do this either via the Vercel UI or by running the following command:

```npm
npx vercel env add POSTGRES_PRISMA_URL
```

At this point, you can get the URL of the deployed application from the Vercel Dashboard and access the edge function via the `/api/edge` route.

### PlanetScale

If you are using a PlanetScale database, you need to:

- use the `@prisma/adapter-planetscale` database adapter (learn more [here](/orm/core-concepts/supported-databases/mysql#planetscale))

#### 1. Configure Prisma schema & database connection

:::note

If you don't have a project to deploy, follow the instructions in the [Prerequisites](#prerequisites) to bootstrap a basic Next.js app with Prisma ORM in it.

:::

First, ensure that the database connection is configured properly. Database connection URLs are configured in `prisma.config.ts`:

```prisma title="schema.prisma" showLineNumbers
generator client {
 provider = "prisma-client"
 output = "./generated"
}

datasource db {
 provider = "mysql"
 relationMode = "prisma" // required for PlanetScale (as by default foreign keys are disabled)
}
```

```ts title="prisma.config.ts"
import "dotenv/config";
import { defineConfig, env } from "prisma/config";

export default defineConfig({
 schema: "prisma/schema.prisma",
 datasource: {
 url: env("DATABASE_URL"),
 },
});
```

Next, you need to set the `DATABASE_URL` environment variable in your `.env` file that's used both by Prisma and Next.js to read your env vars:

```bash title=".env"
DATABASE_URL="mysql://32qxa2r7hfl3102wrccj:password@us-east.connect.psdb.cloud/demo-cf-worker-ps?sslaccept=strict"
```

#### 2. Install dependencies

Next, install the required packages:

```npm
npm install @prisma/adapter-planetscale
```

#### 3. Configure `postinstall` hook

Next, add a new key to the `scripts` section in your `package.json`:

```js title="package.json"
{
 // ...
 "scripts": {
 // ...
 "postinstall": "prisma generate" // [!code ++]
 }
}
```

#### 4. Migrate your database schema (if applicable)

If you ran `npx prisma init` above, you need to migrate your database schema to create the `User` table that's defined in your Prisma schema (if you already have all the tables you need in your database, you can skip this step):

```npm
npx prisma db push
```

#### 5. Use Prisma Client in an Edge Function to send a query to the database

If you created the project from scratch, you can create a new edge function as follows.

First, create a new API route, e.g. by using these commands:

```bash
mkdir src/app/api
mkdir src/app/api/edge
touch src/app/api/edge/route.ts
```

Here is a sample code snippet that you can use to instantiate `PrismaClient` and send a query to your database in the new `app/api/edge/route.ts` file you just created:

```ts title="app/api/edge/route.ts" showLineNumbers
import { NextResponse } from "next/server";
import { PrismaClient } from "./generated/client";
import { PrismaPlanetScale } from "@prisma/adapter-planetscale";

export const runtime = "nodejs"; // can also be set to 'edge'

export async function GET(request: Request) {
 const adapter = new PrismaPlanetScale({ url: process.env.DATABASE_URL });
 const prisma = new PrismaClient({ adapter });

 const users = await prisma.user.findMany();

 return NextResponse.json(users, { status: 200 });
}
```

#### 6. Run the Edge Function locally

Run the app with the following command:

```npm
npm run dev
```

You can now access the Edge Function via `http://localhost:3000/api/edge`.

#### 7. Set the `DATABASE_URL` environment variable and deploy the Edge Function

Run the following command to deploy your project with Vercel:

```npm
npx vercel deploy
```

Note that once the project was created on Vercel, you will need to set the `DATABASE_URL` environment variable (and if this was your first deploy, it likely failed). You can do this either via the Vercel UI or by running the following command:

```npm
npx vercel env add DATABASE_URL
```

At this point, you can get the URL of the deployed application from the Vercel Dashboard and access the edge function via the `/api/edge` route.

### Neon

If you are using a Neon database, you need to:

- use the `@prisma/adapter-neon` database adapter (learn more [here](/orm/core-concepts/supported-databases/postgresql#using-driver-adapters))

#### 1. Configure Prisma schema & database connection

:::note

If you don't have a project to deploy, follow the instructions in the [Prerequisites](#prerequisites) to bootstrap a basic Next.js app with Prisma ORM in it.

:::

First, ensure that the database connection is configured properly. Database connection URLs are configured in `prisma.config.ts`:

```prisma title="schema.prisma" showLineNumbers
generator client {
 provider = "prisma-client"
 output = "./generated"
}

datasource db {
 provider = "postgresql"
}
```

```ts title="prisma.config.ts"
import "dotenv/config";
import { defineConfig, env } from "prisma/config";

export default defineConfig({
 schema: "prisma/schema.prisma",
 datasource: {
 url: env("DATABASE_URL"),
 },
});
```

Next, you need to set the `DATABASE_URL` environment variable in your `.env` file that's used both by Prisma and Next.js to read your env vars:

```bash title=".env"
DATABASE_URL="postgresql://janedoe:password@ep-nameless-pond-a23b1mdz.eu-central-1.aws.neon.tech/neondb?sslmode=require"
```

#### 2. Install dependencies

Next, install the required packages:

```npm
npm install @prisma/adapter-neon
```

#### 3. Configure `postinstall` hook

Next, add a new key to the `scripts` section in your `package.json`:

```js title="package.json"
{
 // ...
 "scripts": {
 // ...
 "postinstall": "prisma generate" // [!code ++]
 }
}
```

#### 4. Migrate your database schema (if applicable)

If you ran `npx prisma init` above, you need to migrate your database schema to create the `User` table that's defined in your Prisma schema (if you already have all the tables you need in your database, you can skip this step):

```npm
npx prisma migrate dev --name init
```

#### 5. Use Prisma Client in an Edge Function to send a query to the database

If you created the project from scratch, you can create a new edge function as follows.

First, create a new API route, e.g. by using these commands:

```bash
mkdir src/app/api
mkdir src/app/api/edge
touch src/app/api/edge/route.ts
```

Here is a sample code snippet that you can use to instantiate `PrismaClient` and send a query to your database in the new `app/api/edge/route.ts` file you just created:

```ts title="app/api/edge/route.ts" showLineNumbers
import { NextResponse } from "next/server";
import { PrismaClient } from "./generated/client";
import { PrismaNeon } from "@prisma/adapter-neon";

export const runtime = "nodejs"; // can also be set to 'edge'

export async function GET(request: Request) {
 const adapter = new PrismaNeon({ connectionString: process.env.DATABASE_URL });
 const prisma = new PrismaClient({ adapter });

 const users = await prisma.user.findMany();

 return NextResponse.json(users, { status: 200 });
}
```

#### 6. Run the Edge Function locally

Run the app with the following command:

```npm
npm run dev
```

You can now access the Edge Function via `http://localhost:3000/api/edge`.

#### 7. Set the `DATABASE_URL` environment variable and deploy the Edge Function

Run the following command to deploy your project with Vercel:

```npm
npx vercel deploy
```

Note that once the project was created on Vercel, you will need to set the `DATABASE_URL` environment variable (and if this was your first deploy, it likely failed). You can do this either via the Vercel UI or by running the following command:

```npm
npx vercel env add DATABASE_URL
```

At this point, you can get the URL of the deployed application from the Vercel Dashboard and access the edge function via the `/api/edge` route.

## Using Prisma ORM with Vercel Fluid

[Fluid compute](https://vercel.com/fluid) is a compute model from Vercel that combines the flexibility of serverless with the stability of servers, making it ideal for dynamic workloads such as streaming data and AI APIs. Vercel's Fluid compute [supports both edge and Node.js runtimes](https://vercel.com/docs/fluid-compute#available-runtime-support). A common challenge in traditional serverless platforms is leaked database connections when functions are suspended and pools can't close idle connections. Fluid provides [`attachDatabasePool`](https://vercel.com/blog/the-real-serverless-compute-to-database-connection-problem-solved) to ensure idle connections are released before a function is suspended.

Use `attachDatabasePool` together with [Prisma's driver adapters](/orm/core-concepts/supported-databases/database-drivers) to safely manage connections in Fluid:

```ts
import { Pool } from "pg";
import { attachDatabasePool } from "@vercel/functions";
import { PrismaPg } from "@prisma/adapter-pg";
import { PrismaClient } from "./generated/client";

const pool = new Pool({ connectionString: process.env.POSTGRES_URL });

attachDatabasePool(pool);

const prisma = new PrismaClient({
 adapter: new PrismaPg(pool),
});
```

---
title: Deploying edge functions with Prisma ORM
description: Learn how to deploy your Prisma-backed apps to edge functions like Cloudflare Workers or Vercel Edge Functions
url: /orm/prisma-client/deployment/edge/overview
metaTitle: 'Overview: Deploy Prisma ORM at the Edge'
metaDescription: Learn how to deploy your Prisma-backed apps to edge functions like Cloudflare Workers or Vercel Edge Functions
---

You can deploy an application that uses Prisma ORM to the edge. Depending on which edge function provider and which database you use, there are different considerations and things to be aware of.

<details>
<summary>Questions answered in this page</summary>

- Which database drivers work on edge?
- How do driver adapters affect connections?
- When to use Prisma Postgres or Accelerate?

</details>

Here is a brief overview of all the edge function providers that are currently supported by Prisma ORM:

| Provider / Product | Supported natively with Prisma ORM | Supported with Prisma Postgres (and Prisma Accelerate) |
| ---------------------- | ------------------------------------------------------- | ------------------------------------------------------ |
| Vercel Edge Functions | ✅ (Preview; only compatible drivers) | ✅ |
| Vercel Edge Middleware | ✅ (Preview; only compatible drivers) | ✅ |
| Cloudflare Workers | ✅ (Preview; only compatible drivers) | ✅ |
| Cloudflare Pages | ✅ (Preview; only compatible drivers) | ✅ |
| Deno Deploy | [Not yet](https://github.com/prisma/prisma/issues/2452) | ✅ |

Deploying edge functions that use Prisma ORM on Cloudflare and Vercel is currently in [Preview](/orm/more/releases#preview).

## Edge-compatibility of database drivers

### Why are there limitations around database drivers in edge functions?

Edge functions typically don't use the standard Node.js runtime. For example, Vercel Edge Functions and Cloudflare Workers are running code in [V8 isolates](https://v8docs.nodesource.com/node-0.8/d5/dda/classv8_1_1_isolate.html). Deno Deploy is using the [Deno](https://deno.com/) JavaScript runtime. As a consequence, these edge functions only have access to a small subset of the standard Node.js APIs and also have constrained computing resources (CPU and memory).

In particular, the constraint of not being able to freely open TCP connections makes it difficult to talk to a traditional database from an edge function. While Cloudflare has introduced a [`connect()`](https://developers.cloudflare.com/workers/runtime-apis/tcp-sockets/) API that enables limited TCP connections, this still only enables database access using specific database drivers that are compatible with that API.

:::note

We recommend using [Prisma Postgres](/postgres). It is fully supported on edge runtimes and does not require a specialized edge-compatible driver.

:::

### Which database drivers are edge-compatible?

Here is an overview of the different database drivers and their compatibility with different edge function offerings:

- [Neon Serverless](https://neon.tech/docs/serverless/serverless-driver) uses HTTP to access the database. It works with Cloudflare Workers and Vercel Edge Functions.
- [PlanetScale Serverless](https://planetscale.com/docs/tutorials/planetscale-serverless-driver) uses HTTP to access the database. It works with Cloudflare Workers and Vercel Edge Functions.
- [`node-postgres`](https://node-postgres.com/) (`pg`) uses Cloudflare's `connect()` (TCP) to access the database. It is only compatible with Cloudflare Workers, not with Vercel Edge Functions.
- [`@libsql/client`](https://github.com/tursodatabase/libsql-client-ts) is used to access Turso databases. It works with Cloudflare Workers and Vercel Edge Functions.
- [Cloudflare D1](https://developers.cloudflare.com/d1/) is used to access D1 databases. It is only compatible with Cloudflare Workers, not with Vercel Edge Functions.
- [Prisma Postgres](/postgres) is used to access a PostgreSQL database built on bare-metal using unikernels. It is supported on both Cloudflare Workers and Vercel.

There's [also work being done](https://github.com/sidorares/node-mysql2/pull/2289) on the `node-mysql2` driver which will enable access to traditional MySQL databases from Cloudflare Workers and Pages in the future as well.

You can use all of these drivers with Prisma ORM using the respective [driver adapters](/orm/core-concepts/supported-databases/database-drivers).

Depending on which deployment provider and database/driver you use, there may be special considerations. Please take a look at the deployment docs for your respective scenario to make sure you can deploy your application successfully:

- Cloudflare
 - [PostgreSQL (traditional)](/orm/prisma-client/deployment/edge/deploy-to-cloudflare#postgresql-traditional)
 - [PlanetScale](/orm/prisma-client/deployment/edge/deploy-to-cloudflare#planetscale)
 - [Neon](/orm/prisma-client/deployment/edge/deploy-to-cloudflare#neon)
 - [Cloudflare D1](/guides/deployment/cloudflare-d1)
 - [Prisma Postgres](https://developers.cloudflare.com/workers/tutorials/using-prisma-postgres-with-workers)
- Vercel
 - [Vercel Postgres](/orm/prisma-client/deployment/edge/deploy-to-vercel#vercel-postgres)
 - [Neon](/orm/prisma-client/deployment/edge/deploy-to-vercel#neon)
 - [PlanetScale](/orm/prisma-client/deployment/edge/deploy-to-vercel#planetscale)
 - [Prisma Postgres](/guides/frameworks/nextjs)

If you want to deploy an app using Turso, you can follow the instructions [here](/orm/core-concepts/supported-databases/sqlite#using-driver-adapters).

---
title: Deploy to AWS Lambda
description: 'Learn how to deploy your Prisma ORM-backed applications to AWS Lambda with AWS SAM, Serverless Framework, or SST'
url: /orm/prisma-client/deployment/serverless/deploy-to-aws-lambda
metaTitle: Deploy your application using Prisma ORM to AWS Lambda
metaDescription: 'Learn how to deploy your Prisma ORM-backed applications to AWS Lambda with AWS SAM, Serverless Framework, or SST'
---

:::info[Quick summary]
This guide explains how to avoid common issues when deploying a project using Prisma ORM to [AWS Lambda](https://aws.amazon.com/lambda/).
:::

<details>
<summary>Questions answered in this page</summary>

- How to deploy Prisma to AWS Lambda?
- Which binaryTargets should I configure?
- How to handle connection pooling on Lambda?

</details>

While a deployment framework is not required to deploy to AWS Lambda, this guide covers deploying with:

- [AWS Serverless Application Model (SAM)](https://aws.amazon.com/serverless/sam/) is an open-source framework from AWS that can be used in the creation of serverless applications. AWS SAM includes the [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-reference.html#serverless-sam-cli), which you can use to build, test, and deploy your application.

- [Serverless Framework](https://www.serverless.com/framework) provides a CLI that helps with workflow automation and AWS resource provisioning. While Prisma ORM works well with the Serverless Framework "out of the box", there are a few improvements that can be made within your project to ensure a smooth deployment and performance. There is also additional configuration that is needed if you are using the [`serverless-webpack`](https://www.npmjs.com/package/serverless-webpack) or [`serverless-bundle`](https://www.npmjs.com/package/serverless-bundle) libraries.

- [SST](https://sst.dev/) provides tools that make it easy for developers to define, test, debug, and deploy their applications. Prisma ORM works well with SST but must be configured so that your schema is correctly packaged by SST.

## General considerations when deploying to AWS Lambda

This section covers changes you will need to make to your application, regardless of framework. After following these steps, follow the steps for your framework.

- [Deploying with AWS SAM](#deploying-with-aws-sam)
- [Deploying with the Serverless Framework](#deploying-with-the-serverless-framework)
- [Deploying with SST](#deploying-with-sst)

### Connection pooling

In a Function as a Service (FaaS) environment, each function invocation typically creates a new database connection. Unlike a continuously running Node.js server, these connections aren't maintained between executions. For better performance in serverless environments, implement connection pooling to reuse existing database connections rather than creating new ones for each function call.

You can use [Prisma Postgres](/postgres), which has built-in connection pooling, to solve this issue. For other solutions, see the [connection management guide for serverless environments](/orm/prisma-client/setup-and-configuration/databases-connections#serverless-environments-faas).

## Deploying with AWS SAM

### Loading environment variables

AWS SAM does not directly support loading values from a `.env` file. You will have to use one of AWS's services to store and retrieve these parameters. [This guide](https://medium.com/bip-xtech/a-practical-guide-to-surviving-aws-sam-d8ab141b3d25) provides a great overview of your options and how to store and retrieve values in Parameters, SSM, Secrets Manager, and more.

## Deploying with the Serverless Framework

### Loading environment variables via a `.env` file

Your functions will need the `DATABASE_URL` environment variable to access the database. The `serverless-dotenv-plugin` will allow you to use your `.env` file in your deployments.

First, make sure that the plugin is installed:

```npm
npm install -D serverless-dotenv-plugin
```

Then, add `serverless-dotenv-plugin` to your list of plugins in `serverless.yml`:

```yaml title="serverless.yml"
plugins:
 - serverless-dotenv-plugin
```

The environment variables in your `.env` file will now be automatically loaded on package or deployment.

```bash
serverless package
```

```bash
Running "serverless" from node_modules
DOTENV: Loading environment variables from .env:
 - DATABASE_URL

Packaging deployment-example-sls for stage dev (us-east-1)
.
```

## Deploying with SST

### Working with environment variables

While SST supports `.env` files, [it is not recommended](https://v2.sst.dev/config#should-i-use-configsecret-or-env-for-secrets). SST recommends using `Config` to access these environment variables in a secure way.

The SST guide [available here](https://v2.sst.dev/config#overview) is a step-by-step guide to get started with `Config`. Assuming you have created a new secret called `DATABASE_URL` and have [bound that secret to your app](https://v2.sst.dev/config#bind-the-config), you can set up `PrismaClient` with the following:

```ts title="prisma.ts" showLineNumbers
import { PrismaClient } from "./generated/client";
import { Config } from "sst/node/config";
import { PrismaPg } from "@prisma/adapter-pg";
const globalForPrisma = global as unknown as { prisma: PrismaClient };

const adapter = new PrismaPg({ connectionString });

export const prisma =
 globalForPrisma.prisma ||
 new PrismaClient({ adapter });

if (process.env.NODE_ENV !== "production") globalForPrisma.prisma = prisma;

export default prisma;
```

---
title: Deploy to Azure Functions
description: Learn how to deploy a Prisma Client based REST API to Azure Functions and connect to an Azure SQL database
url: /orm/prisma-client/deployment/serverless/deploy-to-azure-functions
metaTitle: How to deploy an app using Prisma ORM to Azure Functions
metaDescription: Learn how to deploy a Prisma Client based REST API to Azure Functions and connect to an Azure SQL database
---

This guide explains how to avoid common issues when deploying a Node.js-based function app to Azure using [Azure Functions](https://azure.microsoft.com/en-us/products/functions/).

Azure Functions is a serverless deployment platform. You do not need to maintain infrastructure to deploy your code. With Azure Functions, the fundamental building block is the [function app](https://learn.microsoft.com/en-us/azure/azure-functions/functions-reference?tabs=blob&pivots=programming-language-typescript). A function app provides an execution context in Azure in which your functions run. It is comprised of one or more individual functions that Azure manages, deploys, and scales together. You can organize and collectively manage multiple functions as a single logical unit.

## Prerequisites

- An existing function app project with Prisma ORM

## Things to know

While Prisma ORM works well with Azure functions, there are a few things to take note of before deploying your application.

### Connection pooling

Generally, when you use a FaaS (Function as a Service) environment to interact with a database, every function invocation can result in a new connection to the database. This is not a problem with a constantly running Node.js server. Therefore, it is beneficial to pool DB connections to get better performance. To solve this issue, you can use [Prisma Postgres](/postgres). For other solutions, see the [connection management guide for serverless environments](/orm/prisma-client/setup-and-configuration/databases-connections#serverless-environments-faas).

## Summary

For more insight into Prisma Client's API, explore the function handlers and check out the [Prisma Client API Reference](/orm/reference/prisma-client-reference)

---
title: Deploy to Netlify
description: Learn how to deploy Node.js and TypeScript applications that are using Prisma Client to Netlify
url: /orm/prisma-client/deployment/serverless/deploy-to-netlify
metaTitle: Deploy to Netlify
metaDescription: Learn how to deploy Node.js and TypeScript applications that are using Prisma Client to Netlify.
---

This guide covers the steps you will need to take in order to deploy your application that uses Prisma ORM to [Netlify](https://www.netlify.com/).

Netlify is a cloud platform for continuous deployment, static sites, and serverless functions. Netlify integrates seamlessly with GitHub for automatic deployments upon commits. When you follow the steps below, you will use that approach to create a CI/CD pipeline that deploys your application from a GitHub repository.

## Prerequisites

Before you can follow this guide, you will need to set up your application to begin deploying to Netlify. We recommend the ["Get started with Netlify"](https://docs.netlify.com/get-started/) guide for a quick overview and ["Deploy functions"](https://docs.netlify.com/functions/deploy/?fn-language=ts) for an in-depth look at your deployment options.

## Store environment variables in Netlify

We recommend keeping `.env` files in your `.gitignore` in order to prevent leakage of sensitive connection strings. Instead, you can use the Netlify CLI to [import values into netlify directly](https://docs.netlify.com/environment-variables/get-started/#import-variables-with-the-netlify-cli).

Assuming you have a file like the following:

```text title=".env"
# Connect to DB
DATABASE_URL="postgresql://postgres:__PASSWORD__@__HOST__:__PORT__/__DB_NAME__"
```

You can upload the file as environment variables using the `env:import` command:

```bash no-break-terminal
netlify env:import .env
```

```bash
site: my-very-very-cool-site
---------------------------------------------------------------------------------.
 Imported environment variables |
---------------------------------------------------------------------------------|
 Key | Value |
--------------|------------------------------------------------------------------|
 DATABASE_URL | postgresql://postgres:__PASSWORD__@__HOST__:__PORT__/__DB_NAME__ |
---------------------------------------------------------------------------------'
```

<details>
<summary>If you are not using an `.env` file</summary>

If you are storing your database connection string and other environment variables in a different method, you will need to manually upload your environment variables to Netlify. These options are [discussed in Netlify's documentation](https://docs.netlify.com/environment-variables/get-started/) and one method, uploading via the UI, is described below.

1. Open the Netlify admin UI for the site. You can use Netlify CLI as follows:
 ```bash
 netlify open --admin
 ```
2. Click **Site settings**:
 ![Netlify admin UI](/img/orm/prisma-client/deployment/serverless/images/500-06-deploy-to-netlify-site-settings.png)
3. Navigate to **Build & deploy** in the sidebar on the left and select **Environment**.
4. Click **Edit variables** and create a variable with the key `DATABASE_URL` and set its value to your database connection string.
 ![Netlify environment variables](/img/orm/prisma-client/deployment/serverless/images/500-07-deploy-to-netlify-environment-variables-settings.png)
5. Click **Save**.

</details>

Now start a new Netlify build and deployment so that the new build can use the newly uploaded environment variables.

```bash
netlify deploy
```

You can now test the deployed application.

## Connection pooling

When you use a Function-as-a-Service provider, like Netlify, it is beneficial to pool database connections for performance reasons. This is because every function invocation may result in a new connection to your database which can quickly run out of open connections.

You can use [Prisma Postgres](/postgres), which has built-in connection pooling, to reduce your Prisma Client bundle size, and to avoid cold starts.

For more information on connection management for serverless environments, refer to our [connection management guide](/orm/prisma-client/setup-and-configuration/databases-connections#serverless-environments-faas).

---
title: Deploy to Vercel
description: Learn how to deploy a Next.js application based on Prisma Client to Vercel
url: /orm/prisma-client/deployment/serverless/deploy-to-vercel
metaTitle: Deploy to Vercel
metaDescription: Learn how to deploy a Next.js application based on Prisma Client to Vercel.
---

This guide takes you through the steps to set up and deploy a serverless application that uses Prisma to [Vercel](https://vercel.com/).

Vercel is a cloud platform that hosts static sites, serverless, and edge functions. You can integrate a Vercel project with a GitHub repository to allow you to deploy automatically when you make new commits.

We created an [example application](https://github.com/prisma/deployment-example-vercel) using Next.js you can use as a reference when deploying an application using Prisma to Vercel.

While our examples use Next.js, you can deploy other applications to Vercel. See [Using Express with Vercel](https://vercel.com/guides/using-express-with-vercel) and [Nuxt on Vercel](https://vercel.com/docs/frameworks/nuxt) as examples of other options.

## Build configuration

### Updating Prisma Client during Vercel builds

Vercel will automatically cache dependencies on deployment. For most applications, this will not cause any issues. However, for Prisma ORM, it may result in an outdated version of Prisma Client on a change in your Prisma schema. To avoid this issue, add `prisma generate` to the `postinstall` script of your application:

```json title="package.json" showLineNumbers
{
 ...
 "scripts": {
 "postinstall": "prisma generate" // [!code ++]
 }
 ...
}
```

This will re-generate Prisma Client at build time so that your deployment always has an up-to-date client. Another option to avoid an outdated Prisma Client is to check your client into version control. This way each deployment is guaranteed to include the correct Prisma Client.

:::info
If you see `prisma: command not found` errors during your deployment to Vercel, you are missing `prisma` in your dependencies. By default, `prisma` is a dev dependency and may need to be moved to be a standard dependency.
:::

### CI/CD workflows

In a more sophisticated CI/CD environment, you may additionally want to update the database schema with any migrations you have performed during local development. You can do this using the [`prisma migrate deploy`](/orm/reference/prisma-cli-reference#migrate-deploy) command.

In that case, you could create a custom build command in your `package.json` (e.g. called `vercel-build`) that looks as follows:

```json title="package.json"
{
 ...
 "scripts" {
 "vercel-build": "prisma generate && prisma migrate deploy && next build", // [!code ++]
 }
 ...
}
```

You can invoke this script inside your CI/CD pipeline using the following command:

```npm
npm run vercel-build
```

## Add a separate database for preview deployments

By default, your application will have a single _production_ environment associated with the `main` git branch of your repository. If you open a pull request to change your application, Vercel creates a new _preview_ environment.

Vercel uses the `DATABASE_URL` environment variable you define when you import the project for both the production and preview environments. This causes problems if you create a pull request with a database schema migration because the pull request will change the schema of the production database.

To prevent this, use a _second_ hosted database to handle preview deployments. Once you have that connection string, you can add a `DATABASE_URL` for your preview environment using the Vercel dashboard:

1. Click the **Settings** tab of your Vercel project.

2. Click **Environment variables**.

3. Add an environment variable with a key of `DATABASE_URL` and select only the **Preview** environment option:

 ![Add an environment variable for the preview environment](/img/orm/prisma-client/deployment/serverless/images/300-60-deploy-to-vercel-preview-environment-variable.png)

4. Set the value to the connection string of your second database:

 ```bash
 postgresql://dbUsername:dbPassword@myhost:5432/mydb
 ```

5. Click **Save**.

## Connection pooling

When you use a Function-as-a-Service provider, like Vercel Serverless functions, every invocation may result in a new connection to your database. This can cause your database to quickly run out of open connections and cause your application to stall. For this reason, pooling connections to your database is essential.

You can use [Prisma Postgres](/postgres), which has built-in connection pooling, to reduce your Prisma Client bundle size, and to avoid cold starts.

For more information on connection management for serverless environments, refer to our [connection management guide](/orm/prisma-client/setup-and-configuration/databases-connections#serverless-environments-faas).

## Using Prisma ORM with Vercel Fluid

[Fluid compute](https://vercel.com/fluid) is a compute model from Vercel that combines the flexibility of serverless with the stability of servers, making it ideal for dynamic workloads such as streaming data and AI APIs. Vercel's Fluid compute [supports both edge and Node.js runtimes](https://vercel.com/docs/fluid-compute#available-runtime-support). A common challenge in traditional serverless platforms is leaked database connections when functions are suspended and pools can't close idle connections. Fluid provides [`attachDatabasePool`](https://vercel.com/blog/the-real-serverless-compute-to-database-connection-problem-solved) to ensure idle connections are released before a function is suspended.

Use `attachDatabasePool` together with [Prisma's driver adapters](/orm/core-concepts/supported-databases/database-drivers) to safely manage connections in Fluid:

```ts
import { Pool } from "pg";
import { attachDatabasePool } from "@vercel/functions";
import { PrismaPg } from "@prisma/adapter-pg";
import { PrismaClient } from "./generated/client";

const pool = new Pool({ connectionString: process.env.POSTGRES_URL });

attachDatabasePool(pool);

export const prisma = new PrismaClient({
 adapter: new PrismaPg(pool),
});
```

---
title: Deploy to Fly.io
description: Learn how to deploy a Node.js server that uses Prisma ORM to Fly.io
url: /orm/prisma-client/deployment/traditional/deploy-to-flyio
metaTitle: Deploy a Prisma app to Fly.io
metaDescription: Learn how to deploy a Node.js server that uses Prisma ORM to Fly.io.
---

This guide explains how to deploy a Node.js server that uses Prisma ORM and PostgreSQL to Fly.io.

The [Prisma Render deployment example](https://github.com/prisma/prisma-examples/tree/latest/deployment-platforms/render) contains an Express.js application with REST endpoints and a simple frontend. This app uses Prisma Client to fetch, create, and delete records from its database.
This guide will show you how to deploy the same application, without modification, on Fly.io.

## About Fly.io

[fly.io](https://fly.io/) is a cloud application platform that lets developers easily deploy and scale full-stack applications that start on request near on machines near to users. For this example, it's helpful to know:

- Fly.io lets you deploy long-running, "serverful" full-stack applications in [35 regions around the world](https://fly.io/docs/reference/regions/). By default, applications are configured to to [auto-stop](https://fly.io/docs/launch/autostop-autostart/) when not in use, and auto-start as needed as requests come in.
- Fly.io natively supports a wide variety of [languages and frameworks](https://fly.io/docs/languages-and-frameworks/), including Node.js and Bun. In this guide, we'll use the Node.js runtime.
- Fly.io can [launch apps directly from GitHub](https://fly.io/speedrun). When run from the CLI, `fly launch` will automatically configure applications hosted on GitHub to deploy on push.

## Prerequisites

- Sign up for a [Fly.io](https://fly.io/docs/getting-started/launch/) account

## Get the example code

Download the [example code](https://github.com/prisma/prisma-examples/tree/latest/deployment-platforms/render) to your local machine.

```bash
curl https://codeload.github.com/prisma/prisma-examples/tar.gz/latest | tar -xz --strip=2 prisma-examples-latest/deployment-platforms/render
cd render
```

## Understand the example

Before we deploy the app, let's take a look at the example code.

### Web application

The logic for the Express app is in two files:

- `src/index.js`: The API. The endpoints use Prisma Client to fetch, create, and delete data from the database.
- `public/index.html`: The web frontend. The frontend calls a few of the API endpoints.

### Prisma schema and migrations

The Prisma components of this app are in three files:

- `prisma/schema.prisma`: The data model of this app. This example defines two models, `User` and `Post`. The format of this file follows the [Prisma schema](/orm/prisma-schema/overview).
- `prisma/migrations/<migration name>/migration.sql`: The SQL commands that construct this schema in a PostgreSQL database. You can auto-generate migration files like this one by running [`prisma migrate dev`](/orm/prisma-migrate/understanding-prisma-migrate/mental-model#what-is-prisma-migrate).
- `prisma/seed.js`: defines some test users and postsPrisma, used to [seed the database](/orm/prisma-migrate/workflows/seeding) with starter data.

## Deploy the example

### 1. Run `fly launch` and accept the defaults

That’s it. Your web service will be live at its `fly.dev` URL as soon as the deploy completes. Optionally [scale](https://fly.io/docs/launch/scale-count/) the size, number, and placement of machines as desired. [`fly console`](https://fly.io/docs/flyctl/console/) can be used to ssh into a new or existing machine.

More information can be found on in the [fly.io documentation](https://fly.io/docs/js/prisma/).

---
title: Deploy to Heroku
description: Learn how to deploy a Node.js server that uses Prisma ORM to Heroku
url: /orm/prisma-client/deployment/traditional/deploy-to-heroku
metaTitle: Deploy a Prisma app to Heroku
metaDescription: Learn how to deploy a Node.js server that uses Prisma ORM to Heroku.
---

In this guide, you will set up and deploy a Node.js server that uses Prisma ORM with PostgreSQL to [Heroku](https://www.heroku.com). The application exposes a REST API and uses Prisma Client to handle fetching, creating, and deleting records from a database.

Heroku is a cloud platform as a service (PaaS). In contrast to the popular serverless deployment model, with Heroku, your application is constantly running even if no requests are made to it. This has several benefits due to the connection limits of a PostgreSQL database. For more information, check out the [general deployment documentation](/orm/prisma-client/deployment/deploy-prisma)

Typically Heroku integrates with a Git repository for automatic deployments upon commits. You can deploy to Heroku from a GitHub repository or by pushing your source to a [Git repository that Heroku creates per app](https://devcenter.heroku.com/articles/git). This guide uses the latter approach whereby you push your code to the app's repository on Heroku, which triggers a build and deploys the application.

The application has the following components:

- **Backend**: Node.js REST API built with Express.js with resource endpoints that use Prisma Client to handle database operations against a PostgreSQL database (e.g., hosted on Heroku).
- **Frontend**: Static HTML page to interact with the API.

![Heroku deployment architecture diagram showing a Node.js backend with Prisma Client, static frontend, and PostgreSQL database.](/img/orm/prisma-client/deployment/traditional/images/heroku-architecture.png)

The focus of this guide is showing how to deploy projects using Prisma ORM to Heroku. The starting point will be the [Prisma Heroku example](https://github.com/prisma/prisma-examples/tree/latest/deployment-platforms/heroku), which contains an Express.js server with a couple of preconfigured REST endpoints and a simple frontend.

> **Note:** The various **checkpoints** throughout the guide allowing you to validate whether you performed the steps correctly.

## A note on deploying GraphQL servers to Heroku

While the example uses REST, the same principles apply to a GraphQL server, with the main difference being that you typically have a single GraphQL API endpoint rather than a route for every resource as with REST.

## Prerequisites

- [Heroku](https://www.heroku.com) account.
- [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli) installed.
- Node.js installed.
- PostgreSQL CLI `psql` installed.

> **Note:** Heroku doesn't provide a free plan, so billing information is required.

## Prisma ORM workflow

At the core of Prisma ORM is the [Prisma schema](/orm/prisma-schema/overview) – a declarative configuration where you define your data model and other Prisma ORM-related configuration. The Prisma schema is also a single source of truth for both Prisma Client and Prisma Migrate.

In this guide, you will use [Prisma Migrate](/orm/prisma-migrate) to create the database schema. Prisma Migrate is based on the Prisma schema and works by generating `.sql` migration files that are executed against the database.

Migrate comes with two primary workflows:

- Creating migrations and applying during local development with `prisma migrate dev`
- Applying generated migration to production with `prisma migrate deploy`

For brevity, the guide does not cover how migrations are created with `prisma migrate dev`. Rather, it focuses on the production workflow and uses the Prisma schema and SQL migration that are included in the example code.

You will use Heroku's [release phase](https://devcenter.heroku.com/articles/release-phase) to run the `prisma migrate deploy` command so that the migrations are applied before the application starts.

To learn more about how migrations are created with Prisma Migrate, check out the [start from scratch guide](/prisma-orm/quickstart/postgresql)

## 1. Download the example and install dependencies

Open your terminal and navigate to a location of your choice. Create the directory that will hold the application code and download the example code:

```bash
mkdir prisma-heroku
cd prisma-heroku
curl https://codeload.github.com/prisma/prisma-examples/tar.gz/latest | tar -xz --strip=3 prisma-examples-latest/deployment-platforms/heroku
```

**Checkpoint:** `ls -1` should show:

```bash
ls -1
Procfile
README.md
package.json
prisma
public
src
```

Install the dependencies:

```npm
npm install
```

> **Note:** The `Procfile` tells Heroku the command needed to start the application, i.e. `npm start`, and the command to run during the release phase, i.e., `npx prisma migrate deploy`

## 2. Create a Git repository for the application

In the previous step, you downloaded the code. In this step, you will create a repository from the code so that you can push it to Heroku for deployment.

To do so, run `git init` from the source code folder:

```bash
git init
> Initialized empty Git repository in /Users/alice/prisma-heroku/.git/
```

To use the `main` branch as the default branch, run the following command:

```bash
git branch -M main
```

With the repository initialized, add and commit the files:

```bash
git add .
git commit -m 'Initial commit'
```

**Checkpoint:** `git log -1` should show the commit:

```bash
git log -1
commit 895534590fdd260acee6396e2e1c0438d1be7fed (HEAD -> main)
```

## 3. Heroku CLI login

Make sure you're logged in to Heroku with the CLI:

```bash
heroku login
```

This will allow you to deploy to Heroku from the terminal.

**Checkpoint:** `heroku auth:whoami` should show your username:

```bash
heroku auth:whoami
> your-email
```

## 4. Create a Heroku app

To deploy an application to Heroku, you need to create an app. You can do so with the following command:

```bash
heroku apps:create your-app-name
```

> **Note:** Use a unique name of your choice instead of `your-app-name`.

**Checkpoint:** You should see the URL and the repository for your Heroku app:

```bash
heroku apps:create your-app-name
> Creating ⬢ your-app-name... done
> https://your-app-name.herokuapp.com/ | https://git.heroku.com/your-app-name.git
```

Creating the Heroku app will add the git remote Heroku created to your local repository. Pushing commits to this remote will trigger a deploy.

**Checkpoint:** `git remote -v` should show the Heroku git remote for your application:

```bash
heroku https://git.heroku.com/your-app-name.git (fetch)
heroku https://git.heroku.com/your-app-name.git (push)
```

If you don't see the heroku remote, use the following command to add it:

```bash
heroku git:remote --app your-app-name
```

## 5. Add a PostgreSQL database to your application

Heroku allows your to provision a PostgreSQL database as part of an application.

Create the database with the following command:

```bash
heroku addons:create heroku-postgresql:hobby-dev
```

**Checkpoint:** To verify the database was created you should see the following:

```bash
Creating heroku-postgresql:hobby-dev on ⬢ your-app-name... free
Database has been created and is available
 ! This database is empty. If upgrading, you can transfer
 ! data from another database with pg:copy
Created postgresql-parallel-73780 as DATABASE_URL
```

> **Note:** Heroku automatically sets the `DATABASE_URL` environment variable when the app is running on Heroku. Prisma ORM uses this environment variable because it's declared in the _datasource_ block of the Prisma schema (`prisma/schema.prisma`) with `env("DATABASE_URL")`.

## 6. Push to deploy

Deploy the app by pushing the changes to the Heroku app repository:

```bash
git push heroku main
```

This will trigger a build and deploy your application to Heroku. Heroku will also run the `npx prisma migrate deploy` command which executes the migrations to create the database schema before deploying the app (as defined in the `release` step of the `Procfile`).

**Checkpoint:** `git push` will emit the logs from the build and release phase and display the URL of the deployed app:

```bash
remote: -----> Launching...
remote: ! Release command declared: this new release will not be available until the command succeeds.
remote: Released v5
remote: https://your-app-name.herokuapp.com/ deployed to Heroku
remote:
remote: Verifying deploy... done.
remote: Running release command...
remote:
remote: Prisma schema loaded from prisma/schema.prisma
remote: Datasource "db": PostgreSQL database "your-db-name", schema "public" at "your-db-host.compute-1.amazonaws.com:5432"
remote:
remote: 1 migration found in prisma/migrations
remote:
remote: The following migration have been applied:
remote:
remote: migrations/
remote: └─ 20210310152103_init/
remote: └─ migration.sql
remote:
remote: All migrations have been successfully applied.
remote: Waiting for release.... done.
```

> **Note:** Heroku will also set the `PORT` environment variable to which your application is bound.

## 7. Test your deployed application

You can use the static frontend to interact with the API you deployed via the preview URL.

Open up the preview URL in your browser, the URL should like this: `https://APP_NAME.herokuapp.com`. You should see the following:

![Deployed Prisma app frontend in browser showing Check API status, Seed data, and Load feed buttons.](/img/orm/prisma-client/deployment/traditional/images/heroku-deployed.png)

The buttons allow you to make requests to the REST API and view the response:

- **Check API status**: Will call the REST API status endpoint that returns `{"up":true}`.
- **Seed data**: Will seed the database with a test `user` and `post`. Returns the created users.
- **Load feed**: Will load all `users` in the database with their related `profiles`.

For more insight into Prisma Client's API, look at the route handlers in the `src/index.js` file.

You can view the application's logs with the `heroku logs --tail` command:

```bash
2020-07-07T14:39:07.396544+00:00 app[web.1]:
2020-07-07T14:39:07.396569+00:00 app[web.1]: > prisma-heroku@1.0.0 start /app
2020-07-07T14:39:07.396569+00:00 app[web.1]: > node src/index.js
2020-07-07T14:39:07.396570+00:00 app[web.1]:
2020-07-07T14:39:07.657505+00:00 app[web.1]: 🚀 Server ready at: http://localhost:12516
2020-07-07T14:39:07.657526+00:00 app[web.1]: ⭐️ See sample requests: http://pris.ly/e/ts/rest-express#3-using-the-rest-api
2020-07-07T14:39:07.842546+00:00 heroku[web.1]: State changed from starting to up
```

## Heroku specific notes

There are some implementation details relating to Heroku that this guide addresses and are worth reiterating:

- **Port binding**: web servers bind to a port so that they can accept connections. When deploying to Heroku The `PORT` environment variable is set by Heroku. Ensure you bind to `process.env.PORT` so that your application can accept requests once deployed. A common pattern is to try binding to try `process.env.PORT` and fallback to a preset port as follows:

```js
const PORT = process.env.PORT || 3000;
const server = app.listen(PORT, () => {
 console.log(`app running on port ${PORT}`);
});
```

- **Database URL**: As part of Heroku's provisioning process, a `DATABASE_URL` config var is added to your app’s configuration. This contains the URL your app uses to access the database. Ensure that your `schema.prisma` file uses `env("DATABASE_URL")` so that Prisma Client can successfully connect to the database.

* **Disable SSL certificate validation**: When using the `PrismaPg` adapter with `DATABASE_URL`, make sure to disable SSL certificate validation in line with [Heroku's guidelines](https://devcenter.heroku.com/articles/connecting-heroku-postgres#connecting-in-node-js). Failing to do so can result in a `P1010 DriverAdapterError: DatabaseAccessDenied` error. You can handle this conditionally based on whether the database is local or hosted on Heroku:

```js
const isSecureDb = !process.env.DATABASE_URL.includes("@127.0.0.1");

export const db = new PrismaClient({
 adapter: new PrismaPg({
 connectionString: process.env.DATABASE_URL + (isSecureDb ? `?sslmode=no-verify` : ""),
 }),
});
```

## Summary

Congratulations! You have successfully deployed a Node.js app with Prisma ORM to Heroku.

You can find the source code for the example in [this GitHub repository](https://github.com/prisma/prisma-examples/tree/latest/deployment-platforms/heroku).

For more insight into Prisma Client's API, look at the route handlers in the `src/index.js` file.

---
title: Deploy to Koyeb
description: Learn how to deploy a Node.js server that uses Prisma ORM to Koyeb Serverless Platform
url: /orm/prisma-client/deployment/traditional/deploy-to-koyeb
metaTitle: Deploy a Prisma ORM app to Koyeb
metaDescription: Learn how to deploy a Node.js server that uses Prisma ORM to Koyeb Serverless Platform.
---

In this guide, you will set up and deploy a Node.js server that uses Prisma ORM with PostgreSQL to [Koyeb](https://www.koyeb.com/). The application exposes a REST API and uses Prisma Client to handle fetching, creating, and deleting records from a database.

Koyeb is a developer-friendly serverless platform to deploy apps globally. The platform lets you seamlessly run Docker containers, web apps, and APIs with git-based deployment, TLS encryption, native autoscaling, a global edge network, and built-in service mesh & discovery.

When using the [Koyeb git-driven deployment](https://www.koyeb.com/docs/build-and-deploy/build-from-git) method, each time you push code changes to a GitHub repository a new build and deployment of the application are automatically triggered on the Koyeb Serverless Platform.
This guide uses the latter approach whereby you push your code to the app's repository on GitHub.

The application has the following components:

- **Backend**: Node.js REST API built with Express.js with resource endpoints that use Prisma Client to handle database operations against a PostgreSQL database (e.g., hosted on Heroku).
- **Frontend**: Static HTML page to interact with the API.

![Koyeb deployment architecture diagram showing a Node.js backend with Prisma Client, static frontend, and PostgreSQL database.](/img/orm/prisma-client/deployment/traditional/images/koyeb-architecture.png)

The focus of this guide is showing how to deploy projects using Prisma ORM to Koyeb. The starting point will be the [Prisma Koyeb example](https://github.com/koyeb/example-prisma), which contains an Express.js server with a couple of preconfigured REST endpoints and a simple frontend.

> **Note:** The various **checkpoints** throughout the guide allow you to validate whether you performed the steps correctly.

## Prerequisites

- Hosted PostgreSQL database and a URL from which it can be accessed, e.g. `postgresql://username:password@your_postgres_db.cloud.com/db_identifier` (you can use Supabase, which offers a [free plan](https://dev.to/prisma/set-up-a-free-postgresql-database-on-supabase-to-use-with-prisma-3pk6)).
- [GitHub](https://github.com) account with an empty public repository we will use to push the code.
- [Koyeb](https://www.koyeb.com) account.
- Node.js installed.

## Prisma ORM workflow

At the core of Prisma ORM is the [Prisma schema](/orm/prisma-schema/overview) – a declarative configuration where you define your data model and other Prisma ORM-related configuration. The Prisma schema is also a single source of truth for both Prisma Client and Prisma Migrate.

In this guide, you will create the database schema with [Prisma Migrate](/orm/prisma-migrate) to create the database schema. Prisma Migrate is based on the Prisma schema and works by generating `.sql` migration files that are executed against the database.

Migrate comes with two primary workflows:

- Creating migrations and applying them during local development with `prisma migrate dev`
- Applying generated migration to production with `prisma migrate deploy`

For brevity, the guide does not cover how migrations are created with `prisma migrate dev`. Rather, it focuses on the production workflow and uses the Prisma schema and SQL migration that are included in the example code.

You will use Koyeb's [build step](https://www.koyeb.com/docs/build-and-deploy/build-from-git#the-buildpack-build-process) to run the `prisma migrate deploy` command so that the migrations are applied before the application starts.

To learn more about how migrations are created with Prisma Migrate, check out the [start from scratch guide](/prisma-orm/quickstart/postgresql)

## 1. Download the example and install dependencies

Open your terminal and navigate to a location of your choice. Create the directory that will hold the application code and download the example code:

```bash
mkdir prisma-on-koyeb
cd prisma-on-koyeb
curl https://github.com/koyeb/example-prisma/tarball/main/latest | tar xz --strip=1
```

**Checkpoint:** Executing the `tree` command should show the following directories and files:

```bash
.
├── README.md
├── package.json
├── prisma
│   ├── migrations
│   │   ├── 20210310152103_init
│   │   │   └── migration.sql
│   │   └── migration_lock.toml
│   └── schema.prisma
├── public
│   └── index.html
└── src
 └── index.js

5 directories, 8 files
```

Install the dependencies:

```npm
npm install
```

## 2. Initialize a Git repository and push the application code to GitHub

In the previous step, you downloaded the code. In this step, you will create a repository from the code so that you can push it to a GitHub repository for deployment.

To do so, run `git init` from the source code folder:

```bash
git init
> Initialized empty Git repository in /Users/edouardb/prisma-on-koyeb/.git/
```

With the repository initialized, add and commit the files:

```bash
git add .
git commit -m 'Initial commit'
```

**Checkpoint:** `git log -1` should show the commit:

```bash
git log -1
commit 895534590fdd260acee6396e2e1c0438d1be7fed (HEAD -> main)
```

Then, push the code to your GitHub repository by adding the remote

```bash
git remote add origin git@github.com:<YOUR_GITHUB_USERNAME>/<YOUR_GITHUB_REPOSITORY_NAME>.git
git push -u origin main
```

## 3. Deploy the application on Koyeb

On the [Koyeb Control Panel](https://app.koyeb.com), click the **Create App** button.

You land on the Koyeb App creation page where you are asked for information about the application to deploy such as the deployment method to use, the repository URL, the branch to deploy, the build and run commands to execute.

Pick GitHub as the deployment method and select the GitHub repository containing your application and set the branch to deploy to `main`.

> **Note:** If this is your first time using Koyeb, you will be prompted to install the Koyeb app in your GitHub account.

In the **Environment variables** section, create a new environment variable `DATABASE_URL` that is type Secret. In the value field, click **Create Secret**, name your secret `prisma-pg-url` and set the PostgreSQL database connection string as the secret value which should look as follows: `postgresql://__USER__:__PASSWORD__@__HOST__/__DATABASE__`.
[Koyeb Secrets](https://www.koyeb.com/docs/reference/secrets) allow you to securely store and retrieve sensitive information like API tokens, database connection strings. They enable you to secure your code by removing hardcoded credentials and let you pass environment variables securely to your applications.

Last, give your application a name and click the **Create App** button.

**Checkpoint:** Open the deployed app by clicking on the screenshot of the deployed app. Once the page loads, click on the **Check API status** button, which should return: `{"up":true}`

![Koyeb app creation interface showing environment variables configuration and Create App button.](/img/orm/prisma-client/deployment/traditional/images/koyeb-app-creation.png)

Congratulations! You have successfully deployed the app to Koyeb.

Koyeb will build and deploy the application. Additional commits to your GitHub repository will trigger a new build and deployment on Koyeb.

**Checkpoint:** Once the build and deployment are completed, you can access your application by clicking the App URL ending with koyeb.app in the Koyeb control panel. Once on the app page loads, Once the page loads, click on the **Check API status** button, which should return: `{"up":true}`

## 4. Test your deployed application

You can use the static frontend to interact with the API you deployed via the preview URL.

Open up the preview URL in your browser, the URL should like this: `https://APP_NAME-ORG_NAME.koyeb.app`. You should see the following:

![Deployed Prisma app frontend in browser showing Check API status, Seed data, and Load feed buttons.](/img/orm/prisma-client/deployment/traditional/images/koyeb-deployed.png)

The buttons allow you to make requests to the REST API and view the response:

- **Check API status**: Will call the REST API status endpoint that returns `{"up":true}`.
- **Seed data**: Will seed the database with a test `user` and `post`. Returns the created users.
- **Load feed**: Will load all `users` in the database with their related `profiles`.

For more insight into Prisma Client's API, look at the route handlers in the `src/index.js` file.

You can view the application's logs clicking the `Runtime logs` tab on your app service from the Koyeb control panel:

```bash
node-72d14691	stdout	> prisma-koyeb@1.0.0 start
node-72d14691	stdout	> node src/index.js
node-72d14691	stdout	🚀 Server ready at: http://localhost:8080
node-72d14691	stdout	⭐️ See sample requests: http://pris.ly/e/ts/rest-express#3-using-the-rest-api
```

## Koyeb specific notes

### Build

By default, for applications using the Node.js runtime, if the `package.json` contains a `build` script, Koyeb automatically executes it after the dependencies installation.
In the example, the `build` script is used to run `prisma generate && prisma migrate deploy && next build`.

### Deployment

By default, for applications using the Node.js runtime, if the `package.json` contains a `start` script, Koyeb automatically executes it to launch the application.
In the example, the `start` script is used to run `node src/index.js`.

### Database migrations and deployments

In the example you deployed, migrations are applied using the `prisma migrate deploy` command during the Koyeb build (as defined in the `build` script in `package.json`).

### Additional notes

In this guide, we kept pre-set values for the region, instance size, and horizontal scaling. You can customize them according to your needs.

> **Note:** The Ports section is used to let Koyeb know which port your application is listening to and properly route incoming HTTP requests. A default `PORT` environment variable is set to `8080` and incoming HTTP requests are routed to the `/` path when creating a new application.
> If your application is listening on another port, you can define another port to route incoming HTTP requests.

## Summary

Congratulations! You have successfully deployed a Node.js app with Prisma ORM to Koyeb.

You can find the source code for the example in [this GitHub repository](https://github.com/koyeb/example-prisma).

For more insight into Prisma Client's API, look at the route handlers in the `src/index.js` file.

---
title: Deploy to Railway
description: Learn how to deploy an app that uses Prisma ORM and Prisma Postgres to Railway
url: /orm/prisma-client/deployment/traditional/deploy-to-railway
metaTitle: Deploy a Prisma app to Railway
metaDescription: Learn how to deploy an app that uses Prisma ORM and Prisma Postgres to Railway.
---

This guide explains how to deploy an app that uses Prisma ORM and Prisma Postgres to [Railway](https://railway.com?utm_medium=integration&utm_source=docs&utm_campaign=prisma). The app exposes a REST API and uses Prisma Client to query a Prisma Postgres database. Your app will run on Railway and connect to a managed Prisma Postgres database.

Railway is a deployment platform that simplifies the software development lifecycle with instant deployments, built-in observability, and effortless scaling. It supports code repositories and container images from popular registries. Railway handles configuration management, environment variables, and provides private networking between services.

To get a pre-wired Next.js Project with Prisma ORM, Prisma Postgres, and Railway, use the [Official Prisma Railway Template](https://railway.com/deploy/prisma-postgres?utm_medium=integration&utm_source=docs&utm_campaign=prisma).

This template automates the provisioning and setup of a Prisma Postgres database, linking it directly to your Next.js application upon deployment, making the entire project ready with just one click.

## Prerequisites

To get started, all you need is:

- A [Railway account](https://railway.com?utm_medium=integration&utm_source=docs&utm_campaign=prisma)
- A GitHub repository with your application code.

:::note
If you don't have a project ready, you can use our [example Prisma project](https://github.com/prisma/prisma-examples/tree/latest/deployment-platforms/railway). It's a simple Hono application that uses Prisma ORM and includes a REST API, a frontend for testing endpoints, and a defined Prisma schema with migrations.
:::

## Deploy your application

### 1. Create a new Railway project

1. Go to the [Railway dashboard](https://railway.com/dashboard?utm_medium=integration&utm_source=docs&utm_campaign=prisma)
2. Click **Create a New Project**
3. Select **GitHub Repo**
4. Click **Configure GitHub App** and authorize Railway
5. Select your repository

Your application is now deploying to Railway, but it won't properly run without a database connection.

In the next section, you'll configure a database and set the `DATABASE_URL` environment variable in Railway.

![Railway deploying application](/img/orm/prisma-client/deployment/traditional/images/railway-deploying.png)

## Configure your database

### 1. Get your database connection string

You'll need a Prisma Postgres connection string. There are two ways to obtain one:

- Create a new database on [Prisma Data Platform](https://console.prisma.io)
- Run `npx create-db` for a temporary database _(no account required)_

### 2. Add the database URL to Railway

1. In your Railway project, open your service
2. Go to the **Variables** tab
3. Click **New Variable**
4. Set the name to `DATABASE_URL`
5. Paste your database connection string as the value
6. Click **Deploy** to redeploy your application with the new environment variable

![Railway environment variables setup](/img/orm/prisma-client/deployment/traditional/images/railway-env-vars.png)

### 3. Access your application

Once the deployment completes with the database URL configured:

1. Navigate to the **Settings** tab
2. Under **Networking**, click **Generate Domain**
3. Your application will be available at the generated URL

![Railway networking settings](/img/orm/prisma-client/deployment/traditional/images/railway-networking.png)

Go to the generated URL and you'll see your deployed app!

If you used the example project, you should see three api endpoints already set up:

- Check API status (`/api`)
- Load feed (`/api/feed`)
- Seed data (`/api/seed`)

![Railway deployed application](/img/orm/prisma-client/deployment/traditional/images/railway-final-product.png)

If you see any errors:

- Wait a minute and refresh
- Ensure `DATABASE_URL` is set
- Check the service logs
- Redeploy

To learn more about the various features Railway offers for your application, visit the [Railway documentation](https://docs.railway.app?utm_medium=integration&utm_source=docs&utm_campaign=prisma).

---
title: Deploy to Render
description: Learn how to deploy a Node.js server that uses Prisma ORM to Render
url: /orm/prisma-client/deployment/traditional/deploy-to-render
metaTitle: Deploy a Prisma app to Render
metaDescription: Learn how to deploy a Node.js server that uses Prisma ORM to Render.
---

This guide explains how to deploy a Node.js server that uses Prisma ORM and PostgreSQL to Render.

<details>
<summary>Questions answered in this page</summary>

- How to deploy Prisma apps on Render?
- How to run migrations before start?
- What Render settings are recommended?

</details>

The [Prisma Render deployment example](https://github.com/prisma/prisma-examples/tree/latest/deployment-platforms/render) contains an Express.js application with REST endpoints and a simple frontend. This app uses Prisma Client to fetch, create, and delete records from its database.

## About Render

[Render](https://render.com) is a cloud application platform that lets developers easily deploy and scale full-stack applications. For this example, it's helpful to know:

- Render lets you deploy long-running, "serverful" full-stack applications. You can configure Render services to [autoscale](https://docs.render.com/scaling) based on CPU and/or memory usage. This is one of several [deployment paradigms](/orm/prisma-client/deployment/deploy-prisma) you can choose from.
- Render natively supports [common runtimes](https://docs.render.com/language-support), including Node.js and Bun. In this guide, we'll use the Node.js runtime.
- Render [integrates with Git repos](https://docs.render.com/github) for automatic deployments upon commits. You can deploy to Render from GitHub, GitLab, or Bitbucket. In this guide, we'll deploy from a Git repository.

## Prerequisites

- Sign up for a [Render](https://render.com) account

## Get the example code

Download the [example code](https://github.com/prisma/prisma-examples/tree/latest/deployment-platforms/render) to your local machine.

```bash
curl https://codeload.github.com/prisma/prisma-examples/tar.gz/latest | tar -xz --strip=2 prisma-examples-latest/deployment-platforms/render
cd render
```

## Understand the example

Before we deploy the app, let's take a look at the example code.

### Web application

The logic for the Express app is in two files:

- `src/index.js`: The API. The endpoints use Prisma Client to fetch, create, and delete data from the database.
- `public/index.html`: The web frontend. The frontend calls a few of the API endpoints.

### Prisma schema and migrations

The Prisma components of this app are in two files:

- `prisma/schema.prisma`: The data model of this app. This example defines two models, `User` and `Post`. The format of this file follows the [Prisma schema](/orm/prisma-schema/overview).
- `prisma/migrations/<migration name>/migration.sql`: The SQL commands that construct this schema in a PostgreSQL database. You can auto-generate migration files like this one by running [`prisma migrate dev`](/orm/prisma-migrate/understanding-prisma-migrate/mental-model#what-is-prisma-migrate).

### Render Blueprint

The `render.yaml` file is a [Render blueprint](https://docs.render.com/infrastructure-as-code). Blueprints are Render's Infrastructure as Code format. You can use a Blueprint to programmatically create and modify services on Render.

A `render.yaml` defines the services that will be spun up on Render by a Blueprint. In this `render.yaml`, we see:

- **A web service that uses a Node runtime**: This is the Express app.
- **A PostgreSQL database**: This is the database that the Express app uses.

The format of this file follows the [Blueprint specification](https://docs.render.com/blueprint-spec).

### How Render deploys work with Prisma Migrate

In general, you want all your database migrations to run before your web app is started. Otherwise, the app may hit errors when it queries a database that doesn't have the expected tables and rows.

You can use the Pre-Deploy Command setting in a Render deploy to run any commands, such as database migrations, before the app is started.

For more details about the Pre-Deploy Command, see [Render's deploy guide](https://docs.render.com/deploys#deploy-steps).

In our example code, the `render.yaml` shows the web service's build command, pre-deploy command, and start command. Notably, `npx prisma migrate deploy` (the pre-deploy command) will run before `npm run start` (the start command).

| **Command** | **Value** |
| :----------------- | :------------------------------- |
| Build Command | `npm install --production=false` |
| Pre-Deploy Command | `npx prisma migrate deploy` |
| Start Command | `npm run start` |

## Deploy the example

### 1. Initialize your Git repository

1. Download [the example code](https://github.com/prisma/prisma-examples/tree/latest/deployment-platforms/render) to your local machine.
2. Create a new Git repository on GitHub, GitLab, or BitBucket.
3. Upload the example code to your new repository.

### 2. Deploy manually

1. In the Render Dashboard, click **New** > **PostgreSQL**. Provide a database name, and select a plan. (The Free plan works for this demo.)
2. After your database is ready, look up its [internal URL](https://docs.render.com/postgresql-creating-connecting#internal-connections).
3. In the Render Dashboard, click **New** > **Web Service** and connect the Git repository that contains the example code.
4. Provide the following values during service creation:

| **Setting** | **Value** |
| :----------------------------------------------------------- | :----------------------------------------------------- |
| Language | `Node` |
| Build Command | `npm install --production=false` |
| Pre-Deploy Command (Note: this may be in the "Advanced" tab) | `npx prisma migrate deploy` |
| Start Command | `npm run start` |
| Environment Variables | Set `DATABASE_URL` to the internal URL of the database |

That’s it. Your web service will be live at its `onrender.com` URL as soon as the build finishes.

### 3. (optional) Deploy with Infrastructure as Code

You can also deploy the example using the Render Blueprint. Follow Render's [Blueprint setup guide] and use the `render.yaml` in the example.

## Bonus: Seed the database

Prisma ORM includes a framework for [seeding the database](/orm/prisma-migrate/workflows/seeding) with starter data. In our example, `prisma/seed.js` defines some test users and posts.

To add these users to the database, we can either:

1. Add the seed script to our Pre-Deploy Command, or
2. Manually run the command on our server via an SSH shell

### Method 1: Pre-Deploy Command

If you manually deployed your Render services:

1. In the Render dashboard, navigate to your web service.
2. Select **Settings**.
3. Set the Pre-Deploy Command to: `npx prisma migrate deploy; npx prisma db seed`

If you deployed your Render services using the Blueprint:

1. In your `render.yaml` file, change the `preDeployCommand` to: `npx prisma migrate deploy; npx prisma db seed`
2. Commit the change to your Git repo.

### Method 2: SSH

Render allows you to SSH into your web service.

1. Follow [Render's SSH guide](https://docs.render.com/ssh) to connect to your server.
2. In the shell, run: `npx prisma db seed`

---
title: Deploy to Sevalla
description: Learn how to deploy a Node.js server that uses Prisma ORM to Sevalla
url: /orm/prisma-client/deployment/traditional/deploy-to-sevalla
metaTitle: Deploy a Prisma app to Sevalla
metaDescription: Learn how to deploy a Node.js server that uses Prisma ORM to Sevalla.
---

This guide explains how to deploy a Node.js server that uses Prisma ORM and PostgreSQL to [Sevalla](https://sevalla.com). The app exposes a REST API and uses Prisma Client to query a PostgreSQL database. Both the app and database will be hosted on Sevalla.

Sevalla is a developer-focused PaaS platform designed to simplify application and server deployment. You can easily host your applications, databases, object storage and static sites.

It supports Git-based deployments, Dockerfiles, and Procfiles and runs on Google Kubernetes Engine with global delivery via Cloudflare.

For this example, it's helpful to know:

- Sevalla supports long-running “serverful” applications with built-in autoscaling and flexible pod sizes.
- Git-based deployments are integrated with GitHub, GitLab, and Bitbucket.
- PostgreSQL, MySQL, MariaDB, Redis, and Valkey databases are natively supported.
- Applications and databases can be securely connected through private networking with automatic environment variable injection.

## Prerequisites

To get started, all you need is:

- A [Sevalla account](https://sevalla.com) (comes with $50 in free credits)
- A GitHub repository with your application code.

> **Note:** If you don't have a project ready, you can use our [example Prisma project](https://github.com/sevalla-templates/express-prisma-demo). It's a simple Express.js application that uses Prisma ORM and includes a REST API, a frontend for testing endpoints, a defined Prisma schema with migrations, and an optional database seeding script.

## Create an app on Sevalla

First, create a new application on your Sevalla dashboard. Click **Applications** on the sidebar and then click the **Create an app** button, as shown below.

![Sevalla app creation interface](/img/orm/prisma-client/deployment/traditional/images/sevalla-app-creation.png)

Next, select your Git repository. Sevalla supports deploying from both private and public repositories. If you're using our example project, there's no need to fork — simply enter the URL `https://github.com/sevalla-templates/express-prisma-demo`.

![Choose repository interface in Sevalla](/img/orm/prisma-client/deployment/traditional/images/sevalla-choose-repository.png)

Choose the appropriate branch (usually `main` or `master`), set your application's name, select the desired deployment region, and pick your pod size. (You can start with 0.5 CPU / 1GB RAM using your free credits)

![Name application interface in Sevalla](/img/orm/prisma-client/deployment/traditional/images/sevalla-name-application.png)

Click **Create**, but skip the deploy step for now, so it does not fail since we have not added the database.

## Set up the database on Sevalla

On your Sevalla dashboard, click **Databases** \> **Add database**. Select **PostgreSQL** (or your preferred configured database type), and then provide a database name, username, and password, or simply use the default generated details.

![Create database interface in Sevalla](/img/orm/prisma-client/deployment/traditional/images/sevalla-create-database.png)

Next, give the database a recognizable name for easy identification within your Sevalla dashboard. Ensure you select the **same region** as your application, choose an appropriate database size, and click **Create**.

Once your database has been created, scroll down to the **Connected Applications** section and click **Add Connection**.

Select the application you created previously and enable "Add environment variables". Make sure the variable name is set to `DATABASE_URL` (if it defaults to `DB_URL`, change it accordingly).

![App internal connection interface in Sevalla](/img/orm/prisma-client/deployment/traditional/images/sevalla-app-internal-connection.png)

Finally, click **Add connection**. This securely links your application to your database and configures the necessary environment variable.

## Trigger deployment

Now that your application and database are connected, return to your application's **Deployment** tab and click **Deploy**.

> **Note:** Sevalla automatically builds your application using Nixpacks, detecting your project's `package.json` and `start` script. If you'd rather use a Dockerfile or Buildpacks, navigate to **Settings** \> **Build**, and click **Update Settings** under the **Build environment** section to customize your build method.

## Seed the database (optional)

If your project includes a seed script (typically located at `prisma/seed.js`), you can populate your database with initial or demo data after deploying your application.

To do this, navigate to your application's **Web Terminal** (under **Applications** \> **\[Your App\]** \> **Web Terminal**) and run the following command:

```npm
npx prisma db seed
```

This is what the web terminal looks like:

![Sevalla web terminal interface](/img/orm/prisma-client/deployment/traditional/images/sevalla-web-terminal.png)

Once this is done, you can manage and interact with your database directly via Sevalla's built-in interactive studio. As shown below, we can see the seeded data:

![Sevalla database studio interface](/img/orm/prisma-client/deployment/traditional/images/sevalla-database-studio.png)

Congratulations\! You've successfully deployed a Node.js application using Prisma ORM to Sevalla. Your app and database are now connected, secure, and production-ready\!

---
title: Query optimization
description: How to identify and optimize query performance with Prisma
url: /orm/prisma-client/queries/advanced/query-optimization-performance
metaTitle: Query optimization
metaDescription: How to identify and optimize query performance with Prisma
---

This page covers identifying and optimizing query performance with Prisma ORM.

## Query Insights

[Query Insights](/query-insights) is built into Prisma Postgres and shows you which queries are slow, how expensive they are, and what to fix. It works out of the box for raw SQL, but to see Prisma ORM operations (model name, action, query shape) you need one extra step.

### Enabling Prisma ORM attribution

Install `@prisma/sqlcommenter-query-insights`:

```bash
npm install @prisma/sqlcommenter-query-insights
```

Then pass it to the `comments` option in your `PrismaClient` constructor:

```ts
import { prismaQueryInsights } from "@prisma/sqlcommenter-query-insights";
import { PrismaClient } from "@prisma/client";

const prisma = new PrismaClient({
 adapter: myAdapter, // driver adapter or Accelerate URL required
 comments: [prismaQueryInsights()],
});
```

This adds a SQL comment to every query containing the model, action, and parameterized query shape. Query Insights uses these annotations to trace SQL back to the exact Prisma call that generated it — even when a single Prisma call produces multiple SQL statements.

#### Let your AI agent handle setup

Copy this prompt into your AI coding assistant:

```
Install and configure @prisma/sqlcommenter-query-insights in my project so I can
see Prisma ORM queries in Query Insights. Docs: https://www.prisma.io/docs/query-insights
```

## Debugging performance issues

Common causes of slow queries:
- Over-fetching data
- Missing indexes
- Not caching repeated queries
- Full table scans

Use [Query Insights](/query-insights) to identify which queries are affected and what to change.

## Using bulk queries

It is generally more performant to read and write large amounts of data in bulk - for example, inserting `50,000` records in batches of `1000` rather than as `50,000` separate inserts. `PrismaClient` supports the following bulk queries:

- [`createMany()`](/orm/reference/prisma-client-reference#createmany)
- [`createManyAndReturn()`](/orm/reference/prisma-client-reference#createmanyandreturn)
- [`deleteMany()`](/orm/reference/prisma-client-reference#deletemany)
- [`updateMany()`](/orm/reference/prisma-client-reference#updatemany)
- [`updateManyAndReturn()`](/orm/reference/prisma-client-reference#updatemanyandreturn)
- [`findMany()`](/orm/reference/prisma-client-reference#findmany)

## Reuse `PrismaClient` or use connection pooling to avoid database connection pool exhaustion

Creating multiple instances of `PrismaClient` can exhaust your database connection pool, especially in serverless or edge environments, potentially slowing down other queries. Learn more in the [serverless challenge](/orm/prisma-client/setup-and-configuration/databases-connections#the-serverless-challenge).

For applications with a traditional server, instantiate `PrismaClient` once and reuse it throughout your app instead of creating multiple instances. For example, instead of:

```ts title="query.ts"
async function getPosts() {
 const prisma = new PrismaClient();
 await prisma.post.findMany();
}

async function getUsers() {
 const prisma = new PrismaClient();
 await prisma.user.findMany();
}
```

Define a single `PrismaClient` instance in a dedicated file and re-export it for reuse:

```ts title="db.ts"
export const prisma = new PrismaClient();
```

Then import the shared instance:

```ts title="query.ts"
import { prisma } from "db.ts";

async function getPosts() {
 await prisma.post.findMany();
}

async function getUsers() {
 await prisma.user.findMany();
}
```

For serverless development environments with frameworks that use HMR (Hot Module Replacement), ensure you properly handle a [single instance of Prisma in development](/orm/more/troubleshooting/nextjs#best-practices-for-using-prisma-client-in-development).

## Solving the n+1 problem

The n+1 problem occurs when looping through query results and performing one additional query **per result**.

### Using `findUnique()` with the fluent API

Prisma's dataloader automatically batches `findUnique()` queries in the same tick. Use the fluent API to return related data:

```ts
// Instead of findMany per user, use:
return context.prisma.user
 .findUnique({ where: { id: parent.id } })
 .posts();
```

### Using JOINs with `relationLoadStrategy`

```ts
const posts = await prisma.post.findMany({
 relationLoadStrategy: "join",
 where: { authorId: parent.id },
});
```

- All criteria of the `where` filter are on scalar fields (unique or non-unique) of the same model you're querying.
- All criteria use the `equal` filter, whether that's via the shorthand or explicit syntax `(where: { field: <val>, field1: { equals: <val> } })`.
- No boolean operators or relation filters are present.

Automatic batching of `findUnique()` is particularly useful in a **GraphQL context**. GraphQL runs a separate resolver function for every field, which can make it difficult to optimize a nested query.

For example - the following GraphQL runs the `allUsers` resolver to get all users, and the `posts` resolver **once per user** to get each user's posts (n+1):

```js
query {
 allUsers {
 id,
 posts {
 id
 }
 }
}
```

The `allUsers` query uses `user.findMany(..)` to return all users:

```ts highlight=7;normal
const Query = objectType({
 name: "Query",
 definition(t) {
 t.nonNull.list.nonNull.field("allUsers", {
 type: "User",
 resolve: (_parent, _args, context) => {
 return context.prisma.user.findMany();
 },
 });
 },
});
```

This results in a single SQL query:

```js
{
 timestamp: 2021-02-19T09:43:06.332Z,
 query: 'SELECT `dev`.`User`.`id`, `dev`.`User`.`email`, `dev`.`User`.`name` FROM `dev`.`User` WHERE 1=1 LIMIT ? OFFSET ?',
 params: '[-1,0]',
 duration: 0,
 target: 'quaint::connector::metrics'
}
```

However, the resolver function for `posts` is then invoked **once per user**. This results in a `findMany()` query **✘ per user** rather than a single `findMany()` to return all posts by all users (expand CLI output to see queries).

```ts highlight=10-13;normal;
const User = objectType({
 name: "User",
 definition(t) {
 t.nonNull.int("id");
 t.string("name");
 t.nonNull.string("email");
 t.nonNull.list.nonNull.field("posts", {
 type: "Post",
 resolve: (parent, _, context) => {
 return context.prisma.post.findMany({
 where: { authorId: parent.id || undefined },
 });
 },
 });
 },
});
```

```json
{
 timestamp: 2021-02-19T09:43:06.343Z,
 query: 'SELECT `dev`.`Post`.`id`, `dev`.`Post`.`createdAt`, `dev`.`Post`.`updatedAt`, `dev`.`Post`.`title`, `dev`.`Post`.`content`, `dev`.`Post`.`published`, `dev`.`Post`.`viewCount`, `dev`.`Post`.`authorId` FROM `dev`.`Post` WHERE `dev`.`Post`.`authorId` = ? LIMIT ? OFFSET ?',
 params: '[1,-1,0]',
 duration: 0,
 target: 'quaint::connector::metrics'
}
{
 timestamp: 2021-02-19T09:43:06.347Z,
 query: 'SELECT `dev`.`Post`.`id`, `dev`.`Post`.`createdAt`, `dev`.`Post`.`updatedAt`, `dev`.`Post`.`title`, `dev`.`Post`.`content`, `dev`.`Post`.`published`, `dev`.`Post`.`viewCount`, `dev`.`Post`.`authorId` FROM `dev`.`Post` WHERE `dev`.`Post`.`authorId` = ? LIMIT ? OFFSET ?',
 params: '[3,-1,0]',
 duration: 0,
 target: 'quaint::connector::metrics'
}
{
 timestamp: 2021-02-19T09:43:06.348Z,
 query: 'SELECT `dev`.`Post`.`id`, `dev`.`Post`.`createdAt`, `dev`.`Post`.`updatedAt`, `dev`.`Post`.`title`, `dev`.`Post`.`content`, `dev`.`Post`.`published`, `dev`.`Post`.`viewCount`, `dev`.`Post`.`authorId` FROM `dev`.`Post` WHERE `dev`.`Post`.`authorId` = ? LIMIT ? OFFSET ?',
 params: '[2,-1,0]',
 duration: 0,
 target: 'quaint::connector::metrics'
}
{
 timestamp: 2021-02-19T09:43:06.348Z,
 query: 'SELECT `dev`.`Post`.`id`, `dev`.`Post`.`createdAt`, `dev`.`Post`.`updatedAt`, `dev`.`Post`.`title`, `dev`.`Post`.`content`, `dev`.`Post`.`published`, `dev`.`Post`.`viewCount`, `dev`.`Post`.`authorId` FROM `dev`.`Post` WHERE `dev`.`Post`.`authorId` = ? LIMIT ? OFFSET ?',
 params: '[4,-1,0]',
 duration: 0,
 target: 'quaint::connector::metrics'
}
{
 timestamp: 2021-02-19T09:43:06.348Z,
 query: 'SELECT `dev`.`Post`.`id`, `dev`.`Post`.`createdAt`, `dev`.`Post`.`updatedAt`, `dev`.`Post`.`title`, `dev`.`Post`.`content`, `dev`.`Post`.`published`, `dev`.`Post`.`viewCount`, `dev`.`Post`.`authorId` FROM `dev`.`Post` WHERE `dev`.`Post`.`authorId` = ? LIMIT ? OFFSET ?',
 params: '[5,-1,0]',
 duration: 0,
 target: 'quaint::connector::metrics'
}
// And so on
```

#### Solution 1: Batching queries with the fluent API

Use `findUnique()` in combination with [the fluent API](/orm/prisma-client/queries/relation-queries#fluent-api) (`.posts()`) as shown to return a user's posts. Even though the resolver is called once per user, the Prisma dataloader in Prisma Client **✔ batches the `findUnique()` queries**.

:::info

It may seem counterintuitive to use a `prisma.user.findUnique(...).posts()` query to return posts instead of `prisma.posts.findMany()` - particularly as the former results in two queries rather than one.

The **only** reason you need to use the fluent API (`user.findUnique(...).posts()`) to return posts is that the dataloader in Prisma Client batches `findUnique()` queries and does not currently [batch `findMany()` queries](https://github.com/prisma/prisma/issues/1477).

When the dataloader batches `findMany()` queries or your query has the `relationStrategy` set to `join`, you no longer need to use `findUnique()` with the fluent API in this way.

:::

```ts highlight=13-18;add|10-12;delete
const User = objectType({
 name: "User",
 definition(t) {
 t.nonNull.int("id");
 t.string("name");
 t.nonNull.string("email");
 t.nonNull.list.nonNull.field("posts", {
 type: "Post",
 resolve: (parent, _, context) => {
 return context.prisma.post.findMany({
 // [!code --]
 where: { authorId: parent.id || undefined }, // [!code --]
 }); // [!code --]
 return context.prisma.user // [!code ++]
 .findUnique({
 // [!code ++]
 where: { id: parent.id || undefined }, // [!code ++]
 }) // [!code ++]
 .posts(); // [!code ++]
 }, // [!code ++]
 });
 },
});
```

```json
{
 timestamp: 2021-02-19T09:59:46.340Z,
 query: 'SELECT `dev`.`User`.`id`, `dev`.`User`.`email`, `dev`.`User`.`name` FROM `dev`.`User` WHERE 1=1 LIMIT ? OFFSET ?',
 params: '[-1,0]',
 duration: 0,
 target: 'quaint::connector::metrics'
}
{
 timestamp: 2021-02-19T09:59:46.350Z,
 query: 'SELECT `dev`.`User`.`id` FROM `dev`.`User` WHERE `dev`.`User`.`id` IN (?,?,?) LIMIT ? OFFSET ?',
 params: '[1,2,3,-1,0]',
 duration: 0,
 target: 'quaint::connector::metrics'
}
{
 timestamp: 2021-02-19T09:59:46.350Z,
 query: 'SELECT `dev`.`Post`.`id`, `dev`.`Post`.`createdAt`, `dev`.`Post`.`updatedAt`, `dev`.`Post`.`title`, `dev`.`Post`.`content`, `dev`.`Post`.`published`, `dev`.`Post`.`viewCount`, `dev`.`Post`.`authorId` FROM `dev`.`Post` WHERE `dev`.`Post`.`authorId` IN (?,?,?) LIMIT ? OFFSET ?',
 params: '[1,2,3,-1,0]',
 duration: 0,
 target: 'quaint::connector::metrics'
}
```

If the `posts` resolver is invoked once per user, the dataloader in Prisma Client groups `findUnique()` queries with the same parameters and selection set. Each group is optimized into a single `findMany()`.

#### Solution 2: Using JOINs to perform queries

You can perform the query with a [database join](/orm/prisma-client/queries/relation-queries#relation-load-strategies-preview) by setting `relationLoadStrategy` to `"join"`, ensuring that only **one** query is executed against the database.

```ts
const User = objectType({
 name: "User",
 definition(t) {
 t.nonNull.int("id");
 t.string("name");
 t.nonNull.string("email");
 t.nonNull.list.nonNull.field("posts", {
 type: "Post",
 resolve: (parent, _, context) => {
 return context.prisma.post.findMany({
 relationLoadStrategy: "join",
 where: { authorId: parent.id || undefined },
 });
 },
 });
 },
});
```

### Avoiding n+1 in loops

Don't loop with separate queries:

```ts
// BAD: n+1 queries
const users = await prisma.user.findMany({});
users.forEach(async (usr) => {
 const posts = await prisma.post.findMany({ where: { authorId: usr.id } });
});
```

Use `include` or `in` filter instead:

```ts
// GOOD: 2 queries with include
const usersWithPosts = await prisma.user.findMany({
 include: { posts: true },
});

// GOOD: 2 queries with in filter
const users = await prisma.user.findMany({});
const posts = await prisma.post.findMany({
 where: { authorId: { in: users.map(u => u.id) } },
});

// BEST: 1 query with join
const posts = await prisma.post.findMany({
 relationLoadStrategy: "join",
 where: { authorId: { in: users.map(u => u.id) } },
});
```

This is not an efficient way to query. Instead, you can:

- Use nested reads ([`include`](/orm/reference/prisma-client-reference#include) ) to return users and related posts
- Use the [`in`](/orm/reference/prisma-client-reference#in) filter
- Set the [`relationLoadStrategy`](/orm/prisma-client/queries/relation-queries#relation-load-strategies-preview) to `"join"`

#### Solving n+1 with `include`

You can use `include` to return each user's posts. This only results in **two** SQL queries - one to get users, and one to get posts. This is known as a [nested read](/orm/prisma-client/queries/relation-queries#nested-reads).

```ts
const usersWithPosts = await prisma.user.findMany({
 include: {
 posts: true,
 },
});
```

```sql
SELECT "public"."User"."id", "public"."User"."email", "public"."User"."name" FROM "public"."User" WHERE 1=1 OFFSET $1
SELECT "public"."Post"."id", "public"."Post"."title", "public"."Post"."authorId" FROM "public"."Post" WHERE "public"."Post"."authorId" IN ($1,$2,$3,$4) OFFSET $5
```

#### Solving n+1 with `in`

If you have a list of user IDs, you can use the `in` filter to return all posts where the `authorId` is `in` that list of IDs:

```ts
const users = await prisma.user.findMany({});

const userIds = users.map((x) => x.id);

const posts = await prisma.post.findMany({
 where: {
 authorId: {
 in: userIds,
 },
 },
});
```

```sql
SELECT "public"."User"."id", "public"."User"."email", "public"."User"."name" FROM "public"."User" WHERE 1=1 OFFSET $1
SELECT "public"."Post"."id", "public"."Post"."createdAt", "public"."Post"."updatedAt", "public"."Post"."title", "public"."Post"."content", "public"."Post"."published", "public"."Post"."authorId" FROM "public"."Post" WHERE "public"."Post"."authorId" IN ($1,$2,$3,$4) OFFSET $5
```

#### Solving n+1 with `relationLoadStrategy: "join"`

You can perform the query with a [database join](/orm/prisma-client/queries/relation-queries#relation-load-strategies-preview) by setting `relationLoadStrategy` to `"join"`, ensuring that only **one** query is executed against the database.

```ts
const users = await prisma.user.findMany({});

const userIds = users.map((x) => x.id);

const posts = await prisma.post.findMany({
 relationLoadStrategy: "join",
 where: {
 authorId: {
 in: userIds,
 },
 },
});
```

---
title: Connection management
description: This page explains how database connections are handled with Prisma Client and how to manually connect and disconnect your database
url: /orm/prisma-client/setup-and-configuration/databases-connections/connection-management
metaTitle: Connection management
metaDescription: This page explains how database connections are handled with Prisma Client and how to manually connect and disconnect your database.
---

:::info[Quick summary]
This page explains how Prisma Client manages database connections, including how and when to use the `$connect()` and `$disconnect()` methods, connection pooling behavior, and best practices for both long-running and serverless environments.
:::

`PrismaClient` connects and disconnects from your data source using the following two methods:

- [`$connect()`](/orm/reference/prisma-client-reference)
- [`$disconnect()`](/orm/reference/prisma-client-reference)

In most cases, you **do not need to explicitly call these methods**. `PrismaClient` automatically connects when you run your first query, creates a [connection pool](/orm/prisma-client/setup-and-configuration/databases-connections/connection-pool), and disconnects when the Node.js process ends.

See the [connection management guide](/orm/prisma-client/setup-and-configuration/databases-connections) for information about managing connections for different deployment paradigms (long-running processes and serverless functions).

<details>
<summary>Questions answered in this page</summary>

- When should I call $connect and $disconnect?
- How does Prisma manage connection pools?
- How to handle connections in serverless?

</details>

## `$connect()`

It is not necessary to call [`$connect()`](/orm/reference/prisma-client-reference) thanks to the _lazy connect_ behavior: The `PrismaClient` instance connects lazily when the first request is made to the API (`$connect()` is called for you under the hood).

### Calling `$connect()` explicitly

If you need the first request to respond instantly and cannot wait for a lazy connection to be established, you can explicitly call `prisma.$connect()` to establish a connection to the data source:

```ts
const prisma = new PrismaClient();

// run inside `async` function
await prisma.$connect();
```

## `$disconnect()`

When you call [`$disconnect()`](/orm/reference/prisma-client-reference) , Prisma Client:

1. Runs the [`beforeExit` hook](#exit-hooks)
2. Closes all connections in the pool

In a long-running application such as a GraphQL API, which constantly serves requests, it does not make sense to `$disconnect()` after each request - it takes time to establish a connection, and doing so as part of each request will slow down your application.

:::tip

To avoid too _many_ connections in a long-running application, we recommend that you [use a single instance of `PrismaClient` across your application](/orm/prisma-client/setup-and-configuration/introduction#use-prisma-client-to-send-queries-to-your-database).

:::

### Calling `$disconnect()` explicitly

In most long-running or serverless apps you should **not** call `$disconnect()` after each request, so connections can be reused. In some situations it **does** make sense to call it explicitly—for example, when creating a temporary `PrismaClient` and then immediately releasing its resources (e.g. in [Cloudflare Workers](/orm/prisma-client/deployment/edge/deploy-to-cloudflare), where `ctx.waitUntil(prisma.$disconnect())` is recommended).

Another scenario is a script that:

1. Runs **infrequently** (for example, a scheduled job to send emails each night), which means it does not benefit from a long-running connection to the database _and_
2. Exists in the context of a **long-running application**, such as a background service. If the application never shuts down, Prisma Client never disconnects.

The following script creates a new instance of `PrismaClient`, performs a task, and then disconnects - which closes the connection pool:

```ts
import { PrismaClient } from "../prisma/generated/client";

const prisma = new PrismaClient();
const emailService = new EmailService();

async function main() {
 const allUsers = await prisma.user.findMany();
 const emails = allUsers.map((x) => x.email);

 await emailService.send(emails, "Hello!");
}

main()
 .then(async () => {
 await prisma.$disconnect(); //[!code highlight]
 })
 .catch(async (e) => {
 console.error(e);
 await prisma.$disconnect(); //[!code highlight]
 process.exit(1);
 });
```

If the above script runs multiple times in the context of a long-running application _without_ calling `$disconnect()`, a new connection pool is created with each new instance of `PrismaClient`.

## Exit hooks

The `beforeExit` hook runs when Prisma ORM is triggered externally (e.g. via a `SIGINT` signal) to shut down, and allows you to run code _before_ Prisma Client disconnects - for example, to issue queries as part of a graceful shutdown of a service:

```ts
const prisma = new PrismaClient();

prisma.$on("beforeExit", async () => {
 console.log("beforeExit hook");
 // PrismaClient still available
 await prisma.message.create({
 data: {
 message: "Shutting down server",
 },
 });
});
```

---
title: Connection pool
description: Prisma Client uses a connection pool (from the database driver or driver adapter) to store and manage database connections.
url: /orm/prisma-client/setup-and-configuration/databases-connections/connection-pool
metaDescription: Prisma ORM's query engine creates a connection pool to store and manage database connections.
metaTitle: Connection pool
---

:::info[Quick summary]
This page explains how Prisma ORM manages database connections using a connection pool, and how you can configure limits and timeouts for optimal performance.
:::

Prisma Client uses a **connection pool** of database connections (managed by the database driver when using [driver adapters](/orm/core-concepts/supported-databases/database-drivers)). The pool is created when Prisma Client opens the _first_ connection to the database, which can happen in one of two ways:

- By [explicitly calling `$connect()`](/orm/prisma-client/setup-and-configuration/databases-connections/connection-management#connect) _or_
- By running the first query, which calls `$connect()` under the hood

Relational database connectors use Prisma ORM's own connection pool, and the MongoDB connectors uses the [MongoDB driver connection pool](https://github.com/mongodb/specifications/blob/master/source/connection-monitoring-and-pooling/connection-monitoring-and-pooling.rst).

<details>
<summary>Questions answered in this page</summary>

- How do I size Prisma's connection pool?
- How do I set pool timeouts and limits?
- When should I use PgBouncer with Prisma?

</details>

## Relational databases

Starting with Prisma ORM v7, relational datasources instantiate Prisma Client with [driver adapters](/orm/core-concepts/supported-databases/database-drivers) by default. Driver adapters rely on the Node.js driver you supply, so connection pooling defaults (and configuration) now come from the driver itself.

Use the tables below to translate Prisma ORM v6 connection URL parameters to the Prisma ORM v7 driver adapter fields alongside their defaults.

### Prisma ORM v7 driver adapter defaults

The following tables document the default connection pool settings for each driver adapter.

:::tip[Prisma timeouts]

Prisma ORM also has its own configurable timeouts that are separate from the database driver timeouts. If you see a timeout error and are unsure whether it comes from the driver or from Prisma Client, see the [Prisma Client timeouts and transaction options documentation](/orm/prisma-client/queries/transactions#transaction-isolation-level).

:::

#### PostgreSQL (using the `pg` driver adapter)

Here are the default connection pool settings for the `pg` driver adapter:

| Behavior | v6 URL parameter | v6 default | v7 `pg` config field | v7 default |
| ------------------- | ------------------------------ | ---------------------------------- | ------------------------- | ---------------- |
| Pool size | `connection_limit` | `num_cpus::get_physical() * 2 + 1` | `max` | `10` |
| Acquire timeout | `pool_timeout` | `10s` | `connectionTimeoutMillis` | `0` (no timeout) |
| Connection timeout | `connect_timeout` | `5s` | `connectionTimeoutMillis` | `0` (no timeout) |
| Idle timeout | `max_idle_connection_lifetime` | `300s` | `idleTimeoutMillis` | `10s` |
| Connection lifetime | `max_connection_lifetime` | `0` (no timeout) | `maxLifetimeSeconds` | `0` (no timeout) |

<details>
<summary>Example: Matching Prisma ORM v6 defaults with the `pg` driver adapter</summary>

If you want to preserve the same timeout behavior you had in Prisma ORM v6, pass the following configuration when instantiating the driver adapter:

```ts
import { PrismaPg } from "@prisma/adapter-pg";

const adapter = new PrismaPg({
 connectionString: process.env.DATABASE_URL,
 // Match Prisma ORM v6 defaults:
 connectionTimeoutMillis: 5_000, // v6 connect_timeout was 5s
 idleTimeoutMillis: 300_000, // v6 max_idle_connection_lifetime was 300s
});
```

</details>

:::tip
See the [node-postgres pool documentation](https://node-postgres.com/apis/pool) for details on every available option.
:::

#### MySQL or MariaDB (using the `mariadb` driver)

Here are the default connection pool settings for the `mariadb` driver adapter:

| Behavior | v6 URL parameter | v6 default | v7 `mariadb` config field | v7 default |
| ------------------ | ------------------------------ | ---------------------------------- | ------------------------- | ---------- |
| Pool size | `connection_limit` | `num_cpus::get_physical() * 2 + 1` | `connectionLimit` | `10` |
| Acquire timeout | `pool_timeout` | `10s` | `acquireTimeout` | `10s` |
| Connection timeout | `connect_timeout` | `5s` | `connectTimeout` | `1s` |
| Idle timeout | `max_idle_connection_lifetime` | `300s` | `idleTimeout` | `1800s` |

<details>
<summary>Example: Matching Prisma ORM v6 defaults with the `mariadb` driver adapter</summary>

If you want to preserve the same timeout behavior you had in Prisma ORM v6, pass the following configuration when instantiating the driver adapter:

```ts
import { PrismaMariaDb } from "@prisma/adapter-mariadb";

const adapter = new PrismaMariaDb({
 host: "localhost",
 port: 3306,
 user: process.env.DB_USER,
 password: process.env.DB_PASSWORD,
 database: process.env.DB_NAME,
 // Match Prisma ORM v6 defaults:
 connectTimeout: 5_000, // v6 connect_timeout was 5s
 idleTimeout: 300, // v6 max_idle_connection_lifetime was 300s (note: in seconds, not ms)
});
```

</details>

:::tip
Refer to the [MariaDB Connector/Node.js pool options](https://mariadb.com/docs/connectors/mariadb-connector-nodejs/connector-nodejs-promise-api#pool-options) for configuration and tuning guidance.
:::

#### SQL Server (using the `mssql` driver)

Here are the default connection pool settings for the `mssql` driver adapter:

| Behavior | v6 URL parameter | v6 default | v7 `mssql` config field | v7 default |
| ------------------ | ------------------------------ | ---------------------------------- | ------------------------ | ---------- |
| Pool size | `connection_limit` | `num_cpus::get_physical() * 2 + 1` | `pool.max` | `10` |
| Connection timeout | `connect_timeout` | `5s` | `connectionTimeout` | `15s` |
| Idle timeout | `max_idle_connection_lifetime` | `300s` | `pool.idleTimeoutMillis` | `30s` |

<details>
<summary>Example: Matching Prisma ORM v6 defaults with the `mssql` driver adapter</summary>

If you want to preserve the same timeout behavior you had in Prisma ORM v6, pass the following configuration when instantiating the driver adapter:

```ts
import { PrismaMssql } from "@prisma/adapter-mssql";

const adapter = new PrismaMssql({
 server: "localhost",
 port: 1433,
 database: "mydb",
 user: process.env.DB_USER,
 password: process.env.DB_PASSWORD,
 // Match Prisma ORM v6 defaults:
 connectionTimeout: 5_000, // v6 connect_timeout was 5s
 pool: {
 idleTimeoutMillis: 300_000, // v6 max_idle_connection_lifetime was 300s
 },
});
```

</details>

:::tip
See the [`node-mssql` pool docs](https://tediousjs.github.io/node-mssql/#general-same-for-all-drivers) for details on these fields.
:::

## MongoDB

The MongoDB connector does not use the Prisma ORM connection pool. The connection pool is managed internally by the MongoDB driver and [configured via connection string parameters](https://www.mongodb.com/docs/manual/reference/connection-string-options/#connection-pool-options).

## External connection poolers

The pool size cannot exceed what the underlying database can support. Configure pool size and timeouts via your [driver adapter](/orm/prisma-client/setup-and-configuration/databases-connections/connection-pool) (see the tables above). This is a particular challenge in serverless environments, where each function manages an instance of `PrismaClient` and its own connection pool.

Consider introducing [an external connection pooler like PgBouncer](/orm/prisma-client/setup-and-configuration/databases-connections#pgbouncer) to prevent your application or functions from exhausting the database connection limit.

## Manual database connection handling

When using Prisma Client with a driver adapter, database connections are managed by the driver and its pool. They are not exposed to the developer and it is not possible to manually access individual connections.

---
title: Database connections
description: Learn how to manage database connections and configure connection pools
url: /orm/prisma-client/setup-and-configuration/databases-connections
metaTitle: Database connections
metaDescription: Databases connections
---

Databases can handle a limited number of concurrent connections. Each connection requires RAM, which means that simply increasing the database connection limit without scaling available resources:

- ✔ might allow more processes to connect _but_
- ✘ significantly affects **database performance**, and can result in the database being **shut down** due to **exhaustion of system resources**

The way your application **manages connections** also impacts performance. This guide describes how to approach connection management in [serverless environments](#serverless-environments-faas) and [long-running processes](#long-running-processes).

:::warning

This guide focuses on **relational databases** and how to configure and tune the Prisma ORM connection pool (MongoDB uses the MongoDB driver connection pool).

:::

## Long-running processes

Examples of long-running processes include Node.js applications hosted on a service like Heroku or a virtual machine. Use the following checklist as a guide to connection management in long-running environments:

- Configure [pool size and timeouts](/orm/prisma-client/setup-and-configuration/databases-connections/connection-pool) for your driver adapter (defaults and options are adapter-specific)
- Make sure you have [**one** global instance of `PrismaClient`](#prismaclient-in-long-running-applications)

### `PrismaClient` in long-running applications

In **long-running** applications, we recommend that you:

- ✔ Create **one** instance of `PrismaClient` and re-use it across your application
- ✔ Assign `PrismaClient` to a global variable _in dev environments only_ to [prevent hot reloading from creating new instances](#prevent-hot-reloading-from-creating-new-instances-of-prismaclient)

#### Re-using a single `PrismaClient` instance

To re-use a single instance, create a module that exports a `PrismaClient` object:

```ts title="client.ts"
import { PrismaClient } from "../prisma/generated/client";

let prisma = new PrismaClient();

export default prisma;
```

The object is [cached](https://nodejs.org/api/modules.html#modules_caching) the first time the module is imported. Subsequent requests return the cached object rather than creating a new `PrismaClient`:

```ts title="app.ts"
import prisma from "./client";

async function main() {
 const allUsers = await prisma.user.findMany();
}

main();
```

You do not have to replicate the example above exactly - the goal is to make sure `PrismaClient` is cached. For example, you can [instantiate `PrismaClient` in the `context` object](https://github.com/prisma/prisma-examples/blob/9f1a6b9e7c25b9e1851bd59b273046158d748995/typescript/graphql-express/src/context.ts#L9) that you [pass into an Express app](https://github.com/prisma/prisma-examples/blob/9f1a6b9e7c25b9e1851bd59b273046158d748995/typescript/graphql-express/src/server.ts#L12).

#### Do not explicitly `$disconnect()`

You [do not need to explicitly `$disconnect()`](/orm/prisma-client/setup-and-configuration/databases-connections/connection-management#calling-disconnect-explicitly) in the context of a long-running application that is continuously serving requests. Opening a new connection takes time and can slow down your application if you disconnect after each query.

#### Prevent hot reloading from creating new instances of `PrismaClient`

Frameworks like [Next.js](https://nextjs.org/) support hot reloading of changed files, which enables you to see changes to your application without restarting. However, if the framework refreshes the module responsible for exporting `PrismaClient`, this can result in **additional, unwanted instances of `PrismaClient` in a development environment**.

As a workaround, you can store `PrismaClient` as a global variable in development environments only, as global variables are not reloaded:

```ts title="client.ts"
import { PrismaClient } from "../prisma/generated/client";

const globalForPrisma = globalThis as unknown as { prisma: PrismaClient };

export const prisma = globalForPrisma.prisma || new PrismaClient();

if (process.env.NODE_ENV !== "production") globalForPrisma.prisma = prisma;
```

The way that you import and use Prisma Client does not change:

```ts title="app.ts"
import { prisma } from "./client";

async function main() {
 const allUsers = await prisma.user.findMany();
}

main();
```

## Connections Created per CLI Command

In local tests with Postgres, MySQL, and SQLite, each Prisma CLI command typically uses a single connection. The table below shows the ranges observed in these tests. Your environment _may_ produce slightly different results.

| Command | Connections | Description |
| ------------------------------------------------------------------------------ | ----------- | ------------------------------------------------ |
| [`migrate status`](/orm/reference/prisma-cli-reference#migrate-status) | 1 | Checks the status of migrations |
| [`migrate dev`](/orm/reference/prisma-cli-reference#migrate-dev) | 1–4 | Applies pending migrations in development |
| [`migrate diff`](/orm/reference/prisma-cli-reference#migrate-diff) | 1–2 | Compares database schema with migration history |
| [`migrate reset`](/orm/reference/prisma-cli-reference#migrate-reset) | 1–2 | Resets the database and reapplies migrations |
| [`migrate deploy`](/orm/reference/prisma-cli-reference#migrate-deploy) | 1–2 | Applies pending migrations in production |
| [`db pull`](/orm/reference/prisma-cli-reference#db-pull) | 1 | Pulls the database schema into the Prisma schema |
| [`db push`](/orm/reference/prisma-cli-reference#db-push) | 1–2 | Pushes the Prisma schema to the database |
| [`db execute`](/orm/reference/prisma-cli-reference#db-execute) | 1 | Executes raw SQL commands |
| [`db seed`](/orm/reference/prisma-cli-reference#db-seed) | 1 | Seeds the database with initial data |

## Serverless environments (FaaS)

Examples of serverless environments include Node.js functions hosted on AWS Lambda, Vercel or Netlify Functions. Use the following checklist as a guide to connection management in serverless environments:

- Familiarize yourself with the [serverless connection management challenge](#the-serverless-challenge)
- Configure [pool size and timeouts](/orm/prisma-client/setup-and-configuration/databases-connections/connection-pool) for your driver adapter (defaults and options are adapter-specific)
- [Instantiate `PrismaClient` outside the handler](#instantiate-prismaclient-outside-the-handler) and do not explicitly `$disconnect()`
- Configure [function concurrency](#concurrency-limits) and handle [idle connections](#zombie-connections)

### The serverless challenge

In a serverless environment, each function creates **its own instance** of `PrismaClient`, and each client instance has its own connection pool.

Consider the following example, where a single AWS Lambda function uses `PrismaClient` to connect to a database. The `connection_limit` is **3**:

![An AWS Lambda function connecting to a database.](/img/orm/prisma-client/setup-and-configuration/databases-connections/serverless-connections.png)

A traffic spike causes AWS Lambda to spawn two additional lambdas to handle the increased load. Each lambda creates an instance of `PrismaClient`, each with a `connection_limit` of **3**, which results in a maximum of **9** connections to the database:

![Three AWS Lambda function connecting to a database.](/img/orm/prisma-client/setup-and-configuration/databases-connections/serverless-connections-2.png)

Many _concurrent functions_ responding to a traffic spike 📈 can exhaust the database connection limit very quickly. Furthermore, any functions that are **paused** keep their connections open by default and block them from being used by another function.

1. Configure a small pool size for your [driver adapter](/orm/prisma-client/setup-and-configuration/databases-connections/connection-pool) (adapter-specific; start small when not using a pooler)
2. If you need more connections per function, consider using an [external connection pooler like PgBouncer](#external-connection-poolers)

### `PrismaClient` in serverless environments

#### Instantiate `PrismaClient` outside the handler

Instantiate `PrismaClient` [outside the scope of the function handler](https://github.com/prisma/e2e-tests/blob/5d1041d3f19245d3d237d959eca94d1d796e3a52/platforms/serverless-lambda/index.ts#L3) to increase the chances of reuse. As long as the handler remains 'warm' (in use), the connection is potentially reusable:

```ts highlight=3;normal
import { PrismaClient } from "../prisma/generated/client";

const client = new PrismaClient();

export async function handler() {
 /* ... */
}
```

#### Do not explicitly `$disconnect()`

You [do not need to explicitly `$disconnect()`](/orm/prisma-client/setup-and-configuration/databases-connections/connection-management#calling-disconnect-explicitly) at the end of a function, as there is a possibility that the container might be reused. Opening a new connection takes time and slows down your function's ability to process requests. In some cases (e.g. [Cloudflare Workers](/orm/prisma-client/deployment/edge/deploy-to-cloudflare)), calling `$disconnect()` when releasing a temporary client is recommended—see the [connection management caveat](/orm/prisma-client/setup-and-configuration/databases-connections/connection-management#calling-disconnect-explicitly).

### Other serverless considerations

#### Container reuse

There is no guarantee that subsequent nearby invocations of a function will hit the same container - for example, AWS can choose to create a new container at any time.

Code should assume the container to be stateless and create a connection only if it does not exist - Prisma Client JS already implements this logic.

#### Zombie connections

Containers that are marked "to be removed" and are not being reused still **keep a connection open** and can stay in that state for some time (unknown and not documented from AWS). This can lead to sub-optimal utilization of the database connections.

A potential solution is to **clean up idle connections** ([`serverless-mysql`](https://github.com/jeremydaly/serverless-mysql) implements this idea, but cannot be used with Prisma ORM).

#### Concurrency limits

Depending on your serverless concurrency limit (the number of serverless functions running in parallel), you might still exhaust your database's connection limit. This can happen when too many functions are invoked concurrently, each with its own connection pool, which eventually exhausts the database connection limit. To prevent this, you can [set your serverless concurrency limit](https://docs.aws.amazon.com/lambda/latest/dg/configuration-concurrency.html) to a number lower than the maximum connection limit of your database divided by the number of connections used by each function invocation (as you might want to be able to connect from another client for other purposes).

## Optimizing the connection pool

If Prisma Client cannot obtain a connection from the pool before the adapter's acquire timeout, you will see connection pool timeout exceptions in your log. A connection pool timeout can occur if:

- Many users are accessing your app simultaneously
- You send a large number of queries in parallel (for example, using `await Promise.all()`)

Pool size, acquire timeout, and other pool behavior are **configured per driver adapter**—there are no connection URL parameters for these in Prisma ORM v7. See the [connection pool reference](/orm/prisma-client/setup-and-configuration/databases-connections/connection-pool) for each adapter's pool settings (e.g. `max`, `connectionTimeoutMillis` for `pg`) and the underlying driver documentation. Tune the pool so that:

- Your database can support the total number of concurrent connections (pool size × number of instances)
- Timeouts and queue behavior match your workload (e.g. avoid exhausting system resources if the queue grows unbounded)

## External connection poolers

Connection poolers like PgBouncer prevent your application from exhausting the database's connection limit.

To keep Prisma Client on the pooled connection while allowing Prisma CLI commands (for example, migrations or introspection) to connect directly, define two environment variables:

```bash title=".env"
# Connection URL to your database using PgBouncer.
DATABASE_URL="postgres://root:password@127.0.0.1:54321/postgres?pgbouncer=true"

# Direct connection URL to the database used for Prisma CLI commands. # [!code ++]
DIRECT_URL="postgres://root:password@127.0.0.1:5432/postgres" # [!code ++]
```

Configure `prisma.config.ts` to point to the direct connection string. Prisma CLI commands always read from this configuration.

```ts title="prisma.config.ts" showLineNumbers
import "dotenv/config";
import { defineConfig, env } from "prisma/config";

export default defineConfig({
 schema: "prisma/schema.prisma",
 datasource: {
 url: env("DIRECT_URL"),
 },
});
```

At runtime, instantiate Prisma Client with a driver adapter (for example, `@prisma/adapter-pg`) that uses the pooled connection string:

```ts title="src/db/client.ts" showLineNumbers
import { PrismaClient } from "../prisma/generated/client";
import { PrismaPg } from "@prisma/adapter-pg";

const adapter = new PrismaPg({ connectionString: process.env.DATABASE_URL });
export const prisma = new PrismaClient({ adapter });
```

### PgBouncer

PostgreSQL only supports a certain amount of concurrent connections, and this limit can be reached quite fast when the service usage goes up – especially in [serverless environments](#serverless-environments-faas).

[PgBouncer](https://www.pgbouncer.org/) holds a connection pool to the database and proxies incoming client connections by sitting between Prisma Client and the database. This reduces the number of processes a database has to handle at any given time. PgBouncer passes on a limited number of connections to the database and queues additional connections for delivery when connections become available. To use PgBouncer, see [Configure Prisma Client with PgBouncer](/orm/prisma-client/setup-and-configuration/databases-connections/pgbouncer).

### AWS RDS Proxy

Due to the way AWS RDS Proxy pins connections, [it does not provide any connection pooling benefits](/orm/prisma-client/deployment/caveats-when-deploying-to-aws-platforms#aws-rds-proxy) when used together with Prisma Client.

---
title: Configure Prisma Client with PgBouncer
description: 'Configure Prisma Client with PgBouncer and other poolers: when to use pgbouncer=true, required transaction mode, prepared statements, and Prisma Migrate workarounds'
url: /orm/prisma-client/setup-and-configuration/databases-connections/pgbouncer
metaTitle: Configure Prisma Client with PgBouncer
metaDescription: 'Configure Prisma Client with PgBouncer and other poolers: when to use pgbouncer=true, required transaction mode, prepared statements, and Prisma Migrate workarounds.'
---

An external connection pooler like PgBouncer holds a connection pool to the database, and proxies incoming client connections by sitting between Prisma Client and the database. This reduces the number of processes a database has to handle at any given time.

Usually, this works transparently, but some connection poolers only support a limited set of functionality. One common feature that external connection poolers do not support are named prepared statements, which Prisma ORM uses. For these cases, Prisma ORM can be configured to behave differently.

<details>
<summary>Questions answered in this page</summary>

- How do I configure Prisma with PgBouncer?
- Do I need `pgbouncer=true`, and if so, when?
- How does Prisma Migrate work with PgBouncer?

</details>

:::info

Looking for an easy, infrastructure-free solution? Try [Prisma Accelerate](https://www.prisma.io/accelerate?utm_source=docs&utm_campaign=pgbouncer-help)! It requires little to no setup and works seamlessly with all databases supported by Prisma ORM.

Ready to begin? Get started with Prisma Accelerate by clicking [here](https://console.prisma.io?utm_source=docs&utm_campaign=pgbouncer-help).

:::

## PgBouncer

### Set PgBouncer to transaction mode

For Prisma Client to work reliably, PgBouncer must run in [**Transaction mode**](https://www.pgbouncer.org/features.html).

Transaction mode offers a connection for every transaction – a requirement for the Prisma Client to work with PgBouncer.

### Add `pgbouncer=true` for PgBouncer versions below `1.21.0`

:::warning
We recommend **not** setting `pgbouncer=true` in the database connection string if you're using [PgBouncer `1.21.0`](https://github.com/prisma/prisma/issues/21531#issuecomment-1919059472) or later.
:::

To use Prisma Client with PgBouncer, add the `?pgbouncer=true` flag to the PostgreSQL connection URL:

```shell
postgresql://USER:PASSWORD@HOST:PORT/DATABASE?pgbouncer=true
```

:::info
`PORT` specified for PgBouncer pooling is sometimes different from the default `5432` port. Check your database provider docs for the correct port number.
:::

### Configure `max_prepared_statements` in PgBouncer to be greater than zero

Prisma uses prepared statements, and setting [`max_prepared_statements`](https://www.pgbouncer.org/config.html) to a value greater than `0` enables PgBouncer to use those prepared statements.

:::info
`PORT` specified for PgBouncer pooling is sometimes different from the default `5432` port. Check your database provider docs for the correct port number.
:::

### Prisma Migrate and PgBouncer workaround

Prisma Migrate uses **database transactions** to check out the current state of the database and the migrations table. However, the Schema Engine is designed to use a **single connection to the database**, and does not support connection pooling with PgBouncer. If you attempt to run Prisma Migrate commands in any environment that uses PgBouncer for connection pooling, you might see the following error:

```bash
Error: undefined: Database error
Error querying the database: db error: ERROR: prepared statement "s0" already exists
```

To work around this issue, configure a **direct** connection for Prisma CLI commands in `prisma.config.ts`, while Prisma Client continues to use the PgBouncer URL via a driver adapter.

```bash title=".env"
# PgBouncer (pooled) connection string used by Prisma Client.
DATABASE_URL="postgres://USER:PASSWORD@HOST:PORT/DATABASE?pgbouncer=true"

# Direct database connection string used by Prisma CLI. # [!code ++]
DIRECT_URL="postgres://USER:PASSWORD@HOST:PORT/DATABASE" # [!code ++]
```

```ts title="prisma.config.ts" showLineNumbers
import "dotenv/config";
import { defineConfig, env } from "prisma/config";

export default defineConfig({
 schema: "prisma/schema.prisma",
 datasource: {
 url: env("DIRECT_URL"),
 },
});
```

```ts title="src/db/client.ts" showLineNumbers
import { PrismaClient } from "../prisma/generated/client";
import { PrismaPg } from "@prisma/adapter-pg";

const adapter = new PrismaPg({ connectionString: process.env.DATABASE_URL });
export const prisma = new PrismaClient({ adapter });
```

With this setup, PgBouncer stays in the path for runtime traffic, while Prisma CLI commands (`prisma migrate dev`, `prisma db push`, `prisma db pull`, and so on) always use the direct connection string defined in `prisma.config.ts`.

### PgBouncer with different database providers

There are sometimes minor differences in how to connect directly to a Postgres database that depend on the provider hosting the database.

Below are links to information on how to set up these connections with providers who have setup steps not covered here in our documentation:

- [Connecting directly to a PostgreSQL database hosted on Digital Ocean](https://github.com/prisma/prisma/issues/6157)
- [Connecting directly to a PostgreSQL database hosted on ScaleGrid](https://github.com/prisma/prisma/issues/6701#issuecomment-824387959)

## Supabase Supavisor

Supabase's Supavisor behaves similarly to [PgBouncer](#pgbouncer). You can add `?pgbouncer=true` to your connection pooled connection string available via your [Supabase database settings](https://supabase.com/dashboard/project/_/settings/database).

## Other external connection poolers

Although Prisma ORM does not have explicit support for other connection poolers, if the limitations are similar to the ones of [PgBouncer](#pgbouncer) you can usually also use `pgbouncer=true` in your connection string to put Prisma ORM in a mode that works with them as well.
