# Introduction

In the modern era, software is commonly delivered as a service: called *web apps*, or *software-as-a-service*. The twelve-factor app is a methodology for building software-as-a-service apps that:

- Use **declarative**formats for setup automation, to minimize time and cost for new developers joining the project;
- Have a **clean contract**with the underlying operating system, offering**maximum portability**between execution environments;
- Are suitable for **deployment**on modern**cloud platforms**, obviating the need for servers and systems administration;
- **Minimize divergence**between development and production, enabling- **continuous deployment**for maximum agility;
- And can **scale up**without significant changes to tooling, architecture, or development practices.

The twelve-factor methodology can be applied to apps written in any programming language, and which use any combination of backing services (database, queue, memory cache, etc).

# Background

The contributors to this document have been directly involved in the development and deployment of hundreds of apps, and indirectly witnessed the development, operation, and scaling of hundreds of thousands of apps via our work on the Heroku platform.

This document synthesizes all of our experience and observations on a wide variety of software-as-a-service apps in the wild. It is a triangulation on ideal practices for app development, paying particular attention to the dynamics of the organic growth of an app over time, the dynamics of collaboration between developers working on the app’s codebase, and avoiding the cost of software erosion.

Our motivation is to raise awareness of some systemic problems we’ve seen in modern application development, to provide a shared vocabulary for discussing those problems, and to offer a set of broad conceptual solutions to those problems with accompanying terminology. The format is inspired by Martin Fowler’s books *Patterns of Enterprise Application Architecture* and *Refactoring*.

## Why Intuit is Thrilled About the Evolution of the Twelve-Factor Model

###### 03 Apr, 2025

### Brett Weaver

 At Intuit, we’ve long embraced the twelve-factor app principles as a guiding framework for modern software development. As a company building cutting-edge development tools and runtime platforms for our internal engineers, these principles have been instrumental in unifying service developers, platform engineers, and SREs under a shared philosophy. **read on…**

## Evolving Twelve-Factor: Applications to Modern Cloud-Native Platforms

###### 10 Feb, 2025

### Brian Hammons

 The recent **open sourcing of the Twelve-Factor App Methodology** comes at a transformative moment for cloud-native platforms. As organizations increasingly rely on cloud-native technologies to power mission-critical workloads, the principles behind Twelve-Factor offer timeless foundations that remain relevant for modern platform builders. **read on…**

## December Monthly Updates

###### 3 Dec, 2024

### Vish Abrams

Welcome to our first monthly update! We’re excited to share our progress and what’s coming next.

#### What We’ve Been Working On

In addition to some minor formatting fixes, our initial focus has been on getting organized for larger updates. Here are the key activities: **read on…**

## Twelve-Factor App Methodology is now Open Source

###### 12 Nov, 2024

### Yehuda Katz

 **Join us in modernizing the twelve-factor app manifesto together.** As a community of app, framework and platform developers, we’re working together to refresh this foundational document for the modern era. While it’s not software we’re working on, we’ll use familiar processes like pull requests, issues, and reviews to collaborate together in the twelve-factor project repo.

This initiative builds on a strong foundation laid by Heroku when they originally created “The Twelve-Factor App” all the way back in 2011, a time when container-based deployment was still just emerging. Back then, developers could get apps running on their local machines, but common development mistakes often made it challenging to deploy those apps to production. **read on…**

## Narrow Conduits and the Application-Platform Interface

###### 12 Nov, 2024

### Vish Abrams

Welcome to the twelve-factor maintainters blog. As stated in our announcement, some of our posts will analyze the manifesto more generically. This is the first post in that vein, where we dive into the interface between the application and the platform.

It is well understood that defining a clear contract between parts of a system allows one to shed cognitive load on either side of the contract. This has been called the “narrow-waist” principle which has some unfortunate connotations, so we’ll refer to it as the “narrow-conduit” principle. This principle is especially valuable when the humans on either side of the conduit have dramatically different concerns. **read on…**

# Maintainers

## Vish Abrams

### Heroku/Salesforce

 Vish Abrams is Chief Architect at Heroku, a subsidiary of Salesforce. Formerly he helped Oracle create their cloud, where he focused on virtualization, containerization, and machine learning. He was also NASA Nebula Technical Lead during the creation of Nova, one of the founding OpenStack projects, and was a member of both the OpenStack board and technical committee.

## Evan Anderson

### Stacklok (independent)

 Evan worked for about 15 years in Google’s cloud, with about 10 years in the public cloud. During that time, he was a founding member of the Google Compute Engine team, then worked on App Engine (control plane API), Cloud Functions, and Knative/Cloud Run. In 2019, he moved to VMware, where he spent 4 years on VMware Tandy, VMware’s cloud-native developer platform. He’s also the author of “Designing Serverless Applications with Knative”, and has held many leadership roles in Knative over the years.

## Brian Hammons

### AWS

 Brian Hammons, a Principal Solutions Architect at AWS, is an original member of the launch team for Amazon EKS. He has held crucial roles in growing the service from its inception including co-founding projects such as eksworkshop.com, Data on EKS (DoEKS), and CNOE. Brian leads Application Modernization and Developer Productivity practice areas for AWS Strategic Industries as well as the Open Source Technical Field Community (TFC) of AWS Worldwide Specialists (WWSO).

## Yehuda Katz

### Heroku/Salesforce

 Yehuda is a True Believer in the power of the open web, especially when the web evolves as a collaboration between browser vendors, framework authors and application developers.

He is one of the creators of Ember.js, and a retired member of the Rust, the Ruby on Rails and jQuery Core Teams. He is an occasional member of ECMAScript’s TC39 standards committee, and a former member of the W3C’s TAG (Technical Architecture Group).

He was the co-author of the Extensible Web Manifesto, of which he is still very proud.

## Terence Lee

### Heroku/Salesforce

 Terence is an architect at Heroku where he helped create Classic Buildpacks and then later co-founded Cloud Native Buildpacks, a CNCF Incubation Project. In the Ruby community he’s been a maintainer on projects such as Ruby itself, Bundler, and Resque, but is mostly known for getting people together for #rubykaraoke. When he’s not going to an awesome event, he lives in Austin, TX where it’s acceptable to eat a taco for every meal of the day.

## Brett Weaver

### Intuit

 Brett Weaver is a Distinguished Engineer at Intuit. He has spent the last twenty two odd years at Intuit in various software development and systems engineer roles. He has been focused on building distributed, scalable architectures for Intuit’s flagship products including Quickbooks and TurboTax. Most recently, Brett has been leading architecture for Intuit’s internal services platform.

# Emeritus

## Gail Frederick

### Heroku/Salesforce

 As Heroku CTO, Gail is a Salesforce leader known for her technical excellence and drive to deliver. She stewards the opinionated magic that is Heroku’s developer platform. Prior to Heroku, Gail led engineering for Salesforce DX. Her values are integrity, impact, and joy. Previously, VP Engineering at eBay, where she was the 2019 winner of the John Donahoe Award, eBay’s highest leadership award, for building a new $3B annual GMB business by reinventing eBay’s developer ecosystem. Member of 2024 Curve Power List of 50 LGBTQ+ women and non-binary leaders. Formerly, executive advisory board member of Lesbians Who Tech, Business Governing Board member of OpenAPI Initiative at the Linux Foundation, and represented eBay at the founding of the Facebook-led Libra Initiative. Gail holds 9 software patents.

## Steren Giannini

### Google Cloud

 Steren is an engineer turned product manager. He is leading product management for Google Cloud’s serverless portfolio (Cloud Run, Cloud Run functions, App Engine). He is a founding member of Cloud Run.

https://twitter.com/steren

https://www.linkedin.com/in/steren

https://steren.fr

## Grace Jansen

### IBM

 Grace is a Java Champion and Developer Advocate at IBM. She has been with IBM since graduating with a Degree in Biology. Grace enjoys bringing a varied perspective to her projects and using her knowledge of biological systems to simplify complex software patterns. As a developer advocate, Grace builds POC’s, demos, sample applications and tutorials. She is a regular presenter at international conferences and has authored a book on reactive systems.

## Joe Kutner

### Salesforce

 Joe is co-founder of the Cloud Native Buildpacks project, which aims to make containerization more secure and developer friendly. He started the project in 2018 while working as DX Architect at Salesforce Heroku, and today is the DX Architect for Salesforce’s Hyperforce platform.

## James Ward

### AWS

 Professional software developer since 1997, with much of that time spent helping developers build software that doesn’t suck. A Typed Pure Functional Programming zealot who often compromises on his ideals to just get stuff done. Currently a Developer Advocate for AWS.

## I. Codebase

### One codebase tracked in revision control, many deploys

