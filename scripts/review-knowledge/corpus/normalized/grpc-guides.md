##### Authentication

An overview of gRPC authentication, including built-in auth mechanisms, and how to plug in your own authentication systems.

The documentation covers the following techniques:

An overview of gRPC authentication, including built-in auth mechanisms, and how to plug in your own authentication systems.

gRPC is designed to support high-performance open-source RPCs in many languages. This page describes performance benchmarking tools, scenarios considered by tests, and the testing infrastructure.

Explains how and when to cancel RPCs.

How to compress the data sent over the wire while using gRPC.

A mechanism in the gRPC library that allows users to inject custom metrics at the gRPC server and consume at gRPC clients to make your custom load balancing algorithms.

Explains how custom load balancing policies can help optimize load balancing under unique circumstances.

Explains standard name resolution, the custom name resolver interface, and how to write an implementation.

Explains how deadlines can be used to effectively deal with unreliable backends.

Explains the debugging process of gRPC applications using grpcdebug

How gRPC deals with errors, and gRPC error codes.

Explains what flow control is and how you can manually control it.

Explains how to gracefully shut down a gRPC server to avoid causing RPC failures for connected clients.

Explains how gRPC servers expose a health checking service and how client can be configured to automatically check the health of the server it is connecting to.

Explains how interceptors can be used for implementing generic behavior that applies to many RPC methods.

How to use HTTP/2 PING-based keepalives in gRPC.

Explains what metadata is, how it is transmitted, and what it is used for.

OpenTelemetry Metrics available in gRPC

A user guide of both general and language-specific best practices to improve performance.

Explains how reflection can be used to improve the transparency and interpretability of RPCs.

Explains what request hedging is and how you can configure it.

gRPC takes the stress out of failures! Get fine-grained retry control and detailed insights with OpenCensus and OpenTelemetry support.

How the service config can be used by service owners to control client behavior.

Explains the status codes used in gRPC.

Explains how to configure RPCs to wait for the server to be ready before sending the request.

# Authentication

An overview of gRPC authentication, including built-in auth mechanisms, and how to plug in your own authentication systems.

# Authentication

### Overview

gRPC is designed to work with a variety of authentication mechanisms, making it easy to safely use gRPC to talk to other systems. You can use our supported mechanisms - SSL/TLS with or without Google token-based authentication - or you can plug in your own authentication system by extending our provided code.

gRPC also provides a simple authentication API that lets you provide all the
necessary authentication information as `Credentials` when creating a channel or
making a call.

### Supported auth mechanisms

The following authentication mechanisms are built-in to gRPC:

- **SSL/TLS**: gRPC has SSL/TLS integration and promotes the use of SSL/TLS to authenticate the server, and to encrypt all the data exchanged between the client and the server. Optional mechanisms are available for clients to provide certificates for mutual authentication.
- **ALTS**: gRPC supports ALTS as a transport security mechanism, if the application is running on Compute Engine or Google Kubernetes Engine (GKE). For details, see one of the following language-specific pages: ALTS in C++, ALTS in Go, ALTS in Java, ALTS in Python.
- **Token-based authentication with Google**: gRPC provides a generic mechanism (described below) to attach metadata based credentials to requests and responses. Additional support for acquiring access tokens (typically OAuth2 tokens) while accessing Google APIs through gRPC is provided for certain auth flows: you can see how this works in our code examples below. In general this mechanism must be used- *as well as*SSL/TLS on the channel - Google will not allow connections without SSL/TLS, and most gRPC language implementations will not let you send credentials on an unencrypted channel.

#### Warning

Google credentials should only be used to connect to Google services. Sending a Google issued OAuth2 token to a non-Google service could result in this token being stolen and used to impersonate the client to Google services.### Authentication API

gRPC provides a simple authentication API based around the unified concept of Credentials objects, which can be used when creating an entire gRPC channel or an individual call.

#### Credential types

Credentials can be of two types:

- **Channel credentials**, which are attached to a- `Channel`, such as SSL credentials.
- **Call credentials**, which are attached to a call (or- `ClientContext`in C++).

You can also combine these in a `CompositeChannelCredentials`, allowing you to
specify, for example, SSL details for the channel along with call credentials
for each call made on the channel. A `CompositeChannelCredentials` associates a
`ChannelCredentials` and a `CallCredentials` to create a new
`ChannelCredentials`. The result will send the authentication data associated
with the composed `CallCredentials` with every call made on the channel.

For example, you could create a `ChannelCredentials` from an `SslCredentials`
and an `AccessTokenCredentials`. The result when applied to a `Channel` would
send the appropriate access token for each call on this channel.

Individual `CallCredentials` can also be composed using
`CompositeCallCredentials`. The resulting `CallCredentials` when used in a call
will trigger the sending of the authentication data associated with the two
`CallCredentials`.

#### Using client-side SSL/TLS

Now let’s look at how `Credentials` work with one of our supported auth
mechanisms. This is the simplest authentication scenario, where a client just
wants to authenticate the server and encrypt all data. The example is in C++,
but the API is similar for all languages: you can see how to enable SSL/TLS in
more languages in our Examples section below.

```
// Create a default SSL ChannelCredentials object.
auto channel_creds = grpc::SslCredentials(grpc::SslCredentialsOptions());
// Create a channel using the credentials created in the previous step.
auto channel = grpc::CreateChannel(server_name, channel_creds);
// Create a stub on the channel.
std::unique_ptr<Greeter::Stub> stub(Greeter::NewStub(channel));
// Make actual RPC calls on the stub.
grpc::Status s = stub->sayHello(&context, *request, response);
```
For advanced use cases such as modifying the root CA or using client certs,
the corresponding options can be set in the `SslCredentialsOptions` parameter
passed to the factory method.

#### Note

Non-POSIX-compliant systems (such as Windows) need to specify the root certificates in`SslCredentialsOptions`, since the defaults are only
configured for POSIX filesystems.#### Using OAuth token-based authentication

OAuth 2.0 Protocol is the industry-standard protocol for authorization. It enables websites or applications to obtain limited access to user accounts using OAuth tokens.

gRPC offers a set of simple APIs to integrate OAuth 2.0 into applications, streamlining authentication.

At a high level, using OAuth token-based authentication includes 3 steps:

- Get or generate an OAuth token on client side.- You can generate Google-specific tokens following instructions below.

- Create credentials with the OAuth token.- OAuth token is always part of per-call credentials, you can also attach the per-call credentials to some channel credentials.
- The token will be sent to server, normally as part of HTTP Authorization header.

- Server side verifies the token.- In most implementations, the validation is done using a server side interceptor.

For details of how to use OAuth token in different languages, please refer to our examples below.

#### Using Google token-based authentication

gRPC applications can use a simple API to create a credential that works for authentication with Google in various deployment scenarios. Again, our example is in C++ but you can find examples in other languages in our Examples section.

```
auto creds = grpc::GoogleDefaultCredentials();
// Create a channel, stub and make RPC calls (same as in the previous example)
auto channel = grpc::CreateChannel(server_name, creds);
std::unique_ptr<Greeter::Stub> stub(Greeter::NewStub(channel));
grpc::Status s = stub->sayHello(&context, *request, response);
```
This channel credentials object works for applications using Service Accounts as
well as for applications running in Google Compute Engine
(GCE). In the former case, the service
account’s private keys are loaded from the file named in the environment
variable `GOOGLE_APPLICATION_CREDENTIALS`. The keys are used to generate bearer
tokens that are attached to each outgoing RPC on the corresponding channel.

For applications running in GCE, a default service account and corresponding OAuth2 scopes can be configured during VM setup. At run-time, this credential handles communication with the authentication systems to obtain OAuth2 access tokens and attaches them to each outgoing RPC on the corresponding channel.

#### Extending gRPC to support other authentication mechanisms

The Credentials plugin API allows developers to plug in their own type of credentials. This consists of:

- The `MetadataCredentialsPlugin`abstract class, which contains the pure virtual`GetMetadata`method that needs to be implemented by a sub-class created by the developer.
- The `MetadataCredentialsFromPlugin`function, which creates a`CallCredentials`from the`MetadataCredentialsPlugin`.

Here is example of a simple credentials plugin which sets an authentication ticket in a custom header.

```
class MyCustomAuthenticator : public grpc::MetadataCredentialsPlugin {
 public:
 MyCustomAuthenticator(const grpc::string& ticket) : ticket_(ticket) {}
 grpc::Status GetMetadata(
 grpc::string_ref service_url, grpc::string_ref method_name,
 const grpc::AuthContext& channel_auth_context,
 std::multimap<grpc::string, grpc::string>* metadata) override {
 metadata->insert(std::make_pair("x-custom-auth-ticket", ticket_));
 return grpc::Status::OK;
 }
 private:
 grpc::string ticket_;
};
auto call_creds = grpc::MetadataCredentialsFromPlugin(
 std::unique_ptr<grpc::MetadataCredentialsPlugin>(
 new MyCustomAuthenticator("super-secret-ticket")));
```
A deeper integration can be achieved by plugging in a gRPC credentials implementation at the core level. gRPC internals also allow switching out SSL/TLS with other encryption mechanisms.

### Language guides and examples

These authentication mechanisms will be available in all gRPC’s supported languages. The following table links to examples demonstrating authentication and authorization in various languages.

| Language | Example | Documentation |
|---|---|---|
| C++ | N/A | N/A |
| Go | Go Example | Go Documentation |
| Java | Java Example TLS (Java Example ATLS) | Java Documentation |
| Python | Python Example | Python Documentation |

### Language guides and examples for OAuth token-based authentication

The following table links to examples demonstrating OAuth token-based authentication and authorization in various languages.

| Language | Example | Documentation |
|---|---|---|
| C++ | N/A | N/A |
| Go | Go OAuth Example | Go OAuth Documentation |
| Java | Java OAuth Example | Java OAuth Documentation |
| Python | Python OAuth Example | Python OAuth Documentation |

### Additional Examples

The following sections demonstrate how authentication and authorization features described above appear in other languages not listed above.

#### Ruby

##### Base case - no encryption or authentication

```
stub = Helloworld::Greeter::Stub.new('localhost:50051', :this_channel_is_insecure)
...
```
##### With server authentication SSL/TLS

```
creds = GRPC::Core::ChannelCredentials.new(load_certs) # load_certs typically loads a CA roots file
stub = Helloworld::Greeter::Stub.new('myservice.example.com', creds)
```
##### Authenticate with Google

```
require 'googleauth' # from http://www.rubydoc.info/gems/googleauth/0.1.0
...
ssl_creds = GRPC::Core::ChannelCredentials.new(load_certs) # load_certs typically loads a CA roots file
authentication = Google::Auth.get_application_default()
call_creds = GRPC::Core::CallCredentials.new(authentication.updater_proc)
combined_creds = ssl_creds.compose(call_creds)
stub = Helloworld::Greeter::Stub.new('greeter.googleapis.com', combined_creds)
```
#### Node.js

##### Base case - No encryption/authentication

```
var stub = new helloworld.Greeter('localhost:50051', grpc.credentials.createInsecure());
```
##### With server authentication SSL/TLS

```
const root_cert = fs.readFileSync('path/to/root-cert');
const ssl_creds = grpc.credentials.createSsl(root_cert);
const stub = new helloworld.Greeter('myservice.example.com', ssl_creds);
```
##### Authenticate with Google

```
// Authenticating with Google
var GoogleAuth = require('google-auth-library'); // from https://www.npmjs.com/package/google-auth-library
...
var ssl_creds = grpc.credentials.createSsl(root_certs);
(new GoogleAuth()).getApplicationDefault(function(err, auth) {
 var call_creds = grpc.credentials.createFromGoogleCredential(auth);
 var combined_creds = grpc.credentials.combineChannelCredentials(ssl_creds, call_creds);
 var stub = new helloworld.Greeter('greeter.googleapis.com', combined_credentials);
});
```
##### Authenticate with Google using OAuth2 token (legacy approach)

