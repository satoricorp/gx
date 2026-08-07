# Test Plugins

Bandit supports many different tests to detect various security issues in python code. These tests are created as plugins and new ones can be created to extend the functionality offered by bandit today.

## Writing Tests

- To write a test:
- Identify a vulnerability to build a test for, and create a new file in examples/ that contains one or more cases of that vulnerability.
- Create a new Python source file to contain your test, you can reference existing tests for examples.
- Consider the vulnerability you’re testing for, mark the function with one or more of the appropriate decorators:
 - @checks(‘Call’)
- @checks(‘Import’, ‘ImportFrom’)
- @checks(‘Str’)
 - Register your plugin using the bandit.plugins entry point, see example.
- The function that you create should take a parameter “context” which is an instance of the context class you can query for information about the current element being examined. You can also get the raw AST node for more advanced use cases. Please see the context.py file for more.
- Extend your Bandit configuration file as needed to support your new test.
- Execute Bandit against the test file you defined in examples/ and ensure that it detects the vulnerability. Consider variations on how this vulnerability might present itself and extend the example file and the test function accordingly.

## Config Generation

In Bandit 1.0+ config files are optional. Plugins that need config settings are required to implement a module global gen_config function. This function is called with a single parameter, the test plugin name. It should return a dictionary with keys being the config option names and values being the default settings for each option. An example gen_config might look like the following:

```
def gen_config(name):
 if name == 'try_except_continue':
 return {'check_typed_exception': False}
```
When no config file is specified, or when the chosen file has no section pertaining to a given plugin, gen_config will be called to provide defaults.

The config file generation tool bandit-config-generator will also call gen_config on all discovered plugins to produce template config blocks. If the defaults are acceptable then these blocks may be deleted to create a minimal configuration, or otherwise edited as needed. The above example would produce the following config snippet.

```
try_except_continue: {check_typed_exception: false}
```
## Example Test Plugin

```
@bandit.checks('Call')
def prohibit_unsafe_deserialization(context):
 if 'unsafe_load' in context.call_function_name_qual:
 return bandit.Issue(
 severity=bandit.HIGH,
 confidence=bandit.HIGH,
 text="Unsafe deserialization detected."
 )
```
To register your plugin, you have two options:

- If you’re using setuptools directly, add something like the following to your setup call: - # If you have an imaginary bson formatter in the bandit_bson module # and a function called `formatter`. entry_points={'bandit.formatters': ['bson = bandit_bson:formatter']} # Or a check for using mako templates in bandit_mako that entry_points={'bandit.plugins': ['mako = bandit_mako']}
- If you’re using pbr, add something like the following to your setup.cfg file: - [entry_points] bandit.formatters = bson = bandit_bson:formatter bandit.plugins = mako = bandit_mako

## Plugin ID Groupings

| ID | Description |
|---|---|
| B1xx | misc tests |
| B2xx | application/framework misconfiguration |
| B3xx | blacklists (calls) |
| B4xx | blacklists (imports) |
| B5xx | cryptography |
| B6xx | injection |
| B7xx | XSS |

## Complete Test Plugin Listing

- B101: assert_used
- B102: exec_used
- B103: set_bad_file_permissions
- B104: hardcoded_bind_all_interfaces
- B105: hardcoded_password_string
- B106: hardcoded_password_funcarg
- B107: hardcoded_password_default
- B108: hardcoded_tmp_directory
- B109: password_config_option_not_marked_secret
- B110: try_except_pass
- B111: execute_with_run_as_root_equals_true
- B112: try_except_continue
- B113: request_without_timeout
- B201: flask_debug_true
- B202: tarfile_unsafe_members
- B324: hashlib
- B501: request_with_no_cert_validation
- B502: ssl_with_bad_version
- B503: ssl_with_bad_defaults
- B504: ssl_with_no_version
- B505: weak_cryptographic_key
- B506: yaml_load
- B507: ssh_no_host_key_verification
- B508: snmp_insecure_version
- B509: snmp_weak_cryptography
- B601: paramiko_calls
- B602: subprocess_popen_with_shell_equals_true
- B603: subprocess_without_shell_equals_true
- B604: any_other_function_with_shell_equals_true
- B605: start_process_with_a_shell
- B606: start_process_with_no_shell
- B607: start_process_with_partial_path
- B608: hardcoded_sql_expressions
- B609: linux_commands_wildcard_injection
- B610: django_extra_used
- B611: django_rawsql_used
- B612: logging_config_insecure_listen
- B613: trojansource
- B614: pytorch_load
- B615: huggingface_unsafe_download
- B701: jinja2_autoescape_false
- B702: use_of_mako_templates
- B703: django_mark_safe
- B704: markupsafe_markup_xss

# B101: assert_used

## B101: Test for use of assert

This plugin test checks for the use of the Python `assert` keyword. It was
discovered that some projects used assert to enforce interface constraints.
However, assert is removed with compiling to optimised byte code (python -O
producing *.opt-1.pyc files). This caused various protections to be removed.
Consider raising a semantically meaningful error or `AssertionError` instead.

Please see
https://docs.python.org/3/reference/simple_stmts.html#the-assert-statement for
more info on `assert`.

**Config Options:**

You can configure files that skip this check. This is often useful when you use assert statements in test cases.

```
assert_used:
 skips: ['*_test.py', '*test_*.py']
```
- Example:

```
>> Issue: Use of assert detected. The enclosed code will be removed when
 compiling to optimised byte code.
 Severity: Low Confidence: High
 CWE: CWE-703 (https://cwe.mitre.org/data/definitions/703.html)
 Location: ./examples/assert.py:1
1 assert logged_in
2 display_assets()
```
See also

Added in version 0.11.0.

Changed in version 1.7.3: CWE information added

# B102: exec_used

## B102: Test for the use of exec

This plugin test checks for the use of Python’s exec method or keyword. The Python docs succinctly describe why the use of exec is risky.

- Example:

```
>> Issue: Use of exec detected.
 Severity: Medium Confidence: High
 CWE: CWE-78 (https://cwe.mitre.org/data/definitions/78.html)
 Location: ./examples/exec.py:2
1 exec("do evil")
```
See also

Added in version 0.9.0.

Changed in version 1.7.3: CWE information added

# B103: set_bad_file_permissions

## B103: Test for setting permissive file permissions

POSIX based operating systems utilize a permissions model to protect access to
parts of the file system. This model supports three roles “owner”, “group”
and “world” each role may have a combination of “read”, “write” or “execute”
flags sets. Python provides `chmod` to manipulate POSIX style permissions.

This plugin test looks for the use of `chmod` and will alert when it is used
to set particularly permissive control flags. A MEDIUM warning is generated if
a file is set to group write or executable and a HIGH warning is reported if a
file is set world write or executable. Warnings are given with HIGH confidence.

- Example:

```
>> Issue: Probable insecure usage of temp file/directory.
 Severity: Medium Confidence: Medium
 CWE: CWE-732 (https://cwe.mitre.org/data/definitions/732.html)
 Location: ./examples/os-chmod.py:15
14 os.chmod('/etc/hosts', 0o777)
15 os.chmod('/tmp/oh_hai', 0x1ff)
16 os.chmod('/etc/passwd', stat.S_IRWXU)
>> Issue: Chmod setting a permissive mask 0777 on file (key_file).
 Severity: High Confidence: High
 CWE: CWE-732 (https://cwe.mitre.org/data/definitions/732.html)
 Location: ./examples/os-chmod.py:17
16 os.chmod('/etc/passwd', stat.S_IRWXU)
17 os.chmod(key_file, 0o777)
18
```
See also

Added in version 0.9.0.

Changed in version 1.7.3: CWE information added

Changed in version 1.7.5: Added checks for S_IWGRP and S_IXOTH

# B104: hardcoded_bind_all_interfaces

## B104: Test for binding to all interfaces

Binding to all network interfaces can potentially open up a service to traffic on unintended interfaces, that may not be properly documented or secured. This plugin test looks for a string pattern “0.0.0.0” that may indicate a hardcoded binding to all network interfaces.