A twelve-factor app is always tracked in a version control system, such as Git, Mercurial, or Subversion. A copy of the revision tracking database is known as a *code repository*, often shortened to *code repo* or just *repo*.

A *codebase* is any single repo (in a centralized revision control system like Subversion), or any set of repos who share a root commit (in a decentralized revision control system like Git).

There is always a one-to-one correlation between the codebase and the app:

- If there are multiple codebases, it’s not an app – it’s a distributed system. Each component in a distributed system is an app, and each can individually comply with twelve-factor.
- Multiple apps sharing the same code is a violation of twelve-factor. The solution here is to factor shared code into libraries which can be included through the dependency manager.

There is only one codebase per app, but there will be many deploys of the app. A *deploy* is a running instance of the app. This is typically a production site, and one or more staging sites. Additionally, every developer has a copy of the app running in their local development environment, each of which also qualifies as a deploy.

The codebase is the same across all deploys, although different versions may be active in each deploy. For example, a developer has some commits not yet deployed to staging; staging has some commits not yet deployed to production. But they all share the same codebase, thus making them identifiable as different deploys of the same app.

## II. Dependencies

### Explicitly declare and isolate dependencies

Most programming languages offer a packaging system for distributing support libraries, such as CPAN for Perl or Rubygems for Ruby. Libraries installed through a packaging system can be installed system-wide (known as “site packages”) or scoped into the directory containing the app (known as “vendoring” or “bundling”).

**A twelve-factor app never relies on implicit existence of system-wide packages.** It declares all dependencies, completely and exactly, via a *dependency declaration* manifest. Furthermore, it uses a *dependency isolation* tool during execution to ensure that no implicit dependencies “leak in” from the surrounding system. The full and explicit dependency specification is applied uniformly to both production and development.

For example, Bundler for Ruby offers the `Gemfile` manifest format for dependency declaration and `bundle exec` for dependency isolation. In Python there are two separate tools for these steps – Pip is used for declaration and Virtualenv for isolation. Even C has Autoconf for dependency declaration, and static linking can provide dependency isolation. No matter what the toolchain, dependency declaration and isolation must always be used together – only one or the other is not sufficient to satisfy twelve-factor.

One benefit of explicit dependency declaration is that it simplifies setup for developers new to the app. The new developer can check out the app’s codebase onto their development machine, requiring only the language runtime and dependency manager installed as prerequisites. They will be able to set up everything needed to run the app’s code with a deterministic *build command*. For example, the build command for Ruby/Bundler is `bundle install`, while for Clojure/Leiningen it is `lein deps`.

Twelve-factor apps also do not rely on the implicit existence of any system tools. Examples include shelling out to ImageMagick or `curl`. While these tools may exist on many or even most systems, there is no guarantee that they will exist on all systems where the app may run in the future, or whether the version found on a future system will be compatible with the app. If the app needs to shell out to a system tool, that tool should be vendored into the app.

## III. Config

### Store config in the environment

An app’s *config* is everything that is likely to vary between deploys (staging, production, developer environments, etc). This includes:

- Resource handles to the database, Memcached, and other backing services
- Credentials to external services such as Amazon S3 or Twitter
- Per-deploy values such as the canonical hostname for the deploy

Apps sometimes store config as constants in the code. This is a violation of twelve-factor, which requires **strict separation of config from code**. Config varies substantially across deploys, code does not.

A litmus test for whether an app has all config correctly factored out of the code is whether the codebase could be made open source at any moment, without compromising any credentials.

Note that this definition of “config” does **not** include internal application config, such as `config/routes.rb` in Rails, or how code modules are connected in Spring. This type of config does not vary between deploys, and so is best done in the code.

Another approach to config is the use of config files which are not checked into revision control, such as `config/database.yml` in Rails. This is a huge improvement over using constants which are checked into the code repo, but still has weaknesses: it’s easy to mistakenly check in a config file to the repo; there is a tendency for config files to be scattered about in different places and different formats, making it hard to see and manage all the config in one place. Further, these formats tend to be language- or framework-specific.

**The twelve-factor app stores config in environment variables** (often shortened to

*env vars*or

*env*). Env vars are easy to change between deploys without changing any code; unlike config files, there is little chance of them being checked into the code repo accidentally; and unlike custom config files, or other config mechanisms such as Java System Properties, they are a language- and OS-agnostic standard.

Another aspect of config management is grouping. Sometimes apps batch config into named groups (often called “environments”) named after specific deploys, such as the `development`, `test`, and `production` environments in Rails. This method does not scale cleanly: as more deploys of the app are created, new environment names are necessary, such as `staging` or `qa`. As the project grows further, developers may add their own special environments like `joes-staging`, resulting in a combinatorial explosion of config which makes managing deploys of the app very brittle.

In a twelve-factor app, env vars are granular controls, each fully orthogonal to other env vars. They are never grouped together as “environments”, but instead are independently managed for each deploy. This is a model that scales up smoothly as the app naturally expands into more deploys over its lifetime.

## IV. Backing services

### Treat backing services as attached resources

A *backing service* is any service the app consumes over the network as part of its normal operation. Examples include datastores (such as MySQL or CouchDB), messaging/queueing systems (such as RabbitMQ or Beanstalkd), SMTP services for outbound email (such as Postfix), and caching systems (such as Memcached).

Backing services like the database are traditionally managed by the same systems administrators who deploy the app’s runtime. In addition to these locally-managed services, the app may also have services provided and managed by third parties. Examples include SMTP services (such as Postmark), metrics-gathering services (such as New Relic or Loggly), binary asset services (such as Amazon S3), and even API-accessible consumer services (such as Twitter, Google Maps, or Last.fm).

**The code for a twelve-factor app makes no distinction between local and third party services.** To the app, both are attached resources, accessed via a URL or other locator/credentials stored in the config. A deploy of the twelve-factor app should be able to swap out a local MySQL database with one managed by a third party (such as Amazon RDS) without any changes to the app’s code. Likewise, a local SMTP server could be swapped with a third-party SMTP service (such as Postmark) without code changes. In both cases, only the resource handle in the config needs to change.

Each distinct backing service is a *resource*. For example, a MySQL database is a resource; two MySQL databases (used for sharding at the application layer) qualify as two distinct resources. The twelve-factor app treats these databases as *attached resources*, which indicates their loose coupling to the deploy they are attached to.

Resources can be attached to and detached from deploys at will. For example, if the app’s database is misbehaving due to a hardware issue, the app’s administrator might spin up a new database server restored from a recent backup. The current production database could be detached, and the new database attached – all without any code changes.

## V. Build, release, run

### Strictly separate build and run stages

A codebase is transformed into a (non-development) deploy through three stages:

- The *build stage*is a transform which converts a code repo into an executable bundle known as a*build*. Using a version of the code at a commit specified by the deployment process, the build stage fetches vendors dependencies and compiles binaries and assets.
- The *release stage*takes the build produced by the build stage and combines it with the deploy’s current config. The resulting*release*contains both the build and the config and is ready for immediate execution in the execution environment.
- The *run stage*(also known as “runtime”) runs the app in the execution environment, by launching some set of the app’s processes against a selected release.

**The twelve-factor app uses strict separation between the build, release, and run stages.** For example, it is impossible to make changes to the code at runtime, since there is no way to propagate those changes back to the build stage.

Deployment tools typically offer release management tools, most notably the ability to roll back to a previous release. For example, the Capistrano deployment tool stores releases in a subdirectory named `releases`, where the current release is a symlink to the current release directory. Its `rollback` command makes it easy to quickly roll back to a previous release.

Every release should always have a unique release ID, such as a timestamp of the release (such as `2011-04-06-20:32:17`) or an incrementing number (such as `v100`). Releases are an append-only ledger and a release cannot be mutated once it is created. Any change must create a new release.

Builds are initiated by the app’s developers whenever new code is deployed. Runtime execution, by contrast, can happen automatically in cases such as a server reboot, or a crashed process being restarted by the process manager. Therefore, the run stage should be kept to as few moving parts as possible, since problems that prevent an app from running can cause it to break in the middle of the night when no developers are on hand. The build stage can be more complex, since errors are always in the foreground for a developer who is driving the deploy.

## VI. Processes

### Execute the app as one or more stateless processes

The app is executed in the execution environment as one or more *processes*.

In the simplest case, the code is a stand-alone script, the execution environment is a developer’s local laptop with an installed language runtime, and the process is launched via the command line (for example, `python my_script.py`). On the other end of the spectrum, a production deploy of a sophisticated app may use many process types, instantiated into zero or more running processes.

**Twelve-factor processes are stateless and share-nothing.** Any data that needs to persist must be stored in a stateful backing service, typically a database.