```
var GoogleAuth = require('google-auth-library'); // from https://www.npmjs.com/package/google-auth-library
...
var ssl_creds = grpc.Credentials.createSsl(root_certs); // load_certs typically loads a CA roots file
var scope = 'https://www.googleapis.com/auth/grpc-testing';
(new GoogleAuth()).getApplicationDefault(function(err, auth) {
 if (auth.createScopeRequired()) {
 auth = auth.createScoped(scope);
 }
 var call_creds = grpc.credentials.createFromGoogleCredential(auth);
 var combined_creds = grpc.credentials.combineChannelCredentials(ssl_creds, call_creds);
 var stub = new helloworld.Greeter('greeter.googleapis.com', combined_credentials);
});
```
##### With server authentication SSL/TLS and a custom header with token

```
const rootCert = fs.readFileSync('path/to/root-cert');
const channelCreds = grpc.credentials.createSsl(rootCert);
const metaCallback = (_params, callback) => {
 const meta = new grpc.Metadata();
 meta.add('custom-auth-header', 'token');
 callback(null, meta);
}
const callCreds = grpc.credentials.createFromMetadataGenerator(metaCallback);
const combCreds = grpc.credentials.combineChannelCredentials(channelCreds, callCreds);
const stub = new helloworld.Greeter('myservice.example.com', combCreds);
```
#### PHP

##### Base case - No encryption/authorization

```
$client = new helloworld\GreeterClient('localhost:50051', [
 'credentials' => Grpc\ChannelCredentials::createInsecure(),
]);
```
##### With server authentication SSL/TLS

```
$client = new helloworld\GreeterClient('myservice.example.com', [
 'credentials' => Grpc\ChannelCredentials::createSsl(file_get_contents('roots.pem')),
]);
```
##### Authenticate with Google

```
function updateAuthMetadataCallback($context)
{
 $auth_credentials = ApplicationDefaultCredentials::getCredentials();
 return $auth_credentials->updateMetadata($metadata = [], $context->service_url);
}
$channel_credentials = Grpc\ChannelCredentials::createComposite(
 Grpc\ChannelCredentials::createSsl(file_get_contents('roots.pem')),
 Grpc\CallCredentials::createFromPlugin('updateAuthMetadataCallback')
);
$opts = [
 'credentials' => $channel_credentials
];
$client = new helloworld\GreeterClient('greeter.googleapis.com', $opts);
```
##### Authenticate with Google using OAuth2 token (legacy approach)

```
// the environment variable "GOOGLE_APPLICATION_CREDENTIALS" needs to be set
$scope = "https://www.googleapis.com/auth/grpc-testing";
$auth = Google\Auth\ApplicationDefaultCredentials::getCredentials($scope);
$opts = [
 'credentials' => Grpc\Credentials::createSsl(file_get_contents('roots.pem'));
 'update_metadata' => $auth->getUpdateMetadataFunc(),
];
$client = new helloworld\GreeterClient('greeter.googleapis.com', $opts);
```
#### Dart

##### Base case - no encryption or authentication

```
final channel = new ClientChannel('localhost',
 port: 50051,
 options: const ChannelOptions(
 credentials: const ChannelCredentials.insecure()));
final stub = new GreeterClient(channel);
```
##### With server authentication SSL/TLS

```
// Load a custom roots file.
final trustedRoot = new File('roots.pem').readAsBytesSync();
final channelCredentials =
 new ChannelCredentials.secure(certificates: trustedRoot);
final channelOptions = new ChannelOptions(credentials: channelCredentials);
final channel = new ClientChannel('myservice.example.com',
 options: channelOptions);
final client = new GreeterClient(channel);
```
##### Authenticate with Google

```
// Uses publicly trusted roots by default.
final channel = new ClientChannel('greeter.googleapis.com');
final serviceAccountJson =
 new File('service-account.json').readAsStringSync();
final credentials = new JwtServiceAccountAuthenticator(serviceAccountJson);
final client =
 new GreeterClient(channel, options: credentials.toCallOptions);
```
##### Authenticate a single RPC call

```
// Uses publicly trusted roots by default.
final channel = new ClientChannel('greeter.googleapis.com');
final client = new GreeterClient(channel);
...
final serviceAccountJson =
 new File('service-account.json').readAsStringSync();
final credentials = new JwtServiceAccountAuthenticator(serviceAccountJson);
final response =
 await client.sayHello(request, options: credentials.toCallOptions);
```

gRPC is designed to support high-performance open-source RPCs in many languages. This page describes performance benchmarking tools, scenarios considered by tests, and the testing infrastructure.

gRPC is designed to support high-performance open-source RPCs in many languages. This page describes performance benchmarking tools, scenarios considered by tests, and the testing infrastructure.

Overview

gRPC is designed for both high-performance and high-productivity design of
distributed applications. Continuous performance benchmarking is a critical part
of the gRPC development workflow. Multi-language performance tests run every few
hours against the master branch, and these numbers are reported to a dashboard
for visualization.

Each language implements a performance testing worker that implements a gRPC
WorkerService.
This service directs the worker to act as either a client or a server for the
actual benchmark test, represented as
BenchmarkService.
That service has two methods:

UnaryCall – a unary RPC of a simple request that specifies the number of bytes
to return in the response.

StreamingCall – a streaming RPC that allows repeated ping-pongs of request and
response messages akin to the UnaryCall.

# Cancellation

Explains how and when to cancel RPCs.

# Cancellation

### Overview

When a gRPC client is no longer interested in the result of an RPC call, it may
*cancel* to signal this discontinuation of interest to the server.
Deadline expiration and I/O errors
also trigger cancellation. When an RPC is cancelled, the server should stop
any ongoing computation and end its side of the stream. Often, servers are also
clients to upstream servers, so that cancellation operation should ideally
propagate to all ongoing computation in the system that was initiated due to
the original client RPC call.

A client may cancel an RPC for several reasons. The data it requested may have been made irrelevant or the author of the client may want to be a good citizen of the server and conserve compute resources.

```
sequenceDiagram
 Client ->> Server 1: Cancel
 Server 1 ->> Server 2: Cancel
```
### Cancelling an RPC Call on the Client Side

A client cancels an RPC call by calling a method on the call object or, in some languages, on the accompanying context object. While gRPC clients do not provide additional details to the server about the reason for the cancellation, the cancel API call takes a string describing the reason, which will result in a client-side exception and/or log containing the provided reason. When a server is notified of the cancellation of an RPC, the application-provided server handler may be busy processing the request. The gRPC library in general does not have a mechanism to interrupt the application-provided server handler, so the server handler must coordinate with the gRPC library to ensure that local processing of the request ceases. Therefore, if an RPC is long-lived, its server handler must periodically check if the RPC it is servicing has been cancelled and if it has, cease processing. Some languages will also support automatic cancellation of anyoutgoing RPCs, while in others, the author of the server handler is responsible for this.

```
flowchart LR
 subgraph Client
 end
 subgraph Server1
 direction TB
 cancelled{cancelled?} -->|false| perform("perform some work")
 perform --> cancelled
 cancelled -->|true| cleanup("cancel upstream RPCs")
 cleanup --> exit("exit RPC handler")
 end
 subgraph Server2
 end
 Client -->|CANCEL| Server1
 Server1 -->|CANCEL| Server2
```
### Language Support

| Language | Example | Notes |
|---|---|---|
| Java | Example | Automatically cancels outgoing RPCs |
| Go | Example | Automatically cancels outgoing RPCs |
| C++ | Example | Automatically cancels outgoing RPCs |
| Python | Example |

# Compression

How to compress the data sent over the wire while using gRPC.

# Compression

### Overview

Compression is used to reduce the amount of bandwidth used when communicating between peers and can be enabled or disabled based on call or message level for all languages. For some languages, it is also possible to control compression settings at the channel level. Different languages also support different compression algorithms, including a customized compressor.

### Compression Method Asymmetry Between Peers

gRPC allows asymmetrically compressed communication, whereby a response may be compressed differently with the request, or not compressed at all. A gRPC peer may choose to respond using a different compression method to that of the request, including not performing any compression, regardless of channel and RPC settings (for example, if compression would result in small or negative gains).

If a client message is compressed by an algorithm that is not supported by a server, the message will result in an `UNIMPLEMENTED` error status on the server. The server will include a `grpc-accept-encoding` header to the response which specifies the algorithms that the server accepts.

If the client message is compressed using one of the algorithms from the `grpc-accept-encoding` header and an `UNIMPLEMENTED` error status is returned from the server, the cause of the error won’t be related to compression.

Note that a peer may choose to not disclose all the encodings it supports. However, if it receives a message compressed in an undisclosed but supported encoding, it will include said encoding in the response’s `grpc-accept-encoding` header.

For every message a server is requested to compress using an algorithm it knows the client doesn’t support (as indicated by the last `grpc-accept-encoding` header received from the client), it will send the message uncompressed.

### Specific Disabling of Compression

If the user requests to disable compression, the next message will be sent uncompressed. This is instrumental in preventing BEAST and CRIME attacks. This applies to both the unary and streaming cases.

### Language guides and examples

| Language | Example | Documentation |
|---|---|---|
| C++ | C++ Example | C++ Documentation |
| Go | Go Example | Go Documentation |
| Java | Java Example | Java Documentation |
| Python | Python Example | Python Documentation |

# Custom Backend Metrics

A mechanism in the gRPC library that allows users to inject custom metrics at the gRPC server and consume at gRPC clients to make your custom load balancing algorithms.

# Custom Backend Metrics

### Overview

Simple load balancing decisions can be made by taking into account local or global knowledge of a backend’s load, for example CPU. More sophisticated load balancing decisions are possible with application specific knowledge, e.g. queue depth, or by combining multiple metrics.

The custom backend metrics feature exposes APIs to allow users to implement the metrics feedback in their LB policies.

### Use Cases

The feature is mainly for advanced use cases where a custom LB policy is used to route traffic more intelligently to a list of backend servers to improve the routing performance, e.g. a weighted round robin LB policy.

gRPC traditionally allows users to plug in their own load balancing policies, see guide. For xDS users, custom load balancer can be configured to select the custom LB policy.

### Metrics Reporting

Open Request Cost Aggregation (ORCA) is an open standard for conveying backend metrics information. gRPC uses ORCA service and metrics standards and supports two metrics reporting mechanisms:

- Per-query metrics reporting: the backend server attaches the injected custom metrics in the trailing metadata when the corresponding RPC finishes. This is typically useful for short RPCs like unary calls.
- Out-of-band metrics reporting: the backend server periodically pushes metrics data, e.g. cpu and memory utilization, to the client. This is useful for all situations: unary calls, long RPCs in streaming calls, or no RPCs. However, out-of-band metrics reporting does not send query cost metrics. The metrics emission frequency is user-configurable, and this configuration resides in the custom load balancing policy.

The diagram shows the architecture where a user creates their own LB policy that implements backend metrics feedback.

For more details, please see gRPC proposal A51.## Implementation

## Language Support

Language Example Java Java example Go Go example C++ Example upcoming

# Custom Load Balancing Policies

Explains how custom load balancing policies can help optimize load balancing under unique circumstances.

# Custom Load Balancing Policies

### Overview

One of the key features of gRPC is load balancing, which allows requests from clients to be distributed across multiple servers. This helps prevent any one server from becoming overloaded and allows the system to scale up by adding more servers.

A gRPC load balancing policy is given a list of server IP addresses by the name resolver. The policy is responsible for maintaining connections (subchannels) to the servers and picking a connection to use when an RPC is sent.

