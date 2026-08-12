This page explains how to receive and acknowledge messages using Pub/Sub's exactly-once feature, which makes it possible for you to track and prevent duplicate processing of messages. When the feature is enabled, Pub/Sub provides the following semantics:

- Subscribers can determine if message acknowledgments were successful.
- No redelivery occurs after the message is successfully acknowledged.
- No redelivery occurs while a message is outstanding. A message is considered outstanding until the acknowledgment deadline expires or the message is acknowledged.
- In case of multiple valid deliveries, due to acknowledgment deadline expiration or client-initiated negative acknowledgment, only the latest acknowledgment ID can be used to acknowledge the message. Any requests with a previous acknowledgment ID fail.

With exactly-once enabled, subscribers can ensure messages are processed one time by following these guidelines:

- Acknowledge messages within the acknowledgment deadline.
- Maintain information about the progress of processing a message until it is successfully acknowledged.
- Use the information about the progress of processing a message to prevent duplicate work when an acknowledgment fails.

Only the pull subscription type supports exactly-once delivery, including subscribers that use the StreamingPull API. Push and export subscriptions don't support exactly-once delivery.

Pub/Sub supports exactly-once delivery, within a cloud region, based on a Pub/Sub-defined unique message ID.

## Recommended client library versions

- For best performance, use the latest version of the client library, Python v2.13.6 or higher, Java v1.139.0 or higher, PHP v1.39.0 or higher, C# v3.2.0 or higher, C++ v2.1.0, Go v1.25.1 or higher, Node v3.2.0 or higher and Ruby v2.12.1 or higher.

## Redelivery versus duplicate

It is important to understand the difference between expected and unexpected redeliveries.

- A redelivery can happen either because of client-initiated negative acknowledgment of a message or when the client doesn't extend the acknowledgment deadline of the message before the acknowledgment deadline expires. Redeliveries are considered valid and system working as intended. - To troubleshoot redeliveries, see Dealing with duplicates.
- A duplicate is when a message is re-sent after a successful acknowledgment or before acknowledgment deadline expiration.
- A redelivered message retains the same message ID between redelivery attempts.

Subscriptions with exactly-once delivery enabled don't receive duplicate deliveries.

## Exactly-once delivery support in client libraries

- Supported client libraries have an interface for acknowledgment with response (example: Go). You can use this interface to check if the acknowledgment request succeeded. If the acknowledgment request succeeds, the clients are guaranteed to not receive a re-delivery. If the acknowledgment request fails, the clients can expect a re-delivery.
- Clients can also use the supported client libraries without the acknowledgment interface. However, in such cases, the acknowledgment failures can lead to silent re-deliveries of messages.
- Supported client libraries have interfaces for setting the minimum lease extension time (example: Go). You must set the value for the minimum lease extension to a high number to avoid any network-related acknowledgment expirations. The maximum value is set at 600 seconds.
- If you are using the Java client library and you initialize your subscriber with a custom gRPC Channel using the - `setChannelProvider()`method, it is recommended that you also set- `maxInboundMetadataSize`to at least 1MB when building your- `TransportChannelProvider`. For this configuration, you can use the- `InstantiatingGrpcChannelProvider.Builder.setMaxInboundMetadataSize()`or the- `ManagedChannelBuilder.maxInboundMetadataSize()`method.

The default values and range for the variables related to exactly-once delivery and the names of the variables might differ across client libraries. For example, in the Java client library, the following variables control exactly-once delivery.

| Variable | Description | Value |
|---|---|---|
| `setEnableExactlyOnceDelivery` | Enables or disables exactly-once delivery. | true or false Default=false |
| `minDurationPerAckExtension` | The minimum time in seconds to use for extending the modify acknowledgment deadline. | Range=0 to 600 Default=none |
| `maxDurationPerAckExtension` | The maximum time in seconds to use for extending the modify acknowledgment deadline. | Range=0 to 600 Default=none |

In the case of exactly-once delivery, the `modifyAckDeadline` or `acknowledgment`
request to Pub/Sub fails when the acknowledgment ID is already expired. In such
cases, the service considers the expired acknowledgment ID as invalid, since a
newer delivery might already be in-flight. This is by design for exactly-once
delivery. You then see `acknowledgment` and `ModifyAckDeadline` requests return an
`INVALID_ARGUMENT` response. When exactly-once delivery is disabled, these
requests return `OK` in cases of expired acknowledgment IDs.

