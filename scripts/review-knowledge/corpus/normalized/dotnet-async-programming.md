Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

The Task asynchronous programming (TAP) model provides a layer of abstraction over typical asynchronous coding. In this model, you write code as a sequence of statements, the same as usual. The difference is you can read your task-based code as the compiler processes each statement and before it starts processing the next statement. To accomplish this model, the compiler performs many transformations to complete each task. Some statements can initiate work and return a Task object that represents the ongoing work and the compiler must resolve these transformations. The goal of task asynchronous programming is to enable code that reads like a sequence of statements, but executes in a more complicated order. Execution is based on external resource allocation and when tasks complete.

The task asynchronous programming model is analogous to how people give instructions for processes that include asynchronous tasks. This article uses an example with instructions for making breakfast to show how the `async` and `await` keywords make it easier to reason about code that includes a series of asynchronous instructions. The instructions for making a breakfast might be provided as a list:

- Pour a cup of coffee.
- Heat a pan, then fry two eggs.
- Cook three hash brown patties.
- Toast two pieces of bread.
- Spread butter and jam on the toast.
- Pour a glass of orange juice.

If you have experience with cooking, you might complete these instructions **asynchronously**. You start warming the pan for eggs, then start cooking the hash browns. You put the bread in the toaster, then start cooking the eggs. At each step of the process, you start a task, and then transition to other tasks that are ready for your attention.

Cooking breakfast is a good example of asynchronous work that isn't parallel. One person (or thread) can handle all the tasks. One person can make breakfast asynchronously by starting the next task before the previous task completes. Each cooking task progresses regardless of whether someone is actively watching the process. As soon as you start warming the pan for the eggs, you can begin cooking the hash browns. After the hash browns start to cook, you can put the bread in the toaster.

For a parallel algorithm, you need multiple people who cook (or multiple threads). One person cooks the eggs, another cooks the hash browns, and so on. Each person focuses on their one specific task. Each person who is cooking (or each thread) is blocked synchronously waiting for the current task to complete: Hash browns ready to flip, bread ready to pop up in toaster, and so on.

Consider the same list of synchronous instructions written as C# code statements:

```
using System;
using System.Threading.Tasks;
namespace AsyncBreakfast
{
 // These classes are intentionally empty for the purpose of this example. They are simply marker classes for the purpose of demonstration, contain no properties, and serve no other purpose.
 internal class HashBrown { }
 internal class Coffee { }
 internal class Egg { }
 internal class Juice { }
 internal class Toast { }
 class Program
 {
 static void Main(string[] args)
 {
 Coffee cup = PourCoffee();
 Console.WriteLine("coffee is ready");
 Egg eggs = FryEggs(2);
 Console.WriteLine("eggs are ready");
 HashBrown hashBrown = FryHashBrowns(3);
 Console.WriteLine("hash browns are ready");
 Toast toast = ToastBread(2);
 ApplyButter(toast);
 ApplyJam(toast);
 Console.WriteLine("toast is ready");
 Juice oj = PourOJ();
 Console.WriteLine("oj is ready");
 Console.WriteLine("Breakfast is ready!");
 }
 private static Juice PourOJ()
 {
 Console.WriteLine("Pouring orange juice");
 return new Juice();
 }
 private static void ApplyJam(Toast toast) =>
 Console.WriteLine("Putting jam on the toast");
 private static void ApplyButter(Toast toast) =>
 Console.WriteLine("Putting butter on the toast");
 private static Toast ToastBread(int slices)
 {
 for (int slice = 0; slice < slices; slice++)
 {
 Console.WriteLine("Putting a slice of bread in the toaster");
 }
 Console.WriteLine("Start toasting...");
 Task.Delay(3000).Wait();
 Console.WriteLine("Remove toast from toaster");
 return new Toast();
 }
 private static HashBrown FryHashBrowns(int patties)
 {
 Console.WriteLine($"putting {patties} hash brown patties in the pan");
 Console.WriteLine("cooking first side of hash browns...");
 Task.Delay(3000).Wait();
 for (int patty = 0; patty < patties; patty++)
 {
 Console.WriteLine("flipping a hash brown patty");
 }
 Console.WriteLine("cooking the second side of hash browns...");
 Task.Delay(3000).Wait();
 Console.WriteLine("Put hash browns on plate");
 return new HashBrown();
 }
 private static Egg FryEggs(int howMany)
 {
 Console.WriteLine("Warming the egg pan...");
 Task.Delay(3000).Wait();
 Console.WriteLine($"cracking {howMany} eggs");
 Console.WriteLine("cooking the eggs ...");
 Task.Delay(3000).Wait();
 Console.WriteLine("Put eggs on plate");
 return new Egg();
 }
 private static Coffee PourCoffee()
 {
 Console.WriteLine("Pouring coffee");
 return new Coffee();
 }
 }
}
```
If you interpret these instructions as a computer would, breakfast takes about 30 minutes to prepare. The duration is the sum of the individual task times. The computer blocks for each statement until all work completes, and then it proceeds to the next task statement. This approach can take significant time. In the breakfast example, the computer method creates an unsatisfying breakfast. Later tasks in the synchronous list, like toasting the bread, don't start until earlier tasks complete. Some food gets cold before the breakfast is ready to serve.

If you want the computer to execute instructions asynchronously, you must write asynchronous code. When you write client programs, you want the UI to be responsive to user input. Your application shouldn't freeze all interaction while downloading data from the web. When you write server programs, you don't want to block threads that might be serving other requests. Using synchronous code when asynchronous alternatives exist hurts your ability to scale out less expensively. You pay for blocked threads.

Successful modern apps require asynchronous code. Without language support, writing asynchronous code requires callbacks, completion events, or other means that obscure the original intent of the code. The advantage of synchronous code is the step-by-step action that makes it easy to scan and understand. Traditional asynchronous models force you to focus on the asynchronous nature of the code, not on the fundamental actions of the code.

## Don't block, await instead

The previous code highlights an unfortunate programming practice: Writing synchronous code to perform asynchronous operations. The code blocks the current thread from doing any other work. The code doesn't interrupt the thread while there are running tasks. The outcome of this model is similar to staring at the toaster after you put in the bread. You ignore any interruptions and don't start other tasks until the bread pops up. You don't take the butter and jam out of the fridge. You might miss seeing a fire starting on the stove. You want to both toast the bread and handle other concerns at the same time. The same is true with your code.

You can start by updating the code so the thread doesn't block while tasks are running. The `await` keyword provides a nonblocking way to start a task, then continue execution when the task completes. A simple asynchronous version of the breakfast code looks like the following snippet:

```
static async Task Main(string[] args)
{
 Coffee cup = PourCoffee();
 Console.WriteLine("coffee is ready");
 Egg eggs = await FryEggsAsync(2);
 Console.WriteLine("eggs are ready");
 HashBrown hashBrown = await FryHashBrownsAsync(3);
 Console.WriteLine("hash browns are ready");
 Toast toast = await ToastBreadAsync(2);
 ApplyButter(toast);
 ApplyJam(toast);
 Console.WriteLine("toast is ready");
 Juice oj = PourOJ();
 Console.WriteLine("oj is ready");
 Console.WriteLine("Breakfast is ready!");
}
```
The code updates the original method bodies of `FryEggs`, `FryHashBrowns`, and `ToastBread` to return `Task<Egg>`, `Task<HashBrown>`, and `Task<Toast>` objects, respectively. The updated method names include the "Async" suffix: `FryEggsAsync`, `FryHashBrownsAsync`, and `ToastBreadAsync`. The `Main` method returns the `Task` object, although it doesn't have a `return` expression, which is by design. For more information, see Evaluation of a void-returning async function.

Note

The updated code doesn't yet take advantage of key features of asynchronous programming, which can result in shorter completion times. The code processes the tasks in roughly the same amount of time as the initial synchronous version. For the full method implementations, see the final version of the code later in this article.

Let's apply the breakfast example to the updated code. The thread doesn't block while the eggs or hash browns are cooking, but the code also doesn't start other tasks until the current work completes. You still put the bread in the toaster and stare at the toaster until the bread pops up, but you can now respond to interruptions. In a restaurant where multiple orders are placed, the cook can start a new order while another is already cooking.