### Implementing Your Own Policy

By default the `pick_first` policy will be used. This policy actually does no
load balancing but just tries each address it gets from the name resolver and
uses the first one it can connect to. By updating the gRPC service config you
can also switch to using `round_robin` that connects to every address it gets
and rotates through the connected backends for each RPC. There are also some
other load balancing policies available, but the exact set varies by language.
If the built-in policies do not meet your needs you can also implement your own
custom policy.

This involves implementing a load balancer interface in the language you are using. At a high level, you will have to:

- Register your implementation in the load balancer registry so that it can be referred to from the service config
- Parse the JSON configuration object of your implementation. This allows your load balancer to be configured in the service config with any arbitrary JSON you choose to support
- Manage what backends to maintain a connection with
- Implement a `picker`that will choose which backend to connect to when an RPC is made. Note that this needs to be a fast operation as it is on the RPC call path
- To enable your load balancer, configure it in your service config

The exact steps vary by language, see the language support section for some concrete examples in your language.

```
flowchart TD
NR(Name Resolver) -->|Provides addresses &\nLB config| LB(Load Balancer)
LB --> |Provides a picker| C(Channel)
C -->|Requests\na subchannel| P(Picker)
LB --> |Manages subchannels\nto backends| SC(Subchannel 1..n)
LB -. Creates .-> P
P --> |Picks one| SC
```
### Backend Metrics

What if your load balancing policy needs real-time information about the backend servers? For this you can rely on backend metrics. You can have metrics provided to you either in-band, in the backend RPC responses, or out-of-band as separate RPCs from the backends. Standard metrics like CPU and memory utilization are provided, but you can also implement your own custom metrics.

For more information on this, please see the custom backend metrics guide

### Service Mesh

If you have a service mesh setup where a central control plane is coordinating the configuration of your microservices, you cannot configure your custom load balancer directly via the service config. But support is provided to do this with the xDS protocol that your control plane uses to communicate with your gRPC clients. Please refer to your control plane documentation to determine how custom load balancing configuration is supported.

For more details, please see gRPC proposal A52.

### Language Support

| Language | Example | Notes |
|---|---|---|
| Java | Java example | |
| Go | Go example | |
| C++ | Not yet supported |

# Custom Name Resolution

Explains standard name resolution, the custom name resolver interface, and how to write an implementation.

# Custom Name Resolution

### Overview

Name resolution is fundamentally about service discovery. When sending a gRPC request, the client must determine the IP address of the service name. Name resolution is often thought to be the same as DNS. In practice however, DNS is usually augmented with extensions or completely replaced to enable name resolution.

When making a request with a gRPC client, by default, DNS name resolution is used. However, various other name resolution mechanisms may be used:

| Resolver | Example | Notes |
|---|---|---|
| DNS | `grpc.io:50051` | By default, DNS is assumed. |
| DNS | `dns:///grpc.io:50051` | The extra slash is used to provide an authority |
| Unix Domain Socket | `unix:///run/containerd/containerd.sock` | |
| xDS | `xds:///wallet.grpcwallet.io` | |
| IPv4 | `ipv4:198.51.100.123:50051` | Only supported in some languages |

#### Note

The triple slashes above (`///`) may look unfamiliar if you are used to the
double slashes of HTTP, such as `https://grpc.io`. These *target strings*follow the format for RFC-3986 URIs. The string following the first two slashes and preceding the third (if there is a third at all) is the

*authority*. The authority string identifies a server which contains the URIs of all resources. In the case of a conventional HTTP request, the authority over the URI is the server to which the request will be sent. In other cases, the authority will be the identity of the name resolution server, while the resource itself lives on some other server. Some name resolvers have no need for an authority. In this case, the authority string is left empty, resulting in three slashes in a row.

Several languages support an interface to allow the user to define their own
name resolvers, so that you may define how to resolve any given name. Once
registered, a name resolver with the *scheme* `my-resolver` will be picked up
when a target string begins with `my-resolver:`. For example, requests to
`my-resolver:///my-service` would now use the `my-resolver` name resolver
implementation.

### Custom Name Resolvers

You might consider using a custom name resolver whenever you would like to augment or replace DNS for service discovery. For example, this interface has been used in the past to use Apache Zookeeper to look up service names. It has also been used to directly interface with the Kubernetes API server for service lookup based on headless Service resources.

One reason why it might be particularly useful to use a custom name resolver
rather than standard DNS is that this interface is *reactive*. Within standard
DNS, a client looks up the address for a particular service at the beginning of
the connection and maintains its connection to that address for the lifetime of
the connection. However, custom name resolvers may be watch-based. That is, they
can receive updates from the name server over time and therefore respond
intelligently to backend failure as well as backend scale-ups and backend
scale-downs.

In addition, a custom name resolver may provide the client connection with a
*service config*. A service config is a JSON object that defines
arbitrary configuration specifying how traffic should be routed to and load
balanced across a particular service. At its most basic, this can be used to
specify things like that a particular service should use the round robin load
balancing policy vs. pick first. However, when a custom name resolver is used in
conjunction with arbitrary service config and a *custom load balancing
policy*, very complex
traffic management systems such as xDS may be constructed.

#### Life of a Target String

While the exact interface for custom name resolvers differs from language to
language, the general structure is the same. The client registers an
implementation of a *name resolver provider* to a process-global registry close
to the start of the process. The name resolver provider will be called by the
gRPC library with a target strings intended for the custom name resolver. Given
that target string, the name resolver provider will return an instance of a name
resolver, which will interact with the client connection to direct the request
according to the target string.

```
sequenceDiagram
 Client ->> gRPC: Request to my-resolver:///my-service
 gRPC ->> NameResolverProvider: requests NameResolver
 NameResolverProvider -->> gRPC: returns NameResolver
 gRPC ->> NameResolver: delegates resolution
 NameResolver -->> gRPC: addresses
```
### Language Support

| Language | Example |
|---|---|
| Java | Example |
| Go | Example |
| C++ | Not supported |
| Python | Not supported |

# Deadlines

Explains how deadlines can be used to effectively deal with unreliable backends.

# Deadlines

### Overview

A deadline is used to specify a point in time past which a client is unwilling to wait for a response from a server. This simple idea is very important in building robust distributed systems. Clients that do not wait around unnecessarily and servers that know when to give up processing requests will improve the resource utilization and latency of your system.

Note that while some language APIs have the concept of a **deadline**, others
use the idea of a **timeout**. When an API asks for a deadline, you provide a
point in time which the call should not go past. A timeout is the max duration
of time that the call can take. A timeout can be converted to a deadline by
adding the timeout to the current time when the application starts a call. For
simplicity, we will only refer to deadline in this document.

### Deadlines on the Client

By default, gRPC does not set a deadline which means it is possible for a client to end up waiting for a response effectively forever. To avoid this you should always explicitly set a realistic deadline in your clients. To determine the appropriate deadline you would ideally start with an educated guess based on what you know about your system (network latency, server processing time, etc.), validated by some load testing.

If a server has gone past the deadline when processing a request, the client
will give up and fail the RPC with the `DEADLINE_EXCEEDED` status.

### Deadlines on the Server

A server might receive RPCs from a client with an unrealistically short
deadline that would not give the server enough time to ever respond in time.
This would result in the server just wasting valuable resources and in the worst
case scenario, crash the server. A gRPC server deals with this situation by
automatically cancelling a call (`CANCELLED` status) once a deadline set by the
client has passed.

Please note that the server application is responsible for stopping any activity it has spawned to service the RPC. If your application is running a long-running process you should periodically check if the RPC that initiated it has been cancelled and if so, stop the processing.

#### Deadline Propagation

Your server might need to call another server to produce a response. In these cases where your server also acts as a client you would want to honor the deadline set by the original client. Automatically propagating the deadline from an incoming RPC to an outgoing one is supported by some gRPC implementations. In some languages this behavior needs to be explicitly enabled (e.g. C++) and in others it is enabled by default (e.g. Java and Go). Using this capability lets you avoid the error-prone approach of manually including the deadline for each outgoing RPC.

Since a deadline is set point in time, propagating it as-is to a server can be problematic as the clocks on the two servers might not be synchronized. To address this gRPC converts the deadline to a timeout from which the already elapsed time is already deducted. This shields your system from any clock skew issues.

```
%%{init: { "sequence": { "mirrorActors": false }}}%%
sequenceDiagram
 participant c as Client
 participant us as User Server
 participant bs as Billing Server
 note right of c: Request at 13:00:00<br>Should complete in 2s
 activate c
 c ->> us: GetUserProfile<br>(deadline: 13:00:02)
 activate us
 note right of us: 0.5s spent before<br>calling billing server
 us ->> bs: GetTransactionHistory<br>(timeout: 1.5s)
 activate bs
 bs ->> bs: Retrieve transactions
 note left of bs: It's 13:00:02<br>Time's up!
 note right of c: Stop waiting for server
 c ->> c: Stop waiting for server<br>DEADLINE_EXCEEDED
 deactivate c
 us ->> us: Stop waiting for server
 us -->> c: Cancel
 deactivate us
 bs -->> us: Cancel
 bs ->> bs: Clean up resources<br>(after noticing that the<br>call was cancelled)
 deactivate bs

```
### Language Support

| Language | Example |
|---|---|
| Java | Java example |
| Go | Go example |
| C++ | C++ example |
| Python | Python example |

# Debugging

Explains the debugging process of gRPC applications using grpcdebug

# Debugging

### Overview

grpcdebug is a command line tool within the gRPC ecosystem designed to assist developers in debugging and troubleshooting gRPC services. grpcdebug fetches the internal states of the gRPC library from the application via gRPC protocol and provides a human-friendly UX to browse them. Currently, it supports Channelz/Health Checking/CSDS (aka. admin services). In other words, it can fetch statistics about how many RPCs have being sent or failed on a given gRPC channel, it can inspect address resolution results, it can dump the active xDS configuration that directs the routing of RPCs.

### Language examples

| Language | Example | Notes |
|---|---|---|
| C++ | C++ Example | |
| Go | Go Example | Go test server implementing admin services from grpcdebug docs |
| Java | Java Example |

# Error handling

How gRPC deals with errors, and gRPC error codes.

# Error handling

### Standard error model

As you’ll have seen in our concepts document and examples, when a gRPC call
completes successfully the server returns an `OK` status to the client
(depending on the language the `OK` status may or may not be directly used in
your code). But what happens if the call isn’t successful?

If an error occurs, gRPC returns one of its error status codes instead, with an optional string error message that provides further details about what happened. Error information is available to gRPC clients in all supported languages.

### Richer error model

The error model described above is the official gRPC error model, is supported by all gRPC client/server libraries, and is independent of the gRPC data format (whether protocol buffers or something else). You may have noticed that it’s quite limited and doesn’t include the ability to communicate error details.

If you’re using protocol buffers as your data format, however, you may wish to consider using the richer error model developed and used by Google as described here. This model enables servers to return and clients to consume additional error details expressed as one or more protobuf messages. It further specifies a standard set of error message types to cover the most common needs (such as invalid parameters, quota violations, and stack traces). The protobuf binary encoding of this extra error information is provided as trailing metadata in the response.

This richer error model is already supported in the C++, Go, Java, Python, and Ruby libraries, and at least the grpc-web and Node.js libraries have open issues requesting it. Other language libraries may add support in the future if there’s demand, so check their github repos if interested. Note however that the grpc-core library written in C will not likely ever support it since it is purposely data format agnostic.

You could use a similar approach (put error details in trailing response metadata) if you’re not using protocol buffers, but you’d likely need to find or develop library support for accessing this data in order to make practical use of it in your APIs.

There are important considerations to be aware of when deciding whether to use such an extended error model, however, including:

- Library implementations of the extended error model may not be consistent across languages in terms of requirements for and expectations of the error details payload
- Existing proxies, loggers, and other standard HTTP request processors don’t have visibility into the error details and thus wouldn’t be able to leverage them for monitoring or other purposes
- Additional error detail in the trailers interferes with head-of-line blocking, and will decrease HTTP/2 header compression efficiency due to more frequent cache misses
- Larger error detail payloads may run into protocol limits (like max headers size), effectively losing the original error

### Error status codes

Errors are raised by gRPC under various circumstances, from network failures to unauthenticated connections, each of which is associated with a particular status code. The following error status codes are supported in all gRPC languages.

#### General errors

| Case | Status code |
|---|---|
| Client application cancelled the request | `GRPC_STATUS_CANCELLED` |
| Deadline expired before server returned status | `GRPC_STATUS_DEADLINE_EXCEEDED` |
| Method not found on server | `GRPC_STATUS_UNIMPLEMENTED` |
| Server shutting down | `GRPC_STATUS_UNAVAILABLE` |
| Server threw an exception (or did something other than returning a status code to terminate the RPC) | `GRPC_STATUS_UNKNOWN` |

#### Network failures

| Case | Status code |
|---|---|
| No data transmitted before deadline expires. Also applies to cases where some data is transmitted and no other failures are detected before the deadline expires | `GRPC_STATUS_DEADLINE_EXCEEDED` |
| Some data transmitted (for example, the request metadata has been written to the TCP connection) before the connection breaks | `GRPC_STATUS_UNAVAILABLE` |

#### Protocol errors

| Case | Status code |
|---|---|
| Could not decompress but compression algorithm supported | `GRPC_STATUS_INTERNAL` |
| Compression mechanism used by client not supported by the server | `GRPC_STATUS_UNIMPLEMENTED` |
| Flow-control resource limits reached | `GRPC_STATUS_RESOURCE_EXHAUSTED` |
| Flow-control protocol violation | `GRPC_STATUS_INTERNAL` |
| Error parsing returned status | `GRPC_STATUS_UNKNOWN` |
| Unauthenticated: credentials failed to get metadata | `GRPC_STATUS_UNAUTHENTICATED` |
| Invalid host set in authority metadata | `GRPC_STATUS_UNAUTHENTICATED` |
| Error parsing response protocol buffer | `GRPC_STATUS_INTERNAL` |
| Error parsing request protocol buffer | `GRPC_STATUS_INTERNAL` |

### Language Support

Examples code is available for multiple languages on how to deal with standard errors as well as with the richer error details.

The grpc-errors repo also contains additional error handling examples.

# Flow Control

Explains what flow control is and how you can manually control it.

# Flow Control

### Overview

Flow control is a mechanism to ensure that a receiver of messages does not get overwhelmed by a fast sender. Flow control prevents data loss, improves performance and increases reliability. It applies to streaming RPCs and is not relevant for unary RPCs. By default, gRPC handles the interactions with flow control for you, though some languages allow you to override the default behavior and take explicit control.

gRPC utilizes the underlying transport to detect when it is safe to send more data. As data is read on the receiving side, an acknowledgement is returned to the sender letting it know that the receiver has more capacity.

As needed, the gRPC framework will wait before returning from a write call. In gRPC, when a value is written to a stream, that does not mean that it has gone out over the network. Rather, that it has been passed to the framework which will now take care of the nitty gritty details of buffering it and sending it to the OS on its way over the network.

#### Note

The flow is the same for writing from a Server to a Client as when a Client writes to a Server```
sequenceDiagram
 participant SA as Sender Application
 participant SG as Sender gRPC Framework
 participant RG as Receiver gRPC Framework
 participant RA as Receiver Application

 SA-)+SG: Stream Write
 alt sending too fast
 SG--)SG: Wait
 end
 alt allowed to send
 SG--)-SA: Write call returns
 SG->>RG:Send Msg
 end
 RA->>RG: Request message
 Note right of RA: Request can be done either<br>after or before message arrives
 RG->>RA: Provide message
 RG->>SG: Send Ack w/ msg size
 opt waiting messages
 SG->>RG: Send Next Msg
 end
```
#### Warning

There is the potential for a deadlock if both the client and server are doing synchronous reads or using manual flow control and both try to do a lot of writing without doing any reads.### Language Support

| Language | Example |
|---|---|
| Java | Java Example |

# Graceful Shutdown

Explains how to gracefully shut down a gRPC server to avoid causing RPC failures for connected clients.

# Graceful Shutdown

### Overview

gRPC servers often need to shut down gracefully, ensuring that in-flight RPCs are completed within a reasonable time frame and new RPCs are no longer accepted. The “Graceful shutdown function” facilitates this process, allowing the server to transition smoothly without abruptly terminating active connections.

When the “Graceful shutdown function” is called, the server immediately notifies all clients to stop sending new RPCs. Then after the clients have received that notification, the server will stop accepting new RPCs. In-flight RPCs are allowed to continue until they complete or a specified deadline is reached. Once all active RPCs finish or the deadline expires, the server shuts down completely.

Because graceful shutdown helps prevent clients from encountering RPC failures, it should be used if possible. However, gRPC also provides a forceful shutdown mechanism which will immediately cause the server to stop serving and close all connections, which results in the failure of any in-flight RPCs.

### How to do Graceful Server Shutdown

The exact implementation of the “Graceful shutdown function” varies depending on the programming language you are using. However, the general pattern involves:

- Initiating the graceful shutdown process by calling the “Graceful shutdown function” on your gRPC server object. This function blocks until all currently running RPCs complete. This ensures that in-flight requests are allowed to finish processing.
- Specify a timeout period to limit the time allowed for in-progress RPCs to finish. It’s crucial to separately call the “Forceful shutdown function” on the server object using a timer mechanism (depending on your language) to trigger a forceful shutdown after a predefined duration. This acts as a safety net, ensuring that the server eventually shuts down even if some in-flight RPCs don’t complete within a reasonable time frame. This prevents indefinite blocking.

The following shows the sequence of events that occur during the graceful shutdown process. When a server’s graceful shutdown is invoked, in-flight RPCs continue to process, but new RPCs are rejected. If some in-flight RPCs are not finished in time, the server is forcefully shut down.

```
sequenceDiagram
Client->>Server: New RPC Request 1
Client->>Server: New RPC Request 2
Server-->>Server: Graceful Shutdown Invoked
Server->>Client: Continues Processing In-Flight RPCs
Client->>Client: Detects server shutdown and finds other servers if available
alt RPCs complete within timeout
 Server->>Client: Completes RPC 1
 Server->>Client: Completes RPC 2
 Server-->>Server: Graceful Shutdown Complete
else Timeout reached
 Server->>Client: Forceful Shutdown Invoked, terminating pending RPCs
 Server-->>Server: Forceful Shutdown Complete
end
```
The following is a state based view

```
stateDiagram-v2
 [*] --> SERVING : Server Started
 SERVING --> GRACEFUL_SHUTDOWN : Graceful Shutdown Called (with Timeout)
 GRACEFUL_SHUTDOWN --> TERMINATED : In-Flight RPCs Completed (Before Timeout)
 GRACEFUL_SHUTDOWN --> TIMER_EXPIRED : Timeout Reached
 TIMER_EXPIRED --> TERMINATED : Forceful Shutdown Called
```
### Language Support

| Language | Example |
|---|---|
| C++ | |
| Go | Go Example |
| Java | Java Example |
| Python |

# Health Checking

Explains how gRPC servers expose a health checking service and how client can be configured to automatically check the health of the server it is connecting to.

# Health Checking

### Overview

gRPC specifies a standard service API (health/v1) for performing health check calls against gRPC servers. An implementation of this service is provided, but you are responsible for updating the health status of your services.

On the client side you can have the client automatically communicate with the health services of your backends. This allows the client to avoid services that are considered unhealthy.

### The Server Side Health Service

The health check service on a gRPC server supports two modes of operation:

- Unary calls to the `Check`rpc endpoint- Useful for centralized monitoring or load balancing solutions, but does not scale to support a fleet of gRPC client constantly making health checks

- Streaming health updates by using the `Watch`rpc endpoint- Used by the client side health check feature in gRPC clients

Enabling the health check service on your server involves the following steps:

- Use the provided health check library to create a health check service
- Add the health check service to your server.
- Notify the health check library when the health of one of your services
changes.- `NOT_SERVING`if your service cannot accept requests at the moment
- `SERVING`if your service is open for business
- If you don’t care about the health of individual services, you can use an empty string ("") to represent the health of your whole server.

- Make sure you inform the health check library about server shutdown so that it can notify all the connected clients.

The exact details vary by language, see the **Language Support** section below.

### Enabling Client Health Checking

A gRPC client can be configured to perform health checks against the servers
it connects to by modifying the service config of the channel. E.g. to monitor
the health of the `foo` service you would use (in JSON format):

```
{
 "healthCheckConfig": {
 "serviceName": "foo"
 }
}
```
Note that if your server reports health for the empty string ("") service, signifying the health of the whole server, you can also use an empty string here.

Enabling health checking changes some behavior around calling a server:

- The client will additionally call the `Watch`RPC on the health check service when a connection is established- If the call fails, retries will be made (with exponential backoff), unless the call fails with the status UNIMPLEMENTED, in which case health checking will be disabled.

- Requests won’t be sent until the health check service sends a healthy status for the service being called
- If a healthy service becomes unhealthy the client will no longer send requests for that service
- The calls will resume if the service later becomes healthy
- Some load balancing policies can choose to disable health checking if
the feature does not make sense with the policy (e.g. `pick_first`does this)

More specifically, the state of the subchannel (that represents the physical connection to the server) goes through these states based on the health of the service it is connecting to.

```
stateDiagram-v2
 [*] --> IDLE
 IDLE --> CONNECTING : Connection requested
 CONNECTING --> READY : Health check#colon;\nSERVING
 CONNECTING --> TRANSIENT_FAILURE : Health check#colon;\nNOT_SERVING\nor call fails
 READY --> TRANSIENT_FAILURE : Health check#colon;\nNOT_SERVING
 READY --> IDLE : Connection breaks\nor times out
 TRANSIENT_FAILURE --> READY : Health check#colon;\nSERVING
 note right of TRANSIENT_FAILURE : Allows the load balancer to choose\nanother, working subchannel
```
Again, the specifics on how to enable client side health checking varies by
language, see the examples in the **Language Support** section.

### Language Support

| Language | Example |
|---|---|
| Java | Java example |
| Go | Go example |
| Python | Python example |
| C++ | C++ example |

# Interceptors

Explains how interceptors can be used for implementing generic behavior that applies to many RPC methods.

# Interceptors

### Overview

The core of making gRPC services is implementing RPC methods. But some functionality is independent of the method being run and should apply to all or most RPCs. Interceptors are well suited to this task.

### When to Use Interceptors

You may already be familiar with the concept of interceptors, but may be used to calling them “filters” or “middleware.” Interceptors are very well suited to implementing logic that is not specific to a single RPC method. They are also easy to share across different clients or servers. Interceptors are an important and frequently-used way to extend gRPC. You might find some functionality you want is already available as an interceptor in the wider gRPC ecosystem.

Some example use cases for interceptors are:

- Metadata handling
- Logging
- Fault injection
- Caching
- Metrics
- Policy enforcement
- Server-side Authentication
- Server-side Authorization

#### Note

While*client-side*authentication could be done via an interceptor, gRPC provides a specialized “call credentials” API that is better suited to the task. See the Authentication Guide for details about client-side authentication.

### How to Use Interceptors

Interceptors can be added when building a gRPC channel or server. The interceptor is then called for every RPC on that channel or server. The interceptor APIs are different for client-side than server-side, so an interceptor will either be a “client interceptor” or a “server interceptor.”