To ensure that `acknowledgment` and `ModifyAckDeadline` requests have valid
acknowledgment IDs, consider setting the value for
`minDurationPerAckExtension` to a high number.

## Regional considerations

The exactly-once delivery guarantee only applies when subscribers connect to the service in the same region. If your subscriber application is spread across multiple regions, it can lead to duplicate message delivery, even when exactly-once delivery is enabled. Publishers can send messages to any region and the exactly once guarantee is still maintained.

When you run your application within Google Cloud, by default it connects to the Pub/Sub endpoint in the same region. Therefore, running your application in a single region within Google Cloud generally ensures you are interacting with a single region.

When you are running your subscriber application outside of Google Cloud or in multiple regions, you can guarantee you are connecting to a single region by using a locational endpoint when configuring your Pub/Sub client. All location endpoints for Pub/Sub point to single regions. To learn more about locational endpoints, see Pub/Sub endpoints. For a list of all locational endpoints for Pub/Sub, see List of locational endpoints.

## Create subscriptions with exactly-once delivery

You can create a subscription with exactly-once delivery using the Google Cloud console, the Google Cloud CLI, client library, or the Pub/Sub API.

### Pull subscription

### Console

To create a pull subscription with exactly-once delivery, follow these steps:

- In the Google Cloud console, go to the - **Subscriptions**page.
- Click - **Create subscription**.
- Enter the - **Subscription ID**.
- Choose or create a topic from the drop-down menu. - The subscription receives messages from the topic.
- In the - **Exactly once delivery**section, select- **Enable exactly once delivery**.
- Click - **Create**.

### gcloud

To create a pull subscription with exactly-once delivery, use the
`gcloud pubsub subscriptions create`
command with the `--enable-exactly-once-delivery` flag:

gcloud pubsub subscriptions createSUBSCRIPTION_ID\ --topic=TOPIC_ID\ --enable-exactly-once-delivery

Replace the following:

- `SUBSCRIPTION_ID`: the ID of the subscription to create
- `TOPIC_ID`: the ID of the topic to attach to the subscription

### REST

To create a subscription with exactly-once delivery, use the
`projects.subscriptions.create`
method.

PUT https://pubsub.googleapis.com/v1/projects/PROJECT_ID/subscriptions/SUBSCRIPTION_IDAuthorization: Bearer $(gcloud auth print-access-token)

Replace the following:

- `PROJECT_ID`: the project ID for the project to create the subscription in
- `SUBSCRIPTION_ID`: the ID of the subscription to create

To create a pull subscription with exactly-once delivery, specify this in the request body:

{ "topic": "projects/PROJECT_ID/topics/TOPIC_ID", "enableExactlyOnceDelivery": true, }

Replace the following:

- `PROJECT_ID`: the project ID for the project with the topic
- `TOPIC_ID`: the ID of the topic to attach to the subscription

### C++

Before trying this sample, follow the C++ setup instructions in Quickstart: Using Client Libraries. For more information, see the Pub/Sub C++ API reference documentation.

### C#

Before trying this sample, follow the C# setup instructions in Quickstart: Using Client Libraries. For more information, see the Pub/Sub C# API reference documentation.

### Go

The following sample uses the major version of the Go Pub/Sub client library (v2). If you are still using the v1 library, see the migration guide to v2. To see a list of v1 code samples, see the deprecated code samples.

Before trying this sample, follow the Go setup instructions in Quickstart: Using Client Libraries. For more information, see the Pub/Sub Go API reference documentation.

### Java

Before trying this sample, follow the Java setup instructions in Quickstart: Using Client Libraries. For more information, see the Pub/Sub Java API reference documentation.

### Python

Before trying this sample, follow the Python setup instructions in Quickstart: Using Client Libraries. For more information, see the Pub/Sub Python API reference documentation.

### Node.js

Before trying this sample, follow the Node.js setup instructions in Quickstart: Using Client Libraries. For more information, see the Pub/Sub Node.js API reference documentation.

### Node.js

Before trying this sample, follow the Node.js setup instructions in Quickstart: Using Client Libraries. For more information, see the Pub/Sub Node.js API reference documentation.

### Ruby

The following sample uses Ruby Pub/Sub client library v3. If you are still using the v2 library, see the migration guide to v3. To see a list of Ruby v2 code samples, see the deprecated code samples.