The memory space or filesystem of the process can be used as a brief, single-transaction cache. For example, downloading a large file, operating on it, and storing the results of the operation in the database. The twelve-factor app never assumes that anything cached in memory or on disk will be available on a future request or job – with many processes of each type running, chances are high that a future request will be served by a different process. Even when running only one process, a restart (triggered by code deploy, config change, or the execution environment relocating the process to a different physical location) will usually wipe out all local (e.g., memory and filesystem) state.

Asset packagers like django-assetpackager use the filesystem as a cache for compiled assets. A twelve-factor app prefers to do this compiling during the build stage. Asset packagers such as Jammit and the Rails asset pipeline can be configured to package assets during the build stage.

Some web systems rely on “sticky sessions” – that is, caching user session data in memory of the app’s process and expecting future requests from the same visitor to be routed to the same process. Sticky sessions are a violation of twelve-factor and should never be used or relied upon. Session state data is a good candidate for a datastore that offers time-expiration, such as Memcached or Redis.

## VII. Port binding

### Export services via port binding

Web apps are sometimes executed inside a webserver container. For example, PHP apps might run as a module inside Apache HTTPD, or Java apps might run inside Tomcat.

**The twelve-factor app is completely self-contained** and does not rely on runtime injection of a webserver into the execution environment to create a web-facing service. The web app **exports HTTP as a service by binding to a port**, and listening to requests coming in on that port.

In a local development environment, the developer visits a service URL like `http://localhost:5000/` to access the service exported by their app. In deployment, a routing layer handles routing requests from a public-facing hostname to the port-bound web processes.

This is typically implemented by using dependency declaration to add a webserver library to the app, such as Tornado for Python, Thin for Ruby, or Jetty for Java and other JVM-based languages. This happens entirely in *user space*, that is, within the app’s code. The contract with the execution environment is binding to a port to serve requests.

HTTP is not the only service that can be exported by port binding. Nearly any kind of server software can be run via a process binding to a port and awaiting incoming requests. Examples include ejabberd (speaking XMPP), and Redis (speaking the Redis protocol).

Note also that the port-binding approach means that one app can become the backing service for another app, by providing the URL to the backing app as a resource handle in the config for the consuming app.

## VIII. Concurrency

### Scale out via the process model

Any computer program, once run, is represented by one or more processes. Web apps have taken a variety of process-execution forms. For example, PHP processes run as child processes of Apache, started on demand as needed by request volume. Java processes take the opposite approach, with the JVM providing one massive uberprocess that reserves a large block of system resources (CPU and memory) on startup, with concurrency managed internally via threads. In both cases, the running process(es) are only minimally visible to the developers of the app.

**In the twelve-factor app, processes are a first class citizen.** Processes in the twelve-factor app take strong cues from the unix process model for running service daemons. Using this model, the developer can architect their app to handle diverse workloads by assigning each type of work to a *process type*. For example, HTTP requests may be handled by a web process, and long-running background tasks handled by a worker process.

This does not exclude individual processes from handling their own internal multiplexing, via threads inside the runtime VM, or the async/evented model found in tools such as EventMachine, Twisted, or Node.js. But an individual VM can only grow so large (vertical scale), so the application must also be able to span multiple processes running on multiple physical machines.

The process model truly shines when it comes time to scale out. The share-nothing, horizontally partitionable nature of twelve-factor app processes means that adding more concurrency is a simple and reliable operation. The array of process types and number of processes of each type is known as the *process formation*.

Twelve-factor app processes should never daemonize or write PID files. Instead, rely on the operating system’s process manager (such as systemd, a distributed process manager on a cloud platform, or a tool like Foreman in development) to manage output streams, respond to crashed processes, and handle user-initiated restarts and shutdowns.

## IX. Disposability

### Maximize robustness with fast startup and graceful shutdown

**The twelve-factor app’s processes are disposable, meaning they can be started or stopped at a moment’s notice.** This facilitates fast elastic scaling, rapid deployment of code or config changes, and robustness of production deploys.

Processes should strive to **minimize startup time**. Ideally, a process takes a few seconds from the time the launch command is executed until the process is up and ready to receive requests or jobs. Short startup time provides more agility for the release process and scaling up; and it aids robustness, because the process manager can more easily move processes to new physical machines when warranted.

Processes **shut down gracefully when they receive a SIGTERM** signal from the process manager. For a web process, graceful shutdown is achieved by ceasing to listen on the service port (thereby refusing any new requests), allowing any current requests to finish, and then exiting. Implicit in this model is that HTTP requests are short (no more than a few seconds), or in the case of long polling, the client should seamlessly attempt to reconnect when the connection is lost.

For a worker process, graceful shutdown is achieved by returning the current job to the work queue. For example, on RabbitMQ the worker can send a `NACK`; on Beanstalkd, the job is returned to the queue automatically whenever a worker disconnects. Lock-based systems such as Delayed Job need to be sure to release their lock on the job record. Implicit in this model is that all jobs are reentrant, which typically is achieved by wrapping the results in a transaction, or making the operation idempotent.

Processes should also be **robust against sudden death**, in the case of a failure in the underlying hardware. While this is a much less common occurrence than a graceful shutdown with `SIGTERM`, it can still happen. A recommended approach is use of a robust queueing backend, such as Beanstalkd, that returns jobs to the queue when clients disconnect or time out. Either way, a twelve-factor app is architected to handle unexpected, non-graceful terminations. Crash-only design takes this concept to its logical conclusion.

## X. Dev/prod parity

### Keep development, staging, and production as similar as possible

Historically, there have been substantial gaps between development (a developer making live edits to a local deploy of the app) and production (a running deploy of the app accessed by end users). These gaps manifest in three areas:

- **The time gap**: A developer may work on code that takes days, weeks, or even months to go into production.
- **The personnel gap**: Developers write code, ops engineers deploy it.
- **The tools gap**: Developers may be using a stack like Nginx, SQLite, and OS X, while the production deploy uses Apache, MySQL, and Linux.

**The twelve-factor app is designed for continuous deployment by keeping the gap between development and production small.** Looking at the three gaps described above:

- Make the time gap small: a developer may write code and have it deployed hours or even just minutes later.
- Make the personnel gap small: developers who wrote code are closely involved in deploying it and watching its behavior in production.
- Make the tools gap small: keep development and production as similar as possible.

Summarizing the above into a table:

| Traditional app | Twelve-factor app | |
|---|---|---|
| Time between deploys | Weeks | Hours |
| Code authors vs code deployers | Different people | Same people |
| Dev vs production environments | Divergent | As similar as possible |

Backing services, such as the app’s database, queueing system, or cache, is one area where dev/prod parity is important. Many languages offer libraries which simplify access to the backing service, including *adapters* to different types of services. Some examples are in the table below.

| Type | Language | Library | Adapters |
|---|---|---|---|
| Database | Ruby/Rails | ActiveRecord | MySQL, PostgreSQL, SQLite |
| Queue | Python/Django | Celery | RabbitMQ, Beanstalkd, Redis |
| Cache | Ruby/Rails | ActiveSupport::Cache | Memory, filesystem, Memcached |

Developers sometimes find great appeal in using a lightweight backing service in their local environments, while a more serious and robust backing service will be used in production. For example, using SQLite locally and PostgreSQL in production; or local process memory for caching in development and Memcached in production.

**The twelve-factor developer resists the urge to use different backing services between development and production**, even when adapters theoretically abstract away any differences in backing services. Differences between backing services mean that tiny incompatibilities crop up, causing code that worked and passed tests in development or staging to fail in production. These types of errors create friction that disincentivizes continuous deployment. The cost of this friction and the subsequent dampening of continuous deployment is extremely high when considered in aggregate over the lifetime of an application.

Lightweight local services are less compelling than they once were. Modern backing services such as Memcached, PostgreSQL, and RabbitMQ are not difficult to install and run thanks to modern packaging systems, such as Homebrew and apt-get. Alternatively, declarative provisioning tools such as Chef and Puppet combined with light-weight virtual environments such as Docker and Vagrant allow developers to run local environments which closely approximate production environments. The cost of installing and using these systems is low compared to the benefit of dev/prod parity and continuous deployment.

Adapters to different backing services are still useful, because they make porting to new backing services relatively painless. But all deploys of the app (developer environments, staging, production) should be using the same type and version of each of the backing services.

## XI. Logs

### Treat logs as event streams

*Logs* provide visibility into the behavior of a running app. In server-based environments they are commonly written to a file on disk (a “logfile”); but this is only an output format.