Interceptors are inherently per-call; they are not useful for managing TCP connections, configuring the TCP port, or configuring TLS. While the proper tool for most customization, they can’t be used for everything.

#### Interceptor Order

When using multiple interceptors, their order is significant. You’ll want to make sure to understand the order your gRPC implementation will execute them. It is useful to think about the interceptors as being in a line between the application and the network. Some interceptors will be “closer to the network” and have more control over what is sent and others will be “closer to the application” which have a better view into the application’s behavior.

Suppose you have two client interceptors: a caching interceptor and a logging interceptor. What order should they be in? You might want the logging interceptor closer to the network to better monitor your application’s communication and ignore cached RPCs:

```
flowchart LR
APP(Application) --> INT1
INT1(Caching\nInterceptor) -->|Cache miss| INT2
INT2(Logging\nInterceptor) --> NET
NET(Network)
```
Or you might want it closer to the application to understand your app’s behavior and see what information it is loading:

```
flowchart LR
APP(Application) --> INT2
INT1(Caching\nInterceptor) -->|Cache miss| NET
INT2(Logging\nInterceptor) --> INT1
NET(Network)
```
You can choose between these options by just changing the order of the interceptors.

### Language Support

| Language | Example |
|---|---|
| C++ | C++ example |
| Go | Go example |
| Java | Java example |
| Python | Python example |

# Keepalive

How to use HTTP/2 PING-based keepalives in gRPC.

# Keepalive

### Overview

HTTP/2 PING-based keepalives are a way to keep an HTTP/2 connection alive even when there is no data being transferred. This is done by periodically sending a PING frame to the other end of the connection. HTTP/2 keepalives can improve performance and reliability of HTTP/2 connections, but it is important to configure the keepalive interval carefully.

#### Note

There is a related but separate concern called [Health Checking]. Health checking allows a server to signal whether a*service*is healthy while keepalive is only about the

*connection*.

### Background

TCP keepalive is a well-known method of maintaining connections and detecting broken connections. When TCP keepalive was enabled, either side of the connection can send redundant packets. Once ACKed by the other side, the connection will be considered as good. If no ACK is received after repeated attempts, the connection is deemed broken.

Unlike TCP keepalive, gRPC uses HTTP/2 which provides a mandatory PING frame which can be used to estimate round-trip time, bandwidth-delay product, or test the connection. The interval and retry in TCP keepalive don’t quite apply to PING because the transport is reliable, so they’re replaced with timeout (equivalent to interval * retry) in gRPC PING-based keepalive implementation.

#### Note

It’s not required for service owners to support keepalive.**Client authors must coordinate with service owners**for whether a particular client-side setting is acceptable. Service owners decide what they are willing to support, including whether they are willing to receive keepalives at all (If the service does not support keepalive, the first few keepalive pings will be ignored, and the server will eventually send a

`GOAWAY` message with debug data equal to the ASCII code for `too_many_pings`).### How configuring keepalive affects a call

Keepalive is less likely to be triggered for unary RPCs with quick replies. Keepalive is primarily triggered when there is a long-lived RPC, which will fail if the keepalive check fails and the connection is closed.

For streaming RPCs, if the connection is closed, any in-progress RPCs will fail. If a call is streaming data, the stream will also be closed and any data that has not yet been sent will be lost.

#### Warning

To avoid DDoSing, it’s important to take caution when setting the keepalive configurations. Thus, it is recommended to avoid enabling keepalive without calls and for clients to avoid configuring their keepalive much below one minute.### Common situations where keepalives can be useful

gRPC HTTP/2 keepalives can be useful in a variety of situations, including but not limited to:

- When sending data over a long-lived connection which might be considered as idle by proxy or load balancers.
- When the network is less reliable (For example, mobile applications).
- When using a connection after a long period of inactivity.

### Keepalive configuration specification

| Options | Availability | Description | Client Default | Server Default |
|---|---|---|---|---|
| `KEEPALIVE_TIME` | Client and Server | The interval in milliseconds between PING frames. | INT_MAX (Disabled) | 7200000 (2 hours) |
| `KEEPALIVE_TIMEOUT` | Client and Server | The timeout in milliseconds for a PING frame to be acknowledged. If sender does not receive an acknowledgment within this time, it will close the connection. | 20000 (20 seconds) | 20000 (20 seconds) |
| `KEEPALIVE_WITHOUT_CALLS` | Client | Is it permissible to send keepalive pings from the client without any outstanding streams. | 0 (false) | N/A |
| `PERMIT_KEEPALIVE_WITHOUT_CALLS` | Server | Is it permissible to send keepalive pings from the client without any outstanding streams. | N/A | 0 (false) |
| `PERMIT_KEEPALIVE_TIME` | Server | Minimum allowed time between a server receiving successive ping frames without sending any data/header frame. | N/A | 300000 (5 minutes) |
| `MAX_CONNECTION_IDLE` | Server | Maximum time that a channel may have no outstanding rpcs, after which the server will close the connection. | N/A | INT_MAX (Infinite) |
| `MAX_CONNECTION_AGE` | Server | Maximum time that a channel may exist. | N/A | INT_MAX (Infinite) |
| `MAX_CONNECTION_AGE_GRACE` | Server | Grace period after the channel reaches its max age. | N/A | INT_MAX (Infinite) |

#### Note

Some languages may provide additional options, please refer to language examples and additional resource for more details.### TCP User Timeout

Linux provides a TCP_USER_TIMEOUT socket option that fails a connection when
*any* sent packet fails to receive a TCP acknowledgement before the timeout.
gRPC implementations may enable TCP_USER_TIMEOUT (or equivalent on another
platform) automatically when keepalive is enabled, and use the same
KEEPALIVE_TIMEOUT for TCP_USER_TIMEOUT. This has the advantage of monitoring the
connection more often, with no additional networking cost nor configuration.

If a TCP load balancer is used, then TCP_USER_TIMEOUT will only monitor the connection between gRPC and the load balancer. The regular keepalive PINGs, however, are propagated through TCP load balancers, so they can detect broken connections “hidden” by the TCP load balancer from TCP_USER_TIMEOUT.

### Language guides and examples

| Language | Example | Documentation |
|---|---|---|
| C++ | C++ Example | C++ Documentation |
| Go | Go Example | Go Documentation |
| Java | Java Example | Java Documentation |
| Python | Python Example | Python Documentation |

### Additional Resources

- gRFC for Client-side Keepalive
- gRFC for Server-side Connection Management
- gRFC for TCP User Timeout
- Using gRPC for Long-lived and Streaming RPCs

# Metadata

Explains what metadata is, how it is transmitted, and what it is used for.

# Metadata

### Overview

Metadata is a side channel that allows clients and servers to provide information to each other that is associated with an RPC.

gRPC metadata is a key-value pair of data that is sent with initial or final gRPC requests or responses. It is used to provide additional information about the call, such as authentication credentials, tracing information, or custom headers.

gRPC metadata is implemented using HTTP/2 headers. The keys are ASCII
strings, while the values can be either ASCII strings or binary data. The keys
are case insensitive
and must not start with the prefix `grpc-`, which is reserved for gRPC itself.

gRPC metadata can be sent and received by both the client and the server. Headers are sent from the client to the server before the initial request and from the server to the client before the initial response of an RPC call. Trailers are sent by the server when it closes an RPC.

gRPC metadata is useful for a variety of purposes, such as:

- **Authentication**: gRPC metadata can be used to send authentication credentials to the server. This can be used to implement different authentication schemes, such as- `OAuth2`or- `JWT`using the standard HTTP Authorization header.
- **Tracing**: gRPC metadata can be used to send tracing information to the server. This can be used to track the progress of a request through a distributed system.
- **Custom headers**: gRPC metadata can be used to send custom headers to the server or from the server to the client. This can be used to implement application-specific features, such as load balancing, rate limiting or providing detailed error messages from the server to the client.
- **Internal usages**: gRPC uses HTTP/2 headers and trailers, which will be integrated with the metadata specified by your application.

See Core Concepts

#### Be Aware

```
WARNING: Servers may limit the size of Request-Headers, with a default of 8 KiB suggested.
```
Custom metadata must follow the “Custom-Metadata” format listed in PROTOCOL-HTTP2 , with the exception of binary headers not needing to be base64 encoded.

#### Headers

Headers are sent before the initial request data message from the client to the server and similarly before the initial response data from the server to the client. The header includes things like authentication credentials and how to handle the RPC. Some of the headers, such as authorization, are generated by gRPC for you.

Custom header handling is language dependent, generally through interceptors.

#### Trailers

Trailers are a special kind of header that is sent after the message data. They are used internally to communicate the outcome of an RPC. At the application level, custom trailers can be used to communicate things not directly part of the data, such as server utilization and query cost. Trailers are sent only by the server.

### For more details, please see the following gRFCs

- proposal: G1 true binary metadata
- proposal: L7 go metadata api
- proposal: L48 node metadata options
- proposal: L42 python metadata flags
- proposal: L11 ruby interceptors

### Language Support

| Language | Examples | Notes |
|---|---|---|
| Java | Java Header Java Error Handling | |
| Go | Go Metadata Go Metadata Interceptor | Go Documentation |
| C++ | C++ Metadata | |
| Node | Node Metadata | |
| Python | Python Metadata | |
| Ruby | Example upcoming |

# OpenTelemetry Metrics

OpenTelemetry Metrics available in gRPC

# OpenTelemetry Metrics

## Overview

gRPC provides support for an OpenTelemetry plugin that provides metrics that can help you:

- Troubleshoot your system
- Iterate on improving system performance
- Setup continuous monitoring and alerting.

## Background

OpenTelemetry is an observability framework to create and manage telemetry data. gRPC previously provided observability support through OpenCensus which has been sunsetted in the favor of OpenTelemetry.

## Instruments

The gRPC OpenTelemetry plugin accepts a MeterProvider and depends on the
OpenTelemetry API to create a Meter that identifies the gRPC library being
used, for example, `grpc-c++` at version `1.57.1`. The following listed
instruments are created using this meter. Users should employ the
OpenTelemetry SDK to customize the views exported by OpenTelemetry.

More and more gRPC components are being instrumented for observability. Currently, we have the following components instrumented:

- Per-call : Observe RPCs themselves (for example, latency.)- Client Per-Call (stable, on by default) : Observe a client call
- Client Per-Attempt (stable, on by default) : Observe attempts for a client call, since a call can have multiple attempts due to retry or hedging.
- Client Per-Call Retry (experimental) : Observe retry, transparent retry and hedging,
- Server : Observe a call received at the server.

- LB Policy : Observe various load-balancing policies- Weighted Round Robin (experimental)
- Pick-First (experimental)

- XdsClient (experimental)

Some instruments are off by default and need to be explicitly enabled from the gRPC OpenTelemetry plugin API. Experimental metrics are always off by default. (Reference C++ API)NOTE

### Per-Call Metrics

#### Client Per-Call Instruments

| Name | Type | Unit | Labels (required) | Description |
|---|---|---|---|---|
| grpc.client.call.duration | Histogram | s | grpc.method, grpc.target, grpc.status, grpc.client.call.custom (optional) | This metric aims to measure the end-to-end time the gRPC library takes to complete an RPC from the application’s perspective. |

Refer A66: OpenTelemetry Metrics for details.

#### Client Per-Attempt Instruments

