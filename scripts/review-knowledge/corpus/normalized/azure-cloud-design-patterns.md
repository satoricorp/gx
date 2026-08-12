Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Architects design workloads by integrating platform services, functionality, and code to meet both functional and nonfunctional requirements. To design effective workloads, you must understand these requirements and select topologies and methodologies that address the challenges of your workload's constraints. Cloud design patterns provide solutions to many common challenges.

System design heavily relies on established design patterns. You can design infrastructure, code, and distributed systems by using a combination of these patterns. These patterns are crucial for building reliable, highly secure, cost-optimized, operationally efficient, and high-performing applications in the cloud.

The following cloud design patterns are technology-agnostic, which makes them suitable for any distributed system. You can apply these patterns across Azure, other cloud platforms, on-premises setups, and hybrid environments.

## How cloud design patterns enhance the design process

Cloud workloads are vulnerable to the fallacies of distributed computing, which are common but incorrect assumptions about how distributed systems operate. Examples of these fallacies include:

- The network is reliable.
- Latency is zero.
- Bandwidth is infinite.
- The network is secure.
- Topology doesn't change.
- There's one administrator.
- Component versioning is simple.
- Observability implementation can be delayed.

These misconceptions can result in flawed workload designs. Design patterns don't eliminate these misconceptions, but they help raise awareness, provide compensation strategies, and offer mitigations. Each cloud design pattern has trade-offs. Focus on why you should choose a specific pattern instead of how to implement it.

Consider how to use these industry-standard design patterns as the core building blocks for a well-architected workload design. Each design pattern in the Azure Well-Architected Framework represents one or more of its pillars. Some patterns might introduce trade-offs that affect the goals of other pillars.

## Select a pattern

Choose a pattern based on the problem you need to solve, not the technology you want to use. Begin with a specific constraint or risk in your workload, such as a service that fails under load, a data store that can't keep up with read queries, or a dependency that you can't fully trust. A pattern is a good fit when its problem statement matches the challenge you face and when the trade-offs it introduces are ones you can accept in your workload.

## Pattern catalog

Each pattern in this catalog describes the problem that it addresses, considerations for applying the pattern, and an example based on Microsoft Azure services and tools. Some patterns include code samples or snippets that show how to implement the pattern on Azure.

| Pattern | Summary | Well-Architected Framework pillars |
|---|---|---|
| Ambassador | Create helper services that send network requests on behalf of a consumer service or application. | - Reliability - Security |
| Anti-Corruption Layer | Implement a façade or adapter layer between a modern application and a legacy system. | - Operational Excellence |
| Asynchronous Request-Reply | Decouple back-end processing from a front-end host. This pattern is useful when back-end processing must be asynchronous, but the front end requires a clear and timely response. | - Performance Efficiency |
| Backends for Frontends | Create separate backend services for specific frontend applications or interfaces. | - Reliability - Security - Performance Efficiency |
| Bulkhead | Isolate elements of an application into pools so that if one fails, the others continue to function. | - Reliability - Security - Performance Efficiency |
| Cache-Aside | Load data on demand into a cache from a data store. | - Reliability - Performance Efficiency |
| Choreography | Let individual services decide when and how a business operation is processed, instead of depending on a central orchestrator. | - Operational Excellence - Performance Efficiency |
| Circuit Breaker | Handle faults that might take a variable amount of time to fix when an application connects to a remote service or resource. | - Reliability - Performance Efficiency |
| Claim Check | Split a large message into a claim check and a payload to avoid overwhelming a message bus. | - Reliability - Security - Cost Optimization - Performance Efficiency |
| Compensating Transaction | Undo the work performed by a sequence of steps that collectively form an eventually consistent operation. | - Reliability |
| Competing Consumers | Enable multiple concurrent consumers to process messages that they receive on the same messaging channel. | - Reliability - Cost Optimization - Performance Efficiency |
| Compute Resource Consolidation | Consolidate multiple tasks or operations into a single computational unit. | - Cost Optimization - Operational Excellence - Performance Efficiency |
| CQRS | Separate operations that read data from those that update data by using distinct interfaces. | - Performance Efficiency |
| Deployment Stamps | Deploy multiple independent copies of application components, including data stores. | - Operational Excellence - Performance Efficiency |
| Event Sourcing | Use an append-only store to record a full series of events that describe actions taken on data in a domain. | - Reliability - Performance Efficiency |
| External Configuration Store | Move configuration information out of an application deployment package to a centralized location. | - Operational Excellence |
| Federated Identity | Delegate authentication to an external identity provider. | - Reliability - Security - Performance Efficiency |
| Gatekeeper | Protect applications and services by using a dedicated host instance to validate and sanitize requests before forwarding them to private back ends. | - Security - Performance Efficiency |
| Gateway Aggregation | Use a gateway to aggregate multiple individual requests into a single request. | - Reliability - Security - Operational Excellence - Performance Efficiency |
| Gateway Offloading | Offload shared or specialized service functionality to a gateway proxy. | - Reliability - Security - Cost Optimization - Operational Excellence - Performance Efficiency |
| Gateway Routing | Route requests to multiple services by using a single endpoint. | - Reliability - Operational Excellence - Performance Efficiency |
| Geode | Deploy back-end services across geographically distributed nodes. Each node can handle client requests from any region. | - Reliability - Performance Efficiency |
| Health Endpoint Monitoring | Implement functional checks in an application that external tools can access through exposed endpoints at regular intervals. | - Reliability - Operational Excellence - Performance Efficiency |
| Index Table | Create indexes over the fields in data stores that queries frequently reference. | - Reliability - Performance Efficiency |
| Leader Election | Coordinate actions in a distributed application by electing one instance as the leader. The leader manages a collection of collaborating task instances. | - Reliability |
| Materialized View | Generate prepopulated views over the data in one or more data stores when the data is poorly formatted for required query operations. | - Performance Efficiency |
| Messaging Bridge | Build an intermediary to enable communication between messaging systems that are otherwise incompatible. | - Cost Optimization - Operational Excellence |
| Pipes and Filters | Break down a task that performs complex processing into a series of separate elements that can be reused. | - Reliability |
| Priority Queue | Prioritize requests sent to services so that requests with a higher priority are processed more quickly. | - Reliability - Performance Efficiency |
| Publisher-Subscriber | Enable an application to announce events to multiple consumers asynchronously, without coupling senders to receivers. | - Reliability - Security - Cost Optimization - Operational Excellence - Performance Efficiency |
| Quarantine | Ensure that external assets meet a team-agreed quality level before the workload consumes them. | - Security - Operational Excellence |
| Queue-Based Load Leveling | Use a queue that creates a buffer between a task and a service to smooth intermittent heavy loads. | - Reliability - Cost Optimization - Performance Efficiency |
| Rate Limiting | Avoid or minimize throttling errors by controlling the consumption of resources. | - Reliability |
| Retry | Enable applications to handle anticipated temporary failures by retrying failed operations. | - Reliability |
| Saga | Manage data consistency across microservices in distributed transaction scenarios. | - Reliability |
| Scheduler Agent Supervisor | Coordinate a set of actions across distributed services and resources. | - Reliability - Performance Efficiency |
| Sequential Convoy | Process a set of related messages in a defined order without blocking other message groups. | - Reliability |
| Sharding | Divide a data store into a set of horizontal partitions or shards. | - Reliability - Cost Optimization |
| Sidecar | Deploy components into a separate process or container to provide isolation and encapsulation. | - Security - Operational Excellence |
| Static Content Hosting | Deploy static content to a cloud-based storage service for direct client delivery. | - Cost Optimization |
| Strangler Fig | Incrementally migrate a legacy system by gradually replacing pieces of functionality with new applications and services. | - Reliability - Cost Optimization - Operational Excellence |
| Throttling | Control the consumption of resources from applications, tenants, or services. | - Reliability - Security - Cost Optimization - Performance Efficiency |
| Valet Key | Use a token or key to provide clients with restricted, direct access to a specific resource or service. | - Security - Cost Optimization - Performance Efficiency |

### Combine patterns

Patterns are composable. A single pattern addresses one problem, but a workload usually faces several problems at once, so you often apply multiple patterns together. Some patterns also build on others or pair naturally to cover a gap that one pattern leaves open. Consider these examples:

- Pair Retry with Circuit Breaker so that an application retries transient faults but stops retrying when a fault persists.
- Combine Queue-Based Load Leveling with Competing Consumers to buffer load and then scale the processing of that load.
- Layer the Gateway Routing, Gateway Aggregation, and Gateway Offloading patterns behind a single gateway endpoint.
- Build Saga on Compensating Transaction to maintain data consistency across services when a distributed operation fails partway through.

### Antipatterns

A design pattern describes a practice to apply. An antipattern describes a practice to avoid. Antipatterns often start as reasonable designs that work in testing or at low scale, but they degrade reliability or performance as load increases. Recognizing an antipattern helps you spot a problem in an existing design and choose a pattern that resolves it. For the full list, see Antipatterns for cloud applications.

## AI agent orchestration patterns

The preceding cloud design patterns address common challenges in distributed systems, but AI workloads that use multiple autonomous agents require specialized coordination approaches. Traditional patterns like Scheduler Agent Supervisor or Choreography provide foundational concepts. However, AI agents introduce unique challenges such as nondeterministic outputs, dynamic reasoning capabilities, and the need for intelligent handoffs between specialized components.

For AI workloads that include multiple autonomous agents, see AI agent orchestration patterns. These patterns complement the cloud design patterns in this catalog by addressing the specific coordination requirements of intelligent, autonomous components that work together to accomplish complex outcomes.

## Next steps

Review the design patterns from the perspective of the Well-Architected Framework pillar that the pattern aims to optimize.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Create helper services that send network requests on behalf of a consumer service or application. Think of an ambassador service as an out-of-process proxy that's colocated with the client.

Use the Ambassador pattern to offload common client connectivity tasks like monitoring, logging, routing, security like Transport Layer Security (TLS), and resiliency patterns in a language-agnostic way. Extend the networking capabilities of legacy applications, or other applications that are difficult to modify, by using the Ambassador pattern. Specialized teams can also use the Ambassador pattern to implement those features.

## Context and problem

Resilient cloud-based applications require features like circuit breaking, routing, metering and monitoring, and network-related configuration updates. If the development team doesn't maintain the code or can't easily modify it, it might be difficult or even impossible to update legacy applications or existing code libraries to add these features.

Network calls might also require substantial configuration for connection, authentication, and authorization. When multiple applications use these calls across different languages and frameworks, you must configure the calls separately for each instance. A central team within your organization might need to manage network and security functionality. With a large code base, it can be risky for that team to update unfamiliar application code.

## Solution

Put client frameworks and libraries into an external process that acts as a proxy between your application and external services. To provide control over routing, resiliency, and security features and to avoid host-related access restrictions, deploy the proxy on the same host environment as your application. Use the ambassador pattern to standardize and extend instrumentation. The proxy can monitor performance metrics, like latency or resource usage, in the same host environment as the application.

Diagram that shows a client application and an ambassador proxy colocated on the same host. The client application sends requests to the ambassador instead of calling external services directly. The ambassador forwards those requests to the remote service. Responses from the remote service return through the ambassador and back to the client application.

You can manage features that are offloaded to the ambassador independently of the application. You can update and modify the ambassador without disturbing the application's legacy functionality. Separate, specialized teams can also implement and maintain security, networking, or authentication features that have been moved to the ambassador.

You can deploy ambassador services as a sidecar to accompany the life cycle of a consuming application or service. Alternatively, if multiple separate processes on a common host share an ambassador, you can deploy it as a daemon or Windows service. If the consuming service is containerized, create the ambassador as a separate container on the same host and set up the appropriate links for communication.

## Problems and considerations

Consider the following points when you decide how to implement this pattern:

- The proxy adds some latency overhead. Consider whether a client library that the application directly invokes is a better approach.
- Consider the possible impact of including generalized features in the proxy. For example, the ambassador could handle retries, but that approach might not be safe unless all operations are idempotent.
- Consider a mechanism that the client can use to pass some context to the proxy and back to the client. For example, include HTTP request headers to opt out of retry or specify the maximum number of times to retry.
- Consider how to package and deploy the proxy.
- Consider whether to use a single shared instance for all clients or an instance for each client.

## When to use this pattern

Use this pattern when:

- You must build a common set of client connectivity features for multiple languages or frameworks.
- You must offload cross-cutting client connectivity concerns to infrastructure developers or other more specialized teams.
- You must support cloud or cluster connectivity requirements in a legacy application or an application that's difficult to modify.
- You must support protocols or connectivity patterns that API gateways, service meshes, or standard ingress and egress controls don't handle easily.

This pattern might not be suitable when:

- Network request latency is critical. A proxy introduces minimal overhead, and this overhead might affect the application.
- Client connectivity features are consumed by a single language. In that case, a better option might be a client library that's distributed to the development teams as a package.
- Connectivity features can't be generalized and these features require deeper integration with the client application.
- Your application platform supports prebuilt solutions, like a service mesh, to handle mutual TLS (mTLS), traffic management, and policy capabilities. Use these solutions instead of creating a custom ambassador solution.

## Workload design

Evaluate how to use the Ambassador pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | This pattern introduces a network communications mediation point, so you can add reliability patterns to network communication, like retry or buffering. - RE:07 Self-preservation |
| Security design decisions help ensure the confidentiality,integrity, andavailabilityof your workload's data and systems. | With this pattern, you can implement security on network communications that the client can't handle directly. - SE:06 Network controls - SE:07 Encryption |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

The following diagram shows an application making a request to a remote service via an ambassador proxy. The ambassador provides routing, circuit breaking, and logging. It calls the remote service and then returns the response to the client application.

Diagram that shows a client application sending a request to an ambassador proxy. The application sends a request to the remote service via an ambassador proxy. The ambassador determines the location of the remote services and routes the request appropriately. The ambassador checks the circuit breaker state and enriches request headers with tracing information. The ambassador starts measuring the request latency. The ambassador encrypts and sends the request using mutual certificate-based authentication. The remote service receives the request and sends the response. The ambassador logs the request latency. The ambassador returns the response to the client. The application receives the response.

In containerized environments, this ambassador would run as a sidecar container next to the application container. In noncontainerized environments, you would implement it as a local process or Windows service on the same host.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Implement a facade or adapter layer between different subsystems that don't share the same semantics. This layer translates requests that one subsystem makes to the other subsystem. Use this pattern to ensure that dependencies on outside subsystems don't limit an application's design. Eric Evans first described this pattern in *Domain-Driven Design: Tackling Complexity in the Heart of Software*.

## Context and problem

Most applications rely on other systems for some data or functionality. For example, when you migrate a legacy application to a modern system, the application might continue to use existing legacy resources. New features must be able to call the legacy system. This capability is especially important for gradual migrations in which you move different features of a larger application to a modern system over time.

These legacy systems often have quality problems like convoluted data schemas or obsolete APIs. The features and technologies that legacy systems use can vary widely from more modern systems. To interoperate with the legacy system, the new application might need to support outdated infrastructure, protocols, data models, APIs, or other features that you wouldn't otherwise put into a modern application.

When you maintain access between new and legacy systems, you force the new system to adhere to at least some of the legacy system's APIs or other semantics. When these legacy features have quality problems, this support corrupts what might otherwise be a cleanly designed modern application.

Similar problems can arise with any external system that your development team doesn't control.

## Solution

Isolate the different subsystems by placing an anti-corruption layer between them. This layer translates communication between the two systems. By using this approach, you can keep one system unchanged without compromising the design and technological approach of the other.

*Download a Visio file of this architecture.*

The diagram shows an application that has two subsystems. Subsystem A calls subsystem B through an anti-corruption layer. Communication between subsystem A and the anti-corruption layer always uses the data model and architecture of subsystem A. Calls from the anti-corruption layer to subsystem B conform to that subsystem's data model or methods. The anti-corruption layer contains all the logic necessary to translate between the two systems. You can implement the layer as a component within the application or as an independent service.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- The anti-corruption layer adds latency to calls between the two systems.
- The anti-corruption layer adds an extra service that you must manage and maintain.
- Consider how you plan to scale the anti-corruption layer.
- Consider whether you need more than one anti-corruption layer. For example, you might want to decompose functionality into multiple services that use different technologies or languages.
- Consider how you plan to manage the anti-corruption layer in relation to your other applications or services, and how to integrate it into your monitoring, release, and configuration processes.
- Make sure that you maintain and monitor transaction and data consistency.
- Consider whether the anti-corruption layer needs to handle all communication between different subsystems, or just a subset of features.
- If the anti-corruption layer is part of an application migration strategy, consider whether it's permanent or whether you plan to retire it after you migrate all legacy functionality.
- The previous diagram uses distinct subsystems to illustrate this pattern, but you can also apply it to other service architectures, such as legacy code integration in a monolithic architecture.
- Because the anti-corruption layer mediates systems that might have different trust levels, consider enforcing input validation and sanitization at this boundary.
- Plan for observability, including correlation IDs and structured logging, to diagnose translation failures.

## When to use this pattern

Use this pattern when:

- You plan a migration to happen over multiple stages, but you need to maintain integration between new and legacy systems.
- Two or more subsystems have different semantics, but they need to communicate.

This pattern might not be suitable when:

- The new and legacy systems have no significant semantic differences. In this scenario, it's important to focus the anti-corruption layer on translation logic. Avoid placing business rules or orchestration in the layer.

## Workload design

Evaluate how to use the Anti-Corruption Layer pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | This pattern helps ensure that new component design remains uninfluenced by legacy implementations that might have different data models or business rules when you integrate with these legacy systems. It can reduce technical debt in new components while still supporting existing components. - OE:04 Tools and processes - OE:07 Monitoring system |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

This pattern is conceptual and originates from the domain-driven design software development approach. Azure services like Azure API Management or Azure Functions might assist with protocol handling and translation, but the core purpose of an anti-corruption layer is to protect the domain model, not to prescribe any specific product choice.

In the following example, API Management handles the external exposure and protocol concerns. Azure Functions implements the anti-corruption layer through domain mapping between the new system and the legacy system. Azure Monitor and Application Insights provide the observability that you need to track the success and latency of the translation between the two subsystems.

Beyond this synchronous request-response model, the anti-corruption layer can also use an asynchronous, event-driven approach. By using Azure Service Bus, Azure Event Grid, or Azure Event Hubs, the layer decouples the modern domain from the legacy system's throughput constraints to allow message-based translation for high-throughput or highly decoupled workloads.

## Next steps

- Explore cloud design patterns that help manage distributed transactions and maintain data consistency, such as the Compensating Transaction pattern and Saga distributed transactions pattern.
- Because the anti-corruption layer can become a single point of failure, plan for resilience by using the Retry pattern, Circuit Breaker pattern, Bulkhead pattern, and Health Endpoint Monitoring pattern.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Decouple back-end processing from a front-end host when back-end processing needs to run asynchronously but the front end needs a clear response.

## Context and problem

In modern application development, client applications often depend on remote APIs to provide business logic and compose functionality. Many applications run code in a web browser, and other environments also host client code. The APIs might relate directly to the application or operate as shared services from an external service. Most API calls use HTTP or HTTPS and follow REST semantics.

In most cases, APIs for a client application respond in about 100 milliseconds (ms) or less. Many factors can affect the response latency:

- The application's hosting stack
- Security components
- The relative geographic location of the caller and the back end
- Network infrastructure
- Current load
- The size of the request payload
- Processing queue length
- The time for the back end to process the request

These factors can add latency to the response. You can mitigate some factors by scaling out the back end. Other factors, like network infrastructure, are outside the application developer's control. Most APIs respond quickly enough for the response to return over the same connection. Application code can make a synchronous API call in a nonblocking way to give the appearance of asynchronous processing. We recommend this approach for input and output (I/O)‑bound operations.

In some scenarios, the back end does work that's long-running and takes a few seconds. In other scenarios, the back end does long-running background work for minutes or for extended periods. In these cases, you can't wait for the work to finish before you send a response. This situation can create a problem for synchronous request-reply patterns. For guidance about designing the back-end processing, see Background jobs.

Some architectures solve this problem by using a message broker to separate the request and response stages. Many systems achieve this separation through the Queue-Based Load Leveling pattern. This separation lets the client process and the back-end API scale independently. It also introduces extra complexity when the client requires success notification because that step must also become asynchronous.

Many of the same considerations that apply to client applications also apply to server-to-server REST API calls in distributed systems, like in a microservices architecture.

## Solution

One solution to this problem is to use HTTP polling. Polling works well for client-side code when callback endpoints are unavailable or when long-running connections add too much complexity. Even when callbacks are possible, the extra libraries and services that they require can increase complexity.

The following steps describe the solution:

- The client application makes a synchronous call to the API to trigger a long-running operation on the back end.
- The API responds synchronously as quickly as possible. It returns an HTTP 202 (Accepted) status code to acknowledge that it received the request for processing. - Note - The API should validate the request and the action to be performed before it starts the long-running process. If the request isn't valid, reply immediately with an error code like HTTP 400 (Bad Request).
- The response includes a location reference that points to an endpoint that the client can poll to check the result of the long-running operation.
- The API offloads processing to another component, like a message queue.
- For every successful call to the status endpoint, the endpoint returns HTTP 200 (OK). While the work is in progress, the status endpoint returns a resource that indicates that state. The status response body should include enough information for the client to understand the current state of the operation. - When the work completes, the status endpoint returns a resource that indicates completion or redirects to another resource URL. For example, if the asynchronous operation creates a new resource, the status endpoint redirects to the URL for that resource.

The following diagram shows a typical flow.

- The client sends a request and receives an HTTP 202 (Accepted) response.
- The client sends an HTTP GET request to the status endpoint. This call returns HTTP 200 because the work is pending.
- At some point, the work completes and the status endpoint returns HTTP 303 (See Other) to redirect to the resource.
- The client fetches the resource at the specified URL.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- Multiple ways exist to implement this pattern over HTTP, and upstream services don't always use the same semantics. For example, some implementations don't use a separate status endpoint. Instead, the client polls the target resource URL directly and receives HTTP 404 (Not Found) until the resource is created. This response is generated because the resource doesn't exist yet. However, this approach can be unclear because invalid request IDs also return HTTP 404. A dedicated status endpoint that returns HTTP 200 with a status body, as described in this pattern, avoids this confusion.
- An HTTP 202 response indicates where the client polls and how often. It should include the following headers. - Header - Description - Notes - `Location`- A URL that the client polls for a response status - This URL can be a shared access signature (SAS) token. The Valet Key pattern works well when this location needs access control. The pattern also applies when response polling needs to move to another back end. - `Retry-After`- An estimated completion time for processing - This header helps polling clients avoid sending too many requests to the back end. - Consider expected client behavior when you design this response. A client that you control can follow these response values exactly. Clients that others author, including clients built by using no-code or low-code tools like Azure Logic Apps, can apply their own handling for HTTP 202.
- Consider including the following fields in the status endpoint response. - Field - Description - Notes - `status`- The current state of the operation, such as - *Pending*,- *Running*,- *Succeeded*,- *Failed*, or- *Canceled*- Uses a consistent, documented set of terminal and nonterminal values - `createdAt`- The time that the operation was accepted - Helps clients detect stale or abandoned operations - `lastUpdatedAt`- The time that the status was last updated - Helps clients distinguish between stalled and in-progress operations - `percentComplete`- An optional progress indicator - Useful when the back end can estimate progress - `error`- A structured error object when the status is - *Failed*- For consistency, consider using the RFC 9457 format.
- You might need to use a processing proxy to adjust the response headers or payload, depending on the underlying services that you use.
- If the status endpoint redirects after completion, use HTTP 303 (See Other). A 303 instructs the client to issue a GET request to the redirect URL, regardless of the original request method. This behavior is the correct semantic for this pattern because the client is retrieving a distinct result resource, not resubmitting the original operation. HTTP 302 (Found) doesn't guarantee a method change. Some clients replay the original method on redirect. This behavior can cause unintended side effects, such as duplicate POST requests.
- After the server successfully processes the request, the resource that the - `Location`header specifies returns an HTTP status code like 200, 201 (Created), or 204 (No Content).
- If an error occurs during processing, persist the error at the resource URL that the - `Location`header specifies and return a 4xx status code from the resource that matches the failure. Use a structured error format, such as RFC 9457 (Problem Details for HTTP APIs), so that clients can programmatically parse and handle failures.
- The status resource and any stored results consume storage and compute. Define a retention policy to clean them up after a reasonable period. To inform clients of the retention window, you can add an - `Expires`header to the status response.
- Solutions don't all implement this pattern the same way, and some services include extra or alternate headers. For example, Azure Resource Manager uses a modified variant of this pattern. For more information, see Resource Manager asynchronous operations.
- Legacy clients might not support this pattern. In that case, you might need to place a façade over the asynchronous API to hide the asynchronous processing from the original client. For example, Logic Apps supports this pattern natively, and you can use it as an integration layer between an asynchronous API and a client that makes synchronous calls. For more information, see Asynchronous request-response behavior in Logic Apps.
- To provide a way for clients to cancel a long-running request, expose a DELETE operation on the status endpoint resource. This request should forward a cancellation instruction to the back-end processing component. After the back end handles the cancellation, it should update the status resource to reflect the canceled state. This process helps prevent incomplete work from consuming resources indefinitely. Determine whether the operation supports partial rollback or requires a compensating transaction.
- You can require clients to supply an idempotency key, for example in an - `Idempotency-Key`request header, when they submit the initial request. If the back end receives a duplicate key, it should return the existing status resource instead of enqueuing a second work item. This approach protects against network failures that cause the client to retry a POST that the server already accepted. It's especially important in this pattern because the client can't distinguish between a lost response and a request that was never received.

Note

This pattern describes HTTP polling, in which the client periodically issues new requests to check status. In *long polling*, the client sends a request and the server holds the connection open until new data is available or a timeout occurs. Long polling reduces response latency compared to periodic polling, but it introduces complexity around connection management and timeouts.

## When to use this pattern

Use this pattern when:

- You work with client-side code, like browser applications, and those constraints make callback endpoints difficult to provide, or long-running connections add too much complexity.
- You call a service that uses only the HTTP protocol and the return service can't send callbacks because of firewall restrictions on the client side.
- You integrate with workloads that don't support modern callback mechanisms like WebSockets or webhooks.

This pattern might not be suitable when:

- You can use a service built for asynchronous notifications instead, like Azure Event Grid.
- Responses must stream in real time to the client. Consider using Server-Sent Events (SSEs), which provide a lightweight, HTTP-native, unidirectional push channel from server to client without requiring the client to poll.
- The client needs to collect many results, and the latency of those results is important. Consider using a message broker instead.
- Server-side persistent network connections like WebSockets or SignalR are available. You can use these connections to notify the caller of the result.
- The network design supports open ports to receive asynchronous callbacks or webhooks.

## Workload design

An architect should evaluate how they can use the Asynchronous Request-Reply pattern in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | You improve responsiveness and scalability by decoupling the request and reply phases for processes that don't require an immediate response. An asynchronous approach increases concurrency and lets the server schedule work as capacity becomes available. - PE:05 Scaling and partitioning - PE:07 Code and infrastructure |

As with any design decision, consider trade-offs against the goals of the other pillars that this pattern might introduce.

## Example

The following code shows excerpts from an application that uses Azure Functions to implement this pattern. This solution has three functions:

- The asynchronous API endpoint
- The status endpoint
- A back-end function that takes queued work items and runs them

This sample is available on GitHub.

The implementation uses managed identity to authenticate with Azure Service Bus and Azure Blob Storage, which avoids storing connection strings or account keys. Dependencies are registered in `Program.cs` by using `DefaultAzureCredential` and are injected through primary constructors.

### AsyncProcessingWorkAcceptor function

The `AsyncProcessingWorkAcceptor` function implements an endpoint that accepts work from a client application and enqueues it for processing:

- The function generates a request ID and adds it as metadata to the queue message.
- The HTTP response includes a - `Location`header that points to a status endpoint and a- `Retry-After`header that suggests a polling interval. The request ID appears in the URL path.

```
public class AsyncProcessingWorkAcceptor(ServiceBusClient _serviceBusClient)
{
 [Function("AsyncProcessingWorkAcceptor")]
 public async Task<IActionResult> Run(
 [HttpTrigger(AuthorizationLevel.Anonymous, "post", Route = null)] HttpRequest req,
 [FromBody] CustomerPOCO customer)
 {
 if (string.IsNullOrEmpty(customer.id) || string.IsNullOrEmpty(customer.customername))
 {
 return new BadRequestResult();
 }
 string requestId = Guid.NewGuid().ToString();
 string statusUrl = $"https://{Environment.GetEnvironmentVariable("WEBSITE_HOSTNAME")}/api/RequestStatus/{requestId}";
 var messagePayload = JsonConvert.SerializeObject(customer);
 var message = new ServiceBusMessage(messagePayload);
 message.ApplicationProperties.Add("RequestGUID", requestId);
 message.ApplicationProperties.Add("RequestSubmittedAt", DateTime.UtcNow);
 message.ApplicationProperties.Add("RequestStatusURL", statusUrl);
 var sender = _serviceBusClient.CreateSender("outqueue");
 await sender.SendMessageAsync(message);
 req.HttpContext.Response.Headers["Retry-After"] = "5";
 return new AcceptedResult(statusUrl, null);
 }
}
```
### AsyncProcessingBackgroundWorker function

The `AsyncProcessingBackgroundWorker` function reads the operation from the queue, processes it based on the message payload, and writes the result to a storage account.

```
public class AsyncProcessingBackgroundWorker(BlobContainerClient _blobContainerClient)
{
 [Function("AsyncProcessingBackgroundWorker")]
 public async Task Run(
 [ServiceBusTrigger("outqueue", Connection = "ServiceBusConnection")] ServiceBusReceivedMessage message)
 {
 // Perform an action against the blob data source for the async readers to check against.
 // This is where your service worker processing will be performed.
 var requestGuid = message.ApplicationProperties["RequestGUID"].ToString();
 string blobName = $"{requestGuid}.blobdata";
 var blobClient = _blobContainerClient.GetBlobClient(blobName);
 using (MemoryStream memoryStream = new MemoryStream())
 using (StreamWriter writer = new StreamWriter(memoryStream))
 {
 writer.Write(message.Body.ToString());
 writer.Flush();
 memoryStream.Position = 0;
 await blobClient.UploadAsync(memoryStream, overwrite: true);
 }
 }
}
```
### AsyncOperationStatusChecker function

The `AsyncOperationStatusChecker` function implements the status endpoint. This function checks the status of the request:

- If the request completes, the function returns HTTP 303 (See Other) and redirects the client to a valet key URL for the result.
- If the request is pending, the function returns an HTTP 200 code that includes the current state.

```
public class AsyncOperationStatusChecker(ILogger<AsyncOperationStatusChecker> _logger)
{
 [Function("AsyncOperationStatusChecker")]
 public async Task<IActionResult> Run(
 [HttpTrigger(AuthorizationLevel.Anonymous, "get", Route = "RequestStatus/{requestId}")] HttpRequest req,
 [BlobInput("data/{requestId}.blobdata", Connection = "DataStorage")] BlockBlobClient inputBlob, string requestId)
 {
 OnCompleteEnum OnComplete = Enum.Parse<OnCompleteEnum>(req.Query["OnComplete"].FirstOrDefault() ?? "Redirect");
 OnPendingEnum OnPending = Enum.Parse<OnPendingEnum>(req.Query["OnPending"].FirstOrDefault() ?? "OK");
 _logger.LogInformation("Received status request for {RequestId} - OnComplete {OnComplete} - OnPending {OnPending}",
 requestId, OnComplete, OnPending);
 // Check whether the blob exists.
 if (await inputBlob.ExistsAsync())
 {
 // If the blob exists, the function uses the OnComplete parameter to determine the next action.
 return await OnCompleted(OnComplete, inputBlob, requestId, req);
 }
 else
 {
 // If the blob doesn't exist, the function uses the OnPending parameter to determine the next action.
 switch (OnPending)
 {
 case OnPendingEnum.OK:
 {
 // Return an HTTP 200 status code.
 return new OkObjectResult(new { status = "In progress", Location = rqs });
 }
 case OnPendingEnum.Synchronous:
 {
 // Long polling example: hold the connection open and check for completion
 // using exponential backoff. Time out after approximately one minute.
 int backoff = 250;
 while (!await inputBlob.ExistsAsync() && backoff < 64000)
 {
 _logger.LogInformation("Synchronous mode {RequestId} - retrying in {Backoff} ms", requestId, backoff);
 backoff = backoff * 2;
 await Task.Delay(backoff);
 }
 if (await inputBlob.ExistsAsync())
 {
 _logger.LogInformation("Synchronous mode {RequestId} - completed after {Backoff} ms", requestId, backoff);
 return await OnCompleted(OnComplete, inputBlob, requestId, req);
 }
 else
 {
 _logger.LogInformation("Synchronous mode {RequestId} - NOT FOUND after timeout {Backoff} ms", requestId, backoff);
 return new NotFoundResult();
 }
 }
 default:
 {
 throw new InvalidOperationException($"Unexpected value: {OnPending}");
 }
 }
 }
 }
 private async Task<IActionResult> OnCompleted(OnCompleteEnum OnComplete, BlockBlobClient inputBlob, string requestId, HttpRequest req)
 {
 switch (OnComplete)
 {
 case OnCompleteEnum.Redirect:
 {
 // Generate a user delegation SAS URI by using managed identity credentials.
 BlobServiceClient blobServiceClient = inputBlob.GetParentBlobContainerClient().GetParentBlobServiceClient();
 var userDelegationKey = await blobServiceClient.GetUserDelegationKeyAsync(DateTimeOffset.UtcNow, DateTimeOffset.UtcNow.AddDays(7));
 // Return 303 (See Other) to redirect the client to the result resource.
 // GenerateUserDelegationSasUri is a custom helper. See the full implementation on GitHub.
 req.HttpContext.Response.Headers.Location = GenerateUserDelegationSasUri(inputBlob, userDelegationKey);
 return new StatusCodeResult(StatusCodes.Status303SeeOther);
 }
 case OnCompleteEnum.Stream:
 {
 // Download the file and return it directly to the caller.
 // For larger files, use a stream to minimize RAM usage.
 return new OkObjectResult(await inputBlob.DownloadContentAsync());
 }
 default:
 {
 throw new InvalidOperationException($"Unexpected value: {OnComplete}");
 }
 }
 }
}
public enum OnCompleteEnum
{
 Redirect,
 Stream
}
public enum OnPendingEnum
{
 OK,
 Synchronous
}
```

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Create a separate backend service for each frontend interface, instead of using a single general-purpose backend for all of them. This approach lets you tailor each backend to the needs of its frontend and avoids the bottleneck of a shared backend that serves multiple interfaces. It's based on the Backends for Frontends pattern by Sam Newman.

## Context and problem

Consider an application that's initially designed with a desktop web UI and a corresponding backend service. As business requirements change over time, a mobile interface is added. Both interfaces interact with the same backend service. But the capabilities of a mobile device differ significantly from a desktop browser in terms of screen size, performance, and display limitations.

A backend service frequently encounters competing demands from multiple frontend systems. These demands result in frequent updates and potential development bottlenecks. Conflicting updates and the need to maintain compatibility result in excessive demand on a single deployable resource.

Having a separate team manage the backend service can create a disconnect between frontend and backend teams. This disconnect can cause delays in gaining consensus and balancing requirements. For example, changes requested by one frontend team must be validated with other frontend teams before integration.

## Solution

Introduce a new layer that handles only the requirements that are specific to the interface. This layer, known as the backend-for-frontend (BFF) service, sits between the frontend client and the backend service. If the application supports multiple interfaces, such as a web interface and a mobile app, create a BFF service for each interface.

This pattern customizes the client experience for a specific interface without affecting other interfaces. It also optimizes performance to meet the needs of the frontend environment. Because each BFF service is smaller and less complex than a shared backend service, it can make the application easier to manage.

Frontend teams independently manage their own BFF service, which gives them control over language selection, release cadence, workload prioritization, and feature integration. This autonomy enables them to operate efficiently without depending on a centralized backend development team.

Many BFF services traditionally relied on REST APIs, but GraphQL implementations are emerging as an alternative. With GraphQL, the querying mechanism eliminates the need for a separate BFF layer because it allows clients to request the data that they need without relying on predefined endpoints.

For more information, see Backends for Frontends pattern by Sam Newman.

## Problems and considerations

- Evaluate your optimal number of services depending on the associated costs. Maintaining and deploying more services means increased operational overhead. Each individual service has its own life cycle, deployment and maintenance requirements, and security needs.
- Review the service-level objectives when you add a new service. Increased latency might occur because clients aren't contacting your services directly, and the new service introduces an extra network hop.
- When different interfaces make the same requests, evaluate whether the requests can be consolidated into a single BFF service. Sharing a single BFF service between multiple interfaces can result in different requirements for each client, which can complicate the BFF service's growth and support. - Code duplication is a probable outcome of this pattern. Evaluate the trade-off between code duplication and a better-tailored experience for each client.
- The BFF service should only handle client-specific logic related to a specific user experience. Cross-cutting features, such as monitoring and authorization, should be abstracted to maintain efficiency. Typical features that might surface in the BFF service are handled separately with the Gatekeeping, Rate Limiting, and Routing patterns.
- Consider how learning and implementing this pattern affects the development team. Developing new backend systems requires time and effort, which can result in technical debt. Maintaining the existing backend service adds to this challenge.
- Evaluate whether you need this pattern. For example, if your organization uses GraphQL with frontend specific resolvers, BFF services might not add value to your applications. - Another scenario is an application that combines an API gateway with microservices. This approach might be sufficient for some scenarios where BFF services are typically recommended.

## When to use this pattern

Use this pattern when:

- A shared or general-purpose backend service requires substantial development overhead to maintain.
- You want to optimize the backend for the requirements of specific client interfaces.
- You make customizations to a general-purpose backend to accommodate multiple interfaces.
- A programming language is better suited for the backend of a specific user interface, but not all user interfaces.

This pattern might not be suitable when:

- Interfaces make the same or similar requests to the backend.
- Only one interface interacts with the backend.

## Workload design

Evaluate how to use the Backends for Frontends pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | When you isolate services to a specific frontend interface, you contain malfunctions. The availability of one client doesn't affect the availability of another client's access. When you treat various clients differently, you can prioritize reliability efforts based on expected client access patterns. - RE:02 Critical flows - RE:07 Self-preservation |
| Security design decisions help ensure the confidentiality,integrity, andavailabilityof your workload's data and systems. | The service separation introduced in this pattern allows security and authorization in the service layer to be customized for each client's specific needs. This approach can reduce the API's surface area and limit lateral movement between backends that might expose different capabilities. - SE:04 Segmentation - SE:08 Hardening resources |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | The backend separation enables you to optimize in ways that might not be possible with a shared service layer. When you handle individual clients differently, you can optimize performance for a specific client's constraints and functionality. - PE:02 Capacity planning - PE:09 Critical flows |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

This example demonstrates a use case for the pattern in which two distinct client applications, a mobile app and a desktop application, interact with Azure API Management (data plane gateway). This gateway serves as an abstraction layer and manages common cross-cutting concerns such as:

- **Authorization.**Ensures that only verified identities with the proper access tokens can call protected resources by using API Management with Microsoft Entra ID.
- **Monitoring.**Captures and sends request and response details to Azure Monitor for observability purposes.
- **Request caching.**Optimizes repeated requests by serving responses from cache by built-in features of API Management.
- **Routing and aggregation.**Directs incoming requests to the appropriate BFF services.

Each client has a dedicated BFF service running as an Azure function that serves as an intermediary between the gateway and the underlying microservices. These client-specific BFF services ensure a tailored experience for pagination and other functionalities. The mobile app prioritizes bandwidth efficiency and takes advantage of caching to enhance performance. In contrast, the desktop application retrieves multiple pages in a single request, which creates a more immersive user experience.

### Flow A for the first page request from the mobile client

- The mobile client sends a - `GET`request for page- `1`, including the OAuth 2.0 token in the authorization header.
- The request reaches the API Management gateway, which intercepts it and: - **Checks the authorization status.**API Management implements defense in depth, so it checks the validity of the access token.
- **Streams the request activity as logs to Azure Monitor.**Details of the request are recorded for auditing and monitoring.

- The policies are enforced, then API Management routes the request to the Azure function mobile BFF service.
- The Azure function mobile BFF service then interacts with the necessary microservices to fetch a single page and process the requested data to provide a lightweight experience.
- The response is returned to the client.

### Flow B for the first page cached request from the mobile client

- The mobile client sends the same - `GET`request for page- `1`again, including the OAuth 2.0 token in the authorization header.
- The API Management gateway recognizes that this request was made before and: - **The policies are enforced.**Then the gateway identifies a cached response that matches the request parameters.
- **Returns the cached response immediately.**This quick response eliminates the need to forward the request to the Azure function mobile BFF service.

### Flow C for the first request from the desktop client

- The desktop client sends a - `GET`request for the first time, including the OAuth 2.0 token in the authorization header.
- The request reaches the API Management gateway, where cross-cutting concerns are handled. - **Authorization:**Token validation is always required.
- **Stream the request activity:**Request details are recorded for observability.

- After all policies are enforced, API Management routes the request to the Azure function desktop BFF service, which handles the data-heavy application processing. The desktop BFF service aggregates multiple requests by using underlying microservices calls before responding to the client with a multiple page response.

### Design

- Microsoft Entra ID serves as the cloud-based identity provider. It provides tailored audience claims for both mobile and desktop clients. These claims are then used for authorization.
- API Management serves as a proxy between the clients and their BFF services, which establishes a perimeter. API Management is configured with policies to validate the JSON Web Tokens and rejects requests that lack a token or contain invalid claims for the targeted BFF service. It also streams all the activity logs to Azure Monitor.
- Azure Monitor functions as the centralized monitoring solution. It aggregates all activity logs to ensure comprehensive, end-to-end observability.
- Azure Functions is a serverless solution that efficiently exposes BFF service logic across multiple endpoints, which simplifies development. Azure Functions also minimizes infrastructure overhead and helps lower operational costs.

## Next steps

- Backends for Frontends pattern by Sam Newman
- Authentication and authorization to APIs in API Management
- How to integrate API Management with Application Insights

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Isolate the elements of an application into pools so that if one element fails, the others continue to function. This approach, also known as a *cell-based architecture*, makes an application tolerant of failure and stops a fault in one part of the system from cascading across the rest.

Tip

This pattern is named after the sectioned partitions (bulkheads) of a ship's hull. If the hull of a ship is compromised, only the damaged section fills with water, which prevents the ship from sinking.

## Context and problem

A cloud-based application might include multiple services, and each service has one or more consumers. Excessive load or failure in a service affects all consumers of the service.

Also, a consumer might send requests to multiple services simultaneously and use resources for each request. When the consumer sends a request to a misconfigured or unresponsive service, the resources that the client's request uses might remain unavailable for an extended period. As requests to the service continue, those resources might be exhausted. For example, the client's connection pool might be exhausted. At that point, the consumer's requests to other services are affected. Eventually, the consumer can't send requests to any other services, not only the original unresponsive service.

Resource exhaustion affects services that have multiple consumers. Many requests from one client might exhaust available resources in the service. Resource exhaustion can mean that other consumers can't consume the service, which causes a cascading failure effect.

## Solution

Partition service instances into different groups based on consumer load and availability requirements. This design helps isolate failures. You can sustain service functionality for some consumers, even during a failure.

A consumer can also partition resources to ensure that resources used to call one service don't affect the resources used to call another service. For example, a consumer that calls multiple services might be assigned a connection pool for each service. If a service begins to fail, it only affects the connection pool assigned for that service. The consumer can continue to use other services.

This pattern provides the following benefits:

- Isolates consumers and services from cascading failures. A problem that affects a consumer or service can be isolated within its own bulkhead to prevent the entire solution from failing.
- Preserves some functionality if a service failure occurs. Other services and features of the application continue to work.
- Provides different quality of service levels for consuming applications. You can configure a high-priority consumer pool to use high-priority services.

The following diagram shows bulkheads structured around connection pools that call individual services. If Service A fails or causes a problem, the connection pool is isolated, so only workloads that use the thread pool assigned to Service A are affected. Workloads that use Service B and C aren't affected and can continue working without interruption.

Diagram that shows two workloads, Workload 1 and Workload 2, and three services, Service A, Service B, and Service C. Workload 1 uses a connection pool that's assigned to Service A. Workload 2 uses two connection pools. One connection pool is assigned to Service B, and the other is assigned to Service C. The connection pool that Workload 1 uses is isolated. The connection pools that Workload 2 uses can continue to call Service B and Service C.

The following diagram shows multiple clients that call a single service. Each client is assigned to a separate service instance. Client 1 makes too many requests and overwhelms its instance. Because each service instance is isolated from the others, the other clients can continue to make calls.

Diagram that shows three clients, Client 1, Client 2, and Client 3, and three service instances that each form a part of Service A. Each client connects to its own service instance. The service instances are isolated. If Client 1 overwhelms its instance, Clients 2 and 3 are unaffected.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- Define partitions around the business and technical requirements of the application.
- If you use tactical domain-driven design to design microservices, partition boundaries should align with the bounded contexts.
- When you partition services or consumers into bulkheads, consider the level of isolation offered by the technology and the overhead in terms of cost, performance, and manageability.
- To provide more sophisticated fault handling, consider combining bulkheads with retry, circuit breaker, and throttling patterns.
- When you partition consumers into bulkheads, consider using processes, thread pools, and semaphores. Projects like resilience4j and Polly offer a framework for creating consumer bulkheads.
- When you partition services into bulkheads, consider deploying them into separate virtual machines, containers, or processes. Containers offer a good balance of resource isolation with fairly low overhead.
- Services that communicate by using asynchronous messages can be isolated through different sets of queues. Each queue can have a dedicated set of instances that process messages on the queue or a single group of instances that use an algorithm to dequeue and dispatch processing.
- Determine the level of granularity for the bulkheads. For example, if you want to distribute tenants across partitions, you can place each tenant into a separate partition or put several tenants into one partition.
- Monitor each partition's performance and service-level agreement (SLA).
- Use built-in platform controls, such as Azure API Management rate limits, Azure Cosmos DB request unit (RU) isolation, and resource limits in Azure Kubernetes Service (AKS) or Azure Container Apps. Don't re-create these throttling and isolation mechanisms in your application code.
- AI and inference workloads often require strict bulkheads because of deployment-level quotas and concurrency limits. For example, isolate model deployments or Foundry resources per workload or per tenant.

## When to use this pattern

Use this pattern when:

- You want to isolate resources for specific dependencies so that a disruption in one service doesn't affect the entire application.
- You want to isolate critical consumers from standard consumers.
- You need to protect the application from cascading failures.

This pattern might not be suitable when:

- Less efficient use of resources might not be acceptable in the project.
- The added complexity isn't necessary.

## Workload design

Evaluate how to use the Bulkhead pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | The failure isolation strategy introduced through the intentional and complete segmentation between components attempts to contain faults to the bulkhead that experiences the problem, which prevents impact on other bulkheads. - RE:02 Critical flows - RE:07 Self-preservation |
| Security design decisions help ensure the confidentiality,integrity, andavailabilityof your workload's data and systems. | The segmentation between components helps constrain security incidents to the compromised bulkhead. - SE:04 Segmentation |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | Each bulkhead can be individually scalable to efficiently meet the needs of the task that's encapsulated in the bulkhead. - PE:02 Capacity planning - PE:05 Scaling and partitioning |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

The following Kubernetes configuration file creates an isolated container to run a single service, with its own CPU and memory resources and limits.

```
apiVersion: v1
kind: Pod
metadata:
 name: drone-management
spec:
 containers:
 - name: drone-management-container
 image: drone-service
 resources:
 requests:
 memory: "64Mi"
 cpu: "250m"
 limits:
 memory: "128Mi"
 cpu: "1"
```
## Next steps

- Use API Management rate-limit policies to control request throughput per client.
- Use Azure Functions concurrency controls to limit parallel executions.
- Set Container Apps resource limits to control CPU and memory per workload.
- Assign Azure Cosmos DB RU throughput per container for predictable isolation.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Load data on demand into a cache from a data store. This approach can improve performance and helps maintain consistency between data held in the cache and data in the underlying data store.

## Context and problem

Applications use a cache to improve performance for repeated access to information in a data store. But cached data can't always remain consistent with the data store. Applications should implement a strategy that keeps the data in the cache as up-to-date as possible. The strategy should also detect when cached data becomes stale and handle it appropriately.

## Solution

Many commercial caching systems provide read-through and write-through or write-behind operations. In these systems, an application retrieves data by referencing the cache. If the data isn't in the cache, the application retrieves it from the data store and adds it to the cache. The system automatically writes any changes made to cached data back to the data store.

For caches that don't provide this functionality, the applications that use the cache must maintain the data.

An application can emulate the functionality of read-through caching by implementing the Cache-Aside pattern. This strategy loads data into the cache on demand. The following diagram uses the Cache-Aside pattern to store data in the cache.

- The application determines whether an item currently resides in the cache by attempting to read from the cache.
- If the item isn't in the cache, also known as a - *cache miss*, the application retrieves the item from the data store.
- The application adds the item to the cache and then returns it to the caller.

When an application updates information, it writes the change to the data store and then invalidates the corresponding item in the cache. When the item is needed again, the Cache-Aside pattern retrieves the updated data from the data store and adds it to the cache.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- **Lifetime of cached data:**Many caches use an expiration policy to invalidate data and remove it from the cache if it isn't accessed for a set period. To make cache-aside effective, ensure that the expiration policy matches the pattern of access for applications that use the data. Don't make the expiration period too short because premature expiration can cause applications to continually retrieve data from the data store and add it to the cache. Similarly, don't make the expiration period so long that the cached data becomes stale. Caching works best for relatively static data or data that applications read frequently.
- **Evicting data:**Most caches have a limited size compared to the data store where the data originates. If the cache exceeds its size limit, it evicts data. Most caches adopt a least-recently-used policy to select items for eviction, but some allow customization.
- **Configuration:**You can configure cache behavior globally or per cached item. A single global eviction policy might not suit all items. If an item is expensive to retrieve, configure the cache item individually. In this situation, it makes sense to keep the item in the cache, even if it gets accessed less frequently than cheaper items.
- **Priming the cache:**Many solutions prepopulate the cache with data that an application likely requires as part of the startup processing. The Cache-Aside pattern remains useful when some of this data expires or gets evicted.
- **Consistency:**The Cache-Aside pattern doesn't guarantee consistency between the data store and the cache. For example, an external process can change an item in the data store at any time. This change doesn't appear in the cache until the item loads again. In a system that replicates data across data stores, frequent synchronization can make consistency challenging.
- **Staleness after writes:**The Cache-Aside pattern invalidates the cached entry on write and repopulates it on the next read. Between the write and the next read, a reader can experience a cache miss or briefly see stale data. This behavior distinguishes the Cache-Aside pattern from write-through caching, which updates the data store and the cache in the same write operation so that readers see the new value immediately after a successful write. If read-heavy paths need read-after-write freshness, consider write-through caching instead.
- **Local caching:**A cache can be local to an application instance and be stored in-memory. Cache-aside works well in this environment if an application repeatedly accesses the same data. But a local cache is private, so different application instances can each have a copy of the same cached data. This data can quickly become inconsistent between caches, so you might need to expire data in a private cache and refresh it more frequently. In these scenarios, consider using a shared or distributed caching mechanism.
- **Semantic caching:**Some workloads can benefit from doing cache retrieval based on semantic meaning rather than exact keys. This approach reduces the number of requests and tokens sent to language models. Only use semantic caching when the data supports semantic equivalence, doesn't risk returning unrelated responses, and doesn't contain private and sensitive data. For example, "What is my yearly take home salary?" is semantically similar to "What is my yearly take home pay?" But if different users ask these questions, the answers should differ. You also shouldn't include this sensitive data in your cache.

## When to use this pattern

Use this pattern when:

- A cache doesn't provide native read-through and write-through operations.
- Resource demand is unpredictable. This pattern enables applications to load data on demand. It doesn't assume which data an application requires in advance.

This pattern might not be suitable when:

- The data is sensitive or security related. Storing data in a cache might be inappropriate, especially when multiple applications or users share the cache. Always retrieve this type of data from the primary source.
- The cached data set is static. If the data fits into the available cache space, prime the cache with the data on startup and apply a policy that prevents the data from expiring.
- Most requests don't experience a cache hit. In this situation, the overhead of checking the cache and loading data into it might outweigh the benefits of caching.
- You cache session state information in a web application hosted in a web farm. In this environment, avoid introducing dependencies based on client-server affinity.

## Workload design

Evaluate how to use the Cache-Aside pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | Caching replicates data. In limited ways, it can preserve the availability of frequently accessed data if the origin data store becomes temporarily unavailable. If the cache malfunctions, the workload can fall back to the origin data store. - RE:05 Redundancy |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | Caching improves performance for read-heavy data that changes infrequently and tolerates some staleness. - PE:08 Data performance - PE:12 Continuous performance optimization |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

Consider using Azure Managed Redis to create a distributed cache that multiple application instances can share.

The following example uses the StackExchange.Redis client, which is a Redis client library written for .NET. To connect to an Azure Managed Redis instance, call the static `ConnectionMultiplexer.Connect` method and pass in the connection string. The method returns a `ConnectionMultiplexer` that represents the connection.

One way to share a `ConnectionMultiplexer` instance in your application is to have a static property that returns a connected instance, similar to the following example. This approach provides a thread-safe way to initialize only a single connected instance.

```
private static ConnectionMultiplexer Connection;
// Redis connection string information
private static Lazy<ConnectionMultiplexer> lazyConnection = new Lazy<ConnectionMultiplexer>(() =>
{
 string cacheConnection = ConfigurationManager.AppSettings["CacheConnection"].ToString();
 return ConnectionMultiplexer.Connect(cacheConnection);
});
public static ConnectionMultiplexer Connection => lazyConnection.Value;
```
The `GetMyEntityAsync` method in the following example shows an implementation of the Cache-Aside pattern. This method retrieves an object from the cache by using the read-through approach.

The method identifies an object by using an integer ID as the key. It tries to retrieve an item from the cache by using this key. If the cache contains a matching item, it returns the item. If the cache doesn't contain a match, the `GetMyEntityAsync` method retrieves the object from a data store, adds it to the cache, and then returns it. This example omits the code that reads the data from the data store because that logic depends on the data store. The cached item is configured to expire to prevent it from becoming stale if another service or process updates it.

```
// Set five minute expiration as a default
private const double DefaultExpirationTimeInMinutes = 5.0;
public async Task<MyEntity> GetMyEntityAsync(int id)
{
 // Define a unique key for this method and its parameters.
 var key = $"MyEntity:{id}";
 var cache = Connection.GetDatabase();
 // Try to get the entity from the cache.
 var json = await cache.StringGetAsync(key).ConfigureAwait(false);
 var value = string.IsNullOrWhiteSpace(json)
 ? default(MyEntity)
 : JsonConvert.DeserializeObject<MyEntity>(json);
 if (value == null) // Cache miss
 {
 // If there's a cache miss, get the entity from the original store and cache it.
 // Code has been omitted because it is data store dependent.
 value = ...;
 // Avoid caching a null value.
 if (value != null)
 {
 // Put the item in the cache with a custom expiration time that
 // depends on how critical it is to have stale data.
 await cache.StringSetAsync(key, JsonConvert.SerializeObject(value)).ConfigureAwait(false);
 await cache.KeyExpireAsync(key, TimeSpan.FromMinutes(DefaultExpirationTimeInMinutes)).ConfigureAwait(false);
 }
 }
 return value;
}
```
Note

The examples use Azure Managed Redis to access the store and retrieve information from the cache. For more information, see Create an Azure Managed Redis instance and Use Azure Managed Redis in .NET.

The following `UpdateEntityAsync` method demonstrates how to invalidate an object in the cache when the application changes the value. The code updates the original data store and then removes the cached item from the cache.

```
public async Task UpdateEntityAsync(MyEntity entity)
{
 // Update the object in the original data store.
 await this.store.UpdateEntityAsync(entity).ConfigureAwait(false);
 // Invalidate the current cache object.
 var cache = Connection.GetDatabase();
 var id = entity.Id;
 var key = $"MyEntity:{id}"; // The key for the cached object.
 await cache.KeyDeleteAsync(key).ConfigureAwait(false); // Delete this key from the cache.
}
```
Note

The order of the steps is important. Update the data store *before* removing the item from the cache. If you remove the cached item first, there's a small window of time when a client might fetch the item before the data store is updated. In this situation, the fetch results in a cache miss because the item isn't in the cache. The cache miss causes the application to retrieve the outdated item from the data store and add it back to the cache. This sequence leads to stale data in the cache.

## Next steps

- Data consistency primer: This primer describes problems with consistency across distributed data. It also summarizes how an application can implement eventual consistency to maintain the availability of data. Cloud applications typically store data across multiple data stores and locations. You must efficiently manage and maintain data consistency in this environment, particularly because of concurrency and availability problems that can arise.
- Use Azure Managed Redis as a semantic cache: This tutorial shows you how to implement semantic caching by using Azure Managed Redis.

## Related resources

- Caching guidance: This guidance provides more information about how to cache data in a cloud solution, and problems to consider when you implement a cache.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Have each service decide when and how to process a business operation, instead of depending on a central orchestrator. This approach decentralizes workflow logic and distributes responsibilities across the components of a system.

## Context and problem

You typically divide a cloud-based application into several small services that work together to process an end-to-end business transaction. A single operation within a transaction can result in multiple point-to-point calls among all services. Ideally, those services are loosely coupled. It's challenging to design a distributed, efficient, and scalable workflow because it involves complex interservice communication.

A common pattern for communication is to use a centralized service or an *orchestrator*. Incoming requests flow through the orchestrator as it delegates operations to the respective services. Each service completes their responsibility and isn't aware of the overall workflow.

You typically implement the orchestrator pattern as custom software that has domain knowledge about the responsibilities of the services within the system. One benefit of this approach is that the orchestrator can consolidate the status of a transaction based on the results of individual operations that the downstream services conduct.

This approach also creates some obstacles. Adding or removing services might break existing logic because you need to rewire portions of the communication path. This dependency makes orchestrator implementation complex and hard to maintain. The orchestrator might negatively affect the workload's reliability. Under load, it can introduce performance bottlenecks and be the single point of failure (SPoF). When the orchestrator fails or becomes overloaded, the failure can propagate to all dependent downstream services.

## Solution

Delegate the transaction-handling logic among the services. Let each service participate in the communication workflow for a business operation and decide when and how to process it.

The Choreography pattern minimizes the dependency on custom software that centralizes the communication workflow. The components implement common logic as they choreograph the workflow among themselves without directly communicating with each other.

A common way to implement choreography is to use a message broker that buffers requests until downstream components claim and process them. The following image shows request handling through a publisher-subscriber model.

- Client requests queue as messages in a message broker.
- The services or the subscriber polls the broker to determine whether it can process that message based on its implemented business logic. The broker can also push messages to subscribers interested in that message.
- Each subscribed service does its operation as the message indicates and responds to the broker with an operation success or failure message.
- If the operation is successful, the service can publish a message to the same queue or a different message queue so that another service can continue the workflow if needed. If the operation fails, the service publishes a failure message. Services that subscribe to that message can run predefined compensating actions for the failed operation or the entire transaction.

## Issues and considerations

Consider the following points when deciding how to implement this pattern:

- **Failure handling complexity.**Components in an application might manage atomic tasks and depend on other parts of the system. Failure in one component can affect other components, which might cause delays in completing the overall request.- To handle failures gracefully, you implement failure-handling logic, which introduces complexity. Failure-handling logic, such as compensating transactions, is also prone to failures.
- **Sequential processes.**This pattern suits a workflow that processes independent business operations in parallel. The workflow can become complicated when choreography needs to occur in a sequence. For example, Service D can start its operation only after Service B and Service C complete their operations successfully.
- **Observability at scale.**This pattern presents challenges if the number of services grows rapidly. Many independent moving parts complicates the workflow between services. Without a central orchestrator holding the full transaction state, no single component has a complete view of an in-flight business operation. You must consistently use distributed tracing and correlation identifiers to maintain observability.
- **Resiliency handler communication.**In an orchestrator-led design, the central component can delegate resiliency responsibilities, such as retry handling for transient, nontransient, and timeout failures, to a dedicated resiliency handler.- When you remove the orchestrator in a choreography-based design, downstream components don't assume resiliency responsibilities. They remain centralized in the resiliency handler. But downstream components must communicate with that handler directly, which increases point-to-point communication.
- **Event schema evolution.**Event schema evolution can cause breaking changes in consumers over time. In this pattern, multiple independent services consume the same events. If a producer changes the data structure of an event, it can break downstream consumers that depend on the old schema. Use a schema registry to manage event contracts and use backward-compatible evolution as services evolve independently.
- **Idempotency and event ordering.**At-least-once delivery and retries can produce duplicate messages, and concurrent consumers can process messages out of order. Design consumers to be idempotent by tracking stable message identifiers. When ordered processing is required, use broker features such as Service Bus sessions or include sequence or version data that lets consumers reject stale events and detect gaps.
- **Atomic state and event publication.**A service that updates its data store and publishes an event in separate operations can commit one operation while the other fails. Use the Transactional Outbox pattern or an equivalent atomic mechanism to persist the state change and event together before a separate process publishes the event.
- **Emergent behavior and event storms.**Decentralized event topologies can create emergent behavior at scale. When many services react to each other's events, the system can unintentionally produce feedback loops or event storms. A minor event might trigger a cascade of downstream reactions. To prevent circular event chains, use guardrails like event filtering, consumer concurrency limits, throttling, and explicit rules.

## When to use this pattern

Use this pattern when:

- The downstream components handle atomic operations independently in a - *fire and forget*approach. Each component completes a task and then signals completion to other components through the message broker. The initiating service doesn't actively manage or track the task after dispatching it, but downstream services still communicate outcomes through events.
- You expect to frequently update and replace the components. This pattern lets you modify the application with less effort and minimal disruption to existing services.
- You use serverless architectures for simple workflows. The components can be short-lived and event-driven. When an event occurs, the service creates components that do a task, and the service removes components after they complete that task.
- Communication between bounded contexts requires loose coupling across domain boundaries. For communication inside a single bounded context, consider an orchestrator pattern instead, depending on the complexity and team preference.
- The central orchestrator introduces a performance bottleneck.

This pattern might not be suitable when:

- The application is complex and requires a central component to handle shared logic to keep the downstream components lightweight.
- Point-to-point communication between the components is inevitable.
- You need to use business logic to consolidate all operations that downstream components handle.

## Workload design

Evaluate how to use the Choreography pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | The distributed components in this pattern are autonomous and designed to be replaceable, so you can modify the workload with less overall change to the system. - OE:04 Tools and processes |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | This pattern provides an alternative when performance bottlenecks occur in a centralized orchestration topology. - PE:02 Capacity planning - PE:05 Scaling and partitioning |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

This example shows the Choreography pattern by creating an event-driven, cloud-native workload that runs functions alongside microservices. When a client requests to ship a package, the workload assigns a drone. After the package is ready for pickup by the scheduled drone, the delivery process starts. While the package is in transit, the workload handles the delivery until it receives the shipped status. For the full reference architecture, see Microservices with Azure Container Apps.

The ingestion service receives client requests and converts them into messages that include the delivery details. Business transactions start after services consume those new messages.

A single client business transaction requires three distinct business operations:

- Create or update a package.
- Assign a drone to deliver the package.
- Handle the delivery, including checking and sending a notification when the package ships.

Package, drone scheduler, and delivery microservices perform the business processing. The services use messaging instead of a central orchestrator to communicate with each other. Each service must implement a protocol in advance that coordinates the business workflow in a decentralized way.

### Design

Services process business transactions in a sequence through multiple hops. Each hop shares a single message bus among all the business services.

When a client sends a delivery request through an HTTP endpoint, the ingestion service receives it, converts it into a message, and then publishes the message to the shared message bus. The subscribed business services consume new messages added to the bus. When a business service receives the message, it completes the operation successfully, or the request fails or times out. If the request succeeds, the service responds to the bus with the `Ok` status code, raises a new operation message, and sends it to the message bus. If the request fails or times out, the service reports the failure reason code to the message bus and dead-letters the message through Azure Service Bus. The service also dead-letters messages that it can't receive or process within a specific amount of time.

This design uses multiple message buses to process the entire business transaction. Azure Service Bus and Azure Event Grid provide the messaging service platform for this design. The workload runs on Azure Container Apps. The ingestion service runs as an Azure Function hosted on Container Apps, while the package, drone scheduler, and delivery services run as microservices in the same Container Apps environment. Container Apps handles event-driven processing that runs the business logic.

This design also ensures that the choreography occurs in a sequence. A single Service Bus namespace contains a topic that has two subscriptions and a session-aware queue. The ingestion service publishes messages to the topic. The package service and drone scheduler service subscribe to the topic and publish messages that notify the queue of successful requests. Include a common session identifier that associates a GUID with the delivery identifier so that the delivery service can correlate the two messages it needs for each transaction. One message confirms the package is ready, the other message confirms a drone is scheduled. Without this session-based correlation, the delivery service has no way to associate related messages across independent hops, because no central coordinator tracks the transaction state. The delivery service waits for two related messages for each transaction. The first message indicates that the package is ready to be shipped, and the second message signals that a drone is scheduled.

In this design, Service Bus handles high-value messages that must not be lost or duplicated during the entire delivery process. When the package ships, a change of state publishes to Event Grid. The event sender has no expectation about how the change of state is handled. Downstream organization services that this design doesn't include can listen for this event type and run specific business logic, such as sending an order-status email to the user.

If you deploy this pattern in another compute service, such as AKS, you can deploy an ambassador as a sidecar in the same pod as the business application. Colocation minimizes communication latency, but the proxy adds processing and resource overhead and scales with the application. Use this approach when you need language-independent connectivity concerns that the platform doesn't provide.

To avoid cascading retry operations that might lead to multiple attempts, business services should immediately flag unacceptable messages. Enrich these messages by using common reason codes or a defined application code so that the services can move them to a DLQ. Consider implementing the Saga pattern to manage consistency problems from downstream services. For example, another service handles dead-letter messages for remediation purposes only by running a compensation, retry, or pivot transaction.

The business services are idempotent to ensure that retry operations don't create duplicate resources. For example, the package service uses upsert operations to add data to the data store.

## Next step

- Review asynchronous messaging options in Azure to learn about the different infrastructure choices available for implementing a decentralized workflow.

## Related resources

Consider these patterns in your design for choreography:

- Use the Ambassador pattern to modularize business service communication with the message bus.
- Implement the Queue-Based Load Leveling pattern to handle spikes in the workload.
- Use asynchronous distributed messaging through the Publisher-Subscriber pattern.
- Use compensating transactions to undo a series of successful operations if one or more related operations fail.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Temporarily block access to a remote service or resource after failures reach a threshold, instead of repeatedly retrying an operation that's likely to fail. This approach handles faults that take varying amounts of time to recover from, lets the failing service recover, and improves the stability and resiliency of an application.

## Context and problem

In a distributed environment, calls to remote resources and services can fail because of transient faults. Transient faults include overcommitted or temporarily unavailable resources, slow network connections, or time-outs. These faults typically correct themselves after a short period of time. To help manage these faults, you should design a cloud application to use a strategy, such as the Retry pattern.

Unanticipated events can create faults that take longer to fix. These faults can range in severity from a partial loss of connectivity to a complete service failure. In these situations, an application shouldn't continually retry an operation that's unlikely to succeed. Instead, the application should quickly recognize the failed operation and handle the failure accordingly.

If a service is busy, failure in one part of the system might lead to cascading failures. For example, you can configure an operation that invokes a service to implement a time-out. If the service fails to respond within this period, the operation replies with a failure message.

However, this strategy can block concurrent requests to the same operation until the time-out period expires. These blocked requests might hold critical system resources, such as memory, threads, and database connections. This problem can exhaust resources, which might fail other unrelated parts of the system that need to use the same resources.

In these situations, an operation should fail immediately and only attempt to invoke the service if it's likely to succeed. To resolve this problem, set a shorter time-out. But ensure that the time-out is long enough for the operation to succeed most of the time.

## Solution

The Circuit Breaker pattern helps prevent an application from repeatedly trying to run an operation that's likely to fail. This pattern enables the application to continue running without waiting for the fault to be fixed or wasting CPU cycles on determining that the fault is persistent. The Circuit Breaker pattern also enables an application to detect when the fault is resolved. If the fault is resolved, the application can try to invoke the operation again.

Note

The Circuit Breaker pattern serves a different purpose than the Retry pattern. The Retry pattern enables an application to retry an operation with the expectation that it eventually succeeds. The Circuit Breaker pattern prevents an application from performing an operation that's likely to fail. An application can combine these two patterns by using the Retry pattern to invoke an operation through a circuit breaker. However, the retry logic should be sensitive to any exceptions that the circuit breaker returns and stop retry attempts if the circuit breaker indicates that a fault isn't transient.

A circuit breaker acts as a proxy for operations that might fail. The proxy should monitor the number of recent failures and use this information to decide whether to allow the operation to proceed or to return an exception immediately.

You can implement the proxy as a state machine that includes the following states. These states mimic the functionality of an electrical circuit breaker:

- **Closed:**The request from the application is routed to the operation. The proxy maintains a count of the number of recent failures. If the call to the operation is unsuccessful, the proxy increments this count. If the number of recent failures exceeds a specified threshold within a given time period, the proxy is placed into the- **Open**state and starts a time-out timer. When the timer expires, the proxy is placed into the- **Half-Open**state.- Note - During the time-out, the system tries to fix the problem that caused the failure before it allows the application to attempt the operation again.
- **Open:**The request from the application fails immediately and an exception is returned to the application.
- **Half-Open:**A limited number of requests from the application are allowed to pass through and invoke the operation. If these requests are successful, the circuit breaker assumes that the fault that caused the failure is fixed, and the circuit breaker switches to the- **Closed**state. The failure counter is reset. If any request fails, the circuit breaker assumes that the fault is still present, so it reverts to the- **Open**state. It restarts the time-out timer so that the system can recover from the failure.- Note - The - **Half-Open**state helps prevent a recovering service from suddenly being flooded with requests. As a service recovers, it might be able to support a limited volume of requests until the recovery is complete. But while recovery is in progress, a flood of work can cause the service to time out or fail again.

The following diagram shows the counter operations for each state.

The failure counter for the **Closed** state is time based. It automatically resets at periodic intervals. This design helps prevent the circuit breaker from entering the **Open** state if it experiences occasional failures. The failure threshold triggers the **Open** state only when a specified number of failures occur during a specified interval.

The success counter for the **Half-Open** state records the number of successful attempts to invoke the operation. The circuit breaker reverts to the **Closed** state after a specified number of successful, consecutive operation invocations. If any invocation fails, the circuit breaker enters the **Open** state immediately and the success counter resets the next time it enters the **Half-Open** state.

Note

System recovery is based on external operations, such as restoring or restarting a failed component or repairing a network connection.

The Circuit Breaker pattern provides stability while the system recovers from a failure and minimizes the impact on performance. It can help maintain the response time of the system. This pattern quickly rejects a request for an operation that's likely to fail, rather than waiting for the operation to time out or never return. If the circuit breaker raises an event each time it changes state, this information can help monitor the health of the protected system component or alert an administrator when a circuit breaker switches to the **Open** state.

You can customize and adapt this pattern to different types of failures. For example, you can apply an increasing time-out timer to a circuit breaker. You can place the circuit breaker in the **Open** state for a few seconds initially. If the failure isn't resolved, increase the time-out to a few minutes and adjust accordingly. In some cases, rather than returning a failure and raising an exception, the **Open** state can return a default value that's meaningful to the application.

Note

Traditionally, circuit breakers relied on preconfigured thresholds, such as failure count and time-out duration. This approach resulted in a deterministic but sometimes suboptimal behavior.

Adaptive techniques that use AI and machine learning can dynamically adjust thresholds based on real-time traffic patterns, anomalies, and historical failure rates. This approach improves resiliency and efficiency.

## Problems and considerations

Consider the following factors when you implement this pattern:

- **Exception handling:**An application that invokes an operation through a circuit breaker must be able to handle the exceptions if the operation is unavailable. Exception management is based on the application. For example, an application might temporarily degrade its functionality, invoke an alternative operation to try to perform the same task or obtain the same data, or report the exception to the user and ask them to try again later.
- **Types of exceptions:**The reasons for a request failure can vary in severity. For example, a request might fail because a remote service crashes and requires several minutes to recover, or because an overloaded service causes a time-out. A circuit breaker might be able to examine the types of exceptions that occur and adjust its strategy based on the nature of these exceptions. For example, it might require a larger number of time-out exceptions to trigger the circuit breaker to the- **Open**state compared to the number of failures caused by the unavailable service.
- **Monitoring:**A circuit breaker should provide clear observability into both failed and successful requests so that operations teams can assess system health. Use distributed tracing for end-to-end visibility across services.
- **Recoverability:**You should configure the circuit breaker to match the likely recovery pattern of the operation that it protects. For example, if the circuit breaker remains in the- **Open**state for a long period, it can raise exceptions even if the reason for the failure is resolved. Similarly, a circuit breaker can fluctuate and reduce the response times of applications if it switches from the- **Open**state to the- **Half-Open**state too quickly.
- **Failed operations testing:**In the- **Open**state, rather than using a timer to determine when to switch to the- **Half-Open**state, a circuit breaker can periodically ping the remote service or resource to determine whether it's available. This ping can either attempt to invoke a previously failed operation or use a special health-check operation that the remote service provides. For more information, see Health Endpoint Monitoring pattern.
- **Manual override:**If the recovery time for a failing operation is extremely variable, you should provide a manual reset option that enables an administrator to close a circuit breaker and reset the failure counter. Similarly, an administrator can force a circuit breaker into the- **Open**state and restart the time-out timer if the protected operation is temporarily unavailable.
- **Concurrency:**A large number of concurrent instances of an application can access the same circuit breaker. The implementation shouldn't block concurrent requests or add excessive overhead to each call to an operation.
- **Resource differentiation:**Be careful when you use a single circuit breaker for one type of resource if there might be multiple underlying independent providers. For example, in a data store that contains multiple shards, one shard might be fully accessible while another experiences a temporary problem. If the error responses in these scenarios are merged, an application might try to access some shards even when failure is likely. And access to other shards might be blocked even though it's likely to succeed.
- **Accelerated circuit breaking:**Sometimes a failure response can contain enough information for the circuit breaker to trip immediately and stay tripped for a minimum amount of time. For example, the error response from a shared resource that's overloaded can indicate that the application should instead try again in a few minutes, instead of immediately retrying.
- **Multiregion deployments:**You can design a circuit breaker for single region or multiregion deployments. To design for multiregion deployments, use global load balancers or custom region-aware circuit breaking strategies that help ensure controlled failover, latency optimization, and regulatory compliance.
- **Service mesh circuit breakers:**You can implement circuit breakers at the application layer or as a cross-cutting, abstracted feature. For example, service meshes often support circuit breaking as a sidecar or as a standalone capability without modifying application code.- Note - A service can return HTTP 429 (too many requests) if it's throttling the client or HTTP 503 (service unavailable) if the service isn't available. The response can include other information, such as the anticipated duration of the delay.
- **Failed request replay:**In the- **Open**state, rather than failing immediately, a circuit breaker can record the details of each request in a journal and arrange for these requests to be replayed when the remote resource or service becomes available.
- **Inappropriate time-outs on external services:**A circuit breaker might not fully protect applications from failures in external services that have long time-out periods. If the time-out is too long, a thread that runs a circuit breaker might be blocked for an extended period before the circuit breaker indicates that the operation failed. During this time, many other application instances might also try to invoke the service through the circuit breaker and tie up numerous threads before they all fail.
- **Adaptability to compute diversification:**Circuit breakers should account for different compute environments, from serverless to containerized workloads, where factors like cold starts and scalability affect failure handling. Adaptive approaches can dynamically adjust strategies based on the compute type, which helps ensure resilience across heterogeneous architectures.

## When to use this pattern

Use this pattern when:

- You want to prevent cascading failures by stopping excessive remote service calls or access requests to a shared resource if these operations are likely to fail.
- You want to route traffic intelligently based on real-time failure signals to enhance multiregion resilience.
- You want to protect against slow dependencies so that you can maintain your service-level objectives and avoid performance degradation from high-latency services.
- You want to manage intermittent connectivity problems and reduce request failures in distributed environments.

This pattern might not be suitable when:

- You need to manage access to local private resources in an application, such as in-memory data structures. In this environment, a circuit breaker adds overhead to your system.
- You need to use it as a substitute for handling exceptions in the business logic of your applications.
- Well-known retry algorithms are sufficient and your dependencies are designed to handle retry mechanisms. In this scenario, a circuit breaker in your application might add unnecessary complexity to your system.
- Waiting for a circuit breaker to reset might introduce unacceptable delays.
- You have a message-driven or event-driven architecture, because they often route failed messages to a dead letter queue for manual or deferred processing. Built-in failure isolation and retry mechanisms are often sufficient.
- Failure recovery is managed at the infrastructure or platform level, such as with health checks in global load balancers or service meshes.

## Workload design

Evaluate how to use the Circuit Breaker pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | This pattern helps prevent a faulting dependency from overloading. Use this pattern to trigger graceful degradation in the workload. Couple circuit breakers with automatic recovery to provide self-preservation and self-healing. - RE:03 Failure mode analysis - Transient faults - RE:07 Self-preservation |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | This pattern avoids the retry-on-error approach, which can lead to excessive resource usage during dependency recovery and can overload performance on a dependency that's attempting recovery. - PE:07 Code and infrastructure - PE:11 Live-issues responses |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

This example implements the Circuit Breaker pattern to help prevent quota overrun by using the Azure Cosmos DB lifetime free tier. This tier is primarily for noncritical data and operates under a capacity plan that allocates a specific quota of resource units per second. During seasonal events, demand might exceed the provided capacity, which can result in `429` responses.

When demand spikes occur, Azure Monitor alerts with dynamic thresholds detect and proactively notify the operations and management teams that the database requires more capacity. Simultaneously, a circuit breaker that's tuned by using historical error patterns trips to prevent cascading failures. In this state, the application gracefully degrades by returning default or cached responses. The application informs users of the temporary unavailability of certain data while preserving overall system stability.

This strategy enhances resilience that aligns with business justification. It controls capacity surges so that workload teams can manage cost increases deliberately and maintain service quality without unexpectedly increasing operating expenses. After demand subsides or increased capacity is confirmed, the circuit breaker resets, and the application returns to full functionality that aligns with both technical and budgetary objectives.

*Download a Visio file of this architecture.*

### Flow A: Closed state

- The system operates normally, and all requests reach the database without returning any - `429`HTTP responses.
- The circuit breaker remains closed, and no default or cached responses are necessary.

### Flow B: Open state

- When the circuit breaker receives the first - `429`response, it trips to an- **Open**state.
- Subsequent requests are immediately short-circuited, which returns default or cached responses and informs users of temporary degradation. The application is protected from further overload.
- Azure Monitor receives logs and telemetry data and evaluates them against dynamic thresholds. An alert triggers if the conditions of the alert rule are met.
- An action group proactively notifies the operations team of the overload condition.
- After workload team approval, the operations team can increase the provisioned throughput to alleviate overload or delay scaling if the load subsides naturally.

### Flow C: Half-Open state

- After a predefined time-out, the circuit breaker enters a - **Half-Open**state that permits a limited number of trial requests.
- If these trial requests succeed without returning - `429`responses, the breaker resets to a- **Closed**state, and normal operations restore back to Flow A. If failures persist, the breaker reverts to the- **Open**state, or Flow B.

### Components

- Azure App Service hosts the web application that serves as the primary entry point for client requests. The application code implements the logic that enforces circuit breaker policies and delivers default or cached responses when the circuit is open. This architecture helps prevent overload on downstream systems and maintain the user experience during peak demand or failures.
- Azure Cosmos DB is one of the application's data stores. It serves noncritical data via the free tier, which is ideal for small production workloads. The circuit breaker mechanism helps limit traffic to the database during high-demand periods.
- Azure Monitor functions as the centralized monitoring solution. It aggregates all activity logs to help ensure comprehensive, end-to-end observability. Azure Monitor receives logs and telemetry data from App Service and key metrics from Azure Cosmos DB (like the number of - `429`responses) for aggregation and analysis.
- Azure Monitor alerts weigh alert rules against dynamic thresholds to identify potential outages based on historical data. Predefined alerts notify the operations team when thresholds are breached. - Sometimes, the workload team might approve an increase in provisioned throughput, but the operations team anticipates that the system can recover on its own because the load isn't too high. In these cases, the circuit breaker time-out elapses naturally. During this time, if the - `429`responses cease, the threshold calculation detects the prolonged outages and excludes them from the learning algorithm. As a result, the next time an overload occurs, the threshold waits for a higher error rate in Azure Cosmos DB, which delays the notification. This adjustment allows the circuit breaker to handle the problem without an immediate alert, which improves cost and operational efficiency.

## Related resources

- The Retry pattern describes how an application can handle anticipated temporary failures when it tries to connect to a service or network resource by transparently retrying an operation that previously failed.
- The Health Endpoint Monitoring pattern describes how a circuit breaker can test the health of a service by sending a request to an endpoint that the service exposes. The service should return information that indicates its status.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Store a large message payload in an external data store and send only a reference token, called a *claim check*, through a messaging system. The token is a unique, obscure key. Receiving applications present the token to the external data store to retrieve the payload. This approach lets workloads transfer large payloads without storing them in the messaging system.

## Context and problem

Traditional messaging systems are optimized to manage a high volume of small messages and often have restrictions on the message size they can handle. Large messages not only risk exceeding these limits but can also degrade the performance of the entire system when the messaging system stores them.

## Solution

Use the Claim-Check pattern, and don't send large messages to the messaging system. Instead, send the payload to an external data store and generate a claim-check token for that payload. The messaging system sends a message with the claim-check token to receiving applications so these applications can retrieve the payload from the data store. The messaging system never sees or stores the payload.

- Payload
- Save payload in data store.
- Generate claim-check token and send message with claim-check token.
- Receive message and read claim-check token.
- Retrieve the payload.
- Process the payload.

## Issues and considerations with the Claim-Check pattern

Consider the following recommendations when implementing the Claim-Check pattern:

- *Delete consumed messages.*If you don't need to archive the message, delete the message and payload after the receiving applications consume it. Use either a synchronous or asynchronous deletion strategy:- *Synchronous deletion*: The consuming application deletes the message and payload immediately after consumption. It ties deletion to the message handling workflow and uses messaging-workflow compute capacity.
- *Asynchronous deletion*: A process outside the message processing workflow deletes the message and payload. It decouples the deletion process from the message handling workflow and minimizes use of messaging-workflow compute.

- *Implement the pattern conditionally.*Incorporate logic in the sending application that applies the Claim-Check pattern if the message size surpasses the messaging system's limit. For smaller messages, bypass the pattern and send the smaller message to the messaging system. This conditional approach reduces latency, optimizes resources utilization, and improves throughput.

## When to use the Claim-Check pattern

The following scenarios are the primary use cases for the Claim-Check pattern:

- *Messaging system limitations*: Use the Claim-Check pattern when message sizes surpass the limits of your messaging system. Offload the payload to external storage. Send only the message with its claim-check token to the messaging system.
- *Messaging system performance*: Use the Claim-Check pattern when large messages are straining the messaging system and degrading system performance.

The following scenarios are secondary use cases for the Claim-Check pattern:

- *Sensitive data protection*: Use the Claim-Check pattern when payloads contain sensitive data that you don't want visible to the messaging system. Apply the pattern to all or portions of sensitive information in the payload. Secure the sensitive data without transmitting it directly through the messaging system.
- *Complex routing scenarios*: Messages traversing multiple components can cause performance bottlenecks due to serialization, deserialization, encryption, and decryption tasks. Use the Claim-Check pattern to prevent direct message processing by intermediary components.

## Workload design with the Claim-Check pattern

An architect should evaluate how the Claim-Check pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure it fullyrecoversafter failure. | Messaging systems don't provide the same reliability and disaster recovery that are often present in dedicated data stores. Separating the data from the message can provide increased reliability for the payload. This separation facilitates data redundancy that allows you to recover payloads after a disaster. - RE:03 Failure mode analysis - RE:09 Disaster recovery |
| Security design decisions help ensure the confidentiality,integrity, andavailabilityof workload data and systems. | The Claim-Check pattern can extract sensitive data from messages and store it in a secure data store. This setup allows you to implement tighter access controls, ensuring that only the services intended to use the sensitive data can access it. At the same time, it hides this data from unrelated services, such as those used for queue monitoring. - SE:03 Data classification - SE:04 Segmentation |
| Cost Optimization is focused on sustaining and improvingyour workload'sreturn on investment. | Messaging systems often impose limits on message size, and increased size limits is often a premium feature. Reducing the size of message bodies might enable you to use a cheaper messaging solution. - CO:07 Component costs - CO:09 Flow costs |
| Performance Efficiency helps your workload efficiently meet demandsby optimizing scaling, data transfer, and code execution. | The Claim-Check pattern improves the efficiency of sending and receiving applications and the messaging system by managing large messages more effectively. It reduces the size of messages sent to the messaging system and ensures receiving applications access large messages only when needed. - PE:05 Scaling and partitioning - PE:12 Continuous performance optimization |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Claim-check pattern examples

The following examples demonstrate how Azure facilitates the implementation of the Claim-Check Pattern:

- *Azure messaging systems*: The examples cover four different Azure messaging system scenarios: Azure Queue Storage, Azure Event Hubs (Standard API), Azure Service Bus, and Azure Event Hubs (Kafka API).
- *Automatic vs. manual claim-check token generation*: These examples also show two methods to generate the claim-check token. In code examples 1-3, Azure Event Grid automatically generates the token when the sending application transfers the payload to Azure Blob Storage. Code example 4 shows a manual token generation process using an executable command-line client.

Choose the example that suits your needs and follow the provided link to view the code on GitHub:

| Sample code | Messaging system scenarios | Token generator | Receiving application | Data store |
|---|---|---|---|---|
| Code example 1 | Azure Queue Storage | Azure Event Grid | Function | Azure Blob Storage |
| Code example 2 | Azure Event Hubs (Standard API) | Azure Event Grid | Executable command-line client | Azure Blob Storage |
| Code example 3 | Azure Service Bus | Azure Event Grid | Function | Azure Blob Storage |
| Code example 4 | Azure Event Hubs (Kafka API) | Executable command-line client | Function | Azure Blob Storage |

## Next steps

- The Enterprise Integration Patterns site has a description of this pattern.
- For another example, see Dealing with large Service Bus messages using Claim-Check pattern (blog post).
- An alternative pattern for handling large messages is Split and Aggregate.
- Libraries like NServiceBus provide support for this pattern out-of-the-box with their DataBus feature.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Use this pattern to undo work when one or more steps fail in an eventually consistent operation. Cloud-hosted applications that implement complex business processes and workflows commonly use operations that follow the eventual consistency model.

## Context and problem

Cloud applications frequently modify data that is spread across various data sources in different geographic locations. To avoid contention and improve performance in a distributed environment, applications should implement eventual consistency instead of strong transactional consistency. In the eventual consistency model, a typical business operation consists of a series of separate steps. During these steps, the overall view of the system state might be inconsistent. But the system should become consistent again when all steps finish.

Handling step failures presents a key challenge in the eventual consistency model. After a failure, you might need to undo work from completed operation steps. However, you can't always roll back the data because other concurrent application instances might change the data. Even when concurrent instances don't change the data, it can be more complex to undo a step than to restore the original state. You might need to apply business-specific rules. For an example, see the travel website example.

When an operation that implements eventual consistency spans multiple data stores, you must access each data store to undo the changes. To prevent the system from remaining inconsistent, you must reliably undo the work in every data store.

An operation that implements eventual consistency doesn't always store its affected data in a database. For example, in a service-oriented architecture (SOA) environment, an operation can invoke an action in a service and change the state that the service holds. To undo the operation, you must also undo this state change, which can involve invoking the service again to reverse the first action's effects.

## Solution

Implement a compensating transaction that undoes the effects of completed steps in the original operation. You might think that you can simply restore the system to its original state, but this approach can overwrite changes from other concurrent application instances. Instead, the compensating transaction must intelligently account for concurrent work. This process is usually application specific and depends on the original operation.

You can use a workflow to implement an eventually consistent operation that requires compensation. As the original operation runs, the system records information about each step and how to undo it. If the operation fails, the workflow rewinds through the completed steps and reverses each step.

While each step is a separate action, together they form an eventually consistent operation. The system must perform the steps and the corresponding undo operations for each step. If the customer cancels, these undo operations can run as a compensating transaction.

A single-step failure doesn't always require you to roll back the entire system by using a compensating transaction. For example, in a travel website scenario, a customer books flights F1, F2, and F3 but fails to reserve a room at hotel H1. Offering the customer a room at a different hotel is preferable to canceling the flights. The customer can still choose to cancel, which triggers the compensating transaction to undo the flight bookings. However, the customer should make this decision, not the system. When decisions are high impact or hard to automate reliably, include a human in the decision-making process.

Consider these important points:

- A compensating transaction might not need to undo the work in the exact reverse order of the original operation.
- You might be able to perform some undo steps in parallel.
- You might need to apply business-specific rules. For example, canceling a flight reservation might not entitle the customer to a complete refund.

This approach is similar to the Saga distributed transactions pattern.

Compensating transactions are eventually consistent operations and can fail. The system should record progress so that it can resume the compensating transaction from the point of failure. A step might run multiple times when retried, so design each step as an idempotent command.

Sometimes manual intervention is the only way to recover from a failed step. In these situations, the system should raise an alert that includes detailed information about the reason for the failure.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- It might not be easy to determine when a step in an operation that implements eventual consistency fails. A step might not fail immediately but instead get blocked. You might need to implement a timeout mechanism.
- It's not easy to generalize compensation logic. A compensating transaction is application specific. It relies on the application having sufficient information to undo the effects of each step in a failed operation.
- Compensating transactions don't always work. Define the steps in a compensating transaction as idempotent commands so that you can repeat them if the compensating transaction itself fails.
- The infrastructure that handles the steps must meet the following criteria: - It's resilient in both the original operation and the compensating transaction.
- It doesn't lose the information required to compensate for a failing step.
- It reliably monitors compensation logic progress. Compensating transactions run after the original operations commit, and other transactions might change intermediate states. Therefore, ensure that you can correlate and audit both the original operation and its compensation end-to-end.

- A compensating transaction doesn't necessarily return the system data to its state at the start of the original operation. Instead, the transaction compensates for the work that the operation completes successfully before it failed.
- The compensating transaction steps don't always reverse the original operation in the exact opposite order. For example, if one data store is more sensitive to inconsistencies than another, undo changes to that store first.
- Some measures can help you improve success rates. You can place a short-term lock with a timeout on each resource that's required to complete an operation. You can acquire these resources in advance, and then perform work only after you acquire all resources. Finalize all actions before the locks expire.
- Retry logic that treats more errors as transient can help minimize failures that trigger a compensating transaction. When a step in an operation that implements eventual consistency fails, handle it as a transient exception and retry the step. Only stop the operation and trigger compensation if the step fails repeatedly or you can't recover it. For more information about retry strategies, see Transient fault handling.
- When you implement a compensating transaction, you face many challenges similar to implementing eventual consistency. For more information, see Minimize coordination.
- Define clear - *points of no return*and irreversible steps. In complex workflows, you can't safely or meaningfully undo some operations, such as external side effects or legally binding actions. Identify compensable versus irreversible steps. Design the workflow so that irreversible steps occur only after all critical validations succeed.

## When to use this pattern

Use this pattern when:

- A business operation spans multiple steps, services, or data stores and must be undone if a later step fails. This scenario often occurs in long‑running workflows that follow an eventual consistency model and can't rely on atomic transactions.
- Failure recovery often requires domain-specific logic rather than a simple data rollback. Use compensating actions when undoing work requires you to apply business rules, such as canceling reservations or issuing partial refunds.

This pattern might not be suitable when:

- Operations can be safely retried and most failures are transient. Retry logic alone is often sufficient in these cases, and compensating transactions add unnecessary complexity.
- The system can't tolerate temporary inconsistency, or compensation can't reliably restore a valid state. Use strong consistency mechanisms or atomic transactions across all steps instead.

## Workload design

Evaluate how to use the Compensating Transaction in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | Compensation actions address malfunctions in critical workload paths by using processes like directly rolling back data changes, breaking transaction locks, or even running native system behavior to reverse the effect. - RE:02 Critical flows - RE:09 Disaster recovery |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

The following diagram shows a practical Azure implementation of the Compensating Transaction pattern. Other implementations might also work for your workload requirements. An orchestrator that runs in Azure Container Apps coordinates each step of a long-running workflow by sending commands through Azure Service Bus. As each forward step succeeds, the orchestrator records both execution state and the corresponding compensating action in Azure Cosmos DB so that the workflow can be resumed, correlated, and audited.

This model uses retries first to preserve forward progress. If a step fails, the orchestrator applies retry logic for transient faults and attempts to continue the original operation. Compensation is invoked only when forward progress becomes impossible, such as when retries are exhausted or the failure is classified as nontransient.

Business-specific rules can also prefer forward progress over immediate compensation. If a step fails, the orchestrator can select an alternative path, such as substituting an equivalent service or fallback option, instead of rolling back the workflow. For high-impact or ambiguous cases, you can pause the workflow for human review before you decide whether to continue on an alternative path or trigger compensation. This approach treats compensation as a last resort and lets domain rules drive recovery decisions.

In a typical sequence, the orchestrator sends step messages through Service Bus (steps 1 and 2), receives successful outcomes, and stores forward and compensation metadata in Azure Cosmos DB.

You can trigger compensation in two ways:

- When a later step in the same workload fails and you must undo previously successful steps. This compensation can happen immediately when a step returns a business error such as a rule-validation failure or after technical retries are exhausted and the message is moved to the dead-letter queue.
- When a subsequent client explicitly requests to cancel a completed operation.

In either case, the orchestrator reads the stored compensation records and sends compensation commands to the corresponding service. If a compensation step fails transiently, Service Bus retries can complete it without escalating the incident.

If repeated retries still fail, Service Bus moves the message to a dead-letter queue and preserves failure details. The orchestrator, or a dedicated dead-letter processor, raises an alert and emits structured telemetry, including failure reason and correlation IDs, to Azure Monitor and Log Analytics, which can surface in Application Insights. This operational path helps teams diagnose failures, determine the need for manual intervention, and maintain traceability across both the original and compensating flows.

The workflow can start compensation automatically for clear, low-risk conditions or pause for human review when the situation is ambiguous, high impact, or requires a manual decision.

Use managed identities and Microsoft Entra ID-based authorization between components to avoid shared secrets and enforce least-privilege access. When you create a simplified reference diagram, treat these identity and authorization controls as baseline implementation concerns rather than explicit flow steps. Keep the diagram focused on orchestration, retry, compensation, and failure handling.

## Related resources

- Data considerations for microservices: Learn why eventual consistency and partial failure are inherent in distributed systems. The Compensating Transaction pattern provides a concrete mechanism to handle those failures when operations span multiple services.
- Transactional Outbox pattern with Azure Cosmos DB: Use this pattern when compensating transactions need to publish events or commands reliably. It helps ensure that state changes and messages are recorded atomically, which prevents message loss.
- Design for self-healing: Use compensating transactions as part of a self-healing approach for your applications.
- Scheduler Agent Supervisor pattern: Use this pattern to implement resilient systems that perform business operations across distributed services and resources. These systems sometimes need compensating transactions to undo work.
- Retry pattern: Use this pattern to handle transient failures and minimize the need for compensating transactions.
- Saga distributed transactions pattern: Use this pattern to manage data consistency across microservices in distributed transactions. Saga uses compensating transactions for failure recovery.
- Pipes and Filters pattern: Use this pattern with the Compensating Transaction pattern as an alternative to distributed transactions when you decompose complex tasks into reusable steps.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Enable multiple concurrent consumers to process messages received on the same messaging channel. With multiple concurrent consumers, a system can process multiple messages at the same time to optimize throughput, improve scalability and availability, and balance the workload.

## Context and problem

A cloud application often handles a large number of requests. Instead of processing each request synchronously, the application can pass requests through a messaging system to a consumer service that handles them asynchronously. This strategy helps prevent request processing from blocking application business logic.

The number of requests can vary significantly over time. A sudden increase in user activity or aggregated requests from multiple tenants can create an unpredictable workload. At peak hours, a system might need to process many hundreds of requests per second. At other times, the number might be small. Also, the work required to handle these requests might vary widely. If you use a single consumer service instance, requests can overwhelm that instance. Or an influx of application messages can overload the messaging system.

To handle this fluctuating workload, the system can run multiple consumer service instances. However, the system must coordinate these consumers to ensure that each message is delivered to only one consumer. The system must also balance the workload across consumers to prevent one instance from becoming a bottleneck.

## Solution

Use a message queue to implement the communication channel between the application and the consumer service instances. The application posts requests as messages to the queue, and consumer service instances receive and process messages from the queue. This approach lets the same pool of consumer service instances handle messages from any instance of the application. The following diagram shows how a message queue distributes work to service instances.

Note

Multiple consumers receive these messages, but the Competing Consumers pattern differs from the Publisher-Subscriber pattern. In the Competing Consumers pattern, one consumer receives each message for processing. In the Publisher-Subscriber pattern, **all** consumers receive **every** message.

This solution has the following benefits:

- It provides a load-leveled system that can handle wide variations in request volume from application instances. The queue functions as a buffer between application instances and consumer service instances. This buffer can minimize the effect on availability and responsiveness for both the application and service instances. For more information, see Queue-based Load Leveling pattern. A message that requires some long-running processing doesn't prevent other consumer service instances from processing other messages concurrently.
- It improves reliability. If a producer communicates directly with a consumer instead of using this pattern and doesn't monitor the consumer, it faces a high probability of losing messages or failing to process them when the consumer fails. In this pattern, the system doesn't send messages to a specific service instance. A failed service instance doesn't block a producer, and any working service instance can process messages.
- It doesn't require complex coordination between consumers or between producer and consumer instances. The message queue ensures that each message is delivered at least once.
- It scales. When you apply autoscaling, the system can dynamically increase or decrease the number of consumer service instances as message volume fluctuates.
- It can improve resiliency if the message queue provides transactional read operations. If a consumer service instance reads and processes a message as part of a transactional operation and fails, this pattern can ensure that the message is returned to the queue so that another consumer service instance can process it. To mitigate the risk of continuous message failures, we recommend that you use dead-letter queues.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- **Message ordering:**The order in which consumer service instances receive messages isn't guaranteed and doesn't necessarily show the order in which the messages were created. Design the system so that it processes messages idempotently. This design helps eliminate processing order dependencies.- Azure Service Bus can implement guaranteed first-in-first-out ordering of messages and other patterns by using message sessions.
- **Service resiliency requirements:**If the system detects and restarts failed service instances, it might need to implement the operations that those service instances perform as idempotent to minimize the effects when it retrieves and processes a single message more than once.
- **Poison message detection:**A malformed message or a task that requires access to resources that aren't available can cause a service instance to fail. The system should prevent these messages from returning to the queue indefinitely and instead capture and store their details elsewhere for analysis if necessary. Service Bus can automatically send messages to a dead-letter queue after the delivery count exceeds the configured- `MaxDeliveryCount`threshold.
- **Result handling:**The service instance that handles a message is fully decoupled from the application logic that generates the message, so they might not be able to communicate directly. If the service instance generates results that must go back to the application logic, store this information in a location that both components can access. To prevent the application logic from retrieving incomplete data, the system must indicate when processing completes. A worker process can pass results back to the application logic through a dedicated message reply queue. The application logic must be able to correlate these results with the original message.
- **Messaging system scaling:**In a large-scale solution, high message volume can overwhelm a single message queue and turn it into a system bottleneck. In this situation, consider partitioning the messaging system to send messages from specific producers to a specific queue, or load balance to distribute messages across multiple message queues.
- **Messaging system reliability:**Use a reliable messaging system to guarantee that messages aren't lost after the application enqueues them. This capability is essential to ensure that all messages are delivered at least once.

## When to use this pattern

Use this pattern when:

- The application workload is divided into tasks that can run asynchronously.
- Tasks are independent and can run in parallel.
- The work volume is highly variable and requires a scalable solution.
- The solution must provide high availability and remain resilient when task processing fails.

This pattern might not be suitable when:

- You can't easily separate the application workload into discrete tasks, or there's a high degree of dependence between tasks.
- Tasks must run synchronously, and the application logic must wait for each task to complete before it continues.
- Tasks must run in a specific sequence.

Note

Some messaging systems support sessions that let a producer group messages together and ensure that the same consumer handles all messages in the group. You can use this mechanism with prioritized messages when supported to enforce message ordering and deliver messages in sequence from a producer to a single consumer.

## Workload design

Evaluate how to use the Competing Consumers pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | This pattern builds redundancy in queue processing by treating consumers as replicas, so an instance failure doesn't prevent other consumers from processing queue messages. - RE:05 Redundancy - Background jobs |
| Cost Optimization focuses on sustaining and improvingyour workload'sreturn on investment. | This pattern can help optimize costs because it scales based on queue depth and can scale down to zero when the queue is empty. It can also optimize costs because you can limit the maximum number of concurrent consumer instances. - CO:05 Rate optimization - CO:07 Component costs |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | This pattern distributes load across consumer nodes to increase utilization, and dynamic scaling based on queue depth minimizes overprovisioning. - PE:05 Scaling and partitioning - PE:07 Code and infrastructure |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

Azure provides Service Bus queues and Azure Functions queue triggers that together directly implement this cloud design pattern. Functions integrates with Service Bus through triggers and bindings. This integration lets you build functions that consume queue messages from publishers. Publishing applications post messages to a queue, and consumers implemented as Functions can retrieve and handle those messages.

For resiliency, a Service Bus queue lets a consumer use PeekLock mode when it retrieves a message from the queue. This mode keeps the message but hides it from other consumers. The Functions runtime receives a message in PeekLock mode. If the function completes successfully, the runtime calls `Complete` on the message. If the function fails, the runtime might call `Abandon` and make the message visible again so that another consumer can retrieve it. If the function runs longer than the PeekLock timeout, the runtime automatically renews the lock as long as the function runs.

Functions automatically scales the number of consumer instances based on queue depth and traffic. This scaling lets the solution handle bursts of work while minimizing cost during low-volume periods. If Functions creates multiple instances, they compete by independently pulling and processing messages. For more information, see Service Bus queues, topics, and subscriptions and Service Bus trigger for Functions.

For more information about how to use the Service Bus client library for .NET to send messages to a Service Bus queue, see the published examples.

## Next steps

- Choose a messaging service in Azure: Learn how different Azure messaging services like Service Bus, Azure Storage queues, Azure Event Hubs, and Azure Event Grid support asynchronous communication patterns, and how to choose the right service and messaging model for your scenario.
- Autoscaling best practices: Learn how to design solutions that scale out consumer instances based on workload, like queue length or message throughput, so that you can handle peak load and control cost during periods of low activity.

## Related resources

- Compute Resource Consolidation pattern: You might be able to consolidate multiple instances of a consumer service into a single process to reduce costs and management overhead. The Compute Resource Consolidation pattern describes the benefits and trade-offs of this approach.
- Queue-based Load Leveling pattern: A message queue can add resiliency to the system. Resiliency lets service instances handle widely varying volumes of requests from application instances. The message queue functions as a buffer that levels the load. The Queue-based Load Leveling pattern describes this scenario in more detail.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Consolidate multiple tasks or operations into one computational unit. This pattern can increase compute resource utilization and reduce the costs and management overhead associated with compute processing in cloud-hosted applications.

## Context and problem

A cloud application often implements different types of operations. Initially, you can organize these operations into separate computational units that are hosted and deployed individually. For example, you can deploy separate Azure App Service web apps or separate virtual machines. This strategy can help simplify the logical design of the solution, but if you deploy a large number of computational units as part of the same application, this deployment can increase runtime hosting costs and make system management more complex.

The following figure shows the simplified structure of a cloud-hosted solution that uses multiple computational units. Each computational unit runs in its own virtual environment. The solution implements each function as a separate task that runs in its own computational unit.

Each computational unit consumes chargeable resources even when it's idle or lightly used. This approach isn't always the most cost-effective solution.

## Solution

To help reduce costs, increase utilization, improve communication speed, and reduce management, you can consolidate multiple tasks or operations into a single computational unit.

You can group tasks according to criteria based on the features of the environment and the costs associated with these features. It's common to look for tasks that have similar scalability, lifetime, and processing requirements. To scale the tasks as a unit, you can group them together. Many cloud environments offer elasticity so that you can start and stop extra instances of a computational unit depending on the workload. For example, Azure provides autoscaling that you can apply to App Service and Azure Virtual Machine Scale Sets.

You can also use scalability to determine which operations you shouldn't group together. Consider the following example tasks:

- Task 1 polls a queue of infrequent, time-insensitive messages.
- Task 2 handles high-volume bursts of network traffic.

Task 2 requires elasticity to start and stop a large number of computational units. If you apply the same scaling behavior to Task 1, more tasks listen for infrequent messages on the same queue, which is a waste of resources.

In many cloud environments, you can specify the resources available to a computational unit, such as the number of CPU cores, memory, and disk space. If you specify more resources, the solution usually becomes more expensive. To save money, an expensive computational unit should stay busy and avoid extended periods of inactivity.

If tasks require high CPU power in short bursts, you might consolidate these tasks into a single computational unit that provides the necessary power. However, balance this need to keep expensive resources busy against the contention that could occur if they're stressed. For example, long-running, compute-intensive tasks shouldn't share the same computational unit.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- **Scalability and elasticity:**Many cloud solutions implement computational unit scalability and elasticity by starting and stopping instances of units. Don't group tasks that have conflicting scalability requirements in the same computational unit.
- **Lifetime:**The cloud infrastructure periodically recycles the virtual environment that hosts a computational unit. If there are many long-running tasks inside a computational unit, you might need to prevent the unit from being recycled until these tasks finish. Alternatively, use a checkpointing approach so that the tasks can stop cleanly and continue from where they were interrupted when the computational unit restarts.
- **Release cadence:**If the implementation or configuration of a task changes frequently, you might need to stop the computational unit that hosts the updated code, reconfigure and redeploy the unit, and then restart it. This process also requires you to stop, redeploy, and restart all other tasks within the same computational unit.
- **Security:**Tasks in the same computational unit might share the same security context and might be able to access the same resources. This setup requires a high level of trust between tasks and confidence that one task can't corrupt or adversely affect another. Additionally, if you increase the number of tasks that run in a computational unit, the attack surface of the unit increases. Each task is only as secure as the task with the most vulnerabilities.
- **Fault tolerance:**If one task in a computational unit fails or behaves abnormally, it can affect the other tasks in the same unit. For example, if one task fails to start correctly it can cause the entire startup logic for the computational unit to fail, and it can prevent other tasks in the same unit from running.
- **Contention:**Avoid introducing contention between tasks that compete for resources in the same computational unit. Tasks that share the same computational unit should exhibit different resource utilization characteristics. For example, two compute-intensive tasks shouldn't reside in the same computational unit and neither should two tasks that consume large amounts of memory. However, you can combine a compute-intensive task with a task that requires a large amount of memory.- Note - Consider consolidating compute resources only for systems that are in production long enough so that operators and developers can monitor the system and create a - *heat map*that identifies how each task uses resources. This map helps determine which tasks are good candidates for sharing compute resources.
- **Complexity:**Multiple tasks in a single computational unit adds complexity to the code in the unit, which might make it more difficult to test, debug, and maintain.
- **Stable logical architecture:**Design and implement the code in each task so that it shouldn't need to change, even if the physical environment that the task runs in changes.
- **Other strategies:**Consolidation of compute resources is only one way to help reduce the costs associated with running multiple tasks concurrently. It requires careful planning and monitoring to ensure that it remains an effective approach. Other strategies might be more appropriate, depending on the nature of the work and where the users of the tasks are located.

## When to use this pattern

Use this pattern when:

- Tasks aren't cost effective if they run in their own computational units.
- A task is often idle.
- It would be expensive to run a task in a dedicated unit.

This pattern might not be suitable when:

- Tasks perform critical fault-tolerant operations.
- Tasks process highly sensitive or private data and require their own security context.
- Tasks need to run in their own isolated environment in a separate computational unit.

## Workload design

Evaluate how to use the compute resource consolidation pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Cost Optimization focuses on sustaining and improvingyour workload'sreturn on investment. | This pattern maximizes the utilization of compute resources by avoiding unused provisioned capacity via aggregation of components or even whole workloads on a pooled infrastructure. - CO:14 Consolidation |
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | Consolidation can lead to a more homogeneous compute platform, which can simplify management and observability, reduce disparate approaches to operational tasks, and reduce the amount of tooling required. - OE:07 Monitoring system - OE:10 Automation design |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | Consolidation maximizes compute resource utilization by using spare node capacity and reduces overprovisioning. These infrastructures often use large, vertically scaled compute instances in the resource pool. - PE:02 Capacity planning - PE:03 Select services |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

You can deploy this pattern in different ways, depending on your compute service. For example:

- **App Service**and- **Azure Functions:**Deploy shared App Service plans, which represent the hosting server infrastructure. You can set up one or more apps to run on the same compute resources or in the same App Service plan.
- **Azure Container Apps:**Deploy Container Apps to the same shared environments, especially if you need to manage related services or you need to deploy different applications to the same virtual network.
- **Azure Kubernetes Service (AKS):**AKS is a container-based hosting infrastructure in which you can set up multiple applications or application components to run colocated on the same compute resources (nodes). You can group the compute resources by computational requirements, such as CPU or memory needs (node pools).
- **Virtual machines:**Deploy a single set of virtual machines for all tenants to use so that management costs are shared across the tenants. Virtual Machine Scale Sets supports shared resource management, load balancing, and horizontal scaling of virtual machines.
- **Consumption-based compute (serverless):**Use fully managed, pay‑per‑execution compute models that can scale to zero, such as the Functions Consumption plan and Container Apps. These services run multiple independent workloads on a shared global pool of compute resources. These platforms automatically allocate and release compute so that multiple applications can benefit from elastic scaling and cost efficiency without a dedicated infrastructure.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Segregate the read and write operations for a data store into separate data models. This approach allows you to optimize each model independently and can improve the performance, scalability, and security of an application.

## Context and problem

In a traditional architecture, a single data model is often used for both read and write operations. This approach is straightforward and is suited for basic create, read, update, and delete (CRUD) operations.

As applications grow, it can become increasingly difficult to optimize read and write operations on a single data model. Read and write operations often have different performance and scaling requirements. A traditional CRUD architecture doesn't take this asymmetry into account, which can result in the following challenges:

- **Data mismatch:**The read and write representations of data often differ. Some fields that are required during updates might be unnecessary during read operations.
- **Lock contention:**Parallel operations on the same data set can cause lock contention.
- **Performance problems:**The traditional approach can have a negative effect on performance because of load on the data store and data access layer, and the complexity of queries required to retrieve information.
- **Security challenges:**It can be difficult to manage security when entities are subject to read and write operations. This overlap can expose data in unintended contexts.

Combining these responsibilities can result in an overly complicated model.

## Solution

Use the CQRS pattern to separate write operations, or *commands*, from read operations, or *queries*. Commands update data. Queries retrieve data. The CQRS pattern is useful in scenarios that require a clear separation between commands and reads.

- **Understand commands.**Commands should represent specific business tasks instead of low-level data updates. For example, in a hotel-booking app, use the command "Book hotel room" instead of "Set ReservationStatus to Reserved." This approach better captures the intent of the user and aligns commands with business processes. To help ensure that commands are successful, you might need to refine the user interaction flow and server-side logic and consider asynchronous processing.- Area of refinement - Recommendation - Client-side validation - Validate specific conditions before you send the command to prevent obvious failures. For example, if no rooms are available, disable the "Book" button and provide a clear, user-friendly message in the UI that explains why booking isn’t possible. This setup reduces unnecessary server requests and provides immediate feedback to users, which enhances their experience. - Server-side logic - Enhance the business logic to handle edge cases and failures gracefully. For example, to address race conditions such as multiple users attempting to book the last available room, consider adding users to a waiting list or suggesting alternatives. - Asynchronous processing - Process commands asynchronously by placing them in a queue, instead of handling them synchronously.
- **Understand queries.**Queries never alter data. Instead, they return data transfer objects (DTOs) that present the required data in a convenient format, without any domain logic. This distinct separation of responsibilities simplifies the design and implementation of the system.

### Separate read models and write models

Separating the read model from the write model simplifies system design and implementation by addressing specific concerns for data writes and data reads. This separation improves clarity, scalability, and performance but introduces trade-offs. For example, scaffolding tools like object-relational mapping (O/RM) frameworks can't automatically generate CQRS code from a database schema, so you need custom logic to bridge the gap.

The following sections describe two primary approaches to implement read model and write model separation in CQRS. Each approach has unique benefits and challenges, such as synchronization and consistency management.

#### Separate models in a single data store

This approach represents the foundational level of CQRS, where both the read and write models share a single underlying database but maintain distinct logic for their operations. A basic CQRS architecture allows you to delineate the write model from the read model while relying on a shared data store.

This approach improves clarity, performance, and scalability by defining distinct models for handling read and write concerns.

- **A write model**is designed to handle commands that update or persist data. It includes validation and domain logic, and helps ensure data consistency by optimizing for transactional integrity and business processes.
- **A read model**is designed to serve queries for retrieving data. It focuses on generating DTOs or projections that are optimized for the presentation layer. It enhances query performance and responsiveness by avoiding domain logic.

#### Separate models in different data stores

A more advanced CQRS implementation uses distinct data stores for the read and write models. Separation of the read and write data stores allows you to scale each model to match the load. It also enables you to use a different storage technology for each data store. You can use a document database for the read data store and a relational database for the write data store.

When you use separate data stores, you must ensure that both remain synchronized. A common pattern is to have the write model publish events when it updates the database, which the read model uses to refresh its data. For more information about how to use events, see Event-driven architecture style. Because you usually can't enlist message brokers and databases into a single distributed transaction, challenges in consistency can occur when you update the database and publishing events. For more information, see Idempotent message processing.

The read data store can use its own data schema that's optimized for queries. For example, it can store a materialized view of the data to avoid complex joins or O/RM mappings. The read data store can be a read-only replica of the write store or have a different structure. Deploying multiple read-only replicas can improve performance by reducing latency and increasing availability, especially in distributed scenarios.

### Benefits of CQRS

- **Independent scaling.**CQRS enables the read models and write models to scale independently. This approach can help minimize lock contention and improve system performance under load.
- **Optimized data schemas.**Read operations can use a schema that's optimized for queries. Write operations use a schema that's optimized for updates.
- **Security.**By separating reads and writes, you can ensure that only the appropriate domain entities or operations have permission to perform write actions on the data.
- **Separation of concerns.**Separating the read and write responsibilities results in cleaner, more maintainable models. The write side typically handles complex business logic. The read side can remain simple and focused on query efficiency.
- **Simpler queries.**When you store a materialized view in the read database, the application can avoid complex joins when it queries.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- **Increased complexity.**The core concept of CQRS is straightforward, but it can introduce significant complexity into the application design, specifically when combined with the Event Sourcing pattern.
- **Messaging challenges.**Messaging isn't a requirement for CQRS, but you often use it to process commands and publish update events. When messaging is included, the system must account for potential problems such as message failures, duplicates, and retries. For more information about strategies to handle commands that have varying priorities, see Priority queues.
- **Eventual consistency.**When the read databases and write databases are separated, the read data might not show the most recent changes immediately. This delay results in stale data. Ensuring that the read model store stays up-to-date with changes in the write model store can be challenging. Also, detecting and handling scenarios where a user acts on stale data requires careful consideration.

## When to use this pattern

Use this pattern when:

- **You work in collaborative environments.**In environments where multiple users access and modify the same data simultaneously, CQRS helps reduce merge conflicts. Commands can include enough granularity to prevent conflicts, and the system can resolve any conflicts that occur within the command logic.
- **You have task-based user interfaces.**Applications that guide users through complex processes as a series of steps or with complex domain models benefit from CQRS.- The write model has a full command-processing stack with business logic, input validation, and business validation. The write model might treat a set of associated objects as a single unit for data changes, which is known as an - *aggregate*in domain-driven design terminology. The write model might also help ensure that these objects are always in a consistent state.
- The read model has no business logic or validation stack. It returns a DTO for use in a view model. The read model is eventually consistent with the write model.

- **You need performance tuning.**Systems where the performance of data reads must be fine-tuned separately from performance of data writes benefit from CQRS. This pattern is especially beneficial when the number of reads is greater than the number of writes. The read model scales horizontally to handle large query volumes. The write model runs on fewer instances to minimize merge conflicts and maintain consistency.
- **You have separation of development concerns.**CQRS allows teams to work independently. One team implements the complex business logic in the write model, and another team develops the read model and user interface components.
- **You have evolving systems.**CQRS supports systems that evolve over time. It accommodates new model versions, frequent changes to business rules, or other modifications without affecting existing functionality.
- **You need system integration:**Systems that integrate with other subsystems, especially systems that use the Event Sourcing pattern, remain available even if a subsystem temporarily fails. CQRS isolates failures, which prevents a single component from affecting the entire system.

This pattern might not be suitable when:

- The domain or the business rules are simple.
- A simple CRUD-style user interface and data access operations are sufficient.

## Workload design

Evaluate how to use the CQRS pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of the Performance Efficiency pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Performance Efficiency helps your workload efficiently meet demands through optimizations in scaling, data, and code. | The separation of read operations and write operations in high read-to-write workloads enables targeted performance and scaling optimizations for each operation's specific purpose. - PE:05 Scaling and partitioning - PE:08 Data performance |

Consider any trade-offs against the goals of the other pillars that this pattern might introduce.

## Combine the Event Sourcing and CQRS patterns

Some implementations of CQRS incorporate the Event Sourcing pattern. This pattern stores the system's state as a chronological series of events. Each event captures the changes made to the data at a specific time. To determine the current state, the system replays these events in order. In this setup:

- The event store is the - *write model*and the single source of truth.
- The - *read model*generates materialized views from these events, typically in a highly denormalized form. These views optimize data retrieval by tailoring structures to query and display requirements.

### Benefits of combining the Event Sourcing and CQRS patterns

The same events that update the write model can serve as inputs to the read model. The read model can then build a real-time snapshot of the current state. These snapshots optimize queries by providing efficient and precomputed views of the data.

Instead of directly storing the current state, the system uses a stream of events as the write store. This approach reduces update conflicts on aggregates and enhances performance and scalability. The system can process these events asynchronously to build or update materialized views for the read data store.

Because the event store acts as the single source of truth, you can easily regenerate materialized views or adapt to changes in the read model by replaying historical events. Basically, materialized views function as a durable, read-only cache that's optimized for fast and efficient queries.

### Considerations for how to combine the Event Sourcing and CQRS patterns

Before you combine the CQRS pattern with the Event Sourcing pattern, evaluate the following considerations:

- **Eventual consistency:**Because the write and read data stores are separate, updates to the read data store might lag behind event generation. This delay results in eventual consistency.
- **Increased complexity:**Combining the CQRS pattern with the Event Sourcing pattern requires a different design approach, which can make a successful implementation more challenging. You must write code to generate, process, and handle events, and assemble or update views for the read model. However, the Event Sourcing pattern simplifies domain modeling and allows you to rebuild or create new views easily by preserving the history and intent of all data changes.
- **Performance of view generation:**Generating materialized views for the read model can consume significant time and resources. The same applies to projecting data by replaying and processing events for specific entities or collections. Complexity increases when calculations involve analyzing or summing values over long periods because all related events must be examined. Implement snapshots of the data at regular intervals. For example, store the current state of an entity or periodic snapshots of aggregated totals, which is the number of times a specific action occurs. Snapshots reduce the need to process the full event history repeatedly, which improves performance.

## Example

The following code shows extracts from an example of a CQRS implementation that uses different definitions for the read models and the write models. The model interfaces don't dictate features of the underlying data stores, and they can evolve and be fine-tuned independently because these interfaces are separate.

The following code shows the read model definition.

```
// Query interface
namespace ReadModel
{
 public interface ProductsDao
 {
 ProductDisplay FindById(int productId);
 ICollection<ProductDisplay> FindByName(string name);
 ICollection<ProductInventory> FindOutOfStockProducts();
 ICollection<ProductDisplay> FindRelatedProducts(int productId);
 }
 public class ProductDisplay
 {
 public int Id { get; set; }
 public string Name { get; set; }
 public string Description { get; set; }
 public decimal UnitPrice { get; set; }
 public bool IsOutOfStock { get; set; }
 public double UserRating { get; set; }
 }
 public class ProductInventory
 {
 public int Id { get; set; }
 public string Name { get; set; }
 public int CurrentStock { get; set; }
 }
}
```
The system allows users to rate products. The application code does this by using the `RateProduct` command shown in the following code.

```
public interface ICommand
{
 Guid Id { get; }
}
public class RateProduct : ICommand
{
 public RateProduct()
 {
 this.Id = Guid.NewGuid();
 }
 public Guid Id { get; set; }
 public int ProductId { get; set; }
 public int Rating { get; set; }
 public int UserId {get; set; }
}
```
The system uses the `ProductsCommandHandler` class to handle commands that the application sends. Clients typically send commands to the domain through a messaging system such as a queue. The command handler accepts these commands and invokes methods of the domain interface. The granularity of each command is designed to reduce the chance of conflicting requests. The following code shows an outline of the `ProductsCommandHandler` class.

```
public class ProductsCommandHandler :
 ICommandHandler<AddNewProduct>,
 ICommandHandler<RateProduct>,
 ICommandHandler<AddToInventory>,
 ICommandHandler<ConfirmItemShipped>,
 ICommandHandler<UpdateStockFromInventoryRecount>
{
 private readonly IRepository<Product> repository;
 public ProductsCommandHandler (IRepository<Product> repository)
 {
 this.repository = repository;
 }
 void Handle (AddNewProduct command)
 {
 ...
 }
 void Handle (RateProduct command)
 {
 var product = repository.Find(command.ProductId);
 if (product != null)
 {
 product.RateProduct(command.UserId, command.Rating);
 repository.Save(product);
 }
 }
 void Handle (AddToInventory command)
 {
 ...
 }
 void Handle (ConfirmItemsShipped command)
 {
 ...
 }
 void Handle (UpdateStockFromInventoryRecount command)
 {
 ...
 }
}
```
## Next step

The following information might be relevant when you implement this pattern:

- Data partitioning guidance describes best practices for how to divide data into partitions that you can manage and access separately to improve scalability, reduce contention, and optimize performance.

## Related resources

- Event Sourcing pattern. This pattern describes how to simplify tasks in complex domains and improve performance, scalability, and responsiveness. It also explains how to provide consistency for transactional data while maintaining full audit trails and history that can enable compensating actions.
- Materialized View pattern. This pattern creates prepopulated views, known as - *materialized views*, for efficient querying and data extraction from one or more data stores. The read model of a CQRS implementation can contain materialized views of the write model data, or the read model can be used to generate materialized views.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Deploy multiple independent copies of application components, including data stores, as a single group of resources. Each copy is called a *stamp*, or sometimes a *service unit*, *scale unit*, or *cell*. In a multitenant environment, each stamp serves a predefined number of tenants. Deploy more stamps to scale the solution almost linearly, serve an increasing number of tenants, deploy instances across multiple regions, and separate your customer data.

Note

For more information, see Architect multitenant solutions on Azure.

## Context and problem

When you host an application in the cloud, consider the performance and reliability of your application. If you host a single instance of your solution, the following limitations might apply:

- **Scale limits:**A single instance of your application might reach natural scaling limits. For example, the services that you use might limit the number of inbound connections, host names, Transmission Control Protocol (TCP) sockets, or other resources.
- **Nonlinear scaling or cost:**Some of your solution's components might not scale linearly with the number of requests or the amount of data. Instead, performance can drop or cost can spike after you meet a threshold. For example, you might find that adding more capacity to a database, or scaling up, becomes prohibitively expensive and that scaling out is more cost effective.
- **Separation of customers:**You might need to isolate one customer's data from another customer's data. You might also have customers that consume more system resources than others. You can group them on different sets of infrastructure.
- **Single-tenant and multitenant instances:**Some large customers might need their own independent instances of your solution. Smaller customers can share a multitenant deployment.
- **Complex deployment requirements:**You might need to deploy updates to your service in a controlled manner and deploy to different subsets of your customer base at different times.
- **Update frequency:**Some customers tolerate frequent updates, while risk-averse customers want infrequent updates to the system that serves their requests. You can deploy these customers to isolated environments.
- **Geographical or geopolitical restrictions:**To achieve low latency or comply with data sovereignty requirements, you might deploy some customers to specific regions.

These limitations often apply to software development companies that build software as a service (SaaS), which they typically design as multitenant. The same limitations can also apply to other scenarios.

## Solution

To avoid these problems, consider grouping resources into *scale units* and provisioning multiple copies of your *stamps*. Each scale unit hosts and serves a subset of your tenants. Stamps run independently of each other, and you can deploy and update them independently. A single geographic region might contain one stamp or multiple stamps that scale out horizontally within the region. Each stamp serves a subset of your customers.

Deployment stamps can apply whether your solution uses infrastructure as a service (IaaS) or platform as a service (PaaS) components, or a combination of both. IaaS workloads typically require more intervention to scale, so this pattern can help IaaS-heavy workloads scale out.

You can use stamps to implement deployment rings. If different customers want service updates at different frequencies, group them onto different stamps and deploy updates to each stamp at a different cadence.

Stamps run independently, so they implicitly *shard* your data. A single stamp can also use further sharding internally to scale and remain elastic.

Deploying identical copies of the same components is complex, so good DevOps practices are critical. Describe your infrastructure as code so that the deployment of each stamp is predictable and repeatable.

Deployment stamps relate to but differ from geodes. In a deployment stamp architecture, each independent instance of your system serves a subset of your customers and users. In a geode architecture, every instance can serve requests from any user, but this approach is typically more complex to design and build. You can also combine the two patterns within one solution. The traffic routing approach described later in this article is an example of such a hybrid scenario.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- **Deployment process:**When you deploy multiple stamps, automate and fully repeat your deployment processes. Use Bicep or Terraform modules to declaratively define your stamps and keep the definitions consistent.
- **Cross-stamp operations:**When you deploy your solution independently across multiple stamps, it can be hard to determine how many customers you have across all your stamps. You might need to query each stamp and aggregate the results. Alternatively, you can have all stamps publish data into a centralized data warehouse for consolidated reporting.
- **Scale-out policies:**Stamps have a finite capacity, which you can define by using a proxy metric, such as the number of tenants that you can deploy to the stamp. Monitor the available and used capacity for each stamp, and proactively deploy more stamps to direct new tenants to them.
- **Minimum number of stamps:**When you use the Deployment Stamps pattern, deploy at least two stamps of your solution. If you deploy only a single stamp, you can easily hard-code assumptions into your code or configuration that don't apply when you scale out.
- **Cost:**The Deployment Stamps pattern deploys multiple copies of your infrastructure components, which substantially increases the cost of operating your solution.
- **Moving between stamps:**Each stamp runs independently, so moving tenants between stamps can be difficult. Your application needs custom logic to transmit a customer's information to a different stamp and then remove the tenant's information from the original stamp. This process might require a backplane to communicate between stamps, which further increases the complexity of your solution.
- **Traffic routing:**As described previously in this article, routing traffic to the correct stamp for a given request can require an extra component that resolves tenants to stamps. This component might also need to be highly available.
- **Observability across stamps:**As the number of stamps increases, it becomes harder to understand overall health and detect incidents quickly. Use Azure Monitor to collect and correlate metrics, logs, traces, and alerts across all stamps. Use this data to identify unhealthy stamps and diagnose problems.
- **Regional failure impact:**Stamps run independently, but they aren't inherently redundant across regions. If a region that hosts one or more stamps becomes unavailable, the tenants on those stamps lose access until the region recovers or you migrate the tenants to stamps in another region. To plan for this scenario, document your recovery procedures, set tenant expectations, and consider whether critical tenants need geo-redundant stamp placement.
- **Shared components:**You might have components that you can share across stamps. For example, if you have a shared single-page app for all tenants, deploy it to one region and use Azure Front Door edge caching to replicate it globally.
- **Governance and configuration drift:**As the number of stamps increases, it becomes harder to keep security policies, role-based access control (RBAC) assignments, network controls, observability settings, and service configurations consistent. Use Azure Policy to treat governance as code and continuously validate each stamp for drift to prevent inconsistent behavior and compliance gaps.

## When to use this pattern

Use this pattern when:

- Your solution has natural limits on scalability. For example, if some components can't or shouldn't scale beyond a certain number of customers or requests, use stamps to scale out.
- You need to separate certain tenants from others. If security concerns prevent you from deploying some customers into a multitenant stamp, deploy them onto their own isolated stamp.
- You need to host some tenants on different versions of your solution at the same time.
- You build multiregion applications that need to direct each tenant's data and traffic to a specific region.
- You want to achieve resiliency during outages. Stamps run independently, so if an outage affects a single stamp, tenants on other stamps remain unaffected. This isolation contains the - *blast radius*of an incident or outage.

This pattern might not be suitable when:

- Your solution is simple and doesn't need to scale to a high degree.
- You can scale your system out or up within a single instance, such as by increasing the size of the application layer or by increasing the reserved capacity for databases and the storage tier.
- You need to replicate data across all deployed instances. Consider the Geode pattern for this scenario.
- You only need to scale some components and not others. For example, consider whether you can scale your solution by sharding the data store instead of deploying a new copy of all the solution components.
- Your solution consists solely of static content, such as a front-end JavaScript application. Deliver this content by a Content Delivery Network.

## Workload design

Evaluate how to use the Deployment Stamps pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | Stamps operate independently, so a failure in one stamp is isolated and doesn't affect tenants on other stamps. Deploying multiple stamps across regions also provides a foundation for redundancy and recovery planning, which reduces the blast radius of regional outages. - RE:05 Redundancy - RE:07 Self-preservation |
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | This pattern supports immutable infrastructure goals, advanced deployment models, and can facilitate safe deployment practices. - OE:05 Infrastructure as code - OE:11 Safe deployment practices |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | This pattern often aligns to the defined scale units in your workload. When you need more capacity than a single scale unit provides, you deploy another stamp to scale out. - PE:05 Scaling and partitioning |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

The following example architecture uses Azure Front Door, Azure API Management, and Azure Cosmos DB to route traffic globally to a series of region-specific stamps.

Suppose a user resides in New York. Stamp 3, in the East US region, stores their data.

If the user travels to California and accesses the system, the system routes their connection through the West US 2 region because that region is closest to them when they make the request. However, stamp 3 must ultimately serve the request because it stores their data. The traffic routing system routes the request to the correct stamp.

### Deployment

Describe your infrastructure as code, by using Bicep or Terraform. This approach ensures that the deployment of each stamp is predictable and repeatable. It also reduces the likelihood of human errors such as accidental mismatches in configuration between stamps.

You can deploy updates automatically to all stamps in parallel. Technologies like Bicep can coordinate the deployment of your infrastructure and applications. Alternatively, you might decide to gradually roll out updates to some stamps first, and then progressively to other stamps. Consider using a release management tool like Azure Pipelines or GitHub Actions to orchestrate deployments to each stamp.

Carefully consider the topology of the Azure subscriptions and resource groups for your deployments:

- Typically, a subscription contains all resources for a single solution, so consider using a single subscription for all stamps. However, some Azure services impose subscription-wide quotas. If you use this pattern to allow for a high degree of scale-out, you might need to deploy stamps across different subscriptions.
- Resource groups generally contain components that share the same life cycle. If you plan to deploy updates to all stamps at the same time, you can use a single resource group that contains all components for all stamps. Use resource naming conventions and tags to identify the components that belong to each stamp. Alternatively, if you plan to deploy updates to each stamp independently, you can deploy each stamp into its own resource group.

### Capacity planning

Use load and performance testing to determine the approximate load that a given stamp can accommodate. Load metrics might be based on the number of customers or tenants that a single stamp can accommodate, or on metrics that the services in the stamp emit. Instrument each stamp so that you can measure when it approaches its capacity, and make sure that you can deploy new stamps quickly to respond to demand.

### Traffic routing

The Deployment Stamps pattern works well when you address each stamp independently. For example, if Contoso deploys the same API application across multiple stamps, Contoso might use Domain Name System (DNS) to route traffic to the relevant stamp:

- `unit1.aus.myapi.contoso.com`routes traffic to stamp- `unit1`within an Australian region.
- `unit2.aus.myapi.contoso.com`routes traffic to stamp- `unit2`within an Australian region.
- `unit1.eu.myapi.contoso.com`routes traffic to stamp- `unit1`within a European region.

In Azure, you can host these records in Azure DNS and use a consistent subdomain convention for each region and stamp. This approach maintains predictable routing and operations.

Clients are responsible for connecting to the correct stamp.

If your solution requires a single ingress point for all traffic, you can use a traffic routing service to resolve the stamp for a given request, customer, or tenant. The traffic routing service either directs the client to the relevant URL for the stamp (for example, by returning an HTTP 302 response status code), or it acts as a reverse proxy and forwards the traffic to the relevant stamp without the client being aware.

A centralized traffic routing service can be a complex component to design, especially when a solution runs across multiple regions. Consider deploying the traffic routing service into multiple regions, potentially including every region that hosts stamps, and sync the data store that maps tenants to stamps. The traffic routing component might itself be an instance of the Geode pattern.

For example, you can deploy API Management to act as the traffic routing service. API Management determines the appropriate stamp for a request by looking up data in an Azure Cosmos DB collection that stores the mapping between tenants and stamps. API Management then dynamically sets the back-end URL to the relevant stamp's API service.

To geo-distribute requests and provide geo-redundancy for the traffic routing service, deploy API Management across multiple regions and use Azure Front Door to direct traffic to the closest API Management gateway. In this topology, Azure Front Door uses origin groups, health probes, and an appropriate routing method to route requests away from unhealthy API Management regional gateways. API Management then routes to the appropriate stamp by using the tenant-to-stamp mapping and its back-end configuration (or back-end pools), including failover rules between stamp endpoints as needed. If your application isn't exposed over HTTP or HTTPS, you can use a cross-region Azure load balancer to distribute incoming calls to regional Azure load balancers. Use the global distribution feature of Azure Cosmos DB to keep the mapping information updated across each region.

If your solution includes a traffic routing service, consider whether it acts as a gateway and can perform gateway offloading for the other services, such as token validation, throttling, and authorization.

## Next steps

- Azure Front Door
- Integrate Bicep with Azure Pipelines
- Integrate JSON ARM templates with Azure Pipelines

## Contributors

*Microsoft maintains this article. The following contributors wrote this article.*

Principal author:

- John Downs | Principal Software Engineer, Azure Patterns & Practices

Other contributors:

- Federico Arambarri | Senior Software Developer, Clarius Consulting
- Daniel Larsen | Principal Customer Engineer, FastTrack for Azure
- Angel Lopez | Senior Software Engineer, Azure Patterns and Practices
- Paolo Salvatori | Principal Customer Engineer, FastTrack for Azure
- Arsen Vladimirskiy | Principal Customer Engineer, FastTrack for Azure

*To see nonpublic LinkedIn profiles, sign in to LinkedIn.*

## Related resources

- You can use sharding as another simpler approach to scale out your data tier. Stamps implicitly shard their data, but sharding doesn't require a deployment stamp. For more information, see Sharding pattern.
- If your solution deploys a traffic routing service, you can combine the Gateway Routing and Gateway Offloading patterns to make the best use of this component.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Instead of storing only the current state of the data in a relational database, store the full series of actions taken on an object in an append-only store. The store acts as the system of record that you can use to materialize the domain objects. This approach can improve auditability and write performance in complex systems.

Important

Event sourcing is a complex pattern that introduces significant trade-offs. It changes how you store data, handle concurrency, evolve schemas, and query state. It's costly to migrate to or from an event sourcing solution, and after you adopt the pattern, it constrains future design decisions in the parts of the system that use it. Adopt event sourcing when its benefits, like auditability and historical reconstruction, justify the pattern's complexity. For most systems and most parts of a system, traditional data management is sufficient.

## Context and problem

Most applications work with data. The application typically stores the latest state of the data in a relational database and inserts or updates data as needed. For example, in the traditional create, read, update, and delete (CRUD) model, an application reads data from the store, modifies it, and updates the current state of the data with the new values, typically by using transactions that lock the data.

The CRUD approach is straightforward and fast for most scenarios. However, in high-load systems, this approach presents challenges:

- **Write contention:**Because updates require read-modify-write cycles with row-level locking, concurrent writes to the same entity degrade performance and become a bottleneck under load.
- **Auditability:**CRUD systems only store the latest state of the data. If you don't implement an auditing mechanism that records the details of each operation in a separate log, you lose data history.

## Solution

The Event Sourcing pattern defines an approach to handling operations on data that a sequence of events drive. Each event is recorded in an append-only store. Application code raises events that describe each action taken on the object. It typically sends events to a queue in which a separate process, an event handler, listens to the queue and persists the events in an event store. Each event represents a logical change to the object, such as `AddedItemToOrder` or `OrderCanceled`.

The events persist in an event store that serves as the system of record, or the authoritative data source, about the current state of the data. Extra event handlers can listen for specific events and take action as needed. For example, consumers might initiate tasks that apply operations in the events to other systems or take other associated actions required to finish the operation. The application code that generates the events is decoupled from the systems that subscribe to the events.

Each entity in an event-sourced system has its own eventstream, which is the ordered sequence of events that records every change to that entity. At any point, applications can read the history of events. Applications derive the current state of an entity by replaying all the events in its stream. This process is known as *rehydration*. It can occur on demand when the application handles a request.

Applications typically implement materialized views because it's costly to read and replay events. Materialized views are read-only projections of the event store that are optimized for querying. For example, a system can maintain a materialized view of all customer orders that it uses to populate the UI. When the application adds new orders, adds or removes items in the order, or adds shipping information, the application raises events and a handler updates the materialized view.

The following diagram shows an overview of this pattern combined with the Command Query Responsibility Segregation (CQRS) pattern. The presentation layer reads from a separate read-only store and writes commands to command handlers. The command handlers retrieve the entity's eventstream from the event store, run business logic, and push new events to a queue. Event handlers consume events from the queue and write events to the event store, update the read-only store, or integrate with external systems.

*Download a Visio file of this architecture.*

### Workflow

The following workflow corresponds to the previous diagram:

- The presentation layer calls an object that reads from a read-only store. It uses the returned data to populate the UI.
- The presentation layer calls command handlers to perform actions like - *create a cart*or- *add an item to the cart*.
- The command handler loads the entity by retrieving its eventstream from the event store. For example, it might retrieve all cart events. It replays those events against the entity to reconstruct its current state before any new action occurs.
- The business logic runs and events are raised. In most implementations, the events are pushed to a queue or topic to decouple the event producers and event consumers.
- Event handlers listen for specific events and take the appropriate action for that handler. In this example, the event handlers take the following actions: - Write the events to the event store
- Update a read-only store optimized for queries
- Integrate with external systems

### Pattern advantages

The Event Sourcing pattern provides the following advantages:

- Events are immutable, and you can store them by using an append-only operation. The UI, workflow, or process that initiates an event can continue, and tasks that handle the events can run in the background. Write throughput improves, especially for the presentation layer, because append-only writes avoid the row-level lock contention that update-in-place systems create.
- Events are simple objects that describe an action that occurs along with any associated data required to describe the action that the event represents. Events don't directly update a data store. Event handlers pick up and process recorded events when a handler is available and the system can handle the load. Use events to help simplify implementation and management.
- Events typically have meaning for a domain expert, whereas object-relational impedance mismatch can make complex database tables hard to understand. Tables are artificial constructs that represent the current state of the system, not the events that occur.
- Event sourcing can help prevent concurrent updates from causing conflicts because it avoids the requirement to directly update objects in the data store. Command handlers rehydrate an entity from its eventstream to enforce business rules before they append new events, so two handlers that load the same entity simultaneously can act on the same state. - For example, each handler sees five remaining seats, and both handlers can accept a reservation. Event stores address this scenario by using optimistic concurrency control and reject an append if the stream changed since it was read. Upon rejection, the handler reloads the entity, reevaluates, and retries.
- Append-only event storage provides an audit trail that applications can use to monitor actions taken against a data store. It can regenerate the current state as materialized views or projections by replaying the events at any time, and it can help test and debug the system. - The requirement to use compensating events to cancel changes can provide a history of reversed changes. If the model stores only the current state, this history doesn't exist. You can also use the list of events to analyze application performance, detect user behavior trends, and obtain other useful business information.
- The command handlers raise events, and tasks perform operations in response to those events. This decoupling of the tasks from the events provides flexibility and extensibility. Tasks know about the type of event and the event data, but not about the operation that triggers the event. - Multiple tasks can handle each event, so they can easily integrate with other services and systems that only listen for new events that the event store raises. But the event sourcing events are typically low level, and it might be necessary to generate specific integration events instead.

Tip

Event sourcing is commonly combined with the CQRS pattern by performing the data management tasks in response to the events and by materializing views from the stored events. Use this combination to independently scale reads and writes because append-only event ingestion and query-optimized projections operate separately.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- **Event design:**Design events to capture the business intent behind each change in addition to the resulting state. For example, in the seat-reservation system, an event that records- *two seats were reserved*is more valuable than an event that records- *remaining seats changed to 42*. The first event tells you what happened. The second event only tells you the resulting state. State-focused events reduce the event store to a change log that has no business meaning. Intent-focused events provide more detailed projections, meaningful audit trails, and the flexibility to build new read models from historical events without having to change the write environment.
- **Eventual consistency:**The system is only eventually consistent when it creates materialized views or generates projections of data by replaying events. A delay exists between when an application handles a request and adds events to the event store, when the events publish, and when consumers handle the events. During this period, new events that describe further changes to entities might arrive at the event store. Ensure that your customers understand that data is eventually consistent and that the system is designed to account for eventual consistency in these scenarios.
- **Versioning events:**The event store is the permanent source of information, so you should never update the event data. The only way to update an entity or undo a change is to add a compensating event to the event store. A compensating event is a new event that reverses or corrects the effect of a previous event. For example, a- `ReservationCanceled`event compensates for a prior- `SeatsReserved`event. The original event remains in the stream, and the compensating event records that it was undone.- This immutability also means that if a bug produces incorrect events, those events persist in the store. Fixing the bug in application code doesn't fix the historical events, so you might also need compensating events or upcasters to handle the bad data during replay. If the schema (rather than the data) of the persisted events needs to change, perhaps during a migration, it can be difficult to combine existing events in the store with the new version. - You can use the following strategies individually or in combination: - **Tolerant deserialization:**Design event consumers to ignore unknown fields and use default values for missing fields. This approach handles additive, nonbreaking changes, such as adding an optional field, without requiring any transformation of stored events.
- **Event versioning:**Include a version identifier in each event, either as metadata in the event envelope or as part of the event type name. Consumers use the version to select the appropriate handling logic.
- **Upcasting:**Register transformation functions that convert older event schemas to the current schema during deserialization. You can chain upcasters so that the application code only needs to handle the latest version. The stored events remain unchanged, which preserves immutability.
- **In-place migration:**Rewrite historical events to the new schema directly in the event store. This approach breaks immutability and should be a last resort because it undermines the audit trail.

- **Event ordering:**Multiple-threaded applications and multiple instances of applications might store events in the event store. The consistency of events in the event store and the order of events that affect a specific entity's current state are crucial. Adding a timestamp to every event can help you avoid problems. Another common practice is to annotate each event that results from a request with an incremental identifier. If two actions attempt to add events for the same entity at the same time, the event store can reject an event that matches an existing entity identifier and event identifier.
- **Event querying:**There's no standard approach or existing mechanisms, such as SQL queries, for reading events to obtain information. The only data that you can extract is a stream of events by using an event identifier as the criteria. The event ID typically maps to individual entities. You can determine the current state of an entity only by replaying all of the events that relate to it against the original state of that entity.
- **Event store options:**An event store can be a purpose-built database designed for append-only eventstreams or a general-purpose relational or document database with an append-only table.- Purpose-built event stores provide built-in support for tasks like reading a stream by entity, optimistic concurrency, and snapshots.
- Relational databases are familiar and widely available but require you to build those behaviors yourself.
 - Because each entity has its own independent eventstream, event stores partition naturally by entity ID, which simplifies horizontal scaling or sharding when needed. - Important - Don't confuse an event store with an eventstream message broker. Message brokers such as Apache Kafka typically lack per-entity stream queries and optimistic concurrency. They work well as a distribution layer to fan out events to projections and external consumers, but they aren't a substitute for an event store.
- **Entity state re-creation:**The length of each eventstream affects how you manage and update the system. If the streams are large, replaying every event to rehydrate an entity becomes costly in both time and compute. To mitigate this cost, create snapshots at specific intervals, such as every- *N*events. A snapshot is a serialized representation of the entity's state at a specific point in its eventstream. To rehydrate the entity, load the most recent snapshot and replay only the events that occur after it, rather than replaying the entire stream from the beginning. When you choose a snapshot frequency, balance the storage cost of snapshots against the time saved during rehydration.- Note - Snapshots are an optimization, not a replacement for the eventstream. The eventstream remains the source of truth, and you can regenerate snapshots from it at any time.
- **Conflict handling:**Optimistic concurrency control prevents conflicting writes to the same eventstream, but the application must still handle conflicts that span multiple entities. For example, an event that indicates a reduction in stock inventory might arrive in the data store while a customer places an order for that item. Design the system to reconcile these situations, such as by advising the customer or by creating a back order.
- **Idempotency requirements:**Event delivery to consumers is typically- *at least once*, so consumers can receive the same event more than once. Event handlers must be idempotent so processing a duplicate event doesn't change the outcome. For example, if multiple instances of a consumer process seat-reservation events to maintain an available-seat count, a duplicated reservation event must result in only one decrement. Without idempotency, projections drift from the eventstream and side effects such as payments or notifications trigger more than once. Track the last processed event sequence number for each consumer and skip duplicates, or design state mutations that are inherently safe to repeat.
- **Circular logic:**Be mindful of scenarios in which the processing of one event requires the creation of one or more new events. This sequence can result in an infinite loop.
- **Testing:**A specific testing style best suits event-sourced systems. Set up past events, issue a command, and assert on the new events produced. This- *given-when-then*approach tests business logic without databases, queues, or projections. But you also need integration tests for projections, idempotency behavior, and schema evolution paths, which adds testing surface compared to CRUD systems.
- **Personal data and regulatory compliance:**The append-only, immutable nature of an event store conflicts with data protection regulations that require deletion of personal data, such as the- *right to be forgotten*laws. Deleting events outright breaks stream integrity, so design for this tension from the start.- A common approach is to store personal data outside the event store and reference it by identifier in events. This approach allows deletion to occur independently without affecting the eventstream.
- When you can't separate personal data from events, use crypto-shredding. Encrypt personal data in events by using a per-subject key. Delete the key to render the data unrecoverable while leaving the event structure intact. This approach adds encryption overhead on every read and write and requires robust key management.

## When to use this pattern

Use this pattern when:

- You want to capture intent, purpose, or reason in the data. For example, you can capture changes to a customer entity as a series of specific event types, such as - *Moved home*,- *Closed account*, or- *Deceased*.
- You must minimize or completely avoid conflicting updates to data.
- You want to record events that occur, to replay them to restore the state of a system, to roll back changes, or to keep a history and audit log. For example, when a task consists of multiple steps, you might need to run actions to revert updates and then replay some steps to bring the data back into a consistent state.
- The application already uses events as a natural feature of its operation, and event sourcing requires little extra development or implementation effort.
- You need to decouple the process of inputting or updating data from the tasks required to apply these actions. This change might be to improve UI performance or to distribute events to other listeners that act when the events occur. For example, you can integrate a payroll system with an expense submission website. Both the website and the payroll system consume events that the event store raises in response to data updated on the website.
- You want the flexibility to change the format of materialized models and entity data if requirements change, or when you use CQRS and you need to adapt a read model or the views that expose the data.
- You use CQRS and eventual consistency is acceptable while a read model is updated, or entity and data rehydration from an eventstream results in acceptable performance reduction.

This pattern might not be suitable when:

- Systems have straightforward CRUD operations that don't require auditability, replay, or historical reconstruction of state. The operational overhead of an event store isn't justified if the only requirement is current-state reads and writes.
- Prototypes, minimum viable products (MVPs), or systems have short expected lifespans. The upfront investment in event design, schema evolution strategy, and projection infrastructure rarely yield a return in these scenarios.
- Systems require consistency and real-time updates to the views of the data. Eventual consistency between the event store and projections is inherent to event sourcing.
- Domains in which data is mostly static or for reference, such as lookup tables or catalogs. This type of data changes infrequently and doesn't benefit from change history.
- Teams don't have experience in event-driven architectures. Event sourcing changes how you test, debug, and operate a system. Adopting it without the foundational knowledge increases the risk of antipatterns that are costly to reverse.

Tip

Event sourcing doesn't have to be an all-or-nothing decision for your entire system. Apply it selectively to the parts of your system that it benefits the most, such as a payment ledger or order-processing pipeline. Use traditional CRUD for parts when the complexity isn't justified, such as user profile management or application configuration.

## Workload design

Evaluate how to use the Event Sourcing pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | This pattern can facilitate state reconstruction if you need to recover state stores because you capture a history of changes in complex business processes. - Data partitioning - RE:09 Disaster recovery |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | This pattern, usually combined with CQRS, an appropriate domain design, and strategic snapshotting, can improve workload performance because of atomic append-only operations and the avoidance of database locking for writes and reads. - PE:08 Data performance |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

A conference management system needs to track the number of completed bookings for a conference. By tracking this number, it can check for available seats when a potential attendee tries to make a booking. The system can store the total number of bookings for a conference in at least two ways:

- The system can store information about the total number of bookings as a separate entity in a database that holds booking information. As attendees make or cancel bookings, the system increases or decreases this number. This approach is simple in theory, but it can cause scalability problems if a large number of attendees attempt to book seats during a short period of time. For example, this surge typically occurs on the final day before the booking period closes.
- The system can store information about bookings and cancellations as events held in an event store. It calculates the number of available seats by replaying these events. This approach can be more scalable because of the immutability of events. The system needs to only read data from the event store or append data to the event store. It never modifies event information about bookings and cancellations.

The following diagram shows how you might use event sourcing to implement the seat reservation subsystem of the conference management system.

*Download a Visio file of this architecture.*

### Workflow

The following workflow corresponds to the previous diagram:

- The UI issues a command to reserve seats for two attendees. A separate command handler handles the command. The command handler is a piece of logic that's decoupled from the UI and is responsible for handling requests posted as commands.
- The system constructs an entity that contains information about all reservations for the conference by replaying the events that describe bookings and cancellations. This entity is called - `SeatAvailability`, and it's contained within a domain model that exposes methods for querying and modifying the data in the entity.- Tip - Consider optimizations like snapshots so that you don't need to replay the full list of events to obtain the current state of the entity. Snapshots also maintain a cached copy of the entity in memory.
- The command handler invokes a method that the domain model exposes to make the reservations.
- The - `SeatAvailability`entity raises an event that contains the number of reserved seats. The next time that the entity applies events, it uses all the reservations to compute the number of remaining seats.
- The system appends the new event to the list of events in the event store.

If a user cancels a seat, the system follows a similar process, but the command handler issues a command that generates a seat cancellation event and appends it to the event store.

The system can provide a complete history, or audit trail, of the bookings and cancellations for a conference by using an event store. The events in the event store are the accurate record. You don't need to persist entities in any other way because the system can easily replay the events and restore the state to any point in time.

## Next step

- CQRS pattern: The write store that provides the permanent source of information for a CQRS implementation is typically based on an implementation of the Event Sourcing pattern. The pattern segregates the operations that read data in an application from the operations that update data by using separate interfaces.

## Community resources

- Event Sourcing, by Martin Fowler: The original 2005 description of the pattern that established the foundational vocabulary.
- CQRS Documents (PDF), by Greg Young: The definitive resource about event sourcing and CQRS from the practitioner who formalized both patterns.

## Related resources

The following patterns and guidance might also be relevant when you implement this pattern:

- Materialized View pattern: The data store that you use in an event sourcing system typically isn't suited for efficient querying. Instead, a common approach is to generate prepopulated views of the data at regular intervals or when the data changes.
- Compensating Transaction pattern: The system doesn't update existing data in an event sourcing store. Instead, it adds new entries that transition the state of entities to the new values. To reverse a change, it uses compensating entries because it can't reverse the previous change. The Compensating Transaction pattern article describes how to undo the work that a previous operation performed.
- Domain analysis for microservices: In systems that use domain-driven design (DDD), the entity that owns an eventstream is typically an aggregate, a consistency boundary that receives commands, enforces business rules, and emits events.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Move configuration information out of the application deployment package to a centralized location. This approach provides easier management and control of configuration data, and to share configuration data across applications and application instances.

## Context and problem

Most application runtime environments include configuration information in files that you deploy with the application. In some cases, you can edit these files to change the application's behavior after you deploy the application. However, configuration changes require you to redeploy the application. Redeployment often results in unacceptable downtime and other administrative overhead.

Local configuration files also limit the configuration to a single application. In some scenarios, you might want to share configuration settings across multiple applications. Examples include database connection strings, UI theme information, and the URLs of queues and storage that a related set of applications uses.

Managing changes to local configurations across multiple running instances of the application is challenging. This challenge can result in instances that use different configuration settings while you deploy the update.

Updates to applications and components might also require changes to configuration schemas. Many configuration systems don't support different versions of configuration information.

## Solution

Store the configuration information in external storage, and provide an interface that you can use to quickly and efficiently read and update configuration settings. The type of external store depends on the application's hosting and runtime environment. In a cloud-hosted scenario, external storage is typically a cloud-based storage service or dedicated configuration service. It might also be a hosted database or other custom system.

The backing store that you choose for configuration information should have an interface that provides consistent and easy-to-use access. It should expose the information in a correctly typed and structured format. The implementation might also need to authorize users' access to protect configuration data. It might need to be flexible enough to store multiple versions of the configuration, such as development, staging, and production, including multiple release versions of each configuration.

Many built-in configuration systems read the data when the application starts and then cache the data in memory to provide fast access and minimize the impact on application performance. Depending on the type of backing store that you use and the latency of this store, you might want to implement a caching mechanism within the external configuration store. For more information, see Caching guidance. The following diagram shows an overview of the External Configuration Store pattern with an optional local cache.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- Choose a backing store that provides acceptable performance, high availability, and robustness. Ensure that you can back it up in your application maintenance and administration process. In a cloud-hosted application, use a cloud storage mechanism or a dedicated configuration platform service to meet these requirements.
- Design the schema of the backing store to allow flexibility in the types of information that it can hold. Ensure that it provides capabilities for all configuration requirements, such as typed data, collections of settings, multiple versions of settings, and any other features that the applications require. The schema should be easy to extend to support more settings when requirements change.
- Consider the physical capabilities of the backing store, how they relate to the way it stores configuration information, and the effects on performance. For example, storing an XML document that contains configuration information requires either the configuration interface or the application to parse the document to read individual settings. Parsing complicates how you update settings, but caching the settings can help offset slower read performance.
- Consider how the configuration interface permits control of the scope and inheritance of configuration settings. For example, you might need to scope configuration settings at the organization, application, and machine levels. The configuration interface might need to delegate control over access to different scopes and prevent or allow individual applications to override settings.
- Ensure that the configuration interface can expose the configuration data in the required formats, such as typed values, collections, key-value pairs, and property bags.
- Consider how the configuration store interface behaves when settings contain errors or don't exist in the backing store. You might need to restore default settings and log errors. Also consider the case sensitivity of configuration setting keys or names, how to store and handle binary data, and how to handle null or empty values.
- Consider how to protect the configuration data and give access to only the appropriate users and applications. The configuration store interface typically provides this feature, but you also need to ensure that users and applications can't directly access the data in the backing store without the appropriate permissions. Ensure strict separation between the permissions required to read and write configuration data. Also consider whether you need to encrypt some or all of the configuration settings and how to implement this encryption in the configuration store interface. - You should also turn on audit logging to record who reads or modifies configuration values and when these actions occur. Apply the same audit requirements to any local fallback copies of configuration data.
- Separate nonsensitive configuration values from secrets. Keep routine settings, such as feature flags and endpoints, in the configuration settings. Store secrets, such as connection strings, API keys, certificates, and passwords, in a dedicated secret-management system that provides encryption and controlled access.
- Centrally stored configurations, which change application behavior during runtime, are crucial. Deploy, update, and manage them by using the same mechanisms that you use to deploy application code. For example, you must carry out changes that can affect more than one application by using a fully tested-and-staged deployment approach to ensure that the change suits all applications that use this configuration. If an administrator edits a setting to update one application, it might adversely affect other applications that use the same setting. Products like Azure App Configuration help mitigate this risk through built-in capabilities, such as revision history, point-in-time recovery (PITR), immutable snapshots, and progressive rollout patterns.
- If an application caches configuration information, you need to alert the application when the configuration changes. You might implement an expiration policy for cached configuration data so that this information automatically refreshes periodically. The application sees the changes and implements them.
- Cached configuration data can help address transient connectivity problems that the external configuration store experiences at application runtime, but this approach typically doesn't solve the problem if the external store is down when the application starts. Ensure that your application deployment pipeline can provide the last known set of configuration values in a configuration file to use when your application can't retrieve live values at startup.

## When to use this pattern

Use this pattern when:

- You need to share configuration settings across multiple applications or instances or enforce a standard configuration across them.
- Your standard configuration system doesn't support all required setting types, such as images or complex data structures.
- You need a complementary store for some settings, while allowing applications to override some or all centrally stored values.
- You need to simplify administration across multiple applications and optionally monitor configuration usage by recording access to the configuration store.

This pattern might not be suitable when:

- Your configuration is simple, local to one application, and changes only during normal release cycles. In this case, an external configuration store can add unnecessary operational complexity.

## Workload design

Evaluate how to use the External Configuration Store pattern in a workload design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | This separation of application configuration from application code supports environment-specific configuration and applies versioning to configuration values. External configuration stores are also a common place to manage feature flags to implement safe deployment practices. - OE:10 Automation design - OE:11 Safe deployment practices |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

The following examples show how to implement the External Configuration Store pattern in Azure. The first example uses App Configuration and client libraries. The second example uses a custom backing store for scenarios that require specialized implementation.

### App Configuration

Most applications can use App Configuration instead of a custom configuration store. App Configuration supports key-value pairs that you can apply namespaces to. App Configuration also supports immutable snapshots of configuration so that you can inspect, roll back, or progressively deploy configuration changes without risk to running instances.

Use snapshot references to let applications switch between snapshots at runtime without code changes or redeployment. You can export configuration values so that a copy ships with your application as a backup to use if the service is unreachable when the application starts.

In App Configuration, keys and values are Unicode strings, and each key-value pair has optional metadata, such as label-based variants and content type. Use content type to describe how your application should interpret a value, such as in JSON or in a built-in App Configuration type. App Configuration also keeps a revision history with PITR, which helps you review and recover previous key-value pairs.

For resiliency, provision your store in a region that supports availability zones and turn on geo-replication so that you can configure your applications to read from the nearest replica and switch between replica endpoints during regional outages. Use Azure Key Vault references to keep secrets in Key Vault and reference them from App Configuration, rather than storing credentials directly in the configuration store. Use managed identity and Azure role-based access control (Azure RBAC) instead of connection strings to authenticate applications to App Configuration.

For workloads that run in Azure Kubernetes Service (AKS), the App Configuration Kubernetes Provider can generate ConfigMaps and Secrets directly from your store without requiring code changes in your workload containers. You can also use App Configuration to manage feature flags, including targeted rollout and variant-based experimentation, in your safe deployment practices.

For network isolation, use private endpoints for App Configuration so that client traffic remains on private IP addresses through Azure Private Link. After you set up private access, you can turn off public access to reduce public endpoint exposure. In geo-replicated deployments, a single private endpoint can reach all replicas, but for higher regional resilience, you can provision private endpoints for each replica region and set up the Domain Name System (DNS) accordingly.

#### Client libraries

Client libraries provide many of the preceding features. Client libraries integrate with the application runtime to help fetch and cache values, refresh values when they change, and handle transient outages in App Configuration.

| Runtime | Client library | Notes | Quickstart |
|---|---|---|---|
| .NET | Microsoft.Extensions.Configuration.AzureAppConfiguration | Provider for `Microsoft.Extensions.Configuration` | Quickstart for .NET |
| ASP.NET Core | Microsoft.Azure.AppConfiguration.AspNetCore | Adds request-driven refresh middleware for ASP.NET Core | Quickstart for ASP.NET Core |
| Azure Functions in .NET | Microsoft.Azure.AppConfiguration.Functions.Worker | Provider for the isolated worker model that uses `Program.cs` | Quickstart for Azure Functions |
| .NET Framework | Microsoft.Configuration.ConfigurationBuilders.AzureAppConfiguration | Configuration builder for `System.Configuration` | Quickstart for .NET Framework |
| Java Spring | com.azure.spring > azure-spring-cloud-appconfiguration-config | Supports Spring Framework access via `ConfigurationProperties` | Quickstart for Java Spring |
| Python | azure-appconfiguration-provider | Provider library that provides dynamic refresh and Key Vault reference support | Quickstart for Python |
| JavaScript and Node.js | @azure/app-configuration-provider | Provider library that provides dynamic refresh and Key Vault reference support | Quickstart for JavaScript |

The following App Configuration sync GitHub Action and built-in Azure Pipelines tasks are also available:

### Custom backing store example

In an application that Azure hosts, you can use Azure Storage to store configuration information externally. This approach provides resiliency and high performance. By default, Storage replicates data three times within a single datacenter. For geo-redundancy across regions, you can set up geo-replication with manual failover capabilities. Azure Table Storage provides a key-value store that can use a flexible schema for the values. Azure Blob Storage provides a hierarchical, container-based store that can hold any type of data in individually named blobs.

When you implement this pattern, you need to abstract Blob Storage and expose your settings within your applications. You also need to check for updates at runtime and decide how to respond to those updates.

The following example shows how you can use a simple configuration store and Blob Storage to store and expose configuration information. A `BlobSettingsStore` class abstracts Blob Storage for holding configuration information. It implements a simple `ISettingsStore` interface.

```
public interface ISettingsStore
{
 Task<ETag> GetVersionAsync();
 Task<Dictionary<string, string>> FindAllAsync();
}
```
This interface defines methods for retrieving configuration settings that the configuration store holds and includes a version number that you can use to detect recent configuration setting modifications. A `BlobSettingsStore` class can use the `ETag` property of the blob to implement versioning. The `ETag` property updates automatically each time a blob is written.

Note

By design, this simple illustration exposes all configuration settings as string values rather than typed values.

An `ExternalConfigurationManager` class provides a wrapper around a `BlobSettingsStore` instance. An application can use this class to retrieve configuration information. This class might use a change notification mechanism, such as Microsoft Reactive Extensions, to publish configuration updates while the system runs. It also implements the Cache-Aside pattern for settings to provide better resiliency and performance.

The following example shows how you might implement an `ExternalConfigurationManager` class.

```
static void Main(string[] args)
{
 // Start monitoring configuration changes.
 ExternalConfiguration.Instance.StartMonitor();
 // Get a setting.
 var setting = ExternalConfiguration.Instance.GetAppSetting("someSettingKey");
 …
}
```
## Next steps

- App Configuration samples
- Integrate App Configuration with Kubernetes deployments by using Helm
- Manage feature flags in App Configuration
- Caching guidance
- App Configuration best practices

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Delegate user authentication to an external identity provider (IdP) to simplify development, minimize administrative tasks, and improve application UX.

## Context and problem

Users typically need to work with multiple applications that partner organizations provide and host. They might need to use specific, different sign-in credentials for each application. This requirement can:

- **Cause a disjointed UX.**Employees often forget multiple sign-in credentials.
- **Expose security vulnerabilities.**When an employee leaves the company, the organization must immediately deactivate the account. Large organizations often miss this critical step.
- **Complicate user management.**Admins manage user credentials, issue password reminders, and perform other administrative tasks.

Users typically prefer to use the same sign-in credentials for all applications.

## Solution

Implement a federated identity authentication mechanism. Separate user authentication from the application code and delegate authentication to a trusted IdP. This process simplifies development, minimizes administrative overhead, and provides user authentication via a range of IdPs. Federated identity also separates authentication from authorization.

Trusted IdPs include corporate directories, on-premises federation services, security token services (STSs), and social IdPs like Microsoft, Google, Yahoo!, or Facebook.

The following diagram shows the Federated Identity pattern for a client application that accesses a service that requires authentication. The IdP works with an STS to provide authentication. The IdP issues security tokens that provide information about the authenticated user. This information, called *claims*, includes the user's identity and might also include other claims, like role memberships and more granular access rights.

This model is also called claims-based access control. Applications and services authorize access to features and functionality based on the claims. The service that requires authentication must trust the IdP. The client application contacts the IdP for authentication. If authentication succeeds, the IdP returns a token that contains user-identifying claims to the STS. The IdP and the STS might be part of the same service. The STS can transform and augment the claims based on predefined rules before it returns the token to the client. The client application then passes this token to the service as proof of its identity.

Federated authentication provides a standards-based method to establish trust in identities across domains, and it supports single sign-on (SSO). Many applications, especially cloud-hosted applications, use federated authentication because it supports SSO without a direct network connection to an IdP. This design increases security, because the user doesn't need to create and enter different sign-in credentials for multiple applications. It also limits credential exposure to only the original IdP. Applications see only the authenticated identity information in the token.

Applications and services that use federated authentication don't need to provide identity management features. Instead, the IdP is responsible for identity and credential management. When the corporate directory trusts the IdP, it doesn't need to manage the user identity. This approach eliminates the administrative overhead of directory-based user identity management.

## Problems and considerations

Consider the following points when you decide how to implement this pattern:

- Authentication can be a single point of failure. To maintain application reliability and availability across multiple regions, consider deploying your identity management mechanism across the same regions as your application.
- To configure role-based access control (RBAC), use authentication tools. RBAC supports granular control over feature and resource access.
- Unlike a corporate directory, claims-based authentication that uses social IdPs usually provides only the authenticated user's email address, and sometimes their name. Some social IdPs, such as Microsoft, provide only a unique identifier. The application usually maintains some information about registered users so it can match this information to the identifier in the claims. This task is typically completed during registration, when the user first accesses the application. Information is then injected into the token as new claims after each authentication.
- If multiple IdPs are configured for the STS, the STS must determine which IdP should authenticate the user. This process is called - *home realm discovery*. The STS might determine the IdP automatically based on user-provided information like an email address or a user name, the application subdomain, the user's IP address range, or a cookie stored in the user's browser. For example, if the user enters a Microsoft email address, such as- `user@live.com`, the STS redirects the user to the Microsoft account sign-in page. On subsequent visits, the STS can use a cookie that indicates the user previously signed in by using a Microsoft account. If the STS can't automatically determine the home realm, it displays a home realm discovery page that lists the trusted IdPs. The user then selects an IdP.

## When to use this pattern

Use this pattern when you need:

- **SSO in the enterprise.**In this scenario, you need to authenticate employees for corporate applications that are hosted in the cloud outside the corporate security boundary, without a sign-in every time they visit an application. The user experience matches on-premises applications. Users authenticate when they sign in to the corporate network, and then they can access relevant applications without another sign-in.
- **Federated identity with multiple partners.**In this scenario, you need to authenticate corporate employees and business partners who don't have accounts in the corporate directory. This practice is common in business-to-business applications, applications that integrate with partner services, and in companies that use different IT systems, or merged or shared resources.
- **Federated identity in software as a service (SaaS) applications.**In this scenario, independent software vendors provide a ready-to-use service for multiple clients or tenants. Tenants authenticate by using a suitable IdP. For example, business users use their corporate credentials, while tenant consumers and clients use social identity credentials.
- **Federated identity for workload access.**In this scenario, tenant applications, automation workflows, or continuous integration and continuous delivery systems need to call your APIs without a user present. Tenants authenticate via their own IdPs by using workload identities. The application authorizes access by using tenant-scoped claim validation.

This pattern might not be suitable when you have:

- **One IdP.**In this scenario, application users authenticate by using one IdP, and they don't need to authenticate by using another IdP. This situation is typical in applications that use a corporate directory for authentication, either via a VPN or a virtual network connection between the application and an on-premises directory.
- **Incompatible authentication mechanisms.**In this scenario, the application uses a different authentication mechanism, for example by using custom user stores, or it can't handle claims-based technology negotiation standards. It can be complex and expensive to retrofit claims-based authentication and access control into an existing application.

## Workload design

Evaluate how to use the Federated Identity pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | This pattern offloads user management and authentication to the IdP, which typically has a high service-level objective. During workload disaster recovery (DR), the workload recovery plan doesn't need to address authentication components. - RE:02 Critical flows - RE:09 DR |
| Security design decisions help ensure the confidentiality,integrity, andavailabilityof your workload's data and systems. | This pattern provides advanced identity-based threat detection and prevention capabilities without requiring you to implement them in your workload. External IdPs also use modern interoperable authentication protocols. - SE:02 Secured development lifecycle - SE:10 Monitoring and threat detection |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | This pattern helps you devote application resources to other priorities. - PE:03 Selecting services |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

An organization hosts a multicomponent cloud-based application that includes a web front end and a back-end API. The application delegates authentication to a centralized IdP by using Microsoft Entra ID, rather than implementing authentication logic in each component.

*Download a Visio file of this architecture.*

The following workflow corresponds to the previous diagram.

- The user accesses the web app.
- The web app redirects the user to Microsoft Entra ID for authentication.
- After successful authentication, Microsoft Entra ID redirects the user back to the web app with an authorization code.
- The web app exchanges the authorization code for tokens and sends a POST request to the token endpoint.
- Microsoft Entra ID issues a token that contains claims about the user.
- The web app uses this token to call a back-end API.
- The web app and the back-end API validate the token and enforce their authorization rules based on the claims.
- The API returns the response to the web app.

Key characteristics:

- **Centralized authentication.**Components rely on Microsoft Entra ID to authenticate users, which removes the need for custom authentication logic in the application.
- **Decentralized authorization.**Application components independently enforce authorization decisions based on claims.
- **Claims-based access control.**Access to functionality is determined by using claims, such as roles or scopes.
- **Standards-based protocols.**Components use OAuth 2.0 and OpenID Connect for authentication.
- **Optional MFA enforcement.**If your risk profile requires stronger sign-in assurance, you can enforce multifactor authentication by using Conditional Access policies in Microsoft Entra ID.
- **Optional extensibility through federation.**Microsoft Entra ID can be configured to trust a partner Microsoft Entra tenant by using cross-tenant access settings. Partner users can then access the application without changes to application components.

## Next steps

- What is Microsoft Entra?
- OpenID Connect on the Microsoft identity platform
- Convert a single-tenant app to multitenant by using Microsoft Entra ID
- What is Conditional Access?

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Protect applications and services by using a dedicated component to broker requests between clients and the application or service. The broker validates and sanitizes the requests, and can provide an extra layer of security and limit the system's attack surface.

## Context and problem

Many cloud services expose endpoints that allow client applications to call their APIs across the internet or another untrusted network. The code that implements the APIs triggers or performs several tasks, including but not limited to authentication, authorization, parameter validation, and some or all request processing. The API code is likely to access storage and other services on the client's behalf.

If a malicious user compromises the system and gains access to the application's hosting environment, its security mechanisms and access to data and other services are exposed. As a result, the malicious user can gain unrestricted access to credentials, storage keys, sensitive information, and other services.

## Solution

One solution to this problem is to decouple the code that implements public endpoints from the code that processes requests and accesses storage. Decouple the code by using a façade tier that interacts with clients and routes approved requests through an internal endpoint, queue, or broker to the workload components that handle the business operation. The diagram provides a high-level overview of this pattern.

You can use the Gatekeeper pattern to protect storage, or you can use it as a more comprehensive façade to protect all of the functions of the application. Important factors include:

- **Controlled validation:**The gatekeeper validates all requests, and rejects requests that don't meet validation requirements.
- **Limited risk and exposure:**Risks and exposure are reduced because the gatekeeper doesn't access the credentials or keys that the trusted host uses to access storage and services. If the gatekeeper becomes compromised, attackers can't access these credentials or keys.
- **Appropriate security:**The gatekeeper runs in a limited privilege mode, while the rest of the application runs in the full trust mode required to access storage and services. If the gatekeeper is compromised, it can't directly access the application services or data.

This pattern acts like a firewall in a typical network topography. Unlike a traditional firewall, it allows the gatekeeper to examine requests in detail and make an application-driven decision about whether to pass the request to the trusted host that performs the required tasks. This decision typically requires the gatekeeper to validate and sanitize the request content before it passes it to the trusted host. Gatekeepers might authorize the request, look for unexpected or invalid payload content, perform rate limiting, and perform various other checks.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- Ensure that the trusted hosts expose only internal or protected endpoints that only the gatekeeper uses. The trusted hosts shouldn't expose any external endpoints or interfaces.
- The gatekeeper must run in a limited-privilege mode. In practice, host the gatekeeper and trusted back end on separate compute boundaries, and keep back-end endpoints private.
- The gatekeeper shouldn't perform processing related to the application or services or access data. Its function is solely to validate and sanitize requests. The trusted hosts might need to perform extra request validation, but the gatekeeper should perform the core validation.
- Use a secure communication channel like HTTPS, Secure Sockets Layer (SSL), or Transport Layer Security (TLS) between the gatekeeper and the trusted hosts or tasks where possible. However, some hosting environments don't support HTTPS on internal endpoints.
- Adding the extra layer to implement the Gatekeeper pattern is likely to affect performance because of the extra processing and network communication required.
- The gatekeeper can be a single point of failure (SPoF). To minimize the impact of a failure, consider deploying redundant instances and using an autoscaling mechanism to ensure capacity and maintain availability.

## When to use this pattern

Use this pattern when:

- You handle sensitive information.
- You expose services that require strong protection from malicious traffic.
- You perform mission-critical operations that can't tolerate direct exposure of back-end services.
- You need request validation and sanitization to be separated from core business processing.

This pattern might not be suitable when:

- You can satisfy security and validation requirements through built-in platform controls on the back-end service without adding a dedicated gatekeeper tier.
- Added network hops and validation latency violate strict end-to-end latency requirements.

## Workload design

Evaluate how to use the Gatekeeper pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Security design decisions help ensure the confidentiality,integrity, andavailabilityof your workload's data and systems. | A gatekeeper in the request flow helps you centralize security functionality like web application firewalls, DDoS protection, bot detection, request manipulation, authentication initiation, and authorization checks. - SE:06 Network controls - SE:10 Monitoring and threat detection |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | You can use this pattern to implement throttling at a gatekeeper level rather than implement rate checks at the node level. Rate state coordination among all nodes isn't inherently performant. - PE:03 Select services |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

The Gatekeeper pattern typically implements a layered request path, where each layer has a specific responsibility and a limited trust scope.

*Download a Visio file of this architecture.*

In this design, Azure Application Gateway with Azure Web Application Firewall is the outer gatekeeper. It inspects internet-facing traffic and applies security controls before traffic reaches the API tier. Azure API Management is the inner gatekeeper. It applies API-specific controls and forwards only approved traffic to private back ends.

For example, Azure Web Application Firewall can detect and block SQL injection and cross-site scripting patterns, enforce protocol and request-size rules, and apply bot and IP-based filtering before requests reach API Management or private back ends.

When you use API Management in the inner layer, it applies policies to inbound requests and outbound responses in the gateway pipeline. For more information about how API Management processes requests and responses, see Policies in API Management. For policy options such as JSON Web Token (JWT) validation, rate limiting, header transformation, and response shaping, see API Management policy reference.

Use managed identities for Azure resources consistently for service-to-service authentication in this path. For example, API Management can use the authenticate with managed identity policy to get Microsoft Entra tokens for back-end calls without storing secrets.

The back end remains private. For example, the back end can be an Azure App Service app that uses a private endpoint, so the app can be accessed privately.

For containerized workloads, an alternative can replace the API Management plus App Service inner path with ingress-based compute:

- Azure Kubernetes Service (AKS), which gives you more control over ingress controller choice, Kubernetes policies, network topology, and cluster operations.
- Azure Container Apps, which is a serverless managed container platform that provides ingress capabilities and reduces infrastructure management.

In these alternatives, ingress can route by host or path, terminate TLS, and expose internal-only services. Specific capabilities such as request limits and allow or deny rules depend on the selected ingress implementation. In all cases, keep the gatekeeper boundaries: apply validation and policy enforcement at ingress, and keep back-end services reachable only through that gatekeeper path.

Each layer in this path emits logs and metrics that you should centralize. Azure Web Application Firewall diagnostic logs record matched and blocked rules per request. API Management emits gateway logs that capture request duration, response codes, and policy outcomes. Back-end services emit application-level telemetry. Collect these logs and metrics in Azure Monitor and route them to a Log Analytics workspace for unified querying. Standardize end-to-end request correlation by generating or forwarding a correlation ID at the edge and propagating it through API Management and back-end services (for example, through request headers and distributed trace context) so that a single transaction remains traceable across all layers. Use Microsoft Defender for Cloud to surface security recommendations across the gatekeeper components. Configure alerts on anomalous Azure Web Application Firewall block rates or API Management error spikes to detect threats before they reach private back ends.

## Next steps

The following guidance might be relevant when you implement this pattern:

- Azure Web Application Firewall on Application Gateway
- Policies in API Management
- Use private endpoints for App Service apps

## Related resources

The following cloud design patterns are often used together with the Gatekeeper pattern:

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Use a gateway to aggregate multiple individual requests into a single request. This pattern is useful when a client must make multiple calls to different back-end systems to perform an operation.

## Context and problem

To perform a single task, a client might have to make multiple calls to various back-end services. An application that relies on many services to perform a task must expend resources on each request. When new features or services are added to the application, extra requests are needed, which further increases resource requirements and network calls. This chattiness between a client and a back end can adversely affect the performance and scale of the application. Microservices architectures have made this problem more common because applications built around many smaller services have a higher number of cross-service calls.

In the following diagram, the client sends requests to each service (numbered 1, 2, and 3). Each service processes the request and returns a response to the application (numbered 4, 5, and 6). Sending individual requests in this way over a cellular network that has high latency is inefficient and can cause connectivity loss or incomplete responses. Each request might run in parallel. However, the application still must send, wait for, and process data for each request on separate connections, which increases the chance of failure.

## Solution

Use a gateway to reduce chattiness between the client and the services. The gateway receives client requests, dispatches requests to the various back-end systems, and aggregates the results before it sends them back to the client.

This pattern can reduce the number of requests that the application makes to back-end services and improve application performance over high-latency networks.

In the following diagram, the application sends a request to the gateway (1). The request contains a package of extra requests. The gateway decomposes these requests and processes each request by sending it to the relevant service (2). Each service returns a response to the gateway (3). The gateway combines the responses from each service and sends the response to the application (4). The application makes a single request and receives only a single response from the gateway.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- The gateway shouldn't introduce service coupling across the back-end services.
- The gateway should be located near the back-end services to reduce latency as much as possible.
- The gateway service might introduce a single point of failure (SPoF). Ensure that the gateway is properly designed to meet your application's availability requirements.
- The gateway might introduce a bottleneck. Ensure that the gateway has adequate performance to handle the current load and can be scaled to meet your anticipated growth.
- Perform load testing against the gateway to ensure that you don't introduce cascading failures for services.
- Implement a resilient design by using techniques such as bulkheads, circuit breaking, retry, and timeouts.
- If one or more service calls take too long, it might be acceptable to time out and return a partial set of data. Consider how your application will handle this scenario.
- Use asynchronous input and output (I/O) to ensure that a delay at the back end doesn't cause performance problems in the application.
- Implement distributed tracing by using correlation IDs to track each individual call.
- Monitor request metrics and response sizes.
- Consider returning cached data as a failover strategy to handle failures.
- Rather than build aggregation into the gateway, consider placing an aggregation service behind the gateway. Request aggregation is likely to have different resource requirements than other services in the gateway and might affect the gateway's routing and offloading functionality.

## When to use this pattern

Use this pattern when:

- A client needs to communicate with multiple back-end services to perform an operation.
- The client might use networks that have significant latency, such as cellular networks.

This pattern might not be suitable when:

- You want to reduce the number of calls between a client and a single service across multiple operations. In that scenario, adding a batch operation to the service might be more suitable.
- The client or application is located near the back-end services and latency isn't a significant factor.

## Workload design

Evaluate how to use the Gateway Aggregation pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | With this topology, you can shift transient fault handling from a distributed implementation across clients to a centralized implementation. - Recommendations for handling transient faults |
| Security design decisions help ensure the confidentiality,integrity, andavailabilityof your workload's data and systems. | This topology often reduces the number of touchpoints that a client has with a system, which reduces the public surface area and authentication points. The aggregated back ends can remain fully network-isolated from clients. - SE:04 Segmentation - SE:08 Harden resources |
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | This pattern enables back-end logic to evolve independently from clients. This decoupling gives you the flexibility to change the chained service implementations, or even data sources, without needing to change client touchpoints. - OE:04 Tools and processes |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | This design can incur less latency than a design in which the client establishes multiple connections. Caching in aggregation implementations minimizes calls to back-end systems. - PE:03 Select services - PE:08 Data performance |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

Consider a microservices-based application that provides an order summary experience for a customer. When a user opens an order page, the application must retrieve data from multiple back-end services, such as an order service, a shipment service, and a customer profile service.

In a microservices architecture, these services are implemented and deployed independently. Without aggregation, the client must call each service directly, which increases latency and complexity.

To address this problem, the application uses Azure API Management as the gateway aggregation layer. The client sends a single request to an API Management operation that acts as a collector for order information. API Management then calls the supporting back-end APIs and returns a unified response to the client.

You can implement this lightweight composition by using the API Management send-request policy to retrieve data from multiple services and construct a combined response. In this example, the back-end services run in an Azure Container Apps environment, and you deploy each back-end service as a container app that remains hidden from direct client access.

*Download a Visio file of this architecture.*

The request flow follows these steps:

- The client sends a request to an order summary endpoint exposed through API Management.
- API Management applies a policy that collects the order, shipment, and customer profile data from the back-end services.
- API Management composes the back-end responses into a single order summary payload.
- API Management returns the aggregated response to the client.

By introducing this aggregation layer, the solution reduces client-to-service round trips and simplifies client interactions. This layer becomes responsible for handling unresponsive back-end services gracefully and preventing failures from cascading across the aggregated response. Harden your API Management policies by using per-request timeouts, conditional error handling, and circuit breakers.

If one of the back-end calls times out or returns an error, API Management can apply the behavior that best fits the operation. For example, it might return a partial response when missing data is acceptable, or it might fail the entire request when complete and consistent order data is required. Make this decision explicit in the policy design so that clients experience predictable behavior.

This approach works well when the gateway performs lightweight composition, shaping, and response assembly. If the aggregation requires custom domain logic, complex transformations, or longer-running orchestration, place that functionality in a dedicated custom service behind the gateway.

For monitoring, collect telemetry across the full request path so that you can correlate API Management behavior with back-end latency. This visibility is important in a gateway aggregation pattern because a single client operation depends on multiple back-end calls, and failures or slow responses in any one dependency can affect the final aggregated result. Use Azure Monitor as the central observability platform. Collect API Management logs and metrics for the gateway and policy execution path, and enable monitoring for Container Apps to capture application logs and metrics from the back-end container apps. Route API Management and back-end telemetry to a Log Analytics workspace for unified querying, alerting, and troubleshooting. With this telemetry, you can detect timeout patterns, identify which back-end dependency caused a partial or failed response, and create alerts for elevated latency or error rates.

## Next steps

- Use external services from the API Management service
- Send-request policy
- Container Apps documentation

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Offload shared or specialized service functionality to a gateway proxy. This pattern can simplify application development by moving shared service functionality, such as the use of SSL certificates, from other parts of the application into the gateway.

## Context and problem

Some features are commonly used across multiple services, and these features require configuration, management, and maintenance. A shared or specialized service that is distributed with every application deployment increases the administrative overhead and increases the likelihood of deployment error. Any updates to a shared feature must be deployed across all services that share that feature.

Properly handling security issues (token validation, encryption, SSL certificate management) and other complex tasks can require team members to have highly specialized skills. For example, a certificate needed by an application must be configured and deployed on all application instances. With each new deployment, the certificate must be managed to ensure that it doesn't expire. Any common certificate that is due to expire must be updated, tested, and verified on every application deployment.

Other common services such as authentication, authorization, logging, monitoring, or throttling can be difficult to implement and manage across a large number of deployments. It might be better to consolidate this type of functionality, in order to reduce overhead and the chance of errors.

## Solution

Offload some features into a gateway, particularly cross-cutting concerns such as certificate management, authentication, SSL termination, monitoring, protocol translation, or throttling.

The following diagram shows a gateway that terminates inbound SSL connections. It requests data on behalf of the original requestor from any HTTP server upstream of the gateway.

Benefits of this pattern include:

- Simplify the development of services by removing the need to distribute and maintain supporting resources, such as web server certificates and configuration for secure websites. Simpler configuration results in easier management and scalability and makes service upgrades simpler.
- Allow dedicated teams to implement features that require specialized expertise, such as security. This allows your core team to focus on the application functionality, leaving these specialized but cross-cutting concerns to the relevant experts.
- Provide some consistency for request and response logging and monitoring. Even if a service isn't correctly instrumented, the gateway can be configured to ensure a minimum level of monitoring and logging.
- Centralize carbon-aware traffic management. A gateway can adjust caching, rate limiting, and logging behaviors based on real-time carbon intensity signals.

## Issues and considerations

- Ensure the gateway is highly available and resilient to failure. Avoid single points of failure by running multiple instances of your gateway.
- Ensure the gateway is designed for the capacity and scaling requirements of your application and endpoints. Make sure the gateway doesn't become a bottleneck for the application and is sufficiently scalable.
- Only offload features that are used by the entire application, such as security or data transfer.
- Business logic should never be offloaded to the gateway.
- If you need to track transactions, consider generating correlation IDs for logging purposes.

## When to use this pattern

Use this pattern when:

- An application deployment has a shared concern such as SSL certificates or encryption.
- A feature that is common across application deployments that might have different resource requirements, such as memory resources, storage capacity or network connections.
- You wish to move the responsibility for issues such as network security, throttling, or other network boundary concerns to a more specialized team.

This pattern might not be suitable if it introduces coupling across services.

## Workload design

An architect should evaluate how the Gateway Offloading pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and to ensure that itrecoversto a fully functioning state after a failure occurs. | Offloading this responsibility to a gateway reduces the complexity of application code on backend nodes. In some cases, offloading completely replaces functionality with a reliable platform-provided feature. - RE:01 Simplicity and efficiency |
| Security design decisions help ensure the confidentiality,integrity, andavailabilityof your workload's data and systems. | Adding a gateway into the request flow enables you to centralize security functionality like web application firewalls and TLS connections with clients. Any offloaded functionality that's platform-provided already offers enhanced security. - SE:06 Network controls - SE:08 Hardening resources |
| Cost Optimization is focused on sustaining and improvingyour workload'sreturn on investment. | This pattern enables you to redirect costs from resources that would be spent per-node into the gateway implementation. Costs in the centralized processing model are frequently lower than those of the distributed model. - CO:14 Consolidation |
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | In this pattern, the configuration and upkeep of the offloaded functionality is from single point instead of managing it from multiple nodes. - OE:04 Tools and processes |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, code. | Adding an offloading gateway to the request process enables you to use less resources per-node because functionality is centralized at the gateway. You can optimize the implementation of the offloaded functionality independently of the application code. Offloaded platform-provided functionality is already likely to be highly performant. - PE:03 Selecting services |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

Using Nginx as the SSL offload appliance, the following configuration terminates an inbound SSL connection and distributes the connection to one of three upstream HTTP servers.

```
upstream iis {
 server 10.3.0.10 max_fails=3 fail_timeout=15s;
 server 10.3.0.20 max_fails=3 fail_timeout=15s;
 server 10.3.0.30 max_fails=3 fail_timeout=15s;
}
server {
 listen 443;
 ssl on;
 ssl_certificate /etc/nginx/ssl/domain.cer;
 ssl_certificate_key /etc/nginx/ssl/domain.key;
 location / {
 set $targ iis;
 proxy_pass http://$targ;
 proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
 proxy_set_header X-Forwarded-Proto https;
 proxy_set_header X-Real-IP $remote_addr;
 proxy_set_header Host $host;
 }
}
```
On Azure, this can be achieved by setting up SSL termination on Application Gateway.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Route requests to multiple services or multiple service instances using a single endpoint. The pattern is useful when you want to:

- Expose multiple services on a single endpoint and route to the appropriate service based on the request
- Expose multiple instances of the same service on a single endpoint for load balancing or availability purposes
- Expose differing versions of the same service on a single endpoint and route traffic across the different versions

## Context and problem

When a client needs to consume multiple services, multiple service instances or a combination of both, the client must be updated when services are added or removed. Consider the following scenarios.

- **Multiple disparate services**- An e-commerce application might provide services such as search, reviews, cart, checkout, and order history. Each service has a different API that the client must interact with, and the client must know about each endpoint in order to connect to the services. If an API changes, the client must be updated as well. If you refactor a service into two or more separate services, the code must change in both the service and the client.
- **Multiple instances of the same service**- The system can require running multiple instances of the same service in the same or different regions. Running multiple instances can be done for load balancing purposes or to meet availability requirements. Each time an instance is spun up or down to match demand, the client must be updated.
- **Multiple versions of the same service**- As part of the deployment strategy, new versions of a service can be deployed along side existing versions. This is known as blue green deployments. In these scenarios, the client must be updated each time there are changes to the percentage of traffic being routed to the new version and existing endpoint.

## Solution

Place a gateway in front of a set of applications, services, or deployments. Use application Layer 7 routing to route the request to the appropriate instances.

With this pattern, the client application only needs to know about a single endpoint and communicate with a single endpoint. The following illustrate how the Gateway Routing pattern addresses the three scenarios outlined in the context and problem section.

### Multiple disparate services

The gateway routing pattern is useful in this scenario where a client is consuming multiple services. If a service is consolidated, decomposed or replaced, the client doesn't necessarily require updating. It can continue making requests to the gateway, and only the routing changes.

A gateway also lets you abstract backend services from the clients, allowing you to keep client calls simple while enabling changes in the backend services behind the gateway. Client calls can be routed to whatever service or services need to handle the expected client behavior, allowing you to add, split, and reorganize services behind the gateway without changing the client.

### Multiple instances of the same service

Elasticity is key in cloud computing. Services can be spun up to meet increasing demand or spun down when demand is low to save money. The complexity of registering and unregistering service instances is encapsulated in the gateway. The client is unaware of an increase or decreases in the number of services.

Service instances can be deployed in a single or multiple regions. The Geode pattern details how a multi-region, active-active deployment can improve latency and increase availability of a service.

### Multiple versions of the same service

This pattern can be used for deployments, by allowing you to manage how updates are rolled out to users. When a new version of your service is deployed, it can be deployed in parallel with the existing version. Routing lets you control what version of the service is presented to the clients, giving you the flexibility to use various release strategies, whether incremental, parallel, or complete rollouts of updates. Any issues discovered after the new service is deployed can be quickly reverted by making a configuration change at the gateway, without affecting clients.

## Issues and considerations

- The gateway service can introduce a single point of failure. Ensure it's properly designed to meet your availability requirements. Consider resiliency and fault tolerance capabilities in the implementation.
- The gateway service can introduce a bottleneck. Ensure the gateway has adequate performance to handle load and can easily scale in line with your growth expectations.
- Perform load testing against the gateway to ensure you don't introduce cascading failures for services.
- Gateway routing is level 7. It can be based on IP, port, header, or URL.
- Gateway services can be global or regional. Azure Front Door is a global gateway, while Azure Application Gateway is regional. Use a global gateway if your solution requires multi-region deployments of services. Consider using Application Gateway if you have a regional workload that requires granular control how traffic is balanced. For example, you want to balance traffic between virtual machines.
- The gateway service is the public endpoint for services it sits in front of. Consider limiting public network access to the backend services, by making the services only accessible via the gateway or via a private virtual network.

## When to use this pattern

Use this pattern when:

- A client needs to consume multiple services that can be accessed behind a gateway.
- You want to simplify client applications by using a single endpoint.
- You need to route requests from externally addressable endpoints to internal virtual endpoints, such as exposing ports on a VM to cluster virtual IP addresses.
- A client needs to consume services running in multiple regions for latency or availability benefits.
- You want to route traffic away from regions with high carbon emissions when lower-emission alternatives exist in the backend pool.
- A client needs to consume a variable number of service instances.
- You want to implement a deployment strategy where clients access multiple versions of the service at the same time.

This pattern might not be suitable when you have a simple application that uses only one or two services.

## Workload design

An architect should evaluate how the Gateway Routing pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and to ensure that itrecoversto a fully functioning state after a failure occurs. | Gateway routing enables you to route traffic to only healthy nodes in your system. - RE:05 Redundancy - RE:10 Health monitoring |
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | Gateway routing enables you to decouple requests from backends, which in turn enables your backends to support advanced deployment models, platform transitions, and a single point of management for domain name resolution and encryption in transit. - OE:04 Tools and processes - OE:11 Safe deployment practices |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, code. | Gateway routing enables you to distribute traffic across nodes in your system to balance load. - PE:05 Scaling and partitioning |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

Using Nginx as the router, the following example shows a configuration file for a server that routes requests for applications residing on different virtual directories to different machines at the back end.

```
server {
 listen 80;
 server_name domain.com;
 location /app1 {
 proxy_pass http://10.0.3.10:80;
 }
 location /app2 {
 proxy_pass http://10.0.3.20:80;
 }
 location /app3 {
 proxy_pass http://10.0.3.30:80;
 }
}
```
The following Azure services can be used to implement the gateway routing pattern:

- An Application Gateway instance, which provides regional layer-7 routing.
- An Azure Front Door instance, which provides global layer-7 routing.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Deploy a collection of backend services into a set of **ge**ographical n**ode**s, each of which can service any request from any client in any region. This *active-active* approach distributes request processing around the globe to improve latency and increase availability.

## Context and problem

Many large-scale services have specific challenges around geo-availability and scale. Classic designs often *bring the data to the compute* by storing data in a remote SQL server that serves as the compute tier for that data, relying on scale-up for growth.

The classic approach might cause the following problems:

- Network latency issues for users coming from the other side of the globe to connect to the hosting endpoint
- Traffic management for demand bursts that can overwhelm the services in a single region
- Cost-prohibitive complexity of deploying copies of app infrastructure into multiple regions for a 24x7 service

Modern cloud infrastructure has evolved to enable geographic load balancing of front-end services, while allowing for geographic replication of backend services. For availability and performance, getting data closer to the user is good. When data is geo-distributed across a far-flung user base, the geo-distributed datastores should also be colocated with the compute resources that process the data. The geode pattern *brings the compute to the data*.

## Solution

Deploy the service into multiple satellite deployments spread around the globe, each of which is called a *geode*. The geode pattern harnesses key features of Azure to route traffic via the shortest path to a nearby geode, which improves latency and performance. Each geode is behind a global load balancer, and uses a geo-replicated read-write service like Azure Cosmos DB to host the data plane, ensuring cross-geode data consistency. Data replication services ensure that data stores are identical across geodes, so *all* requests can be served from *all* geodes.

The key difference between a deployment stamp and a geode is that geodes never exist in isolation. There should always be more than one geode in a production platform.

Geodes have the following characteristics:

- Consist of a collection of disparate types of resources, often defined in a template.
- Have no dependencies outside of the geode footprint and are self-contained. No geode is dependent on another to operate, and if one dies, the others continue to operate.
- Are loosely coupled via an edge network and replication backplane. For example, you can use Azure Traffic Manager or Azure Front Door for fronting the geodes, while Azure Cosmos DB can act as the replication backplane. Geodes aren't the same as clusters because they share a replication backplane, so the platform takes care of quorum issues.

The geode pattern occurs in big data architectures that use commodity hardware to process data colocated on the same machine, and MapReduce to consolidate results across machines. Another usage is near-edge compute, which brings compute closer to the intelligent edge of the network to reduce response time.

Services can use this pattern over dozens or hundreds of geodes. Furthermore, the resiliency of the whole solution increases with each added geode, since any geodes can take over if a regional outage takes one or more geodes offline.

It's also possible to augment local availability techniques, such as availability zones or paired regions, with the geode pattern for global availability. This increases complexity, but is useful if your architecture is underpinned by a storage engine such as blob storage that can only replicate to a paired region. You can deploy geodes into a zonal (single zone), multi-zone, or regional footprint, with a mind to regulatory or latency constraints on location.

## Issues and considerations

Use the following techniques and technologies to implement this pattern:

- Modern DevOps practices and tools to produce and rapidly deploy identical geodes across a large number of regions or instances.
- Autoscaling to scale out compute and database throughput instances within a geode. Each geode individually scales out, within the common backplane constraints.
- A front-end service like Azure Front Door that does dynamic content acceleration, routing through an optimal point of presence, and Split TCP.
- A replicating data store like Azure Cosmos DB to control data consistency.
- Serverless technologies where possible, to reduce always-on deployment cost, especially when load is frequently rebalanced around the globe. This strategy allows for many geodes to be deployed with minimal additional investment. Serverless and consumption-based billing technologies reduce waste and cost from duplicate geo-distributed deployments.
- API Management isn't required to implement the design pattern, but can be added to each geode that fronts the region's Azure Function App to provide a more robust API layer, enabling the implementation of additional functionality like rate limiting, for instance.

Consider the following points when deciding how to implement this pattern:

- Choose whether to process data locally in each region, or to distribute aggregations in a single geode and replicate the result across the globe. The Azure Cosmos DB change feed processor offers this granular control using its *lease container*concept, and the*leasecollectionprefix*in the corresponding Azure Functions binding. Each approach has distinct advantages and drawbacks.
- Geodes can work in tandem, using the Azure Cosmos DB change feed and a real-time communication platform like SignalR. Geodes can communicate with remote users via other geodes in a mesh pattern, without knowing or caring where the remote user is located.
- This design pattern implicitly decouples everything, resulting in an ultra-highly distributed and decoupled architecture. Consider how to track different components of the same request as they might execute asynchronously on different instances. A proper monitoring strategy is crucial. Both Azure Front Door and Azure Cosmos DB can be easily integrated with Log Analytics and Azure Functions should be deployed alongside Application Insights to provide a robust monitoring system at each component in the architecture.
- Distributed deployments have a greater number of secrets and ingress points that require property security measures. Key Vault provides a secure layer for secret management and each layer within the API architecture should be properly secured so that the only ingress point for the API is the front-end service like Azure Front Door. The Azure Cosmos DB should restrict traffic to the Azure Function Apps, and the Function apps to Azure Front Door using Microsoft Entra ID or practices like IP restriction.
- Performance is drastically affected by the number of geodes that are deployed and the specific App Service Plans applied to the API technology in each geode. Deployment of additional geodes or movement towards premium tiers come with increased costs for the additional memory and compute, but do not do so on a per transaction basis. Consider load testing the API architecture once deployed and contrast increasing the numbers of geodes with increasing the pricing tier so that the most cost-efficient model is used for your needs.
- Determine the availability requirements for your data. Azure Cosmos DB has optional flags for enabling multi-region write, availability zones, and more. These increase the availability for the Azure Cosmos DB instance and creates a more resilient data layer, but come with additional costs.
- Azure offers various load balancers that provide different functionalities for distribution of traffic. Use the decision tree to help select the right option for your API's front end.

## When to use this pattern

Use this pattern:

- To implement a high-scale platform that has users distributed over a wide area.
- For any service that requires extreme availability and resilience characteristics, because services based on the geode pattern can survive the loss of multiple service regions at the same time.

This pattern might not be suitable for

- Architectures that have constraints so that all geodes can't be equal for data storage. For example, there might be data residency requirements, an application that needs to maintain temporary state for a particular session, or a heavy weighting of requests towards a single region. In this case, consider using deployment stamps in combination with a global routing plane that is aware of where a user's data sits, such as the traffic routing component described within the deployment stamps pattern.
- Situations where there's no geographical distribution required. Instead, consider availability zones and paired regions for clustering.
- Situations where a legacy platform needs to be retrofitted. This pattern works for cloud-native development only, and can be difficult to retrofit.
- Simple architectures and requirements, where geo-redundancy and geo-distribution aren't required or advantageous.

## Workload design

An architect should evaluate how the Geode pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and to ensure that itrecoversto a fully functioning state after a failure occurs. | This pattern uses data replication to support the ideal that any client can connect to any geographical instance and by doing so it can help your workload withstand one or more regional outages. - RE:05 High-availability multi-region design - RE:05 Regions and availability zones |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, code. | You can use this pattern to serve your application from a region that's closest to your distributed user base. Doing so reduces latency by eliminating long-distance traffic and because you share infrastructure only among users that are currently using the same geode. - PE:03 Selecting services |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Examples

- Windows Active Directory implements an early variant of this pattern. Multi-primary replication means all updates and requests can in theory be served from all serviceable nodes, but Flexible Single Master Operation (FSMO) roles mean that all geodes aren't equal.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Implement functional checks in an application that external tools can access at regular intervals through exposed endpoints. This approach helps you verify that applications and services are performing correctly.

## Context and problem

It's a good practice to monitor web applications and back-end services. Monitoring helps ensure that applications and services are available and performing correctly. Business requirements often include monitoring.

It's sometimes more difficult to monitor cloud services than on-premises services. One reason is that you don't have full control of the hosting environment. Another is that the services typically depend on other services that platform vendors and others provide.

Many factors affect cloud-hosted applications. Examples include network latency, the performance and availability of the underlying compute and storage systems, and the network bandwidth between them. A service can fail entirely or partially due to any of these factors. To ensure a required level of availability, you must verify at regular intervals that your service performs correctly. Your service level agreement (SLA) might specify the level that you need to meet.

## Solution

Implement health monitoring by sending requests to an endpoint on your application. The application should perform the necessary checks and then return an indication of its status.

A health monitoring check typically combines two factors:

- The checks (if any) that the application or service performs in response to the request to the health verification endpoint
- The analysis of the results by the tool or framework that performs the health verification check

The response code indicates the status of the application. Optionally, the response code also provides the status of components and services that the app uses. The monitoring tool or framework performs the latency or response time check.

The following figure provides an overview of the pattern.

The health monitoring code in the application might also run other checks to determine:

- The availability and response time of cloud storage or a database.
- The status of other resources or services that the application uses. These resources and services might be in the application or outside it.

Services and tools are available that monitor web applications by submitting a request to a configurable set of endpoints. These services and tools then evaluate the results against a set of configurable rules. It's relatively easy to create a service endpoint for the sole purpose of performing some functional tests on a system.

Typical checks that monitoring tools perform include:

- Validating the response code. For example, an HTTP response of 200 (OK) indicates that the application responded without error. The monitoring system might also check for other response codes to give more comprehensive results.
- Checking the content of the response to detect errors, even when the status code is 200 (OK). By checking the content, you can detect errors that affect only a section of the returned web page or service response. For example, you might check the title of a page or look for a specific phrase that indicates that the app returned the correct page.
- Measuring the response time. The value includes the network latency and the time that the application took to issue the request. An increasing value can indicate an emerging problem with the application or network.
- Checking resources or services that are located outside the application. An example is a content delivery network that the application uses to deliver content from global caches.
- Checking for the expiration of TLS certificates.
- Measuring the response time of a DNS lookup for the URL of the application. This check measures DNS latency and DNS failures.
- Validating the URL that a DNS lookup returns. By validating, you can ensure that entries are correct. You can also help prevent malicious request redirection that might result after an attack on your DNS server.

Where possible, it's also useful to run these checks from different on-premises or hosted locations and then compare response times. Ideally, you should monitor applications from locations that are close to customers. Then you get an accurate view of the performance from each location. This practice provides a more robust checking mechanism. The results can also help you make the following decisions:

- Where to deploy your application
- Whether to deploy it in more than one datacenter

To ensure that your application works correctly for all customers, run tests against all the service instances that customers use. For example, if customer storage is spread across more than one storage account, the monitoring process should check each account.

## Issues and considerations

Consider the following points when you decide how to implement this pattern:

- Think about how to validate the response. For example, determine whether a 200 (OK) status code is sufficient to verify that the application is working correctly. Checking the status code is the minimum implementation of this pattern. A status code provides a basic measure of application availability. But a code supplies little information about the operations, trends, and possible upcoming issues in the application.
- Determine the number of endpoints to expose for an application. One approach is to expose at least one endpoint for the core services that the application uses and another for lower-priority services. With this approach, you can assign different levels of importance to each monitoring result. Also consider exposing extra endpoints. You can expose one for each core service to increase monitoring granularity. For example, a health verification check might check the database, the storage, and an external geocoding service that an application uses. Each might require a different level of uptime and response time. The geocoding service or some other background task might be unavailable for a few minutes. But the application might still be healthy.
- Decide whether to use the same endpoint for monitoring and for general access. You can use the same endpoint for both but design a specific path for health verification checks. For example, you can use - */health*on the general access endpoint. With this approach, monitoring tools can run some functional tests in the application. Examples include registering a new user, signing in, and placing a test order. At the same time, you can also verify that the general access endpoint is available.
- Determine the type of information to collect in the service in response to monitoring requests. You also need to determine how to return this information. Most existing tools and frameworks look only at the HTTP status code that the endpoint returns. To return and validate additional information, you might have to create a custom monitoring utility or service.
- Figure out how much information to collect. Performing excessive processing during the check can overload the application and affect other users. The processing time might also exceed the timeout of the monitoring system. As a result, the system might mark the application as unavailable. Most applications include instrumentation such as error handlers and performance counters. These tools can log performance and detailed error information, which might be sufficient. Consider using this data instead of returning additional information from a health verification check.
- Consider caching the endpoint status. Running the health check frequently might be expensive. For example, if the health status is reported through a dashboard, you don't want every request to the dashboard to trigger a health check. Instead, periodically check the system health, and cache the status. Expose an endpoint that returns the cached status.
- Plan how to configure security for the monitoring endpoints. By configuring security, you can help protect the endpoints from public access, which might: - Expose the application to malicious attacks.
- Risk the exposure of sensitive information.
- Attract denial of service (DoS) attacks.
 - Typically, you configure security in the application configuration. Then you can update the settings easily without restarting the application. Consider using one or more of the following techniques: - Secure the endpoint by requiring authentication. If the monitoring service or tool supports authentication, you can use an authentication security key in the request header. You can also pass credentials with the request. When you use authentication, consider how to access your health check endpoints. As an example, Azure App Service has a built-in health check that integrates with App Service authentication and authorization features.
- Use an obscure or hidden endpoint. For example, expose the endpoint on a different IP address than the one that the default application URL uses. Configure the endpoint on a nonstandard HTTP port. Also, consider using a complex path to your test page. You can usually specify extra endpoint addresses and ports in the application configuration. If necessary, you can add entries for these endpoints to the DNS server. Then you avoid having to specify the IP address directly.
- Expose a method on an endpoint that accepts a parameter such as a key value or an operation mode value. When a request arrives, the code can run specific tests that depend on the value of the parameter. The code can return a 404 (Not Found) error if it doesn't recognize the parameter value. Make it possible to define parameter values in the application configuration.
- Use a separate endpoint that performs basic functional tests without compromising the operation of the application. With this approach, you can help reduce the impact of a DoS attack. Ideally, avoid using a test that might expose sensitive information. Sometimes you must return information that might be useful to an attacker. In this case, consider how to protect the endpoint and the data from unauthorized access. Relying on obscurity isn't enough. Consider also using an HTTPS connection and encrypting sensitive data, although this approach increases the load on the server.

- Decide how to ensure that the monitoring agent is performing correctly. One approach is to expose an endpoint that returns a value from the application configuration or a random value that you can use to test the agent. Also ensure that the monitoring system performs checks on itself. You can use a self-test or built-in test to prevent the monitoring system from issuing false positive results.

## When to use this pattern

This pattern is useful for:

- Monitoring websites and web applications to verify availability.
- Monitoring websites and web applications to check for correct operation.
- Monitoring middle-tier or shared services to detect and isolate failures that can disrupt other applications.
- Complementing existing instrumentation in the application, such as performance counters and error handlers. Health verification checking doesn't replace application requirements for logging and auditing. Instrumentation can provide valuable information for an existing framework that monitors counters and error logs to detect failures or other issues. But instrumentation can't provide information if an application is unavailable.

## Workload design

An architect should evaluate how the Health Endpoint Monitoring pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and to ensure that itrecoversto a fully functioning state after a failure occurs. | These endpoints support a workload's reliability alerting and dashboarding efforts. They can also be used it as a signal for self-healing remediation. - RE:07 Self-healing and self-preservation - RE:10 Monitoring and alerting strategy |
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | Standardizing which health endpoints to expose, and the level of detail in the results, across your workload will help you triage issues. - OE:07 Monitoring system |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, code. | Health endpoints improve load balancing logic by routing traffic to only nodes that are verified as healthy. With additional configuration, you can also get metrics on available node capacity. - PE:05 Scaling and partitioning |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

You can use the ASP.NET Core health checks middleware and libraries to report the health of app infrastructure components. This framework provides a way to report health checks in a consistent way. It implements many of the practices that this article describes. For instance, the ASP.NET Core health checks include external checks like database connectivity and specific concepts like liveness and readiness probes.

Several example implementations that use ASP.NET Core health checks are available on GitHub.

## Monitor endpoints in Azure-hosted applications

Options for monitoring endpoints in Azure applications include:

- Use the built-in monitoring features of Azure, such as Azure Monitor.
- Use a third-party service or a framework like Microsoft System Center Operations Manager.
- Create a custom utility or a service that runs on your own server or a hosted server.

Even though Azure provides comprehensive monitoring options, you can use additional services and tools to provide extra information. Application Insights, a feature of Monitor, is designed for development teams. This feature helps you understand how your app performs and how it's used. Application Insights monitors request rates, response times, failure rates, and dependency rates. It can help you determine whether external services are slowing you down.

The conditions that you can monitor depend on the hosting mechanism that you choose for your application. All options in this section support alert rules. An alert rule uses a web endpoint that you specify in the settings for your service. This endpoint should respond in a timely way so that the alert system can detect that the application is operating correctly. For more information, see Create a new alert rule.

If there's a major outage, client traffic should be routable to an application deployment that's available across other regions or zones. This situation is a good case for cross-premises connectivity and global load balancing. The choice depends on whether the application is internal or external facing. Services such as Azure Front Door, Azure Traffic Manager, or content delivery networks can route traffic across regions based on data that health probes provide.

Traffic Manager is a routing and load-balancing service. It can use a range of rules and settings to distribute requests to specific instances of your application. Besides routing requests, Traffic Manager can regularly ping a URL, port, and relative path. You specify the ping targets with the goal of determining which instances of your application are active and responding to requests. If Traffic Manager detects a status code of 200 (OK), it marks the application as available. Any other status code causes Traffic Manager to mark the application as offline. The Traffic Manager console displays the status of each application. You can configure each rule to reroute requests to other instances of the application that are responding.

Traffic Manager waits for a certain amount of time to receive a response from the monitoring URL. Make sure that your health verification code runs in this time. Allow for network latency for the round trip from Traffic Manager to your application and back again.

## Next steps

The following guidance is useful for implementing this pattern:

- Health monitoring guidance in microservices-based applications
- Monitoring application health for reliability, part of the Azure Well-Architected Framework
- Create a new alert rule

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Create indexes over the fields in data stores that are frequently referenced by queries. This pattern can improve query performance by allowing applications to more quickly locate the data to retrieve from a data store.

## Context and problem

Many data stores organize the data for a collection of entities using the primary key. An application can use this key to locate and retrieve data. The figure shows an example of a data store holding customer information. The primary key is the Customer ID. The figure shows customer information organized by the primary key (Customer ID).

While the primary key is valuable for queries that fetch data based on the value of this key, an application might not be able to use the primary key if it needs to retrieve data based on some other field. In the customers example, an application can't use the Customer ID primary key to retrieve customers if it queries data solely by referencing the value of some other attribute, such as the town in which the customer is located. To perform a query such as this, the application might have to fetch and examine every customer record, which could be a slow process.

Many relational database management systems support secondary indexes. A secondary index is a separate data structure that's organized by one or more nonprimary (secondary) key fields, and it indicates where the data for each indexed value is stored. The items in a secondary index are typically sorted by the value of the secondary keys to enable fast lookup of data. These indexes are usually maintained automatically by the database management system.

You can create as many secondary indexes as you need to support the different queries that your application performs. For example, in a Customers table in a relational database where the Customer ID is the primary key, it's beneficial to add a secondary index over the town field if the application frequently looks up customers by the town where they reside.

However, although secondary indexes are common in relational systems, some NoSQL data stores used by cloud applications don't provide an equivalent feature.

## Solution

If the data store doesn't support secondary indexes, you can emulate them manually by creating your own index tables. An index table organizes the data by a specified key. Three strategies are commonly used for structuring an index table, depending on the number of secondary indexes that are required and the nature of the queries that an application performs.

The first strategy is to duplicate the data in each index table but organize it by different keys (complete denormalization). The next figure shows index tables that organize the same customer information by Town and LastName.

This strategy is appropriate if the data is relatively static compared to the number of times it's queried using each key. If the data is more dynamic, the processing overhead of maintaining each index table becomes too large for this approach to be useful. Also, if the data volume is high, the amount of space required to store duplicate data is significant.

The second strategy is to create normalized index tables organized by different keys and reference the original data by using the primary key rather than duplicating it, as shown in the following figure. The original data is called a fact table.

This technique saves space and reduces the overhead of maintaining duplicate data. The disadvantage is that an application has to perform two lookup operations to find data using a secondary key. It has to find the primary key for the data in the index table, and then use the primary key to look up the data in the fact table.

The third strategy is to create partially normalized index tables organized by different keys that duplicate frequently retrieved fields. Reference the fact table to access less frequently accessed fields. The next figure shows how commonly accessed data is duplicated in each index table.

With this strategy, you can strike a balance between the first two approaches. The data for common queries can be retrieved quickly by using a single lookup, while the space and maintenance overhead isn't as significant as duplicating the entire data set.

If an application frequently queries data by specifying a combination of values (for example, "Find all customers that live in Redmond and that have a last name of Smith"), you could implement the keys to the items in the index table as a concatenation of the Town attribute and the LastName attribute. The next figure shows an index table based on composite keys. The keys are sorted by Town, and then by LastName for records that have the same value for Town.

Index tables can speed up query operations over sharded data, and are especially useful where the shard key is hashed. The next figure shows an example where the shard key is a hash of the Customer ID. The index table can organize data by the nonhashed value (Town and LastName), and provide the hashed shard key as the lookup data. This can save the application from repeatedly calculating hash keys (an expensive operation) if it needs to retrieve data that falls within a range, or it needs to fetch data in order of the nonhashed key. For example, a query such as "Find all customers that live in Redmond" can be quickly resolved by locating the matching items in the index table, where they're all stored in a contiguous block. Then, follow the references to the customer data using the shard keys stored in the index table.

## Issues and considerations

Consider the following points when deciding how to implement this pattern:

- The overhead of maintaining secondary indexes can be significant. You must analyze and understand the queries that your application uses. Only create index tables when they're likely to be used regularly. Don't create speculative index tables to support queries that an application doesn't perform, or performs only occasionally.
- Duplicating data in an index table can add significant overhead in storage costs and the effort required to maintain multiple copies of data.
- Implementing an index table as a normalized structure that references the original data requires an application to perform two lookup operations to find data. The first operation searches the index table to retrieve the primary key, and the second uses the primary key to fetch the data.
- If a system incorporates a number of index tables over large datasets, it can be difficult to maintain consistency between index tables and the original data. It might be possible to design the application around the eventual consistency model. For example, to insert, update, or delete data, an application could post a message to a queue and let a separate task perform the operation and maintain the index tables that reference this data asynchronously. For more information about implementing eventual consistency, see the Data Consistency Primer. - Tip - Microsoft Azure storage tables support transactional updates for changes made to data held in the same partition (referred to as entity group transactions). If you can store the data for a fact table and one or more index tables in the same partition, you can use this feature to help ensure consistency.
- Index tables might themselves be partitioned or sharded.

## When to use this pattern

Use this pattern to improve query performance when an application frequently needs to retrieve data by using a key other than the primary (or shard) key.

This pattern might not be useful when:

- Data is volatile. An index table can become out of date very quickly, making it ineffective or making the overhead of maintaining the index table greater than any savings made by using it.
- A field selected as the secondary key for an index table is nondiscriminating and can only have a small set of values (for example, gender).
- The balance of the data values for a field selected as the secondary key for an index table are highly skewed. For example, if 90% of the records contain the same value in a field, then creating and maintaining an index table to look up data based on this field might create more overhead than scanning sequentially through the data. However, if queries very frequently target values that lie in the remaining 10%, this index can be useful. You should understand the queries that your application is performing, and how frequently they're performed.

## Workload design

An architect should evaluate how the Index Table pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and to ensure that itrecoversto a fully functioning state after a failure occurs. | Because clients are pointed to their shard, partition, or endpoint through a lookup process, you can use this pattern to facilitate a failover approach for data access. - RE:06 Data partitioning - RE:09 Disaster recovery |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, code. | Clients are pointed to their shard, partition, or endpoint, which can enable dynamic data partitioning for performance optimization. - PE:05 Scaling and partitioning - PE:08 Data performance |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

Azure storage tables provide a highly scalable key/value data store for applications that run in the cloud. Applications store and retrieve data values by specifying a key. These data values can contain multiple fields, but the structure of a data item is opaque to table storage, which treats each item as an array of bytes.

Azure storage tables also support sharding. The sharding key includes two elements, a partition key and a row key. Items that have the same partition key are stored in the same partition (shard), and the items are stored in row key order within a shard. Table storage is optimized for performing queries that fetch data falling within a contiguous range of row key values within a partition. If you're building cloud applications that store information in Azure tables, you should structure your data with this feature in mind.

For example, consider an application that stores information about movies. The application frequently queries movies by genre (such as action, documentary, historical, comedy, and drama). You could create an Azure table with partitions for each genre by using the genre as the partition key, and specifying the movie name as the row key, as shown in the next figure.

This approach is less effective if the application also needs to query movies by starring actor. In this case, you can create a separate Azure table that acts as an index table. The partition key is the actor and the row key is the movie name. The data for each actor will be stored in separate partitions. If a movie stars more than one actor, the same movie will occur in multiple partitions.

You can duplicate the movie data in the values held by each partition by adopting the first approach described in the Solution section above. However, it's likely that each movie will be replicated several times (once for each actor), so it might be more efficient to partially denormalize the data to support the most common queries (such as the names of the other actors) and enable an application to retrieve any remaining details by including the partition key necessary to find the complete information in the genre partitions. This approach is described by the third option in the Solution section. The next figure shows this approach.

## Next steps

- Data Consistency Primer. An index table must be maintained as the data that it indexes changes. In the cloud, it might not be possible or appropriate to perform operations that update an index as part of the same transaction that modifies the data. In that case, an eventually consistent approach is more suitable. Provides information on the issues surrounding eventual consistency.

## Related resources

The following patterns might also be relevant when implementing this pattern:

- Sharding pattern. The Index Table pattern is frequently used in conjunction with data partitioned by using shards. The Sharding pattern provides more information on how to divide a data store into a set of shards.
- Materialized View pattern. Instead of indexing data to support queries that summarize data, it might be more appropriate to create a materialized view of the data. Describes how to support efficient summary queries by generating prepopulated views over data.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Coordinate the actions performed by a collection of collaborating instances in a distributed application by electing one instance as the leader that assumes responsibility for managing the others. This can help to ensure that instances don't conflict with each other, cause contention for shared resources, or inadvertently interfere with the work that other instances are performing.

## Context and problem

A typical cloud application has many tasks acting in a coordinated manner. These tasks could all be instances running the same code and requiring access to the same resources, or they might be working together in parallel to perform the individual parts of a complex calculation.

The task instances might run separately for much of the time, but it might also be necessary to coordinate the actions of each instance to ensure that they don't conflict, cause contention for shared resources, or accidentally interfere with the work that other task instances are performing.

For example:

- In a cloud-based system that implements horizontal scaling, multiple instances of the same task could be running at the same time with each instance serving a different user. If these instances write to a shared resource, it's necessary to coordinate their actions to prevent each instance from overwriting the changes made by the others.
- If the tasks are performing individual elements of a complex calculation in parallel, the results need to be aggregated when they all complete.

The task instances are all peers, so there isn't a natural leader that can act as the coordinator or aggregator.

## Solution

A single task instance should be elected to act as the leader, and this instance should coordinate the actions of the other subordinate task instances. If all of the task instances are running the same code, they are each capable of acting as the leader. Therefore, the election process must be managed carefully to prevent two or more instances taking over the leader position at the same time.

The system must provide a robust mechanism for selecting the leader. This method has to cope with events such as network outages or process failures. In many solutions, the subordinate task instances monitor the leader through some type of heartbeat method, or by polling. If the designated leader terminates unexpectedly, or a network failure makes the leader unavailable to the subordinate task instances, it's necessary for them to elect a new leader.

There are multiple strategies for electing a leader among a set of tasks in a distributed environment, including:

- Racing to acquire a shared, distributed mutex. The first task instance that acquires the mutex is the leader. However, the system must ensure that, if the leader terminates or becomes disconnected from the rest of the system, the mutex is released to allow another task instance to become the leader. This strategy is demonstrated in the following example.
- Implementing one of the common leader election algorithms such as the Bully Algorithm, the Raft Consensus Algorithm, or the Chang and Roberts algorithm. These algorithms assume that each candidate in the election has a unique ID, and that it can communicate with the other candidates reliably.

## Issues and considerations

Consider the following points when deciding how to implement this pattern:

- The process of electing a leader should be resilient to transient and persistent failures.
- It must be possible to detect when the leader has failed or has become otherwise unavailable (such as due to a communications failure). How quickly detection is needed is system dependent. Some systems might be able to function for a short time without a leader, during which a transient fault might be fixed. In other cases, it might be necessary to detect leader failure immediately and trigger a new election.
- In a system that implements horizontal autoscaling, the leader could be terminated if the system scales back and shuts down some of the computing resources.
- Using a shared, distributed mutex introduces a dependency on the external service that provides the mutex. The service constitutes a single point of failure. If it becomes unavailable for any reason, the system won't be able to elect a leader.
- Using a single dedicated process as the leader is a straightforward approach. However, if the process fails there could be a significant delay while it's restarted. The resulting latency can affect the performance and response times of other processes if they're waiting for the leader to coordinate an operation.
- Implementing one of the leader election algorithms manually provides the greatest flexibility for tuning and optimizing the code.
- Avoid making the leader a bottleneck in the system. The purpose of the leader is to coordinate the work of the subordinate tasks, and it doesn't necessarily have to participate in this work itself—although it should be able to do so if the task isn't elected as the leader.

## When to use this pattern

Use this pattern when the tasks in a distributed application, such as a cloud-hosted solution, need careful coordination and there's no natural leader.

This pattern might not be useful if:

- There's a natural leader or dedicated process that can always act as the leader. For example, it might be possible to implement a singleton process that coordinates the task instances. If this process fails or becomes unhealthy, the system can shut it down and restart it.
- The coordination between tasks can be achieved by using a more lightweight method. For example, if several task instances simply need coordinated access to a shared resource, a better solution is to use optimistic or pessimistic locking to control access.
- A third-party solution, like Apache Zookeeper might be a more efficient solution.

## Workload design

An architect should evaluate how the Leader Election pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and to ensure that itrecoversto a fully functioning state after a failure occurs. | This pattern mitigates the effect of node malfunctions by reliably redirecting work. It also implements failover via consensus algorithms when a leader malfunctions. - RE:05 Redundancy - RE:07 Self-healing |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

The Leader Election sample on GitHub shows how to use a lease on an Azure Storage blob to provide a mechanism for implementing a shared, distributed mutex. This mutex can be used to elect a leader among a group of available worker instances. The first instance to acquire the lease is elected the leader and remains the leader until it releases the lease or isn't able to renew the lease. Other worker instances can continue to monitor the blob lease in case the leader is no longer available.

A blob lease is an exclusive write lock over a blob. A single blob can be the subject of only one lease at any point in time. A worker instance can request a lease over a specified blob, and it'll be granted the lease if no other worker instance holds a lease over the same blob. Otherwise, the request will throw an exception.

To avoid a faulted leader instance retaining the lease indefinitely, specify a lifetime for the lease. When this expires, the lease becomes available. However, while an instance holds the lease it can request that the lease is renewed, and it'll be granted the lease for a further period of time. The leader instance can continually repeat this process if it wants to retain the lease. For more information on how to lease a blob, see Lease Blob (REST API).

The `BlobDistributedMutex` class in the C# example below contains the `RunTaskWhenMutexAcquired` method that enables a worker instance to attempt to acquire a lease over a specified blob. The details of the blob (the name, container, and storage account) are passed to the constructor in a `BlobSettings` object when the `BlobDistributedMutex` object is created (this object is a simple struct that is included in the sample code). The constructor also accepts a `Task` that references the code that the worker instance should run if it successfully acquires the lease over the blob and is elected the leader. The code that handles the low-level details of acquiring the lease is implemented in a separate helper class named `BlobLeaseManager`.

```
public class BlobDistributedMutex
{
 ...
 private readonly BlobSettings blobSettings;
 private readonly Func<CancellationToken, Task> taskToRunWhenLeaseAcquired;
 ...
 public BlobDistributedMutex(BlobSettings blobSettings,
 Func<CancellationToken, Task> taskToRunWhenLeaseAcquired, ... )
 {
 this.blobSettings = blobSettings;
 this.taskToRunWhenLeaseAcquired = taskToRunWhenLeaseAcquired;
 ...
 }
 public async Task RunTaskWhenMutexAcquired(CancellationToken token)
 {
 var leaseManager = new BlobLeaseManager(blobSettings);
 await this.RunTaskWhenBlobLeaseAcquired(leaseManager, token);
 }
 ...
```
The `RunTaskWhenMutexAcquired` method in the preceding code sample invokes the `RunTaskWhenBlobLeaseAcquired` method shown in the following code sample to actually acquire the lease. The `RunTaskWhenBlobLeaseAcquired` method runs asynchronously. If the lease is successfully acquired, the worker instance has been elected the leader. The purpose of the `taskToRunWhenLeaseAcquired` delegate is to perform the work that coordinates the other worker instances. If the lease isn't acquired, another worker instance has been elected as the leader and the current worker instance remains a subordinate. Note that the `TryAcquireLeaseOrWait` method is a helper method that uses the `BlobLeaseManager` object to acquire the lease.

```
 private async Task RunTaskWhenBlobLeaseAcquired(
 BlobLeaseManager leaseManager, CancellationToken token)
 {
 while (!token.IsCancellationRequested)
 {
 // Try to acquire the blob lease.
 // Otherwise wait for a short time before trying again.
 string? leaseId = await this.TryAcquireLeaseOrWait(leaseManager, token);
 if (!string.IsNullOrEmpty(leaseId))
 {
 // Create a new linked cancellation token source so that if either the
 // original token is canceled or the lease can't be renewed, the
 // leader task can be canceled.
 using (var leaseCts =
 CancellationTokenSource.CreateLinkedTokenSource(new[] { token }))
 {
 // Run the leader task.
 var leaderTask = this.taskToRunWhenLeaseAcquired.Invoke(leaseCts.Token);
 ...
 }
 }
 }
 ...
 }
```
The task started by the leader also runs asynchronously. While this task is running, the `RunTaskWhenBlobLeaseAcquired` method shown in the following code sample periodically attempts to renew the lease. This helps to ensure that the worker instance remains the leader. In the sample solution, the delay between renewal requests is less than the time specified for the duration of the lease in order to prevent another worker instance from being elected the leader. If the renewal fails for any reason, the leader-specific task is canceled.

If the lease fails to be renewed or the task is canceled (possibly as a result of the worker instance shutting down), the lease is released. At this point, this or another worker instance can be elected as the leader. The following code extract shows this part of the process.

```
 private async Task RunTaskWhenBlobLeaseAcquired(
 BlobLeaseManager leaseManager, CancellationToken token)
 {
 while (...)
 {
 ...
 if (...)
 {
 ...
 using (var leaseCts = ...)
 {
 ...
 // Keep renewing the lease in regular intervals.
 // If the lease can't be renewed, then the task completes.
 var renewLeaseTask =
 this.KeepRenewingLease(leaseManager, leaseId, leaseCts.Token);
 // When any task completes (either the leader task itself or when it
 // couldn't renew the lease) then cancel the other task.
 await CancelAllWhenAnyCompletes(leaderTask, renewLeaseTask, leaseCts);
 }
 }
 }
 }
 ...
}
```
The `KeepRenewingLease` method is another helper method that uses the `BlobLeaseManager` object to renew the lease. The `CancelAllWhenAnyCompletes` method cancels the tasks specified as the first two parameters. The following diagram illustrates using the `BlobDistributedMutex` class to elect a leader and run a task that coordinates operations.

The following code example shows how to use the `BlobDistributedMutex` class within a worker instance. This code acquires a lease over a blob named `MyLeaderCoordinatorTask` in the lease's container Azure Blob Storage, and specifies that the code defined in the `MyLeaderCoordinatorTask` method should run if the worker instance is elected the leader.

```
// Create a BlobSettings object with the connection string or managed identity and the name of the blob to use for the lease
BlobSettings blobSettings = new BlobSettings(storageConnStr, "leases", "MyLeaderCoordinatorTask");
// Create a new BlobDistributedMutex object with the BlobSettings object and a task to run when the lease is acquired
var distributedMutex = new BlobDistributedMutex(
 blobSettings, MyLeaderCoordinatorTask);
// Wait for completion of the DistributedMutex and the UI task before exiting
await distributedMutex.RunTaskWhenMutexAcquired(cancellationToken);
...
// Method that runs if the worker instance is elected the leader
private static async Task MyLeaderCoordinatorTask(CancellationToken token)
{
 ...
}
```
Note the following points about the sample solution:

- The blob is a potential single point of failure. If the blob service becomes unavailable, or is inaccessible, the leader won't be able to renew the lease and no other worker instance will be able to acquire the lease. In this case, no worker instance will be able to act as the leader. However, the blob service is designed to be resilient, so complete failure of the blob service is considered to be extremely unlikely.
- If the task being performed by the leader stalls, the leader might continue to renew the lease, preventing any other worker instance from acquiring the lease and taking over the leader position in order to coordinate tasks. In the real world, the health of the leader should be checked at frequent intervals.
- The election process is nondeterministic. You can't make any assumptions about which worker instance will acquire the blob lease and become the leader.
- The blob used as the target of the blob lease shouldn't be used for any other purpose. If a worker instance attempts to store data in this blob, this data won't be accessible unless the worker instance is the leader and holds the blob lease.

## Next steps

The following guidance might also be relevant when implementing this pattern:

- This pattern has a downloadable sample application.
- Autoscaling Guidance. It's possible to start and stop instances of the task hosts as the load on the application varies. Autoscaling can help to maintain throughput and performance during times of peak processing.
- The Task-based Asynchronous pattern.
- Apache Curator a client library for Apache ZooKeeper.
- The article Lease Blob (REST API) on MSDN.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Generate prepopulated views over the data in one or more data stores when the data isn't ideally formatted for required query operations. This can help support efficient querying and data extraction, and improve application performance.

## Context and problem

When storing data, the priority for developers and data administrators is often focused on how the data is stored, as opposed to how it's read. The chosen storage format is usually closely related to the format of the data, requirements for managing data size and data integrity, and the kind of store in use. For example, when you use NoSQL document store, the data is often represented as a series of aggregates, each containing all of the information for that entity.

However, this can have a negative effect on queries. When a query only needs a subset of the data from some entities, such as a summary of orders for several customers without all of the order details, it must extract all of the data for the relevant entities in order to obtain the required information.

## Solution

To support efficient querying, a common solution is to generate, in advance, a view that materializes the data in a format suited to the required results set. The Materialized View pattern describes generating prepopulated views of data in environments where the source data isn't in a suitable format for querying, where generating a suitable query is difficult, or where query performance is poor due to the nature of the data or the data store.

These materialized views, which only contain data required by a query, allow applications to quickly obtain the information they need. In addition to joining tables or combining data entities, materialized views can include the current values of calculated columns or data items, the results of combining values or executing transformations on the data items, and values specified as part of the query. A materialized view can even be optimized for just a single query.

A key point is that a materialized view and the data it contains is completely disposable because it can be entirely rebuilt from the source data stores. A materialized view is never updated directly by an application, and so it's a specialized cache.

When the source data for the view changes, the view must be updated to include the new information. You can schedule this to happen automatically, or when the system detects a change to the original data. In some cases it might be necessary to regenerate the view manually. The figure shows an example of how the Materialized View pattern might be used.

## Issues and considerations

Consider the following points when deciding how to implement this pattern:

How and when the view will be updated. Ideally it'll regenerate in response to an event indicating a change to the source data, although this can lead to excessive overhead if the source data changes rapidly. Alternatively, consider using a scheduled task, an external trigger, or a manual action to regenerate the view.

In some systems, like when you use the Event Sourcing pattern to maintain a store of only the events that modified the data, materialized views are necessary. Prepopulating views by examining all events to determine the current state might be the only way to obtain information from the event store. If you're not using Event Sourcing, you need to consider whether a materialized view is helpful or not. Materialized views tend to be specifically tailored to one, or a small number of queries. If many queries are used, materialized views can result in unacceptable storage capacity requirements and storage cost.

Consider the impact on data consistency when generating the view, and when updating the view if this occurs on a schedule. If the source data is changing at the point when the view is generated, the copy of the data in the view won't be fully consistent with the original data.

Consider where you'll store the view. The view doesn't have to be located in the same store or partition as the original data. It can be a subset from a few different partitions combined.

A view can be rebuilt if lost. Because of that, if the view is transient and is only used to improve query performance by reflecting the current state of the data, or to improve scalability, it can be stored in a cache or in a less reliable location.

When defining a materialized view, maximize its value by adding data items or columns to it based on computation or transformation of existing data items, on values passed in the query, or on combinations of these values when appropriate.

Where the storage mechanism supports it, consider indexing the materialized view to further increase performance. Most relational databases support indexing for views, as do big data solutions based on Apache Hadoop.

## When to use this pattern

This pattern is useful when:

- Creating materialized views over data that's difficult to query directly, or where queries must be very complex to extract data that's stored in a normalized, semi-structured, or unstructured way.
- Creating temporary views that can dramatically improve query performance, or can act directly as source views or data transfer objects for the UI, for reporting, or for display.
- Supporting occasionally connected or disconnected scenarios where connection to the data store isn't always available. The view can be cached locally in this case.
- Simplifying queries and exposing data for experimentation in a way that doesn't require knowledge of the source data format. For example, by joining different tables in one or more databases, or one or more domains in NoSQL stores, and then formatting the data to fit its eventual use.
- Providing access to specific subsets of the source data that, for security or privacy reasons, shouldn't be generally accessible, open to modification, or fully exposed to users.
- Bridging different data stores, to take advantage of their individual capabilities. For example, using a cloud store that's efficient for writing as the reference data store, and a relational database that offers good query and read performance to hold the materialized views.
- When you use microservices, you are recommended to keep them loosely coupled, including their data storage. Therefore, materialized views can help you consolidate data from your services. If materialized views are not appropriate in your microservices architecture or specific scenario, please consider having well-defined boundaries that align to domain driven design (DDD) and aggregate their data when requested.

This pattern isn't useful in the following situations:

- The source data is simple and easy to query.
- The source data changes very quickly, or can be accessed without using a view. In these cases, you should avoid the processing overhead of creating views.
- Consistency is a high priority. The views might not always be fully consistent with the original data.

## Workload design

An architect should evaluate how the Materialized View pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, code. | The materialized views store the results of complex computations or queries without requiring the database engine or client to recompute for every request. This design reduces overall resource consumption. - PE:08 Data performance |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

The following figure shows an example of using the Materialized View pattern to generate a summary of sales. Data in the Order, OrderItem, and Customer tables in separate partitions in an Azure storage account are combined to generate a view containing the total sales value for each product in the Electronics category, along with a count of the number of customers who made purchases of each item.

Creating this materialized view requires complex queries. However, by exposing the query result as a materialized view, users can easily obtain the results and use them directly or incorporate them in another query. The view is likely to be used in a reporting system or dashboard, and can be updated on a scheduled basis such as weekly.

Although this example uses Azure table storage, many relational database management systems also provide native support for materialized views.

## Next steps

- Data Consistency Primer. The summary information in a materialized view has to be maintained so that it reflects the underlying data values. As the data values change, it might not be practical to update the summary data in real time, and instead you'll have to adopt an eventually consistent approach. Summarizes the issues surrounding maintaining consistency over distributed data, and describes the benefits and tradeoffs of different consistency models.

## Related resources

The following patterns might also be relevant when implementing this pattern:

- Command and Query Responsibility Segregation (CQRS) pattern. Use to update the information in a materialized view by responding to events that occur when the underlying data values change.
- Event Sourcing pattern. Use in conjunction with the CQRS pattern to maintain the information in a materialized view. When the data values a materialized view is based on are changed, the system can raise events that describe these changes and save them in an event store.
- Index Table pattern. The data in a materialized view is typically organized by a primary key, but queries might need to retrieve information from this view by examining data in other fields. Use to create secondary indexes over data sets for data stores that don't support native secondary indexes.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Connect messaging systems that are built on different messaging infrastructures and relay messages between them. This approach lets you integrate disparate systems without modifying the systems themselves.

## Context and problem

Many organizations and workloads can inadvertently have IT systems that use multiple messaging infrastructures like Microsoft Message Queueing (MSMQ), RabbitMQ, Azure Service Bus, and Amazon SQS. This problem can occur due to mergers, acquisitions, or due to extending current on-premises systems to cloud-hosted components for cost-effectiveness and the ease of maintenance.

Developers might address this challenge by modifying the systems being integrated to communicate by using HTTP-based web services. However, this approach has drawbacks, including:

- The systems must be modified by adding an HTTP client on one side and an HTTP request handler on the other. The systems must then be retested and redeployed.
- HTTP endpoints must be hosted, which adds complexity when you make web services secure and highly available.
- Frequent network connectivity problems that require custom-built retry mechanisms.

## Solution

If the systems being integrated consist of components that communicate by exchanging messages, the Messaging Bridge pattern improves integration and mitigates drawbacks.

In this scenario, each system connects to one messaging infrastructure. To integrate across different messaging infrastructures, introduce a bridge component that connects to two or more messaging infrastructures at the same time. The bridge pulls messages from one and pushes them to the other without changing the payload.

The systems being integrated don't need to recognize the others or the bridge. The sender system is configured to send specific messages to a designated queue on its native messaging infrastructure. The bridge picks up those messages and forwards them to another queue in a different messaging infrastructure where the receiver system picks them up.

### Benefits

- The systems being integrated via the Messaging Bridge don't have to be modified. Ideally, the endpoints aren't aware that the messages are bridged.
- The integration is more reliable compared to the HTTP alternative due to the at-least-once message delivery mechanism guarantee.
- Migration scenarios can be more flexible. For example, endpoints can be migrated from one messaging infrastructure to another as the schedule permits instead of all at once.

### Drawbacks

- Advanced features of one or both messaging technologies might not be available on the bridged route.
- The bridged route needs to consider both technologies' limitations. For example, the maximum message size might be 4 MB in MSMQ but only 64 KB in Azure Storage queues.

## Issues and considerations

Consider the following points when implementing the Messaging Bridge pattern:

- If one of the integrated systems relies on distributed transactions, for example Microsoft Distributed Transaction Coordinator (DTC), for correctness, you must implement a deduplication mechanism in the bridge.
- If one of the systems being integrated doesn't use any messaging infrastructure and can't be modified, you can build the Messaging Bridge between the infrastructure that's used by the other system and a SQL Server-emulated queue. The legacy system can send messages by using the change data capture feature for SQL Server to push its changes to a dedicated queue table. The bridge can forward these messages to the actual messaging infrastructure.
- You can use a single queue in each messaging infrastructure, designated as the - *bridging queue*. In this topology, configure the sending system to use that specific queue as the destination for message types that are sent to the other system. You can also use multiple pairs of queues in each messaging infrastructure, so the sender is unaware of the bridge. A- *shadow queue*is created for each destination queue in the destination system's messaging infrastructure. The bridge forwards messages between the shadow queues and their counterparts.
- In order to meet the desired availability service-level objectives (SLOs), you might need to scale out the Messaging Bridge by using the Competing consumers approach.
- Regular message-processing components use the Retry pattern to handle transient failures. The retry counter limit enables components to detect - *poison*messages and remove them from the queue to unblock processing. The bridge might require a different retry policy to prevent falsely identifying messages as poison if an infrastructure failure occurs. You might use the Circuit Breaker pattern to pause forwarding.

## When to use this pattern

Use the Messaging Bridge pattern when you need to:

- Integrate existing systems with minimal need for modification.
- Integrate legacy applications that can't use other messaging technologies.
- Extend existing on-premises applications with cloud-hosted components.
- Connect geo-distributed systems when the internet connection isn't stable.
- Migrate a single distributed system from one messaging infrastructure to another incrementally without the need to migrate the whole system in one effort.

This pattern might not be suitable if:

- At least one of the systems involved relies on a feature of one messaging infrastructure that isn't present in the other.
- Integration is synchronous in nature, and the initiating system requires immediate response.
- Integration has specific functional or nonfunctional requirements, such as security or privacy concerns.
- The volume of data for the integration exceeds the capacity of the messaging system or makes messaging an expensive solution to the problem.

## Workload design

An architect should evaluate how the Messaging Bridge pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Cost Optimization is focused on sustaining and improvingyour workload'sreturn on investment. | This intermediary step can increase the longevity of your existing system without the need for rewrites by allowing interoperability with systems that use a different messaging or eventing technology. - CO:07 Component costs |
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | This decoupling provides flexibility when you transition messaging and eventing technology within your workload or when you have heterogeneous requirements from external dependencies. - OE:06 Deploying workload changes |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

There's a .NET Framework application for managing employee scheduling hosted on-premises. The application is well-structured with separate components communicating via MSMQ. The application works, and the workload team has no intention of rewriting it. A new consumer of the scheduling data needs to be built to meet a business need, and the IT strategy calls for building new software as cloud-native applications to optimize the costs and delivery time.

The asynchronous queue-based architecture worked for the workload team in the past, so the team is going to use the same architectural approach but with the modern technology, Service Bus. The workload team doesn't want to introduce synchronous communication between the cloud and the on-premises deployment to mitigate the latency or unavailability of one affecting the other.

The team decides to use the Messaging Bridge pattern to connect the two systems. The pattern consists of two parts. One part receives messages from the existing MSMQ queue and forwards them to Service Bus. The other part takes messages from the Service Bus and forwards them to the existing MSMQ queue.

When the implementation team uses this approach, they utilize existing infrastructure in the existing application to integrate with the new components. The existing application isn't aware that the new components are hosted in Azure. Similarly, the new components communicate with the legacy application in the same way that they communicate between themselves, by sending Service Bus messages. The bridge forwards messages between the two systems.

## Contributors

*This article is maintained by Microsoft. It was originally written by the following contributors.*

Principal authors:

- Rob Bagby | Principal Content Developer - Azure Patterns & Practices
- Kyle Baley | Software Engineer
- Udi Dahan | Founder & CEO of Particular Software
- Chad Kittel | Principal Software Engineer - Azure Patterns & Practices
- Bryan Lamos | Content Developer - Azure Patterns & Practices
- Szymon Pobiega | Engineer

*To see non-public LinkedIn profiles, sign in to LinkedIn.*

## Next steps

- Messaging Bridge pattern description from the enterprise integration patterns community.
- Learn how to implement a Messaging Bridge in the Spring Java framework.
- QPid bridge can be used to bridge AMQP-enabled messaging technologies.
- The NServiceBus Messaging Bridge is a .NET implementation of a queue-to-queue bridge that supports a range of messaging infrastructures including MSMQ, Service Bus, and Azure Queue Storage.
- NServiceBus.Router is an open-source project that implements the Messaging Bridge pattern. It also allows bridging more than two technologies in a single instance and has advanced message-routing capabilities.

## Related resources

- The Competing Consumers pattern ensures the implementation of the Messaging Bridge can handle the load.
- The Retry pattern allows the Messaging Bridge to handle transient failures.
- The Circuit Breaker pattern conserves resources when either side of the bridge experiences downtime.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Decompose a task that performs complex processing into a series of separate elements that can be reused. Doing so can improve performance, scalability, and reusability of initial steps by allowing task elements that perform the processing to be deployed and scaled independently. This approach supports a high level of modularity.

## Context and problem

You have a pipeline of sequential tasks that you need to process. A straightforward but inflexible approach to implement this application is to perform this processing in a monolithic module. However, this approach is likely to reduce the opportunities for refactoring the code, optimizing it, or reusing it if parts of the same processing are required elsewhere in the application.

The following diagram illustrates one of the problems with processing data using a monolithic approach, the inability to reuse code across multiple pipelines. In this example, an application receives and processes data from two sources. A separate module processes the data from each source by performing a series of tasks to transform the data before passing the result to the business logic of the application.

Some of the tasks that the monolithic modules perform are functionally similar, but the code has to be repeated in both modules and is likely tightly coupled within its module. In addition to the inability to reuse logic, this approach introduces a risk when requirements change. You must remember to update the code in both places.

There are other challenges with a monolithic implementation unrelated to multiple pipelines or reuse. With a monolith, you don't have the ability to run specific tasks in different environments or scale them independently. Some tasks might be compute-intensive and would benefit from running on powerful hardware or running multiple instances in parallel. Other tasks might not have the same requirements. Further, with monoliths, it's challenging to reorder tasks or to inject new tasks in the pipeline. These changes require retesting the entire pipeline.

## Solution

Break down the processing required for each stream into a set of separate components (or filters), each performing a single task. Composite tasks should use multiple filters rather than one. The filters are composed into pipelines by connecting the filters with pipes. Filters are independent, self-contained, and typically stateless. Filters receive messages from an inbound pipe and publish messages to a different outbound pipe. Filters can transform the message or test it against one or more criteria to include conditional logic. Pipes don't perform routing or any other logic. They only connect filters, passing the output message from one filter as the input to the next.

Filters act independently and are unaware of other filters. They're only aware of their input and output schemas. As such, the filters can be arranged in any order so long as the input schema for any filter matches the output schema for the previous filter. Using a standardized schema for all filters enhances the ability to reorder filters. Pipes and filters architecture encourages compositional reuse.

The loose coupling of filters makes it easy to:

- Create new pipelines composed of existing filters
- Update or replace logic in individual filters
- Reorder filters, when necessary
- Run filters on differing hardware, where required
- Run filters in parallel

This diagram shows a solution implemented with pipes and filters:

The time it takes to process a single request depends on the speed of the slowest filters in the pipeline. One or more filters could be bottlenecks, especially if a high number of requests appear in a stream from a particular data source. The ability to run parallel instances of slow filters enables the system to spread the load and improve throughput.

The ability to run filters on different compute instances enables them to be scaled independently and take advantage of the elasticity that many cloud environments provide. A filter that's computationally intensive can run on high-performance hardware, while other less-demanding filters can be hosted on less-expensive commodity hardware. The filters don't even need to be in the same datacenter or geographic location, enabling each element in a pipeline to run in an environment that's close to the resources it requires. These efforts require specific design techniques such as messaging and multi-threading to maximize the elasticity of each pipe or filter. This diagram shows an example applied to the pipeline for the data from Source 1:

If the input and output of a filter are structured as a stream, you can perform the processing for each filter in parallel. The first filter in the pipeline can start its work and output its results, which are passed directly to the next filter in the sequence before the first filter completes its work.

Using the Pipes and Filters pattern together with the Compensating Transaction pattern is an alternative approach to implementing distributed transactions. You can break a distributed transaction into separate, compensable tasks, each of which can be implemented via a filter that also implements the Compensating Transaction pattern. You can implement the filters in a pipeline as separate hosted tasks that run close to the data that they maintain.

## Issues and considerations

Consider the following points when you decide how to implement this pattern:

- **Monolithic nature**. This pattern is usually implemented as a monolithic pipeline, so for any change, the entire filter chain should be tested end to end. Also, fault-tolerance for the whole process needs to be considered; if a filter or pipe fails, the whole pipeline is likely to fail.
- **Complexity**. The increased flexibility that this pattern provides can also introduce complexity, especially if the filters in a pipeline are distributed across different servers.
- **Reliability**. Use an infrastructure that ensures that data flowing between filters in a pipe aren't lost.
- **Idempotency**. If a filter in a pipeline fails after receiving a message and the work is rescheduled to another instance of the filter, part of the work might already be complete. If the work updates some aspect of the global state (like information stored in a database), a single update could be repeated. A similar issue might occur if a filter fails after it posts its results to the next filter, but before it indicates that it completed its work successfully. In these cases, another instance of the filter could repeat this work, causing the same results to be posted twice. This scenario could result in subsequent filters in the pipeline processing the same data twice. Therefore, filters in a pipeline should be designed to be idempotent. For more information, see Idempotency Patterns on Jonathan Oliver's blog.
- **Repeated messages**. If a filter in a pipeline fails after it posts a message to the next stage of the pipeline, another instance of the filter might be run, and it would post a copy of the same message to the pipeline. This scenario could cause two instances of the same message to be passed to the next filter. To avoid this problem, the pipeline should detect and eliminate duplicate messages.- Note - If you implement the pipeline by using message queues (like Azure Service Bus queues), the message queuing infrastructure might provide automatic duplicate message detection and removal.
- **Context and state**. In a pipeline, each filter essentially runs in isolation and shouldn't make any assumptions about how it was invoked. Therefore, each filter should be provided with sufficient context to perform its work. This context could include a significant amount of state information. If filters use external state, such as data in a database or external storage, then you must consider the impact on performance. Every filter has to load, operate, and persist that state, which adds overhead over solutions that load the external state a single time.
- **Message tolerance**. Filters must be tolerant of data in the incoming message that they don't operate against. They operate on the data pertinent to them and ignore other data and pass it along unchanged in the output message.
- **Error handling**- Every filter must determine what to do in the case of a breaking error. The filter must determine if it fails the pipeline or propagates the exception.

## When to use this pattern

Use this pattern when:

- The processing required by an application can easily be broken down into a set of independent steps.
- The processing steps performed by an application have different scalability requirements. - Note - You can group filters that should scale together in the same process. For more information, see the Compute Resource Consolidation pattern.
- You require the flexibility to allow reordering of the processing steps the application performs, or to allow the capability to add and remove steps.
- The system can benefit from distributing the processing for steps across different servers.
- You need a reliable solution that minimizes the effects of failure in a step while data is being processed.

This pattern might not be useful when:

- The application follows a request-response pattern.
- The task processing must be completed as part an initial request, such as a request/response scenario.
- The processing steps performed by an application aren't independent, or they have to be performed together as part of a single transaction.
- The amount of context or state information a step requires makes this approach inefficient. You might be able to persist state information to a database, but don't use this strategy if the extra load on the database causes excessive contention.

## Workload design

An architect should evaluate how the Pipes and Filters pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and to ensure that itrecoversto a fully functioning state after a failure occurs. | The single responsibility of each stage enables focused attention and avoids the distraction of commingled data processing. - RE:01 Simplicity - RE:07 Background jobs |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

You can use a sequence of message queues to provide the infrastructure required to implement a pipeline. An initial message queue receives unprocessed messages that become the start of the pipes and filters pattern implementation. A component implemented as a filter task listens for a message on this queue, performs its work, and then posts a new or transformed message to the next queue in the sequence. Another filter task can listen for messages on this queue, process them, post the results to another queue, and complete other steps, until the final step that ends the pipes and filters process. This diagram illustrates a pipeline that uses message queues:

An image processing pipeline could be implemented using this pattern. If your workload takes an image, the image could pass through a series of largely independent and reorderable filters to perform actions such as:

- content moderation
- resizing
- watermarking
- reorientation
- Exif metadata removal
- Content delivery network (CDN) publication

In this example, the filters could be implemented as individually deployed Azure Functions or even a single Azure Function app that contains each filter as an isolated deployment. The use of Azure Function triggers, input bindings, and output bindings can simplify the filter code and work automatically with a queue-based pipe using a claim check to the image to process.

Here's an example of what one filter, implemented as an Azure Function, triggered from a Queue Storage pipe with a claim Check to the image, and writing a new claim check to another Queue Storage pipe might look like. We've replaced the implementation with pseudocode in comments for brevity. More code like this can be found in the demonstration of the Pipes and Filters pattern available on GitHub.

```
// This is the "Resize" filter. It handles claim checks from input pipe, performs the
// resize work, and places a claim check in the next pipe for anther filter to handle.
[Function(nameof(ResizeFilter))]
[QueueOutput("pipe-fjur", Connection = "pipe")] // Destination pipe claim check
public async Task<string> RunAsync(
 [QueueTrigger("pipe-xfty", Connection = "pipe")] string imageFilePath, // Source pipe claim check
 [BlobInput("{QueueTrigger}", Connection = "pipe")] BlockBlobClient imageBlob) // Image to process
{
 _logger.LogInformation("Processing image {uri} for resizing.", imageBlob.Uri);
 // Idempotency checks
 // ...
 // Download image based on claim check in queue message body
 // ...

 // Resize the image
 // ...
 // Write resized image back to storage
 // ...
 // Create claim check for image and place in the next pipe
 // ...

 _logger.LogInformation("Image resizing done or not needed. Adding image {filePath} into the next pipe.", imageFilePath);
 return imageFilePath;
}
```
Note

The Spring Integration Framework has an implementation of the pipes and filters pattern.

## Next steps

You might find the following resources helpful when you implement this pattern:

- A demonstration of the Pipes and Filters Pattern using the image processing scenario is available on GitHub.
- Idempotency patterns, on Jonathan Oliver's blog.

## Related resources

The following patterns might also be relevant when you implement this pattern:

- Claim-Check pattern. A pipeline implemented using a queue might not hold the actual item being sent through the filters, but instead a pointer to the data that needs to be processed. The example uses a claim check in Azure Queue Storage for images stored in Azure Blob Storage.
- Competing Consumers pattern. A pipeline can contain multiple instances of one or more filters. This approach is useful for running parallel instances of slow filters. It enables the system to spread the load and improve throughput. Each instance of a filter competes for input with the other instances, but two instances of a filter shouldn't be able to process the same data. This article explains the approach.
- Compute Resource Consolidation pattern. It might be possible to group filters that should scale together into a single process. This article provides more information about the benefits and tradeoffs of this strategy.
- Compensating Transaction pattern. You can implement a filter as an operation that can be reversed, or that has a compensating operation that restores the state to a previous version if there's a failure. This article explains how you can implement this pattern to maintain or achieve eventual consistency.
- Pipes and Filters - Enterprise Integration Patterns.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Prioritize requests sent to services so that a workload processes high-priority requests more quickly than lower-priority ones. This approach uses messages sent to one or more queues and is useful for applications that provide different service levels or service-level agreements (SLAs) to different request types or customers.

## Context and problem

Workloads might need to manage and process tasks with varying levels of importance and urgency. Some tasks require immediate attention while others can wait. Failure to address high-priority tasks can affect the user experience and breach SLAs.

To handle tasks efficiently based on their priority, workloads need a mechanism to process and run tasks accordingly. By default, most workloads process tasks in the order they arrive, using a first-in, first-out (FIFO) queue structure. This approach doesn't account for varying task importance.

## Solution

Priority queues allow workloads to process tasks based on their priority rather than strictly by their arrival order. The application or *producer* that sends a request assigns a priority value to the message, and *consumers* process the messages by priority. The Priority Queue pattern addresses the following requirements:

- **Handles tasks of varying urgency and importance:**You have tasks with different levels of urgency and importance and need to ensure you process more critical tasks before less critical ones.
- **Handles different SLAs:**You offer different SLAs to different customers and need to ensure high-priority customers receive better performance and availability.
- **Accommodates different workload management needs:**You have a workload that needs to address certain tasks immediately while less urgent tasks can wait.

There are two main approaches to implementing the Priority Queue pattern:

- **Single queue:**Each message is assigned a priority value, and all messages use the same queue.
- **Multiple queues:**Each message is assigned a priority value, and different-priority messages use separate queues.

### Single queue

In a single queue approach, the application assigns a priority to each message and sends all messages to a single queue. The queue orders messages by priority, ensuring that consumers process higher-priority messages before lower-priority ones.

### Multiple queues

Multiple queues separate messages by priority. The application assigns a priority to each message and directs the message to the queue that corresponds to its priority, where consumers process the messages. A multiple-queue solution can use either a single pool of consumers or multiple consumer pools.

#### Single consumer pool

In a single pool setup, all queues share the same consumer pool. Consumers process messages from the highest priority queue first and process messages from lower-priority queues only when there are no more high-priority messages. As a result, single consumer pools always process higher-priority messages before lower-priority ones. This setup can lead to lower-priority messages being continually delayed and potentially never processed.

Use a single consumer pool for the following reasons:

- **Simple management.**Use a single consumer pool when easy setup and maintenance is a priority. A single pool reduces configuration and monitoring complexity.
- **Unified processing needs.**Use a single consumer pool when the incoming tasks are similar in type.

#### Multiple consumer pools

In a multiple consumer pool, each queue has a dedicated consumer pool. Higher-priority queues use more consumers or higher performance tiers to process messages more quickly than lower-priority queues.

Use multiple consumer pools for the following reasons:

- **Strict performance requirements.**Use multiple consumer pools when different task priorities have strict performance requirements that must be met independently.
- **High-reliability needs.**Use multiple consumer pools for applications when reliability and fault isolation are critical, and issues in one queue must not affect other queues.
- **Complex applications.**Use multiple consumer pools for complex applications where different tasks require different processing characteristics and performance guarantees.

## Problems and considerations

Consider the following points when you decide how to implement this pattern:

### General recommendations

- **Define priorities clearly.**Establish distinct and clear priority levels that are relevant to your solution. For example, you might define high-priority messages as those that require processing within 10 seconds. Identify the consumer requirements for handling high-priority items and allocate the necessary resources accordingly.
- **Adjust consumer pools dynamically.**Scale the size of consumer pools based on the length of the queue they're servicing.
- **Monitor queue health.**Track queue depth, processing latency, delivery count, and throughput so you can detect backlogs and slowdowns before they affect work.
- **Use dead-letter queues.**Move poison messages to a dead-letter queue after a configurable number of delivery attempts so one bad message doesn't block the priority path.
- **Prioritize service levels.**Implement priority queues to meet business needs that require prioritized availability or performance. For example, high-priority customers can receive a higher service level so they experience better performance and availability.
- **Consider low-priority processing.**Decide whether all high-priority items must be processed before any lower-priority items. If possible, dynamically increase the priority of old messages to ensure that low-priority messages eventually get processed.
- **Optimize and minimize costs.**Process critical tasks immediately with available consumers. Schedule less critical background tasks during less busy times.- If you use a single queue, optimize costs by scaling back the number of consumers. High-priority messages process first but possibly more slowly, while lower-priority messages might face longer delays.
- **Protect processors from demand peaks.**If the producer arrival rate can exceed consumer processing capacity, combine this pattern with the Queue-Based Load Leveling pattern. This approach buffers traffic bursts and helps prevent downstream processing resources from overloading.

### Multiple queue recommendations

- **Monitor processing speeds.**To ensure that messages are processed at the expected rates, continuously monitor the processing speed of high- and low-priority queues.
- **Implement preemption and suspension.**If you use multiple queues with a single consumer pool, implement an algorithm that ensures high-priority queues are always serviced before lower-priority queues.
- **Consider queue costs.**Be aware of the financial costs associated with checking and processing queues. Some queue services charge fees for posting, retrieving, and querying messages. These fees can increase with the number of queues.

## When to use this pattern

Use this pattern when:

- You must meet different latency or service-level objectives for different classes of work, such as premium versus standard customer requests.
- Work arrives in bursts, and you must protect critical operations by processing high-priority messages first while deferring lower-priority work.

This pattern might not be suitable when:

- All work items have similar business importance, and strict FIFO processing is more important than priority-based scheduling.
- Tasks have strong ordering dependencies across priority levels, and reordering work by priority can cause inconsistent outcomes or require complex coordination logic.

## Workload design

Evaluate how to use the Priority Queue pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and to ensure that itrecoversto a fully functioning state after a failure occurs. | Separating items based on business priority enables you to focus reliability efforts on the most critical work. - RE:02 Critical flows |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | Separating items based on business priority enables you to focus performance efforts on the most time-sensitive work. - PE:09 Critical flows |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

The Priority Queue pattern example on GitHub demonstrates an implementation of the Priority Queue pattern that uses Azure Service Bus topics and subscriptions. The example deploys a secure storage account, an Application Insights resource for monitoring, and a Service Bus namespace to enable communication between the sender and consumer functions.

The deployment includes three function apps: one sender and two consumers. The consumer apps use different maximum instance counts to simulate message prioritization. The `funcPriorityQueueConsumerHigh` function can scale out to 200 instances, while the `funcPriorityQueueConsumerLow` function is limited to 40 instances. All function apps use the Flex consumption plan and are connected to Application Insights for diagnostics and monitoring.

Role assignments grant secure access to Service Bus and storage by using managed identities. All function apps share the same storage account and Application Insights resource. This configuration centralizes observability and logging.

The following diagram shows the priority queue architecture:

In the preceding diagram:

- **Application (producer).**The- `PriorityQueueSender`application creates messages, assigns a custom application property called- `Priority`to each message, and sets the- `Priority`value to- `High`or- `Low`.
- **Message broker and topic.**The Service Bus message broker sends messages to a single Service Bus topic named- `messages`. Service Bus uses SQL filters to route each message to the high-priority or low-priority subscription, based on its- `Priority`value.
- **Multiple consumer pools.**The- `PriorityQueueConsumerHigh`and- `PriorityQueueConsumerLow`consumer pools respond to messages from the high-priority or low-priority subscriptions by using Azure Functions Service Bus triggers.

| Role in example | Azure service in example | Name in example |
|---|---|---|
| Application (producer) | Azure Functions app | PriorityQueueSender |
| Message broker | Azure Service Bus | <your service bus namespace> |
| Message topic | Azure Service Bus topic | `messages` |
| Message subscriptions | Azure Service Bus subscriptions | `highPriority``lowPriority` |
| Consumers | Azure Functions app | PriorityQueueConsumerHigh PriorityQueueConsumerLow |

## Next steps

- Service Bus queues, topics, and subscriptions: Review the Service Bus entities and the differences between queues and topics.
- Duplicate detection: Learn how Service Bus can reject duplicate messages when a sender retries after an uncertain send.
- Dead-letter queues: Learn how Service Bus moves messages that can't be processed to a dead-letter queue for investigation or reprocessing.
- What is Azure Queue Storage?: Review the core concepts of Azure Queue Storage to compare it with Service Bus queues.

## Related resources

The following patterns might be helpful when you implement this pattern:

- Queue-Based Load Leveling pattern: Use a queue as a buffer between request intake and processing. Use it with the Priority Queue pattern when you need both burst protection and differentiated handling.
- Competing Consumers pattern: Implement multiple consumers that listen to the same queue and process tasks in parallel to increase throughput. Only one consumer processes each message.
- Throttling pattern: Implement throttling by using queues to manage request rates. Use priority messaging to prioritize requests from critical applications or high-value customers over less important ones.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Broadcast events asynchronously to multiple interested consumers through an intermediary, without coupling the senders to the receivers. This approach, known as *pub/sub messaging*, lets a sender publish events without knowing which consumers receive them.

## Context and problem

Cloud-based and distributed applications often include system components that send information to other components as events occur. When a sender communicates directly with its consumers, it must know the identity and endpoint of every consumer, deliver messages to each consumer, and manage failures individually. Adding or removing a consumer requires changes to the sender, which limits how independently teams can develop and deploy components.

Message queues decouple senders from consumers and prevent the sender from blocking a process while it waits for a response. A standard queue creates a direct relationship between a sender and a single consumer. To support multiple consumers, the sender must create a separate queue for each consumer, which increases routing complexity and doesn't scale well. Some consumers need only a subset of the information that the sender produces, but queues don't provide built-in ways to filter messages by content or category.

Many scenarios require a sender to announce events to many interested consumers without knowing who those consumers are. Each consumer also needs a way to decide independently which events to receive.

## Solution

Introduce an asynchronous messaging subsystem that includes the following components:

- An input messaging channel that the sender uses. The sender packages events into messages by using a known message format and sends these messages via the input channel. The sender in this pattern is also known as the - *publisher*.- Note - A - *message*is a packet of data. An- *event*is a message that notifies other components about a change or an action that occurs. This pattern typically works with events, but it also carries any type of message, including commands and state notifications.
- One output messaging channel for each consumer. The consumers are known as - *subscribers*.
- A mechanism for copying each message from the input channel to the output channels for all subscribers interested in that message. An intermediary like a message broker or event bus typically handles this operation.

The following diagram shows the logical components of this pattern.

Pub/sub messaging has the following benefits:

- Decouples subsystems that need to communicate. Subsystems support independent management, and the broker retains messages even if one or more receivers are offline.
- Increases scalability and improves sender responsiveness. The sender sends a single message to the input channel and then returns to its core processing responsibilities. The messaging infrastructure routes messages to interested subscribers.
- Isolates faults. A subscriber failure doesn't affect the publisher or other subscribers, and the broker retains messages until a recovered subscriber is ready to process them.
- Supports deferred or scheduled processing. Subscribers can wait to pick up messages until off-peak hours, or the system can route or process messages according to a specific schedule.
- Supports integration between systems that use different platforms, programming languages, and communication protocols, and also connects on-premises systems with applications that run in the cloud.
- Improves testability. Channels support monitoring, and messages are available for inspection or logging as part of an integration test strategy.
- Provides separation of concerns for applications. Each application can focus on its core capabilities while the messaging infrastructure handles the work required to reliably route messages to multiple consumers.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- **Existing technologies:**Use messaging products and services that support a publish-subscribe model instead of building your own. In Azure, consider the following services:- Azure Service Bus for messaging that requires transactions, ordering, sessions, or dead-letter queues.
- Azure Event Grid for event‑based, push‑delivered notifications, especially when Azure resources change state and need to notify subscribed components.
- Azure Event Hubs for high-throughput event streaming scenarios like telemetry ingestion and log aggregation. Event Hubs uses a log-based streaming model rather than traditional pub/sub messaging, but it supports multiple consumer groups that read the same stream independently.
 - For more information, see Choose between Azure services that deliver messages. Other technologies that support pub/sub messaging include Redis, RabbitMQ, and Apache Kafka. - Libraries like MassTransit and NServiceBus provide built-in support for the publish-subscribe model on Service Bus and other messaging technologies.
- **Subscription handling:**The messaging infrastructure must provide mechanisms that consumers use to subscribe to or unsubscribe from available channels.
- **Security:**Authenticate and authorize both publishers and subscribers on a per-topic basis. Unauthorized publishers that inject messages can damage a system as much as unauthorized subscribers that read them. Encrypt messages in transit, and if content is sensitive, encrypt them at rest in the broker to prevent eavesdropping.
- **Subsets of messages:**Subscribers are usually interested in only a subset of messages from a publisher. Messaging services often let subscribers select what they receive through the following mechanisms:- **Topics:**Each topic has a dedicated output channel, and each consumer can subscribe to all relevant topics.
- **Content filtering:**The broker inspects and distributes messages based on their content. Each subscriber can specify the content that it needs.
 - Choose topic granularity deliberately. Broad topics are simpler to manage but require subscribers to filter out messages that they don't need. Narrow topics reduce subscriber-side filtering but increase the number of topics to manage. Some brokers support wildcard subscriptions, like - `orders.*`, which let subscribers match multiple topics without enumerating each one.
- **Bidirectional communication:**Channels in a publish-subscribe system are unidirectional. If a subscriber needs to acknowledge or communicate status back to the publisher, use the Request-Reply pattern. This pattern uses one channel to send a message to the subscriber and a separate reply channel to communicate back to the publisher.
- **Message ordering:**The order in which subscribers receive messages isn't guaranteed and doesn't necessarily reflect the order in which the sender created them. If ordering matters, the broker might support ordered delivery within a partition or session, but that constrains scalability. Design subscribers to handle messages independently of arrival order.
- **Message priority:**Some workloads require that specific messages be processed before others. The Priority Queue pattern provides a mechanism to route higher-priority messages before lower-priority messages.
- **Poison messages:**A malformed message, or a task that requires access to unavailable resources, can cause a service instance to fail. Capture and store these message details elsewhere for analysis. Some message brokers, like Service Bus, support this process through dead-letter queues.
- **Message size:**Brokers enforce message size limits. When payloads are large, store the content, like files or images, in an external data store and include a reference in the message. The Claim-Check pattern describes this approach.
- **Delivery guarantees and duplicate messages:**Messaging systems provide different delivery guarantees that each have trade-offs.- *At-most-once delivery*minimizes overhead but can lose messages if the broker or subscriber fails.
- *At-least-once delivery*ensures message delivery but can result in duplicates, like when a sender fails after it posts a message and a new instance repeats the post.
- *Exactly-once delivery*removes duplicates but adds coordination overhead and latency, and its availability depends on the messaging infrastructure.
 - If your broker doesn't provide deduplication, design subscribers to handle messages idempotently. Different subscribers in the same workload might require different guarantees.
- **Message expiration:**Some messages have a limited lifetime. If a receiver doesn't process a message within that period, the message becomes irrelevant and the system discards it. Set an expiration timestamp in the message data so that receivers can check its relevance before they process it.
- **Message scheduling:**A message might be embargoed and unavailable for processing until a specific date and time. Set a release timestamp so that the messaging system withholds the message until that point.
- **Message schema evolution:**Publishers and subscribers deploy independently, so message schemas change over time. Prefer backward-compatible changes, like adding optional fields, so that existing subscribers continue to work. For breaking changes, version through topic names, like- `orders.v1`and- `orders.v2`, or through a version field in message metadata. Subscribers should ignore fields that they don't recognize.
- **Correlation:**The broker decouples publishers from subscribers, which makes it harder to trace the end-to-end flow of a message. Include a correlation ID in every message so that subscribers and logging systems can connect related operations into a single trace.
- **Backpressure and scaling:**When subscribers can't keep up, unprocessed messages accumulate in the broker and can deplete its resources. Use broker flow control settings to limit unacknowledged messages for each subscriber. Scale out subscribers by using the Competing Consumers pattern when flow control alone isn't sufficient.

## When to use this pattern

Use this pattern when:

- An application needs to broadcast information to a significant number of consumers.
- An application needs to communicate with independently developed applications or services. They might use different platforms, programming languages, or communication protocols.
- An application can send information to consumers without requiring real-time responses from them.
- The systems being integrated are designed to support an eventual consistency model for their data.
- An application needs to communicate information to multiple consumers that have different availability requirements or uptime schedules than the sender.

This pattern might not be suitable when:

- An application has only a few consumers that need significantly different information from the producing application. The overhead of a broker adds complexity without any scaling benefit. Direct communication or separate queues might be more suitable.
- An application requires near real-time interaction with consumers. The pub/sub model introduces latency through the broker. Use a request-reply pattern when the publisher requires a synchronous response.
- Consumers must process messages in a specific, guaranteed order. Pub/sub systems generally don't guarantee ordering across subscribers, and maintaining order adds significant constraints to the broker and consumer design.
- The operation requires a single atomic transaction across the publisher and its consumers. Pub/sub messaging is inherently asynchronous and eventually consistent. If you need transactional guarantees, consider a direct database transaction or the Saga pattern for coordinating distributed transactions.

## Workload design

Evaluate how to use the Publisher-Subscriber pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | This pattern decouples components so that you can set independent reliability targets and remove direct dependencies. - RE:03 Failure mode analysis - RE:07 Background jobs |
| Security design decisions help ensure the confidentiality,integrity, andavailabilityof your workload's data and systems. | This pattern introduces a clear security segmentation boundary. Use it to isolate queue subscribers from the publisher at the network level. - SE:04 Segmentation |
| Cost Optimization focuses on sustainingandimprovingyour workload'sreturn on investment (ROI). | This decoupled design supports event-driven architectures that align with consumption-based billing models and help avoid overprovisioning. - CO:05 Rate optimization - CO:12 Scaling costs |
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | The broker as an intermediary lets you change the implementation on either the publisher or subscriber side without coordinating changes across both components. - OE:06 Workload development - OE:11 Safe deployment practices |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | Decoupling publishers from consumers lets you optimize compute and code for the specific tasks that each consumer handles for a given message type. - PE:02 Capacity planning - PE:05 Scaling and partitioning |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

The following diagram shows an enterprise integration architecture that uses Service Bus to coordinate workflows and Event Grid to notify subsystems of events that occur. For more information, see Enterprise integration on Azure by using message queues and events.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Consume third-party software artifacts in your supply chain only when it's verified and marked as safe-for-use, by well-defined processes. This pattern is an operational sidecar to the development process. The consumer of this pattern invokes this process to verify and block the use of software that could potentially introduce security vulnerabilities.

## Context and problem

Cloud solutions often rely on third-party software obtained from external sources. Open-source binaries, public container images, vendor OS images are some examples of these types of artifacts. All such external artifacts must be treated as *untrusted*.

In a typical workflow, the artifact is retrieved from a store outside the solution's scope and then integrated into the deployment pipeline. There are some potential issues in this approach. The source might not be trusted, the artifact might contain vulnerabilities, or it might not be suitable in some other way for the developer environment.

If these issues aren't addressed, data integrity and confidentiality guarantees of the solution might be compromised, or cause instability due to incompatibility with other components.

Some of those security issues can be avoided by adding checks to each artifact.

## Solution

Have a process that validates the software for security before introducing it in your workload. During the process, each artifact undergoes thorough operational rigor that verifies it against specific conditions. Only after the artifact satisfies those conditions, the process marks it as *trusted*.

The process of quarantining is a security measure, which consists of a series of checkpoints that are employed before an artifact is consumed. Those security checkpoints make sure that an artifact transitions from an untrusted status to a trusted status.

The quarantine process doesn't change the artifact's composition. The process is independent of the software development cycle and is invoked by consumers, as needed. As a consumer of the artifact, block the use of artifacts until they've passed quarantine.

Here's a typical quarantine workflow:

- The consumer signals their intent, specifies the input source of the artifact, and blocks its use.
- The quarantine process validates the origin of the request and gets the artifacts from the specified store.
- A custom verification process is performed as part of quarantine, which includes verifying the input constraints and checking the attributes, source, and type against established standards. - Some of these security checks can be vulnerability scanning, malware detection, and other checks that run on each submitted artifact. - The actual checks depend on the type of artifact. Evaluating an OS image is different from evaluating a NuGet package, for example.
- If the verification process is successful, the artifact is published in a safe store with clear annotations. Otherwise, it remains unavailable to the consumer. - The publishing process can include a cumulative report that shows proof of verification and the criticality of each check. Include expiration in the report beyond which the report should be invalid and the artifact is considered unsafe.
- The process marks the end of the quarantine by signaling an event with state information accompanied by a report. - Based on the information, the consumers can choose to take actions to use the trusted artifact. Those actions are outside the scope of the Quarantine pattern.

## Issues and considerations

- As a team that consumes third-party artifacts, ensure that it's obtained from a trusted source. Your choice must be aligned to organization-approved standards for artifacts that are procured from third-party vendors. The vendors must be able to meet the security requirements of your workload (and your organization). For example, make sure the vendor's responsible disclosure plan meets your organization's security requirements.
- Create segmentation between resources that stores trusted and untrusted artifacts. Use identity and network controls to restrict access to the authorized users.
- Have a reliable way to invoking the quarantine process. Make sure the artifact isn't consumed inadvertently until marked as trusted. The signaling should be automated. For example, tasks related to notifying the responsible parties when an artifact is ingested into the developer environment, committing changes to a GitHub repository, and adding an image to a private registry.
- An alternative to implementing a Quarantine pattern is to outsource it. There are quarantine practitioners who specialize in public asset validation as their business model. Evaluate both the financial and operational costs of implementing the pattern versus outsourcing the responsibility. If your security requirements need more control, implementing your own process is recommended.
- Automate the artifact ingestion process and also the process of publishing the artifact. Because validation tasks can take time, the automation process must be able to continue until all tasks are completed.
- The pattern serves as a first opportunity momentary validation. Successfully passing quarantine doesn't ensure that the artifact remains trustworthy indefinitely. The artifact must continue to undergo continuous scanning, pipeline validation, and other routine security checks that serve as last opportunity validations before promoting the release.
- The pattern can be implemented by central teams of an organization or an individual workload team. If there are many instances or variations of the quarantine process, these operations should be standardized and centralized by the organization. In this case, workload teams share the process and benefit from offloading process management.

## When to use this pattern

Use this pattern when:

- The workload integrates an artifact developed outside the scope of the workload team. Common examples include: - An Open Container Initiative (OCI) artifact from public registries such as, DockerHub, GitHub Container registry, Microsoft container registry
- A software library or package from public sources such as, the NuGet Gallery, Python Package Index, Apache Maven repository
- An external Infrastructure-as-Code (IaC) package such as Terraform modules, Community Chef Cookbooks, Azure Verified Modules
- A vendor-supplied OS image or software installer

- The workload team considers the artifact as a risk that's worth mitigating. The team understands the negative consequences of integrating compromised artifacts and the value of quarantine in assuring a trusted environment.
- The team has a clear and shared understanding of the validation rules that should be applied to a type of artifact. Without consensus, the pattern might not be effective. - For example, if a different set of validation checks are applied each time an OS image is ingested into quarantine, the overall verification process for OS images becomes inconsistent.

This pattern might not be useful when:

- The software artifact is created by the workload team or a trusted partner team.
- The risk of not verifying the artifact is less expensive than the cost of building and maintaining the quarantine process.

## Workload design

An architect and the workload team should evaluate how the Quarantine pattern can be used as part of the workload's DevSecOps practices. The underlying principles are covered in the Azure Well-Architected Framework pillars.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Security design decisions help ensure the confidentiality,integrity, andavailabilityof your workload's data and systems. | The first responsibility of security validation is served by the Quarantine pattern. The validation on an external artifact is conducted in a segmented environment before it's consumed by the development process. - SE:02 Secured development lifecycle - SE:11 Testing and validation |
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | The Quarantine pattern supports safe deployment practices (SDP) by making sure that compromised artifacts aren't consumed by the workload, which could lead to security breaches during progressive exposure deployments. - OE:03 Software development practices - OE:11 Testing and validation |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

This example applies the solution workflow to a scenario where the workload team wants to integrate OCI artifacts from public registries to an Azure Container Registry (ACR) instance, which is owned by the workload team. The team treats that instance as a trusted artifact store.

The workload environment uses Azure Policy for Kubernetes to enforce governance. It restricts container pulls only from their trusted registry instance. Additionally, Azure Monitor alerts are set up to detect any imports into that registry from unexpected sources.

- A request for an external image is made by the workload team through a custom application hosted on Azure Web Apps. The application collects the required information only from authorized users. - *Security checkpoint: The identity of requestor, the destination container registry, and the requested image source, are verified.*
- The request is stored in Azure Cosmos DB. - *Security checkpoint: An audit trail is maintained in the database, keeping track of lineage and validations of the image. This trail is also used for historical reporting.*
- The request is handled by a workflow orchestrator, which is a durable Azure Function. The orchestrator uses a scatter-gather approach for running all validations. - *Security checkpoint: The orchestrator has a managed identity with just-enough access to perform the validation tasks.*
- The orchestrator makes a request to import the image into the quarantine Azure Container Registry (ACR) that is deemed as an untrusted store.
- The import process on the quarantine registry gets the image from the untrusted external repository. If the import is successful, the quarantine registry has local copy of the image to execute validations. - *Security checkpoint: The quarantine registry protects against tampering and workload consumption during the validation process*.
- The orchestrator runs all validation tasks on the local copy of the image. Tasks include checks such as, Common Vulnerabilities and Exposures (CVE) detection, software bill of material (SBOM) evaluation, malware detection, image signing, and others. - The orchestrator decides the type of checks, the order of execution, and the time of execution. In this example, it uses Azure Container Instance as task runners and results are in the Cosmos DB audit database. All tasks can take significant time. - *Security checkpoint: This step is the core of the quarantine process that performs all the validation checks. The type of checks could be custom, open-sourced, or vendor-purchased solutions.*
- The orchestrator makes a decision. If the image passes all validations, the event is noted in the audit database, the image is pushed to the trusted registry, and the local copy is deleted from the quarantine registry. Otherwise, the image is deleted from the quarantine registry to prevent its inadvertent use. - *Security checkpoint: The orchestrator maintains segmentation between trusted and untrusted resource locations.*- Note - Instead of the orchestrator making the decision, the workload team can take on that responsibility. In this alternative, the orchestrator publishes the validation results through an API and keeps the image in the quarantine registry for a period of time. - The workload team makes the decision after reviewing results. If the results meet their risk tolerance, they pull the image from the quarantine repository into their container instance. This pull model is more practical when this pattern is used to support multiple workload teams with different security risk tolerances.

All container registries are covered by Microsoft Defender for Containers, which continuously scans for newly found issues. These issues are shown in Microsoft Defender Vulnerability Management.

## Next steps

The following guidance might be relevant when implementing this pattern:

- Recommendations for securing a development lifecycle provides guidance about using trusted units of code through all stages of the development lifecycle.
- Best practices for a secure software supply chain especially when you have NuGet dependencies in your application.
- Safeguard against malicious public packages using Azure Artifacts.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Use a queue that acts as a buffer between a task and the service that it invokes. This approach smooths intermittent heavy loads that might cause the service to fail or the task to time out. It helps minimize the effect of demand peaks on availability and responsiveness of the task and the service.

## Context and problem

Many solutions in the cloud run tasks that invoke services. In this environment, intermittent heavy loads can cause performance or reliability problems for a service.

A service might be part of the same solution as the tasks that use it, or it might be a partner service that provides access to frequently used resources. Examples of these types of services include a cache or a storage service. When multiple tasks run concurrently and use the same service, it's difficult to predict the volume of requests at any time.

A service might experience demand peaks that overload it and make the service unable to respond to requests quickly. Flooding a service with many concurrent requests can also cause the service to fail if it can't handle the contention that these requests cause.

## Solution

Place a queue between the task and the service. The task and the service run asynchronously. The task posts a message that contains the data that the service requires to the queue. The queue acts as a buffer and stores the message until the service retrieves it. The service retrieves messages from the queue and processes them. Requests from multiple tasks, which can be generated at highly variable rates, can be passed to the service through the same message queue. The following diagram shows how a queue can level the load on a service.

The queue decouples the tasks from the service so that the service can handle the messages at its own pace even when concurrent tasks generate a high volume of requests. Also, tasks aren't delayed if the service isn't available when they post messages to the queue.

This pattern provides the following benefits:

- It helps maximize availability because service delays don't immediately and directly affect the application. The application can continue to post messages to the queue even when the service isn't available or isn't currently processing messages.
- It helps maximize scalability because the number of queues and the number of services can vary to meet demand.
- It helps control costs because you only need enough service instances to meet the requirements for an average load rather than the peak load.

Note

Some services implement throttling when demand reaches a threshold that might cause system failure. Throttling can reduce the available functionality. Implement load leveling in these services to ensure that demand doesn't reach this threshold.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- Implement application logic that controls the rate at which services handle messages to avoid overwhelming the target resource. Avoid passing spikes in demand to the next stage of the system. Test the system under load to ensure that it provides the required leveling. To achieve the required leveling, adjust the number of queues and the number of service instances that handle messages.
- Message queues are a one-way communication mechanism. If a task expects a reply from a service, you might need to implement a mechanism that the service can use to send a response. For more information, see Asynchronous messaging options in Azure.
- Autoscaling without bounding consumers' aggregate downstream rate only moves the overload to downstream dependencies. This overload can increase contention for resources that these services share and diminish the effectiveness of the queue to level the load.
- If your average producer rate exceeds the consumer rate, the queue continues to grow and latency increases. Monitor queue depth and scale consumers within safe limits, or shed work at the producer.
- This pattern depends on queue durability to prevent message loss. If the broker doesn't persist messages to durable storage, a crash or capacity limit can cause enqueued data to be lost before consumers process it. Choose a queue service that persists messages to disk or replicated storage, and understand its size quotas and retention limits. For workloads that require messages to survive regional failures, evaluate geo-disaster recovery options.
- Most queue services deliver messages with at-least-once semantics, which means that consumers can receive the same message more than once. Design consumer logic to be idempotent so that processing the same message multiple times produces the same outcome and avoids problems such as duplicate records or repeated charges.
- Some messages can't be processed because they contain malformed data, reference missing resources, or trigger persistent errors. Rather than letting these messages cycle indefinitely and block the queue, route them to a dead-letter queue. Monitor dead-letter queue depth so that your operations team can investigate failures, fix the underlying problem, and resubmit messages when appropriate.
- Introducing a queue between a producer and consumer doesn't preserve the original submission order under all conditions, especially when multiple consumers process messages in parallel. If your workload requires strict ordering, use features such as message sessions in Azure Service Bus. If strict ordering isn't required, design consumers to handle messages in any order, which simplifies scaling.

## When to use this pattern

Use this pattern when:

- Your workload experiences intermittent spikes that can overwhelm downstream services.
- You need to decouple request intake from processing throughput to improve resilience and cost control.

This pattern might not be suitable when:

- The caller requires a low-latency, synchronous response.
- The workload volume is predictably low and stable, so adding queueing complexity provides little benefit.

## Workload design

Evaluate how to use the Queue-Based Load Leveling pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | The approach that this pattern describes can provide resilience against sudden spikes in demand by decoupling the arrival of tasks from their processing. It can also isolate malfunctions in queue processing so that they don't affect intake. - RE:06 Scaling |
| Cost Optimization focuses on sustaining and improvingyour workload'sreturn on investment. | Because load processing is decoupled from the request or task intake, you can use this approach to reduce the need to overprovision resources to handle peak load. - CO:12 Scaling costs |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | This approach enables intentional design for throughput performance because request intake doesn't need to correlate with the processing rate. - PE:05 Scaling and partitioning |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

A web app writes data to an external data store. If several instances of the web app run concurrently, the data store might be unable to respond to requests quickly enough, which causes requests to time out, be throttled, or otherwise fail. The following diagram shows a data store overwhelmed by concurrent requests from instances of an application.

To resolve this problem, use a queue to level the load between the application instances and the data store. An Azure Functions app reads messages from a Service Bus queue and performs the read/write requests to the data store. Azure Functions can scale instances based on Service Bus backlog by using target-based scaling, within your configured scaling bounds. You can also tune trigger concurrency settings to protect the data store. For implementation guidance, see Target-based scaling and Limit scale-out. Without this tuning, the worker layer can reintroduce back-end contention.

As a technology variation, you can implement the same pattern by using Azure Container Apps instead of Azure Functions. In that approach, a containerized worker consumes messages from Service Bus and writes to the data store. Container Apps scales the worker between configured minimum and maximum replicas based on queue-related scale rules. You can also implement the same approach by using Azure Queue Storage as the event source. For implementation guidance, see Set scaling rules in Container Apps and Deploy an event-driven job by using Container Apps.

## Next steps

The following guidance might also be relevant when implementing this pattern:

- Asynchronous messaging options in Azure: Message queues are inherently asynchronous. You might need to redesign a task's application logic if it communicates directly with a service. Similarly, you might need to refactor a service to accept requests from a message queue.
- Choose between Azure messaging services: Get more information to help you choose a messaging and queuing mechanism in Azure applications.
- Recommendations for developing background jobs: Apply this pattern to background jobs so that message queues can store requests for background tasks when the application experiences high load.
- Web-Queue-Worker architecture style: The web and the worker are both stateless. Session state can be stored in a distributed cache. The worker does long-running work asynchronously and can be triggered by messages on the queue or run on a schedule for batch processing.

## Related resources

- Competing Consumers pattern: It might be possible to run multiple instances of a service, each acting as a message consumer from the load-leveling queue. You can use this approach to adjust the rate at which messages are received and passed to a service.
- Throttling pattern: A simple way to implement throttling in a service is to use queue-based load leveling and route all requests to a service through a message queue. The service can process requests at a rate that ensures it doesn't exhaust the resources it needs and reduces the amount of possible contention.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Control the rate at which your application sends requests to a service so that you stay within the service's throttling limits and overall capacity. This approach helps you avoid or minimize throttling errors and more accurately predict throughput.

Rate limiting is appropriate in many scenarios, but it's particularly helpful for large-scale, repetitive automated tasks such as batch processing.

## Context and problem

Performing large numbers of operations against a throttled service can result in increased traffic and reduced throughput, because you need to track rejected requests and then retry the operations. As the number of operations increases, a throttling limit might require multiple passes of resending data, which results in a larger performance impact.

For example, consider the following problematic retry-on-error process for ingesting data into Azure Cosmos DB:

- Your application needs to ingest 10,000 records into Azure Cosmos DB. Each record costs 10 request units (RUs) to ingest, so a total of 100,000 RUs is required to complete the job.
- Your Azure Cosmos DB instance has 20,000 RUs provisioned capacity.
- You send all 10,000 records to Azure Cosmos DB. 2,000 records are written successfully and 8,000 records are rejected.
- You send the remaining 8,000 records to Azure Cosmos DB. 2,000 records are written successfully and 6,000 records are rejected.
- You send the remaining 6,000 records to Azure Cosmos DB. 2,000 records are written successfully and 4,000 records are rejected.
- You send the remaining 4,000 records to Azure Cosmos DB. 2,000 records are written successfully and 2,000 records are rejected.
- You send the remaining 2,000 records to Azure Cosmos DB. All are written successfully.

The ingestion job completes successfully, but only after sending 30,000 records to Azure Cosmos DB. The entire data set consists of only 10,000 records.

There are other factors to consider in this example:

- Large numbers of errors can also result in extra work to log these errors and process the resulting log data. The preceding approach handles 20,000 errors, and logging these errors might impose a processing, memory, or storage resource cost.
- Because you don't know the throttling limits of the ingestion service, you don't have a way to set expectations for how long data processing takes. Rate limiting can allow you to calculate the time required for ingestion.

## Solution

Rate limiting can reduce your traffic and potentially improve throughput by reducing the number of records sent to a service over a given period of time.

A service can throttle requests based on different metrics over time, such as:

- The number of operations (for example, 20 requests per second).
- The amount of data (for example, 2 GiB per minute).
- The relative cost of operations (for example, 20,000 RUs per second).

Regardless of the metric that you use for throttling, your rate limiting implementation will involve controlling the number and/or size of operations sent to the service during a specific time period. Rate limiting optimizes your use of the service without exceeding its throttling capacity.

In scenarios where your APIs can handle requests faster than throttled ingestion services allow, you must manage how quickly you use the service. Treating throttling only as a data-rate mismatch and buffering ingestion requests until the service recovers creates risk. If your application stops responding in this scenario, any buffered data might be lost.

To avoid this risk, consider sending your records to a durable messaging system that *can* handle your full ingestion rate. (Services such as Azure Event Hubs can handle millions of operations per second.) You can then use one or more job processors to read the records from the messaging system at a controlled rate that's within the throttled service's limits. Submitting records to the messaging system can save internal memory by allowing you to dequeue only the records that can be processed during a given time interval.

Azure provides several durable messaging services that you can use with this pattern, including:

When you send records, the time period that you use for releasing records might be more granular than the period that the service throttles on. Systems often set throttles based on timespans that you can easily comprehend and work with. However, for the computer running a service, these timeframes might be very long compared to how fast it can process information. For instance, a system might throttle per second or per minute, but commonly the code is processing on the order of nanoseconds or milliseconds.

Although it's not required, it's often recommended to send smaller numbers of records more frequently to improve throughput. So, rather than trying to batch records for a release once per second or once per minute, you can be more granular than that to keep your resource consumption (memory, CPU, and network) flowing at a more even rate. This approach prevents potential bottlenecks caused by sudden bursts of requests. For example, if a service allows 100 operations per second, the implementation of a rate limiter might even out requests by releasing 20 operations every 200 milliseconds, as shown in the following graph.

In addition, it's sometimes necessary for multiple uncoordinated processes to share a throttled service. To implement rate limiting in this scenario, you can logically partition the service's capacity and then use a distributed mutual exclusion system to manage exclusive locks on those partitions. The uncoordinated processes can then compete for locks on those partitions whenever they need capacity. For each partition that a process holds a lock for, it's granted a certain amount of capacity.

For example, if the throttled system allows 500 requests per second, you might create 20 partitions worth 25 requests per second each. If a process needed to issue 100 requests, it might ask the distributed mutual exclusion system for four partitions. The system might grant two partitions for 10 seconds. The process would then rate limit to 50 requests per second, complete the task in two seconds, and then release the lock.

One way to implement this pattern is to use Azure Storage. In this scenario, you create one 0-byte blob per logical partition in a container. Your applications can then obtain exclusive leases directly against those blobs for a short period of time (for example, 15 seconds). For every lease an application is granted, it can use that partition's amount of capacity. The application then needs to track the lease time so that, when the time expires, the application can stop using the capacity that it was granted. When you implement this pattern, you'll often want each process to attempt to lease a random partition when it needs capacity.

To further reduce latency, you might allocate a small amount of exclusive capacity for each process. A process would then only seek to obtain a lease on shared capacity if it needed to exceed its reserved capacity.

As an alternative to Azure Storage, you could also implement this kind of lease management system by using technologies such as ZooKeeper, etcd, and Redis/Redsync.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- Although the Rate Limiting pattern can reduce the number of throttling errors, your application still needs to properly handle any throttling errors that might occur.
- Ensure that retries are coordinated with rate limiting. Blind or overly aggressive retries can increase load and create retry storms, so propagate back-pressure signals (for example, HTTP 429 with - `Retry-After`) and use a limited number of retries with small random delays between attempts.
- If your application has multiple workstreams that access the same throttled service, you need to integrate all of them into your rate limiting strategy. For instance, you might support bulk loading records into a database but also querying for records in that same database. You can manage capacity by ensuring that all workstreams are gated through the same rate limiting mechanism. Alternatively, you might reserve separate pools of capacity for each workstream.
- The throttled service might be used in multiple applications. In some cases, it's possible to coordinate that usage (as shown earlier in this article). If you start to see a larger than expected number of throttling errors, that increased number might indicate contention between applications accessing a service. In this case, you might need to consider temporarily reducing the throughput imposed by your rate limiting mechanism until the usage from other applications decreases.

## When to use this pattern

Use this pattern when:

- You need to reduce throttling errors raised by a rate-limited service.
- You want to minimize traffic, as compared to naive retry-on-error approaches.
- You need to reduce memory consumption by dequeuing records only when there's sufficient capacity to process them.

This pattern might not be suitable when:

- The operation requires immediate, synchronous completion with very low latency and can't tolerate queuing or deferred processing.
- The primary bottleneck isn't request rate, but instead is concurrency or resource contention (for example, CPU saturation or long-running in-flight work). In these cases, scaling or concurrency controls are more appropriate.

## Workload design

Evaluate how to use the Rate Limiting pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | This tactic protects the client by acknowledging and honoring the limitations and costs of communicating with a service when the service prefers to avoid excessive usage. - RE:07 Self-preservation |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

The following example application allows users to submit records of various types to an API. Each record type has a unique job processor that performs the following steps:

- Validation
- Enrichment
- Insertion of the record into the database

All components of the application (API, job processor A, and job processor B) are separate processes that can be scaled independently. The processes don't directly communicate with one another.

In this example, each blob lease represents a fixed share of allowed database throughput. A processor can only dequeue and write at the combined rate of the leases that it currently holds. As processors gain or lose leases over time, their allowed write rate changes, which keeps total database traffic within the configured limit while still letting all queued work progress.

This diagram incorporates the following workflow:

- A user submits 10,000 records of type A to the API.
- The API enqueues those 10,000 records in queue A.
- A user submits 5,000 records of type B to the API.
- The API enqueues those 5,000 records in queue B.
- Job processor A sees that queue A has records and attempts to gain an exclusive lease on blob 2.
- Job processor B sees that queue B has records and attempts to gain an exclusive lease on blob 2.
- Job processor A fails to obtain the lease.
- Job processor B obtains the lease on blob 2 for 15 seconds. It can now rate limit requests to the database at a rate of 100 per second.
- Job processor B dequeues 100 records from queue B and writes them.
- One second passes.
- Job processor A sees that queue A has more records and tries to gain an exclusive lease on blob 6.
- Job processor B sees that queue B has more records and tries to gain an exclusive lease on blob 3.
- Job processor A obtains the lease on blob 6 for 15 seconds. It can now rate limit requests to the database at a rate of 100 per second.
- Job processor B obtains the lease on blob 3 for 15 seconds. It can now rate limit requests to the database at a rate of 200 per second. (It also holds the lease for blob 2.)
- Job processor A dequeues 100 records from queue A and writes them.
- Job processor B dequeues 200 records from queue B and writes them.
- One second passes.
- Job processor A sees that queue A has more records and tries to gain an exclusive lease on blob 0.
- Job processor B sees that queue B has more records and tries to gain an exclusive lease on blob 1.
- Job processor A obtains the lease on blob 0 for 15 seconds. It can now rate limit requests to the database at a rate of 200 per second. (It also holds the lease for blob 6.)
- Job processor B obtains the lease on blob 1 for 15 seconds. It can now rate limit requests to the database at a rate of 300 per second. (It also holds the lease for blobs 2 and 3.)
- Job processor A dequeues 200 records from queue A and writes them.
- Job processor B dequeues 300 records from queue B and writes them.
- And so on.

After 15 seconds, one or both jobs still won't be completed. As the leases expire, a processor should also reduce the number of requests that it dequeues and writes.

Implementations of this pattern are available in different programming languages:

## Next steps

The following guidance might also be relevant when you implement this pattern:

- Advanced request throttling with Azure API Management. Use this as complementary edge admission control to enforce per-key call-rate limits and quotas, and to return consistent back-pressure signals to clients.
- Choose between Azure messaging services. Select the best durable messaging backbone for buffering and controlled ingestion.
- Handle transient faults in Azure applications. Design retry behavior so that clients back off correctly when limits are reached.

## Related resources

The following patterns and guidance might also be relevant when you implement this pattern:

- Throttling. The Rate Limiting pattern is typically implemented in response to a throttled service.
- Retry. When requests to a throttled service result in throttling errors, it's generally appropriate to retry those requests after an appropriate interval.
- Queue-Based Load Leveling is similar to the Rate Limiting pattern but differs in several key ways: - Rate limiting doesn't necessarily need to use queues to manage load, but it does need to make use of a durable messaging service. For example, a Rate Limiting pattern can use services like Apache Kafka or Event Hubs.
- The Rate Limiting pattern introduces the concept of a distributed mutual exclusion system on partitions, which allows you to manage capacity for multiple uncoordinated processes that communicate with the same throttled service.
- A Queue-Based Load Leveling pattern is applicable whenever there's a performance mismatch between services or you want to improve resilience. So it's a broader pattern than Rate Limiting, which is more specifically concerned with efficiently accessing a throttled service.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Enable an application to handle transient failures when it tries to connect to a service or network resource, by transparently retrying a failed operation. This can improve the stability of the application.

## Context and problem

An application that communicates with elements running in the cloud has to be sensitive to the transient faults that can occur in this environment. Faults include the momentary loss of network connectivity to components and services, the temporary unavailability of a service, or timeouts that occur when a service is busy.

These faults are typically self-correcting, and if the action that triggered a fault is repeated after a suitable delay it's likely to be successful. For example, a database service that's processing a large number of concurrent requests can implement a throttling strategy that temporarily rejects any further requests until its workload has eased. An application trying to access the database might fail to connect, but if it tries again after a delay it might succeed.

## Solution

In the cloud, transient faults should be expected and an application should be designed to handle them elegantly and transparently. Doing so minimizes the effects faults can have on the business tasks the application is performing. The most common design pattern to address is to introduce a retry mechanism.

The diagram above illustrates invoking an operation in a hosted service using a retry mechanism. If the request is unsuccessful after a predefined number of attempts, the application should treat the fault as an exception and handle it accordingly.

Note

Due to the commonplace nature of transient faults, built-in retry mechanisms are now available in many client libraries and cloud services, with some degree of configurability for the number of maximum retries, the delay between retries, and other parameters. For example, Entity Framework Core provides facilities to retry failed database operations.

### Retry strategies

If an application detects a failure when it tries to send a request to a remote service, it can handle the failure using the following strategies:

- **Cancel**. If the fault indicates that the failure isn't transient or is unlikely to be successful if repeated, the application should cancel the operation and report an exception.
- **Retry immediately**. If the specific fault reported is unusual or rare, like a network packet becoming corrupted while it was being transmitted, the best course of action might be to immediately retry the request.
- **Retry after delay**. If the fault is caused by one of the more commonplace connectivity or busy failures, the network or service might need a short period while the connectivity issues are corrected or the backlog of work is cleared, so programmatically delaying the retry is a good strategy. In many cases, the period between retries should be chosen to spread requests from multiple instances of the application as evenly as possible to reduce the chance of a busy service continuing to be overloaded.

If the request still fails, the application can wait and make another attempt. If necessary, this process can be repeated with increasing delays between retry attempts, until some maximum number of requests have been attempted. The delay can be increased incrementally or exponentially, depending on the type of failure and the probability that it'll be corrected during this time.

The application should wrap all attempts to access a remote service in code that implements a retry policy matching one of the strategies listed above. Requests sent to different services can be subject to different policies.

An application should log the details of faults and failing operations. This information is useful to operators. That being said, in order to avoid flooding operators with alerts on operations where subsequently retried attempts were successful, it is best to log early failures as *informational entries* and only the failure of the last of the retry attempts as an actual error. Here is an example of how this logging model would look like.

If a service is frequently unavailable or busy, it's often because the service has exhausted its resources. You can reduce the frequency of these faults by scaling out the service. For example, if a database service is continually overloaded, it might be beneficial to partition the database and spread the load across multiple servers.

## Issues and considerations

You should consider the following points when deciding how to implement this pattern.

### Impact on performance

The retry policy should be tuned to match the business requirements of the application and the nature of the failure. For some noncritical operations, it's better to fail fast rather than retry several times and affect the throughput of the application. For example, in an interactive web application accessing a remote service, it's better to fail after a smaller number of retries with only a short delay between retry attempts, and display a suitable message to the user (for example, "please try again later"). For a batch application, it might be more appropriate to increase the number of retry attempts with an exponentially increasing delay between attempts.

An aggressive retry policy with minimal delay between attempts, and a large number of retries, could further degrade a busy service that's running close to or at capacity. This retry policy could also affect the responsiveness of the application if it's continually trying to perform a failing operation.

If a request still fails after a significant number of retries, it's better for the application to prevent further requests going to the same resource and report a failure immediately. When the period expires, the application can tentatively allow one or more requests through to see whether they're successful. For more information about this strategy, see Circuit Breaker pattern.

### Idempotency

Consider whether the operation is idempotent. If so, it's inherently safe to retry. Otherwise, retries could cause the operation to be executed more than once, with unintended side effects. For example, a service might receive the request, process the request successfully, but fail to send a response. At that point, the retry logic might re-send the request, assuming that the first request wasn't received.

### Exception type

A request to a service can fail for various reasons raising different exceptions depending on the nature of the failure. Some exceptions indicate a failure that can be resolved quickly, while others indicate that the failure is longer lasting. It's useful for the retry policy to adjust the time between retry attempts based on the type of the exception.

### Transaction consistency

Consider how retrying an operation that's part of a transaction will affect the overall transaction consistency. Fine tune the retry policy for transactional operations to maximize the chance of success and reduce the need to undo all the transaction steps.

## General guidance

- Ensure that all retry code is fully tested against various failure conditions. Check that it doesn't severely affect the performance or reliability of the application, cause excessive load on services and resources, or generate race conditions or bottlenecks.
- Implement retry logic only where the full context of a failing operation is understood. For example, if a task that contains a retry policy invokes another task that also contains a retry policy, this extra layer of retries can add long delays to the processing. It might be better to configure the lower-level task to fail fast and report the reason for the failure back to the task that invoked it. This higher-level task can then handle the failure based on its own policy.
- Log all connectivity failures that cause a retry so that underlying problems with the application, services, or resources can be identified.
- Investigate the faults that are most likely to occur for a service or a resource to discover if they're likely to be long lasting or terminal. If they are, it's better to handle the fault as an exception. The application can report or log the exception, and then try to continue either by invoking an alternative service (if one is available), or by offering degraded functionality. For more information on how to detect and handle long-lasting faults, see the Circuit Breaker pattern.

## When to use this pattern

Use this pattern when an application could experience transient faults as it interacts with a remote service or accesses a remote resource. These faults are expected to be short lived, and repeating a request that has previously failed could succeed on a subsequent attempt.

This pattern might not be useful:

- When a fault is likely to be long lasting, because this can affect the responsiveness of an application. The application might be wasting time and resources trying to repeat a request that's likely to fail.
- For handling failures that aren't due to transient faults, such as internal exceptions caused by errors in the business logic of an application.
- As an alternative to addressing scalability issues in a system. If an application experiences frequent busy faults, it's often a sign that the service or resource being accessed should be scaled up.

## Workload design

An architect should evaluate how the Retry pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and to ensure that itrecoversto a fully functioning state after a failure occurs. | Mitigating transient faults in a distributed system is a core technique for improving a workload's resilience. - RE:07 Self-preservation - RE:07 Transient faults |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

Refer to the Implement a retry policy with .NET guide for a detailed example using the Azure SDK with built-in retry mechanism support.

## Next steps

- Before writing custom retry logic, consider using a general framework such as Polly for .NET or Resilience4j for Java.
- When processing commands that change business data, be aware that retries can result in the action being performed twice, which could be problematic if that action is something like charging a customer's credit card. Using the Idempotence pattern described in this blog post can help deal with these situations.

## Related resources

- For most Azure services, the client SDKs include built-in retry logic.
- Circuit Breaker pattern. If a failure is expected to be more long lasting, it might be more appropriate to implement the Circuit Breaker pattern. Combining the Retry and Circuit Breaker patterns provides a comprehensive approach to handling faults.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Maintain data consistency in distributed systems by coordinating a sequence of local transactions across multiple services. Each service performs its operation and triggers the next step through events or messages. If a step fails, a series of compensating transactions undoes the changes that the completed steps made.

## Context and problem

A *transaction* represents a unit of work, which can include multiple operations. Within a transaction, an *event* refers to a state change that affects an entity. A *command* encapsulates all information needed to perform an action or trigger a subsequent event.

Transactions must adhere to the principles of atomicity, consistency, isolation, and durability (ACID).

- **Atomicity:**All operations succeed or no operations succeed.
- **Consistency:**Data transitions from one valid state to another valid state.
- **Isolation:**Concurrent transactions yield the same results as sequential transactions.
- **Durability:**Changes persist after they're committed, even when failures occur.

In a single service, transactions follow ACID principles because they operate within a single database. However, it can be more complex to achieve ACID compliance across multiple services.

### Challenges in microservices architectures

Microservices architectures typically assign a dedicated database to each microservice. This approach provides several benefits:

- Each service encapsulates its own data.
- Each service can use the most suitable database technology and schema for its specific needs.
- Databases for each service can be scaled independently.
- Failures in one service are isolated from other services.

Despite these advantages, this architecture complicates cross-service data consistency. Traditional database guarantees like ACID aren't directly applicable to multiple independently managed data stores. Because of these limitations, architectures that rely on interprocess communication, or traditional transaction models like two-phase commit protocol, are often better suited for the Saga pattern.

## Solution

The Saga pattern manages transactions by breaking them into a sequence of *local transactions*.

Each local transaction:

- Completes its work atomically within a single service.
- Updates the service's database.
- Initiates the next transaction via an event or message.

If a local transaction fails, the saga performs a series of *compensating transactions* to reverse the changes that the preceding local transactions made.

### Key concepts in the Saga pattern

- **Compensable transactions**can be undone or compensated for by other transactions with the opposite effect. If a step in the saga fails, compensating transactions undo the changes that the compensable transactions made.
- **Pivot transactions**serve as the point of no return in the saga. After a pivot transaction succeeds, compensable transactions are no longer relevant. All subsequent actions must be completed for the system to achieve a consistent final state. A pivot transaction can assume different roles, depending on the flow of the saga:- **Irreversible or noncompensable transactions**can't be undone or retried.
- **The boundary between reversible and committed**means that the pivot transaction can be the last undoable, or compensable, transaction. Or it can be the first retryable operation in the saga.

- **Retryable transactions**follow the pivot transaction. Retryable transactions are idempotent and help ensure that the saga can reach its final state, even if temporary failures occur. They help the saga eventually achieve a consistent state.

### Saga implementation approaches

The two typical saga implementation approaches are *choreography* and *orchestration*. Each approach has its own set of challenges and technologies to coordinate the workflow.

#### Choreography

In the choreography approach, services exchange events without a centralized controller. With choreography, each local transaction publishes domain events that trigger local transactions in other services.

| Benefits of choreography | Drawbacks of choreography |
|---|---|
| Good for simple workflows that have few services and don't need a coordination logic. | Workflow can be confusing when you add new steps. It's difficult to track which commands each saga participant responds to. |
| No other service is required for coordination. | There's a risk of cyclic dependency between saga participants because they have to consume each other's commands. |
| Doesn't introduce a single point of failure because the responsibilities are distributed across the saga participants. | Integration testing is difficult because all services must run to simulate a transaction. |

#### Orchestration

In orchestration, a centralized controller, or *orchestrator*, handles all the transactions and tells the participants which operation to perform based on events. The orchestrator performs saga requests, stores and interprets the states of each task, and handles failure recovery by using compensating transactions.

| Benefits of orchestration | Drawbacks of orchestration |
|---|---|
| Better suited for complex workflows or when you add new services. | Other design complexity requires an implementation of a coordination logic. |
| Avoids cyclic dependencies because the orchestrator manages the flow. | Introduces a point of failure because the orchestrator manages the complete workflow. |
| Clear separation of responsibilities simplifies service logic. |

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- **Shift in design thinking:**Adopting the Saga pattern requires a different mindset. It requires you to focus on transaction coordination and data consistency across multiple microservices.
- **Complexity of debugging sagas:**Debugging sagas can be complex, specifically as the number of participating services grows.
- **Irreversible local database changes:**Data can't be rolled back because saga participants commit changes to their respective databases.
- **Handling transient failures and idempotence:**The system must handle transient failures effectively and ensure idempotence, when repeating the same operation doesn't alter the outcome. For more information, see Idempotent message processing.
- **Need for monitoring and tracking sagas:**Monitoring and tracking the workflow of a saga are essential tasks to maintain operational oversight.
- **Limitations of compensating transactions:**Compensating transactions might not always succeed, which can leave the system in an inconsistent state.

### Potential data anomalies in sagas

Data anomalies are inconsistencies that can occur when sagas operate across multiple services. Because each service manages its own data, called *participant data*, there's no built-in isolation across services. This setup can result in data inconsistencies or durability problems, such as partially applied updates or conflicts between services. Typical problems include:

- **Lost updates:**When one saga modifies data without considering changes made by another saga, it results in overwritten or missing updates.
- **Dirty reads:**When a saga or transaction reads data that another saga has modified, but the modification isn't complete.
- **Fuzzy, or nonrepeatable, reads:**When different steps in a saga read inconsistent data because updates occur between the reads.

### Strategies to address data anomalies

To reduce or prevent these anomalies, consider the following countermeasures:

- **Semantic lock:**Use application-level locks when a saga's compensable transaction uses a semaphore to indicate that an update is in progress.
- **Commutative updates:**Design updates so that they can be applied in any order while still producing the same result. This approach helps reduce conflicts between sagas.
- **Pessimistic view:**Reorder the sequence of the saga so that data updates occur in retryable transactions to eliminate dirty reads. Otherwise, one saga could read dirty data, or- *uncommitted changes*, while another saga simultaneously performs a compensable transaction to roll back its updates.
- **Reread values:**Confirm that data remains unchanged before you make updates. If data changes, stop the current step and restart the saga as needed.
- **Version files:**Maintain a log of all operations performed on a record and ensure that they're performed in the correct sequence to prevent conflicts.
- **Risk-based concurrency based on value:**Dynamically choose the appropriate concurrency mechanism based on the potential business risk. For example, use sagas for low-risk updates and distributed transactions for high-risk updates.

## When to use this pattern

Use this pattern when:

- You need to ensure data consistency in a distributed system without tight coupling.
- You need to roll back or compensate if one of the operations in the sequence fails.

This pattern might not be suitable when:

- Transactions are tightly coupled.
- Compensating transactions occur in earlier participants.
- There are cyclic dependencies.

## Next step

## Related resources

The following patterns might be relevant when you implement this pattern:

- The Choreography pattern has each component of the system participate in the decision-making process about the workflow of a business transaction, instead of relying on a central point of control.
- The Compensating Transaction pattern undoes work performed by a series of steps, and eventually defines a consistent operation if one or more steps fail. Cloud-hosted applications that implement complex business processes and workflows often follow this - *eventual consistency model*.
- The Retry pattern lets an application handle transient failures when it tries to connect to a service or network resource by transparently retrying the failed operation. This pattern can improve the stability of the application.
- The Circuit Breaker pattern handles faults that take a variable amount of time to recover from, when you connect to a remote service or resource. This pattern can improve the stability and resiliency of an application.
- The Health Endpoint Monitoring pattern implements functional checks in an application that external tools can access through exposed endpoints at regular intervals. This pattern can help you verify that applications and services are performing correctly.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Coordinate a set of distributed actions as a single operation. If any of the actions fail, try to handle the failures transparently, or else undo the work that was performed, so the entire operation succeeds or fails as a whole. This can add resiliency to a distributed system, by enabling it to recover and retry actions that fail due to transient exceptions, long-lasting faults, and process failures.

## Context and problem

An application performs tasks that include multiple steps, some of which might invoke remote services or access remote resources. The individual steps might be independent of each other, but they are orchestrated by the application logic that implements the task.

Whenever possible, the application should ensure that the task runs to completion and resolve any failures that might occur when accessing remote services or resources. Failures can occur for many reasons. For example, the network might be down, communications could be interrupted, a remote service might be unresponsive or in an unstable state, or a remote resource might be temporarily inaccessible, perhaps due to resource constraints. In many cases the failures will be transient and can be handled by using the Retry pattern.

If the application detects a more permanent fault it can't easily recover from, it must be able to restore the system to a consistent state and ensure integrity of the entire operation.

## Solution

The Scheduler Agent Supervisor pattern defines the following actors. These actors orchestrate the steps to be performed as part of the overall task.

- The - **Scheduler**arranges for the steps that make up the task to be executed and orchestrates their operation. These steps can be combined into a pipeline or workflow. The Scheduler is responsible for ensuring that the steps in this workflow are performed in the correct order. As each step is performed, the Scheduler records the state of the workflow, such as "step not yet started," "step running," or "step completed." The state information should also include an upper limit of the time allowed for the step to finish, called the complete-by time. If a step requires access to a remote service or resource, the Scheduler invokes the appropriate Agent, passing it the details of the work to be performed. The Scheduler typically communicates with an Agent using asynchronous request/response messaging. This can be implemented using queues, although other distributed messaging technologies could be used instead.- The Scheduler performs a similar function to the Process Manager in the Process Manager pattern. The actual workflow is typically defined and implemented by a workflow engine that's controlled by the Scheduler. This approach decouples the business logic in the workflow from the Scheduler.
- The - **Agent**contains logic that encapsulates a call to a remote service or access to a remote resource that a step in a task references. Each Agent typically wraps calls to a single service or resource and implements the appropriate error handling and retry logic. A timeout constraint, described later in this article, applies to error handling and retry logic. When you implement retry logic, pass a stable identifier across all retry attempts so that the remote service can use it for any deduplication logic that it might have.- If the steps in the workflow that the Scheduler runs use several services and resources across different steps, each step might reference a different Agent. This point is an implementation detail of the pattern. For guidance about how to design retry strategies, see Transient fault handling.
- The - **Supervisor**monitors the status of the steps in the task being performed by the Scheduler. It runs periodically (the frequency will be system-specific), and examines the status of steps maintained by the Scheduler. If it detects any that have timed out or failed, it arranges for the appropriate Agent to recover the step or execute the appropriate remedial action (this might involve modifying the status of a step). The Scheduler and Agents implement the recovery or remedial actions. The Supervisor requests that they carry out these actions.

The Scheduler, Agent, and Supervisor are logical components and their physical implementation depends on the technology being used. For example, several logical agents might be implemented as part of a single web service.

The Scheduler maintains information about the progress of the task and the state of each step in a durable data store, called the state store. The Supervisor can use this information to help determine whether a step has failed. The figure illustrates the relationship between the Scheduler, the Agents, the Supervisor, and the state store.

Note

This diagram shows a simplified version of the pattern. In a real implementation, there might be many instances of the Scheduler running concurrently, each a subset of tasks. Similarly, the system could run multiple instances of each Agent, or even multiple Supervisors. In this case, Supervisors must coordinate their work with each other carefully to ensure that they don't compete to recover the same failed steps and tasks. The Leader Election pattern provides one possible solution to this problem.

When the application is ready to run a task, it submits a request to the Scheduler. The Scheduler records initial state information about the task and its steps (for example, step not yet started) in the state store and then starts performing the operations defined by the workflow. As the Scheduler starts each step, it updates the information about the state of that step in the state store (for example, step running).

If a step references a remote service or resource, the Scheduler sends a message to the appropriate Agent. The message contains the information that the Agent needs to pass to the service or access the resource, in addition to the complete-by time for the operation. If the Agent completes its operation successfully, it returns a response to the Scheduler. The Scheduler can then update the state information in the state store (for example, step completed) and perform the next step. This process continues until the entire task is complete.

An Agent can implement any retry logic that's necessary to perform its work. However, if the Agent doesn't complete its work before the complete-by period expires, the Scheduler will assume that the operation has failed. In this case, the Agent should stop its work and not try to return anything to the Scheduler (not even an error message), or try any form of recovery. The reason for this restriction is that, after a step has timed out or failed, another instance of the Agent might be scheduled to run the failing step (this process is described later).

If the Agent fails, the Scheduler won't receive a response. The pattern doesn't make a distinction between a step that has timed out and one that has genuinely failed.

If a step times out or fails, the state store will contain a record that indicates that the step is running, but the complete-by time will have passed. The Supervisor looks for steps like this and tries to recover them. One possible strategy is for the Supervisor to update the complete-by value to extend the time available to complete the step, and then send a message to the Scheduler identifying the step that has timed out. The Scheduler can then try to repeat this step. However, this design requires the tasks to be idempotent. The system should contain infrastructure to maintain consistency. For more information, see Repeatable Infrastructure, Architect Azure applications for resiliency and availability, and Resource consistency decision guide.

The Supervisor might need to prevent the same step from being retried if it continually fails or times out. To do this, the Supervisor could maintain a retry count for each step, along with the state information, in the state store. If this count exceeds a predefined threshold the Supervisor can adopt a strategy of waiting for an extended period before notifying the Scheduler that it should retry the step, in the expectation that the fault will be resolved during this period. Alternatively, the Supervisor can send a message to the Scheduler to request the entire task be undone by implementing a Compensating Transaction pattern. This approach will depend on the Scheduler and Agents providing the information necessary to implement the compensating operations for each step that completed successfully.

It isn't the purpose of the Supervisor to monitor the Scheduler and Agents, and restart them if they fail. This aspect of the system should be handled by the infrastructure these components are running in. Similarly, the Supervisor shouldn't have knowledge of the actual business operations that the tasks being performed by the Scheduler are running (including how to compensate should these tasks fail). This is the purpose of the workflow logic implemented by the Scheduler. The sole responsibility of the Supervisor is to determine whether a step has failed and arrange either for it to be repeated or for the entire task containing the failed step to be undone.

If the Scheduler is restarted after a failure, or the workflow being performed by the Scheduler terminates unexpectedly, the Scheduler should be able to determine the status of any inflight task that it was handling when it failed, and be prepared to resume this task from that point. The implementation details of this process are likely to be system-specific. If the task can't be recovered, it might be necessary to undo the work already performed by the task. This might also require implementing a compensating transaction.

The key advantage of this pattern is that the system is resilient in the event of unexpected temporary or unrecoverable failures. The system can be constructed to be self-healing. For example, if an Agent or the Scheduler fails, a new one can be started and the Supervisor can arrange for a task to be resumed. If the Supervisor fails, another instance can be started and can take over from where the failure occurred. If the Supervisor is scheduled to run periodically, a new instance can be automatically started after a predefined interval. The state store can be replicated to reach an even greater degree of resiliency.

## Issues and considerations

You should consider the following points when deciding how to implement this pattern:

- This pattern can be difficult to implement and requires thorough testing of each possible failure mode of the system.
- The recovery/retry logic implemented by the Scheduler is complex and dependent on state information held in the state store. It might also be necessary to record the information required to implement a compensating transaction in a durable data store. A compensating transaction might fail as well.
- How often the Supervisor runs will be important. It should run often enough to prevent any failed steps from blocking an application for an extended period, but it shouldn't run so often that it becomes an overhead.
- The steps performed by an Agent could be run more than once. The logic that implements these steps should be idempotent.

## When to use this pattern

Use this pattern when a process that runs in a distributed environment, like the cloud, must remain resilient to both communications failures and operational failures. This pattern is common in background jobs that orchestrate multistep workflows, like order processing or resource provisioning.

This pattern might not be suitable for tasks that don't invoke remote services or access remote resources.

## Workload design

An architect should evaluate how the Scheduler Agent Supervisor pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and to ensure that itrecoversto a fully functioning state after a failure occurs. | This pattern uses health metrics to detect failures and reroute tasks to a healthy agent in order to mitigate the effects of a malfunction. - RE:05 Redundancy - RE:07 Self-healing |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, code. | This pattern uses performance and capacity metrics to detect current utilization and route tasks to an agent that has capacity. You can also use it to prioritize the execution of higher priority work over lower priority work. - PE:05 Scaling and partitioning - PE:09 Critical flows |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

A web application that implements an ecommerce system has been deployed on Microsoft Azure. Users can run this application to browse the available products and to place orders. The user interface runs as a web frontend, and the order processing elements of the application are implemented as a set of background workers. Part of the order processing logic involves accessing a remote service, and this aspect of the system could be prone to transient or more long-lasting faults. For this reason, the designers used the Scheduler Agent Supervisor pattern to implement the order processing elements of the system.

When a customer places an order, the application constructs a message that describes the order and posts this message to a queue. A separate submission process, running as a background worker, retrieves the message, inserts the order details into the orders database, and creates a record for the order process in the state store. The inserts into the orders database and the state store are performed as part of the same operation. The submission process is designed to ensure that both inserts complete together.

The state information that the submission process creates for the order includes:

- **OrderID**. The ID of the order in the orders database.
- **LockedBy**. The instance ID of the worker handling the order. There might be multiple current instances of the worker running the Scheduler, but each order should only be handled by a single instance.
- **CompleteBy**. The time the order should be processed by.
- **ProcessState**. The current state of the task handling the order. The possible states are:- **Pending**. The order has been created but processing hasn't yet been started.
- **Processing**. The order is currently being processed.
- **Processed**. The order has been processed successfully.
- **Error**. The order processing has failed.

- **FailureCount**. The number of times that processing has been tried for the order.

In this state information, the `OrderID` field is copied from the order ID of the new order. The `LockedBy` and `CompleteBy` fields are set to `null`, the `ProcessState` field is set to `Pending`, and the `FailureCount` field is set to 0.

Note

In this example, the order handling logic is relatively simple and only has a single step that invokes a remote service. In a more complex multistep scenario, the submission process would likely involve several steps, and so several records would be created in the state store — each one describing the state of an individual step.

The Scheduler also runs as a background worker and implements the business logic that handles the order. An instance of the Scheduler polling for new orders examines the state store for records where the `LockedBy` field is null and the `ProcessState` field is pending. When the Scheduler finds a new order, it immediately populates the `LockedBy` field with its own instance ID, sets the `CompleteBy` field to an appropriate time, and sets the `ProcessState` field to processing. The code is designed to be exclusive and atomic to ensure that two concurrent instances of the Scheduler can't try to handle the same order simultaneously.

The Scheduler then runs the business workflow to process the order asynchronously, passing it the value in the `OrderID` field from the state store. The workflow handling the order retrieves the details of the order from the orders database and performs its work. When a step in the order processing workflow needs to invoke the remote service, it uses an Agent. The workflow step communicates with the Agent using a pair of Azure Service Bus message queues acting as a request/response channel. The figure shows a high-level view of the solution.

The message sent to the Agent from a workflow step describes the order and includes the complete-by time. If the Agent receives a response from the remote service before the complete-by time expires, it posts a reply message on the Service Bus queue on which the workflow is listening. When the workflow step receives the valid reply message, it completes its processing and the Scheduler sets the `ProcessState` field of the order state to processed. At this point, the order processing has completed successfully.

If the complete-by time expires before the Agent receives a response from the remote service, the Agent halts its processing and terminates handling the order. Similarly, if the workflow handling the order exceeds the complete-by time, it also terminates. In both cases, the state of the order in the state store remains set to processing, but the complete-by time indicates that the time for processing the order has passed and the process is deemed to have failed. If either the Agent that accesses the remote service or the workflow that handles the order terminate unexpectedly, the information in the state store remains set to processing and eventually has an expired complete-by value.

If the Agent detects an unrecoverable, nontransient fault while it's trying to contact the remote service, it can send an error response back to the workflow. The Scheduler can set the status of the order to error and raise an event that alerts an operator. The operator can then try to resolve the reason for the failure manually and resubmit the failed processing step.

The Supervisor periodically examines the state store looking for orders with an expired complete-by value. If the Supervisor finds a record, it increments the `FailureCount` field. If the failure count value is below a specified threshold value, the Supervisor resets the `LockedBy` field to null, updates the `CompleteBy` field with a new expiration time, and sets the `ProcessState` field to pending. An instance of the Scheduler can pick up this order and perform its processing as before. If the failure count value exceeds a specified threshold, the reason for the failure is assumed to be nontransient. The Supervisor sets the status of the order to error and raises an event that alerts an operator.

In this example, the Supervisor is implemented as a separate background worker. You can use various strategies to arrange for the Supervisor task to be run, such as using timer-triggered functions in Azure Functions or a scheduled job in Azure Container Apps.

Although it isn't shown in this example, the Scheduler might need to keep the application that submitted the order informed about the progress and status of the order. The application and the Scheduler are isolated from each other to eliminate any dependencies between them. The application has no knowledge of which instance of the Scheduler is handling the order, and the Scheduler is unaware of which specific application instance posted the order.

To allow the order status to be reported, the application could use its own private response queue. The details of this response queue would be included as part of the request sent to the submission process, which would include this information in the state store. The Scheduler would then post messages to this queue indicating the status of the order (such as request received, order completed, and order failed). It should include the order ID in these messages so they can be correlated with the original request by the application.

## Next steps

The following guidance might also be relevant when implementing this pattern:

- Asynchronous Messaging Primer. The components in the Scheduler Agent Supervisor pattern typically run decoupled from each other and communicate asynchronously. Describes some of the approaches that can be used to implement asynchronous communication based on message queues.
- Reference 6: A Saga on Sagas. An example showing how the CQRS pattern uses a process manager (part of the CQRS Journey guidance).

## Related resources

The following patterns might also be relevant when implementing this pattern:

- Retry pattern. An Agent can use this pattern to transparently retry an operation that accesses a remote service or resource that has previously failed. Use when the expectation is that the cause of the failure is transient and can be corrected.
- Circuit Breaker pattern. An Agent can use this pattern to handle faults that take a variable amount of time to correct when connecting to a remote service or resource.
- Compensating Transaction pattern. If the workflow being performed by a Scheduler can't be completed successfully, it might be necessary to undo any work it's previously performed. The Compensating Transaction pattern describes how this can be achieved for operations that follow the eventual consistency model. These types of operations are commonly implemented by a Scheduler that performs complex business processes and workflows.
- Leader Election pattern. It might be necessary to coordinate the actions of multiple instances of a Supervisor to prevent them from attempting to recover the same failed process. The Leader Election pattern describes how to do this.
- Cloud Architecture: The Scheduler-Agent-Supervisor pattern on Clemens Vasters' blog

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Group related messages by a category key and process each group sequentially, one message at a time, while processing different groups in parallel.

This pattern resolves the tension between maintaining first-in, first-out (FIFO) correctness within each logical group and scaling out concurrent processing across groups. The design ensures that ordering constraints don't become a system-wide bottleneck.

## Context and problem

Applications often need to process related messages in the order they arrive while still scaling out to handle increased load. In a distributed architecture, this requirement is difficult to achieve because workers independently pull messages from a shared queue. When multiple workers compete for messages, as in the Competing Consumers pattern, ordering breaks down.

Consider an order-tracking system that receives a stream of operations, such as creating an order, adding a transaction, modifying a past transaction, and deleting an order. Each order's operations must be processed in FIFO order, because applying them out of sequence would corrupt the order's state. However, the incoming queue interleaves operations across many orders. A single consumer that enforces global ordering becomes a bottleneck, and multiple consumers might process the same order's operations out of sequence.

The straightforward approaches to this problem each break down in a different way:

- **Single consumer.**A single consumer preserves message order because it processes one message at a time, but it can't scale to handle increased throughput.
- **Multiple competing consumers.**Multiple consumers scale throughput by pulling messages in parallel, but they lose per-group ordering guarantees. Two workers can pull consecutive messages for the same order and process them simultaneously or out of sequence, which corrupts the order state.

## Solution

The Sequential Convoy pattern partitions related messages into categories and processes each category sequentially, one message at a time, while categories are processed in parallel.

The pattern works by assigning each message a category key that identifies the group it belongs to. A message broker uses this key to partition messages into logical groups. Within each group, the broker enforces FIFO ordering so that a consumer that locks a group receives messages strictly in the sequence that they were enqueued. Different groups can be processed by different consumers simultaneously, so the system scales horizontally across groups without sacrificing ordering within any single group.

On Azure, Azure Service Bus message sessions provide a built-in implementation of this pattern.

The following diagram shows the general Sequential Convoy pattern.

In the queue, messages for different categories might be interleaved, as shown in the following diagram.

This pattern provides several key benefits:

- **Ordered processing per group.**Messages within each category are processed strictly in sequence, which prevents race conditions, out-of-order state mutations, and the need for reordering workarounds.
- **Horizontal scale across groups.**Each category is an independent unit of concurrency. Adding consumers increases throughput proportionally to the number of active categories, without breaking ordering guarantees.
- **Producer-consumer decoupling.**Producers enqueue messages without knowledge of which consumer will process them or when. Consumers are independently scalable and replaceable.

## Problems and considerations

Consider the following points when you decide how to implement this pattern:

- **Category and scale unit.**Determine what property of your incoming messages you can scale out on. The category key defines the unit of parallelism: each distinct key value becomes an independently processable group. In the order-tracking scenario, this property is the order ID. Choosing a key that is too coarse (for example, a single customer ID for all orders) limits parallelism, while choosing a key that is too fine doesn't provide meaningful ordering benefit.
- **Throughput limits.**Evaluate your target message throughput. Because this pattern enforces sequential processing within each category, throughput per category is bounded by the time to process a single message. Optimize per-message processing time, for example, by using asynchronous I/O or batching downstream writes, because that time directly determines the maximum throughput for each category. If your overall throughput requirement is very high, reconsider whether strict FIFO ordering is necessary for the entire message lifecycle. Alternatives include enforcing a start message and end message to bracket a sequence, or sorting messages by timestamp within a batch window and then sending the batch for parallel processing.
- **Service capabilities.**Verify whether your choice of message broker supports one-at-a-time processing of messages within a queue or category of a queue. Not all messaging services provide session-level locking or FIFO guarantees within a partition. If the broker doesn't natively support this capability, the consumer must implement its own coordination logic, which adds complexity and risks duplicate processing, missed messages, or out-of-order execution. Session support might also constrain the choice of messaging tier or SKU, which affects cost.
- **Evolvability.**Plan how you will add new categories of messages to the system. The pattern must accommodate growth in category cardinality without requiring structural changes to consumers. For example, suppose the ledger system described earlier is specific to one customer. If you need to onboard a new customer, you should be able to add a set of ledger processors that distribute work per customer ID without redesigning the queue topology.
- **Out-of-order message delivery.**Messages can arrive out of order because of variable network latency between the producer and the broker, before the broker's session ordering takes effect. Consider using sequence numbers to verify ordering within each category. You can also include an end-of-sequence flag in the last message of a transaction so consumers can detect when a sequence is complete.
- **Poison message handling.**A message that repeatedly fails processing within a session blocks all subsequent messages in that session because the pattern enforces strict sequential ordering. Design a strategy to detect poison messages, such as tracking delivery attempt counts, and move them to a dead-letter queue after a defined retry threshold so the remaining messages in the session can continue processing.
- **Broker availability.**The message broker is a shared dependency for all categories. Its availability and durability directly affect the pattern's reliability guarantees. Evaluate broker-level resiliency features such as availability zones and geo-disaster recovery based on the workload's availability requirements and budget, because higher-durability configurations typically increase cost.
- **Producer key correctness.**The pattern assumes that producers set the category key (session ID) correctly on every message. If a producer sets an incorrect key, either accidentally or because of a bug, the message routes to the wrong session and corrupts that group's state. Validate that producers assign category keys consistently, and consider adding key-validation logic at the consumer if the consequence of a misrouted message is severe.
- **Operational complexity.**Monitoring session-based processing adds operational overhead above that of standard queue consumption. Operators need visibility into session backlogs (the number of active sessions and the depth of messages waiting in each session) to identify categories that are falling behind. Dead-lettered sessions require a separate monitoring and remediation workflow to investigate failed messages, resolve the root cause, and replay corrected messages back into the session.
- **Session lock contention and latency.**Session locking introduces latency overhead because each consumer must acquire an exclusive lock on a session before processing messages. When a consumer holds a session lock, no other consumer can process messages from that session, even if the consumer is slow or temporarily stalled. If the lock duration is too short, lock expiration can cause message reprocessing. If the lock duration is too long, a stalled consumer delays recovery. Tune the session lock duration based on expected message processing time, and implement lock renewal for longer-running operations.
- **Consumer scaling and cost.**Parallelism across sessions translates to concurrent consumer instances. In a serverless model such as Azure Functions, each active session maps to a concurrent execution, and in a dedicated model, it maps to an instance or thread. The number of active sessions therefore directly influences compute cost. Plan consumer scaling limits and concurrency controls to balance throughput against cost.

## When to use this pattern

Use this pattern when:

- Messages arrive in order and must be processed in the same order.
- Messages can be categorized so that each category becomes an independent unit of scale for the system.

This pattern might not be suitable when:

- You expect extremely high throughput scenarios (millions of messages per minute), because the FIFO requirement limits the scaling that the system can achieve.
- Message ordering isn't required. When messages can be processed independently in any order, the Competing Consumers pattern provides simpler horizontal scaling without the coordination overhead of session locking.

## Workload design

Evaluate how to use Sequential Convoy in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | This pattern uses session-based FIFO ordering to eliminate race conditions, contention-prone message handling logic, and other workarounds for incorrectly ordered messages that can lead to malfunctions. - RE:02 Critical flows - RE:07 Background jobs |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

On Azure, you can implement this pattern by using Service Bus message sessions. For the consumers, you can use either Azure Logic Apps with the Service Bus peek-lock connector or Azure Functions with the Service Bus trigger.

When a producer sets the `SessionId` property on a message, Service Bus groups all messages that share the same session ID into a single logical session. A consumer accepts a session and receives an exclusive lock on it. This lock guarantees that only one consumer processes messages for that session at any time and that messages arrive in FIFO order. Other consumers can simultaneously accept and process different sessions, providing parallel throughput across groups.

In the order-tracking example, the system processes each ledger message in the order it's received and sends each transaction to another queue where the category is set to the order ID. A transaction never spans multiple orders in this scenario, so consumers process each category in parallel but FIFO within the category.

The ledger processor fans out the messages by de-batching the content of each message in the first queue:

The ledger processor performs three steps:

- Walks the ledger one transaction at a time.
- Sets the session ID of the message to match the order ID.
- Sends each ledger transaction to a secondary queue with the session ID set to the order ID.

The consumers listen to the secondary queue and process all messages with matching order IDs in FIFO order. Consumers use peek-lock mode.

The ledger queue is a serial-to-parallel transition point: all transactions pass through it sequentially before fanning out to session-based parallel processing. This serialization stage is the primary scalability bottleneck because it gates the throughput of the entire downstream pipeline. However, after the ledger processor fans out messages to the secondary queue, consumers can scale independently across sessions, one per order ID.

## Supporting technologies

- Service Bus message sessions: Groups messages by session ID and enforces FIFO processing within each session. Message sessions are the primary Azure mechanism for implementing the Sequential Convoy pattern.
- Azure Functions Service Bus trigger: Supports session-based triggers that allow function instances to process messages from a single session at a time.
- Logic Apps Service Bus connector: Provides a Service Bus connector with peek-lock support for consuming session-enabled queues in workflow-based processing.

## Contributors

*Microsoft maintains this article. The following contributors wrote this article.*

Principal author:

- Naga Venkata Cheruvu | Senior Cloud Solution Architect + AI infra

*To see nonpublic LinkedIn profiles, sign in to LinkedIn.*

## Related resources

- Competing Consumers pattern: Multiple consumers pull messages from a shared queue in parallel, which increases throughput but removes per-message ordering guarantees. The Sequential Convoy pattern addresses the ordering gap that Competing Consumers introduces. It addresses this gap by partitioning messages into category-keyed sessions and processing each session sequentially.
- Queue-Based Load Leveling pattern: A queue buffers work between producers and consumers to absorb bursts and smooth uneven load. The Sequential Convoy pattern builds on this buffering by adding session-based partitioning, so that the queue both levels load across categories and preserves FIFO ordering within each category.
- Priority Queue pattern: Messages are routed to separate queues or given priority within a queue so that higher-priority work is processed before lower-priority work. When ordering within a priority level must also be preserved, the Sequential Convoy pattern can be combined with priority queuing to enforce FIFO processing within each priority-keyed session.
- Peek-Lock Message (Non-Destructive Read): This operation atomically retrieves and locks a message from a queue or subscription for processing.
- In order delivery of correlated messages in Logic Apps by using Service Bus sessions: This blog post describes Logic Apps support for the Sequential Convoy pattern.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Divide a data store into a set of horizontal partitions or shards. This approach can improve scalability when you store and access large volumes of data.

## Context and problem

A data store on a single server has the following limitations:

- **Storage space:**A data store for a large-scale cloud application can contain a large volume of data that grows over time. A server provides a finite amount of disk storage, and you can replace existing disks with larger ones or add more disks as data volumes grow. The system eventually reaches a limit where you can't increase storage capacity on a single server.
- **Computing resources:**A cloud application must support a large number of concurrent users who each run queries against the data store. A single server might not provide enough computing power for this load, which results in extended response times and timeouts. You can add memory or upgrade processors, but the system reaches a limit where you can't increase compute resources any further.
- **Network bandwidth:**The rate at which a single server can receive requests and send replies limits data store performance. The volume of network traffic can exceed the capacity of the network connection, which results in failed requests.
- **Geography:**Legal, compliance, or performance requirements might require you to store user data in the same geographic region as the users. If users span across countries/regions, you might not be able to store all the application's data in a single data store.

To postpone these limitations temporarily, you can scale vertically by adding disk capacity, processing power, memory, and network connections. A cloud application that must support a large number of users and high data volumes needs to scale horizontally.

## Solution

Divide the data store into horizontal partitions or shards. Each shard has the same schema but contains its own distinct subset of the data. Each shard is a complete data store that can contain data for many entities of different types. A shard runs on a server that functions as a storage node.

This pattern has the following benefits:

- You can scale out the system by adding more shards on extra storage nodes.
- A system can use prebuilt hardware rather than specialized and expensive computers for each storage node.
- You can reduce contention and improve performance by balancing the workload across shards.
- In the cloud, shards can reside physically close to the users who access the data.

When you divide a data store into shards, decide which data to place in each shard. Each shard typically holds items grouped by one or more data attributes. These attributes form the shard key, sometimes referred to as the *partition key*.

Sharding physically organizes the data. When an application stores and retrieves data, the sharding logic directs it to the appropriate shard. You can implement this logic in the application's data access code or in the data storage system if it transparently supports sharding.

Abstracting the physical location of the data in the sharding logic provides control over which shards contain which data. You can also migrate data between shards without modifying application business logic when you need to redistribute data, such as when shards become unbalanced. The trade-off is the extra data access overhead to determine each data item's location during retrieval.

### Shard key selection

The shard key is the most critical design decision in a sharded system. To change a shard key after you choose it, you typically must migrate all data to a new shard layout, which is an expensive and risky operation on a live system. Make this decision carefully before you write any code.

An effective shard key is immutable, has high cardinality, distributes data and load evenly, and aligns with your dominant query patterns so that most requests resolve against a single shard. Avoid monotonically increasing values (autoincrement integers and sequential timestamps), low-cardinality attributes (booleans and small enum sets), and volatile attributes that change frequently. These attributes lead to hotspots or costly cross-shard data movement.

If no single attribute meets these criteria, define a composite shard key by combining two or more attributes. If queries need to retrieve data by attributes that aren't part of the shard key, use a pattern such as the Index Table pattern to provide secondary lookups.

For more information about how to choose partition keys across Azure services, see Data partitioning guidance and Data partitioning strategies.

## Sharding strategies

Use one of the following strategies when you select the shard key and decide how to distribute data across shards. You don't need a one-to-one correspondence between shards and the servers that host them. A single server can host multiple shards.

### Lookup sharding strategy

In the lookup strategy, also called the *directory-based strategy*, the sharding logic implements a map that routes a data request to the shard that contains that data by using the shard key. In a multitenant application, you might store all the data for a tenant together in a shard by using the tenant ID as the shard key. Multiple tenants might share the same shard, but the data for a single tenant doesn't spread across multiple shards. The following diagram shows sharding tenant data based on tenant IDs.

The mapping between shard key values and physical storage can be direct, where each shard key value maps to a physical partition. A more flexible technique is virtual partitioning, where shard key values map to virtual shards, and the system then maps those virtual shards to fewer physical partitions. An application locates data by using a shard key value that refers to a virtual shard, and the system transparently maps virtual shards to physical partitions. The mapping between a virtual shard and a physical partition can change without requiring application code modifications.

### Range-based sharding strategy

The range-based strategy groups related items together in the same shard and orders them by sequential shard key. This strategy supports applications that frequently retrieve sets of items by using range queries. Range queries return a set of data items for a shard key that falls within a given range.

For example, if an application regularly needs to find all orders placed in a given month, you can retrieve the data faster if you store all orders for a month in date and time order in the same shard. If you store each order in a different shard, the application has to fetch them individually by performing a large number of point queries. The following diagram shows sequential sets, or ranges, of data stored in shards.

In this example, the shard key is a composite key that contains the order month as the most significant element, followed by the order day and time. New orders are naturally sorted as they're created and added to a shard.

Some data stores support two-part shard keys. A partition key identifies the shard, and a row key uniquely identifies an item within the shard. The shard typically stores data in row key order. For items that need range queries and must be grouped together, you can use a shard key that has the same value for the partition key but a unique value for the row key.

### Hash-based sharding strategy

The hash-based strategy reduces the chance of hotspots, which are shards that receive a disproportionate amount of load. This strategy distributes data across shards to balance the size of each shard and the average load that each shard encounters. The sharding logic computes the shard to store an item in based on a hash of one or more attributes of the data. The chosen hashing function should distribute data evenly across the shards. The following diagram shows sharding tenant data based on a hash of tenant IDs.

To understand the advantage of the hash strategy over other sharding strategies, consider how a multitenant application that enrolls new tenants sequentially might assign the tenants to shards in the data store. When you use the range strategy, the data for tenants *1* to *n* are stored in shard A, the data for tenants *n+1* to *m* are stored in shard B, and later tenant ranges map to successive shards. If the most recently registered tenants are also the most active, most data activity occurs in a few shards, which can cause hotspots. In contrast, the hash strategy allocates tenants to shards based on a hash of their tenant ID. The hash usually distributes sequential tenants across different shards, which balances the load. The previous diagram shows this approach for tenants 55 and 56.

### Geographic sharding strategy

The geographic strategy assigns data to shards based on the geographic origin or intended consumption region of that data. In many workloads, users and the data that they generate are concentrated in specific regions. Regulatory requirements such as data residency laws might require that specific data remain within a specific jurisdiction. Even without regulatory drivers, placing data close to the users who access it most frequently reduces network latency for reads and writes.

In this strategy, you derive the shard key from a geographic attribute such as the user's country/region, the originating datacenter region, or a regional tenant identifier. You host each shard in, or pin it to, infrastructure within that geographic boundary.

For example, an application that serves customers in North America, Europe, and Asia-Pacific might maintain three shard groups, one group in each corresponding Azure region. A European application that serves only European users routes a request to the Europe shard. This approach reduces latency and meets data residency requirements.

Geographic sharding introduces the risk of uneven data distribution. If most of your users reside in one region, that region's shard carries a disproportionate share of the load and storage. You can combine geographic sharding with another strategy, such as hash or lookup, within each region to distribute load evenly across multiple shards inside the same geographic boundary.

### Advantages and considerations for each strategy

The four sharding strategies have the following advantages and considerations:

- **The lookup strategy**provides more control over shard configuration. Virtual shards reduce the impact of rebalancing because you can add new physical partitions to balance the workload. You can modify the mapping between a virtual shard and its physical partitions without affecting application code. Looking up shard locations adds overhead.
- **The range strategy**is easy to implement and works well with range queries. Range queries can retrieve multiple data items from a single shard in one operation. Data management is simpler. For example, you can schedule updates per time zone based on local load patterns when users in the same region share a shard. However, this strategy doesn't balance load evenly across shards. Rebalancing is difficult and might not resolve uneven load when most activity concentrates on adjacent shard keys.
- **The hash strategy**provides a better chance of even data and load distribution. You can route requests directly by using the hash function without maintaining a map. Computing the hash adds some overhead. Rebalancing is difficult without consistent hashing.
- **The geographic strategy**meets data residency and sovereignty requirements that other strategies don't inherently address. It reduces read and write latency when users access data in their region. However, geographic sharding can create significant data and load imbalance when user populations aren't evenly distributed across regions. Queries that span regions, such as global reporting, must retrieve data from all geographic shards and incur higher latency. Combine geographic sharding with another strategy within each region when you need both compliance and even load distribution.

Most sharding systems implement one of these approaches, but you should also consider the business requirements of your application and its data usage patterns. For example, in a multitenant application:

- You can shard data based on the workload. Segregate data for highly volatile tenants in separate shards to improve data access speed for other tenants.
- You can shard data based on tenant location. Take tenant data in a specific geographic region offline for backup and maintenance during that region's off-peak hours, while tenant data in other regions remains online during their business hours.
- Assign high-value tenants their own dedicated, lightly loaded shards. Lower-value tenants can share more densely packed shards.
- Store data for tenants that need strong data isolation and privacy on separate servers.

### Scaling and data movement operations for each strategy

Each sharding strategy provides different capabilities and levels of complexity to manage scale in, scale out, data movement, and state maintenance.

- **The lookup strategy**allows scaling and data movement operations at the user level, either online or offline. To move data:- Suspend some or all user activity, typically during off-peak periods.
- Move the data to the new virtual partition or physical shard.
- Update the mappings.
- Invalidate or refresh any caches that hold this data.
- Resume user activity.
 - You can often manage this operation centrally. The lookup strategy requires state to be highly cacheable and replica friendly.
- **The range strategy**limits scaling and data movement operations because you must split and merge data across shards, typically while part or all of the data store is offline. When you move data to rebalance shards, you might not eliminate uneven load if most activity concentrates on adjacent shard keys or data identifiers within the same range. The range strategy might also require state to map ranges to physical partitions.
- **The hash strategy**complicates scaling and data movement operations. The partition keys are hashes of the shard keys or data identifiers. With a standard hash function, such as- `hash(key) mod N`, adding or removing a shard reassigns most keys and triggers large-scale data migration. Consistent hashing reduces this impact by arranging the hash space so that only a small fraction of keys move when the shard count changes. The hash strategy doesn't require maintenance of a separate mapping state.
- **The geographic strategy**directly links scaling operations to regional infrastructure provisioning. Adding capacity in one region doesn't relieve load in another region. Regulatory requirements that mandate geographic sharding can also restrict data movement across geographic boundaries. Within each region, scaling uses the secondary strategy that distributes data across that region's shards.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- Use sharding complementary to other forms of partitioning, such as vertical partitioning and functional partitioning. For example, a single shard can contain vertically partitioned entities, and you can implement a functional partition as multiple shards. For more information, see Horizontal, vertical, and functional data partitioning.
- Keep shards balanced so that they can all handle a similar input/output (I/O) volume. Data skew accumulates over time when records are inserted and deleted, which leads to hotspots. Plan to rebalance periodically. - Rebalancing moves data between shards and often causes downtime or reduced throughput. To rebalance less frequently, use virtual partitions. Map many logical partitions to fewer physical shards. When a shard is overloaded, redistribute its virtual partitions to new physical shards without rehashing the entire dataset. Azure Cosmos DB uses this approach to decouple the partition scheme from the physical infrastructure. - Prefer many small shards over few large ones. Smaller shards migrate faster, balance load more evenly, and provide more flexibility for data redistribution.
- Use stable data for the shard key. If the shard key changes, you might need to move the corresponding data item between shards, which increases update operation overhead. Avoid basing the shard key on potentially volatile information. Choose attributes that are invariant or naturally form a key.
- Ensure that shard keys are unique. For example, avoid using autoincrementing fields as the shard key. In some systems, autoincremented fields can't coordinate across shards, which can result in items in different shards having the same shard key. - Note - Autoincremented values in other fields that aren't shard keys can also cause problems. For example, if you use autoincremented fields to generate unique IDs, two different items in different shards might be assigned the same ID.
- Shard the data to support the most frequently performed queries. You might not be able to design a shard key that matches the requirements of every query against the data. If necessary, create secondary index tables to support queries that retrieve data by attributes that aren't part of the shard key. For more information, see Index Table pattern.
- Design your shard key and data model to keep most operations scoped to a single shard. Queries that access only a single shard are more efficient than queries that retrieve data from multiple shards. Denormalize your data to keep related entities that are commonly queried together, such as customers and their orders, in the same shard to reduce the number of separate reads. - Cross-shard queries add latency, resource consumption, and complexity. When an application must retrieve data from multiple shards, use parallel fan-out queries that run against each shard concurrently and aggregate the results. Even with parallelism, the slowest shard determines overall latency. - Tip - If an entity in one shard references an entity in another shard, include the shard key for the second entity as part of the schema for the first entity. This approach can improve the performance of queries that reference related data across shards.
- Reconsider your shard key or whether sharding fits your needs if your workload requires strong transactional integrity across shard boundaries. Cross-shard transactions present challenges. Distributed coordination protocols, such as two-phase commit, add latency, introduce failure modes, and reduce throughput. Most sharded systems avoid distributed transactions and adopt eventual consistency instead. In this model, each shard updates independently, and the application handles temporary inconsistencies.
- Make sure the resources available to each shard storage node can handle the scalability requirements in terms of data size and throughput. For more information, see Data partitioning strategies.
- Consider replicating reference data to all shards. If a query against a shard also references static or slow-moving data, add this data to the shard. The application can then fetch all data for the query without making a round trip to a separate data store. - Note - If reference data held in multiple shards changes, the system must sync these changes across all shards. Some degree of inconsistency can occur while this synchronization runs. Design your applications to tolerate this inconsistency.
- Sharded systems multiply operational burden. Plan for these concerns: - **Monitoring:**You must aggregate metrics and logs across all shards to get a complete view of system health.
- **Backup and restore:**You must back up each shard independently and design restore procedures to maintain cross-shard consistency. A point-in-time restore of one shard can create inconsistencies with other shards.
- **Schema changes:**You must coordinate Data Definition Language (DDL) changes across every shard.
 - You can implement these tasks by using scripts or other automation solutions.
- You can geolocate shards to place their data near the application instances that use it. This approach can improve performance but requires extra planning for operations that must access multiple shards in different locations.

## When to use this pattern

Tip

Before you design a custom sharding layer, determine which sharding responsibilities your data platform already handles. Some services manage sharding completely. For example, Azure Cosmos DB distributes data across physical partitions, handles splits, and routes queries without application involvement. Other services manage sharding partially. For example, Azure SQL Database provides elastic database tools for shard map management and data-dependent routing, but you design the shard key and manage split operations. Use the Sharding pattern when you build and operate the sharding logic yourself.

Use this pattern when:

- The total data volume exceeds the storage capacity of a single database instance, and no vertical scaling option addresses the shortfall.
- The transaction throughput or query concurrency exceeds what a single instance can sustain, and read replicas alone don't resolve the bottleneck because write load is also high. - Note - Sharding improves the performance and scalability of a system, and it can also improve availability. A failure in one partition doesn't necessarily prevent an application from accessing data in other partitions. And an operator can perform maintenance or recovery of one partition without making all data unavailable. For more information, see Data partitioning guidance.
- Regulatory or compliance requirements mandate that specific data subsets reside in specific geographic jurisdictions, and no single-region deployment can meet all requirements.
- Distinct tenants or customer segments require physical data isolation for security, performance, or contractual reasons. - In scenarios like these, the sharding pattern is sometimes applied beyond traditional data stores. For example, a DNS zone management system could be sharded by team, environment, or region to reduce the blast radius of DNS changes and establish clear ownership boundaries. In that context, the primary motivation is operational segmentation rather than scalability. For more information, see Sharding private DNS zones.

Sharding introduces substantial and permanent complexity into your data architecture. That complexity affects development, operations, testing, query design, and failure recovery for the system's lifetime.

This pattern might not be suitable when:

- Your data volume and throughput fit within a single database instance, even with projected growth. Vertical scaling preserves query simplicity and transactional integrity.
- Your bottleneck is read volume, not write volume or storage capacity. Read replicas and caching layers can offload read traffic without the cross-shard query complexity that sharding introduces.
- Your database engine supports table-level partitioning that meets your performance needs. Partitioning within a single instance doesn't require multiple servers or routing logic.
- Your dominant query patterns require cross-entity joins, multientity transactions, or full-dataset aggregations. Sharding makes these operations expensive, and the overhead of fan-out queries and distributed coordination can outweigh the scaling benefits.

## Workload design

Evaluate how to use the Sharding pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | Data and processing are isolated to the shard, so a malfunction in one shard remains isolated to that shard. - Data partitioning - RE:07 Self-preservation |
| Cost Optimization focuses on sustaining and improvingyour workload'sreturn on investment. | A system that implements shards often benefits from using multiple instances of less expensive compute or storage resources rather than a single more expensive resource. In many cases, this configuration can save you money. - CO:07 Component costs |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | When you use sharding in your scaling strategy, data and processing are isolated to each shard, so requests only compete for resources within their assigned shard. You can also use sharding to optimize based on geography. - PE:05 Scaling and partitioning - PE:08 Data performance |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

Consider a website that surfaces an expansive collection of information about published books worldwide. The number of possible books cataloged in this workload and the typical query and usage patterns exceed what a single relational database can handle. The workload architect decides to shard the data across multiple database instances by using the books' static ISBN as the shard key. Specifically, the architect uses the check digit (0 - 10) of the ISBN, which provides 11 possible logical shards with fairly balanced data distribution.

To start, the architect colocates the 11 logical shards into three physical shard databases. In this virtual partition approach, many logical partitions map to fewer physical nodes. The architect uses the *lookup* sharding approach and stores the key-to-server mapping in a shard map database.

Azure App Service is labeled Book catalog website. It connects to multiple SQL Database instances and an Azure AI Search instance. One of the databases is labeled as the ShardMap database. It includes an example table that mirrors a part of the mapping table, which is listed later in this article. The table includes three shard databases instances: bookdbshard0, bookdbshard1, and bookdbshard2. The other databases include identical example listings of tables under them. The tables include Books, LibraryOfCongressCatalog, and an indicator of more tables. AI Search is used for faceted navigation and site search. Managed identity is associated with the App Service.

### Lookup shard map

The shard map database contains the following shard mapping table and data.

```
SELECT ShardKey, DatabaseServer
FROM BookDataShardMap
```
```
| ShardKey | DatabaseServer |
|----------|----------------|
| 0 | bookdbshard0 |
| 1 | bookdbshard0 |
| 2 | bookdbshard0 |
| 3 | bookdbshard1 |
| 4 | bookdbshard1 |
| 5 | bookdbshard1 |
| 6 | bookdbshard2 |
| 7 | bookdbshard2 |
| 8 | bookdbshard2 |
| 9 | bookdbshard0 |
| 10 | bookdbshard1 |
```
### Example website code: single shard access

The website isn't aware of how many physical shard databases exist (three in this case) or the logic that maps a shard key to a database instance. It only knows that the check digit of a book's ISBN is the shard key. The website has read-only access to the shard map database and read-write access to all shard databases. In this example, the website uses the system managed identity of its Azure App Service host for authorization, which keeps secrets out of connection strings.

The website is configured with the following connection strings either in an `appsettings.json` file, as shown in this example, or through App Service app settings.

```
{
 ...
 "ConnectionStrings": {
 "ShardMapDb": "Data Source=tcp:<database-server-name>.database.windows.net,1433;Initial Catalog=ShardMap;Authentication=Active Directory Default;App=Book Site v1.5a",
 "BookDbFragment": "Data Source=tcp:SHARD.database.windows.net,1433;Initial Catalog=Books;Authentication=Active Directory Default;App=Book Site v1.5a"
 },
 ...
}
```
The following code shows how the website runs an update query against the workload's database shard pool.

```
...
// All data for this book is stored in a shard based on the book's ISBN check digit,
// which is converted to an integer 0 - 10 (special value 'X' becomes 10).
int isbnCheckDigit = book.Isbn.CheckDigitAsInt;
// Establish a pooled connection to the database shard for this specific book.
using (SqlConnection sqlConn = await shardedDatabaseConnections.OpenShardConnectionForKeyAsync(key: isbnCheckDigit, cancellationToken))
{
 // Update the book's Library of Congress catalog information.
 SqlCommand cmd = sqlConn.CreateCommand();
 cmd.CommandText = @"UPDATE LibraryOfCongressCatalog
 SET ControlNumber = @lccn,
 ...
 Classification = @lcc
 WHERE BookID = @bookId";
 cmd.Parameters.AddWithValue("@lccn", book.LibraryOfCongress.Lccn);
 ...
 cmd.Parameters.AddWithValue("@lcc", book.LibraryOfCongress.Lcc);
 cmd.Parameters.AddWithValue("@bookId", book.Id);
 await cmd.ExecuteNonQueryAsync(cancellationToken);
}
...
```
In the previous example code, if `book.Isbn` was **978-8-1130-1024-6**, then `isbnCheckDigit` should be **6**. The `OpenShardConnectionForKeyAsync(6)` call is typically implemented by using a cache-aside approach. If cached shard information for shard key **6** isn't available, the method queries the shard map database identified by the `ShardMapDb` connection string. The method retrieves the value **bookdbshard2** from either the application cache or the shard database and substitutes it for `SHARD` in the `BookDbFragment` connection string. The method then establishes or reestablishes a pooled connection to **bookdbshard2.database.windows.net**, opens it, and returns it to the calling code. The code then updates the existing record on that database instance.

### Example website code: multiple shard access

In the rare case when the website requires a direct, cross-shard query, the application performs a parallel fan-out query across all shards.

```
...
// Retrieve all shard keys.
var shardKeys = shardedDatabaseConnections.GetAllShardKeys();
// Run the query in a fan-out style against each shard in the shard list.
Parallel.ForEachAsync(shardKeys, async (shardKey, cancellationToken) =>
{
 using (SqlConnection sqlConn = await shardedDatabaseConnections.OpenShardConnectionForKeyAsync(key: shardKey, cancellationToken))
 {
 SqlCommand cmd = sqlConn.CreateCommand();
 cmd.CommandText = @"SELECT ...
 FROM ...
 WHERE ...";
 SqlDataReader reader = await cmd.ExecuteReaderAsync(cancellationToken);
 while (await reader.ReadAsync(cancellationToken))
 {
 // Collect the results into a thread-safe data structure.
 }
 reader.Close();
 }
});
...
```
As an alternative to cross-shard queries, this workload can use an externally maintained index in Azure AI Search for site search or faceted navigation.

### Add shard instances

The workload team knows that if the data catalog or its concurrent usage grows significantly, they might require more than three database instances. The workload team doesn't expect to add database servers dynamically, and they accept workload downtime when a new shard comes online. To bring a new shard instance online, they must move data from existing shards into the new shard and update the shard map table. With this fairly static approach, the workload can confidently cache the shard key database mapping in the website code.

The shard key logic in this example has an upper limit of 11 physical shards. If the workload team determines through load estimation that they eventually require more than 11 database instances, they must make an invasive change to the shard key logic. This change involves careful planning of code modifications and data migration to the new key logic.

### SDK functionality

Instead of writing custom code for shard management and query routing to SQL Database instances, evaluate the elastic database client library. This library supports shard map management, data-dependent query routing, and cross-shard queries in both C# and Java.

## Next step

- Consistency levels in Azure Cosmos DB: Distributing data across shards introduces consistency trade-offs. This article describes the spectrum of consistency models, from strong to eventual, and their effects on availability and latency.

## Related resources

- Horizontal, vertical, and functional data partitioning: This article describes other strategies for partitioning data in the cloud to improve scalability, reduce contention, and optimize performance.
- Index Table pattern: Sometimes you can't support all queries through the design of the shard key alone. An application can use the Index Table pattern to retrieve data from a large data store by specifying a key other than the shard key.
- Materialized View pattern: To maintain the performance of some query operations, you can create materialized views that aggregate and summarize data, especially if you distribute that data across shards.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Deploy application components into a process or container separate from the main application to provide isolation and encapsulation. This pattern lets you build applications from diverse components and technologies.

Like a motorcycle sidecar, these components attach to a parent application and share its life cycle, so you create and retire them together. This pattern is also known as the *Sidekick pattern* and supports application decomposition.

## Context and problem

Applications and services often require related functionality, like monitoring, logging, configuration, and networking services. You can implement these peripheral tasks as separate components or services.

Tightly integrated components run in the same process and efficiently use shared resources, but they lack isolation. An outage in one component can affect the entire application. They also require implementation in the parent application's language, which creates interdependence.

If you decompose the application into services, you can build each service by using different languages and technologies. This approach provides more flexibility. But each component has its own dependencies and requires language-specific libraries to access the platform and shared resources. When you deploy these features as separate services, you add latency. Language-specific code and dependencies also increase complexity for hosting and deployment.

## Solution

Deploy a cohesive set of tasks alongside the primary application in a separate process or container. This approach provides a consistent interface for platform services across languages.

A sidecar service connects to the application without being part of it and deploys alongside it. Each application instance gets its own sidecar instance that shares its life cycle.

The Sidecar pattern provides the following advantages:

- **Language independence:**The sidecar runs independently from the primary application's runtime environment and programming language. You can use one sidecar implementation across applications written in different languages.
- **Shared resource access:**The sidecar can access the same resources as the primary application. For example, the sidecar can monitor system resources that both components use.
- **Low latency:**The sidecar's proximity to the primary application minimizes communication latency.
- **Enhanced extensibility:**You can extend applications that lack native extensibility mechanisms by attaching a sidecar as a separate process on the same host or subcontainer.

The most common implementation of this pattern uses containers, which are also called *sidecar containers* or *sidekick containers*.

## Problems and considerations

Consider the following points when you implement this pattern:

- Consider the deployment and packaging format to deploy services, processes, or containers. Containers work well for the Sidecar pattern.
- When you design a sidecar service, carefully choose the interprocess communication mechanism. Use language-agnostic or framework-agnostic technologies unless performance requirements make that approach impractical.
- Before you add functionality to a sidecar, evaluate whether it works better as a separate service or a traditional daemon.
- Consider whether to implement the functionality as a library or through a traditional extension mechanism. Language-specific libraries provide deeper integration and less network overhead.

## When to use this pattern

Use this pattern when:

- Your primary application uses diverse languages and frameworks. Sidecars provide a consistent interface that different applications can use regardless of their language or framework.
- A separate team or external partner owns a component.
- You must deploy a component or feature on the same host as the application.
- You need a service that shares the overall life cycle of your main application but that you can update independently.
- You need fine-grained control over resource limits for a specific resource or component. For example, you can deploy a component as a sidecar to restrict and manage its memory usage independently of the main application.

This pattern might not be suitable when:

- You need to optimize interprocess communication. Sidecars add overhead, especially latency, which makes them unsuitable for applications with frequent communication between components.
- Your application is small. The resource cost of deploying a sidecar for each instance might outweigh the isolation benefits.
- You need to scale the component independently. If you must scale the component differently from the main application, deploy it as a separate service instead.
- Your platform provides equivalent functionality. If your application platform already provides the needed capabilities natively, sidecars add unnecessary complexity.

## Workload design

Evaluate how to use the Sidecar pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Security design decisions help ensure the confidentiality,integrity, andavailabilityof your workload's data and systems. | When you encapsulate these tasks and deploy them in separate processes, you reduce the attack surface to only the necessary code. You can also use sidecars to add cross-cutting security controls to application components that lack native support for these features. - SE:04 Segmentation - SE:07 Encryption |
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | This pattern lets you flexibly integrate observability tools without adding dependencies to your application code. You can update and maintain the sidecar independently of the application. - OE:04 Tools and processes - OE:07 Monitoring system |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | This pattern lets you centralize cross-cutting tasks in sidecars that scale across multiple application instances. You don't need to deploy duplicate functionality for each application instance. - PE:07 Code and infrastructure |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

You can apply the Sidecar pattern to many scenarios. Consider the following examples:

- **Dependency abstraction:**Deploy a custom service alongside each application to provide access to shared dependency capabilities through a consistent API. This approach replaces language-specific client libraries with a sidecar that handles concerns like logging, configuration, service discovery, state management, and health checks.- The Distributed Application Runtime (Dapr) sidecar exemplifies this use case.
- **Service mesh data plane:**Deploy a sidecar proxy alongside each service instance to handle cross-cutting networking concerns like traffic routing, retries, mutual Transport Layer Security (mTLS), policy enforcement, and telemetry.- Service meshes like Istio use sidecar proxies to implement these capabilities without requiring changes to application code.
- **Ambassador sidecar:**Deploy an ambassador service as a sidecar. The application routes calls through the ambassador, which handles request logging, routing, circuit breaking, and other connectivity features.
- **Protocol adapters:**Deploy a sidecar to translate between incompatible protocols or data formats, or to bridge messaging systems. This approach lets the application use simpler or legacy interfaces.
- **Telemetry enrichment:**Deploy a sidecar to preprocess or enrich telemetry data, like metrics, logs, and traces, before it forwards the data to external monitoring systems. Components like the OpenTelemetry Collector can run as sidecars to normalize, enrich, or route telemetry separately from the application.

## Next steps

- Microservice APIs that use Dapr: Learn how Azure Container Apps uses Dapr sidecars to help you build simple, portable, resilient, and secure microservices.
- Native sidecar mode for Istio-based service mesh feature in Azure Kubernetes Service (AKS): Learn how the Istio service mesh feature for AKS uses the Sidecar pattern to address distributed architecture challenges.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Deploy static content to a cloud-based storage service that can deliver them directly to the client. This can reduce the need for potentially expensive compute instances.

## Context and problem

Web applications typically include some elements of static content. This static content might include HTML pages and other resources such as images and documents that are available to the client, either as part of an HTML page (such as inline images, style sheets, and client-side JavaScript files) or as separate downloads (such as PDF documents).

Although web servers are optimized for dynamic rendering and output caching, they still have to handle requests to download static content. This consumes processing cycles that could often be put to better use.

## Solution

In most cloud hosting environments, you can put some of an application's resources and static pages in a storage service. The storage service can serve requests for these resources, reducing load on the compute resources that handle other web requests. The cost for cloud-hosted storage is typically much less than for compute instances.

When hosting some parts of an application in a storage service, the main considerations are related to deployment of the application and to securing resources that aren't intended to be available to anonymous users.

## Issues and considerations

Consider the following points when deciding how to implement this pattern:

- The hosted storage service must expose an HTTP endpoint that users can access to download the static resources. Some storage services also support HTTPS, so it's possible to host resources in storage services that require SSL.
- For maximum performance and availability, consider using a content delivery network (CDN) to cache the contents of the storage container in multiple datacenters around the world. However, you'll likely have to pay for using the CDN.
- Storage accounts are often geo-replicated by default to provide resiliency against events that might affect a datacenter. This means that the IP address might change, but the URL will remain the same.
- When some content is located in a storage account and other content is in a hosted compute instance, it becomes more challenging to deploy and update the application. You might need to deploy the application and its content separately, and version them to manage updates more effectively, especially when the static content includes script files or UI components. But if only static resources need to be updated, you can upload them to the storage account without redeploying the application package.
- Storage services might not support the use of custom domain names. In this case it's necessary to specify the full URL of the resources in links because they'll be in a different domain from the dynamically generated content containing the links.
- The storage containers must be configured for public read access, but it's vital to ensure that they aren't configured for public write access to prevent users being able to upload content.
- Consider using a valet key or token to control access to resources that shouldn't be available anonymously. For more information, see Valet Key pattern.

## When to use this pattern

This pattern is useful for:

- Minimizing the hosting cost for websites and applications that contain some static resources.
- Minimizing the hosting cost for websites that consist of only static content and resources. Depending on the capabilities of the hosting provider's storage system, it might be possible to entirely host a fully static website in a storage account.
- Exposing static resources and content for applications running in other hosting environments or on-premises servers.
- Locating content in more than one geographical area using a content delivery network that caches the contents of the storage account in multiple datacenters around the world.
- Monitoring costs and bandwidth usage. Using a separate storage account for some or all of the static content allows the costs to be more easily separated from hosting and runtime costs.

This pattern might not be useful in the following situations:

- The application needs to perform some processing on the static content before delivering it to the client. For example, it might be necessary to add a timestamp to a document.
- The volume of static content is very small. The overhead of retrieving this content from separate storage can outweigh the cost benefit of separating it out from the compute resource.

## Workload design

An architect should evaluate how the Static Content Hosting pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Cost Optimization is focused on sustaining and improvingyour workload'sreturn on investment. | Dynamic application hosts are usually more expensive than static hosts because dynamic hosts can run your coded business logic. Using an application platform to deliver static content isn't cost-effective. - CO:09 Flow costs - CO:10 Data costs |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, code. | Offloading responsibility to an externalized host helps mitigate congestion and enables you to use your application platform only to deliver business logic. - PE:07 Code an infrastructure |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

Azure Storage supports serving static content directly from a storage container. Files are served through anonymous access requests. By default, files have a URL in a subdomain of `core.windows.net`, such as `https://contoso.z4.web.core.windows.net/image.png`. You can configure a custom domain name, and use Azure CDN to access the files over HTTPS. For more information, see Static website hosting in Azure Storage.

Static website hosting makes the files available for anonymous access. If you need to control who can access the files, you can store files in Azure blob storage and then generate time and scope limited shared access signatures to restrict access. Access signatures generated should use Microsoft Entra user delegated tokens and be short lived.

The links in the pages delivered to the client must specify the full URL of the resource. If the resource is protected with a valet key, such as a shared access signature, this signature must be included in the URL.

A sample application that demonstrates using external storage for static resources is available on GitHub. This sample uses a configuration file to specify the storage account URL and container that holds the public, static content.

The `StaticContentUrlHtmlHelper` class in the file StaticContentUrlHtmlHelper.cs exposes a method named `StaticContentUrl` that generates a URL containing the path to the cloud storage account if the URL passed to it starts with the ASP.NET root path character (~).

```
public static class StaticContentUrlHtmlHelper
{
 public static string StaticContentUrl(this HtmlHelper helper, string contentPath)
 {
 if (contentPath.StartsWith("~"))
 {
 contentPath = contentPath.Substring(1);
 }
 contentPath = string.Format("{0}/{1}", Settings.StaticContentBaseUrl.TrimEnd('/'),
 contentPath.TrimStart('/'));
 var url = new UrlHelper(helper.ViewContext.RequestContext);
 return url.Content(contentPath);
 }
}
```
The file Index.cshtml in the Views\Home folder contains an image element that uses the `StaticContentUrl` method to create the URL for its `src` attribute.

```
<img src="@Html.StaticContentUrl("~/media/orderedList1.png")" alt="Test Image" />
```
## Next steps

- Static Content Hosting sample. A sample application that demonstrates this pattern.

## Related resource

- Valet Key pattern. If the target resources aren't supposed to be available to anonymous users, use this pattern to restrict direct access.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Incrementally migrate a legacy system by gradually replacing specific pieces of functionality with new applications and services. As you replace features from the legacy system, the new system eventually comprises all of the old system's features. This approach suppresses the old system so that you can decommission it.

## Context and problem

As systems age, the development tools, hosting technology, and system architectures that they're built on can become obsolete. As new features and functionality are added, these applications become more complex, which can make them harder to maintain or extend.

It's difficult to replace an entire complex system. Instead, you can migrate to a new system gradually and use the old system for unmigrated features. However, if you run parallel versions of an application, clients must track which version contains each feature. When you migrate a feature or service, you must direct clients to the new location. To address these challenges, adopt an approach that supports incremental migration and minimizes disruptions to clients.

## Solution

After you identify new service boundaries, use an incremental process to replace specific pieces of functionality with new applications and services. Customers continue to use the same interface and are unaware that a migration is in progress.

*Download a Visio file of this architecture.*

The Strangler Fig pattern provides a controlled and phased approach to modernization. It allows the existing application to continue functioning during the modernization effort. A façade (proxy) intercepts requests that go to the back-end legacy system. The façade routes these requests either to the legacy application or to the new services.

This pattern reduces risks in migration by enabling your teams to move forward at a pace that suits the complexity of the project. As you migrate functionality to the new system, the legacy system becomes obsolete, and you decommission the legacy system.

- The Strangler Fig pattern begins by introducing a façade (proxy) between the client app, the legacy system, and the new system. The façade acts as an intermediary. It allows the client app to interact with the legacy system and the new system. Initially, the façade routes most requests to the legacy system.
- As the migration progresses, the façade incrementally shifts requests from the legacy system to the new system. With each iteration, you implement more pieces of functionality in the new system. - This incremental approach gradually reduces the legacy system's responsibilities and expands the scope of the new system. The process is iterative. It allows the team to address complexities and dependencies in manageable stages. These stages help the system remain stable and functional.
- After you migrate all of the functionality and there are no dependencies on the legacy system, you can decommission the legacy system. The façade routes all requests exclusively to the new system.
- You remove the façade and reconfigure the client app to communicate directly with the new system. This step marks the completion of the migration.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- Consider how to handle services and data stores that both the new system and the legacy system might use. Make sure that both systems can access these resources at the same time.
- Structure new applications and services so that you can easily intercept and replace them in future strangler fig migrations. For example, strive to have clear demarcations between parts of your solution so that you can migrate each part individually.
- After the migration is complete, you typically remove the strangler fig façade. Alternatively, you can maintain the façade as an adapter for legacy clients to use while you update the core system for newer clients. - Conceptualize this as transitional architecture, and balance this architecture's risk mitigation benefits against its temporary infrastructural costs.
- Make sure that the façade keeps up with the migration.
- Make sure that the façade doesn't become a single point of failure or a performance bottleneck.
- Plan for cross-system dependencies. During migration, both systems need to coexist and communicate. For example, the new system might need to call unmigrated functionality from the legacy system, and unmigrated legacy components might need to call migrated functionality from the new system. To manage these calls, use the Anti-corruption Layer pattern. An anti-corruption layer acts as an adapter that translates requests between the two systems. This layer protects the new system's design from legacy semantics so that the legacy system can reach new services without significant code changes. Without this adapter, cross-system dependencies can break components or force the new system to adopt legacy conventions.

## When to use this pattern

Use this pattern when:

- You gradually migrate a back-end application to a new architecture, especially when replacing large systems, key components, or complex features introduces risk.
- The original system can continue to exist for an extended period of time during the migration effort.

This pattern might not be suitable when:

- Requests to the back-end system can't be intercepted.
- You can't access the legacy system's source code. To disable migrated features and redirect internal calls, you need to be able to modify the legacy system's source code.
- You migrate a small system and replacing the whole system is simple.
- You need to fully decommission the original solution quickly.

## Workload design

Evaluate how to use the Strangler Fig pattern in a workload's design to address the goals and principles of the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and to ensure that itrecoversto a fully functioning state after a failure occurs. | This pattern's incremental approach can help mitigate risks during a component transition compared to making large systemic changes all at once. - RE:08 Testing |
| Cost Optimization focuses on sustaining and improvingyour workload'sreturn on investment (ROI). | The goal of this approach is to maximize the use of existing investments in the currently running system while modernizing incrementally. It enables you to perform high-ROI replacements before low-ROI replacements. - CO:07 Component costs - CO:08 Environment costs |
| Operational Excellence helps deliver workload qualitythroughstandardized processesand team cohesion. | This pattern provides a continuous improvement approach. Incremental replacements that make small changes over time are preferable to large systemic changes that are riskier to implement. - OE:06 Supply chain for workload development - OE:11 Safe deployment practices |

Consider any trade-offs against the goals of the other pillars that this pattern might introduce.

## Example

Legacy systems typically depend on a centralized monolithic database that serves multiple domains. Over time, this shared database becomes difficult to manage and improve because of its cross-domain dependencies. To address this challenge, the Strangler Fig pattern incrementally extracts domain-specific tables, stored procedures, and related data from the monolithic database into isolated domain databases. Each database contains only one domain. Repeat the extraction process until the monolithic database is fully decomposed.

- Introduce a new system service, which starts to manage requests for its domain. The new system service still reads from and writes to the monolithic database for its domain tables. The legacy system continues to serve all other domains.
- Introduce an isolated domain database for the new system. Migrate the relevant domain tables and their historical data to the new database by using an extract, transform, and load (ETL) process. A change data capture (CDC) process syncs the domain data from the monolithic database to the new domain database. During this phase, the legacy system continues to read from and write to the monolithic database, and the new system writes to the new domain database. Validate consistency between both databases before cutover.
- After validation, the new domain database is the system of record for that domain. The new system performs all read and write operations against the domain database. Remove the corresponding domain tables, stored procedures, and dependencies from the monolithic database. Repeat this process for each domain until the monolithic database is fully decomposed. - You can roll back to the monolithic database during phase 2 and at the start of phase 3, when the domain tables and synchronization processes still exist in the monolithic database. To roll back to the monolithic database after you remove the domain tables, stored procedures, and synchronization processes from the monolithic database, you must restore those objects and replay data changes. However, this process significantly increases effort and risk. Treat the removal of legacy objects as a deliberate final step for each domain. Remove legacy objects only after the new system is validated.

## Contributors

*Microsoft maintains this article. The following contributors wrote this article.*

Principal authors:

- Adnan Khan | Senior Cloud Solutions Architect
- Ovais Mehboob Ahmed Khan | Senior Cloud Solution Architect

*To see nonpublic LinkedIn profiles, sign in to LinkedIn.*

## Next step

- Read Martin Fowler's blog post about Strangler Fig pattern application

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Limit the resources that an application instance, an individual tenant, or an entire service can consume. This lets the system function and meets its service-level objectives (SLOs) under sudden or sustained load.

## Context and problem

The load on a cloud application varies over time based on active users and their activity. More users sign in during business hours, and the system runs computationally expensive analytics at the end of each month. Sudden bursts also occur. If processing demand exceeds available capacity, the system slows or fails. When the system has an agreed service level, that failure violates the SLO.

Several strategies handle varying load, depending on the application's business goals. One strategy is autoscaling, which matches provisioned resources to current demand and controls cost. But provisioning new resources takes time and adds cost. Demand that exceeds capacity growth or budget creates a resource deficit.

## Solution

An alternative to autoscaling is to cap resource use and throttle requests when usage exceeds that cap. The workload monitors its own resource use and throttles requests from one or more users when usage exceeds the threshold. The system continues to function and meet its SLOs.

Throttling is a control loop, not a single admission decision. The system needs low-latency signals at three layers: infrastructure utilization, application state, and per-principal counters. It continuously measures saturation, enforces limits at well-defined boundaries, and adapts those limits as traffic patterns change. Overload is a normal operating mode that a mature system detects and recovers from. Throttling provides self-preservation capabilities in your workload.

The system can implement several throttling or related strategies:

- **Per-principal rate limits:**Reject requests from a user who already exceeded the configured rate over a defined window. This strategy requires the system to attribute each request to a principal and meter resource use against that principal. For multitenant workloads, see Measure the consumption of each tenant.
- **Graceful feature degradation:**Turn off or degrade nonessential features so that essential features have enough resources. This strategy trades response completeness for availability. For example, a video-streaming application can drop to a lower resolution.
- **Load leveling:**Smooth activity volume by using a queue. In a multitenant environment, leveling reduces performance for every tenant. When tenants have different service-level agreements (SLAs), process work for high-value tenants immediately and hold lower-priority work until the backlog eases. Implement this approach by using the Priority Queue pattern or by exposing separate endpoints for each priority tier.
- **Priority-based deferral:**Defer operations on behalf of lower-priority applications or tenants. Suspend or limit operations, and return an exception that tells the tenant to retry later.
- **Outbound rate limits:**Limit your own outbound calls when an external dependency fails or returns errors. Lower the in-flight request count to stop flooding logs and to avoid retry costs against an unhealthy dependency. Restore normal request flow after the dependency recovers. For example, NServiceBus implements this functionality.

The following chart shows resource use (a combination of memory, CPU, bandwidth, and other factors) over time for an application that uses three features, labeled A, B, and C. A feature is a specific area of functionality, such as a component that performs a specific set of tasks, a piece of code that performs a complex calculation, or an element that provides a service such as an in-memory cache.

A line graph plots resource utilization on the y-axis against time on the x-axis. Three colored lines represent Feature A, Feature B, and Feature C, with Feature A's line lowest, Feature B's line in the middle, and Feature C's line highest. A solid horizontal line near the top of the chart marks maximum capacity, and a dashed horizontal line below it marks the soft limit of resource utilization. Two vertical dashed lines mark times T1 and T2. Before T1, all three feature lines fluctuate, and Feature C's line rises and crosses the soft limit. At T1, Feature B's line drops to zero and stays at zero until T2 because Feature B is suspended to free resources for Feature A and Feature C. Feature C's line falls back below the soft limit between T1 and T2 while Feature A continues normally. At T2, Feature B resumes and all three lines continue to fluctuate below the soft limit.

The chart is a stacked area chart. The area below Feature A's line shows the resources that Feature A consumes, the area between Feature A's and Feature B's lines shows the resources that Feature B consumes, and the area between Feature B's and Feature C's lines shows the resources that Feature C consumes. Feature C's line sits at the top of the stack, so it also shows total system resource use over time.

The chart shows graceful feature degradation. Just before time T1, total resource use approaches the threshold and risks exhausting available capacity. Feature B is less critical than Feature A or Feature C, so the system turns off Feature B and releases its resources. Between times T1 and T2, Feature A and Feature C continue normally. By time T2, total resource use drops enough to turn Feature B back on.

You can combine autoscaling, graceful degradation, and throttling to keep applications responsive and within SLAs. When you expect demand to stay high, throttling maintains stability while the system scales out. After scaling completes, the system restores full functionality.

The next chart shows total resource use over time and how throttling combines with autoscaling and other compensating controls.

A line graph plots resource utilization for all applications on the y-axis against time on the x-axis. Two horizontal reference lines mark the soft limit of resource utilization and the maximum capacity before autoscaling. A higher horizontal line, which begins at time T2, marks the maximum capacity after autoscaling. The utilization line rises and fluctuates over time. It crosses the soft limit at time T1, which is the point where autoscaling commences. Between T1 and T2, the system is throttled while autoscaling occurs, and utilization stays below the preautoscaling maximum capacity. At time T2, autoscaling completes, throttling is relaxed, and the utilization line jumps up and continues to fluctuate below the new, higher maximum capacity.

At time T1, the system reaches the soft limit and starts to scale out. If new resources don't arrive in time, demand can exhaust the existing resources, and the system can fail. Throttling rejects excess requests during scale-out to keep resource use below the hard limit, then lifts those restrictions after new capacity comes online.

Tip

Edge controls and the Throttling pattern address different problems. Edge controls, such as Azure DDoS Protection and web application firewall (WAF) rate-limit rules, run at the network boundary and drop volumetric or malicious traffic before it reaches your application. The Throttling pattern runs inside your application and meters *legitimate* traffic against application-defined limits. Use both layers together. DDoS protection doesn't stop a legitimate user from overloading your service, and application throttling doesn't absorb a volumetric attack.

## Problems and considerations

Consider the following points as you decide how to implement this pattern:

- Make throttling decisions early. Throttling is an architectural decision that affects the whole system. Retrofitting it later is expensive.
- Align throttling limits with the component that saturates first. - Request rate is the most familiar dimension to limit, but the real bottleneck is often concurrent in-flight requests, queue depth, CPU or memory utilization, or a downstream dependency's own limits. A requests-per-second limit doesn't protect a system whose bottleneck is concurrency at a fan-out point. - At each throttling enforcement boundary, such as the gateway, the service, a partition, or a downstream dependency, identify what saturates first and set the limit on that dimension. For concurrency-bounded protection at fan-out points, see the Bulkhead pattern, which complements throttling.
- Pick a limiting algorithm intentionally. Match it to the tolerance of the component that you're protecting. - Algorithm - Behavior and best fit - Token bucket - Supports bursts up to a configured size and enforces a steady refill rate. Use for gateways that need to absorb short spikes. - Leaky bucket - Emits at a constant rate. Use for back ends that need a steady ingress rate. - Fixed window - Simple to implement, but admits back-to-back bursts at window boundaries. - Sliding window - Smooths the window-boundary problem of fixed windows at the cost of more state.
- Decide who the limit affects. Throttling at a coarse boundary, such as a regional gateway, can affect many unrelated users when only a few of them drive the load.
- Decide where the counter resides when one limit spans multiple nodes. Local counters are fast but undercount when the same caller reaches multiple replicas. A centralized counter in a shared store like Redis sees every request but adds latency to each decision. To approximate a global rate, divide the limit across replicas and reconcile periodically.
- Make throttling decisions quickly. The system must detect rising load, react, and return to normal after load eases. This process requires continuous performance instrumentation.
- Shed load proactively, not at the edge of collapse. A throttle that only rejects after a component saturates causes latency to spike before callers see any back-pressure. - As utilization approaches the hard limit, start rejecting a growing fraction of requests. Early rejection signals callers to back off and prevents the latency collapse that abrupt limits often trigger. Use p99 latency against your SLO as the primary trigger. Average utilization can look healthy while p99 has already breached. - Where you can distinguish request value, shed lower-value or more retryable work first. For more information, see the Priority Queue pattern.
- Return a status code that tells the client when a temporary rejection is the result of throttling: - **HTTP 429 (Too Many Requests):**The caller exceeds a configured request rate over a defined window.
- **HTTP 503 (Service Unavailable):**The service can't handle the request right now, often because of an unexpected load spike.
 - Include a - `Retry-After`HTTP header so that the client can pick a retry strategy. Return enough context for the caller to retry deliberately instead of guessing. For example, name the limit that the caller exceeds, clarify the affected scope, or suggest a rate that would succeed. Unexplained rejections don't help callers adapt.
- Propagate overload signals from your dependencies instead of absorbing them. A service that throttles its callers must also honor the throttling responses that it receives from its own downstream dependencies. If your service hides a downstream 429 or 503 response by retrying silently or by returning a generic HTTP 500 (Internal Server Error) response, callers can't slow down, retries amplify, and the overload cascades back upstream. The Retry Storm antipattern describes this failure mode. Surface back-pressure to upstream callers so that the entire call chain sheds load together.
- Make rejection cheaper than the work that it prevents. If refusing a request requires heavy authentication, deep parsing, or complex policy evaluation, a flood of rejected requests can still saturate the system. Reject as early in the request pipeline as possible, and load test the rejection path itself.
- Plan for cases where throttling can't buy enough time for autoscale. If demand grows faster than new capacity comes online, even a throttled system can fail. When that outcome is unacceptable, keep larger capacity reserves and configure more aggressive autoscaling.
- Don't use caching as a substitute for throttling. A cache lowers average load on the origin but doesn't bound peak load. Every cache miss passes through to the origin, and when a popular key expires under heavy traffic, many callers can race to refill it. Use caching to reduce normal pressure and throttling to bound the worst case. For more information, see the Cache-Aside pattern.
- Normalize resource costs for different operations because they generally don't carry equal execution costs. For example, throttling limits might be higher for read operations and lower for write operations. Ignoring per-operation cost can exhaust capacity and create an attack vector.
- Make throttling configuration changeable at runtime. When abnormal load arrives, you need to adjust limits without a deployment. Deployments are slow and risky during an incident. The External Configuration Store pattern externalizes the configuration so that you can change it at runtime.
- Consider adaptive limits instead of static limits. Some throttling SDKs react to latency or queue-depth signals so that the limit tracks actual component conditions. Always pair an adaptive limiter with a set maximum.
- Revisit your limits as the workload evolves. Adaptive limiters can't track every kind of drift, such as SLO changes, changes in dependency capacity, or shifts in per-operation cost. Schedule periodic operator review against those inputs.

## When to use this pattern

Use this pattern:

- To keep a system within its SLOs.
- To prevent a single tenant from monopolizing application resources.
- To handle bursts in activity.
- To limit the maximum resource level that a system needs.
- To reduce low-value compute during periods of high grid carbon intensity.

## Workload design

Evaluate how to use the Throttling pattern in a workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. The following table provides guidance about how this pattern supports the goals of each pillar.

| Pillar | How this pattern supports pillar goals |
|---|---|
| Reliability design decisions help your workload become resilientto malfunction and ensure that itrecoversto a fully functioning state after a failure occurs. | You design the limits to help prevent resource exhaustion that might lead to malfunctions. You can also use this pattern as a control mechanism in a graceful degradation plan. - RE:07 Self-preservation |
| Security design decisions help ensure the confidentiality,integrity, andavailabilityof your workload's data and systems. | You can design the limits to help prevent resource exhaustion that could result from automated abuse of the system. - SE:06 Network controls - SE:08 Hardening resources |
| Cost Optimization focuses on sustaining and improvingyour workload'sreturn on investment. | The enforced limits can inform cost modeling and can be directly tied to the business model of your application. They also put clear upper bounds on utilization, which can be factored into resource sizing. - CO:02 Cost model - CO:12 Scaling costs |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, and code. | When the system is under high demand, this pattern helps mitigate congestion that can lead to performance bottlenecks. You can also use it to proactively avoid noisy neighbor scenarios. - PE:02 Capacity planning - PE:05 Scaling and partitioning |

If this pattern introduces trade-offs within a pillar, consider them against the goals of the other pillars.

## Example

The following diagram shows throttling in a multitenant system.

Three labeled users on the left represent tenants of a multitenant Surveys application: Adatum, Fabrikam, and Contoso. Each user sends requests through a tenant-specific custom domain, which the application uses to identify the tenant. Adatum sends 5 requests per second through surveys.adatum.com, Fabrikam sends 10 requests per second through surveys.fabrikam.com, and Contoso sends 150 requests per second through surveys.contoso.com. On the right, the surveys application web role meters the per-second request rate for each tenant. The Adatum and Fabrikam request flows pass through to the application. The Contoso request flow is blocked by an Error: Throttled response because the rate exceeds the per-tenant limit.

Users from several tenant organizations access a cloud-hosted application to fill out and submit surveys. The application contains instrumentation that monitors the rate at which each tenant's users submit requests.

To prevent users from one tenant degrading responsiveness and availability for users in other tenants, the application limits the requests-per-second rate that any single tenant can submit. The application blocks requests that exceed this limit.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

Use a token that provides clients with restricted direct access to a specific resource, in order to offload data transfer from the application. This is particularly useful in applications that use cloud-hosted storage systems or queues, and can minimize cost and maximize scalability and performance.

## Context and problem

Client programs and web browsers often need to read and write files or data streams to and from an application's storage. Typically, the application will handle the movement of the data — either by fetching it from storage and streaming it to the client, or by reading the uploaded stream from the client and storing it in the data store. However, this approach absorbs valuable resources such as compute, memory, and bandwidth.

Data stores have the ability to handle upload and download of data directly, without requiring that the application perform any processing to move this data. But, this typically requires the client to have access to the security credentials for the store. This can be a useful technique to minimize data transfer costs and the requirement to scale out the application, and to maximize performance. It means, though, that the application is no longer able to manage the security of the data. After the client has a connection to the data store for direct access, the application can't act as the gatekeeper. It's no longer in control of the process and can't prevent subsequent uploads or downloads from the data store.

This isn't a realistic approach in distributed systems that need to serve untrusted clients. Instead, applications must be able to securely control access to data in a granular way, but still reduce the load on the server by setting up this connection and then allowing the client to communicate directly with the data store to perform the required read or write operations.

## Solution

You need to resolve the problem of controlling access to a data store where the store can't manage authentication and authorization of clients. One typical solution is to restrict access to the data store's public connection and provide the client with a key or token that the data store can validate.

This key or token is usually referred to as a valet key. It provides time-limited access to specific resources and allows only predefined operations with granular control such as writing to storage but not reading, or uploading and downloading in a web browser. Applications can create and issue valet keys to client devices and web browsers quickly and easily, allowing clients to perform the required operations without requiring the application to directly handle the data transfer. This removes the processing overhead, and the impact on performance and scalability, from the application and the server.

The client uses this token to access a specific resource in the data store for only a specific period, and with specific restrictions on access permissions, as shown in the figure. After the specified period, the key becomes invalid and won't allow access to the resource.

Diagram showing an example of the workflow for a system using the valet key pattern. Step 1 shows the user requesting the target resource. Step 2 shows the valet key application checking the validity of the request and generating an access token. Step 3 shows the token being returned to the user. Step 4 shows the user accessing the target resource using the token.

It's also possible to configure a key that has other dependencies, such as the scope of the data. For example, depending on the data store capabilities, the key can specify a complete table in a data store, or only specific rows in a table. In cloud storage systems the key can specify a container, or just a specific item within a container.

The key can also be invalidated by the application. This is a useful approach if the client notifies the server that the data transfer operation is complete. The server can then invalidate that key to prevent further access.

Using this pattern can simplify managing access to resources because there's no requirement to create and authenticate a user, grant permissions, and then remove the user again, or worse leave that permission as a standing permission. It also lets you limit the location, the permission, and the validity period by generating a key at runtime. The important factors are to limit the validity period, and especially the location of the resource, as tightly as possible so that the recipient can only use it for the intended purpose.

## Issues and considerations

Consider the following points when deciding how to implement this pattern:

**Manage the validity status and period of the key**. If leaked or compromised, the key effectively unlocks the target item and makes it available for malicious use during the validity period. A key can usually be revoked or disabled, depending on how it was issued. Server-side policies can be changed or, the server key it was signed with can be invalidated. Specify a short validity period to minimize the risk of allowing unauthorized operations to take place against the data store. However, if the validity period is too short, the client might not be able to complete the operation before the key expires. Allow authorized users to renew the key before the validity period expires if multiple accesses to the protected resource are required.

**Control the level of access the key will provide**. Typically, the key should allow the user to only perform the actions necessary to complete the operation, such as read-only access if the client shouldn't be able to upload data to the data store. For file uploads, it's common to specify a key that provides write-only permission, as well as the location and the validity period. It's critical to accurately specify the resource or the set of resources to which the key applies.

**Consider how to control users' behavior**. Implementing this pattern means some loss of control over the resources users are granted access to. The level of control that can be exerted is limited by the capabilities of the policies and permissions available for the service or the target data store. For example, it's usually not possible to create a key that limits the size of the data to be written to storage, or the number of times the key can be used to access a file. This can result in huge unexpected costs for data transfer, even when used by the intended client, and might be caused by an error in the code that causes repeated upload or download. To limit the number of times a file can be uploaded, where possible, force the client to notify the application when one operation has completed. For example, some data stores raise events the application code can use to monitor operations and control user behavior. However, it's hard to enforce quotas for individual users in a multitenant scenario where the same key is used by all the users from one tenant. Granting users *create* permissions can help you control the amount of data being updated by making tokens effectively single-use. The *create* permission doesn't allow overwrites, so each token can only be used for one write activity.

**Validate, and optionally sanitize, all uploaded data**. A malicious user that gains access to the key could upload data designed to compromise the system. Alternatively, authorized users might upload data that's invalid and, when processed, could result in an error or system failure. To protect against this, ensure that all uploaded data is validated and checked for malicious content before use.

**Audit all operations**. Many key-based mechanisms can log operations such as uploads, downloads, and failures. These logs can usually be incorporated into an audit process, and also used for billing if the user is charged based on file size or data volume. Use the logs to detect authentication failures that might be caused by issues with the key provider, or accidental removal of a stored access policy.

**Deliver the key securely**. It can be embedded in a URL that the user activates in a web page, or it can be used in a server redirection operation so that the download occurs automatically. Always use HTTPS to deliver the key over a secure channel.

**Protect sensitive data in transit**. Sensitive data delivered through the application will usually take place using TLS, and this should be enforced for clients accessing the data store directly.

Other issues to be aware of when implementing this pattern are:

- If the client doesn't, or can't, notify the server of completion of the operation, and the only limit is the expiration period of the key, the application won't be able to perform auditing operations such as counting the number of uploads or downloads, or preventing multiple uploads or downloads.
- The flexibility of key policies that can be generated might be limited. For example, some mechanisms only allow the use of a timed expiration period. Others aren't able to specify a sufficient granularity of read/write permissions.
- If the start time for the key or token validity period is specified, ensure that it's a little earlier than the current server time to allow for client clocks that might be slightly out of synchronization. The default, if not specified, is usually the current server time.
- The URL containing the key might be recorded in server log files. While the key will typically have expired before the log files are used for analysis, ensure that you limit access to them. If log data is transmitted to a monitoring system or stored in another location, consider implementing a delay to prevent leakage of keys until after their validity period has expired.
- If the client code runs in a web browser, the browser might need to support cross-origin resource sharing (CORS) to enable code that executes within the web browser to access data in a different domain from the one that served the page. Some older browsers and some data stores don't support CORS, and code that runs in these browsers might not be able to use a valet key to provide access to data in a different domain, such as a cloud storage account.
- While the client doesn't need to have pre-configured authentication for to the end resource, the client does need to preestablish means of authentication to the valet key service.
- Keys should only be handed out to authenticated clients with proper authorization.
- The generation of access tokens is privileged action, so the valet key service must be secured with strict access policies. The service might allow access to sensitive systems by third parties, making the security of this service of particular importance.

## When to use this pattern

This pattern is useful for the following situations:

- To minimize resource loading and maximize performance and scalability. Using a valet key doesn't require the resource to be locked, no remote server call is required, there's no limit on the number of valet keys that can be issued, and it avoids a single point of failure resulting from performing the data transfer through the application code. Creating a valet key is typically a simple cryptographic operation of signing a string with a key.
- To minimize operational cost. Enabling direct access to stores and queues is resource and cost efficient, can result in fewer network round trips, and might allow for a reduction in the number of compute resources required.
- When clients regularly upload or download data, particularly where there's a large volume or when each operation involves large files.
- When the application has limited compute resources available, either due to hosting limitations or cost considerations. In this scenario, the pattern is even more helpful if there are many concurrent data uploads or downloads because it relieves the application from handling the data transfer.
- When the data is stored in a remote data store or a different region. If the application was required to act as a gatekeeper, there might be a charge for the additional bandwidth of transferring the data between regions, or across public or private networks between the client and the application, and then between the application and the data store.

This pattern might not be useful in the following situations:

- If clients can already uniquely authenticate to your backend service, with Azure role-based access control (Azure RBAC) for example, don't use this pattern.
- If the application must perform some task on the data before it's stored or before it's sent to the client. For example, if the application needs to perform validation, log access success, or execute a transformation on the data. However, some data stores and clients are able to negotiate and carry out simple transformations such as compression and decompression (for example, a web browser can usually handle gzip formats).
- If the design of an existing application makes it difficult to incorporate the pattern. Using this pattern typically requires a different architectural approach for delivering and receiving data.
- If it's necessary to maintain audit trails or control the number of times a data transfer operation is executed, and the valet key mechanism in use doesn't support notifications that the server can use to manage these operations.
- If it's necessary to limit the size of the data, especially during upload operations. The only solution to this is for the application to check the data size after the operation is complete, or check the size of uploads after a specified period or on a scheduled basis.

## Workload design

An architect should evaluate how the Valet Key pattern can be used in their workload's design to address the goals and principles covered in the Azure Well-Architected Framework pillars. For example:

| Pillar | How this pattern supports pillar goals |
|---|---|
| Security design decisions help ensure the confidentiality,integrity, andavailabilityof your workload's data and systems. | This pattern enables a client to directly access a resource without needing long-lasting or standing credentials. All access requests start with an auditable transaction. The granted access is then limited in both scope and duration. This pattern also makes it easier to revoke the granted access. - SE:05 Identity and access management |
| Cost Optimization is focused on sustaining and improvingyour workload'sreturn on investment. | This design offloads processing as an exclusive relationship between the client and the resource without adding a component to directly handle all client requests. The benefit is most dramatic when client requests are frequent or large enough to require significant proxy resources. - CO:09 Flow costs |
| Performance Efficiency helps your workload efficiently meet demandsthrough optimizations in scaling, data, code. | Not using an intermediary resource to proxy the access offloads processing as an exclusive relationship between the client and the resource without requiring an ambassador component that needs to handle all client requests in a performant way. The benefit of using this pattern is most significant when the proxy doesn't add value to the transaction. - PE:07 Code and infrastructure |

As with any design decision, consider any tradeoffs against the goals of the other pillars that might be introduced with this pattern.

## Example

Azure supports shared access signatures on Azure Storage for granular access control to data in blobs, tables, and queues, and for Service Bus queues and topics. A shared access signature token can be configured to provide specific access rights such as read, write, update, and delete to a specific table; a key range within a table; a queue; a blob; or a blob container. The validity can be a specified time period. This functionality is well suited for using a valet key for access.

Consider a workload that has hundreds of mobile or desktop clients frequently uploading large binaries. Without this pattern, the workload has essentially two options. The first is to provide standing access and configuration to all the clients to perform uploads directly to a storage account. The other is to implement the Gateway Routing pattern to set up an endpoint where clients use proxied access to storage, but this might not be adding additional value to the transaction. Both approaches have problems addressed in the pattern context:

- Long lived, preshared secrets. Potentially without much way to provide different keys to different clients.
- Added expense for running a compute service that has sufficient resources to deal with currently receiving large files.
- Potentially slowing up client interactions by adding an extra layer of compute and network hop to the upload process.

Using the Valet Key pattern addresses the security, cost optimization, and performance concerns.

- Clients, at the last responsible moment, authenticate to a light weight, scale-to-zero Azure Function hosted API to request access.
- The API validates the request and then obtains and returns a time & scope limited SaS token. - The token generated by the API restricts the client to the following limitations: - Which storage account to use. Meaning, the client doesn't need to know this information ahead of time.
- A specific container and filename to use; ensuring that the token can be used with, at most, one file.
- A short, window of operation, such as three minutes. This short time period ensures that tokens have a TTL that doesn't extend past its utility.
- Permissions to only *create*a blob; not download, update, or delete.

- That token is then used by the client, within the narrow time window, to upload the file directly to the storage account.

The API generates these tokens to authorized clients using a *user delegation key* based on the API's own Microsoft Entra ID managed identity. Logging is enabled on both the storage account(s) and the token generation API allow correlation between token requests and token usage. The API can use client authentication information or other data available to it to decide which storage account or container to use, such as in a multitenant situation.

A complete sample is available on GitHub at Valet Key pattern example. The following code snippets are adapted from that example. This first snippet shows how the Azure Function (found in **ValetKey.Web**) generates a user delegated shared access signature token using the Azure Function's own managed identity.

```
[Function("FileServices")]
public async Task<StorageEntitySas> GenerateTokenAsync([HttpTrigger(...)] HttpRequestData req, ...,
 CancellationToken cancellationToken)
{
 // Authorize the caller, select a blob storage account, container, and file name.
 // Authenticate to the storage account with the Azure Function's managed identity.
 ...
 return await GetSharedAccessReferenceForUploadAsync(blobContainerClient, blobName, cancellationToken);
}
/// <summary>
/// Return an access key that allows the caller to upload a blob to this
/// specific destination for about three minutes.
/// </summary>
private async Task<StorageEntitySas> GetSharedAccessReferenceForUploadAsync(BlobContainerClient blobContainerClient,
 string blobName,
 CancellationToken cancellationToken)
{
 var blobServiceClient = blobContainerClient.GetParentBlobServiceClient();
 var blobClient = blobContainerClient.GetBlockBlobClient(blobName);
 // Allows generating a SaS token that is evaluated as the union of the RBAC permissions on the managed identity
 // (for example, Blob Data Contributor) and then narrowed further by the specific permissions in the SaS token.
 var userDelegationKey = await blobServiceClient.GetUserDelegationKeyAsync(DateTimeOffset.UtcNow.AddMinutes(-3),
 DateTimeOffset.UtcNow.AddMinutes(3),
 cancellationToken);
 // Limit the scope of this SaS token to the following parameters:
 var blobSasBuilder = new BlobSasBuilder
 {
 BlobContainerName = blobContainerClient.Name, // - Specific container
 BlobName = blobClient.Name, // - Specific filename
 Resource = "b", // - Blob only
 StartsOn = DateTimeOffset.UtcNow.AddMinutes(-3), // - For about three minutes (+/- for clock drift)
 ExpiresOn = DateTimeOffset.UtcNow.AddMinutes(3), // - For about three minutes (+/- for clock drift)
 Protocol = SasProtocol.Https // - Over HTTPS
 };
 blobSasBuilder.SetPermissions(BlobSasPermissions.Create);
 return new StorageEntitySas
 {
 BlobUri = blobClient.Uri,
 Signature = blobSasBuilder.ToSasQueryParameters(userDelegationKey, blobServiceClient.AccountName).ToString();
 };
}
```
The following snippet is the data transfer object (DTO) used by both the API and the client.

```
public class StorageEntitySas
{
 public Uri? BlobUri { get; internal set; }
 public string? Signature { get; internal set; }
}
```
The client (found in **ValetKey.Client**) then uses the URI and token returned from the API to perform the upload without requiring additional resources and at full client-to-storage performance.

```
...
// Get the SaS token (valet key)
var blobSas = await httpClient.GetFromJsonAsync<StorageEntitySas>(tokenServiceEndpoint);
var sasUri = new UriBuilder(blobSas.BlobUri)
{
 Query = blobSas.Signature
};
// Create a blob client using the SaS token as credentials
var blob = new BlobClient(sasUri.Uri);
// Upload the file directly to blob storage
using (var stream = await GetFileToUploadAsync(cancellationToken))
{
 await blob.UploadAsync(stream, cancellationToken);
}
...
```
## Next steps

The following guidance might be relevant when implementing this pattern:

- The implementation of the example is available on GitHub at Valet Key pattern example.
- Grant limited access to Azure Storage resources using shared access signatures (SAS)
- Shared Access Signature Authentication with Service Bus

## Related resources

The following patterns might also be relevant when implementing this pattern:

- Gatekeeper pattern. This pattern can be used in conjunction with the Valet Key pattern to protect applications and services by using a dedicated host instance that acts as a broker between clients and the application or service. The gatekeeper validates and sanitizes requests, and passes requests and data between the client and the application. Can provide an additional layer of security, and reduce the attack surface of the system.
- Static Content Hosting pattern. Describes how to deploy static resources to a cloud-based storage service that can deliver these resources directly to the client to reduce the requirement for expensive compute instances. Where the resources aren't intended to be publicly available, the Valet Key pattern can be used to secure them.