- Example:

```
>> Issue: Possible binding to all interfaces.
 Severity: Medium Confidence: Medium
 CWE: CWE-605 (https://cwe.mitre.org/data/definitions/605.html)
 Location: ./examples/binding.py:4
3 s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
4 s.bind(('0.0.0.0', 31137))
5 s.bind(('192.168.0.1', 8080))
```
See also

Added in version 0.9.0.

Changed in version 1.7.3: CWE information added

# B105: hardcoded_password_string

-
bandit.plugins.general_hardcoded_password.hardcoded_password_string(*context*)[source]
- **B105: Test for use of hard-coded password strings**- The use of hard-coded passwords increases the possibility of password guessing tremendously. This plugin test looks for all string literals and checks the following conditions: - assigned to a variable that looks like a password
- assigned to a dict key that looks like a password
- assigned to a class attribute that looks like a password
- used in a comparison with a variable that looks like a password
 - Variables are considered to look like a password if they have match any one of: - “password”
- “pass”
- “passwd”
- “pwd”
- “secret”
- “token”
- “secrete”
 - Note: this can be noisy and may generate false positives. - **Config Options:**- None - Example:
 - `>> Issue: Possible hardcoded password '(root)' Severity: Low Confidence: Low CWE: CWE-259 (https://cwe.mitre.org/data/definitions/259.html) Location: ./examples/hardcoded-passwords.py:5 4 def someFunction2(password): 5 if password == "root": 6 print("OK, logged in")`- See also - Added in version 0.9.0. - Changed in version 1.7.3: CWE information added

# B106: hardcoded_password_funcarg