Logs are the stream of aggregated, time-ordered events collected from the output streams of all running processes and backing services. Logs in their raw form are typically a text format with one event per line (though backtraces from exceptions may span multiple lines). Logs have no fixed beginning or end, but flow continuously as long as the app is operating.

**A twelve-factor app never concerns itself with routing or storage of its output stream.** It should not attempt to write to or manage logfiles. Instead, each running process writes its event stream, unbuffered, to `stdout`. During local development, the developer will view this stream in the foreground of their terminal to observe the app’s behavior.

In staging or production deploys, each process’ stream will be captured by the execution environment, collated together with all other streams from the app, and routed to one or more final destinations for viewing and long-term archival. These archival destinations are not visible to or configurable by the app, and instead are completely managed by the execution environment. Open-source log routers (such as Logplex and Fluentd) are available for this purpose.

The event stream for an app can be routed to a file, or watched via realtime tail in a terminal. Most significantly, the stream can be sent to a log indexing and analysis system such as Splunk, or a general-purpose data warehousing system such as Hadoop/Hive. These systems allow for great power and flexibility for introspecting an app’s behavior over time, including:

- Finding specific events in the past.
- Large-scale graphing of trends (such as requests per minute).
- Active alerting according to user-defined heuristics (such as an alert when the quantity of errors per minute exceeds a certain threshold).

## XII. Admin processes

### Run admin/management tasks as one-off processes

The process formation is the array of processes that are used to do the app’s regular business (such as handling web requests) as it runs. Separately, developers will often wish to do one-off administrative or maintenance tasks for the app, such as:

- Running database migrations (e.g. `manage.py migrate`in Django,`rake db:migrate`in Rails).
- Running a console (also known as a REPL shell) to run arbitrary code or inspect the app’s models against the live database. Most languages provide a REPL by running the interpreter without any arguments (e.g. `python`or`perl`) or in some cases have a separate command (e.g.`irb`for Ruby,`rails console`for Rails).
- Running one-time scripts committed into the app’s repo (e.g. `php scripts/fix_bad_records.php`).

One-off admin processes should be run in an identical environment as the regular long-running processes of the app. They run against a release, using the same codebase and config as any process run against that release. Admin code must ship with application code to avoid synchronization issues.

The same dependency isolation techniques should be used on all process types. For example, if the Ruby web process uses the command `bundle exec thin start`, then a database migration should use `bundle exec rake db:migrate`. Likewise, a Python program using Virtualenv should use the vendored `bin/python` for running both the Tornado webserver and any `manage.py` admin processes.

Twelve-factor strongly favors languages which provide a REPL shell out of the box, and which make it easy to run one-off scripts. In a local deploy, developers invoke one-off admin processes by a direct shell command inside the app’s checkout directory. In a production deploy, developers can use ssh or other remote command execution mechanism provided by that deploy’s execution environment to run such a process.

# Úvod

V moderní době je software často dodáván jako služba, což označujeme pojmem *webová aplikace* nebo *software-as-a-service (SaaS)*. Twelve-factor metodika slouží pro vytváření (SaaS) aplikací, které:

- Používají **deklarativní**formáty pro nastavení automatizace, což vede k minimalizaci času a nákladů potřebných pro začlenění nových vývojářů do projektu;
- Mají **jasný kontrakt**s operačním systémem ve kterém běží, čímž je umožněna**maximální přenositelnost**mezi různými běhovými prostředími;
- Jsou vhodné pro **nasazení**na moderních**cloudových platformách**, což eliminuje potřebu správy serverů a podpůrných systémů;
- **Minimalizují rozdíly**mezi vývojovým a produkčním prostředím, čímž umožňují- **průběžné nasazovaní**a maximální flexibilitu;
- Umožňují **vyškálování**bez výrazných změn v nástrojích, architektuře nebo vývojových postupech.

Twelve-factor metodiku lze použít na aplikace napsané v jakémkoliv programovacím jazyce a používající libovolnou kombinaci podpůrných služeb (databáze, fronty, vyrovnávací paměť atd.).

# Pozadí

Přispěvatelé tohoto dokumentu se přímo podíleli na vývoji a nasazení stovek aplikací a byli svědky vývoje, provozu a škálování stovek tisíc aplikací prostřednictvím své práce na platformě Heroku.

Tento dokument shromažďuje všechny naše zkušenosti a postřehy týkající se široké škály aplikací typu software-as-a-service v divočině. Jedná se o sadu ideálních postupů pro vývoj aplikací se zvláštní pozorností věnovanou dynamice organického růstu aplikace v průběhu času, dynamice spolupráce mezi vývojáři pracujícími na kódu aplikace a vyhýbání se nákladům na erozi softwaru.

Naší motivací je zvyšovat povědomí o některých systémových problémech, které jsme zaznamenali v moderním vývoji aplikací, poskytnout společnou slovní zásobu pro diskusi o těchto problémech a nabídnout rozsáhlou sadu koncepčních řešení těchto problémů s doprovodnou terminologií. Formát je inspirován knihami Martina Fowlera *Patterns of Enterprise Application Architecture* a *Refactoring*.

# Kdo by měl číst tento dokument?

Každý vývojář pracující na aplikaci, která běží jako služba. Systémoví inženýři, kteří takové aplikace nasazují nebo spravují.

# Einführung

Heute wird Software oft als Dienst geliefert - auch *Web App* oder *Software-As-A-Service* genannt. Die Zwölf-Faktoren-App ist eine Methode um Software-As-A-Service Apps zu bauen die:

- **deklarative**Formate benutzen für die Automatisierung der Konfiguration, um Zeit und Kosten für neue Entwickler im Projekt zu minimieren;
- einen **sauberen Vertrag**mit dem zugrundeliegenden Betriebssystem haben,**maximale Portierbarkeit**zwischen Ausführungsumgebungen bieten;
- sich für das **Deployment**auf modernen**Cloud-Plattformen**eignen, die Notwendigkeit von Servern und Serveradministration vermeiden;
- die **Abweichung minimieren**zwischen Entwicklung und Produktion, um**Continuous Deployment**für maximale Agilität ermöglichen;
- und **skalieren**können ohne wesentliche Änderungen im Tooling, in der Architektur oder in den Entwicklungsverfahren.

Die Zwölf-Faktoren-Methode kann auf Apps angewendet werden, die in einer beliebigen Programmiersprache geschrieben sind, und die eine beliebige Kombination von unterstützenden Diensten benutzen (Datenbank, Queue, Cache, …)

# Hintergrund

Die Mitwirkenden an diesem Dokument waren direkt beteiligt an der Entwicklung und dem Deployment von hunderten von Apps und wurden Zeugen bei der Entwicklung, beim Betrieb und der Skalierung von hunderttausenden von Apps im Rahmen unserer Arbeit an der Heroku-Plattform.

Dieses Dokument ist eine Synthese all unserer Erfahrungen und der Beobachtungen einer großen Bandbreite von Software-As-A-Service Apps. Es ist eine Bestimmung der idealen Praktiken bei der App-Entwicklung mit besonderem Augenmerk auf die Dynamik des organischen Wachstums einer App über die Zeit, die Dynamik der Zusammenarbeit zwischen den Entwicklern die an einer Codebase zusammenarbeiten und der Vermeidung der Kosten von Software-Erosion.

Unsere Motivation ist, das Bewusstsein zu schärfen für systembedingte Probleme in der aktuellen Applikationsentwicklung, ein gemeinsames Vokabular zur Diskussion dieser Probleme zu liefern und ein Lösungsportfolio zu diesen Problemen mit einer zugehörigen Terminologie anzubieten. Das Format ist angelehnt an Martin Fowlers Bücher *Patterns of Enterprise Application Architecture* und *Refactoring*.

# Wer sollte dieses Dokument lesen?

Jeder Entwickler der Apps baut, die als Dienst laufen. Administratoren, die solche Apps managen oder deployen.

# Εισαγωγή

Στη μοντέρνα εποχή, το λογισμικό συνήθως παρέχεται ως υπηρεσία: καλούμενο *εφαρμογές ιστού* (*web apps*), ή *λογισμικό-ως-υπηρεσία* (*software-as-a-service*). Η εφαρμογή δώδεκα παραγόντων είναι μια μεθοδολογία κατασκευής εφαρμογών λογισμικού-ως-υπηρεσίας όπου:

- Χρησιμοποιεί **δηλωτικές**μορφές (**declarative**formats) για να στήσει τον αυτοματισμό, να ελαχιστοποιήσει το χρόνο και το κόστος για νέους προγραμματιστές να συμμετέχουν στο έργο,
- Έχει ένα **καθαρό συμβόλαιο**(**clean contract**) με το υποκείμενο λειτουργικό σύστημα, προσφέροντας**μέγιστη φορητότητα**(**maximum portability**) μεταξύ περιβαλλόντων εκτέλεσης,
- Είναι κατάλληλη για **ανάπτυξη**(**deployment**) σε μοντέρνες πλατφόρμες υπολογιστικού νέφους (**cloud platforms**), καθιστώντας περιττή την ανάγκη για εξυπηρετητές και διαχείριση συστημάτων,
- **Ελαχιστοποιεί την απόκλιση**μεταξύ υλοποίησης (development) και παραγωγής (production), διευκολύνοντας την- **συνεχή ανάπτυξη**(- **continuous deployment**) για μέγιστη ευκινησία (maximum agility),
- Και μπορεί να **κλιμακωθεί προς τα πάνω**(**scale up**) χωρίς σημαντικές αλλαγές στα εργαλεία, στην αρχιτεκτονική, ή στις πρακτικές υλοποίησης.

Η μεθοδολογία δώδεκα παραγόντων μπορεί να εφαρμοστεί σε εφαρμογές οι οποίες είναι γραμμένες σε οποιαδήποτε γλώσσα προγραμματισμού, και οι οποίες χρησιμοποιούν οποιοδήποτε συνδυασμό από υπηρεσίες υποστήριξης (βάση δεδομένων, ουρά εργασιών, προσωρινή μνήμη, κλπ).

# Introducción

En estos tiempos, el software se está distribuyendo como un servicio: se le denomina *web apps*, o *software as a service* (SaaS). “The twelve-factor app” es una metodología para construir aplicaciones SaaS que:

- Usan formatos **declarativos**para la automatización de la configuración, para minimizar el tiempo y el coste que supone que nuevos desarrolladores se unan al proyecto;
- Tienen un **contrato claro**con el sistema operativo sobre el que trabajan, ofreciendo la**máxima portabilidad**entre los diferentes entornos de ejecución;
- Son apropiadas para **desplegarse**en modernas**plataformas en la nube**, obviando la necesidad de servidores y administración de sistemas;
- **Minimizan las diferencias**entre los entornos de desarrollo y producción, posibilitando un- **despliegue continuo**para conseguir la máxima agilidad;
- Y pueden **escalar**sin cambios significativos para las herramientas, la arquitectura o las prácticas de desarrollo.

La metodología “twelve-factor” puede ser aplicada a aplicaciones escritas en cualquier lenguaje de programación, y cualquier combinación de ‘backing services’ (bases de datos, colas, memoria cache, etc).

# Contexto

Los colaboradores de este documento han estado involucrados directamente en el desarrollo y despliegue de cientos de aplicaciones, y han sido testigos indirectos del desarrollo, las operaciones y el escalado de cientos de miles de aplicaciones mediante nuestro trabajo en la plataforma Heroku.

Este documento sintetiza toda nuestra experiencia y observaciones en una amplia variedad de aplicaciones SaaS. Es la triangulación entre practicas ideales para el desarrollo de aplicaciones, prestando especial atención a las dinámicas del crecimiento natural de una aplicación a lo largo del tiempo, las dinámicas de colaboración entre desarrolladores que trabajan en el código base de las aplicaciones y evitando el coste de la entropía del software.

Nuestra motivación es mejorar la concienciación sobre algunos problemas sistémicos que hemos observado en el desarrollo de las aplicaciones modernas, aportar un vocabulario común que sirva para discutir sobre estos problemas, y ofrecer un conjunto de soluciones conceptualmente robustas para esos problemas acompañados de su correspondiente terminología. El formato está inspirado en los libros de Martin Fowler *Patterns of Enterprise Application Architecture* y *Refactoring*.

# معرفی

در عصر مدرن، نرمافزار معمولاً به عنوان یک سرویس به نام *برنامههای کاربردی تحت وب* یا *نرمافزار به عنوان سرویس* ارائه می شود. برنامهی کاربردی دوازده-سازه روشی برای ساخت برنامههای کاربردی نرمافزاری به عنوان سرویس است که:

- از قالبهای **قابل توصیف**برای خودکارسازی راهاندازی استفاده میکند تا زمان و هزینه را برای توسعهدهندگان جدیدی که به پروژه می پیوندند به حداقل برساند.
- یک **قرارداد تمیز**با سیستم عامل میزبان دارد تا**حداکثر قابلیت حمل**را بین محیطهای اجرا ارائه دهد.
- برای **استقرار**در**سکوهای ابری مدرن**، بدون نیاز به سرورها و مدیریت سیستمها مناسب است.
- **واگرایی**بین توسعه و عملیات را به حداقل میرساند تا- **استقرار مستمر**، حداکثر چابکی را به ارمغان آورد.
- و میتواند بدون تغییر قابل توجه در ابزار، معماری یا شیوههای توسعه، **مقیاسپذیر**شود.

روش دوازده-سازه را می توان برای برنامههایی که به هر زبان برنامهنویسی نوشته شدهاند و از هر ترکیبی از خدمات پشتیبان (پایگاه داده، صف، کش حافظه و غیره) استفاده میکنند، اعمال کرد.

# پسزمینه

مشارکتکنندگان در این نوشتار، به طور مستقیم در توسعه و استقرار صدها برنامهی کاربردی مشارکت داشتهاند و بهطور غیرمستقیم شاهد توسعه، عملیات و مقیاسپذیری صدها هزار برنامهی کاربردی از طریق کار ما در Heroku بودهاند.

این نوشتار، همهی تجربیات و مشاهدههای ما را در مورد طیف گستردهای از برنامههای کاربردی نرم افزاری به عنوان سرویس در صنعت ترکیب میکند. سه ضلعی که بر روی شیوههای ایدهآل برای توسعهی برنامهی کاربردی، توجه ویژه به پویایی رشد ذاتی یک برنامهی کاربردی در طول زمان، پویایی همکاری بین توسعهدهندگانی که روی کدنویسی برنامهی کاربردی کار میکنند، و جلوگیری از هزینهی زوال نرمافزار بنا شده است.

انگیزهی ما افزایش آگاهی از برخی مشکلات سیستمی است که در توسعهی برنامههای کاربردی مدرن دیدهایم، واژگان مشترکی برای بحث در مورد آن مشکلات ارائه می دهیم، و مجموعهای از راهحلهای مفهومی گسترده برای آن مشکلات را همراه با اصطلاحات ارائه میدهیم. این قالب از کتابهای مارتین فاولر به نامهای *الگوهای معماری برنامههای کاربردی سازمانی* و *بازسازی* الهام گرفته شده است.

# Introduction

À l’époque actuelle, les logiciels sont régulièrement délivrés en tant que services : on les appelle des *applications web* (web apps), ou *logiciels en tant que service* (*software-as-a-service*). L’application 12 facteurs est une méthodologie pour concevoir des logiciels en tant que service qui :

- Utilisent des formats **déclaratifs**pour mettre en oeuvre l’automatisation, pour minimiser le temps et les coûts pour que de nouveaux développeurs rejoignent le projet;
- Ont un **contrat propre**avec le système d’exploitation sous-jacent, offrant une**portabilité maximum**entre les environnements d’exécution;
- Sont adaptés à des **déploiements**sur des**plateformes cloud**modernes, rendant inutile le besoin de serveurs et de l’administration de systèmes;
- **Minimisent la divergence**entre le développement et la production, ce qui permet le- **déploiement continu**pour une agilité maximum;
- et peuvent **grossir verticalement**sans changement significatif dans les outils, l’architecture ou les pratiques de développement;

La méthodologie 12 facteurs peut être appliquée à des applications écrites dans tout langage de programmation, et qui utilisent tout type de services externes (base de données, file, cache mémoire, etc.)

# Contexte

Les contributeurs de ce document ont été directement impliqués dans le développement et le déploiement de centaines d’applications, et ont vu, indirectement, le développement, la gestion et le grossissement de centaines de milliers d’applications via le travail fait sur la plateforme Heroku.

Ce document fait la synthèse de toutes nos expériences et observations sur une large variété d’applications software-as-a-service. C’est la triangulation de pratiques idéales pour le développement d’applications, en portant un soin tout particulier aux dynamiques de la croissance organique d’une application au cours du temps, les dynamiques de la collaboration entre les développeurs qui travaillent sur le code de l’application, en évitant le coût de la lente détérioration du logiciel dans un environnement qui évolue (en).

Notre motivation est de faire prendre conscience de certains problèmes systémiques que nous avons rencontrés dans le développement d’applications modernes, afin de fournir un vocabulaire partagé pour discuter ces problèmes, et pour offrir un ensemble de solutions conceptuelles générales à ces problèmes, ainsi que la terminologie correspondante. Le format est inspiré par celui des livres de Martin Fowler *Patterns of Enterprise Application Architecture (en)* et *Refactoring (en)*.

