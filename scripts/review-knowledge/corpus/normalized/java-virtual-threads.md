## Virtual Threads

Virtual threads are lightweight threads that reduce the effort of writing, maintaining, and debugging high-throughput concurrent applications.

For background information about virtual threads, see JEP 444.

A thread is the smallest unit of processing that can be scheduled. It runs concurrently with—and largely independently of—other such units. It's an instance of java.lang.Thread. There are two kinds of threads, platform threads and virtual threads.

### What is a Platform Thread?

A platform thread is implemented as a thin wrapper around an operating system (OS) thread. A platform thread runs Java code on its underlying OS thread, and the platform thread captures its OS thread for the platform thread's entire lifetime. Consequently, the number of available platform threads is limited to the number of OS threads.

Platform threads typically have a large thread stack and other resources that are maintained by the operating system. They are suitable for running all types of tasks but may be a limited resource.

### What is a Virtual Thread?

Like a platform thread, a virtual thread is also an instance of java.lang.Thread. However, a virtual thread isn't tied to a specific OS thread. A virtual thread still runs code on an OS thread. However, when code running in a virtual thread calls a blocking I/O operation, the Java runtime suspends the virtual thread until it can be resumed. The OS thread associated with the suspended virtual thread is now free to perform operations for other virtual threads.

Virtual threads are implemented in a similar way to virtual memory. To simulate a lot of memory, an operating system maps a large virtual address space to a limited amount of RAM. Similarly, to simulate a lot of threads, the Java runtime maps a large number of virtual threads to a small number of OS threads.

Unlike platform threads, virtual threads typically have a shallow call stack, performing as few as a single HTTP client call or a single JDBC query. Although virtual threads support thread-local variables and inheritable thread-local variables, you should carefully consider using them because a single JVM might support millions of virtual threads.

Virtual threads are suitable for running tasks that spend most of the time blocked, often waiting for I/O operations to complete. However, they aren't intended for long-running CPU-intensive operations.

### Why Use Virtual Threads?

Use virtual threads in high-throughput concurrent applications, especially those that consist of a great number of concurrent tasks that spend much of their time waiting. Server applications are examples of high-throughput applications because they typically handle many client requests that perform blocking I/O operations such as fetching resources.

Virtual threads are not faster threads; they do not run code any faster than platform threads. They exist to provide scale (higher throughput), not speed (lower latency).

### Creating and Running a Virtual Thread

The Thread and Thread.Builder APIs provide ways to create both platform and virtual threads. The java.util.concurrent.Executors class also defines methods to create an ExecutorService that starts a new virtual thread for each task.

#### Creating a Virtual Thread with the Thread Class and the Thread.Builder Interface

Call the Thread.ofVirtual() method to create an instance of Thread.Builder for creating virtual threads.

The following example creates and starts a virtual thread that prints a message. It calls the join method to wait for the virtual thread to terminate. (This enables you to see the printed message before the main thread terminates.)

```
Thread thread = Thread.ofVirtual().start(() -> System.out.println("Hello"));
thread.join();
```
The Thread.Builder interface lets you create threads with common Thread properties such as the thread's name. The Thread.Builder.OfPlatform subinterface creates platform threads while Thread.Builder.OfVirtual creates virtual threads.

The following example creates a virtual thread named
 `MyThread` with the Thread.Builder
 interface:

```
Thread.Builder builder = Thread.ofVirtual().name("MyThread");
Runnable task = () -> {
 System.out.println("Running thread");
};
Thread t = builder.start(task);
System.out.println("Thread t name: " + t.getName());
t.join();
```
The following example creates and starts two virtual threads with Thread.Builder:

```
Thread.Builder builder = Thread.ofVirtual().name("worker-", 0);
Runnable task = () -> {
 System.out.println("Thread ID: " + Thread.currentThread().threadId());
};
// name "worker-0"
Thread t1 = builder.start(task);
t1.join();
System.out.println(t1.getName() + " terminated");
// name "worker-1"
Thread t2 = builder.start(task);
t2.join();
System.out.println(t2.getName() + " terminated");
```
This example prints output similar to the following:

```
Thread ID: 21
worker-0 terminated
Thread ID: 24
worker-1 terminated
```
#### Creating and Running a Virtual Thread with the Executors.newVirtualThreadPerTaskExecutor() Method