| Name | Type | Unit | Labels (disposition) | Description |
|---|---|---|---|---|
| grpc.client.attempt. started | Counter | {attempt} | grpc.method (required), grpc.target (required), grpc.client.call.custom (optional) | The total number of RPC attempts started, including those that have not completed. |
| grpc.client.attempt. duration | Histogram | s | grpc.method (required), grpc.target (required), grpc.status (required), grpc.lb.locality (optional), grpc.lb.backend_service (optional), grpc.client.call.custom (optional) | End-to-end time taken to complete an RPC attempt including the time it takes to pick a subchannel. |
| grpc.client.attempt. sent_total_compressed_message_size | Histogram | By | grpc.method (required), grpc.target (required), grpc.status (required), grpc.lb.locality (optional), grpc.lb.backend_service (optional), grpc.client.call.custom (optional) | Total bytes (compressed but not encrypted) sent across all request messages (metadata excluded) per RPC attempt; does not include grpc or transport framing bytes. |
| grpc.client.attempt. rcvd_total_compressed_message_size | Histogram | By | grpc.method (required), grpc.target (required), grpc.status (required), grpc.lb.locality (optional), grpc.lb.backend_service (optional), grpc.client.call.custom (optional) | Total bytes (compressed but not encrypted) received across all response messages (metadata excluded) per RPC attempt; does not include grpc or transport framing bytes. |

Refer A66: OpenTelemetry Metrics for details.

#### Client Per-Call Retry Instruments

| Name | Type | Unit | Labels (required) | Description |
|---|---|---|---|---|
| grpc.client.call.retries | Histogram | {retry} | grpc.method, grpc.target, grpc.client.call.custom (optional) | Number of retries during the client call. If there were no retries, 0 is not reported. |
| grpc.client.call.transparent_retries | Histogram | {transparent_retry} | grpc.method, grpc.target, grpc.client.call.custom (optional) | Number of transparent retries during the client call. If there were no transparent retries, 0 is not reported. |
| grpc.client.call.hedges | Histogram | {hedge} | grpc.method, grpc.target, grpc.client.call.custom (optional) | Number of hedges during the client call. If there were no hedges, 0 is not reported. |
| grpc.client.call.retry_delay | Histogram | s | grpc.method, grpc.target, grpc.client.call.custom (optional) | Total time of delay while there is no active attempt during the client call. |

Refer A96: OTel Metrics for Retries for details.

#### Server Instruments

| Name | Type | Unit | Labels (required) | Description |
|---|---|---|---|---|
| grpc.server.call. started | Counter | {call} | grpc.method | The total number of RPCs started, including those that have not completed. |
| grpc.server.call. sent_total_compressed_message_size | Histogram | By | grpc.method, grpc.status | Total bytes (compressed but not encrypted) sent across all response messages (metadata excluded) per RPC; does not include grpc or transport framing bytes. |
| grpc.server.call. rcvd_total_compressed_message_size | Histogram | By | grpc.method, grpc.status | Total bytes (compressed but not encrypted) received across all request messages (metadata excluded) per RPC; does not include grpc or transport framing bytes. |
| grpc.server.call. duration | Histogram | s | grpc.method, grpc.status | This metric aims to measure the end2end time an RPC takes from the server transport’s (HTTP2/ inproc) perspective. |

Refer A66: OpenTelemetry Metrics for details.

### LB Policy Instruments

#### Weighted Round Robin LB Policy Instruments

| Name | Type | Unit | Labels (disposition) | Description |
|---|---|---|---|---|
| grpc.lb.wrr. rr_fallback | Counter | {update} | grpc.target (required), grpc.lb.locality (optional), grpc.lb.backend_service (optional) | EXPERIMENTAL: Number of scheduler updates in which there were not enough endpoints with valid weight, which caused the WRR policy to fall back to RR behavior. |
| grpc.lb.wrr. endpoint_weight_not_yet_usable | Counter | {endpoint} | grpc.target (required), grpc.lb.locality (optional), grpc.lb.backend_service (optional) | EXPERIMENTAL: Number of endpoints from each scheduler update that don’t yet have usable weight information (i.e., either the load report has not yet been received, or it is within the blackout period). |
| grpc.lb.wrr. endpoint_weight_stale | Counter | {endpoint} | grpc.target (required), grpc.lb.locality (optional), grpc.lb.backend_service (optional) | EXPERIMENTAL: Number of endpoints from each scheduler update whose latest weight is older than the expiration period. |
| grpc.lb.wrr. endpoint_weights | Histogram | {weight} | grpc.target (required), grpc.lb.locality (optional), grpc.lb.backend_service (optional) | EXPERIMENTAL: Weight of an endpoint recorded every scheduler update. |

Refer A78: gRPC OTel Metrics for WRR, Pick First, and XdsClient for details.

#### Pick First LB Policy Instruments

| Name | Type | Unit | Labels (required) | Description |
|---|---|---|---|---|
| grpc.lb.pick_first. disconnections | Counter | {disconnection} | grpc.target | EXPERIMENTAL: Number of times the selected subchannel becomes disconnected. |
| grpc.lb.pick_first. connection_attempts_succeeded | Counter | {attempt} | grpc.target | EXPERIMENTAL: Number of successful connection attempts. |
| grpc.lb.pick_first. connection_attempts_failed | Counter | {attempt} | grpc.target | EXPERIMENTAL: Number of failed connection attempts. |

Refer A78: gRPC OTel Metrics for WRR, Pick First, and XdsClient for details.

### XdsClient Instruments

| Name | Type | Unit | Labels (required) | Description |
|---|---|---|---|---|
| grpc.xds_client. connected | Gauge | {bool} | grpc.target, grpc.xds.server | EXPERIMENTAL: Whether or not the xDS client currently has a working ADS stream to the xDS server. |
| grpc.xds_client. server_failure | Counter | {failure} | grpc.target, grpc.xds.server | EXPERIMENTAL: A counter of xDS servers going from healthy to unhealthy. |
| grpc.xds_client. resource_updates_valid | Counter | {resource} | grpc.target, grpc.xds.server, grpc.xds.resource_type | EXPERIMENTAL: A counter of resources received that were considered valid, even if unchanged. |
| grpc.xds_client. resource_updates_invalid | Counter | {resource} | grpc.target, grpc.xds.server, grpc.xds.resource_type | EXPERIMENTAL: A counter of resources received that were considered invalid. |
| grpc.xds_client. resources | Gauge | {resource} | grpc.target, grpc.xds.authority, grpc.xds.cache_state, grpc.xds.resource_type | EXPERIMENTAL: Number of xDS resources. |

Refer A78: gRPC OTel Metrics for WRR, Pick First, and XdsClient for details.

### Labels/Attributes

With a recorded measurement for an instrument, gRPC might provide some
additional information as attributes or labels. For example,
`grpc.client.attempt.started` has the labels `grpc.method` and `grpc.target`
along with each measurement that tell us the method and the target associated
with the RPC attempt being observed.

Some attributes are marked as optional on the instruments. These need to be explicitly enabled from the gRPC OpenTelemetry Plugin API. (Reference C++ API)NOTE

| Name | Description |
|---|---|
| grpc.method | Full gRPC method name, including package, service and method, e.g. “google.bigtable.v2.Bigtable/CheckAndMutateRow”. |
| grpc.status | gRPC server status code received, e.g. “OK”, “CANCELLED”, “DEADLINE_EXCEEDED”. |
| grpc.target | Canonicalized target URI used when creating gRPC Channel, e.g. “dns:///pubsub.googleapis.com:443”, “xds:///helloworld-gke:8000”. |
| grpc.client.call.custom | EXPERIMENTAL: Client-provided string for custom application use |
| grpc.lb.backend_service | The backend service to which the traffic is being sent. This is relevant when a single channel target can be sent to different sets of servers. When using xDS, this will be the cluster name. When not relevant, the value will be the empty string. |
| grpc.lb.locality | The locality to which the traffic is being sent. |
| grpc.xds.server | For clients, indicates the target of the gRPC channel in which the XdsClient is used. For servers, will be the string “#server”. |
| grpc.xds.authority | The xDS authority. The value will be “#old” for old-style non-xdstp resource names. |
| grpc.xds.cache_state | Indicates the cache state of an xDS resource (“requested”, “does_not_exist”, “acked”, “nacked”, “nacked_but_cached”). |
| grpc.xds.resource_type | xDS resource type, such as “envoy.config.listener.v3.Listener”. |

## FAQ

#### Q. How do I get throughput or QPS (queries per second)?

Use a count aggregation on the latency histogram metrics:
`grpc.client.attempt.duration` / `grpc.client.call.duration` (for clients) or
`grpc.server.call.duration` (for servers).

#### Q. How do I get error rate for RPCs?

Error counts can be calculated by using a filter `grpc.status != OK` value on
the latency histogram metrics `grpc.client.attempt.duration` /
`grpc.client.call.duration` (for clients) or `grpc.server.call.duration` (for
servers).

## Language examples

| Language | Example |
|---|---|
| C++ | C++ Example |
| Go | Go Example |
| Java | Java Example |
| Python | Python Example |

### Additional Resources

- A66: OpenTelemetry Metrics
- A78: gRPC OTel Metrics for WRR, Pick First, and XdsClient
- A79: Non-per-call Metrics Architecture
- A96: OTel Metrics for Retries

# Performance Best Practices

A user guide of both general and language-specific best practices to improve performance.

# Performance Best Practices

### General

