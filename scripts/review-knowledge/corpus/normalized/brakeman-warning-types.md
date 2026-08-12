# Warning Types

Read more about the different warnings Brakeman reports:

- Attribute Restriction
- Authentication
- Basic Authentication
- Command Injection
- Cross-Site Request Forgery
- Cross Site Scripting
- Cross Site Scripting (Content Tag)
- Cross Site Scripting (JSON)
- Dangerous Evaluation
- Dangerous Send
- Default Routes
- Denial of Service
- Divide By Zero
- Dynamic Render Paths
- File Access
- Format Validation
- Information Disclosure
- Mail Link
- Mass Assignment
- Path Traversal
- Remote Code Execution
- Remote Execution in YAML.load
- Session Manipulation
- Session Settings
- SQL Injection
- SSL Verification Bypass
- Unmaintained Dependencies
- Unsafe Deserialization
- Unscoped Find
- Unsafe Redirects
- Weak Hash

# Attribute Restriction

This warning type only applies to Ruby on Rails applications which are not using strong parameters.

Note that disabling mass assignment globally will suppress these warnings.

#### Missing Protection

This warning comes up if a model does not limit what attributes can be set through mass assignment.

In particular, this check looks for `attr_accessible` inside model definitions. If it is not found, this warning will be issued.

#### Use of Blacklist

Brakeman also warns on use of `attr_protected` - especially since it was found to be vulnerable to bypass. Warnings for mass assignment on models using `attr_protected` will be reported, but at a lower confidence level.

#### Suggested Remediation

For newer Ruby on Rails applications, query parameters should be whitelisted before use via strong parameters.

For older Ruby on Rails applications, each model should use `attr_accessible` to carefully whitelist which attributes may be set via mass assignment, if any.

Back to Warning Types

# Authentication

“Authentication” is the act of verifying that a user or client is who they say they are.

Right now, the only Brakeman warning in the authentication category is regarding hardcoded passwords. Brakeman will warn about constants with literal string values that appear to be passwords.

Hardcoded passwords are security issues since they imply a single password and that password is stored in the source code. Typically source code is available to a wide number of people inside an organization, and there have been many instances of source code leaking to the public. Passwords and secrets should be stored in a separate, secure location to limit access.

Additionally, it is recommended not to use a single password for accessing sensitive information. Each user should have their own password to make it easier to audit and revoke access.

Back to Warning Types

# Basic Authentication

In Rails 3.1, a new feature was added to simplify basic authentication.

The example provided in the official Rails Guide looks like this:

```
class PostsController < ApplicationController
 http_basic_authenticate_with :name => "dhh", :password => "secret", :except => :index
 #...
end
```
This warning will be raised if `http_basic_authenticate_with` is used and the password is found to be a string (i.e., stored somewhere in the code).

Back to Warning Types

# Command Injection

Injection is #1 on the 2010 OWASP Top Ten web security risks. Command injection occurs when shell commands unsafely include user-manipulatable values.

There are many ways to run commands in Ruby:

```
`ls #{params[:file]}`
system("ls #{params[:dir]}")
exec("md5sum #{params[:input]}")
```
Brakeman will warn on any method like these that uses user input or unsafely interpolates variables.

See the Ruby Security Guide for details.

Back to Warning Types

# Cross Site Request Forgery

Cross-site request forgery is #5 on the OWASP Top Ten. CSRF allows an attacker to perform actions on a website as if they are an authenticated user.

This warning is raised when no call to `protect_from_forgery` is found in `ApplicationController`. This method prevents CSRF.

For Rails 4 applications, it is recommended that you use `protect_from_forgery :with => :exception`. This code is inserted into newly generated applications. The default is to `nil` out the session object, which has been a source of many CSRF bypasses due to session memoization.

See the Ruby Security Guide for details.

Back to Warning Types

# Cross Site Scripting

Cross site scripting (or XSS) is #2 on the 2010 OWASP Top Ten web security risks and it pops up nearly everywhere.

XSS occurs when a user-manipulatable value is displayed on a web page without escaping it, allowing someone to inject Javascript or HTML into the page.

In Rails 2.x, values need to be explicitly escaped (e.g., by using the `h` method). In Rails 3.x, auto-escaping in views is enabled by default. However, one can still use the `raw` method to output a value directly.

See the Ruby Security Guide for more details.

### Query Parameters and Cookies

Rails 2.x example in ERB:

```
<%= params[:query] %>
```
Brakeman looks for several situations that can allow XSS. The simplest is like the example above: a value from the `params` or `cookies` is being directly output to a view. In such cases, it will issue a warning like:

```
Unescaped parameter value near line 3: params[:query]
```
By default, Brakeman will also warn when a parameter or cookie value is used as an argument to a method, the result of which is output unescaped to a view.

For example:

```
<%= some_method(cookie[:name]) %>
```
This raises a warning like:

```
Unescaped cookie value near line 5: some_method(cookies[:oreo])
```
However, the confidence level for this warning will be weak, because it is not directly outputting the cookie value.

Some methods are known to Brakeman to either be dangerous (`link_to` is one) or safe (`escape_once`). Users can specify safe methods using the `--safe-methods` option. Alternatively, Brakeman can be set to *only* warn when values are used directly with the `--report-direct` option.

### Model Attributes

Because (many) models come from database values, Brakeman mistrusts them by default.

For example, if `@user` is an instance of a model set in an action like

```
def set_user
 @user = User.first