Executors let you to separate thread management and creation from the rest of your application.

The following example creates an ExecutorService with the Executors.newVirtualThreadPerTaskExecutor() method. Whenever ExecutorService.submit(Runnable) is called, a new virtual thread is created and started to run the task. This method returns an instance of Future. Note that the method Future.get() waits for the thread's task to complete. Consequently, this example prints a message once the virtual thread's task is complete.

```
try (ExecutorService myExecutor = Executors.newVirtualThreadPerTaskExecutor()) {
 Future<?> future = myExecutor.submit(() -> System.out.println("Running thread"));
 future.get();
 System.out.println("Task completed");
 // ...
```
#### Multithreaded Client Server Example

The following example consists of two classes. `EchoServer` is
 a server program that listens on a port and starts a new virtual thread for each connection.
 `EchoClient` is a client program that connects to the server and sends
 messages entered on the command line.

`EchoClient` creates a socket, thereby getting a connection to
 `EchoServer`. It reads input from the user on the standard input
 stream, and then forwards that text to `EchoServer` by writing the text
 to the socket. `EchoServer` echoes the input back through the socket to
 the `EchoClient`. `EchoClient` reads and displays the data
 passed back to it from the server. `EchoServer` can service multiple
 clients simultaneously through virtual threads, one thread per each client
 connection.

```
public class EchoServer {

 public static void main(String[] args) throws IOException {

 if (args.length != 1) {
 System.err.println("Usage: java EchoServer <port>");
 System.exit(1);
 }

 int portNumber = Integer.parseInt(args[0]);
 try (
 ServerSocket serverSocket =
 new ServerSocket(Integer.parseInt(args[0]));
 ) {
 while (true) {
 Socket clientSocket = serverSocket.accept();
 // Accept incoming connections
 // Start a service thread
 Thread.ofVirtual().start(() -> {
 try (
 PrintWriter out =
 new PrintWriter(clientSocket.getOutputStream(), true);
 BufferedReader in = new BufferedReader(
 new InputStreamReader(clientSocket.getInputStream()));
 ) {
 String inputLine;
 while ((inputLine = in.readLine()) != null) {
 System.out.println(inputLine);
 out.println(inputLine);
 }

 } catch (IOException e) {
 e.printStackTrace();
 }
 });
 }
 } catch (IOException e) {
 System.out.println("Exception caught when trying to listen on port "
 + portNumber + " or listening for a connection");
 System.out.println(e.getMessage());
 }
 }
}
```
```
public class EchoClient {
 public static void main(String[] args) throws IOException {
 if (args.length != 2) {
 System.err.println(
 "Usage: java EchoClient <hostname> <port>");
 System.exit(1);
 }
 String hostName = args[0];
 int portNumber = Integer.parseInt(args[1]);
 try (
 Socket echoSocket = new Socket(hostName, portNumber);
 PrintWriter out =
 new PrintWriter(echoSocket.getOutputStream(), true);
 BufferedReader in =
 new BufferedReader(
 new InputStreamReader(echoSocket.getInputStream()));
 ) {
 BufferedReader stdIn =
 new BufferedReader(
 new InputStreamReader(System.in));
 String userInput;
 while ((userInput = stdIn.readLine()) != null) {
 out.println(userInput);
 System.out.println("echo: " + in.readLine());
 if (userInput.equals("bye")) break;
 }
 } catch (UnknownHostException e) {
 System.err.println("Don't know about host " + hostName);
 System.exit(1);
 } catch (IOException e) {
 System.err.println("Couldn't get I/O for the connection to " +
 hostName);
 System.exit(1);
 }
 }
}
```
### Scheduling Virtual Threads and Pinned Virtual Threads

The operating system schedules when a platform thread is run. However, the Java runtime schedules when a virtual thread is run. When the Java runtime schedules a virtual thread, it assigns or mounts the virtual thread on a platform thread, then the operating system schedules that platform thread as usual. This platform thread is called a carrier. After running some code, the virtual thread can unmount from its carrier. This usually happens when the virtual thread performs a blocking I/O operation. After a virtual thread unmounts from its carrier, the carrier is free, which means that the Java runtime scheduler can mount a different virtual thread on it.