- Always - **re-use stubs and channels**when possible.
- **Use keepalive pings**to keep HTTP/2 connections alive during periods of inactivity to allow initial RPCs to be made quickly without a delay (i.e. C++ channel arg GRPC_ARG_KEEPALIVE_TIME_MS).
- **Use streaming RPCs**when handling a long-lived logical flow of data from the client-to-server, server-to-client, or in both directions. Streams can avoid continuous RPC initiation, which includes connection load balancing at the client-side, starting a new HTTP/2 request at the transport layer, and invoking a user-defined method handler on the server side.- Streams, however, cannot be load balanced once they have started and can be hard to debug for stream failures. They also might increase performance at a small scale but can reduce scalability due to load balancing and complexity, so they should only be used when they provide substantial performance or simplicity benefit to application logic. Use streams to optimize the application, not gRPC. - **Side note:**- *This does not apply to Python (see Python section for details).*
- *(Special topic)*Each gRPC channel uses 0 or more HTTP/2 connections and each connection usually has a limit on the number of concurrent streams. When the number of active RPCs on the connection reaches this limit, additional RPCs are queued in the client and must wait for active RPCs to finish before they are sent. Applications with high load or long-lived streaming RPCs might see performance issues because of this queueing. There are two possible solutions:- **Create a separate channel for each area of high load**in the application.
- **Use a pool of gRPC channels**to distribute RPCs over multiple connections (channels must have different channel args to prevent re-use so define a use-specific channel arg such as channel number).
 - **Side note:**- *The gRPC team has plans to add a feature to fix these performance issues (see grpc/grpc#21386 for more info), so any solution involving creating multiple channels is a temporary workaround that should eventually not be needed.*

### C++

- **Do not use Sync API for performance sensitive servers.**If performance and/or resource consumption are not concerns, use the Sync API as it is the simplest to implement for low-QPS services.
- **Favor callback API over other APIs for most RPCs**, given that the application can avoid all blocking operations or blocking operations can be moved to a separate thread. The callback API is easier to use than the completion-queue async API but is currently slower for truly high-QPS workloads.
- If having to use the async completion-queue API, the - **best scalability trade-off is having**The ideal number of completion queues in relation to the number of threads can change over time (as gRPC C++ evolves), but as of gRPC 1.41 (Sept 2021), using 2 threads per completion queue seems to give the best performance.- `numcpu`’s threads.
- For the async completion-queue API, make sure to - **register enough server requests for the desired level of concurrency**to avoid the server continuously getting stuck in a slow path that results in essentially serial request processing.
- *(Special topic)*gRPC::GenericStub can be useful in certain cases when there is high contention / CPU time spent on proto serialization. This class allows the application to directly send- **raw gRPC::ByteBuffer as data**rather than serializing from some proto. This can also be helpful if the same data is being sent multiple times, with one explicit proto-to-ByteBuffer serialization followed by multiple ByteBuffer sends.

### Java

- **Use non-blocking stubs**to parallelize RPCs.
- **Provide a custom executor that limits the number of threads, based on your workload**(cached (default), fixed, forkjoin, etc).

### Python

- Streaming RPCs create extra threads for receiving and possibly sending the messages, which makes - **streaming RPCs much slower than unary RPCs**in gRPC Python, unlike the other languages supported by gRPC.
- **Using asyncio**could improve performance.
- Using the future API in the sync stack results in the creation of an extra thread. - **Avoid the future API**if possible.
- *(Experimental)*An experimental- **single-threaded unary-stream implementation**is available via the SingleThreadedUnaryStream channel option, which can save up to 7% latency per message.

# Reflection

Explains how reflection can be used to improve the transparency and interpretability of RPCs.

# Reflection

### Overview

Reflection is a protocol that gRPC servers can use to declare the protobuf-defined APIs they export over a standardized RPC service, including all types referenced by the request and response messages. Clients can then use this information to encode requests and decode responses in a human-readable manner.

Reflection is used heavily by debugging tools such as
`grpcurl` and
Postman.
One coming from the REST world might compare the gRPC reflection API to serving
an OpenAPI document on the HTTP server presenting the REST API being described.

### Transparency and Interpretability

A big contributor to gRPC’s stellar performance is the use of Protobuf for
serialization – a *binary* non-human-readable protocol. While this greatly
speeds up an RPC, it can also make it more difficult to manually interact with a
server. Hypothetically, in order to manually send a gRPC request to a server
over HTTP/2 using `curl`, you would have to:

- Know which RPC services the server exposed.
- Know the protobuf definition of the request message and all the types *it*references.
- Know the protobuf definition of the response message and all the types *it*references.

Then, you’d have to use that knowledge to hand-craft your request message(s) into binary and painstakingly decode the response message(s). This would be time consuming, frustrating, and error prone. Instead, the reflection protocol enables tools to automate this whole process, making it invisible.

### Enabling Reflection on a gRPC Server

Reflection is *not* automatically enabled on a gRPC server. The server author
must call a few additional functions to add a reflection service. These API calls
differ slightly from language to language and, in some languages, require adding
a dependency on a separate package, named something like `grpc-reflection`

Follow these links below for details on your specific language:

| Language | Guide |
|---|---|
| Java | Java example |
| Go | Go example |
| C++ | C++ example |
| Python | Python example |
| Javascript | Javascript example |

### Tips

Reflection works so seamlessly with tools such as `grpcurl` that oftentimes,
people aren’t even aware that it’s happening under the hood. However, if
reflection isn’t exposed, things won’t work seamlessly at all. Instead, the
client will fail with nasty errors. People often run into this when writing the
routing configuration for a gRPC service. The *reflection* service must be
routed to the appropriate backend as well as the application’s main RPC service.

If your gRPC API is accessible to public users, you may *not* want to expose the
reflection service, as you may consider this a security issue. Ultimately, you
will need to make a call here that strikes the best balance between security and
ease-of-use for you and your users.

# Request Hedging

Explains what request hedging is and how you can configure it.

# Request Hedging

### Overview

Hedging is one of two configurable retry policies supported by gRPC. With hedging, a gRPC client sends multiple copies of the same request to different backends and uses the first response it receives. Subsequently, the client cancels any outstanding requests and forwards the response to the application.

Hedging is a technique to reduce tail latency in large scale distributed
systems. While naive implementations could add significant load to the backend
servers, it is possible to get most of the latency reduction effects while
increasing load only modestly. For an in-depth discussion on tail latencies, see the seminal article, The Tail
At Scale, by Jeff Dean and Luiz André
Barroso. Hedging is configurable via gRPC Service Config, at a per-method granularity.
The configuration contains the following knobs: When the application makes an RPC call that contains a When a successful response is received (in response to any of the hedged
requests), all outstanding hedged requests are canceled and the response is
returned to the client application layer. If an error response with a non-fatal status code (controlled by the
 If all instances of a hedged RPC fail, there are no additional retry attempts.
Essentially, hedging can be seen as retrying the original RPC before a failure
is even received. If server pushback that specifies not to retry is received in response to a
hedged request, no further hedged requests should be issued for the call. gRPC provides a way to throttle hedged RPCs to prevent server overload.
Throttling can be configured via the Service Config as well using the
 For each server name, the gRPC client maintains a With hedging, the first request is always sent out, but subsequent hedged
requests are sent only if The only requests that are counted as failures for the throttling policy are the
ones that fail with a status code that qualifies as a non-fatal status code, or
that receive a pushback response indicating not to retry. This avoids conflating
server failure with responses to malformed requests (such as the
 Servers may explicitly pushback by setting metadata in their response to the
client. If the pushback says not to retry, no further hedged requests will be
sent. If the pushback says to retry after a given delay, the next hedged request
(if any) will be issued after the given delay has elapsed. Server pushback is specified using the metadata key, ## Use cases

## Configuring hedging in gRPC

```
"hedgingPolicy": {
 "maxAttempts": INTEGER,
 "hedgingDelay": JSON proto3 Duration type,
 "nonFatalStatusCodes": JSON array of grpc status codes (int or string)
}
```
`maxAttempts`: maximum number of in-flight requests while waiting for a
successful response. This is a mandatory field, and must be specified. If the
specified value is greater than `5`, gRPC uses a value of `5`.`hedgingDelay`: amount of time that needs to elapse before the client sends out
the next request while waiting for a successful response. This field is
optional, and if left unspecified, results in `maxAttempts` number of requests
all sent out at the same time.`nonFatalStatusCodes`: an optional list of grpc status codes. If any of hedged
requests fails with a status code that is not present in this list, all
outstanding requests are canceled and the response is returned to the
application.## Hedging policy

`hedgingPolicy`
configuration in the Service Config, the original RPC is sent immediately, as
with a standard non-hedged call. After `hedgingDelay` has elapsed without a
successful response, the second RPC will be issued. If neither RPC has received
a response after `hedgingDelay` has elapsed again, a third RPC is sent, and so
on, up to `maxAttempts`. gRPC call deadlines apply to the entire chain of hedged
requests. Once the deadline has passed, the operation fails regardless of
in-flight RPCS, and regardless of the hedging configuration.`nonFatalStatusCodes` field) is received from a hedged request, then the next
hedged request in line is sent immediately, shortcutting its hedging delay. If
any other status code is received, all outstanding RPCs are canceled and the
error is returned to the client application layer.## Throttling Hedged RPCs

`RetryThrottlingPolicy` message. The throttling configuration contains the
following:```
"retryThrottling": {
 "maxTokens": 10,
 "tokenRatio": 0.1
}
```
`token_count` which is
initially set to `max_tokens`. Every outgoing RPC (regardless of service or
method invoked) changes `token_count` as follows:`token_count` by `1`.`token_count` by `token_ratio`.`token_count` is greater than the threshold (defined
as `max_tokens / 2`). If `token_count` is less than or equal to the threshold,
hedged requests do not block. Instead they are canceled, and if there are no
other already-sent hedged RPCs the failure is returned to the client
application.`INVALID_ARGUMENT` status code).## Server Pushback

`grpc-retry-pushback-ms`.
The value is an ASCII encoded signed 32-bit integer with no unnecessary leading
zeros that represents how many milliseconds to wait before sending the next
hedged request. If the value for pushback is negative or unparseble, then it
will be seen as the server asking the client not to retry at all.## Resources

## Language Support

Language Example Java Java example C++ Not yet available Go Not yet supported

# Retry

gRPC takes the stress out of failures! Get fine-grained retry control and detailed insights with OpenCensus and OpenTelemetry support.

# Retry

### Overview

Retries are a key pattern for making services more reliable. By re-attempting failed operations, applications can overcome temporary issues like network or server glitches. This is essential for modern cloud applications to handle the inevitable transient faults that occur.

For the best practice, applications should understand what failed operations are suitable for retry, define exponential backoff parameters for retry delay, determine the number of retry attempts, and also monitor retry metrics.

### How gRPC client retry works

gRPC’s built-in retry logic saves the call’s history for potential retries and monitors RPC events. Even if there is no retry policy configured, gRPC still saves the call’s history in case it needs to perform transparent retry (discussed in a later section). Note that ‘retry’ means replacing a failed call with a new call and replaying the call’s history on that newly created call.

If certain criteria are met – the RPC closes with a failure status code matching the retry policy’s retryable status codes and remains within the retry attempt limit – gRPC will create a new retry stream after an exponential backoff delay.

gRPC also supports other features like retry throttling and server push back. See gRFC for client side retry for further details.

Once the response header is received, the RPC is committed. No further retries will be attempted, and gRPC hands over the RPC to the application.

The graph below shows architectural overview of gRPC retry internal.

```
sequenceDiagram
 Application ->> gRPC Client: Configure retry policy. <br> Send request to dns:///my-service
 gRPC Client ->> gRPC Client: Save message
 gRPC Client ->> Server: Create initial attempt
 Server -->> gRPC Client : RPC closed with error
 gRPC Client ->> Server: Create retry attempt 1
 Server -->> gRPC Client: Successful
 gRPC Client ->> Application: No more retry. Proceed.
```
### Retry configuration

Retries are enabled by default, but there is no default retry policy. Without a retry policy, gRPC cannot safely retry RPCs in most cases. Only RPCs that failed due to low-level races are retried, and only if gRPC is certain the RPCs have not been processed by a server. This is known as “transparent retry.” You can configure a retry policy to allow gRPC to retry RPCs in more circumstances and more aggressively. You can also disable retries entirely when creating a channel, which disables transparent retries and any configured retry policies.

#### Transparent Retry

Failure can occur in different stages. Even without an explicit retry policy, gRPC may perform transparent retries. The extent of these retries depends on when the failure happens:

- gRPC may do unlimited transparent retry when RPC never leaves the client.
- gRPC performs a single transparent retry when RPC reaches the gRPC server library, but has never been seen by the server application logic. Be aware of this type of retry, as it adds load to the network.

You can optimize your application’s retry functionality by focusing on key steps and configurations that gRPC supports.

- Max number of retry attempts
- Exponential backoff
- Set of retryable status codes

Retry is configurable via gRPC Service Config, at a per-method granularity. The configuration contains the following knobs:

```
"retryPolicy": {
 "maxAttempts": 4,
 "initialBackoff": "0.1s",
 "maxBackoff": "1s",
 "backoffMultiplier": 2,
 "retryableStatusCodes": [
 "UNAVAILABLE"
 ]
}
```
Jitter of plus or minus 20% is applied to the backoff delay to avoid hammering servers at the same time from a large number of clients. In the example configuration above, `initialBackoff` is set to 100ms, so the actual backoff delay after the first attempt will be for a random time period within the range `[80ms, 120ms]`.

gRPC supports throttle limit that prevents server overload due to retries. Below is an examples of retry throttle configuration:

```
"retryThrottling": {
 "maxTokens": 10,
 "tokenRatio": 0.1
}
```
For each server, the gRPC client tracks a `token_count` (initially set to `maxTokens`). Failed RPCs decrement the count by 1, successful RPCs increment it by `tokenRatio`. If the `token_count` falls below half of `maxTokens`, retries are paused until the count recovers.

Further, hedging is a complementary feature to retries and can be configured similarly. For more details, see the hedging guide.

### Retry Observability

gRPC supports exposing OpenCensus and OpenTelemetry metrics when retry functionality is enabled. Here’s an example of the OpenTelemetry retry attempt statistics available:

- `grpc.client.attempt.started`
- `grpc.client.attempt.duration`
- `grpc.client.attempt.sent_total_compressed_message_size`
- `grpc.client.attempt.rcvd_total_compressed_message_size`

Metrics at per cal level:

- `grpc.client.call.duration`

And server side metrics:

- `grpc.server.call.started`
- `grpc.server.call.sent_total_compressed_message_size`
- `grpc.server.call.rcvd_total_compressed_message_size`
- `grpc.server.call.duration`

Find in-depth metrics and tracing information, along with configuration instructions, in the gRFC for Otel metrics, gRFC for retry status.

### Language guides and examples

| Language | Example | Documentation |
|---|---|---|
| C++ | C++ Example | |
| Go | Go Example | |
| Java | Java Example | Java Documentation |
| Python | Python Example |

### Additional Resources

- gRFC for client side retry
- gRFC for retry status
- Hedging Guide
- gRPC Service Config
- gRFC for Otel metrics

# Service Config

How the service config can be used by service owners to control client behavior.

# Service Config

### Overview

The service config specifies how gRPC clients should behave when interacting with a gRPC server. Service owners can provide a service config with expected behavior of all service clients. The settings in a service config always apply to a specific target string (e.g. “api.myapp.com”), not globally.

### Behavior controlled by the Service Config

The settings in the service config affect client side load balancing, call behavior and health checking.

This page outlines the options in the service config, but the full service config data structure is documented with a protobuf definition.

#### Load Balancing

A service can be composed of multiple servers and the load balancing
configuration specifies how calls from clients should be distributed among
those servers. By default the `pick_first` load balancing policy is utilized,
but another policy can be specified in the service config. E.g. specifying the
`round_robin` policy will make the clients rotate through the servers instead
of repeatedly using the first server.

#### Call Behavior

RPCs can be configured in many ways:

- With wait-for-ready enabled, if a client cannot connect to a backend, the RPC will be delayed instead of immediately failing.
- A call timeout can be provided, indicating the maximum time the client should wait before giving up on the RPC.
- One of:

#### Note

These call behavior settings can be limited to an individual service or a method.

Retry and hedging policies can be further adjusted by setting a *retry
throttling policy* but it will apply across all services and methods.

#### Health Checking

A client can be configured to perform health checking by providing a health checking name. The client will then use the standard gRPC health checking service.

### Acquiring a Service Config

A service config can be provided to a client either via name resolution or programmatically by the client application.

#### Name Resolution

The gRPC name resolution mechanism allows for pluggable name resolver implementations. These implementations return the addresses associated with a name as well as an associated service config. This is the mechanism that service owners can use to distribute their service config out to a fleet of gRPC clients.

- The xDS name resolver converts the xDS configuration it receives from the control plane to a corresponding service config.
- The standard DNS name resolver in the Go implementation supports service configs stored as TXT records on the name server.

#### Note

Even though the service config structure is documented with a protobuf definition the internal representation in the client is JSON. Name resolver implementations are free to store the service config information in any way they prefer as long as they provide it in JSON format at name resolution time.#### Programmatically

The gRPC client API provides a way to specify a service config in JSON format. This is used to provide a default service config that will be used in situations where the name resolver does not provide a service config. It can also be useful in some testing situations.

### Example Service Config

The below example does the following:

- Enables the `round_robin`load balancing policy.
- Sets a default call timeout of 1s that applies to all methods in all services.
- Overrides that timeout to be 2s for the `bar`method in the`foo`service as well as all the methods in the`baz`service.

```
{
 "loadBalancingConfig": [ { "round_robin": {} } ],
 "methodConfig": [
 {
 "name": [{}],
 "timeout": "1s"
 },
 {
 "name": [
 { "service": "foo", "method": "bar" },
 { "service": "baz" }
 ],
 "timeout": "2s"
 }
 ]
}
```

# Status Codes

Explains the status codes used in gRPC.

# Status Codes

### Overview

All RPCs will result in a `status` being returned to the client. A `status` object is composed of an integer
code and a string error description. The server-side (or the gRPC library for library level errors) chooses
the status it returns for a given RPC. Applications should only use values defined below.

When an error situation occurs, the gRPC library may produce a corresponding `status`. The library may do this
either on the client- or the server-side. Only a subset of the pre-defined status codes are generated by the gRPC
libraries. This allows applications to be sure that any other code it sees was actually
returned by the application (although it is also possible for the
server-side to return one of the codes generated by the gRPC libraries).

See the Error handling user guide for how to use status codes.

### Full list of Status codes

gRPC uses a set of well defined status codes as part of the RPC API.

The following status codes are never generated by the library, only by user code:

- INVALID_ARGUMENT
- NOT_FOUND
- ALREADY_EXISTS
- FAILED_PRECONDITION
- ABORTED
- OUT_OF_RANGE
- DATA_LOSS

#### The full list of status codes

| Code | Id | Description |
|---|---|---|
| OK | 0 | Not an error; returned on success. |
| CANCELLED | 1 | The operation was cancelled, typically by the caller. |
| UNKNOWN | 2 | Unknown error. For example, this error may be returned when a `Status`value received from another address space belongs to an error space that is not known in this address space. Also errors raised by APIs that do not return enough error information may be converted to this error. |
| INVALID_ARGUMENT | 3 | The client specified an invalid argument. Note that this differs from `FAILED_PRECONDITION`.`INVALID_ARGUMENT`indicates arguments that are problematic regardless of the state of the system (e.g., a malformed file name). |
| DEADLINE_EXCEEDED | 4 | The deadline expired before the operation could complete. For operations that change the state of the system, this error may be returned even if the operation has completed successfully. For example, a successful response from a server could have been delayed long enough for the deadline to expire. |
| NOT_FOUND | 5 | Some requested entity (e.g., file or directory) was not found. Note to server developers: if a request is denied for an entire class of users, such as gradual feature rollout or undocumented allowlist, `NOT_FOUND`may be used. If a request is denied for some users within a class of users, such as user-based access control,`PERMISSION_DENIED`must be used. |
| ALREADY_EXISTS | 6 | The entity that a client attempted to create (e.g., file or directory) already exists. |
| PERMISSION_DENIED | 7 | The caller does not have permission to execute the specified operation. `PERMISSION_DENIED`must not be used for rejections caused by exhausting some resource (use`RESOURCE_EXHAUSTED`instead for those errors).`PERMISSION_DENIED`must not be used if the caller can not be identified (use`UNAUTHENTICATED`instead for those errors). This error code does not imply the request is valid or the requested entity exists or satisfies other pre-conditions. |
| RESOURCE_EXHAUSTED | 8 | Some resource has been exhausted, perhaps a per-user quota, or perhaps the entire file system is out of space. |
| FAILED_PRECONDITION | 9 | The operation was rejected because the system is not in a state required for the operation’s execution. For example, the directory to be deleted is non-empty, an rmdir operation is applied to a non-directory, etc. Service implementors can use the following guidelines to decide between `FAILED_PRECONDITION`,`ABORTED`, and`UNAVAILABLE`: (a) Use`UNAVAILABLE`if the client can retry just the failing call. (b) Use`ABORTED`if the client should retry at a higher level (e.g., when a client-specified test-and-set fails, indicating the client should restart a read-modify-write sequence). (c) Use`FAILED_PRECONDITION`if the client should not retry until the system state has been explicitly fixed. E.g., if an “rmdir” fails because the directory is non-empty,`FAILED_PRECONDITION`should be returned since the client should not retry unless the files are deleted from the directory. |
| ABORTED | 10 | The operation was aborted, typically due to a concurrency issue such as a sequencer check failure or transaction abort. See the guidelines above for deciding between `FAILED_PRECONDITION`,`ABORTED`, and`UNAVAILABLE`. |
| OUT_OF_RANGE | 11 | The operation was attempted past the valid range. E.g., seeking or reading past end-of-file. Unlike `INVALID_ARGUMENT`, this error indicates a problem that may be fixed if the system state changes. For example, a 32-bit file system will generate`INVALID_ARGUMENT`if asked to read at an offset that is not in the range [0,2^32-1], but it will generate`OUT_OF_RANGE`if asked to read from an offset past the current file size. There is a fair bit of overlap between`FAILED_PRECONDITION`and`OUT_OF_RANGE`. We recommend using`OUT_OF_RANGE`(the more specific error) when it applies so that callers who are iterating through a space can easily look for an`OUT_OF_RANGE`error to detect when they are done. |
| UNIMPLEMENTED | 12 | The operation is not implemented or is not supported/enabled in this service. |
| INTERNAL | 13 | Internal errors. This means that some invariants expected by the underlying system have been broken. This error code is reserved for serious errors. |
| UNAVAILABLE | 14 | The service is currently unavailable. This is most likely a transient condition, which can be corrected by retrying with a backoff. Note that it is not always safe to retry non-idempotent operations. |
| DATA_LOSS | 15 | Unrecoverable data loss or corruption. |
| UNAUTHENTICATED | 16 | The request does not have valid authentication credentials for the operation. |

# Wait-for-Ready

Explains how to configure RPCs to wait for the server to be ready before sending the request.

# Wait-for-Ready

### Overview

This is a feature which can be used on a stub which will cause the RPCs to wait for the server to become available before sending the request. This allows for robust batch workflows since transient server problems won’t cause failures. The deadline still applies, so the wait will be interrupted if the deadline is passed.

When an RPC is created when the channel has failed to connect to the server,
without Wait-for-Ready it will immediately return a failure; with Wait-for-Ready
it will simply be queued until the connection becomes ready. The default is
**without** Wait-for-Ready.

For detailed semantics see this.

### How to use Wait-for-Ready

You can specify for a stub whether or not it should use Wait-for-Ready, which will automatically be passed along when an RPC is created.

#### Note

The RPC can still fail for other reasons besides the server not being ready, so error handling is still necessary.The following shows the sequence of events that occur, when a client sends a message to a server, based upon channel state and whether or not Wait-for-Ready is set.

```
sequenceDiagram
participant A as Application
participant RPC
participant CH as Channel
participant S as Server
A->>RPC: Create RPC using stub
RPC->>CH: Initiate Communication
alt channel state: READY
 CH->>S: Send message
else Channel state: IDLE or CONNECTING
 CH-->>CH: Wait for state change
else Channel state: TRANSIENT_FAILURE
 alt with Wait-for-Ready
 CH-->>CH: Wait for channel<br>becoming READY<br>(or a permanent failure)
 CH->>S: Send message
 else without Wait-for-Ready
 CH->>A: Failure
 end
else Channel state is a Permanent Failure
 CH->>A: Failure
end
```
The following is a state based view

```
stateDiagram-v2
 state "Initiating Communication" as IC
 state "Channel State" as CS
 IC-->CS: Check Channel State
 state CS {
 state "Permanent Failure" as PF
 state "TRANSIENT_FAILURE" as TF
 IDLE --> CONNECTING
 CONNECTING --> READY
 READY-->[*]
 CONNECTING-->TF
 CONNECTING-->PF
 TF-->READY
 TF -->[*]: without\n wait-for-ready
 TF-->PF
 PF-->[*]
 }
 state "MSG sent" as MS
 state "RPC Failed" as RF
 CS-->WAIT:From IDLE /\nCONNECTING
 CS-->WAIT:From Transient\nFailure with\nWait-for-Ready
 WAIT-->CS:State Change
 CS-->MS: From READY
 CS-->RF: From Permanent failure or\nTransient Failure without\nWait-for-Ready
 MS-->[*]
 RF-->[*]
```
### Alternatives

- Loop (with exponential backoff) until the RPC stops returning transient failures.- This could be combined, for efficiency, with implementing an `onReady`Handler*(for languages that support this)*.

- This could be combined, for efficiency, with implementing an
- Accept failures that might have been avoided by waiting because you want to fail fast

### Language Support

| Language | Example |
|---|---|
| Java | Java example |
| Go | Go example |
| Python | Python example |