# Introduzione

Nell’era moderna, il software viene fornito sempre più di frequente come servizio (delivered as a service): si parla di *web app* o *software as a service* (SaaS). La **twelve-factor app** è una metodologia di sviluppo orientata alla costruzione di applicazioni software-as-a-service che:

- Seguono un formato **dichiarativo**per l’automazione della configurazione, minimizzando tempi e costi di ingresso per ogni sviluppatore che si aggiunge al progetto;
- **Si interfacciano in modo pulito**con il sistema operativo sottostante, in modo tale da offrire la- **massima portabilità**sui vari ambienti di esecuzione;
- Sono **adatte allo sviluppo**sulle più recenti**cloud platform**, ovviando alla necessità di server e amministrazioni di sistema;
- **Minimizzano la divergenza**tra sviluppo e produzione, permettendo il- **continuous deployment**per una massima “agilità”;
- Possono **scalare significativamente**senza troppi cambiamenti ai tool, all’architettura e al processo di sviluppo;

La metodologia twelve-factor può essere applicata a ogni software, scritto in qualsiasi linguaggio di programmazione, che fa uso di una serie di servizi come database, code, cache e così via.

# Background

Chi ha scritto questo documento è stato coinvolto direttamente nella realizzazione e nel deployment di centinaia di applicazioni, e ha indirettamente assistito allo sviluppo, le operazioni e lo scaling di centinaia (o migliaia) di app tramite il proprio lavoro sulla piattaforma Heroku.

Questo documento riassume tutta quella che è stata la nostra esperienza, basata sull’osservazione di un grande numero di applicazioni SaaS. Si tratta di una “triangolazione” di pratiche di sviluppo ideali (con una particolare attenzione alla crescita organica dell’app nel tempo), la collaborazione dinamica nel corso del tempo tra gli sviluppatori sulla codebase e la necessità di evitare i costi di software erosion.

La nostra motivazione è di far crescere la consapevolezza intorno ad alcuni problemi sistemici che abbiamo scoperto nello sviluppo di applicazioni moderne, cercando di fornire un vocabolario condiviso per la discussione di tali problemi. Oltre, ovviamente, a offrire delle soluzioni concettuali a queste situazioni (senza però tralasciare il fattore tecnologia). Questo format si rifà ai libri di Martin Fowler *Patterns of Enterprise Application Architecture* e *Refactoring*.

# はじめに

現代では、ソフトウェアは一般にサービスとして提供され、*Webアプリケーション* や *Software as a Service* と呼ばれる。Twelve-Factor Appは、次のようなSoftware as a Serviceを作り上げるための方法論である。

- セットアップ自動化のために **宣言的な**フォーマットを使い、プロジェクトに新しく加わった開発者が要する時間とコストを最小化する。
- 下層のOSへの **依存関係を明確化**し、実行環境間での**移植性を最大化**する。
- モダンな **クラウドプラットフォーム**上への**デプロイ**に適しており、サーバー管理やシステム管理を不要なものにする。
- 開発環境と本番環境の **差異を最小限**にし、アジリティを最大化する**継続的デプロイ**を可能にする。
- ツール、アーキテクチャ、開発プラクティスを大幅に変更することなく **スケールアップ**できる。

Twelve-Factorの方法論は、どのようなプログラミング言語で書かれたアプリケーションにでも適用できる。また、どのようなバックエンドサービス（データベース、メッセージキュー、メモリキャッシュなど）の組み合わせを使っていても適用できる。

# 머리말

최근 소프트웨어를 서비스 형태로 제공하는게 일반화 되면서, 웹앱 혹은 SaaS(Software As A Service)라고 부르게 되었다. Twelve-Factor app은 아래 특징을 가진 SaaS 앱을 만들기 위한 방법론이다.

- 설정 자동화를 위한 **절차(declarative)**를 체계화 하여 새로운 개발자가 프로젝트에 참여하는데 드는 시간과 비용을 최소화한다.
- OS에 따라 **달라지는 부분을 명확히**하고, 실행 환경 사이의**이식성을 극대화**한다.
- 최근 등장한 **클라우드 플랫폼****배포에**적합하고, 서버와 시스템의 관리가 필요없게 된다.
- 개발 환경과 운영 환경의 **차이를 최소화**하고 민첩성을 극대화하기 위해**지속적인 배포**가 가능하다.
- 툴, 아키텍처, 개발 방식을 크게 바꾸지 않고 **확장(scale up)**할 수 있다.

Twelve-Factor 방법론은 어떤 프로그래밍 언어로 작성된 앱에도 적용할 수 있고 백엔드 서비스(데이터베이스, 큐, 메모리 캐시 등)와 다양한 조합으로 사용할 수 있다.

# Wprowadzenie

We współczesnym świecie oprogramowanie jest powszechnie wytwarzane w formie usługi, nazywane *software-as-service (SaaS)* lub aplikacjami internetowymi. Dwanaście aspektów aplikacji jest metodologią budowania aplikacji SaaS, które:

- Używają **deklaratywnego**formatu by zautomatyzować konfigurację aplikacji w celu zmniejszenia czasu i kosztów dołączenia nowych programistów do projektu;
- Mają **czysty kontrakt**z systemem operacyjnym, umożliwiając**jak największą możliwość przenoszenia**pomiędzy środowiskami, w których działają;
- Są dopasowane do **wdrożenia**na nowoczesne**chmury obliczeniowe**, zapobiegając potrzebie użycia serwerów i administracji systemu;
- **Minimalizują rozbieżności**pomiędzy środowiskami developerskimi i produkcyjnymi, umożliwiając- **nieustanne wdrażanie aplikacji**by zmaksymalizować prędkość zmian;
- I mogą **skalować się**bez większej zmiany narzędzi, architektury, czy sposobu pracy zespołu.

Metodologia dwunastu aspektów może być stosowana do aplikacji napisanych w każdym języku programowania i wykorzystujących dowolną kombinację usług wspierających (bazy danych, kolejki, cache pamięci etc).

# Background

Kontrybutorzy tego dokumentu byli bezpośrednio zaangażowani w tworzenie i wdrażanie setek aplikacji i pośrednio byli świadkami produkcji, działania i skalowania setek tysięcy aplikacji dzięki naszej pracy na platformie Heroku.

Ten dokument jest podsumowaniem całego naszego doświadczenia i obserwacji szerokiej gamy aplikacji SaaS. Jest on połączeniem idealnych praktyk developmentu, zwracania szczególnej uwagi na naturalny rozrost aplikacji w czasie, dynamiki współpracy developerów pracujących nad jednym codebase’m, oraz unikania kosztów gnijącego oprogramowania.

Naszym celem jest podniesienie poziomu świadomości o podstawowych problemach, które dostrzegliśmy przy tworzeniu nowoczesnych aplikacji, zapewnienie wspólnego słownictwa do rozmowy o tych problemach oraz zaoferowanie ogólnych rozwiązań dla tych problemów wraz z towarzyszącą terminologią. Format dokumentu jest inspirowany książkami Martina Fowlera *Patterns of Enterprise Application Architecture* oraz *Refactoring*.

# Dla kogo przeznaczony jest ten dokument?

Dla każdego developera tworzącego aplikacje, które działają jako usługa. Dla każdego Dev-opsa, który wdraża i zarządza takimi aplikacjami.

# Introdução

Na era moderna, software é comumente entregue como um serviço: denominados *web apps*, ou *software-como-serviço*. A aplicação doze-fatores é uma metodologia para construir softwares-como-serviço que:

- Usam formatos **declarativos**para automatizar a configuração inicial, minimizar tempo e custo para novos desenvolvedores participarem do projeto;
- Tem um **contrato claro**com o sistema operacional que o suporta, oferecendo**portabilidade máxima**entre ambientes que o executem;
- São adequados para **implantação**em modernas**plataformas em nuvem**, evitando a necessidade por servidores e administração do sistema;
- **Minimizam a divergência**entre desenvolvimento e produção, permitindo a- **implantação contínua**para máxima agilidade;
- E podem **escalar**sem significativas mudanças em ferramentas, arquiteturas, ou práticas de desenvolvimento.

A metodologia doze-fatores pode ser aplicada a aplicações escritas em qualquer linguagem de programação, e que utilizem qualquer combinação de serviços de suportes (banco de dados, filas, cache de memória, etc).

# Experiência

Os contribuidores deste documento estão diretamente envolvidos no desenvolvimento e implantação de centenas de aplicações, e indiretamente testemunhando o desenvolvimento, operação e escalada de centenas de milhares de aplicações através de seu trabalho na plataforma Heroku.