A virtual thread cannot be unmounted during blocking operations when it is
 pinned to its carrier. A virtual thread is pinned when it
 runs a `native` method or a foreign function (see Foreign Function and Memory API). Pinning does not make an application incorrect, but
 it might hinder its scalability.

### Debugging Virtual Threads

Virtual threads are still threads; debuggers can step through them like
 platform threads. JDK Flight Recorder and the `jcmd` tool have
 additional features to help you observe virtual threads in your
 applications.

#### JDK Flight Recorder Events for Virtual Threads

JDK Flight Recorder (JFR) can emit these events related to virtual threads:

- `jdk.VirtualThreadStart`and- `jdk.VirtualThreadEnd`indicate when a virtual thread starts and ends. These events are disabled by default.
- `jdk.VirtualThreadPinned`indicates that a virtual thread was pinned (and its carrier thread wasn’t freed) for longer than the threshold duration. This event is enabled by default with a threshold of 20 ms.
- `jdk.VirtualThreadSubmitFailed`indicates that starting or unparking a virtual thread failed, probably due to a resource issue. Parking a virtual thread releases the underlying carrier thread to do other work, and unparking a virtual thread schedules it to continue. This event is enabled by default.

Enable the events `jdk.VirtualThreadStart` and
 `jdk.VirtualThreadEnd` through JDK Mission Control or with a custom
 JFR configuration as described in Flight Recorder Configurations in Java Platform, Standard Edition Flight Recorder API
 Programmer’s Guide.

To print these events, run the following command, where
 `recording.jfr` is the file name of your recording:

```
jfr print --events
 jdk.VirtualThreadStart,jdk.VirtualThreadEnd,jdk.VirtualThreadPinned,jdk.VirtualThreadSubmitFailed
 recording.jfr
```

#### Viewing Virtual Threads in jcmd Thread Dumps

You can create a thread dump in plain text as well as JSON format:

```
jcmd <PID> Thread.dump_to_file -format=text <file>
jcmd <PID> Thread.dump_to_file -format=json <file>
```
The `jcmd` thread dump lists all threads, including both
 platform and virtual threads. However, it doesn't include object addresses,
 JNI statistics, heap statistics, and other information that appears in
 traditional thread dumps.

### Virtual Threads: An Adoption Guide

Virtual threads are Java threads that are implemented by the Java runtime rather than the OS. The main difference between virtual threads and the traditional threads—which we've come to call platform threads—is that we can easily have a great many active virtual threads, even millions, running in the same Java process. It is their high number that gives virtual threads their power: they can run server applications written in the thread-per-request style more efficiently by allowing the server to process many more requests concurrently, leading to higher throughput and less waste of hardware.

Because virtual threads are an implementation of java.lang.Thread and conform to the same rules that specified java.lang.Thread since Java SE 1.0, developers don't need to learn new concepts to use them. However, the inability to spawn very many platform threads—the only implementation of threads available in Java for many years—has bred practices designed to cope with their high cost. These practices are counterproductive when applied to virtual threads, and must be unlearned. Moreover, the vast difference in cost informs a new way of thinking about threads that may be foreign at first.

This guide is not intended to be comprehensive and cover every important detail of virtual threads. It is meant but to provide an introductory set of guidelines to help those who wish to start using virtual threads make the best of them.

#### Write Simple, Synchronous Code Employing Blocking I/O APIs in the Thread-Per-Request Style

Virtual threads can significantly improve the throughput—not the latency—of servers written in the thread-per-request style. In this style, the server dedicates a thread to processing each incoming request for its entire duration. It dedicates at least one thread because, when processing a single request, you may want to employ more threads to carry some tasks concurrently.

Blocking a platform thread is expensive because it holds on to the thread—a relatively scarce resource—while it is not doing much meaningful work. Because virtual threads can be plentiful, blocking them is cheap and encouraged. Therefore, you should write code in the straightforward synchronous style and use blocking I/O APIs.

For example, the following code, written in the non-blocking, asynchronous style, won't benefit much from virtual threads.

```
CompletableFuture.supplyAsync(info::getUrl, pool)
 .thenCompose(url -> getBodyAsync(url, HttpResponse.BodyHandlers.ofString()))
 .thenApply(info::findImage)
 .thenCompose(url -> getBodyAsync(url, HttpResponse.BodyHandlers.ofByteArray()))
 .thenApply(info::setImageData)
 .thenAccept(this::process)
 .exceptionally(t -> { t.printStackTrace(); return null; });
```
On the other hand, the following code, written in the synchronous style and using simple blocking IO, will benefit greatly:

```
try {
 String page = getBody(info.getUrl(), HttpResponse.BodyHandlers.ofString());
 String imageUrl = info.findImage(page);
 byte[] data = getBody(imageUrl, HttpResponse.BodyHandlers.ofByteArray());
 info.setImageData(data);
 process(info);
} catch (Exception ex) {
 t.printStackTrace();
}
```
Such code is also easier to debug in a debugger, profile in a profiler, or
 observe with thread-dumps. To observe virtual threads, create a thread dump with the
 `jcmd` command:

`jcmd <pid> Thread.dump_to_file -format=json <file>`The more of the stack that's written in this style, the better virtual threads will be for both performance and observability. Programs or frameworks written in other styles that don't dedicate a thread per task should not expect to see a significant benefit from virtual threads. Avoid mixing synchronous, blocking code with asynchronous frameworks.

#### Represent Every Concurrent Task as a Virtual Thread; Never Pool Virtual Threads

The hardest thing to internalize about virtual threads is that, while they have the same behavior as platform threads they should not represent the same program concept.

Platform threads are scarce, and are therefore a precious resource. Precious resources need to be managed, and the most common way to manage platform threads is with thread pools. A question that you then need to answer is, how many threads should we have in the pool?

But virtual threads are plentiful, and so each should represent not some shared, pooled, resource but a task. From a managed resource threads turn into application domain objects. The question of how many virtual threads we should have becomes obvious, just as the question of how many strings we should use to store a set of user names in memory is obvious: The number of virtual threads is always equal to the number of concurrent tasks in your application.

Converting n platform threads to n virtual threads would yield little benefit; rather, it's tasks that need to be converted.

To represent every application task as a thread, don't use a shared thread pool executor like in the following example:

```
Future<ResultA> f1 = sharedThreadPoolExecutor.submit(task1);
Future<ResultB> f2 = sharedThreadPoolExecutor.submit(task2);
// ... use futures
```
Instead, use a virtual thread executor like in the following example:

```
try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
 Future<ResultA> f1 = executor.submit(task1);
 Future<ResultB> f2 = executor.submit(task2);
 // ... use futures
}
```
The code still uses an ExecutorService, but the one returned from Executors.newVirtualThreadPerTaskExecutor() doesn't employ a thread pool. Rather, it creates a new virtual thread for each submitted tasks.

Furthermore, that ExecutorService itself is lightweight, and we can create a new one just as we would with any simple object. That allows us to rely on the newly added ExecutorService.close() method and the try-with-resources construct. The close method, that is implicitly called at the end of the try block will automatically wait for all tasks submitted to the ExecutorService—that is, all virtual threads spawned by the ExecutorService—to terminate.

This is a particularly useful pattern for fanout scenarios, where you wish to concurrently perform multiple outgoing calls to different services like in the following example:

```
void handle(Request request, Response response) {
 var url1 = ...
 var url2 = ...

 try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
 var future1 = executor.submit(() -> fetchURL(url1));
 var future2 = executor.submit(() -> fetchURL(url2));
 response.send(future1.get() + future2.get());
 } catch (ExecutionException | InterruptedException e) {
 response.fail(e);
 }
}

String fetchURL(URL url) throws IOException {
 try (var in = url.openStream()) {
 return new String(in.readAllBytes(), StandardCharsets.UTF_8);
 }
}
```
You should create a new virtual thread, as shown above, for even small, short-lived concurrent tasks.

For even more help writing the fanout pattern and other common concurrency patterns, with better observability, use structured concurrency.

As a rule of thumb, if your application never has 10,000 virtual threads or more, it is unlikely to benefit from virtual threads. Either it experiences too light a load to need better throughput, or you have not represented sufficiently many tasks to virtual threads.

#### Use Semaphores to Limit Concurrency

Sometimes there is a need to limit the concurrency of a certain operation. For example, some external service may not be able to handle more than ten concurrent requests. Because platform threads are a precious resource that is usually managed in a pool, thread pools have become so ubiquitious that they're used for this purpose of restricting concurrency, like in the following example:

```
ExecutorService es = Executors.newFixedThreadPool(10);
...
Result foo() {
 try {
 var fut = es.submit(() -> callLimitedService());
 return f.get();
 } catch (...) { ... }
}
```
This example ensures that there are at most ten concurrent requests to the limited service.

But restricting concurrency is only a side-effect of thread pools' operation. Pools are designed to share scarce resources, and virtual threads aren’t scarce and therefore should never be pooled!

When using virtual threads, if you want to limit the concurrency of accessing some service, you should use a construct designed specifically for that purpose: the Semaphore class. The following example demonstrates this class:

```
Semaphore sem = new Semaphore(10);
...
Result foo() {
 sem.acquire();
 try {
 return callLimitedService();
 } finally {
 sem.release();
 }
}
```
Threads that happen to call `foo` will be throttled, that is, blocked, so
 that only ten of them can make progress at a time, while others will go about their
 business unencumbered.

Simply blocking some virtual threads with a semaphore may appear to be substantially different from submitting tasks to a fixed thread pool, but it isn't. Submitting tasks to a thread pool queues them up for later execution, but the semaphore internally (or any other blocking synchronization construct for that matter) creates a queue of threads that are blocked on it that mirrors the queue of tasks waiting for a pooled thread to execute them. Because virtual threads are tasks, the resulting structure is equivalent:

Figure 14-1 Comparing a Thread Pool with a Semaphore

Even though you can think of a pool of platform threads as workers processing tasks that they pull from a queue and of virtual threads as the tasks themselves, blocked until they may continue, the underlying representation in the computer is virtually identical. Recognizing the equivalence between queued tasks and blocked threads will help you make the most of virtual threads.

Database connection pools themselves serve as a semaphore. A connection pool limited to ten connections would block the eleventh thread attempting to acquire a connection. There is no need to add an additional semaphore on top of the connection pool.

#### Don't Cache Expensive Reusable Objects in Thread-Local Variables

Virtual threads support thread-local variables just as platform threads do. See Thread-Local Variables for more information. Usually, thread-local variables are used to associate some context-specific information with the currently running code, such as the current transaction and user ID. This use of thread-local variables is perfectly reasonable with virtual threads. However, consider using the safer and more efficient scoped values. See Scoped Values for more information.

There is another use of thread-local variables which is fundamentally at odds with virtual threads: caching reusable objects. These objects are typically expensive to create (and consume a significant amount of memory), are mutable, and not thread-safe. They are cached in a thread-local variable to reduce the number of times they are instantiated and their number of instances in memory, but they are reused by the multiple tasks that run on the thread at differerent times.

For example, an instance of SimpleDateFormat is expensive to create and isn't thread-safe. A pattern that emerged is to cache such an instance in a ThreadLocal like in the following example:

```
static final ThreadLocal<SimpleDateFormat> cachedFormatter =
 ThreadLocal.withInitial(SimpleDateFormat::new);
void foo() {
 ...
	cachedFormatter.get().format(...);
	...
}
```
This kind of caching is helpful only when the thread—and therefore the expensive object
 cached in the thread local—is shared and reused by multiple tasks, as would be the case
 when platform threads are pooled. Many tasks may call `foo` when running
 in the thread pool, but because the pool only contains a few threads, the object will
 only be instantiated a few times—once per pool thread—cached, and reused.

However, virtual threads are never pooled and never reused by unrelated tasks. Because
 every task has its own virtual threads, every call to `foo` from a
 different task would trigger the instantiation of a new
 SimpleDateFormat. Moreover, because there may be a great many
 virtual threads running concurrently, the expensive object may consume quite a lot of
 memory. These outcomes are the very opposite of what caching in thread locals intends to
 achieve.

There is no single general alternative to offer, but in the case of SimpleDateFormat, you should replace it with DateTimeFormatter. DateTimeFormatter is immutable, and so a single instance can be shared by all threads:

```
static final DateTimeFormatter formatter = DateTimeFormatter….;
void foo() {
 ...
	formatter.format(...);
	...
}
```
Note that using thread-local variables to cache shared expensive objects is sometimes done behind the scenes by asynchronous frameworks, under their implicit assumption that they are used by a very small number of pooled threads. This is one reason why mixing virtual threads and asynchronous frameworks is not a good idea: a call to a method may result in instantiating costly objects in thread-local variables that were intended to be cached and shared.