In the updated code, the thread working on the breakfast isn't blocked while waiting for any started task that's unfinished. For some applications, this change is all you need. You can enable your app to support user interaction while data downloads from the web. In other scenarios, you might want to start other tasks while waiting for the previous task to complete.

## Start tasks concurrently

For most operations, you want to start several independent tasks immediately. As each task completes, you initiate other work that's ready to start. When you apply this methodology to the breakfast example, you can prepare breakfast more quickly. You also get everything ready close to the same time, so you can enjoy a hot breakfast.

The System.Threading.Tasks.Task class and related types are classes you can use to apply this style of reasoning to tasks that are in progress. This approach enables you to write code that more closely resembles the way you create breakfast in real life. You start cooking the eggs, hash browns, and toast at the same time. As each food item requires action, you turn your attention to that task, take care of the action, and then wait for something else that requires your attention.

In your code, you start a task and hold on to the Task object that represents the work. You use the `await` method on the task to delay acting on the work until the result is ready.

Apply these changes to the breakfast code. The first step is to store the tasks for operations when they start, rather than using the `await` expression:

```
Coffee cup = PourCoffee();
Console.WriteLine("Coffee is ready");
Task<Egg> eggsTask = FryEggsAsync(2);
Egg eggs = await eggsTask;
Console.WriteLine("Eggs are ready");
Task<HashBrown> hashBrownTask = FryHashBrownsAsync(3);
HashBrown hashBrown = await hashBrownTask;
Console.WriteLine("Hash browns are ready");
Task<Toast> toastTask = ToastBreadAsync(2);
Toast toast = await toastTask;
ApplyButter(toast);
ApplyJam(toast);
Console.WriteLine("Toast is ready");
Juice oj = PourOJ();
Console.WriteLine("Oj is ready");
Console.WriteLine("Breakfast is ready!");
```
These revisions don't help to get your breakfast ready any faster. The `await` expression is applied to all tasks as soon as they start. The next step is to move the `await` expressions for the hash browns and eggs to the end of the method, before you serve the breakfast:

```
Coffee cup = PourCoffee();
Console.WriteLine("Coffee is ready");
Task<Egg> eggsTask = FryEggsAsync(2);
Task<HashBrown> hashBrownTask = FryHashBrownsAsync(3);
Task<Toast> toastTask = ToastBreadAsync(2);
Toast toast = await toastTask;
ApplyButter(toast);
ApplyJam(toast);
Console.WriteLine("Toast is ready");
Juice oj = PourOJ();
Console.WriteLine("Oj is ready");
Egg eggs = await eggsTask;
Console.WriteLine("Eggs are ready");
HashBrown hashBrown = await hashBrownTask;
Console.WriteLine("Hash browns are ready");
Console.WriteLine("Breakfast is ready!");
```
You now have an asynchronously prepared breakfast that takes about 20 minutes to prepare. The total cook time is reduced because some tasks run concurrently.

The code updates improve the preparation process by reducing the cook time, but they introduce a regression by burning the eggs and hash browns. You start all the asynchronous tasks at once. You wait on each task only when you need the results. The code might be similar to program in a web application that makes requests to different microservices and then combines the results into a single page. You make all the requests immediately, and then apply the `await` expression on all those tasks and compose the web page.

## Support composition with tasks

The previous code revisions help get everything ready for breakfast at the same time, except the toast. The process of making the toast is a *composition* of an asynchronous operation (toast the bread) with synchronous operations (spread butter and jam on the toast). This example illustrates an important concept about asynchronous programming:

Important

The composition of an asynchronous operation followed by synchronous work is an asynchronous operation. Stated another way, if any portion of an operation is asynchronous, the entire operation is asynchronous.

In the previous updates, you learned how to use Task or Task<TResult> objects to hold running tasks. You wait on each task before you use its result. The next step is to create methods that represent the combination of other work. Before you serve breakfast, you want to wait on the task that represents toasting the bread before you spread the butter and jam.

You can represent this work with the following code:

```
static async Task<Toast> MakeToastWithButterAndJamAsync(int number)
{
 var toast = await ToastBreadAsync(number);
 ApplyButter(toast);
 ApplyJam(toast);
 return toast;
}
```
The `MakeToastWithButterAndJamAsync` method has the `async` modifier in its signature that signals to the compiler that the method contains an `await` expression and contains asynchronous operations. The method represents the task that toasts the bread, then spreads the butter and jam. The method returns a Task<TResult> object that represents the composition of the three operations.

The revised main block of code now looks like this:

```
static async Task Main(string[] args)
{
 Coffee cup = PourCoffee();
 Console.WriteLine("coffee is ready");
 var eggsTask = FryEggsAsync(2);
 var hashBrownTask = FryHashBrownsAsync(3);
 var toastTask = MakeToastWithButterAndJamAsync(2);
 var eggs = await eggsTask;
 Console.WriteLine("eggs are ready");
 var hashBrown = await hashBrownTask;
 Console.WriteLine("hash browns are ready");
 var toast = await toastTask;
 Console.WriteLine("toast is ready");
 Juice oj = PourOJ();
 Console.WriteLine("oj is ready");
 Console.WriteLine("Breakfast is ready!");
}
```
This code change illustrates an important technique for working with asynchronous code. You compose tasks by separating the operations into a new method that returns a task. You can choose when to wait on that task. You can start other tasks concurrently.

## Handle asynchronous exceptions

Up to this point, your code implicitly assumes all tasks complete successfully. Asynchronous methods throw exceptions, just like their synchronous counterparts. The goals for asynchronous support for exceptions and error handling are the same as for asynchronous support in general. The best practice is to write code that reads like a series of synchronous statements. Tasks throw exceptions when they can't complete successfully. The client code can catch those exceptions when the `await` expression is applied to a started task.

In the breakfast example, suppose the toaster catches fire while toasting the bread. You can simulate that problem by modifying the `ToastBreadAsync` method to match the following code:

```
private static async Task<Toast> ToastBreadAsync(int slices)
{
 for (int slice = 0; slice < slices; slice++)
 {
 Console.WriteLine("Putting a slice of bread in the toaster");
 }
 Console.WriteLine("Start toasting...");
 await Task.Delay(2000);
 Console.WriteLine("Fire! Toast is ruined!");
 throw new InvalidOperationException("The toaster is on fire");
 await Task.Delay(1000);
 Console.WriteLine("Remove toast from toaster");
 return new Toast();
}
```
Note

When you compile this code, you see a warning about unreachable code. This error is by design. After the toaster catches fire, operations don't proceed normally and the code returns an error.

After you make the code changes, run the application and check the output:

```
Pouring coffee
Coffee is ready
Warming the egg pan...
putting 3 hash brown patties in the pan
Cooking first side of hash browns...
Putting a slice of bread in the toaster
Putting a slice of bread in the toaster
Start toasting...
Fire! Toast is ruined!
Flipping a hash brown patty
Flipping a hash brown patty
Flipping a hash brown patty
Cooking the second side of hash browns...
Cracking 2 eggs
Cooking the eggs ...
Put hash browns on plate
Put eggs on plate
Eggs are ready
Hash browns are ready
Unhandled exception. System.InvalidOperationException: The toaster is on fire
 at AsyncBreakfast.Program.ToastBreadAsync(Int32 slices) in Program.cs:line 65
 at AsyncBreakfast.Program.MakeToastWithButterAndJamAsync(Int32 number) in Program.cs:line 36
 at AsyncBreakfast.Program.Main(String[] args) in Program.cs:line 24
 at AsyncBreakfast.Program.<Main>(String[] args)
```
Notice that quite a few tasks finish between the time when the toaster catches fire and the system observes the exception. When a task that runs asynchronously throws an exception, that task is **faulted**. The `Task` object holds the exception thrown in the Task.Exception property. Faulted tasks throw an exception when the `await` expression is applied to the task.

There are two important mechanisms to understand about this process:

- How an exception is stored in a faulted task
- How an exception is unpackaged and rethrown when code waits (`await`) on a faulted task

When code running asynchronously throws an exception, the exception is stored in the `Task` object. The Task.Exception property is a System.AggregateException object because more than one exception might be thrown during asynchronous work. Any exception thrown is added to the AggregateException.InnerExceptions collection. If the `Exception` property is null, a new `AggregateException` object is created and the thrown exception is the first item in the collection.

The most common scenario for a faulted task is that the `Exception` property contains exactly one exception. When your code waits on a faulted task, it rethrows the first AggregateException.InnerExceptions exception in the collection. This result is the reason why the output from the example shows an System.InvalidOperationException object rather than an `AggregateException` object. Extracting the first inner exception makes working with asynchronous methods as similar as possible to working with their synchronous counterparts. You can examine the `Exception` property in your code when your scenario might generate multiple exceptions.

Tip

The recommended practice is for any argument validation exceptions to emerge *synchronously* from task-returning methods. For more information and examples, see Exceptions in task-returning methods.

Before you continue to the next section, comment out the following two statements in your `ToastBreadAsync` method. You don't want to start another fire:

```
Console.WriteLine("Fire! Toast is ruined!");
throw new InvalidOperationException("The toaster is on fire");
```
## Apply await expressions to tasks efficiently

You can improve the series of `await` expressions at the end of the previous code by using methods of the `Task` class. One API is the WhenAll method, which returns a Task object that completes when all the tasks in its argument list are complete. The following code demonstrates this method:

```
await Task.WhenAll(eggsTask, hashBrownTask, toastTask);
Console.WriteLine("Eggs are ready");
Console.WriteLine("Hash browns are ready");
Console.WriteLine("Toast is ready");
Console.WriteLine("Breakfast is ready!");
```
Another option is to use the WhenAny method, which returns a `Task<Task>` object that completes when any of its arguments complete. You can wait on the returned task because you know the task is done. The following code shows how you can use the WhenAny method to wait on the first task to finish and then process its result. After you process the result from the completed task, you remove the completed task from the list of tasks passed to the `WhenAny` method.

```
var breakfastTasks = new List<Task> { eggsTask, hashBrownTask, toastTask };
while (breakfastTasks.Count > 0)
{
 Task finishedTask = await Task.WhenAny(breakfastTasks);
 if (finishedTask == eggsTask)
 {
 Console.WriteLine("Eggs are ready");
 }
 else if (finishedTask == hashBrownTask)
 {
 Console.WriteLine("Hash browns are ready");
 }
 else if (finishedTask == toastTask)
 {
 Console.WriteLine("Toast is ready");
 }
 await finishedTask;
 breakfastTasks.Remove(finishedTask);
}
```
Near the end of the code snippet, notice the `await finishedTask;` expression. This line is important because `Task.WhenAny` returns a `Task<Task>` - a wrapper task that contains the completed task. When you `await Task.WhenAny`, you're waiting for the wrapper task to complete, and the result is the actual task that finished first. However, to retrieve that task's result or ensure any exceptions are properly thrown, you must `await` the completed task itself (stored in `finishedTask`). Even though you know the task has finished, awaiting it again allows you to access its result or handle any exceptions that might have caused it to fault.

### Review final code

Here's what the final version of the code looks like:

```
using System;
using System.Collections.Generic;
using System.Threading.Tasks;
namespace AsyncBreakfast
{
 // These classes are intentionally empty for the purpose of this example. They are simply marker classes for the purpose of demonstration, contain no properties, and serve no other purpose.
 internal class HashBrown { }
 internal class Coffee { }
 internal class Egg { }
 internal class Juice { }
 internal class Toast { }
 class Program
 {
 static async Task Main(string[] args)
 {
 Coffee cup = PourCoffee();
 Console.WriteLine("coffee is ready");
 var eggsTask = FryEggsAsync(2);
 var hashBrownTask = FryHashBrownsAsync(3);
 var toastTask = MakeToastWithButterAndJamAsync(2);
 var breakfastTasks = new List<Task> { eggsTask, hashBrownTask, toastTask };
 while (breakfastTasks.Count > 0)
 {
 Task finishedTask = await Task.WhenAny(breakfastTasks);
 if (finishedTask == eggsTask)
 {
 Console.WriteLine("eggs are ready");
 }
 else if (finishedTask == hashBrownTask)
 {
 Console.WriteLine("hash browns are ready");
 }
 else if (finishedTask == toastTask)
 {
 Console.WriteLine("toast is ready");
 }
 await finishedTask;
 breakfastTasks.Remove(finishedTask);
 }
 Juice oj = PourOJ();
 Console.WriteLine("oj is ready");
 Console.WriteLine("Breakfast is ready!");
 }
 static async Task<Toast> MakeToastWithButterAndJamAsync(int number)
 {
 var toast = await ToastBreadAsync(number);
 ApplyButter(toast);
 ApplyJam(toast);
 return toast;
 }
 private static Juice PourOJ()
 {
 Console.WriteLine("Pouring orange juice");
 return new Juice();
 }
 private static void ApplyJam(Toast toast) =>
 Console.WriteLine("Putting jam on the toast");
 private static void ApplyButter(Toast toast) =>
 Console.WriteLine("Putting butter on the toast");
 private static async Task<Toast> ToastBreadAsync(int slices)
 {
 for (int slice = 0; slice < slices; slice++)
 {
 Console.WriteLine("Putting a slice of bread in the toaster");
 }
 Console.WriteLine("Start toasting...");
 await Task.Delay(3000);
 Console.WriteLine("Remove toast from toaster");
 return new Toast();
 }
 private static async Task<HashBrown> FryHashBrownsAsync(int patties)
 {
 Console.WriteLine($"putting {patties} hash brown patties in the pan");
 Console.WriteLine("cooking first side of hash browns...");
 await Task.Delay(3000);
 for (int patty = 0; patty < patties; patty++)
 {
 Console.WriteLine("flipping a hash brown patty");
 }
 Console.WriteLine("cooking the second side of hash browns...");
 await Task.Delay(3000);
 Console.WriteLine("Put hash browns on plate");
 return new HashBrown();
 }
 private static async Task<Egg> FryEggsAsync(int howMany)
 {
 Console.WriteLine("Warming the egg pan...");
 await Task.Delay(3000);
 Console.WriteLine($"cracking {howMany} eggs");
 Console.WriteLine("cooking the eggs ...");
 await Task.Delay(3000);
 Console.WriteLine("Put eggs on plate");
 return new Egg();
 }
 private static Coffee PourCoffee()
 {
 Console.WriteLine("Pouring coffee");
 return new Coffee();
 }
 }
}
```
The code completes the asynchronous breakfast tasks in about 15 minutes. The total time is reduced because some tasks run concurrently. The code simultaneously monitors multiple tasks and takes action only as needed.

The final code is asynchronous. It more accurately reflects how a person might cook breakfast. Compare the final code with the first code sample in the article. The core actions are still clear by reading the code. You can read the final code the same way you read the list of instructions for making a breakfast, as shown at the beginning of the article. The language features for the `async` and `await` keywords provide the translation every person makes to follow the written instructions: Start tasks as you can and don't block while waiting for tasks to complete.

## Async/await vs ContinueWith

The `async` and `await` keywords provide syntactic simplification over using Task.ContinueWith directly. While `async`/`await` and `ContinueWith` have similar semantics for handling asynchronous operations, the compiler doesn't necessarily translate `await` expressions directly into `ContinueWith` method calls. Instead, the compiler generates optimized state machine code that provides the same logical behavior. This transformation provides significant readability and maintainability benefits, especially when chaining multiple asynchronous operations.

Consider a scenario where you need to perform multiple sequential asynchronous operations. Here's how the same logic looks when implemented with `ContinueWith` compared to `async`/`await`:

### Using ContinueWith

With `ContinueWith`, each step in a sequence of asynchronous operations requires nested continuations:

```
// Using ContinueWith - demonstrates the complexity when chaining operations
static Task MakeBreakfastWithContinueWith()
{
 return StartCookingEggsAsync()
 .ContinueWith(eggsTask =>
 {
 var eggs = eggsTask.Result;
 Console.WriteLine("Eggs ready, starting bacon...");
 return StartCookingBaconAsync();
 })
 .Unwrap()
 .ContinueWith(baconTask =>
 {
 var bacon = baconTask.Result;
 Console.WriteLine("Bacon ready, starting toast...");
 return StartToastingBreadAsync();
 })
 .Unwrap()
 .ContinueWith(toastTask =>
 {
 var toast = toastTask.Result;
 Console.WriteLine("Toast ready, applying butter...");
 return ApplyButterAsync(toast);
 })
 .Unwrap()
 .ContinueWith(butteredToastTask =>
 {
 var butteredToast = butteredToastTask.Result;
 Console.WriteLine("Butter applied, applying jam...");
 return ApplyJamAsync(butteredToast);
 })
 .Unwrap()
 .ContinueWith(finalToastTask =>
 {
 var finalToast = finalToastTask.Result;
 Console.WriteLine("Breakfast completed with ContinueWith!");
 });
}
```
### Using async/await

The same sequence of operations using `async`/`await` reads much more naturally:

```
// Using async/await - much cleaner and easier to read
static async Task MakeBreakfastWithAsyncAwait()
{
 var eggs = await StartCookingEggsAsync();
 Console.WriteLine("Eggs ready, starting bacon...");

 var bacon = await StartCookingBaconAsync();
 Console.WriteLine("Bacon ready, starting toast...");

 var toast = await StartToastingBreadAsync();
 Console.WriteLine("Toast ready, applying butter...");

 var butteredToast = await ApplyButterAsync(toast);
 Console.WriteLine("Butter applied, applying jam...");

 var finalToast = await ApplyJamAsync(butteredToast);
 Console.WriteLine("Breakfast completed with async/await!");
}
```
### Why async/await is preferred

The `async`/`await` approach offers several advantages:

- **Readability**: The code reads like synchronous code, making it easier to understand the flow of operations.
- **Maintainability**: Adding or removing steps in the sequence requires minimal code changes.
- **Error handling**: Exception handling with- `try`/- `catch`blocks works naturally, whereas- `ContinueWith`requires careful handling of faulted tasks.
- **Debugging**: The call stack and debugger experience is much better with- `async`/- `await`.
- **Performance**: The compiler optimizations for- `async`/- `await`are more sophisticated than manual- `ContinueWith`chains.

The benefit becomes even more apparent as the number of chained operations increases. While a single continuation might be manageable with `ContinueWith`, sequences of 3-4 or more asynchronous operations quickly become difficult to read and maintain. This pattern, known as "monadic do-notation" in functional programming, allows you to compose multiple asynchronous operations in a sequential, readable manner.

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

You can avoid performance bottlenecks and enhance the overall responsiveness of your application by using asynchronous programming. However, traditional techniques for writing asynchronous applications can be complicated, making them difficult to write, debug, and maintain.

C# supports simplified approach, async programming, that uses asynchronous support in the .NET runtime. The compiler does the difficult work that the developer used to do, and your application retains a logical structure that resembles synchronous code. As a result, you get all the advantages of asynchronous programming with a fraction of the effort.

This article provides an overview of when and how to use async programming and includes links to other articles that contain details and examples.

## Async improves responsiveness

Asynchrony is essential for activities that are potentially blocking, such as web access. Access to a web resource sometimes is slow or delayed. If such an activity is blocked in a synchronous process, the entire application must wait. In an asynchronous process, the application can continue with other work that doesn't depend on the web resource until the potentially blocking task finishes.

The following table shows typical areas where asynchronous programming improves responsiveness. The listed APIs from .NET and the Windows Runtime contain methods that support async programming.

| Application area | .NET types with async methods | Windows Runtime types with async methods |
|---|---|---|
| Web access | HttpClient | Windows.Web.Http.HttpClient SyndicationClient |
| Working with files | JsonSerializer StreamReader StreamWriter XmlReader XmlWriter | StorageFile |
| Working with images | MediaCapture BitmapEncoder BitmapDecoder | |
| WCF programming | Synchronous and Asynchronous Operations |

Asynchrony proves especially valuable for applications that access the UI thread because all UI-related activity usually shares one thread. If any process is blocked in a synchronous application, all are blocked. Your application stops responding, and you might conclude that it failed when instead it's just waiting.

When you use asynchronous methods, the application continues to respond to the UI. You can resize or minimize a window, for example, or you can close the application if you don't want to wait for it to finish.

The async-based approach adds the equivalent of an automatic transmission to the list of options that you can choose from when designing asynchronous operations. That is, you get all the benefits of traditional asynchronous programming but with much less effort from the developer.

## Async methods are easy to write

The async and await keywords in C# are the heart of async programming. By using those two keywords, you can use resources in .NET Framework, .NET Core, or the Windows Runtime to create an asynchronous method almost as easily as you create a synchronous method. Asynchronous methods that you define by using the `async` keyword are referred to as *async methods*.

The following example shows an async method. Almost everything in the code should look familiar to you.

You can find a complete Windows Presentation Foundation (WPF) example available for download from Asynchronous programming with async and await in C#.

```
public async Task<int> GetUrlContentLengthAsync()
{
 using var client = new HttpClient();
 Task<string> getStringTask =
 client.GetStringAsync("https://learn.microsoft.com/dotnet");
 DoIndependentWork();
 string contents = await getStringTask;
 return contents.Length;
}
void DoIndependentWork()
{
 Console.WriteLine("Working...");
}
```
You can learn several practices from the preceding sample. Start with the method signature. It includes the `async` modifier. The return type is `Task<int>` (See "Return Types" section for more options). The method name ends in `Async`. In the body of the method, `GetStringAsync` returns a `Task<string>`. That means that when you `await` the task you get a `string` (`contents`). Before awaiting the task, you can do work that doesn't rely on the `string` from `GetStringAsync`.

Pay close attention to the `await` operator. It suspends `GetUrlContentLengthAsync`:

- `GetUrlContentLengthAsync`can't continue until- `getStringTask`is complete.
- Meanwhile, control returns to the caller of `GetUrlContentLengthAsync`.
- Control resumes here when `getStringTask`is complete.
- The `await`operator then retrieves the`string`result from`getStringTask`.

The return statement specifies an integer result. Any methods that are awaiting `GetUrlContentLengthAsync` retrieve the length value.

If `GetUrlContentLengthAsync` doesn't have any work that it can do between calling `GetStringAsync` and awaiting its completion, you can simplify your code by calling and awaiting in the following single statement.

```
string contents = await client.GetStringAsync("https://learn.microsoft.com/dotnet");
```
The following characteristics summarize what makes the previous example an async method:

- The method signature includes an - `async`modifier.
- The name of an async method, by convention, ends with an "Async" suffix.
- The return type is one of the following types: - Task<TResult> if your method has a return statement in which the operand has type `TResult`.
- Task if your method has no return statement or has a return statement with no operand.
- `void`if you're writing an async event handler.
- Any other type that has a `GetAwaiter`method.
 - For more information, see the Return types and parameters section.
- Task<TResult> if your method has a return statement in which the operand has type
- The method usually includes at least one - `await`expression, which marks a point where the method can't continue until the awaited asynchronous operation is complete. In the meantime, the method is suspended, and control returns to the method's caller. The next section of this article illustrates what happens at the suspension point.

In async methods, you use the provided keywords and types to indicate what you want to do, and the compiler does the rest, including keeping track of what must happen when control returns to an await point in a suspended method. Some routine processes, such as loops and exception handling, can be difficult to handle in traditional asynchronous code. In an async method, you write these elements much as you would in a synchronous solution, and the problem is solved.

For more information about asynchrony in previous versions of .NET Framework, see TPL and traditional .NET Framework asynchronous programming.

## What happens in an async method

The most important thing to understand in asynchronous programming is how the control flow moves from method to method. The following diagram leads you through the process:

The numbers in the diagram correspond to the following steps, initiated when a calling method calls the async method.