Before trying this sample, follow the Ruby setup instructions in Quickstart: Using Client Libraries. For more information, see the Pub/Sub Ruby API reference documentation.

### PHP

Before trying this sample, follow the PHP setup instructions in Quickstart: Using Client Libraries. For more information, see the Pub/Sub PHP API reference documentation.

## Monitor exactly-once delivery subscriptions

The
`subscription/exactly_once_warning_count`
metric records the number of events that
can lead to possible redeliveries (valid or duplicate). This metric counts the
times Pub/Sub fails to process requests associated with
acknowledgment IDs (`ModifyAckDeadline` or `acknowledgment` request). The reasons
for the failure could be server or client-based. For example, if the persistence
layer used to maintain the exactly-once delivery information is unavailable, it
would be a server-based event. If the client tries to acknowledge a message with
an invalid acknowledgment ID, it would be a client-based event.

## Understand the metric

`subscription/exactly_once_warning_count` captures events that might or might not
lead to actual redeliveries and can be noisy based on client behavior. For
example: repeated `acknowledgment` or `ModifyAckDeadline` requests with invalid
acknowledgment IDs increment the metric repeatedly.

The following metrics are also useful in understanding the client behavior:

- `subscription/expired_ack_deadlines_count`metric shows the number of acknowledgment ID expirations. Acknowledgment ID expirations can lead to failures for both- `ModifyAckDeadline`and- `acknowledgment`requests.
- `service.serviceruntime.googleapis.com/api/request_count`metric can be used to capture failures of- `ModifyAckDeadline`or- `acknowledgment`requests in cases where the requests reach Google Cloud but don't reach Pub/Sub. There are failures that this metric won't capture—for example, when clients are disconnected from Google Cloud.

In most cases of failure events that can be retried, supported client libraries retry the request automatically.

## Quotas

Exactly-once delivery subscriptions are subject to additional quota requirements. These quotas are enforced on:

- Number of messages consumed from subscriptions with exactly-once delivery enabled per region.
- Number of messages acknowledged or whose deadline is extended when using subscriptions with exactly-once delivery enabled per region.

For more information regarding these quotas, see the table in the Quotas topic.

## Exactly-once delivery and ordered subscriptions

Pub/Sub supports exactly-once delivery with ordered delivery.

When using ordering with exactly-once delivery, Pub/Sub expects the acknowledgments to be in order. If the acknowledgments are out-of-order, the service fails the requests with temporary errors. If the acknowledgment deadline expires before an in-order acknowledgment for the delivery, the client will receive a redelivery of the message. Due to this, when you use ordering with exactly-once delivery, the client throughput is limited to an order of thousands of messages per second.

## Exactly-once delivery and push subscriptions

Pub/Sub supports exactly-once delivery only with pull subscriptions.

Clients consuming messages from the push subscriptions acknowledge the messages by responding to the push requests with a successful response. However, clients don't know if the Pub/Sub subscription received the response and processed it. This is different from pull subscriptions, where acknowledgment requests are initiated by the clients and the Pub/Sub subscription responds if the request was successfully processed. Because of this, exactly-once delivery semantics don't align well with push subscriptions.

## Things to know

- If the acknowledgment deadline is not specified at CreateSubscription time, exactly-once delivery enabled subscriptions will have a default acknowledgment deadline of 60 seconds.
- Longer default acknowledgment deadlines are beneficial in avoiding redelivery caused by network events. Supported client libraries don't use the default subscription acknowledgment deadline.
- Exactly-once delivery subscriptions have significantly higher publish-to-subscribe latency compared to regular subscriptions.
- If you require high throughput, your exactly-once delivery clients must also use streaming pull.
- A subscription might receive multiple copies of the same message due to publish side duplicates, even with exactly-once delivery enabled. Publish-side duplicates can be due to multiple unique publish retries by the publishing client or the Pub/Sub service. Multiple unique publishes by the publishing client, across retries, leads to redeliveries with - *different*message IDs. Multiple unique publishes by the Pub/Sub service, to respond to a client publish request, leads to redeliveries with the- *same*message IDs.
- You can retry failures in - `subscription/exactly_once_warning_count`and the supported client libraries retry these automatically. However, failures related to invalid acknowledgment IDs cannot be retried.