end
```
and there is a view with

```
<%= @user.name %>
```
Brakeman will raise a warning like

```
Unescaped model attribute near line 3: User.first.name
```
If you trust all your data (although you probably shouldn’t), this can be disabled with `--ignore-model-output`.

Back to Warning Types

# Cross Site Scripting (Content Tag)

Cross site scripting (or XSS) is #2 on the 2010 OWASP Top Ten web security risks and it pops up nearly everywhere. XSS occurs when a user-manipulatable value is displayed on a web page without escaping it, allowing someone to inject Javascript or HTML into the page.

content_tag is a view helper which generates an HTML tag with some content:

```
>> content_tag :p, "Hi!"
=> "<p>Hi!</p>"
```
In Rails 2, this content is unescaped (although attribute values are escaped):

```
>> content_tag :p, "<script>alert(1)</script>"
=> "<p><script>alert(1)</script></p>"
```
In Rails 3, the content is escaped. However, only the *content* and the tag attribute *values* are escaped. The tag and attribute names are never escaped in Rails 2 or 3.

This is more dangerous than a typical method call because `content_tag` marks its output as “HTML safe”, meaning the `rails_xss` plugin and Rails 3 auto-escaping will not escape its output. Due to this, `content_tag` should be used carefully if user input is provided as an argument.

Note that while `content_tag` does have an `escape` parameter, this only applies to tag attribute *values* and is true by default.

Back to Warning Types

# Cross Site Scripting (JSON)

Cross site scripting (or XSS) is #2 on the 2010 OWASP Top Ten web security risks and it pops up nearly everywhere.

XSS occurs when a user-manipulatable value is displayed on a web page without escaping it, allowing someone to inject Javascript or HTML into the page. Calls to `Hash#to_json` can be used to trigger XSS. Brakeman will check to see if there are any calls to `Hash#to_json` with `ActiveSupport#escape_html_entities_in_json` set to false (or if you are running Rails < 2.1.0 which did not have this functionality).

`ActiveSupport#escape_html_entities_in_json` was introduced in the “new_rails_defaults” initializer in Rails 2.1.0 which is set to `false` by default. In Rails 3.0.0, `true` became the default setting. Setting this value to `true` will automatically escape ‘<’, ‘>’, ‘&’ which are commonly used to break out of code generated by a to_json call.

See ActiveSupport#escape_html_entities_in_json for more details.

### Exploiting to_json

Consider the following snippet of Rails 2.x ERB:

```
# controller
@attrs = {:email => '[email protected]</script><script>alert(document.domain)//'}
<!-- view -->
<script>
 var attributes = <%= @attrs.to_json %>
</script>
```
Which generates the following html:

```
<script>
 var attributes = {"email":"[email protected]</script><script>alert(document.domain)//"}
</script>
```
While the generated Javascript appears valid, the browser parses the script tags first, so it sees something like this:

```
<script>
 var attributes = {"email":"[email protected]
</script>
<script>
 alert(document.domain)//"}
</script>
```
The attribute assignment causes a Javascript error, but the alert triggers just fine!

With `escape_html_entities_in_json = true`, you will receive the following innocuous output:

```
<script>
 var attributes = {"email":"[email protected]\u003C/script\u003E\u003Cscript\u003Ealert(document.domain)//"}
</script>
```
Back to Warning Types