- A calling method calls and awaits the - `GetUrlContentLengthAsync`async method.
- `GetUrlContentLengthAsync`creates an HttpClient instance and calls the GetStringAsync asynchronous method to download the contents of a website as a string.
- Something happens in - `GetStringAsync`that suspends its progress. Perhaps it must wait for a website to download or some other blocking activity. To avoid blocking resources,- `GetStringAsync`yields control to its caller,- `GetUrlContentLengthAsync`.- `GetStringAsync`returns a Task<TResult>, where- `TResult`is a string, and- `GetUrlContentLengthAsync`assigns the task to the- `getStringTask`variable. The task represents the ongoing process for the call to- `GetStringAsync`, with a commitment to produce an actual string value when the work is complete.
- Because - `getStringTask`isn't awaited yet,- `GetUrlContentLengthAsync`can continue with other work that doesn't depend on the final result from- `GetStringAsync`. That work is represented by a call to the synchronous method- `DoIndependentWork`.
- `DoIndependentWork`is a synchronous method that does its work and returns to its caller.
- `GetUrlContentLengthAsync`runs out of work that it can do without a result from- `getStringTask`.- `GetUrlContentLengthAsync`next wants to calculate and return the length of the downloaded string, but the method can't calculate that value until the method has the string.- Therefore, - `GetUrlContentLengthAsync`uses an await operator to suspend its progress and to yield control to the method that called- `GetUrlContentLengthAsync`.- `GetUrlContentLengthAsync`returns a- `Task<int>`to the caller. The task represents a promise to produce an integer result that's the length of the downloaded string.- Note - If - `GetStringAsync`(and therefore- `getStringTask`) completes before- `GetUrlContentLengthAsync`awaits it, control remains in- `GetUrlContentLengthAsync`. The expense of suspending and then returning to- `GetUrlContentLengthAsync`would be wasted if the called asynchronous process- `getStringTask`is complete and- `GetUrlContentLengthAsync`doesn't have to wait for the final result.- Inside the calling method the processing pattern continues. The caller might do other work that doesn't depend on the result from - `GetUrlContentLengthAsync`before awaiting that result, or the caller might await immediately. The calling method is waiting for- `GetUrlContentLengthAsync`, and- `GetUrlContentLengthAsync`is waiting for- `GetStringAsync`.
- `GetStringAsync`completes and produces a string result. The string result isn't returned by the call to- `GetStringAsync`in the way that you might expect. (Remember that the method already returned a task in step 3.) Instead, the string result is stored in the task that represents the completion of the method,- `getStringTask`. The await operator retrieves the result from- `getStringTask`. The assignment statement assigns the retrieved result to- `contents`.
- When - `GetUrlContentLengthAsync`has the string result, the method can calculate the length of the string. Then the work of- `GetUrlContentLengthAsync`is also complete, and the waiting event handler can resume. In the full example at the end of the article, you can confirm that the event handler retrieves and prints the value of the length result. If you're new to asynchronous programming, take a minute to consider the difference between synchronous and asynchronous behavior. A synchronous method returns when its work is complete (step 5), but an async method returns a task value when its work is suspended (steps 3 and 6). When the async method eventually completes its work, the task is marked as completed and the result, if any, is stored in the task.

## API async methods

You might be wondering where to find methods such as `GetStringAsync` that support async programming. .NET Framework 4.5 or higher and .NET Core contain many members that work with `async` and `await`. You can recognize them by the "Async" suffix appended to the member name, and by their return type of Task or Task<TResult>. For example, the `System.IO.Stream` class contains methods such as CopyToAsync, ReadAsync, and WriteAsync alongside the synchronous methods CopyTo, Read, and Write.

The Windows Runtime also contains many methods that you can use with `async` and `await` in Windows apps. For more information, see Threading and async programming for UWP development, and Asynchronous programming (Windows Store apps) and Quickstart: Calling asynchronous APIs in C# or Visual Basic if you use earlier versions of the Windows Runtime.

## Threads

Async methods are intended to be non-blocking operations. An `await` expression in an async method doesn't block the current thread while the awaited task is running. Instead, the expression signs up the rest of the method as a continuation and returns control to the caller of the async method.

The `async` and `await` keywords don't cause extra threads to be created. Async methods don't require multithreading because an async method doesn't run on its own thread. The method runs on the current synchronization context and uses time on the thread only when the method is active. You can use Task.Run to move CPU-bound work to a background thread, but a background thread doesn't help with a process that's just waiting for results to become available.

The async-based approach to asynchronous programming is preferable to existing approaches in almost every case. In particular, this approach is better than the BackgroundWorker class for I/O-bound operations because the code is simpler and you don't have to guard against race conditions. In combination with the Task.Run method, async programming is better than BackgroundWorker for CPU-bound operations because async programming separates the coordination details of running your code from the work that `Task.Run` transfers to the thread pool.

## Async and await

If you specify that a method is an async method by using the async modifier, you enable the following two capabilities.

- The marked async method can use await to designate suspension points. The - `await`operator tells the compiler that the async method can't continue past that point until the awaited asynchronous process is complete. In the meantime, control returns to the caller of the async method.- The suspension of an async method at an - `await`expression doesn't constitute an exit from the method, and- `finally`blocks don't run.
- The marked async method can itself be awaited by methods that call it.

An async method typically contains one or more occurrences of an `await` operator, but the absence of `await` expressions doesn't cause a compiler error. If an async method doesn't use an `await` operator to mark a suspension point, the method executes as a synchronous method does, despite the `async` modifier. The compiler issues a warning for such methods.

`async` and `await` are contextual keywords. For more information and examples, see the following articles:

## Return types and parameters

An async method typically returns a Task or a Task<TResult>. Inside an async method, an `await` operator is applied to a task that's returned from a call to another async method.

You specify Task<TResult> as the return type if the method contains a `return` statement that specifies an operand of type `TResult`.

You use Task as the return type if the method has no return statement or has a return statement that doesn't return an operand.

You can also specify any other return type, if the type includes a `GetAwaiter` method. ValueTask<TResult> is an example of such a type. It's available in the System.Threading.Tasks.Extension NuGet package.

The following example shows how you declare and call a method that returns a Task<TResult> or a Task:

```
async Task<int> GetTaskOfTResultAsync()
{
 int hours = 0;
 await Task.Delay(0);
 return hours;
}
Task<int> returnedTaskTResult = GetTaskOfTResultAsync();
int intResult = await returnedTaskTResult;
// Single line
// int intResult = await GetTaskOfTResultAsync();
async Task GetTaskAsync()
{
 await Task.Delay(0);
 // No return statement needed
}
Task returnedTask = GetTaskAsync();
await returnedTask;
// Single line
await GetTaskAsync();
```
Each returned task represents ongoing work. A task encapsulates information about the state of the asynchronous process and, eventually, either the final result from the process or the exception that the process raises if it doesn't succeed.

An async method can also have a `void` return type. This return type is used primarily to define event handlers, where a `void` return type is required. Async event handlers often serve as the starting point for async programs.

An async method that has a `void` return type can't be awaited, and the caller of a void-returning method can't catch any exceptions that the method throws.

An async method can't declare in, ref or out parameters, but the method can call methods that have such parameters. Similarly, an async method can't return a value by reference, although it can call methods with ref return values.