Este documento sintetiza toda nossa experiência e observação em uma variedade de aplicações que operam como software-como-serviço. Isto é a triangulação de práticas ideais ao desenvolvimento de software, com uma atenção particular a respeito das dinâmicas de crescimento orgânico de uma aplicação ao longo do tempo, a dinâmica de colaboração entre desenvolvedores trabalhando em uma base de código, e evitando os custos de erosão de software

Nossa motivação é aumentar a consciência de alguns problemas sistêmicos que temos visto no desenvolvimento de aplicações modernas, prover um vocabulário comum para discussão destes, e oferecer um amplo conjunto de soluções conceituais para esses problemas com a terminologia que os acompanha. O formato é inspirado nos livros de Martin Fowler *Padrões de Arquitetura de Aplicações Enterprise* e *Refatorando*.

# Введение

В наши дни программное обеспечение обычно распространяется в виде сервисов, называемых *веб-приложения* (web apps) или *software-as-a-service* (SaaS). Приложение двенадцати факторов — это методология для создания SaaS-приложений, которые:

- Используют **декларативный**формат для описания процесса установки и настройки, что сводит к минимуму затраты времени и ресурсов для новых разработчиков, подключённых к проекту;
- Имеют **соглашение**с операционной системой, предполагающее**максимальную переносимость**между средами выполнения;
- Подходят для **развёртывания**на современных**облачных платформах**, устраняя необходимость в серверах и системном администрировании;
- **Сводят к минимуму расхождения**между средой разработки и средой выполнения, что позволяет использовать- **непрерывное развёртывание**(continuous deployment) для максимальной гибкости;
- И могут **масштабироваться**без существенных изменений в инструментах, архитектуре и практике разработки.

Методология двенадцати факторов может быть применена для приложений, написанных на любом языке программирования и использующих любые комбинации сторонних служб (backing services) (базы данных, очереди сообщений, кэш-памяти, и т.д.).

# Предпосылки

Участники, внёсшие вклад в этот документ, были непосредственно вовлечены в разработку и развёртывание сотен приложений и косвенно были свидетелями разработки, выполнения и масштабирования сотен тысяч приложений во время нашей работы над платформой Heroku.

В этом документе обобщается весь наш опыт использования и наблюдения за самыми разнообразными SaaS-приложениями в дикой природе. Документ является объединением трёх идеальных подходов к разработке приложений: уделение особого внимания динамике органического роста приложения с течением времени, динамике сотрудничества разработчиков, работающих над кодовой базой приложения, и устранение последствий эрозии программного обеспечения.

Наша мотивация заключается в повышении осведомлённости о некоторых системных проблемах, которые мы встретили в практике разработки современных приложений, а также для того, чтобы предоставить общие основные понятия для обсуждения этих проблем и предложить набор общих концептуальных решений этих проблем с сопутствующей терминологией. Формат навеян книгами Мартина Фаулера (Martin Fowler) *Patterns of Enterprise Application Architecture* и *Refactoring*.

# Кому следует читать этот документ?

Разработчикам, которые создают SaaS-приложения. Ops инженерам, выполняющим развёртывание и управление такими приложениями.

# Úvod

V modernej dobe sa zvyčajne softvér dodáva ako služba: nazýva sa *webová aplikácia*, alebo *software-as-a-service*. Dvanásť faktorová aplikácia je metodológia na budovanie software-as-a-service aplikácií, ktoré:

- Používajú **deklaratívne**formáty na automatizáciu nastavení, a minimalizáciu času a nákladov pre nových developerov, ktorí sa začlenia do projektu;
- Obsahuje **jasnú zmluvu**s operačným systémom, nad ktorým bežia, čím umožňujú**maximálnu portabilitu**medzi rôznymi prostrediami;
- Sú vhodné na **nasadenie**na moderné**cloudové platformy**, čím vylučujú potrebu serverov a systémových administrátorov;
- **Minimalizujú rozdiely**medzi vývojom a produkciou, čím umôžňujú- **continuous deployment**s maximálnou agilnosťou;
- A sú **škálovateľné**bez významných zmien v nástrojoch, architektúre alebo vývojárskych postupoch.

Dvanásť faktorová metodológia sa dá použiť na aplikácie písané v akomkoľvek programovacom jazyku, ktoré používajú akúkoľvek kombináciu podporných služieb (databáza, fronta, pamäťová cache, atď).

# Background

The contributors to this document have been directly involved in the development and deployment of hundreds of apps, and indirectly witnessed the development, operation, and scaling of hundreds of thousands of apps via our work on the Heroku platform.

This document synthesizes all of our experience and observations on a wide variety of software-as-a-service apps in the wild. It is a triangulation on ideal practices for app development, paying particular attention to the dynamics of the organic growth of an app over time, the dynamics of collaboration between developers working on the app’s codebase, and avoiding the cost of software erosion.

Our motivation is to raise awareness of some systemic problems we’ve seen in modern application development, to provide a shared vocabulary for discussing those problems, and to offer a set of broad conceptual solutions to those problems with accompanying terminology. The format is inspired by Martin Fowler’s books *Patterns of Enterprise Application Architecture* and *Refactoring*.

# Kto by si mal prečítať tento dokument?

Každý vývojár pracujúci na aplikácii, ktorá beží ako služba. Systémoví administrátori, ktorý také aplikácie nasadzujú.

# บทนำ

ในยุคสมัยใหม่ ซอฟต์แวร์ถูกส่งมอบทั่วไปเป็นบริการ: เรียกว่า *web apps*, หรือ *software-as-service*. twelve-factor app เป็นหลัการสำหรับสร้างแอพพลิเคชัน software-as-a-service ที่:

- ใช้รูปแบบ **declarative**สำหรับติดตั้งระบบอัตโนมัต เพื่อลดเวลาและค่าใช้จ่ายสำหรับนักพัฒนาใหม่ที่เข้าร่วมกับโครงการ;
- มี **clean contract**กับระบบปฏิบัติการที่แอพพลิเคชันทำงานด้วย นำเสนอ**maximun portibility**ระหว่างสิ่งแวดล้อมที่ระบบทำงาน;
- เหมาะสมสำหรับ **deployment**บน**cloud platforms**สมัยใหม่, ลดความต้องการของเซิร์ฟเวอร์และผู้ดูแลระบบ;
- **Maximized divergence**ระหว่างการพัฒนาและการใช้งานจริง ด้วยการใช้- **continuous deployment**เพื่อเพิ่มความเร็วสูงสุด;
- และสามารถ **scale up**โดยปราศจากการเปลี่ยนแปลงของ เครื่องมือ สถาปัตยกรรม หรือแนวทางปฏิบัตของการพัฒนา

หลักการ twelve-factor สามารถประยุกต์ใช้ได้กับแอพพลิเคชันที่เขียนด้วยภาษาใดๆ และซึ่งใช้ร่วมกับบริการสนับสนุนใดๆ (ฐานข้อมูล, คิว, หน่วยความจำเคช เป็นต้น).

# ประวัติ

ผู้มีส่วนร่วมของเอกสารนี้ได้มีส่วนเกี่ยวข้องโดยตรงกับการพัฒนาและการใช้งานแอพพลิเคชันจำนวนมาก และเกี่ยวข้องทางอ้อมสำหรับการพัฒนา การดำเนินงาน และการขยายขนาดของแอพพลิเคชันจำนวมมหาศาลผ่านงานของเราที่แพลตฟอร์ม Heroku

เอกสารนี้สังเคราะห์จากประสบการณ์และการสังเกตทั้งหมดของพวกเราบนแอพพลิเคชัน software-as-a-service ที่หลากหลายจำนวนมาก เป็นสามเหลียมของแนวทางปฏิบัตในอุดมคติสำหรับการพัฒนาแอพพลิเคชัน ให้ความสนใจเป็นพิเศษกับพลวัตของการเจริญเติบโตของแอพพลิเคชันในช่วงเวลาหนึ่ง พลวัตของการมีส่วนร่วมระหว่างนักพัฒนาที่ทำงานกับ codebase ของแอพพลิเคชัน และหลีกเลี่ยงการใช้จ่ายของซอฟต์แวร์.

แรงจูงใจของเราเพื่อสร้างความตระหนักของปัญหาระบบบางอย่างที่เราเห็นในการพัฒนาแอพพลิเคชันสมัยใหม่ เพื่อให้คำศัพย์ที่ใช้ร่วมกันสำหรับการพูดคุยเกี่ยวกับปัญหาเหล่านี้ และนำเสนอแนวทางแก้ไขแนวกว้างสำหรับปัญหาเหล่านี้พร้อมกับคำศัพท์ที่ใช้ประกอบกัน รูปแบบนี้ได้รับแรงบันดาลใจจากหนังสือของ Martin Fowler *Patterns of Enterprise Application Architecture* และ *Refactoring*.