Brakeman
Documentation
News
Code
Dangerous Evaluation
Content moved to
Dangerous Eval
.

# Dangerous Send

Using unfiltered user data to select a Class or Method to be dynamically sent is dangerous.

It is much safer to whitelist the desired target or method.

Unsafe use of method:

```
method = params[:method]
@result = User.send(method.to_sym)
```
Safe:

```
method = params[:method] == 1 ? :method_a : :method_b
@result = User.send(method, *args)
```
Unsafe use of target:

```
table = params[:table]
model = table.classify.constantize
@result = model.send(:method)
```
Safe:

```
target = params[:target] == 1 ? Account : User
@result = target.send(:method, *args)
```
Including user data in the arguments passed to an Object#send is safe, as long as the method can properly handle potentially bad data.

Safe:

```
args = params["args"] || []
@result = User.send(:method, *args)
```
Back to Warning Types

# Default Routes

The general default routes warning means there is a call to

```
#Rails 2.x
map.connect ":controller/:action/:id"
```
or

```
Rails 3.x
match ':controller(/:action(/:id(.:format)))'
```
in `config/routes.rb`. This allows any public method on any controller to be called as an action.

If this warning is reported for a particular controller, it means there is a route to that controller containing `:action`.

Default routes can be dangerous if methods are made public which are not intended to be used as URLs or actions.

Back to Warning Types

# Denial of Service

Denial of Service (DoS) is any attack which causes a service to become unavailable for legitimate clients.

Denial of Service can be caused by consuming large amounts of network, memory, or CPU resources.

### Regex DoS

If an attacker can control the content of a regular expression, they may be able to construct a regular expression that requires exponential time to run.

Brakeman will warn about dynamic regular expressions that inject user-supplied values.

For example:

```
some.values.any? { |v| v.match /#{params[:query]}/ }
```
More information:

### Symbol DoS

Prior to Ruby 2.2, Symbols were not garbage collected. Creation of large numbers of Symbols could lead to a server running out of memory.

If the application appears to be using an older version of Ruby, Brakeman checks for code where user input which is converted to a Symbol. When this is not restricted, an attacker could create an unlimited number of Symbols.

Note: This is an optional check which can be enabled with `--enable SymbolDoS` or `--run-all-checks`.

Back to Warning Types

# Divide By Zero

Integer division by zero (`0`) in Ruby results in a `ZeroDivisionError` exception.

While not strictly a security issue, if an attacker can trigger a large number of exceptions it can harm site availability.

Brakeman warns when it finds potential division by zero with integers. Dividing a float by zero or `0.0` in Ruby results in `Infinity`, not an exception.

Back to Warning Types

# Dynamic Render Path

When a call to `render` uses a dynamically generated path, template name, file name, or action, there is the possibility that a user can access templates that should be restricted. The issue may be worse if those templates execute code or modify the database.

This warning is shown whenever the path to be rendered is not a static string or symbol.

These warnings are often false positives, however, because it can be difficult to manipulate Rails’ assumptions about paths to perform malicious behavior. Reports of dynamic render paths should be checked carefully to see if they can actually be manipulated maliciously by the user.

Back to Warning Types

# File Access

Using user input when accessing files (local or remote) will raise a warning in Brakeman.

For example

```
File.open("/tmp/#{cookie[:file]}")
```
will raise an error like