-
bandit.plugins.general_hardcoded_password.hardcoded_password_funcarg(*context*)[source]
- **B106: Test for use of hard-coded password function arguments**- The use of hard-coded passwords increases the possibility of password guessing tremendously. This plugin test looks for all function calls being passed a keyword argument that is a string literal. It checks that the assigned local variable does not look like a password. - Variables are considered to look like a password if they have match any one of: - “password”
- “pass”
- “passwd”
- “pwd”
- “secret”
- “token”
- “secrete”
 - Note: this can be noisy and may generate false positives. - **Config Options:**- None - Example:
 - >> Issue: [B106:hardcoded_password_funcarg] Possible hardcoded password: 'blerg' Severity: Low Confidence: Medium CWE: CWE-259 (https://cwe.mitre.org/data/definitions/259.html) Location: ./examples/hardcoded-passwords.py:16 15 16 doLogin(password="blerg") - See also - Added in version 0.9.0. - Changed in version 1.7.3: CWE information added

# B107: hardcoded_password_default

-
bandit.plugins.general_hardcoded_password.hardcoded_password_default(*context*)[source]
- **B107: Test for use of hard-coded password argument defaults**- The use of hard-coded passwords increases the possibility of password guessing tremendously. This plugin test looks for all function definitions that specify a default string literal for some argument. It checks that the argument does not look like a password. - Variables are considered to look like a password if they have match any one of: - “password”
- “pass”
- “passwd”
- “pwd”
- “secret”
- “token”
- “secrete”
 - Note: this can be noisy and may generate false positives. We do not report on None values which can be legitimately used as a default value, when initializing a function or class. - **Config Options:**- None - Example:
 - `>> Issue: [B107:hardcoded_password_default] Possible hardcoded password: 'Admin' Severity: Low Confidence: Medium CWE: CWE-259 (https://cwe.mitre.org/data/definitions/259.html) Location: ./examples/hardcoded-passwords.py:1 1 def someFunction(user, password="Admin"): 2 print("Hi " + user)`- See also - Added in version 0.9.0. - Changed in version 1.7.3: CWE information added

# B108: hardcoded_tmp_directory

## B108: Test for insecure usage of tmp file/directory

Safely creating a temporary file or directory means following a number of rules (see the references for more details). This plugin test looks for strings starting with (configurable) commonly used temporary paths, for example:

/tmp

/var/tmp

/dev/shm

**Config Options:**

This test plugin takes a similarly named config block, hardcoded_tmp_directory. The config block provides a Python list, tmp_dirs, that lists string fragments indicating possible temporary file paths. Any string starting with one of these fragments will report a MEDIUM confidence issue.

```
hardcoded_tmp_directory:
 tmp_dirs: ['/tmp', '/var/tmp', '/dev/shm']
```
- Example:

See also

Added in version 0.9.0.

Changed in version 1.7.3: CWE information added

# B109: password_config_option_not_marked_secret

This plugin has been removed.

B109: Test for a password based config option not marked secret

Passwords are sensitive and must be protected appropriately. In OpenStack Oslo there is an option to mark options “secret” which will ensure that they are not logged. This plugin detects usages of oslo configuration functions that appear to deal with strings ending in ‘password’ and flag usages where they have not been marked secret.

If such a value is found a MEDIUM severity error is generated. If ‘False’ or ‘None’ are explicitly set, Bandit will return a MEDIUM confidence issue. If Bandit can’t determine the value of secret it will return a LOW confidence issue.

**Config Options:**

```
password_config_option_not_marked_secret:
 function_names:
 - oslo.config.cfg.StrOpt
 - oslo_config.cfg.StrOpt
```
- Example:

```
>> Issue: [password_config_option_not_marked_secret] oslo config option
possibly not marked secret=True identified.
 Severity: Medium Confidence: Low
 Location: examples/secret-config-option.py:12
11 help="User's password"),
12 cfg.StrOpt('nova_password',
13 secret=secret,
14 help="Nova user password"),
15 ]
>> Issue: [password_config_option_not_marked_secret] oslo config option not
marked secret=True identified, security issue.
 Severity: Medium Confidence: Medium
 Location: examples/secret-config-option.py:21
20 help="LDAP ubind ser name"),
21 cfg.StrOpt('ldap_password',
22 help="LDAP bind user password"),
23 cfg.StrOpt('ldap_password_attribute',
```
Added in version 0.10.0.

Deprecated since version 1.5.0: This plugin was removed

# B110: try_except_pass

## B110: Test for a pass in the except block

Errors in Python code bases are typically communicated using `Exceptions`.
An exception object is ‘raised’ in the event of an error and can be ‘caught’ at
a later point in the program, typically some error handling or logging action
will then be performed.

However, it is possible to catch an exception and silently ignore it. This is illustrated with the following example

```
try:
 do_some_stuff()
except Exception:
 pass
```
This pattern is considered bad practice in general, but also represents a potential security issue. A larger than normal volume of errors from a service can indicate an attempt is being made to disrupt or interfere with it. Thus errors should, at the very least, be logged.

There are rare situations where it is desirable to suppress errors, but this is
typically done with specific exception types, rather than the base Exception
class (or no type). To accommodate this, the test may be configured to ignore
‘try, except, pass’ where the exception is typed. For example, the following
would not generate a warning if the configuration option
`checked_typed_exception` is set to False:

```
try:
 do_some_stuff()
except ZeroDivisionError:
 pass
```
**Config Options:**

```
try_except_pass:
 check_typed_exception: True
```
- Example:

```
>> Issue: Try, Except, Pass detected.
 Severity: Low Confidence: High
 CWE: CWE-703 (https://cwe.mitre.org/data/definitions/703.html)
 Location: ./examples/try_except_pass.py:4
3 a = 1
4 except:
5 pass
```
Added in version 0.13.0.

Changed in version 1.7.3: CWE information added

# B111: execute_with_run_as_root_equals_true

This plugin has been removed.

B111: Test for the use of rootwrap running as root

Running commands as root dramatically increase their potential risk. Running commands with restricted user privileges provides defense in depth against command injection attacks, or developer and configuration error. This plugin test checks for specific methods being called with a keyword parameter run_as_root set to True, a common OpenStack idiom.

**Config Options:**

This test plugin takes a similarly named configuration block, execute_with_run_as_root_equals_true, providing a list, function_names, of function names. A call to any of these named functions will be checked for a run_as_root keyword parameter, and if True, will report a Low severity issue.

```
execute_with_run_as_root_equals_true:
 function_names:
 - ceilometer.utils.execute
 - cinder.utils.execute
 - neutron.agent.linux.utils.execute
 - nova.utils.execute
 - nova.utils.trycmd
```
- Example:

```
>> Issue: Execute with run_as_root=True identified, possible security
 issue.
 Severity: Low Confidence: Medium
 Location: ./examples/exec-as-root.py:26
25 nova_utils.trycmd('gcc --version')
26 nova_utils.trycmd('gcc --version', run_as_root=True)
27
```
See also

Added in version 0.10.0.

Deprecated since version 1.5.0: This plugin was removed

# B112: try_except_continue

## B112: Test for a continue in the except block

Errors in Python code bases are typically communicated using `Exceptions`.
An exception object is ‘raised’ in the event of an error and can be ‘caught’ at
a later point in the program, typically some error handling or logging action
will then be performed.

However, it is possible to catch an exception and silently ignore it while in a loop. This is illustrated with the following example

```
while keep_going:
 try:
 do_some_stuff()
 except Exception:
 continue
```
This pattern is considered bad practice in general, but also represents a potential security issue. A larger than normal volume of errors from a service can indicate an attempt is being made to disrupt or interfere with it. Thus errors should, at the very least, be logged.

There are rare situations where it is desirable to suppress errors, but this is
typically done with specific exception types, rather than the base Exception
class (or no type). To accommodate this, the test may be configured to ignore
‘try, except, continue’ where the exception is typed. For example, the
following would not generate a warning if the configuration option
`checked_typed_exception` is set to False:

```
while keep_going:
 try:
 do_some_stuff()
 except ZeroDivisionError:
 continue
```
**Config Options:**

```
try_except_continue:
 check_typed_exception: True
```
- Example:

```
>> Issue: Try, Except, Continue detected.
 Severity: Low Confidence: High
 CWE: CWE-703 (https://cwe.mitre.org/data/definitions/703.html)
 Location: ./examples/try_except_continue.py:5
4 a = i
5 except:
6 continue
```
Added in version 1.0.0.

Changed in version 1.7.3: CWE information added

# B113: request_without_timeout

## B113: Test for missing requests timeout

This plugin test checks for `requests` or `httpx` calls without a timeout
specified.

Nearly all production code should use this parameter in nearly all requests, Failure to do so can cause your program to hang indefinitely.

When request methods are used without the timeout parameter set, Bandit will return a MEDIUM severity error.

- Example:

```
>> Issue: [B113:request_without_timeout] Call to requests without timeout
 Severity: Medium Confidence: Low
 CWE: CWE-400 (https://cwe.mitre.org/data/definitions/400.html)
 More Info: https://bandit.readthedocs.io/en/latest/plugins/b113_request_without_timeout.html
 Location: examples/requests-missing-timeout.py:3:0
2
3 requests.get('https://gmail.com')
4 requests.get('https://gmail.com', timeout=None)
--------------------------------------------------
>> Issue: [B113:request_without_timeout] Call to requests with timeout set to None
 Severity: Medium Confidence: Low
 CWE: CWE-400 (https://cwe.mitre.org/data/definitions/400.html)
 More Info: https://bandit.readthedocs.io/en/latest/plugins/b113_request_without_timeout.html
 Location: examples/requests-missing-timeout.py:4:0
3 requests.get('https://gmail.com')
4 requests.get('https://gmail.com', timeout=None)
5 requests.get('https://gmail.com', timeout=5)
```
Added in version 1.7.5.

Changed in version 1.7.10: Added check for httpx module

# B201: flask_debug_true

## B201: Test for use of flask app with debug set to true

Running Flask applications in debug mode results in the Werkzeug debugger being enabled. This includes a feature that allows arbitrary code execution. Documentation for both Flask [1] and Werkzeug [2] strongly suggests that debug mode should never be enabled on production systems.

Operating a production server with debug mode enabled was the probable cause of the Patreon breach in 2015 [3].

- Example:

```
>> Issue: A Flask app appears to be run with debug=True, which exposes
the Werkzeug debugger and allows the execution of arbitrary code.
 Severity: High Confidence: High
 CWE: CWE-94 (https://cwe.mitre.org/data/definitions/94.html)
 Location: examples/flask_debug.py:10
9 #bad
10 app.run(debug=True)
11
```
See also

Added in version 0.15.0.

Changed in version 1.7.3: CWE information added

# B202: tarfile_unsafe_members

## B202: Test for tarfile.extractall

This plugin will look for usage of `tarfile.extractall()`

Severity are set as follows:

- `tarfile.extractall(members=function(tarfile))`- LOW
- `tarfile.extractall(members=?)`- member is not a function - MEDIUM
- `tarfile.extractall()`- members from the archive is trusted - HIGH

Use `tarfile.extractall(members=function_name)` and define a function
that will inspect each member. Discard files that contain a directory
traversal sequences such as `../` or `\..` along with all special filetypes
unless you explicitly need them.

- Example:

```
>> Issue: [B202:tarfile_unsafe_members] tarfile.extractall used without
any validation. You should check members and discard dangerous ones
Severity: High Confidence: High
CWE: CWE-22 (https://cwe.mitre.org/data/definitions/22.html)
Location: examples/tarfile_extractall.py:8
More Info:
https://bandit.readthedocs.io/en/latest/plugins/b202_tarfile_unsafe_members.html
7 tar = tarfile.open(filename)
8 tar.extractall(path=tempfile.mkdtemp())
9 tar.close()
```
See also

Added in version 1.7.5.

Changed in version 1.7.8: Added check for filter parameter

# B324: hashlib

## B324: Test use of insecure md4, md5, or sha1 hash functions in hashlib

This plugin checks for the usage of the insecure MD4, MD5, or SHA1 hash
functions in `hashlib` and `crypt`. The `hashlib.new` function provides
the ability to construct a new hashing object using the named algorithm. This
can be used to create insecure hash functions like MD4 and MD5 if they are
passed as algorithm names to this function.

This check does additional checking for usage of keyword usedforsecurity on all function variations of hashlib.

Similar to `hashlib`, this plugin also checks for usage of one of the
`crypt` module’s weak hashes. `crypt` also permits MD5 among other weak
hash variants.

- Example:

```
>> Issue: [B324:hashlib] Use of weak MD4, MD5, or SHA1 hash for
 security. Consider usedforsecurity=False
 Severity: High Confidence: High
 CWE: CWE-327 (https://cwe.mitre.org/data/definitions/327.html)
 Location: examples/hashlib_new_insecure_functions.py:3:0
 More Info: https://bandit.readthedocs.io/en/latest/plugins/b324_hashlib.html
2
3 hashlib.new('md5')
4
```
Added in version 1.5.0.

Changed in version 1.7.3: CWE information added

Changed in version 1.7.6: Added check for the crypt module weak hashes

# B501: request_with_no_cert_validation

## B501: Test for missing certificate validation

Encryption in general is typically critical to the security of many applications. Using TLS can greatly increase security by guaranteeing the identity of the party you are communicating with. This is accomplished by one or both parties presenting trusted certificates during the connection initialization phase of TLS.

When HTTPS request methods are used, certificates are validated automatically which is the desired behavior. If certificate validation is explicitly turned off Bandit will return a HIGH severity error.

- Example:

```
>> Issue: [request_with_no_cert_validation] Call to requests with
verify=False disabling SSL certificate checks, security issue.
 Severity: High Confidence: High
 CWE: CWE-295 (https://cwe.mitre.org/data/definitions/295.html)
 Location: examples/requests-ssl-verify-disabled.py:4
3 requests.get('https://gmail.com', verify=True)
4 requests.get('https://gmail.com', verify=False)
5 requests.post('https://gmail.com', verify=True)
```
See also

Added in version 0.9.0.

Changed in version 1.7.3: CWE information added

Changed in version 1.7.5: Added check for httpx module

# B502: ssl_with_bad_version

-
bandit.plugins.insecure_ssl_tls.ssl_with_bad_version(*context*,*config*)[source]
- **B502: Test for SSL use with bad version used**- Several highly publicized exploitable flaws have been discovered in all versions of SSL and early versions of TLS. It is strongly recommended that use of the following known broken protocol versions be avoided: - SSL v2
- SSL v3
- TLS v1
- TLS v1.1
 - This plugin test scans for calls to Python methods with parameters that indicate the used broken SSL/TLS protocol versions. Currently, detection supports methods using Python’s native SSL/TLS support and the pyOpenSSL module. A HIGH severity warning will be reported whenever known broken protocol versions are detected. - It is worth noting that native support for TLS 1.2 is only available in more recent Python versions, specifically 2.7.9 and up, and 3.x - A note on ‘SSLv23’: - Amongst the available SSL/TLS versions provided by Python/pyOpenSSL there exists the option to use SSLv23. This very poorly named option actually means “use the highest version of SSL/TLS supported by both the server and client”. This may (and should be) a version well in advance of SSL v2 or v3. Bandit can scan for the use of SSLv23 if desired, but its detection does not necessarily indicate a problem. - When using SSLv23 it is important to also provide flags to explicitly exclude bad versions of SSL/TLS from the protocol versions considered. Both the Python native and pyOpenSSL modules provide the - `OP_NO_SSLv2`and- `OP_NO_SSLv3`flags for this purpose.- **Config Options:**- ssl_with_bad_version: bad_protocol_versions: - PROTOCOL_SSLv2 - SSLv2_METHOD - SSLv23_METHOD - PROTOCOL_SSLv3 # strict option - PROTOCOL_TLSv1 # strict option - SSLv3_METHOD # strict option - TLSv1_METHOD # strict option - Example:
 - >> Issue: ssl.wrap_socket call with insecure SSL/TLS protocol version identified, security issue. Severity: High Confidence: High CWE: CWE-327 (https://cwe.mitre.org/data/definitions/327.html) Location: ./examples/ssl-insecure-version.py:13 12 # strict tests 13 ssl.wrap_socket(ssl_version=ssl.PROTOCOL_SSLv3) 14 ssl.wrap_socket(ssl_version=ssl.PROTOCOL_TLSv1) - See also - `ssl_with_bad_defaults()`
- `ssl_with_no_version()`
- https://security.openstack.org/guidelines/dg_move-data-securely.html
 - Added in version 0.9.0. - Changed in version 1.7.3: CWE information added - Changed in version 1.7.5: Added TLS 1.1

# B503: ssl_with_bad_defaults

-
bandit.plugins.insecure_ssl_tls.ssl_with_bad_defaults(*context*,*config*)[source]
- **B503: Test for SSL use with bad defaults specified**- This plugin is part of a family of tests that detect the use of known bad versions of SSL/TLS, please see ../plugins/ssl_with_bad_version for a complete discussion. Specifically, this plugin test scans for Python methods with default parameter values that specify the use of broken SSL/TLS protocol versions. Currently, detection supports methods using Python’s native SSL/TLS support and the pyOpenSSL module. A MEDIUM severity warning will be reported whenever known broken protocol versions are detected. - **Config Options:**- This test shares the configuration provided for the standard ../plugins/ssl_with_bad_version test, please refer to its documentation. - Example:
 - >> Issue: Function definition identified with insecure SSL/TLS protocol version by default, possible security issue. Severity: Medium Confidence: Medium CWE: CWE-327 (https://cwe.mitre.org/data/definitions/327.html) Location: ./examples/ssl-insecure-version.py:28 27 28 def open_ssl_socket(version=SSL.SSLv2_METHOD): 29 pass - See also - `ssl_with_bad_version()`
- `ssl_with_no_version()`
- https://security.openstack.org/guidelines/dg_move-data-securely.html
 - Added in version 0.9.0. - Changed in version 1.7.3: CWE information added - Changed in version 1.7.5: Added TLS 1.1

# B504: ssl_with_no_version

-
bandit.plugins.insecure_ssl_tls.ssl_with_no_version(*context*)[source]
- **B504: Test for SSL use with no version specified**- This plugin is part of a family of tests that detect the use of known bad versions of SSL/TLS, please see ../plugins/ssl_with_bad_version for a complete discussion. Specifically, This plugin test scans for specific methods in Python’s native SSL/TLS support and the pyOpenSSL module that configure the version of SSL/TLS protocol to use. These methods are known to provide default value that maximize compatibility, but permit use of the aforementioned broken protocol versions. A LOW severity warning will be reported whenever this is detected. - **Config Options:**- This test shares the configuration provided for the standard ../plugins/ssl_with_bad_version test, please refer to its documentation. - Example:
 - >> Issue: ssl.wrap_socket call with no SSL/TLS protocol version specified, the default SSLv23 could be insecure, possible security issue. Severity: Low Confidence: Medium CWE: CWE-327 (https://cwe.mitre.org/data/definitions/327.html) Location: ./examples/ssl-insecure-version.py:23 22 23 ssl.wrap_socket() 24 - See also - `ssl_with_bad_version()`
- `ssl_with_bad_defaults()`
- https://security.openstack.org/guidelines/dg_move-data-securely.html
 - Added in version 0.9.0. - Changed in version 1.7.3: CWE information added

# B505: weak_cryptographic_key

## B505: Test for weak cryptographic key use

As computational power increases, so does the ability to break ciphers with smaller key lengths. The recommended key length size for RSA and DSA algorithms is 2048 and higher. 1024 bits and below are now considered breakable. EC key length sizes are recommended to be 224 and higher with 160 and below considered breakable. This plugin test checks for use of any key less than those limits and returns a high severity error if lower than the lower threshold and a medium severity error for those lower than the higher threshold.

- Example:

```
>> Issue: DSA key sizes below 1024 bits are considered breakable.
 Severity: High Confidence: High
 CWE: CWE-326 (https://cwe.mitre.org/data/definitions/326.html)
 Location: examples/weak_cryptographic_key_sizes.py:36
35 # Also incorrect: without keyword args
36 dsa.generate_private_key(512,
37 backends.default_backend())
38 rsa.generate_private_key(3,
```
See also

Added in version 0.14.0.

Changed in version 1.7.3: CWE information added

# B506: yaml_load

## B506: Test for use of yaml load

This plugin test checks for the unsafe usage of the `yaml.load` function from
the PyYAML package. The yaml.load function provides the ability to construct
an arbitrary Python object, which may be dangerous if you receive a YAML
document from an untrusted source. The function yaml.safe_load limits this
ability to simple Python objects like integers or lists.

Please see
https://pyyaml.org/wiki/PyYAMLDocumentation#LoadingYAML for more information
on `yaml.load` and yaml.safe_load

- Example:

```
>> Issue: [yaml_load] Use of unsafe yaml load. Allows instantiation of
 arbitrary objects. Consider yaml.safe_load().
 Severity: Medium Confidence: High
 CWE: CWE-20 (https://cwe.mitre.org/data/definitions/20.html)
 Location: examples/yaml_load.py:5
4 ystr = yaml.dump({'a' : 1, 'b' : 2, 'c' : 3})
5 y = yaml.load(ystr)
6 yaml.dump(y)
```
See also

Added in version 1.0.0.

Changed in version 1.7.3: CWE information added

# B507: ssh_no_host_key_verification

## B507: Test for missing host key validation

Encryption in general is typically critical to the security of many applications. Using SSH can greatly increase security by guaranteeing the identity of the party you are communicating with. This is accomplished by one or both parties presenting trusted host keys during the connection initialization phase of SSH.

When paramiko methods are used, host keys are verified by default. If host key verification is disabled, Bandit will return a HIGH severity error.

- Example:

```
>> Issue: [B507:ssh_no_host_key_verification] Paramiko call with policy set
to automatically trust the unknown host key.
Severity: High Confidence: Medium
CWE: CWE-295 (https://cwe.mitre.org/data/definitions/295.html)
Location: examples/no_host_key_verification.py:4
3 ssh_client = client.SSHClient()
4 ssh_client.set_missing_host_key_policy(client.AutoAddPolicy)
5 ssh_client.set_missing_host_key_policy(client.WarningPolicy)
```
Added in version 1.5.1.

Changed in version 1.7.3: CWE information added

# B508: snmp_insecure_version

-
bandit.plugins.snmp_security_check.snmp_insecure_version_check(*context*)[source]
- **B508: Checking for insecure SNMP versions**- This test is for checking for the usage of insecure SNMP version like
- v1, v2c
 - Please update your code to use more secure versions of SNMP. - Example:
 - `>> Issue: [B508:snmp_insecure_version_check] The use of SNMPv1 and SNMPv2 is insecure. You should use SNMPv3 if able. Severity: Medium Confidence: High CWE: CWE-319 (https://cwe.mitre.org/data/definitions/319.html) Location: examples/snmp.py:4:4 More Info: https://bandit.readthedocs.io/en/latest/plugins/b508_snmp_insecure_version_check.html 3 # SHOULD FAIL 4 a = CommunityData('public', mpModel=0) 5 # SHOULD FAIL`- See also - Added in version 1.7.2. - Changed in version 1.7.3: CWE information added

# B509: snmp_weak_cryptography

-
bandit.plugins.snmp_security_check.snmp_crypto_check(*context*)[source]
- **B509: Checking for weak cryptography**- This test is for checking for the usage of insecure SNMP cryptography:
- v3 using noAuthNoPriv.
 - Please update your code to use more secure versions of SNMP. For example: - Instead of:
- CommunityData(‘public’, mpModel=0)
- Use (Defaults to usmHMACMD5AuthProtocol and usmDESPrivProtocol
- UsmUserData(“securityName”, “authName”, “privName”)
 - Example:
 - `>> Issue: [B509:snmp_crypto_check] You should not use SNMPv3 without encryption. noAuthNoPriv & authNoPriv is insecure Severity: Medium CWE: CWE-319 (https://cwe.mitre.org/data/definitions/319.html) Confidence: High Location: examples/snmp.py:6:11 More Info: https://bandit.readthedocs.io/en/latest/plugins/b509_snmp_crypto_check.html 5 # SHOULD FAIL 6 insecure = UsmUserData("securityName") 7 # SHOULD FAIL`- See also - Added in version 1.7.2. - Changed in version 1.7.3: CWE information added

# B601: paramiko_calls

## B601: Test for shell injection within Paramiko

Paramiko is a Python library designed to work with the SSH2 protocol for secure (encrypted and authenticated) connections to remote machines. It is intended to run commands on a remote host. These commands are run within a shell on the target and are thus vulnerable to various shell injection attacks. Bandit reports a MEDIUM issue when it detects the use of Paramiko’s “exec_command” method advising the user to check inputs are correctly sanitized.

- Example:

```
>> Issue: Possible shell injection via Paramiko call, check inputs are
 properly sanitized.
 Severity: Medium Confidence: Medium
 CWE: CWE-78 (https://cwe.mitre.org/data/definitions/78.html)
 Location: ./examples/paramiko_injection.py:4
3 # this is not safe
4 paramiko.exec_command('something; really; unsafe')
5
```
See also

Added in version 0.12.0.

Changed in version 1.7.3: CWE information added

# B602: subprocess_popen_with_shell_equals_true

-
bandit.plugins.injection_shell.subprocess_popen_with_shell_equals_true(*context*,*config*)[source]
- **B602: Test for use of popen with shell equals true**- Python possesses many mechanisms to invoke an external executable. However, doing so may present a security issue if appropriate care is not taken to sanitize any user provided or variable input. - This plugin test is part of a family of tests built to check for process spawning and warn appropriately. Specifically, this test looks for the spawning of a subprocess using a command shell. This type of subprocess invocation is dangerous as it is vulnerable to various shell injection attacks. Great care should be taken to sanitize all input in order to mitigate this risk. Calls of this type are identified by a parameter of ‘shell=True’ being given. - Additionally, this plugin scans the command string given and adjusts its reported severity based on how it is presented. If the command string is a simple static string containing no special shell characters, then the resulting issue has low severity. If the string is static, but contains shell formatting characters or wildcards, then the reported issue is medium. Finally, if the string is computed using Python’s string manipulation or formatting operations, then the reported issue has high severity. These severity levels reflect the likelihood that the code is vulnerable to injection. - See also: - ../plugins/linux_commands_wildcard_injection
- ../plugins/subprocess_without_shell_equals_true
- ../plugins/start_process_with_no_shell
- ../plugins/start_process_with_a_shell
- ../plugins/start_process_with_partial_path
 - **Config Options:**- This plugin test shares a configuration with others in the same family, namely shell_injection. This configuration is divided up into three sections, subprocess, shell and no_shell. They each list Python calls that spawn subprocesses, invoke commands within a shell, or invoke commands without a shell (by replacing the calling process) respectively. - This plugin specifically scans for methods listed in subprocess section that have shell=True specified. - shell_injection: # Start a process using the subprocess module, or one of its wrappers. subprocess: - subprocess.Popen - subprocess.call - Example:
 - `>> Issue: subprocess call with shell=True seems safe, but may be changed in the future, consider rewriting without shell Severity: Low Confidence: High CWE: CWE-78 (https://cwe.mitre.org/data/definitions/78.html) Location: ./examples/subprocess_shell.py:21 20 subprocess.check_call(['/bin/ls', '-l'], shell=False) 21 subprocess.check_call('/bin/ls -l', shell=True) 22 >> Issue: call with shell=True contains special shell characters, consider moving extra logic into Python code Severity: Medium Confidence: High CWE: CWE-78 (https://cwe.mitre.org/data/definitions/78.html) Location: ./examples/subprocess_shell.py:26 25 26 subprocess.Popen('/bin/ls *', shell=True) 27 subprocess.Popen('/bin/ls %s' % ('something',), shell=True) >> Issue: subprocess call with shell=True identified, security issue. Severity: High Confidence: High CWE: CWE-78 (https://cwe.mitre.org/data/definitions/78.html) Location: ./examples/subprocess_shell.py:27 26 subprocess.Popen('/bin/ls *', shell=True) 27 subprocess.Popen('/bin/ls %s' % ('something',), shell=True) 28 subprocess.Popen('/bin/ls {}'.format('something'), shell=True)`- See also - Added in version 0.9.0. - Changed in version 1.7.3: CWE information added

# B603: subprocess_without_shell_equals_true

-
bandit.plugins.injection_shell.subprocess_without_shell_equals_true(*context*,*config*)[source]
- **B603: Test for use of subprocess without shell equals true**- Python possesses many mechanisms to invoke an external executable. However, doing so may present a security issue if appropriate care is not taken to sanitize any user provided or variable input. - This plugin test is part of a family of tests built to check for process spawning and warn appropriately. Specifically, this test looks for the spawning of a subprocess without the use of a command shell. This type of subprocess invocation is not vulnerable to shell injection attacks, but care should still be taken to ensure validity of input. - Because this is a lesser issue than that described in subprocess_popen_with_shell_equals_true a LOW severity warning is reported. - See also: - ../plugins/linux_commands_wildcard_injection
- ../plugins/subprocess_popen_with_shell_equals_true
- ../plugins/start_process_with_no_shell
- ../plugins/start_process_with_a_shell
- ../plugins/start_process_with_partial_path
 - **Config Options:**- This plugin test shares a configuration with others in the same family, namely shell_injection. This configuration is divided up into three sections, subprocess, shell and no_shell. They each list Python calls that spawn subprocesses, invoke commands within a shell, or invoke commands without a shell (by replacing the calling process) respectively. - This plugin specifically scans for methods listed in subprocess section that have shell=False specified. - shell_injection: # Start a process using the subprocess module, or one of its wrappers. subprocess: - subprocess.Popen - subprocess.call - Example:
 - >> Issue: subprocess call - check for execution of untrusted input. Severity: Low Confidence: High CWE: CWE-78 (https://cwe.mitre.org/data/definitions/78.html) Location: ./examples/subprocess_shell.py:23 22 23 subprocess.check_output(['/bin/ls', '-l']) 24 - See also - Added in version 0.9.0. - Changed in version 1.7.3: CWE information added

# B604: any_other_function_with_shell_equals_true

-
bandit.plugins.injection_shell.any_other_function_with_shell_equals_true(*context*,*config*)[source]
- **B604: Test for any function with shell equals true**- Python possesses many mechanisms to invoke an external executable. However, doing so may present a security issue if appropriate care is not taken to sanitize any user provided or variable input. - This plugin test is part of a family of tests built to check for process spawning and warn appropriately. Specifically, this plugin test interrogates method calls for the presence of a keyword parameter shell equalling true. It is related to detection of shell injection issues and is intended to catch custom wrappers to vulnerable methods that may have been created. - See also: - ../plugins/linux_commands_wildcard_injection
- ../plugins/subprocess_popen_with_shell_equals_true
- ../plugins/subprocess_without_shell_equals_true
- ../plugins/start_process_with_no_shell
- ../plugins/start_process_with_a_shell
- ../plugins/start_process_with_partial_path
 - **Config Options:**- This plugin test shares a configuration with others in the same family, namely shell_injection. This configuration is divided up into three sections, subprocess, shell and no_shell. They each list Python calls that spawn subprocesses, invoke commands within a shell, or invoke commands without a shell (by replacing the calling process) respectively. - Specifically, this plugin excludes those functions listed under the subprocess section, these methods are tested in a separate specific test plugin and this exclusion prevents duplicate issue reporting. - shell_injection: # Start a process using the subprocess module, or one of its wrappers. subprocess: [subprocess.Popen, subprocess.call, subprocess.check_call, subprocess.check_output execute_with_timeout] - Example:
 - `>> Issue: Function call with shell=True parameter identified, possible security issue. Severity: Medium Confidence: High CWE: CWE-78 (https://cwe.mitre.org/data/definitions/78.html) Location: ./examples/subprocess_shell.py:9 8 pop('/bin/gcc --version', shell=True) 9 Popen('/bin/gcc --version', shell=True) 10`- See also - Added in version 0.9.0. - Changed in version 1.7.3: CWE information added

# B605: start_process_with_a_shell

-
bandit.plugins.injection_shell.start_process_with_a_shell(*context*,*config*)[source]
- **B605: Test for starting a process with a shell**- Python possesses many mechanisms to invoke an external executable. However, doing so may present a security issue if appropriate care is not taken to sanitize any user provided or variable input. - This plugin test is part of a family of tests built to check for process spawning and warn appropriately. Specifically, this test looks for the spawning of a subprocess using a command shell. This type of subprocess invocation is dangerous as it is vulnerable to various shell injection attacks. Great care should be taken to sanitize all input in order to mitigate this risk. Calls of this type are identified by the use of certain commands which are known to use shells. Bandit will report a LOW severity warning. - See also: - ../plugins/linux_commands_wildcard_injection
- ../plugins/subprocess_without_shell_equals_true
- ../plugins/start_process_with_no_shell
- ../plugins/start_process_with_partial_path
- ../plugins/subprocess_popen_with_shell_equals_true
 - **Config Options:**- This plugin test shares a configuration with others in the same family, namely shell_injection. This configuration is divided up into three sections, subprocess, shell and no_shell. They each list Python calls that spawn subprocesses, invoke commands within a shell, or invoke commands without a shell (by replacing the calling process) respectively. - This plugin specifically scans for methods listed in shell section. - shell_injection: shell: - os.system - os.popen - os.popen2 - os.popen3 - os.popen4 - popen2.popen2 - popen2.popen3 - popen2.popen4 - popen2.Popen3 - popen2.Popen4 - commands.getoutput - commands.getstatusoutput - subprocess.getoutput - subprocess.getstatusoutput - Example:
 - `>> Issue: Starting a process with a shell: check for injection. Severity: Low Confidence: Medium CWE: CWE-78 (https://cwe.mitre.org/data/definitions/78.html) Location: examples/os_system.py:3 2 3 os.system('/bin/echo hi')`- See also - Added in version 0.10.0. - Changed in version 1.7.3: CWE information added

# B606: start_process_with_no_shell

-
bandit.plugins.injection_shell.start_process_with_no_shell(*context*,*config*)[source]
- **B606: Test for starting a process with no shell**- Python possesses many mechanisms to invoke an external executable. However, doing so may present a security issue if appropriate care is not taken to sanitize any user provided or variable input. - This plugin test is part of a family of tests built to check for process spawning and warn appropriately. Specifically, this test looks for the spawning of a subprocess in a way that doesn’t use a shell. Although this is generally safe, it maybe useful for penetration testing workflows to track where external system calls are used. As such a LOW severity message is generated. - See also: - ../plugins/linux_commands_wildcard_injection
- ../plugins/subprocess_without_shell_equals_true
- ../plugins/start_process_with_a_shell
- ../plugins/start_process_with_partial_path
- ../plugins/subprocess_popen_with_shell_equals_true
 - **Config Options:**- This plugin test shares a configuration with others in the same family, namely shell_injection. This configuration is divided up into three sections, subprocess, shell and no_shell. They each list Python calls that spawn subprocesses, invoke commands within a shell, or invoke commands without a shell (by replacing the calling process) respectively. - This plugin specifically scans for methods listed in no_shell section. - shell_injection: no_shell: - os.execl - os.execle - os.execlp - os.execlpe - os.execv - os.execve - os.execvp - os.execvpe - os.spawnl - os.spawnle - os.spawnlp - os.spawnlpe - os.spawnv - os.spawnve - os.spawnvp - os.spawnvpe - os.startfile - Example:
 - >> Issue: [start_process_with_no_shell] Starting a process without a shell. Severity: Low Confidence: Medium CWE: CWE-78 (https://cwe.mitre.org/data/definitions/78.html) Location: examples/os-spawn.py:8 7 os.spawnv(mode, path, args) 8 os.spawnve(mode, path, args, env) 9 os.spawnvp(mode, file, args) - See also - Added in version 0.10.0. - Changed in version 1.7.3: CWE information added

# B607: start_process_with_partial_path

-
bandit.plugins.injection_shell.start_process_with_partial_path(*context*,*config*)[source]
- **B607: Test for starting a process with a partial path**- Python possesses many mechanisms to invoke an external executable. If the desired executable path is not fully qualified relative to the filesystem root then this may present a potential security risk. - In POSIX environments, the PATH environment variable is used to specify a set of standard locations that will be searched for the first matching named executable. While convenient, this behavior may allow a malicious actor to exert control over a system. If they are able to adjust the contents of the PATH variable, or manipulate the file system, then a bogus executable may be discovered in place of the desired one. This executable will be invoked with the user privileges of the Python process that spawned it, potentially a highly privileged user. - This test will scan the parameters of all configured Python methods, looking for paths that do not start at the filesystem root, that is, do not have a leading ‘/’ character. - **Config Options:**- This plugin test shares a configuration with others in the same family, namely shell_injection. This configuration is divided up into three sections, subprocess, shell and no_shell. They each list Python calls that spawn subprocesses, invoke commands within a shell, or invoke commands without a shell (by replacing the calling process) respectively. - This test will scan parameters of all methods in all sections. Note that methods are fully qualified and de-aliased prior to checking. - shell_injection: # Start a process using the subprocess module, or one of its wrappers. subprocess: - subprocess.Popen - subprocess.call # Start a process with a function vulnerable to shell injection. shell: - os.system - os.popen - popen2.Popen3 - popen2.Popen4 - commands.getoutput - commands.getstatusoutput # Start a process with a function that is not vulnerable to shell injection. no_shell: - os.execl - os.execle - Example:
 - `>> Issue: Starting a process with a partial executable path Severity: Low Confidence: High CWE: CWE-78 (https://cwe.mitre.org/data/definitions/78.html) Location: ./examples/partial_path_process.py:3 2 from subprocess import Popen as pop 3 pop('gcc --version', shell=False)`- See also - Added in version 0.13.0. - Changed in version 1.7.3: CWE information added

# B608: hardcoded_sql_expressions

## B608: Test for SQL injection

An SQL injection attack consists of insertion or “injection” of a SQL query via the input data given to an application. It is a very common attack vector. This plugin test looks for strings that resemble SQL statements that are involved in some form of string building operation. For example:

“SELECT %s FROM derp;” % var

“SELECT thing FROM “ + tab

“SELECT “ + val + “ FROM “ + tab + …

“SELECT {} FROM derp;”.format(var)

f”SELECT foo FROM bar WHERE id = {product}”

Unless care is taken to sanitize and control the input data when building such SQL statement strings, an injection attack becomes possible. If strings of this nature are discovered, a LOW confidence issue is reported. In order to boost result confidence, this plugin test will also check to see if the discovered string is in use with standard Python DBAPI calls execute or executemany. If so, a MEDIUM issue is reported. For example:

cursor.execute(“SELECT %s FROM derp;” % var)

Use of str.replace in the string construction can also be dangerous. For example:

- “SELECT * FROM foo WHERE id = ‘[VALUE]’”.replace(“[VALUE]”, identifier)

However, such cases are always reported with LOW confidence to compensate for false positives, since valid uses of str.replace can be common.

- Example:

```
>> Issue: Possible SQL injection vector through string-based query
construction.
 Severity: Medium Confidence: Low
 CWE: CWE-89 (https://cwe.mitre.org/data/definitions/89.html)
 Location: ./examples/sql_statements.py:4
3 query = "DELETE FROM foo WHERE id = '%s'" % identifier
4 query = "UPDATE foo SET value = 'b' WHERE id = '%s'" % identifier
5
```
See also

Added in version 0.9.0.

Changed in version 1.7.3: CWE information added

Changed in version 1.7.7: Flag when str.replace is used in the string construction

# B609: linux_commands_wildcard_injection

## B609: Test for use of wildcard injection

Python provides a number of methods that emulate the behavior of standard Linux command line utilities. Like their Linux counterparts, these commands may take a wildcard “*” character in place of a file system path. This is interpreted to mean “any and all files or folders” and can be used to build partially qualified paths, such as “/home/user/*”.

The use of partially qualified paths may result in unintended consequences if an unexpected file or symlink is placed into the path location given. This becomes particularly dangerous when combined with commands used to manipulate file permissions or copy data off of a system.

This test plugin looks for usage of the following commands in conjunction with wild card parameters:

- ‘chown’
- ‘chmod’
- ‘tar’
- ‘rsync’

As well as any method configured in the shell or subprocess injection test configurations.

**Config Options:**

This plugin test shares a configuration with others in the same family, namely shell_injection. This configuration is divided up into three sections, subprocess, shell and no_shell. They each list Python calls that spawn subprocesses, invoke commands within a shell, or invoke commands without a shell (by replacing the calling process) respectively.

This test will scan parameters of all methods in all sections. Note that methods are fully qualified and de-aliased prior to checking.

```
shell_injection:
 # Start a process using the subprocess module, or one of its wrappers.
 subprocess:
 - subprocess.Popen
 - subprocess.call
 # Start a process with a function vulnerable to shell injection.
 shell:
 - os.system
 - os.popen
 - popen2.Popen3
 - popen2.Popen4
 - commands.getoutput
 - commands.getstatusoutput
 # Start a process with a function that is not vulnerable to shell
 injection.
 no_shell:
 - os.execl
 - os.execle
```
- Example:

```
>> Issue: Possible wildcard injection in call: subprocess.Popen
 Severity: High Confidence: Medium
 CWE-78 (https://cwe.mitre.org/data/definitions/78.html)
 Location: ./examples/wildcard-injection.py:8
7 o.popen2('/bin/chmod *')
8 subp.Popen('/bin/chown *', shell=True)
9
>> Issue: subprocess call - check for execution of untrusted input.
 Severity: Low Confidence: High
 CWE-78 (https://cwe.mitre.org/data/definitions/78.html)
 Location: ./examples/wildcard-injection.py:11
10 # Not vulnerable to wildcard injection
11 subp.Popen('/bin/rsync *')
12 subp.Popen("/bin/chmod *")
```
See also

Added in version 0.9.0.

Changed in version 1.7.3: CWE information added

# B610: django_extra_used

-
bandit.plugins.django_sql_injection.django_extra_used(*context*)[source]
- **B610: Potential SQL injection on extra function**- Example:
 - >> Issue: [B610:django_extra_used] Use of extra potential SQL attack vector. Severity: Medium Confidence: Medium CWE: CWE-89 (https://cwe.mitre.org/data/definitions/89.html) Location: examples/django_sql_injection_extra.py:29:0 More Info: https://bandit.readthedocs.io/en/latest/plugins/b610_django_extra_used.html 28 tables_str = 'django_content_type" WHERE "auth_user"."username"="admin' 29 User.objects.all().extra(tables=[tables_str]).distinct() - See also - Added in version 1.5.0. - Changed in version 1.7.3: CWE information added

# B611: django_rawsql_used

-
bandit.plugins.django_sql_injection.django_rawsql_used(*context*)[source]
- **B611: Potential SQL injection on RawSQL function**- Example:
 - >> Issue: [B611:django_rawsql_used] Use of RawSQL potential SQL attack vector. Severity: Medium Confidence: Medium CWE: CWE-89 (https://cwe.mitre.org/data/definitions/89.html) Location: examples/django_sql_injection_raw.py:11:26 More Info: https://bandit.readthedocs.io/en/latest/plugins/b611_django_rawsql_used.html 10 ' WHERE "username"="admin" OR 1=%s --' 11 User.objects.annotate(val=RawSQL(raw, [0])) - See also - Added in version 1.5.0. - Changed in version 1.7.3: CWE information added

# B612: logging_config_insecure_listen

## B612: Test for insecure use of logging.config.listen

This plugin test checks for the unsafe usage of the
`logging.config.listen` function. The logging.config.listen
function provides the ability to listen for external
configuration files on a socket server. Because portions of the
configuration are passed through eval(), use of this function
may open its users to a security risk. While the function only
binds to a socket on localhost, and so does not accept connections
from remote machines, there are scenarios where untrusted code
could be run under the account of the process which calls listen().

logging.config.listen provides the ability to verify bytes received across the socket with signature verification or encryption/decryption.

- Example:

```
>> Issue: [B612:logging_config_listen] Use of insecure
logging.config.listen detected.
 Severity: Medium Confidence: High
 CWE: CWE-94 (https://cwe.mitre.org/data/definitions/94.html)
 Location: examples/logging_config_insecure_listen.py:3:4
2
3 t = logging.config.listen(9999)
```
Added in version 1.7.5.

# B613: trojansource

## B613: TrojanSource - Bidirectional control characters

This plugin checks for the presence of unicode bidirectional control characters in Python source files. Those characters can be embedded in comments and strings to reorder source code characters in a way that changes its logic.

- Example:

```
>> Issue: [B613:trojansource] A Python source file contains bidirectional control characters ('\u202e').
 Severity: High Confidence: Medium
 CWE: CWE-838 (https://cwe.mitre.org/data/definitions/838.html)
 More Info: https://bandit.readthedocs.io/en/1.7.5/plugins/b113_trojansource.html
 Location: examples/trojansource.py:4:25
 3 access_level = "user"
 4 if access_level != 'none': # Check if admin ' and access_level != 'user
 5 print("You are an admin.\n")
```
Added in version 1.7.10.

# B614: pytorch_load

## B614: Test for unsafe PyTorch load

This plugin checks for unsafe use of torch.load and torch.serialization.load. Using torch.load or torch.serialization.load with untrusted data can lead to arbitrary code execution. There are two safe alternatives:

- Use torch.load with weights_only=True where only tensor data is extracted, and no arbitrary Python objects are deserialized
- Use the safetensors library from huggingface, which provides a safe deserialization mechanism

With weights_only=True, PyTorch enforces a strict type check, ensuring that only torch.Tensor objects are loaded.

- Example:

```
>> Issue: Use of unsafe PyTorch load
Severity: Medium Confidence: High
CWE: CWE-502 (https://cwe.mitre.org/data/definitions/502.html)
Location: examples/pytorch_load_save.py:8
7 loaded_model.load_state_dict(torch.load('model_weights.pth'))
8 another_model.load_state_dict(torch.load('model_weights.pth',
 map_location='cpu'))
9
10 print("Model loaded successfully!")
```
See also

Added in version 1.7.10.

# B615: huggingface_unsafe_download

## B615: Test for unsafe Hugging Face Hub downloads

This plugin checks for unsafe downloads from Hugging Face Hub without proper integrity verification. Downloading models, datasets, or files without specifying a revision based on an immmutable revision (commit) can lead to supply chain attacks where malicious actors could replace model files and use an existing tag or branch name to serve malicious content.

The secure approach is to:

- Pin to specific revisions/commits when downloading models, files or datasets

Common unsafe patterns:
- `AutoModel.from_pretrained("org/model-name")`
- `AutoModel.from_pretrained("org/model-name", revision="main")`
- `AutoModel.from_pretrained("org/model-name", revision="v1.0.0")`
- `load_dataset("org/dataset-name")` without revision
- `load_dataset("org/dataset-name", revision="main")`
- `load_dataset("org/dataset-name", revision="v1.0")`
- `AutoTokenizer.from_pretrained("org/model-name")`
- `AutoTokenizer.from_pretrained("org/model-name", revision="main")`
- `AutoTokenizer.from_pretrained("org/model-name", revision="v3.3.0")`
- `hf_hub_download(repo_id="org/model_name", filename="file_name")`
- ``hf_hub_download(repo_id=”org/model_name”,

filename=”file_name”, revision=”main” )``

- ``hf_hub_download(repo_id=”org/model_name”,
- filename=”file_name”, revision=”v2.0.0” - )``

- `snapshot_download(repo_id="org/model_name")`
- `snapshot_download(repo_id="org/model_name", revision="main")`
- `snapshot_download(repo_id="org/model_name", revision="refs/pr/1")`

- Example:

```
>> Issue: Unsafe Hugging Face Hub download without revision pinning
Severity: Medium Confidence: High
CWE: CWE-494 (https://cwe.mitre.org/data/definitions/494.html)
Location: examples/huggingface_unsafe_download.py:8
7 # Unsafe: no revision specified
8 model = AutoModel.from_pretrained("org/model_name")
9
```
See also

Added in version 1.8.6.

# B701: jinja2_autoescape_false

## B701: Test for not auto escaping in jinja2

Jinja2 is a Python HTML templating system. It is typically used to build web applications, though appears in other places well, notably the Ansible automation system. When configuring the Jinja2 environment, the option to use autoescaping on input can be specified. When autoescaping is enabled, Jinja2 will filter input strings to escape any HTML content submitted via template variables. Without escaping HTML input the application becomes vulnerable to Cross Site Scripting (XSS) attacks.

Unfortunately, autoescaping is False by default. Thus this plugin test will warn on omission of an autoescape setting, as well as an explicit setting of false. A HIGH severity warning is generated in either of these scenarios.

- Example:

```
>> Issue: Using jinja2 templates with autoescape=False is dangerous and can
lead to XSS. Use autoescape=True to mitigate XSS vulnerabilities.
 Severity: High Confidence: High
 CWE: CWE-94 (https://cwe.mitre.org/data/definitions/94.html)
 Location: ./examples/jinja2_templating.py:11
10 templateEnv = jinja2.Environment(autoescape=False,
 loader=templateLoader)
11 Environment(loader=templateLoader,
12 load=templateLoader,
13 autoescape=False)
14
>> Issue: By default, jinja2 sets autoescape to False. Consider using
autoescape=True or use the select_autoescape function to mitigate XSS
vulnerabilities.
 Severity: High Confidence: High
 CWE: CWE-94 (https://cwe.mitre.org/data/definitions/94.html)
 Location: ./examples/jinja2_templating.py:15
14
15 Environment(loader=templateLoader,
16 load=templateLoader)
17
18 Environment(autoescape=select_autoescape(['html', 'htm', 'xml']),
19 loader=templateLoader)
```
See also

Added in version 0.10.0.

Changed in version 1.7.3: CWE information added

# B702: use_of_mako_templates

## B702: Test for use of mako templates

Mako is a Python templating system often used to build web applications. It is the default templating system used in Pylons and Pyramid. Unlike Jinja2 (an alternative templating system), Mako has no environment wide variable escaping mechanism. Because of this, all input variables must be carefully escaped before use to prevent possible vulnerabilities to Cross Site Scripting (XSS) attacks.

- Example:

```
>> Issue: Mako templates allow HTML/JS rendering by default and are
inherently open to XSS attacks. Ensure variables in all templates are
properly sanitized via the 'n', 'h' or 'x' flags (depending on context).
For example, to HTML escape the variable 'data' do ${ data |h }.
 Severity: Medium Confidence: High
 CWE: CWE-80 (https://cwe.mitre.org/data/definitions/80.html)
 Location: ./examples/mako_templating.py:10
9
10 mako.template.Template("hern")
11 template.Template("hern")
```
See also

Added in version 0.10.0.

Changed in version 1.7.3: CWE information added

# B703: django_mark_safe

-
bandit.plugins.django_xss.django_mark_safe(*context*)[source]
- **B703: Potential XSS on mark_safe function**- Example:
 - >> Issue: [B703:django_mark_safe] Potential XSS on mark_safe function. Severity: Medium Confidence: High CWE: CWE-80 (https://cwe.mitre.org/data/definitions/80.html) Location: examples/mark_safe_insecure.py:159:4 More Info: https://bandit.readthedocs.io/en/latest/plugins/b703_django_mark_safe.html 158 str_arg = 'could be insecure' 159 safestring.mark_safe(str_arg) - See also - Added in version 1.5.0. - Changed in version 1.7.3: CWE information added

# B704: markupsafe_markup_xss

## B704: Potential XSS on markupsafe.Markup use

`markupsafe.Markup` does not perform any escaping, so passing dynamic
content, like f-strings, variables or interpolated strings will potentially
lead to XSS vulnerabilities, especially if that data was submitted by users.

Instead you should interpolate the resulting `markupsafe.Markup` object,
which will perform escaping, or use `markupsafe.escape`.

**Config Options:**

This plugin allows you to specify additional callable that should be treated
like `markupsafe.Markup`. By default we recognize `flask.Markup` as
an alias, but there are other subclasses or similar classes in the wild
that you may wish to treat the same.

Additionally there is a whitelist for callable names, whose result may
be safely passed into `markupsafe.Markup`. This is useful for escape
functions like e.g. `bleach.clean` which don’t themselves return
`markupsafe.Markup`, so they need to be wrapped. Take care when using
this setting, since incorrect use may introduce false negatives.

These two options can be set in a shared configuration section markupsafe_xss.

```
markupsafe_xss:
 # Recognize additional aliases
 extend_markup_names:
 - webhelpers.html.literal
 - my_package.Markup
 # Allow the output of these functions to pass into Markup
 allowed_calls:
 - bleach.clean
 - my_package.sanitize
```
- Example:

```
>> Issue: [B704:markupsafe_markup_xss] Potential XSS with
 ``markupsafe.Markup`` detected. Do not use ``Markup``
 on untrusted data.
 Severity: Medium Confidence: High
 CWE: CWE-79 (https://cwe.mitre.org/data/definitions/79.html)
 Location: ./examples/markupsafe_markup_xss.py:5:0
4 content = "<script>alert('Hello, world!')</script>"
5 Markup(f"unsafe {content}")
6 flask.Markup("unsafe {}".format(content))
```
See also

Added in version 1.8.3.