For more information and examples, see Async return types (C#).

Asynchronous APIs in Windows Runtime programming have one of the following return types, which are similar to tasks:

- IAsyncOperation<TResult>, which corresponds to Task<TResult>
- IAsyncAction, which corresponds to Task
- IAsyncActionWithProgress<TProgress>
- IAsyncOperationWithProgress<TResult,TProgress>

## Naming convention

By convention, methods that return commonly awaitable types (for example, `Task`, `Task<T>`, `ValueTask`, `ValueTask<T>`) should have names that end with "Async". Methods that start an asynchronous operation but don't return an awaitable type shouldn't have names that end with "Async", but might start with "Begin", "Start", or some other verb to suggest this method doesn't return or throw the result of the operation.

You can ignore the convention where an event, base class, or interface contract suggests a different name. For example, you shouldn't rename common event handlers, such as `OnButtonClick`.

## Related articles (Visual Studio)

| Title | Description |
|---|---|
| How to make multiple web requests in parallel by using async and await (C#) | Demonstrates how to start several tasks at the same time. |
| Async return types (C#) | Illustrates the types that async methods can return, and explains when each type is appropriate. |
| Cancel tasks with a cancellation token as a signaling mechanism. | Shows how to add the following functionality to your async solution: - Cancel a list of tasks (C#) - Cancel tasks after a period of time (C#) - Process asynchronous task as they complete (C#) |
| Using async for file access (C#) | Lists and demonstrates the benefits of using async and await to access files. |
| Task-based asynchronous pattern (TAP) | Describes an asynchronous pattern. The pattern is based on the Task and Task<TResult> types. |
| Async Videos on Channel 9 | Provides links to various videos about async programming. |

Note

Access to this page requires authorization. You can try signing in or changing directories.

Access to this page requires authorization. You can try changing directories.

If your code implements I/O-bound scenarios to support network data requests, database access, or file system read/writes, asynchronous programming is the best approach. You can also write asynchronous code for CPU-bound scenarios like expensive calculations.

C# has a language-level asynchronous programming model that allows you to easily write asynchronous code without having to juggle callbacks or conform to a library that supports asynchrony. The model follows what is known as the Task-based asynchronous pattern (TAP).

## Explore the asynchronous programming model

The `Task` and `Task<T>` objects represent the core of asynchronous programming. These objects are used to model asynchronous operations by supporting the `async` and `await` keywords. In most cases, the model is fairly simple for both I/O-bound and CPU-bound scenarios. Inside an `async` method:

- **I/O-bound code**starts an operation represented by a- `Task`or- `Task<T>`object within the- `async`method.
- **CPU-bound code**starts an operation on a background thread with the Task.Run method.

In both cases, an active `Task` represents an asynchronous operation that might not be complete.

The `await` keyword is where the magic happens. It yields control to the caller of the method that contains the `await` expression, and ultimately allows the UI to be responsive or a service to be elastic. While there are ways to approach asynchronous code other than by using the `async` and `await` expressions, this article focuses on the language-level constructs.

Note

Some examples presented in this article use the System.Net.Http.HttpClient class to download data from a web service. In the example code, the `s_httpClient` object is a static field of type `Program` class:

`private static readonly HttpClient s_httpClient = new();`

For more information, see the complete example code at the end of this article.

### Review underlying concepts

When you implement asynchronous programming in your C# code, the compiler transforms your program into a state machine. This construct tracks various operations and state in your code, such as yielding execution when the code reaches an `await` expression, and resuming execution when a background job completes.

In terms of computer science theory, asynchronous programming is an implementation of the Promise model of asynchrony.

In the asynchronous programming model, there are several key concepts to understand:

- You can use asynchronous code for both I/O-bound and CPU-bound code, but the implementation is different.
- Asynchronous code uses `Task<T>`and`Task`objects as constructs to model work running in the background.
- The `async`keyword declares a method as an asynchronous method, which allows you to use the`await`keyword in the method body.
- When you apply the `await`keyword, the code suspends the calling method and yields control back to its caller until the task completes.
- You can only use the `await`expression in an asynchronous method.

### I/O-bound example: Download data from web service

In this example, when the user selects a button, the app downloads data from a web service. You don't want to block the UI thread for the app during the download process. The following code accomplishes this task:

```
s_downloadButton.Clicked += async (o, e) =>
{
 // This line will yield control to the UI as the request
 // from the web service is happening.
 //
 // The UI thread is now free to perform other work.
 var stringData = await s_httpClient.GetStringAsync(URL);
 DoSomethingWithData(stringData);
};
```
The code expresses the intent (downloading data asynchronously) without getting bogged down in interacting with `Task` objects.

### CPU-bound example: Run game calculation

In the next example, a mobile game inflicts damage on several agents on the screen in response to a button event. Performing the damage calculation can be expensive. Running the calculation on the UI thread can cause display and UI interaction issues during the calculation.

The best way to handle the task is to start a background thread to complete the work with the `Task.Run` method. The operation yields by using an `await` expression. The operation resumes when the task completes. This approach allows the UI to run smoothly while the work completes in the background.

```
static DamageResult CalculateDamageDone()
{
 return new DamageResult()
 {
 // Code omitted:
 //
 // Does an expensive calculation and returns
 // the result of that calculation.
 };
}
s_calculateButton.Clicked += async (o, e) =>
{
 // This line will yield control to the UI while CalculateDamageDone()
 // performs its work. The UI thread is free to perform other work.
 var damageResult = await Task.Run(() => CalculateDamageDone());
 DisplayDamage(damageResult);
};
```
The code clearly expresses the intent of the button `Clicked` event. It doesn't require managing a background thread manually, and it completes the task in a nonblocking manner.

## Recognize CPU-bound and I/O-bound scenarios

The previous examples demonstrate how to use the `async` modifier and `await` expression for I/O-bound and CPU-bound work. An example for each scenario showcases how the code is different based on where the operation is bound. To prepare for your implementation, you need to understand how to identify when an operation is I/O-bound or CPU-bound. Your implementation choice can greatly affect the performance of your code and potentially lead to misusing constructs.

There are two primary questions to address before you write any code:

| Question | Scenario | Implementation |
|---|---|---|
| Should the code wait for a result or action, such as data from a database? | I/O-bound | Use the `async`modifier and`await`expressionwithoutthe`Task.Run`method.Avoid using the Task Parallel Library. |
| Should the code run an expensive computation? | CPU-bound | Use the `async`modifier and`await`expression, but spawn off the work on another thread with the`Task.Run`method. This approach addresses concerns with CPU responsiveness.If the work is appropriate for concurrency and parallelism, also consider using the Task Parallel Library. |

Always measure the execution of your code. You might discover that your CPU-bound work isn't costly enough compared with the overhead of context switches when multithreading. Every choice has tradeoffs. Pick the correct tradeoff for your situation.

## Explore other examples

The examples in this section demonstrate several ways you can write asynchronous code in C#. They cover a few scenarios you might encounter.

### Extract data from a network

The following code downloads HTML from a given URL and counts the number of times the string ".NET" occurs in the HTML. The code uses ASP.NET to define a Web API controller method, which performs the task and returns the count.

Note

If you plan on doing HTML parsing in production code, don't use regular expressions. Use a parsing library instead.

```
[HttpGet, Route("DotNetCount")]
static public async Task<int> GetDotNetCountAsync(string URL)
{
 // Suspends GetDotNetCountAsync() to allow the caller (the web server)
 // to accept another request, rather than blocking on this one.
 var html = await s_httpClient.GetStringAsync(URL);
 return Regex.Matches(html, @"\.NET").Count;
}
```
You can write similar code for a Universal Windows App and perform the counting task after a button press:

```
private readonly HttpClient _httpClient = new HttpClient();
private async void OnSeeTheDotNetsButtonClick(object sender, RoutedEventArgs e)
{
 // Capture the task handle here so we can await the background task later.
 var getDotNetFoundationHtmlTask = _httpClient.GetStringAsync("https://dotnetfoundation.org");
 // Any other work on the UI thread can be done here, such as enabling a Progress Bar.
 // It's important to do the extra work here before the "await" call,
 // so the user sees the progress bar before execution of this method is yielded.
 NetworkProgressBar.IsEnabled = true;
 NetworkProgressBar.Visibility = Visibility.Visible;
 // The await operator suspends OnSeeTheDotNetsButtonClick(), returning control to its caller.
 // This action is what allows the app to be responsive and not block the UI thread.
 var html = await getDotNetFoundationHtmlTask;
 int count = Regex.Matches(html, @"\.NET").Count;
 DotNetCountLabel.Text = $"Number of .NETs on dotnetfoundation.org: {count}";
 NetworkProgressBar.IsEnabled = false;
 NetworkProgressBar.Visibility = Visibility.Collapsed;
}
```
### Wait for multiple tasks to complete

In some scenarios, the code needs to retrieve multiple pieces of data concurrently. The `Task` APIs provide methods that enable you to write asynchronous code that performs a nonblocking wait on multiple background jobs:

- Task.WhenAll method
- Task.WhenAny method

The following example shows how you might grab `User` object data for a set of `userId` objects.

```
private static async Task<User> GetUserAsync(int userId)
{
 // Code omitted:
 //
 // Given a user Id {userId}, retrieves a User object corresponding
 // to the entry in the database with {userId} as its Id.
 return await Task.FromResult(new User() { id = userId });
}
private static async Task<IEnumerable<User>> GetUsersAsync(IEnumerable<int> userIds)
{
 var getUserTasks = new List<Task<User>>();
 foreach (int userId in userIds)
 {
 getUserTasks.Add(GetUserAsync(userId));
 }
 return await Task.WhenAll(getUserTasks);
}
```
You can write this code more succinctly by using LINQ:

```
private static async Task<User[]> GetUsersByLINQAsync(IEnumerable<int> userIds)
{
 var getUserTasks = userIds.Select(id => GetUserAsync(id)).ToArray();
 return await Task.WhenAll(getUserTasks);
}
```
Although you write less code by using LINQ, exercise caution when mixing LINQ with asynchronous code. LINQ uses deferred (or lazy) execution, which means that without immediate evaluation, async calls don't happen until the sequence is enumerated.

The previous example is correct and safe, because it uses the Enumerable.ToArray method to immediately evaluate the LINQ query and store the tasks in an array. This approach ensures the `id => GetUserAsync(id)` calls execute immediately and all tasks start concurrently, just like the `foreach` loop approach. Always use Enumerable.ToArray or Enumerable.ToList when creating tasks with LINQ to ensure immediate execution and concurrent task execution. Here's an example that demonstrates using `ToList()` with `Task.WhenAny` to process tasks as they complete:

```
private static async Task ProcessTasksAsTheyCompleteAsync(IEnumerable<int> userIds)
{
 var getUserTasks = userIds.Select(id => GetUserAsync(id)).ToList();

 while (getUserTasks.Count > 0)
 {
 Task<User> completedTask = await Task.WhenAny(getUserTasks);
 getUserTasks.Remove(completedTask);

 User user = await completedTask;
 Console.WriteLine($"Processed user {user.id}");
 }
}
```
In this example, `ToList()` creates a list that supports the `Remove()` operation, allowing you to dynamically remove completed tasks. This pattern is particularly useful when you want to handle results as soon as they're available, rather than waiting for all tasks to complete.

Although you write less code by using LINQ, exercise caution when mixing LINQ with asynchronous code. LINQ uses deferred (or lazy) execution. Asynchronous calls don't happen immediately as they do in a `foreach` loop, unless you force the generated sequence to iterate with a call to the `.ToList()` or `.ToArray()` method.

You can choose between Enumerable.ToArray and Enumerable.ToList based on your scenario:

- Use `ToArray()`when you plan to process all tasks together, such as with`Task.WhenAll`. Arrays are efficient for scenarios where the collection size is fixed.
- Use `ToList()`when you need to dynamically manage tasks, such as with`Task.WhenAny`where you might remove completed tasks from the collection as they finish.

## Review considerations for asynchronous programming

With asynchronous programming, there are several details to keep in mind that can prevent unexpected behavior.

### Use await inside async() method body

When you use the `async` modifier, you should include one or more `await` expressions in the method body. If the compiler doesn't encounter an `await` expression, the method fails to yield. Although the compiler generates a warning, the code still compiles and the compiler runs the method. The state machine generated by the C# compiler for the asynchronous method doesn't accomplish anything, so the entire process is highly inefficient.

### Add "Async" suffix to asynchronous method names

The .NET style convention is to add the "Async" suffix to all asynchronous method names. This approach helps to more easily differentiate between synchronous and asynchronous methods. Certain methods that aren't explicitly called by your code (such as event handlers or web controller methods) don't necessarily apply in this scenario. Because these items aren't explicitly called by your code, using explicit naming isn't as important.

### Return 'async void' only from event handlers

Event handlers must declare `void` return types and can't use or return `Task` and `Task<T>` objects as other methods do. When you write asynchronous event handlers, you need to use the `async` modifier on a `void` returning method for the handlers. Other implementations of `async void` returning methods don't follow the TAP model and can present challenges:

- Exceptions thrown in an `async void`method can't be caught outside of that method
- `async void`methods are difficult to test
- `async void`methods can cause negative side effects if the caller isn't expecting them to be asynchronous

### Use caution with asynchronous lambdas in LINQ

It's important to use caution when you implement asynchronous lambdas in LINQ expressions. Lambda expressions in LINQ use deferred execution, which means the code can execute at an unexpected time. The introduction of blocking tasks into this scenario can easily result in a deadlock, if the code isn't written correctly. Moreover, the nesting of asynchronous code can also make it difficult to reason about the execution of the code. Async and LINQ are powerful, but these techniques should be used together as carefully and clearly as possible.

### Yield for tasks in a nonblocking manner

If your program needs the result of a task, write code that implements the `await` expression in a nonblocking manner. Blocking the current thread as a means to wait synchronously for a `Task` item to complete can result in deadlocks and blocked context threads. This programming approach can require more complex error-handling. The following table provides guidance on how access results from tasks in a nonblocking way:

| Task scenario | Current code | Replace with 'await' |
|---|---|---|
| Retrieve the result of a background task | `Task.Wait`or`Task.Result` | `await` |
| Continue when any task completes | `Task.WaitAny` | `await Task.WhenAny` |
| Continue when alltasks complete | `Task.WaitAll` | `await Task.WhenAll` |
| Continue after some amount of time | `Thread.Sleep` | `await Task.Delay` |

### Consider using ValueTask type

When an asynchronous method returns a `Task` object, performance bottlenecks might be introduced in certain paths. Because `Task` is a reference type, a `Task` object is allocated from the heap. If a method declared with the `async` modifier returns a cached result or completes synchronously, the extra allocations can accrue significant time costs in performance critical sections of code. This scenario can become costly when the allocations occur in tight loops. For more information, see generalized async return types.

### Understand when to set ConfigureAwait(false)

Developers often inquire about when to use the Task.ConfigureAwait(Boolean) boolean. This API allows for a `Task` instance to configure the context for the state machine that implements any `await` expression. When the boolean isn't set correctly, performance can degrade or deadlocks can occur. For more information, see ConfigureAwait FAQ.

### Write less-stateful code

Avoid writing code that depends on the state of global objects or the execution of certain methods. Instead, depend only on the return values of methods. There are many benefits to writing code that is less-stateful:

- Easier to reason about code
- Easier to test code
- More simple to mix asynchronous and synchronous code
- Able to avoid race conditions in code
- Simple to coordinate asynchronous code that depends on return values
- (Bonus) Works well with dependency injection in code

A recommended goal is to achieve complete or near-complete Referential Transparency in your code. This approach results in a predictable, testable, and maintainable codebase.

### Synchronous access to asynchronous operations

In scenarios, you might need to block on asynchronous operations when the `await` keyword isn't available throughout your call stack. This situation occurs in legacy codebases or when integrating asynchronous methods into synchronous APIs that can't be changed.

Warning

Synchronous blocking on asynchronous operations can lead to deadlocks and should be avoided whenever possible. The preferred solution is to use `async`/`await` throughout your call stack.

When you must block synchronously on a `Task`, here are the available approaches, listed from most to least preferred:

#### Use GetAwaiter().GetResult()

The `GetAwaiter().GetResult()` pattern is generally the preferred approach when you must block synchronously:

```
// When you cannot use await
Task<string> task = GetDataAsync();
string result = task.GetAwaiter().GetResult();
```
This approach:

- Preserves the original exception without wrapping it in an `AggregateException`.
- Blocks the current thread until the task completes.
- Still carries deadlock risk if not used carefully.

#### Use Task.Run for complex scenarios

For complex scenarios where you need to isolate the asynchronous work:

```
// Offload to thread pool to avoid context deadlocks
string result = Task.Run(async () => await GetDataAsync()).GetAwaiter().GetResult();
```
This pattern:

- Executes the asynchronous method on a thread pool thread.
- Can help avoid some deadlock scenarios.
- Adds overhead by scheduling work to the thread pool.

#### Use Wait() and Result

You can use a blocking approach by calling Wait() and Result. However, this approach is discouraged because it wraps exceptions in AggregateException.

```
Task<string> task = GetDataAsync();
task.Wait();
string result = task.Result;
```
Problems with `Wait()` and `Result`:

- Exceptions are wrapped in `AggregateException`, making error handling more complex.
- Higher deadlock risk.
- Less clear intent in code.

#### Additional considerations

- Deadlock prevention: Be especially careful in UI applications or when using a synchronization context.
- Performance impact: Blocking threads reduces scalability.
- Exception handling: Test error scenarios carefully as exception behavior differs between patterns.

For more detailed guidance on the challenges and considerations of synchronous wrappers for asynchronous methods, see Should I expose synchronous wrappers for asynchronous methods?.

## Review the complete example

The following code represents the complete example, which is available in the *Program.cs* example file.

```
using System.Text.RegularExpressions;
using System.Windows;
using Microsoft.AspNetCore.Mvc;
class Button
{
 public Func<object, object, Task>? Clicked
 {
 get;
 internal set;
 }
}
class DamageResult
{
 public int Damage
 {
 get { return 0; }
 }
}
class User
{
 public bool isEnabled
 {
 get;
 set;
 }
 public int id
 {
 get;
 set;
 }
}
public class Program
{
 private static readonly Button s_downloadButton = new();
 private static readonly Button s_calculateButton = new();
 private static readonly HttpClient s_httpClient = new();
 private static readonly IEnumerable<string> s_urlList = new string[]
 {
 "https://learn.microsoft.com",
 "https://learn.microsoft.com/aspnet/core",
 "https://learn.microsoft.com/azure",
 "https://learn.microsoft.com/azure/devops",
 "https://learn.microsoft.com/dotnet",
 "https://learn.microsoft.com/dotnet/desktop/wpf/get-started/create-app-visual-studio",
 "https://learn.microsoft.com/education",
 "https://learn.microsoft.com/shows/net-core-101/what-is-net",
 "https://learn.microsoft.com/enterprise-mobility-security",
 "https://learn.microsoft.com/gaming",
 "https://learn.microsoft.com/graph",
 "https://learn.microsoft.com/microsoft-365",
 "https://learn.microsoft.com/office",
 "https://learn.microsoft.com/powershell",
 "https://learn.microsoft.com/sql",
 "https://learn.microsoft.com/surface",
 "https://dotnetfoundation.org",
 "https://learn.microsoft.com/visualstudio",
 "https://learn.microsoft.com/windows",
 "https://learn.microsoft.com/maui"
 };
 private static void Calculate()
 {
 static DamageResult CalculateDamageDone()
 {
 return new DamageResult()
 {
 // Code omitted:
 //
 // Does an expensive calculation and returns
 // the result of that calculation.
 };
 }
 s_calculateButton.Clicked += async (o, e) =>
 {
 // This line will yield control to the UI while CalculateDamageDone()
 // performs its work. The UI thread is free to perform other work.
 var damageResult = await Task.Run(() => CalculateDamageDone());
 DisplayDamage(damageResult);
 };
 }
 private static void DisplayDamage(DamageResult damage)
 {
 Console.WriteLine(damage.Damage);
 }
 private static void Download(string URL)
 {
 s_downloadButton.Clicked += async (o, e) =>
 {
 // This line will yield control to the UI as the request
 // from the web service is happening.
 //
 // The UI thread is now free to perform other work.
 var stringData = await s_httpClient.GetStringAsync(URL);
 DoSomethingWithData(stringData);
 };
 }
 private static void DoSomethingWithData(object stringData)
 {
 Console.WriteLine($"Displaying data: {stringData}");
 }
 private static async Task<User> GetUserAsync(int userId)
 {
 // Code omitted:
 //
 // Given a user Id {userId}, retrieves a User object corresponding
 // to the entry in the database with {userId} as its Id.
 return await Task.FromResult(new User() { id = userId });
 }
 private static async Task<IEnumerable<User>> GetUsersAsync(IEnumerable<int> userIds)
 {
 var getUserTasks = new List<Task<User>>();
 foreach (int userId in userIds)
 {
 getUserTasks.Add(GetUserAsync(userId));
 }
 return await Task.WhenAll(getUserTasks);
 }
 private static async Task<User[]> GetUsersByLINQAsync(IEnumerable<int> userIds)
 {
 var getUserTasks = userIds.Select(id => GetUserAsync(id)).ToArray();
 return await Task.WhenAll(getUserTasks);
 }
 private static async Task ProcessTasksAsTheyCompleteAsync(IEnumerable<int> userIds)
 {
 var getUserTasks = userIds.Select(id => GetUserAsync(id)).ToList();

 while (getUserTasks.Count > 0)
 {
 Task<User> completedTask = await Task.WhenAny(getUserTasks);
 getUserTasks.Remove(completedTask);

 User user = await completedTask;
 Console.WriteLine($"Processed user {user.id}");
 }
 }
 [HttpGet, Route("DotNetCount")]
 static public async Task<int> GetDotNetCountAsync(string URL)
 {
 // Suspends GetDotNetCountAsync() to allow the caller (the web server)
 // to accept another request, rather than blocking on this one.
 var html = await s_httpClient.GetStringAsync(URL);
 return Regex.Matches(html, @"\.NET").Count;
 }
 static async Task Main()
 {
 Console.WriteLine("Application started.");
 Console.WriteLine("Counting '.NET' phrase in websites...");
 int total = 0;
 foreach (string url in s_urlList)
 {
 var result = await GetDotNetCountAsync(url);
 Console.WriteLine($"{url}: {result}");
 total += result;
 }
 Console.WriteLine("Total: " + total);
 Console.WriteLine("Retrieving User objects with list of IDs...");
 IEnumerable<int> ids = new int[] { 1, 2, 3, 4, 5, 6, 7, 8, 9, 0 };
 var users = await GetUsersAsync(ids);
 foreach (User? user in users)
 {
 Console.WriteLine($"{user.id}: isEnabled={user.isEnabled}");
 }
 Console.WriteLine("Processing tasks as they complete...");
 await ProcessTasksAsTheyCompleteAsync(ids);
 Console.WriteLine("Application ending.");
 }
}
// Example output:
//
// Application started.
// Counting '.NET' phrase in websites...
// https://learn.microsoft.com: 0
// https://learn.microsoft.com/aspnet/core: 57
// https://learn.microsoft.com/azure: 1
// https://learn.microsoft.com/azure/devops: 2
// https://learn.microsoft.com/dotnet: 83
// https://learn.microsoft.com/dotnet/desktop/wpf/get-started/create-app-visual-studio: 31
// https://learn.microsoft.com/education: 0
// https://learn.microsoft.com/shows/net-core-101/what-is-net: 42
// https://learn.microsoft.com/enterprise-mobility-security: 0
// https://learn.microsoft.com/gaming: 0
// https://learn.microsoft.com/graph: 0
// https://learn.microsoft.com/microsoft-365: 0
// https://learn.microsoft.com/office: 0
// https://learn.microsoft.com/powershell: 0
// https://learn.microsoft.com/sql: 0
// https://learn.microsoft.com/surface: 0
// https://dotnetfoundation.org: 16
// https://learn.microsoft.com/visualstudio: 0
// https://learn.microsoft.com/windows: 0
// https://learn.microsoft.com/maui: 6
// Total: 238
// Retrieving User objects with list of IDs...
// 1: isEnabled= False
// 2: isEnabled= False
// 3: isEnabled= False
// 4: isEnabled= False
// 5: isEnabled= False
// 6: isEnabled= False
// 7: isEnabled= False
// 8: isEnabled= False
// 9: isEnabled= False
// 0: isEnabled= False
// Application ending.
```