```
Cookie value used in file name near line 4: File.open("/tmp/#{cookie[:file]}")
```
This type of vulnerability can be used to access arbitrary files on a server (including `/etc/passwd`.

Back to Warning Types

# Format Validation

Calls to `validates_format_of ..., :with => //` which do not use `\A` and `\z` as anchors will cause this warning. Using `^` and `$` is not sufficient, as they will only match up to a new line. This allows an attacker to put whatever malicious input they would like before or after a new line character.

See the Ruby Security Guide for details.

Back to Warning Types

# Information Disclosure

Also known as information leakage or information exposure, this vulnerability refers to system or internal information (such as debugging output, stack traces, error messages, etc.) which is displayed to an end user.

For example, Rails provides detailed exception reports by default in the development environment, but it is turned off by default in production:

```
# Full error reports are disabled
config.consider_all_requests_local = false
```
Brakeman warns if this setting is `true` in production or there is a `show_detailed_exceptions?` method in a controller which does not return `false`.

Back to Warning Types

# Mail Link (CVE-2011-0446)

Certain versions of Rails were vulnerable to a cross-site scripting vulnerability mail_to.

Versions of Rails after 2.3.10 or 3.0.3 are not affected. Updating or removing the mail_to links is advised.

For more details see CVE-2011-0446.

Back to Warning Types

# Mass Assignment

Mass assignment is a feature of Rails which allows an application to create a record from the values of a hash.

Example:

```
User.new(params[:user])
```
Unfortunately, if there is a user field called `admin` which controls administrator access, now any user can make themselves an administrator with a query like

```
?user[admin]=true
```
### Rails With Strong Parameters

In Rails 4 and newer, protection for mass assignment is on by default.

Query parameters must be explicitly whitelisted via `permit` in order to be used in mass assignment:

User.new(params.permit(:name, :password))

Care should be taken to only whitelist values that are safe for a user (or attacker) to set. Foreign keys such as `account_id` are likely unsafe, allowing an attacker to manipulate records belonging to other accounts.

Brakeman will warn on potentially dangerous attributes that are whitelisted.

Brakeman will also warn about uses of `params.permit!`, since that allows everything.

### Rails Without Strong Parameters

In older versions of Rails, `attr_accessible` and `attr_protected` can be used to limit mass assignment.
However, Brakeman will warn unless `attr_accessible` is used, or mass assignment is completely disabled.

There are two different mass assignment warnings which can arise. The first is when mass assignment actually occurs, such as the example above. This results in a warning like

```
Unprotected mass assignment near line 61: User.new(params[:user])
```
The other warning is raised whenever a model is found which does not use `attr_accessible`. This produces generic warnings like

```
Mass assignment is not restricted using attr_accessible
```
with a list of affected models.

In Rails 3.1 and newer, mass assignment can easily be disabled:

```
config.active_record.whitelist_attributes = true
```
Unfortunately, it can also easily be bypassed:

```
User.new(params[:user], :without_protection => true)
```
Brakeman will warn on uses of `without_protection`.

### More Information

Strong Parameters in Rails Security Guide Mass Assignment in Rails Security Guide

Back to Warning Types

# Path Traversal

Path traversal vulnerabilities allow an attacker to access or manipulate files outside the intended directory by providing specially crafted paths as input to read or write sensitive data. This can occur when improperly handling user-supplied input in filesystem-related operations such as image uploads, dynamic content loading, and user file downloads.

An attacker could exploit a path traversal vulnerability to:

- Read sensitive files, including configuration files or other data containing credentials or encryption keys.
- Write files into restricted directories that enables code injection or privilege escalation.
- Download or delete critical system files.
- Gain access to user data and perform unauthorized actions.

## Example

```
# `params[:file][:path]` could contain "../../../../../etc/passwd", e.g.
send_file File.join('some', 'path', params[:file][:path])
```
## Pathname Confusion

`Pathname#join` has some confusing behavior: *any* absolute path segment (e.g. starting with `/`) causes the path to be absolute from that point.

Example:

```
> Pathname.new('a').join("a", "b", "/c", "d")
 => #<Pathname:/c/d>
```
Note that `Rails.root` is a `Pathname`.

Exercise extreme caution when passing user-provided input to this function.

### Additional Protections

Besides coding defensively, there are additional options for protecting against path traversal:

- Use the ActiveStorage module for handling uploaded files and store them in a service like S3, rather than storing user data on the same server or directory as the application.
- Configure permissions on the application server to disallow writing files or reading files outside of the application directory.
- Never include user-provided values in the file path or the file name.

A common pattern is to store files using application-generated file names, but keep a record of the user-provided name. When the user downloads the file, the `download` attribute and/or the Content Disposition header can be used to tell the browser the preferred name of the file, which can be the original user-provided name. Note that libraries like ActiveStorage will handle this for you.

However, be careful if users can download files named by *other* users. Overall, it is safer to generate file names from known-safe values.

Back to Warning Types

# Remote Code Execution

Brakeman reports on several cases of remote code execution, in which a user is able to control and execute code in ways unintended by application authors.

The obvious form of this is the use of `eval` with user input.

However, Brakeman also reports on dangerous uses of `send`, `constantize`, and other methods which allow creation of arbitrary objects or calling of arbitrary methods.

Back to Warning Types

# Remote Code Execution in YAML.Load

As seen in CVE-2013-0156, calling `YAML.load` with user input can lead to remote execution of arbitrary code. (To see a real point-and-fire exploit, see the Metasploit payload). While upgrading Rails, disabling XML parsing, or disabling YAML types in XML request parsing will fix the Rails vulnerability, manually passing user input to `YAML.load` remains unsafe.

For example:

```
#Do not do this!
YAML.load(params[:file])
```
Back to Warning Types

# Session Manipulation

Session manipulation can occur when an application allows user-input in session keys. Since sessions are typically considered a source of truth (e.g. to check the logged-in user or to match CSRF tokens), allowing an attacker to manipulate the session may lead to unintended behavior.

For example:

```
user_id = session[params[:name]]
current_user = User.find(user_id)
```
In this scenario, the attacker can point the `name` parameter to some other session value (for example, `_csrf_token`) that will be interpreted
as a user ID. If the ID matches an existing account, the attacker will now have access to that account.

To prevent this type of session manipulation, avoid using user-supplied input as session keys.

(See here for a tiny, self-contained challenge demonstrating this issue.)

Back to Warning Types

# Session Settings

### HTTP Only

It is recommended that session cookies be set to “http-only”. This helps prevent stealing of cookies via cross site scripting.

### Secret Length

Brakeman will warn if the key length for the session cookies is less than 30 characters.

### Version control inclusion

Brakeman will warn if the config/initializers/secret_token.rb is included in the version control. It is recommended that secret_token.rb is excluded from version control, and included in .gitignore

Back to Warning Types

# SQL Injection

Injection is #1 on the 2010 OWASP Top Ten web security risks. SQL injection is when a user is able to manipulate a value which is used unsafely inside a SQL query. This can lead to data leaks, data loss, elevation of privilege, and other unpleasant outcomes.

Brakeman focuses on ActiveRecord methods dealing with building SQL statements.

A basic (Rails 2.x) example looks like this:

```
User.first(:conditions => "username = '#{params[:username]}'")
```
Brakeman would produce a warning like this:

```
Possible SQL injection near line 30: User.first(:conditions => ("username = '#{params[:username]}'"))
```
The safe way to do this query is to use a parameterized query:

```
User.first(:conditions => ["username = ?", params[:username]])
```
Brakeman also understands the new Rails 3.x way of doing things (and local variables and concatentation):

```
username = params[:user][:name].downcase
password = params[:user][:password]
User.first.where("username = '" + username + "' AND password = '" + password + "'")
```
This results in this kind of warning:

```
Possible SQL injection near line 37:
User.first.where((((("username = '" + params[:user][:name].downcase) + "' AND password = '") + params[:user][:password]) + "'"))
```
See the Ruby Security Guide for more information and Rails-SQLi.org for many examples of SQL injection in Rails.

Back to Warning Types

# SSL Verification Bypass

Simply using SSL isn’t enough to ensure the data you are sending is secure. Man in the middle (MITM) attacks are well known and widely used. In some cases, these attacks rely on the client to establish a connection that doesn’t check the validity of the SSL certificate presented by the server. In this case, the attacker can present their own certificate and act as a man in the middle.

In Ruby, this happens when the OpenSSL verification mode is set to `VERIFY_NONE`

```
require "net/https"
require "uri"
uri = URI.parse("https://ssl-site.com/")
http = Net::HTTP.new(uri.host, uri.port)
http.use_ssl = true
http.verify_mode = OpenSSL::SSL::VERIFY_NONE
request = Net::HTTP::Get.new(uri.request_uri)
response = http.request(request)
```
In this case, if an invalid certificate was presented, no verification would occur, providing an opportunity for attack. When successful, the data transmitted (cookies, request parameters, POST bodies, etc.) would all be able to be intercepted by the MITM.

Brakeman would produce a warning like this:

```
SSL certificate verification was bypassed near line 24: http.verify_mode = OpenSSL::SSL::VERIFY_NONE
```
To ensure that SSL verification happens use the following mode:

```
http.verify_mode = OpenSSL::SSL::VERIFY_PEER
```
If the server certificate is invalid or context.ca_file is not set when verifying peers an OpenSSL::SSL::SSLError will be raised.

For more information on the impact of this issue, see the paper The Most Dangerous Code in the World.

Back to Warning Types

# Unmaintained Dependencies

Unmaintained or “end-of-life” dependencies can present security risks to your application.

When a dependency is no longer maintained, its developers may not release new versions with security patches for known vulnerabilities. This means that any known vulnerabilities in the dependency remain unpatched, leaving your application open to attacks.

In addition to known vulnerabilities, older versions of software are likely to receive less scrutiny are more likely to contain vulnerabilities that are not published and do not receive any public attention.

Maintained dependencies are also more likely to follow security best practices, such as using secure coding practices, regularly testing for vulnerabilities, and providing timely security patches. Outdated dependencies may not follow these best practices, increasing the risk of security vulnerabilities.

As a library ages, it is more likely to be completely abandoned or forgotten by its creator. Abandoned libraries may become target for supply chain attacks, where the attacker takes over an old code repository or an account on a package management server (such as RubyGems) and publishes a malicious version of the software.

## Ruby and Rails

Ruby versions are generally maintained for 3 years and 3 months after release. Check the listing of maintenance branches for more information.

Rails is more complicated, but generally only the current series and the last of the previous series is supported. See the Rails Maintenance Policy.

Back to Warning Types

# Unsafe Deserialization

Objects in Ruby may be serialized to strings. The main method for doing so is the built-in `Marshal` class. The `YAML`, `JSON`, and `CSV` libraries also have methods for dumping Ruby objects into strings, and then creating objects from those strings.

Deserialization of arbitrary objects can lead to remote code execution, as was demonstrated with CVE-2013-0156.

Brakeman warns when loading user input with `Marshal`, `YAML`, or `CSV`. `JSON` is covered by the checks for CVE-2013-0333

Back to Warning Types

# Unscoped Find

Unscoped `find` (and related methods) are a form of Direct Object Reference. Models which belong to another model should typically be accessed via a scoped query.

For example, if an `Account` belongs to a `User`, then this may be an unsafe unscoped find:

```
Account.find(params[:id])
```
Depending on the action, this could allow an attacker to access any account they wish.

Instead, it should be scoped to the currently logged-in user:

```
current_user = User.find(session[:user_id])
current_user.accounts.find(params[:id])
```
Back to Warning Types

# Redirect

Unvalidated redirects and forwards are #10 on the OWASP Top Ten.

Redirects which rely on user-supplied values can be used to “spoof” websites or hide malicious links in otherwise harmless-looking URLs. They can also allow access to restricted areas of a site if the destination is not validated.

Brakeman will raise warnings whenever `redirect_to` appears to be used with a user-supplied value that may allow them to change the `:host` option.

For example,

```
redirect_to params.merge(:action => :home)
```
will create a warning like

```
Possible unprotected redirect near line 46: redirect_to(params)
```
This is because `params` could contain `:host => 'evilsite.com'` which would redirect away from your site and to a malicious site.

If the first argument to `redirect_to` is a hash, then adding `:only_path => true` will limit the redirect to the current host. Another option is to specify the host explicitly.

```
redirect_to params.merge(:only_path => true)
redirect_to params.merge(:host => 'myhost.com')
```
If the first argument is a string, then it is possible to parse the string and extract the path:

```
redirect_to URI.parse(some_url).path
```
**If the URL does not contain a protocol (e.g., http://), then you will probably get unexpected results, as redirect_to will prepend the current host name and a protocol.**

### Rails 7 Updates

If `config.action_controller.raise_on_open_redirects` is `true` (default for *new* Rails 7.0 applications), then Rails will not allow redirecting to a domain that differs from the request.

Even if the configuration setting is not `true`, the protection can be applied by setting `allow_other_host: false` explicitly:

```
redirect_to params[:url], allow_other_host: false
```
The code above will raise an exception if `params[:url]` does not match the current domain.

Brakeman will warn about calls where `allow_other_host` is set to `true`.

To coerce the URL to be “safe”, use `url_from`:

```
redirect_to url_from(params[:url])
```
If the URL is does not match the current domain, then `url_from` returns `false`. The recommended pattern is to provide a fallback:

```
redirect_to url_from(params[:url]) || some_safe_default_url
```
Back to Warning Types

# Weak Hash

Brakeman reports a “Weak Hash” warning when it finds uses of hashing algorithms that should not be used for security-sensitive contexts such as hashing passwords or generating signatures.

Currently, Brakeman warns about the use of SHA1 and MD5, which should not be used for anything outside of interacting with Git.

The confidence level of a “Weak Hash” warning is based on whether the value being hash looks like user-controlled input or a password.

Back to Warning Types