# Giriş

Modern çağda yazılımlar çoğunlukla “web uygulaması” ya da “yazılım hizmeti” olarak isimlendirilen servisler olarak sunulurlar. *On iki faktörlü uygulama*, servis olarak çalışan yazılımlar (İng. *software as a service* veya *SaaS*) geliştirmek için bir yöntembilimdir. Bu yöntembilimin kuralları ve faydaları şunlardır:

- Projenin kurulum otomasyonu için **açıklayıcı**(İng. declarative) biçimler kullanır. Bu şekilde, projeye yeni katılan geliştiricilerin geliştirmeye başlama zamanını ve maliyetinı en aza indirir;
- Üzerinde çalıştığı işletim sistemi ile arasında **basit bir bağlılık**vardır. Bu, tüm çalıştırma ortamlarına (Docker gibi*container*sistemleri ve sıradan işletim sistemlerine)**maksimum uyumluluk**sağlar;
- Sunucu ve sistem yönetimine olan ihtiyacı ortadan kaldıran modern **bulut platformlarına kurulum**için uygundur;
- Geliştirme ve canlı yayın ortamları arasında **farklılıklaşmayı minimize ederek**, maksimum çeviklik ile**sürekli dağıtımı**n önünü açar;
- Kullanılan araçlarda, mimaride veya geliştirme pratiklerinde önemli değişikliklere gerek duymadan **ölçeklenebilir**.

On iki faktör uygulaması herhangi bir programlama dili ile yazılmış ve yardımcı (veritabanları, kuyruk işleyiciler, önbellek, vb. gibi) servislerin herhangi bir kombinasyonuna sahip tüm uygulamalara uygulanabilir.

# Вступ

У наш час програмне забезпечення зазвичай поставляється у вигляді сервісів, що називаються *веб-застосунки* (web-apps) або *software-as-a-service* (SaaS). Застосунок дванадцяти факторів — це методологія для створення SaaS-застосунків, які:

- Використовують **декларативний**формат для автоматизації встановлення та налаштування, що зводить до мінімуму витрати часу і коштів для нових розробників, що приєднуються до проекту;
- Мають **угоду**з операційною системою, пропонуючи**максимальну переносимість**між середовищами виконання;
- Придатні для **розгортання**на сучасних**хмарних платформах**, що усуває необхідність у серверах та їх системному адмініструванні;
- **Мінімізують різницю**між середовищем розробки і production середовищем, що дозволяє- **безперервне розгортання**(continuous deployment) для забезпечення максимальної спритності розробки (agility);
- Можуть **масштабуватися**без значних змін в інструментах, архітектурі і практиці розробки.

Методологію дванадцяти факторів можна використати для застосунків, що написані будь-якою мовою програмування та використовують будь-яку комбінацію із сторонніх служб (бази даних, черги, кеш-пам’ять тощо).

# Передумови

Люди, що працювали над цим документом, брали безпосередню участь в розробці і розгортанні сотень застосунків, і мимоволі стали свідками розвитку, експлуатації та масштабування сотень тисяч застосунків під час нашої роботи над платформою Heroku.

В цьому документі узагальнюється весь наш досвід використання і спостереження за найрізноманітнішими SaaS-застосунками “в дикій природі”. Документ об’єднує ідеальні практики розробки застосунків, особлива увага приділяється динаміці органічного росту застосунку з плином часу, взаємодії між розробниками, які працюють над кодом застосунку, та уникненню витрат при ерозії програмного забезпечення.

Наша мета полягає в тому, щоб підвищити обізнаність про деякі системні проблеми, які ми бачили в практиці розробки сучасних застосунків, а також в тому, щоб сформулювати спільні загальні поняття для обговорення цих проблем, і запропонувати набір загальних концептуальних рішень цих проблем з супутньою термінологією. Формат навіяний книгами Мартіна Фаулера (Martin Fowler) *Patterns of Enterprise Application Architecture* та *Refactoring*.

# Кому слід читати цей документ?

Розробникам, які створюють SaaS-застосунки. Ops-інженерам, які виконують розгортання і керування такими застосунками.

# Giới thiệu

Ngày nay, phần mềm thường được chuyển giao như là một dịch vụ: còn được gọi là *các ứng dụng web*, hay *phần mềm-như-một-dịch vụ (software-as-a-service)*. Ứng dụng 12-hệ số là một phương pháp để xây dựng các ứng dụng phần mềm-như-một-dịch vụ với các tiêu chí sau:

- Sử dụng các định dạng theo kiểu **tường thuật**cho việc thiết lập tự động hoá, để cắt giảm chi phí và thời gian cho lập trình viên mới tham gia dự án;
- Có một **hợp đồng sạch**với hệ điều hành bên dưới, cung cấp**tối đa khả năng dịch chuyển**giữa các môi trường thực thi;
- Phù hợp để **triển khai**trên các**nền tảng đám mây**mới, cắt giảm yêu cầu quản trị cho server và hệ thống;
- **Giảm thiểu sự khác nhau**giữa môi trường phát triển và môi trường sản xuất, cho phép đạt được sự linh hoạt tối đa trong- **triển khai liên tục**;
- Và có thể **mở rộng**mà không cần thay đổi lớn cho các công cụ, kiến trúc, hoặc cách thức phát triển.

Phương pháp 12-hệ số có thể được áp dụng cho các ứng dụng viết bằng bất kì ngôn ngữ lập trình nào, và sử dụng bất kì kết hợp giữa các dịch vụ backend (cơ sở dữ liệu, queue, memory cache, vv.).

# Gốc gác

Tất cả tác giả của tài liệu này đã trực tiếp tham gia vào quá trình phát triển và triển khai của hàng trăm ứng dụng, và gián tiếp theo dõi các quá trình phát triển, vận hành, và mở rộng của hàng nghìn ứng dụng thông qua công việc của chúng tôi trên hệ thống Heroku.

Tài liệu này là cô đọng của tất cả kinh nghiệm và quan sát của chúng tôi trên một số lượng lớn các ứng dụng-như-một-dịch vụ ở ngoài. Đây là kết hợp của kiến thức thực hành chuẩn mực trong việc phát triển ứng dụng, với trọng tâm vào cơ cấu phát triển cơ bản của ứng dụng trong một khoảng thời gian, cơ cấu động của sự hợp tác giữa các lập trình viên đang làm việc trên cùng một mã gốc, và tránh rò rỉ chi phí phát triển phần mềm.

Động lực của chúng tôi là tăng cường nhận thức về các vấn đề hệ thống mà chúng tôi biết với các qui trình phát triển ứng dụng hiện tại, để chia sẻ một kho kiến thức thảo luận về các vấn đề này, và để cung cấp một chuỗi các giải pháp mở chác vấn đề trên và cũng đi kèm với các thuật ngữ chuyên môn. Định dạng này lấy ý tưởng từ cuốn sách *Patterns of Enterprise Application Architecture* và *Refactoring* của ông Martin Fowler.

# 简介

如今，软件通常会作为一种服务来交付，它们被称为网络应用程序，或软件即服务（SaaS）。12-Factor 为构建如下的 SaaS 应用提供了方法论：

- 使用**标准化**流程自动配置，从而使新的开发者花费最少的学习成本加入这个项目。
- 和操作系统之间尽可能的**划清界限**，在各个系统中提供**最大的可移植性**。
- 适合**部署**在现代的**云计算平台**，从而在服务器和系统管理方面节省资源。
- 将开发环境和生产环境的**差异降至最低**，并使用**持续交付**实施敏捷开发。
- 可以在工具、架构和开发流程不发生明显变化的前提下实现**扩展**。

这套理论适用于任意语言和后端服务（数据库、消息队列、缓存等）开发的应用程序。

# 背景

本文的贡献者参与过数以百计的应用程序的开发和部署，并通过 Heroku 平台间接见证了数十万应用程序的开发，运作以及扩展的过程。

本文综合了我们关于 SaaS 应用几乎所有的经验和智慧，是开发此类应用的理想实践标准，并特别关注于应用程序如何保持良性成长，开发者之间如何进行有效的代码协作，以及如何 避免软件污染 。

我们的初衷是分享在现代软件开发过程中发现的一些系统性问题，并加深对这些问题的认识。我们提供了讨论这些问题时所需的共享词汇，同时使用相关术语给出一套针对这些问题的广义解决方案。本文格式的灵感来自于 Martin Fowler 的书籍： *Patterns of Enterprise Application Architecture* ， *Refactoring* 。
