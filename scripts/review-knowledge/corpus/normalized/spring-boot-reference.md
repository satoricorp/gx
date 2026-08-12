Spring Boot Reference Reference This section provides information on using the features and capabilities of Spring Boot. Developing Your First Spring Boot Application Developing with Spring Boot

Spring Boot Reference Reference This section provides information on using the features and capabilities of Spring Boot. Developing Your First Spring Boot Application Developing with Spring Boot

# Developing with Spring Boot

This section goes into more detail about how you should use Spring Boot. It covers topics such as build systems, auto-configuration, and how to run your applications. We also cover some Spring Boot best practices. Although there is nothing particularly special about Spring Boot (it is just another library that you can consume), there are a few recommendations that, when followed, make your development process a little easier.

If you are starting out with Spring Boot, you should probably read the Developing Your First Spring Boot Application tutorial before diving into this section.

# Build Systems

It is strongly recommended that you choose a build system that supports dependency management and that can consume artifacts published to the Maven Central repository. We would recommend that you choose Maven or Gradle. It is possible to get Spring Boot to work with other build systems (Ant, for example), but they are not particularly well supported.

## Dependency Management

Each release of Spring Boot provides a curated list of dependencies that it supports. In practice, you do not need to provide a version for any of these dependencies in your build configuration, as Spring Boot manages that for you. When you upgrade Spring Boot itself, these dependencies are upgraded as well in a consistent way.

| You can still specify a version and override Spring Boot’s recommendations if you need to do so. |

The curated list contains all the Spring modules that you can use with Spring Boot as well as a refined list of third party libraries.
The list is available as a standard Bills of Materials (`spring-boot-dependencies`) that can be used with both Maven and Gradle.

| Each release of Spring Boot is associated with a base version of the Spring Framework.
We highlyrecommend that you do not specify its version. |

## Ant

It is possible to build a Spring Boot project using Apache Ant+Ivy.
The `spring-boot-antlib` “AntLib” module is also available to help Ant create executable jars.

To declare dependencies, a typical `ivy.xml` file looks something like the following example:

```
<ivy-module version="2.0">
	<info organisation="org.springframework.boot" module="spring-boot-sample-ant" />
	<configurations>
 <conf name="compile" description="everything needed to compile this module" />
 <conf name="runtime" extends="compile" description="everything needed to run this module" />
	</configurations>
	<dependencies>
 <dependency org="org.springframework.boot" name="spring-boot-starter"
 rev="${spring-boot.version}" conf="compile" />
	</dependencies>
</ivy-module>
```
A typical `build.xml` looks like the following example:

```
<project
	xmlns:ivy="antlib:org.apache.ivy.ant"
	xmlns:spring-boot="antlib:org.springframework.boot.ant"
	name="myapp" default="build">
	<property name="spring-boot.version" value="4.1.0" />
	<target name="resolve" description="--> retrieve dependencies with ivy">
 <ivy:retrieve pattern="lib/[conf]/[artifact]-[type]-[revision].[ext]" />
	</target>
	<target name="classpaths" depends="resolve">
 <path id="compile.classpath">
 <fileset dir="lib/compile" includes="*.jar" />
 </path>
	</target>
	<target name="init" depends="classpaths">
 <mkdir dir="build/classes" />
	</target>
	<target name="compile" depends="init" description="compile">
 <javac srcdir="src/main/java" destdir="build/classes" classpathref="compile.classpath" />
	</target>
	<target name="build" depends="compile">
 <spring-boot:exejar destfile="build/myapp.jar" classes="build/classes">
 <spring-boot:lib>
 <fileset dir="lib/runtime" />
 </spring-boot:lib>
 </spring-boot:exejar>
	</target>
</project>
```
| If you do not want to use the `spring-boot-antlib`module, see the Build an Executable Archive From Ant without Using spring-boot-antlib section of “How-to Guides”. |

## Starters

Starters are a set of convenient dependency descriptors that you can include in your application.
You get a one-stop shop for all the Spring and related technologies that you need without having to hunt through sample code and copy-paste loads of dependency descriptors.
For example, if you want to get started using Spring and JPA for database access, include the `spring-boot-starter-data-jpa` dependency in your project.

The starters contain a lot of the dependencies that you need to get a project up and running quickly and with a consistent, supported set of managed transitive dependencies.

The following application starters are provided by Spring Boot under the `org.springframework.boot` group:

| Name | Description |
|---|---|
| Core starter, including auto-configuration support, logging and YAML | |
| Starter for using Apache ActiveMQ and JMS | |
| Starter for testing using Apache ActiveMQ and JMS | |
| Starter for testing Spring Boot’s Actuator which provides production ready features to help you monitor and manage your application | |
| Starter for using Spring AMQP and Rabbit MQ | |
| Starter for testing Spring AMQP and Rabbit MQ | |
| Starter for using Apache Artemis and JMS | |
| Starter for testing Apache Artemis and JMS | |
| Starter for using aspect-oriented programming with AspectJ | |
| Starter for testing aspect-oriented programming with AspectJ | |
| Starter for using Spring Batch | |
| Starter for using Spring Batch with Data MongoDB | |
| Starter for testing using Spring Batch with Data MongoDB | |
| Starter for using Spring Batch with JDBC | |
| Starter for testing using Spring Batch with JDBC | |
| Starter for testing using Spring Batch | |
| Starter for using Spring’s caching support | |
| Starter for testing Spring’s caching support | |
| Starter for using Cassandra distributed database | |
| Starter for testing Cassandra distributed database | |
| Core classic starter, including full auto-configuration support, logging and YAML | |
| Starter for using Cloud Foundry | |
| Starter for testing Cloud Foundry | |
| Starter for using Couchbase document-oriented database | |
| Starter for testing Couchbase document-oriented database | |
| Starter for using Cassandra distributed database and Spring Data Cassandra | |
| Starter for using Cassandra distributed database and Spring Data Cassandra Reactive | |
| Starter for testing Cassandra distributed database and Spring Data Cassandra Reactive | |
| Starter for testing Cassandra distributed database and Spring Data Cassandra | |
| Starter for using Couchbase document-oriented database and Spring Data Couchbase | |
| Starter for using Couchbase document-oriented database and Spring Data Couchbase Reactive | |
| Starter for testing Couchbase document-oriented database and Spring Data Couchbase Reactive | |
| Starter for testing Couchbase document-oriented database and Spring Data Couchbase | |
| Starter for using Elasticsearch search and analytics engine and Spring Data Elasticsearch | |
| Starter for testing Elasticsearch search and analytics engine and Spring Data Elasticsearch | |
| Starter for using Spring Data JDBC | |
| Starter for testing Spring Data JDBC | |
| Starter for using Spring Data JPA with Hibernate | |
| Starter for testing Spring Data JPA with Hibernate | |
| Starter for using Spring Data LDAP | |
| Starter for testing Spring Data LDAP | |
| Starter for using MongoDB document-oriented database and Spring Data MongoDB | |
| Starter for using MongoDB document-oriented database and Spring Data MongoDB Reactive | |
| Starter for using MongoDB document-oriented database and Spring Data MongoDB Reactive | |
| Starter for testing MongoDB document-oriented database and Spring Data MongoDB | |
| Starter for using Neo4j graph database and Spring Data Neo4j | |
| Starter for testing Neo4j graph database and Spring Data Neo4j | |
| Starter for using Spring Data R2DBC | |
| Starter for testing Spring Data R2DBC | |
| Starter for using Redis key-value data store with Spring Data Redis and the Lettuce client | |
| Starter for using Redis key-value data store with Spring Data Redis reactive and the Lettuce client | |
| Starter for testing Redis key-value data store with Spring Data Redis reactive and the Lettuce client | |
| Starter for testing Redis key-value data store with Spring Data Redis and the Lettuce client | |
| Starter for using Spring Data repositories exposed over REST using Spring Data REST and Spring MVC | |
| Starter for testing Spring Data repositories exposed over REST using Spring Data REST and Spring MVC | |
| Starter for using Elasticsearch search and analytics engine | |
| Starter for testing Elasticsearch search and analytics engine | |
| Starter for using Flyway database migrations | |
| Starter for testing Flyway database migrations | |
| Starter for using FreeMarker | |
| Starter for testing FreeMarker | |
| Starter using Spring GraphQL | |
| Starter for testing Spring GraphQL | |
| Starter for using Groovy Templates | |
| Starter for testing Groovy Templates | |
| Starter for using Spring gRPC client | |
| Starter for testing gRPC client | |
| Starter for using Spring gRPC server | |
| Starter for testing gRPC server | |
| Starter for using GSON | |
| Starter for testing GSON | |
| Starter for using Spring HATEOS to build hypermedia-based RESTful Spring MVC web applications | |
| Starter for testing Spring HATEOS to build hypermedia-based RESTful Spring MVC web applications | |
| Starter for using Hazelcast | |
| Starter for testing Hazelcast | |
| Starter for using Spring Integration | |
| Starter for testing Spring Integration | |
| Starter for using Jackson | |
| Starter for testing Jackson | |
| Starter for using JDBC with the HikariCP connection pool | |
| Starter for testing JDBC with the HikariCP connection pool | |
| Starter for using JAX-RS and Jersey | |
| Starter for testing JAX-RS and Jersey | |
| Starter for using Jetty as the embedded servlet container | |
| Starter for using JMS | |
| Starter for testing JMS | |
| Starter for using jOOQ to access SQL databases with JDBC | |
| Starter for testing jOOQ to access SQL databases with JDBC | |
| Starter for reading and writing JSON | |
| Starter for using JSON-B | |
| Starter for testing JSON-B | |
| Starter for using Apache Kafka | |
| Starter for testing Apache Kafka | |
| Starter for using Kotlinx Serialization JSON | |
| Starter for testing Kotlinx Serialization JSON | |
| Starter for using LDAP | |
| Starter for testing LDAP | |
| Starter for using Liquibase database migrations | |
| Starter for testing Liquibase database migrations | |
| Starter for using Java Mail and Spring Framework’s email sending support | |
| Starter for testing Java Mail and Spring Framework’s email sending support | |
| Starter for using Micrometer Metrics | |
| Starter for testing Micrometer Metrics | |
| Starter for using MongoDB document-oriented database | |
| Starter for testing MongoDB document-oriented database | |
| Starter for using Mustache | |
| Starter for testing Mustache | |
| Starter for using Neo4j graph database | |
| Starter for testing Neo4j graph database | |
| Starter for using Spring Authorization Server features (deprecated in favor of | |
| Starter for using Spring Security’s OAuth2/OpenID Connect client features (deprecated in favor of | |
| Starter for using Spring Security’s OAuth2 resource server features (deprecated in favor of | |
| Starter for using OpenTelemetry | |
| Starter for testing OpenTelemetry | |
| Starter for using Spring for Apache Pulsar | |
| Starter for testing Spring for Apache Pulsar | |
| Starter for using the Quartz scheduler | |
| Starter for testing the Quartz scheduler | |
| Starter for using R2DBC | |
| Starter for testing R2DBC | |
| Starter for Reactor Netty | |
| Starter using Spring’s blocking HTTP clients (RestClient, RestTemplate and HTTP Service Clients) | |
| Starter for testing Spring’s blocking HTTP clients (RestClient, RestTemplate and HTTP Service Clients) | |
| Starter for using RSocket | |
| Starter for testing RSocket | |
| Starter for using Spring Security | |
| Starter for using Spring Authorization Server features | |
|
 | Starter for testing Spring Authorization Server features |
| Starter for using Spring Security’s OAuth2/OpenID Connect client features | |
| Starter for testing Spring Security’s OAuth2/OpenID Connect client features | |
| Starter for using Spring Security’s OAuth2 resource server features | |
| Starter for testing Spring Security’s OAuth2 resource server features | |
| Starter for using Spring Security with SAML2 | |
| Starter for testing Spring Security with SAML2 | |
| Starter for testing Spring Security | |
| Starter for using Spring Session with Sendgrid | |
| Starter for testing Spring Session with Sendgrid | |
| Starter for using Spring Session with Spring Data Redis | |
| Starter for testing Spring Session with Spring Data Redis | |
| Starter for using Spring Session with JDBC | |
| Starter for testing Spring Session with JDBC | |
| Starter for testing Spring Boot applications with libraries including JUnit Jupiter, Hamcrest and Mockito | |
| Classic starter for testing Spring Boot applications with libraries including JUnit Jupiter, Hamcrest and Mockito | |
| Starter for using Thymeleaf | |
| Starter for testing Thymeleaf | |
| Starter for using Tomcat as the embedded servlet container | |
| Starter for using Java Bean Validation with Hibernate Validator | |
| Starter for testing Java Bean Validation with Hibernate Validator | |
| Starter for building web, including RESTful, applications using Spring MVC. Uses Tomcat as the default embedded container (deprecated in favor of | |
| Starter for testing Spring Web Server | |
| Starter for using Spring Web Services (deprecated in favor of | |
| Starter using Spring’s reactive HTTP clients (WebClient and HTTP Service Clients) | |
| Starter for testing Spring’s reactive HTTP clients (WebClient and HTTP Service Clients) | |
| Starter for using WebFlux and Reactor Netty | |
| Starter for testing WebFlux and Reactor Netty | |
| Starter for using Spring MVC and Tomcat | |
| Starter for testing Spring MVC and Tomcat | |
| Starter for using Spring Web Services | |
| Starter for testing Spring Web Services | |
| Starter for using Spring MVC WebSocket support | |
| Starter for testing Spring MVC WebSocket support | |
| Starter for using Zipkin | |
| Starter for testing Zipkin |

In addition to the application starters, the following starters can be used to add production ready features:

| Name | Description |
|---|---|
| Starter for using Spring Boot’s Actuator which provides production ready features to help you monitor and manage your application |

Finally, Spring Boot also includes the following starters that can be used if you want to exclude or swap specific technical facets:

| Name | Description |
|---|---|
| Starter for the Jetty runtime | |
| Starter for using Log4j2 | |
| Starter for logging using Logback | |
| Starter for logging default logging | |
| Starter for using Spring REST Docs | |
| Starter for the Tomcat runtime |

To learn how to swap technical facets, please see the how-to documentation for swapping web server and logging system.

| For a list of additional community contributed starters, see the README file in the `spring-boot-starters`module on GitHub. |

# Structuring Your Code

Spring Boot does not require any specific code layout to work. However, there are some best practices that help.

| If you wish to enforce a structure based on domains, take a look at Spring Modulith. |

## Using the “default” Package

When a class does not include a `package` declaration, it is considered to be in the “default package”.
The use of the “default package” is generally discouraged and should be avoided.
It can cause particular problems for Spring Boot applications that use the `@ComponentScan`, `@ConfigurationPropertiesScan`, `@EntityScan`, or `@SpringBootApplication` annotations, since every class from every jar is read.

| We recommend that you follow Java’s recommended package naming conventions and use a reversed domain name (for example, `com.example.project`). |

## Locating the Main Application Class

We generally recommend that you locate your main application class in a root package above other classes.
The `@SpringBootApplication` annotation is often placed on your main class, and it implicitly defines a base “search package” for certain items.
For example, if you are writing a JPA application, the package of the `@SpringBootApplication` annotated class is used to search for `@Entity` items.
Using a root package also allows component scan to apply only on your project.

| If you do not want to use `@SpringBootApplication`, the`@EnableAutoConfiguration`and`@ComponentScan`annotations that it imports defines that behavior so you can also use those instead. |

The following listing shows a typical layout:

```
com
 +- example
 +- myapplication
 +- MyApplication.java
 |
 +- customer
 | +- Customer.java
 | +- CustomerController.java
 | +- CustomerService.java
 | +- CustomerRepository.java
 |
 +- order
 +- Order.java
 +- OrderController.java
 +- OrderService.java
 +- OrderRepository.java
```
The `MyApplication.java` file would declare the `main` method, along with the basic `@SpringBootApplication`, as follows:

-
Java
-
Kotlin

```
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
@SpringBootApplication
public class MyApplication {
	public static void main(String[] args) {
 SpringApplication.run(MyApplication.class, args);
	}
}
```
```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
@SpringBootApplication
class MyApplication
fun main(args: Array<String>) {
	runApplication<MyApplication>(*args)
}
```

# Configuration Classes

Spring Boot favors Java-based configuration.
Although it is possible to use `SpringApplication` with XML sources, we generally recommend that your primary source be a single `@Configuration` class.
Usually the class that defines the `main` method is a good candidate as the primary `@Configuration`.

| Many Spring configuration examples have been published on the Internet that use XML configuration.
If possible, always try to use the equivalent Java-based configuration.
Searching for `Enable*`annotations can be a good starting point. |

## Importing Additional Configuration Classes

You need not put all your `@Configuration` into a single class.
The `@Import` annotation can be used to import additional configuration classes.
Alternatively, you can use `@ComponentScan` to automatically pick up all Spring components, including `@Configuration` classes.

## Importing XML Configuration

If you absolutely must use XML based configuration, we recommend that you still start with a `@Configuration` class.
You can then use an `@ImportResource` annotation to load XML configuration files.

# Auto-configuration

Spring Boot auto-configuration attempts to automatically configure your Spring application based on the jar dependencies that you have added.
For example, if `HSQLDB` is on your classpath, and you have not manually configured any database connection beans, then Spring Boot auto-configures an in-memory database.

You need to opt-in to auto-configuration by adding the `@EnableAutoConfiguration` or `@SpringBootApplication` annotations to one of your `@Configuration` classes.

| You should only ever add one `@SpringBootApplication`or`@EnableAutoConfiguration`annotation.
We generally recommend that you add one or the other to your primary`@Configuration`class only. |

## Gradually Replacing Auto-configuration

Auto-configuration is non-invasive.
At any point, you can start to define your own configuration to replace specific parts of the auto-configuration.
For example, if you add your own `DataSource` bean, the default embedded database support backs away.

If you need to find out what auto-configuration is currently being applied, and why, start your application with the `--debug` switch.
Doing so enables debug logs for a selection of core loggers and logs a conditions report to the console.

## Disabling Specific Auto-configuration Classes

If you find that specific auto-configuration classes that you do not want are being applied, you can use the exclude attribute of `@SpringBootApplication` to disable them, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.jdbc.autoconfigure.DataSourceAutoConfiguration;
@SpringBootApplication(exclude = { DataSourceAutoConfiguration.class })
public class MyApplication {
}
```
```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.jdbc.autoconfigure.DataSourceAutoConfiguration
@SpringBootApplication(exclude = [DataSourceAutoConfiguration::class])
class MyApplication
```
If the class is not on the classpath, you can use the `excludeName` attribute of the annotation and specify the fully qualified name instead.
If you prefer to use `@EnableAutoConfiguration` rather than `@SpringBootApplication`, `exclude` and `excludeName` are also available.
Finally, you can also control the list of auto-configuration classes to exclude by using the `spring.autoconfigure.exclude` property.

| You can define exclusions both at the annotation level and by using the property. |

| Even though auto-configuration classes are `public`, the only aspect of the class that is considered public API is the name of the class which can be used for disabling the auto-configuration.
The actual contents of those classes, such as nested configuration classes or bean methods are for internal use only and we do not recommend using those directly. |

## Auto-configuration Packages

Auto-configuration packages are the packages that various auto-configured features look in by default when scanning for things such as entities and Spring Data repositories.
The `@EnableAutoConfiguration` annotation (either directly or through its presence on `@SpringBootApplication`) determines the default auto-configuration package.
Additional packages can be configured using the `@AutoConfigurationPackage` annotation.

# Spring Beans and Dependency Injection

You are free to use any of the standard Spring Framework techniques to define your beans and their injected dependencies.
We generally recommend using constructor injection to wire up dependencies and `@ComponentScan` to find beans.

If you structure your code as suggested above (locating your application class in a top package), you can add `@ComponentScan` without any arguments or use the `@SpringBootApplication` annotation which implicitly includes it.
All of your application components (`@Component`, `@Service`, `@Repository`, `@Controller`, and others) are automatically registered as Spring Beans.

The following example shows a `@Service` Bean that uses constructor injection to obtain a required `RiskAssessor` bean:

-
Java
-
Kotlin

```
import org.springframework.stereotype.Service;
@Service
public class MyAccountService implements AccountService {
	private final RiskAssessor riskAssessor;
	public MyAccountService(RiskAssessor riskAssessor) {
 this.riskAssessor = riskAssessor;
	}
	// ...
}
```
```
import org.springframework.stereotype.Service
@Service
class MyAccountService(private val riskAssessor: RiskAssessor) : AccountService
```
If a bean has more than one constructor, you will need to mark the one you want Spring to use with `@Autowired`:

-
Java
-
Kotlin

```
import java.io.PrintStream;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
@Service
public class MyAccountService implements AccountService {
	private final RiskAssessor riskAssessor;
	private final PrintStream out;
	@Autowired
	public MyAccountService(RiskAssessor riskAssessor) {
 this.riskAssessor = riskAssessor;
 this.out = System.out;
	}
	public MyAccountService(RiskAssessor riskAssessor, PrintStream out) {
 this.riskAssessor = riskAssessor;
 this.out = out;
	}
	// ...
}
```
```
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.stereotype.Service
import java.io.PrintStream
@Service
class MyAccountService : AccountService {
	private val riskAssessor: RiskAssessor
	private val out: PrintStream
	@Autowired
	constructor(riskAssessor: RiskAssessor) {
 this.riskAssessor = riskAssessor
 out = System.out
	}
	constructor(riskAssessor: RiskAssessor, out: PrintStream) {
 this.riskAssessor = riskAssessor
 this.out = out
	}
	// ...
}
```
| Notice how using constructor injection lets the `riskAssessor`field be marked as`final`, indicating that it cannot be subsequently changed. |

# Using the @SpringBootApplication Annotation

Many Spring Boot developers like their apps to use auto-configuration, component scan and be able to define extra configuration on their "application class".
A single `@SpringBootApplication` annotation can be used to enable those three features, that is:

-
`@EnableAutoConfiguration`: enable Spring Boot’s auto-configuration mechanism
-
`@ComponentScan`: enable`@Component`scan on the package where the application is located (see the best practices)
-
`@SpringBootConfiguration`: enable registration of extra beans in the context or the import of additional configuration classes. An alternative to Spring’s standard`@Configuration`that aids configuration detection in your integration tests.

-
Java
-
Kotlin

```
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
// Same as @SpringBootConfiguration @EnableAutoConfiguration @ComponentScan
@SpringBootApplication
public class MyApplication {
	public static void main(String[] args) {
 SpringApplication.run(MyApplication.class, args);
	}
}
```
```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
// Same as @SpringBootConfiguration @EnableAutoConfiguration @ComponentScan
@SpringBootApplication
class MyApplication
fun main(args: Array<String>) {
	runApplication<MyApplication>(*args)
}
```
| `@SpringBootApplication`also provides aliases to customize the attributes of`@EnableAutoConfiguration`and`@ComponentScan`. |

| None of these features are mandatory and you may choose to replace this single annotation by any of the features that it enables. For instance, you may not want to use component scan or configuration properties scan in your application:
 In this example, |

# Running Your Application

One of the biggest advantages of packaging your application as a jar and using an embedded HTTP server is that you can run your application as you would any other. The same applies to debugging Spring Boot applications. You do not need any special IDE plugins or extensions.

| The options below are best suited for running an application locally for development. For production deployment, see Packaging Your Application for Production. |

| This section only covers jar-based packaging. If you choose to package your application as a war file, see your server and IDE documentation. |

## Running From an IDE

You can run a Spring Boot application from your IDE as a Java application.
However, you first need to import your project.
Import steps vary depending on your IDE and build system.
Most IDEs can import Maven projects directly.
For example, Eclipse users can select `Import…` → `Existing Maven Projects` from the `File` menu.

If you cannot directly import your project into your IDE, you may be able to generate IDE metadata by using a build plugin. Maven includes plugins for Eclipse and IntelliJ IDEA. Gradle offers plugins for various IDEs.

| If you accidentally run a web application twice, you see a “Port already in use” error.
Spring Tools users can use the `Relaunch`button rather than the`Run`button to ensure that any existing instance is closed. |

## Running as a Packaged Application

If you use the Spring Boot Maven or Gradle plugins to create an executable jar, you can run your application using `java -jar`, as shown in the following example:

`$ java -jar target/myapplication-0.0.1-SNAPSHOT.jar`It is also possible to run a packaged application with remote debugging support enabled. Doing so lets you attach a debugger to your packaged application, as shown in the following example:

```
$ java -agentlib:jdwp=server=y,transport=dt_socket,address=8000,suspend=n \
 -jar target/myapplication-0.0.1-SNAPSHOT.jar
```
## Using the Maven Plugin

The Spring Boot Maven plugin includes a `run` goal that can be used to quickly compile and run your application.
Applications run in an exploded form, as they do in your IDE.
The following example shows a typical Maven command to run a Spring Boot application:

`$ mvn spring-boot:run`You might also want to use the `MAVEN_OPTS` operating system environment variable, as shown in the following example:

`$ export MAVEN_OPTS=-Xmx1024m`## Using the Gradle Plugin

The Spring Boot Gradle plugin also includes a `bootRun` task that can be used to run your application in an exploded form.
The `bootRun` task is added whenever you apply the `org.springframework.boot` and `java` plugins and is shown in the following example:

`$ gradle bootRun`You might also want to use the `JAVA_OPTS` operating system environment variable, as shown in the following example:

`$ export JAVA_OPTS=-Xmx1024m`## Hot Swapping

Since Spring Boot applications are plain Java applications, JVM hot-swapping should work out of the box. JVM hot swapping is somewhat limited with the bytecode that it can replace. For a more complete solution, JRebel can be used.

The `spring-boot-devtools` module also includes support for quick application restarts.
See the Hot Swapping section in “How-to Guides” for details.

# Developer Tools

Spring Boot includes an additional set of tools that can make the application development experience a little more pleasant.
The `spring-boot-devtools` module can be included in any project to provide additional development-time features.
To include devtools support, add the module dependency to your build, as shown in the following listings for Maven and Gradle:

```
<dependencies>
	<dependency>
 <groupId>org.springframework.boot</groupId>
 <artifactId>spring-boot-devtools</artifactId>
 <optional>true</optional>
	</dependency>
</dependencies>
```
```
dependencies {
	developmentOnly("org.springframework.boot:spring-boot-devtools")
}
```
| Devtools might cause classloading issues, in particular in multi-module projects. Diagnosing Classloading Issues explains how to diagnose and solve them. |

| Developer tools are automatically disabled when running a fully packaged application.
If your application is launched from `java -jar`or if it is started from a special classloader, then it is considered a “production application”.
You can control this behavior by using the`spring.devtools.restart.enabled`system property.
To enable devtools, irrespective of the classloader used to launch your application, set the`-Dspring.devtools.restart.enabled=true`system property.
This must not be done in a production environment where running devtools is a security risk.
To disable devtools, exclude the dependency or set the`-Dspring.devtools.restart.enabled=false`system property. |

| Flagging the dependency as optional in Maven or using the `developmentOnly`configuration in Gradle (as shown above) prevents devtools from being transitively applied to other modules that use your project. |

| Repackaged archives do not contain devtools by default.
If you want to use a certain remote devtools feature, you need to include it.
When using the Maven plugin, opt-in for optional dependencies by setting the `includeOptional`property to`true`.
You also need to set the`excludeDevtools`property to`false`.
When using the Gradle plugin, configure the task’s classpath to include the`developmentOnly`configuration. |

## Diagnosing Classloading Issues

As described in the Restart vs Reload section, restart functionality is implemented by using two classloaders. For most applications, this approach works well. However, it can sometimes cause classloading issues, in particular in multi-module projects.

To diagnose whether the classloading issues are indeed caused by devtools and its two classloaders, try disabling restart. If this solves your problems, customize the restart classloader to include your entire project.

## Property Defaults

Several of the libraries supported by Spring Boot use caches to improve performance. For example, template engines cache compiled templates to avoid repeatedly parsing template files. Also, Spring MVC can add HTTP caching headers to responses when serving static resources.

While caching is very beneficial in production, it can be counter-productive during development, preventing you from seeing the changes you just made in your application. For this reason, spring-boot-devtools disables the caching options by default.

Cache options are usually configured by settings in your `application.properties` file.
For example, Thymeleaf offers the `spring.thymeleaf.cache` property.

The same applies for tracing probability that’s set to 100% as the default may not log all traces used for testing.

Rather than needing to set these properties manually, the `spring-boot-devtools` module automatically applies sensible development-time configuration.

The following table lists all the properties that are applied:

| Name | Default Value |
|---|---|
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |

| If you do not want property defaults to be applied you can set `spring.devtools.add-properties`to`false`in your`application.properties`. |

Because you need more information about web requests while developing Spring MVC and Spring WebFlux applications, developer tools suggests you to enable `DEBUG` logging for the `web` logging group.
This will give you information about the incoming request, which handler is processing it, the response outcome, and other details.
If you wish to log all request details (including potentially sensitive information), you can turn on the `spring.mvc.log-request-details` or `spring.http.codecs.log-request-details` configuration properties.

## Automatic Restart

Applications that use `spring-boot-devtools` automatically restart whenever files on the classpath change.
This can be a useful feature when working in an IDE, as it gives a very fast feedback loop for code changes.
By default, any entry on the classpath that points to a directory is monitored for changes.
Note that certain resources, such as static assets and view templates, do not need to restart the application.

| If you are restarting with Maven or Gradle using the build plugin you must leave the `forking`set to`enabled`.
If you disable forking, the isolated application classloader used by devtools will not be created and restarts will not operate properly. |

| If you use JRebel, automatic restarts are disabled in favor of dynamic class reloading. Other devtools features (such as property overrides) can still be used. |

| DevTools relies on the application context’s shutdown hook to close it during a restart.
It does not work correctly if you have disabled the shutdown hook ( `SpringApplication.setRegisterShutdownHook(false)`). |

| DevTools needs to customize the `ResourceLoader`used by the`ApplicationContext`.
If your application provides one already, it is going to be wrapped.
Direct override of the`getResource`method on the`ApplicationContext`is not supported. |

| Automatic restart is not supported when using AspectJ weaving. |

### Logging Changes in Condition Evaluation

By default, each time your application restarts, a report showing the condition evaluation delta is logged. The report shows the changes to your application’s auto-configuration as you make changes such as adding or removing beans and setting configuration properties.

To disable the logging of the report, set the following property:

-
Properties
-
YAML

`spring.devtools.restart.log-condition-evaluation-delta=false````
spring:
 devtools:
 restart:
 log-condition-evaluation-delta: false
```
### Excluding Resources

Certain resources do not necessarily need to trigger a restart when they are changed.
For example, Thymeleaf templates can be edited in-place.
By default, changing resources in `/META-INF/maven`, `/META-INF/resources`, `/resources`, `/static`, `/public`, or `/templates` does not trigger a restart but does trigger a live reload.
If you want to customize these exclusions, you can use the `spring.devtools.restart.exclude` property.
For example, to exclude only `/static` and `/public` you would set the following property:

-
Properties
-
YAML

`spring.devtools.restart.exclude=static/**,public/**````
spring:
 devtools:
 restart:
 exclude: "static/**,public/**"
```
| If you want to keep those defaults and addadditional exclusions, use the`spring.devtools.restart.additional-exclude`property instead. |

### Watching Additional Paths

You may want your application to be restarted or reloaded when you make changes to files that are not on the classpath.
To do so, use the `spring.devtools.restart.additional-paths` property to configure additional paths to watch for changes.
You can use the `spring.devtools.restart.exclude` property described earlier to control whether changes beneath the additional paths trigger a full restart or a live reload.

### Disabling Restart

If you do not want to use the restart feature, you can disable it by using the `spring.devtools.restart.enabled` property.
In most cases, you can set this property in your `application.properties` (doing so still initializes the restart classloader, but it does not watch for file changes).

If you need to *completely* disable restart support (for example, because it does not work with a specific library), you need to set the `spring.devtools.restart.enabled` `System` property to `false` before calling `SpringApplication.run(…)`, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
@SpringBootApplication
public class MyApplication {
	public static void main(String[] args) {
 System.setProperty("spring.devtools.restart.enabled", "false");
 SpringApplication.run(MyApplication.class, args);
	}
}
```
```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
@SpringBootApplication
class MyApplication
fun main(args: Array<String>) {
	System.setProperty("spring.devtools.restart.enabled", "false")
	runApplication<MyApplication>(*args)
}
```
### Using a Trigger File

If you work with an IDE that continuously compiles changed files, you might prefer to trigger restarts only at specific times. To do so, you can use a “trigger file”, which is a special file that must be modified when you want to actually trigger a restart check.

| Any update to the file will trigger a check, but restart only actually occurs if Devtools has detected it has something to do. |

To use a trigger file, set the `spring.devtools.restart.trigger-file` property to the name (excluding any path) of your trigger file.
The trigger file must appear somewhere on your classpath.

For example, if you have a project with the following structure:

```
src
+- main
 +- resources
 +- .reloadtrigger
```
Then your `trigger-file` property would be:

-
Properties
-
YAML

`spring.devtools.restart.trigger-file=.reloadtrigger````
spring:
 devtools:
 restart:
 trigger-file: ".reloadtrigger"
```
Restarts will now only happen when the `src/main/resources/.reloadtrigger` is updated.

| You might want to set `spring.devtools.restart.trigger-file`as a global setting, so that all your projects behave in the same way. |

Some IDEs have features that save you from needing to update your trigger file manually.
Spring Tools for Eclipse and IntelliJ IDEA (Ultimate Edition) both have such support.
With Spring Tools, you can use the “reload” button from the console view (as long as your `trigger-file` is named `.reloadtrigger`).
For IntelliJ IDEA, you can follow the instructions in their documentation.

### Customizing the Restart Classloader

As described earlier in the Restart vs Reload section, restart functionality is implemented by using two classloaders.
If this causes issues, you can diagnose the problem by using the `spring.devtools.restart.enabled` system property, and if the app works with restart switched off, you might need to customize what gets loaded by which classloader.

By default, any open project in your IDE is loaded with the “restart” classloader, and any regular `.jar` file is loaded with the “base” classloader.
The same is true if you use `mvn spring-boot:run` or `gradle bootRun`: the project containing your `@SpringBootApplication` is loaded with the “restart” classloader, and everything else with the “base” classloader.
The classpath is printed on the console when you start the app, which can help to identify any problematic entries.
Classes used reflectively, especially annotations, can be loaded into the parent (fixed) classloader on startup before the application classes which use them, and this might lead to them not being detected by Spring in the application.

You can instruct Spring Boot to load parts of your project with a different classloader by creating a `META-INF/spring-devtools.properties` file.
The `spring-devtools.properties` file can contain properties prefixed with `restart.exclude` and `restart.include`.
The `include` elements are items that should be pulled up into the “restart” classloader, and the `exclude` elements are items that should be pushed down into the “base” classloader.
The value of the property is a regex pattern that is applied to the classpath passed to the JVM on startup.
Here is an example where some local class files are excluded and some extra libraries are included in the restart class loader:

```
restart.exclude.companycommonlibs="/mycorp-common-[\\w\\d-\\.]/(build|bin|out|target)/"
restart.include.projectcommon="/mycorp-myproj-[\\w\\d-\\.]+\\.jar"
```
| All property keys must be unique.
As long as a property starts with `restart.include.`or`restart.exclude.`it is considered. |

| All `META-INF/spring-devtools.properties`from the classpath are loaded.
You can package files inside your project, or in the libraries that the project consumes.
System properties can not be used, only the properties file. |

### Known Limitations

Restart functionality does not work well with objects that are deserialized by using a standard `ObjectInputStream`.
If you need to deserialize data, you may need to use Spring’s `ConfigurableObjectInputStream` in combination with `Thread.currentThread().getContextClassLoader()`.

Unfortunately, several third-party libraries deserialize without considering the context classloader. If you find such a problem, you need to request a fix with the original authors.

## LiveReload

| Given its decrease in popularity and support, the LiveReload feature is deprecated as of Spring Boot 4.1.0 with no replacement. |

The `spring-boot-devtools` module includes an embedded LiveReload server that can be used to trigger a browser refresh when a resource is changed.
LiveReload browser extensions are freely available for Chrome, Firefox and Safari.
You can find these extensions by searching 'LiveReload' in the marketplace or store of your chosen browser.

If you want to start the LiveReload server when your application runs, you can set the `spring.devtools.livereload.enabled` property to `true`.

| You can only run one LiveReload server at a time. Before starting your application, ensure that no other LiveReload servers are running. If you start multiple applications from your IDE, only the first has LiveReload support. |

| To trigger LiveReload when a file changes, Automatic Restart must be enabled. |

## Global Settings

You can configure global devtools settings by adding any of the following files to the `$HOME/.config/spring-boot` directory:

-
`spring-boot-devtools.properties`
-
`spring-boot-devtools.yaml`
-
`spring-boot-devtools.yml`

Any properties added to these files apply to *all* Spring Boot applications on your machine that use devtools.
For example, to configure restart to always use a trigger file, you would add the following property to your `spring-boot-devtools` file:

-
Properties
-
YAML

`spring.devtools.restart.trigger-file=.reloadtrigger````
spring:
 devtools:
 restart:
 trigger-file: ".reloadtrigger"
```
By default, `$HOME` is the user’s home directory.
To customize this location, set the `SPRING_DEVTOOLS_HOME` environment variable or the `spring.devtools.home` system property.

| If devtools configuration files are not found in `$HOME/.config/spring-boot`, the root of the`$HOME`directory is searched for the presence of a`.spring-boot-devtools.properties`file.
This allows you to share the devtools global configuration with applications that are on an older version of Spring Boot that does not support the`$HOME/.config/spring-boot`location. |

| Profiles are not supported in devtools properties/yaml files. Any profiles activated in |

### Configuring File System Watcher

`FileSystemWatcher` works by polling the class changes with a certain time interval, and then waiting for a predefined quiet period to make sure there are no more changes.
Since Spring Boot relies entirely on the IDE to compile and copy files into the location from where Spring Boot can read them, you might find that there are times when certain changes are not reflected when devtools restarts the application.
If you observe such problems constantly, try increasing the `spring.devtools.restart.poll-interval` and `spring.devtools.restart.quiet-period` parameters to the values that fit your development environment:

-
Properties
-
YAML

```
spring.devtools.restart.poll-interval=2s
spring.devtools.restart.quiet-period=1s
```
```
spring:
 devtools:
 restart:
 poll-interval: "2s"
 quiet-period: "1s"
```
The monitored classpath directories are now polled every 2 seconds for changes, and a 1 second quiet period is maintained to make sure there are no additional class changes.

## Remote Applications

The Spring Boot developer tools are not limited to local development. You can also use several features when running applications remotely. Remote support is opt-in as enabling it can be a security risk. It should only be enabled when running on a trusted network or when secured with SSL. If neither of these options is available to you, you should not use DevTools' remote support. You should never enable support on a production deployment.

To enable it, you need to make sure that `devtools` is included in the repackaged archive, as shown in the following listing:

```
<build>
	<plugins>
 <plugin>
 <groupId>org.springframework.boot</groupId>
 <artifactId>spring-boot-maven-plugin</artifactId>
 <configuration>
 <includeOptional>true</includeOptional>
 <excludeDevtools>false</excludeDevtools>
 </configuration>
 </plugin>
	</plugins>
</build>
```
| Optional dependencies are not included by default, which explains why `includeOptional`is also present. |

Then you need to set the `spring.devtools.remote.secret` property.
Like any important password or secret, the value should be unique and strong such that it cannot be guessed or brute-forced.

Remote devtools support is provided in two parts: a server-side endpoint that accepts connections and a client application that you run in your IDE.
The server component is automatically enabled when the `spring.devtools.remote.secret` property is set.
The client component must be launched manually.

| Remote devtools is not supported for Spring WebFlux applications. |

### Running the Remote Client Application

The remote client application is designed to be run from within your IDE.
You need to run `RemoteSpringApplication` with the same classpath as the remote project that you connect to.
The application’s single required argument is the remote URL to which it connects.

For example, if you are using Eclipse or Spring Tools and you have a project named `my-app` that you have deployed to Cloud Foundry, you would do the following:

-
Select `Run Configurations…`from the`Run`menu.
-
Create a new `Java Application`“launch configuration”.
-
Browse for the `my-app`project.
-
Use `RemoteSpringApplication`as the main class.
-
Add `https://myapp.cfapps.io`to the`Program arguments`(or whatever your remote URL is).

A running remote client might resemble the following listing:

```
 . ____ _ __ _ _
 /\\ / ___'_ __ _ _(_)_ __ __ _ ___ _ \ \ \ \
( ( )\___ | '_ | '_| | '_ \/ _` | | _ \___ _ __ ___| |_ ___ \ \ \ \
 \\/ ___)| |_)| | | | | || (_| []::::::[] / -_) ' \/ _ \ _/ -_) ) ) ) )
 ' |____| .__|_| |_|_| |_\__, | |_|_\___|_|_|_\___/\__\___|/ / / /
 =========|_|==============|___/===================================/_/_/_/
 :: Spring Boot Remote :: (v4.1.0)
2026-06-10T16:33:47.783Z INFO 54045 --- [ main] o.s.b.devtools.RemoteSpringApplication : Starting RemoteSpringApplication v4.1.0 using Java 25.0.3 with PID 54045 (/Users/myuser/.m2/repository/org/springframework/boot/spring-boot-devtools/4.1.0/spring-boot-devtools-4.1.0.jar started by myuser in /opt/apps/)
2026-06-10T16:33:47.793Z INFO 54045 --- [ main] o.s.b.devtools.RemoteSpringApplication : No active profile set, falling back to 1 default profile: "default"
2026-06-10T16:33:48.836Z INFO 54045 --- [ main] o.s.b.d.a.OptionalLiveReloadServer : LiveReload server is running on port 35729
2026-06-10T16:33:48.936Z INFO 54045 --- [ main] o.s.b.devtools.RemoteSpringApplication : Started RemoteSpringApplication in 2.321 seconds (process running for 3.912)
```
| Because the remote client is using the same classpath as the real application it can directly read application properties.
This is how the `spring.devtools.remote.secret`property is read and passed to the server for authentication. |

| It is always advisable to use `https://`as the connection protocol, so that traffic is encrypted and passwords cannot be intercepted. |

| If you need to use a proxy to access the remote application, configure the `spring.devtools.remote.proxy.host`and`spring.devtools.remote.proxy.port`properties. |

### Remote Update

The remote client monitors your application classpath for changes in the same way as the local restart.
Any updated resource is pushed to the remote application and (*if required*) triggers a restart.
This can be helpful if you iterate on a feature that uses a cloud service that you do not have locally.
Generally, remote updates and restarts are much quicker than a full rebuild and deploy cycle.

On a slower development environment, it may happen that the quiet period is not enough, and the changes in the classes may be split into batches. The server is restarted after the first batch of class changes is uploaded. The next batch can’t be sent to the application, since the server is restarting.

This is typically manifested by a warning in the `RemoteSpringApplication` logs about failing to upload some of the classes, and a consequent retry.
But it may also lead to application code inconsistency and failure to restart after the first batch of changes is uploaded.
If you observe such problems constantly, try increasing the `spring.devtools.restart.poll-interval` and `spring.devtools.restart.quiet-period` parameters to the values that fit your development environment.
See the Configuring File System Watcher section for configuring these properties.

| Files are only monitored when the remote client is running. If you change a file before starting the remote client, it is not pushed to the remote server. |

# Packaging Your Application for Production

Once your Spring Boot application is ready for production deployment, there are many options for packaging and optimizing the application. See the Packaging Spring Boot Applications section of the documentation to read about these features.

For additional "production ready" features, such as health, auditing, and metric REST or JMX end-points, consider adding `spring-boot-actuator`.
See Actuator for details.

# Core Features

This section dives into the details of Spring Boot. Here you can learn about the key features that you may want to use and customize. If you have not already done so, you might want to read the Tutorials and Developing with Spring Boot sections, so that you have a good grounding of the basics.

# SpringApplication

The `SpringApplication` class provides a convenient way to bootstrap a Spring application that is started from a `main()` method.
In many situations, you can delegate to the static `SpringApplication.run(Class, String...)` method, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
@SpringBootApplication
public class MyApplication {
	public static void main(String[] args) {
 SpringApplication.run(MyApplication.class, args);
	}
}
```
```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
@SpringBootApplication
class MyApplication
fun main(args: Array<String>) {
	runApplication<MyApplication>(*args)
}
```
When your application starts, you should see something similar to the following output:

```
 . ____ _ __ _ _
 /\\ / ___'_ __ _ _(_)_ __ __ _ \ \ \ \
( ( )\___ | '_ | '_| | '_ \/ _` | \ \ \ \
 \\/ ___)| |_)| | | | | || (_| | ) ) ) )
 ' |____| .__|_| |_|_| |_\__, | / / / /
 =========|_|==============|___/=/_/_/_/
 :: Spring Boot :: (v4.1.0)
2026-06-10T16:33:51.740Z INFO 54551 --- [ main] o.s.b.d.f.logexample.MyApplication : Starting MyApplication using Java 25.0.3 with PID 54551 (/opt/apps/myapp.jar started by myuser in /opt/apps/)
2026-06-10T16:33:51.747Z INFO 54551 --- [ main] o.s.b.d.f.logexample.MyApplication : No active profile set, falling back to 1 default profile: "default"
2026-06-10T16:33:55.110Z INFO 54551 --- [ main] o.s.boot.tomcat.TomcatWebServer : Tomcat initialized with port 8080 (http)
2026-06-10T16:33:55.196Z INFO 54551 --- [ main] o.apache.catalina.core.StandardService : Starting service [Tomcat]
2026-06-10T16:33:55.205Z INFO 54551 --- [ main] o.apache.catalina.core.StandardEngine : Starting Servlet engine: [Apache Tomcat/11.0.22]
2026-06-10T16:33:55.376Z INFO 54551 --- [ main] b.w.c.s.WebApplicationContextInitializer : Root WebApplicationContext: initialization completed in 3421 ms
2026-06-10T16:33:56.597Z INFO 54551 --- [ main] o.s.boot.tomcat.TomcatWebServer : Tomcat started on port 8080 (http) with context path '/'
2026-06-10T16:33:56.628Z INFO 54551 --- [ main] o.s.b.d.f.logexample.MyApplication : Started MyApplication in 6.318 seconds (process running for 7.637)
2026-06-10T16:33:56.668Z INFO 54551 --- [ionShutdownHook] o.s.boot.tomcat.GracefulShutdown : Commencing graceful shutdown. Waiting for active requests to complete
2026-06-10T16:33:56.709Z INFO 54551 --- [tomcat-shutdown] o.s.boot.tomcat.GracefulShutdown : Graceful shutdown complete
```
By default, `INFO` logging messages are shown, including some relevant startup details, such as the user that launched the application.
If you need a log level other than `INFO`, you can set it, as described in Log Levels.
The application version is determined using the implementation version from the main application class’s package.
Startup information logging can be turned off by setting `spring.main.log-startup-info` to `false`.
This will also turn off logging of the application’s active profiles.

| To add additional logging during startup, you can override `logStartupInfo(boolean)`in a subclass of`SpringApplication`. |

## Startup Failure

If your application fails to start, registered `FailureAnalyzer` beans get a chance to provide a dedicated error message and a concrete action to fix the problem.
For instance, if you start a web application on port `8080` and that port is already in use, you should see something similar to the following message:

```
***************************
APPLICATION FAILED TO START
***************************
Description:
Embedded servlet container failed to start. Port 8080 was already in use.
Action:
Identify and stop the process that is listening on port 8080 or configure this application to listen on another port.
```
| Spring Boot provides numerous `FailureAnalyzer`implementations, and you can add your own. |

If no failure analyzers are able to handle the exception, you can still display the full conditions report to better understand what went wrong.
To do so, you need to enable the `debug` property or enable `DEBUG` logging for `ConditionEvaluationReportLoggingListener`.

For instance, if you are running your application by using `java -jar`, you can enable the `debug` property as follows:

`$ java -jar myproject-0.0.1-SNAPSHOT.jar --debug`## Lazy Initialization

`SpringApplication` allows an application to be initialized lazily.
When lazy initialization is enabled, beans are created as they are needed rather than during application startup.
As a result, enabling lazy initialization can reduce the time that it takes your application to start.
In a web application, enabling lazy initialization will result in many web-related beans not being initialized until an HTTP request is received.

A downside of lazy initialization is that it can delay the discovery of a problem with the application. If a misconfigured bean is initialized lazily, a failure will no longer occur during startup and the problem will only become apparent when the bean is initialized. Care must also be taken to ensure that the JVM has sufficient memory to accommodate all of the application’s beans and not just those that are initialized during startup. For these reasons, lazy initialization is not enabled by default and it is recommended that fine-tuning of the JVM’s heap size is done before enabling lazy initialization.

Lazy initialization can be enabled programmatically using the `lazyInitialization` method on `SpringApplicationBuilder` or the `setLazyInitialization` method on `SpringApplication`.
Alternatively, it can be enabled using the `spring.main.lazy-initialization` property as shown in the following example:

-
Properties
-
YAML

`spring.main.lazy-initialization=true````
spring:
 main:
 lazy-initialization: true
```
| If you want to disable lazy initialization for certain beans while using lazy initialization for the rest of the application, you can explicitly set their lazy attribute to false using the `@Lazy(false)`annotation. |

## Customizing the Banner

The banner that is printed on start up can be changed by adding a `banner.txt` file to your classpath or by setting the `spring.banner.location` property to the location of such a file.
If the file has an encoding other than UTF-8, you can set `spring.banner.charset`.

Inside your `banner.txt` file, you can use any key available in the `Environment` as well as any of the following placeholders:

| Variable | Description |
|---|---|
|
 | The version number of your application, as declared in |
|
 | The version number of your application, as declared in |
|
 | The Spring Boot version that you are using.
 For example |
|
 | The Spring Boot version that you are using, formatted for display (surrounded with brackets and prefixed with |
|
 | Where |
|
 | The title of your application, as declared in |

| The `SpringApplication.setBanner(…)`method can be used if you want to generate a banner programmatically.
Use the`Banner`interface and implement your own`printBanner()`method. |

You can also use the `spring.main.banner-mode` property to determine if the banner has to be printed on `System.out` (`console`), sent to the configured logger (`log`), or not produced at all (`off`).

The printed banner is registered as a singleton bean under the following name: `springBootBanner`.

| The To use the |

## Customizing SpringApplication

If the `SpringApplication` defaults are not to your taste, you can instead create a local instance and customize it.
For example, to turn off the banner, you could write:

-
Java
-
Kotlin

```
import org.springframework.boot.Banner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
@SpringBootApplication
public class MyApplication {
	public static void main(String[] args) {
 SpringApplication application = new SpringApplication(MyApplication.class);
 application.setBannerMode(Banner.Mode.OFF);
 application.run(args);
	}
}
```
```
import org.springframework.boot.Banner
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
@SpringBootApplication
class MyApplication
fun main(args: Array<String>) {
	runApplication<MyApplication>(*args) {
 setBannerMode(Banner.Mode.OFF)
	}
}
```
| The constructor arguments passed to `SpringApplication`are configuration sources for Spring beans.
In most cases, these are references to`@Configuration`classes, but they could also be direct references`@Component`classes. |

It is also possible to configure the `SpringApplication` by using an `application.properties` file.
See Externalized Configuration for details.

For a complete list of the configuration options, see the `SpringApplication` API documentation.

## Fluent Builder API

If you need to build an `ApplicationContext` hierarchy (multiple contexts with a parent/child relationship) or if you prefer using a fluent builder API, you can use the `SpringApplicationBuilder`.

The `SpringApplicationBuilder` lets you chain together multiple method calls and includes `parent` and `child` methods that let you create a hierarchy, as shown in the following example:

-
Java
-
Kotlin

```
 new SpringApplicationBuilder().sources(Parent.class)
 .child(Application.class)
 .bannerMode(Banner.Mode.OFF)
 .run(args);
```
```
 SpringApplicationBuilder()
 .sources(Parent::class.java)
 .child(Application::class.java)
 .bannerMode(Banner.Mode.OFF)
 .run(*args)
```
| There are some restrictions when creating an `ApplicationContext`hierarchy.
For example, Web componentsmustbe contained within the child context, and the same`Environment`is used for both parent and child contexts.
See the`SpringApplicationBuilder`API documentation for full details. |

## Application Availability

When deployed on platforms, applications can provide information about their availability to the platform using infrastructure such as Kubernetes Probes. Spring Boot includes out-of-the box support for the commonly used “liveness” and “readiness” availability states. If you are using Spring Boot’s “actuator” support then these states are exposed as health endpoint groups.

In addition, you can also obtain availability states by injecting the `ApplicationAvailability` interface into your own beans.

### Liveness State

The “Liveness” state of an application tells whether its internal state allows it to work correctly, or recover by itself if it is currently failing. A broken “Liveness” state means that the application is in a state that it cannot recover from, and the infrastructure should restart the application.

| In general, the "Liveness" state should not be based on external checks, such as health checks. If it did, a failing external system (a database, a Web API, an external cache) would trigger massive restarts and cascading failures across the platform. |

The internal state of Spring Boot applications is mostly represented by the Spring `ApplicationContext`.
If the application context has started successfully, Spring Boot assumes that the application is in a valid state.
An application is considered live as soon as the context has been refreshed, see Spring Boot application lifecycle and related Application Events.

### Readiness State

The “Readiness” state of an application tells whether the application is ready to handle traffic.
A failing “Readiness” state tells the platform that it should not route traffic to the application for now.
This typically happens during startup, while `CommandLineRunner` and `ApplicationRunner` components are being processed, or at any time if the application decides that it is too busy for additional traffic.

An application is considered ready as soon as application and command-line runners have been called, see Spring Boot application lifecycle and related Application Events.

| Tasks expected to run during startup should be executed by `CommandLineRunner`and`ApplicationRunner`components instead of using Spring component lifecycle callbacks such as`@PostConstruct`. |

### Managing the Application Availability State

Application components can retrieve the current availability state at any time, by injecting the `ApplicationAvailability` interface and calling methods on it.
More often, applications will want to listen to state updates or update the state of the application.

For example, we can export the "Readiness" state of the application to a file so that a Kubernetes "exec Probe" can look at this file:

-
Java
-
Kotlin

```
import org.springframework.boot.availability.AvailabilityChangeEvent;
import org.springframework.boot.availability.ReadinessState;
import org.springframework.context.event.EventListener;
import org.springframework.stereotype.Component;
@Component
public class MyReadinessStateExporter {
	@EventListener
	public void onStateChange(AvailabilityChangeEvent<ReadinessState> event) {
 switch (event.getState()) {
 case ACCEPTING_TRAFFIC -> {
 // create file /tmp/healthy
 }
 case REFUSING_TRAFFIC -> {
 // remove file /tmp/healthy
 }
 }
	}
}
```
```
import org.springframework.boot.availability.AvailabilityChangeEvent
import org.springframework.boot.availability.ReadinessState
import org.springframework.context.event.EventListener
import org.springframework.stereotype.Component
@Component
class MyReadinessStateExporter {
	@EventListener
	fun onStateChange(event: AvailabilityChangeEvent<ReadinessState>) {
 when (event.state) {
 ReadinessState.ACCEPTING_TRAFFIC -> {
 // create file /tmp/healthy
 }
 ReadinessState.REFUSING_TRAFFIC -> {
 // remove file /tmp/healthy
 }
 }
	}
}
```
We can also update the state of the application, when the application breaks and cannot recover:

-
Java
-
Kotlin

```
import org.springframework.boot.availability.AvailabilityChangeEvent;
import org.springframework.boot.availability.LivenessState;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.stereotype.Component;
@Component
public class MyLocalCacheVerifier {
	private final ApplicationEventPublisher eventPublisher;
	public MyLocalCacheVerifier(ApplicationEventPublisher eventPublisher) {
 this.eventPublisher = eventPublisher;
	}
	public void checkLocalCache() {
 try {
 // ...
 }
 catch (CacheCompletelyBrokenException ex) {
 AvailabilityChangeEvent.publish(this.eventPublisher, ex, LivenessState.BROKEN);
 }
	}
}
```
```
import org.springframework.boot.availability.AvailabilityChangeEvent
import org.springframework.boot.availability.LivenessState
import org.springframework.context.ApplicationEventPublisher
import org.springframework.stereotype.Component
@Component
class MyLocalCacheVerifier(private val eventPublisher: ApplicationEventPublisher) {
	fun checkLocalCache() {
 try {
 // ...
 } catch (ex: CacheCompletelyBrokenException) {
 AvailabilityChangeEvent.publish(eventPublisher, ex, LivenessState.BROKEN)
 }
	}
}
```
Spring Boot provides Kubernetes HTTP probes for "Liveness" and "Readiness" with Actuator Health Endpoints. You can get more guidance about deploying Spring Boot applications on Kubernetes in the dedicated section.

## Application Events and Listeners

In addition to the usual Spring Framework events, such as `ContextRefreshedEvent`, a `SpringApplication` sends some additional application events.

| Some events are actually triggered before the If you want those listeners to be registered automatically, regardless of the way the application is created, you can add a |

Application events are sent in the following order, as your application runs:

-
An `ApplicationStartingEvent`is sent at the start of a run but before any processing, except for the registration of listeners and initializers.
-
An `ApplicationEnvironmentPreparedEvent`is sent when the`Environment`to be used in the context is known but before the context is created.
-
An `ApplicationContextInitializedEvent`is sent when the`ApplicationContext`is prepared and ApplicationContextInitializers have been called but before any bean definitions are loaded.
-
An `ApplicationPreparedEvent`is sent just before the refresh is started but after bean definitions have been loaded.
-
An `ApplicationStartedEvent`is sent after the context has been refreshed but before any application and command-line runners have been called.
-
An `AvailabilityChangeEvent`is sent right after with`LivenessState.CORRECT`to indicate that the application is considered as live.
-
An `ApplicationReadyEvent`is sent after any application and command-line runners have been called.
-
An `AvailabilityChangeEvent`is sent right after with`ReadinessState.ACCEPTING_TRAFFIC`to indicate that the application is ready to service requests.
-
An `ApplicationFailedEvent`is sent if there is an exception on startup.

The above list only includes `SpringApplicationEvent`s that are tied to a `SpringApplication`.
In addition to these, the following events are also published after `ApplicationPreparedEvent` and before `ApplicationStartedEvent`:

-
A `WebServerInitializedEvent`is sent after the`WebServer`is ready.`ServletWebServerInitializedEvent`and`ReactiveWebServerInitializedEvent`are the servlet and reactive variants respectively.
-
A `ContextRefreshedEvent`is sent when an`ApplicationContext`is refreshed.

| You often need not use application events, but it can be handy to know that they exist. Internally, Spring Boot uses events to handle a variety of tasks. |

| Event listeners should not run potentially lengthy tasks as they execute in the same thread by default. Consider using application and command-line runners instead. |

Application events are sent by using Spring Framework’s event publishing mechanism.
Part of this mechanism ensures that an event published to the listeners in a child context is also published to the listeners in any ancestor contexts.
As a result of this, if your application uses a hierarchy of `SpringApplication` instances, a listener may receive multiple instances of the same type of application event.

To allow your listener to distinguish between an event for its context and an event for a descendant context, it should request that its application context is injected and then compare the injected context with the context of the event.
The context can be injected by implementing `ApplicationContextAware` or, if the listener is a bean, by using `@Autowired`.

## Web Environment

A `SpringApplication` attempts to create the right type of `ApplicationContext` on your behalf.
The algorithm used to determine a `WebApplicationType` is the following:

-
If Spring MVC is present, an `AnnotationConfigServletWebServerApplicationContext`is used
-
If Spring MVC is not present and Spring WebFlux is present, an `AnnotationConfigReactiveWebServerApplicationContext`is used
-
Otherwise, `AnnotationConfigApplicationContext`is used

This means that if you are using Spring MVC and the new `WebClient` from Spring WebFlux in the same application, Spring MVC will be used by default.
You can override that easily by calling `setWebApplicationType(WebApplicationType)`.

It is also possible to take complete control of the `ApplicationContext` type that is used by calling `setApplicationContextFactory(…)`.

| It is often desirable to call `setWebApplicationType(WebApplicationType.NONE)`when using`SpringApplication`within a JUnit test. |

## Accessing Application Arguments

If you need to access the application arguments that were passed to `SpringApplication.run(…)`, you can inject a `ApplicationArguments` bean.
The `ApplicationArguments` interface provides access to both the raw `String[]` arguments as well as parsed `option` and `non-option` arguments, as shown in the following example:

-
Java
-
Kotlin

```
import java.util.List;
import org.springframework.boot.ApplicationArguments;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	public MyBean(ApplicationArguments args) {
 boolean debug = args.containsOption("debug");
 List<String> files = args.getNonOptionArgs();
 if (debug) {
 System.out.println(files);
 }
 // if run with "--debug logfile.txt" prints ["logfile.txt"]
	}
}
```
```
import org.springframework.boot.ApplicationArguments
import org.springframework.stereotype.Component
@Component
class MyBean(args: ApplicationArguments) {
	init {
 val debug = args.containsOption("debug")
 val files = args.nonOptionArgs
 if (debug) {
 println(files)
 }
 // if run with "--debug logfile.txt" prints ["logfile.txt"]
	}
}
```
| Spring Boot also registers a `CommandLinePropertySource`with the Spring`Environment`.
This lets you also inject single application arguments by using the`@Value`annotation. |

## Using the ApplicationRunner or CommandLineRunner

If you need to run some specific code once the `SpringApplication` has started, you can implement the `ApplicationRunner` or `CommandLineRunner` interfaces.
Both interfaces work in the same way and offer a single `run` method, which is called just before `SpringApplication.run(…)` completes.

| This contract is well suited for tasks that should run after application startup but before it starts accepting traffic. |

The `CommandLineRunner` interfaces provides access to application arguments as a string array, whereas the `ApplicationRunner` uses the `ApplicationArguments` interface discussed earlier.
The following example shows a `CommandLineRunner` with a `run` method:

-
Java
-
Kotlin

```
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Component;
@Component
public class MyCommandLineRunner implements CommandLineRunner {
	@Override
	public void run(String... args) {
 // Do something...
	}
}
```
```
import org.springframework.boot.CommandLineRunner
import org.springframework.stereotype.Component
@Component
class MyCommandLineRunner : CommandLineRunner {
	override fun run(vararg args: String) {
 // Do something...
	}
}
```
If several `CommandLineRunner` or `ApplicationRunner` beans are defined that must be called in a specific order, you can additionally implement the `Ordered` interface or use the `Order` annotation.

## Application Exit

Each `SpringApplication` registers a shutdown hook with the JVM to ensure that the `ApplicationContext` closes gracefully on exit.
All the standard Spring lifecycle callbacks (such as the `DisposableBean` interface or the `@PreDestroy` annotation) can be used.

In addition, beans may implement the `ExitCodeGenerator` interface if they wish to return a specific exit code when `SpringApplication.exit()` is called.
This exit code can then be passed to `System.exit()` to return it as a status code, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.ExitCodeGenerator;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
@SpringBootApplication
public class MyApplication {
	@Bean
	public ExitCodeGenerator exitCodeGenerator() {
 return () -> 42;
	}
	public static void main(String[] args) {
 System.exit(SpringApplication.exit(SpringApplication.run(MyApplication.class, args)));
	}
}
```
```
import org.springframework.boot.ExitCodeGenerator
import org.springframework.boot.SpringApplication
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.context.annotation.Bean
import kotlin.system.exitProcess
@SpringBootApplication
class MyApplication {
	@Bean
	fun exitCodeGenerator() = ExitCodeGenerator { 42 }
}
fun main(args: Array<String>) {
	exitProcess(SpringApplication.exit(
 runApplication<MyApplication>(*args)))
}
```
Also, the `ExitCodeGenerator` interface may be implemented by exceptions.
When such an exception is encountered, Spring Boot returns the exit code provided by the implemented `getExitCode()` method.

If there is more than one `ExitCodeGenerator`, the first non-zero exit code that is generated is used.
To control the order in which the generators are called, additionally implement the `Ordered` interface or use the `Order` annotation.

## Admin Features

It is possible to enable admin-related features for the application by specifying the `spring.application.admin.enabled` property.
This exposes the `SpringApplicationAdminMXBean` on the platform `MBeanServer`.
You could use this feature to administer your Spring Boot application remotely.
This feature could also be useful for any service wrapper implementation.

| If you want to know on which HTTP port the application is running, get the property with a key of `local.server.port`. |

## Application Startup tracking

During the application startup, the `SpringApplication` and the `ApplicationContext` perform many tasks related to the application lifecycle,
the beans lifecycle or even processing application events.
With `ApplicationStartup`, Spring Framework allows you to track the application startup sequence with `StartupStep` objects.
This data can be collected for profiling purposes, or just to have a better understanding of an application startup process.

You can choose an `ApplicationStartup` implementation when setting up the `SpringApplication` instance.
For example, to use the `BufferingApplicationStartup`, you could write:

-
Java
-
Kotlin

```
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.metrics.buffering.BufferingApplicationStartup;
@SpringBootApplication
public class MyApplication {
	public static void main(String[] args) {
 SpringApplication application = new SpringApplication(MyApplication.class);
 application.setApplicationStartup(new BufferingApplicationStartup(2048));
 application.run(args);
	}
}
```
```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.context.metrics.buffering.BufferingApplicationStartup
import org.springframework.boot.runApplication
@SpringBootApplication
class MyApplication
fun main(args: Array<String>) {
	runApplication<MyApplication>(*args) {
 applicationStartup = BufferingApplicationStartup(2048)
	}
}
```
The first available implementation, `FlightRecorderApplicationStartup` is provided by Spring Framework.
It adds Spring-specific startup events to a Java Flight Recorder session and is meant for profiling applications and correlating their Spring context lifecycle with JVM events (such as allocations, GCs, class loading…).
Once configured, you can record data by running the application with the Flight Recorder enabled:

`$ java -XX:StartFlightRecording:filename=recording.jfr,duration=10s -jar demo.jar`Spring Boot ships with the `BufferingApplicationStartup` variant; this implementation is meant for buffering the startup steps and draining them into an external metrics system.
Applications can ask for the bean of type `BufferingApplicationStartup` in any component.

Spring Boot can also be configured to expose a `startup` endpoint that provides this information as a JSON document.

## Virtual threads

Virtual threads require Java 21 or later.
For the best experience, Java 24 or later is strongly recommended.
To enable virtual threads, set the `spring.threads.virtual.enabled` property to `true`.

Before turning on this option for your application, you should consider reading the official Java virtual threads documentation.
In some cases, applications can experience lower throughput because of "Pinned Virtual Threads"; this page also explains how to detect such cases with JDK Flight Recorder or the `jcmd` CLI.

| If virtual threads are enabled, properties which configure thread pools don’t have an effect anymore. That’s because virtual threads are scheduled on a JVM wide platform thread pool and not on dedicated thread pools. |

| One side effect of virtual threads is that they are daemon threads.
A JVM will exit if all of its threads are daemon threads.
This behavior can be a problem when you rely on `@Scheduled`beans, for example, to keep your application alive.
If you use virtual threads, the scheduler thread is a virtual thread and therefore a daemon thread and won’t keep the JVM alive.
This not only affects scheduling and can be the case with other technologies too.
To keep the JVM running in all cases, it is recommended to set the property`spring.main.keep-alive`to`true`.
This ensures that the JVM is kept alive, even if all threads are virtual threads. |

# Externalized Configuration

Spring Boot lets you externalize your configuration so that you can work with the same application code in different environments. You can use a variety of external configuration sources including Java properties files, YAML files, environment variables, and command-line arguments.

Property values can be injected directly into your beans by using the `@Value` annotation, accessed through Spring’s `Environment` abstraction, or be bound to structured objects through `@ConfigurationProperties`.

Spring Boot uses a very particular `PropertySource` order that is designed to allow sensible overriding of values.
Later property sources can override the values defined in earlier ones.

Sources are considered in the following order:

-
Default properties (specified by setting `SpringApplication.setDefaultProperties(Map)`).
-
`@PropertySource`annotations on your`@Configuration`classes. Please note that such property sources are not added to the`Environment`until the application context is being refreshed. This is too late to configure certain properties such as`logging.*`and`spring.main.*`which are read before refresh begins.
-
Config data (such as `application.properties`files).
-
A `RandomValuePropertySource`that has properties only in`random.*`.
-
OS environment variables.
-
Java System properties ( `System.getProperties()`).
-
JNDI attributes from `java:comp/env`.
-
`ServletContext`init parameters.
-
`ServletConfig`init parameters.
-
Properties from `SPRING_APPLICATION_JSON`(inline JSON embedded in an environment variable or system property).
-
Command line arguments.
-
`properties`attribute on your tests. Available on`@SpringBootTest`and the test annotations for testing a particular slice of your application.
-
`@DynamicPropertySource`annotations in your tests.
-
`@TestPropertySource`annotations on your tests.
-
Devtools global settings properties in the `$HOME/.config/spring-boot`directory when devtools is active.

Config data files are considered in the following order:

-
Application properties packaged inside your jar ( `application.properties`and YAML variants).
-
Profile-specific application properties packaged inside your jar ( `application-{profile}.properties`and YAML variants).
-
Application properties outside of your packaged jar ( `application.properties`and YAML variants).
-
Profile-specific application properties outside of your packaged jar ( `application-{profile}.properties`and YAML variants).

| It is recommended to stick with one format for your entire application.
If you have configuration files with both `.properties`and YAML format in the same location,`.properties`takes precedence. |

| If you use environment variables rather than system properties, most operating systems disallow period-separated key names, but you can use underscores instead (for example, `SPRING_CONFIG_NAME`instead of`spring.config.name`).
See Binding From Environment Variables for details. |

| If your application runs in a servlet container or application server, then JNDI properties (in `java:comp/env`) or servlet context initialization parameters can be used instead of, or as well as, environment variables or system properties. |

To provide a concrete example, suppose you develop a `@Component` that uses a `name` property, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	@Value("${name}")
	private String name;
	// ...
}
```
```
import org.springframework.beans.factory.annotation.Value
import org.springframework.stereotype.Component
@Component
class MyBean {
	@Value("\${name}")
	private val name: String? = null
	// ...
}
```
On your application classpath (for example, inside your jar) you can have an `application.properties` file that provides a sensible default property value for `name`.
When running in a new environment, an `application.properties` file can be provided outside of your jar that overrides the `name`.
For one-off testing, you can launch with a specific command line switch (for example, `java -jar app.jar --name="Spring"`).

| The `env`and`configprops`endpoints can be useful in determining why a property has a particular value.
You can use these two endpoints to diagnose unexpected property values.
See the Production ready features section for details. |

## Accessing Command Line Properties

By default, `SpringApplication` converts any command line option arguments (that is, arguments starting with `--`, such as `--server.port=9000`) to a `property` and adds them to the Spring `Environment`.
As mentioned previously, command line properties always take precedence over file-based property sources.

If you do not want command line properties to be added to the `Environment`, you can disable them by using `SpringApplication.setAddCommandLineProperties(false)`.

## JSON Application Properties

Environment variables and system properties often have restrictions that mean some property names cannot be used. To help with this, Spring Boot allows you to encode a block of properties into a single JSON structure.

When your application starts, any `spring.application.json` or `SPRING_APPLICATION_JSON` properties will be parsed and added to the `Environment`.

For example, the `SPRING_APPLICATION_JSON` property can be supplied on the command line in a UN*X shell as an environment variable:

`$ SPRING_APPLICATION_JSON='{"my":{"name":"test"}}' java -jar myapp.jar`In the preceding example, you end up with `my.name=test` in the Spring `Environment`.

The same JSON can also be provided as a system property:

`$ java -Dspring.application.json='{"my":{"name":"test"}}' -jar myapp.jar`Or you could supply the JSON by using a command line argument:

`$ java -jar myapp.jar --spring.application.json='{"my":{"name":"test"}}'`If you are deploying to a classic Application Server, you could also use a JNDI variable named `java:comp/env/spring.application.json`.

| Although `null`values from the JSON will be added to the resulting property source, the`PropertySourcesPropertyResolver`treats`null`properties as missing values.
This means that the JSON cannot override properties from lower order property sources with a`null`value. |

## External Application Properties

Spring Boot will automatically find and load `application.properties` and `application.yaml` files from the following locations when your application starts:

-
From the classpath -
The classpath root
-
The classpath `/config`package

-
-
From the current directory -
The current directory
-
The `config/`subdirectory in the current directory
-
Immediate child directories of the `config/`subdirectory

-

The list is ordered by precedence (with values from lower items overriding earlier ones).
Documents from the loaded files are added as `PropertySource` instances to the Spring `Environment`.

If you do not like `application` as the configuration file name, you can switch to another file name by specifying a `spring.config.name` environment property.
For example, to look for `myproject.properties` and `myproject.yaml` files you can run your application as follows:

`$ java -jar myproject.jar --spring.config.name=myproject`You can also refer to an explicit location by using the `spring.config.location` environment property.
This property accepts a comma-separated list of one or more locations to check.

The following example shows how to specify two distinct files:

```
$ java -jar myproject.jar --spring.config.location=\
	optional:classpath:/default.properties,\
	optional:classpath:/override.properties
```
| Use the prefix `optional:`if the locations are optional and you do not mind if they do not exist. |

| `spring.config.name`,`spring.config.location`, and`spring.config.additional-location`are used very early to determine which files have to be loaded.
They must be defined as an environment property (typically an OS environment variable, a system property, or a command-line argument). |

If `spring.config.location` contains directories (as opposed to files), they should end in `/`.
At runtime they will be appended with the names generated from `spring.config.name` before being loaded.
Files specified in `spring.config.location` are imported directly.

| Both directory and file location values are also expanded to check for profile-specific files.
For example, if you have a `spring.config.location`of`classpath:myconfig.properties`, you will also find appropriate`classpath:myconfig-<profile>.properties`files are loaded. |

In most situations, each `spring.config.location` item you add will reference a single file or directory.
Locations are processed in the order that they are defined and later ones can override the values of earlier ones.

If you have a complex location setup, and you use profile-specific configuration files, you may need to provide further hints so that Spring Boot knows how they should be grouped.
A location group is a collection of locations that are all considered at the same level.
For example, you might want to group all classpath locations, then all external locations.
Items within a location group should be separated with `;`.
See the example in the Profile Specific Files section for more details.

Locations configured by using `spring.config.location` replace the default locations.
For example, if `spring.config.location` is configured with the value `optional:classpath:/custom-config/,optional:file:./custom-config/`, the complete set of locations considered is:

-
`optional:classpath:custom-config/`
-
`optional:file:./custom-config/`

If you prefer to add additional locations, rather than replacing them, you can use `spring.config.additional-location`.
Properties loaded from additional locations can override those in the default locations.
For example, if `spring.config.additional-location` is configured with the value `optional:classpath:/custom-config/,optional:file:./custom-config/`, the complete set of locations considered is:

-
`optional:classpath:/;optional:classpath:/config/`
-
`optional:file:./;optional:file:./config/;optional:file:./config/*/`
-
`optional:classpath:custom-config/`
-
`optional:file:./custom-config/`

This search ordering lets you specify default values in one configuration file and then selectively override those values in another.
You can provide default values for your application in `application.properties` (or whatever other basename you choose with `spring.config.name`) in one of the default locations.
These default values can then be overridden at runtime with a different file located in one of the custom locations.

### Optional Locations

By default, when a specified config data location does not exist, Spring Boot will throw a `ConfigDataLocationNotFoundException` and your application will not start.

If you want to specify a location, but you do not mind if it does not always exist, you can use the `optional:` prefix.
You can use this prefix with the `spring.config.location` and `spring.config.additional-location` properties, as well as with `spring.config.import` declarations.

For example, a `spring.config.import` value of `optional:file:./myconfig.properties` allows your application to start, even if the `myconfig.properties` file is missing.

If you want to ignore all `ConfigDataLocationNotFoundException` errors and always continue to start your application, you can use the `spring.config.on-not-found` property.
Set the value to `ignore` using `SpringApplication.setDefaultProperties(…)` or with a system/environment variable.

### Wildcard Locations

If a config file location includes the `*` character for the last path segment, it is considered a wildcard location.
Wildcards are expanded when the config is loaded so that immediate subdirectories are also checked.
Wildcard locations are particularly useful in an environment such as Kubernetes when there are multiple sources of config properties.

For example, if you have some Redis configuration and some MySQL configuration, you might want to keep those two pieces of configuration separate, while requiring that both those are present in an `application.properties` file.
This might result in two separate `application.properties` files mounted at different locations such as `/config/redis/application.properties` and `/config/mysql/application.properties`.
In such a case, having a wildcard location of `config/*/`, will result in both files being processed.

By default, Spring Boot includes `config/*/` in the default search locations.
It means that all subdirectories of the `/config` directory outside of your jar will be searched.

You can use wildcard locations yourself with the `spring.config.location` and `spring.config.additional-location` properties.

| A wildcard location must contain only one `*`and end with`*/`for search locations that are directories or`*/<filename>`for search locations that are files.
Locations with wildcards are sorted alphabetically based on the absolute path of the file names. |

| Wildcard locations only work with external directories.
You cannot use a wildcard in a `classpath:`location. |

### Profile Specific Files

As well as `application` property files, Spring Boot will also attempt to load profile-specific files using the naming convention `application-{profile}`.
For example, if your application activates a profile named `prod` and uses YAML files, then both `application.yaml` and `application-prod.yaml` will be considered.

Profile-specific properties are loaded from the same locations as standard `application.properties`, with profile-specific files always overriding the non-specific ones.
If several profiles are specified, a last-wins strategy applies.
For example, if profiles `prod,live` are specified by the `spring.profiles.active` property, values in `application-prod.properties` can be overridden by those in `application-live.properties`.

| The last-wins strategy applies at the location group level.
A For example, continuing our /cfg application-live.properties /ext application-live.properties application-prod.properties When we have a
 When we have
 |

The `Environment` has a set of default profiles (by default, `[default]`) that are used if no active profiles are set.
In other words, if no profiles are explicitly activated, then properties from `application-default` are considered.

| Properties files are only ever loaded once. If you have already directly imported a profile specific property files then it will not be imported a second time. |

### Importing Additional Data

Application properties may import further config data from other locations using the `spring.config.import` property.
Imports are processed as they are discovered, and are treated as additional documents inserted immediately below the one that declares the import.

For example, you might have the following in your classpath `application.properties` file:

-
Properties
-
YAML

```
spring.application.name=myapp
spring.config.import=optional:file:./dev.properties
```
```
spring:
 application:
 name: "myapp"
 config:
 import: "optional:file:./dev.properties"
```
This will trigger the import of a `dev.properties` file in current directory (if such a file exists).
Values from the imported `dev.properties` will take precedence over the file that triggered the import.
In the above example, the `dev.properties` could redefine `spring.application.name` to a different value.

| An import will only be imported once no matter how many times it is declared. |

By default, properties files are imported using the ISO-8859-1 charset. To change that, you can use the encoding attribute:

-
Properties
-
YAML

`spring.config.import=classpath:import.properties[encoding=utf-8]````
spring:
 config:
 import: "classpath:import.properties[encoding=utf-8]"
```
The `import.properties` file will now be read in UTF-8 encoding.

#### Using “Fixed” and “Import Relative” Locations

Imports may be specified as *fixed* or *import relative* locations.
A fixed location always resolves to the same underlying resource, regardless of where the `spring.config.import` property is declared.
An import relative location resolves relative to the file that declares the `spring.config.import` property.

A location starting with a forward slash (`/`) or a URL style prefix (`file:`, `classpath:`, etc.) is considered fixed.
All other locations are considered import relative.

| `optional:`prefixes are not considered when determining if a location is fixed or import relative. |

As an example, say we have a `/demo` directory containing our `application.jar` file.
We might add a `/demo/application.properties` file with the following content:

`spring.config.import=optional:core/core.properties`This is an import relative location and so will attempt to load the file `/demo/core/core.properties` if it exists.

If `/demo/core/core.properties` has the following content:

`spring.config.import=optional:extra/extra.properties`It will attempt to load `/demo/core/extra/extra.properties`.
The `optional:extra/extra.properties` is relative to `/demo/core/core.properties` so the full directory is `/demo/core/` + `extra/extra.properties`.

#### Property Ordering

The order an import is defined inside a single document within the properties/yaml file does not matter. For instance, the two examples below produce the same result:

-
Properties
-
YAML

```
spring.config.import=my.properties
my.property=value
```
```
spring:
 config:
 import: "my.properties"
my:
 property: "value"
```
-
Properties
-
YAML

```
my.property=value
spring.config.import=my.properties
```
```
my:
 property: "value"
spring:
 config:
 import: "my.properties"
```
In both of the above examples, the values from the `my.properties` file will take precedence over the file that triggered its import.

Several locations can be specified under a single `spring.config.import` key.
Locations will be processed in the order that they are defined, with later imports taking precedence.

| When appropriate, Profile-specific variants are also considered for import.
The example above would import both `my.properties`as well as any`my-<profile>.properties`variants. |

| Spring Boot includes pluggable API that allows various different location addresses to be supported. By default you can import Java Properties, YAML and configuration trees. Third-party jars can offer support for additional technologies (there is no requirement for files to be local). For example, you can imagine config data being from external stores such as Consul, Apache ZooKeeper or Netflix Archaius. If you want to support your own locations, see the |

### Importing Extensionless Files

Some cloud platforms cannot add a file extension to volume mounted files. To import these extensionless files, you need to give Spring Boot a hint so that it knows how to load them. You can do this by putting an extension hint in square brackets.

For example, suppose you have a `/etc/config/myconfig` file that you wish to import as yaml.
You can import it from your `application.properties` using the following:

-
Properties
-
YAML

`spring.config.import=file:/etc/config/myconfig[.yaml]````
spring:
 config:
 import: "file:/etc/config/myconfig[.yaml]"
```
This is the shorthand for:

-
Properties
-
YAML

`spring.config.import=file:/etc/config/myconfig[extension=.yaml]````
spring:
 config:
 import: "file:/etc/config/myconfig[extension=.yaml]"
```
### File attributes

The `spring.config.import` configuration property supports file attributes, as shown when specifying the encoding or the extension.

If you need to specify multiple attributes, you can use this syntax:

-
Properties
-
YAML

`spring.config.import=file:/etc/config/myconfig[extension=.yaml][encoding=utf-8]````
spring:
 config:
 import: "file:/etc/config/myconfig[extension=.yaml][encoding=utf-8]"
```
### Using Environment Variables

When running applications on a cloud platform (such as Kubernetes) you often need to read config values that the platform supplies. You can either use environment variables for such purpose, or you can use configuration trees.

You can even store whole configurations in properties or yaml format in (multiline) environment variables and load them using the `env:` prefix.
Assume there’s an environment variable called `MY_CONFIGURATION` with this content:

```
my.name=Service1
my.cluster=Cluster1
```
Using the `env:` prefix it is possible to import all properties from this variable:

-
Properties
-
YAML

`spring.config.import=env:MY_CONFIGURATION````
spring:
 config:
 import: "env:MY_CONFIGURATION"
```
| This feature also supports specifying the extension.
The default extension is `.properties`. |

### Using Configuration Trees

Storing config values in environment variables has drawbacks, especially if the value is supposed to be kept secret.

As an alternative to environment variables, many cloud platforms now allow you to map configuration into mounted data volumes.
For example, Kubernetes can volume mount both `ConfigMaps` and `Secrets`.

There are two common volume mount patterns that can be used:

-
A single file contains a complete set of properties (usually written as YAML).
-
Multiple files are written to a directory tree, with the filename becoming the ‘key’ and the contents becoming the ‘value’.

For the first case, you can import the YAML or Properties file directly using `spring.config.import` as described above.
For the second case, you need to use the `configtree:` prefix so that Spring Boot knows it needs to expose all the files as properties.

As an example, let’s imagine that Kubernetes has mounted the following volume:

```
etc/
 config/
 myapp/
 username
 password
```
The contents of the `username` file would be a config value, and the contents of `password` would be a secret.

To import these properties, you can add the following to your `application.properties` or `application.yaml` file:

-
Properties
-
YAML

`spring.config.import=optional:configtree:/etc/config/````
spring:
 config:
 import: "optional:configtree:/etc/config/"
```
You can then access or inject `myapp.username` and `myapp.password` properties from the `Environment` in the usual way.

| The names of the folders and files under the config tree form the property name.
In the above example, to access the properties as `username`and`password`, you can set`spring.config.import`to`optional:configtree:/etc/config/myapp`. |

| Filenames with dot notation are also correctly mapped.
For example, in the above example, a file named `myapp.username`in`/etc/config`would result in a`myapp.username`property in the`Environment`. |

| Configuration tree values can be bound to both string `String`and`byte[]`types depending on the contents expected. |

If you have multiple config trees to import from the same parent folder you can use a wildcard shortcut.
Any `configtree:` location that ends with `/*/` will import all immediate children as config trees.
As with a non-wildcard import, the names of the folders and files under each config tree form the property name.

For example, given the following volume:

```
etc/
 config/
 dbconfig/
 db/
 username
 password
 mqconfig/
 mq/
 username
 password
```
You can use `configtree:/etc/config/*/` as the import location:

-
Properties
-
YAML

`spring.config.import=optional:configtree:/etc/config/*/````
spring:
 config:
 import: "optional:configtree:/etc/config/*/"
```
This will add `db.username`, `db.password`, `mq.username` and `mq.password` properties.

| Directories loaded using a wildcard are sorted alphabetically. If you need a different order, then you should list each location as a separate import |

Configuration trees can also be used for Docker secrets.
When a Docker swarm service is granted access to a secret, the secret gets mounted into the container.
For example, if a secret named `db.password` is mounted at location `/run/secrets/`, you can make `db.password` available to the Spring environment using the following:

-
Properties
-
YAML

`spring.config.import=optional:configtree:/run/secrets/````
spring:
 config:
 import: "optional:configtree:/run/secrets/"
```
### Property Placeholders

The values in `application.properties` and `application.yaml` are filtered through the existing `Environment` when they are used, so you can refer back to previously defined values (for example, from System properties or environment variables).
The standard `${name}` property-placeholder syntax can be used anywhere within a value.
Property placeholders can also specify a default value using a `:` to separate the default value from the property name, for example `${name:default}`.

The use of placeholders with and without defaults is shown in the following example:

-
Properties
-
YAML

```
app.name=MyApp
app.description=${app.name} is a Spring Boot application written by ${username:Unknown}
```
```
app:
 name: "MyApp"
 description: "${app.name} is a Spring Boot application written by ${username:Unknown}"
```
Assuming that the `username` property has not been set elsewhere, `app.description` will have the value `MyApp is a Spring Boot application written by Unknown`.

| You should always refer to property names in the placeholder using their canonical form (kebab-case using only lowercase letters).
This will allow Spring Boot to use the same logic as it does when relaxed binding For example, |

| You can also use this technique to create “short” variants of existing Spring Boot properties. See the Use ‘Short’ Command Line Arguments section in “How-to Guides” for details. |

### Working With Multi-Document Files

Spring Boot allows you to split a single physical file into multiple logical documents which are each added independently. Documents are processed in order, from top to bottom. Later documents can override the properties defined in earlier ones.

For `application.yaml` files, the standard YAML multi-document syntax is used.
Three consecutive hyphens represent the end of one document, and the start of the next.

For example, the following file has two logical documents:

```
spring:
 application:
 name: "MyApp"
---
spring:
 application:
 name: "MyCloudApp"
 config:
 activate:
 on-cloud-platform: "kubernetes"
```
For `application.properties` files a special `#---` or `!---` comment is used to mark the document splits:

```
spring.application.name=MyApp
#---
spring.application.name=MyCloudApp
spring.config.activate.on-cloud-platform=kubernetes
```
| Property file separators must not have any leading whitespace and must have exactly three hyphen characters. The lines immediately before and after the separator must not be same comment prefix. |

| Multi-document property files are often used in conjunction with activation properties such as `spring.config.activate.on-profile`.
See the next section for details. |

| Multi-document property files cannot be loaded by using the `@PropertySource`or`@TestPropertySource`annotations. |

### Activation Properties

It is sometimes useful to only activate a given set of properties when certain conditions are met. For example, you might have properties that are only relevant when a specific profile is active.

You can conditionally activate a properties document using `spring.config.activate.*`.

The following activation properties are available:

| Property | Note |
|---|---|
|
 | A profile expression that must match for the document to be active, or a list of profile expressions of which at least one must match for the document to be active. |
|
 | The |

For example, the following specifies that the second document is only active when running on Kubernetes, and only when either the “prod” or “staging” profiles are active:

-
Properties
-
YAML

```
myprop=always-set
#---
spring.config.activate.on-cloud-platform=kubernetes
spring.config.activate.on-profile=prod | staging
myotherprop=sometimes-set
```
```
myprop:
 "always-set"
---
spring:
 config:
 activate:
 on-cloud-platform: "kubernetes"
 on-profile: "prod | staging"
myotherprop: "sometimes-set"
```
## Encrypting Properties

Spring Boot does not provide any built-in support for encrypting property values, however, it does provide the hook points necessary to modify values contained in the Spring `Environment`.
The `EnvironmentPostProcessor` interface allows you to manipulate the `Environment` before the application starts.
See Customize the Environment or ApplicationContext Before It Starts for details.

If you need a secure way to store credentials and passwords, the Spring Cloud Vault project provides support for storing externalized configuration in HashiCorp Vault.

## Working With YAML

YAML is a superset of JSON and, as such, is a convenient format for specifying hierarchical configuration data.
The `SpringApplication` class automatically supports YAML as an alternative to properties whenever you have the SnakeYAML library on your classpath.

| If you use starters, SnakeYAML is automatically provided by `spring-boot-starter`. |

### Mapping YAML to Properties

YAML documents need to be converted from their hierarchical format to a flat structure that can be used with the Spring `Environment`.
For example, consider the following YAML document:

```
environments:
 dev:
 url: "https://dev.example.com"
 name: "Developer Setup"
 prod:
 url: "https://another.example.com"
 name: "My Cool App"
```
In order to access these properties from the `Environment`, they would be flattened as follows:

```
environments.dev.url=https://dev.example.com
environments.dev.name=Developer Setup
environments.prod.url=https://another.example.com
environments.prod.name=My Cool App
```
Likewise, YAML lists also need to be flattened.
They are represented as property keys with `[index]` dereferencers.
For example, consider the following YAML:

```
 my:
 servers:
 - "dev.example.com"
 - "another.example.com"
```
The preceding example would be transformed into these properties:

```
my.servers[0]=dev.example.com
my.servers[1]=another.example.com
```
| Properties that use the `[index]`notation can be bound to Java`List`or`Set`objects using Spring Boot’s`Binder`class.
For more details see the Type-safe Configuration Properties section below. |

| YAML files cannot be loaded by using the `@PropertySource`or`@TestPropertySource`annotations.
So, in the case that you need to load values that way, you need to use a properties file. |

### Directly Loading YAML

Spring Framework provides two convenient classes that can be used to load YAML documents.
The `YamlPropertiesFactoryBean` loads YAML as `Properties` and the `YamlMapFactoryBean` loads YAML as a `Map`.

You can also use the `YamlPropertySourceLoader` class if you want to load YAML as a Spring `PropertySource`.

## Configuring Random Values

The `RandomValuePropertySource` is useful for injecting random values (for example, into secrets or test cases).
It can produce integers, longs, uuids, or strings, as shown in the following example:

-
Properties
-
YAML

```
my.secret=${random.value}
my.number=${random.int}
my.bignumber=${random.long}
my.uuid=${random.uuid}
my.number-less-than-ten=${random.int(10)}
my.number-in-range=${random.int[1024,65536]}
```
```
my:
 secret: "${random.value}"
 number: "${random.int}"
 bignumber: "${random.long}"
 uuid: "${random.uuid}"
 number-less-than-ten: "${random.int(10)}"
 number-in-range: "${random.int[1024,65536]}"
```
The `random.int*` syntax is `OPEN value (,max) CLOSE` where the `OPEN,CLOSE` are any character and `value,max` are integers.
If `max` is provided, then `value` is the minimum value and `max` is the maximum value (exclusive).

## Configuring System Environment Properties

Spring Boot supports setting a prefix for environment properties.
This is useful if the system environment is shared by multiple Spring Boot applications with different configuration requirements.
The prefix for system environment properties can be set directly on `SpringApplication` by calling the `setEnvironmentPrefix(…)` method before the application is run.

For example, if you set the prefix to `input`, a property such as `remote.timeout` will be resolved as `INPUT_REMOTE_TIMEOUT` in the system environment.

| The prefix onlyapplies to system environment properties.
The example above would continue to use`remote.timeout`when reading properties from other sources. |

## Type-safe Configuration Properties

Using the `@Value("${property}")` annotation to inject configuration properties can sometimes be cumbersome, especially if you are working with multiple properties or your data is hierarchical in nature.
Spring Boot provides an alternative method of working with properties that lets strongly typed beans govern and validate the configuration of your application.

| See also the differences between `@Value`and type-safe configuration properties. |

### JavaBean Properties Binding

It is possible to bind a bean declaring standard JavaBean properties as shown in the following example:

-
Java
-
Kotlin

```
import java.net.InetAddress;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import org.springframework.boot.context.properties.ConfigurationProperties;
@ConfigurationProperties("my.service")
public class MyProperties {
	private boolean enabled;
	private InetAddress remoteAddress;
	private final Security security = new Security();
	// getters / setters...
	public boolean isEnabled() {
 return this.enabled;
	}
	public void setEnabled(boolean enabled) {
 this.enabled = enabled;
	}
	public InetAddress getRemoteAddress() {
 return this.remoteAddress;
	}
	public void setRemoteAddress(InetAddress remoteAddress) {
 this.remoteAddress = remoteAddress;
	}
	public Security getSecurity() {
 return this.security;
	}
	public static class Security {
 private String username;
 private String password;
 private List<String> roles = new ArrayList<>(Collections.singleton("USER"));
 // getters / setters...
 public String getUsername() {
 return this.username;
 }
 public void setUsername(String username) {
 this.username = username;
 }
 public String getPassword() {
 return this.password;
 }
 public void setPassword(String password) {
 this.password = password;
 }
 public List<String> getRoles() {
 return this.roles;
 }
 public void setRoles(List<String> roles) {
 this.roles = roles;
 }
	}
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
import java.net.InetAddress
@ConfigurationProperties("my.service")
class MyProperties {
	var isEnabled = false
	var remoteAddress: InetAddress? = null
	val security = Security()
	class Security {
 var username: String? = null
 var password: String? = null
 var roles: List<String> = ArrayList(setOf("USER"))
	}
}
```
The preceding POJO defines the following properties:

-
`my.service.enabled`, with a value of`false`by default.
-
`my.service.remote-address`, with a type that can be coerced from`String`.
-
`my.service.security.username`, with a nested "security" object whose name is determined by the name of the property. In particular, the type is not used at all there and could have been`SecurityProperties`.
-
`my.service.security.password`.
-
`my.service.security.roles`, with a collection of`String`that defaults to`USER`.

| To use a reserved keyword in the name of a property, such as `my.service.import`, use the`@Name`annotation on the property’s field. |

| The properties that map to `@ConfigurationProperties`classes available in Spring Boot, which are configured through properties files, YAML files, environment variables, and other mechanisms, are public API but the accessors (getters/setters) of the class itself are not meant to be used directly. |

| Such arrangement relies on a default empty constructor and getters and setters are usually mandatory, since binding is through standard Java Beans property descriptors, just like in Spring MVC. A setter may be omitted in the following cases:
 Some people use Project Lombok to add getters and setters automatically. Make sure that Lombok does not generate any particular constructor for such a type, as it is used automatically by the container to instantiate the object. Finally, only standard Java Bean properties are considered and binding on static properties is not supported. |

### Constructor Binding

The example in the previous section can be rewritten in an immutable fashion as shown in the following example:

-
Java
-
Kotlin

```
import java.net.InetAddress;
import java.util.List;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.context.properties.bind.DefaultValue;
@ConfigurationProperties("my.service")
public class MyProperties {
	// fields...
	private final boolean enabled;
	private final InetAddress remoteAddress;
	private final Security security;
	public MyProperties(boolean enabled, InetAddress remoteAddress, Security security) {
 this.enabled = enabled;
 this.remoteAddress = remoteAddress;
 this.security = security;
	}
	// getters...
	public boolean isEnabled() {
 return this.enabled;
	}
	public InetAddress getRemoteAddress() {
 return this.remoteAddress;
	}
	public Security getSecurity() {
 return this.security;
	}
	public static class Security {
 // fields...
 private final String username;
 private final String password;
 private final List<String> roles;
 public Security(String username, String password, @DefaultValue("USER") List<String> roles) {
 this.username = username;
 this.password = password;
 this.roles = roles;
 }
 // getters...
 public String getUsername() {
 return this.username;
 }
 public String getPassword() {
 return this.password;
 }
 public List<String> getRoles() {
 return this.roles;
 }
	}
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.bind.DefaultValue
import java.net.InetAddress
@ConfigurationProperties("my.service")
class MyProperties(val enabled: Boolean, val remoteAddress: InetAddress,
 val security: Security) {
	class Security(val username: String, val password: String,
 @param:DefaultValue("USER") val roles: List<String>)
}
```
In this setup, the presence of a single parameterized constructor implies that constructor binding should be used.
This means that the binder will find a constructor with the parameters that you wish to have bound.
If your class has multiple constructors, the `@ConstructorBinding` annotation can be used to specify which constructor to use for constructor binding.

To opt-out of constructor binding for a class, the parameterized constructor must be annotated with `@Autowired` or made `private`.
Kotlin developers can use an empty primary constructor to opt-out of constructor binding.

For example:

-
Java
-
Kotlin

```
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.context.properties.ConfigurationProperties;
@ConfigurationProperties("my")
public class MyProperties {
	// fields...
	final MyBean myBean;
	private String name;
	@Autowired
	public MyProperties(MyBean myBean) {
 this.myBean = myBean;
	}
	// getters / setters...
	public String getName() {
 return this.name;
	}
	public void setName(String name) {
 this.name = name;
	}
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
@ConfigurationProperties("my")
class MyProperties() {
	constructor(name: String) : this() {
 this.name = name
	}
	// vars...
	var name: String? = null
}
```
Constructor binding can be used with records.
Unless your record has multiple constructors, there is no need to use `@ConstructorBinding`.

Nested members of a constructor bound class (such as `Security` in the example above) will also be bound through their constructor.

| To use constructor binding the class must be enabled using `@EnableConfigurationProperties`or configuration property scanning.
You cannot use constructor binding with beans that are created by the regular Spring mechanisms (for example`@Component`beans, beans created by using`@Bean`methods or beans loaded by using`@Import`) |

| To use constructor binding the class must be compiled with `-parameters`.
This will happen automatically if you use Spring Boot’s Gradle plugin or if you use Maven and`spring-boot-starter-parent`. |

| The use of `Optional`with`@ConfigurationProperties`is not recommended as it is primarily intended for use as a return type.
As such, it is not well-suited to configuration property injection.
For consistency with properties of other types, if you do declare an`Optional`property and it has no value,`null`rather than an empty`Optional`will be bound. |

| To use a reserved keyword in the name of a property, such as `my.service.import`, use the`@Name`annotation on the constructor parameter. |

#### @DefaultValue and Binding

Default values can be specified using `@DefaultValue` on constructor parameters and record components.
The conversion service will be applied to coerce the annotation’s `String` value to the target type of a missing property.

In the `MyProperties` example above, you can see the nested `Security` class uses `@DefaultValue("USER")` for the `roles` parameter.
This means that if `security` properties are defined, but `roles` is not, the default of `"USER"` will be bound.

For example, the following properties:

-
Properties
-
YAML

```
my.service.enabled=true
my.service.security.username=admin
```
```
my:
 service:
 enabled: true
 security:
 username: admin
```
Will be bound as `new MyProperties(true, null, new Security("admin", null, List.of("USER")))`

If the `security` property is not present at all, then the `Security` instance will be `null`.

For example, the following properties:

-
Properties
-
YAML

`my.service.enabled=true````
my:
 service:
 enabled: true
```
Will be bound as `new MyProperties(true, null, null)`

| You can define an empty For YAML, you can use the following syntax: With Will be bound as |

If you want to always bind a non-null instance of `Security`, even when properties are missing, you can use an empty `@DefaultValue` annotation:

-
Java
-
Kotlin

```
	public MyProperties(boolean enabled, InetAddress remoteAddress, @DefaultValue Security security) {
 this.enabled = enabled;
 this.remoteAddress = remoteAddress;
 this.security = security;
	}
```
```
class MyProperties(val enabled: Boolean, val remoteAddress: InetAddress,
 @DefaultValue val security: Security) {
	class Security(val username: String?, val password: String?,
 @param:DefaultValue("USER") val roles: List<String>)
}
```
| When using Kotlin, you will need to declare the `username`and`password`parameters as nullable since they do not have default values |

### Default Values

Default Values defined in configuration properties are not reflected in the `Environment`.
In the examples above, the `enabled` property of `MyProperties` bound to `my.service` is `false` by default.

However, `my.service.enabled` is not available in the `Environment` with a value of `false` if no such property is set by the user.
Concretely, this prevents you to use `@Value(${"my.service.enabled"})` or `my.service.enabled` as a placeholder in configuration properties without explicitly providing a default.
If you need to query the `Environment` for that property, for instance in a `Condition` implementation, the default needs to be provided as well.

### Enabling @ConfigurationProperties-annotated Types

Spring Boot provides infrastructure to bind `@ConfigurationProperties` types and register them as beans.
You can either enable configuration properties on a class-by-class basis or enable configuration property scanning that works in a similar manner to component scanning.

Sometimes, classes annotated with `@ConfigurationProperties` might not be suitable for scanning, for example, if you’re developing your own auto-configuration or you want to enable them conditionally.
In these cases, specify the list of types to process using the `@EnableConfigurationProperties` annotation.
This can be done on any `@Configuration` class, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
@EnableConfigurationProperties(SomeProperties.class)
public class MyConfiguration {
}
```
```
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.context.annotation.Configuration
@Configuration(proxyBeanMethods = false)
@EnableConfigurationProperties(SomeProperties::class)
class MyConfiguration
```
-
Java
-
Kotlin

```
import org.springframework.boot.context.properties.ConfigurationProperties;
@ConfigurationProperties("some.properties")
public class SomeProperties {
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
@ConfigurationProperties("some.properties")
class SomeProperties
```
To use configuration property scanning, add the `@ConfigurationPropertiesScan` annotation to your application.
Typically, it is added to the main application class that is annotated with `@SpringBootApplication` but it can be added to any `@Configuration` class.
By default, scanning will occur from the package of the class that declares the annotation.
If you want to define specific packages to scan, you can do so as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.ConfigurationPropertiesScan;
@SpringBootApplication
@ConfigurationPropertiesScan({ "com.example.app", "com.example.another" })
public class MyApplication {
}
```
```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.context.properties.ConfigurationPropertiesScan
@SpringBootApplication
@ConfigurationPropertiesScan("com.example.app", "com.example.another")
class MyApplication
```
| When the Assuming that it is in the |

We recommend that `@ConfigurationProperties` only deal with the environment and, in particular, does not inject other beans from the context.
For corner cases, setter injection can be used or any of the `*Aware` interfaces provided by the framework (such as `EnvironmentAware` if you need access to the `Environment`).
If you still want to inject other beans using the constructor, the configuration properties bean must be annotated with `@Component` and use JavaBean-based property binding.

### Using @ConfigurationProperties-annotated Types

This style of configuration works particularly well with the `SpringApplication` external YAML configuration, as shown in the following example:

```
my:
 service:
 remote-address: 192.168.1.1
 security:
 username: "admin"
 roles:
 - "USER"
 - "ADMIN"
```
To work with `@ConfigurationProperties` beans, you can inject them in the same way as any other bean, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.stereotype.Service;
@Service
public class MyService {
	private final MyProperties properties;
	public MyService(MyProperties properties) {
 this.properties = properties;
	}
	public void openConnection() {
 Server server = new Server(this.properties.getRemoteAddress());
 server.start();
 // ...
	}
	// ...
}
```
```
import org.springframework.stereotype.Service
@Service
class MyService(val properties: MyProperties) {
	fun openConnection() {
 val server = Server(properties.remoteAddress)
 server.start()
 // ...
	}
	// ...
}
```
| Using `@ConfigurationProperties`also lets you generate metadata files that can be used by IDEs to offer auto-completion for your own keys.
See the appendix for details. |

### Third-party Configuration

As well as using `@ConfigurationProperties` to annotate a class, you can also use it on public `@Bean` methods.
Doing so can be particularly useful when you want to bind properties to third-party components that are outside of your control.

To configure a bean from the `Environment` properties, add `@ConfigurationProperties` to its bean registration, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class ThirdPartyConfiguration {
	@Bean
	@ConfigurationProperties("another")
	public AnotherComponent anotherComponent() {
 return new AnotherComponent();
	}
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
@Configuration(proxyBeanMethods = false)
class ThirdPartyConfiguration {
	@Bean
	@ConfigurationProperties("another")
	fun anotherComponent(): AnotherComponent = AnotherComponent()
}
```
Any JavaBean property defined with the `another` prefix is mapped onto that `AnotherComponent` bean in manner similar to the preceding `SomeProperties` example.

### Relaxed Binding

Spring Boot uses some relaxed rules for binding `Environment` properties to `@ConfigurationProperties` beans, so there does not need to be an exact match between the `Environment` property name and the bean property name.
Common examples where this is useful include dash-separated environment properties (for example, `context-path` binds to `contextPath`), and capitalized environment properties (for example, `PORT` binds to `port`).

As an example, consider the following `@ConfigurationProperties` class:

-
Java
-
Kotlin

```
import org.springframework.boot.context.properties.ConfigurationProperties;
@ConfigurationProperties("my.main-project.person")
public class MyPersonProperties {
	private String firstName;
	public String getFirstName() {
 return this.firstName;
	}
	public void setFirstName(String firstName) {
 this.firstName = firstName;
	}
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
@ConfigurationProperties("my.main-project.person")
class MyPersonProperties {
	var firstName: String? = null
}
```
With the preceding code, the following properties names can all be used:

| Property | Note |
|---|---|
|
 | Kebab case, which is recommended for use in |
|
 | Standard camel case syntax. |
|
 | Underscore notation, which is an alternative format for use in |
|
 | Upper case format, which is recommended when using system environment variables. |

| The `prefix`value for the annotationmustbe in kebab case (lowercase and separated by`-`, such as`my.main-project.person`). |

| Property Source | Simple | List |
|---|---|---|
| Properties Files | Camel case, kebab case, or underscore notation | Standard list syntax using |
| YAML Files | Camel case, kebab case, or underscore notation | Standard YAML list syntax or comma-separated values |
| Environment Variables | Upper case format with underscore as the delimiter (see Binding From Environment Variables). | Numeric values surrounded by underscores (see Binding From Environment Variables) |
| System properties | Camel case, kebab case, or underscore notation | Standard list syntax using |

| We recommend that, when possible, properties are stored in lower-case kebab format, such as `my.person.first-name=Rod`. |

#### Binding Maps

When binding to `Map` properties you may need to use a special bracket notation so that the original `key` value is preserved.
If the key is not surrounded by `[]`, any characters that are not alpha-numeric, `-` or `.` are removed.

For example, consider binding the following properties to a `Map<String,String>`:

-
Properties
-
YAML

```
my.map[/key1]=value1
my.map[/key2]=value2
my.map./key3=value3
```
```
my:
 map:
 "[/key1]": "value1"
 "[/key2]": "value2"
 "/key3": "value3"
```
| For YAML files, the brackets need to be surrounded by quotes for the keys to be parsed properly. |

The properties above will bind to a `Map` with `/key1`, `/key2` and `key3` as the keys in the map.
The slash has been removed from `key3` because it was not surrounded by square brackets.

When binding to scalar values, keys with `.` in them do not need to be surrounded by `[]`.
Scalar values include enums and all types in the `java.lang` package except for `Object`.
Binding `a.b=c` to `Map<String, String>` will preserve the `.` in the key and return a Map with the entry `{"a.b"="c"}`.
For any other types you need to use the bracket notation if your `key` contains a `.`.
For example, binding `a.b=c` to `Map<String, Object>` will return a Map with the entry `{"a"={"b"="c"}}` whereas `[a.b]=c` will return a Map with the entry `{"a.b"="c"}`.

#### Binding From Environment Variables

Most operating systems impose strict rules around the names that can be used for environment variables.
For example, Linux shell variables can contain only letters (`a` to `z` or `A` to `Z`), numbers (`0` to `9`) or the underscore character (`_`).
By convention, Unix shell variables will also have their names in UPPERCASE.

Spring Boot’s relaxed binding rules are, as much as possible, designed to be compatible with these naming restrictions.

To convert a property name in the canonical-form to an environment variable name you can follow these rules:

-
Replace dots ( `.`) with underscores (`_`).
-
Remove any dashes ( `-`).
-
Convert to uppercase.

For example, the configuration property `spring.main.log-startup-info` would be an environment variable named `SPRING_MAIN_LOGSTARTUPINFO`.

Environment variables can also be used when binding to object lists.
To bind to a `List`, the element number should be surrounded with underscores in the variable name.

For example, the configuration property `my.service[0].other` would use an environment variable named `MY_SERVICE_0_OTHER`.

Support for binding from environment variables is applied to the `systemEnvironment` property source and to any additional property source whose name ends with `-systemEnvironment`.

#### Binding Maps From Environment Variables

When Spring Boot binds an environment variable to a property class, it lowercases the environment variable name before binding.
Most of the time this detail isn’t important, except when binding to `Map` properties.

The keys in the `Map` are always in lowercase, as seen in the following example:

-
Java
-
Kotlin

```
import java.util.HashMap;
import java.util.Map;
import org.springframework.boot.context.properties.ConfigurationProperties;
@ConfigurationProperties("my.props")
public class MyMapsProperties {
	private final Map<String, String> values = new HashMap<>();
	public Map<String, String> getValues() {
 return this.values;
	}
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
@ConfigurationProperties("my.props")
class MyMapsProperties {
	val values: Map<String, String> = HashMap()
}
```
When setting `MY_PROPS_VALUES_KEY=value`, the `values` `Map` contains a `{"key"="value"}` entry.

Only the environment variable **name** is lower-cased, not the value.
When setting `MY_PROPS_VALUES_KEY=VALUE`, the `values` `Map` contains a `{"key"="VALUE"}` entry.

#### Caching

Relaxed binding uses a cache to improve performance. By default, this caching is only applied to immutable property sources.
To customize this behavior, for example to enable caching for mutable property sources, use `ConfigurationPropertyCaching`.

### Merging Complex Types

When lists are configured in more than one place, overriding works by replacing the entire list.

For example, assume a `MyPojo` object with `name` and `description` attributes that are `null` by default.
The following example exposes a list of `MyPojo` objects from `MyProperties`:

-
Java
-
Kotlin

```
import java.util.ArrayList;
import java.util.List;
import org.springframework.boot.context.properties.ConfigurationProperties;
@ConfigurationProperties("my")
public class MyProperties {
	private final List<MyPojo> list = new ArrayList<>();
	public List<MyPojo> getList() {
 return this.list;
	}
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
@ConfigurationProperties("my")
class MyProperties {
	val list: List<MyPojo> = ArrayList()
}
```
Consider the following configuration:

-
Properties
-
YAML

```
my.list[0].name=my name
my.list[0].description=my description
#---
spring.config.activate.on-profile=dev
my.list[0].name=my another name
```
```
my:
 list:
 - name: "my name"
 description: "my description"
---
spring:
 config:
 activate:
 on-profile: "dev"
my:
 list:
 - name: "my another name"
```
If the `dev` profile is not active, `MyProperties.list` contains one `MyPojo` entry, as previously defined.
If the `dev` profile is enabled, however, the `list` *still* contains only one entry (with a name of `my another name` and a description of `null`).
This configuration *does not* add a second `MyPojo` instance to the list, and it does not merge the items.

When a `List` is specified in multiple profiles, the one with the highest priority (and only that one) is used.
Consider the following example:

-
Properties
-
YAML

```
my.list[0].name=my name
my.list[0].description=my description
my.list[1].name=another name
my.list[1].description=another description
#---
spring.config.activate.on-profile=dev
my.list[0].name=my another name
```
```
my:
 list:
 - name: "my name"
 description: "my description"
 - name: "another name"
 description: "another description"
---
spring:
 config:
 activate:
 on-profile: "dev"
my:
 list:
 - name: "my another name"
```
In the preceding example, if the `dev` profile is active, `MyProperties.list` contains *one* `MyPojo` entry (with a name of `my another name` and a description of `null`).
For YAML, both comma-separated lists and YAML lists can be used for completely overriding the contents of the list.

For `Map` properties, you can bind with property values drawn from multiple sources.
However, for the same property in multiple sources, the one with the highest priority is used.
The following example exposes a `Map<String, MyPojo>` from `MyProperties`:

-
Java
-
Kotlin

```
import java.util.LinkedHashMap;
import java.util.Map;
import org.springframework.boot.context.properties.ConfigurationProperties;
@ConfigurationProperties("my")
public class MyProperties {
	private final Map<String, MyPojo> map = new LinkedHashMap<>();
	public Map<String, MyPojo> getMap() {
 return this.map;
	}
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
@ConfigurationProperties("my")
class MyProperties {
	val map: Map<String, MyPojo> = LinkedHashMap()
}
```
Consider the following configuration:

-
Properties
-
YAML

```
my.map.key1.name=my name 1
my.map.key1.description=my description 1
#---
spring.config.activate.on-profile=dev
my.map.key1.name=dev name 1
my.map.key2.name=dev name 2
my.map.key2.description=dev description 2
```
```
my:
 map:
 key1:
 name: "my name 1"
 description: "my description 1"
---
spring:
 config:
 activate:
 on-profile: "dev"
my:
 map:
 key1:
 name: "dev name 1"
 key2:
 name: "dev name 2"
 description: "dev description 2"
```
If the `dev` profile is not active, `MyProperties.map` contains one entry with key `key1` (with a name of `my name 1` and a description of `my description 1`).
If the `dev` profile is enabled, however, `map` contains two entries with keys `key1` (with a name of `dev name 1` and a description of `my description 1`) and `key2` (with a name of `dev name 2` and a description of `dev description 2`).

| The preceding merging rules apply to properties from all property sources, and not just files. |

### Properties Conversion

Spring Boot attempts to coerce the external application properties to the right type when it binds to the `@ConfigurationProperties` beans.
If you need custom type conversion, you can provide a `ConversionService` bean (with a bean named `conversionService`) or custom property editors (through a `CustomEditorConfigurer` bean) or custom converters (with bean definitions annotated as `@ConfigurationPropertiesBinding`).

| Beans used for property conversion are requested very early during the application lifecycle so make sure to limit the dependencies that your |

| You may want to rename your custom `ConversionService`if it is not required for configuration keys coercion and only rely on custom converters qualified with`@ConfigurationPropertiesBinding`.
When qualifying a`@Bean`method with`@ConfigurationPropertiesBinding`, the method should be`static`to avoid “bean is not eligible for getting processed by all BeanPostProcessors” warnings. |

#### Converting Durations

Spring Boot has dedicated support for expressing durations.
If you expose a `Duration` property, the following formats in application properties are available:

-
A regular `long`representation (using milliseconds as the default unit unless a`@DurationUnit`has been specified)
-
A more readable format where the value and the unit are coupled ( `10s`means 10 seconds)

Consider the following example:

-
Java
-
Kotlin

```
import java.time.Duration;
import java.time.temporal.ChronoUnit;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.convert.DurationUnit;
@ConfigurationProperties("my")
public class MyProperties {
	@DurationUnit(ChronoUnit.SECONDS)
	private Duration sessionTimeout = Duration.ofSeconds(30);
	private Duration readTimeout = Duration.ofMillis(1000);
	// getters / setters...
	public Duration getSessionTimeout() {
 return this.sessionTimeout;
	}
	public void setSessionTimeout(Duration sessionTimeout) {
 this.sessionTimeout = sessionTimeout;
	}
	public Duration getReadTimeout() {
 return this.readTimeout;
	}
	public void setReadTimeout(Duration readTimeout) {
 this.readTimeout = readTimeout;
	}
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.convert.DurationUnit
import java.time.Duration
import java.time.temporal.ChronoUnit
@ConfigurationProperties("my")
class MyProperties {
	@DurationUnit(ChronoUnit.SECONDS)
	var sessionTimeout = Duration.ofSeconds(30)
	var readTimeout = Duration.ofMillis(1000)
}
```
To specify a session timeout of 30 seconds, `30`, `PT30S` and `30s` are all equivalent.
A read timeout of 500ms can be specified in any of the following form: `500`, `PT0.5S` and `500ms`.

You can also use any of the supported units. These are:

-
`ns`for nanoseconds
-
`us`for microseconds
-
`ms`for milliseconds
-
`s`for seconds
-
`m`for minutes
-
`h`for hours
-
`d`for days

The default unit is milliseconds and can be overridden using `@DurationUnit` as illustrated in the sample above.

If you prefer to use constructor binding, the same properties can be exposed, as shown in the following example:

-
Java
-
Kotlin

```
import java.time.Duration;
import java.time.temporal.ChronoUnit;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.context.properties.bind.DefaultValue;
import org.springframework.boot.convert.DurationUnit;
@ConfigurationProperties("my")
public class MyProperties {
	// fields...
	private final Duration sessionTimeout;
	private final Duration readTimeout;
	public MyProperties(@DurationUnit(ChronoUnit.SECONDS) @DefaultValue("30s") Duration sessionTimeout,
 @DefaultValue("1000ms") Duration readTimeout) {
 this.sessionTimeout = sessionTimeout;
 this.readTimeout = readTimeout;
	}
	// getters...
	public Duration getSessionTimeout() {
 return this.sessionTimeout;
	}
	public Duration getReadTimeout() {
 return this.readTimeout;
	}
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.bind.DefaultValue
import org.springframework.boot.convert.DurationUnit
import java.time.Duration
import java.time.temporal.ChronoUnit
@ConfigurationProperties("my")
class MyProperties(@param:DurationUnit(ChronoUnit.SECONDS) @param:DefaultValue("30s") val sessionTimeout: Duration,
 @param:DefaultValue("1000ms") val readTimeout: Duration)
```
| If you are upgrading a `Long`property, make sure to define the unit (using`@DurationUnit`) if it is not milliseconds.
Doing so gives a transparent upgrade path while supporting a much richer format. |

#### Converting Periods

In addition to durations, Spring Boot can also work with `Period` type.
The following formats can be used in application properties:

-
A regular `int`representation (using days as the default unit unless a`@PeriodUnit`has been specified)
-
A simpler format where the value and the unit pairs are coupled ( `1y3d`means 1 year and 3 days)

The following units are supported with the simple format:

-
`y`for years
-
`m`for months
-
`w`for weeks
-
`d`for days

| The `Period`type never actually stores the number of weeks, it is a shortcut that means “7 days”. |

#### Converting Data Sizes

Spring Framework has a `DataSize` value type that expresses a size in bytes.
If you expose a `DataSize` property, the following formats in application properties are available:

-
A regular `long`representation (using bytes as the default unit unless a`@DataSizeUnit`has been specified)
-
A more readable format where the value and the unit are coupled ( `10MB`means 10 megabytes)

Consider the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.convert.DataSizeUnit;
import org.springframework.util.unit.DataSize;
import org.springframework.util.unit.DataUnit;
@ConfigurationProperties("my")
public class MyProperties {
	@DataSizeUnit(DataUnit.MEGABYTES)
	private DataSize bufferSize = DataSize.ofMegabytes(2);
	private DataSize sizeThreshold = DataSize.ofBytes(512);
	// getters/setters...
	public DataSize getBufferSize() {
 return this.bufferSize;
	}
	public void setBufferSize(DataSize bufferSize) {
 this.bufferSize = bufferSize;
	}
	public DataSize getSizeThreshold() {
 return this.sizeThreshold;
	}
	public void setSizeThreshold(DataSize sizeThreshold) {
 this.sizeThreshold = sizeThreshold;
	}
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.convert.DataSizeUnit
import org.springframework.util.unit.DataSize
import org.springframework.util.unit.DataUnit
@ConfigurationProperties("my")
class MyProperties {
	@DataSizeUnit(DataUnit.MEGABYTES)
	var bufferSize = DataSize.ofMegabytes(2)
	var sizeThreshold = DataSize.ofBytes(512)
}
```
To specify a buffer size of 10 megabytes, `10` and `10MB` are equivalent.
A size threshold of 256 bytes can be specified as `256` or `256B`.

You can also use any of the supported units. These are:

-
`B`for bytes
-
`KB`for kilobytes
-
`MB`for megabytes
-
`GB`for gigabytes
-
`TB`for terabytes

The default unit is bytes and can be overridden using `@DataSizeUnit` as illustrated in the sample above.

If you prefer to use constructor binding, the same properties can be exposed, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.context.properties.bind.DefaultValue;
import org.springframework.boot.convert.DataSizeUnit;
import org.springframework.util.unit.DataSize;
import org.springframework.util.unit.DataUnit;
@ConfigurationProperties("my")
public class MyProperties {
	// fields...
	private final DataSize bufferSize;
	private final DataSize sizeThreshold;
	public MyProperties(@DataSizeUnit(DataUnit.MEGABYTES) @DefaultValue("2MB") DataSize bufferSize,
 @DefaultValue("512B") DataSize sizeThreshold) {
 this.bufferSize = bufferSize;
 this.sizeThreshold = sizeThreshold;
	}
	// getters...
	public DataSize getBufferSize() {
 return this.bufferSize;
	}
	public DataSize getSizeThreshold() {
 return this.sizeThreshold;
	}
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.bind.DefaultValue
import org.springframework.boot.convert.DataSizeUnit
import org.springframework.util.unit.DataSize
import org.springframework.util.unit.DataUnit
@ConfigurationProperties("my")
class MyProperties(@param:DataSizeUnit(DataUnit.MEGABYTES) @param:DefaultValue("2MB") val bufferSize: DataSize,
 @param:DefaultValue("512B") val sizeThreshold: DataSize)
```
| If you are upgrading a `Long`property, make sure to define the unit (using`@DataSizeUnit`) if it is not bytes.
Doing so gives a transparent upgrade path while supporting a much richer format. |

#### Converting Base64 Data

Spring Boot supports resolving binary data that have been Base64 encoded.
If you expose a `Resource` property, the base64 encoded text can be provided as the value with a `base64:` prefix, as shown in the following example:

-
Properties
-
YAML

`my.property=base64:SGVsbG8gV29ybGQ=````
my:
 property: base64:SGVsbG8gV29ybGQ=
```
| The `Resource`property can also be used to provide the path to the resource, making it more versatile. |

### @ConfigurationProperties Validation

Spring Boot attempts to validate `@ConfigurationProperties` classes whenever they are annotated with Spring’s `@Validated` annotation.
You can use JSR-303 `jakarta.validation` constraint annotations directly on your configuration class.
To do so, ensure that a compliant JSR-303 implementation is on your classpath and then add constraint annotations to your fields, as shown in the following example:

-
Java
-
Kotlin

```
import java.net.InetAddress;
import jakarta.validation.constraints.NotNull;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.validation.annotation.Validated;
@ConfigurationProperties("my.service")
@Validated
public class MyProperties {
	@NotNull
	private InetAddress remoteAddress;
	// getters/setters...
	public InetAddress getRemoteAddress() {
 return this.remoteAddress;
	}
	public void setRemoteAddress(InetAddress remoteAddress) {
 this.remoteAddress = remoteAddress;
	}
}
```
```
import jakarta.validation.constraints.NotNull
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.validation.annotation.Validated
import java.net.InetAddress
@ConfigurationProperties("my.service")
@Validated
class MyProperties {
	var remoteAddress: @NotNull InetAddress? = null
}
```
| You can also trigger validation by annotating the `@Bean`method that creates the configuration properties with`@Validated`. |

To cascade validation to nested properties the associated field must be annotated with `@Valid`.
The following example builds on the preceding `MyProperties` example:

-
Java
-
Kotlin

```
import java.net.InetAddress;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.validation.annotation.Validated;
@ConfigurationProperties("my.service")
@Validated
public class MyProperties {
	@NotNull
	private InetAddress remoteAddress;
	@Valid
	private final Security security = new Security();
	// getters/setters...
	public InetAddress getRemoteAddress() {
 return this.remoteAddress;
	}
	public void setRemoteAddress(InetAddress remoteAddress) {
 this.remoteAddress = remoteAddress;
	}
	public Security getSecurity() {
 return this.security;
	}
	public static class Security {
 @NotEmpty
 private String username;
 // getters/setters...
 public String getUsername() {
 return this.username;
 }
 public void setUsername(String username) {
 this.username = username;
 }
	}
}
```
```
import jakarta.validation.Valid
import jakarta.validation.constraints.NotEmpty
import jakarta.validation.constraints.NotNull
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.validation.annotation.Validated
import java.net.InetAddress
@ConfigurationProperties("my.service")
@Validated
class MyProperties {
	var remoteAddress: @NotNull InetAddress? = null
	@Valid
	val security = Security()
	class Security {
 @NotEmpty
 var username: String? = null
	}
}
```
You can also add a custom Spring `Validator` by creating a bean definition called `configurationPropertiesValidator`.
The `@Bean` method should be declared `static`.
The configuration properties validator is created very early in the application’s lifecycle, and declaring the `@Bean` method as static lets the bean be created without having to instantiate the `@Configuration` class.
Doing so avoids any problems that may be caused by early instantiation.

| The `spring-boot-actuator`module includes an endpoint that exposes all`@ConfigurationProperties`beans.
Point your web browser to`/actuator/configprops`or use the equivalent JMX endpoint.
See the Production ready features section for details. |

### @ConfigurationProperties vs. @Value

The `@Value` annotation is a core container feature, and it does not provide the same features as type-safe configuration properties.
The following table summarizes the features that are supported by `@ConfigurationProperties` and `@Value`:

| Feature | `@ConfigurationProperties` | `@Value` |
|---|---|---|
| Yes | Limited (see note below) | |
| Yes | No | |
|
 | No | Yes |

| If you do want to use For example, |

If you define a set of configuration keys for your own components, we recommend you group them in a POJO annotated with `@ConfigurationProperties`.
Doing so will provide you with structured, type-safe object that you can inject into your own beans.

`SpEL` expressions from application property files are not processed at time of parsing these files and populating the environment.
However, it is possible to write a `SpEL` expression in `@Value`.
If the value of a property from an application property file is a `SpEL` expression, it will be evaluated when consumed through `@Value`.

# Profiles

Spring Profiles provide a way to segregate parts of your application configuration and make it be available only in certain environments.
Any `@Component`, `@Configuration` or `@ConfigurationProperties` can be marked with `@Profile` to limit when it is loaded, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Profile;
@Configuration(proxyBeanMethods = false)
@Profile("production")
public class ProductionConfiguration {
	// ...
}
```
```
import org.springframework.context.annotation.Configuration
import org.springframework.context.annotation.Profile
@Configuration(proxyBeanMethods = false)
@Profile("production")
class ProductionConfiguration {
	// ...
}
```
| If `@ConfigurationProperties`beans are registered through`@EnableConfigurationProperties`instead of automatic scanning, the`@Profile`annotation needs to be specified on the`@Configuration`class that has the`@EnableConfigurationProperties`annotation.
In the case where`@ConfigurationProperties`are scanned,`@Profile`can be specified on the`@ConfigurationProperties`class itself. |

You can use a `spring.profiles.active` `Environment` property to specify which profiles are active.
You can specify the property in any of the ways described earlier in this chapter.
For example, you could include it in your `application.properties`, as shown in the following example:

-
Properties
-
YAML

`spring.profiles.active=dev,hsqldb````
spring:
 profiles:
 active: "dev,hsqldb"
```
You could also specify it on the command line by using the following switch: `--spring.profiles.active=dev,hsqldb`.

If no profile is active, a default profile is enabled.
The name of the default profile is `default` and it can be tuned using the `spring.profiles.default` `Environment` property, as shown in the following example:

-
Properties
-
YAML

`spring.profiles.default=none````
spring:
 profiles:
 default: "none"
```
`spring.profiles.active` and `spring.profiles.default` can only be used in non-profile-specific documents.
This means they cannot be included in profile specific files or documents activated by `spring.config.activate.on-profile`.

For example, the second document configuration is invalid:

-
Properties
-
YAML

```
spring.profiles.active=prod
#---
spring.config.activate.on-profile=prod
spring.profiles.active=metrics
```
```
# this document is valid
spring:
 profiles:
 active: "prod"
---
# this document is invalid
spring:
 config:
 activate:
 on-profile: "prod"
 profiles:
 active: "metrics"
```
The `spring.profiles.active` property follows the same ordering rules as other properties.
The highest `PropertySource` wins.
This means that you can specify active profiles in `application.properties` and then **replace** them by using the command line switch.

| See the “Externalized Configuration” for more details on the order in which property sources are considered. |

| By default, profile names in Spring Boot may contain letters, numbers, or permitted characters ( This restriction helps to prevent common parsing issues.
if, however, you prefer more flexible profile names you can set
 |

## Adding Active Profiles

Sometimes, it is useful to have properties that **add** to the active profiles rather than replace them.
The `spring.profiles.include` property can be used to add active profiles on top of those activated by the `spring.profiles.active` property.
The `SpringApplication` entry point also has a Java API for setting additional profiles.
See the `setAdditionalProfiles()` method in `SpringApplication`.

For example, when an application with the following properties is run, the common and local profiles will be activated even when it runs using the `--spring.profiles.active` switch:

-
Properties
-
YAML

```
spring.profiles.include[0]=common
spring.profiles.include[1]=local
```
```
spring:
 profiles:
 include:
 - "common"
 - "local"
```
| Included profiles are added before any `spring.profiles.active`profiles. |

| The `spring.profiles.include`property is processed for each property source, as such the usual complex type merging rules for lists do not apply. |

| Similar to `spring.profiles.active`,`spring.profiles.include`can only be used in non-profile-specific documents.
This means it cannot be included in profile specific files or documents activated by`spring.config.activate.on-profile`. |

Profile groups, which are described in the next section can also be used to add active profiles if a given profile is active.

## Profile Groups

Occasionally the profiles that you define and use in your application are too fine-grained and become cumbersome to use.
For example, you might have `proddb` and `prodmq` profiles that you use to enable database and messaging features independently.

To help with this, Spring Boot lets you define profile groups. A profile group allows you to define a logical name for a related group of profiles.

For example, we can create a `production` group that consists of our `proddb` and `prodmq` profiles.

-
Properties
-
YAML

```
spring.profiles.group.production[0]=proddb
spring.profiles.group.production[1]=prodmq
```
```
spring:
 profiles:
 group:
 production:
 - "proddb"
 - "prodmq"
```
Our application can now be started using `--spring.profiles.active=production` to activate the `production`, `proddb` and `prodmq` profiles in one hit.

| Similar to `spring.profiles.active`and`spring.profiles.include`,`spring.profiles.group`can only be used in non-profile-specific documents.
This means it cannot be included in profile specific files or documents activated by`spring.config.activate.on-profile`. |

## Programmatically Setting Profiles

You can programmatically set active profiles by calling `SpringApplication.setAdditionalProfiles(…)` before your application runs.
It is also possible to activate profiles by using Spring’s `ConfigurableEnvironment` interface.

## Profile-specific Configuration Files

Profile-specific variants of both `application.properties` (or `application.yaml`) and files referenced through `@ConfigurationProperties` are considered as files and loaded.
See Profile Specific Files for details.

# Logging

Spring Boot uses Commons Logging for all internal logging but leaves the underlying log implementation open. Default configurations are provided for Java Util Logging, Log4j2, and Logback. In each case, loggers are pre-configured to use console output with optional file output also available.

By default, if you use the starters, Logback is used for logging. Appropriate Logback routing is also included to ensure that dependent libraries that use Java Util Logging, Commons Logging, Log4J, or SLF4J all work correctly.

| There are a lot of logging frameworks available for Java. Do not worry if the above list seems confusing. Generally, you do not need to change your logging dependencies and the Spring Boot defaults work just fine. |

| When you deploy your application to a servlet container or application server, logging performed with the Java Util Logging API is not routed into your application’s logs. This prevents logging performed by the container or other applications that have been deployed to it from appearing in your application’s logs. |

## Log Format

The default log output from Spring Boot resembles the following example:

```
2026-06-10T16:33:37.305Z INFO 53358 --- [myapp] [ main] o.s.b.d.f.logexample.MyApplication : Starting MyApplication using Java 25.0.3 with PID 53358 (/opt/apps/myapp.jar started by myuser in /opt/apps/)
2026-06-10T16:33:37.313Z INFO 53358 --- [myapp] [ main] o.s.b.d.f.logexample.MyApplication : No active profile set, falling back to 1 default profile: "default"
2026-06-10T16:33:39.683Z INFO 53358 --- [myapp] [ main] o.s.boot.tomcat.TomcatWebServer : Tomcat initialized with port 8080 (http)
2026-06-10T16:33:39.703Z INFO 53358 --- [myapp] [ main] o.apache.catalina.core.StandardService : Starting service [Tomcat]
2026-06-10T16:33:39.704Z INFO 53358 --- [myapp] [ main] o.apache.catalina.core.StandardEngine : Starting Servlet engine: [Apache Tomcat/11.0.22]
2026-06-10T16:33:39.776Z INFO 53358 --- [myapp] [ main] b.w.c.s.WebApplicationContextInitializer : Root WebApplicationContext: initialization completed in 2312 ms
2026-06-10T16:33:40.717Z INFO 53358 --- [myapp] [ main] o.s.boot.tomcat.TomcatWebServer : Tomcat started on port 8080 (http) with context path '/'
2026-06-10T16:33:40.760Z INFO 53358 --- [myapp] [ main] o.s.b.d.f.logexample.MyApplication : Started MyApplication in 4.979 seconds (process running for 6.467)
2026-06-10T16:33:40.804Z INFO 53358 --- [myapp] [ionShutdownHook] o.s.boot.tomcat.GracefulShutdown : Commencing graceful shutdown. Waiting for active requests to complete
2026-06-10T16:33:40.833Z INFO 53358 --- [myapp] [tomcat-shutdown] o.s.boot.tomcat.GracefulShutdown : Graceful shutdown complete
```
The following items are output:

-
Date and Time: Millisecond precision and easily sortable.
-
Log Level: `ERROR`,`WARN`,`INFO`,`DEBUG`, or`TRACE`.
-
Process ID.
-
A `---`separator to distinguish the start of actual log messages.
-
Application name: Enclosed in square brackets (logged by default only if `spring.application.name`is set)
-
Application group: Enclosed in square brackets (logged by default only if `spring.application.group`is set)
-
Thread name: Enclosed in square brackets (may be truncated for console output).
-
Correlation ID: If tracing is enabled (not shown in the sample above)
-
Logger name: This is usually the source class name (often abbreviated).
-
The log message.

| Logback does not have a `FATAL`level.
It is mapped to`ERROR`. |

| If you have a `spring.application.name`property but don’t want it logged you can set`logging.include-application-name`to`false`. |

| If you have a `spring.application.group`property but don’t want it logged you can set`logging.include-application-group`to`false`. |

| For more details about correlation IDs, please see this documentation. |

## Console Output

The default log configuration echoes messages to the console as they are written.
By default, `ERROR`-level, `WARN`-level, and `INFO`-level messages are logged.
You can also enable a “debug” mode by starting your application with a `--debug` flag.

`$ java -jar myapp.jar --debug`| You can also specify `debug=true`in your`application.properties`. |

When the debug mode is enabled, a selection of core loggers (embedded container, Hibernate, and Spring Boot) are configured to output more information.
Enabling the debug mode does *not* configure your application to log all messages with `DEBUG` level.

Alternatively, you can enable a “trace” mode by starting your application with a `--trace` flag (or `trace=true` in your `application.properties`).
Doing so enables trace logging for a selection of core loggers (embedded container, Hibernate schema generation, and the whole Spring portfolio).

If you want to disable console-based logging, you can set the `logging.console.enabled` property to `false`.

### Color-coded Output

If your terminal supports ANSI, color output is used to aid readability.
You can set `spring.output.ansi.enabled` to a supported value to override the auto-detection.

Color coding is configured by using the `%clr` conversion word.
In its simplest form, the converter colors the output according to the log level, as shown in the following example:

`%clr(%5p)`The following table describes the mapping of log levels to colors:

| Level | Color |
|---|---|
|
 | Red |
|
 | Red |
|
 | Yellow |
|
 | Green |
|
 | Green |
|
 | Green |

Alternatively, you can specify the color and styles that should be used by providing them as options to the conversion. For example, to make the text yellow and bold, use the following setting:

`%clr(%d{yyyy-MM-dd'T'HH:mm:ss.SSSXXX}){yellow,bold}`The following text colors are supported:

-
`black`
-
`blue`
-
`bright_black`
-
`bright_blue`
-
`bright_cyan`
-
`bright_green`
-
`bright_magenta`
-
`bright_red`
-
`bright_white`
-
`bright_yellow`
-
`cyan`
-
`green`
-
`magenta`
-
`red`
-
`white`
-
`yellow`

The following background colors are supported:

-
`bg_black`
-
`bg_blue`
-
`bg_bright_black`
-
`bg_bright_blue`
-
`bg_bright_cyan`
-
`bg_bright_green`
-
`bg_bright_magenta`
-
`bg_bright_red`
-
`bg_bright_white`
-
`bg_bright_yellow`
-
`bg_cyan`
-
`bg_green`
-
`bg_magenta`
-
`bg_red`
-
`bg_white`
-
`bg_yellow`

The following styles are supported:

-
`bold`
-
`faint`
-
`italic`
-
`normal`
-
`reverse`
-
`underline`

## File Output

By default, Spring Boot logs only to the console and does not write log files.
If you want to write log files in addition to the console output, you need to set a `logging.file.name` or `logging.file.path` property (for example, in your `application.properties`).
If both properties are set, `logging.file.path` is ignored and only `logging.file.name` is used.

The following table shows how the `logging.*` properties can be used together:

| `logging.file.name` | `logging.file.path` | Description |
|---|---|---|
|
 |
 | Console only logging. |
| Specific file (for example, |
 | Writes to the location specified by |
|
 | Specific directory (for example, | Writes |
| Specific file | Specific directory | Writes to the location specified by |

Log files rotate when they reach 10 MB and, as with console output, `ERROR`-level, `WARN`-level, and `INFO`-level messages are logged by default.
Note that Log4J2 requires `logging.file.path` to be set with such configuration.

| Logging properties are independent of the actual logging infrastructure.
As a result, specific configuration keys (such as `logback.configurationFile`for Logback) are not managed by spring Boot. |

## File Rotation

If you are using Logback of Log4J2, it is possible to fine-tune log rotation settings using your `application.properties` or `application.yaml` file.
For all other logging system, you will need to configure rotation settings directly yourself.

The following rotation policy properties are supported for Logback:

| Name | Description |
|---|---|
|
 | The filename pattern used to create log archives. |
|
 | If log archive cleanup should occur when the application starts. |
|
 | The maximum size of log file before it is archived. |
|
 | The maximum amount of size log archives can take before being deleted. |
|
 | The maximum number of archive log files to keep (defaults to 7). |

If you are using Log4j2, the following rotation policy properties are available:

| Name | Description |
|---|---|
|
 | The filename pattern used to create log archives. |
|
 | The maximum size of log file before it is archived (defaults to 10MB). |
|
 | The maximum number of archive log files to keep (defaults to 7). |
|
 | Rolling policy strategy (defaults to 'size'). |
|
 | Cron expression used when the strategy is 'cron' (default to every day at midnight). |
|
 | Time based triggering interval when the strategy is 'time' or 'size-and-time' (default to 1). |
|
 | Whether to align the next rollover time to occur at the top of the interval when the strategy is time based (defaults to |

## Log Levels

All the supported logging systems can have the logger levels set in the Spring `Environment` (for example, in `application.properties`) by using `logging.level.<logger-name>=<level>` where `level` is one of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, or OFF.
The `root` logger can be configured by using `logging.level.root`.

The following example shows potential logging settings in `application.properties`:

-
Properties
-
YAML

```
logging.level.root=warn
logging.level.org.springframework.web=debug
logging.level.org.hibernate=error
```
```
logging:
 level:
 root: "warn"
 org.springframework.web: "debug"
 org.hibernate: "error"
```
It is also possible to set logging levels using environment variables.
For example, `LOGGING_LEVEL_ORG_SPRINGFRAMEWORK_WEB=DEBUG` will set `org.springframework.web` to `DEBUG`.

| The above approach will only work for package level logging.
Since relaxed binding always converts environment variables to lowercase, it is not possible to configure logging for an individual class in this way.
If you need to configure logging for a class, you can use the `SPRING_APPLICATION_JSON`variable. |

## Log Groups

It is often useful to be able to group related loggers together so that they can all be configured at the same time.
For example, you might commonly change the logging levels for *all* Tomcat related loggers, but you can not easily remember top level packages.

To help with this, Spring Boot allows you to define logging groups in your Spring `Environment`.
For example, here is how you could define a “tomcat” group by adding it to your `application.properties`:

-
Properties
-
YAML

`logging.group.tomcat=org.apache.catalina,org.apache.coyote,org.apache.tomcat````
logging:
 group:
 tomcat: "org.apache.catalina,org.apache.coyote,org.apache.tomcat"
```
Once defined, you can change the level for all the loggers in the group with a single line:

-
Properties
-
YAML

`logging.level.tomcat=trace````
logging:
 level:
 tomcat: "trace"
```
Spring Boot includes the following pre-defined logging groups that can be used out-of-the-box:

| Name | Loggers |
|---|---|
| web |
 |
| sql |
 |

## Using a Log Shutdown Hook

In order to release logging resources when your application terminates, a shutdown hook that will trigger log system cleanup when the JVM exits is provided.
This shutdown hook is registered automatically unless your application is deployed as a war file.
If your application has complex context hierarchies the shutdown hook may not meet your needs.
If it does not, disable the shutdown hook and investigate the options provided directly by the underlying logging system.
For example, Logback offers context selectors which allow each Logger to be created in its own context.
You can use the `logging.register-shutdown-hook` property to disable the shutdown hook.
Setting it to `false` will disable the registration.
You can set the property in your `application.properties` or `application.yaml` file:

-
Properties
-
YAML

`logging.register-shutdown-hook=false````
logging:
 register-shutdown-hook: false
```
## Custom Log Configuration

The various logging systems can be activated by including the appropriate libraries on the classpath and can be further customized by providing a suitable configuration file in the root of the classpath or in a location specified by the following Spring `Environment` property: `logging.config`.

You can force Spring Boot to use a particular logging system by using the `org.springframework.boot.logging.LoggingSystem` system property.
The value should be the fully qualified class name of a `LoggingSystem` implementation.
You can also disable Spring Boot’s logging configuration entirely by using a value of `none`.

| Since logging is initialized beforethe`ApplicationContext`is created, it is not possible to control logging from`@PropertySources`in Spring`@Configuration`files.
The only way to change the logging system or disable it entirely is through System properties. |

Depending on your logging system, the following files are loaded:

| Logging System | Customization |
|---|---|
| Logback |
 |
| Log4j2 |
 |
| JDK (Java Util Logging) |
 |

| When possible, we recommend that you use the `-spring`variants for your logging configuration (for example,`logback-spring.xml`rather than`logback.xml`).
If you use standard configuration locations, Spring cannot completely control log initialization. |

| There are known classloading issues with Java Util Logging that cause problems when running from an 'executable jar'. We recommend that you avoid it when running from an 'executable jar' if at all possible. |

To help with the customization, some other properties are transferred from the Spring `Environment` to System properties.
This allows the properties to be consumed by logging system configuration. For example, setting `logging.file.name` in `application.properties` or `LOGGING_FILE_NAME` as an environment variable will result in the `LOG_FILE` System property being set.
The properties that are transferred are described in the following table:

| Spring Environment | System Property | Comments |
|---|---|---|
|
 |
 | The conversion word used when logging exceptions. |
|
 |
 | If defined, it is used in the default log configuration. |
|
 |
 | If defined, it is used in the default log configuration. |
|
 |
 | The log pattern to use on the console (stdout). |
|
 |
 | Appender pattern for log date format. |
|
 |
 | The charset to use for console logging. |
|
 |
 | The log level threshold to use for console logging. |
|
 |
 | The log pattern to use in a file (if |
|
 |
 | The charset to use for file logging (if |
|
 |
 | The log level threshold to use for file logging. |
|
 |
 | The format to use when rendering the log level (default |
|
 |
 | The structured logging format to use for console logging. |
|
 |
 | The structured logging format to use for file logging. |
|
 |
 | The current process ID (discovered if possible and when not already defined as an OS environment variable). |

If you use Logback, the following properties are also transferred:

| Spring Environment | System Property | Comments |
|---|---|---|
|
 |
 | Pattern for rolled-over log file names (default |
|
 |
 | Whether to clean the archive log files on startup. |
|
 |
 | Maximum log file size. |
|
 |
 | Total size of log backups to be kept. |
|
 |
 | Maximum number of archive log files to keep. |

All the supported logging systems can consult System properties when parsing their configuration files.
See the default configurations in `spring-boot.jar` for examples:

| If you want to use a placeholder in a logging property, you should use Spring Boot’s syntax and not the syntax of the underlying framework.
Notably, if you use Logback, you should use |

| You can add MDC and other ad-hoc content to log lines by overriding only the |

## Structured Logging

Structured logging is a technique where the log output is written in a well-defined, often machine-readable format. Spring Boot supports structured logging and has support for the following JSON formats out of the box:

To enable structured logging, set the property `logging.structured.format.console` (for console output) or `logging.structured.format.file` (for file output) to the id of the format you want to use.

If you are using Custom Log Configuration, update your configuration to respect `CONSOLE_LOG_STRUCTURED_FORMAT` and `FILE_LOG_STRUCTURED_FORMAT` system properties.
Take `CONSOLE_LOG_STRUCTURED_FORMAT` for example:

-
Logback
-
Log4j2

```
<!-- replace your encoder with StructuredLogEncoder -->
<encoder class="org.springframework.boot.logging.logback.StructuredLogEncoder">
	<format>${CONSOLE_LOG_STRUCTURED_FORMAT}</format>
	<charset>${CONSOLE_LOG_CHARSET}</charset>
</encoder>
```
You can also refer to the default configurations included in Spring Boot:

```
<!-- replace your PatternLayout with StructuredLogLayout -->
<StructuredLogLayout format="${sys:CONSOLE_LOG_STRUCTURED_FORMAT}" charset="${sys:CONSOLE_LOG_CHARSET}"/>
```
You can also refer to the default configurations included in Spring Boot:

### Elastic Common Schema

Elastic Common Schema is a JSON based logging format.

To enable the Elastic Common Schema log format, set the appropriate `format` property to `ecs`:

-
Properties
-
YAML

```
logging.structured.format.console=ecs
logging.structured.format.file=ecs
```
```
logging:
 structured:
 format:
 console: ecs
 file: ecs
```
A log line looks like this:

`{"@timestamp":"2024-01-01T10:15:00.067462556Z","log":{"level":"INFO","logger":"org.example.Application"},"process":{"pid":39599,"thread":{"name":"main"}},"service":{"name":"simple"},"message":"No active profile set, falling back to 1 default profile: \"default\"","ecs":{"version":"8.11"}}`This format also adds every key value pair contained in the MDC to the JSON object. You can also use the SLF4J fluent logging API to add key value pairs to the logged JSON object with the addKeyValue method.

The `service` values can be customized using `logging.structured.ecs.service` properties:

-
Properties
-
YAML

```
logging.structured.ecs.service.name=MyService
logging.structured.ecs.service.version=1
logging.structured.ecs.service.environment=Production
logging.structured.ecs.service.node-name=Primary
```
```
logging:
 structured:
 ecs:
 service:
 name: MyService
 version: 1.0
 environment: Production
 node-name: Primary
```
| `logging.structured.ecs.service.name`will default to`spring.application.name`if not specified. |

| `logging.structured.ecs.service.version`will default to`spring.application.version`if not specified. |

### Graylog Extended Log Format (GELF)

Graylog Extended Log Format is a JSON based logging format for the Graylog log analytics platform.

To enable the Graylog Extended Log Format, set the appropriate `format` property to `gelf`:

-
Properties
-
YAML

```
logging.structured.format.console=gelf
logging.structured.format.file=gelf
```
```
logging:
 structured:
 format:
 console: gelf
 file: gelf
```
A log line looks like this:

`{"version":"1.1","short_message":"No active profile set, falling back to 1 default profile: \"default\"","timestamp":1725958035.857,"level":6,"_level_name":"INFO","_process_pid":47649,"_process_thread_name":"main","_log_logger":"org.example.Application"}`This format also adds every key value pair contained in the MDC to the JSON object. You can also use the SLF4J fluent logging API to add key value pairs to the logged JSON object with the addKeyValue method.

Several fields can be customized using `logging.structured.gelf` properties:

-
Properties
-
YAML

```
logging.structured.gelf.host=MyService
logging.structured.gelf.service.version=1
```
```
logging:
 structured:
 gelf:
 host: MyService
 service:
 version: 1.0
```
| `logging.structured.gelf.host`will default to`spring.application.name`if not specified. |

| `logging.structured.gelf.service.version`will default to`spring.application.version`if not specified. |

### Logstash JSON format

The Logstash JSON format is a JSON based logging format.

To enable the Logstash JSON log format, set the appropriate `format` property to `logstash`:

-
Properties
-
YAML

```
logging.structured.format.console=logstash
logging.structured.format.file=logstash
```
```
logging:
 structured:
 format:
 console: logstash
 file: logstash
```
A log line looks like this:

`{"@timestamp":"2024-01-01T10:15:00.111037681+02:00","@version":"1","message":"No active profile set, falling back to 1 default profile: \"default\"","logger_name":"org.example.Application","thread_name":"main","level":"INFO","level_value":20000}`This format also adds every key value pair contained in the MDC to the JSON object. You can also use the SLF4J fluent logging API to add key value pairs to the logged JSON object with the addKeyValue method.

If you add markers, these will show up in a `tags` string array in the JSON.

### Customizing Structured Logging JSON

Spring Boot attempts to pick sensible defaults for the JSON names and values output for structured logging. Sometimes, however, you may want to make small adjustments to the JSON for your own needs. For example, it’s possible that you might want to change some of the names to match the expectations of your log ingestion system. You might also want to filter out certain members since you don’t find them useful.

The following properties allow you to change the way that structured logging JSON is written:

| Property | Description |
|---|---|
|
 | Filters specific paths from the JSON |
|
 | Renames a specific member in the JSON |
|
 | Adds additional members to the JSON |

For example, the following will exclude `log.level`, rename `process.id` to `procid` and add a fixed `corpname` field:

-
Properties
-
YAML

```
logging.structured.json.exclude=log.level
logging.structured.json.rename.process.id=procid
logging.structured.json.add.corpname=mycorp
```
```
logging:
 structured:
 json:
 exclude: log.level
 rename:
 process.id: procid
 add:
 corpname: mycorp
```
| For more advanced customizations, you can use the `StructuredLoggingJsonMembersCustomizer`interface.
You can reference one or more implementations using the`logging.structured.json.customizer`property.
You can also declare implementations by listing them in a`META-INF/spring.factories`file. |

### Customizing Structured Logging Stack Traces

Complete stack traces are included in the JSON output whenever a message is logged with an exception. This amount of information may be costly to process by your log ingestion system, so you may want to tune the way that stack traces are printed.

To do this, you can use one or more of the following properties:

| Property | Description |
|---|---|
|
 | Use |
|
 | The maximum length that should be printed |
|
 | The maximum number of frames to print per stack trace (including common and suppressed frames) |
|
 | If common frames should be included or removed |
|
 | If a hash of the stack trace should be included |

For example, the following will use root first stack traces, limit their length, and include hashes.

-
Properties
-
YAML

```
logging.structured.json.stacktrace.root=first
logging.structured.json.stacktrace.max-length=1024
logging.structured.json.stacktrace.include-common-frames=true
logging.structured.json.stacktrace.include-hashes=true
```
```
logging:
 structured:
 json:
 stacktrace:
 root: first
 max-length: 1024
 include-common-frames: true
 include-hashes: true
```
| If you need complete control over stack trace printing you can set Your |

### Supporting Other Structured Logging Formats

The structured logging support in Spring Boot is extensible, allowing you to define your own custom format.
To do this, implement the `StructuredLogFormatter` interface. The generic type argument has to be `ILoggingEvent` when using Logback and `LogEvent` when using Log4j2 (that means your implementation is tied to a specific logging system).
Your implementation is then called with the log event and returns the `String` to be logged, as seen in this example:

-
Java
-
Kotlin

```
import ch.qos.logback.classic.spi.ILoggingEvent;
import org.springframework.boot.logging.structured.StructuredLogFormatter;
class MyCustomFormat implements StructuredLogFormatter<ILoggingEvent> {
	@Override
	public String format(ILoggingEvent event) {
 return "time=" + event.getInstant() + " level=" + event.getLevel() + " message=" + event.getMessage() + "\n";
	}
}
```
```
import ch.qos.logback.classic.spi.ILoggingEvent
import org.springframework.boot.logging.structured.StructuredLogFormatter
class MyCustomFormat : StructuredLogFormatter<ILoggingEvent> {
	override fun format(event: ILoggingEvent): String {
 return "time=${event.instant} level=${event.level} message=${event.message}\n"
	}
}
```
As you can see in the example, you can return any format, it doesn’t have to be JSON.

To enable your custom format, set the property `logging.structured.format.console` or `logging.structured.format.file` to the fully qualified class name of your implementation.

Your implementation can use some constructor parameters, which are injected automatically.
Please see the JavaDoc of `StructuredLogFormatter` for more details.

## Logback Extensions

Spring Boot includes a number of extensions to Logback that can help with advanced configuration.
You can use these extensions in your `logback-spring.xml` configuration file.

| Because the standard `logback.xml`configuration file is loaded too early, you cannot use extensions in it.
You need to either use`logback-spring.xml`or define a`logging.config`property. |

| The extensions cannot be used with Logback’s configuration scanning. If you attempt to do so, making changes to the configuration file results in an error similar to one of the following being logged: |

```
ERROR in ch.qos.logback.core.joran.spi.Interpreter@4:71 - no applicable action for [springProperty], current ElementPath is [[configuration][springProperty]]
ERROR in ch.qos.logback.core.joran.spi.Interpreter@4:71 - no applicable action for [springProfile], current ElementPath is [[configuration][springProfile]]
```
### Profile-specific Configuration

The `<springProfile>` tag lets you optionally include or exclude sections of configuration based on the active Spring profiles.
Profile sections are supported anywhere within the `<configuration>` element.
Use the `name` attribute to specify which profile accepts the configuration.
The `<springProfile>` tag can contain a profile name (for example `staging`) or a profile expression.
A profile expression allows for more complicated profile logic to be expressed, for example `production & (eu-central | eu-west)`.
Check the Spring Framework reference guide for more details.
The following listing shows three sample profiles:

```
<springProfile name="staging">
	<!-- configuration to be enabled when the "staging" profile is active -->
</springProfile>
<springProfile name="dev | staging">
	<!-- configuration to be enabled when the "dev" or "staging" profiles are active -->
</springProfile>
<springProfile name="!production">
	<!-- configuration to be enabled when the "production" profile is not active -->
</springProfile>
```
### Environment Properties

The `<springProperty>` tag lets you expose properties from the Spring `Environment` for use within Logback.
Doing so can be useful if you want to access values from your `application.properties` file in your Logback configuration.
The tag works in a similar way to Logback’s standard `<property>` tag.
However, rather than specifying a direct `value`, you specify the `source` of the property (from the `Environment`).
If you need to store the property somewhere other than in `local` scope, you can use the `scope` attribute.
If you need a fallback value (in case the property is not set in the `Environment`), you can use the `defaultValue` attribute.
The following example shows how to expose properties for use within Logback:

```
<springProperty scope="context" name="fluentHost" source="myapp.fluentd.host"
 defaultValue="localhost"/>
<appender name="FLUENT" class="ch.qos.logback.more.appenders.DataFluentAppender">
	<remoteHost>${fluentHost}</remoteHost>
	...
</appender>
```
| The `source`must be specified in kebab case (such as`my.property-name`).
However, properties can be added to the`Environment`by using the relaxed rules. |

## Log4j2 Extensions

Spring Boot includes a number of extensions to Log4j2 that can help with advanced configuration.
You can use these extensions in any `log4j2-spring.xml` configuration file.

| Because the standard `log4j2.xml`configuration file is loaded too early, you cannot use extensions in it.
You need to either use`log4j2-spring.xml`or define a`logging.config`property. |

| The extensions supersede the Spring Boot support provided by Log4J.
You should make sure not to include the `org.apache.logging.log4j:log4j-spring-boot`module in your build. |

### Profile-specific Configuration

The `<SpringProfile>` tag lets you optionally include or exclude sections of configuration based on the active Spring profiles.
Profile sections are supported anywhere within the `<Configuration>` element.
Use the `name` attribute to specify which profile accepts the configuration.
The `<SpringProfile>` tag can contain a profile name (for example `staging`) or a profile expression.
A profile expression allows for more complicated profile logic to be expressed, for example `production & (eu-central | eu-west)`.
Check the Spring Framework reference guide for more details.
The following listing shows three sample profiles:

```
<SpringProfile name="staging">
	<!-- configuration to be enabled when the "staging" profile is active -->
</SpringProfile>
<SpringProfile name="dev | staging">
	<!-- configuration to be enabled when the "dev" or "staging" profiles are active -->
</SpringProfile>
<SpringProfile name="!production">
	<!-- configuration to be enabled when the "production" profile is not active -->
</SpringProfile>
```
### Environment Properties Lookup

If you want to refer to properties from your Spring `Environment` within your Log4j2 configuration you can use `spring:` prefixed lookups.
Doing so can be useful if you want to access values from your `application.properties` file in your Log4j2 configuration.

The following example shows how to set Log4j2 properties named `applicationName` and `applicationGroup` that read `spring.application.name` and `spring.application.group` from the Spring `Environment`:

```
<Properties>
	<Property name="applicationName">${spring:spring.application.name}</Property>
	<Property name="applicationGroup">${spring:spring.application.group}</Property>
</Properties>
```
| The lookup key should be specified in kebab case (such as `my.property-name`). |

### Log4j2 System Properties

Log4j2 supports a number of System Properties that can be used to configure various items.
For example, the `log4j2.skipJansi` system property can be used to configure if the `ConsoleAppender` will try to use a Jansi output stream on Windows.

All system properties that are loaded after the Log4j2 initialization can be obtained from the Spring `Environment`.
For example, you could add `log4j2.skipJansi=false` to your `application.properties` file to have the `ConsoleAppender` use Jansi on Windows.

| The Spring `Environment`is only considered when system properties and OS environment variables do not contain the value being loaded. |

| System properties that are loaded during early Log4j2 initialization cannot reference the Spring `Environment`.
For example, the property Log4j2 uses to allow the default Log4j2 implementation to be chosen is used before the Spring Environment is available. |

# Internationalization

Spring Boot supports localized messages so that your application can cater to users of different language preferences.
By default, Spring Boot looks for the presence of a `messages` resource bundle at the root of the classpath.

| The auto-configuration applies when the default properties file for the configured resource bundle is available ( `messages.properties`by default).
If your resource bundle contains only language-specific properties files, you are required to add the default.
If no properties file is found that matches any of the configured base names, there will be no auto-configured`MessageSource`. |

The basename of the resource bundle as well as several other attributes can be configured using the `spring.messages` namespace, as shown in the following example:

-
Properties
-
YAML

```
spring.messages.basename=messages, config.i18n.messages
spring.messages.common-messages=classpath:my-common-messages.properties
spring.messages.fallback-to-system-locale=false
```
```
spring:
 messages:
 basename: "messages, config.i18n.messages"
 common-messages: "classpath:my-common-messages.properties"
 fallback-to-system-locale: false
```
| The `spring.messages.basename`property supports a list of locations, either a package qualifier or a resource resolved from the classpath root.
The`spring.messages.common-messages`property supports a list of property file resources. |

See `MessageSourceProperties` for more supported options.

# Aspect-Oriented Programming

Spring Boot provides auto-configuration for aspect-oriented programming (AOP). You can learn more about AOP with Spring in the Spring Framework reference documentation.

By default, Spring Boot’s auto-configuration configures Spring AOP to use CGLib proxies.
To use JDK proxies instead, set `spring.aop.proxy-target-class` to `false`.

If AspectJ is on the classpath, Spring Boot’s auto-configuration will automatically enable AspectJ auto proxy such that `@EnableAspectJAutoProxy` is not required.

# JSON

Spring Boot provides integration with the following JSON mapping libraries:

-
Jackson 3
-
Jackson 2
-
Gson
-
JSON-B
-
Kotlin Serialization

Jackson 3 is the preferred and default library.

Support for Jackson 2 is deprecated and will be removed in a future Spring Boot 4.x release. It is provided purely to ease the migration from Jackson 2 to Jackson 3 and should not be relied up in the longer term.

## Jackson 3

Auto-configuration for Jackson 3 is provided and Jackson is part of `spring-boot-starter-json`.
When Jackson is on the classpath a `JsonMapper` bean is automatically configured.
Several configuration properties are provided for customizing the configuration of the `JsonMapper`.

### Custom Serializers and Deserializers

If you use Jackson to serialize and deserialize JSON data, you might want to write your own `ValueSerializer` and `ValueDeserializer` classes.
Custom serializers are usually registered with Jackson through a module, but Spring Boot provides an alternative `@JacksonComponent` annotation that makes it easier to directly register Spring Beans.

You can use the `@JacksonComponent` annotation directly on `ValueSerializer`, `ValueDeserializer` or `KeyDeserializer` implementations.
You can also use it on classes that contain serializers/deserializers as inner classes, as shown in the following example:

-
Java
-
Kotlin

```
import tools.jackson.core.JsonGenerator;
import tools.jackson.core.JsonParser;
import tools.jackson.databind.DeserializationContext;
import tools.jackson.databind.JsonNode;
import tools.jackson.databind.SerializationContext;
import tools.jackson.databind.ValueDeserializer;
import tools.jackson.databind.ValueSerializer;
import org.springframework.boot.jackson.JacksonComponent;
@JacksonComponent
public class MyJacksonComponent {
	public static class Serializer extends ValueSerializer<MyObject> {
 @Override
 public void serialize(MyObject value, JsonGenerator jgen, SerializationContext context) {
 jgen.writeStartObject();
 jgen.writeStringProperty("name", value.getName());
 jgen.writeNumberProperty("age", value.getAge());
 jgen.writeEndObject();
 }
	}
	public static class Deserializer extends ValueDeserializer<MyObject> {
 @Override
 public MyObject deserialize(JsonParser jsonParser, DeserializationContext ctxt) {
 JsonNode tree = jsonParser.readValueAsTree();
 String name = tree.get("name").stringValue();
 int age = tree.get("age").intValue();
 return new MyObject(name, age);
 }
	}
}
```
```
import tools.jackson.core.JsonGenerator
import tools.jackson.core.JsonParser
import tools.jackson.databind.DeserializationContext
import tools.jackson.databind.JsonNode
import tools.jackson.databind.SerializationContext
import tools.jackson.databind.ValueDeserializer
import tools.jackson.databind.ValueSerializer
import org.springframework.boot.jackson.JacksonComponent
@JacksonComponent
class MyJacksonComponent {
	class Serializer : ValueSerializer<MyObject>() {
 override fun serialize(value: MyObject, jgen: JsonGenerator, serializers: SerializationContext) {
 jgen.writeStartObject()
 jgen.writeStringProperty("name", value.name)
 jgen.writeNumberProperty("age", value.age)
 jgen.writeEndObject()
 }
	}
	class Deserializer : ValueDeserializer<MyObject>() {
 override fun deserialize(jsonParser: JsonParser, ctxt: DeserializationContext): MyObject {
 val tree = jsonParser.readValueAsTree<JsonNode>()
 val name = tree["name"].stringValue()
 val age = tree["age"].intValue()
 return MyObject(name, age)
 }
	}
}
```
All `@JacksonComponent` beans in the `ApplicationContext` are automatically registered with Jackson.
Because `@JacksonComponent` is meta-annotated with `@Component`, the usual component-scanning rules apply.

Spring Boot also provides `ObjectValueSerializer` and `ObjectValueDeserializer` base classes that provide useful alternatives to the standard Jackson versions when serializing objects.
See `ObjectValueSerializer` and `ObjectValueDeserializer` in the API documentation for details.

The example above can be rewritten to use `ObjectValueSerializer` and `ObjectValueDeserializer` as follows:

-
Java
-
Kotlin

```
import tools.jackson.core.JsonGenerator;
import tools.jackson.core.JsonParser;
import tools.jackson.databind.DeserializationContext;
import tools.jackson.databind.JsonNode;
import tools.jackson.databind.SerializationContext;
import org.springframework.boot.jackson.JacksonComponent;
import org.springframework.boot.jackson.ObjectValueDeserializer;
import org.springframework.boot.jackson.ObjectValueSerializer;
@JacksonComponent
public class MyJacksonComponent {
	public static class Serializer extends ObjectValueSerializer<MyObject> {
 @Override
 protected void serializeObject(MyObject value, JsonGenerator jgen, SerializationContext context) {
 jgen.writeStringProperty("name", value.getName());
 jgen.writeNumberProperty("age", value.getAge());
 }
	}
	public static class Deserializer extends ObjectValueDeserializer<MyObject> {
 @Override
 protected MyObject deserializeObject(JsonParser jsonParser, DeserializationContext context, JsonNode tree) {
 String name = nullSafeValue(tree.get("name"), String.class);
 int age = nullSafeValue(tree.get("age"), Integer.class);
 return new MyObject(name, age);
 }
	}
}
```
```
import tools.jackson.core.JsonGenerator
import tools.jackson.core.JsonParser
import tools.jackson.databind.DeserializationContext
import tools.jackson.databind.JsonNode
import tools.jackson.databind.SerializationContext
import org.springframework.boot.jackson.JacksonComponent;
import org.springframework.boot.jackson.ObjectValueDeserializer
import org.springframework.boot.jackson.ObjectValueSerializer
@JacksonComponent
class MyJacksonComponent {
	class Serializer : ObjectValueSerializer<MyObject>() {
 override fun serializeObject(value: MyObject, jgen: JsonGenerator, context: SerializationContext) {
 jgen.writeStringProperty("name", value.name)
 jgen.writeNumberProperty("age", value.age)
 }
	}
	class Deserializer : ObjectValueDeserializer<MyObject>() {
 override fun deserializeObject(jsonParser: JsonParser, context: DeserializationContext,
 tree: JsonNode): MyObject {
 val name = nullSafeValue(tree["name"], String::class.java) ?: throw IllegalStateException("name is null")
 val age = nullSafeValue(tree["age"], Int::class.java) ?: throw IllegalStateException("age is null")
 return MyObject(name, age)
 }
	}
}
```
### Mixins

Jackson has support for mixins that can be used to mix additional annotations into those already declared on a target class.
Spring Boot’s Jackson auto-configuration will scan your application’s packages for classes annotated with `@JacksonMixin` and register them with the auto-configured `JsonMapper`.
The registration is performed by Spring Boot’s `JacksonMixinModule`.

## Jackson 2

Deprecated auto-configuration for Jackson 2 is provided by the `spring-boot-jackson2` module.
When this module is on the classpath a `ObjectMapper` bean is automatically configured.
Several `spring.jackson2.*` configuration properties are provided for customizing the configuration.
To take more control, define one or more `Jackson2ObjectMapperBuilderCustomizer` beans.

When both Jackson 3 and Jackson 2 are present, various configuration properties can be used to indicate that Jackson 2 is preferred:

-
`spring.graphql.rsocket.preferred-json-mapper`
-
`spring.http.codecs.preferred-json-mapper`(used by Spring WebFlux and reactive HTTP clients)
-
`spring.http.converters.preferred-json-mapper`(used by Spring MVC and imperative HTTP clients)
-
`spring.rsocket.preferred-mapper`
-
`spring.websocket.messaging.preferred-json-mapper`

In each case, set the relevant property to `jackson2` to indicate that Jackson 2 is preferred.

## Gson

Auto-configuration for Gson is provided.
When Gson is on the classpath a `Gson` bean is automatically configured.
Several `spring.gson.*` configuration properties are provided for customizing the configuration.
To take more control, one or more `GsonBuilderCustomizer` beans can be used.

## JSON-B

Auto-configuration for JSON-B is provided.
When the JSON-B API and an implementation are on the classpath a `Jsonb` bean will be automatically configured.
The preferred JSON-B implementation is Eclipse Yasson for which dependency management is provided.

## Kotlin Serialization

Auto-configuration for Kotlin Serialization is provided.
When `kotlinx-serialization-json` is on the classpath a Json bean is automatically configured.
Several `spring.kotlinx.serialization.json.*` configuration properties are provided for customizing the configuration.

# Task Execution and Scheduling

In the absence of an `Executor` bean in the context, Spring Boot auto-configures an `AsyncTaskExecutor`.
When virtual threads are enabled (using Java 21+ and `spring.threads.virtual.enabled` set to `true`) this will be a `SimpleAsyncTaskExecutor` that uses virtual threads.
Otherwise, it will be a `ThreadPoolTaskExecutor` with sensible defaults.

The auto-configured `AsyncTaskExecutor` is used for the following integrations unless a custom `Executor` bean is defined:

-
Execution of asynchronous tasks using `@EnableAsync`, unless a bean of type`AsyncConfigurer`is defined.
-
Asynchronous handling of `Callable`return values from controller methods in Spring for GraphQL.
-
Asynchronous request handling in Spring MVC.
-
Support for blocking execution in Spring WebFlux.
-
Utilized for inbound and outbound message channels in Spring WebSocket.
-
Bootstrap executor for JPA, based on the bootstrap mode of JPA repositories.
-
Bootstrap executor for background initialization of beans in the `ApplicationContext`.

While this approach works in most scenarios, Spring Boot allows you to override the auto-configured `AsyncTaskExecutor`.
By default, when a custom `Executor` bean is registered, the auto-configured `AsyncTaskExecutor` backs off, and the custom `Executor` is used for regular task execution (via `@EnableAsync`).

However, Spring MVC, Spring WebFlux, and Spring GraphQL all require a bean named `applicationTaskExecutor`.
For Spring MVC and Spring WebFlux, this bean must be of type `AsyncTaskExecutor`, whereas Spring GraphQL does not enforce this type requirement.

Spring WebSocket and JPA will use `AsyncTaskExecutor` if either a single bean of this type is available or a bean named `applicationTaskExecutor` is defined.

Finally, the boostrap executor of the `ApplicationContext` uses a bean named `applicationTaskExecutor` unless a bean named `bootstrapExecutor` is defined.

The following code snippet demonstrates how to register a custom `AsyncTaskExecutor` to be used with Spring MVC, Spring WebFlux, Spring GraphQL, Spring WebSocket, JPA, and background initialization of beans.

-
Java
-
Kotlin

```
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.task.SimpleAsyncTaskExecutor;
@Configuration(proxyBeanMethods = false)
public class MyTaskExecutorConfiguration {
	@Bean("applicationTaskExecutor")
	SimpleAsyncTaskExecutor applicationTaskExecutor() {
 return new SimpleAsyncTaskExecutor("app-");
	}
}
```
```
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.core.task.SimpleAsyncTaskExecutor
@Configuration(proxyBeanMethods = false)
class MyTaskExecutorConfiguration {
	@Bean("applicationTaskExecutor")
	fun applicationTaskExecutor(): SimpleAsyncTaskExecutor {
 return SimpleAsyncTaskExecutor("app-")
	}
}
```
| The |

| If neither the auto-configured |

If your application needs multiple `Executor` beans for different integrations, such as one for regular task execution with `@EnableAsync` and other for Spring MVC, Spring WebFlux, Spring WebSocket and JPA, you can configure them as follows.

-
Java
-
Kotlin

```
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.task.SimpleAsyncTaskExecutor;
import org.springframework.scheduling.concurrent.ThreadPoolTaskExecutor;
@Configuration(proxyBeanMethods = false)
public class MyTaskExecutorConfiguration {
	@Bean("applicationTaskExecutor")
	SimpleAsyncTaskExecutor applicationTaskExecutor() {
 return new SimpleAsyncTaskExecutor("app-");
	}
	@Bean("taskExecutor")
	ThreadPoolTaskExecutor taskExecutor() {
 ThreadPoolTaskExecutor threadPoolTaskExecutor = new ThreadPoolTaskExecutor();
 threadPoolTaskExecutor.setThreadNamePrefix("async-");
 return threadPoolTaskExecutor;
	}
}
```
```
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.core.task.SimpleAsyncTaskExecutor
import org.springframework.scheduling.concurrent.ThreadPoolTaskExecutor
@Configuration(proxyBeanMethods = false)
class MyTaskExecutorConfiguration {
	@Bean("applicationTaskExecutor")
	fun applicationTaskExecutor(): SimpleAsyncTaskExecutor {
 return SimpleAsyncTaskExecutor("app-")
	}
	@Bean("taskExecutor")
	fun taskExecutor(): ThreadPoolTaskExecutor {
 val threadPoolTaskExecutor = ThreadPoolTaskExecutor()
 threadPoolTaskExecutor.setThreadNamePrefix("async-")
 return threadPoolTaskExecutor
	}
}
```
| The auto-configured
 |

If a `taskExecutor` named bean is not an option, you can mark your bean as `@Primary` or define an `AsyncConfigurer` bean to specify the `Executor` responsible for handling regular task execution with `@EnableAsync`.
The following example demonstrates how to achieve this.

-
Java
-
Kotlin

```
import java.util.concurrent.Executor;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.scheduling.annotation.AsyncConfigurer;
@Configuration(proxyBeanMethods = false)
public class MyTaskExecutorConfiguration {
	@Bean
	AsyncConfigurer asyncConfigurer(ExecutorService executorService) {
 return new AsyncConfigurer() {
 @Override
 public Executor getAsyncExecutor() {
 return executorService;
 }
 };
	}
	@Bean
	ExecutorService executorService() {
 return Executors.newCachedThreadPool();
	}
}
```
```
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.scheduling.annotation.AsyncConfigurer
import java.util.concurrent.Executor
import java.util.concurrent.ExecutorService
import java.util.concurrent.Executors
@Configuration(proxyBeanMethods = false)
class MyTaskExecutorConfiguration {
	@Bean
	fun asyncConfigurer(executorService: ExecutorService): AsyncConfigurer {
 return object : AsyncConfigurer {
 override fun getAsyncExecutor(): Executor {
 return executorService
 }
 }
	}
	@Bean
	fun executorService(): ExecutorService {
 return Executors.newCachedThreadPool()
	}
}
```
To register a custom `Executor` while keeping the auto-configured `AsyncTaskExecutor`, you can create a custom `Executor` bean and set the `defaultCandidate=false` attribute in its `@Bean` annotation, as demonstrated in the following example:

-
Java
-
Kotlin

```
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyTaskExecutorConfiguration {
	@Bean(defaultCandidate = false)
	@Qualifier("scheduledExecutorService")
	ScheduledExecutorService scheduledExecutorService() {
 return Executors.newSingleThreadScheduledExecutor();
	}
}
```
```
import org.springframework.beans.factory.annotation.Qualifier
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import java.util.concurrent.Executors
import java.util.concurrent.ScheduledExecutorService
@Configuration(proxyBeanMethods = false)
class MyTaskExecutorConfiguration {
	@Bean(defaultCandidate = false)
	@Qualifier("scheduledExecutorService")
	fun scheduledExecutorService(): ScheduledExecutorService {
 return Executors.newSingleThreadScheduledExecutor()
	}
}
```
In that case, you will be able to autowire your custom `Executor` into other components while retaining the auto-configured `AsyncTaskExecutor`.
However, remember to use the `@Qualifier` annotation alongside `@Autowired`.

If this is not possible for you, you can request Spring Boot to auto-configure an `AsyncTaskExecutor` anyway, as follows:

-
Properties
-
YAML

`spring.task.execution.mode=force````
spring:
 task:
 execution:
 mode: force
```
The auto-configured `AsyncTaskExecutor` will be used automatically for all integrations, even if a custom `Executor` bean is registered, including those marked as `@Primary`.
These integrations include:

-
Asynchronous task execution ( `@EnableAsync`), unless an`AsyncConfigurer`bean is present.
-
Spring for GraphQL’s asynchronous handling of `Callable`return values from controller methods.
-
Spring MVC’s asynchronous request processing.
-
Spring WebFlux’s blocking execution support.
-
Utilized for inbound and outbound message channels in Spring WebSocket.
-
Bootstrap executor for JPA, based on the bootstrap mode of JPA repositories.
-
Bootstrap executor for background initialization of beans in the `ApplicationContext`, unless a bean named`bootstrapExecutor`is defined.

| Depending on your target arrangement, you could set |

| When |

When a `ThreadPoolTaskExecutor` is auto-configured, the thread pool uses 8 core threads that can grow and shrink according to the load.
Those default settings can be fine-tuned using the `spring.task.execution` namespace, as shown in the following example:

-
Properties
-
YAML

```
spring.task.execution.pool.max-size=16
spring.task.execution.pool.queue-capacity=100
spring.task.execution.pool.keep-alive=10s
```
```
spring:
 task:
 execution:
 pool:
 max-size: 16
 queue-capacity: 100
 keep-alive: "10s"
```
This changes the thread pool to use a bounded queue so that when the queue is full (100 tasks), the thread pool increases to maximum 16 threads. Shrinking of the pool is more aggressive as threads are reclaimed when they are idle for 10 seconds (rather than 60 seconds by default).

A scheduler can also be auto-configured if it needs to be associated with scheduled task execution (using `@EnableScheduling` for instance).

If virtual threads are enabled (using Java 21+ and `spring.threads.virtual.enabled` set to `true`) this will be a `SimpleAsyncTaskScheduler` that uses virtual threads.
This `SimpleAsyncTaskScheduler` will ignore any pooling related properties.

If virtual threads are not enabled, it will be a `ThreadPoolTaskScheduler` with sensible defaults.
The `ThreadPoolTaskScheduler` uses one thread by default and its settings can be fine-tuned using the `spring.task.scheduling` namespace, as shown in the following example:

-
Properties
-
YAML

```
spring.task.scheduling.thread-name-prefix=scheduling-
spring.task.scheduling.pool.size=2
```
```
spring:
 task:
 scheduling:
 thread-name-prefix: "scheduling-"
 pool:
 size: 2
```
A `ThreadPoolTaskExecutorBuilder` bean, a `SimpleAsyncTaskExecutorBuilder` bean, a `ThreadPoolTaskSchedulerBuilder` bean and a `SimpleAsyncTaskSchedulerBuilder` are made available in the context if a custom executor or scheduler needs to be created.
The `SimpleAsyncTaskExecutorBuilder` and `SimpleAsyncTaskSchedulerBuilder` beans are auto-configured to use virtual threads if they are enabled (using Java 21+ and `spring.threads.virtual.enabled` set to `true`).

# Development-time Services

Development-time services provide external dependencies needed to run the application while developing it. They are only supposed to be used while developing and are disabled when the application is deployed.

Spring Boot offers support for two development time services, Docker Compose and Testcontainers. The next sections will provide more details about them.

## Docker Compose Support

Docker Compose is a popular technology that can be used to define and manage multiple containers for services that your application needs.
A `compose.yml` file is typically created next to your application which defines and configures service containers.

A typical workflow with Docker Compose is to run `docker compose up`, work on your application with it connecting to started services, then run `docker compose down` when you are finished.

The `spring-boot-docker-compose` module can be included in a project to provide support for working with containers using Docker Compose.
Add the module dependency to your build, as shown in the following listings for Maven and Gradle:

```
<dependencies>
	<dependency>
 <groupId>org.springframework.boot</groupId>
 <artifactId>spring-boot-docker-compose</artifactId>
 <optional>true</optional>
	</dependency>
</dependencies>
```
```
dependencies {
	developmentOnly("org.springframework.boot:spring-boot-docker-compose")
}
```
When this module is included as a dependency Spring Boot will do the following:

-
Search for a `compose.yml`and other common compose filenames in your working directory
-
Call `docker compose up`with the discovered`compose.yml`
-
Create service connection beans for each supported container
-
Call `docker compose stop`when the application is shutdown

If the Docker Compose services are already running when starting the application, Spring Boot will only create the service connection beans for each supported container.
It will not call `docker compose up` again and it will not call `docker compose stop` when the application is shutdown.

| Repackaged archives do not contain Spring Boot’s Docker Compose by default.
If you want to use this support, you need to include it.
When using the Maven plugin, set the `excludeDockerCompose`property to`false`.
When using the Gradle plugin, configure the task’s classpath to include the`developmentOnly`configuration. |

### Prerequisites

You need to have the `docker` and `docker compose` (or `docker-compose`) CLI applications on your path.
The minimum supported Docker Compose version is 2.2.0.

### Service Connections

A service connection is a connection to any remote service. Spring Boot’s auto-configuration can consume the details of a service connection and use them to establish a connection to a remote service. When doing so, the connection details take precedence over any connection-related configuration properties.

When using Spring Boot’s Docker Compose support, service connections are established to the port mapped by the container.

| Docker compose is usually used in such a way that the ports inside the container are mapped to ephemeral ports on your computer. For example, a Postgres server may run inside the container using port 5432 but be mapped to a totally different port locally. The service connection will always discover and use the locally mapped port. |

Service connections are established by using the image name of the container. The following service connections are currently supported:

| Connection Details | Matched on |
|---|---|
| Containers named "symptoma/activemq", "apache/activemq-classic" or "apache/activemq" | |
| Containers named "apache/activemq-artemis" or "apache/artemis" | |
| Containers named "cassandra" | |
| Containers named "elasticsearch" or "elasticsearch/elasticsearch" | |
| Containers named "hazelcast/hazelcast". | |
| Containers named "clickhouse/clickhouse-server", "gvenzl/oracle-free", "gvenzl/oracle-xe", "mariadb", "mssql/server", "mysql", or "postgres" | |
| Containers named "osixia/openldap", "lldap/lldap" | |
| Containers named "mongo" | |
| Containers named "neo4j" | |
| Containers named "otel/opentelemetry-collector-contrib", "grafana/otel-lgtm" | |
| Containers named "otel/opentelemetry-collector-contrib", "grafana/otel-lgtm" | |
| Containers named "otel/opentelemetry-collector-contrib", "grafana/otel-lgtm" | |
| Containers named "apachepulsar/pulsar" | |
| Containers named "clickhouse/clickhouse-server", "gvenzl/oracle-free", "gvenzl/oracle-xe", "mariadb", "mssql/server", "mysql", or "postgres" | |
| Containers named "rabbitmq" with container port 5672 mapped | |
| Containers named "rabbitmq" with container port 5552 mapped | |
| Containers named "redis", "redis/redis-stack" or "redis/redis-stack-server" | |
| Containers named "openzipkin/zipkin". |

### SSL support

Some images come with SSL enabled out of the box, or maybe you want to enable SSL for the container to mirror your production setup. Spring Boot supports SSL configuration for supported service connections. Please note that you still have to enable SSL on the service which is running inside the container yourself, this feature only configures SSL on the client side in your application.

SSL is supported for the following service connections:

-
Cassandra
-
Elasticsearch
-
MongoDB
-
RabbitMQ
-
RabbitMQ Streams
-
Redis

To enable SSL support for a service, you can use service labels.

For JKS based keystores and truststores, you can use the following container labels:

-
`org.springframework.boot.sslbundle.jks.key.alias`
-
`org.springframework.boot.sslbundle.jks.key.password`
-
`org.springframework.boot.sslbundle.jks.options.ciphers`
-
`org.springframework.boot.sslbundle.jks.options.enabled-protocols`
-
`org.springframework.boot.sslbundle.jks.protocol`
-
`org.springframework.boot.sslbundle.jks.keystore.type`
-
`org.springframework.boot.sslbundle.jks.keystore.provider`
-
`org.springframework.boot.sslbundle.jks.keystore.location`
-
`org.springframework.boot.sslbundle.jks.keystore.password`
-
`org.springframework.boot.sslbundle.jks.truststore.type`
-
`org.springframework.boot.sslbundle.jks.truststore.provider`
-
`org.springframework.boot.sslbundle.jks.truststore.location`
-
`org.springframework.boot.sslbundle.jks.truststore.password`

These labels mirror the properties available for SSL bundles.

For PEM based keystores and truststores, you can use the following container labels:

-
`org.springframework.boot.sslbundle.pem.key.alias`
-
`org.springframework.boot.sslbundle.pem.key.password`
-
`org.springframework.boot.sslbundle.pem.options.ciphers`
-
`org.springframework.boot.sslbundle.pem.options.enabled-protocols`
-
`org.springframework.boot.sslbundle.pem.protocol`
-
`org.springframework.boot.sslbundle.pem.keystore.type`
-
`org.springframework.boot.sslbundle.pem.keystore.certificate`
-
`org.springframework.boot.sslbundle.pem.keystore.private-key`
-
`org.springframework.boot.sslbundle.pem.keystore.private-key-password`
-
`org.springframework.boot.sslbundle.pem.truststore.type`
-
`org.springframework.boot.sslbundle.pem.truststore.certificate`
-
`org.springframework.boot.sslbundle.pem.truststore.private-key`
-
`org.springframework.boot.sslbundle.pem.truststore.private-key-password`

These labels mirror the properties available for SSL bundles.

The following example enables SSL for a redis container:

```
services:
 redis:
 image: 'redis:latest'
 ports:
 - '6379'
 secrets:
 - ssl-ca
 - ssl-key
 - ssl-cert
 command: 'redis-server --tls-port 6379 --port 0 --tls-cert-file /run/secrets/ssl-cert --tls-key-file /run/secrets/ssl-key --tls-ca-cert-file /run/secrets/ssl-ca'
 labels:
 - 'org.springframework.boot.sslbundle.pem.keystore.certificate=client.crt'
 - 'org.springframework.boot.sslbundle.pem.keystore.private-key=client.key'
 - 'org.springframework.boot.sslbundle.pem.truststore.certificate=ca.crt'
secrets:
 ssl-ca:
 file: 'ca.crt'
 ssl-key:
 file: 'server.key'
 ssl-cert:
 file: 'server.crt'
```
### Custom Images

Sometimes you may need to use your own version of an image to provide a service. You can use any custom image as long as it behaves in the same way as the standard image. Specifically, any environment variables that the standard image supports must also be used in your custom image.

If your image uses a different name, you can use a label in your `compose.yml` file so that Spring Boot can provide a service connection.
Use a label named `org.springframework.boot.service-connection` to provide the service name.

For example:

```
services:
 redis:
 image: 'mycompany/mycustomredis:7.0'
 ports:
 - '6379'
 labels:
 org.springframework.boot.service-connection: redis
```
### Skipping Specific Containers

If you have a container image defined in your `compose.yml` that you don’t want connected to your application you can use a label to ignore it.
Any container with labeled with `org.springframework.boot.ignore` will be ignored by Spring Boot.

For example:

```
services:
 redis:
 image: 'redis:7.0'
 ports:
 - '6379'
 labels:
 org.springframework.boot.ignore: true
```
### Using a Specific Compose File

If your compose file is not in the same directory as your application, or if it’s named differently, you can use `spring.docker.compose.file` in your `application.properties` or `application.yaml` to point to a different file.
Properties can be defined as an exact path or a path that’s relative to your application.

For example:

-
Properties
-
YAML

`spring.docker.compose.file=../my-compose.yml````
spring:
 docker:
 compose:
 file: "../my-compose.yml"
```
### Waiting for Container Readiness

Containers started by Docker Compose may take some time to become fully ready.
The recommended way of checking for readiness is to add a `healthcheck` section under the service definition in your `compose.yml` file.

Since it’s not uncommon for `healthcheck` configuration to be omitted from `compose.yml` files, Spring Boot also checks directly for service readiness.
By default, a container is considered ready when a TCP/IP connection can be established to its mapped port.

You can disable this on a per-container basis by adding a `org.springframework.boot.readiness-check.tcp.disable` label in your `compose.yml` file.

For example:

```
services:
 redis:
 image: 'redis:7.0'
 ports:
 - '6379'
 labels:
 org.springframework.boot.readiness-check.tcp.disable: true
```
You can also change timeout values in your `application.properties` or `application.yaml` file:

-
Properties
-
YAML

```
spring.docker.compose.readiness.tcp.connect-timeout=10s
spring.docker.compose.readiness.tcp.read-timeout=5s
```
```
spring:
 docker:
 compose:
 readiness:
 tcp:
 connect-timeout: 10s
 read-timeout: 5s
```
The overall timeout can be configured using `spring.docker.compose.readiness.timeout`.

### Controlling the Docker Compose Lifecycle

By default Spring Boot calls `docker compose up` when your application starts and `docker compose stop` when it’s shut down.
If you prefer to have different lifecycle management you can use the `spring.docker.compose.lifecycle-management` property.

The following values are supported:

-
`none`- Do not start or stop Docker Compose
-
`start-only`- Start Docker Compose when the application starts and leave it running
-
`start-and-stop`- Start Docker Compose when the application starts and stop it when the JVM exits

In addition you can use the `spring.docker.compose.start.command` property to change whether `docker compose up` or `docker compose start` is used.
The `spring.docker.compose.stop.command` allows you to configure if `docker compose down` or `docker compose stop` is used.

You can also pass additional arguments to Docker Compose commands.
The `spring.docker.compose.arguments` property allows you to specify arguments that are passed to all Docker Compose commands.
The `spring.docker.compose.start.arguments` property allows you to specify arguments that are passed only to the up (or start) command, while the `spring.docker.compose.stop.arguments` property allows you to specify arguments that are passed only to the down (or stop) command.

The following example shows how lifecycle management can be configured:

-
Properties
-
YAML

```
spring.docker.compose.lifecycle-management=start-and-stop
spring.docker.compose.arguments[0]=--project-name=myapp
spring.docker.compose.arguments[1]=--progress=auto
spring.docker.compose.start.command=up
spring.docker.compose.start.arguments[0]=--build
spring.docker.compose.start.arguments[1]=--force-recreate
spring.docker.compose.stop.command=down
spring.docker.compose.stop.timeout=1m
spring.docker.compose.stop.arguments[0]=--volumes
spring.docker.compose.stop.arguments[1]=--remove-orphans
```
```
spring:
 docker:
 compose:
 lifecycle-management: start-and-stop
 arguments:
 - "--project-name=myapp"
 - "--progress=auto"
 start:
 command: up
 arguments:
 - "--build"
 - "--force-recreate"
 stop:
 command: down
 timeout: 1m
 arguments:
 - "--volumes"
 - "--remove-orphans"
```
### Activating Docker Compose Profiles

Docker Compose profiles are similar to Spring profiles in that they let you adjust your Docker Compose configuration for specific environments.
If you want to activate a specific Docker Compose profile you can use the `spring.docker.compose.profiles.active` property in your `application.properties` or `application.yaml` file:

-
Properties
-
YAML

`spring.docker.compose.profiles.active=myprofile````
spring:
 docker:
 compose:
 profiles:
 active: "myprofile"
```
### Using Docker Compose in Tests

By default, Spring Boot’s Docker Compose support is disabled when running tests.

To enable Docker Compose support in tests, set `spring.docker.compose.skip.in-tests` to `false`.

When using Gradle, you also need to change the configuration of the `spring-boot-docker-compose` dependency from `developmentOnly` to `testAndDevelopmentOnly`:

```
dependencies {
	testAndDevelopmentOnly("org.springframework.boot:spring-boot-docker-compose")
}
```
## Testcontainers Support

As well as using Testcontainers for integration testing, it’s also possible to use them at development time. The next sections will provide more details about that.

### Using Testcontainers at Development Time

This approach allows developers to quickly start containers for the services that the application depends on, removing the need to manually provision things like database servers. Using Testcontainers in this way provides functionality similar to Docker Compose, except that your container configuration is in Java rather than YAML.

To use Testcontainers at development time you need to launch your application using your “test” classpath rather than “main”. This will allow you to access all declared test dependencies and give you a natural place to write your test configuration.

To create a test launchable version of your application you should create an “Application” class in the `src/test` directory.
For example, if your main application is in `src/main/java/com/example/MyApplication.java`, you should create `src/test/java/com/example/TestMyApplication.java`

The `TestMyApplication` class can use the `SpringApplication.from(…)` method to launch the real application:

-
Java
-
Kotlin

```
import org.springframework.boot.SpringApplication;
public class TestMyApplication {
	public static void main(String[] args) {
 SpringApplication.from(MyApplication::main).run(args);
	}
}
```
```
import org.springframework.boot.fromApplication
fun main(args: Array<String>) {
	fromApplication<MyApplication>().run(*args)
}
```
You’ll also need to define the `Container` instances that you want to start along with your application.
To do this, you need to make sure that the `spring-boot-testcontainers` module has been added as a `test` dependency.
Once that has been done, you can create a `@TestConfiguration` class that declares `@Bean` methods for the containers you want to start.

You can also annotate your `@Bean` methods with `@ServiceConnection` in order to create `ConnectionDetails` beans.
See the service connections section for details of the supported technologies.

A typical Testcontainers configuration would look like this:

-
Java
-
Kotlin

```
import org.testcontainers.neo4j.Neo4jContainer;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.context.annotation.Bean;
@TestConfiguration(proxyBeanMethods = false)
public class MyContainersConfiguration {
	@Bean
	@ServiceConnection
	public Neo4jContainer neo4jContainer() {
 return new Neo4jContainer("neo4j:5");
	}
}
```
```
import org.springframework.boot.test.context.TestConfiguration
import org.springframework.boot.testcontainers.service.connection.ServiceConnection
import org.springframework.context.annotation.Bean
import org.testcontainers.neo4j.Neo4jContainer
@TestConfiguration(proxyBeanMethods = false)
class MyContainersConfiguration {
	@Bean
	@ServiceConnection
	fun neo4jContainer(): Neo4jContainer {
 return Neo4jContainer("neo4j:5")
	}
}
```
| The lifecycle of `Container`beans is automatically managed by Spring Boot.
Containers will be started and stopped automatically. |

| You can use the `spring.testcontainers.beans.startup`property to change how containers are started.
By default`sequential`startup is used, but you may also choose`parallel`if you wish to start multiple containers in parallel. |

Once you have defined your test configuration, you can use the `with(…)` method to attach it to your test launcher:

-
Java
-
Kotlin

```
import org.springframework.boot.SpringApplication;
public class TestMyApplication {
	public static void main(String[] args) {
 SpringApplication.from(MyApplication::main).with(MyContainersConfiguration.class).run(args);
	}
}
```
```
import org.springframework.boot.fromApplication
import org.springframework.boot.with
fun main(args: Array<String>) {
	fromApplication<MyApplication>().with(MyContainersConfiguration::class).run(*args)
}
```
You can now launch `TestMyApplication` as you would any regular Java `main` method application to start your application and the containers that it needs to run.

| You can use the Maven goal `spring-boot:test-run`or the Gradle task`bootTestRun`to do this from the command line. |

#### Contributing Dynamic Properties at Development Time

If you want to contribute dynamic properties at development time from your `Container` `@Bean` methods, define an additional `DynamicPropertyRegistrar` bean.
The registrar should be defined using a `@Bean` method that injects the container from which the properties will be sourced as a parameter.
This arrangement ensures that container has been started before the properties are used.

A typical configuration would look like this:

-
Java
-
Kotlin

```
import org.testcontainers.mongodb.MongoDBContainer;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.test.context.DynamicPropertyRegistrar;
@TestConfiguration(proxyBeanMethods = false)
public class MyContainersConfiguration {
	@Bean
	public MongoDBContainer mongoDbContainer() {
 return new MongoDBContainer("mongo:5.0");
	}
	@Bean
	public DynamicPropertyRegistrar mongoDbProperties(MongoDBContainer container) {
 return (properties) -> {
 properties.add("spring.mongodb.host", container::getHost);
 properties.add("spring.mongodb.port", container::getFirstMappedPort);
 };
	}
}
```
```
import org.springframework.boot.test.context.TestConfiguration
import org.springframework.context.annotation.Bean
import org.springframework.test.context.DynamicPropertyRegistrar;
import org.testcontainers.mongodb.MongoDBContainer
@TestConfiguration(proxyBeanMethods = false)
class MyContainersConfiguration {
	@Bean
	fun mongoDbContainer(): MongoDBContainer {
 return MongoDBContainer("mongo:5.0")
	}

	@Bean
	fun mongoDbProperties(container: MongoDBContainer): DynamicPropertyRegistrar {
 return DynamicPropertyRegistrar { properties ->
 properties.add("spring.mongodb.host") { container.host }
 properties.add("spring.mongodb.port") { container.firstMappedPort }
 }
	}
}
```
| Using a `@ServiceConnection`is recommended whenever possible, however, dynamic properties can be a useful fallback for technologies that don’t yet have`@ServiceConnection`support. |

#### Importing Testcontainers Declaration Classes

A common pattern when using Testcontainers is to declare `Container` instances as static fields.
Often these fields are defined directly on the test class.
They can also be declared on a parent class or on an interface that the test implements.

For example, the following `MyContainers` interface declares `mongo` and `neo4j` containers:

-
Java
-
Kotlin

```
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.mongodb.MongoDBContainer;
import org.testcontainers.neo4j.Neo4jContainer;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
public interface MyContainers {
	@Container
	@ServiceConnection
	MongoDBContainer mongoContainer = new MongoDBContainer("mongo:5.0");
	@Container
	@ServiceConnection
	Neo4jContainer neo4jContainer = new Neo4jContainer("neo4j:5");
}
```
```
import org.springframework.boot.testcontainers.service.connection.ServiceConnection
import org.testcontainers.junit.jupiter.Container
import org.testcontainers.mongodb.MongoDBContainer
import org.testcontainers.neo4j.Neo4jContainer
interface MyContainers {
	companion object {
 @Container
 @ServiceConnection
 @JvmField
 val mongoContainer = MongoDBContainer("mongo:5.0")
 @Container
 @ServiceConnection
 @JvmField
 val neo4jContainer = Neo4jContainer("neo4j:5")
	}
}
```
If you already have containers defined in this way, or you just prefer this style, you can import these declaration classes rather than defining your containers as `@Bean` methods.
To do so, add the `@ImportTestcontainers` annotation to your test configuration class:

-
Java
-
Kotlin

```
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.boot.testcontainers.context.ImportTestcontainers;
@TestConfiguration(proxyBeanMethods = false)
@ImportTestcontainers(MyContainers.class)
public class MyContainersConfiguration {
}
```
```
import org.springframework.boot.test.context.TestConfiguration
import org.springframework.boot.testcontainers.context.ImportTestcontainers
@TestConfiguration(proxyBeanMethods = false)
@ImportTestcontainers(MyContainers::class)
class MyContainersConfiguration
```
| If you don’t intend to use the service connections feature but want to use `@DynamicPropertySource`instead, remove the`@ServiceConnection`annotation from the`Container`fields.
You can also add`@DynamicPropertySource`annotated methods to your declaration class. |

#### Using DevTools with Testcontainers at Development Time

When using devtools, you can annotate beans and bean methods with `@RestartScope`.
Such beans won’t be recreated when the devtools restart the application.
This is especially useful for `Container` beans, as they keep their state despite the application restart.

-
Java
-
Kotlin

```
import org.testcontainers.mongodb.MongoDBContainer;
import org.springframework.boot.devtools.restart.RestartScope;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.context.annotation.Bean;
@TestConfiguration(proxyBeanMethods = false)
public class MyContainersConfiguration {
	@Bean
	@RestartScope
	@ServiceConnection
	public MongoDBContainer mongoDbContainer() {
 return new MongoDBContainer("mongo:5.0");
	}
}
```
```
import org.springframework.boot.devtools.restart.RestartScope
import org.springframework.boot.test.context.TestConfiguration
import org.springframework.boot.testcontainers.service.connection.ServiceConnection
import org.springframework.context.annotation.Bean
import org.testcontainers.mongodb.MongoDBContainer
@TestConfiguration(proxyBeanMethods = false)
class MyContainersConfiguration {
	@Bean
	@RestartScope
	@ServiceConnection
	fun mongoDbContainer(): MongoDBContainer {
 return MongoDBContainer("mongo:5.0")
	}
}
```
| If you’re using Gradle and want to use this feature, you need to change the configuration of the `spring-boot-devtools`dependency from`developmentOnly`to`testAndDevelopmentOnly`.
With the default scope of`developmentOnly`, the`bootTestRun`task will not pick up changes in your code, as the devtools are not active. |

# Creating Your Own Auto-configuration

If you work in a company that develops shared libraries, or if you work on an open-source or commercial library, you might want to develop your own auto-configuration. Auto-configuration classes can be bundled in external jars and still be picked up by Spring Boot.

Auto-configuration can be associated to a “starter” that provides the auto-configuration code as well as the typical libraries that you would use with it. We first cover what you need to know to build your own auto-configuration and then we move on to the typical steps required to create a custom starter.

## Understanding Auto-configured Beans

Classes that implement auto-configuration are annotated with `@AutoConfiguration`.
This annotation itself is meta-annotated with `@Configuration`, making auto-configurations standard `@Configuration` classes.
Additional `@Conditional` annotations are used to constrain when the auto-configuration should apply.
Usually, auto-configuration classes use `@ConditionalOnClass` and `@ConditionalOnMissingBean` annotations.
This ensures that auto-configuration applies only when relevant classes are found and when you have not declared your own `@Configuration`.

You can browse the source code of `spring-boot-autoconfigure` to see the core `@AutoConfiguration` classes that Spring provides (see the `META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.imports` file).
You can also look at the equivalent file in other modules to see the auto-configurations that they provide.

## Locating Auto-configuration Candidates

Spring Boot checks for the presence of a `META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.imports` file within your published jar.
The file should list your configuration classes, with one class name per line, as shown in the following example:

```
com.mycorp.libx.autoconfigure.LibXAutoConfiguration
com.mycorp.libx.autoconfigure.LibXWebAutoConfiguration
```
| You can add comments to the imports file using the `#`character. |

| In the unusual case that an auto-configuration class is not a top-level class, its class name should use `$`to separate it from its containing class, for example`com.example.Outer$NestedAutoConfiguration`. |

| Auto-configurations must be loaded onlyby being named in the imports file.
Make sure that they are defined in a specific package space and that they are never the target of component scanning.
Furthermore, auto-configuration classes should not enable component scanning to find additional components.
Specific`@Import`annotations should be used instead. |

If your configuration needs to be applied in a specific order, you can use the `before`, `beforeName`, `after` and `afterName` attributes on the `@AutoConfiguration` annotation or the dedicated `@AutoConfigureBefore` and `@AutoConfigureAfter` annotations.
For example, if you provide web-specific configuration, your class may need to be applied after `WebMvcAutoConfiguration`.

If you want to order certain auto-configurations that should not have any direct knowledge of each other, you can also use `@AutoConfigureOrder`.
That annotation has the same semantic as the regular `@Order` annotation but provides a dedicated order for auto-configuration classes.

As with standard `@Configuration` classes, the order in which auto-configuration classes are applied only affects the order in which their beans are defined.
The order in which those beans are subsequently created is unaffected and is determined by each bean’s dependencies and any `@DependsOn` relationships.

### Deprecating and Replacing Auto-configuration Classes

You may need to occasionally deprecate auto-configuration classes and offer an alternative. For example, you may want to change the package name where your auto-configuration class resides.

Since auto-configuration classes may be referenced in `before`/`after` ordering and `excludes`, you’ll need to add an additional file that tells Spring Boot how to deal with replacements.
To define replacements, create a `META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.replacements` file indicating the link between the old class and the new one.

For example:

`com.mycorp.libx.autoconfigure.LibXAutoConfiguration=com.mycorp.libx.autoconfigure.core.LibXAutoConfiguration`| The `AutoConfiguration.imports`file should also be updated toonlyreference the replacement class. |

## Condition Annotations

You almost always want to include one or more `@Conditional` annotations on your auto-configuration class.
The `@ConditionalOnMissingBean` annotation is one common example that is used to allow developers to override auto-configuration if they are not happy with your defaults.

Spring Boot includes a number of `@Conditional` annotations that you can reuse in your own code by annotating `@Configuration` classes or individual `@Bean` methods.
These annotations include:

### Class Conditions

The `@ConditionalOnClass` and `@ConditionalOnMissingClass` annotations let `@Configuration` classes be included based on the presence or absence of specific classes.
Due to the fact that annotation metadata is parsed by using ASM, you can use the `value` attribute to refer to the real class, even though that class might not actually appear on the running application classpath.
You can also use the `name` attribute if you prefer to specify the class name by using a `String` value.

This mechanism does not apply the same way to `@Bean` methods where typically the return type is the target of the condition: before the condition on the method applies, the JVM will have loaded the class and potentially processed method references which will fail if the class is not present.

To handle this scenario, a separate `@Configuration` class can be used to isolate the condition, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.autoconfigure.AutoConfiguration;
import org.springframework.boot.autoconfigure.condition.ConditionalOnClass;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@AutoConfiguration
// Some conditions ...
public final class MyAutoConfiguration {
	// Auto-configured beans ...
	@Configuration(proxyBeanMethods = false)
	@ConditionalOnClass(SomeService.class)
	static class SomeServiceConfiguration {
 @Bean
 @ConditionalOnMissingBean
 SomeService someService() {
 return new SomeService();
 }
	}
}
```
```
import org.springframework.boot.autoconfigure.AutoConfiguration
import org.springframework.boot.autoconfigure.condition.ConditionalOnClass
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
@AutoConfiguration
// Some conditions ...
class MyAutoConfiguration {
	// Auto-configured beans ...
	@Configuration(proxyBeanMethods = false)
	@ConditionalOnClass(SomeService::class)
	class SomeServiceConfiguration {
 @Bean
 @ConditionalOnMissingBean
 fun someService(): SomeService {
 return SomeService()
 }
	}
}
```
| If you use `@ConditionalOnClass`or`@ConditionalOnMissingClass`as a part of a meta-annotation to compose your own composed annotations, you must use`name`as referring to the class in such a case is not handled. |

### Bean Conditions

The `@ConditionalOnBean` and `@ConditionalOnMissingBean` annotations let a bean be included based on the presence or absence of specific beans.
You can use the `value` attribute to specify beans by type or `name` to specify beans by name.
The `search` attribute lets you limit the `ApplicationContext` hierarchy that should be considered when searching for beans.

When placed on a `@Bean` method, the target type defaults to the return type of the method, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.autoconfigure.AutoConfiguration;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.context.annotation.Bean;
@AutoConfiguration
public final class MyAutoConfiguration {
	@Bean
	@ConditionalOnMissingBean
	SomeService someService() {
 return new SomeService();
	}
}
```
```
import org.springframework.boot.autoconfigure.AutoConfiguration
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean
import org.springframework.context.annotation.Bean
@AutoConfiguration
class MyAutoConfiguration {
	@Bean
	@ConditionalOnMissingBean
	fun someService(): SomeService {
 return SomeService()
	}
}
```
In the preceding example, the `someService` bean is going to be created if no bean of type `SomeService` is already contained in the `ApplicationContext`.

| You need to be very careful about the order in which bean definitions are added, as these conditions are evaluated based on what has been processed so far.
For this reason, we recommend using only `@ConditionalOnBean`and`@ConditionalOnMissingBean`annotations on auto-configuration classes (since these are guaranteed to load after any user-defined bean definitions have been added). |

| `@ConditionalOnBean`and`@ConditionalOnMissingBean`do not prevent`@Configuration`classes from being created.
The only difference between using these conditions at the class level and marking each contained`@Bean`method with the annotation is that the former prevents registration of the`@Configuration`class as a bean if the condition does not match. |

| When declaring a `@Bean`method, provide as much type information as possible in the method’s return type.
For example, if your bean’s concrete class implements an interface the bean method’s return type should be the concrete class and not the interface.
Providing as much type information as possible in`@Bean`methods is particularly important when using bean conditions as their evaluation can only rely upon to type information that is available in the method signature. |

### Property Conditions

The `@ConditionalOnProperty` annotation lets configuration be included based on a Spring Environment property.
Use the `prefix` and `name` attributes to specify the property that should be checked.
By default, any property that exists and is not equal to `false` is matched.
There is also a dedicated `@ConditionalOnBooleanProperty` annotation specifically made for boolean properties.
With both annotations you can also create more advanced checks by using the `havingValue` and `matchIfMissing` attributes.

If multiple names are given in the `name` attribute, all of the properties have to pass the test for the condition to match.

### Resource Conditions

The `@ConditionalOnResource` annotation lets configuration be included only when a specific resource is present.
Resources can be specified by using the usual Spring conventions, as shown in the following example: `file:/home/user/test.dat`.

### Web Application Conditions

The `@ConditionalOnWebApplication` and `@ConditionalOnNotWebApplication` annotations let configuration be included depending on whether the application is a web application.
A servlet-based web application is any application that uses a Spring `WebApplicationContext`, defines a `session` scope, or has a `ConfigurableWebEnvironment`.
A reactive web application is any application that uses a `ReactiveWebApplicationContext`, or has a `ConfigurableReactiveWebEnvironment`.

The `@ConditionalOnWarDeployment` and `@ConditionalOnNotWarDeployment` annotations let configuration be included depending on whether the application is a traditional WAR application that is deployed to a servlet container.
This condition will not match for applications that are run with an embedded web server.

### SpEL Expression Conditions

The `@ConditionalOnExpression` annotation lets configuration be included based on the result of a SpEL expression.

| Referencing a bean in the expression will cause that bean to be initialized very early in context refresh processing. As a result, the bean won’t be eligible for post-processing (such as configuration properties binding) and its state may be incomplete. |

## Testing your Auto-configuration

An auto-configuration can be affected by many factors: user configuration (`@Bean` definition and `Environment` customization), condition evaluation (presence of a particular library), and others.
Concretely, each test should create a well defined `ApplicationContext` that represents a combination of those customizations.
`ApplicationContextRunner` provides a great way to achieve that.

| `ApplicationContextRunner`doesn’t work when running the tests in a native image. |

`ApplicationContextRunner` is usually defined as a field of the test class to gather the base, common configuration.
The following example makes sure that `MyServiceAutoConfiguration` is always invoked:

-
Java
-
Kotlin

```
	private final ApplicationContextRunner contextRunner = new ApplicationContextRunner()
 .withConfiguration(AutoConfigurations.of(MyServiceAutoConfiguration.class));
```
```
	val contextRunner = ApplicationContextRunner()
 .withConfiguration(AutoConfigurations.of(MyServiceAutoConfiguration::class.java))
```
| If multiple auto-configurations have to be defined, there is no need to order their declarations as they are invoked in the exact same order as when running the application. |

Each test can use the runner to represent a particular use case.
For instance, the sample below invokes a user configuration (`UserConfiguration`) and checks that the auto-configuration backs off properly.
Invoking `run` provides a callback context that can be used with AssertJ.

-
Java
-
Kotlin

```
	@Test
	void defaultServiceBacksOff() {
 this.contextRunner.withUserConfiguration(UserConfiguration.class).run((context) -> {
 assertThat(context).hasSingleBean(MyService.class);
 assertThat(context).getBean("myCustomService").isSameAs(context.getBean(MyService.class));
 });
	}
	@Configuration(proxyBeanMethods = false)
	static class UserConfiguration {
 @Bean
 MyService myCustomService() {
 return new MyService("mine");
 }
	}
```
```
	@Test
	fun defaultServiceBacksOff() {
 contextRunner.withUserConfiguration(UserConfiguration::class.java)
 .run { context: AssertableApplicationContext ->
 assertThat(context).hasSingleBean(MyService::class.java)
 assertThat(context).getBean("myCustomService")
 .isSameAs(context.getBean(MyService::class.java))
 }
	}
	@Configuration(proxyBeanMethods = false)
	internal class UserConfiguration {
 @Bean
 fun myCustomService(): MyService {
 return MyService("mine")
 }
	}
```
It is also possible to easily customize the `Environment`, as shown in the following example:

-
Java
-
Kotlin

```
	@Test
	void serviceNameCanBeConfigured() {
 this.contextRunner.withPropertyValues("user.name=test123").run((context) -> {
 assertThat(context).hasSingleBean(MyService.class);
 assertThat(context.getBean(MyService.class).getName()).isEqualTo("test123");
 });
	}
```
```
	@Test
	fun serviceNameCanBeConfigured() {
 contextRunner.withPropertyValues("user.name=test123").run { context: AssertableApplicationContext ->
 assertThat(context).hasSingleBean(MyService::class.java)
 assertThat(context.getBean(MyService::class.java).name).isEqualTo("test123")
 }
	}
```
The runner can also be used to display the `ConditionEvaluationReport`.
The report can be printed at `INFO` or `DEBUG` level.
The following example shows how to use the `ConditionEvaluationReportLoggingListener` to print the report in auto-configuration tests.

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.boot.autoconfigure.logging.ConditionEvaluationReportLoggingListener;
import org.springframework.boot.logging.LogLevel;
import org.springframework.boot.test.context.runner.ApplicationContextRunner;
class MyConditionEvaluationReportingTests {
	@Test
	void autoConfigTest() {
 new ApplicationContextRunner()
 .withInitializer(ConditionEvaluationReportLoggingListener.forLogLevel(LogLevel.INFO))
 .run((context) -> {
 // Test something...
 });
	}
}
```
```
import org.junit.jupiter.api.Test
import org.springframework.boot.autoconfigure.logging.ConditionEvaluationReportLoggingListener
import org.springframework.boot.logging.LogLevel
import org.springframework.boot.test.context.assertj.AssertableApplicationContext
import org.springframework.boot.test.context.runner.ApplicationContextRunner
class MyConditionEvaluationReportingTests {
	@Test
	fun autoConfigTest() {
 ApplicationContextRunner()
 .withInitializer(ConditionEvaluationReportLoggingListener.forLogLevel(LogLevel.INFO))
 .run { context: AssertableApplicationContext? ->
 // Test something...
 }
	}
}
```
### Simulating a Web Context

If you need to test an auto-configuration that only operates in a servlet or reactive web application context, use the `WebApplicationContextRunner` or `ReactiveWebApplicationContextRunner` respectively.

### Overriding the Classpath

It is also possible to test what happens when a particular class and/or package is not present at runtime.
Spring Boot ships with a `FilteredClassLoader` that can easily be used by the runner.
In the following example, we assert that if `MyService` is not present, the auto-configuration is properly disabled:

-
Java
-
Kotlin

```
	@Test
	void serviceIsIgnoredIfLibraryIsNotPresent() {
 this.contextRunner.withClassLoader(new FilteredClassLoader(MyService.class))
 .run((context) -> assertThat(context).doesNotHaveBean("myService"));
	}
```
```
	@Test
	fun serviceIsIgnoredIfLibraryIsNotPresent() {
 contextRunner.withClassLoader(FilteredClassLoader(MyService::class.java))
 .run { context: AssertableApplicationContext? ->
 assertThat(context).doesNotHaveBean("myService")
 }
	}
```
## Creating Your Own Starter

A typical Spring Boot starter contains code to auto-configure and customize the infrastructure of a given technology, let’s call that "acme". To make it easily extensible, a number of configuration keys in a dedicated namespace can be exposed to the environment. Finally, a single "starter" dependency is provided to help users get started as easily as possible.

Concretely, a custom starter can contain the following:

-
The `acme-spring-boot`module that contains the auto-configuration code for "acme" as well as any API to use the feature.
-
The `acme-spring-boot-starter`module that provides a dependency to the other starters required by "acme",`acme-spring-boot`, and potentially additional dependencies that are typically useful. In a nutshell, adding the starter should provide everything needed to start using that library.

This separation in two modules is in no way necessary.
If "acme" has several flavors, options or optional features, then it is better to separate the auto-configuration as you can clearly express the fact some features are optional.
Besides, you have the ability to craft a starter that provides an opinion about those optional dependencies.
At the same time, others can rely only on `acme-spring-boot` and craft their own starter with different opinions.

If the auto-configuration is relatively straightforward and does not have optional features, merging the two modules in the starter is definitely an option.

If "acme" has a basic set of dependencies that are required to work, but you want to express a more opinionated view, then having a separate starter is a good option.

When testing "acme" features, you may need test-specific auto-configurations. For example, you could offer a way to replace external dependencies with in-memory alternatives. A separate test-scoped starter can be created for this purpose, following the same principles.

### Naming

You should make sure to provide a proper namespace for your starter.
Do not start your module names with `spring-boot`, even if you use a different Maven `groupId`.
We may offer official support for the thing you auto-configure in the future.

As a rule of thumb, you should name a combined module after the starter.
For example, assume that you are creating a starter for "acme" and that you name the auto-configure module `acme-spring-boot` and the starter `acme-spring-boot-starter`.
If you only have one module that combines the two, name it `acme-spring-boot-starter`.

If "acme" also has a test-scoped starter, name it `acme-spring-boot-starter-test`.

### Configuration keys

If your starter provides configuration keys, use a unique namespace for them.
In particular, do not include your keys in the namespaces that Spring Boot uses (such as `server`, `management`, `spring`, and so on).
If you use the same namespace, we may modify these namespaces in the future in ways that break your modules.
As a rule of thumb, prefix all your keys with a namespace that you own (for example `acme`).

Make sure that configuration keys are documented by adding field Javadoc for each property, as shown in the following example:

-
Java
-
Kotlin

```
import java.time.Duration;
import org.springframework.boot.context.properties.ConfigurationProperties;
@ConfigurationProperties("acme")
public class AcmeProperties {
	/**
 * Whether to check the location of acme resources.
 */
	private boolean checkLocation = true;
	/**
 * Timeout for establishing a connection to the acme server.
 */
	private Duration loginTimeout = Duration.ofSeconds(3);
	// getters/setters ...
	public boolean isCheckLocation() {
 return this.checkLocation;
	}
	public void setCheckLocation(boolean checkLocation) {
 this.checkLocation = checkLocation;
	}
	public Duration getLoginTimeout() {
 return this.loginTimeout;
	}
	public void setLoginTimeout(Duration loginTimeout) {
 this.loginTimeout = loginTimeout;
	}
}
```
```
import org.springframework.boot.context.properties.ConfigurationProperties
import java.time.Duration
@ConfigurationProperties("acme")
class AcmeProperties(
	/**
 * Whether to check the location of acme resources.
 */
	var isCheckLocation: Boolean = true,
	/**
 * Timeout for establishing a connection to the acme server.
 */
	var loginTimeout:Duration = Duration.ofSeconds(3))
```
| You should only use plain text with `@ConfigurationProperties`field Javadoc, since they are not processed before being added to the JSON. |

If you use `@ConfigurationProperties` with record class then record components' descriptions should be provided via class-level Javadoc tag `@param` (there are no explicit instance fields in record classes to put regular field-level Javadocs on).

Here are some rules we follow internally to make sure descriptions are consistent:

-
Do not start the description by "The" or "A".
-
For `boolean`types, start the description with "Whether" or "Enable".
-
For collection-based types, start the description with "Comma-separated list"
-
Use `Duration`rather than`long`and describe the default unit if it differs from milliseconds, such as "If a duration suffix is not specified, seconds will be used".
-
Do not provide the default value in the description unless it has to be determined at runtime.

Make sure to trigger meta-data generation so that IDE assistance is available for your keys as well.
You may want to review the generated metadata (`META-INF/spring-configuration-metadata.json`) to make sure your keys are properly documented.
Using your own starter in a compatible IDE is also a good idea to validate that quality of the metadata.

### The “autoconfigure” Module

The `autoconfigure` module contains everything that is necessary to get started with the library.
It may also contain configuration key definitions (such as `@ConfigurationProperties`) and any callback interface that can be used to further customize how the components are initialized.

| You should mark the dependencies to the library as optional so that you can include the `autoconfigure`module in your projects more easily.
If you do it that way, the library is not provided and, by default, Spring Boot backs off. |

Spring Boot uses an annotation processor to collect the conditions on auto-configurations in a metadata file (`META-INF/spring-autoconfigure-metadata.properties`).
If that file is present, it is used to eagerly filter auto-configurations that do not match, which will improve startup time.

When building with Maven, configure the compiler plugin (3.12.0 or later) to add `spring-boot-autoconfigure-processor` to the annotation processor paths:

```
<project>
	<build>
 <plugins>
 <plugin>
 <groupId>org.apache.maven.plugins</groupId>
 <artifactId>maven-compiler-plugin</artifactId>
 <configuration>
 <annotationProcessorPaths>
 <path>
 <groupId>org.springframework.boot</groupId>
 <artifactId>spring-boot-autoconfigure-processor</artifactId>
 </path>
 </annotationProcessorPaths>
 </configuration>
 </plugin>
 </plugins>
	</build>
</project>
```
With Gradle, the dependency should be declared in the `annotationProcessor` configuration, as shown in the following example:

```
dependencies {
	annotationProcessor "org.springframework.boot:spring-boot-autoconfigure-processor"
}
```
### Starter Module

The starter is really an empty jar. Its only purpose is to provide the necessary dependencies to work with the library. You can think of it as an opinionated view of what is required to get started.

Do not make assumptions about the project in which your starter is added.
If the library you are auto-configuring typically requires other starters, mention them as well.
Providing a proper set of *default* dependencies may be hard if the number of optional dependencies is high, as you should avoid including dependencies that are unnecessary for a typical usage of the library.
In other words, you should not include optional dependencies.

| Either way, your starter must reference the core Spring Boot starter ( `spring-boot-starter`) directly or indirectly (there is no need to add it if your starter relies on another starter).
If a project is created with only your custom starter, Spring Boot’s core features will be honoured by the presence of the core starter. |

# Kotlin Support

Kotlin is a statically-typed language targeting the JVM (and other platforms) which allows writing concise and elegant code while providing interoperability with existing libraries written in Java.

Spring Boot provides Kotlin support by leveraging the support in other Spring projects such as Spring Framework, Spring Data, and Reactor. See the Spring Framework Kotlin support documentation for more information.

The easiest way to start with Spring Boot and Kotlin is to follow this comprehensive tutorial.
You can create new Kotlin projects by using start.spring.io.
Feel free to join the #spring channel of Kotlin Slack or ask a question with the `spring` and `kotlin` tags on Stack Overflow if you need support.

## Requirements

Spring Boot requires at least Kotlin 2.2.x and manages a suitable Kotlin version through dependency management.
To use Kotlin, `org.jetbrains.kotlin:kotlin-stdlib` and `org.jetbrains.kotlin:kotlin-reflect` must be present on the classpath.

Kotlin 2.2.x introduces new defaulting rules for propagating annotations to parameters, fields, and properties. In order to avoid related warnings and use what will likely become the Kotlin default behavior in an upcoming version, it is recommended to configure the `-Xannotation-default-target=param-property` compiler flag.

Since Kotlin classes are final by default, you are likely to want to configure kotlin-spring plugin in order to automatically open Spring-annotated classes so that they can be proxied.

Jackson’s Kotlin module is required for serializing / deserializing JSON data in Kotlin. It is automatically registered when found on the classpath.

| These dependencies and plugins are provided by default if one bootstraps a Kotlin project on start.spring.io. |

## Null-safety

One of Kotlin’s key features is null-safety.
It deals with `null` values at compile time rather than deferring the problem to runtime and encountering a `NullPointerException`.
This helps to eliminate a common source of bugs without paying the cost of wrappers like `Optional`.
Kotlin also allows using functional constructs with nullable values as described in this comprehensive guide to null-safety in Kotlin.

Although Java does not let you express null-safety in its type-system, most Spring projects provide null-safety via JSpecify annotations.

As of Kotlin 2.1, Kotlin enforces strict handling of nullability annotations from the `org.jspecify.annotations` package.

## Kotlin API

### runApplication

Spring Boot provides an idiomatic way to run an application with `runApplication<MyApplication>(*args)` as shown in the following example:

```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
@SpringBootApplication
class MyApplication
fun main(args: Array<String>) {
	runApplication<MyApplication>(*args)
}
```
This is a drop-in replacement for `SpringApplication.run(MyApplication::class.java, *args)`.
It also allows customization of the application as shown in the following example:

```
runApplication<MyApplication>(*args) {
	setBannerMode(OFF)
}
```
### Extensions

Kotlin extensions provide the ability to extend existing classes with additional functionality. The Spring Boot Kotlin API makes use of these extensions to add new Kotlin specific conveniences to existing APIs.

`TestRestTemplate` extensions, similar to those provided by Spring Framework for `RestOperations` in Spring Framework, are provided.
Among other things, the extensions make it possible to take advantage of Kotlin reified type parameters.

## Dependency Management

In order to avoid mixing different versions of Kotlin dependencies on the classpath, Spring Boot imports the Kotlin BOM.

With Maven, the Kotlin version can be customized by setting the `kotlin.version` property and plugin management is provided for `kotlin-maven-plugin`.
With Gradle, the Spring Boot plugin automatically aligns the `kotlin.version` with the version of the Kotlin plugin.

Spring Boot also manages the version of Coroutines dependencies by importing the Kotlin Coroutines BOM.
The version can be customized by setting the `kotlin-coroutines.version` property.

| `org.jetbrains.kotlinx:kotlinx-coroutines-reactor`dependency is provided by default if one bootstraps a Kotlin project with at least one reactive dependency on start.spring.io. |

## @ConfigurationProperties

`@ConfigurationProperties` when used in combination with constructor binding supports data classes with immutable `val` properties as shown in the following example:

```
@ConfigurationProperties("example.kotlin")
data class KotlinExampleProperties(
 val name: String,
 val description: String,
 val myService: MyService
) {
	data class MyService(
 val apiToken: String,
 val uri: URI
	)
}
```
Due to the limitations of their interoperability with Java, support for value classes is limited. In particular, relying upon a value class’s default value will not work with configuration property binding. In such cases, a data class should be used instead.

| To generate your own metadata using the annotation processor, `kapt`should be configured with the`spring-boot-configuration-processor`dependency.
Note that some features (such as detecting the default value or deprecated items) are not working due to limitations in the model kapt provides. |

## Testing

While it is possible to use JUnit 4 to test Kotlin code, JUnit 6 is provided by default and is recommended.
JUnit 6 enables a test class to be instantiated once and reused for all of the class’s tests.
This makes it possible to use `@BeforeAll` and `@AfterAll` annotations on non-static methods, which is a good fit for Kotlin.

To mock Kotlin classes, MockK is recommended.
If you need the `MockK` equivalent of the Mockito specific `@MockitoBean` and `@MockitoSpyBean` annotations, you can use SpringMockK which provides similar `@MockkBean` and `@SpykBean` annotations.

## Resources

### Further Reading

-
Kotlin Slack (with a dedicated #spring channel)
-
Tutorial: building web applications with Spring Boot and Kotlin
-
A Geospatial Messenger with Kotlin, Spring Boot and PostgreSQL

### Examples

-
spring-boot-kotlin-demo: regular Spring Boot + Spring Data JPA project
-
mixit: Spring Boot 2 + WebFlux + Reactive Spring Data MongoDB
-
spring-kotlin-fullstack: WebFlux Kotlin fullstack example with Kotlin2js for frontend instead of JavaScript or TypeScript
-
spring-petclinic-kotlin: Kotlin version of the Spring PetClinic Sample Application
-
spring-kotlin-deepdive: a step by step migration for Boot 1.0 + Java to Boot 2.0 + Kotlin
-
spring-boot-coroutines-demo: Coroutines sample project

# SSL

Spring Boot provides the ability to configure SSL trust material that can be applied to several types of connections in order to support secure communications.
Configuration properties with the prefix `spring.ssl.bundle` can be used to specify named sets of trust material and associated information.

## Configuring SSL With Java KeyStore Files

Configuration properties with the prefix `spring.ssl.bundle.jks` can be used to configure bundles of trust material created with the Java `keytool` utility and stored in Java KeyStore files in the JKS or PKCS12 format.
Each bundle has a user-provided name that can be used to reference the bundle.

When used to secure an embedded web server, a `keystore` is typically configured with a Java KeyStore containing a certificate and private key as shown in this example:

-
Properties
-
YAML

```
spring.ssl.bundle.jks.mybundle.key.alias=application
spring.ssl.bundle.jks.mybundle.keystore.location=classpath:application.p12
spring.ssl.bundle.jks.mybundle.keystore.password=secret
spring.ssl.bundle.jks.mybundle.keystore.type=PKCS12
```
```
spring:
 ssl:
 bundle:
 jks:
 mybundle:
 key:
 alias: "application"
 keystore:
 location: "classpath:application.p12"
 password: "secret"
 type: "PKCS12"
```
When used to secure a client-side connection, a `truststore` is typically configured with a Java KeyStore containing the server certificate as shown in this example:

-
Properties
-
YAML

```
spring.ssl.bundle.jks.mybundle.truststore.location=classpath:server.p12
spring.ssl.bundle.jks.mybundle.truststore.password=secret
```
```
spring:
 ssl:
 bundle:
 jks:
 mybundle:
 truststore:
 location: "classpath:server.p12"
 password: "secret"
```
| Rather than the location to a file, its Base64 encoded content can be provided.
If you chose this option, the value of the property should start with |

See `JksSslBundleProperties` for the full set of supported properties.

| If you’re using environment variables to configure the bundle, the name of the bundle is always converted to lowercase. |

## Configuring SSL With PEM-encoded Certificates

Configuration properties with the prefix `spring.ssl.bundle.pem` can be used to configure bundles of trust material in the form of PEM-encoded text.
Each bundle has a user-provided name that can be used to reference the bundle.

When used to secure an embedded web server, a `keystore` is typically configured with a certificate and private key as shown in this example:

-
Properties
-
YAML

```
spring.ssl.bundle.pem.mybundle.keystore.certificate=classpath:application.crt
spring.ssl.bundle.pem.mybundle.keystore.private-key=classpath:application.key
```
```
spring:
 ssl:
 bundle:
 pem:
 mybundle:
 keystore:
 certificate: "classpath:application.crt"
 private-key: "classpath:application.key"
```
When used to secure a client-side connection, a `truststore` is typically configured with the server certificate as shown in this example:

-
Properties
-
YAML

`spring.ssl.bundle.pem.mybundle.truststore.certificate=classpath:server.crt````
spring:
 ssl:
 bundle:
 pem:
 mybundle:
 truststore:
 certificate: "classpath:server.crt"
```
| Rather than the location to a file, its Base64 encoded content can be provided.
If you chose this option, the value of the property should start with PEM content can also be used directly for both the The following example shows how a truststore certificate can be defined:
 |

See `PemSslBundleProperties` for the full set of supported properties.

| If you’re using environment variables to configure the bundle, the name of the bundle is always converted to lowercase. |

## Applying SSL Bundles

Once configured using properties, SSL bundles can be referred to by name in configuration properties for various types of connections that are auto-configured by Spring Boot. See the sections on embedded web servers, data technologies, and REST clients for further information.

## Using SSL Bundles

Spring Boot auto-configures a bean of type `SslBundles` that provides access to each of the named bundles configured using the `spring.ssl.bundle` properties.

An `SslBundle` can be retrieved from the auto-configured `SslBundles` bean and used to create objects that are used to configure SSL connectivity in client libraries.
The `SslBundle` provides a layered approach of obtaining these SSL objects:

-
`getStores()`provides access to the key store and trust store`KeyStore`instances as well as any required key store password.
-
`getManagers()`provides access to the`KeyManagerFactory`and`TrustManagerFactory`instances as well as the`KeyManager`and`TrustManager`arrays that they create.
-
`createSslContext()`provides a convenient way to obtain a new`SSLContext`instance.

In addition, the `SslBundle` provides details about the key being used, the protocol to use and any option that should be applied to the SSL engine.

The following example shows retrieving an `SslBundle` and using it to create an `SSLContext`:

-
Java
-
Kotlin

```
import javax.net.ssl.SSLContext;
import org.springframework.boot.ssl.SslBundle;
import org.springframework.boot.ssl.SslBundles;
import org.springframework.stereotype.Component;
@Component
public class MyComponent {
	public MyComponent(SslBundles sslBundles) {
 SslBundle sslBundle = sslBundles.getBundle("mybundle");
 SSLContext sslContext = sslBundle.createSslContext();
 // do something with the created sslContext
	}
}
```
```
import org.springframework.boot.ssl.SslBundles
import org.springframework.stereotype.Component
@Component
class MyComponent(sslBundles: SslBundles) {
 init {
 val sslBundle = sslBundles.getBundle("mybundle")
 val sslContext = sslBundle.createSslContext()
 // do something with the created sslContext
 }
}
```
## Reloading SSL bundles

SSL bundles can be reloaded when the key material changes. The component consuming the bundle has to be compatible with reloadable SSL bundles. Currently the following components are compatible:

-
Tomcat web server
-
Netty web server

To enable reloading, you need to opt-in via a configuration property as shown in this example:

-
Properties
-
YAML

```
spring.ssl.bundle.pem.mybundle.reload-on-update=true
spring.ssl.bundle.pem.mybundle.keystore.certificate=file:/some/directory/application.crt
spring.ssl.bundle.pem.mybundle.keystore.private-key=file:/some/directory/application.key
```
```
spring:
 ssl:
 bundle:
 pem:
 mybundle:
 reload-on-update: true
 keystore:
 certificate: "file:/some/directory/application.crt"
 private-key: "file:/some/directory/application.key"
```
A file watcher is then watching the files and if they change, the SSL bundle will be reloaded. This in turn triggers a reload in the consuming component, e.g. Tomcat rotates the certificates in the SSL enabled connectors.

You can configure the quiet period (to make sure that there are no more changes) of the file watcher with the `spring.ssl.bundle.watch.file.quiet-period` property.

### Reloading SSL Bundles With Let’s Encrypt

If you use certificates issued by Let’s Encrypt and renewed by an external tool, such as Certbot, you can configure a PEM bundle to use the generated files and enable reloading.
Certbot typically stores these in `/etc/letsencrypt/live/` under a directory named after your domain.
The following example shows how to configure a PEM bundle for `example.com`:

-
Properties
-
YAML

```
spring.ssl.bundle.pem.webserver.reload-on-update=true
spring.ssl.bundle.pem.webserver.keystore.certificate=file:/etc/letsencrypt/live/example.com/fullchain.pem
spring.ssl.bundle.pem.webserver.keystore.private-key=file:/etc/letsencrypt/live/example.com/privkey.pem
server.ssl.bundle=webserver
```
```
spring:
 ssl:
 bundle:
 pem:
 webserver:
 reload-on-update: true
 keystore:
 certificate: "file:/etc/letsencrypt/live/example.com/fullchain.pem"
 private-key: "file:/etc/letsencrypt/live/example.com/privkey.pem"
server:
 ssl:
 bundle: "webserver"
```
Spring Boot does not request or renew Let’s Encrypt certificates. When Certbot or another ACME client updates the configured files, the SSL bundle is reloaded. Compatible consumers, such as Tomcat and Netty web servers, can then use the updated certificate without restarting the application.

The files in `/etc/letsencrypt/live` are typically symbolic links to files in `/etc/letsencrypt/archive`.
The file watcher follows symbolic links so that updates to the target files can trigger a reload.

# Web

Spring Boot is well suited for web application development.
You can create a self-contained HTTP server by using embedded Tomcat, Jetty, or Netty.
Most web applications use the `spring-boot-starter-web` module to get up and running quickly.
You can also choose to build reactive web applications by using the `spring-boot-starter-webflux` module.

If you have not yet developed a Spring Boot web application, you can follow the “Hello World!” example in the Getting started section.

# Servlet Web Applications

If you want to build servlet-based web applications, you can take advantage of Spring Boot’s auto-configuration for Spring MVC or Jersey.

## The “Spring Web MVC Framework”

The Spring Web MVC framework (often referred to as “Spring MVC”) is a rich “model view controller” web framework.
Spring MVC lets you create special `@Controller` or `@RestController` beans to handle incoming HTTP requests.
Methods in your controller are mapped to HTTP by using `@RequestMapping` annotations.

The following code shows a typical `@RestController` that serves JSON data:

-
Java
-
Kotlin

```
import java.util.List;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
@RestController
@RequestMapping("/users")
public class MyRestController {
	private final UserRepository userRepository;
	private final CustomerRepository customerRepository;
	public MyRestController(UserRepository userRepository, CustomerRepository customerRepository) {
 this.userRepository = userRepository;
 this.customerRepository = customerRepository;
	}
	@GetMapping("/{userId}")
	public User getUser(@PathVariable Long userId) {
 return this.userRepository.findById(userId).get();
	}
	@GetMapping("/{userId}/customers")
	public List<Customer> getUserCustomers(@PathVariable Long userId) {
 return this.userRepository.findById(userId).map(this.customerRepository::findByUser).get();
	}
	@DeleteMapping("/{userId}")
	public void deleteUser(@PathVariable Long userId) {
 this.userRepository.deleteById(userId);
	}
}
```
```
import org.springframework.web.bind.annotation.DeleteMapping
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController
@RestController
@RequestMapping("/users")
class MyRestController(private val userRepository: UserRepository, private val customerRepository: CustomerRepository) {
	@GetMapping("/{userId}")
	fun getUser(@PathVariable userId: Long): User {
 return userRepository.findById(userId).get()
	}
	@GetMapping("/{userId}/customers")
	fun getUserCustomers(@PathVariable userId: Long): List<Customer> {
 return userRepository.findById(userId).map(customerRepository::findByUser).get()
	}
	@DeleteMapping("/{userId}")
	fun deleteUser(@PathVariable userId: Long) {
 userRepository.deleteById(userId)
	}
}
```
“WebMvc.fn”, the functional variant, separates the routing configuration from the actual handling of the requests, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.MediaType;
import org.springframework.web.servlet.function.RequestPredicate;
import org.springframework.web.servlet.function.RouterFunction;
import org.springframework.web.servlet.function.ServerResponse;
import static org.springframework.web.servlet.function.RequestPredicates.accept;
import static org.springframework.web.servlet.function.RouterFunctions.route;
@Configuration(proxyBeanMethods = false)
public class MyRoutingConfiguration {
	private static final RequestPredicate ACCEPT_JSON = accept(MediaType.APPLICATION_JSON);
	@Bean
	public RouterFunction<ServerResponse> routerFunction(MyUserHandler userHandler) {
 return route()
 .GET("/{user}", ACCEPT_JSON, userHandler::getUser)
 .GET("/{user}/customers", ACCEPT_JSON, userHandler::getUserCustomers)
 .DELETE("/{user}", ACCEPT_JSON, userHandler::deleteUser)
 .build();
	}
}
```
```
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.http.MediaType
import org.springframework.web.servlet.function.RequestPredicates.accept
import org.springframework.web.servlet.function.RouterFunction
import org.springframework.web.servlet.function.RouterFunctions
import org.springframework.web.servlet.function.ServerResponse
@Configuration(proxyBeanMethods = false)
class MyRoutingConfiguration {
	@Bean
	fun routerFunction(userHandler: MyUserHandler): RouterFunction<ServerResponse> {
 return RouterFunctions.route()
 .GET("/{user}", ACCEPT_JSON, userHandler::getUser)
 .GET("/{user}/customers", ACCEPT_JSON, userHandler::getUserCustomers)
 .DELETE("/{user}", ACCEPT_JSON, userHandler::deleteUser)
 .build()
	}
	companion object {
 private val ACCEPT_JSON = accept(MediaType.APPLICATION_JSON)
	}
}
```
-
Java
-
Kotlin

```
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.function.ServerRequest;
import org.springframework.web.servlet.function.ServerResponse;
@Component
public class MyUserHandler {
	public ServerResponse getUser(ServerRequest request) {
 ...
	}
	public ServerResponse getUserCustomers(ServerRequest request) {
 ...
	}
	public ServerResponse deleteUser(ServerRequest request) {
 ...
	}
}
```
```
import org.springframework.stereotype.Component
import org.springframework.web.servlet.function.ServerRequest
import org.springframework.web.servlet.function.ServerResponse
@Component
class MyUserHandler {
	fun getUser(request: ServerRequest?): ServerResponse {
 ...
	}
	fun getUserCustomers(request: ServerRequest?): ServerResponse {
 ...
	}
	fun deleteUser(request: ServerRequest?): ServerResponse {
 ...
	}
}
```
Spring MVC is part of the core Spring Framework, and detailed information is available in the reference documentation. There are also several guides that cover Spring MVC available at spring.io/guides.

| You can define as many `RouterFunction`beans as you like to modularize the definition of the router.
Beans can be ordered if you need to apply a precedence. |

### Spring MVC Auto-configuration

Spring Boot provides auto-configuration for Spring MVC that works well with most applications.
It replaces the need for `@EnableWebMvc` and the two cannot be used together.
In addition to Spring MVC’s defaults, the auto-configuration provides the following features:

-
Inclusion of `ContentNegotiatingViewResolver`and`BeanNameViewResolver`beans.
-
Support for serving static resources, including support for WebJars (covered later in this document).
-
Automatic registration of `Converter`,`GenericConverter`, and`Formatter`beans.
-
Support for `HttpMessageConverters`(covered later in this document).
-
Automatic registration of `MessageCodesResolver`(covered later in this document).
-
Static `index.html`support.
-
Automatic use of a `ConfigurableWebBindingInitializer`bean (covered later in this document).

If you want to keep those Spring Boot MVC customizations and make more MVC customizations (interceptors, formatters, view controllers, and other features), you can add your own `@Configuration` class of type `WebMvcConfigurer` but **without** `@EnableWebMvc`.

If you want to provide custom instances of `RequestMappingHandlerMapping`, `RequestMappingHandlerAdapter`, or `ExceptionHandlerExceptionResolver`, and still keep the Spring Boot MVC customizations, you can declare a bean of type `WebMvcRegistrations` and use it to provide custom instances of those components.
The custom instances will be subject to further initialization and configuration by Spring MVC.
To participate in, and if desired, override that subsequent processing, a `WebMvcConfigurer` should be used.

If you do not want to use the auto-configuration and want to take complete control of Spring MVC, add your own `@Configuration` annotated with `@EnableWebMvc`.
Alternatively, add your own `@Configuration`-annotated `DelegatingWebMvcConfiguration` as described in the `@EnableWebMvc` API documentation.

### Spring MVC Conversion Service

Spring MVC uses a different `ConversionService` to the one used to convert values from your `application.properties` or `application.yaml` file.
It means that `Period`, `Duration` and `DataSize` converters are not available and that `@DurationUnit` and `@DataSizeUnit` annotations will be ignored.

If you want to customize the `ConversionService` used by Spring MVC, you can provide a `WebMvcConfigurer` bean with an `addFormatters` method.
From this method you can register any converter that you like, or you can delegate to the static methods available on `ApplicationConversionService`.

Conversion can also be customized using the `spring.mvc.format.*` configuration properties.
When not configured, the following defaults are used:

| Property | `DateTimeFormatter` | Formats |
|---|---|---|
|
 |
 |
 |
|
 |
 | java.time’s |
|
 |
 | java.time’s |

### HttpMessageConverters

Spring MVC uses the `HttpMessageConverter` interface to convert HTTP requests and responses.
Sensible defaults are included out of the box.
For example, objects can be automatically converted to JSON (by using the Jackson library) or XML (by using the Jackson XML extension, if available, or by using JAXB if the Jackson XML extension is not available).
By default, strings are encoded in `UTF-8`.

Any `HttpMessageConverter` bean that is present in the context is added to the list of converters.
You can also override default converters in the same way.

If you need to add or customize converters, you can declare one or more `ClientHttpMessageConvertersCustomizer` or
`ServerHttpMessageConvertersCustomizer` as beans. There, you can choose whether converter instances should be added
before default ones (`addCustomConverter`) or if they should override a specific default converter (like `withJsonConverter`).

See the following listing for an example:

-
Java
-
Kotlin

```
import java.text.SimpleDateFormat;
import tools.jackson.databind.json.JsonMapper;
import org.springframework.boot.http.converter.autoconfigure.ClientHttpMessageConvertersCustomizer;
import org.springframework.boot.http.converter.autoconfigure.ServerHttpMessageConvertersCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.converter.HttpMessageConverters.ClientBuilder;
import org.springframework.http.converter.HttpMessageConverters.ServerBuilder;
import org.springframework.http.converter.json.JacksonJsonHttpMessageConverter;
@Configuration(proxyBeanMethods = false)
public class MyHttpMessageConvertersConfiguration {
	@Bean
	public ClientHttpMessageConvertersCustomizer myClientConvertersCustomizer() {
 return (clientBuilder) -> clientBuilder.addCustomConverter(new AdditionalHttpMessageConverter())
 .addCustomConverter(new AnotherHttpMessageConverter());
	}
	@Bean
	public JacksonConverterCustomizer jacksonConverterCustomizer() {
 JsonMapper jsonMapper = JsonMapper.builder().defaultDateFormat(new SimpleDateFormat("yyyy-MM")).build();
 return new JacksonConverterCustomizer(jsonMapper);
	}
	// contribute a custom JSON converter to both client and server
	static class JacksonConverterCustomizer
 implements ClientHttpMessageConvertersCustomizer, ServerHttpMessageConvertersCustomizer {
 private final JsonMapper jsonMapper;
 JacksonConverterCustomizer(JsonMapper jsonMapper) {
 this.jsonMapper = jsonMapper;
 }
 @Override
 public void customize(ClientBuilder builder) {
 builder.withJsonConverter(new JacksonJsonHttpMessageConverter(this.jsonMapper));
 }
 @Override
 public void customize(ServerBuilder builder) {
 builder.withJsonConverter(new JacksonJsonHttpMessageConverter(this.jsonMapper));
 }
	}
}
```
```
import org.springframework.boot.http.converter.autoconfigure.ClientHttpMessageConvertersCustomizer
import org.springframework.boot.http.converter.autoconfigure.ServerHttpMessageConvertersCustomizer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.http.converter.HttpMessageConverters
import org.springframework.http.converter.json.JacksonJsonHttpMessageConverter
import tools.jackson.databind.json.JsonMapper
import java.text.SimpleDateFormat
@Configuration(proxyBeanMethods = false)
class MyHttpMessageConvertersConfiguration {
	@Bean
	fun myClientConvertersCustomizer(): ClientHttpMessageConvertersCustomizer {
 return ClientHttpMessageConvertersCustomizer { clientBuilder: HttpMessageConverters.ClientBuilder ->
 clientBuilder
 .addCustomConverter(AdditionalHttpMessageConverter())
 .addCustomConverter(AnotherHttpMessageConverter())
 }
	}
	@Bean
	fun jacksonConverterCustomizer(): JacksonConverterCustomizer {
 val jsonMapper = JsonMapper.builder()
 .defaultDateFormat(SimpleDateFormat("yyyy-MM"))
 .build()
 return JacksonConverterCustomizer(jsonMapper)
	}
	// contribute a custom JSON converter to both client and server
	class JacksonConverterCustomizer(private val jsonMapper: JsonMapper) :
 ClientHttpMessageConvertersCustomizer, ServerHttpMessageConvertersCustomizer {
 override fun customize(builder: HttpMessageConverters.ClientBuilder) {
 builder.withJsonConverter(JacksonJsonHttpMessageConverter(this.jsonMapper))
 }
 override fun customize(builder: HttpMessageConverters.ServerBuilder) {
 builder.withJsonConverter(JacksonJsonHttpMessageConverter(this.jsonMapper))
 }
	}
}
```
### MessageCodesResolver

Spring MVC has a strategy for generating error codes for rendering error messages from binding errors: `MessageCodesResolver`.
If you set the `spring.mvc.message-codes-resolver-format` property `PREFIX_ERROR_CODE` or `POSTFIX_ERROR_CODE`, Spring Boot creates one for you (see the enumeration in `DefaultMessageCodesResolver.Format`).

### Static Content

By default, Spring Boot serves static content from a directory called `/static` (or `/public` or `/resources` or `/META-INF/resources`) in the classpath or from the root of the `ServletContext`.
It uses the `ResourceHttpRequestHandler` from Spring MVC so that you can modify that behavior by adding your own `WebMvcConfigurer` and overriding the `addResourceHandlers` method.

In a stand-alone web application, the default servlet from the container is not enabled.
It can be enabled using the `server.servlet.register-default-servlet` property.

The default servlet acts as a fallback, serving content from the root of the `ServletContext` if Spring decides not to handle it.
Most of the time, this does not happen (unless you modify the default MVC configuration), because Spring can always handle requests through the `DispatcherServlet`.

By default, resources are mapped on `/**`, but you can tune that with the `spring.mvc.static-path-pattern` property.
For instance, relocating all resources to `/resources/**` can be achieved as follows:

-
Properties
-
YAML

`spring.mvc.static-path-pattern=/resources/**````
spring:
 mvc:
 static-path-pattern: "/resources/**"
```
You can also customize the static resource locations by using the `spring.web.resources.static-locations` property (replacing the default values with a list of directory locations).
The root servlet context path, `"/"`, is automatically added as a location as well.

In addition to the “standard” static resource locations mentioned earlier, a special case is made for Webjars content.
By default, any resources with a path in `/webjars/**` are served from jar files if they are packaged in the Webjars format.
The path can be customized with the `spring.mvc.webjars-path-pattern` property.

| Do not use the `src/main/webapp`directory if your application is packaged as a jar.
Although this directory is a common standard, it worksonlywith war packaging, and it is silently ignored by most build tools if you generate a jar. |

Spring Boot also supports the advanced resource handling features provided by Spring MVC, allowing use cases such as cache-busting static resources or using version agnostic URLs for Webjars.

To use version agnostic URLs for Webjars, add the `org.webjars:webjars-locator-lite` dependency.
Then declare your Webjar.
Using jQuery as an example, adding `"/webjars/jquery/jquery.min.js"` results in `"/webjars/jquery/x.y.z/jquery.min.js"` where `x.y.z` is the Webjar version.

To use cache busting, the following configuration configures a cache busting solution for all static resources, effectively adding a content hash, such as `<link href="/css/spring-2a2d595e6ed9a0b24f027f2b63b134d6.css"/>`, in URLs:

-
Properties
-
YAML

```
spring.web.resources.chain.strategy.content.enabled=true
spring.web.resources.chain.strategy.content.paths=/**
```
```
spring:
 web:
 resources:
 chain:
 strategy:
 content:
 enabled: true
 paths: "/**"
```
| Links to resources are rewritten in templates at runtime, thanks to a `ResourceUrlEncodingFilter`that is auto-configured for Thymeleaf and FreeMarker.
You should manually declare this filter when using JSPs.
Other template engines are currently not automatically supported but can be with custom template macros/helpers and the use of the`ResourceUrlProvider`. |

When loading resources dynamically with, for example, a JavaScript module loader, renaming files is not an option. That is why other strategies are also supported and can be combined. A "fixed" strategy adds a static version string in the URL without changing the file name, as shown in the following example:

-
Properties
-
YAML

```
spring.web.resources.chain.strategy.content.enabled=true
spring.web.resources.chain.strategy.content.paths=/**
spring.web.resources.chain.strategy.fixed.enabled=true
spring.web.resources.chain.strategy.fixed.paths=/js/lib/
spring.web.resources.chain.strategy.fixed.version=v12
```
```
spring:
 web:
 resources:
 chain:
 strategy:
 content:
 enabled: true
 paths: "/**"
 fixed:
 enabled: true
 paths: "/js/lib/"
 version: "v12"
```
With this configuration, JavaScript modules located under `"/js/lib/"` use a fixed versioning strategy (`"/v12/js/lib/mymodule.js"`), while other resources still use the content one (`<link href="/css/spring-2a2d595e6ed9a0b24f027f2b63b134d6.css"/>`).

See `WebProperties.Resources` for more supported options.

| This feature has been thoroughly described in a dedicated blog post and in Spring Framework’s reference documentation. |

### Welcome Page

Spring Boot supports both static and templated welcome pages.
It first looks for an `index.html` file in the configured static content locations.
If one is not found, it then looks for an `index` template.
If either is found, it is automatically used as the welcome page of the application.

This only acts as a fallback for actual index routes defined by the application.
The ordering is defined by the order of `HandlerMapping` beans which is by default the following:

|
 | Endpoints declared with |
|
 | Endpoints declared in |
|
 | The welcome page support |

### Custom Favicon

As with other static resources, Spring Boot checks for a `favicon.ico` in the configured static content locations.
If such a file is present, it is automatically used as the favicon of the application.

### Path Matching and Content Negotiation

Spring MVC can map incoming HTTP requests to handlers by looking at the request path and matching it to the mappings defined in your application (for example, `@GetMapping` annotations on Controller methods).

Spring Boot chooses to disable suffix pattern matching by default, which means that requests like `"GET /projects/spring-boot.json"` will not be matched to `@GetMapping("/projects/spring-boot")` mappings.
This is considered as a best practice for Spring MVC applications.
This feature was mainly useful in the past for HTTP clients which did not send proper "Accept" request headers; we needed to make sure to send the correct Content Type to the client.
Nowadays, Content Negotiation is much more reliable.

There are other ways to deal with HTTP clients that do not consistently send proper "Accept" request headers.
Instead of using suffix matching, we can use a query parameter to ensure that requests like `"GET /projects/spring-boot?format=json"` will be mapped to `@GetMapping("/projects/spring-boot")`:

-
Properties
-
YAML

`spring.mvc.contentnegotiation.favor-parameter=true````
spring:
 mvc:
 contentnegotiation:
 favor-parameter: true
```
Or if you prefer to use a different parameter name:

-
Properties
-
YAML

```
spring.mvc.contentnegotiation.favor-parameter=true
spring.mvc.contentnegotiation.parameter-name=myparam
```
```
spring:
 mvc:
 contentnegotiation:
 favor-parameter: true
 parameter-name: "myparam"
```
Most standard media types are supported out-of-the-box, but you can also define new ones:

-
Properties
-
YAML

`spring.mvc.contentnegotiation.media-types.markdown=text/markdown````
spring:
 mvc:
 contentnegotiation:
 media-types:
 markdown: "text/markdown"
```
As of Spring Framework 5.3, Spring MVC supports two strategies for matching request paths to controllers.
By default, Spring Boot uses the `PathPatternParser` strategy.
`PathPatternParser` is an optimized implementation but comes with some restrictions compared to the `AntPathMatcher` strategy.
`PathPatternParser` restricts usage of some path pattern variants.
It is also incompatible with configuring the `DispatcherServlet` with a path prefix (`spring.mvc.servlet.path`).

The strategy can be configured using the `spring.mvc.pathmatch.matching-strategy` configuration property, as shown in the following example:

-
Properties
-
YAML

`spring.mvc.pathmatch.matching-strategy=ant-path-matcher````
spring:
 mvc:
 pathmatch:
 matching-strategy: "ant-path-matcher"
```
Spring MVC will throw a `NoHandlerFoundException` if a handler is not found for a request.
Note that, by default, the serving of static content is mapped to `/**` and will, therefore, provide a handler for all requests.
If no static content is available, `ResourceHttpRequestHandler` will throw a `NoResourceFoundException`.
For a `NoHandlerFoundException` to be thrown, set `spring.mvc.static-path-pattern` to a more specific value such as `/resources/**` or set `spring.web.resources.add-mappings` to `false` to disable serving of static content entirely.

### ConfigurableWebBindingInitializer

Spring MVC uses a `WebBindingInitializer` to initialize a `WebDataBinder` for a particular request.
If you create your own `ConfigurableWebBindingInitializer` `@Bean`, Spring Boot automatically configures Spring MVC to use it.

### Template Engines

As well as REST web services, you can also use Spring MVC to serve dynamic HTML content. Spring MVC supports a variety of templating technologies, including Thymeleaf, FreeMarker, and JSPs. Also, many other templating engines include their own Spring MVC integrations.

Spring Boot includes auto-configuration support for the following templating engines:

| If possible, JSPs should be avoided. There are several known limitations when using them with embedded servlet containers. |

When you use one of these templating engines with the default configuration, your templates are picked up automatically from `src/main/resources/templates`.

| Depending on how you run your application, your IDE may order the classpath differently. Running your application in the IDE from its main method results in a different ordering than when you run your application by using Maven or Gradle or from its packaged jar. This can cause Spring Boot to fail to find the expected template. If you have this problem, you can reorder the classpath in the IDE to place the module’s classes and resources first. |

### Error Handling

By default, Spring Boot provides an `/error` mapping that handles all errors in a sensible way, and it is registered as a “global” error page in the servlet container.
For machine clients, it produces a JSON response with details of the error, the HTTP status, and the exception message.
For browser clients, there is a “whitelabel” error view that renders the same data in HTML format (to customize it, add a `View` that resolves to `error`).

There are a number of `spring.web.error` properties that can be set if you want to customize the default error handling behavior.
See the Web Properties section of the Appendix.

To replace the default behavior completely, you can implement `ErrorController` and register a bean definition of that type or add a bean of type `ErrorAttributes` to use the existing mechanism but replace the contents.

| The `BasicErrorController`can be used as a base class for a custom`ErrorController`.
This is particularly useful if you want to add a handler for a new content type (the default is to handle`text/html`specifically and provide a fallback for everything else).
To do so, extend`BasicErrorController`, add a public method with a`@RequestMapping`that has a`produces`attribute, and create a bean of your new type. |

As of Spring Framework 6.0, RFC 9457 Problem Details is supported.
Spring MVC can produce custom error messages with the `application/problem+json` media type, like:

```
{
	"type": "https://example.org/problems/unknown-project",
	"title": "Unknown project",
	"status": 404,
	"detail": "No project found for id 'spring-unknown'",
	"instance": "/projects/spring-unknown"
}
```
This support can be enabled by setting `spring.mvc.problemdetails.enabled` to `true`.

You can also define a class annotated with `@ControllerAdvice` to customize the JSON document to return for a particular controller and/or exception type, as shown in the following example:

-
Java
-
Kotlin

```
import jakarta.servlet.RequestDispatcher;
import jakarta.servlet.http.HttpServletRequest;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.servlet.mvc.method.annotation.ResponseEntityExceptionHandler;
@ControllerAdvice(basePackageClasses = SomeController.class)
public class MyControllerAdvice extends ResponseEntityExceptionHandler {
	@ResponseBody
	@ExceptionHandler(MyException.class)
	public ResponseEntity<?> handleControllerException(HttpServletRequest request, Throwable ex) {
 HttpStatus status = getStatus(request);
 return new ResponseEntity<>(new MyErrorBody(status.value(), ex.getMessage()), status);
	}
	private HttpStatus getStatus(HttpServletRequest request) {
 Integer code = (Integer) request.getAttribute(RequestDispatcher.ERROR_STATUS_CODE);
 HttpStatus status = HttpStatus.resolve(code);
 return (status != null) ? status : HttpStatus.INTERNAL_SERVER_ERROR;
	}
}
```
```
import jakarta.servlet.RequestDispatcher
import jakarta.servlet.http.HttpServletRequest
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.ControllerAdvice
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.ResponseBody
import org.springframework.web.servlet.mvc.method.annotation.ResponseEntityExceptionHandler
@ControllerAdvice(basePackageClasses = [SomeController::class])
class MyControllerAdvice : ResponseEntityExceptionHandler() {
	@ResponseBody
	@ExceptionHandler(MyException::class)
	fun handleControllerException(request: HttpServletRequest, ex: Throwable): ResponseEntity<*> {
 val status = getStatus(request)
 return ResponseEntity(MyErrorBody(status.value(), ex.message), status)
	}
	private fun getStatus(request: HttpServletRequest): HttpStatus {
 val code = request.getAttribute(RequestDispatcher.ERROR_STATUS_CODE) as Int
 val status = HttpStatus.resolve(code)
 return status ?: HttpStatus.INTERNAL_SERVER_ERROR
	}
}
```
In the preceding example, if `MyException` is thrown by a controller defined in the same package as `SomeController`, a JSON representation of the `MyErrorBody` POJO is used instead of the `ErrorAttributes` representation.

In some cases, errors handled at the controller level are not recorded by web observations or the metrics infrastructure. Applications can ensure that such exceptions are recorded with the observations by setting the handled exception on the observation context.

#### Custom Error Pages

If you want to display a custom HTML error page for a given status code, you can add a file to an `/error` directory.
Error pages can either be static HTML (that is, added under any of the static resource directories) or be built by using templates.
The name of the file should be the exact status code or a series mask.

For example, to map `404` to a static HTML file, your directory structure would be as follows:

```
src/
 +- main/
 +- java/
 | + <source code>
 +- resources/
 +- public/
 +- error/
 | +- 404.html
 +- <other public assets>
```
To map all `5xx` errors by using a FreeMarker template, your directory structure would be as follows:

```
src/
 +- main/
 +- java/
 | + <source code>
 +- resources/
 +- templates/
 +- error/
 | +- 5xx.ftlh
 +- <other templates>
```
For more complex mappings, you can also add beans that implement the `ErrorViewResolver` interface, as shown in the following example:

-
Java
-
Kotlin

```
import java.util.Map;
import jakarta.servlet.http.HttpServletRequest;
import org.springframework.boot.webmvc.autoconfigure.error.ErrorViewResolver;
import org.springframework.http.HttpStatus;
import org.springframework.web.servlet.ModelAndView;
public class MyErrorViewResolver implements ErrorViewResolver {
	@Override
	public ModelAndView resolveErrorView(HttpServletRequest request, HttpStatus status, Map<String, Object> model) {
 // Use the request or status to optionally return a ModelAndView
 if (status == HttpStatus.INSUFFICIENT_STORAGE) {
 // We could add custom model values here
 return new ModelAndView("myview");
 }
 return null;
	}
}
```
```
import jakarta.servlet.http.HttpServletRequest
import org.springframework.boot.webmvc.autoconfigure.error.ErrorViewResolver
import org.springframework.http.HttpStatus
import org.springframework.web.servlet.ModelAndView
class MyErrorViewResolver : ErrorViewResolver {
	override fun resolveErrorView(request: HttpServletRequest, status: HttpStatus,
 model: Map<String, Any>): ModelAndView? {
 // Use the request or status to optionally return a ModelAndView
 if (status == HttpStatus.INSUFFICIENT_STORAGE) {
 // We could add custom model values here
 return ModelAndView("myview")
 }
 return null
	}
}
```
You can also use regular Spring MVC features such as `@ExceptionHandler` methods and `@ControllerAdvice`.
The `ErrorController` then picks up any unhandled exceptions.

#### Mapping Error Pages Outside of Spring MVC

For applications that do not use Spring MVC, you can use the `ErrorPageRegistrar` interface to directly register `ErrorPage` instances.
This abstraction works directly with the underlying embedded servlet container and works even if you do not have a Spring MVC `DispatcherServlet`.

-
Java
-
Kotlin

```
import org.springframework.boot.web.error.ErrorPage;
import org.springframework.boot.web.error.ErrorPageRegistrar;
import org.springframework.boot.web.error.ErrorPageRegistry;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpStatus;
@Configuration(proxyBeanMethods = false)
public class MyErrorPagesConfiguration {
	@Bean
	public ErrorPageRegistrar errorPageRegistrar() {
 return this::registerErrorPages;
	}
	private void registerErrorPages(ErrorPageRegistry registry) {
 registry.addErrorPages(new ErrorPage(HttpStatus.BAD_REQUEST, "/400"));
	}
}
```
```
import org.springframework.boot.web.error.ErrorPage
import org.springframework.boot.web.error.ErrorPageRegistrar
import org.springframework.boot.web.error.ErrorPageRegistry
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.http.HttpStatus
@Configuration(proxyBeanMethods = false)
class MyErrorPagesConfiguration {
	@Bean
	fun errorPageRegistrar(): ErrorPageRegistrar {
 return ErrorPageRegistrar { registry: ErrorPageRegistry -> registerErrorPages(registry) }
	}
	private fun registerErrorPages(registry: ErrorPageRegistry) {
 registry.addErrorPages(ErrorPage(HttpStatus.BAD_REQUEST, "/400"))
	}
}
```
| If you register an `ErrorPage`with a path that ends up being handled by a`Filter`(as is common with some non-Spring web frameworks, like Jersey and Wicket), then the`Filter`has to be explicitly registered as an`ERROR`dispatcher, as shown in the following example: |

-
Java
-
Kotlin

```
import java.util.EnumSet;
import jakarta.servlet.DispatcherType;
import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyFilterConfiguration {
	@Bean
	public FilterRegistrationBean<MyFilter> myFilter() {
 FilterRegistrationBean<MyFilter> registration = new FilterRegistrationBean<>(new MyFilter());
 // ...
 registration.setDispatcherTypes(EnumSet.allOf(DispatcherType.class));
 return registration;
	}
}
```
```
import jakarta.servlet.DispatcherType
import org.springframework.boot.web.servlet.FilterRegistrationBean
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import java.util.EnumSet
@Configuration(proxyBeanMethods = false)
class MyFilterConfiguration {
	@Bean
	fun myFilter(): FilterRegistrationBean<MyFilter> {
 val registration = FilterRegistrationBean(MyFilter())
 // ...
 registration.setDispatcherTypes(EnumSet.allOf(DispatcherType::class.java))
 return registration
	}
}
```
Note that the default `FilterRegistrationBean` does not include the `ERROR` dispatcher type.

#### Error Handling in a WAR Deployment

When deployed to a servlet container, Spring Boot uses its error page filter to forward a request with an error status to the appropriate error page. This is necessary as the servlet specification does not provide an API for registering error pages. Depending on the container that you are deploying your war file to and the technologies that your application uses, some additional configuration may be required.

The error page filter can only forward the request to the correct error page if the response has not already been committed.
By default, WebSphere Application Server 8.0 and later commits the response upon successful completion of a servlet’s service method.
You should disable this behavior by setting `com.ibm.ws.webcontainer.invokeFlushAfterService` to `false`.

### CORS Support

Cross-origin resource sharing (CORS) is a W3C specification implemented by most browsers that lets you specify in a flexible way what kind of cross-domain requests are authorized, instead of using some less secure and less powerful approaches such as IFRAME or JSONP.

As of version 4.2, Spring MVC supports CORS.
Using controller method CORS configuration with `@CrossOrigin` annotations in your Spring Boot application does not require any specific configuration.
Global CORS configuration can be defined by registering a `WebMvcConfigurer` bean with a customized `addCorsMappings(CorsRegistry)` method, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.CorsRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
@Configuration(proxyBeanMethods = false)
public class MyCorsConfiguration {
	@Bean
	public WebMvcConfigurer corsConfigurer() {
 return new WebMvcConfigurer() {
 @Override
 public void addCorsMappings(CorsRegistry registry) {
 registry.addMapping("/api/**");
 }
 };
	}
}
```
```
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.web.servlet.config.annotation.CorsRegistry
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer
@Configuration(proxyBeanMethods = false)
class MyCorsConfiguration {
	@Bean
	fun corsConfigurer(): WebMvcConfigurer {
 return object : WebMvcConfigurer {
 override fun addCorsMappings(registry: CorsRegistry) {
 registry.addMapping("/api/**")
 }
 }
	}
}
```
### API Versioning

Spring MVC supports API versioning which can be used to evolve an HTTP API over time.
The same `@Controller` path can be mapped multiple times to support different versions of the API.

For more details see Spring Framework’s reference documentation.

Once mappings have been added, you additionally need to configure Spring MVC so that it is able to use any version information sent with a request. Typically, versions are sent as HTTP headers, query parameters, media type parameters, or as part of the path.

To configure Spring MVC, you can either use a `WebMvcConfigurer` bean and override the `configureApiVersioning(…)` method, or you can use properties.

For example, the following will use an `X-Version` HTTP header to obtain version information and default to `1.0.0` when no header is sent.

-
Properties
-
YAML

```
spring.mvc.apiversion.default=1.0.0
spring.mvc.apiversion.use.header=X-Version
```
```
spring:
 mvc:
 apiversion:
 default: 1.0.0
 use:
 header: X-Version
```
| If your setup requires multiple strategies, such as header and query parameter, consider declaring the order programmatically by overriding the `configureApiVersioning`method. |

For more complete control, you can also define `ApiVersionResolver`, `ApiVersionParser` and `ApiVersionDeprecationHandler` beans which will be injected into the auto-configured Spring MVC configuration.

| API versioning is also supported with both `WebClient`and`RestClient`.
See API Versioning for details. |

## JAX-RS and Jersey

If you prefer the JAX-RS programming model for REST endpoints, you can use one of the available implementations instead of Spring MVC.
Jersey and Apache CXF work quite well out of the box.
CXF requires you to register its `Servlet` or `Filter` as a `@Bean` in your application context.
Jersey has some native Spring support, so we also provide auto-configuration support for it in Spring Boot, together with a starter.

To get started with Jersey, include the `spring-boot-starter-jersey` as a dependency and then you need one `@Bean` of type `ResourceConfig` in which you register all the endpoints, as shown in the following example:

-
Java
-
Kotlin

```
import org.glassfish.jersey.server.ResourceConfig;
import org.springframework.stereotype.Component;
@Component
public class MyJerseyConfig extends ResourceConfig {
	public MyJerseyConfig() {
 register(MyEndpoint.class);
	}
}
```
```
import org.glassfish.jersey.server.ResourceConfig
import org.springframework.stereotype.Component
@Component
class MyJerseyConfig : ResourceConfig() {
	init {
 register(MyEndpoint::class.java)
	}
}
```
| Jersey’s support for scanning executable archives is rather limited.
For example, it cannot scan for endpoints in a package found in a fully executable jar file or in `WEB-INF/classes`when running an executable war file.
To avoid this limitation, the`packages`method should not be used, and endpoints should be registered individually by using the`register`method, as shown in the preceding example. |

For more advanced customizations, you can also register an arbitrary number of beans that implement `ResourceConfigCustomizer`.

All the registered endpoints should be a `@Component` with HTTP resource annotations (`@GET` and others), as shown in the following example:

-
Java
-
Kotlin

```
import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import org.springframework.stereotype.Component;
@Component
@Path("/hello")
public class MyEndpoint {
	@GET
	public String message() {
 return "Hello";
	}
}
```
```
import jakarta.ws.rs.GET
import jakarta.ws.rs.Path
import org.springframework.stereotype.Component
@Component
@Path("/hello")
class MyEndpoint {
	@GET
	fun message(): String {
 return "Hello"
	}
}
```
Since the `@Endpoint` is a Spring `@Component`, its lifecycle is managed by Spring and you can use the `@Autowired` annotation to inject dependencies and use the `@Value` annotation to inject external configuration.
By default, the Jersey servlet is registered and mapped to `/*`.
You can change the mapping by adding `@ApplicationPath` to your `ResourceConfig`.

By default, Jersey is set up as a servlet in a `@Bean` of type `ServletRegistrationBean` named `jerseyServletRegistration`.
By default, the servlet is initialized lazily, but you can customize that behavior by setting `spring.jersey.servlet.load-on-startup`.
You can disable or override that bean by creating one of your own with the same name.
You can also use a filter instead of a servlet by setting `spring.jersey.type=filter` (in which case, the `@Bean` to replace or override is `jerseyFilterRegistration`).
The filter has an `@Order`, which you can set with `spring.jersey.filter.order`.
When using Jersey as a filter, a servlet that will handle any requests that are not intercepted by Jersey must be present.
If your application does not contain such a servlet, you may want to enable the default servlet by setting `server.servlet.register-default-servlet` to `true`.
Both the servlet and the filter registrations can be given init parameters by using `spring.jersey.init.*` to specify a map of properties.

## Embedded Servlet Container Support

For servlet application, Spring Boot includes support for embedded Tomcat and Jetty servers.
Most developers use the appropriate starter to obtain a fully configured instance.
By default, the embedded server listens for HTTP requests on port `8080`.

### Servlets, Filters, and Listeners

When using an embedded servlet container, you can register servlets, filters, and all the listeners (such as `HttpSessionListener`) from the servlet spec, either by using Spring beans or by scanning for servlet components.

#### Registering Servlets, Filters, and Listeners as Spring Beans

Any `Servlet`, `Filter`, or servlet `*Listener` instance that is a Spring bean is registered with the embedded container.
This can be particularly convenient if you want to refer to a value from your `application.properties` during configuration.

By default, if the context contains only a single Servlet, it is mapped to `/`.
In the case of multiple servlet beans, the bean name is used as a path prefix.
Filters map to `/*`.

If convention-based mapping is not flexible enough, you can use the `ServletRegistrationBean`, `FilterRegistrationBean`, and `ServletListenerRegistrationBean` classes for complete control.
If you prefer annotations over `ServletRegistrationBean` and `FilterRegistrationBean`, you can also use `@ServletRegistration` and
`@FilterRegistration` as an alternative.

It is usually safe to leave filter beans unordered.
If a specific order is required, you should annotate the `Filter` with `@Order` or make it implement `Ordered`.
You cannot configure the order of a `Filter` by annotating its bean method with `@Order`.
If you cannot change the `Filter` class to add `@Order` or implement `Ordered`, you must define a `FilterRegistrationBean` for the `Filter` and set the registration bean’s order using the `setOrder(int)` method.
Or, if you prefer annotations, you can also use `@FilterRegistration` and set the `order` attribute.
Avoid configuring a filter that reads the request body at `Ordered.HIGHEST_PRECEDENCE`, since it might go against the character encoding configuration of your application.
If a servlet filter wraps the request, it should be configured with an order that is less than or equal to `OrderedFilter.REQUEST_WRAPPER_FILTER_MAX_ORDER`.

| To see the order of every `Filter`in your application, enable debug level logging for the`web`logging group (`logging.level.web=debug`).
Details of the registered filters, including their order and URL patterns, will then be logged at startup. |

| Take care when registering `Filter`beans since they are initialized very early in the application lifecycle.
If you need to register a`Filter`that interacts with other beans, consider using a`DelegatingFilterProxyRegistrationBean`instead. |

### Servlet Context Initialization

Embedded servlet containers do not directly execute the `ServletContainerInitializer` interface or Spring’s `WebApplicationInitializer` interface.
This is an intentional design decision intended to reduce the risk that third party libraries designed to run inside a war may break Spring Boot applications.

If you need to perform servlet context initialization in a Spring Boot application, you should register a bean that implements the `ServletContextInitializer` interface.
The single `onStartup` method provides access to the `ServletContext` and, if necessary, can easily be used as an adapter to an existing `WebApplicationInitializer`.

#### Init Parameters

Init parameters can be configured on the `ServletContext` using `server.servlet.context-parameters.*` properties.
For example, the property `server.servlet.context-parameters.com.example.parameter=example` will configure a `ServletContext` init parameter named `com.example.parameter` with the value `example`.

#### Scanning for Servlets, Filters, and listeners

When using an embedded container, automatic registration of classes annotated with `@WebServlet`, `@WebFilter`, and `@WebListener` can be enabled by using `@ServletComponentScan`.

| `@ServletComponentScan`has no effect in a standalone container, where the container’s built-in discovery mechanisms are used instead. |

### The ServletWebServerApplicationContext

Under the hood, Spring Boot uses a different type of `ApplicationContext` for embedded servlet container support.
The `ServletWebServerApplicationContext` is a special type of `WebApplicationContext` that bootstraps itself by searching for a single `ServletWebServerFactory` bean.
Usually a `TomcatServletWebServerFactory`, or `JettyServletWebServerFactory` has been auto-configured.

| You usually do not need to be aware of these implementation classes.
Most applications are auto-configured, and the appropriate `ApplicationContext`and`ServletWebServerFactory`are created on your behalf. |

In an embedded container setup, the `ServletContext` is set as part of server startup which happens during application context initialization.
Because of this beans in the `ApplicationContext` cannot be reliably initialized with a `ServletContext`.
One way to get around this is to inject `ApplicationContext` as a dependency of the bean and access the `ServletContext` only when it is needed.
Another way is to use a callback once the server has started.
This can be done using an `ApplicationListener` which listens for the `ApplicationStartedEvent` as follows:

-
Java
-
Kotlin

```
import jakarta.servlet.ServletContext;
import org.springframework.boot.context.event.ApplicationStartedEvent;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ApplicationListener;
import org.springframework.web.context.WebApplicationContext;
public class MyDemoBean implements ApplicationListener<ApplicationStartedEvent> {
	private ServletContext servletContext;
	@Override
	public void onApplicationEvent(ApplicationStartedEvent event) {
 ApplicationContext applicationContext = event.getApplicationContext();
 this.servletContext = ((WebApplicationContext) applicationContext).getServletContext();
	}
}
```
```
import jakarta.servlet.ServletContext
import org.springframework.boot.context.event.ApplicationStartedEvent
import org.springframework.context.ApplicationContext
import org.springframework.context.ApplicationListener
import org.springframework.web.context.WebApplicationContext
class MyDemoBean : ApplicationListener<ApplicationStartedEvent> {
	private var servletContext: ServletContext? = null
	override fun onApplicationEvent(event: ApplicationStartedEvent) {
 val applicationContext: ApplicationContext = event.applicationContext
 this.servletContext = (applicationContext as WebApplicationContext).servletContext
	}
}
```
### Customizing Embedded Servlet Containers

Common servlet container settings can be configured by using Spring `Environment` properties.
Usually, you would define the properties in your `application.properties` or `application.yaml` file.

Common server settings include:

-
Network settings: Listen port for incoming HTTP requests ( `server.port`), interface address to bind to (`server.address`), and so on.
-
Session settings: Whether the session is persistent ( `server.servlet.session.persistent`), session timeout (`server.servlet.session.timeout`), location of session data (`server.servlet.session.store-dir`), and session-cookie configuration (`server.servlet.session.cookie.*`).
-
Error management: Location of the error page ( `spring.web.error.path`) and so on.

Spring Boot tries as much as possible to expose common settings, but this is not always possible.
For those cases, dedicated namespaces offer server-specific customizations (see `server.tomcat`).
For instance, access logs can be configured with specific features of the embedded servlet container.

| See the `ServerProperties`class for a complete list. |

#### SameSite Cookies

The `SameSite` cookie attribute can be used by web browsers to control if and how cookies are submitted in cross-site requests.
The attribute is particularly relevant for modern web browsers which have started to change the default value that is used when the attribute is missing.

If you want to change the `SameSite` attribute of your session cookie, you can use the `server.servlet.session.cookie.same-site` property.
This property is supported by auto-configured Tomcat and Jetty servers.
It is also used to configure Spring Session servlet based `SessionRepository` beans.

For example, if you want your session cookie to have a `SameSite` attribute of `None`, you can add the following to your `application.properties` or `application.yaml` file:

-
Properties
-
YAML

`server.servlet.session.cookie.same-site=none````
server:
 servlet:
 session:
 cookie:
 same-site: "none"
```
If you want to change the `SameSite` attribute on other cookies added to your `HttpServletResponse`, you can use a `CookieSameSiteSupplier`.
The `CookieSameSiteSupplier` is passed a `Cookie` and may return a `SameSite` value, or `null`.

There are a number of convenience factory and filter methods that you can use to quickly match specific cookies.
For example, adding the following bean will automatically apply a `SameSite` of `Lax` for all cookies with a name that matches the regular expression `myapp.*`.

-
Java
-
Kotlin

```
import org.springframework.boot.web.server.servlet.CookieSameSiteSupplier;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MySameSiteConfiguration {
	@Bean
	public CookieSameSiteSupplier applicationCookieSameSiteSupplier() {
 return CookieSameSiteSupplier.ofLax().whenHasNameMatching("myapp.*");
	}
}
```
```
import org.springframework.boot.web.server.servlet.CookieSameSiteSupplier
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
@Configuration(proxyBeanMethods = false)
class MySameSiteConfiguration {
	@Bean
	fun applicationCookieSameSiteSupplier(): CookieSameSiteSupplier {
 return CookieSameSiteSupplier.ofLax().whenHasNameMatching("myapp.*")
	}
}
```
#### Character Encoding

The character encoding behavior of the embedded servlet container for request and response handling can be configured using the `server.servlet.encoding.*` configuration properties.

When a request’s `Accept-Language` header indicates a locale for the request it will be automatically mapped to a charset by the servlet container.
Each container provides default locale to charset mappings and you should verify that they meet your application’s needs.
When they do not, use the `server.servlet.encoding.mapping` configuration property to customize the mappings, as shown in the following example:

-
Properties
-
YAML

`server.servlet.encoding.mapping.ko=UTF-8````
server:
 servlet:
 encoding:
 mapping:
 ko: "UTF-8"
```
In the preceding example, the `ko` (Korean) locale has been mapped to `UTF-8`.
This is equivalent to a `<locale-encoding-mapping-list>` entry in a `web.xml` file of a traditional war deployment.

#### Programmatic Customization

If you need to programmatically configure your embedded servlet container, you can register a Spring bean that implements the `WebServerFactoryCustomizer` interface.
`WebServerFactoryCustomizer` provides access to the `ConfigurableServletWebServerFactory`, which includes numerous customization setter methods.
The following example shows programmatically setting the port:

-
Java
-
Kotlin

```
import org.springframework.boot.web.server.WebServerFactoryCustomizer;
import org.springframework.boot.web.server.servlet.ConfigurableServletWebServerFactory;
import org.springframework.stereotype.Component;
@Component
public class MyWebServerFactoryCustomizer implements WebServerFactoryCustomizer<ConfigurableServletWebServerFactory> {
	@Override
	public void customize(ConfigurableServletWebServerFactory server) {
 server.setPort(9000);
	}
}
```
```
import org.springframework.boot.web.server.servlet.ConfigurableServletWebServerFactory
import org.springframework.boot.web.server.WebServerFactoryCustomizer
import org.springframework.stereotype.Component
@Component
class MyWebServerFactoryCustomizer : WebServerFactoryCustomizer<ConfigurableServletWebServerFactory> {
	override fun customize(server: ConfigurableServletWebServerFactory) {
 server.setPort(9000)
	}
}
```
`TomcatServletWebServerFactory`, and `JettyServletWebServerFactory` are dedicated variants of `ConfigurableServletWebServerFactory` that have additional customization setter methods for Tomcat, and Jetty respectively.
The following example shows how to customize `TomcatServletWebServerFactory` that provides access to Tomcat-specific configuration options:

-
Java
-
Kotlin

```
import java.time.Duration;
import org.springframework.boot.tomcat.servlet.TomcatServletWebServerFactory;
import org.springframework.boot.web.server.WebServerFactoryCustomizer;
import org.springframework.stereotype.Component;
@Component
public class MyTomcatWebServerFactoryCustomizer implements WebServerFactoryCustomizer<TomcatServletWebServerFactory> {
	@Override
	public void customize(TomcatServletWebServerFactory server) {
 server.addConnectorCustomizers((connector) -> connector.setAsyncTimeout(Duration.ofSeconds(20).toMillis()));
	}
}
```
```
import org.springframework.boot.web.server.WebServerFactoryCustomizer
import org.springframework.boot.tomcat.servlet.TomcatServletWebServerFactory
import org.springframework.stereotype.Component
import java.time.Duration
@Component
class MyTomcatWebServerFactoryCustomizer : WebServerFactoryCustomizer<TomcatServletWebServerFactory> {
	override fun customize(server: TomcatServletWebServerFactory) {
 server.addConnectorCustomizers({ connector -> connector.asyncTimeout = Duration.ofSeconds(20).toMillis() })
	}
}
```
#### Customizing ConfigurableServletWebServerFactory Directly

For more advanced use cases that require you to extend from `ServletWebServerFactory`, you can expose a bean of such type yourself.

Setters are provided for many configuration options.
Several protected method “hooks” are also provided should you need to do something more exotic.
See the `ConfigurableServletWebServerFactory` API documentation for details.

| Auto-configured customizers are still applied on your custom factory, so use that option carefully. |

### JSP Limitations

When running a Spring Boot application that uses an embedded servlet container (and is packaged as an executable archive), there are some limitations in the JSP support.

-
With Jetty and Tomcat, it should work if you use war packaging. An executable war will work when launched with `java -jar`, and will also be deployable to any standard container. JSPs are not supported when using an executable jar.
-
Creating a custom `error.jsp`page does not override the default view for error handling. Custom error pages should be used instead.
-
If you run your application using `mvn spring-boot:run`or`gradle bootRun`and you deviate from the standard`src/main/webapp`directory structure you may need to set a`WAR_SOURCE_DIRECTORY`environment variable so that Spring Boot can find your JSPs.

# Reactive Web Applications

Spring Boot simplifies development of reactive web applications by providing auto-configuration for Spring Webflux.

## The “Spring WebFlux Framework”

Spring WebFlux is the new reactive web framework introduced in Spring Framework 5.0. Unlike Spring MVC, it does not require the servlet API, is fully asynchronous and non-blocking, and implements the Reactive Streams specification through the Reactor project.

Spring WebFlux comes in two flavors: functional and annotation-based. The annotation-based one is quite close to the Spring MVC model, as shown in the following example:

-
Java
-
Kotlin

```
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
@RestController
@RequestMapping("/users")
public class MyRestController {
	private final UserRepository userRepository;
	private final CustomerRepository customerRepository;
	public MyRestController(UserRepository userRepository, CustomerRepository customerRepository) {
 this.userRepository = userRepository;
 this.customerRepository = customerRepository;
	}
	@GetMapping("/{userId}")
	public Mono<User> getUser(@PathVariable Long userId) {
 return this.userRepository.findById(userId);
	}
	@GetMapping("/{userId}/customers")
	public Flux<Customer> getUserCustomers(@PathVariable Long userId) {
 return this.userRepository.findById(userId).flatMapMany(this.customerRepository::findByUser);
	}
	@DeleteMapping("/{userId}")
	public Mono<Void> deleteUser(@PathVariable Long userId) {
 return this.userRepository.deleteById(userId);
	}
}
```
```
import org.springframework.web.bind.annotation.DeleteMapping
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono
@RestController
@RequestMapping("/users")
class MyRestController(private val userRepository: UserRepository, private val customerRepository: CustomerRepository) {
	@GetMapping("/{userId}")
	fun getUser(@PathVariable userId: Long): Mono<User> {
 return userRepository.findById(userId)
	}
	@GetMapping("/{userId}/customers")
	fun getUserCustomers(@PathVariable userId: Long): Flux<Customer> {
 return userRepository.findById(userId).flatMapMany { user: User ->
 customerRepository.findByUser(user)
 }
	}
	@DeleteMapping("/{userId}")
	fun deleteUser(@PathVariable userId: Long): Mono<Void> {
 return userRepository.deleteById(userId)
	}
}
```
WebFlux is part of the Spring Framework and detailed information is available in its reference documentation.

“WebFlux.fn”, the functional variant, separates the routing configuration from the actual handling of the requests, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.MediaType;
import org.springframework.web.reactive.function.server.RequestPredicate;
import org.springframework.web.reactive.function.server.RouterFunction;
import org.springframework.web.reactive.function.server.ServerResponse;
import static org.springframework.web.reactive.function.server.RequestPredicates.accept;
import static org.springframework.web.reactive.function.server.RouterFunctions.route;
@Configuration(proxyBeanMethods = false)
public class MyRoutingConfiguration {
	private static final RequestPredicate ACCEPT_JSON = accept(MediaType.APPLICATION_JSON);
	@Bean
	public RouterFunction<ServerResponse> monoRouterFunction(MyUserHandler userHandler) {
 return route()
 .GET("/{user}", ACCEPT_JSON, userHandler::getUser)
 .GET("/{user}/customers", ACCEPT_JSON, userHandler::getUserCustomers)
 .DELETE("/{user}", ACCEPT_JSON, userHandler::deleteUser)
 .build();
	}
}
```
```
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.http.MediaType
import org.springframework.web.reactive.function.server.RequestPredicates.DELETE
import org.springframework.web.reactive.function.server.RequestPredicates.GET
import org.springframework.web.reactive.function.server.RequestPredicates.accept
import org.springframework.web.reactive.function.server.RouterFunction
import org.springframework.web.reactive.function.server.RouterFunctions
import org.springframework.web.reactive.function.server.ServerResponse
@Configuration(proxyBeanMethods = false)
class MyRoutingConfiguration {
	@Bean
	fun monoRouterFunction(userHandler: MyUserHandler): RouterFunction<ServerResponse> {
 return RouterFunctions.route(
 GET("/{user}").and(ACCEPT_JSON), userHandler::getUser).andRoute(
 GET("/{user}/customers").and(ACCEPT_JSON), userHandler::getUserCustomers).andRoute(
 DELETE("/{user}").and(ACCEPT_JSON), userHandler::deleteUser)
	}
	companion object {
 private val ACCEPT_JSON = accept(MediaType.APPLICATION_JSON)
	}
}
```
-
Java
-
Kotlin

```
import reactor.core.publisher.Mono;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.function.server.ServerRequest;
import org.springframework.web.reactive.function.server.ServerResponse;
@Component
public class MyUserHandler {
	public Mono<ServerResponse> getUser(ServerRequest request) {
 ...
	}
	public Mono<ServerResponse> getUserCustomers(ServerRequest request) {
 ...
	}
	public Mono<ServerResponse> deleteUser(ServerRequest request) {
 ...
	}
}
```
```
import org.springframework.stereotype.Component
import org.springframework.web.reactive.function.server.ServerRequest
import org.springframework.web.reactive.function.server.ServerResponse
import reactor.core.publisher.Mono
@Component
class MyUserHandler {
	fun getUser(request: ServerRequest?): Mono<ServerResponse> {
 ...
	}
	fun getUserCustomers(request: ServerRequest?): Mono<ServerResponse> {
 ...
	}
	fun deleteUser(request: ServerRequest?): Mono<ServerResponse> {
 ...
	}
}
```
“WebFlux.fn” is part of the Spring Framework and detailed information is available in its reference documentation.

| You can define as many `RouterFunction`beans as you like to modularize the definition of the router.
Beans can be ordered if you need to apply a precedence. |

To get started, add the `spring-boot-starter-webflux` module to your application.

| Adding both `spring-boot-starter-web`and`spring-boot-starter-webflux`modules in your application results in Spring Boot auto-configuring Spring MVC, not WebFlux.
This behavior has been chosen because many Spring developers add`spring-boot-starter-webflux`to their Spring MVC application to use the reactive`WebClient`.
You can still enforce your choice by setting the chosen application type to`SpringApplication.setWebApplicationType(WebApplicationType.REACTIVE)`. |

### Spring WebFlux Auto-configuration

Spring Boot provides auto-configuration for Spring WebFlux that works well with most applications.

The auto-configuration adds the following features on top of Spring’s defaults:

-
Configuring codecs for `HttpMessageReader`and`HttpMessageWriter`instances (described later in this document).
-
Support for serving static resources, including support for WebJars (described later in this document).

If you want to keep Spring Boot WebFlux features and you want to add additional WebFlux configuration, you can add your own `@Configuration` class of type `WebFluxConfigurer` but **without** `@EnableWebFlux`.

If you want to add additional customization to the auto-configured `HttpHandler`, you can define beans of type `WebHttpHandlerBuilderCustomizer` and use them to modify the `WebHttpHandlerBuilder`.

If you want to take complete control of Spring WebFlux, you can add your own `@Configuration` annotated with `@EnableWebFlux`.

### Spring WebFlux Conversion Service

If you want to customize the `ConversionService` used by Spring WebFlux, you can provide a `WebFluxConfigurer` bean with an `addFormatters` method.

Conversion can also be customized using the `spring.webflux.format.*` configuration properties.
When not configured, the following defaults are used:

| Property | `DateTimeFormatter` | Formats |
|---|---|---|
|
 |
 |
 |
|
 |
 | java.time’s |
|
 |
 | java.time’s |

### HTTP Codecs with HttpMessageReaders and HttpMessageWriters

Spring WebFlux uses the `HttpMessageReader` and `HttpMessageWriter` interfaces to convert HTTP requests and responses.
They are configured with `CodecConfigurer` to have sensible defaults by looking at the libraries available in your classpath.

Spring Boot provides dedicated configuration properties for codecs, `spring.http.codecs.*`.
It also applies further customization by using `CodecCustomizer` instances.
For example, `spring.jackson.*` configuration keys are applied to the Jackson codec.

If you need to add or customize codecs, you can create a custom `CodecCustomizer` component, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.http.codec.CodecCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.codec.ServerSentEventHttpMessageReader;
@Configuration(proxyBeanMethods = false)
public class MyCodecsConfiguration {
	@Bean
	public CodecCustomizer myCodecCustomizer() {
 return (configurer) -> {
 configurer.registerDefaults(false);
 configurer.customCodecs().register(new ServerSentEventHttpMessageReader());
 // ...
 };
	}
}
```
```
import org.springframework.boot.http.codec.CodecCustomizer
import org.springframework.context.annotation.Bean
import org.springframework.http.codec.CodecConfigurer
import org.springframework.http.codec.ServerSentEventHttpMessageReader
class MyCodecsConfiguration {
	@Bean
	fun myCodecCustomizer(): CodecCustomizer {
 return CodecCustomizer { configurer: CodecConfigurer ->
 configurer.registerDefaults(false)
 configurer.customCodecs().register(ServerSentEventHttpMessageReader())
 }
	}
}
```
You can also leverage Boot’s custom JSON serializers and deserializers.

### Static Content

By default, Spring Boot serves static content from a directory called `/static` (or `/public` or `/resources` or `/META-INF/resources`) in the classpath.
It uses the `ResourceWebHandler` from Spring WebFlux so that you can modify that behavior by adding your own `WebFluxConfigurer` and overriding the `addResourceHandlers` method.

By default, resources are mapped on `/**`, but you can tune that by setting the `spring.webflux.static-path-pattern` property.
For instance, relocating all resources to `/resources/**` can be achieved as follows:

-
Properties
-
YAML

`spring.webflux.static-path-pattern=/resources/**````
spring:
 webflux:
 static-path-pattern: "/resources/**"
```
You can also customize the static resource locations by using `spring.web.resources.static-locations`.
Doing so replaces the default values with a list of directory locations.
If you do so, the default welcome page detection switches to your custom locations.
So, if there is an `index.html` in any of your locations on startup, it is the home page of the application.

In addition to the “standard” static resource locations listed earlier, a special case is made for Webjars content.
By default, any resources with a path in `/webjars/**` are served from jar files if they are packaged in the Webjars format.
The path can be customized with the `spring.webflux.webjars-path-pattern` property.

| Spring WebFlux applications do not strictly depend on the servlet API, so they cannot be deployed as war files and do not use the `src/main/webapp`directory. |

### Welcome Page

Spring Boot supports both static and templated welcome pages.
It first looks for an `index.html` file in the configured static content locations.
If one is not found, it then looks for an `index` template.
If either is found, it is automatically used as the welcome page of the application.

This only acts as a fallback for actual index routes defined by the application.
The ordering is defined by the order of `HandlerMapping` beans which is by default the following:

|
 | Endpoints declared with |
|
 | Endpoints declared in |
|
 | The welcome page support |

### Template Engines

As well as REST web services, you can also use Spring WebFlux to serve dynamic HTML content. Spring WebFlux supports a variety of templating technologies, including Thymeleaf, FreeMarker, and Mustache.

Spring Boot includes auto-configuration support for the following templating engines:

| Not all FreeMarker features are supported with WebFlux. For more details, check the description of each property. |

When you use one of these templating engines with the default configuration, your templates are picked up automatically from `src/main/resources/templates`.

### Error Handling

Spring Boot provides a `WebExceptionHandler` that handles all errors in a sensible way.
Its position in the processing order is immediately before the handlers provided by WebFlux, which are considered last.
For machine clients, it produces a JSON response with details of the error, the HTTP status, and the exception message.
For browser clients, there is a “whitelabel” error handler that renders the same data in HTML format.
You can also provide your own HTML templates to display errors (see the next section).

Before customizing error handling in Spring Boot directly, you can leverage the RFC 9457 Problem Details support in Spring WebFlux.
Spring WebFlux can produce custom error messages with the `application/problem+json` media type, like:

```
{
	"type": "https://example.org/problems/unknown-project",
	"title": "Unknown project",
	"status": 404,
	"detail": "No project found for id 'spring-unknown'",
	"instance": "/projects/spring-unknown"
}
```
This support can be enabled by setting `spring.webflux.problemdetails.enabled` to `true`.

The first step to customizing this feature often involves using the existing mechanism but replacing or augmenting the error contents.
For that, you can add a bean of type `ErrorAttributes`.

To change the error handling behavior, you can implement `ErrorWebExceptionHandler` and register a bean definition of that type.
Because an `ErrorWebExceptionHandler` is quite low-level, Spring Boot also provides a convenient `AbstractErrorWebExceptionHandler` to let you handle errors in a WebFlux functional way, as shown in the following example:

-
Java
-
Kotlin

```
import reactor.core.publisher.Mono;
import org.springframework.boot.autoconfigure.web.WebProperties;
import org.springframework.boot.webflux.autoconfigure.error.AbstractErrorWebExceptionHandler;
import org.springframework.boot.webflux.error.ErrorAttributes;
import org.springframework.context.ApplicationContext;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.codec.ServerCodecConfigurer;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.function.server.RouterFunction;
import org.springframework.web.reactive.function.server.RouterFunctions;
import org.springframework.web.reactive.function.server.ServerRequest;
import org.springframework.web.reactive.function.server.ServerResponse;
import org.springframework.web.reactive.function.server.ServerResponse.BodyBuilder;
@Component
public class MyErrorWebExceptionHandler extends AbstractErrorWebExceptionHandler {
	public MyErrorWebExceptionHandler(ErrorAttributes errorAttributes, WebProperties webProperties,
 ApplicationContext applicationContext, ServerCodecConfigurer serverCodecConfigurer) {
 super(errorAttributes, webProperties.getResources(), applicationContext);
 setMessageReaders(serverCodecConfigurer.getReaders());
 setMessageWriters(serverCodecConfigurer.getWriters());
	}
	@Override
	protected RouterFunction<ServerResponse> getRoutingFunction(ErrorAttributes errorAttributes) {
 return RouterFunctions.route(this::acceptsXml, this::handleErrorAsXml);
	}
	private boolean acceptsXml(ServerRequest request) {
 return request.headers().accept().contains(MediaType.APPLICATION_XML);
	}
	public Mono<ServerResponse> handleErrorAsXml(ServerRequest request) {
 BodyBuilder builder = ServerResponse.status(HttpStatus.INTERNAL_SERVER_ERROR);
 // ... additional builder calls
 return builder.build();
	}
}
```
```
import org.springframework.boot.autoconfigure.web.WebProperties
import org.springframework.boot.webflux.error.ErrorAttributes
import org.springframework.boot.webflux.autoconfigure.error.AbstractErrorWebExceptionHandler
import org.springframework.context.ApplicationContext
import org.springframework.http.HttpStatus
import org.springframework.http.MediaType
import org.springframework.http.codec.ServerCodecConfigurer
import org.springframework.stereotype.Component
import org.springframework.web.reactive.function.server.RouterFunction
import org.springframework.web.reactive.function.server.RouterFunctions
import org.springframework.web.reactive.function.server.ServerRequest
import org.springframework.web.reactive.function.server.ServerResponse
import reactor.core.publisher.Mono
@Component
class MyErrorWebExceptionHandler(
 errorAttributes: ErrorAttributes, webProperties: WebProperties,
 applicationContext: ApplicationContext, serverCodecConfigurer: ServerCodecConfigurer
) : AbstractErrorWebExceptionHandler(errorAttributes, webProperties.resources, applicationContext) {
	init {
 setMessageReaders(serverCodecConfigurer.readers)
 setMessageWriters(serverCodecConfigurer.writers)
	}
	override fun getRoutingFunction(errorAttributes: ErrorAttributes): RouterFunction<ServerResponse> {
 return RouterFunctions.route(this::acceptsXml, this::handleErrorAsXml)
	}
	private fun acceptsXml(request: ServerRequest): Boolean {
 return request.headers().accept().contains(MediaType.APPLICATION_XML)
	}
	fun handleErrorAsXml(request: ServerRequest): Mono<ServerResponse> {
 val builder = ServerResponse.status(HttpStatus.INTERNAL_SERVER_ERROR)
 // ... additional builder calls
 return builder.build()
	}
}
```
For a more complete picture, you can also subclass `DefaultErrorWebExceptionHandler` directly and override specific methods.

In some cases, errors handled at the controller level are not recorded by web observations or the metrics infrastructure. Applications can ensure that such exceptions are recorded with the observations by setting the handled exception on the observation context.

#### Custom Error Pages

If you want to display a custom HTML error page for a given status code, you can add views that resolve from `error/*`, for example by adding files to a `/error` directory.
Error pages can either be static HTML (that is, added under any of the static resource directories) or built with templates.
The name of the file should be the exact status code, a status code series mask, or `error` for a default if nothing else matches.
Note that the path to the default error view is `error/error`, whereas with Spring MVC the default error view is `error`.

For example, to map `404` to a static HTML file, your directory structure would be as follows:

```
src/
 +- main/
 +- java/
 | + <source code>
 +- resources/
 +- public/
 +- error/
 | +- 404.html
 +- <other public assets>
```
To map all `5xx` errors by using a Mustache template, your directory structure would be as follows:

```
src/
 +- main/
 +- java/
 | + <source code>
 +- resources/
 +- templates/
 +- error/
 | +- 5xx.mustache
 +- <other templates>
```
### Web Filters

Spring WebFlux provides a `WebFilter` interface that can be implemented to filter HTTP request-response exchanges.
`WebFilter` beans found in the application context will be automatically used to filter each exchange.

Where the order of the filters is important they can implement `Ordered` or be annotated with `@Order`.
Spring Boot auto-configuration may configure web filters for you.
When it does so, the orders shown in the following table will be used:

| Web Filter | Order |
|---|---|
|
 |
 |
|
 |

### API Versioning

Spring WebFlux supports API versioning which can be used to evolve an HTTP API over time.
The same `@Controller` path can be mapped multiple times to support different versions of the API.

For more details see Spring Framework’s reference documentation.

Once mappings have been added, you additionally need to configure Spring WebFlux so that it is able to use any version information sent with a request. Typically, versions are sent as HTTP headers, query parameters, media type parameters, or as part of the path.

To configure Spring WebFlux, you can either use a `WebFluxConfigurer` bean and override the `configureApiVersioning(…)` method, or you can use properties.

For example, the following will use an `X-Version` HTTP header to obtain version information and default to `1.0.0` when no header is sent.

-
Properties
-
YAML

```
spring.webflux.apiversion.default=1.0.0
spring.webflux.apiversion.use.header=X-Version
```
```
spring:
 webflux:
 apiversion:
 default: 1.0.0
 use:
 header: X-Version
```
| If your setup requires multiple strategies, such as header and query parameter, consider declaring the order programmatically by overriding the `configureApiVersioning`method. |

For more complete control, you can also define `ApiVersionResolver`, `ApiVersionParser` and `ApiVersionDeprecationHandler` beans which will be injected into the auto-configured Spring WebFlux configuration.

| API versioning is also supported on the client-side with both `WebClient`and`RestClient`.
See API Versioning for details. |

## Embedded Reactive Server Support

Spring Boot includes support for the following embedded reactive web servers: Reactor Netty, Tomcat, and Jetty. Most developers use the appropriate starter to obtain a fully configured instance. By default, the embedded server listens for HTTP requests on port 8080.

### Customizing Reactive Servers

Common reactive web server settings can be configured by using Spring `Environment` properties.
Usually, you would define the properties in your `application.properties` or `application.yaml` file.

Common server settings include:

-
Network settings: Listen port for incoming HTTP requests ( `server.port`), interface address to bind to (`server.address`), and so on.
-
Error management: Location of the error page ( `spring.web.error.path`) and so on.

Spring Boot tries as much as possible to expose common settings, but this is not always possible.
For those cases, dedicated namespaces such as `server.netty.*` offer server-specific customizations.

| See the `ServerProperties`class for a complete list. |

#### Programmatic Customization

If you need to programmatically configure your reactive web server, you can register a Spring bean that implements the `WebServerFactoryCustomizer` interface.
`WebServerFactoryCustomizer` provides access to the `ConfigurableReactiveWebServerFactory`, which includes numerous customization setter methods.
The following example shows programmatically setting the port:

-
Java
-
Kotlin

```
import org.springframework.boot.web.server.WebServerFactoryCustomizer;
import org.springframework.boot.web.server.reactive.ConfigurableReactiveWebServerFactory;
import org.springframework.stereotype.Component;
@Component
public class MyWebServerFactoryCustomizer implements WebServerFactoryCustomizer<ConfigurableReactiveWebServerFactory> {
	@Override
	public void customize(ConfigurableReactiveWebServerFactory server) {
 server.setPort(9000);
	}
}
```
```
import org.springframework.boot.web.server.WebServerFactoryCustomizer
import org.springframework.boot.web.server.reactive.ConfigurableReactiveWebServerFactory
import org.springframework.stereotype.Component
@Component
class MyWebServerFactoryCustomizer : WebServerFactoryCustomizer<ConfigurableReactiveWebServerFactory> {
	override fun customize(server: ConfigurableReactiveWebServerFactory) {
 server.setPort(9000)
	}
}
```
`JettyReactiveWebServerFactory`, `NettyReactiveWebServerFactory`, and `TomcatReactiveWebServerFactory` are dedicated variants of `ConfigurableReactiveWebServerFactory` that have additional customization setter methods for Jetty, Reactor Netty, and Tomcat respectively.
The following example shows how to customize `NettyReactiveWebServerFactory` that provides access to Reactor Netty-specific configuration options:

-
Java
-
Kotlin

```
import java.time.Duration;
import org.springframework.boot.reactor.netty.NettyReactiveWebServerFactory;
import org.springframework.boot.web.server.WebServerFactoryCustomizer;
import org.springframework.stereotype.Component;
@Component
public class MyNettyWebServerFactoryCustomizer implements WebServerFactoryCustomizer<NettyReactiveWebServerFactory> {
	@Override
	public void customize(NettyReactiveWebServerFactory factory) {
 factory.addServerCustomizers((server) -> server.idleTimeout(Duration.ofSeconds(20)));
	}
}
```
```
import org.springframework.boot.web.server.WebServerFactoryCustomizer
import org.springframework.boot.reactor.netty.NettyReactiveWebServerFactory
import org.springframework.stereotype.Component
import java.time.Duration
@Component
class MyNettyWebServerFactoryCustomizer : WebServerFactoryCustomizer<NettyReactiveWebServerFactory> {
	override fun customize(factory: NettyReactiveWebServerFactory) {
 factory.addServerCustomizers({ server -> server.idleTimeout(Duration.ofSeconds(20)) })
	}
}
```
#### Customizing ConfigurableReactiveWebServerFactory Directly

For more advanced use cases that require you to extend from `ReactiveWebServerFactory`, you can expose a bean of such type yourself.

Setters are provided for many configuration options.
Several protected method “hooks” are also provided should you need to do something more exotic.
See the `ConfigurableReactiveWebServerFactory` API documentation for details.

| Auto-configured customizers are still applied on your custom factory, so use that option carefully. |

## Reactive Server Resources Configuration

When auto-configuring a Reactor Netty or Jetty server, Spring Boot will create specific beans that will provide HTTP resources to the server instance: `ReactorResourceFactory` or `JettyResourceFactory`.

By default, those resources will be also shared with the Reactor Netty and Jetty clients for optimal performances, given:

-
the same technology is used for server and client
-
the client instance is built using the `WebClient.Builder`bean auto-configured by Spring Boot

Developers can override the resource configuration for Jetty and Reactor Netty by providing a custom `ReactorResourceFactory` or `JettyResourceFactory` bean - this will be applied to both clients and servers.

You can learn more about the resource configuration on the client side in the WebClient Runtime section.

# Graceful Shutdown

Graceful shutdown is enabled by default with all three embedded web servers (Jetty, Reactor Netty, and Tomcat) and with both reactive and servlet-based web applications.
It occurs as part of closing the application context and is performed in the earliest phase of stopping `SmartLifecycle` beans.
This stop processing uses a timeout which provides a grace period during which existing requests will be allowed to complete but no new requests will be permitted.

To configure the timeout period, configure the `spring.lifecycle.timeout-per-shutdown-phase` property, as shown in the following example:

-
Properties
-
YAML

`spring.lifecycle.timeout-per-shutdown-phase=20s````
spring:
 lifecycle:
 timeout-per-shutdown-phase: "20s"
```
| Shutdown in your IDE may be immediate rather than graceful if it does not send a proper `SIGTERM`signal.
See the documentation of your IDE for more details. |

## Rejecting Requests During the Grace Period

The exact way in which new requests are not permitted varies depending on the web server that is being used. Implementations may stop accepting requests at the network layer, or they may return a response with a specific HTTP status code or HTTP header. The use of persistent connections can also change the way that requests stop being accepted.

| To learn more about the specific method used with your web server, see the `shutDownGracefully`API documentation for`TomcatWebServer.shutDownGracefully(GracefulShutdownCallback)`,`NettyWebServer.shutDownGracefully(GracefulShutdownCallback)`, or`JettyWebServer.shutDownGracefully(GracefulShutdownCallback)`. |

Jetty, Reactor Netty, and Tomcat will stop accepting new requests at the network layer.

# Spring Security

If Spring Security is on the classpath, then web applications are secured by default.
This includes securing Spring Boot’s `/error` endpoint.
Spring Boot relies on Spring Security’s content-negotiation strategy to determine whether to use `httpBasic` or `formLogin`.
To add method-level security to a web application, you can also add `@EnableMethodSecurity` with your desired settings.
Additional information can be found in the Spring Security Reference Guide.

The default `UserDetailsService` has a single user.
The user name is `user`, and the password is random and is printed at WARN level when the application starts, as shown in the following example:

```
Using generated security password: 78fa095d-3f4c-48b1-ad50-e24c31d5cf35
This generated password is for development use only. Your security configuration must be updated before running your application in production.
```
| If you fine-tune your logging configuration, ensure that the `org.springframework.boot.security.autoconfigure`category is set to log`WARN`-level messages.
Otherwise, the default password is not printed. |

You can change the username and password by providing a `spring.security.user.name` and `spring.security.user.password`.

The basic features you get by default in a web application are:

-
A `UserDetailsService`(or`ReactiveUserDetailsService`in case of a WebFlux application) bean with in-memory store and a single user with a generated password (see`SecurityProperties.User`for the properties of the user).
-
Form-based login or HTTP Basic security (depending on the `Accept`header in the request) for the entire application (including actuator endpoints if actuator is on the classpath).
-
A `DefaultAuthenticationEventPublisher`for publishing authentication events.

You can provide a different `AuthenticationEventPublisher` by adding a bean for it.

## MVC Security

The default security configuration is implemented in `SecurityAutoConfiguration` and `UserDetailsServiceAutoConfiguration`.
`SecurityAutoConfiguration` imports `SpringBootWebSecurityConfiguration` for web security and `UserDetailsServiceAutoConfiguration` for authentication.

To completely switch off the default web application security configuration, including Actuator security, or to combine multiple Spring Security components such as OAuth2 Client and Resource Server, add a bean of type `SecurityFilterChain` (doing so does not disable the `UserDetailsService` configuration).
To also switch off the `UserDetailsService` configuration, add a bean of type `UserDetailsService`, `AuthenticationProvider`, or `AuthenticationManager`.

The auto-configuration of a `UserDetailsService` will also back off when any of the following Spring Security modules is on the classpath:

-
`spring-security-oauth2-client`
-
`spring-security-oauth2-resource-server`
-
`spring-security-saml2-service-provider`

To use `UserDetailsService` in addition to one or more of these dependencies, define your own `InMemoryUserDetailsManager` bean.

Access rules can be overridden by adding a custom `SecurityFilterChain` bean.
Spring Boot provides convenience methods that can be used to override access rules for actuator endpoints and static resources.
`EndpointRequest` can be used to create a `RequestMatcher` that is based on the `management.endpoints.web.base-path` property.
`PathRequest` can be used to create a `RequestMatcher` for resources in commonly used locations.

## WebFlux Security

Similar to Spring MVC applications, you can secure your WebFlux applications by adding the `spring-boot-starter-security` dependency.
The default security configuration is implemented in `ReactiveWebSecurityAutoConfiguration` and `ReactiveUserDetailsServiceAutoConfiguration`.
In addition to reactive web applications, the latter is also auto-configured when RSocket is in use.
`ReactiveWebSecurityAutoConfiguration` imports `WebFluxSecurityConfiguration` for web security.
`ReactiveUserDetailsServiceAutoConfiguration` auto-configures authentication.

To completely switch off the default web application security configuration, including Actuator security, add a bean of type `WebFilterChainProxy` (doing so does not disable the `ReactiveUserDetailsService` configuration).
To also switch off the `ReactiveUserDetailsService` configuration, add a bean of type `ReactiveUserDetailsService` or `ReactiveAuthenticationManager`.

The auto-configuration will also back off when any of the following Spring Security modules is on the classpath:

-
`spring-security-oauth2-client`
-
`spring-security-oauth2-resource-server`

To use `ReactiveUserDetailsService` in addition to one or more of these dependencies, define your own `MapReactiveUserDetailsService` bean.

Access rules and the use of multiple Spring Security components such as OAuth 2 Client and Resource Server can be configured by adding a custom `SecurityWebFilterChain` bean.
Spring Boot provides convenience methods that can be used to override access rules for actuator endpoints and static resources.
`EndpointRequest` can be used to create a `ServerWebExchangeMatcher` that is based on the `management.endpoints.web.base-path` property.

`PathRequest` can be used to create a `ServerWebExchangeMatcher` for resources in commonly used locations.

For example, you can customize your security configuration by adding something like:

```
import org.springframework.boot.security.autoconfigure.web.reactive.PathRequest;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.web.server.SecurityWebFilterChain;
import static org.springframework.security.config.Customizer.withDefaults;
@Configuration(proxyBeanMethods = false)
public class MyWebFluxSecurityConfiguration {
	@Bean
	public SecurityWebFilterChain springSecurityFilterChain(ServerHttpSecurity http) {
 http.authorizeExchange((exchange) -> {
 exchange.matchers(PathRequest.toStaticResources().atCommonLocations()).permitAll();
 exchange.pathMatchers("/foo", "/bar").authenticated();
 });
 http.formLogin(withDefaults());
 return http.build();
	}
}
```
## OAuth2

OAuth2 is a widely used authorization framework. For details of how to configure and use OAuth2 with your web applications, see the “OAuth2” section of under “Security”.

## SAML 2.0

SAML v2.0 is a widely adopted framework for exchanging security information between online business partners. For details of how to configure and use SAML 2.0 with your web applications, see the “SAML 2.0” section under “Security”.

# Spring Session

Spring Boot provides Spring Session auto-configuration for a range of data stores. For each data store, a specific Spring Boot starter is provided.

When building a servlet web application, the following stores can be auto-configured:

-
Redis ( `spring-boot-starter-session-data-redis`)
-
JDBC ( `spring-boot-starter-session-jdbc`)

The servlet auto-configuration replaces the need to use `@Enable*HttpSession`.

When building a reactive web application, the Redis store can be auto-configured by depending on `spring-boot-starter-session-data-redis`.
This replaces the need to use `@EnableRedisWebSession`.

Each store has specific additional settings. For instance, it is possible to customize the name of the table for the JDBC store, as shown in the following example:

-
Properties
-
YAML

`spring.session.jdbc.table-name=SESSIONS````
spring:
 session:
 jdbc:
 table-name: "SESSIONS"
```
For setting the timeout of the session you can use the `spring.session.timeout` property.
If that property is not set with a servlet web application, the auto-configuration falls back to the value of `server.servlet.session.timeout`.
To provide the timeout programmatically, define a `SessionTimeout` bean.

You can take control over Spring Session’s configuration using `@Enable*HttpSession` (servlet) or `@EnableRedisWebSession` (reactive).
This will cause the auto-configuration to back off.
Alternatively, depend on the relevant Spring Session module directly rather than using one of Spring Boot’s starters for Spring Session.
With either approach, Spring Session can then be configured using the annotation’s attributes rather than the previously described configuration properties.

# Spring for GraphQL

If you want to build GraphQL applications, you can take advantage of Spring Boot’s auto-configuration for Spring for GraphQL.
The Spring for GraphQL project is based on GraphQL Java.
You’ll need the `spring-boot-starter-graphql` starter at a minimum.
Because GraphQL is transport-agnostic, you’ll also need to have one or more additional starters in your application to expose your GraphQL API over the web:

| Starter | Transport | Implementation |
|---|---|---|
|
 | HTTP | Spring MVC |
|
 | WebSocket | WebSocket for Servlet apps |
|
 | HTTP, WebSocket | Spring WebFlux |
|
 | TCP, WebSocket | Spring WebFlux on Reactor Netty |

## GraphQL Schema

A Spring GraphQL application requires a defined schema at startup.
By default, you can write ".graphqls" or ".gqls" schema files under `src/main/resources/graphql/**` and Spring Boot will pick them up automatically.
You can customize the locations with `spring.graphql.schema.locations` and the file extensions with `spring.graphql.schema.file-extensions`.

| If you want Spring Boot to detect schema files in all your application modules and dependencies for that location,
you can set `spring.graphql.schema.locations`to`"classpath*:graphql/**/"`(note the`classpath*:`prefix). |

In the following sections, we’ll consider this sample GraphQL schema, defining two types and two queries:

```
type Query {
 greeting(name: String! = "Spring"): String!
 project(slug: ID!): Project
}
""" A Project in the Spring portfolio """
type Project {
 """ Unique string id used in URLs """
 slug: ID!
 """ Project name """
 name: String!
 """ URL of the git repository """
 repositoryUrl: String!
 """ Current support status """
 status: ProjectStatus!
}
enum ProjectStatus {
 """ Actively supported by the Spring team """
 ACTIVE
 """ Supported by the community """
 COMMUNITY
 """ Prototype, not officially supported yet """
 INCUBATING
 """ Project being retired, in maintenance mode """
 ATTIC
 """ End-Of-Lifed """
 EOL
}
```
| By default, field introspection will be allowed on the schema as it is required for tools such as GraphiQL.
If you wish to not expose information about the schema, you can disable introspection by setting `spring.graphql.schema.introspection.enabled`to`false`. |

## GraphQL RuntimeWiring

The GraphQL Java `RuntimeWiring.Builder` can be used to register custom scalar types, directives, type resolvers, `DataFetcher`, and more.
You can declare `RuntimeWiringConfigurer` beans in your Spring config to get access to the `RuntimeWiring.Builder`.
Spring Boot detects such beans and adds them to the GraphQlSource builder.

Typically, however, applications will not implement `DataFetcher` directly and will instead create annotated controllers.
Spring Boot will automatically detect `@Controller` classes with annotated handler methods and register those as `DataFetcher`s.
Here’s a sample implementation for our greeting query with a `@Controller` class:

-
Java
-
Kotlin

```
import org.springframework.graphql.data.method.annotation.Argument;
import org.springframework.graphql.data.method.annotation.QueryMapping;
import org.springframework.stereotype.Controller;
@Controller
public class GreetingController {
	@QueryMapping
	public String greeting(@Argument String name) {
 return "Hello, " + name + "!";
	}
}
```
```
import org.springframework.graphql.data.method.annotation.Argument
import org.springframework.graphql.data.method.annotation.QueryMapping
import org.springframework.stereotype.Controller
@Controller
class GreetingController {
	@QueryMapping
	fun greeting(@Argument name: String): String {
 return "Hello, $name!"
	}
}
```
## Querydsl and QueryByExample Repositories Support

Spring Data offers support for both Querydsl and QueryByExample repositories.
Spring GraphQL can configure Querydsl and QueryByExample repositories as `DataFetcher`.

Spring Data repositories annotated with `@GraphQlRepository` and extending one of:

are detected by Spring Boot and considered as candidates for `DataFetcher` for matching top-level queries.

## Transports

### HTTP and WebSocket

The GraphQL HTTP endpoint is at HTTP POST `/graphql` by default.
It also supports the `"text/event-stream"` media type over Server Sent Events for subscriptions only.
The path can be customized with `spring.graphql.http.path`.

| The HTTP endpoint for both Spring MVC and Spring WebFlux is provided by a `RouterFunction`bean with an`@Order`of`0`.
If you define your own`RouterFunction`beans, you may want to add appropriate`@Order`annotations to ensure that they are sorted correctly. |

The GraphQL WebSocket endpoint is off by default. To enable it:

-
For a Servlet application, add the WebSocket starter `spring-boot-starter-websocket`
-
For a WebFlux application, no additional dependency is required
-
For both, the `spring.graphql.websocket.path`application property must be set

Spring GraphQL provides a Web Interception model.
This is quite useful for retrieving information from an HTTP request header and set it in the GraphQL context or fetching information from the same context and writing it to a response header.
With Spring Boot, you can declare a `WebGraphQlInterceptor` bean to have it registered with the web transport.

Spring MVC and Spring WebFlux support CORS (Cross-Origin Resource Sharing) requests. CORS is a critical part of the web config for GraphQL applications that are accessed from browsers using different domains.

Spring Boot supports many configuration properties under the `spring.graphql.cors.*` namespace; here’s a short configuration sample:

-
Properties
-
YAML

```
spring.graphql.cors.allowed-origins=https://example.org
spring.graphql.cors.allowed-methods=GET,POST
spring.graphql.cors.max-age=1800s
```
```
spring:
 graphql:
 cors:
 allowed-origins: "https://example.org"
 allowed-methods: GET,POST
 max-age: 1800s
```
### RSocket

RSocket is also supported as a transport, on top of WebSocket or TCP.
Once the RSocket server is configured, we can configure our GraphQL handler on a particular route using `spring.graphql.rsocket.mapping`.
For example, configuring that mapping as `"graphql"` means we can use that as a route when sending requests with the `RSocketGraphQlClient`.

Spring Boot auto-configures a `RSocketGraphQlClient.Builder<?>` bean that you can inject in your components:

-
Java
-
Kotlin

```
@Component
public class RSocketGraphQlClientExample {
	private final RSocketGraphQlClient graphQlClient;
	public RSocketGraphQlClientExample(RSocketGraphQlClient.Builder<?> builder) {
 this.graphQlClient = builder.tcp("example.spring.io", 8181).route("graphql").build();
	}
```
```
@Component
class RSocketGraphQlClientExample(private val builder: RSocketGraphQlClient.Builder<*>) {
```
And then send a request: include-code::RSocketGraphQlClientExample[tag=request]

## Exception Handling

Spring GraphQL enables applications to register one or more Spring `DataFetcherExceptionResolver` components that are invoked sequentially.
The Exception must be resolved to a list of `GraphQLError` objects, see Spring GraphQL exception handling documentation.
Spring Boot will automatically detect `DataFetcherExceptionResolver` beans and register them with the `GraphQlSource.Builder`.

## GraphiQL and Schema Printer

Spring GraphQL offers infrastructure for helping developers when consuming or developing a GraphQL API.

Spring GraphQL ships with a default GraphiQL page that is exposed at `"/graphiql"` by default.
This page is disabled by default and can be turned on with the `spring.graphql.graphiql.enabled` property.
Many applications exposing such a page will prefer a custom build.
A default implementation is very useful during development, this is why it is exposed automatically with `spring-boot-devtools` during development.

You can also choose to expose the GraphQL schema in text format at `/graphql/schema` when the `spring.graphql.schema.printer.enabled` property is enabled.

# Spring HATEOAS

If you develop a RESTful API that makes use of hypermedia, Spring Boot provides auto-configuration for Spring HATEOAS that works well with most applications.
The auto-configuration replaces the need to use `@EnableHypermediaSupport` and registers a number of beans to ease building hypermedia-based applications, including a `LinkDiscoverers` (for client side support) and an `JsonMapper` configured to correctly marshal responses into the desired representation.
The `JsonMapper` is customized by setting the various `spring.jackson.*` properties or, if any exist, the `JsonMapperBuilderCustomizer` beans.

You can take control of Spring HATEOAS’s configuration by using `@EnableHypermediaSupport`.
Note that doing so disables the `JsonMapper` customization described earlier.

| `spring-boot-starter-hateoas`is specific to Spring MVC and should not be combined with Spring WebFlux.
In order to use Spring HATEOAS with Spring WebFlux, you can add a direct dependency on`org.springframework.hateoas:spring-hateoas`along with`spring-boot-starter-webflux`. |

By default, requests that accept `application/json` will receive an `application/hal+json` response.
To disable this behavior set `spring.hateoas.use-hal-as-default-json-media-type` to `false` and define a `HypermediaMappingInformation` or `HalConfiguration` to configure Spring HATEOAS to meet the needs of your application and its clients.

Spring Boot Reference Data Data Spring Boot integrates with a number of data technologies, both SQL and NoSQL. Spring HATEOAS SQL Databases

# SQL Databases

The Spring Framework provides extensive support for working with SQL databases, from direct JDBC access using `JdbcClient` or `JdbcTemplate` to complete “object relational mapping” technologies such as Hibernate.
Spring Data provides an additional level of functionality: creating `Repository` implementations directly from interfaces and using conventions to generate queries from your method names.

## Configure a DataSource

Java’s `DataSource` interface provides a standard method of working with database connections.
Traditionally, a `DataSource` uses a `URL` along with some credentials to establish a database connection.

| See the Configure a Custom DataSource section of the “How-to Guides” for more advanced examples, typically to take full control over the configuration of the DataSource. |

### Embedded Database Support

It is often convenient to develop applications by using an in-memory embedded database. Obviously, in-memory databases do not provide persistent storage. You need to populate your database when your application starts and be prepared to throw away data when your application ends.

| The “How-to Guides” section includes a section on how to initialize a database. |

Spring Boot can auto-configure embedded H2, HSQL, and Derby (deprecated) databases.
You need not provide any connection URLs.
You need only include a build dependency to the embedded database that you want to use.
If there are multiple embedded databases on the classpath, set the `spring.datasource.embedded-database-connection` configuration property to control which one is used.
Setting the property to `none` disables auto-configuration of an embedded database.

| If you are using this feature in your tests, you may notice that the same database is reused by your whole test suite regardless of the number of application contexts that you use.
If you want to make sure that each context has a separate embedded database, you should set |

For example, the typical POM dependencies would be as follows:

```
<dependency>
	<groupId>org.springframework.boot</groupId>
	<artifactId>spring-boot-starter-data-jpa</artifactId>
</dependency>
<dependency>
	<groupId>org.hsqldb</groupId>
	<artifactId>hsqldb</artifactId>
	<scope>runtime</scope>
</dependency>
```
| You need a dependency on `spring-jdbc`for an embedded database to be auto-configured.
In this example, it is pulled in transitively through`spring-boot-starter-data-jpa`. |

| If, for whatever reason, you do configure the connection URL for an embedded database, take care to ensure that the database’s automatic shutdown is disabled.
If you use H2, you should use `DB_CLOSE_ON_EXIT=FALSE`to do so.
If you use HSQLDB, you should ensure that`shutdown=true`is not used.
Disabling the database’s automatic shutdown lets Spring Boot control when the database is closed, thereby ensuring that it happens once access to the database is no longer needed. |

### Connection to a Production Database

Production database connections can also be auto-configured by using a pooling `DataSource`.

### DataSource Configuration

DataSource configuration is controlled by external configuration properties in `spring.datasource.*`.
For example, you might declare the following section in `application.properties`:

-
Properties
-
YAML

```
spring.datasource.url=jdbc:mysql://localhost/test
spring.datasource.username=dbuser
spring.datasource.password=dbpass
```
```
spring:
 datasource:
 url: "jdbc:mysql://localhost/test"
 username: "dbuser"
 password: "dbpass"
```
| You should at least specify the URL by setting the `spring.datasource.url`property.
Otherwise, Spring Boot tries to auto-configure an embedded database. |

| Spring Boot can deduce the JDBC driver class for most databases from the URL.
If you need to specify a specific class, you can use the `spring.datasource.driver-class-name`property. |

| For a pooling `DataSource`to be created, we need to be able to verify that a valid`Driver`class is available, so we check for that before doing anything.
In other words, if you set`spring.datasource.driver-class-name=com.mysql.jdbc.Driver`, then that class has to be loadable. |

See `DataSourceProperties` API documentation for more of the supported options.
These are the standard options that work regardless of the actual implementation.
It is also possible to fine-tune implementation-specific settings by using their respective prefix (`spring.datasource.hikari.*`, `spring.datasource.tomcat.*`, `spring.datasource.dbcp2.*`, and `spring.datasource.oracleucp.*`).
See the documentation of the connection pool implementation you are using for more details.

For instance, if you use the Tomcat connection pool, you could customize many additional settings, as shown in the following example:

-
Properties
-
YAML

```
spring.datasource.tomcat.max-wait=10000
spring.datasource.tomcat.max-active=50
spring.datasource.tomcat.test-on-borrow=true
```
```
spring:
 datasource:
 tomcat:
 max-wait: 10000
 max-active: 50
 test-on-borrow: true
```
This will set the pool to wait 10000ms before throwing an exception if no connection is available, limit the maximum number of connections to 50 and validate the connection before borrowing it from the pool.

### Supported Connection Pools

Spring Boot uses the following algorithm for choosing a specific implementation:

-
We prefer HikariCP for its performance and concurrency. If HikariCP is available, we always choose it.
-
Otherwise, if the Tomcat pooling `DataSource`is available, we use it.
-
Otherwise, if Commons DBCP2 is available, we use it.
-
If none of HikariCP, Tomcat, and DBCP2 are available and if Oracle UCP is available, we use it.

| If you use the `spring-boot-starter-jdbc`or`spring-boot-starter-data-jpa`starters, you automatically get a dependency to HikariCP. |

You can bypass that algorithm completely and specify the connection pool to use by setting the `spring.datasource.type` property.
This is especially important if you run your application in a Tomcat container, as `tomcat-jdbc` is provided by default.

Additional connection pools can always be configured manually, using `DataSourceBuilder`.
If you define your own `DataSource` bean, auto-configuration does not occur.
The following connection pools are supported by `DataSourceBuilder`:

-
HikariCP
-
Tomcat pooling `DataSource`
-
Commons DBCP2
-
Oracle UCP & `OracleDataSource`
-
Spring Framework’s `SimpleDriverDataSource`
-
PostgreSQL `PGSimpleDataSource`
-
C3P0
-
Vibur

### Connection to a JNDI DataSource

If you deploy your Spring Boot application to an Application Server, you might want to configure and manage your DataSource by using your Application Server’s built-in features and access it by using JNDI.

The `spring.datasource.jndi-name` property can be used as an alternative to the `spring.datasource.url`, `spring.datasource.username`, and `spring.datasource.password` properties to access the `DataSource` from a specific JNDI location.
For example, the following section in `application.properties` shows how you can access a JBoss AS defined `DataSource`:

-
Properties
-
YAML

`spring.datasource.jndi-name=java:jboss/datasources/customers````
spring:
 datasource:
 jndi-name: "java:jboss/datasources/customers"
```
### Lazy Connection Proxy

When a pooled `DataSource` is auto-configured, it can be wrapped in a proxy that fetches JDBC connections as late as possible by setting `spring.datasource.connection-fetch` to `lazy`, as shown in the following example:

-
Properties
-
YAML

`spring.datasource.connection-fetch=lazy````
spring:
 datasource:
 connection-fetch: "lazy"
```
With this feature enabled, JDBC Connections are only fetched from the pool when actually necessary.
JDBC transaction control can happen without fetching a Connection from the pool or communicating with the database; this will be done lazily on the first creation of a JDBC Statement.
To get access to the target auto-configured `DataSource`, use `DataSource.unwrap`.

| If you are wrapping the target `DataSource`in a`TransactionAwareDataSourceProxy`, make sure to do so with a`BeanPostProcessor`that has a positive order. |

See `LazyConnectionDataSourceProxy` for more details.

## Using JdbcTemplate

Spring’s `JdbcTemplate` and `NamedParameterJdbcTemplate` classes are auto-configured, and you can autowire them directly into your own beans, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final JdbcTemplate jdbcTemplate;
	public MyBean(JdbcTemplate jdbcTemplate) {
 this.jdbcTemplate = jdbcTemplate;
	}
	public void doSomething() {
 this.jdbcTemplate ...
	}
}
```
```
import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.stereotype.Component
@Component
class MyBean(private val jdbcTemplate: JdbcTemplate) {
	fun doSomething() {
 jdbcTemplate.execute("delete from customer")
	}
}
```
You can customize some properties of the template by using the `spring.jdbc.template.*` properties, as shown in the following example:

-
Properties
-
YAML

`spring.jdbc.template.max-rows=500````
spring:
 jdbc:
 template:
 max-rows: 500
```
If tuning of SQL exceptions is required, you can define your own `SQLExceptionTranslator` bean so that it is associated with the auto-configured `JdbcTemplate`.

| The `NamedParameterJdbcTemplate`reuses the same`JdbcTemplate`instance behind the scenes.
If more than one`JdbcTemplate`is defined and no primary candidate exists, the`NamedParameterJdbcTemplate`is not auto-configured. |

## Using JdbcClient

Spring’s `JdbcClient` is auto-configured based on the presence of a `NamedParameterJdbcTemplate`.
You can inject it directly in your own beans as well, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.jdbc.core.simple.JdbcClient;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final JdbcClient jdbcClient;
	public MyBean(JdbcClient jdbcClient) {
 this.jdbcClient = jdbcClient;
	}
	public void doSomething() {
 this.jdbcClient ...
	}
}
```
```
import org.springframework.jdbc.core.simple.JdbcClient
import org.springframework.stereotype.Component
@Component
class MyBean(private val jdbcClient: JdbcClient) {
	fun doSomething() {
 jdbcClient.sql("delete from customer").update()
	}
}
```
If you rely on auto-configuration to create the underlying `JdbcTemplate`, any customization using `spring.jdbc.template.*` properties is taken into account in the client as well.

## JPA and Spring Data JPA

The Java Persistence API is a standard technology that lets you “map” objects to relational databases.
The `spring-boot-starter-data-jpa` POM provides a quick way to get started.
It provides the following key dependencies:

-
Hibernate: One of the most popular JPA implementations.
-
Spring Data JPA: Helps you to implement JPA-based repositories.
-
Spring ORM: Core ORM support from the Spring Framework.

| We do not go into too many details of JPA or Spring Data here. You can follow the Accessing Data with JPA guide from spring.io and read the Spring Data JPA and Hibernate reference documentation. |

### Entity Classes

Traditionally, JPA “Entity” classes are specified in a `persistence.xml` file.
With Spring Boot, this file is not necessary and “Entity Scanning” is used instead.
By default the auto-configuration packages are scanned.

Any classes annotated with `@Entity`, `@Embeddable`, or `@MappedSuperclass` are considered.
A typical entity class resembles the following example:

-
Java
-
Kotlin

```
import java.io.Serializable;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.Id;
@Entity
public class City implements Serializable {
	@Id
	@GeneratedValue
	private Long id;
	@Column(nullable = false)
	private String name;
	@Column(nullable = false)
	private String state;
	// ... additional members, often include @OneToMany mappings
	protected City() {
 // no-args constructor required by JPA spec
 // this one is protected since it should not be used directly
	}
	public City(String name, String state) {
 this.name = name;
 this.state = state;
	}
	public String getName() {
 return this.name;
	}
	public String getState() {
 return this.state;
	}
	// ... etc
}
```
```
import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.GeneratedValue
import jakarta.persistence.Id
import java.io.Serializable
@Entity
class City : Serializable {
	@Id
	@GeneratedValue
	private val id: Long? = null
	@Column(nullable = false)
	var name: String? = null
 private set
	// ... etc
	@Column(nullable = false)
	var state: String? = null
 private set
	// ... additional members, often include @OneToMany mappings
	protected constructor() {
 // no-args constructor required by JPA spec
 // this one is protected since it should not be used directly
	}
	constructor(name: String?, state: String?) {
 this.name = name
 this.state = state
	}
}
```
| You can customize entity scanning locations by using the `@EntityScan`annotation.
See the Separate @Entity Definitions from Spring Configuration section of the “How-to Guides”. |

### Spring Data JPA Repositories

Spring Data JPA repositories are interfaces that you can define to access data.
JPA queries are created automatically from your method names.
For example, a `CityRepository` interface might declare a `findAllByState(String state)` method to find all the cities in a given state.

For more complex queries, you can annotate your method with Spring Data’s `Query` annotation.

Spring Data repositories usually extend from the `Repository` or `CrudRepository` interfaces.
If you use auto-configuration, the auto-configuration packages are searched for repositories.

| You can customize the locations to look for repositories using `@EnableJpaRepositories`. |

The following example shows a typical Spring Data repository interface definition:

-
Java
-
Kotlin

```
import org.springframework.boot.docs.data.sql.jpaandspringdata.entityclasses.City;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.repository.Repository;
public interface CityRepository extends Repository<City, Long> {
	Page<City> findAll(Pageable pageable);
	City findByNameAndStateAllIgnoringCase(String name, String state);
}
```
```
import org.springframework.boot.docs.data.sql.jpaandspringdata.entityclasses.City
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.Repository
interface CityRepository : Repository<City, Long> {
	fun findAll(pageable: Pageable?): Page<City>?
	fun findByNameAndStateAllIgnoringCase(name: String?, state: String?): City?
}
```
Spring Data JPA repositories support three different modes of bootstrapping: default, deferred, and lazy.
To enable deferred or lazy bootstrapping, set the `spring.data.jpa.repositories.bootstrap-mode` property to `deferred` or `lazy` respectively.
When using deferred or lazy bootstrapping, the auto-configured `EntityManagerFactoryBuilder` will use the context’s `AsyncTaskExecutor`, if any, as the bootstrap executor.
If more than one exists, the one named `applicationTaskExecutor` will be used.

| When using deferred or lazy bootstrapping, make sure to defer any access to the JPA infrastructure after the application context bootstrap phase.
You can use |

| We have barely scratched the surface of Spring Data JPA. For complete details, see the Spring Data JPA reference documentation. |

### Spring Data Envers Repositories

If Spring Data Envers is available, JPA repositories are auto-configured to support typical Envers queries.

To use Spring Data Envers, make sure your repository extends from `RevisionRepository` as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.docs.data.sql.jpaandspringdata.entityclasses.Country;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.repository.Repository;
import org.springframework.data.repository.history.RevisionRepository;
public interface CountryRepository extends RevisionRepository<Country, Long, Integer>, Repository<Country, Long> {
	Page<Country> findAll(Pageable pageable);
}
```
```
import org.springframework.boot.docs.data.sql.jpaandspringdata.entityclasses.Country
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.Repository
import org.springframework.data.repository.history.RevisionRepository
interface CountryRepository :
 RevisionRepository<Country, Long, Int>,
 Repository<Country, Long> {
	fun findAll(pageable: Pageable?): Page<Country>?
}
```
| For more details, check the Spring Data Envers reference documentation. |

### Creating and Dropping JPA Databases

By default, JPA databases are automatically created **only** if you use an embedded database, that is H2, HSQL, or Derby (deprecated).
You can explicitly configure JPA settings by using `spring.jpa.*` properties.
For example, to create and drop tables you can add the following line to your `application.properties`:

-
Properties
-
YAML

`spring.jpa.hibernate.ddl-auto=create-drop````
spring:
 jpa:
 hibernate.ddl-auto: "create-drop"
```
| Hibernate’s own internal property name for this (if you happen to remember it better) is `hibernate.hbm2ddl.auto`.
You can set it, along with other Hibernate native properties, by using`spring.jpa.properties.*`(the prefix is stripped before adding them to the entity manager).
The following line shows an example of setting JPA properties for Hibernate: |

-
Properties
-
YAML

`spring.jpa.properties.hibernate.globally_quoted_identifiers=true````
spring:
 jpa:
 properties:
 hibernate:
 "globally_quoted_identifiers": "true"
```
The line in the preceding example passes a value of `true` for the `hibernate.globally_quoted_identifiers` property to the Hibernate entity manager.

By default, the DDL execution (or validation) is deferred until the `ApplicationContext` has started.

### Open EntityManager in View

If you are running a web application, Spring Boot by default registers `OpenEntityManagerInViewInterceptor` to apply the “Open EntityManager in View” pattern, to allow for lazy loading in web views.
If you do not want this behavior, you should set `spring.jpa.open-in-view` to `false` in your `application.properties`.

## Spring Data JDBC

Spring Data includes repository support for JDBC and will automatically generate SQL for the methods on `CrudRepository`.
For more advanced queries, a `@Query` annotation is provided.

Spring Boot will auto-configure Spring Data’s JDBC repositories when the necessary dependencies are on the classpath.
They can be added to your project with a single dependency on `spring-boot-starter-data-jdbc`.
If necessary, you can take control of Spring Data JDBC’s configuration by adding the `@EnableJdbcRepositories` annotation or an `AbstractJdbcConfiguration` subclass to your application.

If you’re using Spring Data JDBC with ahead-of-time processing (targeting either the JVM or a native image), some additional configuration is recommended.
To prevent the need for a DB connection during AOT processing, define a `JdbcDialect` bean that’s appropriate for your application’s database.
For example, if you’re using Postgres, define a `JdbcPostgresDialect` bean.

| For complete details of Spring Data JDBC, see the reference documentation. |

## Using H2’s Web Console

The H2 database provides a browser-based console that Spring Boot can auto-configure for you. The console is auto-configured when the following conditions are met:

-
You are developing a servlet-based web application.
-
`org.springframework.boot:spring-boot-h2console`is on the classpath.
-
You are using Spring Boot’s developer tools.

| If you are not using Spring Boot’s developer tools but would still like to make use of H2’s console, you can configure the `spring.h2.console.enabled`property with a value of`true`. |

| The H2 console is only intended for use during development, so you should take care to ensure that `spring.h2.console.enabled`is not set to`true`in production. |

### Changing the H2 Console’s Path

By default, the console is available at `/h2-console`.
You can customize the console’s path by using the `spring.h2.console.path` property.

### Accessing the H2 Console in a Secured Application

H2 Console uses frames and, as it is intended for development only, does not implement CSRF protection measures. If your application uses Spring Security, you need to configure it to

-
disable CSRF protection for requests against the console,
-
set the header `X-Frame-Options`to`SAMEORIGIN`on responses from the console.

More information on CSRF and the header X-Frame-Options can be found in the Spring Security Reference Guide.

In simple setups, a `SecurityFilterChain` like the following can be used:

-
Java
-
Kotlin

```
import org.springframework.boot.security.autoconfigure.web.servlet.PathRequest;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Profile;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.CsrfConfigurer;
import org.springframework.security.config.annotation.web.configurers.HeadersConfigurer.FrameOptionsConfig;
import org.springframework.security.web.SecurityFilterChain;
@Profile("dev")
@Configuration(proxyBeanMethods = false)
public class DevProfileSecurityConfiguration {
	@Bean
	@Order(Ordered.HIGHEST_PRECEDENCE)
	SecurityFilterChain h2ConsoleSecurityFilterChain(HttpSecurity http) {
 http.securityMatcher(PathRequest.toH2Console());
 http.authorizeHttpRequests(yourCustomAuthorization());
 http.csrf(CsrfConfigurer::disable);
 http.headers((headers) -> headers.frameOptions(FrameOptionsConfig::sameOrigin));
 return http.build();
	}
}
```
```
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.context.annotation.Profile
import org.springframework.core.Ordered
import org.springframework.core.annotation.Order
import org.springframework.security.config.Customizer
import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.web.SecurityFilterChain
@Profile("dev")
@Configuration(proxyBeanMethods = false)
class DevProfileSecurityConfiguration {
	@Bean
	@Order(Ordered.HIGHEST_PRECEDENCE)
	fun h2ConsoleSecurityFilterChain(http: HttpSecurity): SecurityFilterChain {
 return http.authorizeHttpRequests(yourCustomAuthorization())
 .csrf { csrf -> csrf.disable() }
 .headers { headers -> headers.frameOptions { frameOptions -> frameOptions.sameOrigin() } }
 .build()
	}
}
```
| The H2 console is only intended for use during development. In production, disabling CSRF protection or allowing frames for a website may create severe security risks. |

| `PathRequest.toH2Console()`returns the correct request matcher also when the console’s path has been customized. |

## Using jOOQ

jOOQ Object Oriented Querying (jOOQ) is a popular product from Data Geekery which generates Java code from your database and lets you build type-safe SQL queries through its fluent API. Both the commercial and open source editions can be used with Spring Boot. jOOQ requires Java 21 or later.

### Code Generation

In order to use jOOQ type-safe queries, you need to generate Java classes from your database schema.
You can follow the instructions in the jOOQ user manual.
If you use the `jooq-codegen-maven` plugin and you also use the `spring-boot-starter-parent` “parent POM”, you can safely omit the plugin’s `<version>` tag.
You can also use Spring Boot-defined version variables (such as `h2.version`) to declare the plugin’s database dependency.
The following listing shows an example:

```
<plugin>
	<groupId>org.jooq</groupId>
	<artifactId>jooq-codegen-maven</artifactId>
	<executions>
 ...
	</executions>
	<dependencies>
 <dependency>
 <groupId>com.h2database</groupId>
 <artifactId>h2</artifactId>
 <version>${h2.version}</version>
 </dependency>
	</dependencies>
	<configuration>
 <jdbc>
 <driver>org.h2.Driver</driver>
 <url>jdbc:h2:~/yourdatabase</url>
 </jdbc>
 <generator>
 ...
 </generator>
	</configuration>
</plugin>
```
### Using DSLContext

The fluent API offered by jOOQ is initiated through the `DSLContext` interface.
Spring Boot auto-configures a `DSLContext` as a Spring Bean and connects it to your application `DataSource`.
To use the `DSLContext`, you can inject it, as shown in the following example:

-
Java
-
Kotlin

```
import java.util.GregorianCalendar;
import java.util.List;
import org.jooq.DSLContext;
import org.springframework.stereotype.Component;
import static org.springframework.boot.docs.data.sql.jooq.dslcontext.Tables.AUTHOR;
@Component
public class MyBean {
	private final DSLContext create;
	public MyBean(DSLContext dslContext) {
 this.create = dslContext;
	}
}
```
```
import org.jooq.DSLContext
import org.springframework.stereotype.Component
import java.util.GregorianCalendar
@Component
class MyBean(private val create: DSLContext) {
}
```
| The jOOQ manual tends to use a variable named `create`to hold the`DSLContext`. |

You can then use the `DSLContext` to construct your queries, as shown in the following example:

-
Java
-
Kotlin

```
	public List<GregorianCalendar> authorsBornAfter1980() {
 return this.create.selectFrom(AUTHOR)
 .where(AUTHOR.DATE_OF_BIRTH.greaterThan(new GregorianCalendar(1980, 0, 1)))
 .fetch(AUTHOR.DATE_OF_BIRTH);
```
```
	fun authorsBornAfter1980(): List<GregorianCalendar> {
 return create.selectFrom<Tables.TAuthorRecord>(Tables.AUTHOR)
 .where(Tables.AUTHOR?.DATE_OF_BIRTH?.greaterThan(GregorianCalendar(1980, 0, 1)))
 .fetch(Tables.AUTHOR?.DATE_OF_BIRTH)
	}
```
### jOOQ SQL Dialect

Unless the `spring.jooq.sql-dialect` property has been configured, Spring Boot determines the SQL dialect to use for your datasource.
If Spring Boot could not detect the dialect, it uses `DEFAULT`.

| Spring Boot can only auto-configure dialects supported by the open source version of jOOQ. |

### Customizing jOOQ

More advanced customizations can be achieved by defining your own `DefaultConfigurationCustomizer` bean that will be invoked prior to creating the `Configuration` `@Bean`.
This takes precedence to anything that is applied by the auto-configuration.

You can also create your own `Configuration` `@Bean` if you want to take complete control of the jOOQ configuration.

## Using R2DBC

The Reactive Relational Database Connectivity (R2DBC) project brings reactive programming APIs to relational databases.
R2DBC’s `Connection` provides a standard method of working with non-blocking database connections.
Connections are provided by using a `ConnectionFactory`, similar to a `DataSource` with jdbc.

`ConnectionFactory` configuration is controlled by external configuration properties in `spring.r2dbc.*`.
For example, you might declare the following section in `application.properties`:

-
Properties
-
YAML

```
spring.r2dbc.url=r2dbc:postgresql://localhost/test
spring.r2dbc.username=dbuser
spring.r2dbc.password=dbpass
```
```
spring:
 r2dbc:
 url: "r2dbc:postgresql://localhost/test"
 username: "dbuser"
 password: "dbpass"
```
| You do not need to specify a driver class name, since Spring Boot obtains the driver from R2DBC’s Connection Factory discovery. |

| At least the url should be provided.
Information specified in the URL takes precedence over individual properties, that is `name`,`username`,`password`and pooling options. |

| The “How-to Guides” section includes a section on how to initialize a database. |

To customize the connections created by a `ConnectionFactory`, that is, set specific parameters that you do not want (or cannot) configure in your central database configuration, you can use a `ConnectionFactoryOptionsBuilderCustomizer` `@Bean`.
The following example shows how to manually override the database port while the rest of the options are taken from the application configuration:

-
Java
-
Kotlin

```
import io.r2dbc.spi.ConnectionFactoryOptions;
import org.springframework.boot.r2dbc.autoconfigure.ConnectionFactoryOptionsBuilderCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyR2dbcConfiguration {
	@Bean
	public ConnectionFactoryOptionsBuilderCustomizer connectionFactoryPortCustomizer() {
 return (builder) -> builder.option(ConnectionFactoryOptions.PORT, 5432);
	}
}
```
```
import io.r2dbc.spi.ConnectionFactoryOptions
import org.springframework.boot.r2dbc.autoconfigure.ConnectionFactoryOptionsBuilderCustomizer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
@Configuration(proxyBeanMethods = false)
class MyR2dbcConfiguration {
	@Bean
	fun connectionFactoryPortCustomizer(): ConnectionFactoryOptionsBuilderCustomizer {
 return ConnectionFactoryOptionsBuilderCustomizer { builder ->
 builder.option(ConnectionFactoryOptions.PORT, 5432)
 }
	}
}
```
The following examples show how to set some PostgreSQL connection options:

-
Java
-
Kotlin

```
import java.util.HashMap;
import java.util.Map;
import io.r2dbc.postgresql.PostgresqlConnectionFactoryProvider;
import org.springframework.boot.r2dbc.autoconfigure.ConnectionFactoryOptionsBuilderCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyPostgresR2dbcConfiguration {
	@Bean
	public ConnectionFactoryOptionsBuilderCustomizer postgresCustomizer() {
 Map<String, String> options = new HashMap<>();
 options.put("lock_timeout", "30s");
 options.put("statement_timeout", "60s");
 return (builder) -> builder.option(PostgresqlConnectionFactoryProvider.OPTIONS, options);
	}
}
```
```
import io.r2dbc.postgresql.PostgresqlConnectionFactoryProvider
import org.springframework.boot.r2dbc.autoconfigure.ConnectionFactoryOptionsBuilderCustomizer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
@Configuration(proxyBeanMethods = false)
class MyPostgresR2dbcConfiguration {
	@Bean
	fun postgresCustomizer(): ConnectionFactoryOptionsBuilderCustomizer {
 val options: MutableMap<String, String> = HashMap()
 options["lock_timeout"] = "30s"
 options["statement_timeout"] = "60s"
 return ConnectionFactoryOptionsBuilderCustomizer { builder ->
 builder.option(PostgresqlConnectionFactoryProvider.OPTIONS, options)
 }
	}
}
```
When a `ConnectionFactory` bean is available, the regular JDBC `DataSource` auto-configuration backs off.
If you want to retain the JDBC `DataSource` auto-configuration, and are comfortable with the risk of using the blocking JDBC API in a reactive application, add `@Import(DataSourceAutoConfiguration.class)` on a `@Configuration` class in your application to re-enable it.

### Embedded Database Support

Similarly to the JDBC support, Spring Boot can automatically configure an embedded database for reactive usage. You need not provide any connection URLs. You need only include a build dependency to the embedded database that you want to use, as shown in the following example:

```
<dependency>
	<groupId>io.r2dbc</groupId>
	<artifactId>r2dbc-h2</artifactId>
	<scope>runtime</scope>
</dependency>
```
| If you are using this feature in your tests, you may notice that the same database is reused by your whole test suite regardless of the number of application contexts that you use.
If you want to make sure that each context has a separate embedded database, you should set |

### Using DatabaseClient

A `DatabaseClient` bean is auto-configured, and you can autowire it directly into your own beans, as shown in the following example:

-
Java
-
Kotlin

```
import java.util.Map;
import reactor.core.publisher.Flux;
import org.springframework.r2dbc.core.DatabaseClient;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final DatabaseClient databaseClient;
	public MyBean(DatabaseClient databaseClient) {
 this.databaseClient = databaseClient;
	}
	// ...
	public Flux<Map<String, Object>> someMethod() {
 return this.databaseClient.sql("select * from user").fetch().all();
	}
}
```
```
import org.springframework.r2dbc.core.DatabaseClient
import org.springframework.stereotype.Component
import reactor.core.publisher.Flux
@Component
class MyBean(private val databaseClient: DatabaseClient) {
	// ...
	fun someMethod(): Flux<Map<String, Any>> {
 return databaseClient.sql("select * from user").fetch().all()
	}
}
```
### Spring Data R2DBC Repositories

Spring Data R2DBC repositories are interfaces that you can define to access data.
Queries are created automatically from your method names.
For example, a `CityRepository` interface might declare a `findAllByState(String state)` method to find all the cities in a given state.

For more complex queries, you can annotate your method with Spring Data’s `@Query` annotation.

Spring Data repositories usually extend from the `Repository` or `CrudRepository` interfaces.
If you use auto-configuration, the auto-configuration packages are searched for repositories.

The following example shows a typical Spring Data repository interface definition:

-
Java
-
Kotlin

```
import reactor.core.publisher.Mono;
import org.springframework.data.repository.Repository;
public interface CityRepository extends Repository<City, Long> {
	Mono<City> findByNameAndStateAllIgnoringCase(String name, String state);
}
```
```
import org.springframework.data.repository.Repository
import reactor.core.publisher.Mono
interface CityRepository : Repository<City, Long> {
	fun findByNameAndStateAllIgnoringCase(name: String, state: String): Mono<City>
}
```
| We have barely scratched the surface of Spring Data R2DBC. For complete details, see the Spring Data R2DBC reference documentation. |

# Working with NoSQL Technologies

Spring Data provides additional projects that help you access a variety of NoSQL technologies, including:

Of these, Spring Boot provides auto-configuration for Cassandra, Couchbase, Elasticsearch, LDAP, MongoDB, Neo4J and Redis. Additionally, Spring Boot for Apache Geode provides auto-configuration for Apache Geode. You can make use of the other projects, but you must configure them yourself. See the appropriate reference documentation at spring.io/projects/spring-data.

Spring Boot also provides auto-configuration for the InfluxDB client but it is deprecated in favor of the new InfluxDB Java client that provides its own Spring Boot integration.

## Redis

Redis is a cache, message broker, and richly-featured key-value store. Spring Boot offers basic auto-configuration for the Lettuce and Jedis client libraries and the abstractions on top of them provided by Spring Data Redis.

There is a `spring-boot-starter-data-redis` starter for collecting the dependencies in a convenient way.
By default, it uses Lettuce.
That starter handles both traditional and reactive applications.

| We also provide a `spring-boot-starter-data-redis-reactive`starter for consistency with the other stores with reactive support. |

### Connecting to Redis

You can inject an auto-configured `RedisConnectionFactory`, `StringRedisTemplate`, or vanilla `RedisTemplate` instance as you would any other Spring Bean.
The following listing shows an example of such a bean:

-
Java
-
Kotlin

```
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final StringRedisTemplate template;
	public MyBean(StringRedisTemplate template) {
 this.template = template;
	}
	// ...
	public Boolean someMethod() {
 return this.template.hasKey("spring");
	}
}
```
```
import org.springframework.data.redis.core.StringRedisTemplate
import org.springframework.stereotype.Component
@Component
class MyBean(private val template: StringRedisTemplate) {
	// ...
	fun someMethod(): Boolean {
 return template.hasKey("spring")
	}
}
```
By default, the instance tries to connect to a Redis server at `localhost:6379`.
You can specify custom connection details using `spring.data.redis.*` properties, as shown in the following example:

-
Properties
-
YAML

```
spring.data.redis.host=localhost
spring.data.redis.port=6379
spring.data.redis.database=0
spring.data.redis.username=user
spring.data.redis.password=secret
```
```
spring:
 data:
 redis:
 host: "localhost"
 port: 6379
 database: 0
 username: "user"
 password: "secret"
```
You can also specify the url of the Redis server directly. When setting the url, the host, port, username and password properties are ignored. This is shown in the following example:

-
Properties
-
YAML

```
spring.data.redis.url=redis://user:secret@localhost:6379
spring.data.redis.database=0
```
```
spring:
 data:
 redis:
 url: "redis://user:secret@localhost:6379"
 database: 0
```
| You can also register an arbitrary number of beans that implement `LettuceClientConfigurationBuilderCustomizer`for more advanced customizations.`ClientResources`can also be customized using`ClientResourcesBuilderCustomizer`.
If you use Jedis,`JedisClientConfigurationBuilderCustomizer`is also available. |

Alternatively, you can register a bean of type `RedisStandaloneConfiguration`, `RedisSentinelConfiguration`, `RedisClusterConfiguration`, or `RedisStaticMasterReplicaConfiguration` to take full control over the configuration.

| master/replica is not supported by Jedis. |

If you add your own `@Bean` of any of the auto-configured types, it replaces the default (except in the case of `RedisTemplate`, when the exclusion is based on the bean name, `redisTemplate`, not its type).

By default, a pooled connection factory is auto-configured if `commons-pool2` is on the classpath.

The auto-configured `RedisConnectionFactory` can be configured to use SSL for communication with the server by setting the properties as shown in this example:

-
Properties
-
YAML

`spring.data.redis.ssl.enabled=true````
spring:
 data:
 redis:
 ssl:
 enabled: true
```
Custom SSL trust material can be configured in an SSL bundle and applied to the `RedisConnectionFactory` as shown in this example:

-
Properties
-
YAML

`spring.data.redis.ssl.bundle=example````
spring:
 data:
 redis:
 ssl:
 bundle: "example"
```
### Receiving a Message

When the Redis infrastructure is present, any bean can be annotated with `@RedisListener` to create a listener endpoint.
If no `RedisMessageListenerContainer` has been defined, a default one is configured automatically.

The following component creates a listener endpoint on the `someChannel` channel:

-
Java
-
Kotlin

```
import org.springframework.data.redis.annotation.RedisListener;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	@RedisListener("someChannel")
	public void processMessage(String content) {
 // ...
	}
}
```
```
import org.springframework.data.redis.annotation.RedisListener
import org.springframework.stereotype.Component
@Component
class MyBean {
	@RedisListener("someChannel")
	fun processMessage(content: String) {
 // ...
	}
}
```
| See the `@EnableRedisListeners`API documentation for more details. |

If you need to create more `RedisMessageListenerContainer` instances or if you want to override the default, Spring Boot provides a `RedisMessageListenerContainerConfigurer` that you can use to initialize a `RedisMessageListenerContainer` with the same settings as the one that is auto-configured.

For instance, the following example exposes another container that uses a specific `RedisConnectionFactory`:

-
Java
-
Kotlin

```
import org.springframework.boot.data.redis.autoconfigure.RedisMessageListenerContainerConfigurer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.redis.connection.RedisConnectionFactory;
import org.springframework.data.redis.listener.RedisMessageListenerContainer;
@Configuration(proxyBeanMethods = false)
public class MyRedisConfiguration {
	@Bean
	public RedisMessageListenerContainer myRedisMessageListenerContainer(
 RedisMessageListenerContainerConfigurer configurer, RedisConnectionFactory connectionFactory) {
 RedisMessageListenerContainer container = new RedisMessageListenerContainer();
 configurer.configure(container, connectionFactory);
 // ... custom configuration
 return container;
	}
}
```
```
import org.springframework.boot.data.redis.autoconfigure.RedisMessageListenerContainerConfigurer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.data.redis.connection.RedisConnectionFactory
import org.springframework.data.redis.listener.RedisMessageListenerContainer
@Configuration(proxyBeanMethods = false)
class MyRedisConfiguration {
	@Bean
	fun myRedisMessageListenerContainer(
 configurer: RedisMessageListenerContainerConfigurer,
 connectionFactory: RedisConnectionFactory
	): RedisMessageListenerContainer {
 val container = RedisMessageListenerContainer()
 configurer.configure(container, connectionFactory)
 // ... custom configuration
 return container
	}
}
```
Then you can use the container in any `@RedisListener`-annotated method as follows:

-
Java
-
Kotlin

```
import org.springframework.data.redis.annotation.RedisListener;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	@RedisListener(topic = "someChannel", container = "myRedisMessageListenerContainer")
	public void processMessage(String content) {
 // ...
	}
}
```
```
import org.springframework.data.redis.annotation.RedisListener
import org.springframework.stereotype.Component
@Component
class MyBean {
	@RedisListener(topic = "someChannel", container = "myRedisMessageListenerContainer")
	fun processMessage(content: String) {
 // ...
	}
}
```
## MongoDB

MongoDB is an open-source NoSQL document database that uses a JSON-like schema instead of traditional table-based relational data.
Spring Boot offers several conveniences for working with MongoDB, including the `spring-boot-starter-data-mongodb` and `spring-boot-starter-data-mongodb-reactive` starters.

### Connecting to a MongoDB Database

To access MongoDB databases, you can inject an auto-configured `MongoDatabaseFactory`.
By default, the instance tries to connect to a MongoDB server at `mongodb://localhost/test`.
The following example shows how to connect to a MongoDB database:

-
Java
-
Kotlin

```
import com.mongodb.client.MongoCollection;
import com.mongodb.client.MongoDatabase;
import org.bson.Document;
import org.springframework.data.mongodb.MongoDatabaseFactory;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final MongoDatabaseFactory mongo;
	public MyBean(MongoDatabaseFactory mongo) {
 this.mongo = mongo;
	}
	// ...
	public MongoCollection<Document> someMethod() {
 MongoDatabase db = this.mongo.getMongoDatabase();
 return db.getCollection("users");
	}
}
```
```
import com.mongodb.client.MongoCollection
import org.bson.Document
import org.springframework.data.mongodb.MongoDatabaseFactory
import org.springframework.stereotype.Component
@Component
class MyBean(private val mongo: MongoDatabaseFactory) {
	// ...
	fun someMethod(): MongoCollection<Document> {
 val db = mongo.mongoDatabase
 return db.getCollection("users")
	}
}
```
If you have defined your own `MongoClient`, it will be used to auto-configure a suitable `MongoDatabaseFactory`.

The auto-configured `MongoClient` is created using a `MongoClientSettings` bean.
If you have defined your own `MongoClientSettings`, it will be used without modification and the `spring.data.mongodb` properties will be ignored.
Otherwise a `MongoClientSettings` will be auto-configured and will have the `spring.data.mongodb` properties applied to it.
In either case, you can declare one or more `MongoClientSettingsBuilderCustomizer` beans to fine-tune the `MongoClientSettings` configuration.
Each will be called in order with the `MongoClientSettings.Builder` that is used to build the `MongoClientSettings`.

You can set the `spring.mongodb.uri` property to change the URL and configure additional settings such as the *replica set*, as shown in the following example:

-
Properties
-
YAML

`spring.mongodb.uri=mongodb://user:[email protected]:27017,mongoserver2.example.com:23456/test````
spring:
 mongodb:
 uri: "mongodb://user:[email protected]:27017,mongoserver2.example.com:23456/test"
```
Alternatively, you can specify connection details using discrete properties.
For example, you might declare the following settings in your `application.properties`:

-
Properties
-
YAML

```
spring.mongodb.host=mongoserver1.example.com
spring.mongodb.port=27017
spring.mongodb.additional-hosts[0]=mongoserver2.example.com:23456
spring.mongodb.database=test
spring.mongodb.username=user
spring.mongodb.password=secret
```
```
spring:
 mongodb:
 host: "mongoserver1.example.com"
 port: 27017
 additional-hosts:
 - "mongoserver2.example.com:23456"
 database: "test"
 username: "user"
 password: "secret"
```
The auto-configured `MongoClient` can be configured to use SSL for communication with the server by setting the properties as shown in this example:

-
Properties
-
YAML

```
spring.mongodb.uri=mongodb://user:[email protected]:27017,mongoserver2.example.com:23456/test
spring.mongodb.ssl.enabled=true
```
```
spring:
 mongodb:
 uri: "mongodb://user:[email protected]:27017,mongoserver2.example.com:23456/test"
 ssl:
 enabled: true
```
Custom SSL trust material can be configured in an SSL bundle and applied to the `MongoClient` as shown in this example:

-
Properties
-
YAML

```
spring.mongodb.uri=mongodb://user:[email protected]:27017,mongoserver2.example.com:23456/test
spring.mongodb.ssl.bundle=example
```
```
spring:
 mongodb:
 uri: "mongodb://user:[email protected]:27017,mongoserver2.example.com:23456/test"
 ssl:
 bundle: "example"
```
| If You can also specify the port as part of the host address by using the |

| If you do not use Spring Data MongoDB, you can inject a `MongoClient`bean instead of using`MongoDatabaseFactory`.
If you want to take complete control of establishing the MongoDB connection, you can also declare your own`MongoDatabaseFactory`or`MongoClient`bean. |

| If you are using the reactive driver, Netty is required for SSL. The auto-configuration configures this factory automatically if Netty is available and the factory to use has not been customized already. |

### MongoTemplate

Spring Data MongoDB provides a `MongoTemplate` class that is very similar in its design to Spring’s `JdbcTemplate`.
As with `JdbcTemplate`, Spring Boot auto-configures a bean for you to inject the template, as follows:

-
Java
-
Kotlin

```
import com.mongodb.client.MongoCollection;
import org.bson.Document;
import org.springframework.data.mongodb.core.MongoTemplate;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final MongoTemplate mongoTemplate;
	public MyBean(MongoTemplate mongoTemplate) {
 this.mongoTemplate = mongoTemplate;
	}
	// ...
	public MongoCollection<Document> someMethod() {
 return this.mongoTemplate.getCollection("users");
	}
}
```
```
import com.mongodb.client.MongoCollection
import org.bson.Document
import org.springframework.data.mongodb.core.MongoTemplate
import org.springframework.stereotype.Component
@Component
class MyBean(private val mongoTemplate: MongoTemplate) {
	// ...
	fun someMethod(): MongoCollection<Document> {
 return mongoTemplate.getCollection("users")
	}
}
```
See the `MongoOperations` API documentation for complete details.

### Spring Data MongoDB Repositories

Spring Data includes repository support for MongoDB. As with the JPA repositories discussed earlier, the basic principle is that queries are constructed automatically, based on method names.

In fact, both Spring Data JPA and Spring Data MongoDB share the same common infrastructure.
You could take the JPA example from earlier and, assuming that `City` is now a MongoDB data class rather than a JPA `@Entity`, it works in the same way, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.repository.Repository;
public interface CityRepository extends Repository<City, Long> {
	Page<City> findAll(Pageable pageable);
	City findByNameAndStateAllIgnoringCase(String name, String state);
}
```
```
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.Repository
interface CityRepository :
	Repository<City, Long> {
	fun findAll(pageable: Pageable?): Page<City>
	fun findByNameAndStateAllIgnoringCase(name: String, state: String): City?
}
```
Repositories and documents are found through scanning.
By default, the auto-configuration packages are scanned.
You can customize the locations to look for repositories and documents by using `@EnableMongoRepositories` and `@EntityScan` respectively.

| For complete details of Spring Data MongoDB, including its rich object mapping technologies, see its reference documentation. |

## Neo4j

Neo4j is an open-source NoSQL graph database that uses a rich data model of nodes connected by first class relationships, which is better suited for connected big data than traditional RDBMS approaches.
Spring Boot offers several conveniences for working with Neo4j, including the `spring-boot-starter-data-neo4j` starter.

### Connecting to a Neo4j Database

To access a Neo4j server, you can inject an auto-configured `Driver`.
By default, the instance tries to connect to a Neo4j server at `localhost:7687` using the Bolt protocol.
The following example shows how to inject a Neo4j `Driver` that gives you access, amongst other things, to a `Session`:

-
Java
-
Kotlin

```
import org.neo4j.driver.Driver;
import org.neo4j.driver.Session;
import org.neo4j.driver.Values;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final Driver driver;
	public MyBean(Driver driver) {
 this.driver = driver;
	}
	// ...
	public String someMethod(String message) {
 try (Session session = this.driver.session()) {
 return session.executeWrite(
 (transaction) -> transaction
 .run("CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
 Values.parameters("message", message))
 .single()
 .get(0)
 .asString());
 }
	}
}
```
```
import org.neo4j.driver.Driver
import org.neo4j.driver.TransactionContext
import org.neo4j.driver.Values
import org.springframework.stereotype.Component
@Component
class MyBean(private val driver: Driver) {
	// ...
	fun someMethod(message: String?): String {
 driver.session().use { session ->
 return@someMethod session.executeWrite { transaction: TransactionContext ->
 transaction
 .run(
 "CREATE (a:Greeting) SET a.message = \$message RETURN a.message + ', from node ' + id(a)",
 Values.parameters("message", message)
 )
 .single()[0].asString()
 }
 }
	}
}
```
You can configure various aspects of the driver using `spring.neo4j.*` properties.
The following example shows how to configure the uri and credentials to use:

-
Properties
-
YAML

```
spring.neo4j.uri=bolt://my-server:7687
spring.neo4j.authentication.username=neo4j
spring.neo4j.authentication.password=secret
```
```
spring:
 neo4j:
 uri: "bolt://my-server:7687"
 authentication:
 username: "neo4j"
 password: "secret"
```
The auto-configured `Driver` is created using `org.neo4j.driver.Config$ConfigBuilder`.
To fine-tune its configuration, declare one or more `ConfigBuilderCustomizer` beans.
Each will be called in order with the `org.neo4j.driver.Config$ConfigBuilder` that is used to build the `Driver`.

### Spring Data Neo4j Repositories

Spring Data includes repository support for Neo4j. For complete details of Spring Data Neo4j, see the reference documentation.

Spring Data Neo4j shares the common infrastructure with Spring Data JPA as many other Spring Data modules do.
You could take the JPA example from earlier and define `City` as Spring Data Neo4j `@Node` rather than JPA `@Entity` and the repository abstraction works in the same way, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.data.neo4j.repository.Neo4jRepository;
public interface CityRepository extends Neo4jRepository<City, Long> {
	City findOneByNameAndState(String name, String state);
}
```
```
import org.springframework.data.neo4j.repository.Neo4jRepository
interface CityRepository : Neo4jRepository<City, Long> {
	fun findOneByNameAndState(name: String?, state: String?): City?
}
```
The `spring-boot-starter-data-neo4j` starter enables the repository support as well as transaction management.
Spring Boot supports both classic and reactive Neo4j repositories, using the `Neo4jTemplate` or `ReactiveNeo4jTemplate` beans.
When Project Reactor is available on the classpath, the reactive style is also auto-configured.

Repositories and entities are found through scanning.
By default, the auto-configuration packages are scanned.
You can customize the locations to look for repositories and entities by using `@EnableNeo4jRepositories` and `@EntityScan` respectively.

| In an application using the reactive style, a
 |

## Elasticsearch

Elasticsearch is an open source, distributed, RESTful search and analytics engine. Spring Boot offers basic auto-configuration for Elasticsearch clients.

Spring Boot supports several clients:

-
The official low-level REST client
-
The official Java API client
-
The `ReactiveElasticsearchClient`provided by Spring Data Elasticsearch

Spring Boot provides a dedicated starter, `spring-boot-starter-data-elasticsearch`.

### Connecting to Elasticsearch Using REST clients

Elasticsearch ships two different REST clients that you can use to query a cluster: the low-level client and the Java API client.
The Java API client is provided by the `co.elastic.clients:elasticsearch-java` module and
the low-level client is provided by the `co.elastic.clients:elasticsearch-rest5-client` module.
Additionally, Spring Boot provides support for a reactive client from the `org.springframework.data:spring-data-elasticsearch` module.
By default, the clients will target `localhost:9200`.
You can use `spring.elasticsearch.*` properties to further tune how the clients are configured, as shown in the following example:

-
Properties
-
YAML

```
spring.elasticsearch.uris=https://search.example.com:9200
spring.elasticsearch.socket-timeout=10s
spring.elasticsearch.username=user
spring.elasticsearch.password=secret
```
```
spring:
 elasticsearch:
 uris: "https://search.example.com:9200"
 socket-timeout: "10s"
 username: "user"
 password: "secret"
```
#### Connecting to Elasticsearch Using Rest5Client

If you have `co.elastic.clients:elasticsearch-rest5-client` on the classpath, Spring Boot will auto-configure and register a `Rest5Client` bean.
In addition to the properties described previously, to fine-tune the `Rest5Client` you can register an arbitrary number of beans that implement `Rest5ClientBuilderCustomizer` for more advanced customizations.
To take full control over the client’s configuration, define a `Rest5ClientBuilder` bean.

Additionally, a `Sniffer` can be auto-configured to automatically discover nodes from a running Elasticsearch cluster and set them on the `Rest5Client` bean.
You can further tune how `Sniffer` is configured, as shown in the following example:

-
Properties
-
YAML

```
spring.elasticsearch.restclient.sniffer.enabled=true
spring.elasticsearch.restclient.sniffer.interval=10m
spring.elasticsearch.restclient.sniffer.delay-after-failure=30s
```
```
spring:
 elasticsearch:
 restclient:
 sniffer:
 enabled: true
 interval: "10m"
 delay-after-failure: "30s"
```
#### Connecting to Elasticsearch Using ElasticsearchClient

If you use the `spring-boot-starter-elasticsearch` or have added `co.elastic.clients:elasticsearch-java` to the classpath, Spring Boot will auto-configure and register an `ElasticsearchClient` bean.

The `ElasticsearchClient` uses a transport that depends upon the previously described `Rest5Client`.
Therefore, the properties described previously can be used to configure the `ElasticsearchClient`.
Furthermore, you can define a `Rest5ClientOptions` bean to take further control of the behavior of the transport.

#### Connecting to Elasticsearch using ReactiveElasticsearchClient

Spring Data Elasticsearch ships `ReactiveElasticsearchClient` for querying Elasticsearch instances in a reactive fashion.
If you have Spring Data Elasticsearch and Reactor on the classpath, Spring Boot will auto-configure and register a `ReactiveElasticsearchClient`.

The `ReactiveElasticsearchClient` uses a transport that depends upon the previously described `Rest5Client`.
Therefore, the properties described previously can be used to configure the `ReactiveElasticsearchClient`.
Furthermore, you can define a `Rest5ClientOptions` bean to take further control of the behavior of the transport.

### Connecting to Elasticsearch by Using Spring Data

To connect to Elasticsearch, an `ElasticsearchClient` bean must be defined,
auto-configured by Spring Boot or manually provided by the application (see previous sections).
With this configuration in place, an
`ElasticsearchTemplate` can be injected like any other Spring bean,
as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.data.elasticsearch.client.elc.ElasticsearchTemplate;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final ElasticsearchTemplate template;
	public MyBean(ElasticsearchTemplate template) {
 this.template = template;
	}
	// ...
	public boolean someMethod(String id) {
 return this.template.exists(id, User.class);
	}
}
```
```
import org.springframework.stereotype.Component
@Component
class MyBean(private val template: org.springframework.data.elasticsearch.client.elc.ElasticsearchTemplate ) {
	// ...
	fun someMethod(id: String): Boolean {
 return template.exists(id, User::class.java)
	}
}
```
In the presence of `spring-data-elasticsearch` and Reactor, Spring Boot can also auto-configure a `ReactiveElasticsearchClient` and a `ReactiveElasticsearchTemplate` as beans.
They are the reactive equivalent of the other REST clients.

### Spring Data Elasticsearch Repositories

Spring Data includes repository support for Elasticsearch. As with the JPA repositories discussed earlier, the basic principle is that queries are constructed for you automatically based on method names.

In fact, both Spring Data JPA and Spring Data Elasticsearch share the same common infrastructure.
You could take the JPA example from earlier and, assuming that `City` is now an Elasticsearch `@Document` class rather than a JPA `@Entity`, it works in the same way.

Repositories and documents are found through scanning.
By default, the auto-configuration packages are scanned.
You can customize the locations to look for repositories and documents by using `@EnableElasticsearchRepositories` and `@EntityScan` respectively.

| For complete details of Spring Data Elasticsearch, see the reference documentation. |

Spring Boot supports both classic and reactive Elasticsearch repositories, using the `ElasticsearchTemplate` or `ReactiveElasticsearchTemplate` beans.
Most likely those beans are auto-configured by Spring Boot given the required dependencies are present.

If you wish to use your own template for backing the Elasticsearch repositories, you can add your own `ElasticsearchTemplate` or `ElasticsearchOperations` `@Bean`, as long as it is named `"elasticsearchTemplate"`.
Same applies to `ReactiveElasticsearchTemplate` and `ReactiveElasticsearchOperations`, with the bean name `"reactiveElasticsearchTemplate"`.

You can choose to disable the repositories support with the following property:

-
Properties
-
YAML

`spring.data.elasticsearch.repositories.enabled=false````
spring:
 data:
 elasticsearch:
 repositories:
 enabled: false
```
## Cassandra

Cassandra is an open source, distributed database management system designed to handle large amounts of data across many commodity servers.
Spring Boot offers auto-configuration for Cassandra and the abstractions on top of it provided by Spring Data Cassandra.
There is a `spring-boot-starter-data-cassandra` starter for collecting the dependencies in a convenient way.

### Connecting to Cassandra

You can inject an auto-configured `CqlTemplate`, `CassandraTemplate`, or a Cassandra `CqlSession` instance as you would with any other Spring Bean.
The `spring.cassandra.*` properties can be used to customize the connection.
Generally, you provide `keyspace-name` and `contact-points` as well the local datacenter name, as shown in the following example:

-
Properties
-
YAML

```
spring.cassandra.keyspace-name=mykeyspace
spring.cassandra.contact-points=cassandrahost1:9042,cassandrahost2:9042
spring.cassandra.local-datacenter=datacenter1
```
```
spring:
 cassandra:
 keyspace-name: "mykeyspace"
 contact-points: "cassandrahost1:9042,cassandrahost2:9042"
 local-datacenter: "datacenter1"
```
If the port is the same for all your contact points you can use a shortcut and only specify the host names, as shown in the following example:

-
Properties
-
YAML

```
spring.cassandra.keyspace-name=mykeyspace
spring.cassandra.contact-points=cassandrahost1,cassandrahost2
spring.cassandra.local-datacenter=datacenter1
```
```
spring:
 cassandra:
 keyspace-name: "mykeyspace"
 contact-points: "cassandrahost1,cassandrahost2"
 local-datacenter: "datacenter1"
```
| Those two examples are identical as the port default to `9042`.
If you need to configure the port, use`spring.cassandra.port`. |

The auto-configured `CqlSession` can be configured to use SSL for communication with the server by setting the properties as shown in this example:

-
Properties
-
YAML

```
spring.cassandra.keyspace-name=mykeyspace
spring.cassandra.contact-points=cassandrahost1,cassandrahost2
spring.cassandra.local-datacenter=datacenter1
spring.cassandra.ssl.enabled=true
```
```
spring:
 cassandra:
 keyspace-name: "mykeyspace"
 contact-points: "cassandrahost1,cassandrahost2"
 local-datacenter: "datacenter1"
 ssl:
 enabled: true
```
Custom SSL trust material can be configured in an SSL bundle and applied to the `CqlSession` as shown in this example:

-
Properties
-
YAML

```
spring.cassandra.keyspace-name=mykeyspace
spring.cassandra.contact-points=cassandrahost1,cassandrahost2
spring.cassandra.local-datacenter=datacenter1
spring.cassandra.ssl.bundle=example
```
```
spring:
 cassandra:
 keyspace-name: "mykeyspace"
 contact-points: "cassandrahost1,cassandrahost2"
 local-datacenter: "datacenter1"
 ssl:
 bundle: "example"
```
| The Cassandra driver has its own configuration infrastructure that loads an Spring Boot does not look for such a file by default but can load one using For more advanced driver customizations, you can register an arbitrary number of beans that implement |

| If you use `CqlSessionBuilder`to create multiple`CqlSession`beans, keep in mind the builder is mutable so make sure to inject a fresh copy for each session. |

The following code listing shows how to inject a Cassandra bean:

-
Java
-
Kotlin

```
import org.springframework.data.cassandra.core.CassandraTemplate;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final CassandraTemplate template;
	public MyBean(CassandraTemplate template) {
 this.template = template;
	}
	// ...
	public long someMethod() {
 return this.template.count(User.class);
	}
}
```
```
import org.springframework.data.cassandra.core.CassandraTemplate
import org.springframework.stereotype.Component
@Component
class MyBean(private val template: CassandraTemplate) {
	// ...
	fun someMethod(): Long {
 return template.count(User::class.java)
	}
}
```
If you add your own `@Bean` of type `CassandraTemplate`, it replaces the default.

### Spring Data Cassandra Repositories

Spring Data includes basic repository support for Cassandra.
Currently, this is more limited than the JPA repositories discussed earlier and needs `@Query` annotated finder methods.

Repositories and entities are found through scanning.
By default, the auto-configuration packages are scanned.
You can customize the locations to look for repositories and entities by using `@EnableCassandraRepositories` and `@EntityScan` respectively.

| For complete details of Spring Data Cassandra, see the reference documentation. |

## Couchbase

Couchbase is an open-source, distributed, multi-model NoSQL document-oriented database that is optimized for interactive applications.
Spring Boot offers auto-configuration for Couchbase and the abstractions on top of it provided by Spring Data Couchbase.
There are `spring-boot-starter-data-couchbase` and `spring-boot-starter-data-couchbase-reactive` starters for collecting the dependencies in a convenient way.

### Connecting to Couchbase

You can get a `Cluster` by adding the Couchbase SDK and some configuration.
The `spring.couchbase.*` properties can be used to customize the connection.
Generally, you provide the connection string and credentials for authentication. Basic authentication with username and password can be configured as shown in the following example:

-
Properties
-
YAML

```
spring.couchbase.connection-string=couchbase://192.168.1.123
spring.couchbase.username=user
spring.couchbase.password=secret
```
```
spring:
 couchbase:
 connection-string: "couchbase://192.168.1.123"
 username: "user"
 password: "secret"
```
Client certificates can be used for authentication instead of username and password. The location and password for a Java KeyStore containing client certificates can be configured as shown in the following example:

-
Properties
-
YAML

```
spring.couchbase.connection-string=couchbase://192.168.1.123
spring.couchbase.env.ssl.enabled=true
spring.couchbase.authentication.jks.location=classpath:client.p12
spring.couchbase.authentication.jks.password=secret
```
```
spring:
 couchbase:
 connection-string: "couchbase://192.168.1.123"
 env:
 ssl:
 enabled: true
 authentication:
 jks:
 location: "classpath:client.p12"
 password: "secret"
```
PEM-encoded certificates and a private key can be configured as shown in the following example:

-
Properties
-
YAML

```
spring.couchbase.connection-string=couchbase://192.168.1.123
spring.couchbase.env.ssl.enabled=true
spring.couchbase.authentication.pem.certificates=classpath:client.crt
spring.couchbase.authentication.pem.private-key=classpath:client.key
```
```
spring:
 couchbase:
 connection-string: "couchbase://192.168.1.123"
 env:
 ssl:
 enabled: true
 authentication:
 pem:
 certificates: "classpath:client.crt"
 private-key: "classpath:client.key"
```
It is also possible to customize some of the `ClusterEnvironment` settings.
For instance, the following configuration changes the timeout to open a new `Bucket` and enables SSL support with a reference to a configured SSL bundle:

-
Properties
-
YAML

```
spring.couchbase.env.timeouts.connect=3s
spring.couchbase.env.ssl.bundle=example
```
```
spring:
 couchbase:
 env:
 timeouts:
 connect: "3s"
 ssl:
 bundle: "example"
```
| Check the `spring.couchbase.env.*`properties for more details.
To take more control, one or more`ClusterEnvironmentBuilderCustomizer`beans can be used. |

### Spring Data Couchbase Repositories

Spring Data includes repository support for Couchbase.

Repositories and documents are found through scanning.
By default, the auto-configuration packages are scanned.
You can customize the locations to look for repositories and documents by using `@EnableCouchbaseRepositories` and `@EntityScan` respectively.

For complete details of Spring Data Couchbase, see the reference documentation.

You can inject an auto-configured `CouchbaseTemplate` instance as you would with any other Spring Bean, provided a `CouchbaseClientFactory` bean is available.
This happens when a `Cluster` is available, as described above, and a bucket name has been specified:

-
Properties
-
YAML

`spring.data.couchbase.bucket-name=my-bucket````
spring:
 data:
 couchbase:
 bucket-name: "my-bucket"
```
The following examples shows how to inject a `CouchbaseTemplate` bean:

-
Java
-
Kotlin

```
import org.springframework.data.couchbase.core.CouchbaseTemplate;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final CouchbaseTemplate template;
	public MyBean(CouchbaseTemplate template) {
 this.template = template;
	}
	// ...
	public String someMethod() {
 return this.template.getBucketName();
	}
}
```
```
import org.springframework.data.couchbase.core.CouchbaseTemplate
import org.springframework.stereotype.Component
@Component
class MyBean(private val template: CouchbaseTemplate) {
	// ...
	fun someMethod(): String {
 return template.bucketName
	}
}
```
There are a few beans that you can define in your own configuration to override those provided by the auto-configuration:

-
A `CouchbaseMappingContext``@Bean`with a name of`couchbaseMappingContext`.
-
A `CustomConversions``@Bean`with a name of`couchbaseCustomConversions`.
-
A `CouchbaseTemplate``@Bean`with a name of`couchbaseTemplate`.

To avoid hard-coding those names in your own config, you can reuse `BeanNames` provided by Spring Data Couchbase.
For instance, you can customize the converters to use, as follows:

-
Java
-
Kotlin

```
import org.assertj.core.util.Arrays;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.couchbase.config.BeanNames;
import org.springframework.data.couchbase.core.convert.CouchbaseCustomConversions;
@Configuration(proxyBeanMethods = false)
public class MyCouchbaseConfiguration {
	@Bean(BeanNames.COUCHBASE_CUSTOM_CONVERSIONS)
	public CouchbaseCustomConversions myCustomConversions() {
 return new CouchbaseCustomConversions(Arrays.asList(new MyConverter()));
	}
}
```
```
import org.assertj.core.util.Arrays
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.data.couchbase.config.BeanNames
import org.springframework.data.couchbase.core.convert.CouchbaseCustomConversions
@Configuration(proxyBeanMethods = false)
class MyCouchbaseConfiguration {
	@Bean(BeanNames.COUCHBASE_CUSTOM_CONVERSIONS)
	fun myCustomConversions(): CouchbaseCustomConversions {
 return CouchbaseCustomConversions(Arrays.asList(MyConverter()))
	}
}
```
## LDAP

LDAP (Lightweight Directory Access Protocol) is an open, vendor-neutral, industry standard application protocol for accessing and maintaining distributed directory information services over an IP network. Spring Boot offers auto-configuration for any compliant LDAP server as well as support for the embedded in-memory LDAP server from UnboundID.

LDAP abstractions are provided by Spring Data LDAP.
There is a `spring-boot-starter-data-ldap` starter for collecting the dependencies in a convenient way.

### Connecting to an LDAP Server

To connect to an LDAP server, make sure you declare a dependency on the `spring-boot-starter-data-ldap` starter or `spring-ldap-core` and then declare the URLs of your server in your application.properties, as shown in the following example:

-
Properties
-
YAML

```
spring.ldap.urls=ldap://myserver:1235
spring.ldap.username=admin
spring.ldap.password=secret
```
```
spring:
 ldap:
 urls: "ldap://myserver:1235"
 username: "admin"
 password: "secret"
```
If you need to customize connection settings, you can use the `spring.ldap.base` and `spring.ldap.base-environment` properties.

An `LdapContextSource` is auto-configured based on these settings.
If a `DirContextAuthenticationStrategy` bean is available, it is associated to the auto-configured `LdapContextSource`.
If you need to customize it, for instance to use a `PooledContextSource`, you can still inject the auto-configured `LdapContextSource`.
Make sure to flag your customized `ContextSource` as `@Primary` so that the auto-configured `LdapTemplate` uses it.

### Spring Data LDAP Repositories

Spring Data includes repository support for LDAP.

Repositories and documents are found through scanning.
By default, the auto-configuration packages are scanned.
You can customize the locations to look for repositories and documents by using `@EnableLdapRepositories` and `@EntityScan` respectively.

| For complete details of Spring Data LDAP, see the reference documentation. |

You can also inject an auto-configured `LdapTemplate` instance as you would with any other Spring Bean, as shown in the following example:

-
Java
-
Kotlin

```
import java.util.List;
import org.springframework.ldap.core.LdapTemplate;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final LdapTemplate template;
	public MyBean(LdapTemplate template) {
 this.template = template;
	}
	// ...
	public List<User> someMethod() {
 return this.template.findAll(User.class);
	}
}
```
```
import org.springframework.ldap.core.LdapTemplate
import org.springframework.stereotype.Component
@Component
class MyBean(private val template: LdapTemplate) {
	// ...
	fun someMethod(): List<User> {
 return template.findAll(User::class.java)
	}
}
```
### Embedded In-memory LDAP Server

For testing purposes, Spring Boot supports auto-configuration of an in-memory LDAP server from UnboundID.
To configure the server, add a dependency to `com.unboundid:unboundid-ldapsdk` and declare a `spring.ldap.embedded.base-dn` property, as follows:

-
Properties
-
YAML

`spring.ldap.embedded.base-dn=dc=spring,dc=io````
spring:
 ldap:
 embedded:
 base-dn: "dc=spring,dc=io"
```
| It is possible to define multiple base-dn values, however, since distinguished names usually contain commas, they must be defined using the correct notation. In yaml files, you can use the yaml list notation. In properties files, you must include the index as part of the property name:
 |

By default, the server starts on a random port and triggers the regular LDAP support.
There is no need to specify a `spring.ldap.urls` property.

If there is a `schema.ldif` file on your classpath, it is used to initialize the server.
If you want to load the initialization script from a different resource, you can also use the `spring.ldap.embedded.ldif` property.

By default, a standard schema is used to validate `LDIF` files.
You can turn off validation altogether by setting the `spring.ldap.embedded.validation.enabled` property.
If you have custom attributes, you can use `spring.ldap.embedded.validation.schema` to define your custom attribute types or object classes.

#### SSL

The in-memory LDAP server supports SSL (LDAPS).
To enable SSL, configure the SSL bundle to use by setting the `spring.ldap.embedded.ssl.bundle` property, as shown in the following example:

-
Properties
-
YAML

`spring.ldap.embedded.ssl.bundle=example````
spring:
 ldap:
 embedded:
 ssl:
 bundle: "example"
```

# IO

Most applications will need to deal with input and output concerns at some point. Spring Boot provides utilities and integrations with a range of technologies to help when you need IO capabilities. This section covers standard IO features such as caching and validation as well as more advanced topics such as batch, scheduling, and distributed transactions. We will also cover calling remote REST or SOAP services and sending email.

# Caching

The Spring Framework provides support for transparently adding caching to an application. At its core, the abstraction applies caching to methods, thus reducing the number of executions based on the information available in the cache. The caching logic is applied transparently, without any interference to the invoker. For more details, check the relevant section of the Spring Framework reference documentation.

Spring Boot auto-configures the cache infrastructure as long as caching support is enabled by using the `@EnableCaching` annotation.

| Avoid adding `@EnableCaching`to the main method’s application class.
Doing so makes caching a mandatory feature, including when running a test suite. |

To add caching to an operation of your service add the relevant annotation to its method, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Component;
@Component
public class MyMathService {
	@Cacheable("piDecimals")
	public int computePiDecimal(int precision) {
 ...
	}
}
```
```
import org.springframework.cache.annotation.Cacheable
import org.springframework.stereotype.Component
@Component
class MyMathService {
	@Cacheable("piDecimals")
	fun computePiDecimal(precision: Int): Int {
 ...
	}
}
```
This example demonstrates the use of caching on a potentially costly operation.
Before invoking `computePiDecimal`, the abstraction looks for an entry in the `piDecimals` cache that matches the `precision` argument.
If an entry is found, the content in the cache is immediately returned to the caller, and the method is not invoked.
Otherwise, the method is invoked, and the cache is updated before returning the value.

| You can also use the standard JSR-107 (JCache) annotations (such as `@CacheResult`) transparently.
However, we strongly advise you to not mix and match the Spring Cache and JCache annotations. |

If you do not add any specific cache library, Spring Boot auto-configures a simple provider that uses concurrent maps in memory.
When a cache is required (such as `piDecimals` in the preceding example), this provider creates it for you.
The simple provider is not really recommended for production usage, but it is great for getting started and making sure that you understand the features.
When you have made up your mind about the cache provider to use, please make sure to read its documentation to figure out how to configure the caches that your application uses.
Nearly all providers require you to explicitly configure every cache that you use in the application.
Some offer a way to customize the default caches defined by the `spring.cache.cache-names` property.

## Supported Cache Providers

The cache abstraction does not provide an actual store and relies on abstraction materialized by the `Cache` and `CacheManager` interfaces.

If you have not defined a bean of type `CacheManager` or a `CacheResolver` named `cacheResolver` (see `CachingConfigurer`), Spring Boot tries to detect the following providers (in the indicated order):

-
JCache (JSR-107) (EhCache 3, Hazelcast, Infinispan, and others)

| If the `CacheManager`is auto-configured by Spring Boot, it is possible toforcea particular cache provider by setting the`spring.cache.type`property. |

| Use the `spring-boot-starter-cache`starter to quickly add basic caching dependencies.
The starter brings in`spring-context-support`.
If you add dependencies manually, you must include`spring-context-support`in order to use the JCache or Caffeine support. |

If the `CacheManager` is auto-configured by Spring Boot, you can further tune its configuration before it is fully initialized by exposing a bean that implements the `CacheManagerCustomizer` interface.
The following example sets a flag to say that `null` values should not be passed down to the underlying map:

-
Java
-
Kotlin

```
import org.springframework.boot.cache.autoconfigure.CacheManagerCustomizer;
import org.springframework.cache.concurrent.ConcurrentMapCacheManager;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyCacheManagerConfiguration {
	@Bean
	public CacheManagerCustomizer<ConcurrentMapCacheManager> cacheManagerCustomizer() {
 return (cacheManager) -> cacheManager.setAllowNullValues(false);
	}
}
```
```
import org.springframework.boot.cache.autoconfigure.CacheManagerCustomizer
import org.springframework.cache.concurrent.ConcurrentMapCacheManager
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
@Configuration(proxyBeanMethods = false)
class MyCacheManagerConfiguration {
	@Bean
	fun cacheManagerCustomizer(): CacheManagerCustomizer<ConcurrentMapCacheManager> {
 return CacheManagerCustomizer { cacheManager ->
 cacheManager.isAllowNullValues = false
 }
	}
}
```
| In the preceding example, an auto-configured `ConcurrentMapCacheManager`is expected.
If that is not the case (either you provided your own config or a different cache provider was auto-configured), the customizer is not invoked at all.
You can have as many customizers as you want, and you can also order them by using`@Order`or`Ordered`. |

### Generic

Generic caching is used if the context defines *at least* one `Cache` bean.
A `CacheManager` wrapping all beans of that type is created.

### JCache (JSR-107)

JCache is bootstrapped through the presence of a `CachingProvider` on the classpath (that is, a JSR-107 compliant caching library exists on the classpath), and the `JCacheCacheManager` is provided by the `spring-boot-starter-cache` starter.
Various compliant libraries are available, and Spring Boot provides dependency management for Ehcache 3, Hazelcast, and Infinispan.
Any other compliant library can be added as well.

It might happen that more than one provider is present, in which case the provider must be explicitly specified. Even if the JSR-107 standard does not enforce a standardized way to define the location of the configuration file, Spring Boot does its best to accommodate setting a cache with implementation details, as shown in the following example:

-
Properties
-
YAML

```
spring.cache.jcache.provider=com.example.MyCachingProvider
spring.cache.jcache.config=classpath:example.xml
```
```
# Only necessary if more than one provider is present
spring:
 cache:
 jcache:
 provider: "com.example.MyCachingProvider"
 config: "classpath:example.xml"
```
| When a cache library offers both a native implementation and JSR-107 support, Spring Boot prefers the JSR-107 support, so that the same features are available if you switch to a different JSR-107 implementation. |

| Spring Boot has general support for Hazelcast.
If a single `HazelcastInstance`is available, it is automatically reused for the`CacheManager`as well, unless the`spring.cache.jcache.config`property is specified. |

There are two ways to customize the underlying `CacheManager`:

-
Caches can be created on startup by setting the `spring.cache.cache-names`property. If a custom`Configuration`bean is defined, it is used to customize them.
-
`CacheManagerCustomizer`beans are invoked with the reference of the`CacheManager`for full customization.

| If a standard `CacheManager`bean is defined, it is wrapped automatically in an`CacheManager`implementation that the abstraction expects.
No further customization is applied to it. |

### Hazelcast

Spring Boot has general support for Hazelcast.
If a `HazelcastInstance` has been auto-configured and `com.hazelcast:hazelcast-spring` is on the classpath, it is automatically wrapped in a `CacheManager`.

| Hazelcast can be used as a JCache compliant cache or as a Spring `CacheManager`compliant cache.
When setting`spring.cache.type`to`hazelcast`, Spring Boot will use the`CacheManager`based implementation.
If you want to use Hazelcast as a JCache compliant cache, set`spring.cache.type`to`jcache`.
If you have multiple JCache compliant cache providers and want to force the use of Hazelcast, you have to explicitly set the JCache provider. |

### Infinispan

Infinispan has no default configuration file location, so it must be specified explicitly. Otherwise, the default bootstrap is used.

-
Properties
-
YAML

`spring.cache.infinispan.config=infinispan.xml````
spring:
 cache:
 infinispan:
 config: "infinispan.xml"
```
Caches can be created on startup by setting the `spring.cache.cache-names` property.
If a custom `ConfigurationBuilder` bean is defined, it is used to customize the caches.

For more details, see the documentation.

### Couchbase

If Spring Data Couchbase is available and Couchbase is configured, a `CouchbaseCacheManager` is auto-configured.
It is possible to create additional caches on startup by setting the `spring.cache.cache-names` property and cache defaults can be configured by using `spring.cache.couchbase.*` properties.
For instance, the following configuration creates `cache1` and `cache2` caches with an entry *expiration* of 10 minutes:

-
Properties
-
YAML

```
spring.cache.cache-names=cache1,cache2
spring.cache.couchbase.expiration=10m
```
```
spring:
 cache:
 cache-names: "cache1,cache2"
 couchbase:
 expiration: "10m"
```
If you need more control over the configuration, consider registering a `CouchbaseCacheManagerBuilderCustomizer` bean.
The following example shows a customizer that configures a specific entry expiration for `cache1` and `cache2`:

-
Java
-
Kotlin

```
import java.time.Duration;
import org.springframework.boot.cache.autoconfigure.CouchbaseCacheManagerBuilderCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.couchbase.cache.CouchbaseCacheConfiguration;
@Configuration(proxyBeanMethods = false)
public class MyCouchbaseCacheManagerConfiguration {
	@Bean
	public CouchbaseCacheManagerBuilderCustomizer myCouchbaseCacheManagerBuilderCustomizer() {
 return (builder) -> builder
 .withCacheConfiguration("cache1", CouchbaseCacheConfiguration
 .defaultCacheConfig().entryExpiry(Duration.ofSeconds(10)))
 .withCacheConfiguration("cache2", CouchbaseCacheConfiguration
 .defaultCacheConfig().entryExpiry(Duration.ofMinutes(1)));
	}
}
```
```
import org.springframework.boot.cache.autoconfigure.CouchbaseCacheManagerBuilderCustomizer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.data.couchbase.cache.CouchbaseCacheConfiguration
import java.time.Duration
@Configuration(proxyBeanMethods = false)
class MyCouchbaseCacheManagerConfiguration {
	@Bean
	fun myCouchbaseCacheManagerBuilderCustomizer(): CouchbaseCacheManagerBuilderCustomizer {
 return CouchbaseCacheManagerBuilderCustomizer { builder ->
 builder
 .withCacheConfiguration(
 "cache1", CouchbaseCacheConfiguration
 .defaultCacheConfig().entryExpiry(Duration.ofSeconds(10))
 )
 .withCacheConfiguration(
 "cache2", CouchbaseCacheConfiguration
 .defaultCacheConfig().entryExpiry(Duration.ofMinutes(1))
 )
 }
	}
}
```
### Redis

If Redis is available and configured, a `RedisCacheManager` is auto-configured.
It is possible to create additional caches on startup by setting the `spring.cache.cache-names` property and cache defaults can be configured by using `spring.cache.redis.*` properties.
For instance, the following configuration creates `cache1` and `cache2` caches with a *time to live* of 10 minutes:

-
Properties
-
YAML

```
spring.cache.cache-names=cache1,cache2
spring.cache.redis.time-to-live=10m
```
```
spring:
 cache:
 cache-names: "cache1,cache2"
 redis:
 time-to-live: "10m"
```
| By default, a key prefix is added so that, if two separate caches use the same key, Redis does not have overlapping keys and cannot return invalid values.
We strongly recommend keeping this setting enabled if you create your own `RedisCacheManager`. |

| You can take full control of the default configuration by adding a `RedisCacheConfiguration``@Bean`of your own.
This can be useful if you need to customize the default serialization strategy. |

If you need more control over the configuration, consider registering a `RedisCacheManagerBuilderCustomizer` bean.
The following example shows a customizer that configures a specific time to live for `cache1` and `cache2`:

-
Java
-
Kotlin

```
import java.time.Duration;
import org.springframework.boot.cache.autoconfigure.RedisCacheManagerBuilderCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.redis.cache.RedisCacheConfiguration;
@Configuration(proxyBeanMethods = false)
public class MyRedisCacheManagerConfiguration {
	@Bean
	public RedisCacheManagerBuilderCustomizer myRedisCacheManagerBuilderCustomizer() {
 return (builder) -> builder
 .withCacheConfiguration("cache1", RedisCacheConfiguration
 .defaultCacheConfig().entryTtl(Duration.ofSeconds(10)))
 .withCacheConfiguration("cache2", RedisCacheConfiguration
 .defaultCacheConfig().entryTtl(Duration.ofMinutes(1)));
	}
}
```
```
import org.springframework.boot.cache.autoconfigure.RedisCacheManagerBuilderCustomizer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.data.redis.cache.RedisCacheConfiguration
import java.time.Duration
@Configuration(proxyBeanMethods = false)
class MyRedisCacheManagerConfiguration {
	@Bean
	fun myRedisCacheManagerBuilderCustomizer(): RedisCacheManagerBuilderCustomizer {
 return RedisCacheManagerBuilderCustomizer { builder ->
 builder
 .withCacheConfiguration(
 "cache1", RedisCacheConfiguration
 .defaultCacheConfig().entryTtl(Duration.ofSeconds(10))
 )
 .withCacheConfiguration(
 "cache2", RedisCacheConfiguration
 .defaultCacheConfig().entryTtl(Duration.ofMinutes(1))
 )
 }
	}
}
```
### Caffeine

Caffeine is a Java 8 rewrite of Guava’s cache that supersedes support for Guava.
If Caffeine is present, a `CaffeineCacheManager` (provided by the `spring-boot-starter-cache` starter) is auto-configured.
Caches can be created on startup by setting the `spring.cache.cache-names` property and can be customized by one of the following (in the indicated order):

-
A cache spec defined by `spring.cache.caffeine.spec`
-
A `CaffeineSpec`bean is defined
-
A `Caffeine`bean is defined

For instance, the following configuration creates `cache1` and `cache2` caches with a maximum size of 500 and a *time to live* of 10 minutes

-
Properties
-
YAML

```
spring.cache.cache-names=cache1,cache2
spring.cache.caffeine.spec=maximumSize=500,expireAfterAccess=600s
```
```
spring:
 cache:
 cache-names: "cache1,cache2"
 caffeine:
 spec: "maximumSize=500,expireAfterAccess=600s"
```
If a `CacheLoader` bean is defined, it is automatically associated to the `CaffeineCacheManager`.
Since the `CacheLoader` is going to be associated with *all* caches managed by the cache manager, it must be defined as `CacheLoader<Object, Object>`.
The auto-configuration ignores any other generic type.

### Cache2k

Cache2k is an in-memory cache.
If the Cache2k spring integration is present, a `SpringCache2kCacheManager` is auto-configured.

Caches can be created on startup by setting the `spring.cache.cache-names` property.
Cache defaults can be customized using a `Cache2kBuilderCustomizer` bean.
The following example shows a customizer that configures the capacity of the cache to 200 entries, with an expiration of 5 minutes:

-
Java
-
Kotlin

```
import java.util.concurrent.TimeUnit;
import org.springframework.boot.cache.autoconfigure.Cache2kBuilderCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyCache2kDefaultsConfiguration {
	@Bean
	public Cache2kBuilderCustomizer myCache2kDefaultsCustomizer() {
 return (builder) -> builder.entryCapacity(200)
 .expireAfterWrite(5, TimeUnit.MINUTES);
	}
}
```
```
import org.springframework.boot.cache.autoconfigure.Cache2kBuilderCustomizer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import java.util.concurrent.TimeUnit
@Configuration(proxyBeanMethods = false)
class MyCache2kDefaultsConfiguration {
	@Bean
	fun myCache2kDefaultsCustomizer(): Cache2kBuilderCustomizer {
 return Cache2kBuilderCustomizer { builder ->
 builder.entryCapacity(200)
 .expireAfterWrite(5, TimeUnit.MINUTES)
 }
	}
}
```
### Simple

If none of the other providers can be found, a simple implementation using a `ConcurrentHashMap` as the cache store is configured.
This is the default if no caching library is present in your application.
By default, caches are created as needed, but you can restrict the list of available caches by setting the `cache-names` property.
For instance, if you want only `cache1` and `cache2` caches, set the `cache-names` property as follows:

-
Properties
-
YAML

`spring.cache.cache-names=cache1,cache2````
spring:
 cache:
 cache-names: "cache1,cache2"
```
If you do so and your application uses a cache not listed, then it fails at runtime when the cache is needed, but not on startup. This is similar to the way the "real" cache providers behave if you use an undeclared cache.

## Testing

It is generally useful to use a no-op implementation when running a test suite. This section lists a number of strategies that are useful for tests.

When a custom `CacheManager` is defined, the best option is to make sure that caching configuration is defined in an isolated `@Configuration` class.
Doing so makes sure that caching is not required by slice tests.
For tests that enable a full context, such as `@SpringBootTest`, an explicit configuration overriding the regular configuration is required.

If caching is auto-configured, more options are available.
Tests can be annotated with `@AutoConfigureCache` to replace the auto-configured `CacheManager` by a no-op implementation.

-
Java
-
Kotlin

```
import org.springframework.boot.cache.test.autoconfigure.AutoConfigureCache;
import org.springframework.boot.test.context.SpringBootTest;
@SpringBootTest
@AutoConfigureCache
public class MyIntegrationTests {
	// Tests use a no-op cache manager
}
```
```
import org.springframework.boot.cache.test.autoconfigure.AutoConfigureCache
import org.springframework.boot.test.context.SpringBootTest
@SpringBootTest
@AutoConfigureCache
class MyIntegrationTests {
	// Tests use a no-op cache manager
}
```
Another option is to force a no-op implementation for the auto-configured `CacheManager`:

-
Properties
-
YAML

`spring.cache.type=none````
spring:
 cache:
 type: "none"
```

# Spring Batch

Spring Boot offers several conveniences for working with Spring Batch, including running a Job on startup.

When building a batch application, the following stores can be auto-configured:

-
In-memory
-
JDBC
-
MongoDB

Each store has specific additional settings. For instance, it is possible to customize the tables prefix for the JDBC store, as shown in the following example:

-
Properties
-
YAML

`spring.batch.jdbc.table-prefix=CUSTOM_````
spring:
 batch:
 jdbc:
 table-prefix: "CUSTOM_"
```
When using the MongoDB store, you can enable initialization of the Spring Batch job repository schema (collections and indexes):

-
Properties
-
YAML

`spring.batch.data.mongodb.schema.initialize=true````
spring:
 batch:
 data:
 mongodb:
 schema:
 initialize: true
```
To disable Spring Boot’s auto-configuration and take complete control of Spring Batch’s configuration, add `@EnableBatchProcessing` to one of your `@Configuration` classes or extend `DefaultBatchConfiguration`.
This will cause the auto-configuration to back off, including initialization of Spring Batch’s database schema (JDBC or MongoDB).
Spring Batch can then be configured using the `@Enable*JobRepository` annotation’s attributes rather than the previously described configuration properties.

To learn more about manually configuring Spring Batch, see the API documentation of:

For more information about Spring Batch, see the Spring Batch project page.

## Running Spring Batch Jobs on Startup

When Spring Boot auto-configures Spring Batch, and if a single `Job` bean is found in the application context, it is executed on startup (see `JobLauncherApplicationRunner` for details).
If multiple `Job` beans are found, the job that should be executed must be specified using `spring.batch.job.name`.

You can disable running a `Job` found in the application context, as shown in the following example:

-
Properties
-
YAML

`spring.batch.job.enabled=false````
spring:
 batch:
 job:
 enabled: false
```
See `BatchAutoConfiguration`, `BatchJdbcAutoConfiguration`, and `BatchDataMongoAutoConfiguration` for more details.

# gRPC

Google Remote Procedure Call (gRPC) is a high-performance RPC framework that enables client-server communication using binary messages. Spring Boot include support for developing and testing both client and server gRPC applications.

The underling message format used by gRPC is Protocol Buffers which allow messages to be created and consumed by a wide variety of programming languages.

## Service Definitions

To develop a gRPC application you first need a Protocol Buffers service definition file.
A `.proto` file defines the services and messages that your application can consume or provide.

Here’s an example of a typical `.proto` file that uses the `proto3` revision of the protocol buffers language:

```
syntax = "proto3";
option java_package = "com.example.grpc.proto";
option java_multiple_files = true;
service HelloWorld {
 rpc SayHello (HelloRequest) returns (HelloReply) {}
}
message HelloRequest {
 string name = 1;
}
message HelloReply {
 string message = 1;
}
```
This file defines a `HelloWorld` service with a single method that accepts a `HelloReqest` message and return a `HelloReply` message.
The `HelloReqest` message contains a `name` string field.
The `HelloReply` message contains a `message` string field.

With the exception of a the `java_package` and `java_multiple_files` options, there is nothing in the `.proto` file that is specific to the Java programming langage.

### Generating Java Code

Since `.proto` files are language agnostic, we need a process to convert them into usable Java code.
We can then use the generated code to either make a remote procedure call to running service, or implement the service ourselves so that others may call it.

The exact process you use to generate code will depend on your build system. Spring Boot supports for both Maven and Gradle protobuf plugins, but you are free to use whatever solution works best for you.

#### Using the Maven Plugin

Spring Boot include dependency management for the `io.github.ascopes:protobuf-maven-plugin` Maven plugin.
If you are using the the `spring-boot-starter-parent` POM, you’ll also get sensible out-of-the-box configuration.

The following shows a typical Maven POM file that uses the plugin:

```
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
	<modelVersion>4.0.0</modelVersion>
 <parent>
 <groupId>org.springframework.boot</groupId>
 <artifactId>spring-boot-starter-parent</artifactId>
 <version>4.1.0</version>
	</parent>
	<groupId>com.example</groupId>
	<artifactId>myproject</artifactId>
	<version>0.0.1-SNAPSHOT</version>
	<build>
 <plugins>
 <plugin>
 <groupId>io.github.ascopes</groupId>
 <artifactId>protobuf-maven-plugin</artifactId>
 </plugin>
 <plugin>
 <groupId>org.springframework.boot</groupId>
 <artifactId>spring-boot-maven-plugin</artifactId>
 </plugin>
 </plugins>
	</build>
</project>
```
Since POM above extends `spring-boot-starter-parent`, you’ll get the following:

-
Configuration of the `protoc`version.
-
Configuration of the `binary-maven`plugin.
-
Execution configuration for the `generate`goal.

The `.proto` files should be added to `src/main/proto`.

| If you don’t use `spring-boot-starter-parent`, or you want to configure the plugin directly, refer to the protobuf-maven-plugin documentation.
If you use Spring Boot’s dependency management the`${protobuf-java.version}`and`${grpc-java.version}`properties will be useful. |

#### Using the Gradle Plugin

Spring Boot include dependency management for the `com.google.protobuf:protobuf-gradle-plugin` Gradle plugin.
In addition, the `spring-boot-gradle-plugin` will react to the presence of the protobuf plugin and configure it appropriately.

The following shows a typical Gradle file that uses the plugin:

```
plugins {
	id 'java'
	id 'org.springframework.boot' version '4.1.0'
	id 'io.spring.dependency-management' version '1.1.7'
	id 'com.google.protobuf' version '0.9.6'
}
group = 'com.example'
version = '0.0.1-SNAPSHOT'
java {
	toolchain {
 languageVersion = JavaLanguageVersion.of(17)
	}
}
repositories {
	mavenCentral()
}
```
Since this gradle file uses both the `org.springframework.boot` and `com.google.protobuf` plugins, you’ll get the following:

-
Configuration of the `protoc`version.
-
Configuration of the `protoc-gen-grpc-java`version.

The `.proto` files should be added to `src/main/proto`.

| If you don’t use `org.springframework.boot`plugin, or you want to configure the plugin directly, refer to the protobuf-gradle-plugin documentation. |

## Writing a gRPC Server Application

Spring Boot provides a `spring-boot-grpc-server` module and a `spring-boot-starter-grpc-server` starter POM that you can use for server applications.

In order to write the actual server code, you’ll need to extended one or more of base classes generated from your `.proto` file and expose them as Spring beans.
Spring gRPC will automatically expose any bean that implements `BindableService` as a gRPC server.
Since all `.proto` generated classes implement `BindableService`, adding them as beans is enough to expose them over gRPC.

| For more details see the Spring gRPC documentation. |

The following example shows how the `HelloWorld` service from the `.proto` file above could be implemented.
In this example, we’re using the `@GrpcService` annotation and assuming that the code is in a package that will be picked up by component scanning:

-
Java
-
Kotlin

```
import io.grpc.stub.StreamObserver;
import org.springframework.grpc.server.service.GrpcService;
@GrpcService
public class MyHelloWorldService extends HelloWorldGrpc.HelloWorldImplBase {
	@Override
	public void sayHello(HelloRequest request, StreamObserver<HelloReply> responseObserver) {
 String message = "Hello '%s'".formatted(request.getName());
 HelloReply reply = HelloReply.newBuilder().setMessage(message).build();
 responseObserver.onNext(reply);
 responseObserver.onCompleted();
	}
}
```
```
import io.grpc.stub.StreamObserver
import org.springframework.grpc.server.service.GrpcService
@GrpcService
class MyHelloWorldService : HelloWorldGrpc.HelloWorldImplBase() {
	override fun sayHello(request: HelloRequest, responseObserver: StreamObserver<HelloReply>) {
 val message = "Hello '${request.getName()}'"
 val reply = HelloReply.newBuilder().setMessage(message).build()
 responseObserver.onNext(reply)
 responseObserver.onCompleted()
	}
}
```
If the application makes used of `spring-boot-starter-grpc-server`, then Netty will be used as the server implementation listening on port `9090`.

You can test your application using grpcurl:

`$ grpcurl -d '{"name":"Spring"}' -plaintext localhost:9090 HelloWorld.SayHello````
{
 "message": "Hello 'Spring'"
}
```
### Switching to a Netty Shaded Server

If you find that the version of Netty provided by the `spring-boot-starter-grpc-server` starter POM isn’t compatible with other libraries you use, you can switch to a “shaded” version.

To switch, you can excluded `io.grpc:grpc-netty` and include `io.grpc:grpc-netty-shaded`.
For example:

-
Maven
-
Gradle

```
<dependency>
	<groupId>org.springframework.boot</groupId>
	<artifactId>spring-boot-starter-grpc-server</artifactId>
	<exclusions>
 <!-- Exclude the gRPC Netty dependency -->
 <exclusion>
 <groupId>io.grpc</groupId>
 <artifactId>grpc-netty</artifactId>
 </exclusion>
	</exclusions>
</dependency>
<!-- Use gRPC Netty Shaded instead -->
<dependency>
	<groupId>io.grpc</groupId>
	<artifactId>grpc-netty-shaded</artifactId>
</dependency>
```
```
dependencies {
	implementation('org.springframework.boot:spring-boot-starter-grpc-server') {
 // Exclude the gRPC Netty dependency
 exclude group: 'io.grpc', module: 'grpc-netty'
	}
	// Use gRPC Netty Shaded instead
	implementation "io.grpc:grpc-netty-shaded"
}
```
### Switching to a Servlet Container

It’s possible to expose gRPC services using a regular Servlet Container such as Tomcat rather than using Netty. To do so, your Servlet Container must be configured to support HTTP/2.

To switch to the Servlet gRPC implementation, you can exclude `io.grpc:grpc-netty` and include `io.grpc:grpc-servlet-jakarta`.
For example:

-
Maven
-
Gradle

```
<dependency>
	<groupId>org.springframework.boot</groupId>
	<artifactId>spring-boot-starter-webmvc</artifactId>
</dependency>
<dependency>
	<groupId>org.springframework.boot</groupId>
	<artifactId>spring-boot-starter-grpc-server</artifactId>
	<exclusions>
 <!-- Exclude the gRPC Netty dependency -->
 <exclusion>
 <groupId>io.grpc</groupId>
 <artifactId>grpc-netty</artifactId>
 </exclusion>
	</exclusions>
</dependency>
<!-- Use gRPC Servlet Jakarta instead -->
<dependency>
	<groupId>io.grpc</groupId>
	<artifactId>grpc-servlet-jakarta</artifactId>
</dependency>
```
```
dependencies {
	implementation('org.springframework.boot:spring-boot-starter-webmvc') {
	implementation('org.springframework.boot:spring-boot-starter-grpc-server') {
 // Exclude the gRPC Netty dependency
 exclude group: 'io.grpc', module: 'grpc-netty'
	}
	// Use gRPC Servlet Jakarta instead
	implementation "io.grpc:grpc-servlet-jakarta"
}
```
| Remember to include a Servlet Container dependency, for example using `spring-boot-starter-tomcat`, and to set`server.http2.enabled`to`true`. |

| When using a servlet container, certain gRPC server configuration properties are not relevant and will be ignored.
For example, `spring.grpc.server.port`is ignored since`server.port`used used to set a web server port. |

### SSL Support

SSL can be configured for both `grpc-netty` and `grpc-netty-shaded` servers using SSL bundles.
See the SSL core documentation for details on how to declare an SSL bundle.

Once your bundle has been defined, you can use the following properties in your gRPC server application to use it:

-
Properties
-
YAML

`spring.grpc.server.ssl.bundle=mysslbundle````
spring:
 grpc:
 server:
 ssl:
 bundle: mysslbundle
```
Client authentication can also be configured by setting `spring.grpc.server.ssl.client-auth` to `optional` or `require`.

| To temporarily disable server SSL support, for example to aid with testing, you can set `spring.grpc.server.ssl.enabled`to`false`. |

### Using an In-Process Server

You can run an in-process server by including the `io.grpc:grpc-inprocess` dependency on your classpath and defining a `spring.grpc.server.inprocess.name` property.
In this mode, the in-process server factory is auto-configured in addition to the regular server factory.

The name you provide can be used as a client channel target using the form `in-process:<name>`.

### Reflection

When it’s available, Spring Boot will auto-configure the gRPC Reflection service.
This allows clients to browse the metadata of your services and download their `.proto` files.

The reflection service resides in the `io.grpc:grpc-services` library, which is an optional dependency.
You will need to add the dependency to your project in order for auto-configuration to apply.

| If you have the `io.grpc:grpc-services`library but prefer that reflection isn’t auto-configured, you can set`spring.grpc.server.reflection.enabled`to`false`. |

### Server Health

A gRPC server can provide health information using a standard service API (health/v1). This allows clients to check on the health of your server services a route traffic appropriately.

Spring Boot provides a bridge between its own `spring-boot-health` module and the standard gRPC health service.
Health information is provided whenever the `io.grpc:grpc-services` and `org.springframework.boot:spring-boot-health` modules are on your classpath.

| If you don’t want health indicators to be exposed, you can set `spring.grpc.server.health.enabled`to`false`. |

#### Service Specific Health Mappings

By default, health information is provided for the overall server status (`""`) using all available health indicators.

It is also possible to provide fine-grained health information for specific services by including only a sub-set of health indicators. Custom mapping and ordering rules can also be defined on a per-service basis.

For example, the following configuration will provide health for “myservice” using only the `db` and `redis` indicators.

-
Properties
-
YAML

```
spring.grpc.server.health.service.myservice.include[0]=db
spring.grpc.server.health.service.myservice.include[1]=redis
```
```
spring:
 grpc:
 server:
 health:
 service:
 myservice:
 include:
 - db
 - redis
```
| You can set `spring.grpc.server.health.include-overall-health`to`false`to disable the overall server status health if you only want to provide service-specific health. |

#### Push Configuration

Unlike web-based health checks, gRPC health information is periodically pushed rather than pulled. By default, the first health push happens 5 seconds after the application starts and then every subsequent 5 seconds.

To fine-tune this, you can use the following properties:

-
Properties
-
YAML

```
spring.grpc.server.health.schedule.period=5m
spring.grpc.server.health.schedule.delay=2s
```
```
spring:
 grpc:
 server:
 health:
 schedule:
 period: 5m
 delay: 2s
```
| You can also set `spring.grpc.server.health.schedule.enabled`to`false`if want to send health updates in some other way. |

### Securing gRPC Server Applications

#### Netty Based Servers

Spring gRPC includes features that allow you to secure your Netty based server applications declaratively using Spring Security. This follows similar patterns to those you would use to secure a regular web application.

Spring Boot provides auto-configuration for both `GrpcSecurity` and `SecurityGrpcExceptionHandler` beans.
Typically gRPC application are then secured using `@PreAuthorize` annotations on your gRPC service beans, or a `AuthenticationProcessInterceptor` bean.

For more details, please see the Spring gRPC documentation.

#### Servlet Container Based Servers

If your gRPC server is running within a standard Servlet Container, you can use typical web security configuration to secure your application.
Spring Boot will auto-configre `GrpcServerExecutorProvider` and `SecurityContextServerInterceptor` beans to ensure that Spring Security works correctly.

Cross-Site Request Forgery (CSRF) protection is incompatible with the gRPC protocol and will be disabled by default for all gRPC requests.
If you prefer to configure your own CSRF protection, you can switch this off by setting `spring.grpc.server.security.csrf.enabled` to `false`.

To help with manual Security configuration, Spring Boot provides request matchers for gRPC services. Matches are available for both servlet and reactive stacks. For example, the following will include all gRPC services with the exception of “special”.

-
Java
-
Kotlin

```
import org.springframework.boot.grpc.server.autoconfigure.security.web.servlet.GrpcRequest;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.SecurityFilterChain;
import static org.springframework.security.config.Customizer.withDefaults;
@Configuration(proxyBeanMethods = false)
public class MySecurityConfiguration {
	@Bean
	public SecurityFilterChain securityFilterChain(HttpSecurity http) {
 http.securityMatcher(GrpcRequest.toAnyService().excluding("special"));
 http.authorizeHttpRequests((requests) -> requests.anyRequest().hasRole("GRPC_ADMIN"));
 http.httpBasic(withDefaults());
 return http.build();
	}
}
```
```
import org.springframework.boot.grpc.server.autoconfigure.security.web.servlet.GrpcRequest
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.security.config.Customizer.withDefaults
import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.web.SecurityFilterChain
@Configuration(proxyBeanMethods = false)
class MySecurityConfiguration {
	@Bean
	fun securityFilterChain(http: HttpSecurity): SecurityFilterChain {
 http.securityMatcher(GrpcRequest.toAnyService().excluding("special"))
 http.authorizeHttpRequests { requests -> requests.anyRequest().hasRole("GRPC_ADMIN") }
 http.httpBasic(withDefaults())
 return http.build()
	}
}
```
| For reactive matches use `org.springframework.boot.grpc.server.autoconfigure.security.web.servlet.GrpcRequest`. |

#### OAuth2 Resource Server

OAuth2 is a widely used authorization framework. Spring Boot’s OAuth2 Resource Server support in compatible with gRPC and may be configured in the usual way.

For details of how to configure an OAuth2 Resource Server to use with your gRPC server application, see the “OAuth2” section of under “Security”.

## Writing a gRPC Client Application

Spring Boot provides a `spring-boot-grpc-client` module and a `spring-boot-starter-grpc-client` starter POM that you can use for client applications.

Clients can call remote gRPC services by importing one or more of the “stub” classes generated from their `.proto` file.
You can use the `@ImportGrpcClients` annotation to import the stub classes you want to use.

Each import includes a `target` which can either be a logical channel name, or the base URL of the remote server.
We typically recommend using channel names rather than hard-coding targets.

Here’s a typical example:

-
Java
-
Kotlin

```
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.grpc.client.ImportGrpcClients;
@SpringBootApplication(proxyBeanMethods = false)
@ImportGrpcClients(target = "hello", types = HelloWorldGrpc.HelloWorldBlockingStub.class)
public class MyApplication {
	public static void main(String[] args) {
 SpringApplication.run(MyApplication.class, args);
	}
}
```
```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.docs.features.springapplication.MyApplication
import org.springframework.boot.runApplication
import org.springframework.grpc.client.ImportGrpcClients
@SpringBootApplication(proxyBeanMethods = false)
@ImportGrpcClients(target = "hello", types = [HelloWorldGrpc.HelloWorldBlockingStub::class])
class MyApplication
fun main(args: Array<String>) {
	runApplication<MyApplication>(*args)
}
```
| If you don’t specify a `target`then “default” is used. |

| You can use the `basePackageClasses`or`basePackages`attribute of`@ImportGrpcClients`to import all stubs in given package. |

### Channel Properties

When the `target` attribute of `@ImportGrpcClients` uses a logical channel name, you’ll need to provide some properties so that Spring Boot can find the actual gRPC server to call.

To do that, you can add an entry using `spring.grpc.channel.<name>.*` properties.
Typically, you’ll configure the actual real `target`, along with any other settings that are unique to the channel.

For example, the following will configure `myservice` to use the real target of `static://grpc.example.com:9090`.
It also changes the keep-alive timeout and the maximum message size permitted:

-
Properties
-
YAML

```
spring.grpc.client.channel.myservice.target=static://grpc.example.com:9090
spring.grpc.client.channel.myservice.inbound.keepalive.timeout=40s
spring.grpc.client.channel.myservice.inbound.message.max-size=8MB
```
```
spring:
 grpc:
 client:
 channel:
 myservice:
 target: static://grpc.example.com:9090
 inbound:
 keepalive:
 timeout: 40s
 message:
 max-size: 8MB
```
### Using Stub Beans

With your `@ImportGrpcClients` annotations in place, and your properties written, you can use stubs as you would any other bean.

For example, here’s the `HelloWorldStub` being injected into a `ApplicationRunner` bean:

-
Java
-
Kotlin

```
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.boot.docs.io.grpc.client.stubbeans.HelloWorldGrpc.HelloWorldBlockingStub;
import org.springframework.stereotype.Component;
@Component
class MyApplicationRunner implements ApplicationRunner {
	private final HelloWorldBlockingStub helloStub;
	MyApplicationRunner(HelloWorldGrpc.HelloWorldBlockingStub helloStub) {
 this.helloStub = helloStub;
	}
	@Override
	public void run(ApplicationArguments args) throws Exception {
 HelloRequest request = HelloRequest.newBuilder().setName("Spring").build();
 HelloReply reply = this.helloStub.sayHello(request);
 System.out.println(reply.getMessage());
	}
}
```
```
import org.springframework.boot.ApplicationArguments
import org.springframework.boot.ApplicationRunner
class MyApplicationRunner(val helloStub: HelloWorldGrpc.HelloWorldBlockingStub) : ApplicationRunner {
	override fun run(args: ApplicationArguments) {
 val request = HelloRequest.newBuilder().setName("Spring").build()
 val reply: HelloReply = helloStub.sayHello(request)
 println(reply.getMessage())
	}
}
```
### Switching to Netty Shaded Client Transport

Under the hood, remote gRPC network calls are made using Netty.
If you find that the version of Netty provided by the `spring-boot-starter-grpc-client` starter POM isn’t compatible with other libraries you use, you can switch to a “shaded” version.

To switch, you can excluded `io.grpc:grpc-netty` and include `io.grpc:grpc-netty-shaded`.
For example:

-
Maven
-
Gradle

```
<dependency>
	<groupId>org.springframework.boot</groupId>
	<artifactId>spring-boot-starter-grpc-client</artifactId>
	<exclusions>
 <!-- Exclude the gRPC Netty dependency -->
 <exclusion>
 <groupId>io.grpc</groupId>
 <artifactId>grpc-netty</artifactId>
 </exclusion>
	</exclusions>
</dependency>
<!-- Use gRPC Netty Shaded instead -->
<dependency>
	<groupId>io.grpc</groupId>
	<artifactId>grpc-netty-shaded</artifactId>
</dependency>
```
```
dependencies {
	implementation('org.springframework.boot:spring-boot-starter-grpc-client') {
 // Exclude the gRPC Netty dependency
 exclude group: 'io.grpc', module: 'grpc-netty'
	}
	// Use gRPC Netty Shaded instead
	implementation "io.grpc:grpc-netty-shaded"
}
```
### SSL Support

Client gRPC applications can connect to gRPC services using SSL/TSL encrypted connections. You can configure gRPC connections to use standard one-way-TLS, or mutual TLS

#### Standard one-way TLS

To use standard one-way TLS, you can set the `ssl.enabled` property to `true` in your channel properties.
For example, the following will enabled an SSL/TLS connection for the `myservice` channel:

-
Properties
-
YAML

```
spring.grpc.client.channel.myservice.target=static://grpc.example.com:9090
spring.grpc.client.channel.myservice.ssl.enabled=true
```
```
spring:
 grpc:
 client:
 channel:
 myservice:
 target: static://grpc.example.com:9090
 ssl:
 enabled: true
```
#### Mutual TLS

Mutual TLS (mTLS) is a security protocol that requires both the client and the server to present certificates to each other.
To use mutual TLS, you can set the `ssl.bundle` property in your channel properties.
See the SSL core features documentation for details on how to declare an SSL bundle.

Here is an example the configures the `myservice` channel to use the `mybundle` bundle for mutual TLS:

-
Properties
-
YAML

```
spring.grpc.client.channel.myservice.target=static://grpc.example.com:9090
spring.grpc.client.channel.myservice.ssl.bundle=mybundle
```
```
spring:
 grpc:
 client:
 channel:
 myservice:
 target: static://grpc.example.com:9090
 ssl:
 bundle: mybundle
```
| To temporarily disable client SSL support, for example to aid with testing, you can set
 |

### Using In-Process Channels

You can communicate with an in-process server (i.e. not listening on a network port) by including the `io.grpc.grpc-inprocess` dependency on your classpath.

In this mode, the in-process channel factory is auto-configured in addition to the regular channel factories (e.g. Netty). To prevent users from having to deal with multiple channel factories, a composite channel factory is configured as the primary channel factory bean. The composite consults its composed factories to find the first one that supports the channel target.

To use the in-process server the channel target must be set to `in-process:<name>`

| To disable the in-process channel factory, you can set the `spring.grpc.client.inprocess.enabled`property to`false`. |

### Observability

Spring Boot provides auto-configuration of the `ObservationGrpcClientInterceptor` whenever Micrometer is available.
This interceptor provides observability into your gRPC client applications.

| If you use Micrometer, but prefer to not to use it for gRPC, you can set `spring.grpc.client.observation.enabled`to`false`. |

### Channel Customization

If you need to customize your gRPC channel beyond the basic properties, you can use a `GrpcChannelBuilderCustomizer`.
Each customizer is called with the logical target name and the `ManagedChannelBuilder` that will build the channel.
There’s also a convenient `matching(String pattern)` factory method that will limit customizations to targets that match the given regex pattern.

A common customizer use-case is to add security interceptors to the builder.
For example, here we’re adding the `BearerTokenAuthenticationInterceptor` to the target matching “hello”:

-
Java
-
Kotlin

```
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.grpc.client.GrpcChannelBuilderCustomizer;
import org.springframework.grpc.client.interceptor.security.BasicAuthenticationInterceptor;
@Configuration(proxyBeanMethods = false)
public class MyGrpcConfiguration {
	@Bean
	GrpcChannelBuilderCustomizer<?> helloChannelCustomizer() {
 return GrpcChannelBuilderCustomizer.matching("hello",
 (builder) -> builder.intercept(new BasicAuthenticationInterceptor("user", "password")));
	}
}
```
```
import io.grpc.ManagedChannelBuilder
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.grpc.client.GrpcChannelBuilderCustomizer
import org.springframework.grpc.client.interceptor.security.BasicAuthenticationInterceptor
import java.util.function.Consumer
@Configuration(proxyBeanMethods = false)
class MyGrpcConfiguration {
	@Bean
	fun helloChannelCustomizer(): GrpcChannelBuilderCustomizer<*> {
 return GrpcChannelBuilderCustomizer.matching("hello", { builder ->
 builder.intercept(BasicAuthenticationInterceptor("user", "password"))
 })
	}
}
```
## Testing gRPC Applications

To help test your gRPC client and server applications you can use the `spring-boot-grpc-test` module or the `spring-boot-starter-grpc-client-test` / `spring-boot-starter-grpc-server-test` starter POMs.

### Using In-Process Test Transport

The `@AutoConfigureTestGrpcTransport` annotation allows you to quickly replace gRPC communication channels with in-process channels specifically designed for testing.
Unlike regular in-process channels, these test channels to not require any configuration.

Using test gRPC transport means that you don’t need to actually listen on a network port to start your application. This allows your tests to run quickly, whilst still ensuring that your application works as expected.

By default, using `@AutoConfigureTestGrpcTransport` will:

-
Configure test `GrpcServerFactory`/`GrpcChannelFactory`beans
-
Disable any gRPC servelt registration.
-
Disable `GrpcServerFactory`bean auto-configuration.
-
Disable `GrpcChannelFactory`bean auto-configuration.

The following example shows how you can use `@AutoConfigureTestGrpcTransport` to test a gRPC server application:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.grpc.test.autoconfigure.AutoConfigureTestGrpcTransport;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.grpc.client.ImportGrpcClients;
import static org.assertj.core.api.Assertions.assertThat;
@SpringBootTest
@AutoConfigureTestGrpcTransport
@ImportGrpcClients(types = HelloWorldGrpc.HelloWorldBlockingStub.class)
class MyGrpcTests {
	@Autowired
	private HelloWorldGrpc.HelloWorldBlockingStub helloStub;
	@Test
	void sayHello() {
 HelloRequest request = HelloRequest.newBuilder().setName("Spring").build();
 HelloReply reply = this.helloStub.sayHello(request);
 assertThat(reply.getMessage()).isEqualTo("Hello 'Spring'");
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.jooq.DSLContext
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.grpc.test.autoconfigure.AutoConfigureTestGrpcTransport
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.grpc.client.ImportGrpcClients
@SpringBootTest
@AutoConfigureTestGrpcTransport
@ImportGrpcClients(types = [HelloWorldGrpc.HelloWorldBlockingStub::class])
class MyGrpcTests(@Autowired val helloStub: HelloWorldGrpc.HelloWorldBlockingStub) {
	@Test
	fun sayHello() {
 val request = HelloRequest.newBuilder().setName("Spring").build()
 val reply = helloStub.sayHello(request)
 assertThat(reply.getMessage()).isEqualTo("Hello 'Spring'")
	}
}
```
### Testing With a Running Server

If you prefer to test your gRPC application by starting the real server and using the actual network connection, we recommend that you use random ports. This will ensure that you can run your tests in any environment, and that you won’t accidentally call real services.

To start a gRPC server using a random port, set `spring.grpc.server.port` to `0`.
You can use the `@LocalGrpcServerPort` annotation to obtain the actual port that the server started on.

Here’s an example test:

-
Java
-
Kotlin

```
import io.grpc.ManagedChannel;
import io.grpc.netty.NettyChannelBuilder;
import org.junit.jupiter.api.Test;
import org.springframework.boot.docs.io.grpc.testing.localserverport.HelloWorldGrpc.HelloWorldBlockingStub;
import org.springframework.boot.grpc.test.autoconfigure.LocalGrpcServerPort;
import org.springframework.boot.test.context.SpringBootTest;
import static org.assertj.core.api.Assertions.assertThat;
@SpringBootTest(properties = "spring.grpc.server.port=0")
class MyGrpcIntegrationTests {
	@LocalGrpcServerPort
	private int port;
	@Test
	void sayHello() {
 String target = "localhost:%s".formatted(this.port);
 ManagedChannel channel = NettyChannelBuilder.forTarget(target).usePlaintext().build();
 try {
 HelloWorldBlockingStub hello = HelloWorldGrpc.newBlockingStub(channel);
 HelloRequest request = HelloRequest.newBuilder().setName("Spring").build();
 assertThat(hello.sayHello(request).getMessage()).isEqualTo("Hello 'Spring'");
 }
 finally {
 channel.shutdown();
 }
	}
}
```
```
import io.grpc.ManagedChannel
import io.grpc.netty.NettyChannelBuilder
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.boot.grpc.test.autoconfigure.LocalGrpcServerPort
import org.springframework.boot.test.context.SpringBootTest
@SpringBootTest(properties = ["spring.grpc.server.port=0"])
class MyGrpcIntegrationTests {
	@LocalGrpcServerPort
	var port = 0
	@Test
	fun sayHello() {
 val target = "localhost:${port}"
 val channel: ManagedChannel = NettyChannelBuilder.forTarget(target).usePlaintext().build()
 try {
 val hello: HelloWorldGrpc.HelloWorldBlockingStub = HelloWorldGrpc.newBlockingStub(channel)
 val request = HelloRequest.newBuilder().setName("Spring").build()
 assertThat(hello.sayHello(request).getMessage()).isEqualTo("Hello 'Spring'")
 } finally {
 channel.shutdown()
 }
	}
}
```

# Hazelcast

If Hazelcast is on the classpath and a suitable configuration is found, Spring Boot auto-configures a `HazelcastInstance` that you can inject in your application.

Spring Boot first attempts to create a client by checking the following configuration options:

-
The presence of a `ClientConfig`bean.
-
A configuration file defined by the `spring.hazelcast.config`property.
-
The presence of the `hazelcast.client.config`system property.
-
A `hazelcast-client.xml`in the working directory or at the root of the classpath.
-
A `hazelcast-client.yaml`(or`hazelcast-client.yml`) in the working directory or at the root of the classpath.

If a client can not be created, Spring Boot attempts to configure an embedded server.
If you define a `Config` bean, Spring Boot uses it at is.
If your configuration defines an instance name, Spring Boot tries to locate an existing instance rather than creating a new one.

You could also specify the Hazelcast configuration file to use through configuration, as shown in the following example:

-
Properties
-
YAML

`spring.hazelcast.config=classpath:config/my-hazelcast.xml````
spring:
 hazelcast:
 config: "classpath:config/my-hazelcast.xml"
```
Otherwise, Spring Boot tries to find the Hazelcast configuration from the default locations: `hazelcast.xml` in the working directory or at the root of the classpath, or a YAML counterpart in the same locations.
We also check if the `hazelcast.config` system property is set.
See the Hazelcast documentation for more details.

| By default, `@SpringAware`on Hazelcast components is supported.
The`ManagedContext`can be overridden by declaring a`HazelcastConfigCustomizer`bean with an`@Order`higher than zero. |

| Spring Boot also has explicit caching support for Hazelcast.
If caching is enabled, the `HazelcastInstance`is automatically wrapped in a`CacheManager`implementation. |

# Quartz Scheduler

Spring Boot offers several conveniences for working with the Quartz scheduler, including the `spring-boot-starter-quartz` starter.
If Quartz is available, a `Scheduler` is auto-configured (through the `SchedulerFactoryBean` abstraction).

Beans of the following types are automatically picked up and associated with the `Scheduler`:

-
`JobDetail`: defines a particular Job.`JobDetail`instances can be built with the`JobBuilder`API.
-
`Trigger`: defines when a particular job is triggered.

By default, an in-memory `JobStore` is used.
However, it is possible to configure a JDBC-based store if a `DataSource` bean is available in your application and if the `spring.quartz.job-store-type` property is configured accordingly, as shown in the following example:

-
Properties
-
YAML

`spring.quartz.job-store-type=jdbc````
spring:
 quartz:
 job-store-type: "jdbc"
```
When the JDBC store is used, the schema can be initialized on startup, as shown in the following example:

-
Properties
-
YAML

`spring.quartz.jdbc.initialize-schema=always````
spring:
 quartz:
 jdbc:
 initialize-schema: "always"
```
| By default, the database is detected and initialized by using the standard scripts provided with the Quartz library.
These scripts drop existing tables, deleting all triggers on every restart.
To use a custom script, set the `spring.quartz.jdbc.schema`property.
Some of the standard scripts – such as those for SQL Server, Azure SQL, and Sybase – cannot be used without modification.
In these cases, make a copy of the script and edit it as directed in the script’s comments then set`spring.quartz.jdbc.schema`to use your customized script. |

To have Quartz use a `DataSource` other than the application’s main `DataSource`, declare a `DataSource` bean, annotating its `@Bean` method with `@QuartzDataSource`.
Doing so ensures that the Quartz-specific `DataSource` is used by both the `SchedulerFactoryBean` and for schema initialization.
Similarly, to have Quartz use a `TransactionManager` other than the application’s main `TransactionManager` declare a `TransactionManager` bean, annotating its `@Bean` method with `@QuartzTransactionManager`.

By default, jobs created by configuration will not overwrite already registered jobs that have been read from a persistent job store.
To enable overwriting existing job definitions set the `spring.quartz.overwrite-existing-jobs` property.

Quartz Scheduler configuration can be customized using `spring.quartz` properties and `SchedulerFactoryBeanCustomizer` beans, which allow programmatic `SchedulerFactoryBean` customization.
Advanced Quartz configuration properties can be customized using `spring.quartz.properties.*`.

| In particular, an `Executor`bean is not associated with the scheduler as Quartz offers a way to configure the scheduler through`spring.quartz.properties`.
If you need to customize the task executor, consider implementing`SchedulerFactoryBeanCustomizer`. |

Jobs can define setters to inject data map properties. Regular beans can also be injected in a similar manner, as shown in the following example:

-
Java
-
Kotlin

```
import org.quartz.JobExecutionContext;
import org.quartz.JobExecutionException;
import org.springframework.scheduling.quartz.QuartzJobBean;
public class MySampleJob extends QuartzJobBean {
	// fields ...
	private MyService myService;
	private String name;
	// Inject "MyService" bean
	public void setMyService(MyService myService) {
 this.myService = myService;
	}
	// Inject the "name" job data property
	public void setName(String name) {
 this.name = name;
	}
	@Override
	protected void executeInternal(JobExecutionContext context) throws JobExecutionException {
 this.myService.someMethod(context.getFireTime(), this.name);
	}
}
```
```
import org.quartz.JobExecutionContext
import org.springframework.scheduling.quartz.QuartzJobBean
class MySampleJob : QuartzJobBean() {
	// fields ...
	private var myService: MyService? = null
	private var name: String? = null
	// Inject "MyService" bean
	fun setMyService(myService: MyService?) {
 this.myService = myService
	}
	// Inject the "name" job data property
	fun setName(name: String?) {
 this.name = name
	}
	override fun executeInternal(context: JobExecutionContext) {
 myService!!.someMethod(context.fireTime, name)
	}
}
```

# Sending Email

The Spring Framework provides an abstraction for sending email by using the `JavaMailSender` interface, and Spring Boot provides auto-configuration for it as well as a starter module.

| See the reference documentation for a detailed explanation of how you can use `JavaMailSender`. |

If `spring.mail.host` and the relevant libraries (as defined by `spring-boot-starter-mail`) are available, a default `JavaMailSender` is created if none exists.
The sender can be further customized by configuration items from the `spring.mail` namespace.
See `MailProperties` for more details.

In particular, certain default timeout values are infinite, and you may want to change that to avoid having a thread blocked by an unresponsive mail server, as shown in the following example:

-
Properties
-
YAML

```
spring.mail.properties[mail.smtp.connectiontimeout]=5000
spring.mail.properties[mail.smtp.timeout]=3000
spring.mail.properties[mail.smtp.writetimeout]=5000
```
```
spring:
 mail:
 properties:
 "[mail.smtp.connectiontimeout]": 5000
 "[mail.smtp.timeout]": 3000
 "[mail.smtp.writetimeout]": 5000
```
It is also possible to configure a `JavaMailSender` with an existing `Session` from JNDI:

-
Properties
-
YAML

`spring.mail.jndi-name=mail/Session````
spring:
 mail:
 jndi-name: "mail/Session"
```
When a `jndi-name` is set, it takes precedence over all other Session-related settings.

# Validation

The method validation feature supported by Bean Validation 1.1 is automatically enabled as long as a JSR-303 implementation (such as Hibernate Validator, typically provided by `spring-boot-starter-validation`) is on the classpath.
This lets bean methods be annotated with `jakarta.validation` constraints on their parameters and/or on their return value.
Target classes with such annotated methods need to be annotated with the `@Validated` annotation at the type level for their methods to be searched for inline constraint annotations.

For instance, the following service triggers the validation of the first argument, making sure its size is between 8 and 10:

-
Java
-
Kotlin

```
import jakarta.validation.constraints.Size;
import org.springframework.stereotype.Service;
import org.springframework.validation.annotation.Validated;
@Service
@Validated
public class MyBean {
	public Archive findByCodeAndAuthor(@Size(min = 8, max = 10) String code, Author author) {
 return ...
	}
}
```
```
import jakarta.validation.constraints.Size
import org.springframework.stereotype.Service
import org.springframework.validation.annotation.Validated
@Service
@Validated
class MyBean {
	fun findByCodeAndAuthor(code: @Size(min = 8, max = 10) String?, author: Author?): Archive? {
 return null
	}
}
```
The application’s `MessageSource` is used when resolving `{parameters}` in constraint messages.
This allows you to use your application’s `messages.properties` files for Bean Validation messages.
Once the parameters have been resolved, message interpolation is completed using Bean Validation’s default interpolator.

To customize the `Configuration` used to build the `ValidatorFactory`, define a `ValidationConfigurationCustomizer` bean.
When multiple customizer beans are defined, they are called in order based on their `@Order` annotation or `Ordered` implementation.

# Calling REST Services

Spring Boot provides various convenient ways to call remote REST services.
If you are developing a non-blocking reactive application and you’re using Spring WebFlux, then you can use `WebClient`.
If you prefer imperative APIs then you can use `RestClient` or `RestTemplate`.

## WebClient

If you have Spring WebFlux on your classpath we recommend that you use `WebClient` to call remote REST services.
The `WebClient` interface provides a functional style API and is fully reactive.
You can learn more about the `WebClient` in the dedicated section in the Spring Framework docs.

| If you are not writing a reactive Spring WebFlux application you can use the `RestClient`instead of a`WebClient`.
This provides a similar functional API, but is imperative rather than reactive. |

Spring Boot creates and pre-configures a prototype `WebClient.Builder` bean for you.
It is strongly advised to inject it in your components and use it to create `WebClient` instances.
Spring Boot is configuring that builder to share HTTP resources and reflect codecs setup in the same fashion as the server ones (see WebFlux HTTP codecs auto-configuration), and more.

The following code shows a typical example:

-
Java
-
Kotlin

```
import reactor.core.publisher.Mono;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
@Service
public class MyService {
	private final WebClient webClient;
	public MyService(WebClient.Builder webClientBuilder) {
 this.webClient = webClientBuilder.baseUrl("https://example.org").build();
	}
	public Mono<Details> someRestCall(String name) {
 return this.webClient.get().uri("/{name}/details", name).retrieve().bodyToMono(Details.class);
	}
}
```
```
import org.springframework.stereotype.Service
import org.springframework.web.reactive.function.client.WebClient
import reactor.core.publisher.Mono
@Service
class MyService(webClientBuilder: WebClient.Builder) {
	private val webClient: WebClient
	init {
 webClient = webClientBuilder.baseUrl("https://example.org").build()
	}
	fun someRestCall(name: String): Mono<Details> {
 return webClient.get().uri("/{name}/details", name)
 .retrieve().bodyToMono(Details::class.java)
	}
}
```
### WebClient Runtime

Spring Boot will auto-detect which `ClientHttpConnector` to use to drive `WebClient` depending on the libraries available on the application classpath.
In order of preference, the following clients are supported:

-
Reactor Netty
-
Jetty RS client
-
Apache HttpClient
-
JDK HttpClient

If multiple clients are available on the classpath, the most preferred client will be used.

The `spring-boot-starter-webflux` starter depends on `io.projectreactor.netty:reactor-netty` by default, which brings both server and client implementations.
If you choose to use Jetty as a reactive server instead, you should add a dependency on the Jetty Reactive HTTP client library, `org.eclipse.jetty:jetty-reactive-httpclient`.
Using the same technology for server and client has its advantages, as it will automatically share HTTP resources between client and server.

Developers can override the resource configuration for Jetty and Reactor Netty by providing a custom `ReactorResourceFactory` or `JettyResourceFactory` bean - this will be applied to both clients and servers.

If you wish to override that choice for the client, you can define your own `ClientHttpConnector` bean and have full control over the client configuration.

You can learn more about the `WebClient` configuration options in the Spring Framework reference documentation.

### Global HTTP Connector Configuration

If the auto-detected `ClientHttpConnector` does not meet your needs, you can use the `spring.http.clients.reactive.connector` property to pick a specific connector.
For example, if you have Reactor Netty on your classpath, but you prefer Jetty’s `HttpClient` you can add the following:

-
Properties
-
YAML

`spring.http.clients.reactive.connector=jetty````
spring:
 http:
 clients:
 reactive:
 connector: jetty
```
| You can also use global configuration properties which apply to all HTTP clients. |

For more complex customizations, you can use `ClientHttpConnectorBuilderCustomizer` or declare your own `ClientHttpConnectorBuilder` bean which will cause auto-configuration to back off.
This can be useful when you need to customize some of the internals of the underlying HTTP library.

For example, the following will use a JDK client configured with a specific `ProxySelector`:

-
Java
-
Kotlin

```
import java.net.ProxySelector;
import org.springframework.boot.http.client.reactive.ClientHttpConnectorBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyConnectorHttpConfiguration {
	@Bean
	ClientHttpConnectorBuilder<?> clientHttpConnectorBuilder(ProxySelector proxySelector) {
 return ClientHttpConnectorBuilder.jdk().withHttpClientCustomizer((builder) -> builder.proxy(proxySelector));
	}
}
```
```
import org.springframework.boot.http.client.reactive.ClientHttpConnectorBuilder
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import java.net.ProxySelector
@Configuration(proxyBeanMethods = false)
class MyConnectorHttpConfiguration {
	@Bean
	fun clientHttpConnectorBuilder(proxySelector: ProxySelector): ClientHttpConnectorBuilder<*> {
 return ClientHttpConnectorBuilder.jdk().withHttpClientCustomizer { builder -> builder.proxy(proxySelector) }
	}
}
```
### WebClient Customization

There are three main approaches to `WebClient` customization, depending on how broadly you want the customizations to apply.

To make the scope of any customizations as narrow as possible, inject the auto-configured `WebClient.Builder` and then call its methods as required.
`WebClient.Builder` instances are stateful: Any change on the builder is reflected in all clients subsequently created with it.
If you want to create several clients with the same builder, you can also consider cloning the builder with `WebClient.Builder other = builder.clone();`.

To make an application-wide, additive customization to all `WebClient.Builder` instances, you can declare `WebClientCustomizer` beans and change the `WebClient.Builder` locally at the point of injection.

Finally, you can fall back to the original API and use `WebClient.create()`.
In that case, no auto-configuration or `WebClientCustomizer` is applied.

### WebClient SSL Support

If you need custom SSL configuration on the `ClientHttpConnector` used by the `WebClient`, you can inject a `WebClientSsl` instance that can be used with the builder’s `apply` method.

The `WebClientSsl` interface provides access to any SSL bundles that you have defined in your `application.properties` or `application.yaml` file.

The following code shows a typical example:

-
Java
-
Kotlin

```
import reactor.core.publisher.Mono;
import org.springframework.boot.webclient.autoconfigure.WebClientSsl;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
@Service
public class MyService {
	private final WebClient webClient;
	public MyService(WebClient.Builder webClientBuilder, WebClientSsl ssl) {
 this.webClient = webClientBuilder.baseUrl("https://example.org").apply(ssl.fromBundle("mybundle")).build();
	}
	public Mono<Details> someRestCall(String name) {
 return this.webClient.get().uri("/{name}/details", name).retrieve().bodyToMono(Details.class);
	}
}
```
```
import org.springframework.boot.webclient.autoconfigure.WebClientSsl
import org.springframework.stereotype.Service
import org.springframework.web.reactive.function.client.WebClient
import reactor.core.publisher.Mono
@Service
class MyService(webClientBuilder: WebClient.Builder, ssl: WebClientSsl) {
	private val webClient: WebClient
	init {
 webClient = webClientBuilder.baseUrl("https://example.org")
 .apply(ssl.fromBundle("mybundle")).build()
	}
	fun someRestCall(name: String): Mono<Details> {
 return webClient.get().uri("/{name}/details", name)
 .retrieve().bodyToMono(Details::class.java)
	}
}
```
## RestClient

If you are not using Spring WebFlux or Project Reactor in your application we recommend that you use `RestClient` to call remote REST services.

The `RestClient` interface provides a functional style imperative API.

Spring Boot creates and pre-configures a prototype `RestClient.Builder` bean for you.
It is strongly advised to inject it in your components and use it to create `RestClient` instances.
Spring Boot is configuring that builder with `HttpMessageConverters` and an appropriate `ClientHttpRequestFactory`.

The following code shows a typical example:

-
Java
-
Kotlin

```
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;
@Service
public class MyService {
	private final RestClient restClient;
	public MyService(RestClient.Builder restClientBuilder) {
 this.restClient = restClientBuilder.baseUrl("https://example.org").build();
	}
	public Details someRestCall(String name) {
 return this.restClient.get().uri("/{name}/details", name).retrieve().body(Details.class);
	}
}
```
```
import org.springframework.boot.docs.io.restclient.restclient.ssl.Details
import org.springframework.stereotype.Service
import org.springframework.web.client.RestClient
@Service
class MyService(restClientBuilder: RestClient.Builder) {
	private val restClient: RestClient
	init {
 restClient = restClientBuilder.baseUrl("https://example.org").build()
	}
	fun someRestCall(name: String): Details {
 return restClient.get().uri("/{name}/details", name)
 .retrieve().body(Details::class.java)!!
	}
}
```
### RestClient Customization

There are three main approaches to `RestClient` customization, depending on how broadly you want the customizations to apply.

To make the scope of any customizations as narrow as possible, inject the auto-configured `RestClient.Builder` and then call its methods as required.
`RestClient.Builder` instances are stateful: Any change on the builder is reflected in all clients subsequently created with it.
If you want to create several clients with the same builder, you can also consider cloning the builder with `RestClient.Builder other = builder.clone();`.

To make an application-wide, additive customization to all `RestClient.Builder` instances, you can declare `RestClientCustomizer` beans and change the `RestClient.Builder` locally at the point of injection.

Finally, you can fall back to the original API and use `RestClient.create()`.
In that case, no auto-configuration or `RestClientCustomizer` is applied.

| You can also change the global HTTP client configuration. |

### RestClient SSL Support

If you need custom SSL configuration on the `ClientHttpRequestFactory` used by the `RestClient`, you can inject a `RestClientSsl` instance that can be used with the builder’s `apply` method.

The `RestClientSsl` interface provides access to any SSL bundles that you have defined in your `application.properties` or `application.yaml` file.

The following code shows a typical example:

-
Java
-
Kotlin

```
import org.springframework.boot.restclient.autoconfigure.RestClientSsl;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;
@Service
public class MyService {
	private final RestClient restClient;
	public MyService(RestClient.Builder restClientBuilder, RestClientSsl ssl) {
 this.restClient = restClientBuilder.baseUrl("https://example.org").apply(ssl.fromBundle("mybundle")).build();
	}
	public Details someRestCall(String name) {
 return this.restClient.get().uri("/{name}/details", name).retrieve().body(Details.class);
	}
}
```
```
import org.springframework.boot.docs.io.restclient.restclient.ssl.settings.Details
import org.springframework.boot.restclient.autoconfigure.RestClientSsl
import org.springframework.stereotype.Service
import org.springframework.web.client.RestClient
@Service
class MyService(restClientBuilder: RestClient.Builder, ssl: RestClientSsl) {
	private val restClient: RestClient
	init {
 restClient = restClientBuilder.baseUrl("https://example.org")
 .apply(ssl.fromBundle("mybundle")).build()
	}
	fun someRestCall(name: String): Details {
 return restClient.get().uri("/{name}/details", name)
 .retrieve().body(Details::class.java)!!
	}
}
```
If you need to apply other customization in addition to an SSL bundle, you can use the `HttpClientSettings` class with `ClientHttpRequestFactoryBuilder`:

-
Java
-
Kotlin

```
import java.time.Duration;
import org.springframework.boot.http.client.ClientHttpRequestFactoryBuilder;
import org.springframework.boot.http.client.HttpClientSettings;
import org.springframework.boot.ssl.SslBundles;
import org.springframework.http.client.ClientHttpRequestFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;
@Service
public class MyService {
	private final RestClient restClient;
	public MyService(RestClient.Builder restClientBuilder, SslBundles sslBundles) {
 HttpClientSettings settings = HttpClientSettings.ofSslBundle(sslBundles.getBundle("mybundle"))
 .withReadTimeout(Duration.ofMinutes(2));
 ClientHttpRequestFactory requestFactory = ClientHttpRequestFactoryBuilder.detect().build(settings);
 this.restClient = restClientBuilder.baseUrl("https://example.org").requestFactory(requestFactory).build();
	}
	public Details someRestCall(String name) {
 return this.restClient.get().uri("/{name}/details", name).retrieve().body(Details.class);
	}
}
```
```
import org.springframework.boot.http.client.ClientHttpRequestFactoryBuilder;
import org.springframework.boot.http.client.HttpClientSettings
import org.springframework.boot.ssl.SslBundles
import org.springframework.stereotype.Service
import org.springframework.web.client.RestClient
import java.time.Duration
@Service
class MyService(restClientBuilder: RestClient.Builder, sslBundles: SslBundles) {
	private val restClient: RestClient
	init {
 val settings = HttpClientSettings.defaults()
 .withReadTimeout(Duration.ofMinutes(2))
 .withSslBundle(sslBundles.getBundle("mybundle"))
 val requestFactory = ClientHttpRequestFactoryBuilder.detect().build(settings);
 restClient = restClientBuilder
 .baseUrl("https://example.org")
 .requestFactory(requestFactory).build()
	}
	fun someRestCall(name: String): Details {
 return restClient.get().uri("/{name}/details", name).retrieve().body(Details::class.java)!!
	}
}
```
## RestTemplate

Spring Framework’s `RestTemplate` class predates `RestClient` and is the classic way that many applications use to call remote REST services.
You might choose to use `RestTemplate` when you have existing code that you don’t want to migrate to `RestClient`, or because you’re already familiar with the `RestTemplate` API.

Since `RestTemplate` instances often need to be customized before being used, Spring Boot does not provide any single auto-configured `RestTemplate` bean.
It does, however, auto-configure a `RestTemplateBuilder`, which can be used to create `RestTemplate` instances when needed.
The auto-configured `RestTemplateBuilder` ensures that sensible `HttpMessageConverters` and an appropriate `ClientHttpRequestFactory` are applied to `RestTemplate` instances.

The following code shows a typical example:

-
Java
-
Kotlin

```
import org.springframework.boot.restclient.RestTemplateBuilder;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;
@Service
public class MyService {
	private final RestTemplate restTemplate;
	public MyService(RestTemplateBuilder restTemplateBuilder) {
 this.restTemplate = restTemplateBuilder.build();
	}
	public Details someRestCall(String name) {
 return this.restTemplate.getForObject("/{name}/details", Details.class, name);
	}
}
```
```
import org.springframework.boot.restclient.RestTemplateBuilder
import org.springframework.stereotype.Service
import org.springframework.web.client.RestTemplate
@Service
class MyService(restTemplateBuilder: RestTemplateBuilder) {
	private val restTemplate: RestTemplate
	init {
 restTemplate = restTemplateBuilder.build()
	}
	fun someRestCall(name: String): Details {
 return restTemplate.getForObject("/{name}/details", Details::class.java, name)!!
	}
}
```
`RestTemplateBuilder` includes a number of useful methods that can be used to quickly configure a `RestTemplate`.
For example, to add BASIC authentication support, you can use `builder.basicAuthentication("user", "password").build()`.

### RestTemplate Customization

There are three main approaches to `RestTemplate` customization, depending on how broadly you want the customizations to apply.

To make the scope of any customizations as narrow as possible, inject the auto-configured `RestTemplateBuilder` and then call its methods as required.
Each method call returns a new `RestTemplateBuilder` instance, so the customizations only affect this use of the builder.

To make an application-wide, additive customization, use a `RestTemplateCustomizer` bean.
All such beans are automatically registered with the auto-configured `RestTemplateBuilder` and are applied to any templates that are built with it.

The following example shows a customizer that configures the use of a proxy for all hosts except `192.168.0.5`:

-
Java
-
Kotlin

```
import org.apache.hc.client5.http.classic.HttpClient;
import org.apache.hc.client5.http.impl.classic.HttpClientBuilder;
import org.apache.hc.client5.http.impl.routing.DefaultProxyRoutePlanner;
import org.apache.hc.client5.http.routing.HttpRoutePlanner;
import org.apache.hc.core5.http.HttpException;
import org.apache.hc.core5.http.HttpHost;
import org.apache.hc.core5.http.protocol.HttpContext;
import org.springframework.boot.restclient.RestTemplateCustomizer;
import org.springframework.http.client.HttpComponentsClientHttpRequestFactory;
import org.springframework.web.client.RestTemplate;
public class MyRestTemplateCustomizer implements RestTemplateCustomizer {
	@Override
	public void customize(RestTemplate restTemplate) {
 HttpRoutePlanner routePlanner = new CustomRoutePlanner(new HttpHost("proxy.example.com"));
 HttpClient httpClient = HttpClientBuilder.create().setRoutePlanner(routePlanner).build();
 restTemplate.setRequestFactory(new HttpComponentsClientHttpRequestFactory(httpClient));
	}
	static class CustomRoutePlanner extends DefaultProxyRoutePlanner {
 CustomRoutePlanner(HttpHost proxy) {
 super(proxy);
 }
 @Override
 protected HttpHost determineProxy(HttpHost target, HttpContext context) throws HttpException {
 if (target.getHostName().equals("192.168.0.5")) {
 return null;
 }
 return super.determineProxy(target, context);
 }
	}
}
```
```
import org.apache.hc.client5.http.classic.HttpClient
import org.apache.hc.client5.http.impl.classic.HttpClientBuilder
import org.apache.hc.client5.http.impl.routing.DefaultProxyRoutePlanner
import org.apache.hc.client5.http.routing.HttpRoutePlanner
import org.apache.hc.core5.http.HttpException
import org.apache.hc.core5.http.HttpHost
import org.apache.hc.core5.http.protocol.HttpContext
import org.springframework.boot.restclient.RestTemplateCustomizer
import org.springframework.http.client.HttpComponentsClientHttpRequestFactory
import org.springframework.web.client.RestTemplate
class MyRestTemplateCustomizer : RestTemplateCustomizer {
	override fun customize(restTemplate: RestTemplate) {
 val routePlanner: HttpRoutePlanner = CustomRoutePlanner(HttpHost("proxy.example.com"))
 val httpClient: HttpClient = HttpClientBuilder.create().setRoutePlanner(routePlanner).build()
 restTemplate.requestFactory = HttpComponentsClientHttpRequestFactory(httpClient)
	}
	internal class CustomRoutePlanner(proxy: HttpHost?) : DefaultProxyRoutePlanner(proxy) {
 @Throws(HttpException::class)
 public override fun determineProxy(target: HttpHost, context: HttpContext): HttpHost? {
 if (target.hostName == "192.168.0.5") {
 return null
 }
 return super.determineProxy(target, context)
 }
	}
}
```
Finally, you can define your own `RestTemplateBuilder` bean.
Doing so will replace the auto-configured builder.
If you want any `RestTemplateCustomizer` beans to be applied to your custom builder, as the auto-configuration would have done, configure it using a `RestTemplateBuilderConfigurer`.
The following example exposes a `RestTemplateBuilder` that matches what Spring Boot’s auto-configuration would have done, except that custom connect and read timeouts are also specified:

-
Java
-
Kotlin

```
import java.time.Duration;
import org.springframework.boot.restclient.RestTemplateBuilder;
import org.springframework.boot.restclient.autoconfigure.RestTemplateBuilderConfigurer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyRestTemplateBuilderConfiguration {
	@Bean
	public RestTemplateBuilder restTemplateBuilder(RestTemplateBuilderConfigurer configurer) {
 return configurer.configure(new RestTemplateBuilder())
 .connectTimeout(Duration.ofSeconds(5))
 .readTimeout(Duration.ofSeconds(2));
	}
}
```
```
import org.springframework.boot.restclient.autoconfigure.RestTemplateBuilderConfigurer
import org.springframework.boot.restclient.RestTemplateBuilder
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import java.time.Duration
@Configuration(proxyBeanMethods = false)
class MyRestTemplateBuilderConfiguration {
	@Bean
	fun restTemplateBuilder(configurer: RestTemplateBuilderConfigurer): RestTemplateBuilder {
 return configurer.configure(RestTemplateBuilder()).connectTimeout(Duration.ofSeconds(5))
 .readTimeout(Duration.ofSeconds(2))
	}
}
```
The most extreme (and rarely used) option is to create your own `RestTemplateBuilder` bean without using a configurer.
In addition to replacing the auto-configured builder, this also prevents any `RestTemplateCustomizer` beans from being used.

| You can also change the global HTTP client configuration. |

### RestTemplate SSL Support

If you need custom SSL configuration on the `RestTemplate`, you can apply an SSL bundle to the `RestTemplateBuilder` as shown in this example:

-
Java
-
Kotlin

```
import org.springframework.boot.docs.io.restclient.resttemplate.Details;
import org.springframework.boot.restclient.RestTemplateBuilder;
import org.springframework.boot.ssl.SslBundles;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;
@Service
public class MyService {
	private final RestTemplate restTemplate;
	public MyService(RestTemplateBuilder restTemplateBuilder, SslBundles sslBundles) {
 this.restTemplate = restTemplateBuilder.sslBundle(sslBundles.getBundle("mybundle")).build();
	}
	public Details someRestCall(String name) {
 return this.restTemplate.getForObject("/{name}/details", Details.class, name);
	}
}
```
```
import org.springframework.boot.docs.io.restclient.resttemplate.Details
import org.springframework.boot.ssl.SslBundles
import org.springframework.boot.restclient.RestTemplateBuilder
import org.springframework.stereotype.Service
import org.springframework.web.client.RestTemplate
@Service
class MyService(restTemplateBuilder: RestTemplateBuilder, sslBundles: SslBundles) {
 private val restTemplate: RestTemplate
 init {
 restTemplate = restTemplateBuilder.sslBundle(sslBundles.getBundle("mybundle")).build()
 }
 fun someRestCall(name: String): Details {
 return restTemplate.getForObject("/{name}/details", Details::class.java, name)!!
 }
}
```
## HTTP Client Detection for RestClient and RestTemplate

Spring Boot will auto-detect which HTTP client to use with `RestClient` and `RestTemplate` depending on the libraries available on the application classpath.
In order of preference, the following clients are supported:

-
Apache HttpClient
-
Jetty HttpClient
-
Reactor Netty HttpClient
-
JDK client ( `java.net.http.HttpClient`)
-
Simple JDK client ( `java.net.HttpURLConnection`)

If multiple clients are available on the classpath, and no global configuration is provided, the most preferred client will be used.

### Global HTTP Client Configuration

If the auto-detected HTTP client does not meet your needs, you can use the `spring.http.clients.imperative.factory` property to pick a specific factory.
For example, if you have Apache HttpClient on your classpath, but you prefer Jetty’s `HttpClient` you can add the following:

-
Properties
-
YAML

`spring.http.clients.imperative.factory=jetty````
spring:
 http:
 clients:
 imperative:
 factory: jetty
```
| You can also use global configuration properties which apply to all HTTP clients. |

For more complex customizations, you can use `ClientHttpRequestFactoryBuilderCustomizer` or declare your own `ClientHttpRequestFactoryBuilder` bean which will cause auto-configuration to back off.
This can be useful when you need to customize some of the internals of the underlying HTTP library.

For example, the following will use a JDK client configured with a specific `ProxySelector`:

-
Java
-
Kotlin

```
import java.net.ProxySelector;
import org.springframework.boot.http.client.ClientHttpRequestFactoryBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyClientHttpConfiguration {
	@Bean
	ClientHttpRequestFactoryBuilder<?> clientHttpRequestFactoryBuilder(ProxySelector proxySelector) {
 return ClientHttpRequestFactoryBuilder.jdk()
 .withHttpClientCustomizer((builder) -> builder.proxy(proxySelector));
	}
}
```
```
import org.springframework.boot.http.client.ClientHttpRequestFactoryBuilder
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import java.net.ProxySelector
@Configuration(proxyBeanMethods = false)
class MyClientHttpConfiguration {
	@Bean
	fun clientHttpRequestFactoryBuilder(proxySelector: ProxySelector): ClientHttpRequestFactoryBuilder<*> {
 return ClientHttpRequestFactoryBuilder.jdk()
 .withHttpClientCustomizer { builder -> builder.proxy(proxySelector) }
	}
}
```
## API Versioning

Both `WebClient` and `RestClient` support making versioned remote HTTP calls so that APIs can be evolved over time.
Commonly this involves sending an HTTP header, a query parameter or URL path segment that indicates the version of the API that should be used.

You can configure API versioning using methods on `WebClient.Builder` or `RestClient.Builder`.

| API versioning is also supported on the server-side. See the Spring MVC and Spring WebFlux sections for details. |

| The server-side API versioning configuration is not taken into account to auto-configure the client. Clients that should use an API versioning strategy, typically for testing, need to configure it explicitly. |

## HTTP Service Interface Clients

Instead of directly using a `RestClient` or `WebClient` to call an HTTP service, it’s also possible to call them using annotated Java interfaces.

HTTP Service interfaces defines a service contract by using methods that are annotated with `@HttpExchange`, or more typically the method specific variants (`@GetExchange`, `@PostExchange`, `@DeleteExchange`, etc).

For example, the following code defines an HTTP Service for an “echo” API that will return a JSON object containing an echo of the request.

-
Java
-
Kotlin

```
import java.util.Map;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.service.annotation.HttpExchange;
import org.springframework.web.service.annotation.PostExchange;
@HttpExchange(url = "https://echo.zuplo.io")
public interface EchoService {
	@PostExchange
	Map<?, ?> echo(@RequestBody Map<String, String> message);
}
```
```
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.service.annotation.HttpExchange
import org.springframework.web.service.annotation.PostExchange
@HttpExchange(url = "https://echo.zuplo.io")
interface EchoService {
	@PostExchange
	fun echo(@RequestBody message: Map<String, String>): Map<*, *>
}
```
More details about how to develop HTTP Service interface clients can be found in the Spring Framework reference documentation.

### Importing HTTP Services

In order to use an HTTP Service interface as client you need to import it.
One way to achieve this is to use the `@ImportHttpServices` annotation, typically on your main application class.
You can use the annotation to import specific classes, or scan for classes to import from specific packages.

For example, the following configuration will scan for HTTP Service interfaces in the `com.example.myclients` package:

-
Java
-
Kotlin

```
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.service.registry.ImportHttpServices;
@SpringBootApplication
@ImportHttpServices(basePackages = "com.example.myclients")
public class MyApplication {
	public static void main(String[] args) {
 SpringApplication.run(MyApplication.class, args);
	}
}
```
```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.web.service.registry.ImportHttpServices
@SpringBootApplication
@ImportHttpServices(basePackages = ["com.example.myclients"])
class MyApplication
fun main(args: Array<String>) {
	runApplication<MyApplication>(*args)
}
```
### Service Client Groups

Hard-coding absolute URLs in `@HttpExchange` annotations is often not ideal in production applications.
Instead, you will typically want to give the HTTP Service client a logical name in your code, and then lookup a URL from a property based on that name.

HTTP Service clients allow you to do this by registering them into named groups. An HTTP Service group is a collection of HTTP Service interfaces that all share common features.

For example, we may want to define an “echo” group to use for HTTP Service clients that call `https://echo.zuplo.io`.

| HTTP Service groups can be used to define more than just URLs. For example, your group could define connection timeouts and SSL settings. You can also associate client customization logic to a group, such as adding code to insert required authorization headers. |

To associate an HTTP Service interface with a group when using `@ImportHttpServices` you can use the `group` attribute.

For example, if we assume our example above is organized in such a way that all HTTP Service interfaces in the `com.example.myclients` package belong to the `echo` group.
We first remove the hardcoded URL from the service interface:

-
Java
-
Kotlin

```
import java.util.Map;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.service.annotation.PostExchange;
public interface EchoService {
	@PostExchange
	Map<?, ?> echo(@RequestBody Map<String, String> message);
}
```
```
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.service.annotation.HttpExchange
import org.springframework.web.service.annotation.PostExchange
interface EchoService {
	@PostExchange
	fun echo(@RequestBody message: Map<String, String>): Map<*, *>
}
```
We can then write:

-
Java
-
Kotlin

```
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.service.registry.ImportHttpServices;
@SpringBootApplication
@ImportHttpServices(group = "echo", basePackages = "com.example.myclients")
public class MyApplication {
	public static void main(String[] args) {
 SpringApplication.run(MyApplication.class, args);
	}
}
```
```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.web.service.registry.ImportHttpServices
@SpringBootApplication
@ImportHttpServices(group = "echo", basePackages = ["com.example.myclients"])
class MyApplication
fun main(args: Array<String>) {
	runApplication<MyApplication>(*args)
}
```
And finally we can then use a `base-url` property to link the `echo` group to an actual URL:

-
Properties
-
YAML

`spring.http.serviceclient.echo.base-url=https://echo.zuplo.io````
spring:
 http:
 serviceclient:
 echo:
 base-url: "https://echo.zuplo.io"
```
| HTTP Service clients will be associated with a group named “default” if you don’t specify a group. |

| If you have multiple HTTP Service interfaces in the same package that need to be associated with different groups you can list them individually.
The For example:
 |

### Configuration Properties

Configuration properties for HTTP Services can be specified under `spring.http.serviceclient.<group-name>`:

You can use properties to configure aspects such as:

-
The base URL.
-
Any default headers that should be sent.
-
API versioning configuration.
-
Redirect settings.
-
Connection and read timeouts.
-
SSL bundles to use.

| You can also use global configuration properties which apply to all HTTP clients. |

For example, the properties below will:

-
Configure all HTTP clients to use a one second connect timeout (unless otherwise overridden).
-
Configure HTTP Service clients in the “echo” group to: -
Use a specific base URL.
-
Have a one-second connect timeout.
-
Have a two-second read timeout.

-

-
Properties
-
YAML

```
spring.http.clients.connect-timeout=1s
spring.http.serviceclient.echo.base-url=https://echo.zuplo.io
spring.http.serviceclient.echo.connect-timeout=2s
spring.http.serviceclient.echo.read-timeout=2s
```
```
spring:
 http:
 clients:
 connect-timeout: 1s
 serviceclient:
 echo:
 base-url: "https://echo.zuplo.io"
 connect-timeout: 2s
 read-timeout: 2s
```
### Customization

If you need to customize HTTP Service clients beyond basic properties, you can use an HTTP Service group configurer.
For `RestClient` backed HTTP Service clients, you can declare a bean that implements `RestClientHttpServiceGroupConfigurer`.
For `WebClient` backed HTTP Service clients you can declare a bean that implements `WebClientHttpServiceGroupConfigurer`.

Both work in the same way and will be automatically applied by Spring Boot’s auto-configuration.

For example, the following configuration would add a group customizer that adds an HTTP header to each outgoing request containing the group name:

-
Java
-
Kotlin

```
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.client.support.RestClientHttpServiceGroupConfigurer;
@Configuration(proxyBeanMethods = false)
public class MyHttpServiceGroupConfiguration {
	@Bean
	RestClientHttpServiceGroupConfigurer myHttpServiceGroupConfigurer() {
 return (groups) -> groups.forEachClient((group, clientBuilder) -> {
 String groupName = group.name();
 clientBuilder.defaultHeader("service-group", groupName);
 });
	}
}
```
```
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.web.client.support.RestClientHttpServiceGroupConfigurer
@Configuration(proxyBeanMethods = false)
class MyHttpServiceGroupConfiguration {
	@Bean
	fun myHttpServiceGroupConfigurer(): RestClientHttpServiceGroupConfigurer {
 return RestClientHttpServiceGroupConfigurer { groups ->
 groups.forEachClient { group, clientBuilder ->
 val groupName = group.name()
 clientBuilder.defaultHeader("service-group", groupName)
 }
 }
	}
}
```
### Advanced Configuration

As well as the `@ImportHttpServices` annotation, Spring Framework also offers an `AbstractHttpServiceRegistrar` class.
You can `@Import` your own extension of this class to perform programmatic configuration.
For more details, see Spring Framework reference documentation.

Regardless of which method you use to register HTTP Service clients, Spring Boot’s support remains the same.

## Applying Global Configuration to All HTTP Clients

Regardless of the underlying technology being used, all HTTP clients have common settings that can be configured.

These include:

-
Connection Timeouts.
-
Read Timeouts.
-
How HTTP redirects should be handled.
-
Which SSL bundle should be used when connecting.

These common settings are represented by the `HttpClientSettings` class which can be passed into the `build(…)` methods of `ClientHttpConnectorBuilder` and `ClientHttpRequestFactoryBuilder`.

If you want to apply the same configuration to all auto-configured clients, you can use `spring.http.clients` properties to do so:

-
Properties
-
YAML

```
spring.http.clients.connect-timeout=2s
spring.http.clients.read-timeout=1s
spring.http.clients.redirects=dont-follow
```
```
spring:
 http:
 clients:
 connect-timeout: 2s
 read-timeout: 1s
 redirects: dont-follow
```
### InetAddress Filtering and SSRF Protection

It’s sometimes useful to limit the remote addresses that an HTTP client is permitted to call. This technique can be especially useful when hardening your application against Server-Side Request Forgery (SSRF) attacks.

To limit the address that a client can call, you can use an `InetAddressFilter` which will only allow outgoing calls to addresses that match the filter.
The filter is a functional interface and you can either create your own implementation, or use one of the convenient factory methods.

Filters may be applied to the `HttpClientSettings` you used when building a client:

-
Java
-
Kotlin

```
import org.springframework.boot.http.client.ClientHttpRequestFactoryBuilder;
import org.springframework.boot.http.client.HttpClientSettings;
import org.springframework.boot.http.client.InetAddressFilter;
import org.springframework.http.client.ClientHttpRequestFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;
@Service
public class MyService {
	private final RestClient restClient;
	public MyService() {
 InetAddressFilter onlyExternalAddresses = InetAddressFilter.externalAddresses();
 HttpClientSettings settings = HttpClientSettings.defaults().withInetAddressFilter(onlyExternalAddresses);
 ClientHttpRequestFactory requestFactory = ClientHttpRequestFactoryBuilder.jdk().build(settings);
 this.restClient = RestClient.builder().requestFactory(requestFactory).baseUrl("https://example.org").build();
	}
	public Details someRestCall(String name) {
 return this.restClient.get().uri("/{name}/details", name).retrieve().body(Details.class);
	}
}
```
```
import org.springframework.boot.http.client.ClientHttpRequestFactoryBuilder
import org.springframework.boot.http.client.HttpClientSettings
import org.springframework.boot.http.client.InetAddressFilter
import org.springframework.http.client.ClientHttpRequestFactory
import org.springframework.stereotype.Service
import org.springframework.web.client.RestClient
@Service
class MyService {
	private val restClient: RestClient
	init {
 val onlyExternalAddresses = InetAddressFilter.externalAddresses()
 val settings = HttpClientSettings.defaults().withInetAddressFilter(onlyExternalAddresses)
 val requestFactory: ClientHttpRequestFactory = ClientHttpRequestFactoryBuilder.jdk().build(settings)
 restClient = RestClient.builder().requestFactory(requestFactory).baseUrl("https://example.org").build()
	}
	fun someRestCall(name: String?): Details {
 return restClient.get().uri("/{name}/details", name)
 .retrieve().body(Details::class.java)!!
	}
}
```
Or you can also define one as a `@Bean` if you want to apply it to all auto-configured HTTP client builders:

-
Java
-
Kotlin

```
import org.springframework.boot.http.client.InetAddressFilter;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyHttpClientConfiguration {
	@Bean
	public InetAddressFilter httpClientInetAddressFilter() {
 return InetAddressFilter.of("192.168.1.0/24").andNot("192.168.1.1", "192.168.1.10");
	}
}
```
```
import org.springframework.boot.http.client.InetAddressFilter
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
@Configuration(proxyBeanMethods = false)
class MyHttpClientConfiguration {
	@Bean
	fun httpClientInetAddressFilter(): InetAddressFilter {
 return InetAddressFilter.of("192.168.1.0/24").andNot("192.168.1.1", "192.168.1.10")
	}
}
```

# Web Services

Spring Boot provides Web Services auto-configuration so that all you must do is define your `@Endpoint` beans.

The Spring Web Services features can be easily accessed with the `spring-boot-starter-webservices` module.

`SimpleWsdl11Definition` and `SimpleXsdSchema` beans can be automatically created for your WSDLs and XSDs respectively.
To do so, configure their location, as shown in the following example:

-
Properties
-
YAML

`spring.webservices.wsdl-locations=classpath:/wsdl````
spring:
 webservices:
 wsdl-locations: "classpath:/wsdl"
```
## Calling Web Services with WebServiceTemplate

If you need to call remote Web services from your application, you can use the `WebServiceTemplate` class.
Since `WebServiceTemplate` instances often need to be customized before being used, Spring Boot does not provide any single auto-configured `WebServiceTemplate` bean.
It does, however, auto-configure a `WebServiceTemplateBuilder`, which can be used to create `WebServiceTemplate` instances when needed.

The following code shows a typical example:

-
Java
-
Kotlin

```
import org.springframework.boot.webservices.client.WebServiceTemplateBuilder;
import org.springframework.stereotype.Service;
import org.springframework.ws.client.core.WebServiceTemplate;
import org.springframework.ws.soap.client.core.SoapActionCallback;
@Service
public class MyService {
	private final WebServiceTemplate webServiceTemplate;
	public MyService(WebServiceTemplateBuilder webServiceTemplateBuilder) {
 this.webServiceTemplate = webServiceTemplateBuilder.build();
	}
	public SomeResponse someWsCall(SomeRequest detailsReq) {
 return (SomeResponse) this.webServiceTemplate.marshalSendAndReceive(detailsReq,
 new SoapActionCallback("https://ws.example.com/action"));
	}
}
```
```
import org.springframework.boot.webservices.client.WebServiceTemplateBuilder
import org.springframework.stereotype.Service
import org.springframework.ws.client.core.WebServiceTemplate
import org.springframework.ws.soap.client.core.SoapActionCallback
@Service
class MyService(webServiceTemplateBuilder: WebServiceTemplateBuilder) {
	private val webServiceTemplate: WebServiceTemplate
	init {
 webServiceTemplate = webServiceTemplateBuilder.build()
	}
	fun someWsCall(detailsReq: SomeRequest): SomeResponse {
 return webServiceTemplate.marshalSendAndReceive(
 detailsReq,
 SoapActionCallback("https://ws.example.com/action")
 ) as SomeResponse
	}
}
```
By default, `WebServiceTemplateBuilder` detects a suitable HTTP-based `WebServiceMessageSender` using the available HTTP client libraries on the classpath.
You can also customize read and connection timeouts for an individual builder as follows:

-
Java
-
Kotlin

```
import java.time.Duration;
import org.springframework.boot.http.client.HttpClientSettings;
import org.springframework.boot.webservices.client.WebServiceMessageSenderFactory;
import org.springframework.boot.webservices.client.WebServiceTemplateBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.ws.client.core.WebServiceTemplate;
@Configuration(proxyBeanMethods = false)
public class MyWebServiceTemplateConfiguration {
	@Bean
	public WebServiceTemplate webServiceTemplate(WebServiceTemplateBuilder builder) {
 HttpClientSettings settings = HttpClientSettings.defaults()
 .withConnectTimeout(Duration.ofSeconds(2))
 .withReadTimeout(Duration.ofSeconds(2));
 builder.httpMessageSenderFactory(WebServiceMessageSenderFactory.http(settings));
 return builder.build();
	}
}
```
```
import org.springframework.boot.http.client.HttpClientSettings
import org.springframework.boot.webservices.client.WebServiceMessageSenderFactory
import org.springframework.boot.webservices.client.WebServiceTemplateBuilder
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.ws.client.core.WebServiceTemplate
import java.time.Duration
@Configuration(proxyBeanMethods = false)
class MyWebServiceTemplateConfiguration {
	@Bean
	fun webServiceTemplate(builder: WebServiceTemplateBuilder): WebServiceTemplate {
 val settings = HttpClientSettings.defaults()
 .withConnectTimeout(Duration.ofSeconds(2))
 .withReadTimeout(Duration.ofSeconds(2))
 builder.httpMessageSenderFactory(WebServiceMessageSenderFactory.http(settings))
 return builder.build()
	}
}
```
| You can also change the global HTTP client configuration used if not specific template customization code is applied. |

# Distributed Transactions With JTA

Spring Boot supports distributed JTA transactions across multiple XA resources by using a transaction manager retrieved from JNDI.

When a JTA environment is detected, Spring’s `JtaTransactionManager` is used to manage transactions.
Auto-configured JMS, DataSource, and JPA beans are upgraded to support XA transactions.
You can use standard Spring idioms, such as `@Transactional`, to participate in a distributed transaction.
If you are within a JTA environment and still want to use local transactions, you can set the `spring.jta.enabled` property to `false` to disable the JTA auto-configuration.

## Using a Jakarta EE Managed Transaction Manager

If you package your Spring Boot application as a `war` or `ear` file and deploy it to a Jakarta EE application server, you can use your application server’s built-in transaction manager.
Spring Boot tries to auto-configure a transaction manager by looking at common JNDI locations (`java:comp/UserTransaction`, `java:comp/TransactionManager`, and so on).
When using a transaction service provided by your application server, you generally also want to ensure that all resources are managed by the server and exposed over JNDI.
Spring Boot tries to auto-configure JMS by looking for a `ConnectionFactory` at the JNDI path (`java:/JmsXA` or `java:/XAConnectionFactory`), and you can use the `spring.datasource.jndi-name` property to configure your `DataSource`.

## Mixing XA and Non-XA JMS Connections

When using JTA, the primary JMS `ConnectionFactory` bean is XA-aware and participates in distributed transactions.
You can inject into your bean without needing to use any `@Qualifier`:

-
Java
-
Kotlin

```
import jakarta.jms.ConnectionFactory;
public class MyBean {
	public MyBean(ConnectionFactory connectionFactory) {
 // ...
	}
}
```
```
import jakarta.jms.ConnectionFactory
class MyBean(connectionFactory: ConnectionFactory?)
```
In some situations, you might want to process certain JMS messages by using a non-XA `ConnectionFactory`.
For example, your JMS processing logic might take longer than the XA timeout.

If you want to use a non-XA `ConnectionFactory`, you can the `nonXaJmsConnectionFactory` bean:

-
Java
-
Kotlin

```
import jakarta.jms.ConnectionFactory;
import org.springframework.beans.factory.annotation.Qualifier;
public class MyBean {
	public MyBean(@Qualifier("nonXaJmsConnectionFactory") ConnectionFactory connectionFactory) {
 // ...
	}
}
```
```
import jakarta.jms.ConnectionFactory
import org.springframework.beans.factory.annotation.Qualifier
class MyBean(@Qualifier("nonXaJmsConnectionFactory") connectionFactory: ConnectionFactory?)
```
For consistency, the `jmsConnectionFactory` bean is also provided by using the bean alias `xaJmsConnectionFactory`:

-
Java
-
Kotlin

```
import jakarta.jms.ConnectionFactory;
import org.springframework.beans.factory.annotation.Qualifier;
public class MyBean {
	public MyBean(@Qualifier("xaJmsConnectionFactory") ConnectionFactory connectionFactory) {
 // ...
	}
}
```
```
import jakarta.jms.ConnectionFactory
import org.springframework.beans.factory.annotation.Qualifier
class MyBean(@Qualifier("xaJmsConnectionFactory") connectionFactory: ConnectionFactory?)
```
## Supporting an Embedded Transaction Manager

The `XAConnectionFactoryWrapper` and `XADataSourceWrapper` interfaces can be used to support embedded transaction managers.
The interfaces are responsible for wrapping `XAConnectionFactory` and `XADataSource` beans and exposing them as regular `ConnectionFactory` and `DataSource` beans, which transparently enroll in the distributed transaction.
DataSource and JMS auto-configuration use JTA variants, provided you have a `JtaTransactionManager` bean and appropriate XA wrapper beans registered within your `ApplicationContext`.

# Messaging

The Spring Framework provides extensive support for integrating with messaging systems, from simplified use of the JMS API using `JmsClient` to a complete infrastructure to receive messages asynchronously.
Spring AMQP provides a similar feature set for the Advanced Message Queuing Protocol.
Spring Boot also provides auto-configuration options for `RabbitTemplate` and RabbitMQ.
Spring WebSocket natively includes support for STOMP messaging, and Spring Boot has support for that through starters and a small amount of auto-configuration.
Spring Boot also has support for Apache Kafka and Apache Pulsar.

# JMS

The `ConnectionFactory` interface provides a standard method of creating a `Connection` for interacting with a JMS broker.
Although Spring needs a `ConnectionFactory` to work with JMS, you generally need not use it directly yourself and can instead rely on higher level messaging abstractions.
(See the relevant section of the Spring Framework reference documentation for details.)
Spring Boot also auto-configures the necessary infrastructure to send and receive messages.

## ActiveMQ "Classic" Support

When ActiveMQ "Classic" is available on the classpath, Spring Boot can configure a `ConnectionFactory`.
If the broker is present, an embedded broker is automatically started and configured (provided no broker URL is specified through configuration and the embedded broker is not disabled in the configuration).

| If you use `spring-boot-starter-activemq`, the necessary dependencies to connect to an ActiveMQ "Classic" instance are provided, as is the Spring infrastructure to integrate with JMS.
Adding`org.apache.activemq:activemq-broker`to your application lets you use the embedded broker. |

ActiveMQ "Classic" configuration is controlled by external configuration properties in `spring.activemq.*`.

If `activemq-broker` is on the classpath, ActiveMQ "Classic" is auto-configured to use the VM transport, which starts a broker embedded in the same JVM instance.

You can disable the embedded broker by configuring the `spring.activemq.embedded.enabled` property, as shown in the following example:

-
Properties
-
YAML

`spring.activemq.embedded.enabled=false````
spring:
 activemq:
 embedded:
 enabled: false
```
The embedded broker will also be disabled if you configure the broker URL, as shown in the following example:

-
Properties
-
YAML

```
spring.activemq.broker-url=tcp://192.168.1.210:9876
spring.activemq.user=admin
spring.activemq.password=secret
```
```
spring:
 activemq:
 broker-url: "tcp://192.168.1.210:9876"
 user: "admin"
 password: "secret"
```
If you want to take full control over the embedded broker, see the ActiveMQ "Classic" documentation for further information.

By default, a `CachingConnectionFactory` wraps the native `ConnectionFactory` with sensible settings that you can control by external configuration properties in `spring.jms.*`:

-
Properties
-
YAML

`spring.jms.cache.session-cache-size=5````
spring:
 jms:
 cache:
 session-cache-size: 5
```
If you’d rather use native pooling, you can do so by adding a dependency to `org.messaginghub:pooled-jms` and configuring the `JmsPoolConnectionFactory` accordingly, as shown in the following example:

-
Properties
-
YAML

```
spring.activemq.pool.enabled=true
spring.activemq.pool.max-connections=50
```
```
spring:
 activemq:
 pool:
 enabled: true
 max-connections: 50
```
| See `ActiveMQProperties`for more of the supported options.
You can also register an arbitrary number of beans that implement`ActiveMQConnectionFactoryCustomizer`for more advanced customizations. |

By default, ActiveMQ "Classic" creates a destination if it does not yet exist so that destinations are resolved against their provided names.

## ActiveMQ Artemis Support

Spring Boot can auto-configure a `ConnectionFactory` when it detects that ActiveMQ Artemis is available on the classpath.
If the broker is present, an embedded broker is automatically started and configured (unless the mode property has been explicitly set).
The supported modes are `embedded` (to make explicit that an embedded broker is required and that an error should occur if the broker is not available on the classpath) and `native` (to connect to a broker using the `netty` transport protocol).
When the latter is configured, Spring Boot configures a `ConnectionFactory` that connects to a broker running on the local machine with the default settings.

| If you use `spring-boot-starter-artemis`, the necessary dependencies to connect to an existing ActiveMQ Artemis instance are provided, as well as the Spring infrastructure to integrate with JMS.
Adding`org.apache.activemq:artemis-jakarta-server`to your application lets you use embedded mode. |

ActiveMQ Artemis configuration is controlled by external configuration properties in `spring.artemis.*`.
For example, you might declare the following section in `application.properties`:

-
Properties
-
YAML

```
spring.artemis.mode=native
spring.artemis.broker-url=tcp://192.168.1.210:9876
spring.artemis.user=admin
spring.artemis.password=secret
```
```
spring:
 artemis:
 mode: native
 broker-url: "tcp://192.168.1.210:9876"
 user: "admin"
 password: "secret"
```
When embedding the broker, you can choose if you want to enable persistence and list the destinations that should be made available.
These can be specified as a comma-separated list to create them with the default options, or you can define bean(s) of type `JMSQueueConfiguration` or `TopicConfiguration`, for advanced queue and topic configurations, respectively.

By default, a `CachingConnectionFactory` wraps the native `ConnectionFactory` with sensible settings that you can control by external configuration properties in `spring.jms.*`:

-
Properties
-
YAML

`spring.jms.cache.session-cache-size=5````
spring:
 jms:
 cache:
 session-cache-size: 5
```
If you’d rather use native pooling, you can do so by adding a dependency on `org.messaginghub:pooled-jms` and configuring the `JmsPoolConnectionFactory` accordingly, as shown in the following example:

-
Properties
-
YAML

```
spring.artemis.pool.enabled=true
spring.artemis.pool.max-connections=50
```
```
spring:
 artemis:
 pool:
 enabled: true
 max-connections: 50
```
See `ArtemisProperties` for more supported options.

No JNDI lookup is involved, and destinations are resolved against their names, using either the `name` attribute in the ActiveMQ Artemis configuration or the names provided through configuration.

## Using a JNDI ConnectionFactory

If you are running your application in an application server, Spring Boot tries to locate a JMS `ConnectionFactory` by using JNDI.
By default, the `java:/JmsXA` and `java:/XAConnectionFactory` location are checked.
You can use the `spring.jms.jndi-name` property if you need to specify an alternative location, as shown in the following example:

-
Properties
-
YAML

`spring.jms.jndi-name=java:/MyConnectionFactory````
spring:
 jms:
 jndi-name: "java:/MyConnectionFactory"
```
## Sending a Message

Spring’s `JmsClient` is auto-configured, and you can autowire it directly into your own beans, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.jms.core.JmsClient;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final JmsClient jmsClient;
	public MyBean(JmsClient jmsClient) {
 this.jmsClient = jmsClient;
	}
	// ...
	public void someMethod() {
 this.jmsClient.destination("myQueue").send("hello");
	}
}
```
```
import org.springframework.jms.core.JmsClient
import org.springframework.stereotype.Component
@Component
class MyBean(private val jmsClient: JmsClient) {
	// ...
	fun someMethod() {
 jmsClient.destination("myQueue").send("hello")
	}
}
```
| `JmsMessagingTemplate`can be injected in a similar manner, and both use the traditional`JmsTemplate`that can be injected as well.
If a`DestinationResolver`or a`MessageConverter`bean is defined, it is associated automatically to the auto-configured`JmsTemplate`. |

## Receiving a Message

When the JMS infrastructure is present, any bean can be annotated with `@JmsListener` to create a listener endpoint.
If no `JmsListenerContainerFactory` has been defined, a default one is configured automatically.
If a `DestinationResolver`, a `MessageConverter`, or a `ExceptionListener` beans are defined, they are associated automatically with the default factory.

In most scenarios, message listener containers should be configured against the native `ConnectionFactory`.
This way each listener container has its own connection and this gives full responsibility to it in terms of local recovery.
The auto-configuration uses `ConnectionFactoryUnwrapper` to unwrap the native connection factory from the auto-configured one.

| The auto-configuration only unwraps `CachedConnectionFactory`. |

By default, the default factory is transactional.
If you run in an infrastructure where a `JtaTransactionManager` is present, it is associated to the listener container by default.
If not, the `sessionTransacted` flag is enabled.
In that latter scenario, you can associate your local data store transaction to the processing of an incoming message by adding `@Transactional` on your listener method (or a delegate thereof).
This ensures that the incoming message is acknowledged, once the local transaction has completed.
This also includes sending response messages that have been performed on the same JMS session.

The following component creates a listener endpoint on the `someQueue` destination:

-
Java
-
Kotlin

```
import org.springframework.jms.annotation.JmsListener;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	@JmsListener(destination = "someQueue")
	public void processMessage(String content) {
 // ...
	}
}
```
```
import org.springframework.jms.annotation.JmsListener
import org.springframework.stereotype.Component
@Component
class MyBean {
	@JmsListener(destination = "someQueue")
	fun processMessage(content: String?) {
 // ...
	}
}
```
| See the `@EnableJms`API documentation for more details. |

If you need to create more `JmsListenerContainerFactory` instances or if you want to override the default, Spring Boot provides a `DefaultJmsListenerContainerFactoryConfigurer` that you can use to initialize a `DefaultJmsListenerContainerFactory` with the same settings as the one that is auto-configured.

For instance, the following example exposes another factory that uses a specific `MessageConverter`:

-
Java
-
Kotlin

```
import jakarta.jms.ConnectionFactory;
import org.springframework.boot.jms.ConnectionFactoryUnwrapper;
import org.springframework.boot.jms.autoconfigure.DefaultJmsListenerContainerFactoryConfigurer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.jms.config.DefaultJmsListenerContainerFactory;
@Configuration(proxyBeanMethods = false)
public class MyJmsConfiguration {
	@Bean
	public DefaultJmsListenerContainerFactory myFactory(DefaultJmsListenerContainerFactoryConfigurer configurer,
 ConnectionFactory connectionFactory) {
 DefaultJmsListenerContainerFactory factory = new DefaultJmsListenerContainerFactory();
 configurer.configure(factory, ConnectionFactoryUnwrapper.unwrapCaching(connectionFactory));
 factory.setMessageConverter(new MyMessageConverter());
 return factory;
	}
}
```
```
import jakarta.jms.ConnectionFactory
import org.springframework.boot.jms.ConnectionFactoryUnwrapper
import org.springframework.boot.jms.autoconfigure.DefaultJmsListenerContainerFactoryConfigurer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.jms.config.DefaultJmsListenerContainerFactory
@Configuration(proxyBeanMethods = false)
class MyJmsConfiguration {
	@Bean
	fun myFactory(configurer: DefaultJmsListenerContainerFactoryConfigurer,
 connectionFactory: ConnectionFactory): DefaultJmsListenerContainerFactory {
 val factory = DefaultJmsListenerContainerFactory()
 configurer.configure(factory, ConnectionFactoryUnwrapper.unwrapCaching(connectionFactory))
 factory.setMessageConverter(MyMessageConverter())
 return factory
	}
}
```
| In the example above, the customization uses `ConnectionFactoryUnwrapper`to associate the native connection factory to the message listener container the same way the auto-configured factory does. |

Then you can use the factory in any `@JmsListener`-annotated method as follows:

-
Java
-
Kotlin

```
import org.springframework.jms.annotation.JmsListener;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	@JmsListener(destination = "someQueue", containerFactory = "myFactory")
	public void processMessage(String content) {
 // ...
	}
}
```
```
import org.springframework.jms.annotation.JmsListener
import org.springframework.stereotype.Component
@Component
class MyBean {
	@JmsListener(destination = "someQueue", containerFactory = "myFactory")
	fun processMessage(content: String?) {
 // ...
	}
}
```
Analogous to `DefaultJmsListenerContainerFactoryConfigurer`, Spring Boot also provides a `SimpleJmsListenerContainerFactoryConfigurer` that you can use to initialize a `SimpleJmsListenerContainerFactory` and apply the related settings that auto-configuration provides.

| In contrast to `DefaultMessageListenerContainer`that uses a pull-based mechanism (polling) to process messages,`SimpleMessageListenerContainer`uses a push-based mechanism that’s very close to the spirit of the standalone JMS specification.
To learn more about the differences between the two listener containers, consult their respective javadocs and Spring Framework reference documentation. |

# AMQP

The Advanced Message Queuing Protocol (AMQP) is a platform-neutral, wire-level protocol for message-oriented middleware.
The Spring AMQP project applies core Spring concepts to the development of AMQP-based messaging solutions.
Spring Boot offers several conveniences for working with AMQP through RabbitMQ, including the `spring-boot-starter-amqp` starter.

## RabbitMQ Support

RabbitMQ is a lightweight, reliable, scalable, and portable message broker based on the AMQP protocol. Spring uses RabbitMQ to communicate through the AMQP protocol.

RabbitMQ configuration is controlled by external configuration properties in `spring.rabbitmq.*`.
For example, you might declare the following section in `application.properties`:

-
Properties
-
YAML

```
spring.rabbitmq.host=localhost
spring.rabbitmq.port=5672
spring.rabbitmq.username=admin
spring.rabbitmq.password=secret
```
```
spring:
 rabbitmq:
 host: "localhost"
 port: 5672
 username: "admin"
 password: "secret"
```
Alternatively, you could configure the same connection using the `addresses` attribute:

-
Properties
-
YAML

`spring.rabbitmq.addresses=amqp://admin:secret@localhost````
spring:
 rabbitmq:
 addresses: "amqp://admin:secret@localhost"
```
| When specifying addresses that way, the `host`and`port`properties are ignored.
If the address uses the`amqps`protocol, SSL support is enabled automatically. |

See `RabbitProperties` for more of the supported property-based configuration options.
To configure lower-level details of the RabbitMQ `ConnectionFactory` that is used by Spring AMQP, define a `ConnectionFactoryCustomizer` bean.

If a `ConnectionNameStrategy` bean exists in the context, it will be automatically used to name connections created by the auto-configured `CachingConnectionFactory`.

To make an application-wide, additive customization to the `RabbitTemplate`, use a `RabbitTemplateCustomizer` bean.

| See Understanding AMQP, the protocol used by RabbitMQ for more details. |

## Sending a Message

Spring’s `AmqpTemplate` and `AmqpAdmin` are auto-configured, and you can autowire them directly into your own beans, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.amqp.core.AmqpAdmin;
import org.springframework.amqp.core.AmqpTemplate;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final AmqpAdmin amqpAdmin;
	private final AmqpTemplate amqpTemplate;
	public MyBean(AmqpAdmin amqpAdmin, AmqpTemplate amqpTemplate) {
 this.amqpAdmin = amqpAdmin;
 this.amqpTemplate = amqpTemplate;
	}
	// ...
	public void someMethod() {
 this.amqpAdmin.getQueueInfo("someQueue");
	}
	public void someOtherMethod() {
 this.amqpTemplate.convertAndSend("hello");
	}
}
```
```
import org.springframework.amqp.core.AmqpAdmin
import org.springframework.amqp.core.AmqpTemplate
import org.springframework.stereotype.Component
@Component
class MyBean(private val amqpAdmin: AmqpAdmin, private val amqpTemplate: AmqpTemplate) {
	// ...
	fun someMethod() {
 amqpAdmin.getQueueInfo("someQueue")
	}
	fun someOtherMethod() {
 amqpTemplate.convertAndSend("hello")
	}
}
```
| `RabbitMessagingTemplate`can be injected in a similar manner.
If a`MessageConverter`bean is defined, it is associated automatically to the auto-configured`AmqpTemplate`. |

If necessary, any `Queue` that is defined as a bean is automatically used to declare a corresponding queue on the RabbitMQ instance.

To retry operations, you can enable retries on the `AmqpTemplate` (for example, in the event that the broker connection is lost):

-
Properties
-
YAML

```
spring.rabbitmq.template.retry.enabled=true
spring.rabbitmq.template.retry.initial-interval=2s
```
```
spring:
 rabbitmq:
 template:
 retry:
 enabled: true
 initial-interval: "2s"
```
Retries are disabled by default.
You can also customize the `RetryTemplate` programmatically by declaring a `RabbitTemplateRetrySettingsCustomizer` bean.

If you need to create more `RabbitTemplate` instances or if you want to override the default, Spring Boot provides a `RabbitTemplateConfigurer` bean that you can use to initialize a `RabbitTemplate` with the same settings as the factories used by the auto-configuration.

If there’s a bean of type `RabbitTemplateObservationConvention` in the context, it will automatically be configured on the `RabbitTemplate`.

## Sending a Message To A Stream

To send a message to a particular stream, specify the name of the stream, as shown in the following example:

-
Properties
-
YAML

`spring.rabbitmq.stream.name=my-stream````
spring:
 rabbitmq:
 stream:
 name: "my-stream"
```
If a `MessageConverter`, `StreamMessageConverter`, `ProducerCustomizer` or `RabbitStreamTemplateObservationConvention` bean is defined, it is associated automatically to the auto-configured `RabbitStreamTemplate`.

If you need to create more `RabbitStreamTemplate` instances or if you want to override the default, Spring Boot provides a `RabbitStreamTemplateConfigurer` bean that you can use to initialize a `RabbitStreamTemplate` with the same settings as the factories used by the auto-configuration.

### SSL

To use SSL with RabbitMQ Streams, set `spring.rabbitmq.stream.ssl.enabled` to `true` or set `spring.rabbitmq.stream.ssl.bundle` to configure the SSL bundle to use.

## Receiving a Message

When the Rabbit infrastructure is present, any bean can be annotated with `@RabbitListener` to create a listener endpoint.
If no `RabbitListenerContainerFactory` has been defined, a default `SimpleRabbitListenerContainerFactory` is automatically configured and you can switch to a direct container using the `spring.rabbitmq.listener.type` property.
If a `MessageConverter`, a `MessageRecoverer` or a `RabbitListenerObservationConvention` bean is defined, it is automatically associated with the default factory.

The following sample component creates a listener endpoint on the `someQueue` queue:

-
Java
-
Kotlin

```
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	@RabbitListener(queues = "someQueue")
	public void processMessage(String content) {
 // ...
	}
}
```
```
import org.springframework.amqp.rabbit.annotation.RabbitListener
import org.springframework.stereotype.Component
@Component
class MyBean {
	@RabbitListener(queues = ["someQueue"])
	fun processMessage(content: String?) {
 // ...
	}
}
```
| See `@EnableRabbit`for more details. |

If you need to create more `RabbitListenerContainerFactory` instances or if you want to override the default, Spring Boot provides a `SimpleRabbitListenerContainerFactoryConfigurer` and a `DirectRabbitListenerContainerFactoryConfigurer` that you can use to initialize a `SimpleRabbitListenerContainerFactory` and a `DirectRabbitListenerContainerFactory` with the same settings as the factories used by the auto-configuration.

| It does not matter which container type you chose. Those two beans are exposed by the auto-configuration. |

For instance, the following configuration class exposes another factory that uses a specific `MessageConverter`:

-
Java
-
Kotlin

```
import org.springframework.amqp.rabbit.config.SimpleRabbitListenerContainerFactory;
import org.springframework.amqp.rabbit.connection.ConnectionFactory;
import org.springframework.boot.amqp.autoconfigure.SimpleRabbitListenerContainerFactoryConfigurer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyRabbitConfiguration {
	@Bean
	public SimpleRabbitListenerContainerFactory myFactory(SimpleRabbitListenerContainerFactoryConfigurer configurer) {
 SimpleRabbitListenerContainerFactory factory = new SimpleRabbitListenerContainerFactory();
 ConnectionFactory connectionFactory = getCustomConnectionFactory();
 configurer.configure(factory, connectionFactory);
 factory.setMessageConverter(new MyMessageConverter());
 return factory;
	}
	private ConnectionFactory getCustomConnectionFactory() {
 return ...
	}
}
```
```
import org.springframework.amqp.rabbit.config.SimpleRabbitListenerContainerFactory
import org.springframework.amqp.rabbit.connection.CachingConnectionFactory
import org.springframework.amqp.rabbit.connection.ConnectionFactory
import org.springframework.boot.amqp.autoconfigure.SimpleRabbitListenerContainerFactoryConfigurer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
@Configuration(proxyBeanMethods = false)
class MyRabbitConfiguration {
	@Bean
	fun myFactory(configurer: SimpleRabbitListenerContainerFactoryConfigurer): SimpleRabbitListenerContainerFactory {
 val factory = SimpleRabbitListenerContainerFactory()
 val connectionFactory = getCustomConnectionFactory()
 configurer.configure(factory, connectionFactory)
 factory.setMessageConverter(MyMessageConverter())
 return factory
	}
	fun getCustomConnectionFactory() : ConnectionFactory {
 return ...
	}
}
```
Then you can use the factory in any `@RabbitListener`-annotated method, as follows:

-
Java
-
Kotlin

```
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	@RabbitListener(queues = "someQueue", containerFactory = "myFactory")
	public void processMessage(String content) {
 // ...
	}
}
```
```
import org.springframework.amqp.rabbit.annotation.RabbitListener
import org.springframework.stereotype.Component
@Component
class MyBean {
	@RabbitListener(queues = ["someQueue"], containerFactory = "myFactory")
	fun processMessage(content: String?) {
 // ...
	}
}
```
You can enable retries to handle situations where your listener throws an exception.
By default, `RejectAndDontRequeueRecoverer` is used, but you can define a `MessageRecoverer` of your own.
When retries are exhausted, the message is rejected and either dropped or routed to a dead-letter exchange if the broker is configured to do so.
By default, retries are disabled.
You can also customize the `RetryPolicy` programmatically by declaring a `RabbitListenerRetrySettingsCustomizer` bean.

| By default, if retries are disabled and the listener throws an exception, the delivery is retried indefinitely.
You can modify this behavior in two ways: Set the `defaultRequeueRejected`property to`false`so that zero re-deliveries are attempted or throw an`AmqpRejectAndDontRequeueException`to signal the message should be rejected.
The latter is the mechanism used when retries are enabled and the maximum number of delivery attempts is reached. |

# Apache Kafka Support

Apache Kafka is supported by providing auto-configuration of the `spring-kafka` project.

Kafka configuration is controlled by external configuration properties in `spring.kafka.*`.
For example, you might declare the following section in `application.properties`:

-
Properties
-
YAML

```
spring.kafka.bootstrap-servers=localhost:9092
spring.kafka.consumer.group-id=myGroup
```
```
spring:
 kafka:
 bootstrap-servers: "localhost:9092"
 consumer:
 group-id: "myGroup"
```
| To create a topic on startup, add a bean of type `NewTopic`.
If the topic already exists, the bean is ignored. |

See `KafkaProperties` for more supported options.

## Sending a Message

Spring’s `KafkaTemplate` is auto-configured, and you can autowire it directly in your own beans, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final KafkaTemplate<String, String> kafkaTemplate;
	public MyBean(KafkaTemplate<String, String> kafkaTemplate) {
 this.kafkaTemplate = kafkaTemplate;
	}
	// ...
	public void someMethod() {
 this.kafkaTemplate.send("someTopic", "Hello");
	}
}
```
```
import org.springframework.kafka.core.KafkaTemplate
import org.springframework.stereotype.Component
@Component
class MyBean(private val kafkaTemplate: KafkaTemplate<String, String>) {
	// ...
	fun someMethod() {
 kafkaTemplate.send("someTopic", "Hello")
	}
}
```
| If the property `spring.kafka.producer.transaction-id-prefix`is defined, a`KafkaTransactionManager`is automatically configured.
Also, if a`RecordMessageConverter`bean is defined, it is automatically associated to the auto-configured`KafkaTemplate`. |

If there’s a bean of type `KafkaTemplateObservationConvention` in the context, it is automatically registered on the `KafkaTemplate`.

## Receiving a Message

When the Apache Kafka infrastructure is present, any bean can be annotated with `@KafkaListener` to create a listener endpoint.
If no `KafkaListenerContainerFactory` has been defined, a default one is automatically configured with keys defined in `spring.kafka.listener.*`.

The following component creates a listener endpoint on the `someTopic` topic:

-
Java
-
Kotlin

```
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	@KafkaListener(topics = "someTopic")
	public void processMessage(String content) {
 // ...
	}
}
```
```
import org.springframework.kafka.annotation.KafkaListener
import org.springframework.stereotype.Component
@Component
class MyBean {
	@KafkaListener(topics = ["someTopic"])
	fun processMessage(content: String?) {
 // ...
	}
}
```
If a `KafkaTransactionManager` bean is defined, it is automatically associated to the container factory.
Similarly, if a `RecordFilterStrategy`, `CommonErrorHandler`, `AfterRollbackProcessor` or `ConsumerAwareRebalanceListener` bean is defined, it is automatically associated to the default factory.

Depending on the listener type, a `RecordMessageConverter` or `BatchMessageConverter` bean is associated to the default factory.
If only a `RecordMessageConverter` bean is present for a batch listener, it is wrapped in a `BatchMessageConverter`.

| A custom `ChainedKafkaTransactionManager`must be marked`@Primary`as it usually references the auto-configured`KafkaTransactionManager`bean. |

If there’s a bean of type `KafkaListenerObservationConvention` in the context, it is automatically registered on the container factory.

## Kafka Streams

Spring for Apache Kafka provides a factory bean to create a `StreamsBuilder` object and manage the lifecycle of its streams.
Spring Boot auto-configures the required `KafkaStreamsConfiguration` bean as long as `kafka-streams` is on the classpath and Kafka Streams is enabled by the `@EnableKafkaStreams` annotation.

Enabling Kafka Streams means that the application id and bootstrap servers must be set.
The former can be configured using `spring.kafka.streams.application-id`, defaulting to `spring.application.name` if not set.
The latter can be set globally or specifically overridden only for streams.

Several additional properties are available using dedicated properties; other arbitrary Kafka properties can be set using the `spring.kafka.streams.properties` namespace.
See also Additional Kafka Properties for more information.

To use the factory bean, wire `StreamsBuilder` into your `@Bean` as shown in the following example:

-
Java
-
Kotlin

```
import java.util.Locale;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.KeyValue;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.Produced;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.annotation.EnableKafkaStreams;
@Configuration(proxyBeanMethods = false)
@EnableKafkaStreams
public class MyKafkaStreamsConfiguration {
	@Bean
	public KStream<Integer, String> kStream(StreamsBuilder streamsBuilder) {
 KStream<Integer, String> stream = streamsBuilder.stream("ks1In");
 stream.map(this::uppercaseValue)
 .to("ks1Out",
 Produced.with(Serdes.Integer(), new org.springframework.kafka.support.serializer.JsonSerde<>()));
 return stream;
	}
	private KeyValue<Integer, String> uppercaseValue(Integer key, String value) {
 return new KeyValue<>(key, value.toUpperCase(Locale.getDefault()));
	}
}
```
```
import org.apache.kafka.common.serialization.Serdes
import org.apache.kafka.streams.KeyValue
import org.apache.kafka.streams.StreamsBuilder
import org.apache.kafka.streams.kstream.KStream
import org.apache.kafka.streams.kstream.Produced
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.kafka.annotation.EnableKafkaStreams
@Configuration(proxyBeanMethods = false)
@EnableKafkaStreams
class MyKafkaStreamsConfiguration {
	@Bean
	fun kStream(streamsBuilder: StreamsBuilder): KStream<Int, String> {
 val stream = streamsBuilder.stream<Int, String>("ks1In")
 stream.map(this::uppercaseValue).to("ks1Out", Produced.with(Serdes.Integer(),
 org.springframework.kafka.support.serializer.JsonSerde()))
 return stream
	}
	private fun uppercaseValue(key: Int, value: String): KeyValue<Int, String> {
 return KeyValue(key, value.uppercase())
	}
}
```
By default, the streams managed by the `StreamsBuilder` object are started automatically.
You can customize this behavior using the `spring.kafka.streams.auto-startup` property.

| You can also register an arbitrary number of beans that implement `StreamsBuilderFactoryBeanConfigurer`for more advanced customizations. |

## Additional Kafka Properties

The properties supported by auto configuration are shown in the Integration Properties section of the Appendix. Note that, for the most part, these properties (hyphenated or camelCase) map directly to the Apache Kafka dotted properties. See the Apache Kafka documentation for details.

Properties that don’t include a client type (`producer`, `consumer`, `admin`, or `streams`) in their name are considered to be common and apply to all clients.
Most of these common properties can be overridden for one or more of the client types, if needed.

Apache Kafka designates properties with an importance of HIGH, MEDIUM, or LOW. Spring Boot auto-configuration supports all HIGH importance properties, some selected MEDIUM and LOW properties, and any properties that do not have a default value.

Only a subset of the properties supported by Kafka are available directly through the `KafkaProperties` class.
If you wish to configure the individual client types with additional properties that are not directly supported, use the following properties:

-
Properties
-
YAML

```
spring.kafka.properties[prop.one]=first
spring.kafka.admin.properties[prop.two]=second
spring.kafka.consumer.properties[prop.three]=third
spring.kafka.producer.properties[prop.four]=fourth
spring.kafka.streams.properties[prop.five]=fifth
```
```
spring:
 kafka:
 properties:
 "[prop.one]": "first"
 admin:
 properties:
 "[prop.two]": "second"
 consumer:
 properties:
 "[prop.three]": "third"
 producer:
 properties:
 "[prop.four]": "fourth"
 streams:
 properties:
 "[prop.five]": "fifth"
```
This sets the common `prop.one` Kafka property to `first` (applies to producers, consumers, admins, and streams), the `prop.two` admin property to `second`, the `prop.three` consumer property to `third`, the `prop.four` producer property to `fourth` and the `prop.five` streams property to `fifth`.

You can also configure the Spring Kafka `JacksonJsonDeserializer` as follows:

-
Properties
-
YAML

```
spring.kafka.consumer.value-deserializer=org.springframework.kafka.support.serializer.JacksonJsonDeserializer
spring.kafka.consumer.properties[spring.json.value.default.type]=com.example.Invoice
spring.kafka.consumer.properties[spring.json.trusted.packages]=com.example.main,com.example.another
```
```
spring:
 kafka:
 consumer:
 value-deserializer: "org.springframework.kafka.support.serializer.JacksonJsonDeserializer"
 properties:
 "[spring.json.value.default.type]": "com.example.Invoice"
 "[spring.json.trusted.packages]": "com.example.main,com.example.another"
```
Similarly, you can disable the `JacksonJsonSerializer` default behavior of sending type information in headers:

-
Properties
-
YAML

```
spring.kafka.producer.value-serializer=org.springframework.kafka.support.serializer.JacksonJsonSerializer
spring.kafka.producer.properties[spring.json.add.type.headers]=false
```
```
spring:
 kafka:
 producer:
 value-serializer: "org.springframework.kafka.support.serializer.JacksonJsonSerializer"
 properties:
 "[spring.json.add.type.headers]": false
```
| Properties set in this way override any configuration item that Spring Boot explicitly supports. |

## Testing with Embedded Kafka

Spring for Apache Kafka provides a convenient way to test projects with an embedded Apache Kafka broker.
To use this feature, annotate a test class with `@EmbeddedKafka` from the `spring-kafka-test` module.
For more information, please see the Spring for Apache Kafka reference manual.

To make Spring Boot auto-configuration work with the aforementioned embedded Apache Kafka broker, you need to remap a system property for embedded broker addresses (populated by the `EmbeddedKafkaBroker`) into the Spring Boot configuration property for Apache Kafka.
There are several ways to do that:

-
Provide a system property to map embedded broker addresses into `spring.kafka.bootstrap-servers`in the test class:

-
Java
-
Kotlin

```
	static {
 System.setProperty(EmbeddedKafkaBroker.BROKER_LIST_PROPERTY, "spring.kafka.bootstrap-servers");
	}
```
```
	init {
 System.setProperty(EmbeddedKafkaBroker.BROKER_LIST_PROPERTY, "spring.kafka.bootstrap-servers")
	}
```
-
Configure a property name on the `@EmbeddedKafka`annotation:

-
Java
-
Kotlin

```
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.kafka.test.context.EmbeddedKafka;
@SpringBootTest
@EmbeddedKafka(topics = "someTopic", bootstrapServersProperty = "spring.kafka.bootstrap-servers")
class MyTest {
	// ...
}
```
```
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.kafka.test.context.EmbeddedKafka
@SpringBootTest
@EmbeddedKafka(topics = ["someTopic"], bootstrapServersProperty = "spring.kafka.bootstrap-servers")
class MyTest {
	// ...
}
```
-
Use a placeholder in configuration properties:

-
Properties
-
YAML

`spring.kafka.bootstrap-servers=${spring.embedded.kafka.brokers}````
spring:
 kafka:
 bootstrap-servers: "${spring.embedded.kafka.brokers}"
```

# Apache Pulsar Support

Apache Pulsar is supported by providing auto-configuration of the Spring for Apache Pulsar project.

Spring Boot will auto-configure and register the Spring for Apache Pulsar components when `org.springframework.pulsar:spring-pulsar` is on the classpath.

There is the `spring-boot-starter-pulsar` starter for conveniently collecting the dependencies for use.

## Connecting to Pulsar

When you use the Pulsar starter, Spring Boot will auto-configure and register a `PulsarClient` bean.

By default, the application tries to connect to a local Pulsar instance at `pulsar://localhost:6650`.
This can be adjusted by setting the `spring.pulsar.client.service-url` property to a different value.

| The value must be a valid Pulsar Protocol URL |

You can configure the client by specifying any of the `spring.pulsar.client.*` prefixed application properties.

If you need more control over the configuration, consider registering one or more `PulsarClientBuilderCustomizer` beans.

### Authentication

To connect to a Pulsar cluster that requires authentication, you need to specify which authentication plugin to use by setting the `pluginClassName` and any parameters required by the plugin.
You can set the parameters as a map of parameter names to parameter values.
The following example shows how to configure the `AuthenticationOAuth2` plugin.

-
Properties
-
YAML

```
spring.pulsar.client.authentication.plugin-class-name=org.apache.pulsar.client.impl.auth.oauth2.AuthenticationOAuth2
spring.pulsar.client.authentication.param.issuerUrl=https://auth.server.cloud/
spring.pulsar.client.authentication.param.privateKey=file:///Users/some-key.json
spring.pulsar.client.authentication.param.audience=urn:sn:acme:dev:my-instance
```
```
spring:
 pulsar:
 client:
 authentication:
 plugin-class-name: org.apache.pulsar.client.impl.auth.oauth2.AuthenticationOAuth2
 param:
 issuerUrl: https://auth.server.cloud/
 privateKey: file:///Users/some-key.json
 audience: urn:sn:acme:dev:my-instance
```
| You need to ensure that names defined under For example, if you want to configure the issuer url for the This lack of relaxed binding also makes using environment variables for authentication parameters problematic because the case sensitivity is lost during translation. If you use environment variables for the parameters then you will need to follow these steps in the Spring for Apache Pulsar reference documentation for it to work properly. |

### SSL

By default, Pulsar clients communicate with Pulsar services in plain text. You can follow these steps in the Spring for Apache Pulsar reference documentation to enable TLS encryption.

For complete details on the client and authentication see the Spring for Apache Pulsar reference documentation.

## Connecting to Pulsar Administration

Spring for Apache Pulsar’s `PulsarAdministration` client is also auto-configured.

By default, the application tries to connect to a local Pulsar instance at `http://localhost:8080`.
This can be adjusted by setting the `spring.pulsar.admin.service-url` property to a different value in the form `(http|https)://<host>:<port>`.

If you need more control over the configuration, consider registering one or more `PulsarAdminBuilderCustomizer` beans.

### Authentication

When accessing a Pulsar cluster that requires authentication, the admin client requires the same security configuration as the regular Pulsar client.
You can use the aforementioned authentication configuration by replacing `spring.pulsar.client.authentication` with `spring.pulsar.admin.authentication`.

| To create a topic on startup, add a bean of type `PulsarTopic`.
If the topic already exists, the bean is ignored. |

## Sending a Message

Spring’s `PulsarTemplate` is auto-configured, and you can use it to send messages, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.pulsar.core.PulsarTemplate;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final PulsarTemplate<String> pulsarTemplate;
	public MyBean(PulsarTemplate<String> pulsarTemplate) {
 this.pulsarTemplate = pulsarTemplate;
	}
	public void someMethod() {
 this.pulsarTemplate.send("someTopic", "Hello");
	}
}
```
```
import org.apache.pulsar.client.api.PulsarClientException
import org.springframework.pulsar.core.PulsarTemplate
import org.springframework.stereotype.Component
@Component
class MyBean(private val pulsarTemplate: PulsarTemplate<String>) {
	@Throws(PulsarClientException::class)
	fun someMethod() {
 pulsarTemplate.send("someTopic", "Hello")
	}
}
```
The `PulsarTemplate` relies on a `PulsarProducerFactory` to create the underlying Pulsar producer.
Spring Boot auto-configuration also provides this producer factory, which by default, caches the producers that it creates.
You can configure the producer factory and cache settings by specifying any of the `spring.pulsar.producer.*` and `spring.pulsar.producer.cache.*` prefixed application properties.

If you need more control over the producer factory configuration, consider registering one or more `ProducerBuilderCustomizer` beans.
These customizers are applied to all created producers.
You can also pass in a `ProducerBuilderCustomizer` when sending a message to only affect the current producer.

If you need more control over the message being sent, you can pass in a `TypedMessageBuilderCustomizer` when sending a message.

## Receiving a Message

When the Apache Pulsar infrastructure is present, any bean can be annotated with `@PulsarListener` to create a listener endpoint.
The following component creates a listener endpoint on the `someTopic` topic:

-
Java
-
Kotlin

```
import org.springframework.pulsar.annotation.PulsarListener;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	@PulsarListener(topics = "someTopic")
	public void processMessage(String content) {
 // ...
	}
}
```
```
import org.springframework.pulsar.annotation.PulsarListener
import org.springframework.stereotype.Component
@Component
class MyBean {
	@PulsarListener(topics = ["someTopic"])
	fun processMessage(content: String?) {
 // ...
	}
}
```
Spring Boot auto-configuration provides all the components necessary for `PulsarListener`, such as the `PulsarListenerContainerFactory` and the consumer factory it uses to construct the underlying Pulsar consumers.
You can configure these components by specifying any of the `spring.pulsar.listener.*` and `spring.pulsar.consumer.*` prefixed application properties.

If you need more control over the configuration of the consumer factory, consider registering one or more `ConsumerBuilderCustomizer` beans.
These customizers are applied to all consumers created by the factory, and therefore all `@PulsarListener` instances.
You can also customize a single listener by setting the `consumerCustomizer` attribute of the `@PulsarListener` annotation.

If you need more control over the actual container factory configuration, consider registering one or more `PulsarContainerFactoryCustomizer<ConcurrentPulsarListenerContainerFactory<?>>` beans.

## Reading a Message

The Pulsar reader interface enables applications to manually manage cursors. When you use a reader to connect to a topic you need to specify which message the reader begins reading from when it connects to a topic.

When the Apache Pulsar infrastructure is present, any bean can be annotated with `@PulsarReader` to consume messages using a reader.
The following component creates a reader endpoint that starts reading messages from the beginning of the `someTopic` topic:

-
Java
-
Kotlin

```
import org.springframework.pulsar.annotation.PulsarReader;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	@PulsarReader(topics = "someTopic", startMessageId = "earliest")
	public void processMessage(String content) {
 // ...
	}
}
```
```
import org.springframework.pulsar.annotation.PulsarReader
import org.springframework.stereotype.Component
@Component
class MyBean {
	@PulsarReader(topics = ["someTopic"], startMessageId = "earliest")
	fun processMessage(content: String?) {
 // ...
	}
}
```
The `@PulsarReader` relies on a `PulsarReaderFactory` to create the underlying Pulsar reader.
Spring Boot auto-configuration provides this reader factory which can be customized by setting any of the `spring.pulsar.reader.*` prefixed application properties.

If you need more control over the configuration of the reader factory, consider registering one or more `ReaderBuilderCustomizer` beans.
These customizers are applied to all readers created by the factory, and therefore all `@PulsarReader` instances.
You can also customize a single listener by setting the `readerCustomizer` attribute of the `@PulsarReader` annotation.

If you need more control over the actual container factory configuration, consider registering one or more `PulsarContainerFactoryCustomizer<DefaultPulsarReaderContainerFactory<?>>` beans.

| For more details on any of the above components and to discover other available features, see the Spring for Apache Pulsar reference documentation. |

## Transaction Support

Spring for Apache Pulsar supports transactions when using `PulsarTemplate` and `@PulsarListener`.

Setting the `spring.pulsar.transaction.enabled` property to `true` will:

-
Configure a `PulsarTransactionManager`bean
-
Enable transaction support for `PulsarTemplate`
-
Enable transaction support for `@PulsarListener`methods

The `transactional` attribute of `@PulsarListener` can be used to fine-tune when transactions should be used with listeners.

For more control of the Spring for Apache Pulsar transaction features you should define your own `PulsarTemplate` and/or `ConcurrentPulsarListenerContainerFactory` beans.
You can also define a `PulsarAwareTransactionManager` bean if the default auto-configured `PulsarTransactionManager` is not suitable.

## Additional Pulsar Properties

The properties supported by auto-configuration are shown in the Integration Properties section of the Appendix. Note that, for the most part, these properties (hyphenated or camelCase) map directly to the Apache Pulsar configuration properties. See the Apache Pulsar documentation for details.

Only a subset of the properties supported by Pulsar are available directly through the `PulsarProperties` class.
If you wish to tune the auto-configured components with additional properties that are not directly supported, you can use the customizer supported by each aforementioned component.

# RSocket

RSocket is a binary protocol for use on byte stream transports. It enables symmetric interaction models through async message passing over a single connection.

The `spring-messaging` module of the Spring Framework provides support for RSocket requesters and responders, both on the client and on the server side.
See the RSocket section of the Spring Framework reference for more details, including an overview of the RSocket protocol.

## RSocket Strategies Auto-configuration

Spring Boot auto-configures an `RSocketStrategies` bean that provides all the required infrastructure for encoding and decoding RSocket payloads.
By default, the auto-configuration will try to configure the following (in order):

-
CBOR codecs with Jackson
-
JSON codecs with Jackson

The `spring-boot-starter-rsocket` starter provides both dependencies.
See the Jackson support section to know more about customization possibilities.

Developers can customize the `RSocketStrategies` component by creating beans that implement the `RSocketStrategiesCustomizer` interface.
Note that their `@Order` is important, as it determines the order of codecs.

## RSocket Server Auto-configuration

Spring Boot provides RSocket server auto-configuration.
The required dependencies are provided by the `spring-boot-starter-rsocket`.

Spring Boot allows exposing RSocket over WebSocket from a WebFlux server, or standing up an independent RSocket server. This depends on the type of application and its configuration.

For WebFlux application (that is of type `WebApplicationType.REACTIVE`), the RSocket server will be plugged into the Web Server only if the following properties match:

-
Properties
-
YAML

```
spring.rsocket.server.mapping-path=/rsocket
spring.rsocket.server.transport=websocket
```
```
spring:
 rsocket:
 server:
 mapping-path: "/rsocket"
 transport: "websocket"
```
| Plugging RSocket into a web server is only supported with Reactor Netty, as RSocket itself is built with that library. |

Alternatively, an RSocket TCP or websocket server is started as an independent, embedded server. Besides the dependency requirements, the only required configuration is to define a port for that server:

-
Properties
-
YAML

`spring.rsocket.server.port=9898````
spring:
 rsocket:
 server:
 port: 9898
```
## Spring Messaging RSocket Support

Spring Boot will auto-configure the Spring Messaging infrastructure for RSocket.

This means that Spring Boot will create a `RSocketMessageHandler` bean that will handle RSocket requests to your application.

| You can use `@ControllerAdvice`to handle exceptions. |

## Calling RSocket Services with RSocketRequester

Once the `RSocket` channel is established between server and client, any party can send or receive requests to the other.

As a server, you can get injected with an `RSocketRequester` instance on any handler method of an RSocket `@Controller`.
As a client, you need to configure and establish an RSocket connection first.
Spring Boot auto-configures an `RSocketRequester.Builder` for such cases with the expected codecs and applies any `RSocketConnectorConfigurer` bean.

The `RSocketRequester.Builder` instance is a prototype bean, meaning each injection point will provide you with a new instance .
This is done on purpose since this builder is stateful and you should not create requesters with different setups using the same instance.

The following code shows a typical example:

-
Java
-
Kotlin

```
import reactor.core.publisher.Mono;
import org.springframework.messaging.rsocket.RSocketRequester;
import org.springframework.stereotype.Service;
@Service
public class MyService {
	private final RSocketRequester rsocketRequester;
	public MyService(RSocketRequester.Builder rsocketRequesterBuilder) {
 this.rsocketRequester = rsocketRequesterBuilder.tcp("example.org", 9898);
	}
	public Mono<User> someRSocketCall(String name) {
 return this.rsocketRequester.route("user").data(name).retrieveMono(User.class);
	}
}
```
```
import org.springframework.messaging.rsocket.RSocketRequester
import org.springframework.stereotype.Service
import reactor.core.publisher.Mono
@Service
class MyService(rsocketRequesterBuilder: RSocketRequester.Builder) {
	private val rsocketRequester: RSocketRequester
	init {
 rsocketRequester = rsocketRequesterBuilder.tcp("example.org", 9898)
	}
	fun someRSocketCall(name: String): Mono<User> {
 return rsocketRequester.route("user").data(name).retrieveMono(
 User::class.java
 )
	}
}
```

# Spring Integration

Spring Boot offers several conveniences for working with Spring Integration, including the `spring-boot-starter-integration` starter.
Spring Integration provides abstractions over messaging and also other transports such as HTTP, TCP, and others.
If Spring Integration is available on your classpath, it is initialized through the `@EnableIntegration` annotation.

Spring Integration polling logic relies on the auto-configured `TaskScheduler`.
The default `PollerMetadata` (poll unbounded number of messages every second) can be customized with `spring.integration.poller.*` configuration properties.

Spring Boot also configures some features that are triggered by the presence of additional Spring Integration modules.
If `spring-integration-jmx` is also on the classpath, message processing statistics are published over JMX.
If `spring-integration-jdbc` is available, the default database schema can be created on startup, as shown in the following line:

-
Properties
-
YAML

`spring.integration.jdbc.initialize-schema=always````
spring:
 integration:
 jdbc:
 initialize-schema: "always"
```
If `spring-integration-rsocket` is available, developers can configure an RSocket server using `spring.rsocket.server.*` properties and let it use `IntegrationRSocketEndpoint` or `RSocketOutboundGateway` components to handle incoming RSocket messages.
This infrastructure can handle Spring Integration RSocket channel adapters and `@MessageMapping` handlers (given `spring.integration.rsocket.server.message-mapping-enabled` is configured).

Spring Boot can also auto-configure an `ClientRSocketConnector` using configuration properties:

-
Properties
-
YAML

```
spring.integration.rsocket.client.host=example.org
spring.integration.rsocket.client.port=9898
```
```
# Connecting to a RSocket server over TCP
spring:
 integration:
 rsocket:
 client:
 host: "example.org"
 port: 9898
```
-
Properties
-
YAML

`spring.integration.rsocket.client.uri=ws://example.org````
# Connecting to a RSocket Server over WebSocket
spring:
 integration:
 rsocket:
 client:
 uri: "ws://example.org"
```
See the `IntegrationAutoConfiguration` and `IntegrationProperties` classes for more details.

# WebSockets

Spring Boot provides WebSockets auto-configuration for embedded Tomcat and Jetty. If you deploy a war file to a standalone container, Spring Boot assumes that the container is responsible for the configuration of its WebSocket support.

Spring Framework provides rich WebSocket support for MVC web applications that can be easily accessed through the `spring-boot-starter-websocket` module.

WebSocket support is also available for reactive web applications and requires to include the WebSocket API alongside `spring-boot-starter-webflux`:

```
<dependency>
	<groupId>jakarta.websocket</groupId>
	<artifactId>jakarta.websocket-api</artifactId>
</dependency>
```

# Security

Spring Boot provides helpful support for general purpose security features such as OAuth2 and SAML 2.0. This section covers those general purpose security concerns.

If you’re looking to secure web applications, see the web security section instead of this one.

# OAuth2

OAuth2 is a widely used authorization framework.

## Client

If you have `spring-security-oauth2-client` on your classpath, you can take advantage of some auto-configuration to set up OAuth2/Open ID Connect clients.
This configuration makes use of the properties under `OAuth2ClientProperties`.
The same properties are applicable to both servlet and reactive web applications.

Each registration must specify an OAuth 2 provider.
When set, the value of the `spring.security.oauth2.client.registration.<registration-id>.provider` property is used to specify the registration’s provider.
If the `provider` property is not set, the registration’s ID is used instead.
Both approaches are shown in the following example:

-
Properties
-
YAML

```
spring.security.oauth2.client.registration.my-client.client-id=abcd
spring.security.oauth2.client.registration.my-client.client-secret=password
spring.security.oauth2.client.registration.my-client.provider=example
spring.security.oauth2.client.registration.example.client-id=abcd
spring.security.oauth2.client.registration.example.client-secret=password
```
```
spring:
 security:
 oauth2:
 client:
 registration:
 my-client:
 client-id: "abcd"
 client-secret: "password"
 provider: "example"
 example:
 client-id: "abcd"
 client-secret: "password"
```
The registrations `my-client` and `example` will both use the provider with ID `example`.
The former will do so due to the value of the `spring.security.oauth2.client.registration.my-client.provider` property.
The latter will do so due to its ID being `example` and there being no `provider` property configured for the registration.

The specified provider can either be a reference to a provider configured using `spring.security.oauth2.client.provider.<provider-id>.*` properties or one of the known common providers.

You can register multiple OAuth2 clients and providers under the `spring.security.oauth2.client` prefix, as shown in the following example:

-
Properties
-
YAML

```
spring.security.oauth2.client.registration.my-login-client.client-id=abcd
spring.security.oauth2.client.registration.my-login-client.client-secret=password
spring.security.oauth2.client.registration.my-login-client.client-name=Client for OpenID Connect
spring.security.oauth2.client.registration.my-login-client.provider=my-oauth-provider
spring.security.oauth2.client.registration.my-login-client.scope=openid,profile,email,phone,address
spring.security.oauth2.client.registration.my-login-client.redirect-uri={baseUrl}/login/oauth2/code/{registrationId}
spring.security.oauth2.client.registration.my-login-client.client-authentication-method=client_secret_basic
spring.security.oauth2.client.registration.my-login-client.authorization-grant-type=authorization_code
spring.security.oauth2.client.registration.my-client-1.client-id=abcd
spring.security.oauth2.client.registration.my-client-1.client-secret=password
spring.security.oauth2.client.registration.my-client-1.client-name=Client for user scope
spring.security.oauth2.client.registration.my-client-1.provider=my-oauth-provider
spring.security.oauth2.client.registration.my-client-1.scope=user
spring.security.oauth2.client.registration.my-client-1.redirect-uri={baseUrl}/authorized/user
spring.security.oauth2.client.registration.my-client-1.client-authentication-method=client_secret_basic
spring.security.oauth2.client.registration.my-client-1.authorization-grant-type=authorization_code
spring.security.oauth2.client.registration.my-client-2.client-id=abcd
spring.security.oauth2.client.registration.my-client-2.client-secret=password
spring.security.oauth2.client.registration.my-client-2.client-name=Client for email scope
spring.security.oauth2.client.registration.my-client-2.provider=my-oauth-provider
spring.security.oauth2.client.registration.my-client-2.scope=email
spring.security.oauth2.client.registration.my-client-2.redirect-uri={baseUrl}/authorized/email
spring.security.oauth2.client.registration.my-client-2.client-authentication-method=client_secret_basic
spring.security.oauth2.client.registration.my-client-2.authorization-grant-type=authorization_code
spring.security.oauth2.client.provider.my-oauth-provider.authorization-uri=https://my-auth-server.com/oauth2/authorize
spring.security.oauth2.client.provider.my-oauth-provider.token-uri=https://my-auth-server.com/oauth2/token
spring.security.oauth2.client.provider.my-oauth-provider.user-info-uri=https://my-auth-server.com/userinfo
spring.security.oauth2.client.provider.my-oauth-provider.user-info-authentication-method=header
spring.security.oauth2.client.provider.my-oauth-provider.jwk-set-uri=https://my-auth-server.com/oauth2/jwks
spring.security.oauth2.client.provider.my-oauth-provider.user-name-attribute=name
```
```
spring:
 security:
 oauth2:
 client:
 registration:
 my-login-client:
 client-id: "abcd"
 client-secret: "password"
 client-name: "Client for OpenID Connect"
 provider: "my-oauth-provider"
 scope: "openid,profile,email,phone,address"
 redirect-uri: "{baseUrl}/login/oauth2/code/{registrationId}"
 client-authentication-method: "client_secret_basic"
 authorization-grant-type: "authorization_code"
 my-client-1:
 client-id: "abcd"
 client-secret: "password"
 client-name: "Client for user scope"
 provider: "my-oauth-provider"
 scope: "user"
 redirect-uri: "{baseUrl}/authorized/user"
 client-authentication-method: "client_secret_basic"
 authorization-grant-type: "authorization_code"
 my-client-2:
 client-id: "abcd"
 client-secret: "password"
 client-name: "Client for email scope"
 provider: "my-oauth-provider"
 scope: "email"
 redirect-uri: "{baseUrl}/authorized/email"
 client-authentication-method: "client_secret_basic"
 authorization-grant-type: "authorization_code"
 provider:
 my-oauth-provider:
 authorization-uri: "https://my-auth-server.com/oauth2/authorize"
 token-uri: "https://my-auth-server.com/oauth2/token"
 user-info-uri: "https://my-auth-server.com/userinfo"
 user-info-authentication-method: "header"
 jwk-set-uri: "https://my-auth-server.com/oauth2/jwks"
 user-name-attribute: "name"
```
In this example, there are three registrations.
In order of declaration, their IDs are `my-login-client`, `my-client-1`, and `my-client-2`.
There is also a single provider with ID `my-oauth-provider`.

For OpenID Connect providers that support OpenID Connect discovery, the configuration can be further simplified.
The provider needs to be configured with an `issuer-uri` which is the URI that it asserts as its Issuer Identifier.
For example, if the `issuer-uri` provided is `https://example.com`, then an “OpenID Provider Configuration Request” will be made to `https://example.com/.well-known/openid-configuration`.
The result is expected to be an “OpenID Provider Configuration Response”.
The following example shows how an OpenID Connect Provider can be configured with the `issuer-uri`:

-
Properties
-
YAML

`spring.security.oauth2.client.provider.oidc-provider.issuer-uri=https://dev-123456.oktapreview.com/oauth2/default/````
spring:
 security:
 oauth2:
 client:
 provider:
 oidc-provider:
 issuer-uri: "https://dev-123456.oktapreview.com/oauth2/default/"
```
By default, Spring Security’s `OAuth2LoginAuthenticationFilter` only processes URLs matching `/login/oauth2/code/*`.
If you want to customize the `redirect-uri` to use a different pattern, you need to provide configuration to process that custom pattern.
For example, for servlet applications, you can add your own `SecurityFilterChain` that resembles the following:

-
Java
-
Kotlin

```
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.web.SecurityFilterChain;
@Configuration(proxyBeanMethods = false)
@EnableWebSecurity
public class MyOAuthClientConfiguration {
	@Bean
	public SecurityFilterChain securityFilterChain(HttpSecurity http) {
 http.authorizeHttpRequests((requests) ->
 requests.anyRequest().authenticated()
 );
 http.oauth2Login((login) ->
 login.redirectionEndpoint((endpoint) ->
 endpoint.baseUri("/login/oauth2/callback/*")
 )
 );
 return http.build();
	}
}
```
```
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity
import org.springframework.security.config.annotation.web.invoke
import org.springframework.security.web.SecurityFilterChain
@Configuration(proxyBeanMethods = false)
@EnableWebSecurity
open class MyOAuthClientConfiguration {
	@Bean
	open fun securityFilterChain(http: HttpSecurity): SecurityFilterChain {
 http {
 authorizeHttpRequests {
 authorize(anyRequest, authenticated)
 }
 oauth2Login {
 redirectionEndpoint {
 baseUri = "/login/oauth2/callback/*"
 }
 }
 }
 return http.build()
	}
}
```
| Spring Boot auto-configures an `InMemoryOAuth2AuthorizedClientService`which is used by Spring Security for the management of client registrations.
The`InMemoryOAuth2AuthorizedClientService`has limited capabilities and we recommend using it only for development environments.
For production environments, consider using a`JdbcOAuth2AuthorizedClientService`or creating your own implementation of`OAuth2AuthorizedClientService`. |

### OAuth2 Client Registration for Common Providers

For common OAuth2 and OpenID providers (Google, Github, Facebook, and Okta), we provide a set of provider defaults.
The IDs of these common providers are `google`, `github`, `facebook`, and `okta`, respectively.

If you do not need to customize these providers, set the registration’s `provider` property to the ID of one of the common providers.
Alternatively, you can use a registration ID that matches the ID of the provider.
The two configurations in the following example use the common `google` provider:

-
Properties
-
YAML

```
spring.security.oauth2.client.registration.my-client.client-id=abcd
spring.security.oauth2.client.registration.my-client.client-secret=password
spring.security.oauth2.client.registration.my-client.provider=google
spring.security.oauth2.client.registration.google.client-id=abcd
spring.security.oauth2.client.registration.google.client-secret=password
```
```
spring:
 security:
 oauth2:
 client:
 registration:
 my-client:
 client-id: "abcd"
 client-secret: "password"
 provider: "google"
 google:
 client-id: "abcd"
 client-secret: "password"
```
## Resource Server

If you have `spring-security-oauth2-resource-server` on your classpath, Spring Boot can set up an OAuth2 Resource Server.
For JWT configuration, a JWK Set URI or OIDC Issuer URI needs to be specified, as shown in the following examples:

-
Properties
-
YAML

`spring.security.oauth2.resourceserver.jwt.jwk-set-uri=https://example.com/oauth2/default/v1/keys````
spring:
 security:
 oauth2:
 resourceserver:
 jwt:
 jwk-set-uri: "https://example.com/oauth2/default/v1/keys"
```
-
Properties
-
YAML

`spring.security.oauth2.resourceserver.jwt.issuer-uri=https://dev-123456.oktapreview.com/oauth2/default/````
spring:
 security:
 oauth2:
 resourceserver:
 jwt:
 issuer-uri: "https://dev-123456.oktapreview.com/oauth2/default/"
```
| If the authorization server does not support a JWK Set URI, you can configure the resource server with the Public Key used for verifying the signature of the JWT.
This can be done using the `spring.security.oauth2.resourceserver.jwt.public-key-location`property, where the value needs to point to a file containing the public key in the PEM-encoded x509 format. |

The `spring.security.oauth2.resourceserver.jwt.audiences` property can be used to specify the expected values of the aud claim in JWTs.
For example, to require JWTs to contain an aud claim with the value `my-audience`:

-
Properties
-
YAML

`spring.security.oauth2.resourceserver.jwt.audiences[0]=my-audience````
spring:
 security:
 oauth2:
 resourceserver:
 jwt:
 audiences:
 - "my-audience"
```
The same properties are applicable for both servlet and reactive applications.
Alternatively, you can define your own `JwtDecoder` bean for servlet applications or a `ReactiveJwtDecoder` for reactive applications.

In cases where opaque tokens are used instead of JWTs, you can configure the following properties to validate tokens through introspection:

-
Properties
-
YAML

```
spring.security.oauth2.resourceserver.opaquetoken.introspection-uri=https://example.com/check-token
spring.security.oauth2.resourceserver.opaquetoken.client-id=my-client-id
spring.security.oauth2.resourceserver.opaquetoken.client-secret=my-client-secret
```
```
spring:
 security:
 oauth2:
 resourceserver:
 opaquetoken:
 introspection-uri: "https://example.com/check-token"
 client-id: "my-client-id"
 client-secret: "my-client-secret"
```
Again, the same properties are applicable for both servlet and reactive applications.

The result is an auto-configured introspector. Either a `SpringOpaqueTokenIntrospector` or, in a reactive application, a `SpringReactiveOpaqueTokenIntrospector`.
These auto-configured introspectors can be customized using `SpringOpaqueTokenIntrospectorBuilderCustomizer` and `SpringReactiveOpaqueTokenIntrospectorBuilderCustomizer` beans respectively.

To take complete control of the introspection, define your own `OpaqueTokenIntrospector` or `ReactiveOpaqueTokenIntrospector` bean.

## Authorization Server

If you have `spring-security-oauth2-authorization-server` on your classpath, you can take advantage of some auto-configuration to set up a Servlet-based OAuth2 Authorization Server.

You can register multiple OAuth2 clients under the `spring.security.oauth2.authorizationserver.client` prefix, as shown in the following example:

-
Properties
-
YAML

```
spring.security.oauth2.authorizationserver.client.my-client-1.registration.client-id=abcd
spring.security.oauth2.authorizationserver.client.my-client-1.registration.client-secret={noop}secret1
spring.security.oauth2.authorizationserver.client.my-client-1.registration.client-authentication-methods[0]=client_secret_basic
spring.security.oauth2.authorizationserver.client.my-client-1.registration.authorization-grant-types[0]=authorization_code
spring.security.oauth2.authorizationserver.client.my-client-1.registration.authorization-grant-types[1]=refresh_token
spring.security.oauth2.authorizationserver.client.my-client-1.registration.redirect-uris[0]=https://my-client-1.com/login/oauth2/code/abcd
spring.security.oauth2.authorizationserver.client.my-client-1.registration.redirect-uris[1]=https://my-client-1.com/authorized
spring.security.oauth2.authorizationserver.client.my-client-1.registration.scopes[0]=openid
spring.security.oauth2.authorizationserver.client.my-client-1.registration.scopes[1]=profile
spring.security.oauth2.authorizationserver.client.my-client-1.registration.scopes[2]=email
spring.security.oauth2.authorizationserver.client.my-client-1.registration.scopes[3]=phone
spring.security.oauth2.authorizationserver.client.my-client-1.registration.scopes[4]=address
spring.security.oauth2.authorizationserver.client.my-client-1.require-authorization-consent=true
spring.security.oauth2.authorizationserver.client.my-client-1.token.authorization-code-time-to-live=5m
spring.security.oauth2.authorizationserver.client.my-client-1.token.access-token-time-to-live=10m
spring.security.oauth2.authorizationserver.client.my-client-1.token.access-token-format=reference
spring.security.oauth2.authorizationserver.client.my-client-1.token.reuse-refresh-tokens=false
spring.security.oauth2.authorizationserver.client.my-client-1.token.refresh-token-time-to-live=30m
spring.security.oauth2.authorizationserver.client.my-client-2.registration.client-id=efgh
spring.security.oauth2.authorizationserver.client.my-client-2.registration.client-secret={noop}secret2
spring.security.oauth2.authorizationserver.client.my-client-2.registration.client-authentication-methods[0]=client_secret_jwt
spring.security.oauth2.authorizationserver.client.my-client-2.registration.authorization-grant-types[0]=client_credentials
spring.security.oauth2.authorizationserver.client.my-client-2.registration.scopes[0]=user.read
spring.security.oauth2.authorizationserver.client.my-client-2.registration.scopes[1]=user.write
spring.security.oauth2.authorizationserver.client.my-client-2.jwk-set-uri=https://my-client-2.com/jwks
spring.security.oauth2.authorizationserver.client.my-client-2.token-endpoint-authentication-signing-algorithm=RS256
```
```
spring:
 security:
 oauth2:
 authorizationserver:
 client:
 my-client-1:
 registration:
 client-id: "abcd"
 client-secret: "{noop}secret1"
 client-authentication-methods:
 - "client_secret_basic"
 authorization-grant-types:
 - "authorization_code"
 - "refresh_token"
 redirect-uris:
 - "https://my-client-1.com/login/oauth2/code/abcd"
 - "https://my-client-1.com/authorized"
 scopes:
 - "openid"
 - "profile"
 - "email"
 - "phone"
 - "address"
 require-authorization-consent: true
 token:
 authorization-code-time-to-live: 5m
 access-token-time-to-live: 10m
 access-token-format: "reference"
 reuse-refresh-tokens: false
 refresh-token-time-to-live: 30m
 my-client-2:
 registration:
 client-id: "efgh"
 client-secret: "{noop}secret2"
 client-authentication-methods:
 - "client_secret_jwt"
 authorization-grant-types:
 - "client_credentials"
 scopes:
 - "user.read"
 - "user.write"
 jwk-set-uri: "https://my-client-2.com/jwks"
 token-endpoint-authentication-signing-algorithm: "RS256"
```
| The `client-secret`property must be in a format that can be matched by the configured`PasswordEncoder`.
The default instance of`PasswordEncoder`is created via`PasswordEncoderFactories.createDelegatingPasswordEncoder()`. |

The auto-configuration Spring Boot provides for Spring Authorization Server is designed for getting started quickly. Most applications will require customization and will want to define several beans to override auto-configuration.

The following components can be defined as beans to override auto-configuration specific to Spring Authorization Server:

-
`com.nimbusds.jose.jwk.source.JWKSource<com.nimbusds.jose.proc.SecurityContext>`

| Spring Boot auto-configures an `InMemoryRegisteredClientRepository`which is used by Spring Authorization Server for the management of registered clients.
The`InMemoryRegisteredClientRepository`has limited capabilities and we recommend using it only for development environments.
For production environments, consider using a`JdbcRegisteredClientRepository`or creating your own implementation of`RegisteredClientRepository`. |

Additional information can be found in the Getting Started chapter of Spring Security Reference Documentation.

# SAML 2.0

SAML v2.0 is a widely adopted framework for exchanging security information between online business partners.

## Build Configuration

SAML 2.0 support builds off of the OpenSAML library that requires an extra repository configuration.

### Using Maven

With Maven, you need to add an extra `repository` element to your POM as follows:

```
<repositories>
	<repository>
 <id>shibboleth-releases</id>
 <name>Shibboleth Releases Repository</name>
 <url>https://build.shibboleth.net/maven/releases</url>
 <snapshots>
 <enabled>false</enabled>
 </snapshots>
	</repository>
</repositories>
```
## Relying Party

If you have `spring-security-saml2-service-provider` on your classpath, you can take advantage of some auto-configuration to set up a SAML 2.0 Relying Party.
This configuration makes use of the properties under `Saml2RelyingPartyProperties`.

A relying party registration represents a paired configuration between an Identity Provider, IDP, and a Service Provider, SP.
You can register multiple relying parties under the `spring.security.saml2.relyingparty` prefix, as shown in the following example:

-
Properties
-
YAML

```
spring.security.saml2.relyingparty.registration.my-relying-party1.signing.credentials[0].private-key-location=path-to-private-key
spring.security.saml2.relyingparty.registration.my-relying-party1.signing.credentials[0].certificate-location=path-to-certificate
spring.security.saml2.relyingparty.registration.my-relying-party1.decryption.credentials[0].private-key-location=path-to-private-key
spring.security.saml2.relyingparty.registration.my-relying-party1.decryption.credentials[0].certificate-location=path-to-certificate
spring.security.saml2.relyingparty.registration.my-relying-party1.singlelogout.url=https://myapp/logout/saml2/slo
spring.security.saml2.relyingparty.registration.my-relying-party1.singlelogout.response-url=https://remoteidp2.slo.url
spring.security.saml2.relyingparty.registration.my-relying-party1.singlelogout.binding=POST
spring.security.saml2.relyingparty.registration.my-relying-party1.assertingparty.verification.credentials[0].certificate-location=path-to-verification-cert
spring.security.saml2.relyingparty.registration.my-relying-party1.assertingparty.entity-id=remote-idp-entity-id1
spring.security.saml2.relyingparty.registration.my-relying-party1.assertingparty.sso-url=https://remoteidp1.sso.url
spring.security.saml2.relyingparty.registration.my-relying-party2.signing.credentials[0].private-key-location=path-to-private-key
spring.security.saml2.relyingparty.registration.my-relying-party2.signing.credentials[0].certificate-location=path-to-certificate
spring.security.saml2.relyingparty.registration.my-relying-party2.decryption.credentials[0].private-key-location=path-to-private-key
spring.security.saml2.relyingparty.registration.my-relying-party2.decryption.credentials[0].certificate-location=path-to-certificate
spring.security.saml2.relyingparty.registration.my-relying-party2.assertingparty.verification.credentials[0].certificate-location=path-to-other-verification-cert
spring.security.saml2.relyingparty.registration.my-relying-party2.assertingparty.entity-id=remote-idp-entity-id2
spring.security.saml2.relyingparty.registration.my-relying-party2.assertingparty.sso-url=https://remoteidp2.sso.url
spring.security.saml2.relyingparty.registration.my-relying-party2.assertingparty.singlelogout.url=https://remoteidp2.slo.url
spring.security.saml2.relyingparty.registration.my-relying-party2.assertingparty.singlelogout.response-url=https://myapp/logout/saml2/slo
spring.security.saml2.relyingparty.registration.my-relying-party2.assertingparty.singlelogout.binding=POST
```
```
spring:
 security:
 saml2:
 relyingparty:
 registration:
 my-relying-party1:
 signing:
 credentials:
 - private-key-location: "path-to-private-key"
 certificate-location: "path-to-certificate"
 decryption:
 credentials:
 - private-key-location: "path-to-private-key"
 certificate-location: "path-to-certificate"
 singlelogout:
 url: "https://myapp/logout/saml2/slo"
 response-url: "https://remoteidp2.slo.url"
 binding: "POST"
 assertingparty:
 verification:
 credentials:
 - certificate-location: "path-to-verification-cert"
 entity-id: "remote-idp-entity-id1"
 sso-url: "https://remoteidp1.sso.url"
 my-relying-party2:
 signing:
 credentials:
 - private-key-location: "path-to-private-key"
 certificate-location: "path-to-certificate"
 decryption:
 credentials:
 - private-key-location: "path-to-private-key"
 certificate-location: "path-to-certificate"
 assertingparty:
 verification:
 credentials:
 - certificate-location: "path-to-other-verification-cert"
 entity-id: "remote-idp-entity-id2"
 sso-url: "https://remoteidp2.sso.url"
 singlelogout:
 url: "https://remoteidp2.slo.url"
 response-url: "https://myapp/logout/saml2/slo"
 binding: "POST"
```
For SAML2 logout, by default, Spring Security’s `Saml2LogoutRequestFilter` and `Saml2LogoutResponseFilter` only process URLs matching `/logout/saml2/slo`.
If you want to customize the `url` to which AP-initiated logout requests get sent to or the `response-url` to which an AP sends logout responses to, to use a different pattern, you need to provide configuration to process that custom pattern.
For example, for servlet applications, you can add your own `SecurityFilterChain` that resembles the following:

```
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.SecurityFilterChain;
import static org.springframework.security.config.Customizer.withDefaults;
@Configuration(proxyBeanMethods = false)
public class MySamlRelyingPartyConfiguration {
	@Bean
	public SecurityFilterChain securityFilterChain(HttpSecurity http) {
 http.authorizeHttpRequests((requests) -> requests.anyRequest().authenticated());
 http.saml2Login(withDefaults());
 http.saml2Logout((saml2) -> {
 saml2.logoutRequest((request) -> request.logoutUrl("/SLOService.saml2"));
 saml2.logoutResponse((response) -> response.logoutUrl("/SLOService.saml2"));
 });
 return http.build();
	}
}
```

# Testing

Spring Boot provides a number of utilities and annotations to help when testing your application.

Test support is provided by two general-purpose modules – `spring-boot-test` contains core items and `spring-boot-test-autoconfigure` supports auto-configuration for tests – and several focused `-test` modules that provide testing support for a particular feature.

Most developers use the `spring-boot-starter-test` starter, which imports both general-purpose Spring Boot test modules as well as JUnit Jupiter, AssertJ, Hamcrest, and a number of other useful libraries, and the focused `-test` modules that are applicable to their particular application.

| If you have tests that use JUnit 4, JUnit 6’s vintage engine can be used to run them.
To use the vintage engine, add a dependency on |

`hamcrest-core` is excluded in favor of `org.hamcrest:hamcrest` that is part of `spring-boot-starter-test`.

# Test Modules

Spring Boot offers several focused, feature-specific `-test` modules:

| Module | Purpose |
|---|---|
|
 | Testing applications that use Spring Framework’s cache abstraction. |
|
 | Testing applications that use Spring Data Cassandra. Provides the |
|
 | Testing applications that use Spring Data Couchbase. Provides the |
|
 | Testing applications that use Spring Data Elasticsearch. Provides the |
|
 | Testing applications that use Spring Data JDBC. Provides the |
|
 | Testing applications that use Spring Data JPA. Provides the |
|
 | Testing applications that use Spring Data LDAP. Provides the |
|
 | Testing applications that use Spring Data MongoDB. Provides the |
|
 | Testing applications that use Spring Data JPA. Provides the |
|
 | Testing applications that use Spring Data R2DBC. Provides the |
|
 | Testing applications that use Spring Data Redis. Provides the |
|
 | Testing applications that use Spring GraphQL. Provides the |
|
 | Testing applications that use Spring gRPC. |
|
 | Testing applications that using Spring JDBC. Provides the |
|
 | Testing applications that using jOOQ. Provides the |
|
 | Testing applications that use JPA. |
|
 | Testing applications that use Micrometer Metrics. |
|
 | Testing applications that use Micrometer Tracing. |
|
 | Testing applications that use REST clients. Provides the |
|
 | Testing applications that use Spring Security. |
|
 | Testing applications that use |
|
 | Testing applications that use |
|
 | Testing applications that use Spring WebFlux. Provides the |
|
 | Testing applications that use Spring Web MVC. Provides the |
|
 | Testing applications that use Spring Web Services. Provides the |

# Test Scope Dependencies

The `spring-boot-starter-test` starter (in the `test` `scope`) contains the following provided libraries:

-
JUnit: The de-facto standard for unit testing Java applications.
-
Spring Test & Spring Boot Test: Utilities and integration test support for Spring Boot applications.
-
AssertJ: A fluent assertion library.
-
Hamcrest: A library of matcher objects (also known as constraints or predicates).
-
Mockito: A Java mocking framework.
-
JSONassert: An assertion library for JSON.
-
JsonPath: XPath for JSON.
-
Awaitility: A library for testing asynchronous systems.

We generally find these common libraries to be useful when writing tests. If these libraries do not suit your needs, you can add additional test dependencies of your own.

# Testing Spring Applications

One of the major advantages of dependency injection is that it should make your code easier to unit test.
You can instantiate objects by using the `new` operator without even involving Spring.
You can also use *mock objects* instead of real dependencies.

Often, you need to move beyond unit testing and start integration testing (with a Spring `ApplicationContext`).
It is useful to be able to perform integration testing without requiring deployment of your application or needing to connect to other infrastructure.

The Spring Framework includes a dedicated test module for such integration testing.
You can declare a dependency directly to `org.springframework:spring-test` or use the `spring-boot-starter-test` starter to pull it in transitively.

If you have not used the `spring-test` module before, you should start by reading the relevant section of the Spring Framework reference documentation.

# Testing Spring Boot Applications

A Spring Boot application is a Spring `ApplicationContext`, so nothing very special has to be done to test it beyond what you would normally do with a vanilla Spring context.

| External properties, logging, and other features of Spring Boot are installed in the context by default only if you use `SpringApplication`to create it. |

Spring Boot provides a `@SpringBootTest` annotation, which can be used as an alternative to the standard `spring-test` `@ContextConfiguration` annotation when you need Spring Boot features.
The annotation works by creating the `ApplicationContext` used in your tests through `SpringApplication`.
In addition to `@SpringBootTest` a number of other annotations are also provided for testing more specific slices of an application.

| If you are using JUnit 4, do not forget to also add `@RunWith(SpringRunner.class)`to your test, otherwise the annotations will be ignored.
If you are using JUnit 6, there is no need to add the equivalent`@ExtendWith(SpringExtension.class)`as`@SpringBootTest`and the other`@…Test`annotations are already annotated with it. |

By default, `@SpringBootTest` will not start a server.
You can use the `webEnvironment` attribute of `@SpringBootTest` to further refine how your tests run:

-
`MOCK`(Default) : Loads a web`ApplicationContext`and provides a mock web environment. Embedded servers are not started when using this annotation. If a web environment is not available on your classpath, this mode transparently falls back to creating a regular non-web`ApplicationContext`. It can be used in conjunction with`@AutoConfigureMockMvc`, or`@AutoConfigureWebTestClient`] for mock-based testing of your web application.
-
`RANDOM_PORT`: Loads a`WebServerApplicationContext`and provides a real web environment. Embedded servers are started and listen on a random port.
-
`DEFINED_PORT`: Loads a`WebServerApplicationContext`and provides a real web environment. Embedded servers are started and listen on a defined port (from your`application.properties`) or on the default port of`8080`.
-
`NONE`: Loads an`ApplicationContext`by using`SpringApplication`but does not provide*any*web environment (mock or otherwise).

| If your test is `@Transactional`, it rolls back the transaction at the end of each test method by default.
However, as using this arrangement with either`RANDOM_PORT`or`DEFINED_PORT`implicitly provides a real servlet environment, the HTTP client and server run in separate threads and, thus, in separate transactions.
Any transaction initiated on the server does not roll back in this case. |

| `@SpringBootTest`with`webEnvironment = WebEnvironment.RANDOM_PORT`will also start the management server on a separate random port if your application uses a different port for the management server. |

## Detecting Web Application Type

If Spring MVC is available, a regular MVC-based application context is configured. If you have only Spring WebFlux, we will detect that and configure a WebFlux-based application context instead.

If both are present, Spring MVC takes precedence.
If you want to test a reactive web application in this scenario, you must set the `spring.main.web-application-type` property:

-
Java
-
Kotlin

```
import org.springframework.boot.test.context.SpringBootTest;
@SpringBootTest(properties = "spring.main.web-application-type=reactive")
class MyWebFluxTests {
	// ...
}
```
```
import org.springframework.boot.test.context.SpringBootTest
@SpringBootTest(properties = ["spring.main.web-application-type=reactive"])
class MyWebFluxTests {
	// ...
}
```
## Detecting Test Configuration

If you are familiar with the Spring Test Framework, you may be used to using `@ContextConfiguration(classes=…)` in order to specify which Spring `@Configuration` to load.
Alternatively, you might have often used nested `@Configuration` classes within your test.

When testing Spring Boot applications, this is often not required.
Spring Boot’s `@*Test` annotations search for your primary configuration automatically whenever you do not explicitly define one.

The search algorithm works up from the package that contains the test until it finds a class annotated with `@SpringBootApplication` or `@SpringBootConfiguration`.
As long as you structured your code in a sensible way, your main configuration is usually found.

| If you use a test annotation to test a more specific slice of your application, you should avoid adding configuration settings that are specific to a particular area on the main method’s application class. The underlying component scan configuration of |

If you want to customize the primary configuration, you can use a nested `@TestConfiguration` class.
Unlike a nested `@Configuration` class, which would be used instead of your application’s primary configuration, a nested `@TestConfiguration` class is used in addition to your application’s primary configuration.

| Spring’s test framework caches application contexts between tests. Therefore, as long as your tests share the same configuration (no matter how it is discovered), the potentially time-consuming process of loading the context happens only once. |

## Using the Test Configuration Main Method

Typically the test configuration discovered by `@SpringBootTest` will be your main `@SpringBootApplication`.
In most well structured applications, this configuration class will also include the `main` method used to launch the application.

For example, the following is a very common code pattern for a typical Spring Boot application:

-
Java
-
Kotlin

```
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
@SpringBootApplication
public class MyApplication {
	public static void main(String[] args) {
 SpringApplication.run(MyApplication.class, args);
	}
}
```
```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.docs.using.structuringyourcode.locatingthemainclass.MyApplication
import org.springframework.boot.runApplication
@SpringBootApplication
class MyApplication
fun main(args: Array<String>) {
	runApplication<MyApplication>(*args)
}
```
In the example above, the `main` method doesn’t do anything other than delegate to `SpringApplication.run(Class, String...)`.
It is, however, possible to have a more complex `main` method that applies customizations before calling `SpringApplication.run(Class, String...)`.

For example, here is an application that changes the banner mode and sets additional profiles:

-
Java
-
Kotlin

```
import org.springframework.boot.Banner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
@SpringBootApplication
public class MyApplication {
	public static void main(String[] args) {
 SpringApplication application = new SpringApplication(MyApplication.class);
 application.setBannerMode(Banner.Mode.OFF);
 application.setAdditionalProfiles("myprofile");
 application.run(args);
	}
}
```
```
import org.springframework.boot.Banner
import org.springframework.boot.runApplication
import org.springframework.boot.autoconfigure.SpringBootApplication
@SpringBootApplication
class MyApplication
fun main(args: Array<String>) {
	runApplication<MyApplication>(*args) {
 setBannerMode(Banner.Mode.OFF)
 setAdditionalProfiles("myprofile")
	}
}
```
Since customizations in the `main` method can affect the resulting `ApplicationContext`, it’s possible that you might also want to use the `main` method to create the `ApplicationContext` used in your tests.
By default, `@SpringBootTest` will not call your `main` method, and instead the class itself is used directly to create the `ApplicationContext`

If you want to change this behavior, you can change the `useMainMethod` attribute of `@SpringBootTest` to `SpringBootTest.UseMainMethod.ALWAYS` or `SpringBootTest.UseMainMethod.WHEN_AVAILABLE`.
When set to `ALWAYS`, the test will fail if no `main` method can be found.
When set to `WHEN_AVAILABLE` the `main` method will be used if it is available, otherwise the standard loading mechanism will be used.

For example, the following test will invoke the `main` method of `MyApplication` in order to create the `ApplicationContext`.
If the main method sets additional profiles then those will be active when the `ApplicationContext` starts.

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.UseMainMethod;
@SpringBootTest(useMainMethod = UseMainMethod.ALWAYS)
class MyApplicationTests {
	@Test
	void exampleTest() {
 // ...
	}
}
```
```
import org.junit.jupiter.api.Test
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.boot.test.context.SpringBootTest.UseMainMethod
@SpringBootTest(useMainMethod = UseMainMethod.ALWAYS)
class MyApplicationTests {
	@Test
	fun exampleTest() {
 // ...
	}
}
```
## Excluding Test Configuration

If your application uses component scanning (for example, if you use `@SpringBootApplication` or `@ComponentScan`), you may find top-level configuration classes that you created only for specific tests accidentally get picked up everywhere.

As we have seen earlier, `@TestConfiguration` can be used on an inner class of a test to customize the primary configuration.
`@TestConfiguration` can also be used on a top-level class. Doing so indicates that the class should not be picked up by scanning.
You can then import the class explicitly where it is required, as shown in the following example:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Import;
@SpringBootTest
@Import(MyTestsConfiguration.class)
class MyTests {
	@Test
	void exampleTest() {
 // ...
	}
}
```
```
import org.junit.jupiter.api.Test
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.context.annotation.Import
@SpringBootTest
@Import(MyTestsConfiguration::class)
class MyTests {
	@Test
	fun exampleTest() {
 // ...
	}
}
```
| If you directly use `@ComponentScan`(that is, not through`@SpringBootApplication`) you need to register the`TypeExcludeFilter`with it.
See the`TypeExcludeFilter`API documentation for details. |

| An imported `@TestConfiguration`is processed earlier than an inner-class`@TestConfiguration`and an imported`@TestConfiguration`will be processed before any configuration found through component scanning.
Generally speaking, this difference in ordering has no noticeable effect but it is something to be aware of if you’re relying on bean overriding. |

## Using Application Arguments

If your application expects arguments, you can
have `@SpringBootTest` inject them using the `args` attribute.

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.test.context.SpringBootTest;
import static org.assertj.core.api.Assertions.assertThat;
@SpringBootTest(args = "--app.test=one")
class MyApplicationArgumentTests {
	@Test
	void applicationArgumentsPopulated(@Autowired ApplicationArguments args) {
 assertThat(args.getOptionNames()).containsOnly("app.test");
 assertThat(args.getOptionValues("app.test")).containsOnly("one");
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.ApplicationArguments
import org.springframework.boot.test.context.SpringBootTest
@SpringBootTest(args = ["--app.test=one"])
class MyApplicationArgumentTests {
	@Test
	fun applicationArgumentsPopulated(@Autowired args: ApplicationArguments) {
 assertThat(args.optionNames).containsOnly("app.test")
 assertThat(args.getOptionValues("app.test")).containsOnly("one")
	}
}
```
## Testing With a Mock Environment

By default, `@SpringBootTest` does not start the server but instead sets up a mock environment for testing web endpoints.

With Spring MVC, we can query our web endpoints using `MockMvc`.
The following integrations are available:

-
The regular `MockMvc`that uses Hamcrest.
-
`MockMvcTester`that wraps`MockMvc`and uses AssertJ.
-
`RestTestClient`where`MockMvc`is plugged in as the server to handle requests with.
-
`WebTestClient`where`MockMvc`is plugged in as the server to handle requests with.

The following example showcases the available integrations:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.resttestclient.autoconfigure.AutoConfigureRestTestClient;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.webmvc.test.autoconfigure.AutoConfigureMockMvc;
import org.springframework.boot.webtestclient.autoconfigure.AutoConfigureWebTestClient;
import org.springframework.test.web.reactive.server.WebTestClient;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.assertj.MockMvcTester;
import org.springframework.test.web.servlet.client.RestTestClient;
import org.springframework.test.web.servlet.client.RestTestClient.ResponseSpec;
import org.springframework.test.web.servlet.client.assertj.RestTestClientResponse;
import static org.assertj.core.api.Assertions.assertThat;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.content;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;
@SpringBootTest
@AutoConfigureMockMvc
@AutoConfigureRestTestClient
@AutoConfigureWebTestClient
class MyMockMvcTests {
	@Test
	void testWithMockMvc(@Autowired MockMvc mvc) throws Exception {
 mvc.perform(get("/"))
 .andExpect(status().isOk())
 .andExpect(content().string("Hello World"));
	}
	@Test // If AssertJ is on the classpath, you can use MockMvcTester
	void testWithMockMvcTester(@Autowired MockMvcTester mvc) {
 assertThat(mvc.get().uri("/"))
 .hasStatusOk()
 .hasBodyTextEqualTo("Hello World");
	}
	@Test
	void testWithRestTestClient(@Autowired RestTestClient restClient) {
 restClient
 .get().uri("/")
 .exchange()
 .expectStatus().isOk()
 .expectBody(String.class).isEqualTo("Hello World");
	}
	@Test // If you prefer AssertJ, dedicated assertions are available
	void testWithRestTestClientAssertJ(@Autowired RestTestClient restClient) {
 ResponseSpec spec = restClient.get().uri("/").exchange();
 RestTestClientResponse response = RestTestClientResponse.from(spec);
 assertThat(response).hasStatusOk()
 .bodyText().isEqualTo("Hello World");
	}
	@Test // If Spring WebFlux is on the classpath
	void testWithWebTestClient(@Autowired WebTestClient webClient) {
 webClient
 .get().uri("/")
 .exchange()
 .expectStatus().isOk()
 .expectBody(String.class).isEqualTo("Hello World");
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.resttestclient.autoconfigure.AutoConfigureRestTestClient
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.boot.webmvc.test.autoconfigure.AutoConfigureMockMvc
import org.springframework.boot.webtestclient.autoconfigure.AutoConfigureWebTestClient
import org.springframework.test.web.reactive.server.WebTestClient
import org.springframework.test.web.reactive.server.expectBody
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.assertj.MockMvcTester
import org.springframework.test.web.servlet.client.RestTestClient
import org.springframework.test.web.servlet.client.assertj.RestTestClientResponse
import org.springframework.test.web.servlet.client.expectBody
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get
import org.springframework.test.web.servlet.result.MockMvcResultMatchers.content
import org.springframework.test.web.servlet.result.MockMvcResultMatchers.status
@SpringBootTest
@AutoConfigureMockMvc
@AutoConfigureRestTestClient
@AutoConfigureWebTestClient
class MyMockMvcTests {
	@Test
	fun testWithMockMvc(@Autowired mvc: MockMvc) {
 mvc.perform(get("/"))
 .andExpect(status().isOk())
 .andExpect(content().string("Hello World"))
	}
	@Test // If AssertJ is on the classpath, you can use MockMvcTester
	fun testWithMockMvcTester(@Autowired mvc: MockMvcTester) {
 assertThat(mvc.get().uri("/")).hasStatusOk()
 .hasBodyTextEqualTo("Hello World")
	}
	@Test
	fun testWithRestTestClient(@Autowired webClient: RestTestClient) {
 webClient
 .get().uri("/")
 .exchange()
 .expectStatus().isOk
 .expectBody<String>().isEqualTo("Hello World")
	}
	@Test // If you prefer AssertJ, dedicated assertions are available
	fun testWithRestTestClientAssertJ(@Autowired webClient: RestTestClient) {
 val spec = webClient.get().uri("/").exchange()
 val response = RestTestClientResponse.from(spec)
 assertThat(response).hasStatusOk().bodyText().isEqualTo("Hello World")
	}
	@Test // If Spring WebFlux is on the classpath
	fun testWithWebTestClient(@Autowired webClient: WebTestClient) {
 webClient
 .get().uri("/")
 .exchange()
 .expectStatus().isOk
 .expectBody<String>().isEqualTo("Hello World")
	}
}
```
| If you want to focus only on the web layer and not start a complete `ApplicationContext`, consider using`@WebMvcTest`instead. |

With Spring WebFlux endpoints, you can use `WebTestClient` as shown in the following example:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.webtestclient.autoconfigure.AutoConfigureWebTestClient;
import org.springframework.test.web.reactive.server.WebTestClient;
@SpringBootTest
@AutoConfigureWebTestClient
class MyMockWebTestClientTests {
	@Test
	void exampleTest(@Autowired WebTestClient webClient) {
 webClient
 .get().uri("/")
 .exchange()
 .expectStatus().isOk()
 .expectBody(String.class).isEqualTo("Hello World");
	}
}
```
```
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.boot.webtestclient.autoconfigure.AutoConfigureWebTestClient
import org.springframework.test.web.reactive.server.WebTestClient
import org.springframework.test.web.reactive.server.expectBody
@SpringBootTest
@AutoConfigureWebTestClient
class MyMockWebTestClientTests {
	@Test
	fun exampleTest(@Autowired webClient: WebTestClient) {
 webClient
 .get().uri("/")
 .exchange()
 .expectStatus().isOk
 .expectBody<String>().isEqualTo("Hello World")
	}
}
```
| Testing within a mocked environment is usually faster than running with a full servlet container. However, since mocking occurs at the Spring MVC layer, code that relies on lower-level servlet container behavior cannot be directly tested with MockMvc. For example, Spring Boot’s error handling is based on the “error page” support provided by the servlet container. This means that, whilst you can test your MVC layer throws and handles exceptions as expected, you cannot directly test that a specific custom error page is rendered. If you need to test these lower-level concerns, you can start a fully running server as described in the next section. |

## Testing With a Running Server

If you need to start a full running server, we recommend that you use random ports.
If you use `@SpringBootTest(webEnvironment=WebEnvironment.RANDOM_PORT)`, an available port is picked at random each time your test runs.

The `@LocalServerPort` annotation can be used to inject the actual port used into your test.

Tests that need to make REST calls to the started server can autowire a
`RestTestClient` by annotating the test class with `@AutoConfigureRestTestClient`.

The configured client resolves relative links to the running server and comes with a dedicated API for verifying responses, as shown in the following example:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.resttestclient.autoconfigure.AutoConfigureRestTestClient;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment;
import org.springframework.test.web.servlet.client.RestTestClient;
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@AutoConfigureRestTestClient
class MyRandomPortRestTestClientTests {
	@Test
	void exampleTest(@Autowired RestTestClient restClient) {
 restClient
 .get().uri("/")
 .exchange()
 .expectStatus().isOk()
 .expectBody(String.class).isEqualTo("Hello World");
	}
}
```
```
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.resttestclient.autoconfigure.AutoConfigureRestTestClient
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment
import org.springframework.test.web.servlet.client.RestTestClient
import org.springframework.test.web.servlet.client.expectBody
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@AutoConfigureRestTestClient
class MyRandomPortRestTestClientTests {
	@Test
	fun exampleTest(@Autowired webClient: RestTestClient) {
 webClient
 .get().uri("/")
 .exchange()
 .expectStatus().isOk
 .expectBody<String>().isEqualTo("Hello World")
	}
}
```
If you prefer to use AssertJ, dedicated assertions are available from `RestTestClientResponse`, as shown in the following example:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.resttestclient.autoconfigure.AutoConfigureRestTestClient;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment;
import org.springframework.test.web.servlet.client.RestTestClient;
import org.springframework.test.web.servlet.client.RestTestClient.ResponseSpec;
import org.springframework.test.web.servlet.client.assertj.RestTestClientResponse;
import static org.assertj.core.api.Assertions.assertThat;
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@AutoConfigureRestTestClient
class MyRandomPortRestTestClientAssertJTests {
	@Test
	void exampleTest(@Autowired RestTestClient restClient) {
 ResponseSpec spec = restClient.get().uri("/").exchange();
 RestTestClientResponse response = RestTestClientResponse.from(spec);
 assertThat(response).hasStatusOk().bodyText().isEqualTo("Hello World");
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.resttestclient.autoconfigure.AutoConfigureRestTestClient
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment
import org.springframework.test.web.servlet.client.RestTestClient
import org.springframework.test.web.servlet.client.assertj.RestTestClientResponse
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@AutoConfigureRestTestClient
class MyRandomPortRestTestClientAssertJTests {
	@Test
	fun exampleTest(@Autowired webClient: RestTestClient) {
 val exchange = webClient.get().uri("/").exchange()
 val response = RestTestClientResponse.from(exchange)
 assertThat(response).hasStatusOk()
 .bodyText().isEqualTo("Hello World")
	}
}
```
If you have `spring-webflux` on the classpath, you can also autowire a `WebTestClient` by annotating the test class with `@AutoConfigureWebTestClient`.

`WebTestClient` provides a similar API, as shown in the following example:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment;
import org.springframework.boot.webtestclient.autoconfigure.AutoConfigureWebTestClient;
import org.springframework.test.web.reactive.server.WebTestClient;
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@AutoConfigureWebTestClient
class MyRandomPortWebTestClientTests {
	@Test
	void exampleTest(@Autowired WebTestClient webClient) {
 webClient
 .get().uri("/")
 .exchange()
 .expectStatus().isOk()
 .expectBody(String.class).isEqualTo("Hello World");
	}
}
```
```
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment
import org.springframework.boot.webtestclient.autoconfigure.AutoConfigureWebTestClient
import org.springframework.test.web.reactive.server.WebTestClient
import org.springframework.test.web.reactive.server.expectBody
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@AutoConfigureWebTestClient
class MyRandomPortWebTestClientTests {
	@Test
	fun exampleTest(@Autowired webClient: WebTestClient) {
 webClient
 .get().uri("/")
 .exchange()
 .expectStatus().isOk
 .expectBody<String>().isEqualTo("Hello World")
	}
}
```
| `WebTestClient`can also be used with a mock environment, removing the need for a running server, by annotating your test class with`@AutoConfigureWebTestClient`from`spring-boot-webflux-test`. |

| For certain sliced tests you may need to use the `@AutoConfigureWebServer`annotation to auto-configure an embedded web server. |

The `spring-boot-resttestclient` module also provides a `TestRestTemplate` facility:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.resttestclient.TestRestTemplate;
import org.springframework.boot.resttestclient.autoconfigure.AutoConfigureTestRestTemplate;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment;
import static org.assertj.core.api.Assertions.assertThat;
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@AutoConfigureTestRestTemplate
class MyRandomPortTestRestTemplateTests {
	@Test
	void exampleTest(@Autowired TestRestTemplate restTemplate) {
 String body = restTemplate.getForObject("/", String.class);
 assertThat(body).isEqualTo("Hello World");
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment
import org.springframework.boot.resttestclient.TestRestTemplate
import org.springframework.boot.resttestclient.autoconfigure.AutoConfigureTestRestTemplate
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@AutoConfigureTestRestTemplate
class MyRandomPortTestRestTemplateTests {
	@Test
	fun exampleTest(@Autowired restTemplate: TestRestTemplate) {
 val body = restTemplate.getForObject("/", String::class.java)
 assertThat(body).isEqualTo("Hello World")
	}
}
```
To use `TestRestTemplate` a dependency on `spring-boot-restclient` is also required.
Take care when adding this dependency as it will enable auto-configuration for `RestClient.Builder`.
If your main code uses `RestClient.Builder`, declare the `spring-boot-restclient` dependency so that it is on your application’s main classpath and not only on its test classpath.

## Customizing RestTestClient

To customize the `RestTestClient` bean, configure a `RestTestClientBuilderCustomizer` bean.
Any such beans are called with the `RestTestClient.Builder` that is used to create the `RestTestClient`.

## Customizing WebTestClient

To customize the `WebTestClient` bean, configure a `WebTestClientBuilderCustomizer` bean.
Any such beans are called with the `WebTestClient.Builder` that is used to create the `WebTestClient`.

## Using JMX

As the test context framework caches context, JMX is disabled by default to prevent identical components to register on the same domain.
If such test needs access to an `MBeanServer`, consider marking it dirty as well:

-
Java
-
Kotlin

```
import javax.management.MBeanServer;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.annotation.DirtiesContext;
import static org.assertj.core.api.Assertions.assertThat;
@SpringBootTest(properties = "spring.jmx.enabled=true")
@DirtiesContext
class MyJmxTests {
	@Autowired
	private MBeanServer mBeanServer;
	@Test
	void exampleTest() {
 assertThat(this.mBeanServer.getDomains()).contains("java.lang");
 // ...
	}
}
```
```
import javax.management.MBeanServer
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.test.annotation.DirtiesContext
@SpringBootTest(properties = ["spring.jmx.enabled=true"])
@DirtiesContext
class MyJmxTests(@Autowired val mBeanServer: MBeanServer) {
	@Test
	fun exampleTest() {
 assertThat(mBeanServer.domains).contains("java.lang")
 // ...
	}
}
```
## Using Observations

If you annotate a sliced test with `@AutoConfigureTracing` from `spring-boot-micrometer-tracing-test` or with `@AutoConfigureMetrics` from `spring-boot-micrometer-metrics-test`, it auto-configures an `ObservationRegistry`.

## Using Metrics

Regardless of your classpath, meter registries, except the in-memory backed, are not auto-configured when using `@SpringBootTest`.

If you need to export metrics to a different backend as part of an integration test, annotate it with `@AutoConfigureMetrics`.

If you annotate a sliced test with `@AutoConfigureMetrics`, it auto-configures an in-memory `MeterRegistry`.
Data exporting in sliced tests is not supported with the `@AutoConfigureMetrics` annotation.

## Using Tracing

Regardless of your classpath, tracing components which are reporting data are not auto-configured when using `@SpringBootTest`.

If you need those components as part of an integration test, annotate the test with `@AutoConfigureTracing`.

If you have created your own reporting components (e.g. a custom `SpanExporter` or `brave.handler.SpanHandler`) and you don’t want them to be active in tests, you can use the `@ConditionalOnEnabledTracingExport` annotation to disable them.

If you annotate a sliced test with `@AutoConfigureTracing` , it auto-configures a no-op `Tracer`.
Data exporting in sliced tests is not supported with the `@AutoConfigureTracing` annotation.

## Mocking and Spying Beans

When running tests, it is sometimes necessary to mock certain components within your application context. For example, you may have a facade over some remote service that is unavailable during development. Mocking can also be useful when you want to simulate failures that might be hard to trigger in a real environment.

Spring Framework includes a `@MockitoBean` annotation that can be used to define a Mockito mock for a bean inside your `ApplicationContext`.
Additionally, `@MockitoSpyBean` can be used to define a Mockito spy.
Learn more about these features in the Spring Framework documentation.

## Auto-configured Tests

Spring Boot’s auto-configuration system works well for applications but can sometimes be a little too much for tests. It often helps to load only the parts of the configuration that are required to test a “slice” of your application. For example, you might want to test that Spring MVC controllers are mapping URLs correctly, and you do not want to involve database calls in those tests, or you might want to test JPA entities, and you are not interested in the web layer when those tests run.

When combined with `spring-boot-test-autoconfigure`, Spring Boot’s test modules include a number of annotations that can be used to automatically configure such “slices”.
Each of them works in a similar way, providing a `@…Test` annotation that loads the `ApplicationContext` and one or more `@AutoConfigure…` annotations that can be used to customize auto-configuration settings.

| Each slice restricts component scan to appropriate components and loads a very restricted set of auto-configuration classes.
If you need to exclude one of them, most `@…Test`annotations provide an`excludeAutoConfiguration`attribute.
Alternatively, you can use`@ImportAutoConfiguration#exclude`. |

| Including multiple “slices” by using several `@…Test`annotations in one test is not supported.
If you need multiple “slices”, pick one of the`@…Test`annotations and include the`@AutoConfigure…`annotations of the other “slices” by hand. |

| It is also possible to use the `@AutoConfigure…`annotations with the standard`@SpringBootTest`annotation.
You can use this combination if you are not interested in “slicing” your application but you want some of the auto-configured test beans. |

## Auto-configured JSON Tests

To test that object JSON serialization and deserialization is working as expected, you can use the `@JsonTest` annotation from the `spring-boot-test-autoconfigure` module.
`@JsonTest` auto-configures the available supported JSON mapper, which can be one of the following libraries:

-
Jackson `JsonMapper`, any`@JacksonComponent`beans and any Jackson`JacksonModule`
-
Jackson 2 (deprecated) `ObjectMapper`, any`@JsonComponent`beans and any Jackson`Module`
-
`Gson`
-
`Jsonb`

| A list of the auto-configurations that are enabled by `@JsonTest`can be found in the appendix. |

If you need to configure elements of the auto-configuration, you can use the `@AutoConfigureJsonTesters` annotation.

Spring Boot includes AssertJ-based helpers that work with the JSONAssert and JsonPath libraries to check that JSON appears as expected.
The `JacksonTester`, `GsonTester`, `JsonbTester`, and `BasicJsonTester` classes can be used for Jackson, Gson, Jsonb, and Strings respectively.
Any helper fields on the test class can be `@Autowired` when using `@JsonTest`.
The following example shows a test class for Jackson:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.json.JsonTest;
import org.springframework.boot.test.json.JacksonTester;
import static org.assertj.core.api.Assertions.assertThat;
@JsonTest
class MyJsonTests {
	@Autowired
	private JacksonTester<VehicleDetails> json;
	@Test
	void serialize() throws Exception {
 VehicleDetails details = new VehicleDetails("Honda", "Civic");
 // Assert against a `.json` file in the same package as the test
 assertThat(this.json.write(details)).isEqualToJson("expected.json");
 // Or use JSON path based assertions
 assertThat(this.json.write(details)).hasJsonPathStringValue("@.make");
 assertThat(this.json.write(details)).extractingJsonPathStringValue("@.make").isEqualTo("Honda");
	}
	@Test
	void deserialize() throws Exception {
 String content = "{\"make\":\"Ford\",\"model\":\"Focus\"}";
 assertThat(this.json.parse(content)).isEqualTo(new VehicleDetails("Ford", "Focus"));
 assertThat(this.json.parseObject(content).getMake()).isEqualTo("Ford");
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.autoconfigure.json.JsonTest
import org.springframework.boot.test.json.JacksonTester
@JsonTest
class MyJsonTests(@Autowired val json: JacksonTester<VehicleDetails>) {
	@Test
	fun serialize() {
 val details = VehicleDetails("Honda", "Civic")
 // Assert against a `.json` file in the same package as the test
 assertThat(json.write(details)).isEqualToJson("expected.json")
 // Or use JSON path based assertions
 assertThat(json.write(details)).hasJsonPathStringValue("@.make")
 assertThat(json.write(details)).extractingJsonPathStringValue("@.make").isEqualTo("Honda")
	}
	@Test
	fun deserialize() {
 val content = "{\"make\":\"Ford\",\"model\":\"Focus\"}"
 assertThat(json.parse(content)).isEqualTo(VehicleDetails("Ford", "Focus"))
 assertThat(json.parseObject(content).make).isEqualTo("Ford")
	}
}
```
| JSON helper classes can also be used directly in standard unit tests.
To do so, call the `initFields`method of the helper in your`@BeforeEach`method if you do not use`@JsonTest`. |

If you use Spring Boot’s AssertJ-based helpers to assert on a number value at a given JSON path, you might not be able to use `isEqualTo` depending on the type.
Instead, you can use AssertJ’s `satisfies` to assert that the value matches the given condition.
For instance, the following example asserts that the actual number is a float value close to `0.15` within an offset of `0.01`.

-
Java
-
Kotlin

```
	@Test
	void someTest() throws Exception {
 SomeObject value = new SomeObject(0.152f);
 assertThat(this.json.write(value)).extractingJsonPathNumberValue("@.test.numberValue")
 .satisfies((number) -> assertThat(number.floatValue()).isCloseTo(0.15f, within(0.01f)));
	}
```
```
	@Test
	fun someTest() {
 val value = SomeObject(0.152f)
 assertThat(json.write(value)).extractingJsonPathNumberValue("@.test.numberValue")
 .satisfies(ThrowingConsumer { number ->
 assertThat(number.toFloat()).isCloseTo(0.15f, within(0.01f))
 })
	}
```
## Auto-configured Spring MVC Tests

To test whether Spring MVC controllers are working as expected, use the `@WebMvcTest` annotation from the `spring-boot-webmvc-test` module.
`@WebMvcTest` auto-configures the Spring MVC infrastructure and limits scanned beans to `@Controller`, `@ControllerAdvice`, `@JacksonComponent`, `@JsonComponent` (deprecated), `Converter`, `GenericConverter`, `Filter`, `HandlerInterceptor`, `WebMvcConfigurer`, `WebMvcRegistrations`, and `HandlerMethodArgumentResolver`.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@WebMvcTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.

| A list of the auto-configuration settings that are enabled by `@WebMvcTest`can be found in the appendix. |

| If you need to register extra components, such as a `JacksonModule`, you can import additional configuration classes by using`@Import`on your test. |

Often, `@WebMvcTest` is limited to a single controller and is used in combination with `@MockitoBean` to provide mock implementations for required collaborators.

`@WebMvcTest` also auto-configures `MockMvc`.
Mock MVC offers a powerful way to quickly test MVC controllers without needing to start a full HTTP server.
If AssertJ is available, the AssertJ support provided by `MockMvcTester` is auto-configured as well.
If you’d like to use `RestTestClient` in your tests, annotate your test class with `@AutoConfigureRestTestClient`.
A `RestTestClient` that uses the Mock MVC infrastructure will then be auto-configured.

| You can also auto-configure `MockMvc`and`MockMvcTester`in a non-`@WebMvcTest`(such as`@SpringBootTest`) by annotating it with`@AutoConfigureMockMvc`.
The following example uses`MockMvcTester`: |

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.webmvc.test.autoconfigure.WebMvcTest;
import org.springframework.http.MediaType;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.assertj.MockMvcTester;
import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.BDDMockito.given;
@WebMvcTest(UserVehicleController.class)
class MyControllerTests {
	@Autowired
	private MockMvcTester mvc;
	@MockitoBean
	private UserVehicleService userVehicleService;
	@Test
	void testExample() {
 given(this.userVehicleService.getVehicleDetails("sboot"))
 .willReturn(new VehicleDetails("Honda", "Civic"));
 assertThat(this.mvc.get().uri("/sboot/vehicle").accept(MediaType.TEXT_PLAIN))
 .hasStatusOk()
 .hasBodyTextEqualTo("Honda Civic");
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.webmvc.test.autoconfigure.WebMvcTest
import org.springframework.http.MediaType
import org.springframework.test.context.bean.override.mockito.MockitoBean
import org.springframework.test.web.servlet.assertj.MockMvcTester
@WebMvcTest(UserVehicleController::class)
class MyControllerTests(@Autowired val mvc: MockMvcTester) {
	@MockitoBean
	lateinit var userVehicleService: UserVehicleService
	@Test
	fun testExample() {
 given(userVehicleService.getVehicleDetails("sboot"))
 .willReturn(VehicleDetails("Honda", "Civic"))
 assertThat(mvc.get().uri("/sboot/vehicle").accept(MediaType.TEXT_PLAIN))
 .hasStatusOk().hasBodyTextEqualTo("Honda Civic")
	}
}
```
| If you need to configure elements of the auto-configuration (for example, when servlet filters should be applied) you can use attributes in the `@AutoConfigureMockMvc`annotation. |

If you use HtmlUnit and Selenium, auto-configuration also provides an HtmlUnit `WebClient` bean and/or a Selenium `WebDriver` bean.
The following example uses HtmlUnit:

-
Java
-
Kotlin

```
import org.htmlunit.WebClient;
import org.htmlunit.html.HtmlPage;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.webmvc.test.autoconfigure.WebMvcTest;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.BDDMockito.given;
@WebMvcTest(UserVehicleController.class)
class MyHtmlUnitTests {
	@Autowired
	private WebClient webClient;
	@MockitoBean
	private UserVehicleService userVehicleService;
	@Test
	void testExample() throws Exception {
 given(this.userVehicleService.getVehicleDetails("sboot")).willReturn(new VehicleDetails("Honda", "Civic"));
 HtmlPage page = this.webClient.getPage("/sboot/vehicle.html");
 assertThat(page.getBody().getTextContent()).isEqualTo("Honda Civic");
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.htmlunit.WebClient
import org.htmlunit.html.HtmlPage
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.webmvc.test.autoconfigure.WebMvcTest
import org.springframework.test.context.bean.override.mockito.MockitoBean
@WebMvcTest(UserVehicleController::class)
class MyHtmlUnitTests(@Autowired val webClient: WebClient) {
	@MockitoBean
	lateinit var userVehicleService: UserVehicleService
	@Test
	fun testExample() {
 given(userVehicleService.getVehicleDetails("sboot")).willReturn(VehicleDetails("Honda", "Civic"))
 val page = webClient.getPage<HtmlPage>("/sboot/vehicle.html")
 assertThat(page.body.textContent).isEqualTo("Honda Civic")
	}
}
```
| By default, Spring Boot puts `WebDriver`beans in a special “scope” to ensure that the driver exits after each test and that a new instance is injected.
If you do not want this behavior, you can add`@Scope(ConfigurableBeanFactory.SCOPE_SINGLETON)`to your`WebDriver``@Bean`definition. |

| The `webDriver`scope created by Spring Boot will replace any user defined scope of the same name.
If you define your own`webDriver`scope you may find it stops working when you use`@WebMvcTest`. |

If you have Spring Security on the classpath, `@WebMvcTest` will also scan `WebSecurityConfigurer` beans.
Instead of disabling security completely for such tests, you can use Spring Security’s test support.
More details on how to use Spring Security’s `MockMvc` support can be found in this Testing With Spring Security “How-to Guides” section.

| Sometimes writing Spring MVC tests is not enough; Spring Boot can help you run full end-to-end tests with an actual server. |

## Auto-configured Spring WebFlux Tests

To test that Spring WebFlux controllers are working as expected, you can use the `@WebFluxTest` annotation from the `spring-boot-webflux-test` module.
`@WebFluxTest` auto-configures the Spring WebFlux infrastructure and limits scanned beans to `@Controller`, `@ControllerAdvice`, `@JacksonComponent`, `@JsonComponent` (deprecated), `Converter`, `GenericConverter` and `WebFluxConfigurer`.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@WebFluxTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.

| A list of the auto-configurations that are enabled by `@WebFluxTest`can be found in the appendix. |

| If you need to register extra components, such as a `JacksonModule`, you can import additional configuration classes using`@Import`on your test. |

Often, `@WebFluxTest` is limited to a single controller and used in combination with the `@MockitoBean` annotation to provide mock implementations for required collaborators.

`@WebFluxTest` also auto-configures `WebTestClient`, which offers a powerful way to quickly test WebFlux controllers without needing to start a full HTTP server.

| You can also auto-configure `WebTestClient`in a non-`@WebFluxTest`(such as`@SpringBootTest`) by annotating it with`@AutoConfigureWebTestClient`. |

The following example shows a class that uses both `@WebFluxTest` and a `WebTestClient`:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.webflux.test.autoconfigure.WebFluxTest;
import org.springframework.http.MediaType;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.reactive.server.WebTestClient;
import static org.mockito.BDDMockito.given;
@WebFluxTest(UserVehicleController.class)
class MyControllerTests {
	@Autowired
	private WebTestClient webClient;
	@MockitoBean
	private UserVehicleService userVehicleService;
	@Test
	void testExample() {
 given(this.userVehicleService.getVehicleDetails("sboot"))
 .willReturn(new VehicleDetails("Honda", "Civic"));
 this.webClient.get().uri("/sboot/vehicle").accept(MediaType.TEXT_PLAIN).exchange()
 .expectStatus().isOk()
 .expectBody(String.class).isEqualTo("Honda Civic");
	}
}
```
```
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.webflux.test.autoconfigure.WebFluxTest
import org.springframework.http.MediaType
import org.springframework.test.context.bean.override.mockito.MockitoBean
import org.springframework.test.web.reactive.server.WebTestClient
import org.springframework.test.web.reactive.server.expectBody
@WebFluxTest(UserVehicleController::class)
class MyControllerTests(@Autowired val webClient: WebTestClient) {
	@MockitoBean
	lateinit var userVehicleService: UserVehicleService
	@Test
	fun testExample() {
 given(userVehicleService.getVehicleDetails("sboot"))
 .willReturn(VehicleDetails("Honda", "Civic"))
 webClient.get().uri("/sboot/vehicle").accept(MediaType.TEXT_PLAIN).exchange()
 .expectStatus().isOk
 .expectBody<String>().isEqualTo("Honda Civic")
	}
}
```
| This setup is only supported by WebFlux applications as using `WebTestClient`in a mocked web application only works with WebFlux at the moment. |

| `@WebFluxTest`cannot detect routes registered through the functional web framework.
For testing`RouterFunction`beans in the context, consider importing your`RouterFunction`yourself by using`@Import`or by using`@SpringBootTest`. |

| `@WebFluxTest`cannot detect custom security configuration registered as a`@Bean`of type`SecurityWebFilterChain`.
To include that in your test, you will need to import the configuration that registers the bean by using`@Import`or by using`@SpringBootTest`. |

| Sometimes writing Spring WebFlux tests is not enough; Spring Boot can help you run full end-to-end tests with an actual server. |

## Auto-configured Spring GraphQL Tests

Spring GraphQL offers a dedicated testing support module; you’ll need to add it to your project:

```
<dependencies>
	<dependency>
 <groupId>org.springframework.graphql</groupId>
 <artifactId>spring-graphql-test</artifactId>
 <scope>test</scope>
	</dependency>
	<!-- Unless already present in the compile scope -->
	<dependency>
 <groupId>org.springframework.boot</groupId>
 <artifactId>spring-boot-starter-webflux</artifactId>
 <scope>test</scope>
	</dependency>
</dependencies>
```
```
dependencies {
	testImplementation("org.springframework.graphql:spring-graphql-test")
	// Unless already present in the implementation configuration
	testImplementation("org.springframework.boot:spring-boot-starter-webflux")
}
```
This testing module ships the GraphQlTester.
The tester is heavily used in test, so be sure to become familiar with using it.
There are `GraphQlTester` variants and Spring Boot will auto-configure them depending on the type of tests:

-
the `ExecutionGraphQlServiceTester`performs tests on the server side, without a client nor a transport
-
the `HttpGraphQlTester`performs tests with a client that connects to a server, with or without a live server

Spring Boot helps you to test your Spring GraphQL Controllers with the `@GraphQlTest` annotation from the `spring-boot-graphql-test` module.
`@GraphQlTest` auto-configures the Spring GraphQL infrastructure, without any transport nor server being involved.
This limits scanned beans to `@Controller`, `@ControllerAdvice`, `RuntimeWiringConfigurer`, `JacksonComponent`, `@JsonComponent` (deprecated), `Converter`, `GenericConverter`, `DataFetcherExceptionResolver`, `Instrumentation` and `GraphQlSourceBuilderCustomizer`.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@GraphQlTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.

| A list of the auto-configurations that are enabled by `@GraphQlTest`can be found in the appendix. |

Often, `@GraphQlTest` is limited to a set of controllers and used in combination with the `@MockitoBean` annotation to provide mock implementations for required collaborators.

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.docs.web.graphql.runtimewiring.GreetingController;
import org.springframework.boot.graphql.test.autoconfigure.GraphQlTest;
import org.springframework.graphql.test.tester.GraphQlTester;
@GraphQlTest(GreetingController.class)
class GreetingControllerTests {
	@Autowired
	private GraphQlTester graphQlTester;
	@Test
	void shouldGreetWithSpecificName() {
 this.graphQlTester.document("{ greeting(name: \"Alice\") } ")
 .execute()
 .path("greeting")
 .entity(String.class)
 .isEqualTo("Hello, Alice!");
	}
	@Test
	void shouldGreetWithDefaultName() {
 this.graphQlTester.document("{ greeting } ")
 .execute()
 .path("greeting")
 .entity(String.class)
 .isEqualTo("Hello, Spring!");
	}
}
```
```
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.docs.web.graphql.runtimewiring.GreetingController
import org.springframework.boot.graphql.test.autoconfigure.GraphQlTest
import org.springframework.graphql.test.tester.GraphQlTester
@GraphQlTest(GreetingController::class)
internal class GreetingControllerTests {
	@Autowired
	lateinit var graphQlTester: GraphQlTester
	@Test
	fun shouldGreetWithSpecificName() {
 graphQlTester.document("{ greeting(name: \"Alice\") } ").execute().path("greeting").entity(String::class.java)
 .isEqualTo("Hello, Alice!")
	}
	@Test
	fun shouldGreetWithDefaultName() {
 graphQlTester.document("{ greeting } ").execute().path("greeting").entity(String::class.java)
 .isEqualTo("Hello, Spring!")
	}
}
```
`@SpringBootTest` tests are full integration tests and involve the entire application.
A `HttpGraphQlTester` bean can be added by annotating your test class with `@AutoConfigureHttpGraphQlTester` from the `spring-boot-graphql-test` module:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.graphql.test.autoconfigure.tester.AutoConfigureHttpGraphQlTester;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment;
import org.springframework.graphql.test.tester.HttpGraphQlTester;
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@AutoConfigureHttpGraphQlTester
class GraphQlIntegrationTests {
	@Test
	void shouldGreetWithSpecificName(@Autowired HttpGraphQlTester graphQlTester) {
 HttpGraphQlTester authenticatedTester = graphQlTester.mutate()
 .webTestClient((client) -> client.defaultHeaders((headers) -> headers.setBasicAuth("admin", "ilovespring")))
 .build();
 authenticatedTester.document("{ greeting(name: \"Alice\") } ")
 .execute()
 .path("greeting")
 .entity(String.class)
 .isEqualTo("Hello, Alice!");
	}
}
```
```
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.graphql.test.autoconfigure.tester.AutoConfigureHttpGraphQlTester
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment
import org.springframework.graphql.test.tester.HttpGraphQlTester
import org.springframework.http.HttpHeaders
import org.springframework.test.web.reactive.server.WebTestClient
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@AutoConfigureHttpGraphQlTester
class GraphQlIntegrationTests {
	@Test
	fun shouldGreetWithSpecificName(@Autowired graphQlTester: HttpGraphQlTester) {
 val authenticatedTester = graphQlTester.mutate()
 .webTestClient { client: WebTestClient.Builder ->
 client.defaultHeaders { headers: HttpHeaders ->
 headers.setBasicAuth("admin", "ilovespring")
 }
 }.build()
 authenticatedTester.document("{ greeting(name: \"Alice\") } ").execute()
 .path("greeting").entity(String::class.java).isEqualTo("Hello, Alice!")
	}
}
```
The `HttpGraphQlTester` bean uses the relevant transport of the integration test.
When using a random or defined port, the tester is configured against the live server.
To bind the tester to `MockMvc`, make sure to annotate your test class with `@AutoConfigureMockMvc`.

## Auto-configured Data Cassandra Tests

You can use `@DataCassandraTest` from the `spring-boot-data-cassandra-test` module to test Data Cassandra applications.
By default, it configures a `CassandraTemplate`, scans for `@Table` classes, and configures Spring Data Cassandra repositories.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@DataCassandraTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.
(For more about using Cassandra with Spring Boot, see Cassandra.)

| A list of the auto-configuration settings that are enabled by `@DataCassandraTest`can be found in the appendix. |

The following example shows a typical setup for using Cassandra tests in Spring Boot:

-
Java
-
Kotlin

```
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.data.cassandra.test.autoconfigure.DataCassandraTest;
@DataCassandraTest
class MyDataCassandraTests {
	@Autowired
	private SomeRepository repository;
}
```
```
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.data.cassandra.test.autoconfigure.DataCassandraTest
@DataCassandraTest
class MyDataCassandraTests(@Autowired val repository: SomeRepository)
```
## Auto-configured Data Couchbase Tests

You can use `@DataCouchbaseTest` from the `spring-boot-data-couchbase-test` module to test Data Couchbase applications.
By default, it configures a `CouchbaseTemplate` or `ReactiveCouchbaseTemplate`, scans for `@Document` classes, and configures Spring Data Couchbase repositories.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@DataCouchbaseTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.
(For more about using Couchbase with Spring Boot, see Couchbase, earlier in this chapter.)

| A list of the auto-configuration settings that are enabled by `@DataCouchbaseTest`can be found in the appendix. |

The following example shows a typical setup for using Couchbase tests in Spring Boot:

-
Java
-
Kotlin

```
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.data.couchbase.test.autoconfigure.DataCouchbaseTest;
@DataCouchbaseTest
class MyDataCouchbaseTests {
	@Autowired
	private SomeRepository repository;
	// ...
}
```
```
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.data.couchbase.test.autoconfigure.DataCouchbaseTest
@DataCouchbaseTest
class MyDataCouchbaseTests(@Autowired val repository: SomeRepository) {
	// ...
}
```
## Auto-configured Data Elasticsearch Tests

You can use `@DataElasticsearchTest` from the `spring-boot-data-elasticsearch-test` module to test Data Elasticsearch applications.
By default, it configures an `ElasticsearchTemplate`, scans for `@Document` classes, and configures Spring Data Elasticsearch repositories.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@DataElasticsearchTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.
(For more about using Elasticsearch with Spring Boot, see Elasticsearch, earlier in this chapter.)

| A list of the auto-configuration settings that are enabled by `@DataElasticsearchTest`can be found in the appendix. |

The following example shows a typical setup for using Elasticsearch tests in Spring Boot:

-
Java
-
Kotlin

```
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.data.elasticsearch.test.autoconfigure.DataElasticsearchTest;
@DataElasticsearchTest
class MyDataElasticsearchTests {
	@Autowired
	private SomeRepository repository;
	// ...
}
```
```
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.data.elasticsearch.test.autoconfigure.DataElasticsearchTest
@DataElasticsearchTest
class MyDataElasticsearchTests(@Autowired val repository: SomeRepository) {
	// ...
}
```
## Auto-configured Data JPA Tests

You can use the `@DataJpaTest` annotation from the `spring-boot-data-jpa-test` module to test Data JPA applications.
By default, it scans for `@Entity` classes and configures Spring Data JPA repositories.
If an embedded database is available on the classpath, it configures one as well.
SQL queries are logged by default by setting the `spring.jpa.show-sql` property to `true`.
This can be disabled using the `showSql` attribute of the annotation.

Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@DataJpaTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.

| A list of the auto-configuration settings that are enabled by `@DataJpaTest`can be found in the appendix. |

By default, data JPA tests are transactional and roll back at the end of each test. See the relevant section in the Spring Framework Reference Documentation for more details. If that is not what you want, you can disable transaction management for a test or for the whole class as follows:

-
Java
-
Kotlin

```
import org.springframework.boot.data.jpa.test.autoconfigure.DataJpaTest;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;
@DataJpaTest
@Transactional(propagation = Propagation.NOT_SUPPORTED)
class MyNonTransactionalTests {
	// ...
}
```
```
import org.springframework.boot.data.jpa.test.autoconfigure.DataJpaTest
import org.springframework.transaction.annotation.Propagation
import org.springframework.transaction.annotation.Transactional
@DataJpaTest
@Transactional(propagation = Propagation.NOT_SUPPORTED)
class MyNonTransactionalTests {
	// ...
}
```
Data JPA tests may also inject a `TestEntityManager` bean, which provides an alternative to the standard JPA `EntityManager` that is specifically designed for tests.

| `TestEntityManager`can also be auto-configured to any of your Spring-based test class by adding`@AutoConfigureTestEntityManager`.
When doing so, make sure that your test is running in a transaction, for instance by adding`@Transactional`on your test class or method. |

A `JdbcTemplate` is also available if you need that.
The following example shows the `@DataJpaTest` annotation in use:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.data.jpa.test.autoconfigure.DataJpaTest;
import org.springframework.boot.jpa.test.autoconfigure.TestEntityManager;
import static org.assertj.core.api.Assertions.assertThat;
@DataJpaTest
class MyRepositoryTests {
	@Autowired
	private TestEntityManager entityManager;
	@Autowired
	private UserRepository repository;
	@Test
	void testExample() {
 this.entityManager.persist(new User("sboot", "1234"));
 User user = this.repository.findByUsername("sboot");
 assertThat(user.getUsername()).isEqualTo("sboot");
 assertThat(user.getEmployeeNumber()).isEqualTo("1234");
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.data.jpa.test.autoconfigure.DataJpaTest
import org.springframework.boot.jpa.test.autoconfigure.TestEntityManager
@DataJpaTest
class MyRepositoryTests(@Autowired val entityManager: TestEntityManager, @Autowired val repository: UserRepository) {
	@Test
	fun testExample() {
 entityManager.persist(User("sboot", "1234"))
 val user = repository.findByUsername("sboot")
 assertThat(user?.username).isEqualTo("sboot")
 assertThat(user?.employeeNumber).isEqualTo("1234")
	}
}
```
In-memory embedded databases generally work well for tests, since they are fast and do not require any installation.
If, however, you prefer to run tests against a real database you can use the `@AutoConfigureTestDatabase` annotation, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.data.jpa.test.autoconfigure.DataJpaTest;
import org.springframework.boot.jdbc.test.autoconfigure.AutoConfigureTestDatabase;
import org.springframework.boot.jdbc.test.autoconfigure.AutoConfigureTestDatabase.Replace;
@DataJpaTest
@AutoConfigureTestDatabase(replace = Replace.NONE)
class MyRepositoryTests {
	// ...
}
```
```
import org.springframework.boot.jdbc.test.autoconfigure.AutoConfigureTestDatabase
import org.springframework.boot.data.jpa.test.autoconfigure.DataJpaTest
@DataJpaTest
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
class MyRepositoryTests {
	// ...
}
```
## Auto-configured JDBC Tests

`@JdbcTest` from the `spring-boot-jdbc-test` module is similar to `@DataJdbcTest` but is for tests that only require a `DataSource` and do not use Spring Data JDBC.
By default, it configures an in-memory embedded database and a `JdbcTemplate`.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@JdbcTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.

| A list of the auto-configurations that are enabled by `@JdbcTest`can be found in the appendix. |

By default, JDBC tests are transactional and roll back at the end of each test. See the relevant section in the Spring Framework Reference Documentation for more details. If that is not what you want, you can disable transaction management for a test or for the whole class, as follows:

-
Java
-
Kotlin

```
import org.springframework.boot.jdbc.test.autoconfigure.JdbcTest;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;
@JdbcTest
@Transactional(propagation = Propagation.NOT_SUPPORTED)
class MyTransactionalTests {
}
```
```
import org.springframework.boot.jdbc.test.autoconfigure.JdbcTest
import org.springframework.transaction.annotation.Propagation
import org.springframework.transaction.annotation.Transactional
@JdbcTest
@Transactional(propagation = Propagation.NOT_SUPPORTED)
class MyTransactionalTests
```
If you prefer your test to run against a real database, you can use the `@AutoConfigureTestDatabase` annotation in the same way as for `@DataJpaTest`.
(See Auto-configured Data JPA Tests.)

## Auto-configured Data JDBC Tests

`@DataJdbcTest` from the `spring-boot-data-jdbc-test` module is similar to `@JdbcTest` but is for tests that use Spring Data JDBC repositories.
By default, it configures an in-memory embedded database, a `JdbcTemplate`, and Spring Data JDBC repositories.
Only `AbstractJdbcConfiguration` subclasses are scanned when the `@DataJdbcTest` annotation is used, regular `@Component` and `@ConfigurationProperties` beans are not scanned.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.

| A list of the auto-configurations that are enabled by `@DataJdbcTest`can be found in the appendix. |

By default, Data JDBC tests are transactional and roll back at the end of each test. See the relevant section in the Spring Framework Reference Documentation for more details. If that is not what you want, you can disable transaction management for a test or for the whole test class as shown in the JDBC example.

If you prefer your test to run against a real database, you can use the `@AutoConfigureTestDatabase` annotation in the same way as for `@DataJpaTest`.
(See Auto-configured Data JPA Tests.)

## Auto-configured Data R2DBC Tests

`@DataR2dbcTest` from the `spring-boot-data-r2dbc-test` module is similar to `@DataJdbcTest` but is for tests that use Spring Data R2DBC repositories.
By default, it configures an in-memory embedded database, an `R2dbcEntityTemplate`, and Spring Data R2DBC repositories.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@DataR2dbcTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.

| A list of the auto-configurations that are enabled by `@DataR2dbcTest`can be found in the appendix. |

By default, Data R2DBC tests are not transactional.

If you prefer your test to run against a real database, you can use the `@AutoConfigureTestDatabase` annotation in the same way as for `@DataJpaTest`.
(See Auto-configured Data JPA Tests.)

## Auto-configured jOOQ Tests

You can use `@JooqTest` from `spring-boot-jooq-test` in a similar fashion as `@JdbcTest` but for jOOQ-related tests.
As jOOQ relies heavily on a Java-based schema that corresponds with the database schema, the existing `DataSource` is used.
If you want to replace it with an in-memory database, you can use `@AutoConfigureTestDatabase` to override those settings.
(For more about using jOOQ with Spring Boot, see Using jOOQ.)
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@JooqTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.

| A list of the auto-configurations that are enabled by `@JooqTest`can be found in the appendix. |

`@JooqTest` configures a `DSLContext`.
The following example shows the `@JooqTest` annotation in use:

-
Java
-
Kotlin

```
import org.jooq.DSLContext;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.jooq.test.autoconfigure.JooqTest;
@JooqTest
class MyJooqTests {
	@Autowired
	private DSLContext dslContext;
	// ...
}
```
```
import org.jooq.DSLContext
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.jooq.test.autoconfigure.JooqTest
@JooqTest
class MyJooqTests(@Autowired val dslContext: DSLContext) {
	// ...
}
```
JOOQ tests are transactional and roll back at the end of each test by default. If that is not what you want, you can disable transaction management for a test or for the whole test class as shown in the JDBC example.

## Auto-configured Data MongoDB Tests

You can use `@DataMongoTest` from the `spring-boot-data-mongodb-test` module to test MongoDB applications.
By default, it configures a `MongoTemplate`, scans for `@Document` classes, and configures Spring Data MongoDB repositories.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@DataMongoTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.
(For more about using MongoDB with Spring Boot, see MongoDB.)

| A list of the auto-configuration settings that are enabled by `@DataMongoTest`can be found in the appendix. |

The following class shows the `@DataMongoTest` annotation in use:

-
Java
-
Kotlin

```
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.data.mongodb.test.autoconfigure.DataMongoTest;
import org.springframework.data.mongodb.core.MongoTemplate;
@DataMongoTest
class MyDataMongoDbTests {
	@Autowired
	private MongoTemplate mongoTemplate;
	// ...
}
```
```
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.data.mongodb.test.autoconfigure.DataMongoTest
import org.springframework.data.mongodb.core.MongoTemplate
@DataMongoTest
class MyDataMongoDbTests(@Autowired val mongoTemplate: MongoTemplate) {
	// ...
}
```
## Auto-configured Data Neo4j Tests

You can use `@DataNeo4jTest` from the `spring-boot-data-neo4j-test` module to test Neo4j applications.
By default, it scans for `@Node` classes, and configures Spring Data Neo4j repositories.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@DataNeo4jTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.
(For more about using Neo4J with Spring Boot, see Neo4j.)

| A list of the auto-configuration settings that are enabled by `@DataNeo4jTest`can be found in the appendix. |

The following example shows a typical setup for using Neo4J tests in Spring Boot:

-
Java
-
Kotlin

```
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.data.neo4j.test.autoconfigure.DataNeo4jTest;
@DataNeo4jTest
class MyDataNeo4jTests {
	@Autowired
	private SomeRepository repository;
	// ...
}
```
```
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.data.neo4j.test.autoconfigure.DataNeo4jTest
@DataNeo4jTest
class MyDataNeo4jTests(@Autowired val repository: SomeRepository) {
	// ...
}
```
By default, Data Neo4j tests are transactional and roll back at the end of each test. See the relevant section in the Spring Framework Reference Documentation for more details. If that is not what you want, you can disable transaction management for a test or for the whole class, as follows:

-
Java
-
Kotlin

```
import org.springframework.boot.data.neo4j.test.autoconfigure.DataNeo4jTest;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;
@DataNeo4jTest
@Transactional(propagation = Propagation.NOT_SUPPORTED)
class MyDataNeo4jTests {
}
```
```
import org.springframework.boot.data.neo4j.test.autoconfigure.DataNeo4jTest
import org.springframework.transaction.annotation.Propagation
import org.springframework.transaction.annotation.Transactional
@DataNeo4jTest
@Transactional(propagation = Propagation.NOT_SUPPORTED)
class MyDataNeo4jTests
```
| Transactional tests are not supported with reactive access.
If you are using this style, you must configure `@DataNeo4jTest`tests as described above. |

## Auto-configured Data Redis Tests

You can use `@DataRedisTest` from the `spring-boot-data-redis-test` module to test Data Redis applications.
By default, it scans for `@RedisHash` classes and configures Spring Data Redis repositories.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@DataRedisTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.
(For more about using Redis with Spring Boot, see Redis.)

| A list of the auto-configuration settings that are enabled by `@DataRedisTest`can be found in the appendix. |

The following example shows the `@DataRedisTest` annotation in use:

-
Java
-
Kotlin

```
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.data.redis.test.autoconfigure.DataRedisTest;
@DataRedisTest
class MyDataRedisTests {
	@Autowired
	private SomeRepository repository;
	// ...
}
```
```
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.data.redis.test.autoconfigure.DataRedisTest
@DataRedisTest
class MyDataRedisTests(@Autowired val repository: SomeRepository) {
	// ...
}
```
## Auto-configured Data LDAP Tests

You can use `@DataLdapTest` to test Data LDAP applications.
By default, it configures an in-memory embedded LDAP (if available), configures an `LdapTemplate`, scans for `@Entry` classes, and configures Spring Data LDAP repositories.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@DataLdapTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.
(For more about using LDAP with Spring Boot, see LDAP.)

| A list of the auto-configuration settings that are enabled by `@DataLdapTest`can be found in the appendix. |

The following example shows the `@DataLdapTest` annotation in use:

-
Java
-
Kotlin

```
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.data.ldap.test.autoconfigure.DataLdapTest;
import org.springframework.ldap.core.LdapTemplate;
@DataLdapTest
class MyDataLdapTests {
	@Autowired
	private LdapTemplate ldapTemplate;
	// ...
}
```
```
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.data.ldap.test.autoconfigure.DataLdapTest
import org.springframework.ldap.core.LdapTemplate
@DataLdapTest
class MyDataLdapTests(@Autowired val ldapTemplate: LdapTemplate) {
	// ...
}
```
In-memory embedded LDAP generally works well for tests, since it is fast and does not require any developer installation. If, however, you prefer to run tests against a real LDAP server, you should exclude the embedded LDAP auto-configuration, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.data.ldap.test.autoconfigure.DataLdapTest;
import org.springframework.boot.ldap.autoconfigure.embedded.EmbeddedLdapAutoConfiguration;
@DataLdapTest(excludeAutoConfiguration = EmbeddedLdapAutoConfiguration.class)
class MyDataLdapTests {
	// ...
}
```
```
import org.springframework.boot.ldap.autoconfigure.embedded.EmbeddedLdapAutoConfiguration
import org.springframework.boot.data.ldap.test.autoconfigure.DataLdapTest
@DataLdapTest(excludeAutoConfiguration = [EmbeddedLdapAutoConfiguration::class])
class MyDataLdapTests {
	// ...
}
```
## Auto-configured REST Clients

You can use the `@RestClientTest` annotation from the `spring-boot-restclient-test` module to test REST clients.
By default, it auto-configures Jackson, GSON, and Jsonb support, configures a `RestTemplateBuilder` and a `RestClient.Builder`, and adds support for `MockRestServiceServer`.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@RestClientTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.

| A list of the auto-configuration settings that are enabled by `@RestClientTest`can be found in the appendix. |

The specific beans that you want to test should be specified by using the `value` or `components` attribute of `@RestClientTest`.

When using a `RestTemplateBuilder` in the beans under test and `RestTemplateBuilder.rootUri(String rootUri)` has been called when building the `RestTemplate`, then the root URI should be omitted from the `MockRestServiceServer` expectations as shown in the following example:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.restclient.test.autoconfigure.RestClientTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.client.MockRestServiceServer;
import static org.assertj.core.api.Assertions.assertThat;
import static org.springframework.test.web.client.match.MockRestRequestMatchers.requestTo;
import static org.springframework.test.web.client.response.MockRestResponseCreators.withSuccess;
@RestClientTest(org.springframework.boot.docs.testing.springbootapplications.autoconfiguredrestclient.RemoteVehicleDetailsService.class)
class MyRestTemplateServiceTests {
	@Autowired
	private RemoteVehicleDetailsService service;
	@Autowired
	private MockRestServiceServer server;
	@Test
	void getVehicleDetailsWhenResultIsSuccessShouldReturnDetails() {
 this.server.expect(requestTo("/greet/details")).andRespond(withSuccess("hello", MediaType.TEXT_PLAIN));
 String greeting = this.service.callRestService();
 assertThat(greeting).isEqualTo("hello");
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.restclient.test.autoconfigure.RestClientTest
import org.springframework.http.MediaType
import org.springframework.test.web.client.MockRestServiceServer
import org.springframework.test.web.client.match.MockRestRequestMatchers
import org.springframework.test.web.client.response.MockRestResponseCreators
@RestClientTest(RemoteVehicleDetailsService::class)
class MyRestTemplateServiceTests(
	@Autowired val service: RemoteVehicleDetailsService,
	@Autowired val server: MockRestServiceServer) {
	@Test
	fun getVehicleDetailsWhenResultIsSuccessShouldReturnDetails() {
 server.expect(MockRestRequestMatchers.requestTo("/greet/details"))
 .andRespond(MockRestResponseCreators.withSuccess("hello", MediaType.TEXT_PLAIN))
 val greeting = service.callRestService()
 assertThat(greeting).isEqualTo("hello")
	}
}
```
When using a `RestClient.Builder` in the beans under test, or when using a `RestTemplateBuilder` without calling `rootUri(String rootURI)`, the full URI must be used in the `MockRestServiceServer` expectations as shown in the following example:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.restclient.test.autoconfigure.RestClientTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.client.MockRestServiceServer;
import static org.assertj.core.api.Assertions.assertThat;
import static org.springframework.test.web.client.match.MockRestRequestMatchers.requestTo;
import static org.springframework.test.web.client.response.MockRestResponseCreators.withSuccess;
@RestClientTest(RemoteVehicleDetailsService.class)
class MyRestClientServiceTests {
	@Autowired
	private RemoteVehicleDetailsService service;
	@Autowired
	private MockRestServiceServer server;
	@Test
	void getVehicleDetailsWhenResultIsSuccessShouldReturnDetails() {
 this.server.expect(requestTo("https://example.com/greet/details"))
 .andRespond(withSuccess("hello", MediaType.TEXT_PLAIN));
 String greeting = this.service.callRestService();
 assertThat(greeting).isEqualTo("hello");
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.restclient.test.autoconfigure.RestClientTest
import org.springframework.http.MediaType
import org.springframework.test.web.client.MockRestServiceServer
import org.springframework.test.web.client.match.MockRestRequestMatchers
import org.springframework.test.web.client.response.MockRestResponseCreators
@RestClientTest(RemoteVehicleDetailsService::class)
class MyRestClientServiceTests(
	@Autowired val service: RemoteVehicleDetailsService,
	@Autowired val server: MockRestServiceServer) {
	@Test
	fun getVehicleDetailsWhenResultIsSuccessShouldReturnDetails() {
 server.expect(MockRestRequestMatchers.requestTo("https://example.com/greet/details"))
 .andRespond(MockRestResponseCreators.withSuccess("hello", MediaType.TEXT_PLAIN))
 val greeting = service.callRestService()
 assertThat(greeting).isEqualTo("hello")
	}
}
```
## Auto-configured Web Clients

You can use the `@WebClientTest` annotation from the `spring-boot-webclient-test` module to test code that uses `WebClient`.
By default, it auto-configures Jackson, GSON, and Jsonb support, and configures a `WebClient.Builder`.
Regular `@Component` and `@ConfigurationProperties` beans are not scanned when the `@WebClientTest` annotation is used.
`@EnableConfigurationProperties` can be used to include `@ConfigurationProperties` beans.

| A list of the auto-configuration settings that are enabled by `@WebClientTest`can be found in the appendix. |

The specific beans that you want to test should be specified by using the `value` or `components` attribute of `@WebClientTest`.

## Auto-configured Spring REST Docs Tests

You can use the `@AutoConfigureRestDocs` annotation from the `spring-boot-restdocs- module to use Spring REST Docs in your tests with Mock MVC or WebTestClient.
It removes the need for the JUnit extension in Spring REST Docs.

`@AutoConfigureRestDocs` can be used to override the default output directory (`target/generated-snippets` if you are using Maven or `build/generated-snippets` if you are using Gradle).
It can also be used to configure the host, scheme, and port that appears in any documented URIs.

### Auto-configured Spring REST Docs Tests With Mock MVC

`@AutoConfigureRestDocs` customizes the `MockMvc` bean to use Spring REST Docs when testing servlet-based web applications.
You can inject it by using `@Autowired` and use it in your tests as you normally would when using Mock MVC and Spring REST Docs, as shown in the following example:

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.restdocs.test.autoconfigure.AutoConfigureRestDocs;
import org.springframework.boot.webmvc.test.autoconfigure.WebMvcTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.assertj.MockMvcTester;
import static org.assertj.core.api.Assertions.assertThat;
import static org.springframework.restdocs.mockmvc.MockMvcRestDocumentation.document;
@WebMvcTest(UserController.class)
@AutoConfigureRestDocs
class MyUserDocumentationTests {
	@Autowired
	private MockMvcTester mvc;
	@Test
	void listUsers() {
 assertThat(this.mvc.get().uri("/users").accept(MediaType.TEXT_PLAIN)).hasStatusOk()
 .apply(document("list-users"));
	}
}
```
If you prefer to use the AssertJ integration, `MockMvcTester` is available as well, as shown in the following example:

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.restdocs.test.autoconfigure.AutoConfigureRestDocs;
import org.springframework.boot.webmvc.test.autoconfigure.WebMvcTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.assertj.MockMvcTester;
import static org.assertj.core.api.Assertions.assertThat;
import static org.springframework.restdocs.mockmvc.MockMvcRestDocumentation.document;
@WebMvcTest(UserController.class)
@AutoConfigureRestDocs
class MyUserDocumentationTests {
	@Autowired
	private MockMvcTester mvc;
	@Test
	void listUsers() {
 assertThat(this.mvc.get().uri("/users").accept(MediaType.TEXT_PLAIN)).hasStatusOk()
 .apply(document("list-users"));
	}
}
```
Both reuses the same `MockMvc` instance behind the scenes so any configuration to it applies to both.

If you require more control over Spring REST Docs configuration than offered by the attributes of `@AutoConfigureRestDocs`, you can use a `RestDocsMockMvcConfigurationCustomizer` bean, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.restdocs.test.autoconfigure.RestDocsMockMvcConfigurationCustomizer;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.restdocs.mockmvc.MockMvcRestDocumentationConfigurer;
import org.springframework.restdocs.templates.TemplateFormats;
@TestConfiguration(proxyBeanMethods = false)
public class MyRestDocsConfiguration implements RestDocsMockMvcConfigurationCustomizer {
	@Override
	public void customize(MockMvcRestDocumentationConfigurer configurer) {
 configurer.snippets().withTemplateFormat(TemplateFormats.markdown());
	}
}
```
```
import org.springframework.boot.restdocs.test.autoconfigure.RestDocsMockMvcConfigurationCustomizer
import org.springframework.boot.test.context.TestConfiguration
import org.springframework.restdocs.mockmvc.MockMvcRestDocumentationConfigurer
import org.springframework.restdocs.templates.TemplateFormats
@TestConfiguration(proxyBeanMethods = false)
class MyRestDocsConfiguration : RestDocsMockMvcConfigurationCustomizer {
	override fun customize(configurer: MockMvcRestDocumentationConfigurer) {
 configurer.snippets().withTemplateFormat(TemplateFormats.markdown())
	}
}
```
If you want to make use of Spring REST Docs support for a parameterized output directory, you can create a `RestDocumentationResultHandler` bean.
The auto-configuration calls `alwaysDo` with this result handler, thereby causing each `MockMvc` call to automatically generate the default snippets.
The following example shows a `RestDocumentationResultHandler` being defined:

-
Java
-
Kotlin

```
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.restdocs.mockmvc.MockMvcRestDocumentation;
import org.springframework.restdocs.mockmvc.RestDocumentationResultHandler;
@TestConfiguration(proxyBeanMethods = false)
public class MyResultHandlerConfiguration {
	@Bean
	public RestDocumentationResultHandler restDocumentation() {
 return MockMvcRestDocumentation.document("{method-name}");
	}
}
```
```
import org.springframework.boot.test.context.TestConfiguration
import org.springframework.context.annotation.Bean
import org.springframework.restdocs.mockmvc.MockMvcRestDocumentation
import org.springframework.restdocs.mockmvc.RestDocumentationResultHandler
@TestConfiguration(proxyBeanMethods = false)
class MyResultHandlerConfiguration {
	@Bean
	fun restDocumentation(): RestDocumentationResultHandler {
 return MockMvcRestDocumentation.document("{method-name}")
	}
}
```
### Auto-configured Spring REST Docs Tests With WebTestClient

`@AutoConfigureRestDocs` can also be used with `WebTestClient` when testing reactive web applications.
You can inject it by using `@Autowired` and use it in your tests as you normally would when using `@WebFluxTest` and Spring REST Docs, as shown in the following example:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.restdocs.test.autoconfigure.AutoConfigureRestDocs;
import org.springframework.boot.webflux.test.autoconfigure.WebFluxTest;
import org.springframework.test.web.reactive.server.WebTestClient;
import static org.springframework.restdocs.webtestclient.WebTestClientRestDocumentation.document;
@WebFluxTest
@AutoConfigureRestDocs
class MyUsersDocumentationTests {
	@Autowired
	private WebTestClient webTestClient;
	@Test
	void listUsers() {
 this.webTestClient
 .get().uri("/")
 .exchange()
 .expectStatus()
 .isOk()
 .expectBody()
 .consumeWith(document("list-users"));
	}
}
```
```
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.restdocs.test.autoconfigure.AutoConfigureRestDocs
import org.springframework.boot.webflux.test.autoconfigure.WebFluxTest
import org.springframework.restdocs.webtestclient.WebTestClientRestDocumentation
import org.springframework.test.web.reactive.server.WebTestClient
@WebFluxTest
@AutoConfigureRestDocs
class MyUsersDocumentationTests(@Autowired val webTestClient: WebTestClient) {
	@Test
	fun listUsers() {
 webTestClient
 .get().uri("/")
 .exchange()
 .expectStatus()
 .isOk
 .expectBody()
 .consumeWith(WebTestClientRestDocumentation.document("list-users"))
	}
}
```
If you require more control over Spring REST Docs configuration than offered by the attributes of `@AutoConfigureRestDocs`, you can use a `RestDocsWebTestClientConfigurationCustomizer` bean, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.restdocs.test.autoconfigure.RestDocsWebTestClientConfigurationCustomizer;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.restdocs.webtestclient.WebTestClientRestDocumentationConfigurer;
@TestConfiguration(proxyBeanMethods = false)
public class MyRestDocsConfiguration implements RestDocsWebTestClientConfigurationCustomizer {
	@Override
	public void customize(WebTestClientRestDocumentationConfigurer configurer) {
 configurer.snippets().withEncoding("UTF-8");
	}
}
```
```
import org.springframework.boot.restdocs.test.autoconfigure.RestDocsWebTestClientConfigurationCustomizer
import org.springframework.boot.test.context.TestConfiguration
import org.springframework.restdocs.webtestclient.WebTestClientRestDocumentationConfigurer
@TestConfiguration(proxyBeanMethods = false)
class MyRestDocsConfiguration : RestDocsWebTestClientConfigurationCustomizer {
	override fun customize(configurer: WebTestClientRestDocumentationConfigurer) {
 configurer.snippets().withEncoding("UTF-8")
	}
}
```
If you want to make use of Spring REST Docs support for a parameterized output directory, you can use a `WebTestClientBuilderCustomizer` to configure a consumer for every entity exchange result.
The following example shows such a `WebTestClientBuilderCustomizer` being defined:

-
Java
-
Kotlin

```
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.boot.webtestclient.autoconfigure.WebTestClientBuilderCustomizer;
import org.springframework.context.annotation.Bean;
import static org.springframework.restdocs.webtestclient.WebTestClientRestDocumentation.document;
@TestConfiguration(proxyBeanMethods = false)
public class MyWebTestClientBuilderCustomizerConfiguration {
	@Bean
	public WebTestClientBuilderCustomizer restDocumentation() {
 return (builder) -> builder.entityExchangeResultConsumer(document("{method-name}"));
	}
}
```
```
import org.springframework.boot.test.context.TestConfiguration
import org.springframework.boot.webtestclient.autoconfigure.WebTestClientBuilderCustomizer
import org.springframework.context.annotation.Bean
import org.springframework.restdocs.webtestclient.WebTestClientRestDocumentation
import org.springframework.test.web.reactive.server.WebTestClient
@TestConfiguration(proxyBeanMethods = false)
class MyWebTestClientBuilderCustomizerConfiguration {
	@Bean
	fun restDocumentation(): WebTestClientBuilderCustomizer {
 return WebTestClientBuilderCustomizer { builder: WebTestClient.Builder ->
 builder.entityExchangeResultConsumer(
 WebTestClientRestDocumentation.document("{method-name}")
 )
 }
	}
}
```
## Auto-configured Spring Web Services Tests

### Auto-configured Spring Web Services Client Tests

You can use `@WebServiceClientTest` from the `spring-boot-webservices-test` module to test applications that call web services using the Spring Web Services project.
By default, it configures a `MockWebServiceServer` bean and automatically customizes your `WebServiceTemplateBuilder`.
(For more about using Web Services with Spring Boot, see Web Services.)

| A list of the auto-configuration settings that are enabled by `@WebServiceClientTest`can be found in the appendix. |

The following example shows the `@WebServiceClientTest` annotation in use:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.webservices.test.autoconfigure.client.WebServiceClientTest;
import org.springframework.ws.test.client.MockWebServiceServer;
import org.springframework.xml.transform.StringSource;
import static org.assertj.core.api.Assertions.assertThat;
import static org.springframework.ws.test.client.RequestMatchers.payload;
import static org.springframework.ws.test.client.ResponseCreators.withPayload;
@WebServiceClientTest(SomeWebService.class)
class MyWebServiceClientTests {
	@Autowired
	private MockWebServiceServer server;
	@Autowired
	private SomeWebService someWebService;
	@Test
	void mockServerCall() {
 this.server
 .expect(payload(new StringSource("<request/>")))
 .andRespond(withPayload(new StringSource("<response><status>200</status></response>")));
 assertThat(this.someWebService.test())
 .extracting(Response::getStatus)
 .isEqualTo(200);
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.webservices.test.autoconfigure.client.WebServiceClientTest
import org.springframework.ws.test.client.MockWebServiceServer
import org.springframework.ws.test.client.RequestMatchers
import org.springframework.ws.test.client.ResponseCreators
import org.springframework.xml.transform.StringSource
@WebServiceClientTest(SomeWebService::class)
class MyWebServiceClientTests(
 @Autowired val server: MockWebServiceServer, @Autowired val someWebService: SomeWebService) {
	@Test
	fun mockServerCall() {
 server
 .expect(RequestMatchers.payload(StringSource("<request/>")))
 .andRespond(ResponseCreators.withPayload(StringSource("<response><status>200</status></response>")))
 assertThat(this.someWebService.test()).extracting(Response::status).isEqualTo(200)
	}
}
```
### Auto-configured Spring Web Services Server Tests

You can use `@WebServiceServerTest` from the `spring-boot-webservices-test` module to test applications that implement web services using the Spring Web Services project.
By default, it configures a `MockWebServiceClient` bean that can be used to call your web service endpoints.
(For more about using Web Services with Spring Boot, see Web Services.)

| A list of the auto-configuration settings that are enabled by `@WebServiceServerTest`can be found in the appendix. |

The following example shows the `@WebServiceServerTest` annotation in use:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.webservices.test.autoconfigure.server.WebServiceServerTest;
import org.springframework.ws.test.server.MockWebServiceClient;
import org.springframework.ws.test.server.RequestCreators;
import org.springframework.ws.test.server.ResponseMatchers;
import org.springframework.xml.transform.StringSource;
@WebServiceServerTest(ExampleEndpoint.class)
class MyWebServiceServerTests {
	@Autowired
	private MockWebServiceClient client;
	@Test
	void mockServerCall() {
 this.client
 .sendRequest(RequestCreators.withPayload(new StringSource("<ExampleRequest/>")))
 .andExpect(ResponseMatchers.payload(new StringSource("<ExampleResponse>42</ExampleResponse>")));
	}
}
```
```
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.webservices.test.autoconfigure.server.WebServiceServerTest
import org.springframework.ws.test.server.MockWebServiceClient
import org.springframework.ws.test.server.RequestCreators
import org.springframework.ws.test.server.ResponseMatchers
import org.springframework.xml.transform.StringSource
@WebServiceServerTest(ExampleEndpoint::class)
class MyWebServiceServerTests(@Autowired val client: MockWebServiceClient) {
	@Test
	fun mockServerCall() {
 client
 .sendRequest(RequestCreators.withPayload(StringSource("<ExampleRequest/>")))
 .andExpect(ResponseMatchers.payload(StringSource("<ExampleResponse>42</ExampleResponse>")))
	}
}
```
## Additional Auto-configuration and Slicing

Each slice provides one or more `@AutoConfigure…` annotations that namely defines the auto-configurations that should be included as part of a slice.
Additional auto-configurations can be added on a test-by-test basis by creating a custom `@AutoConfigure…` annotation or by adding `@ImportAutoConfiguration` to the test as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.autoconfigure.ImportAutoConfiguration;
import org.springframework.boot.integration.autoconfigure.IntegrationAutoConfiguration;
import org.springframework.boot.jdbc.test.autoconfigure.JdbcTest;
@JdbcTest
@ImportAutoConfiguration(IntegrationAutoConfiguration.class)
class MyJdbcTests {
}
```
```
import org.springframework.boot.autoconfigure.ImportAutoConfiguration
import org.springframework.boot.integration.autoconfigure.IntegrationAutoConfiguration
import org.springframework.boot.jdbc.test.autoconfigure.JdbcTest
@JdbcTest
@ImportAutoConfiguration(IntegrationAutoConfiguration::class)
class MyJdbcTests
```
| Make sure to not use the regular `@Import`annotation to import auto-configurations as they are handled in a specific way by Spring Boot. |

Alternatively, additional auto-configurations can be added for any use of a slice annotation by registering them in a file stored in `META-INF/spring` as shown in the following example:

`com.example.IntegrationAutoConfiguration`In this example, the `com.example.IntegrationAutoConfiguration` is enabled on every test annotated with `@JdbcTest`.

| You can use comments with `#`in this file. |

| A slice or `@AutoConfigure…`annotation can be customized this way as long as it is meta-annotated with`@ImportAutoConfiguration`. |

## User Configuration and Slicing

If you structure your code in a sensible way, your `@SpringBootApplication` class is used by default as the configuration of your tests.

It then becomes important not to litter the application’s main class with configuration settings that are specific to a particular area of its functionality.

Assume that you are using Spring Data MongoDB, you rely on the auto-configuration for it, and you have enabled auditing.
You could define your `@SpringBootApplication` as follows:

-
Java
-
Kotlin

```
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.data.mongodb.config.EnableMongoAuditing;
@SpringBootApplication
@EnableMongoAuditing
public class MyApplication {
	// ...
}
```
```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.data.mongodb.config.EnableMongoAuditing
@SpringBootApplication
@EnableMongoAuditing
class MyApplication {
	// ...
}
```
Because this class is the source configuration for the test, any slice test actually tries to enable Mongo auditing, which is definitely not what you want to do.
A recommended approach is to move that area-specific configuration to a separate `@Configuration` class at the same level as your application, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.context.annotation.Configuration;
import org.springframework.data.mongodb.config.EnableMongoAuditing;
@Configuration(proxyBeanMethods = false)
@EnableMongoAuditing
public class MyMongoConfiguration {
	// ...
}
```
```
import org.springframework.context.annotation.Configuration
import org.springframework.data.mongodb.config.EnableMongoAuditing
@Configuration(proxyBeanMethods = false)
@EnableMongoAuditing
class MyMongoConfiguration {
	// ...
}
```
| Depending on the complexity of your application, you may either have a single `@Configuration`class for your customizations or one class per domain area.
The latter approach lets you enable it in one of your tests, if necessary, with the`@Import`annotation.
See this how-to section for more details on when you might want to enable specific`@Configuration`classes for slice tests. |

Test slices exclude `@Configuration` classes from scanning.
For example, for a `@WebMvcTest`, the following configuration will not include the given `WebMvcConfigurer` bean in the application context loaded by the test slice:

-
Java
-
Kotlin

```
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
@Configuration(proxyBeanMethods = false)
public class MyWebConfiguration {
	@Bean
	public WebMvcConfigurer testConfigurer() {
 return new WebMvcConfigurer() {
 // ...
 };
	}
}
```
```
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer
@Configuration(proxyBeanMethods = false)
class MyWebConfiguration {
	@Bean
	fun testConfigurer(): WebMvcConfigurer {
 return object : WebMvcConfigurer {
 // ...
 }
	}
}
```
The configuration below will, however, cause the custom `WebMvcConfigurer` to be loaded by the test slice.

-
Java
-
Kotlin

```
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
@Component
public class MyWebMvcConfigurer implements WebMvcConfigurer {
	// ...
}
```
```
import org.springframework.stereotype.Component
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer
@Component
class MyWebMvcConfigurer : WebMvcConfigurer {
	// ...
}
```
Another source of confusion is classpath scanning. Assume that, while you structured your code in a sensible way, you need to scan an additional package. Your application may resemble the following code:

-
Java
-
Kotlin

```
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.ComponentScan;
@SpringBootApplication
@ComponentScan({ "com.example.app", "com.example.another" })
public class MyApplication {
	// ...
}
```
```
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.context.annotation.ComponentScan
@SpringBootApplication
@ComponentScan("com.example.app", "com.example.another")
class MyApplication {
	// ...
}
```
Doing so effectively overrides the default component scan directive with the side effect of scanning those two packages regardless of the slice that you chose.
For instance, a `@DataJpaTest` seems to suddenly scan components and user configurations of your application.
Again, moving the custom directive to a separate class is a good way to fix this issue.

| If this is not an option for you, you can create a `@SpringBootConfiguration`somewhere in the hierarchy of your test so that it is used instead.
Alternatively, you can specify a source for your test, which disables the behavior of finding a default one. |

## Using Spock to Test Spring Boot Applications

Spock 2.4 or later can be used to test a Spring Boot application.
To do so, add a dependency on a `-groovy-5.0` version of Spock’s `spock-spring` module to your application’s build.
`spock-spring` integrates Spring’s test framework into Spock.
See the documentation for Spock’s Spring module for further details.

# Testcontainers

The Testcontainers library provides a way to manage services running inside Docker containers. It integrates with JUnit, allowing you to write a test class that can start up a container before any of the tests run. Testcontainers is especially useful for writing integration tests that talk to a real backend service such as MySQL, MongoDB, Cassandra and others.

In following sections we will describe some of the methods you can use to integrate Testcontainers with your tests.

## Using Spring Beans

The containers provided by Testcontainers can be managed by Spring Boot as beans.

To declare a container as a bean, add a `@Bean` method to your test configuration:

-
Java
-
Kotlin

```
import org.testcontainers.mongodb.MongoDBContainer;
import org.testcontainers.utility.DockerImageName;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;
@TestConfiguration(proxyBeanMethods = false)
class MyTestConfiguration {
	@Bean
	MongoDBContainer mongoDbContainer() {
 return new MongoDBContainer(DockerImageName.parse("mongo:5.0"));
	}
}
```
```
import org.springframework.boot.test.context.TestConfiguration
import org.springframework.context.annotation.Bean
import org.testcontainers.mongodb.MongoDBContainer
import org.testcontainers.utility.DockerImageName
@TestConfiguration(proxyBeanMethods = false)
class MyTestConfiguration {
	@Bean
	fun mongoDbContainer(): MongoDBContainer {
 return MongoDBContainer(DockerImageName.parse("mongo:5.0"))
	}
}
```
You can then inject and use the container by importing the configuration class in the test class:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.testcontainers.mongodb.MongoDBContainer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Import;
@SpringBootTest
@Import(MyTestConfiguration.class)
class MyIntegrationTests {
	@Autowired
	private MongoDBContainer mongo;
	@Test
	void myTest() {
 ...
	}
}
```
```
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.context.annotation.Import
import org.testcontainers.mongodb.MongoDBContainer
@SpringBootTest
@Import(MyTestConfiguration::class)
class MyIntegrationTests {
	@Autowired
	private val mongo: MongoDBContainer? = null
	@Test
	fun myTest() {
 ...
	}
}
```
| This method of managing containers is often used in combination with service connection annotations. |

## Using the JUnit Extension

Testcontainers provides a JUnit extension which can be used to manage containers in your tests.
The extension is activated by applying the `@Testcontainers` annotation from Testcontainers to your test class.

You can then use the `@Container` annotation on static container fields.

The `@Testcontainers` annotation can be used on vanilla JUnit tests, or in combination with `@SpringBootTest`:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.neo4j.Neo4jContainer;
import org.springframework.boot.test.context.SpringBootTest;
@Testcontainers
@SpringBootTest
class MyIntegrationTests {
	@Container
	static Neo4jContainer neo4j = new Neo4jContainer("neo4j:5");
	@Test
	void myTest() {
 ...
	}
}
```
```
import org.junit.jupiter.api.Test;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.neo4j.Neo4jContainer;
import org.springframework.boot.test.context.SpringBootTest;
@Testcontainers
@SpringBootTest
class MyIntegrationTests {
	@Test
	fun myTest() {
 ...
	}
	companion object {
 @Container
 @JvmStatic
 val neo4j = Neo4jContainer("neo4j:5");
	}
}
```
The example above will start up a Neo4j container before any of the tests are run. The lifecycle of the container instance is managed by Testcontainers, as described in their official documentation.

When using the JUnit extension, container instances are stopped after the test class has run (for static fields) or after each test method (for non-static fields).
This can cause issues when used with Spring Boot tests, as Spring’s TestContext Framework may cache the `ApplicationContext` beyond that point and reuse it for another test class or method with the same configuration.
If the cached application context contains beans that depend on a container that has already been stopped, later tests or bean destruction callbacks may fail.
For this reason, you should prefer managing containers as Spring beans or importing container declarations when the application context should remain usable for as long as it is cached.

| In most cases, you will additionally need to configure the application to connect to the service running in the container. |

## Importing Container Configuration Interfaces

A common pattern with Testcontainers is to declare the container instances as static fields in an interface.

For example, the following interface declares two containers, one named `mongo` of type `MongoDBContainer` and another named `neo4j` of type `Neo4jContainer`:

-
Java
-
Kotlin

```
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.mongodb.MongoDBContainer;
import org.testcontainers.neo4j.Neo4jContainer;
interface MyContainers {
	@Container
	MongoDBContainer mongoContainer = new MongoDBContainer("mongo:5.0");
	@Container
	Neo4jContainer neo4jContainer = new Neo4jContainer("neo4j:5");
}
```
```
import org.testcontainers.junit.jupiter.Container
import org.testcontainers.mongodb.MongoDBContainer
import org.testcontainers.neo4j.Neo4jContainer
interface MyContainers {
	companion object {
 @Container
 val mongoContainer: MongoDBContainer = MongoDBContainer("mongo:5.0")
 @Container
 val neo4jContainer: Neo4jContainer = Neo4jContainer("neo4j:5")
	}
}
```
When you have containers declared in this way, you can reuse their configuration in multiple tests by having the test classes implement the interface.

It’s also possible to use the same interface configuration in your Spring Boot tests.
To do so, add `@ImportTestcontainers` to your test configuration class:

-
Java
-
Kotlin

```
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.boot.testcontainers.context.ImportTestcontainers;
@TestConfiguration(proxyBeanMethods = false)
@ImportTestcontainers(MyContainers.class)
class MyTestConfiguration {
}
```
```
import org.springframework.boot.test.context.TestConfiguration
import org.springframework.boot.testcontainers.context.ImportTestcontainers
@TestConfiguration(proxyBeanMethods = false)
@ImportTestcontainers(MyContainers::class)
class MyTestConfiguration {
}
```
## Lifecycle of Managed Containers

If you have used the annotations and extensions provided by Testcontainers, then the lifecycle of container instances is managed entirely by Testcontainers. Please refer to the official Testcontainers documentation for the information.

When the containers are managed by Spring as beans, then their lifecycle is managed by Spring:

-
Container beans are created and started before all other beans.
-
Container beans are stopped after the destruction of all other beans.

This process ensures that any beans, which rely on functionality provided by the containers, can use those functionalities. It also ensures that they are cleaned up whilst the container is still available.

| When your application beans rely on functionality of containers, prefer configuring the containers as Spring beans to ensure the correct lifecycle behavior. |

| Having containers managed by Testcontainers instead of as Spring beans provides no guarantee of the order in which beans and containers will shutdown.
It can happen that containers are shutdown before the beans relying on container functionality are cleaned up.
This can lead to exceptions being thrown by client beans, for example, due to loss of connection.
When the application context should remain usable for as long as it is cached, prefer managing containers as Spring beans or importing container declarations with `@ImportTestcontainers`. |

Container beans are created and started once per application context managed by Spring’s TestContext Framework. For details about how the TestContext Framework manages the underlying application contexts and beans therein, please refer to the Spring Framework documentation.

Container beans are stopped as part of the TestContext Framework’s standard application context shutdown process. When the application context gets shutdown, the containers are shutdown as well. This usually happens after all tests using that specific cached application context have finished executing. It may also happen earlier, depending on the caching behavior configured in the TestContext Framework.

| A single test container instance can, and often is, retained across execution of tests from multiple test classes. |

## Service Connections

A service connection is a connection to any remote service. Spring Boot’s auto-configuration can consume the details of a service connection and use them to establish a connection to a remote service. When doing so, the connection details take precedence over any connection-related configuration properties.

When using Testcontainers, connection details can be automatically created for a service running in a container by annotating the container field in the test class.

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.neo4j.Neo4jContainer;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
@Testcontainers
@SpringBootTest
class MyIntegrationTests {
	@Container
	@ServiceConnection
	static Neo4jContainer neo4j = new Neo4jContainer("neo4j:5");
	@Test
	void myTest() {
 ...
	}
}
```
```
import org.junit.jupiter.api.Test;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.neo4j.Neo4jContainer;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
@Testcontainers
@SpringBootTest
class MyIntegrationTests {
	@Test
	fun myTest() {
 ...
	}
	companion object {
 @Container
 @ServiceConnection
 @JvmStatic
 val neo4j = Neo4jContainer("neo4j:5");
	}
}
```
Thanks to `@ServiceConnection`, the above configuration allows Neo4j-related beans in the application to communicate with Neo4j running inside the Testcontainers-managed Docker container.
This is done by automatically defining a `Neo4jConnectionDetails` bean which is then used by the Neo4j auto-configuration, overriding any connection-related configuration properties.

| You’ll need to add the `spring-boot-testcontainers`module as a test dependency in order to use service connections with Testcontainers. |

Service connection annotations are processed by `ContainerConnectionDetailsFactory` classes registered with `spring.factories`.
A `ContainerConnectionDetailsFactory` can create a `ConnectionDetails` bean based on a specific `Container` subclass, or the Docker image name.

The following service connection factories are provided in the `spring-boot-testcontainers` jar:

| Connection Details | Matched on |
|---|---|
| Containers named "symptoma/activemq" or | |
| Containers of type | |
| Containers of type | |
| Containers of type | |
| Containers of type | |
| Containers of type | |
| Containers of type | |
| Containers of type | |
| Containers named "osixia/openldap" or of type | |
| Containers of type | |
| Containers of type | |
| Containers of type | |
| Containers named "otel/opentelemetry-collector-contrib" or of type | |
| Containers named "otel/opentelemetry-collector-contrib" or of type | |
| Containers named "otel/opentelemetry-collector-contrib" or of type | |
| Containers of type | |
| Containers of type
 | |
| Containers of type | |
| Containers of type | |
| Containers of type | |
| Containers named "openzipkin/zipkin" |

| By default, with the exception of If you want to create only a subset of the applicable types, you can use the To create a |

By default `Container.getDockerImageName().getRepository()` is used to obtain the name used to find connection details.
The repository portion of the Docker image name ignores any registry and the version.
This works as long as Spring Boot is able to get the instance of the `Container`, which is the case when using a `static` field like in the example above.

If you’re using a `@Bean` method, Spring Boot won’t call the bean method to get the Docker image name, because this would cause eager initialization issues.
Instead, the return type of the bean method is used to find out which connection detail should be used.
This works as long as you’re using typed containers such as `Neo4jContainer` or `RabbitMQContainer`.
This stops working if you’re using `GenericContainer`, for example with Redis as shown in the following example:

-
Java
-
Kotlin

```
import org.testcontainers.containers.GenericContainer;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.context.annotation.Bean;
@TestConfiguration(proxyBeanMethods = false)
public class MyRedisConfiguration {
	@Bean
	@ServiceConnection(name = "redis")
	public GenericContainer<?> redisContainer() {
 return new GenericContainer<>("redis:7");
	}
}
```
```
import org.springframework.boot.test.context.TestConfiguration
import org.springframework.boot.testcontainers.service.connection.ServiceConnection
import org.springframework.context.annotation.Bean
import org.testcontainers.containers.GenericContainer
@TestConfiguration(proxyBeanMethods = false)
class MyRedisConfiguration {
	@Bean
	@ServiceConnection(name = "redis")
	fun redisContainer(): GenericContainer<*> {
 return GenericContainer("redis:7")
	}
}
```
Spring Boot can’t tell from `GenericContainer` which container image is used, so the `name` attribute from `@ServiceConnection` must be used to provide that hint.

You can also use the `name` attribute of `@ServiceConnection` to override which connection detail will be used, for example when using custom images.
If you are using the Docker image `registry.mycompany.com/mirror/myredis`, you’d use `@ServiceConnection(name="redis")` to ensure `DataRedisConnectionDetails` are created.

### SSL with Service Connections

You can use the `@Ssl`, `@JksKeyStore`, `@JksTrustStore`, `@PemKeyStore` and `@PemTrustStore` annotations on a supported container to enable SSL support for that service connection.
Please note that you still have to enable SSL on the service which is running inside the Testcontainer yourself, the annotations only configure SSL on the client side in your application.

```
import com.redis.testcontainers.RedisContainer;
import org.junit.jupiter.api.Test;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.testcontainers.service.connection.PemKeyStore;
import org.springframework.boot.testcontainers.service.connection.PemTrustStore;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.data.redis.core.RedisOperations;
@Testcontainers
@SpringBootTest
class MyRedisWithSslIntegrationTests {
	@Container
	@ServiceConnection
	@PemKeyStore(certificate = "classpath:client.crt", privateKey = "classpath:client.key")
	@PemTrustStore("classpath:ca.crt")
	static RedisContainer redis = new SecureRedisContainer("redis:latest");
	@Autowired
	private RedisOperations<Object, Object> operations;
	@Test
	void testRedis() {
 // ...
	}
}
```
The above code uses the `@PemKeyStore` annotation to load the client certificate and key into the keystore and the and `@PemTrustStore` annotation to load the CA certificate into the truststore.
This will authenticate the client against the server, and the CA certificate in the truststore makes sure that the server certificate is valid and trusted.

The `SecureRedisContainer` in this example is a custom subclass of `RedisContainer` which copies certificates to the correct places and invokes `redis-server` with commandline parameters enabling SSL.

The SSL annotations are supported for the following service connections:

-
Cassandra
-
Couchbase
-
Elasticsearch
-
Kafka
-
MongoDB
-
RabbitMQ
-
RabbitMQ Streams
-
Redis

The `ElasticsearchContainer` additionally supports automatic detection of server side SSL.
To use this feature, annotate the container with `@Ssl`, as seen in the following example, and Spring Boot takes care of the client side SSL configuration for you:

```
import org.junit.jupiter.api.Test;
import org.testcontainers.elasticsearch.ElasticsearchContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.data.elasticsearch.test.autoconfigure.DataElasticsearchTest;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.boot.testcontainers.service.connection.Ssl;
import org.springframework.data.elasticsearch.client.elc.ElasticsearchTemplate;
@Testcontainers
@DataElasticsearchTest
class MyElasticsearchWithSslIntegrationTests {
	@Ssl
	@Container
	@ServiceConnection
	static ElasticsearchContainer elasticsearch = new ElasticsearchContainer(
 "docker.elastic.co/elasticsearch/elasticsearch:8.17.2");
	@Autowired
	private ElasticsearchTemplate elasticsearchTemplate;
	@Test
	void testElasticsearch() {
 // ...
	}
}
```
## Dynamic Properties

A slightly more verbose but also more flexible alternative to service connections is `@DynamicPropertySource`.
A static `@DynamicPropertySource` method allows adding dynamic property values to the Spring Environment.

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.neo4j.Neo4jContainer;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
@Testcontainers
@SpringBootTest
class MyIntegrationTests {
	@Container
	static Neo4jContainer neo4j = new Neo4jContainer("neo4j:5");
	@Test
	void myTest() {
 // ...
	}
	@DynamicPropertySource
	static void neo4jProperties(DynamicPropertyRegistry registry) {
 registry.add("spring.neo4j.uri", neo4j::getBoltUrl);
	}
}
```
```
import org.junit.jupiter.api.Test
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.test.context.DynamicPropertyRegistry
import org.springframework.test.context.DynamicPropertySource
import org.testcontainers.junit.jupiter.Container
import org.testcontainers.junit.jupiter.Testcontainers
import org.testcontainers.neo4j.Neo4jContainer
@Testcontainers
@SpringBootTest
class MyIntegrationTests {
	@Test
	fun myTest() {
 ...
	}
	companion object {
 @Container
 @JvmStatic
 val neo4j = Neo4jContainer("neo4j:5");
 @DynamicPropertySource
 @JvmStatic
 fun neo4jProperties(registry: DynamicPropertyRegistry) {
 registry.add("spring.neo4j.uri") { neo4j.boltUrl }
 }
	}
}
```
The above configuration allows Neo4j-related beans in the application to communicate with Neo4j running inside the Testcontainers-managed Docker container.

# Test Utilities

A few test utility classes that are generally useful when testing your application are packaged as part of `spring-boot`.

## ConfigDataApplicationContextInitializer

`ConfigDataApplicationContextInitializer` is an `ApplicationContextInitializer` that you can apply to your tests to load Spring Boot `application.properties` files.
You can use it when you do not need the full set of features provided by `@SpringBootTest`, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.test.context.ConfigDataApplicationContextInitializer;
import org.springframework.test.context.ContextConfiguration;
@ContextConfiguration(classes = Config.class, initializers = ConfigDataApplicationContextInitializer.class)
class MyConfigFileTests {
	// ...
}
```
```
import org.springframework.boot.test.context.ConfigDataApplicationContextInitializer
import org.springframework.test.context.ContextConfiguration
@ContextConfiguration(classes = [Config::class], initializers = [ConfigDataApplicationContextInitializer::class])
class MyConfigFileTests {
	// ...
}
```
| Using `ConfigDataApplicationContextInitializer`alone does not provide support for`@Value("${…}")`injection.
Its only job is to ensure that`application.properties`files are loaded into Spring’s`Environment`.
For`@Value`support, you need to either additionally configure a`PropertySourcesPlaceholderConfigurer`or use`@SpringBootTest`, which auto-configures one for you. |

## TestPropertyValues

`TestPropertyValues` lets you quickly add properties to a `ConfigurableEnvironment` or `ConfigurableApplicationContext`.
You can call it with `key=value` strings, as follows:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.util.TestPropertyValues;
import org.springframework.mock.env.MockEnvironment;
import static org.assertj.core.api.Assertions.assertThat;
class MyEnvironmentTests {
	@Test
	void testPropertySources() {
 MockEnvironment environment = new MockEnvironment();
 TestPropertyValues.of("org=Spring", "name=Boot").applyTo(environment);
 assertThat(environment.getProperty("name")).isEqualTo("Boot");
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.boot.test.util.TestPropertyValues
import org.springframework.mock.env.MockEnvironment
class MyEnvironmentTests {
	@Test
	fun testPropertySources() {
 val environment = MockEnvironment()
 TestPropertyValues.of("org=Spring", "name=Boot").applyTo(environment)
 assertThat(environment.getProperty("name")).isEqualTo("Boot")
	}
}
```
## OutputCaptureExtension

`OutputCaptureExtension` is a JUnit `Extension` that you can use to capture `System.out` and `System.err` output.
To use it, add `@ExtendWith(OutputCaptureExtension.class)` and inject `CapturedOutput` as an argument to your test class constructor or test method as follows:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.boot.test.system.CapturedOutput;
import org.springframework.boot.test.system.OutputCaptureExtension;
import static org.assertj.core.api.Assertions.assertThat;
@ExtendWith(OutputCaptureExtension.class)
class MyOutputCaptureTests {
	@Test
	void testName(CapturedOutput output) {
 System.out.println("Hello World!");
 assertThat(output).contains("World");
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.ExtendWith
import org.springframework.boot.test.system.CapturedOutput
import org.springframework.boot.test.system.OutputCaptureExtension
@ExtendWith(OutputCaptureExtension::class)
class MyOutputCaptureTests {
	@Test
	fun testName(output: CapturedOutput?) {
 println("Hello World!")
 assertThat(output).contains("World")
	}
}
```
## TestRestTemplate

`TestRestTemplate` is a convenience alternative to Spring’s `RestTemplate` that is useful in integration tests.
It’s provided by the `spring-boot-resttestclient` module.
A dependency on `spring-boot-restclient` is also required.
Take care when adding this dependency as it will enable auto-configuration for `RestClient.Builder`.
If your main code uses `RestClient.Builder`, declare the `spring-boot-restclient` dependency so that it is on your application’s main classpath and not only on its test classpath.

You can get a vanilla template or one that sends Basic HTTP authentication (with a username and password).
In either case, the template is fault tolerant.
This means that it behaves in a test-friendly way by not throwing exceptions on 4xx and 5xx errors.
Instead, such errors can be detected through the returned `ResponseEntity` and its status code.

If you need fluent API for assertions, consider using `RestTestClient` that works with mock environments and end-to-end tests.

If you are using Spring WebFlux, consider the `WebTestClient` that provides a similar API and works with mock environments, WebFlux integration tests, and end-to-end tests.

It is recommended, but not mandatory, to use the Apache HTTP Client (version 5.1 or better).
If you have that on your classpath, the `TestRestTemplate` responds by configuring the client appropriately.

`TestRestTemplate` can be instantiated directly in your integration tests, as shown in the following example:

-
Java
-
Kotlin

```
import org.junit.jupiter.api.Test;
import org.springframework.boot.resttestclient.TestRestTemplate;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import static org.assertj.core.api.Assertions.assertThat;
class MyTests {
	private final TestRestTemplate template = new TestRestTemplate();
	@Test
	void testRequest() {
 ResponseEntity<String> response = this.template.getForEntity("https://myhost.example.com/example",
 String.class);
 assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);
 // Other assertions to verify the response
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.boot.resttestclient.TestRestTemplate
import org.springframework.http.HttpStatus
class MyTests {
	private val template = TestRestTemplate()
	@Test
	fun testRequest() {
 val response = template.getForEntity("https://myhost.example.com/example", String::class.java)
 assertThat(response.statusCode).isEqualTo(HttpStatus.OK)
 // Other assertions to verify the response
	}
}
```
Alternatively, if you use the `@SpringBootTest` annotation with `WebEnvironment.RANDOM_PORT` or `WebEnvironment.DEFINED_PORT`, you can inject a fully configured `TestRestTemplate` by annotating the test class with `@AutoConfigureTestRestTemplate`.
If necessary, additional customizations can be applied through the `RestTemplateBuilder` bean.

Any URLs that do not specify a host and port automatically connect to the embedded server, as shown in the following example:

-
Java
-
Kotlin

```
import java.time.Duration;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.restclient.RestTemplateBuilder;
import org.springframework.boot.resttestclient.TestRestTemplate;
import org.springframework.boot.resttestclient.autoconfigure.AutoConfigureTestRestTemplate;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import static org.assertj.core.api.Assertions.assertThat;
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@AutoConfigureTestRestTemplate
class MySpringBootTests {
	@Autowired
	private TestRestTemplate template;
	@Test
	void testRequest() {
 ResponseEntity<String> response = this.template.getForEntity("/example", String.class);
 assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);
 // Other assertions to verify the response
	}
	@TestConfiguration(proxyBeanMethods = false)
	static class RestTemplateBuilderConfiguration {
 @Bean
 RestTemplateBuilder restTemplateBuilder() {
 return new RestTemplateBuilder().connectTimeout(Duration.ofSeconds(1)).readTimeout(Duration.ofSeconds(1));
 }
	}
}
```
```
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment
import org.springframework.boot.test.context.TestConfiguration
import org.springframework.boot.restclient.RestTemplateBuilder
import org.springframework.boot.resttestclient.TestRestTemplate
import org.springframework.boot.resttestclient.autoconfigure.AutoConfigureTestRestTemplate
import org.springframework.context.annotation.Bean
import org.springframework.http.HttpStatus
import java.time.Duration
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@AutoConfigureTestRestTemplate
class MySpringBootTests(@Autowired val template: TestRestTemplate) {
	@Test
	fun testRequest() {
 val response = template.getForEntity("/example", String::class.java)
 assertThat(response.statusCode).isEqualTo(HttpStatus.OK)
 // Other assertions to verify the response
	}
	@TestConfiguration(proxyBeanMethods = false)
	internal class RestTemplateBuilderConfiguration {
 @Bean
 fun restTemplateBuilder(): RestTemplateBuilder {
 return RestTemplateBuilder().connectTimeout(Duration.ofSeconds(1))
 .readTimeout(Duration.ofSeconds(1))
 }
	}
}
```

# Packaging Spring Boot Applications

Spring Boot supports several technologies for optimizing applications for deployment, including GraalVM native images, AOT cache, and Checkpoint and Restore.

Spring Boot applications can be packaged in Docker containers using techniques described in Container Images.

# Efficient Deployments

## Unpacking the Executable jar

You can run your application using the executable jar, but loading the classes from nested jars has a small startup cost. Depending on the size of the jar, running the application from an exploded structure is faster and recommended in production. Certain PaaS implementations may also choose to extract archives before they run. For example, Cloud Foundry operates this way.

Spring Boot supports extracting your application to a directory using different layouts. The default layout is the most efficient, and it is AOT cache (and CDS) friendly.

In this layout, the libraries are extracted to a `lib/` folder, and the application jar
contains the application classes and a manifest which references the libraries in the `lib/` folder.

To unpack the executable jar, run this command:

`$ java -Djarmode=tools -jar my-app.jar extract`And then in production, you can run the extracted jar:

`$ java -jar my-app/my-app.jar`After startup, you should not expect any differences in execution time between running an executable jar and running an extracted jar.

| Run `java -Djarmode=tools -jar my-app.jar help extract`to see all possible options. |

# AOT Cache

AOT cache is a JVM feature that can help reduce the startup time and memory footprint of Java applications.

If you are not yet using Java 25 or above, you should read the sections about CDS. CDS is the predecessor of AOT cache, but works similarly.

| Spring Boot supports both CDS and AOT cache, however, we recommend using the AOT cache whenever possible. |

## AOT Cache

| Spring Boot supports the AOT cache for Java 25 and above. If you’re using an earlier version of Java, you have to use CDS instead. |

To use the AOT cache feature, you should first perform a training run on your application in extracted form:

```
$ java -Djarmode=tools -jar my-app.jar extract --destination application
$ cd application
$ java -XX:AOTCacheOutput=app.aot -Dspring.context.exit=onRefresh -jar my-app.jar
```
This creates an `app.aot` cache file that can be reused as long as the application is not updated and the same Java version is used.

To use the cache file, you need to add an extra parameter when starting the application:

`$ java -XX:AOTCache=app.aot -jar my-app.jar`| You have to use the cache file with the extracted form of the application, otherwise it has no effect. |

## CDS

| If you’re using Java 25 or above, please use AOT cache instead of CDS. |

To use CDS, you should first perform a training run on your application in extracted form:

```
$ java -Djarmode=tools -jar my-app.jar extract --destination application
$ cd application
$ java -XX:ArchiveClassesAtExit=application.jsa -Dspring.context.exit=onRefresh -jar my-app.jar
```
This creates an `application.jsa` archive file that can be reused as long as the application is not updated.

To use the archive file, you need to add an extra parameter when starting the application:

`$ java -XX:SharedArchiveFile=application.jsa -jar my-app.jar`| You have to use the cache file with the extracted form of the application, otherwise it has no effect. |

| For more details about CDS, refer to the Class Data Sharing documentation of the JDK. |

# Ahead-of-Time Processing With the JVM

It’s beneficial for the startup time to run your application using the AOT generated initialization code. First, you need to ensure that the jar you are building includes AOT generated code.

| AOT cache and Spring’s AOT can be combined to further improve startup time. |

For Maven, this means that you should build with `-Pnative` to activate the `native` profile:

`$ mvn -Pnative package`For Gradle, you need to ensure that your build includes the `org.springframework.boot.aot` plugin.

When the JAR has been built, run it with `spring.aot.enabled` system property set to `true`. For example:

```
$ java -Dspring.aot.enabled=true -jar myapplication.jar
........ Starting AOT-processed MyApplication ...
```
Beware that using the ahead-of-time processing has drawbacks. It implies the following restrictions:

-
The classpath is fixed and fully defined at build time
-
The beans defined in your application cannot change at runtime, meaning: -
The Spring `@Profile`annotation and profile-specific configuration have limitations.
-
Properties that change if a bean is created are not supported (for example, `@ConditionalOnProperty`and`.enabled`properties).

-

To learn more about ahead-of-time processing, please see the Understanding Spring Ahead-of-Time Processing section.

# GraalVM Native Images

GraalVM Native Images are standalone executables that can be generated by processing compiled Java applications ahead-of-time. Native Images generally have a smaller memory footprint and start faster than their JVM counterparts.

# Introducing GraalVM Native Images

GraalVM Native Images provide a new way to deploy and run Java applications. Compared to the Java Virtual Machine, native images can run with a smaller memory footprint and with much faster startup times.

They are well suited to applications that are deployed using container images and are especially interesting when combined with "Function as a service" (FaaS) platforms.

Unlike traditional applications written for the JVM, GraalVM Native Image applications require ahead-of-time processing in order to create an executable. This ahead-of-time processing involves statically analyzing your application code from its main entry point.

A GraalVM Native Image is a complete, platform-specific executable. You do not need to ship a Java Virtual Machine in order to run a native image.

| If you just want to get started and experiment with GraalVM you can jump to the Developing Your First GraalVM Native Application section and return to this section later. |

## Key Differences with JVM Deployments

The fact that GraalVM Native Images are produced ahead-of-time means that there are some key differences between native and JVM based applications. The main differences are:

-
Static analysis of your application is performed at build-time from the `main`entry point.
-
Code that cannot be reached when the native image is created will be removed and won’t be part of the executable.
-
GraalVM is not directly aware of dynamic elements of your code and must be told about reflection, resources, serialization, and dynamic proxies.
-
The application classpath is fixed at build time and cannot change.
-
There is no lazy class loading, everything shipped in the executables will be loaded in memory on startup.
-
There are some limitations around some aspects of Java applications that are not fully supported.

On top of those differences, Spring uses a process called Spring Ahead-of-Time processing, which imposes further limitations. Please make sure to read at least the beginning of the next section to learn about those.

| The Native Image Compatibility Guide section of the GraalVM reference documentation provides more details about GraalVM limitations. |

## Understanding Spring Ahead-of-Time Processing

Typical Spring Boot applications are quite dynamic and configuration is performed at runtime. In fact, the concept of Spring Boot auto-configuration depends heavily on reacting to the state of the runtime in order to configure things correctly.

Although it would be possible to tell GraalVM about these dynamic aspects of the application, doing so would undo most of the benefit of static analysis. So instead, when using Spring Boot to create native images, a closed-world is assumed and the dynamic aspects of the application are restricted.

A closed-world assumption implies, besides the limitations created by GraalVM itself, the following restrictions:

-
The beans defined in your application cannot change at runtime, meaning: -
The Spring `@Profile`annotation and profile-specific configuration have limitations.
-
Properties that change if a bean is created are not supported (for example, `@ConditionalOnProperty`and`.enabled`properties).

-

When these restrictions are in place, it becomes possible for Spring to perform ahead-of-time processing during build-time and generate additional assets that GraalVM can use. A Spring AOT processed application will typically generate:

-
Java source code
-
Bytecode (for dynamic proxies, etc.)
-
GraalVM JSON hint files in `META-INF/native-image/{groupId}/{artifactId}/`:-
Resource hints ( `resource-config.json`)
-
Reflection hints ( `reflect-config.json`)
-
Serialization hints ( `serialization-config.json`)
-
Java Proxy Hints ( `proxy-config.json`)
-
JNI Hints ( `jni-config.json`)

-

If the generated hints are not sufficient, you can also provide your own.

### Source Code Generation

Spring applications are composed of Spring Beans. Internally, Spring Framework uses two distinct concepts to manage beans. There are bean instances, which are the actual instances that have been created and can be injected into other beans. There are also bean definitions which are used to define attributes of a bean and how its instance should be created.

If we take a typical `@Configuration` class:

```
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyConfiguration {
	@Bean
	public MyBean myBean() {
 return new MyBean();
	}
}
```
The bean definition is created by parsing the `@Configuration` class and finding the `@Bean` methods.
In the above example, we’re defining a `BeanDefinition` for a singleton bean named `myBean`.
We’re also creating a `BeanDefinition` for the `MyConfiguration` class itself.

When the `myBean` instance is required, Spring knows that it must invoke the `myBean()` method and use the result.
When running on the JVM, `@Configuration` class parsing happens when your application starts and `@Bean` methods are invoked using reflection.

When creating a native image, Spring operates in a different way.
Rather than parsing `@Configuration` classes and generating bean definitions at runtime, it does it at build-time.
Once the bean definitions have been discovered, they are processed and converted into source code that can be analyzed by the GraalVM compiler.

The Spring AOT process would convert the configuration class above to code like this:

```
import org.springframework.beans.factory.aot.BeanInstanceSupplier;
import org.springframework.beans.factory.config.BeanDefinition;
import org.springframework.beans.factory.support.RootBeanDefinition;
/**
 * Bean definitions for {@link MyConfiguration}.
 */
public class MyConfiguration__BeanDefinitions {
	/**
 * Get the bean definition for 'myConfiguration'.
 */
	public static BeanDefinition getMyConfigurationBeanDefinition() {
 Class<?> beanType = MyConfiguration.class;
 RootBeanDefinition beanDefinition = new RootBeanDefinition(beanType);
 beanDefinition.setInstanceSupplier(MyConfiguration::new);
 return beanDefinition;
	}
	/**
 * Get the bean instance supplier for 'myBean'.
 */
	private static BeanInstanceSupplier<MyBean> getMyBeanInstanceSupplier() {
 return BeanInstanceSupplier.<MyBean>forFactoryMethod(MyConfiguration.class, "myBean")
 .withGenerator((registeredBean) -> registeredBean.getBeanFactory().getBean(MyConfiguration.class).myBean());
	}
	/**
 * Get the bean definition for 'myBean'.
 */
	public static BeanDefinition getMyBeanBeanDefinition() {
 Class<?> beanType = MyBean.class;
 RootBeanDefinition beanDefinition = new RootBeanDefinition(beanType);
 beanDefinition.setInstanceSupplier(getMyBeanInstanceSupplier());
 return beanDefinition;
	}
}
```
| The exact code generated may differ depending on the nature of your bean definitions. |

You can see above that the generated code creates equivalent bean definitions to the `@Configuration` class, but in a direct way that can be understood by GraalVM.

There is a bean definition for the `myConfiguration` bean, and one for `myBean`.
When a `myBean` instance is required, a `BeanInstanceSupplier` is called.
This supplier will invoke the `myBean()` method on the `myConfiguration` bean.

| During Spring AOT processing, your application is started up to the point that bean definitions are available. Bean instances are not created during the AOT processing phase. |

Spring AOT will generate code like this for all your bean definitions.
It will also generate code when bean post-processing is required (for example, to call `@Autowired` methods).
An `ApplicationContextInitializer` will also be generated which will be used by Spring Boot to initialize the `ApplicationContext` when an AOT processed application is actually run.

| Although AOT generated source code can be verbose, it is quite readable and can be helpful when debugging an application.
Generated source files can be found in `target/spring-aot/main/sources`when using Maven and`build/generated/aotSources`with Gradle. |

### Hint File Generation

In addition to generating source files, the Spring AOT engine will also generate hint files that are used by GraalVM. Hint files contain JSON data that describes how GraalVM should deal with things that it can’t understand by directly inspecting the code.

For example, you might be using a Spring annotation on a private method. Spring will need to use reflection in order to invoke private methods, even on GraalVM. When such situations arise, Spring can write a reflection hint so that GraalVM knows that even though the private method isn’t called directly, it still needs to be available in the native image.

Hint files are generated under `META-INF/native-image` where they are automatically picked up by GraalVM.

| Generated hint files can be found in `target/spring-aot/main/resources`when using Maven and`build/generated/aotResources`with Gradle. |

### Proxy Class Generation

Spring sometimes needs to generate proxy classes to enhance the code you’ve written with additional features. To do this, it uses the cglib library which directly generates bytecode.

When an application is running on the JVM, proxy classes are generated dynamically as the application runs. When creating a native image, these proxies need to be created at build-time so that they can be included by GraalVM.

| Unlike source code generation, generated bytecode isn’t particularly helpful when debugging an application.
However, if you need to inspect the contents of the `.class`files using a tool such as`javap`you can find them in`target/spring-aot/main/classes`for Maven and`build/generated/aotClasses`for Gradle. |

# Advanced Native Images Topics

## Nested Configuration Properties

Reflection hints are automatically created for configuration properties by the Spring ahead-of-time engine.
Nested configuration properties which are not inner classes, however, **must** be annotated with `@NestedConfigurationProperty`, otherwise they won’t be detected and will not be bindable.

```
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.context.properties.NestedConfigurationProperty;
@ConfigurationProperties("my.properties")
public class MyProperties {
	private String name;
	@NestedConfigurationProperty
	private final Nested nested = new Nested();
	// getters / setters...
	public String getName() {
 return this.name;
	}
	public void setName(String name) {
 this.name = name;
	}
	public Nested getNested() {
 return this.nested;
	}
}
```
where `Nested` is:

-
Java
-
Kotlin

```
public class Nested {
	private int number;
	// getters / setters...
	public int getNumber() {
 return this.number;
	}
	public void setNumber(int number) {
 this.number = number;
	}
}
```
```
class Nested {
}
```
The example above produces configuration properties for `my.properties.name` and `my.properties.nested.number`.
Without the `@NestedConfigurationProperty` annotation on the `nested` field, the `my.properties.nested.number` property would not be bindable in a native image.
You can also annotate the getter method.

When using constructor binding, you have to annotate the field with `@NestedConfigurationProperty`:

```
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.context.properties.NestedConfigurationProperty;
@ConfigurationProperties("my.properties")
public class MyPropertiesCtor {
	private final String name;
	@NestedConfigurationProperty
	private final Nested nested;
	public MyPropertiesCtor(String name, Nested nested) {
 this.name = name;
 this.nested = nested;
	}
	// getters / setters...
	public String getName() {
 return this.name;
	}
	public Nested getNested() {
 return this.nested;
	}
}
```
When using records, you have to annotate the parameter with `@NestedConfigurationProperty`:

```
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.context.properties.NestedConfigurationProperty;
@ConfigurationProperties("my.properties")
public record MyPropertiesRecord(String name, @NestedConfigurationProperty Nested nested) {
}
```
When using Kotlin, you need to annotate the parameter of a data class with `@NestedConfigurationProperty`:

```
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.NestedConfigurationProperty
@ConfigurationProperties("my.properties")
data class MyPropertiesKotlin(
	val name: String,
	@NestedConfigurationProperty val nested: Nested
)
```
| Please use public getters and setters in all cases, otherwise the properties will not be bindable. |

## Converting a Spring Boot Executable Jar

It is possible to convert a Spring Boot executable jar into a native image as long as the jar contains the AOT generated assets. This can be useful for a number of reasons, including:

-
You can keep your regular JVM pipeline and turn the JVM application into a native image on your CI/CD platform.
-
As `native-image`does not support cross-compilation, you can keep an OS neutral deployment artifact which you convert later to different OS architectures.

You can convert a Spring Boot executable jar into a native image using Cloud Native Buildpacks, or using the `native-image` tool that is shipped with GraalVM.

| Your executable jar must include AOT generated assets such as generated classes and JSON hint files. |

### Using Buildpacks

Spring Boot applications usually use Cloud Native Buildpacks through the Maven (`mvn spring-boot:build-image`) or Gradle (`gradle bootBuildImage`) integrations.
You can, however, also use `pack` to turn an AOT processed Spring Boot executable jar into a native container image.

| You have to build your application with at least JDK 25, because Buildpacks use the same GraalVM native-image version as the Java version used for compilation. |

First, make sure that a Docker daemon is available (see Get Docker for more details). Configure it to allow non-root user if you are on Linux.

You also need to install `pack` by following the installation guide on buildpacks.io.

Assuming an AOT processed Spring Boot executable jar built as `myproject-0.0.1-SNAPSHOT.jar` is in the `target` directory, run:

```
$ pack build --builder paketobuildpacks/builder-noble-java-tiny \
 --path target/myproject-0.0.1-SNAPSHOT.jar \
 --env 'BP_NATIVE_IMAGE=true' \
 my-application:0.0.1-SNAPSHOT
```
| You do not need to have a local GraalVM installation to generate an image in this way. |

Once `pack` has finished, you can launch the application using `docker run`:

`$ docker run --rm -p 8080:8080 docker.io/library/myproject:0.0.1-SNAPSHOT`### Using GraalVM native-image

Another option to turn an AOT processed Spring Boot executable jar into a native executable is to use the GraalVM `native-image` tool.
For this to work, you’ll need a GraalVM distribution on your machine.
You can either download it manually on the Liberica Native Image Kit page or you can use a download manager like SDKMAN!.

Assuming an AOT processed Spring Boot executable jar built as `myproject-0.0.1-SNAPSHOT.jar` is in the `target` directory, run:

```
$ rm -rf target/native
$ mkdir -p target/native
$ cd target/native
$ jar -xvf ../myproject-0.0.1-SNAPSHOT.jar
$ native-image -H:Name=myproject @META-INF/native-image/argfile -cp .:BOOT-INF/classes:`find BOOT-INF/lib | tr '\n' ':'`
$ mv myproject ../
```
| These commands work on Linux or macOS machines, but you will need to adapt them for Windows. |

| The `@META-INF/native-image/argfile`might not be packaged in your jar.
It is only included when reachability metadata overrides are needed. |

| The `native-image``-cp`flag does not accept wildcards.
You need to ensure that all jars are listed (the command above uses`find`and`tr`to do this). |

## Using the Tracing Agent

The GraalVM native image tracing agent allows you to intercept reflection, resources or proxy usage on the JVM in order to generate the related hints. Spring should generate most of these hints automatically, but the tracing agent can be used to quickly identify the missing entries.

When using the agent to generate hints for a native image, there are a couple of approaches:

-
Launch the application directly and exercise it.
-
Run application tests to exercise the application.

The first option is interesting for identifying the missing hints when a library or a pattern is not recognized by Spring.

The second option sounds more appealing for a repeatable setup, but by default the generated hints will include anything required by the test infrastructure. Some of these will be unnecessary when the application runs for real. To address this problem the agent supports an access-filter file that will cause certain data to be excluded from the generated output.

### Launch the Application Directly

Use the following command to launch the application with the native image tracing agent attached:

```
$ java -Dspring.aot.enabled=true \
 -agentlib:native-image-agent=config-output-dir=/path/to/config-dir/ \
 -jar target/myproject-0.0.1-SNAPSHOT.jar
```
Now you can exercise the code paths you want to have hints for and then stop the application with `ctrl-c`.

On application shutdown the native image tracing agent will write the hint files to the given config output directory.
You can either manually inspect these files, or use them as input to the native image build process.
To use them as input, copy them into the `src/main/resources/META-INF/native-image/` directory.
The next time you build the native image, GraalVM will take these files into consideration.

There are more advanced options which can be set on the native image tracing agent, for example filtering the recorded hints by caller classes, etc. For further reading, please see the official documentation.

## Custom Hints

If you need to provide your own hints for reflection, resources, serialization, proxy usage and so on, you can use the `RuntimeHintsRegistrar` API.
Create a class that implements the `RuntimeHintsRegistrar` interface, and then make appropriate calls to the provided `RuntimeHints` instance:

```
import java.lang.reflect.Method;
import org.springframework.aot.hint.ExecutableMode;
import org.springframework.aot.hint.RuntimeHints;
import org.springframework.aot.hint.RuntimeHintsRegistrar;
import org.springframework.util.ReflectionUtils;
public class MyRuntimeHints implements RuntimeHintsRegistrar {
	@Override
	public void registerHints(RuntimeHints hints, ClassLoader classLoader) {
 // Register method for reflection
 Method method = ReflectionUtils.findMethod(MyClass.class, "sayHello", String.class);
 hints.reflection().registerMethod(method, ExecutableMode.INVOKE);
 // Register type for java serialization
 hints.reflection().registerJavaSerialization(MySerializableClass.class);
 // Register resources
 hints.resources().registerPattern("my-resource.txt");
 // Register proxy
 hints.proxies().registerJdkProxy(MyInterface.class);
	}
}
```
You can then use `@ImportRuntimeHints` on any `@Configuration` class (for example your `@SpringBootApplication` annotated application class) to activate those hints.

If you have classes which need binding (mostly needed when serializing or deserializing JSON), you can use `@RegisterReflectionForBinding` on any bean.
Most of the hints are automatically inferred, for example when accepting or returning data from a `@RestController` method.
But when you work with `WebClient`, `RestClient` or `RestTemplate` directly, you might need to use `@RegisterReflectionForBinding`.

### Testing Custom Hints

The `RuntimeHintsPredicates` API can be used to test your hints.
The API provides methods that build a `Predicate` that can be used to test a `RuntimeHints` instance.

If you’re using AssertJ, your test would look like this:

```
import org.junit.jupiter.api.Test;
import org.springframework.aot.hint.RuntimeHints;
import org.springframework.aot.hint.predicate.RuntimeHintsPredicates;
import org.springframework.boot.docs.packaging.nativeimage.advanced.customhints.MyRuntimeHints;
import static org.assertj.core.api.Assertions.assertThat;
class MyRuntimeHintsTests {
	@Test
	void shouldRegisterHints() {
 RuntimeHints hints = new RuntimeHints();
 new MyRuntimeHints().registerHints(hints, getClass().getClassLoader());
 assertThat(RuntimeHintsPredicates.resource().forResource("my-resource.txt")).accepts(hints);
	}
}
```
### Providing Hints Statically

If you prefer, custom hints can be provided statically in one or more GraalVM JSON hint files.
Such files should be placed in `src/main/resources/` within a `META-INF/native-image/*/*/` directory.
The hints generated during AOT processing are written to a directory named `META-INF/native-image/{groupId}/{artifactId}/`.
Place your static hint files in a directory that does not clash with this location, such as `META-INF/native-image/{groupId}/{artifactId}-additional-hints/`.

## Known Limitations

GraalVM native images are an evolving technology and not all libraries provide support. The GraalVM community is helping by providing reachability metadata for projects that don’t yet ship their own. Spring itself doesn’t contain hints for 3rd party libraries and instead relies on the reachability metadata project.

If you encounter problems when generating native images for Spring Boot applications, please check the Spring Boot with GraalVM page of the Spring Boot wiki. You can also contribute issues to the spring-aot-smoke-tests project on GitHub which is used to confirm that common application types are working as expected.

If you find a library which doesn’t work with GraalVM, please raise an issue on the reachability metadata project.

# Checkpoint and Restore With the JVM

Coordinated Restore at Checkpoint (CRaC) is an OpenJDK project that defines a new Java API to allow you to checkpoint and restore an application on the HotSpot JVM. It is based on CRIU, a project that implements checkpoint/restore functionality on Linux.

The principle is the following: you start your application almost as usual but with a CRaC enabled version of the JDK like BellSoft Liberica JDK with CRaC or Azul Zulu JDK with CRaC.
Then at some point, potentially after some workloads that will warm up your JVM by executing all common code paths, you trigger a checkpoint using an API call, a `jcmd` command, an HTTP endpoint, or a different mechanism.

A memory representation of the running JVM, including its warmness, is then serialized to disk, allowing a fast restoration at a later point, potentially on another machine with a similar operating system and CPU architecture. The restored process retains all the capabilities of the HotSpot JVM, including further JIT optimizations at runtime.

Based on the foundations provided by Spring Framework, Spring Boot provides support for checkpointing and restoring your application, and manages out-of-the-box the lifecycle of resources such as socket, files and thread pools on a limited scope. Additional lifecycle management is expected for other dependencies and potentially for the application code dealing with such resources.

You can find more details about the two modes supported ("on demand checkpoint/restore of a running application" and "automatic checkpoint/restore at startup"), how to enable checkpoint and restore support and some guidelines in the Spring Framework JVM Checkpoint Restore support documentation.

# Container Images

Spring Boot applications can be containerized using Dockerfiles, or by using Cloud Native Buildpacks to create optimized docker compatible container images that you can run anywhere.

Spring Boot applications can be containerized using Dockerfiles, or by using Cloud Native Buildpacks to create optimized docker compatible container images that you can run anywhere.

# Efficient Container Images

It is easily possible to package a Spring Boot uber jar as a Docker image. However, there are various downsides to copying and running the uber jar as-is in the Docker image. There’s always a certain amount of overhead when running an uber jar without unpacking it, and in a containerized environment this can be noticeable. The other issue is that putting your application’s code and all its dependencies in one layer in the Docker image is not optimal. Since you probably recompile your code more often than you upgrade the version of Spring Boot you use, it’s often better to separate things a bit more. If you put jar files in the layer before your application classes, Docker often only needs to change the very bottom layer and can pick others up from its cache.

## Layering Docker Images

To make it easier to create optimized Docker images, Spring Boot supports adding a layer index file to the jar. It provides a list of layers and the parts of the jar that should be contained within them. The list of layers in the index is ordered based on the order in which the layers should be added to the Docker/OCI image. Out-of-the-box, the following layers are supported:

-
`dependencies`(for regular released dependencies)
-
`spring-boot-loader`(for everything under`org/springframework/boot/loader`)
-
`snapshot-dependencies`(for snapshot dependencies)
-
`application`(for application classes and resources)

The following shows an example of a `layers.idx` file:

```
- "dependencies":
 - BOOT-INF/lib/library1.jar
 - BOOT-INF/lib/library2.jar
- "spring-boot-loader":
 - org/springframework/boot/loader/launch/JarLauncher.class
 - ... <other classes>
- "snapshot-dependencies":
 - BOOT-INF/lib/library3-SNAPSHOT.jar
- "application":
 - META-INF/MANIFEST.MF
 - BOOT-INF/classes/a/b/C.class
```
This layering is designed to separate code based on how likely it is to change between application builds. Library code is less likely to change between builds, so it is placed in its own layers to allow tooling to re-use the layers from cache. Application code is more likely to change between builds so it is isolated in a separate layer.

Spring Boot also supports layering for war files with the help of a `layers.idx`.

For Maven, see the packaging layered jar or war section for more details on adding a layer index to the archive. For Gradle, see the packaging layered jar or war section of the Gradle plugin documentation.

# Dockerfiles

While it is possible to convert a Spring Boot uber jar into a Docker image with just a few lines in the `Dockerfile`, using the layering feature will result in an optimized image.
When you create a jar containing the layers index file, the `spring-boot-jarmode-tools` jar will be added as a dependency to your jar.
With this jar on the classpath, you can launch your application in a special mode which allows the bootstrap code to run something entirely different from your application, for example, something that extracts the layers.

Here’s how you can launch your jar with a `tools` jar mode:

`$ java -Djarmode=tools -jar my-app.jar`This will provide the following output:

Usage: java -Djarmode=tools -jar my-app.jar Available commands: extract Extract the contents from the jar list-layers List layers from the jar that can be extracted help Help about any command

The `extract` command can be used to easily split the application into layers to be added to the `Dockerfile`.
Here is an example of a `Dockerfile` using `jarmode`.

```
# Perform the extraction in a separate builder container
FROM bellsoft/liberica-openjre-debian:25-cds AS builder
WORKDIR /builder
# This points to the built jar file in the target folder
# Adjust this to 'build/libs/*.jar' if you're using Gradle
ARG JAR_FILE=target/*.jar
# Copy the jar file to the working directory and rename it to application.jar
COPY ${JAR_FILE} application.jar
# Extract the jar file using an efficient layout
RUN java -Djarmode=tools -jar application.jar extract --layers --destination extracted
# This is the runtime container
FROM bellsoft/liberica-openjre-debian:25-cds
WORKDIR /application
# Copy the extracted jar contents from the builder container into the working directory in the runtime container
# Every copy step creates a new docker layer
# This allows docker to only pull the changes it really needs
COPY --from=builder /builder/extracted/dependencies/ ./
COPY --from=builder /builder/extracted/spring-boot-loader/ ./
COPY --from=builder /builder/extracted/snapshot-dependencies/ ./
COPY --from=builder /builder/extracted/application/ ./
# Start the application jar - this is not the uber jar used by the builder
# This jar only contains application code and references to the extracted jar files
# This layout is efficient to start up and AOT cache (and CDS) friendly
ENTRYPOINT ["java", "-jar", "application.jar"]
```
Assuming the above `Dockerfile` is in the current directory, your Docker image can be built with `docker build .`, or optionally specifying the path to your application jar, as shown in the following example:

`$ docker build --build-arg JAR_FILE=path/to/myapp.jar .`This is a multi-stage `Dockerfile`.
The builder stage extracts the directories that are needed later.
Each of the `COPY` commands relates to the layers extracted by the jarmode.

Of course, a `Dockerfile` can be written without using the `jarmode`.
You can use some combination of `unzip` and `mv` to move things to the right layer but `jarmode` simplifies that.
Additionally, the layout created by the `jarmode` is AOT cache (and CDS) friendly out of the box.

## AOT cache

If you are using Java 25 or above, and want to additionally enable the AOT cache, you can use this `Dockerfile`:

```
# Perform the extraction in a separate builder container
FROM bellsoft/liberica-openjre-debian:25-cds AS builder
WORKDIR /builder
# This points to the built jar file in the target folder
# Adjust this to 'build/libs/*.jar' if you're using Gradle
ARG JAR_FILE=target/*.jar
# Copy the jar file to the working directory and rename it to application.jar
COPY ${JAR_FILE} application.jar
# Extract the jar file using an efficient layout
RUN java -Djarmode=tools -jar application.jar extract --layers --destination extracted
# This is the runtime container
FROM bellsoft/liberica-openjre-debian:25-cds
WORKDIR /application
# Copy the extracted jar contents from the builder container into the working directory in the runtime container
# Every copy step creates a new docker layer
# This allows docker to only pull the changes it really needs
COPY --from=builder /builder/extracted/dependencies/ ./
COPY --from=builder /builder/extracted/spring-boot-loader/ ./
COPY --from=builder /builder/extracted/snapshot-dependencies/ ./
COPY --from=builder /builder/extracted/application/ ./
# Execute the AOT cache training run
RUN java -XX:AOTCacheOutput=app.aot -Dspring.context.exit=onRefresh -jar application.jar
# Start the application jar with AOT cache enabled - this is not the uber jar used by the builder
# This jar only contains application code and references to the extracted jar files
# This layout is efficient to start up and AOT cache friendly
ENTRYPOINT ["java", "-XX:AOTCache=app.aot", "-jar", "application.jar"]
```
This is mostly the same as the above `Dockerfile`.
As the last steps, it creates the AOT cache file by doing a training run and passes the AOT cache parameter to `java -jar`.

## CDS

| If you’re using Java 24 or later, please use AOT cache instead of CDS. |

If you want to additionally enable CDS, you can use this `Dockerfile`:

```
# Perform the extraction in a separate builder container
FROM bellsoft/liberica-openjre-debian:25-cds AS builder
WORKDIR /builder
# This points to the built jar file in the target folder
# Adjust this to 'build/libs/*.jar' if you're using Gradle
ARG JAR_FILE=target/*.jar
# Copy the jar file to the working directory and rename it to application.jar
COPY ${JAR_FILE} application.jar
# Extract the jar file using an efficient layout
RUN java -Djarmode=tools -jar application.jar extract --layers --destination extracted
# This is the runtime container
FROM bellsoft/liberica-openjre-debian:25-cds
WORKDIR /application
# Copy the extracted jar contents from the builder container into the working directory in the runtime container
# Every copy step creates a new docker layer
# This allows docker to only pull the changes it really needs
COPY --from=builder /builder/extracted/dependencies/ ./
COPY --from=builder /builder/extracted/spring-boot-loader/ ./
COPY --from=builder /builder/extracted/snapshot-dependencies/ ./
COPY --from=builder /builder/extracted/application/ ./
# Execute the CDS training run
RUN java -XX:ArchiveClassesAtExit=application.jsa -Dspring.context.exit=onRefresh -jar application.jar
# Start the application jar with CDS enabled - this is not the uber jar used by the builder
# This jar only contains application code and references to the extracted jar files
# This layout is efficient to start up and CDS friendly
ENTRYPOINT ["java", "-XX:SharedArchiveFile=application.jsa", "-jar", "application.jar"]
```
This is mostly the same as the above `Dockerfile`.
As the last steps, it creates the CDS archive by doing a training run and passes the CDS parameter to `java -jar`.

| If you are using Java 25 or above, we recommend using an AOT cache instead of CDS. |

# Cloud Native Buildpacks

Docker images can be built directly from your Maven or Gradle plugin using Cloud Native Buildpacks.
If you’ve ever used an application platform such as Cloud Foundry or Heroku then you’ve probably used a buildpack.
Buildpacks are the part of the platform that takes your application and converts it into something that the platform can actually run.
For example, Cloud Foundry’s Java buildpack will notice that you’re pushing a `.jar` file and automatically add a relevant JRE.

With Cloud Native Buildpacks, you can create Docker compatible images that you can run anywhere. Spring Boot includes buildpack support directly for both Maven and Gradle. This means you can just type a single command and quickly get a sensible image into your locally running Docker daemon.

| The Paketo Spring Boot buildpack supports the `layers.idx`file, so any layer customization that is applied to it will be reflected in the image created by the buildpacks. |

| In order to achieve reproducible builds and container image caching, buildpacks can manipulate the application resources metadata (such as the file "last modified" information).
You should ensure that your application does not rely on that metadata at runtime.
Spring Boot can use that information when serving static resources, but this can be disabled with `spring.web.resources.cache.use-last-modified`. |

# Production-ready Features

Spring Boot includes a number of additional features to help you monitor and manage your application when you push it to production. You can choose to manage and monitor your application by using HTTP endpoints or with JMX. Auditing, health, and metrics gathering can also be automatically applied to your application.

# Enabling Production-ready Features

The `spring-boot-actuator` module provides all of Spring Boot’s production-ready features.
The recommended way to enable the features is to add a dependency on the `spring-boot-starter-actuator` starter.

To add the actuator to a Maven-based project, add the following starter dependency:

```
<dependencies>
	<dependency>
 <groupId>org.springframework.boot</groupId>
 <artifactId>spring-boot-starter-actuator</artifactId>
	</dependency>
</dependencies>
```
For Gradle, use the following declaration:

```
dependencies {
	implementation 'org.springframework.boot:spring-boot-starter-actuator'
}
```

# Endpoints

Actuator endpoints let you monitor and interact with your application.
Spring Boot includes a number of built-in endpoints and lets you add your own.
For example, the `health` endpoint provides basic application health information.

You can control access to each individual endpoint and expose them (make them remotely accessible) over HTTP or JMX.
An endpoint is considered to be available when access to it is permitted and it is exposed.
The built-in endpoints are auto-configured only when they are available.
Most applications choose exposure over HTTP, where the ID of the endpoint and a prefix of `/actuator` is mapped to a URL.
For example, by default, the `health` endpoint is mapped to `/actuator/health`.

| To learn more about the Actuator’s endpoints and their request and response formats, see the API documentation. |

The following technology-agnostic endpoints are available:

| ID | Description |
|---|---|
|
 | Exposes audit events information for the current application.
 Requires an |
|
 | Displays a complete list of all the Spring beans in your application. |
|
 | Exposes available caches. |
|
 | Shows the conditions that were evaluated on configuration and auto-configuration classes and the reasons why they did or did not match. |
|
 | Displays a collated list of all |
|
 | Exposes properties from Spring’s |
|
 | Shows any Flyway database migrations that have been applied.
 Requires one or more |
|
 | Shows application health information. |
|
 | Displays HTTP exchange information (by default, the last 100 HTTP request-response exchanges).
 Requires an |
|
 | Displays arbitrary application info. |
|
 | Shows the Spring Integration graph.
 Requires a dependency on |
|
 | Shows and modifies the configuration of loggers in the application. |
|
 | Shows any Liquibase database migrations that have been applied.
 Requires one or more |
|
 | Shows “metrics” information for the current application to diagnose the metrics the application has recorded. |
|
 | Displays a collated list of all |
|
 | Shows information about Quartz Scheduler jobs. Subject to sanitization. |
|
 | Displays the scheduled tasks in your application. |
|
 | Allows retrieval and deletion of user sessions from a Spring Session-backed session store. Requires a servlet-based web application that uses Spring Session. |
|
 | Lets the application be gracefully shutdown. Only works when using jar packaging. Disabled by default. |
|
 | Shows the startup steps data collected by the |
|
 | Performs a thread dump. |

If your application is a web application (Spring MVC, Spring WebFlux, or Jersey), you can use the following additional endpoints:

| ID | Description |
|---|---|
|
 | Returns a heap dump file.
 On a HotSpot JVM, an |
|
 | Returns the contents of the logfile (if the |
|
 | Exposes metrics in a format that can be scraped by a Prometheus server.
 Requires a dependency on |

## Controlling Access to Endpoints

By default, access to all endpoints except for `shutdown` and `heapdump` is unrestricted.
To configure the permitted access to an endpoint, use its `management.endpoint.<id>.access` property.
The following example allows unrestricted access to the `shutdown` endpoint:

-
Properties
-
YAML

`management.endpoint.shutdown.access=unrestricted````
management:
 endpoint:
 shutdown:
 access: unrestricted
```
If you prefer access to be opt-in rather than opt-out, set the `management.endpoints.access.default` property to `none` and use individual endpoint `access` properties to opt back in.
The following example allows read-only access to the `loggers` endpoint and denies access to all other endpoints:

-
Properties
-
YAML

```
management.endpoints.access.default=none
management.endpoint.loggers.access=read-only
```
```
management:
 endpoints:
 access:
 default: none
 endpoint:
 loggers:
 access: read-only
```
| Inaccessible endpoints are removed entirely from the application context.
If you want to change only the technologies over which an endpoint is exposed, use the `include`and`exclude`properties instead. |

### Limiting Access

Application-wide endpoint access can be limited using the `management.endpoints.access.max-permitted` property.
This property takes precedence over the default access or an individual endpoint’s access level.
Set it to `none` to make all endpoints inaccessible.
Set it to `read-only` to only allow read access to endpoints.

For `@Endpoint`, `@JmxEndpoint`, and `@WebEndpoint`, read access equates to the endpoint methods annotated with `@ReadOperation`.
For `@ControllerEndpoint` and `@RestControllerEndpoint`, read access equates to request mappings that can handle `GET` and `HEAD` requests.
For `@ServletEndpoint`, read access equates to `GET` and `HEAD` requests.

## Exposing Endpoints

By default, only the health endpoint is exposed over HTTP and JMX. Since Endpoints may contain sensitive information, you should carefully consider when to expose them.

To change which endpoints are exposed, use the following technology-specific `include` and `exclude` properties:

| Property | Default |
|---|---|
|
 | |
|
 |
 |
|
 | |
|
 |
 |

The `include` property lists the IDs of the endpoints that are exposed.
The `exclude` property lists the IDs of the endpoints that should not be exposed.
The `exclude` property takes precedence over the `include` property.
You can configure both the `include` and the `exclude` properties with a list of endpoint IDs.

For example, to only expose the `health` and `info` endpoints over JMX, use the following property:

-
Properties
-
YAML

`management.endpoints.jmx.exposure.include=health,info````
management:
 endpoints:
 jmx:
 exposure:
 include: "health,info"
```
`*` can be used to select all endpoints.
For example, to expose everything over HTTP except the `env` and `beans` endpoints, use the following properties:

-
Properties
-
YAML

```
management.endpoints.web.exposure.include=*
management.endpoints.web.exposure.exclude=env,beans
```
```
management:
 endpoints:
 web:
 exposure:
 include: "*"
 exclude: "env,beans"
```
| `*`has a special meaning in YAML, so be sure to add quotation marks if you want to include (or exclude) all endpoints. |

| If your application is exposed publicly, we strongly recommend that you also secure your endpoints. |

| If you want to implement your own strategy for when endpoints are exposed, you can register an `EndpointFilter`bean. |

## Security

For security purposes, only the `/health` endpoint is exposed over HTTP by default.
You can use the `management.endpoints.web.exposure.include` property to configure the endpoints that are exposed.

| Before setting the `management.endpoints.web.exposure.include`, ensure that the exposed actuators do not contain sensitive information, are secured by placing them behind a firewall, or are secured by something like Spring Security. |

If Spring Security is on the classpath and no other `SecurityFilterChain` bean is present, all actuators other than `/health` are secured by Spring Boot auto-configuration.
If you define a custom `SecurityFilterChain` bean, Spring Boot auto-configuration backs off and lets you fully control the actuator access rules.

If you wish to configure custom security for HTTP endpoints (for example, to allow only users with a certain role to access them), Spring Boot provides some convenient `RequestMatcher` objects that you can use in combination with Spring Security.

A typical Spring Security configuration might look something like the following example:

-
Java
-
Kotlin

```
import org.springframework.boot.security.autoconfigure.actuate.web.servlet.EndpointRequest;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.SecurityFilterChain;
import static org.springframework.security.config.Customizer.withDefaults;
@Configuration(proxyBeanMethods = false)
public class MySecurityConfiguration {
	@Bean
	public SecurityFilterChain securityFilterChain(HttpSecurity http) {
 http.securityMatcher(EndpointRequest.toAnyEndpoint());
 http.authorizeHttpRequests((requests) -> requests.anyRequest().hasRole("ENDPOINT_ADMIN"));
 http.httpBasic(withDefaults());
 return http.build();
	}
}
```
```
import org.springframework.boot.security.autoconfigure.actuate.web.servlet.EndpointRequest
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.security.config.Customizer.withDefaults
import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.web.SecurityFilterChain
@Configuration(proxyBeanMethods = false)
class MySecurityConfiguration {
	@Bean
	fun securityFilterChain(http: HttpSecurity): SecurityFilterChain {
 http.securityMatcher(EndpointRequest.toAnyEndpoint()).authorizeHttpRequests { requests ->
 requests.anyRequest().hasRole("ENDPOINT_ADMIN")
 }
 http.httpBasic(withDefaults())
 return http.build()
	}
}
```
The preceding example uses `EndpointRequest.toAnyEndpoint()` to match a request to any endpoint and then ensures that all have the `ENDPOINT_ADMIN` role.
Several other matcher methods are also available on `EndpointRequest`.
See the API documentation for details.

| When matching for Actuator endpoints, `EndpointRequest.to("endpoint")`will consider the endpoint root and all its subpaths,
effectively matching`"/actuator/endpoint/**"`even if the endpoint does not declare nested routes. |

If you deploy applications behind a firewall, you may prefer that all your actuator endpoints can be accessed without requiring authentication.
You can do so by changing the `management.endpoints.web.exposure.include` property, as follows:

-
Properties
-
YAML

`management.endpoints.web.exposure.include=*````
management:
 endpoints:
 web:
 exposure:
 include: "*"
```
Additionally, if Spring Security is present, you would need to add custom security configuration that allows unauthenticated access to the endpoints, as the following example shows:

-
Java
-
Kotlin

```
import org.springframework.boot.security.autoconfigure.actuate.web.servlet.EndpointRequest;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.SecurityFilterChain;
@Configuration(proxyBeanMethods = false)
public class MySecurityConfiguration {
	@Bean
	public SecurityFilterChain securityFilterChain(HttpSecurity http) {
 http.securityMatcher(EndpointRequest.toAnyEndpoint());
 http.authorizeHttpRequests((requests) -> requests.anyRequest().permitAll());
 return http.build();
	}
}
```
```
import org.springframework.boot.security.autoconfigure.actuate.web.servlet.EndpointRequest
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.web.SecurityFilterChain
@Configuration(proxyBeanMethods = false)
class MySecurityConfiguration {
	@Bean
	fun securityFilterChain(http: HttpSecurity): SecurityFilterChain {
 http.securityMatcher(EndpointRequest.toAnyEndpoint()).authorizeHttpRequests { requests ->
 requests.anyRequest().permitAll()
 }
 return http.build()
	}
}
```
| In both of the preceding examples, the configuration applies only to the actuator endpoints.
Since Spring Boot’s security configuration backs off completely in the presence of any `SecurityFilterChain`bean, you need to configure an additional`SecurityFilterChain`bean with rules that apply to the rest of the application. |

### Cross Site Request Forgery Protection

Since Spring Boot relies on Spring Security’s defaults, CSRF protection is turned on by default.
This means that the actuator endpoints that require a `POST` (shutdown and loggers endpoints), a `PUT`, or a `DELETE` get a 403 (forbidden) error when the default security configuration is in use.

| We recommend disabling CSRF protection completely only if you are creating a service that is used by non-browser clients. |

You can find additional information about CSRF protection in the Spring Security Reference Guide.

## Configuring Endpoints

Endpoints automatically cache responses to read operations that do not take any parameters.
To configure the amount of time for which an endpoint caches a response, use its `cache.time-to-live` property.
The following example sets the time-to-live of the `beans` endpoint’s cache to 10 seconds:

-
Properties
-
YAML

`management.endpoint.beans.cache.time-to-live=10s````
management:
 endpoint:
 beans:
 cache:
 time-to-live: "10s"
```
| The `management.endpoint.<name>`prefix uniquely identifies the endpoint that is being configured. |

## Sanitize Sensitive Values

Information returned by the `/env`, `/configprops` and `/quartz` endpoints can be sensitive, so by default values are always fully sanitized (replaced by `******`).

Values can only be viewed in an unsanitized form when:

-
The `show-values`property has been set to something other than`never`
-
No custom `SanitizingFunction`beans apply

The `show-values` property can be configured for sanitizable endpoints to one of the following values:

-
`never`- values are always fully sanitized (replaced by`******`)
-
`always`- values are shown to all users (as long as no`SanitizingFunction`bean applies)
-
`when-authorized`- values are shown only to authorized users (as long as no`SanitizingFunction`bean applies)

For HTTP endpoints, a user is considered to be authorized if they have authenticated and have the roles configured by the endpoint’s roles property. By default, any authenticated user is authorized.

For JMX endpoints, all users are always authorized.

The following example allows all users with the `admin` role to view values from the `/env` endpoint in their original form.
Unauthorized users, or users without the `admin` role, will see only sanitized values.

-
Properties
-
YAML

```
management.endpoint.env.show-values=when-authorized
management.endpoint.env.roles=admin
```
```
management:
 endpoint:
 env:
 show-values: when-authorized
 roles: "admin"
```
| This example assumes that no `SanitizingFunction`beans have been defined. |

## Hypermedia for Actuator Web Endpoints

A “discovery page” is added with links to all the endpoints.
The “discovery page” is available on `/actuator` by default.

To disable the “discovery page”, add the following property to your application properties:

-
Properties
-
YAML

`management.endpoints.web.discovery.enabled=false````
management:
 endpoints:
 web:
 discovery:
 enabled: false
```
When a custom management context path is configured, the “discovery page” automatically moves from `/actuator` to the root of the management context.
For example, if the management context path is `/management`, the discovery page is available from `/management`.
When the management context path is set to `/`, the discovery page is disabled to prevent the possibility of a clash with other mappings.

## CORS Support

Cross-origin resource sharing (CORS) is a W3C specification that lets you specify in a flexible way what kind of cross-domain requests are authorized. If you use Spring MVC or Spring WebFlux, you can configure Actuator’s web endpoints to support such scenarios.

CORS support is disabled by default and is only enabled once you have set the `management.endpoints.web.cors.allowed-origins` property.
The following configuration permits `GET` and `POST` calls from the `example.com` domain:

-
Properties
-
YAML

```
management.endpoints.web.cors.allowed-origins=https://example.com
management.endpoints.web.cors.allowed-methods=GET,POST
```
```
management:
 endpoints:
 web:
 cors:
 allowed-origins: "https://example.com"
 allowed-methods: "GET,POST"
```
| See `CorsEndpointProperties`for a complete list of options. |

## JSON

When working with JSON, Jackson is used for serialization and deserialization.
By default, an isolated `JsonMapper` is used.
This isolation means that it does not share the same configuration as the application’s `JsonMapper` and it is not affected by `spring.jackson.*` properties.
To disable this behavior and configure Actuator to use the application’s `JsonMapper`, set `management.endpoints.jackson.isolated-json-mapper` to `false`.
Alternatively, you can define your own `EndpointJsonMapper` bean that produces a `JsonMapper` that meets your needs.
Actuator will then use it for JSON processing.

## Implementing Custom Endpoints

If you add a `@Bean` annotated with `@Endpoint`, any methods annotated with `@ReadOperation`, `@WriteOperation`, or `@DeleteOperation` are automatically exposed over JMX and, in a web application, over HTTP as well.
Endpoints can be exposed over HTTP by using Jersey, Spring MVC, or Spring WebFlux.
If both Jersey and Spring MVC are available, Spring MVC is used.

The following example exposes a read operation that returns a custom object:

-
Java
-
Kotlin

```
	@ReadOperation
	public CustomData getData() {
 return new CustomData("test", 5);
	}
```
```
	@ReadOperation
	fun getData(): CustomData {
 return CustomData("test", 5)
	}
```
You can also write technology-specific endpoints by using `@JmxEndpoint` or `@WebEndpoint`.
These endpoints are restricted to their respective technologies.
For example, `@WebEndpoint` is exposed only over HTTP and not over JMX.

You can write technology-specific extensions by using `@EndpointWebExtension` and `@EndpointJmxExtension`.
These annotations let you provide technology-specific operations to augment an existing endpoint.
An endpoint may have at most one extension of each type.

Finally, if you need access to web-framework-specific functionality, you can implement servlet or Spring `@Controller` and `@RestController` endpoints at the cost of them not being available over JMX or when using a different web framework.

### Receiving Input

Operations on an endpoint receive input through their parameters.
When exposed over the web, the values for these parameters are taken from the URL’s query parameters and from the JSON request body.
When exposed over JMX, the parameters are mapped to the parameters of the MBean’s operations.
Parameters are required by default.
They can be made optional by annotating them with JSpecify’s `@Nullable`.
Kotlin null safety is also supported.

You can map each root property in the JSON request body to a parameter of the endpoint. Consider the following JSON request body:

```
{
	"name": "test",
	"counter": 42
}
```
You can use this to invoke a write operation that takes `String name` and `int counter` parameters, as the following example shows:

-
Java
-
Kotlin

```
	@WriteOperation
	public void updateData(String name, int counter) {
 // injects "test" and 42
	}
```
```
	@WriteOperation
	fun updateData(name: String?, counter: Int) {
 // injects "test" and 42
	}
```
| Because endpoints are technology agnostic, only simple types can be specified in the method signature.
In particular, declaring a single parameter with a `CustomData`type that defines a`name`and`counter`properties is not supported. |

| To let the input be mapped to the operation method’s parameters, Java code that implements an endpoint should be compiled with `-parameters`.
For Kotlin code, please review the recommendation of the Spring Framework reference.
This will happen automatically if you use Spring Boot’s Gradle plugin or if you use Maven and`spring-boot-starter-parent`. |

#### Input Type Conversion

The parameters passed to endpoint operation methods are, if necessary, automatically converted to the required type.
Before calling an operation method, the input received over JMX or HTTP is converted to the required types by using an instance of `ApplicationConversionService` as well as any `Converter` or `GenericConverter` beans qualified with `@EndpointConverter`.

### Custom Web Endpoints

Operations on an `@Endpoint`, `@WebEndpoint`, or `@EndpointWebExtension` are automatically exposed over HTTP using Jersey, Spring MVC, or Spring WebFlux.
If both Jersey and Spring MVC are available, Spring MVC is used.

#### Web Endpoint Request Predicates

A request predicate is automatically generated for each operation on a web-exposed endpoint.

#### Path

The path of the predicate is determined by the ID of the endpoint and the base path of the web-exposed endpoints.
The default base path is `/actuator`.
For example, an endpoint with an ID of `sessions` uses `/actuator/sessions` as its path in the predicate.

You can further customize the path by annotating one or more parameters of the operation method with `@Selector`.
Such a parameter is added to the path predicate as a path variable.
The variable’s value is passed into the operation method when the endpoint operation is invoked.
If you want to capture all remaining path elements, you can add `@Selector(Match=ALL_REMAINING)` to the last parameter and make it a type that is conversion-compatible with a `String[]`.

#### HTTP method

The HTTP method of the predicate is determined by the operation type, as shown in the following table:

| Operation | HTTP method |
|---|---|
|
 | |
|
 | |
|
 |

#### Consumes

For a `@WriteOperation` (HTTP `POST`) that uses the request body, the `consumes` clause of the predicate is `application/vnd.spring-boot.actuator.v2+json, application/json`.
For all other operations, the `consumes` clause is empty.

#### Produces

The `produces` clause of the predicate can be determined by the `produces` attribute of the `@DeleteOperation`, `@ReadOperation`, and `@WriteOperation` annotations.
The attribute is optional.
If it is not used, the `produces` clause is determined automatically.

#### Web Endpoint Response Status

The default response status for an endpoint operation depends on the operation type (read, write, or delete) and what, if anything, the operation returns.

If a `@ReadOperation` returns a value, the response status will be 200 (OK).
If it does not return a value, the response status will be 404 (Not Found).

If a `@WriteOperation` or `@DeleteOperation` returns a value, the response status will be 200 (OK).
If it does not return a value, the response status will be 204 (No Content).

If an operation is invoked without a required parameter or with a parameter that cannot be converted to the required type, the operation method is not called, and the response status will be 400 (Bad Request).

#### Web Endpoint Range Requests

You can use an HTTP range request to request part of an HTTP resource.
When using Spring MVC or Spring Web Flux, operations that return a `Resource` automatically support range requests.

| Range requests are not supported when using Jersey. |

#### Web Endpoint Security

An operation on a web endpoint or a web-specific endpoint extension can receive the current `Principal` or `SecurityContext` as a method parameter.
The former is typically used in conjunction with `@Nullable` to provide different behavior for authenticated and unauthenticated users.
The latter is typically used to perform authorization checks by using its `isUserInRole(String)` method.

## Health Information

You can use health information to check the status of your running application.
It is often used by monitoring software to alert someone when a production system goes down.
The information exposed by the `health` endpoint depends on the `management.endpoint.health.show-details` and `management.endpoint.health.show-components` properties, which can be configured with one of the following values:

| Name | Description |
|---|---|
|
 | Details are never shown. |
|
 | Details are shown only to authorized users.
 Authorized roles can be configured by using |
|
 | Details are shown to all users. |

The default value is `never`.
A user is considered to be authorized when they are in one or more of the endpoint’s roles.
If the endpoint has no configured roles (the default), all authenticated users are considered to be authorized.
You can configure the roles by using the `management.endpoint.health.roles` property.

| If you have secured your application and wish to use `always`, your security configuration must permit access to the health endpoint for both authenticated and unauthenticated users. |

Health information is collected from the content of a `HealthContributorRegistry` (by default, all `HealthContributor` instances defined in your `ApplicationContext`).
Spring Boot includes a number of auto-configured `HealthContributor` beans, and you can also write your own.

A `HealthContributor` can be either a `HealthIndicator` or a `CompositeHealthContributor`.
A `HealthIndicator` provides actual health information, including a `Status`.
A `CompositeHealthContributor` provides a composite of other `HealthContributor` instances.
Taken together, contributors form a tree structure to represent the overall system health.

By default, the final system health is derived by a `StatusAggregator`, which sorts the statuses from each `HealthIndicator` based on an ordered list of statuses.
The first status in the sorted list is used as the overall health status.
If no `HealthIndicator` returns a status that is known to the `StatusAggregator`, an `UNKNOWN` status is used.

| You can use the `HealthContributorRegistry`to register and unregister health indicators at runtime. |

### Auto-configured HealthIndicators

When appropriate, Spring Boot auto-configures the `HealthIndicator` beans listed in the following table.
You can also enable or disable selected indicators by configuring `management.health.key.enabled`,
with the `key` listed in the following table:

| Key | Name | Description |
|---|---|---|
|
 | Checks that a Cassandra database is up. | |
|
 | Checks that a Couchbase cluster is up. | |
|
 | Checks that a connection to | |
|
 | Checks for low disk space. | |
|
 | Checks that an Elasticsearch cluster is up. | |
|
 | Checks that a Hazelcast server is up. | |
|
 | Checks that a JMS broker is up. | |
|
 | Checks that an LDAP server is up. | |
|
 | Checks that a mail server is up. | |
|
 | Checks that a Mongo database is up. | |
|
 | Checks that a Neo4j database is up. | |
|
 | Always responds with | |
|
 | Checks that a Rabbit server is up. | |
|
 | Checks that a Redis server is up. | |
|
 | Checks that SSL certificates are ok. |

| You can disable them all by setting the `management.health.defaults.enabled`property. |

| The `ssl``HealthIndicator`has a "warning threshold" property named`management.health.ssl.certificate-validity-warning-threshold`.
You can use this threshold to give yourself enough lead time to rotate the soon-to-be-expired certificate.
If an SSL certificate will become invalid within the period defined by this threshold, the`HealthIndicator`will report this in the details section of its response where`details.validChains.certificates.[*].validity.status`will have the value`WILL_EXPIRE_SOON`. |

Additional `HealthIndicator` beans are enabled by default:

| Key | Name | Description |
|---|---|---|
|
 | Exposes the “Liveness” application availability state. | |
|
 | Exposes the “Readiness” application availability state. |

These can be disabled by using the `management.endpoint.health.probes.enabled` configuration property.

### Writing Custom HealthIndicators

To provide custom health information, you can register Spring beans that implement the `HealthIndicator` interface.
You need to provide an implementation of the `health()` method and return a `Health` response.
The `Health` response should include a status and can optionally include additional details to be displayed.
The following code shows a sample `HealthIndicator` implementation:

-
Java
-
Kotlin

```
import org.springframework.boot.health.contributor.Health;
import org.springframework.boot.health.contributor.HealthIndicator;
import org.springframework.stereotype.Component;
@Component
public class MyHealthIndicator implements HealthIndicator {
	@Override
	public Health health() {
 int errorCode = check();
 if (errorCode != 0) {
 return Health.down().withDetail("Error Code", errorCode).build();
 }
 return Health.up().build();
	}
	private int check() {
 // perform some specific health check
 return ...
	}
}
```
```
import org.springframework.boot.health.contributor.Health
import org.springframework.boot.health.contributor.HealthIndicator
import org.springframework.stereotype.Component
@Component
class MyHealthIndicator : HealthIndicator {
	override fun health(): Health {
 val errorCode = check()
 if (errorCode != 0) {
 return Health.down().withDetail("Error Code", errorCode).build()
 }
 return Health.up().build()
	}
	private fun check(): Int {
 // perform some specific health check
 return ...
	}
}
```
| The identifier for a given `HealthIndicator`is the name of the bean without the`HealthIndicator`suffix, if it exists.
In the preceding example, the health information is available in an entry named`my`. |

| Health indicators are usually called over HTTP and need to respond before any connection timeouts.
Spring Boot will log a warning message for any health indicator that takes longer than 10 seconds to respond.
If you want to configure this threshold, you can use the `management.endpoint.health.logging.slow-indicator-threshold`property. |

In addition to Spring Boot’s predefined `Status` types, `Health` can return a custom `Status` that represents a new system state.
In such cases, you also need to provide a custom implementation of the `StatusAggregator` interface, or you must configure the default implementation by using the `management.endpoint.health.status.order` configuration property.

For example, assume a new `Status` with a code of `FATAL` is being used in one of your `HealthIndicator` implementations.
To configure the severity order, add the following property to your application properties:

-
Properties
-
YAML

`management.endpoint.health.status.order=fatal,down,out-of-service,unknown,up````
management:
 endpoint:
 health:
 status:
 order: "fatal,down,out-of-service,unknown,up"
```
The HTTP status code in the response reflects the overall health status.
By default, `OUT_OF_SERVICE` and `DOWN` map to 503.
Any unmapped health statuses, including `UP`, map to 200.
You might also want to register custom status mappings if you access the health endpoint over HTTP.
Configuring a custom mapping disables the defaults mappings for `DOWN` and `OUT_OF_SERVICE`.
If you want to retain the default mappings, you must explicitly configure them, alongside any custom mappings.
For example, the following property maps `FATAL` to 503 (service unavailable) and retains the default mappings for `DOWN` and `OUT_OF_SERVICE`:

-
Properties
-
YAML

```
management.endpoint.health.status.http-mapping.down=503
management.endpoint.health.status.http-mapping.fatal=503
management.endpoint.health.status.http-mapping.out-of-service=503
```
```
management:
 endpoint:
 health:
 status:
 http-mapping:
 down: 503
 fatal: 503
 out-of-service: 503
```
| If you need more control, you can define your own `HttpCodeStatusMapper`bean. |

The following table shows the default status mappings for the built-in statuses:

| Status | Mapping |
|---|---|
|
 |
 |
|
 |
 |
|
 | No mapping by default, so HTTP status is |
|
 | No mapping by default, so HTTP status is |

### Reactive Health Indicators

For reactive applications, such as those that use Spring WebFlux, `ReactiveHealthContributor` provides a non-blocking contract for getting application health.
Similar to a traditional `HealthContributor`, health information is collected from the content of a `ReactiveHealthContributorRegistry` (by default, all `HealthContributor` and `ReactiveHealthContributor` instances defined in your `ApplicationContext`).
Regular `HealthContributor` instances that do not check against a reactive API are executed on the elastic scheduler.

| In a reactive application, you should use the `ReactiveHealthContributorRegistry`to register and unregister health indicators at runtime.
If you need to register a regular`HealthContributor`, you should wrap it with`ReactiveHealthContributor#adapt`. |

To provide custom health information from a reactive API, you can register Spring beans that implement the `ReactiveHealthIndicator` interface.
The following code shows a sample `ReactiveHealthIndicator` implementation:

-
Java
-
Kotlin

```
import reactor.core.publisher.Mono;
import org.springframework.boot.health.contributor.Health;
import org.springframework.boot.health.contributor.ReactiveHealthIndicator;
import org.springframework.stereotype.Component;
@Component
public class MyReactiveHealthIndicator implements ReactiveHealthIndicator {
	@Override
	public Mono<Health> health() {
 return doHealthCheck().onErrorResume((exception) ->
 Mono.just(new Health.Builder().down(exception).build()));
	}
	private Mono<Health> doHealthCheck() {
 // perform some specific health check
 return ...
	}
}
```
```
import org.springframework.boot.health.contributor.Health
import org.springframework.boot.health.contributor.ReactiveHealthIndicator
import org.springframework.stereotype.Component
import reactor.core.publisher.Mono
@Component
class MyReactiveHealthIndicator : ReactiveHealthIndicator {
	override fun health(): Mono<Health> {
 return doHealthCheck().onErrorResume { exception: Throwable ->
 Mono.just(Health.Builder().down(exception).build())
 }
	}
	private fun doHealthCheck(): Mono<Health> {
 // perform some specific health check
 return ...
	}
}
```
| To handle the error automatically, consider extending from `AbstractReactiveHealthIndicator`. |

### Auto-configured ReactiveHealthIndicators

When appropriate, Spring Boot auto-configures the following `ReactiveHealthIndicator` beans:

| Key | Name | Description |
|---|---|---|
|
 | Checks that a Cassandra database is up. | |
|
 | Checks that a Couchbase cluster is up. | |
|
 | Checks that an Elasticsearch cluster is up. | |
|
 | Checks that a Mongo database is up. | |
|
 | Checks that a Neo4j database is up. | |
|
 | Checks that a Redis server is up. |

| If necessary, reactive indicators replace the regular ones.
Also, any `HealthIndicator`that is not handled explicitly is wrapped automatically. |

### Health Groups

It is sometimes useful to organize health indicators into groups that you can use for different purposes.

To create a health indicator group, you can use the `management.endpoint.health.group.<name>` property and specify a list of health indicator IDs to `include` or `exclude`.
For example, to create a group that includes only database indicators you can define the following:

-
Properties
-
YAML

`management.endpoint.health.group.custom.include=db````
management:
 endpoint:
 health:
 group:
 custom:
 include: "db"
```
You can then check the result by hitting `localhost:8080/actuator/health/custom`.

Similarly, to create a group that excludes the database indicators from the group and includes all the other indicators, you can define the following:

-
Properties
-
YAML

`management.endpoint.health.group.custom.exclude=db````
management:
 endpoint:
 health:
 group:
 custom:
 exclude: "db"
```
By default, startup will fail if a health group includes or excludes a health indicator that does not exist.
To disable this behavior set `management.endpoint.health.validate-group-membership` to `false`.

By default, groups inherit the same `StatusAggregator` and `HttpCodeStatusMapper` settings as the system health.
However, you can also define these on a per-group basis.
You can also override the `show-details` and `roles` properties if required:

-
Properties
-
YAML

```
management.endpoint.health.group.custom.show-details=when-authorized
management.endpoint.health.group.custom.roles=admin
management.endpoint.health.group.custom.status.order=fatal,up
management.endpoint.health.group.custom.status.http-mapping.fatal=500
management.endpoint.health.group.custom.status.http-mapping.out-of-service=500
```
```
management:
 endpoint:
 health:
 group:
 custom:
 show-details: "when-authorized"
 roles: "admin"
 status:
 order: "fatal,up"
 http-mapping:
 fatal: 500
 out-of-service: 500
```
| You can use `@Qualifier("groupname")`if you need to register custom`StatusAggregator`or`HttpCodeStatusMapper`beans for use with the group. |

A health group can also include/exclude a `CompositeHealthContributor`.
You can also include/exclude only a certain component of a `CompositeHealthContributor`.
This can be done using the fully qualified name of the component as follows:

```
management.endpoint.health.group.custom.include="test/primary"
management.endpoint.health.group.custom.exclude="test/primary/b"
```
In the example above, the `custom` group will include the `HealthContributor` with the name `primary` which is a component of the composite `test`.
Here, `primary` itself is a composite and the `HealthContributor` with the name `b` will be excluded from the `custom` group.

Health groups can be made available at an additional path on either the main or management port. This is useful in cloud environments such as Kubernetes, where it is quite common to use a separate management port for the actuator endpoints for security purposes. Having a separate port could lead to unreliable health checks because the main application might not work properly even if the health check is successful. The health group can be configured with an additional path as follows:

`management.endpoint.health.group.live.additional-path="server:/healthz"`This would make the `live` health group available on the main server port at `/healthz`.
The prefix is mandatory and must be either `server:` (represents the main server port) or `management:` (represents the management port, if configured.)
The path must be a single path segment.

### DataSource Health

The `DataSource` health indicator shows the health of both standard data sources and routing data source beans.
The health of a routing data source includes the health of each of its target data sources.
In the health endpoint’s response, each of a routing data source’s targets is named by using its routing key.
If you prefer not to include routing data sources in the indicator’s output, set `management.health.db.ignore-routing-data-sources` to `true`.

## Kubernetes Probes

Applications deployed on Kubernetes can provide information about their internal state with Container Probes. Depending on your Kubernetes configuration, the kubelet calls those probes and reacts to the result.

By default, Spring Boot manages your Application Availability state.
If deployed in a Kubernetes environment, actuator gathers the “Liveness” and “Readiness” information from the `ApplicationAvailability` interface and uses that information in dedicated health indicators: `LivenessStateHealthIndicator` and `ReadinessStateHealthIndicator`.
These indicators are shown on the global health endpoint (`"/actuator/health"`).
They are also exposed as separate HTTP Probes by using health groups: `"/actuator/health/liveness"` and `"/actuator/health/readiness"`.

You can then configure your Kubernetes infrastructure with the following endpoint information:

```
livenessProbe:
 httpGet:
 path: "/actuator/health/liveness"
 port: <actuator-port>
 failureThreshold: ...
 periodSeconds: ...
readinessProbe:
 httpGet:
 path: "/actuator/health/readiness"
 port: <actuator-port>
 failureThreshold: ...
 periodSeconds: ...
```
| `<actuator-port>`should be set to the port that the actuator endpoints are available on.
It could be the main web server port or a separate management port if the`"management.server.port"`property has been set. |

These health groups are automatically enabled.
You can disable them by using the `management.endpoint.health.probes.enabled` configuration property.

| If an application takes longer to start than the configured liveness period, Kubernetes mentions the `"startupProbe"`as a possible solution.
Generally speaking, the`"startupProbe"`is not necessarily needed here as the`"readinessProbe"`fails until all startup tasks are done.
This means your application will not receive traffic until it is ready.
However, if your application takes a long time to start, consider configuring a`"startupProbe"`that uses the liveness HTTP probe to make sure that Kubernetes won’t kill your application while it is in the process of starting.
See the section that describes how probes behave during the application lifecycle. |

If your Actuator endpoints are deployed on a separate management context, the endpoints do not use the same web infrastructure (port, connection pools, framework components) as the main application.
In this case, a probe check could be successful even if the main application does not work properly (for example, it cannot accept new connections).
For this reason, it is a good idea to make the `liveness` and `readiness` health groups available on the main server port.
This can be done by setting the following property:

`management.endpoint.health.probes.add-additional-paths=true`This would make the `liveness` group available at `/livez` and the `readiness` group available at `/readyz` on the main server port.
Paths can be customized using the `additional-path` property on each group, see health groups for details.

### Checking External State With Kubernetes Probes

Actuator configures the “liveness” and “readiness” probes as Health Groups. This means that all the health groups features are available for them. You can, for example, configure additional Health Indicators:

-
Properties
-
YAML

`management.endpoint.health.group.readiness.include=readinessState,customCheck````
management:
 endpoint:
 health:
 group:
 readiness:
 include: "readinessState,customCheck"
```
By default, Spring Boot does not add other health indicators to these groups.

The “liveness” probe should not depend on health checks for external systems. If the liveness state of an application is broken, Kubernetes tries to solve that problem by restarting the application instance. This means that if an external system (such as a database, a Web API, or an external cache) fails, Kubernetes might restart all application instances and create cascading failures.

As for the “readiness” probe, the choice of checking external systems must be made carefully by the application developers. For this reason, Spring Boot does not include any additional health checks in the readiness probe. If the readiness state of an application instance is unready, Kubernetes does not route traffic to that instance. Some external systems might not be shared by application instances, in which case they could be included in a readiness probe. Other external systems might not be essential to the application (the application could have circuit breakers and fallbacks), in which case they definitely should not be included. Unfortunately, an external system that is shared by all application instances is common, and you have to make a judgement call: Include it in the readiness probe and expect that the application is taken out of service when the external service is down or leave it out and deal with failures higher up the stack, perhaps by using a circuit breaker in the caller.

| If all instances of an application are unready, a Kubernetes Service with `type=ClusterIP`or`NodePort`does not accept any incoming connections.
There is no HTTP error response (503 and so on), since there is no connection.
A service with`type=LoadBalancer`might or might not accept connections, depending on the provider.
A service that has an explicit ingress also responds in a way that depends on the implementation — the ingress service itself has to decide how to handle the “connection refused” from downstream.
HTTP 503 is quite likely in the case of both load balancer and ingress. |

Also, if an application uses Kubernetes autoscaling, it may react differently to applications being taken out of the load-balancer, depending on its autoscaler configuration.

### Application Lifecycle and Probe States

An important aspect of the Kubernetes Probes support is its consistency with the application lifecycle.
There is a significant difference between the `AvailabilityState` (which is the in-memory, internal state of the application)
and the actual probe (which exposes that state).
Depending on the phase of application lifecycle, the probe might not be available.

Spring Boot publishes application events during startup and shutdown,
and probes can listen to such events and expose the `AvailabilityState` information.

The following tables show the `AvailabilityState` and the state of HTTP connectors at different stages.

When a Spring Boot application starts:

| Startup phase | LivenessState | ReadinessState | HTTP server | Notes |
|---|---|---|---|---|
| Starting |
 |
 | Not started | Kubernetes checks the "liveness" Probe and restarts the application if it takes too long. |
| Started |
 |
 | Refuses requests | The application context is refreshed. The application performs startup tasks and does not receive traffic yet. |
| Ready |
 |
 | Accepts requests | Startup tasks are finished. The application is receiving traffic. |

When a Spring Boot application shuts down:

| Shutdown phase | Liveness State | Readiness State | HTTP server | Notes |
|---|---|---|---|---|
| Running |
 |
 | Accepts requests | Shutdown has been requested. |
| Graceful shutdown |
 |
 | New requests are rejected | If enabled, graceful shutdown processes in-flight requests. HTTP probes also stop accepting traffic, so the availability states are not readily available externally. |
| Shutdown complete | N/A | N/A | Server is shut down | The application context is closed and the application is shut down. |

| See Kubernetes Container Lifecycle for more information about Kubernetes deployment.
In particular, it describes how to use the `preStop`hook to give your application time to shut down gracefully before Kubernetes kills it. |

## Application Information

Application information exposes various information collected from all `InfoContributor` beans defined in your `ApplicationContext`.
Spring Boot includes a number of auto-configured `InfoContributor` beans, and you can write your own.

### Auto-configured InfoContributors

When appropriate, Spring auto-configures the following `InfoContributor` beans:

| ID | Name | Description | Prerequisites |
|---|---|---|---|
|
 | Exposes build information. | A | |
|
 | Exposes any property from the | None. | |
|
 | Exposes git information. | A | |
|
 | Exposes Java runtime information. | None. | |
|
 | Exposes Operating System information. | None. | |
|
 | Exposes process information. | None. | |
|
 | Exposes SSL certificate information. | An SSL Bundle configured. |

Whether an individual contributor is enabled is controlled by its `management.info.<id>.enabled` property.
Different contributors have different defaults for this property, depending on their prerequisites and the nature of the information that they expose.

With no prerequisites to indicate that they should be enabled, the `env`, `java`, `os`, and `process` contributors are disabled by default. The `ssl` contributor has a prerequisite of having an SSL Bundle configured but it is disabled by default.
Each can be enabled by setting its `management.info.<id>.enabled` property to `true`.

The `build` and `git` info contributors are enabled by default.
Each can be disabled by setting its `management.info.<id>.enabled` property to `false`.
Alternatively, to disable every contributor that is usually enabled by default, set the `management.info.defaults.enabled` property to `false`.

### Custom Application Information

When the `env` contributor is enabled, you can customize the data exposed by the `info` endpoint by setting `info.*` Spring properties.
All `Environment` properties under the `info` key are automatically exposed.
For example, you could add the following settings to your `application.properties` file:

-
Properties
-
YAML

```
info.app.encoding=UTF-8
info.app.java.source=17
info.app.java.target=17
```
```
info:
 app:
 encoding: "UTF-8"
 java:
 source: "17"
 target: "17"
```
| Rather than hardcoding those values, you could also expand info properties at build time. Assuming you use Maven, you could rewrite the preceding example as follows:
 |

### Git Commit Information

Another useful feature of the `info` endpoint is its ability to publish information about the state of your `git` source code repository when the project was built.
If a `GitProperties` bean is available, you can use the `info` endpoint to expose these properties.

| A `GitProperties`bean is auto-configured if a`git.properties`file is available at the root of the classpath.
See Generate Git Information for more detail. |

By default, the endpoint exposes `git.branch`, `git.commit.id`, and `git.commit.time` properties, if present.
If you do not want any of these properties in the endpoint response, they need to be excluded from the `git.properties` file.
If you want to display the full git information (that is, the full content of `git.properties`), use the `management.info.git.mode` property, as follows:

-
Properties
-
YAML

`management.info.git.mode=full````
management:
 info:
 git:
 mode: "full"
```
To disable the git commit information from the `info` endpoint completely, set the `management.info.git.enabled` property to `false`, as follows:

-
Properties
-
YAML

`management.info.git.enabled=false````
management:
 info:
 git:
 enabled: false
```
### Build Information

If a `BuildProperties` bean is available, the `info` endpoint can also publish information about your build.
This happens if a `META-INF/build-info.properties` file is available in the classpath.

| The Maven and Gradle plugins can both generate that file. See Generate Build Information for more details. |

### Java Information

The `info` endpoint publishes information about your Java runtime environment, see `JavaInfo` for more details.

### OS Information

The `info` endpoint publishes information about your Operating System, see `OsInfo` for more details.

### Process Information

The `info` endpoint publishes information about your process, see `ProcessInfo` for more details.

### SSL Information

The `info` endpoint publishes information about your SSL certificates (that are configured through SSL Bundles), see `SslInfo` for more details.

### Writing Custom InfoContributors

To provide custom application information, you can register Spring beans that implement the `InfoContributor` interface.

The following example contributes an `example` entry with a single value:

-
Java
-
Kotlin

```
import java.util.Collections;
import org.springframework.boot.actuate.info.Info;
import org.springframework.boot.actuate.info.InfoContributor;
import org.springframework.stereotype.Component;
@Component
public class MyInfoContributor implements InfoContributor {
	@Override
	public void contribute(Info.Builder builder) {
 builder.withDetail("example", Collections.singletonMap("key", "value"));
	}
}
```
```
import org.springframework.boot.actuate.info.Info
import org.springframework.boot.actuate.info.InfoContributor
import org.springframework.stereotype.Component
import java.util.Collections
@Component
class MyInfoContributor : InfoContributor {
	override fun contribute(builder: Info.Builder) {
 builder.withDetail("example", Collections.singletonMap("key", "value"))
	}
}
```
If you reach the `info` endpoint, you should see a response that contains the following additional entry:

```
{
	"example": {
 "key" : "value"
	}
}
```
## Software Bill of Materials (SBOM)

The `sbom` endpoint exposes the Software Bill of Materials.
CycloneDX SBOMs can be auto-detected, but other formats can be manually configured, too.

The `sbom` actuator endpoint will then expose an SBOM called "application", which describes the contents of your application.

| To automatically generate a CycloneDX SBOM at project build time, please see the Generate a CycloneDX SBOM section. |

### Other SBOM formats

If you want to publish an SBOM in a different format, there are some configuration properties which you can use.

The configuration property `management.endpoint.sbom.application.location` sets the location for the application SBOM.
For example, setting this to `classpath:sbom.json` will use the contents of the `/sbom.json` resource on the classpath.

The media type for SBOMs in CycloneDX, SPDX and Syft format is detected automatically.
To override the auto-detected media type, use the configuration property `management.endpoint.sbom.application.media-type`.

### Additional SBOMs

The actuator endpoint can handle multiple SBOMs.
To add SBOMs, use the configuration property `management.endpoint.sbom.additional`, as shown in this example:

-
Properties
-
YAML

```
management.endpoint.sbom.additional.system.location=optional:file:/system.spdx.json
management.endpoint.sbom.additional.system.media-type=application/spdx+json
```
```
management:
 endpoint:
 sbom:
 additional:
 system:
 location: "optional:file:/system.spdx.json"
 media-type: "application/spdx+json"
```
This will add an SBOM called "system", which is stored in `/system.spdx.json`.
The `optional:` prefix can be used to prevent a startup failure if the file doesn’t exist.

# Monitoring and Management Over HTTP

If you are developing a web application, Spring Boot Actuator auto-configures all enabled endpoints to be exposed over HTTP.
The default convention is to use the `id` of the endpoint with a prefix of `/actuator` as the URL path.
For example, `health` is exposed as `/actuator/health`.

| Actuator is supported natively with Spring MVC, Spring WebFlux, and Jersey. If both Jersey and Spring MVC are available, Spring MVC is used. |

| Jackson is a required dependency in order to get the correct JSON responses as documented in the API documentation. Jackson 3 should be used for Spring MVC and Spring WebFlux. Jersey does not yet have a Jackson 3 module, so you will need to use Jackson 2. |

## Customizing the Management Endpoint Paths

Sometimes, it is useful to customize the prefix for the management endpoints.
For example, your application might already use `/actuator` for another purpose.
You can use the `management.endpoints.web.base-path` property to change the prefix for your management endpoint, as the following example shows:

-
Properties
-
YAML

`management.endpoints.web.base-path=/manage````
management:
 endpoints:
 web:
 base-path: "/manage"
```
The preceding `application.properties` example changes the endpoint from `/actuator/{id}` to `/manage/{id}` (for example, `/manage/info`).

| Unless the management port has been configured to expose endpoints by using a different HTTP port, `management.endpoints.web.base-path`is relative to`server.servlet.context-path`(for servlet web applications) or`spring.webflux.base-path`(for reactive web applications).
If`management.server.port`is configured,`management.endpoints.web.base-path`is relative to`management.server.base-path`. |

If you want to map endpoints to a different path, you can use the `management.endpoints.web.path-mapping` property.

The following example remaps `/actuator/health` to `/healthcheck`:

-
Properties
-
YAML

```
management.endpoints.web.base-path=/
management.endpoints.web.path-mapping.health=healthcheck
```
```
management:
 endpoints:
 web:
 base-path: "/"
 path-mapping:
 health: "healthcheck"
```
## Customizing the Management Server Port

Exposing management endpoints by using the default HTTP port is a sensible choice for cloud-based deployments. If, however, your application runs inside your own data center, you may prefer to expose endpoints by using a different HTTP port.

You can set the `management.server.port` property to change the HTTP port, as the following example shows:

-
Properties
-
YAML

`management.server.port=8081````
management:
 server:
 port: 8081
```
| On Cloud Foundry, by default, applications receive requests only on port 8080 for both HTTP and TCP routing. If you want to use a custom management port on Cloud Foundry, you need to explicitly set up the application’s routes to forward traffic to the custom port. |

## Configuring Management-specific SSL

When configured to use a custom port, you can also configure the management server with its own SSL by using the various `management.server.ssl.*` properties.
For example, doing so lets a management server be available over HTTP while the main application uses HTTPS, as the following property settings show:

-
Properties
-
YAML

```
server.port=8443
server.ssl.enabled=true
server.ssl.key-store=classpath:store.jks
server.ssl.key-password=secret
management.server.port=8080
management.server.ssl.enabled=false
```
```
server:
 port: 8443
 ssl:
 enabled: true
 key-store: "classpath:store.jks"
 key-password: "secret"
management:
 server:
 port: 8080
 ssl:
 enabled: false
```
Alternatively, both the main server and the management server can use SSL but with different key stores, as follows:

-
Properties
-
YAML

```
server.port=8443
server.ssl.enabled=true
server.ssl.key-store=classpath:main.jks
server.ssl.key-password=secret
management.server.port=8080
management.server.ssl.enabled=true
management.server.ssl.key-store=classpath:management.jks
management.server.ssl.key-password=secret
```
```
server:
 port: 8443
 ssl:
 enabled: true
 key-store: "classpath:main.jks"
 key-password: "secret"
management:
 server:
 port: 8080
 ssl:
 enabled: true
 key-store: "classpath:management.jks"
 key-password: "secret"
```
## Customizing the Management Server Address

You can customize the address on which the management endpoints are available by setting the `management.server.address` property.
Doing so can be useful if you want to listen only on an internal or ops-facing network or to listen only for connections from `localhost`.

| You can listen on a different address only when the port differs from the main server port. |

The following example `application.properties` does not allow remote management connections:

-
Properties
-
YAML

```
management.server.port=8081
management.server.address=127.0.0.1
```
```
management:
 server:
 port: 8081
 address: "127.0.0.1"
```
## Disabling HTTP Endpoints

If you do not want to expose endpoints over HTTP, you can set the management port to `-1`, as the following example shows:

-
Properties
-
YAML

`management.server.port=-1````
management:
 server:
 port: -1
```
You can also achieve this by using the `management.endpoints.web.exposure.exclude` property, as the following example shows:

-
Properties
-
YAML

`management.endpoints.web.exposure.exclude=*````
management:
 endpoints:
 web:
 exposure:
 exclude: "*"
```

# Monitoring and Management over JMX

Java Management Extensions (JMX) provide a standard mechanism to monitor and manage applications.
By default, this feature is not enabled.
You can turn it on by setting the `spring.jmx.enabled` configuration property to `true`.
Spring Boot exposes the most suitable `MBeanServer` as a bean with an ID of `mbeanServer`.
Any of your beans that are annotated with Spring JMX annotations (`@org.springframework.jmx.export.annotation.ManagedResource`, `@ManagedAttribute`, or `@ManagedOperation`) are exposed to it.

If your platform provides a standard `MBeanServer`, Spring Boot uses that and defaults to the VM `MBeanServer`, if necessary.
If all that fails, a new `MBeanServer` is created.

| `spring.jmx.enabled`affects only the management beans provided by Spring.
Enabling management beans provided by other libraries (for example Log4j2 or Quartz) is independent. |

See the `JmxAutoConfiguration` class for more details.

By default, Spring Boot also exposes management endpoints as JMX MBeans under the `org.springframework.boot` domain.
To take full control over endpoint registration in the JMX domain, consider registering your own `EndpointObjectNameFactory` implementation.

## Customizing MBean Names

The name of the MBean is usually generated from the `id` of the endpoint.
For example, the `health` endpoint is exposed as `org.springframework.boot:type=Endpoint,name=Health`.

If your application contains more than one Spring `ApplicationContext`, you may find that names clash.
To solve this problem, you can set the `spring.jmx.unique-names` property to `true` so that MBean names are always unique.

You can also customize the JMX domain under which endpoints are exposed.
The following settings show an example of doing so in `application.properties`:

-
Properties
-
YAML

```
spring.jmx.unique-names=true
management.endpoints.jmx.domain=com.example.myapp
```
```
spring:
 jmx:
 unique-names: true
management:
 endpoints:
 jmx:
 domain: "com.example.myapp"
```

# Observability

Observability is the ability to observe the internal state of a running system from the outside. It consists of the three pillars: logging, metrics and traces.

For metrics and traces, Spring Boot uses Micrometer Observation.
To create your own observations (which will lead to metrics and traces), you can inject an `ObservationRegistry`.

-
Java
-
Kotlin

```
import io.micrometer.observation.Observation;
import io.micrometer.observation.ObservationRegistry;
import org.springframework.stereotype.Component;
@Component
public class MyCustomObservation {
	private final ObservationRegistry observationRegistry;
	public MyCustomObservation(ObservationRegistry observationRegistry) {
 this.observationRegistry = observationRegistry;
	}
	public void doSomething() {
 Observation.createNotStarted("doSomething", this.observationRegistry)
 .lowCardinalityKeyValue("locale", "en-US")
 .highCardinalityKeyValue("userId", "42")
 .observe(() -> {
 // Execute business logic here
 });
	}
}
```
```
import io.micrometer.observation.Observation
import io.micrometer.observation.ObservationRegistry;
import org.springframework.stereotype.Component
@Component
class MyCustomObservation(private val observationRegistry: ObservationRegistry) {
	fun doSomething() {
 Observation.createNotStarted("doSomething", observationRegistry)
 .lowCardinalityKeyValue("locale", "en-US")
 .highCardinalityKeyValue("userId", "42")
 .observe {
 // Execute business logic here
 }
	}
}
```
| Low cardinality tags will be added to metrics and traces, while high cardinality tags will only be added to traces. |

Beans of type `ObservationPredicate`, `GlobalObservationConvention`, `ObservationFilter` and `ObservationHandler` will be automatically registered on the `ObservationRegistry`.
You can additionally register any number of `ObservationRegistryCustomizer` beans to further configure the registry.

| Observability for JDBC can be configured using a separate project. The Datasource Micrometer project provides a Spring Boot starter which automatically creates observations when JDBC operations are invoked. Read more about it in the reference documentation. |

| Observability for R2DBC is built into Spring Boot.
To enable it, add the `io.r2dbc:r2dbc-proxy`dependency to your project. |

## Context Propagation

Observability support relies on the Context Propagation library for forwarding the current observation across threads and reactive pipelines.
By default, `ThreadLocal` values are not automatically reinstated in reactive operators.
This behavior is controlled with the `spring.reactor.context-propagation` property, which can be set to `auto` to enable automatic propagation.

If you’re working with `@Async` methods and the `AsyncTaskExecutor` is auto-configured, you have to opt-in for context propagation using the `spring.task.execution.propagate-context` property.

If you are configuring the `AsyncTaskExecutor` yourself, then you need to register a `ContextPropagatingTaskDecorator` bean, as shown in the following example:

-
Java
-
Kotlin

```
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.task.support.ContextPropagatingTaskDecorator;
@Configuration(proxyBeanMethods = false)
class ContextPropagationConfiguration {
	@Bean
	ContextPropagatingTaskDecorator contextPropagatingTaskDecorator() {
 return new ContextPropagatingTaskDecorator();
	}
}
```
```
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.core.task.support.ContextPropagatingTaskDecorator
@Configuration(proxyBeanMethods = false)
class ContextPropagationConfiguration {
	@Bean
	fun contextPropagatingTaskDecorator(): ContextPropagatingTaskDecorator {
 return ContextPropagatingTaskDecorator()
	}
}
```
For more details about observations please see the Micrometer Observation documentation.

## Common Tags

Common tags are generally used for dimensional drill-down on the operating environment, such as host, instance, region, stack, and others. Common tags are applied to all observations as low cardinality tags and can be configured, as the following example shows:

-
Properties
-
YAML

```
management.observations.key-values.region=us-east-1
management.observations.key-values.stack=prod
```
```
management:
 observations:
 key-values:
 region: "us-east-1"
 stack: "prod"
```
The preceding example adds `region` and `stack` tags to all observations with a value of `us-east-1` and `prod`, respectively.

## Preventing Observations

If you’d like to prevent some observations from being reported, you can use the `management.observations.enable` properties:

-
Properties
-
YAML

```
management.observations.enable.denied.prefix=false
management.observations.enable.another.denied.prefix=false
```
```
management:
 observations:
 enable:
 denied:
 prefix: false
 another:
 denied:
 prefix: false
```
The preceding example will prevent all observations with a name starting with `denied.prefix` or `another.denied.prefix`.

| If you want to prevent Spring Security from reporting observations, set the property `management.observations.enable.spring.security`to`false`. |

If you need greater control over the prevention of observations, you can register beans of type `ObservationPredicate`.
Observations are only reported if all the `ObservationPredicate` beans return `true` for that observation.

-
Java
-
Kotlin

```
import io.micrometer.observation.Observation.Context;
import io.micrometer.observation.ObservationPredicate;
import org.springframework.stereotype.Component;
@Component
class MyObservationPredicate implements ObservationPredicate {
	@Override
	public boolean test(String name, Context context) {
 return !name.contains("denied");
	}
}
```
```
import io.micrometer.observation.Observation.Context
import io.micrometer.observation.ObservationPredicate
import org.springframework.stereotype.Component
@Component
class MyObservationPredicate : ObservationPredicate {
	override fun test(name: String, context: Context): Boolean {
 return !name.contains("denied")
	}
}
```
The preceding example will prevent all observations whose name contains "denied".

## Micrometer Observation Annotations support

To enable scanning of observability annotations like `@Observed`, `@Timed`, `@Counted`, `@MeterTag` and `@NewSpan`, set the `management.observations.annotations.enabled` property to `true`.
A dependency on `org.aspectj:aspectjweaver`, which is part of `spring-boot-starter-aspectj`, is also required.
This feature is supported by Micrometer directly.
Please refer to the Micrometer, Micrometer Observation and Micrometer Tracing reference docs.

| When you annotate methods or classes which are already instrumented (for example, Spring Data repositories or Spring MVC controllers), you will get duplicate observations.
In that case you can either disable the automatic instrumentation using properties or an `ObservationPredicate`and rely on your annotations, or you can remove your annotations. |

## OpenTelemetry Support

| There are several ways to support OpenTelemetry in your application. You can use the OpenTelemetry Java Agent or the OpenTelemetry Spring Boot Starter, which are supported by the OTel community; the metrics and traces use the semantic conventions defined by OTel libraries. This documentation describes OpenTelemetry as officially supported by the Spring team, using Micrometer and the OTLP exporter; the metrics and traces use the semantic conventions described in the Spring projects documentation, such as Spring Framework. |

Spring Boot’s actuator module includes basic support for OpenTelemetry.

It provides a bean of type `OpenTelemetry`, and if there are beans of type `SdkTracerProvider`, `ContextPropagators`, `SdkLoggerProvider` or `SdkMeterProvider` in the application context, they automatically get registered.
Additionally, it provides a `Resource` bean.
The attributes of the auto-configured `Resource` can be configured via the `management.opentelemetry.resource-attributes` configuration property.
Auto-configured attributes will be merged with attributes from the `OTEL_RESOURCE_ATTRIBUTES` and `OTEL_SERVICE_NAME` environment variables, with attributes configured through the configuration property taking precedence over those from the environment variables.

If you have defined your own `Resource` bean, this will no longer be the case.

| Spring Boot does not provide automatic exporting of OpenTelemetry metrics or logs. Exporting OpenTelemetry traces is only auto-configured when used together with Micrometer Tracing. |

### Disabling OpenTelemetry

The OpenTelemetry support can be disabled by setting the `management.opentelemetry.enabled` property to `false`.
This behaves similarly to the `OTEL_SDK_DISABLED` environment variable (but negated): when the SDK is disabled, metrics, traces, and logging will use no-op implementations.
Context propagators are not affected and continue to function normally.

| Keep in mind that Spring Boot doesn’t use OpenTelemetry’s metrics functionality, so metrics might still be enabled even when disabling OpenTelemetry. |

### Environment variables

Spring Boot supports a subset of the OpenTelemetry SDK environment variables.
These environment variables are automatically mapped to Spring Boot configuration properties at startup.
When an environment variable has a signal-specific variant (for example, `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`) and a general variant (`OTEL_EXPORTER_OTLP_ENDPOINT`), the signal-specific variant takes precedence and is used as-is.
When the general variant is used as a fallback, the signal-specific path (`v1/traces`, `v1/metrics`, or `v1/logs`) is appended to it.
For example, setting `OTEL_EXPORTER_OTLP_ENDPOINT=http://collector:4318` results in `collector:4318/v1/traces` for traces.

This mapping can be disabled by setting `management.opentelemetry.map-environment-variables` to `false`.

#### General

| Environment Variable | Spring Boot Property |
|---|---|
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |

#### Metrics Exporter

| Environment Variable | Spring Boot Property |
|---|---|
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |

#### Traces Exporter

| Environment Variable | Spring Boot Property |
|---|---|
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |

#### Logs Exporter

| Environment Variable | Spring Boot Property |
|---|---|
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |
|
 |
 |

#### Resource

The following environment variables configure the OpenTelemetry resource:

| The `OTEL_RESOURCE_ATTRIBUTES`environment variable consists of a list of key-value pairs.
For example:`key1=value1,key2=value2,key3=spring%20boot`.
All attribute values are treated as strings, and any characters outside the baggage-octet range must bepercent-encoded. |

#### Unsupported Environment Variables

Other environment variables as described in the OpenTelemetry documentation are not supported.

If you want all environment variables specified by OpenTelemetry’s SDK to be effective, you have to supply your own `OpenTelemetry` bean.

| Doing this will switch off Spring Boot’s OpenTelemetry auto-configuration and may break the built-in observability functionality. |

First, add a dependency to `io.opentelemetry:opentelemetry-sdk-extension-autoconfigure` to get OpenTelemetry’s zero-code SDK autoconfigure module, then add this configuration:

```
import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.sdk.autoconfigure.AutoConfiguredOpenTelemetrySdk;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
class AutoConfiguredOpenTelemetrySdkConfiguration {
	@Bean
	OpenTelemetry autoConfiguredOpenTelemetrySdk() {
 return AutoConfiguredOpenTelemetrySdk.initialize().getOpenTelemetrySdk();
	}
}
```
### Logging

The `OpenTelemetryLoggingAutoConfiguration` configures OpenTelemetry’s `SdkLoggerProvider`.
Exporting logs via OTLP is supported through the `OtlpLoggingAutoConfiguration`, which enables OTLP log exporting over HTTP or gRPC.

| If you need to apply advanced customizations to OTLP log record exporters, consider registering `OtlpHttpLogRecordExporterBuilderCustomizer`or`OtlpGrpcLogRecordExporterBuilderCustomizer`beans.
These will be invoked before the creation of the`OtlpHttpLogRecordExporter`or`OtlpGrpcLogRecordExporter`.
The customizers take precedence over anything applied by the auto-configuration. |

However, while there is a `SdkLoggerProvider` bean, Spring Boot doesn’t support bridging logs to this bean out of the box.
This can be done with 3rd-party log bridges, as described in the Logging with OpenTelemetry section.

### Metrics

The choice of metrics in the Spring portfolio is Micrometer, which means that metrics are not collected and exported through the OpenTelemetry’s `SdkMeterProvider`.
Spring Boot doesn’t provide a `SdkMeterProvider` bean.

However, Micrometer metrics can be exported via OTLP to any OpenTelemetry capable backend using the `OtlpMeterRegistry`, as described in the Metrics with OTLP section.

| Micrometer’s OTLP registry doesn’t use the `Resource`bean, but setting`OTEL_RESOURCE_ATTRIBUTES`,`OTEL_SERVICE_NAME`or`management.opentelemetry.resource-attributes`works. |

#### Metrics via the OpenTelemetry API and SDK

If you or a dependency you include make use of OpenTelemetry’s `MeterProvider`, those metrics are not exported.

We strongly recommend that you report your metrics with Micrometer.
If a dependency you include uses OpenTelemetry’s `MeterProvider`, you can include this configuration in your application to configure a `MeterProvider` bean, which you then have to wire into your dependency:

```
import java.time.Duration;
import io.opentelemetry.exporter.otlp.http.metrics.OtlpHttpMetricExporter;
import io.opentelemetry.sdk.metrics.SdkMeterProvider;
import io.opentelemetry.sdk.metrics.export.MetricExporter;
import io.opentelemetry.sdk.metrics.export.MetricReader;
import io.opentelemetry.sdk.metrics.export.PeriodicMetricReader;
import io.opentelemetry.sdk.resources.Resource;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
class OpenTelemetryMetricsConfiguration {
	@Bean
	OtlpHttpMetricExporter metricExporter() {
 String endpoint = "http://localhost:4318/v1/metrics";
 return OtlpHttpMetricExporter.builder().setEndpoint(endpoint).build();
	}
	@Bean
	PeriodicMetricReader metricReader(MetricExporter exporter) {
 Duration interval = Duration.ofMinutes(1);
 return PeriodicMetricReader.builder(exporter).setInterval(interval).build();
	}
	@Bean
	SdkMeterProvider meterProvider(Resource resource, MetricReader metricReader) {
 return SdkMeterProvider.builder().registerMetricReader(metricReader).setResource(resource).build();
	}
}
```
This configuration also enables metrics export via OTLP over HTTP.

### Tracing

If Micrometer tracing is used, the `OpenTelemetryTracingAutoConfiguration` configures OpenTelemetry’s `SdkTracerProvider`.
Exporting traces through OTLP is enabled by the `OtlpTracingAutoConfiguration`, which supports exporting traces with OTLP over HTTP or gRPC.

We strongly recommend using the Micrometer Observation or Tracing API instead of using the OpenTelemetry API directly.

# Loggers

Spring Boot Actuator includes the ability to view and configure the log levels of your application at runtime. You can view either the entire list or an individual logger’s configuration, which is made up of both the explicitly configured logging level as well as the effective logging level given to it by the logging framework. These levels can be one of:

-
`TRACE`
-
`DEBUG`
-
`INFO`
-
`WARN`
-
`ERROR`
-
`FATAL`
-
`OFF`
-
`null`

`null` indicates that there is no explicit configuration.

## Configure a Logger

To configure a given logger, `POST` a partial entity to the resource’s URI, as the following example shows:

```
{
	"configuredLevel": "DEBUG"
}
```
| To “reset” the specific level of the logger (and use the default configuration instead), you can pass a value of `null`as the`configuredLevel`. |

## OpenTelemetry

By default, logging via OpenTelemetry is not configured. You have to provide the location of the OpenTelemetry logs endpoint to configure it:

-
Properties
-
YAML

`management.opentelemetry.logging.export.otlp.endpoint=https://otlp.example.com:4318/v1/logs````
management:
 opentelemetry:
 logging:
 export:
 otlp:
 endpoint: "https://otlp.example.com:4318/v1/logs"
```
The `management.opentelemetry.logging.export.*` configuration properties can be used to configure the `BatchLogRecordProcessor`.
For example, to change the export interval to 15 seconds:

-
Properties
-
YAML

`management.opentelemetry.logging.export.schedule-delay=15s````
management:
 opentelemetry:
 logging:
 export:
 schedule-delay: "15s"
```
The `management.opentelemetry.logging.limits.*` configuration properties can be used to configure log record limits.
For example, to limit the number of attributes per log record to 64 and the maximum attribute value length to 256 characters:

-
Properties
-
YAML

```
management.opentelemetry.logging.limits.max-attributes=64
management.opentelemetry.logging.limits.max-attribute-value-length=256
```
```
management:
 opentelemetry:
 logging:
 limits:
 max-attributes: 64
 max-attribute-value-length: 256
```
If you need full control, you can register a custom `LogLimits` bean.

| The OpenTelemetry Logback appender and Log4j appender are not part of Spring Boot. For more details, see the OpenTelemetry Logback appender or the OpenTelemetry Log4j2 appender in the OpenTelemetry Java instrumentation GitHub repository. |

| You have to configure the appender in your `logback-spring.xml`or`log4j2-spring.xml`configuration to get OpenTelemetry logging working. |

The `OpenTelemetryAppender` for both Logback and Log4j requires access to an `OpenTelemetry` instance to function properly.
This instance must be set programmatically during application startup, which can be done like this:

-
Java
-
Kotlin

```
import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.instrumentation.logback.appender.v1_0.OpenTelemetryAppender;
import org.springframework.beans.factory.InitializingBean;
import org.springframework.stereotype.Component;
@Component
class OpenTelemetryAppenderInitializer implements InitializingBean {
	private final OpenTelemetry openTelemetry;
	OpenTelemetryAppenderInitializer(OpenTelemetry openTelemetry) {
 this.openTelemetry = openTelemetry;
	}
	@Override
	public void afterPropertiesSet() {
 OpenTelemetryAppender.install(this.openTelemetry);
	}
}
```
```
import io.opentelemetry.api.OpenTelemetry
import io.opentelemetry.instrumentation.logback.appender.v1_0.OpenTelemetryAppender
import org.springframework.beans.factory.InitializingBean
import org.springframework.stereotype.Component
@Component
class OpenTelemetryAppenderInitializer(
	private val openTelemetry: OpenTelemetry
) : InitializingBean {
	override fun afterPropertiesSet() = OpenTelemetryAppender.install(openTelemetry)
}
```

# Metrics

Spring Boot Actuator provides dependency management and auto-configuration for Micrometer, an application metrics facade that supports numerous monitoring systems, including:

| To learn more about Micrometer’s capabilities, see its reference documentation, in particular the concepts section. |

## Getting Started

Spring Boot auto-configures a composite `MeterRegistry` and adds a registry to the composite for each of the supported implementations that it finds on the classpath.
Having a dependency on `micrometer-registry-{system}` in your runtime classpath is enough for Spring Boot to configure the registry.

Most registries share common features. For instance, you can disable a particular registry even if the Micrometer registry implementation is on the classpath. The following example disables Datadog:

-
Properties
-
YAML

`management.datadog.metrics.export.enabled=false````
management:
 datadog:
 metrics:
 export:
 enabled: false
```
You can also disable all registries unless stated otherwise by the registry-specific property, as the following example shows:

-
Properties
-
YAML

`management.defaults.metrics.export.enabled=false````
management:
 defaults:
 metrics:
 export:
 enabled: false
```
Spring Boot also adds any auto-configured registries to the global static composite registry on the `Metrics` class, unless you explicitly tell it not to:

-
Properties
-
YAML

`management.metrics.use-global-registry=false````
management:
 metrics:
 use-global-registry: false
```
You can register any number of `MeterRegistryCustomizer` beans to further configure the registry, such as applying common tags, before any meters are registered with the registry:

-
Java
-
Kotlin

```
import io.micrometer.core.instrument.MeterRegistry;
import org.springframework.boot.micrometer.metrics.autoconfigure.MeterRegistryCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyMeterRegistryConfiguration {
	@Bean
	public MeterRegistryCustomizer<MeterRegistry> metricsCommonTags() {
 return (registry) -> registry.config().commonTags("region", "us-east-1");
	}
}
```
```
import io.micrometer.core.instrument.MeterRegistry
import org.springframework.boot.micrometer.metrics.autoconfigure.MeterRegistryCustomizer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
@Configuration(proxyBeanMethods = false)
class MyMeterRegistryConfiguration {
	@Bean
	fun metricsCommonTags(): MeterRegistryCustomizer<MeterRegistry> {
 return MeterRegistryCustomizer { registry ->
 registry.config().commonTags("region", "us-east-1")
 }
	}
}
```
You can apply customizations to particular registry implementations by being more specific about the generic type:

-
Java
-
Kotlin

```
import io.micrometer.core.instrument.Meter;
import io.micrometer.core.instrument.config.NamingConvention;
import io.micrometer.graphite.GraphiteMeterRegistry;
import org.springframework.boot.micrometer.metrics.autoconfigure.MeterRegistryCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyMeterRegistryConfiguration {
	@Bean
	public MeterRegistryCustomizer<GraphiteMeterRegistry> graphiteMetricsNamingConvention() {
 return (registry) -> registry.config().namingConvention(this::name);
	}
	private String name(String name, Meter.Type type, String baseUnit) {
 return ...
	}
}
```
```
import io.micrometer.core.instrument.Meter
import io.micrometer.core.instrument.config.NamingConvention
import io.micrometer.graphite.GraphiteMeterRegistry
import org.springframework.boot.micrometer.metrics.autoconfigure.MeterRegistryCustomizer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
@Configuration(proxyBeanMethods = false)
class MyMeterRegistryConfiguration {
	@Bean
	fun graphiteMetricsNamingConvention(): MeterRegistryCustomizer<GraphiteMeterRegistry> {
 return MeterRegistryCustomizer { registry: GraphiteMeterRegistry ->
 registry.config().namingConvention(this::name)
 }
	}
	private fun name(name: String, type: Meter.Type, baseUnit: String?): String {
 return ...
	}
}
```
Spring Boot also configures built-in instrumentation that you can control through configuration or dedicated annotation markers.

## Supported Monitoring Systems

This section briefly describes each of the supported monitoring systems.

### AppOptics

By default, the AppOptics registry periodically pushes metrics to `api.appoptics.com/v1/measurements`.
To export metrics to SaaS AppOptics, your API token must be provided:

-
Properties
-
YAML

`management.appoptics.metrics.export.api-token=YOUR_TOKEN````
management:
 appoptics:
 metrics:
 export:
 api-token: "YOUR_TOKEN"
```
### Atlas

By default, metrics are exported to Atlas running on your local machine. You can provide the location of the Atlas server:

-
Properties
-
YAML

`management.atlas.metrics.export.uri=https://atlas.example.com:7101/api/v1/publish````
management:
 atlas:
 metrics:
 export:
 uri: "https://atlas.example.com:7101/api/v1/publish"
```
### Datadog

A Datadog registry periodically pushes metrics to datadoghq. To export metrics to Datadog, you must provide your API key:

-
Properties
-
YAML

`management.datadog.metrics.export.api-key=YOUR_KEY````
management:
 datadog:
 metrics:
 export:
 api-key: "YOUR_KEY"
```
If you additionally provide an application key (optional), then metadata such as meter descriptions, types, and base units will also be exported:

-
Properties
-
YAML

```
management.datadog.metrics.export.api-key=YOUR_API_KEY
management.datadog.metrics.export.application-key=YOUR_APPLICATION_KEY
```
```
management:
 datadog:
 metrics:
 export:
 api-key: "YOUR_API_KEY"
 application-key: "YOUR_APPLICATION_KEY"
```
By default, metrics are sent to the Datadog US site (`api.datadoghq.com`).
If your Datadog project is hosted on one of the other sites, or you need to send metrics through a proxy, configure the URI accordingly:

-
Properties
-
YAML

`management.datadog.metrics.export.uri=https://api.datadoghq.eu````
management:
 datadog:
 metrics:
 export:
 uri: "https://api.datadoghq.eu"
```
You can also change the interval at which metrics are sent to Datadog:

-
Properties
-
YAML

`management.datadog.metrics.export.step=30s````
management:
 datadog:
 metrics:
 export:
 step: "30s"
```
### Dynatrace

Dynatrace offers two metrics ingest APIs, both of which are implemented for Micrometer. You can find the Dynatrace documentation on Micrometer metrics ingest here.

Configuration properties in the `v1` namespace apply only when exporting to the Timeseries v1 API.
Support for the V1 API is deprecated.

Configuration properties in the `v2` namespace apply only when exporting to the Metrics v2 API.

Note that this integration can export only to either the `v1` or `v2` version of the API at a time, with `v2` being strongly recommended due to the deprecation of the v1 API.
If the `device-id` (required for v1 but not used in v2) is set in the `v1` namespace, metrics are exported to the `v1` endpoint.
Otherwise, `v2` is assumed.

#### v2 API

You can use the v2 API in two ways.

##### Auto-configuration

Dynatrace auto-configuration is available for hosts that are monitored by the OneAgent or by the Dynatrace Operator for Kubernetes.

**Local OneAgent:** If a OneAgent is running on the host, metrics are automatically exported to the local OneAgent ingest endpoint.
The ingest endpoint forwards the metrics to the Dynatrace backend.

**Dynatrace Kubernetes Operator:** When running in Kubernetes with the Dynatrace Operator installed, the registry will automatically pick up your endpoint URI and API token from the operator instead.

This is the default behavior and requires no special setup beyond a dependency on `io.micrometer:micrometer-registry-dynatrace`.

##### Manual Configuration

If no auto-configuration is available, the endpoint of the Metrics v2 API and an API token are required.
The API token must have the “Ingest metrics” (`metrics.ingest`) permission set.
We recommend limiting the scope of the token to this one permission.
You must ensure that the endpoint URI contains the path (for example, `/api/v2/metrics/ingest`):

The URL of the Metrics API v2 ingest endpoint is different according to your deployment option:

-
SaaS: `https://{your-environment-id}.live.dynatrace.com/api/v2/metrics/ingest`
-
Managed deployments: `https://{your-domain}/e/{your-environment-id}/api/v2/metrics/ingest`

The example below configures metrics export using the `example` environment id:

-
Properties
-
YAML

```
management.dynatrace.metrics.export.uri=https://example.live.dynatrace.com/api/v2/metrics/ingest
management.dynatrace.metrics.export.api-token=YOUR_TOKEN
```
```
management:
 dynatrace:
 metrics:
 export:
 uri: "https://example.live.dynatrace.com/api/v2/metrics/ingest"
 api-token: "YOUR_TOKEN"
```
When using the Dynatrace v2 API, the following optional features are available (more details can be found in the Dynatrace documentation):

-
Metric key prefix: Sets a prefix that is prepended to all exported metric keys.
-
Enrich with Dynatrace metadata: If a OneAgent or Dynatrace operator is running, enrich metrics with additional metadata (for example, about the host, process, or pod).
-
Default dimensions: Specify key-value pairs that are added to all exported metrics. If tags with the same key are specified with Micrometer, they overwrite the default dimensions.
-
Use Dynatrace Summary instruments: In some cases the Micrometer Dynatrace registry created metrics that were rejected. In Micrometer 1.9.x, this was fixed by introducing Dynatrace-specific summary instruments. Setting this toggle to `false`forces Micrometer to fall back to the behavior that was the default before 1.9.x. It should only be used when encountering problems while migrating from Micrometer 1.8.x to 1.9.x.
-
Export meter metadata: Starting from Micrometer 1.12.0, the Dynatrace exporter will also export meter metadata, such as unit and description by default. Use the `export-meter-metadata`toggle to turn this feature off.

It is possible to not specify a URI and API token, as shown in the following example. In this scenario, the automatically configured endpoint is used:

-
Properties
-
YAML

```
management.dynatrace.metrics.export.v2.metric-key-prefix=your.key.prefix
management.dynatrace.metrics.export.v2.enrich-with-dynatrace-metadata=true
management.dynatrace.metrics.export.v2.default-dimensions.key1=value1
management.dynatrace.metrics.export.v2.default-dimensions.key2=value2
management.dynatrace.metrics.export.v2.use-dynatrace-summary-instruments=true
management.dynatrace.metrics.export.v2.export-meter-metadata=true
```
```
management:
 dynatrace:
 metrics:
 export:
 # Specify uri and api-token here if not using the local OneAgent endpoint.
 v2:
 metric-key-prefix: "your.key.prefix"
 enrich-with-dynatrace-metadata: true
 default-dimensions:
 key1: "value1"
 key2: "value2"
 use-dynatrace-summary-instruments: true # (default: true)
 export-meter-metadata: true # (default: true)
```
#### v1 API (Deprecated)

The Dynatrace v1 API metrics registry pushes metrics to the configured URI periodically by using the Timeseries v1 API.
For backwards-compatibility with existing setups, when `device-id` is set (required for v1, but not used in v2), metrics are exported to the Timeseries v1 endpoint.
To export metrics to Dynatrace, your API token, device ID, and URI must be provided:

-
Properties
-
YAML

```
management.dynatrace.metrics.export.uri=https://{your-environment-id}.live.dynatrace.com
management.dynatrace.metrics.export.api-token=YOUR_TOKEN
management.dynatrace.metrics.export.v1.device-id=YOUR_DEVICE_ID
```
```
management:
 dynatrace:
 metrics:
 export:
 uri: "https://{your-environment-id}.live.dynatrace.com"
 api-token: "YOUR_TOKEN"
 v1:
 device-id: "YOUR_DEVICE_ID"
```
For the v1 API, you must specify the base environment URI without a path, as the v1 endpoint path is added automatically.

#### Version-independent Settings

In addition to the API endpoint and token, you can also change the interval at which metrics are sent to Dynatrace.
The default export interval is `60s`.
The following example sets the export interval to 30 seconds:

-
Properties
-
YAML

`management.dynatrace.metrics.export.step=30s````
management:
 dynatrace:
 metrics:
 export:
 step: "30s"
```
You can find more information on how to set up the Dynatrace exporter for Micrometer in the Micrometer documentation and the Dynatrace documentation.

### Elastic

By default, metrics are exported to Elastic running on your local machine. You can provide the location of the Elastic server to use by using the following property:

-
Properties
-
YAML

`management.elastic.metrics.export.host=https://elastic.example.com:8086````
management:
 elastic:
 metrics:
 export:
 host: "https://elastic.example.com:8086"
```
### Ganglia

By default, metrics are exported to Ganglia running on your local machine. You can provide the Ganglia server host and port, as the following example shows:

-
Properties
-
YAML

```
management.ganglia.metrics.export.host=ganglia.example.com
management.ganglia.metrics.export.port=9649
```
```
management:
 ganglia:
 metrics:
 export:
 host: "ganglia.example.com"
 port: 9649
```
### Graphite

By default, metrics are exported to Graphite running on your local machine. You can provide the Graphite server host and port, as the following example shows:

-
Properties
-
YAML

```
management.graphite.metrics.export.host=graphite.example.com
management.graphite.metrics.export.port=9004
```
```
management:
 graphite:
 metrics:
 export:
 host: "graphite.example.com"
 port: 9004
```
Micrometer provides a default `HierarchicalNameMapper` that governs how a dimensional meter ID is mapped to flat hierarchical names.

| To take control over this behavior, define your
 |

### Humio

By default, the Humio registry periodically pushes metrics to cloud.humio.com. To export metrics to SaaS Humio, you must provide your API token:

-
Properties
-
YAML

`management.humio.metrics.export.api-token=YOUR_TOKEN````
management:
 humio:
 metrics:
 export:
 api-token: "YOUR_TOKEN"
```
You should also configure one or more tags to identify the data source to which metrics are pushed:

-
Properties
-
YAML

```
management.humio.metrics.export.tags.alpha=a
management.humio.metrics.export.tags.bravo=b
```
```
management:
 humio:
 metrics:
 export:
 tags:
 alpha: "a"
 bravo: "b"
```
### Influx

By default, metrics are exported to an Influx v1 instance running on your local machine with the default configuration.
To export metrics to InfluxDB v2, configure the `org`, `bucket`, and authentication `token` for writing metrics.
You can provide the location of the Influx server to use by using:

-
Properties
-
YAML

`management.influx.metrics.export.uri=https://influx.example.com:8086````
management:
 influx:
 metrics:
 export:
 uri: "https://influx.example.com:8086"
```
### JMX

Micrometer provides a hierarchical mapping to JMX, primarily as a cheap and portable way to view metrics locally.
By default, metrics are exported to the `metrics` JMX domain.
You can provide the domain to use by using:

-
Properties
-
YAML

`management.jmx.metrics.export.domain=com.example.app.metrics````
management:
 jmx:
 metrics:
 export:
 domain: "com.example.app.metrics"
```
Micrometer provides a default `HierarchicalNameMapper` that governs how a dimensional meter ID is mapped to flat hierarchical names.

| To take control over this behavior, define your
 |

### KairosDB

By default, metrics are exported to KairosDB running on your local machine. You can provide the location of the KairosDB server to use by using:

-
Properties
-
YAML

`management.kairos.metrics.export.uri=https://kairosdb.example.com:8080/api/v1/datapoints````
management:
 kairos:
 metrics:
 export:
 uri: "https://kairosdb.example.com:8080/api/v1/datapoints"
```
### New Relic

A New Relic registry periodically pushes metrics to New Relic. To export metrics to New Relic, you must provide your API key and account ID:

-
Properties
-
YAML

```
management.newrelic.metrics.export.api-key=YOUR_KEY
management.newrelic.metrics.export.account-id=YOUR_ACCOUNT_ID
```
```
management:
 newrelic:
 metrics:
 export:
 api-key: "YOUR_KEY"
 account-id: "YOUR_ACCOUNT_ID"
```
You can also change the interval at which metrics are sent to New Relic:

-
Properties
-
YAML

`management.newrelic.metrics.export.step=30s````
management:
 newrelic:
 metrics:
 export:
 step: "30s"
```
By default, metrics are published through REST calls, but you can also use the Java Agent API if you have it on the classpath:

-
Properties
-
YAML

`management.newrelic.metrics.export.client-provider-type=insights-agent````
management:
 newrelic:
 metrics:
 export:
 client-provider-type: "insights-agent"
```
Finally, you can take full control by defining your own `NewRelicClientProvider` bean.

### OTLP

By default, metrics are exported over the OpenTelemetry protocol (OTLP) to a consumer running on your local machine.
To export to another location, provide the location of the OTLP metrics endpoint using `management.otlp.metrics.export.url`:

-
Properties
-
YAML

`management.otlp.metrics.export.url=https://otlp.example.com:4318/v1/metrics````
management:
 otlp:
 metrics:
 export:
 url: "https://otlp.example.com:4318/v1/metrics"
```
Custom headers, for example for authentication, can also be provided using `management.otlp.metrics.export.headers.*` properties.

If an `OtlpMetricsSender` bean is available, it will be configured on the `OtlpMeterRegistry` that Spring Boot auto-configures.

OTLP Exemplars are also supported.
To enable this feature, an `ExemplarContextProvider` bean should be present.
If you use Micrometer Tracing, this will be auto-configured for you.
By default, only sampled traces are included as exemplars.
You can control this behavior using the `management.tracing.exemplars.include` property.

### Prometheus

Prometheus expects to scrape or poll individual application instances for metrics.
Spring Boot provides an actuator endpoint at `/actuator/prometheus` to present a Prometheus scrape with the appropriate format.

| By default, the endpoint is not available and must be exposed. See exposing endpoints for more details. |

The following example `scrape_config` adds to `prometheus.yml`:

```
scrape_configs:
- job_name: "spring"
 metrics_path: "/actuator/prometheus"
 static_configs:
 - targets: ["HOST:PORT"]
```
Prometheus Exemplars are also supported.
To enable this feature, a `SpanContext` bean should be present.
If you’re using the deprecated Prometheus simpleclient support and want to enable that feature, a `SpanContextSupplier` bean should be present.
If you use Micrometer Tracing, this will be auto-configured for you, but you can always create your own if you want.
By default, only sampled traces are included as exemplars.
You can control this behavior using the `management.tracing.exemplars.include` property.
The value `all` is not supported with Prometheus.
Please check the Prometheus Docs, since this feature needs to be explicitly enabled on Prometheus' side, and it is only supported using the OpenMetrics format.

For ephemeral or batch jobs that may not exist long enough to be scraped, you can use Prometheus Pushgateway support to expose the metrics to Prometheus.

To enable Prometheus Pushgateway support, add the following dependency to your project:

```
<dependency>
	<groupId>io.prometheus</groupId>
	<artifactId>prometheus-metrics-exporter-pushgateway</artifactId>
</dependency>
```
When the Prometheus Pushgateway dependency is present on the classpath and the `management.prometheus.metrics.export.pushgateway.enabled` property is set to `true`, a `PrometheusPushGatewayManager` bean is auto-configured.
This manages the pushing of metrics to a Prometheus Pushgateway.

You can tune the `PrometheusPushGatewayManager` by using properties under `management.prometheus.metrics.export.pushgateway`.
For advanced configuration, you can also provide your own `PrometheusPushGatewayManager` bean.

### Simple

Micrometer ships with a simple, in-memory backend that is automatically used as a fallback if no other registry is configured. This lets you see what metrics are collected in the metrics endpoint.

The in-memory backend disables itself as soon as you use any other available backend. You can also disable it explicitly:

-
Properties
-
YAML

`management.simple.metrics.export.enabled=false````
management:
 simple:
 metrics:
 export:
 enabled: false
```
### Stackdriver

The Stackdriver registry periodically pushes metrics to Stackdriver. To export metrics to SaaS Stackdriver, you must provide your Google Cloud project ID:

-
Properties
-
YAML

`management.stackdriver.metrics.export.project-id=my-project````
management:
 stackdriver:
 metrics:
 export:
 project-id: "my-project"
```
You can also change the interval at which metrics are sent to Stackdriver:

-
Properties
-
YAML

`management.stackdriver.metrics.export.step=30s````
management:
 stackdriver:
 metrics:
 export:
 step: "30s"
```
### StatsD

The StatsD registry eagerly pushes metrics over UDP to a StatsD agent. By default, metrics are exported to a StatsD agent running on your local machine. You can provide the StatsD agent host, port, and protocol to use by using:

-
Properties
-
YAML

```
management.statsd.metrics.export.host=statsd.example.com
management.statsd.metrics.export.port=9125
management.statsd.metrics.export.protocol=udp
```
```
management:
 statsd:
 metrics:
 export:
 host: "statsd.example.com"
 port: 9125
 protocol: "udp"
```
You can also change the StatsD line protocol to use (it defaults to Datadog):

-
Properties
-
YAML

`management.statsd.metrics.export.flavor=etsy````
management:
 statsd:
 metrics:
 export:
 flavor: "etsy"
```
## Supported Metrics and Meters

Spring Boot provides automatic meter registration for a wide variety of technologies. In most situations, the defaults provide sensible metrics that can be published to any of the supported monitoring systems.

### JVM Metrics

Auto-configuration enables JVM Metrics by using core Micrometer classes.
JVM metrics are published under the `jvm.` meter name.

The following JVM metrics are provided:

-
Various memory and buffer pool details
-
Statistics related to garbage collection
-
Thread utilization
-
Virtual threads statistics (for this, `io.micrometer:micrometer-java21`has to be on the classpath)
-
The number of classes loaded and unloaded
-
JVM version information
-
JIT compilation time

### System Metrics

Auto-configuration enables system metrics by using core Micrometer classes.
System metrics are published under the `system.`, `process.`, and `disk.` meter names.

The following system metrics are provided:

-
CPU metrics
-
File descriptor metrics
-
Uptime metrics (both the amount of time the application has been running and a fixed gauge of the absolute start time)
-
Disk space available

### Application Startup Metrics

Auto-configuration exposes application startup time metrics:

-
`application.started.time`: time taken to start the application.
-
`application.ready.time`: time taken for the application to be ready to service requests.

Metrics are tagged by the fully qualified name of the application class.

### Logger Metrics

Auto-configuration enables the event metrics for both Logback and Log4J2.
The details are published under the `log4j2.events.` or `logback.events.` meter names.

### Task Execution and Scheduling Metrics

Auto-configuration enables the instrumentation of all available `ThreadPoolTaskExecutor` and `ThreadPoolTaskScheduler` beans, as long as the underling `ThreadPoolExecutor` is available.
Metrics are tagged by the name of the executor, which is derived from the bean name.

### JMS Metrics

Auto-configuration enables the instrumentation of all available `JmsTemplate` beans and `@JmsListener` annotated methods.
This will produce `"jms.message.publish"` and `"jms.message.process"` metrics respectively.
See the Spring Framework reference documentation for more information on produced observations.

| `JmsClient`and`JmsMessagingTemplate`that uses a`JmsTemplate`bean are also instrumented. |

### Spring MVC Metrics

Auto-configuration enables the instrumentation of all requests handled by Spring MVC controllers and functional handlers.
By default, metrics are generated with the name, `http.server.requests`.
You can customize the name by setting the `management.observations.http.server.requests.name` property.

To add to the default tags, provide a `@Bean` that extends `DefaultServerRequestObservationConvention` from the `org.springframework.http.server.observation` package.
To replace the default tags, provide a `@Bean` that implements `ServerRequestObservationConvention`.

| In some cases, exceptions handled in web controllers are not recorded as request metrics tags. Applications can opt in and record exceptions by setting handled exceptions as request attributes. |

By default, all requests are handled.
To customize the filter, provide a `@Bean` that implements `FilterRegistrationBean<ServerHttpObservationFilter>`.

### Spring WebFlux Metrics

Auto-configuration enables the instrumentation of all requests handled by Spring WebFlux controllers and functional handlers.
By default, metrics are generated with the name, `http.server.requests`.
You can customize the name by setting the `management.observations.http.server.requests.name` property.

To add to the default tags, provide a `@Bean` that extends `DefaultServerRequestObservationConvention` from the `org.springframework.http.server.reactive.observation` package.
To replace the default tags, provide a `@Bean` that implements `ServerRequestObservationConvention`.

| In some cases, exceptions handled in controllers and handler functions are not recorded as request metrics tags. Applications can opt in and record exceptions by setting handled exceptions as request attributes. |

### Jersey Server Metrics

Auto-configuration enables the instrumentation of all requests handled by the Jersey JAX-RS implementation.
By default, metrics are generated with the name, `http.server.requests`.
You can customize the name by setting the `management.observations.http.server.requests.name` property.

By default, Jersey server metrics are tagged with the following information:

| Tag | Description |
|---|---|
|
 | The simple class name of any exception that was thrown while handling the request. |
|
 | The request’s method (for example, |
|
 | The request’s outcome, based on the status code of the response.
 1xx is |
|
 | The response’s HTTP status code (for example, |
|
 | The request’s URI template prior to variable substitution, if possible (for example, |

To customize the tags, provide a `@Bean` that implements `JerseyObservationConvention`.

### SSL Bundle Metrics

Spring Boot Actuator publishes expiry metrics about SSL bundles.
The metric `ssl.chain.expiry` gauges the expiry date of each certificate chain in key stores and trust stores in seconds.
This number will be negative if the chain has already expired.
This metric is tagged with the following information:

| Tag | Description |
|---|---|
|
 | The name of the bundle which contains the certificate chain |
|
 | The serial number (in hex format) of the certificate which is the soonest to expire in the chain |
|
 | The name of the certificate chain. |
|
 | Whether the certificate chain comes from the key store ( |

### HTTP Client Metrics

Spring Boot Actuator manages the instrumentation of `RestTemplate`, `WebClient` and `RestClient`.
For that, you have to inject the auto-configured builder and use it to create instances:

You can also manually apply the customizers responsible for this instrumentation, namely `ObservationRestTemplateCustomizer`, `ObservationWebClientCustomizer` and `ObservationRestClientCustomizer`.

By default, metrics are generated with the name, `http.client.requests`.
You can customize the name by setting the `management.observations.http.client.requests.name` property.

To customize the tags when using `RestTemplate` or `RestClient`, provide a `@Bean` that implements `ClientRequestObservationConvention` from the `org.springframework.http.client.observation` package.
To customize the tags when using `WebClient`, provide a `@Bean` that implements `ClientRequestObservationConvention` from the `org.springframework.web.reactive.function.client` package.

### Tomcat Metrics

Auto-configuration enables the instrumentation of Tomcat only when an MBean `Registry` is enabled.
By default, the MBean registry is disabled, but you can enable it by setting `server.tomcat.mbeanregistry.enabled` to `true`.

Tomcat metrics are published under the `tomcat.` meter name.

### Cache Metrics

Auto-configuration enables the instrumentation of all available `Cache` instances on startup, with metrics prefixed with `cache`.
Cache instrumentation is standardized for a basic set of metrics.
Additional, cache-specific metrics are also available.

The following cache libraries are supported:

-
Cache2k
-
Caffeine
-
Hazelcast
-
Any compliant JCache (JSR-107) implementation
-
Redis

| Metrics should be enabled for the auto-configuration to pick them up. Refer to the documentation of the cache library you are using for more details. |

Metrics are tagged by the name of the cache and by the name of the `CacheManager`, which is derived from the bean name.

| Only caches that are configured on startup are bound to the registry.
For caches not defined in the cache’s configuration, such as caches created on the fly or programmatically after the startup phase, an explicit registration is required.
A `CacheMetricsRegistrar`bean is made available to make that process easier. |

### Spring Batch Metrics

See the Spring Batch reference documentation.

### DataSource Metrics

Auto-configuration enables the instrumentation of all available `DataSource` objects with metrics prefixed with `jdbc.connections`.
Data source instrumentation results in gauges that represent the currently active, idle, maximum allowed, and minimum allowed connections in the pool.

Metrics are also tagged by the name of the `DataSource` computed based on the bean name.

| By default, Spring Boot provides metadata for all supported data sources.
You can add additional `DataSourcePoolMetadataProvider`beans if your favorite data source is not supported.
See`DataSourcePoolMetadataProvidersConfiguration`for examples. |

Also, Hikari-specific metrics are exposed with a `hikaricp` prefix.
Each metric is tagged by the name of the pool (you can control it with `spring.datasource.name`).

### Hibernate Metrics

If `org.hibernate.orm:hibernate-micrometer` is on the classpath, all available Hibernate `EntityManagerFactory` instances that have statistics enabled are instrumented with a metric named `hibernate`.

Metrics are also tagged by the name of the `EntityManagerFactory`, which is derived from the bean name.

To enable statistics, the standard JPA property `hibernate.generate_statistics` must be set to `true`.
You can enable that on the auto-configured `EntityManagerFactory`:

-
Properties
-
YAML

`spring.jpa.properties[hibernate.generate_statistics]=true````
spring:
 jpa:
 properties:
 "[hibernate.generate_statistics]": true
```
### Spring Data Repository Metrics

Auto-configuration enables the instrumentation of all Spring Data `Repository` method invocations.
By default, metrics are generated with the name, `spring.data.repository.invocations`.
You can customize the name by setting the `management.metrics.data.repository.metric-name` property.

The `@Timed` annotation from the `io.micrometer.core.annotation` package is supported on `Repository` interfaces and methods.
If you do not want to record metrics for all `Repository` invocations, you can set `management.metrics.data.repository.autotime.enabled` to `false` and exclusively use `@Timed` annotations instead.

| A `@Timed`annotation with`longTask = true`enables a long task timer for the method.
Long task timers require a separate metric name and can be stacked with a short task timer. |

By default, repository invocation related metrics are tagged with the following information:

| Tag | Description |
|---|---|
|
 | The simple class name of the source |
|
 | The name of the |
|
 | The result state ( |
|
 | The simple class name of any exception that was thrown from the invocation. |

To replace the default tags, provide a `@Bean` that implements `RepositoryTagsProvider`.

### RabbitMQ Metrics

Auto-configuration enables the instrumentation of all available RabbitMQ connection factories with a metric named `rabbitmq`.

### Spring Integration Metrics

Spring Integration automatically provides Micrometer support whenever a `MeterRegistry` bean is available.
Metrics are published under the `spring.integration.` meter name.

### Kafka Metrics

Auto-configuration registers a `MicrometerConsumerListener` and `MicrometerProducerListener` for the auto-configured consumer factory and producer factory, respectively.
It also registers a `KafkaStreamsMicrometerListener` for `StreamsBuilderFactoryBean`.
For more detail, see the Micrometer Native Metrics section of the Spring Kafka documentation.

### MongoDB Metrics

This section briefly describes the available metrics for MongoDB.

#### MongoDB Command Metrics

Auto-configuration registers a `MongoMetricsCommandListener` with the auto-configured `MongoClient`.

A timer metric named `mongodb.driver.commands` is created for each command issued to the underlying MongoDB driver.
Each metric is tagged with the following information by default:

| Tag | Description |
|---|---|
|
 | The name of the command issued. |
|
 | The identifier of the cluster to which the command was sent. |
|
 | The address of the server to which the command was sent. |
|
 | The outcome of the command ( |

To replace the default metric tags, define a `MongoCommandTagsProvider` bean, as the following example shows:

-
Java
-
Kotlin

```
import io.micrometer.core.instrument.binder.mongodb.MongoCommandTagsProvider;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyCommandTagsProviderConfiguration {
	@Bean
	public MongoCommandTagsProvider customCommandTagsProvider() {
 return new CustomCommandTagsProvider();
	}
}
```
```
import io.micrometer.core.instrument.binder.mongodb.MongoCommandTagsProvider
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
@Configuration(proxyBeanMethods = false)
class MyCommandTagsProviderConfiguration {
	@Bean
	fun customCommandTagsProvider(): MongoCommandTagsProvider? {
 return CustomCommandTagsProvider()
	}
}
```
To disable the auto-configured command metrics, set the following property:

-
Properties
-
YAML

`management.metrics.mongodb.command.enabled=false````
management:
 metrics:
 mongodb:
 command:
 enabled: false
```
#### MongoDB Connection Pool Metrics

Auto-configuration registers a `MongoMetricsConnectionPoolListener` with the auto-configured `MongoClient`.

The following gauge metrics are created for the connection pool:

-
`mongodb.driver.pool.size`reports the current size of the connection pool, including idle and in-use members.
-
`mongodb.driver.pool.checkedout`reports the count of connections that are currently in use.
-
`mongodb.driver.pool.waitqueuesize`reports the current size of the wait queue for a connection from the pool.

Each metric is tagged with the following information by default:

| Tag | Description |
|---|---|
|
 | The identifier of the cluster to which the connection pool corresponds. |
|
 | The address of the server to which the connection pool corresponds. |

To replace the default metric tags, define a `MongoConnectionPoolTagsProvider` bean:

-
Java
-
Kotlin

```
import io.micrometer.core.instrument.binder.mongodb.MongoConnectionPoolTagsProvider;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyConnectionPoolTagsProviderConfiguration {
	@Bean
	public MongoConnectionPoolTagsProvider customConnectionPoolTagsProvider() {
 return new CustomConnectionPoolTagsProvider();
	}
}
```
```
import io.micrometer.core.instrument.binder.mongodb.MongoConnectionPoolTagsProvider
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
@Configuration(proxyBeanMethods = false)
class MyConnectionPoolTagsProviderConfiguration {
	@Bean
	fun customConnectionPoolTagsProvider(): MongoConnectionPoolTagsProvider {
 return CustomConnectionPoolTagsProvider()
	}
}
```
To disable the auto-configured connection pool metrics, set the following property:

-
Properties
-
YAML

`management.metrics.mongodb.connectionpool.enabled=false````
management:
 metrics:
 mongodb:
 connectionpool:
 enabled: false
```
### Neo4j Metrics

Auto-configuration registers a `MicrometerObservationProvider` for the auto-configured `Driver`.

To override this behavior, you can register a `ConfigBuilderCustomizer` bean with an order higher than zero.

### Jetty Metrics

Auto-configuration binds metrics for Jetty’s `ThreadPool` by using Micrometer’s `JettyServerThreadPoolMetrics`.
Metrics for Jetty’s `Connector` instances are bound by using Micrometer’s `JettyConnectionMetrics` and, when `server.ssl.enabled` is set to `true`, Micrometer’s `JettySslHandshakeMetrics`.

### Redis Metrics

Auto-configuration registers a `MicrometerTracing` for the auto-configured `LettuceConnectionFactory`.
For more detail, see the Observability section of the Lettuce documentation.

## Registering Custom Metrics

To register custom metrics, inject `MeterRegistry` into your component:

-
Java
-
Kotlin

```
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Tags;
import org.springframework.stereotype.Component;
@Component
public class MyBean {
	private final Dictionary dictionary;
	public MyBean(MeterRegistry registry) {
 this.dictionary = Dictionary.load();
 registry.gauge("dictionary.size", Tags.empty(), this.dictionary.getWords().size());
	}
}
```
```
import io.micrometer.core.instrument.MeterRegistry
import io.micrometer.core.instrument.Tags
import org.springframework.stereotype.Component
@Component
class MyBean(registry: MeterRegistry) {
	private val dictionary: Dictionary
	init {
 dictionary = Dictionary.load()
 registry.gauge("dictionary.size", Tags.empty(), dictionary.words.size)
	}
}
```
If your metrics depend on other beans, we recommend that you use a `MeterBinder` to register them:

-
Java
-
Kotlin

```
import io.micrometer.core.instrument.Gauge;
import io.micrometer.core.instrument.binder.MeterBinder;
import org.springframework.context.annotation.Bean;
public class MyMeterBinderConfiguration {
	@Bean
	public MeterBinder queueSize(Queue queue) {
 return (registry) -> Gauge.builder("queueSize", queue::size).register(registry);
	}
}
```
```
import io.micrometer.core.instrument.Gauge
import io.micrometer.core.instrument.binder.MeterBinder
import org.springframework.context.annotation.Bean
class MyMeterBinderConfiguration {
	@Bean
	fun queueSize(queue: Queue): MeterBinder {
 return MeterBinder { registry ->
 Gauge.builder("queueSize", queue::size).register(registry)
 }
	}
}
```
Using a `MeterBinder` ensures that the correct dependency relationships are set up and that the bean is available when the metric’s value is retrieved.
A `MeterBinder` implementation can also be useful if you find that you repeatedly instrument a suite of metrics across components or applications.

| By default, metrics from all `MeterBinder`beans are automatically bound to the Spring-managed`MeterRegistry`. |

## Customizing Individual Metrics

If you need to apply customizations to specific `Meter` instances, you can use the `MeterFilter` interface.

For example, if you want to rename the `mytag.region` tag to `mytag.area` for all meter IDs beginning with `com.example`, you can do the following:

-
Java
-
Kotlin

```
import io.micrometer.core.instrument.config.MeterFilter;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyMetricsFilterConfiguration {
	@Bean
	public MeterFilter renameRegionTagMeterFilter() {
 return MeterFilter.renameTag("com.example", "mytag.region", "mytag.area");
	}
}
```
```
import io.micrometer.core.instrument.config.MeterFilter
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
@Configuration(proxyBeanMethods = false)
class MyMetricsFilterConfiguration {
	@Bean
	fun renameRegionTagMeterFilter(): MeterFilter {
 return MeterFilter.renameTag("com.example", "mytag.region", "mytag.area")
	}
}
```
| By default, all `MeterFilter`beans are automatically bound to the Spring-managed`MeterRegistry`.
Make sure to register your metrics by using the Spring-managed`MeterRegistry`and not any of the static methods on`Metrics`.
These use the global registry that is not Spring-managed. |

### Common Tags

Common tags are generally used for dimensional drill-down on the operating environment, such as host, instance, region, stack, and others. Commons tags are applied to all meters and can be configured, as the following example shows:

-
Properties
-
YAML

```
management.metrics.tags.region=us-east-1
management.metrics.tags.stack=prod
```
```
management:
 metrics:
 tags:
 region: "us-east-1"
 stack: "prod"
```
The preceding example adds `region` and `stack` tags to all meters with a value of `us-east-1` and `prod`, respectively.

| The order of common tags is important if you use Graphite.
As the order of common tags cannot be guaranteed by using this approach, Graphite users are advised to define a custom `MeterFilter`instead. |

### Per-meter Properties

In addition to `MeterFilter` beans, you can apply a limited set of customization on a per-meter basis using properties.
Per-meter customizations are applied, using Spring Boot’s `PropertiesMeterFilter`, to any meter IDs that start with the given name.
The following example filters out any meters that have an ID starting with `example.remote`.

-
Properties
-
YAML

`management.metrics.enable.example.remote=false````
management:
 metrics:
 enable:
 example:
 remote: false
```
The following properties allow per-meter customization:

| Property | Description |
|---|---|
|
 | Whether to accept meters with certain IDs.
 Meters that are not accepted are filtered from the |
|
 | Whether to publish a histogram suitable for computing aggregable (across dimension) percentile approximations. |
|
 | Publish fewer histogram buckets by clamping the range of expected values. |
|
 | Publish percentile values computed in your application |
|
 | Give greater weight to recent samples by accumulating them in ring buffers which rotate after a configurable expiry, with a configurable buffer length. |
|
 | Publish a cumulative histogram with buckets defined by your service-level objectives. |

For more details on the concepts behind `percentiles-histogram`, `percentiles`, and `slo`, see the Histograms and percentiles section of the Micrometer documentation.

## Metrics Endpoint

Spring Boot provides a `metrics` endpoint that you can use diagnostically to examine the metrics collected by an application.
The endpoint is not available by default and must be exposed.
See exposing endpoints for more details.

Navigating to `/actuator/metrics` displays a list of available meter names.
You can drill down to view information about a particular meter by providing its name as a selector — for example, `/actuator/metrics/jvm.memory.max`.

| The name you use here should match the name used in the code, not the name after it has been naming-convention normalized for a monitoring system to which it is shipped.
In other words, if |

You can also add any number of `tag=KEY:VALUE` query parameters to the end of the URL to dimensionally drill down on a meter — for example, `/actuator/metrics/jvm.memory.max?tag=area:nonheap`.

| The reported measurements are the |

## Integration with Micrometer Observation

A `DefaultMeterObservationHandler` is automatically registered on the `ObservationRegistry`, which creates metrics for every completed observation.

# Tracing

Spring Boot Actuator provides dependency management and auto-configuration for Micrometer Tracing, a facade for popular tracer libraries.

| To learn more about Micrometer Tracing capabilities, see its reference documentation. |

## Supported Tracers

Spring Boot ships auto-configuration for the following tracers:

-
OpenTelemetry with OTLP.
-
OpenZipkin Brave with Zipkin.

## Getting Started

We need an example application that we can use to get started with tracing. For our purposes, the simple “Hello World!” web application that’s covered in the Developing Your First Spring Boot Application section will suffice. We’re going to use the Brave tracer with Zipkin as trace backend.

To recap, our main application code looks like this:

-
Java
-
Kotlin

```
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
@RestController
@SpringBootApplication
public class MyApplication {
	private static final Log logger = LogFactory.getLog(MyApplication.class);
	@RequestMapping("/")
	String home() {
 logger.info("home() has been called");
 return "Hello World!";
	}
	public static void main(String[] args) {
 SpringApplication.run(MyApplication.class, args);
	}
}
```
```
import org.apache.commons.logging.Log
import org.apache.commons.logging.LogFactory
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController
@RestController
@SpringBootApplication
class MyApplication {
	private val logger: Log = LogFactory.getLog(MyApplication::class.java)
	@RequestMapping("/")
	fun home(): String {
 logger.info("home() has been called")
 return "Hello, World!"
	}
}
fun main(args: Array<String>) {
	runApplication<MyApplication>(*args)
}
```
| There’s an added logger statement in the `home()`method, which will be important later. |

Now we have to add the `org.springframework.boot:spring-boot-starter-zipkin` dependency.

Then add the following application properties:

-
Properties
-
YAML

`management.tracing.sampling.probability=1````
management:
 tracing:
 sampling:
 probability: 1.0
```
By default, Spring Boot samples only 10% of requests to prevent overwhelming the trace backend. This property switches it to 100% so that every request is sent to the trace backend.

To collect and visualize the traces, we need a running trace backend. We use Zipkin as our trace backend here. The Zipkin Quickstart guide provides instructions how to start Zipkin locally.

After Zipkin is running, you can start your application.

If you open a web browser to `localhost:8080`, you should see the following output:

`Hello World!`Behind the scenes, an observation has been created for the HTTP request, which in turn gets bridged to Brave, which reports a new trace to Zipkin.

Now open the Zipkin UI at `localhost:9411` and press the "Run Query" button to list all collected traces.
You should see one trace.
Press the "Show" button to see the details of that trace.

## Logging Correlation IDs

Correlation IDs provide a helpful way to link lines in your log files to spans/traces. If you are using Micrometer Tracing, Spring Boot will include correlation IDs in your logs by default.

The default correlation ID is built from `traceId` and `spanId` MDC values.
For example, if Micrometer Tracing has added an MDC `traceId` of `803B448A0489F84084905D3093480352` and an MDC `spanId` of `3425F23BB2432450` the log output will include the correlation ID `[803B448A0489F84084905D3093480352-3425F23BB2432450]`.

If you prefer to use a different format for your correlation ID, you can use the `logging.pattern.correlation` property to define one.
For example, the following will provide a correlation ID for Logback in format previously used by Spring Cloud Sleuth:

-
Properties
-
YAML

```
logging.pattern.correlation=[${spring.application.name:},%X{traceId:-},%X{spanId:-}]
logging.include-application-name=false
```
```
logging:
 pattern:
 correlation: "[${spring.application.name:},%X{traceId:-},%X{spanId:-}] "
 include-application-name: false
```
| In the example above, `logging.include-application-name`is set to`false`to avoid the application name being duplicated in the log messages (`logging.pattern.correlation`already contains it).
It’s also worth mentioning that`logging.pattern.correlation`contains a trailing space so that it is separated from the logger name that comes right after it by default. |

| Correlation IDs rely on context propagation. Please read this documentation for more details. |

## Propagating Traces

To automatically propagate traces over the network, use the auto-configured `RestTemplateBuilder`, `RestClient.Builder` or `WebClient.Builder` to construct the client.

| If you create the `RestTemplate`, the`RestClient`or the`WebClient`without using the auto-configured builders, automatic trace propagation won’t work! |

## Tracer Implementations

As Micrometer Tracer supports multiple tracer implementations, there are multiple dependency combinations possible with Spring Boot. The combinations OpenTelemetry with OTLP and Brave with Zipkin are common and have dedicated starters.

### OpenTelemetry With OTLP

Tracing with OpenTelemetry and reporting using OTLP requires the following dependencies:

-
`org.springframework.boot:spring-boot-starter-opentelemetry`

Use the `management.opentelemetry.tracing.export.otlp.*` configuration properties to configure reporting using OTLP.

| If you need to apply advanced customizations to OTLP span exporters, consider registering `OtlpHttpSpanExporterBuilderCustomizer`or`OtlpGrpcSpanExporterBuilderCustomizer`beans.
These will be invoked before the creation of the`OtlpHttpSpanExporter`or`OtlpGrpcSpanExporter`.
The customizers take precedence over anything applied by the auto-configuration. |

### OpenTelemetry With Zipkin

| OpenTelemetry has deprecated their Zipkin support. The auto-configuration for it will be removed in Spring Boot 4.2. Either switch to Brave or consider using the Zipkin OTel module for ingesting OTLP directly. |

Tracing with OpenTelemetry and reporting to Zipkin requires the following dependencies:

-
`org.springframework.boot:spring-boot-micrometer-tracing-opentelemetry`- Spring Boot’s support for Micrometer Tracing over OpenTelemetry.
-
`io.micrometer:micrometer-tracing-bridge-otel`- bridges the Micrometer Observation API to OpenTelemetry.
-
`org.springframework.boot:spring-boot-zipkin`- Spring Boot’s support for Zipkin.
-
`io.opentelemetry:opentelemetry-exporter-zipkin`- OpenTelemetry exporter that reports traces to Zipkin.

Use the `management.tracing.export.zipkin.*` configuration properties to configure reporting to Zipkin.

## Sampling

By default, Spring Boot samples only 10% of requests to prevent overwhelming the trace backend.
The `management.tracing.sampling.probability` property can be used to configure this.

When using OpenTelemetry, you can also configure which sampler is used via the `management.opentelemetry.tracing.sampler` property.
The following samplers are supported:

| Sampler | Description |
|---|---|
|
 | Samples every trace. |
|
 | Discards every trace. |
|
 | Samples a fraction of traces based on |
|
 | If the parent span is sampled, samples the child span. If there is no parent, samples every trace. |
|
 | If the parent span is sampled, samples the child span. If there is no parent, discards every trace. |
|
 | If the parent span is sampled, samples the child span. If there is no parent, samples a fraction of traces based on |

## Span Limits

When using OpenTelemetry, you can configure span limits via the `management.opentelemetry.tracing.limits.*` configuration properties.
These allow you to control the maximum number of attributes, events, and links per span, as well as the maximum length of string attribute values.

For example, to limit the number of attributes per span to 64 and the maximum attribute value length to 256 characters:

-
Properties
-
YAML

```
management.opentelemetry.tracing.limits.max-attributes=64
management.opentelemetry.tracing.limits.max-attribute-value-length=256
```
```
management:
 opentelemetry:
 tracing:
 limits:
 max-attributes: 64
 max-attribute-value-length: 256
```
If you need full control, you can register a custom `SpanLimits` bean.

## Integration with Micrometer Observation

A `TracingAwareMeterObservationHandler` is automatically registered on the `ObservationRegistry`, which creates spans for every completed observation.

## Creating Custom Spans

You can create your own spans by starting an observation.
For this, inject `ObservationRegistry` into your component:

-
Java
-
Kotlin

```
import io.micrometer.observation.Observation;
import io.micrometer.observation.ObservationRegistry;
import org.springframework.stereotype.Component;
@Component
class CustomObservation {
	private final ObservationRegistry observationRegistry;
	CustomObservation(ObservationRegistry observationRegistry) {
 this.observationRegistry = observationRegistry;
	}
	void someOperation() {
 Observation observation = Observation.createNotStarted("some-operation", this.observationRegistry);
 observation.lowCardinalityKeyValue("some-tag", "some-value");
 observation.observe(() -> {
 // Business logic ...
 });
	}
}
```
```
import io.micrometer.observation.Observation
import io.micrometer.observation.ObservationRegistry
import org.springframework.stereotype.Component
@Component
class CustomObservation(private val observationRegistry: ObservationRegistry) {
	fun someOperation() {
 Observation.createNotStarted("some-operation", observationRegistry)
 .lowCardinalityKeyValue("some-tag", "some-value")
 .observe {
 // Business logic ...
 }
	}
}
```
This will create an observation named "some-operation" with the tag "some-tag=some-value".

| If you want to create a span without creating a metric, you need to use the lower-level `Tracer`API from Micrometer. |

## Baggage

You can create baggage with the `Tracer` API:

-
Java
-
Kotlin

```
import io.micrometer.tracing.BaggageInScope;
import io.micrometer.tracing.Tracer;
import org.springframework.stereotype.Component;
@Component
class CreatingBaggage {
	private final Tracer tracer;
	CreatingBaggage(Tracer tracer) {
 this.tracer = tracer;
	}
	void doSomething() {
 try (BaggageInScope scope = this.tracer.createBaggageInScope("baggage1", "value1")) {
 // Business logic
 }
	}
}
```
```
import io.micrometer.tracing.Tracer
import org.springframework.stereotype.Component
@Component
class CreatingBaggage(private val tracer: Tracer) {
	fun doSomething() {
 tracer.createBaggageInScope("baggage1", "value1").use {
 // Business logic
 }
	}
}
```
This example creates baggage named `baggage1` with the value `value1`.
The baggage is automatically propagated over the network if you’re using W3C propagation.
If you’re using B3 propagation, baggage is not automatically propagated.
To manually propagate baggage over the network, use the `management.tracing.baggage.remote-fields` configuration property (this works for W3C, too).
For the example above, setting this property to `baggage1` results in an HTTP header `baggage1: value1`.

If you want to propagate the baggage to the MDC, use the `management.tracing.baggage.correlation.fields` configuration property.
For the example above, setting this property to `baggage1` results in an MDC entry named `baggage1`.

## Tests

Tracing components which are reporting data are not auto-configured when using `@SpringBootTest`.
See Using Tracing for more details.

# Auditing

Once Spring Security is in play, Spring Boot Actuator has a flexible audit framework that publishes events (by default, “authentication success”, “failure” and “access denied” exceptions). This feature can be very useful for reporting and for implementing a lock-out policy based on authentication failures.

You can enable auditing by providing a bean of type `AuditEventRepository` in your application’s configuration.
For convenience, Spring Boot offers an `InMemoryAuditEventRepository`.
`InMemoryAuditEventRepository` has limited capabilities, and we recommend using it only for development environments.
For production environments, consider creating your own alternative `AuditEventRepository` implementation.

## Custom Auditing

To customize published security events, you can provide your own implementations of `AbstractAuthenticationAuditListener` and `AbstractAuthorizationAuditListener`.

You can also use the audit services for your own business events.
To do so, either inject the `AuditEventRepository` bean into your own components and use that directly or publish an `AuditApplicationEvent` with the Spring `ApplicationEventPublisher` (by implementing `ApplicationEventPublisherAware`).

# Recording HTTP Exchanges

You can enable recording of HTTP exchanges by providing a bean of type `HttpExchangeRepository` in your application’s configuration.
For convenience, Spring Boot offers `InMemoryHttpExchangeRepository`, which, by default, stores the last 100 request-response exchanges.
`InMemoryHttpExchangeRepository` is limited compared to tracing solutions, and we recommend using it only for development environments.
For production environments, we recommend using a production-ready tracing or observability solution, such as Zipkin or OpenTelemetry.
Alternatively, you can create your own `HttpExchangeRepository`.

You can use the `httpexchanges` endpoint to obtain information about the request-response exchanges that are stored in the `HttpExchangeRepository`.

# Process Monitoring

In the `spring-boot` module, you can find two classes to create files that are often useful for process monitoring:

-
`ApplicationPidFileWriter`creates a file that contains the application PID (by default, in the application directory with a file name of`application.pid`).
-
`WebServerPortFileWriter`creates a file (or files) that contain the ports of the running web server (by default, in the application directory with a file name of`application.port`).

By default, these writers are not activated, but you can enable them:

## Extending Configuration

In the `META-INF/spring.factories` file, you can activate the listener (or listeners) that writes a PID file:

```
org.springframework.context.ApplicationListener=\
org.springframework.boot.context.ApplicationPidFileWriter,\
org.springframework.boot.web.server.context.WebServerPortFileWriter
```

# Cloud Foundry Support

Spring Boot’s `spring-boot-cloudfoundry` module (part of `spring-boot-starter-cloudfoundry`) includes additional support that is activated when you deploy to a compatible Cloud Foundry instance.
The `/cloudfoundryapplication` path provides an alternative secured route to all `@Endpoint` beans.

The extended support lets Cloud Foundry management UIs (such as the web application that you can use to view deployed applications) be augmented with Spring Boot actuator information. For example, an application status page can include full health information instead of the typical “running” or “stopped” status.

| The `/cloudfoundryapplication`path is not directly accessible to regular users.
To use the endpoint, you must pass a valid UAA token with the request. |

## Disabling Extended Cloud Foundry Actuator Support

If you want to fully disable the `/cloudfoundryapplication` endpoints, you can add the following setting to your `application.properties` file:

-
Properties
-
YAML

`management.cloudfoundry.enabled=false````
management:
 cloudfoundry:
 enabled: false
```
## Cloud Foundry Self-signed Certificates

By default, the security verification for `/cloudfoundryapplication` endpoints makes SSL calls to various Cloud Foundry services.
If your Cloud Foundry UAA or Cloud Controller services use self-signed certificates, you need to set the following property:

-
Properties
-
YAML

`management.cloudfoundry.skip-ssl-validation=true````
management:
 cloudfoundry:
 skip-ssl-validation: true
```
## Custom Context Path

If the server’s context-path has been configured to anything other than `/`, the Cloud Foundry endpoints are not available at the root of the application.
For example, if `server.servlet.context-path=/app`, Cloud Foundry endpoints are available at `/app/cloudfoundryapplication/*`.

If you expect the Cloud Foundry endpoints to always be available at `/cloudfoundryapplication/*`, regardless of the server’s context-path, you need to explicitly configure that in your application.
The configuration differs, depending on the web server in use.
For Tomcat, you can add the following configuration:

-
Java
-
Kotlin

```
import java.io.IOException;
import java.util.Collections;
import jakarta.servlet.GenericServlet;
import jakarta.servlet.Servlet;
import jakarta.servlet.ServletContainerInitializer;
import jakarta.servlet.ServletContext;
import jakarta.servlet.ServletException;
import jakarta.servlet.ServletRequest;
import jakarta.servlet.ServletResponse;
import org.apache.catalina.Host;
import org.apache.catalina.core.StandardContext;
import org.apache.catalina.startup.Tomcat;
import org.springframework.boot.tomcat.servlet.TomcatServletWebServerFactory;
import org.springframework.boot.web.servlet.ServletContextInitializer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
@Configuration(proxyBeanMethods = false)
public class MyCloudFoundryConfiguration {
	@Bean
	public TomcatServletWebServerFactory servletWebServerFactory() {
 return new TomcatServletWebServerFactory() {
 @Override
 protected void prepareContext(Host host, ServletContextInitializer[] initializers, TempDirs tempDirs) {
 super.prepareContext(host, initializers, tempDirs);
 StandardContext child = new StandardContext();
 child.addLifecycleListener(new Tomcat.FixContextListener());
 child.setPath("/cloudfoundryapplication");
 ServletContainerInitializer initializer = getServletContextInitializer(getContextPath());
 child.addServletContainerInitializer(initializer, Collections.emptySet());
 child.setCrossContext(true);
 host.addChild(child);
 }
 };
	}
	private ServletContainerInitializer getServletContextInitializer(String contextPath) {
 return (classes, context) -> {
 Servlet servlet = new GenericServlet() {
 @Override
 public void service(ServletRequest req, ServletResponse res) throws ServletException, IOException {
 ServletContext context = req.getServletContext().getContext(contextPath);
 context.getRequestDispatcher("/cloudfoundryapplication").forward(req, res);
 }
 };
 context.addServlet("cloudfoundry", servlet).addMapping("/*");
 };
	}
}
```
```
import jakarta.servlet.GenericServlet
import jakarta.servlet.Servlet
import jakarta.servlet.ServletContainerInitializer
import jakarta.servlet.ServletContext
import jakarta.servlet.ServletException
import jakarta.servlet.ServletRequest
import jakarta.servlet.ServletResponse
import org.apache.catalina.Host
import org.apache.catalina.core.StandardContext
import org.apache.catalina.startup.Tomcat.FixContextListener
import org.springframework.boot.tomcat.servlet.TomcatServletWebServerFactory
import org.springframework.boot.web.servlet.ServletContextInitializer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import java.io.IOException
import java.util.Collections.emptySet
@Configuration(proxyBeanMethods = false)
class MyCloudFoundryConfiguration {
	@Bean
	fun servletWebServerFactory(): TomcatServletWebServerFactory {
 return object : TomcatServletWebServerFactory() {
 override fun prepareContext(host: Host, initializers: Array<ServletContextInitializer>, tempDirs: TempDirs) {
 super.prepareContext(host, initializers, tempDirs)
 val child = StandardContext()
 child.addLifecycleListener(FixContextListener())
 child.path = "/cloudfoundryapplication"
 val initializer = getServletContextInitializer(contextPath)
 child.addServletContainerInitializer(initializer, emptySet())
 child.crossContext = true
 host.addChild(child)
 }
 }
	}
	private fun getServletContextInitializer(contextPath: String): ServletContainerInitializer {
 return ServletContainerInitializer { classes: Set<Class<*>?>?, context: ServletContext ->
 val servlet: Servlet = object : GenericServlet() {
 @Throws(ServletException::class, IOException::class)
 override fun service(req: ServletRequest, res: ServletResponse) {
 val servletContext = req.servletContext.getContext(contextPath)
 servletContext.getRequestDispatcher("/cloudfoundryapplication").forward(req, res)
 }
 }
 context.addServlet("cloudfoundry", servlet).addMapping("/*")
 }
	}
}
```
If you’re using a Webflux based application, you can use the following configuration:

-
Java
-
Kotlin

```
import java.util.Map;
import reactor.core.publisher.Mono;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.boot.webflux.autoconfigure.WebFluxProperties;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.server.reactive.ContextPathCompositeHandler;
import org.springframework.http.server.reactive.HttpHandler;
import org.springframework.http.server.reactive.ServerHttpRequest;
import org.springframework.http.server.reactive.ServerHttpResponse;
import org.springframework.web.server.adapter.WebHttpHandlerBuilder;
@Configuration(proxyBeanMethods = false)
@EnableConfigurationProperties(WebFluxProperties.class)
public class MyReactiveCloudFoundryConfiguration {
	@Bean
	public HttpHandler httpHandler(ApplicationContext applicationContext, WebFluxProperties properties) {
 HttpHandler httpHandler = WebHttpHandlerBuilder.applicationContext(applicationContext).build();
 return new CloudFoundryHttpHandler(properties.getBasePath(), httpHandler);
	}
	private static final class CloudFoundryHttpHandler implements HttpHandler {
 private final HttpHandler delegate;
 private final ContextPathCompositeHandler contextPathDelegate;
 private CloudFoundryHttpHandler(String basePath, HttpHandler delegate) {
 this.delegate = delegate;
 this.contextPathDelegate = new ContextPathCompositeHandler(Map.of(basePath, delegate));
 }
 @Override
 public Mono<Void> handle(ServerHttpRequest request, ServerHttpResponse response) {
 // Remove underlying context path first (e.g. Servlet container)
 String path = request.getPath().pathWithinApplication().value();
 if (path.startsWith("/cloudfoundryapplication")) {
 return this.delegate.handle(request, response);
 }
 else {
 return this.contextPathDelegate.handle(request, response);
 }
 }
	}
}
```
```
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.boot.webflux.autoconfigure.WebFluxProperties
import org.springframework.context.ApplicationContext
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.http.server.reactive.ContextPathCompositeHandler
import org.springframework.http.server.reactive.HttpHandler
import org.springframework.http.server.reactive.ServerHttpRequest
import org.springframework.http.server.reactive.ServerHttpResponse
import org.springframework.web.server.adapter.WebHttpHandlerBuilder
import reactor.core.publisher.Mono
@Configuration(proxyBeanMethods = false)
@EnableConfigurationProperties(WebFluxProperties::class)
class MyReactiveCloudFoundryConfiguration {
	@Bean
	fun httpHandler(applicationContext: ApplicationContext, properties: WebFluxProperties): HttpHandler {
 val httpHandler = WebHttpHandlerBuilder.applicationContext(applicationContext).build()
 return CloudFoundryHttpHandler(properties.basePath ?: "/", httpHandler)
	}
	private class CloudFoundryHttpHandler(basePath: String, private val delegate: HttpHandler) : HttpHandler {
 private val contextPathDelegate = ContextPathCompositeHandler(mapOf(basePath to delegate))
 override fun handle(request: ServerHttpRequest, response: ServerHttpResponse): Mono<Void> {
 // Remove underlying context path first (e.g. Servlet container)
 val path = request.path.pathWithinApplication().value()
 return if (path.startsWith("/cloudfoundryapplication")) {
 delegate.handle(request, response)
 } else {
 contextPathDelegate.handle(request, response)
 }
 }
	}
}
```
